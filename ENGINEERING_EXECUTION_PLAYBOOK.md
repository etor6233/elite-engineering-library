# Engineering Execution Playbook — de conocimiento a software verificable

> **Estado:** v1.3 ejecutable auditada el 2026-09-05; contrato JSON 1.1.0 con aseguramiento de implementación, resume cursor NEW/EXISTING, checkpoint log, validadores, pruebas y tres capstones de planificación incluidos.
> **Propósito:** cargar el mínimo corpus necesario, convertir una necesidad en contratos y decisiones, implementar verticalmente y aceptar sólo resultados respaldados por evidencia.
> **Límite:** el kit valida estructura, trazabilidad y artifacts declarados; no sustituye la ejecución real de tests, benchmarks, security reviews, restores ni releases de cada proyecto.

## 1. Resultado

La biblioteca se usa así:

```text
problema real
→ clasificar dominio y riesgo
→ detectar NEW/EXISTING y checkpointar baseline/delta
→ cargar mapa + manual autoridad + fronteras materiales
→ completar project contract
→ validar plan
→ construir vertical slice
→ ejecutar correctness/security/performance/recovery gates
→ capturar evidence por digest
→ validar evidence
→ staged release
→ observar, aprender de errores y actualizar contratos/tests
→ guardar siguiente acción y contexto mínimo reanudable
```

No cargar todos los manuales en cada prompt. La amplitud está en la biblioteca; la precisión operativa viene de seleccionar sólo la autoridad y las fronteras que cambian una decisión.

## 2. Artefactos ejecutables

| Archivo | Uso |
|---|---|
| `engineering_execution_kit/project.schema.json` | contrato formal JSON Schema 2020-12 |
| `engineering_execution_kit/project.example.json` | ejemplo completo, deliberadamente `planned` |
| `engineering_execution_kit/validate_project.py` | gate sin dependencias: referencias, estados, trazabilidad, paths y SHA-256 |
| `engineering_execution_kit/test_validate_project.py` | regresiones del validador |
| `engineering_execution_kit/capstones/ai_production_service.json` | IA: model/data/eval/runtime/canary |
| `engineering_execution_kit/capstones/low_latency_order_book.json` | sistemas: sequence/gap/replay/tails |
| `engineering_execution_kit/capstones/offline_native_client.json` | producto: lifecycle/offline/sync/update |
| `engineering_execution_kit/execution_method_lock.json` | fuentes exactas y límites del método del agente |
| `engineering_execution_kit/execution_state.template.json` | template fail-closed para `PROJECT_EXECUTION_STATE.json` |
| `engineering_execution_kit/validate_execution_state.py` | modo, baseline/delta, tareas, referencias, evidencia y cierre |
| `engineering_execution_kit/checkpoint_execution_state.py` | checkpoint append-only, SHA-256 encadenado y durable |
| `engineering_execution_kit/EXECUTION_STATE_PROTOCOL.md` | flujo nuevo/existente, fallos, reanudación y release |

### 2.1 Comandos

Desde el directorio raíz:

```powershell
python .\engineering_execution_kit\validate_project.py `
  .\engineering_execution_kit\project.example.json --level plan

python .\engineering_execution_kit\validate_project.py `
  .\engineering_execution_kit\capstones --level plan

python -m unittest discover `
  -s .\engineering_execution_kit -p 'test_*.py' -v
```

Para iniciar o reanudar sin depender de memoria conversacional:

```powershell
python .\engineering_execution_kit\validate_execution_state.py `
  .\PROJECT_EXECUTION_STATE.json --project-root . `
  --events .\PROJECT_EXECUTION_EVENTS.jsonl --level resume
```

Cuando un proyecto tenga resultados reales:

```powershell
python .\engineering_execution_kit\validate_project.py `
  C:\ruta\al\proyecto.json --level evidence
```

`plan` verifica que el trabajo esté definido y trazable. `evidence` exige todos los tests/benchmarks en `passed`, evidencia `verified`, archivos existentes, digests SHA-256 correctos y `implementation_assurance=proven` sin bloqueantes.

