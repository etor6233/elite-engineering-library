# Go Lead Candidate Promotion

## 1. Metadata

```yaml
pack_id: "GO-LEAD-CANDIDATE-PROMOTION"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Promueve un candidato durable de proveedor a lead CRM únicamente mediante decisión de contacto explícita, evidencia fijada, mapping versionado, idempotencia, procedencia y outbox atómicos."
stacks: ["Go 1.26.7", "PostgreSQL 18.6"]
compatible_with: ["GO-OMNICHANNEL-LEAD-INGRESS 0.2.x", "GO-ELECTROMOBILITY-APPLICATION 0.1.x", "GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.1.x", "POSTGRES-TRANSACTIONAL-FOUNDATION 0.1.x"]
incompatible_with: ["contacto automático por mera recepción", "mapping inferido", "lead de prueba promovido", "evidencia de consentimiento ausente"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://developers.google.com/google-ads/webhook/docs/implementation", "https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html", "https://www.postgresql.org/docs/18/transaction-iso.html"]
verified_at: "2026-09-04"
```

Ningún archivo upstream se copia. Los cinco bloques `ADAPTED` trasladan contratos oficiales estrechos —evento al menos una vez, deduplicación, transacción serializable y procedencia— al modelo local; los dos `AUTHORED` sólo soportan rollback y regresión. El pack no decide si una campaña permite contacto: el proyecto debe suministrar política, finalidad, evidencia y mapping aprobados.

## 2. Applicability

Úselo después de una ingesta durable admitida y antes de cualquier respuesta comercial. Requiere las migraciones base, 0044 y luego 0045. Permanece `CONDITIONED` hasta demostrar políticas, textos, finalidades, mappings y corpus de la cuenta real.

## 3. Architecture contract

- **Ownership**: `leadpromotion` valida decisiones y mappings; PostgreSQL posee la promoción atómica.
- **Invariantes**: recepción no concede contacto; `is_test=true` nunca llega a CRM; no hay mapping por heurística; un replay exacto devuelve el mismo lead y uno divergente falla; lead, consentimiento, procedencia, idempotencia y outbox se confirman en una transacción.
- **Data flow**: candidato 0044 → decisión aprobada del proyecto → mapping versionado → `crm.lead` + `crm.consent_evidence` + `integration.lead_promotion` + outbox.
- **Failure modes**: candidato ausente, test, evidencia/mapping inválidos y replay divergente hacen rollback completo.
- **Seguridad/privacidad**: el outbox sólo publica identidades/metadatos; los campos no mapeados permanecen en la evidencia de acceso restringido; el proyecto debe fijar cifrado, RLS, retención y borrado.
- **Performance budget**: una transacción serializable por promoción; índices por PK/identidad; sin llamadas externas dentro de la transacción.
- **Operación/rollback**: `0045.down.sql` elimina sólo el vínculo de promoción y requiere antes retirar datos del entorno; no revierte decisiones comerciales reales.

## 4. Exact file manifest

```text
CREATE internal/leadpromotion/promotion.go
CREATE internal/leadpromotion/promotion_test.go
CREATE internal/platform/postgres/lead_promotion.go
CREATE internal/platform/postgres/lead_promotion_integration_test.go
CREATE db/migrations/0045_lead_candidate_promotion.up.sql
CREATE db/migrations/0045_lead_candidate_promotion.down.sql
CREATE db/tests/0045_lead_candidate_promotion.test.sql
```

## 5. Materialization blocks

### FILE: `internal/leadpromotion/promotion.go`
```yaml
block_id: "GO-LEAD-CANDIDATE-PROMOTION:internal/leadpromotion/promotion.go:v1"
operation: CREATE
provenance: ADAPTED
source: "official Google Ads webhook, AWS at-least-once/idempotency and PostgreSQL transaction contracts; local implementation"
license: "LicenseRef-Workspace-Owner"
sha256: "8a32c465cd46f9da0ad322080d0ce629f8c6cd38531129de4828ab9eb5f1d88d"
variables: []
secrets_allowed: false
```
````go
// Package leadpromotion converts a durable provider candidate into a CRM lead
// only after an explicit, evidenced contact-policy decision. Ingestion alone
// never grants permission to contact.
package leadpromotion

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	ErrInvalid     = errors.New("leadpromotion: invalid command")
	ErrNotFound    = errors.New("leadpromotion: candidate not found")
	ErrTestLead    = errors.New("leadpromotion: test lead cannot be promoted")
	ErrConflict    = errors.New("leadpromotion: conflicting replay")
	ErrFieldAbsent = errors.New("leadpromotion: configured field absent")
	providerRe     = regexp.MustCompile(`^[a-z][a-z0-9_]{0,31}$`)
	hex64Re        = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

