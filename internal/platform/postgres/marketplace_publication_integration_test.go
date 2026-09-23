package postgres_test

// AUTHORED local HTTP/PG/provider fixture. No live account or real inventory.
import (
	"bytes"
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/electromobility"
	mb "elite.local/enterprise/internal/marketplacebridge"
	fixture "elite.local/enterprise/internal/marketplacebridge/testfixture"
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
	"strings"
	"sync"
	"testing"
	"time"
)

type marketplaceFixtureToken struct{}

func (marketplaceFixtureToken) MercadoLibreAccessToken(context.Context) (string, error) {
	return fixture.Token, nil
}

type marketplaceFixtureVerifier map[string]identity.Principal
type marketplaceRecorder struct {
	*httptest.ResponseRecorder
	target http.ResponseWriter
}

func (w marketplaceRecorder) Unwrap() http.ResponseWriter { return w.target }

func (v marketplaceFixtureVerifier) Verify(_ context.Context, token string) (identity.Principal, error) {
	p, ok := v[token]
	if !ok {
		return p, identity.ErrUnauthenticated
	}
	return p, nil
}

func TestMarketplacePublicationConnected(t *testing.T) {
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
		return identity.Principal{TenantID: tenant, Subject: subject, Permissions: map[string]struct{}{"catalog:draft": {}, "catalog:read": {}, "catalog:publish": {}, "catalog:review:legal": {}, "catalog:review:technical": {}, "catalog:review:media": {}, "catalog:review:publication": {}, "marketplace:request": {}, "marketplace:approve": {}, "marketplace:read": {}, "marketplace:send": {}, "marketplace:reconcile": {}}, Organizations: map[string]struct{}{org: {}}}
	}
	maker, reviewer := principal("marketplace-maker"), principal("marketplace-reviewer")
	cp := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: tenant, OrganizationID: org, Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}
	catalog, e := db.NewCatalogRelease(pool, cp, db.NewCommerce(pool))
	if e != nil {
		t.Fatal(e)
	}
	source, e := catalog.CreateSource(ctx, maker, cr.SourceRequest{CommandID: "marketplace-source", Model: electromobility.Model{Code: "fixture-bicycle", DisplayName: "Fixture bicycle", VehicleClass: "bicycle", Specification: json.RawMessage("{}")}, Variants: []cr.SourceVariant{{Code: "fixture-bicycle-standard", DisplayName: "Standard", BatterySpecification: json.RawMessage("{}"), AmountMinorUnits: 9007199254740991, TaxMode: "not-applicable"}}, ValidFrom: time.Now().Add(-time.Hour)})
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
	profile := mb.Profile{Schema: "elite-marketplace-publication/v1", TenantID: tenant, OrganizationID: org, SellerID: "123456789", SiteID: "MLA", Currency: "ARS", CurrencyDigits: 2, MediaOrigin: cp.Origin, Bindings: []mb.Binding{{VariantID: source.SourceVariantIDs[0], SKU: "fixture-bicycle-standard", CategoryID: "MLA3530", ListingTypeID: "gold_special", StockMode: "seller_warehouse", StoreID: "123456", NetworkNodeID: "ARP12345", Attributes: []mb.Attribute{{ID: "ITEM_CONDITION", ValueID: "2230284"}}}}}
	receiver := fixture.New(profile)
	defer receiver.Close()
	service, e := db.NewMarketplacePublication(catalog, profile, bytes.Repeat([]byte("f"), 32), marketplaceFixtureToken{}, receiver.HTTP())
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
	httpapi.MarketplaceModule{Service: service, TenantID: tenant, OrganizationID: org}.Register(mux, verifier)
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
		request, _ := http.NewRequest(method, api.URL+"/v1/admin/marketplace"+path, bytes.NewReader(raw))
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
	prepare := func(op, id string) mb.Intent {
		t.Helper()
		var r mb.Intent
		_, e := call("maker", "POST", "/prepare", mb.PrepareRequest{ApprovalID: id, Generation: published.Generation, VariantID: profile.Bindings[0].VariantID, Operation: op, ItemID: receiver.Item.ID, ExpiresAt: time.Now().Add(5 * time.Minute)}, &r)
		if e != nil {
			t.Fatal("prepare", e)
		}
		if r.SourceSHA256 != published.SnapshotSHA256 || r.PriceMinorUnits != 9007199254740991 || r.Quantity != 0 {
			t.Fatal("source/ATP binding", r)
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
		var result db.MarketplaceStatus
		if _, e := call("maker", "GET", "/requests/"+id, nil, &result); e != nil || result.DeliveryState != want {
			t.Fatal("durable state", result.ApprovalID, result.ApprovalState, result.DeliveryState, result.FailureCode, "provider writes", len(receiver.Writes), e)
		}
	}
	price := prepare("PRICE", "marketplace-price-0001")
	for _, token := range []string{"forbidden", "wrong-org", "wrong-tenant"} {
		if code, e := call(token, "POST", "/requests", price, nil); e == nil || code != 403 {
			t.Fatal("scope boundary", token)
		}
	}
	tampered := price
	tampered.PriceMinorUnits = 1
	if _, e := call("maker", "POST", "/requests", tampered, nil); e == nil {
		t.Fatal("source amount replaced")
	}
	hash := submit(price)
	if _, e := call("maker", "POST", "/requests/"+price.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("unapproved effect")
	}
	approve(price, hash)
	var wg sync.WaitGroup
	successes := 0
	var mu sync.Mutex
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, e := call("maker", "POST", "/requests/"+price.ApprovalID+"/send", struct{}{}, nil)
			if e == nil {
				mu.Lock()
				successes++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	status(price.ApprovalID, "accepted")
	if successes == 0 || len(receiver.Writes) != 1 || receiver.PriceMinor != 9007199254740991 {
		t.Fatal("concurrent send", successes, len(receiver.Writes))
	}
	stock := prepare("STOCK", "marketplace-stock-0001")
	approve(stock, submit(stock))
	receiver.DropNext = true
	if _, e := call("maker", "POST", "/requests/"+stock.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("lost provider response accepted")
	}
	status(stock.ApprovalID, "unknown")
	if _, e := call("maker", "POST", "/requests/"+stock.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("unknown resent")
	}
	if _, e := call("maker", "POST", "/requests/"+stock.ApprovalID+"/reconcile", struct{}{}, nil); e != nil {
		t.Fatal("GET recovery", e)
	}
	status(stock.ApprovalID, "accepted")
	if len(receiver.Writes) != 2 || receiver.Quantity != 0 || receiver.Writes[1].Version != "7" {
		t.Fatal("stock recovery changed write/version")
	}
	pause := prepare("PAUSE", "marketplace-pause-0001")
	approve(pause, submit(pause))
	apiMu.Lock()
	dropAPI = true
	apiMu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+pause.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("API loss fixture")
	}
	status(pause.ApprovalID, "accepted")
	if _, e := call("maker", "POST", "/requests/"+pause.ApprovalID+"/send", struct{}{}, nil); e != nil {
		t.Fatal("accepted replay", e)
	}
	resume := prepare("RESUME", "marketplace-resume-0001")
	approve(resume, submit(resume))
	if _, e := call("maker", "POST", "/requests/"+resume.ApprovalID+"/send", struct{}{}, nil); e != nil {
		t.Fatal("resume", e)
	}
	rejection := prepare("PRICE", "marketplace-reject-0001")
	approve(rejection, submit(rejection))
	receiver.RejectNext = true
	if _, e := call("maker", "POST", "/requests/"+rejection.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("rejection accepted")
	}
	status(rejection.ApprovalID, "failed_terminal")
	if _, e := call("maker", "POST", "/requests/"+rejection.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("terminal retry")
	}
	drift := prepare("STOCK", "marketplace-drift-0001")
	approve(drift, submit(drift))
	receiver.Version++
	if _, e := call("maker", "POST", "/requests/"+drift.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("remote version drift accepted")
	}
	status(drift.ApprovalID, "failed_terminal")
	var approvals, effects, attempts int
	if e = pool.QueryRow(ctx, `select count(*) from approval.request where tenant_id=$1 and kind='marketplace_mutation'`, tenant).Scan(&approvals); e != nil {
		t.Fatal(e)
	}
	if e = pool.QueryRow(ctx, `select count(*),coalesce(sum(attempt_count),0) from communication.outbound_delivery where tenant_id=$1 and channel_code='mercadolibre_catalog'`, tenant).Scan(&effects, &attempts); e != nil {
		t.Fatal(e)
	}
	if approvals != 6 || effects != 6 || attempts != 6 || len(receiver.Writes) != 5 || receiver.Effects != 4 {
		t.Fatal("durable counts", approvals, effects, attempts, len(receiver.Writes), receiver.Effects)
	}
	for _, q := range []string{
		`update approval.request set payload=payload||'{"quantity":999}'::jsonb where tenant_id=$1 and kind='marketplace_mutation'`,
		`delete from approval.decision where tenant_id=$1 and request_id='marketplace-price-0001'`,
		`delete from catalog.marketplace_effect where tenant_id=$1`,
	} {
		if _, e = pool.Exec(ctx, q, tenant); e == nil {
			t.Fatal("evidence mutable")
		}
	}
	t.Logf("MARKETPLACE_CONNECTED_PASS source publication=%s; exact price=9007199254740991; existing ATP=0; 6 bound approvals/fences; 5 provider PUT attempts/4 effects; 12 concurrent sends -> 1 PUT; provider/API loss recovered; rejection/drift terminal; API POSTs=%d", published.SnapshotSHA256, apiPosts)
}
