# Software Architecture & System Design — manual operativo

> **Estado:** núcleo general v1 auditado el 2026-08-21; autoridades, integridad documental y cruces recíprocos cerrados. ADR, quality scenarios, trazabilidad y capstones ya tienen contrato ejecutable en `ENGINEERING_EXECUTION_PLAYBOOK.md`; la evidencia se produce por proyecto.
> **Autoridad académica:** Carnegie Mellon Software Engineering Institute (SEI); MIT 6.102 Software Construction.
> **Autoridad normativa:** ISO/IEC/IEEE 42010:2022 e ISO/IEC 25010:2023.
> **Autoridad práctica:** *Software Engineering at Google*; C4 de Simon Brown; arc42 de Gernot Starke y Peter Hruschka; ADR de Michael Nygard; refactoring de Martin Fowler.
> **Propósito:** transformar objetivos y restricciones en una arquitectura justificable, construible, evaluable, desplegable y capaz de evolucionar.

## Cómo usar este manual

```text
drivers y restricciones
→ escenarios de calidad medibles
→ modelo de dominio y límites
→ alternativas y trade-offs
→ vistas + decisiones
→ skeleton/probes
→ implementación incremental
→ tests de arquitectura y sistema
→ operación y feedback
→ evolución o retiro
```

No empezar por tecnología. “Kubernetes + microservices + Kafka” no es una arquitectura mientras no se sepa qué problema resuelve, qué cualidad mejora, qué costo añade y cómo se comprobará.

## Procedencia

- `[SPEC]`: estándar o contrato oficial versionado.
- `[ACADEMIC]`: síntesis de curso, método o investigación identificada.
- `[CODE]`: implementación/repositorio oficial usado como evidencia.
- `[PROD]`: experiencia operacional atribuida.
- `[MEASURED]`: experimento reproducible en entorno declarado.
- `[IMPL]`: traducción propia a decisión, documento, código, test o runbook.
- `[OPEN]`: afirmación todavía no verificada.

## 1. Qué es una decisión arquitectónica

**[ACADEMIC]** SEI describe la arquitectura como decisiones sobre estructura y comportamiento global que permiten razonar temprano sobre cualidades como modifiability, availability y security. La arquitectura no es un diagrama: es el conjunto de decisiones significativas y sus consecuencias observables.

Una decisión es arquitectónicamente significativa si modifica al menos uno:

- límites, responsabilidades o dependencias;
- contratos externos o modelos de datos duraderos;
- distribución, comunicación o ownership;
- seguridad, privacidad o trust boundaries;
- rendimiento, disponibilidad o recovery;
- build, deploy, operación o evolución;
- trabajo de varios equipos o costo difícil de revertir.

### 1.1 Arquitectura, diseño e implementación

| Nivel | Pregunta |
|---|---|
| arquitectura | ¿qué decisiones condicionan el sistema completo y sus cualidades? |
| diseño | ¿cómo colaboran módulos y tipos dentro de esos límites? |
| implementación | ¿cómo se expresa correctamente en lenguaje/runtime concreto? |

La frontera depende del contexto. Elegir un algoritmo puede ser local; en un hot path o formato persistente puede volverse arquitectónico. La reversibilidad y el costo de cambio importan más que el título del documento.

### 1.2 Arquitectura mínima suficiente

Diseñar lo suficiente para:

1. eliminar ambigüedad peligrosa;
2. exponer riesgos costosos;
3. permitir trabajo paralelo sin contratos ficticios;
4. probar las cualidades dominantes;
5. conservar opciones donde todavía hay incertidumbre.

No intentar adivinar cada detalle futuro. Tampoco llamar “emergente” a evitar decisiones irreversibles hasta producción.

## 2. Drivers: propósito, stakeholders y restricciones

### 2.1 Ficha de contexto

```text
problema/oportunidad:
usuarios y tareas:
resultado de negocio/mission:
stakeholders y concerns:
alcance y fuera de alcance:
sistemas externos y owners:
datos, sensibilidad y residencia:
carga actual, horizonte y crecimiento:
SLO/RTO/RPO y riesgo máximo:
presupuesto, plazo y equipo:
tecnología/regulación obligatoria:
decisiones ya tomadas:
incertidumbres y experimentos:
```

No convertir deseos en requisitos. “Escalable”, “seguro”, “flexible” y “tiempo real” son etiquetas hasta que existe estímulo, ambiente, respuesta y medida.

### 2.2 Features, qualities y constraints

| Categoría | Ejemplo | Tratamiento |
|---|---|---|
| feature | aceptar una orden | flujo y outcome |
| quality | 99,99% disponible | escenario medible y táctica |
| constraint | datos sólo en región X | límite no negociable |
| assumption | proveedor entrega en 200 ms | validar y monitorizar |
| risk | rate limit desconocido | probe, fallback o decisión |

