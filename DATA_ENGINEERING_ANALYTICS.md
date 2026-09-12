# Data Engineering & Analytics — manual operativo

> **Estado:** núcleo general v1 auditado el 2026-08-21; fuentes, integridad documental y cruces recíprocos cerrados. Contratos, gates y capstones de planificación disponibles en `ENGINEERING_EXECUTION_PLAYBOOK.md`; pipelines concretos requieren evidencia propia.
> **Autoridad académica:** Stanford CS246 *Mining Massive Data Sets*; Berkeley Data 100; CMU 15-445/645 como frontera de bases de datos.
> **Autoridad técnica:** Apache Beam, Spark, Airflow, Arrow, Parquet e Iceberg; OpenLineage; dbt; Kimball Group.
> **Propósito:** construir productos de datos reproducibles, confiables, gobernados y económicamente operables desde la fuente hasta una decisión, dashboard, API, modelo o archivo.
> **Frontera:** este manual gobierna el ciclo y los contratos del dato. Internals de storage, delivery distribuido, ML y seguridad conservan sus manuales autoridad.

## Cómo usar este manual

```text
pregunta/decisión
→ fuentes y ownership
→ contrato y captura
→ ingestión reproducible
→ transformación por capas
→ calidad + lineage
→ modelo/serving
→ consumo y semántica
→ observación + costo
→ backfill/recovery/evolución
```

No construir un “data lake” antes de saber quién utilizará qué dato, con qué semántica, freshness, calidad, privacidad y costo. Mover bytes no crea información confiable.

## Procedencia

- `[SPEC]`: formato, estándar o documentación oficial versionada.
- `[ACADEMIC]`: síntesis de curso, texto o paper identificado.
- `[CODE]`: implementación/repositorio oficial.
- `[PROD]`: práctica publicada por autores u operadores identificados.
- `[MEASURED]`: medición reproducible dentro del entorno declarado.
- `[IMPL]`: contrato, diseño, código, test o runbook derivado.
- `[OPEN]`: tema no verificado o dependiente del proyecto.

## 1. Contrato de un producto de datos

Un dataset es un producto sólo cuando tiene consumidores, semántica, owner, SLO y lifecycle. Ficha mínima:

```text
nombre y propósito:
decisión/tarea soportada:
owner productor y owner consumidor:
source(s) of truth:
grain/clave y significado de una fila/evento:
schema + unidades + timezone:
freshness y completeness SLO:
quality rules y tolerancias:
history/corrections/deletes:
PII/clasificación/retención:
acceso y tenant scope:
version/compatibility:
lineage y código productor:
cost/capacity:
deprecation/export:
```

### 1.1 Separar niveles de verdad

| Nivel | Ejemplo |
|---|---|
| hecho fuente | orden aceptada por matching engine |
| observación | mensaje recibido por feed handler |
| dato normalizado | evento con schema interno |
| dato derivado | spread o cohort |
| métrica | definición agregada de negocio |
| interpretación | conclusión humana/modelo |

No promover observación incompleta a hecho. No presentar derivado o imputación como fuente cruda.

### 1.2 SLO de datos

- **freshness:** edad desde event/effective time hasta disponibilidad;
- **completeness:** fracción esperada presente;
- **validity:** cumplimiento de tipos/reglas;
- **accuracy:** correspondencia con realidad/oráculo;
- **consistency:** acuerdo entre representaciones;
- **uniqueness:** duplicados bajo clave definida;
- **availability:** consultas/entregas exitosas;
- **lineage coverage:** outputs con origen/código/run identificables.

Un check verde no prueba accuracy si sólo valida formato.

## 2. Diseñar desde el consumo

### 2.1 Preguntas

- ¿La salida alimenta humanos, software, ML o auditoría?
- ¿Necesita detalle, agregado, snapshot o cambios?
- ¿Qué latencia/freshness cambia una decisión?
- ¿Qué errores son tolerables y cuáles causan daño?
- ¿Se reconsulta historia o sólo estado actual?
- ¿Se requiere point-in-time correctness?
- ¿Qué dimensiones, filtros y drill-down necesita?
- ¿Cuánto cuesta equivocarse, llegar tarde o no servir?

### 2.2 Serving patterns

| Consumo | Diseño dominante |
|---|---|
| BI/OLAP | columnar, dimensional/semantic, scan/aggregate |
| API operational | indexed/materialized serving store |
| ML training | versioned point-in-time dataset/features |
| online features | low-latency keyed state + offline parity |
| audit/replay | immutable raw log + schema/metadata |
| exploration | governed sandbox + reproducible notebook/query |

