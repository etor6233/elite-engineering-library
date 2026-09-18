# Third-Party Notices and Provenance Boundary

## V322 encoded license preservation

OFFICIAL_UPSTREAM_ACQUISITION_CORE0.4.78 adds three AUTHORED files (opaque
transport, regression suite and maintenance profile) and three ADAPTED JSON
envelopes encoding the full official Node24.20.0, pnpm11.25.0 and Sigstore4.1.1
license bytes. Base64 is a reversible representation; the decoded bytes retain
their original hashes, including CRLF and absent trailing newlines. The
envelopes are not claimed as VERBATIM source files. Node and pnpm bundled
notices retain their original terms; Sigstore retains Apache-2.0. The opaque
acquisition does not install or redistribute the tool bundles. Evidence:
reconstruction_evidence/GOVERNED_MAINTENANCE_ACQUISITION_V322.md.

## Library contents

The current post-V245 inventory contains 1356 materialization blocks: 1123 `AUTHORED`, 128 `ADAPTED` and 105 `VERBATIM`. V224 omnichannel ingress contributes four `AUTHORED` and six `ADAPTED` blocks; V225 candidate promotion contributes two `AUTHORED` and five `ADAPTED` blocks; V226 Meta lead reconciliation contributes six `AUTHORED` control/artifact files and two `ADAPTED` SDK-bound files; V227 Meta import contributes three `AUTHORED` fixtures and four `ADAPTED` integration files; V228 contributes two `ADAPTED` Responses API contract/test files without embedding OpenAI product source; V236 contributes seven `AUTHORED` connected-runtime/PostgreSQL blocks; V237 contributes eight `AUTHORED` contact-identity/PostgreSQL/E2E blocks; V238 contributes eleven `AUTHORED` outbound-fence/PostgreSQL/E2E blocks; V239 contributes seven `AUTHORED` PostgreSQL/executable/lock files and two `ADAPTED` Meta webhook boundary/tests governed by exact official sample revisions. V240 extends the existing execution validator with seven `AUTHORED` state/lock/protocol/checkpoint/test files; it embeds no GitHub, AWS, NASA, Google or xAI source and records exact identities and narrow method claims separately. V242 adds seven `AUTHORED` lock/profile/parser/test/operator blocks and one explicitly `ADAPTED` Lead client over TikTok's official generic transport and official v1.3 endpoint contracts; V243 adds three `AUTHORED` request/evidence/runner blocks, and V244 adds five `AUTHORED` Go import/test/command/guide blocks that reuse the existing provider-neutral store. V245 reclassifies exactly one Microsoft Durable Task runtime block from `AUTHORED` to `ADAPTED`: it preserves the pinned official 1.9.0 worker start sequence and adds deterministic event-loop cleanup; no Microsoft authorship is claimed for that correction. The TikTok blocks do not claim a provider-generated Lead SDK or authenticated webhook. The Google Ads reporting pack V2 adds two `AUTHORED` integrity artifacts: a complete CPython 3.14 Windows wheel lock and an SDK provenance receipt; the runtime itself remains the exact official Google wheel and dependencies, not embedded source presented as local code. None embeds Google, AWS, OWASP, PostgreSQL, Meta, OpenAI, NIST or TikTok product source. `FRANCHISE_ACCELERATOR.md` and `START_FRANCHISE.md` are governed root entrypoints, not materialization blocks. V174 retains five explicitly `ADAPTED` connected-carrier implementation/migration/test files and one `AUTHORED` derivation record governed narrowly by exact Microsoft BCApps and admitted Amazon adapter boundaries. Provider delivery never becomes automatic customer/business acceptance. No local Go/SQL/control code is represented as Microsoft-, Amazon-, Meta-, Google-, OpenAI-, NIST- or TikTok-authored product code.

Correction to the preceding chronology: the Microsoft Durable Task runtime reclassification and deterministic event-loop cleanup are V246, not V245.

V247 adds seven `AUTHORED` Kiota gate files. They are Elite acquisition, generation, frozen dependency and verification glue; no Microsoft source is embedded or relabeled. Microsoft Kiota 1.35.0 and its Windows x64 asset are acquired from the exact signed release under MIT. Therefore the current inventory is 1,363 blocks, not the older post-V245 total described in the historical paragraph above.

V249 adds four `AUTHORED` Go native-fuzz gate files. They are local fail-closed execution, profile, verification and operator glue governed by the official Go fuzzing documentation; no Google or Go project source is embedded, copied or relabeled. They require real project-owned `FuzzXxx` invariants and do not claim exhaustive security, DAST or production admission.

V252 adds four `AUTHORED` acquisition, source-lock, verification and operator files for `MICROSOFT-DEVSKIM-ADAPTED-SAST-GATE`. They embed neither Microsoft DevSkim nor SharpCompress source or binaries. The produced local runtime is explicitly `ADAPTED`: it acquires Microsoft DevSkim at verified commit `a452aa506f80928b611c4c7bfe306b8115821aae` under MIT, preserves Microsoft's LICENSE/NOTICE, and adds direct SharpCompress 0.48.0 references under MIT to replace the vulnerable transitive 0.40.0 resolution. Its receipt preserves both licenses, the exact adaptation, hashes and dependency evidence. The unchanged official DevSkim upstream remains separately rejected; no Microsoft authorship, official-byte identity or comprehensive-SAST claim is made for the adapted runtime.

V253 adds five `AUTHORED` lock, installer, scanner, verifier and operator files for `GITLAB-OPENGREP-SIGNED-SAST-GATE`. They embed no OpenGrep binary or GitLab rule bytes. Installation acquires the exact signed OpenGrep 1.29.0 Windows asset, verifies it with exact Cosign 3.1.3 and acquires three exact GitLab SAST Rules 2.9.3 files. Their LGPL-2.1-only, Apache-2.0 build-tool, MIT and LGPL-3.0-only licenses/notices remain external and are copied into the generated runtime. GitLab Enterprise and Commons-Clause families are verified and excluded. Running GitLab rules with OpenGrep is an explicitly authored composition, not GitLab's Semgrep analyzer and not attributed to GitLab/OpenGrep as a product runtime.

`Current-Provenance-Counts: AUTHORED=2231; ADAPTED=179; VERBATIM=168; TOTAL=2578`

OpenSSH portable `V_10_2_P1` at signed commit `d01efaa1c9ed84fd9011201dbc3c7cb0a82bcee3` supplies four public `sshsig` regression blobs transported losslessly as Base64 by `PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE`. Their decoded SHA-256 identities and the upstream `LICENCE` blob/hash are retained in the pack notice; no OpenSSH executable, source implementation or private key is redistributed. Windows supplies the executing `ssh-keygen`, whose Microsoft Authenticode status and exact host hash are recorded per release. Google OSV-Scanner 2.5.1 is used by exact locked executable hash to produce SCA and SPDX evidence.

Microsoft SBOM Tool 4.1.5 and its 2026-04-24 `main`, Sigstore Cosign 3.1.3, and in-toto-golang 0.11.0 were evaluated for V254 but are not redistributed or admitted as runtime: current exact dependency/SBOM scans contained affected packages. The former OpenGrep/GitLab reconstruction remains preserved, but its installer now fails closed at `REJECTED_VERIFIER_RUNTIME` until a clean official verifier satisfies the recorded reopen trigger.

