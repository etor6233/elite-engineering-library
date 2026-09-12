# Go Meta Page Publishing Infrastructure

## 1. Metadata

```yaml
pack_id: "GO-META-PAGE-PUBLISHING-INFRASTRUCTURE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Governed Facebook Page text publishing infrastructure: exact human approval, durable scheduled jobs/fence, official SDK bridge, GET-only unknown recovery, public receipt and separately approved revoke."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "CPython 3.14 Windows x86-64 / Meta SDK26.0.1"]
compatible_with: ["GO-HUMAN-APPROVAL-CORE 0.2.0", "PYTHON-OFFICIAL-META-PAGE-WRITE-ADAPTER 0.1.0", "existing jobs/outbound/identity owners"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://github.com/facebook/facebook-python-business-sdk/tree/788f363d15b1269ab5efb7cd00fb5e3b133cd99b"]
verified_at: "2026-09-11"
```

## 2. Applicability

Requires the unchanged seven-file Meta Page adapter, human approval0.2.0, backend/identity/jobs/outbound owners. Local/synthetic reference proven. Actual account/Page/app access remains an activation input; broader social networks/media are not claimed.

## 3. Architecture contract

# Facebook Page publication infrastructure

This materializable scope publishes exact reviewed text to one configured
Facebook Page. It includes a runnable Go API/worker, immutable human approvals,
durable scheduled jobs, a single-write outbound fence, the official Meta Python
SDK, exact GET observation, public receipts and separately approved revocation.
No marketing policy, moderation, audience inference or eligibility is generated.

## Owners and provenance

The existing approval.Registry owns human separation; approval.request/decision
remain the only approval ledger. The new typed binding supplies zero automatic
approval policy. Its three manual kinds are social_publish, social_revoke and
whatsapp_reply. The generic PostgreSQL owner is shared with WhatsApp.
platform.job owns scheduling, bounded retry budgets and claim generations.
communication.outbound_delivery owns the at-most-one write attempt. Public
provider references/receipts live in platform.outbox_event with immutable social
evidence fields. No social ledger or replacement financial logic is introduced.

All new files are AUTHORED configuration/serialization/authorization/persistence
and process glue. The seven-file Meta SDK adapter is an unchanged dependency.
Provider POST/GET/DELETE behavior is DEPENDENCY_PIN facebook-business26.0.1,
fixed commit788f363d15b1269ab5efb7cd00fb5e3b133cd99b / Graphv26.0. The exact
LicenseRef-Meta-Platform bytes remain owned by that adapter. Nothing here is
attributed to a vendor as a locally adapted marketing algorithm.

## Start

1. Compose the admitted backend, jobs, outbound fence, human approval, seven-file
   Meta Page adapter and this publishing integration. Apply migrations through
   0058 plus0060 (0059 belongs to the independently selected WhatsApp delta).
2. Install the exact18-wheel requirements-windows-py314.lock into an isolated
   CPython3.14 Windows environment using --require-hashes --no-deps. Keep all
   source and license notices supplied by its dependency pack.
3. Review a profile derived from deploy/social/profile.reference.json, replacing
   the clearly synthetic tenant/organization/Page identifiers during deployment.
   Select its SHA256 independently. All worker replicas use the same bytes/hash.
   Queue/timing/attempt bounds are technical configuration, not marketing policy.
4. Supply the empty environment inputs from deploy/social/.env.example when
   activating. No secrets are requested or embedded now. Pin the materialized
   meta_page_write/bridge.py SHA256 and isolated Python executable path.
5. Build `go build -mod=readonly ./cmd/social-publishing` and run the resulting
   executable. It serves on loopback127.0.0.1:8097 by default and runs the worker.
   OIDC verifier owns identity; the caller needs the configured organization and
   social:request/social:approve/social:read/social:reconcile permissions as used.
   Front a non-loopback deployment with the target's admitted TLS/auth gateway.

Before serving, startup rejects missing profile/hash, credentials, runtime,
provider SDK source drift, dependency version drift and a changed Python bridge.
The SDK preflight is offline and verifies fixed source bytes; it does not prove
Page ownership, Meta app access approval or live token permissions.

## Operator contract

POST /v1/social/requests submits the Request JSON documented by
internal/socialbridge/contract.go. Server-owned tenant/Page/profile checks reject
cross-scope values. ApprovalID identifies the immutable request; DeliveryKey is
social_ plus SHA256(tenant NUL page NUL approvalID NUL operation). ContentSHA256
binds the exact UTF-8 message. ScheduledAt and ExpiresAt are explicit UTC instants.
The response provides the canonical RequestSHA256 for exact review.

GET /v1/social/requests/{id} exposes the exact pending payload/hash and current
approval, delivery and receipt in the authorized tenant/organization.
POST /v1/social/requests/{id}/decision requires request_sha256, explicit boolean
approve and reason. Its verified human must differ from the requester. Only the
approved transaction creates the scheduled job. Rejection creates no job.

The worker revalidates durable approval, tenant, profile, page, content,
not-before/expiry and current claim generation in the fence transaction. A replay
cannot obtain another write attempt. Its Python child calls the real SDK POST
once, then GET verifies exact id/from.id/message/is_published before success.

After an uncertain response, known post IDs are persisted and retries perform
GET only. POST /v1/social/requests/{id}/reconcile provides the same GET-only
recovery after the queue budget is exhausted. It preserves the attempt number.
A saved receipt recovers a crash before fence/queue acknowledgement without
another SDK call. Restarting the host preserves approvals, jobs and evidence.

Revocation is another request, operation=revoke, separately reviewed. It names
the original approved publication and its exact public receipt/postID/content.
The worker verifies that binding, then the SDK performs GET and DELETE once.
A different Page, post reference, message or original approval cannot authorize
deletion. A rejected pending request never dispatches. Disabling an action in a
new reviewed profile prevents new admission under that profile; do not assume
changing a file can recall a provider operation already in flight.

## Ambiguous external outcomes

A lost POST response with no provider ID remains UNKNOWN; neither the operator
endpoint nor a restarted worker invents a post ID or resends it. A lost DELETE
acknowledgement also remains UNKNOWN: absence or permission denial cannot prove
deletion. These are explicit preserved external outcomes, not successful sends
or missing implementation. Retain the request and seek provider evidence through
the operator process. No credential alone is claimed to resolve such a case.

## Retention, rollback and boundary

Planned publication text is stored in the access-controlled approval payload;
tokens are never stored there. Provider receipts contain public IDs, hashes and
state. Do not send private customer records as publication text. Target retention
configuration must preserve approval/receipt evidence needed for reconciliation.
Stop admission and drain workers before changing executable/profile versions.
The down migrations refuse removal once their governed evidence exists; use a
reviewed forward change rather than deleting that evidence.

Focused PostgreSQL/real-SDK-fixture tests cover exact approval, replay, scheduled
admission, tenant/Page/content isolation, stale claims, one POST, GET-only
recovery, terminal recovery, receipt crash, expiry under a lock and governed
DELETE. These library results do not certify live app/Page permission, target
load/availability, broader social networks or a fresh composition-wide security
review. The parent composition owns its T2803 SCA/SAST and final admission.


## 4. Exact file manifest

```text
CREATE internal/socialbridge/contract.go
CREATE internal/socialbridge/process.go
CREATE internal/socialbridge/contract_test.go
CREATE internal/platform/postgres/social_publishing.go
CREATE internal/platform/postgres/social_publishing_integration_test.go
CREATE internal/platform/httpapi/social_publish.go
CREATE internal/platform/httpapi/social_publish_test.go
CREATE cmd/social-publishing/main.go
CREATE meta_page_write/bridge.py
CREATE db/migrations/0060_social_publish.up.sql
CREATE db/migrations/0060_social_publish.down.sql
CREATE deploy/social/profile.reference.json
CREATE deploy/social/.env.example
CREATE docs/SOCIAL_PUBLISHING.md
```

## 5. Materialization blocks

### FILE: `internal/socialbridge/contract.go`

```yaml
block_id: "GO-META-PAGE-PUBLISHING-INFRASTRUCTURE:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d4eb5f42a4b8285f72eed38c304c8acda2cf818b94d5b35594103f838e732b2f"
variables: []
secrets_allowed: false
```

````go
// Package socialbridge is AUTHORED integration glue. It binds an explicit
// deployment profile and reviewed bytes to existing approval/jobs/fence owners.
// Meta owns provider behavior; no marketing policy or audience logic is inferred.
package socialbridge

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/platform/identity"
)

var ErrBinding = errors.New("social: invalid or unauthorized binding")
var ErrUnknown = errors.New("social: remote outcome unresolved; no repeated write permitted")
var hashRE = regexp.MustCompile(`^[0-9a-f]{64}$`)
var idRE = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)
var pageRE = regexp.MustCompile(`^[1-9][0-9]{0,31}$`)

type Config struct {
	Schema         string `json:"schema"`
	TenantID       string `json:"tenant_id"`
	OrganizationID string `json:"organization_id"`
	PageID         string `json:"page_id"`
	Queue          string `json:"queue"`
	Publish        bool   `json:"publish"`
	Revoke         bool   `json:"revoke"`
	Review         string `json:"review"`
	LeaseSeconds   int    `json:"lease_seconds"`
	RetrySeconds   int    `json:"retry_seconds"`
	PollSeconds    int    `json:"poll_seconds"`
	MaxAttempts    int    `json:"max_attempts"`
}
type Profile struct {
	config Config
	hash   string
}

func Hash(raw []byte) string { v := sha256.Sum256(raw); return hex.EncodeToString(v[:]) }
func Decode(raw []byte, into any) error {
	if len(raw) == 0 || len(raw) > 32768 || !utf8.Valid(raw) {
		return ErrBinding
	}
	// Duplicate/case-alias rejection is protocol glue, not a business policy.
	d := json.NewDecoder(bytes.NewReader(raw))
	if unique(d, 0) != nil {
		return ErrBinding
	}
	if _, e := d.Token(); e != io.EOF {
		return ErrBinding
	}
	d = json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(into) != nil {
		return ErrBinding
	}
	return nil
}
func unique(d *json.Decoder, depth int) error {
	if depth > 16 {
		return ErrBinding
	}
	t, e := d.Token()
	if e != nil {
		return e
	}
	v, ok := t.(json.Delim)
	if !ok {
		return nil
	}
	switch v {
	case '{':
		seen := map[string]bool{}
		for d.More() {
			k, e := d.Token()
			if e != nil {
				return e
			}
			s, ok := k.(string)
			if !ok {
				return ErrBinding
			}
			s = strings.ToLower(s)
			if seen[s] {
				return ErrBinding
			}
			seen[s] = true
			if unique(d, depth+1) != nil {
				return ErrBinding
			}
		}
	case '[':
		for d.More() {
			if unique(d, depth+1) != nil {
				return ErrBinding
			}
		}
	default:
		return ErrBinding
	}
	_, e = d.Token()
	return e
}
func Load(raw []byte, expected string) (*Profile, error) {
	var c Config
	if !hashRE.MatchString(expected) || Hash(raw) != expected || Decode(raw, &c) != nil {
		return nil, ErrBinding
	}
	if c.Schema != "elite.meta-page-publishing.v1" || !idRE.MatchString(c.TenantID) || !idRE.MatchString(c.OrganizationID) || !pageRE.MatchString(c.PageID) || !idRE.MatchString(c.Queue) || c.Review != "one_distinct_human" || c.LeaseSeconds < 60 || c.LeaseSeconds > 600 || c.RetrySeconds < 1 || c.RetrySeconds > 3600 || c.PollSeconds < 1 || c.PollSeconds > 60 || c.MaxAttempts < 2 || c.MaxAttempts > 100 {
		return nil, ErrBinding
	}
	return &Profile{c, expected}, nil
}
func LoadFile(path, hash string) (*Profile, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, ErrBinding
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, 32769))
	if e != nil {
		return nil, ErrBinding
	}
	return Load(b, hash)
}
func (p *Profile) Config() Config { return p.config }
func (p *Profile) SHA256() string { return p.hash }
func (p *Profile) Authorize(a identity.Principal, permission string) bool {
	return p != nil && a.Subject != "" && len(a.Subject) <= 128 && a.TenantID == p.config.TenantID && a.AllowedOrganization(p.config.OrganizationID) && a.Allowed(permission)
}

