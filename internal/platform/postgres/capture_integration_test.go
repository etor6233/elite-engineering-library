package postgres_test

import (
	"bytes"
	"context"
	cap "elite.local/enterprise/internal/capture"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/qr"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCaptureNativeHTTPAuthorizedResolution(t *testing.T) {
	dsn := os.Getenv("CAPTURE_TEST_DB_URL")
	if dsn == "" {
		t.Skip("owned capture database required")
	}
	u, e := url.Parse(dsn)
	if e != nil || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_document_reference_") {
		t.Fatal("owned loopback fixture required")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, dsn)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant := uuid.NewString()
	foreign := uuid.NewString()
	for _, ten := range []string{tenant, foreign} {
		_, e = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,$2,'Fixture','Fixture')`, ten, "capture-"+ten)
		if e != nil {
			t.Fatal(e)
		}
		for _, q := range []string{
			`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'org-a','a','Fixture A','store'),($1,'org-b','b','Fixture B','store')`,
			`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Fixture model','bicycle','active')`,
			`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant-a','model','a','Fixture vehicle','{}','active'),($1,'variant-private','model','private','Private fixture','{}','active')`,
			`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock-a','org-a','variant-a','SKU-1','available',1,clock_timestamp()),($1,'stock-private','org-b','variant-private','PRIVATE-SERIAL','available',1,clock_timestamp())`,
			`insert into partner.supplier(tenant_id,supplier_id,supplier_code,legal_name,status)values($1,'supplier','supplier','Fixture','active')`,
			`insert into procurement.purchase_order(tenant_id,purchase_order_id,supplier_id,destination_organization_id,state,currency,total_minor_units,version)values($1,'po','supplier','org-a','accepted','USD',100,1)`,
			`insert into factory.production_unit(tenant_id,production_unit_id,purchase_order_id,variant_id,serial_number,state)values($1,'factory-unit','po','variant-a','FACTORY-SERIAL','planned')`,
		} {
			if _, e = pool.Exec(ctx, q, ten); e != nil {
				t.Fatal(e)
			}
		}
	}
	principal := func(subject string, perms ...string) identity.Principal {
		p := identity.Principal{TenantID: tenant, Subject: subject, Organizations: map[string]struct{}{"org-a": {}}, Permissions: map[string]struct{}{}}
		for _, s := range perms {
			p.Permissions[s] = struct{}{}
		}
		return p
	}
	maker := principal("operator", "inventory:read")
	none := principal("visitor")
	factory := principal("factory", "factory:read")
	supply := principal("supply", "supply:read")
	wrong := maker
	wrong.Organizations = map[string]struct{}{"elsewhere": {}}
	b, e := os.ReadFile(os.Getenv("CAPTURE_TEST_PROFILE"))
	if e != nil {
		t.Fatal(e)
	}
	var config cap.Config
	if json.Unmarshal(b, &config) != nil {
		t.Fatal("profile")
	}
	decoder, e := cap.NewDecoder(config)
	if e != nil {
		t.Fatal(e)
	}
	store := db.NewCaptureStore(pool)
	mux := http.NewServeMux()
	httpapi.CaptureModule{Decoder: decoder, Store: store}.Register(mux, documentVerifier{"operator": maker, "none": none, "factory": factory, "supply": supply, "wrong": wrong})
	call := func(path, token, content string, b []byte) *httptest.ResponseRecorder {
		r := httptest.NewRequest("POST", "/v1/capture/"+path, bytes.NewReader(b))
		r.Header.Set("Authorization", "Bearer "+token)
		r.Header.Set("Content-Type", content)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, r)
		return rec
	}
	resolve := func(token, text, encoding string) *httptest.ResponseRecorder {
		b, _ := json.Marshal(map[string]string{"text": text, "encoding": encoding, "operation": "lookup"})
		return call("resolve", token, "application/json", b)
	}
	frame, e := os.ReadFile(filepath.Join(os.Getenv("CAPTURE_TEST_FRAMES"), "serial-code128.png.gray8"))
	if e != nil {
		t.Fatal(e)
	}
	rec := call("decode", "operator", "application/octet-stream", frame)
	var decoded cap.Result
	if rec.Code != 200 || json.Unmarshal(rec.Body.Bytes(), &decoded) != nil || decoded.Status != "DECODED" || decoded.Authorized || len(decoded.Candidates) != 1 || decoded.Candidates[0].Text != "SKU-1" {
		t.Fatal("native decode", rec.Code, rec.Body.String())
	}
	rec = resolve("operator", decoded.Candidates[0].Text, "SERIAL")
	var resolved db.CaptureResolution
	if rec.Code != 200 || json.Unmarshal(rec.Body.Bytes(), &resolved) != nil || !resolved.Authorized || resolved.BusinessEffect || len(resolved.Records) != 1 || resolved.Records[0].ID != "stock-a" || resolved.Records[0].OrganizationID != "org-a" {
		t.Fatal("connected resolution", rec.Code, rec.Body.String())
	}
	for _, x := range []struct {
		token, text string
		status      int
	}{{"operator", "PRIVATE-SERIAL", 404}, {"wrong", "SKU-1", 404}, {"none", "SKU-1", 403}, {"unknown", "SKU-1", 401}, {"factory", "FACTORY-SERIAL", 200}, {"supply", "FACTORY-SERIAL", 404}, {"operator", "https://example.test/secret", 400}} {
		r := resolve(x.token, x.text, "SERIAL")
		if r.Code != x.status {
			t.Fatal(x, r.Code, r.Body.String())
		}
	}
	for _, ten := range []string{tenant, foreign} {
		payload, _ := qr.NewPayload(ten, qr.KindProduct, "variant-a")
		text, _ := payload.Encode()
		r := resolve("operator", text, "QR_REFERENCE")
		expected := 200
		if ten == foreign {
			expected = 404
		}
		if r.Code != expected {
			t.Fatal("tenant-bound QR", r.Code, r.Body.String())
		}
	}
	payload, _ := qr.NewPayload(tenant, qr.KindProduct, "variant-private")
	text, _ := payload.Encode()
	if r := resolve("operator", text, "QR_REFERENCE"); r.Code != 404 {
		t.Fatal("forged checksum granted object")
	}
	payload, _ = qr.NewPayload(tenant, qr.KindCustomer, "person")
	text, _ = payload.Encode()
	if r := resolve("operator", text, "QR_REFERENCE"); r.Code != 404 {
		t.Fatal("unsupported target enabled")
	}
	for _, name := range []string{"tenant-a.png.gray8", "two-distinct.png.gray8", "blank.png.gray8"} {
		b, e := os.ReadFile(filepath.Join(os.Getenv("CAPTURE_TEST_FRAMES"), name))
		if e != nil {
			t.Fatal(e)
		}
		r := call("decode", "operator", "application/octet-stream", b)
		var d cap.Result
		if r.Code != 200 || json.Unmarshal(r.Body.Bytes(), &d) != nil || d.Authorized || d.BusinessEffect {
			t.Fatal("image decode truth")
		}
		expected := map[string]string{"tenant-a.png.gray8": "DECODED", "two-distinct.png.gray8": "MULTIPLE", "blank.png.gray8": "NOT_FOUND"}[name]
		if d.Status != expected {
			t.Fatal(name, d.Status)
		}
	}
	if r := call("decode", "none", "application/octet-stream", frame); r.Code != 403 {
		t.Fatal("unauthorized decode")
	}
	if r := call("decode", "operator", "application/octet-stream", []byte("invalid")); r.Code != 400 {
		t.Fatal("invalid frame")
	}
	if r := call("decode", "operator", "application/octet-stream", make([]byte, cap.MaxFrameBytes+1)); r.Code != 413 {
		t.Fatal("oversize frame")
	}
	var events int
	if e = pool.QueryRow(ctx, `select count(*)from platform.outbox_event where tenant_id=$1`, tenant).Scan(&events); e != nil || events != 0 {
		t.Fatal("capture created business effects", e, events)
	}
	// A database failure is unavailable, never misreported as a missing QR target.
	pool.Close()
	payload, _ = qr.NewPayload(tenant, qr.KindProduct, "variant-a")
	text, _ = payload.Encode()
	if r := resolve("operator", text, "QR_REFERENCE"); r.Code != 503 {
		t.Fatal("database failure hidden as not-found", r.Code)
	}
	t.Log("CAPTURE_CONNECTED_PASS: actual browser-derived raster→HTTP authenticated decoder→ZXing native→authorized PG exact serial; cross tenant recomputed checksum and private object rejected; no business effects; QR/multiple/blank exercised; hardware NOT_RUN")
}
