# Connected network role composition

## 1. Metadata

```yaml
pack_id: "GO-NETWORK-ROLE-COMPOSITION"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Recoverable organization/agreement/branch lifecycle by authorized role, using original network owners; local infrastructure reference."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["MARKDOWN-COMPOSITOR0.3.0", "V402 existing owner closure"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-12"
```

## 2. Applicability

LIBRARY_INFRASTRUCTURE T2804/J5: existing fulfillment organization/agreement owners, tenant/hierarchy/territory guards, approval CanonicalPayload and identity/BFF. No in-memory onboarding-core promotion.

## 3. Architecture contract

Four original PG writer SQL bodies extracted exactly into transaction adapter. Same service states and outbox. Immutable command receipt binds actor/scope/hash/exact result; command lock and preallocated aggregate IDs, independent random event IDs. No parent permission inheritance or automatic grants.

## 4. Exact file manifest

```text
CREATE internal/platform/postgres/identity_j5_integration_test.go
CREATE cmd/electromobility-api/network_role.go
CREATE cmd/electromobility-api/network_role_test.go
CREATE db/migrations/0077_network_role_commands.down.sql
CREATE db/migrations/0077_network_role_commands.up.sql
CREATE internal/networkrole/contract.go
CREATE internal/networkrole/contract_test.go
CREATE internal/platform/httpapi/network_role.go
CREATE internal/platform/postgres/network_owner_tx.go
CREATE internal/platform/postgres/network_role.go
CREATE internal/platform/postgres/network_role_browser_integration_test.go
CREATE internal/platform/postgres/network_role_integration_test.go
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/network_role.go`

```yaml
block_id: "GO-NETWORK-ROLE-COMPOSITION:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7696487e8bab2f2f6245aa14fbe765fc8d60593b6d455aff2d48b9336c895ca8"
variables: []
secrets_allowed: false
```

````go
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
````

### FILE: `cmd/electromobility-api/network_role_test.go`

```yaml
block_id: "GO-NETWORK-ROLE-COMPOSITION:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "4e05db3a93e3b05ff450d696b9a5ff999c46a22244fa0f52ca475a4f936ebf72"
variables: []
secrets_allowed: false
```

````go
package main
import("context";"os";"strings";"testing";nr "elite.local/enterprise/internal/networkrole";"elite.local/enterprise/internal/platform/identity";"elite.local/enterprise/internal/platform/postgres";"github.com/google/uuid";"github.com/jackc/pgx/v5/pgxpool")
func TestNetworkRoleHostGuards(t *testing.T){
 ctx:=context.Background();calls:=0
 if m,e:=selectedNetworkRoleModule(ctx,nil,func(k string)string{calls++;if k!="NETWORK_ROLE_ENABLED"{t.Fatal("unexpected lookup")};return"false"});m!=nil||e!=nil||calls!=1{t.Fatal(m,e,calls)}
 if _,e:=selectedNetworkRoleModule(ctx,nil,func(string)string{return"TRUE"});e==nil{t.Fatal("invalid flag")}
 raw:=os.Getenv("PAYMENT_CONNECTED_DB_URL");if raw==""{t.Skip("owned database required")};pool,e:=pgxpool.New(ctx,raw);if e!=nil{t.Fatal(e)};defer pool.Close()
 lookup:=func(string)string{return"true"}
 enabled:=func(want bool){t.Helper();m,e:=selectedNetworkRoleModule(ctx,pool,lookup);if (e==nil&&m!=nil)!=want{t.Fatal("activation mismatch",want,e)}}
 enabled(true)
 pairs:=[][2]string{
 {`alter table franchise.network_command_receipt disable trigger network_receipt_immutable`,`alter table franchise.network_command_receipt enable trigger network_receipt_immutable`},
 {`alter table org.organization disable trigger organization_hierarchy_guard`,`alter table org.organization enable trigger organization_hierarchy_guard`},
 {`alter table franchise.agreement disable trigger franchise_territory_non_overlap`,`alter table franchise.agreement enable trigger franchise_territory_non_overlap`},
 {`alter index franchise.network_receipt_actor_idx rename to network_receipt_actor_saved`,`alter index franchise.network_receipt_actor_saved rename to network_receipt_actor_idx`},
 {`alter table franchise.network_command_receipt drop constraint network_command_receipt_pkey`,`alter table franchise.network_command_receipt add constraint network_command_receipt_pkey primary key(tenant_id,command_id)`},
 }
 for _,pair:=range pairs{func(){if _,e:=pool.Exec(ctx,pair[0]);e!=nil{t.Fatal(e)};defer func(){if _,e:=pool.Exec(context.Background(),pair[1]);e!=nil{t.Error(e)}}();enabled(false)}();enabled(true)}
 down,e:=os.ReadFile("../../db/migrations/0077_network_role_commands.down.sql");if e!=nil{t.Fatal(e)}
 up,e:=os.ReadFile("../../db/migrations/0077_network_role_commands.up.sql");if e!=nil{t.Fatal(e)}
 if _,e=pool.Exec(ctx,string(down));e!=nil{t.Fatal("empty downgrade",e)};enabled(false)
 if _,e=pool.Exec(ctx,string(up));e!=nil{t.Fatal("reapply",e)};enabled(true)
 tenant:=uuid.NewString();if _,e=pool.Exec(ctx,`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'network-host','Synthetic','Synthetic')`,tenant);e!=nil{t.Fatal(e)}
 store,_:=postgres.NewNetworkRole(pool);p:=identity.Principal{TenantID:tenant,Subject:"fixture",Permissions:map[string]struct{}{"*":{}}}
 if _,e=store.Execute(ctx,p,nr.Command{CommandID:"create",Action:"create-organization",EntityID:"root",Code:"root",DisplayName:"Synthetic root",Type:"franchisor"});e!=nil{t.Fatal(e)}
 tx,e:=pool.Begin(ctx);if e!=nil{t.Fatal(e)};defer tx.Rollback(ctx)
 body:=strings.TrimSpace(string(down));body=strings.TrimPrefix(body,"begin;");body=strings.TrimSuffix(body,"commit;")
 if _,e=tx.Exec(ctx,body);e==nil{t.Fatal("populated evidence discarded")};tx.Rollback(ctx);enabled(true)
 t.Log("NETWORK_ROLE_HOST_PASS five_missing_guards_block=true empty_downgrade_reapply=true populated_downgrade_refused=true no_credentials=true")
}
````

