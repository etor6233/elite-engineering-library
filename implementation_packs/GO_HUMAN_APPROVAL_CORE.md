# Go Human Approval Core

## 1. Metadata

```yaml
pack_id: "GO-HUMAN-APPROVAL-CORE"
pack_version: "0.9.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Existing approval policy plus a shared durable exact-payload, distinct-human approval owner for WhatsApp replies and Facebook Page publish/revoke; no provider effects."
stacks: ["Go 1.26.8", "PostgreSQL 18.6"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "PG-TX-FOUNDATION 0.1.x", "GO-CONVERSATIONAL-AGENT 0.1.x"]
incompatible_with: ["aprobación automática de dinero", "revisor igual al solicitante", "sin evidencia (sha) ligada"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://www.postgresql.org/docs/18/"]
verified_at: "2026-09-02"
```

Este pack gobierna el paso humano que confirma reservas, ventas, devoluciones y pagos. No ejecuta el efecto: decide si el efecto puede ejecutarse y deja decisión inmutable. La ejecución del pago/devolución sigue en los workers de dominio; esta cola es la compuerta anti-estafa.

## 2. Applicability

Use este pack cuando reservas/ventas/devoluciones/pagos deban pasar por aprobación humana con separación de deberes y doble control para evitar estafas. Es consumido por el agente (handoff) y por el backend (antes de ejecutar efectos).

Rechace este pack para: efectos automáticos sin revisión; o aprobación sin identidad de revisor.

## 3. Architecture contract

- **Ownership**: `internal/approval` gobierna la cola y la política. La migración `0042` gobierna `approval.request` y `approval.decision`; `0058` protege los payloads/decisiones de los tres nuevos kinds manuales. No se afirma inmutabilidad SQL de los kinds históricos. La ejecución del efecto pertenece a los workers de dominio.
- **Invariantes**: (1) dinero (`payment`/`refund`) nunca auto-aprueba. (2) revisor distinto del solicitante (separación de deberes). (3) montos >= umbral requieren dos revisores distintos (doble control). (4) velocidad por sujeto acotada. (5) decisión inmutable y ligada a request.
- **Data flow**: `Submit` → auto-aprobar (solo bajo umbral y no-dinero) o `pending` → `Approve`/`Reject` (con separación/doble control) → decisión inmutable.
- **Failure modes**: duplicado → `ErrDuplicate`; velocidad → `ErrSubjectFlagged`; revisor = solicitante → `ErrSeparation`; revisor repetido → `ErrDuplicateApprover`; no pendiente → `ErrNotPending`.
- **Seguridad/privacidad**: tenant-scoped; evidencia ligada por sha; decisiones append-only.
- **Performance budget**: O(1) por operación en memoria; índice de pendientes en SQL.
- **Operación/migración/rollback**: migración `0042` up/down atómica; la cola se reconstruye desde las decisiones.

## 4. Exact file manifest

```text
CREATE internal/approval/approval.go
CREATE internal/approval/registry.go
CREATE internal/approval/approval_test.go
CREATE db/migrations/0042_human_approval.up.sql
CREATE db/migrations/0042_human_approval.down.sql
CREATE db/tests/0042_human_approval.test.sql
CREATE internal/approval/payload.go
CREATE internal/approval/payload_test.go
CREATE internal/platform/postgres/human_approval.go
CREATE db/migrations/0058_human_approval_payload.up.sql
CREATE db/migrations/0058_human_approval_payload.down.sql
```

## 5. Materialization blocks

### FILE: `internal/approval/approval.go`
```yaml
block_id: "GO-HUMAN-APPROVAL-CORE:internal/approval/approval.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "59eb5529eddd79f78bcd514d11c7c434d55523de7930691405e693d0e56fced8"
variables: []
secrets_allowed: false
```
````go
// Package approval provides a governed human-approval queue for business
// effects (reservations, sales, refunds, payments): separation of duties,
// dual control for high-value requests, velocity-based fraud flags, immutable
// decisions and fail-closed money handling.
package approval

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Kind is the bounded set of approvable business effects.
type Kind string

const (
	KindReservation          Kind = "reservation"
	KindSale                 Kind = "sale"
	KindRefund               Kind = "refund"
	KindPayment              Kind = "payment"
	KindWhatsAppReply        Kind = "whatsapp_reply"
	KindSocialPublish        Kind = "social_publish"
	KindSocialRevoke         Kind = "social_revoke"
	KindStoredValueOperation Kind = "stored_value_operation"
	KindWarrantyRepair       Kind = "warranty_repair"
	KindSerialQuality        Kind = "serial_quality"
	KindCatalogReview        Kind = "catalog_review"
	KindTrainingAssessment   Kind = "training_assessment"
	KindMarketplaceMutation  Kind = "marketplace_mutation"
	KindWhatsAppSchedule     Kind = "whatsapp_schedule"
	KindDocumentReview       Kind = "document_review"
)

// State is the lifecycle of an approval request.
type State string

const (
	StatePending  State = "pending"
	StateApproved State = "approved"
	StateRejected State = "rejected"
)

