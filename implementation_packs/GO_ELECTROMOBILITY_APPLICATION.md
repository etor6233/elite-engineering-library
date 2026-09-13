# Go Electromobility Application Composition

V314: ambos hosts esperan el drenaje HTTP antes de retornar. Shutdown bloquea
esta ruta hasta completar requests o deadline15s; al vencer cierra conexiones
y devuelve el fallo. Listener/startup errors conservan causa. Cuatro regresiones
HTTP TCP y una de ciclo con seis negativos; evidencia HTTP_HOST_SHUTDOWN_V314.md.
No acredita señales OS, target, WebSockets/hijacked ni handlers que ignoren
cancelación. Sin dependencia ni regla de negocio nueva; no promoción integral.

## 1. Metadata

```yaml
pack_id: "GO-ELECTROMOBILITY-APPLICATION"
pack_version: "1.22.3"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el composition root ejecutable, el worker durable de refresco fiscal y un único generador CSPRNG compartido para órdenes, operaciones, inventario serial/bulk/warehouse y transferencias cross-organization con tránsito/costo exacto, red, regalías, ledger, fiscal, journey, consultas y webhooks."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "OIDC"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "GO-ELECTROMOBILITY-PUBLIC-CRM-API 0.2.x", "GO-SUPPLY-FACTORY-INVENTORY-API 0.12.x", "GO-COMMERCE-PRICING-PAYMENT-API 0.3.x or 0.4.x", "GO-FULFILLMENT-SERVICE-FRANCHISE-API 0.3.x", "GO-ENTERPRISE-QUERY-API 0.1.x", "GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.8.x", "GO-FRANCHISE-ROYALTY-SETTLEMENT-API 0.1.x", "GO-ENTERPRISE-ACCOUNTING-LEDGER-API 0.1.x", "GO-ARCA-FISCAL-ISSUANCE-API 0.5.x", "GO-PROVIDER-INTEGRATION-CORE 0.1.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0 AND MIT dependencies"
upstream_sources: ["https://go.dev", "https://www.postgresql.org/docs/18/", "https://github.com/coreos/go-oidc"]
verified_at: "2026-09-13"
```

## 2. Applicability

V298 cablea PAYMENT_REQUEST_PROVIDER: ausente/vacío desactiva las escrituras HTTP
de pago; stripe o mercadopago habilitan sólo el registro local de intención con
el provider seleccionado. Otros valores fallan sin imprimirlos. No solicita
claves ni ejecuta SDK de cobro. Mantiene ARCA_ENABLED diferida y sus tests.
Nueva configuración AUTHORED, no habilitación productiva. Evidencia V298.

Use as the executable composition root only after the selected Go domain, PostgreSQL, query, worker and provider-edge packs are present. Reject it for a partial composition whose interfaces or migrations are missing. It wires capabilities; it does not select an identity provider, cloud, legal policy or external provider implementation.

## 3. Architecture contract

One process owns HTTP composition while PostgreSQL remains the durable source of truth. Startup validates configuration, OIDC discovery, provider secret references and database connectivity before serving. Public, staff, factory and customer boundaries use distinct route policies. Durable provider ingestion authenticates exact bytes before tenant resolution and transactionally enqueues work. Liveness is not readiness; termination must drain traffic/workers before connection closure.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/http_metrics_profile.go
CREATE cmd/electromobility-api/native_stop_other.go
CREATE cmd/electromobility-api/native_stop_windows.go
CREATE cmd/electromobility-api/native_stop_windows_test.go
CREATE cmd/electromobility-api/portal_activation.go
CREATE cmd/electromobility-api/main.go
CREATE cmd/electromobility-api/main_test.go
CREATE cmd/arca-parameter-worker/main.go
CREATE internal/platform/randomid/generator.go
CREATE internal/platform/randomid/generator_test.go
CREATE docs/electromobility-api-runtime.md
CREATE cmd/electromobility-api/whatsapp_activation.go
CREATE cmd/electromobility-api/whatsapp_activation_test.go
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/main_test.go`

```yaml
block_id: "GO-EM-APP:fiscal-selection-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "a72612f92877918f7284bee95c0abf1fddf9b50779f6d521df6f06eaa4eef5af"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	"elite.local/enterprise/internal/providerintegration"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

type rejectIdentity struct{}

func TestPaymentRequestSelection(t *testing.T) {
	for _, v := range []string{"", "stripe", "mercadopago", "unknown", "true", "STRIPE", " stripe"} {
		got, err := paymentRequestProvider(func(string) string { return v })
		valid := v == "" || v == "stripe" || v == "mercadopago"
		if (err == nil) != valid || (valid && got != v) {
			t.Fatalf("selection %q result %q err%v", v, got, err)
		}
	}
}

func (rejectIdentity) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{}, identity.ErrUnauthenticated
}

func TestOptionalProviderConfigurationFailsClosedWithoutLeaking(t *testing.T) {
	entry := `{"tenant_id":"tenant","connection_id":"primary","provider_code":"sandbox","secret_env":"TEST_SECRET"}`
	for _, tc := range []struct {
		name, config, secret string
		valid, present       bool
	}{
		{"absent", "", "", true, false},
		{"empty-registry", "[]", "", true, false},
		{"malformed", "{", "", false, false},
		{"inline-secret", `[{"secret":"DO_NOT_LEAK_VALUE"}]`, "", false, false},
		{"missing-reference", `[{}]`, "", false, false},
		{"missing-secret", "[" + entry + "]", "", false, false},
		{"short-secret", "[" + entry + "]", "DO_NOT_LEAK_VALUE", false, false},
		{"valid", "[" + entry + "]", strings.Repeat("s", 32), true, true},
		{"duplicate", "[" + entry + "," + entry + "]", strings.Repeat("s", 32), false, false},
		{"trailing", "[" + entry + "]{}", strings.Repeat("s", 32), false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			registry, err := providerintegration.LoadHMACConnections([]byte(tc.config), func(string) (string, bool) { calls++; return tc.secret, tc.secret != "" })
			if (err == nil) != tc.valid {
				t.Fatal("unexpected configuration outcome")
			}
			if err != nil && (strings.Contains(err.Error(), "DO_NOT_LEAK_VALUE") || (tc.secret != "" && strings.Contains(err.Error(), tc.secret))) {
				t.Fatal("secret exposed in error")
			}
			if err == nil {
				_, found := registry.Lookup("sandbox", "primary")
				if found != tc.present {
					t.Fatal("unexpected provider activated")
				}
			}
			if (tc.config == "" || tc.config == "[]") && calls != 0 {
				t.Fatal("inactive integration requested secrets")
			}
		})
	}
}

