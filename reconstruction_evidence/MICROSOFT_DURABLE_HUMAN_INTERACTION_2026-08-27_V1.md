# Microsoft Durable Task Human Interaction — 2026-08-27 V1

## Decisión

`PINNED_CANDIDATE / CONDITIONED` para la primitiva durable; `SAMPLE_ONLY` para las fronteras de aprobación.

Microsoft sí publica código exacto y ejecutable para el patrón mecánico que faltaba: esperar un evento externo de aprobación, competirlo contra un timer durable y reanudar/cancelar el workflow. Ese código no autentica al aprobador ni implementa autorización empresarial. Elite puede reutilizar la primitiva oficial y su SDK; no puede presentar los endpoints de muestra como una bandeja productiva.

## Identidad oficial

- repositorio: `microsoft/durabletask-python`;
- release/tag: `v1.9.0`;
- commit firmado/verificado: `6dfdbac521d9d59d31d7ea975fb669d66f86f4c4`;
- tree Git exacto: `f4fa141c10358729b596fafa302412edf1e26174`;
- archive: 1.418.256 bytes, SHA-256 `c8c941b25804ae9abcc0470f5594cab4929216a981982616344eb2b84a7916e8`;
- licencia MIT: SHA-256 `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`;
- sample standalone `examples/human_interaction.py`: 7.139 bytes, SHA-256 `7f97502962fac35775057651a2081f5e4b7e19bc116dba7145b70c8afb467669`;
- sample Azure Functions `azure-functions-durable/samples/human-interaction/function_app.py`: 1.633 bytes, SHA-256 `42fdd853d67fcec3acc51ee12e6df1f4ac6c221d246f77a24585d023ec3f0e83`.

La API oficial Git devolvió ambos blobs dentro del árbol no truncado. El archive exacto fue descargado y verificado otra vez; los hashes de los dos paths extraídos coincidieron. Ambos sources pasaron `compile(...)` en memoria. El primer `py_compile` de la variante profunda no pudo escribir un `.pyc` por el path Windows; `LIB-FAIL-345` lo separa explícitamente del código upstream.

## Código oficial que sí aporta

El sample standalone:

- modela `Order` y `Approval`;
- llama una actividad para solicitar aprobación;
- crea `approval_event = ctx.wait_for_external_event('approval_received')`;
- crea un timer de 24 horas;
- usa `task.when_any` para decidir aprobación o cancelación;
- reanuda el workflow y ejecuta la actividad posterior sólo si ganó el evento;
- usa `DefaultAzureCredential` cuando el endpoint DTS es HTTPS;
- demuestra `raise_orchestration_event` desde el cliente.

La variante Azure Functions demuestra el mismo patrón con `create_check_status_response`, un endpoint para elevar el evento, timeout y cancelación del timer cuando gana la aprobación.

## Ejecución local del pack endurecido

La biblioteca no instaló nada globalmente. En un venv temporal CPython 3.12.13 se instaló el lock de seis wheels con `--require-hashes`; `pip check` pasó. La suite materializada ejecutó 13/13 pruebas con el backend in-memory oficial Microsoft:

- evento de aprobación hash-bound y persistencia una sola vez;
- rechazo humano sin persistencia;
- timeout sin persistencia;
- evento adelantado bufferizado antes de que el orchestrator empiece a esperarlo;
- efecto ambiguo de persistencia sin retry ciego;
- fallo de evidencia final que conserva la transacción comprometida;
- gates de contrato, idempotencia, paths y claves extra.

Los warnings de shutdown/event-loop ya registrados para el backend de test se conservaron. Esta prueba no convierte el backend in-memory en servicio productivo.

## Límite que permanece abierto

Los samples no verifican quién eleva el evento. El standalone acepta un nombre de CLI; la variante Azure declara `DFApp(http_auth_level=func.AuthLevel.ANONYMOUS)` y entrega `req.get_json()` directamente al evento. Tampoco implementan tenant/resource authorization, nonce/replay, idempotency key de decisión, firma/hash del resultado revisado, bandeja/claim de trabajo, separación de funciones, auditoría inmutable ni reconciliación con negocio.

Esta condición queda en `UP-FAIL-104`. La respuesta correcta no es descartar Durable Task: es conservar el SDK/patrón admitido y bloquear la elevación del evento hasta que un adapter de identidad/revisión elegido por el proyecto demuestre esos controles.

## Integración en la biblioteca

- Los dos paths y hashes se añaden como `additional_artifacts` al source lock del mismo archive Microsoft; no se crea una fuente duplicada.
- `MICROSOFT-DURABLE-DOCUMENT-ORCHESTRATION` conserva el código local endurecido como `AUTHORED` y registra las muestras exactas como autoridad upstream, sin relabeling.
- Ningún endpoint anónimo se materializa ni se habilita por defecto.

## Fuentes oficiales

- <https://github.com/microsoft/durabletask-python/tree/6dfdbac521d9d59d31d7ea975fb669d66f86f4c4/examples/human_interaction.py>
- <https://github.com/microsoft/durabletask-python/tree/6dfdbac521d9d59d31d7ea975fb669d66f86f4c4/azure-functions-durable/samples/human-interaction/function_app.py>