var hex64Re = regexp.MustCompile(`^[0-9a-f]{64}$`)

var (
	ErrInvalidRequest    = errors.New("approval: invalid request")
	ErrDuplicate         = errors.New("approval: duplicate request")
	ErrSubjectFlagged    = errors.New("approval: subject flagged (velocity)")
	ErrNotFound          = errors.New("approval: not found")
	ErrNotPending        = errors.New("approval: not pending")
	ErrSeparation        = errors.New("approval: reviewer must differ from requester")
	ErrDuplicateApprover = errors.New("approval: reviewer already approved")
)

// Policy configures auto-approval and dual-control thresholds.
type Policy struct {
	AutoApproveMinorUnits int64 // <= threshold auto-approves; 0 disables auto-approval
	DualControlMinorUnits int64 // >= threshold requires two distinct reviewers
	MaxOpenPerSubject     int   // velocity bound per subject; 0 disables
}

// Request is an immutable, tenant-scoped approval request.
type Request struct {
	TenantID         string
	ID               string
	Kind             Kind
	SubjectID        string
	AmountMinorUnits int64
	Requester        string
	EvidenceSHA      string
}

// Validate enforces the request contract.
func (r Request) Validate() error {
	if strings.TrimSpace(r.TenantID) == "" {
		return fmt.Errorf("%w: tenant", ErrInvalidRequest)
	}
	if strings.TrimSpace(r.ID) == "" || len(r.ID) > 128 {
		return fmt.Errorf("%w: id", ErrInvalidRequest)
	}
	switch r.Kind {
	case KindReservation, KindSale, KindRefund, KindPayment, KindWhatsAppReply, KindSocialPublish, KindSocialRevoke, KindStoredValueOperation, KindWarrantyRepair, KindSerialQuality, KindCatalogReview, KindTrainingAssessment, KindMarketplaceMutation, KindWhatsAppSchedule, KindDocumentReview:
	default:
		return fmt.Errorf("%w: kind", ErrInvalidRequest)
	}
	if strings.TrimSpace(r.SubjectID) == "" || len(r.SubjectID) > 128 {
		return fmt.Errorf("%w: subject", ErrInvalidRequest)
	}
	if r.AmountMinorUnits < 0 {
		return fmt.Errorf("%w: amount", ErrInvalidRequest)
	}
	if strings.TrimSpace(r.Requester) == "" || len(r.Requester) > 128 {
		return fmt.Errorf("%w: requester", ErrInvalidRequest)
	}
	if !hex64Re.MatchString(strings.ToLower(r.EvidenceSHA)) {
		return fmt.Errorf("%w: evidence sha", ErrInvalidRequest)
	}
	return nil
}

// Decision is an immutable approval/rejection record.
type Decision struct {
	RequestID string
	Reviewer  string
	Approved  bool
	Reason    string
	At        time.Time
}
````

### FILE: `internal/approval/registry.go`
```yaml
block_id: "GO-HUMAN-APPROVAL-CORE:internal/approval/registry.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "d402b0c5d301bf1d3f2c9068a4439df209b77ef84bb25131f1d59bc5075b9d86"
variables: []
secrets_allowed: false
```
````go
package approval

import (
	"sync"
	"time"
)

type record struct {
	req       Request
	state     State
	approvers map[string]bool
}

// Registry is the approval state machine. Money effects (payment/refund) never
// auto-approve; high-value requests require dual control; velocity is bounded.
type Registry struct {
	mu        sync.Mutex
	policy    Policy
	records   map[string]*record
	openCount map[string]int
	decisions []Decision
	clock     func() time.Time
}

// NewRegistry returns a registry with the given policy.
func NewRegistry(p Policy) *Registry {
	return &Registry{
		policy:    p,
		records:   make(map[string]*record),
		openCount: make(map[string]int),
		clock:     time.Now,
	}
}

func (r *Registry) now() time.Time {
	if r.clock != nil {
		return r.clock()
	}
	return time.Now().UTC()
}

func key(tenant, id string) string          { return tenant + "\x00" + id }
func subjKey(tenant, subject string) string { return tenant + "\x00" + subject }

// Submit validates and registers a request. It auto-approves only low-value,
// non-money effects; everything else stays pending. Velocity and duplicates
// fail closed.
func (r *Registry) Submit(req Request) (State, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.records[key(req.TenantID, req.ID)]; ok {
		return "", ErrDuplicate
	}
	sk := subjKey(req.TenantID, req.SubjectID)
	if r.policy.MaxOpenPerSubject > 0 && r.openCount[sk] >= r.policy.MaxOpenPerSubject {
		return "", ErrSubjectFlagged
	}
	rec := &record{req: req, state: StatePending, approvers: map[string]bool{}}
	r.records[key(req.TenantID, req.ID)] = rec
	r.openCount[sk]++
	if r.autoApprove(req) {
		rec.state = StateApproved
		r.openCount[sk]--
		r.decisions = append(r.decisions, Decision{
			RequestID: req.ID, Reviewer: "system", Approved: true, Reason: "auto", At: r.now(),
		})
	}
	return rec.state, nil
}