No obligar un único store a servir todos los patrones.

## 3. Source contract y captura

### 3.1 Inventario de fuente

```text
owner y authority
interface: table/API/file/log/CDC
schema y identifiers
event/effective/ingest time
ordering y delivery
update/delete/correction semantics
snapshot consistency
rate limits y retention
backfill/export
security/privacy
known defects
```

### 3.2 Extracción

| Modo | Ventaja | Riesgo |
|---|---|---|
| full snapshot | simple/oráculo | costo, ventana inconsistente |
| high-water mark | incremental | late update y cursor incorrecto |
| CDC/log | conserva cambios | bootstrap/cut, schema y retention |
| API pagination | contrato externo | drift, rate, cursor expiry |
| files/drop | desacoplado | partial upload, duplicates, discovery |
| events | baja latencia | delivery/order/replay |

Cada extractor guarda cursor/checkpoint junto con el output que confirma. No avanzar posición antes de que el dato sea recuperable según el contrato.

### 3.3 Bootstrap de CDC

Snapshot y log deben formar un corte consistente:

```text
capturar posición L
obtener snapshot asociado a L
cargar snapshot idempotentemente
reproducir cambios > L
deduplicar/ordenar según source semantics
verificar reconciliation
```

Retention del log debe superar peor tiempo de snapshot + outage + recovery. Detalles de WAL/MVCC en `DATABASE_STORAGE_INTERNALS.md`.

## 4. Raw no significa sin contrato

Conservar evidencia fuente con mínima transformación necesaria para recepción segura:

- payload y headers relevantes;
- source, partition, offset/sequence;
- event time e ingest time;
- schema/codec version;
- checksum/size;
- capture run y code version;
- rechazo/quarantine reason.

Raw debe ser immutable o versionado. Corregir mediante nueva versión/evento, no reescribir historia sin evidencia. Cifrado, ACL, retención y delete siguen aplicando.

## 5. Batch y streaming como boundedness

**[SPEC]** Beam representa colecciones bounded o unbounded con transforms comunes. Para unbounded data, resultados dependen además de event time, windows, watermarks, triggers, accumulation y allowed lateness.

### 5.1 Preguntas universales

1. ¿Qué se calcula?
2. ¿Dónde se agrupa? — key/window.
3. ¿Cuándo se emite? — trigger.
4. ¿Cómo cambian resultados previos? — discard/accumulate/retract/upsert.
5. ¿Cuándo se elimina state? — lateness/TTL.

### 5.2 Tiempos

| Tiempo | Significado |
|---|---|
| event/effective | cuando ocurrió o aplica en dominio |
| ingest | cuando entró al pipeline |
| processing | reloj del operador |
| availability | cuando consumidor pudo usarlo |

No sustituir uno por otro silenciosamente. Watermark es una estimación de progreso, no prueba de que jamás llegue un evento anterior.

### 5.3 Windows

- fixed: intervalos no superpuestos;
- sliding: intervalos solapados;
- session: actividad separada por gap;
- global: exige trigger/state explícitos.

Declarar boundaries, timezone/DST, inclusividad, late updates y cómo verá correcciones el sink.

### 5.4 Unified logic, different operation

Compartir semántica batch/stream reduce divergencia, pero ejecución, costo, recovery y sink pueden ser distintos. Reprocesar historia debe producir resultado equivalente bajo el mismo contrato o explicar la diferencia.

## 6. Transformaciones correctas

### 6.1 Propiedades deseables

- deterministic para misma versión/input/config;
- idempotent bajo retry;
- partition-addressable;
- side effects aislados;
- schema explícito;
- output publicado atómicamente;
- code/config/dependencies versionados;
- observabilidad por run/partition.

### 6.2 Capas

Una convención útil, no universal:

```text
source/raw: evidencia recibida
staging: tipos/nombres/normalización mínima
core: entidades y hechos conformados
mart/serving: consumo específico
```

Cada capa añade semántica demostrable. Evitar copias sin owner o tablas intermedias consumidas accidentalmente como contrato público.

### 6.3 SQL como software

- format/lint y review;
- modelos pequeños con nombres semánticos;
- columnas explícitas en interfaces estables;
- joins con cardinalidad esperada;
- null/three-valued logic consciente;
- timezone y decimals definidos;
- tests por input/output;
- query plan/costo medido;
- macros sin ocultar lógica crítica.

## 7. Idempotencia y publicación

