# Go PostgreSQL Outbound Delivery Fence

## 1. Metadata

```yaml
pack_id: "GO-PG-OUTBOUND-DELIVERY-FENCE"
pack_version: "0.2.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un fence durable de salida que invoca al proveedor como máximo una vez automáticamente por DeliveryKey, conserva receipt/evidencia, separa rechazo terminal probado de ambigüedad y exige reconciliación explícita."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "pgx 5.10.0"]
compatible_with: ["GO-CHANNELS-CORE 0.4.x", "GO-APP-WIRING 0.2.x", "GO-PG-CONTACT-CHANNEL-IDENTITY 0.1.x"]
incompatible_with: ["retry ciego después de timeout", "delivery key sin persistencia", "provider receipt no validado", "afirmación exactly-once sobre un proveedor externo"]
license_expression: "LicenseRef-Workspace-Owner AND MIT dependency"
upstream_sources: ["https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/transactional-outbox.html", "https://www.postgresql.org/docs/18/mvcc.html", "https://sre.google/sre-book/reliable-product-launches/"]
verified_at: "2026-09-04"
```

Los once bloques son `AUTHORED`. AWS y PostgreSQL gobiernan los claims estrechos de dual-write, outbox/idempotencia y concurrencia; el código local no se atribuye a AWS, PostgreSQL, Google ni a un proveedor de mensajería.

## 2. Applicability

Use para toda respuesta que cruza hacia WhatsApp, SMS, email u otro proveedor donde perder el ACK podría duplicar contacto o costo. El wrapper requiere un sender que devuelva la identidad y evidencia exactas del proveedor. Si el provider soporta una idempotency key propia, úsela además; este fence no inventa esa capacidad.

`PYTHON-META-WHATSAPP-CLOUD-ADAPTER 0.6.0` incluye `whatsappbridge.Sender`,
que implementa este contrato y conserva el owner Python del POST. V265 prueba
el borde de proceso con fixture sintético y replay/unknown en este PostgreSQL
real. V266 añade ApprovalResolver PostgreSQL ligado al turno/contacto/consentimiento
actuales, probado con fixtures SQL. V267 añade API explícita autenticada usando
este fence: HTTP/RS256/PG y pérdida de respuesta probados con child sintético.
UI/worker, cuenta real y correlación de statuses continúan pendientes; no
constituye notificación live end-to-end.

V268 añade lectura HTTP aislada del resultado local sin invocar Claim ni cambiar
estados/eventos. Un lease vencido se informa como necesidad de reconciliación;
no se confunde accepted con entrega confirmada ni GET con permiso de reenvío.

## 3. Architecture contract

- **Antes del efecto**: inserta `sending` y un evento append-only bajo tenant/canal/DeliveryKey; sólo el creador llama al proveedor.
- **Después del efecto**: receipt válido mueve a `accepted`; replay retorna sin segunda llamada.
- **Rechazo terminal probado**: un sender puede devolver `TerminalFailure` con código y evidencia SHA-256; el store cierra `sending→failed_terminal` sin marcar aceptación ni habilitar retry.
- **Ambigüedad**: error de red, receipt inválido, fallo de commit o lease vencida quedan `unknown`; nunca se reenvían automáticamente.
- **Reconciliación**: un adapter específico puede cerrar `unknown→accepted` o `unknown→failed_terminal` con evidencia; no existe transición que autorice un segundo envío.
- **Privacidad**: recipient y provider message ID se guardan sólo como HMAC tenant/canal; el texto no se persiste en este ledger, sólo su request SHA-256.
- **Producción**: `app.NewProduction` rechaza un registry vacío o con cualquier canal no marcado durable.
- **Semántica honesta**: garantiza un máximo de una invocación automática, no exactamente-once delivery del proveedor.

## 4. Exact file manifest

```text
CREATE internal/outbounddelivery/delivery.go
CREATE internal/outbounddelivery/delivery_test.go
CREATE internal/platform/postgres/outbound_delivery.go
CREATE internal/platform/postgres/outbound_delivery_integration_test.go
CREATE internal/channels/durable_registry.go
CREATE internal/channels/durable_registry_test.go
CREATE internal/app/production.go
CREATE internal/app/durable_delivery_e2e_test.go
CREATE db/migrations/0048_outbound_delivery_fence.up.sql
CREATE db/migrations/0048_outbound_delivery_fence.down.sql
CREATE db/tests/0048_outbound_delivery_fence.test.sql
```

## 5. Materialization blocks

### FILE: `internal/outbounddelivery/delivery.go`
```yaml
block_id: "GO-PG-OUTBOUND-DELIVERY-FENCE:internal/outbounddelivery/delivery.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "f99dc98f1d403f9bbb5def6870f273d0b41a317333f31aa728b3d35a3aec7b6b"
variables: []
secrets_allowed: false
```
````go
// Package outbounddelivery fences provider sends so a lost acknowledgement
// never causes an automatic duplicate. Ambiguous outcomes require explicit
// reconciliation before any further effect.
package outbounddelivery

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"elite.local/enterprise/internal/channels"
)

var (
	ErrInvalid     = errors.New("outbounddelivery: invalid configuration or message")
	ErrConflict    = errors.New("outbounddelivery: payload conflict")
	ErrInProgress  = errors.New("outbounddelivery: delivery already in progress")
	ErrUnknown     = errors.New("outbounddelivery: provider effect unknown; reconcile before retry")
	ErrTerminal    = errors.New("outbounddelivery: delivery terminal")
	hex64RE        = regexp.MustCompile(`^[0-9a-f]{64}$`)
	terminalCodeRE = regexp.MustCompile(`^[A-Z][A-Z0-9_]{2,63}$`)
)

type Claim struct{ Replay bool }

