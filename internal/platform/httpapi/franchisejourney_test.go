package httpapi

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type browserAppointmentClock struct{}

// AUTHORED local integration fixture. Real RS256/JWKS verification, not a
// successful verifier double. No live IdP login or production identity claim.
func confirmationTestIssuer(t *testing.T) (identity.Verifier, func(string, string, []string, []string) string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	encode := base64.RawURLEncoding.EncodeToString
	var issuer *httptest.Server
	issuer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			_ = json.NewEncoder(w).Encode(map[string]any{"issuer": issuer.URL, "jwks_uri": issuer.URL + "/keys", "id_token_signing_alg_values_supported": []string{"RS256"}})
		case "/keys":
			_ = json.NewEncoder(w).Encode(map[string]any{"keys": []any{map[string]any{"kty": "RSA", "use": "sig", "alg": "RS256", "kid": "local-confirmation", "n": encode(key.N.Bytes()), "e": encode(big.NewInt(int64(key.E)).Bytes())}}})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(issuer.Close)
	verifier, err := identity.NewOIDCVerifier(context.Background(), issuer.URL, "confirmation-api")
	if err != nil {
		t.Fatal(err)
	}
	token := func(subject, tenant string, permissions, organizations []string) string {
		header, err := json.Marshal(map[string]string{"alg": "RS256", "kid": "local-confirmation", "typ": "JWT"})
		if err != nil {
			t.Fatal(err)
		}
		claims, err := json.Marshal(map[string]any{"iss": issuer.URL, "aud": "confirmation-api", "sub": subject, "iat": time.Now().Unix(), "exp": time.Now().Add(5 * time.Minute).Unix(), "tenant_id": tenant, "permissions": permissions, "organization_ids": organizations})
		if err != nil {
			t.Fatal(err)
		}
		unsigned := encode(header) + "." + encode(claims)
		digest := sha256.Sum256([]byte(unsigned))
		signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
		if err != nil {
			t.Fatal(err)
		}
		return unsigned + "." + encode(signature)
	}
	return verifier, token
}

func TestAppointmentConfirmationBFFPostgres(t *testing.T) {
	runAppointmentConfirmationPostgres(t, false)
}

func TestAppointmentAgendaBrowserPostgres(t *testing.T) {
	for _, project := range []string{"chromium-desktop", "chromium-mobile", "firefox-desktop", "webkit-desktop"} {
		t.Run(project, func(t *testing.T) {
			t.Setenv("ELITE_AGENDA_BROWSER_PROJECT", project)
			runAppointmentConfirmationPostgres(t, true)
		})
	}
}

