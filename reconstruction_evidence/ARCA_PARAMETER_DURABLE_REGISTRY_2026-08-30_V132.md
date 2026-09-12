# ARCA parameter durable registry — reconstruction evidence V132

## Authority and provenance

V132 connects the seven WSFEv1 parameter families acquired in V131 through the existing Unix-socket worker and records their exact safe response in the existing Go/PostgreSQL fiscal owner. ARCA WSFEv1 manual 4.6 governs the operation families; Microsoft-generated WCF types remain the protocol boundary; PostgreSQL 18.6 provides the durable constraints. All six new files and seven changed integration blocks are local `AUTHORED` code. They are not represented as copied ARCA, Microsoft, Go or PostgreSQL source and contain no parameter table, WSDL, proxy, certificate, token, signature or provider message.

- <https://arca.gob.ar/ws/documentacion/ws-factura-electronica.asp>
- <https://arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG-v4-0.pdf>
- <https://learn.microsoft.com/dotnet/core/additional-tools/dotnet-svcutil-guide>
- <https://www.postgresql.org/docs/18/>

## Exact artifacts

- `GO-ARCA-FISCAL-ISSUANCE-API 0.4.0`: 17 files, 95,650 bytes, SHA-256 `9cc5dad6d1a0add8541abd0eaa671e4099f3f4afa8a1d4cd7a8f07eae6163106`.
- `ARCA-WSFE-UDS-WORKER 0.2.0`: 13 files, 61,664 bytes, SHA-256 `6eadfca2fab528cf4ca0dec428af6a2754144487e9a110570a907d9c76d73372`.
- `GO-ELECTROMOBILITY-APPLICATION 1.4.1`: SHA-256 `57ae10d5f475a35c55f4c4f16dfe5ce6f795f2ab41f625116b5198fac3b081c6`.
- Backend plan SHA-256 `01845776e10271a6b80b9837092081640d4fa224d0c5a9ac0fc0c4539acbea39`; clean record SHA-256 `5d06a5469ed2d31c146d7c23846c67d798fe64d802fb6eb11771461d6517f4f6`.

## Executed gates

1. Fiscal and UDS packs materialized from Markdown 17/17 and 13/13 with exact hashes. The backend composed 223 files from 24 packs into a new directory.
2. Go 1.26.7 full tests, focal `internal/fiscal/...` and `internal/platform/postgres`, vet and both worker/API builds passed. The registry rejects unknown kinds, CUIT/class divergence, duplicates and invalid date ranges.
3. The runner injected the hash-verified homologation proxy rather than distributing it. .NET 10.0.400 locked restore and warning-free build passed; Kestrel performed a native UDS journey and returned `ARCA_WSFE_WORKER_TESTS_PASS tests=12` plus `parameter_bridge=1 production_admitted=false`.
4. PostgreSQL 18.6 applied migrations 0001–0013 twice, once in staging and once from the canonical Markdown composition. Test 0013 proved that an unapproved snapshot is absent from `fiscal.approved_parameter_snapshot`, parameter history cannot mutate and exactly one explicit approve/reject decision is admitted. Both temporary clusters stopped cleanly after their successful gates.
5. `LIB-FAIL-1332` through `LIB-FAIL-1339` preserve incorrect path/build assumptions, the permissive fake, wrapper/patch failures and PostgreSQL lifecycle lessons; none is omitted or converted into a false PASS.
6. The final governing gate returned `VERIFY_LIBRARY_PASS packs=83 materialized_files=868 markdown_files=471`, backend profile 223 and web profile 80. All pack identities, sections, manifests, blocks, hashes, provenance values and 36 composition profiles agreed.

## Admission boundary

This closes reconstructible acquisition transport and durable separation between observed data and business approval. It does not schedule refreshes, expose an authenticated administrator API, prove a live ARCA response or declare an approved snapshot correct for a specific taxpayer. Certificate/service association, homologation, refresh cadence, authorized tax owner, legal/accounting review and target evidence remain mandatory. Production readiness still requires the concrete CDN/WAF, IdP, providers, PostgreSQL recovery, load, offensive security, deployment/rollback and business-acceptance receipts.
