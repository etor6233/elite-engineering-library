# Go Customer Survey API

## 1. Metadata

```yaml
pack_id: "GO-CUSTOMER-SURVEY-API"
pack_version: "0.2.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Connected authenticated customer survey with PostgreSQL replay, scoped NPS and bounded retention"
stacks: ["Go 1.26.8", "PostgreSQL 18.6"]
compatible_with: ["GO-ELECTROMOBILITY-APPLICATION 1.9.4", "GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.x", "GO-ENTERPRISE-QUERY-API 0.x"]
incompatible_with: ["automatic campaign or legal-consent certification", "unscoped response access", "production admission inferred from synthetic fixtures"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://www.netpromotersystem.com/about/measuring-your-net-promoter-score/", "https://www.postgresql.org/docs/18/explicit-locking.html", "https://go.dev/doc/security/fuzz/"]
verified_at: "2026-09-11"
```

## 2. Applicability

Opt-in connected survey reference. Project-owned definition, notice, recommendation question,
retention, threshold and authorized distribution of invitation links are required.
The existing GO-SURVEYS-CORE remains an unchanged isolated reference. No upstream
enterprise implementation is claimed. See docs/customer-surveys.md from the Go pack.

## 3. Architecture contract

One answer per tenant/organization/survey/customer. OIDC identity and active CRM
membership govern writes; customer:self and surveys:read are separate grants.
PostgreSQL owns commit time, conflict/replay and expiry. UI performs no automatic
POST retry and recovers by GET. No raw errors, credentials or foreign identifiers
are returned. Count is visible to an authorized aggregate reader below the NPS
threshold. Synthetic tests do not supply production grants or business policy.

## 4. Exact file manifest

```text
CREATE docs/posthog-nps.md
CREATE third_party/posthog-nps/LICENSE
CREATE third_party/posthog-nps/utils.test.ts.txt
CREATE third_party/posthog-nps/utils.ts.txt
CREATE internal/posthognps/nps.go
CREATE internal/posthognps/nps_fuzz_test.go
CREATE internal/posthognps/nps_test.go
CREATE third_party/posthog-nps/source-lock.json
CREATE internal/customerfeedback/service.go
CREATE internal/customerfeedback/service_test.go
CREATE db/migrations/0054_customer_feedback.up.sql
CREATE db/migrations/0054_customer_feedback.down.sql
CREATE internal/platform/postgres/customerfeedback.go
CREATE internal/platform/httpapi/customerfeedback.go
CREATE internal/platform/httpapi/customerfeedback_test.go
CREATE internal/platform/httpapi/customerfeedback_integration_test.go
CREATE internal/platform/httpapi/customerfeedback_browser_test.go
CREATE internal/platform/httpapi/customerfeedback_fuzz_test.go
CREATE cmd/electromobility-api/customerfeedback_activation.go
CREATE cmd/electromobility-api/customerfeedback_activation_test.go
CREATE cmd/customer-survey-retention/main.go
CREATE cmd/customer-survey-retention/main_test.go
CREATE docs/customer-surveys.md
CREATE internal/platform/httpapi/customerfeedback_metrics.go
```

## 5. Materialization blocks

### FILE: `internal/customerfeedback/service.go`
```yaml
block_id: "GO-CUSTOMER-SURVEY-API:internal/customerfeedback/service.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "cb0f53d269c2a7059c50a0a7c12c4bdb53ca123e015508cbd1ddc32ad64f6620"
variables: []
secrets_allowed: false
```
````go
// AUTHORED customer feedback boundary. Policy texts and survey activation are project-owned.
package customerfeedback

import (
	"context"
	"elite.local/enterprise/internal/posthognps"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalid     = errors.New("invalid survey request")
	ErrUnavailable = errors.New("survey unavailable")
	ErrConflict    = errors.New("survey response conflicts with stored response")
	ErrNotFound    = errors.New("survey response not found")
)

type Scope struct{ Tenant, Organization, Customer string }
type Definition struct {
	ID             string `json:"id"`
	Prompt         string `json:"prompt"`
	ConsentVersion string `json:"consent_version"`
	ConsentNotice  string `json:"consent_notice"`
	Accepting      bool   `json:"accepting"`
}
type Answer struct {
	Score          int       `json:"score"`
	ConsentVersion string    `json:"consent_version"`
	ReceivedAt     time.Time `json:"received_at"`
}
type Submission struct {
	Score          int    `json:"score"`
	Consent        bool   `json:"consent"`
	ConsentVersion string `json:"consent_version"`
}
type Result struct {
	Answer Answer `json:"answer"`
	Replay bool   `json:"replay"`
}
type Summary struct {
	Responses int64    `json:"responses"`
	Available bool     `json:"available"`
	NPS       *float64 `json:"nps"`
}
type Repository interface {
	Definition(context.Context, Scope, string) (Definition, error)
	Submit(context.Context, Scope, string, Submission) (Result, error)
	OwnAnswer(context.Context, Scope, string) (Answer, error)
	Summary(context.Context, Scope, string) (Summary, error)
}
type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }
func valid(value string) bool {
	return strings.TrimSpace(value) == value && value != "" && len(value) <= 128 && !strings.ContainsAny(value, "\x00\r\n")
}
func validate(scope Scope, id string, customer bool) error {
	if !valid(scope.Tenant) || !valid(scope.Organization) || !valid(id) || (customer && !valid(scope.Customer)) {
		return ErrInvalid
	}
	return nil
}
func (s *Service) Definition(ctx context.Context, scope Scope, id string) (Definition, error) {
	if e := validate(scope, id, true); e != nil {
		return Definition{}, e
	}
	return s.repository.Definition(ctx, scope, id)
}
func ValidateSubmission(scope Scope, id string, input Submission) error {
	if e := validate(scope, id, true); e != nil {
		return e
	}
	if input.Score < 0 || input.Score > 10 || !input.Consent || !valid(input.ConsentVersion) {
		return ErrInvalid
	}
	return nil
}
func (s *Service) Submit(ctx context.Context, scope Scope, id string, input Submission) (Result, error) {
	if e := ValidateSubmission(scope, id, input); e != nil {
		return Result{}, e
	}
	return s.repository.Submit(ctx, scope, id, input)
}
func (s *Service) OwnAnswer(ctx context.Context, scope Scope, id string) (Answer, error) {
	if e := validate(scope, id, true); e != nil {
		return Answer{}, e
	}
	return s.repository.OwnAnswer(ctx, scope, id)
}
func (s *Service) Summary(ctx context.Context, scope Scope, id string) (Summary, error) {
	if e := validate(scope, id, false); e != nil {
		return Summary{}, e
	}
	return s.repository.Summary(ctx, scope, id)
}

// NPS requires a nonempty population and an explicit project-owned reporting threshold.
// The selected PostHog MIT owner calculates the score; this caller preserves validation
// and the project-owned disclosure threshold. No representative-sampling claim.
func NPS(count, promoters, detractors, minimum int64) (Summary, error) {
	if count < 0 || promoters < 0 || detractors < 0 || minimum < 1 || promoters > count || detractors > count-promoters {
		return Summary{}, ErrInvalid
	}
	result := Summary{Responses: count}
	if count >= minimum {
		value := posthognps.Calculate(count, promoters, detractors).Score
		result.Available = true
		result.NPS = &value
	}
	return result, nil
}
````

### FILE: `internal/customerfeedback/service_test.go`
```yaml
block_id: "GO-CUSTOMER-SURVEY-API:internal/customerfeedback/service_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "0aee215524094374b0976caaa7ca8426b5c2cd178dffa9318374b017c478b916"
variables: []
secrets_allowed: false
```
````go
package customerfeedback

import (
	"context"
	"testing"
)

type neverRepository struct{}

func (neverRepository) Definition(context.Context, Scope, string) (Definition, error) {
	panic("invalid request reached persistence")
}
func (neverRepository) Submit(context.Context, Scope, string, Submission) (Result, error) {
	panic("invalid request reached persistence")
}
func (neverRepository) OwnAnswer(context.Context, Scope, string) (Answer, error) {
	panic("invalid request reached persistence")
}
func (neverRepository) Summary(context.Context, Scope, string) (Summary, error) {
	panic("invalid request reached persistence")
}
func TestFeedbackNPSGoldenAndThreshold(t *testing.T) {
	for _, c := range []struct {
		n, p, d, m int64
		available  bool
		want       float64
	}{
		{0, 0, 0, 1, false, 0}, {1, 1, 0, 2, false, 0}, {5, 1, 3, 1, true, -40}, {3, 2, 1, 2, true, 100.0 / 3}, {2, 0, 2, 1, true, -100}, {2, 2, 0, 1, true, 100},
	} {
		got, e := NPS(c.n, c.p, c.d, c.m)
		if e != nil || got.Available != c.available || got.Responses != c.n {
			t.Fatalf("%+v %v", got, e)
		}
		if c.available {
			if got.NPS == nil || *got.NPS-c.want > 1e-12 || c.want-*got.NPS > 1e-12 {
				t.Fatal(got)
			}
		} else if got.NPS != nil {
			t.Fatal("score exposed below threshold")
		}
	}
	for _, c := range [][4]int64{{-1, 0, 0, 1}, {1, 2, 0, 1}, {1, 1, 1, 1}, {1, 0, 0, 0}, {1, -1, 0, 1}} {
		if _, e := NPS(c[0], c[1], c[2], c[3]); e != ErrInvalid {
			t.Fatal(c, e)
		}
	}
}
func TestFeedbackRejectsBeforePersistence(t *testing.T) {
	service := NewService(neverRepository{})
	scope := Scope{"tenant", "org", "customer"}
	for _, input := range []Submission{{-1, true, "v1"}, {11, true, "v1"}, {1, false, "v1"}, {1, true, ""}, {1, true, "v1\n"}} {
		if _, e := service.Submit(context.Background(), scope, "survey", input); e != ErrInvalid {
			t.Fatal(e)
		}
	}
	for _, s := range []Scope{{}, {"tenant", "org", ""}, {"tenant", "\x00", "customer"}} {
		if _, e := service.Definition(context.Background(), s, "survey"); e != ErrInvalid {
			t.Fatal(e)
		}
	}
}
````

### FILE: `db/migrations/0054_customer_feedback.up.sql`
```yaml
block_id: "GO-CUSTOMER-SURVEY-API:db/migrations/0054_customer_feedback.up.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "f396c089ede507812a3d427b9e826eff0ada703b72d70d578f30a4af439a2cb5"
variables: []
secrets_allowed: false
```
````sql
begin;
-- AUTHORED: no campaign, invitation, consent policy or activation is inferred.
create table crm.survey_definition (
 tenant_id uuid not null,
 organization_id text not null,
 survey_id text not null check(survey_id ~ '^[A-Za-z0-9][A-Za-z0-9_-]{0,127}$'),
 prompt text not null check(length(prompt) between 1 and 1000),
 consent_version text not null check(consent_version ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$'),
 consent_notice text not null check(length(consent_notice) between 1 and 4000),
 opens_at timestamptz not null,
 closes_at timestamptz not null,
 retain_until timestamptz not null,
 minimum_responses integer not null check(minimum_responses between 1 and 1000000),
 active boolean not null default false,
 primary key(tenant_id,organization_id,survey_id),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 check(opens_at<closes_at and closes_at<retain_until)
);
create table crm.survey_response (
 tenant_id uuid not null,
 organization_id text not null,
 survey_id text not null,
 customer_principal_id text not null,
 score smallint not null check(score between 0 and 10),
 consent_version text not null,
 received_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,organization_id,survey_id,customer_principal_id),
 foreign key(tenant_id,organization_id,survey_id) references crm.survey_definition(tenant_id,organization_id,survey_id),
 foreign key(tenant_id,customer_principal_id) references crm.customer_profile(tenant_id,customer_principal_id)
);
-- Definitions are versioned by survey ID. Activation may change; historical terms may not.
create function crm.survey_definition_immutable() returns trigger language plpgsql as $$
begin
 if (to_jsonb(new)-'active') is distinct from (to_jsonb(old)-'active') then
  raise exception using errcode='23514',message='survey definition is immutable; create a new survey id';
 end if;
 return new;
end $$;
create trigger survey_definition_immutable before update on crm.survey_definition
 for each row execute function crm.survey_definition_immutable();
commit;
````

### FILE: `db/migrations/0054_customer_feedback.down.sql`
```yaml
block_id: "GO-CUSTOMER-SURVEY-API:db/migrations/0054_customer_feedback.down.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "c2298982d4727d7a35014ccc55be548db77bee3b7f4733c2bd4ba3080795a969"
variables: []
secrets_allowed: false
```
````sql
begin;
drop table crm.survey_response;
drop trigger survey_definition_immutable on crm.survey_definition;
drop function crm.survey_definition_immutable();
drop table crm.survey_definition;

commit;
````

### FILE: `internal/platform/postgres/customerfeedback.go`
```yaml
block_id: "GO-CUSTOMER-SURVEY-API:internal/platform/postgres/customerfeedback.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "9be6e41b375ba1fe28259a167e12de0940b2d2603ec0e218e658f0bbf860d1aa"
variables: []
secrets_allowed: false
```
````go
// AUTHORED PostgreSQL implementation; authentication belongs to the existing HTTP identity boundary.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/customerfeedback"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerFeedback struct{ pool *pgxpool.Pool }

func NewCustomerFeedback(pool *pgxpool.Pool) *CustomerFeedback { return &CustomerFeedback{pool: pool} }

const feedbackEligible = `exists(select 1 from crm.customer_profile c where c.tenant_id=d.tenant_id and c.customer_principal_id=$4 and c.status='active')`

func (r *CustomerFeedback) Definition(ctx context.Context, s customerfeedback.Scope, id string) (customerfeedback.Definition, error) {
	var d customerfeedback.Definition
	e := r.pool.QueryRow(ctx, `select d.survey_id,d.prompt,d.consent_version,d.consent_notice,
 d.active and clock_timestamp()>=d.opens_at and clock_timestamp()<d.closes_at
 from crm.survey_definition d where d.tenant_id=$1 and d.organization_id=$2 and d.survey_id=$3
 and clock_timestamp()<d.retain_until
 and (d.active or exists(select 1 from crm.survey_response a where
 (a.tenant_id,a.organization_id,a.survey_id,a.customer_principal_id)=(d.tenant_id,d.organization_id,d.survey_id,$4)))
 and `+feedbackEligible, s.Tenant, s.Organization, id, s.Customer).
		Scan(&d.ID, &d.Prompt, &d.ConsentVersion, &d.ConsentNotice, &d.Accepting)
	if errors.Is(e, pgx.ErrNoRows) {
		return d, customerfeedback.ErrUnavailable
	}
	return d, e
}
func (r *CustomerFeedback) OwnAnswer(ctx context.Context, s customerfeedback.Scope, id string) (customerfeedback.Answer, error) {
	var a customerfeedback.Answer
	e := r.pool.QueryRow(ctx, `select a.score,a.consent_version,a.received_at from crm.survey_response a
 join crm.survey_definition d using(tenant_id,organization_id,survey_id)
 where a.tenant_id=$1 and a.organization_id=$2 and a.survey_id=$3 and a.customer_principal_id=$4
 and clock_timestamp()<d.retain_until and `+feedbackEligible, s.Tenant, s.Organization, id, s.Customer).
		Scan(&a.Score, &a.ConsentVersion, &a.ReceivedAt)
	if errors.Is(e, pgx.ErrNoRows) {
		return a, customerfeedback.ErrNotFound
	}
	return a, e
}
func (r *CustomerFeedback) Submit(ctx context.Context, s customerfeedback.Scope, id string, in customerfeedback.Submission) (customerfeedback.Result, error) {
	var result customerfeedback.Result
	if e := customerfeedback.ValidateSubmission(s, id, in); e != nil {
		return result, e
	}
	tx, e := r.pool.Begin(ctx)
	if e != nil {
		return result, e
	}
	defer tx.Rollback(ctx)
	// Acquire both definition and customer locks before the final eligibility check.
	// No accepting boolean computed before a lock wait is trusted for admission.
	var version string
	e = tx.QueryRow(ctx, `select d.consent_version
 from crm.survey_definition d join crm.customer_profile c on c.tenant_id=d.tenant_id
 where d.tenant_id=$1 and d.organization_id=$2 and d.survey_id=$3
 and c.customer_principal_id=$4 and c.status='active' for share of d,c`,
		s.Tenant, s.Organization, id, s.Customer).Scan(&version)
	if errors.Is(e, pgx.ErrNoRows) {
		return result, customerfeedback.ErrUnavailable
	}
	if e != nil {
		return result, e
	}
	readExisting := func() error {
		return tx.QueryRow(ctx, `select a.score,a.consent_version,a.received_at
 from crm.survey_response a join crm.survey_definition d using(tenant_id,organization_id,survey_id)
 where a.tenant_id=$1 and a.organization_id=$2 and a.survey_id=$3 and a.customer_principal_id=$4
 and clock_timestamp()<d.retain_until`, s.Tenant, s.Organization, id, s.Customer).
			Scan(&result.Answer.Score, &result.Answer.ConsentVersion, &result.Answer.ReceivedAt)
	}
	checkReplay := func() (customerfeedback.Result, error) {
		if result.Answer.Score != in.Score || result.Answer.ConsentVersion != in.ConsentVersion {
			return result, customerfeedback.ErrConflict
		}
		result.Replay = true
		return result, tx.Commit(ctx)
	}
	// Recovery is a read of the original receipt, including after deactivation/closing,
	// while the customer is still eligible and the retention deadline has not passed.
	e = readExisting()
	if e == nil {
		return checkReplay()
	}
	if !errors.Is(e, pgx.ErrNoRows) {
		return result, e
	}
	e = tx.QueryRow(ctx, `insert into crm.survey_response(tenant_id,organization_id,survey_id,customer_principal_id,score,consent_version)
 select $1,$2,$3,$4,$5,d.consent_version from crm.survey_definition d
 where d.tenant_id=$1 and d.organization_id=$2 and d.survey_id=$3 and d.consent_version=$6
 and d.active and clock_timestamp()>=d.opens_at and clock_timestamp()<d.closes_at
 and clock_timestamp()<d.retain_until
 on conflict do nothing returning score,consent_version,received_at`,
		s.Tenant, s.Organization, id, s.Customer, in.Score, in.ConsentVersion).
		Scan(&result.Answer.Score, &result.Answer.ConsentVersion, &result.Answer.ReceivedAt)
	if errors.Is(e, pgx.ErrNoRows) {
		e = readExisting()
		if errors.Is(e, pgx.ErrNoRows) {
			return result, customerfeedback.ErrUnavailable
		}
		if e != nil {
			return result, e
		}
		return checkReplay()
	}
	if e != nil {
		return result, e
	}
	return result, tx.Commit(ctx)
}
func (r *CustomerFeedback) Summary(ctx context.Context, s customerfeedback.Scope, id string) (customerfeedback.Summary, error) {
	var n, p, d, m int64
	e := r.pool.QueryRow(ctx, `select count(a.customer_principal_id),count(*) filter(where a.score>=9),
 count(*) filter(where a.score<=6),d.minimum_responses
 from crm.survey_definition d left join crm.survey_response a using(tenant_id,organization_id,survey_id)
 where d.tenant_id=$1 and d.organization_id=$2 and d.survey_id=$3 and clock_timestamp()<d.retain_until
 group by d.minimum_responses`, s.Tenant, s.Organization, id).Scan(&n, &p, &d, &m)
	if errors.Is(e, pgx.ErrNoRows) {
		return customerfeedback.Summary{}, customerfeedback.ErrUnavailable
	}
	if e != nil {
		return customerfeedback.Summary{}, e
	}
	return customerfeedback.NPS(n, p, d, m)
}

// PurgeExpired removes only expired responses in one explicit authorized tenant/org.
// Caller must schedule this bounded operation under its approved retention policy.
func (r *CustomerFeedback) PurgeExpired(ctx context.Context, tenant, organization string, limit int) (int64, error) {
	if tenant == "" || organization == "" || limit < 1 || limit > 1000 {
		return 0, customerfeedback.ErrInvalid
	}
	tag, e := r.pool.Exec(ctx, `with doomed as (
 select a.tenant_id,a.organization_id,a.survey_id,a.customer_principal_id
 from crm.survey_response a join crm.survey_definition d using(tenant_id,organization_id,survey_id)
 where a.tenant_id=$1 and a.organization_id=$2 and d.retain_until<=clock_timestamp()
 order by a.survey_id,a.customer_principal_id for update of a skip locked limit $3)
 delete from crm.survey_response a using doomed x where
 (a.tenant_id,a.organization_id,a.survey_id,a.customer_principal_id)=
 (x.tenant_id,x.organization_id,x.survey_id,x.customer_principal_id)`, tenant, organization, limit)
	return tag.RowsAffected(), e
}
````

### FILE: `internal/platform/httpapi/customerfeedback.go`
```yaml
block_id: "GO-CUSTOMER-SURVEY-API:internal/platform/httpapi/customerfeedback.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "626295bc6ab87c14dc48d1faee96058a5e1fdb189e42739e86f750587dbd0ecd"
variables: []
secrets_allowed: false
```
````go
package httpapi

import (
	"context"
	"elite.local/enterprise/internal/customerfeedback"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"time"
)

// CustomerFeedbackModule is opt-in. The existing verifier remains the identity authority.
type CustomerFeedbackModule struct{ Service *customerfeedback.Service }

func (m CustomerFeedbackModule) Register(mux *http.ServeMux, v identity.Verifier) {
	if m.Service == nil {
		return
	}
	a := feedbackAPI{service: m.Service, verifier: v}
	a.registerMetrics(mux)
	mux.HandleFunc("GET /v1/customer/surveys/{survey}", a.definition)
	mux.HandleFunc("POST /v1/customer/surveys/{survey}/response", a.submit)
	mux.HandleFunc("GET /v1/customer/surveys/{survey}/response", a.ownAnswer)
	mux.HandleFunc("GET /v1/admin/surveys/{survey}/summary", a.summary)
}

type feedbackAPI struct {
	service  *customerfeedback.Service
	verifier identity.Verifier
}

func (a feedbackAPI) scope(w http.ResponseWriter, r *http.Request, permission string) (customerfeedback.Scope, bool) {
	p, org, ok := (enterpriseQueryAPI{verifier: a.verifier}).scope(w, r, permission)
	if !ok {
		return customerfeedback.Scope{}, false
	}
	query := r.URL.Query()
	if len(query) != 1 || len(query["organization_id"]) != 1 {
		writeProblem(w, 400, "INVALID_QUERY", "one organization_id is required")
		return customerfeedback.Scope{}, false
	}
	return customerfeedback.Scope{Tenant: p.TenantID, Organization: org, Customer: p.Subject}, true
}
func feedbackError(w http.ResponseWriter, e error) {
	switch {
	case errors.Is(e, customerfeedback.ErrInvalid):
		writeProblem(w, 400, "INVALID_SURVEY_REQUEST", "request does not match the survey contract")
	case errors.Is(e, customerfeedback.ErrConflict):
		writeProblem(w, 409, "SURVEY_RESPONSE_CONFLICT", "a different response is already stored; read your response")
	case errors.Is(e, customerfeedback.ErrUnavailable), errors.Is(e, customerfeedback.ErrNotFound):
		writeProblem(w, 404, "SURVEY_NOT_AVAILABLE", "survey or response is not available")
	default:
		writeProblem(w, 503, "SURVEY_UNAVAILABLE", "survey operation is unavailable; recover with a read")
	}
}
func (a feedbackAPI) definition(w http.ResponseWriter, r *http.Request) {
	s, ok := a.scope(w, r, "customer:self")
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	d, e := a.service.Definition(ctx, s, r.PathValue("survey"))
	if e != nil {
		feedbackError(w, e)
		return
	}
	writeJSON(w, 200, d)
}
func (a feedbackAPI) ownAnswer(w http.ResponseWriter, r *http.Request) {
	s, ok := a.scope(w, r, "customer:self")
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	d, e := a.service.OwnAnswer(ctx, s, r.PathValue("survey"))
	if e != nil {
		feedbackError(w, e)
		return
	}
	writeJSON(w, 200, d)
}
func (a feedbackAPI) summary(w http.ResponseWriter, r *http.Request) {
	s, ok := a.scope(w, r, "surveys:read")
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	d, e := a.service.Summary(ctx, s, r.PathValue("survey"))
	if e != nil {
		feedbackError(w, e)
		return
	}
	writeJSON(w, 200, d)
}
func readSurveySubmission(reader io.Reader) (customerfeedback.Submission, error) {
	var input customerfeedback.Submission
	d := json.NewDecoder(reader)
	token, e := d.Token()
	if e != nil || token != json.Delim('{') {
		return input, customerfeedback.ErrInvalid
	}
	seen := map[string]bool{}
	for d.More() {
		token, e = d.Token()
		key, ok := token.(string)
		if e != nil || !ok || seen[key] {
			return input, customerfeedback.ErrInvalid
		}
		seen[key] = true
		switch key {
		case "score":
			var value *int
			if e = d.Decode(&value); e != nil || value == nil {
				return input, customerfeedback.ErrInvalid
			}
			input.Score = *value
		case "consent":
			var value *bool
			if e = d.Decode(&value); e != nil || value == nil {
				return input, customerfeedback.ErrInvalid
			}
			input.Consent = *value
		case "consent_version":
			var value *string
			if e = d.Decode(&value); e != nil || value == nil {
				return input, customerfeedback.ErrInvalid
			}
			input.ConsentVersion = *value
		default:
			return input, customerfeedback.ErrInvalid
		}
	}
	token, e = d.Token()
	if e != nil || token != json.Delim('}') || len(seen) != 3 {
		return input, customerfeedback.ErrInvalid
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		return input, customerfeedback.ErrInvalid
	}
	return input, nil
}
func (a feedbackAPI) submit(w http.ResponseWriter, r *http.Request) {
	s, ok := a.scope(w, r, "customer:self")
	if !ok {
		return
	}
	kind, _, e := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if e != nil || kind != "application/json" || (r.Header.Get("Content-Encoding") != "" && r.Header.Get("Content-Encoding") != "identity") {
		writeProblem(w, 415, "INVALID_CONTENT_TYPE", "application/json is required")
		return
	}
	input, e := readSurveySubmission(http.MaxBytesReader(w, r.Body, 2048))
	if e != nil {
		feedbackError(w, e)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	result, e := a.service.Submit(ctx, s, r.PathValue("survey"), input)
	if e != nil {
		feedbackError(w, e)
		return
	}
	status := 201
	if result.Replay {
		status = 200
		w.Header().Set("Idempotency-Replayed", "true")
	}
	writeJSON(w, status, result)
}
````

### FILE: `internal/platform/httpapi/customerfeedback_test.go`
```yaml
block_id: "GO-CUSTOMER-SURVEY-API:internal/platform/httpapi/customerfeedback_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "4647abffd3a756221a9e73ad6241d42ba587599961613a6a7e526feba22abdf1"
variables: []
secrets_allowed: false
```
````go
package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFeedbackUnconfiguredModuleHasNoRoutes(t *testing.T) {
	mux := http.NewServeMux()
	CustomerFeedbackModule{}.Register(mux, nil)
	for _, path := range []string{"/v1/customer/surveys/survey", "/v1/customer/surveys/survey/response", "/v1/admin/surveys/survey/summary"} {
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, httptest.NewRequest("GET", path, nil))
		if response.Code != 404 {
			t.Fatal(path, response.Code)
		}
	}
}

