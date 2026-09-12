# Redes, sistemas distribuidos y streaming — manual operativo

> **Estado:** núcleo general v1 auditado el 2026-08-21 dentro del alcance declarado: Stanford CS144 Fall 2025, MIT 6.5840 Spring 2026, RFC/IETF seleccionadas, Kafka 4.3, Flink 2.3 y documentación oficial Linux/DPDK/RDMA. Los cruces con sistemas/rendimiento, backend/API, bases de datos, GPU, seguridad/SRE/cloud, microestructura y frontend están cerrados; extensiones dependientes de versión, hardware, venue o manuales futuros se declaran al final.
>
> **Autoridad académica:** Stanford CS144 (Keith Winstein) y MIT 6.5840 (Robert Morris, Frans Kaashoek y equipo PDOS), complementados por los papers originales.
>
> **Autoridad normativa:** IETF/RFC Editor. **Autoridad de implementación:** documentación y código oficial de Apache Kafka, Apache Flink, Linux, DPDK y RDMA.
>
> **Uso:** diseñar, implementar, revisar, probar y operar comunicaciones, servicios distribuidos y pipelines de eventos. No reemplaza la especificación versionada de un protocolo, producto, kernel, NIC o exchange concreto.

## 0. Cómo usar este manual

Etiquetas de evidencia:

- **[SPEC]** requisito o semántica de estándar oficial.
- **[ACADEMIC]** modelo, algoritmo o resultado académico.
- **[CODE]** comportamiento comprobable en código/documentación oficial.
- **[PROD]** práctica de producción con restricciones explícitas.
- **[MEASURED]** afirmación válida sólo para el entorno medido.
- **[IMPL]** decisión propuesta para un sistema concreto.
- **[OPEN]** incógnita que debe resolverse antes de depender de ella.

Reglas:

1. No convertir una propiedad local en garantía end-to-end.
2. No llamar “ordenado”, “durable”, “exactly once” o “disponible” a algo sin nombrar alcance y fallos tolerados.
3. Separar seguridad (*nunca ocurre algo incorrecto*) de vivacidad (*eventualmente progresa bajo supuestos declarados*).
4. Tratar timeout como evidencia de incertidumbre, no de fracaso remoto.
5. Probar bajo pérdida, duplicación, reordenamiento, partición, pausa, restart y version skew.
6. Optimizar después de preservar contrato, corrección y evidencia reproducible.

Este documento enseña semántica y disciplina. Los números concretos —timeouts, buffers, batches, réplicas, particiones, núcleos— requieren medición y threat/failure model del proyecto.

## 1. Contrato mínimo de un sistema distribuido

Antes de elegir protocolo o producto, completar:

```yaml
system:
  purpose: ""
  clients_and_operators: []
  data_ownership: ""
  trust_boundaries: []

correctness:
  invariants: []
  consistency_model: ""
  ordering_scope: ""
  duplicate_policy: ""
  loss_policy: ""
  time_semantics: ""

failure_model:
  crash_stop: true
  crash_recovery: true
  packet_loss_reorder_duplicate: true
  network_partition: true
  clock_skew_and_jump: true
  disk_loss_or_corruption: false
  byzantine_or_adversarial: false

service:
  slo: ""
  offered_load: ""
  goodput_target: ""
  latency_percentiles: [p50, p95, p99, p99_9]
  recovery_objectives: {rto: "", rpo: ""}

operations:
  overload_policy: ""
  retry_owner: ""
  reconciliation: ""
  upgrade_and_rollback: ""
  observability: []
```

**[ACADEMIC]** Un sistema distribuido no ofrece una observación instantánea y perfecta de todos sus componentes. Mensajes tardan, nodos se detienen y relojes difieren. **[PROD]** Por eso el contrato debe decir qué respuesta sigue siendo correcta cuando no se sabe si una operación remota ocurrió.

## 2. Modelo de red de punta a punta

### 2.1 Capas útiles, no cajas mágicas

```text
contrato de aplicación y framing
        ↓
HTTP / WebSocket / protocolo propio / Kafka / gRPC
        ↓
TLS y transporte: TCP, QUIC o UDP
        ↓
IP, routing, MTU, ICMP
        ↓
Ethernet/Wi-Fi, ARP/ND, NIC y colas
        ↓
driver, kernel, sockets, scheduler, memoria
        ↓
runtime, parser, colas y lógica de negocio
```

Cada capa agrega estado, buffers, límites y modos de fallo. La latencia observada incluye como mínimo:

```text
serialización + enqueue local + kernel/NIC + propagación + switching/routing
+ cola remota + parsing + scheduling + aplicación + respuesta
```

**[MEASURED]** Un RTT de `ping`, una conexión loopback o un benchmark del handler no demuestra latencia de negocio.

### 2.2 Datagrama, stream y mensaje

- IP ofrece datagramas; no promete entrega, orden ni ausencia de duplicados.
- UDP conserva límites de datagrama, pero no agrega confiabilidad, orden ni protección contra duplicados.
- TCP ofrece a la aplicación un byte stream confiable y ordenado; no conserva límites de `send()`.
- QUIC ofrece conexiones seguras y múltiples byte streams; tampoco conserva fronteras de frames de aplicación dentro de un stream.
- WebSocket sí define mensajes y frames sobre su transporte, con fragmentación y control frames.

**Regla:** todo protocolo sobre byte stream necesita framing explícito y límites antes de reservar memoria.

Opciones comunes:

```text
longitud fija | length-prefix | delimiter escapado | TLV | schema binario
```

El parser debe tolerar lecturas parciales, múltiples mensajes por lectura, EOF a mitad de frame y longitudes hostiles.

### 2.3 Bandwidth-delay product y colas

```text
BDP = bandwidth × RTT
```

El BDP aproxima los bytes necesarios en vuelo para llenar un camino. Buffers menores pueden limitar throughput; buffers enormes pueden ocultar congestión y elevar latencia (*bufferbloat*). Medir ocupación y edad de cola, no sólo tamaño configurado.

**[PROD]** Una cola sin límite transforma overload transitorio en agotamiento de memoria y latencia sin cota. Una cola acotada obliga a decidir: bloquear, rechazar, degradar, samplear o descartar según semántica.

## 3. IP, routing, MTU y DNS

### 3.1 IP y forwarding

Un router examina destino, selecciona ruta por prefijo más largo, decrementa hop limit/TTL y reenvía por la interfaz elegida. Routing calcula o distribuye rutas; forwarding aplica la tabla por paquete.

Invariantes:

- nunca asumir que request y response siguen el mismo camino;
- diseñar para cambios de ruta, ECMP y NAT;
- no identificar una sesión únicamente por IP origen;
- observar pérdida y latencia por path/zone/provider, no sólo global.

En el enlace local, ARP (IPv4) o Neighbor Discovery (IPv6) resuelve next hop a dirección de enlace. La caché expira y una resolución pendiente no autoriza una cola infinita: limitar datagramas/bytes/edad, espaciar solicitudes y descartar o devolver error según el protocolo. No confundir destino IP final con next hop de la interfaz.

### 3.2 MTU y fragmentación

La MTU útil cambia por enlaces, túneles y encapsulación. Path MTU Discovery puede fallar si ICMP necesario es bloqueado. Datagram PLPMTUD (RFC 8899) hace probes de tamaño desde el propio transporte/protocolo y exige recuperación ante black holes; no convierte un tamaño supuesto en universal. Evitar depender de fragmentación IP para tráfico crítico; limitar datagramas y comprobar tamaños reales después de headers y cifrado.

**[IMPL]** Registrar tamaño serializado p50/p99/máximo y errores por MTU. Para UDP/QUIC, seguir la RFC y la biblioteca concreta; no inventar un “MTU seguro universal”.

### 3.3 DNS es un sistema distribuido cacheado

**[SPEC]** DNS delega autoridad y permite cachear resource records según TTL. TTL no obliga a una aplicación o conexión existente a cambiar inmediatamente de destino.

Diseño:

- definir conducta ante múltiples A/AAAA y cambios de respuesta;
- considerar negative caching y fallos de resolver;
- no usar DNS como único mecanismo de failover con RTO estricto;
- observar resolución separada de connect/TLS/application;
- validar DNS y destino en controles SSRF después de redirects o re-resolución;
- tratar TTL bajo como costo/carga y no como consistencia instantánea.

## 4. UDP: mínimo mecanismo, máxima responsabilidad

**[SPEC]** RFC 768 define un servicio de datagramas con mecanismo mínimo; entrega y protección contra duplicados no están garantizadas. RFC 8085 reúne pautas de uso actuales.

UDP es apropiado cuando la aplicación puede:

- tolerar pérdida o reparar selectivamente;
- manejar duplicación y reordenamiento;
- delimitar mensajes;
- aplicar congestion control o limitar tasa de manera responsable;
- autenticar y proteger contra spoofing/amplificación;
- controlar fragmentación y tamaño.

**[SPEC]** El congestion control debe limitar el tráfico agregado que una aplicación dirige al mismo path, incluso si usa varios sockets o workers; repartirlo no elimina la obligación. Usar checksum según IP/protocolo y threat model. El keepalive puede conservar estado en middleboxes, pero no reemplaza ACK, detección de pérdida ni recuperación de sesión.

Contrato UDP mínimo:

```text
version | message_type | stream/session_id | sequence | timestamp_domain
payload_length | integrity/authentication | payload
```

El protocolo debe definir:

- ventana de secuencia y wraparound;
- gap detection y política de recovery;
- deduplicación y caducidad;
- heartbeat/liveness;
- replay protection;
- límites por peer y amplificación;
- conducta ante mensaje desconocido o truncado.

**[PROD]** Un feed UDP multicast de mercado puede priorizar oportunidad sobre reparación inmediata; la corrección se recupera con sequence numbers, snapshot/retransmission y reconciliación. La especificación del venue, no este manual, define el algoritmo exacto.

## 5. TCP: byte stream confiable, no transacción

### 5.1 Semántica esencial

