# Codex Elite Project Bootstrap

> **Estado:** v1 operativo, 2026-08-25.
> **Propósito:** convertir una idea o repositorio existente en un proyecto ejecutable usando los manuales autoridad y únicamente arquitecturas/código público que hayan superado el gate correspondiente.
> **Salida:** un pack de decisión específico del proyecto antes de expandir implementación.

## 1. Principio

“Usar todos los Markdown” significa que Codex conoce el mapa completo, selecciona conscientemente las autoridades necesarias y demuestra por qué las demás no cambian la decisión. No significa cargar 200.000 palabras en cada turno ni mezclar todos los patrones.

```text
biblioteca completa como índice de autoridad
→ contexto mínimo por decisión
→ contrato específico del proyecto
→ implementación vertical
→ evidencia
```

## 2. Orden de lectura obligatorio

1. `CODEX_ELITE_PROJECT_BOOTSTRAP.md`.
2. `AI_ENGINEERING_MASTER_MAP.md` y `SYSTEMS_ENGINEERING_MASTER_MAP.md` para routing.
3. `ENTERPRISE_FULL_STACK_BLUEPRINT.md` cuando se diseña una aplicación empresarial completa.
4. `LEADING_COMPANY_PUBLIC_CODE_MATRIX.md`, `PUBLIC_CODE_ARCHITECTURE_MASTER_MAP.md` y `PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md` si se considera código público.
5. `ENGINEERING_EXECUTION_PLAYBOOK.md` para contrato, gates y evidencia.
6. Sólo los manuales autoridad y fronteras materiales de la tarea.
7. Los `REUSABLE_PACK` compatibles pueden incorporarse directamente. Un pack `REBUILD_VERIFIED / CONDITIONED` también puede materializarse cuando el agente demuestre sus condiciones en el proyecto. Los `CANDIDATE_PACK` sirven para planificar o completar admisión, no para adopción ciega.
8. Al materializar `ENGINEERING_EXECUTION_VALIDATOR.md`, crear `PROJECT_EXECUTION_STATE.json` desde la plantilla y `PROJECT_EXECUTION_EVENTS.jsonl` con `checkpoint_execution_state.py`. En cada reanudación validar primero el estado; cargar sólo `must_read_refs` y reutilizar por hash `reuse_without_reload_refs`.

## 3. Router por preocupación

| Pregunta material | Manual autoridad |
|---|---|
| drivers, bounded contexts, trade-offs, vistas, ADR y evolución | `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md` |
| backend, HTTP/gRPC, idempotencia, errores, auth y contracts | `SOFTWARE_BACKEND_API_ENGINEERING.md` |
| schema, constraints, transacciones, índices, WAL y recovery | `DATABASE_STORAGE_INTERNALS.md` |
| red, RPC, eventos, consistency, delivery y backpressure | `NETWORKING_DISTRIBUTED_STREAMING.md` |
| CPU, memoria, concurrencia, profiling y tails | `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md` |
| seguridad, IAM, cloud, supply chain, SRE, incident y restore | `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md` |
| frontend, HCI, accesibilidad, browser, UX y performance perceptible | `FRONTEND_PRODUCT_ENGINEERING_UX.md` |
| mobile/desktop, lifecycle, offline, firma y updates | `NATIVE_MOBILE_DESKTOP_ENGINEERING.md` |
| pipelines, formats, lineage, quality, metrics y analytics | `DATA_ENGINEERING_ANALYTICS.md` |
| builds, packages, ABI/FFI, reproducibility y releases | `TOOLCHAINS_BUILDS_PACKAGING_FFI.md` |
| algoritmos, estructuras, corrección y complejidad | `ALGORITHMS_DATA_STRUCTURES_PROBLEM_SOLVING.md` |
| GPU, kernels, memory, precision y multi-device | `GPU_ACCELERATED_COMPUTING.md` |
| order books, feeds, matching y riesgo de venue | `MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md` |
| ML/DL, modelos, RAG, agentes, inferencia y seguridad IA | seleccionar desde `AI_ENGINEERING_MASTER_MAP.md` |

## 4. Fase 0 — Intake antes de arquitectura

Codex debe obtener o descubrir:

```yaml
business:
  problem: ""
  users_buyers_operators: []
  value_and_failure_cost: ""
  jurisdictions_and_markets: []
scope:
  must_have_journeys: []
  later: []
  out_of_scope: []
domain:
  entities_and_invariants: []
  source_of_truth: []
  roles_and_organizations: []
quality:
  availability: ""
  latency: ""
  capacity: ""
  privacy_security: ""
  recovery: ""
delivery:
  team_and_skills: []
  environments: []
  budget_and_deadline: ""
  existing_assets_integrations: []
```

Si faltan datos, Codex puede construir un discovery skeleton, pero no inventar reglas comerciales, regulatorias o de seguridad irreversibles.

