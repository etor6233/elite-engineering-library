# Go Commerce, Pricing and Payment API

## 1. Metadata

```yaml
pack_id: "GO-COMMERCE-PRICING-PAYMENT-API"
pack_version: "0.8.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Añade price books, líneas de pedido con precio server-side, colocación, asignación de stock ligada a la reserva canónica y payment attempts durables."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "pgx 5.10.0"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "GO-ELECTROMOBILITY-PUBLIC-CRM-API 0.2.x", "GO-SUPPLY-FACTORY-INVENTORY-API 0.12.x", "GO-BC-EXACT-AMOUNT-ADAPTER 0.1.0", "GO-BC-SALES-CONTRACT-ADAPTER 0.1.0"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner AND MIT dependency"
upstream_sources: ["https://go.dev", "https://www.postgresql.org/docs/18/", "https://github.com/jackc/pgx"]
verified_at: "2026-09-07"
```

Todos los bloques son `AUTHORED`. El pack crea una intención de pago y evento durable; no captura dinero sin un adapter oficial y contract tests del proveedor.

V300 / 0.6.1: la proyección Operations rechaza handovers cuyos vínculos con
pedido/cliente/stock no coinciden dentro de su snapshot Repeatable Read existente.
No devuelve resultados parciales ni modifica la reserva, el pago o la entrega.
Evidencia y límites: reconstruction_evidence/DELIVERY_READ_RECOVERY_V300.md.

V296 añade GET /v1/commerce/orders al owner actual: snapshot PostgreSQL read-only
Repeatable Read de pedidos/líneas, stock disponible, estados de pago y entrega.
Sin PII de cliente ni referencias provider. Exige organización y permiso operativo;
límites 25 pedidos/100 detalles/200 unidades explícitos con truncated. La UI bloquea
reservas con snapshot parcial; todavía no sustituye paginación para grandes volúmenes.
No cambia la transacción de reserva ni autoriza cobros/entregas por lectura.
Evidencia y límites: reconstruction_evidence/ORDER_STOCK_CONNECTED_V296.md.

V297 conserva un recibo idempotente local de pago: mismo tenant/proveedor/clave e
intención devuelve el ID/estado durable; cambios de intención devuelven 409.
El actor HTTP se guarda en el outbox original y una falla del outbox revierte
el intento. La migración 0053 permite varias referencias aún NULL sin permitir
referencias no nulas duplicadas. No introduce proveedor, variable, SDK ni cobro.
Fuente normativa: AWS Builders' Library, Making retries safe with idempotent APIs;
PostgreSQL 18 INSERT/ON CONFLICT y Unique Constraints. Implementación AUTHORED,
no código copiado ni aprobado por AWS/PostgreSQL. Ver PAYMENT_INTENT_RECOVERY_V297.md.

V402 / 0.6.2: calls GO-BC-EXACT-AMOUNT-ADAPTER0.1.0 for the precisely mapped amount/balance/reversal functions. Existing caller blocks remain AUTHORED; other business semantics are not reclassified. Source derivation, FAIL808 and exact local tests: reconstruction_evidence/BC_EXACT_AMOUNT_ADAPTATION_V402.md. Composition now requires that adapter.

V402 staging candidate / 0.6.3: delegates the precisely mapped price/quote transformations to GO-BC-SALES-CONTRACT-ADAPTER0.1.0. Existing caller provenance stays AUTHORED; scheduling, pricebook exclusivity, consent and financial state policy are not reclassified. No canonical publication is asserted by this candidate.

## 2. Applicability

V298 añade POST /v1/commerce/orders/{id}/payment-request en el mismo owner.
Sólo solicita el primer pago del total: organización/clave provienen del comando;
importe/moneda se derivan bajo lock del pedido. Otra clave o proveedor no permite
un segundo intento inicial. La exclusión vive en la transacción compartida y
se aplica también a la ruta previa; ningún caller del owner puede eludirla con
otra clave. Replay válido se conserva; múltiples cobros/reintentos financieros
necesitan política y workflow separado aprobado, no borrar el intento anterior.
CommerceModule.PaymentProvider vacío/inválido desactiva requests/transiciones HTTP;
stripe/mercadopago son selecciones explícitas, no credenciales ni adapters live.
GET publica esa selección para la UI. Error interno del nuevo recorder→503 sin
detalles privados; client input no puede elegir dinero/proveedor/actor.
Fuente AUTHORED, no upstream nuevo. Evidencia: PAYMENT_REQUEST_PORTAL_V298.md.

Use when the Go/PostgreSQL composition needs server-owned price books, order pricing, stock allocation and durable payment intent. Reject it when an external commerce platform is the source of truth or jurisdictional tax/accounting and the selected payment provider have not been modeled. It never claims to capture funds without an admitted adapter.

## 3. Architecture contract

Pricing is resolved server-side; clients cannot submit authoritative totals. Order, allocation, payment attempt and outbox changes share transactional boundaries and optimistic/unique guards. Tenant and organization scope comes from verified authorization. Provider effects occur asynchronously through idempotent adapters and reconciliation. Invalid transitions, stale versions and cross-scope access fail closed.

## 4. Exact file manifest

```text
CREATE internal/commerce/service.go
CREATE internal/commerce/service_test.go
CREATE internal/platform/postgres/commerce.go
CREATE internal/platform/postgres/commerce_integration_test.go
CREATE internal/platform/httpapi/commerce.go
CREATE internal/platform/httpapi/commerce_test.go
CREATE db/migrations/0053_payment_intent_reference.up.sql
CREATE db/migrations/0053_payment_intent_reference.down.sql
CREATE db/migrations/0066_order_funding_projection.up.sql
CREATE db/migrations/0066_order_funding_projection.down.sql
CREATE internal/platform/postgres/order_funding.go
CREATE internal/platform/postgres/local_funding.go
```

## 5. Materialization blocks

### FILE: `internal/commerce/service.go`

```yaml
block_id: "GO-COMMERCE-API:internal-commerce-service-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "32100fdc384284acef4db78d034121f86e5d290f48e88360245e6f9eab4cb983"
variables: []
secrets_allowed: false
```

````go
package commerce

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"
)

var ErrConflict = errors.New("commerce conflict")
var ErrPaymentUnavailable = errors.New("payment intent unavailable")
var currencyPattern = regexp.MustCompile(`^[A-Z]{3}$`)
var marketPattern = regexp.MustCompile(`^[A-Z]{2}$`)

type PriceEntry struct {
	VariantID        string `json:"variant_id"`
	AmountMinorUnits int64  `json:"amount_minor_units"`
	TaxMode          string `json:"tax_mode"`
}
type PriceBook struct {
	ID         string       `json:"id"`
	Market     string       `json:"market"`
	Currency   string       `json:"currency"`
	ValidFrom  time.Time    `json:"valid_from"`
	ValidUntil *time.Time   `json:"valid_until,omitempty"`
	Status     string       `json:"status"`
	Entries    []PriceEntry `json:"entries"`
}
type OrderLine struct {
	OrderID             string `json:"order_id"`
	LineID              string `json:"line_id"`
	OrganizationID      string `json:"organization_id"`
	PriceBookID         string `json:"price_book_id"`
	VariantID           string `json:"variant_id"`
	Quantity            int    `json:"quantity"`
	UnitPriceMinorUnits int64  `json:"unit_price_minor_units"`
	OrderVersion        int64  `json:"order_version"`
}
type PaymentAttempt struct {
	ID               string `json:"id"`
	OrderID          string `json:"order_id"`
	OrganizationID   string `json:"organization_id"`
	ProviderCode     string `json:"provider_code"`
	State            string `json:"state"`
	Currency         string `json:"currency"`
	AmountMinorUnits int64  `json:"amount_minor_units"`
	Version          int64  `json:"version"`
}
type Repository interface {
	CreatePriceBook(context.Context, string, string, PriceBook) error
	ActivatePriceBook(context.Context, string, string, string) error
	PublicPrice(context.Context, string, string, string) (PriceEntry, error)
	AddOrderLine(context.Context, string, string, OrderLine, int64) (OrderLine, error)
	PlaceOrder(context.Context, string, string, string, int64, string) error
	AllocateStock(context.Context, string, string, string, string, string, int64, int64, string) error
	CreatePaymentAttempt(context.Context, string, string, string, PaymentAttempt) error
	TransitionPayment(context.Context, string, string, string, string, string, int64, string, string) error
}
type IDGenerator interface{ New() string }

