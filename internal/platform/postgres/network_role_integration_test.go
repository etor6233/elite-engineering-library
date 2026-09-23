package postgres_test

// AUTHORED new network/recovery proof. No existing commerce/warranty suites rerun.
import (
	"context"
	nr "elite.local/enterprise/internal/networkrole"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestNetworkRoleAtomic(t *testing.T) {
	raw := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if raw == "" {
		t.Skip("owned database required")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, raw)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant := uuid.NewString()
	if _, e = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'network-role','Synthetic','Synthetic')`, tenant); e != nil {
		t.Fatal(e)
	}
	store, e := db.NewNetworkRole(pool)
	if e != nil {
		t.Fatal(e)
	}
	actor := identity.Principal{TenantID: tenant, Subject: "bootstrap", Permissions: map[string]struct{}{"*": {}}, Organizations: map[string]struct{}{}}
	root := nr.Command{CommandID: "root-create", Action: "create-organization", EntityID: "root", Code: "root", DisplayName: "Synthetic franchisor", Type: "franchisor"}
	var wg sync.WaitGroup
	answers := make(chan nr.Receipt, 4)
	failures := make(chan error, 4)
	start := make(chan struct{})
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); <-start; v, e := store.Execute(ctx, actor, root); answers <- v; failures <- e }()
	}
	close(start)
	wg.Wait()
	close(answers)
	close(failures)
	for e := range failures {
		if e != nil {
			t.Fatal(e)
		}
	}
	fresh, replay := 0, 0
	for v := range answers {
		if v.Entity.ID != "root" || v.Entity.Version != "1" {
			t.Fatal(v)
		}
		if v.Replay {
			replay++
		} else {
			fresh++
		}
	}
	if fresh != 1 || replay != 3 {
		t.Fatal(fresh, replay)
	}
	execute := func(c nr.Command) nr.Receipt {
		t.Helper()
		v, e := store.Execute(ctx, actor, c)
		if e != nil {
			t.Fatal(c.Action, c.EntityID, e)
		}
		if v.Replay {
			t.Fatal("new command replayed")
		}
		return v
	}
	transition := func(kind, scope, id, from, to, version, key string) nr.Receipt {
		return execute(nr.Command{CommandID: key, Action: "transition-" + kind, ScopeOrganizationID: scope, EntityID: id, Current: from, Target: to, Version: version})
	}
	transition("organization", "root", "root", "provisioning", "active", "1", "root-activate")
	execute(nr.Command{CommandID: "franchise-create", Action: "create-organization", ScopeOrganizationID: "root", EntityID: "franchise", Code: "franchise", DisplayName: "Synthetic franchise", Type: "franchisee"})
	transition("organization", "franchise", "franchise", "provisioning", "active", "1", "franchise-activate")
	agreement := nr.Command{CommandID: "agreement-create", Action: "create-agreement", ScopeOrganizationID: "franchise", EntityID: "agreement", TerritoryCode: "SYNTHETIC", TermsVersion: "fixture-v1", StartsOn: "2026-01-01", EndsOn: "2027-01-01"}
	execute(agreement)
	transition("agreement", "franchise", "agreement", "draft", "active", "1", "agreement-active")
	branch := nr.Command{CommandID: "branch-create", Action: "create-organization", ScopeOrganizationID: "franchise", EntityID: "branch", Code: "branch", DisplayName: "Synthetic branch", Type: "store"}
	execute(branch)
	transition("organization", "branch", "branch", "provisioning", "active", "1", "branch-active")
	state := func() string {
		t.Helper()
		var s string
		e := pool.QueryRow(ctx, `select jsonb_build_array((select count(*)from org.organization where tenant_id=$1),(select count(*)from franchise.agreement where tenant_id=$1),(select count(*)from franchise.network_command_receipt where tenant_id=$1),(select count(*)from platform.outbox_event where tenant_id=$1))::text`, tenant).Scan(&s)
		if e != nil {
			t.Fatal(e)
		}
		return s
	}
	if state() != "[3, 1, 8, 8]" {
		t.Fatal(state())
	}
	original := state()
	deny := func(c nr.Command, p identity.Principal) {
		t.Helper()
		if _, e := store.Execute(ctx, p, c); e == nil {
			t.Fatal("unexpected command admitted", c)
		}
		if state() != original {
			t.Fatal("denial leaked durable effect", state(), original)
		}
	}
	changed := root
	changed.DisplayName = "Different payload"
	deny(changed, actor)
	other := actor
	other.Subject = "other-actor"
	deny(root, other)
	if _, e = store.Result(ctx, other, "", "root-create"); !errors.Is(e, nr.ErrNotFound) {
		t.Fatal("foreign actor recovered", e)
	}
	foreign := actor
	foreign.TenantID = uuid.NewString()
	if _, e = store.Result(ctx, foreign, "", "root-create"); !errors.Is(e, nr.ErrNotFound) {
		t.Fatal(e)
	}
	scoped := identity.Principal{TenantID: tenant, Subject: "limited", Permissions: map[string]struct{}{"network:admin": {}}, Organizations: map[string]struct{}{"branch": {}}}
	deny(nr.Command{CommandID: "wrong-scope", Action: "transition-organization", ScopeOrganizationID: "franchise", EntityID: "franchise", Current: "active", Target: "suspended", Version: "2"}, scoped)
	deny(nr.Command{CommandID: "stale", Action: "transition-organization", ScopeOrganizationID: "branch", EntityID: "branch", Current: "provisioning", Target: "active", Version: "1"}, actor)
	deny(nr.Command{CommandID: "premature-close", Action: "transition-organization", ScopeOrganizationID: "franchise", EntityID: "franchise", Current: "active", Target: "closed", Version: "2"}, actor)
	if _, e = pool.Exec(ctx, `create function franchise.network_fixture_fail()returns trigger language plpgsql as $$begin if new.command_id='late-failure'then raise exception 'late receipt failure';end if;return new;end$$;create trigger network_fixture_failure before insert on franchise.network_command_receipt for each row execute function franchise.network_fixture_fail()`); e != nil {
		t.Fatal(e)
	}
	late := branch
	late.CommandID = "late-failure"
	late.EntityID = "late"
	late.Code = "late"
	deny(late, actor)
	if _, e = pool.Exec(ctx, `drop trigger network_fixture_failure on franchise.network_command_receipt;drop function franchise.network_fixture_fail()`); e != nil {
		t.Fatal(e)
	}
	conflict := agreement
	conflict.CommandID = "second-agreement"
	conflict.EntityID = "second"
conflict.StartsOn = "2026-02-01"
	execute(conflict)
	original = state()
	deny(nr.Command{CommandID: "territory-conflict", Action: "transition-agreement", ScopeOrganizationID: "franchise", EntityID: "second", Current: "draft", Target: "active", Version: "1"}, actor)
	if _, e = pool.Exec(ctx, `update franchise.network_command_receipt set actor_subject='changed'where tenant_id=$1`, tenant); e == nil {
		t.Fatal("immutable receipt updated")
	}
	transition("agreement", "franchise", "agreement", "active", "terminated", "2", "agreement-terminate")
	transition("organization", "branch", "branch", "active", "closed", "2", "branch-close")
	transition("organization", "franchise", "franchise", "active", "closed", "2", "franchise-close")
	v, e := store.Result(ctx, actor, "franchise", "branch-create")
	if e != nil || v.Entity.State != "provisioning" || v.Entity.Version != "1" {
		t.Fatal("historical receipt", v, e)
	}
	current, e := store.Entity(ctx, actor, "organization", "branch", "branch")
	if e != nil || current.State != "closed" || current.Version != "3" {
		t.Fatal("current projection", current, e)
	}
	if _, e = store.Entity(ctx, scoped, "agreement", "franchise", "agreement"); !errors.Is(e, nr.ErrNotFound) {
		t.Fatal("view escaped organization", e)
	}
	var p json.RawMessage
	if e = pool.QueryRow(ctx, `select entity_payload from franchise.network_command_receipt where tenant_id=$1 and command_id='agreement-create'`, tenant).Scan(&p); e != nil || !strings.Contains(string(p), "2026-01-01") {
		t.Fatal("explicit dates retained", e)
	}
	t.Logf("NETWORK_ROLE_ATOMIC_PASS concurrency=1new3replay rollback=4tables territory_non_overlap=true no_grants=true historical_GET=true final=%s", state())
}
