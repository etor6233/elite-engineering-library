# Go Business Policy Profile

## 1. Metadata

```yaml
pack_id: "GO-BUSINESS-POLICY-PROFILE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Explicit hash-locked reference scheduling/pricebook contract, bounded immutable configuration loader, shared host owner binding and durable audit/replay enforcement of compatible values."
stacks: ["Go 1.26.8", "PostgreSQL 18.6"]
compatible_with: ["GO-COMMERCE-PRICING-PAYMENT-API 0.6.4", "GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.10.21", "GO-ELECTROMOBILITY-APPLICATION with supplied policy host patch"]
incompatible_with: ["unlocked or incomplete custom profile", "mixed profile hashes across host owners", "changed modes incompatible with migrations0007/0009/0015", "claims of complete scheduling or BC policy equivalence"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-11"
```

V402 library composition component. All ten files are AUTHORED configuration,
contract enforcement, composition glue or tests. No vendor source or architecture
is attributed. No new third-party license/notice owner is introduced. Internal
library use is the scope; LicenseRef-Workspace-Owner is not an OSI license or an
automatic permission to relicense/distribute the owner's code publicly.

Conditions: compose both exact updated owners and integrate the supplied main
and environment deltas into their existing owners. D has this wiring and passed
the host checks; publication/profile promotion remains a separate root operation.
The PBC-CORE relationship is a documented sidecar extension, not a duplicated
portablebusinessconfig runtime. Its base schema is not changed or validated here.

## 2. Applicability

Existing reference semantics are explicit in JSON: 30-minute lead, 8-hour maximum,
capacity1..100, requested/confirmed occupancy, working containment, unavailable
overlap rejection, single active tenant/market/currency price source and half-open
validity. Compatible numeric changes are selected by exact file hash without
reprogramming Go. Unsupported modes are rejected before startup; they need a
separate admitted schema/owner change, not credentials.

## 3. Architecture contract

Trusted deployment configuration owns an immutable Profile. The host selects it
before OIDC/database work, and shares the same pointer with Commerce and Journey.
Existing PostgreSQL transactions/locks own quotas and single-source enforcement.
The policy does not compute staffing, price priority, compliance or overbooking.
New effects carry the profile hash. Reference request hashes preserve historical
identity; custom-profile hashes are domain-separated without changing keys.

## 4. Exact file manifest

```text
CREATE contracts/business-policy.schema.json
CREATE docs/provenance/BUSINESS_POLICY_PROFILE.md
CREATE internal/businesspolicy/profile.go
CREATE internal/businesspolicy/profile_test.go
CREATE internal/businesspolicy/reference-profile.json
CREATE internal/businesspolicy/reference.go
CREATE internal/franchisejourney/policy_profile_test.go
CREATE internal/platform/postgres/policy_profile_integration_test.go
CREATE cmd/electromobility-api/business_policy.go
CREATE cmd/electromobility-api/business_policy_test.go
```

## 5. Materialization blocks

### FILE: `contracts/business-policy.schema.json`
```yaml
block_id: "GO-BUSINESS-POLICY-PROFILE:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "Local explicit reference contract, configuration validation, composition/binding glue or tests; no external algorithm attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "b9dfd47b5d97601df53f030d18a821a7b0485d196a8eed0702a657d7ee546b09"
variables: []
secrets_allowed: false
```
````json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "urn:elite:business-policy:v1",
  "$comment": "AUTHORED compatibility sidecar to PBC-CORE/1.0.0. The authoritative Go loader additionally rejects duplicate/case-aliased keys, invalid UTF-8, more than16384 bytes, minimum capacity greater than maximum, and bytes not matching the independently selected lowercase SHA-256.",
  "type": "object",
  "additionalProperties": false,
  "required": [
    "schema",
    "extends",
    "profile_id",
    "revision",
    "appointments",
    "pricebooks"
  ],
  "properties": {
    "schema": {
      "const": "elite-business-policy/v1"
    },
    "extends": {
      "const": "PBC-CORE/1.0.0"
    },
    "profile_id": {
      "type": "string",
      "pattern": "^[a-z][a-z0-9_-]{1,63}$"
    },
    "revision": {
      "type": "integer",
      "minimum": 1,
      "maximum": 9223372036854775807
    },
    "appointments": {
      "type": "object",
      "additionalProperties": false,
      "required": [
        "lead_time_seconds",
        "maximum_slot_seconds",
        "minimum_slot_capacity",
        "maximum_slot_capacity",
        "occupying_states",
        "working_window",
        "unavailable_window"
      ],
      "properties": {
        "lead_time_seconds": {
          "type": "integer",
          "minimum": 0,
          "maximum": 9223372036
        },
        "maximum_slot_seconds": {
          "type": "integer",
          "minimum": 1,
          "maximum": 28800
        },
        "minimum_slot_capacity": {
          "type": "integer",
          "minimum": 1,
          "maximum": 100
        },
        "maximum_slot_capacity": {
          "type": "integer",
          "minimum": 1,
          "maximum": 100
        },
        "occupying_states": {
          "const": [
            "requested",
            "confirmed"
          ]
        },
        "working_window": {
          "const": "contains_slot"
        },
        "unavailable_window": {
          "const": "reject_overlap"
        }
      }
    },
    "pricebooks": {
      "type": "object",
      "additionalProperties": false,
      "required": [
        "selection",
        "validity"
      ],
      "properties": {
        "selection": {
          "const": "single_active_per_market_currency"
        },
        "validity": {
          "const": "half_open"
        }
      }
    }
  }
}
````

### FILE: `docs/provenance/BUSINESS_POLICY_PROFILE.md`
```yaml
block_id: "GO-BUSINESS-POLICY-PROFILE:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "Local explicit reference contract, configuration validation, composition/binding glue or tests; no external algorithm attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "5bbb7d830e6b4ebf261d35efbe1f10f2eae43a59593a5bf724d79b0c53926da5"
variables: []
secrets_allowed: false
```
````markdown
# Business policy profile: explicit reference contract

Scope: local library configuration infrastructure. All files added by this delta
are AUTHORED contract data, validation, audit/binding glue, or tests. No BC,
Microsoft or other vendor authorship is claimed. Existing owner files retain
AUTHORED provenance; this does not reclassify their remaining business logic.

## Contract and exact defaults