No esconder una constraint como preferencia. No presentar una assumption como hecho.

### 2.3 Priorización

Máximo tres a cinco quality drivers dominantes. Registrar tensiones:

```text
consistency ↔ availability durante partición
latencia ↔ batching/eficiencia
aislamiento ↔ costo/utilización
autonomía de equipo ↔ duplicación/consistencia
flexibilidad ↔ complejidad accidental
time-to-market ↔ reversibilidad/deuda
```

## 3. Escenarios de atributos de calidad

**[ACADEMIC]** El formato SEI usa seis partes: fuente del estímulo, estímulo, entorno, artefacto afectado, respuesta y medida de respuesta. Esto vuelve falsable una cualidad y conecta requisito con táctica.

```text
Fuente:
Estímulo:
Entorno:
Artefacto:
Respuesta:
Medida:
```

Ejemplo:

```text
Fuente: cliente autenticado
Estímulo: envía request válida
Entorno: pico 2× con una instancia caída
Artefacto: API de lectura y su dependencia primaria
Respuesta: responde dato consistente o error explícito sin retry storm
Medida: ≥99,9% dentro de 250 ms; 0 respuestas corruptas; queue age <100 ms
```

### 3.1 Catálogo de cualidades

ISO/IEC 25010:2023 proporciona nueve características de calidad como checklist; no sustituye métricas específicas. Añadir según sistema:

- functional suitability;
- performance efficiency;
- compatibility;
- interaction capability;
- reliability;
- security;
- maintainability;
- flexibility;
- safety.

Cruzar además privacidad, observabilidad, operabilidad, deployability, auditabilidad, costo y sostenibilidad cuando sean drivers reales.

### 3.2 Tipos de escenario obligatorios

- uso normal y pico;
- cambio de requisito o dependencia;
- fault/failure y recovery;
- ataque/abuso;
- deploy/rollback/mixed version;
- operación humana e incidente;
- crecimiento de datos/equipo/tráfico;
- retiro/migración/exportación.

## 4. Modelo de dominio y lenguaje

### 4.1 Antes de dividir software

Modelar:

- conceptos y definiciones;
- invariantes;
- comandos, eventos y estados;
- lifecycle y outcomes terminales;
- autoridad de cada dato;
- límites transaccionales;
- actores y políticas;
- términos ambiguos entre áreas.

Un nombre idéntico puede representar modelos distintos. Un límite debe explicarse por lenguaje, invariantes y ritmo de cambio, no por tabla o pantalla.

### 4.2 Bounded context como herramienta, no dogma

Dentro de un contexto, modelo y vocabulario son coherentes. Entre contextos, definir traducción y ownership. No asumir:

```text
bounded context = microservice = repositorio = equipo = base de datos
```

Pueden coincidir si los drivers lo justifican; no es una identidad universal.

### 4.3 Invariantes y agregados

Una boundary transaccional protege invariantes que deben cumplirse juntas. Mantenerla pequeña sin romper reglas. Si una regla cruza límites distribuidos, elegir explícitamente entre:

- coordinación/consenso;
- reserva/escrow;
- saga/compensación;
- reconciliación eventual;
- rediseño del invariante.

“Eventually consistent” no explica qué estados intermedios son válidos ni cuándo convergen.

## 5. Modularidad, abstracción y dependencias

**[ACADEMIC]** MIT 6.102 conecta specifications, abstract data types, representation invariants, interfaces, subtyping, concurrencia y message passing. Un módulo fuerte oculta decisiones y expone un contrato pequeño que puede probarse independientemente.

### 5.1 Criterios

| Propiedad | Pregunta |
|---|---|
| cohesión | ¿las responsabilidades cambian por la misma razón? |
| acoplamiento | ¿qué conocimiento atraviesa el límite? |
| encapsulación | ¿la representación puede cambiar sin consumidores? |
| estabilidad | ¿dependencias apuntan hacia contratos más estables? |
| testability | ¿se puede verificar sin montar todo el sistema? |
| replaceability | ¿el costo de sustitución coincide con el riesgo? |

### 5.2 Contrato de módulo

```text
responsabilidad y no-responsabilidades
API y tipos
pre/postcondiciones
invariantes
errores y lifecycle
concurrencia/ownership
datos persistentes
dependencias permitidas
observabilidad
compatibilidad/versionado
```

### 5.3 Dependency rules

- dominio no depende de detalles de transporte o storage por comodidad;
- adapters traducen contratos externos a internos;
- ciclos necesitan rediseño o una interfaz estable compartida;
- global mutable state amplifica acoplamiento temporal;
- abstraer sólo una variación real o un boundary que necesita test/control;
- no crear interfaces de una sola implementación por ceremonia.

## 6. Estilos y topologías

### 6.1 Modular monolith

Buen default cuando un equipo necesita transacciones simples, cambios rápidos y operación contenida. Exigir límites comprobables dentro del proceso; de otro modo es un monolito accidental.

