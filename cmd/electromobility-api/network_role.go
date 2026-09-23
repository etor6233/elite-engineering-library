package main

// AUTHORED optional host guard; organization activation is not production readiness.
import (
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errNetworkRoleConfiguration = errors.New("network role activation invalid")

func init() { networkRoleModuleFactory = selectedNetworkRoleModule }
func selectedNetworkRoleModule(ctx context.Context, pool *pgxpool.Pool, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	if lookup == nil {
		return nil, errNetworkRoleConfiguration
	}
	switch lookup("NETWORK_ROLE_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errNetworkRoleConfiguration
	}
	if pool == nil {
		return nil, errNetworkRoleConfiguration
	}
	var ready bool
	e := pool.QueryRow(ctx, `select
 (select count(*)from pg_trigger where not tgisinternal and tgenabled in('O','A')and(
 (tgrelid=to_regclass('franchise.network_command_receipt')and tgname='network_receipt_immutable')or
 (tgrelid=to_regclass('org.organization')and tgname='organization_hierarchy_guard')or
 (tgrelid=to_regclass('franchise.agreement')and tgname='franchise_territory_non_overlap')))=3
 and exists(select 1 from pg_constraint where conrelid=to_regclass('franchise.network_command_receipt')and contype='p'and convalidated)
 and exists(select 1 from pg_index where indexrelid=to_regclass('franchise.network_receipt_actor_idx')and indisvalid and indisready)`).Scan(&ready)
	if e != nil || !ready {
		return nil, errNetworkRoleConfiguration
	}
	store, e := postgres.NewNetworkRole(pool)
	if e != nil {
		return nil, errNetworkRoleConfiguration
	}
	return httpapi.NetworkRoleModule{Service: store}, nil
}
