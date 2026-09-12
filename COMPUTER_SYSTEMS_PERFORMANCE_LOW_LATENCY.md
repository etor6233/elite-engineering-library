# Computer Systems, Performance & Low Latency — manual operativo

> **Estado:** núcleo general v1 auditado el 2026-08-21; los cruces con backend/API, redes/distribuidos, datos/almacenamiento, GPU, seguridad/SRE/cloud, microestructura y frontend están cerrados; las extensiones dependientes de hardware/proyecto o manuales futuros se declaran al final.
> **Autoridad académica:** MIT 6.172; MIT 6.S081/6.828; CMU 15-213/CS:APP; OSTEP.
> **Autoridad normativa:** Intel SDM/Optimization Manual; C/C++ memory model; LLVM atomics; Rust Project; Linux documentation.
> **Autoridad operativa:** metodología de Brendan Gregg y evidencia reproducible del sistema objetivo.
> **Propósito:** razonar desde código hasta hardware, construir software eficiente y diagnosticar latencia, throughput, memoria, concurrencia y variabilidad sin adivinar.

## Cómo usar este manual

1. Definir primero corrección, carga, SLO y entorno.
2. Medir el sistema completo antes de aislar componentes.
3. Localizar el recurso o dependencia limitante.
4. Formular una hipótesis falsable.
5. Cambiar una variable material por vez.
6. Verificar corrección, rendimiento y regresiones en el mismo artefacto desplegable.
7. Conservar datos crudos, configuración, binario y procedimiento de rollback.

No usar este archivo como una lista de trucos. Una optimización sin perfil, contrato de corrección y benchmark reproducible es solamente una apuesta.

## Procedencia

- `[SPEC]`: comportamiento normativo de estándar, ISA o documentación oficial versionada.
- `[ACADEMIC]`: síntesis fiel de curso, texto o paper identificado; no es transcripción literal.
- `[PROD]`: práctica publicada por un operador o experto de producción identificado.
- `[MEASURED]`: resultado reproducible dentro del entorno declarado.
- `[IMPL]`: traducción propia a diseño, código, pruebas o runbook.
- `[OPEN]`: tema todavía no auditado completamente contra la fuente.

## 1. Contrato de rendimiento

### 1.1 No existe “rápido” sin una métrica

**[ACADEMIC]** MIT 6.172 enseña rendimiento como ingeniería de sistemas escalables mediante análisis, medición, optimización algorítmica, arquitectura, caché y paralelismo. Brendan Gregg sitúa la disciplina entre aplicación, sistema operativo, kernel y hardware, con dos metas generales: mejorar precio/rendimiento y reducir outliers de latencia.

**[IMPL]** Antes de tocar código, escribir:

```text
unidad de trabajo:
distribución de entrada:
concurrencia ofrecida:
latencia objetivo p50/p95/p99/p99.9/máxima:
throughput o goodput objetivo:
presupuesto CPU/memoria/red/storage/energía/costo:
restricciones de corrección:
entorno y versión:
```

Métricas mínimas:

| Dimensión | Medida |
|---|---|
| latencia | distribución completa; percentiles y máximo observado |
| capacidad | requests, mensajes, bytes, eventos u operaciones por segundo |
| goodput | trabajo válido que cumple SLO, no sólo intentos procesados |
| utilización | CPU por core, memoria, bandwidth, disco, red, acelerador |
| saturación | run queue, queue depth, backlog, outstanding I/O, throttling |
| errores | timeout, drop, retry, rechazo, corrupción o resultado inválido |
| eficiencia | trabajo útil por core, byte, joule o unidad monetaria |

### 1.2 Ecuaciones de control

```text
speedup = T_baseline / T_candidate
throughput = operaciones_completadas / tiempo
goodput = operaciones_correctas_dentro_del_SLO / tiempo
Little: L = λW
```

Para una fracción serial `s` y `P` procesadores, el límite ideal de Amdahl es:

```text
S(P) <= 1 / (s + (1-s)/P)
```

**[IMPL]** No usar Amdahl como predicción exacta: colas, coordinación, coherencia, NUMA, ancho de banda y desbalance añaden costos que la fórmula básica no modela.

## 2. Método de investigación

### 2.1 Secuencia obligatoria

```text
correctness baseline
→ workload representativo
→ medición end-to-end
→ perfil por recurso y stack
→ hipótesis causal
→ microbenchmark si aísla la hipótesis
→ cambio mínimo
→ verificación diferencial
→ carga sostenida y tails
→ decisión y registro
```

**[PROD]** La metodología USE de Brendan Gregg examina cada recurso mediante utilización, saturación y errores. Es una estrategia de exploración, no una garantía de causa: debe conectarse con stack traces, eventos, colas y experimentos.

### 2.2 Preguntas que evitan optimizar el lugar equivocado

- ¿El tiempo es CPU ejecutando, espera de scheduler, page fault, bloqueo, I/O o red?
- ¿La carga es compute-bound, bandwidth-bound, latency-bound o contention-bound?
- ¿El cuello está en trabajo útil, serialización, copias, allocation, syscall o coordinación?
- ¿Aumentar throughput empeora edad de cola o percentiles?
- ¿El benchmark conserva la misma corrección y calidad?
- ¿La muestra incluye warm-up, steady state, throttling y degradación prolongada?

## 3. Protocolo de benchmark reproducible

### 3.1 Manifiesto del entorno

Registrar al menos:

- CPU exacta, stepping/microarquitectura, cores/threads y topología NUMA;
- RAM, canales, frecuencia efectiva y huge-page policy;
- SO, kernel, microcode, firmware y mitigaciones relevantes;
- governor, turbo/boost, C-states, límites de potencia y temperatura;
- compilador, versión, target, flags, LTO/PGO y bibliotecas;
- afinidad de CPU/IRQ, aislamiento, container/cgroup y vecinos;
- tamaño/distribución de datos, seed, concurrencia y duración;
- commit, configuración, binario y hash del artefacto.

### 3.2 Medición correcta

**[ACADEMIC]** MIT 6.172 dedica una clase a medición y timing, además de tareas de profiling, vectorización, cachegrind, análisis dinámico y sincronización.

**[IMPL]** Exigir:

1. reloj monotónico apropiado;
2. calentamiento separado de la muestra;
3. suficientes operaciones por muestra para superar resolución y overhead;
4. repetición intercalada de baseline/candidate cuando el entorno deriva;
5. resultados crudos, no sólo el promedio;
6. percentiles y dispersión; bootstrap o intervalos cuando corresponda;
7. detección de outliers explicada, no borrado automático;
8. control de dead-code elimination y constant folding;
9. verificación de output para impedir que el compilador elimine trabajo;
10. carga sostenida para observar thermal throttling, GC, reclaim y colas.

### 3.3 Microbenchmark versus sistema

| Microbenchmark | End-to-end |
|---|---|
| aísla una operación | conserva interacción real |
| facilita perfiles estables | revela colas, copias y coordinación |
| puede caber artificialmente en caché | incluye working set real |
| tiende a omitir cold paths | incluye startup, fallos y tail |

Un microbenchmark acepta o rechaza una hipótesis local. El sistema end-to-end decide si el producto mejoró.

## 4. Del código a la máquina

### 4.1 Cadena de transformación

**[ACADEMIC]** MIT 6.172 recorre source → compilador → machine code → hardware → ejecución; CS:APP estudia representación, assembly, arquitectura, optimización, linking, control excepcional, memoria virtual, I/O, red y concurrencia desde la perspectiva del programador.

```text
source
→ frontend/IR
→ optimizaciones
→ instrucciones y layout
→ loader/linker/runtime
→ microarquitectura
→ memoria y dispositivos
→ kernel/scheduler
→ resultado observable
```

**[IMPL]** Una revisión de rendimiento debe inspeccionar, según el caso:

- complejidad y cantidad de trabajo;
- IR/vectorization report;
- assembly generado;
- instrucciones, branches y dependencias;
- cache/TLB misses y bandwidth;
- syscalls, faults, context switches y locks;
- colas e I/O externos.

### 4.2 Lo que el compilador puede y no puede demostrar

- El compilador optimiza bajo la semántica del lenguaje, no bajo la intención del programador.
- Undefined behavior permite transformaciones que vuelven inválido el razonamiento intuitivo.
- Alias potencial, alineación desconocida o dependencias aparentes pueden impedir vectorización.
- `volatile` no es sincronización entre threads y no sustituye atomics.
- Una abstracción puede compilar sin costo, pero sólo IR/assembly y medición lo comprueban.
- Flags agresivos pueden cambiar semántica numérica o portabilidad.

**[IMPL] Gate:** no aceptar “el compilador lo hará” ni “el compilador es malo”; producir el diagnóstico, el assembly y el impacto medido.

## 5. Trabajo antes que trucos

**[ACADEMIC]** Las reglas de Bentley presentadas en MIT 6.172 incluyen packing/encoding, inicialización en compile time, loop unrolling, short-circuiting, fast paths y combinación de tests. El principio dominante es reducir o reorganizar trabajo antes de microoptimizar instrucciones.

Orden preferido:

1. eliminar trabajo innecesario;
2. elegir algoritmo/estructura apropiados;
3. evitar conversiones, copias y allocations;
4. mejorar locality y representación;
5. permitir optimización del compilador;
6. paralelizar trabajo suficiente;
7. usar SIMD/intrinsics cuando el perfil lo justifique;
8. escribir assembly sólo con una razón extraordinaria y pruebas por target.

Correcciones:

- loop unrolling puede aumentar code size y presión de instruction cache;
- packing reduce bytes pero añade encode/decode y puede empeorar acceso;
- fast path introduce ramas y una segunda implementación que debe permanecer equivalente;
- bit hacks pueden ser menos claros y peores que una instrucción reconocida por el compilador;
- una tabla precomputada intercambia cómputo por memoria y caché.

## 6. Jerarquía de memoria y localidad

### 6.1 Modelo operativo

```text
registers → L1 → L2 → LLC → DRAM local → DRAM remota → storage/red
```

Cada nivel cambia por procesador; el modelo ordena el razonamiento, no asigna latencias universales.

Conceptos obligatorios:

- localidad temporal y espacial;
- cache line, asociatividad, sets y evictions;
- working set e instruction footprint;
- TLB, page walks y page size;
- hardware/software prefetch;
- write allocation y coherencia;
- bandwidth, latency y memory-level parallelism;
- NUMA y first-touch;
- false sharing.

### 6.2 Estructuras cache-friendly

| Decisión | Beneficio posible | Riesgo |
|---|---|---|
| contiguous storage | prefetch y menos indirections | movimientos/reallocation |
| SoA en vez de AoS | cargar sólo campos usados, vectorizar | acceso por objeto más complejo |
| tiling/blocking | reutilizar datos en caché | tuning dependiente del target |
| compact IDs | menos bytes y mejor locality | conversión y límites de rango |
| arena/pool | allocation predecible | lifetime grueso, memoria retenida |
| padding/alignment | evitar false sharing o permitir SIMD | mayor footprint |

**[IMPL]** El layout forma parte de la API del hot path. Medir bytes tocados por operación, working set, misses y bandwidth; no inferir locality sólo por la sintaxis.

## 7. Procesos, memoria virtual y kernel

**[ACADEMIC]** MIT 6.S081 usa xv6/RISC-V para estudiar procesos, system calls, traps, page tables, filesystem, locking e interacción hardware/software. OSTEP organiza el sistema operativo alrededor de virtualización, concurrencia y persistencia.

### 7.1 Costos que una API de alto nivel puede ocultar

- transición user/kernel y validación;
- scheduling, preemption y migration;
- page fault minor/major y copy-on-write;
- allocation, reclaim, compaction y swap;
- filesystem cache, writeback y fsync;
- interrupciones, softirqs y driver queues;
- copias y cambios de representación;
- locks internos y límites de descriptor.

### 7.2 Reglas

- Page cache no equivale a persistencia; distinguir write completion, flush y durable commit.
- `mmap` no elimina I/O ni faults; cambia la interfaz y el patrón de acceso.
- Huge pages reducen presión TLB pero pueden aumentar fragmentación, latencia de allocation y desperdicio.
- Afinidad puede reducir migración y mejorar locality, pero puede crear hotspots o impedir balanceo.
- Busy polling reduce wake-up latency en ciertos entornos a cambio de cores, energía y fairness.
- Real-time scheduling mal configurado puede impedir progreso de tareas críticas del sistema.

## 8. Concurrencia y paralelismo

### 8.1 Distinguir problemas

| Concepto | Pregunta |
|---|---|
| concurrencia | ¿qué tareas pueden estar activas simultáneamente? |
| paralelismo | ¿qué trabajo puede ejecutarse al mismo tiempo? |
| asincronía | ¿cómo progresa una tarea mientras espera? |
| distribución | ¿cómo coordinar fallos y estado entre máquinas? |

### 8.2 Correctitud primero

Una optimización concurrente debe declarar:

- estado compartido y ownership;
- invariantes;
- punto de linearización si corresponde;
- orden permitido de eventos;
- política de cancelación/shutdown;
- backpressure;
- progreso esperado: blocking, lock-free u otro;
- reclamación segura de memoria;
- comportamiento ante overload y fallo parcial.

### 8.3 Memory model y atomics

**[SPEC]** LLVM modela atomics para C/C++ y otros lenguajes mediante loads/stores atómicos, `cmpxchg`, `atomicrmw` y fences. Los orderings ofrecen garantías distintas; acquire obtiene sentido de sincronización al emparejarse con release sobre una relación válida. `seq_cst` añade un orden total entre operaciones secuencialmente consistentes.

```text
relaxed: atomicidad del objeto, sin publicar otros datos
release: operaciones anteriores deben publicarse antes
acquire: operaciones posteriores observan una publicación correspondiente
acq_rel: ambas direcciones en read-modify-write/fence
seq_cst: acquire/release + orden total SC
```

**[IMPL] Correcciones obligatorias:**

- atomicidad no implica ordering suficiente;
- x86 fuerte no convierte código C++ con data race en válido;
- `volatile` no reemplaza `atomic`;
- lock-free no significa wait-free ni más rápido;
- CAS loop puede sufrir contention, starvation y ABA;
- reclamación requiere una estrategia: epoch, hazard pointers, ownership u otra demostrable;
- validar también en ARM/RISC-V si el producto pretende portabilidad.

### 8.4 Locks versus lock-free

Preferir la implementación más simple que cumpla el SLO. Un mutex puede superar a un algoritmo lock-free con baja contención, producir código más mantenible y permitir mejor espera. Lock-free se justifica cuando la medición y requisitos de progreso muestran que blocking/convoy/preemption es el limitante.

## 9. Allocation y lifetime

**[ACADEMIC]** MIT 6.172 incluye storage allocation serial y paralelo; OSTEP cubre free-space management y estructuras sincronizadas.

Investigar:

- tasa y tamaño de allocation;
- lifetime distribution;
- fragmentation interna/externa;
- contention por arena;
- remote free y NUMA;
- zeroing y page faults;
- pico/steady-state y memoria retenida.

Alternativas condicionadas:

- stack/value semantics;
- reservar/reutilizar buffers;
- slab/pool por clase de tamaño;
- arena por request/fase;
- object cache por thread;
- allocator general alternativo.

No sustituir un allocator por fama. Medir workload real, RSS, faults, fragmentation, tail latency y shutdown/leaks.

## 10. Lenguajes de sistemas

### 10.1 C++

**[SPEC/IMPL]** Usar RAII y tipos para expresar ownership; evitar raw owning pointers; declarar lifetime y thread-safety; preferir vistas (`span` o equivalente) para buffers no propietarios; comprobar overflow, bounds y alineación; limitar `unsafe` implícito generado por casts, aliasing y lifetime incorrecto.

Para performance:

- medir allocations y copies, no asumir “zero-cost”;
- controlar layout sólo donde sea parte del contrato;
- revisar exception/RTTI policy por artefacto, no dogma;
- inspeccionar vectorization y assembly;
- usar atomics sólo con prueba de memory order;
- ejecutar ASan, UBSan, TSan y análisis estático en configuraciones apropiadas.

### 10.2 Rust

**[SPEC]** Ownership y el sistema de tipos convierten muchas violaciones de memoria y concurrencia en errores de compilación. `Send` y `Sync` expresan capacidades entre threads. Threads y async tienen costos y casos de uso diferentes; async no vuelve compute-bound work automáticamente paralelo.

**[IMPL]** Auditar especialmente:

- todo bloque `unsafe` y su safety invariant;
- FFI, layout y ownership cruzando ABI;
- pinning/self-reference;
- atomics y interior mutability;
- cancel safety de futures;
- blocking dentro del executor;
- allocator y panic strategy;
- bounds checks realmente eliminados o no.

La seguridad de Rust reduce clases de bugs; no prueba ausencia de deadlocks, starvation, errores lógicos, latencia o uso incorrecto de `unsafe`.

## 11. Diseño de hot path de baja latencia

### 11.1 Propiedades deseables

- trabajo acotado por evento;
- layout compacto y predecible;
- allocations y syscalls controladas;
- no logging síncrono en el camino crítico;
- ownership y afinidad explícitos;
- backpressure antes de memoria infinita;
- snapshots fuera del hot path;
- timestamps con reloj y dominio declarados;
- secuencia/gap detection cuando la fuente lo permita;
- degradación y shutdown deterministas.

### 11.2 Presupuesto de latencia

```text
ingest
+ parse/validate
+ normalize
+ state update
+ decision/compute
+ serialize
+ queue/network
+ scheduling/contention
= end-to-end
```

Medir cada etapa y también el total. La suma de medianas no es el p99 end-to-end; las dependencias y colas correlacionan tails.

### 11.3 Batching

Batching amortiza overhead y mejora vectorización/throughput, pero espera a completar el lote y puede aumentar latencia/edad de evento. Definir batch máximo, timeout/deadline y flush por prioridad. Comparar goodput dentro de SLO, no sólo throughput pico.

## 12. Observabilidad con costo conocido

Capas:

1. métricas de producto y SLO;
2. métricas del proceso/runtime;
3. scheduler, VM, filesystem y red;
4. perfiles CPU/off-CPU;
5. PMU: cycles, instructions, branches, cache/TLB;
6. tracing dinámico/estático;
7. logs y artefactos de incidente.

**[IMPL]** Antes de instrumentar, declarar overhead, sampling, cardinalidad, buffers y comportamiento si el collector falla. La observabilidad no debe bloquear el hot path ni filtrar secretos.

## 13. Seguridad de las optimizaciones

- `unsafe`, intrinsics y assembly amplían superficie de UB y portabilidad.
- Parsers rápidos siguen necesitando bounds, length y overflow checks.
- Desactivar mitigaciones puede cambiar el threat model; no es una optimización local inocua.
- Shared memory exige permisos, lifetime, limpieza y validación de mensajes.
- Huge pages, pinned memory y kernel bypass alteran aislamiento y operación.
- Cache keys, pools y buffers deben respetar tenant y clasificación de datos.