// Command contains only decisions supplied by the project. It deliberately
// has no default policy, field mapping, purpose or consent version.
type Command struct {
	TenantID        string
	Provider        string
	ProviderEventID string
	LeadID          string
	ConsentID       string
	PurposeCode     string
	PolicyVersion   string
	EvidenceSHA256  string
	DecisionAt      time.Time
	MappingVersion  string
	FieldMapping    map[string]string // provider field ID -> email|phone|name
	IdempotencyKey  string
}

func (c Command) Validate() error {
	if strings.TrimSpace(c.TenantID) == "" || !providerRe.MatchString(c.Provider) || strings.TrimSpace(c.ProviderEventID) == "" {
		return fmt.Errorf("%w: source identity", ErrInvalid)
	}
	if strings.TrimSpace(c.LeadID) == "" || strings.TrimSpace(c.ConsentID) == "" || strings.TrimSpace(c.IdempotencyKey) == "" {
		return fmt.Errorf("%w: target identity", ErrInvalid)
	}
	if len(c.IdempotencyKey) < 8 || len(c.IdempotencyKey) > 200 {
		return fmt.Errorf("%w: idempotency key", ErrInvalid)
	}
	if strings.TrimSpace(c.PurposeCode) == "" || strings.TrimSpace(c.PolicyVersion) == "" || strings.TrimSpace(c.MappingVersion) == "" || c.DecisionAt.IsZero() {
		return fmt.Errorf("%w: policy evidence", ErrInvalid)
	}
	if !hex64Re.MatchString(strings.ToLower(c.EvidenceSHA256)) || len(c.FieldMapping) == 0 {
		return fmt.Errorf("%w: evidence or mapping", ErrInvalid)
	}
	seenTargets := map[string]bool{}
	for source, target := range c.FieldMapping {
		if strings.TrimSpace(source) == "" || (target != "email" && target != "phone" && target != "name") || seenTargets[target] {
			return fmt.Errorf("%w: field mapping", ErrInvalid)
		}
		seenTargets[target] = true
	}
	if !seenTargets["email"] && !seenTargets["phone"] {
		return fmt.Errorf("%w: email or phone mapping required", ErrInvalid)
	}
	return nil
}

type Field struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

// ContactPayload applies only the project's admitted mapping. Unknown fields
// remain in the immutable source candidate and are not guessed into CRM.
func ContactPayload(fields []Field, mapping map[string]string) (json.RawMessage, error) {
	values := map[string]string{}
	for _, field := range fields {
		if target, ok := mapping[field.ID]; ok {
			if strings.TrimSpace(field.Value) == "" || values[target] != "" {
				return nil, fmt.Errorf("%w: %s", ErrInvalid, target)
			}
			values[target] = strings.TrimSpace(field.Value)
		}
	}
	for source, target := range mapping {
		_ = source
		if (target == "email" || target == "phone") && values[target] == "" {
			return nil, fmt.Errorf("%w: %s", ErrFieldAbsent, target)
		}
	}
	return json.Marshal(values)
}