### FILE: `db/migrations/0077_network_role_commands.down.sql`

```yaml
block_id: "GO-NETWORK-ROLE-COMPOSITION:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "288892b037cddfa09a3a917f11019ebba215a607aec0fc7c886dc615e3442819"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$begin
 if exists(select 1 from franchise.network_command_receipt)then raise exception 'cannot discard durable network command evidence';end if;
end$$;
drop table franchise.network_command_receipt;
drop function franchise.reject_network_receipt_mutation();
commit;
````

### FILE: `db/migrations/0077_network_role_commands.up.sql`

```yaml
block_id: "GO-NETWORK-ROLE-COMPOSITION:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "61b9986d1bc630dc83782422f0b0347c13b6de01ef86eb666af7cd07fd49944d"
variables: []
secrets_allowed: false
```

````sql
begin;
create table franchise.network_command_receipt(
 tenant_id uuid not null references platform.tenant(tenant_id),
 command_id text not null check(length(command_id)between 1 and 128),
 action text not null check(action in('create-organization','transition-organization','create-agreement','transition-agreement')),
 scope_organization_id text not null check(length(scope_organization_id)<=128),
 actor_subject text not null check(length(actor_subject)between 1 and 256),
 request_sha256 text not null check(request_sha256~'^[0-9a-f]{64}$'),
 entity_payload jsonb not null check(jsonb_typeof(entity_payload)='object'),
 entity_sha256 text not null check(entity_sha256~'^[0-9a-f]{64}$'),
 recorded_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,command_id)
);
create function franchise.reject_network_receipt_mutation()returns trigger language plpgsql as $$begin raise exception 'network receipt is immutable';end$$;
create trigger network_receipt_immutable before update or delete on franchise.network_command_receipt for each row execute function franchise.reject_network_receipt_mutation();
create index network_receipt_actor_idx on franchise.network_command_receipt(tenant_id,actor_subject,scope_organization_id,recorded_at,command_id);
commit;
````

### FILE: `internal/networkrole/contract.go`

```yaml
block_id: "GO-NETWORK-ROLE-COMPOSITION:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "372fc404e687fe964abc949d5f9efb67cb21645a76c6e32a802318c4a779ac67"
variables: []
secrets_allowed: false
```

````go
// AUTHORED bounded transport and durable-result contract around existing fulfillment.
package networkrole

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

var ErrInvalid = errors.New("invalid network command")
var ErrNotFound = errors.New("network result unavailable")
var ErrConflict = errors.New("network command conflict")

// Existing organization_organization_code_check, surfaced before transport.
var organizationCode = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)
var idPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

func ID(s string) bool { return idPattern.MatchString(s) }
func Text(s string, max int) bool {
	if len(s) < 2 || len(s) > max || !utf8.ValidString(s) || strings.TrimSpace(s) != s {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}

type Command struct {
	CommandID           string `json:"command_id"`
	Action              string `json:"action"`
	ScopeOrganizationID string `json:"scope_organization_id"`
	EntityID            string `json:"entity_id"`
	Code                string `json:"code,omitempty"`
	DisplayName         string `json:"display_name,omitempty"`
	Type                string `json:"type,omitempty"`
	TerritoryCode       string `json:"territory_code,omitempty"`
	TermsVersion        string `json:"terms_version,omitempty"`
	StartsOn            string `json:"starts_on,omitempty"`
	EndsOn              string `json:"ends_on,omitempty"`
	Current             string `json:"current,omitempty"`
	Target              string `json:"target,omitempty"`
	Version             string `json:"version,omitempty"`
}

func Permission(action string) string {
	switch action {
	case "create-organization", "transition-organization":
		return "network:admin"
	case "create-agreement", "transition-agreement":
		return "franchise:write"
	}
	return ""
}
func (c Command) Valid() bool {
	if !ID(c.CommandID) || !ID(c.EntityID) || (c.ScopeOrganizationID != "" && !ID(c.ScopeOrganizationID)) {
		return false
	}
	switch c.Action {
	case "create-organization":
		types := map[string]bool{"enterprise": true, "franchisor": true, "franchisee": true, "factory": true, "warehouse": true, "store": true, "service_center": true}
		root := c.Type == "enterprise" || c.Type == "franchisor"
		if !types[c.Type] || (root && c.ScopeOrganizationID != "") || (!root && c.ScopeOrganizationID == "") {
			return false
		}
		return len(c.Code) <= 128 && organizationCode.MatchString(c.Code) && Text(c.DisplayName, 100) && ID(c.Type) && c.TerritoryCode == "" && c.TermsVersion == "" && c.StartsOn == "" && c.EndsOn == "" && c.Current == "" && c.Target == "" && c.Version == ""
	case "create-agreement":
		start, e := time.Parse("2006-01-02", c.StartsOn)
		if e != nil || start.Year() < 1 {
			return false
		}
		if c.EndsOn != "" {
			end, e := time.Parse("2006-01-02", c.EndsOn)
			if e != nil || !end.After(start) {
				return false
			}
		}
		return ID(c.ScopeOrganizationID) && ID(c.TerritoryCode) && ID(c.TermsVersion) && c.Code == "" && c.DisplayName == "" && c.Type == "" && c.Current == "" && c.Target == "" && c.Version == ""
	case "transition-organization", "transition-agreement":
		v, e := strconv.ParseInt(c.Version, 10, 64)
		return e == nil && v > 0 && v < math.MaxInt64 && strconv.FormatInt(v, 10) == c.Version && ID(c.ScopeOrganizationID) && ID(c.Current) && ID(c.Target) && c.Code == "" && c.DisplayName == "" && c.Type == "" && c.TerritoryCode == "" && c.TermsVersion == "" && c.StartsOn == "" && c.EndsOn == "" && (c.Action != "transition-organization" || c.ScopeOrganizationID == c.EntityID)
	}
	return false
}
func Authorized(p identity.Principal, action, scope string) bool {
	perm := Permission(action)
	if p.TenantID == "" || p.Subject == "" || perm == "" || !p.Allowed(perm) {
		return false
	}
	if scope == "" {
		return action == "create-organization" && p.Allowed("network:bootstrap")
	}
	return p.AllowedOrganization(scope)
}
func Canonical(v any) (json.RawMessage, string, error) {
	raw, e := json.Marshal(v)
	if e != nil {
		return nil, "", e
	}
	return approval.CanonicalPayload(raw)
}

type Entity struct {
	Kind                 string `json:"kind"`
	ID                   string `json:"id"`
	OrganizationID       string `json:"organization_id"`
	Version              string `json:"version"`
	State                string `json:"state"`
	Code                 string `json:"code,omitempty"`
	DisplayName          string `json:"display_name,omitempty"`
	Type                 string `json:"type,omitempty"`
	ParentOrganizationID string `json:"parent_organization_id,omitempty"`
	TerritoryCode        string `json:"territory_code,omitempty"`
	TermsVersion         string `json:"terms_version,omitempty"`
	StartsOn             string `json:"starts_on,omitempty"`
	EndsOn               string `json:"ends_on,omitempty"`
}
type Receipt struct {
	CommandID           string    `json:"command_id"`
	Action              string    `json:"action"`
	ScopeOrganizationID string    `json:"scope_organization_id"`
	Actor               string    `json:"actor"`
	RequestSHA256       string    `json:"request_sha256"`
	Entity              Entity    `json:"entity"`
	RecordedAt          time.Time `json:"recorded_at"`
	Replay              bool      `json:"replay"`
}
````

### FILE: `internal/networkrole/contract_test.go`

```yaml
block_id: "GO-NETWORK-ROLE-COMPOSITION:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "4bea0ef65e7a52e69fd807ef85e41feae8de5dd460b4dc330e01502860f55e18"
variables: []
secrets_allowed: false
```

````go
package networkrole

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestNetworkRoleWireGoldens(t *testing.T) {
	raw, e := os.ReadFile("../../deploy/network/role-goldens.json")
	if e != nil {
		t.Fatal(e)
	}
	var rows []struct {
		Command   Command `json:"command"`
		Canonical string  `json:"canonical"`
		SHA       string  `json:"sha256"`
	}
	if json.Unmarshal(raw, &rows) != nil || len(rows) != 4 {
		t.Fatal("four shared goldens required")
	}
	for _, r := range rows {
		if !r.Command.Valid() {
			t.Fatal("invalid golden", r.Command)
		}
		body, hash, e := Canonical(r.Command)
		if e != nil || string(body) != r.Canonical || hash != r.SHA {
			t.Fatal(r.Command.Action, string(body), e)
		}
	}
	for _, s := range []string{"0", "01", "-1", "9223372036854775807", "9223372036854775808"} {
		c := rows[2].Command
		c.Version = s
		if c.Valid() {
			t.Fatal("unsafe version", s)
		}
	}
	for _, s := range []string{"UPPER", "space code", strings.Repeat("x", 129)} {
		c := rows[0].Command
		c.Code = s
		if c.Valid() {
			t.Fatal("schema code", s)
		}
	}
}
func FuzzNetworkRoleWire(f *testing.F) {
	for _, s := range []string{`{"command_id":"x","action":"transition-organization","scope_organization_id":"org","entity_id":"org","current":"active","target":"closed","version":"1"}`, `{"action":"create-organization"}`, `{}`} {
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
			t.Fatal("valid command is not canonical", e)
		}
		var again Command
		if json.Unmarshal(body, &again) != nil || again != c || !again.Valid() {
			t.Fatal("canonical round trip")
		}
	})
}
````

### FILE: `internal/platform/httpapi/network_role.go`

```yaml
block_id: "GO-NETWORK-ROLE-COMPOSITION:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d84807a4a9df7c04c6ce34b83556ab23a2ef26ec6eee3d54725949198c36a4f6"
variables: []
secrets_allowed: false
```

````go
// AUTHORED bounded transport for original organization/agreement owners.
package httpapi

import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/approval"
	nr "elite.local/enterprise/internal/networkrole"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type NetworkRoleService interface {
	Execute(context.Context, identity.Principal, nr.Command) (nr.Receipt, error)
	Result(context.Context, identity.Principal, string, string) (nr.Receipt, error)
	Entity(context.Context, identity.Principal, string, string, string) (nr.Entity, error)
}
type NetworkRoleModule struct{ Service NetworkRoleService }

func networkReply(w http.ResponseWriter, v any, e error) {
	if e != nil {
		switch {
		case errors.Is(e, nr.ErrInvalid):
			writeProblem(w, 400, "NETWORK_INVALID", "invalid network command")
		case errors.Is(e, nr.ErrNotFound):
			writeProblem(w, 404, "NETWORK_UNAVAILABLE", "result unavailable in this scope")
		case errors.Is(e, nr.ErrConflict):
			writeProblem(w, 409, "NETWORK_CONFLICT", "consult saved result and current version")
		default:
			writeProblem(w, 503, "NETWORK_UNCONFIRMED", "result not confirmed; consult saved command")
		}
		return
	}
	writeJSON(w, 200, v)
}
func (m NetworkRoleModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	auth := func(w http.ResponseWriter, r *http.Request) (identity.Principal, bool) {
		p, e := authenticate(r.Context(), r.Header.Get("Authorization"), verifier)
		if e != nil {
			writeProblem(w, 401, "UNAUTHENTICATED", "valid session required")
			return p, false
		}
		if !p.Allowed("network:admin") && !p.Allowed("franchise:write") {
			writeProblem(w, 403, "FORBIDDEN", "network or franchise permission required")
			return p, false
		}
		return p, true
	}
	mux.HandleFunc("POST /v1/franchise/network/commands", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r)
		if !ok {
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			writeProblem(w, 415, "CONTENT_TYPE", "application/json required")
			return
		}
		raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
		if e != nil {
			writeProblem(w, 413, "BODY_LIMIT", "network command exceeds32KiB")
			return
		}
		if _, _, e = approval.CanonicalPayload(raw); e != nil {
			networkReply(w, nil, nr.ErrInvalid)
			return
		}
		var c nr.Command
		decoder := json.NewDecoder(bytes.NewReader(raw))
		decoder.DisallowUnknownFields()
		if decoder.Decode(&c) != nil || decoder.Decode(new(any)) != io.EOF || !c.Valid() {
			networkReply(w, nil, nr.ErrInvalid)
			return
		}
		if !nr.Authorized(p, c.Action, c.ScopeOrganizationID) {
			writeProblem(w, 403, "NETWORK_FORBIDDEN", "action or organization unavailable")
			return
		}
		v, e := m.Service.Execute(r.Context(), p, c)
		networkReply(w, v, e)
	})
	mux.HandleFunc("GET /v1/franchise/network/commands/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r)
		if !ok {
			return
		}
		q := r.URL.Query()
		// protectedGet omits an empty root scope. Result still requires original root authority.
		if (len(q) != 0 && (len(q) != 1 || len(q["scope_organization_id"]) != 1)) || !nr.ID(r.PathValue("id")) {
			networkReply(w, nil, nr.ErrInvalid)
			return
		}
		v, e := m.Service.Result(r.Context(), p, q.Get("scope_organization_id"), r.PathValue("id"))
		networkReply(w, v, e)
	})
	mux.HandleFunc("GET /v1/franchise/network/entities/{kind}/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth(w, r)
		if !ok {
			return
		}
		q := r.URL.Query()
		if len(q) != 1 || len(q["scope_organization_id"]) != 1 || !nr.ID(r.PathValue("id")) {
			networkReply(w, nil, nr.ErrInvalid)
			return
		}
		v, e := m.Service.Entity(r.Context(), p, r.PathValue("kind"), q.Get("scope_organization_id"), r.PathValue("id"))
		networkReply(w, v, e)
	})
}
````

### FILE: `internal/platform/postgres/network_owner_tx.go`

```yaml
block_id: "GO-NETWORK-ROLE-COMPOSITION:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "17ae0881ee4aa83e0d7be5bd70629b877254f6f82ffd8398c3b5e847f1108f92"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED narrow adapter: only four existing network writers are exposed.
import (
	"context"
	"elite.local/enterprise/internal/fulfillment"
	"github.com/jackc/pgx/v5"
)

