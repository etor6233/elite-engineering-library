# Go ARCA Fiscal Issuance API

## 1. Metadata

```yaml
pack_id: "GO-ARCA-FISCAL-ISSUANCE-API"
pack_version: "0.6.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Añade solicitud fiscal, comprobante asociado durable para la lane estricta factura/nota de crédito A-B-C y registro de parámetros WSFEv1 por organización/CUIT, manteniendo adquisición y aprobación explícitamente separadas."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "ARCA WSFEv1 hash-locked homologation WSDL"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "GO-COMMERCE-PRICING-PAYMENT-API 0.3.x", "GO-ENTERPRISE-ACCOUNTING-LEDGER-API 0.1.x", "MICROSOFT-ARCA-WSAA-CREDENTIAL-CORE 0.1.x", "MICROSOFT-ARCA-WSFE-GENERATED-CLIENT 0.2.x", "GO-ELECTROMOBILITY-APPLICATION 1.9.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner AND ARCA public specification AND MIT upstream evidence"
upstream_sources: ["https://www.arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG.pdf", "https://www.arca.gob.ar/ws/documentacion/wsaa.asp", "https://github.com/microsoft/BCApps/tree/31a860b527f0dc72c7a44a255d7e7d403cfa4789", "https://www.postgresql.org/docs/18/"]
verified_at: "2026-09-13"
```

Los veintiocho bloques son `AUTHORED`: no se atribuye Go o SQL local a ARCA, Microsoft, PostgreSQL ni Google. ARCA gobierna el protocolo fiscal, `CbtesAsoc` y sus operaciones paramétricas; Microsoft BCApps fijado aporta evidencia de separación entre documento fiscal, VAT entries, posting y reversión. PostgreSQL gobierna transacciones, locks y `SKIP LOCKED`. La implementación local conserva su procedencia real y no embebe tablas fiscales.

## 2. Applicability

Use después de pedidos, pagos capturados, organizaciones y outbox. Resuelve la propiedad durable de una solicitud WSFEv1 y su recuperación ante fallos de comunicación. Rechácelo para moneda distinta de ARS, cálculo tributario no aprobado, comprobantes no soportados por WSFEv1 o uso sin CUIT/certificado/relación/punto de venta de homologación reales. No autoriza producción ni reemplaza asesoramiento fiscal, tablas paramétricas vigentes, adapter SOAP exacto o aceptación del contribuyente.

## 3. Architecture contract

Pedido y pago capturado son la fuente económica; el cliente no puede alterar moneda o total. Para notas de crédito la única lane admitida relaciona exactamente `3→1`, `8→6` o `13→11` con una factura local ya autorizada del mismo contribuyente y organización; PostgreSQL carga y vuelve inmutable su CUIT, tipo, punto de venta, número y fecha, sin confiar esa identidad al cliente. La condición IVA del receptor es una entrada explícita obligatoria y nunca se deriva del tipo o número de documento. La composición fiscal, cada base/alícuota IVA y cada tributo deben cuadrar en unidades menores antes de persistir; PostgreSQL revalida sus sumas mediante constraint triggers diferidos y vuelve inmutables los detalles. Una lane durable por CUIT, ambiente, punto de venta y tipo impide asignaciones concurrentes. El procesador consulta `FECompUltimoAutorizado`, asigna exactamente el siguiente número y llama `FECAESolicitar`. Si la respuesta puede haberse perdido, nunca reenvía a ciegas: conserva número y estado `reconcile_required`, luego ejecuta `FECompConsultar`; sólo una consulta negativa válida permite reenviar la misma solicitud. CAE, vencimiento, hashes, códigos y cada intento quedan auditados sin Token, Sign, certificado ni payload secreto. El adapter SOAP y WSAA son transportes separados; PostgreSQL sigue siendo el dueño del workflow. El evento autorizado alimenta el ledger existente sin duplicarlo.

## 4. Exact file manifest

```text
CREATE internal/platform/httpapi/fiscal_connected_reference_test.go
CREATE internal/fiscal/service.go
CREATE internal/fiscal/service_test.go
CREATE internal/fiscal/processor.go
CREATE internal/fiscal/processor_test.go
CREATE internal/platform/postgres/fiscal.go
CREATE internal/platform/postgres/fiscal_integration_test.go
CREATE internal/platform/httpapi/fiscal.go
CREATE internal/platform/httpapi/fiscal_test.go
CREATE db/migrations/0012_arca_fiscal_issuance.up.sql
CREATE db/migrations/0012_arca_fiscal_issuance.down.sql
CREATE db/tests/0012_arca_fiscal_issuance.test.sql
CREATE internal/fiscal/parameters.go
CREATE internal/fiscal/parameters_test.go
CREATE internal/platform/postgres/fiscal_parameters.go
CREATE db/migrations/0013_arca_fiscal_parameters.up.sql
CREATE db/migrations/0013_arca_fiscal_parameters.down.sql
CREATE db/tests/0013_arca_fiscal_parameters.test.sql
CREATE internal/fiscal/parameter_refresh.go
CREATE internal/fiscal/parameter_refresh_test.go
CREATE internal/platform/postgres/fiscal_parameter_schedule_integration_test.go
CREATE internal/platform/httpapi/fiscal_parameters.go
CREATE internal/platform/httpapi/fiscal_parameters_test.go
CREATE db/migrations/0014_arca_parameter_refresh_admin.up.sql
CREATE db/migrations/0014_arca_parameter_refresh_admin.down.sql
CREATE db/tests/0014_arca_parameter_refresh_admin.test.sql
CREATE db/migrations/0023_arca_associated_voucher.up.sql
CREATE db/migrations/0023_arca_associated_voucher.down.sql
CREATE db/tests/0023_arca_associated_voucher.test.sql
```

## 5. Materialization blocks

### FILE: `internal/fiscal/service.go`
```yaml
block_id: "GO-ARCA-FISCAL:service:v1"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by ARCA WSFEv1 4.6 and pinned Microsoft BCApps evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "3b343c66cead67bb3bca420b6b3a8a203a683f6e85919ed5d720fc412cf44c10"
variables: []
secrets_allowed: false
```
````go
package fiscal

import (
	"context"
	"errors"
	"regexp"
	"time"
)

var ErrInvalid = errors.New("invalid fiscal command")
var ErrConflict = errors.New("fiscal conflict")
var ErrNotFound = errors.New("fiscal record not found")

var digitsPattern = regexp.MustCompile(`^[0-9]+$`)
var hashPattern = regexp.MustCompile(`^[a-f0-9]{64}$`)
var currencyPattern = regexp.MustCompile(`^[A-Z]{3}$`)

type PointOfSale struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	TaxpayerCUIT   string `json:"taxpayer_cuit"`
	Environment    string `json:"environment"`
	Number         int    `json:"number"`
	Active         bool   `json:"active"`
	Version        int64  `json:"version"`
}

type VATLine struct {
	ID               int   `json:"id"`
	BaseMinorUnits   int64 `json:"base_minor_units"`
	AmountMinorUnits int64 `json:"amount_minor_units"`
}

type OtherTaxLine struct {
	ID               int    `json:"id"`
	Description      string `json:"description"`
	BaseMinorUnits   int64  `json:"base_minor_units"`
	RateBasisPoints  int    `json:"rate_basis_points"`
	AmountMinorUnits int64  `json:"amount_minor_units"`
}

type AssociatedVoucher struct {
	InvoiceID    string    `json:"invoice_id"`
	TaxpayerCUIT string    `json:"taxpayer_cuit,omitempty"`
	VoucherType  int       `json:"voucher_type,omitempty"`
	PointOfSale  int       `json:"point_of_sale,omitempty"`
	Number       int64     `json:"number,omitempty"`
	IssuedOn     time.Time `json:"issued_on,omitempty"`
}

type Invoice struct {
	TenantID                string              `json:"-"`
	ID                      string              `json:"id"`
	OrganizationID          string              `json:"organization_id"`
	OrderID                 string              `json:"order_id"`
	PaymentAttemptID        string              `json:"payment_attempt_id"`
	PointOfSaleID           string              `json:"point_of_sale_id"`
	TaxpayerCUIT            string              `json:"taxpayer_cuit,omitempty"`
	Environment             string              `json:"environment,omitempty"`
	PointOfSaleNumber       int                 `json:"point_of_sale_number,omitempty"`
	VoucherType             int                 `json:"voucher_type"`
	Concept                 int                 `json:"concept"`
	RecipientDocumentType   int                 `json:"recipient_document_type"`
	RecipientDocument       string              `json:"recipient_document"`
	RecipientVATConditionID int                 `json:"recipient_vat_condition_id"`
	Currency                string              `json:"currency,omitempty"`
	TotalMinorUnits         int64               `json:"total_minor_units"`
	NetMinorUnits           int64               `json:"net_minor_units"`
	VATMinorUnits           int64               `json:"vat_minor_units"`
	ExemptMinorUnits        int64               `json:"exempt_minor_units"`
	NonTaxedMinorUnits      int64               `json:"non_taxed_minor_units"`
	OtherTaxMinorUnits      int64               `json:"other_tax_minor_units"`
	VATLines                []VATLine           `json:"vat_lines,omitempty"`
	OtherTaxLines           []OtherTaxLine      `json:"other_tax_lines,omitempty"`
	AssociatedVouchers      []AssociatedVoucher `json:"associated_vouchers,omitempty"`
	IssuedOn                time.Time           `json:"issued_on"`
	ServiceFrom             *time.Time          `json:"service_from,omitempty"`
	ServiceUntil            *time.Time          `json:"service_until,omitempty"`
	PaymentDueOn            *time.Time          `json:"payment_due_on,omitempty"`
	Status                  string              `json:"status"`
	VoucherNumber           int64               `json:"voucher_number,omitempty"`
	CAE                     string              `json:"cae,omitempty"`
	CAEExpiresOn            *time.Time          `json:"cae_expires_on,omitempty"`
	Version                 int64               `json:"version"`
	ClaimedFrom             string              `json:"-"`
}

type Authorization struct {
	Found         bool
	Authorized    bool
	CAE           string
	CAEExpiresOn  *time.Time
	ResponseHash  string
	ProviderCodes []string
}

type Repository interface {
	ConfigurePointOfSale(context.Context, string, string, PointOfSale) error
	RequestInvoice(context.Context, string, string, string, string, Invoice) (Invoice, bool, error)
	GetInvoice(context.Context, string, string, string) (Invoice, error)
}

type IDGenerator interface{ New() string }

type Service struct {
	repository Repository
	ids        IDGenerator
}

func NewService(repository Repository, ids IDGenerator) *Service {
	return &Service{repository: repository, ids: ids}
}

func (s *Service) ConfigurePointOfSale(ctx context.Context, tenant string, value PointOfSale) (PointOfSale, error) {
	if tenant == "" || value.OrganizationID == "" || !ValidCUIT(value.TaxpayerCUIT) || (value.Environment != "homologation" && value.Environment != "production") || value.Number < 1 || value.Number > 99999 {
		return value, ErrInvalid
	}
	value.ID = s.ids.New()
	value.Active = true
	value.Version = 1
	if err := s.repository.ConfigurePointOfSale(ctx, tenant, s.ids.New(), value); err != nil {
		return value, err
	}
	return value, nil
}

func (s *Service) RequestInvoice(ctx context.Context, tenant, idempotencyKey, requestHash string, value Invoice) (Invoice, bool, error) {
	if tenant == "" || value.OrganizationID == "" || value.OrderID == "" || value.PaymentAttemptID == "" || value.PointOfSaleID == "" || value.Concept < 1 || value.Concept > 3 || value.RecipientDocumentType < 0 || value.RecipientDocumentType > 999 || len(value.RecipientDocument) > 20 || (value.RecipientDocument != "" && !digitsPattern.MatchString(value.RecipientDocument)) || value.RecipientVATConditionID < 1 || value.RecipientVATConditionID > 999 || !validAmounts(value) || value.IssuedOn.IsZero() || len(idempotencyKey) < 16 || len(idempotencyKey) > 128 || !hashPattern.MatchString(requestHash) || !validServiceDates(value) || !validAssociationRequest(value) {
		return value, false, ErrInvalid
	}
	value.ID = s.ids.New()
	value.Status = "queued"
	value.Version = 1
	return s.repository.RequestInvoice(ctx, tenant, idempotencyKey, requestHash, s.ids.New(), value)
}

func AssociatedOriginalVoucherType(voucherType int) (int, bool) {
	switch voucherType {
	case 3:
		return 1, true
	case 8:
		return 6, true
	case 13:
		return 11, true
	default:
		return 0, false
	}
}

func SupportedInvoiceVoucherType(voucherType int) bool {
	return voucherType == 1 || voucherType == 6 || voucherType == 11
}

func validAssociationRequest(value Invoice) bool {
	_, credit := AssociatedOriginalVoucherType(value.VoucherType)
	if !credit {
		return SupportedInvoiceVoucherType(value.VoucherType) && len(value.AssociatedVouchers) == 0
	}
	if len(value.AssociatedVouchers) != 1 {
		return false
	}
	associated := value.AssociatedVouchers[0]
	return associated.InvoiceID != "" && associated.TaxpayerCUIT == "" && associated.VoucherType == 0 && associated.PointOfSale == 0 && associated.Number == 0 && associated.IssuedOn.IsZero()
}

func (s *Service) GetInvoice(ctx context.Context, tenant, organization, id string) (Invoice, error) {
	if tenant == "" || organization == "" || id == "" {
		return Invoice{}, ErrInvalid
	}
	return s.repository.GetInvoice(ctx, tenant, organization, id)
}

func validAmounts(value Invoice) bool {
	if value.NetMinorUnits < 0 || value.VATMinorUnits < 0 || value.ExemptMinorUnits < 0 || value.NonTaxedMinorUnits < 0 || value.OtherTaxMinorUnits < 0 {
		return false
	}
	total := value.NetMinorUnits + value.VATMinorUnits + value.ExemptMinorUnits + value.NonTaxedMinorUnits + value.OtherTaxMinorUnits
	if total <= 0 || (value.TotalMinorUnits != 0 && value.TotalMinorUnits != total) {
		return false
	}
	seenVAT := map[int]struct{}{}
	var vatBase, vatAmount int64
	for _, line := range value.VATLines {
		if line.ID < 1 || line.ID > 999 || line.BaseMinorUnits <= 0 || line.AmountMinorUnits < 0 {
			return false
		}
		if _, ok := seenVAT[line.ID]; ok {
			return false
		}
		seenVAT[line.ID] = struct{}{}
		vatBase += line.BaseMinorUnits
		vatAmount += line.AmountMinorUnits
	}
	if vatAmount != value.VATMinorUnits || (value.VATMinorUnits > 0 && (len(value.VATLines) == 0 || vatBase != value.NetMinorUnits)) || (value.VATMinorUnits == 0 && len(value.VATLines) != 0) {
		return false
	}
	seenTax := map[int]struct{}{}
	var otherAmount int64
	for _, line := range value.OtherTaxLines {
		if line.ID < 1 || line.ID > 999 || len(line.Description) < 1 || len(line.Description) > 80 || line.BaseMinorUnits < 0 || line.RateBasisPoints < 0 || line.RateBasisPoints > 100000 || line.AmountMinorUnits < 0 {
			return false
		}
		if _, ok := seenTax[line.ID]; ok {
			return false
		}
		seenTax[line.ID] = struct{}{}
		otherAmount += line.AmountMinorUnits
	}
	return otherAmount == value.OtherTaxMinorUnits && ((value.OtherTaxMinorUnits == 0 && len(value.OtherTaxLines) == 0) || (value.OtherTaxMinorUnits > 0 && len(value.OtherTaxLines) > 0))
}

func validServiceDates(value Invoice) bool {
	if value.Concept == 1 {
		return value.ServiceFrom == nil && value.ServiceUntil == nil && value.PaymentDueOn == nil
	}
	return value.ServiceFrom != nil && value.ServiceUntil != nil && value.PaymentDueOn != nil && !value.ServiceUntil.Before(*value.ServiceFrom) && !value.PaymentDueOn.Before(*value.ServiceUntil)
}

func ValidCUIT(value string) bool {
	if len(value) != 11 || !digitsPattern.MatchString(value) {
		return false
	}
	weights := [...]int{5, 4, 3, 2, 7, 6, 5, 4, 3, 2}
	sum := 0
	for index, weight := range weights {
		sum += int(value[index]-'0') * weight
	}
	check := 11 - sum%11
	if check == 11 {
		check = 0
	} else if check == 10 {
		check = 9
	}
	return check == int(value[10]-'0')
}

func validAuthorization(value Authorization) bool {
	if !hashPattern.MatchString(value.ResponseHash) {
		return false
	}
	if !value.Found || !value.Authorized {
		return value.CAE == "" && value.CAEExpiresOn == nil
	}
	return len(value.CAE) >= 8 && len(value.CAE) <= 20 && digitsPattern.MatchString(value.CAE) && value.CAEExpiresOn != nil && !value.CAEExpiresOn.IsZero()
}

func ValidProviderCurrency(value string) bool { return currencyPattern.MatchString(value) }
````

### FILE: `internal/fiscal/service_test.go`
```yaml
block_id: "GO-ARCA-FISCAL:service-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "53c5d2d25924493b730d7819cdfbff1837f18bf287529bf5b776b091a752bff9"
variables: []
secrets_allowed: false
```
````go
package fiscal

import (
	"context"
	"errors"
	"testing"
	"time"
)

type serviceRepo struct{ requested int }

func (*serviceRepo) ConfigurePointOfSale(context.Context, string, string, PointOfSale) error {
	return nil
}
func (r *serviceRepo) RequestInvoice(_ context.Context, _, _, _, _ string, value Invoice) (Invoice, bool, error) {
	r.requested++
	value.TotalMinorUnits = value.NetMinorUnits + value.VATMinorUnits + value.ExemptMinorUnits + value.NonTaxedMinorUnits + value.OtherTaxMinorUnits
	return value, false, nil
}
func (*serviceRepo) GetInvoice(context.Context, string, string, string) (Invoice, error) {
	return Invoice{ID: "invoice"}, nil
}

type testIDs struct{ n int }

func (i *testIDs) New() string { i.n++; return "id" + string(rune('0'+i.n)) }

