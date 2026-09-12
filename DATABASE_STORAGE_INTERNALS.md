# Bases de datos y almacenamiento — manual operativo

> **Estado:** núcleo general v1 auditado dentro del alcance declarado el 2026-08-21. Verificados: CMU 15-445/645 Fall 2025, Berkeley CS186 Spring 2026, selección analítica CMU 15-721 Spring 2024, papers originales, PostgreSQL 18.6, RocksDB, SQLite, Faiss y pgvector. Auditoría transversal cerrada con sistemas/rendimiento, backend/API, redes/distribuidos, GPU, seguridad/SRE/cloud, microestructura y frontend.
> **Autoridad académica:** CMU Database Group, UC Berkeley y papers originales seleccionados.
> **Autoridad de implementación:** documentación y código oficial versionado de PostgreSQL, RocksDB y SQLite.
> **Propósito:** diseñar, implementar, elegir, medir y operar almacenamiento correcto y eficiente, desde páginas e índices hasta transacciones, recovery, optimización y distribución.

## 0. Cómo usar este manual

Este archivo sirve para:

- convertir invariantes de negocio en schema, constraints y transacciones;
- elegir estructura de almacenamiento e índice desde el workload, no desde una moda;
- razonar sobre planes, cardinalidad, memoria, I/O, concurrencia y recuperación;
- diseñar pruebas de anomalías, crash, restore, migración y degradación;
- dar a Codex un contrato de implementación y revisión profesional.

Etiquetas:

- `[SPEC]`: comportamiento documentado por estándar o implementación oficial;
- `[ACADEMIC]`: modelo, algoritmo o evidencia de curso/paper;
- `[CODE]`: detalle comprobable en código/documentación oficial;
- `[PROD]`: práctica operativa publicada y situada;
- `[MEASURED]`: resultado reproducible del sistema objetivo;
- `[IMPL]`: regla de diseño derivada para este manual;
- `[OPEN]`: pendiente o dependiente de versión/proyecto.

Reglas:

1. “ACID”, “SQL”, “NoSQL”, “MVCC” o “distributed” no especifican una garantía concreta.
2. La documentación de la versión desplegada manda sobre una simplificación académica.
3. Un índice mejora un patrón y cobra en writes, memoria, WAL, storage y mantenimiento.
4. Un plan estimado no es ejecución observada; un benchmark pequeño no representa producción.
5. Durabilidad se demuestra mediante crash y restore, no por recibir un `COMMIT` en un test feliz.

## 1. Contrato del sistema de datos

Antes de elegir producto, schema o índice, escribir:

```yaml
workload:
  operations: []
  read_write_ratio: ""
  keys_and_distributions: ""
  working_set: ""
  object_sizes: {p50: "", p99: "", max: ""}
  concurrency: ""
correctness:
  invariants: []
  isolation_required: ""
  ordering: ""
  idempotency: ""
durability:
  ack_means: ""
  rpo: ""
  rto: ""
  retention: ""
performance:
  latency_slo: {read: "", write: "", transaction: ""}
  throughput: ""
  burst: ""
  growth: ""
operations:
  backup_restore: ""
  migration: ""
  observability: ""
  failure_domains: []
```

Preguntas de diseño:

- ¿Cuál es la unidad atómica de cambio?
- ¿Qué estado constituye la fuente de verdad?
- ¿Qué lecturas pueden ser stale y cuánto?
- ¿Qué conflicto debe bloquear, abortar o reconciliar?
- ¿Qué significa confirmar: memoria, page cache, WAL durable, réplica o quorum?
- ¿Qué se reconstruye desde log y qué debe restaurarse desde backup?
- ¿Qué crecimiento domina: datos vivos, versiones muertas, índices, WAL o compaction?

## 2. Modelo relacional, schema e invariantes

### 2.1 Modelo antes que tablas accidentales

El modelo relacional representa relaciones y permite operadores cerrados sobre ellas. SQL añade bags, `NULL`, tipos y semánticas propias; no confundirlo con álgebra relacional pura.

Diseñar desde:

- entidades y valores, no objetos del framework;
- identidad estable y claves candidatas;
- dependencias funcionales e invariantes;
- cardinalidad y opcionalidad reales;
- historial, vigencia y ownership;
- consultas y mutaciones críticas.

### 2.2 Constraints son código de corrección

Preferir que la base rechace estados inválidos mediante:

- `NOT NULL` cuando ausencia no es válida;
- `CHECK` para predicados locales;
- `UNIQUE` para identidad lógica;
- `PRIMARY KEY` y `FOREIGN KEY` para referencias;
- exclusion/constraints especializadas cuando la implementación las ofrece;
- transacción o serialización explícita para invariantes entre varias filas.

**[IMPL]** Validar también en la API para devolver errores útiles, pero no depender de dos procesos que “miren antes de insertar”. Sólo la constraint/serialización situada en el punto de concurrencia cierra la carrera.

### 2.3 Normalización y desnormalización

Normalizar reduce redundancia y anomalías de actualización. Desnormalizar puede reducir joins o estabilizar snapshots, pero crea una obligación de sincronización.

Por cada duplicación declarar:

```text
source of truth → mecanismo de actualización → ventana de inconsistencia
→ detector de drift → reparación → rollback
```

No almacenar un valor derivado sin definir cuándo queda stale y cómo se recompone.

### 2.4 Tipos y tiempo

- dinero: precisión, escala, moneda y regla de redondeo;
- tiempo: instante, zona, calendario, resolución y clock source;
- identificador: identidad, ordenabilidad, exposición y distribución;
- decimal/float: error permitido y operaciones;
- JSON/blob: schema, límites, evolución e indexación.

`NULL` significa marcador desconocido/ausente según el modelo, no cadena vacía ni cero. La lógica ternaria afecta predicados, joins y constraints; probarla explícitamente.

## 3. Del dispositivo a páginas y registros

### 3.1 Jerarquía y unidad de I/O

```text
CPU cache → DRAM/buffer pool → OS page cache → block layer
→ device cache/controller → SSD/HDD → réplica/backup
```

Ninguna capa implica por sí sola persistencia. `write()` completado, page cache y energía estable tienen garantías distintas; documentar `fsync`/FUA/barriers y hardware reales.

### 3.2 Páginas

Un DBMS organiza datos e índices en páginas para amortizar I/O y gestionar espacio. Una página suele contener header, slots/offsets, tuples/records y espacio libre.

El esquema slotted-page permite mover records dentro de la página sin cambiar el identificador lógico que apunta al slot. Decisiones:

- tamaño de página versus random I/O, fan-out y amplificación;
- fixed versus variable-length;
- inline versus overflow/TOAST;
- checksums y detección de torn/corrupt pages;
- free-space y visibility metadata;
- compresión y encoding.

**[CODE]** PostgreSQL 18 documenta páginas de tabla e índice con header, item identifiers, free space, items y special space; su comportamiento concreto no debe extrapolarse a RocksDB o SQLite.

### 3.3 Row, column y formatos híbridos

| Layout | Ventaja típica | Costo típico |
|---|---|---|
| row/N-ary | acceso y update de fila completa | scans analíticos leen atributos innecesarios |
| column/decomposition | projection, compresión, SIMD | reconstrucción de filas y writes pequeños |
| PAX/híbrido | localidad por columna dentro de página | complejidad y beneficio dependiente de workload |

Compresión reduce bytes/I/O y puede aumentar throughput, pero consume CPU y cambia latencia aleatoria. Medir ratio, decode bandwidth, branch/cache behavior y costo de update.

### 3.4 Device-aware sin hardcodear folklore

SSD no vuelve gratis el random I/O; mantiene colas, erase/program, garbage collection, endurance y tail latency. HDD sigue importando en almacenamiento frío/secuencial. Cloud block/object storage añade red, throttling y semántica del proveedor.

## 4. Buffer pool y memoria

El buffer pool mapea page IDs a frames de memoria, coordina pin/unpin, dirty state, eviction y flush. Correctness requiere que una página usada no sea evicted y que WAL durable preceda al flush de datos correspondiente.

Políticas como LRU, CLOCK o variantes sólo son útiles respecto del patrón real. Scans grandes pueden expulsar el working set —*sequential flooding*—; separar o detectar accesos secuenciales puede evitar contaminación. **[ACADEMIC]** LRU-K estima reutilización con las últimas `K` referencias y ARC adapta espacio entre recencia y frecuencia; ninguno elimina la obligación de medir el workload. El proyecto de buffer pool de CMU usa ARC y hace explícitas cuatro piezas separables: page/frame table, replacer, disk scheduler y page guards.

Medir:

- hit/miss por relación/índice/operación;
- working set y reuse distance;
- dirty pages y flush rate;
- eviction útil versus churn;
- pin wait y lock contention;
- memoria fuera del cache principal;
- double caching con OS;
- NUMA locality.

**[IMPL]** Un “cache de 10 GB” no implica RSS de 10 GB. Índices, filters, memtables, pinned blocks, metadata, allocator fragmentation, conexiones y query memory cuentan por separado.

## 5. Índices y filtros

### 5.1 Contrato de un índice

Un índice es una copia organizada de keys y referencias/versiones. Elegirlo desde:

```text
predicados + orden + cardinalidad/selectividad + update rate
+ tamaño + memoria + concurrencia + maintenance + recovery
```

### 5.2 B+ tree

B+ tree mantiene keys ordenadas; nodos internos guían y hojas contienen entradas, normalmente enlazadas para range scans. Fan-out alto reduce altura. Insert/delete puede split/merge/redistribute; concurrencia exige latch protocol correcto.