type networkTxRepository struct{ tx pgx.Tx }

var _ fulfillment.Repository = networkTxRepository{}

func (r networkTxRepository) CreateShipment(a0 context.Context, a1 string, a2 string, a3 fulfillment.Shipment) error {
	return fulfillment.ErrConflict
}
func (r networkTxRepository) TransitionShipment(a0 context.Context, a1 string, a2 string, a3 string, a4 string, a5 string, a6 string, a7 string, a8 string) error {
	return fulfillment.ErrConflict
}
func (r networkTxRepository) CreateCustomerTransport(a0 context.Context, a1 string, a2 fulfillment.TransportIDs, a3 fulfillment.CustomerTransportCommand) (fulfillment.CustomerTransport, error) {
	return fulfillment.CustomerTransport{}, fulfillment.ErrConflict
}
func (r networkTxRepository) RecordProviderReport(a0 context.Context, a1 string, a2 string, a3 string, a4 fulfillment.ProviderReportCommand) (fulfillment.CustomerTransport, error) {
	return fulfillment.CustomerTransport{}, fulfillment.ErrConflict
}
func (r networkTxRepository) OpenServiceCase(a0 context.Context, a1 string, a2 string, a3 fulfillment.ServiceCase) error {
	return fulfillment.ErrConflict
}
func (r networkTxRepository) TransitionServiceCase(a0 context.Context, a1 string, a2 string, a3 string, a4 string, a5 string, a6 int64, a7 string) error {
	return fulfillment.ErrConflict
}
func (r networkTxRepository) CreateRecall(a0 context.Context, a1 string, a2 string, a3 fulfillment.Recall) error {
	return fulfillment.ErrConflict
}
func (r networkTxRepository) ActivateRecall(a0 context.Context, a1 string, a2 string, a3 string) error {
	return fulfillment.ErrConflict
}
func (r networkTxRepository) AddRecallUnit(a0 context.Context, a1 string, a2 string, a3 string, a4 string, a5 string) error {
	return fulfillment.ErrConflict
}
func (r networkTxRepository) CreateOrganization(a0 context.Context, a1 string, a2 string, a3 fulfillment.Organization) error {
	return createOrganizationInTx(a0, r.tx, a1, a2, a3)
}
func (r networkTxRepository) TransitionOrganization(a0 context.Context, a1 string, a2 string, a3 string, a4 string, a5 int64, a6 string) error {
	return transitionOrganizationInTx(a0, r.tx, a1, a2, a3, a4, a5, a6)
}
func (r networkTxRepository) CreateAgreement(a0 context.Context, a1 string, a2 string, a3 fulfillment.Agreement) error {
	return createAgreementInTx(a0, r.tx, a1, a2, a3)
}
func (r networkTxRepository) TransitionAgreement(a0 context.Context, a1 string, a2 string, a3 string, a4 string, a5 string, a6 int64, a7 string) error {
	return transitionAgreementInTx(a0, r.tx, a1, a2, a3, a4, a5, a6, a7)
}
func (r networkTxRepository) QueueMessage(a0 context.Context, a1 string, a2 string, a3 fulfillment.Message) error {
	return fulfillment.ErrConflict
}
func (r networkTxRepository) TransitionMessage(a0 context.Context, a1 string, a2 string, a3 string, a4 string, a5 string, a6 string, a7 string) error {
	return fulfillment.ErrConflict
}
````

