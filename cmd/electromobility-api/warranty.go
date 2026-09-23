package main

// AUTHORED optional local activation; credentials are neither read nor needed.
import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errWarrantyConfiguration = errors.New("warranty activation is invalid")

func init() { warrantyModuleFactory = selectedWarrantyModule }
func selectedWarrantyModule(ctx context.Context, pool *pgxpool.Pool, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	if lookup == nil {
		return nil, errWarrantyConfiguration
	}
	switch lookup("WARRANTY_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errWarrantyConfiguration
	}
	if pool == nil {
		return nil, errWarrantyConfiguration
	}
	path := lookup("WARRANTY_PROFILE_FILE")
	if !filepath.IsAbs(path) {
		return nil, errWarrantyConfiguration
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, errWarrantyConfiguration
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() < 1 || info.Size() > 32768 {
		return nil, errWarrantyConfiguration
	}
	raw, err := io.ReadAll(io.LimitReader(file, 32769))
	if err != nil {
		return nil, errWarrantyConfiguration
	}
	profile, err := wc.LoadProfile(raw, lookup("WARRANTY_PROFILE_SHA256"))
	if err != nil {
		return nil, errWarrantyConfiguration
	}
	tenant, org := profile.Scope()
	if tenant != lookup("WARRANTY_TENANT_ID") || org != lookup("WARRANTY_ORGANIZATION_ID") {
		return nil, errWarrantyConfiguration
	}
	store, err := postgres.NewWarranty(pool, profile)
	if err != nil {
		return nil, errWarrantyConfiguration
	}
	// Fail before mounting routes if the selected composition lacks the schema
	// or its immutable-fact/connected-command guards.
	var ready bool
	err = pool.QueryRow(ctx, `select to_regclass('service_ops.warranty_offer') is not null and to_regclass('service_ops.warranty_activation') is not null
 and to_regclass('service_ops.warranty_claim_step') is not null and to_regclass('service_ops.warranty_part_reservation') is not null
 and (select count(*) from pg_trigger where not tgisinternal and tgenabled in ('O','A') and
 (tgrelid=to_regclass('service_ops.service_case') and tgname='warranty_claim_case_guard'
 or tgrelid=to_regclass('sales.quotation') and tgname='warranty_quote_acceptance_guard'
 or tgrelid=to_regclass('service_ops.warranty_work_plan') and tgname='warranty_plan_commit_expiry'
 or tgrelid=to_regclass('service_ops.warranty_claim_step') and tgname='warranty_approval_commit_expiry'))=4`).Scan(&ready)
	if err != nil || !ready {
		return nil, errWarrantyConfiguration
	}
	return httpapi.WarrantyModule{Service: store}, nil
}