type Receipt struct {
	ProviderMessageID string
	EvidenceSHA256    string
	AcceptedAt        time.Time
}

type TerminalFailure struct {
	Code           string
	EvidenceSHA256 string
	Cause          error
}

func (e *TerminalFailure) Error() string {
	if e == nil {
		return ErrTerminal.Error()
	}
	return ErrTerminal.Error() + ": " + e.Code
}

func (e *TerminalFailure) Unwrap() error { return e.Cause }

func NewTerminalFailure(code, evidenceSHA256 string, cause error) error {
	if !terminalCodeRE.MatchString(code) || !hex64RE.MatchString(evidenceSHA256) {
		return ErrInvalid
	}
	return &TerminalFailure{Code: code, EvidenceSHA256: evidenceSHA256, Cause: cause}
}

func (r Receipt) Validate() error {
	if strings.TrimSpace(r.ProviderMessageID) == "" || len(r.ProviderMessageID) > 256 || !hex64RE.MatchString(r.EvidenceSHA256) || r.AcceptedAt.IsZero() {
		return ErrInvalid
	}
	return nil
}

type Sender interface {
	SendWithReceipt(context.Context, channels.Message) (Receipt, error)
}

type Store interface {
	Claim(context.Context, channels.Message, string) (Claim, error)
	Complete(context.Context, channels.Message, string, Receipt) error
	MarkUnknown(context.Context, channels.Message, string, string) error
	MarkFailed(context.Context, channels.Message, string, string, string) error
}

type Channel struct {
	CodeValue string
	Receiver  channels.Channel
	Sender    Sender
	Store     Store
}

func (c *Channel) Code() string { return c.CodeValue }

func (c *Channel) DurableDelivery() bool {
	return c != nil && c.CodeValue != "" && c.Receiver != nil && c.Sender != nil && c.Store != nil
}

func (c *Channel) Receive(ctx context.Context) ([]channels.Message, error) {
	if !c.DurableDelivery() || c.Receiver.Code() != c.CodeValue {
		return nil, ErrInvalid
	}
	return c.Receiver.Receive(ctx)
}

func (c *Channel) Send(ctx context.Context, message channels.Message) error {
	if !c.DurableDelivery() || c.Receiver.Code() != c.CodeValue || message.ChannelCode != c.CodeValue || message.Direction != channels.DirectionOut {
		return ErrInvalid
	}
	if err := message.Validate(); err != nil {
		return err
	}
	hash, err := MessageSHA256(message)
	if err != nil {
		return err
	}
	claim, err := c.Store.Claim(ctx, message, hash)
	if err != nil || claim.Replay {
		return err
	}
	receipt, sendErr := c.Sender.SendWithReceipt(ctx, message)
	if sendErr != nil {
		var terminal *TerminalFailure
		if errors.As(sendErr, &terminal) {
			markErr := c.Store.MarkFailed(ctx, message, hash, terminal.EvidenceSHA256, terminal.Code)
			return errors.Join(ErrTerminal, sendErr, markErr)
		}
		markErr := c.Store.MarkUnknown(ctx, message, hash, "PROVIDER_CALL_UNCERTAIN")
		return errors.Join(ErrUnknown, sendErr, markErr)
	}
	if err = receipt.Validate(); err != nil {
		markErr := c.Store.MarkUnknown(ctx, message, hash, "PROVIDER_RECEIPT_INVALID")
		return errors.Join(ErrUnknown, err, markErr)
	}
	if err = c.Store.Complete(ctx, message, hash, receipt); err != nil {
		return errors.Join(ErrUnknown, err)
	}
	return nil
}

func MessageSHA256(message channels.Message) (string, error) {
	if err := message.Validate(); err != nil || message.Direction != channels.DirectionOut {
		return "", ErrInvalid
	}
	canonical := struct {
		ChannelCode, TenantID, ExternalID, ThreadID, DeliveryKey, Text string
	}{message.ChannelCode, message.TenantID, message.ExternalID, message.ThreadID, message.DeliveryKey, message.Text}
	raw, err := json.Marshal(canonical)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}
````

### FILE: `internal/outbounddelivery/delivery_test.go`
```yaml
block_id: "GO-PG-OUTBOUND-DELIVERY-FENCE:internal/outbounddelivery/delivery_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c3213a9ea101c3c1d08eed03d1e4cc81549aeff716d7f98a0a109a7614cfea64"
variables: []
secrets_allowed: false
```
````go
package outbounddelivery

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/channels"
)

type memoryStore struct {
	hash, state string
	receipt     Receipt
}

func (s *memoryStore) Claim(_ context.Context, _ channels.Message, hash string) (Claim, error) {
	if s.hash != "" && s.hash != hash {
		return Claim{}, ErrConflict
	}
	if s.state == "accepted" {
		return Claim{Replay: true}, nil
	}
	if s.state == "sending" {
		return Claim{}, ErrInProgress
	}
	if s.state == "unknown" {
		return Claim{}, ErrUnknown
	}
	if s.state == "failed_terminal" {
		return Claim{}, ErrTerminal
	}
	s.hash, s.state = hash, "sending"
	return Claim{}, nil
}
func (s *memoryStore) Complete(_ context.Context, _ channels.Message, hash string, receipt Receipt) error {
	if s.hash != hash || s.state != "sending" {
		return ErrConflict
	}
	s.state, s.receipt = "accepted", receipt
	return nil
}
func (s *memoryStore) MarkUnknown(_ context.Context, _ channels.Message, hash, _ string) error {
	if s.hash != hash || s.state != "sending" {
		return ErrConflict
	}
	s.state = "unknown"
	return nil
}
func (s *memoryStore) MarkFailed(_ context.Context, _ channels.Message, hash, evidence, code string) error {
	if s.hash != hash || s.state != "sending" || !hex64RE.MatchString(evidence) || code == "" {
		return ErrConflict
	}
	s.state = "failed_terminal"
	return nil
}