func TestValidCUITAndInvoiceBoundary(t *testing.T) {
	if !ValidCUIT("30715117564") || ValidCUIT("30715117565") || ValidCUIT("not-a-cuit") {
		t.Fatal("CUIT checksum contract failed")
	}
	repo := &serviceRepo{}
	service := NewService(repo, &testIDs{})
	base := Invoice{OrganizationID: "franchise", OrderID: "order", PaymentAttemptID: "payment", PointOfSaleID: "pos", VoucherType: 6, Concept: 1, RecipientDocumentType: 99, RecipientDocument: "0", RecipientVATConditionID: 5, NetMinorUnits: 8264, VATMinorUnits: 1736, VATLines: []VATLine{{ID: 5, BaseMinorUnits: 8264, AmountMinorUnits: 1736}}, IssuedOn: time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC)}
	if _, _, err := service.RequestInvoice(context.Background(), "tenant", "1234567890abcdef", string(make([]byte, 64)), base); !errors.Is(err, ErrInvalid) {
		t.Fatalf("non-hex hash accepted: %v", err)
	}
	value, replay, err := service.RequestInvoice(context.Background(), "tenant", "1234567890abcdef", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", base)
	if err != nil || replay || value.TotalMinorUnits != 10000 || repo.requested != 1 {
		t.Fatalf("value=%+v replay=%v err=%v count=%d", value, replay, err, repo.requested)
	}
	base.Concept = 2
	if _, _, err = service.RequestInvoice(context.Background(), "tenant", "1234567890abcdef", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", base); !errors.Is(err, ErrInvalid) {
		t.Fatalf("service dates omission accepted: %v", err)
	}
	base.Concept = 1
	base.VATLines[0].AmountMinorUnits = 1735
	if _, _, err = service.RequestInvoice(context.Background(), "tenant", "1234567890abcdef", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", base); !errors.Is(err, ErrInvalid) {
		t.Fatalf("VAT mismatch accepted: %v", err)
	}
	base.VATLines[0].AmountMinorUnits = 1736
	base.RecipientVATConditionID = 0
	if _, _, err = service.RequestInvoice(context.Background(), "tenant", "1234567890abcdef", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", base); !errors.Is(err, ErrInvalid) {
		t.Fatalf("missing recipient VAT condition accepted: %v", err)
	}
	base.RecipientVATConditionID = 5
	base.VoucherType = 8
	if _, _, err = service.RequestInvoice(context.Background(), "tenant", "1234567890abcdef", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", base); !errors.Is(err, ErrInvalid) {
		t.Fatalf("credit without associated invoice accepted: %v", err)
	}
	base.AssociatedVouchers = []AssociatedVoucher{{InvoiceID: "original-invoice"}}
	if _, _, err = service.RequestInvoice(context.Background(), "tenant", "1234567890abcdef", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", base); err != nil {
		t.Fatalf("valid credit association rejected: %v", err)
	}
	base.AssociatedVouchers[0].Number = 7
	if _, _, err = service.RequestInvoice(context.Background(), "tenant", "1234567890abcdef", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", base); !errors.Is(err, ErrInvalid) {
		t.Fatalf("client-supplied associated identity accepted: %v", err)
	}
}
````

### FILE: `internal/fiscal/processor.go`
```yaml
block_id: "GO-ARCA-FISCAL:processor:v1"
operation: CREATE
provenance: AUTHORED
source: "local implementation of ARCA documented timeout reconciliation"
license: "LicenseRef-Workspace-Owner"
sha256: "77fe23edacdc2235e8a521e8447feee16336f612cd0315224f58a21c20a4a313"
variables: []
secrets_allowed: false
```
````go
package fiscal

import (
	"context"
	"errors"
	"time"
)

var ErrNoWork = errors.New("no fiscal work available")

type Provider interface {
	LastAuthorized(context.Context, Invoice) (int64, string, error)
	Consult(context.Context, Invoice) (Authorization, error)
	Authorize(context.Context, Invoice) (Authorization, error)
}

type WorkRepository interface {
	ClaimInvoice(context.Context, string, time.Duration, string) (Invoice, error)
	AssignVoucherNumber(context.Context, Invoice, string, int64, string, string) (Invoice, error)
	FinishInvoice(context.Context, Invoice, string, Authorization, string) (Invoice, error)
	DeferInvoice(context.Context, Invoice, string, string, string, bool) error
}

type Processor struct {
	repository WorkRepository
	provider   Provider
	ids        IDGenerator
	workerID   string
	lease      time.Duration
}

func NewProcessor(repository WorkRepository, provider Provider, ids IDGenerator, workerID string, lease time.Duration) (*Processor, error) {
	if repository == nil || provider == nil || ids == nil || workerID == "" || lease < 30*time.Second || lease > 10*time.Minute {
		return nil, ErrInvalid
	}
	return &Processor{repository: repository, provider: provider, ids: ids, workerID: workerID, lease: lease}, nil
}

func (p *Processor) ProcessOne(ctx context.Context) (Invoice, error) {
	invoice, err := p.repository.ClaimInvoice(ctx, p.workerID, p.lease, p.ids.New())
	if err != nil {
		return invoice, err
	}
	if invoice.VoucherNumber > 0 {
		consulted, consultErr := p.provider.Consult(ctx, invoice)
		if consultErr != nil {
			_ = p.repository.DeferInvoice(ctx, invoice, p.workerID, "consult", "", true)
			return invoice, consultErr
		}
		if !validAuthorization(consulted) {
			_ = p.repository.DeferInvoice(ctx, invoice, p.workerID, "consult-invalid", consulted.ResponseHash, true)
			return invoice, ErrInvalid
		}
		if consulted.Found {
			return p.repository.FinishInvoice(ctx, invoice, p.workerID, consulted, p.ids.New())
		}
	}
	if invoice.VoucherNumber == 0 {
		last, responseHash, sequenceErr := p.provider.LastAuthorized(ctx, invoice)
		if sequenceErr != nil || last < 0 || !hashPattern.MatchString(responseHash) {
			_ = p.repository.DeferInvoice(ctx, invoice, p.workerID, "last-authorized", responseHash, false)
			if sequenceErr != nil {
				return invoice, sequenceErr
			}
			return invoice, ErrInvalid
		}
		invoice, err = p.repository.AssignVoucherNumber(ctx, invoice, p.workerID, last+1, responseHash, p.ids.New())
		if err != nil {
			return invoice, err
		}
	}
	authorized, authorizationErr := p.provider.Authorize(ctx, invoice)
	if authorizationErr != nil {
		_ = p.repository.DeferInvoice(ctx, invoice, p.workerID, "authorize", "", true)
		return invoice, authorizationErr
	}
	if !validAuthorization(authorized) || !authorized.Found {
		_ = p.repository.DeferInvoice(ctx, invoice, p.workerID, "authorize-invalid", authorized.ResponseHash, true)
		return invoice, ErrInvalid
	}
	return p.repository.FinishInvoice(ctx, invoice, p.workerID, authorized, p.ids.New())
}
````

### FILE: `internal/fiscal/processor_test.go`
```yaml
block_id: "GO-ARCA-FISCAL:processor-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c09eeeaefa54177bdb6e6c4912ac6501210feb315bd95b237f1ada7c2843a7f5"
variables: []
secrets_allowed: false
```
````go
package fiscal

import (
	"context"
	"errors"
	"testing"
	"time"
)

type workRepo struct {
	invoice  Invoice
	assigned int64
	finished int
	deferred string
}

func (r *workRepo) ClaimInvoice(context.Context, string, time.Duration, string) (Invoice, error) {
	return r.invoice, nil
}
func (r *workRepo) AssignVoucherNumber(_ context.Context, value Invoice, _ string, number int64, _, _ string) (Invoice, error) {
	r.assigned = number
	value.VoucherNumber = number
	value.Status = "authorizing"
	return value, nil
}
func (r *workRepo) FinishInvoice(_ context.Context, value Invoice, _ string, result Authorization, _ string) (Invoice, error) {
	r.finished++
	if result.Authorized {
		value.Status = "authorized"
		value.CAE = result.CAE
	} else {
		value.Status = "rejected"
	}
	return value, nil
}
func (r *workRepo) DeferInvoice(_ context.Context, _ Invoice, _ string, phase, _ string, _ bool) error {
	r.deferred = phase
	return nil
}

type providerFake struct {
	last         int64
	consult      Authorization
	authorize    Authorization
	consultErr   error
	authorizeErr error
}

func (p providerFake) LastAuthorized(context.Context, Invoice) (int64, string, error) {
	return p.last, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", nil
}
func (p providerFake) Consult(context.Context, Invoice) (Authorization, error) {
	return p.consult, p.consultErr
}
func (p providerFake) Authorize(context.Context, Invoice) (Authorization, error) {
	return p.authorize, p.authorizeErr
}

func TestProcessorSequencesAuthorizesAndReconcilesBeforeRetry(t *testing.T) {
	expires := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)
	ok := Authorization{Found: true, Authorized: true, CAE: "12345678901234", CAEExpiresOn: &expires, ResponseHash: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"}
	repo := &workRepo{invoice: Invoice{ID: "invoice", Status: "queued"}}
	processor, err := NewProcessor(repo, providerFake{last: 40, authorize: ok}, &testIDs{}, "worker", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	value, err := processor.ProcessOne(context.Background())
	if err != nil || repo.assigned != 41 || repo.finished != 1 || value.Status != "authorized" {
		t.Fatalf("value=%+v assigned=%d finished=%d err=%v", value, repo.assigned, repo.finished, err)
	}
	repo = &workRepo{invoice: Invoice{ID: "ambiguous", Status: "reconcile_required", VoucherNumber: 41}}
	processor, _ = NewProcessor(repo, providerFake{consult: ok, authorizeErr: errors.New("must not authorize")}, &testIDs{}, "worker", time.Minute)
	value, err = processor.ProcessOne(context.Background())
	if err != nil || repo.finished != 1 || value.CAE != ok.CAE {
		t.Fatalf("reconciliation failed: value=%+v err=%v", value, err)
	}
	repo = &workRepo{invoice: Invoice{ID: "timeout", Status: "queued"}}
	processor, _ = NewProcessor(repo, providerFake{last: 9, authorizeErr: errors.New("timeout")}, &testIDs{}, "worker", time.Minute)
	if _, err = processor.ProcessOne(context.Background()); err == nil || repo.deferred != "authorize" || repo.assigned != 10 {
		t.Fatalf("ambiguous authorization was not deferred: assigned=%d deferred=%s err=%v", repo.assigned, repo.deferred, err)
	}
}
````

### FILE: `internal/platform/postgres/fiscal.go`
```yaml
block_id: "GO-ARCA-FISCAL:postgres:v1"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by PostgreSQL 18 transaction and locking documentation"
license: "LicenseRef-Workspace-Owner"
sha256: "fd4777085eaf06802543fdc86e7761768edad81bbeec0c62337889342573627e"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"elite.local/enterprise/internal/fiscal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Fiscal struct{ pool *pgxpool.Pool }

func NewFiscal(pool *pgxpool.Pool) *Fiscal { return &Fiscal{pool: pool} }

func fiscalConflict(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23503" || pgErr.Code == "23505" || pgErr.Code == "23514" || pgErr.Code == "40001") {
		return fiscal.ErrConflict
	}
	return err
}

func (r *Fiscal) ConfigurePointOfSale(ctx context.Context, tenant, eventID string, value fiscal.PointOfSale) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into fiscal.point_of_sale(tenant_id,point_of_sale_id,organization_id,taxpayer_cuit,environment,point_of_sale_number,active,version) select $1,$2,$3,$4,$5,$6,true,1 where exists(select 1 from org.organization where tenant_id=$1 and organization_id=$3 and status='active')`, tenant, value.ID, value.OrganizationID, value.TaxpayerCUIT, value.Environment, value.Number)
	if err != nil {
		return fiscalConflict(err)
	}
	if result.RowsAffected() != 1 {
		return fiscal.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'fiscal-point-of-sale',$3,1,'fiscal-point-of-sale.configured',1,clock_timestamp(),jsonb_build_object('organization_id',$4::text,'environment',$5::text,'point_of_sale_number',$6::integer))`, tenant, eventID, value.ID, value.OrganizationID, value.Environment, value.Number)
	if err != nil {
		return fiscalConflict(err)
	}
	return fiscalConflict(tx.Commit(ctx))
}

func (r *Fiscal) RequestInvoice(ctx context.Context, tenant, idempotency, requestHash, eventID string, value fiscal.Invoice) (fiscal.Invoice, bool, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return value, false, err
	}
	defer tx.Rollback(ctx)
	var existingHash string
	err = tx.QueryRow(ctx, `select request_hash from fiscal.invoice where tenant_id=$1 and idempotency_key=$2`, tenant, idempotency).Scan(&existingHash)
	if err == nil {
		if existingHash != requestHash {
			return value, false, fiscal.ErrConflict
		}
		existing, loadErr := scanFiscalInvoice(tx.QueryRow(ctx, invoiceSelect+` where i.tenant_id=$1 and i.idempotency_key=$2`, tenant, idempotency))
		if loadErr == nil {
			loadErr = loadFiscalLines(ctx, tx, &existing)
		}
		return existing, true, loadErr
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return value, false, err
	}
	originalType, credit := fiscal.AssociatedOriginalVoucherType(value.VoucherType)
	if (!credit && (!fiscal.SupportedInvoiceVoucherType(value.VoucherType) || len(value.AssociatedVouchers) != 0)) || (credit && (len(value.AssociatedVouchers) != 1 || value.AssociatedVouchers[0].InvoiceID == "")) {
		return value, false, fiscal.ErrConflict
	}
	if credit {
		associated := &value.AssociatedVouchers[0]
		err = tx.QueryRow(ctx, `select original_pos.taxpayer_cuit,original.voucher_type,original_pos.point_of_sale_number,original.voucher_number,original.issued_on from fiscal.invoice original join fiscal.point_of_sale original_pos on original_pos.tenant_id=original.tenant_id and original_pos.point_of_sale_id=original.point_of_sale_id join fiscal.point_of_sale credit_pos on credit_pos.tenant_id=original.tenant_id and credit_pos.point_of_sale_id=$5 and credit_pos.organization_id=original.organization_id and credit_pos.taxpayer_cuit=original_pos.taxpayer_cuit and credit_pos.active where original.tenant_id=$1 and original.organization_id=$2 and original.invoice_id=$3 and original.status='authorized' and original.voucher_type=$4 and original.voucher_number is not null for key share of original,original_pos,credit_pos`, tenant, value.OrganizationID, associated.InvoiceID, originalType, value.PointOfSaleID).Scan(&associated.TaxpayerCUIT, &associated.VoucherType, &associated.PointOfSale, &associated.Number, &associated.IssuedOn)
		if errors.Is(err, pgx.ErrNoRows) {
			return value, false, fiscal.ErrConflict
		}
		if err != nil {
			return value, false, err
		}
		if associated.IssuedOn.After(value.IssuedOn) {
			return value, false, fiscal.ErrConflict
		}
	}
	var orderCurrency, paymentCurrency string
	var orderTotal, paymentTotal int64
	err = tx.QueryRow(ctx, `select o.currency,o.total_minor_units,p.currency,p.amount_minor_units from sales.customer_order o join payment.payment_attempt p on p.tenant_id=o.tenant_id and p.order_id=o.order_id where o.tenant_id=$1 and o.order_id=$2 and o.organization_id=$3 and o.state<>'cancelled' and p.payment_attempt_id=$4 and p.state='captured' for update of o,p`, tenant, value.OrderID, value.OrganizationID, value.PaymentAttemptID).Scan(&orderCurrency, &orderTotal, &paymentCurrency, &paymentTotal)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, false, fiscal.ErrConflict
	}
	if err != nil {
		return value, false, err
	}
	total := value.NetMinorUnits + value.VATMinorUnits + value.ExemptMinorUnits + value.NonTaxedMinorUnits + value.OtherTaxMinorUnits
	if orderCurrency != paymentCurrency || orderCurrency != "ARS" || orderTotal != paymentTotal || paymentTotal != total || !fiscalDetailsValid(value) {
		return value, false, fiscal.ErrConflict
	}
	result, err := tx.Exec(ctx, `insert into fiscal.invoice(tenant_id,invoice_id,organization_id,order_id,payment_attempt_id,point_of_sale_id,voucher_type,concept,recipient_document_type,recipient_document,recipient_vat_condition_id,currency,provider_currency,total_minor_units,net_minor_units,vat_minor_units,exempt_minor_units,non_taxed_minor_units,other_tax_minor_units,issued_on,service_from,service_until,payment_due_on,status,request_hash,idempotency_key,version) select $1,$2,$3,$4,$5,p.point_of_sale_id,$6,$7,$8,$9,$10,$11,'PES',$12,$13,$14,$15,$16,$17,($18::timestamptz at time zone 'UTC')::date,($19::timestamptz at time zone 'UTC')::date,($20::timestamptz at time zone 'UTC')::date,($21::timestamptz at time zone 'UTC')::date,'queued',$22,$23,1 from fiscal.point_of_sale p where p.tenant_id=$1 and p.point_of_sale_id=$24 and p.organization_id=$3 and p.active`, tenant, value.ID, value.OrganizationID, value.OrderID, value.PaymentAttemptID, value.VoucherType, value.Concept, value.RecipientDocumentType, value.RecipientDocument, value.RecipientVATConditionID, orderCurrency, total, value.NetMinorUnits, value.VATMinorUnits, value.ExemptMinorUnits, value.NonTaxedMinorUnits, value.OtherTaxMinorUnits, value.IssuedOn, value.ServiceFrom, value.ServiceUntil, value.PaymentDueOn, requestHash, idempotency, value.PointOfSaleID)
	if err != nil {
		return value, false, fiscalConflict(err)
	}
	if result.RowsAffected() != 1 {
		return value, false, fiscal.ErrConflict
	}
	if credit {
		associated := value.AssociatedVouchers[0]
		if _, err = tx.Exec(ctx, `insert into fiscal.invoice_associated_voucher(tenant_id,invoice_id,associated_invoice_id,taxpayer_cuit,voucher_type,point_of_sale_number,voucher_number,issued_on) values($1,$2,$3,$4,$5,$6,$7,($8::timestamptz at time zone 'UTC')::date)`, tenant, value.ID, associated.InvoiceID, associated.TaxpayerCUIT, associated.VoucherType, associated.PointOfSale, associated.Number, associated.IssuedOn); err != nil {
			return value, false, fiscalConflict(err)
		}
	}
	for _, line := range value.VATLines {
		if _, err = tx.Exec(ctx, `insert into fiscal.invoice_vat(tenant_id,invoice_id,vat_id,base_minor_units,amount_minor_units) values($1,$2,$3,$4,$5)`, tenant, value.ID, line.ID, line.BaseMinorUnits, line.AmountMinorUnits); err != nil {
			return value, false, fiscalConflict(err)
		}
	}
	for _, line := range value.OtherTaxLines {
		if _, err = tx.Exec(ctx, `insert into fiscal.invoice_other_tax(tenant_id,invoice_id,tax_id,description,base_minor_units,rate_basis_points,amount_minor_units) values($1,$2,$3,$4,$5,$6,$7)`, tenant, value.ID, line.ID, line.Description, line.BaseMinorUnits, line.RateBasisPoints, line.AmountMinorUnits); err != nil {
			return value, false, fiscalConflict(err)
		}
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'fiscal-invoice',$3,1,'fiscal-invoice.requested',1,clock_timestamp(),jsonb_build_object('organization_id',$4::text,'order_id',$5::text,'payment_attempt_id',$6::text,'point_of_sale_id',$7::text,'voucher_type',$8::integer,'total_minor_units',$9::bigint,'currency',$10::text))`, tenant, eventID, value.ID, value.OrganizationID, value.OrderID, value.PaymentAttemptID, value.PointOfSaleID, value.VoucherType, total, orderCurrency)
	if err != nil {
		return value, false, fiscalConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return value, false, fiscalConflict(err)
	}
	value.TenantID = tenant
	value.TotalMinorUnits = total
	value.Currency = orderCurrency
	value.Status = "queued"
	value.Version = 1
	return value, false, nil
}