**[ACADEMIC]** En *latch crabbing/coupling*, el recorrido adquiere el latch del hijo antes de soltar el ancestro y puede liberar ancestros cuando el hijo es seguro para la operación. La variante optimista baja con read latches y reinicia con protección más fuerte si encuentra una página insegura. El orden de adquisición debe ser global y coherente; scans entre hojas requieren una política que evite el ciclo entre sibling y recorrido descendente. Esto protege la estructura física, no sustituye transaction locks ni garantiza aislamiento.

Útil para:

- equality y range;
- prefix de índice compuesto compatible;
- `ORDER BY` y min/max bajo condiciones;
- uniqueness;
- index-only/covering cuando visibility y columnas lo permiten.

Costos:

- write amplification y WAL;
- page split/bloat;
- random access y cache misses;
- mantenimiento de cada índice secundario;
- versiones múltiples bajo MVCC.

El layout de una página de índice también es contrato: tipo, ocupación y capacidad en header; nodos internos con separadores y `m+1` hijos; hojas con key/value y enlace de rango. No introducir contenedores con construcción no trivial dentro de una página persistida ni depender del layout ABI accidental. Serialización, endianness, checksums y versionado deben ser explícitos en un motor real.

### 5.3 Hash

Hash ofrece igualdad esperada eficiente, pero no orden/range. Chaining, open addressing, extendible/linear hashing y resize tienen perfiles distintos. Medir load factor, probe length, skew y resize stalls.

### 5.4 Bloom y filtros probabilísticos

Bloom filter puede responder “definitivamente no” o “posiblemente sí”; no produce falsos negativos bajo uso correcto, sí falsos positivos. Bits/key y hashes intercambian memoria, CPU y false-positive rate.

No usar como autoridad de existencia. Es valioso para evitar I/O a SSTs/partitions que no pueden contener la key.

### 5.5 Índices compuestos, parciales y covering

- El orden de columnas responde a predicados y orden reales, no a una regla universal de “más selectiva primero”.
- Un partial index sólo sirve cuando el planner puede demostrar su predicate.
- Covering reduce heap access, pero aumenta índice y write cost.
- Expression index fija semántica/collation/function; su cambio requiere migración.
- Índice no usado todavía consume writes, WAL, vacuum y backup.

**[PROD]** PostgreSQL recomienda `ANALYZE` y datos realistas antes de juzgar index usage. Una tabla diminuta puede favorecer scan aunque el mismo porcentaje en producción no.

## 6. LSM trees y almacenamiento write-optimized

### 6.1 Write path

```text
write → WAL → mutable memtable → immutable memtable
→ flush a SST sorted → compaction entre runs/levels
```

Lectura puede consultar memtables y múltiples SSTs; indexes y filters reducen candidatos. Tombstones representan deletes hasta que compaction pueda descartarlos con seguridad.

### 6.2 Tres amplificaciones

- **write amplification:** bytes escritos físicamente / bytes lógicos;
- **read amplification:** runs/files/blocks consultados por lectura;
- **space amplification:** almacenamiento físico / datos vivos.

Leveled compaction suele limitar read/space amplification a cambio de más reescritura. Tiered/universal reduce parte del write cost a cambio de más runs y espacio. FIFO sólo es correcto para semántica tipo cache/retención compatible.

### 6.3 Compaction es trabajo de primer nivel

Si compaction no sigue write rate:

- crecen L0/runs y read amplification;
- sube espacio y puede dispararse write stall;
- compaction compite por CPU/I/O/cache;
- aumentan tails.

Medir backlog, pending bytes, stalls, compaction bandwidth/time, bytes read/written, L0 files, tombstones y latency por operación. Rate limiting puede suavizar interferencia, pero no crea capacidad.

## 7. Ejecución de consultas

### 7.1 Pipeline de consulta

```text
SQL → parse/bind → logical plan → rewrites
→ physical alternatives + cost → chosen plan
→ operators → storage/index → result
```

La semántica del resultado pertenece al contrato; el plan es una decisión reemplazable.

### 7.2 Modelos de ejecución

- iterator/Volcano: operador solicita próxima tuple; modular, con overhead por tuple;
- materialization: produce intermediarios; simple, usa memoria/I/O;
- vectorized: procesa batches y amortiza dispatch, favorece cache/SIMD;
- compiled: genera código/pipeline especializado; reduce interpretación, aumenta compile/complexity.

Pull/iterator pide datos desde el consumidor; push propaga desde el productor. Un operador que necesita ver toda o gran parte de la entrada —sort, build de hash join, aggregate global— rompe el pipeline. Esa frontera determina first-row latency, memoria, spill, cancellation y backpressure; no es sólo detalle interno.

Batch size demasiado pequeño desperdicia overhead; demasiado grande aumenta cache pressure y latency. **[MEASURED]** decidir por operator mix y hardware.

### 7.3 Operadores

| Operador | Variables dominantes |
|---|---|
| scan | pages/columns, predicates, zone maps, prefetch |
| sort | cardinalidad, key width, memory, runs/merge, spill |
| hash aggregate | groups, skew, table load, spill |
| nested-loop join | outer rows × inner access; excelente sólo con escala/acceso adecuados |
| hash join | build size, key width, skew, memory, partition/spill |
| sort-merge join | orden disponible, sort cost, ranges/duplicates |

Un spill no es una excepción rara: definir memory quota, temporary storage, cancellation, cleanup y observabilidad.

### 7.4 Paralelismo

Separar:

- intra-operator versus inter-operator;
- partitioning y exchange;
- local versus remote;
- work stealing/scheduling;
- skew y stragglers;
- NUMA y shared-state contention.

Más workers pueden empeorar por coordination, memory bandwidth, cache, I/O y reparto desigual.

Los *exchange operators* hacen explícito el movimiento: `gather`, `distribute/broadcast` y `repartition`. En ejecución distribuida o NUMA, el costo de serializar, copiar y mover una fila puede dominar el operador local. Paralelismo sólo escala hasta el recurso limitante; si I/O o memoria ya saturan, más workers agregan cola.

### 7.5 Ejecución analítica moderna

**[ACADEMIC]** El bloque seleccionado de CMU 15-721 muestra un diseño recurrente, no una receta única:

```text
columnar/open format + encoding/compression
→ pruning/pushdown + late materialization
→ vector batches/selection representation
→ cache/SIMD-friendly operators o compiled pipelines
→ morsel scheduling + NUMA-aware placement
→ exchanges hacia compute/storage desagregado
```

- Formato abierto no garantiza performance ni semántica transaccional: medir metadata pruning, small files, decode, schema evolution y object-store requests.
- Late materialization evita reconstruir filas antes de necesitar columnas, pero random gathers y columnas muy selectivas pueden cambiar el balance.
- SIMD pierde eficiencia con control flow divergente; selection vector/bitmap, batch size y encoding deben medirse juntos.
- Code generation reduce dispatch y materialización, pero compilation latency, code size/instruction cache y skew pueden dominar consultas breves.
- Morsels pequeños mejoran balance/stealing; demasiado pequeños elevan scheduling/atomic contention. Colocar datos y trabajo con NUMA-awareness evita remote memory traffic.
- Separar storage y compute mejora elasticidad/aislamiento, pero añade red, cache consistency, retry, metadata coordination y costo de egress/request. El SLO debe incluir cold start y cache-cold scans.

## 8. Optimización y cardinalidad

El optimizer explora equivalencias lógicas y alternativas físicas usando statistics y cost model. El mayor error suele propagarse desde cardinality estimation.

**[ACADEMIC]** No existe un plan físico universalmente mejor: scan versus index, hash versus sort-merge y orden de joins dependen de cardinalidad, orden disponible, memoria, I/O y propiedades requeridas. La búsqueda completa crece explosivamente. System R popularizó programación dinámica sobre alternativas y órdenes *left-deep*; optimizadores modernos amplían o podan ese espacio. Toda poda, timeout o regla fija es parte del contrato de calidad del plan.

Rewrites como predicate/projection pushdown, constant folding o convertir nested-loop equijoin en hash join sólo son válidos si preservan semántica de `NULL`, outer joins, duplicados, volatilidad, collation y errores observables. El orden de reglas importa: una transformación puede habilitar u ocultar la siguiente.

### 8.1 Statistics

Histogramas, distinct counts, most-common values y correlation aproximan distribuciones. Independencia entre columnas suele ser falsa; estadísticas multivariadas pueden reducir error cuando la implementación las soporta.

Registrar por nodo:

```text
estimated rows vs actual rows
estimated cost vs actual time
loops
buffers/I/O
memory/spill
workers and skew
```

El primer gran desvío de rows suele ser más causal que el último nodo lento.

### 8.2 Plan stability

Plan puede cambiar por data growth, stats, parámetros, versión, memory, parallelism o index availability. No congelar un plan por reflejo; corregir modelo/stats/schema cuando sea posible y mantener regression workload.

### 8.3 EXPLAIN con disciplina

- `EXPLAIN` muestra plan estimado.
- `EXPLAIN ANALYZE` ejecuta: no usar sin control sobre writes/carga.
- medir con datos, parameters y cache state representativos;
- distinguir startup, total, rows, loops y cumulative child work;
- incluir buffers/WAL/timing según implementación y costo aceptable.

## 9. Transacciones y anomalías

### 9.1 ACID con significado preciso

- **Atomicity:** commit completo o aborto lógico; recovery elimina efectos parciales según protocolo.
- **Consistency:** invariantes definidos por schema/aplicación, no una propiedad mágica del motor.
- **Isolation:** resultados permitidos entre transacciones concurrentes.
- **Durability:** efectos confirmados sobreviven los fallos incluidos en el modelo.

Una transacción no puede deshacer por sí sola un email, una llamada HTTP o una orden enviada fuera del DBMS. Esos efectos requieren idempotencia, outbox/inbox, compensación o un protocolo distribuido explícito.

