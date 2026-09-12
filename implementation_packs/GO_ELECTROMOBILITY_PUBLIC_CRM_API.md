# Go Electromobility Public Catalog and CRM API

## 1. Metadata

```yaml
pack_id: "GO-ELECTROMOBILITY-PUBLIC-CRM-API"
pack_version: "0.2.3"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Añade catálogo público, alta administrativa de modelos y captura pública de leads con consentimiento, idempotencia durable y outbox."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "pgx 5.10.0"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "ELECTROMOBILITY-FRANCHISE-MODULES 0.1.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner AND MIT dependency"
upstream_sources: ["https://go.dev", "https://www.postgresql.org/docs/18/", "https://github.com/jackc/pgx"]
verified_at: "2026-08-25"
```

Bloques `AUTHORED`. El tenant público se resuelve por `tenant_code`; la organización receptora por `organization_code`. El endpoint administrativo deriva tenant y permisos del token verificado. La captura pública exige `Idempotency-Key`, persiste request hash/replay y todavía necesita rate limiting/antiabuse en el edge del proyecto.

## 2. Applicability

Use for a public vehicle catalog, authorized model administration and consented lead capture backed by the enterprise Go/PostgreSQL composition. Reject it when another CRM/catalog is authoritative unless a synchronization contract is designed. Public reachability requires a project-specific antiabuse and consent policy.

## 3. Architecture contract

Public tenant/organization codes are resolved server-side; protected tenant and permissions derive only from the verified token. Lead creation uses a required idempotency key, request hash, transaction and outbox. Catalog publication is explicit. Cross-tenant and cross-organization access fails in HTTP and repository predicates. Contact data is minimized and remains subject to retention/deletion policy.

## 4. Exact file manifest

```text
CREATE internal/electromobility/service.go
CREATE internal/platform/postgres/electromobility.go
CREATE internal/platform/postgres/electromobility_integration_test.go
CREATE internal/platform/httpapi/electromobility.go
CREATE internal/platform/httpapi/electromobility_test.go
```

## 5. Materialization blocks

### FILE: `internal/electromobility/service.go`

```yaml
block_id: "GO-EM-API:service:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "e4c13302f6dab9a1651411a00909984da11de2c594aef5b6483f5b166e9a3204"
variables: []
secrets_allowed: false
```

````go
package electromobility

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

var ErrNotFound = errors.New("resource not found")
var ErrConflict = errors.New("electromobility conflict")
var codePattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`)

type Model struct {
	ID            string          `json:"id"`
	Code          string          `json:"code"`
	DisplayName   string          `json:"displayName"`
	VehicleClass  string          `json:"vehicleClass"`
	Specification json.RawMessage `json:"specification"`
}
type Lead struct {
	ID, TenantID, OrganizationID, ModelID, SourceCode string
	Contact                                           json.RawMessage
	ConsentID, ConsentEvidenceHash                    string
}
type Repository interface {
	ListPublicModels(context.Context, string) ([]Model, error)
	CreateModel(context.Context, string, string, Model) error
	ResolvePublicOrganization(context.Context, string, string) (string, string, error)
	CreateLead(context.Context, Lead, string, string) (string, bool, error)
}
type IDGenerator interface{ New() string }
type Service struct {
	repository Repository
	ids        IDGenerator
}

func NewService(repository Repository, ids IDGenerator) *Service {
	return &Service{repository: repository, ids: ids}
}

func (s *Service) ListPublicModels(ctx context.Context, tenantCode string) ([]Model, error) {
	if !codePattern.MatchString(tenantCode) {
		return nil, fmt.Errorf("invalid tenant code")
	}
	return s.repository.ListPublicModels(ctx, tenantCode)
}
func (s *Service) CreateModel(ctx context.Context, tenantID string, input Model) (Model, error) {
	classes := map[string]bool{"motorcycle": true, "bicycle": true, "scooter": true, "utility": true, "other": true}
	if tenantID == "" || !codePattern.MatchString(input.Code) || len(input.DisplayName) < 2 || !classes[input.VehicleClass] || !validObject(input.Specification) {
		return Model{}, fmt.Errorf("invalid model")
	}
	input.ID = s.ids.New()
	eventID := s.ids.New()
	if err := s.repository.CreateModel(ctx, tenantID, eventID, input); err != nil {
		return Model{}, err
	}
	return input, nil
}
func (s *Service) CaptureLead(ctx context.Context, tenantCode, organizationCode, modelID, source string, contact json.RawMessage, consentHash, idempotencyKey string) (Lead, bool, error) {
	if !codePattern.MatchString(tenantCode) || !codePattern.MatchString(organizationCode) || !codePattern.MatchString(source) || !validObject(contact) || !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(consentHash) || len(idempotencyKey) < 16 || len(idempotencyKey) > 128 {
		return Lead{}, false, fmt.Errorf("invalid lead")
	}
	tenantID, organizationID, err := s.repository.ResolvePublicOrganization(ctx, tenantCode, organizationCode)
	if err != nil {
		return Lead{}, false, err
	}
	lead := Lead{ID: s.ids.New(), TenantID: tenantID, OrganizationID: organizationID, ModelID: modelID, SourceCode: source, Contact: contact, ConsentID: s.ids.New(), ConsentEvidenceHash: consentHash}
	resourceID, replayed, err := s.repository.CreateLead(ctx, lead, s.ids.New(), idempotencyKey)
	if err != nil {
		return Lead{}, false, err
	}
	lead.ID = resourceID
	return lead, replayed, nil
}
func validObject(value json.RawMessage) bool {
	var object map[string]any
	return len(value) > 1 && json.Unmarshal(value, &object) == nil && object != nil
}
````

### FILE: `internal/platform/postgres/electromobility.go`

```yaml
block_id: "GO-EM-API:repository:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "45977bd13125515b62ca85ca838e69ff6305a5e3d44ea2837784e33fdb00fe6b"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"elite.local/enterprise/internal/electromobility"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Electromobility struct{ pool *pgxpool.Pool }