type Intent struct {
	TenantID      string `json:"tenant_id"`
	PageID        string `json:"page_id"`
	ProfileSHA256 string `json:"profile_sha256"`
	ApprovalID    string `json:"approval_id"`
	DeliveryKey   string `json:"delivery_key"`
	Operation     string `json:"operation"`
	Message       string `json:"message"`
	ContentSHA256 string `json:"content_sha256"`
}
type Request struct {
	Intent             Intent    `json:"intent"`
	ScheduledAt        time.Time `json:"scheduled_at"`
	ExpiresAt          time.Time `json:"expires_at"`
	OriginalApprovalID string    `json:"original_approval_id"`
	ProviderReference  string    `json:"provider_reference"`
}

func (p *Profile) Validate(r Request) error {
	i := r.Intent
	c := p.config
	if i.TenantID != c.TenantID || i.PageID != c.PageID || i.ProfileSHA256 != p.hash || !idRE.MatchString(i.ApprovalID) || i.DeliveryKey != DeliveryKey(i) || !utf8.ValidString(i.Message) || len(i.Message) > 16384 || strings.TrimSpace(i.Message) == "" || Hash([]byte(i.Message)) != i.ContentSHA256 || r.ScheduledAt.IsZero() || !r.ExpiresAt.After(r.ScheduledAt) {
		return ErrBinding
	}
	if i.Operation == "publish" {
		if !c.Publish || r.OriginalApprovalID != "" || r.ProviderReference != "" {
			return ErrBinding
		}
	} else if i.Operation == "revoke" {
		if !c.Revoke || !idRE.MatchString(r.OriginalApprovalID) || !ValidReference(c.PageID, r.ProviderReference) {
			return ErrBinding
		}
	} else {
		return ErrBinding
	}
	return nil
}
func DeliveryKey(i Intent) string {
	return "social_" + Hash([]byte(i.TenantID+"\x00"+i.PageID+"\x00"+i.ApprovalID+"\x00"+i.Operation))
}
func ValidReference(page, ref string) bool {
	return regexp.MustCompile(`^` + regexp.QuoteMeta(page) + `_[1-9][0-9]{0,63}$`).MatchString(ref)
}
func (r Request) Hash() string { b, _ := json.Marshal(r); return Hash(b) }
func (r Request) Message() channels.Message {
	return channels.Message{ChannelCode: "facebook_pages", TenantID: r.Intent.TenantID, ExternalID: r.Intent.PageID, ThreadID: r.Hash(), DeliveryKey: r.Intent.DeliveryKey, Direction: channels.DirectionOut, Text: r.Intent.Message}
}

type Receipt struct {
	TenantID          string `json:"tenant_id"`
	PageID            string `json:"page_id"`
	ProfileSHA256     string `json:"profile_sha256"`
	ApprovalID        string `json:"approval_id"`
	DeliveryKey       string `json:"delivery_key"`
	ProviderReference string `json:"provider_reference"`
	State             string `json:"state"`
	ContentSHA256     string `json:"content_sha256"`
	EvidenceSHA256    string `json:"evidence_sha256"`
}

func (r Receipt) Validate(i Intent) error {
	want := "published"
	if i.Operation == "revoke" {
		want = "revoked"
	}
	if r.TenantID != i.TenantID || r.PageID != i.PageID || r.ProfileSHA256 != i.ProfileSHA256 || r.ApprovalID != i.ApprovalID || r.DeliveryKey != i.DeliveryKey || r.ContentSHA256 != i.ContentSHA256 || r.State != want || !ValidReference(i.PageID, r.ProviderReference) || !hashRE.MatchString(r.EvidenceSHA256) {
		return ErrBinding
	}
	return nil
}

type Result struct {
	Receipt           *Receipt `json:"receipt,omitempty"`
	Unknown           bool     `json:"unknown,omitempty"`
	ProviderReference string   `json:"provider_reference,omitempty"`
	Code              string   `json:"code,omitempty"`
}
type Adapter interface {
	Execute(context.Context, Request, bool) (Result, error)
}
````

### FILE: `internal/socialbridge/process.go`

```yaml
block_id: "GO-META-PAGE-PUBLISHING-INFRASTRUCTURE:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "4d4c7defc1c0f9532cc2e686f325aa1aea7a9fb295c5c88c1c167cc04ef006a8"
variables: []
secrets_allowed: false
```

````go
package socialbridge

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Process invokes only exact materialized code. Credentials are environment
// values, never command arguments, receipts, logs or approval payloads.
type Process struct {
	Python, Script, ScriptSHA256 string
	Environment                  []string
}

func (p Process) Preflight(ctx context.Context) error {
	if p.Validate() != nil {
		return ErrBinding
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, p.Python, "-I", "-B", p.Script, "--preflight")
	for _, k := range []string{"SYSTEMROOT", "WINDIR", "TEMP", "TMP"} {
		if v := os.Getenv(k); v != "" {
			c.Env = append(c.Env, k+"="+v)
		}
	}
	for _, v := range p.Environment {
		k, _, ok := strings.Cut(v, "=")
		if !ok || (k != "META_APP_ID" && k != "META_APP_SECRET" && k != "META_PAGE_ACCESS_TOKEN") {
			return ErrBinding
		}
		c.Env = append(c.Env, v)
	}
	var out boundedOutput
	c.Stdout = &out
	c.Stderr = io.Discard
	if c.Run() != nil || out.exceeded || out.String() != `{"state":"PASS"}` {
		return ErrBinding
	}
	return nil
}

func (p Process) Validate() error {
	if !filepath.IsAbs(p.Python) || !filepath.IsAbs(p.Script) {
		return ErrBinding
	}
	if _, e := os.Stat(p.Python); e != nil {
		return ErrBinding
	}
	b, e := os.ReadFile(p.Script)
	if e != nil || !hashRE.MatchString(p.ScriptSHA256) || Hash(b) != p.ScriptSHA256 {
		return ErrBinding
	}
	return nil
}

type boundedOutput struct {
	bytes.Buffer
	exceeded bool
}

func (b *boundedOutput) Write(v []byte) (int, error) {
	if b.Len()+len(v) > 32768 {
		b.exceeded = true
		return 0, io.ErrShortBuffer
	}
	return b.Buffer.Write(v)
}
func (p Process) Execute(ctx context.Context, r Request, reconcile bool) (Result, error) {
	if p.Validate() != nil {
		return Result{}, ErrBinding
	}
	input, _ := json.Marshal(struct {
		Request   Request `json:"request"`
		Reconcile bool    `json:"reconcile"`
	}{r, reconcile})
	ctx, cancel := context.WithTimeout(ctx, 35*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, p.Python, "-I", "-B", p.Script)
	c.Stdin = bytes.NewReader(input)
	// Explicit allowlist. Proxy/SSL overrides, PYTHONPATH, debug and arbitrary
	// inherited service credentials never cross this child process boundary.
	for _, k := range []string{"SYSTEMROOT", "WINDIR", "TEMP", "TMP"} {
		if v := os.Getenv(k); v != "" {
			c.Env = append(c.Env, k+"="+v)
		}
	}
	for _, v := range p.Environment {
		k, _, ok := strings.Cut(v, "=")
		if !ok || (k != "META_APP_ID" && k != "META_APP_SECRET" && k != "META_PAGE_ACCESS_TOKEN") {
			return Result{}, ErrBinding
		}
		c.Env = append(c.Env, v)
	}
	var out boundedOutput
	c.Stdout = &out
	c.Stderr = io.Discard
	if c.Run() != nil || out.exceeded {
		return Result{Unknown: true, Code: "PROCESS_RESULT_UNCONFIRMED"}, nil
	}
	var result Result
	if Decode(out.Bytes(), &result) != nil {
		return Result{Unknown: true, Code: "PROCESS_RESULT_INVALID"}, nil
	}
	if result.Receipt != nil {
		if result.Unknown || result.Receipt.Validate(r.Intent) != nil {
			return Result{Unknown: true, Code: "PROCESS_RECEIPT_INVALID"}, nil
		}
	} else if !result.Unknown {
		return Result{Unknown: true, Code: "PROCESS_RESULT_INVALID"}, nil
	}
	if result.ProviderReference != "" && !ValidReference(r.Intent.PageID, result.ProviderReference) {
		return Result{Unknown: true, Code: "PROCESS_REFERENCE_INVALID"}, nil
	}
	return result, nil
}
````

### FILE: `internal/socialbridge/contract_test.go`

```yaml
block_id: "GO-META-PAGE-PUBLISHING-INFRASTRUCTURE:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "efb8dcedd6b325eadd72fa7833c30003aa6015e995f21d8ec52a837bf83e3b00"
variables: []
secrets_allowed: false
```

