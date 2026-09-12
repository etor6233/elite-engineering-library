# AWS IDP review authority — evidencia V93

Fecha: 2026-08-28  
Estado: `REBUILD_VERIFIED / CONDITIONED`  
Pack: `AWS-IDP-IMMUTABLE-EVALUATION-HANDOFF` 0.1.2  
Pack SHA-256: `68d12314d2b5853eca1898988f88bc1c9a3e31c49d6304bf86c0da6d1da074af`

## 1. Fuente oficial y límite observado

Se inspeccionó AWS GenAI IDP Accelerator tag `v0.6.5`, commit `1b5fd74454e593de233342a02ee911af8ee38359`, especialmente `src/lambda/complete_section_review/index.py` y `idp_common.hooks.load_hook_document`. El código AWS conserva RBAC, claim condicional, revisión por sección e historial; no se reatribuyen a AWS los controles locales.

El handler oficial de review actualiza `HITLReviewHistory` mediante read/append/update, escribe ediciones S3 por bucket/key y permite `Review Skipped` con `HITLCompleted=True`. No emite un receipt inmutable unido a documento/config/revisores. Por eso `HITLCompleted` aislado no se admite como autorización de persistencia. La condición queda registrada como `UP-FAIL-179`.

## 2. Cambio ADAPTED verificable

El handoff 0.1.2 añade lectura DynamoDB `ConsistentRead=True` sobre la fila oficial `PK=doc#<document-id>, SK=none`. Para `Completed` exige:

- `HITLCompleted=True` y estado distinto de `Review Skipped`;
- listas pending/skipped vacías;
- cobertura exacta y no duplicada de todas las secciones tanto en el `Document` como en tracking;
- exactamente un registro de review por sección, reviewer conocido y timestamp con zona;
- prueba normalizada con hashes SHA-256 de subject/email, nunca identidad raw, y binding autocontenido a document SHA, config version, schema profile y execution ARN hash;
- objeto review-evidence versionado, SSE-KMS y Object Lock; su hash/VersionId se incorpora al evaluation id, mensaje y dispatch receipt.

`Skipped` se entrega sólo como `SKIPPED_NOT_REVIEWED`, sin review evidence. Todas las respuestas, mensajes y receipts conservan `automaticBusinessPersistenceAuthorized=false`.

## 3. Archivos materializados

| Archivo | Bytes | SHA-256 |
|---|---:|---|
| `deploy.ps1` | 2.384 | `5c13dc48bcf10c74df1ab6c02c300eae5e90b1af33c4faca383726b9f6515dc4` |
| `handler.py` | 18.759 | `f3e27bd5f60b5665e8f5e5ba48051d346408830ebb3e813b7a31fa4d42e6b7c4` |
| `hook-config.template.yaml` | 393 | `2239a38c0454171109bc6efe417c299d6c28e9c2399563165d5fab589519359a` |
| `provider-profile.template.json` | 1.233 | `946302965c2c247c3988e320c6638917ae89fd05a6429147e08baceb7d3f5f71` |
| `README.md` | 2.674 | `c20429a6a02cafeabc4bc945b2532ea5909e24a7f6e227ca876ebbf30e78cb80` |
| `template.yaml` | 6.711 | `fdd7ec0e6576d2e403ae6f8ac91d08a1c39a3b159e21435fcbf4345c1438d7f7` |
| `test_handler.py` | 14.402 | `f92aacbab93c3dcec3c30d188e7e941b55c1d76f644c01450db88d7f61acaf4f` |
| `verify_contract.ps1` | 2.058 | `0f9e48d16201cfa48b0732e821bda98a1833080fb4e41655c80c232f8b304ab9` |

## 4. Gates ejecutados

1. Materialización canónica 8/8 y comparación SHA staging→round-trip: PASS.
2. Compilación en memoria sin bytecode: PASS.
3. `python -B -m unittest -v test_handler.py`: 20/20 PASS.
4. Negativos: tracking ausente, skip-all, pending/skipped residual, cobertura incompleta, reviewer unknown, timestamp sin zona, sección desconocida/duplicada y fuga de PII: rechazados.
5. `cfn-lint` 1.55.1 por import root fijado: exit 0.
6. Profile fail-closed, IAM `dynamodb:GetItem` exacto, hook `onError: fail`/no mutation y primitives Object Lock/FIFO/no-storage: PASS.

## 5. Fallos convertidos en memoria

`LIB-FAIL-1043` a `LIB-FAIL-1050` registran nombre de pack inferido, resolución/import/metadata de cfn-lint, binding de array del updater, staging detectado por el gate global, filenames abreviados y la falta de binding autocontenido detectada antes de diseñar el consumidor. El resultado aceptado se reconstruyó después de cada corrección. La memoria vigente queda en 1.050 fallos locales + 179 condiciones upstream = 1.229 IDs únicos, cero abiertos locales.

Después del cleanup exacto, `VERIFY_LIBRARY_PASS` confirmó 66 packs/622 archivos/414 Markdown/35 perfiles y `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` confirmó 66 packs/111 fuentes, incluida la suite 20/20 del handoff dentro de la auditoría integral.

## 6. Límites honestos

No se ejecutó AWS live ni se afirma que el tracking mutable sea una firma AWS. Una carrera concurrente después del consistent read sigue siendo condición upstream; el receipt prueba la vista leída y congelada. Faltan cuenta/región target, KMS/Object Lock/DynamoDB/SQS reales, corpus y políticas aprobados, consumidor idempotente, persistencia transaccional, reconciliación, carga, DLQ/redrive, restore y rollback. Nada en V93 autoriza almacenamiento empresarial automático.