## 3. Selección mínima del corpus

### 3.1 Router

1. Consultar `AI_ENGINEERING_MASTER_MAP.md` si el valor central depende de modelo, datos de entrenamiento/evaluación, retrieval o agente.
2. Consultar `SYSTEMS_ENGINEERING_MASTER_MAP.md` para software, sistemas, datos, seguridad, producto o plataforma.
3. Elegir un único manual autoridad por decisión.
4. Añadir sólo manuales frontera que el flujo realmente atraviesa.
5. En el manifest, guardar esos nombres en `codex.authority_docs`.

### 3.2 Paquetes frecuentes

| Trabajo | Autoridad mínima |
|---|---|
| API transaccional | Arquitectura + Backend + DB; Seguridad si hay usuarios/datos |
| servicio distribuido | Arquitectura + Backend + Redes + Seguridad/SRE |
| pipeline analítico | Data Engineering + source/storage/stream correspondiente |
| optimización de CPU | Sistemas + Algoritmos + manual del dominio |
| GPU/IA serving | manual IA + Inference + GPU + Sistemas + Seguridad |
| agente con tools | Agentic AI + Backend + Seguridad + dominio de cada tool |
| control de seguridad para agentes | AI Security + Agentic AI + Arquitectura + Backend + Seguridad/SRE |
| RAG | NLP/RAG + ML Production + Backend + DB + Seguridad |
| frontend web | Frontend + Backend + Seguridad + dominio visible |
| app mobile/desktop | Nativo + Frontend + Backend/DB/Seguridad según flujo |
| order book monitor | Mercados + Redes + Sistemas; Backend/Data/Frontend si existen |

### 3.3 Regla de arbitraje

Si dos manuales parecen ordenar cosas distintas:

```text
dominio define verdad e invariantes
arquitectura define boundaries y quality trade-offs
plataforma define semántica disponible
sistemas define costo físico medido
seguridad define acciones permitidas y evidencia
producto define task/user outcome
operación define SLO, recovery y sostenibilidad
```

No resolver una tensión borrando una capa. Escribir el supuesto y convertirlo en ADR, test o benchmark.

## 4. Dos gates distintos

Antes de ambos gates, `implementation_assurance` declara procedencia, tier no inferior al riesgo observado, sources del método y nueve dimensiones. Un proyecto high/critical debe cubrir corrección, integración, seguridad, fallos, rendimiento, operación, recovery y release; UX/accesibilidad se exige cuando hay una persona en el journey o se excluye con razón verificable. Código `ADAPTED|VERBATIM|MIXED` requiere source locks. En `evidence`, artefacto y release usan evidencias verificadas distintas para impedir que un build local se presente como paquete promovido.

### 4.1 Gate de plan

Responde:

- ¿qué problema y scope se autorizan?;
- ¿qué requisito tiene qué fuente?;
- ¿qué test acepta cada requisito?;
- ¿qué calidad se mide en qué escenario?;
- ¿qué decisiones y riesgos existen?;
- ¿qué contratos y compatibilidad se prometen?;
- ¿qué artifacts y evidencia deberán producirse?;
- ¿qué procedencia, riesgo y dimensiones demostrarán que la implementación es adecuada al claim?;

Pasar `plan` no significa que el software funcione.

### 4.2 Gate de evidencia

Responde:

- ¿se ejecutaron tests y benchmarks?;
- ¿están en `passed`?;
- ¿cada afirmación tiene artifact verificable?;
- ¿el archivo existe dentro del root permitido?;
- ¿el SHA-256 coincide?;
- ¿el ambiente/herramienta/fecha están declarados?;
- ¿la implementación está `proven`, sin bloqueantes, con source locks cuando corresponden y artefacto/release distintos?;

Pasar `evidence` tampoco prueba universos no ensayados. Prueba el contrato y el ambiente declarados.

## 5. Intake profesional

Antes de diseñar, capturar:

```text
objetivo y usuario/operador:
journey o workflow crítico:
resultado observable:
invariantes que nunca pueden romperse:
datos/activos/personas afectadas:
volumen, distribución, ráfagas y crecimiento:
latencia/availability/durability/quality:
fallos y adversarios plausibles:
plataformas, integraciones y constraints:
compatibilidad/migración/fecha:
equipo, ownership y operación:
fuera de alcance:
```

Preguntar sólo lo que cambia scope, riesgo o decisión. Lo descubrible en código, métricas, especificaciones o configuración se inspecciona.

### 5.1 Bloqueadores reales

Detener implementación si falta:

- autorización para una mutación externa/material;
- definición de una operación irreversible;
- fuente de verdad o identidad de dato;
- criterio de corrección en un dominio de alto riesgo;
- secreto/entorno que no puede sustituirse con fake;
- elección entre dos scopes materialmente distintos.

No detenerse por nombres, detalles reversibles o elecciones locales que una ADR puede registrar.

### 5.2 Gate de oportunidad para productos nuevos

Un plan técnico no demuestra que exista una empresa. Antes de una construcción comercial amplia, separar hechos, inferencias e hipótesis sobre:

```text
usuario y comprador económico:
problema/workflow actual y frecuencia:
costo, riesgo o ingreso afectado:
alternativas y presupuesto existente:
evento que vuelve urgente la compra:
ventaja inicial y razón para cambiar:
canal de llegada y ciclo de venta:
precio/costo/margen hipotéticos:
design partners y acceso a evidencia:
umbral explícito de continuar, pivotar o abandonar:
```

Orden de evidencia creciente:

1. fuentes públicas y comportamiento observable;
2. entrevistas con el perfil correcto y problema concreto;
3. acceso a workflow/dataset/sandbox representativo;
4. compromiso verificable —tiempo, datos, responsable, LOI o piloto—;
5. pago, uso y retención.

Comentarios favorables o tamaño de mercado no sustituyen compromiso. Tampoco construir un MVP valida por sí mismo disposición a pagar. Definir antes de contactar cuántas señales independientes, qué objeciones y qué umbral decidirán el siguiente paso para no reinterpretar toda respuesta como confirmación.

La planificación se abre por etapas: `descubrimiento → validación → vertical slice → piloto → productización → escala`. Cada etapa tiene hipótesis, evidencia requerida, presupuesto/tiempo acotado y kill criteria. Codex puede investigar, preparar outreach, prototipos y análisis; la demanda sólo queda validada por comportamiento autorizado de compradores reales.

## 6. Requirements y trazabilidad

Cada requirement debe ser:

- observable;
- singular;
- independiente de implementación salvo que ésta sea constraint;
- etiquetado con fuente;
- conectado a uno o más tests;
- priorizado `must/should/could`.

Forma:

```text
ID: REQ-01
Statement: condición observable
Source: [SPEC|ACADEMIC|CODE|PROD|MEASURED|IMPL|OPEN] referencia
Acceptance: TEST-01, TEST-04
```

Una expectativa sin test queda `[OPEN]`; un test sin requirement suele ser trabajo sin propósito o una obligación no documentada.

### 6.1 Gate para prácticas derivadas de empresas líderes

La reputación de una organización selecciona candidatos para inspección; no convierte una afirmación en evidencia. Toda práctica industrial incorporada debe registrar:

```text
organization + repository/document
revision/commit + path/section + observed_at
license/terms
evidence_type: CODE|SPEC|ACADEMIC|PROD|MEASURED
claim_supported
claim_not_supported
transfer_conditions
acceptance_tests/metrics
```

Reglas:

1. Preferir código+tests y especificaciones; usar engineering posts para decisiones/mediciones que declaran, no para completar internals ausentes.
2. No convertir una cifra de producción en promesa: conservar workload, hardware, baseline, percentil, fecha y método; sin ellos queda `[PROD]`, no `[MEASURED]` reusable.
3. No copiar defaults de otra empresa sin comparar constraints y baseline propios.
4. Distinguir patrón generalizable, detalle específico del sistema y componente omitido/no público.
5. Revisar licencia antes de reutilizar código; una síntesis técnica no autoriza copiar una implementación con términos incompatibles.
6. Una práctica sólo se vuelve normativa en un proyecto cuando tiene requirement, test/benchmark, riesgo, owner y rollback aplicables a ese proyecto.
7. Guardar la referencia fijada; enlaces móviles o ramas default no bastan para reproducibilidad.
8. Si dos líderes difieren, formular el trade-off y ensayar ambas opciones bajo el mismo contrato; no decidir por marca.

## 7. Quality scenario de seis partes

Usar:

```text
source
stimulus
environment
artifact
response
measure
```

Ejemplo correcto:

```text
Durante canary y carga pico, si un dependency excede su deadline,
el request path cancela trabajo downstream, no reintenta fuera de budget
y responde/encola según contrato; p99 y queue age permanecen bajo X.
```

“Debe ser rápido, seguro y escalable” no es verificable.

## 8. Arquitectura mínima suficiente

El manifest exige:

- context: actores/sistemas/trust boundaries;
- components: responsabilidad y ownership;
- data flows: qué cruza qué límite;
- ADRs: context, opciones, decisión, consecuencias y rollback trigger.

### 8.1 Orden de diseño

1. invariantes y fuente de verdad;
2. lifecycle/state machine;
3. synchronous/async boundaries;
4. failure/timeout/ambiguous result;
5. identity/authz/tenant boundary;
6. data ownership/retention/recovery;
7. capacity/performance budgets;
8. deploy/version/rollback;
9. observability/incident path.

La topología se deriva de esto; no empezar por microservices, Kubernetes, framework o database favorita.

### 8.2 ADR útil

Una ADR se acepta sólo si:

- presenta al menos dos opciones reales;
- explica drivers y evidencia;
- enumera consecuencias negativas;
- declara compatibilidad/migración;
- define qué observación reabre o revierte la decisión.

## 9. Contratos ejecutables

### 9.1 API/protocolo

- schema, límites y versionado;
- identity/authentication/authorization;
- timeout/deadline/cancellation;
- ordering/idempotency/dedupe;
- errors y ambiguous outcome;
- pagination/resume/replay;
- rate/admission/backpressure;
- compatibility y deprecation.

### 9.2 Datos

- owner y source of truth;
- grain, IDs, event/effective time;
- schema/null/order/precision;
- transaction/isolation/durability;
- correction/delete/retention;
- lineage/access/privacy;
- migration/backfill/reconcile;
- freshness/quality SLO.

### 9.3 Lifecycle

Toda entidad o proceso crítico tiene estados, transiciones permitidas, comandos idempotentes, terminalidad, recovery y observabilidad. Rechazar estados imposibles; no inferirlos de campos independientes contradictorios.

## 10. Vertical slice antes de expansión

Implementar una trayectoria del input real al resultado real:

```text
entrypoint
→ validation/authz
→ domain decision
→ state/data effect
→ response/publication/UI
→ metrics/log/trace
→ failure/recovery
→ package/deploy
```

La slice incluye un caso feliz, un fallo material, un test y una métrica. Sólo después se multiplica features, endpoints, models, platforms o partitions.

## 11. Estrategia de tests

### 11.1 Por obligación

| Pregunta | Test |
|---|---|
| función pura correcta | unit/property/model-based |
| módulos encajan | integration |
| consumidor/proveedor compatibles | contract |
| journey real funciona | system/E2E |
| falla/retry/recovery seguros | fault/recovery |
| adversario no cruza boundary | security/fuzz/red-team |
| budget real se cumple | performance/load/soak |
| usuario completa tarea | usability/accessibility |

### 11.2 Oracle

El oracle debe decidir, no sólo observar “no crash”:

- expected exact result;
- invariant/history model;
- reference implementation;
- checksum/digest;
- quality tolerance/slices;
- durable state after recovery;
- semantic accessibility/task completion.

