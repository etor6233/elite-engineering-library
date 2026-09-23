package wsfeipc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"time"

	"elite.local/enterprise/internal/fiscal"
)

var ErrUnavailable = errors.New("ARCA WSFE bridge unavailable")
var ErrProtocol = errors.New("ARCA WSFE bridge protocol violation")

type Provider struct {
	client *http.Client
	base   string
}

func NewUnixProvider(socketPath string, timeout time.Duration) (*Provider, error) {
	if socketPath == "" || !filepath.IsAbs(socketPath) || timeout < 5*time.Second || timeout > time.Minute {
		return nil, fiscal.ErrInvalid
	}
	dialer := &net.Dialer{Timeout: 5 * time.Second, KeepAlive: 30 * time.Second}
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, "unix", socketPath)
		},
		DisableCompression:  true,
		MaxIdleConns:        2,
		MaxIdleConnsPerHost: 2,
		IdleConnTimeout:     30 * time.Second,
	}
	client := &http.Client{Transport: transport, Timeout: timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	return &Provider{client: client, base: "http://arca-wsfe"}, nil
}

type wireInvoice struct {
	TaxpayerCUIT            string                  `json:"taxpayer_cuit"`
	PointOfSale             int                     `json:"point_of_sale"`
	VoucherType             int                     `json:"voucher_type"`
	Concept                 int                     `json:"concept"`
	RecipientDocumentType   int                     `json:"recipient_document_type"`
	RecipientDocument       string                  `json:"recipient_document"`
	RecipientVATConditionID int                     `json:"recipient_vat_condition_id"`
	Currency                string                  `json:"currency"`
	VoucherNumber           int64                   `json:"voucher_number"`
	TotalMinorUnits         int64                   `json:"total_minor_units"`
	NetMinorUnits           int64                   `json:"net_minor_units"`
	VATMinorUnits           int64                   `json:"vat_minor_units"`
	ExemptMinorUnits        int64                   `json:"exempt_minor_units"`
	NonTaxedMinorUnits      int64                   `json:"non_taxed_minor_units"`
	OtherTaxMinorUnits      int64                   `json:"other_tax_minor_units"`
	IssuedOn                string                  `json:"issued_on"`
	ServiceFrom             *string                 `json:"service_from,omitempty"`
	ServiceUntil            *string                 `json:"service_until,omitempty"`
	PaymentDueOn            *string                 `json:"payment_due_on,omitempty"`
	VATLines                []fiscal.VATLine        `json:"vat_lines"`
	OtherTaxLines           []fiscal.OtherTaxLine   `json:"other_tax_lines"`
	AssociatedVouchers      []wireAssociatedVoucher `json:"associated_vouchers"`
}

type wireAssociatedVoucher struct {
	TaxpayerCUIT string `json:"taxpayer_cuit"`
	VoucherType  int    `json:"voucher_type"`
	PointOfSale  int    `json:"point_of_sale"`
	Number       int64  `json:"number"`
	IssuedOn     string `json:"issued_on"`
}

type safeResult struct {
	Found         bool     `json:"found"`
	Authorized    bool     `json:"authorized"`
	CAE           *string  `json:"cae"`
	CAEExpiresOn  *string  `json:"cae_expires_on"`
	ResponseHash  string   `json:"response_hash"`
	ProviderCodes []string `json:"provider_codes"`
}

type lastResult struct {
	Number        int64    `json:"number"`
	ResponseHash  string   `json:"response_hash"`
	ProviderCodes []string `json:"provider_codes"`
}

type parameterRequest struct {
	TaxpayerCUIT string  `json:"taxpayer_cuit"`
	Kind         string  `json:"kind"`
	VoucherClass *string `json:"voucher_class"`
}

type parameterResult struct {
	Kind          string                 `json:"kind"`
	VoucherClass  *string                `json:"voucher_class"`
	Items         []fiscal.ParameterItem `json:"items"`
	ResponseHash  string                 `json:"response_hash"`
	ProviderCodes []string               `json:"provider_codes"`
}

func (p *Provider) FetchParameters(ctx context.Context, taxpayerCUIT, kind string, voucherClass *string) (fiscal.ParameterSnapshot, error) {
	var result parameterResult
	if err := p.postAny(ctx, "/v1/parameters", parameterRequest{TaxpayerCUIT: taxpayerCUIT, Kind: kind, VoucherClass: voucherClass}, &result); err != nil {
		return fiscal.ParameterSnapshot{}, err
	}
	if result.Kind != kind || !validHash(result.ResponseHash) || len(result.Items) == 0 {
		return fiscal.ParameterSnapshot{}, ErrProtocol
	}
	return fiscal.ParameterSnapshot{TaxpayerCUIT: taxpayerCUIT, Kind: result.Kind, VoucherClass: result.VoucherClass, Items: result.Items, ResponseHash: result.ResponseHash, ProviderCodes: append([]string{}, result.ProviderCodes...)}, nil
}