// OperationReader is the read side of the existing commerce owner. It never
// creates a second inventory, payment ledger or delivery state machine.
type OperationReader interface {
	Operations(context.Context, string, string) (OperationSnapshot, error)
}
type AuditedAllocator interface {
	AllocateStockAs(context.Context, string, string, string, string, string, int64, int64, string, string) error
}
type PaymentIntentRecorder interface {
	RecordPaymentIntent(context.Context, string, string, string, PaymentAttempt, string) (PaymentAttempt, error)
}
type OrderPaymentRecorder interface {
	RecordOrderPayment(context.Context, string, string, string, PaymentAttempt, string) (PaymentAttempt, error)
}
type OperationSnapshot struct {
	PaymentProvider string           `json:"payment_provider"`
	Orders          []OperationOrder `json:"orders"`
	Stock           []OperationStock `json:"stock"`
	Truncated       bool             `json:"truncated"`
}
type OperationOrder struct {
	ID        string            `json:"id"`
	State     string            `json:"state"`
	Version   int64             `json:"version"`
	Lines     []OperationLine   `json:"lines"`
	Payments  []OperationStatus `json:"payments"`
	Handovers []OperationStatus `json:"handovers"`
}
type OperationLine struct {
	ID        string `json:"id"`
	VariantID string `json:"variant_id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	StockID   string `json:"stock_id"`
}
type OperationStock struct {
	ID        string `json:"id"`
	VariantID string `json:"variant_id"`
	Serial    string `json:"serial"`
	Version   int64  `json:"version"`
}
type OperationStatus struct {
	ID    string `json:"id"`
	State string `json:"state"`
}
type Service struct {
	repository Repository
	ids        IDGenerator
}

func NewService(r Repository, ids IDGenerator) *Service { return &Service{repository: r, ids: ids} }
func (s *Service) Operations(ctx context.Context, tenant, organization string) (OperationSnapshot, error) {
	reader, ok := s.repository.(OperationReader)
	if !ok || tenant == "" || organization == "" {
		return OperationSnapshot{}, errors.New("operation reader unavailable or invalid scope")
	}
	return reader.Operations(ctx, tenant, organization)
}
func (s *Service) CreatePriceBook(ctx context.Context, tenant string, input PriceBook) (PriceBook, error) {
	if tenant == "" || !marketPattern.MatchString(input.Market) || !currencyPattern.MatchString(input.Currency) || input.ValidFrom.IsZero() || len(input.Entries) == 0 {
		return PriceBook{}, fmt.Errorf("invalid price book")
	}
	seen := map[string]bool{}
	for _, entry := range input.Entries {
		if entry.VariantID == "" || entry.AmountMinorUnits < 0 || !map[string]bool{"exclusive": true, "inclusive": true, "not-applicable": true}[entry.TaxMode] || seen[entry.VariantID] {
			return PriceBook{}, fmt.Errorf("invalid price entry")
		}
		seen[entry.VariantID] = true
	}
	input.ID = s.ids.New()
	input.Status = "draft"
	if err := s.repository.CreatePriceBook(ctx, tenant, s.ids.New(), input); err != nil {
		return PriceBook{}, err
	}
	return input, nil
}
func (s *Service) ActivatePriceBook(ctx context.Context, tenant, id string) error {
	if tenant == "" || id == "" {
		return fmt.Errorf("invalid price book")
	}
	return s.repository.ActivatePriceBook(ctx, tenant, id, s.ids.New())
}
func (s *Service) PublicPrice(ctx context.Context, tenantCode, market, variant string) (PriceEntry, error) {
	if tenantCode == "" || !marketPattern.MatchString(market) || variant == "" {
		return PriceEntry{}, fmt.Errorf("invalid price query")
	}
	return s.repository.PublicPrice(ctx, tenantCode, market, variant)
}
func (s *Service) AddOrderLine(ctx context.Context, tenant string, input OrderLine, expectedVersion int64) (OrderLine, error) {
	if tenant == "" || input.OrderID == "" || input.OrganizationID == "" || input.PriceBookID == "" || input.VariantID == "" || input.Quantity < 1 || expectedVersion < 1 {
		return OrderLine{}, fmt.Errorf("invalid order line")
	}
	input.LineID = s.ids.New()
	return s.repository.AddOrderLine(ctx, tenant, s.ids.New(), input, expectedVersion)
}
func (s *Service) PlaceOrder(ctx context.Context, tenant, organization, orderID string, version int64) error {
	if tenant == "" || organization == "" || orderID == "" || version < 1 {
		return fmt.Errorf("invalid order placement")
	}
	return s.repository.PlaceOrder(ctx, tenant, organization, orderID, version, s.ids.New())
}
func (s *Service) AllocateStock(ctx context.Context, tenant, organization, orderID, lineID, stockID string, orderVersion, stockVersion int64) error {
	if tenant == "" || organization == "" || orderID == "" || lineID == "" || stockID == "" || orderVersion < 1 || stockVersion < 1 {
		return fmt.Errorf("invalid allocation")
	}
	return s.repository.AllocateStock(ctx, tenant, organization, orderID, lineID, stockID, orderVersion, stockVersion, s.ids.New())
}

func (s *Service) AllocateStockAs(ctx context.Context, tenant, organization, orderID, lineID, stockID string, orderVersion, stockVersion int64, actor string) error {
	repo, ok := s.repository.(AuditedAllocator)
	if !ok || actor == "" || tenant == "" || organization == "" || orderID == "" || lineID == "" || stockID == "" || orderVersion < 1 || stockVersion < 1 {
		return fmt.Errorf("invalid audited allocation")
	}
	return repo.AllocateStockAs(ctx, tenant, organization, orderID, lineID, stockID, orderVersion, stockVersion, s.ids.New(), actor)
}

var paymentTransitions = map[string]map[string]bool{"created": {"pending": true, "authorized": true, "failed": true}, "pending": {"authorized": true, "failed": true}, "authorized": {"captured": true, "failed": true}, "captured": {"refunded": true, "disputed": true}, "disputed": {"refunded": true}}

// RequestOrderPayment is the initial whole-order intent only. It does not
// authorize a provider call, another attempt, partial payment or credit policy.
func (s *Service) RequestOrderPayment(ctx context.Context, tenant, organization, orderID, provider, key, actor string) (PaymentAttempt, error) {
	repo, ok := s.repository.(OrderPaymentRecorder)
	if !ok {
		return PaymentAttempt{}, ErrPaymentUnavailable
	}
	if tenant == "" || organization == "" || orderID == "" || actor == "" || (provider != "stripe" && provider != "mercadopago") || len(key) < 16 || len(key) > 128 {
		return PaymentAttempt{}, fmt.Errorf("invalid order payment request")
	}
	value, err := repo.RecordOrderPayment(ctx, tenant, s.ids.New(), key, PaymentAttempt{ID: s.ids.New(), OrderID: orderID, OrganizationID: organization, ProviderCode: provider}, actor)
	if err != nil && !errors.Is(err, ErrConflict) {
		return PaymentAttempt{}, ErrPaymentUnavailable
	}
	return value, err
}

func (s *Service) CreatePaymentAttempt(ctx context.Context, tenant, idempotency string, input PaymentAttempt) (PaymentAttempt, error) {
	return s.createPaymentAttempt(ctx, tenant, idempotency, input, "")
}
func (s *Service) CreatePaymentAttemptAs(ctx context.Context, tenant, idempotency string, input PaymentAttempt, actor string) (PaymentAttempt, error) {
	if actor == "" {
		return PaymentAttempt{}, fmt.Errorf("authenticated payment actor required")
	}
	return s.createPaymentAttempt(ctx, tenant, idempotency, input, actor)
}
func (s *Service) createPaymentAttempt(ctx context.Context, tenant, idempotency string, input PaymentAttempt, actor string) (PaymentAttempt, error) {
	if tenant == "" || len(idempotency) < 16 || len(idempotency) > 128 || input.OrderID == "" || input.OrganizationID == "" || input.ProviderCode == "" || !currencyPattern.MatchString(input.Currency) || input.AmountMinorUnits <= 0 {
		return PaymentAttempt{}, fmt.Errorf("invalid payment attempt")
	}
	recorder, ok := s.repository.(PaymentIntentRecorder)
	if !ok {
		return PaymentAttempt{}, fmt.Errorf("durable payment recorder unavailable")
	}
	input.ID = s.ids.New()
	input.State = "created"
	input.Version = 1
	return recorder.RecordPaymentIntent(ctx, tenant, s.ids.New(), idempotency, input, actor)
}
func (s *Service) TransitionPayment(ctx context.Context, tenant, organization, id, current, target string, version int64, providerReference string) error {
	if organization == "" || !paymentTransitions[current][target] || version < 1 {
		return fmt.Errorf("%w: invalid payment transition", ErrConflict)
	}
	return s.repository.TransitionPayment(ctx, tenant, organization, id, current, target, version, providerReference, s.ids.New())
}
````

### FILE: `internal/commerce/service_test.go`

```yaml
block_id: "GO-COMMERCE-API:internal-commerce-service_test-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "ba1911f4d2db42c14ef6c46693bcaf7d674db693ddf3af156e1f306d3bc89d60"
variables: []
secrets_allowed: false
```

````go
package commerce

import (
	"context"
	"testing"
	"time"
)

type fakeRepo struct{ calls int }

func (f *fakeRepo) CreatePriceBook(context.Context, string, string, PriceBook) error {
	f.calls++
	return nil
}
func (f *fakeRepo) ActivatePriceBook(context.Context, string, string, string) error { return nil }
func (f *fakeRepo) PublicPrice(context.Context, string, string, string) (PriceEntry, error) {
	return PriceEntry{}, nil
}
func (f *fakeRepo) AddOrderLine(_ context.Context, _ string, _ string, v OrderLine, version int64) (OrderLine, error) {
	v.OrderVersion = version + 1
	return v, nil
}
func (f *fakeRepo) PlaceOrder(context.Context, string, string, string, int64, string) error {
	return nil
}
func (f *fakeRepo) AllocateStock(context.Context, string, string, string, string, string, int64, int64, string) error {
	return nil
}
func (f *fakeRepo) CreatePaymentAttempt(context.Context, string, string, string, PaymentAttempt) error {
	return nil
}
func (f *fakeRepo) TransitionPayment(context.Context, string, string, string, string, string, int64, string, string) error {
	return nil
}

type ids struct{ n int }

func (i *ids) New() string {
	i.n++
	return []string{"018f4d4a-7b36-7a21-8d10-2f4c54c28301", "018f4d4a-7b36-7a21-8d10-2f4c54c28302", "018f4d4a-7b36-7a21-8d10-2f4c54c28303", "018f4d4a-7b36-7a21-8d10-2f4c54c28304"}[i.n-1]
}
func TestCommerceRejectsInvalidPolicies(t *testing.T) {
	service := NewService(&fakeRepo{}, &ids{})
	ctx := context.Background()
	_, err := service.CreatePriceBook(ctx, "tenant", PriceBook{Market: "AR", Currency: "ARS", ValidFrom: time.Now(), Entries: []PriceEntry{{VariantID: "v", AmountMinorUnits: 1, TaxMode: "inclusive"}, {VariantID: "v", AmountMinorUnits: 2, TaxMode: "inclusive"}}})
	if err == nil {
		t.Fatal("duplicate variant accepted")
	}
	if err := service.TransitionPayment(ctx, "tenant", "o", "p", "created", "captured", 1, ""); err == nil {
		t.Fatal("payment skipped authorization")
	}
	if err := service.AllocateStock(ctx, "tenant", "org", "o", "l", "s", 0, 1); err == nil {
		t.Fatal("allocation without order version")
	}
}

func TestPaymentRecorderAndActorFailClosed(t *testing.T) {
	s := NewService(&fakeRepo{}, &ids{})
	in := PaymentAttempt{OrderID: "order", OrganizationID: "org", ProviderCode: "synthetic", Currency: "ARS", AmountMinorUnits: 100}
	if _, err := s.CreatePaymentAttempt(context.Background(), "tenant", "valid-payment-key", in); err == nil {
		t.Fatal("legacy-only repository silently used")
	}
	if _, err := s.CreatePaymentAttemptAs(context.Background(), "tenant", "valid-payment-key", in, ""); err == nil {
		t.Fatal("missing authenticated actor accepted")
	}
}
````

### FILE: `internal/platform/postgres/commerce.go`

```yaml
block_id: "GO-COMMERCE-API:internal-platform-postgres-commerce-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "a07be25f185586d6472b58901ef87ef9a2557548c6bfca74ea056385dc22f9a9"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"elite.local/enterprise/internal/bcamounts"
	"elite.local/enterprise/internal/businesspolicy"
	"elite.local/enterprise/internal/commerce"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Commerce struct {
	pool   *pgxpool.Pool
	policy *businesspolicy.Profile
}

func NewCommerce(pool *pgxpool.Pool) *Commerce {
	return &Commerce{pool: pool, policy: businesspolicy.Reference()}
}

// The v1 profile admits the existing single-active-source contract only.
// A mode that would need price ranking is rejected by businesspolicy.Load.
func NewCommerceWithProfile(pool *pgxpool.Pool, policy *businesspolicy.Profile) (*Commerce, error) {
	if pool == nil || !policy.Valid() {
		return nil, businesspolicy.ErrProfile
	}
	return &Commerce{pool: pool, policy: policy}, nil
}
func (r *Commerce) BusinessPolicySHA256() string {
	if r == nil {
		return ""
	}
	return r.policy.SHA256()
}

func (r *Commerce) Operations(ctx context.Context, tenant, organization string) (commerce.OperationSnapshot, error) {
	value := commerce.OperationSnapshot{Orders: []commerce.OperationOrder{}, Stock: []commerce.OperationStock{}}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(ctx, `select order_id,state,version from sales.customer_order where tenant_id=$1 and organization_id=$2 order by created_at desc,order_id limit 26`, tenant, organization)
	if err != nil {
		return value, err
	}
	for rows.Next() {
		var o commerce.OperationOrder
		if err = rows.Scan(&o.ID, &o.State, &o.Version); err != nil {
			rows.Close()
			return value, err
		}
		o.Lines = []commerce.OperationLine{}
		o.Payments = []commerce.OperationStatus{}
		o.Handovers = []commerce.OperationStatus{}
		value.Orders = append(value.Orders, o)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return value, err
	}
	if len(value.Orders) > 25 {
		value.Truncated = true
		value.Orders = value.Orders[:25]
	}
	for i := range value.Orders {
		o := &value.Orders[i]
		rows, err = tx.Query(ctx, `select l.line_id,l.variant_id,v.display_name,l.quantity,coalesce(l.allocated_stock_unit_id,'') from sales.customer_order_line l join catalog.vehicle_variant v on v.tenant_id=l.tenant_id and v.variant_id=l.variant_id where l.tenant_id=$1 and l.order_id=$2 order by l.line_id limit 101`, tenant, o.ID)
		if err != nil {
			return value, err
		}
		for rows.Next() {
			var l commerce.OperationLine
			if err = rows.Scan(&l.ID, &l.VariantID, &l.Name, &l.Quantity, &l.StockID); err != nil {
				rows.Close()
				return value, err
			}
			o.Lines = append(o.Lines, l)
		}
		rows.Close()
		if err = rows.Err(); err != nil {
			return value, err
		}
		if len(o.Lines) > 100 {
			value.Truncated = true
			o.Lines = o.Lines[:100]
		}
		for _, query := range []struct {
			sql    string
			target *[]commerce.OperationStatus
		}{
			{`select payment_attempt_id,state from payment.payment_attempt where tenant_id=$1 and order_id=$2 order by payment_attempt_id limit 101`, &o.Payments},
			{`select h.handover_id,h.state,exists(select 1 from sales.customer_order co join inventory.stock_unit su on su.tenant_id=co.tenant_id and su.stock_unit_id=h.stock_unit_id and su.organization_id=h.organization_id where co.tenant_id=h.tenant_id and co.order_id=h.order_id and co.organization_id=h.organization_id and co.customer_principal_id=h.customer_principal_id) from sales.delivery_handover h where h.tenant_id=$1 and h.order_id=$2 and h.organization_id=$3 order by h.handover_id limit 101`, &o.Handovers},
		} {
			args := []any{tenant, o.ID}
			if query.target == &o.Handovers {
				args = append(args, organization)
			}
			rows, err = tx.Query(ctx, query.sql, args...)
			if err != nil {
				return value, err
			}
			for rows.Next() {
				var s commerce.OperationStatus
				coherent := true
				columns := []any{&s.ID, &s.State}
				if query.target == &o.Handovers {
					columns = append(columns, &coherent)
				}
				if err = rows.Scan(columns...); err != nil {
					rows.Close()
					return value, err
				}
				if !coherent {
					rows.Close()
					return commerce.OperationSnapshot{}, errors.New("operation handover scope inconsistent")
				}
				*query.target = append(*query.target, s)
			}
			rows.Close()
			if err = rows.Err(); err != nil {
				return value, err
			}
			if len(*query.target) > 100 {
				value.Truncated = true
				*query.target = (*query.target)[:100]
			}
		}
	}
	rows, err = tx.Query(ctx, `select stock_unit_id,variant_id,serial_number,version from inventory.stock_unit where tenant_id=$1 and organization_id=$2 and state='available' order by stock_unit_id limit 201`, tenant, organization)
	if err != nil {
		return value, err
	}
	for rows.Next() {
		var s commerce.OperationStock
		if err = rows.Scan(&s.ID, &s.VariantID, &s.Serial, &s.Version); err != nil {
			rows.Close()
			return value, err
		}
		value.Stock = append(value.Stock, s)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return value, err
	}
	if len(value.Stock) > 200 {
		value.Truncated = true
		value.Stock = value.Stock[:200]
	}
	return value, tx.Commit(ctx)
}
func (r *Commerce) CreatePriceBook(ctx context.Context, tenant, eventID string, value commerce.PriceBook) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = r.createPriceBookTx(ctx, tx, tenant, eventID, value); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Commerce) ActivatePriceBook(ctx context.Context, tenant, id, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = r.activatePriceBookTx(ctx, tx, tenant, id, eventID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Commerce) PublicPrice(ctx context.Context, tenantCode, market, variant string) (commerce.PriceEntry, error) {
	rows, err := r.pool.Query(ctx, `select e.variant_id,e.amount_minor_units,e.tax_mode from pricing.price_book b join pricing.price_book_entry e on e.tenant_id=b.tenant_id and e.price_book_id=b.price_book_id join platform.tenant t on t.tenant_id=b.tenant_id `+bcPriceTimeContextSQL+` where t.tenant_code=$1 and t.status='active' and b.market=$2 and `+bcPriceEligibilitySQL+` and e.variant_id=$3 limit 2`, tenantCode, market, variant)
	if err != nil {
		return commerce.PriceEntry{}, err
	}
	defer rows.Close()
	values := []commerce.PriceEntry{}
	for rows.Next() {
		var value commerce.PriceEntry
		if err := rows.Scan(&value.VariantID, &value.AmountMinorUnits, &value.TaxMode); err != nil {
			return value, err
		}
		values = append(values, value)
	}
	if len(values) != 1 {
		return commerce.PriceEntry{}, commerce.ErrConflict
	}
	return values[0], rows.Err()
}
func (r *Commerce) AddOrderLine(ctx context.Context, tenant, eventID string, line commerce.OrderLine, expectedVersion int64) (commerce.OrderLine, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return line, err
	}
	defer tx.Rollback(ctx)
	var currency string
	var total int64
	err = tx.QueryRow(ctx, `select currency,total_minor_units from sales.customer_order where tenant_id=$1 and order_id=$2 and organization_id=$3 and state='draft' and version=$4 for update`, tenant, line.OrderID, line.OrganizationID, expectedVersion).Scan(&currency, &total)
	if errors.Is(err, pgx.ErrNoRows) {
		return line, commerce.ErrConflict
	}
	if err != nil {
		return line, err
	}
	var amount int64
	err = tx.QueryRow(ctx, `select e.amount_minor_units from pricing.price_book b join pricing.price_book_entry e on e.tenant_id=b.tenant_id and e.price_book_id=b.price_book_id `+bcPriceTimeContextSQL+` where b.tenant_id=$1 and b.price_book_id=$2 and `+bcPriceCurrencyBindingSQL+` and `+bcPriceEligibilitySQL+` and e.variant_id=$4`, tenant, line.PriceBookID, currency, line.VariantID).Scan(&amount)
	if errors.Is(err, pgx.ErrNoRows) {
		return line, commerce.ErrConflict
	}
	if err != nil {
		return line, err
	}
	var current int64
	if err = tx.QueryRow(ctx, `select coalesce(sum(quantity*unit_price_minor_units),0) from sales.customer_order_line where tenant_id=$1 and order_id=$2`, tenant, line.OrderID).Scan(&current); err != nil {
		return line, err
	}
	lineAmount, amountErr := bcamounts.LineAmount(int64(line.Quantity), amount, 0)
	if amountErr != nil || lineAmount <= 0 || current < 0 || total < 0 || lineAmount > total || current > total-lineAmount {
		return line, commerce.ErrConflict
	}
	line.UnitPriceMinorUnits = amount
	line.OrderVersion = expectedVersion + 1
	_, err = tx.Exec(ctx, `insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units)values($1,$2,$3,$4,$5,$6)`, tenant, line.OrderID, line.LineID, line.VariantID, line.Quantity, amount)
	if err != nil {
		return line, err
	}
	result, err := tx.Exec(ctx, `update sales.customer_order set version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and order_id=$2 and version=$3`, tenant, line.OrderID, expectedVersion)
	if err != nil {
		return line, err
	}
	if result.RowsAffected() != 1 {
		return line, commerce.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'order',$3,$4,'order.line-added',1,clock_timestamp(),jsonb_build_object('line_id',$5::text,'variant_id',$6::text))`, tenant, eventID, line.OrderID, line.OrderVersion, line.LineID, line.VariantID)
	if err != nil {
		return line, err
	}
	return line, tx.Commit(ctx)
}
func (r *Commerce) PlaceOrder(ctx context.Context, tenant, organization, orderID string, version int64, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var total int64
	err = tx.QueryRow(ctx, `select total_minor_units from sales.customer_order where tenant_id=$1 and organization_id=$2 and order_id=$3 and state='draft' and version=$4 for update`, tenant, organization, orderID, version).Scan(&total)
	if errors.Is(err, pgx.ErrNoRows) {
		return commerce.ErrConflict
	}
	if err != nil {
		return err
	}
	var lines int64
	if err = tx.QueryRow(ctx, `select coalesce(sum(quantity*unit_price_minor_units),0) from sales.customer_order_line where tenant_id=$1 and order_id=$2`, tenant, orderID).Scan(&lines); err != nil {
		return err
	}
	if total != lines || lines <= 0 {
		return commerce.ErrConflict
	}
	result, err := tx.Exec(ctx, `update sales.customer_order set state='placed',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and order_id=$2 and version=$3`, tenant, orderID, version)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return commerce.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'order',$3,$4,'order.placed',1,clock_timestamp(),'{}')`, tenant, eventID, orderID, version+1)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Commerce) AllocateStock(ctx context.Context, tenant, authorizedOrganization, orderID, lineID, stockID string, orderVersion, stockVersion int64, eventID string) error {
	return r.allocateStock(ctx, tenant, authorizedOrganization, orderID, lineID, stockID, orderVersion, stockVersion, eventID, "")
}
func (r *Commerce) AllocateStockAs(ctx context.Context, tenant, organization, orderID, lineID, stockID string, orderVersion, stockVersion int64, eventID, actor string) error {
	if actor == "" {
		return commerce.ErrConflict
	}
	return r.allocateStock(ctx, tenant, organization, orderID, lineID, stockID, orderVersion, stockVersion, eventID, actor)
}
func (r *Commerce) allocateStock(ctx context.Context, tenant, authorizedOrganization, orderID, lineID, stockID string, orderVersion, stockVersion int64, eventID, actor string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var variant, organization string
	err = tx.QueryRow(ctx, `select l.variant_id,o.organization_id from sales.customer_order_line l join sales.customer_order o on o.tenant_id=l.tenant_id and o.order_id=l.order_id where l.tenant_id=$1 and o.organization_id=$2 and l.order_id=$3 and l.line_id=$4 and l.allocated_stock_unit_id is null and l.quantity=1 and o.state in ('placed','confirmed') and o.version=$5 for update of o,l`, tenant, authorizedOrganization, orderID, lineID, orderVersion).Scan(&variant, &organization)
	if errors.Is(err, pgx.ErrNoRows) {
		return commerce.ErrConflict
	}
	if err != nil {
		return err
	}
	result, err := tx.Exec(ctx, `update inventory.stock_unit set state='reserved',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and stock_unit_id=$2 and organization_id=$3 and variant_id=$4 and state='available' and version=$5`, tenant, stockID, organization, variant, stockVersion)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return commerce.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into inventory.serial_reservation(tenant_id,reservation_id,organization_id,stock_unit_id,variant_id,demand_kind,demand_id,demand_line_id,status,cancellation_disallowed,version) values($1,$2,$3,$4,$5,'customer-order',$6,$7,'reservation',true,1)`, tenant, eventID, organization, stockID, variant, orderID, lineID)
	if err != nil {
		return err
	}
	result, err = tx.Exec(ctx, `update sales.customer_order_line set allocated_stock_unit_id=$4 where tenant_id=$1 and order_id=$2 and line_id=$3 and allocated_stock_unit_id is null`, tenant, orderID, lineID, stockID)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return commerce.ErrConflict
	}
	_, err = tx.Exec(ctx, `update sales.customer_order set version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and order_id=$2 and version=$3`, tenant, orderID, orderVersion)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'order',$3,$4,'order.stock-allocated',1,clock_timestamp(),jsonb_build_object('line_id',$5::text,'stock_unit_id',$6::text) || case when $7::text<>'' then jsonb_build_object('actor_subject',$7::text) else '{}'::jsonb end)`, tenant, eventID, orderID, orderVersion+1, lineID, stockID, actor)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Commerce) CreatePaymentAttempt(ctx context.Context, tenant, eventID, idempotency string, value commerce.PaymentAttempt) error {
	_, err := r.recordPaymentIntent(ctx, tenant, eventID, idempotency, value, "", false, false)
	return err
}
func (r *Commerce) RecordPaymentIntent(ctx context.Context, tenant, eventID, idempotency string, value commerce.PaymentAttempt, actor string) (commerce.PaymentAttempt, error) {
	return r.recordPaymentIntent(ctx, tenant, eventID, idempotency, value, actor, true, false)
}
func (r *Commerce) RecordOrderPayment(ctx context.Context, tenant, eventID, idempotency string, value commerce.PaymentAttempt, actor string) (commerce.PaymentAttempt, error) {
	return r.recordPaymentIntent(ctx, tenant, eventID, idempotency, value, actor, true, true)
}
func (r *Commerce) recordPaymentIntent(ctx context.Context, tenant, eventID, idempotency string, value commerce.PaymentAttempt, actor string, allowReplay, deriveOrder bool) (commerce.PaymentAttempt, error) {
	empty := commerce.PaymentAttempt{}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return empty, err
	}
	defer tx.Rollback(ctx)
	var currency, orderState string
	var total int64
	// Lock the authorized order before reading its amount/state and writing the intent.
	// Existing receipts remain recoverable after a later order transition.
	err = tx.QueryRow(ctx, `select currency,total_minor_units,state from sales.customer_order where tenant_id=$1 and organization_id=$2 and order_id=$3 for update`, tenant, value.OrganizationID, value.OrderID).Scan(&currency, &total, &orderState)
	if errors.Is(err, pgx.ErrNoRows) {
		return empty, commerce.ErrConflict
	}
	if err != nil {
		return empty, err
	}
	total, err = orderProviderDue(ctx, tx, tenant, value.OrganizationID, value.OrderID, currency, total)
	if err != nil {
		return empty, err
	}
	if deriveOrder {
		value.Currency = currency
		value.AmountMinorUnits = total
	}
	readReceipt := func() (commerce.PaymentAttempt, error) {
		var saved commerce.PaymentAttempt
		err := tx.QueryRow(ctx, `select p.payment_attempt_id,p.order_id,o.organization_id,p.provider_code,p.state,p.currency,p.amount_minor_units,p.version from payment.payment_attempt p join sales.customer_order o on o.tenant_id=p.tenant_id and o.order_id=p.order_id where p.tenant_id=$1 and p.provider_code=$2 and p.idempotency_key=$3 and o.organization_id=$4`, tenant, value.ProviderCode, idempotency, value.OrganizationID).Scan(&saved.ID, &saved.OrderID, &saved.OrganizationID, &saved.ProviderCode, &saved.State, &saved.Currency, &saved.AmountMinorUnits, &saved.Version)
		if err != nil {
			return empty, err
		}
		if !allowReplay || saved.OrderID != value.OrderID || saved.Currency != value.Currency || saved.AmountMinorUnits != value.AmountMinorUnits {
			return empty, commerce.ErrConflict
		}
		return saved, nil
	}
	if saved, err := readReceipt(); err == nil {
		return saved, tx.Commit(ctx)
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return empty, err
	}
	// Every public creation path requests the remaining provider due. Enforce initial
	// intent exclusivity here, not just in the new frontend route.
	var exists bool
	if err := tx.QueryRow(ctx, `select exists(select 1 from payment.payment_attempt where tenant_id=$1 and order_id=$2)`, tenant, value.OrderID).Scan(&exists); err != nil {
		return empty, err
	}
	if exists {
		return empty, commerce.ErrConflict
	}
	if total <= 0 || currency != value.Currency || total != value.AmountMinorUnits || (orderState != "placed" && orderState != "confirmed") {
		return empty, commerce.ErrConflict
	}
	result, err := tx.Exec(ctx, `insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,idempotency_key,state,currency,amount_minor_units,version)values($1,$2,$3,$4,$5,'created',$6,$7,1) on conflict(tenant_id,provider_code,idempotency_key) do nothing`, tenant, value.ID, value.OrderID, value.ProviderCode, idempotency, value.Currency, value.AmountMinorUnits)
	if err != nil {
		return empty, err
	}
	if result.RowsAffected() != 1 {
		// Another order can contend for this provider/key. Read Committed gives
		// this subsequent statement visibility of the committed winning intent.
		saved, err := readReceipt()
		if errors.Is(err, pgx.ErrNoRows) {
			err = commerce.ErrConflict
		}
		if err != nil {
			return empty, err
		}
		return saved, tx.Commit(ctx)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'payment',$3,1,'payment.requested',1,clock_timestamp(),jsonb_build_object('provider_code',$4::text,'order_id',$5::text) || case when $6::text<>'' then jsonb_build_object('actor_subject',$6::text) else '{}'::jsonb end)`, tenant, eventID, value.ID, value.ProviderCode, value.OrderID, actor)
	if err != nil {
		return empty, err
	}
	value.State = "created"
	value.Version = 1
	return value, tx.Commit(ctx)
}
func (r *Commerce) TransitionPayment(ctx context.Context, tenant, organization, id, current, target string, version int64, providerReference, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update payment.payment_attempt p set state=$6,version=p.version+1,provider_reference=coalesce(nullif($7,''),p.provider_reference),updated_at=clock_timestamp() from sales.customer_order o where p.tenant_id=$1 and p.payment_attempt_id=$3 and p.state=$4 and p.version=$5 and o.tenant_id=p.tenant_id and o.order_id=p.order_id and o.organization_id=$2`, tenant, organization, id, current, version, target, providerReference)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return commerce.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'payment',$3,$4,$5,1,clock_timestamp(),'{}')`, tenant, eventID, id, version+1, "payment."+target)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED transaction extraction; original activation SQL retained.
func (r *Commerce) activatePriceBookTx(ctx context.Context, tx pgx.Tx, tenant, id, eventID string) error {
	var market, currency string
	var from time.Time
	var until *time.Time
	err := tx.QueryRow(ctx, `select market,currency,valid_from,valid_until from pricing.price_book where tenant_id=$1 and price_book_id=$2 and status='draft' for update`, tenant, id).Scan(&market, &currency, &from, &until)
	if errors.Is(err, pgx.ErrNoRows) {
		return commerce.ErrConflict
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":"+market+":"+currency)
	if err != nil {
		return err
	}
	var overlap bool
	err = tx.QueryRow(ctx, `select exists(select 1 from pricing.price_book where tenant_id=$1 and market=$2 and currency=$3 and status='active' and valid_from<coalesce($5::timestamptz,'infinity'::timestamptz) and coalesce(valid_until,'infinity'::timestamptz)>$4)`, tenant, market, currency, from, until).Scan(&overlap)
	if err != nil {
		return err
	}
	if overlap {
		return commerce.ErrConflict
	}
	_, err = tx.Exec(ctx, `update pricing.price_book set status='active' where tenant_id=$1 and price_book_id=$2 and status='draft'`, tenant, id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'price-book',$3,2,'price-book.activated',1,clock_timestamp(),jsonb_build_object('policy_sha256',$4::text))`, tenant, eventID, id, r.policy.SHA256())
	if err != nil {
		return err
	}
	return nil
}

