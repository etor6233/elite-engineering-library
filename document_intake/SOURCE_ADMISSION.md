# V403-DOCS-FINAL-01: local document intake, review and preview

Selected scope: existing 21-class owner, versioned channel receipts, real browser/BFF/Go/PostgreSQL review, and bounded PDF preview. Glue is AUTHORED; it is not Mozilla/Microsoft/AWS code. No production claim.

| Gate | Evidence / decision |
|---|---|
| G0 | Mozilla PDF.js v6.3.289 tag commit equals npm gitHead (`1c8020a7d4e43668ac287a3ecf9a8dbea17e4c56`). Registry archive SHA512 and included files SHA256 verified in SOURCE_LOCK.json. |
| G1 | Exact Apache-2.0 license and bundle headers preserved. Two VERBATIM browser bundles only; no optional Node canvas, fonts, WASM, viewer brand assets, or models. Local glue remains LicenseRef-Workspace-Owner. Existing AVM subpath notices preserved. |
| G2 | R2 U3 independent review, original hash, same inbox and recovery. PDF/JPEG input contract unchanged; HTML original rejected, reviewed fields rendered as escaped React HTML. No original HTML execution. |
| G3 | Authenticated original read and SHA validation, then bytes transferred to opaque sandbox and pinned worker; canvas output only. No viewer scripting manager, annotations, XFA, eval or external fetch. Intake only creates QUARANTINED originals; processing/approval remains explicit and independent. |
| G4 | Real API build, TypeScript noEmit, 12 browser checks, 18 inherited MIME tests, 8 intake boundary tests. Red/green StrictMode load bug preserved. |
| G5 | Auth401, fixed asset allowlist404, sandbox without same-origin, source/nonce checks, CSP connect-src none, PDF OpenAction JavaScript no network/popups, HTML rejection, existing owner tenant/organization/permission checks preserved. |
| G6 | 2MiB, 30pages, 16,777,216 total pixels, 4096 dimension; 15s renderer and 20s host deadline. Byte/page/pixel excess rejection executed. Zoom reuses bounded canvas pixels. Local timings in JOURNEYS.json, not a production throughput/SLO claim. |
| G7 | Loading/rendered/rejected, cancellation destroys frame/worker, authorized download fallback; immutable overlay rollback without DB migration. Reopen on advisory, runtime, corpus or renderer changes. Live monitoring/deployment NOT_RUN. |
| G8 | Guarded manifest, source lock, harness, canonical exact-byte materialization and reconstruction receipt. Selected local composition proved; root must validate integration after merging other overlays. |

Official sources checked: https://github.com/mozilla/pdf.js ; https://mozilla.github.io/pdf.js/api/draft/module-pdfjsLib.html ; https://react.dev/reference/react/StrictMode ; https://react.dev/reference/react/useEffect ; https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/iframe ; https://learn.microsoft.com/en-us/azure/event-grid/event-schema-blob-storage . Existing MIME and Microsoft AVM SFTP packs reconstructed exactly with canonical materializer; pack hashes in SOURCE_LOCK.json.

Advisories snapshotted in evidence/pdfjs-advisories.json: GHSA-wgrm-67xf-hhpq fixed4.2.67 and GHSA-hq66-cqwq-w95j fixed6.2.108. Selected6.3.289 is above both patches. No assertion that future vulnerabilities are impossible. No viewer scripting manager is instantiated; CSP and opaque origin add containment.

Final runtime: evidence/20260922T143539Z. Next/BFF, production Go owner, fixture OIDC, fresh owned PostgreSQL and current migrations. Upload/process/propose/different reviewer/commit, lost real response followed by GET/reload recovery. SQL: 3 originals, 1 committed row, distinct reviewer. Email packing list and SFTP proforma enter that same inbox; duplicate delivery produces zero extra mutation. Existing revision8 21-class PostgreSQL proof is reused; local typed extraction is not a real OCR claim.

The original detail effect failed under StrictMode because an aborted old read changed the new mount's state. Successor binds completion to the effect generation; StrictMode stays enabled. Red evidence 20260922T141527Z, green20260922T141720Z and final runtime above.

SFTP normalizer requires authenticated upstream delivery, exact topic, SftpCommit, URL/etag/size/scope and retained-object binding. No event URL is fetched. Descriptors come from a trusted provider/retention owner, not anonymous HTTP. CLI takes bearer via environment only. Raw MIME extractor and AVM infrastructure are included; live SMTP/SES/SFTP/EventGrid identity, delivery and enforced retention remain NOT_RUN. Static AVM test passes; pinned Bicep compile was not rerun.

Original and review data share desktop space; mobile preview has fit/zoom. Screenshots remain PENDING_USER_APPROVAL; phone and screen reader NOT_RUN. Canvas plus HTML fields/download is not a claim of a full accessible PDF editor or signature verifier. Unsupported/encrypted/resource-heavy PDFs use the download fallback.

Apply DOCS_FINAL_MANIFEST.json to an exact compatible predecessor. Two parser files in third_party/pdfjs are traced into Next standalone output; package/lock unchanged. Roll back by selecting the previous layer, with no DB mutation. qualification/run.py uses pinned local fixture toolchains and new evidence directories; portable root integration is a separate root-owned gate.