Of the 585 materialization blocks present at the current 2026-08-28 V88 audit, 478 declare `AUTHORED`, eight declare `ADAPTED` and ninety-nine declare `VERBATIM`. The Microsoft pg_durable handoff pack contributes 26 exact PostgreSQL License files from release v0.2.6; `live_smoke_check.sh` is not embedded because its missing final LF cannot be represented byte-exactly by the Markdown fence contract. The Dapr outbox pack contributes thirteen exact Apache-2.0 files from signed Dapr Runtime/Docs commits and eight local admission/acquisition files; the rejected Dapr Go SDK is not redistributed. The Meta WhatsApp pack contains three exact official Meta code examples plus its recorded executable adaptation/license normalization. The Microsoft Durable document pack contains the exact official Azure Durable Task Python skill as one additional `VERBATIM` block and its MIT license with final-newline-only normalization as one additional `ADAPTED` block; its orchestration code remains local `AUTHORED` glue and is never presented as Microsoft product code. The Microsoft Azure Document Intelligence invoice pack contains the exact upstream sample sync and root MIT license as two `VERBATIM` blocks; its six approval/acquisition/test/guide files are local `AUTHORED` integration code, while the official JPEG fixture is hash-locked but not redistributed. The three Microsoft Azure Content Understanding sample packs contribute eight `VERBATIM` blocks in total: invoice sample+license, analyzer-copy sample+official test+license, and binary-document sample+official test+license; their nine README/locks/contracts are `AUTHORED` packaging and never presented as Microsoft product code. The Google Document AI process pack contains Google's exact standalone sample, Custom Document Extractor per-request schema sample, both live tests and Apache-2.0 license as five `VERBATIM` blocks; its eight local approval/acquisition/validation/test/guide files do not claim Google authorship, and its invoice/packing-list/output/license fixtures remain external hash-locked Google artifacts. The Google Document AI lifecycle pack adds sixteen exact Google source/tests plus the Apache license as seventeen `VERBATIM` blocks, the single decoded Cloud documentation code block as one `ADAPTED` block, and three local `AUTHORED` guide/lock/verifier blocks; the adapted classification records DOM extraction and never claims Git byte identity. The Amazon Textractor pack contains three AWS files as `VERBATIM`, four official test-support files as `ADAPTED` solely to add the final LF representable by the fence contract, and five local `AUTHORED` packaging/lock/test files; original and packaged hashes remain separate. The Transparency.dev Tessera pack contains five exact Apache-2.0 upstream files—license, authors, POSIX design and two CLI sources—as `VERBATIM`; its profile, lock, runners, tests and README are six local `AUTHORED` integration files and are never presented as upstream product code. Exact revisions, paths, raw/packaged hashes, changes and licenses are preserved in the packs and reconstruction evidence; official references are not presented as the production execution path. The Firebase and Google Cloud Storage adapters are authored integration code importing official Apache-2.0 Google modules, not copied Google source. The Azure Blob evidence adapter is authored integration code importing the exact official Microsoft MIT wheel; it embeds no Microsoft SDK source. The Mercado Libre adapter is authored against official HTTP documentation and deliberately excludes the archived official SDK. The MarkItDown local runtime contains authored wrapper/lock/tests that import the exact official Microsoft MIT wheel; no Microsoft source is embedded as an authored block. The document routing gate is authored deterministic validation informed by locked AWS/Microsoft/Google patterns; it copies no vendor classifier or pipeline. The secure local file gate contains authored orchestration/config/tests around exact Google Magika, Cisco Talos ClamAV and VirusTotal/Google YARA-X binaries; it embeds no upstream engine or detection rules. The strict field evaluation gate contains authored fail-closed orchestration importing the exact AWS Labs Stickler release; it embeds or modifies no AWS engine source. The Business Central pack contains only authored acquisition/approval orchestration; Microsoft product code and BcContainerHelper are acquired under their own MIT/product terms and are not embedded as local blocks. The AWS authenticated-approval, Cloudflare network-edge and Google grounded-text-extraction source profiles are local `AUTHORED` selection boundaries; downloaded upstream files retain their own licenses and are not embedded as authored product code. The Microsoft AVM SFTP pack contains the exact AVM MIT license, Local User module and WAF test as three `VERBATIM` blocks; its five composition/lock/verification files are local `AUTHORED` code and never presented as Microsoft product source. The secure email MIME core adds five local `AUTHORED` blocks and explicitly claims no AWS authorship. Other official archives, package artifacts and modules acquired later remain under their recorded upstream licenses and are not part of the authored block count.

`INSTALL_AGENT_BRIDGE.ps1` and the Codex/Claude Skills it emits are local `AUTHORED` orchestration based on documented OpenAI and Anthropic discovery/progressive-disclosure mechanisms. They embed no OpenAI or Anthropic product source and do not change the 1380 materialization-block provenance count.

Public availability is not treated as permission to copy. Any future `ADAPTED` or `VERBATIM` block must preserve its exact upstream revision, path, license and required notices under `PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md` and `markdown_system/PACK_CONTRACT.md`.

## Tooling and runtime dependencies

OpenGrep 1.29.0 is LGPL-2.1-only external software whose exact Windows x64 release asset is signed with Cosign. GitLab SAST Rules 2.9.3 is an external mixed-license package; only exact MIT and LGPL-3.0-only `eslint`, `gosec` and `nodejs_scan` outputs are selected. Cosign 3.1.3 is Apache-2.0 build tooling and is not redistributed in the generated runtime. `GITLAB_OPENGREP_SIGNED_SAST_GATE.md` embeds only authored glue, preserves exact license texts and excludes GitLab Enterprise and Commons-Clause rules. It is a conditioned SAST lane, not a security certification or byte-identical GitLab analyzer.

Microsoft DevSkim at commit `a452aa506f80928b611c4c7bfe306b8115821aae` is MIT-licensed external software with NOTICE. `MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md` embeds no DevSkim or SharpCompress source/binary; it acquires and verifies the signed commit tree, applies the declared two-project SharpCompress 0.48.0 adaptation, rebuilds, runs all 300 upstream tests, scans the resolved dependency graphs and publishes licenses plus provenance. SharpCompress 0.48.0 is MIT. This security-lint lane is conditioned and does not replace project-specific threat modeling, code review, DAST, fuzzing or offensive testing.

Microsoft Kiota 1.35.0 is MIT-licensed external software. `MICROSOFT_KIOTA_OPENAPI_CLIENT_GATE.md` embeds no Kiota source or binary: it acquires the exact signed-release Windows x64 asset, verifies archive/executable identity and invokes it through seven `AUTHORED` files. The official `ToDoApi.yaml` remains external and hash-locked as a generator fixture. Generated clients and their dependency licenses belong to the target project; Kiota does not confer production admission on auth, idempotency, reconciliation or provider behavior.