**[SPEC]** RFC 9293 reemplaza RFC 793 como especificación base. TCP ofrece un byte stream bidireccional, confiable y en orden. Usa sequence numbers, ACK, checksum, retransmission, flow control y puertos.

No promete:

- fronteras de mensaje;
- que un `write()` sea una unidad de envío o un `read()` la recupere;
- que datos aceptados por el socket hayan sido procesados por la aplicación remota;
- timeout de negocio;
- liveness inherente;
- exactly-once de efectos externos.

### 5.2 Estado y cierre

Entender al menos: `LISTEN`, `SYN-SENT`, `SYN-RECEIVED`, `ESTABLISHED`, half-close, `FIN-WAIT`, `CLOSE-WAIT`, `LAST-ACK`, `TIME-WAIT` y reset.

Errores frecuentes:

- no leer hasta EOF y confundir half-close con caída;
- olvidar cerrar el descriptor que deja `CLOSE-WAIT`;
- asumir que `RST` conserva datos pendientes;
- reconectar sin reconciliar una operación cuyo resultado fue ambiguo;
- usar keepalive del SO como sustituto de heartbeat/deadline de aplicación.

TCP mantiene al menos tres espacios que no deben confundirse: sequence number de 32 bits con wrap e ISN, absolute sequence conceptual y stream index de bytes de aplicación. SYN y FIN consumen espacio de secuencia aunque no sean bytes del stream. Comparar valores wrapping requiere una referencia/checkpoint cercano; una comparación entera ingenua falla alrededor del wrap.

### 5.3 Flow control versus congestion control

- **Flow control:** protege al receptor; ventana anunciada limita bytes que acepta.
- **Congestion control:** protege el camino compartido; ventana de congestión responde a señales del network path.

El emisor efectivo queda limitado por ambas. RFC 5681 especifica slow start, congestion avoidance, fast retransmit y fast recovery; RFC 6298 especifica el cálculo base de retransmission timeout.

El receptor anuncia el siguiente byte requerido y una ventana de aceptación. El reassembler debe deduplicar/combinar solapamientos y contar capacidad real una sola vez; bytes fuera de la capacidad se descartan. Ante ventana cero, TCP usa probing regulado para descubrir que volvió a abrirse, no envío ilimitado.

**Regla:** no desactivar o combatir congestion control para obtener una cifra de laboratorio. Seleccionar algoritmo y pacing con evidencia del entorno y equidad requerida.

### 5.4 Head-of-line y conexiones

Pérdida de un segmento TCP retrasa la entrega posterior en ese byte stream aunque otras respuestas lógicas estén listas. Multiplexar muchas operaciones en una conexión reduce handshakes y sockets pero comparte backpressure y failure domain.

Medir:

- RTT y retransmissions;
- cwnd, bytes in flight y receive window;
- connect/TLS time;
- send/receive queue;
- resets, timeouts y reconnects;
- latencia por operación, no sólo por conexión.

### 5.5 Nagle, delayed ACK y mensajes pequeños

No aplicar `TCP_NODELAY` por reflejo. Nagle intenta reducir segmentos pequeños; delayed ACK puede interactuar con patrones request/response. **[MEASURED]** Comparar distribución de latencia, paquetes/segundo, bytes/paquete y CPU con el patrón real. Batching mejora eficiencia pero agrega espera: definir un máximo por bytes y tiempo.

## 6. QUIC y HTTP/3

**[SPEC]** QUIC v1 está definido por RFC 9000, seguridad con TLS por RFC 9001 y pérdida/congestión por RFC 9002. HTTP/3 se especifica en RFC 9114.

Propiedades útiles:

- transporte sobre UDP implementado normalmente en user space;
- TLS integrado;
- múltiples streams con flow control por stream y conexión;
- pérdida en un stream no bloquea la entrega de otros streams del mismo modo que un único byte stream TCP;
- connection IDs permiten ciertos cambios de path y migración;
- loss detection usa packet numbers y PTO.

Límites:

- cada stream continúa siendo un byte stream: requiere framing de aplicación;
- flow control de conexión aún puede acoplar streams;
- 0-RTT puede ser replayable y sólo sirve para operaciones seguras frente a replay;
- middleboxes, UDP blocking, CPU criptográfica y calidad de implementación importan;
- QUIC no elimina colas, congestión ni backpressure.

**[PROD]** Comparar HTTP/2 y HTTP/3 bajo pérdida, movilidad, CPU, handshake y tráfico reales; “más nuevo” no demuestra menor p99.

## 7. WebSocket y conexiones persistentes

**[SPEC]** RFC 6455 define handshake, frames, mensajes, masking cliente→servidor, fragmentation, ping/pong y cierre. RFC 8441 permite bootstrap mediante Extended CONNECT en HTTP/2 cuando fue negociado.

Contrato de aplicación requerido:

```yaml
protocol_version: ""
message_types: []
schema_and_limits: ""
authentication_refresh: ""
sequence_scope: ""
heartbeat: {ping_interval: "", dead_after: ""}
reconnect: {backoff: "", jitter: "", max: ""}
resume: {cursor_or_sequence: "", snapshot_rule: ""}
backpressure: {queue_bound: "", slow_consumer_policy: ""}
close_codes: {}
```

Reglas:

- autenticar handshake y autorizar cada suscripción/acción;
- limitar frame, mensaje, compresión, frecuencia y subscripciones;
- no depender de fronteras de fragmento como fronteras semánticas;
- aceptar que control frames pueden intercalarse dentro de un mensaje fragmentado; no son fragmentables y su payload está limitado por la RFC a 125 bytes;
- ping/pong sólo prueba la ruta hasta el peer, no que su estado de negocio esté actualizado;
- reconectar con jitter para evitar tormenta;
- recuperar estado mediante cursor/sequence/snapshot, no suponiendo continuidad;
- cerrar o degradar consumidores lentos antes de acumular memoria sin límite;
- implementar el closing handshake, timeout y códigos de cierre; un socket TCP cerrado no demuestra que el peer procesó el último mensaje;
- proteger compresión frente a bombs y filtraciones side-channel según threat model.

## 8. Sockets, event loops y backpressure

### 8.1 Partial I/O es normal

`send`/`recv` pueden transferir menos bytes que lo solicitado; nonblocking I/O puede indicar “intentar más tarde”. Un evento de readiness no garantiza que la siguiente operación complete todo.

State machine mínima por conexión:

```text
ACCEPT/CONNECT → TLS_HANDSHAKE → READ_HEADER → READ_BODY
→ DISPATCH → WRITE_QUEUE → HALF_CLOSE/CLOSE/ERROR
```

Cada transición necesita límites, timeout, cancelación y cleanup idempotente.

### 8.2 Modelos de concurrencia

| Modelo | Ventaja | Riesgo principal |
|---|---|---|
| thread por conexión | simple | stacks, scheduling, contention |
| pool + blocking I/O | control de recursos | pool exhaustion |
| event loop/readiness | muchas conexiones | bloquear loop, state machines complejas |
| completion I/O | batching y async kernel | lifetime/cancellation complejos |
| core/queue affinity | localidad y jitter menor | imbalance y operación rígida |

No elegir por moda. Medir conexiones, tasa de mensajes, CPU por mensaje, allocations, syscalls, context switches, migraciones de CPU y tail latency.

### 8.3 Backpressure end-to-end

```text
NIC RX → kernel receive queue → socket buffer → parser → work queue
→ dependency → output queue → socket buffer → NIC TX
```

Si una etapa no puede reducir demanda, propagar presión o rechazar, otra cola crece. Registrar para cada borde:

- capacidad y unidad;
- productor/consumidor;
- high/low watermark;
- política de overflow;
- edad máxima;
- señal de saturación;
- cómo se recupera o reconcilia lo descartado.

### 8.4 Connection pools

Un pool es una cola y un límite de concurrencia. Debe tener acquire deadline, máximo, health semantics y observabilidad. Pool demasiado grande transfiere saturación a la dependencia; demasiado pequeño impone cola local. Separar pools por failure domain o prioridad sólo cuando reduzca interferencia demostrable.

## 9. El problema fundamental: incertidumbre y fallo parcial

En un proceso local suele distinguirse retorno de función y crash. En red pueden ocurrir:

```text
request perdido
request procesado + response perdido
request en cola mientras cliente vence timeout
peer pausado, no muerto
partición asimétrica
respuesta tardía después de retry
```

Por eso:

- timeout no implica que la operación no ocurrió;
- retry puede duplicar efecto;
- failure detector es sospecha basada en tiempo;
- una respuesta tardía necesita request/operation ID y regla de obsolescencia;
- reconciliación es parte del protocolo, no soporte posterior.

**[ACADEMIC]** En Raft, la seguridad no depende de timing; la disponibilidad sí requiere relaciones entre broadcast time, election timeout y MTBF. Esta separación debe conservarse en diseños propios.

## 10. Tiempo, orden y causalidad

### 10.1 Relojes físicos

Wall clock puede saltar por sincronización o ajuste. Usar reloj monotónico para duraciones/deadlines dentro de un host. Para timestamps entre hosts declarar sincronización, incertidumbre y precisión requerida.

Nunca derivar orden causal perfecto de timestamps físicos sin supuestos adicionales.

### 10.2 Orden lógico

Lamport clocks preservan: si `a → b`, entonces `L(a) < L(b)`; la inversa no prueba causalidad. Vector clocks pueden detectar concurrencia a costo proporcional al conjunto representado. Sequence numbers dan orden sólo en su stream/partición/epoch declarado.

Toda secuencia operativa necesita:

- scope y productor autorizado;
- epoch/generation para reinicios y cambio de líder;
- regla de comparación y wrap;
- gap y duplicate behavior;
- persistencia o reconstrucción.

### 10.3 TrueTime no es “usar timestamps”

**[ACADEMIC]** Spanner expone un intervalo de incertidumbre temporal y espera cuando es necesario para cumplir external consistency. No copiar su conclusión sin su infraestructura de relojes y sus límites demostrados.

## 11. Consistencia: nombrar exactamente la propiedad

### 11.1 Modelos frecuentes

