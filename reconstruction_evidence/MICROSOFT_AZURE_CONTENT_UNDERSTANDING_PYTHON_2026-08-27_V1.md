# Microsoft Azure Content Understanding Python 1.1.0 — source, artifact and runtime evidence

Date: 2026-08-27  
Decision: `PINNED_CANDIDATE / CONDITIONED`  
Authority boundary: the SDK, models, poller and service calls are Microsoft code under MIT. Elite adds only visibly `AUTHORED` acquisition, profile/security chaining and evidence orchestration; it does not relabel that glue as Microsoft code and does not implement OCR/extraction logic.

## Exact Microsoft source

- repository: `Azure/azure-sdk-for-python`;
- tag: `azure-ai-contentunderstanding_1.1.0`;
- tag target and commit: `129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb`;
- GitHub commit verification: `verified=true`, `reason=valid`, committer date `2026-04-20T21:43:08Z`;
- archive: `191429136` bytes, SHA-256 `218d4d555fb8b565ffb220ff00d8ae80dd51bf7c16bc626dfaf582e4ab3284a8`;
- root `LICENSE`: `1074` bytes, SHA-256 `7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744`;
- root `NOTICE.txt`: `614` bytes, SHA-256 `8783fc47f3420556eb8b53f746441ecb36b37f118ebec73d630e3468bde7f876`;
- package path: `sdk/contentunderstanding/azure-ai-contentunderstanding`;
- package `LICENSE`: SHA-256 `fd532481d828e13a0b13ccb598e02338a3617740675a862ee6bdc1541b68e93d`;
- package `pyproject.toml`: SHA-256 `2b74694ad1b5bdb001c81e9ace393c4a7e9b34a1d70211257ad6d00d033b3e08`;
- package `assets.json`: SHA-256 `97968ace7481e6e92924560cb40a3038d56c9932bf590fd5abeaa7074ec077c2`;
- service `ci.yml`: SHA-256 `8eb1c40762b2886cc95821e179c25b0b9c1915fdc0251c88b1b7fb0c1ec08d39`;
- service `tests.yml`: SHA-256 `fe8ac8bc282c8f168b601cdedb7d66cd52a236a018b72ae9431d19be5cd0da62`.

Primary sources: <https://github.com/Azure/azure-sdk-for-python/tree/129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb/sdk/contentunderstanding/azure-ai-contentunderstanding>, <https://github.com/Azure/azure-sdk-for-python/releases/tag/azure-ai-contentunderstanding_1.1.0> and <https://pypi.org/project/azure-ai-contentunderstanding/1.1.0/>.

## Exact PyPI artifacts and source relationship

PyPI reports version `1.1.0`, `Requires-Python >=3.9` and classifier `Development Status :: 5 - Production/Stable`.

| Artifact | Bytes | SHA-256 |
|---|---:|---|
| `azure_ai_contentunderstanding-1.1.0-py3-none-any.whl` | 101987 | `d1d6bdeffe02f5c8cc5309f0e120a1338f51af04638605f9617b7b501e2d1dd6` |
| `azure_ai_contentunderstanding-1.1.0.tar.gz` | 230330 | `00f39c7cf2ba50c8d4586315f4295a45fad0881daaff07d13d33aea2bafc7ec8` |

PyPI did not expose an attestation for these files. That absence is not hidden. The independent byte comparison linked the sdist to the signed tag: all `114/114` sdist files that correspond to repository package files were byte-identical. The other seven sdist entries are generated packaging metadata (`PKG-INFO`, `setup.cfg` and egg-info). This is strong identity evidence, not a substitute for a missing PyPI attestation.

## Published CI and tests

The signed source contains 43 test Python files plus fixtures/sample data. Microsoft `ci.yml` uses the Azure SDK client archetype with `TestProxy: true`; `tests.yml` sets a 200-minute timeout, public cloud/eastus and playback variables. `assets.json` points to `Azure/azure-sdk-assets` tag `python/contentunderstanding/azure-ai-contentunderstanding_83f2e32128`.

On Windows x86-64 with CPython `3.12.13`, uv `0.11.30` and the exact wheel:

```text
183 passed in 1.29s
IMPORT_PASS 1.1.0 1.41.0
All installed packages are compatible
```

The 183 tests are the two official unit files that do not require service playback after being copied byte-for-byte outside the package `conftest.py`: model exports/value behavior and analyze poller operation/usage behavior. They were first run under pytest 9.0.2, OSV exposed its known tmpdir issue, then repeated PASS under fixed pytest 9.0.3. No result from the vulnerable test pin is used as the final environment claim.

The complete 268-test collection was also attempted from the signed archive with the exact monorepo `eng/tools/azure-sdk-tools` and `sdk/core/azure-core`. All 268 aborted at session setup before a test ran because `devtools_testutils.ascend_to_root` requires a `.git` directory to start test-proxy. Elite does not create a Git repository merely to manufacture a PASS. This retained upstream condition is `UP-FAIL-074`; a project may close it only with the official checkout/test-proxy assets or authorized live Azure resources.