The materialized `internal/businesspolicy/reference-profile.json` preserves the
previous service values: lead 1800 seconds, maximum slot 28800 seconds, capacity 1..100.
It explicitly selects occupying states requested+confirmed, a working window
containing the entire slot, rejection of any unavailable overlap, one active
pricebook per tenant/market/currency at each instant, and half-open validity.

Reference SHA256: `da3feda4d66a19806f20234773720c2f98f6c7d4251640d6fb12cb58b24be5c5`.
JSON Schema SHA256: `b9dfd47b5d97601df53f030d18a821a7b0485d196a8eed0702a657d7ee546b09`.

This is a sidecar extension of PBC-CORE schema 1.0.0. The existing PBC-CORE 0.1.0
pack explicitly leaves workflow/domain invariants to tested extensions. Its base
schema has no policy-extension field; it rejects unknown fields and `customFields`
are UI data declarations. We therefore neither modify the base schema nor put
executable policy in customFields. There is no `portablebusinessconfig` Go package
in the selected composition to import. `extends` identifies this compatibility
relationship; it does not claim that the base PBC profile has been loaded or
validated by this loader. Base-profile validation remains its original owner.
PBC pack SHA256 at inspection: `e5cf5b37c9dbcf709c87076c085a8909dfd5fffbc4055ca714922347e36ff2ce`.

## Compatible changes without changing Go

Copy the reference JSON to a deployment configuration file. Change the values,
profile_id and revision in a reviewable configuration revision, record its exact
SHA256 in the deployment's independent artifact lock, and use:

```go
profile, err := businesspolicy.LoadFile(configPath, lockedSHA256)
// Reject err before serving requests. Do not derive the expected hash from the
// same untrusted file inside this loading operation.
journeyRepo, err := postgres.NewFranchiseJourneyWithProfile(pool, profile)
journeyService, err := franchisejourney.NewServiceWithProfile(journeyRepo, ids, clock, profile)
commerceRepo, err := postgres.NewCommerceWithProfile(pool, profile)
```

All errors must be handled; these statements illustrate the underlying API.
The connected host now implements selection in `cmd/electromobility-api/business_policy.go`.
Both BUSINESS_POLICY_PROFILE_FILE and BUSINESS_POLICY_PROFILE_SHA256 absent/empty
select the embedded reference. A custom file requires both fields; missing,
unreadable, incompatible or hash-mismatched configuration rejects startup. The
host validates the file before OIDC discovery, database access or serving HTTP.
It constructs Commerce and Journey from the same immutable *Profile. These fields
are non-secret and independent of the web presentation BUSINESS_CONFIG_FILE.
The host's existing payment, fiscal and handover activation fields are preserved.

Lead time may be any nonnegative number of whole seconds that fits Go duration.
Maximum slot length may be 1..28800 seconds; capacity minimum/maximum must satisfy
1<=minimum<=maximum<=100. These outer duration/capacity ceilings are compatibility
limits of migration 0007, not secretly retained defaults. The declared modes and
occupying-state list have exactly one supported value in v1: other values are
rejected before construction because migrations 0007/0009/0015 and existing owners
would otherwise contradict the profile. Changing them requires a separately
admitted schema/owner delta; changing credentials cannot resolve that constraint.

Loading rejects unknown/duplicate/case-aliased keys, invalid UTF-8, trailing JSON,
missing required numeric fields, excessive size/depth, incompatible values and
hash mismatch. Profile internals are immutable. Explicit service construction
requires the repository's exact hash; the legacy reference constructor rejects a
custom bound repository at startup, rather than mixing policies. Legacy unbound
fakes remain supported by that legacy constructor for existing tests.

## Persistence and replay

Slot creation enforces configured bounds using the database clock, including for
direct repository calls. Previously direct callers could bypass the service lead
time; this delta closes that inconsistency. Public availability and booking now
use the same minimum start threshold. A listing can become stale before booking;
booking always rechecks under the existing transaction/locks. Requested/confirmed
counts, interval operators and atomic exclusivity remain existing SQL enforcement
of the explicitly selected compatible contract. No pricing priority, staffing,
overtime, overbooking, fiscal or compliance rule is introduced.

New slot/appointment events and pricebook-activation events include the profile
SHA256. Appointment and slot-creation idempotency keys are unchanged. Exact
reference-profile request hashes retain historical byte identity; other profiles
bind request hashes with SHA256(domain-separator,profile-hash,request-hash). Same
profile replays return the existing result; a changed profile conflicts rather
than pretending the historical decision used the new policy. Existing rows and
receipts are never rewritten. Lowering maximum capacity applies to newly created
slots and does not retroactively shrink existing slot quotas.

Profiles are selected per composed instance, not per public request. There is no
new central policy registry, live reload, coordinated multi-instance rollout or
automatic rollback. The deployment must select one locked revision consistently
and restart its owners; this scope proves constructors and persistent effect
binding, not a control-plane rollout. Existing service time validation still runs
before repository replay; this delta proves durable repository replay and does
not expand expired-request replay semantics in the transport/service layer.

## Remaining source/policy gaps are not reclassified

The following preexisting business choices remain outside this narrow contract:
appointment/resource kind enumerations; valid lifecycle transitions; resource
skills/count/identity coupling; slot open/closed lifecycle and non-overlap grouping;
pricebook draft/active/archive lifecycle; consent, quote retention and handover
acceptance. See `internal/franchisejourney/service.go` (appointmentTransitions,
appointmentKinds, resourceKinds, CreateServiceResource) and migrations 0007/0009/
0015. They are not BC adaptations, credential blockers, or closed by placing the
values handled here in a JSON file. The next change must either bind each required
policy explicitly or derive a real algorithm from an admitted exact source.

## Verification scope

Targeted tests cover changed compatible configuration, invalid/hash-unbound input,
service/repository mismatch, prospective bounds, reference/custom public listing,
new request/replay/configuration conflict and persistent profile event hashes.
PostgreSQL 18.6 applied 54 unchanged admitted migrations; the new focused integration
test passed without skips. The final PG receipt is policy-config-pg-2/result.json.
Policy parser fuzz: Go 1.26.8, 3 seconds, 2 workers, 53626 executions, 3 seeds, PASS. Fuzz ran
an exact-copy minimal module to avoid rerunning unrelated domain suites. The
Commerce/Accounting end-to-end suites were not rerun. Final service binding test
and full Go build passed after the legacy-constructor guard was added. This is
local/synthetic evidence, not production readiness or a source audit of all owners.