func (r *Registry) autoApprove(req Request) bool {
	if r.policy.AutoApproveMinorUnits <= 0 {
		return false
	}
	if req.Kind == KindPayment || req.Kind == KindRefund || req.Kind == KindWhatsAppReply || req.Kind == KindSocialPublish || req.Kind == KindSocialRevoke || req.Kind == KindStoredValueOperation || req.Kind == KindWarrantyRepair || req.Kind == KindSerialQuality || req.Kind == KindCatalogReview || req.Kind == KindTrainingAssessment || req.Kind == KindMarketplaceMutation || req.Kind == KindWhatsAppSchedule || req.Kind == KindDocumentReview {
		return false // money movements always require a human
	}
	return req.AmountMinorUnits <= r.policy.AutoApproveMinorUnits
}

// Approve records a distinct reviewer. Dual-control requests stay pending until
// two distinct reviewers have approved.
func (r *Registry) Approve(tenant, id, reviewer, reason string) (State, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.records[key(tenant, id)]
	if !ok {
		return "", ErrNotFound
	}
	if rec.state != StatePending {
		return "", ErrNotPending
	}
	if reviewer == "" || reviewer == rec.req.Requester {
		return "", ErrSeparation
	}
	if rec.approvers[reviewer] {
		return "", ErrDuplicateApprover
	}
	rec.approvers[reviewer] = true
	needsDual := r.policy.DualControlMinorUnits > 0 && rec.req.AmountMinorUnits >= r.policy.DualControlMinorUnits
	if needsDual && len(rec.approvers) < 2 {
		return StatePending, nil // awaiting a second reviewer
	}
	rec.state = StateApproved
	r.openCount[subjKey(tenant, rec.req.SubjectID)]--
	r.decisions = append(r.decisions, Decision{
		RequestID: id, Reviewer: reviewer, Approved: true, Reason: reason, At: r.now(),
	})
	return StateApproved, nil
}

// Reject records a rejection (separation of duties enforced).
func (r *Registry) Reject(tenant, id, reviewer, reason string) (State, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.records[key(tenant, id)]
	if !ok {
		return "", ErrNotFound
	}
	if rec.state != StatePending {
		return "", ErrNotPending
	}
	if reviewer == "" || reviewer == rec.req.Requester {
		return "", ErrSeparation
	}
	rec.state = StateRejected
	r.openCount[subjKey(tenant, rec.req.SubjectID)]--
	r.decisions = append(r.decisions, Decision{
		RequestID: id, Reviewer: reviewer, Approved: false, Reason: reason, At: r.now(),
	})
	return StateRejected, nil
}

// State returns the current state of a request.
func (r *Registry) State(tenant, id string) (State, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.records[key(tenant, id)]
	if !ok {
		return "", false
	}
	return rec.state, true
}

// Audit returns a copy of the decision log.
func (r *Registry) Audit() []Decision {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Decision, len(r.decisions))
	copy(out, r.decisions)
	return out
}
````

### FILE: `internal/approval/approval_test.go`
```yaml
block_id: "GO-HUMAN-APPROVAL-CORE:internal/approval/approval_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "a42ae28920611105f5c9b6aa13b248511d0112fcb0fa3ffdb3b756d715696eea"
variables: []
secrets_allowed: false
```
````go
package approval

import (
	"errors"
	"strings"
	"testing"
)

func req(id, kind string, amount int64) Request {
	return Request{
		TenantID: "tenant-a", ID: id, Kind: Kind(kind), SubjectID: "customer-1",
		AmountMinorUnits: amount, Requester: "seller-1",
		EvidenceSHA: strings.Repeat("a", 64),
	}
}

func policy() Policy {
	return Policy{
		AutoApproveMinorUnits: 1000,
		DualControlMinorUnits: 50000,
		MaxOpenPerSubject:     3,
	}
}

func TestSubmitAutoApprovesLowValueNonMoney(t *testing.T) {
	r := NewRegistry(policy())
	st, err := r.Submit(req("r1", "sale", 500))
	if err != nil || st != StateApproved {
		t.Fatalf("expected auto-approve, got %q/%v", st, err)
	}
	if len(r.Audit()) != 1 || r.Audit()[0].Reviewer != "system" {
		t.Fatalf("expected one system decision, got %v", r.Audit())
	}
}

func TestMoneyNeverAutoApproves(t *testing.T) {
	r := NewRegistry(policy())
	for _, k := range []string{"payment", "refund"} {
		st, err := r.Submit(req("r-"+k, k, 1))
		if err != nil || st != StatePending {
			t.Fatalf("money %s should stay pending, got %q/%v", k, st, err)
		}
	}
}

func TestVelocityFlag(t *testing.T) {
	r := NewRegistry(policy())
	for i := 0; i < 3; i++ {
		if _, err := r.Submit(req("v"+string(rune('0'+i)), "sale", 10000)); err != nil {
			t.Fatalf("unexpected error before bound: %v", err)
		}
	}
	if _, err := r.Submit(req("v3", "sale", 10000)); !errors.Is(err, ErrSubjectFlagged) {
		t.Fatalf("expected velocity flag, got %v", err)
	}
}