## Product dependency graph and known-vulnerability query

The exact wheel resolved nine product packages on CPython 3.12.13/Windows x86-64:

```text
azure-ai-contentunderstanding==1.1.0
azure-core==1.41.0
certifi==2026.7.22
charset-normalizer==3.5.1
idna==3.19
isodate==0.7.2
requests==2.34.2
typing-extensions==4.16.0
urllib3==2.7.0
```

The resolved list SHA-256 is `b0e1ceab557d8ff8519eb6bae2e9eb39330683ae263769e7c99ac487f77f9c15`. Google OSV-Scanner `2.5.1`, official Windows binary SHA-256 `25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6`, returned exit `0` and zero vulnerability records for those nine package-versions. Report SHA-256: `54345d5194d238b6017d0f3181f392109ce35023cc6b4bdfa8c134943dcb578e`. This dated query does not prove absence of unknown, future or unreachable vulnerabilities.

## Elite runtime correction and proof

`AZURE-CONTENT-UNDERSTANDING-DOCUMENT-RUNTIME 0.1.2` remains authored glue around Microsoft's `ContentUnderstandingClient.begin_analyze_binary`. It now rejects the two bypasses found in V0.1.0:

1. callers choose a document class, not an arbitrary analyzer ID; only a unique class with `decision=REQUIRED` can supply the analyzer;
2. every analysis requires an `elite-secure-local-file-receipt/v1` with `decision=ADMITTED`, `business_storage_authorized=false` and exact SHA-256/bytes/MIME equality with the local input.

The analysis receipt V2 hash-links input, security receipt, document profile and complete provider result; rejects provider analyzer/API divergence; writes atomically; and still sets `automatic_storage_authorized=false`. Five wrapper tests pass, including blocked/unknown classes, rejected/tampered security receipt, wrong MIME, preview API response, analyzer mismatch, empty content, occupied output and unsupported extension. These gates prevent bypass; they do not assert document-field accuracy.

The pack now also materializes `requirements-windows-py312.lock`: nine exact package versions and one admitted wheel SHA-256 per distribution for CPython 3.12/Windows x86-64. The global executable verifier creates a new venv, installs it with `pip --require-hashes`, runs `pip check` and executes the five tests when `-AllowNetwork` and the exact target runtime are supplied. The dated run passed from an empty venv. Offline audit validates the lock and emits an explicit `AZURE_CONTENT_UNDERSTANDING_RUNTIME_GATES_SKIPPED`; it never turns absence of network into an execution PASS.

`AZURE_DOCUMENT_RUNTIME_PACK_PLAN.md` composes the five required layers: exact upstream acquisition, local file-security gate, SDK artifact lock, Azure runtime and strict field evaluation. No layer downloads, creates credentials or authorizes storage by itself.

## Acquisition integration

`OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.24` now locks 70 exact sources and validates optional package/license/assets/CI path-hash pairs. The Python SDK is included in both Microsoft and document-intelligence profiles with the upstream limitation visible before acquisition.

The exact source was then acquired through the rebuilt runner using a one-source, hash-linked audit approval. It produced `UPSTREAM_SOURCE_ACQUIRED`, a receipt for commit `129a5cbb…`/archive `218d4d55…`/`191429136` bytes and a source tree of `52748` files. The root license/notice and all five package/CI artifacts matched the hashes recorded above. No Git repository was created.

```text
UPSTREAM_LOCK_VALID sources=70 selected=70
UPSTREAM_ACQUISITION_TEST_PASS negatives=6
SOURCE_PROFILE_VALID profile=microsoft-secure-document-pipeline selected=13
SOURCE_PROFILE_VALID profile=document-intelligence-leaders selected=31
SOURCE_PROFILE_TEST_PASS valid=9 negatives=5 positives=1
```

## What remains unproven

- no Azure account, managed identity/key policy, region, quota, model deployment or cost boundary was supplied;
- no live/sandbox analyzer call was made;
- no approved REVESTEX corpus or independently produced ground truth exists;
- no proforma, packing list, bill of lading, customs declaration or other custom analyzer has been created/evaluated;
- no field/class has demonstrated zero false acceptance, false discovery, false negative and false positive on target data;
- no business storage is authorized.

Therefore this work adds real Microsoft leader code and a reproducible execution boundary, but it does not make the overall library ready under the user's expanded end-to-end standard.

## Failure memory

Retained local failures `LIB-FAIL-218..228` cover package/service path confusion, root license naming, wheel filename identity, tooling build-file assumptions, audit/test graph separation, materializer/path drift, optional-artifact regression output, safe diagnostic cleanup and the missing Azure lane in the global executable runner. Retained upstream condition `UP-FAIL-074` covers the Git/test-proxy dependency of the full suite. None was erased after correction.