````go
package socialbridge

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func testProfile(t *testing.T) *Profile {
	t.Helper()
	raw, _ := json.Marshal(Config{Schema: "elite.meta-page-publishing.v1", TenantID: "tenant", OrganizationID: "organization", PageID: "123", Queue: "social", Publish: true, Revoke: true, Review: "one_distinct_human", LeaseSeconds: 60, RetrySeconds: 1, PollSeconds: 1, MaxAttempts: 2})
	p, e := Load(raw, Hash(raw))
	if e != nil {
		t.Fatal(e)
	}
	return p
}
func TestSocialProfileAdmissionAndChange(t *testing.T) {
	p := testProfile(t)
	c := p.Config()
	c.PollSeconds = 5
	c.Queue = "different_queue"
	raw, _ := json.Marshal(c)
	changed, e := Load(raw, Hash(raw))
	if e != nil || changed.SHA256() == p.SHA256() || changed.Config().PollSeconds != 5 {
		t.Fatal("profile configuration not admitted")
	}
	if _, e = Load(raw, p.SHA256()); e == nil {
		t.Fatal("hash drift")
	}
	c.Review = "automatic"
	raw, _ = json.Marshal(c)
	if _, e = Load(raw, Hash(raw)); e == nil {
		t.Fatal("unimplemented policy admitted")
	}
	for _, raw := range []string{`{"schema":"a","schema":"b"}`, `{"queue":"x","QUEUE":"y"}`, `{} {}`, `{"unknown":true}`, `null`} {
		if _, e = Load([]byte(raw), Hash([]byte(raw))); e == nil {
			t.Fatal("bad profile", raw)
		}
	}
}
func TestSocialProcessPreflight(t *testing.T) {
	python := os.Getenv("SOCIAL_TEST_PYTHON")
	if python == "" {
		t.Skip("locked Python runtime required")
	}
	script, e := filepath.Abs("../../meta_page_write/bridge.py")
	if e != nil {
		t.Fatal(e)
	}
	raw, e := os.ReadFile(script)
	if e != nil {
		t.Fatal(e)
	}
	p := Process{Python: python, Script: script, ScriptSHA256: Hash(raw), Environment: []string{"META_APP_ID=321", "META_APP_SECRET=fixture-secret", "META_PAGE_ACCESS_TOKEN=fixture-token"}}
	if e = p.Preflight(context.Background()); e != nil {
		t.Fatal("offline exact SDK preflight", e)
	}
	// A changed lock must not remove the installed-version checks silently.
	copyRoot := filepath.Join(t.TempDir(), "meta_page_write")
	if e = os.Mkdir(copyRoot, 0700); e != nil {
		t.Fatal(e)
	}
	for _, name := range []string{"bridge.py", "adapter.py", "requirements-windows-py314.lock"} {
		b, e := os.ReadFile(filepath.Join(filepath.Dir(script), name))
		if e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(filepath.Join(copyRoot, name), b, 0600); e != nil {
			t.Fatal(e)
		}
	}
	changed := p
	changed.Script = filepath.Join(copyRoot, "bridge.py")
	if e = os.WriteFile(filepath.Join(copyRoot, "requirements-windows-py314.lock"), []byte("# removed dependencies\n"), 0600); e != nil {
		t.Fatal(e)
	}
	if e = changed.Preflight(context.Background()); e == nil {
		t.Fatal("changed dependency lock disabled preflight checks")
	}
	p.Environment = nil
	if e = p.Preflight(context.Background()); e == nil {
		t.Fatal("missing credentials accepted")
	}
	p.ScriptSHA256 = Hash([]byte("different"))
	if e = p.Validate(); e == nil {
		t.Fatal("script drift accepted")
	}
}
func FuzzSocialProfileAndPayload(f *testing.F) {
	for _, s := range []string{`{"a":1}`, `{"a":{"b":[1,null,true]}}`, `{"a":1,"A":2}`, `null`, `{} {}`} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		var v any
		if Decode(raw, &v) == nil {
			canonical, e := json.Marshal(v)
			if e != nil {
				t.Fatal(e)
			}
			var again any
			if Decode(canonical, &again) != nil {
				t.Fatal("admitted JSON failed canonical decode")
			}
		}
		if p, e := Load(raw, Hash(raw)); e == nil && p.SHA256() != Hash(raw) {
			t.Fatal("unbound profile")
		}
	})
}
````

### FILE: `internal/platform/postgres/social_publishing.go`

```yaml
block_id: "GO-META-PAGE-PUBLISHING-INFRASTRUCTURE:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "73b70a28c85cf512beef8052863aa728ea65810fd540dcc88d931ef603343478"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED binding glue. Request/decision, scheduling, single-write fence and
// public provider observations retain their existing shared persistence owners.
import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/socialbridge"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SocialPublishing struct {
	pool      *pgxpool.Pool
	profile   *socialbridge.Profile
	approvals *HumanApprovals
	jobs      *Jobs
	fence     *OutboundDeliveryStore
	adapter   socialbridge.Adapter
}
type socialJob struct {
	ApprovalID    string `json:"approval_id"`
	ProfileSHA256 string `json:"profile_sha256"`
	PageID        string `json:"page_id"`
	RequestSHA256 string `json:"request_sha256"`
}

func NewSocialPublishing(pool *pgxpool.Pool, p *socialbridge.Profile, key []byte, a socialbridge.Adapter) (*SocialPublishing, error) {
	if p == nil || a == nil {
		return nil, socialbridge.ErrBinding
	}
	f, e := NewOutboundDeliveryStore(pool, key, time.Duration(p.Config().LeaseSeconds)*time.Second)
	if e != nil {
		return nil, e
	}
	return &SocialPublishing{pool, p, NewHumanApprovals(pool), NewJobs(pool), f, a}, nil
}
func socialKind(r socialbridge.Request) approval.Kind {
	if r.Intent.Operation == "revoke" {
		return approval.KindSocialRevoke
	}
	return approval.KindSocialPublish
}
func socialJobID(tenant, id string) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte("elite.social.job\x00"+tenant+"\x00"+id)).String()
}
func (s *SocialPublishing) request(ctx context.Context, id string) (socialbridge.Request, string, error) {
	var raw []byte
	var hash string
	c := s.profile.Config()
	e := s.pool.QueryRow(ctx, `select payload,evidence_sha from approval.request where tenant_id=$1 and request_id=$2 and organization_id=$3 and kind in ('social_publish','social_revoke')`, c.TenantID, id, c.OrganizationID).Scan(&raw, &hash)
	var r socialbridge.Request
	if e != nil {
		return r, "", e
	}
	if socialbridge.Decode(raw, &r) != nil || s.profile.Validate(r) != nil || r.Intent.ApprovalID != id {
		return r, "", socialbridge.ErrBinding
	}
	return r, hash, nil
}
func (s *SocialPublishing) guard(ctx context.Context, tx pgx.Tx, r socialbridge.Request, needDue bool) error {
	if s.profile.Validate(r) != nil {
		return socialbridge.ErrBinding
	}
	c := s.profile.Config()
	var valid bool
	e := tx.QueryRow(ctx, `with locked as materialized(select status from platform.tenant where tenant_id=$1 for share)
 select status='active' and $2::timestamptz>clock_timestamp() and (not $3::boolean or $4::timestamptz<=clock_timestamp()) from locked`, c.TenantID, r.ExpiresAt, needDue, r.ScheduledAt).Scan(&valid)
	if e != nil || !valid {
		return socialbridge.ErrBinding
	}
	if r.Intent.Operation == "revoke" {
		// The exact publish receipt belongs to this page, original approved content
		// and tenant. A user-supplied post id alone never authorizes a DELETE.
		var raw []byte
		e = tx.QueryRow(ctx, `select payload from platform.outbox_event where tenant_id=$1 and aggregate_type='social_approval' and aggregate_id=$2 and event_type='social.receipt' and aggregate_version=1 for share`, c.TenantID, r.OriginalApprovalID).Scan(&raw)
		var receipt socialbridge.Receipt
		if e != nil || socialbridge.Decode(raw, &receipt) != nil || receipt.State != "published" || receipt.PageID != c.PageID || receipt.TenantID != c.TenantID || receipt.ProfileSHA256 != s.profile.SHA256() || receipt.ContentSHA256 != r.Intent.ContentSHA256 || receipt.ProviderReference != r.ProviderReference {
			return socialbridge.ErrBinding
		}
	}
	return nil
}
func (s *SocialPublishing) Submit(ctx context.Context, p identity.Principal, r socialbridge.Request) (string, bool, error) {
	if !s.profile.Authorize(p, "social:request") || s.profile.Validate(r) != nil {
		return "", false, socialbridge.ErrBinding
	}
	raw, _ := json.Marshal(r)
	payload, hash, e := approval.CanonicalPayload(raw)
	if e != nil {
		return "", false, e
	}
	spec := HumanApprovalSpec{Request: approval.Request{TenantID: p.TenantID, ID: r.Intent.ApprovalID, Kind: socialKind(r), SubjectID: r.Intent.PageID, Requester: p.Subject, EvidenceSHA: hash}, OrganizationID: s.profile.Config().OrganizationID, Payload: payload}
	replay, e := s.approvals.Submit(ctx, p, spec, "social:request", func(ctx context.Context, tx pgx.Tx) error { return s.guard(ctx, tx, r, false) })
	return hash, replay, e
}
func (s *SocialPublishing) Decide(ctx context.Context, p identity.Principal, id, hash string, approve bool, reason string) (approval.State, error) {
	if !s.profile.Authorize(p, "social:approve") {
		return "", socialbridge.ErrBinding
	}
	r, stored, e := s.request(ctx, id)
	if e != nil || stored != hash {
		return "", socialbridge.ErrBinding
	}
	c := s.profile.Config()
	return s.approvals.Decide(ctx, p, c.TenantID, id, c.OrganizationID, hash, approve, reason, "social:approve", func(ctx context.Context, tx pgx.Tx) error {
		if !approve {
			return nil
		}
		if e := s.guard(ctx, tx, r, false); e != nil {
			return e
		}
		payload, _ := json.Marshal(socialJob{id, s.profile.SHA256(), c.PageID, hash})
		_, e := tx.Exec(ctx, `insert into platform.job(tenant_id,job_id,queue,job_type,schema_version,payload,max_attempts,available_at)values($1,$2,$3,'social.publish',1,$4,$5,$6)`, c.TenantID, socialJobID(c.TenantID, id), c.Queue, payload, c.MaxAttempts, r.ScheduledAt)
		return e
	})
}

type SocialStatus struct {
	ApprovalID    string          `json:"approval_id"`
	State         string          `json:"approval_state"`
	RequestSHA256 string          `json:"request_sha256"`
	Payload       json.RawMessage `json:"payload"`
	DeliveryState string          `json:"delivery_state"`
	Receipt       json.RawMessage `json:"receipt,omitempty"`
	TerminalCode  string          `json:"terminal_code,omitempty"`
}

func (s *SocialPublishing) Status(ctx context.Context, p identity.Principal, id string) (SocialStatus, error) {
	var v SocialStatus
	if !s.profile.Authorize(p, "social:read") {
		return v, socialbridge.ErrBinding
	}
	c := s.profile.Config()
	e := s.pool.QueryRow(ctx, `select a.request_id,a.state,a.evidence_sha,a.payload,coalesce(d.state,''),e.payload,coalesce(j.terminal_error_code,'') from approval.request a
 left join communication.outbound_delivery d on d.tenant_id=a.tenant_id and d.channel_code='facebook_pages' and d.delivery_key=a.payload->'intent'->>'delivery_key'
 left join platform.outbox_event e on e.tenant_id=a.tenant_id and e.aggregate_type='social_approval' and e.aggregate_id=a.request_id and e.event_type='social.receipt' and e.aggregate_version=1
 left join platform.job j on j.tenant_id=a.tenant_id and j.job_id=$4
 where a.tenant_id=$1 and a.request_id=$2 and a.organization_id=$3 and a.kind in ('social_publish','social_revoke')`, c.TenantID, id, c.OrganizationID, socialJobID(c.TenantID, id)).Scan(&v.ApprovalID, &v.State, &v.RequestSHA256, &v.Payload, &v.DeliveryState, &v.Receipt, &v.TerminalCode)
	return v, e
}
func (s *SocialPublishing) Claim(ctx context.Context, worker string) ([]Job, error) {
	c := s.profile.Config()
	match, _ := json.Marshal(map[string]string{"profile_sha256": s.profile.SHA256(), "page_id": c.PageID})
	return s.jobs.ClaimScoped(ctx, c.Queue, worker, time.Duration(c.LeaseSeconds)*time.Second, 1, JobScope{TenantID: c.TenantID, JobType: "social.publish", SchemaVersion: 1, PayloadMatch: match})
}
func (s *SocialPublishing) event(ctx context.Context, r socialbridge.Request, kind string, payload any) error {
	raw, _ := json.Marshal(payload)
	id := uuid.NewSHA1(uuid.NameSpaceOID, []byte("elite.social.event\x00"+r.Intent.TenantID+"\x00"+r.Intent.ApprovalID+"\x00"+kind)).String()
	_, e := s.pool.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'social_approval',$3,1,$4,1,clock_timestamp(),$5) on conflict do nothing`, r.Intent.TenantID, id, r.Intent.ApprovalID, kind, raw)
	if e != nil {
		return e
	}
	var same bool
	e = s.pool.QueryRow(ctx, `select payload=$3::jsonb from platform.outbox_event where tenant_id=$1 and event_id=$2`, r.Intent.TenantID, id, raw).Scan(&same)
	if e != nil {
		return e
	}
	if !same {
		return socialbridge.ErrBinding
	}
	return nil
}
func (s *SocialPublishing) observation(ctx context.Context, r socialbridge.Request) (string, *socialbridge.Receipt, error) {
	var raw []byte
	e := s.pool.QueryRow(ctx, `select payload from platform.outbox_event where tenant_id=$1 and aggregate_type='social_approval' and aggregate_id=$2 and event_type='social.receipt' and aggregate_version=1`, r.Intent.TenantID, r.Intent.ApprovalID).Scan(&raw)
	if e == nil {
		var v socialbridge.Receipt
		if socialbridge.Decode(raw, &v) != nil || v.Validate(r.Intent) != nil {
			return "", nil, socialbridge.ErrBinding
		}
		return v.ProviderReference, &v, nil
	}
	if !errors.Is(e, pgx.ErrNoRows) {
		return "", nil, e
	}
	var ref string
	e = s.pool.QueryRow(ctx, `select payload->>'provider_reference' from platform.outbox_event where tenant_id=$1 and aggregate_type='social_approval' and aggregate_id=$2 and event_type='social.observation' and aggregate_version=1`, r.Intent.TenantID, r.Intent.ApprovalID).Scan(&ref)
	if errors.Is(e, pgx.ErrNoRows) {
		return "", nil, nil
	}
	return ref, nil, e
}
func (s *SocialPublishing) Process(ctx context.Context, job Job, worker string) error {
	var payload socialJob
	c := s.profile.Config()
	if job.TenantID != c.TenantID || job.JobType != "social.publish" || job.Queue != c.Queue || job.SchemaVersion != 1 || socialbridge.Decode(job.Payload, &payload) != nil || payload.ProfileSHA256 != s.profile.SHA256() || payload.PageID != c.PageID || job.JobID != socialJobID(c.TenantID, payload.ApprovalID) {
		return socialbridge.ErrBinding
	}
	r, hash, e := s.request(ctx, payload.ApprovalID)
	if e != nil || hash != payload.RequestSHA256 {
		return socialbridge.ErrBinding
	}
	message := r.Message()
	messageHash, e := outbounddelivery.MessageSHA256(message)
	if e != nil {
		return e
	}
	claim, e := s.fence.claimWithAdmission(ctx, message, messageHash, func(ctx context.Context, tx pgx.Tx) error {
		if _, e := LookupHumanApproval(ctx, tx, c.TenantID, payload.ApprovalID, c.OrganizationID, socialKind(r), hash); e != nil {
			return e
		}
		if e := s.guard(ctx, tx, r, true); e != nil {
			return e
		}
		// An obsolete job generation cannot acquire the remote write fence.
		var live bool
		e := tx.QueryRow(ctx, `with locked as materialized(select claimed_by,attempts,claimed_until,completed_at,terminal_error_code,payload from platform.job where tenant_id=$1 and job_id=$2 for update)
 select claimed_by=$3 and attempts=$4 and claimed_until>clock_timestamp() and completed_at is null and terminal_error_code is null and payload=$5::jsonb from locked`, job.TenantID, job.JobID, worker, job.Attempts, job.Payload).Scan(&live)
		if e != nil || !live {
			return ErrJobClaimLost
		}
		return nil
	})
	if e == nil && claim.Replay {
		return s.jobs.Complete(ctx, job.TenantID, job.JobID, worker, job.Attempts)
	}
	unknown := errors.Is(e, outbounddelivery.ErrUnknown)
	if e != nil && !unknown {
		return e
	}
	ref, receipt, e := s.observation(ctx, r)
	if e != nil {
		return e
	}
	if receipt == nil {
		// Every ambiguous delivery is read-only from here onward. No reference means
		// there is no safe GET identity; retain the case for provider evidence.
		if unknown && ref == "" {
			_, _ = s.jobs.Fail(ctx, job.TenantID, job.JobID, worker, "PROVIDER_REFERENCE_REQUIRED", job.Attempts, time.Duration(c.RetrySeconds)*time.Second)
			return socialbridge.ErrUnknown
		}
		call := r
		if unknown {
			call.ProviderReference = ref
		}
		result, callErr := s.adapter.Execute(ctx, call, unknown)
		if callErr != nil || result.Receipt == nil || result.Receipt.Validate(r.Intent) != nil {
			if result.ProviderReference != "" && socialbridge.ValidReference(c.PageID, result.ProviderReference) {
				if e = s.event(ctx, r, "social.observation", map[string]string{"provider_reference": result.ProviderReference, "request_sha256": hash}); e != nil {
					return e
				}
			}
			if !unknown {
				if e = s.fence.MarkUnknown(ctx, message, messageHash, "PROVIDER_OUTCOME_UNCONFIRMED"); e != nil {
					return e
				}
			}
			_, _ = s.jobs.Fail(ctx, job.TenantID, job.JobID, worker, "PROVIDER_OUTCOME_UNCONFIRMED", job.Attempts, time.Duration(c.RetrySeconds)*time.Second)
			return socialbridge.ErrUnknown
		}
		receipt = result.Receipt
		if e = s.event(ctx, r, "social.receipt", receipt); e != nil {
			return e
		}
	}
	durable := outbounddelivery.Receipt{ProviderMessageID: receipt.ProviderReference, EvidenceSHA256: receipt.EvidenceSHA256, AcceptedAt: time.Now().UTC()}
	if unknown {
		e = s.fence.ReconcileAccepted(ctx, message, messageHash, durable)
	} else {
		e = s.fence.Complete(ctx, message, messageHash, durable)
	}
	if e != nil {
		return e
	}
	return s.jobs.Complete(ctx, job.TenantID, job.JobID, worker, job.Attempts)
}