func TestDuplicateRejected(t *testing.T) {
	r := NewRegistry(policy())
	_, _ = r.Submit(req("d1", "sale", 500))
	if _, err := r.Submit(req("d1", "sale", 500)); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("expected duplicate, got %v", err)
	}
}

func TestApproveSeparationOfDuties(t *testing.T) {
	r := NewRegistry(policy())
	_, _ = r.Submit(req("a1", "sale", 10000))
	if _, err := r.Approve("tenant-a", "a1", "seller-1", "ok"); !errors.Is(err, ErrSeparation) {
		t.Fatalf("expected separation error, got %v", err)
	}
	st, err := r.Approve("tenant-a", "a1", "manager-1", "ok")
	if err != nil || st != StateApproved {
		t.Fatalf("expected approve, got %q/%v", st, err)
	}
}

func TestDualControlRequiresTwoReviewers(t *testing.T) {
	r := NewRegistry(policy())
	_, _ = r.Submit(req("dc1", "sale", 60000)) // >= DualControlMinorUnits
	st, err := r.Approve("tenant-a", "dc1", "manager-1", "first")
	if err != nil || st != StatePending {
		t.Fatalf("expected still pending after one reviewer, got %q/%v", st, err)
	}
	if _, err := r.Approve("tenant-a", "dc1", "manager-1", "dup"); !errors.Is(err, ErrDuplicateApprover) {
		t.Fatalf("expected duplicate approver error, got %v", err)
	}
	st, err = r.Approve("tenant-a", "dc1", "manager-2", "second")
	if err != nil || st != StateApproved {
		t.Fatalf("expected approved after two reviewers, got %q/%v", st, err)
	}
}

func TestReject(t *testing.T) {
	r := NewRegistry(policy())
	_, _ = r.Submit(req("rj1", "sale", 10000))
	st, err := r.Reject("tenant-a", "rj1", "manager-1", "fraudulent")
	if err != nil || st != StateRejected {
		t.Fatalf("expected reject, got %q/%v", st, err)
	}
	if _, err := r.Approve("tenant-a", "rj1", "manager-2", "late"); !errors.Is(err, ErrNotPending) {
		t.Fatalf("expected not-pending after reject, got %v", err)
	}
}
````

### FILE: `db/migrations/0042_human_approval.up.sql`
```yaml
block_id: "GO-HUMAN-APPROVAL-CORE:db/migrations/0042_human_approval.up.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "4b412bb3cbab53c54b3a23aa8f9f8e2d7b60424d293cd5f814c8e9bee21af31b"
variables: []
secrets_allowed: false
```
````sql
begin;

create schema if not exists approval;

create table approval.request (
  tenant_id uuid not null,
  request_id text not null,
  kind text not null check (kind in ('reservation', 'sale', 'refund', 'payment')),
  subject_id text not null,
  amount_minor_units bigint not null check (amount_minor_units >= 0),
  requester text not null,
  evidence_sha text not null check (evidence_sha ~ '^[0-9a-f]{64}$'),
  state text not null default 'pending' check (state in ('pending', 'approved', 'rejected')),
  created_at timestamptz not null default clock_timestamp(),
  decided_at timestamptz,
  primary key (tenant_id, request_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (length(request_id) between 1 and 128),
  check (length(subject_id) between 1 and 128),
  check (length(requester) between 1 and 128),
  check ((state <> 'pending') = (decided_at is not null))
);

create table approval.decision (
  tenant_id uuid not null,
  request_id text not null,
  reviewer text not null,
  approved boolean not null,
  reason text not null default '',
  decided_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, request_id, reviewer, decided_at),
  foreign key (tenant_id, request_id) references approval.request (tenant_id, request_id),
  check (length(reviewer) between 1 and 128)
);

create index approval_request_pending_idx
  on approval.request (tenant_id, subject_id) where state = 'pending';

commit;
````

### FILE: `db/migrations/0042_human_approval.down.sql`
```yaml
block_id: "GO-HUMAN-APPROVAL-CORE:db/migrations/0042_human_approval.down.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "d456ea67ca5311540603caf6def7b7f5b44be8847e7dfb0ea27445fde8e01355"
variables: []
secrets_allowed: false
```
````sql
begin;

drop table if exists approval.decision;
drop table if exists approval.request;
drop schema if exists approval;

commit;
````

### FILE: `db/tests/0042_human_approval.test.sql`
```yaml
block_id: "GO-HUMAN-APPROVAL-CORE:db/tests/0042_human_approval.test.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "24e3a6741e7d8241e1eff57839294cf771fa226e9326db55bcd0c0b3a04458ce"
variables: []
secrets_allowed: false
```
````sql
-- 0042_human_approval.test.sql — verifica cola de aprobación: estado pendiente
-- con decided_at nulo, decisión inmutable ligada a request, y dinero (payment)
-- que permanece pendiente (nunca auto-aprobado).
-- Precondición: migración 0001 (platform.tenant) y 0042 aplicadas.

begin;

