# AWS IDP evaluation decision worker — evidencia V94

Fecha: 2026-08-28  
Estado: `REBUILD_VERIFIED / CONDITIONED`  
Pack: `AWS-IDP-EVALUATION-DECISION-WORKER` 0.1.1  
Pack SHA-256: `2fa04488966fc9939a0110ea43f498c7af39cf8e6af33901c1979f88d432edd2`

## Fuente y decisión de admisión

Base oficial: AWS Powertools for Lambda Python 3.34.0 release/commit `376757161b002f0c2f5d19d5cdf9def5b8c45704`, cuyos ejemplos fijados combinan idempotencia por registro, timeout context y partial-batch. AWS documenta entrega SQS al menos una vez, clave de negocio elegida por productor, payload validation, clave obligatoria y expiración configurable.

El sample oficial `aws-samples/transactional-outbox-pattern` main `23e0519c7a8048c6db4acd8f3d7de534a9d5deac` fue auditado y rechazado para copia inmediata: commit 2024 no firmado, dependencias antiguas, sólo context tests y borrado completo de filas después de `sendMessageBatch` sin comprobar fallos parciales. `UP-FAIL-180/181` preservan ambos límites. La guía AWS de outbox se conserva como contrato; el código sample no se incorpora.

## Resultado materializado

| Archivo | Bytes | SHA-256 |
|---|---:|---|
| `deploy.ps1` | 1.963 | `a3823849ba2a902cc58c5f57b84110e70ac20214cdbacb9dc620dfce5b1c4c6d` |
| `evaluation-policy.template.json` | 405 | `6a1d92736fa3c49dfde1491c0df57699d4308f9dfa3c557d5503e3cdb5048fcd` |
| `handler.py` | 17.602 | `e99ee5468fd90cc44f07cde097c612901290b5aa30e438422dcaf1613a57a11a` |
| `provider-profile.template.json` | 1.038 | `be795f524c12533ad659549f8872953221bbcec90c2fb411dfa75c8c0b15e542` |
| `README.md` | 1.601 | `15bf6b267a87194954432dbbb08ae51ac1ca08c23f76415b67ee296056be856c` |
| `template.yaml` | 4.952 | `d9fdcfcd57ceabae0db5b8a271b67505c62a66253b08c668da24514a2d7f46d8` |
| `test_handler.py` | 12.584 | `5124b7270c909297fd44b4b1dce188fd65258e29610bab8b30e2e5ff973c23e6` |
| `verify_contract.ps1` | 1.921 | `1563020506885f7518a6f90a83ea7cbbe6a7a2e425c15b7543f92b20ce0f56da` |

Seis archivos son `AUTHORED`, dos `ADAPTED`, cero `VERBATIM`. No se atribuye a AWS el handler, policy ni template combinados.

## Gates focales

- materialización y SHA staging→round-trip 8/8: PASS;
- compilación sin bytecode y 13/13 tests: PASS;
- cfn-lint 1.55.1: exit 0;
- negativos: evaluationId/payload/snapshot/review/policy/storage tamper, referencia o review cross-tenant, schema/clase/sección, skipped/no-review, replay y colisión divergente;
- `SchemaValidationError` oficial se convierte en decisión durable `SCHEMA_VALIDATION_FAILED`; errores de infraestructura continúan siendo retry;
- decision receipt sin documento, valores/reviewer y con `automaticBusinessPersistenceAuthorized=false`;
- profile/policy distribuidos comienzan cerrados.

## Memoria y límites

`LIB-FAIL-1051` a `1057` registran conteo de objeto JSON, ciclo CloudFormation, patch documental atómico, namespace cross-tenant, schema inválido inicialmente tratado como retry y un `$LASTEXITCODE` residual del primer round-trip. Memoria vigente: 1.057 locales + 181 upstream = 1.238 IDs, cero abiertos locales antes de gates globales.

No hubo AWS live. Faltan cuenta/IAM/layer exacto, policy/corpus/HITL target, KMS/Object Lock/Dynamo/SQS reales, duplicado/race/load, DLQ/redrive, restore/rollback/costo y consumidor de persistencia negocio+outbox transaccional. Una decisión positiva no autoriza por sí sola escritura empresarial.

## Gates globales

`VERIFY_LIBRARY_PASS` confirmó 67 packs/630 archivos/416 Markdown/35 perfiles, incluido email 10/92. `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` confirmó 67 packs/111 fuentes y `aws_idp_evaluation_decision_workers=1`; las 13 pruebas del worker 0.1.1 ejecutaron dentro de la auditoría integral.
