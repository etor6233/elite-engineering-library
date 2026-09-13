# ARCA WSFE Unix-Socket Fiscal Worker

## 1. Metadata

```yaml
pack_id: "ARCA-WSFE-UDS-WORKER"
pack_version: "0.3.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Conecta emisión, comprobantes asociados y siete consultas paramétricas del owner fiscal Go con WSAA/WSFE .NET por Unix socket local, sin TCP, con JSON estricto, secretos fuera y aprobación separada."
stacks: ["Go 1.26.8", ".NET SDK 10.0.400", "ASP.NET Core Kestrel UDS", "WCF 10.0.652802", "PostgreSQL 18.6", "ARCA WSFEv1 hash-locked homologation WSDL"]
compatible_with: ["MICROSOFT-ARCA-WSFE-GENERATED-CLIENT 0.2.x", "MICROSOFT-ARCA-WSAA-CREDENTIAL-CORE 0.1.x", "MICROSOFT-ARCA-WSFE-SOAP-ADAPTER 0.3.x", "GO-ARCA-FISCAL-ISSUANCE-API 0.6.x", "GO-ELECTROMOBILITY-APPLICATION 1.5.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner AND Microsoft dependency licenses AND ARCA public specification"
upstream_sources: ["https://learn.microsoft.com/en-us/dotnet/api/microsoft.aspnetcore.server.kestrel.core.kestrelserveroptions.listenunixsocket?view=aspnetcore-10.0", "https://learn.microsoft.com/en-us/aspnet/core/fundamentals/servers/kestrel/security-considerations?view=aspnetcore-10.0", "https://pkg.go.dev/net", "https://github.com/dotnet/wcf", "https://www.arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG.pdf"]
verified_at: "2026-09-13"
```

Los trece bloques son `AUTHORED`: componen APIs públicas oficiales, el proxy generado y packs previamente verificados; no son código copiado de Microsoft, Google, ARCA o Go. No contienen WSDL, proxy generado, certificado, clave privada, Token, Sign, credenciales ni mensajes del proveedor.

## 2. Applicability

Use después de materializar el owner fiscal 0.4.x, aplicación 1.4.x, credential core, adapter y un cliente oficial generado de homologación. Adóptelo cuando Go deba emitir o adquirir parámetros sin incorporar WCF al proceso HTTP ni exponer el bridge por red. Rechácelo ante producción, TCP público, certificado en archivo, tablas inventadas o aprobación automática. El claim termina en IPC local reconstruido y probado; no demuestra credenciales, asociación, homologación live ni aprobación fiscal.

## 3. Architecture contract

PostgreSQL y `fiscal.Processor` conservan emisión; `ParameterRegistry` conserva snapshots y decisión separada. Go llama emisión y `/v1/parameters` por `net.Dialer.DialContext("unix", ...)`; no reintenta HTTP dentro del adapter. Kestrel escucha sólo UDS absoluto, limita requests a 64 KiB y JSON rechaza campos desconocidos/tenant/secrets. .NET despacha sólo siete kinds cerrados a operaciones generadas, admite homologación fija, carga la clave por thumbprint, limita SOAP a 1 MiB/30 s y sanitiza errores. Producción, valores inventados, aprobación automática y cambio de endpoint permanecen bloqueados.

## 4. Exact file manifest

```text
CREATE arca/fiscal/worker/fixtures/Elite.Arca.Wsfe.Fixture/Elite.Arca.Wsfe.Fixture.csproj
CREATE arca/fiscal/worker/fixtures/Elite.Arca.Wsfe.Fixture/FixtureParameters.cs
CREATE arca/fiscal/worker/fixtures/Elite.Arca.Wsfe.Fixture/Program.cs
CREATE arca/fiscal/worker/fixtures/Elite.Arca.Wsfe.Fixture/packages.lock.json
CREATE docs/ARCA_LOCAL_REFERENCE.md
CREATE internal/fiscal/wsfeipc/serialization_test.go
CREATE cmd/arca-fiscal-worker/main.go
CREATE internal/fiscal/wsfeipc/provider.go
CREATE internal/fiscal/wsfeipc/provider_test.go
CREATE arca/fiscal/worker/README.md
CREATE arca/fiscal/worker/Elite.Arca.Wsfe.Worker/Elite.Arca.Wsfe.Worker.csproj
CREATE arca/fiscal/worker/Elite.Arca.Wsfe.Worker/packages.lock.json
CREATE arca/fiscal/worker/Elite.Arca.Wsfe.Worker/Program.cs
CREATE arca/fiscal/worker/Elite.Arca.Wsfe.Worker/SoapTransports.cs
CREATE arca/fiscal/worker/Elite.Arca.Wsfe.Worker/WorkerHost.cs
CREATE arca/fiscal/worker/tests/Elite.Arca.Wsfe.Worker.Tests/Elite.Arca.Wsfe.Worker.Tests.csproj
CREATE arca/fiscal/worker/tests/Elite.Arca.Wsfe.Worker.Tests/packages.lock.json
CREATE arca/fiscal/worker/tests/Elite.Arca.Wsfe.Worker.Tests/Program.cs
CREATE tools/test-arca-wsfe-uds-worker.ps1
```

## 5. Materialization blocks

### FILE: `cmd/arca-fiscal-worker/main.go`
```yaml
block_id: "ARCA-WSFE-UDS-WORKER:go-command:v1"
operation: CREATE
path: "cmd/arca-fiscal-worker/main.go"
sha256: "df53732678241d8b3063bffb0524ec0e6e67c7014f41905bedc00aa768c0bf31"
provenance: AUTHORED
source: "local worker governed by Go net and database/sql-style cancellation contracts"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````go
package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"elite.local/enterprise/internal/fiscal"
	"elite.local/enterprise/internal/fiscal/wsfeipc"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var workerIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$`)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	databaseURL := os.Getenv("DATABASE_URL")
	socketPath := os.Getenv("ARCA_WSFE_SOCKET")
	workerID := os.Getenv("FISCAL_WORKER_ID")
	if databaseURL == "" || socketPath == "" || !workerIDPattern.MatchString(workerID) {
		slog.Error("DATABASE_URL, ARCA_WSFE_SOCKET and a valid FISCAL_WORKER_ID are required")
		os.Exit(2)
	}

	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		slog.Error("database configuration is invalid")
		os.Exit(2)
	}
	poolConfig.MaxConns = 4
	poolConfig.MinConns = 0
	poolConfig.MaxConnLifetime = 30 * time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		slog.Error("database pool creation failed")
		os.Exit(1)
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		slog.Error("database unavailable")
		os.Exit(1)
	}

	provider, err := wsfeipc.NewUnixProvider(socketPath, 30*time.Second)
	if err != nil {
		slog.Error("ARCA WSFE socket configuration is invalid")
		os.Exit(2)
	}
	processor, err := fiscal.NewProcessor(postgres.NewFiscal(pool), provider, randomid.Generator{}, workerID, 2*time.Minute)
	if err != nil {
		slog.Error("fiscal processor configuration is invalid")
		os.Exit(2)
	}

	if err = run(ctx, processor); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("fiscal worker stopped unexpectedly")
		os.Exit(1)
	}
}

func run(ctx context.Context, processor *fiscal.Processor) error {
	idle := time.NewTimer(0)
	if !idle.Stop() {
		<-idle.C
	}
	defer idle.Stop()
	for {
		_, err := processor.ProcessOne(ctx)
		if err == nil {
			continue
		}
		if errors.Is(err, context.Canceled) {
			return err
		}
		wait := 2 * time.Second
		if errors.Is(err, fiscal.ErrNoWork) {
			wait = 500 * time.Millisecond
		} else {
			slog.Warn("fiscal work deferred after a recoverable processing failure")
		}
		idle.Reset(wait)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-idle.C:
		}
	}
}
````

### FILE: `internal/fiscal/wsfeipc/provider.go`
```yaml
block_id: "ARCA-WSFE-UDS-WORKER:go-provider:v1"
operation: CREATE
path: "internal/fiscal/wsfeipc/provider.go"
sha256: "bc9a92f73e025ec60bf90a02c73886461288414317d967b7c304c0dc903db577"
provenance: AUTHORED
source: "local adapter governed by official Go net.Dialer DialContext documentation"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````go
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
````

### FILE: `internal/fiscal/wsfeipc/provider_test.go`
```yaml
block_id: "ARCA-WSFE-UDS-WORKER:go-provider-tests:v1"
operation: CREATE
path: "internal/fiscal/wsfeipc/provider_test.go"
sha256: "70ed3aae1f6c90739c53e373d44854c81341b0139245e62bab21ea46ba537cda"
provenance: AUTHORED
source: "local Unix-socket contract tests"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````go
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
````