insert into platform.tenant (tenant_id, tenant_code, legal_name, display_name) values
  ('11111111-1111-1111-1111-111111111111', 'tenant-a', 'A', 'Tenant A');

insert into approval.request
  (tenant_id, request_id, kind, subject_id, amount_minor_units, requester, evidence_sha)
values
  ('11111111-1111-1111-1111-111111111111', 'req-1', 'sale', 'sub-1', 500, 'seller-1', repeat('a',64)),
  ('11111111-1111-1111-1111-111111111111', 'req-2', 'payment', 'sub-2', 10000, 'seller-1', repeat('b',64));

-- dos pendientes con decided_at nulo
do $$
declare n int;
begin
  select count(*) into n from approval.request
   where tenant_id = '11111111-1111-1111-1111-111111111111'
     and state = 'pending' and decided_at is null;
  if n <> 2 then raise exception 'expected 2 pending, got %', n; end if;
end $$;

-- aprobar req-1 y registrar decisión
update approval.request set state = 'approved', decided_at = clock_timestamp()
 where tenant_id = '11111111-1111-1111-1111-111111111111' and request_id = 'req-1';

insert into approval.decision (tenant_id, request_id, reviewer, approved, reason) values
  ('11111111-1111-1111-1111-111111111111', 'req-1', 'manager-1', true, 'ok');

do $$
declare n int;
begin
  select count(*) into n from approval.decision
   where tenant_id = '11111111-1111-1111-1111-111111111111' and request_id = 'req-1';
  if n <> 1 then raise exception 'expected 1 decision, got %', n; end if;
end $$;

-- el dinero (payment) permanece pendiente
do $$
declare st text;
begin
  select state into st from approval.request
   where tenant_id = '11111111-1111-1111-1111-111111111111' and request_id = 'req-2';
  if st <> 'pending' then raise exception 'payment should remain pending, got %', st; end if;
end $$;

rollback;
````


### FILE: `internal/approval/payload.go`

```yaml
block_id: "HUMAN-APPROVAL-EXTENSION:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "cd344c29efe4272efb8eeaa16c4ff132c88c70c04d6d7a539c3e0e1abbf28356"
variables: []
secrets_allowed: false
```

````go
package approval

// AUTHORED serialization glue for durable approval of exact payloads.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"strings"
	"unicode/utf8"
)