const invoiceSelect = `select i.tenant_id,i.invoice_id,i.organization_id,i.order_id,i.payment_attempt_id,i.point_of_sale_id,p.taxpayer_cuit,p.environment,p.point_of_sale_number,i.voucher_type,i.concept,i.recipient_document_type,i.recipient_document,i.recipient_vat_condition_id,i.currency,i.total_minor_units,i.net_minor_units,i.vat_minor_units,i.exempt_minor_units,i.non_taxed_minor_units,i.other_tax_minor_units,i.issued_on,i.service_from,i.service_until,i.payment_due_on,i.status,coalesce(i.voucher_number,0),coalesce(i.cae,''),i.cae_expires_on,i.version from fiscal.invoice i join fiscal.point_of_sale p on p.tenant_id=i.tenant_id and p.point_of_sale_id=i.point_of_sale_id`

type fiscalRow interface{ Scan(...any) error }

func scanFiscalInvoice(row fiscalRow) (fiscal.Invoice, error) {
	var v fiscal.Invoice
	err := row.Scan(&v.TenantID, &v.ID, &v.OrganizationID, &v.OrderID, &v.PaymentAttemptID, &v.PointOfSaleID, &v.TaxpayerCUIT, &v.Environment, &v.PointOfSaleNumber, &v.VoucherType, &v.Concept, &v.RecipientDocumentType, &v.RecipientDocument, &v.RecipientVATConditionID, &v.Currency, &v.TotalMinorUnits, &v.NetMinorUnits, &v.VATMinorUnits, &v.ExemptMinorUnits, &v.NonTaxedMinorUnits, &v.OtherTaxMinorUnits, &v.IssuedOn, &v.ServiceFrom, &v.ServiceUntil, &v.PaymentDueOn, &v.Status, &v.VoucherNumber, &v.CAE, &v.CAEExpiresOn, &v.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, fiscal.ErrNotFound
	}
	return v, err
}
func (r *Fiscal) GetInvoice(ctx context.Context, tenant, organization, id string) (fiscal.Invoice, error) {
	v, err := scanFiscalInvoice(r.pool.QueryRow(ctx, invoiceSelect+` where i.tenant_id=$1 and i.organization_id=$2 and i.invoice_id=$3`, tenant, organization, id))
	if err != nil {
		return v, err
	}
	err = loadFiscalLines(ctx, r.pool, &v)
	return v, err
}

func (r *Fiscal) ClaimInvoice(ctx context.Context, worker string, lease time.Duration, attemptID string) (fiscal.Invoice, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fiscal.Invoice{}, err
	}
	defer tx.Rollback(ctx)
	query := invoiceSelect + ` where i.status in ('queued','reconcile_required') and p.active and (p.active_invoice_id is null or p.lease_until < clock_timestamp()) order by i.created_at,i.tenant_id,i.invoice_id for update of i,p skip locked limit 1`
	v, err := scanFiscalInvoice(tx.QueryRow(ctx, query))
	if errors.Is(err, fiscal.ErrNotFound) {
		return v, fiscal.ErrNoWork
	}
	if err != nil {
		return v, err
	}
	if err = loadFiscalLines(ctx, tx, &v); err != nil {
		return v, err
	}
	v.ClaimedFrom = v.Status
	result, err := tx.Exec(ctx, `update fiscal.point_of_sale set active_invoice_id=$3,lease_owner=$4,lease_until=clock_timestamp()+$5::interval where tenant_id=$1 and point_of_sale_id=$2 and (active_invoice_id is null or lease_until<clock_timestamp())`, v.TenantID, v.PointOfSaleID, v.ID, worker, lease.String())
	if err != nil {
		return v, err
	}
	if result.RowsAffected() != 1 {
		return v, fiscal.ErrConflict
	}
	err = tx.QueryRow(ctx, `update fiscal.invoice set attempt_count=attempt_count+1,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and invoice_id=$2 and point_of_sale_id=$3 returning version`, v.TenantID, v.ID, v.PointOfSaleID).Scan(&v.Version)
	if err != nil {
		return v, err
	}
	_, err = tx.Exec(ctx, `insert into fiscal.issuance_attempt(tenant_id,attempt_id,invoice_id,phase,voucher_number,request_hash,outcome) select tenant_id,$2,invoice_id,'claim',voucher_number,request_hash,'started' from fiscal.invoice where tenant_id=$1 and invoice_id=$3`, v.TenantID, attemptID, v.ID)
	if err != nil {
		return v, fiscalConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return v, fiscalConflict(err)
	}
	return v, nil
}

func (r *Fiscal) AssignVoucherNumber(ctx context.Context, value fiscal.Invoice, worker string, number int64, responseHash, eventID string) (fiscal.Invoice, error) {
	if number < 1 {
		return value, fiscal.ErrInvalid
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	var tenant, requestHash string
	err = tx.QueryRow(ctx, `select i.tenant_id,i.request_hash from fiscal.invoice i join fiscal.point_of_sale p on p.tenant_id=i.tenant_id and p.point_of_sale_id=i.point_of_sale_id where i.tenant_id=$5 and i.invoice_id=$1 and i.point_of_sale_id=$2 and i.status='queued' and i.voucher_number is null and i.version=$3 and p.active_invoice_id=i.invoice_id and p.lease_owner=$4 and p.lease_until>=clock_timestamp() for update of i,p`, value.ID, value.PointOfSaleID, value.Version, worker, value.TenantID).Scan(&tenant, &requestHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, fiscal.ErrConflict
	}
	if err != nil {
		return value, err
	}
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":"+value.PointOfSaleID+":"+string(rune(value.VoucherType))); err != nil {
		return value, err
	}
	err = tx.QueryRow(ctx, `update fiscal.invoice set voucher_number=$3,status='authorizing',last_response_hash=$4,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and invoice_id=$2 returning version,status`, tenant, value.ID, number, responseHash).Scan(&value.Version, &value.Status)
	if err != nil {
		return value, fiscalConflict(err)
	}
	value.VoucherNumber = number
	_, err = tx.Exec(ctx, `insert into fiscal.issuance_attempt(tenant_id,attempt_id,invoice_id,phase,voucher_number,request_hash,response_hash,outcome) values($1,$2,$3,'assign-number',$4,$5,$6,'observed')`, tenant, eventID, value.ID, number, requestHash, responseHash)
	if err != nil {
		return value, fiscalConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return value, fiscalConflict(err)
	}
	return value, nil
}

func (r *Fiscal) FinishInvoice(ctx context.Context, value fiscal.Invoice, worker string, authorization fiscal.Authorization, eventID string) (fiscal.Invoice, error) {
	codes, err := json.Marshal(append([]string{}, authorization.ProviderCodes...))
	if err != nil {
		return value, fiscal.ErrInvalid
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	status := "rejected"
	outcome := "rejected"
	if authorization.Authorized {
		status = "authorized"
		outcome = "authorized"
	}
	var tenant, requestHash string
	err = tx.QueryRow(ctx, `select i.tenant_id,i.request_hash from fiscal.invoice i join fiscal.point_of_sale p on p.tenant_id=i.tenant_id and p.point_of_sale_id=i.point_of_sale_id where i.tenant_id=$5 and i.invoice_id=$1 and i.point_of_sale_id=$2 and i.status in ('authorizing','reconcile_required') and i.voucher_number=$3 and p.active_invoice_id=i.invoice_id and p.lease_owner=$4 and p.lease_until>=clock_timestamp() for update of i,p`, value.ID, value.PointOfSaleID, value.VoucherNumber, worker, value.TenantID).Scan(&tenant, &requestHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, fiscal.ErrConflict
	}
	if err != nil {
		return value, err
	}
	err = tx.QueryRow(ctx, `update fiscal.invoice set status=$3,cae=$4,cae_expires_on=($5::timestamptz at time zone 'UTC')::date,last_response_hash=$6,provider_codes=$7,authorized_at=case when $3='authorized' then clock_timestamp() else null end,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and invoice_id=$2 returning version,status,coalesce(cae,''),cae_expires_on`, tenant, value.ID, status, nullIfEmpty(authorization.CAE), authorization.CAEExpiresOn, authorization.ResponseHash, codes).Scan(&value.Version, &value.Status, &value.CAE, &value.CAEExpiresOn)
	if err != nil {
		return value, fiscalConflict(err)
	}
	_, err = tx.Exec(ctx, `update fiscal.point_of_sale set active_invoice_id=null,lease_owner=null,lease_until=null where tenant_id=$1 and point_of_sale_id=$2 and active_invoice_id=$3 and lease_owner=$4`, tenant, value.PointOfSaleID, value.ID, worker)
	if err != nil {
		return value, err
	}
	_, err = tx.Exec(ctx, `insert into fiscal.issuance_attempt(tenant_id,attempt_id,invoice_id,phase,voucher_number,request_hash,response_hash,outcome,provider_codes) values($1,$2,$3,'finish',$4,$5,$6,$7,$8)`, tenant, eventID, value.ID, value.VoucherNumber, requestHash, authorization.ResponseHash, outcome, codes)
	if err != nil {
		return value, fiscalConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,gen_random_uuid(),'fiscal-invoice',$2,$3,$4,1,clock_timestamp(),jsonb_build_object('organization_id',$5::text,'order_id',$6::text,'voucher_number',$7::bigint,'voucher_type',$8::integer,'cae',$9::text))`, tenant, value.ID, value.Version, "fiscal-invoice."+status, value.OrganizationID, value.OrderID, value.VoucherNumber, value.VoucherType, value.CAE)
	if err != nil {
		return value, fiscalConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return value, fiscalConflict(err)
	}
	return value, nil
}

func (r *Fiscal) DeferInvoice(ctx context.Context, value fiscal.Invoice, worker, phase, responseHash string, ambiguous bool) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	status := "queued"
	if ambiguous || value.VoucherNumber > 0 {
		status = "reconcile_required"
	}
	var tenant, requestHash string
	err = tx.QueryRow(ctx, `select i.tenant_id,i.request_hash from fiscal.invoice i join fiscal.point_of_sale p on p.tenant_id=i.tenant_id and p.point_of_sale_id=i.point_of_sale_id where i.tenant_id=$3 and i.invoice_id=$1 and p.active_invoice_id=i.invoice_id and p.lease_owner=$2 for update of i,p`, value.ID, worker, value.TenantID).Scan(&tenant, &requestHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return fiscal.ErrConflict
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `update fiscal.invoice set status=$3,last_response_hash=case when $4='' then last_response_hash else $4 end,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and invoice_id=$2`, tenant, value.ID, status, responseHash)
	if err != nil {
		return fiscalConflict(err)
	}
	_, err = tx.Exec(ctx, `update fiscal.point_of_sale set active_invoice_id=null,lease_owner=null,lease_until=null where tenant_id=$1 and point_of_sale_id=$2 and active_invoice_id=$3 and lease_owner=$4`, tenant, value.PointOfSaleID, value.ID, worker)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into fiscal.issuance_attempt(tenant_id,attempt_id,invoice_id,phase,voucher_number,request_hash,response_hash,outcome,error_class) values($1,gen_random_uuid(),$2,'defer',$3,$4,nullif($5,''),'deferred',$6)`, tenant, value.ID, nullIfZero(value.VoucherNumber), requestHash, responseHash, phase)
	if err != nil {
		return fiscalConflict(err)
	}
	return fiscalConflict(tx.Commit(ctx))
}

func nullIfEmpty(value string) any {
	if value == "" {
		return nil
	}
	return value
}
func nullIfZero(value int64) any {
	if value == 0 {
		return nil
	}
	return value
}

type fiscalQueryer interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

