# GPU y cómputo acelerado — manual operativo

> **Estado:** núcleo general v1 auditado dentro del alcance declarado el 2026-08-21. Corpus verificado: schedule, slides y prácticas públicas de CMU 15-418/15-618 Fall 2025; material público oficial seleccionado de Berkeley CS267; CUDA 13.3; Nsight Compute 13.3 y Nsight Systems; NCCL 2.31.2; CUTLASS 4.7.0; TensorRT 11.2.1; PyTorch 2.13; Triton; AMD HIP/ROCm 7.x; SYCL 2020 rev. 12 y papers primarios inventariados en §26. Cruces recíprocos cerrados con sistemas, backend, redes, datos, seguridad/SRE/cloud, microestructura, frontend y los tres manuales de IA relacionados.
> **Propósito:** decidir si acelerar, construir kernels correctos, localizar límites reales y operar entrenamiento/inferencia/multi-GPU sin confundir pico teórico con rendimiento útil.
> **Alcance:** principios transferibles y contratos vigentes; toda cifra, instrucción o capacidad se vuelve específica al declarar GPU, compute capability/ISA, driver, toolkit, framework y shape.

## 0. Cómo usar este manual

Sirve para:

- decidir CPU, GPU, librería, compiler/DSL o kernel propio;
- mapear un algoritmo a threads/warps/blocks/tiles y jerarquía de memoria;
- distinguir compute, bandwidth, latency, launch, transfer y synchronization bound;
- optimizar con evidencia de timeline, kernel profiler y resultados validados;
- construir entrenamiento e inferencia con precisión, memoria y comunicación controladas;
- entregar a Codex un contrato reproducible de implementación, benchmark y rollback.

Etiquetas:

- `[SPEC]`: especificación o guía oficial versionada;
- `[ACADEMIC]`: curso o paper primario;
- `[CODE]`: API/código/documentación de implementación;
- `[PROD]`: práctica operativa publicada y situada;
- `[MEASURED]`: evidencia reproducida en hardware/workload objetivo;
- `[IMPL]`: obligación derivada para esta biblioteca;
- `[OPEN]`: depende del proyecto o falta auditar.

Reglas:

1. GPU no acelera automáticamente trabajo pequeño, secuencial, irregular o dominado por transferencias.
2. Un kernel más rápido no implica una aplicación más rápida; medir extremo a extremo.
3. Occupancy alta no equivale a utilization ni performance alta.
4. Menor precisión no se acepta sin prueba numérica y de calidad del producto.
5. Una operación encolada no terminó; asincronía cambia errores, lifetime y medición.
6. Elegir primero librería/framework optimizado; escribir kernel cuando el profile y el contrato lo justifican.

## 1. Contrato del problema

Antes de tocar código:

```yaml
operation: ""
inputs:
  shapes: []
  dtype_layout_alignment: []
  size_distribution: ""
outputs_and_tolerance: ""
baseline_cpu_or_framework: ""
latency_slo: {p50: "", p99: "", p999: ""}
throughput_goodput: ""
batching_dynamic_shapes: ""
hardware_topology: ""
software_versions: ""
memory_budget: {host: "", device: "", workspace: ""}
transfer_path: ""
concurrency_multitenancy: ""
determinism_reproducibility: ""
failure_and_fallback: ""
power_cost_budget: ""
```

Preguntas de gate:

- ¿Qué porcentaje end-to-end es acelerable según profile, no intuición?
- ¿Hay suficiente trabajo paralelo para amortizar launch, setup y movement?
- ¿Dónde viven hoy y mañana los datos?
- ¿La salida requiere exactitud bitwise, tolerancia o métrica de calidad?
- ¿Shapes y control flow son estables o altamente dinámicos?
- ¿El SLO prioriza single-request latency, batch throughput o ambos?
- ¿Qué ocurre ante OOM, reset, sticky error, driver mismatch o GPU ausente?

## 2. Modelo mental de extremo a extremo

```text
input/storage/network
→ CPU parsing/preprocess
→ host allocation/pinning
→ H2D or shared/system memory
→ launch/queue
→ GPU kernels + device memory traffic
→ collectives/peer traffic
→ D2H/postprocess
→ response/persist
```

Tiempo aproximado:

```text
T_total = T_cpu + T_queue + T_transfer + T_launch
        + T_kernels + T_sync + T_communication + T_post
```

Overlap sólo reduce el camino crítico cuando existen recursos independientes, dependencias correctas y trabajo suficiente. Sumar tiempos de timeline no equivale siempre a wall time por concurrencia; medir ambos.

### 2.1 Amdahl, Gustafson y costo de offload

Si una fracción `P` se acelera `S` veces:

```text
speedup ≤ 1 / ((1 - P) + P/S)
```

El término real incluye movimiento, launch y sincronización. Weak scaling puede aumentar el problema para aprovechar más recursos, pero no responde al SLO de un input fijo. Declarar cuál se mide.

### 2.2 GPU como throughput machine

La GPU dedica muchos recursos a ejecutar y ocultar latencia con múltiples warps; la CPU privilegia control, cache y single-thread latency. No es una jerarquía de “mejor/peor”. El algoritmo, batch y residencia de datos deciden.

## 3. Arquitectura y modelo SIMT

CUDA organiza trabajo como grid → thread blocks → threads. Blocks se asignan a SMs; threads se ejecutan en grupos hardware llamados warps. La arquitectura exacta, warp scheduling, caches, Tensor Cores, shared memory y límites dependen de compute capability.

### 3.1 Obligaciones de mapeo

- Cada thread identifica un elemento/tile/fragmento sin salir de bounds.
- Un block debe ser ejecutable en cualquier SM y no depender del orden entre blocks.
- Sin mecanismo especial, no existe barrera global dentro de un kernel común.
- Un warp ejecuta instrucciones sobre lanes activos; branches divergentes serializan paths según arquitectura/modelo.
- Grid/block dimensions cubren distribución real, incluido tail incompleto.

Grid-stride loop general:

```cpp
for (std::size_t i = blockIdx.x * blockDim.x + threadIdx.x;
     i < n;
     i += static_cast<std::size_t>(blockDim.x) * gridDim.x) {
  output[i] = f(input[i]);
}
```

No copiar este patrón si cada elemento no es independiente o si `n`, stride o multiplicaciones pueden overflow.

### 3.2 Divergencia

Branches uniformes por warp cuestan distinto de branches con lanes tomando paths opuestos. Antes de branchless tricks:

- medir branch efficiency/stalls e instruction count;
- considerar agrupar/ordenar datos si preserva semántica;
- evaluar predication versus trabajo extra;
- no leer/escribir fuera de bounds bajo máscara incorrecta;
- recordar que divergencia también surge de loops de longitud variable.

### 3.3 Tiling

Un tile cambia tráfico global por reuse local:

```text
global tile load → shared/register reuse → compute → global store
```

Elegir tile desde shape, dtype, layout, bank mapping, registers, shared memory, Tensor Core fragment y tail handling. El tile óptimo no es universal entre GPU generations.

## 4. Jerarquía de memoria y movimiento

Modelo conceptual:

```text
registers/thread
→ shared memory/block y L1
→ L2/device
→ HBM/GDDR
→ NVLink/PCIe
→ host DRAM/storage/network
```

Latencia, bandwidth, capacidad, coherencia y scope difieren. “Memoria unificada” no elimina topology ni page migration.

### 4.1 Coalescing