type testReceiver struct{ code string }

func (r testReceiver) Code() string                                      { return r.code }
func (testReceiver) Receive(context.Context) ([]channels.Message, error) { return nil, nil }
func (testReceiver) Send(context.Context, channels.Message) error {
	return errors.New("direct send must not be used")
}

type testSender struct {
	calls int
	err   error
}

func (s *testSender) SendWithReceipt(context.Context, channels.Message) (Receipt, error) {
	s.calls++
	if s.err != nil {
		return Receipt{}, s.err
	}
	return Receipt{ProviderMessageID: "provider-1", EvidenceSHA256: strings.Repeat("a", 64), AcceptedAt: time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)}, nil
}

func outboundMessage() channels.Message {
	return channels.Message{ChannelCode: "whatsapp", TenantID: "tenant", ExternalID: "contact", ThreadID: "thread", DeliveryKey: strings.Repeat("d", 64), Direction: channels.DirectionOut, Text: "respuesta"}
}

func TestChannelSendsOnceAndReplaysWithoutProvider(t *testing.T) {
	store, sender := &memoryStore{}, &testSender{}
	c := &Channel{CodeValue: "whatsapp", Receiver: testReceiver{code: "whatsapp"}, Sender: sender, Store: store}
	message := outboundMessage()
	if err := c.Send(context.Background(), message); err != nil {
		t.Fatal(err)
	}
	if err := c.Send(context.Background(), message); err != nil {
		t.Fatal(err)
	}
	if sender.calls != 1 || store.state != "accepted" {
		t.Fatalf("calls=%d state=%s", sender.calls, store.state)
	}
	message.Text = "divergent"
	if err := c.Send(context.Background(), message); !errors.Is(err, ErrConflict) {
		t.Fatalf("divergent replay=%v", err)
	}
}

func TestAmbiguousProviderFailureNeverRetriesAutomatically(t *testing.T) {
	store, sender := &memoryStore{}, &testSender{err: errors.New("timeout after write")}
	c := &Channel{CodeValue: "whatsapp", Receiver: testReceiver{code: "whatsapp"}, Sender: sender, Store: store}
	message := outboundMessage()
	if err := c.Send(context.Background(), message); !errors.Is(err, ErrUnknown) {
		t.Fatalf("first=%v", err)
	}
	if err := c.Send(context.Background(), message); !errors.Is(err, ErrUnknown) {
		t.Fatalf("retry=%v", err)
	}
	if sender.calls != 1 || store.state != "unknown" {
		t.Fatalf("calls=%d state=%s", sender.calls, store.state)
	}
}

func TestProvenTerminalProviderFailureClosesWithoutUnknown(t *testing.T) {
	store := &memoryStore{}
	sender := &testSender{err: NewTerminalFailure("PROVIDER_REJECTED", strings.Repeat("e", 64), errors.New("invalid request"))}
	c := &Channel{CodeValue: "whatsapp", Receiver: testReceiver{code: "whatsapp"}, Sender: sender, Store: store}
	message := outboundMessage()
	if err := c.Send(context.Background(), message); !errors.Is(err, ErrTerminal) {
		t.Fatalf("first=%v", err)
	}
	if err := c.Send(context.Background(), message); !errors.Is(err, ErrTerminal) {
		t.Fatalf("replay=%v", err)
	}
	if sender.calls != 1 || store.state != "failed_terminal" {
		t.Fatalf("calls=%d state=%s", sender.calls, store.state)
	}
}

func TestMessageHashBindsRecipientAndPayload(t *testing.T) {
	a, err := MessageSHA256(outboundMessage())
	if err != nil {
		t.Fatal(err)
	}
	b := outboundMessage()
	b.ExternalID = "other"
	bh, _ := MessageSHA256(b)
	if a == bh || len(a) != 64 {
		t.Fatalf("a=%s b=%s", a, bh)
	}
}
````

### FILE: `internal/platform/postgres/outbound_delivery.go`
```yaml
block_id: "GO-PG-OUTBOUND-DELIVERY-FENCE:internal/platform/postgres/outbound_delivery.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "e7161276b124f7bb1a5c0666df59d33b9c0663beb8a8b22d65d47b20ce26a257"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"errors"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/outbounddelivery"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboundDeliveryStore struct {
	pool    *pgxpool.Pool
	hmacKey []byte
	lease   time.Duration
}

func NewOutboundDeliveryStore(pool *pgxpool.Pool, hmacKey []byte, lease time.Duration) (*OutboundDeliveryStore, error) {
	if pool == nil || len(hmacKey) < 32 || lease < time.Second || lease > 10*time.Minute {
		return nil, outbounddelivery.ErrInvalid
	}
	return &OutboundDeliveryStore{pool: pool, hmacKey: append([]byte(nil), hmacKey...), lease: lease}, nil
}

func (s *OutboundDeliveryStore) Claim(ctx context.Context, message channels.Message, requestHash string) (outbounddelivery.Claim, error) {
	return s.claimWithAdmission(ctx, message, requestHash, nil)
}

func (s *OutboundDeliveryStore) claimWithAdmission(ctx context.Context, message channels.Message, requestHash string, admit func(context.Context, pgx.Tx) error) (outbounddelivery.Claim, error) {
	if s == nil || s.pool == nil || len(requestHash) != 64 {
		return outbounddelivery.Claim{}, outbounddelivery.ErrInvalid
	}
	recipientHash, err := contactidentity.ExternalDigest(s.hmacKey, message.TenantID, message.ChannelCode, message.ExternalID)
	if err != nil {
		return outbounddelivery.Claim{}, err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return outbounddelivery.Claim{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err := tx.Exec(ctx, `insert into communication.outbound_delivery
