# Public Code & Architecture Admission Standard

> **Estado:** v1, creado el 2026-08-24.
> **Propósito:** impedir que código público, popular o promocional se incorpore a esta biblioteca como arquitectura de élite sin demostrar licencia, corrección, adecuación, seguridad, operación y posibilidad real de reutilización.
> **Autoridad interna:** `SYSTEMS_ENGINEERING_MASTER_MAP.md`, `AI_ENGINEERING_MASTER_MAP.md`, `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md`, `ENGINEERING_EXECUTION_PLAYBOOK.md` y los manuales de cada frontera atravesada.
> **Límite:** esta norma organiza una auditoría técnica y de procedencia. No sustituye asesoramiento legal sobre una combinación, distribución, jurisdicción o modelo comercial concreto.

## 1. Regla absoluta

```text
Público != open source.
Open source != adecuado.
Adecuado != probado.
Probado por su autor != reproducido por nosotros.
Arquitectura de una empresa líder != arquitectura correcta para cualquier proyecto.
```

Un artefacto sólo puede incorporarse si mejora materialmente una decisión, implementación, prueba, medición, operación o recuperación y no contradice los contratos existentes.

Si el código no alcanza el nivel de los manuales, los manuales prevalecen y el código se rechaza o queda condicionado.

## 2. Estados permitidos

| Estado | Significado | Uso permitido por Codex |
|---|---|---|
| `DISCOVERED` | fuente localizada, todavía no auditada | ninguno; sólo backlog de investigación |
| `LICENSE_VERIFIED` | revisión y términos observados, sin aprobación técnica | citar licencia; no copiar ni adoptar arquitectura |
| `CANDIDATE` | licencia y relevancia inicial aceptables | estudiar; generar preguntas y probes |
| `CONDITIONED` | útil sólo bajo licencia, stack, escala o integración explícitos | usar como referencia con decisión y revisión legal/técnica |
| `ELITE_REFERENCE` | supera G0–G7 y aporta evidencia arquitectónica real | reutilizar patrones dentro de sus límites |
| `REUSABLE_PACK` | además posee adaptación, contratos, tests y gate de proyecto | Codex puede proponerlo para implementación inmediata |
| `REJECTED` | licencia, calidad, seguridad, evidencia o adecuación insuficiente | no usar; conservar razón para evitar redescubrimiento |
| `RETIRED` | antes admitido, pero licencia, mantenimiento o evidencia cambió | no iniciar adopciones; migrar consumidores según plan |

Ningún README puede declarar por sí solo `ELITE_REFERENCE`. Ningún pack se considera `passed` sin evidencia del proyecto que lo adopta.

Si el catálogo no contiene una opción compatible para una capability requerida, usar `markdown_system/CAPABILITY_GAP_RESOLUTION_PACK_PLAN.md`: investigación primaria actual, evidencia local, identidad/licencia/revisión/hash y G0–G8 por candidato. Un expediente `RESEARCH_INCOMPLETE` no es una decisión y debe continuar. `NO_ADMISSIBLE_SOURCE` exige alcance y agotamiento demostrados, blockers y trigger de reapertura; no autoriza presentar código `AUTHORED` como upstream. Sólo un resultado `USE_REUSABLE_PACK` permite implementación inmediata, todavía sujeto a gates del target.

## 3. Unidad de auditoría

La unidad no es el nombre del proyecto. Es:

```text
owner/repository
@ commit o tag inmutable
+ path/subproject exacto
+ archivo de licencia aplicable
+ third-party notices/SBOM
+ documentación y versión
+ configuración/runtime/hardware relevante
+ claim concreto que se quiere reutilizar
```

Un monorepo puede contener MIT, Apache-2.0, GPL/AGPL, código empresarial y assets con términos diferentes. La licencia de la raíz nunca se extrapola silenciosamente a todos los paths.

## 4. Gates obligatorios

### G0 — Identidad y procedencia

- propietario oficial del repositorio y relación con el sistema;
- commit/tag inmutable y fecha de observación;
- repositorio activo, archivado, fork o export parcial;
- release estable frente a branch de desarrollo;
- archivos generados, vendorizados, submódulos y terceros identificados;
- para archives GitHub generados desde un commit, no confundir bytes del contenedor comprimido con identidad del source: GitHub advierte que el layout de compresión puede cambiar. Fijar commit, tree/firma cuando exista, raíz única segura, inventario y hash canónico de path+tamaño+SHA-256 por archivo; conservar el hash de transporte observado en el receipt. Un release asset inmutable sí conserva su digest de bytes exacto;
- diferencia entre código público y arquitectura propietaria no publicada.

**Falla automática:** mirror no verificable, autoría ambigua, tarball sin provenance o claim sobre internals ausentes.

### G1 — Licencia y reutilización

