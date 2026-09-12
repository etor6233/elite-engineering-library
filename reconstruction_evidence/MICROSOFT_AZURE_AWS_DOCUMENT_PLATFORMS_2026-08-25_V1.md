# Microsoft Azure and AWS document platforms evidence V1

Verified on 2026-08-25. This is evidence for exact official source snapshots, not a claim that either sample is a finished REVESTEX implementation or universally accurate on unknown documents.

## Microsoft Azure Content Understanding Python

| Property | Exact value |
|---|---|
| repository | `Azure-Samples/azure-ai-content-understanding-python` |
| commit | `fbc18880179397729abfa9506bda865d009dc711` |
| archive bytes | `398347` |
| archive SHA-256 | `c38a543b437fcf8c76bb0c3574e8c5c4e7e70ba0eea604ce03ad8227c41f641e` |
| license | MIT; `LICENSE` SHA-256 `d494b116cb354b153d2291ce672a4a67ffda779fa4f385ce6ac62b7eb2e5c934` |
| inventory | 153 files; 18 Python; 8 files under a test path |

The exact snapshot targets the GA Content Understanding API and contains Bicep infrastructure, Python REST helpers, analyzer templates, classification/splitting, field extraction, trained-analyzer examples, migration tooling from Document Intelligence and official input/label/result fixtures. Its migration suite ran from the exact archive in an isolated Python 3.14 environment:

```text
70 passed, 1 skipped in 0.74s
```

The repository `requirements.txt` is not a complete frozen environment for that suite. Initial collection failed because `rich` was absent; after adding it, collection failed because `python-dateutil` was absent. Installing those two undeclared requirements allowed the result above. Therefore this source is `SAMPLE_ONLY / CONDITIONED`, not `REBUILD_VERIFIED`. Live analysis also requires an Azure AI resource, supported region, identity and service terms.

The official production Python SDK was separately fixed as `azure-ai-contentunderstanding==1.1.0`: wheel 101,987 bytes, SHA-256 `d1d6bdeffe02f5c8cc5309f0e120a1338f51af04638605f9617b7b501e2d1dd6`, MIT `License-Expression`, Python `>=3.9`. The exact wheel installed in an isolated Python 3.14 environment and imported `ContentUnderstandingClient`, `AnalysisInput` and `ContentAnalyzer`; 70 public model symbols were observed. Microsoft documents this GA SDK for API `2025-11-01` and production use. The exact artifact and install URL are in `OFFICIAL_DOCUMENT_SDK_ARTIFACT_LOCK.md`.

## Microsoft configurable data-extraction solution

The official `Azure-Samples/data-extraction-using-azure-content-understanding` snapshot is fixed at commit `0461a5549104ca769b8ec082c05894468997f300`: ZIP 808,247 bytes, SHA-256 `23fb059d172cc9ebbe76cebc5461cf9150222fbc91830f35c0519ddaf2b7f7cc`, MIT `LICENSE` SHA-256 `10d8ad187611652f75ab44ba79ad4d5c118f94746b03266a4b102fc38a2e35c5`. It contains 191 files: 87 Python, 55 Terraform and 35 files under test paths. The implementation includes Azure Functions, JSON extraction schemas, document classification/routing, Content Understanding extraction, Cosmos persistence, citations/source geometry, configuration hashes, Key Vault, health checks, OpenTelemetry/Application Insights and reproducible infrastructure definitions.

The exact dependency graph does not install on Python 3.14 because its Pydantic/PyO3 generation supports at most 3.13; the repository explicitly requires Python 3.12. On Python 3.12.13 its requirements installed, but test collection exposed an omitted development dependency, `parameterized`. After installing that missing package, the complete official suite passed:

```text
166 passed, 19 warnings in 5.43s
```

The source is admitted as `SAMPLE_ONLY / CONDITIONED`, not a finished production product: there is no dependency lock, the README explicitly requires production security/monitoring/error-handling adaptation, and live deployment requires Azure subscription, Content Understanding, OpenAI, Functions, Cosmos DB, Key Vault, Storage, Terraform and Azure CLI. The separately located Microsoft MCP sample was not admitted because its repository exposes no explicit reusable license.