func (c Command) RequestSHA256() (string, error) {
	if err := c.Validate(); err != nil {
		return "", err
	}
	keys := make([]string, 0, len(c.FieldMapping))
	for key := range c.FieldMapping {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	type pair struct{ Source, Target string }
	pairs := make([]pair, 0, len(keys))
	for _, key := range keys {
		pairs = append(pairs, pair{key, c.FieldMapping[key]})
	}
	payload, err := json.Marshal(struct {
		TenantID, Provider, ProviderEventID, LeadID, ConsentID, PurposeCode, PolicyVersion, EvidenceSHA256, DecisionAt, MappingVersion string
		Mapping                                                                                                                        []pair
	}{c.TenantID, c.Provider, c.ProviderEventID, c.LeadID, c.ConsentID, c.PurposeCode, c.PolicyVersion, strings.ToLower(c.EvidenceSHA256), c.DecisionAt.UTC().Format(time.RFC3339Nano), c.MappingVersion, pairs})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

type Receipt struct {
	LeadID   string
	Replayed bool
}

type Store interface {
	Promote(context.Context, Command) (Receipt, error)
}

type Service struct{ Store Store }

func (s Service) Promote(ctx context.Context, command Command) (Receipt, error) {
	if s.Store == nil {
		return Receipt{}, errors.New("leadpromotion: nil store")
	}
	if err := command.Validate(); err != nil {
		return Receipt{}, err
	}
	return s.Store.Promote(ctx, command)
}
````

### FILE: `internal/leadpromotion/promotion_test.go`
```yaml
block_id: "GO-LEAD-CANDIDATE-PROMOTION:internal/leadpromotion/promotion_test.go:v1"
operation: CREATE
provenance: ADAPTED
source: "official replay and provider test-lead conditions; local regression"
license: "LicenseRef-Workspace-Owner"
sha256: "8b9aadc43a95efccb4bdd705b9cfc6c8d8c35baf52e13eff58698a4de79c423b"
variables: []
secrets_allowed: false
```
````go
package leadpromotion

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func validCommand() Command {
	return Command{TenantID: "tenant", Provider: "google_ads", ProviderEventID: "event-1", LeadID: "lead-1", ConsentID: "consent-1", PurposeCode: "sales-contact", PolicyVersion: "policy-7", EvidenceSHA256: strings.Repeat("a", 64), DecisionAt: time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC), MappingVersion: "google-form-v1", FieldMapping: map[string]string{"EMAIL": "email"}, IdempotencyKey: "promote-0001"}
}

func TestCommandRequiresExplicitPolicyAndContactMapping(t *testing.T) {
	c := validCommand()
	c.PolicyVersion = ""
	if !errors.Is(c.Validate(), ErrInvalid) {
		t.Fatal("missing policy accepted")
	}
	c = validCommand()
	c.FieldMapping = map[string]string{"NAME": "name"}
	if !errors.Is(c.Validate(), ErrInvalid) {
		t.Fatal("mapping without contact accepted")
	}
}

func TestContactPayloadUsesOnlyConfiguredFields(t *testing.T) {
	raw, err := ContactPayload([]Field{{ID: "EMAIL", Value: "person@example.test"}, {ID: "UNMAPPED", Value: "secret"}}, map[string]string{"EMAIL": "email"})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != `{"email":"person@example.test"}` || strings.Contains(string(raw), "secret") {
		t.Fatalf("unexpected payload %s", raw)
	}
}

func TestContactPayloadRejectsMissingConfiguredField(t *testing.T) {
	_, err := ContactPayload([]Field{{ID: "PHONE", Value: "+541100000000"}}, map[string]string{"EMAIL": "email"})
	if !errors.Is(err, ErrFieldAbsent) {
		t.Fatalf("expected absent field, got %v", err)
	}
}

type fakeStore struct{ called bool }

func (f *fakeStore) Promote(_ context.Context, c Command) (Receipt, error) {
	f.called = true
	return Receipt{LeadID: c.LeadID}, nil
}

func TestServiceValidatesBeforeStore(t *testing.T) {
	f := &fakeStore{}
	s := Service{Store: f}
	c := validCommand()
	c.EvidenceSHA256 = "bad"
	if _, err := s.Promote(context.Background(), c); !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid, got %v", err)
	}
	if f.called {
		t.Fatal("store called for invalid command")
	}
}
````

