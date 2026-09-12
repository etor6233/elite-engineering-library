# All Implementation Packs Materialization — 2026-08-27 V44

## Alcance del snapshot

Este snapshot sucede a V43 sin reescribirlo. Conserva 46 packs, 413 archivos materializables, 23 perfiles y 80 sources; fija los dos ejemplos oficiales Microsoft Durable Task Human Interaction y prueba de nuevo toda la biblioteca sin convertir un endpoint anónimo en autorización productiva.

Estado canónico al cierre:

- 46 implementation packs;
- 413 bloques `CREATE` materializables: 401 `AUTHORED`, 3 `ADAPTED`, 9 `VERBATIM`;
- 23 perfiles de composición y 10 source profiles;
- 80 upstream source archives fijados por identidad, bytes, SHA-256 y licencia;
- 5 artefactos SDK documentales y 1 artefacto Business Central;
- 9 provider adapters, 1 orquestador documental y 1 log de evidencia;
- 348 fallos locales y 104 condiciones upstream preservados al corte;
- 299 Markdown y 8 scripts PowerShell canónicos después de incorporar este expediente.

## Cambio material

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` sube a `0.4.35` y `MICROSOFT-DURABLE-DOCUMENT-ORCHESTRATION` a `0.1.1`. El archive oficial `microsoft/durabletask-python` permanece como una sola source y ahora liga, mediante paths y SHA-256 exactos, estos artefactos del release `v1.9.0`, commit firmado `6dfdbac521d9d59d31d7ea975fb669d66f86f4c4`:

- `examples/human_interaction.py`, SHA-256 `7f97502962fac35775057651a2081f5e4b7e19bc116dba7145b70c8afb467669`;
- `azure-functions-durable/samples/human-interaction/function_app.py`, SHA-256 `42fdd853d67fcec3acc51ee12e6df1f4ac6c221d246f77a24585d023ec3f0e83`.

Los ejemplos demuestran `wait_for_external_event`, timer durable, `when_any` y entrega de evento desde el cliente. No demuestran identidad ni autorización del aprobador: el ejemplo Azure declara `AuthLevel.ANONYMOUS` y entrega JSON sin una decisión empresarial autenticada. Por ello se admiten como evidencia mecánica y nunca como frontera HTTP productiva. La condición queda preservada en el lock, pack, readiness, catálogo y `UP-FAIL-104`.

## Pruebas ejecutadas

Desde materializaciones limpias:

```text
Materialized 19 files
UPSTREAM_ACQUISITION_TEST_PASS negatives=7
UPSTREAM_LOCK_VALID sources=80 selected=80
SOURCE_PROFILE_TEST_PASS valid=10 negatives=5 positives=1

Materialized 13 files
13/13 tests PASS sobre backend oficial Microsoft in-memory
6 wheels instalados con --require-hashes
pip check PASS
0 OSV
ambos sources oficiales: identidad, SHA-256 y sintaxis PASS
```

El runtime oficial in-memory continúa limitado a pruebas. El servicio Durable Task Scheduler, credenciales, red, identidad del revisor, autorización por recurso, deduplicación del evento y evidencia de la decisión siguen siendo condiciones del proyecto y no se infieren de los samples.

## Auditoría global

La verificación previa a este expediente produjo:

```text
VERIFY_LIBRARY_PASS
packs=46 materialized_files=413 markdown_files=298

VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=46 upstream_sources=80 document_sdk_artifacts=5 business_central_artifacts=1 provider_adapters=9 document_orchestrators=1 evidence_logs=1
```

La verificación estructural final posterior debe producir exactamente:

```text
VERIFY_LIBRARY_PASS
packs=46 materialized_files=413 markdown_files=299
```

## Lecciones preservadas

`LIB-FAIL-343` a `LIB-FAIL-348` conservan los fallos reales de interfaz supuesta, runtime base sin dependencia, compilación que intentó escribir bajo un path profundo, glob inválido en Windows, patch duplicado y referencia narrativa que duplicó un ID del ledger. Todos requieren y reciben regresión antes del cierre; ninguno se borra del historial.

## Estado honesto

La biblioteca continúa incompleta bajo el estándar ampliado del propietario. Durable Task aporta mecánica durable oficial reutilizable, pero todavía falta una referencia oficial, ejecutable y admisible para revisión humana autenticada/autorizada de punta a punta. También permanecen condicionados los runtimes que necesitan red o cuenta real y la evidencia productiva específica de cada negocio.

No existe autorización para Git, ZIP final, credenciales, cuentas, despliegue ni servicios pagos.

## Fuentes oficiales

- <https://github.com/microsoft/durabletask-python/tree/6dfdbac521d9d59d31d7ea975fb669d66f86f4c4/examples>
- <https://github.com/microsoft/durabletask-python/tree/6dfdbac521d9d59d31d7ea975fb669d66f86f4c4/azure-functions-durable/samples/human-interaction>