func NewElectromobility(pool *pgxpool.Pool) *Electromobility { return &Electromobility{pool: pool} }

func (r *Electromobility) ListPublicModels(ctx context.Context, tenantCode string) ([]electromobility.Model, error) {
	rows, err := r.pool.Query(ctx, `select m.model_id,m.model_code,m.display_name,m.vehicle_class,m.specification from catalog.vehicle_model m join platform.tenant t on t.tenant_id=m.tenant_id where t.tenant_code=$1 and t.status='active' and m.lifecycle_state='active' and m.publicly_visible=true order by m.display_name,m.model_id`, tenantCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	models := []electromobility.Model{}
	for rows.Next() {
		var model electromobility.Model
		if err := rows.Scan(&model.ID, &model.Code, &model.DisplayName, &model.VehicleClass, &model.Specification); err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	return models, rows.Err()
}
func (r *Electromobility) CreateModel(ctx context.Context, tenantID, eventID string, model electromobility.Model) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state,publicly_visible,specification)values($1,$2,$3,$4,$5,'draft',false,$6)`, tenantID, model.ID, model.Code, model.DisplayName, model.VehicleClass, model.Specification)
	if err != nil {
		return err
	}
	payload, _ := json.Marshal(model)
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'vehicle-model',$3,1,'vehicle-model.created',1,clock_timestamp(),$4)`, tenantID, eventID, model.ID, payload)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Electromobility) ResolvePublicOrganization(ctx context.Context, tenantCode, organizationCode string) (string, string, error) {
	var tenantID, organizationID string
	err := r.pool.QueryRow(ctx, `select t.tenant_id,o.organization_id from platform.tenant t join org.organization o on o.tenant_id=t.tenant_id where t.tenant_code=$1 and t.status='active' and o.organization_code=$2 and o.status='active'`, tenantCode, organizationCode).Scan(&tenantID, &organizationID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", electromobility.ErrNotFound
	}
	return tenantID, organizationID, err
}
func (r *Electromobility) CreateLead(ctx context.Context, lead electromobility.Lead, eventID, idempotencyKey string) (string, bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", false, err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at)values($1,'public-lead',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict do nothing`, lead.TenantID, idempotencyKey, lead.ConsentEvidenceHash)
	if err != nil {
		return "", false, err
	}
	if result.RowsAffected() == 0 {
		var requestHash, status, resourceID string
		err = tx.QueryRow(ctx, `select request_sha256_hex,status,coalesce(resource_id,'') from platform.idempotency_record where tenant_id=$1 and scope='public-lead' and idempotency_key=$2`, lead.TenantID, idempotencyKey).Scan(&requestHash, &status, &resourceID)
		if err != nil {
			return "", false, err
		}
		if requestHash != lead.ConsentEvidenceHash || status != "completed" || resourceID == "" {
			return "", false, electromobility.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return "", false, err
		}
		return resourceID, true, nil
	}
	var model any = nil
	if lead.ModelID != "" {
		model = lead.ModelID
	}
	_, err = tx.Exec(ctx, `insert into crm.lead(tenant_id,lead_id,organization_id,model_id,lifecycle_state,source_code,contact_payload,consent_required)values($1,$2,$3,$4,'new',$5,$6,true)`, lead.TenantID, lead.ID, lead.OrganizationID, model, lead.SourceCode, lead.Contact)
	if err != nil {
		return "", false, err
	}
	_, err = tx.Exec(ctx, `insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex)values($1,$2,$3,'sales-contact','v1','granted',clock_timestamp(),$4)`, lead.TenantID, lead.ConsentID, lead.ID, lead.ConsentEvidenceHash)
	if err != nil {
		return "", false, err
	}
	payload, _ := json.Marshal(map[string]string{"lead_id": lead.ID, "organization_id": lead.OrganizationID, "source_code": lead.SourceCode})
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'lead',$3,1,'lead.captured',1,clock_timestamp(),$4)`, lead.TenantID, eventID, lead.ID, payload)
	if err != nil {
		return "", false, fmt.Errorf("write lead outbox: %w", err)
	}
	_, err = tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=202,response_body=jsonb_build_object('lead_id',$3::text,'status','accepted'),resource_type='lead',resource_id=$3,locked_until=null where tenant_id=$1 and scope='public-lead' and idempotency_key=$2 and status='processing'`, lead.TenantID, idempotencyKey, lead.ID)
	if err != nil {
		return "", false, err
	}
	if err = tx.Commit(ctx); err != nil {
		return "", false, err
	}
	return lead.ID, false, nil
}
````