## 5. Fase 1 — Authority selection report

Antes de programar se crea `PROJECT_AUTHORITY_MAP.md`:

```md
# Project Authority Map

## Objetivo y alcance
## Journeys y failure cost
## Manuales autoridad seleccionados
## Fronteras materiales
## Manuales evaluados pero no cargados y razón
## Reusable packs candidatos
## Licencias y compatibilidad esperada
## Incertidumbres que bloquean garantías
```

Un candidato `DISCOVERED`, `LICENSE_VERIFIED` o `CANDIDATE` no puede entrar como dependencia por decisión automática. `CONDITIONED` exige comprobar condiciones; si además es `REBUILD_VERIFIED`, el agente puede hacerlo autónomamente y registrar evidencia. `REJECTED` está prohibido salvo investigación legal nueva que cambie el estado.

### 5.1 Criterio de aceptación de generación compacta — V329

La automatización pendiente debe producir las ocho secciones anteriores desde
decisiones trazables del scope activo. Una lista de `codex.authority_docs`, un
enlace al plan o el descubrimiento de archivos no justifican la selección.
El owner del trabajo sigue siendo T2801 del roadmap; esta sección precisa su
aceptación y no constituye un pack implementado ni una aprobación de proyecto.

Cada decisión debe conservar ID, requisito/journey de origen, owner, manual y
sección aplicable, motivo, frontera afectada, evidencia y estado explícito. Los
manuales evaluados pero no cargados requieren motivo por scope; los candidatos
conservan revisión/licencia/condiciones, y las incertidumbres mantienen bloqueo.
La selección razonada corresponde al agente bajo las autoridades del proyecto;
un generador determinista sólo puede comprobar y representar esa selección,
sin resolver por sí mismo decisiones de negocio ni declarar adecuación semántica.

Para cerrar la generación deben demostrarse:

- NEW: salida completa con trazabilidad, sin respuestas ni aprobaciones inventadas.
- EXISTING: selección limitada al delta, evidencia previa preservada y sucesor
  explícito cuando cambia una decisión; no reemplazo silencioso del mapa manual.
- Rechazo de requisitos sin cobertura, decisiones duplicadas o contradictorias,
  referencias ausentes/fuera de raíz, revisiones vencidas y motivos vacíos.
- Resultado estable para la misma entrada fijada, sin seleccionar el plan por
  mtime ni leer todo el corpus; presupuesto finito declarado antes del ensayo.
- Publicación que conserva archivos existentes e identifica bytes/versión de
  entrada y salida; recuperación probada según el riesgo del escritor.
- Evaluación de pertinencia frente a casos revisados y evidencia del journey
  de bootstrap; estructura válida no equivale a selección correcta. Ahorro de
  tokens requiere BENCH01, independiente de tamaño de archivos o tiempo de CLI.

Investigación acotada y bloqueos: `reconstruction_evidence/AUTHORITY_MAP_GAP_V329.md`.
El resolver de capabilities debe alcanzar `USE_REUSABLE_PACK` antes de implementar
este gap en el proyecto; `NO_ADMISSIBLE_SOURCE` conserva la brecha abierta.

## 6. Fase 2 — Project engineering contract

Materializar primero `implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md` en el root del proyecto; esto crea `engineering_execution_kit/` como output reproducible, no como fuente paralela. Crear un manifest desde `engineering_execution_kit/project.example.json`, crear el estado desde `execution_state.template.json` y validarlos con:

```powershell
python .\engineering_execution_kit\validate_project.py `
  .\path\project.json --level plan

python .\engineering_execution_kit\validate_execution_state.py `
  .\PROJECT_EXECUTION_STATE.json `
  --project-root . `
  --events .\PROJECT_EXECUTION_EVENTS.jsonl `
  --level resume