The implementation packs currently reference or materialize projects under their own licenses, including Go, PostgreSQL, pgx, go-oidc, Next.js, React, pnpm, jose, openid-client, PGlite, OpenTelemetry, Prometheus, CycloneDX tooling, Red Hat/Debezium core, Quarkus and examples, Microsoft Durable Task, Azure Content Understanding and Azure Storage Blob SDKs, Amazon `amzn-sp-api`, AWS SDK for Go v2 S3, AWS Powertools TypeScript batch/idempotency, AWS GenAI IDP Accelerator, NVIDIA NeMo Retriever, Google Merchant/Ads/Document AI clients, Meta Business/CAPI/WhatsApp examples, TikTok Business API SDK and container base images. Exact runtime/build pins and license expressions live in each pack and its lockfiles.

`GO_OMNICHANNEL_LEAD_INGRESS` embeds no Google, Mercado Libre, AWS, OWASP, CloudEvents or PostgreSQL source. Its seven `ADAPTED` and five `AUTHORED` files are workspace-owned implementations constrained by those official public contracts. Google and Mercado Libre documentation/code-sample licensing remains attached to the linked pages; `pgx`, Go and PostgreSQL retain their dependency licenses. No local classification overrides any upstream terms.

`GO_LEAD_CANDIDATE_PROMOTION` embeds no Google, AWS or PostgreSQL source. Its five `ADAPTED` and two `AUTHORED` files are workspace-owned code constrained by official webhook, at-least-once/idempotency and transaction contracts. It does not convert an upstream sample into a production claim and it never treats provider delivery as contact permission.

Microsoft `durabletask` 1.9.0, its two exact Human Interaction sample paths and the official Azure Durable Task Python skill are MIT; signed commit/archive, wheel, sdist, license, sample and skill bytes are fixed. The samples prove event/timer mechanics but not approver identity—the Azure HTTP sample is anonymous—while the in-memory backend is test-only and production DTS remains project-conditioned. Microsoft Agents Human Oversight commit `9eca38d…` is also MIT and hash-locked as an archived rejected reference: it contains Office 365 option-email code but no implemented Entra policy and its audited graph retains advisories. Azure.AI.ContentUnderstanding .NET 1.1.0 is MIT; its source tag/commit and author/repository-signed NuGet were verified, but the service remains project-conditioned. Microsoft `azure-storage-blob` 12.30.0 is MIT; wheel, sdist, source archive and package license are hash-locked, while its unsigned tag condition is preserved and no live Azure authority is inferred. AWS Powertools TypeScript 2.35.0 is MIT-0 with NOTICE; its exact source archive and batch/idempotency npm tarballs are hash-locked as external conditioned runtime components, never embedded or relicensed as Elite code. AWS GenAI IDP 0.6.5 is MIT-0 with NOTICE and is acquired as an external conditioned source, never relicensed as Elite code. AWS Transactional Outbox commit `23e0519…`, AWS RAPID 1.25.1 and AWS Assess Workbench commit `79a57b5…` are MIT-0; IBM Document Extraction Toolkit commit `d9144ec…` is Apache-2.0; Microsoft Expense Submission MCP y Approvals Box pertenecen al mismo archive MIT del commit `8c2cb6e…`. Their exact external archives, licenses and manifests/locks or additional artifact hashes where published are recorded only as rejected research references, and none of their code is embedded or composed. NVIDIA NeMo Retriever is Apache-2.0 with `THIRD_PARTY_LICENSES.md`; Elite fixes a signed unreleased commit and preserves both, but two dependency advisories keep it conditioned.

AWS SES Mail Manager Attachment Pipeline commit `79314a9…` is MIT-0 with NOTICE and remains an external exact source, not redistributed product code. It is rejected for immediate adoption: no tests/lock, vulnerable `aws-cdk-lib 2.251.0`, nonexistent npm CLI pin, paid service/add-on conditions and application-level retention/tenant/idempotency/KMS gaps. A clean separately resolved CDK graph and synth do not transfer production admission to the upstream sample.

AWS `serverless-mail` commit `70ac931…` and AWS Security IR email integration commit `08f11d4…` are MIT-0 external exact sources, not embedded product code. Both remain rejected: the first has no focal tests/lock and incomplete fail-closed MIME handling; the second passes 42 isolated unit tests but does not extract attachments, enforce receipt authentication verdicts or connect its declared DLQ. AWS `sample-bda-redaction` commit `2c89832…` publishes no LICENSE/NOTICE in its archive; it is not acquired, copied or treated as reusable despite public visibility.

AWS `sample-secure-transfer-family-code` commit `474ca07…` is an external MIT-0 exact source and is not embedded as product code. It remains rejected for immediate adoption: zero tests/locks/workflows, external GuardDuty plan, shared upload root/role, obsolete selected security policy, no object-version idempotency and a routing handler that returns normally on copy failure without retry/DLQ/reconciliation. Exact archive, license and template hashes are retained in the source lock and `AWS_TRANSFER_FAMILY_SFTP_CANDIDATE_2026-08-28_V1.md`.

AWS `file-transfer-sync-solution` commit `5ee79ee…` is an external MIT-0 exact source and is not embedded as product code. It remains rejected for immediate adoption: zero tests/locks/workflows, dependency installation during CDK synthesis, timestamp-only change detection, discarded asynchronous transfer IDs/status, absorbed errors and a first-copy flag written after dispatch even when the report is corrupt. Exact archive, license and focal source hashes are retained in the source lock and `AWS_FILE_TRANSFER_SYNC_CANDIDATE_2026-08-28_V1.md`.

AWS `sample-ai-receipt-processing-methods` commit `6eb72bf…` is an external MIT-0 exact source and is not embedded as product code. It remains rejected for immediate adoption: its PWA and infrastructure build, but the declared tests are absent, dated production SCA reports 14 root and 13 frontend vulnerabilities, status exposes pre-auth diagnostics, the demo disables versioning/PITR, POST omits checksum/content security, S3-event failures are swallowed and automatic persistence lacks CAS/field evidence/review. Exact archive, license, locks and focal hashes are retained in the source lock and `AWS_RECEIPT_MOBILE_INGESTION_CANDIDATE_2026-08-28_V1.md`.

AWS `amazon-s3-endedupe` commit `7209a02…` is an external MIT-0 exact source and is not embedded as product code. It remains rejected for immediate adoption: 17 unit tests pass, but measured coverage is 92% without importing `app.py`, current SCA reports 30 occurrences/23 unique IDs for each exact manifest, Python 3.9 is deprecated, variable-length sequencers are compared without left-padding, unlock/rollback do not fence by owner, locks never expire and the target lacks DLQ/retry/reconciliation. Exact archive, license, manifests and focal hashes are retained in the source lock and `AWS_S3_ENDEDUPE_CANDIDATE_2026-08-28_V1.md`.

AWS `amazon-sqs-best-practices-cdk` commit `54143c7…` is an external MIT-0 exact source and is not embedded as product code. It remains rejected for immediate adoption: application code is unchanged from unsigned 2023 source, has no tests/workflows/lock, pins vulnerable CDK 2.76.0 and Python 3.8, ignores `SendMessageBatch` partial failures, handles only `Records[0]`, writes a new random UUID on retry and omits `ReportBatchItemFailures`, conditional business effects and DLQ replay/reconciliation. Exact archive, license, focal hashes, synth/SCA and semantic harness are retained in the source lock and `AWS_SQS_BEST_PRACTICES_SAMPLE_2026-08-28_V1.md`.