func TestFiscalActivationRequiresExplicitValidConfiguration(t *testing.T) {
	for _, tc := range []struct {
		name, flag, socket    string
		wantModule, wantError bool
	}{
		{"absent", "", "", false, false},
		{"deferred", "false", "", false, false},
		{"deferred-ignores-unused-socket", "false", "must-not-connect", false, false},
		{"invalid", "yes", "", false, true},
		{"whitespace", " true ", "", false, true},
		{"enabled-missing-socket", "true", "", false, true},
		{"enabled-relative-socket", "true", "relative.sock", false, true},
		{"enabled-valid-socket", "true", filepath.Join(t.TempDir(), "wsfe.sock"), true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			socketLookups := 0
			m, err := selectedFiscalModule(nil, randomid.Generator{}, func(k string) string {
				if k == "ARCA_ENABLED" {
					return tc.flag
				}
				socketLookups++
				return tc.socket
			})
			if (err != nil) != tc.wantError || (m != nil) != tc.wantModule {
				t.Fatalf("unexpected module/error presence")
			}
			if !tc.wantModule && !tc.wantError && socketLookups != 0 {
				t.Fatal("deferred path read fiscal configuration")
			}
			mux := http.NewServeMux()
			mux.HandleFunc("GET /independent", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
			if m != nil {
				m.Register(mux, rejectIdentity{})
			}
			independent := httptest.NewRecorder()
			mux.ServeHTTP(independent, httptest.NewRequest("GET", "/independent", nil))
			if independent.Code != 204 {
				t.Fatal("independent route affected")
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, httptest.NewRequest("POST", "/v1/fiscal/invoices", nil))
			wantStatus := http.StatusNotFound
			if tc.wantModule {
				wantStatus = http.StatusUnauthorized
			}
			if w.Code != wantStatus {
				t.Fatalf("fiscal route status = %d, want %d", w.Code, wantStatus)
			}
		})
	}
}
````

### FILE: `cmd/electromobility-api/main.go`