(tenant_id,channel_code,delivery_key,request_sha256_hex,recipient_hmac,state,attempt_count,locked_until)
values($1,$2,$3,$4,$5,'sending',1,clock_timestamp()+$6*interval '1 millisecond') on conflict do nothing`, message.TenantID, message.ChannelCode, message.DeliveryKey, requestHash, recipientHash, s.lease.Milliseconds())
	if err != nil {
		return outbounddelivery.Claim{}, outboundWriteError(err)
	}
	if result.RowsAffected() == 1 {
		if admit != nil {
			if err = admit(ctx, tx); err != nil {
				return outbounddelivery.Claim{}, err
			}
		}
		_, err = tx.Exec(ctx, `insert into communication.outbound_delivery_event(tenant_id,channel_code,delivery_key,sequence,state) values($1,$2,$3,1,'sending')`, message.TenantID, message.ChannelCode, message.DeliveryKey)
		if err != nil {
			return outbounddelivery.Claim{}, outboundWriteError(err)
		}
		if err = tx.Commit(ctx); err != nil {
			return outbounddelivery.Claim{}, outboundWriteError(err)
		}
		return outbounddelivery.Claim{}, nil
	}
	var storedHash, state string
	var lockedUntil *time.Time
	err = tx.QueryRow(ctx, `select request_sha256_hex,state,locked_until from communication.outbound_delivery where tenant_id=$1 and channel_code=$2 and delivery_key=$3 for update`, message.TenantID, message.ChannelCode, message.DeliveryKey).Scan(&storedHash, &state, &lockedUntil)
	if err != nil {
		return outbounddelivery.Claim{}, err
	}
	if storedHash != requestHash {
		return outbounddelivery.Claim{}, outbounddelivery.ErrConflict
	}
	switch state {
	case "accepted":
		if err = tx.Commit(ctx); err != nil {
			return outbounddelivery.Claim{}, err
		}
		return outbounddelivery.Claim{Replay: true}, nil
	case "unknown":
		return outbounddelivery.Claim{}, outbounddelivery.ErrUnknown
	case "failed_terminal":
		return outbounddelivery.Claim{}, outbounddelivery.ErrTerminal
	case "sending":
		if lockedUntil != nil && lockedUntil.After(time.Now()) {
			return outbounddelivery.Claim{}, outbounddelivery.ErrInProgress
		}
		_, err = tx.Exec(ctx, `update communication.outbound_delivery set state='unknown',locked_until=null,failure_code='LEASE_EXPIRED',updated_at=clock_timestamp() where tenant_id=$1 and channel_code=$2 and delivery_key=$3 and state='sending'`, message.TenantID, message.ChannelCode, message.DeliveryKey)
		if err != nil {
			return outbounddelivery.Claim{}, err
		}
		_, err = tx.Exec(ctx, `insert into communication.outbound_delivery_event(tenant_id,channel_code,delivery_key,sequence,state,failure_code) select tenant_id,channel_code,delivery_key,2,'unknown','LEASE_EXPIRED' from communication.outbound_delivery where tenant_id=$1 and channel_code=$2 and delivery_key=$3`, message.TenantID, message.ChannelCode, message.DeliveryKey)
		if err != nil {
			return outbounddelivery.Claim{}, outboundWriteError(err)
		}
		if err = tx.Commit(ctx); err != nil {
			return outbounddelivery.Claim{}, outboundWriteError(err)
		}
		return outbounddelivery.Claim{}, outbounddelivery.ErrUnknown
	default:
		return outbounddelivery.Claim{}, outbounddelivery.ErrInvalid
	}
}

func (s *OutboundDeliveryStore) Complete(ctx context.Context, message channels.Message, requestHash string, receipt outbounddelivery.Receipt) error {
	if err := receipt.Validate(); err != nil {
		return err
	}
	providerHash, err := contactidentity.ExternalDigest(s.hmacKey, message.TenantID, message.ChannelCode, receipt.ProviderMessageID)
	if err != nil {
		return err
	}
	return s.transition(ctx, message, requestHash, "sending", "accepted", providerHash, receipt.EvidenceSHA256, "", receipt.AcceptedAt.UTC())
}

func (s *OutboundDeliveryStore) MarkUnknown(ctx context.Context, message channels.Message, requestHash, code string) error {
	if code == "" {
		return outbounddelivery.ErrInvalid
	}
	return s.transition(ctx, message, requestHash, "sending", "unknown", "", "", code, time.Time{})
}

func (s *OutboundDeliveryStore) MarkFailed(ctx context.Context, message channels.Message, requestHash, evidenceHash, code string) error {
	if len(evidenceHash) != 64 || code == "" {
		return outbounddelivery.ErrInvalid
	}
	return s.transition(ctx, message, requestHash, "sending", "failed_terminal", "", evidenceHash, code, time.Time{})
}

func (s *OutboundDeliveryStore) ReconcileAccepted(ctx context.Context, message channels.Message, requestHash string, receipt outbounddelivery.Receipt) error {
	if err := receipt.Validate(); err != nil {
		return err
	}
	providerHash, err := contactidentity.ExternalDigest(s.hmacKey, message.TenantID, message.ChannelCode, receipt.ProviderMessageID)
	if err != nil {
		return err
	}
	return s.transition(ctx, message, requestHash, "unknown", "accepted", providerHash, receipt.EvidenceSHA256, "", receipt.AcceptedAt.UTC())
}

func (s *OutboundDeliveryStore) ReconcileFailed(ctx context.Context, message channels.Message, requestHash, evidenceHash, code string) error {
	if len(evidenceHash) != 64 || code == "" {
		return outbounddelivery.ErrInvalid
	}
	return s.transition(ctx, message, requestHash, "unknown", "failed_terminal", "", evidenceHash, code, time.Time{})
}

func (s *OutboundDeliveryStore) transition(ctx context.Context, message channels.Message, requestHash, from, to, providerHash, evidenceHash, code string, acceptedAt time.Time) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var storedHash, state string
	var sequence int64
	err = tx.QueryRow(ctx, `select request_sha256_hex,state from communication.outbound_delivery where tenant_id=$1 and channel_code=$2 and delivery_key=$3 for update`, message.TenantID, message.ChannelCode, message.DeliveryKey).Scan(&storedHash, &state)
	if errors.Is(err, pgx.ErrNoRows) {
		return outbounddelivery.ErrConflict
	}
	if err != nil {
		return err
	}
	if storedHash != requestHash || state != from {
		return outbounddelivery.ErrConflict
	}
	if err = tx.QueryRow(ctx, `select coalesce(max(sequence),0) from communication.outbound_delivery_event where tenant_id=$1 and channel_code=$2 and delivery_key=$3`, message.TenantID, message.ChannelCode, message.DeliveryKey).Scan(&sequence); err != nil {
		return err
	}
	result, err := tx.Exec(ctx, `update communication.outbound_delivery set state=$4,locked_until=null,provider_message_hmac=nullif($5,''),evidence_sha256_hex=nullif($6,''),failure_code=nullif($7,''),accepted_at=$8,updated_at=clock_timestamp() where tenant_id=$1 and channel_code=$2 and delivery_key=$3 and state=$9 and request_sha256_hex=$10`, message.TenantID, message.ChannelCode, message.DeliveryKey, to, providerHash, evidenceHash, code, nullableTime(acceptedAt), from, requestHash)
	if err != nil {
		return outboundWriteError(err)
	}
	if result.RowsAffected() != 1 {
		return outbounddelivery.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into communication.outbound_delivery_event(tenant_id,channel_code,delivery_key,sequence,state,evidence_sha256_hex,failure_code) values($1,$2,$3,$4,$5,nullif($6,''),nullif($7,''))`, message.TenantID, message.ChannelCode, message.DeliveryKey, sequence+1, to, evidenceHash, code)
	if err != nil {
		return outboundWriteError(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return outboundWriteError(err)
	}
	return nil
}

func nullableTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value
}

func outboundWriteError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23505" || pgErr.Code == "40001") {
		return outbounddelivery.ErrConflict
	}
	return err
}
````

### FILE: `internal/platform/postgres/outbound_delivery_integration_test.go`
```yaml
block_id: "GO-PG-OUTBOUND-DELIVERY-FENCE:internal/platform/postgres/outbound_delivery_integration_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "14051a303a11bba752751eb65e91f6ddd0a0973976ff1b6096b1f60deed3f43a"
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

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/leadstream"
	"elite.local/enterprise/internal/outbounddelivery"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestOutboundDeliveryFenceReplayUnknownAndReconciliation(t *testing.T) {
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
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Outbound','Outbound')`, tenant, "outbound-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	key := []byte("0123456789abcdef0123456789abcdef")
	store, err := NewOutboundDeliveryStore(pool, key, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	message := channels.Message{ChannelCode: "whatsapp", TenantID: tenant, ExternalID: "recipient-secret", ThreadID: "thread", DeliveryKey: strings.Repeat("d", 64), Direction: channels.DirectionOut, Text: "respuesta"}
	hash, err := outbounddelivery.MessageSHA256(message)
	if err != nil {
		t.Fatal(err)
	}
	claim, err := store.Claim(ctx, message, hash)
	if err != nil || claim.Replay {
		t.Fatalf("claim=%+v err=%v", claim, err)
	}
	if _, err = store.Claim(ctx, message, hash); !errors.Is(err, outbounddelivery.ErrInProgress) {
		t.Fatalf("second claim=%v", err)
	}
	receipt := outbounddelivery.Receipt{ProviderMessageID: "wamid.secret", EvidenceSHA256: strings.Repeat("a", 64), AcceptedAt: time.Now().UTC()}
	if err = store.Complete(ctx, message, hash, receipt); err != nil {
		t.Fatal(err)
	}
	claim, err = store.Claim(ctx, message, hash)
	if err != nil || !claim.Replay {
		t.Fatalf("replay=%+v err=%v", claim, err)
	}
	changed := message
	changed.Text = "different"
	changedHash, _ := outbounddelivery.MessageSHA256(changed)
	if _, err = store.Claim(ctx, changed, changedHash); !errors.Is(err, outbounddelivery.ErrConflict) {
		t.Fatalf("divergent=%v", err)
	}

	unknown := message
	unknown.DeliveryKey = strings.Repeat("e", 64)
	unknownHash, _ := outbounddelivery.MessageSHA256(unknown)
	if _, err = store.Claim(ctx, unknown, unknownHash); err != nil {
		t.Fatal(err)
	}
	if err = store.MarkUnknown(ctx, unknown, unknownHash, "PROVIDER_TIMEOUT"); err != nil {
		t.Fatal(err)
	}
	if _, err = store.Claim(ctx, unknown, unknownHash); !errors.Is(err, outbounddelivery.ErrUnknown) {
		t.Fatalf("unknown replay=%v", err)
	}
	if err = store.ReconcileAccepted(ctx, unknown, unknownHash, receipt); err != nil {
		t.Fatal(err)
	}
	claim, err = store.Claim(ctx, unknown, unknownHash)
	if err != nil || !claim.Replay {
		t.Fatalf("reconciled=%+v err=%v", claim, err)
	}
	var attempts, events int
	var leaked bool
	if err = pool.QueryRow(ctx, `select d.attempt_count,count(e.*),d.recipient_hmac like '%recipient-secret%' or d.provider_message_hmac like '%wamid%' from communication.outbound_delivery d join communication.outbound_delivery_event e using(tenant_id,channel_code,delivery_key) where d.tenant_id=$1 and d.delivery_key=$2 group by d.attempt_count,d.recipient_hmac,d.provider_message_hmac`, tenant, message.DeliveryKey).Scan(&attempts, &events, &leaked); err != nil {
		t.Fatal(err)
	}
	if attempts != 1 || events != 2 || leaked {
		t.Fatalf("attempts=%d events=%d leaked=%v", attempts, events, leaked)
	}
	expired := message
	expired.DeliveryKey = strings.Repeat("f", 64)
	expiredHash, _ := outbounddelivery.MessageSHA256(expired)
	if _, err = store.Claim(ctx, expired, expiredHash); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `update communication.outbound_delivery set locked_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and delivery_key=$2`, tenant, expired.DeliveryKey); err != nil {
		t.Fatal(err)
	}
	if _, err = store.Claim(ctx, expired, expiredHash); !errors.Is(err, outbounddelivery.ErrUnknown) {
		t.Fatalf("expired lease=%v", err)
	}
	if err = pool.QueryRow(ctx, `select attempt_count from communication.outbound_delivery where tenant_id=$1 and delivery_key=$2`, tenant, expired.DeliveryKey).Scan(&attempts); err != nil || attempts != 1 {
		t.Fatalf("expired attempts=%d err=%v", attempts, err)
	}
	failed := message
	failed.DeliveryKey = strings.Repeat("9", 64)
	failedHash, _ := outbounddelivery.MessageSHA256(failed)
	if _, err = store.Claim(ctx, failed, failedHash); err != nil {
		t.Fatal(err)
	}
	if err = store.MarkUnknown(ctx, failed, failedHash, "PROVIDER_TIMEOUT"); err != nil {
		t.Fatal(err)
	}
	if err = store.ReconcileFailed(ctx, failed, failedHash, strings.Repeat("8", 64), "PROVIDER_CONFIRMED_FAILED"); err != nil {
		t.Fatal(err)
	}
	if _, err = store.Claim(ctx, failed, failedHash); !errors.Is(err, outbounddelivery.ErrTerminal) {
		t.Fatalf("terminal=%v", err)
	}
}
````

