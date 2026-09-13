# Connected durable help CMS

## 1. Metadata

```yaml
pack_id: "GO-CONNECTED-HELP-CMS"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Organization-scoped durable help authoring/publication/history and recovery, with six operation guides shared by UI/help and explicitly versioned human training."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["MARKDOWN-COMPOSITOR0.3.0", "V402 existing owner closure"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-12"
```

## 2. Applicability

LIBRARY_INFRASTRUCTURE T2804: existing GO-HELP-CENTER-CORE pure article lifecycle, PG/identity/outbox/approval owners. Public operation guides remain release bound; private CMS content requires explicit reviewed curriculum incorporation.

## 3. Architecture contract

Organization/locale immutable; draft edits only, publish/archive preserve immutable version and actor/command/hash receipt plus outbox in one transaction. Archive withdraws prior published versions from ordinary readers. Explicit current read before another UI write; lost response GET recovery. No public CMS or automatic grant.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/help_cms.go
CREATE cmd/electromobility-api/help_cms_test.go
CREATE db/migrations/0078_help_cms.down.sql
CREATE db/migrations/0078_help_cms.up.sql
CREATE internal/helpcms/contract.go
CREATE internal/helpcms/contract_test.go
CREATE internal/platform/httpapi/help_cms.go
CREATE internal/platform/httpapi/help_cms_test.go
CREATE internal/platform/postgres/help_cms.go
CREATE internal/platform/postgres/help_cms_browser_integration_test.go
CREATE internal/platform/postgres/help_cms_integration_test.go
CREATE internal/trainingbridge/release_guides_connected_test.go
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/help_cms.go`

```yaml
block_id: "GO-CONNECTED-HELP-CMS:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8ca0f4f570b076d849c5f02ef5eed76ce9e396c5e28d922f2b170351ba3da78c"
variables: []
secrets_allowed: false
```

````go
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
````

### FILE: `cmd/electromobility-api/help_cms_test.go`

```yaml
block_id: "GO-CONNECTED-HELP-CMS:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "0bc8151608d7106b566b9bd90aa25a4803e2328e8c2270de55bc6059fb74a67b"
variables: []
secrets_allowed: false
```

````go
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
````

### FILE: `db/migrations/0078_help_cms.down.sql`

```yaml
block_id: "GO-CONNECTED-HELP-CMS:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "290a8e1f14d759f8b5e57122103446411e73d7942ad1128eefb35ac3df513485"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$begin if exists(select 1 from help.revision)then raise exception 'cannot discard help content history';end if;end$$;
alter table help.article drop constraint help_current_revision_fk;
drop table help.revision;
drop table help.article;
drop function help.guard_revision();
drop function help.guard_head();
drop schema help;
commit;
````

### FILE: `db/migrations/0078_help_cms.up.sql`

```yaml
block_id: "GO-CONNECTED-HELP-CMS:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "25ad123d9038e535405849dccb8771a95441f6c48ca42b974969d4965cc188c6"
variables: []
secrets_allowed: false
```

````sql
begin;
create schema help;
create table help.article(
 tenant_id uuid not null,
 article_id text not null check(length(article_id)between 1 and 128),
 organization_id text not null,
 locale text not null check(locale in('es','en')),
 category text not null check(category~'^[a-z][a-z0-9_-]{0,63}$'),
 current_version bigint not null check(current_version>0),
 primary key(tenant_id,article_id),
 foreign key(tenant_id,organization_id)references org.organization(tenant_id,organization_id)
);
create table help.revision(
 tenant_id uuid not null,
 article_id text not null,
 version bigint not null check(version>0),
 command_id text not null check(length(command_id)between 1 and 128),
 action text not null check(action in('create','update','publish','archive')),
 actor_subject text not null check(length(actor_subject)between 1 and 256),
 request_sha256 text not null check(request_sha256~'^[a-f0-9]{64}$'),
 article_sha256 text not null check(article_sha256~'^[a-f0-9]{64}$'),
 snapshot jsonb not null check(jsonb_typeof(snapshot)='object'),
 title text not null check(octet_length(title)between 1 and 200),
 body text not null check(octet_length(body)between 1 and 16384),
 state text not null check(state in('draft','published','archived')),
 primary key(tenant_id,article_id,version),
 unique(tenant_id,command_id),
 foreign key(tenant_id,article_id)references help.article(tenant_id,article_id)
);
alter table help.article add constraint help_current_revision_fk foreign key(tenant_id,article_id,current_version)references help.revision(tenant_id,article_id,version)deferrable initially deferred;
create function help.guard_revision()returns trigger language plpgsql as $$begin
 if tg_op<>'INSERT'then raise exception 'help revision is immutable';end if;
 if new.snapshot->>'id' is distinct from new.article_id or new.snapshot->>'version' is distinct from new.version::text
 or new.snapshot->>'state' is distinct from new.state or new.snapshot->>'title' is distinct from new.title or new.snapshot->>'body' is distinct from new.body
 or not exists(select 1 from help.article h where h.tenant_id=new.tenant_id and h.article_id=new.article_id
 and h.organization_id=new.snapshot->>'organization_id'and h.locale=new.snapshot->>'locale'and h.category=new.snapshot->>'category')
 then raise exception 'help snapshot differs from head';end if;return new;
end$$;
create trigger help_revision_guard before insert or update or delete on help.revision for each row execute function help.guard_revision();
create function help.guard_head()returns trigger language plpgsql as $$begin
 if tg_op='DELETE'then raise exception 'help head evidence cannot be deleted';end if;
 if row(new.tenant_id,new.article_id,new.organization_id,new.locale,new.category)is distinct from row(old.tenant_id,old.article_id,old.organization_id,old.locale,old.category)
 or new.current_version<>old.current_version+1 then raise exception 'help head scope or version changed';end if;return new;
end$$;
create trigger help_head_guard before update or delete on help.article for each row execute function help.guard_head();
create index help_scope_page_idx on help.article(tenant_id,organization_id,locale,article_id);
commit;
````

### FILE: `internal/helpcms/contract.go`