```yaml
block_id: "GO-EM-APP:main:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "057a3f33cdb105e67c8b3026467c432cbb7c9f6afafc16dd54137cdfa9153e6c"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/electromobility"
	"elite.local/enterprise/internal/enterprisequery"
	"elite.local/enterprise/internal/fiscal"
	"elite.local/enterprise/internal/fiscal/wsfeipc"
	"elite.local/enterprise/internal/fulfillment"
	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/order"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"elite.local/enterprise/internal/providerintegration"
	"elite.local/enterprise/internal/royalty"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, nativeCleanup, err := nativeStopContext(ctx, os.Getenv("ELITE_STOP_EVENT_HANDLE"))
	if err != nil {
		slog.Error("native shutdown binding is invalid")
		os.Exit(2)
	}
	defer func() {
		if err := nativeCleanup(); err != nil {
			slog.Error("native shutdown cleanup failed")
		}
	}()
	if ctx.Err() != nil {
		return
	}
	businessPolicy, err := selectedBusinessPolicy(os.Getenv)
	if err != nil {
		slog.Error("business policy profile configuration is invalid")
		os.Exit(2)
	}
	databaseURL := os.Getenv("DATABASE_URL")
	issuer := os.Getenv("OIDC_ISSUER")
	audience := os.Getenv("OIDC_AUDIENCE")
	if databaseURL == "" || issuer == "" || audience == "" {
		slog.Error("DATABASE_URL, OIDC_ISSUER and OIDC_AUDIENCE are required")
		os.Exit(2)
	}
	verifier, err := identity.NewOIDCVerifier(ctx, issuer, audience)
	if err != nil {
		slog.Error("OIDC discovery failed")
		os.Exit(1)
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		slog.Error("database configuration failed")
		os.Exit(2)
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		slog.Error("database unavailable")
		os.Exit(1)
	}
	generator := randomid.Generator{}
	orders := order.NewService(postgres.NewOrders(pool), generator)
	catalog := electromobility.NewService(postgres.NewElectromobility(pool), generator)
	operationsService := operations.NewService(postgres.NewOperations(pool), generator)
	inventoryRepository := postgres.NewInventoryControl(pool)
	inventoryControlService := inventorycontrol.NewService(inventoryRepository, generator)
	bulkInventoryService := inventorycontrol.NewBulkService(inventoryRepository, generator)
	warehouseService := inventorycontrol.NewWarehouseService(inventoryRepository, generator)
	bulkTransferService := inventorycontrol.NewBulkTransferService(inventoryRepository, generator)
	policyRuntime, err := prepareBusinessPolicyRuntime(pool, generator, systemClock{}, businessPolicy)
	if err != nil {
		slog.Error("business policy owner binding is invalid")
		os.Exit(2)
	}
	commerceService := policyRuntime.commerceService
	paymentHost, err := preparePaymentRuntime(ctx, pool, os.Getenv, nil)
	if err != nil {
		slog.Error("payment checkout configuration or credential binding is unavailable")
		os.Exit(2)
	}
	paymentProvider := ""
	paymentTenant, paymentOrganization := "", ""
	if paymentHost != nil {
		paymentProvider = paymentHost.worker.Scope.ProviderCode
		paymentTenant, paymentOrganization = paymentHost.worker.Scope.TenantID, paymentHost.worker.Scope.OrganizationID
	}
	initialHandoverModule, err := selectedInitialHandoverModule(pool, paymentHost, os.Getenv)
	if err != nil {
		slog.Error("initial handover profile or payment scope is invalid")
		os.Exit(2)
	}
	fulfillmentService := fulfillment.NewService(postgres.NewFulfillment(pool), generator)
	queryService := enterprisequery.NewService(postgres.NewEnterpriseQuery(pool))
	roleMetrics, err := postgres.NewRoleMetrics(pool)
	if err != nil {
		slog.Error("role metrics initialization failed")
		os.Exit(2)
	}
	journeyService := policyRuntime.journeyService
	royaltyService := royalty.NewService(postgres.NewRoyalty(pool), generator)
	accountingService := accounting.NewService(postgres.NewAccounting(pool), generator)
	fiscalModule, err := selectedFiscalModule(pool, generator, os.Getenv)
	if err != nil {
		slog.Error("fiscal activation configuration is invalid")
		os.Exit(2)
	}
	providerRegistry, err := providerintegration.LoadHMACConnections([]byte(os.Getenv("PROVIDER_WEBHOOK_CONNECTIONS_JSON")), os.LookupEnv)
	if err != nil {
		slog.Error("provider webhook configuration is invalid")
		os.Exit(2)
	}
	providerService := providerintegration.NewService(postgres.NewProviderIntegration(pool), providerRegistry, providerintegration.HMACSHA256{Now: time.Now, Tolerance: 5 * time.Minute})
	modules := []httpapi.EnterpriseModule{httpapi.OperationsModule{Service: operationsService}, httpapi.InventoryControlModule{Service: inventoryControlService, Bulk: bulkInventoryService, Warehouse: warehouseService, Transfer: bulkTransferService}, httpapi.CommerceModule{Service: commerceService, PaymentProvider: paymentProvider, PaymentTenantID: paymentTenant, PaymentOrganizationID: paymentOrganization, ProviderObservedPayments: true}, httpapi.FulfillmentModule{Service: fulfillmentService}, httpapi.EnterpriseQueryModule{Service: queryService, Metrics: roleMetrics}, httpapi.FranchiseJourneyModule{Service: journeyService}, httpapi.RoyaltyModule{Service: royaltyService}, httpapi.AccountingModule{Service: accountingService}, httpapi.ProviderIntegrationModule{Service: providerService}}
	if paymentHost != nil {
		modules = append(modules, paymentHost.module)
	}
	if initialHandoverModule != nil {
		modules = append(modules, initialHandoverModule)
	}
	if fiscalModule != nil {
		modules = append(modules, fiscalModule)
	}
	feedbackModule, err := selectedCustomerFeedbackModule(pool, os.Getenv)
	if err != nil {
		slog.Error("customer survey activation configuration is invalid")
		os.Exit(2)
	}
	if feedbackModule != nil {
		modules = append(modules, feedbackModule)
	}
	fxModule, err := selectedFXConversionModule(pool, os.Getenv)
	if err != nil {
		slog.Error("FX snapshot configuration is invalid")
		os.Exit(2)
	}
	if fxModule != nil {
		modules = append(modules, fxModule)
	}
	whatsappHost, err := selectedWhatsAppHost(ctx, pool, verifier, os.Getenv)
	if err != nil {
		slog.Error("WhatsApp host activation is unavailable")
		os.Exit(2)
	}
	if whatsappHost != nil {
		defer whatsappHost.close()
		modules = append(modules, whatsappHost)
	}

	if storedValueModuleFactory != nil {
		module, e := storedValueModuleFactory(ctx, pool, os.Getenv)
		if e != nil {
			slog.Error("stored-value activation failed")
			os.Exit(1)
		}
		if module != nil {
			modules = append(modules, module)
		}
	} else if enabled := os.Getenv("STORED_VALUE_ENABLED"); enabled != "" && enabled != "false" {
		slog.Error("stored-value pack is not selected")
		os.Exit(1)
	}
	if warrantyModuleFactory != nil {
		module, e := warrantyModuleFactory(ctx, pool, os.Getenv)
		if e != nil {
			slog.Error("warranty activation failed")
			os.Exit(1)
		}
		if module != nil {
			modules = append(modules, module)
		}
	} else if enabled := os.Getenv("WARRANTY_ENABLED"); enabled != "" && enabled != "false" {
		slog.Error("warranty pack is not selected")
		os.Exit(1)
	}
	if serialSupplyModuleFactory != nil {
		module, e := serialSupplyModuleFactory(ctx, pool, os.Getenv)
		if e != nil {
			slog.Error("serial supply activation failed")
			os.Exit(1)
		}
		if module != nil {
			modules = append(modules, module)
		}
	} else if enabled := os.Getenv("SERIAL_SUPPLY_ENABLED"); enabled != "" && enabled != "false" {
		slog.Error("serial supply pack is not selected")
		os.Exit(1)
	}
	if catalogReleaseModuleFactory != nil {
		module, e := catalogReleaseModuleFactory(ctx, pool, policyRuntime.commerceRepository, os.Getenv)
		if e != nil {
			slog.Error("catalog publication activation failed")
			os.Exit(1)
		}
		if module != nil {
			modules = append(modules, module)
		}
	} else if enabled := os.Getenv("CATALOG_RELEASE_ENABLED"); enabled != "" && enabled != "false" {
		slog.Error("catalog publication pack is not selected")
		os.Exit(1)
	}
	if marketplaceModuleFactory != nil {
		module, e := marketplaceModuleFactory(ctx, pool, policyRuntime.commerceRepository, os.Getenv)
		if e != nil {
			slog.Error("marketplace mutation activation failed")
			os.Exit(1)
		}
		if module != nil {
			modules = append(modules, module)
		}
	} else if enabled := os.Getenv("MARKETPLACE_ENABLED"); enabled != "" && enabled != "false" {
		slog.Error("marketplace mutation pack is not selected")
		os.Exit(1)
	}
	if documentModuleFactory != nil {
		module, e := documentModuleFactory(ctx, pool, os.Getenv)
		if e != nil {
			slog.Error("document configuration rejected")
			os.Exit(2)
		}
		if module != nil {
			modules = append(modules, module)
		}
	} else if enabled := os.Getenv("DOCUMENTS_ENABLED"); enabled != "" && enabled != "false" {
		slog.Error("document pack is not selected")
		os.Exit(2)
	}
	if merchantModuleFactory != nil {
		module, e := merchantModuleFactory(ctx, pool, policyRuntime.commerceRepository, os.Getenv)
		if e != nil {
			slog.Error("Google Merchant activation failed")
			os.Exit(1)
		}
		if module != nil {
			modules = append(modules, module)
		}
	} else if enabled := os.Getenv("MERCHANT_ENABLED"); enabled != "" && enabled != "false" {
		slog.Error("Google Merchant pack is not selected")
		os.Exit(1)
	}
	if trainingModuleFactory != nil {
		module, e := trainingModuleFactory(ctx, pool, os.Getenv)
		if e != nil {
			slog.Error("training activation failed")
			os.Exit(1)
		}
		if module != nil {
			modules = append(modules, module)
		}
	} else if os.Getenv("TRAINING_ENABLED") != "" && os.Getenv("TRAINING_ENABLED") != "false" {
		slog.Error("training pack not selected")
		os.Exit(1)
	}
	if networkRoleModuleFactory != nil {
		module, e := networkRoleModuleFactory(ctx, pool, os.Getenv)
		if e != nil {
			slog.Error("network role activation failed")
			os.Exit(1)
		}
		if module != nil {
			modules = append(modules, module)
		}
	} else if os.Getenv("NETWORK_ROLE_ENABLED") != "" && os.Getenv("NETWORK_ROLE_ENABLED") != "false" {
		slog.Error("network role pack not selected")
		os.Exit(1)
	}
	if helpCMSModuleFactory != nil {
		module, e := helpCMSModuleFactory(ctx, pool, os.Getenv)
		if e != nil {
			slog.Error("help CMS activation failed")
			os.Exit(1)
		}
		if module != nil {
			modules = append(modules, module)
		}
	} else if os.Getenv("HELP_CMS_ENABLED") != "" && os.Getenv("HELP_CMS_ENABLED") != "false" {
		slog.Error("help CMS pack not selected")
		os.Exit(1)
	}
	portalHost, err := selectedPortalHost(ctx, pool, os.Getenv)
	if err != nil {
		slog.Error("portal lifecycle activation is unavailable")
		os.Exit(2)
	}
	if portalHost != nil {
		modules = append(modules, portalHost)
	}
	handler := httpapi.NewEnterprise(orders, verifier, catalog, modules...)
	handler, metricsRun, metricsClose, err := prepareHTTPMetrics(handler, os.Getenv)
	if err != nil {
		slog.Error("HTTP metrics configuration is invalid")
		os.Exit(2)
	}
	var metricsDone chan error
	if metricsRun != nil {
		metricsDone = make(chan error, 1)
		go func() {
			e := metricsRun(ctx)
			metricsDone <- e
			if e != nil {
				stop()
			}
		}()
		defer func() {
			c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if e := metricsClose(c); e != nil {
				slog.Error("metrics provider shutdown failed")
			}
		}()
	}

	server := &http.Server{Addr: envDefault("HTTP_ADDRESS", ":8080"), Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}
	if paymentHost != nil {
		workerContext, cancelWorkers := context.WithCancel(ctx)
		workersDone := make(chan struct{})
		go func() { defer close(workersDone); paymentHost.run(workerContext) }()
		defer func() { cancelWorkers(); <-workersDone }()
	}
	if whatsappHost != nil {
		workerContext, cancelWorkers := context.WithCancel(ctx)
		workersDone := make(chan struct{})
		go func() {
			defer close(workersDone)
			if err := whatsappHost.run(workerContext); err != nil && !errors.Is(err, context.Canceled) {
				slog.Error("WhatsApp worker reporting unavailable; stopping host")
				stop()
			}
		}()
		defer func() { cancelWorkers(); <-workersDone }()
	}
	slog.Info("electromobility api starting", "address", server.Addr)
	if portalHost != nil {
		workerContext, cancelWorkers := context.WithCancel(ctx)
		workersDone := make(chan struct{})
		go func() { defer close(workersDone); portalHost.run(workerContext) }()
		defer func() { cancelWorkers(); <-workersDone }()
	}
	if err = httpapi.ServeUntilShutdown(ctx, server, 15*time.Second); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("server failed")
		os.Exit(1)
	}
	if metricsDone != nil {
		stop()
		if e := <-metricsDone; e != nil && !errors.Is(e, context.Canceled) {
			slog.Error("metrics listener failed")
			os.Exit(1)
		}
	}
	if err := nativeStopFailure(ctx); err != nil {
		slog.Error("native stop protocol failed")
		os.Exit(1)
	}
}
func envDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

// No fiscal credentials, provider or routes are required for an independent
// non-fiscal journey. Activation is explicit; workers remain separate processes.
func paymentRequestProvider(lookup func(string) string) (string, error) {
	switch v := lookup("PAYMENT_REQUEST_PROVIDER"); v {
	case "", "stripe", "mercadopago":
		return v, nil
	default:
		return "", errors.New("PAYMENT_REQUEST_PROVIDER must be absent, stripe or mercadopago")
	}
}

func selectedFiscalModule(pool *pgxpool.Pool, generator randomid.Generator, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	switch lookup("ARCA_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
		provider, err := wsfeipc.NewUnixProvider(lookup("ARCA_WSFE_SOCKET"), 30*time.Second)
		if err != nil {
			return nil, errors.New("ARCA_ENABLED=true requires a valid ARCA_WSFE_SOCKET")
		}
		repository := postgres.NewFiscal(pool)
		return httpapi.FiscalModule{Service: fiscal.NewService(repository, generator), Parameters: fiscal.NewParameterRegistry(repository, provider, generator, time.Now)}, nil
	default:
		return nil, errors.New("ARCA_ENABLED must be true or false")
	}
}

// Optional pack hook preserves independent profiles without stored-value imports.
var storedValueModuleFactory func(context.Context, *pgxpool.Pool, func(string) string) (httpapi.EnterpriseModule, error)

// Optional warranty pack hook preserves independent profiles.
var warrantyModuleFactory func(context.Context, *pgxpool.Pool, func(string) string) (httpapi.EnterpriseModule, error)

// Optional serial supply pack preserves independent profile dependency closure.
var serialSupplyModuleFactory func(context.Context, *pgxpool.Pool, func(string) string) (httpapi.EnterpriseModule, error)

var catalogReleaseModuleFactory func(context.Context, *pgxpool.Pool, *postgres.Commerce, func(string) string) (httpapi.EnterpriseModule, error)

var trainingModuleFactory func(context.Context, *pgxpool.Pool, func(string) string) (httpapi.EnterpriseModule, error)

var networkRoleModuleFactory func(context.Context, *pgxpool.Pool, func(string) string) (httpapi.EnterpriseModule, error)

var helpCMSModuleFactory func(context.Context, *pgxpool.Pool, func(string) string) (httpapi.EnterpriseModule, error)

// Optional provider mutation composition; absence cannot silently enable effects.
var marketplaceModuleFactory func(context.Context, *pgxpool.Pool, *postgres.Commerce, func(string) string) (httpapi.EnterpriseModule, error)

// Optional Google Merchant composition; missing pack fails activation closed.
var merchantModuleFactory func(context.Context, *pgxpool.Pool, *postgres.Commerce, func(string) string) (httpapi.EnterpriseModule, error)

var documentModuleFactory func(context.Context, *pgxpool.Pool, func(string) string) (httpapi.EnterpriseModule, error)
````

### FILE: `internal/platform/randomid/generator.go`

```yaml
block_id: "GO-EM-APP:random-id:v1"
operation: CREATE
provenance: AUTHORED
source: "local generator governed by official Go crypto/rand documentation and RFC 4122 UUID v4 bit layout"
license: "LicenseRef-Workspace-Owner"
sha256: "a5894d48b94c44149b6bd5966d2a9423a23021b52897b2a2c660d1b3ed195120"
variables: []
secrets_allowed: false
```

````go
package randomid

import (
	"crypto/rand"
	"encoding/hex"
)

// Generator creates RFC 4122 version 4 identifiers from the operating system CSPRNG.
// A CSPRNG failure is unrecoverable because continuing could break durable identity.
type Generator struct{}

func (Generator) New() string {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		panic(err)
	}
	value[6] = (value[6] & 0x0f) | 0x40
	value[8] = (value[8] & 0x3f) | 0x80
	encoded := hex.EncodeToString(value[:])
	return encoded[0:8] + "-" + encoded[8:12] + "-" + encoded[12:16] + "-" + encoded[16:20] + "-" + encoded[20:32]
}
````

### FILE: `internal/platform/randomid/generator_test.go`

```yaml
block_id: "GO-EM-APP:random-id-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local UUID version and uniqueness regression"
license: "LicenseRef-Workspace-Owner"
sha256: "e5675ea7fd971cf35a48f9e31891e5d7d4eca555fb1d3fd7235a5218fca52098"
variables: []
secrets_allowed: false
```

````go
package randomid

import (
	"regexp"
	"testing"
)

var uuidV4Pattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestGeneratorProducesUniqueVersion4Identifiers(t *testing.T) {
	generator := Generator{}
	seen := make(map[string]struct{}, 256)
	for range 256 {
		value := generator.New()
		if !uuidV4Pattern.MatchString(value) {
			t.Fatalf("invalid UUID v4: %q", value)
		}
		if _, exists := seen[value]; exists {
			t.Fatalf("duplicate UUID v4: %q", value)
		}
		seen[value] = struct{}{}
	}
}
````

### FILE: `docs/electromobility-api-runtime.md`

```yaml
block_id: "GO-EM-APP:runtime-doc:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "3d488453bf876e751d00bffa2a6e4c54982c2255c5a223a58f4cf323af17a061"
variables: []
secrets_allowed: false
```

````markdown
# Electromobility API runtime

Required configuration: `DATABASE_URL` (secret when it contains credentials), `OIDC_ISSUER` and `OIDC_AUDIENCE`; optional `HTTP_ADDRESS` defaults to `:8080`. Apply the complete migration manifest of the selected composition before starting (not only the historical first twelve migrations). Inject database credentials through the selected secret mechanism rather than committing them.

Fiscal activation is explicit. `ARCA_ENABLED` absent or exactly `false` omits FiscalModule and all its HTTP routes without reading `ARCA_WSFE_SOCKET`. Exactly `true` requires a valid absolute socket path and registers the existing authenticated fiscal routes. Other flag values fail closed. The flag controls this API process only: while ARCA is deferred, do not deploy/start the ARCA parameter/fiscal workers, return-fiscal-worker or WSAA/WSFE sidecars. Keep their migrations and pending fiscal dependencies, but do not label invoicing or fiscal returns enabled. Re-enable only after the existing fiscal admission, credentials, country/POS and transport tests; no credential is needed to postpone ARCA. Other modules are not dynamically disabled by business JSON flags.

`PROVIDER_WEBHOOK_CONNECTIONS_JSON` is optional and contains only `tenant_id`, `connection_id`, `provider_code` and `secret_env`. Every named secret variable is resolved separately at startup and must contain at least 32 characters; inline secrets and unknown fields are rejected. The reference HMAC adapter signs `unix_timestamp + "." + exact_body`, uses the `v1=<hex>` signature form and admits at most five minutes of clock skew. Replace it with a provider-specific verifier when the official provider contract differs.

The process fails closed when configuration, OIDC discovery or PostgreSQL connectivity is unavailable. It serves liveness and order routes plus public catalog/lead/price/location/appointment, administrative catalog, procurement, factory, inventory, price-book, allocation, payment-intent, logistics, service, recall, franchise, royalty policy/accrual/settlement/reversal/reconciliation, accounting account/period/journal/post/reversal/trial-balance, fiscal point-of-sale/request/read, communication and visitor-to-customer journey routes. It also exposes bounded queries for administrative overview/orders/service/leads, factory units and customer-owned orders/service/appointments/quotes/handovers. Tenant identity and the `organization_ids` resource scope for protected routes come from a verified token; wildcard permission is the explicit tenant-wide administrative escape hatch. Organization-scoped mutations and reads bind the organization in PostgreSQL predicates; customer reads and delivery acceptance additionally bind the verified subject. Public tenant and organization codes are resolved server-side. Provider webhooks resolve tenant from a preconfigured connection, authenticate the exact body, deduplicate by provider event and atomically enqueue durable processing. Fiscal requests derive currency and total from the captured payment/order, persist a lane and outbox, and require the separate ARCA WSAA/generated-WSFE transport worker for actual authorization and reconciliation. Payment and communication routes still require provider-specific outbound adapters and reconciliation workers.

Production must place the process behind a trusted TLS edge, define proxy/header policy, distributed antiabuse for public routes, workload identity, telemetry export, readiness distinct from liveness, database pool budgets, graceful-drain coordination and the deployment/rollback pack selected by the project.
````

### FILE: `cmd/arca-parameter-worker/main.go`

```yaml
block_id: "GO-EM-APP:arca-parameter-worker:v1"
operation: CREATE
provenance: AUTHORED
source: "local composition using the fiscal parameter owner and UDS provider"
license: "LicenseRef-Workspace-Owner"
sha256: "a4e745c1754c2427f0e44de0966688223534c9ef45e2dfa0f72e576aeeb3ddf4"
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

var parameterWorkerPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$`)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	databaseURL, socketPath, workerID := os.Getenv("DATABASE_URL"), os.Getenv("ARCA_WSFE_SOCKET"), os.Getenv("FISCAL_PARAMETER_WORKER_ID")
	if databaseURL == "" || socketPath == "" || !parameterWorkerPattern.MatchString(workerID) {
		slog.Error("DATABASE_URL, ARCA_WSFE_SOCKET and a valid FISCAL_PARAMETER_WORKER_ID are required")
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
	repository, ids := postgres.NewFiscal(pool), randomid.Generator{}
	registry := fiscal.NewParameterRegistry(repository, provider, ids, time.Now)
	processor, err := fiscal.NewParameterRefreshProcessor(repository, registry, ids, workerID, 2*time.Minute, 5*time.Minute)
	if err != nil {
		slog.Error("parameter processor configuration is invalid")
		os.Exit(2)
	}
	if err = runParameterWorker(ctx, processor); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("parameter worker stopped unexpectedly")
		os.Exit(1)
	}
}

func runParameterWorker(ctx context.Context, processor *fiscal.ParameterRefreshProcessor) error {
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()
	for {
		err := processor.ProcessOne(ctx)
		if err == nil {
			continue
		}
		if errors.Is(err, context.Canceled) {
			return err
		}
		wait := 2 * time.Second
		if errors.Is(err, fiscal.ErrNoParameterWork) {
			wait = 30 * time.Second
		} else {
			slog.Warn("parameter refresh deferred after recoverable failure")
		}
		timer.Reset(wait)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
}
````

### FILE: `cmd/electromobility-api/whatsapp_activation.go`

```yaml
block_id: "COMMUNICATIONS-DELTA-GO_ELECTROMOBILITY_APPLICATION:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8338bbe16067463721259f4da4fab799161707326e366b9f125068d6f58ae58e"
variables: []
secrets_allowed: false
```

````go
package main

// AUTHORED optional composition hook. This file belongs to the base API pack;
// it does not import WhatsApp/AI owners into a backend-only composition.
import (
	"context"
	"errors"

	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type whatsappRuntime interface {
	httpapi.EnterpriseModule
	run(context.Context) error
	close() error
}

var whatsappRuntimeFactory func(context.Context, *pgxpool.Pool, identity.Verifier, func(string) string) (whatsappRuntime, error)

func selectedWhatsAppHost(ctx context.Context, pool *pgxpool.Pool, verifier identity.Verifier, lookup func(string) string) (whatsappRuntime, error) {
	if lookup == nil {
		return nil, errors.New("WhatsApp activation configuration unavailable")
	}
	switch lookup("WHATSAPP_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
		if whatsappRuntimeFactory == nil {
			return nil, errors.New("WhatsApp host pack is not selected")
		}
		return whatsappRuntimeFactory(ctx, pool, verifier, lookup)
	default:
		return nil, errors.New("WHATSAPP_ENABLED must be true, false or absent")
	}
}
````

### FILE: `cmd/electromobility-api/whatsapp_activation_test.go`

```yaml
block_id: "COMMUNICATIONS-DELTA-GO_ELECTROMOBILITY_APPLICATION:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c1a5c1a4d216a7fb8aa0db9f0ee29659e15b81b78d44dd2a2aa8ec9a0be1afc8"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"testing"
)