**[SPEC]** Las best practices de Airflow tratan tasks como transacciones: retries no deben dejar resultados parciales y un rerun debe producir el mismo outcome. Recomiendan particiones específicas en lugar de “latest” o `now()` para cómputo crítico.

Patrones:

- escribir staging con run ID y hacer atomic rename/swap;
- MERGE/UPSERT bajo key correcta;
- replace partition, no append ciego;
- manifest/commit marker después de checks;
- optimistic commit con snapshot/table format;
- sink idempotent o dedupe retention suficiente.

Un task exitoso no equivale a pipeline completo publicado. Separar compute, validate y promote.

## 8. Schemas y compatibilidad

### 8.1 Contrato de campo

```text
nombre estable
tipo físico y lógico
nullable/required
unidad/scale/precision
timezone
domain/enum
default y missing semantics
PII/classification
owner/descripción
```

### 8.2 Cambio

Clasificar add/remove/rename/type/nullability/semantic change. La compatibilidad depende de readers/writers y del formato. Un cambio físicamente legible puede ser semánticamente breaking.

Proceso:

1. inventariar consumidores por lineage + telemetry;
2. agregar y dual-write/derive;
3. backfill y comparar;
4. migrar consumidores;
5. bloquear uso viejo;
6. retirar con evidencia.

### 8.3 Schema drift

No aceptar nuevas columnas/tipos en producción por defecto sin policy. Opciones: reject, quarantine, allow additive, map unknown, version dataset. Alertar por semántica, no sólo parser failure.

## 9. Formatos: fila, columna y memoria

### 9.1 Selección

| Necesidad | Candidato |
|---|---|
| intercambio humano/pequeño | CSV/JSON con schema externo y límites |
| mensajes tipados | Avro/Protobuf según evolution/ecosystem |
| analítica persistente | Parquet/ORC |
| memoria/interprocess columnar | Arrow |
| tabla sobre object storage | Iceberg/Delta/Hudi según contrato/ecosistema |

No confundir file format, table format, catalog, compute engine y storage service.

### 9.2 Parquet

**[SPEC]** Un archivo Parquet organiza row groups; cada row group contiene un column chunk por columna, dividido en pages. File/row group, column chunk y page son unidades diferentes de paralelización, I/O y encoding/compression.

Decisiones:

- row group/file size y small-files problem;
- sort/cluster para pruning;
- statistics/bloom availability;
- codec y CPU/ratio;
- logical types;
- reader/writer compatibility;
- encryption/footer support;
- schema evolution real.

### 9.3 Arrow

**[SPEC]** Arrow define layout columnar en memoria, metadata e IPC language-agnostic, con adjacency, random access generalmente constante, SIMD friendliness y acceso relocatable/zero-copy bajo condiciones. No es por sí mismo base de datos ni protocolo de mutación.

Validar null bitmap, offsets, dictionary lifecycle, endianness, alignment, ownership/lifetime y conversiones de tipo.

## 10. Table formats y lakehouse

Un table format añade metadata/commit semantics sobre archivos:

- snapshots y time travel;
- atomic commits;
- schema/partition evolution;
- manifests/statistics;
- concurrent writers;
- delete/update representation;
- maintenance/expiration.

**[SPEC]** Iceberg separa partitioning lógico de layout visible, soporta evolution y snapshots; la compatibilidad real depende de catálogo, engine, versión y feature usada.

### 10.1 Operación obligatoria

- compact small files;
- rewrite manifests/data/delete files;
- expire snapshots sin romper readers/rollback;
- remove orphans con safety window;
- monitor metadata growth;
- reconcile catalog/storage;
- test concurrent commits;
- conservar legal hold/audit/retention.

Object storage no convierte archivos en tabla transaccional automáticamente.

## 11. Particionamiento, clustering y skew

Elegir por filtros, joins, maintenance y write pattern. Una partición demasiado fina crea archivos pequeños/metadata; demasiado gruesa lee datos inútiles.

Registrar:

```text
partition transform/key
cardinality y distribution
files/partition y size distribution
top keys/skew
pruning effectiveness
shuffle bytes/spill
evolution strategy
```

Hashing distribuye bajo supuestos; no elimina hot keys. Salting puede repartir cómputo y complica agrupación/recomposición.

## 12. Procesamiento distribuido

**[ACADEMIC]** Stanford CS246 usa MapReduce/Spark para razonar sobre algoritmos de datos masivos, incluyendo itemsets, LSH, grafos, clustering, recomendadores y streams. La lección transferible es diseñar alrededor de partición, shuffle, comunicación y memoria, no sólo traducir un loop.