Fortalezas:

- llamadas y transacciones locales;
- deploy/observabilidad más simples;
- refactoring atómico;
- menor costo de plataforma.

Riesgos:

- coupling invisible;
- build/deploy creciente;
- ownership difuso;
- failure domain grande.

### 6.2 Layered

Separa presentación, aplicación, dominio e infraestructura. Evitar que cada cambio atraviese todas las capas o que “service” se vuelva un módulo sin cohesión. Las dependencias deben obedecer reglas, no sólo carpetas.

### 6.3 Ports and adapters

El núcleo define puertos; adapters implementan transporte, persistencia y proveedores. Útil para testabilidad y reemplazo de dependencias. No implica que toda función requiera interfaz, DTO y mapper separados.

### 6.4 Microservices

Adoptar sólo si independencia de despliegue/escala/ownership supera:

- latencia y partial failure;
- contratos/version skew;
- transacciones y consistencia;
- observabilidad distribuida;
- seguridad entre servicios;
- plataforma/on-call;
- duplicación y costo cognitivo.

Un servicio necesita dueño, SLO, datos, API/eventos, deploy, dashboards, runbook y retiro. “Pequeño” no es un criterio suficiente.

### 6.5 Event-driven

Útil para desacoplar tiempo y fan-out, pero introduce delivery, orden, replay, schema evolution, lag y efectos ambiguos. Un evento debe declarar:

```text
hecho ocurrido vs intención
owner y schema
key/partition/order
timestamp/epoch/version
delivery y dedupe
retention/replay
PII y ACL
```

### 6.6 CQRS y event sourcing

CQRS separa modelos de lectura/escritura cuando sus necesidades divergen materialmente. Event sourcing hace del log de eventos la fuente para reconstruir estado. Ninguno es prerequisito del otro. Antes de adoptar, probar:

- evolución/corrección de eventos;
- rebuild y tiempo de replay;
- privacidad/delete;
- snapshots;
- queries/auditoría;
- side effects;
- operación y capacitación.

### 6.7 Pipeline, batch y streaming

Elegir por boundedness, freshness, volumen, orden y recuperación. Un pipeline no es sólo flechas: cada etapa declara input/output, checkpoint, retry, idempotencia, backpressure y dueño.

## 7. Comunicación y contratos

### 7.1 Local versus remoto

Una llamada remota puede fallar antes o después de causar efecto. Diseñar:

- deadline total;
- idempotencia/deduplicación;
- partial response;
- retry budget;
- backpressure;
- circuit/fallback sólo si preserva semántica;
- reconciliación de outcome ambiguo.

No ocultar una red detrás de una interfaz que parece local sin expresar latencia y fallos.

### 7.2 Sync versus async

| Sincrónico | Asincrónico |
|---|---|
| feedback inmediato | desacopla tiempo y absorbe bursts acotados |
| dependencia temporal visible | lag/replay/order se vuelven parte del contrato |
| cadena larga amplifica latencia/fallo | cola puede esconder saturación |

Elegir según workflow, no según tecnología favorita. La asincronía no elimina dependencia; cambia cómo se observa y recupera.

### 7.3 Compatibilidad

- reader/writer coexistentes;
- additive antes que breaking;
- defaults y unknown fields;
- consumer-driven assumptions visibles;
- deprecation con telemetría;
- schema y rollout coordinados;
- rollback que soporte datos nuevos.

Hyrum’s Law advierte que todo comportamiento observable puede adquirir consumidores; reducir superficie accidental y medir uso antes de retirar.

## 8. Datos y ownership

### 8.1 Una autoridad por hecho

Para cada dato:

```text
owner/source of truth
writers autorizados
invariante
freshness
retención/delete
replicas/caches/views
reconciliation
export/migration
```

Shared database puede simplificar transacciones y complicar ownership/evolución. Database-per-service puede aislar y complicar joins/consistencia. Elegir por drivers y probar failure paths.

### 8.2 Transacción, evento y side effect

No prometer atomicidad que el sistema no posee. Para DB + broker/API:

- transactional outbox/inbox;
- idempotent consumer;
- saga/compensation;
- reconciliation job;
- operation ID y audit trail.

Exactly-once de un componente no crea exactly-once end-to-end.

### 8.3 Cache

Toda caché declara key completa, autoridad, TTL/freshness, invalidation, consistency, negative caching, stampede policy, tenant/security scope y comportamiento durante partición. Cachear es introducir una réplica.

## 9. Tácticas por quality driver

### 9.1 Performance

- reducir trabajo y bytes;
- cache/batching con freshness y espera acotadas;
- concurrencia/parallelism con ownership;
- locality y índices;
- admission/backpressure;
- presupuestos por fase.

Medir p50/p95/p99/p99.9, goodput, queue age y recursos; no sólo promedio.

