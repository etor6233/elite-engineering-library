# Microsoft-platform ARCA WSAA Credential Core — Reconstruction Evidence V125

## Scope and provenance

V125 converts the generated WSAA transport into immediately composable credential code. All eighteen files are local `AUTHORED` integration code governed by the official ARCA WSAA request/response contract, the generated Microsoft WCF interface and Microsoft platform cryptography. No line is represented as copied ARCA or Microsoft product source.

The implementation uses the exact package graph already proven in V123/V124: .NET SDK 10.0.400, System.Security.Cryptography.Pkcs/Xml 10.0.11 and System.ServiceModel.Http 10.0.652802. Four NuGet lockfiles preserve direct/transitive content identities.

## Materialized result

- `MICROSOFT-ARCA-WSAA-CREDENTIAL-CORE 0.1.0`, SHA-256 `e2c25fdcce5ffe1c3192fabfd79c40637e79b847fa5f7b87e6ff65905b8b9b55`, materializes 18/18 files.
- Official `loginTicketRequest` XML uses a positive externally supplied ID, UTC skew/expiration, service `wsfe` and a twelve-hour maximum.
- Attached CMS/PKCS#7 is created by Microsoft `SignedCms`; certificate validation requires an active RSA private key of at least 2048 bits.
- Certificate loading is exact-thumbprint from the platform certificate store. PFX bytes, passwords, token and sign are prohibited from configuration/logs/receipts.
- Response parsing rejects DTD, external resolution, payloads above one MiB, missing/oversized secrets and expiration outside the admitted interval.
- Credential refresh is memory-only and single-flight. Thirty-two concurrent requests produce one WSAA transport call.
- `GeneratedWsaaTransport` implements the exact generated `LoginCMS.loginCmsAsync(loginCmsRequest)` signature; three integration tests prove mapping and reject blank CMS/response. Manual glue is zero.
- Current inventory after close: 80 packs / 822 materializable files; provenance 680 `AUTHORED`, 37 `ADAPTED`, 105 `VERBATIM`. Backend remains 19/165 and web 6/80.

## Exact executable proof

- Clean Markdown materialization: 18/18 PASS.
- Four locked restores, core and generated-interface builds: zero warnings/errors.
- Nine credential tests: options, request shape, attached verified CMS, thumbprint, response, DTD, expiry, 32-way single-flight and no-private-key rejection PASS.
- Three generated transport tests PASS; `manual_glue=0`.
- Current NuGet vulnerability scan: no vulnerable packages reported.
- Final receipt: `tests=12 build_warnings=0 vulnerable_packages=0 secrets_logged=0 manual_glue=0 production_admitted=false`.

## Failure memory

The first test fixture crossed top-level C# scope boundaries and was corrected. More importantly, the first streaming parser skipped sibling elements and rejected a valid response; the suite caught it before admission. The parser now uses secure bounded XML loading. Both failures are preserved in the library ledger.

## Production boundary

The credential core and generated transport are ready to compose, but live WSAA still requires the project's ARCA-issued homologation certificate, service association and target endpoint. Fiscal rules, durable invoice state, sequential numbering, idempotency, CAE reconciliation, accounting mapping and operational homologation are separate required capabilities. A concrete project remains unready for production until CDN/WAF, IdP, providers, PostgreSQL/recovery, load, offensive security, deployment/rollback and business acceptance are demonstrated.
