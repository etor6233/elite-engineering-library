package postgres_test

// AUTHORED real HTTP/PostgreSQL fixture tests. No OCR engine or live business posting.
import (
	"bytes"
	"context"
	doc "elite.local/enterprise/internal/documentbridge"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
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

func TestDocumentTyped21Connected(t *testing.T) {
	dsn := os.Getenv("DOCUMENT_TYPED_DB_URL")
	if dsn == "" {
		t.Skip("owned typed document fixture database required")
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
	org := "typed-fixture-store"
	_, e = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,$2,'Fixture','Fixture')`, tenant, "typed-"+tenant)
	if e != nil {
		t.Fatal(e)
	}
	_, e = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,$2,$2,'Fixture','store')`, tenant, org)
	if e != nil {
		t.Fatal(e)
	}
	principal := func(subject string, perms ...string) identity.Principal {
		p := identity.Principal{TenantID: tenant, Subject: subject, Organizations: map[string]struct{}{org: {}}, Permissions: map[string]struct{}{}}
		for _, s := range perms {
			p.Permissions[s] = struct{}{}
		}
		return p
	}
	maker := principal("typed-maker", "documents:write", "documents:process")
	reviewer := principal("typed-reviewer", "documents:review")
	reader := principal("typed-maker", "documents:read")
	other := principal("typed-other", "documents:read")
	wrong := maker
	wrong.Organizations = map[string]struct{}{"other": {}}
	alien := maker
	alien.TenantID = uuid.NewString()
	profile := filepath.Join(os.Getenv("DOCUMENT_TEST_ROOT"), "config/documents/typed-profile.json")
	processor, e := doc.NewTypedFixturePipeline(profile, doc.TypedProfileSHA())
	if e != nil {
		t.Fatal(e)
	}
	scope := doc.Scope{TenantID: tenant, OrganizationID: org, ProfileSHA: doc.TypedProfileSHA(), Mode: "TYPED_FIXTURE"}
	store, e := db.NewDocuments(pool, scope, processor)
	if e != nil {
		t.Fatal(e)
	}
	mux := http.NewServeMux()
	httpapi.DocumentModule{Store: store}.Register(mux, documentVerifier{"maker": maker, "reviewer": reviewer, "reader": reader, "other": other, "wrong": wrong, "alien": alien})
	call := func(method, path, token string, body []byte, original *doc.Original) *httptest.ResponseRecorder {
		t.Helper()
		r := httptest.NewRequest(method, "/v1/documents"+path, bytes.NewReader(body))
		r.Header.Set("Authorization", "Bearer "+token)
		r.Header.Set("Content-Type", "application/json")
		if original != nil {
			r.Header.Set("Content-Type", "application/octet-stream")
			r.Header.Set("X-Document-Name", original.Name)
			r.Header.Set("X-Document-SHA256", original.SHA256)
			r.Header.Set("X-Document-Profile-SHA256", scope.ProfileSHA)
			r.Header.Set("X-Document-Class", original.ClassID)
			r.Header.Set("X-Document-Schema-Version", original.SchemaVersion)
		}
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, r)
		return rec
	}
	raw := func(v any) []byte {
		b, e := json.Marshal(v)
		if e != nil {
			t.Fatal(e)
		}
		return b
	}
	view := func(rec *httptest.ResponseRecorder, status int) db.DocumentView {
		t.Helper()
		if rec.Code != status {
			t.Fatalf("status %d want %d: %s", rec.Code, status, rec.Body.String())
		}
		var v db.DocumentView
		if status < 300 && json.Unmarshal(rec.Body.Bytes(), &v) != nil {
			t.Fatal("view")
		}
		return v
	}
	catalogRec := call("GET", "/classes", "reader", nil, nil)
	if catalogRec.Code != 200 {
		t.Fatal(catalogRec.Code)
	}
	var catalog doc.ClassCatalog
	if json.Unmarshal(catalogRec.Body.Bytes(), &catalog) != nil || len(catalog.Classes) != 21 || catalog.ProfileSHA != scope.ProfileSHA {
		t.Fatal("catalog")
	}
	ids := []string{}
	for _, class := range catalog.Classes {
		t.Run(class.ID, func(t *testing.T) {
			id := uuid.NewString()
			ids = append(ids, id)
			original, e := doc.TypedFixture(class.ID)
			if e != nil {
				t.Fatal(e)
			}
			v := view(call("PUT", "/"+id+"/original", "maker", original.Bytes, &original), 201)
			if v.State != "QUARANTINED" || v.ClassID != class.ID || v.SchemaVersion != "1" {
				t.Fatal(v)
			}
			view(call("PUT", "/"+id+"/original", "maker", original.Bytes, &original), 201)
			if _, e = store.Read(ctx, reader, id); e != nil {
				t.Fatal("read-only owner scope rejected", e)
			}
			view(call("GET", "/"+id, "other", nil, nil), 404)
			v = view(call("POST", "/"+id+"/process", "maker", []byte(`{}`), nil), 200)
			if v.State != "REVIEW_REQUIRED" || doc.ValidateTypedFields(v.ClassID, v.SchemaVersion, v.Suggested) != nil {
				t.Fatal(v)
			}
			// No approval/commit is fabricated by receive or process.
			var count int
			if e = pool.QueryRow(ctx, `select count(*)from document.committed where tenant_id=$1 and document_id=$2`, tenant, id).Scan(&count); e != nil || count != 0 {
				t.Fatal("automatic persistence", e, count)
			}
			view(call("POST", "/"+id+"/process", "maker", []byte(`{}`), nil), 200)
			corrected := map[string]string{}
			if json.Unmarshal(v.Suggested, &corrected) != nil {
				t.Fatal("suggested")
			}
			for _, f := range class.Fields {
				if f.Type == "text" {
					corrected[f.Key] = "Reviewed fixture value"
					break
				}
			}
			proposal := raw(map[string]any{"evidence_sha256": v.EvidenceSHA, "fields": corrected})
			v = view(call("POST", "/"+id+"/review", "maker", proposal, nil), 200)
			view(call("POST", "/"+id+"/review", "maker", proposal, nil), 200)
			if v.Proposal == nil || v.Proposal.Schema != "document-review/v2" || v.Proposal.ClassID != class.ID {
				t.Fatal("typed review binding")
			}
			decision := raw(map[string]any{"payload_sha256": v.PayloadSHA, "approved": true, "reason": "Independent comparison of synthetic original and reviewed fields"})
			view(call("POST", "/"+id+"/decision", "maker", decision, nil), 403)
			// A principal holding both permissions still cannot approve their own proposal.
			both := maker
			both.Permissions = map[string]struct{}{"documents:write": {}, "documents:process": {}, "documents:review": {}}
			if _, e = store.Decide(ctx, both, id, v.PayloadSHA, true, "self approval"); e == nil {
				t.Fatal("same actor approved")
			}
			discarded := call("POST", "/"+id+"/decision", "reviewer", decision, nil)
			if discarded.Code != 200 {
				t.Fatal("decision", discarded.Body.String())
			}
			// Ignore the response and construct a fresh service, then recover only by GET.
			restarted, e := db.NewDocuments(pool, scope, processor)
			if e != nil {
				t.Fatal(e)
			}
			recovered, e := restarted.Read(ctx, reader, id)
			if e != nil || recovered.State != "PERSISTED" {
				t.Fatal("restart recovery", e, recovered)
			}
			v = view(call("POST", "/"+id+"/decision", "reviewer", decision, nil), 200)
			if v.State != "PERSISTED" {
				t.Fatal(v)
			}
			var stored []byte
			e = pool.QueryRow(ctx, `select fields from document.committed where tenant_id=$1 and document_id=$2`, tenant, id).Scan(&stored)
			if e != nil {
				t.Fatal(e)
			}
			var storedFields map[string]string
			if json.Unmarshal(stored, &storedFields) != nil || len(storedFields) != len(corrected) {
				t.Fatal("fields")
			}
			for k, val := range corrected {
				if storedFields[k] != val {
					t.Fatal("reviewed field changed")
				}
			}
			download := call("GET", "/"+id+"/original", "reader", nil, nil)
			if download.Code != 200 || !bytes.Equal(download.Body.Bytes(), original.Bytes) || download.Header().Get("Content-Disposition") != "attachment; filename=\"original.bin\"" || download.Header().Get("Content-Type") != "application/octet-stream" || download.Header().Get("Cache-Control") != "no-store" || download.Header().Get("X-Content-Type-Options") != "nosniff" {
				t.Fatal("unsafe or changed original")
			}
			for _, part := range []string{"security", "provider", "analysis"} {
				rec := call("GET", "/"+id+"/evidence/"+part, "reviewer", nil, nil)
				if rec.Code != 200 || !strings.Contains(rec.Header().Get("Content-Disposition"), "attachment") {
					t.Fatal("evidence", part)
				}
			}
			if e = pool.QueryRow(ctx, `select count(*)from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='document.committed'`, tenant, id).Scan(&count); e != nil || count != 1 {
				t.Fatal("outbox duplicate", e, count)
			}
		})
	}
	if len(ids) != 21 {
		t.Fatal("missing classes")
	}
	for _, token := range []string{"wrong", "alien"} {
		view(call("GET", "/"+ids[0], token, nil, nil), 403)
		if rec := call("GET", "?limit=2", token, nil, nil); rec.Code != 403 {
			t.Fatal("scope list", token)
		}
	}
	empty := call("GET", "?limit=5", "other", nil, nil)
	var inbox db.DocumentInbox
	if json.Unmarshal(empty.Body.Bytes(), &inbox) != nil || len(inbox.Items) != 0 {
		t.Fatal("list disclosed other uploader")
	}
	seen := map[string]bool{}
	cursor := ""
	for {
		path := "?limit=4"
		if cursor != "" {
			path += "&cursor=" + url.QueryEscape(cursor)
		}
		rec := call("GET", path, "reader", nil, nil)
		if rec.Code != 200 || json.Unmarshal(rec.Body.Bytes(), &inbox) != nil {
			t.Fatal("inbox", rec.Body.String())
		}
		for _, v := range inbox.Items {
			if seen[v.ID] {
				t.Fatal("duplicate page row")
			}
			seen[v.ID] = true
		}
		if inbox.NextCursor == nil {
			break
		}
		cursor = *inbox.NextCursor
	}
	if len(seen) != 21 {
		t.Fatal("missing pagination rows", len(seen))
	}
	for _, query := range []string{"?limit=0", "?limit=51", "?limit=1&limit=2", "?cursor=bad", "?unexpected=1"} {
		view(call("GET", query, "reader", nil, nil), 400)
	}
	if rec := call("GET", "?limit=4&cursor="+url.QueryEscape(cursor), "other", nil, nil); rec.Code != 400 {
		t.Fatal("cursor scope mismatch not rejected")
	}
	original, _ := doc.TypedFixture("receipt")
	wrongOriginal := original
	wrongOriginal.ClassID = "packing-list"
	if _, e = store.Receive(ctx, maker, uuid.NewString(), scope.ProfileSHA, wrongOriginal); !errors.Is(e, doc.ErrContract) {
		t.Fatal("cross class admitted", e)
	}
	// Direct SQL must not bypass exact fixture class/hash and review field contracts.
	_, e = pool.Exec(ctx, `insert into document.original(tenant_id,document_id,organization_id,uploader,name,original_sha256,profile_sha256,mode,content,class_id,schema_version)values($1,$2,$3,$4,$5,$6,$7,'TYPED_FIXTURE',$8,'packing-list','1')`, tenant, uuid.NewString(), org, maker.Subject, original.Name, original.SHA256, scope.ProfileSHA, original.Bytes)
	if e == nil {
		t.Fatal("DB accepted wrong fixture class")
	}
	id := uuid.NewString()
	if _, e = store.Receive(ctx, maker, id, scope.ProfileSHA, original); e != nil {
		t.Fatal(e)
	}
	v, e := store.Process(ctx, maker, id)
	if e != nil {
		t.Fatal(e)
	}
	malformed := raw(map[string]any{"schema": "document-review/v2", "document_id": id, "organization_id": org, "original_sha256": v.OriginalSHA, "profile_sha256": v.ProfileSHA, "evidence_sha256": v.EvidenceSHA, "mode": "TYPED_FIXTURE", "class_id": "receipt", "schema_version": "1", "fields": map[string]string{"receipt_number": "missing fields"}})
	_, e = pool.Exec(ctx, `insert into approval.request(tenant_id,request_id,kind,subject_id,amount_minor_units,requester,evidence_sha,organization_id,payload)values($1,$2,'document_review',$3,0,$4,$5,$6,$7)`, tenant, "document:"+id, id, maker.Subject, doc.Hash(malformed), org, malformed)
	if e == nil {
		t.Fatal("DB accepted invalid class fields")
	}
	v, e = store.SubmitFields(ctx, maker, id, v.EvidenceSHA, v.Suggested)
	if e != nil {
		t.Fatal(e)
	}
	v, e = store.Decide(ctx, reviewer, id, v.PayloadSHA, false, "Synthetic rejection retained")
	if e != nil || v.State != "REJECTED" {
		t.Fatal(e)
	}
	if _, e = store.Decide(ctx, reviewer, id, v.PayloadSHA, true, "reverse rejection"); e == nil {
		t.Fatal("reversed decision")
	}
	down, e := os.ReadFile(filepath.Join(os.Getenv("DOCUMENT_TEST_ROOT"), "db/migrations/0089_document_classes.down.sql"))
	if e != nil {
		t.Fatal(e)
	}
	conn, e := pool.Acquire(ctx)
	if e != nil {
		t.Fatal(e)
	}
	_, e = conn.Exec(ctx, string(down))
	_, _ = conn.Exec(ctx, "rollback")
	conn.Release()
	if e == nil {
		t.Fatal("populated typed rollback erased evidence")
	}
	t.Log("TYPED21_CONNECTED_PASS: 21 exact PDF originals,21 typed proposals,21 independent review commits, scoped inbox, durable recovery, exact attachment download, DB negative guards, no automatic financial effect")
}