```yaml
block_id: "GO-CONNECTED-HELP-CMS:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b9bca3646f504883761852a20f2099edfb68268181f59641e65a3999d19b3156"
variables: []
secrets_allowed: false
```

````go
// AUTHORED bounded persistence/authorization contract; article rules live in helpcenter.
package helpcms

import (
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var ErrInvalid = errors.New("invalid help CMS command")
var ErrNotFound = errors.New("help article unavailable")
var ErrConflict = errors.New("help CMS command conflict")
var idPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
var categoryPattern = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)

const MaxBodyBytes = 16384

func ID(s string) bool { return idPattern.MatchString(s) }
func Text(s string, max int, multiline bool) bool {
	if strings.TrimSpace(s) == "" || len(s) > max || !utf8.ValidString(s) {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) && (!multiline || r != '\n' && r != '\r' && r != '\t') {
			return false
		}
	}
	return true
}

type Command struct {
	CommandID      string `json:"command_id"`
	Action         string `json:"action"`
	ArticleID      string `json:"article_id"`
	OrganizationID string `json:"organization_id"`
	Version        string `json:"version,omitempty"`
	Locale         string `json:"locale,omitempty"`
	Category       string `json:"category,omitempty"`
	Title          string `json:"title,omitempty"`
	Body           string `json:"body,omitempty"`
}

func (c Command) Valid() bool {
	raw, e := json.Marshal(c)
	if e != nil || len(raw) > 24576 {
		return false
	}
	if !ID(c.CommandID) || !ID(c.ArticleID) || !ID(c.OrganizationID) {
		return false
	}
	switch c.Action {
	case "create":
		return c.Version == "" && (c.Locale == "es" || c.Locale == "en") && categoryPattern.MatchString(c.Category) && Text(c.Title, 200, false) && Text(c.Body, MaxBodyBytes, true)
	case "update", "publish", "archive":
		v, e := strconv.ParseInt(c.Version, 10, 64)
		if e != nil || v < 1 || v >= math.MaxInt64 || strconv.FormatInt(v, 10) != c.Version || c.Locale != "" || c.Category != "" {
			return false
		}
		if c.Action == "update" {
			return Text(c.Title, 200, false) && Text(c.Body, MaxBodyBytes, true)
		}
		return c.Title == "" && c.Body == ""
	}
	return false
}
func Permission(action string) string {
	switch action {
	case "create", "update":
		return "help:write"
	case "publish", "archive":
		return "help:publish"
	}
	return ""
}
func Authorized(p identity.Principal, permission, org string) bool {
	return permission != "" && p.TenantID != "" && p.Subject != "" && ID(org) && p.Allowed(permission) && p.AllowedOrganization(org)
}
func Editor(p identity.Principal, org string) bool {
	return Authorized(p, "help:write", org) || Authorized(p, "help:publish", org)
}
func Reader(p identity.Principal, org string) bool {
	return Editor(p, org) || Authorized(p, "help:read", org)
}
func Canonical(v any) (json.RawMessage, string, error) {
	raw, e := json.Marshal(v)
	if e != nil {
		return nil, "", ErrInvalid
	}
	body, hash, e := approval.CanonicalPayload(raw)
	if e != nil {
		return nil, "", ErrInvalid
	}
	return body, hash, nil
}

type Article struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Locale         string    `json:"locale"`
	Category       string    `json:"category"`
	Title          string    `json:"title"`
	Body           string    `json:"body"`
	State          string    `json:"state"`
	Version        string    `json:"version"`
	UpdatedAt      time.Time `json:"updated_at"`
}
type Receipt struct {
	CommandID     string  `json:"command_id"`
	Action        string  `json:"action"`
	Actor         string  `json:"actor"`
	RequestSHA256 string  `json:"request_sha256"`
	ArticleSHA256 string  `json:"article_sha256"`
	Article       Article `json:"article"`
	Replay        bool    `json:"replay"`
}
type Summary struct {
	ID       string `json:"id"`
	Locale   string `json:"locale"`
	Category string `json:"category"`
	Title    string `json:"title"`
	State    string `json:"state"`
	Version  string `json:"version"`
}
type Page struct {
	Items []Summary `json:"items"`
	Next  string    `json:"next,omitempty"`
}
````

### FILE: `internal/helpcms/contract_test.go`

```yaml
block_id: "GO-CONNECTED-HELP-CMS:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "fd4852fb2a98f81bfe089ed45ac509f54dbf093faadcadf15d8351e852e8800a"
variables: []
secrets_allowed: false
```

````go
package helpcms

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestHelpCMSWireGoldens(t *testing.T) {
	raw, e := os.ReadFile("../../deploy/help/cms-goldens.json")
	if e != nil {
		t.Fatal(e)
	}
	var rows []struct {
		Command   Command `json:"command"`
		Canonical string  `json:"canonical"`
		SHA       string  `json:"sha256"`
	}
	if json.Unmarshal(raw, &rows) != nil || len(rows) != 4 {
		t.Fatal("four typed goldens required")
	}
	for _, r := range rows {
		body, hash, e := Canonical(r.Command)
		if !r.Command.Valid() || e != nil || string(body) != r.Canonical || hash != r.SHA {
			t.Fatal(r.Command.Action, e)
		}
	}
	for _, s := range []string{"0", "01", "-1", "9223372036854775807", "9223372036854775808"} {
		c := rows[1].Command
		c.Version = s
		if c.Valid() {
			t.Fatal("unsafe version", s)
		}
	}
	for _, body := range []string{strings.Repeat("a", 16385), strings.Repeat("<", 5000), "bad\x00body", string([]byte{0xff})} {
		c := rows[0].Command
		c.Body = body
		if c.Valid() {
			t.Fatal("unbounded/invalid body")
		}
	}
}
func FuzzHelpCMSWire(f *testing.F) {
	for _, s := range []string{`{"command_id":"create","action":"create","article_id":"guide","organization_id":"org","locale":"es","category":"operations","title":"Guía","body":"Contenido"}`, `{"command_id":"publish","action":"publish","article_id":"guide","organization_id":"org","version":"9007199254740993"}`, `{}`} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 32768 {
			return
		}
		var c Command
		if json.Unmarshal(raw, &c) != nil || !c.Valid() {
			return
		}
		body, hash, e := Canonical(c)
		if e != nil || len(hash) != 64 {
			t.Fatal("valid command cannot canonicalize", e)
		}
		var again Command
		if json.Unmarshal(body, &again) != nil || again != c || !again.Valid() {
			t.Fatal("canonical round trip")
		}
	})
}
````

