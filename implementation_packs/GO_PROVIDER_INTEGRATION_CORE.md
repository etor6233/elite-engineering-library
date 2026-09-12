# Go Provider Integration Core

## 1. Metadata

```yaml
pack_id: "GO-PROVIDER-INTEGRATION-CORE"
pack_version: "0.1.3"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un borde reusable para webhooks firmados, secretos referenciados, deduplicación transaccional y jobs de reconciliación/procesamiento por proveedor."
stacks: ["Go 1.26.7", "PostgreSQL 18.6"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "PG-TX-FOUNDATION 0.1.x", "GO-RELIABLE-ASYNC-WORKERS 0.1.x"]
incompatible_with: ["provider signature protocol without a matching admitted verifier"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0 AND MIT dependencies"
upstream_sources: ["https://go.dev", "https://www.postgresql.org/docs/18/"]
verified_at: "2026-09-06"
```

Este núcleo no afirma que todos los proveedores firman igual. `HMACSHA256` es un adapter de referencia verificable; Stripe, Mercado Pago, Mercado Libre, Amazon, Google o Meta deben aportar su adapter oficial y sus contract tests. El tenant se resuelve desde una conexión preconfigurada, nunca desde el cuerpo o un header del webhook. Los secretos no pueden vivir en el JSON de conexiones: sólo se permiten nombres de variables entregadas por el secret manager del target.

## 2. Applicability

Use as the common inbound boundary for providers that support an admitted authentication/replay contract and durable asynchronous processing. Reject the reference HMAC verifier when the provider protocol differs; implement and contract-test the official scheme instead. This core does not claim universal compatibility with Stripe, marketplaces or ad platforms.

## 3. Architecture contract

Connection identity is preconfigured and maps to one tenant; request body/header data never chooses the tenant. Exact bytes are authenticated before parsing. Receipt, hash-based replay/conflict handling, inbox event and durable job commit atomically. Secrets are references resolved outside source control. Provider-specific handlers are idempotent, effect remote calls outside the receipt transaction and reconcile periodically.

V273 / 0.1.2: AcceptWebhook rechecks the exact active tenant/connection/provider
row under FOR SHARE through commit. Disabled or missing connections cannot accept
new events or acknowledge replays, even if a static registry still contains them.
A replay also requires the same provider, event type, body hash and JSONB payload;
JSON whitespace/key order is immaterial, changed semantics are a conflict. The
existing inbox and job transaction remain the only durable admission mechanism.
This correction is AUTHORED, not Meta or Microsoft source. PostgreSQL 18
Explicit Locking governs lock semantics; Microsoft Azure API design supports
consistent idempotent request handling, not certification of this implementation.
Run DATABASE_URL against a dedicated loopback elite_provider_ database with the
profile migrations, then go test ./internal/platform/postgres -run '^TestProviderWebhook' -v -count=1.
V273 evidence: reconstruction_evidence/PROVIDER_INBOX_ADMISSION_SCOPE_V273.md.
Disabling a connection is serialized with admission, not retroactive cancellation
of already committed jobs. Each worker must independently recheck current authority
before its own effects. This is not the Meta signature protocol or a mounted
WhatsApp callback; provider-specific verification, secure raw-byte retention,
worker routing, acknowledgement semantics and live delivery remain separate gates.

## 4. Exact file manifest

```text
CREATE db/migrations/0004_provider_integration.up.sql
CREATE db/migrations/0004_provider_integration.down.sql
CREATE internal/providerintegration/service.go
CREATE internal/providerintegration/config.go
CREATE internal/providerintegration/service_test.go
CREATE internal/providerintegration/reconciliation.go
CREATE internal/providerintegration/reconciliation_test.go
CREATE internal/platform/postgres/providerintegration.go
CREATE internal/platform/postgres/providerintegration_integration_test.go
CREATE internal/platform/httpapi/providerintegration.go
CREATE internal/platform/httpapi/providerintegration_test.go
```

## 5. Materialization blocks

### FILE: `db/migrations/0004_provider_integration.up.sql`
```yaml
block_id: "GO-PROVIDER-INTEGRATION:migration-up:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "7477ca14e4b4ddff9b53479cd171bcb8d8852bd39937d324b5258a84680d742c"
variables: []
secrets_allowed: false
```
````sql
begin;

create table integration.provider_connection (
  tenant_id uuid not null,
  connection_id text not null,
  provider_code text not null,
  organization_id text,
  secret_ref text not null,
  state text not null default 'active' check (state in ('active', 'disabled')),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, connection_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  unique (tenant_id, connection_id, provider_code),
  check (connection_id ~ '^[a-z0-9][a-z0-9._-]{0,127}$'),
  check (provider_code ~ '^[a-z0-9][a-z0-9._-]{0,63}$'),
  check (secret_ref ~ '^[A-Z][A-Z0-9_]{2,127}$'),
  check (updated_at >= created_at)
);

create table integration.webhook_event (
  tenant_id uuid not null,
  connection_id text not null,
  provider_code text not null,
  provider_event_id text not null,
  event_type text not null,
  body_sha256_hex text not null,
  payload jsonb not null,
  received_at timestamptz not null default clock_timestamp(),
  processed_at timestamptz,
  state text not null default 'received' check (state in ('received', 'processed', 'failed')),
  last_error_code text,
  primary key (tenant_id, connection_id, provider_event_id),
  foreign key (tenant_id, connection_id, provider_code)
    references integration.provider_connection (tenant_id, connection_id, provider_code),
  check (provider_event_id <> '' and length(provider_event_id) <= 200),
  check (event_type <> '' and length(event_type) <= 200),
  check (body_sha256_hex ~ '^[0-9a-f]{64}$'),
  check (jsonb_typeof(payload) = 'object'),
  check ((state = 'processed') = (processed_at is not null))
);