Threads vecinos deberían acceder direcciones que el hardware pueda atender con pocas transacciones. Medir requested versus transferred bytes y sector utilization. Un layout AoS cómodo para CPU puede desperdiciar GPU bandwidth; SoA o tiles pueden mejorar acceso a costa de conversiones.

### 4.2 Shared memory

Shared memory permite reuse y cooperación dentro del block. Riesgos:

- bank conflicts;
- barrera ausente o condicional;
- ocupación reducida por bytes/block;
- tiles con padding incorrecto;
- leer datos no inicializados;
- asumir persistencia más allá del block/kernel.

### 4.3 Registers y local memory

Registers son rápidos pero finitos. Mayor register count/thread puede reducir blocks/SM; spilling coloca valores en local memory respaldada fuera del register file y puede añadir tráfico caro. No imponer un register limit sólo para subir occupancy sin medir spills e instrucciones.

### 4.4 Host-device transfers

- pageable host memory puede requerir staging;
- pinned/page-locked memory habilita transferencias eficientes/async bajo condiciones, pero reduce memoria pageable y debe limitarse;
- muchos transfers pequeños suelen perder contra batching;
- `cudaMemcpyAsync` sólo se solapa si memoria, stream, dependencias y copy engines lo permiten;
- peer access/NVLink/PCIe dependen de topology;
- GPUDirect/RDMA/Storage requiere hardware, driver, filesystem/NIC y configuración compatibles.

La mejor transferencia es la evitada: mantener datos residentes, fusionar etapas y mover resultados compactos.

### 4.5 Unified/System Memory

Unified Memory simplifica addressability, no garantiza residencia óptima. Page faults, migration, oversubscription, first-touch, prefetch y access hints cambian tails. En producción medir fault/migration bytes y topology; no juzgar sólo un segundo run ya migrado.

## 5. Diseño de kernels

### 5.1 Orden de decisión

1. referencia correcta simple;
2. librería/framework existente;
3. layout y algoritmo;
4. fusion/batching;
5. mapping y memory traffic;
6. occupancy/latency hiding;
7. instruction-level tuning;
8. assembly-specific sólo con hardware congelado.

### 5.2 Patrones

| Patrón | Riesgo dominante |
|---|---|
| map/elementwise | launch y bandwidth; fusion |
| reduction | associativity, synchronization, tail y precision |
| scan | work efficiency, bank conflicts, multi-block composition |
| histogram | contention/skew; privatization y merge |
| transpose | coalescing, shared tile y bank conflicts |
| stencil | halo, reuse y boundary conditions |
| GEMM/conv | tiling jerárquico, data movement, precision y epilogue fusion |
| gather/scatter | irregular access, duplicates y atomics |
| graph/sparse | load imbalance, locality y frontier/degree skew |

### 5.3 Fusion

Fusion puede eliminar intermediate writes y launches, pero aumenta register pressure, code size, compile time y recomputation. Una fusión demasiado grande reduce occupancy o impide elegir el mejor kernel por operador. Comparar end-to-end, no sólo bytes teóricos.

### 5.4 Atomics

Atomicity aplica a operación, address y scope especificados; no crea un protocolo completo ni ordena automáticamente otros datos. Medir contention y distribución. Alternativas: privatización por thread/warp/block, sort/reduce, sharding, hierarchical aggregation o un único owner.

**[SPEC]** En el memory model de CUDA 13.3, una operación sólo sincroniza threads incluidos por su scope. Mezclar un `store` block-scope en un block con un `load` device-scope en otro no amplía el alcance del primero: existe data race y el comportamiento es indefinido. Para GPU↔GPU o CPU↔GPU, verificar además soporte de atomics peer/system mediante los atributos documentados; un puntero accesible no implica atomicidad nativa.

## 6. Sincronización y memory model

Separar:

- barrier: quién debe llegar;
- visibility/order: qué writes se observan;
- atomicity: qué actualización es indivisible;
- scope: thread/block/device/system;
- lifetime: cuándo puede reutilizarse memoria.

`__syncthreads()` coordina threads de un block cuando todos alcanzan la barrera conforme al contrato. Colocarla en branch divergente donde no todos participan puede deadlock o ser indefinido. Warp-synchronous assumptions deben usar primitives y masks documentados; no depender de lockstep histórico implícito.

**[SPEC]** CUDA extiende C++ con scopes `thread`, `block`, `device` y `system`; el costo y el alcance crecen con la distancia. El orden correcto exige una relación *happens-before* y operaciones acquire/release/fences/atomics cuyo scope incluya a todos los participantes. Barrier, fence y atomic resuelven dimensiones distintas: no sustituir uno por otro.

Streams y events coordinan comandos, no protegen estructuras host ni validan ownership. Un buffer no puede reciclarse hasta que todos los streams que lo usan hayan completado; frameworks con caching allocators pueden requerir registrar el stream/lifetime.

## 7. Asincronía, streams, events y graphs

### 7.1 Asincronía

Launch y muchas APIs retornan antes de completar trabajo. Consecuencias:

- error puede aparecer en sincronización/API posterior;
- temporización host sin sync mide enqueue;
- input/output lifetime se extiende;
- destructores/cancelación de request no detienen automáticamente GPU;
- una sync accidental destruye overlap y cambia tails.

En desarrollo usar `cudaGetLastError()` tras launch y sincronización controlada donde se necesite localizar fallos; producción requiere política para errores async/sticky.

`cudaPeekAtLastError()` observa sin limpiar y `cudaGetLastError()` devuelve y limpia el último error del contexto. Ninguno transforma trabajo encolado en completado: los fallos de ejecución pueden emerger recién al esperar un event/stream. Asociar el fallo a una frontera conocida y no continuar usando un contexto cuya validez ya no puede garantizarse.

### 7.2 Streams y events

Dentro de un stream el orden es lineal; entre streams sólo hay orden con dependencia explícita o semántica documentada. Events permiten dependencia y timing device-side. Diseñar DAG:

```text
H2D(stream A) → event ready → kernel(stream B)
→ event done → D2H(stream C)
```

Probar carreras con inputs compartidos, allocator reuse y callbacks. Más streams no crean más engines/SMs ni garantizan concurrencia si un kernel ocupa todos los recursos.

**[SPEC]** La prioridad de un stream es una sugerencia, no una garantía de orden ni de preemption. Para copias CPU↔GPU realmente asíncronas, el buffer host debe ser page-locked; con memoria pageable la API puede funcionar correctamente pero perder el overlap. Una host function encolada con `cudaLaunchHostFunc` no debe llamar APIs CUDA.

### 7.3 CUDA Graphs

Graphs amortizan CPU/driver launch overhead para secuencias repetidas. Requieren capturabilidad, lifetimes y shapes/control suficientemente estables. En PyTorch, replay reutiliza direcciones virtuales y omite CPU work; violar restricciones puede causar error o resultado silenciosamente incorrecto. Medir capture/instantiate/update y memoria además de replay.

Cambios de topología o tipo de nodo requieren re-instanciar; parámetros compatibles pueden actualizarse sobre el graph executable. Una actualización rige lanzamientos futuros, no los ya en vuelo. En PyTorch, compartir pools privados entre graphs sólo es seguro si se reproducen en el orden capturado y nunca concurrentemente; un replay puede sobrescribir outputs de otro graph que comparte pool.

## 8. Occupancy, latency hiding y recursos

Occupancy es warps activos respecto del máximo; la limitan threads/block, registers, shared memory, blocks/SM y arquitectura. Puede ayudar a ocultar latency, pero:

- un kernel compute-bound puede rendir con occupancy moderada;
- reducir registers puede causar spills;
- blocks demasiado pequeños desperdician scheduling;
- blocks enormes reducen flexibilidad/residencia;
- dependencia serial larga no se arregla sólo con más warps.

Medir achieved occupancy, eligible/active warps, issue rate y razones de stall; correlacionar con bandwidth/compute utilization y duración.

## 9. Roofline y límites

Arithmetic intensity:

```text
AI = useful operations / bytes moved from chosen memory level
```

Roofline aproxima:

```text
attainable_perf ≤ min(peak_compute, AI × attainable_bandwidth)
```

Usar pico sostenible y el nivel de memoria relevante; contar bytes/ops útiles con semántica consistente. Un kernel puede ser L1-bound, L2-bound, HBM-bound, instruction/latency-bound o launch-bound sin encajar en una única línea global.

Diagnóstico:

| Síntoma | Hipótesis a probar |
|---|---|
| bajo GPU active y muchos gaps | CPU/launch/input/locks |
| HBM cerca del sostenible | reducir bytes, mejorar locality/fusion/compression |
| compute pipe alto | algoritmo, precision, Tensor Cores, instructions |
| baja utilization y stalls dependency | más ILP/warps o algoritmo |
| alta replay/transactions | coalescing/alignment/bank/conflicts |
| copies dominan | residency, batching, pinned/peer/topology |

## 10. Correctitud numérica y precisión

Floating point no es real exacto. Parallel reduction cambia orden y rounding; FMA, flush-to-zero, fast math, Tensor Cores y distintas libraries pueden cambiar bits.

**[SPEC]** CUDA 13.3 documenta que suma y multiplicación de precisión finita no son asociativas. FMA realiza producto+suma con un solo redondeo y normalmente mejora precisión, pero cambia el resultado frente a dos operaciones; `--use_fast_math` también habilita intrinsics más rápidos y menos exactos y `-ftz=true` elimina subnormales. Cada flag forma parte del contrato numérico y del artifact.

### 10.1 Contrato numérico

Definir por output:

- dtype de input, accumulation y output;
- rango, escala y distribución;
- error absoluto/relativo/ULP apropiado;
- NaN/Inf/subnormal/signed zero;
- tolerancia por elemento y métrica global;
- determinism requerido;
- impacto en loss, accuracy, recall o decisión de negocio.

### 10.2 Mixed precision

FP16, BF16, TF32, FP8, FP4/INT8 y variantes tienen rangos/precisión distintos. Tensor Core availability y acumulación dependen de hardware/API. AMP selecciona precisiones; gradient scaling ayuda con underflow FP16, no prueba convergencia.

Gate:

1. baseline FP32/FP64 según dominio;
2. comparar intermediate outputs y gradients;
3. detectar overflow/underflow/NaN/Inf;
4. evaluar entrenamiento completo o corpus representativo;
5. probar seeds y casos extremos;
6. documentar layers/ops forzados a mayor precisión;
7. conservar rollback.

### 10.3 Determinismo y reproducibilidad

Seed no basta. Algoritmos, atomics, reduction order, autotuning, driver, library, compiler, hardware y scheduling afectan. PyTorch advierte que reproducibilidad completa no se garantiza entre releases/plataformas ni CPU/GPU. Separar:

- bitwise repeatability en ambiente fijo;
- statistical reproducibility;
- model quality equivalence;
- performance determinism/tail stability.

## 11. Profiling disciplinado

### 11.1 Orden correcto

```text
wall-clock + service metrics
→ system timeline
→ kernels dominantes
→ kernel counters/source
→ hipótesis única
→ cambio
→ correctness + remeasure
```

**[SPEC]** Nsight Systems responde cuándo CPU, runtime, copies, kernels, NVTX, OS y collectives se solapan o esperan. Nsight Compute inspecciona kernels y puede replay/recoger múltiples métricas, alterando ejecución. No empezar con cientos de counters antes de encontrar el camino crítico.

La utilization temporal de Nsight Systems responde si hubo alguna actividad, no qué fracción de recursos ocupó: una memcpy pequeña y un kernel que satura la GPU cuentan como intervalos activos. GPU Metrics muestrea SM/Tensor/IO, pero no atribuye por sí solo cada muestra a proceso/contexto. Delimitar la ventana con NVTX/capture ranges y comprobar overhead y missing samples.

### 11.2 Timeline

Buscar:

- gaps entre kernels y CPU launch overhead;
- implicit synchronizations;
- H2D/D2H y ausencia de overlap;
- stream dependencies y serialization inesperada;
- dataloader/preprocess starvation;
- NCCL overlap o exposed communication;
- first-run initialization/JIT/autotune;
- allocator/OOM/retry;
- frecuencia, power y thermal throttling.

### 11.3 Kernel profile

Registrar GPU/SM exacto y analizar:

- duration y launch geometry;
- registers/thread, shared/block y occupancy limiters;
- SM/Tensor/ALU utilization;
- DRAM/L2/L1 throughput y hit behavior;
- global load/store efficiency;
- branch/warp state y stall reasons;
- instruction mix, spills y source correlation.

Un contador es evidencia dentro del modelo del profiler; no causa por sí solo. Cambiar una cosa y confirmar impacto end-to-end.

**[SPEC]** Nsight Compute puede serializar launches y ejecutar múltiples replay passes; kernel replay guarda/restaura memoria escrita, mientras application replay reejecuta toda la aplicación y exige determinismo de kernels, devices, contexts, streams y rangos. No inferir wall time desde timers host/events ejecutados bajo el profiler. Recoger primero un set pequeño y filtrar kernels/rangos reduce perturbación.

## 12. Benchmarking reproducible

### 12.1 Manifiesto

```yaml
gpu: {model: "", count: 0, clocks_power_mode: "", topology: ""}
host: {cpu_numa: "", memory: "", pcie: ""}
software: {driver: "", toolkit: "", libraries: "", framework: "", build: ""}
workload: {shapes: [], dtypes: [], layouts: [], batch: "", distribution: ""}
state: {warmup: "", cache: "", jit_autotune: "", concurrency: ""}
semantics: {precision: "", determinism: "", output_check: ""}
measurement: {wall: "", events: "", repetitions: 0, duration: ""}
```

### 12.2 Timing

- Warmup incluye context, module load, JIT, autotune y memory pools según objetivo.
- CUDA events miden intervalo device; wall-clock mide servicio completo.
- Sincronizar sólo donde define frontera de medida.
- Reportar distribución, no mejor sample.
- Separar first/cold, steady-state, shape change y concurrent load.
- Evitar dead-code/elided result; validar output dentro o fuera según contrato, pero declararlo.

### 12.3 Resultado

Reportar latency p50/p95/p99/p99.9/max, throughput y goodput, CPU/GPU utilization, memory peak/reserved, transfers, energy/power si importa, precision/quality y error/timeout/OOM. Comparar con idéntica semántica.

## 13. Librería, compiler/DSL o kernel propio

Orden preferido:

| Nivel | Usar cuando | Evidencia |
|---|---|---|
| framework op | cubre semántica y shapes | profile + correctness |
| cuBLAS/cuDNN/cuFFT/oneDNN/rocBLAS | primitive estable | benchmark de algoritmos/workspace |
| graph/compiler | fusion y specialization útiles | graph breaks, compile/recompile, quality |
| Triton/CUTLASS | operación/fusión especializada | tests diferenciales + sweep |
| CUDA/HIP/SYCL propio | control requerido y hotspot probado | review memory/sync/numeric + profiler |
| PTX/ISA | último recurso, target fijo | versión/arch guards + fallback |