- **Linearizability:** cada operación parece ocurrir atómicamente entre invocación y respuesta, respetando precedencia real para el objeto.
- **Sequential consistency:** existe un orden secuencial compatible con el orden de cada proceso, no necesariamente con tiempo real.
- **Serializability:** transacciones equivalen a algún orden serial; no exige por sí sola tiempo real.
- **Strict serializability:** serializability más orden de tiempo real.
- **Causal consistency:** procesos observan causas antes que efectos.
- **Eventual consistency:** sin nuevas actualizaciones y bajo supuestos de comunicación, réplicas convergen; no define qué se lee durante convergencia.

No usar “strong consistency” sin traducirlo a una propiedad comprobable.

### 11.2 Safety, liveness y availability

- Safety se falsifica con una historia finita incorrecta.
- Liveness necesita supuestos de scheduling, red y tiempo para garantizar progreso.
- Durability pregunta qué sobrevive y bajo qué fallos.
- Availability debe incluir denominador, operación, región, timeout y respuestas válidas.

### 11.3 CAP sin eslogan

**[ACADEMIC]** Ante una partición de red, un sistema que exige una consistencia equivalente a registro atómico no puede garantizar respuesta exitosa de todos los nodos. CAP no dice que en operación normal sólo puedan elegirse “dos de tres”, ni decide automáticamente el modelo adecuado.

Registrar por operación:

| Operación | Durante partición | Consistencia | Reconciliación |
|---|---|---|---|
| lectura | rechazar / stale / quorum | modelo exacto | evidencia |
| escritura | rechazar / aceptar local | conflicto permitido | merge/owner |

## 12. Replicación y quórums

### 12.1 Replicar datos no replica automáticamente verdad

Definir:

- quién acepta writes;
- cómo se ordenan;
- cuándo se consideran committed;
- qué replica puede servir reads;
- cómo se detecta y repara divergence;
- cómo se incorpora un nodo con estado viejo;
- qué sucede si se pierde almacenamiento estable.

### 12.2 Primary/backup

Un primary ordena updates y los backups replican. Falta resolver elección única, fencing del primary viejo, commit, replay, snapshots y membership. Un heartbeat por sí solo no evita split brain.

**[ACADEMIC]** Chain Replication envía updates por la cadena y sirve queries desde el tail para combinar throughput con una garantía fuerte bajo su modelo fail-stop. Cambiar el punto de lectura puede cambiar la garantía.

### 12.3 Quórums

En un modelo simple con `N` réplicas, write quorum `W` y read quorum `R`, `R + W > N` produce intersección, pero no basta por sí solo para linearizability: importan versiones, coordinación, fallos, sloppy quorums, repairs y resolución de conflictos.

**[ACADEMIC]** Dynamo prioriza alta disponibilidad para workloads específicos y utiliza consistent hashing, object versioning, quorum-like techniques y reconciliación. No es permiso genérico para aceptar conflictos invisibles.

## 13. Consensus y replicated state machines

### 13.1 Qué resuelve consensus

Procesos proponen valores y acuerdan uno bajo un modelo de fallos. **[ACADEMIC]** Paxos exige que sólo se elija un valor propuesto, que no se elijan dos valores distintos y que nadie aprenda un valor no elegido.

Consensus no implementa por sí solo:

- schema o lógica de negocio;
- almacenamiento durable correcto;
- deduplicación de clientes;
- snapshots/compaction;
- reconfiguration segura;
- control de overload;
- reparación de corrupción;
- operaciones multi-shard.

### 13.2 Raft operativo

Raft separa leader election, log replication y safety; usa terms, randomized election timeouts y liderazgo fuerte.

Invariantes principales del paper:

- como máximo un líder por term;
- el líder sólo agrega a su log;
- logs con mismo index/term comparten el prefijo correspondiente;
- una entrada committed aparece en líderes futuros;
- dos state machines no aplican comandos distintos en el mismo índice.

Implementación profesional:

1. persistir term/vote/log en el orden requerido antes de responder;
2. aplicar sólo entradas committed y en orden;
3. deduplicar requests de cliente con session/request ID persistente;
4. usar read protocol válido; “leer del líder” no basta si no confirmó liderazgo actual;
5. snapshot y truncado deben preservar log matching y estado de sesión;
6. membership change requiere algoritmo conjunto/seguro, no edición simultánea informal;
7. limitar AppendEntries, snapshots y colas para que recovery no destruya servicio;
8. hacer fencing de procesos/leases/recursos externos.

### 13.3 Consensus no debe estar en el hot path sin necesidad

Consensus agrega coordinación y latencia. Mantener decisiones por shard/partition cuando la semántica lo permita; usar estado local derivable y caches donde un resultado stale sea explícitamente correcto. Nunca eliminar coordinación si rompe el invariante.

### 13.4 Fallos Byzantine y servidor no confiable

Crash-fault tolerance supone que un nodo deja de responder o reinicia según el modelo; no tolera comportamiento arbitrario. **[ACADEMIC]** PBFT estudia state machine replication cuando hasta `f` réplicas pueden ser Byzantine y requiere supuestos, autenticación y tamaños de quorum distintos de Raft/Paxos. No agregar “una réplica más” y declarar BFT.

**[ACADEMIC]** SUNDR muestra otra elección: con servidor de archivos no confiable, fork consistency hace que vistas incompatibles no puedan volver a unirse sin detección bajo sus supuestos criptográficos. Es una garantía distinta de availability o linearizability.

El threat model debe decidir explícitamente:

- crash, corrupción accidental o adversario activo;
- independencia real de fallos y diversidad de implementación;
- identidad/autenticidad de mensajes;
- replay/equivocation;
- trusted computing base y key management;
- propiedad de seguridad y supuesto de progreso.

### 13.5 De modelo a implementación verificada

**[ACADEMIC]** IronFleet combina refinement de state machines con verificación de código imperativo y separa pruebas de safety y liveness. La lección operativa no es “la verificación elimina todo bug”: especificación, trusted base, compilador, runtime, SO, hardware y supuestos de red permanecen en el argumento.

Para un protocolo crítico:

```text
spec abstracta → invariantes/model checking → protocolo → implementación
→ conformance/fuzz/fault injection → historia observable
```

Cada flecha puede contener un semantic gap; documentarlo y probarlo.

## 14. RPC, idempotencia y efectos ambiguos

Un RPC local-looking oculta red y ejecución remota. El cliente debe conocer:

```text
deadline + retry class + idempotency + request ID + error model + cancellation
```

Estados de resultado:

- definitivamente no ejecutado;
- ejecutado y confirmado;
- rechazado de manera definitiva;
- ambiguo: pudo ejecutarse.

Para operaciones mutantes:

```text
operation_id → mismo intento lógico → mismo resultado lógico
```

Persistir outcome o dato suficiente para reconstruirlo durante una ventana definida. Un retry con nuevo ID es una nueva operación. Cancelar espera local no revierte automáticamente el efecto remoto.

**[PROD]** Un único nivel debe ser dueño de retries por operación. Presupuesto total, backoff, jitter y límite de concurrencia evitan amplificación durante fallos.

## 15. Transacciones distribuidas, sagas y fencing

### 15.1 Atomic commit no es consensus completo

Two-phase commit coordina prepare/commit, pero puede bloquear si el coordinator falla y necesita logs/recovery. Consensus puede replicar decisiones, aunque el sistema completo aún debe integrar participants, storage y reconfiguration.

### 15.2 Sagas

Una saga divide trabajo en transacciones locales con compensaciones. Compensación no equivale a rollback perfecto: emails enviados, trades ejecutados o interacciones externas pueden ser irreversibles.

Por paso definir:

- idempotency key;
- precondition/version;
- resultado persistido;
- compensación y sus fallos;
- retry owner;
- estado terminal/manual;
- audit trail.

### 15.3 Leases y fencing tokens

Una lease expira según un modelo temporal; un cliente pausado puede despertar creyéndose owner. Entregar un fencing token monotónico y exigir que el recurso rechace tokens viejos. Si el recurso no valida fencing, la lease sola no protege el invariante.

## 16. Particionamiento, sharding y rebalancing

Una partition es simultáneamente unidad de orden, paralelismo, ownership, recuperación y a menudo disponibilidad. La clave de partición decide qué puede procesarse localmente y qué requiere coordinación.

Preguntas:

- ¿qué operaciones deben compartir orden?
- ¿distribución de keys es estable y no sesgada?
- ¿cómo se manejan hot keys?
- ¿qué cambia al aumentar particiones?
- ¿cómo se mueve estado sin doble owner ni pérdida?
- ¿cómo se versiona el mapping?
- ¿quién hace fencing durante rebalance?

**[PROD]** Rebalance es un evento de corrección, no sólo de capacidad. Pausar, drenar, checkpoint, revocar y reasignar en orden explícito.

### 16.1 Cómputo distribuido y data movement

**[ACADEMIC]** MapReduce vuelve explícitos map, shuffle/partition y reduce; reejecuta trabajo para tolerar fallos y explota data locality. Determinismo y commit atómico de outputs facilitan razonar sobre reexecution. Straggler mitigation consume capacidad extra y debe evitar duplicar efectos externos.

Distributed futures y actors, como en el paper de Ray estudiado por MIT, permiten expresar dependencias, ubicación y paralelismo; no eliminan ownership, garbage collection distribuida, recovery ni movimiento de datos. Pass-by-reference reduce copias cuando la ubicación ayuda, pero agrega metadata y coordinación. Elegir por tamaño y locality medidos.

El caso AWS Lambda del corpus ilustra otra regla: a gran escala, reducir data movement puede requerir cache, deduplicación y carga sparse/on-demand. Los resultados dependen de la distribución real —commonality, cacheability y sparsity— y no justifican lazy loading universal.

## 17. Logs y event streaming

### 17.1 Log distribuido

Un log particionado ofrece posiciones/offsets y orden total dentro de cada partición, no orden global entre particiones. Retention elimina historia por política; compaction conserva estado por key bajo reglas del producto, no todos los eventos.

Diseño de evento:

```yaml
event_id: "estable"
event_type: ""
schema_version: 1
producer: ""
aggregate_or_partition_key: ""
sequence_or_version: ""
event_time: ""
ingest_time: ""
causation_id: ""
correlation_id: ""
payload: {}
```

No mezclar:

- event time, ingest time y processing time;
- command e event;
- log de integración y audit log legal;
- tombstone y evento de negocio “deleted”.

### 17.2 Schema evolution

Definir compatibilidad backward/forward/full según productores y consumidores reales. Agregar un campo “opcional” puede ser incompatible si cambia semántica o un consumidor lo supone presente. Probar old producer→new consumer y new producer→old consumer con datos de borde.

## 18. Kafka: garantías y operación

> Versión observada al construir este núcleo: Apache Kafka 4.3 figura como release más reciente en el sitio oficial el 2026-08-21. Verificar versión exacta del cluster antes de aplicar configuraciones.

### 18.1 Modelo

- topic dividido en partitions;
- cada partition es un log ordenado con offsets;
- producer elige partition directa o por key/partitioner;
- en un consumer group, ownership distribuye partitions;
- replication mantiene copias y un líder atiende la partition;
- KRaft administra metadata mediante un quorum de controllers.

Garantía de orden: dentro de una topic-partition y bajo las condiciones del producer/protocolo. No existe orden total del topic completo.

### 18.2 Durabilidad

`acks=all` significa ACK de todas las réplicas actualmente in-sync, no de todo replica set asignado. Combinar con `min.insync.replicas` y evitar unclean leader election cuando la pérdida no sea aceptable. Declarar qué ocurre si todas las réplicas en sync se pierden.

Medir:

- under-replicated/offline partitions;
- ISR shrink/expand;
- controller/quorum health;
- produce/fetch request latency y errors;
- disk, page cache, network y replication lag;
- storage growth contra retention.

KRaft requiere mayoría viva del quorum de controllers. Tres controllers toleran uno caído y cinco toleran dos; separar roles broker/controller en despliegues críticos permite aislar carga, escalar y actualizar cada plano. Kafka soporta quorum dinámico desde la línea 4.1, pero migración, alta/baja de controllers y feature levels son operaciones explícitas: observar primero catch-up y ejecutar el procedimiento de la versión. No formatear o reemplazar metadata storage como reparación improvisada.

### 18.3 Delivery semantics

- **At most once:** commit de offset antes del efecto puede perder procesamiento.
- **At least once:** efecto antes del commit puede repetirse.
- **Kafka exactly-once:** producer idempotente + transactions + offsets y outputs Kafka en la misma transacción + consumer `read_committed`, bajo configuración válida.

**Regla crítica:** exactly-once Kafka no vuelve transaccional a un sink externo. Para DB/API/archivo se necesita cooperación transaccional, idempotencia/deduplicación o guardar offset junto al efecto.

### 18.4 Producer

Revisar como conjunto:

```text
acks | enable.idempotence | retries | max.in.flight
delivery.timeout | request.timeout | linger | batch.size | compression
transactional.id | partitioner | serialization
```

Batching y compresión intercambian CPU/espera por throughput y bytes. Una key equivocada genera skew o rompe orden requerido. No cambiar partitioner sin plan de compatibilidad.

En Kafka 4.3, `enable.idempotence` es `true` por defecto si no hay configuraciones incompatibles. Su validez exige `acks=all`, `retries>0` y `max.in.flight.requests.per.connection<=5`; si se desactiva idempotencia, retries con más de una request in-flight pueden reordenar batches. `delivery.timeout.ms` limita el tiempo total de éxito o fallo de un envío e incluye espera previa, requests y retries: debe ser coherente con `request.timeout.ms` y `linger.ms`.

### 18.5 Consumer

Controlar:

- poll/processing time versus session y rebalance;
- offset commit sólo después del efecto que representa;
- revoke/assign callbacks y cleanup;
- poison records y DLQ con auditabilidad;
- lag por partition, edad de evento y end-to-end latency;
- capacidad de replay y side effects;
- static/cooperative membership sólo con semántica comprendida.

`read_committed` oculta records abortados y sólo entrega hasta el last stable offset; una transacción abierta puede retrasar progreso visible. `max.poll.interval.ms` limita el tiempo entre llamadas a `poll()` y debe cubrir procesamiento real o separar poll de trabajo con ownership correcto.

Kafka 4.3 conserva `group.protocol=classic` como default del cliente; el protocolo `consumer`, GA desde 4.0, es incremental, elimina la barrera global de sincronización y mueve heartbeat, session timeout y assignors al broker. Migrarlo exige verificar APIs/configs que dejan de aplicar y probar upgrade/downgrade, churn, callbacks y efectos durante revoke/assign.

Lag cero no significa datos correctos; lag alto no identifica por sí solo CPU, skew, sink lento o rebalance.

### 18.6 Hot partitions y operación

Más partitions aumentan paralelismo pero también metadata, files, recovery y coordinación. La unidad de escala debe seguir el key distribution real. Probar fallo de broker, controller, disco lento, leader movement, rolling upgrade, quota y consumer churn.

### 18.7 Share Groups no son consumer groups con otro nombre

Desde Kafka 4.2, Share Groups son production-ready para trabajo por registro tipo cola: una partition puede abastecer a varios consumidores y la cantidad de consumidores puede superar la de partitions. Cada record adquirido queda bajo un lock temporal y se confirma individualmente como `ACCEPT`, `RELEASE`, `REJECT` o, cuando se habilita, `RENEW`; no hacer nada deja vencer el lock. El broker limita intentos de entrega y locks por partition.

**[SPEC]** Elegir Share Group sólo si no se necesita procesar la partition como stream ordenado por un dueño exclusivo. El lock no convierte el side effect externo en exactly-once. Dimensionar lock/renewal contra tiempo p99 de procesamiento y observar expiraciones, intentos, rejects, renewals y acknowledgements fallidos.

### 18.8 Código Kafka — idempotencia requiere estado y fencing

**[PROD]** Código upstream fijado en `e28414ee4a1dfa3d093dc10168f944f1a759434e` (Apache-2.0). El `TransactionManager` no implementa “retry mágico”: conserva producer ID/epoch, sequence por partition, requests in-flight, coordinator y una máquina con estados ready/in-transaction/committing/aborting/abortable/fatal. Una secuencia ambigua bloquea nuevos números hasta resolverla o renovar identidad; un estado fatal se envenena para impedir continuar y corromper garantías. En broker, `ProducerStateManager` persiste snapshots de epochs/sequences/transacciones, descarta snapshots corruptos y reconstruye desde snapshot + log dentro de offsets válidos.

Gate derivado: reiniciar producer/broker y cortar conexiones alrededor de append/ACK/commit/abort; duplicar y reordenar responses; truncar log; corromper/omitir snapshot; cambiar leader/epoch; dejar transacción abierta; expirar batch con requests in-flight. Verificar por partition: output, orden, duplicados, last stable/high watermark, fencing del productor viejo y convergencia después de replay. Si el efecto sale de Kafka, repetir la misma matriz contra el protocolo del sink; Kafka no puede reconstruir una transacción externa que éste no registra.

## 19. Flink: estado, tiempo y recuperación

> La documentación estable mostraba Flink 2.3 al construir este núcleo. DataStream API V2 y otras áreas pueden cambiar; fijar versión del job.

### 19.1 Modelo de tiempo

- **Event time:** timestamp del evento en su dominio de origen.
- **Processing time:** reloj del operador que procesa.
- **Watermark:** afirmación del pipeline sobre progreso de event time, no garantía de que ningún evento anterior llegará.

Definir generación de timestamps, idleness, out-of-orderness, lateness aceptada y destino de late data. Un watermark demasiado agresivo reduce espera pero clasifica más datos como tardíos.

Asignar la estrategia en la source cuando sea posible: allí puede aprovechar conocimiento por split/partition. Como el watermark combinado es el mínimo de sus entradas, una partition inactiva puede detenerlo; marcar idleness sólo tras un umbral compatible con el comportamiento real. En el extremo opuesto, una entrada rápida puede inflar sin límite el estado de joins/windows mientras espera a la lenta. Watermark alignment pausa fuentes rápidas dentro de un drift configurado, pero agrega coordinación y requiere que el source/splits soporten pause/resume; medir estado, drift y tiempo pausado.

### 19.2 Estado

Keyed state está particionado junto a keys; operator state pertenece a instancias de operador. Las keys deben derivarse determinísticamente. Estado necesita serializer/schema estable, TTL semántico y plan de rescale/upgrade. El cleanup por TTL es best effort: `NeverReturnExpired` puede impedir devolver un valor vencido aunque su eliminación física todavía esté pendiente; no interpretar TTL como borrado físico instantáneo.

Preguntar:

- ¿qué parte es working state y qué parte es source of truth?
- ¿cuánto tarda snapshot y restore en p99?
- ¿qué key groups/hot keys dominan?
- ¿qué cambia con paralelismo máximo?
- ¿TTL puede borrar estado necesario para un evento tardío?

### 19.3 Checkpoints y savepoints

**[CODE]** Flink usa snapshots distribuidos y replay. Checkpoint es automático para recovery; savepoint es activado para operaciones como upgrade/rescale. Barrier alignment puede agregar tail latency. Unaligned checkpoint incluye datos in-flight para que barriers adelanten buffers: estabiliza duración bajo backpressure, pero aumenta I/O/estado y puede empeorar si el storage ya es el cuello. Buffer debloating reduce datos in-flight; checkpoints incrementales guardan cambios respecto del checkpoint completado anterior y pueden bajar mucho el costo de estado grande. Ninguna opción sustituye eliminar el origen de backpressure.

Medir:

- checkpoint duration, alignment y bytes;
- failed/aborted checkpoints;
- time since last successful checkpoint;
- state size y growth;
- restore/recovery time;
- backpressure e input/output rate por subtask.

### 19.4 Exactly once preciso