### FILE: `internal/platform/postgres/lead_promotion.go`
```yaml
block_id: "GO-LEAD-CANDIDATE-PROMOTION:internal/platform/postgres/lead_promotion.go:v1"
operation: CREATE
provenance: ADAPTED
source: "official PostgreSQL serializable transaction and AWS idempotent-consumer contracts; local schema"
license: "LicenseRef-Workspace-Owner"
sha256: "31001b1c9476e0d5bd2cd4d13cd3edfa5d0f4e832a03fff58dc9996e3624e6f2"
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

	"elite.local/enterprise/internal/leadpromotion"
	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LeadPromotion struct{ pool *pgxpool.Pool }

func NewLeadPromotion(pool *pgxpool.Pool) *LeadPromotion { return &LeadPromotion{pool: pool} }

func (r *LeadPromotion) Promote(ctx context.Context, c leadpromotion.Command) (leadpromotion.Receipt, error) {
	if r == nil || r.pool == nil {
		return leadpromotion.Receipt{}, errors.New("postgres lead promotion: nil pool")
	}
	requestHash, err := c.RequestSHA256()
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at) values($1,'lead-promotion',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict do nothing`, c.TenantID, c.IdempotencyKey, requestHash)
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	if result.RowsAffected() == 0 {
		var priorHash, status, resourceID string
		if err = tx.QueryRow(ctx, `select request_sha256_hex,status,coalesce(resource_id,'') from platform.idempotency_record where tenant_id=$1 and scope='lead-promotion' and idempotency_key=$2`, c.TenantID, c.IdempotencyKey).Scan(&priorHash, &status, &resourceID); err != nil {
			return leadpromotion.Receipt{}, err
		}
		if priorHash != requestHash || status != "completed" || resourceID == "" {
			return leadpromotion.Receipt{}, leadpromotion.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return leadpromotion.Receipt{}, err
		}
		return leadpromotion.Receipt{LeadID: resourceID, Replayed: true}, nil
	}
	var organization string
	var isTest bool
	var fieldsJSON []byte
	var sourceHash string
	err = tx.QueryRow(ctx, `select c.organization_id,c.is_test,c.fields,r.source_payload_sha256 from integration.lead_candidate c join integration.lead_ingress_raw r using(tenant_id,provider,provider_event_id) where c.tenant_id=$1 and c.provider=$2 and c.provider_event_id=$3 and c.contact_eligibility='pending_policy' for update of c`, c.TenantID, c.Provider, c.ProviderEventID).Scan(&organization, &isTest, &fieldsJSON, &sourceHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return leadpromotion.Receipt{}, leadpromotion.ErrNotFound
	}
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	if isTest {
		return leadpromotion.Receipt{}, leadpromotion.ErrTestLead
	}
	var fields []leadpromotion.Field
	if err = json.Unmarshal(fieldsJSON, &fields); err != nil {
		return leadpromotion.Receipt{}, fmt.Errorf("decode candidate fields: %w", err)
	}
	contact, err := leadpromotion.ContactPayload(fields, c.FieldMapping)
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	mappingJSON, err := json.Marshal(c.FieldMapping)
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	_, err = tx.Exec(ctx, `insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required) values($1,$2,$3,'new',$4,$5,true)`, c.TenantID, c.LeadID, organization, c.Provider, contact)
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	_, err = tx.Exec(ctx, `insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex) values($1,$2,$3,$4,$5,'granted',$6,$7)`, c.TenantID, c.ConsentID, c.LeadID, c.PurposeCode, c.PolicyVersion, c.DecisionAt.UTC(), c.EvidenceSHA256)
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	_, err = tx.Exec(ctx, `insert into integration.lead_promotion(tenant_id,provider,provider_event_id,lead_id,consent_id,mapping_version,field_mapping,source_payload_sha256,decision_at) values($1,$2,$3,$4,$5,$6,$7,$8,$9)`, c.TenantID, c.Provider, c.ProviderEventID, c.LeadID, c.ConsentID, c.MappingVersion, mappingJSON, sourceHash, c.DecisionAt.UTC())
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	eventID := leadstream.StableUUID(c.TenantID, c.Provider, c.ProviderEventID, "promoted")
	payload, _ := json.Marshal(map[string]string{"provider": c.Provider, "provider_event_id": c.ProviderEventID, "lead_id": c.LeadID, "organization_id": organization, "purpose_code": c.PurposeCode})
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'lead',$3,1,'lead.promoted',1,clock_timestamp(),$4)`, c.TenantID, eventID, c.LeadID, payload)
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	_, err = tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('lead_id',$3::text),resource_type='lead',resource_id=$3,locked_until=null where tenant_id=$1 and scope='lead-promotion' and idempotency_key=$2 and status='processing'`, c.TenantID, c.IdempotencyKey, c.LeadID)
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return leadpromotion.Receipt{}, err
	}
	return leadpromotion.Receipt{LeadID: c.LeadID}, nil
}
````