// AUTHORED transaction extraction; original creation SQL/event retained.
func (r *Commerce) createPriceBookTx(ctx context.Context, tx pgx.Tx, tenant, eventID string, value commerce.PriceBook) error {
	_, err := tx.Exec(ctx, `insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,valid_until,status)values($1,$2,$3,$4,$5,$6,'draft')`, tenant, value.ID, value.Market, value.Currency, value.ValidFrom, value.ValidUntil)
	if err != nil {
		return err
	}
	for _, entry := range value.Entries {
		if _, err = tx.Exec(ctx, `insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)values($1,$2,$3,$4,$5)`, tenant, value.ID, entry.VariantID, entry.AmountMinorUnits, entry.TaxMode); err != nil {
			return err
		}
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'price-book',$3,1,'price-book.created',1,clock_timestamp(),jsonb_build_object('market',$4::text,'currency',$5::text))`, tenant, eventID, value.ID, value.Market, value.Currency)
	if err != nil {
		return err
	}
	return nil
}
````

### FILE: `internal/platform/postgres/commerce_integration_test.go`

```yaml
block_id: "GO-COMMERCE-API:internal-platform-postgres-commerce_integration_test-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "91e8c31cbf133f81044bcad51e7f382c5f21fff18303d99a1571635585663945"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"elite.local/enterprise/internal/commerce"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
	"time"
)