No reimplementar GEMM/attention/sort/reduction sólo porque el ejemplo pequeño parece simple. Libraries contienen dispatch y kernels por arquitectura/shape.

## 14. CUDA profesional

### 14.1 Versiones y compatibilidad

Congelar:

```text
GPU architecture/compute capability
+ driver
+ CUDA runtime/toolkit
+ PTX/SASS targets
+ cuBLAS/cuDNN/NCCL/CUTLASS/framework
+ compiler/host ABI/container
```

Fat binaries/JIT/fallback tienen tiempos y compatibilidad distintos. No compilar sólo para la GPU del desarrollador. Leer release notes, deprecations y known issues.

### 14.2 Error handling

- comprobar return status de cada API;
- comprobar launch error y completion donde corresponda;
- adjuntar operation/shape/device/stream sin datos sensibles;
- diferenciar recoverable allocation/input error de context/device sticky failure;
- no seguir sirviendo resultados después de error que invalida contexto;
- tener health, drain/restart y CPU/otro-device fallback cuando el producto lo exige.

### 14.3 Allocation

`cudaMalloc/free` frecuente puede sincronizar y fragmentar. Pools y stream-ordered allocator amortizan costo pero requieren lifetime/dependency correctos. Medir allocated, reserved, fragmentation, high watermark y headroom para workspace/collectives/JIT.

**[SPEC]** `cudaMallocAsync` y `cudaFreeAsync` se ordenan en un stream. Si otro stream usa el allocation, events/dependencies deben ordenar allocation → usos → free; acceder después de iniciado el free es comportamiento indefinido. El release threshold del pool intercambia footprint por menos llamadas al SO. Distinguir bytes *used* de *reserved* y registrar ambos high-water marks.

## 15. CUTLASS y Triton

### 15.1 CUTLASS

**[CODE]** CUTLASS 4.7.0 ofrece abstracciones C++ y DSLs Python/CuTe para GEMM y operaciones relacionadas, con jerarquía explícita de layouts, copies y MMA/Tensor Cores. Es potente cuando epilogue/fusion/layout no está cubierto por una librería de alto nivel.

CuTe DSL seguía en public beta en la documentación 4.7.0 y Primitives/Task Scheduling figuraban como experimentales. Además, las features architecture-accelerated con sufijo `a` no tienen la forward compatibility general del PTX: por ejemplo, `sm_90a`/`sm_100a` deben coincidir con el target soportado. Aislar estas APIs y conservar kernel/library fallback.

Obligaciones:

- target architecture y toolkit exactos;
- alignment/layout/stride y tail;
- accumulator/output precision;
- workspace y schedule;
- referencia numérica y unit tests oficiales;
- sweep representativo, no un único tamaño múltiplo perfecto;
- beta/experimental API aislada detrás de adapter.

### 15.2 Triton

Triton expresa programas por bloques con Python y compiler. Sus tutoriales oficiales progresan por vector add, fused softmax, matmul, dropout, layer norm, attention, grouped/persistent/block-scaled matmul.

Gate de kernel Triton:

- mask correcta en todos los loads/stores tails;
- strides y layouts no contiguos probados;
- launch grid y meta-parameters por shape;
- numerical reference y gradcheck cuando aplica;
- autotune keys incluyen toda dimensión que cambia óptimo/correctitud;
- cache/JIT/compile latency medida;
- resource/spill/profile y fallback.

DSL no elimina hardware; hace más productivo expresar tiling y specialization.

El intérprete Triton sirve para inspección funcional secuencial en CPU, pero no reproduce scheduling ni performance GPU y tiene limitaciones —por ejemplo `bfloat16` e indirect addressing—. Usar `static_assert/device_assert`, tests diferenciales y Compute Sanitizer/ASan del backend; autotuning debe ocurrir fuera del hot path o con cache/artifact controlado.

## 16. Portabilidad: HIP/ROCm y SYCL

### 16.1 HIP/ROCm

HIP ofrece un modelo C++ familiar para AMD/NVIDIA según toolchain, pero source similarity no implica idéntica wave size, memory hierarchy, library behavior, performance counters o tuning. HIP API 7.0 introdujo cambios incompatibles que pueden exigir recompilar; usar documentación HIP/ROCm 7.x y profiler del target y no traducir warp assumptions mecánicamente.

Porting gate:

- compile y tests en ambos targets reales;
- subgroup/warp size y masks;
- atomics/memory scope;
- shared/LDS bank/layout;
- math intrinsics y precision;
- libraries y algorithm availability;
- topology/collectives;
- profiler y performance por target.

### 16.2 SYCL

**[SPEC]** SYCL 2020 rev. 12 es un modelo C++ single-source y cross-platform con queues, events, buffers/accessors o USM, groups/subgroups y backend interop. Portabilidad funcional usa el feature set común; performance portable puede requerir specialization y extensiones controladas.

Declarar device selection, required aspects, subgroup assumptions, USM lifetime, queue ordering/error handler y backend-specific code. Ejecutar conformance/functional/performance matrix, no asumir que compilar equivale a soportar.

## 17. Multi-GPU y collectives

### 17.1 Topology primero

Inventariar GPU↔GPU, GPU↔NIC, PCIe switches/root complexes, NVLink/NVSwitch/XGMI, NUMA CPU/memory y container visibility. Logical device IDs no describen distancia.

### 17.2 NCCL

**[CODE]** NCCL 2.31.2 ofrece collectives topology-aware; no es framework distribuido completo. AllReduce, AllGather, ReduceScatter, Broadcast y point-to-point mueven/combinen datos con contratos diferentes.

Gate:

- versión driver/CUDA/NCCL/framework y NIC plugin;
- rank/device mapping único;
- count/dtype/order iguales entre participantes;
- communicator init/error/abort y timeout;
- stream dependency antes y después de collective;
- topology y transport observados;
- payload sweep y overlap real;
- fallo de rank/link y recovery del job.

Una discrepancia en orden/count puede hang o corrupción. Async error no se arregla ignorando watchdog. Aislar logs diagnósticos porque variables debug intensas cambian volumen/timing.

La devolución de una collective sólo confirma que fue encolada, no completada; usar semántica CUDA para observarla. `ncclGroupEnd()` normal confirma enqueue, pero con communicator nonblocking puede devolver `ncclInProgress`: no usar el stream hasta que `ncclCommGetAsyncError()` indique `ncclSuccess`. Tras un error del communicator no se puede asumir completion ni correctitud de operaciones previas: coordinar revoke/abort/restart según API y framework. Mezclar múltiples streams dentro de un mismo group introduce un punto global de sincronización entre esos streams.

### 17.3 Paralelismo de ML

- data parallel replica modelo y reduce gradients;
- tensor parallel particiona operaciones/tensores y comunica frecuentemente;
- pipeline parallel particiona layers y crea bubbles/schedules;
- sequence/context/expert parallel resuelven dimensiones distintas;
- sharded optimizer/parameters reduce memoria a costa de comunicación/complexity.

Elegir desde memory footprint, compute/communication ratio, topology, batch/microbatch y failure recovery. Más GPUs pueden reducir eficiencia o empeorar single-request latency.

## 18. Entrenamiento de deep learning

Presupuesto por step:

```text
input + parameters + gradients + optimizer state
+ activations/saved tensors + temporary workspace
+ communication buffers + allocator fragmentation
```

### 18.1 Throughput correcto

Medir samples/tokens útiles por segundo, step time distribution, data wait, forward/backward/optimizer/collective, memory peak, overflow/skipped steps y model quality. No contar padding, failed/retried work o samples descartados como goodput.

### 18.2 Memoria

- gradient accumulation cambia effective batch y optimizer cadence;
- activation checkpointing cambia memoria por recomputation;
- sharding/offload mueve memoria a comunicación/CPU/storage;
- mixed precision reduce ciertos tensors, no todos;
- allocator reserved no equivale a live tensors;
- sequence length suele escalar activaciones/attention no linealmente según algoritmo.

OOM gate: worst shape, first step, optimizer step, validation, checkpoint save, collective peak y coexistencia con otros procesos.

### 18.3 Pipeline de datos

GPU idle puede venir de decode/tokenize/augment/storage/network. Medir queue depth, host CPU, pinned memory, worker lifecycle, NUMA y H2D. Prefetch sin límite transforma starvation en OOM/latency.

### 18.4 Convergencia

Optimización de kernel no se acepta sólo por loss de pocas iteraciones. Comparar curvas, final metrics, gradient statistics y estabilidad sobre seeds/datasets adecuados. Cualquier cambio de precision, reduction o fused optimizer necesita gate del manual de Deep Learning.

### 18.5 Semántica PyTorch 2.13

**[CODE]** Las operaciones CUDA son asíncronas; timing correcto requiere `torch.cuda.synchronize()` o events. En streams no-default, el usuario debe ordenar productores/consumidores con `wait_stream`/events y extender el lifetime del tensor con `record_stream` cuando el caching allocator podría reciclar su memoria. Esto aplica incluso si el nuevo kernel sobrescribe el buffer sin leerlo: puede existir trabajo pendiente del propietario anterior.

Cada op CUDA del backward se ejecuta en el stream de su forward correspondiente. Consumir gradients desde otro stream exige una dependencia explícita; el default stream ya no aporta la sincronización histórica de PyTorch 1.9 y anteriores. Separar `memory_allocated` —tensors vivos— de `memory_reserved` —pool del allocator—; `empty_cache()` no libera tensors activos ni aumenta la memoria utilizable por el propio proceso para esos tensors.

## 19. Inferencia

### 19.1 Objetivos separados

- offline throughput;
- online batch throughput;
- single-request latency;
- time to first token;
- time per output token/inter-token latency;
- deadline goodput y costo por resultado.

Dynamic batching intercambia wait por utilization. Definir max wait, batch/shape policy, fairness, cancellation y memory. Medir por clase de request, no promedio global.

### 19.2 TensorRT

**[CODE]** TensorRT 11.2.1 compila modelos a engines específicos con tactics, fusion, precision y optimization profiles. El engine serializado cruza una trust boundary y ejecuta native code: sólo cargar artifacts propios/firmados.

Obligaciones:

- ONNX/framework reference y intermediate comparison;
- hardware/software/build flags exactos;
- profiles para shapes reales;
- workspace/runtime memory acotada;
- tactic/timing cache compatible con target;
- warmup y first inference tras profile/shape change;
- un execution context por thread según guía;
- sticky-error isolation/restart;
- engine rebuild/rollback.

TensorRT 11.2 usa redes strongly typed por defecto y eliminó weak typing/implicit calibration heredados; plugins nuevos deben usar V3. Objetos runtime/engine pueden compartirse sólo para operaciones no mutantes; un `IExecutionContext` no es thread-safe y el patrón canónico es uno por worker/thread. Auxiliary streams pueden reducir latency intra-inference, pero consumen más memoria y TensorRT inserta dependencias con el mainstream; concurrencia de contexts cambia los recursos disponibles respecto del build y puede invalidar la tactic óptima. Construir y medir bajo un nivel de contención representativo.

### 19.3 LLM inference

Separar prefill compute-heavy de decode pequeño/memory/launch-sensitive. Presupuestar weights, KV cache, temporary workspace y fragmentation. Continuous batching, paged KV, speculative decoding, quantization, tensor/pipeline parallelism y prefix caching cambian calidad, fairness, latency y memoria; medir corpus y concurrencia real.

No aplicar GPU a lógica de order book “atómica” por reflejo: batches pequeños, branching y host/NIC round trips suelen favorecer CPU. GPU puede ser útil para analytics, feature batches o matrices residentes, siempre fuera del path de corrección secuencial salvo evidencia.

## 20. Baja latencia

Para p99.9:

- mantener context/model/buffers calientes cuando el costo lo permite;
- evitar JIT/autotune/allocation/profile switch en request path;
- predefinir shapes/buckets y memory pools;
- reducir CPU launch con fusion/graphs sólo si graph-safe;
- controlar batching wait y head-of-line blocking;
- aislar tenants/workloads o usar mecanismos hardware soportados tras medir;
- observar clocks, P-state, thermal/power, ECC y background activity;
- incluir queue y transfers, no sólo CUDA event del kernel.

Polling/busy waiting en CPU puede reducir wake latency pero consume core/energía y compite con feeders/NCCL. Integrar con el manual de sistemas y medir end-to-end.

## 21. Testing, debugging y confiabilidad

### 21.1 Capas

1. reference implementation;
2. unit tests de shapes/dtypes/layout/tails;
3. randomized/property/differential tests;
4. race/memory sanitizers cuando soportados;
5. numerical/adversarial extremes;
6. multi-stream/multi-thread stress;
7. long soak y memory leak/fragmentation;
8. device reset/OOM/worker death;
9. performance regression con correctness gate.

### 21.2 Casos mínimos

- `n=0`, 1, warp−1/warp/warp+1, block boundaries y tamaños grandes;
- non-contiguous strides, misalignment permitido y aliasing;
- NaN/Inf/subnormal/max/min/zero;
- duplicate indices y contention;
- dynamic shapes y profile changes;
- stream interleavings y early buffer reuse;
- device/rank absent, link degraded y timeout;
- compile cache missing/corrupt y cold start.

### 21.3 Diagnóstico

Herramientas sync/blocking son para localizar; pueden ocultar races o cambiar performance. Reducir a kernel/input mínimo, guardar versions y seed, comparar referencia, usar sanitizer/profiler apropiado y quitar instrumentation antes de medir release.

### 21.4 Segunda pasada de código NVIDIA

**[PROD]** Snapshots: CUTLASS `7107b05535f8977f5ecb9d01ee203205b1fd9bc4` (BSD-3-Clause), NCCL `7b83616df3ae082a1f32bb74c27458bfe8153a13` (licencias Apache-2.0/BSD declaradas por el repo) y TensorRT `10d15ae2f3b21437caf8248ba317cf76e0c0ebcb` (Apache-2.0 y notices de componentes). Complementan las versiones documentales de este manual; no sustituyen la matrix exacta de driver/toolkit/GPU ni describen componentes NVIDIA no publicados.

**CUTLASS.** El adapter separa `can_implement`, cálculo de workspace, initialize/update y launch; exige workspace no nulo cuando el problema lo necesita y comprueba el error inmediato del launch. El caller aún debe esperar el stream/event para conocer fallos de ejecución y conservar arguments/workspace/buffers hasta completion. Gate de kernel: capability/layout/alignment/shape → `can_implement`; workspace exacto + headroom; launch check; completion check; referencia diferencial y sweep de tails/dtypes/strides, no sólo shape óptima.