```

Después de cada cambio material se ejecuta `checkpoint_execution_state.py`; el cursor no es un segundo plan y no puede contradecir los owners que referencia.

El contrato debe incluir como mínimo:

- requirements ↔ acceptance tests;
- quality scenarios de seis partes;
- context/component/runtime/deployment views;
- ADR con alternativas y rollback trigger;
- API/data/lifecycle contracts;
- threat/failure model y riesgos residuales;
- SLOs, benchmarks y operation plan;
- release identity, compatibility, migration y rollback;
- authority docs y forbidden assumptions.
- connected journeys release-bound: persona/rol → interfaz → autorización → contrato → dominio → dato/efecto → auditoría → respuesta/error → ayuda → capacitación → soporte → update → E2E/negativos → recovery.

## 7. Fase 3 — Architecture decision ladder

### 7.0 Baseline de composición

Crear primero `PROJECT_BLUEPRINT.md` y seleccionar implementation packs mediante `AGENT_SYSTEM_START.md` y `markdown_system/COMPOSITION_PROTOCOL.md`. No elegir framework por costumbre. Cada stack o sustitución exige ADR, compatibilidad demostrada, costo de migración, evidencia y rollback.

Todo proyecto que vaya a staging o producción compone `implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md` y reemplaza sus placeholders con evidencias reales. Si el blueprint selecciona PostgreSQL, compone además `implementation_packs/POSTGRES_BACKUP_RESTORE_CORE.md` y ejecuta el restore drill en el entorno objetivo. Estos packs no eligen cloud, CI ni proveedor; obligan a que esas selecciones queden explícitas y probadas.

Un prototipo ejecutable no gobierna todos los proyectos. Sólo los bloques `RECONSTRUCTIBLE` o `REBUILD_VERIFIED` se materializan automáticamente, y sólo un `REUSABLE_PACK` supera el gate de admisión para su claim estrecho.

### 7.1 Topología

Default para un producto empresarial nuevo:

```text
modular monolith
+ PostgreSQL source of truth
+ explicit module contracts
+ transactional outbox
+ background workers
+ versioned APIs
+ one observable deployable
```

Separar servicios sólo si existe evidencia de al menos una fuerza:

- ownership/equipo verdaderamente independiente;
- deployment cadence incompatible;
- SLO o escala muy diferente;
- isolation/security boundary real;
- tecnología/hardware irreconciliable;
- failure containment medible;
- requisito externo de integración.

No usar microservices, Kafka, Kubernetes, Temporal, service mesh o vector DB como símbolos de calidad.

### 7.2 Build versus adopt

Para cada capacidad:

| Opción | Licencia | Ajuste dominio | Riesgo/TCB | Operación | Exit cost | Evidencia |
|---|---|---:|---:|---:|---:|---|
| construir | propia | | | | | |
| adoptar | exacta | | | | | |
| servicio externo | terms/SLA | | | | | |

Se adopta sólo si el costo total y el riesgo son menores que construir el contrato necesario.

## 8. Fase 4 — Domain slice

Antes de ampliar features se implementa un vertical slice real:

```text
UI/accessibility
→ API/schema/authn/authz
→ domain invariant
→ transaction/source of truth
→ event/outbox if needed
→ telemetry
→ tests
→ package/deploy
→ rollback
```

Debe atravesar el camino desplegable, no mocks de todas las fronteras.

## 9. Fase 5 — Gates

### Correctness

- invariantes de dominio;
- property y boundary tests;
- concurrency/idempotency histories;
- migrations y version skew;
- oráculo independiente.

### Security

- threat model;
- identity y authorization por recurso/tenant;
- validación y rate/resource limits;
- secrets y supply chain;
- negative/abuse tests;
- audit sin payload sensible.

### Performance

- workload representativo;
- offered load y goodput;
- p50/p95/p99/p99.9 donde sea material;
- resource high-water;
- overload y recovery;
- comparación baseline/candidate.

### Operations

- SLI/SLO/error budget;
- logs/metrics/traces correlacionados;
- alertas accionables;
- backup/restore;
- incident runbook;
- canary y rollback.

## 10. Fase 6 — Evidencia

No escribir “passed” sin artifacts. Capturar resultados con tool, environment, fecha y SHA-256; después ejecutar:

```powershell
python .\engineering_execution_kit\validate_project.py `
  .\path\project.json --level evidence
```

## 11. Instrucción maestra reutilizable

Copiar y completar:

```md
Actúa como arquitecto e implementador principal de este proyecto.

OBJETIVO
[resultado empresarial y usuarios]

ALCANCE
[journeys must-have, later y out-of-scope]

RESTRICCIONES
[jurisdicción, equipo, plazo, stack existente, presupuesto, deployment]

CALIDAD
[seguridad/privacidad, disponibilidad, latencia, capacidad, RPO/RTO]

INTEGRACIONES
[proveedores, fábricas, marketplaces, pagos, ads, CRM, ERP, logística]

PROTOCOLO OBLIGATORIO
1. Lee CODEX_ELITE_PROJECT_BOOTSTRAP.md.
2. Usa AI_ENGINEERING_MASTER_MAP.md y SYSTEMS_ENGINEERING_MASTER_MAP.md como router; no cargues toda la biblioteca indiscriminadamente.
3. Si consideras código externo, aplica PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md y consulta PUBLIC_CODE_ARCHITECTURE_MASTER_MAP.md.
4. No incorpores nada por fama, estrellas o marketing. Usa un `REUSABLE_PACK`, un `REBUILD_VERIFIED / CONDITIONED` cuyas condiciones puedas demostrar, o crea trabajo `AUTHORED` desde los contratos autoridad con tests y procedencia. Abre un expediente cuando reutilices código público no admitido.
5. Produce PROJECT_AUTHORITY_MAP.md y el manifest JSON del Engineering Execution Kit antes de expansión material.
6. Parte de arquitectura mínima suficiente; modular monolith es el baseline salvo evidencia contraria.
7. Implementa un vertical slice desplegable con contratos, tests, seguridad, telemetría y rollback.
8. Mantén autenticación, autorización, dominio, persistencia, integraciones y presentación en fronteras explícitas.
9. Toda integración externa usa adapter, schema versionado, idempotencia, deadline, retry budget, reconciliation, rate limit, secret isolation y contract tests.
10. No declares low latency, secure, scalable, exactly-once, production-ready ni passed sin evidencia del entorno.
11. Registra decisiones, contradicciones con los manuales, riesgos residuales y condiciones de reversión.
12. Detente ante licencia incompatible, secreto, requisito legal no resuelto o acción irreversible sin autoridad.

ENTREGA INICIAL
- problem/journey map;
- authority selection;
- domain model e invariantes;
- architecture alternatives y ADR inicial;
- threat/failure model;
- project manifest validado a nivel plan;
- vertical slice seleccionado;
- archivos y tests a crear;
- candidatos públicos admitidos y licencias;
- riesgos OPEN y próximos gates.
```

