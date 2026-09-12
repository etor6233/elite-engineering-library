# Go PostgreSQL Contact Channel Identity

## 1. Metadata

```yaml
pack_id: "GO-PG-CONTACT-CHANNEL-IDENTITY"
pack_version: "0.1.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa la asociación verificable, versionada y revocable entre una identidad externa de canal y un lead CRM; resuelve scope dinámico sin persistir el identificador externo en claro y demuestra lead→conversación→cotización con replay."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "pgx 5.10.0"]
compatible_with: ["GO-CONNECTED-CONVERSATION-RUNTIME 0.1.x", "GO-APP-WIRING 0.2.x", "GO-OMNICHANNEL-LEAD-INGRESS 0.2.x", "GO-LEAD-CANDIDATE-PROMOTION 0.1.x"]
incompatible_with: ["asociación inferida por texto o coincidencia difusa", "identificador externo persistido en claro", "consentimiento supuesto", "reasignación silenciosa de identidad"]
license_expression: "LicenseRef-Workspace-Owner AND MIT dependency"
upstream_sources: ["https://pages.nist.gov/800-63-4/sp800-63a/accounts/", "https://pages.nist.gov/800-63-4/sp800-63/model/", "https://www.postgresql.org/docs/18/mvcc.html"]
verified_at: "2026-09-04"
```

Todos los bloques son `AUTHORED`. NIST gobierna el claim estrecho de cuenta/sujeto, registros de proofing/consentimiento y bindings; PostgreSQL gobierna serialización, locking y manejo de conflictos. Este código no se atribuye a NIST, PostgreSQL, Meta, Google ni TikTok.

## 2. Applicability

Use este pack después de promover un candidato a `crm.lead` y antes de permitir acciones conversacionales sobre ese lead. El canal o un workflow autorizado debe proporcionar la identidad exacta, la decisión de política y su evidencia. Un lead de publicidad no queda vinculado automáticamente a WhatsApp, SMS, email ni otra identidad por similitud de nombre, teléfono o correo.

## 3. Architecture contract

- **Ownership**: `contactidentity` valida y deriva identidades; PostgreSQL posee la proyección vigente y el ledger de decisiones; `conversationruntime` sólo consume un `ContactResolver`.
- **Privacidad**: la identidad externa se transforma con HMAC-SHA-256 tenant+canal scoped mediante una clave de al menos 32 bytes; ni la proyección ni la evidencia durable almacenan el valor crudo.
- **Autorización**: el binding referencia un lead existente y el resolver obtiene la organización desde ese lead; nunca la recibe del modelo ni del mensaje.
- **Cambio**: creación exige versión esperada cero; política/PII/revocación usan compare-and-set; lead y subject no pueden reasignarse silenciosamente; cada decisión es append-only e idempotente por `request_id`+hash.
- **Concurrencia**: transacciones serializables y constraints convierten dos creaciones rivales en un único éxito y un conflicto explícito.
- **E2E**: el test de app usa el resolver PostgreSQL real; un replay conserva una sola cotización y una sola cadena LLM de dos llamadas.

## 4. Exact file manifest

```text
CREATE internal/contactidentity/identity.go
CREATE internal/contactidentity/identity_test.go
CREATE internal/platform/postgres/contact_identity.go
CREATE internal/platform/postgres/contact_identity_integration_test.go
CREATE internal/app/contact_identity_e2e_test.go
CREATE db/migrations/0047_contact_channel_identity.up.sql
CREATE db/migrations/0047_contact_channel_identity.down.sql
CREATE db/tests/0047_contact_channel_identity.test.sql
```

## 5. Materialization blocks

### FILE: `internal/contactidentity/identity.go`
```yaml
block_id: "GO-PG-CONTACT-CHANNEL-IDENTITY:internal/contactidentity/identity.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "3ff24f2a1350da62b880302c0b93500d6027de872d80bb08b7a36fd713b44155"
variables: []
secrets_allowed: false
```
````go
// Package contactidentity defines the evidence-bound association between a
// provider channel identity and an existing CRM lead. It never infers identity
// or consent from message text.
package contactidentity

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalid  = errors.New("contactidentity: invalid command")
	ErrConflict = errors.New("contactidentity: version or replay conflict")
	ErrNotFound = errors.New("contactidentity: active binding not found")
	channelRE   = regexp.MustCompile(`^[a-z][a-z0-9_]{0,31}$`)
	hexRE       = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