### FILE: `internal/platform/postgres/lead_promotion_integration_test.go`
```yaml
block_id: "GO-LEAD-CANDIDATE-PROMOTION:internal/platform/postgres/lead_promotion_integration_test.go:v1"
operation: CREATE
provenance: ADAPTED
source: "official at-least-once/replay conditions; local PostgreSQL integration regression"
license: "LicenseRef-Workspace-Owner"
sha256: "fa60af0cf11f615e402657bbd4a8239971221fd3dac807eb325dd10f580e2cb4"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/leadpromotion"
	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestLeadPromotionRequiresEvidenceAndIsAtomicAndIdempotent(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := leadstream.StableUUID(t.Name(), time.Now().UTC().Format(time.RFC3339Nano))
	organization := "promotion-store"
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Promotion Test','Promotion Test')`, tenant, "promotion-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,$2,$2,'Promotion Store','store')`, tenant, organization); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = pool.Exec(context.Background(), `update platform.outbox_event set published_at=greatest(clock_timestamp(),occurred_at),claimed_by=null,claimed_until=null where tenant_id=$1 and published_at is null`, tenant)
	}()

	repo := NewLeadIngress(pool)
	when := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)
	raw, candidate := promotionCandidate(t, tenant, organization, "provider-event-1", false, when)
	if _, err = repo.Record(ctx, raw, candidate, ""); err != nil {
		t.Fatal(err)
	}
	command := leadpromotion.Command{TenantID: tenant, Provider: "google_ads", ProviderEventID: "provider-event-1", LeadID: "crm-lead-1", ConsentID: "consent-1", PurposeCode: "sales-contact", PolicyVersion: "policy-approved-7", EvidenceSHA256: strings.Repeat("a", 64), DecisionAt: when.Add(time.Minute), MappingVersion: "google-form-42-v3", FieldMapping: map[string]string{"EMAIL": "email", "FULL_NAME": "name"}, IdempotencyKey: "promotion-request-0001"}
	promoter := NewLeadPromotion(pool)
	receipt, err := promoter.Promote(ctx, command)
	if err != nil || receipt.Replayed || receipt.LeadID != "crm-lead-1" {
		t.Fatalf("first promotion=%+v err=%v", receipt, err)
	}
	receipt, err = promoter.Promote(ctx, command)
	if err != nil || !receipt.Replayed {
		t.Fatalf("replay=%+v err=%v", receipt, err)
	}
	changed := command
	changed.PolicyVersion = "different"
	if _, err = promoter.Promote(ctx, changed); !errors.Is(err, leadpromotion.ErrConflict) {
		t.Fatalf("expected conflicting replay, got %v", err)
	}
	var leads, consents, promotions, events int
	if err = pool.QueryRow(ctx, `select count(*) from crm.lead where tenant_id=$1 and lead_id='crm-lead-1' and contact_payload->>'email'='person@example.test'`, tenant).Scan(&leads); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from crm.consent_evidence where tenant_id=$1 and consent_id='consent-1' and decision='granted' and policy_version='policy-approved-7'`, tenant).Scan(&consents); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from integration.lead_promotion where tenant_id=$1 and provider_event_id='provider-event-1' and source_payload_sha256=$2`, tenant, raw.SourceSHA256).Scan(&promotions); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='lead.promoted' and aggregate_id='crm-lead-1'`, tenant).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if leads != 1 || consents != 1 || promotions != 1 || events != 1 {
		t.Fatalf("leads=%d consents=%d promotions=%d events=%d", leads, consents, promotions, events)
	}

	testRaw, testCandidate := promotionCandidate(t, tenant, organization, "provider-event-test", true, when)
	if _, err = repo.Record(ctx, testRaw, testCandidate, ""); err != nil {
		t.Fatal(err)
	}
	testCommand := command
	testCommand.ProviderEventID = "provider-event-test"
	testCommand.LeadID = "must-not-exist"
	testCommand.ConsentID = "must-not-exist"
	testCommand.IdempotencyKey = "promotion-request-test"
	if _, err = promoter.Promote(ctx, testCommand); !errors.Is(err, leadpromotion.ErrTestLead) {
		t.Fatalf("test lead promotion=%v", err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from crm.lead where tenant_id=$1 and lead_id='must-not-exist'`, tenant).Scan(&leads); err != nil || leads != 0 {
		t.Fatalf("test lead persisted=%d err=%v", leads, err)
	}
}