**NCCL.** Group mantiene profundidad y error por thread, agrupa jobs async y rechaza mezclar communicators blocking/nonblocking en el mismo group. Init/finalize/abort publican async result y abort flags; destroy coordina streams, callbacks, proxy y ranks locales. Gate multi-rank: todos los ranks invocan mismo collective/count/order; polling de async error con deadline; ante error, no reutilizar communicator/output incierto; abort/revoke coordinado y crear estado nuevo. Ensayar rank que muere durante init/enqueue/kernel/finalize, red degradada y un rank que omite o desordena una llamada.

**TensorRT.** Un engine puede tener varios execution contexts; el contexto posee profile/shapes/addresses y device memory con lifetime hasta que `enqueueV3` completa. Reutilizar esa memoria en otro contexto paralelo es comportamiento indefinido. Profile async y enqueue streams requieren dependencia explícita; auxiliary streams no deben ser default ni repetirse, y el runtime inserta events con el mainstream. Gate serving: un context por ejecución concurrente/worker, profile compatible, direcciones y workspace válidos, warmup por shape/profile, completion antes de recycle, OOM/error recorder/sticky-error policy, y benchmark bajo la misma concurrency/aux streams con que se construyó y eligió tactics.

Regla común de las tres implementaciones:

```text
validación de argumentos ≠ launch aceptado ≠ operación completada
≠ resultado numéricamente válido ≠ servicio publicado con SLO
```

## 22. Seguridad y multitenancy

- GPU memory puede contener prompts, embeddings, weights y datos regulados; zeroization/isolation depende de runtime/proveedor.
- No deserializar engines, cubins, PTX, models o plugins no confiables.
- Custom kernels/plugins son native code y amplían supply-chain/memory-safety risk.
- Limitar input shapes para evitar OOM/compile bomb/DoS.
- Compartir GPU crea interferencia de cache, bandwidth, memory y scheduler; quota de VRAM no garantiza latency isolation.
- Telemetría/NVTX no debe incluir secrets ni payloads.
- Containers necesitan driver/device permissions mínimos y versiones compatibles.

Cruce cerrado con `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md`: workload identity y device access son mínimos/efímeros; image/model/engine/plugin/kernel exigen digest+provenance+policy y parsing/sandbox según trust; Kubernetes RBAC/PSS/admission/device-plugin/node boundary se prueban en runtime. MIG/MPS/time-slicing y quota de VRAM declaran qué aíslan y qué comparten; ningún modo implica por sí solo confidencialidad, performance isolation o borrado de memoria. Confidential computing verifica attestation chain/measurement/freshness antes de entregar claves y conserva side-channel, patch, observability y recovery assumptions. OOM/ECC/Xid/reset/driver o communicator dudoso aíslan el failure domain, invalidan outputs inciertos y activan fallback/rollback e incidente.

## 23. Despliegue y operación

Artifact manifest:

```yaml
source_revision: ""
model_kernel_hash: ""
compiler_flags: []
target_architectures: []
driver_toolkit_runtime: ""
framework_libraries: {}
precision_profiles: []
reference_quality: ""
benchmark_report: ""
fallback_rollback: ""
```

Canary por GPU/shape/tenant, comparar output y SLO, observar JIT/rebuild/cache/OOM/ECC/reset, y revertir artifact+config juntos. Un engine/kernel válido en un SKU no se promueve a otro sin compatibility y benchmark.

Capacity planning incluye concurrency útil, VRAM headroom, workspace peak, collectives, replicas para fallo/maintenance, thermal/power y host/NIC/storage feeders. GPU utilization alta puede coexistir con goodput bajo.

## 24. Gate profesional

Antes de producción:

1. baseline y output contract;
2. GPU/topology/software manifest;
3. timeline end-to-end y bottleneck demostrado;
4. memory/transfer/launch/compute/communication budget;
5. kernels/libraries con tests de tails, strides y numérica;
6. asincronía, streams y lifetime probados;
7. precision y model quality aprobadas;
8. first-run, dynamic shape y steady-state medidos;
9. OOM/reset/sticky error/fallback ensayados;
10. multi-GPU failure y communicator recovery probados;
11. p99.9/goodput bajo concurrencia y background work;
12. artifact trust, rollout, observabilidad y rollback.

## 25. Contrato para Codex

Prompt mínimo:

```text
Inspecciona primero el baseline, shapes/dtypes/layouts, tests, hardware/topology,
versiones y profile. No optimices por folklore. Conserva una referencia correcta,
explicita precision/determinism, asincronía y lifetimes. Prefiere librerías maduras;
si escribes kernel, prueba tails/strides/extremos y compara diferencialmente.
Mide wall-clock y device time con warmup y sincronización correctos; reporta
goodput/percentiles/memoria, evidencia del profiler, riesgos y fallback.
```

Codex debe devolver:

1. bottleneck e hipótesis con evidencia;
2. decisión CPU/library/compiler/custom;
3. invariantes de mapping, memory, sync y numeric;
4. cambio mínimo;
5. tests correctness/race/numerical;
6. benchmark manifest y resultados;
7. compatibility matrix;
8. deployment, monitoring y rollback;
9. `[OPEN]` específicos del hardware/proyecto.

## 26. Fuentes primarias inventariadas

### 26.1 Academia

- CMU 15-418/15-618 Fall 2025 — home, schedule, slides y assignments: <https://www.cs.cmu.edu/~418/>, <https://www.cs.cmu.edu/~418/schedule.html>, <https://www.cs.cmu.edu/~418/assignments.html>
- CMU Assignment 2 público — CUDA con árboles `saxpy`, `scan` y `render`: <https://github.com/cmu15418f25/asst2>
- Berkeley CS267 — landing oficial y material público histórico de Jim Demmel: <https://inst.eecs.berkeley.edu/~cs267/>, <https://people.eecs.berkeley.edu/~demmel/cs267/>

Matriz de cobertura curricular:

| Fuente | Bloques verificados | Incorporación operativa |
|---|---|---|
| CMU 15-418/618, clases 1–5 | motivación, ILP, multicore, modelos paralelos, CUDA/GPU | §§2–4, Amdahl, jerarquía, SIMT y mapping |
| CMU, clases 6–8 + A2 | fundamentos paralelos, scheduling, locality/communication/contention; SAXPY, scan y renderer | §§5, 8–12; ruta práctica §28 |
| CMU, clases 9–17 | coherence, interconnect, VM, consistency, synchronization, lock-free | §§4, 6, 17 y cruces con sistemas/redes |
| CMU, clases 18–20 | heterogeneidad, specialization, DSLs y graph DSL | §§13, 15–16 |
| CMU, clases 21–23 | DL básico, data parallel, model/pipeline parallel | §§17–19 |
| Berkeley CS267 público | memory hierarchy, locality, performance, modelos/máquinas, communication avoidance, álgebra y aplicaciones | §§2, 4, 9, 13, 17; no se presenta como auditoría del offering actual restringido |

Límite declarado: el schedule y materiales enlazados de CMU fueron recorridos; algunos videos requieren cuenta `andrew.cmu.edu`, por lo que no se atribuyen afirmaciones a videos no accesibles. El portal vigente de CS267 redirige parte de su archivo a CalNet; se usó únicamente material oficial públicamente accesible y no se afirma haber auditado una edición cerrada.

### 26.2 NVIDIA