create index provider_webhook_ready_idx
  on integration.webhook_event (tenant_id, provider_code, received_at, provider_event_id)
  where state = 'received';

commit;
````

### FILE: `db/migrations/0004_provider_integration.down.sql`
```yaml
block_id: "GO-PROVIDER-INTEGRATION:migration-down:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "114dd0cebdf86587db757079a72edde0463a159ab2741c35b754f084ce8fe187"
variables: []
secrets_allowed: false
```
````sql
begin;
drop table integration.webhook_event;
drop table integration.provider_connection;
commit;
````

### FILE: `internal/providerintegration/service.go`
```yaml
block_id: "GO-PROVIDER-INTEGRATION:service:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "2639ceb1e734d8615886c714051a9a8049c6dc5d49b072505a7571d856825be1"
variables: []
secrets_allowed: false
```
````go
package providerintegration

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalid         = errors.New("invalid webhook")
	ErrUnauthenticated = errors.New("webhook authentication failed")
	ErrConnection      = errors.New("provider connection not found")
	ErrConflict        = errors.New("provider event conflict")
)

type Connection struct {
	TenantID, ConnectionID, ProviderCode, SecretRef string
	Secret                                          []byte
}

type Registry interface {
	Lookup(providerCode, connectionID string) (Connection, bool)
}

type Receipt struct {
	TenantID, ConnectionID, ProviderCode string
	ProviderEventID, EventType, BodyHash string
	Payload                              json.RawMessage
}

type Repository interface {
	AcceptWebhook(context.Context, Receipt) (replayed bool, err error)
}

type SignatureVerifier interface {
	Verify(timestamp, signature string, body, secret []byte) error
}

type Service struct {
	repository Repository
	registry   Registry
	verifier   SignatureVerifier
}

func NewService(repository Repository, registry Registry, verifier SignatureVerifier) *Service {
	return &Service{repository: repository, registry: registry, verifier: verifier}
}

func (s *Service) Receive(ctx context.Context, providerCode, connectionID, timestamp, signature string, body []byte) (bool, error) {
	if s == nil || s.repository == nil || s.registry == nil || s.verifier == nil || len(body) == 0 {
		return false, ErrInvalid
	}
	connection, ok := s.registry.Lookup(providerCode, connectionID)
	if !ok {
		return false, ErrConnection
	}
	if err := s.verifier.Verify(timestamp, signature, body, connection.Secret); err != nil {
		return false, err
	}
	var envelope struct {
		EventID string          `json:"event_id"`
		Type    string          `json:"type"`
		Data    json.RawMessage `json:"data"`
	}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&envelope); err != nil {
		return false, ErrInvalid
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) || !bounded(envelope.EventID, 200) || !bounded(envelope.Type, 200) || !jsonObject(envelope.Data) {
		return false, ErrInvalid
	}
	hash := sha256.Sum256(body)
	return s.repository.AcceptWebhook(ctx, Receipt{TenantID: connection.TenantID, ConnectionID: connection.ConnectionID, ProviderCode: connection.ProviderCode, ProviderEventID: envelope.EventID, EventType: envelope.Type, BodyHash: hex.EncodeToString(hash[:]), Payload: envelope.Data})
}

func bounded(value string, limit int) bool { return value != "" && len(value) <= limit }
func jsonObject(value json.RawMessage) bool {
	var object map[string]json.RawMessage
	return len(value) > 0 && json.Unmarshal(value, &object) == nil && object != nil
}

type HMACSHA256 struct {
	Now       func() time.Time
	Tolerance time.Duration
}

func (v HMACSHA256) Verify(timestamp, signature string, body, secret []byte) error {
	if len(secret) < 32 || v.Now == nil || v.Tolerance <= 0 || v.Tolerance > 15*time.Minute {
		return ErrUnauthenticated
	}
	seconds, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return ErrUnauthenticated
	}
	delta := v.Now().Sub(time.Unix(seconds, 0))
	if delta < -v.Tolerance || delta > v.Tolerance {
		return ErrUnauthenticated
	}
	if !strings.HasPrefix(signature, "v1=") || strings.Contains(signature[3:], ",") {
		return ErrUnauthenticated
	}
	provided, err := hex.DecodeString(signature[3:])
	if err != nil || len(provided) != sha256.Size {
		return ErrUnauthenticated
	}
	mac := hmac.New(sha256.New, secret)
	_, _ = fmt.Fprintf(mac, "%s.", timestamp)
	_, _ = mac.Write(body)
	if !hmac.Equal(provided, mac.Sum(nil)) {
		return ErrUnauthenticated
	}
	return nil
}
````

### FILE: `internal/providerintegration/config.go`
```yaml
block_id: "GO-PROVIDER-INTEGRATION:config:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "beb3b892d5843fb95a044dd6b13a7cb04294670cae09f6a135addae521060895"
variables: []
secrets_allowed: false
```
````go
package providerintegration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
)

type StaticRegistry struct{ connections map[string]Connection }