func promotionCandidate(t *testing.T, tenant, organization, eventID string, isTest bool, when time.Time) (leadstream.RawEvent, *leadstream.LeadCandidate) {
	t.Helper()
	payload := []byte(`{"lead_id":"` + eventID + `"}`)
	raw, err := leadstream.NewRawEvent(tenant, organization, "google_ads", eventID, "google.ads.lead.v3", "https://googleads.googleapis.com/lead-form", "3", when, when.Add(time.Second), payload)
	if err != nil {
		t.Fatal(err)
	}
	return raw, &leadstream.LeadCandidate{TenantID: tenant, OrganizationID: organization, Provider: "google_ads", ProviderLeadID: eventID, FormID: "42", SubmittedAt: when, IsTest: isTest, Fields: []leadstream.Field{{ID: "EMAIL", Value: "person@example.test"}, {ID: "FULL_NAME", Value: "Ada Lovelace"}}, ContactEligibility: "pending_policy"}
}
````

### FILE: `db/migrations/0045_lead_candidate_promotion.up.sql`
```yaml
block_id: "GO-LEAD-CANDIDATE-PROMOTION:db/migrations/0045_lead_candidate_promotion.up.sql:v1"
operation: CREATE
provenance: ADAPTED
source: "official PostgreSQL constraints/FK/transaction contracts; local schema"
license: "LicenseRef-Workspace-Owner"
sha256: "2d6b404be0e8611bbc7951b8759ac70ecd326390fb53b496047dae3782ea139d"
variables: []
secrets_allowed: false
```
````sql
begin;

create table integration.lead_promotion (
  tenant_id uuid not null,
  provider text not null,
  provider_event_id text not null,
  lead_id text not null,
  consent_id text not null,
  mapping_version text not null,
  field_mapping jsonb not null,
  source_payload_sha256 text not null check (source_payload_sha256 ~ '^[0-9a-f]{64}$'),
  decision_at timestamptz not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, provider, provider_event_id),
  unique (tenant_id, lead_id),
  unique (tenant_id, consent_id),
  foreign key (tenant_id, provider, provider_event_id) references integration.lead_candidate (tenant_id, provider, provider_event_id),
  foreign key (tenant_id, lead_id) references crm.lead (tenant_id, lead_id),
  foreign key (tenant_id, consent_id) references crm.consent_evidence (tenant_id, consent_id),
  check (length(mapping_version) between 1 and 128),
  check (jsonb_typeof(field_mapping)='object')
);

create trigger lead_promotion_immutable before update or delete on integration.lead_promotion
for each row execute function integration.reject_lead_ingress_mutation();

commit;
````

### FILE: `db/migrations/0045_lead_candidate_promotion.down.sql`
```yaml
block_id: "GO-LEAD-CANDIDATE-PROMOTION:db/migrations/0045_lead_candidate_promotion.down.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local rollback"
license: "LicenseRef-Workspace-Owner"
sha256: "9988d288b62cc7e09b2e755790f3e4719e35dad4816167b840b00537e6b16c9f"
variables: []
secrets_allowed: false
```
````sql
begin;
drop table integration.lead_promotion;
commit;
````