func runAppointmentConfirmationPostgres(t *testing.T, browser bool) {
	t.Helper()
	if os.Getenv("ELITE_CONFIRMATION_E2E") != "1" {
		t.Skip("explicit disposable confirmation gate not requested")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	config, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal("invalid test database configuration")
	}
	// Immutable audit rows deliberately remain for inspection until the caller
	// discards this dedicated database. Never disable their protective triggers.
	if config.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(config.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires a disposable loopback elite_confirmation_* database, never a project database")
	}
	web := os.Getenv("ELITE_WEB_ROOT")
	if !filepath.IsAbs(web) {
		t.Fatal("ELITE_WEB_ROOT must be absolute")
	}
	cli := filepath.Join(web, "node_modules", "vitest", "vitest.mjs")
	if _, err := os.Stat(cli); err != nil {
		t.Fatal("materialize web and install its exact frozen lock first")
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := (randomid.Generator{}).New()
	code := "confirmation-" + strings.ReplaceAll(tenant, "-", "")
	if _, err := pool.Exec(ctx, "insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,$2,'Synthetic confirmation','Synthetic confirmation')", tenant, code); err != nil {
		t.Fatal(err)
	}
	for _, sql := range []string{
		"insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic Store','store')",
		"insert into org.public_location(tenant_id,organization_id,city,region,country,published)values($1,'store','Fixture City','Fixture Region','AR',true)",
		"insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic Customer','confirmation@example.invalid')",
		"insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','customer','new','public-web','{}')",
	} {
		if _, err := pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	verifier, token := confirmationTestIssuer(t)
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: franchisejourney.NewService(postgres.NewFranchiseJourney(pool), randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
	api := httptest.NewServer(mux)
	defer api.Close()
	operator := token("operator", tenant, []string{"appointment:manage", "resource:manage", "availability:manage"}, []string{"store"})
	post := func(path, bearer string, payload any, expected int) map[string]any {
		t.Helper()
		body, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequestWithContext(ctx, "POST", api.URL+path, strings.NewReader(string(body)))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+bearer)
		req.Header.Set("Idempotency-Key", "confirmation-request-key")
		if path == "/v1/franchise/availability" {
			req.Header.Set("Idempotency-Key", "availability-"+(randomid.Generator{}).New())
		}
		res, err := api.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		var value map[string]any
		if err := json.NewDecoder(res.Body).Decode(&value); err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != expected {
			t.Fatalf("%s expected=%d actual=%d response=%v", path, expected, res.StatusCode, value)
		}
		return value
	}
	start := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	end := start.Add(30 * time.Minute)
	working := map[string]any{"organization_id": "store", "entry_type": "working", "starts_at": start.Add(-time.Hour), "ends_at": end.Add(time.Hour)}
	post("/v1/franchise/availability", operator, working, 201)
	post("/v1/franchise/appointment-slots", operator, map[string]any{"organization_id": "store", "kind": "consultation", "starts_at": start, "ends_at": end, "capacity": 1}, 201)
	appointment := post("/v1/public/"+code+"/store/appointments", "", map[string]any{"lead_id": "lead", "kind": "consultation", "starts_at": start}, 202)
	id, ok := appointment["id"].(string)
	if !ok || id == "" {
		t.Fatal("missing appointment receipt")
	}
	transitionPath := "/v1/franchise/appointments/" + id + "/transitions"
	confirmation := map[string]any{"organization_id": "store", "current": "requested", "target": "confirmed", "version": 1}
	for _, negative := range []struct {
		name, bearer string
		status       int
	}{
		{"invalid-signature", operator[:strings.LastIndex(operator, ".")+1] + "AAAA", 401},
		{"no-permission", token("observer", tenant, []string{"lead:read"}, []string{"store"}), 403},
		{"other-organization", token("operator", tenant, []string{"appointment:manage"}, []string{"other"}), 403},
		{"other-tenant", token("operator", (randomid.Generator{}).New(), []string{"appointment:manage"}, []string{"store"}), 409},
		{"no-resource", operator, 409},
	} {
		t.Run(negative.name, func(t *testing.T) { post(transitionPath, negative.bearer, confirmation, negative.status) })
	}
	resource := post("/v1/franchise/resources", operator, map[string]any{"organization_id": "store", "principal_subject": "technician", "display_name": "Synthetic technician", "kind": "employee", "skills": []string{"consultation"}}, 201)
	assignment := map[string]any{"organization_id": "store", "resource_id": resource["id"], "version": 1}
	assignmentPath := "/v1/franchise/appointments/" + id + "/resources"
	post(assignmentPath, operator, assignment, 409) // no resource working window
	working["resource_id"] = resource["id"]
	post("/v1/franchise/availability", operator, working, 201)
	assignment["version"] = 2
	post(assignmentPath, operator, assignment, 409) // stale version, no effect
	assignment["version"] = 1
	if !browser {
		assigned := post(assignmentPath, operator, assignment, 200)
		if assigned["version"] != float64(2) || assigned["state"] != "requested" {
			t.Fatalf("invalid assignment receipt: %v", assigned)
		}
	}
	// Real server-side BFF client forwards signed tokens to this API. The test
	// races confirmations then reads the customer's durable, scoped timeline.
	cmd := exec.CommandContext(ctx, "node", cli, "run", "src/platform/backend/protected-client.test.ts", "-t", "connected appointment confirmation")
	cmd.Dir = web
	if browser {
		cmd = exec.CommandContext(ctx, "node", filepath.Join(web, "microsoft_playwright_browser_gate", "node_modules", "@playwright", "test", "cli.js"), "test", "tests/enterprise-web.spec.mjs", "--grep", "operator agenda confirms", "--project", os.Getenv("ELITE_AGENDA_BROWSER_PROJECT"), "--workers=1", "--retries=0", "--max-failures=1")
		cmd.Dir = filepath.Join(web, "microsoft_playwright_browser_gate")
	}
	for _, item := range os.Environ() {
		name := strings.ToUpper(strings.SplitN(item, "=", 2)[0])
		if name == "TEST_DATABASE_URL" || name == "DATABASE_URL" || strings.HasPrefix(name, "ELITE_CONFIRMATION_") || name == "ENTERPRISE_API_BASE_URL" || name == "AUTH_SESSION_SECRET" || name == "ELITE_BASE_URL" || name == "ELITE_RUNTIME_ONLY" || name == "ELITE_OPERATOR_AGENDA_E2E" || name == "ELITE_WEB_ROOT" {
			continue
		}
		cmd.Env = append(cmd.Env, item)
	}
	cmd.Env = append(cmd.Env, "ENTERPRISE_API_BASE_URL="+api.URL, "ELITE_CONFIRMATION_E2E=1", "ELITE_CONFIRMATION_ID="+id, "ELITE_CONFIRMATION_TENANT="+tenant, "ELITE_CONFIRMATION_OPERATOR="+operator, "ELITE_CONFIRMATION_CUSTOMER="+token("customer", tenant, []string{"customer:self"}, []string{"store"}), "ELITE_CONFIRMATION_STRANGER="+token("stranger", tenant, []string{"customer:self"}, []string{"store"}))
	if browser {
		cmd.Env = append(cmd.Env, "ELITE_OPERATOR_AGENDA_E2E=1", "ELITE_WEB_ROOT="+web, "APP_BASE_URL=https://127.0.0.1:4173", "ELITE_CONFIRMATION_DAY="+start.Format("2006-01-02"), "AUTH_SESSION_SECRET=synthetic-"+(randomid.Generator{}).New()+(randomid.Generator{}).New())
	}
	if browser {
		target, _ := url.Parse("http://127.0.0.1:4173")
		edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
		defer edge.Close()
		filtered := cmd.Env[:0]
		for _, item := range cmd.Env {
			if !strings.HasPrefix(item, "APP_BASE_URL=") && !strings.HasPrefix(item, "ELITE_BASE_URL=") {
				filtered = append(filtered, item)
			}
		}
		cmd.Env = append(filtered, "APP_BASE_URL="+edge.URL, "ELITE_BASE_URL="+edge.URL)
		server := exec.CommandContext(ctx, "node", filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", "4173")
		server.Dir = web
		server.Env = cmd.Env
		if err := server.Start(); err != nil {
			t.Fatal(err)
		}
		defer func() { _ = server.Process.Kill(); _ = server.Wait() }()
		client := &http.Client{Timeout: time.Second}
		ready := false
		for deadline := time.Now().Add(15 * time.Second); time.Now().Before(deadline); {
			response, err := client.Get(target.String() + "/icon.svg")
			if err == nil {
				response.Body.Close()
				if response.StatusCode == 200 {
					ready = true
					break
				}
			}
			select {
			case <-ctx.Done():
				t.Fatal("web fixture startup cancelled")
			case <-time.After(100 * time.Millisecond):
			}
		}
		if !ready {
			t.Fatal("web fixture did not become ready")
		}
	}
	output, err := cmd.CombinedOutput()
	// Never log signed fixture tokens, even on assertion failures.
	clean := string(output)
	for _, item := range cmd.Env {
		parts := strings.SplitN(item, "=", 2)
		if strings.HasPrefix(parts[0], "ELITE_CONFIRMATION_") && len(parts) == 2 && strings.Count(parts[1], ".") == 2 {
			clean = strings.ReplaceAll(clean, parts[1], "[REDACTED_TEST_TOKEN]")
		}
	}
	t.Log(clean)
	if err != nil {
		t.Fatalf("BFF confirmation gate failed: %v", err)
	}
	var state, actor string
	var version int
	if err := pool.QueryRow(ctx, "select a.state,a.version,x.actor_subject from crm.appointment a join crm.appointment_transition x on x.tenant_id=a.tenant_id and x.appointment_id=a.appointment_id where a.tenant_id=$1 and a.appointment_id=$2", tenant, id).Scan(&state, &version, &actor); err != nil {
		t.Fatal(err)
	}
	if state != "confirmed" || version != 3 || actor != "operator" {
		t.Fatalf("durable confirmation differs: %s %d %s", state, version, actor)
	}
	for _, sql := range []string{
		"select count(*) from crm.appointment_transition where tenant_id=$1 and appointment_id=$2 and from_state='requested' and to_state='confirmed'",
		"select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='appointment.confirmed'",
		"select count(*) from crm.appointment_resource where tenant_id=$1 and appointment_id=$2",
	} {
		var count int
		if err := pool.QueryRow(ctx, sql, tenant, id).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("expected exactly one durable effect, got %d", count)
		}
	}
	if _, err := pool.Exec(ctx, "delete from crm.appointment_transition where tenant_id=$1 and appointment_id=$2", tenant, id); err == nil {
		t.Fatal("immutable audit accepted deletion")
	}
	readRepo := postgres.NewFranchiseJourney(pool)
	agenda, err := readRepo.AppointmentAgenda(ctx, tenant, "store", start.Add(-time.Hour), end.Add(time.Hour))
	if err != nil || len(agenda.Appointments) != 1 || len(agenda.Resources) != 1 || agenda.Appointments[0].State != "confirmed" || agenda.Resources[0].PrincipalSubject != "" || agenda.Truncated {
		t.Fatalf("agenda snapshot mismatch: %v", err)
	}
	for _, scope := range [][2]string{{tenant, "other"}, {(randomid.Generator{}).New(), "store"}} {
		isolated, err := readRepo.AppointmentAgenda(ctx, scope[0], scope[1], start.Add(-time.Hour), end.Add(time.Hour))
		if err != nil || len(isolated.Appointments) != 0 || len(isolated.Resources) != 0 {
			t.Fatal("agenda scope leaked or query failed")
		}
	}
	if _, err := pool.Exec(ctx, `insert into crm.service_resource(tenant_id,resource_id,organization_id,display_name,resource_kind,status,version) select $1,'overflow-'||n,'store','Synthetic overflow '||n,'service-bay','active',1 from generate_series(1,201) n`, tenant); err != nil {
		t.Fatal(err)
	}
	bounded, err := readRepo.AppointmentAgenda(ctx, tenant, "store", start.Add(-time.Hour), end.Add(time.Hour))
	if err != nil || !bounded.Truncated || len(bounded.Resources) != 200 {
		t.Fatal("agenda silently truncated or exceeded resource budget")
	}
	if browser {
		t.Log("APPOINTMENT_AGENDA_BROWSER_POSTGRES_PASS signature=RS256 state=confirmed version=3 actor=operator audit=1 outbox=1 scoped_read=true bounded_read=true")
	} else {
		t.Log("APPOINTMENT_CONFIRMATION_BFF_POSTGRES_PASS signature=RS256 state=confirmed version=3 actor=operator confirmations=1 audit=1 outbox=1 duplicate=409 customer_scoped=true")
	}
}

func (browserAppointmentClock) Now() time.Time { return time.Now() }

func TestPublicAppointmentBrowserPostgres(t *testing.T) {
	runPublicBrowserPostgres(t, publicBrowserFixture{
		testName:      "public appointment follows captured lead",
		cleanupTables: []string{"crm.appointment", "crm.appointment_slot", "crm.availability_entry", "org.public_location"},
		setup: func(ctx context.Context, t *testing.T, pool *pgxpool.Pool, tenant string) []EnterpriseModule {
			if _, err := pool.Exec(ctx, `insert into org.public_location(tenant_id,organization_id,city,region,country,published)values($1,'browser-org','Fixture City','Fixture Region','AR',true)`, tenant); err != nil {
				t.Fatal(err)
			}
			start := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Minute)
			if _, err := pool.Exec(ctx, `insert into crm.availability_entry(tenant_id,availability_id,organization_id,entry_type,starts_at,ends_at,state,version,created_by_subject)values($1,'browser-working','browser-org','working',$2,$3,'active',1,'synthetic-fixture')`, tenant, start.Add(-time.Hour), start.Add(5*time.Hour)); err != nil {
				t.Fatal(err)
			}
			for i, project := range []string{"chromium-desktop", "chromium-mobile", "firefox-desktop", "webkit-desktop"} {
				at := start.Add(time.Duration(i) * time.Hour)
				if _, err := pool.Exec(ctx, `insert into crm.appointment_slot(tenant_id,slot_id,organization_id,appointment_kind,starts_at,ends_at,capacity,state,version)values($1,$2,'browser-org','consultation',$3,$4,1,'open',1)`, tenant, "browser-slot-"+project, at, at.Add(30*time.Minute)); err != nil {
					t.Fatal(err)
				}
			}
			return []EnterpriseModule{FranchiseJourneyModule{Service: franchisejourney.NewService(postgres.NewFranchiseJourney(pool), randomid.Generator{}, browserAppointmentClock{})}}
		},
		verify: func(ctx context.Context, t *testing.T, pool *pgxpool.Pool, tenant string) {
			var appointments, events, receipts, violations int
			for _, check := range []struct {
				sql    string
				target *int
			}{
				{`select count(*) from crm.appointment a join crm.lead l on l.tenant_id=a.tenant_id and l.lead_id=a.lead_id join crm.appointment_slot s on s.tenant_id=a.tenant_id and s.slot_id=a.slot_id where a.tenant_id=$1 and a.organization_id=l.organization_id and a.model_id=l.model_id and a.starts_at=s.starts_at and a.ends_at=s.ends_at and a.state='requested' and a.version=1`, &appointments},
				{`select count(*) from platform.outbox_event where tenant_id=$1 and event_type='appointment.requested'`, &events},
				{`select count(*) from platform.idempotency_record where tenant_id=$1 and scope='public-appointment' and status='completed'`, &receipts},
				{`select count(*) from (select s.slot_id from crm.appointment_slot s left join crm.appointment a on a.tenant_id=s.tenant_id and a.slot_id=s.slot_id and a.state in ('requested','confirmed') where s.tenant_id=$1 group by s.slot_id,s.capacity having count(a.appointment_id)<>s.capacity) q`, &violations},
			} {
				if err := pool.QueryRow(ctx, check.sql, tenant).Scan(check.target); err != nil {
					t.Fatal(err)
				}
			}
			if appointments != 4 || events != 4 || receipts != 4 || violations != 0 {
				t.Fatalf("appointment invariants: appointments=%d events=%d receipts=%d capacity_violations=%d", appointments, events, receipts, violations)
			}
			t.Log("PUBLIC_APPOINTMENT_BROWSER_POSTGRES_PASS browsers=4 appointments=4 outbox=4 receipts=4 capacity_violations=0")
		},
	})
}

type journeyIDs struct{ n int }

func (i *journeyIDs) New() string { i.n++; return "journey-id" }

type journeyClock struct{}

func (journeyClock) Now() time.Time { return time.Date(2026, 8, 29, 20, 0, 0, 0, time.UTC) }

type journeyVerifier struct {
	principal identity.Principal
	err       error
}

func (v journeyVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return v.principal, v.err
}

type journeyRepo struct {
	readFailure  bool
	appointment  franchisejourney.Appointment
	slot         franchisejourney.AppointmentSlot
	customer     string
	conflict     bool
	resource     franchisejourney.ServiceResource
	resourced    franchisejourney.Appointment
	availability franchisejourney.AvailabilityEntry
	quoteActor   string
	quoteKey     string
	quoteHash    string
	quoteReplay  bool
}

func (r *journeyRepo) PublicLocations(context.Context, string) ([]franchisejourney.Location, error) {
	return []franchisejourney.Location{{Code: "store", Name: "Store"}}, nil
}
func (r *journeyRepo) PublicAppointmentSlots(context.Context, string, string, string, time.Time, time.Time) ([]franchisejourney.AppointmentSlot, error) {
	return []franchisejourney.AppointmentSlot{{ID: "slot", Capacity: 2}}, nil
}
func (r *journeyRepo) CreateAppointmentSlotOnce(ctx context.Context, tenant, subject, key, hash string, v franchisejourney.AppointmentSlot, event string) (franchisejourney.AppointmentSlot, bool, error) {
	value, err := r.CreateAppointmentSlot(ctx, tenant, v, event)
	return value, false, err
}
func (r *journeyRepo) CreateAppointmentSlot(_ context.Context, _ string, value franchisejourney.AppointmentSlot, _ string) (franchisejourney.AppointmentSlot, error) {
	r.slot = value
	return value, nil
}
func (r *journeyRepo) RequestAppointment(_ context.Context, _, _, _ string, v franchisejourney.Appointment, _, _ string) (franchisejourney.Appointment, bool, error) {
	r.appointment = v
	return v, false, nil
}
func (r *journeyRepo) CreateServiceResourceOnce(ctx context.Context, tenant, subject, key, hash string, v franchisejourney.ServiceResource, event string) (franchisejourney.ServiceResource, bool, error) {
	value, err := r.CreateServiceResource(ctx, tenant, v, event)
	return value, false, err
}
func (r *journeyRepo) CreateServiceResource(_ context.Context, _ string, v franchisejourney.ServiceResource, _ string) (franchisejourney.ServiceResource, error) {
	r.resource = v
	return v, nil
}
func (r *journeyRepo) AssignAppointmentResource(_ context.Context, _, _, appointment, resource string, version int64, _ string) (franchisejourney.Appointment, error) {
	r.resourced = franchisejourney.Appointment{ID: appointment, ResourceID: resource, State: "requested", Version: version + 1}
	return r.resourced, nil
}
func (r *journeyRepo) CreateAvailability(_ context.Context, _, _ string, value franchisejourney.AvailabilityEntry, _ string) (franchisejourney.AvailabilityEntry, error) {
	r.availability = value
	return value, nil
}
func (r *journeyRepo) CancelAvailability(_ context.Context, _, _, entry string, version int64, _, _, _ string) (franchisejourney.AvailabilityEntry, error) {
	return franchisejourney.AvailabilityEntry{ID: entry, State: "cancelled", Version: version + 1}, nil
}
func (r *journeyRepo) Availability(context.Context, string, string, string, time.Time, time.Time) ([]franchisejourney.AvailabilityEntry, error) {
	return []franchisejourney.AvailabilityEntry{{ID: "working-window", EntryType: "working", State: "active", Version: 1}}, nil
}
func (r *journeyRepo) TransitionAppointment(_ context.Context, _, _, appointment, _, target string, version int64, _, _, _ string) (franchisejourney.Appointment, error) {
	return franchisejourney.Appointment{ID: appointment, State: target, Version: version + 1}, nil
}
func (r *journeyRepo) CancelCustomerAppointment(_ context.Context, _, _, customer, appointment string, version int64, _, _, _ string) (franchisejourney.Appointment, error) {
	r.customer = customer
	return franchisejourney.Appointment{ID: appointment, State: "cancelled", Version: version + 1}, nil
}
func (r *journeyRepo) Leads(context.Context, string, string, int, string) (franchisejourney.Page[franchisejourney.Lead], error) {
	return franchisejourney.Page[franchisejourney.Lead]{Items: []franchisejourney.Lead{{ID: "lead", Version: 1}}}, nil
}
func (r *journeyRepo) AssignLead(context.Context, string, string, string, string, int64, string) (franchisejourney.Lead, error) {
	if r.conflict {
		return franchisejourney.Lead{}, franchisejourney.ErrConflict
	}
	return franchisejourney.Lead{ID: "lead", Version: 2}, nil
}

func (r *journeyRepo) AssignLeadAs(ctx context.Context, tenant, organization, lead, subject string, version int64, eventID, actor string) (franchisejourney.Lead, error) {
	if actor == "" {
		return franchisejourney.Lead{}, franchisejourney.ErrInvalid
	}
	return r.AssignLead(ctx, tenant, organization, lead, subject, version, eventID)
}
func (r *journeyRepo) TransitionLeadAs(ctx context.Context, tenant, organization, lead, current, target string, version int64, eventID, actor string) (franchisejourney.Lead, error) {
	if actor == "" {
		return franchisejourney.Lead{}, franchisejourney.ErrInvalid
	}
	return r.TransitionLead(ctx, tenant, organization, lead, current, target, version, eventID)
}
func (r *journeyRepo) TransitionLead(context.Context, string, string, string, string, string, int64, string) (franchisejourney.Lead, error) {
	return franchisejourney.Lead{ID: "lead", State: "contacted", Version: 2}, nil
}
func (r *journeyRepo) CreateQuote(_ context.Context, _, key string, v franchisejourney.Quote, hash, _ string) (franchisejourney.Quote, bool, error) {
	r.quoteKey = key
	r.quoteHash = hash
	v.Currency = "ARS"
	v.TotalMinorUnits = 1000
	return v, r.quoteReplay, nil
}
func (r *journeyRepo) AcceptQuote(_ context.Context, _, _, customer, _ string, version int64, _ string, orderID, _, _, _ string) (franchisejourney.Quote, error) {
	r.customer = customer
	return franchisejourney.Quote{State: "accepted", Version: version + 1, OrderID: orderID}, nil
}
func (r *journeyRepo) PublishDeliveryChecklist(_ context.Context, _, subject string, value franchisejourney.DeliveryChecklist, _ string) (franchisejourney.DeliveryChecklist, error) {
	r.customer = subject
	value.State = "published"
	return value, nil
}
func (r *journeyRepo) CompleteDeliveryChecklist(_ context.Context, _, _, subject, _ string, version int64, checklistID string, checklistVersion int64, _ []franchisejourney.ChecklistResponse, _ string) (franchisejourney.Handover, error) {
	r.customer = subject
	return franchisejourney.Handover{State: "presented", Version: version + 1, ChecklistID: checklistID, ChecklistVersion: checklistVersion}, nil
}
func (r *journeyRepo) RejectHandover(_ context.Context, _, organization, customer, handover string, _ int64, reason, details, _ string, exceptionID, _ string) (franchisejourney.DeliveryException, error) {
	r.customer = customer
	return franchisejourney.DeliveryException{ID: exceptionID, OrganizationID: organization, HandoverID: handover, CustomerSubject: customer, ReasonCode: reason, Details: details, State: "open", Version: 1}, nil
}
func (r *journeyRepo) DeliveryExceptions(context.Context, string, string, int) ([]franchisejourney.DeliveryException, error) {
	return []franchisejourney.DeliveryException{{ID: "exception", State: "open", Version: 1}}, nil
}
func (r *journeyRepo) ResolveDeliveryException(_ context.Context, _, _, subject, exceptionID string, version int64, action, _ string, successorID, authorizationID, _, _ string) (franchisejourney.DeliveryResolution, error) {
	r.customer = subject
	result := franchisejourney.DeliveryResolution{Exception: franchisejourney.DeliveryException{ID: exceptionID, State: "resolved", Version: version + 1, ResolutionAction: action}}
	if action == "correct-and-represent" {
		result.SuccessorHandover = &franchisejourney.Handover{ID: successorID, State: "prepared", Version: 1}
	} else {
		result.ReturnAuthorizationID = authorizationID
		result.Disposition = action
	}
	return result, nil
}
func (r *journeyRepo) ReturnCases(context.Context, string, string, int) ([]franchisejourney.ReturnCase, error) {
	return []franchisejourney.ReturnCase{{AuthorizationID: "authorization", AuthorizedAction: "return"}}, nil
}
func (r *journeyRepo) ReceiveReturn(_ context.Context, _, organization, subject, authorizationID, serial, condition, notes, evidence, receiptID, _ string) (franchisejourney.ReturnReceipt, error) {
	r.customer = subject
	return franchisejourney.ReturnReceipt{ID: receiptID, AuthorizationID: authorizationID, OrganizationID: organization, ReceivedSerialNumber: serial, ConditionCode: condition, Notes: notes, EvidenceSHA256: evidence, ReceivedBySubject: subject}, nil
}
func (r *journeyRepo) DecideReturn(_ context.Context, _, _, subject, receiptID, inventoryAction, notes, dispositionID, inventoryID, remedyID, accountingID, _, _ string) (franchisejourney.ReturnDisposition, error) {
	r.customer = subject
	return franchisejourney.ReturnDisposition{ID: dispositionID, ReceiptID: receiptID, InventoryAction: inventoryAction, CustomerRemedy: "refund", Notes: notes, Effects: []franchisejourney.ReturnEffectRequest{{ID: inventoryID}, {ID: remedyID}, {ID: accountingID}}}, nil
}
func (r *journeyRepo) CustomerJourney(_ context.Context, _, _, customer string) (franchisejourney.CustomerJourney, error) {
	r.customer = customer
	if r.readFailure {
		return franchisejourney.CustomerJourney{Handovers: []franchisejourney.Handover{{ID: "PRIVATE-PARTIAL-RESULT"}}}, errors.New("PRIVATE-INTERNAL-DETAIL")
	}
	return franchisejourney.CustomerJourney{}, nil
}

func TestCustomerJourneyReadFailureIsRedactedAndRecoverable(t *testing.T) {
	verifier, token := confirmationTestIssuer(t)
	repo := &journeyRepo{readFailure: true}
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: franchisejourney.NewService(repo, randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
	for _, tc := range []struct {
		name, bearer string
		status       int
	}{
		{"anonymous", "", 401},
		{"wrong-permission", token("customer", "018f4d4a-7b36-7a21-8d10-2f4c54c29b01", []string{"admin:read"}, []string{"store"}), 403},
		{"wrong-organization", token("customer", "018f4d4a-7b36-7a21-8d10-2f4c54c29b01", []string{"customer:self"}, []string{"other"}), 403},
		{"read-failure", token("customer", "018f4d4a-7b36-7a21-8d10-2f4c54c29b01", []string{"customer:self"}, []string{"store"}), 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/v1/customer/journey?organization_id=store", nil)
			if tc.bearer != "" {
				req.Header.Set("Authorization", "Bearer "+tc.bearer)
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != tc.status || strings.Contains(w.Body.String(), "PRIVATE-") {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if tc.status == 500 && !strings.Contains(w.Body.String(), "QUERY_FAILED") {
				t.Fatal("missing controlled error")
			}
		})
	}
	repo.readFailure = false
	req := httptest.NewRequest("GET", "/v1/customer/journey?organization_id=store", nil)
	req.Header.Set("Authorization", "Bearer "+token("customer", "018f4d4a-7b36-7a21-8d10-2f4c54c29b01", []string{"customer:self"}, []string{"store"}))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 || repo.customer != "customer" {
		t.Fatalf("recovery status=%d", w.Code)
	}
}

func (r *journeyRepo) AppointmentAgenda(context.Context, string, string, time.Time, time.Time) (franchisejourney.AppointmentAgenda, error) {
	return franchisejourney.AppointmentAgenda{Appointments: []franchisejourney.Appointment{}, Resources: []franchisejourney.ServiceResource{}}, nil
}

func TestAppointmentAgendaAuthorizationAndRange(t *testing.T) {
	principal := identity.Principal{TenantID: "tenant", Subject: "operator", Permissions: map[string]struct{}{"appointment:manage": {}}, Organizations: map[string]struct{}{"store": {}}}
	valid := "/v1/franchise/agenda?organization_id=store&from=2026-09-06T00:00:00Z&to=2026-09-07T00:00:00Z"
	for _, tc := range []struct {
		name, path string
		principal  identity.Principal
		err        error
		status     int
	}{
		{"authorized", valid, principal, nil, 200},
		{"unauthenticated", valid, principal, identity.ErrUnauthenticated, 401},
		{"wrong-role", valid, identity.Principal{TenantID: "tenant", Organizations: principal.Organizations}, nil, 403},
		{"wrong-organization", strings.Replace(valid, "organization_id=store", "organization_id=other", 1), principal, nil, 403},
		{"missing-dates", "/v1/franchise/agenda?organization_id=store", principal, nil, 400},
		{"unbounded", strings.Replace(valid, "2026-09-07", "2026-12-07", 1), principal, nil, 400},
		{"reversed", strings.Replace(valid, "2026-09-07", "2026-09-05", 1), principal, nil, 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			journeyHandler(&journeyRepo{}, tc.principal, tc.err).ServeHTTP(w, request("GET", tc.path, ""))
			if w.Code != tc.status {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if tc.status == 200 && (w.Header().Get("Cache-Control") != "no-store" || !strings.Contains(w.Body.String(), "\"appointments\":[]")) {
				t.Fatal("read contract/cache mismatch")
			}
		})
	}
}
func (r *journeyRepo) AcceptHandover(_ context.Context, _, _, customer, _ string, _ int64, _, checklistID string, checklistVersion int64, evidence, _ string) (franchisejourney.Handover, error) {
	r.customer = customer
	return franchisejourney.Handover{State: "accepted", AcceptanceEvidence: evidence, ChecklistID: checklistID, ChecklistVersion: checklistVersion}, nil
}

func journeyHandler(repository *journeyRepo, principal identity.Principal, verifyErr error) http.Handler {
	service := franchisejourney.NewService(repository, &journeyIDs{}, journeyClock{})
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: service}.Register(mux, journeyVerifier{principal: principal, err: verifyErr})
	return mux
}
func request(method, path, body string) *http.Request {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		r.Header.Set("Content-Type", "application/json")
	}
	r.Header.Set("Authorization", "Bearer token")
	return r
}

func TestPublicLocationsAndAppointment(t *testing.T) {
	repository := &journeyRepo{}
	handler := journeyHandler(repository, identity.Principal{}, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("GET", "/v1/public/tenant/locations", ""))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "Store") {
		t.Fatalf("locations status=%d body=%s", w.Code, w.Body.String())
	}
	body := `{"lead_id":"lead","model_id":"model","kind":"test-drive","starts_at":"2026-08-29T21:00:00Z"}`
	r := request("POST", "/v1/public/tenant/store/appointments", body)
	r.Header.Set("Idempotency-Key", "appointment-key-1")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != 202 || repository.appointment.State != "requested" {
		t.Fatalf("appointment status=%d body=%s value=%+v", w.Code, w.Body.String(), repository.appointment)
	}
}