// Reconcile is an authorized GET-only recovery path even after a job's retry
// budget is exhausted. It never resets a generation or reissues POST/DELETE.
func (s *SocialPublishing) Reconcile(ctx context.Context, p identity.Principal, id string) error {
	if !s.profile.Authorize(p, "social:reconcile") {
		return socialbridge.ErrBinding
	}
	r, _, e := s.request(ctx, id)
	if e != nil {
		return e
	}
	message := r.Message()
	hash, _ := outbounddelivery.MessageSHA256(message)
	var state string
	e = s.pool.QueryRow(ctx, `select state from communication.outbound_delivery where tenant_id=$1 and channel_code='facebook_pages' and delivery_key=$2`, r.Intent.TenantID, r.Intent.DeliveryKey).Scan(&state)
	if e != nil {
		return e
	}
	if state != "accepted" {
		_, e = s.fence.Claim(ctx, message, hash)
		if !errors.Is(e, outbounddelivery.ErrUnknown) {
			return socialbridge.ErrUnknown
		}
		ref, receipt, e := s.observation(ctx, r)
		if e != nil {
			return e
		}
		if receipt == nil {
			if ref == "" {
				return socialbridge.ErrUnknown
			}
			call := r
			call.ProviderReference = ref
			result, e := s.adapter.Execute(ctx, call, true)
			if e != nil || result.Receipt == nil || result.Receipt.Validate(r.Intent) != nil {
				return socialbridge.ErrUnknown
			}
			receipt = result.Receipt
			if e = s.event(ctx, r, "social.receipt", receipt); e != nil {
				return e
			}
		}
		if e = s.fence.ReconcileAccepted(ctx, message, hash, outbounddelivery.Receipt{ProviderMessageID: receipt.ProviderReference, EvidenceSHA256: receipt.EvidenceSHA256, AcceptedAt: time.Now().UTC()}); e != nil {
			return e
		}
	}
	// Reconcile completion is only allowed without a live worker. The accepted
	// fence and exact original payload remain authority; attempts are unchanged.
	_, e = s.pool.Exec(ctx, `update platform.job j set completed_at=clock_timestamp(),terminal_error_code=null,claimed_by=null,claimed_until=null
 where j.tenant_id=$1 and j.job_id=$2 and j.job_type='social.publish' and j.payload->>'approval_id'=$3 and j.payload->>'profile_sha256'=$4
 and j.completed_at is null and (j.claimed_until is null or j.claimed_until<=clock_timestamp())
 and exists(select 1 from communication.outbound_delivery d where d.tenant_id=j.tenant_id and d.channel_code='facebook_pages' and d.delivery_key=$5 and d.state='accepted' and d.request_sha256_hex=$6)`, r.Intent.TenantID, socialJobID(r.Intent.TenantID, id), id, s.profile.SHA256(), r.Intent.DeliveryKey, hash)
	return e
}
````

### FILE: `internal/platform/postgres/social_publishing_integration_test.go`

```yaml
block_id: "GO-META-PAGE-PUBLISHING-INFRASTRUCTURE:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "05ddfaf3ecada8b00c9d9587d07adc4d0c02d026efe0cd481f95ea7bcefe1c3b"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/socialbridge"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type sdkStep struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Value   any    `json:"value,omitempty"`
	Timeout bool   `json:"timeout,omitempty"`
}

func sdkProcess(t *testing.T, steps []sdkStep) (socialbridge.Process, func() []map[string]any) {
	t.Helper()
	python := os.Getenv("SOCIAL_TEST_PYTHON")
	if python == "" {
		t.Fatal("SOCIAL_TEST_PYTHON exact admitted runtime required")
	}
	root, e := filepath.Abs("../../..")
	if e != nil {
		t.Fatal(e)
	}
	d := t.TempDir()
	state := filepath.Join(d, "steps.json")
	log := filepath.Join(d, "calls.jsonl")
	script := filepath.Join(d, "fixture.py")
	raw, _ := json.Marshal(steps)
	if os.WriteFile(state, raw, 0600) != nil {
		t.Fatal("fixture state")
	}
	quote := func(v string) string { b, _ := json.Marshal(v); return string(b) }
	code := `# Test-only real SDK transport; no production flag enables this transport.