func loadFiscalLines(ctx context.Context, q fiscalQueryer, value *fiscal.Invoice) error {
	rows, err := q.Query(ctx, `select vat_id,base_minor_units,amount_minor_units from fiscal.invoice_vat where tenant_id=$1 and invoice_id=$2 order by vat_id`, value.TenantID, value.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	value.VATLines = []fiscal.VATLine{}
	for rows.Next() {
		var line fiscal.VATLine
		if err = rows.Scan(&line.ID, &line.BaseMinorUnits, &line.AmountMinorUnits); err != nil {
			return err
		}
		value.VATLines = append(value.VATLines, line)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	rows.Close()
	rows, err = q.Query(ctx, `select tax_id,description,base_minor_units,rate_basis_points,amount_minor_units from fiscal.invoice_other_tax where tenant_id=$1 and invoice_id=$2 order by tax_id`, value.TenantID, value.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	value.OtherTaxLines = []fiscal.OtherTaxLine{}
	for rows.Next() {
		var line fiscal.OtherTaxLine
		if err = rows.Scan(&line.ID, &line.Description, &line.BaseMinorUnits, &line.RateBasisPoints, &line.AmountMinorUnits); err != nil {
			return err
		}
		value.OtherTaxLines = append(value.OtherTaxLines, line)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	rows.Close()
	rows, err = q.Query(ctx, `select associated_invoice_id,taxpayer_cuit,voucher_type,point_of_sale_number,voucher_number,issued_on from fiscal.invoice_associated_voucher where tenant_id=$1 and invoice_id=$2`, value.TenantID, value.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	value.AssociatedVouchers = []fiscal.AssociatedVoucher{}
	for rows.Next() {
		var associated fiscal.AssociatedVoucher
		if err = rows.Scan(&associated.InvoiceID, &associated.TaxpayerCUIT, &associated.VoucherType, &associated.PointOfSale, &associated.Number, &associated.IssuedOn); err != nil {
			return err
		}
		value.AssociatedVouchers = append(value.AssociatedVouchers, associated)
	}
	return rows.Err()
}
func fiscalDetailsValid(value fiscal.Invoice) bool {
	var base, vat, other int64
	seenVAT := map[int]bool{}
	for _, line := range value.VATLines {
		if line.ID < 1 || line.ID > 999 || line.BaseMinorUnits <= 0 || line.AmountMinorUnits < 0 || seenVAT[line.ID] {
			return false
		}
		seenVAT[line.ID] = true
		base += line.BaseMinorUnits
		vat += line.AmountMinorUnits
	}
	seenTax := map[int]bool{}
	for _, line := range value.OtherTaxLines {
		if line.ID < 1 || line.ID > 999 || len(line.Description) < 1 || len(line.Description) > 80 || line.BaseMinorUnits < 0 || line.RateBasisPoints < 0 || line.RateBasisPoints > 100000 || line.AmountMinorUnits < 0 || seenTax[line.ID] {
			return false
		}
		seenTax[line.ID] = true
		other += line.AmountMinorUnits
	}
	return vat == value.VATMinorUnits && other == value.OtherTaxMinorUnits && ((vat == 0 && len(value.VATLines) == 0) || (vat > 0 && base == value.NetMinorUnits)) && ((other == 0 && len(value.OtherTaxLines) == 0) || (other > 0 && len(value.OtherTaxLines) > 0))
}
````

### FILE: `internal/platform/postgres/fiscal_integration_test.go`
```yaml
block_id: "GO-ARCA-FISCAL:postgres-integration:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c4268a20e4a7436672c998a0e1573ff053848546345881b9dbf2e68c41676f88"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"elite.local/enterprise/internal/fiscal"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestFiscalPostgresSourceBindingLeaseReconciliationAndImmutability(t *testing.T) {
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
	var raw [16]byte
	if _, err = rand.Read(raw[:]); err != nil {
		t.Fatal(err)
	}
	raw[6] = (raw[6] & 0x0f) | 0x40
	raw[8] = (raw[8] & 0x3f) | 0x80
	tenant := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", raw[0:4], raw[4:6], raw[6:8], raw[8:10], raw[10:16])
	defer func() {
		tx, e := pool.Begin(ctx)
		if e != nil {
			t.Errorf("cleanup begin: %v", e)
			return
		}
		defer tx.Rollback(ctx)
		if _, e = tx.Exec(ctx, `set local session_replication_role=replica`); e != nil {
			t.Errorf("cleanup role: %v", e)
			return
		}
		for _, q := range []string{`delete from platform.outbox_event where tenant_id=$1`, `delete from fiscal.issuance_attempt where tenant_id=$1`, `delete from fiscal.invoice_associated_voucher where tenant_id=$1`, `delete from fiscal.invoice_vat where tenant_id=$1`, `delete from fiscal.invoice_other_tax where tenant_id=$1`, `delete from fiscal.invoice where tenant_id=$1`, `delete from fiscal.point_of_sale where tenant_id=$1`, `delete from payment.payment_attempt where tenant_id=$1`, `delete from sales.customer_order where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`} {
			if _, e = tx.Exec(ctx, q, tenant); e != nil {
				t.Errorf("cleanup: %v", e)
				return
			}
		}
		if e = tx.Commit(ctx); e != nil {
			t.Errorf("cleanup commit: %v", e)
		}
	}()
	fixtures := []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'fiscal-'||substring($1::text,1,8),'Fiscal','Fiscal')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'franchise','franchise','Franchise','franchisee'),($1,'other','other','Other','franchisee')`, `insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','franchise','customer','paid','ARS',12100,2),($1,'wrong-order','franchise','customer','paid','ARS',12000,1)`, `insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version)values($1,'payment','order','provider','reference','payment-fiscal','captured','ARS',12100,3),($1,'wrong-payment','wrong-order','provider','wrong-reference','wrong-payment-fiscal','captured','ARS',12000,1)`}
	for _, q := range fixtures {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFiscal(pool)
	pos := fiscal.PointOfSale{ID: "pos", OrganizationID: "franchise", TaxpayerCUIT: "30715117564", Environment: "homologation", Number: 1, Active: true, Version: 1}
	if err = repo.ConfigurePointOfSale(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c2a001", pos); err != nil {
		t.Fatal(err)
	}
	request := fiscal.Invoice{ID: "invoice", OrganizationID: "franchise", OrderID: "order", PaymentAttemptID: "payment", PointOfSaleID: "pos", VoucherType: 6, Concept: 1, RecipientDocumentType: 99, RecipientDocument: "0", RecipientVATConditionID: 5, NetMinorUnits: 10000, VATMinorUnits: 2100, VATLines: []fiscal.VATLine{{ID: 5, BaseMinorUnits: 10000, AmountMinorUnits: 2100}}, IssuedOn: time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC), Status: "queued", Version: 1}
	requestHash := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	created, replayed, err := repo.RequestInvoice(ctx, tenant, "1234567890abcdef", requestHash, "018f4d4a-7b36-7a21-8d10-2f4c54c2a002", request)
	if err != nil || replayed || created.TotalMinorUnits != 12100 {
		t.Fatalf("created=%+v replay=%v err=%v", created, replayed, err)
	}
	if len(created.VATLines) != 1 || created.VATLines[0].ID != 5 {
		t.Fatalf("VAT detail not retained: %+v", created.VATLines)
	}
	if created.RecipientVATConditionID != 5 {
		t.Fatalf("recipient VAT condition not retained: %+v", created)
	}
	if _, replayed, err = repo.RequestInvoice(ctx, tenant, "1234567890abcdef", requestHash, "unused", request); err != nil || !replayed {
		t.Fatalf("replay=%v err=%v", replayed, err)
	}
	request.ID = "wrong"
	request.PaymentAttemptID = "wrong-payment"
	if _, _, err = repo.RequestInvoice(ctx, tenant, "1234567890abcdeg", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "unused", request); !errors.Is(err, fiscal.ErrConflict) {
		t.Fatalf("source mismatch accepted: %v", err)
	}
	claimed, err := repo.ClaimInvoice(ctx, "worker", time.Minute, "attempt-claim")
	if err != nil || claimed.TenantID != tenant || claimed.ClaimedFrom != "queued" {
		t.Fatalf("claimed=%+v err=%v", claimed, err)
	}
	if _, err = repo.ClaimInvoice(ctx, "other-worker", time.Minute, "attempt-other"); !errors.Is(err, fiscal.ErrNoWork) {
		t.Fatalf("second lane claim accepted: %v", err)
	}
	claimed, err = repo.AssignVoucherNumber(ctx, claimed, "worker", 1, "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc", "attempt-sequence")
	if err != nil || claimed.VoucherNumber != 1 || claimed.Status != "authorizing" {
		t.Fatalf("assigned=%+v err=%v", claimed, err)
	}
	if err = repo.DeferInvoice(ctx, claimed, "worker", "authorize-timeout", "", true); err != nil {
		t.Fatal(err)
	}
	claimed, err = repo.ClaimInvoice(ctx, "worker", time.Minute, "attempt-reconcile")
	if err != nil || claimed.ClaimedFrom != "reconcile_required" || claimed.VoucherNumber != 1 {
		t.Fatalf("reclaimed=%+v err=%v", claimed, err)
	}
	expires := time.Date(2026, 9, 9, 0, 0, 0, 0, time.UTC)
	finished, err := repo.FinishInvoice(ctx, claimed, "worker", fiscal.Authorization{Found: true, Authorized: true, CAE: "12345678901234", CAEExpiresOn: &expires, ResponseHash: "dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd", ProviderCodes: nil}, "attempt-finish")
	if err != nil || finished.Status != "authorized" || finished.CAE != "12345678901234" {
		t.Fatalf("finished=%+v err=%v", finished, err)
	}
	creditRequest := fiscal.Invoice{ID: "credit", OrganizationID: "franchise", OrderID: "order", PaymentAttemptID: "payment", PointOfSaleID: "pos", VoucherType: 8, Concept: 1, RecipientDocumentType: 99, RecipientDocument: "0", RecipientVATConditionID: 5, NetMinorUnits: 10000, VATMinorUnits: 2100, VATLines: []fiscal.VATLine{{ID: 5, BaseMinorUnits: 10000, AmountMinorUnits: 2100}}, AssociatedVouchers: []fiscal.AssociatedVoucher{{InvoiceID: "invoice"}}, IssuedOn: time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC), Status: "queued", Version: 1}
	credit, replayed, err := repo.RequestInvoice(ctx, tenant, "1234567890credit", "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", "018f4d4a-7b36-7a21-8d10-2f4c54c2a003", creditRequest)
	if err != nil || replayed || len(credit.AssociatedVouchers) != 1 || credit.AssociatedVouchers[0].VoucherType != 6 || credit.AssociatedVouchers[0].PointOfSale != 1 || credit.AssociatedVouchers[0].Number != 1 || credit.AssociatedVouchers[0].TaxpayerCUIT != pos.TaxpayerCUIT {
		t.Fatalf("credit=%+v replay=%v err=%v", credit, replayed, err)
	}
	credit, err = repo.ClaimInvoice(ctx, "credit-worker", time.Minute, "attempt-credit-claim")
	if err != nil || credit.ID != "credit" || len(credit.AssociatedVouchers) != 1 || credit.AssociatedVouchers[0].InvoiceID != "invoice" {
		t.Fatalf("claimed credit=%+v err=%v", credit, err)
	}
	if _, err = pool.Exec(ctx, `update fiscal.issuance_attempt set outcome='deferred' where tenant_id=$1 and attempt_id='attempt-claim'`, tenant); err == nil {
		t.Fatal("immutable attempt update accepted")
	}
	var attempts, authorizedEvents int
	if err = pool.QueryRow(ctx, `select count(*) from fiscal.issuance_attempt where tenant_id=$1`, tenant).Scan(&attempts); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='fiscal-invoice.authorized'`, tenant).Scan(&authorizedEvents); err != nil {
		t.Fatal(err)
	}
	if attempts != 6 || authorizedEvents != 1 {
		t.Fatalf("attempts=%d authorized_events=%d", attempts, authorizedEvents)
	}
}
````

### FILE: `internal/platform/httpapi/fiscal.go`
```yaml
block_id: "GO-ARCA-FISCAL:http:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "41a6084789c802ad5bf3c000032656613b999b3d5e6ce446db4104b05161fc6e"
variables: []
secrets_allowed: false
```
````go
package httpapi

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"elite.local/enterprise/internal/fiscal"
	"elite.local/enterprise/internal/platform/identity"
)

type FiscalModule struct {
	Service    *fiscal.Service
	Parameters *fiscal.ParameterRegistry
}

func (m FiscalModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	api := fiscalAPI{service: m.Service, verifier: verifier}
	mux.HandleFunc("POST /v1/fiscal/points-of-sale", api.configurePointOfSale)
	mux.HandleFunc("POST /v1/fiscal/invoices", api.requestInvoice)
	mux.HandleFunc("GET /v1/fiscal/invoices/{id}", api.getInvoice)
	if m.Parameters != nil {
		parameterAPI{registry: m.Parameters, verifier: verifier}.Register(mux)
	}
}

type fiscalAPI struct {
	service  *fiscal.Service
	verifier identity.Verifier
}

func (a fiscalAPI) authorize(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return p, false
	}
	return p, true
}
func fiscalScope(w http.ResponseWriter, p identity.Principal, organization string) bool {
	if !p.AllowedOrganization(organization) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return false
	}
	return true
}
func (a fiscalAPI) configurePointOfSale(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "fiscal:configure")
	if !ok {
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return
	}
	var input fiscal.PointOfSale
	if !decodeStrict(w, r, &input) || !fiscalScope(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.ConfigurePointOfSale(r.Context(), p.TenantID, input)
	if err != nil {
		writeFiscalResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a fiscalAPI) requestInvoice(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "fiscal:issue")
	if !ok {
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 64<<10))
	if err != nil {
		writeProblem(w, 400, "INVALID_BODY", "body is invalid or too large")
		return
	}
	var input fiscal.Invoice
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&input) != nil {
		writeProblem(w, 400, "INVALID_BODY", "body does not match the contract")
		return
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		writeProblem(w, 400, "INVALID_BODY", "body must contain exactly one JSON value")
		return
	}
	if !fiscalScope(w, p, input.OrganizationID) {
		return
	}
	key := r.Header.Get("Idempotency-Key")
	sum := sha256.Sum256(body)
	value, replayed, err := a.service.RequestInvoice(r.Context(), p.TenantID, key, hex.EncodeToString(sum[:]), input)
	if err != nil {
		writeFiscalResult(w, err)
		return
	}
	if replayed {
		w.Header().Set("Idempotency-Replayed", "true")
		writeJSON(w, 200, value)
		return
	}
	writeJSON(w, 202, value)
}
func (a fiscalAPI) getInvoice(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "fiscal:read")
	if !ok {
		return
	}
	organization := r.URL.Query().Get("organization_id")
	if !fiscalScope(w, p, organization) {
		return
	}
	value, err := a.service.GetInvoice(r.Context(), p.TenantID, organization, r.PathValue("id"))
	if err != nil {
		writeFiscalResult(w, err)
		return
	}
	writeJSON(w, 200, value)
}
func writeFiscalResult(w http.ResponseWriter, err error) {
	if errors.Is(err, fiscal.ErrInvalid) {
		writeProblem(w, 400, "INVALID_FISCAL_COMMAND", "fiscal command does not match contract")
		return
	}
	if errors.Is(err, fiscal.ErrNotFound) {
		writeProblem(w, 404, "FISCAL_NOT_FOUND", "fiscal record was not found in this scope")
		return
	}
	if errors.Is(err, fiscal.ErrConflict) {
		writeProblem(w, 409, "FISCAL_CONFLICT", "fiscal state, source, sequence, idempotency or version conflict")
		return
	}
	writeProblem(w, 500, "INTERNAL", "request failed")
}
````

### FILE: `internal/platform/httpapi/fiscal_test.go`
```yaml
block_id: "GO-ARCA-FISCAL:http-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "5e5fe4e8dc27fad1225c4f4a7810129dc8165eba38f1c1410f4ab529e034deb9"
variables: []
secrets_allowed: false
```
````go
package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"elite.local/enterprise/internal/fiscal"
	"elite.local/enterprise/internal/platform/identity"
)

type fiscalRepoFake struct{ requests int }

func (*fiscalRepoFake) ConfigurePointOfSale(context.Context, string, string, fiscal.PointOfSale) error {
	return nil
}
func (r *fiscalRepoFake) RequestInvoice(_ context.Context, _, _, _, _ string, value fiscal.Invoice) (fiscal.Invoice, bool, error) {
	r.requests++
	value.Status = "queued"
	value.TotalMinorUnits = value.NetMinorUnits + value.VATMinorUnits
	return value, false, nil
}
func (*fiscalRepoFake) GetInvoice(context.Context, string, string, string) (fiscal.Invoice, error) {
	return fiscal.Invoice{ID: "invoice"}, nil
}

type fiscalVerifier struct{ p identity.Principal }

func (v fiscalVerifier) Verify(context.Context, string) (identity.Principal, error) { return v.p, nil }

type fiscalIDs struct{ n int }

func (i *fiscalIDs) New() string { i.n++; return "id" + string(rune('0'+i.n)) }
func TestFiscalHTTPRejectsCrossScopeAndQueuesAuthorizedScope(t *testing.T) {
	repo := &fiscalRepoFake{}
	service := fiscal.NewService(repo, &fiscalIDs{})
	mux := http.NewServeMux()
	principal := identity.Principal{Subject: "controller", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c29f12", Permissions: map[string]struct{}{"fiscal:issue": {}}, Organizations: map[string]struct{}{"franchise": {}}}
	FiscalModule{Service: service}.Register(mux, fiscalVerifier{p: principal})
	body := `{"organization_id":"other","order_id":"order","payment_attempt_id":"payment","point_of_sale_id":"pos","voucher_type":6,"concept":1,"recipient_document_type":99,"recipient_document":"0","recipient_vat_condition_id":5,"net_minor_units":10000,"vat_minor_units":2100,"vat_lines":[{"id":5,"base_minor_units":10000,"amount_minor_units":2100}],"issued_on":"2026-08-30T00:00:00Z"}`
	request := httptest.NewRequest("POST", "/v1/fiscal/invoices", strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "1234567890abcdef")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 || repo.requests != 0 {
		t.Fatalf("status=%d requests=%d", response.Code, repo.requests)
	}
	request = httptest.NewRequest("POST", "/v1/fiscal/invoices", strings.NewReader(strings.Replace(body, `"other"`, `"franchise"`, 1)))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "1234567890abcdef")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 202 || repo.requests != 1 {
		t.Fatalf("status=%d requests=%d body=%s", response.Code, repo.requests, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/fiscal/invoices", strings.NewReader(strings.Replace(body, `"other"`, `"franchise"`, 1)+` {}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "1234567890abcdef")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 400 || repo.requests != 1 {
		t.Fatalf("trailing JSON status=%d requests=%d", response.Code, repo.requests)
	}
}
````

### FILE: `db/migrations/0012_arca_fiscal_issuance.up.sql`
```yaml
block_id: "GO-ARCA-FISCAL:migration-up:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "aafbe0f1eb46b8f528a06002803c47869fb162918bf91c1dd8b083ddf71f4ecc"
variables: []
secrets_allowed: false
```
````sql
create schema if not exists fiscal;

create table fiscal.point_of_sale (
  tenant_id uuid not null,
  point_of_sale_id text not null,
  organization_id text not null,
  taxpayer_cuit text not null check (taxpayer_cuit ~ '^[0-9]{11}$'),
  environment text not null check (environment in ('homologation','production')),
  point_of_sale_number integer not null check (point_of_sale_number between 1 and 99999),
  active boolean not null default true,
  version bigint not null check (version > 0),
  active_invoice_id text,
  lease_owner text,
  lease_until timestamptz,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,point_of_sale_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  unique (tenant_id,taxpayer_cuit,environment,point_of_sale_number),
  check ((active_invoice_id is null and lease_owner is null and lease_until is null) or (active_invoice_id is not null and length(lease_owner) between 1 and 128 and lease_until is not null))
);

create table fiscal.invoice (
  tenant_id uuid not null,
  invoice_id text not null,
  organization_id text not null,
  order_id text not null,
  payment_attempt_id text not null,
  point_of_sale_id text not null,
  voucher_type integer not null check (voucher_type between 1 and 999),
  concept integer not null check (concept between 1 and 3),
  recipient_document_type integer not null check (recipient_document_type between 0 and 999),
  recipient_document text not null check (recipient_document ~ '^[0-9]{0,20}$'),
  recipient_vat_condition_id integer not null check (recipient_vat_condition_id between 1 and 999),
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  provider_currency text not null check (provider_currency ~ '^[A-Z]{3}$'),
  total_minor_units bigint not null check (total_minor_units > 0),
  net_minor_units bigint not null check (net_minor_units >= 0),
  vat_minor_units bigint not null check (vat_minor_units >= 0),
  exempt_minor_units bigint not null check (exempt_minor_units >= 0),
  non_taxed_minor_units bigint not null check (non_taxed_minor_units >= 0),
  other_tax_minor_units bigint not null check (other_tax_minor_units >= 0),
  issued_on date not null,
  service_from date,
  service_until date,
  payment_due_on date,
  status text not null check (status in ('queued','authorizing','reconcile_required','authorized','rejected')),
  voucher_number bigint check (voucher_number > 0),
  cae text check (cae ~ '^[0-9]{8,20}$'),
  cae_expires_on date,
  request_hash text not null check (request_hash ~ '^[a-f0-9]{64}$'),
  idempotency_key text not null check (length(idempotency_key) between 16 and 128),
  last_response_hash text check (last_response_hash ~ '^[a-f0-9]{64}$'),
  provider_codes jsonb not null default '[]'::jsonb check (jsonb_typeof(provider_codes)='array'),
  attempt_count integer not null default 0 check (attempt_count >= 0),
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  authorized_at timestamptz,
  primary key (tenant_id,invoice_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
  foreign key (tenant_id,payment_attempt_id) references payment.payment_attempt(tenant_id,payment_attempt_id),
  foreign key (tenant_id,point_of_sale_id) references fiscal.point_of_sale(tenant_id,point_of_sale_id),
  unique (tenant_id,idempotency_key),
  unique (tenant_id,payment_attempt_id,voucher_type),
  check (total_minor_units=net_minor_units+vat_minor_units+exempt_minor_units+non_taxed_minor_units+other_tax_minor_units),
  check ((concept=1 and service_from is null and service_until is null and payment_due_on is null) or (concept in (2,3) and service_from is not null and service_until >= service_from and payment_due_on >= service_until)),
  check ((status='authorized' and voucher_number is not null and cae is not null and cae_expires_on is not null and authorized_at is not null) or status<>'authorized'),
  check (status not in ('authorizing','reconcile_required') or voucher_number is not null)
);

create unique index fiscal_voucher_sequence_idx on fiscal.invoice(tenant_id,point_of_sale_id,voucher_type,voucher_number) where voucher_number is not null;
create index fiscal_invoice_work_idx on fiscal.invoice(status,created_at,tenant_id,invoice_id) where status in ('queued','reconcile_required');
create index fiscal_invoice_order_idx on fiscal.invoice(tenant_id,organization_id,order_id);

create table fiscal.invoice_vat (
  tenant_id uuid not null,
  invoice_id text not null,
  vat_id integer not null check (vat_id between 1 and 999),
  base_minor_units bigint not null check (base_minor_units > 0),
  amount_minor_units bigint not null check (amount_minor_units >= 0),
  primary key (tenant_id,invoice_id,vat_id),
  foreign key (tenant_id,invoice_id) references fiscal.invoice(tenant_id,invoice_id)
);

create table fiscal.invoice_other_tax (
  tenant_id uuid not null,
  invoice_id text not null,
  tax_id integer not null check (tax_id between 1 and 999),
  description text not null check (length(description) between 1 and 80),
  base_minor_units bigint not null check (base_minor_units >= 0),
  rate_basis_points integer not null check (rate_basis_points between 0 and 100000),
  amount_minor_units bigint not null check (amount_minor_units >= 0),
  primary key (tenant_id,invoice_id,tax_id),
  foreign key (tenant_id,invoice_id) references fiscal.invoice(tenant_id,invoice_id)
);

create or replace function fiscal.validate_invoice_details() returns trigger language plpgsql as $$
declare
  target_tenant uuid := case when tg_op='DELETE' then old.tenant_id else new.tenant_id end;
  target_invoice text := case when tg_op='DELETE' then old.invoice_id else new.invoice_id end;
  expected_net bigint; expected_vat bigint; expected_other bigint;
  actual_base bigint; actual_vat bigint; actual_other bigint; vat_rows bigint; tax_rows bigint;
begin
  select net_minor_units,vat_minor_units,other_tax_minor_units into expected_net,expected_vat,expected_other from fiscal.invoice where tenant_id=target_tenant and invoice_id=target_invoice;
  if not found then return null; end if;
  select coalesce(sum(base_minor_units),0),coalesce(sum(amount_minor_units),0),count(*) into actual_base,actual_vat,vat_rows from fiscal.invoice_vat where tenant_id=target_tenant and invoice_id=target_invoice;
  select coalesce(sum(amount_minor_units),0),count(*) into actual_other,tax_rows from fiscal.invoice_other_tax where tenant_id=target_tenant and invoice_id=target_invoice;
  if actual_vat<>expected_vat or actual_other<>expected_other or (expected_vat>0 and (vat_rows=0 or actual_base<>expected_net)) or (expected_vat=0 and vat_rows<>0) or (expected_other>0 and tax_rows=0) or (expected_other=0 and tax_rows<>0) then raise exception 'fiscal detail totals do not match invoice'; end if;
  return null;
end; $$;
create constraint trigger fiscal_invoice_details_valid after insert or update on fiscal.invoice deferrable initially deferred for each row execute function fiscal.validate_invoice_details();
create constraint trigger fiscal_invoice_vat_valid after insert or update or delete on fiscal.invoice_vat deferrable initially deferred for each row execute function fiscal.validate_invoice_details();
create constraint trigger fiscal_invoice_other_tax_valid after insert or update or delete on fiscal.invoice_other_tax deferrable initially deferred for each row execute function fiscal.validate_invoice_details();

create or replace function fiscal.prevent_fiscal_detail_mutation() returns trigger language plpgsql as $$ begin raise exception 'immutable fiscal detail'; end; $$;
create trigger fiscal_invoice_vat_immutable before update or delete on fiscal.invoice_vat for each row execute function fiscal.prevent_fiscal_detail_mutation();
create trigger fiscal_invoice_other_tax_immutable before update or delete on fiscal.invoice_other_tax for each row execute function fiscal.prevent_fiscal_detail_mutation();

create table fiscal.issuance_attempt (
  tenant_id uuid not null,
  attempt_id text not null,
  invoice_id text not null,
  phase text not null check (phase in ('claim','last-authorized','assign-number','consult','authorize','finish','defer')),
  voucher_number bigint,
  request_hash text not null check (request_hash ~ '^[a-f0-9]{64}$'),
  response_hash text check (response_hash ~ '^[a-f0-9]{64}$'),
  outcome text not null check (outcome in ('started','observed','authorized','rejected','deferred')),
  provider_codes jsonb not null default '[]'::jsonb check (jsonb_typeof(provider_codes)='array'),
  error_class text,
  recorded_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,attempt_id),
  foreign key (tenant_id,invoice_id) references fiscal.invoice(tenant_id,invoice_id)
);

create index fiscal_attempt_invoice_idx on fiscal.issuance_attempt(tenant_id,invoice_id,recorded_at,attempt_id);
create or replace function fiscal.prevent_attempt_mutation() returns trigger language plpgsql as $$ begin raise exception 'immutable fiscal attempt history'; end; $$;
create trigger fiscal_attempt_immutable before update or delete on fiscal.issuance_attempt for each row execute function fiscal.prevent_attempt_mutation();
````

### FILE: `db/migrations/0012_arca_fiscal_issuance.down.sql`
```yaml
block_id: "GO-ARCA-FISCAL:migration-down:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "5dc38808804f75fcca12f9a596abacd45280e1fa37c9ec446f9f9e493cdb6385"
variables: []
secrets_allowed: false
```
````sql
drop schema if exists fiscal cascade;
````

### FILE: `db/tests/0012_arca_fiscal_issuance.test.sql`
```yaml
block_id: "GO-ARCA-FISCAL:sql-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c75fe8b064437e3064ee5da682b868d863a4b09aeadbaa0e37ae4140c565465f"
variables: []
secrets_allowed: false
```
````sql
begin;
insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','fiscal-sql','Fiscal','Fiscal');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','franchise','franchise','Franchise','franchisee');
insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','order','franchise','customer','paid','ARS',12100,1);
insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','payment','order','provider','reference','fiscal-sql-payment','captured','ARS',12100,1);
insert into fiscal.point_of_sale(tenant_id,point_of_sale_id,organization_id,taxpayer_cuit,environment,point_of_sale_number,active,version) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','pos','franchise','30715117564','homologation',1,true,1);
insert into fiscal.invoice(tenant_id,invoice_id,organization_id,order_id,payment_attempt_id,point_of_sale_id,voucher_type,concept,recipient_document_type,recipient_document,recipient_vat_condition_id,currency,provider_currency,total_minor_units,net_minor_units,vat_minor_units,exempt_minor_units,non_taxed_minor_units,other_tax_minor_units,issued_on,status,request_hash,idempotency_key,version) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','invoice','franchise','order','payment','pos',6,1,99,'0',5,'ARS','PES',12100,10000,2100,0,0,0,'2026-08-30','queued',repeat('a',64),'1234567890abcdef',1);
insert into fiscal.invoice_vat(tenant_id,invoice_id,vat_id,base_minor_units,amount_minor_units) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','invoice',5,10000,2100);
do $$ begin
  begin insert into fiscal.invoice(tenant_id,invoice_id,organization_id,order_id,payment_attempt_id,point_of_sale_id,voucher_type,concept,recipient_document_type,recipient_document,recipient_vat_condition_id,currency,provider_currency,total_minor_units,net_minor_units,vat_minor_units,exempt_minor_units,non_taxed_minor_units,other_tax_minor_units,issued_on,status,request_hash,idempotency_key,version) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','bad-total','franchise','order','payment','pos',11,1,99,'0',5,'ARS','PES',12099,10000,2100,0,0,0,'2026-08-30','queued',repeat('b',64),'1234567890abcdeg',1); raise exception 'unbalanced fiscal total accepted'; exception when check_violation then null; end;
  insert into fiscal.issuance_attempt(tenant_id,attempt_id,invoice_id,phase,request_hash,outcome) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','attempt','invoice','claim',repeat('a',64),'started');
  begin update fiscal.invoice_vat set amount_minor_units=1 where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29f12'; raise exception 'VAT detail mutation accepted'; exception when raise_exception then if sqlerrm <> 'immutable fiscal detail' then raise; end if; end;
  begin update fiscal.issuance_attempt set outcome='deferred' where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29f12'; raise exception 'attempt mutation accepted'; exception when raise_exception then if sqlerrm <> 'immutable fiscal attempt history' then raise; end if; end;