### FILE: `internal/platform/postgres/network_role.go`

```yaml
block_id: "GO-NETWORK-ROLE-COMPOSITION:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1b588b8bbc040cca926aba135d2d634fd311fe7165765c4bb2e9f450ebb7b311"
variables: []
secrets_allowed: false
```

````go
// AUTHORED transaction/recovery glue; all business mutations use existing fulfillment writers.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/fulfillment"
	nr "elite.local/enterprise/internal/networkrole"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"time"
)

type NetworkRole struct{ pool *pgxpool.Pool }

func NewNetworkRole(pool *pgxpool.Pool) (*NetworkRole, error) {
	if pool == nil {
		return nil, nr.ErrInvalid
	}
	return &NetworkRole{pool}, nil
}

type networkRowReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}
type networkCommandIDs struct{ first string }

func (g *networkCommandIDs) New() string {
	if g.first != "" {
		v := g.first
		g.first = ""
		return v
	}
	return randomid.Generator{}.New()
}
func networkError(e error) error {
	if errors.Is(e, pgx.ErrNoRows) {
		return nr.ErrNotFound
	}
	if errors.Is(e, fulfillment.ErrConflict) {
		return nr.ErrConflict
	}
	if e = fulfillmentConstraint(e); errors.Is(e, fulfillment.ErrConflict) {
		return nr.ErrConflict
	}
	return e
}
func readNetworkReceipt(ctx context.Context, q networkRowReader, tenant, id string) (nr.Receipt, error) {
	var v nr.Receipt
	var raw []byte
	var hash string
	e := q.QueryRow(ctx, `select command_id,action,scope_organization_id,actor_subject,request_sha256,entity_payload,entity_sha256,recorded_at from franchise.network_command_receipt where tenant_id=$1 and command_id=$2`, tenant, id).Scan(&v.CommandID, &v.Action, &v.ScopeOrganizationID, &v.Actor, &v.RequestSHA256, &raw, &hash, &v.RecordedAt)
	if e != nil {
		return v, networkError(e)
	}
	if json.Unmarshal(raw, &v.Entity) != nil {
		return v, nr.ErrConflict
	}
	_, actual, e := nr.Canonical(v.Entity)
	if e != nil || actual != hash {
		return v, nr.ErrConflict
	}
	return v, nil
}
func readNetworkEntity(ctx context.Context, q networkRowReader, tenant, kind, id string) (nr.Entity, error) {
	var v nr.Entity
	var version int64
	var e error
	v.Kind = kind
	switch kind {
	case "organization":
		e = q.QueryRow(ctx, `select organization_id,organization_id,coalesce(parent_organization_id,''),organization_code,display_name,organization_type,status,version from org.organization where tenant_id=$1 and organization_id=$2`, tenant, id).Scan(&v.ID, &v.OrganizationID, &v.ParentOrganizationID, &v.Code, &v.DisplayName, &v.Type, &v.State, &version)
	case "agreement":
		e = q.QueryRow(ctx, `select agreement_id,franchise_organization_id,territory_code,terms_version,starts_on::text,coalesce(ends_on::text,''),status,version from franchise.agreement where tenant_id=$1 and agreement_id=$2`, tenant, id).Scan(&v.ID, &v.OrganizationID, &v.TerritoryCode, &v.TermsVersion, &v.StartsOn, &v.EndsOn, &v.State, &version)
	default:
		return v, nr.ErrInvalid
	}
	v.Version = strconv.FormatInt(version, 10)
	return v, networkError(e)
}
func (s *NetworkRole) Result(ctx context.Context, p identity.Principal, scope, id string) (nr.Receipt, error) {
	if !nr.ID(id) || p.TenantID == "" || p.Subject == "" {
		return nr.Receipt{}, nr.ErrNotFound
	}
	v, e := readNetworkReceipt(ctx, s.pool, p.TenantID, id)
	if e != nil {
		return v, e
	}
	if v.Actor != p.Subject || v.ScopeOrganizationID != scope || !nr.Authorized(p, v.Action, scope) {
		return nr.Receipt{}, nr.ErrNotFound
	}
	v.Replay = true
	return v, nil
}
func (s *NetworkRole) Entity(ctx context.Context, p identity.Principal, kind, scope, id string) (nr.Entity, error) {
	action := "transition-" + kind
	if !nr.ID(id) || !nr.Authorized(p, action, scope) {
		return nr.Entity{}, nr.ErrNotFound
	}
	v, e := readNetworkEntity(ctx, s.pool, p.TenantID, kind, id)
	if e != nil {
		return v, e
	}
	if v.OrganizationID != scope {
		return nr.Entity{}, nr.ErrNotFound
	}
	return v, nil
}
func (s *NetworkRole) Execute(ctx context.Context, p identity.Principal, c nr.Command) (nr.Receipt, error) {
	var empty nr.Receipt
	if !c.Valid() {
		return empty, nr.ErrInvalid
	}
	if !nr.Authorized(p, c.Action, c.ScopeOrganizationID) {
		return empty, nr.ErrNotFound
	}
	_, hash, e := nr.Canonical(c)
	if e != nil {
		return empty, e
	}
	tx, e := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if e != nil {
		return empty, e
	}
	defer tx.Rollback(ctx)
	if _, e = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended('network-command:'||$1::text||':'||$2::text,0))`, p.TenantID, c.CommandID); e != nil {
		return empty, e
	}
	old, e := readNetworkReceipt(ctx, tx, p.TenantID, c.CommandID)
	if e == nil {
		if old.Actor != p.Subject || old.Action != c.Action || old.ScopeOrganizationID != c.ScopeOrganizationID || old.RequestSHA256 != hash {
			return empty, nr.ErrConflict
		}
		old.Replay = true
		return old, tx.Commit(ctx)
	}
	if !errors.Is(e, nr.ErrNotFound) {
		return empty, e
	}
	ids := &networkCommandIDs{}
	if c.Action == "create-organization" || c.Action == "create-agreement" {
		ids.first = c.EntityID
	}
	service := fulfillment.NewService(networkTxRepository{tx: tx}, ids)
	kind := ""
	switch c.Action {
	case "create-organization":
		kind = "organization"
		_, e = service.CreateOrganization(ctx, p.TenantID, fulfillment.Organization{ParentOrganizationID: c.ScopeOrganizationID, Code: c.Code, DisplayName: c.DisplayName, Type: c.Type})
	case "transition-organization":
		kind = "organization"
		version, _ := strconv.ParseInt(c.Version, 10, 64)
		e = service.TransitionOrganization(ctx, p.TenantID, c.EntityID, c.Current, c.Target, version)
	case "create-agreement":
		kind = "agreement"
		start, _ := time.Parse("2006-01-02", c.StartsOn)
		var end *time.Time
		if c.EndsOn != "" {
			v, _ := time.Parse("2006-01-02", c.EndsOn)
			end = &v
		}
		_, e = service.CreateAgreement(ctx, p.TenantID, fulfillment.Agreement{OrganizationID: c.ScopeOrganizationID, TerritoryCode: c.TerritoryCode, TermsVersion: c.TermsVersion, StartsOn: start, EndsOn: end})
	case "transition-agreement":
		kind = "agreement"
		version, _ := strconv.ParseInt(c.Version, 10, 64)
		e = service.TransitionAgreement(ctx, p.TenantID, c.ScopeOrganizationID, c.EntityID, c.Current, c.Target, version)
	}
	if e != nil {
		return empty, networkError(e)
	}
	entity, e := readNetworkEntity(ctx, tx, p.TenantID, kind, c.EntityID)
	if e != nil {
		return empty, e
	}
	raw, entityHash, e := nr.Canonical(entity)
	if e != nil {
		return empty, e
	}
	receipt := nr.Receipt{CommandID: c.CommandID, Action: c.Action, ScopeOrganizationID: c.ScopeOrganizationID, Actor: p.Subject, RequestSHA256: hash, Entity: entity}
	e = tx.QueryRow(ctx, `insert into franchise.network_command_receipt(tenant_id,command_id,action,scope_organization_id,actor_subject,request_sha256,entity_payload,entity_sha256)values($1,$2,$3,$4,$5,$6,$7::jsonb,$8) returning recorded_at`, p.TenantID, c.CommandID, c.Action, c.ScopeOrganizationID, p.Subject, hash, string(raw), entityHash).Scan(&receipt.RecordedAt)
	if e != nil {
		return empty, networkError(e)
	}
	return receipt, networkError(tx.Commit(ctx))
}
````

### FILE: `internal/platform/postgres/network_role_browser_integration_test.go`

```yaml
block_id: "GO-NETWORK-ROLE-COMPOSITION:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "157d6d77607ab0e541fffd8088b5fb7d8d4e3d01253009bc2935e21966c6ee37"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED role browser: only tenant/org/supplier/catalog masters seeded; all J2 writes from forms.
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