### 11.3 Errores como datos de aprendizaje

Para cada fallo:

```text
observación → clasificación → hipótesis → experimento diferenciador
→ causa → corrección → test de regresión → prevención sistémica
```

Registrar si el error provino de data, objective, code, interface, environment, operation o assumption. Ésta operacionaliza la disciplina de error analysis: no arreglar ejemplos al azar ni cambiar varias variables a la vez.

En un producto multi-tenant, un fallo observado no autoriza reutilizar el payload del cliente. Separar tres carriles con contratos distintos:

1. **telemetría operacional:** mínima, redactada, con propósito, owner, acceso y retención;
2. **corpus de evaluación:** casos reproducibles sintetizados o expresamente autorizados, con provenance, tenant boundary y deletion lineage;
3. **datos de entrenamiento:** incorporación independiente, nunca implícita, con base contractual/consentimiento aplicable, evaluación de leakage/poisoning y posibilidad de exclusión o borrado según el compromiso asumido.

El aprendizaje global preferido es primero estructural: taxonomía de causa, invariante violado, patrón de ataque, control faltante y regresión sintética. Promover cambios sólo después de comparar baseline/candidate sobre holdout independiente, slices por cliente/riesgo, seguridad, costo y rollback. Más clientes amplían cobertura potencial; no convierten automáticamente los datos en utilizables, representativos ni confiables.

## 12. Benchmark profesional

Un benchmark válido fija:

- hipótesis;
- artifact/config digest;
- hardware/OS/runtime/toolchain;
- workload/distribución/skew/burst;
- correctness/quality gate;
- warmup/duración/repeticiones;
- offered load y goodput;
- p50/p95/p99/p99.9/max según riesgo;
- resource/energy/cost;
- baseline y candidate;
- raw evidence.

### 12.1 Secuencia

```text
correctness gate
→ baseline profile
→ localizar cuello
→ una modificación
→ repeated before/after
→ tails + resources + quality
→ sustained/recovery
→ accept/reject
```

No aceptar una mejora de kernel, handler, query o render si el camino end-to-end, la semántica o la recuperación empeoran.

## 13. Riesgo, seguridad y privacidad

Cada riesgo enlaza amenaza/fallo, impacto, probabilidad, controles, tests y riesgo residual.

Checklist mínima:

- assets y adversarios;
- trust boundaries;
- least privilege/authz por objeto/acción;
- input limits y parsers hostiles;
- secrets/keys/session lifecycle;
- tenant/data isolation;
- propósito y autorización para reutilizar telemetría, evals o datos de entrenamiento;
- poisoning, contribución desproporcionada y leakage entre tenants en loops de aprendizaje;
- dependency/build/update chain;
- logs/telemetry/redaction;
- backup/restore y breakglass;
- incident/revocation/rollback.

Un scanner, sandbox, firma, guardrail o WAF es un control parcial; ninguno prueba por sí solo que el sistema sea seguro.

## 14. SLO y operación

Un SLO contiene indicator, objective, window, budget y measurement exacta. Distinguir:

- service SLO;
- data freshness/quality SLO;
- model quality/safety gate;
- client crash/performance SLO;
- internal capacity target.

Alertar sobre impacto o burn, no sobre cualquier métrica ruidosa. Todo alert tiene owner, severidad, diagnóstico y acción; si no, es telemetry, no alerta.

### 14.1 Runbook mínimo

```text
trigger/symptom
safety actions
scope/impact query
diagnostic branches
mitigation
rollback/failover
data validation/reconciliation
recovery criteria
evidence preservation
communication/owners
follow-up tests and prevention
```

## 15. Release y compatibilidad

No hay “done” antes de:

- artifact identity completa;
- package/container/app instalado y probado;
- schema/config/protocol compatibility;
- secrets/IAM/admission correctos;
- staged rollout y gates;
- rollback o forward recovery ensayado;
- symbols/mappings/runbooks disponibles;
- field telemetry ligada al digest.