### FILE: `internal/platform/httpapi/help_cms.go`

```yaml
block_id: "GO-CONNECTED-HELP-CMS:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c25aed65abde98826f17fd987f5e368edaed4ae13fe8ecd4f2a181db0b4cb48c"
variables: []
secrets_allowed: false
```

````go
// AUTHORED bounded transport; persisted article rules belong to helpcenter/helpcms.
package httpapi

import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/approval"
	cms "elite.local/enterprise/internal/helpcms"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type HelpCMSService interface {
	Execute(context.Context, identity.Principal, cms.Command) (cms.Receipt, error)
	Result(context.Context, identity.Principal, string, string) (cms.Receipt, error)
	Article(context.Context, identity.Principal, string, string, int64) (cms.Article, error)
	List(context.Context, identity.Principal, string, string, string, string) (cms.Page, error)
	History(context.Context, identity.Principal, string, string, int64) (cms.Page, error)
}
type HelpCMSModule struct{ Service HelpCMSService }

func helpCMSReply(w http.ResponseWriter, v any, e error) {
	if e != nil {
		switch {
		case errors.Is(e, cms.ErrInvalid):
			writeProblem(w, 400, "HELP_INVALID", "invalid article request")
		case errors.Is(e, cms.ErrNotFound):
			writeProblem(w, 404, "HELP_UNAVAILABLE", "article unavailable in this scope")
		case errors.Is(e, cms.ErrConflict):
			writeProblem(w, 409, "HELP_CONFLICT", "consult saved result and current version")
		default:
			writeProblem(w, 503, "HELP_UNCONFIRMED", "result unconfirmed; consult saved command")
		}
		return
	}
	writeJSON(w, 200, v)
}
func helpCMSQuery(r *http.Request, allowed ...string) (url.Values, bool) {
	if len(r.URL.RawQuery) > 2048 {
		return nil, false
	}
	q, e := url.ParseQuery(r.URL.RawQuery)
	if e != nil {
		return nil, false
	}
	if len(q["organization_id"]) != 1 || !cms.ID(q.Get("organization_id")) {
		return nil, false
	}
	for k, v := range q {
		found := k == "organization_id"
		for _, a := range allowed {
			found = found || k == a
		}
		if !found || len(v) != 1 {
			return nil, false
		}
	}
	return q, true
}
func helpCMSVersion(s string) (int64, bool) {
	if s == "" {
		return 0, true
	}
	v, e := strconv.ParseInt(s, 10, 64)
	return v, e == nil && v > 0 && strconv.FormatInt(v, 10) == s
}
func (m HelpCMSModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	auth := func(w http.ResponseWriter, r *http.Request, write bool) (identity.Principal, bool) {
		w.Header().Set("Cache-Control", "private, no-store")
		w.Header().Set("Vary", "Authorization")
		p, e := authenticate(r.Context(), r.Header.Get("Authorization"), verifier)
		if e != nil {
			writeProblem(w, 401, "UNAUTHENTICATED", "valid session required")
			return p, false
		}
		if !p.Allowed("help:write") && !p.Allowed("help:publish") && (write || !p.Allowed("help:read")) {
			writeProblem(w, 403, "FORBIDDEN", "help permission required")
			return p, false
		}
		return p, true
	}
	mux.HandleFunc("POST /v1/help/cms/commands", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r, true)
		if !ok {
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			writeProblem(w, 415, "CONTENT_TYPE", "application/json required")
			return
		}
		raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
		if e != nil {
			writeProblem(w, 413, "BODY_LIMIT", "article command exceeds32KiB")
			return
		}
		if _, _, e = approval.CanonicalPayload(raw); e != nil {
			helpCMSReply(w, nil, cms.ErrInvalid)
			return
		}
		var c cms.Command
		decoder := json.NewDecoder(bytes.NewReader(raw))
		decoder.DisallowUnknownFields()
		if decoder.Decode(&c) != nil || decoder.Decode(new(any)) != io.EOF || !c.Valid() {
			helpCMSReply(w, nil, cms.ErrInvalid)
			return
		}
		if !cms.Authorized(p, cms.Permission(c.Action), c.OrganizationID) {
			writeProblem(w, 403, "HELP_FORBIDDEN", "action or organization unavailable")
			return
		}
		v, e := m.Service.Execute(r.Context(), p, c)
		helpCMSReply(w, v, e)
	})
	mux.HandleFunc("GET /v1/help/cms/commands/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r, true)
		if !ok {
			return
		}
		q, ok := helpCMSQuery(r)
		if !ok || !cms.ID(r.PathValue("id")) {
			helpCMSReply(w, nil, cms.ErrInvalid)
			return
		}
		v, e := m.Service.Result(r.Context(), p, q.Get("organization_id"), r.PathValue("id"))
		helpCMSReply(w, v, e)
	})
	mux.HandleFunc("GET /v1/help/cms/articles", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r, false)
		if !ok {
			return
		}
		q, ok := helpCMSQuery(r, "locale", "q", "after")
		if !ok {
			helpCMSReply(w, nil, cms.ErrInvalid)
			return
		}
		v, e := m.Service.List(r.Context(), p, q.Get("organization_id"), q.Get("locale"), q.Get("q"), q.Get("after"))
		helpCMSReply(w, v, e)
	})
	mux.HandleFunc("GET /v1/help/cms/articles/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r, false)
		if !ok {
			return
		}
		q, ok := helpCMSQuery(r, "version")
		version, valid := helpCMSVersion(q.Get("version"))
		if !ok || !valid || !cms.ID(r.PathValue("id")) {
			helpCMSReply(w, nil, cms.ErrInvalid)
			return
		}
		v, e := m.Service.Article(r.Context(), p, q.Get("organization_id"), r.PathValue("id"), version)
		helpCMSReply(w, v, e)
	})
	mux.HandleFunc("GET /v1/help/cms/articles/{id}/history", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r, true)
		if !ok {
			return
		}
		q, ok := helpCMSQuery(r, "before")
		before, valid := helpCMSVersion(q.Get("before"))
		if !ok || !valid || !cms.ID(r.PathValue("id")) {
			helpCMSReply(w, nil, cms.ErrInvalid)
			return
		}
		v, e := m.Service.History(r.Context(), p, q.Get("organization_id"), r.PathValue("id"), before)
		helpCMSReply(w, v, e)
	})
}
````