import sys,json
from pathlib import Path
sys.path.insert(0, ` + quote(root) + `)
from meta_page_write.bridge import main
from meta_page_write.adapter import create_api
from requests.adapters import BaseAdapter
from requests import Response, Timeout
from urllib.parse import urlparse,parse_qs
state=Path(` + quote(state) + `)
log=Path(` + quote(log) + `)
class Fixture(BaseAdapter):
 def send(self,request,**kwargs):
  url=urlparse(request.url)
  assert url.scheme=='https' and url.netloc=='graph.facebook.com'
  assert kwargs['timeout']==10
  steps=json.loads(state.read_text()); assert steps
  step=steps.pop(0)
  assert request.method==step['method'] and url.path.rstrip('/')==step['path'],(request.method,url.path,step)
  state.write_text(json.dumps(steps))
  body=request.body.decode() if isinstance(request.body,bytes) else request.body
  with log.open('a') as f:f.write(json.dumps({'method':request.method,'path':url.path,'body':parse_qs(body or '')})+'\n')
  if step.get('timeout'):raise Timeout('fixture')
  response=Response();response.status_code=200;response.url=request.url;response.request=request;response._content=json.dumps(step.get('value',{})).encode();response.headers['Content-Type']='application/json';return response
 def close(self):pass
def api():
 value=create_api('321','fixture-secret','fixture-token');value._session.requests.mount('https://',Fixture());value._session.requests.mount('http://',Fixture());return value
main(api)
`
	if e = os.WriteFile(script, []byte(code), 0600); e != nil {
		t.Fatal(e)
	}
	read := func() []map[string]any {
		raw, e := os.ReadFile(log)
		if os.IsNotExist(e) {
			return nil
		}
		if e != nil {
			t.Fatal(e)
		}
		var values []map[string]any
		for _, line := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
			var v map[string]any
			if json.Unmarshal([]byte(line), &v) != nil {
				t.Fatal("fixture log invalid")
			}
			values = append(values, v)
		}
		return values
	}
	return socialbridge.Process{Python: python, Script: script, ScriptSHA256: socialbridge.Hash([]byte(code))}, read
}
func TestSocialPublishingConnectedSDKApprovalAndRecovery(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("dedicated PostgreSQL required")
	}
	cfg, e := pgxpool.ParseConfig(url)
	if e != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_social_") {
		t.Fatal("isolated local database required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	pool, e := pgxpool.NewWithConfig(ctx, cfg)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant, org := uuid.NewString(), uuid.NewString()
	_, e = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,$2,'Synthetic Social','Fixture')`, tenant, "social-"+uuid.NewString())
	if e != nil {
		t.Fatal(e)
	}
	config := socialbridge.Config{Schema: "elite.meta-page-publishing.v1", TenantID: tenant, OrganizationID: org, PageID: "123", Queue: "social_fixture", Publish: true, Revoke: true, Review: "one_distinct_human", LeaseSeconds: 60, RetrySeconds: 1, PollSeconds: 1, MaxAttempts: 2}
	raw, _ := json.Marshal(config)
	profile, e := socialbridge.Load(raw, socialbridge.Hash(raw))
	if e != nil {
		t.Fatal(e)
	}
	actor := func(subject string) identity.Principal {
		return identity.Principal{TenantID: tenant, Subject: subject, Permissions: map[string]struct{}{"social:request": {}, "social:approve": {}, "social:read": {}, "social:reconcile": {}}, Organizations: map[string]struct{}{org: {}}}
	}
	requester, reviewer := actor("requester"), actor("reviewer")
	makeRequest := func(id string) socialbridge.Request {
		i := socialbridge.Intent{TenantID: tenant, PageID: "123", ProfileSHA256: profile.SHA256(), ApprovalID: id, Operation: "publish", Message: "Exact reviewed fixture bytes"}
		i.ContentSHA256 = socialbridge.Hash([]byte(i.Message))
		i.DeliveryKey = socialbridge.DeliveryKey(i)
		return socialbridge.Request{Intent: i, ScheduledAt: time.Now().UTC().Add(-time.Second), ExpiresAt: time.Now().UTC().Add(time.Hour)}
	}
	observed := func(ref string) any {
		return map[string]any{"id": ref, "from": map[string]string{"id": "123"}, "message": "Exact reviewed fixture bytes", "is_published": true}
	}
	process, calls := sdkProcess(t, []sdkStep{
		{Method: "POST", Path: "/v26.0/123/feed", Value: map[string]string{"id": "123_100"}}, {Method: "GET", Path: "/v26.0/123_100", Value: observed("123_100")},
		{Method: "POST", Path: "/v26.0/123/feed", Value: map[string]string{"id": "123_101"}}, {Method: "GET", Path: "/v26.0/123_101", Timeout: true}, {Method: "GET", Path: "/v26.0/123_101", Value: observed("123_101")},
		{Method: "GET", Path: "/v26.0/123_100", Value: observed("123_100")}, {Method: "DELETE", Path: "/v26.0/123_100", Value: map[string]bool{"success": true}},
		{Method: "POST", Path: "/v26.0/123/feed", Timeout: true},
		{Method: "POST", Path: "/v26.0/123/feed", Value: map[string]string{"id": "123_102"}}, {Method: "GET", Path: "/v26.0/123_102", Timeout: true}, {Method: "GET", Path: "/v26.0/123_102", Timeout: true}, {Method: "GET", Path: "/v26.0/123_102", Value: observed("123_102")},
	})
	service, e := NewSocialPublishing(pool, profile, []byte(strings.Repeat("k", 32)), process)
	if e != nil {
		t.Fatal(e)
	}
	claim := func() Job {
		t.Helper()
		v, e := service.Claim(ctx, "worker")
		if e != nil || len(v) != 1 {
			t.Fatalf("claim %+v %v", v, e)
		}
		return v[0]
	}
	approve := func(r socialbridge.Request) string {
		t.Helper()
		hash, replay, e := service.Submit(ctx, requester, r)
		if e != nil || replay {
			t.Fatalf("submit %v %v", replay, e)
		}
		if _, e = service.Decide(ctx, reviewer, r.Intent.ApprovalID, hash, true, "Reviewed exact bytes"); e != nil {
			t.Fatal(e)
		}
		return hash
	}
	r := makeRequest("published")
	hash, replay, e := service.Submit(ctx, requester, r)
	if e != nil || replay {
		t.Fatal(e)
	}
	if jobs, e := service.Claim(ctx, "worker"); e != nil || len(jobs) != 0 {
		t.Fatal("unapproved request enqueued")
	}
	if _, e = service.Decide(ctx, requester, r.Intent.ApprovalID, hash, true, ""); !errors.Is(e, approval.ErrSeparation) {
		t.Fatalf("self review %v", e)
	}
	if _, e = service.Decide(ctx, reviewer, r.Intent.ApprovalID, strings.Repeat("0", 64), true, ""); e == nil {
		t.Fatal("wrong hash accepted")
	}
	_, replay, e = service.Submit(ctx, requester, r)
	if e != nil || !replay {
		t.Fatal("exact submit replay rejected")
	}
	changed := r
	changed.Intent.Message = "other"
	changed.Intent.ContentSHA256 = socialbridge.Hash([]byte("other"))
	if _, _, e = service.Submit(ctx, requester, changed); e == nil {
		t.Fatal("changed content under same id")
	}
	foreign := requester
	foreign.TenantID = uuid.NewString()
	if _, _, e = service.Submit(ctx, foreign, r); e == nil {
		t.Fatal("foreign tenant")
	}
	badPage := r
	badPage.Intent.PageID = "999"
	badPage.Intent.DeliveryKey = socialbridge.DeliveryKey(badPage.Intent)
	if _, _, e = service.Submit(ctx, requester, badPage); e == nil {
		t.Fatal("foreign page")
	}
	if _, e = service.Decide(ctx, reviewer, r.Intent.ApprovalID, hash, true, "Exact reviewed bytes"); e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `update approval.request set payload='{}' where tenant_id=$1 and request_id=$2`, tenant, r.Intent.ApprovalID); e == nil {
		t.Fatal("approval payload mutable")
	}
	if _, e = pool.Exec(ctx, `update approval.decision set reviewer='intruder' where tenant_id=$1 and request_id=$2`, tenant, r.Intent.ApprovalID); e == nil {
		t.Fatal("decision mutable")
	}
	job := claim()
	if jobs, e := service.Claim(ctx, "other-worker"); e != nil || len(jobs) != 0 {
		t.Fatal("concurrent claim won")
	}
	stale := job
	stale.Attempts++
	if e = service.Process(ctx, stale, "worker"); !errors.Is(e, ErrJobClaimLost) {
		t.Fatalf("wrong generation %v", e)
	}
	if len(calls()) != 0 {
		t.Fatal("stale generation sent")
	}
	if e = service.Process(ctx, job, "worker"); e != nil {
		t.Fatal(e)
	}
	if len(calls()) != 2 {
		t.Fatal("SDK POST/GET not observed")
	}
	_ = service.Process(ctx, job, "worker")
	if len(calls()) != 2 {
		t.Fatal("replay repeated POST")
	}
	status, e := service.Status(ctx, reviewer, r.Intent.ApprovalID)
	if e != nil || status.DeliveryState != "accepted" || len(status.Receipt) == 0 {
		t.Fatalf("status %+v %v", status, e)
	}
	if _, e = service.Status(ctx, foreign, r.Intent.ApprovalID); e == nil {
		t.Fatal("foreign read")
	}
	uncertain := makeRequest("unknown-known-id")
	approve(uncertain)
	first := claim()
	if e = service.Process(ctx, first, "worker"); !errors.Is(e, socialbridge.ErrUnknown) {
		t.Fatalf("expected unknown %v", e)
	}
	if len(calls()) != 4 {
		t.Fatal("unknown POST/GET count")
	}
	// A new instance proves receipt identity and retry behavior are in PostgreSQL.
	restarted, e := NewSocialPublishing(pool, profile, []byte(strings.Repeat("k", 32)), process)
	if e != nil {
		t.Fatal(e)
	}
	service = restarted
	if _, e = pool.Exec(ctx, `update platform.job set available_at=clock_timestamp()-interval '1 second' where tenant_id=$1 and job_id=$2`, tenant, first.JobID); e != nil {
		t.Fatal(e)
	}
	if e = service.Process(ctx, claim(), "worker"); e != nil {
		t.Fatal(e)
	}
	if len(calls()) != 5 || calls()[4]["method"] != "GET" {
		t.Fatal("recovery repeated POST")
	}
	revoke := makeRequest("revoke-reviewed")
	revoke.Intent.Operation = "revoke"
	revoke.Intent.DeliveryKey = socialbridge.DeliveryKey(revoke.Intent)
	revoke.OriginalApprovalID = r.Intent.ApprovalID
	revoke.ProviderReference = "123_999"
	if _, _, e = service.Submit(ctx, requester, revoke); e == nil {
		t.Fatal("arbitrary delete reference")
	}
	revoke.ProviderReference = "123_100"
	approve(revoke)
	if e = service.Process(ctx, claim(), "worker"); e != nil {
		t.Fatal(e)
	}
	if len(calls()) != 7 || calls()[6]["method"] != "DELETE" {
		t.Fatal("governed revoke absent")
	}
	unknown := makeRequest("unknown-no-id")
	approve(unknown)
	lost := claim()
	if e = service.Process(ctx, lost, "worker"); !errors.Is(e, socialbridge.ErrUnknown) {
		t.Fatalf("POST uncertainty %v", e)
	}
	if _, e = pool.Exec(ctx, `update platform.job set available_at=clock_timestamp()-interval '1 second' where tenant_id=$1 and job_id=$2`, tenant, lost.JobID); e != nil {
		t.Fatal(e)
	}
	if e = service.Process(ctx, claim(), "worker"); !errors.Is(e, socialbridge.ErrUnknown) {
		t.Fatal("missing reference guessed")
	}
	if e = service.Reconcile(ctx, reviewer, unknown.Intent.ApprovalID); !errors.Is(e, socialbridge.ErrUnknown) {
		t.Fatal("missing ref reconciled")
	}
	if len(calls()) != 8 {
		t.Fatal("unknown write duplicated")
	}
	terminal := makeRequest("terminal-known-id")
	approve(terminal)
	attempt := claim()
	if e = service.Process(ctx, attempt, "worker"); !errors.Is(e, socialbridge.ErrUnknown) {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `update platform.job set available_at=clock_timestamp()-interval '1 second' where tenant_id=$1 and job_id=$2`, tenant, attempt.JobID); e != nil {
		t.Fatal(e)
	}
	if e = service.Process(ctx, claim(), "worker"); !errors.Is(e, socialbridge.ErrUnknown) {
		t.Fatal(e)
	}
	var attempts int
	var terminalCode string
	if e = pool.QueryRow(ctx, `select attempts,terminal_error_code from platform.job where tenant_id=$1 and job_id=$2`, tenant, attempt.JobID).Scan(&attempts, &terminalCode); e != nil || attempts != 2 || terminalCode == "" {
		t.Fatalf("terminal budget %d %s %v", attempts, terminalCode, e)
	}
	if e = service.Reconcile(ctx, reviewer, terminal.Intent.ApprovalID); e != nil {
		t.Fatal("terminal GET recovery", e)
	}
	var complete bool
	if e = pool.QueryRow(ctx, `select attempts=2 and completed_at is not null and terminal_error_code is null from platform.job where tenant_id=$1 and job_id=$2`, tenant, attempt.JobID).Scan(&complete); e != nil || !complete {
		t.Fatal("reconciliation reset generation or failed ack", e)
	}
	if len(calls()) != 12 || calls()[11]["method"] != "GET" {
		t.Fatal("terminal reconciliation wrote provider")
	}
	crash := makeRequest("crash-after-durable-receipt")
	approve(crash)
	crashedJob := claim()
	m := crash.Message()
	mh, _ := outbounddelivery.MessageSHA256(m)
	if _, e = service.fence.Claim(ctx, m, mh); e != nil {
		t.Fatal(e)
	}
	crashReceipt := socialbridge.Receipt{TenantID: tenant, PageID: "123", ProfileSHA256: profile.SHA256(), ApprovalID: crash.Intent.ApprovalID, DeliveryKey: crash.Intent.DeliveryKey, ProviderReference: "123_103", State: "published", ContentSHA256: crash.Intent.ContentSHA256, EvidenceSHA256: strings.Repeat("a", 64)}
	if e = service.event(ctx, crash, "social.receipt", crashReceipt); e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `update communication.outbound_delivery set locked_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and channel_code='facebook_pages' and delivery_key=$2`, tenant, crash.Intent.DeliveryKey); e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `update platform.job set claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and job_id=$2`, tenant, crashedJob.JobID); e != nil {
		t.Fatal(e)
	}
	if e = service.Reconcile(ctx, reviewer, crash.Intent.ApprovalID); e != nil {
		t.Fatal("durable receipt crash recovery", e)
	}
	if len(calls()) != 12 {
		t.Fatal("durable receipt caused provider call")
	}
	expiring := makeRequest("expiry-while-lock-wait")
	expiring.ExpiresAt = time.Now().UTC().Add(500 * time.Millisecond)
	approve(expiring)
	expiryJob := claim()
	block, e := pool.Begin(ctx)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = block.Exec(ctx, `select tenant_id from platform.tenant where tenant_id=$1 for update`, tenant); e != nil {
		t.Fatal(e)
	}
	done := make(chan error, 1)
	go func() { done <- service.Process(ctx, expiryJob, "worker") }()
	time.Sleep(time.Until(expiring.ExpiresAt) + 100*time.Millisecond)
	if e = block.Commit(ctx); e != nil {
		t.Fatal(e)
	}
	if e = <-done; e == nil {
		t.Fatal("expiry during authority lock wait admitted")
	}
	if len(calls()) != 12 {
		t.Fatal("expired approval sent")
	}
	// This rejected claim is intentionally released without changing its request.
	_, _ = service.jobs.Fail(ctx, tenant, expiryJob.JobID, "worker", "APPROVAL_EXPIRED", expiryJob.Attempts, time.Hour)
	scheduled := makeRequest("scheduled-future")
	scheduled.ScheduledAt = time.Now().UTC().Add(30 * time.Minute)
	approve(scheduled)
	if jobs, e := service.Claim(ctx, "worker"); e != nil || len(jobs) != 0 {
		t.Fatal("scheduled early")
	}
	// Tampering the queue timestamp is still constrained by the approved schedule.
	if _, e = pool.Exec(ctx, `update platform.job set available_at=clock_timestamp()-interval '1 second' where tenant_id=$1 and job_id=$2`, tenant, socialJobID(tenant, scheduled.Intent.ApprovalID)); e != nil {
		t.Fatal(e)
	}
	if e = service.Process(ctx, claim(), "worker"); e == nil {
		t.Fatal("approved not-before bypassed")
	}
	if len(calls()) != 12 {
		t.Fatal("early schedule effect")
	}
	t.Logf("real official SDK: POST=4 GET=7 DELETE=1; known-reference/terminal recovery GET-only; receipt-before-fence crash and expiry-under-lock covered; no network; tenant=%s", tenant)
	_ = fmt.Sprintf("%s", hash)
}
````

