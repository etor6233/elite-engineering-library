# GO ARCA Fiscal Recipient VAT — V128 reconstruction evidence

Date: 2026-08-30  
Scope: mandatory recipient VAT-condition preservation before exact WSFEv1 bridge construction; no production claim.

## Official authority and discovered gate

- Current ARCA WSFEv1 developer manual 4.6: `https://www.arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG.pdf`.
- Official PDF: 4,081,946 bytes; SHA-256 `809a68756c06a5cfa011a5e6db850e402d4fb33156ad0cb0029259a358d5305f`.
- Revision 4.0 in that manual adds `CondicionIVAReceptorId`; error 10246 states the recipient VAT condition is mandatory under RG 5616. The generated Microsoft svcutil proxy exposes the exact field plus its `Specified` flag.
- V127 did not persist this value. V128 adds `recipient_vat_condition_id` to the request, service validation, PostgreSQL owner, replay/claim/read path, HTTP body and tests. It is explicit input; no document type, CUIT or customer classification is used to invent it.
- Page 21 also fixes common error 602 as “No existen datos en nuestros registros”; pages 75 and 187–193 govern timeout recovery via `FECompConsultar`. Those facts are retained for the subsequent SOAP bridge, not yet claimed as implemented here.

## Reconstructed artifact

- `GO-ARCA-FISCAL-ISSUANCE-API` 0.3.0: 75,748 bytes; SHA-256 `8c322bab35cd1f190468f724a19f0b5ca87f850780b4ffea313b11a118f30809`.
- `GO-ELECTROMOBILITY-APPLICATION` 1.3.0: SHA-256 `a8697be78c9512e62726289fe616c8fca012ec931c339a5681c5b9a0fc90a127`; compatibility now requires fiscal 0.3.x.
- Backend plan SHA-256 `f36b1074f3f2e8882a29cce7fab2923b2639e13c26967a2320f5383643dcec94`.
- Exact composition remains 20 packs / 176 implementation files. Provenance remains 691 `AUTHORED`, 37 `ADAPTED`, 105 `VERBATIM` over 833 blocks.

| File | SHA-256 |
|---|---|
| `db/migrations/0012_arca_fiscal_issuance.down.sql` | `5dc38808804f75fcca12f9a596abacd45280e1fa37c9ec446f9f9e493cdb6385` |
| `db/migrations/0012_arca_fiscal_issuance.up.sql` | `aafbe0f1eb46b8f528a06002803c47869fb162918bf91c1dd8b083ddf71f4ecc` |
| `db/tests/0012_arca_fiscal_issuance.test.sql` | `c75fe8b064437e3064ee5da682b868d863a4b09aeadbaa0e37ae4140c565465f` |
| `internal/fiscal/processor_test.go` | `c09eeeaefa54177bdb6e6c4912ac6501210feb315bd95b237f1ada7c2843a7f5` |
| `internal/fiscal/processor.go` | `77fe23edacdc2235e8a521e8447feee16336f612cd0315224f58a21c20a4a313` |
| `internal/fiscal/service_test.go` | `4be0e208f722b4886379b9d234928517fc0a3ecbb4a5e754eebccf38cabf94b1` |
| `internal/fiscal/service.go` | `02435da4fe50042cacd36d4e7f95d9f8571e931109db095ed0e738bdc6d63ff1` |
| `internal/platform/httpapi/fiscal_test.go` | `5e5fe4e8dc27fad1225c4f4a7810129dc8165eba38f1c1410f4ab529e034deb9` |
| `internal/platform/httpapi/fiscal.go` | `76e2450f7a954ca34ee619a24b0ae9bd550b56d2e9ccbb72079dd2ea2b151e14` |
| `internal/platform/postgres/fiscal_integration_test.go` | `0625d26cd8ddbd9e76c9c523471993e9304c2dbfd8c0864133a9d10ae6f2953d` |
| `internal/platform/postgres/fiscal.go` | `249f20587c537dfe1d3659d33ba2b1abe7ff03777ebb54a5cac0dc16a29bb15b` |

## Executed gates

- Markdown materialization: fiscal 11/11 and backend 20/176, exact hashes and record: PASS.
- Official Go 1.26.7 over the canonical composition: formatting clean, full tests twice with PostgreSQL enabled, vet and build: PASS.
- PostgreSQL 18.6 over canonical files: migration 0012 down/up, SQL invariant test and the exact fiscal integration test twice with two asserted PASS events: PASS.
- Tests reject a missing recipient VAT condition and prove the admitted value survives create/replay/claim/read unchanged.
- The global structural/library gate is rerun after the evidence and documentation update; its result governs current status.

## Failures retained as lessons

- `LIB-FAIL-1313`: the isolated pack tree had no `go.mod`; compilation moved to a full composed backend instead of fabricating dependencies.
- `LIB-FAIL-1314`: quoting collapsed PostgreSQL `-p` and `-h` into one port value; the server log supplied the exact cause and the arguments were passed as one correct `-o` value.
- `LIB-FAIL-1315`: launching a PowerShell child process flattened the `-Files` array; the updater was invoked directly in the current PowerShell process.
- `LIB-FAIL-1316`: first PDF output used cp1252 and failed on a glyph; extraction was repeated with UTF-8 and no truncated text was used as authority.
- `LIB-FAIL-1317`: one multi-file patch used a non-exact test context and was rejected atomically; exact file slices were read and smaller patches applied.

## Exact condition

This closes one more representational prerequisite; it does not validate whether a chosen condition ID is legally correct for a voucher. Production remains fail-closed until the project fetches and pins current ARCA parameter tables, applies approved tax/voucher rules, supplies real credentials and homologation receipts, and proves the full SOAP bridge, QR/PDF/notes/contingencies, load, recovery, security and accounting/business acceptance.