type State string

const (
	StateActive  State = "active"
	StateRevoked State = "revoked"
)

type Command struct {
	TenantID        string
	ChannelCode     string
	ExternalID      string
	LeadID          string
	SubjectID       string
	PIIAllowed      bool
	State           State
	PolicyVersion   string
	EvidenceSHA256  string
	EffectiveAt     time.Time
	ExpectedVersion int64
	RequestID       string
}

type Receipt struct {
	Version  int64
	State    State
	Replayed bool
}

func (c Command) Validate() error {
	if strings.TrimSpace(c.TenantID) == "" || !channelRE.MatchString(c.ChannelCode) || strings.TrimSpace(c.ExternalID) == "" || len(c.ExternalID) > 256 || strings.TrimSpace(c.LeadID) == "" || strings.TrimSpace(c.SubjectID) == "" || len(c.SubjectID) > 256 || strings.TrimSpace(c.PolicyVersion) == "" || len(c.PolicyVersion) > 128 || !hexRE.MatchString(c.EvidenceSHA256) || c.EffectiveAt.IsZero() || c.ExpectedVersion < 0 || strings.TrimSpace(c.RequestID) == "" || len(c.RequestID) > 128 {
		return ErrInvalid
	}
	if c.State != StateActive && c.State != StateRevoked {
		return ErrInvalid
	}
	if c.State == StateRevoked && c.PIIAllowed {
		return ErrInvalid
	}
	return nil
}

func ExternalDigest(secret []byte, tenantID, channelCode, externalID string) (string, error) {
	if len(secret) < 32 || strings.TrimSpace(tenantID) == "" || !channelRE.MatchString(channelCode) || strings.TrimSpace(externalID) == "" || len(externalID) > 256 {
		return "", ErrInvalid
	}
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(tenantID))
	_, _ = mac.Write([]byte{0})
	_, _ = mac.Write([]byte(channelCode))
	_, _ = mac.Write([]byte{0})
	_, _ = mac.Write([]byte(externalID))
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func (c Command) RequestSHA256(secret []byte) (string, string, error) {
	if err := c.Validate(); err != nil {
		return "", "", err
	}
	digest, err := ExternalDigest(secret, c.TenantID, c.ChannelCode, c.ExternalID)
	if err != nil {
		return "", "", err
	}
	canonical := struct {
		TenantID, ChannelCode, ExternalDigest, LeadID, SubjectID, PolicyVersion, EvidenceSHA256, EffectiveAt, RequestID string
		PIIAllowed                                                                                                      bool
		State                                                                                                           State
		ExpectedVersion                                                                                                 int64
	}{c.TenantID, c.ChannelCode, digest, c.LeadID, c.SubjectID, c.PolicyVersion, c.EvidenceSHA256, c.EffectiveAt.UTC().Format(time.RFC3339Nano), c.RequestID, c.PIIAllowed, c.State, c.ExpectedVersion}
	b, err := json.Marshal(canonical)
	if err != nil {
		return "", "", err
	}
	sum := sha256.Sum256(b)
	return digest, hex.EncodeToString(sum[:]), nil
}
````

### FILE: `internal/contactidentity/identity_test.go`
```yaml
block_id: "GO-PG-CONTACT-CHANNEL-IDENTITY:internal/contactidentity/identity_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "a7059e66f5d2f089aaa0f3084a042b511055342e1f454d9cafe04e5b537b3b12"
variables: []
secrets_allowed: false
```
````go
package contactidentity

import (
	"strings"
	"testing"
	"time"
)

func validCommand() Command {
	return Command{TenantID: "tenant-a", ChannelCode: "whatsapp", ExternalID: "+5491112345678", LeadID: "lead-1", SubjectID: "lead:lead-1", PIIAllowed: true, State: StateActive, PolicyVersion: "sales-contact-v1", EvidenceSHA256: strings.Repeat("a", 64), EffectiveAt: time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC), RequestID: "bind-1"}
}

func TestDigestIsScopedDeterministicAndDoesNotRevealExternalID(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	a, err := ExternalDigest(key, "tenant-a", "whatsapp", "+5491112345678")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := ExternalDigest(key, "tenant-a", "whatsapp", "+5491112345678")
	c, _ := ExternalDigest(key, "tenant-b", "whatsapp", "+5491112345678")
	if a != b || a == c || strings.Contains(a, "12345678") || len(a) != 64 {
		t.Fatalf("a=%q b=%q c=%q", a, b, c)
	}
}

