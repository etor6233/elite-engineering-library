# Official Document SDK Source Distributions — 2026-08-27 V1

## Scope and decision

This evidence adds exact source distributions, not locally invented extraction logic, for stable SDK releases already admitted as integration candidates. `OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE` advances from 0.1.1/seven artifacts to 0.2.0/ten artifacts.

Decision: the Microsoft Azure Document Intelligence 1.0.2, Microsoft Azure Content Understanding 1.1.0, Google Cloud Document AI 3.15.0 and OpenAI 3.3.1 wheel/source pairs are `INTEGRATION_ONLY`. They make official implementation source and samples immediately acquirable after hash-linked approval. They do not authorize credentials, provider cost, automatic field persistence or a claim of correctness on project documents. Google Document AI Toolbox 0.17.3 remains a separate `PINNED_CANDIDATE` blocked by Alpha status and `UP-FAIL-122`.

## Official artifacts

PyPI JSON and immutable `files.pythonhosted.org` objects supplied these identities:

| ID | Bytes | SHA-256 | Official publication |
|---|---:|---|---|
| `azure-ai-documentintelligence-1.0.2-sdist` | 170,940 | `4d75a2513f2839365ebabc0e0e1772f5601b3a8c9a71e75da12440da13b63484` | Microsoft/PyPI, 2025-03-27 UTC |
| `google-cloud-documentai-3.15.0-sdist` | 362,736 | `d50b69a8a62aaf803b0d926f93b6579df440f5f779f45fec2c5972768a1cdbbd` | Google/PyPI, 2026-06-03 UTC |
| `openai-3.3.1-sdist` | 1,282,113 | `6f22807de1a976c932cecda620e8172a8c3fdbaeed29c7f21564e0c2410edf56` | OpenAI/PyPI, 2026-08-19 UTC |

The already locked Azure Content Understanding 1.1.0 sdist remains 230,330 bytes with SHA-256 `00f39c7cf2ba50c8d4586315f4295a45fad0881daaff07d13d33aea2bafc7ec8`.

The artifact runner acquired the three new sdists together from an empty destination and emitted:

```text
DOCUMENT_SDK_ACQUISITION_PASS artifacts=3
```

Receipt, filenames, sizes and hashes matched the ten-artifact lock exactly. Each new ID also passed `ValidateOnly`; the pack regression remained:

```text
DOCUMENT_SDK_ACQUISITION_TEST_PASS negatives=4 offline_verified=1
DOCUMENT_SDK_LOCK_VALID artifacts=10 selected=1
```

## Source contents and compilation

The tar archives were read without path rewriting. None contained an absolute path, `..`, hard link or symlink. Every Python source member was compiled directly from archive bytes:

| Distribution | Entries | `.py` | Official samples | Official tests | Compile errors |
|---|---:|---:|---:|---:|---:|
| Azure Document Intelligence 1.0.2 | 131 | 103 | 56 | 20 | 0 |
| Google Cloud Document AI 3.15.0 | 109 | 71 | 0 | 8 | 0 |
| OpenAI 3.3.1 | 1,821 | 1,760 | 0 | 189 | 0 |
| **Total** | **2,061** | **1,934** | **56** | **217** | **0** |

Azure Document Intelligence includes four upstream invoice paths:

```text
samples/sample_analyze_invoices.py
samples/sample_analyze_invoices_from_bytes_source.py
samples/async_samples/sample_analyze_invoices_async.py
samples/async_samples/sample_analyze_invoices_from_bytes_source_async.py
```

Those Microsoft samples use `prebuilt-invoice`, expose invoice fields and confidence values and demonstrate URL/stream/bytes variants. They remain samples: their prints and API key environment variables are not a persistence boundary and are never relabeled as Elite production code.

## Build and import proof

The exact Azure Document Intelligence sdist was installed into a new Python 3.14 environment. The build completed, `pip check` returned no broken requirements and the installed metadata/symbol probe returned:

```text
AZURE_DOCUMENTINTELLIGENCE_SDIST_PROBE_PASS version=1.0.2 symbols=DocumentIntelligenceClient,AnalyzeDocumentRequest,AnalyzeResult
```

No Azure endpoint, key, billed request or private document was used. A live extraction gate remains blocked until the user supplies the account, approved document classes, schema, representative corpus and ground truth.

## Bootstrap and vulnerability correction

The initial Python 3.14 venv inherited pip 26.0.1. Google OSV-Scanner 2.5.1 exact —Windows binary SHA-256 `25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6`— found four advisory groups. pip 26.1.2 removed three but retained `CVE-2026-13346`/`PYSEC-2026-3721`; the advisory range fixes it in 26.2.

The official pip 26.2 wheel was then fixed and verified:

```text
bytes: 1,816,475
sha256: 931c303696af6fa3417112103b1cad26890e5a07eccb5b99783700e33f2b8aad
```

After installation, separate scans returned zero known findings for:

- the nine-package Azure Document Intelligence runtime graph; report SHA-256 `54345d5194d238b6017d0f3181f392109ce35023cc6b4bdfa8c134943dcb578e`;
- the pip 26.2 toolchain entry; report SHA-256 `54345d5194d238b6017d0f3181f392109ce35023cc6b4bdfa8c134943dcb578e`.

The identical hash is expected because both empty OSV JSON reports have identical canonical bytes. This dated query is not a guarantee against unknown or future vulnerabilities. `UP-FAIL-123` preserves the failed 26.0.1/26.1.2 bootstrap versions so an agent cannot regress silently.

## Archive portability finding

The separate `Azure-Samples/document-intelligence-code-samples` commit `8ac994c2712f3104d00ca4224b0ded6fd1ecbf70` remains exact head, signed/verified and MIT. Its ZIP is 57,555,887 bytes with SHA-256 `018915d5ef37cba5fa3421a8716eff293127d5bc7492cafbf8cf141e6484f086`, but contains two filenames with literal `*` under tax schemas. Windows `Expand-Archive` cannot materialize them and returned a misleading zero exit after errors. Therefore V53 uses the portable official PyPI sdist for immediate SDK source/samples and preserves the repository ZIP only under the existing source acquisition condition. No upstream filename was renamed or silently dropped.

## Failure learning

`LIB-FAIL-481` through `LIB-FAIL-488` retain the failed extraction, array binding, exit-code wrapper, multiline probe, OSV filename/argument and inventory/date-query attempts. The corrected paths prove archive reading, three-artifact acquisition, source compilation, SDK build/import, fixed bootstrap and exact plan alignment before those incidents can close.

## Official authorities

- <https://pypi.org/project/azure-ai-documentintelligence/1.0.2/>
- <https://pypi.org/project/azure-ai-contentunderstanding/1.1.0/>
- <https://pypi.org/project/google-cloud-documentai/3.15.0/>
- <https://pypi.org/project/openai/3.3.1/>
- <https://github.com/Azure-Samples/document-intelligence-code-samples/tree/8ac994c2712f3104d00ca4224b0ded6fd1ecbf70>
- <https://github.com/google/osv-scanner/releases/tag/v2.5.1>
- <https://pypi.org/project/pip/26.2/>