### FILE: `db/tests/0045_lead_candidate_promotion.test.sql`
```yaml
block_id: "GO-LEAD-CANDIDATE-PROMOTION:db/tests/0045_lead_candidate_promotion.test.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local schema regression"
license: "LicenseRef-Workspace-Owner"
sha256: "02b6e88cab6ccf47431c9412b8950f834341229cddb0e4f3a56199245e56c498"
variables: []
secrets_allowed: false
```
````sql
begin;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values('45454545-4545-4545-8545-454545454545','lead-promotion-test','Lead Promotion Test','Lead Promotion Test');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
values('45454545-4545-4545-8545-454545454545','store-1','store-1','Store 1','store');
insert into integration.lead_ingress_raw
(tenant_id,organization_id,provider,provider_event_id,event_type,source,schema_version,occurred_at,received_at,payload_redacted,source_payload_sha256,stored_payload_sha256,redaction_profile,state)
values('45454545-4545-4545-8545-454545454545','store-1','google_ads','lead-1','google.ads.lead.v3','https://googleads.googleapis.com/lead-form','3',clock_timestamp(),clock_timestamp(),convert_to('{"lead_id":"lead-1"}','UTF8'),repeat('a',64),repeat('b',64),'google_ads:remove-google_key:v1','normalized');
insert into integration.lead_candidate
(tenant_id,provider,provider_event_id,organization_id,provider_lead_id,submitted_at,is_test,fields,contact_eligibility)
values('45454545-4545-4545-8545-454545454545','google_ads','lead-1','store-1','lead-1',clock_timestamp(),false,'[{"id":"EMAIL","value":"person@example.test"}]','pending_policy');
insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required)
values('45454545-4545-4545-8545-454545454545','crm-lead-1','store-1','new','google_ads','{"email":"person@example.test"}',true);
insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex)
values('45454545-4545-4545-8545-454545454545','consent-1','crm-lead-1','sales-contact','policy-7','granted',clock_timestamp(),repeat('c',64));
insert into integration.lead_promotion
(tenant_id,provider,provider_event_id,lead_id,consent_id,mapping_version,field_mapping,source_payload_sha256,decision_at)
values('45454545-4545-4545-8545-454545454545','google_ads','lead-1','crm-lead-1','consent-1','google-form-v1','{"EMAIL":"email"}',repeat('a',64),clock_timestamp());

do $test$
begin
  begin
    update integration.lead_promotion set mapping_version='changed'
    where tenant_id='45454545-4545-4545-8545-454545454545' and provider='google_ads' and provider_event_id='lead-1';
    raise exception 'promotion mutation unexpectedly succeeded';
  exception when sqlstate '55000' then null;
  end;
  begin
    delete from integration.lead_promotion
    where tenant_id='45454545-4545-4545-8545-454545454545' and provider='google_ads' and provider_event_id='lead-1';
    raise exception 'promotion delete unexpectedly succeeded';
  exception when sqlstate '55000' then null;
  end;
end;
$test$;

rollback;
````

## 6. Configuration surface

| Entrada | Obligatoria | Origen |
|---|---:|---|
| tenant/provider/event/lead/consent IDs | sí | identidad durable del proyecto |
| purpose/policy/evidence/decision time | sí | decisión humana o política aprobada |
| mapping versionado de campos | sí | formulario/cuenta/corpus admitidos |
| idempotency key | sí | worker/orquestador |

## 7. Dependency bill

| Dependencia | Pin | Uso | Licencia |
|---|---|---|---|
| Go | 1.26.7 | dominio/tests | BSD-3-Clause |
| PostgreSQL | 18.6 | transacción/constraints | PostgreSQL |
| pgx/v5 | lock del backend compuesto | adapter | MIT |

## 8. Apply order

1. Materializar foundation, CRM/journeys e ingesta 0044.
2. Aplicar `0045_lead_candidate_promotion.up.sql`.
3. Inyectar una decisión y mapping aprobados; no crearlos por defecto.
4. Ejecutar pruebas de dominio, SQL, PostgreSQL, replay y rollback.

## 9. Verification

- dominio: validación de política/mapping, exclusión de campos y store fail-closed;
- PostgreSQL real: promoción + consentimiento + procedencia + outbox atómicos, replay exacto, divergencia y test lead bloqueado;
- SQL: fila de promoción inmutable;
- pendiente por proyecto: evidencia jurídica/comercial, mapping/corpus real, RLS/cifrado/retención y carga.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_LEAD_CANDIDATE_PROMOTION_2026-09-04_V225.md`.