### FILE: `internal/platform/postgres/electromobility_integration_test.go`

```yaml
block_id: "GO-EM-API:repository-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "6debd5bae85e045983dda112fd8b48d5bd0731ae101b15d4769dab934ec859b2"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"elite.local/enterprise/internal/electromobility"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
)

func TestElectromobilityRepositoryCatalogAndLead(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28e01"
	org := "em-api-org"
	cleanup := func() {
		_, _ = pool.Exec(ctx, `delete from platform.outbox_event where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from platform.idempotency_record where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from crm.consent_evidence where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from crm.lead where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from catalog.vehicle_model where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from org.organization where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	}
	cleanup()
	defer cleanup()
	_, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'em-api','EM API','EM')`, tenant)
	if err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,$2,'central-store','Central','store')`, tenant, org)
	if err != nil {
		t.Fatal(err)
	}
	repo := NewElectromobility(pool)
	model := electromobility.Model{ID: "model-1", Code: "urban-one", DisplayName: "Urban One", VehicleClass: "bicycle", Specification: json.RawMessage(`{"range_km":80}`)}
	if err := repo.CreateModel(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28e02", model); err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(ctx, `update catalog.vehicle_model set lifecycle_state='active',publicly_visible=true where tenant_id=$1 and model_id=$2`, tenant, model.ID)
	if err != nil {
		t.Fatal(err)
	}
	models, err := repo.ListPublicModels(ctx, "em-api")
	if err != nil || len(models) != 1 {
		t.Fatalf("models=%+v err=%v", models, err)
	}
	tenantID, orgID, err := repo.ResolvePublicOrganization(ctx, "em-api", "central-store")
	if err != nil || tenantID != tenant || orgID != org {
		t.Fatalf("resolve %s %s %v", tenantID, orgID, err)
	}
	lead := electromobility.Lead{ID: "lead-1", TenantID: tenant, OrganizationID: org, ModelID: model.ID, SourceCode: "public-web", Contact: json.RawMessage(`{"email":"person@example.test"}`), ConsentID: "consent-1", ConsentEvidenceHash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
	resourceID, replayed, err := repo.CreateLead(ctx, lead, "018f4d4a-7b36-7a21-8d10-2f4c54c28e03", "lead-request-00000001")
	if err != nil || replayed || resourceID != lead.ID {
		t.Fatal(err)
	}
	resourceID, replayed, err = repo.CreateLead(ctx, lead, "018f4d4a-7b36-7a21-8d10-2f4c54c28e04", "lead-request-00000001")
	if err != nil || !replayed || resourceID != lead.ID {
		t.Fatalf("replay id=%s replayed=%t err=%v", resourceID, replayed, err)
	}
	changed := lead
	changed.ConsentEvidenceHash = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	if _, _, err = repo.CreateLead(ctx, changed, "018f4d4a-7b36-7a21-8d10-2f4c54c28e05", "lead-request-00000001"); !errors.Is(err, electromobility.ErrConflict) {
		t.Fatalf("changed request err=%v", err)
	}
	var leads, consents, events int
	_ = pool.QueryRow(ctx, `select count(*) from crm.lead where tenant_id=$1`, tenant).Scan(&leads)
	_ = pool.QueryRow(ctx, `select count(*) from crm.consent_evidence where tenant_id=$1`, tenant).Scan(&consents)
	_ = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1`, tenant).Scan(&events)
	if leads != 1 || consents != 1 || events != 2 {
		t.Fatalf("lead=%d consent=%d events=%d", leads, consents, events)
	}
}
````

### FILE: `internal/platform/httpapi/electromobility.go`

```yaml
block_id: "GO-EM-API:http:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "bbcc571cb56f363cdc36b29992592252c5cb4ad7a70c9771054493f7dd23741c"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"crypto/sha256"
	"elite.local/enterprise/internal/electromobility"
	"elite.local/enterprise/internal/order"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type ElectromobilityAPI struct {
	service  *electromobility.Service
	verifier identity.Verifier
}

type EnterpriseModule interface {
	Register(*http.ServeMux, identity.Verifier)
}

func NewEnterprise(orders *order.Service, verifier identity.Verifier, service *electromobility.Service, modules ...EnterpriseModule) http.Handler {
	api := &ElectromobilityAPI{service: service, verifier: verifier}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/public/{tenantCode}/models", api.listModels)
	mux.HandleFunc("POST /v1/public/{tenantCode}/{organizationCode}/leads", api.captureLead)
	mux.HandleFunc("POST /v1/catalog/models", api.createModel)
	for _, module := range modules {
		if module != nil {
			module.Register(mux, verifier)
		}
	}
	mux.Handle("/", New(orders, verifier))
	return recoverMiddleware(mux)
}
func (a *ElectromobilityAPI) listModels(w http.ResponseWriter, r *http.Request) {
	models, err := a.service.ListPublicModels(r.Context(), r.PathValue("tenantCode"))
	if err != nil {
		writeProblem(w, 400, "INVALID_TENANT", "tenant code is invalid")
		return
	}
	writeJSON(w, 200, map[string]any{"models": models})
}
func (a *ElectromobilityAPI) createModel(w http.ResponseWriter, r *http.Request) {
	principal, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return
	}
	if !principal.Allowed("catalog:write") {
		writeProblem(w, 403, "FORBIDDEN", "catalog:write permission is required")
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return
	}
	var input electromobility.Model
	if !decodeStrict(w, r, &input) {
		return
	}
	model, err := a.service.CreateModel(r.Context(), principal.TenantID, input)
	if err != nil {
		writeProblem(w, 400, "INVALID_MODEL", "model does not match the contract")
		return
	}
	writeJSON(w, 201, model)
}
func (a *ElectromobilityAPI) captureLead(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 64<<10))
	if err != nil {
		writeProblem(w, 400, "INVALID_BODY", "body is invalid or too large")
		return
	}
	idempotencyKey := r.Header.Get("Idempotency-Key")
	if len(idempotencyKey) < 16 || len(idempotencyKey) > 128 {
		writeProblem(w, 400, "IDEMPOTENCY_KEY_REQUIRED", "Idempotency-Key must contain 16 to 128 characters")
		return
	}
	var input struct {
		ModelID        string          `json:"model_id"`
		SourceCode     string          `json:"source_code"`
		Contact        json.RawMessage `json:"contact"`
		ConsentGranted bool            `json:"consent_granted"`
	}
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&input) != nil || !input.ConsentGranted {
		writeProblem(w, 400, "CONSENT_REQUIRED", "valid explicit consent is required")
		return
	}
	hash := sha256.Sum256(body)
	lead, replayed, err := a.service.CaptureLead(r.Context(), r.PathValue("tenantCode"), r.PathValue("organizationCode"), input.ModelID, input.SourceCode, input.Contact, hex.EncodeToString(hash[:]), idempotencyKey)
	if err != nil {
		if errors.Is(err, electromobility.ErrConflict) {
			writeProblem(w, 409, "IDEMPOTENCY_CONFLICT", "idempotency key is processing or belongs to another request")
			return
		}
		writeProblem(w, 400, "INVALID_LEAD", "lead does not match the contract")
		return
	}
	status := 202
	if replayed {
		status = 200
		w.Header().Set("Idempotency-Replayed", "true")
	}
	writeJSON(w, status, map[string]string{"lead_id": lead.ID, "status": "accepted"})
}
func decodeStrict(w http.ResponseWriter, r *http.Request, target any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
	decoder.DisallowUnknownFields()
	if decoder.Decode(target) != nil {
		writeProblem(w, 400, "INVALID_BODY", "body does not match the contract")
		return false
	}
	return true
}
````