### 9.2 Historias y serializabilidad

Conflict serializability razona con precedence graph; serializability exige equivalencia a algún orden serial. Linearizability añade orden de tiempo real para operaciones y no es sinónimo de serializable transactions.

Anomalías a probar:

- dirty/non-repeatable/phantom read;
- lost update;
- read/write skew;
- fractured read;
- serialization failure;
- stale replica read;
- efecto externo duplicado tras respuesta/commit ambiguo.

### 9.3 Niveles son implementación-específicos

No inferir comportamiento sólo del nombre SQL. **[SPEC]** PostgreSQL 18:

- trata `READ UNCOMMITTED` como `READ COMMITTED`;
- `READ COMMITTED` usa snapshot por statement y es default;
- `REPEATABLE READ` implementa snapshot isolation y aún permite serialization anomalies;
- `SERIALIZABLE` usa Serializable Snapshot Isolation y puede abortar con `40001`;
- la aplicación debe reintentar la transacción completa, no sólo el statement fallido.

Una secuencia/identity puede no revertirse con la transacción. No usar ausencia de gaps como invariante.

## 10. Control de concurrencia

### 10.1 Locks y 2PL

Two-phase locking separa fase de adquisición y liberación. 2PL produce historias conflict-serializable, pero la forma básica aún puede permitir cascading aborts y deadlocks. Strict 2PL conserva locks de write hasta commit/abort; strong strict 2PL conserva todos los locks, simplificando recoverability. Granularidad más fina aumenta concurrencia y metadata; escalamiento reduce metadata y aumenta conflicto. Intention locks permiten coordinar jerarquías sin inspeccionar cada lock descendiente.

Deadlock:

- prevenir con orden global/timeout;
- detectar mediante wait-for graph;
- abortar víctima y reintentar idempotentemente;
- observar lock wait, blockers, edad y objetos, no sólo count.

`wait-die` y `wound-wait` usan edad/timestamp para prevenir ciclos con políticas opuestas; elección de víctima y backoff deben evitar starvation. `SKIP LOCKED` cambia semántica —omite filas bloqueadas— y sirve para colas cooperativas, no como aceleración transparente de cualquier consulta.

Latches protegen estructuras internas cortas; locks protegen transacciones/datos. Confundirlos produce diseños incorrectos.

### 10.2 Timestamp ordering y OCC

Timestamp ordering valida contra orden lógico. Optimistic concurrency control trabaja sin locks largos y valida al commit; funciona bien con conflicto bajo, pero desperdicia trabajo y genera retries con contención/skew.

En OCC clásico, cada transacción atraviesa read, validation y write; su workspace privado no hace desaparecer phantoms de rangos. La validación debe cubrir el conjunto realmente leído, incluidas ausencias/rangos cuando afectan el resultado.

### 10.3 MVCC

MVCC mantiene versiones para que readers y writers se bloqueen menos. Requiere:

- regla de visibilidad/snapshot;
- transaction IDs/timestamps;
- almacenamiento de versiones y tombstones;
- garbage collection/vacuum;
- protección frente a transacciones/snapshots largos;
- índices coherentes con versiones.

MVCC no equivale a serializabilidad ni elimina locks. Constraints, DDL, unique checks y write conflicts aún coordinan.

Una implementación formativa útil mantiene por tuple una versión actual y una cadena de deltas/undo. El `read timestamp` selecciona la versión visible; el `commit timestamp` publica la nueva versión; un watermark igual al menor timestamp de lectores activos limita qué historia puede reclamarse. **[IMPL]** GC nunca puede inferirse sólo de la edad de pared: debe probar que ninguna transacción/snapshot/replica/backup puede alcanzar la versión.

### 10.4 Elegir desde conflicto

Medir matriz read/write por key/range, hot keys, duración, abort rate y wasted work. Sharding no elimina contención si todas las operaciones caen en una key lógica.

## 11. WAL, commit y crash recovery

### 11.1 Regla WAL

**[SPEC]** La regla central de PostgreSQL es: el log que describe un cambio debe quedar durable antes de escribir la página de datos modificada. Esto permite replay sin forzar todas las data pages en cada commit.

Separar:

```text
logical update
→ WAL record generated
→ WAL flushed to durability boundary
→ commit acknowledged
→ dirty page later flushed
```

`NO-FORCE` evita flush de todas las pages al commit; `STEAL` permite expulsar dirty pages de transacciones no confirmadas. Juntas exigen REDO y UNDO o un diseño equivalente.

### 11.2 ARIES como modelo

ARIES usa log sequence numbers, physiological logging, write-ahead rule, fuzzy checkpoints y recovery en fases:

1. **Analysis:** reconstruir transaction table y dirty-page state.
2. **Redo:** repetir historia desde el punto necesario de manera idempotente.
3. **Undo:** revertir losers y registrar compensation log records para que recovery mismo sea recuperable.

El modelo enlaza `LSN` del registro, `pageLSN` de la página y `flushedLSN` del log: una página dirty sólo puede persistirse cuando el WAL requerido ya es durable. Un fuzzy checkpoint registra, sin detener todo el sistema, la *active transaction table* y la *dirty page table*; `recLSN` ayuda a elegir desde dónde rehacer. Los compensation log records hacen que un crash durante undo pueda continuar sin volver a deshacer lo ya compensado.

No afirmar que un motor “usa ARIES” sin documentación/código; muchos adoptan partes o protocolos distintos.

### 11.3 Group commit y durability latency

Agrupar commits amortiza flush, pero agrega espera y acopla tails. Medir commit batch, WAL bytes, fsync latency, queue depth y storage behavior. `synchronous_commit` o equivalentes cambian qué pérdida acepta el ACK; documentarlo como semántica, no tuning inocuo.

### 11.4 Checkpoint

Checkpoint limita recovery work, pero puede crear burst de I/O. Debe coordinar dirty pages, WAL retention y replica/backup. Medir checkpoint duration, bytes, write rate, stalls y recovery time real.

## 12. Vacuum, garbage collection y mantenimiento

Versiones muertas no desaparecen por quedar invisibles. Vacuum/GC/compaction recupera espacio lógico/físico según motor.

En PostgreSQL, vacuum también es necesario para:

- reutilizar espacio de dead tuples;
- actualizar visibility map para index-only scans;
- sostener planner statistics mediante analyze;
- evitar transaction ID wraparound.

**[SPEC]** Son cuatro trabajos con frecuencias diferentes. `VACUUM` estándar reutiliza espacio dentro de la relación y puede coexistir con DML; `VACUUM FULL` reescribe, recupera más espacio al sistema operativo y toma `ACCESS EXCLUSIVE`, por lo que no es mantenimiento rutinario equivalente. La visibility map permite saltar heap fetches en index-only scans sólo cuando la página está marcada all-visible. Congelar XIDs/MXIDs antiguos es corrección, no optimización opcional.

Transacciones largas, replication slots o snapshots retenidos pueden impedir limpieza y acumular WAL/bloat. Alertar por edad y capacidad restante, no sólo tamaño actual.

Maintenance comparte CPU, I/O, memory y locks con tráfico. Presupuestarlo dentro de capacidad sostenida.

## 13. Replicación y distribución

Este capítulo aplica el modelo de `NETWORKING_DISTRIBUTED_STREAMING.md` al dato.

### 13.1 Replication contract

Declarar:

- líder, multi-leader o leaderless;
- synchronous/asynchronous y ACK set;
- consistency/read policy;
- failover y fencing;
- lag en bytes/tiempo/LSN;
- pérdida aceptable y rejoin;
- backup independiente.

Replica no es backup: replica corrupción, delete lógico y credencial comprometida.

### 13.2 Sharding

Partition key decide locality, parallelism y hot spots. Cross-shard query/transaction agrega red, coordinación y failure modes. Rebalancing mueve datos mientras sirve tráfico: medir amplification y disponibilidad.

### 13.3 Atomic commit y consensus

Two-phase commit coordina decisión commit/abort entre participantes pero puede quedar bloqueado si el coordinator falla sin protocolo adicional. Consensus acuerda estado/orden bajo otro modelo; no son intercambiables.

Usar saga sólo si las acciones tienen compensaciones de negocio válidas. Compensar no borra observaciones externas ya realizadas.

### 13.4 CDC y outbox

Dual-write DB+broker sin transacción común puede divergir. Transactional outbox guarda cambio y evento pendiente en la misma transacción local; un relay publica al menos una vez y consumidores deduplican. Definir orden, retention, replay y poison handling.

## 14. Schema evolution y migraciones

Una migración es un protocolo entre versiones de writers, readers, schema y datos.

Patrón expand/contract:

1. añadir forma compatible;
2. desplegar readers tolerantes;
3. dual-write/backfill con checkpoint y rate limit;
4. validar equivalencia;
5. cambiar source of truth;
6. detener forma antigua;
7. retirar después de rollback window.

Gate:

- lock level y duración;
- table rewrite/index build;
- WAL/replication lag;
- disk headroom;
- resumability/idempotency;
- old/new binary compatibility;
- rollback que preserve datos escritos por la versión nueva.

## 15. Conexiones, pools y prepared work

Una conexión consume memoria, process/thread/state y capacidad de locks. Pool grande no crea database throughput; puede transformar backpressure en colapso.

Dimensionar con:

- cores/IOPS y concurrency útil;
- transaction duration;
- queue wait y timeout;
- per-connection memory;
- serverless fan-out;
- failover/reconnect storm.

Prepared statements reducen parse/plan repetido, pero plan genérico versus específico puede cambiar rendimiento con parámetros skewed. Medir ambos y versionar invalidation behavior.

## 16. PostgreSQL 18: mapa operativo

