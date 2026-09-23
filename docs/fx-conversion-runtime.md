# Exact supplied-snapshot FX conversion receipts

This infrastructure converts declared integer minor units using the selected immutable rates. It stores a conversion receipt under the existing accounting owner and commits its outbox event in the same transaction. It does not post journals, settle balances, collect fees or choose rates for a user.

The supported profile is `BC_DIRECT_BASE_V1`: each foreign currency has a positive decimal Exchange Rate Amount and Relational Exchange Rate Amount against an explicit local currency. An explicit conversion date selects the latest eligible starting date. Relational currency graphs are rejected. Currency codes must be declared in the profile; this is not a built-in ISO4217 catalog. Destination decimals and positive integral minor-unit rounding quantum are explicit. The separate mathematical rounding profile is `NEAREST_TIES_AWAY_FROM_ZERO`; same-currency and zero conversions retain the upstream early-return behavior. See the documented negative-example discrepancy and limited source claim in `provenance/BC_FX_DERIVATION.md`.

The profile bytes are externally SHA256-locked. They identify the source document by its SHA256, ID and revision. The source document carries immutable row identities, dates and exact decimal strings. Both documents reject unknown/duplicate JSON keys, duplicate identities/currency dates, missing decimals, invalid precision, unsupported relationships and oversized/deep content. The active snapshot owns private copied bytes and values. Its availability window is checked at startup and during each new conversion, including after an outbox write that waited. The requested conversion date and maximum age of the selected rates are independent explicit bounds.

`config/fx/reference-profile.json` and `reference-rates.json` contain synthetic examples; the default host does not activate them. A materialized project selects its authorized source data and profile. The byte lock proves which supplied data was used, not its economic truth. No live account, secret, rate recommendation or automatic feed is required by this supplied-snapshot lane.

Set these non-secret activation fields in the selected runtime configuration:

- `FX_ENABLED=true`
- `FX_PROFILE_FILE`, `FX_RATES_FILE`: paths to the exact supplied documents.
- `FX_PROFILE_ID`, `FX_PROFILE_REVISION`, `FX_PROFILE_SHA256`: exact activated profile identity and bytes.
- `FX_TENANT_ID`, `FX_ORGANIZATION_ID`, `FX_LOCAL_CURRENCY`: deployment scope, matched to the profile.

Missing/inconsistent/expired activation fails startup. `FX_ENABLED=false` mounts no FX routes. Rate refresh is explicit replacement with a newly reviewed profile revision/hash and source bytes; it never silently changes historical receipts. Reusing an existing profile ID/revision with different bytes is rejected by PostgreSQL.

The API requires the existing verified bearer identity. POST requires `accounting:write`; GET requires `accounting:read`. Both enforce organization scope and the activated tenant/organization. Recovery is bound to the original principal and request key. Money values are JSON strings of integer minor units, avoiding a binary64 boundary.

```http
POST /v1/accounting/fx/conversions
Authorization: Bearer <runtime identity>
Content-Type: application/json
Idempotency-Key: <stable 16-128 character key>

{"organization_id":"store","from_currency":"EUR","to_currency":"GBP","amount_minor":"1000","conversion_date":"2026-01-01"}
```

The reply is an immutable receipt with exact input/output, rational result, selected rate identities/dates, precision, source/profile identities, actor and recording time. The effect is always `CONVERSION_RECEIPT_ONLY`. A repeated identical request returns the same receipt; a divergent request or principal under that key fails. After a lost response, recover through:

```http
GET /v1/accounting/fx/conversions/result?organization_id=store
Authorization: Bearer <same principal identity>
Idempotency-Key: <original key>
```

Recovery is historical; it does not claim current rate validity or authorize money movement. A retained receipt prevents duplicates even after cleanup of the common transient idempotency record. No second ledger, queue or independent idempotency service exists.

Apply migration0061 with the existing accounting schema and platform idempotency/outbox. Its down migration is a destructive rollback for an empty/rejected deployment; do not run it over retained business receipts without the project's recovery procedure. Source snapshots and receipts use the existing accounting immutability trigger. The reference fixture uses only a fresh loopback PostgreSQL database named `elite_fx_*`, selected by `FX_CONNECTED_DB_URL`.

Focused gates use the admitted Go runtime with module downloads disabled:

```text
go test -mod=readonly ./internal/bcfx -run "^(TestSnapshot.*|TestRoundMinorIndependentDecimalOracle)$" -count=1
go test -mod=readonly ./cmd/electromobility-api -run "^TestFXHostActivationExactSnapshot$" -count=1
go test -mod=readonly ./internal/platform/postgres -run "^TestFXConnectedHTTPReceipts$" -count=1 -v
go test -mod=readonly ./internal/bcfx -run "^$" -fuzz "^FuzzSnapshotExactByteBoundary$" -fuzztime=3s -parallel=2
```

The connected test reuses the reference composition's local `handoverBrowserIssuer` helper for genuine RS256/JWKS verification. It does not rerun the handover journey. Missing database configuration explicitly skips; acceptance must reject skips. The retained 525-vector independent oracle can be regenerated into an absent file with `tools/generate_fx_rounding_oracle.py --output <new-path>` and compared byte for byte.

No new Go module dependency is introduced. These local gates do not close the reference composition's pending global SCA/security/release gates or authorize production.