### 9.2 Availability y resilience

- fault isolation y bulkheads;
- redundancy sin shared fate oculto;
- health semántico;
- timeout/cancellation;
- bounded retry con jitter;
- graceful degradation;
- failover con fencing;
- backup/restore probado.

Redundancia sin detección, capacidad ni recovery probado sólo duplica componentes.

#### Práctica pública AWS — estabilidad estática y aislamiento celular

**[PROD]** Amazon define *static stability* como conservar el servicio útil cuando una dependencia/control plane se deteriora, sin requerir cambios inmediatos para recuperarse. La capacidad de supervivencia, configuración y estado necesarios ya existen antes del incidente; el data plane puede continuar con estado previamente válido aunque no reciba actualizaciones.

**[IMPL] Pregunta de design review:** “Si esta zona, control plane, discovery, autoscaler o configuration service falla ahora, ¿qué acción nueva debe tener éxito para que el camino crítico siga funcionando?”. Cada respuesta es una dependencia de recovery que debe eliminarse, preprovisionarse o probarse bajo fallo conjunto.

Patrones combinables:

- **control/data plane:** el data plane conserva localmente el estado requerido y tiene un SLO superior; updates pueden quedar stale dentro de un límite explícito;
- **headroom preexistente:** capacidad restante soporta la pérdida del failure domain sin crear instancias, credenciales o registros durante el incidente;
- **cells:** réplicas completas, de tamaño acotado e independientes atienden particiones de tenants/tráfico; routing y rebalancing no introducen un shared fate mayor;
- **shuffle sharding:** cada tenant usa un subconjunto pequeño de workers; combinaciones distintas reducen coimpacto sin dedicar una flota completa por tenant;
- **rollout por fault boundary:** one-box/canary y expansión por célula/zona impiden que un mismo cambio alcance todos los dominios simultáneamente;
- **evitar fallback bimodal:** una ruta rara y poco ejercitada que aumenta trabajo durante el fallo suele amplificarlo; preferir comportamiento continuo, acotado y probado.

Si existen `N` workers y cada tenant recibe `k`, hay `C(N,k)` combinaciones posibles, pero la matemática no prueba aislamiento: importan distribución, hotspots, capacidad del shard, tenants grandes, rebalancing y dependencias compartidas.

**Gate:** simular pérdida de cell/zona/control plane; mantener offered load; demostrar goodput dentro del SLO, recursos acotados y blast radius esperado; probar routing stale, cell llena, tenant desproporcionado, deploy defectuoso, rebalancing y recuperación sin dependencia circular. Registrar costo/headroom y una alternativa más simple: cells no son default para sistemas cuyo riesgo no justifica su complejidad.