### FILE: `arca/fiscal/worker/README.md`
```yaml
block_id: "ARCA-WSFE-UDS-WORKER:readme:v1"
operation: CREATE
path: "arca/fiscal/worker/README.md"
sha256: "7e315b51ef64b3587b5943dd7ff62b0615817dbc69775e3fd64d214784b849ec"
provenance: AUTHORED
source: "local runtime contract governed by Microsoft Kestrel Unix socket guidance"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````text
# ARCA WSFE local Unix-socket worker

This process is the narrow boundary between the Go fiscal owner and the Microsoft-generated WSAA/WSFE clients. It listens only on an absolute Unix-domain socket, accepts a strict snake-case JSON contract, invokes homologation endpoints and returns only safe result fields. It has no TCP listener and never returns or logs Token, Sign, certificate material, private keys or provider messages.

Required runtime inputs are `ARCA_ENVIRONMENT=homologation`, `ARCA_WSFE_SOCKET`, `ARCA_CERTIFICATE_THUMBPRINT` and `ARCA_CERTIFICATE_STORE=CurrentUser|LocalMachine`. Provision the socket directory with least privilege before startup and ensure the socket path does not exist. The process creates the socket and sets mode `0660` on Unix; Kestrel does not create or secure the parent directory.

Production is deliberately blocked. Before promotion, generate current official clients, prove certificate association and current ARCA parameter tables, run the owner/worker reconciliation flow against homologation, approve tax/accounting policy, and execute deployment, monitoring and rollback gates for the selected runtime.
````

### FILE: `arca/fiscal/worker/Elite.Arca.Wsfe.Worker/Elite.Arca.Wsfe.Worker.csproj`
```yaml
block_id: "ARCA-WSFE-UDS-WORKER:worker-project:v1"
operation: CREATE
path: "arca/fiscal/worker/Elite.Arca.Wsfe.Worker/Elite.Arca.Wsfe.Worker.csproj"
sha256: "68c57bf0bb38a34e90e5e2d8272d17e090660b0e0519ac6ec4b153c12e31942e"
provenance: AUTHORED
source: "local project governed by Microsoft .NET 10 SDK"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````xml
<Project Sdk="Microsoft.NET.Sdk.Web">
  <PropertyGroup>
    <OutputType>Exe</OutputType>
    <TargetFramework>net10.0</TargetFramework>
    <ImplicitUsings>enable</ImplicitUsings>
    <Nullable>enable</Nullable>
    <TreatWarningsAsErrors>true</TreatWarningsAsErrors>
    <AnalysisLevel>latest</AnalysisLevel>
    <EnforceCodeStyleInBuild>true</EnforceCodeStyleInBuild>
    <RestorePackagesWithLockFile>true</RestorePackagesWithLockFile>
    <NuGetAudit>true</NuGetAudit>
    <NuGetAuditMode>all</NuGetAuditMode>
  </PropertyGroup>
  <ItemGroup>
    <ProjectReference Include="../../bridge/src/Elite.Arca.Wsfe.Bridge/Elite.Arca.Wsfe.Bridge.csproj" />
    <ProjectReference Include="../../generated/Elite.Arca.Wsfe.Generated.csproj" />
    <ProjectReference Include="../../../credentials/src/Elite.Arca.Credentials/Elite.Arca.Credentials.csproj" />
  </ItemGroup>
</Project>
````

### FILE: `arca/fiscal/worker/Elite.Arca.Wsfe.Worker/packages.lock.json`
```yaml
block_id: "ARCA-WSFE-UDS-WORKER:worker-lock:v1"
operation: CREATE
path: "arca/fiscal/worker/Elite.Arca.Wsfe.Worker/packages.lock.json"
sha256: "cc28c4c6193a2f96790da7d6b178ffbb6c5788377dfd2a6ed22e57b3f187d82d"
provenance: AUTHORED
source: "lock generated by Microsoft .NET SDK 10.0.400"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````json
{
  "version": 1,
  "dependencies": {
    "net10.0": {
      "System.Security.Cryptography.Pkcs": {
        "type": "Transitive",
        "resolved": "10.0.11",
        "contentHash": "8IV+rI3xN/Mkq9MsSX7VZTv9T9Wt+tyhHLwBX9VNZBk8m8kknEs1JeTSjI2x9M2L/fSja+c2WSE5n5Yy0GXdoQ=="
      },
      "System.ServiceModel.Http": {
        "type": "Transitive",
        "resolved": "10.0.652802",
        "contentHash": "G02XZvmccf42QCU5MjviBIg69MSMAVHwL1inVPsNSpfp5g+t5BkQM3DyvWRLN4qmeFDWSF/mA1rIYONIDu/6Dg==",
        "dependencies": {
          "System.ServiceModel.Primitives": "10.0.652802"
        }
      },
      "System.ServiceModel.Primitives": {
        "type": "Transitive",
        "resolved": "10.0.652802",
        "contentHash": "ULfGNl75BNXkpF42wNV2CDXJ64dUZZEa8xO2mBsc4tqbW9QjruxjEB6bAr4Z/T1rNU+leOztIjCJQYsBGFWYlw=="
      },
      "elite.arca.credentials": {
        "type": "Project",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "[10.0.11, )"
        }
      },
      "elite.arca.wsfe.bridge": {
        "type": "Project",
        "dependencies": {
          "Elite.Arca.Credentials": "[1.0.0, )",
          "Elite.Arca.Wsfe.Generated": "[1.0.0, )"
        }
      },
      "elite.arca.wsfe.generated": {
        "type": "Project",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "[10.0.11, )",
          "System.ServiceModel.Http": "[10.0.652802, )"
        }
      }
    }
  }
}
````

### FILE: `arca/fiscal/worker/Elite.Arca.Wsfe.Worker/Program.cs`
```yaml
block_id: "ARCA-WSFE-UDS-WORKER:worker-program:v1"
operation: CREATE
path: "arca/fiscal/worker/Elite.Arca.Wsfe.Worker/Program.cs"
sha256: "d58967f975d85da77e853c40c4c1b1d3d3854288be3aa61872f6b6302ab83aed"
provenance: AUTHORED
source: "local composition governed by ARCA homologation and Microsoft WCF contracts"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````csharp
using System.Security.Cryptography.X509Certificates;
using Elite.Arca.Credentials;
using Elite.Arca.Wsfe.Bridge;

namespace Elite.Arca.Wsfe.Worker;

public static class Program
{
    private static readonly Uri WsaaHomologation = new("https://wsaahomo.afip.gov.ar/ws/services/LoginCms");
    private static readonly Uri WsfeHomologation = new("https://wswhomo.afip.gov.ar/wsfev1/service.asmx");

    public static async Task<int> Main(string[] args)
    {
        try
        {
            if (!string.Equals(Environment.GetEnvironmentVariable("ARCA_ENVIRONMENT"), "homologation", StringComparison.Ordinal))
                throw new InvalidOperationException("Only the ARCA homologation environment is admitted by this worker.");
            var socketPath = Required("ARCA_WSFE_SOCKET");
            var thumbprint = Required("ARCA_CERTIFICATE_THUMBPRINT");
            var storeLocation = Required("ARCA_CERTIFICATE_STORE") switch
            {
                "CurrentUser" => StoreLocation.CurrentUser,
                "LocalMachine" => StoreLocation.LocalMachine,
                _ => throw new InvalidOperationException("ARCA_CERTIFICATE_STORE must be CurrentUser or LocalMachine.")
            };
            using var certificate = CertificateStoreLoader.LoadByThumbprint(thumbprint, storeLocation);
            using var credentials = new WsaaCredentialProvider(
                WsaaOptions.ElectronicInvoice,
                new SystemUtcClock(),
                new LoginTicketRequestFactory(WsaaOptions.ElectronicInvoice, new SystemUtcClock(), new MonotonicLoginTicketIdSource()),
                new CmsRequestSigner(),
                certificate,
                new WsaaSoapTransport(WsaaHomologation));
            var access = new CredentialAccessSource(credentials);
            var transport = new WsfeSoapTransport(WsfeHomologation);
            var operations = new BridgeFiscalOperations(new WsfeBridge(access, transport), new WsfeParameterBridge(access, transport));
            await using var app = WorkerHost.Build(args, socketPath, operations);
            await app.StartAsync().ConfigureAwait(false);
            WorkerHost.SecureSocket(socketPath);
            await app.WaitForShutdownAsync().ConfigureAwait(false);
            return 0;
        }
        catch
        {
            Console.Error.WriteLine("ARCA WSFE worker startup or runtime failed; inspect controlled platform diagnostics.");
            return 1;
        }
    }

    private static string Required(string name) => Environment.GetEnvironmentVariable(name) is { Length: > 0 } value ? value : throw new InvalidOperationException($"{name} is required.");
}
````