El modo `EXACTLY_ONCE` asegura que, al recuperar, cada evento afecte una vez el estado administrado de operadores. No significa que el código físico ejecute una vez. End-to-end exige source replayable y sink transactional o idempotent con protocolo compatible.

Unaligned checkpoints tienen límites operativos: no interrumpen el procesamiento de un record largo, no admiten checkpoints unaligned concurrentes y pueden cambiar la disponibilidad del último watermark durante recovery si el operador no lo persiste. Probar explícitamente rescaling, broadcast/pointwise edges y recovery con channel state.

**[IMPL]** Inyectar fallo después del side effect y antes del checkpoint/commit. Si el resultado se duplica, la garantía declarada estaba fuera de alcance.

### 19.5 Backpressure y estado grande

Una dependencia lenta propaga presión upstream, demora watermarks y checkpoints y puede aumentar buffers. No “resolver” aumentando memoria sin límite. Separar:

- skew/hot key;
- operador CPU-bound;
- serialization/GC;
- network shuffle;
- state backend/storage;
- sink throttling;
- checkpoint alignment.

Las métricas `busyTimeMsPerSecond`, `idleTimeMsPerSecond` y `backPressuredTimeMsPerSecond` suman aproximadamente un segundo por subtask, pero son promedios de una ventana corta: igual promedio puede ocultar carga estable o ráfagas. Correlacionarlas con throughput, colas, checkpoint start delay, alignment, GC, storage y sink.

Gate Flink de producción:

- fijar max parallelism y UIDs estables antes de depender de savepoints;
- elegir state backend y storage durable, y verificar restore con el tamaño real;
- derivar checkpoint interval de RPO/reprocesamiento, visibilidad del sink y costo medido;
- configurar restart/failover y JobManager HA; con checkpointing sin estrategia explícita, la versión actual usa exponential delay por defecto;
- restringir cluster, REST y ejecución remota con TLS, autenticación y RBAC; no exponerlo públicamente;
- ensayar upgrade desde savepoint, rollback compatible, pérdida de TaskManager/JobManager y agotamiento del destino de checkpoints.

### 19.6 Código Flink — un barrier es parte de un protocolo

**[PROD]** Código upstream fijado en `2245026dff55e1348356ff472ee42c716710bbba` (Apache-2.0). El handler auditado identifica checkpoint y canales, ignora barriers obsoletos, cancela el checkpoint subsumido al avanzar, aborta ante cancellation/end-of-input según configuración y registra start delay, alignment duration y bytes durante alignment. El modo unaligned añade channel state; no elimina la máquina de estados ni el costo de recovery. Los tests del sink 2PC ejercitan notificaciones de checkpoints fuera de orden, fallo antes de notify, restore, abort y transaction timeout.

Gate derivado: para cada source/operator/sink, inyectar barrier viejo/duplicado/cancelado, canal lento/cerrado, checkpoint subsumido, task failure antes/después de snapshot y antes/después de commit, restore múltiple y timeout del sink. El test acepta solamente estado administrado convergente y side effects acordes al contrato; además limita tiempo/bytes de alignment, tamaño de channel state y antigüedad de la transacción. “Checkpoint completed” no basta si el commit del destino puede vencer, quedar ambiguo o no ser idempotente.

## 20. Cache y consistencia distribuida

Patrones:

- cache-aside;
- read-through/write-through;
- write-behind;
- invalidation/event-driven;
- versioned keys.

Fallos a modelar:

- fill concurrente y stampede;
- write DB → invalidation perdida;
- invalidation antes de replicación visible;
- negative cache obsoleta;
- eviction de metadata usada para corrección;
- region replica lag;
- hot key.

**[ACADEMIC/PROD]** El paper de Memcache en Facebook muestra que ordering entre replicación e invalidación cambia la probabilidad de recachear datos viejos. “Eliminar cache después del write” no es universalmente suficiente en multi-region.

Mitigaciones: version/CAS, request coalescing, jittered TTL, stale-while-revalidate cuando sea correcto, leases, invalidation durable y fallback con límites.

## 21. Baja latencia y datapath Linux

### 21.1 Primero localizar costo

Descomponer:

```text
wire → NIC queue → interrupt/NAPI → kernel stack → socket
→ wakeup/scheduler → userspace → parse → state update → decision
```

Instrumentar timestamps en dominios conocidos. Hardware timestamp, kernel timestamp y application timestamp no son intercambiables; documentar sincronización y punto de captura.

### 21.2 RSS, RPS, RFS y XPS

**[CODE]** Linux documenta:

- RSS distribuye flows en hardware receive queues;
- RPS realiza steering en software;
- RFS intenta dirigir procesamiento a la CPU de la aplicación para localidad;
- XPS selecciona transmit queue según CPU/RX queue.

Configurar IRQ, queues, CPU affinity y NUMA como conjunto. Si RSS ya da una cola por CPU local, RPS puede ser redundante. Medir reordering, imbalance, cache misses, IPIs e IRQ load.

No maximizar queues por regla fija. Más queues pueden acortar una cola saturada, pero también elevan IRQs y trabajo; para throughput alto suele convenir el menor número que no desborde ninguna CPU. Verificar `/proc/interrupts`, afinidad efectiva e interferencia de `irqbalance`. El hash RSS simétrico facilita colocar ambas direcciones juntas, pero reduce entropía y puede ampliar superficie de ataque: usarlo sólo con necesidad declarada.

### 21.3 Syscalls, batching y zero-copy

- `recvmmsg`/`sendmmsg` amortizan syscalls para datagramas.
- `sendfile`/`splice` pueden evitar roundtrip innecesario por userspace para casos compatibles.
- `MSG_ZEROCOPY` reemplaza copy cost por pinning/accounting/completions; Linux indica que normalmente empieza a ser efectivo alrededor de writes mayores a 10 KiB, cifra orientativa que debe medirse.
- io_uring puede reducir transiciones y soportar completion/batching; su beneficio depende del workload y kernel. Su zero-copy RX conserva el TCP stack del kernel, pero exige NIC con header/data split, flow steering y RSS adecuados.
- NAPI busy poll permite consultar antes de la interrupción y cambia CPU/energía por latencia. Afinidad, NAPI ID, budget, coalescing e IRQ suspension deben tratarse como un sistema; un core dedicado puede quedar al 100% aun sin trabajo útil.

“Zero-copy” siempre debe indicar qué copia se evita; headers, DMA, NIC y ownership de buffers aún existen.

### 21.4 AF_XDP

**[CODE]** AF_XDP conecta XDP con sockets user space mediante RX/TX rings y UMEM. FILL y COMPLETION transfieren ownership de frames entre kernel y aplicación; RX/TX transportan descriptors. Un programa XDP y `XSKMAP` dirigen cada queue al socket compatible. `XDP_ZEROCOPY` exige zero-copy o falla; sin forzarlo el bind puede caer a copy mode. La aplicación adquiere responsabilidad sobre rings, buffer lifecycle, backpressure, parsing seguro y observabilidad. Reusar simultáneamente el mismo frame en más de un ring puede corromper paquetes, y una completion TX no demuestra entrega por la red.

Gate antes de adoptar:

- el kernel stack medido es realmente el cuello;
- NIC/driver/kernel soportados y fijados;
- fallback y deployment reproducibles;
- aislamiento de queues/CPU/NUMA;
- protección, rate control y recovery reimplementados;
- mejora p99/p99.9 end-to-end, no sólo Mpps sintético.

### 21.5 DPDK

**[CODE]** Los Poll Mode Drivers consultan RX/TX descriptors en userspace, normalmente sin interrupciones, y usan bursts, mbufs/mempools, hugepages y recursos por core. Modelos comunes: run-to-completion —recibir, procesar y transmitir en el mismo lcore— y pipeline —pasar trabajo entre lcores mediante rings—.

Tradeoffs:

- núcleos dedicados y consumo energético;
- operación/seguridad más complejas;
- NUMA y queue ownership estrictos;
- necesidad de reimplementar o integrar funciones del stack;
- batching puede elevar latencia de mensajes aislados;
- un pipeline agrega rings y cache transfers.

Usar una RX queue por lcore owner cuando corresponda; las APIs PMD lock-free suponen que dos lcores no operan concurrentemente sobre la misma queue salvo capacidad explícita. Mantener NIC, lcore, descriptors, mbufs y mempool en el mismo nodo NUMA. Los bursts amortizan costos y aprovechan prefetch/vectorización, pero aumentan espera hasta reunir trabajo: comparar run-to-completion contra pipeline y tamaños de burst con el tráfico real.

### 21.6 RDMA

RDMA separa un slow path de creación/gestión vía kernel de un fast path que normalmente escribe directamente registros de hardware mapeados en userspace, sin syscall/context switch por operación. Ofrece send/receive y, según transporte/hardware, one-sided read/write/atomics. Reduce intervención de CPU remota, pero no elimina coherencia, autorización, memory registration/keys, QP/CQ lifecycle, congestion, failure recovery ni versioning. Un `completion` prueba la semántica definida por el verb/transport, no que la lógica remota haya consumido el dato.

**[ACADEMIC]** FaRM demuestra una combinación específica de OCC, replication y one-sided RDMA. No generalizar sus resultados sin hardware, topology, workload y protocolo equivalentes.

## 22. Streaming para estado de order book

Este manual cubre el pipeline, no la microestructura ni el protocolo de un venue.

State machine genérica:

```text
DISCONNECTED
  → CONNECTING/AUTH
  → BUFFERING_DELTAS
  → FETCHING_SNAPSHOT
  → RECONCILING(snapshot, buffered deltas)
  → LIVE
  → GAP_DETECTED
  → RECOVERING or FULL_RESYNC
```

Invariantes a especializar con la especificación oficial:

- cada update se valida antes de mutar;
- sequence/epoch sólo se compara en su scope;
- duplicate no se vuelve segundo efecto;
- gap impide declarar estado `LIVE` correcto;
- snapshot y deltas se combinan con regla exacta del venue;
- reconnect nunca supone continuidad;
- estado publicado incluye freshness y health;
- parsing o recovery nunca bloquea indefinidamente el receive path;
- pérdida local por overflow es observable y fuerza política segura.

