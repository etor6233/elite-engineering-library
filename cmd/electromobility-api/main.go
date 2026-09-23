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
	modules := []httpapi.EnterpriseModule{httpapi.OperationsModule{Service: operationsService}, selectedInventoryModule(inventoryControlService, bulkInventoryService, warehouseService, bulkTransferService), httpapi.CommerceModule{Service: commerceService, PaymentProvider: paymentProvider, PaymentTenantID: paymentTenant, PaymentOrganizationID: paymentOrganization, ProviderObservedPayments: true}, httpapi.FulfillmentModule{Service: fulfillmentService}, httpapi.EnterpriseQueryModule{Service: queryService, Metrics: roleMetrics}, httpapi.FranchiseJourneyModule{Service: journeyService}, selectedRoyaltyModule(royaltyService), httpapi.AccountingModule{Service: accountingService}, httpapi.ProviderIntegrationModule{Service: providerService}}
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
	modules = append(modules, selectedFinanceWorkspaceModule(pool, generator, os.Getenv))
	modules = append(modules, httpapi.WarehouseWorkspaceModule{Reader: inventoryRepository, Reservations: inventoryRepository})
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