func TestCommercePriceOrderAllocationPaymentFlow(t *testing.T) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28401"
	cleanup := func() {
		for _, q := range []string{`delete from platform.outbox_event where tenant_id=$1`, `delete from payment.payment_attempt where tenant_id=$1`, `delete from inventory.serial_reservation where tenant_id=$1`, `delete from sales.customer_order_line where tenant_id=$1`, `delete from sales.customer_order where tenant_id=$1`, `delete from inventory.stock_unit where tenant_id=$1`, `delete from pricing.price_book_entry where tenant_id=$1`, `delete from pricing.price_book where tenant_id=$1`, `delete from catalog.vehicle_variant where tenant_id=$1`, `delete from catalog.vehicle_model where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`} {
			_, _ = pool.Exec(ctx, q, tenant)
		}
	}
	cleanup()
	defer cleanup()
	fixtures := []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'commerce-api','Commerce','Commerce')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'org','store','Store','store')`, `insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Model','bicycle','active')`, `insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Variant','{}','active')`, `insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','org','variant','STOCK-SERIAL','available',1,clock_timestamp())`, `insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','org','customer','draft','ARS',1000,1)`}
	for _, q := range fixtures {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewCommerce(pool)
	book := commerce.PriceBook{ID: "book", Market: "AR", Currency: "ARS", ValidFrom: time.Now().Add(-time.Hour), Status: "draft", Entries: []commerce.PriceEntry{{VariantID: "variant", AmountMinorUnits: 1000, TaxMode: "inclusive"}}}
	if err := repo.CreatePriceBook(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28402", book); err != nil {
		t.Fatal(err)
	}
	if err := repo.ActivatePriceBook(ctx, tenant, "book", "018f4d4a-7b36-7a21-8d10-2f4c54c28403"); err != nil {
		t.Fatal(err)
	}
	price, err := repo.PublicPrice(ctx, "commerce-api", "AR", "variant")
	if err != nil || price.AmountMinorUnits != 1000 {
		t.Fatalf("price=%+v err=%v", price, err)
	}
	line := commerce.OrderLine{OrderID: "order", LineID: "line", OrganizationID: "org", PriceBookID: "book", VariantID: "variant", Quantity: 1}
	line, err = repo.AddOrderLine(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28404", line, 1)
	if err != nil || line.UnitPriceMinorUnits != 1000 || line.OrderVersion != 2 {
		t.Fatalf("line=%+v err=%v", line, err)
	}
	if err := repo.PlaceOrder(ctx, tenant, "org", "order", 2, "018f4d4a-7b36-7a21-8d10-2f4c54c28405"); err != nil {
		t.Fatal(err)
	}
	if err := repo.AllocateStock(ctx, tenant, "org-other", "order", "line", "stock", 3, 1, "018f4d4a-7b36-7a21-8d10-2f4c54c28411"); err == nil {
		t.Fatal("allocation crossed organization scope")
	}
	if err := repo.AllocateStock(ctx, tenant, "org", "order", "line", "stock", 3, 1, "018f4d4a-7b36-7a21-8d10-2f4c54c28406"); err != nil {
		t.Fatal(err)
	}
	if err := repo.AllocateStock(ctx, tenant, "org", "order", "line", "stock", 3, 1, "018f4d4a-7b36-7a21-8d10-2f4c54c28407"); err == nil {
		t.Fatal("duplicate allocation succeeded")
	}
	payment := commerce.PaymentAttempt{ID: "payment", OrderID: "order", OrganizationID: "org", ProviderCode: "sandbox", Currency: "ARS", AmountMinorUnits: 1000}
	if err := repo.CreatePaymentAttempt(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28408", "payment-key-00000001", payment); err != nil {
		t.Fatal(err)
	}
	if err := repo.TransitionPayment(ctx, tenant, "org-other", "payment", "created", "authorized", 1, "provider-ref", "018f4d4a-7b36-7a21-8d10-2f4c54c28412"); err == nil {
		t.Fatal("payment crossed organization scope")
	}
	if err := repo.TransitionPayment(ctx, tenant, "org", "payment", "created", "authorized", 1, "provider-ref", "018f4d4a-7b36-7a21-8d10-2f4c54c28409"); err != nil {
		t.Fatal(err)
	}
	if err := repo.TransitionPayment(ctx, tenant, "org", "payment", "created", "failed", 1, "", "018f4d4a-7b36-7a21-8d10-2f4c54c28410"); err == nil {
		t.Fatal("stale payment transition succeeded")
	}
	var state string
	var version int64
	if err := pool.QueryRow(ctx, `select state,version from inventory.stock_unit where tenant_id=$1 and stock_unit_id='stock'`, tenant).Scan(&state, &version); err != nil || state != "reserved" || version != 2 {
		t.Fatalf("stock=%s version=%d err=%v", state, version, err)
	}
	var reservationStatus string
	if err := pool.QueryRow(ctx, `select status from inventory.serial_reservation where tenant_id=$1 and stock_unit_id='stock' and demand_kind='customer-order' and demand_id='order' and demand_line_id='line'`, tenant).Scan(&reservationStatus); err != nil || reservationStatus != "reservation" {
		t.Fatalf("reservation=%s err=%v", reservationStatus, err)
	}
	view, err := repo.Operations(ctx, tenant, "org")
	if err != nil || view.Truncated || len(view.Orders) != 1 || len(view.Stock) != 0 {
		t.Fatalf("operation snapshot mismatch err=%v", err)
	}
	if len(view.Orders[0].Lines) != 1 || view.Orders[0].Lines[0].StockID != "stock" || len(view.Orders[0].Payments) != 1 || view.Orders[0].Payments[0].State != "authorized" || len(view.Orders[0].Handovers) != 0 {
		t.Fatal("operation projection differs from durable owners")
	}
	for _, scope := range []struct{ tenant, org string }{{tenant, "other"}, {"018f4d4a-7b36-7a21-8d10-2f4c54c28999", "org"}} {
		other, err := repo.Operations(ctx, scope.tenant, scope.org)
		if err != nil || len(other.Orders) != 0 || len(other.Stock) != 0 {
			t.Fatal("operation query crosses scope")
		}
	}
	// Every exposed truncation is explicit; incomplete data cannot authorize a UI write.
	if _, err := pool.Exec(ctx, `insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) select $1,'limit-'||s::text,'org','customer','draft','ARS',1000,1 from generate_series(1,26)s`, tenant); err != nil {
		t.Fatal(err)
	}
	view, err = repo.Operations(ctx, tenant, "org")
	if err != nil || !view.Truncated || len(view.Orders) != 25 {
		t.Fatal("operation order limit not explicit")
	}
}
````

### FILE: `internal/platform/httpapi/commerce.go`

```yaml
block_id: "GO-COMMERCE-API:internal-platform-httpapi-commerce-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "5ba7a61bd62b41d03514b392062d02cdad514c7d017940f140add1779c430809"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"net/http"
)