Host integration validation: two focused Go tests and six executable-startup
cases passed. The reference and a valid file proceed to the existing required
connection-configuration check; partial fields, wrong hash and an incompatible
hash-matching file exit at policy validation first. No listener, live account,
OIDC discovery or database connection is needed by these cases. Host vet/build
passed. Receipts: policy-host-tests/result.json. Prior PostgreSQL and fuzz evidence
is reused by exact unchanged implementation hashes; PostgreSQL was not rerun for
host wiring. Candidate publication and global profile integration remain separate.
````

### FILE: `internal/businesspolicy/profile.go`
```yaml
block_id: "GO-BUSINESS-POLICY-PROFILE:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "Local explicit reference contract, configuration validation, composition/binding glue or tests; no external algorithm attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "01690c642562ce13877aa12d098e2bcd58d9025c69b51fc01d0e5ab94388b571"
variables: []
secrets_allowed: false
```
````go
package businesspolicy

// AUTHORED configuration-validation and contract-enforcement glue. No external
// algorithm or vendor authorship is claimed. SQL owner compatibility is explicit.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var ErrProfile = errors.New("business policy profile invalid, unbound or incompatible")

const MaximumProfileBytes = 16384

// These are compatibility limits of existing migration0007, not defaults.
const databaseMaximumSlotSeconds = 28800
const databaseMaximumSlotCapacity = 100

var profileID = regexp.MustCompile(`^[a-z][a-z0-9_-]{1,63}$`)

type document struct {
	Schema       string              `json:"schema"`
	Extends      string              `json:"extends"`
	ProfileID    string              `json:"profile_id"`
	Revision     int64               `json:"revision"`
	Appointments appointmentContract `json:"appointments"`
	Pricebooks   pricebookContract   `json:"pricebooks"`
}
type appointmentContract struct {
	LeadTimeSeconds     *int64   `json:"lead_time_seconds"`
	MaximumSlotSeconds  *int64   `json:"maximum_slot_seconds"`
	MinimumSlotCapacity *int     `json:"minimum_slot_capacity"`
	MaximumSlotCapacity *int     `json:"maximum_slot_capacity"`
	OccupyingStates     []string `json:"occupying_states"`
	WorkingWindow       string   `json:"working_window"`
	UnavailableWindow   string   `json:"unavailable_window"`
}
type pricebookContract struct {
	Selection string `json:"selection"`
	Validity  string `json:"validity"`
}

// Profile is immutable outside this package. Use Load with an independently
// selected SHA256, then give the same profile to both service and repository.
type Profile struct {
	hash  string
	value document
}

func validHash(value string) bool {
	decoded, err := hex.DecodeString(value)
	return err == nil && len(decoded) == sha256.Size && value == strings.ToLower(value)
}

func LoadFile(path, expectedSHA256 string) (*Profile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, ErrProfile
	}
	defer f.Close()
	raw, err := io.ReadAll(io.LimitReader(f, MaximumProfileBytes+1))
	if err != nil {
		return nil, ErrProfile
	}
	return Load(raw, expectedSHA256)
}

