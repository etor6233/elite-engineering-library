package postgres_test

// AUTHORED real PostgreSQL/HTTP and original SDK fixture composition. Simulated
// detector/provider output cannot establish private document or OCR accuracy.
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
	"io"
	"net/http"
	"net/http/httptest"
	urlpkg "net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type documentVerifier map[string]identity.Principal

func (v documentVerifier) Verify(_ context.Context, token string) (identity.Principal, error) {
	p, ok := v[token]
	if !ok {
		return p, identity.ErrUnauthenticated
	}
	return p, nil
}

type documentProcessorFunc func(context.Context, doc.Original) (doc.Evidence, error)

func (f documentProcessorFunc) Run(ctx context.Context, o doc.Original) (doc.Evidence, error) {
	return f(ctx, o)
}

type documentFixture struct {
	pool            *pgxpool.Pool
	scope           doc.Scope
	maker, reviewer identity.Principal
	original        doc.Original
	pipeline        *doc.Pipeline
}

func newDocumentFixture(t *testing.T) documentFixture {
	t.Helper()
	url := os.Getenv("DOCUMENT_CONNECTED_DB_URL")
	if url == "" {
		t.Skip("owned local document DB required")
	}
	parsed, parseErr := urlpkg.Parse(url)
	if parseErr != nil || (parsed.Hostname() != "127.0.0.1" && parsed.Hostname() != "localhost") || (!strings.HasPrefix(parsed.Path, "/elite_document_reference_") && !strings.HasPrefix(parsed.Path, "/elite_payment_connected_")) {
		t.Fatal("owned loopback fixture database required")
	}
	pool, e := pgxpool.New(context.Background(), url)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(pool.Close)
	tenant := uuid.NewString()
	org := "document-store"
	_, e = pool.Exec(context.Background(), `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,$2,'Fixture','Fixture')`, tenant, "document-"+tenant)
	if e != nil {
		t.Fatal(e)
	}
	_, e = pool.Exec(context.Background(), `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,$2,$2,'Fixture','store')`, tenant, org)
	if e != nil {
		t.Fatal(e)
	}
	raw, e := os.ReadFile(os.Getenv("ELITE_DOCUMENT_REFERENCE_FIXTURE"))
	if e != nil || doc.Hash(raw) != doc.PublicFixtureSHA {
		t.Fatal("exact public fixture missing")
	}
	profile := filepath.Join(os.Getenv("DOCUMENT_TEST_ROOT"), "config", "documents", "reference-profile.json")
	b, e := os.ReadFile(profile)
	if e != nil {
		t.Fatal(e)
	}
	hash := doc.Hash(b)
	processor, e := doc.NewFixturePipeline(profile, hash, t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	principal := func(subject string) identity.Principal {
		return identity.Principal{TenantID: tenant, Subject: subject, Organizations: map[string]struct{}{org: {}}, Permissions: map[string]struct{}{"documents:write": {}, "documents:review": {}, "documents:process": {}}}
	}
	return documentFixture{pool, doc.Scope{TenantID: tenant, OrganizationID: org, ProfileSHA: hash, Mode: "FIXTURE"}, principal("uploader"), principal("reviewer"), doc.Original{Name: "invoice.jpg", SHA256: doc.PublicFixtureSHA, Bytes: raw}, processor}
}
func TestDocumentConnectedReference(t *testing.T) {
	f := newDocumentFixture(t)
	ctx := context.Background()
	var calls atomic.Int32
	service, e := db.NewDocuments(f.pool, f.scope, documentProcessorFunc(func(c context.Context, o doc.Original) (doc.Evidence, error) {
		calls.Add(1)
		return f.pipeline.Run(c, o)
	}))
	if e != nil {
		t.Fatal(e)
	}
	wrong := f.maker
	wrong.Organizations = map[string]struct{}{"other": {}}
	alien := f.maker
	alien.TenantID = uuid.NewString()
	unprivileged := f.maker
	unprivileged.Permissions = map[string]struct{}{}
	mux := http.NewServeMux()
	httpapi.DocumentModule{Store: service}.Register(mux, documentVerifier{"maker": f.maker, "reviewer": f.reviewer, "wrong": wrong, "alien": alien, "none": unprivileged})
	var drop atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, r)
		if drop.Load() && r.Method == "POST" && rec.Code == 200 {
			drop.Store(false)
			conn, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				t.Error(err)
			} else {
				_ = conn.Close()
			}
			return
		}
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		_, _ = w.Write(rec.Body.Bytes())
	}))
	defer server.Close()
	call := func(method, path, token string, body []byte) (int, []byte, error) {
		r, _ := http.NewRequest(method, server.URL+"/v1/documents/"+path, bytes.NewReader(body))
		r.Header.Set("Authorization", "Bearer "+token)
		r.Header.Set("Content-Type", "application/json")
		if method == "PUT" {
			r.Header.Set("Content-Type", "application/octet-stream")
			r.Header.Set("X-Document-Name", f.original.Name)
			r.Header.Set("X-Document-SHA256", f.original.SHA256)
			r.Header.Set("X-Document-Profile-SHA256", f.scope.ProfileSHA)
		}
		response, err := server.Client().Do(r)
		if err != nil {
			return 0, nil, err
		}
		defer response.Body.Close()
		b, err := io.ReadAll(response.Body)
		return response.StatusCode, b, err
	}
	expect := func(method, path, token string, body []byte, status int) db.DocumentView {
		t.Helper()
		code, b, err := call(method, path, token, body)
		if err != nil || code != status {
			t.Fatalf("%s %s: status=%d want=%d err=%v body=%s", method, path, code, status, err, b)
		}
		var v db.DocumentView
		if status < 300 && json.Unmarshal(b, &v) != nil {
			t.Fatal("bad view")
		}
		return v
	}
	raw := func(v any) []byte {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	id := uuid.NewString()
	v := expect("PUT", id+"/original", "maker", f.original.Bytes, 201)
	if v.State != "QUARANTINED" || v.Mode != "FIXTURE" {
		t.Fatal(v)
	}
	expect("PUT", id+"/original", "maker", f.original.Bytes, 201)
	expect("PUT", id+"/original", "reviewer", f.original.Bytes, 409)
	expect("PUT", uuid.NewString()+"/original", "maker", []byte("changed bytes"), 400)
	for _, token := range []string{"wrong", "alien", "none"} {
		expect("GET", id, token, nil, 403)
	}
	expect("GET", id, "unknown", nil, 401)
	code, original, e := call("GET", id+"/original", "maker", nil)
	if e != nil || code != 200 || !bytes.Equal(original, f.original.Bytes) {
		t.Fatal("original retention")
	}
	drop.Store(true)
	if _, _, e = call("POST", id+"/process", "maker", []byte(`{}`)); e == nil {
		t.Fatal("response loss not injected")
	}
	v = expect("POST", id+"/process", "maker", []byte(`{}`), 200)
	if v.State != "REVIEW_REQUIRED" || calls.Load() != 1 || v.Suggested == nil {
		t.Fatal("durable extraction replay", v, calls.Load())
	}
	var storedSecurity []byte
	if e = f.pool.QueryRow(ctx, `select security_receipt from document.extraction where tenant_id=$1 and document_id=$2`, f.scope.TenantID, id).Scan(&storedSecurity); e != nil || !bytes.Contains(storedSecurity, []byte("SIMULATED")) {
		t.Fatal("lost fixture provenance")
	}
	var fields doc.Fields
	if json.Unmarshal(v.Suggested, &fields) != nil {
		t.Fatal("bad suggested fields")
	}
	fields.Vendor = "Explicitly corrected fixture vendor"
	proposal := map[string]any{"evidence_sha256": v.EvidenceSHA, "fields": fields}
	expect("POST", id+"/review", "maker", []byte(`{"evidence_sha256":"a","EVIDENCE_SHA256":"b"}`), 400)
	v = expect("POST", id+"/review", "maker", raw(proposal), 200)
	var reviewed doc.Fields
	if v.Proposal == nil || json.Unmarshal(v.Proposal.Fields, &reviewed) != nil || reviewed.Vendor != fields.Vendor {
		t.Fatal("correction not retained")
	}
	decision := map[string]any{"payload_sha256": v.PayloadSHA, "approved": true, "reason": "Compared four fields with public original; fixture only"}
	expect("POST", id+"/decision", "maker", raw(decision), 409)
	stale := map[string]any{"payload_sha256": strings.Repeat("0", 64), "approved": true, "reason": "stale"}
	expect("POST", id+"/decision", "reviewer", raw(stale), 409)
	// Force an outbox failure and verify decision+commit roll back together.
	_, e = f.pool.Exec(ctx, `create function document.fixture_fail_outbox()returns trigger language plpgsql as $$begin if new.event_type='document.committed'then raise exception 'fixture outbox unavailable';end if;return new;end$$;create trigger fixture_outbox_failure before insert on platform.outbox_event for each row execute function document.fixture_fail_outbox()`)
	if e != nil {
		t.Fatal(e)
	}
	expect("POST", id+"/decision", "reviewer", raw(decision), 503)
	var state string
	var n int
	e = f.pool.QueryRow(ctx, `select state,(select count(*)from document.committed where tenant_id=$1 and document_id=$2)from approval.request where tenant_id=$1 and request_id=$3`, f.scope.TenantID, id, "document:"+id).Scan(&state, &n)
	if e != nil || state != "pending" || n != 0 {
		t.Fatal("partial approval/commit", e, state, n)
	}
	_, e = f.pool.Exec(ctx, `drop trigger fixture_outbox_failure on platform.outbox_event;drop function document.fixture_fail_outbox()`)
	if e != nil {
		t.Fatal(e)
	}
	drop.Store(true)
	if _, _, e = call("POST", id+"/decision", "reviewer", raw(decision)); e == nil {
		t.Fatal("approval response loss not injected")
	}
	v = expect("POST", id+"/decision", "reviewer", raw(decision), 200)
	if v.State != "PERSISTED" || v.Reviewer != f.reviewer.Subject {
		t.Fatal(v)
	}
	decision["reason"] = "changed replay"
	expect("POST", id+"/decision", "reviewer", raw(decision), 409)
	for _, q := range []string{
		`update document.original set name='changed'where tenant_id=$1 and document_id=$2`,
		`delete from document.extraction where tenant_id=$1 and document_id=$2`,
		`update document.committed set reviewer='changed'where tenant_id=$1 and document_id=$2`,
	} {
		if _, e = f.pool.Exec(ctx, q, f.scope.TenantID, id); e == nil {
			t.Fatal("immutable document changed")
		}
	}
	if _, e = f.pool.Exec(ctx, `update approval.request set payload=payload||'{"unbound":true}'::jsonb where tenant_id=$1 and request_id=$2`, f.scope.TenantID, "document:"+id); e == nil {
		t.Fatal("immutable proposal changed")
	}
	e = f.pool.QueryRow(ctx, `select count(*)from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='document.committed'`, f.scope.TenantID, id).Scan(&n)
	if e != nil || n != 1 {
		t.Fatal("duplicate commit outbox", e, n)
	}
	// A rejection is terminal and has no persisted business document.
	rejected := uuid.NewString()
	if _, e = service.Receive(ctx, f.maker, rejected, f.scope.ProfileSHA, f.original); e != nil {
		t.Fatal(e)
	}
	rv, e := service.Process(ctx, f.maker, rejected)
	if e != nil {
		t.Fatal(e)
	}
	rv, e = service.Submit(ctx, f.maker, rejected, rv.EvidenceSHA, fields)
	if e != nil {
		t.Fatal(e)
	}
	rv, e = service.Decide(ctx, f.reviewer, rejected, rv.PayloadSHA, false, "fixture rejected")
	if e != nil || rv.State != "REJECTED" {
		t.Fatal(e, rv)
	}
	if _, e = service.Decide(ctx, f.reviewer, rejected, rv.PayloadSHA, true, "changed decision"); e == nil {
		t.Fatal("rejection reversed")
	}
	down, e := os.ReadFile(filepath.Join(os.Getenv("DOCUMENT_TEST_ROOT"), "db", "migrations", "0085_document_review.down.sql"))
	if e != nil {
		t.Fatal(e)
	}
	conn, e := f.pool.Acquire(ctx)
	if e != nil {
		t.Fatal(e)
	}
	_, e = conn.Exec(ctx, string(down))
	_, _ = conn.Exec(ctx, "rollback")
	conn.Release()
	if e == nil {
		t.Fatal("populated down migration erased evidence")
	}
	t.Log("DOCUMENT_CONNECTED_PASS: original+job atomic; scoped HTTP; original SDK; explicit simulated detectors; corrected immutable review; independent reviewer; commit+approval+outbox atomic; lost response replay; rejection; immutable evidence; populated rollback refused")
}
func TestDocumentRecoveryAndFencing(t *testing.T) {
	f := newDocumentFixture(t)
	ctx := context.Background()
	var calls atomic.Int32
	wrapper := documentProcessorFunc(func(c context.Context, o doc.Original) (doc.Evidence, error) {
		if calls.Add(1) == 1 {
			return doc.Evidence{}, doc.ErrExtraction
		}
		return f.pipeline.Run(c, o)
	})
	store, e := db.NewDocuments(f.pool, f.scope, wrapper)
	if e != nil {
		t.Fatal(e)
	}
	id := uuid.NewString()
	if _, e = store.Receive(ctx, f.maker, id, f.scope.ProfileSHA, f.original); e != nil {
		t.Fatal(e)
	}
	if _, e = store.Process(ctx, f.maker, id); !errors.Is(e, doc.ErrExtraction) {
		t.Fatal(e)
	}
	ready := func(id string) {
		t.Helper()
		_, err := f.pool.Exec(ctx, `update platform.job set available_at=clock_timestamp()-interval '1 second' where tenant_id=$1 and payload->>'document_id'=$2`, f.scope.TenantID, id)
		if err != nil {
			t.Fatal(err)
		}
	}
	ready(id)
	v, e := store.Process(ctx, f.maker, id)
	if e != nil || v.State != "REVIEW_REQUIRED" || calls.Load() != 2 {
		t.Fatal(v, e, calls.Load())
	}
	var count int
	e = f.pool.QueryRow(ctx, `select count(*)from document.attempt_failure where tenant_id=$1 and document_id=$2`, f.scope.TenantID, id).Scan(&count)
	if e != nil || count != 1 {
		t.Fatal("failure receipt lost", e, count)
	}
	// A delayed old generation may finish its extraction, but cannot commit it.
	started, release := make(chan struct{}), make(chan struct{})
	var attempts atomic.Int32
	fenced, e := db.NewDocuments(f.pool, f.scope, documentProcessorFunc(func(c context.Context, o doc.Original) (doc.Evidence, error) {
		if attempts.Add(1) == 1 {
			close(started)
			select {
			case <-release:
			case <-c.Done():
				return doc.Evidence{}, c.Err()
			}
		}
		return f.pipeline.Run(c, o)
	}))
	if e != nil {
		t.Fatal(e)
	}
	stale := uuid.NewString()
	if _, e = fenced.Receive(ctx, f.maker, stale, f.scope.ProfileSHA, f.original); e != nil {
		t.Fatal(e)
	}
	done := make(chan error, 1)
	go func() { _, err := fenced.Process(ctx, f.maker, stale); done <- err }()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("worker not started")
	}
	_, e = f.pool.Exec(ctx, `update platform.job set claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and payload->>'document_id'=$2`, f.scope.TenantID, stale)
	if e != nil {
		t.Fatal(e)
	}
	v, e = fenced.Process(ctx, f.maker, stale)
	close(release)
	if e != nil || v.State != "REVIEW_REQUIRED" {
		t.Fatal(e, v)
	}
	select {
	case err := <-done:
		if !errors.Is(err, db.ErrJobClaimLost) {
			t.Fatal("stale worker committed", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("stale worker hung")
	}
	e = f.pool.QueryRow(ctx, `select count(*)from document.extraction where tenant_id=$1 and document_id=$2 and attempt=2`, f.scope.TenantID, stale).Scan(&count)
	if e != nil || count != 1 {
		t.Fatal("fence lost", e, count)
	}
	// Security refusal never reaches the official provider; finite attempts stop.
	blocked, e := db.NewDocuments(f.pool, f.scope, documentProcessorFunc(func(context.Context, doc.Original) (doc.Evidence, error) { return doc.Evidence{}, doc.ErrSecurity }))
	if e != nil {
		t.Fatal(e)
	}
	bad := uuid.NewString()
	if _, e = blocked.Receive(ctx, f.maker, bad, f.scope.ProfileSHA, f.original); e != nil {
		t.Fatal(e)
	}
	for i := 0; i < 3; i++ {
		ready(bad)
		if _, e = blocked.Process(ctx, f.maker, bad); !errors.Is(e, doc.ErrSecurity) {
			t.Fatal(e)
		}
	}
	v, e = blocked.Read(ctx, f.maker, bad)
	if e != nil || v.State != "QUARANTINE_TERMINAL" {
		t.Fatal(e, v)
	}
	// A crashed last attempt is exhausted by the original job owner, not reset.
	exhausted := uuid.NewString()
	if _, e = blocked.Receive(ctx, f.maker, exhausted, f.scope.ProfileSHA, f.original); e != nil {
		t.Fatal(e)
	}
	_, e = f.pool.Exec(ctx, `update platform.job set attempts=max_attempts,claimed_by='crashed-fixture',claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and payload->>'document_id'=$2`, f.scope.TenantID, exhausted)
	if e != nil {
		t.Fatal(e)
	}
	v, e = blocked.Process(ctx, f.maker, exhausted)
	if e != nil || v.State != "QUARANTINE_TERMINAL" {
		t.Fatal(e, v)
	}
	t.Log("DOCUMENT_RECOVERY_PASS: durable transient failure; bounded retry; generation fencing; security terminal quarantine; exhausted crashed attempt retained")
}

func TestDocumentProfileHistoryAndDatabaseReviewContract(t *testing.T) {
	f := newDocumentFixture(t)
	ctx := context.Background()
	store, e := db.NewDocuments(f.pool, f.scope, f.pipeline)
	if e != nil {
		t.Fatal(e)
	}
	id := uuid.NewString()
	if _, e = store.Receive(ctx, f.maker, id, f.scope.ProfileSHA, f.original); e != nil {
		t.Fatal(e)
	}
	changed := f.scope
	changed.ProfileSHA = strings.Repeat("b", 64)
	newStore, e := db.NewDocuments(f.pool, changed, f.pipeline)
	if e != nil {
		t.Fatal(e)
	}
	v, e := newStore.Read(ctx, f.maker, id)
	if e != nil || v.State != "QUARANTINED" || v.ProfileSHA != f.scope.ProfileSHA {
		t.Fatal("stored profile history unreadable", e, v)
	}
	if _, e = newStore.Process(ctx, f.maker, id); !errors.Is(e, db.ErrDocumentConflict) {
		t.Fatal("stale profile processed", e)
	}
	v, e = store.Process(ctx, f.maker, id)
	if e != nil {
		t.Fatal(e)
	}
	payload := map[string]any{"schema": "document-review/v1", "document_id": id, "organization_id": f.scope.OrganizationID, "original_sha256": v.OriginalSHA, "profile_sha256": v.ProfileSHA, "evidence_sha256": v.EvidenceSHA, "mode": "FIXTURE", "fields": map[string]string{"invoice_number": "only one field"}}
	raw, _ := json.Marshal(payload)
	if _, e = f.pool.Exec(ctx, `insert into approval.request(tenant_id,request_id,kind,subject_id,amount_minor_units,requester,evidence_sha,organization_id,payload)values($1,$2,'document_review',$3,0,'fixture-internal-caller',$4,$5,$6)`, f.scope.TenantID, "document:"+id, id, doc.Hash(raw), f.scope.OrganizationID, raw); e == nil {
		t.Fatal("database admitted incomplete fields")
	}
	// Restart evidence: earlier connected run rows and corrected commit survive.
	var n int
	e = f.pool.QueryRow(ctx, `select count(*)from document.committed c join document.original o using(tenant_id,document_id)where o.mode='FIXTURE'and encode(sha256(o.content),'hex')=o.original_sha256 and c.fields->>'vendor'='Explicitly corrected fixture vendor'`).Scan(&n)
	if e != nil || n < 1 {
		t.Fatal("prior committed reference did not survive host restart", e, n)
	}

	mux := http.NewServeMux()
	httpapi.DocumentModule{Store: store}.Register(mux, documentVerifier{"maker": f.maker})
	for _, part := range []string{"security", "provider", "analysis"} {
		req := httptest.NewRequest("GET", "/v1/documents/"+id+"/evidence/"+part, nil)
		req.Header.Set("Authorization", "Bearer maker")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		exact, err := store.EvidencePart(ctx, f.maker, id, part)
		if err != nil || rec.Code != 200 || !bytes.Equal(exact, rec.Body.Bytes()) || rec.Header().Get("X-Document-Evidence-SHA256") != doc.Hash(exact) || rec.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Fatal("exact evidence review failed", part, err, rec.Code)
		}
	}
	if _, e = store.EvidencePart(ctx, f.maker, id, "../original"); !errors.Is(e, doc.ErrContract) {
		t.Fatal("unbounded evidence part", e)
	}
	t.Log("DOCUMENT_REVIEW_DELTA_PASS: stored-profile history; stale execution refusal; database field contract; original and commit survive process restart")
}
