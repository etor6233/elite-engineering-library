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
