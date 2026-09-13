# Go Exact FX Snapshot Accounting Receipt

## 1. Metadata

```yaml
pack_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Immutable direct-base FX snapshot and exact dated conversion, plus typed durable draft binding to the existing accounting posting/reversal lifecycle. Scoped HTTP, shared idempotency/outbox and explicit host. No live feed, automatic valuation policy, AL runtime or global release approval."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "pgx 5.10.0"]
compatible_with: ["Reference composition with accounting0011, platform idempotency/outbox, verified identity and handoverBrowserIssuer test fixture"]
incompatible_with: ["Relational-currency graphs", "Binary64 monetary boundary", "Unreviewed or expired source snapshots", "GO-FX-CORE as conversion authority"]
license_expression: "LicenseRef-Workspace-Owner AND MIT"
upstream_sources: ["https://github.com/microsoft/BCApps/tree/2eae56d704a1fd035d104f333602aea7091b7749", "https://learn.microsoft.com/en-us/dynamics365/business-central/dev-itpro/developer/methods-auto/system/system-round-method"]
verified_at: "2026-09-12"
```

## 2. Applicability

Use only the declared direct-base profile. ExchangeExact and FindLast are ADAPTED from fixed BCApps methods. RoundMinor is AUTHORED exact representation/rounding glue under an explicitly selected mathematical profile, with the official documentation discrepancy recorded. The source/notice lock and derivation are materialized below. Remaining snapshot/accounting/HTTP/host/test code is AUTHORED integration glue. Supplied data is hash-bound, not certified as an economic rate recommendation.

Local snapshot/oracle/host/fuzz, actual HTTP/RS256/JWKS/PostgreSQL concurrency/recovery/expiry, vet/build and exact materialization gates pass. G5-G8 global composition SCA/security/release remain pending the root frozen-composition gates; no current zero-SCA claim is made here. No new Go dependency is introduced. Future live/target acceptance is separate.

## 3. Architecture contract

The selected accounting owner remains singular. Migration0061 adds immutable snapshots and conversion receipts, reusing the existing accounting immutability trigger and platform idempotency/outbox. Migration0065 adds an immutable conversion-to-journal binding. The existing accounting owner supplies the draft writer in the same transaction and retains its separate posting/reversal permissions and algorithms. No parallel ledger is introduced. The existing application owner must include the provided selectedFXConversionModule call and optional environment/notice overlays recorded separately by prehash. They are excluded from this pack to avoid duplicate main/environment ownership. The shared exact Microsoft MIT notice may already be materialized and must remain byte-identical.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/fx.go
CREATE cmd/electromobility-api/fx_test.go
CREATE config/fx/reference-profile.json
CREATE config/fx/reference-rates.json
CREATE db/migrations/0061_fx_conversion_receipt.down.sql
CREATE db/migrations/0061_fx_conversion_receipt.up.sql
CREATE docs/fx-conversion-runtime.md
CREATE docs/provenance/BC_FX_DERIVATION.md
CREATE docs/provenance/BC_FX_SOURCE_LOCK.json
CREATE docs/provenance/BC_FX_THIRD_PARTY_NOTICES.md
CREATE internal/accounting/fx.go
CREATE internal/bcfx/exchange.go
CREATE internal/bcfx/exchange_test.go
CREATE internal/bcfx/rounding.go
CREATE internal/bcfx/selection.go
CREATE internal/bcfx/snapshot.go
CREATE internal/bcfx/snapshot_test.go
CREATE internal/bcfx/testdata/rounding-oracle.json
CREATE internal/platform/httpapi/fx_conversion.go
CREATE internal/platform/postgres/fx_connected_integration_test.go
CREATE internal/platform/postgres/fx_conversion.go
CREATE licenses/Microsoft-BCApps-MIT.txt
CREATE tools/generate_fx_rounding_oracle.py
CREATE internal/accounting/fx_journal.go
CREATE internal/platform/postgres/fx_journal.go
CREATE internal/platform/httpapi/fx_journal.go
CREATE internal/platform/postgres/fx_journal_integration_test.go
CREATE db/migrations/0065_fx_journal_receipt.up.sql
CREATE db/migrations/0065_fx_journal_receipt.down.sql
CREATE docs/fx-journal-runtime.md
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/fx.go`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:cmd/electromobility-api/fx.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "240b6bbb3b71f0f917b750316d0f5898b734793cd578fdc2d522f1cef9995b75"
variables: []
secrets_allowed: false
```
````go
package main

// AUTHORED explicit FX host activation. Disabled means no FX route is mounted.
import (
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcfx"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"time"
)

var errFXConfiguration = errors.New("FX snapshot activation is invalid")

func loadFXSnapshot(lookup func(string) string, now time.Time) (*bcfx.Snapshot, error) {
	if lookup == nil {
		return nil, errFXConfiguration
	}
	switch lookup("FX_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errFXConfiguration
	}
	revision, e := strconv.Atoi(lookup("FX_PROFILE_REVISION"))
	if e != nil || revision < 1 {
		return nil, errFXConfiguration
	}
	s, e := bcfx.LoadSnapshotFiles(lookup("FX_PROFILE_FILE"), lookup("FX_RATES_FILE"), bcfx.Activation{Enabled: true, ProfileID: lookup("FX_PROFILE_ID"), Revision: revision, ProfileSHA256: lookup("FX_PROFILE_SHA256"), TenantID: lookup("FX_TENANT_ID"), OrganizationID: lookup("FX_ORGANIZATION_ID"), LocalCurrency: lookup("FX_LOCAL_CURRENCY")})
	if e != nil || !s.Current(now) {
		return nil, errFXConfiguration
	}
	return s, nil
}
func selectedFXConversionModule(pool *pgxpool.Pool, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	snapshot, e := loadFXSnapshot(lookup, time.Now())
	if e != nil || snapshot == nil {
		return nil, e
	}
	if pool == nil {
		return nil, errFXConfiguration
	}
	service, e := accounting.NewFXService(postgres.NewAccounting(pool), randomid.Generator{}, snapshot)
	if e != nil {
		return nil, errFXConfiguration
	}
	return httpapi.FXConversionModule{Service: service}, nil
}
````

### FILE: `cmd/electromobility-api/fx_test.go`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:cmd/electromobility-api/fx_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "1b7c28ea0fa90bc3477a1a5f50bfadc29c73b0f6e24ab7539fb03104eaf89dc7"
variables: []
secrets_allowed: false
```
````go
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFXHostActivationExactSnapshot(t *testing.T) {
	p := filepath.Join("..", "..", "config", "fx", "reference-profile.json")
	r := filepath.Join("..", "..", "config", "fx", "reference-rates.json")
	raw, e := os.ReadFile(p)
	if e != nil {
		t.Fatal(e)
	}
	sum := sha256.Sum256(raw)
	config := map[string]string{"FX_ENABLED": "true", "FX_PROFILE_FILE": p, "FX_RATES_FILE": r, "FX_PROFILE_ID": "fx-reference", "FX_PROFILE_REVISION": "1", "FX_PROFILE_SHA256": hex.EncodeToString(sum[:]), "FX_TENANT_ID": "018f4d4a-7b36-7a21-8d10-2f4c54c28101", "FX_ORGANIZATION_ID": "store", "FX_LOCAL_CURRENCY": "USD"}
	now := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	lookup := func(k string) string { return config[k] }
	if s, e := loadFXSnapshot(lookup, now); e != nil || s == nil {
		t.Fatal(e)
	}
	for _, key := range []string{"FX_PROFILE_FILE", "FX_RATES_FILE", "FX_PROFILE_ID", "FX_PROFILE_REVISION", "FX_PROFILE_SHA256", "FX_TENANT_ID", "FX_ORGANIZATION_ID", "FX_LOCAL_CURRENCY"} {
		saved := config[key]
		config[key] = ""
		if _, e := loadFXSnapshot(lookup, now); e == nil {
			t.Fatal("incomplete", key)
		}
		config[key] = saved
	}
	if _, e := loadFXSnapshot(lookup, time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)); e == nil {
		t.Fatal("expired")
	}
	config["FX_ENABLED"] = "false"
	if v, e := selectedFXConversionModule(nil, lookup); e != nil || v != nil {
		t.Fatal("disabled mounted")
	}
	config["FX_ENABLED"] = "invalid"
	if _, e := loadFXSnapshot(lookup, now); e == nil {
		t.Fatal("flag")
	}
}
````

### FILE: `config/fx/reference-profile.json`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:config/fx/reference-profile.json:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "4686c982987ae98d21579b6c1b2de0cea2ad6de67eb9f7bfce0b853129e06801"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-fx-direct-base-profile/v1",
  "profile_id": "fx-reference",
  "revision": 1,
  "algorithm": "BC_DIRECT_BASE_V1",
  "rounding": "NEAREST_TIES_AWAY_FROM_ZERO",
  "tenant_id": "018f4d4a-7b36-7a21-8d10-2f4c54c28101",
  "organization_id": "store",
  "local_currency": "USD",
  "valid_from": "2026-01-01T00:00:00Z",
  "valid_until": "2027-01-01T00:00:00Z",
  "conversion_date_from": "2026-01-01",
  "conversion_date_through": "2026-12-31",
  "maximum_rate_age_days": 366,
  "source_id": "synthetic-reference",
  "source_revision": "fixture-1",
  "source_sha256": "e8e8f92a1841e979b090ed5c2e1657dd3575b9e7d418969f72c4781a6cd22b18",
  "currencies": [
    {
      "code": "USD",
      "minor_unit_decimals": 2,
      "rounding_precision_minor": "1"
    },
    {
      "code": "EUR",
      "minor_unit_decimals": 2,
      "rounding_precision_minor": "5"
    },
    {
      "code": "GBP",
      "minor_unit_decimals": 3,
      "rounding_precision_minor": "1"
    }
  ]
}
````

### FILE: `config/fx/reference-rates.json`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:config/fx/reference-rates.json:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "e8e8f92a1841e979b090ed5c2e1657dd3575b9e7d418969f72c4781a6cd22b18"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-fx-direct-base-rates/v1",
  "source_id": "synthetic-reference",
  "source_revision": "fixture-1",
  "kind": "SUPPLIED_SNAPSHOT",
  "rates": [
    {
      "rate_id": "eur-2026-01-01",
      "currency": "EUR",
      "starting_date": "2026-01-01",
      "exchange_rate_amount": "100",
      "relational_exchange_rate_amount": "125",
      "relational_currency": ""
    },
    {
      "rate_id": "eur-2026-01-02",
      "currency": "EUR",
      "starting_date": "2026-01-02",
      "exchange_rate_amount": "100",
      "relational_exchange_rate_amount": "130",
      "relational_currency": ""
    },
    {
      "rate_id": "gbp-2026-01-01",
      "currency": "GBP",
      "starting_date": "2026-01-01",
      "exchange_rate_amount": "200",
      "relational_exchange_rate_amount": "300",
      "relational_currency": ""
    }
  ]
}
````

### FILE: `db/migrations/0061_fx_conversion_receipt.down.sql`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:db/migrations/0061_fx_conversion_receipt.down.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "b5253ec3744884881186b7624a59d57646e7c501996f0d9c8bb6b5b5ea866375"
variables: []
secrets_allowed: false
```
````sql
begin;
drop table accounting.fx_conversion_receipt;
drop table accounting.fx_rate_snapshot;
commit;
````

### FILE: `db/migrations/0061_fx_conversion_receipt.up.sql`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:db/migrations/0061_fx_conversion_receipt.up.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "74f7777537eebcdaeac95dfdb32c5cdb1ed7375ab26e7f8810d99fc6c2f3d471"
variables: []
secrets_allowed: false
```
````sql
begin;
-- AUTHORED records at the existing accounting owner. No new ledger or posting.
create table accounting.fx_rate_snapshot(
 tenant_id uuid not null,
 profile_id text not null,
 profile_revision integer not null check(profile_revision>0),
 organization_id text not null,
 profile_sha256_hex text not null check(profile_sha256_hex ~ '^[0-9a-f]{64}$'),
 source_sha256_hex text not null check(source_sha256_hex ~ '^[0-9a-f]{64}$'),
 profile_raw bytea not null check(octet_length(profile_raw) between 1 and 1048576),
 source_raw bytea not null check(octet_length(source_raw) between 1 and 1048576),
 valid_from timestamptz not null,
 valid_until timestamptz not null check(valid_until>valid_from),
 primary key(tenant_id,profile_id,profile_revision),
 unique(tenant_id,profile_id,profile_revision,profile_sha256_hex),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id)
);
create trigger fx_rate_snapshot_immutable before update or delete on accounting.fx_rate_snapshot for each row execute function accounting.prevent_posted_history_mutation();
create table accounting.fx_conversion_receipt(
 tenant_id uuid not null,
 conversion_id text not null,
 organization_id text not null,
 requested_by_subject text not null check(length(requested_by_subject) between 1 and 256),
 request_key text not null check(request_key ~ '^[A-Za-z0-9_-]{16,128}$'),
 request_sha256_hex text not null check(request_sha256_hex ~ '^[0-9a-f]{64}$'),
 profile_id text not null,
 profile_revision integer not null,
 profile_sha256_hex text not null,
 receipt_raw bytea not null check(octet_length(receipt_raw) between 1 and 16384),
 receipt_sha256_hex text not null check(receipt_sha256_hex ~ '^[0-9a-f]{64}$'),
 recorded_at timestamptz not null,
 primary key(tenant_id,conversion_id),
 unique(tenant_id,request_key),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 foreign key(tenant_id,profile_id,profile_revision,profile_sha256_hex) references accounting.fx_rate_snapshot(tenant_id,profile_id,profile_revision,profile_sha256_hex)
);
create trigger fx_conversion_receipt_immutable before update or delete on accounting.fx_conversion_receipt for each row execute function accounting.prevent_posted_history_mutation();
commit;
````

### FILE: `docs/fx-conversion-runtime.md`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:docs/fx-conversion-runtime.md:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "846a5650402b06f28b925cd7ed89db99f5f7f92faf8596f9b9f1112a38482165"
variables: []
secrets_allowed: false
```
````markdown
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
````

### FILE: `docs/provenance/BC_FX_DERIVATION.md`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:docs/provenance/BC_FX_DERIVATION.md:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "e49d10dc59298c20743bcc6efbee51b955c2e87016e353ed98700e33d3c7a01d"
variables: []
secrets_allowed: false
```
````markdown
# Business Central FX derivation

Algorithm source: Microsoft BCApps commit `2eae56d704a1fd035d104f333602aea7091b7749`, `src/Layers/W1/BaseApp/Finance/Currency/CurrencyExchangeRate.Table.al`, SHA256 `7f25d2996d4f2098cdfff16dfc8cb35be6101ed5052e12fbfebffab826759ac6`. Exact Git blob, byte length and SHA256 were independently checked before this connected implementation. `BC_FX_SOURCE_LOCK.json` records that file, Currency.Table.al, the enum and three inspected upstream test codeunits.

| Upstream method | Narrow adaptation | Local boundary |
|---|---|---|
| ExchangeAmtFCYToFCY, lines371–430 | ExchangeExact divides/multiplies exact Exchange Rate Amount and Relational Exch. Rate Amount through the explicit local base. | Only empty Relational Currency Code is accepted. Relational chains, event subscribers, adjustment factors and BC runtime are not implemented. |
| FindCurrency2, lines459–469 | FindLast chooses the last effective Starting Date at or before the explicit request date. | No WorkDate fallback. Duplicate identities/dates and unsupported rows fail closed. |
| ExchangeAmount, lines329–346 | Calls explicit destination precision and preserves same-currency/zero bypass. | Currency minor-unit decimals and rounding quantum are explicit configuration. Same-currency and zero inputs preserve the upstream early return. |

`internal/bcfx/exchange.go` and `selection.go` are ADAPTED, retain Microsoft copyright/MIT and declare Go representation, bounds and error handling. `rounding.go` is separately AUTHORED representation/rounding glue, not an adapted AL built-in implementation. The surrounding snapshot, hash/identity/expiry checks, API, persistence and host are AUTHORED integration glue. The upstream test codeunits are inspected source, not an executed AL suite; in particular CurrencyUT constrains currency codes, ERMCurrencyFactor verifies division by the currency factor and explicit precision, and ERMChangeExchangeRate exercises editable factor behavior that this immutable profile does not expose.

The initial raw GitHub test-source request failed with HTTP503. The same fixed Git blobs were subsequently obtained through the official GitHub blob API and verified; the original failure remains in the admission history. No branch or revision changed.

The profile accepts a supplied immutable source document, with source ID/revision and rate identities, whose exact bytes are SHA256-bound by the separately activated profile. This verifies configured bytes and identities; it does not certify the economic truth or legal acceptability of a supplied rate. Synthetic fixture rates are explicitly marked synthetic. A future materialized project selects its authorized source and dates without requiring live credentials to exercise this infrastructure.

The accounting conversion receipt is historical conversion evidence and creates no journal, posting, settlement, tax calculation or movement of money. All such effects remain with their existing owners. No signature, financial recommendation, live feed or complete Business Central runtime equivalence is claimed.

License: `licenses/Microsoft-BCApps-MIT.txt` is the exact upstream MIT notice. No new external dependency is added by this adapter. The obsolete binary64 in-memory FX pack remains outside this profile and is not relabeled.

Source correction FX-DOC-ROUND-001: the official System.Round page inspected for this candidate (https://learn.microsoft.com/en-us/dynamics365/business-central/dev-itpro/developer/methods-auto/system/system-round-method, last-updated2025-01-28) has an inconsistent negative example: -1234.56789 with precision1 is printed as -1234, while its nearest rule and other negative examples imply -1235. The root candidate comment asserting AL built-in equality was corrected before publication. This profile explicitly selects nearest/ties-away-from-zero and verifies it with an independent Python Decimal oracle; it makes no AL runtime-equivalence claim for that contradiction. No arbitrary exclusion of negative amounts is introduced.
````

### FILE: `docs/provenance/BC_FX_SOURCE_LOCK.json`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:docs/provenance/BC_FX_SOURCE_LOCK.json:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "d7acac8cc299f75da6f7470fa8d44b4785d1d75b20aa31b6acf0fc9f77414927"
variables: []
secrets_allowed: false
```
````json
{
  "algorithm_revision": "2eae56d704a1fd035d104f333602aea7091b7749",
  "source_files": [
    {
      "path": "src/Layers/W1/BaseApp/Finance/Currency/Currency.Table.al",
      "local": "Currency.Table.al",
      "source_url": "https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Finance/Currency/Currency.Table.al",
      "revision": "2eae56d704a1fd035d104f333602aea7091b7749",
      "git_blob": "d293d2b9924c32aa2b95400c375b09e250037124",
      "sha256": "4aa140ceb4a08c7923640e50ccedaec4e15b26b4a5b920d160109b270cf4182a",
      "bytes": 66950
    },
    {
      "path": "src/Layers/W1/BaseApp/Finance/Currency/CurrencyExchangeRate.Table.al",
      "local": "CurrencyExchangeRate.Table.al",
      "source_url": "https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Finance/Currency/CurrencyExchangeRate.Table.al",
      "revision": "2eae56d704a1fd035d104f333602aea7091b7749",
      "git_blob": "770743ebc00e447da68cde0b273bf1f1e6def848",
      "sha256": "7f25d2996d4f2098cdfff16dfc8cb35be6101ed5052e12fbfebffab826759ac6",
      "bytes": 32932
    },
    {
      "path": "src/Layers/W1/BaseApp/Finance/Currency/FixExchRateAmountType.Enum.al",
      "local": "FixExchRateAmountType.Enum.al",
      "source_url": "https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Finance/Currency/FixExchRateAmountType.Enum.al",
      "revision": "2eae56d704a1fd035d104f333602aea7091b7749",
      "git_blob": "44b306684924c6b05e05edb649f1b878f1ae3505",
      "sha256": "9a0cb99bf3b4b94c07ce333b779e6158be7cd3b678c734772cdbff7694f06285",
      "bytes": 1407
    }
  ],
  "inspected_upstream_test_files": [
    {
      "name": "CurrencyUT.Codeunit.al",
      "path": "src/Layers/W1/Tests/ERM-Finance/CurrencyUT.Codeunit.al",
      "sha": "5d2e7f28e153230c0e5c2edfb37bf1cd23f3523c",
      "size": 13903,
      "download_url": "https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/Tests/ERM-Finance/CurrencyUT.Codeunit.al",
      "retrieved_via": "https://api.github.com/repos/microsoft/BCApps/git/blobs/5d2e7f28e153230c0e5c2edfb37bf1cd23f3523c",
      "sha256": "1509f34bbbcf6271a4bf48f2b29982f0966a1c36fcb8330a3b6c9890f11f844b"
    },
    {
      "name": "ERMChangeExchangeRate.Codeunit.al",
      "path": "src/Layers/W1/Tests/General Journal/ERMChangeExchangeRate.Codeunit.al",
      "sha": "47e3865f764bc0bc3466afeacbc515288e2300b5",
      "size": 7734,
      "download_url": "https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/Tests/General%20Journal/ERMChangeExchangeRate.Codeunit.al",
      "retrieved_via": "https://api.github.com/repos/microsoft/BCApps/git/blobs/47e3865f764bc0bc3466afeacbc515288e2300b5",
      "sha256": "2783f0728a07df4c27d9d1d9f0f188dd1dbc522556c13327599b93b4503ea381"
    },
    {
      "name": "ERMCurrencyFactor.Codeunit.al",
      "path": "src/Layers/W1/Tests/General Journal/ERMCurrencyFactor.Codeunit.al",
      "sha": "aaffa6228157318b351f4b47b63c5163a14832cc",
      "size": 24205,
      "download_url": "https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/Tests/General%20Journal/ERMCurrencyFactor.Codeunit.al",
      "retrieved_via": "https://api.github.com/repos/microsoft/BCApps/git/blobs/aaffa6228157318b351f4b47b63c5163a14832cc",
      "sha256": "010f2b7942ada7fa1ace1bfded260c090cfec23c5a30648074a609793548eef7"
    }
  ],
  "license_sha256": "c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383",
  "upstream_tests_executed": false,
  "local_adaptation": "DIRECT_BASE_EXACT_V1",
  "license_source": {
    "url": "https://api.github.com/repos/microsoft/BCApps/license?ref=2eae56d704a1fd035d104f333602aea7091b7749",
    "sha256": "c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383",
    "bytes": 1141,
    "git_blob": "9e841e7a26e4eb057b24511e7b92d42b257a80e5",
    "git_blob_matched": true,
    "prior_raw_fetch": "HTTP503 Backend.max_conn reached; not used"
  }
}
````