### FILE: `internal/platform/httpapi/electromobility_test.go`

```yaml
block_id: "GO-EM-API:http-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "0439847e0cf95c9b8dc40853226484b2332381b5cc311c9fddcb7622983d951f"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"elite.local/enterprise/internal/electromobility"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// AUTHORED integration harness: real admitted handler/repository and Microsoft
// Playwright runtime. No successful authentication double, provider or production claim.
type browserDenyVerifier struct{}

func (browserDenyVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{}, identity.ErrUnauthenticated
}

type publicBrowserFixture struct {
	testName      string
	setup         func(context.Context, *testing.T, *pgxpool.Pool, string) []EnterpriseModule
	verify        func(context.Context, *testing.T, *pgxpool.Pool, string)
	cleanupTables []string
}

func TestPublicLeadBrowserPostgres(t *testing.T) {
	runPublicBrowserPostgres(t, publicBrowserFixture{testName: "public lead crosses"})
}

func runPublicBrowserPostgres(t *testing.T, fixture publicBrowserFixture) {
	t.Helper()
	if os.Getenv("ELITE_PUBLIC_LEAD_E2E") != "1" {
		t.Skip("explicit disposable browser/Go/PostgreSQL gate not requested")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	config, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal("invalid disposable database configuration")
	}
	if config.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(config.ConnConfig.Database, "elite_browser_") {
		t.Fatal("requires explicit loopback database named elite_browser_*; never use a project database")
	}
	webRoot := os.Getenv("ELITE_WEB_ROOT")
	if !filepath.IsAbs(webRoot) {
		t.Fatal("ELITE_WEB_ROOT must be an absolute, freshly built enterprise-web composition")
	}
	gate := filepath.Join(webRoot, "microsoft_playwright_browser_gate")
	cli := filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js")
	for _, path := range []string{cli, filepath.Join(webRoot, ".next", "BUILD_ID")} {
		if _, err := os.Stat(path); err != nil {
			t.Fatal("missing built web or installed pinned Playwright runtime")
		}
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal("cannot initialize disposable database pool")
	}
	defer pool.Close()
	generator := randomid.Generator{}
	tenant := generator.New()
	code := "browser-" + strings.ReplaceAll(tenant, "-", "")
	_, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,$2,'Browser Fixture','Browser Fixture')`, tenant, code)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		cleanup, done := context.WithTimeout(context.Background(), 10*time.Second)
		defer done()
		// Only rows belonging to the random tenant created above are removed.
		for _, table := range append(fixture.cleanupTables, "platform.outbox_event", "platform.idempotency_record", "crm.consent_evidence", "crm.lead", "catalog.vehicle_model", "org.organization", "platform.tenant") {
			if _, err := pool.Exec(cleanup, "delete from "+table+" where tenant_id=$1", tenant); err != nil {
				t.Errorf("fixture cleanup failed for %s: %v", table, err)
			}
		}
	}()
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'browser-org','browser-store','Browser Store','store')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state,publicly_visible,specification)values($1,'browser-model','browser-model','Browser Fixture Model','bicycle','active',true,'{}')`, tenant); err != nil {
		t.Fatal(err)
	}
	var modules []EnterpriseModule
	if fixture.setup != nil {
		modules = fixture.setup(ctx, t, pool, tenant)
	}
	server := httptest.NewServer(NewEnterprise(nil, browserDenyVerifier{}, electromobility.NewService(postgres.NewElectromobility(pool), generator), modules...))
	defer server.Close()
	command := exec.CommandContext(ctx, "node", cli, "test", "tests/enterprise-web.spec.mjs", "--grep", fixture.testName, "--workers=1", "--retries=0", "--max-failures=1")
	command.Dir = gate
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		switch strings.ToUpper(name) {
		case "DATABASE_URL", "TEST_DATABASE_URL", "ELITE_BASE_URL", "ELITE_RUNTIME_ONLY", "ENTERPRISE_API_BASE_URL", "ENTERPRISE_TENANT_CODE", "ENTERPRISE_ORGANIZATION_CODE":
			continue
		}
		command.Env = append(command.Env, entry)
	}
	command.Env = append(command.Env, "ENTERPRISE_API_BASE_URL="+server.URL, "ENTERPRISE_TENANT_CODE="+code, "ENTERPRISE_ORGANIZATION_CODE=browser-store")
	output, err := command.CombinedOutput()
	t.Log(string(output))
	if err != nil {
		t.Fatalf("connected browser gate failed: %v", err)
	}
	var leads, consents, events, receipts int
	for _, check := range []struct {
		sql   string
		count *int
	}{
		{`select count(*) from crm.lead where tenant_id=$1 and organization_id='browser-org' and model_id='browser-model' and lifecycle_state='new' and source_code='public-web' and contact_payload->>'email' like '%@example.invalid'`, &leads},
		{`select count(*) from crm.consent_evidence where tenant_id=$1 and decision='granted'`, &consents},
		{`select count(*) from platform.outbox_event where tenant_id=$1 and event_type='lead.captured' and not (payload ? 'email') and not (payload ? 'contact_payload')`, &events},
		{`select count(*) from platform.idempotency_record where tenant_id=$1 and scope='public-lead' and status='completed'`, &receipts},
	} {
		if err := pool.QueryRow(ctx, check.sql, tenant).Scan(check.count); err != nil {
			t.Fatal(err)
		}
	}
	if leads != 4 || consents != 4 || events != 4 || receipts != 4 {
		t.Fatalf("durable invariant failed: leads=%d consents=%d events=%d receipts=%d; want four, no duplicates", leads, consents, events, receipts)
	}
	t.Log("PUBLIC_LEAD_BROWSER_POSTGRES_PASS browsers=4 durable_leads=4 consent=4 outbox=4 idempotency=4")
	if fixture.verify != nil {
		fixture.verify(ctx, t, pool, tenant)
	}
}

type emRepo struct {
	models  []electromobility.Model
	created int
	lead    int
}

func (r *emRepo) ListPublicModels(context.Context, string) ([]electromobility.Model, error) {
	return r.models, nil
}
func (r *emRepo) CreateModel(_ context.Context, _ string, _ string, _ electromobility.Model) error {
	r.created++
	return nil
}
func (r *emRepo) ResolvePublicOrganization(context.Context, string, string) (string, string, error) {
	return "018f4d4a-7b36-7a21-8d10-2f4c54c28f01", "org", nil
}
func (r *emRepo) CreateLead(_ context.Context, lead electromobility.Lead, _, _ string) (string, bool, error) {
	r.lead++
	return lead.ID, false, nil
}

type emIDs struct{ n int }

func (i *emIDs) New() string {
	i.n++
	return "018f4d4a-7b36-7a21-8d10-" + []string{"2f4c54c28f10", "2f4c54c28f11", "2f4c54c28f12", "2f4c54c28f13", "2f4c54c28f14", "2f4c54c28f15"}[i.n-1]
}

type emVerifier struct{ principal identity.Principal }

func (v emVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return v.principal, nil
}
func TestElectromobilityHTTPBoundaries(t *testing.T) {
	repo := &emRepo{models: []electromobility.Model{{ID: "m", Code: "urban", DisplayName: "Urban", VehicleClass: "bicycle", Specification: json.RawMessage(`{}`)}}}
	ids := &emIDs{}
	service := electromobility.NewService(repo, ids)
	verifier := emVerifier{principal: identity.Principal{Subject: "admin", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c28f01", Permissions: map[string]struct{}{"catalog:write": {}}}}
	handler := NewEnterprise(nil, verifier, service)
	request := httptest.NewRequest("GET", "/v1/public/acme/models", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatalf("list status=%d", response.Code)
	}
	request = httptest.NewRequest("POST", "/v1/catalog/models", strings.NewReader(`{"code":"urban-one","displayName":"Urban One","vehicleClass":"bicycle","specification":{}}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer valid")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 201 || repo.created != 1 {
		t.Fatalf("create status=%d created=%d body=%s", response.Code, repo.created, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/public/acme/store/leads", strings.NewReader(`{"model_id":"m","source_code":"public-web","contact":{"email":"person@example.test"},"consent_granted":true}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "lead-request-00000001")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 202 || repo.lead != 1 {
		t.Fatalf("lead status=%d leads=%d body=%s", response.Code, repo.lead, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/public/acme/store/leads", strings.NewReader(`{"source_code":"public-web","contact":{},"consent_granted":false}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "lead-request-00000002")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 400 || repo.lead != 1 {
		t.Fatalf("consent status=%d leads=%d", response.Code, repo.lead)
	}
}
````