func (r *StaticRegistry) Lookup(providerCode, connectionID string) (Connection, bool) {
	if r == nil {
		return Connection{}, false
	}
	value, ok := r.connections[providerCode+"\x00"+connectionID]
	return value, ok
}

var (
	codePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,127}$`)
	envPattern  = regexp.MustCompile(`^[A-Z][A-Z0-9_]{2,127}$`)
)

func LoadHMACConnections(configuration []byte, lookupEnv func(string) (string, bool)) (*StaticRegistry, error) {
	registry := &StaticRegistry{connections: map[string]Connection{}}
	if len(bytes.TrimSpace(configuration)) == 0 {
		return registry, nil
	}
	if len(configuration) > 64<<10 || lookupEnv == nil {
		return nil, fmt.Errorf("invalid provider connection configuration")
	}
	var entries []struct {
		TenantID     string `json:"tenant_id"`
		ConnectionID string `json:"connection_id"`
		ProviderCode string `json:"provider_code"`
		SecretEnv    string `json:"secret_env"`
	}
	decoder := json.NewDecoder(bytes.NewReader(configuration))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&entries); err != nil || len(entries) > 100 {
		return nil, fmt.Errorf("invalid provider connection configuration")
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return nil, fmt.Errorf("invalid provider connection configuration")
	}
	for _, entry := range entries {
		if entry.TenantID == "" || !codePattern.MatchString(entry.ConnectionID) || len(entry.ProviderCode) > 64 || !codePattern.MatchString(entry.ProviderCode) || !envPattern.MatchString(entry.SecretEnv) {
			return nil, fmt.Errorf("invalid provider connection metadata")
		}
		secret, ok := lookupEnv(entry.SecretEnv)
		if !ok || len(secret) < 32 {
			return nil, fmt.Errorf("provider secret is missing or too short")
		}
		key := entry.ProviderCode + "\x00" + entry.ConnectionID
		if _, duplicate := registry.connections[key]; duplicate {
			return nil, fmt.Errorf("duplicate provider connection")
		}
		registry.connections[key] = Connection{TenantID: entry.TenantID, ConnectionID: entry.ConnectionID, ProviderCode: entry.ProviderCode, SecretRef: entry.SecretEnv, Secret: []byte(secret)}
	}
	return registry, nil
}
````

### FILE: `internal/providerintegration/service_test.go`
```yaml
block_id: "GO-PROVIDER-INTEGRATION:service-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "5e585e3a6e3dbed94f3360559ae27b4e652b642d51ff6ea294f4977c168cdf06"
variables: []
secrets_allowed: false
```
````go
package providerintegration

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"
)

type receiptRepository struct {
	receipt Receipt
	replay  bool
}

func (r *receiptRepository) AcceptWebhook(_ context.Context, value Receipt) (bool, error) {
	r.receipt = value
	return r.replay, nil
}

func signed(secret []byte, timestamp string, body []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(timestamp + "."))
	_, _ = mac.Write(body)
	return "v1=" + hex.EncodeToString(mac.Sum(nil))
}

func TestReceiveAuthenticatesAndNormalizes(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	secret := []byte("01234567890123456789012345678901")
	registry, err := LoadHMACConnections([]byte(`[{"tenant_id":"tenant","connection_id":"primary","provider_code":"sandbox","secret_env":"SANDBOX_SECRET"}]`), func(string) (string, bool) { return string(secret), true })
	if err != nil {
		t.Fatal(err)
	}
	repository := &receiptRepository{}
	service := NewService(repository, registry, HMACSHA256{Now: func() time.Time { return now }, Tolerance: 5 * time.Minute})
	body := []byte(`{"event_id":"evt-1","type":"payment.updated","data":{"state":"paid"}}`)
	replayed, err := service.Receive(context.Background(), "sandbox", "primary", "1800000000", signed(secret, "1800000000", body), body)
	if err != nil || replayed || repository.receipt.ProviderEventID != "evt-1" || repository.receipt.TenantID != "tenant" || len(repository.receipt.BodyHash) != 64 {
		t.Fatalf("replayed=%v receipt=%+v err=%v", replayed, repository.receipt, err)
	}
}

func TestReceiveRejectsTamperReplayWindowAndShape(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	secret := []byte("01234567890123456789012345678901")
	registry := &StaticRegistry{connections: map[string]Connection{"sandbox\x00primary": {TenantID: "tenant", ConnectionID: "primary", ProviderCode: "sandbox", Secret: secret}}}
	service := NewService(&receiptRepository{}, registry, HMACSHA256{Now: func() time.Time { return now }, Tolerance: 5 * time.Minute})
	body := []byte(`{"event_id":"evt-1","type":"payment.updated","data":{}}`)
	if _, err := service.Receive(context.Background(), "sandbox", "primary", "1799999000", signed(secret, "1799999000", body), body); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("old timestamp err=%v", err)
	}
	if _, err := service.Receive(context.Background(), "sandbox", "primary", "1800000000", signed(secret, "1800000000", body), append(body, ' ')); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("tamper err=%v", err)
	}
	invalid := []byte(`{"event_id":"evt-1","type":"payment.updated","data":[],"extra":1}`)
	if _, err := service.Receive(context.Background(), "sandbox", "primary", "1800000000", signed(secret, "1800000000", invalid), invalid); !errors.Is(err, ErrInvalid) {
		t.Fatalf("shape err=%v", err)
	}
}

func TestConnectionConfigurationRejectsInlineOrMissingSecrets(t *testing.T) {
	lookup := func(string) (string, bool) { return "", false }
	if _, err := LoadHMACConnections([]byte(`[{"tenant_id":"t","connection_id":"c","provider_code":"p","secret":"inline"}]`), lookup); err == nil {
		t.Fatal("inline secret accepted")
	}
	if _, err := LoadHMACConnections([]byte(`[{"tenant_id":"t","connection_id":"c","provider_code":"p","secret_env":"PROVIDER_SECRET"}]`), lookup); err == nil {
		t.Fatal("missing secret accepted")
	}
}
````

### FILE: `internal/platform/postgres/providerintegration.go`
```yaml
block_id: "GO-PROVIDER-INTEGRATION:postgres:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "77fd3573d0cf246ec068af19757a8cb58c148d977e874088f31681cb104b8dbd"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"errors"
	"fmt"

	"elite.local/enterprise/internal/providerintegration"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProviderIntegration struct{ pool *pgxpool.Pool }

