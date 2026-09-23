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