type CommerceModule struct {
	Service                                *commerce.Service
	PaymentProvider                        string
	PaymentTenantID, PaymentOrganizationID string
	ProviderObservedPayments               bool
}

func (m CommerceModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	api := commerceAPI{service: m.Service, verifier: verifier, paymentProvider: m.PaymentProvider, paymentTenantID: m.PaymentTenantID, paymentOrganizationID: m.PaymentOrganizationID, providerObservedPayments: m.ProviderObservedPayments}
	mux.HandleFunc("GET /v1/commerce/orders", api.operations)
	mux.HandleFunc("GET /v1/public/{tenantCode}/prices", api.publicPrice)
	mux.HandleFunc("POST /v1/pricing/books", api.createBook)
	mux.HandleFunc("POST /v1/pricing/books/{id}/activate", api.activateBook)
	mux.HandleFunc("POST /v1/commerce/orders/{id}/lines", api.addLine)
	mux.HandleFunc("POST /v1/commerce/orders/{id}/place", api.placeOrder)
	mux.HandleFunc("POST /v1/commerce/orders/{id}/allocations", api.allocate)
	mux.HandleFunc("POST /v1/commerce/orders/{id}/payments", api.createPayment)
	mux.HandleFunc("POST /v1/commerce/orders/{id}/payment-request", api.requestOrderPayment)
	mux.HandleFunc("POST /v1/payments/{id}/transitions", api.transitionPayment)
}

type commerceAPI struct {
	service                                *commerce.Service
	verifier                               identity.Verifier
	paymentProvider                        string
	paymentTenantID, paymentOrganizationID string
	providerObservedPayments               bool
}

func (a commerceAPI) operations(w http.ResponseWriter, r *http.Request) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return
	}
	if !p.Allowed("inventory:allocate") && !p.Allowed("payment:create") && !p.Allowed("handover:manage") && !p.Allowed("admin:read") {
		writeProblem(w, 403, "FORBIDDEN", "operational permission is required")
		return
	}
	organization := r.URL.Query().Get("organization_id")
	if organization == "" || !p.AllowedOrganization(organization) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "organization scope is required")
		return
	}
	value, err := a.service.Operations(r.Context(), p.TenantID, organization)
	if err != nil {
		writeProblem(w, 503, "OPERATIONS_UNAVAILABLE", "operation snapshot unavailable")
		return
	}
	if (a.paymentProvider == "stripe" || a.paymentProvider == "mercadopago") && (a.paymentTenantID == "" || (a.paymentTenantID == p.TenantID && a.paymentOrganizationID == organization)) {
		value.PaymentProvider = a.paymentProvider
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, 200, value)
}