func NewProviderIntegration(pool *pgxpool.Pool) *ProviderIntegration {
	return &ProviderIntegration{pool: pool}
}

func (r *ProviderIntegration) AcceptWebhook(ctx context.Context, value providerintegration.Receipt) (bool, error) {
	return r.acceptWebhookWithEffect(ctx, value, "provider-events", nil)
}

// The optional local effect commits with a newly inserted inbox entry/job.
// Replays never repeat the effect; failure rolls back all three components.
func (r *ProviderIntegration) acceptWebhookWithEffect(ctx context.Context, value providerintegration.Receipt, queue string, effect func(context.Context, pgx.Tx) error) (bool, error) {
	if r == nil || r.pool == nil {
		return false, providerintegration.ErrConnection
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	// A startup registry is not durable authority. Lock the current connection
	// through the event/job commit; a concurrent disable must serialize with it.
	var connection string
	err = tx.QueryRow(ctx, `select connection_id from integration.provider_connection
where tenant_id=$1 and connection_id=$2 and provider_code=$3 and state='active'
for share`, value.TenantID, value.ConnectionID, value.ProviderCode).Scan(&connection)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, providerintegration.ErrConnection
	}
	if err != nil {
		return false, fmt.Errorf("resolve provider connection: %w", err)
	}
	var inserted bool
	err = tx.QueryRow(ctx, `with accepted as (
insert into integration.webhook_event(tenant_id,connection_id,provider_code,provider_event_id,event_type,body_sha256_hex,payload)
values($1,$2,$3,$4,$5,$6,$7) on conflict do nothing returning true)
select coalesce((select true from accepted),false)`, value.TenantID, value.ConnectionID, value.ProviderCode, value.ProviderEventID, value.EventType, value.BodyHash, value.Payload).Scan(&inserted)
	if err != nil {
		return false, fmt.Errorf("accept provider webhook: %w", err)
	}
	if !inserted {
		var equivalent bool
		if err = tx.QueryRow(ctx, `select provider_code=$4 and event_type=$5 and body_sha256_hex=$6 and payload=$7::jsonb
from integration.webhook_event where tenant_id=$1 and connection_id=$2 and provider_event_id=$3`, value.TenantID, value.ConnectionID, value.ProviderEventID, value.ProviderCode, value.EventType, value.BodyHash, value.Payload).Scan(&equivalent); err != nil {
			return false, err
		}
		if !equivalent {
			return false, providerintegration.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return false, err
		}
		return true, nil
	}
	if effect != nil {
		if err = effect(ctx, tx); err != nil {
			return false, err
		}
	}
	_, err = tx.Exec(ctx, `insert into platform.job(tenant_id,job_id,queue,job_type,schema_version,payload,max_attempts)
values($1,gen_random_uuid(),$6,'provider.webhook.received',1,jsonb_build_object('connection_id',$2::text,'provider_code',$3::text,'provider_event_id',$4::text,'event_type',$5::text),10)`, value.TenantID, value.ConnectionID, value.ProviderCode, value.ProviderEventID, value.EventType, queue)
	if err != nil {
		return false, fmt.Errorf("enqueue provider webhook: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return false, err
	}
	return false, nil
}

func (r *ProviderIntegration) UpsertMapping(ctx context.Context, value providerintegration.Mapping) error {
	_, err := r.pool.Exec(ctx, `insert into integration.external_mapping(tenant_id,provider_code,resource_type,internal_id,external_id,version_token)
values($1,$2,$3,$4,$5,nullif($6,''))
on conflict(tenant_id,provider_code,resource_type,internal_id) do update
set external_id=excluded.external_id,version_token=excluded.version_token,updated_at=clock_timestamp()`, value.TenantID, value.ProviderCode, value.ResourceType, value.InternalID, value.ExternalID, value.VersionToken)
	return err
}

func (r *ProviderIntegration) RecordReconciliation(ctx context.Context, value providerintegration.Reconciliation) error {
	_, err := r.pool.Exec(ctx, `insert into integration.reconciliation_item(tenant_id,reconciliation_id,provider_code,resource_type,internal_id,external_id,state,evidence,observed_at)
values($1,$2,$3,$4,nullif($5,''),nullif($6,''),$7,$8,$9)`, value.TenantID, value.ID, value.ProviderCode, value.ResourceType, value.InternalID, value.ExternalID, value.State, value.Evidence, value.ObservedAt)
	return err
}

func (r *ProviderIntegration) ResolveReconciliation(ctx context.Context, tenantID, providerCode, reconciliationID string) error {
	result, err := r.pool.Exec(ctx, `update integration.reconciliation_item set state='resolved',resolved_at=clock_timestamp() where tenant_id=$1 and provider_code=$2 and reconciliation_id=$3 and state<>'resolved'`, tenantID, providerCode, reconciliationID)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return providerintegration.ErrConnection
	}
	return nil
}
````