Promover el mismo artifact. Rebuild “equivalente” requiere una prueba de reproducibilidad, no confianza.

## 16. Contrato de tarea para Codex

Usar este encabezado al pedir implementación:

```text
Objetivo:
Workspace y scope autorizado:
Manifest del proyecto:
Manuales autoridad:
Invariantes:
Quality/SLO budgets:
Threat/failure model:
Compatibility/migration:
Entregables:
Comandos de verificación:
Acciones externas/destructivas no autorizadas:

Reglas:
- inspeccionar antes de modificar;
- preservar cambios ajenos;
- implementar una vertical slice verificable;
- no declarar passed sin artifact/evidencia;
- no optimizar antes de correctness/profile;
- registrar supuestos y ADR material;
- ejecutar tests proporcionales al riesgo;
- informar límites y remaining OPEN items.
```

### 16.1 Salida esperada del agente

1. outcome concreto;
2. archivos/contratos cambiados;
3. decisiones y trade-offs;
4. tests/benchmarks ejecutados con resultados;
5. evidencia/artifacts;
6. riesgos o gaps reales;
7. próxima acción sólo si queda trabajo autorizado.

## 17. Flujo para proyecto nuevo

1. copiar `project.example.json` fuera del kit;
2. declarar modo `NEW`, capturar baseline `EMPTY|SCAFFOLD`, crear el failure ledger y checkpointar revisión 1;
3. completar intake/scope;
4. seleccionar authority docs;
5. escribir requirements/tests juntos;
6. definir quality scenarios, risks y SLO;
7. diseñar context/state/data/ADRs;
8. completar `implementation_assurance` y ejecutar `--level plan`;
9. enlazar `tasks.md`, clasificar cada task exactamente una vez y seleccionar vertical slice;
10. implementar vertical slice;
11. ejecutar fault/security/performance/recovery;
12. guardar artifacts y SHA-256 en `evidence`;
13. cambiar estados sólo con evidencia y agregar checkpoint durable;
14. ejecutar `--level evidence`;
15. canary/observe/expand o rollback.

## 18. Flujo para mejorar un sistema existente

1. no reescribir por intuición;
2. declarar modo `EXISTING`, inventariar código/IaC/contratos/datos/evidencia y checkpointar baseline;
3. marcar hechos `PROVEN`, contradicciones como blockers y enlazar el delta autorizado sin crear owners paralelos;
4. capturar current architecture y deployed artifact;
5. reproducir síntoma con timeline/trace;
6. declarar invariante/SLO afectado;
7. formar hipótesis diferenciables;
8. añadir regression test antes o junto con fix;
9. modificar la capa responsable, no ocultar en UI/retry/cache;
10. comparar baseline/candidate;
11. probar version skew, migration y rollback;
12. actualizar ADR/runbook/contract si cambió comportamiento y checkpointar sólo el delta demostrado.

## 19. Flujo de incidente

```text
protect people/data/correctness
→ establish incident command and timeline
→ bound blast radius
→ preserve evidence
→ mitigate with reversible control
→ validate state/data
→ recover against explicit criteria
→ find contributing system conditions
→ add regression/fault test and structural prevention
```

No usar postmortem para encontrar una persona culpable. Buscar qué señal, límite, review, test, isolation o recovery faltó.

## 20. Capstones

### 20.1 IA en producción

`engineering_execution_kit/capstones/ai_production_service.json` obliga a unir dataset/model/eval/runtime identity, quality/safety slices, serving tails, canary y rollback. No acepta promedio de benchmark como calidad del producto.

### 20.2 Order book de baja latencia

`engineering_execution_kit/capstones/low_latency_order_book.json` obliga sequence/gap/session state, snapshot+replay, bounded publication, deterministic oracle y end-to-end tails. No introduce predicción ni order entry.

### 20.3 Cliente nativo offline-first

`engineering_execution_kit/capstones/offline_native_client.json` obliga local source-of-truth, transactional outbox, process death, conflictos, accesibilidad, device performance y signed staged update.

