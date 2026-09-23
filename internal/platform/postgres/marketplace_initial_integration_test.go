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

func TestMarketplaceInitialConnected(t *testing.T) {
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
	receiver := fixture.NewInitial(profile)
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
		_, e := call("maker", "POST", "/prepare", mb.PrepareRequest{ApprovalID: id, Generation: published.Generation, VariantID: profile.Bindings[0].VariantID, Operation: op, MediaApprovalID: mediaID, ExpiresAt: time.Now().Add(5 * time.Minute)}, &r)
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

	mediaRequest := prepare("MEDIA", "initial-media-approved", "")
	for _, token := range []string{"forbidden", "wrong-org", "wrong-tenant"} {
		if code, e := call(token, "POST", "/requests", mediaRequest, nil); e == nil || code != 403 {
			t.Fatal("scope", token, e)
		}
	}
	tampered := mediaRequest
	tampered.MediaSHA256 = strings.Repeat("f", 64)
	if _, e := call("maker", "POST", "/requests", tampered, nil); e == nil {
		t.Fatal("foreign source hash admitted")
	}
	mediaHash := submit(mediaRequest)
	if _, e := call("maker", "POST", "/requests/"+mediaRequest.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("unapproved upload")
	}
	approve(mediaRequest, mediaHash)
	// Concurrent clicks share the original fence; only one multipart upload.
	var wg sync.WaitGroup
	for range 12 {
		wg.Go(func() { call("maker", "POST", "/requests/"+mediaRequest.ApprovalID+"/send", struct{}{}, nil) })
	}
	wg.Wait()
	status(mediaRequest.ApprovalID, "accepted")
	if receiver.Uploads != 1 || receiver.UploadedSHA != mediaRequest.MediaSHA256 {
		t.Fatal("upload/fence/source", receiver.Uploads)
	}
	creation := prepare("CREATE", "initial-create-approved", mediaRequest.ApprovalID)
	if creation.Expected.Pictures[0].ID != fixture.PictureID || creation.SourceSHA256 != mediaRequest.SourceSHA256 {
		t.Fatal("accepted upload binding")
	}
	other := creation
	other.ApprovalID = "initial-create-invalid"
	other.MediaApprovalID = "missing-media-receipt"
	if _, e := call("maker", "POST", "/requests", other, nil); e == nil {
		t.Fatal("unaccepted media")
	}
	approve(creation, submit(creation))
	receiver.Mu.Lock()
	receiver.DropCreate = true
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+creation.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("lost provider acknowledgement")
	}
	status(creation.ApprovalID, "unknown")
	if receiver.Creates != 1 || receiver.Effects != 1 {
		t.Fatal("initial write count")
	}
	if _, e := call("maker", "POST", "/requests/"+creation.ApprovalID+"/send", struct{}{}, nil); e == nil || receiver.Creates != 1 {
		t.Fatal("unknown retried")
	}
	if _, e := call("maker", "POST", "/requests/"+creation.ApprovalID+"/reconcile", struct{}{}, nil); e != nil {
		t.Fatal("read-only reconcile", e)
	}
	status(creation.ApprovalID, "accepted")
	if _, e := call("maker", "POST", "/requests/"+creation.ApprovalID+"/send", struct{}{}, nil); e != nil || receiver.Creates != 1 {
		t.Fatal("accepted replay", e)
	}
	if receiver.Item.Quantity != 0 || receiver.Item.Status != "paused" || receiver.PriceMinor != 9007199254740991 {
		t.Fatal("ATP/price readback")
	}
	// No known provider image after a lost multipart acknowledgement: explicit
	// unknown, no forged SHA equivalence, no upload retry, new approval required.
	lost := prepare("MEDIA", "initial-media-lost-ack", "")
	approve(lost, submit(lost))
	receiver.Mu.Lock()
	receiver.DropUpload = true
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+lost.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("lost upload")
	}
	status(lost.ApprovalID, "unknown")
	if _, e := call("maker", "POST", "/requests/"+lost.ApprovalID+"/reconcile", struct{}{}, nil); e == nil {
		t.Fatal("invented image recovery")
	}
	if _, e := call("maker", "POST", "/requests/"+lost.ApprovalID+"/send", struct{}{}, nil); e == nil || receiver.Uploads != 2 {
		t.Fatal("upload replay")
	}
	// Another separately approved upload can replace an unassociated orphan,
	// without releasing or resetting the unknown operation.
	replacement := prepare("MEDIA", "initial-media-replacement", "")
	approve(replacement, submit(replacement))
	apiMu.Lock()
	dropAPI = true
	apiMu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+replacement.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("lost local API acknowledgement")
	}
	status(replacement.ApprovalID, "accepted")
	if _, e := call("maker", "POST", "/requests/"+replacement.ApprovalID+"/send", struct{}{}, nil); e != nil || receiver.Uploads != 3 {
		t.Fatal("API recovery replay")
	}
	// A second CREATE stays rejected even if provider search is temporarily empty.
	receiver.Mu.Lock()
	receiver.Exists = false
	receiver.Mu.Unlock()
	duplicate := prepare("CREATE", "initial-create-duplicate", replacement.ApprovalID)
	approve(duplicate, submit(duplicate))
	if _, e := call("maker", "POST", "/requests/"+duplicate.ApprovalID+"/send", struct{}{}, nil); e == nil || receiver.Creates != 1 {
		t.Fatal("historical local SKU exclusion")
	}
	for _, q := range []string{
		`update catalog.marketplace_observation set provider_id='fake' where tenant_id=$1`,
		`delete from catalog.marketplace_observation where tenant_id=$1`,
		`update catalog.marketplace_effect set operation='MEDIA' where tenant_id=$1`,
	} {
		if _, e := pool.Exec(ctx, q, tenant); e == nil {
			t.Fatal("immutable observation/effect")
		}
	}
	var approvals, effects, observations, attempts int
	for _, v := range []struct {
		sql string
		out *int
	}{
		{`select count(*) from approval.request where tenant_id=$1 and kind='marketplace_mutation'`, &approvals},
		{`select count(*) from catalog.marketplace_effect where tenant_id=$1`, &effects},
		{`select count(*) from catalog.marketplace_observation where tenant_id=$1`, &observations},
		{`select coalesce(sum(attempt_count),0) from communication.outbound_delivery where tenant_id=$1 and channel_code='mercadolibre_catalog'`, &attempts},
	} {
		if e := pool.QueryRow(ctx, v.sql, tenant).Scan(v.out); e != nil {
			t.Fatal(e)
		}
	}
	if approvals != 5 || effects != 4 || observations != 2 || attempts != 4 {
		t.Fatal("durable counts", approvals, effects, observations, attempts)
	}
	t.Logf("MARKETPLACE_INITIAL_PASS approvals=%d effects=%d observations=%d attempts=%d multipart=%d item_posts=%d api_posts=%d source=%s image=%s", approvals, effects, observations, attempts, receiver.Uploads, receiver.Creates, apiPosts, published.SnapshotSHA256, mediaRequest.MediaSHA256)
}