func TestCommandValidationFailsClosed(t *testing.T) {
	c := validCommand()
	if err := c.Validate(); err != nil {
		t.Fatal(err)
	}
	c.State, c.PIIAllowed = StateRevoked, true
	if err := c.Validate(); err == nil {
		t.Fatal("revoked binding cannot allow PII")
	}
	c = validCommand()
	c.EvidenceSHA256 = "unproven"
	if err := c.Validate(); err == nil {
		t.Fatal("evidence hash required")
	}
	if _, err := ExternalDigest([]byte("short"), "tenant-a", "whatsapp", "id"); err == nil {
		t.Fatal("weak key accepted")
	}
}

func TestRequestHashChangesWithPolicyAndHidesRawIdentity(t *testing.T) {
	c := validCommand()
	key := []byte("0123456789abcdef0123456789abcdef")
	_, a, err := c.RequestSHA256(key)
	if err != nil {
		t.Fatal(err)
	}
	c.PolicyVersion = "sales-contact-v2"
	_, b, err := c.RequestSHA256(key)
	if err != nil {
		t.Fatal(err)
	}
	if a == b || strings.Contains(a, c.ExternalID) {
		t.Fatalf("a=%q b=%q", a, b)
	}
}
````

### FILE: `internal/platform/postgres/contact_identity.go`
```yaml
block_id: "GO-PG-CONTACT-CHANNEL-IDENTITY:internal/platform/postgres/contact_identity.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c82e352a434a987dcbbd72adfd5f1421bc079e835c81609a9f849fe549952a72"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"errors"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/conversationruntime"
	"elite.local/enterprise/internal/domainbind"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactIdentityStore struct {
	pool    *pgxpool.Pool
	hmacKey []byte
}

func NewContactIdentityStore(pool *pgxpool.Pool, hmacKey []byte) (*ContactIdentityStore, error) {
	if pool == nil || len(hmacKey) < 32 {
		return nil, contactidentity.ErrInvalid
	}
	key := append([]byte(nil), hmacKey...)
	return &ContactIdentityStore{pool: pool, hmacKey: key}, nil
}

func (s *ContactIdentityStore) Resolve(ctx context.Context, message channels.Message) (conversationruntime.Contact, error) {
	if s == nil || s.pool == nil || message.Direction != channels.DirectionIn {
		return conversationruntime.Contact{}, contactidentity.ErrInvalid
	}
	digest, err := contactidentity.ExternalDigest(s.hmacKey, message.TenantID, message.ChannelCode, message.ExternalID)
	if err != nil {
		return conversationruntime.Contact{}, err
	}
	var leadID, organizationID, subjectID string
	var piiAllowed bool
	err = s.pool.QueryRow(ctx, `select b.lead_id,l.organization_id,b.subject_id,b.pii_allowed