### 12.1 Costos

- input bytes leídos;
- shuffle bytes y records;
- serialization/deserialization;
- spill y disk I/O;
- task startup/scheduler;
- skew/stragglers;
- memory per executor;
- output files/commits;
- retries/speculation.

### 12.2 Transformación

Preferir projection/filter temprano, combiners/partial aggregation cuando sean algebraicamente válidos, broadcast sólo bajo tamaño acotado, y joins con estrategia basada en cardinalidad medida.

No usar `groupByKey` donde una combinación incremental conserva la semántica. No `collect` al driver sin bound.

## 13. Spark como engine, no semántica universal

Spark 4.2.0 era la documentación estable observada el 2026-08-21. DataFrames/SQL permiten optimización sobre schema; RDD/API/streaming difieren en contratos y operación.

Checklist:

- plan lógico/físico y statistics;
- partition count/size;
- adaptive behavior y version;
- join strategy/skew;
- cache/persist lifecycle;
- serialization/UDF boundary;
- spill/GC/executor loss;
- checkpoint y sink semantics;
- driver/executor capacity;
- output file layout.

Un job verde puede producir datos duplicados, stale o semánticamente erróneos.

## 14. Orquestación

El orchestrator coordina trabajo; no debe contener toda la transformación ni mover datasets grandes por su metadata DB.

### 14.1 DAG contract

```text
logical date/data interval
inputs versionados
outputs/partitions
dependencies por dato, no sleeps
retry/idempotencia
timeout/SLA/freshness
backfill/catchup
concurrency/resource pool
failure/cancel behavior
owner/runbook
```

### 14.2 Airflow

Airflow 3.3.1 separa scheduler, DAG processor, API server, metadata DB y workers/triggerer opcionales. DAG code puede ejecutarse repetidamente durante parsing: evitar I/O/cómputo top-level. Tasks distribuidas no comparten filesystem local confiable.

No usar XCom para payloads grandes. No guardar secrets en DAGs. Pin de DAG bundle/code version y run ID forman parte de reproducibilidad.

### 14.3 Backfill

Un backfill es una migración de datos:

- rango/particiones y code version;
- inputs congelados o correcciones conocidas;
- isolated compute/quotas;
- output namespace/version;
- validation/diff;
- promotion/cutover;
- rollback y cleanup;
- efecto sobre downstream y costo.

## 15. Data quality por capas

### 15.1 Gates

| Capa | Checks |
|---|---|
| ingest | bytes/checksum/schema/sequence/row count |
| staging | types, parse, domain, nulls |
| core | key, uniqueness, referential/invariants |
| mart | grain, totals, reconciliation, metric semantics |
| serving | freshness, availability, latency, access |

### 15.2 Tipos de test

- schema/contract;
- row/column constraints;
- property/invariant;
- reconciliation contra source/control totals;
- differential contra versión anterior;
- distribution/anomaly con threshold justificado;
- freshness/completeness;
- point-in-time leakage;
- privacy/access;
- recovery/backfill.

### 15.3 Fail, quarantine o warn

La acción depende del daño:

- **fail closed:** corrupción peligrosa o publicación contractual;
- **quarantine:** subset identificable sin contaminar output;
- **warn/degrade:** señal no crítica con owner y deadline;
- **continue with annotation:** consumidor acepta incompletitud explícita.

Nunca descartar silenciosamente records inválidos.

## 16. Observability y lineage

### 16.1 Tres planos

- **pipeline:** run/task, duration, retry, resource/cost;
- **data:** volume, schema, quality, freshness, distribution;
- **business:** métricas/invariantes de dominio.

### 16.2 OpenLineage

**[SPEC]** OpenLineage 1.52.0 observado modela Job, Run y Dataset; eventos runtime y design-time pueden portar facets de schema, source code, version, quality y statistics. Lineage describe cómo surge el dataset, no prueba que sea correcto.

Identity estable necesita namespace/name/version/run ID bien definidos. Capturar column lineage sólo donde precisión y costo justifican.

### 16.3 Impact analysis

Antes de cambiar:

1. recorrer consumidores downstream;
2. clasificar contractual/exploratorio/abandonado;
3. verificar uso real;
4. simular schema/metric change;
5. coordinar migration;
6. observar después del cutover.

Lineage incompleto debe declararse, no asumirse completo.

## 17. Modelado dimensional

**[PROD]** Kimball comienza por proceso de negocio, grain, dimensions y facts. El grain se declara antes de elegir medidas; mezclar grains crea agregados incorrectos.