Arquitectura de baja latencia posible:

```text
NIC/RX owner → framing+validation → sequence gate → book state owner
→ immutable/versioned publication → consumers
```

Un solo owner por book/partition puede simplificar corrección y locality. Sólo introducir paralelismo dentro del mismo orden si conserva determinismo y demuestra beneficio.

La semántica exacta pertenece a `MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md`: snapshot/delta de Binance/Coinbase/Kraken, ITCH/OUCH, FIX session, MDP/EOBI y estados de derivados no se deducen de TCP/UDP/WebSocket. Cruce recíproco auditado el 2026-08-21:

- A/B multicast son feeds redundantes bajo contrato, no particiones para concatenar;
- payload sequence y transport sequence pueden ser espacios distintos;
- retransmission FIX no equivale a application resend ni garantiza idempotencia de orden;
- reconnect/resubscribe crea una nueva epoch hasta demostrar continuidad;
- timestamp no reemplaza sequence ni crea orden total entre venues;
- local overflow se trata como pérdida del feed aunque la red haya entregado correctamente.

Mediciones:

- wire/kernel/app receive timestamps;
- decode y apply latency;
- sequence gaps/duplicates/out-of-order;
- queue depth/age y drops;
- resync count/duration;
- book freshness;
- CPU migrations, cache misses y allocation;
- p50/p99/p99.9/max bajo bursts reales.

## 23. Observabilidad y diagnóstico

### 23.1 Señales por capa

| Capa | Señales |
|---|---|
| aplicación | rate, goodput, error class, deadlines, queue age |
| stream/broker | partition lag, event age, rebalances, checkpoint health |
| socket/TCP | RTT, retransmits, cwnd, send/receive queues, resets |
| IP/path | loss, MTU/ICMP, route/zone/provider |
| host | CPU, run queue, IRQ/softirq, NUMA, memory, GC |
| NIC | RX/TX queues, drops, errors, ring starvation, timestamps |

Una métrica agregada puede esconder un shard muerto. Mantener cardinalidad controlada pero permitir drill-down por cluster/topic/partition/zone/error.

### 23.2 Correlación

Propagar IDs de request, operation, trace, event, causation y partition/offset según necesidad. No usar trace ID como idempotency key sin contrato. Logs deben registrar decisión y estado, no payload secreto completo.

### 23.3 Diagnóstico científico

```text
síntoma → timeline → hipótesis → predicción diferenciadora
→ experimento/trace → resultado → corrección → regresión → prevención
```

Conservar artifacts: configs, versiones, topología, packet capture acotado, métricas, logs, clock sync y fault injection. “La red estuvo lenta” no es diagnóstico.

## 24. Testing de red y distribución

### 24.1 Capas de prueba

1. parser/framing con unit, property y fuzz tests;
2. state machines deterministas con model-based tests;
3. protocolos con simulación de pérdida/reordenamiento/duplicación;
4. integración multi-process;
5. fault injection en red, proceso, disco y clock;
6. performance bajo carga ofrecida y burst;
7. soak/recovery/rolling upgrade;
8. historia de operaciones comprobada contra el modelo de consistencia.

### 24.2 Matriz mínima de fallos

| Falla | Momento crítico | Evidencia esperada |
|---|---|---|
| pérdida de request | antes de commit | retry/reconcile correcto |
| pérdida de response | después de commit | no duplica efecto lógico |
| caída de líder | antes/después de replicate | historia válida |
| partición | minoría/mayoría/asimétrica | política declarada |
| disco lento/lleno | snapshot/log/commit | backpressure, no corrupción |
| pausa larga | owner/lease/consumer | fencing/rebalance seguro |
| reorder/duplicate | secuencia | detector y dedupe |
| clock jump/skew | deadline/event time | no rompe invariante |
| version skew | rolling deploy | compatibilidad demostrada |
| overload | burst/recovery | rechazo acotado y recuperación |

### 24.3 Jepsen y model checking

**[CODE/PROD]** Jepsen ejecuta histories mientras inyecta fallos y las compara con modelos; encuentra interacciones que unit tests sanos no ven. No prueba todas las ejecuciones. Para algoritmos propios, complementar con especificación/model checking (por ejemplo TLA+) y tests de implementación.

## 25. Benchmarking profesional

Registrar:

```yaml
hardware: {cpu: "", numa: "", nic: "", link: "", switches: ""}
software: {os: "", kernel: "", firmware: "", runtime: "", versions: ""}
network: {rtt: "", loss: "", mtu: "", tls: "", topology: ""}
workload: {connections: 0, msg_sizes: [], key_skew: "", burst: ""}
semantics: {acks: "", durability: "", consistency: ""}
measurement: {clock: "", warmup: "", duration: "", samples: 0}
results: {offered_load: "", goodput: "", percentiles: {}, drops: 0}
```

Reglas:

- reportar distribución y máximos con contexto, no sólo promedio;
- separar coordinated omission;
- mantener offered load aunque el sistema se atrase para observar overload real;
- incluir recovery y steady state;
- comparar misma corrección/durabilidad;
- medir desde el límite que importa al usuario;
- repetir y publicar variabilidad/configuración.

Para paths de red, añadir una campaña temporal suficientemente larga: delivery rate, ráfaga máxima de pérdidas, racha de éxitos, min/max RTT, CDF, serie temporal y correlación/autocorrelación. Pérdidas con el mismo promedio pueden tener consecuencias radicalmente distintas si son independientes o bursty. `ping`/`traceroute` caracterizan ICMP y una ruta observada; no sustituyen medición del protocolo de aplicación.

## 26. Gates de diseño y producción

### 26.1 Gate de protocolo

1. ¿Framing, versión, tamaños y unknown fields están definidos?
2. ¿Orden, duplicados, gaps, retry y replay tienen scope exacto?
3. ¿Timeout conduce a estado inequívoco o a reconciliación?
4. ¿Autenticación, autorización, integridad y límites cubren cada mensaje?
5. ¿Cierre, reconnect, resume y version skew están probados?

### 26.2 Gate distribuido

1. ¿Invariantes y modelo de consistencia son comprobables?
2. ¿Failure model incluye pausas, particiones y storage?
3. ¿Commit y durability dicen cuántas réplicas/qué medio/qué ACK?
4. ¿Elección, fencing, snapshot y membership son seguros?
5. ¿Cada efecto ambiguo se deduplica o reconcilia?
6. ¿Recovery tiene presupuesto de recursos y no causa overload recursivo?

### 26.3 Gate streaming

1. ¿Partition key preserva orden necesario y distribuye carga?
2. ¿Schema evolution se probó en ambas direcciones relevantes?
3. ¿Offsets/checkpoints y side effects forman una garantía real end-to-end?
4. ¿Late data, gaps, poison records y replay tienen política?
5. ¿Lag, edad, backpressure, checkpoint y restore son observables?

### 26.4 Gate baja latencia

1. ¿Se midió el camino completo y se localizó el costo?
2. ¿Affinity, queues y memoria respetan NUMA/locality?
3. ¿Batching/copies/syscalls mejoran p99 sin romper freshness?
4. ¿Kernel bypass conserva seguridad, flow/congestion control y operación?
5. ¿La optimización sobrevive bursts, fallos y competencia de recursos?

## 27. Contrato de instrucciones para Codex

Usar este bloque al pedir implementación:

```md
Objetivo y usuarios:
Invariantes:
Topología y trust boundaries:
Protocolo y versión normativa:
Framing/schema/límites:
Modelo de consistencia y alcance del orden:
Failure model:
Timeout/retry/idempotency/reconciliation:
Backpressure y overload:
Durabilidad y recovery:
SLO/carga/percentiles:
Lenguaje/runtime/kernel/hardware fijados:
Tests de fallos obligatorios:
Métricas/traces/runbook:
Cambios fuera de alcance:
Incertidumbres [OPEN]:
```

Instrucciones permanentes al agente:

- citar RFC/documentación versionada para semántica;
- no inventar exactly-once, orden global ni atomicidad remota;
- implementar parser incremental con límites antes de allocations;
- hacer cleanup/cancellation idempotentes;
- acotar toda cola, pool, retry y buffer;
- preservar request/event IDs a través de retries;
- crear tests de efecto ambiguo y fault injection;
- entregar benchmark reproducible, no adjetivos de rendimiento;
- registrar `[OPEN]` cuando hardware/proveedor/protocolo no esté fijado.

## 28. Auditoría transversal cerrada

### 28.1 Con sistemas y rendimiento

Cruce completado el 2026-08-21 con `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md`:

| Decisión de red/distribución | Obligación de sistemas/rendimiento |
|---|---|
| SLO end-to-end | medir carga ofrecida, goodput y p50/p95/p99/p99.9 por fase |
| queue/worker/shard | ligar ownership a CPU, cache, NUMA, IRQ/NAPI y capacidad |
| batching/compresión | medir espera agregada, CPU, bytes y efecto sobre tails |
| polling/zero-copy/bypass | demostrar copia/syscall evitado y costo de pinning, completions o core dedicado |
| backpressure | límites en cada queue/ring/pool, edad, rechazo y propagación |
| timestamp | dominio, resolución, sincronización y punto exacto de captura |
| retransmisión/retry | trabajo amplificado bajo saturación y presupuesto total |
| benchmark | mismo path, payload, seguridad, afinidad, warmup y duración que producción |

Una mejora de Mpps, syscall rate o promedio no autoriza una regresión de corrección o p99.9. Perfilar primero determina si el cuello está en NIC/IRQ, kernel, scheduler, copia, parsing, state machine, coordinación o dependencia.

### 28.2 Con backend y APIs

Cruce completado el 2026-08-21 con `SOFTWARE_BACKEND_API_ENGINEERING.md`:

| Contrato backend/API | Obligación de red/distribución |
|---|---|
| método/operación | declarar safe/idempotent o deduplicación antes de retry/hedging |
| deadline/cancelación | presupuestar DNS, connect, TLS, queue, write, service y read; propagar restante |
| HTTP/gRPC/WebSocket | preservar framing, flow control, partial I/O, cierre y límites de mensaje |
| respuesta perdida | modelar efecto ambiguo y ofrecer clave/estado/reconciliación |
| auth/TLS | autenticar peer y autorizar recurso/acción; rotación/reconnect no debe abrir bypass |
| schema/version | compatibilidad producer/consumer, unknown fields y límites antes de parsear |
| evento/stream | definir partition key, orden, delivery, commit y side effect externo |
| observabilidad | correlación sin secretos entre request, RPC, mensaje, partition y estado |

Cambiar HTTP por gRPC/QUIC, una llamada local por RPC o una cola por Kafka cambia fallos y latencia, no los invariantes de negocio. El contrato debe sobrevivir retransmisión, duplicado, reorder, disconnect, rebalance y recovery.

### 28.3 Con bases de datos y almacenamiento

Cruce completado el 2026-08-21 con `DATABASE_STORAGE_INTERNALS.md`:

| Contrato distribuido | Obligación de datos |
|---|---|
| ACK/quorum | nombrar WAL/data durable, réplicas que confirmaron y visibilidad; “replicado” no equivale a backup |
| líder/failover | fencing del líder anterior, commit index/LSN y regla de stale reads/rejoin |
| 2PC | logs de prepare/decision, recovery y bloqueo; consensus no reemplaza atomic commit de participantes |
| evento externo | outbox/inbox o sink transaccional/idempotente; commit Kafka/Flink no atomiciza una API/DB ajena |
| CDC/bootstrap | snapshot y cursor/LSN/offset forman un corte; retention debe cubrir caída y replay más lentos |
| backpressure | pool/locks/WAL/compaction/checkpoint propagan saturación al stream; acotar buffers y lag |
| cache/replica | invalidation y write visibility tienen orden explícito; medir ventana stale y failure path |
| recuperación | partition, failover, disk full, missing WAL/segment y restore deben preservar invariantes |

La prueba end-to-end inyecta caída antes/después de persistir, replicar, responder, publicar y confirmar offset. Se acepta cada historia sólo si satisface el contrato, no porque el sistema finalmente “parezca converger”.

### 28.4 Con GPU y cómputo acelerado

Cruce completado el 2026-08-21 con `GPU_ACCELERATED_COMPUTING.md`:

| Frontera red↔GPU | Contrato compartido |
|---|---|
| NIC/RDMA/GPUDirect | inventariar NIC–PCIe–NUMA–GPU y compatibilidad real; acceso directo no elimina ownership, orden ni fallback |
| collective NCCL | retorno de API/group confirma enqueue, no completion; observar estado async, timeout y error de communicator |
| backpressure | acotar bytes, descriptors, buffers pinned/device y requests; propagar saturación hasta admisión |
| orden/lifetime | registrar stream/event que protege cada buffer; no reciclar RX/TX/device memory antes de completion observable |
| rank/link failure | abort/revoke/restart coordinado; un resultado parcial o de communicator fallido no se publica |
| medición | correlacionar event/receive time, queue, transfer, kernel, collective y entrega; separar throughput de goodput |

RDMA o GPUDirect reducen copias potenciales; no crean exactly-once, durabilidad ni consistencia. La prueba de fallo corta la ruta antes/después de recepción, DMA, kernel, collective y respuesta y verifica que no haya hang, reuse prematuro ni output inválido.

### 28.5 Con seguridad, SRE y cloud

Cruce completado el 2026-08-21 con `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md`:

| Frontera de red/distribución | Contrato compartido |
|---|---|
| peer y canal | TLS/mTLS autentica identidad esperada, valida lifecycle de cert/keys y no reemplaza authz del recurso |
| ingress/egress | source/forwarded headers sólo desde proxies confiables; DNS/redirect/IPv4/IPv6/metadata y destinations bajo policy |
| DoS/overload | limitar estado y asimetría pre-auth; admission, byte/work/concurrency quotas, retry budget y shedding por prioridad |
| segmentación | policy por flujo e identidad realmente aplicada por firewall/CNI; deny-all futuro, metadata y control planes no públicos |
| protocolo/stream | malformed/replay/duplicate/reorder/slowloris y amplification dentro del threat/failure model; secrets/PII fuera de payload/telemetría |
| cloud/failure domain | LB/DNS/zone/region/egress/quota y provider control plane modelados; failover conserva identidad, fencing, capacidad y data cut |
| incidente | packet/flow/app evidence correlacionable con retención/acceso acotados; isolation y revocation no destruyen evidencia ni recovery |

Zero trust no cambia TCP, delivery ni consenso: elimina confianza por ubicación y exige identity/policy/freshness/failure mode. Cifrar no oculta necesariamente IP, endpoints, timing, volumen o patrones; privacidad y anonimato requieren threat model propio.

### 28.6 Con frontend y experiencia en navegador

Cruce completado el 2026-08-21 con `FRONTEND_PRODUCT_ENGINEERING_UX.md`:

| Frontera de red/stream | Obligación frontend |
|---|---|
| connect/DNS/TLS/HTTP | medir waterfall y deadline real; distinguir offline, timeout, abort y respuesta inválida |
| WebSocket/SSE | reconnect implica nueva sesión; epoch/sequence/resubscribe/resync según protocolo |
| ordering/duplicate/gap | state machine descarta stale/duplicate y muestra `GAPPED/RESYNCING`, nunca continuidad inventada |
| backpressure | limitar buffers y frecuencia de render; coalescing declara si conserva estado o historial |
| cache/CDN/service worker | keys, `Vary`, freshness, version y rollback; no servir shell/data incompatibles |
| telemetry | timings browser↔edge↔service correlacionables sin exponer payload/identidad sensible |

“Connected” no demuestra reachability de la dependencia ni freshness del dato. La UI conserva estados de transporte separados del outcome de negocio y no convierte delivery al menos una vez en una experiencia exactamente una vez ficticia.

### 28.7 Con algoritmos y estructuras de datos

Cruce completado el 2026-08-21 con `ALGORITHMS_DATA_STRUCTURES_PROBLEM_SOLVING.md`:

| Problema de red/distribución | Obligación algorítmica |
|---|---|
| buffers, retransmisión y timers | estructura acotada, expiración, orden y costo adversarial explícitos |
| routing/topología/dependencias | modelo de grafo correcto; condición de pesos para shortest path y ciclo |
| partition/rebalance | función, skew, movimiento y estabilidad medidos; hash no implica balance perfecto |
| stream/window/dedup | memoria, pasadas, error, lateness y mergeability declarados |
| consensus/replicación | no sustituir proof de safety/liveness por similitud con BFS, queue o log local |
| overload | analizar secuencia y peor caso; amortizado no garantiza cada tail ni memoria acotada |

El manual algorítmico gobierna estructura, prueba y costo abstracto; éste gobierna partial failure, comunicación, tiempo, delivery y recuperación. Una simulación secuencial correcta no demuestra corrección distribuida.

Cruce con `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md` cerrado el 2026-08-21: arquitectura elige distribución sólo por drivers medibles; este manual verifica que sync/async, service/event boundaries, replicas y regions expresen correctamente delivery, ordering, partition, consensus y recovery. Microservices o event-driven no eliminan coupling: lo trasladan a contratos, tiempo, operación y version skew, que deben aparecer en runtime/deployment views y scenarios.

Cruce con `DATA_ENGINEERING_ANALYTICS.md` cerrado el 2026-08-21: red/streaming gobierna delivery, partition, offset, checkpoint, watermark transport y backpressure; datos gobierna grain, event/effective time, window result, late correction, lineage y publication. Broker/Flink checkpoint no atomiciza sink externo; replay y backfill deben producir output contractual o una divergencia versionada, con backlog/retention suficientes.

Cruce con `NATIVE_MOBILE_DESKTOP_ENGINEERING.md` cerrado el 2026-08-21: una interfaz disponible no prueba reachability; clientes móviles/escritorio manejan cambio de red, captive portal, background, reconnect, auth refresh, sequence/gap y bounded buffers. Red define delivery/ordering/resume/backpressure; nativo decide lifecycle y frecuencia visible. Push es señal no confiable y un stream no puede hacer render por evento si excede frame/battery budget.

### 28.8 Práctica pública Cloudflare/AWS — proxy evolutivo y control de escala

**[CODE] Evidencia fijada:** `cloudflare/pingora@0046038bd402bc82912da862dadf9a479f31e9f1`, Apache-2.0. El workspace Rust publica core/proxy, pools, load balancing, health checks, límites, timeouts, cache y observabilidad, con pruebas de shutdown, pooling y reconstrucción de selectors. Cloudflare declara uso de la familia Pingora a gran escala; el repositorio no reproduce su configuración, tráfico ni todos sus servicios internos.

Patrones transferibles:

- lifecycle del request dividido en fases/hooks tipados en lugar de lógica monolítica;
- pool de conexiones con límite global, identidad completa del peer, idle validation/timeout y eviction;
- discovery, selector y health checking con schedules/ownership separados y estado previo válido durante refresh fallido;
- graceful reload que coordina readiness, deja de aceptar y drena dentro de deadline;
- límites de inflight/rate y cache lock para evitar stampede;
- retry marcado por clase de error y presupuesto, no repetición indiscriminada.

**[PROD] Scale inversion de AWS:** cuando una flota grande consulta/controla a un servicio mucho menor, recovery o polling sincronizado puede sobrecargarlo. El componente pequeño debe controlar el ritmo: particionar pull, distribuir estado por un medio escalable, aplicar leases/tokens o hacer que clientes consuman cambios a una cadencia acotada. Añadir clientes no puede multiplicar linealmente trabajo en un control plane central durante recuperación.

**Gate de proxy/control plane:**

