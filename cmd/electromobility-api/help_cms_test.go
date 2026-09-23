package main

import (
	"context"
	cms "elite.local/enterprise/internal/helpcms"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strings"
	"testing"
)

func TestHelpCMSHostGuards(t *testing.T) {
	ctx := context.Background()
	calls := 0
	if m, e := selectedHelpCMSModule(ctx, nil, func(k string) string {
		calls++
		if k != "HELP_CMS_ENABLED" {
			t.Fatal("unexpected lookup")
		}
		return "false"
	}); m != nil || e != nil || calls != 1 {
		t.Fatal(m, e, calls)
	}
	if _, e := selectedHelpCMSModule(ctx, nil, func(string) string { return "TRUE" }); e == nil {
		t.Fatal("invalid flag")
	}
	raw := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if raw == "" {
		t.Skip("owned database required")
	}
	pool, e := pgxpool.New(ctx, raw)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	lookup := func(string) string { return "true" }
	enabled := func(want bool) {
		t.Helper()
		m, e := selectedHelpCMSModule(ctx, pool, lookup)
		if (e == nil && m != nil) != want {
			t.Fatal("activation mismatch", want, e)
		}
	}
	enabled(true)
	pairs := [][2]string{
		{`alter table help.article disable trigger help_head_guard`, `alter table help.article enable trigger help_head_guard`},
		{`alter table help.revision disable trigger help_revision_guard`, `alter table help.revision enable trigger help_revision_guard`},
		{`alter table help.article drop constraint help_current_revision_fk`, `alter table help.article add constraint help_current_revision_fk foreign key(tenant_id,article_id,current_version)references help.revision(tenant_id,article_id,version)deferrable initially deferred`},
		{`alter table help.revision drop constraint revision_tenant_id_command_id_key`, `alter table help.revision add constraint revision_tenant_id_command_id_key unique(tenant_id,command_id)`},
		{`alter index help.help_scope_page_idx rename to help_scope_page_saved`, `alter index help.help_scope_page_saved rename to help_scope_page_idx`},
	}
	for _, pair := range pairs {
		func() {
			if _, e := pool.Exec(ctx, pair[0]); e != nil {
				t.Fatal(e)
			}
			defer func() {
				if _, e := pool.Exec(context.Background(), pair[1]); e != nil {
					t.Error(e)
				}
			}()
			enabled(false)
		}()
		enabled(true)
	}
	down, e := os.ReadFile("../../db/migrations/0078_help_cms.down.sql")
	if e != nil {
		t.Fatal(e)
	}
	up, e := os.ReadFile("../../db/migrations/0078_help_cms.up.sql")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, string(down)); e != nil {
		t.Fatal("empty downgrade", e)
	}
	enabled(false)
	if _, e = pool.Exec(ctx, string(up)); e != nil {
		t.Fatal("reapply", e)
	}
	enabled(true)
	tenant := uuid.NewString()
	if _, e = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'help-host','Synthetic','Synthetic')`, tenant); e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'root','root','Synthetic','franchisor')`, tenant); e != nil {
		t.Fatal(e)
	}
	store, _ := postgres.NewHelpCMS(pool)
	p := identity.Principal{TenantID: tenant, Subject: "fixture", Permissions: map[string]struct{}{"*": {}}}
	if _, e = store.Execute(ctx, p, cms.Command{CommandID: "create", Action: "create", ArticleID: "guide", OrganizationID: "root", Locale: "es", Category: "operations", Title: "Guía", Body: "Contenido revisado"}); e != nil {
		t.Fatal(e)
	}
	tx, e := pool.Begin(ctx)
	if e != nil {
		t.Fatal(e)
	}
	defer tx.Rollback(ctx)
	body := strings.TrimSpace(string(down))
	body = strings.TrimPrefix(body, "begin;")
	body = strings.TrimSuffix(body, "commit;")
	if _, e = tx.Exec(ctx, body); e == nil {
		t.Fatal("populated evidence discarded")
	}
	tx.Rollback(ctx)
	enabled(true)
	t.Log("HELP_CMS_HOST_PASS five_missing_guards_block=true empty_downgrade_reapply=true populated_downgrade_refused=true no_credentials=true")
}