> Versión actual observada el 2026-08-21: PostgreSQL 18.6; PostgreSQL 19 estaba beta. Fijar major/minor desplegado.

### 16.1 Storage y MVCC

- heap tuples contienen metadata de visibilidad;
- índices contienen entradas separadas y no “entienden” por sí solos la identidad lógica entre versiones;
- HOT puede evitar nuevas index entries bajo condiciones concretas;
- TOAST maneja valores grandes;
- free-space y visibility maps sostienen allocation y scans;
- autovacuum no es opcional en una operación seria.

**[SPEC]** HOT sólo es posible si el `UPDATE` no cambia columnas referenciadas por índices —excepto índices de resumen como BRIN— y la página original tiene espacio. Puede evitar nuevas entradas de índice y podar versiones intermedias; bajar `fillfactor` aumenta espacio para HOT a costa de densidad. Medir `n_tup_hot_upd`, non-HOT updates, dead tuples, edad XID/MXID y páginas all-visible.

### 16.2 Planner

Usar `ANALYZE`, estadísticas extendidas cuando hay correlación, `EXPLAIN (ANALYZE, BUFFERS, WAL)` en ambiente controlado y vistas de estadísticas. No empezar cambiando cost constants globales para ocultar estimaciones defectuosas.

`ANALYZE` implica ejecución; en PostgreSQL 18 activa `BUFFERS` salvo que se deshabilite. Los buffers informados por nodo incluyen trabajo de hijos y no son bloques distintos. Leer `actual rows × loops`, primer desvío de estimación, rows removed/recheck, index searches, heap fetches, sort method/disk, hash batches/memory y planning versus execution. Forzar un scan sirve como experimento comparativo, no demuestra que ese plan deba fijarse.

### 16.3 Backup/PITR

PITR combina base backup consistente con archivo continuo de WAL. Definir:

- retención y prueba de cadena;
- credenciales/cifrado/inmutabilidad;
- timeline y target de recovery;
- RPO según archivo real;
- restore drill y tiempo;
- verificación lógica posterior.

Un backup que nunca se restauró es evidencia incompleta.

**[SPEC]** PITR físico necesita una base backup y una secuencia continua de WAL desde el comienzo de esa copia hasta el target. `pg_dump`/`pg_dumpall` son copias lógicas y no contienen lo necesario para replay físico. El intervalo entre base backups intercambia almacenamiento de WAL por tiempo de restore; PostgreSQL 18 también soporta base backups incrementales vinculados por manifests. Probar pérdida de un segmento, timeline derivada, target anterior/posterior y restauración en infraestructura aislada.

### 16.4 Locks, checkpoints y ACK durable

- PostgreSQL detecta deadlocks y aborta una víctima no predecible; adquirir múltiples objetos en orden consistente reduce ciclos, pero el código debe reintentar la transacción completa.
- Sin deadlock ni timeout configurado, un lock incompatible puede esperar indefinidamente. Nunca sostener una transacción mientras se espera input humano o una RPC no acotada.
- Con `full_page_writes` activo, la primera modificación de una página tras checkpoint registra una imagen completa para proteger contra torn pages. Checkpoints más frecuentes acortan redo pero elevan flush y WAL; medir no sólo su duración sino el pico posterior.
- `synchronous_commit=off` puede reconocer antes del flush local: un crash puede perder commits recientemente reconocidos sin dejar inconsistencia interna. No equivale a `fsync=off`, que cambia el riesgo de corrupción.
- Con standby síncrono, `local`, `remote_write`, `on` y `remote_apply` esperan fronteras distintas. El SLO debe nombrar exactamente qué host escribió, hizo durable y/o aplicó antes del ACK.

## 17. RocksDB: mapa operativo

RocksDB es un storage engine embebido basado en LSM, no un servidor SQL completo. La aplicación conserva responsabilidades de schema, coordinación, distribución y operación.

### 17.1 Paths

Write va a memtable y WAL; flush produce SST. Read consulta memtables/cache/SSTs guiado por indexes/filters. Compaction reorganiza y descarta versiones/tombstones cuando es seguro.

Cada write recibe sequence number; snapshots, iterators y transacciones usan esa historia para visibilidad y GC. El MANIFEST registra cambios del conjunto de archivos/versions para reconstruir un estado coherente: WAL de datos, MANIFEST de metadata y SST inmutables cumplen papeles distintos.

### 17.2 Memory budget real

Contabilizar:

- block cache;
- indexes y filters dentro/fuera del cache;
- mutable/immutable memtables;
- blocks pinned por iterators;
- compaction buffers y allocator;
- OS page cache si no se usa direct I/O.

### 17.3 Correctness gate

- `WriteOptions.sync`, WAL enabled/disabled y filesystem definen durability: `sync=false` —default— puede sobrevivir process crash por WAL en page cache, pero no promete machine-crash durability; `sync=true` sincroniza WAL antes de retornar; `disableWAL=true` admite perder memtable no flusheada;
- snapshot fija sequence visibility, no congela recursos gratis;
- snapshot no persiste a través de reopen; iterator implícitamente fija una vista al crearse y su lifetime puede retener memtables, blocks y SSTs ya compactadas;
- prefix extractor/comparator forman parte del formato lógico;
- column families comparten WAL pero tienen configs/resources propios;
- `WriteBatch` da atomicidad de múltiples keys, pero el conflicto transaccional requiere `TransactionDB`/`OptimisticTransactionDB`; reads ordinarios no crean automáticamente read-write conflicts —usar la API y el isolation contract correctos—;
- checkpoint crea una copia consistente; puede hard-linkear SSTs en el mismo filesystem y copiar MANIFEST/CURRENT/WAL necesario. Backup/checkpoint debe incluir y validar wrapper/application metadata y restaurarse como DB independiente.

### 17.4 Compaction, stalls y admission

Si flush/compaction no sigue al ingreso, RocksDB ralentiza o detiene writers por demasiadas immutable memtables, demasiados L0 files o excessive pending compaction bytes. Los umbrales pueden ser por column family mientras el stall afecta al DB completo. Subir límites sólo aplaza deuda y puede cambiar un stall visible por agotamiento de disco y read amplification.

Medir por CF y agregado:

- incoming/flush/compaction bytes y write amplification;
- L0 files, pending compaction bytes y compaction score;
- stall count/time/cause y `Status::Incomplete` con `no_slowdown`;
- read amplification, Bloom usefulness, block-cache hit y iterator-pinned bytes/files;
- space amplification y headroom durante compaction.

Leveled, universal/tiered y FIFO intercambian read, write y space amplification; no son niveles de “calidad”. Cambiar compaction style, comparator, prefix extractor, merge operator o timestamp semantics exige revisar compatibilidad del formato y procedimiento de migración.

Nunca copiar opciones de tuning sin misma versión, storage, key/value distribution y SLO.

## 18. SQLite: mapa operativo

SQLite es una biblioteca transaccional embebida con un archivo de base; su ausencia de servidor cambia failure/locking/deployment, no elimina ingeniería.

### 18.1 Rollback journal versus WAL mode

Rollback journal preserva contenido previo para deshacer. WAL appendea cambios y readers pueden seguir viendo snapshot mientras un writer agrega; checkpoint mueve WAL hacia la base. WAL mejora ciertos patrones de concurrencia, pero sólo existe un writer a la vez y el checkpoint puede afectar tails.

El wal-index usa shared memory, por lo que todos los participantes deben estar en la misma máquina y un network filesystem no es un soporte transparente para WAL. Cada reader conserva su end mark; el checkpointer no puede sobrepasar la versión de un reader activo. El auto-checkpoint default ocurre alrededor de 1000 páginas y puede cargar latencia sobre el `COMMIT` que cruza el umbral. Separarlo a otro thread cambia la distribución de latencia, no elimina I/O ni los bloqueos de un checkpoint agresivo.

`PASSIVE` avanza sin interferir y puede quedar incompleto; `FULL`, `RESTART` y `TRUNCATE` tienen requisitos/bloqueos más fuertes. Medir WAL pages/bytes, checkpointed frames, busy result, reader age, commit p99.9 y tiempo de reset. Un WAL grande degrada espacio y puede aumentar trabajo de lectura.

### 18.2 Atomicidad depende del entorno

La explicación oficial modela locks, journal, flush y propiedades del filesystem/device. Copiar el archivo mientras hay actividad o ignorar archivos auxiliares puede producir backup inválido. Usar backup API/procedimiento documentado.

En WAL, `synchronous=FULL` sincroniza en cada commit; `NORMAL` puede perder transacciones tras power loss/hard reset aunque conserve consistencia lógica bajo el modelo documentado. `OFF` permite reordenamientos que pueden corromper tras fallo de energía. SQLite sólo puede confiar en lo que VFS, filesystem y dispositivo reportan; un medio que miente sobre flush invalida la garantía.

**[SPEC][PROD] Gate de versión 2026:** el *WAL-reset bug* puede corromper raramente con dos o más conexiones y write/checkpoint concurrentes. Afecta probablemente 3.7.0–3.51.2; está corregido desde 3.51.3 y en backports 3.44.6/3.50.7. No desplegar WAL multi-connection sobre una versión vulnerable.

### 18.3 Gate embebido

- threads/processes y locking mode;
- busy timeout versus backpressure real;
- transaction scope corto;
- WAL growth/checkpoint;
- filesystem local soportado;
- corruption checks y backup restore;
- migrations compatibles con binaries desplegados.

SQLite serializa writes y ofrece aislamiento serializable salvo el caso explícito de shared cache + `read_uncommitted`. `BEGIN DEFERRED` puede empezar leyendo y fallar con `SQLITE_BUSY` al promover a writer; `BEGIN IMMEDIATE` intenta reservar write desde el inicio. `busy_timeout` no reemplaza admission/backpressure y todo `BUSY` debe tener política acotada, observable e idempotente.