## 12. Extensión para franquicias y movilidad eléctrica

Cuando el proyecto sea una red de franquicias de motos/bicicletas eléctricas, el intake debe resolver además:

- organización matriz, franquiciado, sucursal, territorio y delegación;
- usuarios cliente/empleado/manager/admin/proveedor/fábrica y permisos por objeto;
- modelo, variante, VIN/serie, batería, lote, repuesto y documentación;
- catálogo público frente a información interna;
- price books, moneda, impuestos, promociones y vigencia;
- stock por estado/localización, reserva, transferencia y conteo;
- lead, consentimiento, campaña, atribución y lifecycle CRM;
- compra a proveedor, producción, hitos, logística, recepción y calidad;
- venta/reserva/pago/factura/entrega;
- garantía, service, mantenimiento, recall y trazabilidad;
- marketplaces/ads con adapter por canal y source of truth explícita;
- auditoría, privacidad, retención y jurisdicción;
- modo degradado/offline para sucursales si se justifica;
- analytics sin convertir eventos en contabilidad o fuente de verdad.

## 13. Anti-patrones prohibidos

- empezar por framework o cloud;
- copiar un repositorio completo como arquitectura;
- importar un sample educativo como producción;
- usar roles globales para acceso multi-tenant por objeto;
- compartir base sin ownership de writes;
- llamar evento a una transacción;
- reintentar side effects sin idempotencia/reconciliation;
- telemetry con PII/secrets;
- cola, cache, batch o retry sin límite;
- microservices sin ownership/SLO independiente;
- benchmark sin correctness gate;
- licencia de raíz aplicada a paths Enterprise;
- `main`, `master`, `develop` o `latest` como pin productivo;
- marcar terminado porque compila.

## 14. Definition of Ready para programar

- [ ] usuarios, journeys y failure cost definidos;
- [ ] dominio, invariantes y sources of truth;
- [ ] authority map completo;
- [ ] opciones públicas con licencia y estado admitido;
- [ ] arquitectura y ADR inicial;
- [ ] threat/failure model;
- [ ] quality scenarios y SLOs;
- [ ] requirements ↔ tests;
- [ ] dependency/runtime/upstream inventory, update policy, EOL, SBOM y rollback;
- [ ] authorities temporales verificadas, sin stale/conflict blocker;
- [ ] migration/release/rollback;
- [ ] manifest pasa `--level plan`;
- [ ] vertical slice escogido.

Si este gate no pasa, Codex sigue descubriendo y especificando. No compensa ambigüedad con más código.

## 15. Activación desde la primera interacción

Esta biblioteca incluye un `AGENTS.md` raíz conciso. Codex lo descubre antes de trabajar cuando la tarea se inicia con esta carpeta como workspace; ese archivo obliga a usar este bootstrap y los mapas como router sin intentar cargar todo el corpus.

Para un proyecto creado fuera de esta carpeta hay dos opciones válidas:

1. iniciar el proyecto dentro de este workspace y mantener la cadena de instrucciones; o
2. copiar al root del nuevo repositorio un `AGENTS.md` de proyecto que apunte a una copia/version fijada de esta biblioteca y conservar accesibles los archivos autoridad.

Verificación recomendada al abrir una nueva sesión:

```text
Resume las instrucciones de proyecto que cargaste y enumera sus archivos fuente.
Luego ejecuta el intake de CODEX_ELITE_PROJECT_BOOTSTRAP.md para este proyecto.
```

La activación se reconstruye al iniciar una nueva ejecución/sesión de Codex. Crear o modificar `AGENTS.md` dentro de una sesión ya iniciada no demuestra que esa misma sesión lo haya recargado.