El manual `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md` gobernará los controles generales; este archivo conserva los riesgos propios de bajo nivel.

### 13.1 Cruce con microestructura y exchanges

`MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md` posee sequence, book, matching, order lifecycle y venue rules; este manual posee el costo físico. Para un feed/order path:

- un owner por book/shard debe alinearse con RX queue, CPU, NUMA y memoria sin inventar paralelismo dentro de la secuencia;
- fixed-point, endian, bounds y decoder correctness preceden SIMD/zero-copy;
- rings/pools tienen capacidad, edad y overflow policy; drop local invalida lineage y fuerza resync;
- medir capture→decode→apply→publish y order intent→wire/ACK/fill, con p99.9/max y clock quality;
- logging/checkpoint/analytics no bloquean el state owner; degradan explícitamente o propagan backpressure;
- busy-poll, bypass, huge pages o FPGA conservan recovery, observabilidad, aislamiento y risk gate.

Cruce recíproco auditado el 2026-08-21: optimizar messages/s sin hash/invariantes del libro no es un resultado válido.

## 14. Gates de código y rendimiento

### Corrección

- tests unitarios, property-based y diferenciales;
- invariantes y oráculos independientes;
- sanitizers/model checkers cuando corresponda;
- fault injection y shutdown/restart;
- misma salida o tolerancia numérica explícita.

### Rendimiento

- baseline versionado;
- carga y hardware declarados;
- muestras crudas conservadas;
- p50/p95/p99/p99.9 y máximos;
- CPU, memoria, bandwidth, faults, context switches y colas;
- comparación estadística y tamaño de efecto;
- carga sostenida y overload;
- performance budget en CI sólo en runners controlados.

### Mantenibilidad

- comentario explica por qué existe la optimización;
- test captura equivalencia y límites;
- benchmark reproduce el beneficio;
- fallback o rollback disponible;
- owner y fecha/arquitectura de validación registrados.

## 15. Runbook de diagnóstico

```text
1. Confirmar síntoma y ventana temporal.
2. Separar demanda, errores, latencia y goodput.
3. Comparar despliegue/configuración/host/región.
4. Aplicar utilización–saturación–errores por recurso.
5. Obtener perfil on-CPU y off-CPU.
6. Examinar colas, locks, faults, I/O y red.
7. Correlacionar traces con percentiles, no sólo promedios.
8. Formular una causa y predecir qué métrica cambiará.
9. Reproducir o ejecutar experimento seguro.
10. Mitigar/rollback antes de optimizar si el SLO está en riesgo.
11. Incorporar el fallo a tests, benchmark y observabilidad.
```

## 16. Contrato para agentes Codex

Al recibir una tarea de sistemas o rendimiento, el agente debe producir:

1. contrato funcional y no funcional;
2. diagrama del camino de datos;
3. presupuesto de latencia/recursos;
4. baseline y procedimiento reproducible;
5. hipótesis priorizadas por evidencia;
6. cambio mínimo y riesgos semánticos;
7. tests de equivalencia y concurrencia;
8. benchmark antes/después;
9. límites de portabilidad y seguridad;
10. rollback y seguimiento.

Prohibiciones:

- no inventar cifras de latencia de hardware;
- no declarar causalidad por una correlación aislada;
- no usar promedio para ocultar tails;
- no recomendar lock-free, SIMD, huge pages o kernel bypass sin workload;
- no atribuir `[IMPL]` a un profesor o estándar;
- no declarar mejora si cambió la semántica, calidad o carga.

## 17. MIT 6.172 — extracción operativa completa

Esta sección sintetiza las 23 clases y las prácticas públicas del curso de 2018. `[ACADEMIC]` significa paráfrasis fiel de sus slides y materiales oficiales; las frases literales se mantienen breves, entre comillas y en inglés. Cilk, Tapir, detalles de Haswell y cifras de experimentos del curso son casos históricos concretos, no recomendaciones universales para un sistema actual.

### 17.1 Cobertura lección por lección

| # | Tema oficial | Enseñanza que debe conservar un agente |
|---:|---|---|
| 1 | Introduction & Matrix Multiplication | El rendimiento no llega automáticamente con hardware moderno: algoritmo, lenguaje, representación, vectorización y paralelismo cambian el resultado. Optimizar no autoriza sacrificar corrección, confiabilidad o mantenibilidad. |
| 2 | Bentley Rules for Optimizing Work | Reducir trabajo mediante representación, precomputación, caché, loops, lógica y estructura de funciones; partir de código correcto y conservar regression tests. Menos operaciones es una heurística fuerte, no una garantía de menor tiempo. |
| 3 | Bit Hacks | Aprovechar word-level parallelism y operaciones de bits sólo cuando semántica, portabilidad, claridad y código generado lo justifican; primero comprobar qué reconoce el compilador. |
| 4 | Assembly & Architecture | Entender pipelines, dependencias verdaderas, hazards, superscalar/out-of-order, renaming, branch prediction y SIMD. Una cifra de una microarquitectura no debe convertirse en ley general. |
| 5 | C to Assembly | Seguir C/C++ → LLVM IR/SSA/CFG → assembly; usar inlining, LTO y reportes del compilador con evidencia. Los passes cooperan, pero aliasing, UB, tamaño y límites de análisis restringen transformaciones. |
| 6 | Multicore Programming | Expresar paralelismo de tareas a un nivel que permita scheduling y análisis; Cilk ilustra abstracción y work stealing, pero sus ideas deben traducirse al runtime moderno elegido. |
| 7 | Races and Parallelism | Modelar DAG, work `T1`, span `T∞` y paralelismo `T1/T∞`; detectar determinacy races y entender bounds del scheduler bajo sus supuestos. |
| 8 | Analysis of Multithreaded Algorithms | Resolver recurrences de work/span y analizar loops paralelos; una API paralela no garantiza suficiente paralelismo ni bajo overhead. |
| 9 | What Compilers Can and Cannot Do | El compilador aplica transformaciones mecánicas y conservadoras; la abstracción puede desaparecer, pero sólo IR, assembly, diagnósticos y medición lo demuestran. |
| 10 | Measurement and Timing | Controlar variabilidad, elegir reloj/herramienta/estadística según la pregunta, comparar hipótesis con ensayos y validar modelos fuera de la muestra usada para ajustarlos. |
| 11 | Storage Allocation | Distinguir stack/heap, free lists por tamaño, fragmentación y familias de GC. Costos amortizados no son garantías de tail latency. |
| 12 | Parallel Storage Allocation | Examinar heaps locales, ownership, remote free y GC stop-the-world/incremental/parallel/concurrent; read y write barriers trasladan costos, no los eliminan. |
| 13 | Cilk Runtime System | Aplicar el principio *work-first*: mantener barata la ejecución ordinaria y trasladar overhead al camino de steal cuando el modelo lo favorece; conocer lo que realmente hace el runtime. |
| 14 | Caching and Cache-Efficient Algorithms | Usar associativity, blocking y divide-and-conquer bajo un modelo de caché explícito; minimizar misses seriales puede beneficiar ejecución paralela, pero hardware real añade NUMA, coherencia y prefetch. |
| 15 | Cache-Oblivious Algorithms | Diseñar recursión/localidad que se adapte a varios niveles sin fijar un tile, y comparar sus bounds con cache-aware. Tall-cache y modelo ideal son supuestos, no propiedades automáticas del host. |
| 16 | Nondeterministic Parallel Programming | Preferir determinismo. La regla literal central es “Never write nondeterministic parallel programs”; si el requisito obliga, encapsularlo y diseñar una estrategia de prueba deliberada. |
| 17 | Synchronization Without Locks | Razonar desde el memory model; fences y atomics tienen semántica y costo. CAS permite estructuras lock-free, pero trae coherencia, contención, starvation, reclamación y ABA. |
| 18 | Domain-Specific Languages and Autotuning | Separar algoritmo de schedule permite expresar intención, codificar conocimiento experto y explorar implementaciones. Un autotuner necesita espacio, objetivo, presupuesto y validación fuera del conjunto de tuning. |
| 19 | Leiserchess Codewalk | Comprender primero representación, flujo, invariantes, search y tests de un sistema real. Probar con frecuencia; empezar por optimizaciones que no alteran la búsqueda y mejorar heurísticas existentes antes de inventar otras. |
| 20 | Speculative Parallelism & Leiserchess | La especulación intercambia paralelismo por trabajo posiblemente inútil. Generarla sólo si falta paralelismo no especulativo y es probable necesitarla; medir aborts, wasted work y dependencia del scheduling. |
| 21 | Tuning a TSP Algorithm | La mejora mayor proviene de ingeniería algorítmica incremental: partial sums, pruning, orden de pruebas, precomputación, lazy evaluation, bit sets y mejores bounds, sostenidos por perfiles y cuaderno experimental. |
| 22 | Graph Optimization | Grafos reales suelen ser grandes, sparse, irregulares y memory-bound; representación, frontier direction, compresión y reordenamiento explotan localidad, pero la mejor opción depende del grafo. |
| 23 | High Performance in Dynamic Languages | El enfoque de dos lenguajes añade complejidad; specialization, JIT, tipos inferibles, metaprogramación y generación de código pueden conservar productividad y lograr kernels rápidos. Vectorizar no resuelve todo. |

### 17.2 Disciplina de optimización

**[ACADEMIC]** El curso combina tres movimientos, en este orden lógico:

1. establecer semántica correcta y tests;
2. reducir/reorganizar trabajo y escoger representación;
3. adaptar el trabajo restante a compilador, memoria y paralelismo.

Las reglas de Bentley quedan operacionalizadas así:

| Familia | Técnicas | Pregunta de revisión |
|---|---|---|
| datos | packing/encoding, augmentation, precomputation, compile-time initialization, caching, lazy evaluation, sparsity | ¿qué bytes y resultados podemos no producir, no mover o no recalcular? |
| loops | hoisting, sentinels, unrolling, fusion, eliminar iteraciones inútiles | ¿qué trabajo del loop es invariante, redundante o puede terminar antes? |
| lógica | folding/propagation, CSE, identidades, short-circuit, reordenar/componer tests, fast path | ¿qué condición barata descarta antes la mayor cantidad de trabajo? |
| funciones | inlining, tail-recursion elimination, coarsening | ¿la frontera de abstracción impide una transformación material y demostrable? |

**[IMPL]** Para cada técnica registrar cantidad de trabajo, bytes tocados, branch behavior, code size y tiempo. Si una reduce instrucciones pero aumenta misses o mispredictions, la medición decide.

### 17.3 Del paralelismo teórico al runtime

```text
T1 = trabajo total en un procesador
T∞ = longitud del camino crítico (span)
paralelismo disponible = T1 / T∞
límite inferior ideal: TP >= max(T1/P, T∞)
```

**[ACADEMIC]** Work/span separa trabajo insuficiente de scheduling deficiente. Work stealing y los bounds presentados en el curso dependen de un modelo; no incluyen automáticamente NUMA, bandwidth, oversubscription, preemption ni I/O. La lección del runtime culmina en el breve principio literal “work-first principle”: optimizar el camino común serial aun si un steal raro requiere trabajo adicional.

**[IMPL]** Antes de paralelizar:

- calcular o estimar `T1`, `T∞`, granularidad y máximo speedup útil;
- medir `T_serial`, `T_1_worker` y `T_P`; la diferencia entre los dos primeros expone overhead del runtime;
- inspeccionar task creation, steals, imbalance, synchronization, bandwidth y false sharing;
- aumentar grain size o coarsen sólo con una curva que muestre el punto de equilibrio;
- tratar Cilk/Tapir como caso de estudio; documentar scheduler y garantías de la biblioteca real.

### 17.4 Medición: lo que la lección 10 permite y lo que no

**[ACADEMIC]** Las fuentes de ruido enumeradas incluyen daemons, interrupts, alignment, placement, scheduler, SMT, multitenancy, DVFS y Turbo. Quiescing, pinning o desactivar mecanismos puede aislar una hipótesis, pero describe un entorno experimental concreto.

El curso usa el mínimo repetido para aproximar el tiempo sin ruido en un programa determinista cuando el ruido sólo puede elevar la duración. **[IMPL] Corrección:** esa elección no representa experiencia de usuario, colas, jitter ni tail latency. Para producto se informa la distribución y se selecciona la estadística que corresponde al SLO.

- usar media geométrica para combinar ratios positivos comparables, no media aritmética de speedups;
- intercalar variantes cuando el entorno deriva;
- separar exploratory profiling de confirmatory benchmark;
- reservar workloads/inputs de validación para evitar optimizar el benchmark;
- no interpretar `R²` alto como causalidad ni extrapolación válida;
- conservar intentos que empeoraron: son evidencia sobre el modelo causal.

### 17.5 Allocation, caché y localidad

**[ACADEMIC]** Las lecciones 11–15 conectan lifetime, allocator y jerarquía de memoria. Un free list o heap local cambia búsqueda, sincronización, fragmentación y ownership. Un algoritmo cache-aware fija parámetros del target; uno cache-oblivious intenta obtener localidad recursiva entre niveles, bajo ideal-cache y tall-cache assumptions.

**[IMPL]** Para una optimización de memoria exigir:

```text
layout + lifetime + access pattern + working set
→ bytes/op + allocations/op + cache/TLB + bandwidth + NUMA
→ latencia completa, incluidos faults, reclaim, GC y remote free
```

La compresión de grafos ilustra un trade-off esencial: decodificar añade cómputo, pero puede acelerar un workload memory-bound al reducir tráfico y contención de memoria. No generalizarlo sin medir por dataset.

### 17.6 Determinismo, locks y lock-free

**[ACADEMIC]** Una ausencia de determinacy races en el modelo del curso da determinismo para ese input; estar libre de data races no elimina todo nondeterminismo introducido por el orden de locks. Las estrategias presentadas para nondeterminismo intencional son: desactivarlo en test, encapsularlo, sustituir una alternativa determinista o aplicar herramientas de análisis.

**[SPEC/IMPL]** La exposición de 2018 sobre C11/x86 debe reconciliarse con el estándar y compilador actuales antes de escribir código. Las reglas permanentes son:

- separar atomicidad, mutual exclusion, ordering y progress;
- especificar la relación happens-before y el punto de linearización;
- usar el ordering mínimo sólo después de demostrar su suficiencia;
- medir el costo de coherencia: CAS sobre una línea compartida invalida copias y puede escalar mal;
- distinguir lock-free de wait-free; un thread puede starve en el primero;
- impedir ABA mediante versioning apropiado o, más generalmente, reclamación segura que impida reutilización prematura;
- someter estructuras no bloqueantes a model checking/stress, sanitizers compatibles y pruebas de reclamación.

### 17.7 DSL, schedule y autotuning

**[ACADEMIC]** GraphIt y Halide ejemplifican la separación entre *qué* computar y *cómo* programarlo sobre el hardware. OpenTuner ejemplifica búsqueda en espacios demasiado grandes para exhaustive search. Esto permite codificar transformations específicas del dominio sin obligar al programador de la aplicación a deshacer decisiones low-level prematuras.

**[IMPL]** Contrato mínimo de autotuning:

```text
semántica/invariantes fijas
espacio de schedules y restricciones válidas
función objetivo multi-métrica
hardware, datos y presupuesto de búsqueda
repeticiones y control de ruido
holdout de entradas y entorno
fallback portable + caché versionada de resultados
```

No aceptar el ganador del tuner si sólo memoriza inputs, altera precisión/corrección o empeora tails, energía, memoria o compilación fuera del objetivo declarado.

### 17.8 Casos integradores y aprendizaje de errores

Leiserchess, TSP, grafos y lenguajes dinámicos enseñan a optimizar sistemas completos:

- **Leiserchess:** representación → move generation → evaluation → search → ordering → transposition cache → tests; primero serial y correcto, después paralelismo especulativo.
- **TSP:** una sucesión de algoritmos medidos supera al ajuste aislado; pruning y mejores bounds evitan factoriales enteros de trabajo antes de tocar instrucciones.
- **grafos:** cambiar traversal direction, layout, compression y IDs según densidad/frontier/dataset; ninguna configuración domina todos los grafos.
- **lenguaje dinámico:** specialization/code generation puede eliminar despacho y unboxing en kernels; mantener escape hatch e inspeccionar el código generado.

**[IMPL] Registro obligatorio de aprendizaje:**

```text
hipótesis:
cambio y diff:
predicción observable:
correctness oracle:
workload/entorno:
resultado y datos crudos:
por qué funcionó o falló:
qué supuesto quedó refutado:
test/benchmark añadido:
decisión: conservar | revertir | investigar
```

Un intento fallido es útil si reduce el espacio de hipótesis y queda reproducible. No racionalizar después del resultado: escribir la predicción antes de ejecutar.

### 17.9 Práctica profesional extraída de tareas y proyectos

**[ACADEMIC]** Las 10 tareas progresan por tooling/C, profiling, vectorización, Cilk/reducers, teoría, allocator, análisis dinámico, cachegrind, determinismo y sincronización. Los cuatro proyectos cubren bit hacks, collision detection, allocator y Leiserchess.

El patrón pedagógico que se adopta para trabajo profesional es:

1. beta funcional y correctness suite;
2. perfil y modelo del cuello de botella;
3. design/code review por otra persona o agente;
4. optimización incremental con medición antes/después;
5. informe de bottlenecks, cambios, speedups y técnicas fallidas;
6. evaluación sobre artefacto final y carga no usada para ajustar.

Por tanto, un agente no entrega sólo código “rápido”: entrega evidencia, límites, pruebas y explicación causal revisable.

## 18. Sistemas operativos — triangulación MIT 6.S081, CS:APP y OSTEP

### 18.1 Qué aporta cada fuente

| Fuente | Perspectiva | Cobertura verificada | Límite que debe respetarse |
|---|---|---|---|
| MIT 6.S081/6.828 (2021) | implementar un OS pequeño sobre RISC-V | 25 clases; xv6 y labs de utilities, syscalls, page tables, traps, COW, threads, driver de red, locks, filesystem y `mmap` | xv6 privilegia claridad pedagógica; no representa toda la complejidad, seguridad ni rendimiento de Linux/Windows/BSD |
| CS:APP 3e / CMU | el sistema completo desde el programador | bits, machine code, processor, optimización, memoria, linking, exceptional control, VM, allocator, I/O, sockets y concurrencia; labs sobre programas reales | sus casos x86-64/Linux y edición concreta deben versionarse frente al target actual |
| OSTEP 1.10 | conceptos, mecanismos y políticas de OS | 57 capítulos en virtualización, concurrencia, persistencia, distribución y seguridad, con simuladores, code, homework y proyectos | explica modelos y diseños; una política real del kernel exige su documentación/código y medición |