### FILE: `internal/channels/durable_registry.go`
```yaml
block_id: "GO-PG-OUTBOUND-DELIVERY-FENCE:internal/channels/durable_registry.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "ce04d2df69ae78147e8b46f0a9689003b7144a50d3f1ce3ebccafe0ab2b00410"
variables: []
secrets_allowed: false
```
````go
package channels

// DurableDeliveryChannel marks an adapter whose Send path is guarded by a
// durable attempt/receipt store. It is checked by app.NewProduction.
type DurableDeliveryChannel interface {
	Channel
	DurableDelivery() bool
}

// AllDurable reports whether every registered channel proves the durable
// marker. An empty registry is never production-ready.
func (r *Registry) AllDurable() bool {
	if r == nil {
		return false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.channels) == 0 {
		return false
	}
	for _, channel := range r.channels {
		durable, ok := channel.(DurableDeliveryChannel)
		if !ok || !durable.DurableDelivery() {
			return false
		}
	}
	return true
}
````

### FILE: `internal/channels/durable_registry_test.go`
```yaml
block_id: "GO-PG-OUTBOUND-DELIVERY-FENCE:internal/channels/durable_registry_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "2db87b48457778b571d81f2b88b7e7442ae9c0a187a54cd619db6823e8b6d5b8"
variables: []
secrets_allowed: false
```
````go
package channels

import "testing"

type durableFake struct{ *fakeChannel }

func (d durableFake) DurableDelivery() bool { return d.fakeChannel != nil }

func TestRegistryRequiresEveryProductionChannelDurable(t *testing.T) {
	r := NewRegistry()
	if r.AllDurable() {
		t.Fatal("empty registry accepted")
	}
	if err := r.Register(&fakeChannel{code: "email"}); err != nil {
		t.Fatal(err)
	}
	if r.AllDurable() {
		t.Fatal("plain channel accepted")
	}
	r = NewRegistry()
	if err := r.Register(durableFake{fakeChannel: &fakeChannel{code: "whatsapp"}}); err != nil {
		t.Fatal(err)
	}
	if !r.AllDurable() {
		t.Fatal("durable channel rejected")
	}
}
````

### FILE: `internal/app/production.go`
```yaml
block_id: "GO-PG-OUTBOUND-DELIVERY-FENCE:internal/app/production.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "48c2628528b8f5f467d5827af88df5290aa2d564a25b71d5a4c4176b6b0acf4e"
variables: []
secrets_allowed: false
```
````go
package app

import "fmt"

// NewProduction adds the delivery invariant that every configured channel is
// fenced by a durable attempt/receipt store. New remains available for local
// tests and compositions that intentionally do not send externally.
func NewProduction(cfg Config) (*App, error) {
	if cfg.ChannelRegistry == nil || !cfg.ChannelRegistry.AllDurable() {
		return nil, fmt.Errorf("%w: all channels require durable delivery", ErrIncomplete)
	}
	return New(cfg)
}
````