func (a commerceAPI) auth(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return p, false
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return p, false
	}
	return p, true
}
func (a commerceAPI) publicPrice(w http.ResponseWriter, r *http.Request) {
	value, err := a.service.PublicPrice(r.Context(), r.PathValue("tenantCode"), r.URL.Query().Get("market"), r.URL.Query().Get("variant_id"))
	if err != nil {
		writeProblem(w, 404, "PRICE_NOT_FOUND", "one current price was not found")
		return
	}
	writeJSON(w, 200, value)
}
func (a commerceAPI) createBook(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "pricing:write")
	if !ok {
		return
	}
	var input commerce.PriceBook
	if !decodeStrict(w, r, &input) {
		return
	}
	value, err := a.service.CreatePriceBook(r.Context(), p.TenantID, input)
	if err != nil {
		writeProblem(w, 400, "INVALID_PRICE_BOOK", "price book does not match contract")
		return
	}
	writeJSON(w, 201, value)
}
func (a commerceAPI) activateBook(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "pricing:write")
	if !ok {
		return
	}
	writeCommerceResult(w, a.service.ActivatePriceBook(r.Context(), p.TenantID, r.PathValue("id")))
}
func (a commerceAPI) addLine(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "order:write")
	if !ok {
		return
	}
	var input struct {
		OrganizationID  string `json:"organization_id"`
		PriceBookID     string `json:"price_book_id"`
		VariantID       string `json:"variant_id"`
		Quantity        int    `json:"quantity"`
		ExpectedVersion int64  `json:"expected_version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.AddOrderLine(r.Context(), p.TenantID, commerce.OrderLine{OrderID: r.PathValue("id"), OrganizationID: input.OrganizationID, PriceBookID: input.PriceBookID, VariantID: input.VariantID, Quantity: input.Quantity}, input.ExpectedVersion)
	if err != nil {
		writeCommerceResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a commerceAPI) placeOrder(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "order:write")
	if !ok {
		return
	}
	var input struct {
		OrganizationID  string `json:"organization_id"`
		ExpectedVersion int64  `json:"expected_version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeCommerceResult(w, a.service.PlaceOrder(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.ExpectedVersion))
}
func (a commerceAPI) allocate(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "inventory:allocate")
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		LineID         string `json:"line_id"`
		StockUnitID    string `json:"stock_unit_id"`
		OrderVersion   int64  `json:"order_version"`
		StockVersion   int64  `json:"stock_version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeCommerceResult(w, a.service.AllocateStockAs(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.LineID, input.StockUnitID, input.OrderVersion, input.StockVersion, p.Subject))
}
func (a commerceAPI) paymentEnabled(w http.ResponseWriter) bool {
	if a.paymentProvider != "stripe" && a.paymentProvider != "mercadopago" {
		writeProblem(w, 503, "PAYMENT_REQUEST_DISABLED", "payment provider selection is required; no payment effect enabled")
		return false
	}
	return true
}
func (a commerceAPI) paymentScope(w http.ResponseWriter, p identity.Principal, organization string) bool {
	if (a.paymentTenantID != "" || a.paymentOrganizationID != "") && (a.paymentTenantID != p.TenantID || a.paymentOrganizationID != organization) {
		writeProblem(w, 503, "PAYMENT_SCOPE_DISABLED", "payment checkout is not enabled for this scope")
		return false
	}
	return true
}
func (a commerceAPI) requestOrderPayment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "payment:create")
	if !ok {
		return
	}
	if !a.paymentEnabled(w) {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		RequestKey     string `json:"request_key"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "organization scope is required")
		return
	}
	if !a.paymentScope(w, p, input.OrganizationID) {
		return
	}
	v, err := a.service.RequestOrderPayment(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), a.paymentProvider, input.RequestKey, p.Subject)
	if err != nil {
		writeCommerceResult(w, err)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, 202, v)
}
func (a commerceAPI) createPayment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "payment:create")
	if !ok {
		return
	}
	if !a.paymentEnabled(w) {
		return
	}
	var input commerce.PaymentAttempt
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	input.OrderID = r.PathValue("id")
	if !a.paymentScope(w, p, input.OrganizationID) {
		return
	}
	if input.ProviderCode != a.paymentProvider {
		writeProblem(w, 400, "PAYMENT_PROVIDER_MISMATCH", "provider must match configured selection")
		return
	}
	value, err := a.service.CreatePaymentAttemptAs(r.Context(), p.TenantID, r.Header.Get("Idempotency-Key"), input, p.Subject)
	if err != nil {
		writeCommerceResult(w, err)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, 202, value)
}
func (a commerceAPI) transitionPayment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "payment:write")
	if !ok {
		return
	}
	if !a.paymentEnabled(w) {
		return
	}
	var input struct {
		OrganizationID    string `json:"organization_id"`
		Current           string `json:"current"`
		Target            string `json:"target"`
		Version           int64  `json:"version"`
		ProviderReference string `json:"provider_reference"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	if !a.paymentScope(w, p, input.OrganizationID) {
		return
	}
	// In the provider-observed profile, an operator's JSON cannot establish
	// authorization, capture, refund or dispute facts. Only a never-dispatched
	// created request may be cancelled here; the repository CAS rejects it if
	// the atomic send claim has already changed it to pending.
	if a.providerObservedPayments && (input.Current != "created" || input.Target != "failed" || input.ProviderReference != "") {
		writeProblem(w, 403, "PROVIDER_OBSERVATION_REQUIRED", "financial provider state requires authenticated reconciliation")
		return
	}
	writeCommerceResult(w, a.service.TransitionPayment(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Current, input.Target, input.Version, input.ProviderReference))
}
func writeCommerceResult(w http.ResponseWriter, err error) {
	if errors.Is(err, commerce.ErrPaymentUnavailable) {
		writeProblem(w, 503, "PAYMENT_INTENT_UNAVAILABLE", "result unavailable; consult persisted state before retrying")
		return
	}
	if errors.Is(err, commerce.ErrConflict) {
		writeProblem(w, 409, "COMMERCE_CONFLICT", "state, version or business invariant conflict")
		return
	}
	if err != nil {
		writeProblem(w, 400, "INVALID_COMMERCE_COMMAND", "command does not match contract")
		return
	}
	writeJSON(w, 200, map[string]string{"status": "accepted"})
}
````

### FILE: `internal/platform/httpapi/commerce_test.go`

```yaml
block_id: "GO-COMMERCE-API:internal-platform-httpapi-commerce_test-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "8c62000fb647c78217fe60b45ac92974e0006f937a6c450848c92d126919fcb4"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

type commerceRepo struct{ books int }

func TestPaymentRequestsDisabledWithoutSelection(t *testing.T) {
	for _, provider := range []string{"", "unsupported"} {
		mux := http.NewServeMux()
		p := identity.Principal{Subject: "operator", TenantID: "tenant", Permissions: map[string]struct{}{"payment:create": {}, "payment:write": {}}, Organizations: map[string]struct{}{"org": {}}}
		CommerceModule{Service: commerce.NewService(&commerceRepo{}, &commerceIDs{}), PaymentProvider: provider}.Register(mux, paymentPrincipal{p})
		for _, path := range []string{"/v1/commerce/orders/order/payment-request", "/v1/commerce/orders/order/payments", "/v1/payments/payment/transitions"} {
			r := httptest.NewRequest("POST", path, strings.NewReader(`{}`))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Authorization", "Bearer synthetic")
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			if w.Code != 503 || !strings.Contains(w.Body.String(), "PAYMENT_REQUEST_DISABLED") {
				t.Fatalf("disabled route %s status%d", path, w.Code)
			}
		}
	}
}

// The HTTP boundary uses a synthetic verified principal; real OIDC is covered
// by the connected portal gate, not claimed by this focused payment regression.
type paymentPrincipal struct{ identity.Principal }

func (v paymentPrincipal) Verify(context.Context, string) (identity.Principal, error) {
	return v.Principal, nil
}

type paymentFailureIDs struct {
	values []string
	n      int
}

func (g *paymentFailureIDs) New() string { v := g.values[g.n]; g.n++; return v }

func TestPaymentHTTPDurableReplayPostgres(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	u, err := url.Parse(dsn)
	if err != nil || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_confirmation_") {
		t.Fatal("requires explicit loopback confirmation database")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { pool.Close() }()
	ids := randomid.Generator{}
	tenant := ids.New()
	for _, sql := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'test-'||($1::uuid)::text,'Synthetic payment test','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'org','store','Store','store'),($1,'other','other','Other','store')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','org','synthetic','placed','ARS',1000,1),($1,'order-2','org','synthetic','placed','ARS',1000,1),($1,'order-3','org','synthetic','placed','ARS',1000,1),($1,'private','other','synthetic','placed','ARS',1000,1)`,
	} {
		if _, err := pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("synthetic tenant=%s", tenant)
	principal := identity.Principal{Subject: "payment-operator", TenantID: tenant, Permissions: map[string]struct{}{"payment:create": {}}, Organizations: map[string]struct{}{"org": {}}}
	var handler http.Handler
	rebind := func(p identity.Principal) {
		mux := http.NewServeMux()
		CommerceModule{Service: commerce.NewService(postgres.NewCommerce(pool), ids), PaymentProvider: "stripe"}.Register(mux, paymentPrincipal{p})
		handler = mux
	}
	rebind(principal)
	body := `{"organization_id":"org","provider_code":"stripe","currency":"ARS","amount_minor_units":1000}`
	request := func(target, key, payload string, auth bool) *httptest.ResponseRecorder {
		r := httptest.NewRequest("POST", "/v1/commerce/orders/"+target+"/payments", strings.NewReader(payload))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Idempotency-Key", key)
		if auth {
			r.Header.Set("Authorization", "Bearer synthetic")
		}
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		return w
	}
	key := "payment-replay-00000001"
	first := request("order", key, body, true)
	if first.Code != 202 {
		t.Fatalf("first status=%d body=%s", first.Code, first.Body.String())
	}
	var original commerce.PaymentAttempt
	if err := json.Unmarshal(first.Body.Bytes(), &original); err != nil {
		t.Fatal(err)
	}
	// Discard the first response as a client that did not retain its receipt.
	// This is not claimed as an injected transport drop.
	second := request("order", key, body, true)
	if second.Code != 202 {
		t.Fatalf("durable replay status=%d want202 body=%s", second.Code, second.Body.String())
	}
	if second.Body.String() != first.Body.String() {
		t.Fatal("replay differs from durable original")
	}
	var wg sync.WaitGroup
	responses := make(chan *httptest.ResponseRecorder, 16)
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); responses <- request("order", key, body, true) }()
	}
	wg.Wait()
	close(responses)
	for w := range responses {
		if w.Code != 202 || w.Body.String() != first.Body.String() {
			t.Fatalf("concurrent replay status=%d", w.Code)
		}
	}
	for _, c := range []struct {
		name, target, key, body string
		want                    int
	}{
		{"amount", "order", key, strings.Replace(body, ":1000", ":999", 1), 409},
		{"currency", "order", key, strings.Replace(body, "ARS", "USD", 1), 409},
		{"other-order", "order-2", key, body, 409},
		{"private-order", "private", key, body, 409},
		{"organization", "order", key, strings.Replace(body, `"org"`, `"other"`, 1), 403},
		{"invalid-key", "order", "short", body, 400},
	} {
		t.Run(c.name, func(t *testing.T) {
			if w := request(c.target, c.key, c.body, true); w.Code != c.want {
				t.Fatalf("status%d want%d", w.Code, c.want)
			}
		})
	}
	if w := request("order", key, body, false); w.Code != 401 {
		t.Fatalf("anonymous%d", w.Code)
	}
	denied := principal
	denied.Permissions = map[string]struct{}{}
	rebind(denied)
	if w := request("order", key, body, true); w.Code != 403 {
		t.Fatalf("permission%d", w.Code)
	}
	other := principal
	other.TenantID = ids.New()
	rebind(other)
	if w := request("order", key, body, true); w.Code != 409 {
		t.Fatalf("tenant%d", w.Code)
	}
	rebind(principal)
	var rows, events int
	if err := pool.QueryRow(ctx, `select count(*) from payment.payment_attempt where tenant_id=$1`, tenant).Scan(&rows); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='payment.requested' and payload->>'actor_subject'='payment-operator'`, tenant).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if rows != 1 || events != 1 {
		t.Fatalf("rows%d audited events%d", rows, events)
	}
	// A different order must be able to await a provider reference concurrently.
	responses = make(chan *httptest.ResponseRecorder, 16)
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); responses <- request("order-2", "payment-replay-00000002", body, true) }()
	}
	wg.Wait()
	close(responses)
	receipt := ""
	for w := range responses {
		if w.Code != 202 {
			t.Fatalf("fresh concurrent status%d", w.Code)
		}
		if receipt == "" {
			receipt = w.Body.String()
		}
		if receipt != w.Body.String() {
			t.Fatal("fresh concurrent IDs differ")
		}
	}
	if err := pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='payment.requested'`, tenant).Scan(&events); err != nil || events != 2 {
		t.Fatalf("fresh race events%d err%v", events, err)
	}
	if w := request("order-2", "another-distinct-key-0002", body, true); w.Code != 409 {
		t.Fatalf("legacy route bypassed initial payment exclusivity status%d", w.Code)
	}
	// Force an outbox unique failure after the payment INSERT, not before it.
	var duplicateEvent string
	if err := pool.QueryRow(ctx, `select event_id::text from platform.outbox_event where tenant_id=$1 and aggregate_id=$2`, tenant, original.ID).Scan(&duplicateEvent); err != nil {
		t.Fatal(err)
	}
	faultMux := http.NewServeMux()
	CommerceModule{Service: commerce.NewService(postgres.NewCommerce(pool), &paymentFailureIDs{values: []string{ids.New(), duplicateEvent}}), PaymentProvider: "stripe"}.Register(faultMux, paymentPrincipal{principal})
	handler = faultMux
	if w := request("order-3", "payment-replay-00000003", body, true); w.Code != 400 {
		t.Fatalf("outbox failure status%d", w.Code)
	}
	if err := pool.QueryRow(ctx, `select count(*) from payment.payment_attempt where tenant_id=$1 and order_id='order-3'`, tenant).Scan(&rows); err != nil || rows != 0 {
		t.Fatalf("partial intent after outbox failure rows%d err%v", rows, err)
	}
	rebind(principal)
	if w := request("order-3", "payment-replay-00000003", body, true); w.Code != 202 {
		t.Fatalf("retry after rollback status%d", w.Code)
	}
	if _, err := pool.Exec(ctx, `update payment.payment_attempt set state='pending',version=2 where tenant_id=$1 and payment_attempt_id=$2`, tenant, original.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `update sales.customer_order set state='cancelled',version=2 where tenant_id=$1 and order_id='order'`, tenant); err != nil {
		t.Fatal(err)
	}
	pool.Close()
	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	rebind(principal)
	w := request("order", key, body, true)
	var recovered commerce.PaymentAttempt
	if w.Code != 202 || json.Unmarshal(w.Body.Bytes(), &recovered) != nil || recovered.ID != original.ID || recovered.State != "pending" || recovered.Version != 2 {
		t.Fatalf("reconnect replay status%d", w.Code)
	}
	if _, err := pool.Exec(ctx, `update payment.payment_attempt set provider_reference='same-external-id' where tenant_id=$1`, tenant); err == nil {
		t.Fatal("duplicate non-null provider reference accepted")
	}
	t.Log("PAYMENT_HTTP_REPLAY_PASS sequential+16 concurrent replay+16 fresh race; one audited event per intent; changed intent/scope denied; independent null references; atomic outbox failure/retry; reconnect/current-state recovery")
}

func (c *commerceRepo) CreatePriceBook(context.Context, string, string, commerce.PriceBook) error {
	c.books++
	return nil
}
func (c *commerceRepo) ActivatePriceBook(context.Context, string, string, string) error {
	return commerce.ErrConflict
}
func (c *commerceRepo) PublicPrice(context.Context, string, string, string) (commerce.PriceEntry, error) {
	return commerce.PriceEntry{VariantID: "v", AmountMinorUnits: 100, TaxMode: "inclusive"}, nil
}
func (c *commerceRepo) AddOrderLine(_ context.Context, _ string, _ string, v commerce.OrderLine, _ int64) (commerce.OrderLine, error) {
	return v, nil
}
func (c *commerceRepo) PlaceOrder(context.Context, string, string, string, int64, string) error {
	return nil
}
func (c *commerceRepo) AllocateStock(context.Context, string, string, string, string, string, int64, int64, string) error {
	return nil
}
func (c *commerceRepo) CreatePaymentAttempt(context.Context, string, string, string, commerce.PaymentAttempt) error {
	return nil
}
func (c *commerceRepo) TransitionPayment(context.Context, string, string, string, string, string, int64, string, string) error {
	return nil
}

type commerceIDs struct{ n int }

func (i *commerceIDs) New() string {
	i.n++
	return []string{"018f4d4a-7b36-7a21-8d10-2f4c54c28501", "018f4d4a-7b36-7a21-8d10-2f4c54c28502"}[i.n-1]
}

type commerceVerifier struct{}

type operationVerifier struct{ permission, organization string }

func (v operationVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{Subject: "operator", TenantID: "tenant", Permissions: map[string]struct{}{v.permission: {}}, Organizations: map[string]struct{}{v.organization: {}}}, nil
}

func TestOperationReadFailsClosedWithoutScopeOrReader(t *testing.T) {
	for _, c := range []struct {
		name, permission, organization, query string
		auth                                  bool
		status                                int
	}{
		{"anonymous", "inventory:allocate", "store", "store", false, 401},
		{"customer", "customer:self", "store", "store", true, 403},
		{"other-org", "inventory:allocate", "other", "store", true, 403},
		{"missing-org", "inventory:allocate", "store", "", true, 403},
		{"reader-unavailable", "inventory:allocate", "store", "store", true, 503},
	} {
		t.Run(c.name, func(t *testing.T) {
			mux := http.NewServeMux()
			CommerceModule{Service: commerce.NewService(&commerceRepo{}, &commerceIDs{})}.Register(mux, operationVerifier{c.permission, c.organization})
			r := httptest.NewRequest("GET", "/v1/commerce/orders?organization_id="+c.query, nil)
			if c.auth {
				r.Header.Set("Authorization", "Bearer fixture")
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			if w.Code != c.status {
				t.Fatalf("got%d expected%d", w.Code, c.status)
			}
		})
	}
}

func (commerceVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{Subject: "admin", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c28500", Permissions: map[string]struct{}{"pricing:write": {}}}, nil
}
func TestCommerceHTTPPublicAndProtected(t *testing.T) {
	repo := &commerceRepo{}
	service := commerce.NewService(repo, &commerceIDs{})
	mux := http.NewServeMux()
	CommerceModule{Service: service}.Register(mux, commerceVerifier{})
	request := httptest.NewRequest("GET", "/v1/public/acme/prices?market=AR&variant_id=v", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatalf("price status=%d", response.Code)
	}
	request = httptest.NewRequest("POST", "/v1/pricing/books", strings.NewReader(`{"market":"AR","currency":"ARS","valid_from":"`+time.Now().UTC().Format(time.RFC3339)+`","entries":[{"variant_id":"v","amount_minor_units":100,"tax_mode":"inclusive"}]}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.books != 1 {
		t.Fatalf("book status=%d books=%d body=%s", response.Code, repo.books, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/commerce/orders/o/payments", strings.NewReader(`{}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 {
		t.Fatalf("payment permission status=%d", response.Code)
	}
}
````