Fuentes: [Static stability](https://aws.amazon.com/builders-library/static-stability-using-availability-zones/), [shuffle sharding](https://aws.amazon.com/builders-library/workload-isolation-using-shuffle-sharding/), [avoiding fallback](https://aws.amazon.com/builders-library/avoiding-fallback-in-distributed-systems/) y [cell-based architecture](https://docs.aws.amazon.com/solutions/cell-based-architecture-on-aws/).

### 9.3 Modifiability y deployability

- boundaries por decisiones que cambian juntas;
- stable contracts;
- branch by abstraction/expand-contract;
- feature flag con owner/expiry;
- backward/forward compatibility;
- small reversible rollout;
- automated conformance tests.

### 9.4 Security y privacy

- trust boundaries y least privilege;
- authn separada de authz;
- minimización y clasificación de datos;
- secrets/keys fuera de artefactos;
- input/output limits;
- isolation/multitenancy;
- audit e incident response;
- supply-chain provenance.

Security no es un proxy añadido al final; cambia límites, identidad, datos y operación.

### 9.5 Observability y operability

- señales derivadas de SLO/failure modes;
- correlation/operation IDs;
- logs estructurados sin secretos;
- métricas bounded-cardinality;
- traces con sampling consciente;
- admin/control plane separado;
- runbooks y safe diagnostics.

### 9.6 Testability

- contratos deterministas;
- dependency control;
- clocks/IDs/randomness inyectables donde sea necesario;
- oráculos e invariantes;
- test environments representativos;
- fault injection;
- arquitectura comprobable por código.

## 10. Estimación de capacidad

### 10.1 Back-of-envelope con unidades

```text
requests/s pico y sostenido
bytes/request y response
fan-out y amplification
CPU ms/request
working set y bytes/objeto
read/write ratio
retención × ingestión
replication/encoding overhead
concurrencia ≈ throughput × latency
headroom y failure capacity
```

Toda cifra lleva fuente, fecha, percentil y rango de incertidumbre. La estimación detecta imposibles; el load test valida el diseño.

### 10.2 Carga y crecimiento

- promedio no dimensiona pico;
- compresión nominal no equivale a memoria pico;
- réplicas consumen red/storage/operación;
- maintenance compite con serving;
- failover necesita capacidad disponible durante el fallo;
- cardinalidad/skew domina particiones;
- usuarios, equipos y deploys también escalan.

## 11. Diseñar fallos antes del happy path final

### 11.1 Failure inventory

Por dependency y fase:

```text
no responde / lento / respuesta inválida
respuesta perdida tras efecto
duplicado / reorder / stale
partial write / disk full / corruption
process/node/zone/region/control-plane loss
credential/cert expiry
quota/cost exhaustion
bad deploy/config/schema
operator error / compromised actor
```

### 11.2 Estados explícitos

No colapsar `PENDING`, `SUCCEEDED`, `FAILED`, `CANCELLED`, `EXPIRED` y `UNKNOWN`. `UNKNOWN` exige reconciliación, no retry ciego.

### 11.3 RTO, RPO y evidence

Declarar recovery time y recovery point por flujo/dato, no por sistema abstracto. Probar restore funcional, credenciales, schema, dependencias, capacidad y cut de datos. Backup sin restore drill es una hipótesis.

## 12. Documentación arquitectónica

### 12.1 ISO 42010: descripción no es arquitectura

**[SPEC]** ISO/IEC/IEEE 42010:2022 distingue la arquitectura de la arquitectura description y estructura la descripción alrededor de entity of interest, stakeholders, concerns, viewpoints, views y model kinds.

**[IMPL]** Cada vista debe indicar:

- audiencia/concerns;
- alcance y nivel;
- elementos/relaciones;
- leyenda y semántica;
- fuente de verdad y owner;
- fecha/versión;
- links a decisiones/tests;
- inconsistencias conocidas.

### 12.2 C4

El C4 oficial usa jerarquía system → container → component → code y agrega landscape, dynamic y deployment views. “Container” en C4 es una unidad ejecutable/data store, no necesariamente Docker.

Reglas:

- empezar por context;
- mostrar personas/sistemas externos;
- cada relación tiene verbo, dirección y tecnología relevante;
- un diagrama por historia y audiencia;
- no mezclar niveles;
- component/code sólo donde aportan una decisión.

### 12.3 arc42 compacto

Las doce secciones oficiales proporcionan un receptor práctico:

1. introduction/goals;
2. constraints;
3. context/scope;
4. solution strategy;
5. building blocks;
6. runtime;
7. deployment;
8. crosscutting concepts;
9. decisions;
10. quality requirements;
11. risks/technical debt;
12. glossary.

No completar por ceremonia. Mantener lo que ayuda a una decisión, implementación, operación, auditoría o onboarding.

### 12.4 Vista de runtime

Para cada escenario crítico mostrar:

- actor y trigger;
- orden/correlation;
- componentes y boundaries;
- sync/async;
- transaction/durable points;
- timeout/retry/cancel;
- errores y compensación;
- señales operativas.

### 12.5 Vista de deployment

Mostrar artefacto, runtime, node/zone/region, red/trust, storage, secrets, autoscaling, failure domains, observabilidad y rollout. Un diagrama lógico no demuestra aislamiento físico.

## 13. Architecture Decision Records

**[PROD]** Michael Nygard propuso ADRs cortos en el repositorio para conservar fuerzas, decisión, estado y consecuencias, incluyendo las negativas; las decisiones reemplazadas se conservan como superseded.

Plantilla operativa:

```markdown
# ADR-NNN: decisión

Fecha/owner/status:

## Contexto y fuerzas
Drivers, restricciones, evidencia, incertidumbre.

## Opciones
Alternativas reales y opción de no cambiar.

## Decisión
Qué se hará, alcance y fecha/trigger.

## Consecuencias
Positivas, negativas, neutras; costo operativo y migración.

## Validación
Escenarios, tests, métricas y fecha de revisión.

## Reversión/sucesión
Rollback, exit y ADR que reemplaza.
```

Una ADR no es una minuta ni propaganda posterior. Una decisión por registro; no editar la historia para fingir que siempre se supo la respuesta.

## 14. Proceso de system design

### Fase A — discovery

1. Clarificar usuario, negocio y daño.
2. Inventariar contexto y dependencias.
3. Cuantificar carga, datos y SLO.
4. Priorizar escenarios de calidad.
5. Registrar incertidumbres.

### Fase B — alternativas

1. Diseñar opción simple/base.
2. Proponer máximo dos o tres alternativas materiales.
3. Evaluar cada driver, no una lista genérica de pros/cons.
4. Identificar decisiones reversibles e irreversibles.
5. Crear probes para unknowns dominantes.

### Fase C — skeleton

1. Context/container/runtime/deployment views.
2. Interfaces y schemas mínimos.
3. Walking skeleton end-to-end.
4. Auth, observabilidad, deploy y rollback desde el inicio.
5. Prueba de capacidad/fallo de mayor riesgo.

### Fase D — implementación incremental

1. Vertical slices pequeños.
2. Tests de contrato/invariante.
3. Migraciones expand-contract.
4. Canary/feature flag con expiración.
5. Documentar divergencias reales.

### Fase E — validación

1. Escenarios funcionales y de calidad.
2. Load/soak/fault/security/recovery.
3. Architecture conformance.
4. Operability/game day.
5. Revisión de objetivos y costos.

## 15. Evaluación de trade-offs

### 15.1 Mini-ATAM

**[ACADEMIC]** ATAM del SEI evalúa decisiones contra quality attribute requirements y expone sensitivity points, tradeoff points, riesgos y non-risks.

Procedimiento compacto:

1. presentar objetivos y stakeholders;
2. presentar arquitectura y decisiones;
3. priorizar quality scenarios;
4. recorrer escenario por componentes/mecanismos;
5. identificar supuesto, sensibilidad y trade-off;
6. registrar riesgo, evidencia y acción;
7. repetir después de probes o cambios.

Tabla:

| Escenario | Decisión/táctica | Evidencia | Sensibilidad | Trade-off | Riesgo/acción |
|---|---|---|---|---|---|

### 15.2 Matriz de alternativa

No sumar scores arbitrarios sin pesos ni incertidumbre. Usar:

```text
driver + umbral
evidencia por alternativa
consecuencia negativa
confidence
experimento que cambia la decisión
```

## 16. Build versus buy

Evaluar:

- strategic differentiation;
- fit funcional y de calidad;
- integración/operación;
- seguridad/compliance/data;
- lock-in y export;
- roadmap/EOL;
- licencia/support;
- costo total y skills;
- failure/exit plan.

Un managed service transfiere tareas, no responsabilidad por contrato, datos o continuidad. Un componente propio crea mantenimiento permanente.

## 17. Evolución, refactoring y modernización

**[PROD]** Fowler define refactoring como transformaciones pequeñas que preservan comportamiento externo. La arquitectura también debe evolucionar mediante pasos comprobables.

### 17.1 Secuencia segura

```text
characterization tests
→ observabilidad
→ seam/adapter
→ cambio pequeño
→ dual/read-shadow si aplica
→ compare/reconcile
→ migrate/backfill
→ cutover gradual
→ retirar camino viejo
```

### 17.2 Patrones de migración

- branch by abstraction;
- expand/migrate/contract;
- strangler con routing claro;
- parallel run/shadow traffic sin side effects dobles;
- anti-corruption layer;
- change-data capture con cut consistente;
- façade antes de extracción.

### 17.3 Deuda técnica

Registrar principal, interés, riesgo, owner, evidencia y trigger. No llamar deuda a todo código desagradable; deuda es una decisión/circunstancia que encarece cambios o riesgo futuro. Priorizar por impacto, frecuencia de cambio y blast radius.

## 18. Organización y ownership

La arquitectura crea interfaces humanas. Para cada componente/servicio:

- owner y reviewers;
- consumers y SLO;
- on-call/escalation;
- cambios permitidos;
- dependency policy;
- lifecycle y deprecation;
- documentación y runbook.

No fragmentar el sistema más fino que la capacidad de operar y coordinar. La independencia técnica sin autoridad de producto/datos/deploy es parcial.

## 19. Tests y fitness arquitectónico

### 19.1 Estáticos

- dependency direction;
- forbidden imports/cycles;
- public API surface;
- schema/lint/policy;
- ownership y layering;
- artifact/dependency provenance.

### 19.2 Dinámicos

- contract/compatibility;
- architecture scenario end-to-end;
- performance budget;
- failover/recovery;
- security boundaries;
- isolation/tenant;
- migration/mixed version;
- observability assertions.

### 19.3 Conformance

Comparar arquitectura declarada con build graph, runtime topology, network policy, deployment manifests y traces. Un diagrama desactualizado no es evidencia; una divergencia puede exigir corregir código o decisión.

## 20. Observabilidad de decisiones

Cada driver necesita indicador:

| Driver | Evidencia |
|---|---|
| performance | latency distribution, goodput, saturation |
| availability | SLI/SLO, failure coverage, MTTR |
| modifiability | lead time, blast radius, dependency/change graph |
| deployability | frequency, change failure, rollback time |
| security | control effectiveness, exposure, detection/recovery |
| cost | unit economics y headroom |

No convertir métricas en targets ciegos. Revisar si miden la cualidad o un proxy manipulable.

## 21. Anti-patrones

- architecture astronautics sin walking skeleton;
- distributed monolith;
- shared database sin ownership;
- event soup sin schemas ni consumidores conocidos;
- sync chain profunda;
- queue infinita como “resilience”;
- abstraction por cada clase;
- framework como dominio;
- diagramas sin relaciones/verbos/owners;
- ADR que sólo documenta la opción elegida;
- NFR sin medida;
- HA sin failure domain/capacity test;
- cache sin invalidation/security scope;
- platform team que obliga una solución no justificada;
- reescritura total sin cutover/reconcile;
- premature multi-region/microservices.

## 22. Gate de design review

- [ ] Problema, alcance, stakeholders y constraints explícitos.
- [ ] Quality scenarios priorizados y medibles.
- [ ] Carga/datos/capacidad con unidades y supuestos.
- [ ] Dominio, invariantes y source of truth.
- [ ] Context/container/runtime/deployment suficientes.
- [ ] Alternativas reales y trade-offs.
- [ ] Security/privacy/threat boundaries.
- [ ] Failure model, RTO/RPO y outcome ambiguo.
- [ ] Compatibility, rollout, migration y rollback.
- [ ] Observabilidad, on-call y costos.
- [ ] Tests de arquitectura y aceptación.
- [ ] ADRs para decisiones significativas.
- [ ] Riesgos/unknowns con probes y owner.
- [ ] Exit/deprecation cuando existe dependencia duradera.

## 23. Contrato para Codex

```text
Actuá como arquitecto implementador, no como selector de tecnologías.

Entrada:
- objetivo, usuarios y daño;
- alcance/constraints;
- carga, datos, SLO/RTO/RPO;
- equipo, plazo, presupuesto y stack existente;
- dependencias y decisiones previas.

Entregar:
1. preguntas/supuestos bloqueantes y no bloqueantes;
2. quality scenarios medibles;
3. modelo de dominio, invariantes y ownership;
4. baseline simple + 1–2 alternativas;
5. trade-offs por driver y recomendación condicionada;
6. C4 context/container + runtime/deployment críticos;
7. APIs/events/data contracts y failure semantics;
8. capacity estimate con unidades;
9. security, observability, rollout, rollback y recovery;
10. ADRs, probes, tests y plan incremental.

No inventar requisitos.
No afirmar escalabilidad, seguridad o alta disponibilidad sin escenario y evidencia.
No elegir microservices/event sourcing/multi-region por defecto.
No ocultar partial failure ni transacción distribuida.
No producir diagrama sin texto verificable ni ownership.
```

## 24. Pack mínimo por proyecto

Mantener pocos artefactos vivos:

```text
README: propósito, ejecución, owners y links
SYSTEM_DESIGN: drivers, vistas, contratos, quality/failure scenarios
docs/adr/: decisiones significativas
API/schema files: contratos ejecutables
tests/architecture/: conformance y quality gates
runbooks/: operación, incidente, rollback y restore
```

Para sistemas pequeños, un único documento de system design puede incluir todo arc42/C4 compacto. Separar sólo cuando audiencia, ownership o tamaño lo justifiquen.

## 25. Capstones de validación

### A. Modular monolith evolutivo

Diseñar billing/orders/catalog con límites internos, una transacción crítica y outbox. Probar dependency rules, invariantes, migration y extracción simulada sin extraer prematuramente.

### B. Servicio de alta disponibilidad

Definir SLO, capacity y fault domains. Ejecutar node/zone/dependency failure, overload, cert expiry y rollback. Demostrar goodput y recovery, no sólo réplicas.

### C. Pipeline de datos/IA

Modelar ingestión, validación, lineage, training, registry, serving y feedback; separar calidad del dato/modelo/sistema. Probar replay, schema evolution, privacy y rollback.

### D. Order-book platform

Separar feed handler/hot state, durable raw log, analytics, API y UI. Preservar epoch/sequence/freshness y mantener IA fuera del hot path salvo evidencia. Probar gap, replay, backpressure y venue disconnect.

## 26. Inventario de autoridades

| Fuente | Cobertura usada | Límite |
|---|---|---|
| SEI Software Architecture | quality-driven architecture, documentation, evaluation, ATAM, risk | método; adaptar al contexto |
| MIT 6.102 Spring 2025, 19 readings | specs, ADT, rep invariants, interfaces, concurrency, message passing | construcción modular, no arquitectura enterprise completa |
| ISO/IEC/IEEE 42010:2022 | arquitectura description, stakeholders, concerns, views/viewpoints | no prescribe arquitectura concreta |
| ISO/IEC 25010:2023 | checklist/modelo de calidad de producto | categorías no sustituyen escenarios/umbrales |
| Google SWE book | tiempo, escala, trade-offs, design docs, cambios, testing/dependencies | experiencia Google, no ley universal |
| C4 | abstracciones y diagramas jerárquicos | comunicación, no proceso de decisión |
| arc42 | 12 receptores documentales y calidad | template adaptable, no checklist ceremonial |
| Nygard ADR | contexto, decisión, status, consecuencias | registro, no prueba de acierto |
| Fowler Refactoring | cambio pequeño que preserva comportamiento | requiere tests y estrategia de migración sistémica |

## 27. Fuentes oficiales verificadas

- SEI Software Architecture: <https://www.sei.cmu.edu/our-work/software-architecture/>
- SEI ATAM: <https://www.sei.cmu.edu/documents/629/2000_005_001_13706.pdf>
- SEI Quality Attribute Scenarios: <https://www.sei.cmu.edu/documents/704/2003_005_001_14213.pdf>
- MIT 6.102 Spring 2025: <https://web.mit.edu/6.102/www/sp25/>
- ISO/IEC/IEEE 42010:2022: <https://www.iso.org/standard/74393.html>
- ISO/IEC 25010:2023: <https://www.iso.org/standard/78176.html>
- *Software Engineering at Google*: <https://abseil.io/resources/swe-book>
- Google design docs chapter: <https://abseil.io/resources/swe-book/html/ch10.html>
- C4 official: <https://c4model.com/>
- arc42 overview: <https://arc42.org/overview/>
- arc42 documentation: <https://docs.arc42.org/home/>
- Michael Nygard, ADR original: <https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions>
- Martin Fowler, Refactoring: <https://refactoring.com/>

## 28. Cruces iniciales

| Manual | Autoridad de frontera |
|---|---|
| `ALGORITHMS_DATA_STRUCTURES_PROBLEM_SOLVING.md` | modelado, corrección y costo abstracto |
| `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md` | costo físico, concurrencia y medición |
| `SOFTWARE_BACKEND_API_ENGINEERING.md` | contratos de servicio y construcción backend |
| `NETWORKING_DISTRIBUTED_STREAMING.md` | comunicación, partial failure, consenso y streams |
| `DATABASE_STORAGE_INTERNALS.md` | persistencia, isolation, recovery e índices |
| `GPU_ACCELERATED_COMPUTING.md` | arquitectura y operación de aceleradores |
| `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md` | threat model, plataforma, SLO e incidentes |
| `MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md` | semántica de mercados/exchanges |
| `FRONTEND_PRODUCT_ENGINEERING_UX.md` | interacción, navegador y experiencia |
| `AI_ENGINEERING_MASTER_MAP.md` | modelos, datos, agentes y evaluación de IA |
| `NATIVE_MOBILE_DESKTOP_ENGINEERING.md` | lifecycle del OS, activación, offline, permisos y distribución instalada |

Cruces recíprocos cerrados el 2026-08-21: cada receptor incorpora sólo la obligación que modifica una decisión y conserva autoridad sobre su dominio. Este manual gobierna drivers, boundaries, views y trade-offs; no reemplaza semántica algorítmica, física, distribuida, de datos, seguridad, producto, mercados o IA.

Cruce con `DATA_ENGINEERING_ANALYTICS.md` cerrado el 2026-08-21: arquitectura decide boundaries, ownership, quality/failure scenarios y build/buy; datos concreta source-of-truth, grain, time, lineage, freshness, quality, retention y consumer contracts. Lake/warehouse/mesh/lakehouse no son objetivos: cada topology se justifica por escenarios y conserva backfill, delete, cost, recovery y exit plan.

Cruce con `TOOLCHAINS_BUILDS_PACKAGING_FFI.md` cerrado el 2026-08-21: artifact/package/plugin/FFI boundaries, supported targets y compatibility son decisiones arquitectónicas; toolchains demuestra que se construyen, instalan, versionan y revierten. Build graph/manifest/lock/provenance deben corresponder a C4/deployment/ADR; shared library, native extension o package manager no son detalles si condicionan lifecycle y ABI.

Cruce con `NATIVE_MOBILE_DESKTOP_ENGINEERING.md` cerrado el 2026-08-21: proceso, ventana/escena, lifecycle, activación, background, permisos, offline source-of-truth y update son drivers arquitectónicos. La arquitectura fija boundaries y qualities; el manual nativo prueba su comportamiento bajo process death, conectividad variable, resource pressure y reglas de distribución. Compartir UI/core sólo se decide tras una vertical slice medida y un escape hatch de plataforma explícito.

## 29. Puerta de cierre

- [x] Autoridades académicas, normativas y operativas.
- [x] Drivers, quality scenarios, dominio y modularidad.
- [x] Estilos, contratos, datos y tácticas.
- [x] Capacity, failure model, views y ADR.
- [x] Proceso, ATAM, evolución, tests y Codex contract.
- [x] Fuentes e inventario con límites.
- [x] Auditoría Markdown/referencias: fences, encabezados, caracteres y referencias internas.
- [x] Cruces recíprocos y arbitraje de solapamientos.
- [x] Contrato ejecutable de proyecto, ADR, quality scenarios y capstones disponible en `ENGINEERING_EXECUTION_PLAYBOOK.md`.

El conocimiento general v1 y su contrato de planificación quedan auditados. No declarar una arquitectura implementada ni una suite `passed` hasta que exista evidencia del proyecto real.