**[ACADEMIC]** El principio explícito de CS:APP es estudiar hardware, OS, compiler y network como sistema completo, y aprender desarrollando/evaluando programas reales sobre máquinas reales. La consecuencia operativa es directa: una explicación que termina en la API está incompleta si el costo o la corrección dependen de traducción, kernel, runtime, linker o red.

### 18.2 Mapa curricular reconciliado

| Problema | MIT 6.S081 | CS:APP | OSTEP | Resultado para implementación |
|---|---|---|---|---|
| proceso y syscall | organización, entry/exit, traps, interrupts | exceptional control flow, procesos, signals | process API, direct execution | trazar user→kernel→user, errores, blocking, copia y scheduling |
| direcciones y memoria | page tables, faults, COW, `mmap` | address translation, TLB, mapping, allocator | segmentation, paging, TLB, page tables, swapping | separar reserva, commit, resident set, fault, reclaim y protección |
| CPU y scheduling | context switch, sleep/wakeup, multiprocessor | concurrency/parallelism y control flow | FIFO/SJF/STCF/RR, MLFQ, lottery y multi-CPU | medir runnable time, wait, preemption, migration y fairness |
| sincronización | spinlocks, sleeplocks, lock lab, RCU | threads, semaphores, synchronization | locks, CV, semaphores, bugs, event-based concurrency | demostrar invariantes, ordering, progress y shutdown; perfilar contención |
| almacenamiento | buffer cache, filesystem, log, crash recovery | Unix I/O, files y robust I/O | devices, disks, RAID, FFS, journaling, LFS, SSD, integrity | distinguir cache, completion, ordering, atomicity y durability |
| red | NIC driver, interrupts, receive livelock | sockets/client-server | distributed systems, NFS y AFS | diseñar bounded queues, backpressure y política frente a overload |
| aislamiento/seguridad | privilege, Meltdown, VM/OS organization | buffer overflow, memory bugs, modes | auth, access control y cryptography | convertir boundary, validation y least privilege en requisitos, no anexos |
| práctica | labs que modifican xv6 | Data/Bomb/Attack/Arch/Perf/Cache/Shell/Malloc/Proxy labs | simuladores y proyectos | entregar código, pruebas, medición, fallos y explicación del mecanismo |

### 18.3 Modelo de ejecución user–kernel

```text
user code
→ runtime/libc
→ syscall instruction / exception / interrupt
→ trap entry + cambio de privilegio
→ validación y copia de argumentos
→ subsistema del kernel
→ posible bloqueo, wakeup y context switch
→ driver/dispositivo u otro proceso
→ retorno y restauración de contexto
```

**[IMPL]** Para un camino crítico registrar:

- syscalls por operación y bytes por llamada;
- copies y cambios de representación;
- fallos de página, faults evitables y TLB behavior;
- runnable/wait/off-CPU time y migrations;
- locks del proceso y del kernel visibles por tracing;
- interrupciones/softirqs, queue depth y completions;
- errores parciales, cancelación y señalización.

Agrupar llamadas o usar shared memory puede reducir transiciones, pero cambia latencia, aislamiento, backpressure, lifetime y threat model. La medición y la especificación completa deciden.

### 18.4 Memoria virtual sin mitos

**[ACADEMIC]** Las tres fuentes conectan address space, page tables, TLB, faults, protection y allocation. MIT 6.S081 permite observarlos en código; CS:APP sigue la traducción end-to-end; OSTEP separa mecanismos de paging/translation de políticas de swapping.

Reglas operativas:

- memoria virtual reservada no equivale a RAM residente;
- un acceso puede provocar page walk, minor fault, major fault, copy-on-write o protección;
- `fork`/COW difiere costos hasta la primera escritura y puede crear tails posteriores;
- `mmap` cambia ownership/faulting/writeback, pero no elimina I/O;
- un allocator puede retener páginas aunque los objetos se hayan liberado;
- RSS, PSS, committed, mapped y working set responden preguntas distintas;
- TLB reach depende de page size y working set; huge pages tienen beneficios y costos;
- NUMA placement y first-touch no aparecen en el modelo simple de xv6.

**[IMPL] Gate:** antes de culpar a “memoria”, correlacionar allocation rate, residency, faults, reclaim, page-table/TLB, NUMA y bandwidth con la ventana del síntoma.

### 18.5 Scheduling y coordinación

**[ACADEMIC]** xv6 muestra scheduler, context switch y sleep/wakeup; OSTEP compara políticas y explica multi-CPU scheduling; CS:APP conecta procesos, signals y threads con programas reales.

Objetivos que pueden entrar en conflicto:

| Objetivo | Posible costo |
|---|---|
| baja latencia interactiva | menor throughput o más preemption |
| fairness | peor cache affinity o tail de prioridad alta |
| afinidad/localidad | desbalance y hotspots |
| busy wait | cores/energía y starvation del vecino |
| work conservation | interferencia sobre tarea crítica |
| prioridad estricta | priority inversion o starvation |

**[IMPL]** No diagnosticar una pausa sólo con CPU%. Obtener on-CPU, runnable, sleeping, blocked, involuntary context switches, wakeup latency y dependency chain. Un thread puede usar poca CPU porque espera un lock, I/O, memoria, quota o productor upstream.

### 18.6 Persistencia y crash consistency

**[ACADEMIC]** MIT 6.S081 recorre filesystem, logging, crash recovery y fast recovery; OSTEP contrasta fsck, journaling, LFS, SSD e integridad; CS:APP fija la semántica de Unix I/O observada por el programa.

Separar siempre:

```text
write aceptado por una API
≠ datos copiados al page cache
≠ request enviado al dispositivo
≠ completion del dispositivo
≠ orden persistente requerido
≠ transacción recuperable tras crash
```

Contrato de datos persistentes:

- unidad de atomicidad y orden entre datos/metadata;
- qué significa ACK o commit;
- política de flush/barrier y comportamiento del hardware;
- checksum y detección de torn/corrupt writes;
- replay/idempotencia del log;
- recuperación tras power loss, process crash y partial write;
- prueba por fault injection en puntos de corte relevantes.

Un benchmark que omite `fsync`/durability no puede compararse con otro que confirma persistencia.

### 18.7 Del modelo pedagógico al kernel real

Antes de trasladar una conclusión de xv6 o de un capítulo al target:

1. identificar mecanismo conceptual y supuesto;
2. localizar documentación y, cuando sea necesario, código de la versión real del kernel/runtime;
3. comprobar configuración, cgroup/container, filesystem, driver y hardware;
4. observar el camino con tracing/perfil, no inferirlo sólo por nombres de API;
5. diseñar un experimento que cambie una causa predicha;
6. conservar diferencia entre comportamiento normativo, detalle de implementación y resultado medido.

### 18.8 Prácticas que convierten teoría en capacidad profesional

- **xv6 labs:** modificar el mecanismo y pasar tests hace visibles invariantes que una lectura no revela.
- **Shell Lab:** procesos, job control, signals y carreras de control.
- **Malloc Lab:** layout, free lists, coalescing y trade-off tiempo/espacio.
- **Cache/Performance Lab:** dirección real, locality y transformación medida.
- **Proxy Lab:** byte order, I/O, procesos/threads, sockets y concurrencia integrados.
- **OSTEP simulators:** variar política/carga y contrastar predicción con ejecución.

**[IMPL]** Un agente que use este manual debe poder construir al menos un shell con job control, un allocator instrumentado, un cache simulator, un servicio concurrente con backpressure y una prueba de crash recovery. Conocer definiciones sin poder preservar invariantes en código no cierra la competencia.

### 18.9 Correcciones contra Linux actual

**[SPEC/IMPL]** La documentación vigente del kernel confirma que la implementación real es una familia de subsistemas y configuraciones, no “el scheduler” o “la memoria” en abstracto:

- Linux comenzó la transición de CFS hacia EEVDF desde 6.6; EEVDF usa lag/eligibility y virtual deadlines. Diagnósticos basados sólo en textos históricos de CFS deben revisar kernel/configuración efectivos.
- El índice de scheduling también contiene deadline, real-time groups, bandwidth control, capacity/energy awareness, utilization clamping y extensible scheduler; la clase de la tarea importa.
- Memory Management incluye page tables, reclaim, swap, page cache, OOM, THP/HugeTLB, migration, NUMA y allocators; “available memory” no identifica por sí sola el mecanismo activo.
- La guía de memory barriers del kernel se declara guía incompleta, no especificación infalible; las primitivas Linux y su formal memory model gobiernan código kernel, no sustituyen el estándar C++/Rust en userspace.
- JBD2/ext4 protege la consistencia de las transacciones que pasan por el journal; no convierte cualquier `write()` de aplicación en commit durable ni garantiza que data y metadata tengan idéntica política.

**[IMPL] Manifiesto Linux mínimo:** `uname/kernel build`, config relevante, scheduler class/policy, cgroup v2 y quotas, NUMA/THP, filesystem/mount options, block device/cache, container/VM y mitigaciones. Sin ese manifiesto, una afirmación sobre “Linux” queda incompleta.

### 18.10 Código Linux — `io_uring` como protocolo de ownership

**[CODE] Evidencia fijada:** `torvalds/linux@26260251022fbc2f248a3d747a9b2b961b18d2d8`, GPL-2.0 con Linux-syscall-note para UAPI. La implementación mantiene submission/completion rings, recursos registrados, cancelación, polling y rutas explícitas de overflow. El patrón transferible no es “menos syscalls”, sino un protocolo de estados entre aplicación y kernel:

```text
preparar SQE + publicar tail
→ kernel acepta/rechaza/ejecuta
→ CQE identifica resultado por user_data
→ aplicación consume head y libera/reutiliza recursos
```