### FILE: `docs/provenance/BC_FX_THIRD_PARTY_NOTICES.md`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:docs/provenance/BC_FX_THIRD_PARTY_NOTICES.md:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "ef4d4390345af5bfe76f1b056af1dd69f39cd72c6ee82e44416273d6b66476f7"
variables: []
secrets_allowed: false
```
````markdown
# BC FX adapted source notice

Copyright (c) Microsoft Corporation.

The two arithmetic/selection files in internal/bcfx are adapted from the fixed BCApps source documented in BC_FX_DERIVATION.md and licensed under the exact MIT notice in licenses/Microsoft-BCApps-MIT.txt. Local integration files are AUTHORED and are not Microsoft code. No AL runtime or upstream suite execution is claimed.

Source correction FX-DOC-ROUND-001: the official System.Round page inspected for this candidate (https://learn.microsoft.com/en-us/dynamics365/business-central/dev-itpro/developer/methods-auto/system/system-round-method, last-updated2025-01-28) has an inconsistent negative example: -1234.56789 with precision1 is printed as -1234, while its nearest rule and other negative examples imply -1235. The root candidate comment asserting AL built-in equality was corrected before publication. This profile explicitly selects nearest/ties-away-from-zero and verifies it with an independent Python Decimal oracle; it makes no AL runtime-equivalence claim for that contradiction. No arbitrary exclusion of negative amounts is introduced.
````

### FILE: `internal/accounting/fx.go`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:internal/accounting/fx.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "60b863c71f2e39501f35540f10761075487d3bbd0e23bd0d46a06c0d14e69dfa"
variables: []
secrets_allowed: false
```
````go
package accounting

// AUTHORED conversion-receipt glue. RecordConversion creates only a receipt.
// PrepareJournal binds that receipt to the existing separate posting lifecycle.
import (
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/bcfx"
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"
	"time"
)

var ErrFXNotFound = errors.New("FX conversion receipt not found")

type FXCommand struct {
	OrganizationID string `json:"organization_id"`
	FromCurrency   string `json:"from_currency"`
	ToCurrency     string `json:"to_currency"`
	AmountMinor    *int64 `json:"amount_minor,string"`
	ConversionDate string `json:"conversion_date"`
	IdempotencyKey string `json:"-"`
}
type FXReceipt struct {
	ID          string                `json:"conversion_id"`
	RequestedBy string                `json:"requested_by"`
	RequestKey  string                `json:"request_key"`
	RecordedAt  time.Time             `json:"recorded_at"`
	Snapshot    bcfx.SnapshotIdentity `json:"snapshot"`
	Conversion  bcfx.Conversion       `json:"conversion"`
	Effect      string                `json:"effect"`
}
type FXRepository interface {
	RecordFXConversion(context.Context, string, string, string, string, FXCommand, *bcfx.Snapshot, string) (FXReceipt, bool, error)
	FXConversionResult(context.Context, string, string, string, string) (FXReceipt, error)
	PrepareFXJournal(context.Context, string, string, string, string, string, FXJournalCommand, string) (FXJournalResult, bool, error)
	FXJournalResult(context.Context, string, string, string, string) (FXJournalResult, error)
}
type FXService struct {
	repository FXRepository
	ids        IDGenerator
	snapshot   *bcfx.Snapshot
}

var fxKey = regexp.MustCompile(`^[A-Za-z0-9_-]{16,128}$`)

func NewFXService(repository FXRepository, ids IDGenerator, snapshot *bcfx.Snapshot) (*FXService, error) {
	if repository == nil || ids == nil || snapshot == nil || snapshot.Identity().ProfileSHA256 == "" {
		return nil, ErrInvalid
	}
	return &FXService{repository, ids, snapshot}, nil
}
func (s *FXService) Record(ctx context.Context, tenant, actor string, c FXCommand) (FXReceipt, bool, error) {
	if s == nil || c.AmountMinor == nil || len(actor) < 1 || len(actor) > 256 || !fxKey.MatchString(c.IdempotencyKey) || !s.snapshot.Allows(tenant, c.OrganizationID) || s.snapshot.ValidateRequest(c.FromCurrency, c.ToCurrency, c.ConversionDate) != nil {
		return FXReceipt{}, false, ErrInvalid
	}
	amount := *c.AmountMinor
	c.AmountMinor = &amount
	identity := s.snapshot.Identity()
	payload, _ := json.Marshal(struct {
		Tenant, Actor, Profile, Source string
		Command                        FXCommand
	}{tenant, actor, identity.ProfileSHA256, identity.SourceSHA256, c})
	sum := sha256.Sum256(payload)
	return s.repository.RecordFXConversion(ctx, tenant, actor, s.ids.New(), s.ids.New(), c, s.snapshot, hex.EncodeToString(sum[:]))
}
func (s *FXService) Result(ctx context.Context, tenant, organization, actor, key string) (FXReceipt, error) {
	if s == nil || !s.snapshot.Allows(tenant, organization) || len(actor) < 1 || len(actor) > 256 || !fxKey.MatchString(key) {
		return FXReceipt{}, ErrInvalid
	}
	return s.repository.FXConversionResult(ctx, tenant, organization, actor, key)
}
````

### FILE: `internal/bcfx/exchange.go`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:internal/bcfx/exchange.go:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/microsoft/BCApps/blob/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Finance/Currency/CurrencyExchangeRate.Table.al"
license: "MIT"
sha256: "6d75adf6eee5c73bd92e46891e17974080cab6281de76a2deb577e1a3f2084bb"
variables: []
secrets_allowed: false
```
````go
// Copyright (c) Microsoft Corporation.
// SPDX-License-Identifier: MIT
// ADAPTED from BCApps 2eae56d704a1fd035d104f333602aea7091b7749,
// CurrencyExchangeRate.Table.al, ExchangeAmtFCYToFCY and ExchangeAmount.
// Go representation, bounds and errors are declared local modifications.
// See docs/provenance/BC_FX_DERIVATION.md and licenses/Microsoft-BCApps-MIT.txt.
package bcfx

import (
	"errors"
	"math/big"
)

var ErrRate = errors.New("unsupported or invalid exact exchange rate")
var ErrAmount = errors.New("converted amount is outside int64 minor units")

// DirectRate is the supported BC profile: Relational Currency Code is empty,
// so Exchange/Relational describe one currency against the selected local base.
// Currency identity and dated selection belong to the snapshot boundary.
type DirectRate struct{ Exchange, Relational *big.Rat }

func validRate(r DirectRate) bool {
	return r.Exchange != nil && r.Relational != nil && r.Exchange.Sign() > 0 && r.Relational.Sign() > 0 &&
		r.Exchange.Num().BitLen() <= 128 && r.Exchange.Denom().BitLen() <= 128 &&
		r.Relational.Num().BitLen() <= 128 && r.Relational.Denom().BitLen() <= 128
}

// ExchangeExact translates the direct-base branches at source lines383–387,
// 401–405 and425–433. Nil means the explicitly configured local currency.
// Callers supply major units; no floats, quote spread, fees or tax are added.
func ExchangeExact(amount *big.Rat, from, to *DirectRate) (*big.Rat, error) {
	if amount == nil || amount.Num().BitLen() > 128 || amount.Denom().BitLen() > 128 {
		return nil, ErrAmount
	}
	if from != nil && !validRate(*from) || to != nil && !validRate(*to) {
		return nil, ErrRate
	}
	out := new(big.Rat).Set(amount)
	if from != nil {
		out.Quo(out, from.Exchange)
		out.Mul(out, from.Relational)
	}
	if to != nil {
		out.Quo(out, to.Relational)
		out.Mul(out, to.Exchange)
	}
	return out, nil
}
````

### FILE: `internal/bcfx/exchange_test.go`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:internal/bcfx/exchange_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "d2b15af5bdb26b6de8ef56f6ac6550e47f38778912da95c4a816850f8b77ed8d"
variables: []
secrets_allowed: false
```
````go
// AUTHORED verification vectors for the declared narrow adaptation.
// These do not claim execution of the Microsoft AL suite.
package bcfx

import (
	"math"
	"math/big"
	"testing"
	"time"
)

func rat(s string) *big.Rat {
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		panic(s)
	}
	return r
}
func direct(e, r string) *DirectRate { return &DirectRate{rat(e), rat(r)} }
func TestExactDirectBaseConversions(t *testing.T) {
	for _, c := range []struct {
		name, amount string
		from, to     *DirectRate
		want         string
	}{
		{"local-to-local", "123.45", nil, nil, "2469/20"},
		{"foreign-to-local", "10", direct("100", "125"), nil, "25/2"},
		{"local-to-foreign", "10", nil, direct("100", "125"), "8"},
		{"foreign-cross", "10", direct("100", "125"), direct("200", "300"), "25/3"},
		{"negative", "-10", direct("100", "125"), direct("200", "300"), "-25/3"},
		{"zero", "0", direct("1", "3"), direct("1", "7"), "0"},
	} {
		t.Run(c.name, func(t *testing.T) {
			a := rat(c.amount)
			before := a.RatString()
			got, err := ExchangeExact(a, c.from, c.to)
			if err != nil || got.Cmp(rat(c.want)) != 0 {
				t.Fatalf("got %v/%v", got, err)
			}
			if a.RatString() != before {
				t.Fatal("mutated input")
			}
		})
	}
	for _, r := range []*DirectRate{direct("0", "1"), direct("1", "0"), direct("-1", "1"), {nil, rat("1")}} {
		if _, err := ExchangeExact(rat("1"), r, nil); err == nil {
			t.Fatal("invalid rate admitted")
		}
	}
}
func TestDestinationRoundingAndBounds(t *testing.T) {
	for _, c := range []struct {
		amount  string
		dec     uint8
		p, want int64
	}{
		{"1.234", 2, 1, 123}, {"1.235", 2, 1, 124}, {"-1.235", 2, 1, -124},
		{"0.025", 2, 5, 5}, {"-0.025", 2, 5, -5}, {"1.024", 2, 5, 100},
		{"9223372036854775807", 0, 1, math.MaxInt64}, {"-9223372036854775808", 0, 1, math.MinInt64},
	} {
		got, err := RoundMinor(rat(c.amount), c.dec, c.p)
		if err != nil || got != c.want {
			t.Fatalf("%s got %d %v", c.amount, got, err)
		}
	}
	for _, s := range []string{"9223372036854775807.5", "-9223372036854775808.5"} {
		if _, err := RoundMinor(rat(s), 0, 1); err == nil {
			t.Fatal("overflow")
		}
	}
	if _, err := RoundMinor(rat("1"), 10, 1); err == nil {
		t.Fatal("scale")
	}
	if _, err := RoundMinor(rat("1"), 2, 0); err == nil {
		t.Fatal("precision")
	}
}
func TestFindLastEffectiveDate(t *testing.T) {
	day := func(s string) time.Time {
		v, e := time.Parse("2006-01-02", s)
		if e != nil {
			t.Fatal(e)
		}
		return v
	}
	a := DatedRate{"USD", day("2026-01-01"), *direct("1", "10")}
	b := DatedRate{"USD", day("2026-02-01"), *direct("1", "12")}
	rows := []DatedRate{b, a, {"EUR", day("2026-01-20"), *direct("1", "20")}}
	for _, c := range []struct{ date, want string }{{"2026-01-01", "10"}, {"2026-01-31", "10"}, {"2026-02-01", "12"}} {
		got, e := FindLast(rows, "USD", day(c.date))
		if e != nil || got.Amounts.Relational.Cmp(rat(c.want)) != 0 {
			t.Fatalf("%v %v", got, e)
		}
		got.Amounts.Relational.SetInt64(99)
	}
	for _, c := range []struct {
		rows     []DatedRate
		currency string
		date     time.Time
	}{{rows, "USD", day("2025-12-31")}, {rows, "GBP", day("2026-02-01")}, {append(rows, a), "USD", day("2026-02-01")}, {rows, "USD", day("2026-02-01").Add(time.Hour)}} {
		if _, e := FindLast(c.rows, c.currency, c.date); e == nil {
			t.Fatal("ambiguous/missing/invalid date")
		}
	}
}
func FuzzDirectConversionExactInverse(f *testing.F) {
	f.Add(int64(123456), uint32(100), uint32(127))
	f.Add(int64(-1), uint32(1), uint32(3))
	f.Add(int64(math.MinInt64), uint32(1), uint32(1))
	f.Fuzz(func(t *testing.T, amount int64, exchange, relational uint32) {
		if exchange == 0 || relational == 0 {
			return
		}
		r := &DirectRate{new(big.Rat).SetInt64(int64(exchange)), new(big.Rat).SetInt64(int64(relational))}
		a := new(big.Rat).SetInt64(amount)
		converted, e := ExchangeExact(a, r, nil)
		if e != nil {
			t.Fatal(e)
		}
		back, e := ExchangeExact(converted, nil, r)
		if e != nil || a.Cmp(back) != 0 {
			t.Fatalf("inverse %v %v", back, e)
		}
		n, e := RoundMinor(back, 0, 1)
		if e != nil || n != amount {
			t.Fatalf("roundtrip %d %v", n, e)
		}
	})
}
````

### FILE: `internal/bcfx/rounding.go`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:internal/bcfx/rounding.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "f604e241f1931b61a1e616dddc45163696d2bb7fb255cefd237d3b95d7625d4e"
variables: []
secrets_allowed: false
```
````go
package bcfx

// AUTHORED exact minor-unit representation/rounding glue. The explicitly selected
// profile is nearest, ties away from zero. This is not an AL built-in equivalence
// claim: the inspected Learn table has a contradictory negative example.
import "math/big"

func RoundMinor(major *big.Rat, decimals uint8, precisionMinor int64) (int64, error) {
	if major == nil || decimals > 9 || precisionMinor <= 0 || major.Num().BitLen() > 640 || major.Denom().BitLen() > 640 {
		return 0, ErrAmount
	}
	factor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	units := new(big.Rat).Mul(major, new(big.Rat).SetInt(factor))
	units.Quo(units, new(big.Rat).SetInt64(precisionMinor))
	n := new(big.Int).Abs(units.Num())
	q, rem := new(big.Int), new(big.Int)
	q.QuoRem(n, units.Denom(), rem)
	if rem.Lsh(rem, 1).Cmp(units.Denom()) >= 0 {
		q.Add(q, big.NewInt(1))
	}
	if units.Sign() < 0 {
		q.Neg(q)
	}
	q.Mul(q, big.NewInt(precisionMinor))
	if !q.IsInt64() {
		return 0, ErrAmount
	}
	return q.Int64(), nil
}
````

### FILE: `internal/bcfx/selection.go`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:internal/bcfx/selection.go:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/microsoft/BCApps/blob/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Finance/Currency/CurrencyExchangeRate.Table.al"
license: "MIT"
sha256: "0db470b03cfc8249e7f40d98c2b7a3d273119396aed9c054baec7ae4b553afac"
variables: []
secrets_allowed: false
```
````go
// Copyright (c) Microsoft Corporation.
// SPDX-License-Identifier: MIT
// ADAPTED from BCApps 2eae56d704a1fd035d104f333602aea7091b7749,
// CurrencyExchangeRate.Table.al FindCurrency2 lines459–470.
// Slice representation, duplicate rejection and explicit date are local deltas.
package bcfx

import (
	"math/big"
	"time"
)

type DatedRate struct {
	Currency     string
	StartingDate time.Time
	Amounts      DirectRate
}

// FindLast translates SetRange(currency), SetRange(0D,date), FindLast().
// An explicit date is required: this boundary never substitutes WorkDate/Now.
// Duplicate source identities fail instead of depending on input slice order.
func FindLast(rows []DatedRate, currency string, date time.Time) (DatedRate, error) {
	var result DatedRate
	if currency == "" || !dayOnly(date) || len(rows) > 4096 {
		return result, ErrRate
	}
	seen := make(map[int64]struct{})
	for _, row := range rows {
		if row.Currency != currency {
			continue
		}
		if !dayOnly(row.StartingDate) || !validRate(row.Amounts) {
			return DatedRate{}, ErrRate
		}
		key := row.StartingDate.Unix()
		if _, ok := seen[key]; ok {
			return DatedRate{}, ErrRate
		}
		seen[key] = struct{}{}
		if !row.StartingDate.After(date) && (result.StartingDate.IsZero() || row.StartingDate.After(result.StartingDate)) {
			result = row
		}
	}
	if result.StartingDate.IsZero() {
		return DatedRate{}, ErrRate
	}
	result.Amounts = DirectRate{new(big.Rat).Set(result.Amounts.Exchange), new(big.Rat).Set(result.Amounts.Relational)}
	return result, nil
}

func dayOnly(t time.Time) bool {
	return !t.IsZero() && t.Equal(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC))
}
````