func CanonicalPayload(raw []byte) (json.RawMessage, string, error) {
	if len(raw) == 0 || len(raw) > 32768 || !utf8.Valid(raw) {
		return nil, "", ErrInvalidRequest
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	if payloadKeys(d, 0) != nil {
		return nil, "", ErrInvalidRequest
	}
	if _, e := d.Token(); e != io.EOF {
		return nil, "", ErrInvalidRequest
	}
	d = json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	var value map[string]any
	if d.Decode(&value) != nil || value == nil {
		return nil, "", ErrInvalidRequest
	}
	canonical, e := json.Marshal(value)
	if e != nil || len(canonical) > 32768 {
		return nil, "", ErrInvalidRequest
	}
	h := sha256.Sum256(canonical)
	return canonical, hex.EncodeToString(h[:]), nil
}
func payloadKeys(d *json.Decoder, depth int) error {
	if depth > 16 {
		return ErrInvalidRequest
	}
	t, e := d.Token()
	if e != nil {
		return e
	}
	v, ok := t.(json.Delim)
	if !ok {
		return nil
	}
	switch v {
	case '{':
		seen := map[string]bool{}
		for d.More() {
			k, e := d.Token()
			if e != nil {
				return e
			}
			s, ok := k.(string)
			if !ok {
				return ErrInvalidRequest
			}
			s = strings.ToLower(s)
			if seen[s] {
				return ErrInvalidRequest
			}
			seen[s] = true
			if payloadKeys(d, depth+1) != nil {
				return ErrInvalidRequest
			}
		}
	case '[':
		for d.More() {
			if payloadKeys(d, depth+1) != nil {
				return ErrInvalidRequest
			}
		}
	default:
		return ErrInvalidRequest
	}
	_, e = d.Token()
	return e
}
````

### FILE: `internal/approval/payload_test.go`

```yaml
block_id: "HUMAN-APPROVAL-EXTENSION:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "dde4c772d3a09831a0c8642e7d9b44b52feda1e38f4f97fcb5edda7df64c2ab8"
variables: []
secrets_allowed: false
```

````go
package approval

import (
	"strings"
	"testing"
)

func TestBoundPayloadAndManualKinds(t *testing.T) {
	a, h, e := CanonicalPayload([]byte(`{"b":2,"a":"reviewed"}`))
	if e != nil || string(a) != `{"a":"reviewed","b":2}` || len(h) != 64 {
		t.Fatal(string(a), h, e)
	}
	for _, raw := range []string{`{"a":1,"a":2}`, `{"a":1,"A":2}`, `{} {}`, `null`, `[]`, `{"a":{"x":1,"X":2}}`} {
		if _, _, e = CanonicalPayload([]byte(raw)); e == nil {
			t.Fatal("invalid payload", raw)
		}
	}
	for _, kind := range []Kind{KindWhatsAppReply, KindSocialPublish, KindSocialRevoke} {
		r := NewRegistry(Policy{AutoApproveMinorUnits: 100})
		state, e := r.Submit(Request{TenantID: "tenant", ID: "request", Kind: kind, SubjectID: "page", Requester: "requester", EvidenceSHA: strings.Repeat("a", 64)})
		if e != nil || state != StatePending {
			t.Fatal("manual kind auto-approved", kind, state, e)
		}
	}
}
````

### FILE: `internal/platform/postgres/human_approval.go`

```yaml
block_id: "HUMAN-APPROVAL-EXTENSION:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "511c953727e80f52e8a7d84582146e48e7d812eeb378ad3b45f68f01790123c3"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED persistence/transaction glue around approval.Registry. The typed
// caller supplies its permission and same-transaction business binding guard.
import (
	"context"
	"encoding/json"
	"errors"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HumanApprovalSpec struct {
	Request        approval.Request
	OrganizationID string
	Payload        json.RawMessage
}
type HumanApprovalGuard func(context.Context, pgx.Tx) error
type HumanApprovals struct{ pool *pgxpool.Pool }

func NewHumanApprovals(pool *pgxpool.Pool) *HumanApprovals { return &HumanApprovals{pool} }
func boundApprovalKind(k approval.Kind) bool {
	return k == approval.KindWhatsAppReply || k == approval.KindSocialPublish || k == approval.KindSocialRevoke || k == approval.KindStoredValueOperation || k == approval.KindWarrantyRepair || k == approval.KindSerialQuality || k == approval.KindCatalogReview || k == approval.KindTrainingAssessment || k == approval.KindMarketplaceMutation || k == approval.KindWhatsAppSchedule || k == approval.KindDocumentReview
}
func approvalPrincipal(p identity.Principal, tenant, org, permission string) bool {
	return permission != "" && p.Subject != "" && len(p.Subject) <= 128 && p.TenantID == tenant && p.Allowed(permission) && p.AllowedOrganization(org)
}
func (s *HumanApprovals) Submit(ctx context.Context, p identity.Principal, v HumanApprovalSpec, permission string, guard HumanApprovalGuard) (bool, error) {

	if s == nil || s.pool == nil {
		return false, approval.ErrInvalidRequest
	}
	tx, e := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if e != nil {
		return false, e
	}
	defer tx.Rollback(ctx)
	replay, e := s.submitTx(ctx, tx, p, v, permission, guard)
	if e != nil {
		return false, e
	}
	return replay, tx.Commit(ctx)
}

// submitTx reuses the same approval admission under a caller-owned local transaction.
// It does not commit, approve, or run external provider effects.
func (s *HumanApprovals) submitTx(ctx context.Context, tx pgx.Tx, p identity.Principal, v HumanApprovalSpec, permission string, guard HumanApprovalGuard) (bool, error) {
	canonical, hash, e := approval.CanonicalPayload(v.Payload)
	if s == nil || s.pool == nil || e != nil || v.Request.Validate() != nil || !boundApprovalKind(v.Request.Kind) || v.Request.AmountMinorUnits != 0 || v.Request.Requester != p.Subject || v.Request.EvidenceSHA != hash || !approvalPrincipal(p, v.Request.TenantID, v.OrganizationID, permission) {
		return false, approval.ErrInvalidRequest
	}
	if guard != nil {
		if e = guard(ctx, tx); e != nil {
			return false, e
		}
	}
	tag, e := tx.Exec(ctx, `insert into approval.request(tenant_id,request_id,kind,subject_id,amount_minor_units,requester,evidence_sha,organization_id,payload)
 values($1,$2,$3,$4,0,$5,$6,$7,$8) on conflict do nothing`, v.Request.TenantID, v.Request.ID, v.Request.Kind, v.Request.SubjectID, v.Request.Requester, hash, v.OrganizationID, canonical)
	if e != nil {
		return false, e
	}
	var same bool
	e = tx.QueryRow(ctx, `select kind=$3 and subject_id=$4 and amount_minor_units=0 and requester=$5 and evidence_sha=$6 and organization_id=$7 and payload=$8::jsonb from approval.request where tenant_id=$1 and request_id=$2`, v.Request.TenantID, v.Request.ID, v.Request.Kind, v.Request.SubjectID, v.Request.Requester, hash, v.OrganizationID, canonical).Scan(&same)
	if e != nil || !same {
		return false, approval.ErrDuplicate
	}
	return tag.RowsAffected() == 0, nil
}
func (s *HumanApprovals) Decide(ctx context.Context, p identity.Principal, tenant, id, org, expectedSHA string, approved bool, reason, permission string, guard HumanApprovalGuard) (approval.State, error) {
	if s == nil || s.pool == nil || len(reason) > 2048 || !approvalPrincipal(p, tenant, org, permission) {
		return "", approval.ErrInvalidRequest
	}
	tx, e := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if e != nil {
		return "", e
	}
	defer tx.Rollback(ctx)
	state, e := s.decideTx(ctx, tx, p, tenant, id, org, expectedSHA, approved, reason, permission, guard)
	if e != nil {
		return "", e
	}
	return state, tx.Commit(ctx)
}

// AUTHORED caller-controlled transaction composition; registry decision owner retained.
func (s *HumanApprovals) decideTx(ctx context.Context, tx pgx.Tx, p identity.Principal, tenant, id, org, expectedSHA string, approved bool, reason, permission string, guard HumanApprovalGuard) (approval.State, error) {
	if s == nil || s.pool == nil || len(reason) > 2048 || !approvalPrincipal(p, tenant, org, permission) {
		return "", approval.ErrInvalidRequest
	}
	v, state, e := readHumanApproval(ctx, tx, tenant, id, true)
	if e != nil {
		return "", e
	}
	if v.OrganizationID != org || v.Request.EvidenceSHA != expectedSHA || !boundApprovalKind(v.Request.Kind) {
		return "", approval.ErrInvalidRequest
	}
	if state != approval.StatePending {
		return state, approval.ErrNotPending
	}
	// Registry owns separation and the one-human decision transition. Zero policy
	// intentionally disables automatic/threshold/velocity decisions for this glue.
	registry := approval.NewRegistry(approval.Policy{})
	if _, e = registry.Submit(v.Request); e != nil {
		return "", e
	}
	if approved {
		state, e = registry.Approve(tenant, id, p.Subject, reason)
	} else {
		state, e = registry.Reject(tenant, id, p.Subject, reason)
	}
	if e != nil {
		return "", e
	}
	if guard != nil {
		if e = guard(ctx, tx); e != nil {
			return "", e
		}
	}
	_, e = tx.Exec(ctx, `insert into approval.decision(tenant_id,request_id,reviewer,approved,reason) values($1,$2,$3,$4,$5)`, tenant, id, p.Subject, approved, reason)
	if e != nil {
		return "", e
	}
	_, e = tx.Exec(ctx, `update approval.request set state=$3,decided_at=clock_timestamp() where tenant_id=$1 and request_id=$2 and state='pending'`, tenant, id, state)
	if e != nil {
		return "", e
	}
	return state, nil
}
func readHumanApproval(ctx context.Context, tx pgx.Tx, tenant, id string, lock bool) (HumanApprovalSpec, approval.State, error) {
	var v HumanApprovalSpec
	var state approval.State
	q := `select tenant_id::text,request_id,kind,subject_id,amount_minor_units,requester,evidence_sha,organization_id,payload,state from approval.request where tenant_id=$1 and request_id=$2`
	if lock {
		q += ` for update`
	}
	e := tx.QueryRow(ctx, q, tenant, id).Scan(&v.Request.TenantID, &v.Request.ID, &v.Request.Kind, &v.Request.SubjectID, &v.Request.AmountMinorUnits, &v.Request.Requester, &v.Request.EvidenceSHA, &v.OrganizationID, &v.Payload, &state)
	if errors.Is(e, pgx.ErrNoRows) {
		e = approval.ErrNotFound
	}
	return v, state, e
}

// LookupHumanApproval holds the request row through the caller's fence commit.
// It verifies durable approval only; domain/expiry/provider checks stay typed.
func LookupHumanApproval(ctx context.Context, tx pgx.Tx, tenant, id, org string, kind approval.Kind, expectedSHA string) (HumanApprovalSpec, error) {
	v, state, e := readHumanApproval(ctx, tx, tenant, id, true)
	if e != nil {
		return v, e
	}
	if state != approval.StateApproved || v.OrganizationID != org || v.Request.Kind != kind || v.Request.EvidenceSHA != expectedSHA {
		return HumanApprovalSpec{}, approval.ErrInvalidRequest
	}
	var reviewed bool
	e = tx.QueryRow(ctx, `select count(*)=1 and coalesce(bool_and(approved and reviewer<>$3),false) from approval.decision where tenant_id=$1 and request_id=$2`, tenant, id, v.Request.Requester).Scan(&reviewed)
	if e != nil {
		return HumanApprovalSpec{}, e
	}
	if !reviewed {
		return HumanApprovalSpec{}, approval.ErrSeparation
	}
	return v, nil
}
````

### FILE: `db/migrations/0058_human_approval_payload.up.sql`

```yaml
block_id: "HUMAN-APPROVAL-EXTENSION:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "94500514edcc534167b08472cc40fff25a46f70c414979db2ba4572b89d84c6b"
variables: []
secrets_allowed: false
```

````sql
begin;
-- AUTHORED persistence bindings for existing approval.Request/Registry. No
-- second approval ledger, provider success claim, or auto-approval is added.
alter table approval.request add column organization_id text;
alter table approval.request add column payload jsonb;
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke'));
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke') or new.kind in ('whatsapp_reply','social_publish','social_revoke') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create trigger bound_request_immutable before update or delete on approval.request
 for each row execute function approval.guard_bound_request();
create function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;
create trigger bound_decision_immutable before update or delete on approval.decision
 for each row execute function approval.guard_bound_decision();
commit;
````

### FILE: `db/migrations/0058_human_approval_payload.down.sql`

```yaml
block_id: "HUMAN-APPROVAL-EXTENSION:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c30a8752780e98be6677df241805c002c97d9469202e11a17400d5ae46404e12"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$begin
 if exists(select 1 from approval.request where kind in ('whatsapp_reply','social_publish','social_revoke')) then
  raise exception 'rollback requires preserving bound approval evidence; no destructive downgrade';
 end if;
end $$;
drop trigger bound_decision_immutable on approval.decision;
drop function approval.guard_bound_decision();
drop trigger bound_request_immutable on approval.request;
drop function approval.guard_bound_request();
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in ('reservation','sale','refund','payment'));
alter table approval.request drop column organization_id;
alter table approval.request drop column payload;
commit;
````


