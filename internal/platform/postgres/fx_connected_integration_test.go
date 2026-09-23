package postgres_test

// AUTHORED local fixture: actual HTTP/RS256/JWKS, adapted arithmetic and PG.
import (
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcfx"
	"elite.local/enterprise/internal/platform/httpapi"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/hex"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func fxConnectedSnapshot(t testing.TB, tenant string, alter func(*bcfx.ProfileDocument)) *bcfx.Snapshot {
	t.Helper()
	p, e := os.ReadFile("../../../config/fx/reference-profile.json")
	if e != nil {
		t.Fatal(e)
	}
	r, e := os.ReadFile("../../../config/fx/reference-rates.json")
	if e != nil {
		t.Fatal(e)
	}
	var d bcfx.ProfileDocument
	if json.Unmarshal(p, &d) != nil {
		t.Fatal("profile")
	}
	d.TenantID = tenant
	d.ValidFrom = time.Now().UTC().Add(-time.Hour).Truncate(time.Second).Format(time.RFC3339)
	d.ValidUntil = time.Now().UTC().Add(time.Hour).Truncate(time.Second).Format(time.RFC3339)
	if alter != nil {
		alter(&d)
	}
	p, e = json.Marshal(d)
	if e != nil {
		t.Fatal(e)
	}
	sum := sha256.Sum256(p)
	s, e := bcfx.LoadSnapshot(p, r, bcfx.Activation{Enabled: true, ProfileID: d.ID, Revision: d.Revision, ProfileSHA256: hex.EncodeToString(sum[:]), TenantID: tenant, OrganizationID: d.OrganizationID, LocalCurrency: d.LocalCurrency})
	if e != nil {
		t.Fatal(e)
	}
	return s
}
func TestFXConnectedHTTPReceipts(t *testing.T) {
	dbURL := os.Getenv("FX_CONNECTED_DB_URL")
	if dbURL == "" {
		t.Skip("explicit FX fixture DB required")
	}
	u, e := url.Parse(dbURL)
	if e != nil || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_fx_") {
		t.Fatal("isolated loopback FX DB required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, e := pgxpool.New(ctx, dbURL)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	ids := randomid.Generator{}
	tenant := ids.New()
	if _, e = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'fx-connected','Fixture','Fixture')`, tenant); e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Store','store'),($1,'other','other','Other','store')`, tenant); e != nil {
		t.Fatal(e)
	}
	snapshot := fxConnectedSnapshot(t, tenant, nil)
	repo := db.NewAccounting(pool)
	service, e := accounting.NewFXService(repo, ids, snapshot)
	if e != nil {
		t.Fatal(e)
	}
	verifier, token := handoverBrowserIssuer(t)
	mux := http.NewServeMux()
	httpapi.FXConversionModule{Service: service}.Register(mux, verifier)
	server := httptest.NewServer(mux)
	defer server.Close()
	good := token("accountant", tenant, []string{"accounting:read", "accounting:write"}, []string{"store"})
	other := token("other", tenant, []string{"accounting:read", "accounting:write"}, []string{"store"})
	readonly := token("reader", tenant, []string{"accounting:read"}, []string{"store"})
	foreign := token("foreign", tenant, []string{"accounting:read", "accounting:write"}, []string{"other"})
	request := func(method, key, bearer, body string) (int, http.Header, []byte) {
		path := "/v1/accounting/fx/conversions"
		if method == "GET" {
			path += "/result?organization_id=store"
		}
		r, e := http.NewRequestWithContext(ctx, method, server.URL+path, strings.NewReader(body))
		if e != nil {
			t.Fatal(e)
		}
		r.Header.Set("Authorization", "Bearer "+bearer)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Idempotency-Key", key)
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
	command := `{"organization_id":"store","from_currency":"EUR","to_currency":"GBP","amount_minor":"1000","conversion_date":"2026-01-01"}`
	key := "fx-connected-first-0001"
	status, _, first := request("POST", key, good, command)
	if status != 201 {
		t.Fatalf("POST %d %s", status, first)
	}
	// Treat the mutation response as lost; recover through the authenticated GET.
	status, headers, recovered := request("GET", key, good, "")
	if status != 200 || headers.Get("Cache-Control") != "no-store" || string(recovered) != string(first) {
		t.Fatal("recovery", status, string(recovered))
	}
	var receipt accounting.FXReceipt
	if json.Unmarshal(recovered, &receipt) != nil || receipt.Conversion.OutputMinor != 8333 || receipt.Conversion.FromRateID != "eur-2026-01-01" || receipt.Conversion.ToRateID != "gbp-2026-01-01" || receipt.Effect != "CONVERSION_RECEIPT_ONLY" {
		t.Fatal("receipt", receipt)
	}
	status, headers, replay := request("POST", key, good, command)
	if status != 201 || headers.Get("Idempotent-Replay") != "true" || string(replay) != string(first) {
		t.Fatal("replay")
	}
	for _, v := range []struct {
		method, key, actor, body string
		want                     int
	}{{"POST", key, good, strings.Replace(command, `"1000"`, `"1001"`, 1), 409}, {"POST", key, other, command, 409}, {"GET", key, other, "", 404}, {"POST", "fx-no-permission-01", readonly, command, 403}, {"POST", "fx-foreign-org-0001", foreign, command, 403}, {"POST", "fx-unknown-curr-01", good, strings.Replace(command, "EUR", "ZZZ", 1), 400}, {"POST", "fx-missing-money-01", good, strings.Replace(command, `"amount_minor":"1000",`, "", 1), 400}, {"POST", "fx-duplicate-json-01", good, strings.Replace(command, `"amount_minor":"1000"`, `"amount_minor":"1","amount_minor":"1000"`, 1), 400}} {
		status, _, raw := request(v.method, v.key, v.actor, v.body)
		if status != v.want {
			t.Fatal(v, status, string(raw))
		}
	}
	var wg sync.WaitGroup
	results := make(chan []byte, 12)
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			status, _, raw := request("POST", "fx-concurrent-key-0001", good, command)
			if status != 201 {
				t.Errorf("concurrent status=%d body=%s", status, raw)
			}
			results <- raw
		}()
	}
	wg.Wait()
	close(results)
	var parallel []byte
	for raw := range results {
		if parallel == nil {
			parallel = raw
		} else if string(parallel) != string(raw) {
			t.Fatal("concurrent receipts diverged")
		}
	}
	// Shared idempotency expiry does not permit a second durable conversion.
	if _, e = pool.Exec(ctx, `delete from platform.idempotency_record where tenant_id=$1 and scope='accounting-fx-conversion' and idempotency_key=$2`, tenant, key); e != nil {
		t.Fatal(e)
	}
	status, _, retained := request("POST", key, good, command)
	if status != 201 || string(retained) != string(first) {
		t.Fatal("retention replay")
	}
	if _, e = pool.Exec(ctx, `update accounting.fx_rate_snapshot set source_raw=source_raw where tenant_id=$1`, tenant); e == nil {
		t.Fatal("snapshot mutable")
	}
	if _, e = pool.Exec(ctx, `delete from accounting.fx_conversion_receipt where tenant_id=$1`, tenant); e == nil {
		t.Fatal("receipt mutable")
	}
	amount := int64(1000)
	c := accounting.FXCommand{OrganizationID: "store", FromCurrency: "EUR", ToCurrency: "USD", AmountMinor: &amount, ConversionDate: "2026-01-01", IdempotencyKey: "fx-policy-collision-01"}
	changed := fxConnectedSnapshot(t, tenant, func(p *bcfx.ProfileDocument) { p.Currencies[0].RoundingPrecisionMinor = "5" })
	changedService, _ := accounting.NewFXService(repo, ids, changed)
	if _, _, e = changedService.Record(ctx, tenant, "accountant", c); e == nil {
		t.Fatal("profile id/revision changed bytes")
	}
	// Expiration while waiting on an outbox write must roll back all new records.
	if _, e = pool.Exec(ctx, `create function public.fx_wait_fixture() returns trigger language plpgsql as $$ begin if NEW.event_type='fx-conversion.recorded' and NEW.payload->'snapshot'->>'profile_id'='expiry-wait' then perform pg_sleep(2.2); end if; return NEW; end; $$; create trigger fx_wait_fixture before insert on platform.outbox_event for each row execute function public.fx_wait_fixture();`); e != nil {
		t.Fatal(e)
	}
	expiry := fxConnectedSnapshot(t, tenant, func(p *bcfx.ProfileDocument) {
		p.ID = "expiry-wait"
		p.ValidUntil = time.Now().UTC().Add(2 * time.Second).Truncate(time.Second).Format(time.RFC3339)
	})
	expiryService, _ := accounting.NewFXService(repo, ids, expiry)
	c.IdempotencyKey = "fx-expiring-wait-0001"
	if _, _, e = expiryService.Record(ctx, tenant, "accountant", c); e == nil {
		t.Fatal("expired before commit")
	}
	var receipts, events, snapshots, journals, entries, claims int
	e = pool.QueryRow(ctx, `select (select count(*) from accounting.fx_conversion_receipt where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='fx-conversion.recorded'),(select count(*) from accounting.fx_rate_snapshot where tenant_id=$1),(select count(*) from accounting.journal where tenant_id=$1),(select count(*) from accounting.entry where tenant_id=$1),(select count(*) from platform.idempotency_record where tenant_id=$1 and idempotency_key in ('fx-expiring-wait-0001','fx-policy-collision-01'))`, tenant).Scan(&receipts, &events, &snapshots, &journals, &entries, &claims)
	if e != nil || receipts != 2 || events != 2 || snapshots != 1 || journals != 0 || entries != 0 || claims != 0 {
		t.Fatal(receipts, events, snapshots, journals, entries, claims, e)
	}
	t.Log("FX_CONNECTED_PASS HTTP_RS256_JWKS=true receipts=2 outbox=2 snapshots=1 concurrent12_one_receipt=true scoped_recovery=true immutable=true expires_while_outbox_wait_rolls_back=true retention_replay=true journal_entries=0")
}