### FILE: `arca/fiscal/worker/Elite.Arca.Wsfe.Worker/SoapTransports.cs`
```yaml
block_id: "ARCA-WSFE-UDS-WORKER:soap-transports:v1"
operation: CREATE
path: "arca/fiscal/worker/Elite.Arca.Wsfe.Worker/SoapTransports.cs"
sha256: "99a623004e35812b2fbd8ee09105ce118c15d60c5e37811469637ec4337af86b"
provenance: AUTHORED
source: "local bounded transport governed by Microsoft WCF client APIs"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````csharp
using System.Security.Cryptography;
using System.ServiceModel;
using System.ServiceModel.Channels;
using Elite.Arca.Credentials;
using Elite.Arca.Wsfe.Bridge;
using WsaaGenerated = Elite.Arca.Wsaa.Generated;
using WsfeGenerated = Elite.Arca.Wsfe.Generated;

namespace Elite.Arca.Wsfe.Worker;

public sealed class MonotonicLoginTicketIdSource : ILoginTicketIdSource
{
    private long current = RandomNumberGenerator.GetInt32(1, int.MaxValue - 1_000_000);

    public long Next()
    {
        var value = Interlocked.Increment(ref current);
        if (value >= int.MaxValue) throw new InvalidOperationException("Login ticket ID space exhausted; restart is required.");
        return value;
    }
}

public sealed class WsaaSoapTransport(Uri endpoint) : IWsaaTransport
{
    public async Task<string> LoginCmsAsync(string cmsBase64, CancellationToken cancellationToken)
    {
        var client = new WsaaGenerated.LoginCMSClient(SoapBinding(), new EndpointAddress(endpoint));
        try
        {
            var response = await client.loginCmsAsync(cmsBase64).WaitAsync(cancellationToken).ConfigureAwait(false);
            if (string.IsNullOrWhiteSpace(response.loginCmsReturn)) throw new InvalidDataException("WSAA returned an empty response.");
            Close(client);
            return response.loginCmsReturn;
        }
        catch { Abort(client); throw; }
    }

    internal static BasicHttpBinding SoapBinding() => new(BasicHttpSecurityMode.Transport)
    {
        OpenTimeout = TimeSpan.FromSeconds(10),
        CloseTimeout = TimeSpan.FromSeconds(5),
        SendTimeout = TimeSpan.FromSeconds(30),
        ReceiveTimeout = TimeSpan.FromSeconds(30),
        MaxBufferSize = 1 << 20,
        MaxReceivedMessageSize = 1 << 20,
        TransferMode = TransferMode.Buffered,
        ReaderQuotas = new System.Xml.XmlDictionaryReaderQuotas
        {
            MaxDepth = 32,
            MaxStringContentLength = 1 << 20,
            MaxArrayLength = 65_536,
            MaxBytesPerRead = 4_096,
            MaxNameTableCharCount = 16_384
        }
    };

    internal static void Close(ICommunicationObject client)
    {
        try
        {
            if (client.State == CommunicationState.Faulted) client.Abort();
            else client.Close();
        }
        catch { client.Abort(); }
    }

    internal static void Abort(ICommunicationObject client)
    {
        try { client.Abort(); } catch { }
    }
}

public sealed class WsfeSoapTransport(Uri endpoint) : IWsfeSoap, IWsfeParameterSoap
{
    public Task<WsfeGenerated.FERecuperaLastCbteResponse> LastAuthorizedAsync(WsfeGenerated.FEAuthRequest auth, int pointOfSale, int voucherType, CancellationToken cancellationToken) =>
        Invoke(client => client.FECompUltimoAutorizadoAsync(auth, pointOfSale, voucherType), cancellationToken);

    public Task<WsfeGenerated.FECompConsultaResponse> ConsultAsync(WsfeGenerated.FEAuthRequest auth, WsfeGenerated.FECompConsultaReq request, CancellationToken cancellationToken) =>
        Invoke(client => client.FECompConsultarAsync(auth, request), cancellationToken);

    public Task<WsfeGenerated.FECAEResponse> AuthorizeAsync(WsfeGenerated.FEAuthRequest auth, WsfeGenerated.FECAERequest request, CancellationToken cancellationToken) =>
        Invoke(client => client.FECAESolicitarAsync(auth, request), cancellationToken);

    public Task<WsfeGenerated.CbteTipoResponse> VoucherTypesAsync(WsfeGenerated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        Invoke(client => client.FEParamGetTiposCbteAsync(auth), cancellationToken);
    public Task<WsfeGenerated.ConceptoTipoResponse> ConceptsAsync(WsfeGenerated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        Invoke(client => client.FEParamGetTiposConceptoAsync(auth), cancellationToken);
    public Task<WsfeGenerated.DocTipoResponse> DocumentTypesAsync(WsfeGenerated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        Invoke(client => client.FEParamGetTiposDocAsync(auth), cancellationToken);
    public Task<WsfeGenerated.IvaTipoResponse> VatRatesAsync(WsfeGenerated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        Invoke(client => client.FEParamGetTiposIvaAsync(auth), cancellationToken);
    public Task<WsfeGenerated.FETributoResponse> OtherTaxesAsync(WsfeGenerated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        Invoke(client => client.FEParamGetTiposTributosAsync(auth), cancellationToken);
    public Task<WsfeGenerated.FEPtoVentaResponse> PointsOfSaleAsync(WsfeGenerated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        Invoke(client => client.FEParamGetPtosVentaAsync(auth), cancellationToken);
    public Task<WsfeGenerated.CondicionIvaReceptorResponse> RecipientVatConditionsAsync(WsfeGenerated.FEAuthRequest auth, string voucherClass, CancellationToken cancellationToken) =>
        Invoke(client => client.FEParamGetCondicionIvaReceptorAsync(auth, voucherClass), cancellationToken);

    private async Task<T> Invoke<T>(Func<WsfeGenerated.ServiceSoapClient, Task<T>> call, CancellationToken cancellationToken)
    {
        var client = new WsfeGenerated.ServiceSoapClient(WsaaSoapTransport.SoapBinding(), new EndpointAddress(endpoint));
        try
        {
            var response = await call(client).WaitAsync(cancellationToken).ConfigureAwait(false);
            WsaaSoapTransport.Close(client);
            return response;
        }
        catch { WsaaSoapTransport.Abort(client); throw; }
    }
}
````

### FILE: `arca/fiscal/worker/Elite.Arca.Wsfe.Worker/WorkerHost.cs`
```yaml
block_id: "ARCA-WSFE-UDS-WORKER:worker-host:v1"
operation: CREATE
path: "arca/fiscal/worker/Elite.Arca.Wsfe.Worker/WorkerHost.cs"
sha256: "971d1aaeb4dcced501fdeccd7cf4e3dca5996a76efa459fe0e7f6484eac6b206"
provenance: AUTHORED
source: "local UDS host governed by Microsoft Kestrel security guidance"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````csharp
using System.Text.Json.Serialization;
using Elite.Arca.Wsfe.Bridge;
using Microsoft.AspNetCore.Http.Json;

namespace Elite.Arca.Wsfe.Worker;

public interface IFiscalOperations
{
    Task<LastAuthorizedResult> LastAuthorizedAsync(FiscalInvoice invoice, CancellationToken cancellationToken);
    Task<SafeWsfeResult> ConsultAsync(FiscalInvoice invoice, CancellationToken cancellationToken);
    Task<SafeWsfeResult> AuthorizeAsync(FiscalInvoice invoice, CancellationToken cancellationToken);
    Task<SafeParameterResult> ParametersAsync(ParameterRequest request, CancellationToken cancellationToken);
}

public sealed record ParameterRequest(string TaxpayerCuit, string Kind, string? VoucherClass);

public sealed class BridgeFiscalOperations(WsfeBridge bridge, WsfeParameterBridge parameters) : IFiscalOperations
{
    public Task<LastAuthorizedResult> LastAuthorizedAsync(FiscalInvoice invoice, CancellationToken cancellationToken) => bridge.LastAuthorizedAsync(invoice, cancellationToken);
    public Task<SafeWsfeResult> ConsultAsync(FiscalInvoice invoice, CancellationToken cancellationToken) => bridge.ConsultAsync(invoice, cancellationToken);
    public Task<SafeWsfeResult> AuthorizeAsync(FiscalInvoice invoice, CancellationToken cancellationToken) => bridge.AuthorizeAsync(invoice, cancellationToken);
    public Task<SafeParameterResult> ParametersAsync(ParameterRequest request, CancellationToken cancellationToken) => request.Kind switch
    {
        "voucher_type" when request.VoucherClass is null => parameters.VoucherTypesAsync(request.TaxpayerCuit, cancellationToken),
        "concept" when request.VoucherClass is null => parameters.ConceptsAsync(request.TaxpayerCuit, cancellationToken),
        "document_type" when request.VoucherClass is null => parameters.DocumentTypesAsync(request.TaxpayerCuit, cancellationToken),
        "vat_rate" when request.VoucherClass is null => parameters.VatRatesAsync(request.TaxpayerCuit, cancellationToken),
        "other_tax" when request.VoucherClass is null => parameters.OtherTaxesAsync(request.TaxpayerCuit, cancellationToken),
        "point_of_sale" when request.VoucherClass is null => parameters.PointsOfSaleAsync(request.TaxpayerCuit, cancellationToken),
        "recipient_vat_condition" when request.VoucherClass is not null => parameters.RecipientVatConditionsAsync(request.TaxpayerCuit, request.VoucherClass, cancellationToken),
        _ => throw new ArgumentException("Unsupported parameter request.", nameof(request))
    };
}