Obligaciones que el código hace visibles:

- dimensionar y drenar la CQ: existen contadores/flags de `sq_dropped` y `cq_overflow`; una aplicación que deja de consumir completions pierde control de progreso aunque la SQ siga aceptando trabajo;
- `user_data` identifica una operación lógica, no debe apuntar a memoria cuyo lifetime termina antes del CQE;
- fixed files, buffers registrados y buffer selection reducen lookup/pinning sólo si registration, update, unregister y shutdown tienen ownership inequívoco;
- linked requests cambian cancelación y propagación de error; probar fallo en cada eslabón y no asumir una transacción;
- `SQPOLL`, `IOPOLL`, `SINGLE_ISSUER` y `DEFER_TASKRUN` alteran threads, wakeups y requisitos de progreso; fijar combinaciones soportadas por kernel/liburing y medir CPU/energía además de latencia;
- cancelar no implica que la operación nunca empezó ni que un side effect externo sea reversible. Reconciliar outcomes ambiguos igual que en cualquier API asíncrona;
- close/exec/teardown debe impedir nuevas submissions, cancelar lo cancelable, seguir drenando completions y liberar sólo después de demostrar ausencia de usuarios.

**Gate:** comparar contra `epoll`/threads o API simple con idéntica semántica; medir submit→complete p50/p99.9, completions por wakeup, ring occupancy/overflow, CPU por rol, allocations, syscalls, cancel latency y shutdown bajo burst, slow device, partial I/O, timeout y consumer pausado. `io_uring` se adopta sólo si mejora el cuello real sin debilitar lifetime, backpressure, observabilidad o recovery.

## 19. Metodología de producción — Brendan Gregg y Linux

### 19.1 Empezar por preguntas, no por herramientas

**[PROD]** La USE Method construye una lista de recursos y, para cada uno, busca **utilization, saturation y errors**. Su ventaja es exponer preguntas sin herramienta o métrica disponible como *known unknowns*. Su propio autor la presenta como una herramienta temprana para bottlenecks/errores, no como método universal ni prueba de causalidad.

| Método | Pregunta que responde | Señal inicial |
|---|---|---|
| Problem Statement | ¿qué evidencia define el problema y qué cambió? | síntoma, alcance, versiones, baseline sano |
| Workload Characterization | ¿quién genera qué carga, con qué distribución? | rate, tamaño, mezcla, locality, concurrencia |
| USE | ¿qué recurso está ocupado, saturado o fallando? | CPU, memoria, red, storage, buses/interconnects |
| RED | ¿qué servicio degrada solicitudes? | rate, errors, duration |
| CPU Profile | ¿dónde se consume CPU? | stacks ponderados por samples |
| Off-CPU | ¿por qué esperan los threads? | stacks y duración de espera |
| TSA | ¿en qué estados transcurre el tiempo del thread? | executing, runnable, sleep/I/O, lock, idle |
| Active Benchmarking | ¿qué limita el benchmark mientras corre? | perfiles y recursos simultáneos al ensayo |
| Time/phase division | ¿en qué ventana o fase aparece el costo? | timeline alineada con eventos/cambios |

**[IMPL] Secuencia de incidente:**

```text
definir síntoma y boundary
→ caracterizar workload
→ RED por servicio + USE por recurso
→ CPU/off-CPU/TSA según el tiempo observado
→ tracing dirigido por hipótesis
→ experimento/mitigación
→ verificar SLO, efectos secundarios y causa
```

Encontrar un recurso saturado no demuestra que sea la única causa. Continuar hasta conectar demanda → cola/espera → stack/evento → resultado de usuario.

### 19.2 Métrica, perfil y trace no son intercambiables

| Instrumento | Conserva | Pierde o arriesga |
|---|---|---|
| counter/metric | tendencia barata y agregable | orden causal y distribución interna |
| histogram/heatmap | distribución y cambio temporal | stack causal si no se adjunta |
| sampled profile | atribución aproximada por stack | eventos raros; sesgo/frecuencia de sample |
| trace | orden, duración y correlación por evento | volumen, overhead, drops y privacidad |
| flame graph | stacks agregados y ancho proporcional | timeline, orden de eventos y, solo, causalidad |

Reglas:

- CPU flame graph para on-CPU; off-CPU flame graph para esperas: no mezclarlos sin etiqueta;
- perf counters se multiplexan si se piden más eventos que slots; registrar escala y precisión;
- tracepoints son hooks explícitos del kernel; kprobes/uprobes dependen más de implementación;
- BPF amplía observabilidad, pero verifier, mapas, sampling y attach points forman parte del experimento;
- instrumentar puede cambiar scheduling, caché y tails; medir overhead y eventos perdidos;
- `perf_events` puede exponer paths, PIDs, direcciones y actividad sensible; respetar `CAP_PERFMON`, scopes y política de datos.

### 19.3 Checklist de 60 segundos y escalamiento

Una primera pasada debe ser corta y repetible:

1. demanda, errores y latencia del servicio;
2. CPU por core, run queue y throttling;
3. memoria disponible, faults, reclaim y swapping;
4. red por interfaz: rate, drops, retransmisiones y backlog;
5. storage: utilization, queue/wait, latency y errores;
6. despliegues/configuración/vecinos durante la ventana;
7. perfiles on/off-CPU si lo anterior no explica el síntoma.

No congelar nombres de comandos como conocimiento eterno. Resolver herramienta y campo contra la versión desplegada; el checklist debe declarar cómo obtener cada métrica y qué significa en ese kernel/runtime.

## 20. Arquitectura, microarquitectura y PMU

### 20.1 Tres contratos distintos

| Nivel | Define | Fuente dominante |
|---|---|---|
| ISA/arquitectura | instrucciones, registros, exceptions, memory/system semantics | Intel SDM, AMD APM, Arm Architecture Reference Manual |
| microarquitectura | pipelines, ports, caches, predictors, buffers y ejecución concreta | optimization manual y guía del modelo/familia |
| plataforma | NUMA, firmware, memory channels, devices, power/thermal | fabricante de CPU/board, ACPI/firmware y medición |

**[SPEC]** Intel y AMD documentan x86-64 en familias de manuales; Arm documenta A-profile y su memory model. La compatibilidad de ISA no implica mismos counters, latencias, cache topology ni mejores sequences.

**Verificado el 2026-08-21:** Intel SDM v092 y Optimization Reference Manual v050; AMD64 APM combinado rev. 4.09 (2026-03-09), con Vol. 2 System Programming rev. 3.44 y Vol. 4 incluyendo 128/256/512-bit media instructions rev. 3.26. Las revisiones se deben volver a verificar al usar el manual.

### 20.2 Procedimiento de optimización por CPU

```text
feature discovery (CPUID/HWCAP)
→ target/ABI y dispatch
→ compiler report + IR + assembly
→ benchmark por tamaño/distribución
→ cycles/instructions + top-down/eventos pertinentes
→ cache/TLB/bandwidth/NUMA
→ frecuencia, energía, temperatura y throttling
→ validar en cada familia soportada
```

Reglas:

- compilar `-march=native` sólo para un artefacto ligado al host; distribución requiere baseline y dispatch/versiones;
- una instrucción SIMD más ancha puede reducir instrucciones y aun bajar frecuencia, aumentar tail o no mejorar por bandwidth;
- IPC no es objetivo aislado: puede bajar al hacer menos trabajo y mejorar el tiempo;
- raw PMU event codes son model-specific; preferir eventos nombrados/versionados y documentar family/model/stepping;
- separar frontend bound, bad speculation, backend bound y retiring como clasificación inicial, no como causa final;
- `cycles`, reference cycles y wall time responden preguntas distintas bajo DVFS/turbo;
- SMT comparte recursos; medir throughput total y tail por tarea con SMT on/off antes de decidir;
- binding sin topología puede cruzar NUMA o compartir core físico accidentalmente.

### 20.3 Memory ordering portátil

El hardware puede reordenar y propagar memoria de modo distinto entre x86-64 y Arm. El compilador también transforma bajo el memory model del lenguaje.

```text
semántica C++/Rust
→ IR atomics/fences de LLVM
→ lowering por target/ABI
→ instrucción, fence o libcall
→ coherencia y memory system del procesador
```

**[IMPL]** Nunca razonar sólo desde assembly x86. Primero probar happens-before en el lenguaje; después comprobar lowering y costo en x86-64/AArch64 u otros targets. Una operación atomic de tamaño/alineación no soportados nativamente puede requerir runtime/libcall; comprobar las garantías de lock-freedom del tipo/target.

## 21. C++, LLVM y Rust — contrato de concurrencia de bajo nivel

### 21.1 C++: ownership y sincronización antes de atomics débiles

**[PROD]** Las C++ Core Guidelines son guía evolutiva, no el estándar ISO. Sus reglas pertinentes son evitar data races, minimizar writable sharing, pensar en tasks, no usar `volatile` para sincronización y validar con herramientas. Para locks: RAII, adquisición conjunta de múltiples mutexes, no llamar código desconocido bajo lock, esperar con predicado, minimizar critical section y asociar mutex con los datos protegidos.

**[IMPL] Orden de preferencia:**

1. ownership único e inmutabilidad;
2. message passing/queues con backpressure;
3. estructuras de alto nivel y locks con invariantes claras;
4. atomics `seq_cst` si hace falta compartir a bajo nivel;
5. orderings débiles sólo con prueba escrita y tests específicos;
6. lock-free custom sólo si progress y benchmark lo exigen.