func TestFeedbackStrictSubmission(t *testing.T) {
	good := `{"score":0,"consent":true,"consent_version":"v1"}`
	got, e := readSurveySubmission(strings.NewReader(good))
	if e != nil || got.Score != 0 || !got.Consent || got.ConsentVersion != "v1" {
		t.Fatal(got, e)
	}
	for _, body := range []string{
		`{"consent":true,"consent_version":"v1"}`, `{"score":null,"consent":true,"consent_version":"v1"}`,
		`{"score":1.5,"consent":true,"consent_version":"v1"}`, `{"score":1,"score":2,"consent":true,"consent_version":"v1"}`,
		`{"score":1,"\u0073core":2,"consent":true,"consent_version":"v1"}`,
		good + `{}`, good + "false", `[]`, `null`, `{"score":1,"consent":null,"consent_version":"v1"}`,
		`{"score":1,"consent":true,"consent_version":null}`, `{"score":1,"consent":true,"consent_version":"v1","customer":"another"}`,
	} {
		t.Run(body, func(t *testing.T) {
			if _, e := readSurveySubmission(strings.NewReader(body)); e == nil {
				t.Fatal("accepted ambiguous body")
			}
		})
	}
}
````

### FILE: `internal/platform/httpapi/customerfeedback_integration_test.go`
```yaml
block_id: "GO-CUSTOMER-SURVEY-API:internal/platform/httpapi/customerfeedback_integration_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "a82893a0c0ad50cdb6ba4e76fcc4985a7c7328a28cce12c0edccdc1dcefcf1a0"
variables: []
secrets_allowed: false
```
````go
package httpapi