func (p *Provider) LastAuthorized(ctx context.Context, invoice fiscal.Invoice) (int64, string, error) {
	var result lastResult
	if err := p.post(ctx, "/v1/last-authorized", request(invoice), &result); err != nil {
		return 0, "", err
	}
	if result.Number < 0 || !validHash(result.ResponseHash) {
		return 0, "", ErrProtocol
	}
	return result.Number, result.ResponseHash, nil
}

func (p *Provider) Consult(ctx context.Context, invoice fiscal.Invoice) (fiscal.Authorization, error) {
	return p.authorization(ctx, "/v1/consult", invoice)
}

func (p *Provider) Authorize(ctx context.Context, invoice fiscal.Invoice) (fiscal.Authorization, error) {
	return p.authorization(ctx, "/v1/authorize", invoice)
}

func (p *Provider) authorization(ctx context.Context, path string, invoice fiscal.Invoice) (fiscal.Authorization, error) {
	var result safeResult
	if err := p.post(ctx, path, request(invoice), &result); err != nil {
		return fiscal.Authorization{}, err
	}
	value := fiscal.Authorization{Found: result.Found, Authorized: result.Authorized, ResponseHash: result.ResponseHash, ProviderCodes: append([]string{}, result.ProviderCodes...)}
	if result.CAE != nil {
		value.CAE = *result.CAE
	}
	if result.CAEExpiresOn != nil {
		parsed, err := time.Parse("2006-01-02", *result.CAEExpiresOn)
		if err != nil {
			return fiscal.Authorization{}, ErrProtocol
		}
		value.CAEExpiresOn = &parsed
	}
	if !validHash(value.ResponseHash) {
		return fiscal.Authorization{}, ErrProtocol
	}
	return value, nil
}

func (p *Provider) post(ctx context.Context, path string, input wireInvoice, output any) error {
	return p.postAny(ctx, path, input, output)
}

func (p *Provider) postAny(ctx context.Context, path string, input any, output any) error {
	body, err := json.Marshal(input)
	if err != nil {
		return ErrProtocol
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.base+path, bytes.NewReader(body))
	if err != nil {
		return ErrProtocol
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	response, err := p.client.Do(req)
	if err != nil {
		return ErrUnavailable
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 64<<10))
		return fmt.Errorf("%w: status %d", ErrUnavailable, response.StatusCode)
	}
	payload, err := io.ReadAll(io.LimitReader(response.Body, (64<<10)+1))
	if err != nil || len(payload) > 64<<10 {
		return ErrProtocol
	}
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(output); err != nil {
		return ErrProtocol
	}
	if err = decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return ErrProtocol
	}
	return nil
}

func request(value fiscal.Invoice) wireInvoice {
	associated := make([]wireAssociatedVoucher, len(value.AssociatedVouchers))
	for index, voucher := range value.AssociatedVouchers {
		associated[index] = wireAssociatedVoucher{TaxpayerCUIT: voucher.TaxpayerCUIT, VoucherType: voucher.VoucherType, PointOfSale: voucher.PointOfSale, Number: voucher.Number, IssuedOn: date(voucher.IssuedOn)}
	}
	return wireInvoice{
		TaxpayerCUIT: value.TaxpayerCUIT, PointOfSale: value.PointOfSaleNumber, VoucherType: value.VoucherType, Concept: value.Concept,
		RecipientDocumentType: value.RecipientDocumentType, RecipientDocument: value.RecipientDocument, RecipientVATConditionID: value.RecipientVATConditionID,
		Currency: value.Currency, VoucherNumber: value.VoucherNumber, TotalMinorUnits: value.TotalMinorUnits, NetMinorUnits: value.NetMinorUnits,
		VATMinorUnits: value.VATMinorUnits, ExemptMinorUnits: value.ExemptMinorUnits, NonTaxedMinorUnits: value.NonTaxedMinorUnits, OtherTaxMinorUnits: value.OtherTaxMinorUnits,
		IssuedOn: date(value.IssuedOn), ServiceFrom: optionalDate(value.ServiceFrom), ServiceUntil: optionalDate(value.ServiceUntil), PaymentDueOn: optionalDate(value.PaymentDueOn),
		VATLines: append([]fiscal.VATLine{}, value.VATLines...), OtherTaxLines: append([]fiscal.OtherTaxLine{}, value.OtherTaxLines...),
		AssociatedVouchers: associated,
	}
}

func date(value time.Time) string { return value.UTC().Format("2006-01-02") }
func optionalDate(value *time.Time) *string {
	if value == nil {
		return nil
	}
	result := date(*value)
	return &result
}
func validHash(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, character := range value {
		if !(character >= '0' && character <= '9') && !(character >= 'a' && character <= 'f') {
			return false
		}
	}
	return true
}