func TestWhatsAppActivationWithoutOptionalPack(t *testing.T) {
	saved := whatsappRuntimeFactory
	whatsappRuntimeFactory = nil
	defer func() { whatsappRuntimeFactory = saved }()
	if h, e := selectedWhatsAppHost(context.Background(), nil, nil, func(string) string { return "true" }); e == nil || h != nil {
		t.Fatal("activation without selected implementation")
	}
	if h, e := selectedWhatsAppHost(context.Background(), nil, nil, func(k string) string {
		if k != "WHATSAPP_ENABLED" {
			t.Fatal("read inactive secret/config")
		}
		return "false"
	}); e != nil || h != nil {
		t.Fatal("disabled optional owner")
	}
}
````

## 6. Configuration surface

| Variable | Type/default | Secret | Validation/effect |
|---|---|---|---|
| `DATABASE_URL` | PostgreSQL URL / none | yes | required; startup connectivity fails closed |
| `OIDC_ISSUER` | HTTPS URL / none | no | required discovery issuer |
| `OIDC_AUDIENCE` | non-empty string / none | no | required token audience |
| `ARCA_WSFE_SOCKET` | absolute Unix socket path / none | no | required by API parameter administration and both fiscal workers |
| `FISCAL_PARAMETER_WORKER_ID` | bounded worker identifier / none | no | required by parameter refresh worker; owns leases only, never approval |
| `HTTP_ADDRESS` | address / `:8080` | no | validated bind target |
| `PROVIDER_WEBHOOK_CONNECTIONS_JSON` | strict JSON / omitted | no | connection metadata and secret variable names only |
| referenced provider secret variables | at least 32 chars / none | yes | resolved separately; missing/weak values fail startup |