- CUDA Toolkit 13.3: <https://docs.nvidia.com/cuda/>
- CUDA Programming Guide: <https://docs.nvidia.com/cuda/cuda-programming-guide/contents.html>
- CUDA Best Practices: <https://docs.nvidia.com/cuda/cuda-c-best-practices-guide/>
- Nsight Systems: <https://docs.nvidia.com/nsight-systems/>
- Nsight Compute 13.3: <https://docs.nvidia.com/nsight-compute/>
- NCCL 2.31.2: <https://docs.nvidia.com/deeplearning/nccl/>
- CUTLASS 4.7.0: <https://docs.nvidia.com/cutlass/latest/>
- TensorRT 11.2.1: <https://docs.nvidia.com/deeplearning/tensorrt/latest/>
- CUTLASS fijado `7107b05535f8977f5ecb9d01ee203205b1fd9bc4`: [GEMM universal adapter](https://github.com/NVIDIA/cutlass/blob/7107b05535f8977f5ecb9d01ee203205b1fd9bc4/include/cutlass/gemm/device/gemm_universal_adapter.h) y [unit-test corpus](https://github.com/NVIDIA/cutlass/tree/7107b05535f8977f5ecb9d01ee203205b1fd9bc4/test/unit/gemm/device).
- NCCL fijado `7b83616df3ae082a1f32bb74c27458bfe8153a13`: [group/async jobs](https://github.com/NVIDIA/nccl/blob/7b83616df3ae082a1f32bb74c27458bfe8153a13/src/group.cc), [enqueue](https://github.com/NVIDIA/nccl/blob/7b83616df3ae082a1f32bb74c27458bfe8153a13/src/enqueue/enqueue.cc) e [init/finalize/abort](https://github.com/NVIDIA/nccl/blob/7b83616df3ae082a1f32bb74c27458bfe8153a13/src/init.cc).
- TensorRT fijado `10d15ae2f3b21437caf8248ba317cf76e0c0ebcb`: [runtime interfaces](https://github.com/NVIDIA/TensorRT/blob/10d15ae2f3b21437caf8248ba317cf76e0c0ebcb/include/NvInferRuntime.h), [sample inference lifecycle](https://github.com/NVIDIA/TensorRT/blob/10d15ae2f3b21437caf8248ba317cf76e0c0ebcb/samples/common/sampleInference.cpp) y [engine build/load](https://github.com/NVIDIA/TensorRT/blob/10d15ae2f3b21437caf8248ba317cf76e0c0ebcb/samples/common/sampleEngines.cpp).

### 26.3 Lenguajes y portabilidad

- Triton: <https://triton-lang.org/main/>
- AMD HIP/ROCm 7.x: <https://rocm.docs.amd.com/projects/HIP/en/latest/>
- ROCm profilers: <https://rocm.docs.amd.com/projects/rocprofiler/en/latest/>, <https://rocm.docs.amd.com/projects/rocprofiler-compute/en/latest/>
- SYCL 2020 rev. 12: <https://registry.khronos.org/SYCL/specs/sycl-2020/html/sycl-2020.html>
- oneAPI 1.3 rev. 1: <https://oneapi-spec.uxlfoundation.org/specifications/oneapi/v1.3-rev-1/>

### 26.4 Framework

- PyTorch 2.13 CUDA semantics: <https://docs.pytorch.org/docs/2.13/notes/cuda.html>
- PyTorch 2.13 reproducibility: <https://docs.pytorch.org/docs/2.13/notes/randomness.html>
- PyTorch AMP: <https://docs.pytorch.org/docs/2.13/notes/amp_examples.html>
- PyTorch performance tuning: <https://docs.pytorch.org/tutorials/recipes/recipes/tuning_guide.html>

### 26.5 Papers fundacionales y de sistemas

- Nickolls et al., *Scalable Parallel Programming with CUDA* (2008): <https://doi.org/10.1145/1365490.1365500>
- Williams, Waterman y Patterson, *Roofline* (2009): <https://escholarship.org/uc/item/5tz795vq>
- Tillet, Kung y Cox, *Triton* (2019): <https://doi.org/10.1145/3315508.3329973>
- Rajbhandari et al., *ZeRO* (2019): <https://arxiv.org/abs/1910.02054>
- Huang et al., *GPipe* (2019): <https://arxiv.org/abs/1811.06965>
- Narayanan et al., *PipeDream* (SOSP 2019): <https://www.microsoft.com/en-us/research/publication/pipedream-generalized-pipeline-parallelism-for-dnn-training/>
- Narayanan et al., *Efficient Large-Scale Language Model Training on GPU Clusters Using Megatron-LM* (2021): <https://arxiv.org/abs/2104.04473>
- Dao et al., *FlashAttention* (2022): <https://arxiv.org/abs/2205.14135>

Síntesis transferible: Roofline relaciona trabajo con tráfico, Triton eleva el tile a abstracción del compiler, FlashAttention demuestra que cambiar IO puede superar una optimización puramente FLOP, y ZeRO/GPipe/PipeDream/Megatron muestran que memoria, comunicación, scheduling y semántica de actualización deben diseñarse conjuntamente. Sus cifras experimentales no son promesas para otro hardware.

## 27. Auditoría de trazabilidad

| Área | Autoridad contrastada | Contrato que queda en el manual | Estado |
|---|---|---|---|
| SIMT/memoria/kernels | CMU + CUDA 13.3 | mapping, coalescing, banks, registers, atomics y tails | cerrado |
| ejecución/orden | CUDA memory model, streams, graphs y allocator | scopes, happens-before, lifetime, capture/update y async free | cerrado |
| numérica | CUDA floating point + PyTorch reproducibility/AMP | dtype de acumulación, FMA/fast-math, determinismo y quality gate | cerrado |
| medición | Best Practices, Nsight Systems y Compute | APOD, timeline→kernel, replay/overhead y wall-clock separado | cerrado |
| librerías/DSL | CUTLASS 4.7, Triton | target, API maturity, autotune, debug y fallback | cerrado |
| portabilidad | HIP/ROCm 7.x, SYCL rev. 12, oneAPI | conformance + benchmark por backend; no performance portability asumida | cerrado |
| multi-GPU | CMU, CUDA multi-GPU, NCCL 2.31.2 y papers | topology, collective completion/error y estrategia híbrida | cerrado |
| framework/training | PyTorch 2.13 + ZeRO/GPipe/PipeDream/Megatron | streams/autograd, allocator, memoria y convergence gate | cerrado |
| inferencia | TensorRT 11.2.1 + FlashAttention + manual DLAI | engine trust, contexts, shapes, contention, TTFT/TPOT/goodput | cerrado |
| operación/baja latencia | sistemas + backend + redes + datos | end-to-end, queues, isolation, failure y rollback | cerrado dentro de los manuales existentes |
| seguridad/SRE/cloud | Stanford CS155 + NIST/OWASP + Google SRE + Kubernetes | TCB/attestation, artifact trust, IAM, tenant boundary, SLO/incident/recovery | cruce recíproco cerrado |

La auditoría técnica de este archivo y sus cruces están cerrados dentro del corpus declarado. Microestructura ya existe en `MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md`: la GPU queda fuera del mutation/risk path secuencial salvo prueba formal y medición; sí puede acelerar analytics, simulación, calibración y batches históricos con lineage/freshness explícitos.

## 28. Ruta práctica y suites de aceptación

La cadena mínima reproduce la disciplina de CMU sin copiar evaluaciones privadas:

1. **SAXPY/map:** baseline CPU y framework; kernel simple; validar tails; calcular effective bandwidth; separar transfer/launch/kernel.
2. **Scan/reduction:** versión de referencia; jerarquía warp/block/grid; probar no-potencias-de-dos, orden numérico y work efficiency.
3. **Renderer/irregular:** medir divergencia, acceso, contention y balance; comparar binning/tiling/privatización sin alterar blending/ordering.
4. **GEMM/fused op:** librería → Triton/CUTLASS sólo si existe gap; sweep de shapes/layouts/dtypes y roofline.
5. **Training step:** forward/backward/optimizer/collective; quality parity, memoria pico, tokens/s y fallo de rank.
6. **Serving:** cold/warm, shape profiles, dynamic batching, p99.9, OOM, context isolation, artifact trust y rollback.

Suites ejecutables exigidas:

| Suite | Casos | Gate |
|---|---|---|
| `gpu_correctness` | empty, 1, warp/block±1, strides, tails, aliasing, extremes | referencia + tolerancia explícita |
| `gpu_race_lifetime` | streams permutados, early free/reuse, events ausentes, contention | sanitizer limpio + DAG probado |
| `gpu_numeric` | dtypes, reduction orders, FMA/fast-math, NaN/Inf/subnormal | error/quality dentro de presupuesto |
| `gpu_perf` | cold/warm, shape sweep, concurrency, clocks/thermal | goodput y percentiles reproducibles |
| `gpu_memory` | live/reserved/workspace, fragmentation, graph pools, OOM | high-water bajo presupuesto + fallback |
| `gpu_multi_rank` | rank/link loss, timeout, communicator error/restart | sin hang; resultado fallido no se publica |
| `gpu_deploy` | driver/toolkit/arch matrix, corrupt/untrusted artifact, canary | compatibilidad y rollback demostrados |
| `gpu_security_recovery` | unauthorized device/artifact, cross-tenant assumptions, attestation/identity expiry, OOM/ECC/Xid/reset y restore | acceso/provenance válidos; containment, output invalidation y recovery demostrados |

## 29. Auditorías transversales

| Manual relacionado | Autoridad de frontera | Obligación recíproca |
|---|---|---|
| `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md` | CPU/NUMA/OS/allocator/thermal y medición host | GPU incluye queue, feeder, topology, power y p99.9 end-to-end |
| `SOFTWARE_BACKEND_API_ENGINEERING.md` | admission, deadline, cancellation, idempotency y lifecycle de request | GPU define bounded batching, context ownership, cancellation/fallback y no publica output tras fallo |
| `NETWORKING_DISTRIBUTED_STREAMING.md` | NIC/network/backpressure/failure distribuido | GPU trata RDMA/collectives como red con orden, timeout, topology y recovery |
| `DATABASE_STORAGE_INTERNALS.md` | persistencia, WAL/checkpoints, pool/IO y crash recovery | GPU no confunde completion device con persistencia ni acelera un pipeline dominado por storage |
| `DEEP_LEARNING_ANDREW_NG.md` | error analysis, train/dev/test, optimización y calidad | toda precisión/kernel/parallelism conserva experimento y métrica de modelo |
| `PYTORCH_TRANSFORMERS_GENERATIVE_MODELS.md` | semántica de modelo, training, serving y eval | este manual posee kernel, stream, allocator, topology y profiler |
| `AI_INFERENCE_PERFORMANCE_HARDWARE.md` | scheduler/KV/cache/quantization/SLO de inferencia | este manual posee engine/device correctness y capacidad física |
| `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md` | IAM, provenance, Kubernetes, TCB/attestation, SLO, incident y restore | GPU prueba device/artifact/tenant boundary, invalidación de output y failure/rollback físico |
| `MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md` | sequence/book/matching/order/risk y venue semantics | GPU no reordena state; batch/async output conserva epoch+sequence y se descarta si llega stale |
| `FRONTEND_PRODUCT_ENGINEERING_UX.md` | browser lifecycle, render, interaction, accessibility y user outcome | WebGPU/canvas/media declaran adapter/device loss, async completion, memory budget, fallback semántico y no bloquean main thread |
| `NATIVE_MOBILE_DESKTOP_ENGINEERING.md` | app/scene lifecycle, frame budget, resource pressure, battery y thermal | GPU libera/pausa recursos correctamente, trata device loss y conserva fallback/user state |

Regla de arbitraje: el manual de dominio define corrección del producto; este manual define corrección y rendimiento del cómputo acelerado; sistemas/redes/backend/datos definen los límites externos. Una optimización sólo se acepta cuando las cinco capas concuerdan.

Cruce de mercados auditado recíprocamente el 2026-08-21. Benchmark obligatorio: comparar CPU y GPU end-to-end incluyendo transfer, batching wait, launch, synchronization, freshness y fallback; kernel time aislado no justifica adopción.

Cruce frontend cerrado el 2026-08-21. El retorno de una API WebGPU o enqueue no equivale a pixel presentado ni tarea completada; medir CPU encode, upload, queue, GPU, readback/composite y presentation. Device loss/OOM/backgrounding/resizing no puede dejar output silenciosamente stale. Toda visualización acelerada conserva alternativa textual/semántica y pruebas de teclado/foco cuando comunica datos o permite acciones.

Cruce con `ALGORITHMS_DATA_STRUCTURES_PROBLEM_SOLVING.md` cerrado el 2026-08-21. La selección comienza por corrección, trabajo/span, dependencia y volumen de datos; recién después se mapea a SIMT, memoria y collectives. Scan, reduction, sort, histogram, sparse traversal, DP, graph y tiling exigen invariante y oráculo CPU pequeño. Menor work asintótico no garantiza menor tiempo GPU; medir transfers, divergence, occupancy, bandwidth, sincronización, precisión y camino completo.

Cruce con `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md` cerrado el 2026-08-21. Arquitectura decide si el acelerador pertenece a batch, serving o hot path mediante quality/cost/failure scenarios; GPU decide device topology, memory, async completion y fallback viable. Mostrar feeder, queues, transfers, model/kernel, collectives y recovery en runtime/deployment views. “GPU-enabled” no prueba capacity, portability ni operability.

Cruce con `DATA_ENGINEERING_ANALYTICS.md` cerrado el 2026-08-21. ETL/analytics GPU conserva schema, nulls, order donde aplique, precision, determinism, lineage y output parity contra oráculo CPU pequeño. Medir decode→host/device transfer→kernel→shuffle/collective→write, incluyendo batches, pinned/device memory y fallback. Acelerar kernel no justifica convertir Arrow/Parquet tipos ni aggregates con semántica distinta.

Cruce con `TOOLCHAINS_BUILDS_PACKAGING_FFI.md` cerrado el 2026-08-21. Toolkit/compiler/driver API, GPU arch/code objects, host compiler, libc/C++ runtime, PTX/JIT policy y native package tags forman la matrix del artefacto. Build hermético no incluye driver/kernel del host de ejecución por magia: compatibility y smoke test se hacen en targets reales. Extensions exponen ownership/stream/device/error y no dejan exception/panic cruzar ABI.

Cruce con `NATIVE_MOBILE_DESKTOP_ENGINEERING.md` cerrado el 2026-08-21. Render, media e inferencia on-device incluyen encode/upload/queue/device/composite/presentation en el presupuesto, no sólo kernel time. Background/minimize/memory warning/device loss deben cancelar, drenar o invalidar trabajo y liberar recursos según lifecycle. Medir batería, thermal throttling y oldest supported device; el fallback conserva semántica y accesibilidad.