- texto de licencia leído en la revisión fijada;
- identificador SPDX cuando exista;
- aprobación OSI comprobada cuando se afirme “open source”;
- copyleft, network copyleft, file-level copyleft, excepciones y source-available distinguidos;
- paths de edición empresarial y componentes con otra licencia separados;
- notices, atribución, patentes, marcas, datos, modelos y assets revisados;
- compatibilidad con el modo esperado: interno, SaaS, on-premise, distribución, plugin o derivado.

Fuentes de control: [lista SPDX](https://spdx.org/licenses/) y [licencias aprobadas por OSI](https://opensource.org/licenses).

**Falla automática:** `NOASSERTION` tratado como licencia, “GitHub público” usado como permiso, términos no localizados o restricción incompatible con el proyecto.

### G2 — Alineación con los manuales

Por cada claim se seleccionan:

- manual autoridad;
- secciones e invariantes aplicables;
- fronteras afectadas;
- contradicciones;
- delta real que el código aporta.

Matriz mínima:

| Claim del código | Autoridad interna | Coincide | Diverge | Decisión |
|---|---|---|---|---|
| comportamiento | manual/sección | evidencia | riesgo | aceptar, adaptar o rechazar |

**Falla automática:** adoptar porque “lo usa una gran empresa” sin demostrar que carga, consistencia, amenaza, organización y costo coinciden.

### G3 — Arquitectura observable

Debe poder reconstruirse desde código y artefactos:

- context, runtime y deployment views;
- módulos, responsabilidades y dependency rules;
- fuentes de verdad y ownership de datos;
- protocolos, schemas y compatibilidad;
- estados y transiciones;
- concurrencia, idempotencia, backpressure y failure model;
- configuración, extensión y lifecycle;
- migración, upgrade y rollback.

Un diagrama sin código o un repositorio sin límites explicados no supera este gate.

### G4 — Correctitud y calidad de implementación

- invariantes visibles en tipos, constraints o validadores;
- tests unitarios, property, contract, integration y system según riesgo;
- casos negativos y de error;
- determinismo/replay cuando corresponda;
- compatibilidad y migraciones probadas;
- ausencia de secretos o defaults evidentemente inseguros;
- dependencias y toolchain reproducibles;
- issues/advisories materiales revisados.

Las estrellas, cobertura agregada o un CI verde no sustituyen inspección de los caminos críticos.

### G5 — Seguridad y privacidad

- threat model y trust boundaries;
- autenticación separada de autorización;
- deny-by-default y least privilege;
- autorización por objeto/tenant/acción donde aplica;
- validación en cada frontera;
- secretos, cifrado, supply chain y provenance;
- aislamiento multi-tenant;
- auditoría de intención, decisión y efecto;
- retención, borrado, minimización y datos sensibles;
- abuse cases, rate limits y resource exhaustion;
- advisory process y política de actualización.

Un framework de identidad, scanner o sandbox aislado no demuestra seguridad end-to-end.

### G6 — Rendimiento y resiliencia

- workload y SLO declarados;
- offered load, goodput y percentiles, no sólo promedio;
- queues, pools, caches y buffers acotados;
- overload, timeout, cancellation y retry budget;
- failover, restart, partition y version skew;
- capacidad, costo y resource high-water;
- benchmark reproducible sobre artefacto desplegable;
- fallback y rollback probados.

“Low latency” se acepta sólo con camino, hardware, distribución de carga y evidencia. No se exige baja latencia a componentes donde consistencia, costo o mantenibilidad dominan.

### G7 — Operación y evolución

- métricas, logs y traces accionables;
- alertas ligadas a SLO o impacto;
- runbooks y ownership;
- backup/restore y disaster recovery;
- releases inmutables, canary y rollback;
- deprecation y compatibility window;
- mantenimiento activo y gobernanza;
- evidencia de incident learning o failure injection cuando sea material.

### G8 — Conversión a `REUSABLE_PACK`

Para que Codex pueda usarlo inmediatamente se requiere además:

- problema y no-objetivos;
- condiciones de selección y rechazo;
- arquitectura portable, no atada accidentalmente al producto fuente;
- código/adaptadores propios o paths reutilizables identificados;
- contracts y schemas;
- configuración segura;
- migrations/seed sólo cuando sean reproducibles;
- test harness y oráculos;
- benchmark manifest;
- observabilidad, despliegue y rollback;
- manifest de licencias/notices;
- checklist de integración con los manuales;
- estado `planned` hasta ejecutarse dentro de un proyecto real.

## 5. Scoring auxiliar

El scoring ayuda a ordenar; nunca compensa un gate obligatorio fallido.

| Dimensión | 0 | 1 | 2 | 3 |
|---|---|---|---|---|
| procedencia | dudosa | indirecta | oficial | oficial + revisión fijada |
| licencia | ausente/incompatible | ambigua/mixta | clara condicionada | clara y compatible |
| arquitectura | marketing | diagrama | docs + código | código/tests/operación |
| correctitud | sin tests | happy path | capas principales | fallos/property/recovery |
| seguridad | no declarada | checklist | controles | threat model + pruebas |
| rendimiento | adjetivos | microbench | workload | end-to-end reproducible |
| operación | ausente | configuración | observabilidad | SLO/runbooks/restore |
| portabilidad | acoplamiento oculto | vendor-specific | adaptadores | contrato portable probado |

Requisitos mínimos:

- `ELITE_REFERENCE`: ningún cero, G0–G7 aprobados y al menos 20/24;
- `REUSABLE_PACK`: `ELITE_REFERENCE` + G8 completo;
- seguridad, licencia y correctitud nunca se promedian.

## 6. Regla de composición

Dos componentes admitidos individualmente pueden formar un sistema incorrecto. Antes de combinarlos se prueba:

```text
identity propagation
+ authorization semantics
+ transaction/event boundaries
+ schema/version compatibility
+ timeout/cancellation/retry interaction
+ data ownership and deletion
+ observability correlation
+ license compatibility
+ rollout/rollback order
```

Ejemplo: un IdP no reemplaza autorización por objeto; un workflow engine no reemplaza la base de negocio; un event broker no crea exactly-once end-to-end; un commerce engine no constituye ERP, CRM o contabilidad completa.

## 7. Reglas para snippets y código extraído

1. Preferir interfaces y patrones propios sobre copiar grandes bloques.
2. Todo bloque copiado conserva origen, commit, path, licencia y modificaciones.
3. No copiar assets, datasets, prompts, modelos o ejemplos comerciales por herencia implícita.
4. No introducir código sin tests del comportamiento requerido.
5. No introducir una dependencia para evitar escribir una función pequeña si aumenta TCB o costo operacional.
6. No portar bugs, defaults de demo ni configuración educativa.
7. No mezclar código GPL/AGPL/MPL con un producto de otra licencia sin una decisión explícita.
8. No usar branch head como dependency pin de producción.

## 8. Plantilla de expediente

```yaml
source:
  owner_repo: ""
  revision: ""
  observed_at: ""
  release_status: "stable|development|archived"
  paths: []
license:
  root_expression: ""
  path_overrides: []
  notices: []
  osi_approved: "yes|no|unknown"
  intended_use: "internal|saas|distribution|plugin|derivative"
claim:
  capability: ""
  evidence_label: "CODE"
  exact_claim: ""
  non_claims: []
alignment:
  authority_docs: []
  invariants: []
  contradictions: []
architecture:
  boundaries: []
  data_ownership: []
  protocols: []
  failure_model: []
evidence:
  code_paths: []
  tests: []
  benchmarks: []
  operations: []
decision:
  state: "DISCOVERED"
  approved_uses: []
  forbidden_uses: []
  remaining_gates: []
```

## 9. Reauditoría

Reabrir un expediente si cambia cualquiera:

- licencia o estructura de ediciones;
- propietario, gobernanza o estado archivado;
- major release o contrato público;
- advisory material;
- runtime/toolchain soportado;
- necesidad del proyecto;
- evidencia que contradiga el claim anterior.

Cada proyecto ejecuta un preflight actual. Una revisión histórica puede enseñar arquitectura, pero no autoriza adoptar hoy su dependencia o licencia.

Una vulnerabilidad, compromiso de maintainer/registry/signing key o EOL reabre automáticamente la admisión de código ya incorporado. Hasta resolver el expediente, su estado baja a `CONDITIONED` o `BLOCKED` según exposición. El agente aplica `markdown_system/DEPENDENCY_UPDATE_CONTRACT.md`: primero fija artefacto/consumers y autoridad oficial; luego usa release corregida, patch oficial aún no publicado, reemplazo o contención explícita. Un backport/local fix conserva enlace al diff upstream pero declara provenance propia `ADAPTED_PATCH`/`AUTHORED_PATCH`; jamás hereda falsamente el sello del proveedor.

## 10. Gate final

- [ ] revisión y paths exactos fijados;
- [ ] licencia por path y notices verificados;
- [ ] claim estrecho y non-claims escritos;
- [ ] manuales autoridad seleccionados;
- [ ] arquitectura reconstruida desde código;
- [ ] correctness/security/failure gates comprobados;
- [ ] operación, compatibilidad y rollback existentes;
- [ ] delta justifica el costo de incorporación;
- [ ] estado asignado sin exageración;
- [ ] pack no se marca `passed` sin evidencia del proyecto.

La calidad de esta biblioteca se protege más por lo que rechaza que por la cantidad de repositorios que acumula.
