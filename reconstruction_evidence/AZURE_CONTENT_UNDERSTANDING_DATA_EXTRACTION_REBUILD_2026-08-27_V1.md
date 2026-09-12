# Azure Content Understanding Data Extraction Rebuild 2026-08-27 V1

## Resultado

La solución oficial [`Azure-Samples/data-extraction-using-azure-content-understanding`](https://github.com/Azure-Samples/data-extraction-using-azure-content-understanding/tree/0461a5549104ca769b8ec082c05894468997f300) queda fijada como fuente Microsoft/Azure auténtica pero clasificada **`SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`**. No se copió ningún componente de este repositorio a un runtime Elite ni se lo presentó como código productivo.

El rechazo no contradice que la fuente sea oficial o valiosa como arquitectura de muestra. Significa que el commit exacto no satisface la frontera solicitada para recibir archivos empresariales, almacenar datos automáticamente o desplegarse sin adaptación. Su propio README declara que es una demostración y exige adaptar seguridad, monitoreo y manejo de errores para producción.

## Identidad, licencia y actualidad

- repositorio oficial: `Azure-Samples/data-extraction-using-azure-content-understanding`;
- commit: `0461a5549104ca769b8ec082c05894468997f300`;
- GitHub API al 2026-08-27: sigue siendo la cabeza exacta de `main`; repositorio activo, no archivado ni deshabilitado;
- verificación del commit: `verified=true`, razón `valid`; fecha del commit 2025-07-25;
- releases oficiales: respuesta HTTP 200 con JSON `[]`; no existe release/tag publicable que reemplace este snapshot;
- ZIP: 808.247 bytes, SHA-256 `23fb059d172cc9ebbe76cebc5461cf9150222fbc91830f35c0519ddaf2b7f7cc`;
- MIT `LICENSE`: SHA-256 `10d8ad187611652f75ab44ba79ad4d5c118f94746b03266a4b102fc38a2e35c5`;
- inventario del archive: 191 archivos, 87 Python, 55 Terraform y 35 bajo rutas de tests.

El repositorio no contiene workflow bajo `.github/workflows`, lock Python integral ni `.terraform.lock.hcl`.

## Reconstrucción de la suite

Se usaron `uv 0.11.30` oficial, verificado por su checksum publicado, y CPython 3.12.13 administrado por `uv`. Se instaló exactamente lo declarado por `requirements.txt` y `requirements_dev.txt` en un entorno nuevo.

Primer resultado, sin reparar upstream:

```text
ModuleNotFoundError: No module named 'parameterized'
1 error during collection
```

`requirements_dev.txt` no declara el paquete que importa `tests/models/test_extracted_collection_documents.py`. Después de agregar **sólo al entorno de auditoría** `parameterized==0.9.0`, sin editar source ni tests Microsoft:

```text
166 passed, 19 warnings in 6.65s
```

Este PASS mide la suite unitaria disponible. No prueba upload hostil, autorización, Azure real, aislamiento, exactitud de campos, carga, restauración ni despliegue.

## Grafo y seguridad conocida

Los manifests usan pins parciales y rangos; no existe lock integral. La resolución fechada produjo 133 paquetes. El freeze de auditoría tiene SHA-256 `e0c13bbd57e409c434313455191605e71b512439616abfc9bf179362c0e5b032` y no se incorpora como lock Microsoft ni como recomendación de actualización.

Google OSV-Scanner 2.5.1 oficial analizó ese freeze. Reporte JSON SHA-256 `8638d4418cd9ab51ddfa09d3dc714d2395384dd7254f622b31c8f7b3a0dbe357`, exit 1: 17 registros OSV con ID sobre tres package-version fijados/resueltos:

- `cryptography==44.0.1`: 11 registros, incluidos aliases de CVE/GHSA/PYSEC;
- `semantic-kernel==1.22.0`: 4 registros;
- `pymongo==3.12.3`: 2 registros.

Los registros pueden representar el mismo advisory bajo aliases distintos; no se los infla como 17 vulnerabilidades independientes. Tampoco se realizó reachability productiva porque no existe un target admitido. Los findings, el grafo móvil y la falta de lock bastan para impedir adopción inmediata.

## Frontera HTTP y archivo

`src/function_app.py` crea la aplicación con `func.AuthLevel.ANONYMOUS` y registra configuración, consulta e ingesta. No se encontró verificación server-side de token/claims/roles en `src`. La consulta confía en un header `x-user` aportado por el cliente. Los tests que describen un usuario "authorized" asignan `enable_authorization` a un mock, pero la ruta no consume esa propiedad.

La ruta de ingesta obtiene `req.get_body()` y sólo rechaza cuerpo vacío. Antes de enviarlo a Content Understanding no demuestra:

- límite de bytes o rate limit;
- tipo real por contenido y allowlist;
- antimalware, YARA, quarantine y release;
- nombre/tenant normalizados;
- autenticación/autorización de la carga.

Por eso el upload no puede sustituir el gate local seguro ya separado en Elite.

## Exactitud, persistencia y concurrencia

El controlador llama al analyzer/classifier y entrega inmediatamente la respuesta al servicio que escribe markdown en Blob y campos en Cosmos/Mongo. `_process_extracted_field` acepta un campo configurado si contiene alguna clave `value*`; no aplica umbral de confianza, campos críticos, schema cerrado recursivo, verificación cruzada, corpus ground-truth ni aprobación humana antes del upsert.

Dos probes externos, sin modificar source Microsoft, demostraron:

```text
LOCK_TIMEOUT_IGNORED_PROVEN
ZERO_CONFIDENCE_FIELD_ACCEPTED_PROVEN
```

El primero hace que `MongoLockManager.wait()` devuelva `False`: `ingest_analyzer_output()` ignora el retorno y llega igualmente a `_upsert_document`. El segundo prueba que un `invoice_total` con `confidence: 0.0` y valor `999999.99` se acepta y queda agregado al modelo que se persiste.

Existe deduplicación por collection/config hash/lease/path, pero no se demostró transacción, outbox o compensación entre Blob y Cosmos. Esa deduplicación parcial no convierte la escritura multi-store en una operación empresarial consistente.

## API, runtime, observabilidad e infraestructura

- el cliente Python propio fija `_DEFAULT_API_VERSION = "2025-05-01-preview"`;
- el [quickstart REST oficial vigente de Microsoft](https://learn.microsoft.com/en-us/azure/ai-services/content-understanding/quickstart/use-rest-api) respondió HTTP 200 y declara `2025-11-01`; el source no usa el SDK GA `azure-ai-contentunderstanding==1.1.0` ya verificado separadamente;
- README exige Python 3.12 o posterior y el devcontainer usa Python 3.12 móvil, pero Terraform despliega Function Python 3.11;
- `set_up_monitoring()` está comentado en la entrada de la Function;
- el registro de las rutas classifier también está comentado aunque el README anuncia clasificación;
- Terraform habilita acceso público en Key Vault, AI/OpenAI, Cosmos, Function y Storage; varios ACL usan `default_action = "Allow"`;
- aparecen `Key Vault Administrator`, `DocumentDB Account Contributor` y `Storage Blob Data Owner`, sin demostración target de least privilege;
- el quick start propone `curl` desde `raw.githubusercontent.com/.../main/deploy.sh | bash`, una referencia móvil ejecutada por pipe, no el commit/hash fijado.

## Decisión de reutilización

No se aisló un componente productivo listo de este repositorio. El cliente usa preview; la frontera HTTP no es admisible; la persistencia no tiene puerta de exactitud; el lock falla; observabilidad está desactivada y el baseline IaC es de demostración.

Se conserva el archive exacto en el lock de fuentes para estudio y trazabilidad, con la condición de rechazo visible antes de adquirir. Para una selección Azure, el agente debe usar el SDK GA oficial ya fijado, los gates oficiales separados de tipo/antimalware/evaluación y una integración de proyecto declarada como `AUTHORED`; no puede atribuir esa integración a Microsoft ni saltar cuentas, costos, corpus, schemas y pruebas reales.

## Condiciones de reapertura

Reauditar sólo una revisión oficial posterior que, como mínimo:

1. use API/SDK GA y runtime soportado coherente;
2. publique lock integral, CI y suite reproducible sin dependencia omitida;
3. cierre findings conocidos o documente fix/reachability oficial verificable;
4. implemente authn/authz, límites, tipo real, antimalware y cuarentena;
5. haga fail-closed el lock y pruebe concurrencia/idempotencia multi-store;
6. impida persistencia de confianza cero y demuestre schema/evidencia/corpus/review;
7. habilite monitoreo seguro, redes privadas/least privilege, backup/restore y rollback;
8. pase tests offline, Azure sandbox, corpus adversarial, carga y seguridad.

Hasta entonces, el estado global continúa **`NOT_READY_UNDER_EXPANDED_USER_STANDARD`**.