En WAL, `-wal` forma parte del estado persistente: separarlo del archivo principal puede perder commits o corromper. Para una copia viva usar Online Backup API, `VACUUM INTO` o herramienta oficial; para una copia offline demostrar que no hay transacción y preservar cualquier hot journal/WAL necesario. Verificar el resultado con apertura, `integrity_check`, schema/version y muestreo lógico.

## 19. Vector y similarity indexes

Un vector index no reemplaza metadata relacional, autorización ni filtros exactos. Separar:

- métrica: cosine, dot product, L2 u otra;
- exact search versus approximate nearest neighbor;
- recall@k/precision, latency y throughput;
- build/update/delete cost;
- memory y storage;
- filtering antes/durante/después;
- freshness y compaction/rebuild.

HNSW, IVF y product quantization tienen trade-offs diferentes. La auditoría siguiente fija principios comunes con papers originales, Faiss y pgvector; un motor administrado aún requiere su propia documentación/version exacta antes de dar parámetros.

### 19.1 Contrato matemático y de datos

Congelar juntos: embedding model + revisión, dimension, dtype/quantization, preprocessing, métrica y política de normalización. Cosine, inner product y L2 no son intercambiables sin condiciones; por ejemplo, con vectores unit-normalized, maximizar inner product induce el mismo ranking que cosine, pero los scores y tolerancias siguen necesitando contrato.

El embedding es dato derivado. Guardar `source_id`, model/version, created time y estado de backfill; cambiar el modelo crea otro espacio y normalmente exige dual-read/evaluación/reindex, no mezclar vectores silenciosamente. `NULL`, zero vector, NaN/Inf y dimensión errónea deben tener política explícita.

### 19.2 Baseline exacto antes de ANN

**[CODE]** Faiss `Flat`/`IndexFlatL2` o `IndexFlatIP` hace búsqueda exhaustiva y es la referencia exacta. Conservar un corpus/query set estratificado y ground truth exacto para medir:

```text
recall@k por segmento/filtro/tenant
+ latency p50/p95/p99 y QPS bajo concurrencia
+ build/add/delete/update time
+ RAM/VRAM/storage y bytes/vector
+ freshness, failure/reload y costo total
```

Medir queries individuales y batches; cold/warm cache; misma distribución y también drift/out-of-distribution. Recall global puede ocultar cero recall en un tenant o filtro raro. La métrica final del producto —calidad RAG, clasificación o retrieval— acompaña a recall ANN, no es sustituida por ella.

### 19.3 HNSW

**[ACADEMIC]** HNSW construye capas de proximity graphs: pocos elementos alcanzan capas altas y la búsqueda desciende hasta la capa base. `M`/grado cambia memoria, build y conectividad; `efConstruction` cambia calidad/costo de construcción; `efSearch` cambia candidatos, recall y latency. Los nombres/límites exactos son de cada implementación.

Ventaja típica: high recall con búsqueda rápida sin entrenamiento de centroides. Costos: grafo en memoria, build/update y comportamiento de deletes/compaction/replication dependiente del motor. Nunca prometer complejidad o recall del paper sobre la distribución propia sin medición.

### 19.4 IVF y product quantization

IVF entrena centroides/listas, asigna vectores y busca sólo `nprobe` listas cercanas. Pocos probes pierden vecinos; aproximarse a todas las listas elimina el beneficio y con IVFFlat se acerca al exhaustivo. Training set debe representar la distribución; crear IVFFlat con pocos datos o después de fuerte drift degrada asignación y puede requerir retraining/rebuild.

PQ divide el vector en subespacios y codifica cada parte con un codebook. Reduce bytes y distance bandwidth a cambio de quantization error. OPQ/pre-transform puede facilitar compresión. Recuperar más candidatos y re-rankear con vectores originales o mayor precisión permite intercambiar candidate cost por recall final.

### 19.5 Filtering, multitenancy y pgvector

**[CODE]** En pgvector 0.8.x, filtros ordinarios pueden aplicarse después del scan ANN y devolver menos de `k`. Iterative scans siguen explorando hasta suficientes resultados o hasta límites de tuples/probes/memoria; `strict_order` preserva orden exacto de distancia dentro de lo explorado y `relaxed_order` intercambia orden por recall/performance. Esto no convierte ANN en exacto.

Elegir según selectividad y cardinalidad:

- B-tree/partial index y exact scan para subconjunto pequeño;
- partitioning/índice por tenant cuando aislamiento y recall lo exigen;
- ANN compartido + iterative scan cuando el coste y mezcla se validaron;
- candidate retrieval + exact re-ranking para calidad final.

Autorización debe restringir candidatos antes de exponer IDs/payloads; filtrar sólo al final puede filtrar resultados visibles, pero también desperdicia presupuesto y crear canales laterales de timing/count. Probar tenants skewed y ausencia de permiso.

### 19.6 Gate de producción

1. Registrar versión exacta de pgvector/FAISS/motor y formato serializado.
2. Construir exact baseline y thresholds por segmento, no sólo promedio.
3. Barrer parámetros bajo presupuesto de RAM/VRAM, build y p99; no copiar defaults.
4. Probar add/update/delete, tombstones, vacuum/reindex, restart y replica/backup.
5. Validar filtros, `k` insuficiente, duplicates, ties y orden determinista cuando importa.
6. Probar model migration con índice antiguo/nuevo coexistiendo y rollback.
7. Observar drift de embeddings/queries, recall canary, index age y rebuild progress.

Fuentes: HNSW original <https://arxiv.org/abs/1603.09320>; PQ original <https://doi.org/10.1109/TPAMI.2010.57>; Faiss oficial <https://github.com/facebookresearch/faiss/wiki/Guidelines-to-choose-an-index>; pgvector oficial <https://github.com/pgvector/pgvector>.

## 20. Observabilidad

### 20.1 Señales

| Capa | Señales mínimas |
|---|---|
| workload | QPS/TPS, mix, key/cardinality/skew, object size |
| latency | queue + execute + commit por operación y percentiles |
| planner | estimate/actual, plan changes, stats age |
| memory | buffer/cache hit, working set, eviction, query memory, RSS |
| I/O | read/write bytes, IOPS, fsync, queue, device tails |
| concurrency | active/idle transactions, locks, waits, abort/retry |
| WAL/recovery | generation, flush, archive, replay, checkpoint, restore |
| maintenance | vacuum/compaction backlog, bloat, tombstones, stalls |
| replication | lag, quorum/leader, slot retention, conflict, failover |
| capacity | data/index/WAL/temp growth y headroom |

### 20.2 Correlación

Correlacionar request/trace con query fingerprint, transaction, shard y plan sin registrar secrets ni parámetros sensibles. No usar SQL crudo con valores como label de métrica.

### 20.3 Diagnóstico

```text
síntoma → ventana/carga → wait class/recurso
→ query/transaction/maintenance causante
→ hipótesis falsable → experimento controlado
→ corrección → comparación y rollback
```

## 21. Testing y failure injection

### 21.1 Capas

- model/property tests para invariantes;
- schema/constraint tests;
- isolation tests con schedules coordinados;
- differential tests contra modelo simple;
- query plan/performance regression con distribución realista;
- crash/restart en cada frontera de WAL/flush/checkpoint;
- disk full, I/O error, corruption y slow storage;
- replica lag, partition, failover y split-brain defenses;
- backup/restore/PITR drills;
- migration mixed-version y rollback;
- load/soak con maintenance concurrente.

### 21.2 Prueba de transacción

Para cada invariante concurrente:

```yaml
initial_state: ""
transactions:
  A: []
  B: []
controlled_interleavings: []
allowed_outcomes: []
forbidden_outcomes: []
isolation_and_retry: ""
evidence: ""
```

Un test que ejecuta dos threads sin controlar el schedule puede pasar siempre sin probar la anomalía.

### 21.3 Crash consistency

Confirmar ACK, matar proceso/VM o cortar dispositivo en puntos seleccionados, reiniciar, verificar invariantes y repetir. Fault injection debe respetar seguridad del entorno; nunca experimentar destructivamente sobre la única copia.

## 22. Benchmarking

### 22.1 Manifiesto

Registrar:

- motor/build/config y schema/indexes;
- hardware, NUMA, memory, filesystem/device/cloud class;
- dataset y distribución/skew/cardinality;
- cache warm/cold y preload;
- clientes, conexiones y offered load;
- transaction/query mix y payload;
- durability/isolation/replication;
- maintenance y background work;
- duración, warmup y repetitions.

### 22.2 Resultados

Reportar:

- throughput y goodput;
- p50/p95/p99/p99.9/max y timeout/error/abort;
- queue wait versus service;
- CPU, memory, I/O, WAL y network;
- amplification y storage growth;
- recovery/restore time;
- plan y cardinality errors.

No comparar motores con durability, isolation o replication diferentes sin destacarlo como resultado semánticamente distinto.

## 23. Almacenamiento para order books y baja latencia

El estado vivo de un book secuenciado suele pertenecer a una state machine en memoria con ownership claro; una transacción SQL por delta rara vez es el hot path correcto. Persistencia puede dividirse:

```text
raw feed / normalized append log
→ deterministic in-memory state
→ snapshots/checkpoints
→ analytical/history store
```

Obligaciones:

- sequence/epoch del venue y timestamps por dominio;
- append ordering y checksum;
- backpressure sin ocultar gaps;
- snapshot consistente con log offset;
- replay determinista;
- retención/compaction sin destruir auditabilidad requerida;
- storage writer aislado para no bloquear receive/state update;
- recuperación medida hasta `LIVE` correcto.