public static class WorkerHost
{
    public static WebApplication Build(string[] args, string socketPath, IFiscalOperations operations)
    {
        ArgumentNullException.ThrowIfNull(operations);
        ValidateSocket(socketPath);
        var builder = WebApplication.CreateSlimBuilder(args);
        builder.Logging.ClearProviders();
        builder.Services.Configure<JsonOptions>(options =>
        {
            options.SerializerOptions.PropertyNamingPolicy = System.Text.Json.JsonNamingPolicy.SnakeCaseLower;
            options.SerializerOptions.DictionaryKeyPolicy = System.Text.Json.JsonNamingPolicy.SnakeCaseLower;
            options.SerializerOptions.UnmappedMemberHandling = JsonUnmappedMemberHandling.Disallow;
        });
        builder.WebHost.ConfigureKestrel(options =>
        {
            options.AddServerHeader = false;
            options.Limits.MaxRequestBodySize = 64 << 10;
            options.Limits.MaxRequestHeaderCount = 32;
            options.Limits.MaxRequestHeadersTotalSize = 16 << 10;
            options.Limits.RequestHeadersTimeout = TimeSpan.FromSeconds(5);
            options.Limits.KeepAliveTimeout = TimeSpan.FromSeconds(30);
            options.ListenUnixSocket(socketPath);
        });
        var app = builder.Build();
        app.MapGet("/health/live", static () => Results.Ok(new { status = "alive" }));
        app.MapPost("/v1/last-authorized", (FiscalInvoice invoice, CancellationToken token) => Execute(() => operations.LastAuthorizedAsync(invoice, token)));
        app.MapPost("/v1/consult", (FiscalInvoice invoice, CancellationToken token) => Execute(() => operations.ConsultAsync(invoice, token)));
        app.MapPost("/v1/authorize", (FiscalInvoice invoice, CancellationToken token) => Execute(() => operations.AuthorizeAsync(invoice, token)));
        app.MapPost("/v1/parameters", (ParameterRequest request, CancellationToken token) => Execute(() => operations.ParametersAsync(request, token)));
        return app;
    }

    public static void SecureSocket(string socketPath)
    {
        if (!Path.Exists(socketPath)) throw new IOException("Kestrel did not create the Unix socket.");
        if (!OperatingSystem.IsWindows())
            File.SetUnixFileMode(socketPath, UnixFileMode.UserRead | UnixFileMode.UserWrite | UnixFileMode.GroupRead | UnixFileMode.GroupWrite);
    }

    private static async Task<IResult> Execute<T>(Func<Task<T>> operation)
    {
        try { return Results.Ok(await operation().ConfigureAwait(false)); }
        catch (WsfeResponseException error)
        {
            return Results.Json(new SafeFailure("provider_rejected", error.ResponseHash, error.ProviderCodes), statusCode: StatusCodes.Status502BadGateway);
        }
        catch (ArgumentException) { return Results.Json(new SafeFailure("invalid_request", null, []), statusCode: StatusCodes.Status400BadRequest); }
        catch (OperationCanceledException) { return Results.Json(new SafeFailure("cancelled", null, []), statusCode: 499); }
        catch { return Results.Json(new SafeFailure("provider_unavailable", null, []), statusCode: StatusCodes.Status503ServiceUnavailable); }
    }

    private static void ValidateSocket(string socketPath)
    {
        if (string.IsNullOrWhiteSpace(socketPath) || !Path.IsPathFullyQualified(socketPath)) throw new ArgumentException("An absolute Unix socket path is required.", nameof(socketPath));
        var parent = Path.GetDirectoryName(socketPath);
        if (parent is null || !Directory.Exists(parent)) throw new DirectoryNotFoundException("The Unix socket parent directory must be provisioned before startup.");
        if (Path.Exists(socketPath)) throw new IOException("The Unix socket path already exists; stale sockets must be investigated before startup.");
        if ((File.GetAttributes(parent) & FileAttributes.ReparsePoint) != 0) throw new IOException("The Unix socket parent directory cannot be a reparse point.");
    }

    private sealed record SafeFailure(string Code, string? ResponseHash, IReadOnlyList<string> ProviderCodes);
}
````

### FILE: `arca/fiscal/worker/tests/Elite.Arca.Wsfe.Worker.Tests/Elite.Arca.Wsfe.Worker.Tests.csproj`
```yaml
block_id: "ARCA-WSFE-UDS-WORKER:test-project:v1"
operation: CREATE
path: "arca/fiscal/worker/tests/Elite.Arca.Wsfe.Worker.Tests/Elite.Arca.Wsfe.Worker.Tests.csproj"
sha256: "f1f2e0664386aec472e683629d97f530e23c1cf3391c1afea706ceae69d0cb1a"
provenance: AUTHORED
source: "local executable contract test project"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````xml
<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <OutputType>Exe</OutputType>
    <TargetFramework>net10.0</TargetFramework>
    <ImplicitUsings>enable</ImplicitUsings>
    <Nullable>enable</Nullable>
    <TreatWarningsAsErrors>true</TreatWarningsAsErrors>
    <AnalysisLevel>latest</AnalysisLevel>
    <RestorePackagesWithLockFile>true</RestorePackagesWithLockFile>
    <NuGetAudit>true</NuGetAudit>
    <NuGetAuditMode>all</NuGetAuditMode>
  </PropertyGroup>
  <ItemGroup>
    <ProjectReference Include="../../../worker/Elite.Arca.Wsfe.Worker/Elite.Arca.Wsfe.Worker.csproj" />
  </ItemGroup>
</Project>
````

### FILE: `arca/fiscal/worker/tests/Elite.Arca.Wsfe.Worker.Tests/packages.lock.json`
```yaml
block_id: "ARCA-WSFE-UDS-WORKER:test-lock:v1"
operation: CREATE
path: "arca/fiscal/worker/tests/Elite.Arca.Wsfe.Worker.Tests/packages.lock.json"
sha256: "55e434a478b434a6a2180f8ab587263657f990b72441b8eff1c38004b29fa431"
provenance: AUTHORED
source: "lock generated by Microsoft .NET SDK 10.0.400"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````json
{
  "version": 1,
  "dependencies": {
    "net10.0": {
      "Microsoft.Extensions.ObjectPool": {
        "type": "Transitive",
        "resolved": "10.0.0",
        "contentHash": "bpeCq0IYmVLACyEUMzFIOQX+zZUElG1t+nu1lSxthe7B+1oNYking7b91305+jNB6iwojp9fqTY9O+Nh7ULQxg=="
      },
      "System.Security.Cryptography.Pkcs": {
        "type": "Transitive",
        "resolved": "10.0.11",
        "contentHash": "8IV+rI3xN/Mkq9MsSX7VZTv9T9Wt+tyhHLwBX9VNZBk8m8kknEs1JeTSjI2x9M2L/fSja+c2WSE5n5Yy0GXdoQ=="
      },
      "System.Security.Cryptography.Xml": {
        "type": "Transitive",
        "resolved": "10.0.11",
        "contentHash": "TokfVsaU2fmcwvsU6ihLCkseVfwhLX6Tmbl3qh3nOBIqHSX4yqFllpslA+PuHERH/ENu3kwqo4deZiyT0Wvh3Q==",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "10.0.11"
        }
      },
      "System.ServiceModel.Http": {
        "type": "Transitive",
        "resolved": "10.0.652802",
        "contentHash": "G02XZvmccf42QCU5MjviBIg69MSMAVHwL1inVPsNSpfp5g+t5BkQM3DyvWRLN4qmeFDWSF/mA1rIYONIDu/6Dg==",
        "dependencies": {
          "System.ServiceModel.Primitives": "10.0.652802"
        }
      },
      "System.ServiceModel.Primitives": {
        "type": "Transitive",
        "resolved": "10.0.652802",
        "contentHash": "ULfGNl75BNXkpF42wNV2CDXJ64dUZZEa8xO2mBsc4tqbW9QjruxjEB6bAr4Z/T1rNU+leOztIjCJQYsBGFWYlw==",
        "dependencies": {
          "Microsoft.Extensions.ObjectPool": "10.0.0",
          "System.Security.Cryptography.Xml": "10.0.0"
        }
      },
      "elite.arca.credentials": {
        "type": "Project",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "[10.0.11, )"
        }
      },
      "elite.arca.wsfe.bridge": {
        "type": "Project",
        "dependencies": {
          "Elite.Arca.Credentials": "[1.0.0, )",
          "Elite.Arca.Wsfe.Generated": "[1.0.0, )"
        }
      },
      "elite.arca.wsfe.generated": {
        "type": "Project",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "[10.0.11, )",
          "System.Security.Cryptography.Xml": "[10.0.11, )",
          "System.ServiceModel.Http": "[10.0.652802, )"
        }
      },
      "elite.arca.wsfe.worker": {
        "type": "Project",
        "dependencies": {
          "Elite.Arca.Credentials": "[1.0.0, )",
          "Elite.Arca.Wsfe.Bridge": "[1.0.0, )",
          "Elite.Arca.Wsfe.Generated": "[1.0.0, )"
        }
      }
    }
  }
}
````