import (
	"context"
	"elite.local/enterprise/internal/customerfeedback"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const feedbackTenant = "018f4d4a-7b36-7a21-8d10-2f4c54c24001"
const feedbackOtherTenant = "018f4d4a-7b36-7a21-8d10-2f4c54c24002"

func feedbackPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is required")
	}
	config, e := pgxpool.ParseConfig(url)
	if e != nil {
		t.Fatal(e)
	}
	if config.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(config.ConnConfig.Database, "elite_feedback_") {
		t.Fatal("requires an owned disposable loopback elite_feedback_* database")
	}
	p, e := pgxpool.NewWithConfig(context.Background(), config)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(p.Close)
	return p
}
func TestFeedbackConnectedPostgres(t *testing.T) {
	pool := feedbackPool(t)
	ctx := context.Background()
	exec := func(q string, args ...any) {
		t.Helper()
		if _, e := pool.Exec(ctx, q, args...); e != nil {
			t.Fatal(e)
		}
	}
	for _, tenant := range []string{feedbackTenant, feedbackOtherTenant} {
		exec(`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Feedback fixture','Feedback fixture')`, tenant, "feedback-"+tenant)
		exec(`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
  values($1,'org-a','a','A','store'),($1,'org-b','b','B','store')`, tenant)
		for _, customer := range []string{"alice", "bob", "carol", "inactive"} {
			state := "active"
			if customer == "inactive" {
				state = "restricted"
			}
			exec(`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized,status)
   values($1,$2,$2,$2||'@example.test',$3)`, tenant, customer, state)
		}
		exec(`insert into crm.survey_definition(tenant_id,organization_id,survey_id,prompt,consent_version,consent_notice,opens_at,closes_at,retain_until,minimum_responses,active)
  values($1,'org-a','satisfaction','Synthetic NPS question','fixture-v1','Synthetic fixture acknowledgement; no real collection',clock_timestamp()-interval '1 hour',clock_timestamp()+interval '1 hour',clock_timestamp()+interval '1 day',2,true)`, tenant)
	}
	verifier, token := confirmationTestIssuer(t)
	service := customerfeedback.NewService(postgres.NewCustomerFeedback(pool))
	mux := http.NewServeMux()
	CustomerFeedbackModule{Service: service}.Register(mux, verifier)
	var drop atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && drop.Swap(false) {
			record := httptest.NewRecorder()
			mux.ServeHTTP(record, r)
			if record.Code != 201 {
				t.Errorf("loss fixture status=%d", record.Code)
			}
			conn, _, e := w.(http.Hijacker).Hijack()
			if e != nil {
				t.Error(e)
				return
			}
			conn.Close()
			return
		}
		mux.ServeHTTP(w, r)
	}))
	defer server.Close()
	client := &http.Client{Timeout: 5 * time.Second}
	customerToken := func(tenant, subject string) string {
		return token(subject, tenant, []string{"customer:self"}, []string{"org-a"})
	}
	alice := customerToken(feedbackTenant, "alice")
	admin := token("operator", feedbackTenant, []string{"surveys:read"}, []string{"org-a"})
	path := "/v1/customer/surveys/satisfaction"
	send := func(method, path, bearer, body string) (int, []byte, error) {
		request, e := http.NewRequest(method, server.URL+path, strings.NewReader(body))
		if e != nil {
			return 0, nil, e
		}
		if bearer != "" {
			request.Header.Set("Authorization", "Bearer "+bearer)
		}
		if method == "POST" {
			request.Header.Set("Content-Type", "application/json")
		}
		response, e := client.Do(request)
		if e != nil {
			return 0, nil, e
		}
		defer response.Body.Close()
		data, e := io.ReadAll(response.Body)
		return response.StatusCode, data, e
	}
	must := func(method, p, bearer, body string, status int) []byte {
		t.Helper()
		code, data, e := send(method, p, bearer, body)
		if e != nil || code != status {
			t.Fatalf("%s %s: status=%d wanted=%d body=%s error=%v", method, p, code, status, data, e)
		}
		return data
	}
	answer := func(score int) string {
		return fmt.Sprintf(`{"score":%d,"consent":true,"consent_version":"fixture-v1"}`, score)
	}
	own := path + "/response?organization_id=org-a"
	must("GET", path+"?organization_id=org-a", "", "", 401)
	must("GET", path+"?organization_id=org-b", alice, "", 403)
	must("GET", path+"?organization_id=org-a&customer_id=bob", alice, "", 400)
	must("GET", path+"?organization_id=org-a", admin, "", 403)
	must("GET", path+"?organization_id=org-a", alice, "", 200)
	must("POST", own, alice, `{"score":9,"consent":false,"consent_version":"fixture-v1"}`, 400)
	must("POST", own, alice, `{"score":9,"consent":true,"consent_version":"wrong"}`, 404)
	must("GET", own, alice, "", 404)
	first := must("POST", own, alice, answer(9), 201)
	var saved customerfeedback.Result
	if e := json.Unmarshal(first, &saved); e != nil {
		t.Fatal(e)
	}
	replay := must("POST", own, alice, answer(9), 200)
	var repeated customerfeedback.Result
	json.Unmarshal(replay, &repeated)
	if !repeated.Replay || !saved.Answer.ReceivedAt.Equal(repeated.Answer.ReceivedAt) {
		t.Fatal("replay changed original")
	}
	must("POST", own, alice, answer(0), 409)
	summaryPath := "/v1/admin/surveys/satisfaction/summary?organization_id=org-a"
	var summary customerfeedback.Summary
	json.Unmarshal(must("GET", summaryPath, admin, "", 200), &summary)
	if summary.Responses != 1 || summary.Available || summary.NPS != nil {
		t.Fatal(summary)
	}
	bob := customerToken(feedbackTenant, "bob")
	results := make(chan int, 24)
	var workers sync.WaitGroup
	for i := 0; i < 24; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			code, _, e := send("POST", own, bob, answer(0))
			if e != nil {
				results <- 0
			} else {
				results <- code
			}
		}()
	}
	workers.Wait()
	close(results)
	created := 0
	for code := range results {
		if code == 201 {
			created++
		} else if code != 200 {
			t.Fatal("concurrent status", code)
		}
	}
	if created != 1 {
		t.Fatal("created", created)
	}
	drop.Store(true)
	if _, _, e := send("POST", own, customerToken(feedbackTenant, "carol"), answer(10)); e == nil {
		t.Fatal("lost response was not lost")
	}
	var recovered customerfeedback.Answer
	json.Unmarshal(must("GET", own, customerToken(feedbackTenant, "carol"), "", 200), &recovered)
	if recovered.Score != 10 {
		t.Fatal(recovered)
	}
	json.Unmarshal(must("GET", summaryPath, admin, "", 200), &summary)
	if summary.Responses != 3 || summary.NPS == nil || *summary.NPS < 33.3333333 || *summary.NPS > 33.3333334 {
		t.Fatal(summary)
	}
	must("GET", summaryPath, alice, "", 403)
	must("POST", own, customerToken(feedbackOtherTenant, "alice"), answer(5), 201)
	must("POST", own, customerToken(feedbackTenant, "inactive"), answer(7), 404)
	if _, e := pool.Exec(ctx, `update crm.survey_definition set prompt='changed' where tenant_id=$1`, feedbackTenant); e == nil {
		t.Fatal("mutable historical question")
	}
	exec(`update crm.survey_definition set active=false where tenant_id=$1`, feedbackTenant)
	must("POST", own, alice, answer(9), 200)
	t.Run("close_window_while_waiting_for_definition_lock", func(t *testing.T) {
		exec(`insert into crm.survey_definition(tenant_id,organization_id,survey_id,prompt,consent_version,consent_notice,opens_at,closes_at,retain_until,minimum_responses,active)
  values($1,'org-a','closing','Synthetic','fixture-v1','Synthetic',clock_timestamp()-interval '1 hour',clock_timestamp()+interval '1500 milliseconds',clock_timestamp()+interval '1 day',1,true)`, feedbackTenant)
		tx, e := pool.Begin(ctx)
		if e != nil {
			t.Fatal(e)
		}
		defer tx.Rollback(ctx)
		if _, e = tx.Exec(ctx, `select 1 from crm.survey_definition where tenant_id=$1 and survey_id='closing' for update`, feedbackTenant); e != nil {
			t.Fatal(e)
		}
		done := make(chan int, 1)
		go func() {
			code, _, _ := send("POST", "/v1/customer/surveys/closing/response?organization_id=org-a", alice, answer(9))
			done <- code
		}()
		deadline := time.Now().Add(1200 * time.Millisecond)
		waiting := false
		for time.Now().Before(deadline) {
			if e = pool.QueryRow(ctx, `select exists(select 1 from pg_stat_activity where datname=current_database() and wait_event_type='Lock' and query like '%crm.survey_definition d%')`).Scan(&waiting); e != nil {
				t.Fatal(e)
			}
			if waiting {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if !waiting {
			t.Fatal("request never waited on definition")
		}
		for {
			var closed bool
			if e = pool.QueryRow(ctx, `select clock_timestamp()>=closes_at from crm.survey_definition where tenant_id=$1 and survey_id='closing'`, feedbackTenant).Scan(&closed); e != nil {
				t.Fatal(e)
			}
			if closed {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if e = tx.Commit(ctx); e != nil {
			t.Fatal(e)
		}
		if code := <-done; code != 404 {
			t.Fatalf("expired while lock held: got %d; want404", code)
		}
		var n int
		if e = pool.QueryRow(ctx, `select count(*) from crm.survey_response where tenant_id=$1 and survey_id='closing'`, feedbackTenant).Scan(&n); e != nil || n != 0 {
			t.Fatal(n, e)
		}
	})
	t.Run("unpublished_definition_is_not_disclosed", func(t *testing.T) {
		_, e := pool.Exec(ctx, `insert into crm.survey_definition(tenant_id,organization_id,survey_id,prompt,consent_version,consent_notice,opens_at,closes_at,retain_until,minimum_responses)values($1,'org-a','draft','Unpublished','v1','Synthetic',clock_timestamp()-interval '1 hour',clock_timestamp()+interval '1 hour',clock_timestamp()+interval '1 day',2)`, feedbackTenant)
		if e != nil {
			t.Fatal(e)
		}
		repository := postgres.NewCustomerFeedback(pool)
		s := customerfeedback.Scope{Tenant: feedbackTenant, Organization: "org-a", Customer: "alice"}
		if _, e = repository.Definition(ctx, s, "draft"); e != customerfeedback.ErrUnavailable {
			t.Fatalf("draft exposed: %v", e)
		}
		if d, e := repository.Definition(ctx, s, "satisfaction"); e != nil || d.Accepting {
			t.Fatalf("closed receipt context lost: %+v %v", d, e)
		}
	})
	t.Log("FEEDBACK_CONNECTED_HTTP_OIDC_POSTGRES: identity, consent, replay,24concurrent requests, lost response, aggregation and close-window fence checked")
}
func TestFeedbackAfterDatabaseRestart(t *testing.T) {
	pool := feedbackPool(t)
	service := customerfeedback.NewService(postgres.NewCustomerFeedback(pool))
	ctx := context.Background()
	s := customerfeedback.Scope{Tenant: feedbackTenant, Organization: "org-a", Customer: "carol"}
	got, e := service.OwnAnswer(ctx, s, "satisfaction")
	if e != nil || got.Score != 10 {
		t.Fatal(got, e)
	}
	stats, e := service.Summary(ctx, s, "satisfaction")
	if e != nil || stats.Responses != 3 || stats.NPS == nil {
		t.Fatal(stats, e)
	}
	for _, customer := range []string{"alice", "bob", "carol"} {
		s.Customer = customer
		if _, e := service.OwnAnswer(ctx, s, "satisfaction"); e != nil {
			t.Fatal(e)
		}
	}
	t.Log("FEEDBACK_POSTGRES_RESTART_PASS responses=3 no_mutations=true")
}
````

### FILE: `internal/platform/httpapi/customerfeedback_browser_test.go`
```yaml
block_id: "GO-CUSTOMER-SURVEY-API:internal/platform/httpapi/customerfeedback_browser_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "cecf6fcd1fc5caf26ff8c4cb3a12724dd1085af2566f32345255e7885c86bdf6"
variables: []
secrets_allowed: false
```
````go
package httpapi

import (
	"context"
	"elite.local/enterprise/internal/customerfeedback"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
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
)

func TestFeedbackBrowserPostgres(t *testing.T) {
	if os.Getenv("ELITE_FEEDBACK_E2E") != "1" {
		t.Skip("explicit disposable survey browser gate not requested")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cfg, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_feedback_") {
		t.Fatal("disposable loopback elite_feedback_* database required")
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'feedback-browser','Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'org-a','org-a','A','store'),($1,'org-b','org-b','B','store')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized,status) values($1,'alice','Synthetic Alice','a@example.test','active'),($1,'bob','Synthetic Bob','b@example.test','active')`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	for _, id := range []string{"chromium-desktop", "chromium-mobile", "firefox-desktop", "webkit-desktop"} {
		_, err = pool.Exec(ctx, `insert into crm.survey_definition(tenant_id,organization_id,survey_id,prompt,consent_version,consent_notice,opens_at,closes_at,retain_until,minimum_responses,active) values($1,'org-a',$2,'Synthetic survey <img src=x onerror=alert(1)>','fixture-v1','Synthetic fixture only',clock_timestamp()-interval '1 hour',clock_timestamp()+interval '1 hour',clock_timestamp()+interval '1 day',2,true)`, tenant, id)
		if err != nil {
			t.Fatal(err)
		}
	}
	verifier, token := confirmationTestIssuer(t)
	identities := map[string]any{}
	for _, name := range []string{"alice", "bob", "admin"} {
		permissions := []string{"customer:self"}
		if name == "admin" {
			permissions = []string{"surveys:read"}
		}
		identities[name] = map[string]any{"subject": name, "tenantId": tenant, "permissions": permissions, "organizations": []string{"org-a"}, "accessToken": token(name, tenant, permissions, []string{"org-a"})}
	}
	mux := http.NewServeMux()
	repository := postgres.NewCustomerFeedback(pool)
	CustomerFeedbackModule{Service: customerfeedback.NewService(repository)}.Register(mux, verifier)
	var writes atomic.Int64
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			writes.Add(1)
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
	web := os.Getenv("ELITE_WEB_ROOT")
	if !filepath.IsAbs(web) {
		t.Fatal("absolute web root required")
	}
	env := []string{}
	for _, item := range os.Environ() {
		name := strings.ToUpper(strings.SplitN(item, "=", 2)[0])
		if strings.HasPrefix(name, "ELITE_") || name == "DATABASE_URL" || name == "TEST_DATABASE_URL" || name == "APP_BASE_URL" || name == "ENTERPRISE_API_BASE_URL" || name == "AUTH_SESSION_SECRET" {
			continue
		}
		env = append(env, item)
	}
	encoded, err := json.Marshal(identities)
	if err != nil {
		t.Fatal(err)
	}
	env = append(env, "ELITE_FEEDBACK_E2E=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_FEEDBACK_IDENTITIES="+string(encoded), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL, "APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-"+randomid.Generator{}.New()+randomid.Generator{}.New())
	artifacts, err := os.MkdirTemp(web, "survey-browser-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	nextlog, err := os.Create(filepath.Join(artifacts, "next.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer nextlog.Close()
	server := exec.CommandContext(ctx, "node", filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", strings.Split(address, ":")[1])
	server.Dir, server.Env = web, env
	server.Stdout, server.Stderr = nextlog, nextlog
	if err = server.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { server.Process.Kill(); server.Wait() }()
	ready := false
	client := &http.Client{Timeout: time.Second}
	for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
		response, e := client.Get(target.String() + "/icon.svg")
		if e == nil {
			response.Body.Close()
			if response.StatusCode == 200 {
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
	command := exec.CommandContext(ctx, "node", filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/customer-survey.spec.mjs", "--timeout=35000", "--output="+filepath.Join(artifacts, "browsers"))
	command.Dir, command.Env = gate, env
	output, err := command.CombinedOutput()
	if e := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); e != nil {
		t.Fatal(e)
	}
	if err != nil || !strings.Contains(string(output), "4 passed") {
		t.Fatalf("browser: %v\n%s", err, output)
	}
	var rows int
	if err = pool.QueryRow(ctx, `select count(*) from crm.survey_response where tenant_id=$1`, tenant).Scan(&rows); err != nil || rows != 8 || writes.Load() != 8 {
		t.Fatalf("responses=%d writes=%d err=%v", rows, writes.Load(), err)
	}
	t.Logf("FEEDBACK_BROWSER_POSTGRES_PASS browsers=4 durable_responses=8 writes=8 artifacts=%s", artifacts)
	t.Run("retention_scope_and_limit", func(t *testing.T) {
		for _, org := range []string{"org-a", "org-b"} {
			if _, e := pool.Exec(ctx, `insert into crm.survey_definition(tenant_id,organization_id,survey_id,prompt,consent_version,consent_notice,opens_at,closes_at,retain_until,minimum_responses)values($1,$2,'expired','Synthetic expired','v1','Synthetic',clock_timestamp()-interval '3 day',clock_timestamp()-interval '2 day',clock_timestamp()-interval '1 day',2)`, tenant, org); e != nil {
				t.Fatal(e)
			}
			if _, e := pool.Exec(ctx, `insert into crm.survey_response(tenant_id,organization_id,survey_id,customer_principal_id,score,consent_version)values($1,$2,'expired','alice',9,'v1'),($1,$2,'expired','bob',0,'v1')`, tenant, org); e != nil {
				t.Fatal(e)
			}
		}
		scope := customerfeedback.Scope{Tenant: tenant, Organization: "org-a", Customer: "alice"}
		if _, e := repository.OwnAnswer(ctx, scope, "expired"); e != customerfeedback.ErrNotFound {
			t.Fatalf("expired answer=%v", e)
		}
		if _, e := repository.Definition(ctx, scope, "expired"); e != customerfeedback.ErrUnavailable {
			t.Fatalf("expired definition=%v", e)
		}
		if _, e := repository.Summary(ctx, scope, "expired"); e != customerfeedback.ErrUnavailable {
			t.Fatalf("expired summary=%v", e)
		}
		for _, limit := range []int{0, 1001} {
			if _, e := repository.PurgeExpired(ctx, tenant, "org-a", limit); e != customerfeedback.ErrInvalid {
				t.Fatalf("limit: %v", e)
			}
		}
		for _, want := range []int64{1, 1, 0} {
			n, e := repository.PurgeExpired(ctx, tenant, "org-a", 1)
			if e != nil || n != want {
				t.Fatalf("purge=%d want=%d err=%v", n, want, e)
			}
		}
		var other, live int
		if e := pool.QueryRow(ctx, `select count(*) filter(where organization_id='org-b'),count(*) filter(where survey_id<>'expired') from crm.survey_response where tenant_id=$1`, tenant).Scan(&other, &live); e != nil || other != 2 || live != 8 {
			t.Fatalf("retention crossed boundary other=%d live=%d err=%v", other, live, e)
		}
		t.Log("FEEDBACK_RETENTION_PASS expired_reads_blocked=true bounded=1 foreign_org_unchanged=true future_unchanged=true")
	})
}
````

### FILE: `internal/platform/httpapi/customerfeedback_fuzz_test.go`
```yaml
block_id: "GO-CUSTOMER-SURVEY-API:internal/platform/httpapi/customerfeedback_fuzz_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "348d0994d660894e4f15836e9efedd8c85d41b17d65d9ab14a129d933be06e01"
variables: []
secrets_allowed: false
```
````go
package httpapi

import (
	"encoding/json"
	"strings"
	"testing"
)

func FuzzFeedbackSubmission(f *testing.F) {
	for _, seed := range []string{
		`{"score":0,"consent":true,"consent_version":"v1"}`,
		`{"score":10,"consent":true,"consent_version":"v1"}`,
		`{"score":-1,"consent":true,"consent_version":"v1"}`,
		`{"score":0,"score":10,"consent":true,"consent_version":"v1"}`,
		`{"score":null,"consent":true,"consent_version":"v1"}`,
		`{"score":0,"consent":false,"consent_version":"v1"}`,
		`{}`, `[]`, `null`, `{"score":1e0}`, `{"scor\u0065":0,"score":1}`,
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		if len(raw) > 2048 {
			return
		}
		got, err := readSurveySubmission(strings.NewReader(raw))
		if err != nil {
			return
		}
		var wire struct {
			Score   int    `json:"score"`
			Consent bool   `json:"consent"`
			Version string `json:"consent_version"`
		}
		if e := json.Unmarshal([]byte(raw), &wire); e != nil {
			t.Fatalf("accepted non-JSON: %v", e)
		}
		if got.Score != wire.Score || got.Consent != wire.Consent || got.ConsentVersion != wire.Version {
			t.Fatal("wire and strict parser disagree")
		}
	})
}
````

### FILE: `cmd/electromobility-api/customerfeedback_activation.go`
```yaml
block_id: "GO-CUSTOMER-SURVEY-API:cmd/electromobility-api/customerfeedback_activation.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "8f0d7811b5e75b56e0c219fbb06a9f41d3e6ff05e4f5c8ed46ba4fae74fad0ef"
variables: []
secrets_allowed: false
```
````go
package main

import (
	"elite.local/enterprise/internal/customerfeedback"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func selectedCustomerFeedbackModule(pool *pgxpool.Pool, getenv func(string) string) (httpapi.EnterpriseModule, error) {
	switch getenv("CUSTOMER_SURVEYS_ENABLED") {
	case "", "0":
		return nil, nil
	case "1":
		if pool == nil {
			return nil, errors.New("customer survey database is missing")
		}
		return httpapi.CustomerFeedbackModule{Service: customerfeedback.NewService(postgres.NewCustomerFeedback(pool))}, nil
	default:
		return nil, errors.New("customer survey activation must be 0 or 1")
	}
}
````

### FILE: `cmd/electromobility-api/customerfeedback_activation_test.go`
```yaml
block_id: "GO-CUSTOMER-SURVEY-API:cmd/electromobility-api/customerfeedback_activation_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "e50299271f65cefc82764643dde517fa7f97aab0f8914b40dfe8309c1db013ba"
variables: []
secrets_allowed: false
```
````go
package main

import "testing"

func TestCustomerFeedbackActivation(t *testing.T) {
	for _, value := range []string{"", "0", "1", "true", "false", " 1", "2"} {
		t.Run(value, func(t *testing.T) {
			module, err := selectedCustomerFeedbackModule(nil, func(key string) string {
				if key != "CUSTOMER_SURVEYS_ENABLED" {
					t.Fatal(key)
				}
				return value
			})
			if module != nil {
				t.Fatal("nil database mounted")
			}
			if (err != nil) != (value != "" && value != "0") {
				t.Fatalf("activation %q: %v", value, err)
			}
		})
	}
}
````

### FILE: `cmd/customer-survey-retention/main.go`
```yaml
block_id: "GO-CUSTOMER-SURVEY-API:cmd/customer-survey-retention/main.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "f492d6493d93e1d79fc85f901541657d35c58184883cf09d6eb52bf7dd7b6c68"
variables: []
secrets_allowed: false
```
````go
// AUTHORED operator command. Scope and retention policy belong to the project.
package main

import (
	"context"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"os"
	"strings"
	"time"
)

func run(args []string, getenv func(string) string, out, stderr io.Writer) int {
	flags := flag.NewFlagSet("customer-survey-retention", flag.ContinueOnError)
	flags.SetOutput(stderr)
	tenant := flags.String("tenant", "", "explicit tenant UUID")
	organization := flags.String("organization", "", "explicit organization ID")
	limit := flags.Int("limit", 0, "required maximum expired responses per invocation (1..1000)")
	if flags.Parse(args) != nil {
		return 2
	}
	if flags.NArg() != 0 || strings.TrimSpace(*tenant) != *tenant || *tenant == "" || strings.TrimSpace(*organization) != *organization || *organization == "" || *limit < 1 || *limit > 1000 {
		fmt.Fprintln(stderr, "tenant, organization and limit 1..1000 are required")
		return 2
	}
	database := getenv("DATABASE_URL")
	if database == "" {
		fmt.Fprintln(stderr, "database configuration is missing")
		return 2
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, database)
	if err != nil {
		fmt.Fprintln(stderr, "database configuration is invalid")
		return 2
	}
	defer pool.Close()
	count, err := postgres.NewCustomerFeedback(pool).PurgeExpired(ctx, *tenant, *organization, *limit)
	if err != nil {
		fmt.Fprintln(stderr, "retention operation failed; inspect authorized database diagnostics")
		return 1
	}
	if json.NewEncoder(out).Encode(map[string]int64{"deleted": count}) != nil {
		return 1
	}
	return 0
}
func main() { os.Exit(run(os.Args[1:], os.Getenv, os.Stdout, os.Stderr)) }
````

### FILE: `cmd/customer-survey-retention/main_test.go`
```yaml
block_id: "GO-CUSTOMER-SURVEY-API:cmd/customer-survey-retention/main_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "744090c0f7f009c4fc0bab147c278d8949f9d408c31046463db417a3f2421552"
variables: []
secrets_allowed: false
```
````go
package main

import (
	"bytes"
	"testing"
)

func TestRetentionRequiresScopeAndBound(t *testing.T) {
	for _, args := range [][]string{nil, {"-tenant", "tenant"}, {"-tenant", "tenant", "-organization", "org"}, {"-tenant", "tenant", "-organization", "org", "-limit", "1001"}, {"-tenant", "tenant", "-organization", "org", "-limit", "1", "extra"}} {
		var out, err bytes.Buffer
		code := run(args, func(string) string { t.Fatal("database consulted before scope validation"); return "" }, &out, &err)
		if code != 2 || out.Len() != 0 {
			t.Fatalf("invalid invocation code=%d output=%s", code, out.String())
		}
	}
}
````

### FILE: `docs/customer-surveys.md`
```yaml
block_id: "GO-CUSTOMER-SURVEY-API:docs/customer-surveys.md:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "af87a4625d3f52d79fe53b5a96a4efc1a286def0e2f27bfa3ac48ab6bd7c8b02"
variables: []
secrets_allowed: false
```
````markdown
# Customer surveys 1.0.0

The connected reference accepts one versioned 0–10 response from an authenticated
active CRM customer, stores it in PostgreSQL, recovers the original receipt and
exposes a separately authorized aggregate. This is AUTHORED code. It replaces the
in-memory store only for this narrow journey; it does not change or certify the
GO-SURVEYS-CORE source.

## Enable and configure

Apply migrations 0001–0054 in numeric order with the existing database migration
procedure. 0054 creates two CRM tables and a definition immutability trigger.
The down migration deletes these tables and their responses: use only on an
empty disposable database, or after the project's approved recovery procedure.

Set CUSTOMER_SURVEYS_ENABLED=1 on the API and customer_surveys=true in the
business configuration features map. Both are disabled when absent.
The BFF must use its existing configured API origin, OIDC session and HTTPS
APP_BASE_URL. A flag does not grant any permission.

An authorized database operator creates a crm.survey_definition with all of:
tenant_id, organization_id, survey_id, prompt, consent_version, consent_notice,
opens_at, closes_at, retain_until, minimum_responses. The default active=false
allows review before setting active=true. Use ASCII identifiers (letters,
digits, underscore or hyphen; initial letter/digit); versions additionally
permit dot and colon. Questions and notices are plain text.
closes_at must follow opens_at, and retain_until must follow closes_at.
Use a new survey_id when terms, question, dates or reporting threshold change.
Only active can be changed on an existing definition.

The project supplies the real question, notice, dates and reporting threshold.
Use an appropriate recommendation question before describing its result as NPS.
The metric does not prove representativeness, consent sufficiency or legal
compliance. No real policy is supplied by the synthetic test fixture.

Customer invitation URL:
/customer/surveys?organizationId=ORG&surveyId=SURVEY

Operator results URL:
/admin/surveys?organizationId=ORG&surveyId=SURVEY

The project distributes these links through its authorized existing workflow.
No outbound messaging, campaign, invitation list or purchase verification is
implemented here. The customer needs customer:self and the exact organization;
the result reader needs surveys:read and the exact organization. API identity
and active CRM membership are checked again on the server.

## Response and recovery

The response key is tenant/organization/survey/customer. Equal resubmission
returns the original receipt; a different score/version conflicts. Consent must
be explicit and match the definition. The database checks the admission window
after acquiring customer/definition locks. Closing or deactivating a survey
prevents new responses but preserves recovery until retention expires.

If the POST outcome is unknown, the UI disables submission and offers a GET-only
recovery button. Reload first reads the existing response. No automatic POST
retry or replacement of a stored answer occurs. The UI shows the actual stored
score, including zero. A new attempt with changed input cannot overwrite it.

All supplied query/body fields are strictly bounded. Customers cannot submit
tenant/customer IDs, read another customer's response or use aggregate routes
without the separate grant. Below the configured minimum, the aggregate returns
the response count and no NPS; this is not a promise of count suppression.

## Retention and operations

Reads stop exposing responses at retain_until. Physical deletion requires the
operator command, with credentials authorized for the selected tenant/org:

    go run ./cmd/customer-survey-retention -tenant TENANT_UUID -organization ORG -limit 100

DATABASE_URL is read only by the process. The command deletes at most 1–1000
expired responses and prints {"deleted":N}. It has a ten-second deadline.
Repeat under the project's approved schedule until deleted=0; failed or
interrupted runs may be repeated because only expired rows are eligible.
No scheduler, backup erasure, legal retention policy or production SLA is
inferred. Record successful runs and failures using the project's monitoring.
Definitions contain no individual answer and are retained separately.

Customer deletion/restriction, backups, audit access, runtime database grants,
load budgets and deployment/restore remain project gates. The reference
database account in tests is synthetic and is not a production role design.

## Evidence and scope

The focal fixture runs real Go HTTP, RS256/JWKS verification, PostgreSQL, Next
JWE sessions and Chromium desktop/mobile, Firefox and WebKit. It verifies
normal submit, duplicate-click fencing, zero score, a lost POST response,
GET-only recovery, permission denial, minimum response threshold and NPS.
Separate probes verify concurrent commits, deadline/lock ordering and actual
database restart. Retention checks bound deletion and preserve other scopes.

The library's broader integration, dependency admission and signed final release
gates remain independent. No 48/48 or production approval follows from this pack.

V402319: the numerical NPS owner is now an explicitly declared PostHog MIT adaptation. See docs/posthog-nps.md and third_party/posthog-nps/source-lock.json. Original caller glue remains AUTHORED; no corporate authorship is assigned to the API, database, UI or privacy policy. The API still returns unrounded float64 with the existing explicit reporting minimum.
````

## 6. Configuration surface

CUSTOMER_SURVEYS_ENABLED=1 explicitly mounts API routes; features.customer_surveys=true
exposes BFF/pages. Both default to disabled. Existing OIDC/session/API-origin config
is reused. Definitions are inactive by default. Terms/dates/thresholds are immutable;
new policy versions require a new survey ID. DATABASE_URL is an operator secret.

## 7. Dependency bill

V400 historical dependency statement (superseded for the selected source in V402319): no dependency, version or external source added. Exact unchanged go.mod/go.sum and
pnpm lock/workspace/package.json come from the selected reference composition.
Go1.26.8, PostgreSQL18.6, Node24.20.0 and frozen pnpm11.25.0 tool identity is recorded
in the V400 evidence; existing dependency/provenance/redistribution conditions remain.

## 8. Apply order

Compose the compatible Go application/query/customer-journey foundations and, for
UI, the matching web/session/BFF and Playwright gate. Apply SQL0001–0054. Configure
project-owned definitions and permissions; enable only after target review.
Disable both flags to roll back exposure. A down migration destroys responses and
must not be used as a data-preserving production rollback. Use the documented
bounded retention command under the project retention/backup policy.

## 9. Verification

Focal HTTP/OIDC/PostgreSQL negatives,24 concurrent submissions, lost response and
actual database restart; production Next browser execution in4projects;8browser
writes produce8responses. Empty migration down/up and bounded retention CLI are
checked. UI typecheck/build and38focal tests PASS. Fuzz and reconstruction receipts
are bound in the V400 record. Readiness, full-domain coverage, runtime monitoring,
load/security/admission and signed final release remain independent gates.

## 10. Reconstruction evidence

See reconstruction_evidence/CONNECTED_CUSTOMER_SURVEYS_V400.md from the library
root. Exact file manifest and hashes bind the rebuilt sources to tested bytes.
This is CONDITIONED for this narrow reference, not 48/48 or target approval.

V400 final correction0.1.1: both SQL migrations are explicit transactions. The original partial-failure fixture leaves2tables; corrected migration rolls back to0. Successful up/down/up passes in a disposable database. No Go/TS behavior or dependency changes; retained original receipts are not relabelled. See V400 atomicity appendix.

V402 composed delta: T2804 role metrics use existing domain read models with exact strings, organization/customer/factory/program permissions, NPS minimum/retention and no zero on unavailable. FAIL868 converted lead and FAIL457 bounded generic body corrected. ROLE_METRICS_RELEASE_V402.md/json; no new dependency or corporate authorship.

### FILE: `internal/platform/httpapi/customerfeedback_metrics.go`

```yaml
block_id: "GO-CUSTOMER-SURVEY-API-METRIC-DELTA:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f8d95a258c391142f6a5db22797cd3fbfb23342506715f4b26c7703c2d2d6144"
variables: []
secrets_allowed: false
```

````go
package httpapi

// AUTHORED reporting wire. Existing Summary retains threshold/retention and calls the admitted PostHog NPS owner.
import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var metricSurveyID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

func (a feedbackAPI) registerMetrics(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/reporting/surveys/{survey}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "private, no-store")
		w.Header().Set("Vary", "Authorization")
		s, ok := a.scope(w, r, "surveys:read")
		if !ok {
			return
		}
		q, e := url.ParseQuery(r.URL.RawQuery)
		if e != nil || len(r.URL.RawQuery) > 256 || len(q) != 1 || !metricSurveyID.MatchString(s.Organization) || !metricSurveyID.MatchString(r.PathValue("survey")) {
			writeProblem(w, 400, "METRIC_INVALID", "bounded scope and survey required")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		v, e := a.service.Summary(ctx, s, r.PathValue("survey"))
		if e != nil {
			feedbackError(w, e)
			return
		}
		writeJSON(w, 200, map[string]any{"organization_id": s.Organization, "survey_id": r.PathValue("survey"), "source": "crm.survey_definition + crm.survey_response", "basis": "retained_survey_population_with_configured_minimum", "responses": strconv.FormatInt(v.Responses, 10), "available": v.Available, "nps": v.NPS, "observed_at": time.Now().UTC()})
	})
}
````


### FILE: `docs/posthog-nps.md`

```yaml
block_id: "GO-CUSTOMER-SURVEY-API-NPS:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2a48b97048c3a11d9ae8a5501c0e980cbff6f619f1cabc89e8bd526f780e0964"
variables: []
secrets_allowed: false
```

````markdown
# PostHog NPS: bounded Go adaptation

This owner translates only calculateNPSFromRawData from PostHog's MIT core,
commit 6fafbb9081bd15e79448af5650e02a4f9ea435cc. It does not import PostHog's
product, enterprise code, dependencies, analytics runtime, or survey service.
Original utils.ts, utils.test.ts and root LICENSE are retained byte-for-byte
as non-executable evidence under third_party/posthog-nps. See source-lock.json.

The original sums 0..6 detractors, 7..8 passives, 9..10 promoters, then computes
(promoters - detractors) / total * 100. Existing SQL groups these same ranges;
the adapted Go owner takes validated grouped counts and derives passives from
total. It preserves zero-population behavior. int64 counts replace JS Number;
subtraction happens before float conversion, preserving the existing Go API.
The original toFixed(1) produces a display string. This adapter deliberately
returns unrounded float64 because changing the public API to one decimal would
be a compatibility change. Six meaningful typed equivalents of the original
nine cases are included; all nine original TS cases were separately executed.

The AUTHORED customerfeedback caller retains overflow-safe validation, explicit
project-owned minimum-response threshold, consent/configuration and access
boundaries. It rejects inconsistent/negative counts before calling Calculate;
the upstream legacy parser is not used. No sampling, legal consent, business
policy, or whole-PostHog certification is implied by this numerical adaptation.
Caller glue, fuzz oracle and source-lock metadata are local work. The algorithm
and six translated examples are ADAPTED MIT with the original notice retained.
No new runtime or module is required. Fixed official provenance is mandatory
before adoption; local fixture qualification does not certify production.
````

### FILE: `third_party/posthog-nps/LICENSE`

```yaml
block_id: "GO-CUSTOMER-SURVEY-API-NPS:file2:v1"
operation: CREATE
provenance: VERBATIM
source: "PostHog/posthog 6fafbb9081bd15e79448af5650e02a4f9ea435cc MIT core; exact source and adaptation boundaries in third_party/posthog-nps/source-lock.json"
license: "MIT"
sha256: "6d82d67dba42eb94ba10f1e986d2eec338c22fb7c5216c2c0ebdecd83d53a029"
variables: []
secrets_allowed: false
```

````text
Copyright (c) 2020-2026 PostHog Inc.

Portions of this software are licensed as follows:

* All content that resides under the "ee/" directory of this repository, if that directory exists, is licensed under the license defined in "ee/LICENSE".
* All third party components incorporated into the PostHog Software are licensed under the original license provided by the owner of the applicable component.
* Content outside of the above mentioned directories or restrictions above is available under the "MIT Expat" license as defined below.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````

### FILE: `third_party/posthog-nps/utils.test.ts.txt`

```yaml
block_id: "GO-CUSTOMER-SURVEY-API-NPS:file3:v1"
operation: CREATE
provenance: VERBATIM
source: "PostHog/posthog 6fafbb9081bd15e79448af5650e02a4f9ea435cc MIT core; exact source and adaptation boundaries in third_party/posthog-nps/source-lock.json"
license: "MIT"
sha256: "7b3b24fe2d6d4525113657629991e0d0056b2c8fcf91e6bd2a906b1c280bd565"
variables: []
secrets_allowed: false
```

````text
import { getAppContext } from 'lib/utils/getAppContext'
import { SurveyRatingResults } from 'scenes/surveys/surveyLogic'
import { urls } from 'scenes/urls'

import {
    EventPropertyFilter,
    FeatureFlagFilters,
    PropertyOperator,
    PropertyFilterType,
    Survey,
    SurveyAppearance,
    SurveyDisplayConditions,
    SurveyEventName,
    SurveyEventProperties,
    SurveyQuestion,
    SurveyQuestionType,
    SurveySchedule,
    SurveyType,
    SurveyWidgetType,
} from '~/types'

import {
    buildAggregateQuery,
    buildOpenEndedQuery,
    buildSurveyExampleInvocationGlobals,
    buildSurveyOptionalBooleanPropertyFilter,
    buildSurveyTimestampFilter,
    calculateNpsBreakdown,
    createAnswerFilterHogQLExpression,
    doesSurveyRepeatOnEveryEvent,
    getExpressionCommentForQuestion,
    getSurveyNotificationFilters,
    getRecurringSurveyScheduleInfo,
    getResolvedSurveyDateRange,
    getSurveyAudienceSummaryValue,
    getSurveyDisplayConditionsSummary,
    getSurveyEndDateForQuery,
    getSurveyResponse,
    getSurveyResponseOutcomeBreakdown,
    getSurveyResponseStatus,
    transformSurveyResponseRows,
    getSurveyStartDateForQuery,
    isSimpleSurveyAudienceTargeting,
    sanitizeColor,
    sanitizeSurvey,
    sanitizeSurveyAppearance,
    sanitizeSurveyDisplayConditions,
    splitChoicesOnPaste,
    surveyEmitsPartialSentEvents,
    validateCSSProperty,
    validateSurveyAppearance,
} from './utils'
import type { SurveyQueryFilters } from './utils'

jest.mock('lib/utils/getAppContext', () => ({
    getAppContext: jest.fn(() => undefined),
}))

const mockedGetAppContext = getAppContext as jest.MockedFunction<typeof getAppContext>

afterEach(() => {
    mockedGetAppContext.mockReturnValue(undefined)
})

describe('survey utils', () => {
    it.each<{ counts: [number, number, number]; percentages: number[] }>([
        { counts: [2, 1, 2], percentages: [0.4, 0.2, 0.4] },
        { counts: [0, 1, 3], percentages: [0, 0.25, 0.75] },
        { counts: [3, 0, 0], percentages: [1, 0, 0] },
        { counts: [0, 0, 0], percentages: [0, 0, 0] },
    ])('calculates response outcome shares for $counts', ({ counts, percentages }) => {
        expect(getSurveyResponseOutcomeBreakdown(counts)).toEqual(
            ['Completed', 'Dismissed', 'Abandoned'].map((label, index) => ({
                label,
                count: counts[index],
                percentage: percentages[index],
            }))
        )
    })

    it.each([
        ['survey sent', { $survey_completed: false }, 'Abandoned'],
        ['survey dismissed', { $survey_partially_completed: true }, 'Dismissed'],
        ['survey abandoned', { $survey_partially_completed: 'true' }, 'Abandoned'],
        ['survey sent', {}, null],
        ['survey dismissed', { $survey_completed: true, $survey_partially_completed: true }, null],
    ])('labels %s using completion and dismissal status', (event, properties, expected) => {
        expect(getSurveyResponseStatus(event, properties)).toBe(expected)
    })

    it.each(['completed', 'abandoned'])('renders merged answers and the %s outcome', (outcome) => {
        const survey = {
            questions: [
                { id: 'rating', type: SurveyQuestionType.Rating },
                { id: 'text', type: SurveyQuestionType.Open },
            ],
        } as Survey
        const rows = [
            {
                result: [
                    [
                        'event-id',
                        'respondent',
                        '2026-09-08T12:00:00Z',
                        'person-id',
                        '{}',
                        JSON.stringify({
                            $survey_id: 'survey-id',
                            $survey_response_text: 'Final answer',
                            $survey_completed: false,
                        }),
                        outcome,
                        ['9', 'Final answer'],
                        SurveyEventName.SENT,
                    ],
                ],
            },
        ]
        const [row] = transformSurveyResponseRows(rows, survey)
        expect(Array.isArray(row.result) ? row.result[0] : null).toMatchObject({
            uuid: 'event-id',
            event: SurveyEventName.SENT,
            properties: {
                $survey_response_rating: '9',
                $survey_response_text: 'Final answer',
                $survey_completed: outcome === 'completed',
                $survey_partially_completed: outcome !== 'completed',
            },
        })
    })

    beforeAll(() => {
        // Mock CSS.supports
        global.CSS = {
            supports: (property: string, value: string): boolean => {
                // Basic color validation - this is a simplified version
                if (property === 'color') {
                    // Helper to validate RGB/HSL number ranges
                    const isValidRGBNumber = (n: string): boolean => {
                        const num = parseInt(n)
                        return !isNaN(num) && num >= 0 && num <= 255
                    }

                    const isValidAlpha = (n: string): boolean => {
                        const num = parseFloat(n)
                        return !isNaN(num) && num >= 0 && num <= 1
                    }

                    // Hex colors (3, 4, 6 or 8 digits)
                    if (value.match(/^#([0-9A-Fa-f]{3}){1,2}$/) || value.match(/^#([0-9A-Fa-f]{4}){1,2}$/)) {
                        return true
                    }

                    // RGB colors
                    const rgbMatch = value.match(/^rgb\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)$/)
                    if (rgbMatch) {
                        return (
                            isValidRGBNumber(rgbMatch[1]) &&
                            isValidRGBNumber(rgbMatch[2]) &&
                            isValidRGBNumber(rgbMatch[3])
                        )
                    }

                    // RGBA colors
                    const rgbaMatch = value.match(/^rgba\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*,\s*([\d.]+)\s*\)$/)
                    if (rgbaMatch) {
                        return (
                            isValidRGBNumber(rgbaMatch[1]) &&
                            isValidRGBNumber(rgbaMatch[2]) &&
                            isValidRGBNumber(rgbaMatch[3]) &&
                            isValidAlpha(rgbaMatch[4])
                        )
                    }

                    // HSL colors
                    if (value.match(/^hsl\(\s*\d+\s*,\s*\d+%\s*,\s*\d+%\s*\)$/)) {
                        return true
                    }

                    // HSLA colors
                    if (value.match(/^hsla\(\s*\d+\s*,\s*\d+%\s*,\s*\d+%\s*,\s*[\d.]+\s*\)$/)) {
                        return true
                    }

                    // Named colors - extend the list with more common colors
                    return ['red', 'blue', 'green', 'transparent', 'black', 'white'].includes(value)
                }
                return false
            },
        } as unknown as typeof CSS
    })

    describe('validateColor', () => {
        it('returns undefined for valid colors in different formats', () => {
            // Hex colors
            expect(validateCSSProperty('color', '#ff0000')).toBeUndefined()
            expect(validateCSSProperty('color', '#f00')).toBeUndefined()
            expect(validateCSSProperty('color', '#ff000080')).toBeUndefined() // With alpha

            // RGB/RGBA colors
            expect(validateCSSProperty('color', 'rgb(255, 0, 0)')).toBeUndefined()
            expect(validateCSSProperty('color', 'rgba(255, 0, 0, 0.5)')).toBeUndefined()

            // HSL/HSLA colors
            expect(validateCSSProperty('color', 'hsl(0, 100%, 50%)')).toBeUndefined()
            expect(validateCSSProperty('color', 'hsla(0, 100%, 50%, 0.5)')).toBeUndefined()

            // Named colors
            expect(validateCSSProperty('color', 'red')).toBeUndefined()
            expect(validateCSSProperty('color', 'transparent')).toBeUndefined()
        })

        it('returns error message for invalid colors', () => {
            expect(validateCSSProperty('color', 'not-a-color')).toBe('not-a-color is not a valid property for color.')
        })

        it('returns undefined for undefined input', () => {
            expect(validateCSSProperty('color', undefined)).toBeUndefined()
        })
    })

    describe('sanitizeColor', () => {
        it('returns undefined for falsy values', () => {
            expect(sanitizeColor(undefined)).toBeUndefined()
            expect(sanitizeColor('')).toBeUndefined()
        })

        it('adds # prefix to valid hex colors without it', () => {
            expect(sanitizeColor('ff0000')).toBe('#ff0000')
            expect(sanitizeColor('123456')).toBe('#123456')
        })

        it('returns original value for already valid colors', () => {
            expect(sanitizeColor('#ff0000')).toBe('#ff0000')
            expect(sanitizeColor('rgb(255, 0, 0)')).toBe('rgb(255, 0, 0)')
            expect(sanitizeColor('red')).toBe('red')
        })
    })

    describe('validateSurveyAppearance', () => {
        const invalidAppearance: SurveyAppearance = {
            backgroundColor: 'not-a-color',
            borderColor: 'also-not-a-color',
            maxWidth: 'definitely-not-a-width',
        }

        it('skips all appearance validation for API surveys', () => {
            // API surveys are rendered by the customer, so PostHog's appearance CSS is not applied
            // and the Customization section is hidden in the editor — flagging errors here would
            // route submitSurveyFailure to a non-existent section and silently block saves.
            expect(validateSurveyAppearance(invalidAppearance, false, SurveyType.API)).toEqual({})
        })

        it('validates appearance CSS for Popover surveys', () => {
            const result = validateSurveyAppearance(invalidAppearance, false, SurveyType.Popover)
            expect(result.backgroundColor).toBe('not-a-color is not a valid property for background-color.')
            expect(result.borderColor).toBe('also-not-a-color is not a valid property for border-color.')
            expect(result.maxWidth).toBe('definitely-not-a-width is not a valid property for width.')
        })
    })

    describe('getSurveyNotificationFilters', () => {
        it('builds survey-specific notification filters', () => {
            expect(getSurveyNotificationFilters('survey-123', true)).toEqual({
                events: [
                    {
                        id: SurveyEventName.SENT,
                        type: 'events',
                        properties: [
                            {
                                key: SurveyEventProperties.SURVEY_ID,
                                type: PropertyFilterType.Event,
                                value: 'survey-123',
                                operator: PropertyOperator.Exact,
                            },
                            {
                                key: SurveyEventProperties.SURVEY_COMPLETED,
                                type: PropertyFilterType.Event,
                                value: true,
                                operator: PropertyOperator.Exact,
                            },
                        ],
                    },
                    {
                        id: SurveyEventName.DISMISSED,
                        type: 'events',
                        properties: [
                            {
                                key: SurveyEventProperties.SURVEY_ID,
                                type: PropertyFilterType.Event,
                                value: 'survey-123',
                                operator: PropertyOperator.Exact,
                            },
                            {
                                key: SurveyEventProperties.SURVEY_PARTIALLY_COMPLETED,
                                type: PropertyFilterType.Event,
                                value: true,
                                operator: PropertyOperator.Exact,
                            },
                        ],
                    },
                ],
            })
        })

        it('also matches a sent event with no completion flag when partial responses are off', () => {
            const sentBranches = getSurveyNotificationFilters('survey-123', false).events?.filter(
                (event) => event.id === SurveyEventName.SENT
            )

            expect(sentBranches).toHaveLength(2)
            expect(sentBranches?.[1].properties).toContainEqual({
                key: SurveyEventProperties.SURVEY_COMPLETED,
                type: PropertyFilterType.Event,
                value: PropertyOperator.IsNotSet,
                operator: PropertyOperator.IsNotSet,
            })
        })
    })

    // An API survey's `survey sent` events come from the integrator's own code, which has no reason
    // to set `$survey_completed` — so requiring it left the notification silently matching nothing.
    describe('surveyEmitsPartialSentEvents', () => {
        it.each([
            [SurveyType.Popover, true, true],
            [SurveyType.Popover, false, false],
            [SurveyType.Widget, true, true],
            [SurveyType.API, true, false],
            [SurveyType.API, false, false],
        ])('%s with partial responses %s', (type, enablePartialResponses, expected) => {
            expect(surveyEmitsPartialSentEvents({ type, enable_partial_responses: enablePartialResponses })).toBe(
                expected
            )
        })
    })

    describe('buildSurveyExampleInvocationGlobals', () => {
        it('builds a survey sent example payload with question response properties', () => {
            const globals = buildSurveyExampleInvocationGlobals({
                survey: {
                    id: 'survey-123',
                    name: 'Onboarding survey',
                    questions: [
                        { id: 'q1', type: SurveyQuestionType.Open, question: 'Tell us more' },
                        {
                            id: 'q2',
                            type: SurveyQuestionType.SingleChoice,
                            question: 'How did you hear about us?',
                            choices: ['Twitter', 'Word of mouth'],
                        },
                        {
                            id: 'q3',
                            type: SurveyQuestionType.MultipleChoice,
                            question: 'What do you use most?',
                            choices: ['Funnels', 'Session replay', 'Feature flags'],
                        },
                        {
                            id: 'q4',
                            type: SurveyQuestionType.Rating,
                            question: 'How satisfied are you?',
                            scale: 10,
                            display: 'number',
                            lowerBoundLabel: 'Low',
                            upperBoundLabel: 'High',
                        },
                    ],
                } as Survey,
                projectId: 1,
                projectName: 'Project',
                projectUrl: 'https://app.posthog.com/project/1',
                timestamp: '2026-04-13T12:00:00.000Z',
                eventUuid: 'event-uuid',
                distinctId: 'person-distinct-id',
            })

            expect(globals.event.event).toEqual(SurveyEventName.SENT)
            expect(globals.event.properties).toEqual({
                $survey_id: 'survey-123',
                $survey_name: 'Onboarding survey',
                $survey_completed: true,
                $survey_submission_id: 'survey-submission-id',
                $survey_response_q1: 'Tell us more',
                $survey_response_q2: 'Twitter',
                $survey_response_q3: ['Funnels', 'Session replay'],
                $survey_response_q4: '9',
            })
        })
    })

    describe('getSurveyResponse', () => {
        it('uses the backend HogQL helper for single-value questions', () => {
            const question = {
                id: 'question-123',
                type: SurveyQuestionType.Rating,
                question: 'How satisfied are you?',
                scale: 10,
                display: 'number',
                lowerBoundLabel: 'Low',
                upperBoundLabel: 'High',
            } as SurveyQuestion

            expect(getSurveyResponse(question, 0)).toBe("getSurveyResponse(0, 'question-123')")
        })

        it('uses the backend HogQL helper for multiple choice questions', () => {
            const question = {
                id: 'question-456',
                type: SurveyQuestionType.MultipleChoice,
                question: 'Which features do you use?',
                choices: ['Insights', 'Session replay'],
            } as SurveyQuestion

            expect(getSurveyResponse(question, 1)).toBe("getSurveyResponse(1, 'question-456', true)")
        })
    })

    describe('getExpressionCommentForQuestion', () => {
        const makeQuestion = (question: string): SurveyQuestion =>
            ({ id: 'q-1', type: SurveyQuestionType.Open, question }) as SurveyQuestion

        it('returns single-line question text unchanged', () => {
            expect(getExpressionCommentForQuestion(makeQuestion('¿Cómo calificarías tu experiencia?'), 0)).toBe(
                '¿Cómo calificarías tu experiencia?'
            )
        })

        it('collapses newlines so multi-line question text stays on a single line', () => {
            // Regression: a newline followed by a non-ASCII char used to leak past the `--` HogQL
            // comment and crash the Survey Results query with "Unexpected character U+00E9".
            const result = getExpressionCommentForQuestion(
                makeQuestion('Queremos compensar tu experiencia.\nDéjanos tu correo y te contactaremos para ayudarte.'),
                4
            )
            expect(result).not.toMatch(/[\r\n]/)
            expect(result).toBe(
                'Queremos compensar tu experiencia. Déjanos tu correo y te contactaremos para ayudarte.'
            )
        })

        it('collapses CRLF and surrounding whitespace', () => {
            expect(getExpressionCommentForQuestion(makeQuestion('line one \r\n  line two'), 0)).toBe(
                'line one line two'
            )
        })

        it('falls back to a positional label when the question is empty or whitespace', () => {
            expect(getExpressionCommentForQuestion(makeQuestion('   '), 2)).toBe('Question 3')
        })
    })

    describe('sanitizeSurveyAppearance', () => {
        it('returns null for null input', () => {
            expect(sanitizeSurveyAppearance(null)).toBeNull()
        })

        it('sanitizes all color fields in the appearance object', () => {
            const input: SurveyAppearance = {
                backgroundColor: 'ff0000',
                borderColor: '00ff00',
                ratingButtonActiveColor: '0000ff',
                ratingButtonColor: 'ffffff',
                submitButtonColor: '000000',
                submitButtonTextColor: 'cccccc',
                // Add other required fields from SurveyAppearance type as needed
            }

            const result = sanitizeSurveyAppearance(input)

            expect(result?.backgroundColor).toBe('#ff0000')
            expect(result?.borderColor).toBe('#00ff00')
            expect(result?.ratingButtonActiveColor).toBe('#0000ff')
            expect(result?.ratingButtonColor).toBe('#ffffff')
            expect(result?.submitButtonColor).toBe('#000000')
            expect(result?.submitButtonTextColor).toBe('#cccccc')
        })

        it('removes surveyPopupDelaySeconds for external surveys', () => {
            const input: SurveyAppearance = {
                backgroundColor: '#ffffff',
                surveyPopupDelaySeconds: 5,
                submitButtonColor: '#000000',
            }

            const result = sanitizeSurveyAppearance(input, false, SurveyType.ExternalSurvey)

            expect(result?.backgroundColor).toBe('#ffffff')
            expect(result?.submitButtonColor).toBe('#000000')
            expect(result?.surveyPopupDelaySeconds).toBeUndefined()
        })

        it('preserves surveyPopupDelaySeconds for non-external surveys', () => {
            const input: SurveyAppearance = {
                backgroundColor: '#ffffff',
                surveyPopupDelaySeconds: 5,
                submitButtonColor: '#000000',
            }

            const result = sanitizeSurveyAppearance(input, false, SurveyType.Popover)

            expect(result?.backgroundColor).toBe('#ffffff')
            expect(result?.submitButtonColor).toBe('#000000')
            expect(result?.surveyPopupDelaySeconds).toBe(5)
        })
    })

    describe('sanitizeSurveyDisplayConditions', () => {
        it('returns null for null input with non-external survey', () => {
            expect(sanitizeSurveyDisplayConditions(null, SurveyType.Popover)).toBeNull()
        })

        it('returns empty conditions object for external surveys with populated input', () => {
            const input: SurveyDisplayConditions = {
                url: 'https://example.com',
                actions: { values: [{ id: 123, name: 'test' }] },
                deviceTypes: ['mobile'],
                seenSurveyWaitPeriodInDays: 7,
                events: { values: [{ name: 'test' }] },
            }

            const result = sanitizeSurveyDisplayConditions(input, SurveyType.ExternalSurvey)

            expect(result).toEqual({
                actions: { values: [] },
                events: { values: [] },
                deviceTypes: undefined,
                deviceTypesMatchType: undefined,
                linkedFlagVariant: undefined,
                seenSurveyWaitPeriodInDays: undefined,
                url: undefined,
                urlMatchType: undefined,
            })
        })

        it('preserves conditions for non-external surveys', () => {
            const input: SurveyDisplayConditions = {
                url: 'https://example.com',
                actions: { values: [{ id: 123, name: 'test' }] },
                events: { values: [{ name: 'test' }] },
                deviceTypes: ['mobile'],
            }

            const result = sanitizeSurveyDisplayConditions(input, SurveyType.Popover)

            expect(result?.url).toBe('https://example.com')
            expect(result?.actions).toEqual({ values: [{ id: 123, name: 'test' }] })
            expect(result?.events).toEqual({ values: [{ name: 'test' }] })
            expect(result?.deviceTypes).toEqual(['mobile'])
        })
    })

    describe('sanitizeSurvey', () => {
        it('sanitizes external survey by removing prohibited fields', () => {
            const inputSurvey = {
                type: SurveyType.ExternalSurvey,
                name: 'Test External Survey',
                questions: [],
                linked_flag_id: 123,
                targeting_flag_filters: { groups: [{ rollout_percentage: 50 }] },
                conditions: {
                    url: 'https://example.com',
                    actions: { values: [{ id: 123, name: 'test' }] },
                    events: { values: [{ name: 'test' }] },
                },
                appearance: {
                    backgroundColor: '#ffffff',
                    surveyPopupDelaySeconds: 5,
                    submitButtonColor: '#000000',
                },
            }

            const result = sanitizeSurvey(inputSurvey)

            // Should remove prohibited fields
            expect(result.linked_flag_id).toBeNull()
            expect(result.targeting_flag_filters).toBeUndefined()
            expect(result.remove_targeting_flag).toBe(true)

            // Should sanitize conditions to empty values
            expect(result.conditions).toEqual({
                actions: { values: [] },
                events: { values: [] },
                deviceTypes: undefined,
                deviceTypesMatchType: undefined,
                linkedFlagVariant: undefined,
                seenSurveyWaitPeriodInDays: undefined,
                url: undefined,
                urlMatchType: undefined,
            })

            // Should remove surveyPopupDelaySeconds from appearance
            expect(result.appearance?.surveyPopupDelaySeconds).toBeUndefined()
            expect(result.appearance?.backgroundColor).toBe('#ffffff')
            expect(result.appearance?.submitButtonColor).toBe('#000000')
        })

        it('preserves fields for non-external surveys', () => {
            const inputSurvey = {
                type: SurveyType.Popover,
                name: 'Test Popover Survey',
                questions: [],
                linked_flag_id: 123,
                targeting_flag_filters: { groups: [{ rollout_percentage: 50 }] },
                conditions: {
                    url: 'https://example.com',
                    actions: { values: [{ id: 123, name: 'test' }] },
                    events: { values: [{ name: 'test' }] },
                },
                appearance: {
                    backgroundColor: '#ffffff',
                    surveyPopupDelaySeconds: 5,
                    submitButtonColor: '#000000',
                },
            }

            const result = sanitizeSurvey(inputSurvey)

            // Should preserve all fields for non-external surveys
            expect(result.linked_flag_id).toBe(123)
            expect(result.targeting_flag_filters).toEqual({ groups: [{ rollout_percentage: 50 }] })
            expect(result.remove_targeting_flag).toBeUndefined()

            // Should preserve conditions
            expect(result.conditions?.url).toBe('https://example.com')
            expect(result.conditions?.actions).toEqual({ values: [{ id: 123, name: 'test' }] })
            expect(result.conditions?.events).toEqual({ values: [{ name: 'test' }] })

            // Should preserve surveyPopupDelaySeconds
            expect(result.appearance?.surveyPopupDelaySeconds).toBe(5)
            expect(result.appearance?.backgroundColor).toBe('#ffffff')
            expect(result.appearance?.submitButtonColor).toBe('#000000')
        })

        it('removes widget-specific fields for non-widget surveys', () => {
            const inputSurvey: Partial<Survey> = {
                type: SurveyType.Popover,
                name: 'Test Survey',
                questions: [],
                appearance: {
                    backgroundColor: '#ffffff',
                    widgetType: SurveyWidgetType.Tab,
                    widgetLabel: 'Feedback',
                    widgetColor: '#ff0000',
                },
            }

            const result = sanitizeSurvey(inputSurvey)

            // Should remove widget-specific fields for non-widget surveys
            expect(result.appearance?.backgroundColor).toBe('#ffffff')
            expect(result.appearance).not.toHaveProperty('widgetType')
            expect(result.appearance).not.toHaveProperty('widgetLabel')
            expect(result.appearance).not.toHaveProperty('widgetColor')
        })

        it('removing conditions object makes it go back to the empty conditions object', () => {
            const inputSurvey = {
                type: SurveyType.ExternalSurvey,
                name: 'Test Survey',
                questions: [],
                conditions: {
                    actions: { values: [] },
                    events: { values: [] },
                },
            }

            const result = sanitizeSurvey(inputSurvey)

            // Should remove empty conditions object
            expect(result.conditions).toEqual({
                actions: {
                    values: [],
                },
                events: {
                    values: [],
                },
                deviceTypes: undefined,
                deviceTypesMatchType: undefined,
                linkedFlagVariant: undefined,
                seenSurveyWaitPeriodInDays: undefined,
                url: undefined,
                urlMatchType: undefined,
            })
        })

        it('Remove conditions key if its value is null', () => {
            const inputSurvey = {
                type: SurveyType.ExternalSurvey,
                name: 'Test Survey',
                questions: [],
                conditions: null,
            }

            const result = sanitizeSurvey(inputSurvey)

            expect(result.conditions).toBeUndefined()
        })

        it('Keep conditions key even if its value is null when option is present', () => {
            const inputSurvey = {
                type: SurveyType.ExternalSurvey,
                name: 'Test Survey',
                questions: [],
                conditions: null,
            }

            const result = sanitizeSurvey(inputSurvey, { keepEmptyConditions: true })

            expect(result.conditions).toBeNull()
        })
    })

    it.each([undefined, 'test-1'])('keeps the linked flag and variant %s in display conditions', (variant) => {
        const survey = {
            linked_flag: { id: 123, key: 'survey-test-flag' },
            conditions: { linkedFlagVariant: variant },
        } as Survey

        expect(getSurveyDisplayConditionsSummary(survey)).toEqual([
            { type: 'flag', label: 'Feature flag', value: 'survey-test-flag', href: urls.featureFlag(123) },
            ...(variant ? [{ type: 'flag_variant', label: 'Variant', value: variant }] : []),
        ])
    })

    describe('audience targeting summaries', () => {
        const baseSurvey = {
            id: 'survey-id',
            created_at: '2024-01-01T00:00:00Z',
            end_date: null,
            conditions: null,
            linked_flag: null,
            linked_flag_id: null,
            targeting_flag: null,
            targeting_flag_filters: undefined,
        } as unknown as Survey

        it('summarizes simple person-property audience rules', () => {
            const survey = {
                ...baseSurvey,
                targeting_flag_filters: {
                    groups: [
                        {
                            properties: [
                                {
                                    key: 'email',
                                    value: ['@posthog.com'],
                                    operator: 'icontains',
                                    type: PropertyFilterType.Person,
                                },
                                {
                                    key: 'plan',
                                    value: ['paid'],
                                    operator: 'exact',
                                    type: PropertyFilterType.Person,
                                },
                            ],
                            rollout_percentage: 100,
                            variant: null,
                        },
                    ],
                },
            } as Survey

            expect(getSurveyAudienceSummaryValue(survey)).toBe('2 audience rules')
            expect(getSurveyDisplayConditionsSummary(survey)).toContainEqual({
                type: 'targeting',
                label: 'Targeting',
                value: '2 audience rules',
            })
        })

        it('supports simple cohort targeting with rollout', () => {
            const survey = {
                ...baseSurvey,
                targeting_flag_filters: {
                    groups: [
                        {
                            properties: [
                                {
                                    key: 'id',
                                    value: 17,
                                    type: PropertyFilterType.Cohort,
                                },
                            ],
                            rollout_percentage: 50,
                            variant: null,
                        },
                    ],
                },
            } as Survey

            expect(getSurveyAudienceSummaryValue(survey)).toBe('1 audience rule · 50% shown')
            expect(isSimpleSurveyAudienceTargeting(survey.targeting_flag_filters)).toBe(true)
        })

        it('summarizes rollout-only targeting', () => {
            const survey = {
                ...baseSurvey,
                targeting_flag_filters: {
                    groups: [
                        {
                            properties: [],
                            rollout_percentage: 50,
                            variant: null,
                        },
                    ],
                },
            } as Survey

            expect(getSurveyAudienceSummaryValue(survey)).toBe('50% of matching users')
        })

        it('detects advanced audience targeting', () => {
            const filters: FeatureFlagFilters = {
                groups: [
                    {
                        properties: [
                            {
                                key: 'email',
                                value: ['@posthog.com'],
                                operator: PropertyOperator.IContains,
                                type: PropertyFilterType.Person,
                            },
                        ],
                        rollout_percentage: 100,
                    },
                    {
                        properties: [],
                        rollout_percentage: 100,
                    },
                ],
            }

            expect(isSimpleSurveyAudienceTargeting(filters)).toBe(false)
            expect(
                getSurveyAudienceSummaryValue({
                    ...baseSurvey,
                    targeting_flag_filters: filters,
                } as Survey)
            ).toBe('Advanced audience targeting')
        })
    })

    describe('calculateNpsBreakdown', () => {
        it('returns all zeros when surveyRatingResults is empty', () => {
            const surveyResults: SurveyRatingResults[number] = {
                data: [],
                total: 0,
            }

            const result = calculateNpsBreakdown(surveyResults)

            expect(result).toBeNull()
        })

        it('returns all zeros when data array is missing', () => {
            const surveyResults: SurveyRatingResults[number] = {
                data: [],
                total: 0,
            }

            const result = calculateNpsBreakdown(surveyResults)

            expect(result).toBeNull()
        })

        it('returns all zeros when data array has incorrect length', () => {
            const surveyResults: SurveyRatingResults[number] = {
                data: [1, 2, 3], // Less than 11 elements
                total: 6,
            }

            const result = calculateNpsBreakdown(surveyResults)

            expect(result).toBeNull()
        })

        it('returns early with all zeros when total is 0', () => {
            const surveyResults: SurveyRatingResults[number] = {
                data: [1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0], // despite having some counts in data
                total: 0, // total is explicitly 0
            }

            const result = calculateNpsBreakdown(surveyResults)

            expect(result).toEqual({
                detractors: 0,
                passives: 0,
                promoters: 0,
                score: '0.0',
                total: 0,
            })
        })

        it('correctly calculates NPS breakdown with all categories present', () => {
            const surveyResults: SurveyRatingResults[number] = {
                data: [
                    1,
                    1,
                    1,
                    1,
                    1,
                    1,
                    1, // 7 detractors (0-6)
                    2,
                    2, // 4 passives (7-8)
                    3,
                    3,
                ], // 6 promoters (9-10)
                total: 17,
            }

            const result = calculateNpsBreakdown(surveyResults)

            expect(result).toEqual({
                detractors: 7,
                passives: 4,
                promoters: 6,
                score: '-5.9',
                total: 17,
            })
        })

        it('handles all zeros', () => {
            const surveyResults: SurveyRatingResults[number] = {
                data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
                total: 0,
            }

            const result = calculateNpsBreakdown(surveyResults)

            expect(result).toEqual({
                detractors: 0,
                passives: 0,
                promoters: 0,
                score: '0.0',
                total: 0,
            })
        })

        it('handles only promoters', () => {
            const surveyResults: SurveyRatingResults[number] = {
                data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5], // only 9s and 10s
                total: 10,
            }

            const result = calculateNpsBreakdown(surveyResults)

            expect(result).toEqual({
                detractors: 0,
                passives: 0,
                promoters: 10,
                score: '100.0',
                total: 10,
            })
        })

        it('handles only passives', () => {
            const surveyResults: SurveyRatingResults[number] = {
                data: [0, 0, 0, 0, 0, 0, 0, 5, 5, 0, 0], // only 7s and 8s
                total: 10,
            }

            const result = calculateNpsBreakdown(surveyResults)

            expect(result).toEqual({
                detractors: 0,
                passives: 10,
                promoters: 0,
                score: '0.0',
                total: 10,
            })
        })

        it('handles only detractors', () => {
            const surveyResults: SurveyRatingResults[number] = {
                data: [2, 2, 2, 2, 2, 2, 2, 0, 0, 0, 0], // only 0-6
                total: 14,
            }

            const result = calculateNpsBreakdown(surveyResults)

            expect(result).toEqual({
                detractors: 14,
                passives: 0,
                promoters: 0,
                score: '-100.0',
                total: 14,
            })
        })
    })

    describe('buildSurveyTimestampFilter', () => {
        it('uses survey default dates when no date range provided', () => {
            const survey = { created_at: '2024-08-27T15:30:00Z', end_date: '2024-08-30T10:00:00Z' }
            const result = buildSurveyTimestampFilter(survey)

            expect(result).toBe(`AND timestamp >= '2024-08-27T00:00:00'
    AND timestamp <= '2024-08-30T23:59:59'`)
        })

        it('respects user date range when provided', () => {
            const survey = { created_at: '2024-08-27T15:30:00Z', end_date: '2024-08-30T10:00:00Z' }
            const dateRange = { date_from: '2024-08-28', date_to: '2024-08-29' }
            const result = buildSurveyTimestampFilter(survey, dateRange)

            expect(result).toBe(`AND timestamp >= '2024-08-28T00:00:00'
    AND timestamp <= '2024-08-29T23:59:59'`)
        })

        it('uses user dates even when before survey creation (no clamping)', () => {
            const survey = { created_at: '2024-08-27T15:30:00Z', end_date: null }
            const dateRange = { date_from: '2024-08-25', date_to: '2024-08-29' } // Earlier than survey creation
            const result = buildSurveyTimestampFilter(survey, dateRange)

            // User's dates are used directly - query will return empty results if no data exists
            expect(result).toContain(`timestamp >= '2024-08-25T00:00:00'`)
            expect(result).toContain(`timestamp <= '2024-08-29T23:59:59'`)
        })

        it('uses team timezone for date boundaries', () => {
            mockedGetAppContext.mockReturnValue({
                current_team: { timezone: 'America/New_York' },
            } as any)

            const survey = { created_at: '2024-08-27T15:30:00Z', end_date: '2024-08-30T10:00:00Z' }
            const dateRange = { date_from: '2024-08-28T12:00:00Z', date_to: '2024-08-29T12:00:00Z' }
            const result = buildSurveyTimestampFilter(survey, dateRange)

            expect(result).toBe(`AND timestamp >= '2024-08-28T00:00:00'
    AND timestamp <= '2024-08-29T23:59:59'`)
        })

        it('defaults to UTC when no team timezone is set', () => {
            mockedGetAppContext.mockReturnValue(undefined)

            const survey = { created_at: '2024-08-27T15:30:00Z', end_date: '2024-08-30T10:00:00Z' }
            const result = buildSurveyTimestampFilter(survey)

            expect(result).toBe(`AND timestamp >= '2024-08-27T00:00:00'
    AND timestamp <= '2024-08-30T23:59:59'`)
        })

        it('handles date_to with time component from date picker', () => {
            const survey = { created_at: '2024-08-27T15:30:00Z', end_date: '2024-08-30T10:00:00Z' }
            // Date picker provides date_to with T23:59:59
            const dateRange = { date_from: '2024-08-28', date_to: '2024-08-28T23:59:59' }
            const result = buildSurveyTimestampFilter(survey, dateRange)

            expect(result).toContain(`timestamp >= '2024-08-28T00:00:00'`)
            expect(result).toContain(`timestamp <= '2024-08-28T23:59:59'`)
        })

        it('uses survey defaults when only date_from provided', () => {
            const survey = { created_at: '2024-08-27T15:30:00Z', end_date: '2024-08-30T10:00:00Z' }
            const dateRange = { date_from: '2024-08-28', date_to: null }
            const result = buildSurveyTimestampFilter(survey, dateRange)

            expect(result).toContain(`timestamp >= '2024-08-28T00:00:00'`)
            expect(result).toContain(`timestamp <= '2024-08-30T23:59:59'`) // Survey end date
        })

        it('ignores date_to when date_from not provided (avoids impossible ranges)', () => {
            const survey = { created_at: '2024-08-27T15:30:00Z', end_date: '2024-08-30T10:00:00Z' }
            const dateRange = { date_from: null, date_to: '2024-08-29' }
            const result = buildSurveyTimestampFilter(survey, dateRange)

            // Uses survey defaults since date_to only could create impossible range
            expect(result).toContain(`timestamp >= '2024-08-27T00:00:00'`) // Survey start date
            expect(result).toContain(`timestamp <= '2024-08-30T23:59:59'`) // Survey end date
        })

        it('prefers survey start_date over created_at for the lower bound', () => {
            const survey = {
                created_at: '2024-08-20T15:30:00Z',
                start_date: '2024-08-27T09:00:00Z',
                end_date: '2024-08-30T10:00:00Z',
            }
            const result = buildSurveyTimestampFilter(survey)

            expect(result).toContain(`timestamp >= '2024-08-27T00:00:00'`)
        })
    })

    describe('getResolvedSurveyDateRange', () => {
        it('does not shift dates due to timezone conversion', () => {
            const survey = { created_at: '2024-11-19T00:00:00Z', end_date: '2024-11-25T00:00:00Z' }
            // This datetime should NOT be shifted to Nov 21 due to local timezone conversion
            const dateRange = { date_from: '2024-11-20', date_to: '2024-11-20T23:59:59' }

            const result = getResolvedSurveyDateRange(survey, dateRange)

            expect(result.fromDate).toBe('2024-11-20T00:00:00')
            expect(result.toDate).toBe('2024-11-20T23:59:59')
        })
    })

    describe('submission merging in the results queries', () => {
        const buildSurvey = (enablePartialResponses: boolean): Survey =>
            ({
                id: 'test-survey-id',
                created_at: '2024-11-19T00:00:00Z',
                end_date: null,
                enable_partial_responses: enablePartialResponses,
                questions: [
                    { id: 'q-rating', type: SurveyQuestionType.Rating, question: 'How was it?' },
                    { id: 'q-open', type: SurveyQuestionType.Open, question: 'Why?' },
                    {
                        id: 'q-multi',
                        type: SurveyQuestionType.MultipleChoice,
                        question: 'Which ones?',
                        choices: ['a', 'b'],
                    },
                ],
            }) as Survey

        const buildFilters = (survey: Survey, overrides: Partial<SurveyQueryFilters> = {}): SurveyQueryFilters => ({
            timestampFilter: buildSurveyTimestampFilter(survey),
            answerFilters: [],
            archivedResponsesFilter: '',
            ...overrides,
        })

        it.each([
            ['rating', 0, 'isNotNull(q0_raw)'],
            ['open', 1, 'isNotNull(q1_raw)'],
            // Multiple-choice answers are arrays, so an `isNotNull` merge condition would be true on
            // every event and re-elect the latest one, dropping choices made on an earlier event.
            ['multiple choice', 2, 'length(q2_raw) > 0'],
        ])('merges the %s answer across the submission with argMaxIf', (_type, index, presenceExpr) => {
            const survey = buildSurvey(true)

            const query = buildAggregateQuery(survey, buildFilters(survey))

            expect(query).toContain(
                `argMaxIf(q${index}_raw, tuple(timestamp, event_uuid), ${presenceExpr}) AS q${index}_answer`
            )
            expect(query).toContain('GROUP BY submission_key')
        })

        it.each([true, false])('includes captured answers with partial collection set to %s', (enabled) => {
            const survey = buildSurvey(enabled)
            for (const query of [
                buildAggregateQuery(survey, buildFilters(survey)),
                buildOpenEndedQuery(survey, buildFilters(survey))?.query,
            ]) {
                expect(query).toContain("event = 'survey sent'")
                expect(query).toContain("'survey dismissed', 'survey abandoned'")
                expect(query).toContain(SurveyEventProperties.SURVEY_PARTIALLY_COMPLETED)
                expect(query).toContain('GROUP BY submission_key')
                expect(query).not.toContain('HAVING countIf(is_completed_event) > 0')
            }
        })

        it('applies answer and archive filters to the merged answer, not to single events', () => {
            const survey = buildSurvey(true)
            const filters = buildFilters(survey, {
                answerFilters: [
                    {
                        type: PropertyFilterType.Event,
                        key: '$survey_response_q-rating',
                        operator: PropertyOperator.Exact,
                        value: '2',
                    } as EventPropertyFilter,
                ],
                archivedResponsesFilter: "AND uuid NOT IN ('archived-uuid')",
            })

            const query = buildAggregateQuery(survey, filters)

            // Filtering on the raw event expression would discard a submission whose matching
            // answer arrived on a non-final event.
            expect(query).toContain("(q0_answer = '2')")
            expect(query).not.toContain("getSurveyResponse(0, 'q-rating') = '2'")
            expect(query).toContain("uuid NOT IN ('archived-uuid')")
        })

        it('counts an unanswered optional single choice when the merge yields null instead of an empty string', () => {
            const survey = {
                ...buildSurvey(true),
                questions: [
                    {
                        id: 'q-choice',
                        type: SurveyQuestionType.SingleChoice,
                        question: 'Pick one',
                        choices: ['a', 'b'],
                        optional: true,
                    },
                ],
            } as Survey

            const query = buildAggregateQuery(survey, buildFilters(survey))

            // argMaxIf returns the type default when no event answered the question, so the old
            // `= ''` test silently missed those submissions.
            expect(query).toContain("length(trim(coalesce(q0_answer, ''))) = 0")
        })

        it('reads the merged submissions once rather than once per question', () => {
            const survey = buildSurvey(true)

            const query = buildAggregateQuery(survey, buildFilters(survey))

            // ClickHouse inlines a CTE instead of materializing it, so counting each question in
            // its own UNION ALL branch re-runs the whole merge per branch. Measured at roughly
            // twice the runtime on a four-question survey before this collapsed to one arrayJoin.
            expect(query).not.toContain('UNION ALL')
            expect(query!.match(/argMaxIf\(q0_raw/g)).toHaveLength(1)
        })

        it.each([
            ['rating', { id: 'q-rating', type: SurveyQuestionType.Rating, question: 'How was it?' }],
            ['open', { id: 'q-open', type: SurveyQuestionType.Open, question: 'Why?' }],
            [
                'single choice',
                {
                    id: 'q-choice',
                    type: SurveyQuestionType.SingleChoice,
                    question: 'Pick one',
                    choices: ['a', 'b'],
                },
            ],
        ])('builds a valid query for a survey with only one required %s question', (_type, question) => {
            const survey = { ...buildSurvey(true), questions: [question] } as Survey

            const query = buildAggregateQuery(survey, buildFilters(survey))

            // These questions each emit one label-pair expression, and HogQL rejects arrayConcat
            // with a single argument, so the results tab failed to load.
            expect(query).not.toContain('arrayConcat')
            expect(query).toContain('arrayJoin(if(isNotNull(q0_answer)')
        })

        it('concatenates the label pairs when a survey emits more than one expression', () => {
            const survey = buildSurvey(true)

            const query = buildAggregateQuery(survey, buildFilters(survey))

            expect(query).toContain('arrayJoin(arrayConcat(')
        })

        it('does not alias the merged timestamp back onto the column the merge orders by', () => {
            const survey = buildSurvey(true)

            const query = buildAggregateQuery(survey, buildFilters(survey))

            // `max(timestamp) AS timestamp` makes every sibling `argMax(..., timestamp)` resolve
            // its ordering argument to that aggregate, and ClickHouse rejects the nesting with
            // "Aggregate function ... is found inside another aggregate function".
            expect(query).toContain('max(timestamp) AS submitted_at')
            expect(query).not.toContain('max(timestamp) AS timestamp')
        })

        it('keeps respondent metadata after the open columns so positional parsing still lines up', () => {
            const survey = buildSurvey(true)

            const result = buildOpenEndedQuery(survey, buildFilters(survey))

            const openColumnIndex = result!.query.indexOf('q1_answer AS q1_response')
            expect(openColumnIndex).toBeGreaterThan(-1)
            expect(
                result!.query.indexOf('distinct_id,\n            submitted_at,\n            session_id')
            ).toBeGreaterThan(openColumnIndex)
            expect(result!.columnMap['q-open']).toEqual({
                columnIndex: 0,
                questionIndex: 1,
                type: SurveyQuestionType.Open,
            })
        })
    })

    describe('buildSurveyOptionalBooleanPropertyFilter', () => {
        it('builds a null-safe comparison for optional survey booleans', () => {
            expect(
                buildSurveyOptionalBooleanPropertyFilter(SurveyEventProperties.SURVEY_PARTIALLY_COMPLETED, 'true')
            ).toBe(`coalesce(JSONExtractString(properties, '$survey_partially_completed'), '') != 'true'`)
        })
    })
})

describe('createAnswerFilterHogQLExpression', () => {
    const mockSurvey = {
        questions: [{ id: 'q1' }, { id: 'q2' }, { id: 'q3' }],
    } as any as Survey

    it('returns empty string for empty filters array', () => {
        expect(createAnswerFilterHogQLExpression([], mockSurvey)).toBe('')
    })

    it('returns empty string for null or undefined filters', () => {
        expect(createAnswerFilterHogQLExpression(null as any, mockSurvey)).toBe('')
        expect(createAnswerFilterHogQLExpression(undefined as any, mockSurvey)).toBe('')
    })

    it('handles single exact filter', () => {
        const filters = [
            { key: '$survey_response_q1', value: 'yes', operator: 'exact', type: PropertyFilterType.Event },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(`AND (${getSurveyResponse(mockSurvey.questions[0], 0)} = 'yes')`)
    })

    it('handles filter for a different question', () => {
        const filters = [
            { key: '$survey_response_q2', value: 'no', operator: 'exact', type: PropertyFilterType.Event },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(`AND (${getSurveyResponse(mockSurvey.questions[1], 1)} = 'no')`)
    })

    it('skips filters with empty values', () => {
        const filters = [
            { key: '$survey_response_q1', value: '', operator: 'exact', type: PropertyFilterType.Event },
            { key: '$survey_response_q2', value: null, operator: 'exact', type: PropertyFilterType.Event },
            { key: '$survey_response_q3', value: undefined, operator: 'exact', type: PropertyFilterType.Event },
        ] as EventPropertyFilter[]

        expect(createAnswerFilterHogQLExpression(filters, mockSurvey)).toBe('')
    })

    it('skips filters with empty arrays', () => {
        const filters = [
            { key: '$survey_response_q1', value: [], operator: 'exact', type: PropertyFilterType.Event },
        ] as EventPropertyFilter[]

        expect(createAnswerFilterHogQLExpression(filters, mockSurvey)).toBe('')
    })

    it('skips icontains filters with empty search patterns', () => {
        const filters = [
            { key: '$survey_response_q1', value: '%', operator: 'icontains', type: PropertyFilterType.Event },
            { key: '$survey_response_q2', value: '%%', operator: 'icontains', type: PropertyFilterType.Event },
            { key: '$survey_response_q3', value: '   ', operator: 'icontains', type: PropertyFilterType.Event },
        ] as EventPropertyFilter[]

        expect(createAnswerFilterHogQLExpression(filters, mockSurvey)).toBe('')
    })

    it('handles exact operator with single value', () => {
        const filters = [
            { key: '$survey_response_q1', value: 'test', operator: 'exact', type: PropertyFilterType.Event },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(`AND (${getSurveyResponse(mockSurvey.questions[0], 0)} = 'test')`)
    })

    it('handles exact operator with array values', () => {
        const filters = [
            {
                key: '$survey_response_q1',
                value: ['option1', 'option2'],
                operator: 'exact',
                type: PropertyFilterType.Event,
            },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(`AND (${getSurveyResponse(mockSurvey.questions[0], 0)} IN ('option1', 'option2'))`)
    })

    it('handles is_not operator with single value', () => {
        const filters = [
            { key: '$survey_response_q1', value: 'test', operator: 'is_not', type: PropertyFilterType.Event },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(`AND (${getSurveyResponse(mockSurvey.questions[0], 0)} != 'test')`)
    })

    it('handles is_not operator with array values', () => {
        const filters = [
            {
                key: '$survey_response_q1',
                value: ['option1', 'option2'],
                operator: 'is_not',
                type: PropertyFilterType.Event,
            },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(`AND (${getSurveyResponse(mockSurvey.questions[0], 0)} NOT IN ('option1', 'option2'))`)
    })

    it('handles icontains operator', () => {
        const filters = [
            { key: '$survey_response_q1', value: 'search', operator: 'icontains', type: PropertyFilterType.Event },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(`AND (${getSurveyResponse(mockSurvey.questions[0], 0)} ILIKE '%search%')`)
    })

    it('handles not_icontains operator', () => {
        const filters = [
            { key: '$survey_response_q1', value: 'search', operator: 'not_icontains', type: PropertyFilterType.Event },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(`AND (NOT ${getSurveyResponse(mockSurvey.questions[0], 0)} ILIKE '%search%')`)
    })

    it('handles regex operator', () => {
        const filters = [
            { key: '$survey_response_q1', value: '.*test.*', operator: 'regex', type: PropertyFilterType.Event },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(`AND (match(${getSurveyResponse(mockSurvey.questions[0], 0)}, '.*test.*'))`)
    })

    it('handles not_regex operator', () => {
        const filters = [
            { key: '$survey_response_q1', value: '.*test.*', operator: 'not_regex', type: PropertyFilterType.Event },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(`AND (NOT match(${getSurveyResponse(mockSurvey.questions[0], 0)}, '.*test.*'))`)
    })

    it('combines multiple filters with AND', () => {
        const filters = [
            { key: '$survey_response_q1', value: 'yes', operator: 'exact', type: PropertyFilterType.Event },
            { key: '$survey_response_q2', value: 'no', operator: 'exact', type: PropertyFilterType.Event },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(
            `AND (${getSurveyResponse(mockSurvey.questions[0], 0)} = 'yes') AND (${getSurveyResponse(
                mockSurvey.questions[1],
                1
            )} = 'no')`
        )
    })

    it('skips filters with invalid question keys', () => {
        const filters = [
            { key: '$survey_response_invalid', value: 'test', operator: 'exact', type: PropertyFilterType.Event },
            { key: '$survey_response_q4', value: 'test2', operator: 'exact', type: PropertyFilterType.Event }, // q4 doesn't exist in mockSurvey
        ] as EventPropertyFilter[]

        expect(createAnswerFilterHogQLExpression(filters, mockSurvey)).toBe('')
    })

    it('handles array values for regex and not_regex operators', () => {
        const filters = [
            { key: '$survey_response_q1', value: ['.*pattern.*'], operator: 'regex', type: PropertyFilterType.Event },
            {
                key: '$survey_response_q2',
                value: ['.*pattern.*'],
                operator: 'not_regex',
                type: PropertyFilterType.Event,
            },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(
            `AND (match(${getSurveyResponse(
                mockSurvey.questions[0],
                0
            )}, '.*pattern.*')) AND (NOT match(${getSurveyResponse(mockSurvey.questions[1], 1)}, '.*pattern.*'))`
        )
    })

    it('handles array values for icontains operator', () => {
        const filters = [
            {
                key: '$survey_response_q1',
                value: ['searchterm'],
                operator: 'icontains',
                type: PropertyFilterType.Event,
            },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(`AND (${getSurveyResponse(mockSurvey.questions[0], 0)} ILIKE '%searchterm%')`)
    })

    it('handles unsupported operators', () => {
        const filters = [
            {
                key: '$survey_response_q1',
                value: "O'Reilly",
                operator: 'exact',
                type: PropertyFilterType.Event,
            },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(`AND (${getSurveyResponse(mockSurvey.questions[0], 0)} = 'O\\'Reilly')`)
    })

    it('escapes backslashes in values', () => {
        const filters = [
            {
                key: '$survey_response_q1',
                value: 'C:\\\\path\\\\to\\\\file',
                operator: 'exact',
                type: PropertyFilterType.Event,
            },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(
            `AND (${getSurveyResponse(mockSurvey.questions[0], 0)} = 'C:\\\\\\\\path\\\\\\\\to\\\\\\\\file')`
        )
    })

    it('escapes SQL injection attempts in array values', () => {
        const filters = [
            {
                key: '$survey_response_q1',
                value: ['normal', "'; DROP TABLE users; --", "Robert'); DROP TABLE students; --"],
                operator: 'exact',
                type: PropertyFilterType.Event,
            },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(
            `AND (${getSurveyResponse(
                mockSurvey.questions[0],
                0
            )} IN ('normal', '\\'; DROP TABLE users; --', 'Robert\\'); DROP TABLE students; --'))`
        )
    })

    it('escapes complex SQL injection patterns', () => {
        const filters = [
            {
                key: '$survey_response_q1',
                value: "' UNION SELECT * FROM users; --",
                operator: 'exact',
                type: PropertyFilterType.Event,
            },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(
            `AND (${getSurveyResponse(mockSurvey.questions[0], 0)} = '\\' UNION SELECT * FROM users; --')`
        )
    })

    it('handles regex patterns with special characters', () => {
        const filters = [
            {
                key: '$survey_response_q1',
                value: ".*'; DROP TABLE.*",
                operator: 'regex',
                type: PropertyFilterType.Event,
            },
        ] as EventPropertyFilter[]

        const result = createAnswerFilterHogQLExpression(filters, mockSurvey)
        expect(result).toBe(`AND (match(${getSurveyResponse(mockSurvey.questions[0], 0)}, '.*\\'; DROP TABLE.*'))`)
    })

    describe('multiple choice questions', () => {
        const surveyWithMultipleChoiceQuestion = {
            ...mockSurvey,
            questions: [
                {
                    ...mockSurvey.questions[0],
                    type: SurveyQuestionType.MultipleChoice,
                    choices: [
                        { id: 'c1', label: 'test' },
                        { id: 'c2', label: 'test2' },
                    ],
                },
            ],
        } as any as Survey

        it('handles icontains operator', () => {
            const filters = [
                { key: '$survey_response_q1', value: 'test', operator: 'icontains', type: PropertyFilterType.Event },
            ] as EventPropertyFilter[]

            const result = createAnswerFilterHogQLExpression(filters, surveyWithMultipleChoiceQuestion)
            expect(result).toBe(
                `AND (arrayExists(x -> x ilike '%test%', ${getSurveyResponse(surveyWithMultipleChoiceQuestion.questions[0], 0)}))`
            )
        })

        it('handles not_icontains operator for multiple choice question', () => {
            const filters = [
                {
                    key: '$survey_response_q1',
                    value: 'test',
                    operator: 'not_icontains',
                    type: PropertyFilterType.Event,
                },
            ] as EventPropertyFilter[]

            const result = createAnswerFilterHogQLExpression(filters, surveyWithMultipleChoiceQuestion)
            expect(result).toBe(
                `AND (NOT arrayExists(x -> x ilike '%test%', ${getSurveyResponse(surveyWithMultipleChoiceQuestion.questions[0], 0)}))`
            )
        })

        it('handles regex operator', () => {
            const filters = [
                { key: '$survey_response_q1', value: '.*test.*', operator: 'regex', type: PropertyFilterType.Event },
            ] as EventPropertyFilter[]

            const result = createAnswerFilterHogQLExpression(filters, surveyWithMultipleChoiceQuestion)
            expect(result).toBe(
                `AND (arrayExists(x -> match(x, '.*test.*'), ${getSurveyResponse(surveyWithMultipleChoiceQuestion.questions[0], 0)}))`
            )
        })

        it('handles not_regex operator', () => {
            const filters = [
                {
                    key: '$survey_response_q1',
                    value: '.*test.*',
                    operator: 'not_regex',
                    type: PropertyFilterType.Event,
                },
            ] as EventPropertyFilter[]

            const result = createAnswerFilterHogQLExpression(filters, surveyWithMultipleChoiceQuestion)
            expect(result).toBe(
                `AND (NOT arrayExists(x -> match(x, '.*test.*'), ${getSurveyResponse(surveyWithMultipleChoiceQuestion.questions[0], 0)}))`
            )
        })
    })
})

describe('timezone handling in survey date queries', () => {
    const createMockSurvey = (
        createdAt: string,
        endDate?: string,
        startDate?: string
    ): Pick<Survey, 'created_at' | 'end_date'> & Partial<Pick<Survey, 'start_date'>> => ({
        created_at: createdAt,
        end_date: endDate || null,
        start_date: startDate,
    })

    afterEach(() => {
        mockedGetAppContext.mockReset()
    })

    it('uses team timezone to compute date boundaries', () => {
        mockedGetAppContext.mockReturnValue({
            current_team: { timezone: 'Asia/Tokyo' },
        } as any)

        // 2024-08-27T15:30:00Z = 2024-08-28T00:30:00 JST
        const survey = createMockSurvey('2024-08-27T15:30:00Z', '2024-08-30T10:00:00Z')

        const startDate = getSurveyStartDateForQuery(survey)
        const endDate = getSurveyEndDateForQuery(survey)

        // In JST (UTC+9), the created_at falls on Aug 28, not Aug 27
        expect(startDate).toBe('2024-08-28T00:00:00')
        expect(endDate).toBe('2024-08-30T23:59:59')
    })

    it('defaults to UTC when no team context', () => {
        mockedGetAppContext.mockReturnValue(undefined)

        const survey = createMockSurvey('2024-08-27T15:30:00Z', '2024-08-30T10:00:00Z')

        const startDate = getSurveyStartDateForQuery(survey)
        const endDate = getSurveyEndDateForQuery(survey)

        expect(startDate).toBe('2024-08-27T00:00:00')
        expect(endDate).toBe('2024-08-30T23:59:59')
    })

    it('handles null end_date correctly', () => {
        mockedGetAppContext.mockReturnValue({
            current_team: { timezone: 'America/Chicago' },
        } as any)

        const survey = createMockSurvey('2024-08-27T15:30:00Z')
        const result = getSurveyEndDateForQuery(survey)

        expect(result).toMatch(/^\d{4}-\d{2}-\d{2}T23:59:59$/)
    })
})

describe('splitChoicesOnPaste', () => {
    it('returns null when only one segment is pasted', () => {
        expect(splitChoicesOnPaste('single value', [''], 0, false)).toBeNull()
        expect(splitChoicesOnPaste('Yes, sometimes', [''], 0, false)).toBeNull()
    })

    it('splits newline-separated values into the choices array', () => {
        expect(splitChoicesOnPaste('one\ntwo\nthree', [''], 0, false)).toEqual(['one', 'two', 'three'])
    })

    it('splits tab-separated values (spreadsheet rows)', () => {
        expect(splitChoicesOnPaste('one\ttwo\tthree', [''], 0, false)).toEqual(['one', 'two', 'three'])
    })

    it('trims and drops empty segments', () => {
        expect(splitChoicesOnPaste('  one  \n\n  two  \n', [''], 0, false)).toEqual(['one', 'two'])
    })

    it('inserts segments in place of the target choice and keeps surrounding choices', () => {
        expect(splitChoicesOnPaste('two\nthree', ['one', 'placeholder', 'four'], 1, false)).toEqual([
            'one',
            'two',
            'three',
            'four',
        ])
    })

    it('preserves the open-ended "Other" entry when pasting into a regular slot', () => {
        expect(splitChoicesOnPaste('two\nthree', ['one', '', 'Other'], 1, true)).toEqual([
            'one',
            'two',
            'three',
            'Other',
        ])
    })

    it('preserves the open-ended "Other" entry when pasting into the open-ended slot itself', () => {
        expect(splitChoicesOnPaste('two\nthree', ['one', 'Other'], 1, true)).toEqual(['one', 'two', 'three', 'Other'])
    })
})

describe('doesSurveyRepeatOnEveryEvent', () => {
    it.each([
        [
            'repeated activation with a trigger event',
            true,
            { values: [{ name: 'purchase' }], repeatedActivation: true },
        ],
        ['repeated activation without trigger events', false, { values: [], repeatedActivation: true }],
        [
            'trigger events without repeated activation',
            false,
            { values: [{ name: 'purchase' }], repeatedActivation: false },
        ],
        ['no events object', false, null],
    ])('%s -> %s', (_name, expected, events) => {
        const survey = { conditions: events ? { events } : null } as Pick<Survey, 'conditions'>
        expect(doesSurveyRepeatOnEveryEvent(survey)).toBe(expected)
    })
})

describe('getRecurringSurveyScheduleInfo', () => {
    it('computes the total run duration as count * frequency days', () => {
        const info = getRecurringSurveyScheduleInfo({
            schedule: SurveySchedule.Recurring,
            iteration_count: 2,
            iteration_frequency_days: 30,
            start_date: null,
            end_date: null,
        })
        expect(info).not.toBeNull()
        expect(info?.totalDurationDays).toBe(60)
        expect(info?.autoCloseDate).toBeNull()
    })

    it('computes the auto-close date from the start date in UTC', () => {
        const info = getRecurringSurveyScheduleInfo({
            schedule: SurveySchedule.Recurring,
            iteration_count: 2,
            iteration_frequency_days: 30,
            // Time-of-day near a UTC midnight boundary must not shift the calendar day the backend uses
            start_date: '2026-01-01T01:00:00Z',
            end_date: null,
        })
        // 2 iterations of 30 days -> closes 60 days after launch
        expect(info?.autoCloseDate?.format('YYYY-MM-DD')).toBe('2026-03-02')
    })

    it('returns null once the survey has already ended', () => {
        const info = getRecurringSurveyScheduleInfo({
            schedule: SurveySchedule.Recurring,
            iteration_count: 2,
            iteration_frequency_days: 30,
            start_date: '2026-01-01T00:00:00Z',
            end_date: '2026-01-15T00:00:00Z',
        })
        expect(info).toBeNull()
    })

    it('returns null for a non-recurring survey even with leftover iteration fields', () => {
        const info = getRecurringSurveyScheduleInfo({
            schedule: SurveySchedule.Once,
            iteration_count: 2,
            iteration_frequency_days: 30,
            start_date: '2026-01-01T00:00:00Z',
            end_date: null,
        })
        expect(info).toBeNull()
    })

    it('clamps the run duration to the backend iteration cap', () => {
        const info = getRecurringSurveyScheduleInfo({
            schedule: SurveySchedule.Recurring,
            // Above MAX_ITERATION_COUNT (500) — the backend only generates 500 windows
            iteration_count: 1000,
            iteration_frequency_days: 30,
            start_date: null,
            end_date: null,
        })
        expect(info?.totalDurationDays).toBe(500 * 30)
    })

    it.each([
        [
            'zero count',
            {
                schedule: SurveySchedule.Recurring,
                iteration_count: 0,
                iteration_frequency_days: 30,
                start_date: null,
                end_date: null,
            },
        ],
        [
            'zero frequency',
            {
                schedule: SurveySchedule.Recurring,
                iteration_count: 2,
                iteration_frequency_days: 0,
                start_date: null,
                end_date: null,
            },
        ],
        [
            'null count',
            {
                schedule: SurveySchedule.Recurring,
                iteration_count: null,
                iteration_frequency_days: 30,
                start_date: null,
                end_date: null,
            },
        ],
        [
            'null frequency',
            {
                schedule: SurveySchedule.Recurring,
                iteration_count: 2,
                iteration_frequency_days: null,
                start_date: null,
                end_date: null,
            },
        ],
    ])('returns null for %s', (_name, survey) => {
        expect(getRecurringSurveyScheduleInfo(survey)).toBeNull()
    })
})
````

### FILE: `third_party/posthog-nps/utils.ts.txt`

```yaml
block_id: "GO-CUSTOMER-SURVEY-API-NPS:file4:v1"
operation: CREATE
provenance: VERBATIM
source: "PostHog/posthog 6fafbb9081bd15e79448af5650e02a4f9ea435cc MIT core; exact source and adaptation boundaries in third_party/posthog-nps/source-lock.json"
license: "MIT"
sha256: "a43ed06c7eaffbe5bde3fe515177e512ebc4e6bce24dbef253ed100cbb1a08c1"
variables: []
secrets_allowed: false
```

````text
import DOMPurify from 'dompurify'
import { DeepPartialMap, ValidationErrorType } from 'kea-forms'
import posthog from 'posthog-js'

import { dayjs } from 'lib/dayjs'
import { dateStringToDayJs } from 'lib/utils/dateFilters'
import { getAppContext } from 'lib/utils/getAppContext'
import {
    MAX_ITERATION_COUNT,
    NEW_SURVEY,
    NewSurvey,
    SURVEY_CREATED_SOURCE,
    SURVEY_RATING_SCALE,
} from 'scenes/surveys/constants'
import { SurveyRatingResults } from 'scenes/surveys/surveyLogic'
import { urls } from 'scenes/urls'

import type { DataTableRow } from '~/queries/nodes/DataTable/dataTableLogic'
import {
    BasicSurveyQuestion,
    CyclotronJobInvocationGlobals,
    CyclotronJobFiltersType,
    EventPropertyFilter,
    EventType,
    FeatureFlagFilters,
    LinkSurveyQuestion,
    MultipleSurveyQuestion,
    PropertyFilterType,
    PropertyOperator,
    QuestionProcessedResponses,
    RatingSurveyQuestion,
    Survey,
    SurveyAppearance,
    SurveyDisplayConditions,
    SurveyEventName,
    SurveyEventProperties,
    SurveyQuestion,
    SurveyQuestionType,
    SurveyRates,
    SurveySchedule,
    SurveyStats,
    SurveyType,
} from '~/types'

const sanitizeConfig = { ADD_ATTR: ['target'] }

export function sanitizeHTML(html: string): string {
    return DOMPurify.sanitize(html, sanitizeConfig)
}

export function sanitizeColor(color: string | undefined): string | undefined {
    if (!color) {
        return undefined
    }

    // test if the color is valid by adding a # to the beginning of the string
    if (CSS.supports('color', `#${color}`)) {
        return `#${color}`
    }

    return color
}

export function validateCSSProperty(property: string, value: string | undefined): string | undefined {
    if (!value) {
        return undefined
    }
    const isValidCSSProperty = CSS.supports(property, value)
    return !isValidCSSProperty ? `${value} is not a valid property for ${property}.` : undefined
}

export function validateSurveyAppearance(
    appearance: SurveyAppearance,
    hasRatingQuestions: boolean,
    surveyType: SurveyType
): DeepPartialMap<SurveyAppearance, ValidationErrorType> {
    // API surveys are rendered by the customer, so PostHog's appearance CSS is not applied.
    // The Customization section is also hidden in the editor for API surveys (SurveyEdit.tsx),
    // so flagging appearance errors would route submitSurveyFailure to a non-existent section
    // and silently block saves.
    if (surveyType === SurveyType.API) {
        return {}
    }
    return {
        backgroundColor: validateCSSProperty('background-color', appearance.backgroundColor),
        borderColor: validateCSSProperty('border-color', appearance.borderColor),
        textColor: validateCSSProperty('color', appearance.textColor),
        inputBackground: validateCSSProperty('background-color', appearance.inputBackground),
        inputTextColor: validateCSSProperty('color', appearance.inputTextColor),
        // Only validate rating button colors if there's a rating question
        ...(hasRatingQuestions && {
            ratingButtonActiveColor: validateCSSProperty('background-color', appearance.ratingButtonActiveColor),
            ratingButtonColor: validateCSSProperty('background-color', appearance.ratingButtonColor),
        }),
        submitButtonColor: validateCSSProperty('background-color', appearance.submitButtonColor),
        submitButtonTextColor: validateCSSProperty('color', appearance.submitButtonTextColor),
        maxWidth: validateCSSProperty('width', appearance.maxWidth),
        boxPadding: validateCSSProperty('padding', appearance.boxPadding),
        boxShadow: validateCSSProperty('box-shadow', appearance.boxShadow),
        borderRadius: validateCSSProperty('border-radius', appearance.borderRadius),
        zIndex: validateCSSProperty('z-index', appearance.zIndex),
        widgetSelector:
            surveyType === SurveyType.Widget && appearance?.widgetType === 'selector' && !appearance.widgetSelector
                ? 'Please enter a CSS selector.'
                : undefined,
    }
}

export function getSurveyResponseKey(questionIndex: number): string {
    return questionIndex === 0
        ? SurveyEventProperties.SURVEY_RESPONSE
        : `${SurveyEventProperties.SURVEY_RESPONSE}_${questionIndex}`
}

export function getSurveyIdBasedResponseKey(questionId: string): string {
    return `${SurveyEventProperties.SURVEY_RESPONSE}_${questionId}`
}

type SurveyExampleContext = Pick<Survey, 'id' | 'name' | 'questions'> | null | undefined

function getExampleSurveyResponseValue(question: SurveyQuestion, index: number): string | string[] | undefined {
    switch (question.type) {
        case SurveyQuestionType.Open:
            return question.question || `Example answer ${index + 1}`
        case SurveyQuestionType.Rating:
            return String(question.scale >= 10 ? 9 : Math.min(question.scale, 4))
        case SurveyQuestionType.SingleChoice:
            return question.choices[0] || `Option ${index + 1}`
        case SurveyQuestionType.MultipleChoice:
            return question.choices.slice(0, Math.min(question.choices.length, 2))
        case SurveyQuestionType.Link:
            return undefined
    }
}

export function buildSurveyExampleInvocationGlobals({
    survey,
    projectId,
    projectName,
    projectUrl,
    source,
    timestamp = new Date().toISOString(),
    eventUuid = '00000000-0000-0000-0000-000000000000',
    distinctId = 'example-distinct-id',
    personId = 'person-id',
    personName = 'Jane Doe',
    personEmail = 'jane@example.com',
}: {
    survey: SurveyExampleContext
    projectId: number
    projectName: string
    projectUrl: string
    source?: CyclotronJobInvocationGlobals['source']
    timestamp?: string
    eventUuid?: string
    distinctId?: string
    personId?: string
    personName?: string
    personEmail?: string
}): CyclotronJobInvocationGlobals {
    const responseProperties = Object.fromEntries(
        (survey?.questions ?? [])
            .filter((question) => question.id && question.type !== SurveyQuestionType.Link)
            .map((question, index) => [
                getSurveyIdBasedResponseKey(question.id!),
                getExampleSurveyResponseValue(question, index),
            ])
            .filter(([, value]) => value !== undefined)
    )

    return {
        project: {
            id: projectId,
            name: projectName,
            url: projectUrl,
        },
        event: {
            event: SurveyEventName.SENT,
            uuid: eventUuid,
            distinct_id: distinctId,
            timestamp,
            elements_chain: '',
            properties: {
                [SurveyEventProperties.SURVEY_ID]: survey?.id && survey.id !== NEW_SURVEY.id ? survey.id : 'survey-id',
                $survey_name: survey?.name || 'Survey',
                [SurveyEventProperties.SURVEY_COMPLETED]: true,
                [SurveyEventProperties.SURVEY_SUBMISSION_ID]: 'survey-submission-id',
                ...responseProperties,
            },
            url: `${projectUrl}/events/${encodeURIComponent(eventUuid)}/${encodeURIComponent(timestamp)}`,
        },
        person: {
            id: personId,
            name: personName,
            url: `${projectUrl}/person/${encodeURIComponent(distinctId)}`,
            properties: {
                email: personEmail,
            },
        },
        groups: {},
        ...(source ? { source } : {}),
    }
}

// Helper function to generate the response field keys with proper typing
export const getResponseFieldWithId = (
    questionIndex: number,
    questionId?: string
): { indexBasedKey: string; idBasedKey: string | undefined } => {
    return {
        indexBasedKey: getSurveyResponseKey(questionIndex),
        idBasedKey: questionId ? getSurveyIdBasedResponseKey(questionId) : undefined,
    }
}

export function getSurveyResponseValue(
    eventProperties: Record<string, any>,
    questionIndex: number,
    questionId?: string
): any {
    const { indexBasedKey, idBasedKey } = getResponseFieldWithId(questionIndex, questionId)
    return (idBasedKey && eventProperties[idBasedKey]) ?? eventProperties[indexBasedKey]
}

export function sanitizeSurveyDisplayConditions(
    displayConditions?: SurveyDisplayConditions | null,
    surveyType?: SurveyType
): SurveyDisplayConditions | null {
    if (!displayConditions) {
        return null
    }

    if (surveyType === SurveyType.ExternalSurvey) {
        return {
            actions: {
                values: [],
            },
            events: {
                values: [],
            },
            deviceTypes: undefined,
            deviceTypesMatchType: undefined,
            linkedFlagVariant: undefined,
            seenSurveyWaitPeriodInDays: undefined,
            url: undefined,
            urlMatchType: undefined,
        }
    }

    const trimmedUrl = displayConditions.url?.trim()
    const trimmedSelector = displayConditions.selector?.trim()
    const trimmedLinkedFlagVariant = displayConditions.linkedFlagVariant?.trim()

    const sanitized: SurveyDisplayConditions = {
        ...displayConditions,
        ...(trimmedUrl && { url: trimmedUrl }),
        ...(trimmedSelector && { selector: trimmedSelector }),
        ...(trimmedLinkedFlagVariant && { linkedFlagVariant: trimmedLinkedFlagVariant }),
    }

    // Remove the original keys if they were empty after trimming
    if (!trimmedUrl) {
        delete sanitized.url
    }
    if (!trimmedSelector) {
        delete sanitized.selector
    }
    if (!trimmedLinkedFlagVariant) {
        delete sanitized.linkedFlagVariant
    }

    return sanitized
}

export function sanitizeSurveyAppearance(
    appearance?: SurveyAppearance | null,
    isPartialResponsesEnabled = false,
    surveyType?: SurveyType
): SurveyAppearance | null {
    if (!appearance) {
        return null
    }

    return {
        ...appearance,
        shuffleQuestions: isPartialResponsesEnabled ? false : appearance.shuffleQuestions,
        backgroundColor: sanitizeColor(appearance.backgroundColor),
        borderColor: sanitizeColor(appearance.borderColor),
        ratingButtonActiveColor: sanitizeColor(appearance.ratingButtonActiveColor),
        ratingButtonColor: sanitizeColor(appearance.ratingButtonColor),
        submitButtonColor: sanitizeColor(appearance.submitButtonColor),
        submitButtonTextColor: sanitizeColor(appearance.submitButtonTextColor),
        thankYouMessageHeader: sanitizeHTML(appearance.thankYouMessageHeader ?? ''),
        thankYouMessageDescription: sanitizeHTML(appearance.thankYouMessageDescription ?? ''),
        surveyPopupDelaySeconds:
            surveyType === SurveyType.ExternalSurvey ? undefined : appearance.surveyPopupDelaySeconds,
    }
}

export type NPSBreakdown = {
    total: number
    promoters: number
    passives: number
    detractors: number
    score: string
}

// NPS calculation constants
const NPS_SCALE_SIZE = 11 // 0-10 scale
const NPS_PROMOTER_MIN = 9 // 9-10 are promoters
const NPS_PASSIVE_MIN = 7 // 7-8 are passives. 0-6 are detractors but we don't need a variable for that.

interface NPSRawData {
    values: number[]
    total: number
}

/**
 * Extracts raw NPS data from processed survey data
 */
function extractNPSRawData(processedData: QuestionProcessedResponses): NPSRawData | null {
    if (
        !processedData?.data ||
        processedData.type !== SurveyQuestionType.Rating ||
        !Array.isArray(processedData.data) ||
        processedData.data.length !== NPS_SCALE_SIZE
    ) {
        return null
    }

    return {
        values: processedData.data.map((item) => item.value),
        total: processedData.totalResponses,
    }
}

/**
 * Extracts raw NPS data from legacy survey rating results
 */
function extractNPSRawDataFromLegacy(surveyRatingResults: SurveyRatingResults[number]): NPSRawData | null {
    if (!surveyRatingResults?.data || surveyRatingResults.data.length !== NPS_SCALE_SIZE) {
        return null
    }

    return {
        values: surveyRatingResults.data,
        total: surveyRatingResults.total,
    }
}

/**
 * Core NPS calculation logic - works with raw data arrays
 */
function calculateNPSFromRawData(rawData: NPSRawData): NPSBreakdown {
    if (rawData.total === 0) {
        return { total: 0, promoters: 0, passives: 0, detractors: 0, score: '0.0' }
    }

    const promoters = rawData.values.slice(NPS_PROMOTER_MIN, NPS_SCALE_SIZE).reduce((acc, curr) => acc + curr, 0)
    const passives = rawData.values.slice(NPS_PASSIVE_MIN, NPS_PROMOTER_MIN).reduce((acc, curr) => acc + curr, 0)
    const detractors = rawData.values.slice(0, NPS_PASSIVE_MIN).reduce((acc, curr) => acc + curr, 0)

    const score = ((promoters - detractors) / rawData.total) * 100

    return {
        total: rawData.total,
        promoters,
        passives,
        detractors,
        score: score.toFixed(1),
    }
}

export function calculateNpsBreakdownFromProcessedData(processedData: QuestionProcessedResponses): NPSBreakdown | null {
    const rawData = extractNPSRawData(processedData)
    return rawData ? calculateNPSFromRawData(rawData) : null
}

export function calculateNpsBreakdown(surveyRatingResults: SurveyRatingResults[number]): NPSBreakdown | null {
    const rawData = extractNPSRawDataFromLegacy(surveyRatingResults)
    return rawData ? calculateNPSFromRawData(rawData) : null
}

// Helper to escape special characters in SQL strings
function escapeSqlString(value: string): string {
    return value.replace(/['\\]/g, '\\$&')
}

export function getSurveyResponse(question: SurveyQuestion, index: number): string {
    // Delegate to the backend HogQL helper so survey response typing stays
    // consistent with PropertyDefinition metadata and materialized column rules.
    if (question.type === SurveyQuestionType.MultipleChoice) {
        return question.id
            ? `getSurveyResponse(${index}, '${question.id}', true)`
            : `getSurveyResponse(${index}, '', true)`
    }

    return question.id ? `getSurveyResponse(${index}, '${question.id}')` : `getSurveyResponse(${index})`
}

/**
 * Creates a HogQL expression for survey answer filters that handles both index-based and ID-based property keys
 * using OR logic between the alternative formats for each question.
 *
 * @param filters - The answer filters to convert to HogQL expressions
 * @param survey - The survey object (needed to access question IDs)
 * @param resolveResponseExpr - How to address a question's answer. Defaults to the event-level
 * `getSurveyResponse(...)` accessor. Callers querying merged submissions pass a resolver
 * returning the merged column alias instead, since `getSurveyResponse` is not in scope there.
 * @returns A HogQL expression string that can be used in queries. If there are no filters, it returns an empty string.
 *
 * TODO: Consider leveraging the backend query builder instead of duplicating this logic in the frontend.
 * ClickHouse has powerful functions like match(), multiIf(), etc. that could be used more effectively.
 */
export function createAnswerFilterHogQLExpression(
    filters: EventPropertyFilter[],
    survey: Survey,
    resolveResponseExpr: (question: SurveyQuestion, questionIndex: number) => string = getSurveyResponse
): string {
    if (!filters || !filters.length) {
        return ''
    }

    // Build the filter expression as a string
    let filterExpression = ''
    let hasValidFilter = false

    // Process each filter
    for (const filter of filters) {
        // Skip filters with empty or undefined values
        if (filter.value === undefined || filter.value === null || filter.value === '') {
            continue
        }

        // Skip empty arrays
        if (Array.isArray(filter.value) && filter.value.length === 0) {
            continue
        }

        // Skip ILIKE filters with empty search patterns
        if (
            filter.operator === 'icontains' &&
            (filter.value === '%' ||
                filter.value === '%%' ||
                (typeof filter.value === 'string' && filter.value.trim() === ''))
        ) {
            continue
        }

        // split the string '$survey_response_' and take the last part, as that's the question id
        const questionId = filter.key.split(`${SurveyEventProperties.SURVEY_RESPONSE}_`).at(-1)
        const question = survey.questions.find((question) => question.id === questionId)
        if (!questionId || !question) {
            continue
        }

        const questionIndex = survey.questions.findIndex((question) => question.id === questionId)

        // Create the condition for this filter
        let condition = ''
        const escapedValue = escapeSqlString(String(filter.value))

        // Handle different operators
        switch (filter.operator) {
            case 'exact':
            case 'is_not':
                if (Array.isArray(filter.value)) {
                    const valueList = filter.value.map((v) => `'${escapeSqlString(String(v))}'`).join(', ')
                    condition = `(${resolveResponseExpr(question, questionIndex)} ${
                        filter.operator === 'is_not' ? 'NOT IN' : 'IN'
                    } (${valueList}))`
                } else {
                    condition = `(${resolveResponseExpr(question, questionIndex)} ${
                        filter.operator === 'is_not' ? '!=' : '='
                    } '${escapedValue}')`
                }
                break
            case 'icontains':
                if (question.type !== SurveyQuestionType.MultipleChoice) {
                    condition = `(${resolveResponseExpr(question, questionIndex)} ILIKE '%${escapedValue}%')`
                } else {
                    condition = `(arrayExists(x -> x ilike '%${escapedValue}%', ${resolveResponseExpr(question, questionIndex)}))`
                }
                break
            case 'not_icontains':
                if (question.type !== SurveyQuestionType.MultipleChoice) {
                    condition = `(NOT ${resolveResponseExpr(question, questionIndex)} ILIKE '%${escapedValue}%')`
                } else {
                    condition = `(NOT arrayExists(x -> x ilike '%${escapedValue}%', ${resolveResponseExpr(question, questionIndex)}))`
                }
                break
            case 'regex':
                if (question.type !== SurveyQuestionType.MultipleChoice) {
                    condition = `(match(${resolveResponseExpr(question, questionIndex)}, '${escapedValue}'))`
                } else {
                    condition = `(arrayExists(x -> match(x, '${escapedValue}'), ${resolveResponseExpr(question, questionIndex)}))`
                }
                break
            case 'not_regex':
                if (question.type !== SurveyQuestionType.MultipleChoice) {
                    condition = `(NOT match(${resolveResponseExpr(question, questionIndex)}, '${escapedValue}'))`
                } else {
                    condition = `(NOT arrayExists(x -> match(x, '${escapedValue}'), ${resolveResponseExpr(question, questionIndex)}))`
                }
                break
            // Add more operators as needed
            default:
                continue // Skip unsupported operators
        }

        // Add this condition to the overall expression
        if (condition) {
            if (hasValidFilter) {
                filterExpression += ' AND '
            }
            filterExpression += condition
            hasValidFilter = true
        }
    }

    return hasValidFilter ? `AND ${filterExpression}` : ''
}

export function isSurveyRunning(survey: Pick<Survey, 'start_date' | 'end_date'>): boolean {
    return !!(survey.start_date && !survey.end_date)
}

// Auto-submit only makes sense for questions where a single selection is a complete
// answer: any rating, or a single-choice question without a free-text "open" option.
export function canQuestionSkipSubmitButton(
    question: SurveyQuestion
): question is RatingSurveyQuestion | MultipleSurveyQuestion {
    return (
        question.type === SurveyQuestionType.Rating ||
        (question.type === SurveyQuestionType.SingleChoice && !question.hasOpenChoice)
    )
}

// Some fields can only be edited in the full editor — opening such a survey
// in the wizard would hide those values from the user, so we route them to
// the full editor regardless of their general editor preference. Keep this
// list in sync with what the wizard's steps actually expose.
export function canUseSurveyWizard(survey: Survey | NewSurvey): boolean {
    if (survey.type !== SurveyType.Popover) {
        return false
    }
    // SurveySchedule.Always — the wizard offers Once + recurring frequencies, but not "every time
    // the display conditions are met". Keep Always surveys in the legacy editor where the option
    // is actually visible, so the wizard never silently misrepresents the cadence.
    if (survey.schedule === SurveySchedule.Always) {
        return false
    }
    // Adaptive sampling — WhenStep exposes a simple responses_limit but not
    // the adaptive sampling controls
    if (survey.response_sampling_limit || survey.response_sampling_start_date) {
        return false
    }
    // Property-based targeting filters — WhereStep handles linked_flag
    // (release conditions) but not targeting_flag_filters
    if (survey.targeting_flag_filters && Object.keys(survey.targeting_flag_filters).length > 0) {
        return false
    }
    return true
}

export function doesSurveyRepeatOnEveryEvent(survey: Pick<Survey, 'conditions'>): boolean {
    return !!(survey.conditions?.events?.repeatedActivation && (survey.conditions?.events?.values?.length ?? 0) > 0)
}

export interface RecurringSurveyScheduleInfo {
    /** Total number of days the survey runs from its launch date before auto-closing. */
    totalDurationDays: number
    /** The date the survey will automatically close, or null if it hasn't been launched yet. */
    autoCloseDate: dayjs.Dayjs | null
}

/**
 * A recurring survey ("Repeat on a schedule") auto-closes once its final iteration window has passed.
 * The last iteration starts on `start_date + (count - 1) * frequency` days and lasts `frequency` more days,
 * so the survey runs for `count * frequency` days total and closes at the end of that span.
 * Mirrors the backend logic in posthog/tasks/update_survey_iteration.py, which computes iteration windows
 * on the UTC calendar day — so we do the arithmetic in UTC too.
 *
 * Returns null once the survey has already ended: it then shows its real end date, so a projected one would
 * only contradict it.
 */
export function getRecurringSurveyScheduleInfo(
    survey: Pick<Survey, 'schedule' | 'iteration_count' | 'iteration_frequency_days' | 'start_date' | 'end_date'>
): RecurringSurveyScheduleInfo | null {
    const count = survey.iteration_count
    const frequency = survey.iteration_frequency_days
    if (
        survey.schedule !== SurveySchedule.Recurring ||
        survey.end_date ||
        !count ||
        !frequency ||
        count < 1 ||
        frequency < 1
    ) {
        return null
    }
    // The backend caps the generated iteration windows at MAX_ITERATION_COUNT, so anything above that never
    // extends the schedule — mirror the cap here to match the real close date.
    const effectiveCount = Math.min(count, MAX_ITERATION_COUNT)
    const totalDurationDays = effectiveCount * frequency
    const autoCloseDate = survey.start_date ? dayjs.utc(survey.start_date).add(totalDurationDays, 'day') : null
    return { totalDurationDays, autoCloseDate }
}

export function doesSurveyHaveDisplayConditions(survey: Survey | NewSurvey): boolean {
    const conditions = sanitizeSurveyDisplayConditions(survey.conditions)
    if (!conditions) {
        return false
    }

    // check string fields
    if (conditions.url) {
        return true
    }

    if (conditions.selector) {
        return true
    }

    if (conditions.linkedFlagVariant) {
        return true
    }

    // check numeric fields
    if (conditions.seenSurveyWaitPeriodInDays !== undefined && conditions.seenSurveyWaitPeriodInDays !== null) {
        return true
    }

    // check array fields
    if (conditions.deviceTypes && conditions.deviceTypes.length > 0) {
        return true
    }

    // check enum fields
    if (conditions.urlMatchType !== undefined && conditions.urlMatchType !== null) {
        return true
    }

    if (conditions.deviceTypesMatchType !== undefined && conditions.deviceTypesMatchType !== null) {
        return true
    }

    // check complex object fields
    if (conditions.actions && conditions.actions.values && conditions.actions.values.length > 0) {
        return true
    }

    if (conditions.events && conditions.events.values && conditions.events.values.length > 0) {
        return true
    }

    if (conditions.events?.repeatedActivation !== undefined && conditions.events.repeatedActivation !== null) {
        return true
    }

    return false
}

export function buildSurveyOptionalBooleanPropertyFilter(
    propertyName: SurveyEventProperties,
    excludedValue: 'true' | 'false'
): string {
    return `coalesce(JSONExtractString(properties, '${propertyName}'), '') != '${excludedValue}'`
}

export function buildSurveyResponseEventFilter(): string {
    return `(event = '${SurveyEventName.SENT}' OR (
        event IN ('${SurveyEventName.DISMISSED}', '${SurveyEventName.ABANDONED}')
        AND coalesce(JSONExtractString(properties, '${SurveyEventProperties.SURVEY_PARTIALLY_COMPLETED}'), '') = 'true'
    ))`
}

export interface SurveyQueryFilters {
    timestampFilter: string
    answerFilters: EventPropertyFilter[]
    archivedResponsesFilter: string
}

/**
 * HogQL expression collapsing a submission's response events into one group. An event with
 * no `$survey_submission_id` is keyed by its own uuid, so it stays a distinct response the way it
 * did before submission IDs existed.
 *
 * Must stay identical to `SUBMISSION_GROUPING_KEY` in
 * `products/surveys/backend/responses/fetch_rows.py`, otherwise the Results tab and the responses
 * API disagree about what counts as one submission.
 */
const SUBMISSION_GROUPING_KEY = `if(
    coalesce(properties.\`${SurveyEventProperties.SURVEY_SUBMISSION_ID}\`, '') = '',
    toString(uuid),
    properties.\`${SurveyEventProperties.SURVEY_SUBMISSION_ID}\`
)`

/** Alias holding a question's merged answer in the submission-merge subquery. */
function mergedAnswerAlias(questionIndex: number): string {
    return `q${questionIndex}_answer`
}

/** Alias holding a question's raw per-event answer in the submission-merge subquery. */
function rawAnswerAlias(questionIndex: number): string {
    return `q${questionIndex}_raw`
}

/**
 * True when the event actually carries an answer to this question, so the merge can pick the event
 * that answered it rather than whichever event in the submission happens to be latest.
 *
 * Multiple-choice answers are arrays, which are never null, so they need a length check instead.
 */
function buildAnswerPresenceExpr(rawAlias: string, question: SurveyQuestion): string {
    return question.type === SurveyQuestionType.MultipleChoice ? `length(${rawAlias}) > 0` : `isNotNull(${rawAlias})`
}

/** True when the merged answer holds no content, covering both the null and empty-string cases. */
export function buildAnswerIsEmptyExpr(mergedAlias: string, question: SurveyQuestion): string {
    return question.type === SurveyQuestionType.MultipleChoice
        ? `length(${mergedAlias}) = 0`
        : `length(trim(coalesce(${mergedAlias}, ''))) = 0`
}

interface QuestionWithIndex {
    question: SurveyQuestion
    index: number
}

/**
 * Builds a subquery emitting one row per submission, with every question's answer merged across
 * that submission's events.
 *
 * A submission can span several `survey sent` events that don't each repeat the answers given
 * earlier. The AI feedback flow produces exactly that shape: the rating arrives on one event and
 * the free-text follow-up on another, joined by `$survey_submission_id`. Electing a single event
 * per submission therefore drops every answer that only ever lived on a non-elected event, which
 * is why this merges per question with `argMaxIf` instead, keeping the latest answer to each.
 *
 * This mirrors the responses API in `products/surveys/backend/responses/fetch_rows.py`, including
 * its use of `isNotNull` rather than a stricter emptiness test, so both surfaces resolve a
 * re-answered question the same way.
 */
function buildMergedSubmissionsSubquery(
    survey: Survey,
    filters: SurveyQueryFilters,
    questions: QuestionWithIndex[],
    { includeRespondentMetadata = false }: { includeRespondentMetadata?: boolean } = {}
): string {
    const completedEventExpr = `event = '${SurveyEventName.SENT}' AND ${buildSurveyOptionalBooleanPropertyFilter(SurveyEventProperties.SURVEY_COMPLETED, 'false')}`

    const innerColumns = [
        'uuid AS event_uuid',
        'timestamp',
        'person_id',
        ...(includeRespondentMetadata
            ? [
                  'distinct_id',
                  'properties.`$session_id` AS session_id',
                  'properties AS event_properties',
                  'person.properties AS person_properties',
              ]
            : []),
        `${completedEventExpr} AS is_completed_event`,
        'event',
        ...questions.map(({ question, index }) => `${getSurveyResponse(question, index)} AS ${rawAnswerAlias(index)}`),
        `${SUBMISSION_GROUPING_KEY} AS submission_key`,
    ]

    const outerColumns = [
        'argMax(event_uuid, tuple(timestamp, event_uuid)) AS uuid',
        'argMax(person_id, tuple(timestamp, event_uuid)) AS person_id',
        `if(countIf(is_completed_event) > 0, 'completed', if(argMax(event, tuple(timestamp, event_uuid)) = '${SurveyEventName.DISMISSED}', 'dismissed', 'abandoned')) AS outcome`,
        // Aliased away from `timestamp` because every other aggregate here orders by that column,
        // and an alias of the same name would resolve to this aggregate instead, nesting them.
        'max(timestamp) AS submitted_at',
        ...(includeRespondentMetadata
            ? [
                  'argMax(distinct_id, tuple(timestamp, event_uuid)) AS distinct_id',
                  'argMax(session_id, tuple(timestamp, event_uuid)) AS session_id',
                  'argMax(event_properties, tuple(timestamp, event_uuid)) AS event_properties',
                  'argMax(person_properties, tuple(timestamp, event_uuid)) AS person_properties',
                  'argMax(event, tuple(timestamp, event_uuid)) AS latest_event',
              ]
            : []),
        ...questions.map(({ question, index }) => {
            const raw = rawAnswerAlias(index)
            return `argMaxIf(${raw}, tuple(timestamp, event_uuid), ${buildAnswerPresenceExpr(raw, question)}) AS ${mergedAnswerAlias(index)}`
        }),
    ]

    // Answer and archive filters read the merged answer, so they belong in HAVING. The archive
    // filter names `uuid`, which resolves to the representative uuid aliased above — the same one
    // the responses table archives.
    const havingConditions: string[] = []
    const mergedAnswerFilter = createAnswerFilterHogQLExpression(filters.answerFilters, survey, (_, index) =>
        mergedAnswerAlias(index)
    )
    if (mergedAnswerFilter !== '') {
        havingConditions.push(stripLeadingAnd(mergedAnswerFilter))
    }
    if (filters.archivedResponsesFilter !== '') {
        havingConditions.push(stripLeadingAnd(filters.archivedResponsesFilter))
    }

    return `SELECT ${outerColumns.join(',\n            ')}
        FROM (
            SELECT ${innerColumns.join(',\n                ')}
            FROM events
            WHERE ${buildSurveyResponseEventFilter()}
                AND properties.\`${SurveyEventProperties.SURVEY_ID}\` = '${survey.id}'
                ${filters.timestampFilter}
                AND {filters}
        )
        GROUP BY submission_key${havingConditions.length > 0 ? `\n        HAVING ${havingConditions.join(' AND ')}` : ''}`
}

export function getSurveyResponseStatus(
    eventName: string | undefined,
    properties: Record<string, unknown>
): string | null {
    const completed = properties[SurveyEventProperties.SURVEY_COMPLETED]
    if (completed === true || completed === 'true') {
        return null
    }
    const partial = properties[SurveyEventProperties.SURVEY_PARTIALLY_COMPLETED]
    if (completed !== false && completed !== 'false' && partial !== true && partial !== 'true') {
        return null
    }
    if (eventName === SurveyEventName.DISMISSED) {
        return 'Dismissed'
    }
    return 'Abandoned'
}

export function isSurveyResponseEvent(eventName: string, properties: Record<string, unknown>): boolean {
    return (
        !!properties[SurveyEventProperties.SURVEY_ID] &&
        (eventName === SurveyEventName.SENT ||
            (([SurveyEventName.DISMISSED, SurveyEventName.ABANDONED] as string[]).includes(eventName) &&
                [true, 'true'].includes(
                    properties[SurveyEventProperties.SURVEY_PARTIALLY_COMPLETED] as boolean | string
                )))
    )
}

export function transformSurveyResponseRows(rows: DataTableRow[], survey: Pick<Survey, 'questions'>): DataTableRow[] {
    return rows.map((row) => {
        if (!Array.isArray(row.result) || !Array.isArray(row.result[0])) {
            return row
        }
        const [
            uuid,
            distinctId,
            timestamp,
            personId,
            personProperties,
            eventProperties,
            outcome,
            answers,
            latestEvent,
        ] = row.result[0]
        const properties = { ...JSON.parse(eventProperties || '{}') }
        survey.questions.forEach((question, index) => {
            const answer = answers[index]
            if (answer !== null && answer !== undefined) {
                properties[getSurveyResponseKey(index)] = answer
                if (question.id) {
                    properties[`$survey_response_${question.id}`] = answer
                }
            }
        })
        properties[SurveyEventProperties.SURVEY_COMPLETED] = outcome === 'completed'
        properties[SurveyEventProperties.SURVEY_PARTIALLY_COMPLETED] = outcome !== 'completed'
        const event: EventType = {
            id: uuid,
            uuid,
            distinct_id: distinctId,
            timestamp,
            event: outcome === 'completed' ? SurveyEventName.SENT : latestEvent,
            properties,
            person_id: personId,
            person: {
                is_identified: false,
                distinct_ids: [distinctId],
                properties: JSON.parse(personProperties || '{}'),
            },
            elements: [],
        }
        return { ...row, result: [event, ...row.result.slice(1)] }
    })
}

export function buildSurveyResponsesQuery(survey: Survey, filters: SurveyQueryFilters): string {
    const questions = getAnswerableQuestions(survey)
    const merged = buildMergedSubmissionsSubquery(survey, filters, questions, { includeRespondentMetadata: true })
    const answers = survey.questions.map((question, index) =>
        question.type !== SurveyQuestionType.Link ? mergedAnswerAlias(index) : 'NULL'
    )
    const columns = [
        `tuple(uuid, distinct_id, submitted_at, person_id, person_properties, event_properties, outcome, tuple(${answers.length ? answers.join(', ') : 'NULL'}), latest_event) AS response`,
        ...survey.questions.map(
            (question, index) =>
                `${question.type === SurveyQuestionType.MultipleChoice ? `arrayStringConcat(${answers[index]}, ', ')` : answers[index]} AS answer_${index}`
        ),
        'outcome AS status',
        'submitted_at AS timestamp',
        'distinct_id AS respondent',
        'uuid AS actions',
    ]
    return `SELECT ${columns.join(',\n')} FROM (${merged}) ORDER BY submitted_at DESC`
}

export function buildSurveyResponseSQLQuery(
    survey: Survey,
    filters: SurveyQueryFilters,
    questionIndex?: number
): string {
    const merged = buildMergedSubmissionsSubquery(survey, filters, getAnswerableQuestions(survey), {
        includeRespondentMetadata: true,
    }).replaceAll('{filters}', '1 = 1')
    const columns = survey.questions.flatMap((question, index) => {
        if (question.type === SurveyQuestionType.Link || (questionIndex !== undefined && index !== questionIndex)) {
            return []
        }
        const title = (question.question || `Question ${index + 1}`).replace(/\s*[\r\n]+\s*/g, ' ').replace(/"/g, '""')
        return [`${mergedAnswerAlias(index)} AS "${title}"`]
    })
    return `SELECT distinct_id, ${columns.length ? columns.join(', ') + ', ' : ''}outcome, submitted_at
        FROM (${merged}) ORDER BY submitted_at DESC LIMIT 100`
}

export function buildSurveyResponseStatsQuery(survey: Survey, filters: SurveyQueryFilters): string {
    const merged = buildMergedSubmissionsSubquery(survey, filters, getAnswerableQuestions(survey))
    return `SELECT '${SurveyEventName.SENT}' AS event_name, count() AS total_count,
        count(DISTINCT person_id) AS unique_persons,
        if(count() > 0, min(submitted_at), null) AS first_seen,
        if(count() > 0, max(submitted_at), null) AS last_seen,
        tuple(countIf(outcome = 'completed'), countIf(outcome = 'dismissed'), countIf(outcome = 'abandoned')) AS outcome_counts
        FROM (${merged})`
}

export interface SurveyResponseOutcome {
    label: string
    count: number
    percentage: number
}

export function getSurveyResponseOutcomeBreakdown(counts: [number, number, number]): SurveyResponseOutcome[] {
    const total = counts.reduce((sum, count) => sum + count, 0)
    return ['Completed', 'Dismissed', 'Abandoned'].map((label, index) => ({
        label,
        count: counts[index],
        percentage: total > 0 ? counts[index] / total : 0,
    }))
}

export function buildSurveyRespondentQuery(survey: Survey, filters: SurveyQueryFilters): string {
    return `SELECT person_id FROM (${buildMergedSubmissionsSubquery(survey, filters, getAnswerableQuestions(survey))})`
}

function stripLeadingAnd(expression: string): string {
    return expression.replace(/^\s*AND\s+/, '')
}

/** Questions that can hold an answer. Link questions never produce a response. */
function getAnswerableQuestions(survey: Survey): QuestionWithIndex[] {
    return survey.questions
        .map((question, index) => ({ question, index }))
        .filter(({ question }) => question.type !== SurveyQuestionType.Link)
}

export interface OpenEndedColumnMap {
    [questionId: string]: {
        columnIndex: number
        questionIndex: number
        type: SurveyQuestionType.Open | SurveyQuestionType.SingleChoice | SurveyQuestionType.MultipleChoice
    }
}

export function buildAggregateQuery(survey: Survey, filters: SurveyQueryFilters): string | null {
    const questions = getAnswerableQuestions(survey)
    if (questions.length === 0) {
        return null
    }

    // Each entry emits the (question_id, label) pairs one submission contributes to one question.
    // They are concatenated and unrolled with a single arrayJoin so the merge below is read once,
    // rather than once per question: a ClickHouse CTE is inlined, so a UNION ALL branch per
    // question would re-run the whole merge per branch.
    const labelPairs: string[] = []
    const noPairs = '[]'

    for (const { question, index } of questions) {
        const answer = mergedAnswerAlias(index)
        const questionId = `'${question.id}'`
        // The merged answer is nullable, and a nullable label would not match the literal pairs
        // below when the arrays are concatenated.
        const answerLabel = `coalesce(toString(${answer}), '')`

        if (question.type === SurveyQuestionType.Rating || question.type === SurveyQuestionType.SingleChoice) {
            labelPairs.push(`if(isNotNull(${answer}), [(${questionId}, ${answerLabel})], ${noPairs})`)

            if (question.type === SurveyQuestionType.SingleChoice && question.optional) {
                labelPairs.push(
                    `if(${buildAnswerIsEmptyExpr(answer, question)}, [(${questionId}, '__no_response__')], ${noPairs})`
                )
            }
        } else if (question.type === SurveyQuestionType.MultipleChoice) {
            labelPairs.push(
                `arrayMap(choice -> (${questionId}, choice),
                    arrayFilter(choice -> choice != '',
                        arrayMap(choice -> trim(BOTH '"\\'' FROM choice), ${answer})))`
            )
            labelPairs.push(`if(length(${answer}) > 0, [(${questionId}, '__total__')], ${noPairs})`)

            if (question.optional) {
                labelPairs.push(`if(length(${answer}) = 0, [(${questionId}, '__no_response__')], ${noPairs})`)
            }
        } else if (question.type === SurveyQuestionType.Open) {
            labelPairs.push(`if(isNotNull(${answer}), [(${questionId}, '__total__')], ${noPairs})`)
        }
    }

    if (labelPairs.length === 0) {
        return null
    }

    const mergedSubmissions = buildMergedSubmissionsSubquery(survey, filters, questions)

    // arrayConcat needs two arguments or more. A survey that emits one pair expression, such as a
    // single rating question, goes straight to arrayJoin.
    const allLabelPairs =
        labelPairs.length === 1
            ? labelPairs[0]
            : `arrayConcat(\n                ${labelPairs.join(',\n                ')}\n            )`

    return `SELECT
            tupleElement(question_label, 1) AS question_id,
            tupleElement(question_label, 2) AS label,
            count() AS cnt
        FROM (
            SELECT arrayJoin(${allLabelPairs}) AS question_label
            FROM (
                ${mergedSubmissions}
            )
        )
        GROUP BY question_id, label
        LIMIT 50000`
}

export function buildOpenEndedQuery(
    survey: Survey,
    filters: SurveyQueryFilters,
    limit: number = 50000
): { query: string; columnMap: OpenEndedColumnMap } | null {
    const questions = getAnswerableQuestions(survey)
    const openColumns: string[] = []
    const columnMap: OpenEndedColumnMap = {}
    let columnIndex = 0

    for (const { question, index } of questions) {
        const isOpen = question.type === SurveyQuestionType.Open
        const hasOpenChoice =
            (question.type === SurveyQuestionType.SingleChoice ||
                question.type === SurveyQuestionType.MultipleChoice) &&
            (question as MultipleSurveyQuestion).hasOpenChoice

        if (isOpen || hasOpenChoice) {
            openColumns.push(`${mergedAnswerAlias(index)} AS q${index}_response`)
            columnMap[question.id!] = { columnIndex, questionIndex: index, type: question.type }
            columnIndex++
        }
    }

    if (openColumns.length === 0) {
        return null
    }

    // The merge needs every answerable question, not just the open ones, because answer filters in
    // HAVING can reference a question that has no open column of its own.
    const mergedSubmissions = buildMergedSubmissionsSubquery(survey, filters, questions, {
        includeRespondentMetadata: true,
    })

    // Column order stays open columns, then distinct_id, timestamp, session_id — processOpenEndedResults
    // reads the metadata positionally from the end.
    const query = `SELECT
            ${openColumns.join(',\n')},
            distinct_id,
            submitted_at,
            session_id
        FROM (
            ${mergedSubmissions}
        )
        ORDER BY submitted_at DESC
        LIMIT ${limit}`

    return { query, columnMap }
}

interface SanitizeSurveyOptions {
    keepEmptyConditions?: boolean
}

export function sanitizeSurvey(survey: Partial<Survey>, options?: SanitizeSurveyOptions): Partial<Survey> {
    const sanitizedQuestions =
        survey.questions?.map((question) => {
            const sanitized = {
                ...question,
                question: sanitizeHTML(question.question ?? ''),
                description: sanitizeHTML(question.description ?? ''),
            }
            if (
                (sanitized.type === SurveyQuestionType.SingleChoice ||
                    sanitized.type === SurveyQuestionType.MultipleChoice) &&
                sanitized.choices
            ) {
                sanitized.choices = sanitized.choices.map((choice) => choice.trim())
            }
            // Drop a stale auto-submit flag if the question is no longer eligible for it
            // (e.g. an open-ended choice was added, or the type was switched).
            if ('skipSubmitButton' in sanitized && !canQuestionSkipSubmitButton(sanitized)) {
                delete (sanitized as { skipSubmitButton?: boolean }).skipSubmitButton
            }
            return sanitized
        }) || []

    const sanitizedAppearance = sanitizeSurveyAppearance(
        survey.appearance,
        survey.enable_partial_responses ?? false,
        survey.type
    )

    // Remove widget-specific fields if survey type is not Widget
    if (survey.type !== SurveyType.Widget && sanitizedAppearance) {
        delete sanitizedAppearance.widgetType
        delete sanitizedAppearance.widgetLabel
        delete sanitizedAppearance.widgetColor
    }

    const conditions = sanitizeSurveyDisplayConditions(survey.conditions, survey.type)
    const sanitized: Partial<Survey> = {
        ...survey,
        conditions: conditions,
        questions: sanitizedQuestions,
        appearance: sanitizedAppearance,
    }

    if (survey.type === SurveyType.ExternalSurvey) {
        sanitized.remove_targeting_flag = true
        sanitized.linked_flag_id = null
        sanitized.targeting_flag_filters = undefined
    }

    if (options?.keepEmptyConditions !== true && (!conditions || Object.keys(conditions).length === 0)) {
        delete sanitized.conditions
    }
    if (!sanitizedAppearance || Object.keys(sanitizedAppearance).length === 0) {
        delete sanitized.appearance
    }

    return sanitized
}

export function calculateSurveyRates(stats: SurveyStats | null): SurveyRates {
    const defaultRates: SurveyRates = {
        response_rate: 0.0,
        dismissal_rate: 0.0,
        unique_users_response_rate: 0.0,
        unique_users_dismissal_rate: 0.0,
    }

    if (!stats) {
        return defaultRates
    }

    const shownCount = stats[SurveyEventName.SHOWN].total_count
    if (shownCount > 0) {
        const sentCount = stats[SurveyEventName.SENT].total_count
        const dismissedCount = stats[SurveyEventName.DISMISSED].total_count
        const uniqueUsersShownCount = stats[SurveyEventName.SHOWN].unique_persons
        const uniqueUsersSentCount = stats[SurveyEventName.SENT].unique_persons
        const uniqueUsersDismissedCount = stats[SurveyEventName.DISMISSED].unique_persons

        return {
            response_rate: parseFloat(((sentCount / shownCount) * 100).toFixed(2)),
            dismissal_rate: parseFloat(((dismissedCount / shownCount) * 100).toFixed(2)),
            unique_users_response_rate: parseFloat(((uniqueUsersSentCount / uniqueUsersShownCount) * 100).toFixed(2)),
            unique_users_dismissal_rate: parseFloat(
                ((uniqueUsersDismissedCount / uniqueUsersShownCount) * 100).toFixed(2)
            ),
        }
    }
    return defaultRates
}

export function captureMaxAISurveyCreationException(error?: string, source?: SURVEY_CREATED_SOURCE): void {
    posthog.captureException(error || 'Undefined error when creating MaxAI survey', {
        action: 'max-ai-survey-creation-failed',
        source: source,
    })
}

export const DATE_FORMAT = 'YYYY-MM-DDTHH:mm:ss'

function getTeamTimezone(): string {
    return getAppContext()?.current_team?.timezone || 'UTC'
}

export function getSurveyStartDateForQuery(
    survey: Pick<Survey, 'created_at'> & Partial<Pick<Survey, 'start_date'>>
): string {
    const tz = getTeamTimezone()
    return dayjs
        .tz(survey.start_date ?? survey.created_at, tz)
        .startOf('day')
        .format(DATE_FORMAT)
}

export function getSurveyEndDateForQuery(survey: Pick<Survey, 'end_date'>): string {
    const tz = getTeamTimezone()
    return survey.end_date
        ? dayjs.tz(survey.end_date, tz).endOf('day').format(DATE_FORMAT)
        : dayjs.tz(undefined, tz).endOf('day').format(DATE_FORMAT)
}

export interface SurveyDateRange {
    date_from: string | null
    date_to: string | null
}

export function getResolvedSurveyDateRange(
    survey: Pick<Survey, 'created_at' | 'end_date'> & Partial<Pick<Survey, 'start_date'>>,
    dateRange?: SurveyDateRange | null
): { fromDate: string; toDate: string } {
    let fromDate = getSurveyStartDateForQuery(survey)
    let toDate = getSurveyEndDateForQuery(survey)

    // date_from only is valid ("from custom date until now")
    // date_to only is ignored to avoid impossible ranges
    if (dateRange?.date_from) {
        const tz = getTeamTimezone()
        fromDate = dateStringToDayJs(dateRange.date_from, tz)?.startOf('day').format(DATE_FORMAT) ?? fromDate

        if (dateRange.date_to) {
            toDate = dateStringToDayJs(dateRange.date_to, tz)?.endOf('day').format(DATE_FORMAT) ?? toDate
        }
    }

    return { fromDate, toDate }
}

export function buildSurveyTimestampFilter(
    survey: Pick<Survey, 'created_at' | 'end_date'> & Partial<Pick<Survey, 'start_date'>>,
    dateRange?: SurveyDateRange | null
): string {
    const { fromDate, toDate } = getResolvedSurveyDateRange(survey, dateRange)

    return `AND timestamp >= '${fromDate}'
    AND timestamp <= '${toDate}'`
}

export function getExpressionCommentForQuestion(
    q: BasicSurveyQuestion | LinkSurveyQuestion | RatingSurveyQuestion | MultipleSurveyQuestion,
    questionIndex: number
): string {
    const question = q.question.trim()
    if (question.length > 0) {
        // This is appended after `--` in the generated HogQL, and HogQL `--` comments are
        // single-line. Collapse any newlines so multi-line question text can't leak past the
        // comment and break the query (e.g. a stray non-ASCII char -> "Unexpected character").
        return question.replace(/\s*[\r\n]+\s*/g, ' ')
    }
    return `Question ${questionIndex + 1}`
}

export function getSurveyForFeatureFlagVariant(variantKey: string, surveys?: Survey[]): Survey | undefined {
    return surveys?.find((survey) => survey.conditions?.linkedFlagVariant === variantKey)
}

export function duplicateExistingSurvey(survey: Survey | NewSurvey): Partial<Survey> {
    return {
        ...survey,
        questions: survey.questions.map((question) => ({
            ...question,
            id: undefined,
        })),
        id: NEW_SURVEY.id,
        name: `${survey.name} (duplicated at ${dayjs().format('YYYY-MM-DD HH:mm:ss')})`,
        archived: false,
        start_date: null,
        end_date: null,
        targeting_flag_filters: survey.targeting_flag?.filters ?? NEW_SURVEY.targeting_flag_filters,
        linked_flag_id: survey.linked_flag?.id ?? NEW_SURVEY.linked_flag_id,
    }
}

export const isThumbQuestion = (question: SurveyQuestion): boolean => {
    return (
        question.type === SurveyQuestionType.Rating &&
        question.display === 'emoji' &&
        question.scale === SURVEY_RATING_SCALE.THUMB_2_POINT
    )
}

/**
 * A 2-point rating question always represents a binary thumbs up / thumbs down regardless of `display`,
 * so we render the icon + label in response views to make the value readable at a glance.
 */
export const isScaleTwoRating = (question: SurveyQuestion): boolean => {
    return question.type === SurveyQuestionType.Rating && question.scale === SURVEY_RATING_SCALE.THUMB_2_POINT
}

/**
 * Splits text pasted into a choice input on newlines or tabs (spreadsheet rows).
 * Returns the merged choices array, or `null` if there's nothing to split (the caller
 * should let the paste fall through to the default input behavior).
 *
 * Always keeps the open-ended ("Other") entry as the last item when `hasOpenChoice`
 * is true — including when the paste happens into the open-ended slot itself.
 */
export function splitChoicesOnPaste(
    pasted: string,
    choices: string[],
    choiceIndex: number,
    hasOpenChoice: boolean
): string[] | null {
    const segments = pasted
        .split(/[\n\t]+/)
        .map((segment) => segment.trim())
        .filter((segment) => segment.length > 0)

    if (segments.length <= 1) {
        return null
    }

    const openTail = hasOpenChoice ? [choices[choices.length - 1]] : []
    const head = choices.slice(0, choiceIndex)
    const tailStart = choiceIndex + 1
    const tailEnd = hasOpenChoice ? choices.length - 1 : choices.length
    const tail = choices.slice(tailStart, tailEnd)
    return [...head, ...segments, ...tail, ...openTail]
}

export type SurveyConditionType =
    | 'url'
    | 'selector'
    | 'device'
    | 'events'
    | 'actions'
    | 'flag'
    | 'flag_variant'
    | 'targeting'
    | 'wait_period'

export interface SurveyConditionSummary {
    type: SurveyConditionType
    label: string
    value: string
    href?: string
}

export interface SurveyCollectionLimitSummary {
    label: 'Response limit' | 'Sampling limit'
    value: string
}

export function getSurveyTargetingFilters(survey: Survey | NewSurvey): FeatureFlagFilters | undefined {
    if (survey.targeting_flag_filters) {
        return survey.targeting_flag_filters
    }

    return survey.targeting_flag?.filters || undefined
}

export function getSurveyAudienceRuleCount(filters?: FeatureFlagFilters | null): number {
    return filters?.groups.reduce((count, group) => count + (group.properties?.length ?? 0), 0) ?? 0
}

export function getSurveyAudienceRolloutPercentage(filters?: FeatureFlagFilters | null): number | null {
    if (!filters || filters.groups.length !== 1) {
        return null
    }

    return filters.groups[0].rollout_percentage ?? 100
}

export function isSimpleSurveyAudienceTargeting(filters?: FeatureFlagFilters | null): boolean {
    if (!filters) {
        return true
    }

    if (filters.groups.length !== 1 || filters.aggregation_group_type_index != null || filters.feature_enrollment) {
        return false
    }

    if (filters.multivariate?.variants?.length) {
        return false
    }

    const [group] = filters.groups

    if (group.aggregation_group_type_index != null || group.variant != null) {
        return false
    }

    return (group.properties || []).every(
        (property) => property.type === PropertyFilterType.Person || property.type === PropertyFilterType.Cohort
    )
}

export function getSurveyAudienceSummaryValue(survey: Survey | NewSurvey): string | null {
    const filters = getSurveyTargetingFilters(survey)

    if (!filters) {
        return null
    }

    if (!isSimpleSurveyAudienceTargeting(filters)) {
        return 'Advanced audience targeting'
    }

    const ruleCount = getSurveyAudienceRuleCount(filters)
    const rolloutPercentage = getSurveyAudienceRolloutPercentage(filters)
    const isPartialRollout = rolloutPercentage != null && rolloutPercentage < 100

    if (ruleCount === 0 && !isPartialRollout) {
        return null
    }

    if (ruleCount === 0) {
        return `${rolloutPercentage}% of matching users`
    }

    if (isPartialRollout) {
        return `${ruleCount} audience rule${ruleCount === 1 ? '' : 's'} · ${rolloutPercentage}% shown`
    }

    return `${ruleCount} audience rule${ruleCount === 1 ? '' : 's'}`
}

export function getSurveyCollectionLimitSummary(survey: Survey | NewSurvey): SurveyCollectionLimitSummary | null {
    if (survey.responses_limit && survey.responses_limit > 0) {
        return {
            label: 'Response limit',
            value: String(survey.responses_limit),
        }
    }

    if (
        survey.response_sampling_limit &&
        survey.response_sampling_limit > 0 &&
        survey.response_sampling_interval &&
        survey.response_sampling_interval > 0 &&
        survey.response_sampling_interval_type
    ) {
        return {
            label: 'Sampling limit',
            value: `${survey.response_sampling_limit} / ${survey.response_sampling_interval} ${survey.response_sampling_interval_type}`,
        }
    }

    return null
}

export function getSurveyDisplayConditionsSummary(survey: Survey | NewSurvey): SurveyConditionSummary[] {
    const parts: SurveyConditionSummary[] = []
    const conditions = survey.conditions

    if (conditions?.url) {
        parts.push({
            type: 'url',
            label: 'URL',
            value: `${conditions.urlMatchType === 'exact' ? 'is' : 'contains'} "${conditions.url}"`,
        })
    }
    if (conditions?.selector) {
        parts.push({ type: 'selector', label: 'Selector', value: conditions.selector })
    }
    if (conditions?.deviceTypes?.length) {
        parts.push({ type: 'device', label: 'Device', value: conditions.deviceTypes.join(', ') })
    }
    if ((conditions?.events?.values?.length ?? 0) > 0) {
        parts.push({
            type: 'events',
            label: 'Events',
            value: conditions!.events!.values.map((e) => e.name).join(', '),
        })
    }
    if ((conditions?.actions?.values?.length ?? 0) > 0) {
        parts.push({
            type: 'actions',
            label: 'Actions',
            value: conditions!.actions!.values.map((a) => a.name).join(', '),
        })
    }
    if (survey.linked_flag?.key) {
        parts.push({
            type: 'flag',
            label: 'Feature flag',
            value: survey.linked_flag.key,
            href: urls.featureFlag(survey.linked_flag.id),
        })
    } else if (survey.linked_flag_id) {
        parts.push({ type: 'flag', label: 'Feature flag', value: 'Linked' })
    }
    if ((survey.linked_flag || survey.linked_flag_id) && conditions?.linkedFlagVariant) {
        parts.push({ type: 'flag_variant', label: 'Variant', value: conditions.linkedFlagVariant })
    }
    const audienceSummary = getSurveyAudienceSummaryValue(survey)
    if (audienceSummary) {
        parts.push({ type: 'targeting', label: 'Targeting', value: audienceSummary })
    }
    if (conditions?.seenSurveyWaitPeriodInDays) {
        parts.push({
            type: 'wait_period',
            label: 'Wait period',
            value: `${conditions.seenSurveyWaitPeriodInDays} days since last survey`,
        })
    }

    return parts
}

/**
 * True when posthog-js emits an intermediate `survey sent` event per answered question, sharing one
 * `$survey_submission_id`, with only the last carrying `$survey_completed: true`. Requiring the
 * property to be `true` is what keeps a notification from firing once per question, so it is only
 * worth requiring here.
 *
 * An API survey has no posthog-js rendering it. The integrator sends one event per submission from
 * their own code and marks a partial one with an explicit `$survey_completed: false`, the way
 * posthog-js does, so absent means completed there whatever `enable_partial_responses` says.
 */
export function surveyEmitsPartialSentEvents(survey: Pick<Survey, 'type' | 'enable_partial_responses'>): boolean {
    return (survey.enable_partial_responses ?? false) && survey.type !== SurveyType.API
}

/**
 * Without intermediate partial events, posthog-js has no partial submission to distinguish a
 * complete one from, so it never sets `$survey_completed` and requiring `= true` matches nothing.
 * Accept the property being absent as completed too, the same way the response summary counts
 * legacy events. An explicit `false` stays excluded from sent-event notifications.
 */
export function getSurveyNotificationFilters(
    surveyId: string,
    emitsPartialSentEvents: boolean,
    extraSentEventProperties: EventPropertyFilter[] = []
): CyclotronJobFiltersType {
    const surveyIdProperty: EventPropertyFilter = {
        key: SurveyEventProperties.SURVEY_ID,
        type: PropertyFilterType.Event,
        value: surveyId,
        operator: PropertyOperator.Exact,
    }
    const sentEventProperties: EventPropertyFilter[] = [
        surveyIdProperty,
        {
            key: SurveyEventProperties.SURVEY_COMPLETED,
            type: PropertyFilterType.Event,
            value: true,
            operator: PropertyOperator.Exact,
        },
        ...extraSentEventProperties,
    ]
    // Event entries are OR'd, so a second branch is how "absent or true" is expressed with
    // plain property filters rather than a hand-written HogQL predicate.
    const completedUnsetEventProperties: EventPropertyFilter[] = [
        surveyIdProperty,
        {
            key: SurveyEventProperties.SURVEY_COMPLETED,
            type: PropertyFilterType.Event,
            value: PropertyOperator.IsNotSet,
            operator: PropertyOperator.IsNotSet,
        },
        ...extraSentEventProperties,
    ]

    return {
        events: [
            {
                id: SurveyEventName.SENT,
                type: 'events',
                properties: sentEventProperties,
            },
            ...(emitsPartialSentEvents
                ? []
                : [
                      {
                          id: SurveyEventName.SENT,
                          type: 'events' as const,
                          properties: completedUnsetEventProperties,
                      },
                  ]),
            {
                id: SurveyEventName.DISMISSED,
                type: 'events',
                properties: [
                    {
                        key: SurveyEventProperties.SURVEY_ID,
                        type: PropertyFilterType.Event,
                        value: surveyId,
                        operator: PropertyOperator.Exact,
                    },
                    {
                        key: SurveyEventProperties.SURVEY_PARTIALLY_COMPLETED,
                        type: PropertyFilterType.Event,
                        value: true,
                        operator: PropertyOperator.Exact,
                    },
                ],
            },
        ],
    }
}
````

### FILE: `internal/posthognps/nps.go`

```yaml
block_id: "GO-CUSTOMER-SURVEY-API-NPS:file5:v1"
operation: CREATE
provenance: ADAPTED
source: "PostHog/posthog 6fafbb9081bd15e79448af5650e02a4f9ea435cc MIT core; exact source and adaptation boundaries in third_party/posthog-nps/source-lock.json"
license: "MIT"
sha256: "0132848b6af8b978e739726fc2f57e04df8d65e7b62814e472e0a95d134c032f"
variables: []
secrets_allowed: false
```

````go
// Adapted from PostHog/posthog, MIT, commit 6fafbb9081bd15e79448af5650e02a4f9ea435cc.
// Original: frontend/src/scenes/surveys/utils.ts, calculateNPSFromRawData.
// See docs/posthog-nps.md for the exact adaptation boundary and preserved notices.
package posthognps

// Breakdown uses the existing database's nonnegative int64 grouped counts.
type Breakdown struct {
	Total, Promoters, Passives, Detractors int64
	Score                                  float64
}

// Calculate translates the upstream aggregation after SQL has grouped 0..6,
// 7..8, and 9..10. Caller must validate total >= promoters + detractors >= 0.
// Presentation rounding is deliberately left to the consumer: the existing API
// returns an unrounded number, unlike upstream's one-decimal display string.
func Calculate(total, promoters, detractors int64) Breakdown {
	if total == 0 {
		return Breakdown{}
	}
	return Breakdown{
		Total: total, Promoters: promoters, Passives: total - promoters - detractors,
		Detractors: detractors,
		Score:      float64(promoters-detractors) / float64(total) * 100,
	}
}
````

### FILE: `internal/posthognps/nps_fuzz_test.go`

```yaml
block_id: "GO-CUSTOMER-SURVEY-API-NPS:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "4b75d99b116377ac26896c399f720b9a3ad2cd7d0bf3b2706cb356b60e9f429d"
variables: []
secrets_allowed: false
```

````go
// AUTHORED domain oracle and boundary qualification for the declared Go adaptation.
package posthognps

import (
	"math"
	"math/big"
	"testing"
)

func checkOracle(t *testing.T, n, p, d int64) {
	t.Helper()
	g := Calculate(n, p, d)
	if g.Total != n || g.Promoters != p || g.Detractors != d || g.Passives != n-p-d {
		t.Fatalf("lost population: %+v", g)
	}
	if n == 0 {
		if g.Score != 0 {
			t.Fatal(g)
		}
		return
	}
	exact := new(big.Rat).SetFrac(big.NewInt(p-d), big.NewInt(n))
	exact.Mul(exact, big.NewRat(100, 1))
	want, _ := exact.Float64()
	if math.IsNaN(g.Score) || math.IsInf(g.Score, 0) || g.Score < -100 || g.Score > 100 || math.Abs(g.Score-want) > 3e-14 {
		t.Fatalf("score %.17g oracle %.17g", g.Score, want)
	}
}
func TestContractPrecisionAndInt64(t *testing.T) {
	for _, v := range [][3]int64{{3, 1, 0}, {3, 0, 1}, {math.MaxInt64, math.MaxInt64, 0}, {math.MaxInt64, 0, math.MaxInt64}, {math.MaxInt64, math.MaxInt64 / 2, math.MaxInt64 / 2}} {
		checkOracle(t, v[0], v[1], v[2])
	}
	if Calculate(3, 1, 0).Score == 33.3 {
		t.Fatal("presentation rounding changed API")
	}
}
func FuzzGroupedCounts(f *testing.F) {
	for _, v := range [][3]uint64{{0, 0, 0}, {17, 6, 7}, {3, 1, 0}, {10, 10, 0}, {14, 0, 14}, {math.MaxInt64, math.MaxInt64, 0}} {
		f.Add(v[0], v[1], v[2])
	}
	f.Fuzz(func(t *testing.T, a, b, c uint64) {
		n := int64(a & math.MaxInt64)
		p := int64(b % (uint64(n) + 1))
		d := int64(c % (uint64(n-p) + 1))
		checkOracle(t, n, p, d)
	})
}
````

### FILE: `internal/posthognps/nps_test.go`

```yaml
block_id: "GO-CUSTOMER-SURVEY-API-NPS:file7:v1"
operation: CREATE
provenance: ADAPTED
source: "PostHog/posthog 6fafbb9081bd15e79448af5650e02a4f9ea435cc MIT core; exact source and adaptation boundaries in third_party/posthog-nps/source-lock.json"
license: "MIT"
sha256: "3e606d755bab78d6d8efab67c868258c200e81edb42da72b3c19c346f31afb8b"
variables: []
secrets_allowed: false
```

````go
// Adapted selected PostHog MIT calculateNpsBreakdown examples; see source-lock.
// Invalid/missing histogram parsing cases belong to the upstream TS harness,
// because the local caller receives typed, validated PostgreSQL counts.
package posthognps

import (
	"fmt"
	"testing"
)

func TestUpstreamBreakdownExamples(t *testing.T) {
	cases := []struct {
		name        string
		total, p, d int64
		want        Breakdown
		display     string
	}{
		{"zero total", 0, 0, 0, Breakdown{}, "0.0"},
		{"mixed", 17, 6, 7, Breakdown{17, 6, 4, 7, -100.0 / 17}, "-5.9"},
		{"all zero bins", 0, 0, 0, Breakdown{}, "0.0"},
		{"only promoters", 10, 10, 0, Breakdown{10, 10, 0, 0, 100}, "100.0"},
		{"only passives", 10, 0, 0, Breakdown{10, 0, 10, 0, 0}, "0.0"},
		{"only detractors", 14, 0, 14, Breakdown{14, 0, 0, 14, -100}, "-100.0"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := Calculate(c.total, c.p, c.d)
			if g.Total != c.want.Total || g.Promoters != c.want.Promoters || g.Passives != c.want.Passives || g.Detractors != c.want.Detractors || fmt.Sprintf("%.1f", g.Score) != c.display {
				t.Fatalf("got %+v want %+v display %s", g, c.want, c.display)
			}
		})
	}
}
````

### FILE: `third_party/posthog-nps/source-lock.json`

```yaml
block_id: "GO-CUSTOMER-SURVEY-API-NPS:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "536596b594484460b396be0214d78e0bd81bd45f3e9ea4f6582374a5dfcbe153"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-posthog-nps-source-lock/v1",
  "commit": "6fafbb9081bd15e79448af5650e02a4f9ea435cc",
  "release": "desktop-v0.61.382",
  "license_expression": "MIT",
  "upstream_sources": [
    {
      "path": "frontend/src/scenes/surveys/utils.ts",
      "url": "https://raw.githubusercontent.com/PostHog/posthog/6fafbb9081bd15e79448af5650e02a4f9ea435cc/frontend/src/scenes/surveys/utils.ts",
      "git_blob_sha1": "98066cd409963769cc829dfbe79a6f3afb7f6138",
      "bytes": 62358,
      "sha256": "a43ed06c7eaffbe5bde3fe515177e512ebc4e6bce24dbef253ed100cbb1a08c1",
      "git_blob_verified": true,
      "observed_at": "2026-09-13T14:46:44.868169+00:00"
    },
    {
      "path": "frontend/src/scenes/surveys/utils.test.ts",
      "url": "https://raw.githubusercontent.com/PostHog/posthog/6fafbb9081bd15e79448af5650e02a4f9ea435cc/frontend/src/scenes/surveys/utils.test.ts",
      "git_blob_sha1": "339afd4b6b6c8b8fa34eff405cc212f9090552d2",
      "bytes": 74584,
      "sha256": "7b3b24fe2d6d4525113657629991e0d0056b2c8fcf91e6bd2a906b1c280bd565",
      "git_blob_verified": true,
      "observed_at": "2026-09-13T14:46:45.204514+00:00"
    },
    {
      "path": "LICENSE",
      "url": "https://raw.githubusercontent.com/PostHog/posthog/6fafbb9081bd15e79448af5650e02a4f9ea435cc/LICENSE",
      "git_blob_sha1": "847d9d3def03458035447bdae131e1981e2a63da",
      "bytes": 1563,
      "sha256": "6d82d67dba42eb94ba10f1e986d2eec338c22fb7c5216c2c0ebdecd83d53a029",
      "git_blob_verified": true,
      "observed_at": "2026-09-13T14:46:45.501129+00:00"
    }
  ],
  "selected_implementation_lines": [
    [
      300,
      316
    ],
    [
      340,
      350
    ],
    [
      354,
      372
    ],
    [
      379,
      382
    ]
  ],
  "selected_original_test_lines": [
    843,
    990
  ],
  "original_selected_cases": 9,
  "original_test_harness_sha256": "0d118812a072a6e7dcfb30f50080e3e6b24c08801fa36ca9bb4ea53717c20d34",
  "original_test_log_sha256": "fa8153d15625aabc91e3e7e82343767d81eac4bab9dc629ea94811cedc25068d",
  "adaptation": "Grouped int64 SQL counts; derive passives; numeric unrounded output, no legacy histogram parser. Strict invalid-count validation and reporting minimum remain local caller glue.",
  "runtime_dependencies_added": 0,
  "files": {
    "docs/posthog-nps.md": "2a48b97048c3a11d9ae8a5501c0e980cbff6f619f1cabc89e8bd526f780e0964",
    "third_party/posthog-nps/LICENSE": "6d82d67dba42eb94ba10f1e986d2eec338c22fb7c5216c2c0ebdecd83d53a029",
    "third_party/posthog-nps/utils.test.ts.txt": "7b3b24fe2d6d4525113657629991e0d0056b2c8fcf91e6bd2a906b1c280bd565",
    "third_party/posthog-nps/utils.ts.txt": "a43ed06c7eaffbe5bde3fe515177e512ebc4e6bce24dbef253ed100cbb1a08c1",
    "internal/posthognps/nps.go": "0132848b6af8b978e739726fc2f57e04df8d65e7b62814e472e0a95d134c032f",
    "internal/posthognps/nps_fuzz_test.go": "4b75d99b116377ac26896c399f720b9a3ad2cd7d0bf3b2706cb356b60e9f429d",
    "internal/posthognps/nps_test.go": "3e606d755bab78d6d8efab67c868258c200e81edb42da72b3c19c346f31afb8b"
  },
  "conditions": [
    "Only validated nonnegative grouped counts; no inference about representative population/consent",
    "No whole upstream runtime, analytics, enterprise code or production certification",
    "Rebuild complete reference and run connected customer feedback regression before promotion"
  ]
}
````


V402319: PostHog MIT fixed-source NPS qualified; composed integration evidence governs target use. See docs/posthog-nps.md.

V402 composed delta: PostHog NPS G0-G8 USE_REUSABLE_PACK; bounded arithmetic selected, existing API caller/threshold unchanged. NPS_SOURCE_ADAPTATION_V402.md/json.