AWS `ocr-with-aws-ai-services` commit `5341910…` and the Acme Bikes purchase-order validator under `aws-ai-intelligent-document-processing` commit `31d671e…` are external MIT-0 exact research sources and are not embedded or composed. The OCR runtime is rejected because its current graph has 31 unique vulnerability IDs, its declared Ruff gate reports 633 findings and its suite is not portable across archive/Windows symlink contexts. The purchase-order validator is rejected because it has no focal tests or lock and contains business-specific/non-strict validation behavior. Exact revision, archive or focal-file hashes and decisions are retained in `AWS_OCR_PURCHASE_ORDER_AND_AZURE_BINARY_REAUDIT_2026-08-28_V80.md`.

AWS Powertools for Lambda Python 3.34.0 is an MIT-0 conditioned runtime component. Exact GitHub archive, PyPI wheel/sdist and SLSA bundle hashes are fixed; 314 package files match tag↔sdist, 189 focal tests pass and the dated application runtime scan has zero known advisories. `AWS_POWERTOOLS_IDEMPOTENT_SQS_BATCH_COMPONENT.md` redistributes five exact upstream files under MIT-0 — LICENSE, documented idempotency+batch handler, payload and two separate SAM templates — without relicensing or claiming AWS published a merged deployment. It is not represented as a complete worker: tag/commit are unsigned, provenance says non-reproducible, Windows integral retains 11 failures, live Redis/AWS E2E are pending, build pip retains 7/6 advisories, DynamoDB update/delete lack owner fencing and S3 ordering/replay/reconciliation remain project gates. Evidence: `AWS_POWERTOOLS_PYTHON_IDEMPOTENCY_BATCH_2026-08-28_V1.md` and V75.

AWS Lambda Durable Execution SDK for Python 1.7.0 is Apache-2.0 with NOTICE. `AWS_LAMBDA_DURABLE_EXECUTION_COMPONENT.md` redistributes eight byte-verbatim files from verified commit `075b65a…`: LICENSE, NOTICE and three official example/test pairs for retry, `AT_MOST_ONCE_PER_RETRY` and replay-aware logging. The exact core wheel/sdist and runtime graph are hash-locked. The published testing wheel 1.2.1 is not redistributed or treated as compatible evidence because it differs from the selected tag source and lacks a symbol required by that tree. The component is conditioned; it does not claim exactly-once, S3 ordering, a business effect, IAM/IaC, DLQ/redrive/reconciliation or live AWS authorization.

Amazon Textractor 1.10.0 from `aws-samples/amazon-textract-textractor` is Apache-2.0 with NOTICE. `PYTHON_AWS_TEXTRACTOR_OFFICIAL_COMPONENT.md` redistributes LICENSE, NOTICE, the official Queries test and four support/fixture files from commit `8ea5f9a…`; final-newline-only adaptations preserve original bytes/hashes and packaged hashes separately. The official PyPI wheel is hash-locked but not embedded. Linux deterministic tests pass 69/16 skips, while two Windows portability failures and an upstream unmarked flaky fuzz remain explicit. The component does not claim ground-truth accuracy, universal reading order, automatic storage or live AWS authorization.

AWS DynamoDB Lock Client 1.5.0 is an external Apache-2.0 conditioned component and is not embedded or relicensed as Elite code. Exact commit archive and Maven JAR/sources/POM hashes are fixed; the three Maven artifacts validate under fingerprint `79D6089D38724F977A7DFD922B60A375F5B66ABF`, 18/18 main sources match the sources JAR and 112 upstream tests pass on Temurin JDK 8. The published dependency graph retains eleven runtime and four test advisory IDs, seven HIGH; PowerMock tests do not run on JDK 21, eight rebuilt class bytes differ from the signed JAR, no live DynamoDB test ran and its random UUID record version is not monotonic fencing for external effects. A separately labeled non-upstream SDK/Log4j update passed 112 tests and OSV 0 but is not an AWS lock-client release. Evidence: `AWS_DYNAMODB_LOCK_CLIENT_2026-08-28_V1.md`.

The V64 source lock records eight exact human-review surfaces without embedding their code separately: AWS human-in-loop patterns, AWS agentic insurance claims, AWS Nitro multi-approver and AWS scalable BDA+A2I are MIT-0; AWS timecards, Microsoft Azure agents escalation and Microsoft Durable Task JS are MIT; Google Document AI `hitl-custom-review` is Apache-2.0 inside the already locked Google samples archive. Exact commits, archive/license hashes and critical artifact hashes are in `OFFICIAL_UPSTREAM_ACQUISITION_CORE 0.4.49`. AWS Nitro is conditioned after 67 official tests; five other AWS/Microsoft samples, the Microsoft Durable Task JS human demos and Google HITL are rejected for incomplete contracts, findings or service deprecation. Downloaded files retain their upstream licenses and must never be relabeled as Elite-authored code.

Transparency.dev Tessera at signed commit `a8f33c56b80be808712c0537cb1b71f2fd5846ea` is Apache-2.0. Its official `AUTHORS` names Anthropic PBC, Google LLC and Internet Security Research Group. Google Trillian recommends Tessera for new transparency systems, but Elite attributes it to Transparency.dev/Tessera authors. The selected commit is a conditioned post-V1.0.4 snapshot; five OSV records remain in three modules not called by the POSIX CLI, and Linux/POSIX, key custody, independent checkpoints or witnesses, backups and access controls remain project gates.

The Amazon SP-API SDK source/wheel is Apache-2.0 with NOTICE; the verified Google Merchant wheel/source is Apache-2.0 and its 20-wheel Windows/Python graph additionally contains MIT, MPL-2.0, BSD-3-Clause and PSF-2.0 components. Meta Business 26.0.1, CAPI Parameter Builder 1.3.0 and WhatsApp API examples use a non-OSI platform license restricted to Facebook web services/APIs; their exact official LICENSE must be preserved. The CAPI wheel omits `LICENSE`; its 18-wheel graph also includes LGPL-2.1-only `pycountry` and the licenses recorded in that pack. TikTok Business API source at signed commit `f809c396520df2d7b201a9ccc5378d822b728ed3` is MIT; no release/tag/wheel was admitted, so use requires the exact approved source receipt plus its root license and critical-file hashes. Its four-wheel runtime graph retains MIT, Apache/BSD and MPL-2.0 components.

Google SafeValues 1.2.0 is an external Apache-2.0 npm runtime with zero runtime dependencies. Elite embeds no Google source: `TYPESCRIPT_GO_API_WEB_BRIDGE` adds two `AUTHORED` adapter/test files and imports the exact tarball fixed by npm integrity. Source commit/tree/archive, license, manifest and package hashes plus its conditioned scope are recorded in the V106 evidence. The Google `safety-web` alpha linter is not distributed or admitted.