## AWS sample document processing platform

| Property | Exact value |
|---|---|
| repository | `aws-samples/sample-document-processing` |
| commit | `9ca2eb1bf7b94f25c33926222817e0cd2d7748b4` |
| archive bytes | `7757260` |
| archive SHA-256 | `a0df347babc50b9cdee6168cd716839d6453923832322f126446482d5b8a810f` |
| license | MIT-0; `LICENSE` SHA-256 `7cb750713252efd1d578837ba8785b61319109906f0c10b87adab7cf4badfc42` |
| notice | `NOTICE` SHA-256 `d4290ed64c2edd0fce1d84e3f9dfb2881240fe534def76b8cd29ed6af683e287` |
| inventory | 148 files; 35 Python; 37 TypeScript; 2 files under a test path |

The official sample provides a React/Next.js review UI, REST and WebSocket backend, S3 document storage, DynamoDB state, Step Functions orchestration, Bedrock/AgentCore extraction into a target JSON schema, GuardDuty malware scanning, VPC, LLM gateway and CDK deployment stacks. Python source compilation passed. With fresh dependency resolution from the exact package manifests, the backend infrastructure, LLM-gateway infrastructure, Next.js production frontend, UI infrastructure and VPC infrastructure builds passed. The workflow infrastructure build failed because the repository has no lockfiles and resolved incompatible current Node type definitions against TypeScript 5.4.5 (`IteratorObject`, `AsyncIteratorObject` and `BuiltinIteratorReturn` errors).

This platform is therefore `SAMPLE_ONLY / CONDITIONED`. It is immediately acquirable and materially useful when AWS is explicitly selected, but it is not the portable Go/PostgreSQL base, and no deployment/runtime claim is made without AWS account, region, Bedrock model access, AgentCore, IAM, CDK bootstrap and a corrected frozen dependency graph.

## Admission decision

Both snapshots are admitted to the exact upstream lock because they add non-duplicative official code from Microsoft and AWS: Azure adds the current GA classifier/analyzer/training path; AWS adds a recent end-to-end upload, review, orchestration, security and deployment platform. Their upstream defects and cloud prerequisites remain blocking conditions. Neither source authorizes a claim of 100% field accuracy for proformas, packing lists, purchase orders, bills of lading, customs documents or a private supplier layout without project corpus and provider evaluation.

## AWS agentic IDP evaluation framework

| Property | Exact value |
|---|---|
| repository | `aws-samples/sample-agentic-idp-evaluation-framework` |
| commit | `8aaa758db972c951ce5654622f516598cc03a595` |
| archive bytes | `1644195` |
| archive SHA-256 | `386f3d12b658e0c54e9d14023cb7cdef36bd7b3798aca02829b5e18a36e6bd61` |
| license | MIT-0; `LICENSE` SHA-256 `5d0ed32d916046277a19e21df9cabc50e88de80e62bfff7f4df38b89800355ef` |
| root dependency lock | `package-lock.json` SHA-256 `fa1ad1c5c12c6e10291c012b5da895b44b958ab258515b6b032b2f0d7c3dcc23` |
| CDK dependency lock | `infrastructure-cdk/package-lock.json` SHA-256 `a8adabef5d2983225b6445e2fb2707512a33fae47afb693c6309f6286ba10335` |
| inventory | 344 files; 208 TypeScript; 15 Terraform; 41 files under test paths |

Frozen `npm ci` succeeded. The exact source generated 33 processing skills and 100 model-output limits, then built shared, backend and the Vite production frontend. Its independent CDK tree also installed frozen and compiled. The backend suite produced:

```text
35 test files passed, 3 failed, 2 skipped
356 tests passed, 7 failed, 12 skipped
```

The seven failures are material: three audio-routing expectations and four PII/Guardrails routing expectations fail against current source behavior. Installation also reports the upstream `multer@1.4.5-lts.2` security deprecation warning. The repository is admitted as `SAMPLE_ONLY / CONDITIONED`: it provides the strongest official evaluation harness found in this pass, but PII/audio decisions and dependency security must be corrected and revalidated before adoption. No live AWS evaluation or deployment was possible without account, model access, IAM and project documents.
