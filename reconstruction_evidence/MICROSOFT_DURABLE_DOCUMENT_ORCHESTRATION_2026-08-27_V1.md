# Microsoft Durable Document Orchestration — Reconstruction Evidence V1

Date: 2026-08-27. Scope: exact Microsoft Durable Task source/package identity, official Azure agent skill, clean Markdown reconstruction, frozen runtime, deterministic document orchestration and fail-closed effect handling. No Azure account, credential, paid API, production DTS backend, Docker image or external storage was contacted.

## Official authority, artifacts and license

- package: Microsoft `durabletask 1.9.0`;
- release: `v1.9.0`, published 2026-07-30;
- exact commit: `6dfdbac521d9d59d31d7ea975fb669d66f86f4c4`; GitHub reports its commit signature verified and valid;
- exact commit archive: 1,418,256 bytes, SHA-256 `c8c941b25804ae9abcc0470f5594cab4929216a981982616344eb2b84a7916e8`;
- release-tag archive retained for comparison: 1,387,666 bytes, SHA-256 `5dbfad8cf7d47cc3d86ffd2d9cc78d25daa0ee8b2210ec678715eb042350e045`;
- PyPI wheel: 212,413 bytes, SHA-256 `a759af4ad8e6897922575886e93bf40c6b3af936eaed83798d492ee561607185`;
- PyPI sdist: 184,873 bytes, SHA-256 `55460fbfda8941e721096c4639e72111b03638708df1556af466160ee649478d`;
- MIT `LICENSE`: 1,141 bytes, SHA-256 `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`;
- root `pyproject.toml`: SHA-256 `a32e6ea59bc34568111697c3a8abb5e38a528e8c8f053e573b3d86eb3ad97902`.

The release-tag and exact-commit archives contain 368/368 corresponding files byte-identical after removing the different archive root directory name. Wheel and sdist contain 60/60 corresponding Python files byte-identical; those 60 files also match the exact commit archive. The acquisition lock therefore names the commit URL directly and does not infer a tag URL or PyPI artifact URL.

The official core source suite was executed from the exact source with the complete upstream test requirements: 1,034 passed, one skipped and 14 deliberately deselected in 399.20 seconds. DTS/Azurite integration groups were not executed. Twelve `datetime.utcnow()` deprecation warnings are preserved as `UP-FAIL-081`; the worker event loop warning is preserved as `UP-FAIL-082`. Neither was patched or hidden.

The official Azure Samples authority is `Azure-Samples/Durable-Task-Scheduler` at signed commit `5636c25ffbdaabaca1062c3b6c6e77072f2b2eb2`; archive 8,425,148 bytes, SHA-256 `f9754c7d0f9d957aa01ec541d78f838fe097b0cdb0fcb57d7decb8b0c0e55347`; MIT `LICENSE.md` SHA-256 `d9a1b1e30d633d5732ea18e3cba9538d293ebc53e1a9e4e96ab739e0c5c4f1cb`. Its official `.github/skills/durable-task-python/SKILL.md` is packaged byte-verbatim: 15,429 bytes, SHA-256 `553f0dc278c63bb129df5557a7d6f7d645bf827efba4cb077588688f88a2e9e2`. The materialized license adds only the final newline and is declared `ADAPTED_FINAL_NEWLINE_ONLY`, never verbatim.

## Materialized artifact

Pack SHA-256: `d37c0e17b67f68c40108803d15bab2aea869bf94fa819bda9cab5272c0155395`. Acquisition pack V0.4.28 SHA-256: `8126c41fbde16912059746521eb41159f5a401cba97850813e5d3b2f8fdedde3`. Vendor-neutral plan SHA-256: `10bfa0b9eeb46e975a6d8ccc690a540d89ac8e51a8c378a8a17721fa3ea5d2ba`.