### FILE: `internal/bcfx/snapshot.go`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:internal/bcfx/snapshot.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "523b377edd668f76628efd695a18afa65cd1176c2f69d315f82414030caca22a"
variables: []
secrets_allowed: false
```
````go
package bcfx

// AUTHORED immutable input/identity boundary around the declared BC adaptation.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"time"
)

var ErrSnapshot = errors.New("invalid, expired or unsupported FX snapshot")

const ProfileSchema = "elite-fx-direct-base-profile/v1"
const SourceSchema = "elite-fx-direct-base-rates/v1"
const Algorithm = "BC_DIRECT_BASE_V1"

type CurrencyPrecision struct {
	Code                   string `json:"code"`
	MinorUnitDecimals      *uint8 `json:"minor_unit_decimals"`
	RoundingPrecisionMinor string `json:"rounding_precision_minor"`
}
type ProfileDocument struct {
	Schema             string              `json:"schema"`
	ID                 string              `json:"profile_id"`
	Revision           int                 `json:"revision"`
	Algorithm          string              `json:"algorithm"`
	Rounding           string              `json:"rounding"`
	TenantID           string              `json:"tenant_id"`
	OrganizationID     string              `json:"organization_id"`
	LocalCurrency      string              `json:"local_currency"`
	ValidFrom          string              `json:"valid_from"`
	ValidUntil         string              `json:"valid_until"`
	DateFrom           string              `json:"conversion_date_from"`
	DateThrough        string              `json:"conversion_date_through"`
	MaximumRateAgeDays int                 `json:"maximum_rate_age_days"`
	SourceID           string              `json:"source_id"`
	SourceRevision     string              `json:"source_revision"`
	SourceSHA256       string              `json:"source_sha256"`
	Currencies         []CurrencyPrecision `json:"currencies"`
}
type RateRow struct {
	ID                 string  `json:"rate_id"`
	Currency           string  `json:"currency"`
	StartingDate       string  `json:"starting_date"`
	Exchange           string  `json:"exchange_rate_amount"`
	Relational         string  `json:"relational_exchange_rate_amount"`
	RelationalCurrency *string `json:"relational_currency"`
}
type SourceDocument struct {
	Schema   string    `json:"schema"`
	ID       string    `json:"source_id"`
	Revision string    `json:"source_revision"`
	Kind     string    `json:"kind"`
	Rates    []RateRow `json:"rates"`
}
type Activation struct {
	Enabled        bool
	ProfileID      string
	Revision       int
	ProfileSHA256  string
	TenantID       string
	OrganizationID string
	LocalCurrency  string
}
type precision struct {
	decimals uint8
	quantum  int64
}
type Snapshot struct {
	profile                                      ProfileDocument
	profileRaw, sourceRaw                        []byte
	hash                                         string
	currencies                                   map[string]precision
	rows                                         []DatedRate
	identities                                   map[string]RateRow
	validFrom, validUntil, dateFrom, dateThrough time.Time
}
type SnapshotIdentity struct {
	ProfileID      string    `json:"profile_id"`
	Revision       int       `json:"profile_revision"`
	ProfileSHA256  string    `json:"profile_sha256"`
	SourceID       string    `json:"source_id"`
	SourceRevision string    `json:"source_revision"`
	SourceSHA256   string    `json:"source_sha256"`
	TenantID       string    `json:"tenant_id"`
	OrganizationID string    `json:"organization_id"`
	LocalCurrency  string    `json:"local_currency"`
	ValidFrom      time.Time `json:"valid_from"`
	ValidUntil     time.Time `json:"valid_until"`
}
type Conversion struct {
	FromCurrency           string `json:"from_currency"`
	ToCurrency             string `json:"to_currency"`
	InputMinor             int64  `json:"input_minor,string"`
	OutputMinor            int64  `json:"output_minor,string"`
	ConversionDate         string `json:"conversion_date"`
	FromDecimals           uint8  `json:"from_decimals"`
	ToDecimals             uint8  `json:"to_decimals"`
	RoundingPrecisionMinor int64  `json:"rounding_precision_minor,string"`
	RoundingApplied        bool   `json:"rounding_applied"`
	ExactMajorNumerator    string `json:"exact_major_numerator"`
	ExactMajorDenominator  string `json:"exact_major_denominator"`
	FromRateID             string `json:"from_rate_id,omitempty"`
	ToRateID               string `json:"to_rate_id,omitempty"`
	FromRateDate           string `json:"from_rate_date,omitempty"`
	ToRateDate             string `json:"to_rate_date,omitempty"`
}

var snapshotID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
var snapshotTenant = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
var currencyCode = regexp.MustCompile(`^[A-Z]{3}$`)
var decimalRate = regexp.MustCompile(`^(0|[1-9][0-9]{0,17})(\.[0-9]{1,18})?$`)

func digest(raw []byte) string { sum := sha256.Sum256(raw); return hex.EncodeToString(sum[:]) }
func hexHash(value string) bool {
	if len(value) != 64 {
		return false
	}
	raw, e := hex.DecodeString(value)
	return e == nil && hex.EncodeToString(raw) == value
}
func exactDate(raw string) (time.Time, error) {
	v, e := time.Parse("2006-01-02", raw)
	if e != nil || len(raw) != 10 || v.IsZero() || v.Format("2006-01-02") != raw {
		return time.Time{}, ErrSnapshot
	}
	return v, nil
}
func exactInstant(raw string) (time.Time, error) {
	v, e := time.Parse(time.RFC3339, raw)
	if e != nil || v.IsZero() || v.UTC().Format(time.RFC3339) != raw {
		return time.Time{}, ErrSnapshot
	}
	return v, nil
}
func decodeSnapshot(raw []byte, value any) error {
	if len(raw) == 0 || len(raw) > 1048576 || !uniqueSnapshotJSON(raw) {
		return ErrSnapshot
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(value) != nil {
		return ErrSnapshot
	}
	if d.Decode(new(any)) != io.EOF {
		return ErrSnapshot
	}
	return nil
}

// DecodeExactJSON enforces the same unambiguous FX transport representation.
func DecodeExactJSON(raw []byte, value any) error { return decodeSnapshot(raw, value) }
func LoadSnapshot(profileRaw, sourceRaw []byte, a Activation) (*Snapshot, error) {
	if !a.Enabled || !hexHash(a.ProfileSHA256) || digest(profileRaw) != a.ProfileSHA256 {
		return nil, ErrSnapshot
	}
	var p ProfileDocument
	var source SourceDocument
	if decodeSnapshot(profileRaw, &p) != nil || decodeSnapshot(sourceRaw, &source) != nil {
		return nil, ErrSnapshot
	}
	if p.Schema != ProfileSchema || p.Algorithm != Algorithm || p.Rounding != "NEAREST_TIES_AWAY_FROM_ZERO" || !snapshotID.MatchString(p.ID) || p.Revision < 1 || p.ID != a.ProfileID || p.Revision != a.Revision || !snapshotTenant.MatchString(p.TenantID) || p.TenantID != a.TenantID || !snapshotID.MatchString(p.OrganizationID) || p.OrganizationID != a.OrganizationID || !currencyCode.MatchString(p.LocalCurrency) || p.LocalCurrency != a.LocalCurrency {
		return nil, ErrSnapshot
	}
	if source.Schema != SourceSchema || source.Kind != "SUPPLIED_SNAPSHOT" || !snapshotID.MatchString(source.ID) || !snapshotID.MatchString(source.Revision) || source.ID != p.SourceID || source.Revision != p.SourceRevision || !hexHash(p.SourceSHA256) || digest(sourceRaw) != p.SourceSHA256 || len(source.Rates) > 4096 || len(p.Currencies) < 1 || len(p.Currencies) > 64 || p.MaximumRateAgeDays < 1 || p.MaximumRateAgeDays > 3660 {
		return nil, ErrSnapshot
	}
	from, e1 := exactInstant(p.ValidFrom)
	until, e2 := exactInstant(p.ValidUntil)
	dateFrom, e3 := exactDate(p.DateFrom)
	dateThrough, e4 := exactDate(p.DateThrough)
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil || !until.After(from) || dateThrough.Before(dateFrom) {
		return nil, ErrSnapshot
	}
	s := &Snapshot{profile: p, hash: a.ProfileSHA256, profileRaw: bytes.Clone(profileRaw), sourceRaw: bytes.Clone(sourceRaw), currencies: map[string]precision{}, identities: map[string]RateRow{}, validFrom: from, validUntil: until, dateFrom: dateFrom, dateThrough: dateThrough}
	for _, v := range p.Currencies {
		q, e := strconv.ParseInt(v.RoundingPrecisionMinor, 10, 64)
		if !currencyCode.MatchString(v.Code) || v.MinorUnitDecimals == nil || *v.MinorUnitDecimals > 9 || e != nil || q <= 0 || strconv.FormatInt(q, 10) != v.RoundingPrecisionMinor {
			return nil, ErrSnapshot
		}
		if _, ok := s.currencies[v.Code]; ok {
			return nil, ErrSnapshot
		}
		s.currencies[v.Code] = precision{*v.MinorUnitDecimals, q}
	}
	if _, ok := s.currencies[p.LocalCurrency]; !ok {
		return nil, ErrSnapshot
	}
	seen := map[string]bool{}
	for _, r := range source.Rates {
		start, e := exactDate(r.StartingDate)
		if !snapshotID.MatchString(r.ID) || seen[r.ID] || r.Currency == p.LocalCurrency || r.RelationalCurrency == nil || *r.RelationalCurrency != "" || e != nil || !decimalRate.MatchString(r.Exchange) || !decimalRate.MatchString(r.Relational) {
			return nil, ErrSnapshot
		}
		if _, ok := s.currencies[r.Currency]; !ok {
			return nil, ErrSnapshot
		}
		key := r.Currency + "@" + r.StartingDate
		if _, ok := s.identities[key]; ok {
			return nil, ErrSnapshot
		}
		exchange, ok1 := new(big.Rat).SetString(r.Exchange)
		relational, ok2 := new(big.Rat).SetString(r.Relational)
		rate := DirectRate{exchange, relational}
		if !ok1 || !ok2 || !validRate(rate) {
			return nil, ErrSnapshot
		}
		seen[r.ID] = true
		s.identities[key] = r
		s.rows = append(s.rows, DatedRate{r.Currency, start, rate})
	}
	return s, nil
}
func LoadSnapshotFiles(profilePath, sourcePath string, a Activation) (*Snapshot, error) {
	read := func(path string) ([]byte, error) {
		f, e := os.Open(path)
		if e != nil {
			return nil, ErrSnapshot
		}
		defer f.Close()
		raw, e := io.ReadAll(io.LimitReader(f, 1048577))
		if e != nil || len(raw) > 1048576 {
			return nil, ErrSnapshot
		}
		return raw, nil
	}
	p, e := read(profilePath)
	if e != nil {
		return nil, e
	}
	r, e := read(sourcePath)
	if e != nil {
		return nil, e
	}
	return LoadSnapshot(p, r, a)
}
func (s *Snapshot) Allows(tenant, organization string) bool {
	return s != nil && s.hash != "" && s.profile.TenantID == tenant && s.profile.OrganizationID == organization
}
func (s *Snapshot) Current(now time.Time) bool {
	return s != nil && !now.IsZero() && !now.Before(s.validFrom) && now.Before(s.validUntil)
}
func (s *Snapshot) Identity() SnapshotIdentity {
	if s == nil {
		return SnapshotIdentity{}
	}
	return SnapshotIdentity{s.profile.ID, s.profile.Revision, s.hash, s.profile.SourceID, s.profile.SourceRevision, s.profile.SourceSHA256, s.profile.TenantID, s.profile.OrganizationID, s.profile.LocalCurrency, s.validFrom, s.validUntil}
}
func (s *Snapshot) Bytes() ([]byte, []byte) {
	if s == nil {
		return nil, nil
	}
	return bytes.Clone(s.profileRaw), bytes.Clone(s.sourceRaw)
}
func (s *Snapshot) ValidateRequest(from, to, date string) error {
	if s == nil {
		return ErrSnapshot
	}
	d, e := exactDate(date)
	_, f := s.currencies[from]
	_, t := s.currencies[to]
	if e != nil || !f || !t || d.Before(s.dateFrom) || d.After(s.dateThrough) {
		return ErrSnapshot
	}
	return nil
}
func (s *Snapshot) Convert(from, to, date string, minor int64, now time.Time) (Conversion, error) {
	var v Conversion
	if s.ValidateRequest(from, to, date) != nil || !s.Current(now) {
		return v, ErrSnapshot
	}
	f, t := s.currencies[from], s.currencies[to]
	day, _ := exactDate(date)
	factor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(f.decimals)), nil)
	major := new(big.Rat).SetFrac(big.NewInt(minor), factor)
	v = Conversion{FromCurrency: from, ToCurrency: to, InputMinor: minor, ConversionDate: date, FromDecimals: f.decimals, ToDecimals: t.decimals, RoundingPrecisionMinor: t.quantum}
	if from == to || minor == 0 {
		v.OutputMinor = minor
		v.ExactMajorNumerator = major.Num().String()
		v.ExactMajorDenominator = major.Denom().String()
		return v, nil
	}
	get := func(currency string) (*DirectRate, string, string, error) {
		if currency == s.profile.LocalCurrency {
			return nil, "", "", nil
		}
		row, e := FindLast(s.rows, currency, day)
		if e != nil || row.StartingDate.Before(day.AddDate(0, 0, -s.profile.MaximumRateAgeDays)) {
			return nil, "", "", ErrSnapshot
		}
		rate := s.identities[currency+"@"+row.StartingDate.Format("2006-01-02")]
		return &row.Amounts, rate.ID, rate.StartingDate, nil
	}
	fr, fid, fd, e := get(from)
	if e != nil {
		return Conversion{}, e
	}
	tr, tid, td, e := get(to)
	if e != nil {
		return Conversion{}, e
	}
	out, e := ExchangeExact(major, fr, tr)
	if e != nil {
		return Conversion{}, e
	}
	rounded, e := RoundMinor(out, t.decimals, t.quantum)
	if e != nil {
		return Conversion{}, e
	}
	v.OutputMinor = rounded
	v.RoundingApplied = true
	v.ExactMajorNumerator = out.Num().String()
	v.ExactMajorDenominator = out.Denom().String()
	v.FromRateID = fid
	v.ToRateID = tid
	v.FromRateDate = fd
	v.ToRateDate = td
	return v, nil
}

// Reuses the existing profile duplicate-key rejection; adds a depth bound.
func uniqueSnapshotJSON(raw []byte) bool {
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	var visit func(int) bool
	visit = func(depth int) bool {
		if depth > 16 {
			return false
		}
		token, e := d.Token()
		if e != nil {
			return false
		}
		delim, ok := token.(json.Delim)
		if !ok {
			return true
		}
		switch delim {
		case '{':
			seen := map[string]bool{}
			for d.More() {
				k, e := d.Token()
				key, ok := k.(string)
				if e != nil || !ok || seen[key] {
					return false
				}
				seen[key] = true
				if !visit(depth + 1) {
					return false
				}
			}
			end, e := d.Token()
			return e == nil && end == json.Delim('}')
		case '[':
			for d.More() {
				if !visit(depth + 1) {
					return false
				}
			}
			end, e := d.Token()
			return e == nil && end == json.Delim(']')
		default:
			return false
		}
	}
	if !visit(0) {
		return false
	}
	_, e := d.Token()
	return e == io.EOF
}
````

### FILE: `internal/bcfx/snapshot_test.go`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:internal/bcfx/snapshot_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "233f4215a87270ff71c1a21c9cea901c06992cafa314d592ca2ff4ef076da260"
variables: []
secrets_allowed: false
```
````go
package bcfx

import (
	"bytes"
	"encoding/json"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func fxFixture(t testing.TB) ([]byte, []byte, Activation) {
	t.Helper()
	p, e := os.ReadFile("../../config/fx/reference-profile.json")
	if e != nil {
		t.Fatal(e)
	}
	r, e := os.ReadFile("../../config/fx/reference-rates.json")
	if e != nil {
		t.Fatal(e)
	}
	var d ProfileDocument
	if json.Unmarshal(p, &d) != nil {
		t.Fatal("profile")
	}
	return p, r, Activation{true, d.ID, d.Revision, digest(p), d.TenantID, d.OrganizationID, d.LocalCurrency}
}
func TestSnapshotExactBindingsAndConversion(t *testing.T) {
	p, r, a := fxFixture(t)
	s, e := LoadSnapshot(p, r, a)
	if e != nil {
		t.Fatal(e)
	}
	now := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	for _, v := range []struct {
		from, to, date string
		amount, want   int64
		rate           string
		rounded        bool
	}{{"EUR", "USD", "2026-01-01", 1000, 1250, "eur-2026-01-01", true}, {"EUR", "USD", "2026-01-15", 1000, 1300, "eur-2026-01-02", true}, {"USD", "EUR", "2026-01-01", 100, 80, "", true}, {"EUR", "GBP", "2026-01-01", 1000, 8333, "eur-2026-01-01", true}, {"EUR", "EUR", "2026-01-01", 103, 103, "", false}, {"EUR", "GBP", "2026-01-01", 0, 0, "", false}, {"EUR", "USD", "2026-01-01", -1000, -1250, "eur-2026-01-01", true}} {
		got, e := s.Convert(v.from, v.to, v.date, v.amount, now)
		if e != nil || got.OutputMinor != v.want || got.FromRateID != v.rate || got.RoundingApplied != v.rounded {
			t.Fatal(v, got, e)
		}
	}
	p[0] = 'x'
	r[0] = 'x'
	raw, _ := s.Bytes()
	raw[0] = 'x'
	if !s.Allows(a.TenantID, a.OrganizationID) || s.Identity().ProfileSHA256 != a.ProfileSHA256 {
		t.Fatal("mutable")
	}
	if _, e = s.Convert("EUR", "USD", "2026-01-01", 1000, now); e != nil {
		t.Fatal(e)
	}
	for _, v := range []struct {
		from, to, date string
		now            time.Time
	}{{"ZZZ", "USD", "2026-01-01", now}, {"EUR", "ZZZ", "2026-01-01", now}, {"EUR", "USD", "2025-12-31", now}, {"EUR", "USD", "2027-01-01", now}, {"EUR", "USD", "2026-01-01", time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)}} {
		if _, e := s.Convert(v.from, v.to, v.date, 0, v.now); e == nil {
			t.Fatal("invalid admitted", v)
		}
	}
}
func TestSnapshotRejectsAmbiguousAndUnsupported(t *testing.T) {
	for _, mode := range []string{"hash", "scope", "duplicate-json", "unknown", "precision", "decimals-missing", "duplicate-currency", "relational", "duplicate-rate", "duplicate-id", "source-hash", "source-unknown", "bad-date", "zero-rate", "exponent", "missing-local", "bad-rounding", "stale-rate"} {
		t.Run(mode, func(t *testing.T) {
			p, r, a := fxFixture(t)
			var d ProfileDocument
			var src SourceDocument
			_ = json.Unmarshal(p, &d)
			_ = json.Unmarshal(r, &src)
			switch mode {
			case "hash":
				a.ProfileSHA256 = strings.Repeat("0", 64)
			case "scope":
				a.OrganizationID = "foreign"
			case "precision":
				d.Currencies[0].RoundingPrecisionMinor = "0"
			case "decimals-missing":
				d.Currencies[0].MinorUnitDecimals = nil
			case "duplicate-currency":
				d.Currencies = append(d.Currencies, d.Currencies[0])
			case "relational":
				value := "USD"
				src.Rates[0].RelationalCurrency = &value
			case "duplicate-rate":
				src.Rates = append(src.Rates, src.Rates[0])
			case "duplicate-id":
				src.Rates[1].ID = src.Rates[0].ID
			case "source-hash":
				d.SourceSHA256 = strings.Repeat("0", 64)
			case "source-unknown":
				src.Rates[0].Currency = "JPY"
			case "bad-date":
				src.Rates[0].StartingDate = "2026-01-01T00:00:00Z"
			case "zero-rate":
				src.Rates[0].Exchange = "0"
			case "exponent":
				src.Rates[0].Exchange = "1e2"
			case "missing-local":
				d.Currencies = d.Currencies[1:]
			case "bad-rounding":
				d.Rounding = "DEFAULT"
			case "stale-rate":
				d.MaximumRateAgeDays = 1
			}
			r, _ = json.Marshal(src)
			if mode != "source-hash" {
				d.SourceSHA256 = digest(r)
			}
			p, _ = json.Marshal(d)
			if mode == "duplicate-json" {
				p = append([]byte(`{"schema":"other",`), p[1:]...)
			}
			if mode == "unknown" {
				p = append([]byte(`{"unknown":true,`), p[1:]...)
			}
			if mode != "hash" {
				a.ProfileSHA256 = digest(p)
			}
			s, e := LoadSnapshot(p, r, a)
			if mode == "stale-rate" {
				if e != nil {
					t.Fatal(e)
				}
				_, e = s.Convert("EUR", "USD", "2026-01-15", 100, time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC))
			}
			if e == nil {
				t.Fatal("invalid snapshot admitted")
			}
		})
	}
}
func TestRoundMinorIndependentDecimalOracle(t *testing.T) {
	raw, e := os.ReadFile("testdata/rounding-oracle.json")
	if e != nil {
		t.Fatal(e)
	}
	var doc struct {
		Vectors []struct {
			Numerator, Denominator string
			Decimals               uint8
			Precision, Want        string
			Overflow               bool
		}
	}
	if json.Unmarshal(raw, &doc) != nil || len(doc.Vectors) < 500 {
		t.Fatal("oracle")
	}
	for i, v := range doc.Vectors {
		n, _ := new(big.Int).SetString(v.Numerator, 10)
		d, _ := new(big.Int).SetString(v.Denominator, 10)
		q, _ := strconv.ParseInt(v.Precision, 10, 64)
		got, e := RoundMinor(new(big.Rat).SetFrac(n, d), v.Decimals, q)
		if v.Overflow {
			if e == nil {
				t.Fatal(i, "overflow accepted")
			}
		} else if e != nil || strconv.FormatInt(got, 10) != v.Want {
			t.Fatal(i, got, e, v.Want)
		}
	}
	t.Logf("independent Decimal oracle vectors=%d", len(doc.Vectors))
}
func FuzzSnapshotExactByteBoundary(f *testing.F) {
	p, r, a := fxFixture(f)
	f.Add(p, uint8(0))
	f.Add(append([]byte(`{"schema":"duplicate",`), p[1:]...), uint8(1))
	f.Fuzz(func(t *testing.T, raw []byte, mode uint8) {
		if len(raw) > 1048576 {
			return
		}
		activation := a
		if mode%2 == 0 {
			activation.ProfileSHA256 = digest(raw)
		}
		s, e := LoadSnapshot(raw, r, activation)
		if e != nil {
			return
		}
		if !bytes.Equal(raw, s.profileRaw) || s.Identity().ProfileSHA256 != digest(raw) || !s.Allows(a.TenantID, a.OrganizationID) {
			t.Fatal("binding invariant")
		}
		mutated, _ := s.Bytes()
		if len(mutated) > 0 {
			mutated[0] ^= 1
		}
		if !bytes.Equal(raw, s.profileRaw) {
			t.Fatal("mutable bytes")
		}
	})
}
````

