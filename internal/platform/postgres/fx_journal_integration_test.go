package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/platform/httpapi"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestFXJournalConnectedPostingRecoveryAndReversal(t *testing.T) {
	rawURL := os.Getenv("FX_CONNECTED_DB_URL")
	if rawURL == "" {
		t.Skip("explicit local FX fixture required")
	}
	u, e := url.Parse(rawURL)
	if e != nil || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_fx_") {
		t.Fatal("owned loopback FX database required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	pool, e := pgxpool.New(ctx, rawURL)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	ids := randomid.Generator{}
	tenant := ids.New()
	if _, e = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'fx-journal','Fixture','Fixture')`, tenant); e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Store','store'),($1,'other','other','Other','store')`, tenant); e != nil {
		t.Fatal(e)
	}
	repo := db.NewAccounting(pool)
	ledger := accounting.NewService(repo, ids)
	fx, e := accounting.NewFXService(repo, ids, fxConnectedSnapshot(t, tenant, nil))
	if e != nil {
		t.Fatal(e)
	}
	for _, account := range []accounting.Account{{Code: "FX_ASSET", Name: "Fixture asset", Type: "asset"}, {Code: "FX_CLEARING", Name: "Fixture clearing", Type: "liability"}} {
		if _, e = ledger.CreateAccount(ctx, tenant, account); e != nil {
			t.Fatal(e)
		}
	}
	period, e := ledger.OpenPeriod(ctx, tenant, accounting.Period{StartsOn: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), EndsOn: time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)})
	if e != nil {
		t.Fatal(e)
	}
	verifier, token := handoverBrowserIssuer(t)
	mux := http.NewServeMux()
	httpapi.FXConversionModule{Service: fx}.Register(mux, verifier)
	httpapi.AccountingModule{Service: ledger}.Register(mux, verifier)
	server := httptest.NewServer(mux)
	defer server.Close()
	good := token("accountant", tenant, []string{"accounting:write", "accounting:read"}, []string{"store"})
	other := token("other", tenant, []string{"accounting:write", "accounting:read"}, []string{"store"})
	readonly := token("reader", tenant, []string{"accounting:read"}, []string{"store"})
	foreign := token("foreign", tenant, []string{"accounting:write", "accounting:read"}, []string{"other"})
	controller := token("controller", tenant, []string{"accounting:post", "accounting:reverse", "accounting:read"}, []string{"store"})
	request := func(method, path, key, bearer, body string) (int, http.Header, []byte) {
		r, e := http.NewRequestWithContext(ctx, method, server.URL+path, strings.NewReader(body))
		if e != nil {
			t.Error(e)
			return 0, nil, nil
		}
		r.Header.Set("Authorization", "Bearer "+bearer)
		r.Header.Set("Content-Type", "application/json")
		if key != "" {
			r.Header.Set("Idempotency-Key", key)
		}
		response, e := server.Client().Do(r)
		if e != nil {
			t.Error(e)
			return 0, nil, nil
		}
		defer response.Body.Close()
		raw, e := io.ReadAll(response.Body)
		if e != nil {
			t.Error(e)
		}
		return response.StatusCode, response.Header, raw
	}
	conversion := func(key, to string, amount int64) accounting.FXReceipt {
		body := fmt.Sprintf(`{"organization_id":"store","from_currency":"EUR","to_currency":%q,"amount_minor":"%d","conversion_date":"2026-01-01"}`, to, amount)
		status, _, raw := request("POST", "/v1/accounting/fx/conversions", key, good, body)
		var out accounting.FXReceipt
		if status != 201 || json.Unmarshal(raw, &out) != nil {
			t.Fatalf("conversion %d %s", status, raw)
		}
		return out
	}
	command := func(receipt accounting.FXReceipt) accounting.FXJournalCommand {
		return accounting.FXJournalCommand{OrganizationID: "store", ConversionID: receipt.ID, ConversionRequestKey: receipt.RequestKey, PeriodID: period.ID, DebitAccount: "FX_ASSET", CreditAccount: "FX_CLEARING", Description: "Synthetic fixture conversion"}
	}
	encode := func(v any) string {
		raw, e := json.Marshal(v)
		if e != nil {
			t.Fatal(e)
		}
		return string(raw)
	}
	converted := conversion("fx-journal-conversion-01", "USD", 1000)
	if converted.Conversion.OutputMinor != 1250 {
		t.Fatal("independent fixture amount", converted)
	}
	c := command(converted)
	body := encode(c)
	key := "fx-journal-prepare-0001"
	var group sync.WaitGroup
	results := make(chan accounting.FXJournalResult, 12)
	for range 12 {
		group.Add(1)
		go func() {
			defer group.Done()
			status, _, raw := request("POST", "/v1/accounting/fx/journals", key, good, body)
			var result accounting.FXJournalResult
			if status != 201 || json.Unmarshal(raw, &result) != nil {
				t.Errorf("prepare %d %s", status, raw)
				return
			}
			results <- result
		}()
	}
	group.Wait()
	close(results)
	var initial accounting.FXJournalResult
	count := 0
	for result := range results {
		count++
		if initial.Receipt.JournalID == "" {
			initial = result
		}
		if encode(result) != encode(initial) {
			t.Fatal("concurrent draft receipts differ")
		}
	}
	if count != 12 || initial.CurrentStatus != "draft" || initial.Receipt.AmountMinor != 1250 || initial.Receipt.Currency != "USD" || initial.Receipt.PostingDate != "2026-01-01" || initial.Receipt.Effect != "FX_JOURNAL_PREPARED" {
		t.Fatal("draft receipt", initial, count)
	}
	resultPath := "/v1/accounting/fx/journals/result?organization_id=store"
	status, headers, recovered := request("GET", resultPath, key, good, "")
	if status != 200 || headers.Get("Cache-Control") != "no-store" || strings.TrimSpace(string(recovered)) != encode(initial) {
		t.Fatal("lost prepare response recovery", status, string(recovered))
	}
	var n int
	if e = pool.QueryRow(ctx, `select count(*) from accounting.journal where tenant_id=$1`, tenant).Scan(&n); e != nil || n != 1 {
		t.Fatal("journal count", n, e)
	}
	if e = pool.QueryRow(ctx, `select count(*) from accounting.entry where tenant_id=$1`, tenant).Scan(&n); e != nil || n != 0 {
		t.Fatal("draft posted by preparation", n, e)
	}
	for _, tc := range []struct {
		method, path, key, bearer, body string
		want                            int
	}{
		{"POST", "/v1/accounting/fx/journals", key, good, strings.Replace(body, "FX_ASSET", "FX_CLEARING", 1), 409},
		{"POST", "/v1/accounting/fx/journals", key, other, body, 409},
		{"GET", resultPath, key, other, "", 404},
		{"POST", "/v1/accounting/fx/journals", "fx-journal-no-scope-01", readonly, body, 403},
		{"POST", "/v1/accounting/fx/journals", "fx-journal-foreign-01", foreign, body, 403},
		{"POST", "/v1/accounting/fx/journals", "fx-journal-reuse-0001", good, body, 409},
		{"POST", "/v1/accounting/fx/journals", "fx-journal-inject-01", good, strings.Replace(body, "{", `{"amount_minor":"1",`, 1), 400},
		{"POST", "/v1/accounting/fx/journals", "fx-journal-duplicate-01", good, strings.Replace(body, "{", `{"organization_id":"store",`, 1), 400},
	} {
		status, _, raw := request(tc.method, tc.path, tc.key, tc.bearer, tc.body)
		if status != tc.want {
			t.Errorf("negative wanted%d got%d %s", tc.want, status, raw)
		}
	}
	forged := accounting.Journal{OrganizationID: "store", PeriodID: period.ID, SourceType: "FX_CONVERSION", SourceID: "forged", Currency: "USD", PostingDate: period.StartsOn, Lines: []accounting.Line{{LineNo: 1, AccountCode: "FX_ASSET", DebitMinorUnits: 1250}, {LineNo: 2, AccountCode: "FX_CLEARING", CreditMinorUnits: 1250}}}
	status, _, raw := request("POST", "/v1/accounting/journals", "", good, encode(forged))
	if status != 400 {
		t.Fatal("generic source spoof", status, string(raw))
	}
	for i, tc := range []struct {
		to     string
		amount int64
	}{{"GBP", 1000}, {"USD", 0}, {"USD", -1000}} {
		receipt := conversion(fmt.Sprintf("fx-journal-negative-%02d", i), tc.to, tc.amount)
		status, _, raw := request("POST", "/v1/accounting/fx/journals", fmt.Sprintf("fx-journal-reject-%03d", i), good, encode(command(receipt)))
		if status != 409 {
			t.Errorf("invalid journal conversion %d %s", status, raw)
		}
	}
	unused := conversion("fx-journal-rollback-conv", "USD", 1500)
	badCommand := command(unused)
	badCommand.DebitAccount = "MISSING_ACCOUNT"
	status, _, raw = request("POST", "/v1/accounting/fx/journals", "fx-journal-invalid-account", good, encode(badCommand))
	if status != 409 {
		t.Fatal("invalid account", status, string(raw))
	}
	if e = pool.QueryRow(ctx, `select count(*) from accounting.journal where tenant_id=$1 and source_id=$2`, tenant, unused.ID).Scan(&n); e != nil || n != 0 {
		t.Fatal("orphan after account rejection", n, e)
	}
	// Block the actual shared outbox while the existing draft writer is running.
	blocker, e := pool.Begin(ctx)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = blocker.Exec(ctx, `lock table platform.outbox_event in access exclusive mode`); e != nil {
		t.Fatal(e)
	}
	blockedCtx, blockedCancel := context.WithTimeout(ctx, 150*time.Millisecond)
	bounded := command(unused)
	bounded.IdempotencyKey = "fx-journal-outbox-cancel"
	_, _, blockedErr := fx.PrepareJournal(blockedCtx, tenant, "accountant", bounded)
	blockedCancel()
	if e = blocker.Rollback(ctx); e != nil {
		t.Fatal(e)
	}
	if blockedErr == nil {
		t.Fatal("blocked write committed")
	}
	if e = pool.QueryRow(ctx, `select count(*) from accounting.journal where tenant_id=$1 and source_id=$2`, tenant, unused.ID).Scan(&n); e != nil || n != 0 {
		t.Fatal("orphan after outbox cancellation", n, e)
	}
	if e = pool.QueryRow(ctx, `select count(*) from platform.idempotency_record where tenant_id=$1 and scope='accounting-fx-journal' and idempotency_key=$2`, tenant, bounded.IdempotencyKey).Scan(&n); e != nil || n != 0 {
		t.Fatal("orphan idempotency", n, e)
	}
	// Post/reverse through the existing separately authorized accounting endpoints.
	postPath := "/v1/accounting/journals/" + initial.Receipt.JournalID + "/post"
	postBody := `{"organization_id":"store","expected_version":1}`
	status, _, raw = request("POST", postPath, "", good, postBody)
	if status != 403 {
		t.Fatal("write permission posted journal", status, string(raw))
	}
	status, _, raw = request("POST", postPath, "", controller, postBody)
	if status != 200 {
		t.Fatal("post", status, string(raw))
	}
	status, _, raw = request("GET", resultPath, key, good, "")
	var posted accounting.FXJournalResult
	if status != 200 || json.Unmarshal(raw, &posted) != nil || posted.CurrentStatus != "posted" || posted.CurrentVersion != 2 || encode(posted.Receipt) != encode(initial.Receipt) {
		t.Fatal("post response lost; current projection", status, string(raw))
	}
	status, _, raw = request("GET", "/v1/accounting/trial-balance?organization_id=store&period_id="+url.QueryEscape(period.ID), "", controller, "")
	var balance []accounting.Balance
	if status != 200 || json.Unmarshal(raw, &balance) != nil || len(balance) != 2 {
		t.Fatal("trial balance", status, string(raw))
	}
	for _, line := range balance {
		if line.NetMinorUnits != 1250 && line.NetMinorUnits != -1250 {
			t.Fatal("converted posting amount", balance)
		}
	}
	status, _, raw = request("POST", "/v1/accounting/journals/"+initial.Receipt.JournalID+"/reverse", "", controller, `{"organization_id":"store","expected_version":2,"reason":"Synthetic fixture reversal"}`)
	if status != 201 {
		t.Fatal("reverse", status, string(raw))
	}
	status, _, raw = request("GET", resultPath, key, good, "")
	var reversed accounting.FXJournalResult
	if status != 200 || json.Unmarshal(raw, &reversed) != nil || reversed.CurrentStatus != "reversed" || reversed.CurrentVersion != 3 || encode(reversed.Receipt) != encode(initial.Receipt) {
		t.Fatal("historical/current split", status, string(raw))
	}
	status, _, raw = request("GET", "/v1/accounting/trial-balance?organization_id=store&period_id="+url.QueryEscape(period.ID), "", controller, "")
	if status != 200 || json.Unmarshal(raw, &balance) != nil {
		t.Fatal("reversed balance")
	}
	for _, line := range balance {
		if line.NetMinorUnits != 0 {
			t.Fatal("reversal not neutral", balance)
		}
	}
	if _, e = pool.Exec(ctx, `delete from platform.idempotency_record where tenant_id=$1 and scope='accounting-fx-journal' and idempotency_key=$2`, tenant, key); e != nil {
		t.Fatal(e)
	}
	status, headers, raw = request("POST", "/v1/accounting/fx/journals", key, good, body)
	if status != 201 || headers.Get("Idempotent-Replay") != "true" || strings.TrimSpace(string(raw)) != encode(reversed) {
		t.Fatal("retained receipt replay", status, string(raw))
	}
	if _, e = pool.Exec(ctx, `update accounting.fx_journal_receipt set requested_by_subject='other' where tenant_id=$1`, tenant); e == nil {
		t.Fatal("mutable FX journal receipt")
	}
	// A closed period prevents new FX preparation through the same existing writer.
	if _, e = ledger.ClosePeriod(ctx, tenant, period.ID, 1); e != nil {
		t.Fatal(e)
	}
	status, _, raw = request("POST", "/v1/accounting/fx/journals", "fx-journal-closed-period", good, encode(command(unused)))
	if status != 409 {
		t.Fatal("closed period", status, string(raw))
	}
	if e = pool.QueryRow(ctx, `select count(*) from accounting.journal where tenant_id=$1`, tenant).Scan(&n); e != nil || n != 2 {
		t.Fatal("expected original and reversal only", n, e)
	}
	if e = pool.QueryRow(ctx, `select count(*) from accounting.entry where tenant_id=$1`, tenant).Scan(&n); e != nil || n != 4 {
		t.Fatal("expected four immutable entries", n, e)
	}
}
