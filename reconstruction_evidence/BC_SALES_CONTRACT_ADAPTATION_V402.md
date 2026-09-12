# V402 staged BC sales derivation — exact narrow gate

PASS_STAGING_ONLY. No canonical publication, profile promotion, TEST02 closure,
production or whole Commerce/Journey provenance claim. Parent owns the final merge.

## Provenance and scope

GO-BC-SALES-CONTRACT-ADAPTER0.1.0 candidate adds6files:2ADAPTED and4AUTHORED local
clock/tests/docs. It requires the existing GO-BC-EXACT-AMOUNT-ADAPTER0.1.0 owner
for the verbatim MIT notice; there is no duplicate license path. Commerce0.6.3
and Journey0.10.20 are staging candidate versions only, with one changed caller
file each. They retain their AUTHORED provenance and current admission conditions.

The exact source map and representation limits are materialized at
sales-derivation-reference/docs/provenance/BC_SALES_DERIVATION.md. Existing source
files were rehashed against fixed receipts, not fetched again. No AL code/suite
was executed. The source transformations are price active/start/end/currency
eligibility and quote type/header/line filtering/copy/retarget. Consent, expiry,
pricebook exclusivity, scheduling and financial states remain outside the claim.

## Results

-180 PostgreSQL differential vectors: status3 x starting5 x ending6 x currency2;
  source inclusive ending mapped to the existing microsecond half-open domain.
  Includes exact boundaries, future starting date (official T176 scenario), wrong
  currency, inactive status and unbounded ending. No price ranking or currency
  conversion invented.
-Real PostgreSQL18.6 with54 hash-identical migrations: existing Commerce
  price/order/stock/payment-request and Journey persistence/isolation/replay
  suites PASS; source functions connected, zero skips. Existing Journey fixture
  cleanup uses its documented disposable-database trigger bypass; no live data.
-Quote conservation/scoping unit tests PASS: wrong document type rejected,
  foreign quote lines filtered, copied fields unchanged, output storage separate.
  Native fuzz gate3seeds/70,545executions/3s budget PASS. This is local source-
  conservation verification, not Microsoft-authored test execution.
-Vet affected packages and full reference module build PASS. No new module
  dependency, runtime, schema or API. PostgreSQL stopped on completion. Accounting
  tests were not repeated; the preexisting amount fix is retained by baseline.
-Three candidate packs materialized47/47 byte-identical blocks. Only6new files
  and2existing callers differ from the exact business-reference baseline. No D
  root or canonical file was modified. Parent must merge concurrent D changes.

## G0–G8 for the narrow candidate

| Gate | Result and evidence |
|---|---|
| G0 | PASS: same official BCApps commit2eae56d704a1fd035d104f333602aea7091b7749; source code and T176 file hashes rechecked; no reputation-based ownership claim. |
| G1 | PASS: MIT bytes c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383 exist through the explicit prior adapter dependency; source files ADAPTED, glue/tests/docs AUTHORED; no duplicate owner. |
| G2 | PASS: exact function/line mapping; Date→microsecond half-open endpoint conversion; nonblank currency excludes BC fallback; copied subset and omitted BC policy stated. |
| G3 | PASS: three existing price queries share the predicate; existing AcceptQuote invokes the transferred header/line inside the same transaction. No parallel domain or invented transition. |
| G4 | PASS: 180 differential negatives/edges, field-conservation fixtures, connected Commerce/Journey, vet and build; not an AL suite claim. |
| G5 | PASS narrow delta: tenant/org/customer and explicit book/variant/order-currency scopes unchanged; immutable quote input preserved; wrong source type/foreign line scope tested; no secret/network/dependency. Broader TEST03 remains separate. |
| G6 | PASS local: constant SQL predicate/context and linear source-line copy; current caller supplies one line; finite3s fuzz3seeds/70,545executions. No production-load guarantee. |
| G7 | PASS: unchanged durable format/migrations, existing replay and isolation tests; rollback callers before removing adapter, preserving prior bcamounts fix. |
| G8 | PASS staging: three candidate packs reconstruct47/47 exact; change manifest8files, patch, source map and receipts preserved. Canonical merge/profile update remains parent work. |

## Hash receipts

| Artifact | SHA-256 |
|---|---|
| `sales-derivation-manifest.json` | `1155967d7251d68a00ec7ed270d109b7f8410237459dcb37b8841a3a4ec93ad5` |
| `sales-derivation.patch` | `2006847171d14a108e1e1e3f0ec83d18c65b41b33e6c62cb7d30d5425f49a50a` |
| `sales-roundtrip-receipt.json` | `dfdc3b8d7914c275327d9fb19bc125de049fb9232ca84994568743edb5ffd19e` |
| `sales-derivation-pg-1/connected.log` | `38aaa0ad9cc6a6c22ce72cd80e46df06920a41ecd2718ba2ffbfca8319c1cfdf` |
| `sales-derivation-pg-1/quote-unit.log` | `19f8cc848faa0b36ec3accb3530f7e1bcff93093006c42f18d182fdb9d4eafbd` |
| `sales-derivation-pg-1/vet.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `sales-derivation-pg-1/build.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `sales-derivation-pg-1/result.json` | `621f6ac47f76ccb732231af1067945c0821a387591d9dcc7d27c091bde693e29` |
| `sales-fuzz.log` | `24ad49b527fb549a3b8756a36ab79597a8db21707db5e63a53ab696dc398deb3` |
| `sales-fuzz-receipt.json` | `3ad383ba606fb23ce4df20072fb57537a631303cbebd03ae789c7993d52c3c99` |
| `sales-constraint-expedients.md` | `ac190f9ec6f30fb904346b4efff89c79dc04bbdae5ce7e2e2eeb2c5b5bd7b579` |

Remaining-policy evidence is in sales-constraint-expedients.md: locks/CAS and mechanical quota enforcement can be glue; hardcoded scheduling30min/8h/cap100 and occupancy/calendar rules are not reclassified. No source acquisition or new policy implementation occurred.

## Canonical integration

The preceding receipt describes its original staging run. Parent subsequently verified every before/after hash, published the three identical packs and merged all8outputs into payment-connected-reference with no conflict. All four Commerce/Journey plans select the new adapter and prior amount/license owner. Current reference72packs/834files; the separate HTTP metrics profile adds8files. The underlying47file reconstruction and tests are reused because source bytes match; no repeated suite or full TEST02 claim is inferred.
