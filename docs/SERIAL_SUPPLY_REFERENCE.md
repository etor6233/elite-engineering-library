# Connected serial supply reference — J2

Scope: LIBRARY_INFRASTRUCTURE, strict-serial-reference/v1. AUTHORED orchestration
of existing procurement, factory, shared human approval and inventory owners.
No new pricing or stock ledger. No third-party authorship is claimed for this
binding. BC-derived inventory owners retain their own provenance and notices.

The existing draft purchase order must already identify its active supplier,
destination, currency and total. Bind a demand reference, active factory and
1–32 unique variant lines, total1–1000 whole serialized units. The reference
permits no overreceipt. Its purchase total is unchanged; no unit price is inferred.
Missing VIN/battery is permitted; all present identifiers and serials remain unique.

## Activation

Apply all migrations in order through0072. Enable SERIAL_SUPPLY_ENABLED=true
and set SERIAL_SUPPLY_POLICY_SHA256=02767ecb3efb11f8ddaed3bdf237f3111024ad175e6001526f8f028c3171b8fa.
The SHA binds exact UTF-8 PolicyJSON bytes in internal/serialsupply/contract.go,
without an added newline. This explicitly selects the built-in synthetic policy:
{"schema":"elite-serial-supply-policy/v1","scope":"LIBRARY_INFRASTRUCTURE_REFERENCE","quantity":"whole-serialized-units","overreceipt":0,"receipt_state":"quarantine","release":"distinct-human-with-evidence","pricing":"existing-purchase-order-total-unchanged","production_authorized":false}

Disabled mode reads no policy. Enabled mode requires a pool, exact policy hash
and15enabled database guards. No credential, remote endpoint or account is used.
The existing host supplies database configuration and the selected identity
verifier. A target later assigns real organization-scoped permissions explicitly.

## Connected flow and roles

Every route includes organization_id once. Identity comes from the verifier;
request bodies cannot choose actor or tenant. Each route narrows a broad principal
to its one permission and selected organization. Versions are JSON decimal strings.

POST /v1/franchise/supply/orders/{id}/plan binds the immutable plan (supply:plan).
Franchise submit/cancel use supply:plan. Factory confirm/start/register/milestone/
ship use supply:factory under the bound factory organization. A released or rejected
quality milestone requires a distinct human and evidence. Rejected serial history
is retained; a replacement may consume the released plan slot.

Ship writes the immutable ASN/manifest and existing stock in-transit entries.
Ship/receive each take1–100 unique units, allowing split shipments and receipts.
Franchise receive requires supply:receive and manifest membership; stock enters
quarantine, excluded from existing available-to-promise. quality/quality-reject
require supply:release and a distinct reviewer. A rejected receipt stays quarantined.
reinspect requires supply:inspect, a prior rejection and a new evidence hash.
Only approved receipt quality makes the existing stock available. After approval,
ordinary existing commerce/transfer/service transitions remain authoritative.

GET /v1/{franchise|factory}/supply/orders/{id} uses supply:read or
supply:factory-read. Units are bounded100/page, with next_unit_id/after_unit.
GET with command_id recovers an immutable command receipt. Following an uncertain
POST, compare actor, purchase order, command ID and request SHA before deciding
whether an explicit identical retry is appropriate. Never generate a new command
ID merely because a response was lost. Actor or payload changes cannot replay.

All effects, approval decisions and outbox rows commit together. Deferred database
guards reject generic API/SQL bypass before release. Down0072 refuses populated
history; on an empty database down/up is supported. It never erases audit facts.

## Evidence and limits

SERIAL_SUPPLY_CONNECTED_RELEASE_V402.md/json in the canonical library binds exact
source, G0–G8 and reconstruction. Actual PG18.6 fixture: planned3, registered4,
one rejected/replaced, two ASNs, three received/quarantined/reviewed/available;
30versions. Eight concurrent identical commands: one new effect and seven replays.
Outbox failure rolls back15table snapshots. Two committed HTTP responses lost and
recovered by GET without another accepted POST. Generic approval/stock/PO bypass
rejected. Host policy/guard validation,13transport boundaries, populated/empty
downgrade and finite identity fuzz pass. JWT/provider/live factory certification
and role UI are not inferred; UI is the separate T2804 work item.