### FILE: `internal/providerintegration/reconciliation.go`
```yaml
block_id: "GO-PROVIDER-INTEGRATION:reconciliation:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "77f7edb91a51bcbb9f48331303a2e5972eba65d7bacd6270252e3d172166a93d"
variables: []
secrets_allowed: false
```
````go
package providerintegration

import (
	"context"
	"encoding/json"
	"time"
)

type Mapping struct {
	TenantID, ProviderCode, ResourceType string
	InternalID, ExternalID, VersionToken string
}

type Reconciliation struct {
	TenantID, ID, ProviderCode, ResourceType string
	InternalID, ExternalID, State            string
	Evidence                                 json.RawMessage
	ObservedAt                               time.Time
}

type ReconciliationRepository interface {
	UpsertMapping(context.Context, Mapping) error
	RecordReconciliation(context.Context, Reconciliation) error
	ResolveReconciliation(context.Context, string, string, string) error
}

type ReconciliationService struct{ repository ReconciliationRepository }

func NewReconciliationService(repository ReconciliationRepository) *ReconciliationService {
	return &ReconciliationService{repository: repository}
}

func (s *ReconciliationService) UpsertMapping(ctx context.Context, value Mapping) error {
	if s == nil || s.repository == nil || !bounded(value.TenantID, 200) || !bounded(value.ProviderCode, 64) || !bounded(value.ResourceType, 100) || !bounded(value.InternalID, 200) || !bounded(value.ExternalID, 200) || len(value.VersionToken) > 500 {
		return ErrInvalid
	}
	return s.repository.UpsertMapping(ctx, value)
}

func (s *ReconciliationService) Record(ctx context.Context, value Reconciliation) error {
	allowedState := map[string]bool{"matched": true, "missing-internal": true, "missing-external": true, "amount-mismatch": true, "state-mismatch": true}
	if s == nil || s.repository == nil || !bounded(value.TenantID, 200) || !bounded(value.ID, 200) || !bounded(value.ProviderCode, 64) || !bounded(value.ResourceType, 100) || (!bounded(value.InternalID, 200) && !bounded(value.ExternalID, 200)) || !allowedState[value.State] || !jsonObject(value.Evidence) || value.ObservedAt.IsZero() || value.ObservedAt.After(time.Now().Add(24*time.Hour)) {
		return ErrInvalid
	}
	return s.repository.RecordReconciliation(ctx, value)
}

func (s *ReconciliationService) Resolve(ctx context.Context, tenantID, providerCode, reconciliationID string) error {
	if s == nil || s.repository == nil || !bounded(tenantID, 200) || !bounded(providerCode, 64) || !bounded(reconciliationID, 200) {
		return ErrInvalid
	}
	return s.repository.ResolveReconciliation(ctx, tenantID, providerCode, reconciliationID)
}
````

### FILE: `internal/providerintegration/reconciliation_test.go`
```yaml
block_id: "GO-PROVIDER-INTEGRATION:reconciliation-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "8acd034d013a6c2e2e89b00cf5286c525833d5098d29009c383b51104288cbb7"
variables: []
secrets_allowed: false
```
````go
package providerintegration

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

type reconciliationRepository struct{ mappings, records, resolutions int }

func (r *reconciliationRepository) UpsertMapping(context.Context, Mapping) error {
	r.mappings++
	return nil
}
func (r *reconciliationRepository) RecordReconciliation(context.Context, Reconciliation) error {
	r.records++
	return nil
}
func (r *reconciliationRepository) ResolveReconciliation(context.Context, string, string, string) error {
	r.resolutions++
	return nil
}