### FILE: `arca/fiscal/worker/tests/Elite.Arca.Wsfe.Worker.Tests/Program.cs`
```yaml
block_id: "ARCA-WSFE-UDS-WORKER:worker-tests:v1"
operation: CREATE
path: "arca/fiscal/worker/tests/Elite.Arca.Wsfe.Worker.Tests/Program.cs"
sha256: "6ffabf6f70ece8f6ffae8e7c2d97a72e6b5bf7b601bfba03f34f0b5d26020eba"
provenance: AUTHORED
source: "local Kestrel Unix-socket executable contract tests"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````csharp
using System.Net;
using System.Net.Http.Json;
using System.Net.Sockets;
using System.Text;
using Elite.Arca.Wsfe.Bridge;
using Elite.Arca.Wsfe.Worker;

const string InvoiceJson = """
{"taxpayer_cuit":"33693450239","point_of_sale":12,"voucher_type":1,"concept":1,"recipient_document_type":80,"recipient_document":"20111111112","recipient_vat_condition_id":1,"currency":"ARS","voucher_number":41,"total_minor_units":18405,"net_minor_units":15000,"vat_minor_units":2625,"exempt_minor_units":0,"non_taxed_minor_units":0,"other_tax_minor_units":780,"issued_on":"2010-09-03","vat_lines":[{"id":5,"base_minor_units":10000,"amount_minor_units":2100},{"id":4,"base_minor_units":5000,"amount_minor_units":525}],"other_tax_lines":[{"id":99,"description":"Impuesto Municipal Matanza","base_minor_units":15000,"rate_basis_points":520,"amount_minor_units":780}]}
""";

var root = Path.Combine(Path.GetTempPath(), $"elite-wsfe-worker-{Environment.ProcessId}");
Directory.CreateDirectory(root);
var socketPath = Path.Combine(root, "worker.sock");
var operations = new FakeOperations();
await using var app = WorkerHost.Build([], socketPath, operations);
try
{
    await app.StartAsync();
    WorkerHost.SecureSocket(socketPath);
    Assert(Path.Exists(socketPath), "Kestrel created Unix socket");
    Assert(app.Urls.Count > 0 && app.Urls.All(value => value.Contains("unix:", StringComparison.OrdinalIgnoreCase) || value.Contains(socketPath, StringComparison.OrdinalIgnoreCase)), "only the Unix socket address was configured");
    using var handler = new SocketsHttpHandler
    {
        ConnectCallback = async (_, cancellationToken) =>
        {
            var socket = new Socket(AddressFamily.Unix, SocketType.Stream, ProtocolType.Unspecified);
            try
            {
                await socket.ConnectAsync(new UnixDomainSocketEndPoint(socketPath), cancellationToken);
                return new NetworkStream(socket, ownsSocket: true);
            }
            catch { socket.Dispose(); throw; }
        }
    };
    using var client = new HttpClient(handler) { BaseAddress = new Uri("http://arca-wsfe") };

    var health = await client.GetStringAsync("/health/live");
    Assert(health == "{\"status\":\"alive\"}", "liveness response is narrow");

    using var last = await client.PostAsync("/v1/last-authorized", Json(InvoiceJson));
    Assert(last.StatusCode == HttpStatusCode.OK, "last-authorized status");
    Assert((await last.Content.ReadAsStringAsync()).Contains("\"number\":40", StringComparison.Ordinal), "last-authorized result");
    Assert(operations.LastInvoice is { RecipientVatConditionId: 1, VatLines.Count: 2, OtherTaxLines.Count: 1 }, "exact tax details crossed IPC");

    using var authorized = await client.PostAsync("/v1/authorize", Json(InvoiceJson));
    var authorizedBody = await authorized.Content.ReadAsStringAsync();
    Assert(authorized.StatusCode == HttpStatusCode.OK && authorizedBody.Contains("\"cae_expires_on\":\"2010-09-13\"", StringComparison.Ordinal), "authorization result uses strict snake case");

    using var parameters = await client.PostAsync("/v1/parameters", Json("""{"taxpayer_cuit":"33693450239","kind":"vat_rate","voucher_class":null}"""));
    var parameterBody = await parameters.Content.ReadAsStringAsync();
    Assert(parameters.StatusCode == HttpStatusCode.OK && parameterBody.Contains("\"kind\":\"vat_rate\"", StringComparison.Ordinal) && parameterBody.Contains("\"code\":\"5\"", StringComparison.Ordinal), "parameter snapshot crosses UDS");
    Assert(!parameterBody.Contains("secret", StringComparison.OrdinalIgnoreCase) && !parameterBody.Contains("message", StringComparison.OrdinalIgnoreCase), "parameter response is safe");

    using var invalidParameters = await client.PostAsync("/v1/parameters", Json("""{"taxpayer_cuit":"33693450239","kind":"invented","voucher_class":null}"""));
    Assert(invalidParameters.StatusCode == HttpStatusCode.BadRequest, "unknown parameter kind fails closed");

    using var unknown = await client.PostAsync("/v1/authorize", Json(InvoiceJson.Replace("\"taxpayer_cuit\"", "\"unknown\":true,\"taxpayer_cuit\"", StringComparison.Ordinal)));
    Assert(unknown.StatusCode == HttpStatusCode.BadRequest, "unknown request field fails closed");

    operations.FailConsult = true;
    using var failed = await client.PostAsync("/v1/consult", Json(InvoiceJson));
    var failedBody = await failed.Content.ReadAsStringAsync();
    Assert(failed.StatusCode == HttpStatusCode.BadGateway && failedBody.Contains("\"code\":\"provider_rejected\"", StringComparison.Ordinal), "provider failure is classified");
    Assert(!failedBody.Contains("secret-token-sign", StringComparison.Ordinal) && !failedBody.Contains("message", StringComparison.OrdinalIgnoreCase), "provider message and credentials are absent");

    using var oversized = await client.PostAsync("/v1/authorize", Json("{\"padding\":\"" + new string('x', 70_000) + "\"}"));
    Assert(oversized.StatusCode == HttpStatusCode.RequestEntityTooLarge, "oversized request is rejected by Kestrel");
    Console.WriteLine("ARCA_WSFE_WORKER_TESTS_PASS tests=12");
}
finally
{
    await app.StopAsync();
    if (Path.Exists(socketPath)) File.Delete(socketPath);
    Directory.Delete(root);
}

static StringContent Json(string value) => new(value, Encoding.UTF8, "application/json");
static void Assert(bool condition, string name)
{
    if (!condition) throw new InvalidOperationException($"FAILED: {name}");
}

sealed class FakeOperations : IFiscalOperations
{
    internal FiscalInvoice? LastInvoice { get; private set; }
    internal bool FailConsult { get; set; }
    public Task<LastAuthorizedResult> LastAuthorizedAsync(FiscalInvoice invoice, CancellationToken cancellationToken)
    {
        LastInvoice = invoice;
        return Task.FromResult(new LastAuthorizedResult(40, new string('a', 64), []));
    }
    public Task<SafeWsfeResult> ConsultAsync(FiscalInvoice invoice, CancellationToken cancellationToken)
    {
        if (FailConsult) throw new WsfeResponseException("secret-token-sign provider message", new string('b', 64), ["E:100"]);
        return Task.FromResult(new SafeWsfeResult(false, false, null, null, new string('b', 64), ["E:602"]));
    }
    public Task<SafeWsfeResult> AuthorizeAsync(FiscalInvoice invoice, CancellationToken cancellationToken) =>
        Task.FromResult(new SafeWsfeResult(true, true, "41124578989845", new DateOnly(2010, 9, 13), new string('c', 64), []));
    public Task<SafeParameterResult> ParametersAsync(ParameterRequest request, CancellationToken cancellationToken)
    {
        if (request.Kind != "vat_rate" || request.VoucherClass is not null) throw new ArgumentException("unsupported parameter fixture");
        return Task.FromResult(new SafeParameterResult(request.Kind, request.VoucherClass, [new SafeParameterItem("5", "21%", new DateOnly(2009, 2, 20), null, null, null, null, null)], new string('d', 64), []));
    }
}
````

### FILE: `tools/test-arca-wsfe-uds-worker.ps1`
```yaml
block_id: "ARCA-WSFE-UDS-WORKER:runner:v1"
operation: CREATE
path: "tools/test-arca-wsfe-uds-worker.ps1"
sha256: "d4dd61553b86eb038c7bc0fbc84b41003386e28ec52866e8c0c1341d308ede71"
provenance: AUTHORED
source: "local locked runner governed by Microsoft .NET and Go toolchain documentation"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param(
    [Parameter(Mandatory)][string]$BackendRoot,
    [Parameter(Mandatory)][string]$GeneratedClientRoot,
    [Parameter(Mandatory)][string]$DotnetPath,
    [Parameter(Mandatory)][string]$GoPath,
    [Parameter(Mandatory)][string]$WorkDirectory
)
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

function Require-Absolute([string]$PathValue, [string]$Name) {
    if (-not [IO.Path]::IsPathFullyQualified($PathValue)) { throw "$Name must be absolute" }
}
foreach ($entry in @(@($BackendRoot,'BackendRoot'), @($GeneratedClientRoot,'GeneratedClientRoot'), @($DotnetPath,'DotnetPath'), @($GoPath,'GoPath'), @($WorkDirectory,'WorkDirectory'))) {
    Require-Absolute $entry[0] $entry[1]
}
foreach ($file in @($DotnetPath,$GoPath,(Join-Path $BackendRoot 'go.mod'),(Join-Path $GeneratedClientRoot 'Elite.Arca.Wsfe.Generated.csproj'),(Join-Path $GeneratedClientRoot 'packages.lock.json'),(Join-Path $GeneratedClientRoot 'Generated/WsaaReference.cs'),(Join-Path $GeneratedClientRoot 'Generated/WsfeV1Reference.cs'))) {
    if (-not (Test-Path -LiteralPath $file -PathType Leaf)) { throw "Required file missing: $file" }
}
if (-not (Test-Path -LiteralPath $WorkDirectory)) { New-Item -ItemType Directory -Path $WorkDirectory | Out-Null }
if ((Get-ChildItem -LiteralPath $WorkDirectory -Force | Measure-Object).Count -ne 0) { throw 'WorkDirectory must be empty' }