### FILE: `db/migrations/0053_payment_intent_reference.up.sql`

```yaml
block_id: "GO-COMMERCE-API:payment-reference-up:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "44398870ddd011c716d966e32a7d6dd7e1e26d519aaf0a47d817c3413ef1f140"
variables: []
secrets_allowed: false
```

````sql
begin;
-- AUTHORED: absence of a provider receipt is not a shared provider identity.
-- Keep the original migration immutable and preserve non-null uniqueness.
alter table payment.payment_attempt
  drop constraint payment_attempt_tenant_id_provider_code_provider_reference_key,
  add constraint payment_attempt_tenant_id_provider_code_provider_reference_key
    unique nulls distinct (tenant_id,provider_code,provider_reference);
commit;
````

### FILE: `db/migrations/0053_payment_intent_reference.down.sql`

```yaml
block_id: "GO-COMMERCE-API:payment-reference-down:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "0857363ca840794c480074e45798a5123ae4b89ef69056e6c12a58d18e975079"
variables: []
secrets_allowed: false
```

````sql
begin;
-- Fails atomically if multiple pending references now exist. Never delete,
-- merge or manufacture payment references to make a downgrade pass.
alter table payment.payment_attempt
  drop constraint payment_attempt_tenant_id_provider_code_provider_reference_key,
  add constraint payment_attempt_tenant_id_provider_code_provider_reference_key
    unique nulls not distinct (tenant_id,provider_code,provider_reference);
commit;
````

## 6. Configuration surface

No standalone environment variables are introduced. The pack consumes the composed PostgreSQL pool and verified identity/authorization. Provider selection, allowlist and execution adapter remain project integration requirements: accepting a local intent does not authorize a remote payment effect. Market, currency, tax mode, permissions and state transitions remain validated domain inputs. A UI must retain its original Idempotency-Key for the same intent; a new key is not deduplicated against an earlier key merely because the order matches.

## 7. Dependency bill

| Dependency | Pin | Use | License | Scope | Source |
|---|---|---|---|---|---|
| Go | `1.26.7` | service/HTTP/tests | BSD-3-Clause | build/runtime | `go.dev` |
| PostgreSQL | `18.6` verified baseline | transactional source of truth | PostgreSQL | runtime/test | `postgresql.org` |
| `pgx` | `5.10.0` | PostgreSQL adapter | MIT | build/runtime | `github.com/jackc/pgx` |

## 8. Apply order

Compose after foundation, domain migration, CRM and supply packs; then wire the module in the application root. Apply SQL before serving routes and execute negative authorization/concurrency tests. Existing systems must map monetary units, order ownership and provider idempotency explicitly. Rollback disables routes/workers first and preserves durable records; schema reversal is development-only unless proven safe.

## 9. Verification

Compose with Go application and migrations 0001/0002/0003/0053. Run format, unit/HTTP/PostgreSQL integration, vet and build. La revisión 0.2 limita órdenes, asignaciones y pagos por organización en HTTP y SQL, con pruebas negativas. V297 prueba el recibo idempotente mediante handlers/principal sintético y PostgreSQL; no sustituye OIDC/BFF/navegador del pago. Producción requiere reglas fiscales/contables, aprobación de precios, descuentos/financiación, adapters oficiales con firma de webhook/inbox/reconciliación, pagos parciales/refunds, decisión PCI, fraude/riesgo y carga concurrente.

## 10. Reconstruction evidence

Reconstruction, database integration and HTTP authorization evidence is recorded in `reconstruction_evidence/GO_COMMERCE_PRICING_PAYMENT_API_2026-08-24_V1.md`; the final audit revalidates current hashes and the full backend profile.

V402 composed delta: Payment/initial-handover composition. New behavior and tests belong to the explicit runtime/portal packs; source and library release claims remain bounded to their evidence. Existing source provenance is preserved.