### FILE: `internal/platform/httpapi/social_publish.go`

```yaml
block_id: "GO-META-PAGE-PUBLISHING-INFRASTRUCTURE:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7e96189a2eae33e51f84449d422a7ff6b3440366c57c075edd7460f20986b298"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"io"
	"net/http"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/socialbridge"
)

type SocialPublishingService interface {
	Submit(context.Context, identity.Principal, socialbridge.Request) (string, bool, error)
	Decide(context.Context, identity.Principal, string, string, bool, string) (approval.State, error)
	Status(context.Context, identity.Principal, string) (postgres.SocialStatus, error)
	Reconcile(context.Context, identity.Principal, string) error
}

// NewSocialPublishing exposes reviewed requests and read-only reconciliation.
// Principal is supplied only by the existing OIDC verifier, never request JSON.
func NewSocialPublishing(service SocialPublishingService, verifier identity.Verifier) http.Handler {
	mux := http.NewServeMux()
	principal := func(w http.ResponseWriter, r *http.Request) (identity.Principal, bool) {
		p, e := authenticate(r.Context(), r.Header.Get("Authorization"), verifier)
		if e != nil {
			writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
			return p, false
		}
		return p, true
	}
	decode := func(w http.ResponseWriter, r *http.Request, v any) bool {
		if r.Header.Get("Content-Type") != "application/json" {
			writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "application/json required")
			return false
		}
		raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
		if e != nil || socialbridge.Decode(raw, v) != nil {
			writeProblem(w, 400, "INVALID_BODY", "exact contract required")
			return false
		}
		return true
	}
	mux.HandleFunc("POST /v1/social/requests", func(w http.ResponseWriter, r *http.Request) {
		p, ok := principal(w, r)
		if !ok {
			return
		}
		var v socialbridge.Request
		if !decode(w, r, &v) {
			return
		}
		hash, replay, e := service.Submit(r.Context(), p, v)
		if e != nil {
			writeProblem(w, 409, "SOCIAL_BINDING_REJECTED", "request is unauthorized or conflicts with its exact binding")
			return
		}
		code := 201
		if replay {
			code = 200
		}
		writeJSON(w, code, map[string]any{"approval_id": v.Intent.ApprovalID, "request_sha256": hash, "replayed": replay})
	})
	mux.HandleFunc("POST /v1/social/requests/{id}/decision", func(w http.ResponseWriter, r *http.Request) {
		p, ok := principal(w, r)
		if !ok {
			return
		}
		var v struct {
			RequestSHA256 string `json:"request_sha256"`
			Approve       *bool  `json:"approve"`
			Reason        string `json:"reason"`
		}
		if !decode(w, r, &v) {
			return
		}
		if v.Approve == nil {
			writeProblem(w, 400, "INVALID_BODY", "explicit decision required")
			return
		}
		state, e := service.Decide(r.Context(), p, r.PathValue("id"), v.RequestSHA256, *v.Approve, v.Reason)
		if e != nil {
			writeProblem(w, 409, "SOCIAL_DECISION_REJECTED", "decision is unauthorized or no longer applicable")
			return
		}
		writeJSON(w, 200, map[string]any{"state": state})
	})
	mux.HandleFunc("GET /v1/social/requests/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := principal(w, r)
		if !ok {
			return
		}
		v, e := service.Status(r.Context(), p, r.PathValue("id"))
		if e != nil {
			writeProblem(w, 404, "SOCIAL_REQUEST_UNAVAILABLE", "request is unavailable in this scope")
			return
		}
		writeJSON(w, 200, v)
	})
	mux.HandleFunc("POST /v1/social/requests/{id}/reconcile", func(w http.ResponseWriter, r *http.Request) {
		p, ok := principal(w, r)
		if !ok {
			return
		}
		if e := service.Reconcile(r.Context(), p, r.PathValue("id")); e != nil {
			writeProblem(w, 409, "SOCIAL_RECONCILIATION_UNCONFIRMED", "no further write permitted; provider observation or an active lease remains unresolved")
			return
		}
		writeJSON(w, 200, map[string]string{"state": "reconciled"})
	})
	return mux
}
````

### FILE: `internal/platform/httpapi/social_publish_test.go`

```yaml
block_id: "GO-META-PAGE-PUBLISHING-INFRASTRUCTURE:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "4f0a3ba7a39952c0530104ec2163f6d3352bfa270eabc648d2b07927764022b9"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/socialbridge"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type socialTestVerifier struct{}

func (socialTestVerifier) Verify(_ context.Context, token string) (identity.Principal, error) {
	if token != "verified" {
		return identity.Principal{}, identity.ErrUnauthenticated
	}
	return identity.Principal{Subject: "verified-operator", TenantID: "tenant"}, nil
}