### FILE: `internal/bcfx/testdata/rounding-oracle.json`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:internal/bcfx/testdata/rounding-oracle.json:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "feffa6b70d5229e27fa6adc3fa409489115e258f0177d19290a646363aaa5be9"
variables: []
secrets_allowed: false
```
````json
{
  "oracle": "Python decimal 200 digits, ROUND_HALF_UP, independent of Go math/big quotient/remainder",
  "vectors": [
    {
      "numerator": "-123456789",
      "denominator": "1",
      "decimals": 0,
      "precision": "1",
      "want": "-123456789",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "1",
      "decimals": 2,
      "precision": "1",
      "want": "-12345678900",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "1",
      "decimals": 2,
      "precision": "5",
      "want": "-12345678900",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "1",
      "decimals": 3,
      "precision": "25",
      "want": "-123456789000",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "1",
      "decimals": 9,
      "precision": "1",
      "want": "-123456789000000000",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "2",
      "decimals": 0,
      "precision": "1",
      "want": "-61728395",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "2",
      "decimals": 2,
      "precision": "1",
      "want": "-6172839450",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "2",
      "decimals": 2,
      "precision": "5",
      "want": "-6172839450",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "2",
      "decimals": 3,
      "precision": "25",
      "want": "-61728394500",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "2",
      "decimals": 9,
      "precision": "1",
      "want": "-61728394500000000",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "4",
      "decimals": 0,
      "precision": "1",
      "want": "-30864197",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "4",
      "decimals": 2,
      "precision": "1",
      "want": "-3086419725",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "4",
      "decimals": 2,
      "precision": "5",
      "want": "-3086419725",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "4",
      "decimals": 3,
      "precision": "25",
      "want": "-30864197250",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "4",
      "decimals": 9,
      "precision": "1",
      "want": "-30864197250000000",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "5",
      "decimals": 0,
      "precision": "1",
      "want": "-24691358",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "5",
      "decimals": 2,
      "precision": "1",
      "want": "-2469135780",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "5",
      "decimals": 2,
      "precision": "5",
      "want": "-2469135780",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "5",
      "decimals": 3,
      "precision": "25",
      "want": "-24691357800",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "5",
      "decimals": 9,
      "precision": "1",
      "want": "-24691357800000000",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "8",
      "decimals": 0,
      "precision": "1",
      "want": "-15432099",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "8",
      "decimals": 2,
      "precision": "1",
      "want": "-1543209863",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "8",
      "decimals": 2,
      "precision": "5",
      "want": "-1543209865",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "8",
      "decimals": 3,
      "precision": "25",
      "want": "-15432098625",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "8",
      "decimals": 9,
      "precision": "1",
      "want": "-15432098625000000",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "10",
      "decimals": 0,
      "precision": "1",
      "want": "-12345679",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "10",
      "decimals": 2,
      "precision": "1",
      "want": "-1234567890",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "10",
      "decimals": 2,
      "precision": "5",
      "want": "-1234567890",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "10",
      "decimals": 3,
      "precision": "25",
      "want": "-12345678900",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "10",
      "decimals": 9,
      "precision": "1",
      "want": "-12345678900000000",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "20",
      "decimals": 0,
      "precision": "1",
      "want": "-6172839",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "20",
      "decimals": 2,
      "precision": "1",
      "want": "-617283945",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "20",
      "decimals": 2,
      "precision": "5",
      "want": "-617283945",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "20",
      "decimals": 3,
      "precision": "25",
      "want": "-6172839450",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "20",
      "decimals": 9,
      "precision": "1",
      "want": "-6172839450000000",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "100",
      "decimals": 0,
      "precision": "1",
      "want": "-1234568",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "100",
      "decimals": 2,
      "precision": "1",
      "want": "-123456789",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "100",
      "decimals": 2,
      "precision": "5",
      "want": "-123456790",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "100",
      "decimals": 3,
      "precision": "25",
      "want": "-1234567900",
      "overflow": false
    },
    {
      "numerator": "-123456789",
      "denominator": "100",
      "decimals": 9,
      "precision": "1",
      "want": "-1234567890000000",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "1",
      "decimals": 0,
      "precision": "1",
      "want": "-103",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "1",
      "decimals": 2,
      "precision": "1",
      "want": "-10300",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "1",
      "decimals": 2,
      "precision": "5",
      "want": "-10300",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "1",
      "decimals": 3,
      "precision": "25",
      "want": "-103000",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "1",
      "decimals": 9,
      "precision": "1",
      "want": "-103000000000",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "2",
      "decimals": 0,
      "precision": "1",
      "want": "-52",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "2",
      "decimals": 2,
      "precision": "1",
      "want": "-5150",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "2",
      "decimals": 2,
      "precision": "5",
      "want": "-5150",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "2",
      "decimals": 3,
      "precision": "25",
      "want": "-51500",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "2",
      "decimals": 9,
      "precision": "1",
      "want": "-51500000000",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "4",
      "decimals": 0,
      "precision": "1",
      "want": "-26",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "4",
      "decimals": 2,
      "precision": "1",
      "want": "-2575",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "4",
      "decimals": 2,
      "precision": "5",
      "want": "-2575",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "4",
      "decimals": 3,
      "precision": "25",
      "want": "-25750",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "4",
      "decimals": 9,
      "precision": "1",
      "want": "-25750000000",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "5",
      "decimals": 0,
      "precision": "1",
      "want": "-21",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "5",
      "decimals": 2,
      "precision": "1",
      "want": "-2060",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "5",
      "decimals": 2,
      "precision": "5",
      "want": "-2060",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "5",
      "decimals": 3,
      "precision": "25",
      "want": "-20600",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "5",
      "decimals": 9,
      "precision": "1",
      "want": "-20600000000",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "8",
      "decimals": 0,
      "precision": "1",
      "want": "-13",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "8",
      "decimals": 2,
      "precision": "1",
      "want": "-1288",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "8",
      "decimals": 2,
      "precision": "5",
      "want": "-1290",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "8",
      "decimals": 3,
      "precision": "25",
      "want": "-12875",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "8",
      "decimals": 9,
      "precision": "1",
      "want": "-12875000000",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "10",
      "decimals": 0,
      "precision": "1",
      "want": "-10",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "10",
      "decimals": 2,
      "precision": "1",
      "want": "-1030",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "10",
      "decimals": 2,
      "precision": "5",
      "want": "-1030",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "10",
      "decimals": 3,
      "precision": "25",
      "want": "-10300",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "10",
      "decimals": 9,
      "precision": "1",
      "want": "-10300000000",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "20",
      "decimals": 0,
      "precision": "1",
      "want": "-5",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "20",
      "decimals": 2,
      "precision": "1",
      "want": "-515",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "20",
      "decimals": 2,
      "precision": "5",
      "want": "-515",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "20",
      "decimals": 3,
      "precision": "25",
      "want": "-5150",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "20",
      "decimals": 9,
      "precision": "1",
      "want": "-5150000000",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "100",
      "decimals": 0,
      "precision": "1",
      "want": "-1",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "100",
      "decimals": 2,
      "precision": "1",
      "want": "-103",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "100",
      "decimals": 2,
      "precision": "5",
      "want": "-105",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "100",
      "decimals": 3,
      "precision": "25",
      "want": "-1025",
      "overflow": false
    },
    {
      "numerator": "-103",
      "denominator": "100",
      "decimals": 9,
      "precision": "1",
      "want": "-1030000000",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "1",
      "decimals": 0,
      "precision": "1",
      "want": "-101",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "1",
      "decimals": 2,
      "precision": "1",
      "want": "-10100",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "1",
      "decimals": 2,
      "precision": "5",
      "want": "-10100",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "1",
      "decimals": 3,
      "precision": "25",
      "want": "-101000",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "1",
      "decimals": 9,
      "precision": "1",
      "want": "-101000000000",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "2",
      "decimals": 0,
      "precision": "1",
      "want": "-51",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "2",
      "decimals": 2,
      "precision": "1",
      "want": "-5050",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "2",
      "decimals": 2,
      "precision": "5",
      "want": "-5050",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "2",
      "decimals": 3,
      "precision": "25",
      "want": "-50500",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "2",
      "decimals": 9,
      "precision": "1",
      "want": "-50500000000",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "4",
      "decimals": 0,
      "precision": "1",
      "want": "-25",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "4",
      "decimals": 2,
      "precision": "1",
      "want": "-2525",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "4",
      "decimals": 2,
      "precision": "5",
      "want": "-2525",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "4",
      "decimals": 3,
      "precision": "25",
      "want": "-25250",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "4",
      "decimals": 9,
      "precision": "1",
      "want": "-25250000000",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "5",
      "decimals": 0,
      "precision": "1",
      "want": "-20",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "5",
      "decimals": 2,
      "precision": "1",
      "want": "-2020",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "5",
      "decimals": 2,
      "precision": "5",
      "want": "-2020",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "5",
      "decimals": 3,
      "precision": "25",
      "want": "-20200",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "5",
      "decimals": 9,
      "precision": "1",
      "want": "-20200000000",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "8",
      "decimals": 0,
      "precision": "1",
      "want": "-13",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "8",
      "decimals": 2,
      "precision": "1",
      "want": "-1263",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "8",
      "decimals": 2,
      "precision": "5",
      "want": "-1265",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "8",
      "decimals": 3,
      "precision": "25",
      "want": "-12625",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "8",
      "decimals": 9,
      "precision": "1",
      "want": "-12625000000",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "10",
      "decimals": 0,
      "precision": "1",
      "want": "-10",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "10",
      "decimals": 2,
      "precision": "1",
      "want": "-1010",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "10",
      "decimals": 2,
      "precision": "5",
      "want": "-1010",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "10",
      "decimals": 3,
      "precision": "25",
      "want": "-10100",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "10",
      "decimals": 9,
      "precision": "1",
      "want": "-10100000000",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "20",
      "decimals": 0,
      "precision": "1",
      "want": "-5",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "20",
      "decimals": 2,
      "precision": "1",
      "want": "-505",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "20",
      "decimals": 2,
      "precision": "5",
      "want": "-505",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "20",
      "decimals": 3,
      "precision": "25",
      "want": "-5050",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "20",
      "decimals": 9,
      "precision": "1",
      "want": "-5050000000",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "100",
      "decimals": 0,
      "precision": "1",
      "want": "-1",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "100",
      "decimals": 2,
      "precision": "1",
      "want": "-101",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "100",
      "decimals": 2,
      "precision": "5",
      "want": "-100",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "100",
      "decimals": 3,
      "precision": "25",
      "want": "-1000",
      "overflow": false
    },
    {
      "numerator": "-101",
      "denominator": "100",
      "decimals": 9,
      "precision": "1",
      "want": "-1010000000",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "1",
      "decimals": 0,
      "precision": "1",
      "want": "-25",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "1",
      "decimals": 2,
      "precision": "1",
      "want": "-2500",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "1",
      "decimals": 2,
      "precision": "5",
      "want": "-2500",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "1",
      "decimals": 3,
      "precision": "25",
      "want": "-25000",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "1",
      "decimals": 9,
      "precision": "1",
      "want": "-25000000000",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "2",
      "decimals": 0,
      "precision": "1",
      "want": "-13",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "2",
      "decimals": 2,
      "precision": "1",
      "want": "-1250",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "2",
      "decimals": 2,
      "precision": "5",
      "want": "-1250",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "2",
      "decimals": 3,
      "precision": "25",
      "want": "-12500",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "2",
      "decimals": 9,
      "precision": "1",
      "want": "-12500000000",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "4",
      "decimals": 0,
      "precision": "1",
      "want": "-6",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "4",
      "decimals": 2,
      "precision": "1",
      "want": "-625",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "4",
      "decimals": 2,
      "precision": "5",
      "want": "-625",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "4",
      "decimals": 3,
      "precision": "25",
      "want": "-6250",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "4",
      "decimals": 9,
      "precision": "1",
      "want": "-6250000000",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "5",
      "decimals": 0,
      "precision": "1",
      "want": "-5",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "5",
      "decimals": 2,
      "precision": "1",
      "want": "-500",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "5",
      "decimals": 2,
      "precision": "5",
      "want": "-500",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "5",
      "decimals": 3,
      "precision": "25",
      "want": "-5000",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "5",
      "decimals": 9,
      "precision": "1",
      "want": "-5000000000",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "8",
      "decimals": 0,
      "precision": "1",
      "want": "-3",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "8",
      "decimals": 2,
      "precision": "1",
      "want": "-313",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "8",
      "decimals": 2,
      "precision": "5",
      "want": "-315",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "8",
      "decimals": 3,
      "precision": "25",
      "want": "-3125",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "8",
      "decimals": 9,
      "precision": "1",
      "want": "-3125000000",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "10",
      "decimals": 0,
      "precision": "1",
      "want": "-3",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "10",
      "decimals": 2,
      "precision": "1",
      "want": "-250",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "10",
      "decimals": 2,
      "precision": "5",
      "want": "-250",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "10",
      "decimals": 3,
      "precision": "25",
      "want": "-2500",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "10",
      "decimals": 9,
      "precision": "1",
      "want": "-2500000000",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "20",
      "decimals": 0,
      "precision": "1",
      "want": "-1",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "20",
      "decimals": 2,
      "precision": "1",
      "want": "-125",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "20",
      "decimals": 2,
      "precision": "5",
      "want": "-125",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "20",
      "decimals": 3,
      "precision": "25",
      "want": "-1250",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "20",
      "decimals": 9,
      "precision": "1",
      "want": "-1250000000",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "100",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "100",
      "decimals": 2,
      "precision": "1",
      "want": "-25",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "100",
      "decimals": 2,
      "precision": "5",
      "want": "-25",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "100",
      "decimals": 3,
      "precision": "25",
      "want": "-250",
      "overflow": false
    },
    {
      "numerator": "-25",
      "denominator": "100",
      "decimals": 9,
      "precision": "1",
      "want": "-250000000",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "1",
      "decimals": 0,
      "precision": "1",
      "want": "-5",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "1",
      "decimals": 2,
      "precision": "1",
      "want": "-500",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "1",
      "decimals": 2,
      "precision": "5",
      "want": "-500",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "1",
      "decimals": 3,
      "precision": "25",
      "want": "-5000",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "1",
      "decimals": 9,
      "precision": "1",
      "want": "-5000000000",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "2",
      "decimals": 0,
      "precision": "1",
      "want": "-3",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "2",
      "decimals": 2,
      "precision": "1",
      "want": "-250",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "2",
      "decimals": 2,
      "precision": "5",
      "want": "-250",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "2",
      "decimals": 3,
      "precision": "25",
      "want": "-2500",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "2",
      "decimals": 9,
      "precision": "1",
      "want": "-2500000000",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "4",
      "decimals": 0,
      "precision": "1",
      "want": "-1",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "4",
      "decimals": 2,
      "precision": "1",
      "want": "-125",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "4",
      "decimals": 2,
      "precision": "5",
      "want": "-125",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "4",
      "decimals": 3,
      "precision": "25",
      "want": "-1250",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "4",
      "decimals": 9,
      "precision": "1",
      "want": "-1250000000",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "5",
      "decimals": 0,
      "precision": "1",
      "want": "-1",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "5",
      "decimals": 2,
      "precision": "1",
      "want": "-100",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "5",
      "decimals": 2,
      "precision": "5",
      "want": "-100",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "5",
      "decimals": 3,
      "precision": "25",
      "want": "-1000",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "5",
      "decimals": 9,
      "precision": "1",
      "want": "-1000000000",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "8",
      "decimals": 0,
      "precision": "1",
      "want": "-1",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "8",
      "decimals": 2,
      "precision": "1",
      "want": "-63",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "8",
      "decimals": 2,
      "precision": "5",
      "want": "-65",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "8",
      "decimals": 3,
      "precision": "25",
      "want": "-625",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "8",
      "decimals": 9,
      "precision": "1",
      "want": "-625000000",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "10",
      "decimals": 0,
      "precision": "1",
      "want": "-1",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "10",
      "decimals": 2,
      "precision": "1",
      "want": "-50",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "10",
      "decimals": 2,
      "precision": "5",
      "want": "-50",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "10",
      "decimals": 3,
      "precision": "25",
      "want": "-500",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "10",
      "decimals": 9,
      "precision": "1",
      "want": "-500000000",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "20",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "20",
      "decimals": 2,
      "precision": "1",
      "want": "-25",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "20",
      "decimals": 2,
      "precision": "5",
      "want": "-25",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "20",
      "decimals": 3,
      "precision": "25",
      "want": "-250",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "20",
      "decimals": 9,
      "precision": "1",
      "want": "-250000000",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "100",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "100",
      "decimals": 2,
      "precision": "1",
      "want": "-5",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "100",
      "decimals": 2,
      "precision": "5",
      "want": "-5",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "100",
      "decimals": 3,
      "precision": "25",
      "want": "-50",
      "overflow": false
    },
    {
      "numerator": "-5",
      "denominator": "100",
      "decimals": 9,
      "precision": "1",
      "want": "-50000000",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "1",
      "decimals": 0,
      "precision": "1",
      "want": "-1",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "1",
      "decimals": 2,
      "precision": "1",
      "want": "-100",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "1",
      "decimals": 2,
      "precision": "5",
      "want": "-100",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "1",
      "decimals": 3,
      "precision": "25",
      "want": "-1000",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "1",
      "decimals": 9,
      "precision": "1",
      "want": "-1000000000",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "2",
      "decimals": 0,
      "precision": "1",
      "want": "-1",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "2",
      "decimals": 2,
      "precision": "1",
      "want": "-50",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "2",
      "decimals": 2,
      "precision": "5",
      "want": "-50",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "2",
      "decimals": 3,
      "precision": "25",
      "want": "-500",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "2",
      "decimals": 9,
      "precision": "1",
      "want": "-500000000",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "4",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "4",
      "decimals": 2,
      "precision": "1",
      "want": "-25",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "4",
      "decimals": 2,
      "precision": "5",
      "want": "-25",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "4",
      "decimals": 3,
      "precision": "25",
      "want": "-250",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "4",
      "decimals": 9,
      "precision": "1",
      "want": "-250000000",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "5",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "5",
      "decimals": 2,
      "precision": "1",
      "want": "-20",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "5",
      "decimals": 2,
      "precision": "5",
      "want": "-20",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "5",
      "decimals": 3,
      "precision": "25",
      "want": "-200",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "5",
      "decimals": 9,
      "precision": "1",
      "want": "-200000000",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "8",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "8",
      "decimals": 2,
      "precision": "1",
      "want": "-13",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "8",
      "decimals": 2,
      "precision": "5",
      "want": "-15",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "8",
      "decimals": 3,
      "precision": "25",
      "want": "-125",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "8",
      "decimals": 9,
      "precision": "1",
      "want": "-125000000",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "10",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "10",
      "decimals": 2,
      "precision": "1",
      "want": "-10",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "10",
      "decimals": 2,
      "precision": "5",
      "want": "-10",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "10",
      "decimals": 3,
      "precision": "25",
      "want": "-100",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "10",
      "decimals": 9,
      "precision": "1",
      "want": "-100000000",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "20",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "20",
      "decimals": 2,
      "precision": "1",
      "want": "-5",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "20",
      "decimals": 2,
      "precision": "5",
      "want": "-5",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "20",
      "decimals": 3,
      "precision": "25",
      "want": "-50",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "20",
      "decimals": 9,
      "precision": "1",
      "want": "-50000000",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "100",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "100",
      "decimals": 2,
      "precision": "1",
      "want": "-1",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "100",
      "decimals": 2,
      "precision": "5",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "100",
      "decimals": 3,
      "precision": "25",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "-1",
      "denominator": "100",
      "decimals": 9,
      "precision": "1",
      "want": "-10000000",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "1",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "1",
      "decimals": 2,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "1",
      "decimals": 2,
      "precision": "5",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "1",
      "decimals": 3,
      "precision": "25",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "1",
      "decimals": 9,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "2",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "2",
      "decimals": 2,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "2",
      "decimals": 2,
      "precision": "5",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "2",
      "decimals": 3,
      "precision": "25",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "2",
      "decimals": 9,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "4",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "4",
      "decimals": 2,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "4",
      "decimals": 2,
      "precision": "5",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "4",
      "decimals": 3,
      "precision": "25",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "4",
      "decimals": 9,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "5",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "5",
      "decimals": 2,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "5",
      "decimals": 2,
      "precision": "5",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "5",
      "decimals": 3,
      "precision": "25",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "5",
      "decimals": 9,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "8",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "8",
      "decimals": 2,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "8",
      "decimals": 2,
      "precision": "5",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "8",
      "decimals": 3,
      "precision": "25",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "8",
      "decimals": 9,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "10",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "10",
      "decimals": 2,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "10",
      "decimals": 2,
      "precision": "5",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "10",
      "decimals": 3,
      "precision": "25",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "10",
      "decimals": 9,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "20",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "20",
      "decimals": 2,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "20",
      "decimals": 2,
      "precision": "5",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "20",
      "decimals": 3,
      "precision": "25",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "20",
      "decimals": 9,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "100",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "100",
      "decimals": 2,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "100",
      "decimals": 2,
      "precision": "5",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "100",
      "decimals": 3,
      "precision": "25",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "0",
      "denominator": "100",
      "decimals": 9,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "1",
      "decimals": 0,
      "precision": "1",
      "want": "1",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "1",
      "decimals": 2,
      "precision": "1",
      "want": "100",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "1",
      "decimals": 2,
      "precision": "5",
      "want": "100",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "1",
      "decimals": 3,
      "precision": "25",
      "want": "1000",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "1",
      "decimals": 9,
      "precision": "1",
      "want": "1000000000",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "2",
      "decimals": 0,
      "precision": "1",
      "want": "1",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "2",
      "decimals": 2,
      "precision": "1",
      "want": "50",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "2",
      "decimals": 2,
      "precision": "5",
      "want": "50",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "2",
      "decimals": 3,
      "precision": "25",
      "want": "500",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "2",
      "decimals": 9,
      "precision": "1",
      "want": "500000000",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "4",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "4",
      "decimals": 2,
      "precision": "1",
      "want": "25",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "4",
      "decimals": 2,
      "precision": "5",
      "want": "25",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "4",
      "decimals": 3,
      "precision": "25",
      "want": "250",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "4",
      "decimals": 9,
      "precision": "1",
      "want": "250000000",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "5",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "5",
      "decimals": 2,
      "precision": "1",
      "want": "20",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "5",
      "decimals": 2,
      "precision": "5",
      "want": "20",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "5",
      "decimals": 3,
      "precision": "25",
      "want": "200",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "5",
      "decimals": 9,
      "precision": "1",
      "want": "200000000",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "8",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "8",
      "decimals": 2,
      "precision": "1",
      "want": "13",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "8",
      "decimals": 2,
      "precision": "5",
      "want": "15",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "8",
      "decimals": 3,
      "precision": "25",
      "want": "125",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "8",
      "decimals": 9,
      "precision": "1",
      "want": "125000000",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "10",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "10",
      "decimals": 2,
      "precision": "1",
      "want": "10",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "10",
      "decimals": 2,
      "precision": "5",
      "want": "10",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "10",
      "decimals": 3,
      "precision": "25",
      "want": "100",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "10",
      "decimals": 9,
      "precision": "1",
      "want": "100000000",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "20",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "20",
      "decimals": 2,
      "precision": "1",
      "want": "5",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "20",
      "decimals": 2,
      "precision": "5",
      "want": "5",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "20",
      "decimals": 3,
      "precision": "25",
      "want": "50",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "20",
      "decimals": 9,
      "precision": "1",
      "want": "50000000",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "100",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "100",
      "decimals": 2,
      "precision": "1",
      "want": "1",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "100",
      "decimals": 2,
      "precision": "5",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "100",
      "decimals": 3,
      "precision": "25",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "1",
      "denominator": "100",
      "decimals": 9,
      "precision": "1",
      "want": "10000000",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "1",
      "decimals": 0,
      "precision": "1",
      "want": "5",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "1",
      "decimals": 2,
      "precision": "1",
      "want": "500",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "1",
      "decimals": 2,
      "precision": "5",
      "want": "500",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "1",
      "decimals": 3,
      "precision": "25",
      "want": "5000",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "1",
      "decimals": 9,
      "precision": "1",
      "want": "5000000000",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "2",
      "decimals": 0,
      "precision": "1",
      "want": "3",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "2",
      "decimals": 2,
      "precision": "1",
      "want": "250",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "2",
      "decimals": 2,
      "precision": "5",
      "want": "250",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "2",
      "decimals": 3,
      "precision": "25",
      "want": "2500",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "2",
      "decimals": 9,
      "precision": "1",
      "want": "2500000000",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "4",
      "decimals": 0,
      "precision": "1",
      "want": "1",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "4",
      "decimals": 2,
      "precision": "1",
      "want": "125",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "4",
      "decimals": 2,
      "precision": "5",
      "want": "125",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "4",
      "decimals": 3,
      "precision": "25",
      "want": "1250",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "4",
      "decimals": 9,
      "precision": "1",
      "want": "1250000000",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "5",
      "decimals": 0,
      "precision": "1",
      "want": "1",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "5",
      "decimals": 2,
      "precision": "1",
      "want": "100",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "5",
      "decimals": 2,
      "precision": "5",
      "want": "100",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "5",
      "decimals": 3,
      "precision": "25",
      "want": "1000",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "5",
      "decimals": 9,
      "precision": "1",
      "want": "1000000000",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "8",
      "decimals": 0,
      "precision": "1",
      "want": "1",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "8",
      "decimals": 2,
      "precision": "1",
      "want": "63",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "8",
      "decimals": 2,
      "precision": "5",
      "want": "65",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "8",
      "decimals": 3,
      "precision": "25",
      "want": "625",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "8",
      "decimals": 9,
      "precision": "1",
      "want": "625000000",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "10",
      "decimals": 0,
      "precision": "1",
      "want": "1",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "10",
      "decimals": 2,
      "precision": "1",
      "want": "50",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "10",
      "decimals": 2,
      "precision": "5",
      "want": "50",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "10",
      "decimals": 3,
      "precision": "25",
      "want": "500",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "10",
      "decimals": 9,
      "precision": "1",
      "want": "500000000",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "20",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "20",
      "decimals": 2,
      "precision": "1",
      "want": "25",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "20",
      "decimals": 2,
      "precision": "5",
      "want": "25",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "20",
      "decimals": 3,
      "precision": "25",
      "want": "250",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "20",
      "decimals": 9,
      "precision": "1",
      "want": "250000000",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "100",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "100",
      "decimals": 2,
      "precision": "1",
      "want": "5",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "100",
      "decimals": 2,
      "precision": "5",
      "want": "5",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "100",
      "decimals": 3,
      "precision": "25",
      "want": "50",
      "overflow": false
    },
    {
      "numerator": "5",
      "denominator": "100",
      "decimals": 9,
      "precision": "1",
      "want": "50000000",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "1",
      "decimals": 0,
      "precision": "1",
      "want": "25",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "1",
      "decimals": 2,
      "precision": "1",
      "want": "2500",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "1",
      "decimals": 2,
      "precision": "5",
      "want": "2500",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "1",
      "decimals": 3,
      "precision": "25",
      "want": "25000",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "1",
      "decimals": 9,
      "precision": "1",
      "want": "25000000000",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "2",
      "decimals": 0,
      "precision": "1",
      "want": "13",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "2",
      "decimals": 2,
      "precision": "1",
      "want": "1250",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "2",
      "decimals": 2,
      "precision": "5",
      "want": "1250",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "2",
      "decimals": 3,
      "precision": "25",
      "want": "12500",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "2",
      "decimals": 9,
      "precision": "1",
      "want": "12500000000",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "4",
      "decimals": 0,
      "precision": "1",
      "want": "6",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "4",
      "decimals": 2,
      "precision": "1",
      "want": "625",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "4",
      "decimals": 2,
      "precision": "5",
      "want": "625",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "4",
      "decimals": 3,
      "precision": "25",
      "want": "6250",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "4",
      "decimals": 9,
      "precision": "1",
      "want": "6250000000",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "5",
      "decimals": 0,
      "precision": "1",
      "want": "5",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "5",
      "decimals": 2,
      "precision": "1",
      "want": "500",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "5",
      "decimals": 2,
      "precision": "5",
      "want": "500",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "5",
      "decimals": 3,
      "precision": "25",
      "want": "5000",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "5",
      "decimals": 9,
      "precision": "1",
      "want": "5000000000",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "8",
      "decimals": 0,
      "precision": "1",
      "want": "3",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "8",
      "decimals": 2,
      "precision": "1",
      "want": "313",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "8",
      "decimals": 2,
      "precision": "5",
      "want": "315",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "8",
      "decimals": 3,
      "precision": "25",
      "want": "3125",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "8",
      "decimals": 9,
      "precision": "1",
      "want": "3125000000",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "10",
      "decimals": 0,
      "precision": "1",
      "want": "3",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "10",
      "decimals": 2,
      "precision": "1",
      "want": "250",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "10",
      "decimals": 2,
      "precision": "5",
      "want": "250",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "10",
      "decimals": 3,
      "precision": "25",
      "want": "2500",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "10",
      "decimals": 9,
      "precision": "1",
      "want": "2500000000",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "20",
      "decimals": 0,
      "precision": "1",
      "want": "1",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "20",
      "decimals": 2,
      "precision": "1",
      "want": "125",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "20",
      "decimals": 2,
      "precision": "5",
      "want": "125",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "20",
      "decimals": 3,
      "precision": "25",
      "want": "1250",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "20",
      "decimals": 9,
      "precision": "1",
      "want": "1250000000",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "100",
      "decimals": 0,
      "precision": "1",
      "want": "0",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "100",
      "decimals": 2,
      "precision": "1",
      "want": "25",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "100",
      "decimals": 2,
      "precision": "5",
      "want": "25",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "100",
      "decimals": 3,
      "precision": "25",
      "want": "250",
      "overflow": false
    },
    {
      "numerator": "25",
      "denominator": "100",
      "decimals": 9,
      "precision": "1",
      "want": "250000000",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "1",
      "decimals": 0,
      "precision": "1",
      "want": "101",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "1",
      "decimals": 2,
      "precision": "1",
      "want": "10100",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "1",
      "decimals": 2,
      "precision": "5",
      "want": "10100",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "1",
      "decimals": 3,
      "precision": "25",
      "want": "101000",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "1",
      "decimals": 9,
      "precision": "1",
      "want": "101000000000",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "2",
      "decimals": 0,
      "precision": "1",
      "want": "51",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "2",
      "decimals": 2,
      "precision": "1",
      "want": "5050",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "2",
      "decimals": 2,
      "precision": "5",
      "want": "5050",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "2",
      "decimals": 3,
      "precision": "25",
      "want": "50500",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "2",
      "decimals": 9,
      "precision": "1",
      "want": "50500000000",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "4",
      "decimals": 0,
      "precision": "1",
      "want": "25",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "4",
      "decimals": 2,
      "precision": "1",
      "want": "2525",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "4",
      "decimals": 2,
      "precision": "5",
      "want": "2525",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "4",
      "decimals": 3,
      "precision": "25",
      "want": "25250",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "4",
      "decimals": 9,
      "precision": "1",
      "want": "25250000000",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "5",
      "decimals": 0,
      "precision": "1",
      "want": "20",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "5",
      "decimals": 2,
      "precision": "1",
      "want": "2020",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "5",
      "decimals": 2,
      "precision": "5",
      "want": "2020",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "5",
      "decimals": 3,
      "precision": "25",
      "want": "20200",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "5",
      "decimals": 9,
      "precision": "1",
      "want": "20200000000",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "8",
      "decimals": 0,
      "precision": "1",
      "want": "13",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "8",
      "decimals": 2,
      "precision": "1",
      "want": "1263",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "8",
      "decimals": 2,
      "precision": "5",
      "want": "1265",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "8",
      "decimals": 3,
      "precision": "25",
      "want": "12625",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "8",
      "decimals": 9,
      "precision": "1",
      "want": "12625000000",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "10",
      "decimals": 0,
      "precision": "1",
      "want": "10",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "10",
      "decimals": 2,
      "precision": "1",
      "want": "1010",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "10",
      "decimals": 2,
      "precision": "5",
      "want": "1010",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "10",
      "decimals": 3,
      "precision": "25",
      "want": "10100",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "10",
      "decimals": 9,
      "precision": "1",
      "want": "10100000000",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "20",
      "decimals": 0,
      "precision": "1",
      "want": "5",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "20",
      "decimals": 2,
      "precision": "1",
      "want": "505",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "20",
      "decimals": 2,
      "precision": "5",
      "want": "505",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "20",
      "decimals": 3,
      "precision": "25",
      "want": "5050",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "20",
      "decimals": 9,
      "precision": "1",
      "want": "5050000000",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "100",
      "decimals": 0,
      "precision": "1",
      "want": "1",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "100",
      "decimals": 2,
      "precision": "1",
      "want": "101",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "100",
      "decimals": 2,
      "precision": "5",
      "want": "100",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "100",
      "decimals": 3,
      "precision": "25",
      "want": "1000",
      "overflow": false
    },
    {
      "numerator": "101",
      "denominator": "100",
      "decimals": 9,
      "precision": "1",
      "want": "1010000000",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "1",
      "decimals": 0,
      "precision": "1",
      "want": "103",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "1",
      "decimals": 2,
      "precision": "1",
      "want": "10300",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "1",
      "decimals": 2,
      "precision": "5",
      "want": "10300",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "1",
      "decimals": 3,
      "precision": "25",
      "want": "103000",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "1",
      "decimals": 9,
      "precision": "1",
      "want": "103000000000",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "2",
      "decimals": 0,
      "precision": "1",
      "want": "52",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "2",
      "decimals": 2,
      "precision": "1",
      "want": "5150",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "2",
      "decimals": 2,
      "precision": "5",
      "want": "5150",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "2",
      "decimals": 3,
      "precision": "25",
      "want": "51500",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "2",
      "decimals": 9,
      "precision": "1",
      "want": "51500000000",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "4",
      "decimals": 0,
      "precision": "1",
      "want": "26",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "4",
      "decimals": 2,
      "precision": "1",
      "want": "2575",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "4",
      "decimals": 2,
      "precision": "5",
      "want": "2575",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "4",
      "decimals": 3,
      "precision": "25",
      "want": "25750",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "4",
      "decimals": 9,
      "precision": "1",
      "want": "25750000000",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "5",
      "decimals": 0,
      "precision": "1",
      "want": "21",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "5",
      "decimals": 2,
      "precision": "1",
      "want": "2060",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "5",
      "decimals": 2,
      "precision": "5",
      "want": "2060",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "5",
      "decimals": 3,
      "precision": "25",
      "want": "20600",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "5",
      "decimals": 9,
      "precision": "1",
      "want": "20600000000",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "8",
      "decimals": 0,
      "precision": "1",
      "want": "13",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "8",
      "decimals": 2,
      "precision": "1",
      "want": "1288",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "8",
      "decimals": 2,
      "precision": "5",
      "want": "1290",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "8",
      "decimals": 3,
      "precision": "25",
      "want": "12875",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "8",
      "decimals": 9,
      "precision": "1",
      "want": "12875000000",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "10",
      "decimals": 0,
      "precision": "1",
      "want": "10",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "10",
      "decimals": 2,
      "precision": "1",
      "want": "1030",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "10",
      "decimals": 2,
      "precision": "5",
      "want": "1030",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "10",
      "decimals": 3,
      "precision": "25",
      "want": "10300",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "10",
      "decimals": 9,
      "precision": "1",
      "want": "10300000000",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "20",
      "decimals": 0,
      "precision": "1",
      "want": "5",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "20",
      "decimals": 2,
      "precision": "1",
      "want": "515",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "20",
      "decimals": 2,
      "precision": "5",
      "want": "515",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "20",
      "decimals": 3,
      "precision": "25",
      "want": "5150",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "20",
      "decimals": 9,
      "precision": "1",
      "want": "5150000000",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "100",
      "decimals": 0,
      "precision": "1",
      "want": "1",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "100",
      "decimals": 2,
      "precision": "1",
      "want": "103",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "100",
      "decimals": 2,
      "precision": "5",
      "want": "105",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "100",
      "decimals": 3,
      "precision": "25",
      "want": "1025",
      "overflow": false
    },
    {
      "numerator": "103",
      "denominator": "100",
      "decimals": 9,
      "precision": "1",
      "want": "1030000000",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "1",
      "decimals": 0,
      "precision": "1",
      "want": "123456789",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "1",
      "decimals": 2,
      "precision": "1",
      "want": "12345678900",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "1",
      "decimals": 2,
      "precision": "5",
      "want": "12345678900",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "1",
      "decimals": 3,
      "precision": "25",
      "want": "123456789000",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "1",
      "decimals": 9,
      "precision": "1",
      "want": "123456789000000000",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "2",
      "decimals": 0,
      "precision": "1",
      "want": "61728395",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "2",
      "decimals": 2,
      "precision": "1",
      "want": "6172839450",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "2",
      "decimals": 2,
      "precision": "5",
      "want": "6172839450",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "2",
      "decimals": 3,
      "precision": "25",
      "want": "61728394500",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "2",
      "decimals": 9,
      "precision": "1",
      "want": "61728394500000000",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "4",
      "decimals": 0,
      "precision": "1",
      "want": "30864197",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "4",
      "decimals": 2,
      "precision": "1",
      "want": "3086419725",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "4",
      "decimals": 2,
      "precision": "5",
      "want": "3086419725",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "4",
      "decimals": 3,
      "precision": "25",
      "want": "30864197250",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "4",
      "decimals": 9,
      "precision": "1",
      "want": "30864197250000000",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "5",
      "decimals": 0,
      "precision": "1",
      "want": "24691358",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "5",
      "decimals": 2,
      "precision": "1",
      "want": "2469135780",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "5",
      "decimals": 2,
      "precision": "5",
      "want": "2469135780",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "5",
      "decimals": 3,
      "precision": "25",
      "want": "24691357800",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "5",
      "decimals": 9,
      "precision": "1",
      "want": "24691357800000000",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "8",
      "decimals": 0,
      "precision": "1",
      "want": "15432099",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "8",
      "decimals": 2,
      "precision": "1",
      "want": "1543209863",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "8",
      "decimals": 2,
      "precision": "5",
      "want": "1543209865",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "8",
      "decimals": 3,
      "precision": "25",
      "want": "15432098625",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "8",
      "decimals": 9,
      "precision": "1",
      "want": "15432098625000000",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "10",
      "decimals": 0,
      "precision": "1",
      "want": "12345679",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "10",
      "decimals": 2,
      "precision": "1",
      "want": "1234567890",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "10",
      "decimals": 2,
      "precision": "5",
      "want": "1234567890",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "10",
      "decimals": 3,
      "precision": "25",
      "want": "12345678900",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "10",
      "decimals": 9,
      "precision": "1",
      "want": "12345678900000000",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "20",
      "decimals": 0,
      "precision": "1",
      "want": "6172839",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "20",
      "decimals": 2,
      "precision": "1",
      "want": "617283945",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "20",
      "decimals": 2,
      "precision": "5",
      "want": "617283945",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "20",
      "decimals": 3,
      "precision": "25",
      "want": "6172839450",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "20",
      "decimals": 9,
      "precision": "1",
      "want": "6172839450000000",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "100",
      "decimals": 0,
      "precision": "1",
      "want": "1234568",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "100",
      "decimals": 2,
      "precision": "1",
      "want": "123456789",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "100",
      "decimals": 2,
      "precision": "5",
      "want": "123456790",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "100",
      "decimals": 3,
      "precision": "25",
      "want": "1234567900",
      "overflow": false
    },
    {
      "numerator": "123456789",
      "denominator": "100",
      "decimals": 9,
      "precision": "1",
      "want": "1234567890000000",
      "overflow": false
    },
    {
      "numerator": "9223372036854775807",
      "denominator": "1",
      "decimals": 0,
      "precision": "1",
      "want": "9223372036854775807",
      "overflow": false
    },
    {
      "numerator": "-9223372036854775808",
      "denominator": "1",
      "decimals": 0,
      "precision": "1",
      "want": "-9223372036854775808",
      "overflow": false
    },
    {
      "numerator": "18446744073709551615",
      "denominator": "2",
      "decimals": 0,
      "precision": "1",
      "want": "9223372036854775808",
      "overflow": true
    },
    {
      "numerator": "-18446744073709551617",
      "denominator": "2",
      "decimals": 0,
      "precision": "1",
      "want": "-9223372036854775809",
      "overflow": true
    },
    {
      "numerator": "-123456789",
      "denominator": "100000",
      "decimals": 0,
      "precision": "1",
      "want": "-1235",
      "overflow": false
    }
  ]
}
````

### FILE: `internal/platform/httpapi/fx_conversion.go`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:internal/platform/httpapi/fx_conversion.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "9da7fb9aeea6bdc906b70cdcec478dd852e96b8db755824f9cbc4f39b0367afc"
variables: []
secrets_allowed: false
```
````go
package httpapi