func Load(raw []byte, expectedSHA256 string) (*Profile, error) {
	if len(raw) == 0 || len(raw) > MaximumProfileBytes || !utf8.Valid(raw) || !validHash(expectedSHA256) {
		return nil, ErrProfile
	}
	sum := sha256.Sum256(raw)
	if hex.EncodeToString(sum[:]) != expectedSHA256 {
		return nil, ErrProfile
	}
	// encoding/json rejects unknown fields below but normally accepts duplicate
	// keys. Reject duplicates before decoding a contract that controls behavior.
	tokens := json.NewDecoder(bytes.NewReader(raw))
	if err := uniqueJSON(tokens, 0); err != nil {
		return nil, ErrProfile
	}
	if _, err := tokens.Token(); err != io.EOF {
		return nil, ErrProfile
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var value document
	if decoder.Decode(&value) != nil {
		return nil, ErrProfile
	}
	a := value.Appointments
	if value.Schema != "elite-business-policy/v1" || value.Extends != "PBC-CORE/1.0.0" || !profileID.MatchString(value.ProfileID) || value.Revision < 1 || a.LeadTimeSeconds == nil || a.MaximumSlotSeconds == nil || a.MinimumSlotCapacity == nil || a.MaximumSlotCapacity == nil {
		return nil, ErrProfile
	}
	if *a.LeadTimeSeconds < 0 || *a.LeadTimeSeconds > math.MaxInt64/int64(time.Second) || *a.MaximumSlotSeconds < 1 || *a.MaximumSlotSeconds > databaseMaximumSlotSeconds || *a.MinimumSlotCapacity < 1 || *a.MaximumSlotCapacity > databaseMaximumSlotCapacity || *a.MinimumSlotCapacity > *a.MaximumSlotCapacity {
		return nil, ErrProfile
	}
	// Changed modes would contradict existing lifecycle/triggers0007/0009/0015.
	// They require a separately admitted schema/owner change, not silent config.
	if len(a.OccupyingStates) != 2 || a.OccupyingStates[0] != "requested" || a.OccupyingStates[1] != "confirmed" || a.WorkingWindow != "contains_slot" || a.UnavailableWindow != "reject_overlap" || value.Pricebooks.Selection != "single_active_per_market_currency" || value.Pricebooks.Validity != "half_open" {
		return nil, ErrProfile
	}
	return &Profile{hash: expectedSHA256, value: value}, nil
}

func uniqueJSON(decoder *json.Decoder, depth int) error {
	if depth > 16 {
		return ErrProfile
	}
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		seen := map[string]bool{}
		for decoder.More() {
			key, err := decoder.Token()
			if err != nil {
				return err
			}
			name, ok := key.(string)
			if !ok || seen[name] || strings.ToLower(name) != name {
				return ErrProfile
			}
			seen[name] = true
			if err := uniqueJSON(decoder, depth+1); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil || end != json.Delim('}') {
			return ErrProfile
		}
	case '[':
		for decoder.More() {
			if err := uniqueJSON(decoder, depth+1); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil || end != json.Delim(']') {
			return ErrProfile
		}
	default:
		return ErrProfile
	}
	return nil
}

func (p *Profile) Valid() bool { return p != nil && validHash(p.hash) }
func (p *Profile) SHA256() string {
	if !p.Valid() {
		return ""
	}
	return p.hash
}
func (p *Profile) LeadTime() time.Duration {
	return time.Duration(*p.value.Appointments.LeadTimeSeconds) * time.Second
}
func (p *Profile) MaximumSlotDuration() time.Duration {
	return time.Duration(*p.value.Appointments.MaximumSlotSeconds) * time.Second
}
func (p *Profile) MinimumSlotCapacity() int { return *p.value.Appointments.MinimumSlotCapacity }
func (p *Profile) MaximumSlotCapacity() int { return *p.value.Appointments.MaximumSlotCapacity }

func (p *Profile) AllowsSlot(starts, ends, now time.Time, capacity int) bool {
	return p.Valid() && !starts.Before(now.Add(p.LeadTime())) && ends.After(starts) && ends.Sub(starts) <= p.MaximumSlotDuration() && capacity >= p.MinimumSlotCapacity() && capacity <= p.MaximumSlotCapacity()
}

// BindRequestHash preserves historical reference-profile idempotency receipts.
// Nonreference profiles domain-separate the same caller hash. A replay under a
// changed profile therefore conflicts instead of claiming the old decision was
// made using new policy. It does not rewrite any existing key or stored receipt.
func (p *Profile) BindRequestHash(requestHash string) (string, error) {
	if !p.Valid() || !validHash(requestHash) {
		return "", ErrProfile
	}
	if p.hash == ReferenceSHA256 {
		return requestHash, nil
	}
	sum := sha256.Sum256([]byte("business-policy-request/v1\x00" + p.hash + "\x00" + requestHash))
	return hex.EncodeToString(sum[:]), nil
}
````

### FILE: `internal/businesspolicy/profile_test.go`
```yaml
block_id: "GO-BUSINESS-POLICY-PROFILE:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "Local explicit reference contract, configuration validation, composition/binding glue or tests; no external algorithm attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "6e74ae5a8c32f129b5cf1941d0733d9ee80b6e165f9b9836687c6ecd2e7dcee2"
variables: []
secrets_allowed: false
```
````go
package businesspolicy

// AUTHORED contract tests. No vendor attribution.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func profileHash(raw []byte) string { sum := sha256.Sum256(raw); return hex.EncodeToString(sum[:]) }
func changedProfile(t *testing.T, old, new string) *Profile {
	t.Helper()
	raw := bytes.ReplaceAll(ReferenceJSON(), []byte(old), []byte(new))
	p, err := Load(raw, profileHash(raw))
	if err != nil {
		t.Fatal(err)
	}
	return p
}
func TestPolicyProfileBehaviorAndBindings(t *testing.T) {
	base := Reference()
	now := time.Date(2026, 9, 11, 12, 0, 0, 0, time.UTC)
	if base.SHA256() != ReferenceSHA256 || base.LeadTime() != 30*time.Minute || base.MaximumSlotDuration() != 8*time.Hour || base.MinimumSlotCapacity() != 1 || base.MaximumSlotCapacity() != 100 {
		t.Fatal("reference drift")
	}
	near := now.Add(time.Minute)
	if base.AllowsSlot(near, near.Add(time.Hour), now, 2) {
		t.Fatal("reference lead time bypass")
	}
	short := changedProfile(t, `"lead_time_seconds": 1800`, `"lead_time_seconds": 0`)
	if !short.AllowsSlot(near, near.Add(time.Hour), now, 2) {
		t.Fatal("compatible config did not change admission")
	}
	capTwo := changedProfile(t, `"maximum_slot_capacity": 100`, `"maximum_slot_capacity": 2`)
	start := now.Add(time.Hour)
	if !capTwo.AllowsSlot(start, start.Add(time.Hour), now, 2) || capTwo.AllowsSlot(start, start.Add(time.Hour), now, 3) {
		t.Fatal("capacity configuration ignored")
	}
	oneHour := changedProfile(t, `"maximum_slot_seconds": 28800`, `"maximum_slot_seconds": 3600`)
	if !oneHour.AllowsSlot(start, start.Add(time.Hour), now, 1) || oneHour.AllowsSlot(start, start.Add(time.Hour+time.Nanosecond), now, 1) {
		t.Fatal("duration boundary")
	}
	request := strings.Repeat("a", 64)
	historical, err := base.BindRequestHash(request)
	if err != nil || historical != request {
		t.Fatal("historical identity changed")
	}
	first, err := short.BindRequestHash(request)
	again, _ := short.BindRequestHash(request)
	other, _ := capTwo.BindRequestHash(request)
	if err != nil || first != again || first == historical || first == other {
		t.Fatal("policy request binding")
	}
	if _, err = short.BindRequestHash("invalid"); !errors.Is(err, ErrProfile) {
		t.Fatal("unbound request hash")
	}
	raw := ReferenceJSON()
	p, err := Load(raw, ReferenceSHA256)
	if err != nil {
		t.Fatal(err)
	}
	raw[0] = 'x'
	if !p.Valid() || ReferenceJSON()[0] != '{' {
		t.Fatal("caller mutated immutable profile")
	}
	path := filepath.Join(t.TempDir(), "policy.json")
	if err = os.WriteFile(path, ReferenceJSON(), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = LoadFile(path, ReferenceSHA256); err != nil {
		t.Fatal(err)
	}
	if _, err = LoadFile(path, first); !errors.Is(err, ErrProfile) {
		t.Fatal("file hash mismatch")
	}
}
func TestPolicyProfileRejectsInvalidOrIncompatible(t *testing.T) {
	ref := string(ReferenceJSON())
	cases := map[string]string{
		"unknown":             strings.Replace(ref, `"revision": 1`, `"revision": 1, "extra": true`, 1),
		"duplicate":           strings.Replace(ref, `"revision": 1`, `"revision": 1, "revision": 2`, 1),
		"case-alias":          strings.Replace(ref, `"revision": 1`, `"revision": 1, "Revision": 2`, 1),
		"missing":             strings.Replace(ref, `"lead_time_seconds": 1800,`, "", 1),
		"null":                strings.Replace(ref, `"lead_time_seconds": 1800`, `"lead_time_seconds": null`, 1),
		"negative":            strings.Replace(ref, `"lead_time_seconds": 1800`, `"lead_time_seconds": -1`, 1),
		"overflow":            strings.Replace(ref, `"lead_time_seconds": 1800`, `"lead_time_seconds": 9223372037`, 1),
		"duration":            strings.Replace(ref, `"maximum_slot_seconds": 28800`, `"maximum_slot_seconds": 28801`, 1),
		"capacity":            strings.Replace(ref, `"maximum_slot_capacity": 100`, `"maximum_slot_capacity": 101`, 1),
		"minmax":              strings.Replace(strings.Replace(ref, `"minimum_slot_capacity": 1`, `"minimum_slot_capacity": 3`, 1), `"maximum_slot_capacity": 100`, `"maximum_slot_capacity": 2`, 1),
		"occupancy":           strings.Replace(ref, `"confirmed"`, `"completed"`, 1),
		"occupancy-duplicate": strings.Replace(ref, `"confirmed"`, `"requested"`, 1),
		"working":             strings.Replace(ref, `"contains_slot"`, `"ignore"`, 1),
		"unavailable":         strings.Replace(ref, `"reject_overlap"`, `"ignore"`, 1),
		"exclusive":           strings.Replace(ref, `"single_active_per_market_currency"`, `"best_price"`, 1),
		"halfopen":            strings.Replace(ref, `"half_open"`, `"inclusive"`, 1),
		"trailing":            ref + "{}", "oversize": strings.Repeat(" ", MaximumProfileBytes+1), "utf8": ref + string([]byte{0xff}),
	}
	for name, value := range cases {
		t.Run(name, func(t *testing.T) {
			raw := []byte(value)
			if _, err := Load(raw, profileHash(raw)); !errors.Is(err, ErrProfile) {
				t.Fatal("invalid profile accepted", err)
			}
		})
	}
	for _, hash := range []string{"", strings.Repeat("0", 64), strings.ToUpper(ReferenceSHA256)} {
		if _, err := Load(ReferenceJSON(), hash); !errors.Is(err, ErrProfile) {
			t.Fatal("unverified SHA accepted")
		}
	}
	if (*Profile)(nil).Valid() || new(Profile).Valid() {
		t.Fatal("zero profile accepted")
	}
}
func FuzzPolicyProfile(f *testing.F) {
	f.Add(ReferenceJSON())
	f.Add([]byte(`{"schema":"elite-business-policy/v1","revision":1,"revision":2}`))
	f.Add([]byte("null"))
	f.Fuzz(func(t *testing.T, raw []byte) {
		p, err := Load(raw, profileHash(raw))
		if err == nil {
			if !p.Valid() || p.SHA256() != profileHash(raw) {
				t.Fatal("unbound accepted profile")
			}
			replay, e := Load(raw, p.SHA256())
			if e != nil || replay.SHA256() != p.SHA256() {
				t.Fatal("nondeterministic loader")
			}
		}
	})
}
````

### FILE: `internal/businesspolicy/reference-profile.json`
```yaml
block_id: "GO-BUSINESS-POLICY-PROFILE:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "Local explicit reference contract, configuration validation, composition/binding glue or tests; no external algorithm attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "da3feda4d66a19806f20234773720c2f98f6c7d4251640d6fb12cb58b24be5c5"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-business-policy/v1",
  "extends": "PBC-CORE/1.0.0",
  "profile_id": "library-reference",
  "revision": 1,
  "appointments": {
    "lead_time_seconds": 1800,
    "maximum_slot_seconds": 28800,
    "minimum_slot_capacity": 1,
    "maximum_slot_capacity": 100,
    "occupying_states": [
      "requested",
      "confirmed"
    ],
    "working_window": "contains_slot",
    "unavailable_window": "reject_overlap"
  },
  "pricebooks": {
    "selection": "single_active_per_market_currency",
    "validity": "half_open"
  }
}
````

### FILE: `internal/businesspolicy/reference.go`
```yaml
block_id: "GO-BUSINESS-POLICY-PROFILE:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "Local explicit reference contract, configuration validation, composition/binding glue or tests; no external algorithm attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "810cb077fc9b23065b7e936988cafd094342577b4b3857b6b4b768a70a65a396"
variables: []
secrets_allowed: false
```
````go
package businesspolicy

// AUTHORED reference-configuration binding. The original values live in JSON,
// not business logic. A changed embedded artifact fails startup validation.
import _ "embed"

//go:embed reference-profile.json
var referenceBytes []byte

const ReferenceSHA256 = "da3feda4d66a19806f20234773720c2f98f6c7d4251640d6fb12cb58b24be5c5"

func Reference() *Profile {
	p, err := Load(referenceBytes, ReferenceSHA256)
	if err != nil {
		panic("invalid embedded business policy profile")
	}
	return p
}

func ReferenceJSON() []byte { return append([]byte(nil), referenceBytes...) }
````

### FILE: `internal/franchisejourney/policy_profile_test.go`
```yaml
block_id: "GO-BUSINESS-POLICY-PROFILE:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "Local explicit reference contract, configuration validation, composition/binding glue or tests; no external algorithm attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "2d83dc5992210d47ba61da017c72360c3cf8842e423662445a672899ad27518d"
variables: []
secrets_allowed: false
```
````go
package franchisejourney

// AUTHORED profile-to-owner binding regression tests.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/businesspolicy"
	"fmt"
	"testing"
	"time"
)

type policyBoundFake struct {
	fakeRepository
	hash string
}

func (r *policyBoundFake) BusinessPolicySHA256() string { return r.hash }
func TestJourneyPolicyConfigurationAndBinding(t *testing.T) {
	raw := bytes.ReplaceAll(businesspolicy.ReferenceJSON(), []byte(`"lead_time_seconds": 1800`), []byte(`"lead_time_seconds": 0`))
	p, err := businesspolicy.Load(raw, fmt.Sprintf("%x", sha256.Sum256(raw)))
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 9, 11, 12, 0, 0, 0, time.UTC)
	repo := &policyBoundFake{hash: p.SHA256()}
	service, err := NewServiceWithProfile(repo, &fixedIDs{}, fixedClock{now}, p)
	if err != nil {
		t.Fatal(err)
	}
	slot := AppointmentSlot{OrganizationID: "org", Kind: "consultation", StartsAt: now.Add(time.Minute), EndsAt: now.Add(61 * time.Minute), Capacity: 2}
	if _, err = NewService(&fakeRepository{}, &fixedIDs{}, fixedClock{now}).CreateAppointmentSlot(context.Background(), "tenant", slot); err == nil {
		t.Fatal("reference lead time bypass")
	}
	if _, err = service.CreateAppointmentSlot(context.Background(), "tenant", slot); err != nil {
		t.Fatal("configuration ignored", err)
	}
	t.Run("legacy-constructor-rejects-custom-policy", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("legacy service allowed a different repository policy")
			}
		}()
		NewService(repo, &fixedIDs{}, fixedClock{now})
	})
	repo.hash = businesspolicy.ReferenceSHA256
	if _, err = NewServiceWithProfile(repo, &fixedIDs{}, fixedClock{now}, p); err == nil {
		t.Fatal("mismatched service/repository profile")
	}
	if _, err = NewServiceWithProfile(&fakeRepository{}, &fixedIDs{}, fixedClock{now}, p); err == nil {
		t.Fatal("repository without binding")
	}
	if _, err = NewServiceWithProfile(repo, &fixedIDs{}, fixedClock{now}, nil); err == nil {
		t.Fatal("nil profile")
	}
}
````

### FILE: `internal/platform/postgres/policy_profile_integration_test.go`
```yaml
block_id: "GO-BUSINESS-POLICY-PROFILE:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "Local explicit reference contract, configuration validation, composition/binding glue or tests; no external algorithm attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "3e96b414051fd13893861ced4b4442aca0f6a129600507106b4b866b71bb5e25"
variables: []
secrets_allowed: false
```
````go
package postgres

// AUTHORED delta tests; isolated PostgreSQL fixtures, no external provider.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/businesspolicy"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/randomid"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strings"
	"testing"
	"time"
)

func policyFixture(t *testing.T, old, new string) *businesspolicy.Profile {
	t.Helper()
	raw := bytes.ReplaceAll(businesspolicy.ReferenceJSON(), []byte(old), []byte(new))
	p, err := businesspolicy.Load(raw, fmt.Sprintf("%x", sha256.Sum256(raw)))
	if err != nil {
		t.Fatal(err)
	}
	return p
}
func TestBusinessPolicyPersistenceBindingsAndReplay(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	tenantCode := "policy-" + strings.ReplaceAll(tenant, "-", "")
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'policy-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into org.public_location(tenant_id,organization_id,city,region,country,published)values($1,'store','Synthetic','Synthetic','AR',true)`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','new','fixture','{}')`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	p := policyFixture(t, `"lead_time_seconds": 1800`, `"lead_time_seconds": 0`)
	repo, err := NewFranchiseJourneyWithProfile(pool, p)
	if err != nil {
		t.Fatal(err)
	}
	strict := policyFixture(t, `"maximum_slot_capacity": 100`, `"maximum_slot_capacity": 1`)
	strictRepo, _ := NewFranchiseJourneyWithProfile(pool, strict)
	if repo.BusinessPolicySHA256() != p.SHA256() {
		t.Fatal("repository policy binding")
	}
	if _, err = NewFranchiseJourneyWithProfile(pool, nil); err == nil {
		t.Fatal("invalid policy accepted")
	}
	start := time.Now().UTC().Add(10 * time.Minute).Truncate(time.Second)
	end := start.Add(time.Hour)
	if _, err = repo.CreateAvailability(ctx, tenant, "scheduler", franchisejourney.AvailabilityEntry{ID: "working", OrganizationID: "store", EntryType: "working", StartsAt: start, EndsAt: end.Add(4 * time.Hour)}, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	slot := franchisejourney.AppointmentSlot{ID: "slot", OrganizationID: "store", Kind: "consultation", StartsAt: start, EndsAt: end, Capacity: 2}
	if _, err = NewFranchiseJourney(pool).CreateAppointmentSlot(ctx, tenant, slot, randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("repository bypassed reference lead time", err)
	}
	tooMany := slot
	tooMany.ID = "too-many"
	tooMany.StartsAt = end
	tooMany.EndsAt = end.Add(time.Hour)
	if _, err = strictRepo.CreateAppointmentSlot(ctx, tenant, tooMany, randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("repository bypassed configured capacity", err)
	}
	short := policyFixture(t, `"maximum_slot_seconds": 28800`, `"maximum_slot_seconds": 3600`)
	shortRepo, _ := NewFranchiseJourneyWithProfile(pool, short)
	long := tooMany
	long.ID = "too-long"
	long.EndsAt = long.StartsAt.Add(2 * time.Hour)
	if _, err = shortRepo.CreateAppointmentSlot(ctx, tenant, long, randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("repository bypassed duration", err)
	}
	hash := strings.Repeat("a", 64)
	event := randomid.Generator{}.New()
	first, replay, err := repo.CreateAppointmentSlotOnce(ctx, tenant, "scheduler", "slot-key", hash, slot, event)
	if err != nil || replay || first.ID != slot.ID {
		t.Fatal(first, replay, err)
	}
	again, replay, err := repo.CreateAppointmentSlotOnce(ctx, tenant, "scheduler", "slot-key", hash, slot, randomid.Generator{}.New())
	if err != nil || !replay || again.ID != first.ID {
		t.Fatal("same policy replay", again, replay, err)
	}
	if _, _, err = strictRepo.CreateAppointmentSlotOnce(ctx, tenant, "scheduler", "slot-key", hash, slot, randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("changed policy replay accepted", err)
	}
	if items, e := repo.PublicAppointmentSlots(ctx, tenantCode, "store", "consultation", start.Add(-time.Minute), end); e != nil || len(items) != 1 {
		t.Fatal("custom policy listing ignored", len(items), e)
	}
	if items, e := NewFranchiseJourney(pool).PublicAppointmentSlots(ctx, tenantCode, "store", "consultation", start.Add(-time.Minute), end); e != nil || len(items) != 0 {
		t.Fatal("unbookable reference slot advertised", len(items), e)
	}
	appointment := franchisejourney.Appointment{ID: "appointment", LeadID: "lead", Kind: "consultation", StartsAt: start}
	if _, _, err = NewFranchiseJourney(pool).RequestAppointment(ctx, tenantCode, "store", "reference-booking", appointment, hash, randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("direct booking bypassed reference lead time", err)
	}
	booked, replay, err := repo.RequestAppointment(ctx, tenantCode, "store", "booking-key", appointment, hash, randomid.Generator{}.New())
	if err != nil || replay || booked.ID != appointment.ID {
		t.Fatal(booked, replay, err)
	}
	booked, replay, err = repo.RequestAppointment(ctx, tenantCode, "store", "booking-key", appointment, hash, randomid.Generator{}.New())
	if err != nil || !replay || booked.ID != appointment.ID {
		t.Fatal("booking replay", booked, replay, err)
	}
	if _, _, err = strictRepo.RequestAppointment(ctx, tenantCode, "store", "booking-key", appointment, hash, randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("booking policy mismatch", err)
	}
	var count int
	var bound, receipt string
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type in ('appointment-slot.created','appointment.requested')`, tenant).Scan(&count); err != nil || count != 2 {
		t.Fatal("duplicate or partial effects", count, err)
	}
	if err = pool.QueryRow(ctx, `select payload->>'policy_sha256' from platform.outbox_event where tenant_id=$1 and event_id=$2`, tenant, event).Scan(&bound); err != nil || bound != p.SHA256() {
		t.Fatal("event hash binding", bound, err)
	}
	expected, _ := p.BindRequestHash(hash)
	if err = pool.QueryRow(ctx, `select request_sha256_hex from platform.idempotency_record where tenant_id=$1 and scope='franchise-slot' and idempotency_key='slot-key'`, tenant).Scan(&receipt); err != nil || receipt != expected || receipt == hash {
		t.Fatal("receipt hash binding", receipt, err)
	}
	// Existing records are prospective: lowering maximum capacity does not rewrite slots.
	if err = pool.QueryRow(ctx, `select capacity from crm.appointment_slot where tenant_id=$1 and slot_id='slot'`, tenant).Scan(&count); err != nil || count != 2 {
		t.Fatal("configuration rewrote persisted quota")
	}
	// Narrow Commerce delta: exact profile binding in activation audit; no pricing/order suite rerun.
	cr, err := NewCommerceWithProfile(pool, p)
	if err != nil || cr.BusinessPolicySHA256() != p.SHA256() {
		t.Fatal("commerce binding", err)
	}
	if _, err = NewCommerceWithProfile(pool, nil); err == nil {
		t.Fatal("commerce invalid profile")
	}
	for _, q := range []string{
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	b := commerce.PriceBook{ID: "book", Market: "AR", Currency: "ARS", ValidFrom: time.Now().Add(-time.Hour), Status: "draft", Entries: []commerce.PriceEntry{{VariantID: "variant", AmountMinorUnits: 1000, TaxMode: "inclusive"}}}
	if err = cr.CreatePriceBook(ctx, tenant, randomid.Generator{}.New(), b); err != nil {
		t.Fatal(err)
	}
	activation := randomid.Generator{}.New()
	if err = cr.ActivatePriceBook(ctx, tenant, b.ID, activation); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select payload->>'policy_sha256' from platform.outbox_event where tenant_id=$1 and event_id=$2`, tenant, activation).Scan(&bound); err != nil || bound != p.SHA256() {
		t.Fatal("pricebook policy audit", bound, err)
	}
}
````

### FILE: `cmd/electromobility-api/business_policy.go`
```yaml
block_id: "GO-BUSINESS-POLICY-PROFILE:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "Local explicit reference contract, configuration validation, composition/binding glue or tests; no external algorithm attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "f826e5cc4c77bd611f7b910e23c74d3957c715df7ceda08c19363ada8a04e093"
variables: []
secrets_allowed: false
```
````go
package main

// AUTHORED configuration and owner-construction glue. No vendor attribution.
import (
	"elite.local/enterprise/internal/businesspolicy"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func selectedBusinessPolicy(lookup func(string) string) (*businesspolicy.Profile, error) {
	if lookup == nil {
		return nil, businesspolicy.ErrProfile
	}
	path, hash := lookup("BUSINESS_POLICY_PROFILE_FILE"), lookup("BUSINESS_POLICY_PROFILE_SHA256")
	if path == "" && hash == "" {
		return businesspolicy.Reference(), nil
	}
	if path == "" || hash == "" {
		return nil, businesspolicy.ErrProfile
	}
	return businesspolicy.LoadFile(path, hash)
}

type businessPolicyRuntime struct {
	profile            *businesspolicy.Profile
	commerceRepository *postgres.Commerce
	journeyRepository  *postgres.FranchiseJourney
	commerceService    *commerce.Service
	journeyService     *franchisejourney.Service
}

func prepareBusinessPolicyRuntime(pool *pgxpool.Pool, ids franchisejourney.IDGenerator, clock franchisejourney.Clock, profile *businesspolicy.Profile) (*businessPolicyRuntime, error) {
	commerceRepository, err := postgres.NewCommerceWithProfile(pool, profile)
	if err != nil {
		return nil, err
	}
	journeyRepository, err := postgres.NewFranchiseJourneyWithProfile(pool, profile)
	if err != nil {
		return nil, err
	}
	journeyService, err := franchisejourney.NewServiceWithProfile(journeyRepository, ids, clock, profile)
	if err != nil {
		return nil, err
	}
	return &businessPolicyRuntime{
		profile: profile, commerceRepository: commerceRepository, journeyRepository: journeyRepository,
		commerceService: commerce.NewService(commerceRepository, ids), journeyService: journeyService,
	}, nil
}
````

### FILE: `cmd/electromobility-api/business_policy_test.go`
```yaml
block_id: "GO-BUSINESS-POLICY-PROFILE:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "Local explicit reference contract, configuration validation, composition/binding glue or tests; no external algorithm attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "37f00cdbeb90fe13df273a724cd45557f04d09da28f55a5e8c4152c831249c9c"
variables: []
secrets_allowed: false
```
````go
package main

// AUTHORED startup tests. No listener, database connection or provider is used.
import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/businesspolicy"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func policyEnvironment(values map[string]string) func(string) string {
	return func(key string) string { return values[key] }
}

func TestBusinessPolicyHostReferenceAndCompleteFile(t *testing.T) {
	reference, err := selectedBusinessPolicy(policyEnvironment(nil))
	if err != nil || reference.SHA256() != businesspolicy.ReferenceSHA256 || reference.LeadTime() != 30*time.Minute {
		t.Fatal("default reference changed", err)
	}
	raw := bytes.ReplaceAll(businesspolicy.ReferenceJSON(), []byte(`"lead_time_seconds": 1800`), []byte(`"lead_time_seconds": 0`))
	hash := fmt.Sprintf("%x", sha256.Sum256(raw))
	file := filepath.Join(t.TempDir(), "business-policy.json")
	if err = os.WriteFile(file, raw, 0600); err != nil {
		t.Fatal(err)
	}
	custom, err := selectedBusinessPolicy(policyEnvironment(map[string]string{"BUSINESS_POLICY_PROFILE_FILE": file, "BUSINESS_POLICY_PROFILE_SHA256": hash}))
	if err != nil || custom.SHA256() != hash || custom.LeadTime() != 0 {
		t.Fatal("custom profile not selected", err)
	}
	// A zero pool is a nonconnecting handle: constructors must not perform I/O.
	runtime, err := prepareBusinessPolicyRuntime(&pgxpool.Pool{}, randomid.Generator{}, systemClock{}, custom)
	if err != nil || runtime == nil || runtime.profile != custom || runtime.commerceService == nil || runtime.journeyService == nil {
		t.Fatal("owner construction", err)
	}
	if runtime.commerceRepository.BusinessPolicySHA256() != custom.SHA256() || runtime.journeyRepository.BusinessPolicySHA256() != custom.SHA256() {
		t.Fatal("owners received different profiles")
	}
}

func TestBusinessPolicyHostRejectsIncompleteInvalidAndMissingFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "business-policy.json")
	if err := os.WriteFile(file, businesspolicy.ReferenceJSON(), 0600); err != nil {
		t.Fatal(err)
	}
	cases := map[string]map[string]string{
		"path-only":       {"BUSINESS_POLICY_PROFILE_FILE": file},
		"hash-only":       {"BUSINESS_POLICY_PROFILE_SHA256": businesspolicy.ReferenceSHA256},
		"wrong-hash":      {"BUSINESS_POLICY_PROFILE_FILE": file, "BUSINESS_POLICY_PROFILE_SHA256": strings.Repeat("0", 64)},
		"missing-file":    {"BUSINESS_POLICY_PROFILE_FILE": file + ".absent", "BUSINESS_POLICY_PROFILE_SHA256": businesspolicy.ReferenceSHA256},
		"whitespace-path": {"BUSINESS_POLICY_PROFILE_FILE": " ", "BUSINESS_POLICY_PROFILE_SHA256": businesspolicy.ReferenceSHA256},
	}
	for name, values := range cases {
		t.Run(name, func(t *testing.T) {
			if p, err := selectedBusinessPolicy(policyEnvironment(values)); err == nil || p != nil {
				t.Fatal("invalid configuration fell back to reference")
			}
		})
	}
	invalid := bytes.ReplaceAll(businesspolicy.ReferenceJSON(), []byte(`"maximum_slot_capacity": 100`), []byte(`"maximum_slot_capacity": 101`))
	if err := os.WriteFile(file, invalid, 0600); err != nil {
		t.Fatal(err)
	}
	if p, err := selectedBusinessPolicy(policyEnvironment(map[string]string{"BUSINESS_POLICY_PROFILE_FILE": file, "BUSINESS_POLICY_PROFILE_SHA256": fmt.Sprintf("%x", sha256.Sum256(invalid))})); err == nil || p != nil {
		t.Fatal("hash-matching incompatible policy admitted")
	}
	if _, err := selectedBusinessPolicy(nil); err == nil {
		t.Fatal("nil environment accessor")
	}
	for _, p := range []*businesspolicy.Profile{nil, new(businesspolicy.Profile)} {
		if runtime, err := prepareBusinessPolicyRuntime(&pgxpool.Pool{}, randomid.Generator{}, systemClock{}, p); err == nil || runtime != nil {
			t.Fatal("invalid profile reached constructed owners")
		}
	}
	if _, err := prepareBusinessPolicyRuntime(nil, randomid.Generator{}, systemClock{}, businesspolicy.Reference()); err == nil {
		t.Fatal("nil persistence handle")
	}
	if _, err := prepareBusinessPolicyRuntime(&pgxpool.Pool{}, nil, systemClock{}, businesspolicy.Reference()); err == nil {
		t.Fatal("nil id generator")
	}
	if _, err := prepareBusinessPolicyRuntime(&pgxpool.Pool{}, randomid.Generator{}, nil, businesspolicy.Reference()); err == nil {
		t.Fatal("nil clock")
	}
}
````