func TestNetworkRoleBrowser(t *testing.T) {
	if os.Getenv("ELITE_NETWORK_ROLE_BROWSER") != "1" {
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
	store, e := db.NewNetworkRole(pool)
	if e != nil {
		t.Fatal(e)
	}
	verifier, token := catalogRoleIssuer(t)
	identities := map[string]any{}

	for _, name := range []string{"bootstrap", "limited", "foreign", "unprivileged"} {
		permissions := []string{"network:admin"}
		orgs := []string{"foreign"}
		switch name {
		case "bootstrap":
			permissions = []string{"*"}
			orgs = []string{"fixture-scope"}
		case "unprivileged":
			permissions = []string{}
		}
		identities[name] = map[string]any{"subject": name, "tenantId": tenant, "permissions": permissions, "organizations": orgs, "accessToken": token(name, tenant, permissions, orgs)}
	}
	mux := http.NewServeMux()
	httpapi.NetworkRoleModule{Service: store}.Register(mux, verifier)
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
	env = append(env, "ELITE_NETWORK_ROLE_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_NETWORK_IDENTITIES="+string(identitiesJSON), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL,
		"APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-network-role-"+uuid.NewString(), "BUSINESS_CONFIG_FILE=network.role.fixture.json",
		"CATALOG_RELEASE_ENABLED=false", "PUBLIC_SITE_ORIGIN=https://catalog.example.invalid", "PUBLIC_INDEXING_ENABLED=0", "ENTERPRISE_TENANT_CODE=role-catalog", "ENTERPRISE_ORGANIZATION_CODE=role-org", "NEXT_TELEMETRY_DISABLED=1")
	artifacts, e := os.MkdirTemp(web, "network-role-artifacts-")
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
	command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/network-role-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers"))
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
 (select count(*)from org.organization where tenant_id=$1),
 (select count(*)from franchise.agreement where tenant_id=$1 and status='terminated'),
 (select count(*)from franchise.network_command_receipt where tenant_id=$1),
 (select count(*)from platform.outbox_event where tenant_id=$1))::text`, tenant).Scan(&actual)
	if e != nil || actual != "[3, 1, 13, 13]" {
		t.Fatal("durable effects", actual, e)
	}
	mu.Lock()
	countJSON, _ := json.Marshal(counts)
	total := 0
	for _, n := range counts {
		total += n
	}
	mu.Unlock()
	if total != 13 {
		t.Fatal("unexpected backend writes", total)
	}
	if e = os.WriteFile(filepath.Join(artifacts, "api-post-counts.json"), countJSON, 0600); e != nil {
		t.Fatal(e)
	}
	t.Logf("NETWORK_ROLE_BROWSER_PASS JWE_RS256_JWKS=true root_create_and_branch_close_response_loss_GET_only=true bootstrap_structure_agreement_branch_lifecycle=true backend_POSTs=13 artifacts=%s", artifacts)
}
````

### FILE: `internal/platform/postgres/network_role_integration_test.go`

```yaml
block_id: "GO-NETWORK-ROLE-COMPOSITION:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7834314f813337babff61a5ddae4f5997423b9f4ce33c65c8bd52b90b1e7b562"
variables: []
secrets_allowed: false
```

````go
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
````

## 6. Configuration surface

docs/NETWORK_ROLE_REFERENCE.md. NETWORK_ROLE_ENABLED optional guard and features.network_portal opt-in. Existing network:admin/franchise:write and organization scope. Root authority remains existing wildcard permission. Explicit terms/territory/dates; org.active is not production readiness.

## 7. Dependency bill

21new and6changed AUTHORED glue/fixtures/config. No new source/dependency/runtime. Migration0077 preserves evidence with immutable receipt and guarded downgrade. Source license/notices unchanged.

## 8. Apply order

GO-FULFILLMENT-SERVICE-FRANCHISE-API0.6.0, GO-ELECTROMOBILITY-APPLICATION1.17.0, GO-HUMAN-APPROVAL-CORE and identity required. Browser fixture reuses catalog synthetic issuer. Companion TS-NETWORK-ROLE-PORTAL owns Go/TS goldens. Full reference required for connected proof.

## 9. Verification

Atomic proof1new3replay,4table late-failure rollback, actor/hash/tenant/org/stale/territory/parent-close guards and historical/current results. Actual Next/BFF/Go/PG/Chromium13POST create/activate root/franchise/agreement/branch, suspend/resume and close in order. Two lost responses GET after reload,0extraPOST. Five host guards, empty downgrade/reapply and populated downgrade refusal.10web tests,4Go/TSgoldens, Next/types, fuzz2s3seeds286419executions. Desktop390px screenshots inspected.

## 10. Reconstruction evidence

NETWORK_ROLE_RELEASE_V402.md/json includes source SQL hash correspondence, RED/PASS and exact profile rebuilds. FAIL861 template extraction,862 qualified types,863 retained SQL lowercase-code fixture corrected without weakening owner guards. WholeT2804 and J5 identity/readiness remain open.


### FILE: `internal/platform/postgres/identity_j5_integration_test.go`

```yaml
block_id: "J5-EXTENSION-GO_NETWORK_ROLE_COMPOSITION:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e77e2d4799a27ed5c0add5bdd4b7a19a0013006d9d372eda83e64940da4b81e5"
variables: []
secrets_allowed: false
```

````go
package postgres_test

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	nr "elite.local/enterprise/internal/networkrole"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func j5RoleIssuer(t *testing.T) (identity.Verifier, func(string, string, []string, []string, bool) string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	encode := base64.RawURLEncoding.EncodeToString
	var issuer *httptest.Server
	issuer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			_ = json.NewEncoder(w).Encode(map[string]any{"issuer": issuer.URL, "jwks_uri": issuer.URL + "/keys", "id_token_signing_alg_values_supported": []string{"RS256"}})
		case "/keys":
			_ = json.NewEncoder(w).Encode(map[string]any{"keys": []any{map[string]any{"kty": "RSA", "use": "sig", "alg": "RS256", "kid": "local-confirmation", "n": encode(key.N.Bytes()), "e": encode(big.NewInt(int64(key.E)).Bytes())}}})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(issuer.Close)
	verifier, err := identity.NewOIDCVerifier(context.Background(), issuer.URL, "confirmation-api")
	if err != nil {
		t.Fatal(err)
	}
	token := func(subject, tenant string, permissions, organizations []string, expired bool) string {
		header, err := json.Marshal(map[string]string{"alg": "RS256", "kid": "local-confirmation", "typ": "JWT"})
		if err != nil {
			t.Fatal(err)
		}
		issued, expires := time.Now(), time.Now().Add(120*time.Second)
		if expired {
			issued = time.Now().Add(-300 * time.Second)
			expires = time.Now().Add(-120 * time.Second)
		}
		claims, err := json.Marshal(map[string]any{"iss": issuer.URL, "aud": "confirmation-api", "sub": subject, "iat": issued.Unix(), "exp": expires.Unix(), "tenant_id": tenant, "permissions": permissions, "organization_ids": organizations})
		if err != nil {
			t.Fatal(err)
		}
		unsigned := encode(header) + "." + encode(claims)
		digest := sha256.Sum256([]byte(unsigned))
		signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
		if err != nil {
			t.Fatal(err)
		}
		return unsigned + "." + encode(signature)
	}
	return verifier, token
}

func TestJ5IdentityScopedBootstrap(t *testing.T) {
	dsn := os.Getenv("PORTAL_SESSION_DB_URL")
	if dsn == "" {
		t.Skip("explicit owned PG fixture required")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := uuid.NewString()
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,$2,'Synthetic','Synthetic')`, tenant, "j5-"+tenant); err != nil {
		t.Fatal(err)
	}
	verifier, mint := j5RoleIssuer(t)
	store, err := db.NewNetworkRole(pool)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	httpapi.NetworkRoleModule{Service: store}.Register(mux, verifier)
	api := httptest.NewServer(mux)
	defer api.Close()
	post := func(token string, command nr.Command, status int) nr.Receipt {
		t.Helper()
		raw, _ := json.Marshal(command)
		request, _ := http.NewRequest("POST", api.URL+"/v1/franchise/network/commands", bytes.NewReader(raw))
		request.Header.Set("Authorization", "Bearer "+token)
		request.Header.Set("Content-Type", "application/json")
		response, err := api.Client().Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		if response.StatusCode != status {
			t.Fatalf("%s expected%d got%d %s", command.CommandID, status, response.StatusCode, body)
		}
		var receipt nr.Receipt
		if status == 200 && json.Unmarshal(body, &receipt) != nil {
			t.Fatal("receipt")
		}
		return receipt
	}
	get := func(token, path string, status int) {
		t.Helper()
		request, _ := http.NewRequest("GET", api.URL+path, nil)
		request.Header.Set("Authorization", "Bearer "+token)
		response, err := api.Client().Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
		if response.StatusCode != status {
			t.Fatal("read authority", status, response.StatusCode)
		}
	}
	admin := mint("initial-admin", tenant, []string{"network:admin", "network:bootstrap"}, []string{"root"}, false)
	root := nr.Command{CommandID: "j5-root-create", Action: "create-organization", EntityID: "root", Code: "root", DisplayName: "Synthetic root", Type: "franchisor"}
	post(mint("unprivileged", tenant, []string{"network:bootstrap"}, []string{"root"}, false), root, 403)
	post(mint("normal-admin", tenant, []string{"network:admin"}, []string{"root"}, false), root, 403)
	receipt := post(admin, root, 200)
	if receipt.Actor != "initial-admin" || receipt.Entity.ID != "root" {
		t.Fatal("issuer actor not retained")
	}
	current := mint("initial-admin", tenant, []string{"network:admin"}, []string{"root"}, false)
	activation := nr.Command{CommandID: "j5-root-active", Action: "transition-organization", ScopeOrganizationID: "root", EntityID: "root", Current: "provisioning", Target: "active", Version: "1"}
	post(current, activation, 200)
	get(current, "/v1/franchise/network/entities/organization/root?scope_organization_id=root", 200)
	other := root
	other.CommandID = "forbidden-extra-root"
	other.EntityID = "other"
	other.Code = "other"
	post(current, other, 403)
	scoped := nr.Command{CommandID: "j5-branch-create", Action: "create-organization", ScopeOrganizationID: "root", EntityID: "branch", Code: "branch", DisplayName: "Synthetic branch", Type: "store"}
	post(current, scoped, 200)
	foreignScope := mint("initial-admin", tenant, []string{"network:admin"}, []string{"foreign"}, false)
	post(foreignScope, activation, 403)
	withdrawn := mint("initial-admin", tenant, []string{}, []string{"root"}, false)
	post(withdrawn, activation, 403)
	get(withdrawn, "/v1/franchise/network/entities/organization/root?scope_organization_id=root", 403)
	expired := mint("initial-admin", tenant, []string{"network:admin"}, []string{"root"}, true)
	post(expired, activation, 401)
	post(current+"bad", activation, 401)
	get(mint("other-admin", tenant, []string{"network:admin"}, []string{"root"}, false), "/v1/franchise/network/commands/j5-root-active?scope_organization_id=root", 404)
	get(mint("initial-admin", uuid.NewString(), []string{"network:admin", "network:bootstrap"}, []string{"root"}, false), "/v1/franchise/network/commands/j5-root-create", 404)
	var organizations, receipts, events int
	err = pool.QueryRow(ctx, `select (select count(*)from org.organization where tenant_id=$1),(select count(*)from franchise.network_command_receipt where tenant_id=$1),(select count(*)from platform.outbox_event where tenant_id=$1)`, tenant).Scan(&organizations, &receipts, &events)
	if err != nil || organizations != 2 || receipts != 3 || events != 3 {
		t.Fatal("denial mutated state", organizations, receipts, events, err)
	}
	t.Log("J5_IDENTITY_PASS signed_OIDC=true wildcard=false bootstrap_scoped=true explicit_role_removal=true expiry_and_tamper_rejected=true tenant_actor_isolated=true organizations=2 immutable_receipts=3 outbox=3; existing bearer remains valid until fixed expiry, no immediate live IdP revocation claim")
}
````


V402 composed delta: Explicit network:admin + network:bootstrap, original scoped ongoing access; IDENTITY_J5_RELEASE_V402.md/json. IdP owns role grants and deprovisioning.
