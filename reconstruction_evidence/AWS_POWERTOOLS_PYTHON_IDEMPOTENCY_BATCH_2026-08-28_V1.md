# AWS Powertools Python Idempotency and Batch — V1

Fecha: 2026-08-28  
Clasificación: `PINNED_CANDIDATE_CONDITIONED`

## Identidad oficial

- repositorio: `aws-powertools/powertools-lambda-python`;
- release vigente por API oficial: `v3.34.0`, publicada 2026-08-10;
- tag annotated unsigned `c4490bca94f6d59191c87123ee70348fcbc1cf3e`;
- commit unsigned `376757161b002f0c2f5d19d5cdf9def5b8c45704`, tree `411d3dbf377058aaf3eedfd3daee92e76ad20261`;
- source ZIP 9.650.888 bytes, SHA-256 `b7b8b0599b039275b68138837a4644061f84fbedba624c2f434d4d9e7b5c4742`;
- MIT-0; `LICENSE` SHA-256 `bcc1d2d9e0609b03d420c3dbecb438a7b89bfc53d54176ab56301ab1f96b4f0a`;
- `pyproject.toml` SHA-256 `4d97425bbfa9b27e3ae9c2a29e36e81d77cc98bac09a06ccff32dfdf758c1203`;
- `poetry.lock` SHA-256 `5085ef379180666ccd470973998c474141cd8e7fa66a250d6ed911f358ab70f7`;
- archive: 2.366 files/633 directories, 1.311 Python, 29 workflows.

## Distribución y provenance

- wheel PyPI oficial `aws_lambda_powertools-3.34.0-py3-none-any.whl`: 957.321 bytes/SHA-256 `ab1354c58085ccecf92e9c548b1336ad5bae72c7ae23e61d2aa0ff15520d679c`;
- sdist PyPI: 800.151 bytes/SHA-256 `75f2c65a5997630666c9c3c495002c1f9d1e9defd3cb6d5d9fc7272b7dc9ef54`;
- release asset SLSA `multiple.intoto.jsonl`: 24.011 bytes/SHA-256 `76e168c0d535396e801957999871749af8021fa156661c8a171cd298058e50d6`;
- la attestation nombra exactamente ambos hashes y builder SLSA GitHub Generator 2.1.0;
- material atestado: `develop@715012fe0b998d18cb639b4a5b5d8e99543296e3`; el commit release es su hijo de version bump. La provenance declara `reproducible=false`;
- los 314 archivos de `aws_lambda_powertools/` son byte-idénticos entre tag source y sdist atestado; cero faltantes/diferencias.

## Ejecución

- Poetry 2.2.1 validó el lock e instaló el grafo en venv aislado después de fijar `VIRTUAL_ENV`;
- subset idempotency DynamoDB/Pydantic y batch: `189 passed`, 64 warnings;
- suite local unit+functional: `2.492 passed`, `11 failed`, `3 skipped`, 233 warnings;
- ocho fallos son portabilidad Windows reproducible: tempfile ZIP abierto, paths POSIX `/`/`/var/task`/`PurePath`, dos asserts `os.linesep` y dos tests de carrera Redis que Windows `spawn` no puede serializar (`KeyedRef`);
- tres fallos de logger/metrics pasan aislados y revelan dependencia de orden/estado de la suite;
- 11 archivos E2E idempotency requieren recursos AWS y no se ejecutaron sin cuenta/autoridad;
- wheel oficial + extra real `aws-sdk` instala/importa 3.34.0 y `pip check` pasa;
- SCA de dependencias de aplicación: cero advisories conocidos. El `pip 25.0.1` de tooling del venv registra 7 ocurrencias/6 IDs; se preserva y no se atribuye a Powertools;
- Pydantic no es un extra publicado por 3.34.0: debe declararse por separado si el proyecto lo utiliza.

## Código focal

| Path | SHA-256 |
|---|---|
| `aws_lambda_powertools/utilities/idempotency/persistence/dynamodb.py` | `1f4766f14d5dc0d3f3007b660110683c8abe6a32358a1eb66e04055f5472571a` |
| `aws_lambda_powertools/utilities/idempotency/persistence/base.py` | `08351f32c160d3c7805e5e3bb6fcf24ef4215c40cf82a5fe0f6b2899e85f2c74` |
| `aws_lambda_powertools/utilities/idempotency/config.py` | `87c52b68c6dec08dfd2a820887b16907f7d7b48bffbf45a3375511b6f2df7894` |
| `aws_lambda_powertools/utilities/batch/base.py` | `3f3466f53e92a5972e898042ccb38a97042befd318ef46a784ff5fc036f609de` |
| `aws_lambda_powertools/utilities/batch/decorators.py` | `3f27f636eb60949b265db0d74befc9c6034bf8b3cb0c41f67b97b1d60c13529e` |

La persistencia oficial implementa conditional create, expiración general e `in_progress_expiration`, payload validation y retorno de resultados. En esta revisión, `_update_record` y `_delete_record` operan por key sin `ConditionExpression` ni owner token. Por tanto, la utilidad no constituye por sí sola fencing completo después de expirar/reclamar un INPROGRESS.

## Admisión

Admitir como componente oficial pinneado para idempotency, partial batch y observabilidad, nunca como worker object-event completo. El proyecto todavía debe aportar y probar: key exacta bucket/key/versionId/eTag según evento, comparación/order cuando aplique, owner/fencing o semántica equivalente, SQS visibility, partial failures, retries/DLQ, replay/reconciliation, side effects idempotentes, IaC/IAM/KMS/TTL/PITR/alarms, live DynamoDB/Redis y E2E AWS. No se atribuye a AWS ningún wrapper local.