func TestReconciliationContracts(t *testing.T) {
	repository := &reconciliationRepository{}
	service := NewReconciliationService(repository)
	if err := service.UpsertMapping(context.Background(), Mapping{TenantID: "tenant", ProviderCode: "provider", ResourceType: "offer", InternalID: "model-1", ExternalID: "external-1"}); err != nil {
		t.Fatal(err)
	}
	if err := service.Record(context.Background(), Reconciliation{TenantID: "tenant", ID: "recon-1", ProviderCode: "provider", ResourceType: "offer", InternalID: "model-1", State: "state-mismatch", Evidence: json.RawMessage(`{"internal":"active","external":"paused"}`), ObservedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if err := service.Resolve(context.Background(), "tenant", "provider", "recon-1"); err != nil {
		t.Fatal(err)
	}
	if repository.mappings != 1 || repository.records != 1 || repository.resolutions != 1 {
		t.Fatalf("repository=%+v", repository)
	}
	if err := service.Record(context.Background(), Reconciliation{TenantID: "tenant", ID: "bad", ProviderCode: "provider", ResourceType: "offer", State: "invented", Evidence: json.RawMessage(`{}`), ObservedAt: time.Now()}); err == nil {
		t.Fatal("invalid reconciliation accepted")
	}
}
````

### FILE: `internal/platform/postgres/providerintegration_integration_test.go`
```yaml
block_id: "GO-PROVIDER-INTEGRATION:postgres-integration-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "cb81f9228899be5b7f83b7d24c05b09a07dbf3a9562ae7ec252818a48d28d997"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/providerintegration"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestProviderWebhookConnectionAndReplayScope(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	if !strings.HasPrefix(url, "postgres://postgres@127.0.0.1:55959/elite_provider_") {
		t.Fatal("requires dedicated loopback elite_provider_ test database")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28d41"
	defer func() {
		for _, q := range []string{
			`delete from platform.job where tenant_id=$1`,
			`delete from integration.webhook_event where tenant_id=$1`,
			`delete from integration.provider_connection where tenant_id=$1`,
			`delete from platform.tenant where tenant_id=$1`,
		} {
			if _, err := pool.Exec(ctx, q, tenant); err != nil {
				t.Error("fixture cleanup", err)
			}
		}
	}()
	if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'provider-scope','Synthetic Scope','Synthetic Scope')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref) values($1,'primary','sandbox','SANDBOX_SECRET')`, tenant); err != nil {
		t.Fatal(err)
	}
	repo := NewProviderIntegration(pool)
	base := providerintegration.Receipt{TenantID: tenant, ConnectionID: "primary", ProviderCode: "sandbox", ProviderEventID: "evt-scope", EventType: "payment.updated", BodyHash: strings.Repeat("a", 64), Payload: json.RawMessage(`{"state":"paid"}`)}
	if replay, err := repo.AcceptWebhook(ctx, base); replay || err != nil {
		t.Fatal(replay, err)
	}
	for _, mode := range []string{"disabled_new", "disabled_replay", "event_type", "payload", "provider", "connection"} {
		t.Run(mode, func(t *testing.T) {
			if _, err := pool.Exec(ctx, `update integration.provider_connection set state='active' where tenant_id=$1`, tenant); err != nil {
				t.Fatal(err)
			}
			value := base
			want := providerintegration.ErrConflict
			switch mode {
			case "disabled_new", "disabled_replay":
				if _, err := pool.Exec(ctx, `update integration.provider_connection set state='disabled' where tenant_id=$1`, tenant); err != nil {
					t.Fatal(err)
				}
				if mode == "disabled_new" {
					value.ProviderEventID = "evt-new"
				}
				want = providerintegration.ErrConnection
			case "event_type":
				value.EventType = "refund.created"
			case "payload":
				value.Payload = json.RawMessage(`{"state":"refunded"}`)
			case "provider":
				value.ProviderCode = "other"
				want = providerintegration.ErrConnection
			case "connection":
				value.ConnectionID = "missing"
				want = providerintegration.ErrConnection
			}
			if replay, err := repo.AcceptWebhook(ctx, value); replay || !errors.Is(err, want) {
				t.Errorf("scope accepted or wrong failure: replay=%v error=%v want=%v", replay, err, want)
			}
		})
	}
	if _, err := pool.Exec(ctx, `update integration.provider_connection set state='active' where tenant_id=$1`, tenant); err != nil {
		t.Fatal(err)
	}
	base.Payload = json.RawMessage(`{ "state" : "paid" }`)
	if replay, err := repo.AcceptWebhook(ctx, base); !replay || err != nil {
		t.Fatal("equivalent JSON replay rejected", replay, err)
	}
	t.Run("concurrent_disable", func(t *testing.T) {
		bounded, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		change, err := pool.Begin(bounded)
		if err != nil {
			t.Fatal(err)
		}
		defer change.Rollback(context.Background())
		if _, err := change.Exec(bounded, `update integration.provider_connection set state='disabled' where tenant_id=$1`, tenant); err != nil {
			t.Fatal(err)
		}
		value := base
		value.ProviderEventID = "evt-concurrent-disable"
		type outcome struct {
			replay bool
			err    error
		}
		done := make(chan outcome, 1)
		go func() { replay, err := repo.AcceptWebhook(bounded, value); done <- outcome{replay, err} }()
		ticker := time.NewTicker(20 * time.Millisecond)
		defer ticker.Stop()
		probe, stop := context.WithTimeout(ctx, 5*time.Second)
		defer stop()
		for {
			var blocked bool
			if err := pool.QueryRow(probe, `select exists(select 1 from pg_stat_activity where datname=current_database() and wait_event_type='Lock' and query like 'select connection_id from integration.provider_connection%')`).Scan(&blocked); err != nil {
				t.Fatal("lock probe", err)
			}
			if blocked {
				break
			}
			select {
			case result := <-done:
				t.Fatal("receipt did not wait for current connection", result)
			case <-probe.Done():
				t.Fatal("connection lock not observed")
			case <-ticker.C:
			}
		}
		if err := change.Commit(bounded); err != nil {
			t.Fatal(err)
		}
		select {
		case result := <-done:
			if result.replay || !errors.Is(result.err, providerintegration.ErrConnection) {
				t.Fatal("concurrent disable bypassed", result)
			}
		case <-bounded.Done():
			t.Fatal("receipt did not resume after disable")
		}
	})
	var events, jobs int
	if err := pool.QueryRow(ctx, `select (select count(*) from integration.webhook_event where tenant_id=$1),(select count(*) from platform.job where tenant_id=$1)`, tenant).Scan(&events, &jobs); err != nil || events != 1 || jobs != 1 {
		t.Fatal("scope writes escaped", events, jobs, err)
	}
}