Google CSP Evaluator 1.1.8 is an external Apache-2.0 npm development dependency with zero runtime dependencies. Elite embeds no CSP Evaluator source: the bridge imports the exact tarball fixed by registry integrity solely in tests. Source commit/tree/archive, package/license hashes, npm signature and disclaimer are recorded in V107 evidence. Its output is never treated as a universal XSS/security proof. Microsoft `@microsoft/eslint-plugin-sdl` 1.1.0 is not distributed or admitted because its only current peer baseline is ESLint 9, which reached end of life before this snapshot.

Grafana k6 2.2.0 is an external AGPL-3.0-only load-testing engine. Elite does not redistribute its source archive, Windows ZIP, SPDX document or executable; it records their exact official IDs, sizes and SHA-256 values and executes only a user-provided matching binary after explicit authorization. Elite-authored adapters bind argv, workload script, target, summaries and error budget. Neither a zero exit nor a threshold PASS alone is represented as production-capacity proof.

PostgreSQL 18.6 is external PostgreSQL-licensed software. Elite records exact PGDG mirror and official ftp source/checksum identities but redistributes none of them. The recovery adapter is local orchestration and never represents `pg_verifybackup` alone as proof of restoration or PITR.

OWASP ZAP 2.17.0 is external Apache-2.0 software. Elite records exact source, Core ZIP and SBOM identities but redistributes none of them. The local offensive-security adapter requires explicit target authorization, authenticated web/API plans, governed alerts and manual review; a zero-alert result is not represented as absence of vulnerabilities.

The latest verified dependency evidence is recorded in `reconstruction_evidence/DEPENDENCY_LICENSE_EVIDENCE_2026-08-24_V1.md`. It includes a real Go `go-licenses` report/check/save run and a pnpm production license report. The optional Windows sharp/libvips artifact reported `Apache-2.0 AND LGPL-3.0-or-later`; the example policy intentionally does not preapprove that redistribution condition.

## Distribution rule

Red Hat/Debezium core and Quarkus release 3.6.1.Final and the exact Debezium Examples outbox snapshot are Apache-2.0 external sources. Elite embeds only one authored source-selection profile; it does not embed or relicense Debezium code. Their unsigned Git identities, exact archive/license/artifact hashes, 48 focal tests, two compiled runtime JARs, example package result, OSV scope and missing Docker/PostgreSQL/Kafka integration are retained in `DEBEZIUM_RELATIONAL_OUTBOX_2026-08-27_V1.md`. The example remains `SAMPLE_ONLY`; outbox delivery is not exactly-once business execution or an immutable audit ledger.

Cloudflare Pingora 0.8.1 at commit `719ef6cd54e40b530127751bab6c1afc5ae815a8` is an external Apache-2.0 source. Elite embeds no Pingora source: it adds one authored selection profile and an exact lock record. The unsigned tag/commit, archive/license/artifact hashes, two POSIX symlinks, failed HTTP/2 test, 11 dependency advisories and failed release workflow remain explicit. Acquisition is fail-closed on Windows and Pingora is not composed or presented as production-ready.

Google LangExtract 1.6.0 at commit `62a25764745c0970f322b04bd89c844e5554bce1` is an external Apache-2.0 source. Elite embeds no LangExtract source: it adds one authored fail-closed selection profile and an exact lock record. Archive, license, package manifest, wheel/sdist and critical-file hashes are preserved in `GOOGLE_LANGEXTRACT_1_6_0_REAUDIT_2026-08-27_V1.md`. The release is rejected for automatic trusted persistence because its Gemini realtime and batch paths reproduced successful empty outputs for no-text/blocked responses; open PR code is not copied or treated as a fix.

Google Cloud Document AI Toolbox 0.17.3 at commit `3e29682282ac03d929fbfb93b5f94736fd420249` is external Apache-2.0 source and PyPI wheel/sdist. Elite embeds no Toolbox source or binary: it adds authored lock/acquisition/probe metadata and exact hashes. The source is `PINNED_CANDIDATE_CONDITIONED`, not automatically composable; Google labels it Alpha/experimental and its declared PyArrow range cannot select the official fix for CVE-2026-25087. Identity, tests, symlink and dependency evidence are preserved in `GOOGLE_DOCUMENT_AI_TOOLBOX_0_17_3_REAUDIT_2026-08-27_V1.md`.

Before distributing a generated application:

1. scan the exact artifact for every target OS and architecture;
2. apply the project's approved license policy;
3. include all required LICENSE, NOTICE, source or offer artifacts;
4. retain the resulting SBOM, provenance and hashed license evidence with the release;
5. block promotion when a license is unknown, denied or not reviewed.

`LicenseRef-Workspace-Owner` applies only to original workspace material. It never overrides a dependency or upstream license. This document is not legal advice.

## V320 Node official advisory adapter

Adds7 AUTHORED wrapper/bridge/test/profile/lock/operator files,1 ADAPTED
is-vulnerable.js with only getJson changed to verified local reads, and2 VERBATIM
ascii.js/LICENSE files from nodejs/is-my-node-vulnerable1.6.1 at
c37a56bad56e34fe5223ddd3cb223cc4158136ae under MIT. Original copyright and full
MIT terms remain in the materialized official/LICENSE. Local guard code is not
attributed to Node/OpenJS/npm. External semver7.8.5 is ISC; its official tarball,
53files and unchanged LICENSE are pinned but not embedded in this library.
External Node executable/LICENSE and JSON snapshots are also not redistributed
as pack source. The ten-file pack carries its own explicit provenance notice.
30tests/4CLI and10/10 rebuild qualify a compatible instance only; no release
or production monitoring admission. Evidence NODE_OFFICIAL_ADVISORY_GATE_V320.md.

## V321 tool artifact scope

The three active web tool consumers pin pnpm11.25.0; the library does not embed its455published files. Preserve the official MIT root license and all embedded notices.473fixed package-license declarations include unchanged LGPL-3.0+/MPL-2.0/CC-BY-3.0 components; pnpm MIT is not a relicensing of those dependencies. Local verification tool sigstore4.1.1 Apache-2.0 is external with its50package lock. No new product implementation-file count or provenance split:1447files remain1206AUTHORED/134ADAPTED/107VERBATIM. Evidence PNPM_BUNDLE_SECURITY_V321.md; profile/redistribution conditions remain.

## V334 authored artifact projection tooling

Four AUTHORED blocks add the selector, verifier/tests, operator guide and factual
hash/member policy. No pnpm implementation or native bytes are embedded in the
library. Projected output is ADAPTED_SELECTION_UNCHANGED_FILES with22original
notices retained; runtime/redistribution admission remains false. Original
counts1209/137/107/1453 are preserved in the V333 evidence history.

V335 adds2AUTHORED files for consumer recipe planning and tests; existing guide changes. No upstream code or runtime redistribution incorporated. Historical counts remain dated.

V337 adds1AUTHORED source profile and1ADAPTED reversible base64 envelope of the unchanged npm-lifecycle Artistic2.0 licence. Decoded notice remains byte-identical to the exact source/archive. No package implementation bytes incorporated into the library. Historical totals retained.

V338: no new third-party implementation files or licence envelopes enter the portable payload. Existing1216AUTHORED/138ADAPTED/107VERBATIM=1461 counts remain unchanged. A separate local32-text draft evidence collection records observed notices and licence gaps; it is not a cleared runtime distribution. See reconstruction_evidence/PNPM_NOTICE_COVERAGE_V338.md.