// AUTHORED authenticated conversion and draft-binding routes.
// The existing accounting endpoints retain posting and reversal authority.
import (
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcfx"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"io"
	"net/http"
)

type FXConversionModule struct{ Service *accounting.FXService }

func (m FXConversionModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	a := fxAPI{m.Service, verifier}
	mux.HandleFunc("POST /v1/accounting/fx/conversions", a.record)
	mux.HandleFunc("GET /v1/accounting/fx/conversions/result", a.result)
	mux.HandleFunc("POST /v1/accounting/fx/journals", a.prepareJournal)
	mux.HandleFunc("GET /v1/accounting/fx/journals/result", a.journalResult)
}

type fxAPI struct {
	service  *accounting.FXService
	verifier identity.Verifier
}

func (a fxAPI) principal(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
	p, e := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if e != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "valid bearer required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", "permission required")
		return p, false
	}
	return p, true
}
func (a fxAPI) record(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "accounting:write")
	if !ok {
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "JSON required")
		return
	}
	if len(r.Header.Values("Idempotency-Key")) != 1 {
		writeProblem(w, 400, "INVALID_KEY", "one request key required")
		return
	}
	var c accounting.FXCommand
	raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 4096))
	if e != nil {
		var limit *http.MaxBytesError
		if errors.As(e, &limit) {
			writeProblem(w, 413, "BODY_TOO_LARGE", "request is too large")
		} else {
			writeProblem(w, 400, "INVALID_REQUEST", "request is invalid")
		}
		return
	}
	if bcfx.DecodeExactJSON(raw, &c) != nil {
		writeProblem(w, 400, "INVALID_REQUEST", "request is invalid")
		return
	}
	c.IdempotencyKey = r.Header.Get("Idempotency-Key")
	if !accountingOrg(w, p, c.OrganizationID) {
		return
	}
	v, replay, e := a.service.Record(r.Context(), p.TenantID, p.Subject, c)
	if e != nil {
		fxProblem(w, e)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	if replay {
		w.Header().Set("Idempotent-Replay", "true")
	}
	writeJSON(w, 201, v)
}
func (a fxAPI) result(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "accounting:read")
	if !ok {
		return
	}
	q := r.URL.Query()
	if len(q) != 1 || len(q["organization_id"]) != 1 || len(r.Header.Values("Idempotency-Key")) != 1 {
		writeProblem(w, 400, "INVALID_QUERY", "exact organization and request key required")
		return
	}
	organization := q.Get("organization_id")
	if !accountingOrg(w, p, organization) {
		return
	}
	v, e := a.service.Result(r.Context(), p.TenantID, organization, p.Subject, r.Header.Get("Idempotency-Key"))
	if e != nil {
		fxProblem(w, e)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, 200, v)
}
func fxProblem(w http.ResponseWriter, e error) {
	if errors.Is(e, accounting.ErrFXNotFound) {
		writeProblem(w, 404, "NOT_FOUND", "conversion receipt not found")
		return
	}
	writeAccountingResult(w, e)
}
````