No usar una base para “predecir” lo que el protocolo no dice. Primero reconstruir estado correcto; después consultar o analizar.

Cruce recíproco con `MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md` auditado el 2026-08-21:

- raw bytes/eventos conservan provenance, schema/version, channel, epoch y sequence antes de normalizar;
- un checkpoint referencia exactamente el último evento aplicado y se publica atómicamente;
- trade/order corrections, busts y late events son eventos, no updates destructivos sin historia;
- executions/positions se reconcilian contra fuente privada autorizada; intención de orden no crea posición;
- time-series partitions y retention nunca borran la ventana necesaria para replay/audit declarados;
- analytics separa intervalos `LIVE` de `GAPPED/STALE` y no rellena state correctness con interpolación.

El storage puede demostrar qué se recibió y reprodujo; no puede completar hidden liquidity, paquetes perdidos no recuperables o reglas no documentadas.

## 24. Gate profesional

Antes de producción:

1. invariantes en schema/transacciones y tests concurrentes;
2. nivel de aislamiento y retry contract explícitos;
3. índices justificados por workload y write cost;
4. plans medidos con cardinalidad/data realistas;
5. memory/connection/query budgets limitados;
6. WAL/durability/ACK definidos;
7. backup restaurado y PITR ensayado;
8. vacuum/compaction/checkpoint dentro de capacidad sostenida;
9. replication/failover/fencing probados;
10. migrations resumibles, observables y reversibles;
11. disk full/corruption/slow I/O con respuesta segura;
12. SLO y alertas ligados a síntomas e invariantes.

## 25. Contrato para Codex

Al pedir implementación o revisión, proporcionar:

```yaml
engine_and_version: ""
schema_and_scale: ""
operations_and_queries: []
invariants: []
isolation_durability_replication: ""
latency_throughput_slo: ""
growth_retention: ""
hardware_storage: ""
failure_model: []
```

Codex debe devolver:

1. supuestos y cuestiones `[OPEN]`;
2. invariantes y failure model;
3. schema/index/transaction design;
4. costo por read/write/memory/WAL/storage;
5. código/migración mínima;
6. tests de correctness, concurrency y crash;
7. observabilidad y benchmark;
8. rollout, backfill, rollback y restore;
9. divergencias específicas de versión citadas.

### 25.1 Auditoría transversal cerrada

Cruce completado el 2026-08-21 con `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md`, `SOFTWARE_BACKEND_API_ENGINEERING.md`, `NETWORKING_DISTRIBUTED_STREAMING.md`, `GPU_ACCELERATED_COMPUTING.md` y `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md`:

| Frontera | Invariante compartido | Prueba/medición obligatoria |
|---|---|---|
| ACK ↔ storage | `write()`/page cache no equivale a durable; nombrar WAL/data/replica/quorum alcanzado | process crash, OS/VM crash y power-loss model; pérdida permitida por ACK |
| buffer/cache ↔ memoria | buffer pool, OS cache, query memory, memtables, pinned blocks y allocator compiten | RSS/cgroup/OOM, hit/eviction, NUMA, reclaim y tails con maintenance |
| query ↔ API deadline | pool wait + lock wait + execute + commit + response caben en deadline | timeout por fase, cancellation real y ausencia de trabajo huérfano |
| retry ↔ transaction | retry sólo con operación/transacción completa idempotente o deduplicada; `40001`, deadlock y commit ambiguo son distintos | history test + duplicate request + respuesta perdida antes/después de commit |
| schema ↔ rollout | expand/contract soporta binaries coexistentes; rollback no destruye writes nuevos | mixed-version, backfill resumible, replica lag/WAL y rollback drill |
| DB ↔ evento | outbox hace atómico dato+intención local; relay/consumer siguen al menos una vez | crash en insert/commit/publish/ack, replay y dedupe retention |
| replication ↔ consenso | replication, 2PC, consensus y backup resuelven problemas distintos | partition, leader fencing, quorum loss, stale read, rejoin y restore independiente |
| CDC/stream ↔ MVCC/WAL | cursor/LSN/offset, snapshot y retention deben formar un corte recuperable | bootstrap snapshot + concurrent changes + disconnect + truncation/replay |
| low latency ↔ maintenance | p99.9 incluye fsync, checkpoint, vacuum/compaction, page fault y queueing | carga sostenida con background work; no sólo benchmark warm y corto |
| observabilidad ↔ seguridad | correlacionar request/query/txn/LSN sin SQL/PII/secrets de alta cardinalidad | redaction tests, bounded labels y trazabilidad de incidente |
| GPU ↔ persistencia | event/kernel/collective completion no equivale a WAL, fsync, commit ni réplica durable; storage y feeder pueden dominar | crash entre compute/write/WAL/sync/ACK, checksum/result parity, presión conjunta de pinned/NUMA/I/O y benchmark end-to-end |
| IAM/authz ↔ datos | identidad de workload corta y authz por operación/tenant/objeto; roles de admin/backup/replication no se heredan por comodidad | cross-tenant/role escalation, stale/revoked credential, breakglass y audit efectivo |
| cifrado/keys ↔ recovery | TLS/at-rest/backup cifran bajo lifecycle y trust roots recuperables; key loss o revocation forma parte del disaster model | rotate/rekey, cert expiry, KMS/control-plane loss y restore aislado con keys correctas |
| deletion ↔ integridad | soft-delete/backup/validator protegen contra corrupción sin retener indefinidamente ni prometer purge instantáneo | delete/export/retention, corrupt backup, ransomware, clean point y functional restore |
| cloud ↔ failure domain | servicio administrado no elimina quotas, region, control plane, engine semantics, shared fate ni exit | AZ/region/quota/provider outage, capacity for failover y portable recovery evidence |
| client cache ↔ autoridad | URL/store/service worker no sustituyen MVCC, constraints ni source of truth; freshness y tenant scope explícitos | reload/back-forward/offline, invalidation, permission change, stale replica y cross-tenant cache-key tests |
| optimistic UI ↔ commit | estado tentativo porta operation ID y se reconcilia con resultado durable; cancelación browser no revierte commit | respuesta perdida antes/después de commit, duplicate submit, conflict y recovery tras reload |
| pagination/search ↔ query | cursor/order/filter/version estables; límites y cardinalidad no se ocultan en infinite scroll | concurrent inserts/deletes, cursor expiry, deep navigation, plan/cost y accesibilidad de resultados |
| algoritmo ↔ storage | hash/tree/sort/join se eligen por operaciones y prueba; I/O, pages, cardinalidad, skew y maintenance determinan costo real | oráculo pequeño, EXPLAIN/trace, datasets adversariales y crossover con working set/compaction |

Resultado: una optimización local no está aceptada si cambia aislamiento/durabilidad, desplaza cola a otra capa, rompe recovery o mejora promedio degradando goodput/p99.9.

Cruce con `ALGORITHMS_DATA_STRUCTURES_PROBLEM_SOLVING.md` cerrado el 2026-08-21: el manual algorítmico aporta invariantes, B-trees/hash, sorting, grafos, DP y análisis external-memory; este manual decide páginas, buffer pool, WAL/MVCC, recovery y efectos de concurrencia. Big-O RAM no predice I/O ni write amplification. Todo índice o plan se valida con cardinalidad/skew reales, maintenance sostenido y la misma isolation/durability.

Cruce con `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md` cerrado el 2026-08-21: arquitectura declara source of truth, ownership, consistency, RPO/RTO y límites transaccionales; este manual valida si engine, schema, index, WAL/MVCC, replication y recovery cumplen. Shared DB y database-per-service son opciones con trade-offs, no mandatos. Cache/view/event no adquiere autoridad sin freshness, invalidation y reconciliation explícitos.

Cruce con `DATA_ENGINEERING_ANALYTICS.md` cerrado el 2026-08-21: DB internals decide snapshot/WAL/CDC cut, execution, index, transaction y recovery; datos decide raw/core/mart, grain, quality, lineage y consumer semantics. OLTP replica no es warehouse ni backup. Extract/CDC respeta MVCC y retention; analytical layout se valida con query shape, cardinality, pruning, maintenance y point-in-time correctness.

Cruce con `NATIVE_MOBILE_DESKTOP_ENGINEERING.md` cerrado el 2026-08-21: la base local sostiene la source-of-truth offline mediante transacciones, migrations, tombstones, outbox y sync metadata; cache borrable y datos del usuario quedan separados. Storage gobierna durability/recovery/schema; nativo gobierna lifecycle, UX pending/confirmed y conflict policy. Probar kill/disk-full/corruption durante commit, migration y reconcile sobre el package real.

Prohibido:

- recomendar índice sin query shape y distribución;
- afirmar exactly-once por una transacción local cuando existe side effect externo;
- envolver en retry una transacción no idempotente sin revisar el callback completo;
- usar `EXPLAIN ANALYZE` destructivo en producción sin control;
- desactivar fsync/WAL/constraints para “ganar benchmark” sin declarar pérdida semántica;
- tratar replica como backup;
- ocultar bloat/compaction/WAL bajo tamaño de datos lógico.

## 25.1 Práctica pública Apple/FoundationDB — simulación determinista de fallos

**[CODE] Evidencia fijada:** `apple/foundationdb@cd9bf922c51160e1759bbd3f7ab5eaafd9db8fca` más documentación oficial. FoundationDB ejecuta un cluster simulado dentro de un proceso, usa RNG/tiempo deterministas e inyecta fallos de red, proceso, disco y datacenter; una seed permite reproducir la historia si el workload no introduce nondeterminismo externo. Complementa simulación con performance tests y fallos sobre hardware real.

La práctica generalizable es diseñar el sistema para que tiempo, aleatoriedad, scheduling, red, storage y procesos puedan sustituirse por un ambiente controlado:

```text
same code paths + virtual time + seeded choices
+ fault schedule + invariant workload
→ reproducible failing history + trace
```

**[IMPL] Gate para un storage/distributed component:**

1. oracle/invariantes verificables durante y después del workload, no sólo “no crash”;
2. seed, config, binary digest y fault schedule persistidos;
3. faults en commit/WAL/flush/checkpoint/election/reconfiguration/backup y recovery;
4. partición, delay, duplicate/reorder permitido por el modelo, disk full/slow, reboot y proceso resurrectado;
5. reducir/minimizar la historia fallida sin perder reproducción;
6. code probes o coverage semántica para demostrar que se alcanzaron caminos raros;
7. ejecutar también real-cluster, multithread, filesystem/device y performance tests: el simulador sólo prueba su modelo.

**Buggify para clientes:** inyectar errores retriables/permanentes en puntos instrumentados permite comprobar manejo de transacciones y retry. Debe estar deshabilitado en producción y no sustituye faults de servidor/red reales.

La determinización tiene costo arquitectónico: threads, reloj real, syscalls no virtualizadas y fuentes externas rompen replay. Declarar qué queda fuera en vez de atribuir cobertura total. Esta técnica no vuelve correcta a FoundationDB por autoridad; ofrece un método reusable para encontrar historias que tests convencionales rara vez alcanzan.

