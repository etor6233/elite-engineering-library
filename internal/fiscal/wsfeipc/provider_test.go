package wsfeipc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/fiscal"
)

func TestUnixProviderExactSchemaAndFailClosedResponses(t *testing.T) {
	socket := filepath.Join(os.TempDir(), fmt.Sprintf("elite-wsfe-%d.sock", os.Getpid()))
	_ = os.Remove(socket)
	listener, err := net.Listen("unix", socket)
	if err != nil {
		t.Fatalf("unix listener unavailable: %v", err)
	}
	defer func() { listener.Close(); _ = os.Remove(socket) }()
	server := &http.Server{Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var body map[string]any
		decoder := json.NewDecoder(request.Body)
		if err := decoder.Decode(&body); err != nil {
			t.Error(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, found := body["tenant_id"]; found {
			t.Error("tenant leaked to bridge")
		}
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/v1/last-authorized":
			assertInvoiceBody(t, body)
			_, _ = writer.Write([]byte(`{"number":40,"response_hash":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","provider_codes":[]}`))
		case "/v1/consult":
			assertInvoiceBody(t, body)
			_, _ = writer.Write([]byte(`{"found":false,"authorized":false,"cae":null,"cae_expires_on":null,"response_hash":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","provider_codes":["E:602"]}`))
		case "/v1/authorize":
			assertInvoiceBody(t, body)
			_, _ = writer.Write([]byte(`{"found":true,"authorized":true,"cae":"41124578989845","cae_expires_on":"2010-09-13","response_hash":"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc","provider_codes":[]}`))
		case "/v1/parameters":
			if body["taxpayer_cuit"] != "33693450239" || body["kind"] != "vat_rate" || body["voucher_class"] != nil {
				t.Errorf("parameter scope changed: %+v", body)
			}
			_, _ = writer.Write([]byte(`{"kind":"vat_rate","voucher_class":null,"items":[{"code":"5","description":"21%","valid_from":"2009-02-20","valid_until":null,"voucher_class":null,"emission_type":null,"blocked":null,"deregistered_on":null}],"response_hash":"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd","provider_codes":[]}`))
		default:
			writer.WriteHeader(http.StatusNotFound)
		}
	})}
	go func() { _ = server.Serve(listener) }()
	defer server.Shutdown(context.Background())
	provider, err := NewUnixProvider(socket, 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	invoice := fixtureInvoice()
	last, hash, err := provider.LastAuthorized(context.Background(), invoice)
	if err != nil || last != 40 || !strings.HasPrefix(hash, "aaaa") {
		t.Fatalf("last=%d hash=%s err=%v", last, hash, err)
	}
	consulted, err := provider.Consult(context.Background(), invoice)
	if err != nil || validateAuthorization(consulted) != nil {
		t.Fatalf("consulted=%+v err=%v", consulted, err)
	}
	authorized, err := provider.Authorize(context.Background(), invoice)
	if err != nil || !authorized.Found || !authorized.Authorized || authorized.CAE != "41124578989845" || authorized.CAEExpiresOn == nil || authorized.ProviderCodes == nil {
		t.Fatalf("authorized=%+v err=%v", authorized, err)
	}
	parameters, err := provider.FetchParameters(context.Background(), "33693450239", "vat_rate", nil)
	if err != nil || parameters.Kind != "vat_rate" || len(parameters.Items) != 1 || parameters.Items[0].Code != "5" || !strings.HasPrefix(parameters.ResponseHash, "dddd") {
		t.Fatalf("parameters=%+v err=%v", parameters, err)
	}
}

func assertInvoiceBody(t *testing.T, body map[string]any) {
	t.Helper()
	if body["recipient_vat_condition_id"] != float64(1) || len(body["vat_lines"].([]any)) != 2 || len(body["other_tax_lines"].([]any)) != 1 {
		t.Errorf("exact fiscal details missing: %+v", body)
	}
	associated := body["associated_vouchers"].([]any)
	if len(associated) != 1 || associated[0].(map[string]any)["voucher_type"] != float64(6) || associated[0].(map[string]any)["number"] != float64(40) || associated[0].(map[string]any)["issued_on"] != "2010-09-02" {
		t.Errorf("associated voucher identity missing: %+v", associated)
	}
}

func TestUnixProviderRejectsUnsafeConfigurationAndProtocol(t *testing.T) {
	if _, err := NewUnixProvider("relative.sock", time.Second); !errors.Is(err, fiscal.ErrInvalid) {
		t.Fatalf("unsafe socket accepted: %v", err)
	}
	provider := &Provider{client: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: ioNopCloser(`{"number":1,"response_hash":"bad","provider_codes":[],"unknown":true}`), Header: make(http.Header)}, nil
	})}, base: "http://test"}
	if _, _, err := provider.LastAuthorized(context.Background(), fixtureInvoice()); !errors.Is(err, ErrProtocol) {
		t.Fatalf("unknown/invalid response accepted: %v", err)
	}
}

func fixtureInvoice() fiscal.Invoice {
	return fiscal.Invoice{TaxpayerCUIT: "33693450239", PointOfSaleNumber: 12, VoucherType: 8, Concept: 1, RecipientDocumentType: 80, RecipientDocument: "20111111112", RecipientVATConditionID: 1, Currency: "ARS", VoucherNumber: 1, TotalMinorUnits: 18405, NetMinorUnits: 15000, VATMinorUnits: 2625, OtherTaxMinorUnits: 780, IssuedOn: time.Date(2010, 9, 3, 0, 0, 0, 0, time.UTC), VATLines: []fiscal.VATLine{{ID: 5, BaseMinorUnits: 10000, AmountMinorUnits: 2100}, {ID: 4, BaseMinorUnits: 5000, AmountMinorUnits: 525}}, OtherTaxLines: []fiscal.OtherTaxLine{{ID: 99, Description: "Impuesto Municipal Matanza", BaseMinorUnits: 15000, RateBasisPoints: 520, AmountMinorUnits: 780}}, AssociatedVouchers: []fiscal.AssociatedVoucher{{InvoiceID: "original", TaxpayerCUIT: "33693450239", VoucherType: 6, PointOfSale: 12, Number: 40, IssuedOn: time.Date(2010, 9, 2, 0, 0, 0, 0, time.UTC)}}}
}

func validateAuthorization(value fiscal.Authorization) error {
	if value.Found || value.Authorized || value.CAE != "" || value.CAEExpiresOn != nil || !strings.HasPrefix(value.ResponseHash, "bbbb") {
		return errors.New("invalid not-found authorization")
	}
	return nil
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (function roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

type stringReadCloser struct{ *strings.Reader }

func (stringReadCloser) Close() error { return nil }

func ioNopCloser(value string) *stringReadCloser { return &stringReadCloser{strings.NewReader(value)} }