Canonical V402 integration: selected by the current profile with exact dependencies and caller overlays. Metadata promotion records byte reconstruction, not closure of every admission/release gate. Payload provenance is unchanged.

V402 composed delta: V402 source-backed stored-value integration: exact remaining provider due, explicit payment/funding XOR, shared approval, bounded browser transport and optional host. See STORED_VALUE_OPERATOR_FLOW_V402.md; source/pack admission successor governs final claim. Existing provider-only behavior retained.


### FILE: `db/migrations/0066_order_funding_projection.up.sql`

```yaml
block_id: "GO-COMMERCE-PRICING-PAYMENT-API:db/migrations/0066_order_funding_projection.up.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "Local typed source/transaction/transport/UI/recovery glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "dc853c12918a17d3eea1912ceb34245bdf29ebab80d2cf98b34564d3db2abd75"
variables: []
secrets_allowed: false
```

````sql
begin;
-- AUTHORED base projection. Gross order pricing remains with sales; this view
-- introduces no discount rule, alternate order, payment or accounting entry.
create view payment.order_funding as
 select tenant_id,order_id,organization_id,currency,version as order_version,
 total_minor_units as gross_minor_units,0::bigint as gift_minor_units,
 0::bigint as discount_minor_units,total_minor_units as provider_due_minor_units,
 encode(sha256(convert_to('[]','UTF8')),'hex') as contributions_sha256
 from sales.customer_order;
create table payment.local_funding_receipt (
 tenant_id uuid not null,
 funding_id text not null,
 request_key text not null check(length(request_key) between 16 and 128),
 request_sha256 text not null check(request_sha256 ~ '^[0-9a-f]{64}$'),
 requested_by_subject text not null check(length(requested_by_subject) between 1 and 256),
 receipt_raw bytea not null check(octet_length(receipt_raw) between 1 and 65536),
 order_id text not null,
 organization_id text not null,
 currency text not null,
 gross_minor_units bigint not null check(gross_minor_units>0),
 gift_minor_units bigint not null check(gift_minor_units>=0),
 discount_minor_units bigint not null check(discount_minor_units>=0),
 provider_minor_units bigint not null check(provider_minor_units>=0),
 payment_attempt_id text,
 provider_evidence_sha256 text,
 allocation_sha256 text not null check(allocation_sha256 ~ '^[0-9a-f]{64}$'),
 receipt_sha256 text not null check(receipt_sha256 ~ '^[0-9a-f]{64}$'),
 receipt jsonb not null check(jsonb_typeof(receipt)='object' and pg_column_size(receipt)<=65536),
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,funding_id),
 unique(tenant_id,request_key),
 foreign key(tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 foreign key(tenant_id,payment_attempt_id) references payment.payment_attempt(tenant_id,payment_attempt_id),
 check(gross_minor_units::numeric=gift_minor_units::numeric+discount_minor_units::numeric+provider_minor_units::numeric),
 check((provider_minor_units=0 and payment_attempt_id is null and provider_evidence_sha256 is null) or
       (provider_minor_units>0 and payment_attempt_id is not null and provider_evidence_sha256 is not null and provider_evidence_sha256 ~ '^[0-9a-f]{64}$'))
);
create view payment.local_funding_evidence as select tenant_id,funding_id,order_id,organization_id,currency,gross_minor_units,gift_minor_units,discount_minor_units,allocation_sha256,receipt_raw,receipt_sha256,created_at from payment.local_funding_receipt where provider_minor_units=0 and payment_attempt_id is null and provider_evidence_sha256 is null;
create index local_funding_order on payment.local_funding_receipt(tenant_id,order_id,created_at desc);
create function payment.local_funding_immutable() returns trigger language plpgsql as $$
begin raise exception 'local funding observation is immutable';end $$;
create trigger local_funding_immutable before update or delete on payment.local_funding_receipt for each row execute function payment.local_funding_immutable();
commit;
````

### FILE: `db/migrations/0066_order_funding_projection.down.sql`

```yaml
block_id: "GO-COMMERCE-PRICING-PAYMENT-API:db/migrations/0066_order_funding_projection.down.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "Local typed source/transaction/transport/UI/recovery glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "82373b270f2c171ad97ea86b7e469272c5a4b1b88cd325c3117305e4a808fa33"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$begin if exists(select 1 from payment.local_funding_receipt) then raise exception 'preserve local funding evidence';end if;end $$;
-- No CASCADE: dependent handover columns must be rolled back explicitly first.
drop view payment.local_funding_evidence;
drop view payment.order_funding;
drop table payment.local_funding_receipt;
drop function payment.local_funding_immutable();
commit;
````

### FILE: `internal/platform/postgres/order_funding.go`

```yaml
block_id: "GO-COMMERCE-PRICING-PAYMENT-API:internal/platform/postgres/order_funding.go:v1"
operation: CREATE
provenance: AUTHORED
source: "Local typed source/transaction/transport/UI/recovery glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7aa822ccc948733553a852b3cc7934561719c7f8fd55b9f9068a46f5a2cb5097"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED exact projection/hash glue. The selected source-derived writer owns
// the gift/discount rules; all readers share this representation under order lock.
import (
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/commerce"
	"encoding/hex"
	"encoding/json"
	"github.com/jackc/pgx/v5"
)

type orderFundingReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}
type OrderFundingSnapshot struct {
	TenantID            string `json:"tenant_id"`
	OrganizationID      string `json:"organization_id"`
	OrderID             string `json:"order_id"`
	Currency            string `json:"currency"`
	GrossMinor          int64  `json:"gross_minor_units"`
	GiftMinor           int64  `json:"gift_minor_units"`
	DiscountMinor       int64  `json:"discount_minor_units"`
	ProviderMinor       int64  `json:"provider_minor_units"`
	ContributionsSHA256 string `json:"contributions_sha256"`
}

func fundingHash(raw []byte) string           { h := sha256.Sum256(raw); return hex.EncodeToString(h[:]) }
func (v OrderFundingSnapshot) SHA256() string { raw, _ := json.Marshal(v); return fundingHash(raw) }
func readOrderFunding(ctx context.Context, q orderFundingReader, tenant, org, order, currency string, gross int64) (OrderFundingSnapshot, error) {
	v := OrderFundingSnapshot{TenantID: tenant, OrganizationID: org, OrderID: order}
	e := q.QueryRow(ctx, `select currency,gross_minor_units,gift_minor_units,discount_minor_units,provider_due_minor_units,contributions_sha256 from payment.order_funding where tenant_id=$1 and organization_id=$2 and order_id=$3`, tenant, org, order).Scan(&v.Currency, &v.GrossMinor, &v.GiftMinor, &v.DiscountMinor, &v.ProviderMinor, &v.ContributionsSHA256)
	if e != nil {
		return v, e
	}
	hash, e := hex.DecodeString(v.ContributionsSHA256)
	if currency != v.Currency || gross != v.GrossMinor || gross <= 0 || v.GiftMinor < 0 || v.GiftMinor > gross || v.DiscountMinor < 0 || v.DiscountMinor > gross-v.GiftMinor || v.ProviderMinor < 0 || v.ProviderMinor != gross-v.GiftMinor-v.DiscountMinor || e != nil || len(hash) != 32 || hex.EncodeToString(hash) != v.ContributionsSHA256 {
		return v, commerce.ErrConflict
	}
	return v, nil
}
func orderProviderDue(ctx context.Context, q orderFundingReader, tenant, org, order, currency string, gross int64) (int64, error) {
	v, e := readOrderFunding(ctx, q, tenant, org, order, currency, gross)
	return v.ProviderMinor, e
}
````

### FILE: `internal/platform/postgres/local_funding.go`

```yaml
block_id: "GO-COMMERCE-PRICING-PAYMENT-API:internal/platform/postgres/local_funding.go:v1"
operation: CREATE
provenance: AUTHORED
source: "Local typed source/transaction/transport/UI/recovery glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7afb4efeb0c1aa24130ac55ccc4bdb418c71f85211468ab62ccadf2eecd073bc"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED immutable observation of already approved contributions. This is
// not a provider capture, cash movement, general ledger or gift-card rule.
import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/commerce"
	"encoding/json"
	"io"
	"time"
)

const LocalFundingEffect = "ORDER_FULLY_FUNDED_BY_STORED_VALUE"

type LocalFundingReceipt struct {
	ID               string               `json:"funding_receipt_id"`
	RequestKey       string               `json:"request_key"`
	RequestSHA256    string               `json:"request_sha256"`
	RequestedBy      string               `json:"requested_by"`
	Allocation       OrderFundingSnapshot `json:"allocation"`
	AllocationSHA256 string               `json:"allocation_sha256"`
	ObservedAt       time.Time            `json:"observed_at"`
	Effect           string               `json:"effect"`
}
type LocalFundingResult struct {
	Receipt LocalFundingReceipt `json:"receipt"`
	SHA256  string              `json:"receipt_sha256"`
}

// Order is locked by the caller; contributions cannot change during validation.
func readLocalFunding(ctx context.Context, q orderFundingReader, tenant, org, order, id, evidence, currency string, gross int64) (LocalFundingResult, error) {
	var out LocalFundingResult
	var raw []byte
	var storedAllocation string
	var amount, gift, discount int64
	var at time.Time
	e := q.QueryRow(ctx, `select receipt_raw,receipt_sha256,allocation_sha256,gross_minor_units,gift_minor_units,discount_minor_units,created_at from payment.local_funding_evidence where tenant_id=$1 and organization_id=$2 and order_id=$3 and funding_id=$4 and currency=$5`, tenant, org, order, id, currency).Scan(&raw, &out.SHA256, &storedAllocation, &amount, &gift, &discount, &at)
	if e != nil {
		return out, e
	}
	if len(raw) > 65536 || fundingHash(raw) != out.SHA256 || out.SHA256 != evidence {
		return out, commerce.ErrConflict
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&out.Receipt) != nil || decoder.Decode(new(any)) != io.EOF {
		return out, commerce.ErrConflict
	}
	current, e := readOrderFunding(ctx, q, tenant, org, order, currency, gross)
	if e != nil {
		return out, e
	}
	r := out.Receipt
	if current.ProviderMinor != 0 || current.GiftMinor+current.DiscountMinor <= 0 || r.ID != id || r.Effect != LocalFundingEffect || !r.ObservedAt.Equal(at) || r.Allocation != current || r.AllocationSHA256 != current.SHA256() || r.AllocationSHA256 != storedAllocation || amount != gross || gift != current.GiftMinor || discount != current.DiscountMinor {
		return out, commerce.ErrConflict
	}
	return out, nil
}
````

V402 composed delta: J3 immutable approved catalog publication reuses original Commerce SQL/shared approval and optional host/public model owner; existing Next storefront consumes a validated published projection. No new dependencies or corporate attribution. CATALOG_CONNECTED_RELEASE_V402.md.
