# AWS IDP immutable evaluation handoff — evidencia V92

Fecha de corte: 2026-08-28. Esta evidencia gobierna únicamente el pack V92 y no sustituye pruebas LIVE del proyecto.

## Identidad oficial

- Repositorio: `aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws`.
- Tag exacto: `v0.6.5`.
- Commit exacto: `1b5fd74454e593de233342a02ee911af8ee38359`.
- Licencia upstream: MIT-0.
- `idp_common` importado realmente: 0.6.5; Boto3/Botocore del lock: 1.42.97/1.42.97; CPython 3.12.14.
- Contratos focales: `lib/idp_common_pkg/idp_common/hooks/__init__.py`, sus tests unitarios, `patterns/unified/src/pipeline_hooks_function/index.py`, `patterns/unified/statemachine/workflow.asl.json` y `docs/feature-platform.md`.

AWS publica que el hook `postprocessing` es síncrono, posterior a evaluation y anterior al terminal; puede usar `onError: fail`, recibe `Document` inline o comprimido, y se ejecuta dos veces cuando HITL queda pendiente y luego concluye. El helper oficial `load_hook_document` resuelve ambas formas. El código local no se atribuye a AWS: cuatro bloques son `ADAPTED` y cuatro `AUTHORED`.

## Resultado materializable

- Pack: `implementation_packs/AWS_IDP_IMMUTABLE_EVALUATION_HANDOFF.md`.
- SHA-256 canónico: `e0023a57a9e94fa25ee0573bf1ca1a80f7cf03f448db3d7e3e36ea9142fb792a`.
- Manifest/bloques: 8/8, operaciones `CREATE`, hashes exactos.
- Materialización aceptada: `%LOCALAPPDATA%\Temp\elite-v92-clean-roundtrip-9664386c35d143969c23fd10f3e5db00`.

| archivo | bytes | SHA-256 |
|---|---:|---|
| `deploy.ps1` | 2.264 | `c7485ba38378fd7909f5db612cbb710a4949e813882418f1b00c781a8df74abd` |
| `handler.py` | 12.955 | `a7607d5080f141e327f10c0ad3902a1347eced2ad17e89a32a40104162578d8d` |
| `hook-config.template.yaml` | 393 | `2239a38c0454171109bc6efe417c299d6c28e9c2399563165d5fab589519359a` |
| `provider-profile.template.json` | 1.086 | `f191ef103dd341ab56b4439eb22c446f3cd97b9987290d29cef994ae4890c789` |
| `README.md` | 2.012 | `fd6348fe92bb623354acde2354e4c234f9586b974eee73569edf8022f0cfe694` |
| `template.yaml` | 6.351 | `2d6268ddb0c41b229db9a54d57a929494553450b3aaf4730388b8d1c2d24c2ee` |
| `test_handler.py` | 9.805 | `fa74178b429c5dffac05e6fc419c04d34716544eb5a7026c88825b41dde77645` |
| `verify_contract.ps1` | 1.969 | `30f42440c3fabde58f4669b836cdc0023a9e93599099efe4e41546432fab8ac0` |

## Gates ejecutados

1. Materialización desde Markdown: `Materialized 8 files`.
2. Compilación en memoria de `handler.py` y `test_handler.py` con CPython `compile(...)`: PASS y cero `.pyc`/`__pycache__`.
3. `python -B -m unittest -v test_handler.py`: 14/14 PASS. Cubre terminal, HITL pending/in-progress/completed/skipped, replay, colisión divergente, error SQS, config/clase/ejecución incorrectas, issues, redacted superseded, oversize y storage authority siempre false.
4. `cfn-lint` 1.55.1 sobre `template.yaml`: exit 0.
5. Round-trip real con `idp_common` 0.6.5: `Document(Status.EVALUATING, evaluation_status=COMPLETED)` → `load_hook_document` → `to_dict`, con config, execution y classification preservados: PASS.
6. Profile distribuido: acknowledgement false, allow-list vacía, retención/tamaño cero y todas las evidencias/owners vacíos: PASS fail-closed.
7. Config de hook: `onError: fail`, `allowDocumentUpdate: false`, feature/tenant/config/schema fijados: PASS.
8. `VERIFY_LIBRARY_PASS`: 66 packs, 622 archivos materializables, 413 Markdown y los 35 perfiles; email 7/71.
9. `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 66 packs/111 fuentes, `aws_idp_immutable_evaluation_handoffs=1` y las 14 pruebas del handoff dentro de la auditoría integral.

## Garantías y límites honestos

El handler carga el `Document` mediante el helper AWS, valida config/args/ejecución/clases/issues/HITL, canonicaliza los bytes observados, crea snapshot y receipts con `IfNoneMatch:*`, SHA-256, VersionId, SSE-KMS y Object Lock COMPLIANCE, entrega por SQS FIFO y compromete el dispatch receipt al final. Un replay con receipt final no reenvía. HITL pendiente sólo crea un receipt diferido. Cada respuesta y mensaje fija `automaticBusinessPersistenceAuthorized=false`.

No se afirma exactamente-once empresarial: una caída después de `SendMessage` y antes del receipt puede producir un reenvío fuera de la ventana FIFO, por lo que el consumidor debe deduplicar durablemente por `evaluationId`. Tampoco se afirma deploy AWS, precisión universal, corpus target, política HITL correcta, KMS/Object Lock/alarms/DLQ/redrive/restore/costo aprobados ni persistencia transaccional. La salida queda lista para el siguiente gate ejecutable, no para escribir directamente en ERP/negocio.

## Fallos convertidos en memoria

Quedaron registrados `LIB-FAIL-1036` a `LIB-FAIL-1042`: resolución de uv, cache offline incompleta, patch multiarchivo inválido, recurrencia de backticks PowerShell, manifest sin LF, cleanup bloqueado y bytecode de `py_compile`. `UP-FAIL-178` conserva la limitación upstream que impide tratar el ejemplo post-processing ERP o una referencia comprimida no versionada como frontera de persistencia.
