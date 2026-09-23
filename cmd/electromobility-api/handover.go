package main

// AUTHORED activation glue. Profile algorithms and transaction invariants remain
// with the existing handover owner; this host binds its activation to payments.
import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errHandoverConfiguration = errors.New("initial handover activation or payment scope is invalid")

func loadInitialHandoverContract(paymentHost *paymentRuntime, lookup func(string) string) (*franchisejourney.HandoverReleaseContract, error) {
	if lookup == nil {
		return nil, errHandoverConfiguration
	}
	switch lookup("HANDOVER_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errHandoverConfiguration
	}
	if paymentHost == nil || paymentHost.worker.Scope.Validate() != nil {
		return nil, errHandoverConfiguration
	}
	scope := paymentHost.worker.Scope
	revision, err := strconv.Atoi(lookup("HANDOVER_PROFILE_REVISION"))
	if err != nil || revision < 1 {
		return nil, errHandoverConfiguration
	}
	path := lookup("HANDOVER_PROFILE_FILE")
	if path == "" || len(path) > 4096 || strings.TrimSpace(path) != path || strings.ContainsAny(path, "\x00\r\n") {
		return nil, errHandoverConfiguration
	}
	activation := franchisejourney.HandoverProfileActivation{Enabled: true, ProfileID: lookup("HANDOVER_PROFILE_ID"), ProfileRevision: revision, DocumentSHA256: lookup("HANDOVER_PROFILE_SHA256"), TenantID: scope.TenantID, OrganizationID: scope.OrganizationID, ExpectedLiveMode: scope.LiveMode, ProviderCode: scope.ProviderCode, ProviderAccountRef: scope.AccountRef, ProviderConnectionID: scope.ConnectionID}
	contract, err := franchisejourney.LoadHandoverProfileFile(path, activation)
	if err != nil || !contract.Valid() || contract.Scope != "MATERIALIZED_PROFILE" || !contract.AllowsScope(scope.TenantID, scope.OrganizationID) || contract.ExpectedLiveMode != scope.LiveMode {
		return nil, errHandoverConfiguration
	}
	return &contract, nil
}
func selectedInitialHandoverModule(pool *pgxpool.Pool, paymentHost *paymentRuntime, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	contract, err := loadInitialHandoverContract(paymentHost, lookup)
	if err != nil || contract == nil {
		return nil, err
	}
	if pool == nil {
		return nil, errHandoverConfiguration
	}
	service, err := franchisejourney.NewHandoverPreparationService(postgres.NewFranchiseJourney(pool), randomid.Generator{}, *contract)
	if err != nil {
		return nil, errHandoverConfiguration
	}
	return initialHandoverHostModule{service: service, commercial: contract.AllowsCommercialRelease(paymentHost.worker.Scope.TenantID, paymentHost.worker.Scope.OrganizationID)}, nil
}

// Both routes share the exact profile-bound owner. The read-only profile never
// mounts a commercial mutation route, even if a caller guesses its URL.
type initialHandoverHostModule struct {
	service    *franchisejourney.HandoverPreparationService
	commercial bool
}

func (m initialHandoverHostModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	httpapi.InitialHandoverModule{Service: m.service}.Register(mux, verifier)
	httpapi.HandoverContextModule{Service: m.service}.Register(mux, verifier)
	if m.commercial {
		httpapi.CommercialReleaseModule{Service: m.service}.Register(mux, verifier)
	}
}