### 21.2 LLVM IR: puente, no especificación del programa

**[SPEC]** LLVM documenta `load/store atomic`, `cmpxchg`, `atomicrmw` y `fence`, con orderings desde `unordered/monotonic` hasta acquire/release/acq_rel/seq_cst. Su guía actual declara intención de implementar atomics de C++20. `volatile` y atomic son propiedades ortogonales en IR.

Consecuencias:

- optimizaciones permitidas cambian según ordering;
- acquire sólo sincroniza mediante una relación válida con release;
- `relaxed` preserva atomicidad/modification order del objeto, no publica el resto del estado;
- `cmpxchg` tiene ordering de éxito y de fallo;
- el target decide fences/instrucciones/libcalls y puede hacer imposible cierto atomic inline;
- inspeccionar IR sirve para explicar lowering, no reemplaza la prueba semántica del source language.

### 21.3 Rust: safety del tipo y frontera `unsafe`

**[SPEC]** La documentación estable de Rust declara que sus atomics siguen actualmente las reglas de C++20 sin `consume`; la Rust Reference advierte que el memory model general está incompleto. Ownership, `Send` y `Sync` previenen muchas clases de sharing inválido, pero no prueban algoritmo, progress, deadlock ni invariantes dentro de `unsafe`.

Checklist para `unsafe`/concurrencia:

- safety invariant documentada junto al bloque/API;
- aliasing, provenance, initialization, alignment y lifetime;
- atomic sizes disponibles en cada target;
- ordering de éxito/fallo y publicación del payload;
- reclamación y destrucción durante shutdown/cancel;
- unwind/panic y callbacks bajo locks;
- Loom/Miri/sanitizers/stress cuando sean aplicables, más revisión manual.

### 21.4 Plantilla de prueba happens-before

```text
estado protegido:
writer: operaciones ordinarias → atomic release sobre X
reader: atomic acquire de X que observa la publicación → lecturas ordinarias
invariante publicada:
punto de linearización:
orden de fallo de compare_exchange:
progress: blocking | lock-free | wait-free (alcance exacto):
reclamación:
targets/alineaciones/tamaños:
test/model checker/benchmark:
```

Si no puede completarse esta plantilla, usar una primitiva más fuerte. Debilitar orderings después de medir y demostrar, no antes.

### 21.5 Reclamación segura: el nodo removido todavía puede estar vivo

Desvincular un nodo de una estructura no demuestra que ningún reader conserve su dirección. Liberarlo o reutilizarlo prematuramente produce use-after-free y puede convertir un cambio `A→B→A` en un CAS aparentemente válido (ABA).

| Esquema | Protección | Fortaleza | Riesgo operativo |
|---|---|---|---|
| reference counting | cada referencia viva incrementa ownership | lifetime local y comprensible | RMW/cache traffic; ciclos; decremento final/destructor en hot path |
| hazard pointers | cada reader publica los objetos que podría dereferenciar; el reclaimer escanea antes de liberar | reclamación portable y progreso del core diseñado para HP | slots/scan, fences, retire threshold y protocolo de revalidación |
| epoch/QSBR | reader anuncia región/epoch o quiescent state; se libera tras avance global seguro | read path amortizado y buen throughput | reader detenido/pinneado retrasa reclaim y puede hacer crecer memoria |
| RCU | removal y reclamation separadas por grace period; API define publicación/dereference | reads extremadamente baratos en escenarios read-mostly | variantes y reglas de blocking distintas; update/reclaim y stalls pueden ser costosos |

**[ACADEMIC/SPEC]** El paper de Maged Michael formaliza hazard pointers para reclamación de objetos lock-free. P2530R3 llevó una interfaz de hazard pointers al working paper de C++26; la disponibilidad real se comprueba mediante standard/library/toolchain y feature-test, no por fecha. La documentación de Linux define RCU como removal primero y reclamation sólo después de un grace period que cubra a readers anteriores.

Obligaciones comunes:

1. **publication:** el reader debe publicar protección antes de poder usar el objeto;
2. **revalidation:** después de publicar, volver a leer el enlace y reintentar si cambió;
3. **retirement:** remover lógicamente antes de colocar en retire list;
4. **safe point:** demostrar qué observación prueba que ningún reader antiguo lo usa;
5. **reuse:** no devolver al allocator/pool antes del safe point;
6. **shutdown:** drenar callbacks/retired nodes y threads registrados;
7. **boundedness:** definir qué pasa si un thread pausa, muere o no reporta quiescence;
8. **progress:** no atribuir al conjunto un grado de progress superior al componente más débil;
9. **testing:** forzar pausas entre load/protect/revalidate/CAS/free y ejecutar reclamación agresiva;
10. **measurement:** medir memory high-water, retire backlog, scan/grace time, throughput y tails.

RCU es especialmente apropiado para read-mostly; la documentación del kernel advierte que puede no ser adecuado en cargas write-heavy. Elegir reclamación es parte del algoritmo y del SLO de memoria, no una tarea posterior.

## 21.6 Práctica pública Meta/Linux — scheduling específico del workload

**[PROD]** Meta documentó una política `sched_ext`/BPF para ads serving que separa dinámicamente threads del request path crítico y trabajo menos sensible, buscando localidad de LLC y prioridad de tails. En su rollout sobre un server type AMD Bergamo, al pasar de Linux 6.4+CFS a 6.9+`sched_ext`, reportó 28 % menos p99 en la etapa de retrieval; esa cifra pertenece a ese workload, hardware, comparación y metodología, no a `sched_ext` en general.

**[SPEC]/[CODE]** `sched_ext`, incorporado upstream desde Linux 6.12, permite políticas de scheduling BPF cargadas desde userspace y posee mecanismos de salida/fallback si la política falla. La oportunidad es iterar sin mantener un fork completo del kernel; el costo es convertir la política y sus hints de aplicación en parte crítica del sistema.

**[IMPL] Protocolo de adopción:**

1. demostrar con perfiles que run-queue delay, preemption, migraciones o locality dominan el tail;
2. clasificar threads/requests con una señal estable, autenticada dentro del proceso y de cardinalidad acotada;
3. definir fairness, starvation bound, affinity/NUMA, prioridad de background y comportamiento ante clasificación ausente;
4. comparar contra el scheduler general bajo igual kernel, mitigations, CPU set, IRQ placement, carga y energía;
5. medir offered load, goodput, p50/p99/p99.9/max, run-queue delay, migrations, LLC/DRAM, CPU/energía y progreso de background;
6. canary por hosts, watchdog/fallback y unload/rollback probado; una policy caída no puede dejar CPUs o trabajo inaccesibles;
7. repetir con overload, burst, mixed tenants, cgroup throttling, CPU hotplug y actualización de kernel.

No introducir scheduling personalizado si batching, pools, locks, allocation, I/O o backpressure explican el problema: aumenta superficie privilegiada y acopla aplicación, kernel y operación.