### 17.1 Facts

- transaction fact: un evento;
- periodic snapshot: estado por intervalo;
- accumulating snapshot: lifecycle actualizado;
- factless fact: ocurrencia/cobertura sin medida numérica.

Medidas additive, semi-additive o non-additive según dimensiones. Ratios se calculan desde numerador/denominador, no promediando promedios.

### 17.2 Dimensions

- conformed dimensions entre procesos;
- surrogate key cuando historia/integración lo requieren;
- degenerate dimension para identifier sin tabla propia;
- role-playing dimension;
- slowly changing dimension con política explícita.

SCD:

| Tipo | Semántica |
|---|---|
| 1 | sobrescribe; sin historia |
| 2 | nueva fila/version con vigencia |
| 3 | conserva valor anterior limitado |

No usar Type 2 sin intervalos, current flag consistente y join point-in-time.

## 18. Métricas y capa semántica

Contrato de métrica:

```text
nombre/owner/propósito
entidad y grain
numerador/denominador
filtros e inclusiones
time dimension/timezone
aggregation y dimensions válidas
late/correction policy
source models/lineage
tests/reconciliation
version/deprecation
```

“Revenue”, “active user” o “latency” sin definición versionada genera múltiples verdades. La semantic layer ayuda a reutilizar definición; no resuelve desacuerdo de negocio por sí misma.

## 19. Analytics y EDA reproducible

**[ACADEMIC]** Berkeley Data 100 organiza trabajo alrededor de un ciclo iterativo con adquisición, limpieza, exploración, visualización, inferencia/modelado y comunicación, y exige pensamiento crítico sobre datos y efectos humanos.

### 19.1 Notebook no es producto final

Un análisis entregable conserva:

- pregunta e hipótesis;
- snapshot/query y permisos;
- environment/dependencies;
- cleaning decisions;
- exclusions/missingness;
- code/seed;
- figures/tables generadas;
- uncertainty y limitations;
- revisión y fecha.

Promover lógica estable del notebook a módulos/SQL/tests. El notebook sigue como narrativa reproducible, no como scheduler oculto.

### 19.2 Evitar conclusiones falsas

- selection/survivorship bias;
- leakage y future information;
- Simpson’s paradox;
- multiple comparisons;
- seasonality y nonstationarity;
- missing not at random;
- instrumentation change;
- causal claim desde correlación;
- averages que ocultan tails/slices.

## 20. Privacidad, seguridad y gobierno

Por dataset/campo:

- purpose y lawful/authorized use;
- classification/PII;
- least-privilege access;
- encryption/key lifecycle;
- row/column/tenant policy;
- masking/tokenization limits;
- retention/delete/legal hold;
- export/share constraints;
- audit y anomaly detection.

Una copia derivada conserva obligaciones aunque cambie nombre o formato. Lineage ayuda a localizar copias; no ejecuta deletion ni autorización.

### 20.1 Entornos

No copiar producción completa a development. Usar synthetic, sampled/minimized o masked data con evaluación de reidentification y utilidad. Separar service accounts, buckets/schemas, catalogs y egress.

## 21. CI/CD de datos

### 21.1 PR gate

- parse/compile/lint/type;
- unit SQL/transform;
- schema/contract diff;
- lineage/impact;
- sample/ephemeral build;
- data tests;
- query plan/cost guard;
- security policy;
- docs/owner.

### 21.2 Deployment

Promover mismo código/artefacto; configuración y secrets externos. Para breaking evolution usar dual outputs, shadow comparison y consumer migration. Rollback de código puede no revertir datos escritos: diseñar forward fix, restore o versioned output.

### 21.3 Environments y state

CI necesita datasets pequeños representativos y namespaces aislados. No permitir que tests escriban producción. Conservar fixtures con edge cases y contracts; rotar si contienen datos reales.

## 22. Operación y SRE de datos

SLIs:

- freshness lag/event age;
- success/good partitions;
- completeness/control totals;
- invalid/quarantined rate;
- runtime/queue/backlog;
- cost per run/GB/consumer;
- small files/metadata growth;
- downstream impact.

Alertar sobre síntomas que requieren acción, con dataset/partition/run/code version/owner. Un “DAG failed” sin impacto/freshness ni runbook es insuficiente.

### 22.1 Incident response

```text
detectar y detener publicación/consumo peligroso
preservar raw evidence y lineage
identificar ventana/datasets/consumidores
comunicar calidad/freshness real
reparar o versionar output
backfill/reconcile
validar downstream
retirar datos incorrectos/caches
postmortem y control preventivo
```