Los tres manifests pasan hoy el gate de `plan`. Permanecen intencionalmente sin evidencia: son especificaciones ejecutables de proyectos futuros, no sistemas ya implementados.

## 21. Definition of Ready

- [ ] para producto nuevo, problema/comprador y umbral de validación definidos antes de ampliar construcción;
- [ ] objetivo, owner, scope y out-of-scope;
- [ ] invariantes y fuente de verdad;
- [ ] requirements ↔ tests;
- [ ] quality scenarios medibles;
- [ ] architecture/context/state/data;
- [ ] riesgos, controles y residual;
- [ ] política de datos/telemetría/aprendizaje y tenant boundary declarados cuando corresponda;
- [ ] SLO/budgets y measurement;
- [ ] compatibility/migration/release;
- [ ] authority docs mínimos;
- [ ] procedencia, risk tier, nueve dimensiones y source locks aplicables en `implementation_assurance`;
- [ ] `--level plan` pasa.
- [ ] estado NEW/EXISTING, baseline/delta, tasks y próximo paso pasan `--level resume` con checkpoint íntegro.

## 22. Definition of Done

- [ ] código/config/schema/artifact versionados;
- [ ] requisitos y journeys aceptados;
- [ ] fault/security/performance/recovery según riesgo;
- [ ] package/deploy real probado;
- [ ] compatibility/migration/rollback ejecutados;
- [ ] telemetry/alerts/runbooks operables;
- [ ] retención, borrado y promociones del learning loop probados cuando usa datos reales;
- [ ] artifacts con environment/tool/date/SHA-256;
- [ ] tests/benchmarks sólo `passed` con evidence;
- [ ] `implementation_assurance=proven`, cero blockers y evidencias separadas para artefacto/release;
- [ ] `--level evidence` pasa;
- [ ] `PROJECT_EXECUTION_STATE.json --level complete` pasa contra el checkpoint final y el target declarado;
- [ ] remaining `[OPEN]` no contradice el objetivo.

## 23. Qué valida y qué no

| El kit sí valida | El kit no puede validar solo |
|---|---|
| campos y estructura | corrección del código no ejecutado |
| IDs únicos | exhaustividad de un threat model |
| requirement/test/risk/evidence/assurance refs | representatividad del workload/dataset |
| `passed` con evidence ref | que el oracle sea semánticamente correcto |
| authority docs existentes | que se hayan leído bien |
| evidence path contenido | autenticidad externa del productor |
| digest del archivo | verdad universal de la conclusión |
| todos passed en evidence gate | producción futura sin fallos |

Por eso el gate automático complementa revisión humana, tests de dominio, medición y operación; no los suplanta.

## 24. Puerta de cierre

- [x] Selección mínima del corpus.
- [x] Intake, requirements, quality scenarios y ADR.
- [x] Contratos API/data/lifecycle.
- [x] Tests, benchmark, seguridad, SLO y runbook.
- [x] Release, compatibility y rollback.
- [x] Contrato reusable para Codex.
- [x] JSON Schema y ejemplo.
- [x] Validador plan/evidence sin dependencies.
- [x] Once tests del contrato: planes válidos, schema/refs, IDs, trazabilidad, evidencia/digest, riesgo, source locks y rechazo compile-only.
- [x] Tres capstones con gate plan aprobado.
- [x] Resume cursor NEW/EXISTING que no duplica readiness, Spec Kit ni tasks.
- [x] Checkpoint log append-only con SHA-256, `fsync`, transición y confinamiento de paths.
- [x] Trece regresiones adicionales de baseline/delta, tareas, contexto, evidencia, tamper y cierre; 24/24 totales.
- [ ] Evidencia de implementación de capstones: deliberadamente pendiente hasta ejecutarlos como proyectos reales.

La biblioteca ya no depende sólo de “leer buenos Markdown”: contiene una ruta verificable para convertirlos en decisiones, código, pruebas y evidencia sin afirmar resultados que aún no existen.