## 7. Dependency bill

| Dependency | Pin | Use | License | Scope | Source |
|---|---|---|---|---|---|
| Go | `1.26.7` | composition/runtime | BSD-3-Clause | build/runtime | `go.dev` |
| PostgreSQL | `18.6` verified baseline | source of truth | PostgreSQL | runtime/test | `postgresql.org` |
| `go-oidc` and transitive OAuth packages | versions fixed by composed `go.mod` | token verification/discovery | Apache-2.0/BSD | build/runtime | official repositories |

## 8. Apply order

Materialize last among backend code packs, after migrations 0001–0027 and module interfaces; resolve collisions, inject external configuration, run migrations, execute all tests/vet/build, then start API and workers behind the selected trusted edge. Existing applications integrate routes one capability at a time. Rollback stops new traffic/workers and redeploys the previous binary compatible with the expanded schema.

## 9. Verification

The clean scoped composition must pass formatting, all Go tests against PostgreSQL migrations 0001–0027, vet and every build. Version 1.8.0 wires the same inventory repository into serial, bulk/lot/bin/FIFO-specific and operational warehouse services without duplicating ownership. Fiscal and provider owners retain their existing gates. This proves durable local wiring, not full Business Central WMS/planning/cost adjustment, live ARCA homologation, approved accounting policy, production identity, edge, telemetry or deployment.

## 10. Reconstruction evidence

Component history through fiscal parameter administration remains recorded in V128–V133; V159 records serial reservation/ATP/transfer wiring, V160 records bulk/lot/bin/FIFO-specific composition and V161 records the operational receipt/put-away/pick/FEFO composition and backend reconstruction.

V400 wires the explicit CUSTOMER_SURVEYS_ENABLED selector from GO-CUSTOMER-SURVEY-API0.1.0. All application compositions select that dependency; absent/0 mounts no routes,1 requires database,other values fail. See CONNECTED_CUSTOMER_SURVEYS_V400.md. No provider or policy activation is inferred.

V402 composed delta: Payment/initial-handover composition. New behavior and tests belong to the explicit runtime/portal packs; source and library release claims remain bounded to their evidence. Existing source provenance is preserved.

Canonical V402 integration: selected by the current profile with exact dependencies and caller overlays. Metadata promotion records byte reconstruction, not closure of every admission/release gate. Payload provenance is unchanged.