1. key de pool/cache incluye todos los factores de autoridad y seguridad —host/`:authority`, scheme/TLS peer, método y variantes aplicables—;
2. probar peer stale/half-closed, DNS/discovery stale, health flapping, selector rebuild fallido y generación superseded;
3. slowloris, body/header límites, backpressure, timeout por fase y cancelación;
4. retry sólo antes/después del commit point permitido y con request body replayable;
5. graceful reload bajo conexiones largas/streaming y deadline de drain;
6. cache stampede, poisoning, cross-tenant isolation, purge/version y bypass behavior;
7. offered load constante durante backend/control-plane down/recovery; recursos y goodput permanecen acotados.

**Aprendizaje de error real:** una advisory de Pingora corrigió una cache key default que omitía factores como host, con riesgo de poisoning/cross-origin. La lección normativa no es “Pingora es inseguro”, sino que memoria segura y experiencia de escala no prueban semántica de caché: eliminar defaults ambiguos, actualizar a una revisión corregida y convertir cada dimensión de autoridad en negative tests.

Fuentes: [Pingora fijado](https://github.com/cloudflare/pingora/tree/0046038bd402bc82912da862dadf9a479f31e9f1), [arquitectura de producción](https://blog.cloudflare.com/how-we-built-pingora-the-proxy-that-connects-cloudflare-to-the-internet/), [advisory de cache key](https://github.com/cloudflare/pingora/security/advisories/GHSA-f93w-pcj3-rggc) y [scale inversion de AWS](https://aws.amazon.com/builders-library/avoiding-overload-in-distributed-systems-by-putting-the-smaller-service-in-control/).

## 29. Fuentes primarias auditadas

### 29.1 Academia y papers

- Stanford CS144 Fall 2025: <https://cs144.github.io/>. Se auditaron sus ocho checkpoints públicos: socket/byte stream, reassembly, TCP receiver/sender, medición de paths, interfaz Ethernet/ARP, router y capstone.
- MIT 6.5840 Spring 2026: <https://pdos.csail.mit.edu/6.824/>
- Schedule y corpus 6.5840: <https://pdos.csail.mit.edu/6.824/schedule.html>
- MapReduce: <https://pdos.csail.mit.edu/6.824/papers/mapreduce.pdf>
- Google File System: <https://pdos.csail.mit.edu/6.824/papers/gfs.pdf>
- Paxos Made Simple: <https://pdos.csail.mit.edu/6.824/papers/paxos-simple.pdf>
- Raft extended: <https://pdos.csail.mit.edu/6.824/papers/raft-extended.pdf>
- Linearizability: <https://pdos.csail.mit.edu/6.824/papers/p463-herlihy.pdf>
- ZooKeeper: <https://pdos.csail.mit.edu/6.824/papers/zookeeper.pdf>
- Spanner: <https://pdos.csail.mit.edu/6.824/papers/spanner.pdf>
- Chain Replication: <https://pdos.csail.mit.edu/6.824/papers/cr-osdi04.pdf>
- FaRM: <https://pdos.csail.mit.edu/6.824/papers/farm-2015.pdf>
- Memcache at Facebook: <https://pdos.csail.mit.edu/6.824/papers/memcache-fb.pdf>
- IronFleet: <https://pdos.csail.mit.edu/6.824/papers/ironfleet.pdf>
- AWS Lambda on-demand container loading: <https://pdos.csail.mit.edu/6.824/papers/atc23-brooker.pdf>
- Ray distributed futures/ownership: <https://pdos.csail.mit.edu/6.824/papers/ray.pdf>
- SUNDR: <https://pdos.csail.mit.edu/6.824/papers/li-sundr.pdf>
- Bitcoin: <https://pdos.csail.mit.edu/6.824/papers/bitcoin.pdf>
- Practical Byzantine Fault Tolerance: <https://pdos.csail.mit.edu/6.824/papers/castro-practicalbft.pdf>
- Dynamo, Amazon Science: <https://www.amazon.science/publications/dynamo-amazons-highly-available-key-value-store>
- Chandy–Lamport, Microsoft Research: <https://www.microsoft.com/en-us/research/people/lamport/publications/>
- Jepsen: <https://github.com/jepsen-io/jepsen>

### 29.2 IETF/RFC

- UDP RFC 768, usage guidelines RFC 8085 y Datagram PLPMTUD RFC 8899: <https://www.rfc-editor.org/rfc/rfc768.html>, <https://www.rfc-editor.org/rfc/rfc8085.html>, <https://www.rfc-editor.org/rfc/rfc8899.html>
- IPv6 RFC 8200: <https://www.rfc-editor.org/rfc/rfc8200.html>
- TCP RFC 9293, congestion RFC 5681 y RTO RFC 6298: <https://www.rfc-editor.org/rfc/rfc9293.html>, <https://www.rfc-editor.org/rfc/rfc5681.html>, <https://www.rfc-editor.org/rfc/rfc6298.html>
- QUIC RFC 9000/9001/9002: <https://www.rfc-editor.org/rfc/rfc9000.html>, <https://www.rfc-editor.org/rfc/rfc9001.html>, <https://www.rfc-editor.org/rfc/rfc9002.html>
- WebSocket RFC 6455 y HTTP/2 CONNECT RFC 8441: <https://www.rfc-editor.org/rfc/rfc6455.html>, <https://www.rfc-editor.org/rfc/rfc8441.html>
- DNS RFC 1034/1035: <https://www.rfc-editor.org/rfc/rfc1034.html>, <https://www.rfc-editor.org/rfc/rfc1035.html>
- HTTP/2 y HTTP/3 RFC 9113/9114: <https://www.rfc-editor.org/rfc/rfc9113.html>, <https://www.rfc-editor.org/rfc/rfc9114.html>

### 29.3 Implementaciones oficiales

- Apache Kafka 4.3: <https://kafka.apache.org/43/>
- Kafka design: <https://kafka.apache.org/43/design/design/>
- Kafka producer/consumer configs, KRaft, rebalance y upgrades: <https://kafka.apache.org/43/generated/producer_config.html>, <https://kafka.apache.org/43/generated/consumer_config.html>, <https://kafka.apache.org/43/operations/kraft/>, <https://kafka.apache.org/43/operations/consumer-rebalance-protocol/>, <https://kafka.apache.org/43/getting-started/upgrade/>
- Kafka fijado `e28414ee4a1dfa3d093dc10168f944f1a759434e`: [TransactionManager](https://github.com/apache/kafka/blob/e28414ee4a1dfa3d093dc10168f944f1a759434e/clients/src/main/java/org/apache/kafka/clients/producer/internals/TransactionManager.java), [ProducerStateManager](https://github.com/apache/kafka/blob/e28414ee4a1dfa3d093dc10168f944f1a759434e/storage/src/main/java/org/apache/kafka/storage/internals/log/ProducerStateManager.java) y [log recovery tests](https://github.com/apache/kafka/blob/e28414ee4a1dfa3d093dc10168f944f1a759434e/core/src/test/scala/unit/kafka/server/LogRecoveryTest.scala).
- Apache Flink stable: <https://nightlies.apache.org/flink/flink-docs-stable/>
- Flink watermarks, state, checkpoints/backpressure y production: <https://nightlies.apache.org/flink/flink-docs-stable/docs/dev/datastream/event-time/generating_watermarks/>, <https://nightlies.apache.org/flink/flink-docs-stable/docs/dev/datastream/fault-tolerance/state/>, <https://nightlies.apache.org/flink/flink-docs-stable/docs/ops/state/checkpoints/>, <https://nightlies.apache.org/flink/flink-docs-stable/docs/ops/state/checkpointing_under_backpressure/>, <https://nightlies.apache.org/flink/flink-docs-stable/docs/ops/production_ready/>
- Flink fijado `2245026dff55e1348356ff472ee42c716710bbba`: [barrier handler](https://github.com/apache/flink/blob/2245026dff55e1348356ff472ee42c716710bbba/flink-runtime/src/main/java/org/apache/flink/streaming/runtime/io/checkpointing/SingleCheckpointBarrierHandler.java) y [2PC sink recovery tests](https://github.com/apache/flink/blob/2245026dff55e1348356ff472ee42c716710bbba/flink-streaming-java/src/test/java/org/apache/flink/streaming/api/functions/sink/legacy/TwoPhaseCommitSinkFunctionTest.java).
- Linux networking scaling: <https://docs.kernel.org/networking/scaling.html>
- Linux NAPI: <https://docs.kernel.org/networking/napi.html>
- Linux AF_XDP: <https://docs.kernel.org/networking/af_xdp.html>
- Linux `MSG_ZEROCOPY`: <https://docs.kernel.org/networking/msg_zerocopy.html>
- Linux io_uring zero-copy receive: <https://docs.kernel.org/networking/iou-zcrx.html>
- DPDK Programmer’s Guide 26.07 observado y PMD: <https://doc.dpdk.org/guides/prog_guide/>, <https://doc.dpdk.org/guides/prog_guide/ethdev/ethdev.html>
- Linux userspace RDMA verbs: <https://docs.kernel.org/infiniband/user_verbs.html>
- rdma-core: <https://github.com/linux-rdma/rdma-core>

## 30. Extensiones deliberadamente condicionadas

Este núcleo ya permite implementar y revisar redes, coordinación y streaming generales. La auditoría cerró las afirmaciones operativas incorporadas; no promete haber reproducido cada párrafo de cada corpus. Permanecen condicionados:

- completar el pase de las lecture notes de CS144; sus ocho checkpoints públicos ya fueron auditados;
- completar la lectura integral de cada paper 6.5840 seleccionado; el schedule completo y los modelos/claims centrales ya fueron contrastados;
- completar RFC complementarias y errata para casos raros de TCP/QUIC/DNS; el núcleo operativo UDP/TCP/QUIC/WebSocket/MTU ya fue reconciliado con las RFC citadas;
- extender, cuando un proyecto lo requiera, Kafka Connect/Streams y Flink SQL/connectors; el núcleo Kafka 4.3 y Flink 2.3 de garantías, configuración, recovery, backpressure, upgrade y readiness ya fue auditado;
- los cruces con bases de datos, seguridad/SRE/cloud, microestructura y frontend ya quedaron cerrados;
- añadir reglas de venue/NIC/kernel/runtime únicamente cuando un proyecto declare versiones concretas.