| Path | Bytes | SHA-256 |
|---|---:|---|
| `durable_document_orchestration/requirements-direct.in` | 19 | `43faf914614ecd498cd0613ed87f5e6eb309ccd38db223c7872472d7044e7aae` |
| `durable_document_orchestration/requirements-windows-py312.lock` | 618 | `4d7cca8b66e29db5226811ef09e7532fdac1e0afbe2ab1e8fe247bcbc4e4ed7f` |
| `durable_document_orchestration/official-artifact-lock.json` | 1,721 | `1fcfef02326a5308b221363ab074051a91a93a88416e63193cbc7ccd5d714259` |
| `durable_document_orchestration/pipeline-profile.template.json` | 555 | `6df767c2e34ab66ecad033f41a9ba4b3ae2d184a6888e859fea7f4b0262e5426` |
| `durable_document_orchestration/documentflow/__init__.py` | 426 | `7301e1ee94d4d18b855b410a1365d2a8bc977dc3dc6417e8d3de0568faf384e9` |
| `durable_document_orchestration/documentflow/contracts.py` | 8,238 | `f591190226211095c9ebef22652737491232ceefa5d8c1b3be0ab8fb1559d6dd` |
| `durable_document_orchestration/documentflow/orchestration.py` | 6,118 | `af72a431a24bbd6403d94681d4a075542987dd5b7ef2318f180d181f7c68caee` |
| `durable_document_orchestration/documentflow/runtime.py` | 1,728 | `06ff429a1f2c45d94e6bc20871d76103d2a7c54eb25cc239a13b6358c367e6f4` |
| `durable_document_orchestration/test_contracts.py` | 2,280 | `79e9156721ed9dd051a930854bb906fb4b5a1a37a076b8de2d8f19aaff4dd801` |
| `durable_document_orchestration/test_pipeline.py` | 9,432 | `bad71391dc699fb96a155ad136765d88cee5537bb4e2bc5513186a7ace30ebe3` |
| `durable_document_orchestration/README.md` | 1,662 | `ee2af1910afe57ddd1e52aadbb605fb6f4c1076cac79b67a1cefc121c1e0632e` |
| `durable_document_orchestration/upstream/LICENSE.md` | 1,141 | `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383` |
| `durable_document_orchestration/upstream/durable-task-python-SKILL.md` | 15,429 | `553f0dc278c63bb129df5557a7d6f7d645bf827efba4cb077588688f88a2e9e2` |

The root materializer recreated all 13 files and verified manifest↔FILE↔SHA equality. `DURABLE_DOCUMENT_PIPELINE_PACK_PLAN.md` composes five packs/46 files. The provider-specific profiles compose Azure seven/63, Google seven/63 and AWS six/58, each with a valid final-hash materialization record.

## Executable behavior

The frozen CPython 3.12 Windows x86-64 graph contains six exact wheels with hashes: asyncio 4.0.0, durabletask 1.9.0, grpcio 1.83.0, packaging 26.3, protobuf 7.36.0 and typing-extensions 4.16.0. A clean environment installed them with `pip --only-binary=:all: --require-hashes`; `pip check`, `compileall` and all reconstructed tests passed:

```text
python -m unittest discover -s . -p test_*.py -v
Ran 13 tests in 9.210s
OK
python -m pip check
No broken requirements found.
```

The flow is security → retained original → extraction → strict evaluation → optional durable human review → idempotent business persistence → retained final evidence. Orchestration history accepts opaque references, hashes and receipts rather than document bytes or extracted fields. The distributed profile enables no backend, class, storage lane, review timeout or business effect.

Tests prove closed request schemas, hash-bound idempotency and review events, no persistence after security rejection/human rejection/review timeout, exactly-once activity invocation in the in-memory harness, reconciliation after ambiguous original storage, no blind retry after ambiguous business persistence, and preservation of a committed transaction when final evidence becomes a partial effect.

## Dependency security and limits

An official OSV batch query for all six locked package/version pairs returned zero findings on 2026-08-27. This is dated evidence, not a future guarantee. Any dependency, Python target, backend or source revision change reopens identity, license, equivalence, OSV, install and behavior gates.

Result: `REBUILD_VERIFIED / CONDITIONED`. The code is immediately materializable and executable against Microsoft’s official in-memory test backend. That backend is test-only and stores state in memory. Production remains blocked until the project supplies an exact DTS/service or self-hosted backend identity, digest/package lock, TLS, managed identity/RBAC, task-hub isolation, durable storage, versioning strategy, load/soak, recovery, observability and live adapter evidence. No claim of perfect extraction, exactly-once distributed effects or production readiness is made.