V339 adds documentary evidence only: seven fixed-source licences are kept in a separate local39-text draft, not incorporated as new implementation payloads. Counts1216AUTHORED/138ADAPTED/107VERBATIM=1461 remain unchanged. Noble full ciphers attribution, yallist scope wording and gyp/packaging nested notices are preserved; no licence or runtime promotion. See PNPM_NESTED_LICENSE_SOURCES_V339.md.

V340 documentary research keeps eight new version-bound notice/attribution texts in a separate local47-copy draft. No implementation payload or source roster changes; provenance1216AUTHORED/138ADAPTED/107VERBATIM=1461 remains. Yarn registry/tag linkage unresolved, Undici6file proof distinct from7bundle markers. See PNPM_YARN_UNDICI_NOTICES_V340.md.

V369 adds1AUTHORED source-inspection profile and2ADAPTED base64 envelopes of byte-identical official Apache-2.0 LICENSE texts (TRL1.12.0 and PEFT0.20.0). The two training source archives remain outside the portable library; their download and license correspondence do not admit the frameworks, models or dependency graph.

V371 adds3AUTHORED files: quarantined PyPI wheel transport, its negative/regression tests and an isolated acquisition profile. Downloaded wheels and inspected metadata are external evidence, not redistributed executable library code; no training dependency or model admission.

V372 adds22AUTHORED blocks: one source-inspection profile, eight refund telemetry/lifecycle/reference CLI files and thirteen finite Windows reference acceptance files. It embeds no Collector/Prometheus/vendor product source as authored code. Their exact official sources, locally patched dependency locks, test-only Windows adaptation and full license/notice evidence remain explicitly external runtime qualification in INTEGRATED_TELEMETRY_CONTROL_V372.md. No binary release, provider payment or production certification.

V373 adds17AUTHORED history-training orchestration/fixture/lock files; no copied HF trainer, wheel or private model is distributed. Read the pack notices and HISTORY_TRAINING_SDK_QUALIFICATION_V373.md for exact dependencies, tokenizers missing wheel LICENSE paired from its fixed source, explicit safetensors prerelease and retained paste maintenance advisory. Runtime/model/corpus licenses and native completeness remain target gates.

## V395 — QRCode notice carried by pnpm recipe tooling

PNPM-ARTIFACT-SELECTION-GATE0.4.0 retains6files. Its plan_install.py block is now
ADAPTED (AUTHORED control logic plus VERBATIM vendor notice text), licensed
LicenseRef-Workspace-Owner AND MIT. This changes the exact split by one file:
AUTHORED1307→1306, ADAPTED147→148, VERBATIM107 unchanged, total1561 unchanged.
No QRCode algorithm or pnpm executable bytes enter the portable library.

The embedded header retains Copyright(c)2009 Kazuhiko Arase, the MIT declaration,
trademark statement and Node modification note from gtanner/qrcode-terminal
commit90f66cf5c6b10bcb4358df96a9580f9eb383307b vendor/QRCode/index.js,
Git blob10eb8eb0a06aa50d0d3c508f886a7728f52e5d98. Full MIT permission text is
supplied from https://opensource.org/license/mit with the header's copyright.
The assembled supplemental notice does not replace the containing package's
Apache license. Keep all original notices and both generated supplements with
their relevant works; no whole pnpm runtime or redistribution admission.
Evidence: reconstruction_evidence/QRCODE_VENDOR_NOTICE_DELIVERY_V395.md.

V396: one AUTHORED semver license quarantine profile increases materializable files to1562. Existing AUTHORED transport/tests are extended; no third-party algorithm added. The existing ADAPTED pnpm planner retains exact original semver-utils1.1.4 package/LICENSE1839bytes, SHA25680b98c1b20edfc51abd2c802ee7a1d3c5561151d108a24f19c207b6935beaed8, from npm artifact fa6980458971864f25f2afc7d7f6632b015d8b767296d48421e6ff13221bd2b3. Original text offers MIT OR Apache-2.0; complete MIT option selected, entire original text and copyright2013AJONeal retained. Block remains ADAPTED / LicenseRef-Workspace-Owner AND MIT; original APACHEv2 declaration unchanged. No whole-pnpm license or redistribution promotion.


V397: retained-notice-catalog.json is one new ADAPTED envelope carrying47
unchanged license/attribution text copies (151253bytes before base64).
LicenseRef-Pnpm-Bundled-Licenses identifies the collection's individual retained
terms, not a new grant or blanket relicensing.22texts come from the exact pnpm
selected payload;25come from earlier fixed research evidence. Each retains its
origin/scope/hash and portable locator; GPL/LGPL/MPL/Artistic/CC/BlueOak/BSD/MIT/
Apache and other original texts remain in their own scopes. No source offer,
relinking fulfillment, complete pnpm license closure or upstream executable code
is claimed by transporting these texts. The authored loader/control code remains
in the existing ADAPTED planner alongside its earlier notice constants.

V398: next-path-source-evidence.json is a new ADAPTED evidence envelope under
LicenseRef-Workspace-Owner AND MPL-2.0. It carries unchanged base64 copies of
next-path1.0.0 index.js, package.json and the full MPL-2.0 LICENSE from Seth
Holladay's fixed official commit ce2c1386836339bc473b76c9888947899ead8d56.
Original authorship and MPL terms remain attached to those originals; the
authored evidence envelope does not relicense them. The planner delivers all
three as data in NEXT-PATH-MPL-SOURCE.md with source URLs and named bundling
transformations. Neither original package source nor its manifest is installed
or executed. This closes only the evidenced local source/license delivery,
not whole-pnpm redistribution, runtime-loader equivalence or other components'
source/relinking obligations.

V400 customer survey packs add23AUTHORED source/test/document files under LicenseRef-Workspace-Owner. No external implementation or dependency version is copied/admitted. Bain NPS and PostgreSQL locking documentation are semantic authorities only; existing dependency licenses remain unchanged.

V401 adds one AUTHORED contain_zip.py exact-byte transformation utility. It embeds no upstream source or exploit. Its separate ADAPTED output retains original pnpm notices and records removal of adm-zip; no global source/relinking or redistribution admission. reconstruction_evidence/PNPM_ZIP_CONTAINMENT_V401.md

V402: readiness validator0.7.0 adds three AUTHORED orchestration/template/test files (11 total). Library preparation and release scopes are explicit; no external business code is attributed by that validator. Historical product AUTHORED classifications remain until real derivation/admission is demonstrated.

V402 published delta: GO-BC-EXACT-AMOUNT-ADAPTER0.1.0 contains two genuinely ADAPTED narrow BCApps functions/files, three AUTHORED tests/derivation records and the exact Microsoft MIT license (VERBATIM); Commerce0.6.2 and Accounting0.1.1 callers retain their provenance. GO-OFFICIAL-PAYMENT-WEBHOOK-ADAPTERS0.3.0 uses dependency-pinned Stripe-Go86.3.0 and MercadoPago1.14.0; its new hosted-checkout, GET and normalization wrappers/tests are AUTHORED glue, not vendor code. Refund worker0.1.6 adds two local regression files and a declared provider-code compatibility mapping. See BC_EXACT_AMOUNT_ADAPTATION_V402.md, PAYMENT_SDK_RECONCILIATION_V402.md and REFUND_PROVIDER_ALIAS_COMPATIBILITY_V402.md in reconstruction_evidence.