func TestCreateQuoteCarriesIdempotencyAndReportsReplay(t *testing.T) {
	repository := &journeyRepo{}
	principal := identity.Principal{Subject: "writer", TenantID: "tenant", Permissions: map[string]struct{}{"quote:write": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	body := `{"organization_id":"store","lead_id":"lead","variant_id":"variant","price_book_id":"retail","valid_until":"2026-08-30T20:00:00Z"}`
	r := request("POST", "/v1/franchise/quotes", body)
	r.Header.Set("Idempotency-Key", "quote-request-0001")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != 201 || repository.quoteActor != "writer" || repository.quoteKey != "quote-request-0001" || len(repository.quoteHash) != 64 {
		t.Fatalf("quote status=%d key=%q hash=%q body=%s", w.Code, repository.quoteKey, repository.quoteHash, w.Body.String())
	}
	noKey := httptest.NewRecorder()
	handler.ServeHTTP(noKey, request("POST", "/v1/franchise/quotes", body))
	if noKey.Code != 400 {
		t.Fatalf("missing key status=%d", noKey.Code)
	}
	repository.quoteReplay = true
	r = request("POST", "/v1/franchise/quotes", body)
	r.Header.Set("Idempotency-Key", "quote-request-0001")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("quote replay status=%d body=%s", w.Code, w.Body.String())
	}
	r = request("POST", "/v1/franchise/quotes", body)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != 400 {
		t.Fatalf("quote without key status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestPublicAndProtectedAppointmentCapacity(t *testing.T) {
	repository := &journeyRepo{}
	principal := identity.Principal{Subject: "slot-manager", TenantID: "tenant", Permissions: map[string]struct{}{"appointment:manage": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("GET", "/v1/public/tenant/store/appointment-slots?kind=service&from=2026-08-29T20:01:00Z&to=2026-08-30T20:00:00Z", ""))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "slot") {
		t.Fatalf("public slots status=%d body=%s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	slotRequest := request("POST", "/v1/franchise/appointment-slots", `{"organization_id":"store","kind":"service","starts_at":"2026-08-29T21:00:00Z","ends_at":"2026-08-29T22:00:00Z","capacity":2}`)
	slotRequest.Header.Set("Idempotency-Key", "slot-unit-fixture-key")
	handler.ServeHTTP(w, slotRequest)
	if w.Code != 201 || repository.slot.Capacity != 2 || repository.slot.State != "open" {
		t.Fatalf("create slot status=%d body=%s slot=%+v", w.Code, w.Body.String(), repository.slot)
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/appointment-slots", `{"organization_id":"store","kind":"service","starts_at":"2026-08-29T21:00:00Z","ends_at":"2026-08-29T22:00:00Z","capacity":2,"state":"open"}`))
	if w.Code != 400 {
		t.Fatalf("server-managed slot field accepted status=%d", w.Code)
	}
}

func TestProtectedAppointmentResourceOperations(t *testing.T) {
	repository := &journeyRepo{}
	principal := identity.Principal{Subject: "resource-manager", TenantID: "tenant", Permissions: map[string]struct{}{"resource:manage": {}, "appointment:manage": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	w := httptest.NewRecorder()
	resourceRequest := request("POST", "/v1/franchise/resources", `{"organization_id":"store","principal_subject":"technician","display_name":"Technician","kind":"employee","skills":["service"]}`)
	resourceRequest.Header.Set("Idempotency-Key", "resource-unit-fixture-key")
	handler.ServeHTTP(w, resourceRequest)
	if w.Code != 201 || repository.resource.Status != "active" || repository.resource.ID == "" {
		t.Fatalf("resource status=%d body=%s value=%+v", w.Code, w.Body.String(), repository.resource)
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/appointments/appointment/resources", `{"organization_id":"store","resource_id":"resource","version":1}`))
	if w.Code != 200 || repository.resourced.ResourceID != "resource" {
		t.Fatalf("assignment status=%d body=%s value=%+v", w.Code, w.Body.String(), repository.resourced)
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/resources", `{"organization_id":"store","display_name":"Bay","kind":"service-bay","skills":["service"],"status":"active"}`))
	if w.Code != 400 {
		t.Fatalf("server-managed resource state accepted status=%d", w.Code)
	}
}

func TestProtectedScopeAndCustomerSubject(t *testing.T) {
	repository := &journeyRepo{}
	principal := identity.Principal{Subject: "customer-1", TenantID: "tenant", Permissions: map[string]struct{}{"customer:self": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("GET", "/v1/customer/journey?organization_id=other", ""))
	if w.Code != 403 {
		t.Fatalf("cross organization status=%d", w.Code)
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("GET", "/v1/customer/journey?organization_id=store", ""))
	if w.Code != 200 || repository.customer != "customer-1" {
		t.Fatalf("customer status=%d bound=%s", w.Code, repository.customer)
	}
}

func TestAvailabilityAndCustomerCancellationBindActorAndScope(t *testing.T) {
	repository := &journeyRepo{}
	principal := identity.Principal{Subject: "scheduler", TenantID: "tenant", Permissions: map[string]struct{}{"availability:manage": {}, "availability:read": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	w := httptest.NewRecorder()
	creationRequest := request("POST", "/v1/franchise/availability", `{"organization_id":"store","resource_id":"technician","entry_type":"unavailable","reason_code":"annual-leave","starts_at":"2026-08-30T21:00:00Z","ends_at":"2026-08-30T22:00:00Z"}`)
	creationRequest.Header.Set("Idempotency-Key", "availability-test-key")
	handler.ServeHTTP(w, creationRequest)
	if w.Code != 201 || repository.availability.State != "active" || repository.availability.Version != 1 {
		t.Fatalf("create availability status=%d body=%s value=%+v", w.Code, w.Body.String(), repository.availability)
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("GET", "/v1/franchise/availability?organization_id=store&resource_id=technician&from=2026-08-30T20:00:00Z&to=2026-08-31T20:00:00Z", ""))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "working-window") {
		t.Fatalf("availability status=%d body=%s", w.Code, w.Body.String())
	}

	principal = identity.Principal{Subject: "customer-1", TenantID: "tenant", Permissions: map[string]struct{}{"customer:self": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler = journeyHandler(repository, principal, nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/customer/appointments/appointment/cancel", `{"organization_id":"store","version":2,"reason_code":"customer-request"}`))
	if w.Code != 200 || repository.customer != "customer-1" || !strings.Contains(w.Body.String(), `"state":"cancelled"`) {
		t.Fatalf("customer cancel status=%d customer=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
}

func TestLeadPermissionsConflictAndStrictJSON(t *testing.T) {
	repository := &journeyRepo{conflict: true}
	principal := identity.Principal{Subject: "sales", TenantID: "tenant", Permissions: map[string]struct{}{"lead:assign": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/leads/lead/assign", `{"organization_id":"store","assigned_subject":"sales","version":1}`))
	if w.Code != 409 {
		t.Fatalf("conflict status=%d body=%s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/leads/lead/assign", `{"organization_id":"store","assigned_subject":"sales","version":1,"unknown":true}`))
	if w.Code != 400 {
		t.Fatalf("strict json status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthenticationFailure(t *testing.T) {
	handler := journeyHandler(&journeyRepo{}, identity.Principal{}, errors.New("bad token"))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("GET", "/v1/franchise/leads?organization_id=store", ""))
	if w.Code != 401 {
		t.Fatalf("status=%d", w.Code)
	}
}

func TestAcceptQuoteBindsIdentityAndRejectsCrossOrganization(t *testing.T) {
	repository := &journeyRepo{}
	principal := identity.Principal{Subject: "customer-1", TenantID: "tenant", Permissions: map[string]struct{}{"customer:self": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/customer/quotes/quote/accept", `{"organization_id":"other","version":1}`))
	if w.Code != 403 {
		t.Fatalf("cross organization status=%d body=%s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/customer/quotes/quote/accept", `{"organization_id":"store","version":1}`))
	if w.Code != 200 || repository.customer != "customer-1" || !strings.Contains(w.Body.String(), `"state":"accepted"`) {
		t.Fatalf("accept status=%d customer=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
}

func TestAcceptHandoverHashesExactCommandAndBindsIdentity(t *testing.T) {
	repository := &journeyRepo{}
	principal := identity.Principal{Subject: "customer-1", TenantID: "tenant", Permissions: map[string]struct{}{"customer:self": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/customer/handovers/handover/accept", `{"organization_id":"store","version":2,"confirmed_received":true,"serial_number":"SERIAL","checklist_id":"standard-delivery","checklist_version":1}`))
	if w.Code != 200 || repository.customer != "customer-1" || !strings.Contains(w.Body.String(), `"state":"accepted"`) {
		t.Fatalf("accept status=%d customer=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/customer/handovers/handover/accept", `{"organization_id":"store","version":2,"confirmed_received":false,"serial_number":"SERIAL","checklist_id":"standard-delivery","checklist_version":1}`))
	if w.Code != 400 {
		t.Fatalf("unconfirmed handover accepted status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestDeliveryChecklistPublicationAndCompletionBindOperator(t *testing.T) {
	repository := &journeyRepo{}
	principal := identity.Principal{Subject: "operator-1", TenantID: "tenant", Permissions: map[string]struct{}{"handover:manage": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/delivery-checklists", `{"organization_id":"store","checklist_id":"standard-delivery","version":1,"title":"Entrega estándar","items":[{"id":"serial-observed","prompt":"Verificar serie","response_type":"serial","required":true}]}`))
	if w.Code != 201 || repository.customer != "operator-1" || !strings.Contains(w.Body.String(), `"state":"published"`) {
		t.Fatalf("publish status=%d subject=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/handovers/handover/complete-checklist", `{"organization_id":"store","version":1,"checklist_id":"standard-delivery","checklist_version":1,"responses":[{"item_id":"serial-observed","response_text":"SERIAL"}]}`))
	if w.Code != 200 || repository.customer != "operator-1" || !strings.Contains(w.Body.String(), `"state":"presented"`) {
		t.Fatalf("complete status=%d subject=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
}

func TestDeliveryRejectionAndResolutionBindVerifiedActors(t *testing.T) {
	repository := &journeyRepo{}
	customer := identity.Principal{Subject: "customer-1", TenantID: "tenant", Permissions: map[string]struct{}{"customer:self": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, customer, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/customer/handovers/handover/reject", `{"organization_id":"store","version":2,"reason_code":"visible-damage","details":"Rayón visible"}`))
	if w.Code != 200 || repository.customer != "customer-1" || !strings.Contains(w.Body.String(), `"state":"open"`) {
		t.Fatalf("reject status=%d subject=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
	repository = &journeyRepo{}
	operator := identity.Principal{Subject: "operator-1", TenantID: "tenant", Permissions: map[string]struct{}{"handover:manage": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler = journeyHandler(repository, operator, nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/delivery-exceptions/exception/resolve", `{"organization_id":"store","version":1,"action":"correct-and-represent","notes":"Corregir preparación"}`))
	if w.Code != 200 || repository.customer != "operator-1" || !strings.Contains(w.Body.String(), `"successor_handover"`) {
		t.Fatalf("resolve status=%d subject=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("GET", "/v1/franchise/delivery-exceptions?organization_id=store", ""))
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"exception"`) {
		t.Fatalf("list status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestReturnReceiptAndDispositionBindOperatorAndClosedFields(t *testing.T) {
	repository := &journeyRepo{}
	operator := identity.Principal{Subject: "operator-1", TenantID: "tenant", Permissions: map[string]struct{}{"handover:manage": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, operator, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("GET", "/v1/franchise/returns?organization_id=store", ""))
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"authorization"`) {
		t.Fatalf("return list status=%d body=%s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/return-authorizations/authorization/receive", `{"organization_id":"store","serial_number":"SERIAL","condition_code":"damaged","notes":"Daño confirmado","evidence_sha256":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`))
	if w.Code != 200 || repository.customer != "operator-1" || !strings.Contains(w.Body.String(), `"condition_code":"damaged"`) {
		t.Fatalf("return receive status=%d subject=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/return-receipts/receipt/decide", `{"organization_id":"store","inventory_action":"quarantine","notes":"Separar hasta reconciliar"}`))
	if w.Code != 200 || repository.customer != "operator-1" || !strings.Contains(w.Body.String(), `"customer_remedy":"refund"`) {
		t.Fatalf("return decision status=%d subject=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/return-receipts/receipt/decide", `{"organization_id":"store","inventory_action":"restock","notes":"No","refund_now":true}`))
	if w.Code != 400 {
		t.Fatalf("return decision accepted browser-owned effect status=%d body=%s", w.Code, w.Body.String())
	}
}

// TestQuoteAcceptanceBrowserPostgres is an opt-in AUTHORED connected fixture.
// Reuses the real journey API, PostgreSQL repository and RS256/JWKS issuer.
// No live login, real customer, provider, payment or business acceptance claim.
func TestQuoteAcceptanceBrowserPostgres(t *testing.T) {
	if os.Getenv("ELITE_QUOTE_E2E") != "1" {
		t.Skip("explicit disposable quote gate not requested")
	}
	for _, project := range []string{"chromium-desktop", "chromium-mobile", "firefox-desktop", "webkit-desktop"} {
		if selected := os.Getenv("ELITE_QUOTE_PROJECT"); selected != "" && selected != project {
			continue
		}
		t.Run(project, func(t *testing.T) { runQuoteAcceptanceBrowser(t, project) })
	}
}

func runQuoteAcceptanceBrowser(t *testing.T, project string) {
	budget := 3 * time.Minute
	if os.Getenv("ELITE_PAYMENT_E2E") == "1" {
		budget = 5 * time.Minute
	}
	ctx, cancel := context.WithTimeout(context.Background(), budget)
	defer cancel()
	config, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil || config.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(config.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback elite_confirmation_* database")
	}
	if os.Getenv("ELITE_RETURN_MULTITAB_E2E") == "1" && (os.Getenv("ELITE_RETURN_RECOVERY_E2E") != "1" || os.Getenv("ELITE_DELIVERY_READ_E2E") != "1" || os.Getenv("ELITE_ORDER_E2E") != "1") {
		t.Fatal("return multi-tab gate requires return, order and delivery-read fixtures")
	}
	if os.Getenv("ELITE_RETURN_SESSION_E2E") == "1" && os.Getenv("ELITE_RETURN_MULTITAB_E2E") != "1" {
		t.Fatal("return session gate requires multi-tab fixture")
	}
	web := os.Getenv("ELITE_WEB_ROOT")
	if !filepath.IsAbs(web) {
		t.Fatal("absolute materialized web root required")
	}
	artifacts, err := os.MkdirTemp(web, "quote-browser-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	tenant := (randomid.Generator{}).New()
	fixtures := []string{
		"insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'quote-'||replace(($1::uuid)::text,'-',''),'Synthetic quote','Synthetic quote')",
		"insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic Store','store')",
		"insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic Customer','quote@example.invalid')",
		"insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','customer','new','fixture','{}')",
		"insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic bicycle','bicycle','active')",
		"insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic variant','{}','active')",
		"insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,status)values($1,'book','AR','ARS',clock_timestamp()-interval '1 day','active')",
		"insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)values($1,'book','variant',250000,'inclusive')",
		"insert into sales.quotation(tenant_id,quotation_id,organization_id,lead_id,customer_principal_id,variant_id,price_book_id,currency,total_minor_units,valid_until,state,version)select $1,q,'store','lead','customer','variant','book','ARS',250000,clock_timestamp()+interval '1 day','issued',1 from unnest(array['quote-success','quote-lost','quote-race','quote-denied'])q",
		"insert into sales.quotation(tenant_id,quotation_id,organization_id,lead_id,customer_principal_id,variant_id,price_book_id,currency,total_minor_units,created_at,valid_until,state,version)values($1,'quote-expired','store','lead','customer','variant','book','ARS',250000,clock_timestamp()-interval '2 days',clock_timestamp()-interval '1 day','issued',1)",
	}
	for _, sql := range fixtures {
		if _, err = pool.Exec(ctx, sql, tenant); err != nil {
			pool.Close()
			t.Fatal(err)
		}
	}
	pool.Close()
	verifier, token := confirmationTestIssuer(t)
	phases := []string{"initial", "read-failure", "restart"}
	if os.Getenv("ELITE_ORDER_E2E") == "1" {
		phases = append(phases, "operation", "operation-read-failure", "operation-restart")
	}
	if os.Getenv("ELITE_PAYMENT_E2E") == "1" {
		phases = append(phases, "payment", "payment-disabled", "payment-restart")
	}
	if os.Getenv("ELITE_DELIVERY_READ_E2E") == "1" {
		if os.Getenv("ELITE_ORDER_E2E") != "1" {
			t.Fatal("delivery read gate requires connected order/stock gate")
		}
		phases = append(phases, "delivery-read-corrupt", "delivery-read-restored")
	}
	if os.Getenv("ELITE_DELIVERY_ACTION_E2E") == "1" {
		if os.Getenv("ELITE_DELIVERY_READ_E2E") != "1" {
			t.Fatal("delivery action gate requires delivery read fixture")
		}
		phases = append(phases, "delivery-action-loss")
	}
	if os.Getenv("ELITE_OPERATOR_SECTIONS_E2E") == "1" {
		if os.Getenv("ELITE_ORDER_E2E") != "1" {
			t.Fatal("operator sections require connected order gate")
		}
		phases = append(phases, "operation-sections")
	}
	if os.Getenv("ELITE_RESOLUTION_RECOVERY_E2E") == "1" {
		if os.Getenv("ELITE_DELIVERY_ACTION_E2E") != "1" {
			t.Fatal("resolution recovery requires actual rejected delivery fixture")
		}
		phases = append(phases, "operation-resolution")
	}
	if os.Getenv("ELITE_LEAD_RECOVERY_E2E") == "1" {
		phases = append(phases, "operation-lead")
	}
	if os.Getenv("ELITE_AVAILABILITY_RECOVERY_E2E") == "1" {
		phases = append(phases, "operation-availability")
	}
	if os.Getenv("ELITE_QUOTE_CREATE_E2E") == "1" {
		phases = append(phases, "operation-quote")
	}
	if os.Getenv("ELITE_RETURN_RECOVERY_E2E") == "1" {
		phases = append(phases, "operation-returns")
	}
	if os.Getenv("ELITE_CHECKLIST_COMPLETE_E2E") == "1" {
		phases = append(phases, "operation-complete-checklist")
	}
	if os.Getenv("ELITE_CHECKLIST_PUBLISH_E2E") == "1" {
		phases = append(phases, "operation-publish-checklist")
	}
	if os.Getenv("ELITE_SLOT_CREATE_E2E") == "1" {
		phases = append(phases, "operation-create-slot")
	}
	if os.Getenv("ELITE_RESOURCE_CREATE_E2E") == "1" {
		phases = append(phases, "operation-create-resource")
	}
	if os.Getenv("ELITE_AVAILABILITY_CREATE_E2E") == "1" {
		phases = append(phases, "operation-create-availability")
	}
	for _, phase := range phases {
		func() {
			// New pool, repository, API process fixture and Next process each phase.
			p, err := pgxpool.NewWithConfig(ctx, config.Copy())
			if err != nil {
				t.Fatal(err)
			}
			defer p.Close()
			if phase == "operation-quote" {
				for i := 0; i < 3; i++ {
					if _, err := p.Exec(ctx, `insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,lifecycle_state,source_code,contact_payload)values($1,$2,'store','customer','new','fixture','{}')`, tenant, fmt.Sprintf("quote-create-%d", i)); err != nil {
						t.Fatal(err)
					}
				}
			}
			if phase == "operation-availability" {
				repo := postgres.NewFranchiseJourney(p)
				if _, err := p.Exec(ctx, `insert into org.public_location(tenant_id,organization_id,city,region,country,published)values($1,'store','Synthetic','Synthetic','AR',true)`, tenant); err != nil {
					t.Fatal(err)
				}
				for i := 0; i < 5; i++ {
					start := time.Now().UTC().Add(time.Duration(72+i*24) * time.Hour).Truncate(time.Second)
					entryType, reason := "working", ""
					if i%2 == 1 {
						entryType, reason = "unavailable", "synthetic-absence"
					}
					if _, err := repo.CreateAvailability(ctx, tenant, "fixture", franchisejourney.AvailabilityEntry{ID: fmt.Sprintf("cancel-window-%d", i), OrganizationID: "store", EntryType: entryType, ReasonCode: reason, StartsAt: start, EndsAt: start.Add(time.Hour)}, randomid.Generator{}.New()); err != nil {
						t.Fatal(err)
					}
					if i == 4 {
						if _, err := p.Exec(ctx, `insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,appointment_kind,starts_at,state,version)values($1,'blocked-window-appointment','store','lead','consultation',$2,'requested',1)`, tenant, start.Add(15*time.Minute)); err != nil {
							t.Fatal(err)
						}
					}
					if i%2 == 0 {
						if _, err := repo.CreateAppointmentSlot(ctx, tenant, franchisejourney.AppointmentSlot{ID: fmt.Sprintf("window-slot-%d", i), OrganizationID: "store", Kind: "consultation", StartsAt: start, EndsAt: start.Add(time.Hour), Capacity: 2}, randomid.Generator{}.New()); err != nil {
							t.Fatal(err)
						}
					}
				}
			}
			if phase == "operation-lead" {
				if _, err := p.Exec(ctx, `insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,lifecycle_state,source_code,contact_payload) values($1,'operator-lead','store','customer','new','fixture','{}')`, tenant); err != nil {
					t.Fatal(err)
				}
			}
			if phase == "operation" {
				if _, err := p.Exec(ctx, `insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) select $1,s,'store','variant',s,s,s,'available',1,clock_timestamp() from unnest(array['stock-lost','stock-race'])s`, tenant); err != nil {
					t.Fatal(err)
				}
			}
			deliverySerial := ""
			mux := http.NewServeMux()
			if phase == "delivery-read-corrupt" {
				for _, q := range []string{
					`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'foreign-delivery','foreign-delivery','Synthetic foreign store','store')`,
					`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) values($1,'PRIVATE-FOREIGN-STOCK','foreign-delivery','variant','PRIVATE-FOREIGN-STOCK','PRIVATE-FOREIGN-STOCK','PRIVATE-FOREIGN-STOCK','available',1,clock_timestamp())`,
					`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) select $1,'read-fixture-handover','store',o.order_id,o.customer_principal_id,'PRIVATE-FOREIGN-STOCK','prepared',1 from sales.customer_order o join sales.customer_order_line l on l.tenant_id=o.tenant_id and l.order_id=o.order_id where o.tenant_id=$1 and o.organization_id='store' and l.allocated_stock_unit_id is not null order by o.order_id limit 1`,
				} {
					if tag, err := p.Exec(ctx, q, tenant); err != nil || tag.RowsAffected() != 1 {
						t.Fatalf("delivery read fixture failed: %v", err)
					}
				}
			}
			if phase == "delivery-read-restored" {
				if tag, err := p.Exec(ctx, `update sales.delivery_handover h set stock_unit_id=l.allocated_stock_unit_id from sales.customer_order_line l where h.tenant_id=$1 and h.handover_id='read-fixture-handover' and l.tenant_id=h.tenant_id and l.order_id=h.order_id and l.allocated_stock_unit_id is not null`, tenant); err != nil || tag.RowsAffected() != 1 {
					t.Fatalf("restore synthetic relation: %v", err)
				}
			}
			if phase == "delivery-action-loss" {
				if err := p.QueryRow(ctx, `select i.serial_number from sales.delivery_handover h join inventory.stock_unit i on i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id where h.tenant_id=$1 and h.handover_id='read-fixture-handover'`, tenant).Scan(&deliverySerial); err != nil {
					t.Fatal(err)
				}
				if _, err := p.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) select tenant_id,'reject-fixture-handover',organization_id,order_id,customer_principal_id,stock_unit_id,'prepared',1 from sales.delivery_handover where tenant_id=$1 and handover_id='read-fixture-handover'`, tenant); err != nil {
					t.Fatal(err)
				}
				repo := postgres.NewFranchiseJourney(p)
				checklist, err := repo.PublishDeliveryChecklist(ctx, tenant, "fixture-operator", franchisejourney.DeliveryChecklist{ID: "delivery-ui-test", OrganizationID: "store", Version: 1, Title: "Synthetic preparation", State: "published", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Verify serial", ResponseType: "serial", Required: true}}}, randomid.Generator{}.New())
				if err != nil {
					t.Fatal(err)
				}
				for _, id := range []string{"read-fixture-handover", "reject-fixture-handover"} {
					if _, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "fixture-operator", id, 1, checklist.ID, 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: deliverySerial}}, randomid.Generator{}.New()); err != nil {
						t.Fatal(err)
					}
				}
			}
			if phase == "operation-resolution" {
				repo := postgres.NewFranchiseJourney(p)
				for _, id := range []string{"return-fixture-handover", "exchange-fixture-handover"} {
					if _, err := p.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) select tenant_id,$2,organization_id,order_id,customer_principal_id,stock_unit_id,'prepared',1 from sales.delivery_handover where tenant_id=$1 and handover_id='reject-fixture-handover'`, tenant, id); err != nil {
						t.Fatal(err)
					}
					var serial string
					if err := p.QueryRow(ctx, `select i.serial_number from sales.delivery_handover h join inventory.stock_unit i on i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id where h.tenant_id=$1 and h.handover_id=$2`, tenant, id).Scan(&serial); err != nil {
						t.Fatal(err)
					}
					if _, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "fixture-operator", id, 1, "delivery-ui-test", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: serial}}, randomid.Generator{}.New()); err != nil {
						t.Fatal(err)
					}
					if _, err := repo.RejectHandover(ctx, tenant, "store", "customer", id, 2, "asset-condition", "Synthetic fixture rejection", strings.Repeat("a", 64), randomid.Generator{}.New(), randomid.Generator{}.New()); err != nil {
						t.Fatal(err)
					}
				}
			}
			FranchiseJourneyModule{Service: franchisejourney.NewService(postgres.NewFranchiseJourney(p), randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
			paymentProvider := ""
			if strings.HasPrefix(phase, "payment") && phase != "payment-disabled" {
				paymentProvider = "stripe"
			}
			CommerceModule{Service: commerce.NewService(postgres.NewCommerce(p), randomid.Generator{}), PaymentProvider: paymentProvider}.Register(mux, verifier)
			var returnWritten atomic.Bool
			var returnReadFailures atomic.Int32
			var exceptionReadFailures atomic.Int32
			var resolutionWritten atomic.Bool
			var leadWritten atomic.Bool
			var availabilityWritten atomic.Bool
			var availabilityReadFailures atomic.Int32
			var leadReadFailures atomic.Int32
			var resolutionReadFailures atomic.Int32
			api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if phase == "operation-returns" && r.Method == http.MethodPost && r.URL.Path == "/v1/franchise/return-authorizations/ui-return-loss/receive" {
					defer returnWritten.Store(true)
				}
				if phase == "operation-returns" && r.Method == http.MethodGet && r.URL.Path == "/v1/franchise/returns" && returnWritten.Load() && returnReadFailures.Add(1) == 1 {
					writeProblem(w, 503, "INJECTED_RETURN_LIST_FAILURE", "PRIVATE-SYNTHETIC-RETURN-DETAIL")
					return
				}

				if phase == "operation-availability" && r.URL.Path == "/v1/franchise/availability" && r.Method == http.MethodGet && availabilityWritten.Load() && availabilityReadFailures.Add(1) == 1 {
					writeProblem(w, 503, "INJECTED_AVAILABILITY_QUERY_FAILURE", "PRIVATE-SYNTHETIC-AVAILABILITY-DETAIL")
					return
				}
				if phase == "operation-availability" && r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/availability/cancel-window-") {
					defer availabilityWritten.Store(true)
				}
				if phase == "operation-lead" && r.URL.Path == "/v1/franchise/leads" && leadWritten.Load() && leadReadFailures.Add(1) == 1 {
					writeProblem(w, 503, "INJECTED_LEAD_QUERY_FAILURE", "PRIVATE-SYNTHETIC-LEAD-DETAIL")
					return
				}
				if phase == "operation-lead" && r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/leads/operator-lead/") {
					defer leadWritten.Store(true)
				}
				if phase == "operation-resolution" && r.URL.Path == "/v1/franchise/delivery-exceptions" && resolutionWritten.Load() && resolutionReadFailures.Add(1) == 1 {
					writeProblem(w, 503, "INJECTED_RESOLUTION_QUERY_FAILURE", "PRIVATE-SYNTHETIC-RESOLUTION-DETAIL")
					return
				}
				if phase == "operation-resolution" && r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/resolve") {
					defer resolutionWritten.Store(true)
				}
				if phase == "operation-sections" && r.URL.Path == "/v1/franchise/delivery-exceptions" && exceptionReadFailures.Add(1) == 1 {
					writeProblem(w, 503, "INJECTED_SECTION_UNAVAILABLE", "PRIVATE-SYNTHETIC-SECTION-DETAIL")
					return
				}
				if phase == "operation-read-failure" && r.URL.Path == "/v1/commerce/orders" {
					writeProblem(w, 503, "INJECTED_READ_UNAVAILABLE", "synthetic outage")
					return
				}
				if phase == "read-failure" && r.URL.Path == "/v1/customer/journey" {
					w.Header().Set("Content-Type", "application/problem+json")
					w.WriteHeader(http.StatusServiceUnavailable)
					_, _ = w.Write([]byte(`{"code":"INJECTED_READ_UNAVAILABLE"}`))
					return
				}
				mux.ServeHTTP(w, r)
			}))
			defer api.Close()
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			address := listener.Addr().String()
			listener.Close()
			target, _ := url.Parse("http://" + address)
			edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
			defer edge.Close()
			env := []string{}
			for _, item := range os.Environ() {
				name := strings.ToUpper(strings.SplitN(item, "=", 2)[0])
				if name == "TEST_DATABASE_URL" || name == "DATABASE_URL" || strings.HasPrefix(name, "ELITE_") || name == "ENTERPRISE_API_BASE_URL" || name == "AUTH_SESSION_SECRET" || name == "APP_BASE_URL" {
					continue
				}
				env = append(env, item)
			}

			if phase == "operation-returns" {
				if e := p.QueryRow(ctx, `select i.serial_number from sales.delivery_handover h join inventory.stock_unit i on i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id where h.tenant_id=$1 and h.handover_id='read-fixture-handover'`, tenant).Scan(&deliverySerial); e != nil {
					t.Fatal(e)
				}
				repo := postgres.NewFranchiseJourney(p)
				checklist := franchisejourney.DeliveryChecklist{ID: "ui-return-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic returns", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}}}
				if _, e := repo.PublishDeliveryChecklist(ctx, tenant, "fixture-author", checklist, randomid.Generator{}.New()); e != nil {
					t.Fatal(e)
				}
				for i, id := range []string{"ui-return-loss", "ui-exchange-invalid", "ui-return-race"} {
					if _, e := p.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) select tenant_id,$2,organization_id,order_id,customer_principal_id,stock_unit_id,'prepared',1 from sales.delivery_handover where tenant_id=$1 and handover_id='read-fixture-handover'`, tenant, id); e != nil {
						t.Fatal(e)
					}
					if _, e := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "fixture", id, 1, checklist.ID, 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: deliverySerial}}, randomid.Generator{}.New()); e != nil {
						t.Fatal(e)
					}
					x, e := repo.RejectHandover(ctx, tenant, "store", "customer", id, 2, "serial-mismatch", "Synthetic", strings.Repeat("a", 64), randomid.Generator{}.New(), randomid.Generator{}.New())
					if e != nil {
						t.Fatal(e)
					}
					action := "return"
					if i == 1 {
						action = "exchange"
					}
					if _, e = repo.ResolveDeliveryException(ctx, tenant, "store", "fixture", x.ID, 1, action, "Synthetic", "", id, randomid.Generator{}.New(), randomid.Generator{}.New()); e != nil {
						t.Fatal(e)
					}
				}
			}
			if phase == "operation-complete-checklist" {
				if e := p.QueryRow(ctx, `select i.serial_number from sales.delivery_handover h join inventory.stock_unit i on i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id and i.organization_id=h.organization_id where h.tenant_id=$1 and h.handover_id='read-fixture-handover'`, tenant).Scan(&deliverySerial); e != nil {
					t.Fatal(e)
				}
				repo := postgres.NewFranchiseJourney(p)
				if _, e := repo.PublishDeliveryChecklist(ctx, tenant, "fixture-author", franchisejourney.DeliveryChecklist{ID: "ui-completion-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic completion", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}, {ID: "confirmed", Ordinal: 2, Prompt: "Confirmación", ResponseType: "confirmation", Required: true}}}, randomid.Generator{}.New()); e != nil {
					t.Fatal(e)
				}
				for _, id := range []string{"ui-completion-loss", "ui-completion-invalid", "ui-completion-race"} {
					if _, e := p.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) select tenant_id,$2,organization_id,order_id,customer_principal_id,stock_unit_id,'prepared',1 from sales.delivery_handover where tenant_id=$1 and handover_id='read-fixture-handover'`, tenant, id); e != nil {
						t.Fatal(e)
					}
				}
			}
			if phase == "operation-create-slot" {
				for _, days := range []int{30, 31, 200} {
					start := time.Now().UTC().Add(time.Duration(days) * 24 * time.Hour)
					if _, e := postgres.NewFranchiseJourney(p).CreateAvailability(ctx, tenant, "planner", franchisejourney.AvailabilityEntry{ID: fmt.Sprintf("slot-create-working-%d", days), OrganizationID: "store", EntryType: "working", StartsAt: start.Add(-time.Hour), EndsAt: start.Add(2 * time.Hour)}, randomid.Generator{}.New()); e != nil {
						t.Fatal(e)
					}
				}
			}
			env = append(env, "ELITE_RETURN_SESSION_E2E="+os.Getenv("ELITE_RETURN_SESSION_E2E"), "ELITE_RETURN_MULTITAB_E2E="+os.Getenv("ELITE_RETURN_MULTITAB_E2E"), "ELITE_DELIVERY_SERIAL="+deliverySerial, "ELITE_QUOTE_E2E=1", "ELITE_QUOTE_PHASE="+phase, "ELITE_QUOTE_TENANT="+tenant,
				"ELITE_QUOTE_WRITER="+token("writer", tenant, []string{"lead:read", "quote:write"}, []string{"store"}),
				"ELITE_QUOTE_HANDOVER="+token("handover", tenant, []string{"handover:manage"}, []string{"store"}),
				"ELITE_QUOTE_HANDOVER_OTHER="+token("handover_other", tenant, []string{"handover:manage"}, []string{"store"}),
				"ELITE_QUOTE_LEAD="+token("lead-operator", tenant, []string{"lead:read", "lead:assign", "lead:update"}, []string{"store"}),
				"ELITE_QUOTE_SCHEDULER="+token("scheduler", tenant, []string{"availability:read", "availability:manage", "inventory:allocate"}, []string{"store"}),
				"ELITE_QUOTE_PLANNER="+token("planner", tenant, []string{"appointment:manage"}, []string{"store"}),
				"ELITE_QUOTE_RESOURCE="+token("resource", tenant, []string{"resource:manage"}, []string{"store"}),
				"ELITE_QUOTE_AVAILABILITY="+token("availability", tenant, []string{"availability:read"}, []string{"store"}),
				"ELITE_QUOTE_PAYMENT="+token("payment", tenant, []string{"payment:create"}, []string{"store"}),
				"ELITE_QUOTE_OPERATOR="+token("operator", tenant, []string{"inventory:allocate"}, []string{"store"}),
				"ELITE_QUOTE_READER="+token("reader", tenant, []string{"admin:read"}, []string{"store"}),
				"ELITE_QUOTE_CUSTOMER="+token("customer", tenant, []string{"customer:self"}, []string{"store"}),
				"ELITE_QUOTE_STRANGER="+token("stranger", tenant, []string{"customer:self"}, []string{"store"}),
				"ELITE_QUOTE_OBSERVER="+token("customer", tenant, []string{"lead:read"}, []string{"store"}),
				"ELITE_QUOTE_OTHER_ORG="+token("customer", tenant, []string{"customer:self"}, []string{"other"}),
				"ELITE_QUOTE_OTHER_TENANT="+token("customer", (randomid.Generator{}).New(), []string{"customer:self"}, []string{"store"}),
				"ENTERPRISE_API_BASE_URL="+api.URL, "ELITE_WEB_ROOT="+web,
				"AUTH_SESSION_SECRET=synthetic-"+(randomid.Generator{}).New()+(randomid.Generator{}).New(),
				"APP_BASE_URL="+edge.URL, "ELITE_BASE_URL="+edge.URL)
			server := exec.CommandContext(ctx, "node", filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", strings.Split(address, ":")[1])
			server.Dir, server.Env = web, env
			if err = server.Start(); err != nil {
				t.Fatal(err)
			}
			defer func() { _ = server.Process.Kill(); _ = server.Wait() }()
			ready := false
			client := &http.Client{Timeout: time.Second}
			for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
				res, err := client.Get(target.String() + "/icon.svg")
				if err == nil {
					res.Body.Close()
					if res.StatusCode == 200 {
						ready = true
						break
					}
				}
				select {
				case <-ctx.Done():
					t.Fatal("startup cancelled")
				case <-time.After(100 * time.Millisecond):
				}
			}
			if !ready {
				t.Fatal("web did not start")
			}
			gate := filepath.Join(web, "microsoft_playwright_browser_gate")
			pattern := "customer quote acceptance connects"
			if strings.HasPrefix(phase, "operation") {
				pattern = "operator stock continues existing orders"
			}
			if strings.HasPrefix(phase, "payment") {
				pattern = "operator payment request connects"
			}
			cmd := exec.CommandContext(ctx, "node", filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/enterprise-web.spec.mjs", "-g", pattern, "--project="+project, "--output="+filepath.Join(artifacts, project, phase))
			cmd.Dir, cmd.Env = gate, env
			output, err := cmd.CombinedOutput()
			clean := string(output)
			for _, item := range env {
				parts := strings.SplitN(item, "=", 2)
				if len(parts) == 2 && (parts[0] == "AUTH_SESSION_SECRET" || (strings.HasPrefix(parts[0], "ELITE_QUOTE_") && strings.Count(parts[1], ".") == 2)) {
					clean = strings.ReplaceAll(clean, parts[1], "[REDACTED_FIXTURE]")
				}
			}
			t.Log(clean)
			if err != nil {
				t.Fatalf("quote browser phase %s failed: %v", phase, err)
			}
			if phase == "operation-returns" && os.Getenv("ELITE_RETURN_SESSION_E2E") == "1" && !strings.Contains(clean, "RETURN_SESSION_BOUNDARY_PASS transitions=8 denied=12 authorized_reads=2 restored=2") {
				t.Fatal("return session gate did not produce required evidence")
			}
			if phase == "operation-returns" && os.Getenv("ELITE_RETURN_MULTITAB_E2E") == "1" && !strings.Contains(clean, "RETURN_MULTITAB_BROWSER_PASS races=6 posts=12 wins=6 conflicts=6 losses=2 opener_copy=1") {
				t.Fatal("multi-tab browser proof missing; a skipped test is not a PASS")
			}
			checks := []string{
				"select count(*) from sales.quotation where tenant_id=$1 and state='accepted' and order_id is not null",
				"select count(*) from sales.customer_order where tenant_id=$1 and state='placed' and total_minor_units=250000 and customer_principal_id='customer'",
				"select count(*) from sales.customer_order_line where tenant_id=$1 and quantity=1 and unit_price_minor_units=250000",
				"select count(*) from sales.quotation_acceptance where tenant_id=$1 and customer_principal_id='customer' and evidence_sha256_hex ~ '^[a-f0-9]{64}$'",
				"select count(*) from platform.outbox_event where tenant_id=$1 and event_type='quotation.accepted'",
				"select count(*) from platform.outbox_event where tenant_id=$1 and event_type='customer-order.placed'",
			}
			for _, sql := range checks {
				var n int
				if err := p.QueryRow(ctx, sql, tenant).Scan(&n); err != nil {
					t.Fatal(err)
				}
				if n != 3 {
					t.Fatalf("durable invariant: expected3 got%d", n)
				}
			}
			var denied int
			if err := p.QueryRow(ctx, "select count(*) from sales.quotation where tenant_id=$1 and quotation_id in ('quote-denied','quote-expired') and state='issued' and order_id is null", tenant).Scan(&denied); err != nil || denied != 2 {
				t.Fatal("negative requests changed quotes")
			}
			t.Log("QUOTE_ACCEPTANCE_BROWSER_POSTGRES_PASS phase=" + phase + " orders=3 acceptances=3 outbox=6 denied_unchanged=2")

			if phase == "operation-complete-checklist" {
				var rows, answers, events, actors, accepted, rejected int
				e := p.QueryRow(ctx, `select count(*),
 (select count(*) from sales.delivery_checklist_response where tenant_id=$1 and handover_id like 'ui-completion-%'),
 (select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id like 'ui-completion-%' and event_type='delivery-handover.checklist-completed'),
 (select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id like 'ui-completion-%' and event_type='delivery-handover.checklist-completed' and payload->>'actor_subject'='handover'),
 count(*)filter(where state='accepted'),count(*)filter(where state='rejected')
 from sales.delivery_handover where tenant_id=$1 and handover_id like 'ui-completion-%' and checklist_completed_at is not null and checklist_completed_by_subject='handover'`, tenant).Scan(&rows, &answers, &events, &actors, &accepted, &rejected)
				if e != nil || rows != 3 || answers != 6 || events != 3 || actors != 3 || accepted != 1 || rejected != 1 {
					t.Fatalf("completion counts=%d/%d/%d/%d/%d/%d error=%v", rows, answers, events, actors, accepted, rejected, e)
				}
				t.Log("CHECKLIST_COMPLETION_BROWSER_PASS handovers=3 answers=6 events=3 actors=3 accepted=1 rejected=1 presented=1")
			}
			if phase == "operation-returns" {
				var receipts, decisions, requests, events int
				e := p.QueryRow(ctx, `select (select count(*) from sales.return_receipt where tenant_id=$1 and authorization_id in ('ui-return-loss','ui-exchange-invalid','ui-return-race') and received_by_subject='handover'),(select count(*) from sales.return_disposition d join sales.return_receipt r on r.tenant_id=d.tenant_id and r.receipt_id=d.receipt_id where d.tenant_id=$1 and r.authorization_id in ('ui-return-loss','ui-exchange-invalid','ui-return-race') and d.decided_by_subject='handover'),(select count(*) from sales.return_effect_request e join sales.return_disposition d on d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id join sales.return_receipt r on r.tenant_id=d.tenant_id and r.receipt_id=d.receipt_id where e.tenant_id=$1 and r.authorization_id in ('ui-return-loss','ui-exchange-invalid','ui-return-race') and e.state='requested'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type in ('return.received','return.disposition-requested') and payload->>'actor_subject'='handover')`, tenant).Scan(&receipts, &decisions, &requests, &events)
				if e != nil || receipts != 3 || decisions != 3 || requests != 11 || events != 6 {
					t.Fatal("return browser effects", receipts, decisions, requests, events, e)
				}
				if os.Getenv("ELITE_RETURN_MULTITAB_E2E") == "1" {
					t.Log("RETURN_OPERATIONS_BROWSER_PASS receipts=3 decisions=3 requests=11 events=6 actors=6 multitab=PASS")
				} else {
					t.Log("RETURN_OPERATIONS_BROWSER_PASS receipts=3 decisions=3 requests=11 events=6 actors=6 replay=409 exact_recovery=PASS list_failure=PASS")
				}
			}

			if phase == "operation-publish-checklist" {
				var versions, items, events, actors int
				e := p.QueryRow(ctx, `select count(*),
 (select count(*) from sales.delivery_checklist_item where tenant_id=$1 and checklist_id like 'ui-publication-%'),
 (select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id like 'ui-publication-%' and event_type='delivery-checklist.published'),
 (select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id like 'ui-publication-%' and event_type='delivery-checklist.published' and payload->>'actor_subject'='handover')
 from sales.delivery_checklist_template where tenant_id=$1 and checklist_id like 'ui-publication-%' and state='published' and created_by_subject='handover'`, tenant).Scan(&versions, &items, &events, &actors)
				if e != nil || versions != 3 || items != 6 || events != 3 || actors != 3 {
					t.Fatalf("checklist publication counts=%d/%d/%d/%d err=%v", versions, items, events, actors, e)
				}
				t.Log("CHECKLIST_PUBLICATION_BROWSER_PASS versions=3 items=6 events=3 actors=3")
			}
			if phase == "operation-create-slot" {
				var rows, keys, events, actors int
				e := p.QueryRow(ctx, `select count(*),
  (select count(*) from platform.idempotency_record where tenant_id=$1 and scope='franchise-slot'),
  (select count(*) from platform.outbox_event e join platform.idempotency_record i on i.tenant_id=e.tenant_id and i.resource_id=e.aggregate_id and i.scope='franchise-slot' where e.tenant_id=$1 and e.event_type='appointment-slot.created'),
  (select count(*) from platform.outbox_event e join platform.idempotency_record i on i.tenant_id=e.tenant_id and i.resource_id=e.aggregate_id and i.scope='franchise-slot' where e.tenant_id=$1 and e.event_type='appointment-slot.created' and e.payload->>'actor_subject'='planner')
  from crm.appointment_slot a join platform.idempotency_record i on i.tenant_id=a.tenant_id and i.resource_id=a.slot_id and i.scope='franchise-slot' where a.tenant_id=$1`, tenant).Scan(&rows, &keys, &events, &actors)
				if e != nil || rows != 3 || keys != 3 || events != 3 || actors != 3 {
					t.Fatalf("slot creation counts=%d/%d/%d/%d err=%v", rows, keys, events, actors, e)
				}
				t.Log("SLOT_CREATION_BROWSER_PASS slots=3 keys=3 events=3 actors=3")
			}
			if phase == "operation-create-resource" {
				var rows, keys, events, actors, skills int
				e := p.QueryRow(ctx, `select count(*),
  (select count(*) from platform.idempotency_record where tenant_id=$1 and scope='franchise-resource'),
  (select count(*) from platform.outbox_event e join platform.idempotency_record i on i.tenant_id=e.tenant_id and i.resource_id=e.aggregate_id and i.scope='franchise-resource' where e.tenant_id=$1 and e.event_type='service-resource.created'),
  (select count(*) from platform.outbox_event e join platform.idempotency_record i on i.tenant_id=e.tenant_id and i.resource_id=e.aggregate_id and i.scope='franchise-resource' where e.tenant_id=$1 and e.event_type='service-resource.created' and e.payload->>'actor_subject'='resource'),
  (select count(*) from crm.resource_skill rs join platform.idempotency_record i on i.tenant_id=rs.tenant_id and i.resource_id=rs.resource_id and i.scope='franchise-resource' where rs.tenant_id=$1)
  from crm.service_resource r join platform.idempotency_record i on i.tenant_id=r.tenant_id and i.resource_id=r.resource_id and i.scope='franchise-resource' where r.tenant_id=$1`, tenant).Scan(&rows, &keys, &events, &actors, &skills)
				if e != nil || rows != 3 || keys != 3 || events != 3 || actors != 3 || skills != 6 {
					t.Fatalf("resource creation counts=%d/%d/%d/%d/%d err=%v", rows, keys, events, actors, skills, e)
				}
				t.Log("RESOURCE_CREATION_BROWSER_PASS resources=3 keys=3 events=3 actors=3 skills=6")
			}
			if phase == "operation-create-availability" {
				var entries, keys, events, actors, cancelled int
				e := p.QueryRow(ctx, `select count(*),count(*)filter(where a.state='cancelled'),(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='franchise-availability'),(select count(*) from platform.outbox_event e join platform.idempotency_record i on i.tenant_id=e.tenant_id and i.resource_id=e.aggregate_id and i.scope='franchise-availability' where e.tenant_id=$1 and e.event_type='availability-entry.created'),(select count(*) from platform.outbox_event e join platform.idempotency_record i on i.tenant_id=e.tenant_id and i.resource_id=e.aggregate_id and i.scope='franchise-availability' where e.tenant_id=$1 and e.event_type='availability-entry.created' and e.payload->>'actor_subject'='scheduler') from crm.availability_entry a join platform.idempotency_record i on i.tenant_id=a.tenant_id and i.resource_id=a.availability_id and i.scope='franchise-availability' where a.tenant_id=$1`, tenant).Scan(&entries, &cancelled, &keys, &events, &actors)
				if e != nil || entries != 3 || keys != 3 || events != 3 || actors != 3 || cancelled != 1 {
					t.Fatalf("availability creation counts=%d/%d/%d/%d cancelled=%d err=%v", entries, keys, events, actors, cancelled, e)
				}
				t.Log("AVAILABILITY_CREATION_BROWSER_PASS entries=3 keys=3 created_events=3 actors=3 cancelled=1")
			}
			if phase == "operation-quote" {
				var quotes, records, events, actors int
				err := p.QueryRow(ctx, `select (select count(*) from sales.quotation where tenant_id=$1 and lead_id like 'quote-create-%'),(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='franchise-quote' and status='completed'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='quotation.issued'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='quotation.issued' and payload->>'actor_subject'='writer')`, tenant).Scan(&quotes, &records, &events, &actors)
				if err != nil || quotes != 3 || records != 3 || events != 3 || actors != 3 {
					t.Fatalf("quote creation effects quotes=%d keys=%d events=%d err=%v", quotes, records, events, err)
				}
				t.Log("QUOTE_CREATION_DURABLE_PASS quotes=3 keys=3 events=3 actors=3")
			}
			if phase == "operation-lead" {
				var state, assignee string
				var version, events, actors int
				if err := p.QueryRow(ctx, `select lifecycle_state,assigned_subject,version from crm.lead where tenant_id=$1 and lead_id='operator-lead'`, tenant).Scan(&state, &assignee, &version); err != nil || state != "lost" || assignee != "sales-owner" || version != 5 {
					t.Fatalf("lead final %s %s %d %v", state, assignee, version, err)
				}
				if err := p.QueryRow(ctx, `select count(*),count(*) filter(where payload->>'actor_subject'='lead-operator') from platform.outbox_event where tenant_id=$1 and aggregate_id='operator-lead'`, tenant).Scan(&events, &actors); err != nil || events != 4 || actors != 4 {
					t.Fatalf("lead audit events=%d actors=%d err=%v", events, actors, err)
				}
				t.Log("LEAD_RECOVERY_PASS version=5 events=4 authenticated_actors=4 no_duplicates")
			}
			if phase == "operation-availability" {
				var entries, actors, events, eventActors int
				if err := p.QueryRow(ctx, `select count(*),count(*) filter(where cancelled_by_subject='scheduler' and cancellation_reason_code='schedule-correction') from crm.availability_entry where tenant_id=$1 and availability_id like 'cancel-window-%' and state='cancelled' and version=2`, tenant).Scan(&entries, &actors); err != nil || entries != 4 || actors != 4 {
					t.Fatalf("availability entries=%d actors=%d err=%v", entries, actors, err)
				}
				if err := p.QueryRow(ctx, `select count(*),count(*) filter(where payload->>'actor_subject'='scheduler') from platform.outbox_event where tenant_id=$1 and aggregate_id like 'cancel-window-%' and event_type='availability-entry.cancelled'`, tenant).Scan(&events, &eventActors); err != nil || events != 4 || eventActors != 4 {
					t.Fatalf("availability events=%d actors=%d err=%v", events, eventActors, err)
				}
				t.Log("AVAILABILITY_RECOVERY_PASS cancelled=4 version=2 events=4 authenticated_actors=4")
				var blockedState string
				var blockedVersion, blockedEvents int
				if err := p.QueryRow(ctx, `select state,version,(select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='cancel-window-4' and event_type='availability-entry.cancelled') from crm.availability_entry where tenant_id=$1 and availability_id='cancel-window-4'`, tenant).Scan(&blockedState, &blockedVersion, &blockedEvents); err != nil || blockedState != "active" || blockedVersion != 1 || blockedEvents != 0 {
					t.Fatalf("blocked interval=%s/%d events=%d err=%v", blockedState, blockedVersion, blockedEvents, err)
				}
				t.Log("AVAILABILITY_BOOKING_GUARD_PASS active=1 version=1 cancelled_events=0")
				var bookingEffects int
				if err := p.QueryRow(ctx, `select (select count(*) from crm.appointment where tenant_id=$1 and appointment_id<>'blocked-window-appointment')+(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='public-appointment')+(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='appointment.requested')`, tenant).Scan(&bookingEffects); err != nil || bookingEffects != 0 {
					t.Fatalf("booking outside working availability effects=%d err=%v", bookingEffects, err)
				}
				t.Log("AVAILABILITY_PUBLIC_HTTP_PASS withdrawn_slots=2 rejected_bookings=2 durable_booking_effects=0")
			}
			if phase == "operation-resolution" {
				// Existing delivery recovery remains a separate phase.
				for _, q := range []string{
					`select count(*) from sales.delivery_exception where tenant_id=$1 and handover_id='reject-fixture-handover' and state='resolved' and version=2`,
					`select count(*) from sales.delivery_exception_resolution r join sales.delivery_exception e on e.tenant_id=r.tenant_id and e.exception_id=r.exception_id where r.tenant_id=$1 and e.handover_id='reject-fixture-handover' and r.action='correct-and-represent' and r.resolved_by_subject='handover'`,
					`select count(*) from sales.delivery_handover where tenant_id=$1 and supersedes_handover_id='reject-fixture-handover' and state='prepared' and version=1`,
					`select count(*) from platform.outbox_event o join sales.delivery_exception e on e.tenant_id=o.tenant_id and e.exception_id=o.aggregate_id where o.tenant_id=$1 and e.handover_id='reject-fixture-handover' and o.event_type='delivery-exception.resolved'`,
				} {
					var n int
					if err := p.QueryRow(ctx, q, tenant).Scan(&n); err != nil || n != 1 {
						t.Fatalf("resolution durable count=%d err=%v", n, err)
					}
				}
				for _, action := range []string{"return", "exchange"} {
					for _, q := range []string{
						`select count(*) from sales.delivery_exception where tenant_id=$1 and handover_id=$2 and state='resolved' and version=2`,
						`select count(*) from sales.delivery_exception_resolution r join sales.delivery_exception e on e.tenant_id=r.tenant_id and e.exception_id=r.exception_id where r.tenant_id=$1 and e.handover_id=$2 and r.resolved_by_subject='handover'`,
						`select count(*) from sales.return_authorization where tenant_id=$1 and handover_id=$2 and state='authorized'`,
						`select count(*) from platform.outbox_event o join sales.delivery_exception e on e.tenant_id=o.tenant_id and e.exception_id=o.aggregate_id where o.tenant_id=$1 and e.handover_id=$2 and o.event_type='delivery-exception.resolved'`,
					} {
						var n int
						if err := p.QueryRow(ctx, q, tenant, action+"-fixture-handover").Scan(&n); err != nil || n != 1 {
							t.Fatalf("resolution %s count=%d err=%v", action, n, err)
						}
					}
				}
				t.Log("RESOLUTION_RECOVERY_PASS resolutions=3 successor=1 authorizations=2 events=3 actor=handover")
			}
			if phase == "operation-sections" {
				var forbidden int
				if err := p.QueryRow(ctx, `select count(*) from sales.delivery_checklist_template where tenant_id=$1 and checklist_id=$2`, tenant, "denied-probe").Scan(&forbidden); err != nil || forbidden != 0 {
					t.Fatalf("forbidden checklist effect: %d %v", forbidden, err)
				}
				t.Log("OPERATOR_SECTIONS_PASS roles=4 outage_recovered no_forbidden_checklist")
			}
			if phase == "delivery-action-loss" {
				for _, q := range []string{
					`select count(*) from sales.delivery_handover where tenant_id=$1 and handover_id='read-fixture-handover' and state='accepted' and version=3`,
					`select count(*) from sales.delivery_handover where tenant_id=$1 and handover_id='reject-fixture-handover' and state='rejected' and version=3`,
					`select count(*) from sales.delivery_exception where tenant_id=$1 and handover_id='reject-fixture-handover' and state='open'`,
					`select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='read-fixture-handover' and event_type='delivery-handover.accepted'`,
					`select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='reject-fixture-handover' and event_type='delivery-handover.rejected'`,
				} {
					var n int
					if err := p.QueryRow(ctx, q, tenant).Scan(&n); err != nil || n != 1 {
						t.Fatalf("delivery effect count=%d err=%v", n, err)
					}
				}
				t.Log("DELIVERY_ACTION_LOSS_PASS accepted=1 rejected=1 events=2 recovered_by_GET no_payment_or_shipment")
			}
			if strings.HasPrefix(phase, "delivery-read-") {
				var state string
				var events int
				if err := p.QueryRow(ctx, `select state,(select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='read-fixture-handover') from sales.delivery_handover where tenant_id=$1 and handover_id='read-fixture-handover'`, tenant).Scan(&state, &events); err != nil || state != "prepared" || events != 0 {
					t.Fatalf("read mutated fixture state=%s events=%d err=%v", state, events, err)
				}
				t.Log("DELIVERY_READ_HTTP_POSTGRES_PASS phase=" + phase + " no_acceptance_no_outbox")
			}
			if phase == "restart" && os.Getenv("ELITE_ORDER_E2E") != "1" && os.Getenv("ELITE_PAYMENT_E2E") != "1" {
				verifyQuoteCommerceContinuation(t, ctx, p, tenant)
			}
			if strings.HasPrefix(phase, "operation") {
				for _, sql := range []string{
					`select count(*) from inventory.stock_unit where tenant_id=$1 and state='reserved'`,
					`select count(*) from inventory.serial_reservation where tenant_id=$1 and status='reservation'`,
					`select count(*) from sales.customer_order_line where tenant_id=$1 and allocated_stock_unit_id is not null`,
					`select count(*) from platform.outbox_event where tenant_id=$1 and event_type='order.stock-allocated' and payload->>'actor_subject'='operator'`,
				} {
					var n int
					if err := p.QueryRow(ctx, sql, tenant).Scan(&n); err != nil || n != 2 {
						t.Fatalf("operation durable invariant expected2 got%d err=%v", n, err)
					}
				}
				t.Log("ORDER_STOCK_HTTP_BROWSER_PASS phase=" + phase + " reservations=2 allocations=2 outbox=2")
			}
			if strings.HasPrefix(phase, "payment") {
				for _, sql := range []string{`select count(*) from payment.payment_attempt where tenant_id=$1 and state='created' and amount_minor_units=250000 and currency='ARS' and provider_code='stripe'`, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='payment.requested' and payload->>'actor_subject'='payment'`} {
					var n int
					if err := p.QueryRow(ctx, sql, tenant).Scan(&n); err != nil || n != 2 {
						t.Fatalf("payment invariant expected2 got%d err%v", n, err)
					}
				}
				t.Log("PAYMENT_UI_HTTP_POSTGRES_PASS phase=" + phase + " intents=2 audited_events=2 no_provider_calls")
			}
		}()
	}
}

// AUTHORED continuation of the SAME orders created through the customer browser.
// Exercises existing repository contracts, not operator UI or live payment approval.
func verifyQuoteCommerceContinuation(t *testing.T, ctx context.Context, pool *pgxpool.Pool, tenant string) {
	t.Helper()
	ids := randomid.Generator{}
	repo := postgres.NewCommerce(pool)
	type demand struct{ order, line string }
	demands := make([]demand, 2)
	for i, quote := range []string{"quote-success", "quote-race"} {
		if err := pool.QueryRow(ctx, `select o.order_id,l.line_id from sales.quotation q join sales.customer_order o on o.tenant_id=q.tenant_id and o.order_id=q.order_id join sales.customer_order_line l on l.tenant_id=o.tenant_id and l.order_id=o.order_id where q.tenant_id=$1 and q.quotation_id=$2`, tenant, quote).Scan(&demands[i].order, &demands[i].line); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := pool.Exec(ctx, `insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at) values($1,'quote-stock','store','variant','QUOTE-SERIAL','available',1,clock_timestamp())`, tenant); err != nil {
		t.Fatal(err)
	}
	if err := repo.AllocateStock(ctx, tenant, "other", demands[0].order, demands[0].line, "quote-stock", 1, 1, ids.New()); !errors.Is(err, commerce.ErrConflict) {
		t.Fatalf("cross-org allocation=%v", err)
	}
	if err := repo.AllocateStock(ctx, ids.New(), "store", demands[0].order, demands[0].line, "quote-stock", 1, 1, ids.New()); !errors.Is(err, commerce.ErrConflict) {
		t.Fatalf("cross-tenant allocation=%v", err)
	}
	if err := repo.AllocateStock(ctx, tenant, "store", demands[0].order, demands[0].line, "quote-stock", 99, 1, ids.New()); !errors.Is(err, commerce.ErrConflict) {
		t.Fatalf("stale allocation=%v", err)
	}
	start := make(chan struct{})
	results := make(chan error, 2)
	for _, d := range demands {
		go func(d demand) {
			<-start
			results <- repo.AllocateStock(ctx, tenant, "store", d.order, d.line, "quote-stock", 1, 1, ids.New())
		}(d)
	}
	close(start)
	success, conflict := 0, 0
	for range demands {
		err := <-results
		if err == nil {
			success++
		} else if errors.Is(err, commerce.ErrConflict) {
			conflict++
		} else {
			t.Fatal(err)
		}
	}
	if success != 1 || conflict != 1 {
		t.Fatalf("stock race success=%d conflict=%d", success, conflict)
	}
	var winner, line string
	if err := pool.QueryRow(ctx, `select order_id,line_id from sales.customer_order_line where tenant_id=$1 and allocated_stock_unit_id='quote-stock'`, tenant).Scan(&winner, &line); err != nil {
		t.Fatal(err)
	}
	if err := repo.AllocateStock(ctx, tenant, "store", winner, line, "quote-stock", 2, 2, ids.New()); !errors.Is(err, commerce.ErrConflict) {
		t.Fatalf("duplicate allocation=%v", err)
	}
	payment := commerce.PaymentAttempt{ID: ids.New(), OrderID: winner, OrganizationID: "other", ProviderCode: "synthetic-no-dispatch", Currency: "ARS", AmountMinorUnits: 250000}
	if err := repo.CreatePaymentAttempt(ctx, tenant, ids.New(), "quote-payment-key", payment); !errors.Is(err, commerce.ErrConflict) {
		t.Fatalf("cross-org payment=%v", err)
	}
	payment.OrganizationID = "store"
	payment.AmountMinorUnits = 1
	if err := repo.CreatePaymentAttempt(ctx, tenant, ids.New(), "quote-payment-key", payment); !errors.Is(err, commerce.ErrConflict) {
		t.Fatalf("wrong amount payment=%v", err)
	}
	payment.AmountMinorUnits = 250000
	payment.Currency = "USD"
	if err := repo.CreatePaymentAttempt(ctx, tenant, ids.New(), "quote-payment-key", payment); !errors.Is(err, commerce.ErrConflict) {
		t.Fatalf("wrong currency payment=%v", err)
	}
	payment.Currency = "ARS"
	if err := repo.CreatePaymentAttempt(ctx, tenant, ids.New(), "quote-payment-key", payment); err != nil {
		t.Fatal(err)
	}
	payment.ID = ids.New()
	if err := repo.CreatePaymentAttempt(ctx, tenant, ids.New(), "quote-payment-key", payment); err == nil {
		t.Fatal("duplicate payment key accepted")
	}
	// A fresh pool proves durable reads, not a PostgreSQL process restart/PITR.
	reopened, err := pgxpool.NewWithConfig(ctx, pool.Config().Copy())
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	for _, sql := range []string{
		`select count(*) from inventory.serial_reservation where tenant_id=$1 and stock_unit_id='quote-stock' and status='reservation'`,
		`select count(*) from inventory.stock_unit where tenant_id=$1 and stock_unit_id='quote-stock' and state='reserved' and version=2`,
		`select count(*) from sales.customer_order_line where tenant_id=$1 and allocated_stock_unit_id='quote-stock'`,
		`select count(*) from payment.payment_attempt where tenant_id=$1 and state='created' and currency='ARS' and amount_minor_units=250000`,
		`select count(*) from platform.outbox_event where tenant_id=$1 and event_type='order.stock-allocated'`,
		`select count(*) from platform.outbox_event where tenant_id=$1 and event_type='payment.requested'`,
	} {
		var n int
		if err := reopened.QueryRow(ctx, sql, tenant).Scan(&n); err != nil || n != 1 {
			t.Fatalf("commerce durable count=%d err=%v", n, err)
		}
	}
	t.Log("QUOTE_COMMERCE_CONTINUATION_PASS stock=1 competing_orders=2 payment_created=1 stock_events=1 payment_events=1 provider_calls=0")
}

func (r *journeyRepo) CreateQuoteAs(ctx context.Context, tenant, key string, v franchisejourney.Quote, hash, event, actor string) (franchisejourney.Quote, bool, error) {
	r.quoteActor = actor
	return r.CreateQuote(ctx, tenant, key, v, hash, event)
}

func (r *journeyRepo) CreateAvailabilityOnce(ctx context.Context, tenant, subject, key, hash string, v franchisejourney.AvailabilityEntry, event string) (franchisejourney.AvailabilityEntry, bool, error) {
	value, err := r.CreateAvailability(ctx, tenant, subject, v, event)
	return value, false, err
}

func TestResourceCreationHTTPRecovery(t *testing.T) {
	if os.Getenv("TEST_DATABASE_URL") == "" {
		t.Skip("disposable database required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cfg, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'resource-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
	} {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	verifier, token := confirmationTestIssuer(t)
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: franchisejourney.NewService(postgres.NewFranchiseJourney(pool), randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
	api := httptest.NewServer(mux)
	defer api.Close()
	bearer := token("resource-manager", tenant, []string{"resource:manage"}, []string{"store"})
	body := `{"organization_id":"store","display_name":"Synthetic bay","kind":"service-bay","skills":["service"]}`
	post := func() (int, map[string]any) {
		t.Helper()
		req, err := http.NewRequestWithContext(ctx, "POST", api.URL+"/v1/franchise/resources", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("content-type", "application/json")
		req.Header.Set("authorization", "Bearer "+bearer)
		req.Header.Set("idempotency-key", "resource-recovery-fixture-key")
		res, err := api.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		var value map[string]any
		if err = json.NewDecoder(res.Body).Decode(&value); err != nil {
			t.Fatal(err)
		}
		return res.StatusCode, value
	}
	first, a := post()
	second, b := post()
	var count int
	if err = pool.QueryRow(ctx, `select count(*) from crm.service_resource where tenant_id=$1`, tenant).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if first != 201 || second != 200 || a["id"] != b["id"] || count != 1 {
		t.Fatalf("resource replay first=%d second=%d same_id=%v durable_resources=%d", first, second, a["id"] == b["id"], count)
	}
	var actor string
	if err = pool.QueryRow(ctx, `select payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and event_type='service-resource.created'`, tenant).Scan(&actor); err != nil || actor != "resource-manager" {
		t.Fatalf("actor not retained: %v", err)
	}
}

func TestSlotCreationHTTPRecovery(t *testing.T) {
	if os.Getenv("TEST_DATABASE_URL") == "" {
		t.Skip("disposable database required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cfg, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'resource-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
	} {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	verifier, token := confirmationTestIssuer(t)
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: franchisejourney.NewService(postgres.NewFranchiseJourney(pool), randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
	api := httptest.NewServer(mux)
	defer api.Close()
	bearer := token("slot-manager", tenant, []string{"appointment:manage"}, []string{"store"})
	start := time.Now().UTC().Add(72 * time.Hour).Truncate(time.Second)
	if _, err := postgres.NewFranchiseJourney(pool).CreateAvailability(ctx, tenant, "slot-manager", franchisejourney.AvailabilityEntry{ID: "working", OrganizationID: "store", EntryType: "working", StartsAt: start.Add(-time.Hour), EndsAt: start.Add(2 * time.Hour)}, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	payload, err := json.Marshal(map[string]any{"organization_id": "store", "kind": "service", "starts_at": start, "ends_at": start.Add(time.Hour), "capacity": 2})
	if err != nil {
		t.Fatal(err)
	}
	body := string(payload)
	post := func() (int, map[string]any) {
		t.Helper()
		req, err := http.NewRequestWithContext(ctx, "POST", api.URL+"/v1/franchise/appointment-slots", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("content-type", "application/json")
		req.Header.Set("authorization", "Bearer "+bearer)
		req.Header.Set("idempotency-key", "slot-recovery-fixture-key")
		res, err := api.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		var value map[string]any
		if err = json.NewDecoder(res.Body).Decode(&value); err != nil {
			t.Fatal(err)
		}
		return res.StatusCode, value
	}
	first, a := post()
	second, b := post()
	var count int
	if err = pool.QueryRow(ctx, `select count(*) from crm.appointment_slot where tenant_id=$1`, tenant).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if first != 201 || second != 200 || a["id"] != b["id"] || count != 1 {
		t.Fatalf("slot replay first=%d second=%d same_id=%v durable_slots=%d", first, second, a["id"] == b["id"], count)
	}
	var actor string
	if err = pool.QueryRow(ctx, `select payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and event_type='appointment-slot.created'`, tenant).Scan(&actor); err != nil || actor != "slot-manager" {
		t.Fatalf("actor not retained: %v", err)
	}
}

func TestChecklistPublicationHTTPRecovery(t *testing.T) {
	if os.Getenv("TEST_DATABASE_URL") == "" {
		t.Skip("disposable database required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cfg, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'checklist-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
	} {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	verifier, token := confirmationTestIssuer(t)
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: franchisejourney.NewService(postgres.NewFranchiseJourney(pool), randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
	api := httptest.NewServer(mux)
	defer api.Close()
	bearer := token("checklist-publisher", tenant, []string{"handover:manage"}, []string{"store"})
	body := `{"organization_id":"store","checklist_id":"fixture-checklist","version":1,"title":"Synthetic checklist","items":[{"id":"serial","prompt":"Verify synthetic serial","response_type":"serial","required":true},{"id":"safe","prompt":"Confirm synthetic review","response_type":"confirmation","required":true}]}`
	req, err := http.NewRequestWithContext(ctx, "POST", api.URL+"/v1/franchise/delivery-checklists", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+bearer)
	res, err := api.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	var published franchisejourney.DeliveryChecklist
	if res.StatusCode != 201 {
		res.Body.Close()
		t.Fatalf("publication status=%d", res.StatusCode)
	}
	if err = json.NewDecoder(res.Body).Decode(&published); err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	var actor string
	if err = pool.QueryRow(ctx, `select coalesce(payload->>'actor_subject','') from platform.outbox_event where tenant_id=$1 and event_type='delivery-checklist.published'`, tenant).Scan(&actor); err != nil || actor != "checklist-publisher" {
		t.Errorf("publication actor=%q error=%v", actor, err)
	}
	req, err = http.NewRequestWithContext(ctx, "GET", api.URL+"/v1/franchise/delivery-checklists/result?organization_id=store&checklist_id=fixture-checklist&version=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("authorization", "Bearer "+bearer)
	res, err = api.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("published version recovery status=%d expected=200", res.StatusCode)
	}
	var recovered franchisejourney.DeliveryChecklist
	if err = json.NewDecoder(res.Body).Decode(&recovered); err != nil {
		t.Fatal(err)
	}
	if recovered.ID != published.ID || recovered.Version != 1 || recovered.State != "published" || len(recovered.Items) != 2 || recovered.Items[0].Ordinal != 1 || recovered.Items[1].Ordinal != 2 {
		t.Fatalf("incoherent recovered checklist: %+v", recovered)
	}
}

func TestChecklistCompletionHTTPRecovery(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database required")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'complete-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','synthetic@example.test')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','SYNTHETIC-SERIAL','sold',1,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','delivered','ARS',123456,1)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'completion-fixture','store','order','customer','stock','prepared',1)`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := postgres.NewFranchiseJourney(pool)
	checklist := franchisejourney.DeliveryChecklist{ID: "completion-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}, {ID: "confirmed", Ordinal: 2, Prompt: "Confirmación", ResponseType: "confirmation", Required: true}}}
	if _, err = repo.PublishDeliveryChecklist(ctx, tenant, "fixture-author", checklist, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	verifier, token := confirmationTestIssuer(t)
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: franchisejourney.NewService(repo, randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
	api := httptest.NewServer(mux)
	defer api.Close()
	bearer := token("checklist-operator", tenant, []string{"handover:manage"}, []string{"store"})
	body := `{"organization_id":"store","version":1,"checklist_id":"completion-checklist","checklist_version":1,"responses":[{"item_id":"serial","response_text":"SYNTHETIC-SERIAL"},{"item_id":"confirmed","response_text":"confirmed"}]}`
	for i := 0; i < 2; i++ {
		req, e := http.NewRequestWithContext(ctx, "POST", api.URL+"/v1/franchise/handovers/completion-fixture/complete-checklist", strings.NewReader(body))
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Set("content-type", "application/json")
		req.Header.Set("authorization", "Bearer "+bearer)
		res, e := api.Client().Do(req)
		if e != nil {
			t.Fatal(e)
		}
		res.Body.Close()
		expected := 200
		if i == 1 {
			expected = 409
		}
		if res.StatusCode != expected {
			t.Fatalf("complete/replay status=%d want=%d", res.StatusCode, expected)
		}
	}
	var answers, events int
	if err = pool.QueryRow(ctx, `select (select count(*) from sales.delivery_checklist_response where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='delivery-handover.checklist-completed' and payload->>'actor_subject'='checklist-operator')`, tenant).Scan(&answers, &events); err != nil || answers != 2 || events != 1 {
		t.Fatal(answers, events, err)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", api.URL+"/v1/franchise/handovers/completion-fixture/checklist-result?organization_id=store", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("authorization", "Bearer "+bearer)
	res, err := api.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("completion lookup status=%d; durable answers=%d events=%d", res.StatusCode, answers, events)
	}
	var value struct {
		HandoverID  string                               `json:"handover_id"`
		State       string                               `json:"state"`
		Version     int64                                `json:"version"`
		ChecklistID string                               `json:"checklist_id"`
		Actor       string                               `json:"actor_subject"`
		Responses   []franchisejourney.ChecklistResponse `json:"responses"`
	}
	if err = json.NewDecoder(res.Body).Decode(&value); err != nil || value.HandoverID != "completion-fixture" || value.State != "presented" || value.Version != 2 || value.ChecklistID != checklist.ID || value.Actor != "checklist-operator" || len(value.Responses) != 2 {
		t.Fatal(value, err)
	}
	t.Log("CHECKLIST_COMPLETION_HTTP_PASS complete=200 replay=409 answers=2 events=1 lookup=200 actor=checklist-operator")
}

func TestReturnOperationsHTTPRecovery(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database required")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'complete-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','synthetic@example.test')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','SYNTHETIC-SERIAL','sold',1,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','delivered','ARS',123456,1)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'completion-fixture','store','order','customer','stock','prepared',1)`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := postgres.NewFranchiseJourney(pool)
	checklist := franchisejourney.DeliveryChecklist{ID: "completion-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}, {ID: "confirmed", Ordinal: 2, Prompt: "Confirmación", ResponseType: "confirmation", Required: true}}}
	if _, err = repo.PublishDeliveryChecklist(ctx, tenant, "fixture-author", checklist, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}

	responses := []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SYNTHETIC-SERIAL"}, {ItemID: "confirmed", ResponseText: "confirmed"}}
	if _, err = repo.CompleteDeliveryChecklist(ctx, tenant, "store", "fixture-operator", "completion-fixture", 1, checklist.ID, 1, responses, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	exception, err := repo.RejectHandover(ctx, tenant, "store", "customer", "completion-fixture", 2, "serial-mismatch", "Synthetic rejection", strings.Repeat("a", 64), randomid.Generator{}.New(), randomid.Generator{}.New())
	if err != nil {
		t.Fatal(err)
	}
	resolution, err := repo.ResolveDeliveryException(ctx, tenant, "store", "fixture-manager", exception.ID, 1, "return", "Synthetic authorization", "", "return-authorization", randomid.Generator{}.New(), randomid.Generator{}.New())
	if err != nil || resolution.ReturnAuthorizationID != "return-authorization" {
		t.Fatal(resolution, err)
	}
	verifier, token := confirmationTestIssuer(t)
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: franchisejourney.NewService(repo, randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
	api := httptest.NewServer(mux)
	defer api.Close()
	bearer := token("return-operator", tenant, []string{"handover:manage"}, []string{"store"})
	post := func(path, body string, target any) int {
		req, e := http.NewRequestWithContext(ctx, "POST", api.URL+path, strings.NewReader(body))
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Set("content-type", "application/json")
		req.Header.Set("authorization", "Bearer "+bearer)
		res, e := api.Client().Do(req)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		if res.StatusCode == 200 && target != nil {
			if e = json.NewDecoder(res.Body).Decode(target); e != nil {
				t.Fatal(e)
			}
		}
		return res.StatusCode
	}
	body := fmt.Sprintf(`{"organization_id":"store","serial_number":"SYNTHETIC-SERIAL","condition_code":"sealed","notes":"Synthetic receipt","evidence_sha256":"%s"}`, strings.Repeat("b", 64))
	var receipt franchisejourney.ReturnReceipt
	if code := post("/v1/franchise/return-authorizations/return-authorization/receive", body, &receipt); code != 200 || receipt.ID == "" {
		t.Fatal("receive", code, receipt)
	}
	if code := post("/v1/franchise/return-authorizations/return-authorization/receive", body, nil); code != 409 {
		t.Fatal("receive replay", code)
	}
	lookup := func(label string, decided bool) {
		req, e := http.NewRequestWithContext(ctx, "GET", api.URL+"/v1/franchise/returns/result?organization_id=store&authorization_id=return-authorization", nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Set("authorization", "Bearer "+bearer)
		res, e := api.Client().Do(req)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			t.Errorf("%s lookup=%d after durable commit", label, res.StatusCode)
			return
		}
		var value franchisejourney.ReturnCase
		if e = json.NewDecoder(res.Body).Decode(&value); e != nil {
			t.Fatal(e)
		}
		if value.AuthorizationID != "return-authorization" || value.OrganizationID != "store" || value.Receipt == nil || value.Receipt.ID != receipt.ID || value.Receipt.ReceivedBySubject != "return-operator" {
			t.Fatal("receipt lookup", value)
		}
		if decided && (value.Disposition == nil || value.Disposition.ReceiptID != receipt.ID || value.Disposition.DecidedBySubject != "return-operator" || len(value.Disposition.Effects) != 4) {
			t.Fatal("decision lookup", value)
		}
		if !decided && value.Disposition != nil {
			t.Fatal("uncreated decision exposed")
		}
	}
	lookup("receipt", false)
	body = `{"organization_id":"store","inventory_action":"quarantine","notes":"Synthetic decision"}`
	var disposition franchisejourney.ReturnDisposition
	if code := post("/v1/franchise/return-receipts/"+receipt.ID+"/decide", body, &disposition); code != 200 || disposition.ID == "" || len(disposition.Effects) != 4 {
		t.Fatal("decision", code, disposition)
	}
	if code := post("/v1/franchise/return-receipts/"+receipt.ID+"/decide", body, nil); code != 409 {
		t.Fatal("decision replay", code)
	}
	lookup("decision", true)
	var receipts, decisions, effects, events int
	err = pool.QueryRow(ctx, `select (select count(*) from sales.return_receipt where tenant_id=$1),(select count(*) from sales.return_disposition where tenant_id=$1),(select count(*) from sales.return_effect_request where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and payload->>'actor_subject'='return-operator')`, tenant).Scan(&receipts, &decisions, &effects, &events)
	if err != nil || receipts != 1 || decisions != 1 || effects != 4 || events != 2 {
		t.Fatal(receipts, decisions, effects, events, err)
	}
	if !t.Failed() {
		t.Log("RETURN_OPERATIONS_HTTP_PASS receive=200/replay409 decide=200/replay409 receipts=1 decisions=1 effects=4 events=2 lookup=200")
	}
}
