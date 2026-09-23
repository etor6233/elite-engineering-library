package postgres_test

// AUTHORED HTTP/PG/provider fixture for initial catalog publication only.
import (
	"bytes"
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/electromobility"
	mb "elite.local/enterprise/internal/merchantbridge"
	fixture "elite.local/enterprise/internal/merchantbridge/testfixture"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"image"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMerchantCatalogConnected(t *testing.T) {
	database := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if database == "" {
		t.Skip("owned fixture DB required")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, database)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant := "d84c64f9-1252-4db5-bf5a-3af8026190f1"
	org := "marketplace-store"
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'marketplace-fixture','Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'marketplace-store','marketplace-store','Fixture','store')`,
	} {
		if _, e = pool.Exec(ctx, q, tenant); e != nil {
			t.Fatal(e)
		}
	}
	principal := func(subject string) identity.Principal {
		return identity.Principal{TenantID: tenant, Subject: subject, Permissions: map[string]struct{}{"catalog:draft": {}, "catalog:read": {}, "catalog:publish": {}, "catalog:review:legal": {}, "catalog:review:technical": {}, "catalog:review:media": {}, "catalog:review:publication": {}, "merchant:request": {}, "merchant:approve": {}, "merchant:read": {}, "merchant:send": {}, "merchant:reconcile": {}}, Organizations: map[string]struct{}{org: {}}}
	}
	maker, reviewer := principal("marketplace-maker"), principal("marketplace-reviewer")
	cp := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: tenant, OrganizationID: org, Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}
	catalog, e := db.NewCatalogRelease(pool, cp, db.NewCommerce(pool))
	if e != nil {
		t.Fatal(e)
	}
	source, e := catalog.CreateSource(ctx, maker, cr.SourceRequest{CommandID: "marketplace-source", Model: electromobility.Model{Code: "fixture-bicycle", DisplayName: "Fixture bicycle", VehicleClass: "bicycle", Specification: json.RawMessage(`{"merchant_description":"Descripción sintética revisada del catálogo."}`)}, Variants: []cr.SourceVariant{{Code: "fixture-bicycle-standard", DisplayName: "Standard", BatterySpecification: json.RawMessage("{}"), AmountMinorUnits: 1234567890123, TaxMode: "not-applicable"}}, ValidFrom: time.Now().Add(-time.Hour)})
	if e != nil {
		t.Fatal("source", e)
	}
	var pngRaw bytes.Buffer
	if e = png.Encode(&pngRaw, image.NewNRGBA(image.Rect(0, 0, 2, 2))); e != nil {
		t.Fatal(e)
	}
	media, e := catalog.UploadPNG(ctx, maker, "marketplace-media", pngRaw.Bytes())
	if e != nil {
		t.Fatal(e)
	}
	draft, e := catalog.CreateDraft(ctx, maker, cr.DraftRequest{CommandID: "marketplace-draft", PriceBookID: source.SourcePriceBookID, Models: []cr.ModelReference{{ModelID: source.ResourceID, MediaID: media.ResourceID}}})
	if e != nil {
		t.Fatal("draft", e)
	}
	for _, stage := range []string{"legal", "technical", "media", "publication"} {
		_, e = catalog.Review(ctx, reviewer, cr.ReviewRequest{CommandID: "marketplace-review-" + stage, DraftID: draft.ResourceID, Stage: stage, SnapshotSHA256: draft.SnapshotSHA256, Approved: true, Reason: "synthetic source review", EvidenceSHA256: strings.Repeat("a", 64)})
		if e != nil {
			t.Fatal("review", e)
		}
	}
	published, e := catalog.Publish(ctx, reviewer, cr.PublishRequest{CommandID: "marketplace-publish", DraftID: draft.ResourceID, SnapshotSHA256: draft.SnapshotSHA256, ExpectedGeneration: 0, Reason: "fixture publication"})
	if e != nil {
		t.Fatal("publish", e)
	}

	profile := mb.Profile{Schema: "elite-connected-merchant/v1", TenantID: tenant, OrganizationID: org, AccountID: "123456", DataSourceID: "789012", Language: "es", FeedLabel: "AR", Currency: "ARS", CurrencyDigits: 2, Origin: cp.Origin, RefreshCadenceDays: 30, Bindings: []mb.Binding{{VariantID: source.SourceVariantIDs[0], OfferID: "fixture-bicycle-standard", Condition: "NEW", GTINs: []string{}}}}
	receiver := fixture.New(profile)
	defer receiver.Close()
	root := os.Getenv("MERCHANT_TEST_ROOT")
	python := os.Getenv("MERCHANT_TEST_PYTHON")
	fileHash := func(path string) string {
		t.Helper()
		raw, e := os.ReadFile(path)
		if e != nil {
			t.Fatal(e)
		}
		return mb.Hash(raw)
	}
	script := filepath.Join(root, "google_merchant_product_sync", "connected_worker.py")
	process := mb.Process{Python: python, PythonSHA256: fileHash(python), Script: script, ScriptSHA256: fileHash(script), OwnerSHA256: fileHash(filepath.Join(filepath.Dir(script), "sync_product.py")), RuntimeLockSHA256: fileHash(filepath.Join(filepath.Dir(script), "connected-runtime.lock.json")), Mode: "LOCAL_FIXTURES", FixtureOrigin: receiver.HTTP.URL}
	service, e := db.NewMerchantPublication(catalog, profile, bytes.Repeat([]byte("f"), 32), process)
	if e != nil {
		t.Fatal(e)
	}

	forbidden := principal("unprivileged")
	forbidden.Permissions = map[string]struct{}{}
	wrongOrg := principal("wrong-org")
	wrongOrg.Organizations = map[string]struct{}{"other": {}}
	wrongTenant := principal("wrong-tenant")
	wrongTenant.TenantID = "33b9e628-aa94-4c70-9c24-6ac4a3fb805e"
	verifier := marketplaceFixtureVerifier{"maker": maker, "reviewer": reviewer, "forbidden": forbidden, "wrong-org": wrongOrg, "wrong-tenant": wrongTenant}
	mux := http.NewServeMux()
	httpapi.MerchantModule{Service: service, TenantID: tenant, OrganizationID: org}.Register(mux, verifier)
	var apiMu sync.Mutex
	dropAPI := false
	apiPosts := 0
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(marketplaceRecorder{recorder, w}, r)
		apiMu.Lock()
		defer apiMu.Unlock()
		if r.Method == "POST" {
			apiPosts++
		}
		if dropAPI && strings.HasSuffix(r.URL.Path, "/send") && recorder.Code == 200 {
			dropAPI = false
			conn, _, e := w.(http.Hijacker).Hijack()
			if e == nil {
				conn.Close()
			}
			return
		}
		for k, values := range recorder.Header() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(recorder.Code)
		_, _ = w.Write(recorder.Body.Bytes())
	}))
	defer api.Close()
	call := func(token, method, path string, in, out any) (int, error) {
		var raw []byte
		if in != nil {
			raw, _ = json.Marshal(in)
		}
		request, _ := http.NewRequest(method, api.URL+"/v1/admin/merchant"+path, bytes.NewReader(raw))
		request.Header.Set("Authorization", "Bearer "+token)
		request.Header.Set("Content-Type", "application/json")
		response, e := api.Client().Do(request)
		if e != nil {
			return 0, e
		}
		defer response.Body.Close()
		data, e := io.ReadAll(io.LimitReader(response.Body, 65537))
		if e != nil {
			return response.StatusCode, e
		}
		if response.Header.Get("Cache-Control") != "no-store" {
			return response.StatusCode, errors.New("private response cacheable")
		}
		if response.StatusCode >= 400 {
			return response.StatusCode, fmt.Errorf("HTTP %d %s", response.StatusCode, data)
		}
		if out != nil {
			e = json.Unmarshal(data, out)
		}
		return response.StatusCode, e
	}

	prepare := func(id string) mb.Intent {
		t.Helper()
		var r mb.Intent
		_, e := call("maker", "POST", "/prepare", mb.PrepareRequest{ApprovalID: id, Generation: published.Generation, VariantID: profile.Bindings[0].VariantID, ExpiresAt: time.Now().Add(5 * time.Minute)}, &r)
		if e != nil {
			t.Fatal("prepare", e)
		}
		if r.SourceSHA256 != published.SnapshotSHA256 || r.Product.Price.AmountMicros != "12345678901230000" || r.Quantity != 0 || r.Product.Availability != "OUT_OF_STOCK" {
			t.Fatal("source/amount/ATP", r)
		}
		return r
	}

	submit := func(r mb.Intent) string {
		t.Helper()
		var result struct {
			Hash string `json:"request_sha256"`
		}
		if _, e := call("maker", "POST", "/requests", r, &result); e != nil {
			t.Fatal("submit", e)
		}
		return result.Hash
	}
	approve := func(r mb.Intent, hash string) {
		t.Helper()
		body := map[string]any{"request_sha256": hash, "approve": true, "reason": "exact synthetic provider mutation"}
		if status, e := call("maker", "POST", "/requests/"+r.ApprovalID+"/decision", body, nil); e == nil || status != 409 {
			t.Fatal("self approval")
		}
		if _, e := call("reviewer", "POST", "/requests/"+r.ApprovalID+"/decision", body, nil); e != nil {
			t.Fatal("review", e)
		}
	}
	status := func(id, want string) {
		t.Helper()
		var result db.MerchantStatus
		if _, e := call("maker", "GET", "/requests/"+id, nil, &result); e != nil || result.DeliveryState != want {
			t.Fatal("durable state", result.ApprovalID, result.ApprovalState, result.DeliveryState, result.FailureCode, "provider writes", receiver.Inserts, e)
		}
	}

	first := prepare("merchant-lost-input")
	for _, token := range []string{"forbidden", "wrong-org", "wrong-tenant"} {
		if code, e := call(token, "POST", "/requests", first, nil); e == nil || code != 403 {
			t.Fatal("scope", token, e)
		}
	}
	tampered := first
	tampered.Product.Price.AmountMicros = "1"
	if _, e := call("maker", "POST", "/requests", tampered, nil); e == nil {
		t.Fatal("client price override")
	}
	approve(first, submit(first))
	receiver.Mu.Lock()
	receiver.DropNext = true
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+first.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("lost SDK acknowledgement")
	}
	status(first.ApprovalID, "unknown")
	if receiver.Inserts != 1 {
		t.Fatal("SDK retries", receiver.Inserts)
	}
	receiver.Mu.Lock()
	receiver.Pending = true
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+first.ApprovalID+"/reconcile", struct{}{}, nil); e == nil {
		t.Fatal("pending inferred success")
	}
	receiver.Mu.Lock()
	receiver.Pending = false
	receiver.ForeignSource = true
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+first.ApprovalID+"/reconcile", struct{}{}, nil); e == nil {
		t.Fatal("foreign data source accepted")
	}
	receiver.Mu.Lock()
	receiver.ForeignSource = false
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+first.ApprovalID+"/reconcile", struct{}{}, nil); e != nil {
		t.Fatal("GET source recovery", e)
	}
	status(first.ApprovalID, "accepted")
	var observed db.MerchantStatus
	if _, e := call("maker", "GET", "/requests/"+first.ApprovalID, nil, &observed); e != nil {
		t.Fatal(e)
	}
	if observed.InputAcknowledged || observed.RefreshDueAt != nil || !observed.ProcessingConfirmed || !bytes.Contains(observed.ProcessingResponse, []byte("pending_countries")) {
		t.Fatal("GET invented input refresh or Google approval", observed)
	}
	// A separately approved refresh retains the official source generation.
	refresh := prepare("merchant-approved-refresh")
	approve(refresh, submit(refresh))
	apiMu.Lock()
	dropAPI = true
	apiMu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+refresh.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("lost local API acknowledgement")
	}
	status(refresh.ApprovalID, "accepted")
	var wg sync.WaitGroup
	for range 12 {
		wg.Go(func() { call("maker", "POST", "/requests/"+refresh.ApprovalID+"/send", struct{}{}, nil) })
	}
	wg.Wait()
	if receiver.Inserts != 2 {
		t.Fatal("refresh replay", receiver.Inserts)
	}
	if _, e := call("maker", "POST", "/requests/"+refresh.ApprovalID+"/reconcile", struct{}{}, nil); e != nil {
		t.Fatal("processing status", e)
	}
	if _, e := call("maker", "GET", "/requests/"+refresh.ApprovalID, nil, &observed); e != nil || !observed.InputAcknowledged || observed.RefreshDueAt == nil || !observed.ProcessingConfirmed {
		t.Fatal("input/status separation", e)
	}
	rejected := prepare("merchant-provider-rejected")
	approve(rejected, submit(rejected))
	receiver.Mu.Lock()
	receiver.RejectNext = true
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+rejected.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("provider rejection")
	}
	status(rejected.ApprovalID, "failed_terminal")
	if _, e := call("maker", "POST", "/requests/"+rejected.ApprovalID+"/send", struct{}{}, nil); e == nil || receiver.Inserts != 3 {
		t.Fatal("rejected retry")
	}
	for _, q := range []string{
		`update catalog.merchant_observation set response='{}'::bytea where tenant_id=$1`,
		`delete from catalog.merchant_observation where tenant_id=$1`,
		`update catalog.merchant_effect set offer_id='forged' where tenant_id=$1`,
	} {
		if _, e := pool.Exec(ctx, q, tenant); e == nil {
			t.Fatal("immutable provider evidence")
		}
	}
	var approvals, effects, observations, attempts int
	for _, v := range []struct {
		sql string
		out *int
	}{
		{`select count(*) from approval.request where tenant_id=$1 and kind='marketplace_mutation'`, &approvals},
		{`select count(*) from catalog.merchant_effect where tenant_id=$1`, &effects},
		{`select count(*) from catalog.merchant_observation where tenant_id=$1`, &observations},
		{`select sum(attempt_count) from communication.outbound_delivery where tenant_id=$1 and channel_code='google_merchant'`, &attempts},
	} {
		if e := pool.QueryRow(ctx, v.sql, tenant).Scan(v.out); e != nil {
			t.Fatal(e)
		}
	}
	if approvals != 3 || effects != 3 || observations != 5 || attempts != 3 {
		t.Fatal("durable counts", approvals, effects, observations, attempts)
	}
	t.Logf("MERCHANT_CONNECTED_PASS approvals=%d effects=%d observations=%d attempts=%d actual_sdk_inserts=%d actual_sdk_gets=%d api_posts=%d source=%s", approvals, effects, observations, attempts, receiver.Inserts, receiver.Gets, apiPosts, published.SnapshotSHA256)
}