### FILE: `internal/platform/postgres/fx_connected_integration_test.go`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:internal/platform/postgres/fx_connected_integration_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "931f2f38bc03da3316eda57a5ffc3feb49e12a8c126123304b7c673cb0fa723a"
variables: []
secrets_allowed: false
```
````go
package postgres_test

// AUTHORED local fixture: actual HTTP/RS256/JWKS, adapted arithmetic and PG.
import (
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcfx"
	"elite.local/enterprise/internal/platform/httpapi"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/hex"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func fxConnectedSnapshot(t testing.TB, tenant string, alter func(*bcfx.ProfileDocument)) *bcfx.Snapshot {
	t.Helper()
	p, e := os.ReadFile("../../../config/fx/reference-profile.json")
	if e != nil {
		t.Fatal(e)
	}
	r, e := os.ReadFile("../../../config/fx/reference-rates.json")
	if e != nil {
		t.Fatal(e)
	}
	var d bcfx.ProfileDocument
	if json.Unmarshal(p, &d) != nil {
		t.Fatal("profile")
	}
	d.TenantID = tenant
	d.ValidFrom = time.Now().UTC().Add(-time.Hour).Truncate(time.Second).Format(time.RFC3339)
	d.ValidUntil = time.Now().UTC().Add(time.Hour).Truncate(time.Second).Format(time.RFC3339)
	if alter != nil {
		alter(&d)
	}
	p, e = json.Marshal(d)
	if e != nil {
		t.Fatal(e)
	}
	sum := sha256.Sum256(p)
	s, e := bcfx.LoadSnapshot(p, r, bcfx.Activation{Enabled: true, ProfileID: d.ID, Revision: d.Revision, ProfileSHA256: hex.EncodeToString(sum[:]), TenantID: tenant, OrganizationID: d.OrganizationID, LocalCurrency: d.LocalCurrency})
	if e != nil {
		t.Fatal(e)
	}
	return s
}
func TestFXConnectedHTTPReceipts(t *testing.T) {
	dbURL := os.Getenv("FX_CONNECTED_DB_URL")
	if dbURL == "" {
		t.Skip("explicit FX fixture DB required")
	}
	u, e := url.Parse(dbURL)
	if e != nil || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_fx_") {
		t.Fatal("isolated loopback FX DB required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, e := pgxpool.New(ctx, dbURL)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	ids := randomid.Generator{}
	tenant := ids.New()
	if _, e = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'fx-connected','Fixture','Fixture')`, tenant); e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Store','store'),($1,'other','other','Other','store')`, tenant); e != nil {
		t.Fatal(e)
	}
	snapshot := fxConnectedSnapshot(t, tenant, nil)
	repo := db.NewAccounting(pool)
	service, e := accounting.NewFXService(repo, ids, snapshot)
	if e != nil {
		t.Fatal(e)
	}
	verifier, token := handoverBrowserIssuer(t)
	mux := http.NewServeMux()
	httpapi.FXConversionModule{Service: service}.Register(mux, verifier)
	server := httptest.NewServer(mux)
	defer server.Close()
	good := token("accountant", tenant, []string{"accounting:read", "accounting:write"}, []string{"store"})
	other := token("other", tenant, []string{"accounting:read", "accounting:write"}, []string{"store"})
	readonly := token("reader", tenant, []string{"accounting:read"}, []string{"store"})
	foreign := token("foreign", tenant, []string{"accounting:read", "accounting:write"}, []string{"other"})
	request := func(method, key, bearer, body string) (int, http.Header, []byte) {
		path := "/v1/accounting/fx/conversions"
		if method == "GET" {
			path += "/result?organization_id=store"
		}
		r, e := http.NewRequestWithContext(ctx, method, server.URL+path, strings.NewReader(body))
		if e != nil {
			t.Fatal(e)
		}
		r.Header.Set("Authorization", "Bearer "+bearer)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Idempotency-Key", key)
		response, e := server.Client().Do(r)
		if e != nil {
			t.Error(e)
			return 0, nil, nil
		}
		defer response.Body.Close()
		raw, e := io.ReadAll(response.Body)
		if e != nil {
			t.Error(e)
		}
		return response.StatusCode, response.Header, raw
	}
	command := `{"organization_id":"store","from_currency":"EUR","to_currency":"GBP","amount_minor":"1000","conversion_date":"2026-01-01"}`
	key := "fx-connected-first-0001"
	status, _, first := request("POST", key, good, command)
	if status != 201 {
		t.Fatalf("POST %d %s", status, first)
	}
	// Treat the mutation response as lost; recover through the authenticated GET.
	status, headers, recovered := request("GET", key, good, "")
	if status != 200 || headers.Get("Cache-Control") != "no-store" || string(recovered) != string(first) {
		t.Fatal("recovery", status, string(recovered))
	}
	var receipt accounting.FXReceipt
	if json.Unmarshal(recovered, &receipt) != nil || receipt.Conversion.OutputMinor != 8333 || receipt.Conversion.FromRateID != "eur-2026-01-01" || receipt.Conversion.ToRateID != "gbp-2026-01-01" || receipt.Effect != "CONVERSION_RECEIPT_ONLY" {
		t.Fatal("receipt", receipt)
	}
	status, headers, replay := request("POST", key, good, command)
	if status != 201 || headers.Get("Idempotent-Replay") != "true" || string(replay) != string(first) {
		t.Fatal("replay")
	}
	for _, v := range []struct {
		method, key, actor, body string
		want                     int
	}{{"POST", key, good, strings.Replace(command, `"1000"`, `"1001"`, 1), 409}, {"POST", key, other, command, 409}, {"GET", key, other, "", 404}, {"POST", "fx-no-permission-01", readonly, command, 403}, {"POST", "fx-foreign-org-0001", foreign, command, 403}, {"POST", "fx-unknown-curr-01", good, strings.Replace(command, "EUR", "ZZZ", 1), 400}, {"POST", "fx-missing-money-01", good, strings.Replace(command, `"amount_minor":"1000",`, "", 1), 400}, {"POST", "fx-duplicate-json-01", good, strings.Replace(command, `"amount_minor":"1000"`, `"amount_minor":"1","amount_minor":"1000"`, 1), 400}} {
		status, _, raw := request(v.method, v.key, v.actor, v.body)
		if status != v.want {
			t.Fatal(v, status, string(raw))
		}
	}
	var wg sync.WaitGroup
	results := make(chan []byte, 12)
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			status, _, raw := request("POST", "fx-concurrent-key-0001", good, command)
			if status != 201 {
				t.Errorf("concurrent status=%d body=%s", status, raw)
			}
			results <- raw
		}()
	}
	wg.Wait()
	close(results)
	var parallel []byte
	for raw := range results {
		if parallel == nil {
			parallel = raw
		} else if string(parallel) != string(raw) {
			t.Fatal("concurrent receipts diverged")
		}
	}
	// Shared idempotency expiry does not permit a second durable conversion.
	if _, e = pool.Exec(ctx, `delete from platform.idempotency_record where tenant_id=$1 and scope='accounting-fx-conversion' and idempotency_key=$2`, tenant, key); e != nil {
		t.Fatal(e)
	}
	status, _, retained := request("POST", key, good, command)
	if status != 201 || string(retained) != string(first) {
		t.Fatal("retention replay")
	}
	if _, e = pool.Exec(ctx, `update accounting.fx_rate_snapshot set source_raw=source_raw where tenant_id=$1`, tenant); e == nil {
		t.Fatal("snapshot mutable")
	}
	if _, e = pool.Exec(ctx, `delete from accounting.fx_conversion_receipt where tenant_id=$1`, tenant); e == nil {
		t.Fatal("receipt mutable")
	}
	amount := int64(1000)
	c := accounting.FXCommand{OrganizationID: "store", FromCurrency: "EUR", ToCurrency: "USD", AmountMinor: &amount, ConversionDate: "2026-01-01", IdempotencyKey: "fx-policy-collision-01"}
	changed := fxConnectedSnapshot(t, tenant, func(p *bcfx.ProfileDocument) { p.Currencies[0].RoundingPrecisionMinor = "5" })
	changedService, _ := accounting.NewFXService(repo, ids, changed)
	if _, _, e = changedService.Record(ctx, tenant, "accountant", c); e == nil {
		t.Fatal("profile id/revision changed bytes")
	}
	// Expiration while waiting on an outbox write must roll back all new records.
	if _, e = pool.Exec(ctx, `create function public.fx_wait_fixture() returns trigger language plpgsql as $$ begin if NEW.event_type='fx-conversion.recorded' and NEW.payload->'snapshot'->>'profile_id'='expiry-wait' then perform pg_sleep(2.2); end if; return NEW; end; $$; create trigger fx_wait_fixture before insert on platform.outbox_event for each row execute function public.fx_wait_fixture();`); e != nil {
		t.Fatal(e)
	}
	expiry := fxConnectedSnapshot(t, tenant, func(p *bcfx.ProfileDocument) {
		p.ID = "expiry-wait"
		p.ValidUntil = time.Now().UTC().Add(2 * time.Second).Truncate(time.Second).Format(time.RFC3339)
	})
	expiryService, _ := accounting.NewFXService(repo, ids, expiry)
	c.IdempotencyKey = "fx-expiring-wait-0001"
	if _, _, e = expiryService.Record(ctx, tenant, "accountant", c); e == nil {
		t.Fatal("expired before commit")
	}
	var receipts, events, snapshots, journals, entries, claims int
	e = pool.QueryRow(ctx, `select (select count(*) from accounting.fx_conversion_receipt where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='fx-conversion.recorded'),(select count(*) from accounting.fx_rate_snapshot where tenant_id=$1),(select count(*) from accounting.journal where tenant_id=$1),(select count(*) from accounting.entry where tenant_id=$1),(select count(*) from platform.idempotency_record where tenant_id=$1 and idempotency_key in ('fx-expiring-wait-0001','fx-policy-collision-01'))`, tenant).Scan(&receipts, &events, &snapshots, &journals, &entries, &claims)
	if e != nil || receipts != 2 || events != 2 || snapshots != 1 || journals != 0 || entries != 0 || claims != 0 {
		t.Fatal(receipts, events, snapshots, journals, entries, claims, e)
	}
	t.Log("FX_CONNECTED_PASS HTTP_RS256_JWKS=true receipts=2 outbox=2 snapshots=1 concurrent12_one_receipt=true scoped_recovery=true immutable=true expires_while_outbox_wait_rolls_back=true retention_replay=true journal_entries=0")
}
````

### FILE: `internal/platform/postgres/fx_conversion.go`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:internal/platform/postgres/fx_conversion.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "2f278736d928f31004e1b4f11f511f69b875ed93a03b65d3d1f814fbef7f4fb3"
variables: []
secrets_allowed: false
```
````go
package postgres

// AUTHORED immutable receipt, shared idempotency and transactional outbox glue.
import (
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcfx"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"
)

type fxReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func fxHash(raw []byte) string { v := sha256.Sum256(raw); return hex.EncodeToString(v[:]) }
func readFXReceipt(ctx context.Context, q fxReader, tenant, organization, actor, key string) (accounting.FXReceipt, string, error) {
	var v accounting.FXReceipt
	var raw []byte
	var digest, requestHash, id string
	err := q.QueryRow(ctx, `select conversion_id,receipt_raw,receipt_sha256_hex,request_sha256_hex from accounting.fx_conversion_receipt where tenant_id=$1 and organization_id=$2 and requested_by_subject=$3 and request_key=$4`, tenant, organization, actor, key).Scan(&id, &raw, &digest, &requestHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, "", accounting.ErrFXNotFound
	}
	if err != nil {
		return v, "", err
	}
	if fxHash(raw) != digest || json.Unmarshal(raw, &v) != nil || v.ID != id || v.RequestedBy != actor || v.RequestKey != key || v.Snapshot.TenantID != tenant || v.Snapshot.OrganizationID != organization || v.Effect != "CONVERSION_RECEIPT_ONLY" {
		return accounting.FXReceipt{}, "", accounting.ErrConflict
	}
	return v, requestHash, nil
}
func (r *Accounting) FXConversionResult(ctx context.Context, tenant, organization, actor, key string) (accounting.FXReceipt, error) {
	v, _, err := readFXReceipt(ctx, r.pool, tenant, organization, actor, key)
	return v, err
}
func (r *Accounting) RecordFXConversion(ctx context.Context, tenant, actor, id, event string, c accounting.FXCommand, s *bcfx.Snapshot, hash string) (accounting.FXReceipt, bool, error) {
	var empty accounting.FXReceipt
	if s == nil || c.AmountMinor == nil || !s.Allows(tenant, c.OrganizationID) || len(hash) != 64 {
		return empty, false, accounting.ErrInvalid
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return empty, false, err
	}
	defer tx.Rollback(ctx)
	// A durable receipt still prevents duplicate conversion if shared transient
	// idempotency rows have been retained for a shorter period by the host.
	if prior, saved, e := readFXReceipt(ctx, tx, tenant, c.OrganizationID, actor, c.IdempotencyKey); e == nil {
		if saved != hash {
			return empty, false, accounting.ErrConflict
		}
		return prior, true, tx.Commit(ctx)
	} else if !errors.Is(e, accounting.ErrFXNotFound) {
		return empty, false, e
	}
	claim, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at) values($1,'accounting-fx-conversion',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict do nothing`, tenant, c.IdempotencyKey, hash)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	if claim.RowsAffected() == 0 {
		var saved string
		if e := tx.QueryRow(ctx, `select request_sha256_hex from platform.idempotency_record where tenant_id=$1 and scope='accounting-fx-conversion' and idempotency_key=$2`, tenant, c.IdempotencyKey).Scan(&saved); e != nil || saved != hash {
			return empty, false, accounting.ErrConflict
		}
		v, receiptHash, e := readFXReceipt(ctx, tx, tenant, c.OrganizationID, actor, c.IdempotencyKey)
		if e != nil || receiptHash != hash {
			return empty, false, accounting.ErrConflict
		}
		return v, true, tx.Commit(ctx)
	}
	var organization string
	if err = tx.QueryRow(ctx, `select organization_id from org.organization where tenant_id=$1 and organization_id=$2 and status='active' for share`, tenant, c.OrganizationID).Scan(&organization); err != nil {
		return empty, false, accounting.ErrConflict
	}
	identity := s.Identity()
	profileRaw, sourceRaw := s.Bytes()
	_, err = tx.Exec(ctx, `insert into accounting.fx_rate_snapshot(tenant_id,profile_id,profile_revision,organization_id,profile_sha256_hex,source_sha256_hex,profile_raw,source_raw,valid_from,valid_until) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) on conflict do nothing`, tenant, identity.ProfileID, identity.Revision, c.OrganizationID, identity.ProfileSHA256, identity.SourceSHA256, profileRaw, sourceRaw, identity.ValidFrom, identity.ValidUntil)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	var savedProfile, savedSource, savedOrg string
	if err = tx.QueryRow(ctx, `select profile_sha256_hex,source_sha256_hex,organization_id from accounting.fx_rate_snapshot where tenant_id=$1 and profile_id=$2 and profile_revision=$3`, tenant, identity.ProfileID, identity.Revision).Scan(&savedProfile, &savedSource, &savedOrg); err != nil || savedProfile != identity.ProfileSHA256 || savedSource != identity.SourceSHA256 || savedOrg != c.OrganizationID {
		return empty, false, accounting.ErrConflict
	}
	var now time.Time
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
		return empty, false, err
	}
	value, err := s.Convert(c.FromCurrency, c.ToCurrency, c.ConversionDate, *c.AmountMinor, now)
	if err != nil {
		return empty, false, accounting.ErrConflict
	}
	receipt := accounting.FXReceipt{ID: id, RequestedBy: actor, RequestKey: c.IdempotencyKey, RecordedAt: now.UTC(), Snapshot: identity, Conversion: value, Effect: "CONVERSION_RECEIPT_ONLY"}
	raw, err := json.Marshal(receipt)
	if err != nil || len(raw) > 16384 {
		return empty, false, accounting.ErrInvalid
	}
	_, err = tx.Exec(ctx, `insert into accounting.fx_conversion_receipt(tenant_id,conversion_id,organization_id,requested_by_subject,request_key,request_sha256_hex,profile_id,profile_revision,profile_sha256_hex,receipt_raw,receipt_sha256_hex,recorded_at) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`, tenant, id, c.OrganizationID, actor, c.IdempotencyKey, hash, identity.ProfileID, identity.Revision, identity.ProfileSHA256, raw, fxHash(raw), now)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'fx-conversion',$3,1,'fx-conversion.recorded',1,$4,$5)`, tenant, event, id, now, raw)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	// Expiry is rechecked after blocking inserts/outbox, before any commit.
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil || !s.Current(now) {
		return empty, false, accounting.ErrConflict
	}
	done, err := tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('conversion_id',$3::text),resource_type='fx-conversion',resource_id=$3,locked_until=null where tenant_id=$1 and scope='accounting-fx-conversion' and idempotency_key=$2 and status='processing'`, tenant, c.IdempotencyKey, id)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	if done.RowsAffected() != 1 {
		return empty, false, accounting.ErrConflict
	}
	if err = tx.Commit(ctx); err != nil {
		return empty, false, accountingConflict(err)
	}
	return receipt, false, nil
}
````

