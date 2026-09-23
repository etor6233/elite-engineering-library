package postgres_test

// AUTHORED connected immutable CMS proof; synthetic owned database only.
import (
	"context"
	cms "elite.local/enterprise/internal/helpcms"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"sync"
	"testing"
)

func TestHelpCMSAtomic(t *testing.T) {
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
	if _, e = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'help-cms','Synthetic','Synthetic');`, tenant); e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type,status,version)values($1,'root','root','Synthetic','franchisor','active',1)`, tenant); e != nil {
		t.Fatal(e)
	}
	store, e := db.NewHelpCMS(pool)
	if e != nil {
		t.Fatal(e)
	}
	principal := func(subject string, permissions ...string) identity.Principal {
		p := identity.Principal{TenantID: tenant, Subject: subject, Permissions: map[string]struct{}{}, Organizations: map[string]struct{}{"root": {}}}
		for _, v := range permissions {
			p.Permissions[v] = struct{}{}
		}
		return p
	}
	actor := principal("editor", "help:write", "help:publish")
	reader := principal("reader", "help:read")
	c := cms.Command{CommandID: "create", Action: "create", ArticleID: "safety", OrganizationID: "root", Locale: "es", Category: "operations", Title: "Recepción segura", Body: "Verificar serie y preservar evidencia. <script>literal</script>"}
	var wg sync.WaitGroup
	answers := make(chan cms.Receipt, 4)
	errs := make(chan error, 4)
	start := make(chan struct{})
	for range 4 {
		wg.Add(1)
		go func() { defer wg.Done(); <-start; v, e := store.Execute(ctx, actor, c); answers <- v; errs <- e }()
	}
	close(start)
	wg.Wait()
	close(answers)
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	fresh, replay := 0, 0
	for v := range answers {
		if v.Replay {
			replay++
		} else {
			fresh++
		}
	}
	if fresh != 1 || replay != 3 {
		t.Fatal(fresh, replay)
	}
	state := func() string {
		t.Helper()
		var s string
		if e := pool.QueryRow(ctx, `select jsonb_build_array((select count(*)from help.article where tenant_id=$1),(select count(*)from help.revision where tenant_id=$1),(select count(*)from platform.outbox_event where tenant_id=$1))::text`, tenant).Scan(&s); e != nil {
			t.Fatal(e)
		}
		return s
	}
	deny := func(p identity.Principal, c cms.Command) {
		t.Helper()
		before := state()
		if _, e := store.Execute(ctx, p, c); e == nil {
			t.Fatal("command admitted", c)
		}
		if state() != before {
			t.Fatal("denial leaked writes")
		}
	}
	deny(reader, c)
	deny(principal("other", "help:write"), c)
	changed := c
	changed.Body = "changed"
	deny(actor, changed)
	wrong := actor
	wrong.Organizations = map[string]struct{}{"elsewhere": {}}
	deny(wrong, c)
	foreign := reader
	foreign.TenantID = uuid.NewString()
	for _, p := range []identity.Principal{reader, wrong, foreign} {
		if _, e := store.Article(ctx, p, "root", "safety", 0); !errors.Is(e, cms.ErrNotFound) {
			t.Fatal("draft visible", e)
		}
	}
	if _, e := store.Result(ctx, principal("other", "help:write"), "root", "create"); !errors.Is(e, cms.ErrNotFound) {
		t.Fatal("foreign recovery", e)
	}
	updated := cms.Command{CommandID: "update", Action: "update", ArticleID: "safety", OrganizationID: "root", Version: "1", Title: "Recepción revisada", Body: "Respetar el serial y la aprobación humana."}
	if v, e := store.Execute(ctx, actor, updated); e != nil || v.Article.Version != "2" {
		t.Fatal(v, e)
	}
	stale := updated
	stale.CommandID = "stale"
	deny(actor, stale)
	// Inject a failure at the final durable effect, after both article and revision inserts.
	if _, e = pool.Exec(ctx, `create function help.fixture_fail()returns trigger language plpgsql as $$begin if new.aggregate_id='late'then raise exception 'late outbox failure';end if;return new;end$$;create trigger help_fixture_failure before insert on platform.outbox_event for each row execute function help.fixture_fail()`); e != nil {
		t.Fatal(e)
	}
	late := c
	late.CommandID = "late"
	late.ArticleID = "late"
	deny(actor, late)
	if _, e = pool.Exec(ctx, `drop trigger help_fixture_failure on platform.outbox_event;drop function help.fixture_fail()`); e != nil {
		t.Fatal(e)
	}
	pub := cms.Command{CommandID: "publish", Action: "publish", ArticleID: "safety", OrganizationID: "root", Version: "2"}
	deny(principal("writer", "help:write"), pub)
	if v, e := store.Execute(ctx, actor, pub); e != nil || v.Article.Version != "3" {
		t.Fatal(v, e)
	}
	if a, e := store.Article(ctx, reader, "root", "safety", 0); e != nil || a.Title != updated.Title || a.Body != updated.Body {
		t.Fatal(a, e)
	}
	if _, e := store.Article(ctx, reader, "root", "safety", 1); !errors.Is(e, cms.ErrNotFound) {
		t.Fatal("draft history leaked", e)
	}
	if page, e := store.List(ctx, reader, "root", "es", "APROBACIÓN", ""); e != nil || len(page.Items) != 1 {
		t.Fatal(page, e)
	}
	if page, e := store.List(ctx, reader, "root", "en", "", ""); e != nil || len(page.Items) != 0 {
		t.Fatal("locale leaked", page, e)
	}
	if _, e := store.History(ctx, reader, "root", "safety", 0); !errors.Is(e, cms.ErrNotFound) {
		t.Fatal("reader history", e)
	}
	editedPublished := updated
	editedPublished.CommandID = "bad-edit"
	editedPublished.Version = "3"
	deny(actor, editedPublished)
	if _, e = pool.Exec(ctx, `update help.revision set body='tamper'where tenant_id=$1`, tenant); e == nil {
		t.Fatal("revision not immutable")
	}
	if _, e = pool.Exec(ctx, `delete from help.revision where tenant_id=$1`, tenant); e == nil {
		t.Fatal("evidence deleted")
	}
	if _, e = pool.Exec(ctx, `update help.article set locale='en'where tenant_id=$1`, tenant); e == nil {
		t.Fatal("scope mutable")
	}
	archive := pub
	archive.CommandID = "archive"
	archive.Action = "archive"
	archive.Version = "3"
	if v, e := store.Execute(ctx, actor, archive); e != nil || v.Article.Version != "4" {
		t.Fatal(v, e)
	}
	for _, version := range []int64{0, 3} {
		if _, e := store.Article(ctx, reader, "root", "safety", version); !errors.Is(e, cms.ErrNotFound) {
			t.Fatal("withdrawn published body visible", e)
		}
	}
	if page, e := store.List(ctx, reader, "root", "es", "", ""); e != nil || len(page.Items) != 0 {
		t.Fatal(page, e)
	}
	if page, e := store.History(ctx, actor, "root", "safety", 0); e != nil || len(page.Items) != 4 || page.Items[0].Version != "4" {
		t.Fatal(page, e)
	}
	if v, e := store.Result(ctx, actor, "root", "create"); e != nil || v.Article.Version != "1" || !v.Replay {
		t.Fatal("original recovery changed", v, e)
	}
	if state() != "[1, 4, 4]" {
		t.Fatal(state())
	}
	// Keyset pagination includes only the selected organization/locale, without duplicates.
	for i := range 51 {
		v := c
		v.CommandID = fmt.Sprintf("page-%02d", i)
		v.ArticleID = v.CommandID
		if _, e := store.Execute(ctx, actor, v); e != nil {
			t.Fatal(e)
		}
	}
	page, e := store.List(ctx, actor, "root", "es", "", "")
	if e != nil || len(page.Items) != 50 || page.Next == "" {
		t.Fatal(page, e)
	}
	rest, e := store.List(ctx, actor, "root", "es", "", page.Next)
	if e != nil || len(rest.Items) != 2 || rest.Next != "" {
		t.Fatal(rest, e)
	}
	seen := map[string]bool{}
	for _, v := range append(page.Items, rest.Items...) {
		if seen[v.ID] {
			t.Fatal("duplicate page")
		}
		seen[v.ID] = true
	}
	var grants int
	if e = pool.QueryRow(ctx, `select count(*)from platform.outbox_event where tenant_id=$1 and event_type not like 'help.article-%'`, tenant).Scan(&grants); e != nil || grants != 0 {
		t.Fatal("CMS side effects outside article", e, grants)
	}
	rawReceipt, e := json.Marshal(c)
	if e != nil || len(rawReceipt) == 0 {
		t.Fatal(e)
	}
	t.Logf("HELP_CMS_ATOMIC_PASS concurrency=1new3replay rollback=3tables immutable=true scope=true publication_withdrawal=true pagination=52 no_grants=true final=%s", state())
}