### FILE: `internal/app/durable_delivery_e2e_test.go`
```yaml
block_id: "GO-PG-OUTBOUND-DELIVERY-FENCE:internal/app/durable_delivery_e2e_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "cfd51aaf676c28e2624f15a35b4a234ff5cc496d97938023e2e65ffe4098cc34"
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
	"elite.local/enterprise/internal/outbounddelivery"
	platformpostgres "elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type appReceiptSender struct{ calls int }

func (s *appReceiptSender) SendWithReceipt(context.Context, channels.Message) (outbounddelivery.Receipt, error) {
	s.calls++
	return outbounddelivery.Receipt{ProviderMessageID: "wamid.outbound.1", EvidenceSHA256: strings.Repeat("f", 64), AcceptedAt: time.Date(2026, 9, 4, 12, 1, 0, 0, time.UTC)}, nil
}

func TestProductionJourneyFencesOutboundReplay(t *testing.T) {
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
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Durable E2E','Durable E2E')`, tenant, "durable-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'org','org','Store','store')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required) values($1,'lead','org','new','meta','{}',true)`, tenant); err != nil {
		t.Fatal(err)
	}
	key := []byte("0123456789abcdef0123456789abcdef")
	identityStore, err := platformpostgres.NewContactIdentityStore(pool, key)
	if err != nil {
		t.Fatal(err)
	}
	when := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)
	_, err = identityStore.Apply(ctx, contactidentity.Command{TenantID: tenant, ChannelCode: "whatsapp", ExternalID: "wa-durable", LeadID: "lead", SubjectID: "lead:lead", PIIAllowed: true, State: contactidentity.StateActive, PolicyVersion: "sales-contact-v1", EvidenceSHA256: strings.Repeat("a", 64), EffectiveAt: when, RequestID: "durable-binding"})
	if err != nil {
		t.Fatal(err)
	}
	turnStore, err := platformpostgres.NewConversationStore(pool, 5*time.Second, 24*time.Hour, 3)
	if err != nil {
		t.Fatal(err)
	}
	deliveryStore, err := platformpostgres.NewOutboundDeliveryStore(pool, key, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	llmCalls := 0
	llm := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		llmCalls++
		w.Header().Set("content-type", "application/json")
		if llmCalls == 1 {
			_, _ = w.Write([]byte(`{"id":"r1","status":"completed","output":[{"type":"function_call","call_id":"c1","name":"create_quote","arguments":"{\"product\":\"scooter\",\"quantity\":1}"}],"usage":{"input_tokens":1,"output_tokens":1,"total_tokens":2}}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"r2","status":"completed","output":[{"type":"message","content":[{"type":"output_text","text":"Cotización creada."}]}],"usage":{"input_tokens":1,"output_tokens":1,"total_tokens":2}}`))
	}))
	defer llm.Close()
	domainCalls := 0
	domain := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		domainCalls++
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"quotation_id":"q1"}`))
	}))
	defer domain.Close()
	receiver := &appChannel{}
	sender := &appReceiptSender{}
	durable := &outbounddelivery.Channel{CodeValue: "whatsapp", Receiver: receiver, Sender: sender, Store: deliveryStore}
	registry := channels.NewRegistry()
	if err = registry.Register(durable); err != nil {
		t.Fatal(err)
	}
	cfg := fullConfig()
	cfg.TenantID, cfg.LLMBaseURL, cfg.DomainBaseURL = tenant, llm.URL, domain.URL
	cfg.ChannelRegistry, cfg.ConversationStore, cfg.ContactResolver = registry, turnStore, identityStore
	a, err := NewProduction(cfg)
	if err != nil {
		t.Fatal(err)
	}
	message := channels.Message{ChannelCode: "whatsapp", TenantID: tenant, ExternalID: "wa-durable", ThreadID: "thread", ProviderMessageID: "wamid.inbound.1", OccurredAt: when, Direction: channels.DirectionIn, Text: "Cotizame un scooter"}
	if err = a.Dispatcher.Handle(ctx, message); err != nil {
		t.Fatal(err)
	}
	if err = a.Dispatcher.Handle(ctx, message); err != nil {
		t.Fatal(err)
	}
	var state string
	var attempts, events int
	if err = pool.QueryRow(ctx, `select d.state,d.attempt_count,count(e.*) from communication.outbound_delivery d join communication.outbound_delivery_event e using(tenant_id,channel_code,delivery_key) where d.tenant_id=$1 group by d.state,d.attempt_count`, tenant).Scan(&state, &attempts, &events); err != nil {
		t.Fatal(err)
	}
	if state != "accepted" || attempts != 1 || events != 2 || sender.calls != 1 || llmCalls != 2 || domainCalls != 1 {
		t.Fatalf("state=%s attempts=%d events=%d sends=%d llm=%d domain=%d", state, attempts, events, sender.calls, llmCalls, domainCalls)
	}
}

func TestNewProductionRejectsUnfencedChannel(t *testing.T) {
	cfg := fullConfig()
	if _, err := NewProduction(cfg); err == nil {
		t.Fatal("unfenced channel accepted")
	}
}
````

### FILE: `db/migrations/0048_outbound_delivery_fence.up.sql`
```yaml
block_id: "GO-PG-OUTBOUND-DELIVERY-FENCE:db/migrations/0048_outbound_delivery_fence.up.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "f7c69aa2f005653bf1fbfe8d3097e2bde0b96f72b494b5011c5753f1a29db9f9"
variables: []
secrets_allowed: false
```
````sql
begin;

create table communication.outbound_delivery (
  tenant_id uuid not null references platform.tenant(tenant_id),
  channel_code text not null check (channel_code ~ '^[a-z][a-z0-9_]{0,31}$'),
  delivery_key text not null check (length(delivery_key) between 16 and 128),
  request_sha256_hex text not null check (request_sha256_hex ~ '^[0-9a-f]{64}$'),
  recipient_hmac text not null check (recipient_hmac ~ '^[0-9a-f]{64}$'),
  state text not null check (state in ('sending','accepted','unknown','failed_terminal')),
  attempt_count integer not null check (attempt_count=1),
  locked_until timestamptz,
  provider_message_hmac text check (provider_message_hmac is null or provider_message_hmac ~ '^[0-9a-f]{64}$'),
  evidence_sha256_hex text check (evidence_sha256_hex is null or evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  failure_code text,
  accepted_at timestamptz,
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, channel_code, delivery_key),
  check ((state='sending' and locked_until is not null) or state<>'sending'),
  check ((state='accepted' and provider_message_hmac is not null and evidence_sha256_hex is not null and accepted_at is not null) or state<>'accepted'),
  check (failure_code is null or failure_code ~ '^[A-Z][A-Z0-9_]{0,127}$')
);

create table communication.outbound_delivery_event (
  tenant_id uuid not null,
  channel_code text not null,
  delivery_key text not null,
  sequence bigint not null check (sequence > 0),
  state text not null check (state in ('sending','accepted','unknown','failed_terminal')),
  evidence_sha256_hex text check (evidence_sha256_hex is null or evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  failure_code text,
  recorded_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,channel_code,delivery_key,sequence),
  foreign key (tenant_id,channel_code,delivery_key) references communication.outbound_delivery(tenant_id,channel_code,delivery_key),
  check (failure_code is null or failure_code ~ '^[A-Z][A-Z0-9_]{0,127}$')
);

create function communication.reject_outbound_delivery_event_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='outbound delivery event is immutable';
end;
$function$;

create trigger outbound_delivery_event_immutable before update or delete on communication.outbound_delivery_event
for each row execute function communication.reject_outbound_delivery_event_mutation();

create index outbound_delivery_unknown_idx on communication.outbound_delivery(updated_at) where state='unknown';

commit;
````

### FILE: `db/migrations/0048_outbound_delivery_fence.down.sql`
```yaml
block_id: "GO-PG-OUTBOUND-DELIVERY-FENCE:db/migrations/0048_outbound_delivery_fence.down.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "65c2f7f8110261080274cbe88bbb765744d85365aa5e32e4d1810e808dbdef16"
variables: []
secrets_allowed: false
```
````sql
begin;
drop table communication.outbound_delivery_event;
drop table communication.outbound_delivery;
drop function communication.reject_outbound_delivery_event_mutation();
commit;
````

### FILE: `db/tests/0048_outbound_delivery_fence.test.sql`
```yaml
block_id: "GO-PG-OUTBOUND-DELIVERY-FENCE:db/tests/0048_outbound_delivery_fence.test.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c9211698f9ba69eabed24a172679cc3e358f5a802d6d5609bc6a9185a7cbe666"
variables: []
secrets_allowed: false
```
````sql
begin;
insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values('48484848-4848-4848-8848-484848484848','outbound-test','Outbound Test','Outbound Test');
insert into communication.outbound_delivery
(tenant_id,channel_code,delivery_key,request_sha256_hex,recipient_hmac,state,attempt_count,locked_until)
values('48484848-4848-4848-8848-484848484848','whatsapp',repeat('d',64),repeat('a',64),repeat('b',64),'sending',1,clock_timestamp()+interval '1 minute');
insert into communication.outbound_delivery_event(tenant_id,channel_code,delivery_key,sequence,state)
values('48484848-4848-4848-8848-484848484848','whatsapp',repeat('d',64),1,'sending');
do $test$
begin
  begin
    update communication.outbound_delivery_event set state='accepted'
    where tenant_id='48484848-4848-4848-8848-484848484848';
    raise exception 'delivery event mutation unexpectedly succeeded';
  exception when sqlstate '55000' then null;
  end;
  if exists(select 1 from communication.outbound_delivery where attempt_count<>1) then
    raise exception 'automatic retry became possible';
  end if;
end;
$test$;
rollback;
````

## 6. Configuration surface

- HMAC key de al menos 32 bytes desde secret manager; puede compartir ciclo gobernado con contact identity sólo si la política del proyecto lo aprueba.
- Lease entre 1 segundo y 10 minutos, mayor al timeout máximo del sender y probado con fault injection.
- El sender debe devolver `ProviderMessageID`, `EvidenceSHA256` y `AcceptedAt`; ningún campo se sintetiza.
- Los estados `unknown` exigen job/runbook de reconciliación específico del proveedor.

## 7. Dependency bill

- Go stdlib.
- `github.com/jackc/pgx/v5 v5.10.0`, ya fijado por backend.
- PostgreSQL 18.6; depende de tenant, communication y del pack de identidad para HMAC.
- El provider sender live y sus términos/licencia permanecen en su pack específico.

## 8. Apply order

Aplicar `0001`–`0047`, luego `0048`. En producción envolver cada adapter registrado con `outbounddelivery.Channel` y construir la app mediante `app.NewProduction`. No registrar un sender directo como sustituto.

## 9. Verification

1. Materializar once archivos y comprobar hash/gofmt.
2. Aplicar 48 migraciones limpias y el SQL focal 0048.
3. Ejecutar suites de outbound, channels, app y PostgreSQL dos veces.
4. Ejecutar full Go test/vet/build.
5. Inyectar timeout posterior al write y demostrar una sola llamada, estado unknown y cierre sólo por reconciliación.
6. Provider live: probar receipt/status webhook/polling, rate/costo y no duplicación en sandbox; no se hereda del PASS local.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_PG_OUTBOUND_DELIVERY_FENCE_2026-09-04_V238.md`. El pack no contiene código de producto AWS/Meta/Google ni demuestra entrega live.

V402 composed delta: Payment/initial-handover composition. New behavior and tests belong to the explicit runtime/portal packs; source and library release claims remain bounded to their evidence. Existing source provenance is preserved.

Canonical V402 integration: selected by the current profile with exact dependencies and caller overlays. Metadata promotion records byte reconstruction, not closure of every admission/release gate. Payload provenance is unchanged.