### FILE: `licenses/Microsoft-BCApps-MIT.txt`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:licenses/Microsoft-BCApps-MIT.txt:v1"
operation: CREATE
provenance: VERBATIM
source: "https://api.github.com/repos/microsoft/BCApps/license?ref=2eae56d704a1fd035d104f333602aea7091b7749"
license: "MIT"
sha256: "c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383"
variables: []
secrets_allowed: false
```
````text
    MIT License

    Copyright (c) Microsoft Corporation.

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE
````

### FILE: `tools/generate_fx_rounding_oracle.py`
```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:tools/generate_fx_rounding_oracle.py:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration glue; see docs/provenance/BC_FX_DERIVATION.md"
license: "LicenseRef-Workspace-Owner"
sha256: "d5c122dfdda81eea8b63433ae3bd7852d2ae65c7ca20cb1e9c7d81ad607aa114"
variables: []
secrets_allowed: false
```
````python
"""AUTHORED independent Decimal oracle; creates only an absent output file."""
from decimal import Decimal, localcontext, ROUND_HALF_UP
from pathlib import Path
import argparse
import json

parser = argparse.ArgumentParser()
parser.add_argument("--output", type=Path, required=True)
args = parser.parse_args()
cases = [(str(n), str(d), dec, q)
         for n in [-123456789, -103, -101, -25, -5, -1, 0, 1, 5, 25, 101, 103, 123456789]
         for d in [1, 2, 4, 5, 8, 10, 20, 100]
         for dec, q in [(0, 1), (2, 1), (2, 5), (3, 25), (9, 1)]]
cases += [(n, d, 0, 1) for n, d in [
    ("9223372036854775807", "1"), ("-9223372036854775808", "1"),
    ("18446744073709551615", "2"), ("-18446744073709551617", "2"),
    ("-123456789", "100000")]]
vectors = []
with localcontext() as context:
    context.prec = 200
    for n, d, decimals, quantum in cases:
        units = Decimal(n) / Decimal(d) * Decimal(10) ** decimals / Decimal(quantum)
        result = int(units.quantize(Decimal("1"), rounding=ROUND_HALF_UP) * Decimal(quantum))
        vectors.append({"numerator": n, "denominator": d, "decimals": decimals,
                        "precision": str(quantum), "want": str(result),
                        "overflow": not (-(2 ** 63) <= result < 2 ** 63)})
document = {"oracle": "Python decimal 200 digits, ROUND_HALF_UP, independent of Go math/big quotient/remainder",
            "vectors": vectors}
with args.output.open("x", encoding="utf-8", newline="\n") as stream:
    stream.write(json.dumps(document, indent=2) + "\n")
````

V402 composed delta: Connect the immutable FX conversion to the existing draft writer; preserve posting/reversal, explicit actor/account/period selection and no corporate authorship. Evidence FX_JOURNAL_CONNECTION_V402.md.

### FILE: `internal/accounting/fx_journal.go`

```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:internal/accounting/fx_journal.go:v1"
operation: CREATE
provenance: AUTHORED
source: "typed transaction, authorization, HTTP and verification glue around the existing source-admitted FX and accounting owners; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "dfe80352679004ad25108e6c575fca7b30839660b78db4020dd2c68bc326f114"
variables: []
secrets_allowed: false
```

````go
package accounting

// AUTHORED binding of an existing immutable conversion to existing accounting.
// The accountant selects accounts; this adapter selects no business/fiscal rule.
import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
	"unicode/utf8"
)

type FXJournalCommand struct {
	OrganizationID       string `json:"organization_id"`
	ConversionID         string `json:"conversion_id"`
	ConversionRequestKey string `json:"conversion_request_key"`
	PeriodID             string `json:"period_id"`
	DebitAccount         string `json:"debit_account"`
	CreditAccount        string `json:"credit_account"`
	Description          string `json:"description"`
	IdempotencyKey       string `json:"-"`
}
type FXJournalReceipt struct {
	JournalID      string    `json:"journal_id"`
	ConversionID   string    `json:"conversion_id"`
	OrganizationID string    `json:"organization_id"`
	RequestedBy    string    `json:"requested_by"`
	RequestKey     string    `json:"request_key"`
	PeriodID       string    `json:"period_id"`
	PostingDate    string    `json:"posting_date"`
	Currency       string    `json:"currency"`
	AmountMinor    int64     `json:"amount_minor,string"`
	DebitAccount   string    `json:"debit_account"`
	CreditAccount  string    `json:"credit_account"`
	Description    string    `json:"description"`
	RecordedAt     time.Time `json:"recorded_at"`
	Effect         string    `json:"effect"`
}
type FXJournalResult struct {
	Receipt        FXJournalReceipt `json:"receipt"`
	CurrentStatus  string           `json:"current_status"`
	CurrentVersion int64            `json:"current_version"`
}

func (s *FXService) PrepareJournal(ctx context.Context, tenant, actor string, c FXJournalCommand) (FXJournalResult, bool, error) {
	if s == nil || !s.snapshot.Allows(tenant, c.OrganizationID) || len(actor) < 1 || len(actor) > 256 || !fxKey.MatchString(c.IdempotencyKey) || !fxKey.MatchString(c.ConversionRequestKey) || len(c.ConversionID) < 1 || len(c.ConversionID) > 128 || len(c.PeriodID) < 1 || len(c.PeriodID) > 128 || !codePattern.MatchString(c.DebitAccount) || !codePattern.MatchString(c.CreditAccount) || !utf8.ValidString(c.Description) || len(c.Description) > 250 {
		return FXJournalResult{}, false, ErrInvalid
	}
	raw, err := json.Marshal(struct {
		Tenant, Actor string
		Command       FXJournalCommand
	}{tenant, actor, c})
	if err != nil {
		return FXJournalResult{}, false, ErrInvalid
	}
	hash := sha256.Sum256(raw)
	return s.repository.PrepareFXJournal(ctx, tenant, actor, s.ids.New(), s.ids.New(), s.ids.New(), c, hex.EncodeToString(hash[:]))
}
func (s *FXService) JournalResult(ctx context.Context, tenant, organization, actor, key string) (FXJournalResult, error) {
	if s == nil || !s.snapshot.Allows(tenant, organization) || len(actor) < 1 || len(actor) > 256 || !fxKey.MatchString(key) {
		return FXJournalResult{}, ErrInvalid
	}
	return s.repository.FXJournalResult(ctx, tenant, organization, actor, key)
}
````

### FILE: `internal/platform/postgres/fx_journal.go`

```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:internal/platform/postgres/fx_journal.go:v1"
operation: CREATE
provenance: AUTHORED
source: "typed transaction, authorization, HTTP and verification glue around the existing source-admitted FX and accounting owners; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "5f05998186adcc1d4a5b94f1b2b8b48c09ab1f3072fd810d9b483c6972c8dec3"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED atomic linkage to the existing draft writer. Conversion, posting and
// reversal algorithms remain with their admitted existing owners.
import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcamounts"
	"github.com/jackc/pgx/v5"
)

func readFXJournal(ctx context.Context, q fxReader, tenant, organization, actor, key string) (accounting.FXJournalResult, string, error) {
	var out accounting.FXJournalResult
	var raw []byte
	var hash, requestHash, journal, conversion, currency, period, sourceType, sourceID string
	var amount, credit int64
	var date time.Time
	err := q.QueryRow(ctx, `select r.receipt_raw,r.receipt_sha256_hex,r.request_sha256_hex,r.journal_id,r.conversion_id,j.status,j.version,j.currency,j.period_id,j.posting_date,j.total_debit_minor_units,j.total_credit_minor_units,j.source_type,j.source_id from accounting.fx_journal_receipt r join accounting.journal j on j.tenant_id=r.tenant_id and j.journal_id=r.journal_id and j.organization_id=r.organization_id where r.tenant_id=$1 and r.organization_id=$2 and r.requested_by_subject=$3 and r.request_key=$4`, tenant, organization, actor, key).Scan(&raw, &hash, &requestHash, &journal, &conversion, &out.CurrentStatus, &out.CurrentVersion, &currency, &period, &date, &amount, &credit, &sourceType, &sourceID)
	if errors.Is(err, pgx.ErrNoRows) {
		return out, "", accounting.ErrFXNotFound
	}
	if err != nil {
		return out, "", err
	}
	if fxHash(raw) != hash || json.Unmarshal(raw, &out.Receipt) != nil {
		return accounting.FXJournalResult{}, "", accounting.ErrConflict
	}
	r := out.Receipt
	if r.JournalID != journal || r.ConversionID != conversion || r.OrganizationID != organization || r.RequestedBy != actor || r.RequestKey != key || r.Effect != "FX_JOURNAL_PREPARED" || r.AmountMinor != amount || amount != credit || r.Currency != currency || r.PeriodID != period || r.PostingDate != date.Format("2006-01-02") || sourceType != "FX_CONVERSION" || sourceID != conversion {
		return accounting.FXJournalResult{}, "", accounting.ErrConflict
	}
	return out, requestHash, nil
}
func (r *Accounting) FXJournalResult(ctx context.Context, tenant, organization, actor, key string) (accounting.FXJournalResult, error) {
	out, _, err := readFXJournal(ctx, r.pool, tenant, organization, actor, key)
	return out, err
}

func (r *Accounting) PrepareFXJournal(ctx context.Context, tenant, actor, journalID, journalEvent, bindingEvent string, c accounting.FXJournalCommand, hash string) (accounting.FXJournalResult, bool, error) {
	var empty accounting.FXJournalResult
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return empty, false, err
	}
	defer tx.Rollback(ctx)
	if prior, saved, e := readFXJournal(ctx, tx, tenant, c.OrganizationID, actor, c.IdempotencyKey); e == nil {
		if saved != hash {
			return empty, false, accounting.ErrConflict
		}
		return prior, true, tx.Commit(ctx)
	} else if !errors.Is(e, accounting.ErrFXNotFound) {
		return empty, false, e
	}
	claim, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at) values($1,'accounting-fx-journal',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict do nothing`, tenant, c.IdempotencyKey, hash)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	if claim.RowsAffected() == 0 {
		var saved string
		if e := tx.QueryRow(ctx, `select request_sha256_hex from platform.idempotency_record where tenant_id=$1 and scope='accounting-fx-journal' and idempotency_key=$2`, tenant, c.IdempotencyKey).Scan(&saved); e != nil || saved != hash {
			return empty, false, accounting.ErrConflict
		}
		prior, saved, e := readFXJournal(ctx, tx, tenant, c.OrganizationID, actor, c.IdempotencyKey)
		if e != nil || saved != hash {
			return empty, false, accounting.ErrConflict
		}
		return prior, true, tx.Commit(ctx)
	}
	var conversionID string
	if err = tx.QueryRow(ctx, `select c.conversion_id from accounting.fx_conversion_receipt c join org.organization o on o.tenant_id=c.tenant_id and o.organization_id=c.organization_id where c.tenant_id=$1 and c.organization_id=$2 and c.requested_by_subject=$3 and c.request_key=$4 and c.conversion_id=$5 and o.status='active' for update of c for share of o`, tenant, c.OrganizationID, actor, c.ConversionRequestKey, c.ConversionID).Scan(&conversionID); err != nil {
		return empty, false, accounting.ErrConflict
	}
	conversion, _, err := readFXReceipt(ctx, tx, tenant, c.OrganizationID, actor, c.ConversionRequestKey)
	if err != nil {
		return empty, false, err
	}
	if conversion.ID != conversionID || conversion.Conversion.ToCurrency != conversion.Snapshot.LocalCurrency || conversion.Conversion.InputMinor <= 0 || bcamounts.CheckBalance(conversion.Conversion.OutputMinor, conversion.Conversion.OutputMinor) != nil {
		return empty, false, accounting.ErrConflict
	}
	var used bool
	if err = tx.QueryRow(ctx, `select exists(select 1 from accounting.fx_journal_receipt where tenant_id=$1 and conversion_id=$2)`, tenant, conversionID).Scan(&used); err != nil {
		return empty, false, err
	}
	if used {
		return empty, false, accounting.ErrConflict
	}
	date, err := time.Parse("2006-01-02", conversion.Conversion.ConversionDate)
	if err != nil {
		return empty, false, accounting.ErrConflict
	}
	amount := conversion.Conversion.OutputMinor
	journal := accounting.Journal{ID: journalID, OrganizationID: c.OrganizationID, PeriodID: c.PeriodID, SourceType: "FX_CONVERSION", SourceID: conversionID, Currency: conversion.Conversion.ToCurrency, PostingDate: date, Status: "draft", Version: 1, TotalDebitMinorUnits: amount, TotalCreditMinorUnits: amount, Lines: []accounting.Line{{LineNo: 1, AccountCode: c.DebitAccount, Description: c.Description, DebitMinorUnits: amount}, {LineNo: 2, AccountCode: c.CreditAccount, Description: c.Description, CreditMinorUnits: amount}}}
	if err = createJournalInTx(ctx, tx, tenant, journalEvent, journal); err != nil {
		return empty, false, err
	}
	var now time.Time
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
		return empty, false, err
	}
	receipt := accounting.FXJournalReceipt{JournalID: journalID, ConversionID: conversionID, OrganizationID: c.OrganizationID, RequestedBy: actor, RequestKey: c.IdempotencyKey, PeriodID: c.PeriodID, PostingDate: conversion.Conversion.ConversionDate, Currency: journal.Currency, AmountMinor: amount, DebitAccount: c.DebitAccount, CreditAccount: c.CreditAccount, Description: c.Description, RecordedAt: now.UTC(), Effect: "FX_JOURNAL_PREPARED"}
	raw, err := json.Marshal(receipt)
	if err != nil || len(raw) > 16384 {
		return empty, false, accounting.ErrInvalid
	}
	_, err = tx.Exec(ctx, `insert into accounting.fx_journal_receipt(tenant_id,organization_id,conversion_id,journal_id,requested_by_subject,request_key,request_sha256_hex,receipt_raw,receipt_sha256_hex,recorded_at) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, tenant, c.OrganizationID, conversionID, journalID, actor, c.IdempotencyKey, hash, raw, fxHash(raw), now)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'journal',$3,1,'fx-journal.prepared',1,$4,$5)`, tenant, bindingEvent, journalID, now, raw)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	done, err := tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('journal_id',$3::text),resource_type='fx-journal',resource_id=$3,locked_until=null where tenant_id=$1 and scope='accounting-fx-journal' and idempotency_key=$2 and status='processing'`, tenant, c.IdempotencyKey, journalID)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	if done.RowsAffected() != 1 {
		return empty, false, accounting.ErrConflict
	}
	if err = tx.Commit(ctx); err != nil {
		return empty, false, accountingConflict(err)
	}
	return accounting.FXJournalResult{Receipt: receipt, CurrentStatus: "draft", CurrentVersion: 1}, false, nil
}
````

### FILE: `internal/platform/httpapi/fx_journal.go`

```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:internal/platform/httpapi/fx_journal.go:v1"
operation: CREATE
provenance: AUTHORED
source: "typed transaction, authorization, HTTP and verification glue around the existing source-admitted FX and accounting owners; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1153dce4e0214feeb794e5deb04a546e7b748e2c8fd11cd36d4c70a33b31bcf7"
variables: []
secrets_allowed: false
```

````go
package httpapi

// AUTHORED typed HTTP binding. Existing posting/reversal endpoints keep authority.
import (
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcfx"
	"io"
	"net/http"
)

func (a fxAPI) prepareJournal(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "accounting:write")
	if !ok {
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "JSON required")
		return
	}
	if len(r.Header.Values("Idempotency-Key")) != 1 {
		writeProblem(w, 400, "INVALID_KEY", "one request key required")
		return
	}
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 4096))
	if err != nil {
		writeProblem(w, 413, "BODY_TOO_LARGE", "bounded JSON required")
		return
	}
	var c accounting.FXJournalCommand
	if bcfx.DecodeExactJSON(raw, &c) != nil {
		writeProblem(w, 400, "INVALID_REQUEST", "exact FX journal request required")
		return
	}
	c.IdempotencyKey = r.Header.Get("Idempotency-Key")
	if !accountingOrg(w, p, c.OrganizationID) {
		return
	}
	result, replay, err := a.service.PrepareJournal(r.Context(), p.TenantID, p.Subject, c)
	if err != nil {
		fxProblem(w, err)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	if replay {
		w.Header().Set("Idempotent-Replay", "true")
	}
	writeJSON(w, 201, result)
}
func (a fxAPI) journalResult(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "accounting:read")
	if !ok {
		return
	}
	q := r.URL.Query()
	if len(q) != 1 || len(q["organization_id"]) != 1 || len(r.Header.Values("Idempotency-Key")) != 1 {
		writeProblem(w, 400, "INVALID_QUERY", "exact organization and request key required")
		return
	}
	org := q.Get("organization_id")
	if !accountingOrg(w, p, org) {
		return
	}
	result, err := a.service.JournalResult(r.Context(), p.TenantID, org, p.Subject, r.Header.Get("Idempotency-Key"))
	if err != nil {
		fxProblem(w, err)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, 200, result)
}
````

### FILE: `internal/platform/postgres/fx_journal_integration_test.go`

```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:internal/platform/postgres/fx_journal_integration_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "typed transaction, authorization, HTTP and verification glue around the existing source-admitted FX and accounting owners; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "249983207f974d5486a699b74086ea96d6a857f3fc83bc3ebcd23f4d53c96525"
variables: []
secrets_allowed: false
```

