# Python Azure Blob Immutable Evidence Adapter — Reconstruction Evidence V1

Date: 2026-08-27. Scope: exact public Microsoft SDK identity, source/artifact equivalence, local materialization, fail-closed retained-evidence behavior and executable gates. No Azure account, credential, paid API, control-plane mutation or production container was contacted.

## Official authority, artifacts and license

- package: `azure-storage-blob 12.30.0`;
- official tag: `azure-storage-blob_12.30.0`;
- commit: `245d2024e6764bb4dd5186b5cd30a776c46993cb`; GitHub reports `verified=false`, reason `unsigned`;
- source archive: 205,202,391 bytes, SHA-256 `7c4743ccfe1cb81687ddc949df8343901c0b9a9e3c4e782ffaf4ad5f68c4a8fd`;
- wheel: 435,610 bytes, SHA-256 `d415ac50b67a8da6b3ae7e9f1014b1b55cd7aafa0b8d4ca9b380568dc7360423`;
- sdist: 618,229 bytes, SHA-256 `2cd74d4d5731e5eb6b8d5c5056ee115a5e88f8fdf22517b739836fda685018be`;
- package license: MIT, `LICENSE` SHA-256 `fd532481d828e13a0b13ccb598e02338a3617740675a862ee6bdc1541b68e93d`;
- repository root `LICENSE` SHA-256 `7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744`;
- package `NOTICE` SHA-256 `8783fc47f3420556eb8b53f746441ecb36b37f118ebec73d630e3468bde7f876`.

The PyPI wheel and sdist contain 79/79 corresponding Python files byte-identical. The selected package, tests and key metadata in the sdist match 193/193 corresponding files in the exact commit archive. The official package contains 58 `test_*.py` files, but its complete monorepo suite depends on Azure devtools, Git checkout state and live test-proxy assets; this evidence therefore does not claim that suite as executed.

The exact official SDK source establishes the primitives used by the adapter: `BlobClient.upload_blob` accepts `immutability_policy` and `legal_hold` on creation, `overwrite=False`, `validate_content=True`, and returns `version_id`; properties and bytes can be read for that exact version. Microsoft documentation also establishes version-level immutable storage prerequisites and irreversible locked retention. The local `azureblob` package is explicitly `AUTHORED` integration glue; it is not represented as Microsoft application code.

## Materialized artifact

Pack SHA-256: `23890dae2e775e0f0922541c00912fd25bb7b7670219b1ce54905df8d2bac31b`.

| Path | Bytes | SHA-256 |
|---|---:|---|
| `azure_blob_storage_adapter/azureblob/__init__.py` | 414 | `48425acd30bdf81d93b72a71150910ca0986294d4208e4ff1a4bdc4a8c3ef98c` |
| `azure_blob_storage_adapter/azureblob/azure_backend.py` | 3,727 | `ece9bdeeb3c4f3d467258f6354c3cae8a2a4c56a96e5dff0d5784a1ed418450e` |
| `azure_blob_storage_adapter/azureblob/storage.py` | 16,869 | `913d55459f70b6d2e83ec57e11783f8c62d2e2cf9cdad2872f5b6025cbfa63ab` |
| `azure_blob_storage_adapter/official-artifact-lock.json` | 1,038 | `bdd491579586eb07026a497fc547467f5eaad4b2ca824b97298517a33770c6d2` |
| `azure_blob_storage_adapter/provider-profile.template.json` | 673 | `645e14a65074995eafebd3a759b3e4cbb6a861ce6c0caef16ec1fb306dbc9e52` |
| `azure_blob_storage_adapter/README.md` | 2,320 | `a07c02e4564ca8ff8ae7c65f4547871eabf4de9717d3244088e72cb04c645b51` |
| `azure_blob_storage_adapter/requirements-direct.in` | 244 | `6b5a20c1db334ba9c03f477f5b8d28e7018bc74bd4624c13add2532e96a79bbb` |
| `azure_blob_storage_adapter/requirements-windows-py312.lock` | 1,172 | `2156fc146c023a35e110422f8c766948e84f68ef50a7e2231b8720db23eb3bf8` |
| `azure_blob_storage_adapter/test_storage.py` | 13,392 | `088f90d6ecc37037560a39f13c9c21c625f5badf08c73d65c128472255ec1853` |

The root materializer recreated all nine files from Markdown and verified every declared SHA-256. `AZURE_DOCUMENT_RUNTIME_PACK_PLAN.md` composes acquisition, secure local ingestion, official SDK artifacts, Azure Content Understanding runtime, strict field evaluation and Azure Blob retained evidence as 6 packs / 50 implementation files.

## Executable behavior

An exact CPython 3.12 environment installed the 12-wheel Windows x64 lock with `pip --require-hashes`. `pip check`, `compileall` and the reconstructed suite passed:

```text
python -m unittest discover -s azure_blob_storage_adapter -p test_*.py -v
Ran 15 tests in 0.003s
OK
python -m compileall -q azure_blob_storage_adapter/azureblob
PASS
python -m pip check
No broken requirements found.
```

The tests cover exact official SDK call signatures, private versioned immutable container prerequisites, non-overridable encryption scope, target/authority hash binding and freshness, access/cost/retention/irreversibility approvals, separate legal-hold authorization, size/hash/name validation, deterministic create-only writes, MD5 transport validation, SHA-256 evidence verification, exact-version property/download reads, enum normalization, unknown create effect, partial verification and reconciliation without a second write.

The adapter never creates or modifies accounts, containers, RBAC, networking, encryption scopes, versioning or WORM policy. A create exception becomes `UNKNOWN_EFFECT`, not a blind retry. A mismatch after creation becomes `PARTIAL_VERIFICATION` with the exact `version_id` required for reconciliation.

## Dependency security evidence

The locked runtime graph contains 12 exact wheels: Azure Core 1.41.0, Azure Storage Blob 12.30.0, certifi 2026.7.22, cffi 2.1.1, charset-normalizer 3.5.1, cryptography 50.0.1, idna 3.19, isodate 0.7.2, pycparser 3.0, requests 2.34.2, typing-extensions 4.16.0 and urllib3 2.7.0. An OSV batch query for the exact package/version pairs returned zero findings. Any dependency or platform change invalidates this result and must repeat hash acquisition, install, OSV, compile and behavior gates.

## Result and limits

Result: `REBUILD_VERIFIED / CONDITIONED`. The code is immediately materializable and locally executable. It cannot be enabled for a project until the user supplies and approves a real Azure account, exact resource ID, private container, RBAC/network configuration, version-level immutability/versioning, non-overridable encryption scope, retention/legal policy, cost/quota, security admission, corpus and recovery evidence. No live WORM write, provider concurrency/load or deletion/retention recovery was exercised. The unsigned release-tag condition remains open and is not converted into a signature claim.