end $$;
rollback;
````

### FILE: `internal/fiscal/parameters.go`
```yaml
block_id: "GO-ARCA-FISCAL:parameters:v1"
operation: CREATE
path: "internal/fiscal/parameters.go"
sha256: "184b09d1612e84e0b7fb2eab35ed2511c88796a1a7061176ea69115febe4d3f1"
provenance: AUTHORED
source: "local registry governed by ARCA WSFEv1 parameter operations"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````go
package fiscal

import (
	"context"
	"time"
)

var parameterKinds = map[string]bool{
	"voucher_type": true, "concept": true, "document_type": true, "vat_rate": true,
	"other_tax": true, "point_of_sale": true, "recipient_vat_condition": true,
}

type ParameterItem struct {
	Code           string  `json:"code"`
	Description    *string `json:"description"`
	ValidFrom      *string `json:"valid_from"`
	ValidUntil     *string `json:"valid_until"`
	VoucherClass   *string `json:"voucher_class"`
	EmissionType   *string `json:"emission_type"`
	Blocked        *string `json:"blocked"`
	DeregisteredOn *string `json:"deregistered_on"`
}

type ParameterSnapshot struct {
	TenantID       string          `json:"-"`
	ID             string          `json:"id"`
	OrganizationID string          `json:"organization_id"`
	TaxpayerCUIT   string          `json:"taxpayer_cuit"`
	Kind           string          `json:"kind"`
	VoucherClass   *string         `json:"voucher_class"`
	Items          []ParameterItem `json:"items"`
	ResponseHash   string          `json:"response_hash"`
	ProviderCodes  []string        `json:"provider_codes"`
	FetchedAt      time.Time       `json:"fetched_at"`
}

type ParameterSchedule struct {
	TenantID        string    `json:"-"`
	ID              string    `json:"id"`
	OrganizationID  string    `json:"organization_id"`
	TaxpayerCUIT    string    `json:"taxpayer_cuit"`
	Kind            string    `json:"kind"`
	VoucherClass    *string   `json:"voucher_class"`
	IntervalSeconds int64     `json:"interval_seconds"`
	NextRunAt       time.Time `json:"next_run_at"`
	Active          bool      `json:"active"`
	Version         int64     `json:"version"`
	LeaseOwner      string    `json:"-"`
}

type ParameterProvider interface {
	FetchParameters(context.Context, string, string, *string) (ParameterSnapshot, error)
}

type ParameterRepository interface {
	StoreParameterSnapshot(context.Context, ParameterSnapshot, string) (ParameterSnapshot, bool, error)
	GetParameterSnapshot(context.Context, string, string, string) (ParameterSnapshot, error)
	DecideParameterSnapshot(context.Context, string, string, string, string, bool, string, string) error
	ConfigureParameterSchedule(context.Context, ParameterSchedule, string, string) (ParameterSchedule, error)
}

type ParameterRegistry struct {
	repository ParameterRepository
	provider   ParameterProvider
	ids        IDGenerator
	now        func() time.Time
}

func NewParameterRegistry(repository ParameterRepository, provider ParameterProvider, ids IDGenerator, now func() time.Time) *ParameterRegistry {
	return &ParameterRegistry{repository: repository, provider: provider, ids: ids, now: now}
}

func (r *ParameterRegistry) Refresh(ctx context.Context, tenant, organization, taxpayerCUIT, kind string, voucherClass *string) (ParameterSnapshot, bool, error) {
	if tenant == "" || organization == "" || !ValidCUIT(taxpayerCUIT) || !validParameterScope(kind, voucherClass) {
		return ParameterSnapshot{}, false, ErrInvalid
	}
	value, err := r.provider.FetchParameters(ctx, taxpayerCUIT, kind, voucherClass)
	if err != nil {
		return ParameterSnapshot{}, false, err
	}
	if value.TaxpayerCUIT != taxpayerCUIT || value.Kind != kind || !sameOptional(value.VoucherClass, voucherClass) || !hashPattern.MatchString(value.ResponseHash) || len(value.Items) == 0 || !validParameterItems(value.Items, voucherClass) {
		return ParameterSnapshot{}, false, ErrInvalid
	}
	value.TenantID, value.OrganizationID, value.ID, value.FetchedAt = tenant, organization, r.ids.New(), r.now().UTC()
	return r.repository.StoreParameterSnapshot(ctx, value, r.ids.New())
}

func (r *ParameterRegistry) Snapshot(ctx context.Context, tenant, organization, snapshotID string) (ParameterSnapshot, error) {
	if tenant == "" || organization == "" || snapshotID == "" {
		return ParameterSnapshot{}, ErrInvalid
	}
	return r.repository.GetParameterSnapshot(ctx, tenant, organization, snapshotID)
}

func (r *ParameterRegistry) Decide(ctx context.Context, tenant, organization, snapshotID, subject string, approved bool, reason string) error {
	if tenant == "" || organization == "" || snapshotID == "" || subject == "" || len(reason) < 3 || len(reason) > 500 {
		return ErrInvalid
	}
	return r.repository.DecideParameterSnapshot(ctx, tenant, organization, snapshotID, subject, approved, reason, r.ids.New())
}

func (r *ParameterRegistry) ConfigureSchedule(ctx context.Context, tenant, subject string, value ParameterSchedule) (ParameterSchedule, error) {
	if tenant == "" || subject == "" || value.OrganizationID == "" || !ValidCUIT(value.TaxpayerCUIT) || !validParameterScope(value.Kind, value.VoucherClass) || value.IntervalSeconds < 900 || value.IntervalSeconds > 2592000 || value.NextRunAt.IsZero() {
		return value, ErrInvalid
	}
	value.TenantID, value.ID, value.Active, value.Version = tenant, r.ids.New(), true, 1
	value.NextRunAt = value.NextRunAt.UTC()
	return r.repository.ConfigureParameterSchedule(ctx, value, subject, r.ids.New())
}

func validParameterScope(kind string, voucherClass *string) bool {
	if !parameterKinds[kind] {
		return false
	}
	if kind != "recipient_vat_condition" {
		return voucherClass == nil
	}
	return voucherClass != nil && (*voucherClass == "A" || *voucherClass == "B" || *voucherClass == "C" || *voucherClass == "M")
}

func validParameterItems(items []ParameterItem, requestedClass *string) bool {
	seen := map[string]bool{}
	for _, item := range items {
		if item.Code == "" || len(item.Code) > 16 || seen[item.Code] || !validOptionalDate(item.ValidFrom) || !validOptionalDate(item.ValidUntil) || !validOptionalDate(item.DeregisteredOn) || (item.ValidFrom != nil && item.ValidUntil != nil && *item.ValidUntil < *item.ValidFrom) {
			return false
		}
		seen[item.Code] = true
		if requestedClass != nil && (item.VoucherClass == nil || *item.VoucherClass != *requestedClass) {
			return false
		}
	}
	return true
}

func validOptionalDate(value *string) bool {
	if value == nil {
		return true
	}
	_, err := time.Parse("2006-01-02", *value)
	return err == nil
}

func sameOptional(left, right *string) bool {
	return left == nil && right == nil || left != nil && right != nil && *left == *right
}
````

### FILE: `internal/fiscal/parameters_test.go`
```yaml
block_id: "GO-ARCA-FISCAL:parameter-tests:v1"
operation: CREATE
path: "internal/fiscal/parameters_test.go"
sha256: "268c8c86da47af16c2f9aef23c0e80db73aa2752cb3d2e61a7b3e28764dd1a5e"
provenance: AUTHORED
source: "local contracts governed by ARCA WSFEv1 parameter operations"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````go
package fiscal

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestParameterRegistryRefreshAndApprovalAreSeparated(t *testing.T) {
	description := "21%"
	provider := &parameterProvider{value: ParameterSnapshot{TaxpayerCUIT: "33693450239", Kind: "vat_rate", Items: []ParameterItem{{Code: "5", Description: &description}}, ResponseHash: string(make([]byte, 64))}}
	provider.value.ResponseHash = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	repository := &parameterRepository{}
	registry := NewParameterRegistry(repository, provider, sequenceIDs{}, func() time.Time { return time.Date(2026, 8, 30, 12, 0, 0, 0, time.UTC) })
	value, replay, err := registry.Refresh(context.Background(), "tenant", "organization", "33693450239", "vat_rate", nil)
	if err != nil || replay || value.ID == "" || repository.stored == nil || repository.decided {
		t.Fatalf("refresh=%+v replay=%v err=%v", value, replay, err)
	}
	if err = registry.Decide(context.Background(), "tenant", "organization", value.ID, "tax-owner", true, "validated against current ARCA configuration"); err != nil || !repository.decided {
		t.Fatalf("approval was not explicit: %v", err)
	}
}

func TestParameterRegistryRejectsInventedScopeAndDivergentResponse(t *testing.T) {
	registry := NewParameterRegistry(&parameterRepository{}, &parameterProvider{err: errors.New("must not call")}, sequenceIDs{}, time.Now)
	if _, _, err := registry.Refresh(context.Background(), "tenant", "organization", "33693450239", "invented", nil); !errors.Is(err, ErrInvalid) {
		t.Fatalf("invented kind accepted: %v", err)
	}
	class := "A"
	provider := &parameterProvider{value: ParameterSnapshot{TaxpayerCUIT: "33693450239", Kind: "recipient_vat_condition", VoucherClass: &class, Items: []ParameterItem{{Code: "1", VoucherClass: ptr("B")}}, ResponseHash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}}
	registry = NewParameterRegistry(&parameterRepository{}, provider, sequenceIDs{}, time.Now)
	if _, _, err := registry.Refresh(context.Background(), "tenant", "organization", "33693450239", "recipient_vat_condition", &class); !errors.Is(err, ErrInvalid) {
		t.Fatalf("divergent class accepted: %v", err)
	}
}

type parameterProvider struct {
	value ParameterSnapshot
	err   error
}

func (p *parameterProvider) FetchParameters(context.Context, string, string, *string) (ParameterSnapshot, error) {
	return p.value, p.err
}

type parameterRepository struct {
	stored  *ParameterSnapshot
	decided bool
}

func (r *parameterRepository) StoreParameterSnapshot(_ context.Context, value ParameterSnapshot, _ string) (ParameterSnapshot, bool, error) {
	r.stored = &value
	return value, false, nil
}
func (r *parameterRepository) GetParameterSnapshot(context.Context, string, string, string) (ParameterSnapshot, error) {
	return ParameterSnapshot{}, nil
}
func (r *parameterRepository) DecideParameterSnapshot(context.Context, string, string, string, string, bool, string, string) error {
	r.decided = true
	return nil
}
func (r *parameterRepository) ConfigureParameterSchedule(_ context.Context, value ParameterSchedule, _, _ string) (ParameterSchedule, error) {
	return value, nil
}

type sequenceIDs struct{}

func (sequenceIDs) New() string { return "018f4d4a-7b36-7a21-8d10-2f4c54c29f99" }
func ptr(value string) *string  { return &value }
````

### FILE: `internal/platform/postgres/fiscal_parameters.go`
```yaml
block_id: "GO-ARCA-FISCAL:parameter-repository:v1"
operation: CREATE
path: "internal/platform/postgres/fiscal_parameters.go"
sha256: "666127b53d5a112705d5b37970bace620d9a0884693ce9cd3f3d0d7691d267b1"
provenance: AUTHORED
source: "local durable repository governed by PostgreSQL 18 contracts"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"elite.local/enterprise/internal/fiscal"
	"github.com/jackc/pgx/v5"
)

func (r *Fiscal) GetParameterSnapshot(ctx context.Context, tenant, organization, snapshotID string) (fiscal.ParameterSnapshot, error) {
	var value fiscal.ParameterSnapshot
	var codes []byte
	err := r.pool.QueryRow(ctx, `select tenant_id::text,snapshot_id,organization_id,taxpayer_cuit,parameter_kind,voucher_class,response_hash,provider_codes,fetched_at from fiscal.parameter_snapshot where tenant_id=$1 and organization_id=$2 and snapshot_id=$3`, tenant, organization, snapshotID).Scan(&value.TenantID, &value.ID, &value.OrganizationID, &value.TaxpayerCUIT, &value.Kind, &value.VoucherClass, &value.ResponseHash, &codes, &value.FetchedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, fiscal.ErrNotFound
	}
	if err != nil {
		return value, err
	}
	if err = json.Unmarshal(codes, &value.ProviderCodes); err != nil {
		return value, err
	}
	rows, err := r.pool.Query(ctx, `select parameter_code,description,valid_from::text,valid_until::text,voucher_class,emission_type,blocked,deregistered_on::text from fiscal.parameter_item where tenant_id=$1 and snapshot_id=$2 order by parameter_code`, tenant, snapshotID)
	if err != nil {
		return value, err
	}
	defer rows.Close()
	for rows.Next() {
		var item fiscal.ParameterItem
		if err = rows.Scan(&item.Code, &item.Description, &item.ValidFrom, &item.ValidUntil, &item.VoucherClass, &item.EmissionType, &item.Blocked, &item.DeregisteredOn); err != nil {
			return value, err
		}
		value.Items = append(value.Items, item)
	}
	return value, rows.Err()
}

func (r *Fiscal) StoreParameterSnapshot(ctx context.Context, value fiscal.ParameterSnapshot, eventID string) (fiscal.ParameterSnapshot, bool, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return value, false, err
	}
	defer tx.Rollback(ctx)
	class := any(value.VoucherClass)
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1::text||'|'||$2::text||'|'||$3::text||'|'||coalesce($4::text,''),0))`, value.TenantID, value.OrganizationID, value.Kind, class); err != nil {
		return value, false, err
	}
	var existing string
	err = tx.QueryRow(ctx, `select snapshot_id from fiscal.parameter_snapshot where tenant_id=$1 and organization_id=$2 and taxpayer_cuit=$3 and parameter_kind=$4 and voucher_class is not distinct from $5 and response_hash=$6`, value.TenantID, value.OrganizationID, value.TaxpayerCUIT, value.Kind, class, value.ResponseHash).Scan(&existing)
	if err == nil {
		value.ID = existing
		return value, true, tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return value, false, err
	}
	codes, err := json.Marshal(append([]string{}, value.ProviderCodes...))
	if err != nil {
		return value, false, fiscal.ErrInvalid
	}
	result, err := tx.Exec(ctx, `insert into fiscal.parameter_snapshot(tenant_id,snapshot_id,organization_id,taxpayer_cuit,parameter_kind,voucher_class,response_hash,provider_codes,fetched_at) select $1,$2,$3,$4,$5,$6,$7,$8,$9 where exists(select 1 from org.organization where tenant_id=$1 and organization_id=$3 and status='active')`, value.TenantID, value.ID, value.OrganizationID, value.TaxpayerCUIT, value.Kind, class, value.ResponseHash, codes, value.FetchedAt)
	if err != nil {
		return value, false, fiscalConflict(err)
	}
	if result.RowsAffected() != 1 {
		return value, false, fiscal.ErrConflict
	}
	for _, item := range value.Items {
		_, err = tx.Exec(ctx, `insert into fiscal.parameter_item(tenant_id,snapshot_id,parameter_code,description,valid_from,valid_until,voucher_class,emission_type,blocked,deregistered_on) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, value.TenantID, value.ID, item.Code, item.Description, item.ValidFrom, item.ValidUntil, item.VoucherClass, item.EmissionType, item.Blocked, item.DeregisteredOn)
		if err != nil {
			return value, false, fiscalConflict(err)
		}
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'fiscal-parameter-snapshot',$3,1,'fiscal-parameter-snapshot.recorded',1,clock_timestamp(),jsonb_build_object('organization_id',$4::text,'taxpayer_cuit',$5::text,'parameter_kind',$6::text,'voucher_class',$7::text,'response_hash',$8::text,'item_count',$9::integer,'approved',false))`, value.TenantID, eventID, value.ID, value.OrganizationID, value.TaxpayerCUIT, value.Kind, class, value.ResponseHash, len(value.Items))
	if err != nil {
		return value, false, fiscalConflict(err)
	}
	return value, false, fiscalConflict(tx.Commit(ctx))
}