V402 BC sales delta: GO-BC-SALES-CONTRACT-ADAPTER0.1.0 adds2ADAPTED price predicate/quote transfer files and4AUTHORED clock/test/document files. Source and field/precision deltas are explicit at fixed BCApps revision2eae56d704a1fd035d104f333602aea7091b7749. The Microsoft MIT notice is reused through its existing owner, never duplicated or attributed to caller glue. See BC_SALES_CONTRACT_ADAPTATION_V402.md.

V402 connected delta:69additional AUTHORED configuration/transport/persistence/UI/verification files;8existing owner packs retain their original source categories. Stripe-Go86.3.0 and MercadoPago SDK1.14.0 remain DEPENDENCY_PIN with exact MIT notices and unchanged dependency graph; direct/indirect Go metadata was normalized offline. No local domain algorithm or glue is attributed to those companies. BC narrow derivations retain their separate MIT source evidence. Current tally is173packs/1679files; the franchise selects77/915.

V402 handover browser delta:18additional AUTHORED integration/UI/test files, no new dependency or copied upstream implementation. Existing SDK/runtime notices and exact revisions are unchanged. Current library173packs/1697files; franchise77/933. The Microsoft Playwright name identifies the pinned dependency, not authorship of local scenarios.

V402 factory browser delta:10AUTHORED read projection/BFF/UI/test files, no new dependency or domain/writer change. Current library173/1707; franchise77/943. Existing operations graph remains locally authored and is not attributed to an upstream company.

V402 current composition:178packs/1797blocks; franchise82packs/1032uniquefiles. The FX pack contains the same BC MIT license as the prior exact-amount pack; its selection excludes that duplicate path, without removing the standalone license. Two FX functions are ADAPTED from microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749; exact rounding/profile and persistence are local glue, not an AL built-in equivalence claim. Meta Page SDK remains26.0.1@788f363d15b1269ab5efb7cd00fb5e3b133cd99b; WhatsApp's original source mapping remains fixed. Service identity uses existing coreos/go-oidc3.20.0, oauth2.36.0 and their existing jose dependency; three exact license files and a dedicated notice are emitted by the broker pack. New source maps/notices are included in their packs. No dependency version or corporate authorship is inferred from application wiring. See COMMUNICATIONS_RUNTIME_V402.md and EXACT_FX_SNAPSHOT_V402.md.

V402 FX journal binding: seven new AUTHORED glue files, no new dependency or source license; exact existing BC source/notice retained. FX0.2.0 and accounting0.1.2 preserve the original posting/reversal algorithm. Current178packs/1804blocks, franchise82/1039. See FX_JOURNAL_CONNECTION_V402.md.

V402 stored-value connected release: Odoo Community19.0 commit99edb6dd82b7b560930c00b03b694ba700785370. Twenty-four official complete source files plus LICENSE/COPYRIGHT copies are VERBATIM; engine and fixture expectations are ADAPTED LGPL-3.0-only. Source paths/SHA/Gitblob and source-replacement/notices ship with PYTHON_ODOO_STORED_VALUE_CALCULATOR0.1.0. The authored Python module glue is LGPL; Go/SQL/TS wrappers are AUTHORED workspace integration. No company authorship inferred. LGPL LICENSE43529bytes SHAabc09dad5f84a76e1b0279237053cae16c03228ab27d8d467677054c2bd17eeb is preserved without final LF. Full corresponding selected source, modifications and rights are described in docs/provenance/ODOO_STORED_VALUE_NOTICES.md. Counts180packs/1877blocks; franchise84/1112.

V402 warranty: BCApps2eae56d704a1fd035d104f333602aea7091b7749 MIT, complete ServiceItemLine.Table.al32b4fcf60b25d363d31ef6a1927b1f051eb4ef4ef48d9ca6b2ecc04844a854ff and license VERBATIM. Only date predicates ADAPTED; connected sold-policy/PG/HTTP/host/materializer is AUTHORED glue. Product notices docs/provenance/BC_WARRANTY_NOTICES.md preserve the distinction. No AL runtime or corporate workflow claim. Two new packs, one shared selected license;182packs/1911blocks,86/1145franchise.

V402 optional identifiers: three AUTHORED migration/regression blocks, no new third-party material; SERIAL_OPTIONAL_IDENTIFIERS_V402.md. Current182/1914, franchise86/1148.

V402 serial J2:18new AUTHORED blocks and5AUTHORED owner revisions; no added third-party material. Existing inventory provenance remains. SERIAL_SUPPLY_CONNECTED_RELEASE_V402.md/json.

V402 catalog J3:36new AUTHORED blocks and11owner revisions; no new third-party material or dependency. Existing source/SDK/license locks unchanged; CATALOG_CONNECTED_RELEASE_V402.md/json.

V402 training:26new AUTHORED glue blocks and9owner revisions; existing15public guides retain local content authorship. No new dependencies, upstream code or corporate attribution. TRAINING_CONNECTED_RELEASE_V402.md/json.

V402 catalog authoring:20new AUTHORED glue blocks and10owner revisions; original model/price writers retained, no new dependency/source attribution. CATALOG_ROLE_AUTHORING_RELEASE_V402.md/json.

V402 supply role:16new AUTHORED glue blocks,7owner revisions; original Operations/BindPlan retained, no new dependency/source attribution. SUPPLY_ROLE_RELEASE_V402.md/json.

V402 warranty role:14new AUTHORED projections/forms/fixtures,7owner revisions; original warranty/coverage/approval/FIFO retained. No new dependency or corporate attribution. WARRANTY_ROLE_RELEASE_V402.md/json.

V402 network role:21new AUTHORED glue blocks,6revisions. Original four Fulfillment SQL bodies preserved byhash; no new dependency/corporate attribution. NETWORK_ROLE_RELEASE_V402.md/json.

V402 helpCMS:23newAUTHORED glue blocks, existing kernel extraction and21shared guides. No new dependency/corporate attribution; original15guides retained. HELP_CMS_RELEASE_V402.md/json.

V402 role KPIs:19newAUTHORED glue/proof blocks,7existing outputs revised; no new dependency/upstream/company attribution. ROLE_METRICS_RELEASE_V402.md/json.

V402 private locale:20new AUTHORED glue/translation/proof blocks,45revised outputs. No new dependency or corporate attribution; existing source pins and notices unchanged. PRIVATE_LOCALE_RELEASE_V402.md/json.

V402 marketplace mutation:16new AUTHORED HTTP/SQL/host/fixture/proof blocks. No new external code or dependency; prior source/license/notices retained. MARKETPLACE_MUTATION_RELEASE_V402.md/json.

V402304:8AUTHORED initial-marketplace composition/proof files; no third-party code/dependency added. Official contract URL/SHA in MARKETPLACE_INITIAL_RELEASE_V402.json; all prior license obligations retained.