V402 composed delta: Optional WhatsApp host activation hook plus exact FX module composition. Disabled host reads no external inputs; absent optional pack fails closed only when explicitly enabled. No AI dependencies added to the base-only host.

V402 composed delta: V402 source-backed stored-value integration: exact remaining provider due, explicit payment/funding XOR, shared approval, bounded browser transport and optional host. See STORED_VALUE_OPERATOR_FLOW_V402.md; source/pack admission successor governs final claim. Existing provider-only behavior retained.

V402 composed delta: Connected warranty reuses existing transaction/approval/stock/service owners; SQL ordering and public wrapper behavior retained. Optional host factory fails closed. Exact source tested in WARRANTY_INTERFACE_AND_PORTABILITY_V402.md; no new dependency or corporate attribution.

V402 composed delta: Connected J2 uses existing transaction owners and preserves public operations transitions; serial quality remains in the shared distinct-human approval owner. Optional host hook keeps narrower profiles compatible. No dependency added or corporate attribution. SERIAL_SUPPLY_CONNECTED_RELEASE_V402.md.

V402 composed delta: J3 immutable approved catalog publication reuses original Commerce SQL/shared approval and optional host/public model owner; existing Next storefront consumes a validated published projection. No new dependencies or corporate attribution. CATALOG_CONNECTED_RELEASE_V402.md.

V402 composed delta: T2804 connected training: existing versioned help/audit/shared approvals/outbox/BFF; opt-in host and navigation; bounded body retains original default. No new dependency or automatic grant. TRAINING_CONNECTED_RELEASE_V402.md.

V402 composed delta: T2804 network role: original four fulfillment SQL bodies extracted unchanged into one transaction with immutable result; optional host, forms and GET recovery. Migration0077, no dependency/domain-rule change. NETWORK_ROLE_RELEASE_V402.md.

V402 composed delta: T2804 help CMS and21same-release guides;15oldguides unchanged,5bounded training curricula revision2, no automatic grants. New optional host, shared inline text, actual browser/PG proofs. HELP_CMS_RELEASE_V402.md.

V402 composed delta: T2804 role metrics use existing domain read models with exact strings, organization/customer/factory/program permissions, NPS minimum/retention and no zero on unavailable. FAIL868 converted lead and FAIL457 bounded generic body corrected. ROLE_METRICS_RELEASE_V402.md/json; no new dependency or corporate authorship.

V402 composed delta: T2805 narrow connected Mercado Libre PRICE/STOCK/PAUSE/RESUME for existing User Products item: immutable current catalog + existing serial ATP + exact distinct human approval + shared one-attempt fence + GET-only recovery. MARKETPLACE_MUTATION_RELEASE_V402.md/json. AUTHORED HTTP/SQL/host/proof glue; initial publication/media and other T2805 work remain open.

V402 composed delta: Optional Google Merchant host activation slot; disabled profiles unchanged; MERCHANT_CONNECTED_RELEASE_V402.md/json.

### FILE: `cmd/electromobility-api/portal_activation.go`

```yaml
block_id: "PORTAL-EXTENSION-1:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "eb8d00579454becb4a6b430d24c720558a1c61edc4031e304d5e95b4d9c38cad"
variables: []
secrets_allowed: false
```

````go
package main

// AUTHORED optional hook owned by the base API pack. No frontend/portal import.
import (
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type portalRuntime interface {
	httpapi.EnterpriseModule
	run(context.Context)
}

var portalRuntimeFactory func(context.Context, *pgxpool.Pool, func(string) string) (portalRuntime, error)

func selectedPortalHost(ctx context.Context, pool *pgxpool.Pool, getenv func(string) string) (portalRuntime, error) {
	if getenv == nil {
		return nil, errors.New("portal activation unavailable")
	}
	switch getenv("OIDC_PORTAL_LIFECYCLE_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
		if portalRuntimeFactory == nil {
			return nil, errors.New("portal lifecycle pack is not selected")
		}
		return portalRuntimeFactory(ctx, pool, getenv)
	default:
		return nil, errors.New("portal lifecycle enabled must be true or false")
	}
}
````


V402 composed delta: Portal lifecycle optional host/BFF integration; IDENTITY_PORTAL_RELEASE_V402.md/json; no change to closed business journeys.

V402 composed delta: T2806 connected document reference. Optional host hook, explicit bound document-review kind, exact AWS module closure and preserved security-floor checksums; no unchanged business policy modified. DOCUMENT_REFERENCE_RELEASE_V402.md/json.

V402 composed delta: Local delivery316: native stop owner reused, Next standalone exact build identity, container template without mutable defaults and complete local module context. Local fixture qualification only; docs/LOCAL_REFERENCE_DELIVERY.md.

### FILE: `cmd/electromobility-api/native_stop_other.go`

```yaml
block_id: "GO-ELECTROMOBILITY-APPLICATION-LOCAL-DELIVERY:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9bade57a4db708cdc464970eec0e7d129cba9cffdf48b69db89351683bbf9efa"
variables: []
secrets_allowed: false
```

````go
//go:build !windows

package main

import (
	"context"
	"errors"
)

// Native handles are a Windows-only trusted-launcher protocol.
func nativeStopContext(parent context.Context, raw string) (context.Context, func() error, error) {
	if raw != "" {
		return parent, func() error { return nil }, errors.New("NATIVE_STOP_PROTOCOL_FAILED")
	}
	ctx, cancel := context.WithCancel(parent)
	return ctx, func() error { cancel(); return nil }, nil
}

func nativeStopFailure(context.Context) error { return nil }
````

### FILE: `cmd/electromobility-api/native_stop_windows.go`

```yaml
block_id: "GO-ELECTROMOBILITY-APPLICATION-LOCAL-DELIVERY:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "97d705c2674f7d364d80d22a1029ab27978ebb5dfa15f72fad94c5a0b9668a91"
variables: []
secrets_allowed: false
```

````go
//go:build windows

package main

// AUTHORED qualification only. The launcher supplies a trusted manual-reset
// event; this is neither object-type authentication nor hostile-worker isolation.
import (
	"context"
	"errors"
	"strconv"
	"sync"
	"syscall"
	"unsafe"
)

var errNativeStop = errors.New("NATIVE_STOP_REQUESTED")
var errNativeProtocol = errors.New("NATIVE_STOP_PROTOCOL_FAILED")
var stopKernel = syscall.NewLazyDLL("kernel32.dll")
var stopDuplicate = stopKernel.NewProc("DuplicateHandle")
var stopWait = stopKernel.NewProc("WaitForSingleObject")
var stopClose = stopKernel.NewProc("CloseHandle")

func parseStopHandle(raw string) (uintptr, error) {
	n, err := strconv.ParseUint(raw, 10, strconv.IntSize)
	if err != nil || n == 0 || n > uint64(^uintptr(0)>>1) || strconv.FormatUint(n, 10) != raw {
		return 0, errNativeProtocol
	}
	return uintptr(n), nil
}

// Duplicate before waiting. Cleanup joins the waiter BEFORE closing its handle:
// CloseHandle during a pending Windows wait has undefined behavior.
// Empty raw retains the ordinary parent-context lifecycle. Cancellation stops
// claims through the real host loop; it does not guarantee a domain commit.
func nativeStopContext(parent context.Context, raw string) (context.Context, func() error, error) {
	ctx, cancel := context.WithCancelCause(parent)
	if raw == "" {
		return ctx, func() error { cancel(context.Canceled); return nil }, nil
	}
	source, err := parseStopHandle(raw)
	if err != nil {
		cancel(err)
		return ctx, func() error { return nil }, err
	}
	var owned uintptr
	ok, _, _ := stopDuplicate.Call(^uintptr(0), source, ^uintptr(0), uintptr(unsafe.Pointer(&owned)), 0x100000, 0, 0)
	if ok == 0 {
		cancel(errNativeProtocol)
		return ctx, func() error { return nil }, errNativeProtocol
	}
	done, joined := make(chan struct{}), make(chan struct{})
	var once sync.Once
	var closeErr error
	cleanup := func() error {
		once.Do(func() {
			close(done)
			<-joined
			ok, _, _ := stopClose.Call(owned)
			if ok == 0 {
				closeErr = errNativeProtocol
			}
			cancel(context.Canceled)
		})
		return closeErr
	}
	// A pre-signaled event must not expose a live context to the first claim.
	// The trusted launcher supplies a manual-reset event, so this read is not
	// destructive. Other waitable object types are outside the protocol.
	initial, _, _ := stopWait.Call(owned, 0)
	if initial != 258 {
		close(joined)
		if initial == 0 {
			cancel(errNativeStop)
			return ctx, cleanup, nil
		}
		cancel(errNativeProtocol)
		_ = cleanup()
		return ctx, cleanup, errNativeProtocol
	}
	go func() {
		defer close(joined)
		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			default:
			}
			status, _, _ := stopWait.Call(owned, 20)
			switch status {
			case 0:
				cancel(errNativeStop)
				return
			case 258:
			default:
				cancel(errNativeProtocol)
				return
			}
		}
	}()
	return ctx, cleanup, nil
}