from communication.contact_channel_binding b
join crm.lead l on l.tenant_id=b.tenant_id and l.lead_id=b.lead_id
where b.tenant_id=$1 and b.channel_code=$2 and b.external_id_hmac=$3 and b.state='active' and b.effective_at<=statement_timestamp()`, message.TenantID, message.ChannelCode, digest).Scan(&leadID, &organizationID, &subjectID, &piiAllowed)
	if errors.Is(err, pgx.ErrNoRows) {
		return conversationruntime.Contact{}, contactidentity.ErrNotFound
	}
	if err != nil {
		return conversationruntime.Contact{}, err
	}
	return conversationruntime.Contact{Scope: domainbind.Scope{OrganizationID: organizationID, LeadID: leadID}, SubjectID: subjectID, PIIAllowed: piiAllowed}, nil
}

func (s *ContactIdentityStore) Apply(ctx context.Context, command contactidentity.Command) (contactidentity.Receipt, error) {
	if s == nil || s.pool == nil {
		return contactidentity.Receipt{}, contactidentity.ErrInvalid
	}
	digest, requestHash, err := command.RequestSHA256(s.hmacKey)
	if err != nil {
		return contactidentity.Receipt{}, identityWriteError(err)
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return contactidentity.Receipt{}, identityWriteError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var replayHash string
	var replayVersion int64
	var replayState contactidentity.State
	err = tx.QueryRow(ctx, `select request_sha256_hex,binding_version,state from communication.contact_channel_binding_decision where tenant_id=$1 and request_id=$2`, command.TenantID, command.RequestID).Scan(&replayHash, &replayVersion, &replayState)
	if err == nil {
		if replayHash != requestHash {
			return contactidentity.Receipt{}, contactidentity.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return contactidentity.Receipt{}, err
		}
		return contactidentity.Receipt{Version: replayVersion, State: replayState, Replayed: true}, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return contactidentity.Receipt{}, err
	}
	var currentVersion int64
	var currentLead, currentSubject string
	err = tx.QueryRow(ctx, `select version,lead_id,subject_id from communication.contact_channel_binding where tenant_id=$1 and channel_code=$2 and external_id_hmac=$3 for update`, command.TenantID, command.ChannelCode, digest).Scan(&currentVersion, &currentLead, &currentSubject)
	exists := err == nil
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return contactidentity.Receipt{}, err
	}
	if command.ExpectedVersion == 0 {
		if exists || command.State != contactidentity.StateActive {
			return contactidentity.Receipt{}, contactidentity.ErrConflict
		}
		currentVersion = 1
		_, err = tx.Exec(ctx, `insert into communication.contact_channel_binding
(tenant_id,channel_code,external_id_hmac,lead_id,subject_id,pii_allowed,state,policy_version,evidence_sha256_hex,effective_at,version)
values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`, command.TenantID, command.ChannelCode, digest, command.LeadID, command.SubjectID, command.PIIAllowed, command.State, command.PolicyVersion, command.EvidenceSHA256, command.EffectiveAt.UTC(), currentVersion)
	} else {
		if !exists || currentVersion != command.ExpectedVersion || currentLead != command.LeadID || currentSubject != command.SubjectID {
			return contactidentity.Receipt{}, contactidentity.ErrConflict
		}
		currentVersion++
		_, err = tx.Exec(ctx, `update communication.contact_channel_binding set pii_allowed=$4,state=$5,policy_version=$6,evidence_sha256_hex=$7,effective_at=$8,version=$9,updated_at=clock_timestamp()
where tenant_id=$1 and channel_code=$2 and external_id_hmac=$3 and version=$10`, command.TenantID, command.ChannelCode, digest, command.PIIAllowed, command.State, command.PolicyVersion, command.EvidenceSHA256, command.EffectiveAt.UTC(), currentVersion, command.ExpectedVersion)
	}
	if err != nil {
		return contactidentity.Receipt{}, identityWriteError(err)
	}
	_, err = tx.Exec(ctx, `insert into communication.contact_channel_binding_decision