Fuentes: [caso de producción Meta](https://engineering.fb.com/2026/07/13/ml-applications/modernizing-the-meta-ads-service-with-an-open-source-kernel-scheduler/) y [documentación upstream de sched_ext](https://docs.kernel.org/scheduler/sched-ext.html).

## 22. Inventario curricular y fuentes pendientes

### MIT 6.172 — Performance Engineering of Software Systems

**Verificado:** 23 videoclases, 10 tareas, 4 proyectos, recitaciones, slides y quizzes públicos en MIT OpenCourseWare. Las evaluaciones se inventarían por su función; no se reproducen soluciones.

| Bloque | Lecciones | Resultado esperado |
|---|---:|---|
| introducción y reducción de trabajo | 1–3 | baseline, matrix multiply, Bentley rules y bit techniques |
| arquitectura y compilador | 4–5, 9 | source→assembly→hardware, vectorización y límites del compilador |
| paralelismo | 6–8, 13, 16–17, 20 | work/span, races, runtime, nondeterminismo, atomics/locks y especulación |
| medición y memoria | 10–12, 14–15 | timing, allocators, caché y cache-oblivious algorithms |
| optimización aplicada | 18–19, 21–23 | autotuning, codewalks, TSP, grafos y lenguajes dinámicos |

Prácticas verificadas:

- HW1–HW10: herramientas/C, profiling, vectorización, Cilk/reducers, teoría, allocator, análisis dinámico, cachegrind, ejecución determinista y sincronización.
- Proyectos: bit hacks, collision detection, allocator dinámico serial y Leiserchess.
- El proceso pedagógico exige beta correcta, tests, revisión de diseño/código, medición, informe de optimizaciones y resultado final.

### OSTEP

**Verificado:** versión 1.10; capítulos públicos organizados en virtualización, concurrencia, persistencia y seguridad, con code, homework y proyectos asociados. El inventario completo consta de 57 capítulos principales, además de diálogos/apéndices.

### Intel

**Verificado el 2026-08-21:** Intel SDM v092; Optimization Reference Manual volumen 1 y 2 v050. La ISA y los eventos PMU son específicos de familia/modelo; este manual no fijará cifras universales.

### Estado de cobertura

- [x] síntesis lección por lección de las 23 clases MIT 6.172 contra slides, inventario, tareas y proyectos oficiales;
- [x] mapa completo y labs de MIT 6.S081 (2021), con límites de xv6 explicitados;
- [x] mapa completo CS:APP 3e/CMU y función profesional de sus labs;
- [x] mapa completo OSTEP 1.10 y reconciliación conceptual con xv6/CS:APP;
- [x] extracción técnica selectiva de VM, scheduler, barriers, I/O y crash recovery contra Linux actual;
- [x] metodología de Brendan Gregg reconciliada con perf/tracepoints y seguridad de observabilidad de Linux;
- [x] marco ISA/microarquitectura/PMU/SIMD/topología reconciliado para Intel/AMD/Arm;
- [x] C++/LLVM/Rust atomics y memory-order proof template;
- [x] hazard pointers, epochs/QSBR y RCU reconciliados con proof obligations, failure modes y métricas;

### Extensiones deliberadamente condicionadas

- Revisar pasajes de video sólo cuando aporten una enseñanza no contenida en slides/materiales escritos o una atribución puntual necesaria.
- Extraer PMU events y tuning por familia/modelo sólo cuando el proyecto declare hardware target.
- Elegir e implementar reclamación concreta sólo con lenguaje, biblioteca, workload y progress guarantee definidos.
- La auditoría transversal con `SOFTWARE_BACKEND_API_ENGINEERING.md` quedó completada el 2026-08-21: SLO/goodput, parsing/serialización, pools, backpressure, retries, caché y observabilidad deben conservar semántica y medirse en el camino real.
- La auditoría transversal con `NETWORKING_DISTRIBUTED_STREAMING.md` quedó completada el 2026-08-21: NIC/IRQ/NAPI, queue ownership, CPU/cache/NUMA, syscalls/copias, batching, polling, timestamps, retransmisión y backpressure forman un único camino medible; una mejora local no puede degradar corrección ni p99.9 end-to-end.
- La auditoría transversal con `DATABASE_STORAGE_INTERNALS.md` quedó completada el 2026-08-21: page cache no es durability; buffer/query/LSM/allocator memory compite bajo cgroup/NUMA; fsync, WAL, checkpoint, vacuum/compaction y recovery pertenecen al camino medido. Comparar sólo con igual isolation/durability/replication y ejecutar crash/power-loss model, carga sostenida con maintenance y p99.9/goodput.
- La auditoría transversal con `GPU_ACCELERATED_COMPUTING.md` quedó completada el 2026-08-21: incluir CPU feeder, queue, launch, transfers, device, collective y readback en el camino real; fijar NUMA/topología, clocks/thermal y estados cold/warm; distinguir CUDA event de wall-clock. Streams, pinned buffers, allocator pools y callbacks tienen ownership/lifetime explícito, y una mejora de kernel no se acepta si empeora goodput, p99.9, consumo, precisión o recovery end-to-end.
- La auditoría transversal con `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md` quedó completada el 2026-08-21: profiler, debug socket, PMU, core dump y eBPF son superficies privilegiadas con acceso temporal/auditado y redaction; artifact/toolchain/config conservan provenance; hardening, microcode y mitigaciones de side channels se validan por threat model y se miden porque cambian performance. CPU/NUMA/pools/queues/cgroups tienen headroom y failure domains; overload usa admission, bounded queues, deadlines/cancelación y retry budget. Ninguna optimización puede ampliar TCB/privilegios, exponer datos o romper SLO, canary, rollback y recovery.
- La auditoría con microestructura quedó cerrada. Cruce con `FRONTEND_PRODUCT_ENGINEERING_UX.md` cerrado el 2026-08-21: la UI mide input delay, handler/render/presentation, memory growth, parse/compile/hydration y p75/tails en dispositivos reales; Web Worker, WASM, OffscreenCanvas o SIMD sólo se adoptan tras localizar CPU/copia/cola y medir transferencia/lifecycle. Ninguna optimización de main thread puede romper semántica, focus, accessibility tree, freshness, cancelación o recovery; FPS/promedio/bundle comprimido aislados no prueban task success ni latencia perceptible.
- Cruce con `ALGORITHMS_DATA_STRUCTURES_PROBLEM_SOLVING.md` cerrado el 2026-08-21: algoritmos decide corrección, trabajo asintótico y estructura; este manual decide cómo layout, caché, NUMA, SIMD, allocation, sincronización y runtime convierten ese trabajo en costo físico. Primero descartar complejidad inadecuada y después perfilar; comparar alternativas con idéntica semántica. Work/span, RAM, I/O y cache-oblivious son modelos declarados, no predicciones automáticas del host.
- Cruce con `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md` cerrado el 2026-08-21: arquitectura asigna quality budgets, boundaries y failure domains; este manual valida capacidad y costo físico del path real. Un diagrama de replicas/queues no prueba headroom, isolation ni p99.9, y una optimización local no autoriza violar contratos, deployment/recovery o reversibilidad.
- Cruce con `DATA_ENGINEERING_ANALYTICS.md` cerrado el 2026-08-21: file/row-group/column/page, Arrow buffers, partitions, shuffle, spill y compression se traducen a bytes, CPU, memoria, I/O y red medibles. Data engineering decide semántica/layout; sistemas perfila scan, serialization, locality, NUMA, page cache y tails. Benchmark con igual output/quality y workload; un codec o zero-copy nominal no prueba menor costo end-to-end.
- Cruce con `TOOLCHAINS_BUILDS_PACKAGING_FFI.md` cerrado el 2026-08-21: compiler/linker/flags/ISA/libc/runtime/LTO/PGO/sanitizers forman parte de artifact identity y benchmark manifest. Toolchains fija hermeticidad, host/target, ABI y symbols; sistemas mide machine code y runtime reales. Comparar mismo digest/config y no distribuir accidentalmente una ISA del builder; FFI conserva allocator, unwind, lifetime y memory model.
- Cruce con `NATIVE_MOBILE_DESKTOP_ENGINEERING.md` cerrado el 2026-08-21: startup, input-to-presentation, frames, memory pressure, background, battery, thermal, radio y storage se miden como journeys en release build y hardware real. Sistemas localiza CPU/GPU/I/O/locks/queues; nativo conserva lifecycle, platform policy y user outcome. Mover trabajo del main thread no autoriza colas ilimitadas, priority inversion, wakeups o resultados stale tras cerrar la pantalla.

## Fuentes primarias

- MIT 6.172: <https://ocw.mit.edu/courses/6-172-performance-engineering-of-software-systems-fall-2018/>
- MIT 6.172 — slides: <https://ocw.mit.edu/courses/6-172-performance-engineering-of-software-systems-fall-2018/pages/lecture-slides/>
- MIT 6.S081/6.828: <https://pdos.csail.mit.edu/6.828/2021/>
- CS:APP/CMU: <https://csapp.cs.cmu.edu/>
- OSTEP: <https://pages.cs.wisc.edu/~remzi/OSTEP/>
- Brendan Gregg: <https://www.brendangregg.com/systems-performance-2nd-edition-book.html>
- Brendan Gregg — metodologías: <https://www.brendangregg.com/methodology.html>
- Brendan Gregg — USE Method: <https://www.brendangregg.com/usemethod.html>
- Intel SDM: <https://www.intel.com/content/www/us/en/developer/articles/technical/intel-sdm.html>
- Intel Optimization: <https://www.intel.com/content/www/us/en/developer/articles/technical/intel64-and-ia32-architectures-optimization.html>
- AMD Developer Documentation (buscar publication `40332` y la revisión vigente): <https://www.amd.com/en/developer/browse-by-resource-type/documentation.html>
- Armv8-A memory model guide: <https://developer.arm.com/-/media/Arm%20Developer%20Community/PDF/Learn%20the%20Architecture/Armv8-A%20memory%20model%20guide.pdf?revision=58b1dd0a-3800-4218-b21a-f95a0332034c>
- Linux perf security: <https://www.kernel.org/doc/html/latest/admin-guide/perf-security.html>
- Linux tracepoints: <https://www.kernel.org/doc/html/latest/trace/tracepoints.html>
- Linux scheduler: <https://www.kernel.org/doc/html/latest/scheduler/index.html>
- Linux memory management: <https://www.kernel.org/doc/html/latest/mm/index.html>
- Linux memory barriers: <https://www.kernel.org/doc/html/latest/core-api/wrappers/memory-barriers.html>
- Linux source, snapshot auditado: <https://github.com/torvalds/linux/tree/26260251022fbc2f248a3d747a9b2b961b18d2d8>
- Linux `io_uring` UAPI y core fijados: <https://github.com/torvalds/linux/blob/26260251022fbc2f248a3d747a9b2b961b18d2d8/include/uapi/linux/io_uring.h>, <https://github.com/torvalds/linux/blob/26260251022fbc2f248a3d747a9b2b961b18d2d8/io_uring/io_uring.c>
- ext4/JBD2 journal: <https://www.kernel.org/doc/html/latest/filesystems/ext4/journal.html>
- C++ Core Guidelines: <https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines>
- LLVM Atomics: <https://llvm.org/docs/Atomics.html>
- Rust Book: <https://doc.rust-lang.org/stable/book/>
- Rustonomicon: <https://doc.rust-lang.org/stable/nomicon/>
- Rust atomics: <https://doc.rust-lang.org/std/sync/atomic/>
- Rust memory model status: <https://doc.rust-lang.org/reference/memory-model.html>
- Linux RCU Handbook: <https://www.kernel.org/doc/html/latest/RCU/index.html>
- WG21 P2530R3 — Hazard Pointers for C++26: <https://www.open-std.org/jtc1/sc22/wg21/docs/papers/2023/p2530r3.pdf>
- Maged M. Michael — Hazard Pointers (DOI): <https://doi.org/10.1109/TPDS.2004.8>

## Regla final

```text
Correctness is a gate.
The workload defines relevance.
The profiler localizes cost.
The specification defines semantics.
The benchmark decides performance only for its declared environment.
Production decides whether the trade-off is worth operating.
```