### FILE: `internal/platform/httpapi/help_cms_test.go`

```yaml
block_id: "GO-CONNECTED-HELP-CMS:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "bda72999eb9b21fc12ffb6e895186d044c935dc17b0b7d3e3832ce47f74639f4"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type cmsVerifier struct{ p identity.Principal }

func (v cmsVerifier) Verify(context.Context, string) (identity.Principal, error) { return v.p, nil }

type cmsUnread struct{ reads int }

func (v *cmsUnread) Read([]byte) (int, error) { v.reads++; return 0, io.EOF }
func (v *cmsUnread) Close() error             { return nil }
func TestHelpCMSHTTPBoundaries(t *testing.T) {
	for _, permissions := range []map[string]struct{}{{"help:read": {}}, {}} {
		mux := http.NewServeMux()
		HelpCMSModule{}.Register(mux, cmsVerifier{identity.Principal{TenantID: "tenant", Subject: "reader", Permissions: permissions}})
		body := &cmsUnread{}
		r := httptest.NewRequest("POST", "/v1/help/cms/commands", body)
		r.Header.Set("Authorization", "Bearer fixture")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != 403 || body.reads != 0 {
			t.Fatal("unauthorized body consumed", w.Code, body.reads)
		}
	}
	p := identity.Principal{TenantID: "tenant", Subject: "writer", Permissions: map[string]struct{}{"help:write": {}}, Organizations: map[string]struct{}{"org": {}}}
	cases := []struct {
		method, path, body string
		status             int
	}{
		{"POST", "/v1/help/cms/commands", strings.Repeat("x", 32769), 413},
		{"POST", "/v1/help/cms/commands", `{"action":"create","action":"archive"}`, 400},
		{"POST", "/v1/help/cms/commands", `{"command_id":"x","action":"publish","article_id":"guide","organization_id":"org","version":"1"}`, 403},
		{"POST", "/v1/help/cms/commands", `{"command_id":"x","action":"update","article_id":"guide","organization_id":"org","version":1,"title":"Title","body":"Body"}`, 400},
		{"GET", "/v1/help/cms/articles/guide?organization_id=org&organization_id=other", "", 400},
		{"GET", "/v1/help/cms/articles/guide?organization_id=org&version=01", "", 400},
		{"GET", "/v1/help/cms/articles/guide?organization_id=org&version=9223372036854775808", "", 400},
	}
	for _, c := range cases {
		mux := http.NewServeMux()
		HelpCMSModule{}.Register(mux, cmsVerifier{p})
		r := httptest.NewRequest(c.method, c.path, strings.NewReader(c.body))
		r.Header.Set("Authorization", "Bearer fixture")
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != c.status {
			t.Fatal(c.path, w.Code, c.status, w.Body.String())
		}
		if !strings.Contains(w.Header().Get("Cache-Control"), "no-store") {
			t.Fatal("private content cached")
		}
	}
}
````

### FILE: `internal/platform/postgres/help_cms.go`

```yaml
block_id: "GO-CONNECTED-HELP-CMS:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "31e0375fa7128a24ba3ac36a1ad725a33bb6ce8476967d304a65068e73f1a721"
variables: []
secrets_allowed: false
```

````go
// AUTHORED durable/authenticated adapter around the existing helpcenter article kernel.
package postgres

import (
	"context"
	hc "elite.local/enterprise/internal/helpcenter"
	cms "elite.local/enterprise/internal/helpcms"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"
	"time"
)

type HelpCMS struct{ pool *pgxpool.Pool }

func NewHelpCMS(pool *pgxpool.Pool) (*HelpCMS, error) {
	if pool == nil {
		return nil, cms.ErrInvalid
	}
	return &HelpCMS{pool}, nil
}

type cmsReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}
type cmsHead struct {
	org, locale, category string
	version               int64
}

func cmsError(e error) error {
	if errors.Is(e, pgx.ErrNoRows) {
		return cms.ErrNotFound
	}
	if errors.Is(e, hc.ErrInvalidArticle) {
		return cms.ErrInvalid
	}
	for _, v := range []error{hc.ErrDuplicate, hc.ErrNotFound, hc.ErrVersion, hc.ErrBadTransition, hc.ErrVersionExhausted} {
		if errors.Is(e, v) {
			return cms.ErrConflict
		}
	}
	return e
}
func readCMSHead(ctx context.Context, q cmsReader, tenant, id string, locked bool) (cmsHead, error) {
	var h cmsHead
	sql := `select organization_id,locale,category,current_version from help.article where tenant_id=$1 and article_id=$2`
	if locked {
		sql += " for update"
	}
	e := q.QueryRow(ctx, sql, tenant, id).Scan(&h.org, &h.locale, &h.category, &h.version)
	return h, cmsError(e)
}
func readCMSReceipt(ctx context.Context, q cmsReader, tenant, id string, version int64, command string) (cms.Receipt, error) {
	var v cms.Receipt
	var raw []byte
	sql := `select command_id,action,actor_subject,request_sha256,article_sha256,snapshot from help.revision where tenant_id=$1 and article_id=$2 and version=$3`
	args := []any{tenant, id, version}
	if command != "" {
		sql = `select command_id,action,actor_subject,request_sha256,article_sha256,snapshot from help.revision where tenant_id=$1 and command_id=$2`
		args = []any{tenant, command}
	}
	e := q.QueryRow(ctx, sql, args...).Scan(&v.CommandID, &v.Action, &v.Actor, &v.RequestSHA256, &v.ArticleSHA256, &raw)
	if e != nil {
		return v, cmsError(e)
	}
	if json.Unmarshal(raw, &v.Article) != nil {
		return v, cms.ErrConflict
	}
	_, hash, e := cms.Canonical(v.Article)
	if e != nil || hash != v.ArticleSHA256 {
		return v, cms.ErrConflict
	}
	return v, nil
}
func cmsScope(h cmsHead, a cms.Article, id, org string) bool {
	return h.org == org && a.OrganizationID == org && a.ID == id && a.Locale == h.locale && a.Category == h.category
}
func (s *HelpCMS) Result(ctx context.Context, p identity.Principal, org, command string) (cms.Receipt, error) {
	if !cms.ID(command) || !cms.Editor(p, org) {
		return cms.Receipt{}, cms.ErrNotFound
	}
	v, e := readCMSReceipt(ctx, s.pool, p.TenantID, "", 0, command)
	if e != nil {
		return v, e
	}
	if v.Article.OrganizationID != org || v.Actor != p.Subject || !cms.Authorized(p, cms.Permission(v.Action), org) {
		return cms.Receipt{}, cms.ErrNotFound
	}
	v.Replay = true
	return v, nil
}
func (s *HelpCMS) Article(ctx context.Context, p identity.Principal, org, id string, version int64) (cms.Article, error) {
	var empty cms.Article
	if !cms.ID(id) || version < 0 || !cms.Reader(p, org) {
		return empty, cms.ErrNotFound
	}
	tx, e := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if e != nil {
		return empty, e
	}
	defer tx.Rollback(ctx)
	h, e := readCMSHead(ctx, tx, p.TenantID, id, false)
	if e != nil {
		return empty, e
	}
	if h.org != org {
		return empty, cms.ErrNotFound
	}
	current, e := readCMSReceipt(ctx, tx, p.TenantID, id, h.version, "")
	if e != nil {
		return empty, e
	}
	if !cmsScope(h, current.Article, id, org) {
		return empty, cms.ErrConflict
	}
	editor := cms.Editor(p, org)
	if !editor && current.Article.State != "published" {
		return empty, cms.ErrNotFound
	}
	selected := current
	if version != 0 && version != h.version {
		selected, e = readCMSReceipt(ctx, tx, p.TenantID, id, version, "")
		if e != nil {
			return empty, e
		}
	}
	if !cmsScope(h, selected.Article, id, org) {
		return empty, cms.ErrConflict
	}
	if !editor && selected.Article.State != "published" {
		return empty, cms.ErrNotFound
	}
	return selected.Article, tx.Commit(ctx)
}
func (s *HelpCMS) Execute(ctx context.Context, p identity.Principal, c cms.Command) (cms.Receipt, error) {
	var empty cms.Receipt
	if !c.Valid() {
		return empty, cms.ErrInvalid
	}
	if !cms.Authorized(p, cms.Permission(c.Action), c.OrganizationID) {
		return empty, cms.ErrNotFound
	}
	_, hash, e := cms.Canonical(c)
	if e != nil {
		return empty, e
	}
	tx, e := s.pool.Begin(ctx)
	if e != nil {
		return empty, e
	}
	defer tx.Rollback(ctx)
	if _, e = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended('help-command:'||$1::text||':'||$2::text,0))`, p.TenantID, c.CommandID); e != nil {
		return empty, e
	}
	old, e := readCMSReceipt(ctx, tx, p.TenantID, "", 0, c.CommandID)
	if e == nil {
		if old.Actor != p.Subject || old.Action != c.Action || old.Article.OrganizationID != c.OrganizationID || old.Article.ID != c.ArticleID || old.RequestSHA256 != hash {
			return empty, cms.ErrConflict
		}
		old.Replay = true
		return old, tx.Commit(ctx)
	}
	if !errors.Is(e, cms.ErrNotFound) {
		return empty, e
	}
	if _, e = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended('help-article:'||$1::text||':'||$2::text,0))`, p.TenantID, c.ArticleID); e != nil {
		return empty, e
	}
	h, e := readCMSHead(ctx, tx, p.TenantID, c.ArticleID, true)
	if e != nil && !errors.Is(e, cms.ErrNotFound) {
		return empty, e
	}
	exists := e == nil
	var article hc.Article
	if c.Action == "create" {
		if exists {
			return empty, cms.ErrConflict
		}
		h = cmsHead{org: c.OrganizationID, locale: c.Locale, category: c.Category, version: 1}
		article = hc.Article{TenantID: p.TenantID, ID: c.ArticleID, Category: c.Category, Title: c.Title, Body: c.Body}
	} else {
		if !exists || h.org != c.OrganizationID {
			return empty, cms.ErrNotFound
		}
		saved, e := readCMSReceipt(ctx, tx, p.TenantID, c.ArticleID, h.version, "")
		if e != nil {
			return empty, e
		}
		if !cmsScope(h, saved.Article, c.ArticleID, c.OrganizationID) {
			return empty, cms.ErrConflict
		}
		v := saved.Article
		article = hc.Article{TenantID: p.TenantID, ID: c.ArticleID, Category: v.Category, Title: v.Title, Body: v.Body, State: hc.State(v.State), Version: h.version, UpdatedAt: v.UpdatedAt}
	}
	var now time.Time
	if e = tx.QueryRow(ctx, "select clock_timestamp()").Scan(&now); e != nil {
		return empty, e
	}
	expected, _ := strconv.ParseInt(c.Version, 10, 64)
	var next hc.Article
	switch c.Action {
	case "create":
		next, e = hc.NewDraft(article, now)
	case "update":
		// Narrow CMS profile: published/archived snapshots are never edited.
		if article.State != hc.StateDraft {
			return empty, cms.ErrConflict
		}
		next, e = hc.ReviseArticle(article, c.Title, c.Body, expected, now)
	case "publish":
		next, e = hc.PublishArticle(article, expected, now)
	case "archive":
		next, e = hc.ArchiveArticle(article, expected, now)
	}
	if e != nil {
		return empty, cmsError(e)
	}
	view := cms.Article{ID: next.ID, OrganizationID: h.org, Locale: h.locale, Category: next.Category, Title: next.Title, Body: next.Body, State: string(next.State), Version: strconv.FormatInt(next.Version, 10), UpdatedAt: next.UpdatedAt}
	payload, payloadHash, e := cms.Canonical(view)
	if e != nil {
		return empty, e
	}
	if !exists {
		if _, e = tx.Exec(ctx, `insert into help.article(tenant_id,article_id,organization_id,locale,category,current_version)values($1,$2,$3,$4,$5,1)`, p.TenantID, c.ArticleID, h.org, h.locale, h.category); e != nil {
			return empty, e
		}
	}
	_, e = tx.Exec(ctx, `insert into help.revision(tenant_id,article_id,version,command_id,action,actor_subject,request_sha256,article_sha256,snapshot,title,body,state)values($1,$2,$3,$4,$5,$6,$7,$8,$9::jsonb,$10,$11,$12)`, p.TenantID, c.ArticleID, next.Version, c.CommandID, c.Action, p.Subject, hash, payloadHash, string(payload), view.Title, view.Body, view.State)
	if e != nil {
		return empty, e
	}
	if exists {
		tag, e := tx.Exec(ctx, `update help.article set current_version=$3 where tenant_id=$1 and article_id=$2 and current_version=$4`, p.TenantID, c.ArticleID, next.Version, h.version)
		if e != nil {
			return empty, e
		}
		if tag.RowsAffected() != 1 {
			return empty, cms.ErrConflict
		}
	}
	if e = outbox(ctx, tx, p.TenantID, randomid.Generator{}.New(), "help-article", c.ArticleID, next.Version, "help.article-"+c.Action, `jsonb_build_object('organization_id',$7::text,'article_sha256',$8::text)`, h.org, payloadHash); e != nil {
		return empty, e
	}
	receipt := cms.Receipt{CommandID: c.CommandID, Action: c.Action, Actor: p.Subject, RequestSHA256: hash, ArticleSHA256: payloadHash, Article: view}
	return receipt, tx.Commit(ctx)
}
func (s *HelpCMS) List(ctx context.Context, p identity.Principal, org, locale, query, after string) (cms.Page, error) {
	out := cms.Page{Items: []cms.Summary{}}
	if !cms.Reader(p, org) || (locale != "es" && locale != "en") || (after != "" && !cms.ID(after)) || len(query) > 160 || (query != "" && !cms.Text(query, 160, false)) {
		return out, cms.ErrNotFound
	}
	editor := cms.Editor(p, org)
	rows, e := s.pool.Query(ctx, `select h.article_id,h.locale,h.category,r.title,r.state,r.version::text
 from help.article h join help.revision r on r.tenant_id=h.tenant_id and r.article_id=h.article_id and r.version=h.current_version
 where h.tenant_id=$1 and h.organization_id=$2 and h.locale=$3 and h.article_id>$4
 and($5::boolean or r.state='published')and($6=''or strpos(lower((r.title||E'\n'||r.body)collate pg_catalog.pg_unicode_fast),lower($6::text collate pg_catalog.pg_unicode_fast))>0)
 order by h.article_id limit 51`, p.TenantID, org, locale, after, editor, strings.TrimSpace(query))
	if e != nil {
		return out, e
	}
	defer rows.Close()
	for rows.Next() {
		var v cms.Summary
		if e = rows.Scan(&v.ID, &v.Locale, &v.Category, &v.Title, &v.State, &v.Version); e != nil {
			return out, e
		}
		out.Items = append(out.Items, v)
	}
	if e = rows.Err(); e != nil {
		return out, e
	}
	if len(out.Items) > 50 {
		out.Items = out.Items[:50]
		out.Next = out.Items[49].ID
	}
	return out, nil
}
func (s *HelpCMS) History(ctx context.Context, p identity.Principal, org, id string, before int64) (cms.Page, error) {
	out := cms.Page{Items: []cms.Summary{}}
	if !cms.ID(id) || before < 0 || !cms.Editor(p, org) {
		return out, cms.ErrNotFound
	}
	rows, e := s.pool.Query(ctx, `select h.article_id,h.locale,h.category,r.title,r.state,r.version::text
 from help.article h join help.revision r on r.tenant_id=h.tenant_id and r.article_id=h.article_id
 where h.tenant_id=$1 and h.organization_id=$2 and h.article_id=$3 and($4::bigint=0 or r.version<$4)order by r.version desc limit 51`, p.TenantID, org, id, before)
	if e != nil {
		return out, e
	}
	defer rows.Close()
	for rows.Next() {
		var v cms.Summary
		if e = rows.Scan(&v.ID, &v.Locale, &v.Category, &v.Title, &v.State, &v.Version); e != nil {
			return out, e
		}
		out.Items = append(out.Items, v)
	}
	if e = rows.Err(); e != nil {
		return out, e
	}
	if len(out.Items) > 50 {
		out.Items = out.Items[:50]
		out.Next = out.Items[49].Version
	}
	return out, nil
}
````

### FILE: `internal/platform/postgres/help_cms_browser_integration_test.go`

```yaml
block_id: "GO-CONNECTED-HELP-CMS:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a1020823ca19507782637bca35091351d60e3a4de1e88fcdf3c9626717a8f15f"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED CMS browser: only tenant/org seeded; all six article mutations from forms.
import (
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestHelpCMSBrowser(t *testing.T) {
	if os.Getenv("ELITE_HELP_CMS_BROWSER") != "1" {
		t.Skip("explicit local browser fixture required")
	}
	raw := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if raw == "" {
		t.Fatal("owned database required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	pool, e := pgxpool.New(ctx, raw)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant := uuid.NewString()

	if _, e = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'network-browser','Synthetic','Synthetic')`, tenant); e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'root','root','Synthetic','franchisor')`, tenant); e != nil {
		t.Fatal(e)
	}
	store, e := db.NewHelpCMS(pool)
	if e != nil {
		t.Fatal(e)
	}
	verifier, token := catalogRoleIssuer(t)
	identities := map[string]any{}

	for _, name := range []string{"editor", "reader", "foreign", "unprivileged"} {
		permissions := []string{"help:read"}
		orgs := []string{"root"}
		switch name {
		case "editor":
			permissions = []string{"help:write", "help:publish"}
		case "foreign":
			orgs = []string{"foreign"}
		case "unprivileged":
			permissions = []string{}
		}
		identities[name] = map[string]any{"subject": name, "tenantId": tenant, "permissions": permissions, "organizations": orgs, "accessToken": token(name, tenant, permissions, orgs)}
	}
	mux := http.NewServeMux()
	httpapi.HelpCMSModule{Service: store}.Register(mux, verifier)
	var mu sync.Mutex
	counts := map[string]int{}
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			mu.Lock()
			counts[r.URL.Path]++
			mu.Unlock()
		}
		mux.ServeHTTP(w, r)
	}))
	defer api.Close()
	web, node := os.Getenv("ELITE_WEB_ROOT"), os.Getenv("ELITE_NODE_BIN")
	if !filepath.IsAbs(web) || !filepath.IsAbs(node) {
		t.Fatal("fixed runtime/root required")
	}
	listener, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	address := listener.Addr().String()
	listener.Close()
	target, _ := url.Parse("http://" + address)
	edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
	defer edge.Close()
	env := []string{}
	for _, v := range os.Environ() {
		k := strings.ToUpper(strings.SplitN(v, "=", 2)[0])
		if strings.HasPrefix(k, "ELITE_") || strings.HasPrefix(k, "CATALOG_") || strings.HasPrefix(k, "PUBLIC_") || strings.HasPrefix(k, "ENTERPRISE_") || strings.Contains(k, "DATABASE_URL") || k == "APP_BASE_URL" || k == "AUTH_SESSION_SECRET" || k == "BUSINESS_CONFIG_FILE" {
			continue
		}
		env = append(env, v)
	}
	identitiesJSON, _ := json.Marshal(identities)
	env = append(env, "ELITE_HELP_CMS_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_HELP_IDENTITIES="+string(identitiesJSON), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL,
		"APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-help-cms-"+uuid.NewString(), "BUSINESS_CONFIG_FILE=help.cms.fixture.json",
		"CATALOG_RELEASE_ENABLED=false", "PUBLIC_SITE_ORIGIN=https://catalog.example.invalid", "PUBLIC_INDEXING_ENABLED=0", "ENTERPRISE_TENANT_CODE=role-catalog", "ENTERPRISE_ORGANIZATION_CODE=role-org", "NEXT_TELEMETRY_DISABLED=1")
	artifacts, e := os.MkdirTemp(web, "help-cms-artifacts-")
	if e != nil {
		t.Fatal(e)
	}
	log, e := os.Create(filepath.Join(artifacts, "next.log"))
	if e != nil {
		t.Fatal(e)
	}
	defer log.Close()
	server := exec.CommandContext(ctx, node, filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", strings.Split(address, ":")[1])
	server.Dir, server.Env, server.Stdout, server.Stderr = web, env, log, log
	if e = server.Start(); e != nil {
		t.Fatal(e)
	}
	defer func() { server.Process.Kill(); server.Wait() }()
	client := &http.Client{Timeout: time.Second}
	ready := false
	for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
		r, e := client.Get(target.String() + "/icon.svg")
		if e == nil {
			r.Body.Close()
			if r.StatusCode == 200 {
				ready = true
				break
			}
		}
		select {
		case <-ctx.Done():
			t.Fatal("startup cancelled")
		case <-time.After(100 * time.Millisecond):
		}
	}
	if !ready {
		t.Fatal("Next startup", artifacts)
	}
	gate := filepath.Join(web, "microsoft_playwright_browser_gate")
	command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/help-cms-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers"))
	command.Dir, command.Env = gate, env
	output, e := command.CombinedOutput()
	if x := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); x != nil {
		t.Fatal(x)
	}
	if e != nil || !strings.Contains(string(output), "1 passed") {
		t.Fatalf("catalog role browser %v\n%s\nartifacts=%s", e, output, artifacts)
	}
	var actual string

	e = pool.QueryRow(ctx, `select jsonb_build_array(
 (select count(*)from help.article where tenant_id=$1),
 (select count(*)from help.revision where tenant_id=$1),
 (select count(*)from platform.outbox_event where tenant_id=$1))::text`, tenant).Scan(&actual)
	if e != nil || actual != "[2, 6, 6]" {
		t.Fatal("durable effects", actual, e)
	}
	mu.Lock()
	countJSON, _ := json.Marshal(counts)
	total := 0
	for _, n := range counts {
		total += n
	}
	mu.Unlock()
	if total != 6 {
		t.Fatal("unexpected backend writes", total)
	}
	if e = os.WriteFile(filepath.Join(artifacts, "api-post-counts.json"), countJSON, 0600); e != nil {
		t.Fatal(e)
	}
	t.Logf("HELP_CMS_BROWSER_PASS JWE_RS256_JWKS=true create_and_archive_response_loss_GET_only=true draft_revision_publication_reader_history_withdrawal=true backend_POSTs=6 artifacts=%s", artifacts)
}
````