(tenant_id,request_id,request_sha256_hex,channel_code,external_id_hmac,binding_version,lead_id,subject_id,pii_allowed,state,policy_version,evidence_sha256_hex,effective_at)
values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`, command.TenantID, command.RequestID, requestHash, command.ChannelCode, digest, currentVersion, command.LeadID, command.SubjectID, command.PIIAllowed, command.State, command.PolicyVersion, command.EvidenceSHA256, command.EffectiveAt.UTC())
	if err != nil {
		return contactidentity.Receipt{}, identityWriteError(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return contactidentity.Receipt{}, identityWriteError(err)
	}
	return contactidentity.Receipt{Version: currentVersion, State: command.State}, nil
}

func identityWriteError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23505" || pgErr.Code == "40001") {
		return contactidentity.ErrConflict
	}
	return err
}
````

### FILE: `internal/platform/postgres/contact_identity_integration_test.go`
```yaml
block_id: "GO-PG-CONTACT-CHANNEL-IDENTITY:internal/platform/postgres/contact_identity_integration_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "fc84707ab1ae6b3a325089f7cf1d8bda2c3981a501f2f3da533ea7f8d10ec4f2"
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
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestContactIdentityEvidenceVersioningResolutionAndRevocation(t *testing.T) {
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
	seed := time.Now().UTC().Format(time.RFC3339Nano)
	tenant := leadstream.StableUUID(t.Name(), seed, "tenant")
	otherTenant := leadstream.StableUUID(t.Name(), seed, "other")
	for _, id := range []string{tenant, otherTenant} {
		if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Identity Test','Identity Test')`, id, "identity-"+id[len(id)-4:]); err != nil {
			t.Fatal(err)
		}
		if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store-1','store-1','Store 1','store')`, id); err != nil {
			t.Fatal(err)
		}
		if _, err = pool.Exec(ctx, `insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required) values($1,'lead-1','store-1','new','meta','{}',true)`, id); err != nil {
			t.Fatal(err)
		}
	}
	key := []byte("0123456789abcdef0123456789abcdef")
	store, err := NewContactIdentityStore(pool, key)
	if err != nil {
		t.Fatal(err)
	}
	when := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)
	command := contactidentity.Command{TenantID: tenant, ChannelCode: "whatsapp", ExternalID: "+5491112345678", LeadID: "lead-1", SubjectID: "lead:lead-1", PIIAllowed: true, State: contactidentity.StateActive, PolicyVersion: "sales-contact-v1", EvidenceSHA256: strings.Repeat("a", 64), EffectiveAt: when, RequestID: "binding-1"}
	receipt, err := store.Apply(ctx, command)
	if err != nil || receipt.Version != 1 || receipt.Replayed {
		t.Fatalf("receipt=%+v err=%v", receipt, err)
	}
	replay, err := store.Apply(ctx, command)
	if err != nil || !replay.Replayed || replay.Version != 1 {
		t.Fatalf("replay=%+v err=%v", replay, err)
	}
	changed := command
	changed.PolicyVersion = "changed"
	if _, err = store.Apply(ctx, changed); !errors.Is(err, contactidentity.ErrConflict) {
		t.Fatalf("expected replay conflict, got %v", err)
	}
	message := channels.Message{ChannelCode: "whatsapp", TenantID: tenant, ExternalID: command.ExternalID, ProviderMessageID: "wamid.1", OccurredAt: when, Direction: channels.DirectionIn, Text: "hola"}
	contact, err := store.Resolve(ctx, message)
	if err != nil || contact.Scope.OrganizationID != "store-1" || contact.Scope.LeadID != "lead-1" || contact.SubjectID != "lead:lead-1" || !contact.PIIAllowed {
		t.Fatalf("contact=%+v err=%v", contact, err)
	}
	message.TenantID = otherTenant
	if _, err = store.Resolve(ctx, message); !errors.Is(err, contactidentity.ErrNotFound) {
		t.Fatalf("cross-tenant resolve=%v", err)
	}
	command.ExpectedVersion, command.RequestID, command.PolicyVersion, command.EvidenceSHA256, command.PIIAllowed = 1, "binding-2", "sales-contact-v2", strings.Repeat("b", 64), false
	receipt, err = store.Apply(ctx, command)
	if err != nil || receipt.Version != 2 {
		t.Fatalf("policy update=%+v err=%v", receipt, err)
	}
	message.TenantID = tenant
	contact, err = store.Resolve(ctx, message)
	if err != nil || contact.PIIAllowed {
		t.Fatalf("updated contact=%+v err=%v", contact, err)
	}
	command.ExpectedVersion, command.RequestID, command.State, command.PolicyVersion, command.EvidenceSHA256 = 2, "binding-3", contactidentity.StateRevoked, "sales-contact-revoked-v1", strings.Repeat("c", 64)
	receipt, err = store.Apply(ctx, command)
	if err != nil || receipt.Version != 3 || receipt.State != contactidentity.StateRevoked {
		t.Fatalf("revoke=%+v err=%v", receipt, err)
	}
	if _, err = store.Resolve(ctx, message); !errors.Is(err, contactidentity.ErrNotFound) {
		t.Fatalf("revoked resolve=%v", err)
	}
	var decisions int
	var leaked bool
	if err = pool.QueryRow(ctx, `select count(*),bool_or(external_id_hmac like '%12345678%') from communication.contact_channel_binding_decision where tenant_id=$1`, tenant).Scan(&decisions, &leaked); err != nil {
		t.Fatal(err)
	}
	if decisions != 3 || leaked {
		t.Fatalf("decisions=%d leaked=%v", decisions, leaked)
	}
	concurrent := contactidentity.Command{TenantID: tenant, ChannelCode: "sms", ExternalID: "+5491187654321", LeadID: "lead-1", SubjectID: "lead:lead-1", State: contactidentity.StateActive, PolicyVersion: "sales-contact-v1", EvidenceSHA256: strings.Repeat("d", 64), EffectiveAt: when, ExpectedVersion: 0}
	results := make(chan error, 2)
	var start sync.WaitGroup
	start.Add(2)
	for i := 0; i < 2; i++ {
		go func(i int) {
			defer start.Done()
			copy := concurrent
			copy.RequestID = []string{"concurrent-a", "concurrent-b"}[i]
			_, applyErr := store.Apply(context.Background(), copy)
			results <- applyErr
		}(i)
	}
	start.Wait()
	close(results)
	succeeded, conflicted := 0, 0
	for applyErr := range results {
		switch {
		case applyErr == nil:
			succeeded++
		case errors.Is(applyErr, contactidentity.ErrConflict):
			conflicted++
		default:
			t.Fatalf("unexpected concurrent result: %v", applyErr)
		}
	}
	if succeeded != 1 || conflicted != 1 {
		t.Fatalf("concurrent succeeded=%d conflicted=%d", succeeded, conflicted)
	}
}
````