type socialTestService struct {
	decisions int
	principal string
	approved  bool
}

func (s *socialTestService) Submit(context.Context, identity.Principal, socialbridge.Request) (string, bool, error) {
	return "hash", false, nil
}
func (s *socialTestService) Decide(_ context.Context, p identity.Principal, _, _ string, a bool, _ string) (approval.State, error) {
	s.decisions++
	s.principal = p.Subject
	s.approved = a
	return approval.StateApproved, nil
}
func (s *socialTestService) Status(context.Context, identity.Principal, string) (postgres.SocialStatus, error) {
	return postgres.SocialStatus{}, nil
}
func (s *socialTestService) Reconcile(context.Context, identity.Principal, string) error { return nil }
func TestSocialHTTPExplicitDecisionAndVerifiedIdentity(t *testing.T) {
	s := &socialTestService{}
	h := NewSocialPublishing(s, socialTestVerifier{})
	call := func(token, body string) int {
		r := httptest.NewRequest(http.MethodPost, "/v1/social/requests/request/decision", strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w.Code
	}
	for _, body := range []string{`{"request_sha256":"hash","approve":null,"reason":"x"}`, `{"request_sha256":"hash","reason":"x"}`, `{"request_sha256":"hash","approve":true,"Approve":false}`, `{"request_sha256":"hash","approve":true,"requester":"attacker"}`} {
		if code := call("verified", body); code != 400 {
			t.Fatalf("body %s code %d", body, code)
		}
	}
	if s.decisions != 0 {
		t.Fatal("invalid decision reached owner")
	}
	valid := `{"request_sha256":"hash","approve":true,"reason":"reviewed"}`
	if code := call("unverified", valid); code != 401 {
		t.Fatal(code)
	}
	if code := call("verified", valid); code != 200 || s.decisions != 1 || s.principal != "verified-operator" || !s.approved {
		t.Fatal(code, s)
	}
}
````

### FILE: `cmd/social-publishing/main.go`

```yaml
block_id: "GO-META-PAGE-PUBLISHING-INFRASTRUCTURE:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8fa415290944b8326dc9b26aa68a40800c1d0e8df2e36117196ff765b1499b81"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/socialbridge"
	"github.com/jackc/pgx/v5/pgxpool"
)

func run(ctx context.Context) error {
	p, e := socialbridge.LoadFile(os.Getenv("SOCIAL_PROFILE_FILE"), os.Getenv("SOCIAL_PROFILE_SHA256"))
	if e != nil {
		return socialbridge.ErrBinding
	}
	key, e := hex.DecodeString(os.Getenv("SOCIAL_FENCE_HMAC_KEY_HEX"))
	if e != nil || len(key) < 32 {
		return socialbridge.ErrBinding
	}
	adapter := socialbridge.Process{Python: os.Getenv("SOCIAL_PYTHON"), Script: os.Getenv("SOCIAL_BRIDGE_FILE"), ScriptSHA256: os.Getenv("SOCIAL_BRIDGE_SHA256")}
	for _, k := range []string{"META_APP_ID", "META_APP_SECRET", "META_PAGE_ACCESS_TOKEN"} {
		v := os.Getenv(k)
		if v == "" {
			return socialbridge.ErrBinding
		}
		adapter.Environment = append(adapter.Environment, k+"="+v)
	}
	if adapter.Validate() != nil || os.Getenv("DATABASE_URL") == "" || os.Getenv("OIDC_ISSUER") == "" || os.Getenv("OIDC_AUDIENCE") == "" {
		return socialbridge.ErrBinding
	}
	if adapter.Preflight(ctx) != nil {
		return socialbridge.ErrBinding
	}
	verifier, e := identity.NewOIDCVerifier(ctx, os.Getenv("OIDC_ISSUER"), os.Getenv("OIDC_AUDIENCE"))
	if e != nil {
		return e
	}
	pool, e := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if e != nil {
		return e
	}
	defer pool.Close()
	if e = pool.Ping(ctx); e != nil {
		return e
	}
	service, e := postgres.NewSocialPublishing(pool, p, key, adapter)
	if e != nil {
		return e
	}
	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	workerDone := make(chan struct{})
	go func() {
		defer close(workerDone)
		lastState := ""
		report := func(state string) {
			if state != lastState {
				slog.Info("social worker state", "state", state)
				lastState = state
			}
		}
		ticker := time.NewTicker(time.Duration(p.Config().PollSeconds) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-workerCtx.Done():
				return
			case <-ticker.C:
				jobs, e := service.Claim(workerCtx, "social-publishing")
				if e != nil {
					report("QUEUE_UNAVAILABLE")
					continue
				}
				if len(jobs) == 0 {
					report("IDLE")
				}
				for _, job := range jobs {
					e := service.Process(workerCtx, job, "social-publishing")
					if e == nil {
						report("ACKNOWLEDGED")
					} else if errors.Is(e, socialbridge.ErrUnknown) {
						report("PROVIDER_UNRESOLVED")
					} else {
						report("DISPATCH_NOT_ADMITTED")
					}
				}
			}
		}
	}()
	defer func() { cancel(); <-workerDone }()
	addr := os.Getenv("SOCIAL_HTTP_ADDRESS")
	if addr == "" {
		addr = "127.0.0.1:8097"
	}
	mux := http.NewServeMux()
	mux.Handle("/", httpapi.NewSocialPublishing(service, verifier))
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if pool.Ping(ctx) != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	server := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 45 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 16384}
	return httpapi.ServeUntilShutdown(ctx, server, 40*time.Second)
}
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if e := run(ctx); e != nil && !errors.Is(e, context.Canceled) {
		slog.Error("social publishing configuration or runtime unavailable")
		os.Exit(2)
	}
}
````

### FILE: `meta_page_write/bridge.py`

```yaml
block_id: "GO-META-PAGE-PUBLISHING-INFRASTRUCTURE:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f8dcb6a2c66ed70cef4fbdd7e2ea542e95ebd063a1602eabaab0c96c3c5d0b0b"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED bounded process bridge; the seven-file SDK adapter is unchanged."""
from dataclasses import asdict
import json
import hashlib
import importlib.metadata
import os
from pathlib import Path
import sys

# -I excludes cwd/user packages. Import only this materialized package root.
sys.path.insert(0, str(Path(__file__).resolve().parent.parent))
from meta_page_write.adapter import (
    ApprovedIntent, FacebookPageWriteAdapter, Scope, UnknownDelivery, create_api,
)


def preflight():
    # Source bytes already compared to the fixed official wheel by G0-G8.
    locked = {
        "facebook_business/adobjects/page.py": "0d08e317fe078de5ad2805ab5501423fe51fea39ef90ae4a0160122a5914ee51",
        "facebook_business/adobjects/pagepost.py": "b0b2009a0667ad57914423284877b104068dca1af4bdfb7f6543b0bec41f0b16",
        "facebook_business/api.py": "2f2e537354859758faba0d89a910b2212b478e0b2aaee771bb1c239abf4c2c72",
        "facebook_business/session.py": "72dd2e43efacc27728222d6c176202f6ff8919a85ba19811ac71fa357bc01587",
        "facebook_business/adobjects/objectparser.py": "9634dee1363b3930469b74942a5b5f58a53bfef5f21bd91519a2e9a25f1de9be",
    }
    distribution = importlib.metadata.distribution("facebook-business")
    for path, expected in locked.items():
        if hashlib.sha256(distribution.locate_file(path).read_bytes()).hexdigest() != expected:
            raise ValueError("installed SDK source differs")
    root = Path(__file__).resolve().parent
    if hashlib.sha256((root / "adapter.py").read_bytes()).hexdigest() != "a397d84a5f92fc8d9d3a44850ad568875da837e59800aeab324fcb5adc6a9e1e":
        raise ValueError("materialized adapter differs")
    requirements = (root / "requirements-windows-py314.lock").read_bytes()
    if hashlib.sha256(requirements).hexdigest() != "a4b4b65d2a19be4bb51892bcb25a4e8127f7b2c16b8f49a2f40b319797314e19":
        raise ValueError("dependency lock differs")
    for line in requirements.decode("utf-8").splitlines():
        if not line or line.startswith("#"):
            continue
        name, url = line.split(" @ ", 1)
        version = url.rsplit("/", 1)[1].split("-")[1]
        if importlib.metadata.version(name) != version:
            raise ValueError("dependency version differs")
    create_api(os.environ.get("META_APP_ID"), os.environ.get("META_APP_SECRET"), os.environ.get("META_PAGE_ACCESS_TOKEN"))


def main(api_factory=None):
    try:
        raw = sys.stdin.buffer.read(32769)
        if not raw or len(raw) > 32768:
            raise ValueError("protocol bounds")
        value = json.loads(raw.decode("utf-8"))
        if set(value) != {"request", "reconcile"} or type(value["reconcile"]) is not bool:
            raise ValueError("protocol fields")
        request = value["request"]
        if set(request) != {"intent", "scheduled_at", "expires_at", "original_approval_id", "provider_reference"}:
            raise ValueError("request fields")
        intent = ApprovedIntent(**request["intent"])
        factory = api_factory or (lambda: create_api(os.environ.get("META_APP_ID"), os.environ.get("META_APP_SECRET"), os.environ.get("META_PAGE_ACCESS_TOKEN")))
        adapter = FacebookPageWriteAdapter(Scope(intent.tenant_id, intent.page_id, intent.profile_sha256), factory())
        if value["reconcile"]:
            if intent.operation != "publish":
                # A missing/permission-denied post cannot prove a DELETE. Never
                # turn absence into success or issue another DELETE implicitly.
                raise UnknownDelivery("REVOKE_REQUIRES_PROVIDER_PROOF", request["provider_reference"])
            receipt = adapter.reconcile(intent, request["provider_reference"])
        elif intent.operation == "publish":
            receipt = adapter.publish_once(intent)
        else:
            receipt = adapter.revoke_once(intent, request["provider_reference"])
        result = {"receipt": asdict(receipt)}
    except UnknownDelivery as error:
        result = {"unknown": True, "code": error.code}
        if error.provider_reference:
            result["provider_reference"] = error.provider_reference
    except Exception:
        # The caller conservatively fences any process/protocol exception.
        # SDK exceptions and secrets are never printed.
        result = {"unknown": True, "code": "PROCESS_BOUNDARY_UNCONFIRMED"}
    sys.stdout.write(json.dumps(result, separators=(",", ":"), ensure_ascii=True))


if __name__ == "__main__":
    if sys.argv[1:] == ["--preflight"]:
        try:
            preflight()
            sys.stdout.write('{"state":"PASS"}')
        except Exception:
            sys.stdout.write('{"state":"FAIL"}')
            sys.exit(2)
    elif len(sys.argv) == 1:
        main()
    else:
        sys.exit(2)
````

### FILE: `db/migrations/0060_social_publish.up.sql`

```yaml
block_id: "GO-META-PAGE-PUBLISHING-INFRASTRUCTURE:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f2a82e5051966fff87f450c76d0c946e9f88e5a9d450c115b9e141ceae0e6a96"
variables: []
secrets_allowed: false
```

````sql
begin;
-- Provider references and receipts remain public outbox evidence. Protect the
-- fields which authorize GET-only recovery while allowing normal outbox leases.
create function communication.guard_social_observation() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.aggregate_type='social_approval' then raise exception 'social evidence is immutable';end if;return old;
 end if;
 if (old.aggregate_type='social_approval' or new.aggregate_type='social_approval') and
 row(new.tenant_id,new.event_id,new.aggregate_type,new.aggregate_id,new.aggregate_version,new.event_type,new.schema_version,new.occurred_at,new.payload)
 is distinct from row(old.tenant_id,old.event_id,old.aggregate_type,old.aggregate_id,old.aggregate_version,old.event_type,old.schema_version,old.occurred_at,old.payload)
 then raise exception 'social evidence is immutable';end if;return new;
end $$;
create trigger social_observation_immutable before update or delete on platform.outbox_event
 for each row execute function communication.guard_social_observation();
commit;
````

### FILE: `db/migrations/0060_social_publish.down.sql`

```yaml
block_id: "GO-META-PAGE-PUBLISHING-INFRASTRUCTURE:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "925923406cca0152d3341395fb203b216eb83a5fa4ac3fa34b54d26b079c89be"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$begin
 if exists(select 1 from platform.outbox_event where aggregate_type='social_approval') then
  raise exception 'social provider evidence exists; retain schema or use reviewed forward migration';
 end if;
end $$;
drop trigger social_observation_immutable on platform.outbox_event;
drop function communication.guard_social_observation();
commit;
````

### FILE: `deploy/social/profile.reference.json`

```yaml
block_id: "GO-META-PAGE-PUBLISHING-INFRASTRUCTURE:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "dd5e3c0a487ef6c95d43ac57bd954c98718a3d893ae6b50e7ee0df513bddf2dd"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite.meta-page-publishing.v1",
  "tenant_id": "11111111-1111-4111-8111-111111111111",
  "organization_id": "22222222-2222-4222-8222-222222222222",
  "page_id": "123",
  "queue": "social_publishing",
  "publish": true,
  "revoke": true,
  "review": "one_distinct_human",
  "lease_seconds": 60,
  "retry_seconds": 60,
  "poll_seconds": 5,
  "max_attempts": 10
}
````

### FILE: `deploy/social/.env.example`

```yaml
block_id: "GO-META-PAGE-PUBLISHING-INFRASTRUCTURE:file13:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "01e333bbd69bbaa6ec84c9446f6fc8e50a7e5c188ecbcce0d6f1cf6aba3c90df"
variables: []
secrets_allowed: false
```

````text
# Empty deployment inputs. Never put credentials in the materialized profile.
DATABASE_URL=
OIDC_ISSUER=
OIDC_AUDIENCE=
SOCIAL_PROFILE_FILE=
SOCIAL_PROFILE_SHA256=
SOCIAL_FENCE_HMAC_KEY_HEX=
SOCIAL_PYTHON=
SOCIAL_BRIDGE_FILE=
SOCIAL_BRIDGE_SHA256=
SOCIAL_HTTP_ADDRESS=
META_APP_ID=
META_APP_SECRET=
META_PAGE_ACCESS_TOKEN=
````

### FILE: `docs/SOCIAL_PUBLISHING.md`

```yaml
block_id: "GO-META-PAGE-PUBLISHING-INFRASTRUCTURE:file14:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b0a1ae6f81abcd303e3d41f5553150ba32df1cf8eced4d926920e4e591c849e5"
variables: []
secrets_allowed: false
```

````markdown
# Facebook Page publication infrastructure

This materializable scope publishes exact reviewed text to one configured
Facebook Page. It includes a runnable Go API/worker, immutable human approvals,
durable scheduled jobs, a single-write outbound fence, the official Meta Python
SDK, exact GET observation, public receipts and separately approved revocation.
No marketing policy, moderation, audience inference or eligibility is generated.

## Owners and provenance

The existing approval.Registry owns human separation; approval.request/decision
remain the only approval ledger. The new typed binding supplies zero automatic
approval policy. Its three manual kinds are social_publish, social_revoke and
whatsapp_reply. The generic PostgreSQL owner is shared with WhatsApp.
platform.job owns scheduling, bounded retry budgets and claim generations.
communication.outbound_delivery owns the at-most-one write attempt. Public
provider references/receipts live in platform.outbox_event with immutable social
evidence fields. No social ledger or replacement financial logic is introduced.

All new files are AUTHORED configuration/serialization/authorization/persistence
and process glue. The seven-file Meta SDK adapter is an unchanged dependency.
Provider POST/GET/DELETE behavior is DEPENDENCY_PIN facebook-business26.0.1,
fixed commit788f363d15b1269ab5efb7cd00fb5e3b133cd99b / Graphv26.0. The exact
LicenseRef-Meta-Platform bytes remain owned by that adapter. Nothing here is
attributed to a vendor as a locally adapted marketing algorithm.

## Start

1. Compose the admitted backend, jobs, outbound fence, human approval, seven-file
   Meta Page adapter and this publishing integration. Apply migrations through
   0058 plus0060 (0059 belongs to the independently selected WhatsApp delta).
2. Install the exact18-wheel requirements-windows-py314.lock into an isolated
   CPython3.14 Windows environment using --require-hashes --no-deps. Keep all
   source and license notices supplied by its dependency pack.
3. Review a profile derived from deploy/social/profile.reference.json, replacing
   the clearly synthetic tenant/organization/Page identifiers during deployment.
   Select its SHA256 independently. All worker replicas use the same bytes/hash.
   Queue/timing/attempt bounds are technical configuration, not marketing policy.
4. Supply the empty environment inputs from deploy/social/.env.example when
   activating. No secrets are requested or embedded now. Pin the materialized
   meta_page_write/bridge.py SHA256 and isolated Python executable path.
5. Build `go build -mod=readonly ./cmd/social-publishing` and run the resulting
   executable. It serves on loopback127.0.0.1:8097 by default and runs the worker.
   OIDC verifier owns identity; the caller needs the configured organization and
   social:request/social:approve/social:read/social:reconcile permissions as used.
   Front a non-loopback deployment with the target's admitted TLS/auth gateway.

Before serving, startup rejects missing profile/hash, credentials, runtime,
provider SDK source drift, dependency version drift and a changed Python bridge.
The SDK preflight is offline and verifies fixed source bytes; it does not prove
Page ownership, Meta app access approval or live token permissions.

## Operator contract

POST /v1/social/requests submits the Request JSON documented by
internal/socialbridge/contract.go. Server-owned tenant/Page/profile checks reject
cross-scope values. ApprovalID identifies the immutable request; DeliveryKey is
social_ plus SHA256(tenant NUL page NUL approvalID NUL operation). ContentSHA256
binds the exact UTF-8 message. ScheduledAt and ExpiresAt are explicit UTC instants.
The response provides the canonical RequestSHA256 for exact review.

GET /v1/social/requests/{id} exposes the exact pending payload/hash and current
approval, delivery and receipt in the authorized tenant/organization.
POST /v1/social/requests/{id}/decision requires request_sha256, explicit boolean
approve and reason. Its verified human must differ from the requester. Only the
approved transaction creates the scheduled job. Rejection creates no job.

The worker revalidates durable approval, tenant, profile, page, content,
not-before/expiry and current claim generation in the fence transaction. A replay
cannot obtain another write attempt. Its Python child calls the real SDK POST
once, then GET verifies exact id/from.id/message/is_published before success.

After an uncertain response, known post IDs are persisted and retries perform
GET only. POST /v1/social/requests/{id}/reconcile provides the same GET-only
recovery after the queue budget is exhausted. It preserves the attempt number.
A saved receipt recovers a crash before fence/queue acknowledgement without
another SDK call. Restarting the host preserves approvals, jobs and evidence.

Revocation is another request, operation=revoke, separately reviewed. It names
the original approved publication and its exact public receipt/postID/content.
The worker verifies that binding, then the SDK performs GET and DELETE once.
A different Page, post reference, message or original approval cannot authorize
deletion. A rejected pending request never dispatches. Disabling an action in a
new reviewed profile prevents new admission under that profile; do not assume
changing a file can recall a provider operation already in flight.

## Ambiguous external outcomes

A lost POST response with no provider ID remains UNKNOWN; neither the operator
endpoint nor a restarted worker invents a post ID or resends it. A lost DELETE
acknowledgement also remains UNKNOWN: absence or permission denial cannot prove
deletion. These are explicit preserved external outcomes, not successful sends
or missing implementation. Retain the request and seek provider evidence through
the operator process. No credential alone is claimed to resolve such a case.

## Retention, rollback and boundary

Planned publication text is stored in the access-controlled approval payload;
tokens are never stored there. Provider receipts contain public IDs, hashes and
state. Do not send private customer records as publication text. Target retention
configuration must preserve approval/receipt evidence needed for reconciliation.
Stop admission and drain workers before changing executable/profile versions.
The down migrations refuse removal once their governed evidence exists; use a
reviewed forward change rather than deleting that evidence.

Focused PostgreSQL/real-SDK-fixture tests cover exact approval, replay, scheduled
admission, tenant/Page/content isolation, stale claims, one POST, GET-only
recovery, terminal recovery, receipt crash, expiry under a lock and governed
DELETE. These library results do not certify live app/Page permission, target
load/availability, broader social networks or a fresh composition-wide security
review. The parent composition owns its T2803 SCA/SAST and final admission.
````

## 6. Configuration surface

deploy/social/.env.example contains only empty values; profile.reference.json identifiers are synthetic fixtures. Independent profile/bridge SHA required. No secrets are in packs.

## 7. Dependency bill

DEPENDENCY_PIN Meta SDK26.0.1 and unchanged18-wheel lock from PYTHON-OFFICIAL-META-PAGE-WRITE-ADAPTER. Shared Go1.26.8/pgx5.10.0/PG18.6 and google/uuid1.6.0 pins unchanged. All new files AUTHORED inevitable contract/transport/persistence glue; no algorithm vendor attribution.

## 8. Apply order

Apply shared0058, independently selected WhatsApp0059 if used, then0060. Build cmd/social-publishing. Separate runnable host composes the typed owners without replacing the franchise main. Stop/drain before profile/executable changes; preserve unknown/receipt evidence. Down refuses evidence loss.

## 9. Verification

social-connected-pg-final/result.json: real PG59migrations, real official SDK fixture POST4/GET7/DELETE1, new-instance recovery, terminal budget recovery, crash after receipt, locked-expiry rejection. HTTP/profile/runtime preflight tests and finite3s fuzz57,718exec PASS. Parent owns remaining composition-wide T2803 SCA/SAST/threat model and profile promotion.

## 10. Reconstruction evidence

21new/modified source files across this14-file pack and shared11-file human owner plus unchanged7-file leaf; exact materialization roundtrip. No full-core or live production certification; SOURCE AUTHORED glue versus DEPENDENCY_PIN SDK remains explicit.