func TestProviderWebhookDeduplicationAndJobAreAtomic(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28d40"
	defer func() {
		_, _ = pool.Exec(ctx, `delete from platform.job where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from integration.webhook_event where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from integration.reconciliation_item where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from integration.external_mapping where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from integration.provider_connection where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	}()
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'provider-test','Provider Test','Provider')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref) values($1,'primary','sandbox','SANDBOX_SECRET')`, tenant); err != nil {
		t.Fatal(err)
	}
	repository := NewProviderIntegration(pool)
	receipt := providerintegration.Receipt{TenantID: tenant, ConnectionID: "primary", ProviderCode: "sandbox", ProviderEventID: "evt-1", EventType: "payment.updated", BodyHash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Payload: json.RawMessage(`{"state":"paid"}`)}
	if replayed, err := repository.AcceptWebhook(ctx, receipt); err != nil || replayed {
		t.Fatalf("first replay=%v err=%v", replayed, err)
	}
	if replayed, err := repository.AcceptWebhook(ctx, receipt); err != nil || !replayed {
		t.Fatalf("second replay=%v err=%v", replayed, err)
	}
	receipt.BodyHash = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	if _, err := repository.AcceptWebhook(ctx, receipt); !errors.Is(err, providerintegration.ErrConflict) {
		t.Fatalf("conflict err=%v", err)
	}
	var events, jobs int
	if err = pool.QueryRow(ctx, `select count(*) from integration.webhook_event where tenant_id=$1`, tenant).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.job where tenant_id=$1 and job_type='provider.webhook.received'`, tenant).Scan(&jobs); err != nil {
		t.Fatal(err)
	}
	if events != 1 || jobs != 1 {
		t.Fatalf("events=%d jobs=%d", events, jobs)
	}
	reconciliation := providerintegration.NewReconciliationService(repository)
	if err = reconciliation.UpsertMapping(ctx, providerintegration.Mapping{TenantID: tenant, ProviderCode: "sandbox", ResourceType: "offer", InternalID: "variant-1", ExternalID: "external-1", VersionToken: "v1"}); err != nil {
		t.Fatal(err)
	}
	if err = reconciliation.Record(ctx, providerintegration.Reconciliation{TenantID: tenant, ID: "recon-1", ProviderCode: "sandbox", ResourceType: "offer", InternalID: "variant-1", ExternalID: "external-1", State: "state-mismatch", Evidence: json.RawMessage(`{"internal":"active","external":"paused"}`), ObservedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if err = reconciliation.Resolve(ctx, tenant, "sandbox", "recon-1"); err != nil {
		t.Fatal(err)
	}
	var mapping, resolved int
	if err = pool.QueryRow(ctx, `select count(*) from integration.external_mapping where tenant_id=$1 and internal_id='variant-1' and external_id='external-1'`, tenant).Scan(&mapping); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from integration.reconciliation_item where tenant_id=$1 and reconciliation_id='recon-1' and state='resolved' and resolved_at is not null`, tenant).Scan(&resolved); err != nil {
		t.Fatal(err)
	}
	if mapping != 1 || resolved != 1 {
		t.Fatalf("mapping=%d resolved=%d", mapping, resolved)
	}
}
````

### FILE: `internal/platform/httpapi/providerintegration.go`
```yaml
block_id: "GO-PROVIDER-INTEGRATION:http:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "11f14b4b0e835461a083cacc41a2bc70ef105a1171c2960864f5ea3855c0ced8"
variables: []
secrets_allowed: false
```
````go
package httpapi

import (
	"errors"
	"io"
	"mime"
	"net/http"

	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/providerintegration"
)

type ProviderIntegrationModule struct{ Service *providerintegration.Service }

func (m ProviderIntegrationModule) Register(mux *http.ServeMux, _ identity.Verifier) {
	mux.HandleFunc("POST /v1/integrations/{provider}/connections/{connection}/webhooks", m.receive)
}

func (m ProviderIntegrationModule) receive(w http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 256<<10))
	if err != nil {
		writeProblem(w, 400, "INVALID_WEBHOOK", "webhook body is invalid or too large")
		return
	}
	replayed, err := m.Service.Receive(r.Context(), r.PathValue("provider"), r.PathValue("connection"), r.Header.Get("X-Elite-Webhook-Timestamp"), r.Header.Get("X-Elite-Webhook-Signature"), body)
	switch {
	case errors.Is(err, providerintegration.ErrUnauthenticated):
		writeProblem(w, 401, "WEBHOOK_UNAUTHENTICATED", "webhook signature is invalid")
	case errors.Is(err, providerintegration.ErrConnection):
		writeProblem(w, 404, "PROVIDER_CONNECTION_NOT_FOUND", "provider connection was not found")
	case errors.Is(err, providerintegration.ErrConflict):
		writeProblem(w, 409, "PROVIDER_EVENT_CONFLICT", "provider event identifier belongs to another body")
	case err != nil:
		writeProblem(w, 400, "INVALID_WEBHOOK", "webhook does not match the admitted contract")
	case replayed:
		w.Header().Set("Idempotency-Replayed", "true")
		writeJSON(w, 200, map[string]string{"status": "accepted"})
	default:
		writeJSON(w, 202, map[string]string{"status": "accepted"})
	}
}
````