## 23. Costo y FinOps

Descomponer:

- ingest/network egress;
- stored bytes × replicas/versions;
- compute scan/shuffle/spill;
- orchestrator/catalog/observability;
- maintenance/backfill;
- idle clusters/warehouses;
- human operation.

Optimizar después de atribuir costo por dataset/job/team/consumer. Medidas:

- projection/pruning;
- file sizing/compaction;
- incremental computation correcto;
- materialization/caching con reuse real;
- autosuspend/right-size;
- retention/snapshot expiry;
- evitar duplicación sin consumidores.

No bajar costo rompiendo reproducibilidad, recovery o calidad.

## 24. Real-time: decidir por valor

Preguntar:

```text
¿qué decisión mejora entre T+1 día, T+1 hora, T+1 min y T+1 s?
¿event time o processing time?
¿qué late correction acepta el consumidor?
¿cómo se repara historia?
¿qué costo/on-call agrega?
```

Near-real-time microbatch suele ser suficiente. Streaming continuo se justifica cuando freshness tiene valor medible y el equipo puede operar state, backlog, replay y schema evolution.

## 25. Datos para ML e IA

Añadir:

- dataset/version/split lineage;
- label definition y delay;
- point-in-time joins;
- entity keys;
- leakage tests;
- training-serving parity;
- slice coverage;
- consent/retention;
- feature freshness;
- feedback contamination;
- model/eval version.

Este manual entrega datos confiables; `ML_PRODUCTION_LLMOPS_EVALUATION.md` decide ciclo de modelo/evals y `NLP_RAG_RETRIEVAL_DATA.md` retrieval/RAG.

## 26. Order books y datos temporales

Separar:

- raw venue frames;
- decoded normalized events;
- validated local-book state;
- snapshots/derived metrics;
- analytical aggregates.

Cada output porta venue/instrument/epoch/sequence/event/receive time/freshness. No rellenar gaps con forward-fill como si fueran observaciones. Usar replay determinista y reconciliation contra snapshot/checksum según protocolo. El hot path pertenece a sistemas/red/microestructura; el warehouse no controla el matching state.

## 27. Gate profesional

- [ ] Consumidor, decisión, owner y source of truth.
- [ ] Grain, keys, schema, unidades y tiempos.
- [ ] SLO de freshness/completeness/quality.
- [ ] Capture/cursor/bootstrap/replay.
- [ ] Idempotencia y atomic publish.
- [ ] Batch/stream/window/late semantics.
- [ ] Partition/file/table/engine con costo medido.
- [ ] Quality gates y quarantine/reconcile.
- [ ] Lineage code→run→dataset→consumer.
- [ ] Privacy/access/retention/delete.
- [ ] Backfill/migration/rollback.
- [ ] Metrics contract y point-in-time correctness.
- [ ] Observability/runbook/cost.
- [ ] Tests con edge, skew, late, duplicate y failure.

## 28. Contrato para Codex

```text
Objetivo/consumidor:
Fuentes/authority:
Grain/keys/schema/time:
Volumen, rate, retention y freshness:
Corrections/deletes/late/order:
Quality/privacy/access:
Stack y restricciones:

Entregar:
1) contratos de source y dataset;
2) arquitectura por capas y ownership;
3) batch/stream semantics;
4) schema/evolution y file/table layout;
5) transforms idempotentes con atomic publish;
6) tests de schema, invariantes, reconciliation y point-in-time;
7) lineage/observability/SLO;
8) backfill/recovery/migration;
9) capacity/cost y benchmark;
10) riesgos, unknowns y runbook.

No usar latest/now implícito en cómputo reproducible.
No afirmar exactly-once end-to-end por garantía local.
No descartar datos inválidos silenciosamente.
No llamar raw a datos sin provenance.
No definir métricas sin grain, tiempo y agregación.
```

## 29. Suites de aceptación

| Suite | Inyección | Aceptación |
|---|---|---|
| `source_capture` | pagination drift, late update, CDC disconnect | sin pérdida/duplicado fuera de contrato; cursor recuperable |
| `schema_evolution` | add/remove/type/null/semantic | consumers compatibles o bloqueo/migración explícita |
| `idempotent_backfill` | retry/crash/partial publish | mismo output/version; sin parciales visibles |
| `time_windows` | late/out-of-order/DST | panes/correcciones según contrato |
| `quality_reconcile` | null/duplicate/skew/control-total mismatch | fail/quarantine/alert correcto |
| `lineage_impact` | cambio upstream | consumidores y run/code/data versions identificables |
| `privacy_delete` | delete/access revoke | propagación demostrada dentro de SLA |
| `cost_scale` | 10× bytes/skew/small files | SLO/costo/headroom medidos |
| `restore` | metadata/catalog/data loss | snapshot/cut correcto y serving validado |