V402305:7AUTHORED existing-marketplace content bindings/proof files. No external code or dependency update. Eight original official HTTP contract pins and all prior license obligations retained.

V402306:24AUTHORED Go/Python/SQL/runtime/proof binding files around unchanged Google Merchant SDK1.8.0 and20exact wheels. No Google source reattribution. Original wheel license payloads and Apache-2.0/MIT/MPL-2.0/BSD/PSF obligations retained;1041installed files hash verified. MERCHANT_CONNECTED_RELEASE_V402.json.

V402307:17new AUTHORED schedule/SQL/host/proof files;9revised glue files, original Meta source and third-party license bytes unchanged. SCHEDULED_COMMUNICATIONS_RELEASE_V402.json.

V402308:16new AUTHORED campaign/source/SQL/host/proof files;7revised glue files, original Meta source and third-party license bytes unchanged. CAMPAIGN_CONNECTED_RELEASE_V402.json.

V402309:26new AUTHORED identity/session/profile/host/test/doc glue files;6revised. Original SDK/dependency/license bytes unchanged. IDENTITY_PORTAL_RELEASE_V402.json.

V402310:4new AUTHORED J5 fixture/contract files;3revised permission-routing/fixture files. Original SDK/dependency/license bytes unchanged. IDENTITY_J5_RELEASE_V402.json.

V402311:5new AUTHORED orchestration/catalogue/SBOM/test files and2revised wrapper/documentation files. Embedded279upstream notice texts remain VERBATIM under their own licenses, with archive/locator/SHA per entry. No local corporate authorship claim or altered pnpm redistribution. PNPM_LOCAL_RUNTIME_V402.json.

V402312: Sin nuevos bloques ni cambios de licencias:1950AUTHORED/159ADAPTED/141VERBATIM. Backend0.4.6 preserva floors oficiales admitidos; DevSkim gate0.1.1 corrige sólo glue de diagnóstico. Fuentes/licencias originales intactas. COMPOSITION_SECURITY_RELEASE_V402.md/json.

V402313: 205packs/2270blocks:1966AUTHORED/159ADAPTED/145VERBATIM. Cuatro textos SDK AWS/Smithy originales y mapa15módulos; MIT fixture original. No scanner effectiveness ni atribución de glue a proveedores. DOCUMENT_REFERENCE_RELEASE_V402.md/json.

V402314: Cinco bloques AUTHORED nuevos de gobernanza de conversación, sin dependencia o notice externo nuevo.205/2275;1971AUTHORED159ADAPTED145VERBATIM. AI_RUNTIME_GOVERNANCE_V402.md/json.

V402315: 205/2282:1978AUTHORED159ADAPTED145VERBATIM. Siete archivos AUTHORED nuevos; notices/dependencias oficiales intactos. Histórico opt-in retiene su README/notices/pins originales. AI_CONNECTED_REFERENCE_RELEASE_V402.md/json.

V402316: 205/2294:1990AUTHORED159ADAPTED145VERBATIM. Doce bloques nuevos AUTHORED;18archivos de owners admitidos seleccionados exactos. Metadata generada Next ADAPTED explícita; herramientas/dependencias sin cambio. Excepciones de notices npm del artefacto se resuelven antes de TEST07, no ocultas. LOCAL_REFERENCE_DELIVERY_V402.md/json.

V402317: 205/2300:1996AUTHORED159ADAPTED145VERBATIM. Seis bloques AUTHORED de ensamblaje/ensayo/guía;17archivos existentes seleccionados. Middleware oficial sin modificación; lock distingue ejecutable histórico de actual. Obligaciones de notices de telemetría conservadas; distribución final de ejecutables requiere cierre T2810. LOCAL_REFERENCE_OPERATIONS_V402.md/json.

V402319: PostHog Inc.2020–2026 MIT core, commit6fafbb9081bd15e79448af5650e02a4f9ea435cc. Two ADAPTED numerical/test blocks and three VERBATIM source/test/LICENSE; exact root license6d82d67dba42eb94ba10f1e986d2eec338c22fb7c5216c2c0ebdecd83d53a029 retained at third_party/posthog-nps/LICENSE. No ee code/runtime. Three AUTHORED glue blocks, no corporate attribution.205/2308 and116/1591. NPS_SOURCE_ADAPTATION_V402.md/json.

V402321 release tooling: strict public-input selector/metadata, case-correct home-path gate, fixture regressions and current audit reconciliation are AUTHORED maintenance glue. Four pack outer-document repairs change zero payload bytes. Pack provenance1999AUTHORED/161ADAPTED/148VERBATIM (2308blocks) remains unchanged; no enterprise attribution for tooling.

V402322 reference runtime notices:16VERBATIM original license/README/manifest files,3explicit ADAPTED notice renderings (no invented upstream file/owner/year),5AUTHORED metadata/collector/test/docs files. Full original notices/terms and exact source mappings materialize underlicenses/reference/.155unique text artifacts account for129runtime manifests; native/build-only exclusions remain explicit. REFERENCE_NOTICES_RELEASE_V402.md/json.

V402323 signed gate: one additional AUTHORED regression file, no upstream license/fixture changes, no Microsoft authorship assigned to local orchestration. Two synthetic fixture-key families remain outside product/source and are not release identities.

V402324: five AUTHORED local signed-build glue/test/template/docs blocks; no external licensing/source change, no corporate attribution. Local signing key excluded from redistribution.

V402325 ARCA: four fixed original legal texts VERBATIM (MIT/Apache-2.0),18new AUTHORED glue/fixture/lock/docs blocks. External runtime composition is ADAPTED: original signed Microsoft svcutil binaries plus nine signed official NuGet6.12.5 assemblies. This is not an official Microsoft tool release; local glue is never attributed to Microsoft/ARCA. Package notices and339runtime files bind exact SDK10.0.400 archive/current10.0.11 runtime. See ARCA_CONNECTED_INFRA_V402 and its provenance table.

V402326: two existing AUTHORED release orchestration/test blocks corrected for UTC epoch. Provenance counts unchanged2028AUTHORED164ADAPTED168VERBATIM; no new upstream attribution or dependency.

V402328: two existing AUTHORED packaging/test blocks preserve original source plus declared Next-generated metadata. Counts unchanged2028AUTHORED164ADAPTED168VERBATIM; no new corporate/source attribution.

V402329: three existing AUTHORED release/report-parser/test blocks corrected. Provenance counts unchanged2028AUTHORED164ADAPTED168VERBATIM; no dependency or upstream attribution change.

V402330 five existing AUTHORED glue/test blocks corrected.2028AUTHORED/164ADAPTED/168VERBATIM unchanged. SPDX unversioned records are local metadata, never attributed third-party code.

V402331:19notice metadata sources now bind their existing official receipts and exact payload hashes; no changed license, notice text, source mode or dependency.

V402 final source330/metadata331: actual signed artifact notice collectors and complete SPDX coverage verified.1653selected files/116owners:1473AUTHORED glue,116ADAPTED,64VERBATIM. Current per-file/pinned-source ledger is reconstruction_evidence/LIBRARY_REFERENCE_PROVENANCE_V402.json. No installer/private key redistribution or corporate authorship inferred.