func nativeStopFailure(ctx context.Context) error {
	if errors.Is(context.Cause(ctx), errNativeProtocol) {
		return errNativeProtocol
	}
	return nil
}
````

### FILE: `cmd/electromobility-api/native_stop_windows_test.go`

```yaml
block_id: "GO-ELECTROMOBILITY-APPLICATION-LOCAL-DELIVERY:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "4eba1980e6cf549db8636f6bc9c71065918b3048dd7b3dc0e216e997c58949a9"
variables: []
secrets_allowed: false
```

````go
//go:build windows

package main

import (
	"context"
	"strconv"
	"testing"
	"time"
)

func TestAPINativeStopLifecycle(t *testing.T) {
	h, _, _ := stopKernel.NewProc("CreateEventW").Call(0, 1, 0, 0)
	if h == 0 {
		t.Fatal("event creation")
	}
	defer stopClose.Call(h)
	ctx, cleanup, err := nativeStopContext(context.Background(), strconv.FormatUint(uint64(h), 10))
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	if ctx.Err() != nil {
		t.Fatal("premature stop")
	}
	stopKernel.NewProc("SetEvent").Call(h)
	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("stop did not cancel host")
	}
	if err := cleanup(); err != nil {
		t.Fatal(err)
	}
	ctx, cleanup, err = nativeStopContext(context.Background(), strconv.FormatUint(uint64(h), 10))
	if err != nil || ctx.Err() == nil {
		t.Fatal("pre-signaled stop must precede host startup")
	}
	if err := cleanup(); err != nil {
		t.Fatal(err)
	}
}

func TestAPINativeStopInvalidHandle(t *testing.T) {
	for _, raw := range []string{"0", "-1", "+4", "04", "999999999999999999999999"} {
		if _, _, err := nativeStopContext(context.Background(), raw); err == nil {
			t.Fatal("invalid handle accepted")
		}
	}
}
````


V402316: optional Windows local reference commands require PNPM-ARTIFACT-SELECTION-GATE0.9.0 and four exact WINDOW​​S-REFERENCE-TELEMETRY-RUNTIME0.1.0 supervisor files. Original portable runner unchanged. Read docs/LOCAL_REFERENCE_DELIVERY.md. No live production or corporate attribution.

V402 composed delta: V402317 connected local API/Next telemetry, current OIDC, fixed official middleware, finite supervised alert/fault/load/WAL recovery. Historical lock kept separate; docs/LOCAL_REFERENCE_OPERATIONS.md. No production admission.

### FILE: `cmd/electromobility-api/http_metrics_profile.go`

```yaml
block_id: "GO-ELECTROMOBILITY-APPLICATION-LOCAL-OPERATIONS:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "cb6d02fc9f57ae944d963648aa8ae8061b49e153ef9bf3db4dbe68c0d876155f"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"errors"
	"net/http"
)

// Optional owner hook: smaller profiles keep the same host without importing
// metrics dependencies. Explicit activation without its owner fails closed.
var httpMetricsFactory func(http.Handler, func(string) string) (http.Handler, func(context.Context) error, func(context.Context) error, error)

func prepareHTTPMetrics(h http.Handler, lookup func(string) string) (http.Handler, func(context.Context) error, func(context.Context) error, error) {
	switch lookup("HTTP_METRICS_ENABLED") {
	case "", "false":
		return h, nil, nil, nil
	case "true":
		if httpMetricsFactory == nil {
			return nil, nil, nil, errors.New("HTTP metrics owner is not selected")
		}
		return httpMetricsFactory(h, lookup)
	default:
		return nil, nil, nil, errors.New("HTTP_METRICS_ENABLED must be true or false")
	}
}
````


V402317: current host profile is documented in LOCAL_REFERENCE_OPERATIONS.md; source-lock separates historical synthetic principal from current OIDC host. Narrow local proof only.