### FILE: `internal/platform/postgres/help_cms_integration_test.go`

```yaml
block_id: "GO-CONNECTED-HELP-CMS:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "0b0e9aa221f4356e56827a67761e95857fc921d2b484b97ec0a11092d2af3098"
variables: []
secrets_allowed: false
```

````go
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
````

### FILE: `internal/trainingbridge/release_guides_connected_test.go`

```yaml
block_id: "GO-CONNECTED-HELP-CMS:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "aca6f34c134188e7de80725d00347f75af9c8707f5ce4f7c002755a088aac083"
variables: []
secrets_allowed: false
```

````go
package trainingbridge

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"testing"
)

func TestTrainingReleaseGuidesConnected(t *testing.T) {
	f := connectedFixture(t)
	ctx := context.Background()
	hash := f.profile.Hash()
	id := uuid.NewString()
	if _, e := f.store.Start(ctx, f.learner, id, "content-onboarding", hash); e != nil {
		t.Fatal(e)
	}
	v, e := f.store.Read(ctx, f.learner, id)
	if e != nil || v.View.ProfileRevision != 2 || len(v.View.Articles) != 3 {
		t.Fatal(v, e)
	}
	expected := []string{"help-cms-view", "catalog-role-view", "training-role-view"}
	for i, article := range v.View.Articles {
		if article.ID != expected[i] || article.Version != "1.0.0" || len(article.Paragraphs) != 3 {
			t.Fatal("wrong same-release content", article)
		}
	}
	answers := map[string]string{"review": "Consulto el comando pendiente, reviso el contenido exacto y preparo una práctica nueva; la evaluación no concede permisos."}
	if _, e = f.store.Submit(ctx, f.learner, id, hash, answers); e == nil {
		t.Fatal("unread curriculum admitted")
	}
	for _, lesson := range expected {
		if _, e = f.store.Acknowledge(ctx, f.learner, id, lesson, hash); e != nil {
			t.Fatal(e)
		}
	}
	a, e := f.store.Submit(ctx, f.learner, id, hash, answers)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = f.store.Assess(ctx, f.learner, a.RequestID, a.PayloadSHA256, true, "Self review"); e == nil {
		t.Fatal("self approval")
	}
	if _, e = f.store.Assess(ctx, f.reviewer, a.RequestID, a.PayloadSHA256, true, "Tres guías de la misma revisión contrastadas con la respuesta"); e != nil {
		t.Fatal(e)
	}
	// Future activation keeps this immutable, approved attempt on revision2.
	d, b := profileFixture(t)
	d.TenantID = f.learner.TenantID
	d.Revision = 3
	raw, _ := json.Marshal(d)
	next, e := ParseProfile(raw, b, Activation{true, d.ID, d.Revision, digest(raw), d.TenantID, d.OrganizationID})
	if e != nil {
		t.Fatal(e)
	}
	restarted, e := NewStore(f.pool, next)
	if e != nil {
		t.Fatal(e)
	}
	v, e = restarted.Read(ctx, f.learner, id)
	if e != nil || v.CurrentProfile || v.View.ProfileRevision != 2 || v.Assessment == nil || v.Assessment.State != "approved" || v.Assessment.Reviewer != "reviewer" {
		t.Fatal("prior curriculum changed", v, e)
	}
	if _, e = restarted.Acknowledge(ctx, f.learner, id, expected[0], hash); e == nil {
		t.Fatal("stale profile mutated")
	}
	var facts, events, requests, decisions, resources int
	e = f.pool.QueryRow(ctx, `select(select count(*)from audit.event where tenant_id=$1 and resource_type='training-attempt'),(select count(*)from platform.outbox_event where tenant_id=$1 and aggregate_type='training-evidence'),(select count(*)from approval.request where tenant_id=$1 and kind='training_assessment'),(select count(*)from approval.decision where tenant_id=$1),(select count(*)from crm.service_resource where tenant_id=$1)`, f.learner.TenantID).Scan(&facts, &events, &requests, &decisions, &resources)
	if e != nil || facts != 6 || events != 6 || requests != 1 || decisions != 1 || resources != 0 {
		t.Fatal("effects", facts, events, requests, decisions, resources, e)
	}
	if f.learner.Allowed("help:publish") || f.learner.Allowed("network:admin") {
		t.Fatal("training granted authority")
	}
	t.Log("HELP_RELEASE_TRAINING_PASS revision2_three_new_guides=true reading_submission_distinct_human_review=true six_facts_six_events=true future_profile_preserves_history=true permission_grants=0")
}
````