### FILE: `internal/platform/httpapi/providerintegration_test.go`
```yaml
block_id: "GO-PROVIDER-INTEGRATION:http-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "9cd0b886e66c46f01ea24bfc1c0d26f4d9d0a8b5fc8c3b0289a10c26f1f0f518"
variables: []
secrets_allowed: false
```
````go
package httpapi

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/providerintegration"
)

type webhookRepository struct{ calls int }

func (r *webhookRepository) AcceptWebhook(_ context.Context, _ providerintegration.Receipt) (bool, error) {
	r.calls++
	return r.calls > 1, nil
}

func TestProviderWebhookHTTPAdmissionAndReplay(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	secret := []byte("01234567890123456789012345678901")
	registry, err := providerintegration.LoadHMACConnections([]byte(`[{"tenant_id":"tenant","connection_id":"primary","provider_code":"sandbox","secret_env":"SANDBOX_SECRET"}]`), func(string) (string, bool) { return string(secret), true })
	if err != nil {
		t.Fatal(err)
	}
	repository := &webhookRepository{}
	service := providerintegration.NewService(repository, registry, providerintegration.HMACSHA256{Now: func() time.Time { return now }, Tolerance: 5 * time.Minute})
	mux := http.NewServeMux()
	ProviderIntegrationModule{Service: service}.Register(mux, nil)
	body := `{"event_id":"evt-1","type":"payment.updated","data":{}}`
	for index, expected := range []int{202, 200} {
		request := httptest.NewRequest("POST", "/v1/integrations/sandbox/connections/primary/webhooks", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-Elite-Webhook-Timestamp", "1800000000")
		request.Header.Set("X-Elite-Webhook-Signature", providerSignature(secret, "1800000000", []byte(body)))
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		if response.Code != expected {
			t.Fatalf("index=%d status=%d body=%s", index, response.Code, response.Body.String())
		}
	}
	request := httptest.NewRequest("POST", "/v1/integrations/sandbox/connections/primary/webhooks", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Elite-Webhook-Timestamp", "1800000000")
	request.Header.Set("X-Elite-Webhook-Signature", "v1=00")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 401 {
		t.Fatalf("tampered status=%d", response.Code)
	}
}

func providerSignature(secret []byte, timestamp string, body []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(timestamp + "."))
	_, _ = mac.Write(body)
	return "v1=" + hex.EncodeToString(mac.Sum(nil))
}
````

## 6. Configuration surface

| Input | Type/default | Secret | Validation/effect |
|---|---|---|---|
| connection registry | strict JSON / omitted | no | known tenant, connection, provider and secret variable only |
| referenced verifier secret | bytes, minimum 32 / none | yes | external injection; missing/weak fails startup |
| timestamp/signature/body | provider request | signature is sensitive | exact verifier, clock and body-size policy |
| queue/event mapping | project adapter | no | unknown provider/event fails without side effect |

## 7. Dependency bill

| Dependency | Pin | Use | License | Scope | Source |
|---|---|---|---|---|---|
| Go | `1.26.7` | verifier/HTTP/tests | BSD-3-Clause | build/runtime | `go.dev` |
| PostgreSQL | `18.6` verified baseline | connections/inbox/jobs | PostgreSQL | runtime/test | `postgresql.org` |
| composed OIDC/pgx modules | versions fixed by backend `go.mod` | shared runtime infrastructure | Apache-2.0/MIT | build/runtime | official repositories |

## 8. Apply order

Apply migration 0004 after foundation/domain migrations, materialize provider packages, register admitted connections, inject secrets, wire the application and workers, then run signature/replay/conflict/atomicity tests. Existing webhooks require a cutover/replay plan. Rollback stops ingress and workers while retaining inbox/jobs; remove schema only in a disposable environment.

## 9. Verification

- Provisionar cada `integration.provider_connection` mediante migración/seed auditado; `secret_ref` contiene sólo el nombre de la referencia.
- Inyectar el secreto real desde el secret manager del target; mínimo 32 bytes para el adapter HMAC.
- Usar una URL de conexión no predecible además de la firma, HTTPS, rate limiting en el edge y límites de cuerpo.
- Conservar el cuerpo exacto hasta verificar la firma; después persistir sólo el payload admitido y su hash.
- Un mismo `provider_event_id` con igual hash es replay exitoso; con otro hash es conflicto.
- La recepción y el job durable se confirman en una sola transacción.
- El worker específico del proveedor debe mapear el evento a comandos internos idempotentes y ejecutar reconciliación periódica.
- Gates mínimos: migración up/down/up, unit tests de firma/timestamp/tamper/config, PostgreSQL duplicate/conflict/atomic-job, HTTP 202/200/401, `gofmt`, `go test`, `go vet` y builds.

## 10. Reconstruction evidence

Clean rebuild, migration and negative verifier/replay gates are recorded in `reconstruction_evidence/GO_PROVIDER_INTEGRATION_CORE_2026-08-24_V1.md`; final evidence rechecks the current integrated application.

V402 composed delta: Payment/initial-handover composition. New behavior and tests belong to the explicit runtime/portal packs; source and library release claims remain bounded to their evidence. Existing source provenance is preserved.

Canonical V402 integration: selected by the current profile with exact dependencies and caller overlays. Metadata promotion records byte reconstruction, not closure of every admission/release gate. Payload provenance is unchanged.