### FILE: `internal/app/contact_identity_e2e_test.go`
```yaml
block_id: "GO-PG-CONTACT-CHANNEL-IDENTITY:internal/app/contact_identity_e2e_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "52754b240ac2c89bc15b1b98154649fc7f466b3e2ade4611fb5c4525006aaadb"
variables: []
secrets_allowed: false
```
````go
package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/leadstream"
	platformpostgres "elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCapturedLeadIdentityReachesConversationAndQuoteExactlyOnce(t *testing.T) {
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
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Contact E2E','Contact E2E')`, tenant, "contact-e2e-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'org','org','Contact Store','store')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required) values($1,'lead','org','new','meta','{}',true)`, tenant); err != nil {
		t.Fatal(err)
	}
	identityStore, err := platformpostgres.NewContactIdentityStore(pool, []byte("0123456789abcdef0123456789abcdef"))
	if err != nil {
		t.Fatal(err)
	}
	when := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)
	_, err = identityStore.Apply(ctx, contactidentity.Command{TenantID: tenant, ChannelCode: "whatsapp", ExternalID: "wa-captured-lead", LeadID: "lead", SubjectID: "lead:lead", PIIAllowed: true, State: contactidentity.StateActive, PolicyVersion: "sales-contact-v1", EvidenceSHA256: strings.Repeat("d", 64), EffectiveAt: when.Add(-time.Minute), RequestID: "contact-e2e-binding-1"})
	if err != nil {
		t.Fatal(err)
	}
	turnStore, err := platformpostgres.NewConversationStore(pool, 5*time.Second, 24*time.Hour, 3)
	if err != nil {
		t.Fatal(err)
	}
	llmCalls := 0
	llm := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		llmCalls++
		w.Header().Set("content-type", "application/json")
		if llmCalls == 1 {
			_, _ = w.Write([]byte(`{"id":"resp-contact-1","status":"completed","output":[{"type":"function_call","call_id":"call-contact-1","name":"create_quote","arguments":"{\"product\":\"scooter\",\"quantity\":1}"}],"usage":{"input_tokens":12,"output_tokens":4,"total_tokens":16}}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"resp-contact-2","status":"completed","output":[{"type":"message","content":[{"type":"output_text","text":"Cotización creada."}]}],"usage":{"input_tokens":8,"output_tokens":5,"total_tokens":13}}`))
	}))
	defer llm.Close()
	domainCalls := 0
	domain := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainCalls++
		if r.URL.Path != "/v1/franchise/quotes" || r.Header.Get("Authorization") != "Bearer token" || len(r.Header.Get("Idempotency-Key")) != 64 {
			t.Errorf("domain request path=%q auth=%q key=%q", r.URL.Path, r.Header.Get("Authorization"), r.Header.Get("Idempotency-Key"))
		}
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"quotation_id":"quote-contact-1"}`))
	}))
	defer domain.Close()
	channel := &appChannel{}
	registry := channels.NewRegistry()
	if err = registry.Register(channel); err != nil {
		t.Fatal(err)
	}
	cfg := fullConfig()
	cfg.TenantID, cfg.LLMBaseURL, cfg.DomainBaseURL = tenant, llm.URL, domain.URL
	cfg.ChannelRegistry, cfg.ConversationStore, cfg.ContactResolver = registry, turnStore, identityStore
	a, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	message := channels.Message{ChannelCode: "whatsapp", TenantID: tenant, ExternalID: "wa-captured-lead", ThreadID: "thread-contact", ProviderMessageID: "wamid.contact.1", OccurredAt: when, Direction: channels.DirectionIn, Text: "Cotizame un scooter"}
	if err = a.Dispatcher.Handle(ctx, message); err != nil {
		t.Fatal(err)
	}
	if err = a.Dispatcher.Handle(ctx, message); err != nil {
		t.Fatal(err)
	}
	var state, toolName string
	var attempts int
	if err = pool.QueryRow(ctx, `select state,tool_name,attempt_count from communication.conversation_turn where tenant_id=$1 and channel_code=$2 and provider_message_id=$3`, tenant, message.ChannelCode, message.ProviderMessageID).Scan(&state, &toolName, &attempts); err != nil {
		t.Fatal(err)
	}
	if state != "completed" || toolName != "create_quote" || attempts != 1 || llmCalls != 2 || domainCalls != 1 || len(channel.sent) != 2 {
		t.Fatalf("state=%s tool=%s attempts=%d llm=%d domain=%d sent=%d", state, toolName, attempts, llmCalls, domainCalls, len(channel.sent))
	}
}
````

### FILE: `db/migrations/0047_contact_channel_identity.up.sql`
```yaml
block_id: "GO-PG-CONTACT-CHANNEL-IDENTITY:db/migrations/0047_contact_channel_identity.up.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "fc284e307ca355a0e4caf32e2dc9d1c51ba56f14836fd28074ef7afdfe9fb30c"
variables: []
secrets_allowed: false
```
````sql
begin;

create table communication.contact_channel_binding (
  tenant_id uuid not null references platform.tenant(tenant_id),
  channel_code text not null check (channel_code ~ '^[a-z][a-z0-9_]{0,31}$'),
  external_id_hmac text not null check (external_id_hmac ~ '^[0-9a-f]{64}$'),
  lead_id text not null,
  subject_id text not null check (length(subject_id) between 1 and 256),
  pii_allowed boolean not null,
  state text not null check (state in ('active','revoked')),
  policy_version text not null check (length(policy_version) between 1 and 128),
  evidence_sha256_hex text not null check (evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  effective_at timestamptz not null,
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, channel_code, external_id_hmac),
  foreign key (tenant_id, lead_id) references crm.lead(tenant_id, lead_id),
  check (state <> 'revoked' or pii_allowed = false),
  check (updated_at >= created_at)
);

create table communication.contact_channel_binding_decision (
  tenant_id uuid not null,
  request_id text not null check (length(request_id) between 1 and 128),
  request_sha256_hex text not null check (request_sha256_hex ~ '^[0-9a-f]{64}$'),
  channel_code text not null,
  external_id_hmac text not null,
  binding_version bigint not null check (binding_version > 0),
  lead_id text not null,
  subject_id text not null,
  pii_allowed boolean not null,
  state text not null check (state in ('active','revoked')),
  policy_version text not null,
  evidence_sha256_hex text not null check (evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  effective_at timestamptz not null,
  recorded_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, request_id),
  unique (tenant_id, channel_code, external_id_hmac, binding_version),
  foreign key (tenant_id, channel_code, external_id_hmac) references communication.contact_channel_binding(tenant_id, channel_code, external_id_hmac),
  check (state <> 'revoked' or pii_allowed = false)
);

create function communication.reject_contact_binding_decision_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='contact binding evidence is immutable';
end;
$function$;

create trigger contact_binding_decision_immutable before update or delete on communication.contact_channel_binding_decision
for each row execute function communication.reject_contact_binding_decision_mutation();

create index contact_binding_lead_idx on communication.contact_channel_binding(tenant_id, lead_id) where state='active';

commit;
````

### FILE: `db/migrations/0047_contact_channel_identity.down.sql`
```yaml
block_id: "GO-PG-CONTACT-CHANNEL-IDENTITY:db/migrations/0047_contact_channel_identity.down.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "19237159847b26cd1dd97b547dc5891f674123680848c7a78210be947988f3d2"
variables: []
secrets_allowed: false
```
````sql
begin;
drop table communication.contact_channel_binding_decision;
drop table communication.contact_channel_binding;
drop function communication.reject_contact_binding_decision_mutation();
commit;
````

### FILE: `db/tests/0047_contact_channel_identity.test.sql`
```yaml
block_id: "GO-PG-CONTACT-CHANNEL-IDENTITY:db/tests/0047_contact_channel_identity.test.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "b6bd61cbc7f53cb643d40e3e89a13a14e1a07340abb169c0dc6cac4171c7f1df"
variables: []
secrets_allowed: false
```
````sql
begin;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values('47474747-4747-4747-8747-474747474747','contact-id-test','Contact Identity Test','Contact Identity Test');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
values('47474747-4747-4747-8747-474747474747','store-1','store-1','Store 1','store');
insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required)
values('47474747-4747-4747-8747-474747474747','lead-1','store-1','new','meta','{}',true);
insert into communication.contact_channel_binding
(tenant_id,channel_code,external_id_hmac,lead_id,subject_id,pii_allowed,state,policy_version,evidence_sha256_hex,effective_at,version)
values('47474747-4747-4747-8747-474747474747','whatsapp',repeat('a',64),'lead-1','lead:lead-1',true,'active','policy-1',repeat('b',64),clock_timestamp(),1);
insert into communication.contact_channel_binding_decision
(tenant_id,request_id,request_sha256_hex,channel_code,external_id_hmac,binding_version,lead_id,subject_id,pii_allowed,state,policy_version,evidence_sha256_hex,effective_at)
values('47474747-4747-4747-8747-474747474747','request-1',repeat('c',64),'whatsapp',repeat('a',64),1,'lead-1','lead:lead-1',true,'active','policy-1',repeat('b',64),clock_timestamp());

do $test$
begin
  begin
    update communication.contact_channel_binding_decision set policy_version='changed'
    where tenant_id='47474747-4747-4747-8747-474747474747' and request_id='request-1';
    raise exception 'decision mutation unexpectedly succeeded';
  exception when sqlstate '55000' then null;
  end;
  if exists(select 1 from communication.contact_channel_binding where external_id_hmac like '%54911%') then
    raise exception 'raw channel identity leaked';
  end if;
end;
$test$;

rollback;
````

## 6. Configuration surface

- `CONTACT_IDENTITY_HMAC_KEY`: secreto aleatorio de al menos 32 bytes, entregado por el secret manager; rotación requiere un plan de doble lectura/rebinding aprobado, no reemplazo ciego.
- Cada alta, cambio de política o revocación exige `policy_version`, `evidence_sha256`, timestamp efectivo, versión esperada y request ID estable.
- `PIIAllowed` representa una decisión externa ya probada; el store no inventa finalidad, base legal ni consentimiento.

## 7. Dependency bill

- Go stdlib: HMAC/SHA-256, JSON, validación y tiempo.
- `github.com/jackc/pgx/v5 v5.10.0`, ya fijado por la composición backend.
- PostgreSQL 18.6; depende de `platform.tenant`, `org.organization`, `crm.lead` y del schema `communication` creados por packs anteriores.

## 8. Apply order

Aplicar migraciones base y CRM, luego `0044`/ `0045`, `0046` y finalmente `0047`. Construir `ContactIdentityStore` con el pool y la clave secreta; usarlo como binder en el workflow autorizado y como `app.Config.ContactResolver` en el runtime.

## 9. Verification

1. Materializar los ocho archivos en una composición que incluya sus packs compatibles.
2. Aplicar `0047_contact_channel_identity.up.sql` y ejecutar el SQL focal con `ON_ERROR_STOP=1`.
3. Ejecutar `go test -count=1 ./internal/contactidentity ./internal/platform/postgres ./internal/app` dos veces con `TEST_DATABASE_URL`.
4. Ejecutar `go test ./...`, `go vet ./...` y `go build ./...` sobre la composición completa.
5. Verificar round-trip byte exacto y `gofmt` no-op. El race detector queda como gate de host con CGO; no se afirma en Windows sin CGO.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_PG_CONTACT_CHANNEL_IDENTITY_2026-09-04_V237.md`. La reconstrucción local no prueba cuentas live, propiedad real de una identidad de canal, política jurídica, entrega outbound ni producción.


V402 composed delta: Reject contact bindings whose effective_at is in the future. Real negative test previously reached model/domain; predicate corrected and focused PostgreSQL regression preserved.
