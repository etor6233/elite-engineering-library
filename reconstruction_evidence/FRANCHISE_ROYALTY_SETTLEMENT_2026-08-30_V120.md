# Franchise Royalty Settlement — Reconstruction Evidence V120

## Scope and provenance

V120 adds one bounded context after the existing franchise agreement and commerce/payment owners. All new Go/SQL is local `AUTHORED` workspace-owner code; no line is represented as Microsoft source. The governing public authority is the exact MIT snapshot `microsoft/BCApps@31a860b527f0dc72c7a44a255d7e7d403cfa4789` (archive SHA-256 `e3151b040df39cada83d41fcf5d9e6cdff1d8fddf226934c1ad21a2bdebc7210`, root license SHA-256 `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`):

- `src/Layers/W1/BaseApp/Sales/Reports/SalespersonCommission.Report.al`, SHA-256 `4cff73d35c960415e6c77149b3a6e2533b0925f90fa1ba312a6de52f127b36fd`;
- `src/Layers/W1/BaseApp/Finance/GeneralLedger/Ledger/GLEntry.Table.al`, SHA-256 `08655960b539eb5afbea0c91f6189af07b4f83c68f7de08d5d590da6591fd9dd`;
- `src/Layers/W1/BaseApp/Finance/GeneralLedger/Reversal/ReversalEntry.Table.al`, SHA-256 `13e17b83c3a2d685060ac151a9f5ef0b156238036f034515ccb260d0b2eb835d`;
- `src/Layers/W1/BaseApp/Bank/Reconciliation/BankAccReconciliation.Table.al`, SHA-256 `2fac6264945000836bd27f0d48e74be6db1dd1bab1cf6675bc3ae32d06ab4b62`;
- `src/Apps/W1/Shopify/App/src/Payments/Tables/ShpfyPayout.Table.al`, SHA-256 `7503e221f987cecdc3401aec1d3b41e73e00e4f04dcc2f877a555b351ec75f31`;
- `src/Apps/W1/Shopify/App/src/Payments/Tables/ShpfyPaymentTransaction.Table.al`, SHA-256 `e8f8cb2bf3d6b1b276bc2cc96f7cfbcfb6538e7bc2084499dc3b0f25a7c42db6`.

These exact sources establish commission basis/rate, durable entry identity and amount/currency, compensating reversal, payout transaction detail and reconciliation difference patterns. They do not supply REVESTEX contract rates, Argentine tax/accounting rules or Go/PostgreSQL source. Google SRE governs proportional launch/release gates and NASA governs component verification versus target validation; neither is claimed as code authorship.

## Materialized result

- `GO-FRANCHISE-ROYALTY-SETTLEMENT-API 0.1.0`, SHA-256 `eb968e8141961dcb598386b50dbd7e826afca6b6dc09bcabbe96cde02e15bda7`, materializes nine files.
- Migration 0010 creates one royalty owner with immutable effective policies, captured/refunded accruals, settlement runs/lines and reconciliation evidence.
- The caller cannot submit amount, currency, agreement, policy or royalty result. PostgreSQL resolves an exact organization-scoped payment state/version, one active agreement and one effective currency policy, then calculates minor units with numeric rounding.
- Duplicate provider event/payment state, overlap policy, cross-scope payment, stale version, empty close and second reversal fail closed.
- Closing is serializable and advisory-locked; two simultaneous close attempts yield exactly one success and one conflict.
- Reversal creates equal-and-opposite lines and a separate outbox identity; it never edits accrual/line history. Reconciliation stores expected, actual, difference and `matched`/`mismatch` explicitly.
- `GO-ELECTROMOBILITY-APPLICATION 1.1.0`, SHA-256 `813080529824ad381fbb6630b49a2863eac6c741ec2bb39c57f69ba5d31779e1`, wires the protected module.
- Current inventory after global close: 76 packs / 782 files; backend 18/156; provenance 645 `AUTHORED`, 32 `ADAPTED`, 105 `VERBATIM`.

## Exact executable proof

- Markdown-only composition produced 156 files from 18 packs; the new pack separately materialized 9/9 with zero SHA mismatches.
- Official Go 1.26.7: two consecutive full `go test ./...` runs, `go vet ./...` and both commands build PASS.
- Official PostgreSQL 18.6: fresh database, migrations 0001–0010 and SQL invariant test PASS with `ON_ERROR_STOP=1`.
- Live repository integration proved rate calculation `10001 × 650 / 10000 = 650`, idempotency, scope, policy overlap, one-winner close, reversal net zero and reconciliation difference.
- Migration 0010 down removed the schema; up and SQL test passed again.
- The global library verifier is the final inventory/contract gate recorded after canonical documentation alignment.

## Failure memory

The library ledger records identities 1277–1285: compositor path, workdir-relative staging, timezone-dependent civil date, invalid derived event identity, global fixture key, PostgreSQL parameter inference, residual immutable test evidence, rejected narrative patch and blocked compound cleanup. Each has a proved correction; none was suppressed or converted into a false PASS.

## Production boundary

This proves a reusable contractual royalty subledger and API, not the correctness of a project’s agreement/rate, tax, general ledger, statutory invoice, bank/provider payout, authorization model or jurisdiction. A concrete project remains unready for production until those values and adapters are approved and its CDN/WAF, IdP, PostgreSQL/recovery, load, offensive security, deployment/rollback and business acceptance are demonstrated in the target environment.