func (r *Fiscal) DecideParameterSnapshot(ctx context.Context, tenant, organization, snapshotID, subject string, approved bool, reason, eventID string) error {
	decision := "rejected"
	if approved {
		decision = "approved"
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into fiscal.parameter_decision(tenant_id,snapshot_id,decision,decided_by_subject,decision_reason) select s.tenant_id,s.snapshot_id,$4,$5,$6 from fiscal.parameter_snapshot s where s.tenant_id=$1 and s.organization_id=$2 and s.snapshot_id=$3`, tenant, organization, snapshotID, decision, subject, reason)
	if err != nil {
		return fiscalConflict(err)
	}
	if result.RowsAffected() != 1 {
		return fiscal.ErrNotFound
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'fiscal-parameter-snapshot',$3,2,'fiscal-parameter-snapshot.decided',1,clock_timestamp(),jsonb_build_object('organization_id',$4::text,'decision',$5::text,'decided_by_subject',$6::text,'decision_reason',$7::text))`, tenant, eventID, snapshotID, organization, decision, subject, reason)
	if err != nil {
		return fiscalConflict(err)
	}
	return fiscalConflict(tx.Commit(ctx))
}

func (r *Fiscal) ConfigureParameterSchedule(ctx context.Context, value fiscal.ParameterSchedule, subject, eventID string) (fiscal.ParameterSchedule, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	class := any(value.VoucherClass)
	err = tx.QueryRow(ctx, `insert into fiscal.parameter_refresh_schedule(tenant_id,schedule_id,organization_id,taxpayer_cuit,parameter_kind,voucher_class,interval_seconds,next_run_at,active,version,configured_by_subject) select $1,$2,$3,$4,$5,$6,$7,$8,true,1,$9 where exists(select 1 from org.organization where tenant_id=$1 and organization_id=$3 and status='active') on conflict (tenant_id,organization_id,taxpayer_cuit,parameter_kind,voucher_class) do update set interval_seconds=excluded.interval_seconds,next_run_at=excluded.next_run_at,active=true,version=fiscal.parameter_refresh_schedule.version+1,configured_by_subject=excluded.configured_by_subject,configured_at=clock_timestamp(),lease_owner=null,lease_until=null returning schedule_id,version`, value.TenantID, value.ID, value.OrganizationID, value.TaxpayerCUIT, value.Kind, class, value.IntervalSeconds, value.NextRunAt, subject).Scan(&value.ID, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, fiscal.ErrConflict
	}
	if err != nil {
		return value, fiscalConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'fiscal-parameter-schedule',$3,$4,'fiscal-parameter-schedule.configured',1,clock_timestamp(),jsonb_build_object('organization_id',$5::text,'taxpayer_cuit',$6::text,'parameter_kind',$7::text,'voucher_class',$8::text,'interval_seconds',$9::bigint,'configured_by_subject',$10::text))`, value.TenantID, eventID, value.ID, value.Version, value.OrganizationID, value.TaxpayerCUIT, value.Kind, class, value.IntervalSeconds, subject)
	if err != nil {
		return value, fiscalConflict(err)
	}
	return value, fiscalConflict(tx.Commit(ctx))
}

func (r *Fiscal) ClaimParameterSchedule(ctx context.Context, worker string, lease time.Duration) (fiscal.ParameterSchedule, error) {
	var value fiscal.ParameterSchedule
	err := r.pool.QueryRow(ctx, `with candidate as (select tenant_id,schedule_id from fiscal.parameter_refresh_schedule where active and next_run_at<=clock_timestamp() and (lease_until is null or lease_until<clock_timestamp()) order by next_run_at,tenant_id,schedule_id for update skip locked limit 1) update fiscal.parameter_refresh_schedule s set lease_owner=$1,lease_until=clock_timestamp()+$2::interval,attempts=attempts+1 from candidate c where s.tenant_id=c.tenant_id and s.schedule_id=c.schedule_id returning s.tenant_id::text,s.schedule_id,s.organization_id,s.taxpayer_cuit,s.parameter_kind,s.voucher_class,s.interval_seconds,s.next_run_at,s.active,s.version,s.lease_owner`, worker, lease.String()).Scan(&value.TenantID, &value.ID, &value.OrganizationID, &value.TaxpayerCUIT, &value.Kind, &value.VoucherClass, &value.IntervalSeconds, &value.NextRunAt, &value.Active, &value.Version, &value.LeaseOwner)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, fiscal.ErrNoParameterWork
	}
	return value, err
}

func (r *Fiscal) CompleteParameterSchedule(ctx context.Context, value fiscal.ParameterSchedule, worker, snapshotID, eventID string) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update fiscal.parameter_refresh_schedule set next_run_at=clock_timestamp()+make_interval(secs=>interval_seconds::double precision),lease_owner=null,lease_until=null,last_snapshot_id=$4,last_succeeded_at=clock_timestamp(),last_error_code=null,version=version+1 where tenant_id=$1 and schedule_id=$2 and lease_owner=$3 and lease_until>=clock_timestamp()`, value.TenantID, value.ID, worker, snapshotID)
	if err != nil {
		return fiscalConflict(err)
	}
	if result.RowsAffected() != 1 {
		return fiscal.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'fiscal-parameter-schedule',$3,$4+1,'fiscal-parameter-schedule.refreshed',1,clock_timestamp(),jsonb_build_object('organization_id',$5::text,'snapshot_id',$6::text))`, value.TenantID, eventID, value.ID, value.Version, value.OrganizationID, snapshotID)
	if err != nil {
		return fiscalConflict(err)
	}
	return fiscalConflict(tx.Commit(ctx))
}

func (r *Fiscal) DeferParameterSchedule(ctx context.Context, value fiscal.ParameterSchedule, worker, code string, retry time.Duration) error {
	result, err := r.pool.Exec(ctx, `update fiscal.parameter_refresh_schedule set next_run_at=clock_timestamp()+$4::interval,lease_owner=null,lease_until=null,last_error_code=$5,version=version+1 where tenant_id=$1 and schedule_id=$2 and lease_owner=$3`, value.TenantID, value.ID, worker, retry.String(), code)
	if err != nil {
		return fiscalConflict(err)
	}
	if result.RowsAffected() != 1 {
		return fiscal.ErrConflict
	}
	return nil
}
````

### FILE: `db/migrations/0013_arca_fiscal_parameters.up.sql`
```yaml
block_id: "GO-ARCA-FISCAL:migration-0013-up:v1"
operation: CREATE
path: "db/migrations/0013_arca_fiscal_parameters.up.sql"
sha256: "ee03c836d59f0c96f08ff98f3f49c5536b170c219837a02ec0d53385674eb506"
provenance: AUTHORED
source: "local schema governed by ARCA WSFEv1 and PostgreSQL 18"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````sql
create table fiscal.parameter_snapshot (
  tenant_id uuid not null,
  snapshot_id text not null,
  organization_id text not null,
  taxpayer_cuit text not null check (taxpayer_cuit ~ '^[0-9]{11}$'),
  parameter_kind text not null check (parameter_kind in ('voucher_type','concept','document_type','vat_rate','other_tax','point_of_sale','recipient_vat_condition')),
  voucher_class text,
  response_hash text not null check (response_hash ~ '^[a-f0-9]{64}$'),
  provider_codes jsonb not null check (jsonb_typeof(provider_codes)='array'),
  fetched_at timestamptz not null,
  recorded_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,snapshot_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  check ((parameter_kind='recipient_vat_condition' and voucher_class in ('A','B','C','M')) or (parameter_kind<>'recipient_vat_condition' and voucher_class is null)),
  unique nulls not distinct (tenant_id,organization_id,taxpayer_cuit,parameter_kind,voucher_class,response_hash)
);

create table fiscal.parameter_item (
  tenant_id uuid not null,
  snapshot_id text not null,
  parameter_code text not null check (length(parameter_code) between 1 and 16),
  description text,
  valid_from date,
  valid_until date,
  voucher_class text,
  emission_type text,
  blocked text,
  deregistered_on date,
  primary key (tenant_id,snapshot_id,parameter_code),
  foreign key (tenant_id,snapshot_id) references fiscal.parameter_snapshot(tenant_id,snapshot_id),
  check (valid_until is null or valid_from is null or valid_until>=valid_from)
);

create table fiscal.parameter_decision (
  tenant_id uuid not null,
  snapshot_id text not null,
  decision text not null check (decision in ('approved','rejected')),
  decided_by_subject text not null check (length(decided_by_subject) between 1 and 255),
  decided_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,snapshot_id),
  foreign key (tenant_id,snapshot_id) references fiscal.parameter_snapshot(tenant_id,snapshot_id)
);

create or replace function fiscal.prevent_parameter_history_mutation() returns trigger language plpgsql as $$ begin raise exception 'immutable fiscal parameter history'; end; $$;
create trigger fiscal_parameter_snapshot_immutable before update or delete on fiscal.parameter_snapshot for each row execute function fiscal.prevent_parameter_history_mutation();
create trigger fiscal_parameter_item_immutable before update or delete on fiscal.parameter_item for each row execute function fiscal.prevent_parameter_history_mutation();
create trigger fiscal_parameter_decision_immutable before update or delete on fiscal.parameter_decision for each row execute function fiscal.prevent_parameter_history_mutation();

create view fiscal.approved_parameter_snapshot as
select s.* from fiscal.parameter_snapshot s join fiscal.parameter_decision d using(tenant_id,snapshot_id) where d.decision='approved';
````

### FILE: `db/migrations/0013_arca_fiscal_parameters.down.sql`
```yaml
block_id: "GO-ARCA-FISCAL:migration-0013-down:v1"
operation: CREATE
path: "db/migrations/0013_arca_fiscal_parameters.down.sql"
sha256: "a11009f4d3a073bea148314352d3afd5957d9893f546e620c31e9218c47d361e"
provenance: AUTHORED
source: "local rollback governed by PostgreSQL 18"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````sql
drop view if exists fiscal.approved_parameter_snapshot;
drop table if exists fiscal.parameter_decision;
drop table if exists fiscal.parameter_item;
drop table if exists fiscal.parameter_snapshot;
drop function if exists fiscal.prevent_parameter_history_mutation();
````

### FILE: `db/tests/0013_arca_fiscal_parameters.test.sql`
```yaml
block_id: "GO-ARCA-FISCAL:test-0013:v1"
operation: CREATE
path: "db/tests/0013_arca_fiscal_parameters.test.sql"
sha256: "f9c188f97b0613d36250b1690e77dd3cdcb74fea4e2877eae4f9848cd0d9e733"
provenance: AUTHORED
source: "local PostgreSQL regression governed by the parameter approval contract"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````sql
begin;
insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f13','fiscal-param','Fiscal Param','Fiscal Param');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f13','franchise','franchise-param','Franchise','franchisee');
insert into fiscal.parameter_snapshot(tenant_id,snapshot_id,organization_id,taxpayer_cuit,parameter_kind,response_hash,provider_codes,fetched_at) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f13','snapshot','franchise','33693450239','vat_rate',repeat('a',64),'[]','2026-08-30T12:00:00Z');
insert into fiscal.parameter_item(tenant_id,snapshot_id,parameter_code,description,valid_from) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f13','snapshot','5','21%','2009-02-20');
do $$ begin
  if exists(select 1 from fiscal.approved_parameter_snapshot where snapshot_id='snapshot') then raise exception 'unapproved snapshot exposed'; end if;
  begin update fiscal.parameter_item set description='invented' where snapshot_id='snapshot'; raise exception 'parameter history mutation accepted'; exception when raise_exception then if sqlerrm <> 'immutable fiscal parameter history' then raise; end if; end;
end $$;
insert into fiscal.parameter_decision(tenant_id,snapshot_id,decision,decided_by_subject,decision_reason) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f13','snapshot','approved','tax-owner','validated against current ARCA configuration');
do $$ begin
  if (select count(*) from fiscal.approved_parameter_snapshot where snapshot_id='snapshot')<>1 then raise exception 'approved snapshot absent'; end if;
  begin insert into fiscal.parameter_decision(tenant_id,snapshot_id,decision,decided_by_subject,decision_reason) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f13','snapshot','rejected','other','separate rejected decision'); raise exception 'second decision accepted'; exception when unique_violation then null; end;
end $$;
rollback;
````

### FILE: `internal/fiscal/parameter_refresh.go`
```yaml
block_id: "GO-ARCA-FISCAL:parameter-refresh:v1"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by ARCA WSFEv1 4.6 and PostgreSQL 18 locking semantics"
license: "LicenseRef-Workspace-Owner"
sha256: "fcee3a2ab7bc5f89cdc31daa40676b7f6b7bca95813aebe30e3874ebfb3e21ce"
variables: []
secrets_allowed: false
```
````go
package fiscal

import (
	"context"
	"errors"
	"time"
)

var ErrNoParameterWork = errors.New("no parameter refresh work")

type ParameterScheduleRepository interface {
	ClaimParameterSchedule(context.Context, string, time.Duration) (ParameterSchedule, error)
	CompleteParameterSchedule(context.Context, ParameterSchedule, string, string, string) error
	DeferParameterSchedule(context.Context, ParameterSchedule, string, string, time.Duration) error
}

type ParameterRefreshProcessor struct {
	repository ParameterScheduleRepository
	registry   *ParameterRegistry
	workerID   string
	lease      time.Duration
	retry      time.Duration
	ids        IDGenerator
}

func NewParameterRefreshProcessor(repository ParameterScheduleRepository, registry *ParameterRegistry, ids IDGenerator, workerID string, lease, retry time.Duration) (*ParameterRefreshProcessor, error) {
	if repository == nil || registry == nil || ids == nil || workerID == "" || lease < 30*time.Second || lease > 10*time.Minute || retry < time.Minute || retry > time.Hour {
		return nil, ErrInvalid
	}
	return &ParameterRefreshProcessor{repository: repository, registry: registry, workerID: workerID, lease: lease, retry: retry, ids: ids}, nil
}

func (p *ParameterRefreshProcessor) ProcessOne(ctx context.Context) error {
	schedule, err := p.repository.ClaimParameterSchedule(ctx, p.workerID, p.lease)
	if err != nil {
		return err
	}
	snapshot, _, err := p.registry.Refresh(ctx, schedule.TenantID, schedule.OrganizationID, schedule.TaxpayerCUIT, schedule.Kind, schedule.VoucherClass)
	if err != nil {
		_ = p.repository.DeferParameterSchedule(ctx, schedule, p.workerID, "PARAMETER_PROVIDER_FAILURE", p.retry)
		return err
	}
	return p.repository.CompleteParameterSchedule(ctx, schedule, p.workerID, snapshot.ID, p.ids.New())
}
````

### FILE: `internal/fiscal/parameter_refresh_test.go`
```yaml
block_id: "GO-ARCA-FISCAL:parameter-refresh-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local regression tests"
license: "LicenseRef-Workspace-Owner"
sha256: "5891af45539435d8b8b45c45e8f6bbc838721d0230df31e1f46d8d9380643ce3"
variables: []
secrets_allowed: false
```
````go
package fiscal

import (
	"context"
	"errors"
	"testing"
	"time"
)

type scheduleRepository struct {
	schedule  ParameterSchedule
	completed bool
	deferred  bool
}

func (r *scheduleRepository) ClaimParameterSchedule(context.Context, string, time.Duration) (ParameterSchedule, error) {
	return r.schedule, nil
}
func (r *scheduleRepository) CompleteParameterSchedule(_ context.Context, _ ParameterSchedule, _, snapshotID, _ string) error {
	r.completed = snapshotID != ""
	return nil
}
func (r *scheduleRepository) DeferParameterSchedule(context.Context, ParameterSchedule, string, string, time.Duration) error {
	r.deferred = true
	return nil
}

func TestParameterRefreshProcessorStoresButNeverApproves(t *testing.T) {
	repository := &parameterRepository{}
	provider := &parameterProvider{value: ParameterSnapshot{TaxpayerCUIT: "33693450239", Kind: "vat_rate", Items: []ParameterItem{{Code: "5"}}, ResponseHash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}}
	registry := NewParameterRegistry(repository, provider, sequenceIDs{}, time.Now)
	schedules := &scheduleRepository{schedule: ParameterSchedule{TenantID: "tenant", OrganizationID: "organization", TaxpayerCUIT: "33693450239", Kind: "vat_rate"}}
	processor, err := NewParameterRefreshProcessor(schedules, registry, sequenceIDs{}, "worker", time.Minute, 5*time.Minute)
	if err != nil || processor.ProcessOne(context.Background()) != nil || !schedules.completed || schedules.deferred || repository.decided {
		t.Fatalf("processor err=%v completed=%v deferred=%v approved=%v", err, schedules.completed, schedules.deferred, repository.decided)
	}
	provider.err = errors.New("provider unavailable")
	provider.value = ParameterSnapshot{}
	schedules.completed, schedules.deferred = false, false
	if err = processor.ProcessOne(context.Background()); err == nil || !schedules.deferred || schedules.completed {
		t.Fatalf("failure err=%v completed=%v deferred=%v", err, schedules.completed, schedules.deferred)
	}
}
````