## 6. Configuration surface

Sin variables ni secretos. Los umbrales (`AutoApproveMinorUnits`, `DualControlMinorUnits`, `MaxOpenPerSubject`) se configuran en código como valores exactos.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | core Go | BSD-3-Clause | runtime | https://go.dev |
| PostgreSQL 18.6 | 18.6 | cola/decisiones/índices | PostgreSQL License | runtime | https://www.postgresql.org |

## 8. Apply order

1. Componer `PG-TX-FOUNDATION` (migración `0001`) y el backend.
2. Aplicar migración `0042` sobre PostgreSQL 18.6.
3. Colocar los tres archivos Go bajo `internal/approval/` y los tres SQL bajo `db/`.
4. Verificar: `psql ... -v ON_ERROR_STOP=1 -f db/tests/0042_human_approval.test.sql` y `go test ./... -count=1`.
5. Rollback: migración `0042` down y eliminar `internal/approval/`.

## 9. Verification

- `go test ./internal/approval/ -count=1`: 7/7 PASS (auto-aprobar bajo umbral no-dinero, dinero nunca auto-aprueba, velocidad, duplicado, separación de deberes, doble control con dos revisores, rechazo, no-pendiente).
- `go test ./... -count=1` (7 paquetes): PASS.
- `go vet ./...`: exit 0.
- PostgreSQL 18.6 real (initdb → up → test → down): `ON_ERROR_STOP=1` exit 0; aserciones DO: 2 pendientes, decisión inmutable, dinero permanece pendiente; down deja `approval_ns=true`.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_HUMAN_APPROVAL_CORE_2026-09-02_V182.md`.

V402 composed delta: New durable exact-payload manual approvals for WhatsApp/social; immutable scoped bindings, expected-hash decision, existing Registry separation and shared transaction guards. No automatic new-kind approval or provider side effect.

New durable scope: WhatsApp/social only, amount=0, one_distinct_human. Existing amount thresholds/velocity policies are not newly attributed, promoted or implemented for durable historical money kinds. G0-G8 and final composition conditions are in social-connected-evidence.md.

V402 composed delta: V402 source-backed stored-value integration: exact remaining provider due, explicit payment/funding XOR, shared approval, bounded browser transport and optional host. See STORED_VALUE_OPERATOR_FLOW_V402.md; source/pack admission successor governs final claim. Existing provider-only behavior retained.

V402 composed delta: Connected warranty reuses existing transaction/approval/stock/service owners; SQL ordering and public wrapper behavior retained. Optional host factory fails closed. Exact source tested in WARRANTY_INTERFACE_AND_PORTABILITY_V402.md; no new dependency or corporate attribution.

V402 composed delta: Connected J2 uses existing transaction owners and preserves public operations transitions; serial quality remains in the shared distinct-human approval owner. Optional host hook keeps narrower profiles compatible. No dependency added or corporate attribution. SERIAL_SUPPLY_CONNECTED_RELEASE_V402.md.

V402 composed delta: J3 immutable approved catalog publication reuses original Commerce SQL/shared approval and optional host/public model owner; existing Next storefront consumes a validated published projection. No new dependencies or corporate attribution. CATALOG_CONNECTED_RELEASE_V402.md.

V402 composed delta: T2804 connected training: existing versioned help/audit/shared approvals/outbox/BFF; opt-in host and navigation; bounded body retains original default. No new dependency or automatic grant. TRAINING_CONNECTED_RELEASE_V402.md.

V402 composed delta: T2805 narrow connected Mercado Libre PRICE/STOCK/PAUSE/RESUME for existing User Products item: immutable current catalog + existing serial ATP + exact distinct human approval + shared one-attempt fence + GET-only recovery. MARKETPLACE_MUTATION_RELEASE_V402.md/json. AUTHORED HTTP/SQL/host/proof glue; initial publication/media and other T2805 work remain open.

V402 composed delta: Scheduled WhatsApp exact source/approval/job/host glue and template receipt recovery; SCHEDULED_COMMUNICATIONS_RELEASE_V402.md/json.

V402 composed delta: T2806 connected document reference. Optional host hook, explicit bound document-review kind, exact AWS module closure and preserved security-floor checksums; no unchanged business policy modified. DOCUMENT_REFERENCE_RELEASE_V402.md/json.