## 6. Configuration surface

docs/HELP_CMS_REFERENCE.md. HELP_CMS_ENABLED host and features.help_cms UI opt-in, explicit help:read/write/publish and organization. PostgreSQL18 UTF8 pg_unicode_fast,5guard checks, bounded24KiB encoded command/16KiB body/32KiBHTTP,50row keysets.

## 7. Dependency bill

23new AUTHORED glue blocks;15selected owner revisions plus pure kernel extraction. No new dependency/runtime/upstream. Existing licenses/notices preserved. Immutable history and populated downgrade guard; original isolated kernel12tests16seeds retained.

## 8. Apply order

GO-HELP-CENTER-CORE0.2.0, GO-CONNECTED-HUMAN-TRAINING0.2.0, GO-HUMAN-APPROVAL-CORE, GO-ELECTROMOBILITY-APPLICATION1.18.0 and existing PG/identity owners. Browser uses catalog synthetic issuer. TS-CONNECTED-HELP-CMS-PORTAL owns shared goldens; full reference required for connected proof.

## 9. Verification

CMS PG1new3replay,3table rollback, role/org/actor/tenant/stale/immutability/history/withdrawal/pagination. Browser actual Next/BFF/Go/PG/Chromium6POST2articles6revisions6events, create/archive lost response GET afterreload0extraPOST; desktop390px inspected.5host guards/down/reapply.4goldens/11webtests/7guidealignment checks, HTTP boundaries and Next/types. Fuzz2s3seeds328160executions. Training3newlessons→humanreview6facts6events1request1decision0grants, futureprofile retains old attempt.

## 10. Reconstruction evidence

HELP_CMS_RELEASE_V402.md/json and exact profile rebuilds. FAIL864Unicode search,865export module root,866test indices,867bounded curriculum corrected with RED history and no relaxed constraints. T2804 remains active forKPI/privatei18n; later controls ordered.