$dotnetVersion = (& $DotnetPath --version).Trim()
$goVersion = (& $GoPath version).Trim()
if ($dotnetVersion -ne '10.0.400') { throw "Unexpected .NET SDK: $dotnetVersion" }
if ($goVersion -notmatch '\bgo1\.26\.7\b') { throw "Unexpected Go toolchain: $goVersion" }

Copy-Item -Path (Join-Path $BackendRoot '*') -Destination $WorkDirectory -Recurse -Force
$generatedDestination = Join-Path $WorkDirectory 'arca/fiscal/generated'
New-Item -ItemType Directory -Path (Join-Path $generatedDestination 'Generated') -Force | Out-Null
Copy-Item -LiteralPath (Join-Path $GeneratedClientRoot 'Elite.Arca.Wsfe.Generated.csproj') -Destination $generatedDestination
Copy-Item -LiteralPath (Join-Path $GeneratedClientRoot 'packages.lock.json') -Destination $generatedDestination
Copy-Item -LiteralPath (Join-Path $GeneratedClientRoot 'Generated/WsaaReference.cs') -Destination (Join-Path $generatedDestination 'Generated')
Copy-Item -LiteralPath (Join-Path $GeneratedClientRoot 'Generated/WsfeV1Reference.cs') -Destination (Join-Path $generatedDestination 'Generated')

$previousToolchain = $env:GOTOOLCHAIN
$env:GOTOOLCHAIN = 'local'
Push-Location $WorkDirectory
try {
    & $GoPath test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'Go tests failed' }
    & $GoPath vet ./...
    if ($LASTEXITCODE -ne 0) { throw 'Go vet failed' }
    & $GoPath build ./cmd/electromobility-api ./cmd/arca-fiscal-worker
    if ($LASTEXITCODE -ne 0) { throw 'Go build failed' }
    $testProject = Join-Path $WorkDirectory 'arca/fiscal/worker/tests/Elite.Arca.Wsfe.Worker.Tests/Elite.Arca.Wsfe.Worker.Tests.csproj'
    & $DotnetPath restore $testProject --locked-mode
    if ($LASTEXITCODE -ne 0) { throw '.NET locked restore failed' }
    & $DotnetPath build $testProject -c Release --no-restore
    if ($LASTEXITCODE -ne 0) { throw '.NET build failed' }
    & $DotnetPath run --project $testProject -c Release --no-build
    if ($LASTEXITCODE -ne 0) { throw '.NET worker contract tests failed' }
}
finally {
    Pop-Location
    $env:GOTOOLCHAIN = $previousToolchain
}
Write-Output 'ARCA_WSFE_UDS_WORKER_PASS go=1.26.7 dotnet=10.0.400 uds_tests=12 parameter_bridge=1 production_admitted=false'
````

## 6. Configuration surface

| Input | Type/default | Secret | Validation and effect |
|---|---|---:|---|
| `DATABASE_URL` | PostgreSQL DSN / none | yes | required by Go worker; parsed before connection and never logged |
| `ARCA_WSFE_SOCKET` | absolute path / none | no | required by both processes; parent exists, is not a reparse point and target must not exist at .NET startup |
| `FISCAL_WORKER_ID` | 1–64 safe characters / none | no | required stable operational identity; no customer data |
| `ARCA_ENVIRONMENT` | literal `homologation` / none | no | any other value fails before socket creation |
| `ARCA_CERTIFICATE_THUMBPRINT` | SHA-1/SHA-256 hex / none | no | resolves to exactly one valid RSA >=2048 certificate with private key |
| `ARCA_CERTIFICATE_STORE` | `CurrentUser|LocalMachine` / none | no | any other store fails startup |
| certificate/private key | OS certificate store | yes | external only; never accepted through JSON/env as PFX bytes/password |
| WSAA/WSFE endpoints | fixed homologation URLs | no | not configurable in this admission; prevents silent production switch |

## 7. Dependency bill

| Package/tool | Exact pin | Use | License | Runtime/build | Official source |
|---|---:|---|---|---|---|
| Go | 1.26.7 | UDS provider and durable worker | BSD-3-Clause | runtime/build | go.dev |
| .NET SDK/ASP.NET Core | 10.0.400 / 10.0.11 runtime | Kestrel UDS worker | MIT/.NET terms | runtime/build | Microsoft |
| System.ServiceModel.Http | 10.0.652802 | generated SOAP clients | MIT | runtime | Microsoft NuGet |
| System.Security.Cryptography.Pkcs/Xml | 10.0.11 | WSAA signing/generated graph | MIT | runtime | Microsoft NuGet |
| pgx | 5.10.0 | durable PostgreSQL worker | MIT | runtime | jackc/pgx |
| PostgreSQL | 18.6 | fiscal source of truth | PostgreSQL | runtime | postgresql.org |
| ARCA WSAA/WSFEv1 | manual WSFE 4.6 + live generated homologation contracts | protocol authority | public specification | integration | ARCA |

## 8. Apply order

Materialize Go foundation/domain/PostgreSQL, fiscal API 0.4.x and application 1.4.x. Generate the pinned homologation clients, then materialize credential core, SOAP adapter and this pack. Provision a least-privilege socket directory and certificate-store identity outside source control. Execute the locked runner in an empty directory; configure API and worker as separate processes. Rollback stops both processes and preserves invoices, parameter snapshots and decisions; never delete fiscal evidence or renumber invoices.

## 9. Verification

Require 13/13 hash-exact materialization; Go full test, focal UDS contract, vet/builds; locked .NET restore; warning-free worker/test build; twelve Kestrel UDS contracts covering actual socket, no TCP, liveness, exact tax shape, parameter snapshot, unknown kind 400, unknown field, sanitized error and body limit. Re-run NuGet audit. Live homologation, credential association, refresh/policy approval and tax/accounting acceptance remain project gates; production remains false.

## 10. Reconstruction evidence

V130 records issuance IPC; V132 records parameter UDS, persistence and separate approval. V145 records associated-voucher transport in `reconstruction_evidence/ARCA_ASSOCIATED_VOUCHER_EXECUTION_INVENTORY_2026-08-31_V145.md`.

V402 composed delta: V402325 connected local ARCA fixture, fixed offline generated-tool adaptation, explicit legal/source pins and portable .NET/Go build. See ARCA_CONNECTED_INFRA_V402 evidence; only future homologation credentials conditioned, no productive fiscal claim.

### FILE: `arca/fiscal/worker/fixtures/Elite.Arca.Wsfe.Fixture/Elite.Arca.Wsfe.Fixture.csproj`

```yaml
block_id: "ARCA_WSFE_UDS_WORKER-V402325:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f84ce0b3b0ae4dccfba4dc0bf3bed475fe888b7e6ccc917ace3b01d7b58fca33"
variables: []
secrets_allowed: false
```

````text
<Project Sdk="Microsoft.NET.Sdk.Web">
  <PropertyGroup>
    <OutputType>Exe</OutputType><TargetFramework>net10.0</TargetFramework>
    <ImplicitUsings>enable</ImplicitUsings><Nullable>enable</Nullable>
    <TreatWarningsAsErrors>true</TreatWarningsAsErrors><RestorePackagesWithLockFile>true</RestorePackagesWithLockFile>
    <NuGetAudit>true</NuGetAudit><NuGetAuditMode>all</NuGetAuditMode>
  </PropertyGroup>
  <ItemGroup><ProjectReference Include="../../Elite.Arca.Wsfe.Worker/Elite.Arca.Wsfe.Worker.csproj" /></ItemGroup>
</Project>
````

### FILE: `arca/fiscal/worker/fixtures/Elite.Arca.Wsfe.Fixture/FixtureParameters.cs`

```yaml
block_id: "ARCA_WSFE_UDS_WORKER-V402325:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6a144ed9a2b53b65d6d3cfa21c13234357616db766378484c1dc823006fcd002"
variables: []
secrets_allowed: false
```

````text
using Elite.Arca.Wsfe.Bridge;
using Generated = Elite.Arca.Wsfe.Generated;

sealed class FakeParameters : IWsfeParameterSoap
{
    internal Generated.FEAuthRequest? LastAuth { get; private set; }
    internal bool Error { get; set; }
    internal bool Duplicate { get; set; }
    internal bool MismatchedClass { get; set; }
    private Generated.Err[]? Errors => Error ? [new Generated.Err { Code = 500, Msg = "provider message secret" }] : null;
    private static Generated.Evt[] Events => [new Generated.Evt { Code = 10, Msg = "provider message" }];
    private T Capture<T>(Generated.FEAuthRequest auth, T response) { LastAuth = auth; return response; }