### FILE: `internal/platform/postgres/fiscal_parameter_schedule_integration_test.go`
```yaml
block_id: "GO-ARCA-FISCAL:parameter-schedule-pg-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL integration regression"
license: "LicenseRef-Workspace-Owner"
sha256: "711b488f5bfa12aa0aba3803186b7ca5ec40c5f222b9e88542a41de9fa70251f"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"elite.local/enterprise/internal/fiscal"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestFiscalParameterScheduleLeaseSnapshotAndExplicitDecision(t *testing.T) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c2f133"
	cleanup := []string{
		`delete from platform.outbox_event where tenant_id=$1`, `delete from fiscal.parameter_refresh_schedule where tenant_id=$1`, `delete from fiscal.parameter_decision where tenant_id=$1`, `delete from fiscal.parameter_item where tenant_id=$1`, `delete from fiscal.parameter_snapshot where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`,
	}
	defer func() {
		tx, e := pool.Begin(ctx)
		if e != nil {
			t.Errorf("cleanup begin: %v", e)
			return
		}
		defer tx.Rollback(ctx)
		if _, e = tx.Exec(ctx, `set local session_replication_role=replica`); e != nil {
			t.Errorf("cleanup role: %v", e)
			return
		}
		for _, query := range cleanup {
			if _, e = tx.Exec(ctx, query, tenant); e != nil {
				t.Errorf("cleanup: %v", e)
				return
			}
		}
		if e = tx.Commit(ctx); e != nil {
			t.Errorf("cleanup commit: %v", e)
		}
	}()
	for _, query := range []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'fiscal-v133-integration','Fiscal V133','Fiscal V133')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'franchise','franchise-v133','Franchise','franchisee')`} {
		if _, err = pool.Exec(ctx, query, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repository := NewFiscal(pool)
	now := time.Now().UTC()
	snapshot := fiscal.ParameterSnapshot{TenantID: tenant, ID: "snapshot-v133", OrganizationID: "franchise", TaxpayerCUIT: "33693450239", Kind: "vat_rate", Items: []fiscal.ParameterItem{{Code: "5"}}, ResponseHash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", FetchedAt: now}
	if _, _, err = repository.StoreParameterSnapshot(ctx, snapshot, "018f4d4a-7b36-7a21-8d10-2f4c54c2f134"); err != nil {
		t.Fatal(err)
	}
	schedule := fiscal.ParameterSchedule{TenantID: tenant, ID: "schedule-v133", OrganizationID: "franchise", TaxpayerCUIT: "33693450239", Kind: "vat_rate", IntervalSeconds: 900, NextRunAt: now.Add(-time.Minute), Active: true, Version: 1}
	schedule, err = repository.ConfigureParameterSchedule(ctx, schedule, "tax-owner", "018f4d4a-7b36-7a21-8d10-2f4c54c2f135")
	if err != nil {
		t.Fatal(err)
	}
	claimed, err := repository.ClaimParameterSchedule(ctx, "worker-a", time.Minute)
	if err != nil || claimed.ID != schedule.ID {
		t.Fatalf("claim=%+v err=%v", claimed, err)
	}
	if _, err = repository.ClaimParameterSchedule(ctx, "worker-b", time.Minute); !errors.Is(err, fiscal.ErrNoParameterWork) {
		t.Fatalf("second claim=%v", err)
	}
	if err = repository.CompleteParameterSchedule(ctx, claimed, "worker-a", snapshot.ID, "018f4d4a-7b36-7a21-8d10-2f4c54c2f136"); err != nil {
		t.Fatal(err)
	}
	if err = repository.DecideParameterSnapshot(ctx, tenant, "franchise", snapshot.ID, "tax-owner", true, "validated against current ARCA configuration", "018f4d4a-7b36-7a21-8d10-2f4c54c2f137"); err != nil {
		t.Fatal(err)
	}
	loaded, err := repository.GetParameterSnapshot(ctx, tenant, "franchise", snapshot.ID)
	if err != nil || len(loaded.Items) != 1 || loaded.Items[0].Code != "5" {
		t.Fatalf("loaded=%+v err=%v", loaded, err)
	}
}
````

### FILE: `internal/platform/httpapi/fiscal_parameters.go`
```yaml
block_id: "GO-ARCA-FISCAL:parameter-admin-http:v1"
operation: CREATE
provenance: AUTHORED
source: "local OIDC and resource-scope integration"
license: "LicenseRef-Workspace-Owner"
sha256: "e474b1cbd96307333c4ee0a657e75d17adba3b87890f6004382b8a12cfe4676b"
variables: []
secrets_allowed: false
```
````go
package httpapi

import (
	"encoding/json"
	"io"
	"net/http"

	"elite.local/enterprise/internal/fiscal"
	"elite.local/enterprise/internal/platform/identity"
)

type parameterAPI struct {
	registry *fiscal.ParameterRegistry
	verifier identity.Verifier
}

func (a parameterAPI) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/fiscal/parameter-schedules", a.configureSchedule)
	mux.HandleFunc("POST /v1/fiscal/parameters/refresh", a.refresh)
	mux.HandleFunc("GET /v1/fiscal/parameters/{id}", a.getSnapshot)
	mux.HandleFunc("POST /v1/fiscal/parameters/{id}/decisions", a.decide)
}

func (a parameterAPI) principal(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return p, false
	}
	return p, true
}

func parameterJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return false
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
	decoder.DisallowUnknownFields()
	if decoder.Decode(target) != nil || decoder.Decode(&struct{}{}) != io.EOF {
		writeProblem(w, 400, "INVALID_BODY", "body must contain exactly one value matching the contract")
		return false
	}
	return true
}

func (a parameterAPI) configureSchedule(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "fiscal:parameters:configure")
	if !ok {
		return
	}
	var input fiscal.ParameterSchedule
	if !parameterJSON(w, r, &input) || !fiscalScope(w, p, input.OrganizationID) {
		return
	}
	value, err := a.registry.ConfigureSchedule(r.Context(), p.TenantID, p.Subject, input)
	if err != nil {
		writeFiscalResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}

func (a parameterAPI) refresh(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "fiscal:parameters:refresh")
	if !ok {
		return
	}
	var input struct {
		OrganizationID string  `json:"organization_id"`
		TaxpayerCUIT   string  `json:"taxpayer_cuit"`
		Kind           string  `json:"kind"`
		VoucherClass   *string `json:"voucher_class"`
	}
	if !parameterJSON(w, r, &input) || !fiscalScope(w, p, input.OrganizationID) {
		return
	}
	value, replayed, err := a.registry.Refresh(r.Context(), p.TenantID, input.OrganizationID, input.TaxpayerCUIT, input.Kind, input.VoucherClass)
	if err != nil {
		writeFiscalResult(w, err)
		return
	}
	if replayed {
		w.Header().Set("Idempotency-Replayed", "true")
	}
	writeJSON(w, 201, value)
}

func (a parameterAPI) getSnapshot(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "fiscal:parameters:read")
	if !ok {
		return
	}
	organization := r.URL.Query().Get("organization_id")
	if !fiscalScope(w, p, organization) {
		return
	}
	value, err := a.registry.Snapshot(r.Context(), p.TenantID, organization, r.PathValue("id"))
	if err != nil {
		writeFiscalResult(w, err)
		return
	}
	writeJSON(w, 200, value)
}

func (a parameterAPI) decide(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "fiscal:parameters:approve")
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		Approved       bool   `json:"approved"`
		Reason         string `json:"reason"`
	}
	if !parameterJSON(w, r, &input) || !fiscalScope(w, p, input.OrganizationID) {
		return
	}
	if err := a.registry.Decide(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), p.Subject, input.Approved, input.Reason); err != nil {
		writeFiscalResult(w, err)
		return
	}
	w.WriteHeader(204)
}
````

### FILE: `internal/platform/httpapi/fiscal_parameters_test.go`
```yaml
block_id: "GO-ARCA-FISCAL:parameter-admin-http-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local authorization regression tests"
license: "LicenseRef-Workspace-Owner"
sha256: "cf85a989ed029c62277ed3988540e520cb0bb3958c007daf0ad121baa85a26ea"
variables: []
secrets_allowed: false
```
````go
package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/fiscal"
	"elite.local/enterprise/internal/platform/identity"
)

type parameterHTTPRepository struct {
	decisions int
	schedules int
	snapshots int
}

func (r *parameterHTTPRepository) StoreParameterSnapshot(_ context.Context, value fiscal.ParameterSnapshot, _ string) (fiscal.ParameterSnapshot, bool, error) {
	r.snapshots++
	return value, false, nil
}
func (r *parameterHTTPRepository) GetParameterSnapshot(context.Context, string, string, string) (fiscal.ParameterSnapshot, error) {
	return fiscal.ParameterSnapshot{ID: "snapshot"}, nil
}
func (r *parameterHTTPRepository) DecideParameterSnapshot(context.Context, string, string, string, string, bool, string, string) error {
	r.decisions++
	return nil
}
func (r *parameterHTTPRepository) ConfigureParameterSchedule(_ context.Context, value fiscal.ParameterSchedule, _, _ string) (fiscal.ParameterSchedule, error) {
	r.schedules++
	return value, nil
}

type parameterHTTPProvider struct{}

func (parameterHTTPProvider) FetchParameters(context.Context, string, string, *string) (fiscal.ParameterSnapshot, error) {
	return fiscal.ParameterSnapshot{TaxpayerCUIT: "33693450239", Kind: "vat_rate", Items: []fiscal.ParameterItem{{Code: "5"}}, ResponseHash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}, nil
}

func TestFiscalParameterAdminSeparatesPermissionsAndScope(t *testing.T) {
	repository := &parameterHTTPRepository{}
	registry := fiscal.NewParameterRegistry(repository, parameterHTTPProvider{}, &fiscalIDs{}, time.Now)
	p := identity.Principal{Subject: "tax-owner", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c29f12", Organizations: map[string]struct{}{"franchise": {}}, Permissions: map[string]struct{}{"fiscal:parameters:refresh": {}, "fiscal:parameters:approve": {}, "fiscal:parameters:configure": {}}}
	mux := http.NewServeMux()
	FiscalModule{Service: fiscal.NewService(&fiscalRepoFake{}, &fiscalIDs{}), Parameters: registry}.Register(mux, fiscalVerifier{p: p})
	cases := []struct {
		path, body string
		want       int
	}{
		{"/v1/fiscal/parameters/refresh", `{"organization_id":"other","taxpayer_cuit":"33693450239","kind":"vat_rate","voucher_class":null}`, 403},
		{"/v1/fiscal/parameters/refresh", `{"organization_id":"franchise","taxpayer_cuit":"33693450239","kind":"vat_rate","voucher_class":null}`, 201},
		{"/v1/fiscal/parameters/snapshot/decisions", `{"organization_id":"franchise","approved":true,"reason":"validated against ARCA configuration"}`, 204},
		{"/v1/fiscal/parameter-schedules", `{"organization_id":"franchise","taxpayer_cuit":"33693450239","kind":"vat_rate","voucher_class":null,"interval_seconds":900,"next_run_at":"2026-08-31T00:00:00Z"}`, 201},
	}
	for _, tc := range cases {
		request := httptest.NewRequest("POST", tc.path, strings.NewReader(tc.body))
		request.Header.Set("Authorization", "Bearer valid")
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		if response.Code != tc.want {
			t.Fatalf("%s status=%d body=%s", tc.path, response.Code, response.Body.String())
		}
	}
	if repository.snapshots != 1 || repository.decisions != 1 || repository.schedules != 1 {
		t.Fatalf("snapshots=%d decisions=%d schedules=%d", repository.snapshots, repository.decisions, repository.schedules)
	}
}
````

### FILE: `db/migrations/0014_arca_parameter_refresh_admin.up.sql`
```yaml
block_id: "GO-ARCA-FISCAL:migration-0014-up:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL schema governed by PostgreSQL 18 documentation"
license: "LicenseRef-Workspace-Owner"
sha256: "2e7ba695ccce73b259e9ac2a9120a40d2abc830a425fa42575ccd484472341bb"
variables: []
secrets_allowed: false
```
````sql
alter table fiscal.parameter_decision add column decision_reason text check (decision_reason is null or length(decision_reason) between 3 and 500);
create or replace function fiscal.require_parameter_decision_reason() returns trigger language plpgsql as $$ begin if new.decision_reason is null then raise exception 'decision reason required'; end if; return new; end; $$;
create trigger fiscal_parameter_decision_reason before insert on fiscal.parameter_decision for each row execute function fiscal.require_parameter_decision_reason();

create table fiscal.parameter_refresh_schedule (
  tenant_id uuid not null,
  schedule_id text not null,
  organization_id text not null,
  taxpayer_cuit text not null check (taxpayer_cuit ~ '^[0-9]{11}$'),
  parameter_kind text not null check (parameter_kind in ('voucher_type','concept','document_type','vat_rate','other_tax','point_of_sale','recipient_vat_condition')),
  voucher_class text,
  interval_seconds bigint not null check (interval_seconds between 900 and 2592000),
  next_run_at timestamptz not null,
  active boolean not null,
  attempts bigint not null default 0,
  lease_owner text,
  lease_until timestamptz,
  last_snapshot_id text,
  last_succeeded_at timestamptz,
  last_error_code text,
  version bigint not null check (version>0),
  configured_by_subject text not null check (length(configured_by_subject) between 1 and 255),
  configured_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,schedule_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,last_snapshot_id) references fiscal.parameter_snapshot(tenant_id,snapshot_id),
  check ((parameter_kind='recipient_vat_condition' and voucher_class in ('A','B','C','M')) or (parameter_kind<>'recipient_vat_condition' and voucher_class is null)),
  check ((lease_owner is null)=(lease_until is null)),
  unique nulls not distinct (tenant_id,organization_id,taxpayer_cuit,parameter_kind,voucher_class)
);

create index fiscal_parameter_refresh_due_idx on fiscal.parameter_refresh_schedule(next_run_at) where active;
````

### FILE: `db/migrations/0014_arca_parameter_refresh_admin.down.sql`
```yaml
block_id: "GO-ARCA-FISCAL:migration-0014-down:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL rollback"
license: "LicenseRef-Workspace-Owner"
sha256: "457a4cf937c732d4ffa19184ebdede8591d35c20853963bfd89adba2241d78dd"
variables: []
secrets_allowed: false
```
````sql
drop table if exists fiscal.parameter_refresh_schedule;
drop trigger if exists fiscal_parameter_decision_reason on fiscal.parameter_decision;
drop function if exists fiscal.require_parameter_decision_reason();
alter table fiscal.parameter_decision drop column if exists decision_reason;
````