## 30. Capstones

### A. Plataforma batch reproducible

API/DB → raw immutable → Parquet/Iceberg → transforms → dimensional mart → dashboard. Exigir schema evolution, idempotent backfill, quality/lineage, privacy y cost report.

### B. Pipeline streaming

Events → Beam/Flink semantics → windows/state → upsert sink. Inyectar late, duplicate, gap, restart y sink ambiguity. Comparar replay batch.

### C. Feature pipeline ML

Offline/online feature con point-in-time join, version, parity, freshness y deletion. Entrenar/evaluar y reproducir dataset exacto.

### D. Market-data analytics

Raw frames y book replay separados de marts. Calcular métricas no predictivas con sequence/freshness y demostrar que un gap invalida outputs hasta resync.

## 31. Inventario de fuentes y cobertura

### 31.1 Academia

| Fuente | Cobertura |
|---|---|
| Stanford CS246 2026 | MapReduce/Spark, itemsets, LSH, recommenders, grafos, large-scale ML, streams |
| Berkeley Data 100 | lifecycle, cleaning/EDA, sampling, visualization, modeling, inferencia y ética |
| CMU 15-445/645 | frontera: storage, execution, transactions, recovery |

### 31.2 Especificaciones/proyectos

| Fuente | Cobertura auditada |
|---|---|
| Beam Programming Guide | pipeline/PCollection/transforms/I/O/schema/coders/windows/triggers/metrics/state/timers/SDF/multilanguage/batching |
| Spark 4.2.0 | SQL/DataFrame engine y ejecución distribuida |
| Airflow 3.3.1 | architecture, DAG/tasks, idempotencia, parsing, testing y operación |
| Arrow 25.0.1 docs / format 1.5 | layout columnar, tipos, buffers, IPC y compatibility |
| Parquet | row groups/chunks/pages, types, metadata, codecs y format compatibility |
| Iceberg latest observado | snapshots, schema/partition evolution, reliability y maintenance |
| OpenLineage 1.52.0 | Job/Run/Dataset, runtime/design events y facets |
| dbt current docs | structure, metrics, materialization, idempotence, tests y workflows |
| Kimball Group | grain, facts, dimensions y dimensional techniques |

Las versiones son observaciones al 2026-08-21, no promesa futura. Todo proyecto fija engine/format/catalog/provider y ejecuta compatibility tests reales.

### 31.3 Práctica pública Uber/Airbnb — datos operados como producto

**[PROD]** Uber documenta que cada artefacto de datos necesita owner, propósito, tier, calidad/SLA, incidentes y deprecation; su catálogo une metadata básica, lineage, uso, quality checks, costo y bugs. Distingue **freshness** de **completeness**: llegar temprano con una fracción inestable puede producir decisiones peores que esperar un watermark defendible.

**[PROD]** Airbnb documenta una capa métrica que centraliza definiciones de metrics/dimensions y selecciona fuentes válidas para distintos consumidores. La práctica no es construir una plataforma gigantesca desde el inicio, sino impedir que dashboard, experimento y finanzas reimplementen silenciosamente grain, filtros, joins y ventanas incompatibles.

**[IMPL] Contrato operacional por dataset/métrica:**

```text
owner + purpose + criticality + consumers
grain + keys + event/effective/processing time
schema + semantic definition + allowed dimensions
freshness target + completeness watermark
accuracy/consistency/uniqueness checks
lineage + code/run revisions + cost/usage
incident/bug SLA + deprecation/migration
```

Gate: recomputar una métrica desde facts certificados; comparar API/dashboard/experiment output; backfill bajo definición versionada; simular late/missing/duplicate/corrected data; demostrar qué consumers quedan afectados mediante lineage; alertar sólo con owner y acción. Un quality score o certificación es evidencia resumida, no reemplaza resultados descompuestos ni adecuación al uso.

