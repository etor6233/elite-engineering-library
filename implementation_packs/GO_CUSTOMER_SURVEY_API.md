# Go Customer Survey API

## 1. Metadata

```yaml
pack_id: "GO-CUSTOMER-SURVEY-API"
pack_version: "0.1.1"
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
```

## 5. Materialization blocks

### FILE: `internal/customerfeedback/service.go`
```yaml
block_id: "GO-CUSTOMER-SURVEY-API:internal/customerfeedback/service.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "1d11d2c4a1f8388c18ba46e7510266f2b4c8bce8654cf23dfc28fc6aa7642b79"
variables: []
secrets_allowed: false
```
````go
// AUTHORED customer feedback boundary. Policy texts and survey activation are project-owned.
package customerfeedback

import (
	"context"
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
// This is arithmetic, not proof of representative sampling or legally sufficient consent.
func NPS(count, promoters, detractors, minimum int64) (Summary, error) {
	if count < 0 || promoters < 0 || detractors < 0 || minimum < 1 || promoters > count || detractors > count-promoters {
		return Summary{}, ErrInvalid
	}
	result := Summary{Responses: count}
	if count >= minimum {
		value := float64(promoters-detractors) / float64(count) * 100
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
sha256: "e6560b6e36912494580ea5cff17a863039d8b70280cd7a6777a090890f86630d"
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
sha256: "1bfa151b9d9c3d9f09f8982d44cfbe3743b74c4c89ccc4fbda7c34cb78f7890b"
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
````

## 6. Configuration surface

CUSTOMER_SURVEYS_ENABLED=1 explicitly mounts API routes; features.customer_surveys=true
exposes BFF/pages. Both default to disabled. Existing OIDC/session/API-origin config
is reused. Definitions are inactive by default. Terms/dates/thresholds are immutable;
new policy versions require a new survey ID. DATABASE_URL is an operator secret.

## 7. Dependency bill

No dependency, version or external source added. Exact unchanged go.mod/go.sum and
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