### FILE: `db/tests/0014_arca_parameter_refresh_admin.test.sql`
```yaml
block_id: "GO-ARCA-FISCAL:migration-0014-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL contract regression"
license: "LicenseRef-Workspace-Owner"
sha256: "1797dc8134350d562699bcf1096af70a029fafbb42646ce2286fa1c8ad09adf2"
variables: []
secrets_allowed: false
```
````sql
begin;
insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name,status) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','fiscal-v133','Fiscal V133 SA','Fiscal V133','active');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type,status) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','org-v133','org-v133','Organization V133','franchisee','active');
insert into fiscal.parameter_refresh_schedule(tenant_id,schedule_id,organization_id,taxpayer_cuit,parameter_kind,interval_seconds,next_run_at,active,version,configured_by_subject) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','schedule-v133','org-v133','33693450239','vat_rate',900,clock_timestamp(),true,1,'controller');
do $$ begin
  if (select count(*) from fiscal.parameter_refresh_schedule where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29f12' and active and next_run_at<=clock_timestamp())<>1 then raise exception 'due schedule missing'; end if;
  begin
    insert into fiscal.parameter_refresh_schedule(tenant_id,schedule_id,organization_id,taxpayer_cuit,parameter_kind,interval_seconds,next_run_at,active,version,configured_by_subject) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','bad-v133','org-v133','33693450239','vat_rate',899,clock_timestamp(),true,1,'controller');
    raise exception 'unsafe cadence accepted';
  exception when check_violation then null; end;
end $$;
rollback;
````

### FILE: `db/migrations/0023_arca_associated_voucher.up.sql`
```yaml
block_id: "GO-ARCA-FISCAL:migration-0023-up:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL enforcement governed by the official ARCA WSFEv1 associated-voucher contract"
license: "LicenseRef-Workspace-Owner"
sha256: "0d12d958a561c5d14262a38116d25d09c6f3bf8b47550eee1ddab046d7290014"
variables: []
secrets_allowed: false
```
````sql
create table fiscal.invoice_associated_voucher (
  tenant_id uuid not null,
  invoice_id text not null,
  associated_invoice_id text not null,
  taxpayer_cuit text not null check (taxpayer_cuit ~ '^[0-9]{11}$'),
  voucher_type integer not null check (voucher_type in (1,6,11)),
  point_of_sale_number integer not null check (point_of_sale_number between 1 and 99999),
  voucher_number bigint not null check (voucher_number > 0),
  issued_on date not null,
  primary key (tenant_id,invoice_id),
  foreign key (tenant_id,invoice_id) references fiscal.invoice(tenant_id,invoice_id),
  foreign key (tenant_id,associated_invoice_id) references fiscal.invoice(tenant_id,invoice_id),
  unique (tenant_id,associated_invoice_id),
  check (invoice_id<>associated_invoice_id)
);

create or replace function fiscal.validate_associated_voucher() returns trigger language plpgsql as $$
declare
  target_tenant uuid := case when tg_table_name='invoice' then new.tenant_id else case when tg_op='DELETE' then old.tenant_id else new.tenant_id end end;
  target_invoice text := case when tg_table_name='invoice' then new.invoice_id else case when tg_op='DELETE' then old.invoice_id else new.invoice_id end end;
  child_type integer;
  association_count bigint;
  valid_count bigint;
begin
  select voucher_type into child_type from fiscal.invoice where tenant_id=target_tenant and invoice_id=target_invoice;
  if not found then return null; end if;
  select count(*) into association_count from fiscal.invoice_associated_voucher where tenant_id=target_tenant and invoice_id=target_invoice;
  if child_type in (3,8,13) then
    if association_count<>1 then raise exception 'credit note requires exactly one associated voucher'; end if;
    select count(*) into valid_count
    from fiscal.invoice_associated_voucher a
    join fiscal.invoice child on child.tenant_id=a.tenant_id and child.invoice_id=a.invoice_id
    join fiscal.invoice original on original.tenant_id=a.tenant_id and original.invoice_id=a.associated_invoice_id
    join fiscal.point_of_sale original_pos on original_pos.tenant_id=original.tenant_id and original_pos.point_of_sale_id=original.point_of_sale_id
    join fiscal.point_of_sale child_pos on child_pos.tenant_id=child.tenant_id and child_pos.point_of_sale_id=child.point_of_sale_id
    where a.tenant_id=target_tenant and a.invoice_id=target_invoice
      and original.status='authorized' and original.organization_id=child.organization_id
      and original.voucher_type=case child.voucher_type when 3 then 1 when 8 then 6 when 13 then 11 end
      and original.voucher_number=a.voucher_number and original.issued_on=a.issued_on
      and original_pos.taxpayer_cuit=a.taxpayer_cuit and original_pos.taxpayer_cuit=child_pos.taxpayer_cuit
      and original_pos.point_of_sale_number=a.point_of_sale_number and a.issued_on<=child.issued_on;
    if valid_count<>1 then raise exception 'associated voucher identity is invalid'; end if;
  elsif child_type in (1,6,11) then
    if association_count<>0 then raise exception 'invoice cannot contain an associated voucher'; end if;
  else
    raise exception 'voucher type outside admitted invoice and credit-note lane';
  end if;
  return null;
end; $$;

create constraint trigger fiscal_invoice_association_valid after insert or update on fiscal.invoice deferrable initially deferred for each row execute function fiscal.validate_associated_voucher();
create constraint trigger fiscal_association_valid after insert or update or delete on fiscal.invoice_associated_voucher deferrable initially deferred for each row execute function fiscal.validate_associated_voucher();
create trigger fiscal_association_immutable before update or delete on fiscal.invoice_associated_voucher for each row execute function fiscal.prevent_fiscal_detail_mutation();
````

### FILE: `db/migrations/0023_arca_associated_voucher.down.sql`
```yaml
block_id: "GO-ARCA-FISCAL:migration-0023-down:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL rollback"
license: "LicenseRef-Workspace-Owner"
sha256: "85e54f2f9952dfd3aa6e2db2ff34f6ad9a4d6389963878948adf9e46b9a5a0a1"
variables: []
secrets_allowed: false
```
````sql
drop trigger if exists fiscal_association_immutable on fiscal.invoice_associated_voucher;
drop trigger if exists fiscal_association_valid on fiscal.invoice_associated_voucher;
drop trigger if exists fiscal_invoice_association_valid on fiscal.invoice;
drop function if exists fiscal.validate_associated_voucher();
drop table if exists fiscal.invoice_associated_voucher;
````

### FILE: `db/tests/0023_arca_associated_voucher.test.sql`
```yaml
block_id: "GO-ARCA-FISCAL:migration-0023-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL contract regression"
license: "LicenseRef-Workspace-Owner"
sha256: "3fc2b474686a11b9cb4a310c5857cb4c78b8ce7727c464e1a367e4c0523f8697"
variables: []
secrets_allowed: false
```
````sql
begin;
insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values('018f4d4a-7b36-7a21-8d10-2f4c54c2a023','fiscal-assoc','Fiscal','Fiscal');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values('018f4d4a-7b36-7a21-8d10-2f4c54c2a023','franchise','franchise','Franchise','franchisee');
insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values('018f4d4a-7b36-7a21-8d10-2f4c54c2a023','order','franchise','customer','paid','ARS',12100,1);
insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version) values('018f4d4a-7b36-7a21-8d10-2f4c54c2a023','payment','order','provider','reference','fiscal-assoc-payment','captured','ARS',12100,1);
insert into fiscal.point_of_sale(tenant_id,point_of_sale_id,organization_id,taxpayer_cuit,environment,point_of_sale_number,active,version) values('018f4d4a-7b36-7a21-8d10-2f4c54c2a023','pos','franchise','30715117564','homologation',1,true,1);
insert into fiscal.invoice(tenant_id,invoice_id,organization_id,order_id,payment_attempt_id,point_of_sale_id,voucher_type,concept,recipient_document_type,recipient_document,recipient_vat_condition_id,currency,provider_currency,total_minor_units,net_minor_units,vat_minor_units,exempt_minor_units,non_taxed_minor_units,other_tax_minor_units,issued_on,status,voucher_number,cae,cae_expires_on,authorized_at,request_hash,idempotency_key,version) values('018f4d4a-7b36-7a21-8d10-2f4c54c2a023','original','franchise','order','payment','pos',6,1,99,'0',5,'ARS','PES',12100,10000,2100,0,0,0,'2026-08-30','authorized',10,'12345678901234','2026-09-09',clock_timestamp(),repeat('a',64),'1234567890abcdef',1);
insert into fiscal.invoice_vat(tenant_id,invoice_id,vat_id,base_minor_units,amount_minor_units) values('018f4d4a-7b36-7a21-8d10-2f4c54c2a023','original',5,10000,2100);
insert into fiscal.invoice(tenant_id,invoice_id,organization_id,order_id,payment_attempt_id,point_of_sale_id,voucher_type,concept,recipient_document_type,recipient_document,recipient_vat_condition_id,currency,provider_currency,total_minor_units,net_minor_units,vat_minor_units,exempt_minor_units,non_taxed_minor_units,other_tax_minor_units,issued_on,status,request_hash,idempotency_key,version) values('018f4d4a-7b36-7a21-8d10-2f4c54c2a023','credit','franchise','order','payment','pos',8,1,99,'0',5,'ARS','PES',12100,10000,2100,0,0,0,'2026-08-31','queued',repeat('b',64),'1234567890abcdeg',1);
insert into fiscal.invoice_vat(tenant_id,invoice_id,vat_id,base_minor_units,amount_minor_units) values('018f4d4a-7b36-7a21-8d10-2f4c54c2a023','credit',5,10000,2100);
insert into fiscal.invoice_associated_voucher(tenant_id,invoice_id,associated_invoice_id,taxpayer_cuit,voucher_type,point_of_sale_number,voucher_number,issued_on) values('018f4d4a-7b36-7a21-8d10-2f4c54c2a023','credit','original','30715117564',6,1,10,'2026-08-30');
set constraints all immediate;
do $$ begin
  begin update fiscal.invoice_associated_voucher set voucher_number=11 where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2a023'; raise exception 'associated identity mutation accepted'; exception when raise_exception then if sqlerrm <> 'immutable fiscal detail' then raise; end if; end;
end $$;
rollback;
````

## 6. Configuration surface

| Variable | Type/default | Secret | Validation/effect |
|---|---|---|---|
| project point of sale | CUIT checksum + environment + integer | no | must be configured per organization and enabled in ARCA |
| worker identities | fiscal + parameter IDs / none | no | each owns a bounded 30s–10m lease |
| parameter refresh cadence | 900–2592000 seconds / none | no | configured explicitly per organization/CUIT/kind; no statutory cadence is invented |
| WSAA certificate/private key | external secret reference / none | yes | belongs to the separate credential pack; never stored here |
| ARCA environment | `homologation` or `production` | no | production requires separate approval and evidence |

## 7. Dependency bill

| Dependency | Pin | Use | License | Scope | Source |
|---|---|---|---|---|---|
| Go | `1.26.7` | domain, worker orchestration and HTTP | BSD-3-Clause | build/runtime | `go.dev` |
| PostgreSQL | `18.6` | durable workflow, lane, audit and outbox | PostgreSQL | runtime/test | `postgresql.org` |
| pgx | `5.10.0` | transactions and SQL | MIT | build/runtime | `github.com/jackc/pgx` |
| ARCA WSFEv1 | manual `4.6` current at verification | fiscal protocol authority | public specification | design/runtime contract | `arca.gob.ar` |

## 8. Apply order

Materialize after migrations 0001–0011 and their Go owners. Apply 0012–0014 y luego 0023 en el perfil integrado, execute every SQL and Go test, wire `FiscalModule` and both workers in the application, then provide the UDS transport backed by generated WSFE plus WSAA. Acquire or refresh a snapshot before a separately authorized subject approves or rejects it with a reason; no scheduler response becomes policy automatically. Start only in homologation. Rollback drains both workers and removes 0023 antes de 0014 y las migraciones fiscales anteriores only when retained history permits it.

## 9. Verification

Require 28/28 exact materialization/hash, gofmt, full Go tests, vet and builds; PostgreSQL 18.6 migrations 0001–0014 plus 0023 and all SQL invariants. Prove payment/total/scope rejection, replay, one active lane, exact sequence, timeout→reconciliation, immutable attempts, CAE, exact associated-voucher identity for `3→1`, `8→6`, `13→11`, seven parameter scopes, response hash, duplicate rejection, immutable history, reasoned decision, exclusive schedule lease, bounded explicit cadence, retry and exclusion of every unapproved snapshot from policy. Live promotion additionally requires official ARCA homologation credentials, exact generated SOAP contracts, approved tax/currency rules, CAE/QR/PDF output, error-code policy, load, recovery and accountant/business acceptance.

## 10. Reconstruction evidence

V126–V128 preserve issuance; V132 adds the durable registry. V133 adds authenticated administration and durable scheduled refresh while preserving the separate approval boundary. V145 adds the associated-voucher lane required for full credit notes A/B/C. Current hashes, commands, results and remaining live conditions are recorded in `reconstruction_evidence/ARCA_ASSOCIATED_VOUCHER_EXECUTION_INVENTORY_2026-08-31_V145.md`; earlier evidence remains historical.

V402 composed delta: V402325 connected local ARCA fixture, fixed offline generated-tool adaptation, explicit legal/source pins and portable .NET/Go build. See ARCA_CONNECTED_INFRA_V402 evidence; only future homologation credentials conditioned, no productive fiscal claim.

### FILE: `internal/platform/httpapi/fiscal_connected_reference_test.go`

```yaml
block_id: "GO_ARCA_FISCAL_ISSUANCE_API-V402325:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "42fc0e164001a0681b344683f010978dd2cc49ed265fa1701b3dda69916c0b56"
variables: []
secrets_allowed: false
```

````go
package httpapi

// AUTHORED local composition fixture. Paid-source rows and provider responses are
// synthetic; actual HTTP handlers, PostgreSQL, Go worker, UDS and .NET owners run.
import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/fiscal"
	"elite.local/enterprise/internal/fiscal/wsfeipc"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type fiscalReferenceReceipt struct {
	Tenant  string         `json:"tenant"`
	Invoice string         `json:"invoice"`
	Credit  string         `json:"credit"`
	Stats   map[string]any `json:"stats"`
}

func TestFiscalConnectedReference(t *testing.T) {
	dsn, socket, path := os.Getenv("TEST_DATABASE_URL"), os.Getenv("ARCA_FIXTURE_SOCKET"), os.Getenv("ARCA_FIXTURE_RECEIPT")
	if dsn == "" || socket == "" || path == "" {
		t.Skip("explicit local fixture environment required")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	ids := randomid.Generator{}
	tenant := ids.New()
	for _, sql := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'arca-'||$1::text,'Local ARCA Fixture','Local ARCA Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'franchise','franchise','Fixture','franchisee')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'paid-order','franchise','fixture-customer','paid','ARS',12100,2)`,
		`insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version)values($1,'captured-payment','paid-order','fixture','fixture-provider-reference','fixture-payment-key','captured','ARS',12100,3)`,
	} {
		if _, err = pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := postgres.NewFiscal(pool)
	service := fiscal.NewService(repo, ids)
	principal := identity.Principal{Subject: "fixture-controller", TenantID: tenant, Permissions: map[string]struct{}{"fiscal:configure": {}, "fiscal:issue": {}, "fiscal:read": {}}, Organizations: map[string]struct{}{"franchise": {}}}
	mux := http.NewServeMux()
	FiscalModule{Service: service}.Register(mux, fiscalVerifier{p: principal})
	server := httptest.NewServer(mux)
	defer server.Close()
	call := func(method, url, key string, body []byte, want int) []byte {
		t.Helper()
		request, er := http.NewRequestWithContext(ctx, method, server.URL+url, bytes.NewReader(body))
		if er != nil {
			t.Fatal(er)
		}
		request.Header.Set("Authorization", "Bearer fixture-principal")
		request.Header.Set("Content-Type", "application/json")
		if key != "" {
			request.Header.Set("Idempotency-Key", key)
		}
		response, er := server.Client().Do(request)
		if er != nil {
			t.Fatal(er)
		}
		defer response.Body.Close()
		raw, er := io.ReadAll(io.LimitReader(response.Body, 65537))
		if er != nil || response.StatusCode != want {
			t.Fatalf("fiscal HTTP status=%d expected=%d error=%v body=%s", response.StatusCode, want, er, raw)
		}
		return raw
	}
	posRaw := call("POST", "/v1/fiscal/points-of-sale", "", []byte(`{"organization_id":"franchise","taxpayer_cuit":"30715117564","environment":"homologation","number":12}`), 201)
	var pos fiscal.PointOfSale
	if err = json.Unmarshal(posRaw, &pos); err != nil || pos.ID == "" {
		t.Fatalf("point-of-sale %v", err)
	}
	invoice := fiscal.Invoice{OrganizationID: "franchise", OrderID: "paid-order", PaymentAttemptID: "captured-payment", PointOfSaleID: pos.ID, VoucherType: 6, Concept: 1, RecipientDocumentType: 99, RecipientDocument: "0", RecipientVATConditionID: 5, NetMinorUnits: 10000, VATMinorUnits: 2100, VATLines: []fiscal.VATLine{{ID: 5, BaseMinorUnits: 10000, AmountMinorUnits: 2100}}, IssuedOn: time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC)}
	body, _ := json.Marshal(invoice)
	raw := call("POST", "/v1/fiscal/invoices", "fixture-issue-once", body, 202)
	var queued fiscal.Invoice
	if err = json.Unmarshal(raw, &queued); err != nil {
		t.Fatal(err)
	}
	replay := call("POST", "/v1/fiscal/invoices", "fixture-issue-once", body, 200)
	var same fiscal.Invoice
	if err = json.Unmarshal(replay, &same); err != nil || same.ID != queued.ID {
		t.Fatal("invoice replay changed identity")
	}
	call("POST", "/v1/fiscal/invoices", "fixture-denied-org", bytes.Replace(body, []byte(`"franchise"`), []byte(`"other"`), 1), 403)
	provider, err := wsfeipc.NewUnixProvider(socket, 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	processor, err := fiscal.NewProcessor(repo, provider, ids, "fixture-worker-a", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = processor.ProcessOne(ctx); err == nil {
		t.Fatal("lost provider response unexpectedly accepted")
	}
	pending, err := repo.GetInvoice(ctx, tenant, "franchise", queued.ID)
	if err != nil || pending.Status != "reconcile_required" || pending.VoucherNumber != 1 {
		t.Fatalf("ambiguous request was not durable: %+v %v", pending, err)
	}
	resumed, err := fiscal.NewProcessor(postgres.NewFiscal(pool), provider, ids, "fixture-worker-restarted", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	authorized, err := resumed.ProcessOne(ctx)
	if err != nil || authorized.ID != queued.ID || authorized.Status != "authorized" || authorized.CAE != "41124578989845" {
		t.Fatalf("consult/reconcile failed: %+v %v", authorized, err)
	}
	invoice.VoucherType = 8
	invoice.AssociatedVouchers = []fiscal.AssociatedVoucher{{InvoiceID: queued.ID}}
	invoice.IssuedOn = invoice.IssuedOn.AddDate(0, 0, 1)
	creditBody, _ := json.Marshal(invoice)
	raw = call("POST", "/v1/fiscal/invoices", "fixture-credit-once", creditBody, 202)
	var credit fiscal.Invoice
	if err = json.Unmarshal(raw, &credit); err != nil {
		t.Fatal(err)
	}
	credit, err = resumed.ProcessOne(ctx)
	if err != nil || credit.Status != "authorized" || len(credit.AssociatedVouchers) != 1 || credit.AssociatedVouchers[0].Number != 1 || credit.AssociatedVouchers[0].VoucherType != 6 {
		t.Fatalf("credit association did not cross actual bridge: %+v %v", credit, err)
	}
	if _, err = resumed.ProcessOne(ctx); !errors.Is(err, fiscal.ErrNoWork) {
		t.Fatalf("duplicate work left: %v", err)
	}
	for _, id := range []string{queued.ID, credit.ID} {
		raw = call("GET", "/v1/fiscal/invoices/"+id+"?organization_id=franchise", "", nil, 200)
		if !bytes.Contains(raw, []byte(`"status":"authorized"`)) {
			t.Fatal("authorized invoice not readable")
		}
	}
	for _, kind := range []string{"voucher_type", "concept", "document_type", "vat_rate", "other_tax", "point_of_sale", "recipient_vat_condition"} {
		var class *string
		if kind == "recipient_vat_condition" {
			v := "A"
			class = &v
		}
		snapshot, er := provider.FetchParameters(ctx, "30715117564", kind, class)
		if er != nil || len(snapshot.Items) != 1 || len(snapshot.ResponseHash) != 64 {
			t.Fatalf("parameter %s: %v", kind, er)
		}
	}
	dialer := net.Dialer{Timeout: time.Second}
	client := http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{DialContext: func(c context.Context, _, _ string) (net.Conn, error) { return dialer.DialContext(c, "unix", socket) }}}
	response, err := client.Get("http://arca-wsfe/fixture/stats")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var stats map[string]any
	if err = json.NewDecoder(response.Body).Decode(&stats); err != nil {
		t.Fatal(err)
	}
	if stats["authorize"] != float64(2) || stats["last"] != float64(2) || stats["consult"] != float64(1) || stats["cms"] != float64(1) {
		t.Fatalf("unexpected provider side effects: %+v", stats)
	}
	var rows, events int
	if err = pool.QueryRow(ctx, `select count(*) from fiscal.invoice where tenant_id=$1 and status='authorized'`, tenant).Scan(&rows); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='fiscal-invoice.authorized'`, tenant).Scan(&events); err != nil || rows != 2 || events != 2 {
		t.Fatalf("durable effects rows=%d events=%d %v", rows, events, err)
	}
	receipt, _ := json.Marshal(fiscalReferenceReceipt{Tenant: tenant, Invoice: queued.ID, Credit: credit.ID, Stats: stats})
	if err = os.WriteFile(path, receipt, 0600); err != nil {
		t.Fatal(err)
	}
	t.Log("ARCA_CONNECTED_HTTP_PG_GO_UDS_DOTNET_CMS_PASS invoices=2 authorized_events=2 authorize_calls=2 lost_response_reconciled=1 credential_refreshes=1 parameter_kinds=7")
	stopped, stopErr := client.Post("http://arca-wsfe/fixture/stop", "application/json", strings.NewReader("{}"))
	if stopErr != nil {
		t.Fatal(stopErr)
	}
	defer stopped.Body.Close()
	if stopped.StatusCode != 200 {
		t.Fatal("fixture graceful stop failed")
	}
}

func TestFiscalReferenceRecovery(t *testing.T) {
	path, dsn := os.Getenv("ARCA_FIXTURE_RECEIPT"), os.Getenv("TEST_DATABASE_URL")
	if path == "" || dsn == "" {
		t.Skip("explicit recovery fixture required")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var r fiscalReferenceReceipt
	if err = json.Unmarshal(raw, &r); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	repo := postgres.NewFiscal(pool)
	for _, id := range []string{r.Invoice, r.Credit} {
		v, er := repo.GetInvoice(ctx, r.Tenant, "franchise", id)
		if er != nil || v.Status != "authorized" || v.CAE != "41124578989845" || !strings.HasPrefix(v.Environment, "homo") {
			t.Fatalf("durable fiscal recovery %+v %v", v, er)
		}
	}
	t.Log("ARCA_POSTGRES_RESTART_PASS invoices=2 provider_calls=0")
}
````