````go
package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/platform/httpapi"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestFXJournalConnectedPostingRecoveryAndReversal(t *testing.T) {
	rawURL := os.Getenv("FX_CONNECTED_DB_URL")
	if rawURL == "" {
		t.Skip("explicit local FX fixture required")
	}
	u, e := url.Parse(rawURL)
	if e != nil || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_fx_") {
		t.Fatal("owned loopback FX database required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	pool, e := pgxpool.New(ctx, rawURL)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	ids := randomid.Generator{}
	tenant := ids.New()
	if _, e = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'fx-journal','Fixture','Fixture')`, tenant); e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Store','store'),($1,'other','other','Other','store')`, tenant); e != nil {
		t.Fatal(e)
	}
	repo := db.NewAccounting(pool)
	ledger := accounting.NewService(repo, ids)
	fx, e := accounting.NewFXService(repo, ids, fxConnectedSnapshot(t, tenant, nil))
	if e != nil {
		t.Fatal(e)
	}
	for _, account := range []accounting.Account{{Code: "FX_ASSET", Name: "Fixture asset", Type: "asset"}, {Code: "FX_CLEARING", Name: "Fixture clearing", Type: "liability"}} {
		if _, e = ledger.CreateAccount(ctx, tenant, account); e != nil {
			t.Fatal(e)
		}
	}
	period, e := ledger.OpenPeriod(ctx, tenant, accounting.Period{StartsOn: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), EndsOn: time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)})
	if e != nil {
		t.Fatal(e)
	}
	verifier, token := handoverBrowserIssuer(t)
	mux := http.NewServeMux()
	httpapi.FXConversionModule{Service: fx}.Register(mux, verifier)
	httpapi.AccountingModule{Service: ledger}.Register(mux, verifier)
	server := httptest.NewServer(mux)
	defer server.Close()
	good := token("accountant", tenant, []string{"accounting:write", "accounting:read"}, []string{"store"})
	other := token("other", tenant, []string{"accounting:write", "accounting:read"}, []string{"store"})
	readonly := token("reader", tenant, []string{"accounting:read"}, []string{"store"})
	foreign := token("foreign", tenant, []string{"accounting:write", "accounting:read"}, []string{"other"})
	controller := token("controller", tenant, []string{"accounting:post", "accounting:reverse", "accounting:read"}, []string{"store"})
	request := func(method, path, key, bearer, body string) (int, http.Header, []byte) {
		r, e := http.NewRequestWithContext(ctx, method, server.URL+path, strings.NewReader(body))
		if e != nil {
			t.Error(e)
			return 0, nil, nil
		}
		r.Header.Set("Authorization", "Bearer "+bearer)
		r.Header.Set("Content-Type", "application/json")
		if key != "" {
			r.Header.Set("Idempotency-Key", key)
		}
		response, e := server.Client().Do(r)
		if e != nil {
			t.Error(e)
			return 0, nil, nil
		}
		defer response.Body.Close()
		raw, e := io.ReadAll(response.Body)
		if e != nil {
			t.Error(e)
		}
		return response.StatusCode, response.Header, raw
	}
	conversion := func(key, to string, amount int64) accounting.FXReceipt {
		body := fmt.Sprintf(`{"organization_id":"store","from_currency":"EUR","to_currency":%q,"amount_minor":"%d","conversion_date":"2026-01-01"}`, to, amount)
		status, _, raw := request("POST", "/v1/accounting/fx/conversions", key, good, body)
		var out accounting.FXReceipt
		if status != 201 || json.Unmarshal(raw, &out) != nil {
			t.Fatalf("conversion %d %s", status, raw)
		}
		return out
	}
	command := func(receipt accounting.FXReceipt) accounting.FXJournalCommand {
		return accounting.FXJournalCommand{OrganizationID: "store", ConversionID: receipt.ID, ConversionRequestKey: receipt.RequestKey, PeriodID: period.ID, DebitAccount: "FX_ASSET", CreditAccount: "FX_CLEARING", Description: "Synthetic fixture conversion"}
	}
	encode := func(v any) string {
		raw, e := json.Marshal(v)
		if e != nil {
			t.Fatal(e)
		}
		return string(raw)
	}
	converted := conversion("fx-journal-conversion-01", "USD", 1000)
	if converted.Conversion.OutputMinor != 1250 {
		t.Fatal("independent fixture amount", converted)
	}
	c := command(converted)
	body := encode(c)
	key := "fx-journal-prepare-0001"
	var group sync.WaitGroup
	results := make(chan accounting.FXJournalResult, 12)
	for range 12 {
		group.Add(1)
		go func() {
			defer group.Done()
			status, _, raw := request("POST", "/v1/accounting/fx/journals", key, good, body)
			var result accounting.FXJournalResult
			if status != 201 || json.Unmarshal(raw, &result) != nil {
				t.Errorf("prepare %d %s", status, raw)
				return
			}
			results <- result
		}()
	}
	group.Wait()
	close(results)
	var initial accounting.FXJournalResult
	count := 0
	for result := range results {
		count++
		if initial.Receipt.JournalID == "" {
			initial = result
		}
		if encode(result) != encode(initial) {
			t.Fatal("concurrent draft receipts differ")
		}
	}
	if count != 12 || initial.CurrentStatus != "draft" || initial.Receipt.AmountMinor != 1250 || initial.Receipt.Currency != "USD" || initial.Receipt.PostingDate != "2026-01-01" || initial.Receipt.Effect != "FX_JOURNAL_PREPARED" {
		t.Fatal("draft receipt", initial, count)
	}
	resultPath := "/v1/accounting/fx/journals/result?organization_id=store"
	status, headers, recovered := request("GET", resultPath, key, good, "")
	if status != 200 || headers.Get("Cache-Control") != "no-store" || strings.TrimSpace(string(recovered)) != encode(initial) {
		t.Fatal("lost prepare response recovery", status, string(recovered))
	}
	var n int
	if e = pool.QueryRow(ctx, `select count(*) from accounting.journal where tenant_id=$1`, tenant).Scan(&n); e != nil || n != 1 {
		t.Fatal("journal count", n, e)
	}
	if e = pool.QueryRow(ctx, `select count(*) from accounting.entry where tenant_id=$1`, tenant).Scan(&n); e != nil || n != 0 {
		t.Fatal("draft posted by preparation", n, e)
	}
	for _, tc := range []struct {
		method, path, key, bearer, body string
		want                            int
	}{
		{"POST", "/v1/accounting/fx/journals", key, good, strings.Replace(body, "FX_ASSET", "FX_CLEARING", 1), 409},
		{"POST", "/v1/accounting/fx/journals", key, other, body, 409},
		{"GET", resultPath, key, other, "", 404},
		{"POST", "/v1/accounting/fx/journals", "fx-journal-no-scope-01", readonly, body, 403},
		{"POST", "/v1/accounting/fx/journals", "fx-journal-foreign-01", foreign, body, 403},
		{"POST", "/v1/accounting/fx/journals", "fx-journal-reuse-0001", good, body, 409},
		{"POST", "/v1/accounting/fx/journals", "fx-journal-inject-01", good, strings.Replace(body, "{", `{"amount_minor":"1",`, 1), 400},
		{"POST", "/v1/accounting/fx/journals", "fx-journal-duplicate-01", good, strings.Replace(body, "{", `{"organization_id":"store",`, 1), 400},
	} {
		status, _, raw := request(tc.method, tc.path, tc.key, tc.bearer, tc.body)
		if status != tc.want {
			t.Errorf("negative wanted%d got%d %s", tc.want, status, raw)
		}
	}
	forged := accounting.Journal{OrganizationID: "store", PeriodID: period.ID, SourceType: "FX_CONVERSION", SourceID: "forged", Currency: "USD", PostingDate: period.StartsOn, Lines: []accounting.Line{{LineNo: 1, AccountCode: "FX_ASSET", DebitMinorUnits: 1250}, {LineNo: 2, AccountCode: "FX_CLEARING", CreditMinorUnits: 1250}}}
	status, _, raw := request("POST", "/v1/accounting/journals", "", good, encode(forged))
	if status != 400 {
		t.Fatal("generic source spoof", status, string(raw))
	}
	for i, tc := range []struct {
		to     string
		amount int64
	}{{"GBP", 1000}, {"USD", 0}, {"USD", -1000}} {
		receipt := conversion(fmt.Sprintf("fx-journal-negative-%02d", i), tc.to, tc.amount)
		status, _, raw := request("POST", "/v1/accounting/fx/journals", fmt.Sprintf("fx-journal-reject-%03d", i), good, encode(command(receipt)))
		if status != 409 {
			t.Errorf("invalid journal conversion %d %s", status, raw)
		}
	}
	unused := conversion("fx-journal-rollback-conv", "USD", 1500)
	badCommand := command(unused)
	badCommand.DebitAccount = "MISSING_ACCOUNT"
	status, _, raw = request("POST", "/v1/accounting/fx/journals", "fx-journal-invalid-account", good, encode(badCommand))
	if status != 409 {
		t.Fatal("invalid account", status, string(raw))
	}
	if e = pool.QueryRow(ctx, `select count(*) from accounting.journal where tenant_id=$1 and source_id=$2`, tenant, unused.ID).Scan(&n); e != nil || n != 0 {
		t.Fatal("orphan after account rejection", n, e)
	}
	// Block the actual shared outbox while the existing draft writer is running.
	blocker, e := pool.Begin(ctx)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = blocker.Exec(ctx, `lock table platform.outbox_event in access exclusive mode`); e != nil {
		t.Fatal(e)
	}
	blockedCtx, blockedCancel := context.WithTimeout(ctx, 150*time.Millisecond)
	bounded := command(unused)
	bounded.IdempotencyKey = "fx-journal-outbox-cancel"
	_, _, blockedErr := fx.PrepareJournal(blockedCtx, tenant, "accountant", bounded)
	blockedCancel()
	if e = blocker.Rollback(ctx); e != nil {
		t.Fatal(e)
	}
	if blockedErr == nil {
		t.Fatal("blocked write committed")
	}
	if e = pool.QueryRow(ctx, `select count(*) from accounting.journal where tenant_id=$1 and source_id=$2`, tenant, unused.ID).Scan(&n); e != nil || n != 0 {
		t.Fatal("orphan after outbox cancellation", n, e)
	}
	if e = pool.QueryRow(ctx, `select count(*) from platform.idempotency_record where tenant_id=$1 and scope='accounting-fx-journal' and idempotency_key=$2`, tenant, bounded.IdempotencyKey).Scan(&n); e != nil || n != 0 {
		t.Fatal("orphan idempotency", n, e)
	}
	// Post/reverse through the existing separately authorized accounting endpoints.
	postPath := "/v1/accounting/journals/" + initial.Receipt.JournalID + "/post"
	postBody := `{"organization_id":"store","expected_version":1}`
	status, _, raw = request("POST", postPath, "", good, postBody)
	if status != 403 {
		t.Fatal("write permission posted journal", status, string(raw))
	}
	status, _, raw = request("POST", postPath, "", controller, postBody)
	if status != 200 {
		t.Fatal("post", status, string(raw))
	}
	status, _, raw = request("GET", resultPath, key, good, "")
	var posted accounting.FXJournalResult
	if status != 200 || json.Unmarshal(raw, &posted) != nil || posted.CurrentStatus != "posted" || posted.CurrentVersion != 2 || encode(posted.Receipt) != encode(initial.Receipt) {
		t.Fatal("post response lost; current projection", status, string(raw))
	}
	status, _, raw = request("GET", "/v1/accounting/trial-balance?organization_id=store&period_id="+url.QueryEscape(period.ID), "", controller, "")
	var balance []accounting.Balance
	if status != 200 || json.Unmarshal(raw, &balance) != nil || len(balance) != 2 {
		t.Fatal("trial balance", status, string(raw))
	}
	for _, line := range balance {
		if line.NetMinorUnits != 1250 && line.NetMinorUnits != -1250 {
			t.Fatal("converted posting amount", balance)
		}
	}
	status, _, raw = request("POST", "/v1/accounting/journals/"+initial.Receipt.JournalID+"/reverse", "", controller, `{"organization_id":"store","expected_version":2,"reason":"Synthetic fixture reversal"}`)
	if status != 201 {
		t.Fatal("reverse", status, string(raw))
	}
	status, _, raw = request("GET", resultPath, key, good, "")
	var reversed accounting.FXJournalResult
	if status != 200 || json.Unmarshal(raw, &reversed) != nil || reversed.CurrentStatus != "reversed" || reversed.CurrentVersion != 3 || encode(reversed.Receipt) != encode(initial.Receipt) {
		t.Fatal("historical/current split", status, string(raw))
	}
	status, _, raw = request("GET", "/v1/accounting/trial-balance?organization_id=store&period_id="+url.QueryEscape(period.ID), "", controller, "")
	if status != 200 || json.Unmarshal(raw, &balance) != nil {
		t.Fatal("reversed balance")
	}
	for _, line := range balance {
		if line.NetMinorUnits != 0 {
			t.Fatal("reversal not neutral", balance)
		}
	}
	if _, e = pool.Exec(ctx, `delete from platform.idempotency_record where tenant_id=$1 and scope='accounting-fx-journal' and idempotency_key=$2`, tenant, key); e != nil {
		t.Fatal(e)
	}
	status, headers, raw = request("POST", "/v1/accounting/fx/journals", key, good, body)
	if status != 201 || headers.Get("Idempotent-Replay") != "true" || strings.TrimSpace(string(raw)) != encode(reversed) {
		t.Fatal("retained receipt replay", status, string(raw))
	}
	if _, e = pool.Exec(ctx, `update accounting.fx_journal_receipt set requested_by_subject='other' where tenant_id=$1`, tenant); e == nil {
		t.Fatal("mutable FX journal receipt")
	}
	// A closed period prevents new FX preparation through the same existing writer.
	if _, e = ledger.ClosePeriod(ctx, tenant, period.ID, 1); e != nil {
		t.Fatal(e)
	}
	status, _, raw = request("POST", "/v1/accounting/fx/journals", "fx-journal-closed-period", good, encode(command(unused)))
	if status != 409 {
		t.Fatal("closed period", status, string(raw))
	}
	if e = pool.QueryRow(ctx, `select count(*) from accounting.journal where tenant_id=$1`, tenant).Scan(&n); e != nil || n != 2 {
		t.Fatal("expected original and reversal only", n, e)
	}
	if e = pool.QueryRow(ctx, `select count(*) from accounting.entry where tenant_id=$1`, tenant).Scan(&n); e != nil || n != 4 {
		t.Fatal("expected four immutable entries", n, e)
	}
}
````

### FILE: `db/migrations/0065_fx_journal_receipt.up.sql`

```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:db/migrations/0065_fx_journal_receipt.up.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "typed transaction, authorization, HTTP and verification glue around the existing source-admitted FX and accounting owners; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f0c415554e8267d7125c4836b3cff45b43df41f50b706bdee4e686c64b47cb57"
variables: []
secrets_allowed: false
```

````sql
begin;
-- AUTHORED relationship/receipt at the existing accounting owner, not a ledger.
create table accounting.fx_journal_receipt(
 tenant_id uuid not null,
 organization_id text not null,
 conversion_id text not null,
 journal_id text not null,
 requested_by_subject text not null check(length(requested_by_subject) between 1 and 256),
 request_key text not null check(request_key ~ '^[A-Za-z0-9_-]{16,128}$'),
 request_sha256_hex text not null check(request_sha256_hex ~ '^[0-9a-f]{64}$'),
 receipt_raw bytea not null check(octet_length(receipt_raw) between 1 and 16384),
 receipt_sha256_hex text not null check(receipt_sha256_hex ~ '^[0-9a-f]{64}$'),
 recorded_at timestamptz not null,
 primary key(tenant_id,conversion_id),
 unique(tenant_id,request_key),
 unique(tenant_id,journal_id),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 foreign key(tenant_id,conversion_id) references accounting.fx_conversion_receipt(tenant_id,conversion_id),
 foreign key(tenant_id,journal_id) references accounting.journal(tenant_id,journal_id)
);
create trigger fx_journal_receipt_immutable before update or delete on accounting.fx_journal_receipt for each row execute function accounting.prevent_posted_history_mutation();
commit;
````

### FILE: `db/migrations/0065_fx_journal_receipt.down.sql`

```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:db/migrations/0065_fx_journal_receipt.down.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "typed transaction, authorization, HTTP and verification glue around the existing source-admitted FX and accounting owners; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b6cb77d1929a4e5ce20595c95766aa191cdf72f92c93b7f8a770e9b3adb855b3"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$ begin
 if exists(select 1 from accounting.fx_journal_receipt) then
  raise exception 'FX journal receipts exist; preserve accounting history';
 end if;
end $$;
drop table accounting.fx_journal_receipt;
commit;
````

### FILE: `docs/fx-journal-runtime.md`

```yaml
block_id: "GO-EXACT-FX-SNAPSHOT-ACCOUNTING:docs/fx-journal-runtime.md:v1"
operation: CREATE
provenance: AUTHORED
source: "typed transaction, authorization, HTTP and verification glue around the existing source-admitted FX and accounting owners; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "88aeccfcdaaac54223a9bab4cd134f7d9b19f265c372e732697c790a11eda8fb"
variables: []
secrets_allowed: false
```

````markdown
# FX conversion to accounting journal

This optional module uses one immutable direct-base conversion receipt and the existing accounting draft, posting and reversal writer. Enable the existing hash-bound FX profile; no additional account or service credential is introduced. AUTHORED code binds admitted owners and selects no account, exchange-rate feed, gain/loss, tax or valuation policy.

An accountant with accounting:write creates a conversion using POST /v1/accounting/fx/conversions and retains its request key and returned ID. POST /v1/accounting/fx/journals then accepts organization_id, conversion_id, conversion_request_key, period_id, debit_account, credit_account and description, with a fresh Idempotency-Key. The actor must own the conversion receipt. Both input and output must be positive; output currency must equal the snapshot local currency. Accounts and period are supplied explicitly. The amount and date come from the immutable receipt and cannot be overridden in this request.

Preparation creates one balanced draft journal, its immutable binding, shared idempotency record and outbox events in one PostgreSQL transaction. It creates no posted ledger entries. One conversion receipt backs at most one journal; after reversal a new accounting operation requires a new conversion receipt. The generic journal creation endpoint reserves source_type FX_CONVERSION for this typed path. The conversion receipt remains CONVERSION_RECEIPT_ONLY. The new receipt remains FX_JOURNAL_PREPARED even after posting or reversal.

If the response is lost, GET /v1/accounting/fx/journals/result?organization_id=... with the same Idempotency-Key retrieves the original receipt plus current_status/current_version. It requires accounting:read and the same actor and organization. Exact replay survives expiry/removal of the transient idempotency row. A changed request, actor or second request key for an already bound conversion is rejected. Never infer current posting from the historical receipt effect.

Post and reverse through the existing accounting endpoints and their separate accounting:post/accounting:reverse permissions and expected version. No posting authorization is granted by preparing a draft. Use the existing trial balance and journal queries to inspect the result. The demonstrated fixture converts EUR 1000 minor units to USD 1250, prepares a draft, explicitly posts, reverses through the existing writer, and returns net zero in the trial balance.

Migration 0065 follows accounting0011 and FX0061. Downgrade refuses to remove nonempty FX journal history. Roll back application deployment while preserving the schema and immutable rows; do not erase receipts to force a downgrade. Retention follows the existing accounting owner, not the short-lived idempotency table. Other currencies, zero/negative conversions, remeasurement and automatic realized/unrealized FX gains are outside this adapter's claim and require an explicit source-admitted policy if selected later.

Evidence: reconstruction_evidence/FX_JOURNAL_CONNECTION_V402.md in the source library binds actual HTTP/RS256/JWKS/PostgreSQL concurrency, recovery, permission, invalid input, outbox-cancellation rollback, posting and reversal tests, exact reconstruction and unchanged upstream arithmetic. Global composition security/SCA, delivery and production acceptance are separate gates.
````

## 6. Configuration surface

Read config/fx/reference-profile.json and reference-rates.json with the exact
selected SHA; docs/fx-conversion-runtime.md describes supported direct-base,
dated snapshots and explicitly configured rounding. No live rate is inferred.
Accounts/period and authorized actor are required by docs/fx-journal-runtime.md.

## 7. Dependency bill

No added Go module. Selected Go1.26.8 math/big, pgx5.10.0 and PostgreSQL18.6
provide representation/persistence. The fixed BCApps source lock, adaptation
and notices are in docs/provenance/BC_FX_SOURCE_LOCK.json and its adjacent files.
The complete composition supplies identity and the existing accounting owner.

## 8. Apply order

Compose with accounting0011, identity, shared idempotency/outbox and the existing
host, then apply0061 and0065 in order. Journal preparation is separate from
authorized posting/reversal. Preserve immutable conversion/journal history;
nonempty downgrade is rejected. Application rollback keeps that schema/history.

## 9. Verification

The original source/oracle/fuzz and connected HTTP/JWKS/PostgreSQL concurrency,
expiry, replay, outbox-rollback, posting and reversal results remain governed by
FX_JOURNAL_CONNECTION_V402.md and the current T2802 control receipt. Rebuild
the selected composition and compare its exact manifest before use. Outer
documentation repair321 changes no executable, fixture, schema or source hash.

## 10. Reconstruction evidence

Current reference319 contains every declared selected payload byte unchanged
by this heading repair. Existing source/target receipts retain their original
revision and narrow scope. Global signed release acceptance remains T2810;
no source algorithm or production approval is created by a documentation change.
