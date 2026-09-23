package main

// AUTHORED optional CMS host, no user credentials or automatic grants.
import (
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errHelpCMSConfiguration = errors.New("help CMS activation invalid")

func init() { helpCMSModuleFactory = selectedHelpCMSModule }
func selectedHelpCMSModule(ctx context.Context, pool *pgxpool.Pool, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	if lookup == nil {
		return nil, errHelpCMSConfiguration
	}
	switch lookup("HELP_CMS_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errHelpCMSConfiguration
	}
	if pool == nil {
		return nil, errHelpCMSConfiguration
	}
	var ready bool
	e := pool.QueryRow(ctx, `select current_setting('server_encoding')='UTF8'and to_regcollation('pg_catalog.pg_unicode_fast')is not null and
 (select count(*)from pg_trigger where not tgisinternal and tgenabled in('O','A')and(
 (tgrelid=to_regclass('help.article')and tgname='help_head_guard')or(tgrelid=to_regclass('help.revision')and tgname='help_revision_guard')))=2
 and exists(select 1 from pg_constraint where conrelid=to_regclass('help.article')and conname='help_current_revision_fk'and contype='f'and convalidated and condeferrable and condeferred)
 and exists(select 1 from pg_constraint where conrelid=to_regclass('help.revision')and conname='revision_tenant_id_command_id_key'and contype='u'and convalidated)
 and exists(select 1 from pg_index where indexrelid=to_regclass('help.help_scope_page_idx')and indisvalid and indisready)`).Scan(&ready)
	if e != nil || !ready {
		return nil, errHelpCMSConfiguration
	}
	store, e := postgres.NewHelpCMS(pool)
	if e != nil {
		return nil, errHelpCMSConfiguration
	}
	return httpapi.HelpCMSModule{Service: store}, nil
}
