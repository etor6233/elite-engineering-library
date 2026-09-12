# Evidencia V246 — lifecycle limpio del worker Microsoft Durable Task

Fecha: 2026-09-04

## Hallazgo

El Audit V245 hizo visible `ResourceWarning: unclosed event loop` y sockets gRPC sin cerrar al crear/detener repetidamente `TaskHubGrpcWorker` de `durabletask==1.9.0` sobre Windows/CPython 3.12. El exit 0 previo no se aceptó como prueba de cleanup.

La fuente instalada exacta en `durabletask/worker.py` demuestra que `TaskHubGrpcWorker.start` crea `asyncio.new_event_loop()` dentro del thread y llama `run_until_complete(self._async_run_loop())`, pero no ejecuta `loop.close()` al retornar. El artefacto sigue fijado a wheel SHA-256 `a759af4ad8e6897922575886e93bf40c6b3af936eaed83798d492ee561607185` y commit Microsoft `6dfdbac521d9d59d31d7ea975fb669d66f86f4c4`.

## Corrección atribuida

`MICROSOFT-DURABLE-DOCUMENT-ORCHESTRATION 0.1.2` incorpora `ManagedTaskHubGrpcWorker` dentro de `documentflow/runtime.py`. El bloque completo se clasifica `ADAPTED`, no Microsoft `VERBATIM` ni `AUTHORED`: conserva la secuencia oficial del worker fijado y añade, después del retorno del loop, `shutdown_asyncgens`, `shutdown_default_executor`, `set_event_loop(None)` y `loop.close()`.

El test usa ese lifecycle en vez de instanciar directamente el worker oficial. `VERIFY_EXECUTABLE_LIBRARY.ps1` exige los controles de cleanup y ejecuta la suite con `-W error::ResourceWarning`.

## Resultado

| Gate | Resultado |
|---|---|
| materialización canónica | 13/13 archivos |
| contratos + pipeline | 13/13 PASS |
| pipeline focal | 8/8 PASS |
| política de warnings | `ResourceWarning` elevado a error, cero event loops o sockets reportados |
| compileall | PASS |
| root posterior | `VERIFY_LIBRARY_PASS`: 153 packs, 1.356 bloques, 666 Markdown antes de este expediente |
| Audit global posterior | `PASS`: 153 packs/121 fuentes/16 adapters; `ResourceWarning` fue error y no apareció |

El warning de log gRPC `StatusCode.CANCELLED` durante `stop()` es emitido por el SDK oficial como parte de su apagado y es seguido por `Worker shutdown completed`; no es el `ResourceWarning` de recursos sin cerrar corregido aquí. Backend productivo DTS/Azure, identidad, carga, recovery y operación target permanecen condicionados.
