package postgres_test

// AUTHORED HTTP/PG/provider fixture for initial catalog publication only.
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

func TestMarketplaceContentConnected(t *testing.T) {
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
	profile := mb.Profile{Schema: "elite-marketplace-publication/v1", TenantID: tenant, OrganizationID: org, SellerID: "123456789", SiteID: "MLA", Currency: "ARS", CurrencyDigits: 2, MediaOrigin: cp.Origin, Bindings: []mb.Binding{{VariantID: source.SourceVariantIDs[0], SKU: "fixture-bicycle-standard", CategoryID: "MLA3530", ListingTypeID: "gold_special", StockMode: "item", Attributes: []mb.Attribute{{ID: "ITEM_CONDITION", ValueID: "2230284"}}}}}
	receiver := fixture.NewContent(profile)
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
	prepare := func(op, id, mediaID string) mb.Intent {
		t.Helper()
		var r mb.Intent
		item := ""
		if op == "CONTENT" {
			item = receiver.Item.ID
		}
		_, e := call("maker", "POST", "/prepare", mb.PrepareRequest{ApprovalID: id, Generation: published.Generation, VariantID: profile.Bindings[0].VariantID, Operation: op, MediaApprovalID: mediaID, ItemID: item, ExpiresAt: time.Now().Add(5 * time.Minute)}, &r)
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

	mediaRequest := prepare("MEDIA", "content-approved-media", "")
	approve(mediaRequest, submit(mediaRequest))
	if _, e := call("maker", "POST", "/requests/"+mediaRequest.ApprovalID+"/send", struct{}{}, nil); e != nil {
		t.Fatal("media", e)
	}
	status(mediaRequest.ApprovalID, "accepted")
	rejected := prepare("CONTENT", "content-scope-reject", mediaRequest.ApprovalID)
	approve(rejected, submit(rejected))
	receiver.Mu.Lock()
	receiver.Multiple = true
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+rejected.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("multi-item update")
	}
	status(rejected.ApprovalID, "failed_terminal")
	if receiver.ContentWrites != 0 {
		t.Fatal("unapproved family affected")
	}
	receiver.Mu.Lock()
	receiver.Multiple = false
	receiver.Mu.Unlock()
	content := prepare("CONTENT", "content-approved-loss", mediaRequest.ApprovalID)
	if content.Expected.FamilyName != "Fixture bicycle" || content.Expected.Pictures[0].ID != fixture.PictureID {
		t.Fatal("source mapping")
	}
	tampered := content
	tampered.Expected.FamilyName = "Client invented title"
	if _, e := call("maker", "POST", "/requests", tampered, nil); e == nil {
		t.Fatal("source tamper")
	}
	for _, token := range []string{"forbidden", "wrong-org", "wrong-tenant"} {
		if code, e := call(token, "POST", "/requests", content, nil); e == nil || code != 403 {
			t.Fatal("scope", token, e)
		}
	}
	approve(content, submit(content))
	receiver.Mu.Lock()
	receiver.DropNext = true
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+content.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("missing write ack")
	}
	status(content.ApprovalID, "unknown")
	if _, e := call("maker", "POST", "/requests/"+content.ApprovalID+"/send", struct{}{}, nil); e == nil || receiver.ContentWrites != 1 {
		t.Fatal("unknown retried")
	}
	if _, e := call("maker", "POST", "/requests/"+content.ApprovalID+"/reconcile", struct{}{}, nil); e != nil {
		t.Fatal("GET reconcile", e)
	}
	status(content.ApprovalID, "accepted")
	var wg sync.WaitGroup
	for range 12 {
		wg.Go(func() { call("maker", "POST", "/requests/"+content.ApprovalID+"/send", struct{}{}, nil) })
	}
	wg.Wait()
	if receiver.ContentWrites != 1 || receiver.Item.FamilyName != "Fixture bicycle" || receiver.Item.Pictures[0].ID != fixture.PictureID || receiver.PriceMinor != 9007199254740993 || receiver.Quantity != 4 {
		t.Fatal("content changed price/stock or duplicate", receiver.ContentWrites)
	}
	for _, state := range []string{"sold", "missing", "multiple"} {
		receiver.Mu.Lock()
		receiver.Sales = 0
		receiver.MissingSales = false
		receiver.Multiple = false
		if state == "sold" {
			receiver.Sales = 1
		}
		if state == "missing" {
			receiver.MissingSales = true
		}
		if state == "multiple" {
			receiver.Multiple = true
		}
		receiver.Mu.Unlock()
		input := mb.PrepareRequest{ApprovalID: "content-denied-" + state, MediaApprovalID: mediaRequest.ApprovalID, Generation: published.Generation, VariantID: profile.Bindings[0].VariantID, Operation: "CONTENT", ItemID: receiver.Item.ID, ExpiresAt: time.Now().Add(time.Minute)}
		if _, e := call("maker", "POST", "/prepare", input, nil); e == nil {
			t.Fatal("unsafe preparation", state)
		}
	}
	var approvals, effects, attempts int
	for _, v := range []struct {
		sql string
		out *int
	}{
		{`select count(*) from approval.request where tenant_id=$1 and kind='marketplace_mutation'`, &approvals},
		{`select count(*) from catalog.marketplace_effect where tenant_id=$1`, &effects},
		{`select coalesce(sum(attempt_count),0) from communication.outbound_delivery where tenant_id=$1 and channel_code='mercadolibre_catalog'`, &attempts},
	} {
		if e := pool.QueryRow(ctx, v.sql, tenant).Scan(v.out); e != nil {
			t.Fatal(e)
		}
	}
	if approvals != 3 || effects != 3 || attempts != 3 || receiver.Uploads != 1 || receiver.ContentWrites != 1 {
		t.Fatal("durable counts", approvals, effects, attempts)
	}
	t.Logf("MARKETPLACE_CONTENT_PASS approvals=%d effects=%d attempts=%d multipart=%d content_puts=%d api_posts=%d source=%s", approvals, effects, attempts, receiver.Uploads, receiver.ContentWrites, apiPosts, published.SnapshotSHA256)
}