## 6. Configuration surface

No independent environment variables are introduced. Required inputs are the composed database pool, ID generator and token verifier. Tenant/organization codes, consent flag/version, idempotency key and bounded contact fields are validated per request; policy-specific retention, locales and channels belong to project configuration.

## 7. Dependency bill

| Dependency | Pin | Use | License | Scope | Source |
|---|---|---|---|---|---|
| Go | `1.26.7` | service/HTTP/tests | BSD-3-Clause | build/runtime | `go.dev` |
| PostgreSQL | `18.6` verified baseline | catalog/CRM/outbox | PostgreSQL | runtime/test | `postgresql.org` |
| `pgx` | `5.10.0` | PostgreSQL adapter | MIT | build/runtime | `github.com/jackc/pgx` |

## 8. Apply order

Apply foundation and electromobility migrations first, then materialize service/repository/HTTP files and wire the module. Execute public validation, idempotency replay/conflict and negative scope tests before exposure. Existing systems must map their catalog/CRM ownership. Rollback removes routing first while preserving leads/outbox records; destructive schema rollback is not a production default.

## 9. Verification

V260 añade `TestPublicLeadBrowserPostgres` al test HTTP existente: harness opt-in con API/repositorio reales, BFF de producción construido y Microsoft Playwright 1.62.1 fijado. Cuatro proyectos de navegador envían el formulario y prueban replay/conflicto/consentimiento/idempotencia; PostgreSQL confirma exactamente cuatro leads, consentimientos, eventos y receipts. Configuración/comando/limpieza en el README materializado de `MICROSOFT_PLAYWRIGHT_BROWSER_GATE` 0.1.3. Sólo admite base descartable `elite_browser_*` en 127.0.0.1; el verifier del harness rechaza toda autenticación, no reemplaza el IdP productivo. Sin opt-in el test se omite explícitamente, nunca cuenta como PASS conectado. Evidencia: `reconstruction_evidence/PUBLIC_LEAD_BROWSER_POSTGRES_V260.md`.

La composición limpia con Go core y migrations 0001/0002/0003 pasó `gofmt` sin diff, unit/HTTP/integration tests, vet y build; evidencia en `reconstruction_evidence/GO_ELECTROMOBILITY_PUBLIC_CRM_API_2026-08-24_V1.md`. Condiciones productivas: rate limit/antiabuse distribuido, política de consentimiento y versión por jurisdicción, normalización/validación de contacto, sesión/CORS/CSRF si se usan cookies, authorization E2E, pagination/cache/ETag, observabilidad, roles mínimos y reglas de publicación del catálogo.

## 10. Reconstruction evidence

Clean rebuild and database/HTTP gates are recorded in `reconstruction_evidence/GO_ELECTROMOBILITY_PUBLIC_CRM_API_2026-08-24_V2.md`; final library evidence rechecks version 0.2.1 in the integrated profile.