Fuentes: [Uber data culture y quality](https://www.uber.com/us/en/blog/ubers-journey-toward-better-data-culture-from-first-principles/) y [Airbnb Minerva](https://medium.com/airbnb-engineering/how-airbnb-achieved-metric-consistency-at-scale-f23cc53dea70).

## 32. Fuentes oficiales verificadas

- Stanford CS246 2026: <https://web.stanford.edu/class/cs246/>
- Berkeley Data 100: <https://ds100.org/course-notes/>
- Apache Beam Programming Guide: <https://beam.apache.org/documentation/programming-guide/>
- Apache Spark: <https://spark.apache.org/docs/latest/>
- Apache Airflow architecture: <https://airflow.apache.org/docs/apache-airflow/stable/core-concepts/overview.html>
- Apache Airflow best practices: <https://airflow.apache.org/docs/apache-airflow/stable/best-practices.html>
- Apache Arrow columnar format: <https://arrow.apache.org/docs/format/Columnar.html>
- Apache Parquet: <https://parquet.apache.org/docs/>
- Apache Iceberg: <https://iceberg.apache.org/docs/latest/>
- OpenLineage object model: <https://openlineage.io/docs/spec/object-model/>
- dbt best practices: <https://docs.getdbt.com/best-practices>
- Kimball dimensional techniques: <https://www.kimballgroup.com/data-warehouse-business-intelligence-resources/kimball-techniques/dimensional-modeling-techniques/>

## 33. Cruces iniciales

| Manual | Autoridad de frontera |
|---|---|
| `ALGORITHMS_DATA_STRUCTURES_PROBLEM_SOLVING.md` | external/streaming algorithms, sketches y costo abstracto |
| `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md` | drivers, boundaries, ownership y quality scenarios |
| `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md` | CPU/memory/I/O y benchmark físico |
| `SOFTWARE_BACKEND_API_ENGINEERING.md` | source/serving APIs y contratos |
| `NETWORKING_DISTRIBUTED_STREAMING.md` | Kafka/Flink, delivery, state/checkpoints y backpressure |
| `DATABASE_STORAGE_INTERNALS.md` | indexes, execution, MVCC/WAL/recovery |
| `GPU_ACCELERATED_COMPUTING.md` | ETL/analytics acelerado, memoria y parity numérica |
| `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md` | IAM/privacy/cloud/SLO/incident |
| `MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md` | book/feed/order semantics |
| `AI_ENGINEERING_MASTER_MAP.md` | model data/evals/RAG/features |
| `FRONTEND_PRODUCT_ENGINEERING_UX.md` | dashboards, visualización y estado visible |

Cruces recíprocos cerrados el 2026-08-21: cada receptor incorpora sólo la obligación que cambia una decisión y conserva su autoridad. Este manual gobierna provenance, grain, tiempo, quality, lineage, publication y consumo; no reemplaza delivery, internals, costo físico, device, seguridad, semántica de dominio ni evaluación de IA.

Cruce con `TOOLCHAINS_BUILDS_PACKAGING_FFI.md` cerrado el 2026-08-21: DAG/job/SQL/codegen/connector y schema generator fijan code, toolchain, dependencies, config y target; backfill/replay registra artifact digest. Build reproducible no garantiza data reproducible si inputs/time cambian, y dataset version no demuestra que el binary/package pueda recrearse: conservar ambas chains y unirlas por run lineage.

Cruce con `NATIVE_MOBILE_DESKTOP_ENGINEERING.md` cerrado el 2026-08-21: telemetry móvil/escritorio fija schema/version, event time, app build, OS/device class, consent/purpose y offline upload/dedupe; no captura payload, token, texto, ubicación exacta o identifier persistente por comodidad. Datos gobierna ingestion/quality/lineage/metrics; nativo gobierna lifecycle, collection context y user controls. Crash-free, startup, jank y sync se segmentan sin convertir cardinalidad o fingerprints en tracking indebido.

## 34. Puerta de cierre

- [x] Contrato/ciclo del producto de datos.
- [x] Captura, raw, batch/stream y tiempos.
- [x] Transformación, idempotencia y schemas.
- [x] Parquet/Arrow/table formats/partitioning.
- [x] Distributed compute y orchestration.
- [x] Quality, lineage y modelado dimensional.
- [x] Analytics, gobierno, CI/CD, SRE y costo.
- [x] ML/order-book frontiers, Codex y suites.
- [x] Inventario y fuentes versionadas.
- [x] Auditoría Markdown/referencias: fences, encabezados, caracteres y referencias internas.
- [x] Cruces recíprocos y arbitraje.
- [x] Contrato ejecutable y capstones de planificación disponibles en `ENGINEERING_EXECUTION_PLAYBOOK.md`; outputs de datos reales siguen condicionados a source/workload/proyecto.

El conocimiento general v1 queda auditado. No declarar suites ejecutadas hasta que existan sus assets y resultados.