    public Task<Generated.CbteTipoResponse> VoucherTypesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken) => Task.FromResult(Capture(auth, new Generated.CbteTipoResponse { ResultGet = [new Generated.CbteTipo { Id = 1, Desc = "Factura A", FchDesde = "20100101", FchHasta = "" }], Errors = Errors, Events = Events }));
    public Task<Generated.ConceptoTipoResponse> ConceptsAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken) => Task.FromResult(Capture(auth, new Generated.ConceptoTipoResponse { ResultGet = [new Generated.ConceptoTipo { Id = 1, Desc = "Productos", FchDesde = "20100101", FchHasta = "" }], Errors = Errors }));
    public Task<Generated.DocTipoResponse> DocumentTypesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken)
    {
        var items = Duplicate ? new[] { Doc(), Doc() } : [Doc()];
        return Task.FromResult(Capture(auth, new Generated.DocTipoResponse { ResultGet = items, Errors = Errors }));
    }
    public Task<Generated.IvaTipoResponse> VatRatesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken) => Task.FromResult(Capture(auth, new Generated.IvaTipoResponse { ResultGet = [new Generated.IvaTipo { Id = "5", Desc = "21%", FchDesde = "20100101", FchHasta = "" }], Errors = Errors }));
    public Task<Generated.FETributoResponse> OtherTaxesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken) => Task.FromResult(Capture(auth, new Generated.FETributoResponse { ResultGet = [new Generated.TributoTipo { Id = 99, Desc = "Otros", FchDesde = "20100101", FchHasta = "" }], Errors = Errors }));
    public Task<Generated.FEPtoVentaResponse> PointsOfSaleAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken) => Task.FromResult(Capture(auth, new Generated.FEPtoVentaResponse { ResultGet = [new Generated.PtoVenta { Nro = 12, EmisionTipo = "CAE", Bloqueado = "N", FchBaja = "" }], Errors = Errors }));
    public Task<Generated.CondicionIvaReceptorResponse> RecipientVatConditionsAsync(Generated.FEAuthRequest auth, string voucherClass, CancellationToken cancellationToken) => Task.FromResult(Capture(auth, new Generated.CondicionIvaReceptorResponse { ResultGet = [new Generated.CondicionIvaReceptor { Id = 1, Desc = "IVA Responsable Inscripto", Cmp_Clase = MismatchedClass ? "B" : voucherClass }], Errors = Errors }));
    private static Generated.DocTipo Doc() => new() { Id = 80, Desc = "CUIT", FchDesde = "20100101", FchHasta = "" };
}
````

### FILE: `arca/fiscal/worker/fixtures/Elite.Arca.Wsfe.Fixture/Program.cs`

```yaml
block_id: "ARCA_WSFE_UDS_WORKER-V402325:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "45be2119f23ab658117fed445a890e6e0608a73c456a1782eb801b126adf30a3"
variables: []
secrets_allowed: false
```

````text
// AUTHORED local fixture composition. No live transport or user certificate.
using System.Security.Cryptography;
using System.Security.Cryptography.Pkcs;
using System.Security.Cryptography.X509Certificates;
using System.Xml;
using Elite.Arca.Credentials;
using Elite.Arca.Wsfe.Bridge;
using Elite.Arca.Wsfe.Worker;
using Generated = Elite.Arca.Wsfe.Generated;

if (args.Length != 1 || Environment.GetEnvironmentVariable("ELITE_ARCA_FIXTURE") != "1")
    throw new InvalidOperationException("Explicit local fixture mode and one Unix socket path required.");
using var rsa = RSA.Create(2048);
var request = new CertificateRequest("CN=Elite Local ARCA Fixture Only", rsa, HashAlgorithmName.SHA256, RSASignaturePadding.Pkcs1);
using var certificate = request.CreateSelfSigned(DateTimeOffset.UtcNow.AddMinutes(-1), DateTimeOffset.UtcNow.AddHours(2));
var wsaa = new FixtureWsaa();
using var credentials = new WsaaCredentialProvider(WsaaOptions.ElectronicInvoice, new SystemUtcClock(),
    new LoginTicketRequestFactory(WsaaOptions.ElectronicInvoice, new SystemUtcClock(), new MonotonicLoginTicketIdSource()),
    new CmsRequestSigner(), certificate, wsaa);
var access = new CredentialAccessSource(credentials);
var soap = new FixtureSoap();
var operations = new BridgeFiscalOperations(new WsfeBridge(access, soap), new WsfeParameterBridge(access, new FakeParameters()));
await using var app = WorkerHost.Build([], args[0], operations);
app.MapGet("/fixture/stats", () => new { scope = "LOCAL_FIXTURES", last = soap.LastCalls, authorize = soap.AuthorizeCalls, consult = soap.ConsultCalls, cms = wsaa.Calls });
app.MapPost("/fixture/stop", (IHostApplicationLifetime lifetime) => { lifetime.StopApplication(); return Results.Ok(new { stopping = true }); });
await app.StartAsync();
WorkerHost.SecureSocket(args[0]);
if (app.Urls.Count == 0 || app.Urls.Any(value => !value.Contains("unix:", StringComparison.OrdinalIgnoreCase) && !value.Contains(args[0], StringComparison.OrdinalIgnoreCase)))
    throw new InvalidOperationException("Fixture must expose only the Unix socket.");
Console.WriteLine("ARCA_CONNECTED_FIXTURE_READY uds_only=true live_transport=false");
await app.WaitForShutdownAsync();

sealed class FixtureWsaa : IWsaaTransport
{
    public int Calls;
    public Task<string> LoginCmsAsync(string cmsBase64, CancellationToken cancellationToken)
    {
        cancellationToken.ThrowIfCancellationRequested();
        var cms = new SignedCms(); cms.Decode(Convert.FromBase64String(cmsBase64)); cms.CheckSignature(verifySignatureOnly: true);
        if (cms.ContentInfo.Content.Length == 0) throw new InvalidDataException("Empty signed fixture request.");
        Interlocked.Increment(ref Calls);
        var expiration = XmlConvert.ToString(DateTimeOffset.UtcNow.AddHours(1));
        return Task.FromResult($"<loginTicketResponse><header><expirationTime>{expiration}</expirationTime></header><credentials><token>public-fixture-token</token><sign>public-fixture-sign</sign></credentials></loginTicketResponse>");
    }
}

sealed class FixtureSoap : IWsfeSoap
{
    private readonly object sync = new();
    private readonly Dictionary<(long Cuit, int Point, int Type, long Number), Generated.FECompConsResponse> issued = [];
    public int LastCalls, AuthorizeCalls, ConsultCalls;
    private static void Auth(Generated.FEAuthRequest auth)
    {
        if (auth.Token != "public-fixture-token" || auth.Sign != "public-fixture-sign") throw new InvalidDataException("Fixture credential mapping failed.");
    }
    public Task<Generated.FERecuperaLastCbteResponse> LastAuthorizedAsync(Generated.FEAuthRequest auth, int pointOfSale, int voucherType, CancellationToken cancellationToken)
    {
        Auth(auth); Interlocked.Increment(ref LastCalls);
        lock (sync) return Task.FromResult(new Generated.FERecuperaLastCbteResponse { PtoVta = pointOfSale, CbteTipo = voucherType, CbteNro = checked((int)issued.Keys.Where(k => k.Cuit == auth.Cuit && k.Point == pointOfSale && k.Type == voucherType).Select(k => k.Number).DefaultIfEmpty(0).Max()) });
    }
    public Task<Generated.FECompConsultaResponse> ConsultAsync(Generated.FEAuthRequest auth, Generated.FECompConsultaReq request, CancellationToken cancellationToken)
    {
        Auth(auth); Interlocked.Increment(ref ConsultCalls);
        lock (sync) return Task.FromResult(issued.TryGetValue((auth.Cuit, request.PtoVta, request.CbteTipo, request.CbteNro), out var result)
            ? new Generated.FECompConsultaResponse { ResultGet = result }
            : new Generated.FECompConsultaResponse { Errors = [new Generated.Err { Code = 602, Msg = "fixture absent" }] });
    }
    public Task<Generated.FECAEResponse> AuthorizeAsync(Generated.FEAuthRequest auth, Generated.FECAERequest request, CancellationToken cancellationToken)
    {
        Auth(auth); var count = Interlocked.Increment(ref AuthorizeCalls);
        var detail = request.FeDetReq.Single(); var header = request.FeCabReq;
        if (detail.ImpTotal != 121 || detail.ImpNeto != 100 || detail.ImpIVA != 21 || !detail.CondicionIVAReceptorIdSpecified || detail.MonId != "PES")
            throw new InvalidDataException("Exact fixture invoice mapping failed.");
        lock (sync)
        {
            if (header.CbteTipo == 8 && (detail.CbtesAsoc is not { Length: 1 } || !issued.ContainsKey((auth.Cuit, detail.CbtesAsoc[0].PtoVta, detail.CbtesAsoc[0].Tipo, detail.CbtesAsoc[0].Nro))))
                throw new InvalidDataException("Credit fixture must reference issued original.");
            issued.Add((auth.Cuit, header.PtoVta, header.CbteTipo, detail.CbteDesde), new Generated.FECompConsResponse { PtoVta = header.PtoVta, CbteTipo = header.CbteTipo, CbteDesde = detail.CbteDesde, CbteHasta = detail.CbteHasta, Resultado = "A", CodAutorizacion = "41124578989845", FchVto = "20300913" });
        }
        // Provider commits before losing its first response. Retry must consult.
        if (count == 1) throw new IOException("Deliberate fixture response loss after authorization.");
        return Task.FromResult(new Generated.FECAEResponse { FeCabResp = new Generated.FECAECabResponse { PtoVta = header.PtoVta, CbteTipo = header.CbteTipo, CantReg = 1, Resultado = "A" }, FeDetResp = [new Generated.FECAEDetResponse { CbteDesde = detail.CbteDesde, CbteHasta = detail.CbteHasta, Resultado = "A", CAE = "41124578989845", CAEFchVto = "20300913" }] });
    }
}
````