## 6. Configuration surface

BUSINESS_POLICY_PROFILE_FILE and BUSINESS_POLICY_PROFILE_SHA256 are both absent
for the embedded reference; both required for a reviewed custom JSON revision.
No configuration data is printed on failure. See the materialized contract for
the reference digest and exact compatibility bounds. Both fields are non-secret.

## 7. Dependency bill

Go standard library only for the loader; host glue uses already selected Commerce,
Journey and pgx owners. go.mod/go.sum, runtimes, migrations and upstream locks do
not change. No new upstream is acquired or admitted by this pack.

## 8. Apply order

Compose this pack with Commerce0.6.4 and Journey0.10.21, then apply the provided
main/environment changes in their existing owners. Avoid duplicate file ownership
for the two new host helper files, which belong to this pack. Select one reviewed
profile revision across instances before serving. Keep previous profile bytes,
hash and receipts when reverting configuration; do not rewrite idempotency rows.
Requests replayed under another profile fail closed; existing scoped result reads
remain the recovery surface. There is no automatic live rollout/control plane.

## 9. Verification

Focused configuration invalid/behavior/hash/replay tests, PostgreSQL18.6 with54
unchanged migrations, Go parser fuzz3seconds/53,626executions. Exact implementation
hashes reuse the prior PG/fuzz receipts; no PG rerun for host integration. Two host
unit tests, six executable startup cases, host vet/build pass. G0-G8 evidence is
policy-host-evidence.md in the staging review bundle. No whole-owner source,
live target, SAST/DAST, production-load or financial policy certification.

## 10. Reconstruction evidence

Materialize these three candidates in absent destinations; compare every output
with the candidate block hash and changed/new files with D. Unchanged blocks
retain the canonical owner bytes. The parent must also merge its independently
pending payment/handover host overlays; these candidates do not overwrite them.

Canonical V402 integration: selected by the current profile with exact dependencies and caller overlays. Metadata promotion records byte reconstruction, not closure of every admission/release gate. Payload provenance is unchanged.
