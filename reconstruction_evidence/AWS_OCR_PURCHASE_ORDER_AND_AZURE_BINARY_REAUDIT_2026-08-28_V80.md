# AWS OCR, purchase order and Azure binary re-audit — V80

Date: 2026-08-28  
Scope: official document-ingestion candidates from AWS and Microsoft; no live cloud execution.

## Decision

One official Microsoft component was admitted as a conditioned, immediately materializable pack. Two official AWS candidates were rejected for immediate adoption. Publication by a leading vendor was not treated as proof of production readiness.

## AWS OCR comparison sample

Official repository: <https://github.com/aws-samples/ocr-with-aws-ai-services>.

| Property | Exact value |
|---|---|
| commit | `5341910c3aee2da872f2966caed0fe0246b8a55c` |
| archive bytes | `4,127,057` |
| archive SHA-256 | `0e419c085a1957a37de791050a72e004b3d54b91f392ccf0311a4a5be1d2ac8f` |
| license | MIT-0; 947 bytes; SHA-256 `7cb750713252efd1d578837ba8785b61319109906f0c10b87adab7cf4badfc42` |
| inventory | 95 files; 55 Python; 26 under tests; 20 JSON; 3 PDF |
| sample bundles | 10/10 with `schema.json` and `truth.json` |

The code provides Textract, Bedrock and Bedrock Data Automation comparison, schema-shaped outputs, per-field comparison, BDA explainability geometry, run records, strict schema/truth pairing and deterministic synthetic fixtures. It contains no packing-list, purchase-order, proforma, bill-of-lading, commercial-invoice or customs schema. Its ten fixtures are graphics, handwriting, nutrition, sheets/schedules and three synthetic claim forms.

The exact archive on Python 3.12 produced 663 PASS/5 FAIL: four tests require a Git checkout and one requires Windows symlink privilege. A clean detached checkout of the same commit produced 667 PASS/1 FAIL; the remaining failure is WinError 1314 before the security assertion executes. Synthetic generation check passed for all three tracked bundles.

The README instructs `ruff check .` without a version/config pin. Current `ruff 0.16.5` reports 633 findings in 25 categories, 496 auto-fixable. Its `requirements.txt` fixes Gradio 5.31.0, Pillow 11.2.1 and NumPy 2.2.6 but leaves requests, boto3, pandas and PyMuPDF open. Current Python 3.12 resolution contains 59 dependencies; structured `pip-audit` found 45 vulnerability records/31 unique IDs in Gradio 5.31.0, Pillow 11.2.1 and Starlette 0.52.1, including one ID without a listed fix. The JSON audit is 93,159 bytes with SHA-256 `3b2f2209c23ae4b6e686dc724b7f0bcc64ae0a898c4ca9fb55059b375d1cfef6`.

Decision: full runtime `REJECTED_COMPONENT`. Its pure-looking helpers import `shared.config`, which imports the Gradio theme, so they were not falsely isolated from the rejected graph. Conditions are retained as `UP-FAIL-161` through `UP-FAIL-163`.

## AWS purchase-order validator

Official repository: <https://github.com/aws-samples/aws-ai-intelligent-document-processing>.

| Property | Exact value |
|---|---|
| commit | `31d671ef54c6694fea6431ee4102e3de85b982c1` |
| tree | 536 entries; `guidance/agentic-orchestration` 278 entries |
| root license | 927 bytes; SHA-256 `ef47b4ae2a1a8d38ef87b38dc5927967e7b472ffa3f7d125e89dbc8353a1e7af` |
| validator Python | 17,510 bytes; SHA-256 `59b9b6fdccb0836023eee49e4eaf21fff4b660f4f1f0862a13d5aaf4c3265ecf` |
| purchase-order schema | 5,154 bytes; SHA-256 `c9c18661e424858afb9feb31a1308d9e90f267435d07e0ac9e7e52aea06f71a5` |
| requirements | 55 bytes; SHA-256 `77727ba6d09f66fb94f742e6ed55d874c22ad5fcd005ceccac03fbef5734311f` |
| Dockerfile | 741 bytes; SHA-256 `9cbc3703b7c2c5227473a091368f627e103f5f194e8054eb27829e4f4bce9054` |

Python AST and JSON parsing passed, but the component has no focal tests or dependency lock. Its three dependencies use open lower bounds. The schema is explicitly `acmebikes.com`, fixes an `AB-XX-NNN` SKU pattern and business-specific limits. The handler logs the complete event, invokes JSON Schema without `FormatChecker`, skips price comparison when the supplied price is zero and accepts the first retailer from an ambiguous partial-name query. It is a useful official scenario, not a universal enterprise purchase-order validator.

Decision: `REJECTED_COMPONENT`, retained as `UP-FAIL-164`; no code was copied into an implementation pack.

## Microsoft Azure Content Understanding binary sample

Official repository: <https://github.com/Azure/azure-sdk-for-python/tree/129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb/sdk/contentunderstanding/azure-ai-contentunderstanding>. Official package: <https://pypi.org/project/azure-ai-contentunderstanding/1.1.0/>.

| Artifact | Exact evidence |
|---|---|
| PyPI sdist | 230,330 bytes; SHA-256 `00f39c7cf2ba50c8d4586315f4295a45fad0881daaff07d13d33aea2bafc7ec8` |
| `sample_analyze_binary.py` | 6,979 bytes; SHA-256 `dce0d3c2684bb0d5bca01f1015e601d4f15407b8272c72b773e2cfb822588893`; commit=sdist |
| official test | 20,749 bytes; SHA-256 `88f8cb38d39a351639c835f246da5c443d5c50cc068a363032ab43ee229b276f`; commit=sdist |
| Microsoft license | SHA-256 `7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744` |

`MICROSOFT_AZURE_CONTENT_UNDERSTANDING_OFFICIAL_BINARY_DOCUMENT_SAMPLE.md` materializes six files. The three upstream files are `VERBATIM`; README, source lock and offline contracts are explicitly `AUTHORED`. Materialization 6/6, three official hashes, five offline contracts and three Python syntax checks passed. The exact sample performs three `begin_analyze_binary` calls with `prebuilt-documentSearch`, including page ranges `3-` and `1-3,5,9-`.

This component preserves local document bytes, markdown, layout, pages and tables. It does not claim procurement-field extraction. Projects must select `prebuilt-procurement`, `prebuilt-purchaseOrder`, another documented analyzer or a custom analyzer through their decision record, then prove field behavior against their corpus.

## Global gates

Before adding this evidence file:

```text
VERIFY_LIBRARY_PASS
packs=59 materialized_files=548 markdown_files=390
profile=AZURE_DOCUMENT_RUNTIME_PACK_PLAN.md implementation_files=85
```

```text
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=59 upstream_sources=110 official_azure_cu_samples=1 official_azure_cu_binary_samples=1 official_azure_cu_copy_samples=1
```

The failure memory contains 895 local lessons and 164 upstream conditions: 1,059 unique IDs and zero duplicates. Three additional local lessons retain the policy-safe cleanup regression; all seven audit temporaries were removed. Live provider accuracy, private project corpus, cloud access, costs and productive persistence remain deliberately unclaimed.
