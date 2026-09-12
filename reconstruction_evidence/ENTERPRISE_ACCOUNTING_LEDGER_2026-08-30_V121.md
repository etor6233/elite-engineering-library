# Enterprise Accounting Ledger — Reconstruction Evidence V121

## Scope and provenance

V121 adds one operational general-ledger boundary after the existing business/subledger owners. All new Go/SQL is local `AUTHORED`; no line is represented as Microsoft code. Exact design authority is the MIT snapshot `microsoft/BCApps@31a860b527f0dc72c7a44a255d7e7d403cfa4789` (archive `e3151b040df39cada83d41fcf5d9e6cdff1d8fddf226934c1ad21a2bdebc7210`, root license `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`):

- `GenJournalLine.Table.al` `8fb790053a04517d02f425f45e0d02dacecaf8b5ef2bd2bc23fa6b398342769b`;
- `GLEntry.Table.al` `08655960b539eb5afbea0c91f6189af07b4f83c68f7de08d5d590da6591fd9dd`;
- `GLRegister.Table.al` `2f7a870f2868a39c53d89b97ed831741579afef506458a5d701bd8bf6c6c807d`;
- `AccountingPeriod.Table.al` `d0b2628475e1abb64208139256ee3ec9f7d45e81564ab94c250564b2779e1ed6`;
- `VATEntry.Table.al` `c7aff50430b505d30ebf2ad6b3ee7f0a213510a2bf32828f5ffcf54a6d8d1fe4`;
- `GeneralPostingSetup.Table.al` `033fac045efae2df61be9d1d7ab22006acbc9f53c87a828bf74d0e96388ac580`;
- `ReversalEntry.Table.al` `13e17b83c3a2d685060ac151a9f5ef0b156238036f034515ccb260d0b2eb835d`;
- `GenJnlPostLine.Codeunit.al` `fdecc9a5e52831552e9addfe4c7e66d1b53fa3425c95b2d49d606478a24d54f3`.

These sources govern journals, entries, registers, periods, posting setup and reversal boundaries. VAT source is evidence that tax is a separate ledger concern; V121 intentionally does not claim Argentine tax or statutory localization.

## Materialized result

- `GO-ENTERPRISE-ACCOUNTING-LEDGER-API 0.1.0`, SHA-256 `d6a788b518904c56b6c2c0a35c745827964f9d3caafa9567719784dec36b9daa`, materializes nine files.
- Migration 0011 creates tenant accounts, non-overlapping periods, source-unique balanced journals, immutable lines/registers/entries and trial-balance index.
- Protected commands create accounts/periods/journals, post, reverse and close; reads are tenant/organization/period scoped.
- PostgreSQL rechecks open period, posting date, active account, exact balance, source, organization and optimistic version. Serializable posting yields one success/one conflict.
- Reversal creates a new posted journal with swapped debit/credit and per-entry origin; original history remains. A closed period rejects further journals.
- `GO-ELECTROMOBILITY-APPLICATION 1.2.0`, SHA-256 `67c751b2c9693e5a15b43984829984e8c14919a5fbc979d78f8be8efd0de4b3c`, wires the module.
- Current inventory after global close: 77 packs / 791 files; backend 19/165; provenance 654 `AUTHORED`, 32 `ADAPTED`, 105 `VERBATIM`.

## Exact executable proof

- Markdown-only composition produced 165 files from 19 packs; accounting separately materialized 9/9.
- Official Go 1.26.7: two consecutive full test runs, vet and both command builds PASS.
- Official PostgreSQL 18.6: fresh database, migrations/tests 0001–0011 PASS with `ON_ERROR_STOP=1`.
- Repository integration proved overlapping-period rejection, exact balance, one-winner posting, organization-scoped trial balance, equal/opposite reversal, second reversal rejection and closed-period rejection.
- Migration 0011 down/up/test PASS.
- Global verifier is recorded after canonical inventory alignment.

## Failure memory

V121 detected the already-known `NULLS NOT DISTINCT` trap before execution and replaced it with a partial unique index only for non-null original-entry links. Two atomic app/document patches were rejected because they used an abbreviated long line; exact context was reread and applied. These are recorded in the library ledger and not hidden.

## Production boundary

This is an operational double-entry primitive, not an audited or statutory accounting product. A concrete project still needs an approved chart/mapping, fiscal localization, taxes, e-invoicing, currency/FX, bank/provider feeds, consolidation, access segregation, period-close policy, audit retention and target acceptance. Production remains unready until its CDN/WAF, IdP, PostgreSQL/recovery, load, offensive security, deployment/rollback and business acceptance are demonstrated.