Fuentes: [Simulation and Testing](https://apple.github.io/foundationdb/testing.html), [Client Testing/Buggify](https://apple.github.io/foundationdb/client-testing.html), [paper de FoundationDB](https://www.foundationdb.org/files/fdb-paper.pdf) y [código fijado](https://github.com/apple/foundationdb/tree/cd9bf922c51160e1759bbd3f7ab5eaafd9db8fca).

## 25.2 Código PostgreSQL/Meta — probar historias, no endpoints aislados

**[CODE] Evidencia fijada:** `postgres/postgres@7788fcf39bf98c3403f052f3edd73d86fe1501a7` (PostgreSQL License) y `facebook/rocksdb@ad81b654a4dec07f7067f1ec3749301e403ee857` (Apache-2.0 más componentes con avisos propios).

PostgreSQL incluye un `isolationtester` que abre sesiones concurrentes y ejecuta permutaciones explícitas o todas las intercalaciones compatibles con el orden de cada sesión. Sus specs convierten anomalías como write skew, deadlock y conflictos de locks en historias reproducibles. SSI observa estructuras peligrosas de dependencias y puede abortar con `serialization_failure`, incluso con falsos positivos conservadores: la aplicación debe reintentar la **transacción completa**, no continuar desde la última sentencia.

RocksDB aporta `db_stress`, expected-state tracking, snapshots compartidos, checksums, reopen, fault-injection filesystem, cambios de configuración y variaciones de compaction/recovery. Su valor transferible es combinar operaciones aleatorias con un oracle independiente y conservar seed/config; “no hubo crash” no demuestra que valores, visibilidad entre column families, MANIFEST/WAL o snapshots sean coherentes.

**Puerta de pruebas de un storage component:**

1. modelar cada actor como sesión y enumerar/intercalar steps alrededor de read, lock, write, commit, abort, timeout y retry;
2. definir outcomes permitidos por isolation level; comparar estado y errores, no timing accidental;
3. mantener oracle/expected state independiente de la ruta bajo prueba, junto con seed y trace de operaciones;
4. inyectar I/O error, short write, sync failure, disk full, corruption, process kill y restart en WAL/MANIFEST/SST/checkpoint/archive;
5. ejecutar reopen/recovery y verificar checksums, invariantes lógicas, visibility por snapshot y ausencia de valores fantasma;
6. estresar compaction/flush mientras hay reads, snapshots y writes; observar stalls, pending bytes, amplification y espacio;
7. repetir con versiones mixed durante upgrade/rollback y con backup/replica restaurada como instancia independiente.

No copiar el harness interno de un motor como oracle universal: adaptar las historias al contrato de la aplicación y ejecutar además tests físicos sobre filesystem/device reales.

Fuentes fijadas: [PostgreSQL source](https://github.com/postgres/postgres/tree/7788fcf39bf98c3403f052f3edd73d86fe1501a7), [isolation test harness](https://github.com/postgres/postgres/blob/7788fcf39bf98c3403f052f3edd73d86fe1501a7/src/test/isolation/README), [SSI implementation notes](https://github.com/postgres/postgres/blob/7788fcf39bf98c3403f052f3edd73d86fe1501a7/src/backend/storage/lmgr/README-SSI) y [RocksDB `db_stress`](https://github.com/facebook/rocksdb/tree/ad81b654a4dec07f7067f1ec3749301e403ee857/db_stress_tool).

## 26. Fuentes primarias inventariadas

### 26.1 Cursos

- CMU 15-445/645 Fall 2025: <https://15445.courses.cs.cmu.edu/fall2025/>
- Schedule CMU 15-445/645, 25 clases: <https://15445.courses.cs.cmu.edu/fall2025/schedule.html>
- Berkeley CS186 Spring 2026: <https://cs186berkeley.net/>
- CMU 15-721 Spring 2024, analítica moderna: <https://15721.courses.cs.cmu.edu/spring2024/schedule.html>

#### Cobertura auditada de CMU 15-445/645 Fall 2025

**[ACADEMIC]** Se verificaron el schedule completo, las notas 01–24 y las especificaciones oficiales P0–P4. La secuencia compacta, sin repetir cada clase, es:

| Bloque | Material | Obligación transferible |
|---|---|---|
| modelo y SQL | 01–02 | separar semántica lógica de layout y plan físico; SQL opera con bags además de `NULL` y orden sólo cuando se pide |
| storage | 03–06 | páginas/slotted records, buffer pool, row/column/PAX y compresión según access path |
| acceso | 07–10 | hash, B+ tree, filters y concurrencia estructural; coste de writes y recuperación incluido |
| operadores | 11–14 | sort/aggregate/join con spill; iterator, materialization, vectorization, pipelines y paralelismo |
| optimización | 15–16 | equivalencias, estadísticas, cardinalidad, cost model y búsqueda acotada |
| concurrencia | 17–20 | serializabilidad, 2PL/deadlocks, timestamp/OCC y MVCC como alternativas con fallos distintos |
| durabilidad | 21–22 | STEAL/NO-FORCE, WAL, checkpoints, ARIES y crash durante recovery |
| distribución | 23–24 | partición, replicación, atomic commit y consenso como problemas diferentes |
| integración | 25 | razonar extremo a extremo; ninguna capa compensa silenciosamente una anterior incorrecta |

#### Ruta práctica equivalente — no copiar entregas

La utilidad profesional está en reproducir los invariantes en una implementación propia o laboratorio autorizado, no en entregar soluciones del curso:

1. **P0 — estructura probabilística concurrente:** Count-Min Sketch en C++17, thread safety, memoria fija, tests, sanitizers y medición; evitar un único latch global si se exige escalabilidad.
2. **P1 — buffer pool:** page/frame table, ARC, disk scheduler y RAII page guards; probar pin/eviction, dirty flush, fallos y carreras.
3. **P2 — B+ tree:** layout de páginas internas/hoja, insert/delete/search, iterator y optimistic latch crabbing; probar splits/merges/root changes y lecturas/escrituras concurrentes.
4. **P3 — ejecución:** scans y writes que mantienen índices, aggregate/joins, hash join + rules, external merge sort, limit y window functions; comparar `EXPLAIN` y resultado antes/después de cada rewrite.
5. **P4 — MVOCC:** timestamps, watermark, version chains con undo, visibility, executors, índices, abort y GC; probar write conflicts, self-modification, readers antiguos y recuperación de versiones.

Dependencia real:

```text
correctitud de páginas/buffer
  → correctitud y liveness del índice
    → resultados correctos de ejecutores y rewrites
      → visibilidad/abort/GC correctos bajo concurrencia
```

Un test feliz al final no localiza el error. Cada flecha exige unit tests, randomized/differential tests, sanitizers, invariants internos, concurrency stress y benchmark con resultado validado.

#### Contraste independiente — Berkeley CS186 Spring 2026

Se auditó el schedule oficial completo como prueba de cobertura. Confirma SQL; disks/buffers/files; cost models, B+ tree y spatial indexes; sorting/hashing/joins/iterators; plan space y cost search; locking; recovery; parallel query; 2PC/Paxos; NoSQL/MapReduce/Spark. Sus proyectos B+ tree, joins+query optimization, locking, recovery y NoSQL cubren la misma cadena con otra implementación. **Resultado:** no apareció un vacío de núcleo; añadió spatial indexing y exige mantener separado 2PC de consensus.

#### Selección auditada — CMU 15-721 Spring 2024

No se traslada el curso entero. Se incorporaron sólo bloques que cambian construcción de sistemas analíticos: lakehouse/composable systems; formatos columnar, encoding y compression; hyper-pipelining/Velox; vectorized y compiled execution; morsel-driven/NUMA scheduling; hash y multi-way joins; optimizer/Cascades, unnesting y adaptive cost; y análisis de Dremel/BigQuery, Photon/Delta, Snowflake, DuckDB, Yellowbrick y Redshift. La obligación común quedó en §7.5: medir fronteras de materialización, datos/movimiento, memoria, skew y desacople storage/compute.

### 26.2 PostgreSQL 18

- Manual actual e internals: <https://www.postgresql.org/docs/current/>, <https://www.postgresql.org/docs/current/internals.html>
- Indexes y planner: <https://www.postgresql.org/docs/current/indexes.html>, <https://www.postgresql.org/docs/current/using-explain.html>, <https://www.postgresql.org/docs/current/planner-stats.html>
- Concurrency/isolation: <https://www.postgresql.org/docs/current/mvcc.html>, <https://www.postgresql.org/docs/current/transaction-iso.html>
- WAL y PITR: <https://www.postgresql.org/docs/current/wal-intro.html>, <https://www.postgresql.org/docs/current/continuous-archiving.html>
- Vacuum y storage: <https://www.postgresql.org/docs/current/routine-vacuuming.html>, <https://www.postgresql.org/docs/current/storage-page-layout.html>
- HOT y locks: <https://www.postgresql.org/docs/current/storage-hot.html>, <https://www.postgresql.org/docs/current/explicit-locking.html>
- WAL/checkpoints/configuración: <https://www.postgresql.org/docs/current/wal-configuration.html>, <https://www.postgresql.org/docs/current/runtime-config-wal.html>
- Monitoring: <https://www.postgresql.org/docs/current/monitoring-stats.html>

### 26.3 RocksDB

- Repositorio: <https://github.com/facebook/rocksdb>
- Overview: <https://github.com/facebook/rocksdb/wiki/RocksDB-Overview>
- WAL: <https://github.com/facebook/rocksdb/wiki/Write-Ahead-Log-(WAL)>
- Compaction: <https://github.com/facebook/rocksdb/wiki/Compaction>
- Write stalls: <https://github.com/facebook/rocksdb/wiki/Write-Stalls>
- Memory y block cache: <https://github.com/facebook/rocksdb/wiki/Memory-usage-in-RocksDB>, <https://github.com/facebook/rocksdb/wiki/Block-Cache>
- Tuning guide: <https://github.com/facebook/rocksdb/wiki/RocksDB-Tuning-Guide>
- Transactions, snapshots e iterators: <https://github.com/facebook/rocksdb/wiki/Transactions>, <https://github.com/facebook/rocksdb/wiki/Snapshot>, <https://github.com/facebook/rocksdb/wiki/Iterator>
- Checkpoints y backup: <https://github.com/facebook/rocksdb/wiki/Checkpoints>, <https://github.com/facebook/rocksdb/wiki/How-to-backup-RocksDB>

### 26.4 SQLite

- Atomic commit: <https://www.sqlite.org/atomiccommit.html>
- WAL: <https://www.sqlite.org/wal.html>
- Isolation: <https://www.sqlite.org/isolation.html>
- Query planner: <https://www.sqlite.org/queryplanner.html>
- File format: <https://www.sqlite.org/fileformat.html>
- Transactions y Backup API: <https://www.sqlite.org/lang_transaction.html>, <https://www.sqlite.org/backup.html>
- Modos documentados de corrupción y gate de versión WAL: <https://www.sqlite.org/howtocorrupt.html>, <https://www.sqlite.org/wal.html#the_wal_reset_bug>

### 26.5 Papers fundacionales auditados

| Trabajo original | Afirmación central incorporada | Obligación de ingeniería |
|---|---|---|
| Selinger et al., *Access Path Selection in a Relational DBMS* (System R, 1979) | SQL no prescribe access path ni join order; el optimizer compara alternativas mediante costo y propiedades | estimar no basta: contrastar cardinalidad/plan/costo con ejecución y controlar el espacio/tiempo de búsqueda |
| Mohan et al., *ARIES* (1992) | WAL, repeating history, `pageLSN`, CLRs y recovery analysis/redo/undo soportan STEAL/NO-FORCE y rollback parcial | probar crash en cada frontera, incluido otro crash durante recovery; no atribuir ARIES completo a motores distintos |
| O’Neil et al., *The LSM-Tree* (1996) | merges amortizan mantenimiento de índices bajo alta inserción trasladando trabajo | presupuestar read/write/space amplification, compaction y stalls; “write optimized” no significa trabajo eliminado |
| Hellerstein, Stonebraker y Hamilton, *Architecture of a Database System* (2007) | process model, admission, query processor, storage, transactions y utilities forman un sistema acoplado | medir colas y recursos de extremo a extremo; no optimizar un operador ignorando buffer, WAL, maintenance y OS |
| Berenson et al., *A Critique of ANSI SQL Isolation Levels* (1995) | tres fenómenos SQL no caracterizan todas las implementaciones; define Snapshot Isolation y muestra que puede no ser serializable | especificar historias/anomalías reales y probarlas; nunca inferir garantías por el nombre del nivel |
| Ports y Grittner, *Serializable Snapshot Isolation in PostgreSQL* (2012) | SSI rastrea dependencias read-write y aborta estructuras peligrosas preservando gran parte de la concurrencia MVCC | toda transacción serializable debe tolerar retry completo; predicate/SIREAD locks son detección, no bloqueo ordinario |

Fuentes primarias:

- System R: <https://doi.org/10.1145/582095.582099>
- ARIES, IBM Research: <https://research.ibm.com/publications/aries-a-transaction-recovery-method-supporting-fine-granularity-locking-and-partial-rollbacks-using-write-ahead-logging>
- LSM tree: <https://dsf.berkeley.edu/cs286/papers/lsm-acta1996.pdf>
- Arquitectura DBMS: <https://cs.uwaterloo.ca/~david/cs848/background/Architecture%20of%20a%20Database%20System-Hellerstein%20Stonebraker.pdf>
- Crítica de isolation SQL, Microsoft Research: <https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/tr-95-51.pdf>
- SSI en PostgreSQL: <https://arxiv.org/abs/1208.4179> y código/README vigente <https://github.com/postgres/postgres/blob/master/src/backend/storage/lmgr/README-SSI>

### 26.6 Extensiones avanzadas

- Bw-tree, Masstree y modern index concurrency: incorporar sólo si el proyecto requiere construir/reemplazar un índice concurrente;
- HNSW, IVF/PQ, Faiss y pgvector quedaron auditados en §19. Antes de parametrizar un proyecto administrado, añadir su documentación/version exacta.

## 27. Cierre de auditoría v1

### 27.1 Matriz de trazabilidad

| Capacidad | Autoridad académica | Implementación/operación | Salida verificable |
|---|---|---|---|
| modelo/schema | CMU 01–02; CS186 SQL | PostgreSQL/SQLite | constraints + migrations + invariant tests |
| páginas/buffer | CMU 03–06, P1; CS186 storage | PostgreSQL layout; RocksDB memory | page/frame invariants + eviction/flush tests |
| hash/B+ tree | CMU 07–10, P2; CS186 B+ tree | PostgreSQL indexes | split/merge/range/concurrency + plan evidence |
| LSM | O’Neil et al. | RocksDB | amplification/stall/capacity benchmark |
| ejecución | CMU 11–14, P3; CMU 15-721 | PostgreSQL `EXPLAIN`; Faiss exact path | result equivalence + spill/memory/NUMA profile |
| optimizer | CMU 15–16; System R; CS186 | PostgreSQL planner/stats | estimate/actual + plan regression corpus |
| aislamiento | CMU 17–20, P4; Berenson/SSI | PostgreSQL, RocksDB, SQLite | controlled histories + retry/abort contract |
| WAL/recovery | CMU 21–22; ARIES; CS186 | tres motores auditados | crash matrix + restart invariants |
| distribución | CMU 23–24; CS186 2PC/Paxos | PG replication/CDC + manual de red | partition/fencing/replay/restore histories |
| vector ANN | HNSW y PQ originales | Faiss + pgvector | exact ground truth + recall/latency/memory gates |
| operación | arquitectura DBMS | vacuum/compaction/checkpoint/backup docs | sustained load + maintenance + restore drill |
| extremo a extremo | cuatro cruces transversales | backend, OS, red, GPU y storage | goodput/p99.9 sin degradar semántica |

### 27.2 Suites de aceptación ejecutables

Cada proyecto deriva casos concretos de estas suites; un ítem sin evidencia mantiene el proyecto —no este núcleo general— como `[OPEN]`:

1. **Schema:** generar casos válidos/inválidos y probar constraints, cascades, `NULL`, collations y límites.
2. **Query:** comparar resultado contra modelo simple; capturar plan, estimate/actual, buffers, spills y parámetros.
3. **Isolation:** orquestar interleavings permitidos/prohibidos; clasificar deadlock, serialization y commit ambiguo.
4. **Crash:** fallar antes/después de WAL, sync, commit, flush y checkpoint; repetir crash durante recovery.
5. **Restore:** restaurar backup/PITR/checkpoint aislado, verificar chain/timeline/schema/checksum e invariantes lógicas.
6. **Migration:** coexistencia old/new, backfill checkpointable, cancel/restart, rollback y writes durante transición.
7. **Capacity:** carga sostenida con vacuum/compaction/checkpoint/backup; disk-full/slow I/O y memory pressure.
8. **Distribution:** partition/failover/fencing, lag, stale reads, CDC snapshot+cursor, duplicate/reorder y rejoin.
9. **Vector:** exact ground truth, recall por filtro/tenant, drift, insufficient `k`, rebuild, delete y model migration.
10. **Order book:** gap/snapshot/delta/replay determinista, writer lento/caído y tiempo hasta estado `LIVE` correcto.

### 27.3 Extensiones condicionadas

- Motor/proveedor no cubierto: añadir documentación y versión exacta antes de fijar garantías o parámetros.
- Bw-tree, Masstree u otro índice concurrente: auditar paper+código sólo si se va a construir o adoptar.
- Los cruces con sistemas, backend, red, GPU y seguridad/SRE/cloud están cerrados; mercados se cerrará recíprocamente cuando exista ese manual, sin invalidar el alcance interno de bases de datos v1.
- Revisar videos completos sólo si aportan contenido que no está en notas, papers o especificaciones, o si se necesita una atribución literal puntual.