### FILE: `arca/fiscal/worker/fixtures/Elite.Arca.Wsfe.Fixture/packages.lock.json`

```yaml
block_id: "ARCA_WSFE_UDS_WORKER-V402325:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "491bd6344247876ff0cc67e35da754618abb4ea48128d3df248e2e8c4488be90"
variables: []
secrets_allowed: false
```

````json
{
  "version": 1,
  "dependencies": {
    "net10.0": {
      "System.Security.Cryptography.Pkcs": {
        "type": "Transitive",
        "resolved": "10.0.11",
        "contentHash": "8IV+rI3xN/Mkq9MsSX7VZTv9T9Wt+tyhHLwBX9VNZBk8m8kknEs1JeTSjI2x9M2L/fSja+c2WSE5n5Yy0GXdoQ=="
      },
      "System.ServiceModel.Http": {
        "type": "Transitive",
        "resolved": "10.0.652802",
        "contentHash": "G02XZvmccf42QCU5MjviBIg69MSMAVHwL1inVPsNSpfp5g+t5BkQM3DyvWRLN4qmeFDWSF/mA1rIYONIDu/6Dg==",
        "dependencies": {
          "System.ServiceModel.Primitives": "10.0.652802"
        }
      },
      "System.ServiceModel.Primitives": {
        "type": "Transitive",
        "resolved": "10.0.652802",
        "contentHash": "ULfGNl75BNXkpF42wNV2CDXJ64dUZZEa8xO2mBsc4tqbW9QjruxjEB6bAr4Z/T1rNU+leOztIjCJQYsBGFWYlw=="
      },
      "elite.arca.credentials": {
        "type": "Project",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "[10.0.11, )"
        }
      },
      "elite.arca.wsfe.bridge": {
        "type": "Project",
        "dependencies": {
          "Elite.Arca.Credentials": "[1.0.0, )",
          "Elite.Arca.Wsfe.Generated": "[1.0.0, )"
        }
      },
      "elite.arca.wsfe.generated": {
        "type": "Project",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "[10.0.11, )",
          "System.ServiceModel.Http": "[10.0.652802, )"
        }
      },
      "elite.arca.wsfe.worker": {
        "type": "Project",
        "dependencies": {
          "Elite.Arca.Credentials": "[1.0.0, )",
          "Elite.Arca.Wsfe.Bridge": "[1.0.0, )",
          "Elite.Arca.Wsfe.Generated": "[1.0.0, )"
        }
      }
    }
  }
}
````

### FILE: `docs/ARCA_LOCAL_REFERENCE.md`

```yaml
block_id: "ARCA_WSFE_UDS_WORKER-V402325:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f4f919cf51a2d532756a4427281ee118c75a4543a59a31db2efe65c2fdfdc50a"
variables: []
secrets_allowed: false
```

````markdown
# ARCA: referencia local y conexión de homologación

El artefacto contiene los dos workers Go, el worker SOAP/UDS .NET, un host de
fixtures explícito y el runtime Microsoft .NET/ASP.NET 10.0.11 fijado. Los clientes
WSAA/WSFE se generan desde los dos WSDL oficiales fijados; sus fuentes y receipt
quedan en `arca/generated-source`. No contiene certificados ni secretos de usuario.

## Ensayo local

Usar únicamente una base PostgreSQL de fixtures con las migraciones del mismo
artefacto. El API debe apuntar a esa misma base. No conectar el fixture fiscal a
datos reales. La identidad fiscal y las operaciones de punto de venta/factura
son las del contrato HTTP materializado y mantienen autorización por organización.

En tres terminales PowerShell abiertas en la raíz del artefacto:

```powershell
# Terminal 1: host de homologación simulada; genera su certificado efímero local.
$env:ELITE_ARCA_FIXTURE = '1'
$env:ARCA_WSFE_SOCKET = Join-Path $env:TEMP 'elite-arca-reference.sock'
& .\arca\dotnet\dotnet.exe .\arca\fixture\Elite.Arca.Wsfe.Fixture.dll $env:ARCA_WSFE_SOCKET
```

```powershell
# Terminal 2: configurar DATABASE_URL para la base local de fixtures.
$env:ARCA_WSFE_SOCKET = Join-Path $env:TEMP 'elite-arca-reference.sock'
$env:FISCAL_WORKER_ID = 'reference-fiscal-1'
& .\arca\go\arca-fiscal-worker.exe
```

```powershell
# Terminal 3: la misma DATABASE_URL y el mismo socket.
$env:ARCA_WSFE_SOCKET = Join-Path $env:TEMP 'elite-arca-reference.sock'
$env:FISCAL_PARAMETER_WORKER_ID = 'reference-parameters-1'
& .\arca\go\arca-parameter-worker.exe
```

El fixture reproduce una respuesta perdida después de autorizar; el worker Go
consulta el comprobante y reconcilia su estado durable sin volver a emitirlo.
También acepta la nota de crédito asociada y las siete consultas de parámetros.
Sus valores son datos sintéticos para el ensayo, no reglas fiscales aplicables.
Los procesos interactivos se detienen con Ctrl+C; la aceptación automatizada
cierra el fixture por su endpoint de control privado UDS y verifica PostgreSQL
después de un reinicio.

## Cuando el usuario aporte sus credenciales

La conexión de homologación usa el mismo socket y los mismos workers Go. Reemplazar
el host de fixtures por `arca/worker/Elite.Arca.Wsfe.Worker.dll` ejecutado con el
runtime incluido. Sus entradas son `ARCA_ENVIRONMENT=homologation`,
`ARCA_WSFE_SOCKET`, `ARCA_CERTIFICATE_THUMBPRINT` y
`ARCA_CERTIFICATE_STORE=CurrentUser|LocalMachine`; el certificado con clave privada
debe existir en el almacén indicado. El CUIT/punto de venta pertenece a la
configuración de la organización en PostgreSQL, nunca al código generado.

Estado de la conexión externa: `CONDITIONED_USER_CREDENTIALS`. El worker falla
cerrado si faltan sus entradas. Este artefacto no habilita facturación productiva;
las reglas fiscales, autorización del contribuyente y aceptación real se aplican
cuando se materialice el proyecto concreto.

## Reconstrucción

`ci/build_signed_reference.ps1` llama al builder local, que incluye
`ci/build_arca_reference.py`. Los inputs externos del build se fijan por SHA-256:
SDK/runtime, PowerShell, archivo original svcutil, nueve parches NuGet oficiales,
los dos WSDL y cinco paquetes de runtime. Cada build crea su propia caché NuGet
vacía y restaura sólo desde el directorio local comprobado. No ejecutar un
`dotnet tool restore` sobre el manifiesto histórico svcutil8.0.0: el generador
admitido es la composición adaptada definida en `svcutil-adaptation.lock.json`.
La adaptación es glue local declarado; no es una publicación oficial de Microsoft.

El release conserva licencias originales, fuentes correspondientes y el grafo
de dependencias. La generación offline no afirma haber consultado vulnerabilidades:
el gate de release ejecuta SCA sobre los manifiestos explícitos de esta revisión.
````

### FILE: `internal/fiscal/wsfeipc/serialization_test.go`

```yaml
block_id: "ARCA_WSFE_UDS_WORKER-V402325:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e82673059e72910fd8b21e43cc4c577c43a34adc585b7b0c656eb5db42c4cc5f"
variables: []
secrets_allowed: false
```

````go
package wsfeipc

import (
 "encoding/json"
 "testing"
)

func TestEmptyFiscalWireCollectionsRemainArrays(t *testing.T) {
 invoice := fixtureInvoice()
 invoice.VATLines = nil
 invoice.OtherTaxLines = nil
 invoice.AssociatedVouchers = nil
 raw,err := json.Marshal(request(invoice));if err!=nil{t.Fatal(err)}
 var payload map[string]json.RawMessage;if err=json.Unmarshal(raw,&payload);err!=nil{t.Fatal(err)}
 for _,field:=range []string{"vat_lines","other_tax_lines","associated_vouchers"}{if string(payload[field])!="[]"{t.Fatalf("%s=%s: .NET collection contract requires an array",field,payload[field])}}
}
````

