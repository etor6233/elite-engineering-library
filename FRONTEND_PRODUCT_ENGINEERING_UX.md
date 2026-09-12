# Frontend, producto y UX — manual operativo

> **Estado:** núcleo general v1 auditado y cruces recíprocos cerrados el 2026-08-21. Las 22 lecturas públicas y la secuencia práctica de MIT 6.813/6.831 Spring 2016 fueron auditadas de punta a punta. Fuentes de plataforma verificadas: WHATWG HTML/DOM/Fetch/URL; ECMAScript 2026; TypeScript 6.0; W3C CSS/WCAG 2.2/WAI-ARIA 1.2/i18n; Chromium; Web Vitals; React 19.2 y Web Components; Web Platform Tests, Testing Library y Playwright. Versiones y fronteras se fijan en §24.
> **Propósito:** convertir necesidades humanas en interfaces correctas, accesibles, rápidas, seguras, medibles y mantenibles sin confundir un framework con la plataforma web.
> **Relación:** backend conserva negocio/API; seguridad conserva trust/IAM; redes conserva delivery; este manual conserva interacción, estado visible, semántica del navegador y experiencia del usuario.

## 0. Cómo usar este manual

Toda tarea sigue este orden:

```text
persona/tarea/contexto
→ outcome y failure cost
→ content + interaction model
→ semantic/accessibility contract
→ state machine + API contract
→ rendering/performance budget
→ implementation
→ tests with real browser/users
→ field measurement + iteration
```

Etiquetas: `[SPEC]` norma/plataforma; `[ACADEMIC]` HCI/estudio; `[CODE]` implementación oficial; `[PROD]` práctica operativa; `[MEASURED]` medición reproducible; `[IMPL]` traducción propia; `[OPEN]` evidencia faltante.

No declarar “intuitivo”, “rápido”, “responsive”, “accesible” o “seguro” sin tarea, población, dispositivo, test y métrica.

## 1. Contrato de producto antes de componentes

```yaml
product:
  users_and_contexts: []
  jobs_and_critical_paths: []
  outcomes_and_non_goals: []
  harmful_or_irreversible_actions: []
platform:
  browsers_devices_input_modes: []
  locales_writing_modes: []
  connectivity_memory_cpu_constraints: []
interaction:
  information_architecture: {}
  states_and_transitions: {}
  empty_loading_error_offline_partial: []
data:
  sources_authority_freshness: {}
  optimistic_and_reconciliation: {}
quality:
  wcag_target: "2.2 AA unless project says otherwise"
  performance_budgets: {}
  security_privacy: {}
  test_and_field_metrics: []
```

Definir el resultado del usuario, no “hacer una pantalla”. Una interfaz puede renderizar correctamente y fracasar porque oculta el estado, exige memoria innecesaria, no permite recuperación o mide clicks en vez del outcome.

## 2. HCI: diseñar, observar, corregir

**[ACADEMIC]** MIT 6.813/6.831 organiza la disciplina en principios de diseño, task analysis/prototyping/user testing, implementación y métodos de investigación. Su lección central es iterativa: prototipos baratos y feedback temprano antes de pagar el costo del sistema final.

### 2.1 Task analysis

Para cada tarea crítica:

- quién la realiza, con qué conocimiento y frecuencia;
- objetivo real y evidencia de éxito;
- información disponible y que debe recordar;
- secuencia normal, atajos y dependencias;
- errores probables, costo y recuperación;
- entorno: móvil/escritorio, ruido, movimiento, presión, baja conectividad, assistive technology;
- excepciones que hoy resuelve fuera del producto.

Observar trabajo real y artefactos existentes. Entrevistas describen percepciones; logs describen eventos; ninguna fuente por sí sola demuestra el workflow.

### 2.2 Principios operativos

- **learnability:** conceptos, mapping y señales permiten comenzar sin memorizar el sistema;
- **visibility:** estado, progreso, scope y resultado de acciones son observables;
- **efficiency:** usuarios frecuentes reducen pasos sin volver invisible el control;
- **error prevention:** constraints y defaults seguros eliminan clases de error;
- **recovery:** undo, cancel, drafts, confirmation selectiva y mensajes accionables;
- **consistency:** mismo concepto conserva nombre, comportamiento y posición razonable;
- **user control:** no secuestrar foco, navegación, clipboard, zoom o decisiones.

Reconocimiento suele ser menos costoso que recuerdo, pero una UI saturada también impone búsqueda. Progressive disclosure conserva el camino principal y hace descubribles opciones avanzadas.

Tratar la eficiencia humana como propiedad medible, no como estética:

- **Fitts:** para apuntar, `T ≈ a + b·log2(D/W + 1)`; aumentar el target efectivo y acercar acciones frecuentes reduce dificultad. Los coeficientes dependen de dispositivo/población y se calibran, no se copian.
- **Steering:** atravesar un canal estrecho cuesta aproximadamente con `D/W`; evitar menús o gestos que obligan a mantener una trayectoria frágil.
- **KLM/GOMS:** puede comparar procedimientos estables de expertos mediante operadores; no predice descubribilidad, errores, fatiga ni tareas nuevas.
- **chunking:** agrupar información según unidades que el usuario ya reconoce; decorar una cadena arbitraria no crea comprensión.

**[ACADEMIC]** Son modelos para generar hipótesis y encontrar costos, no garantías universales. Una mejora calculada se valida en la tarea y dispositivo reales.

Taxonomía de seguridad de interacción:

| Fallo | Qué ocurre | Respuesta de diseño |
|---|---|---|
| slip | intención correcta, ejecución incorrecta | targets distinguibles, constraints, feedback, undo |
| lapse | se pierde un paso/intención | estado visible, drafts, resumibilidad, checklist contextual |
| mistake | modelo/objetivo equivocado | mejor mapping, explicación, preview, límites del dominio |
| mode error | acción interpretada bajo modo no percibido | modo visible, temporal cuando sea posible, salida clara |

Confirmar repetidamente crea habituación. Para acciones reversibles, preferir undo semántico; para irreversibles, mostrar objeto, alcance y consecuencia exactos. Automatización crítica conserva estado observable y escape/manual override autorizado.

### 2.3 Prototipo y evaluación

```text
sketch/paper → task walkthrough → interactive prototype
→ moderated/unmoderated usability test → implementation slice
→ production telemetry + qualitative follow-up
```

Separar problemas de **concepto**, **flujo**, **contenido**, **feedback**, **motor/percepción**, **accesibilidad** y **performance**. Priorizar por frecuencia, severidad, alcance y costo de recuperación; convertir hallazgos materiales en cambio + regression test o design-system rule.

La fidelidad tiene varias dimensiones: amplitud de tareas, profundidad funcional, apariencia y conducta. Elegir sólo la necesaria para responder la pregunta. Prototipos horizontales comparan arquitectura/recorrido; cortes verticales exponen integración y riesgo técnico. Las primeras alternativas deben ser baratas y desechables: apegarse al código temprano sesga la decisión.

No confundir tres instrumentos:

| Instrumento | Pregunta | Protocolo mínimo |
|---|---|---|
| cognitive walkthrough | ¿un usuario nuevo descubrirá cada acción? | tarea + camino; subgoal, señal, mapping y feedback por paso |
| heuristic evaluation | ¿qué principios/contratos viola la interfaz? | 3–5 evaluadores independientes; consolidar, severidad y debrief |
| formative user test | ¿qué hacen usuarios representativos en tareas representativas? | briefing estable; observar sin enseñar; think-aloud con sesgo reconocido |

Severidad heurística combina frecuencia, impacto y persistencia. Registrar `problema → evidencia → heurística/contrato → severidad → recomendación → captura`. En prototipos estáticos buscar activamente elementos ausentes: una pantalla dibujada no bloquea al evaluador cuando falta una acción real.

### 2.4 Contrato experimental

```yaml
hypothesis: "testable and directional when justified"
population_and_sampling_bias: {}
independent_variables: []
dependent_variables: [task_success, time, errors, recovery, satisfaction]
confounds_and_controls: []
assignment_and_order: "randomized/counterbalanced"
primary_metric_and_guardrails: {}
sample_size_and_stopping: {}
analysis: {effect_size: "", uncertainty: "", assumptions: []}
privacy_ethics: {}
decision_rule: ""
```

Validez interna pregunta causalidad; externa, generalización; reliability, repetibilidad. Controlar o aleatorizar confusores, reportar sesgo de muestra y contrabalancear aprendizaje/fatiga. Visualizar distribuciones/outliers antes del test; distinguir significancia estadística de importancia práctica; justificar supuestos, múltiples comparaciones y abandono de datos.

### 2.5 Cobertura MIT 6.813/6.831 auditada

| Bloque | Lecturas | Transferencia al manual |
|---|---:|---|
| calidad humana | 01–05 | usability contextual; learnability; affordance/feedback; eficiencia; safety/undo |
| arquitectura e interacción | 06, 08–09 | view tree, listener, model/view; output; raw/translated events, dispatch y foco |
| proceso y evidencia | 07, 10–13, 20 | UCD, prototipos, experimentos, user testing/análisis, evaluación heurística |
| diseño perceptual | 14–17 | reducción/regularidad/contraste, grouping, grid, color, tipografía |
| inclusión/contexto | 18–19, 22 | accesibilidad universal, i18n/localización y diseño bajo restricciones reales |
| frontera de investigación | 21 | nuevas técnicas requieren modelo, prototipo y evaluación; premio/demo no equivale a generalidad |

La secuencia práctica `análisis → alternativas → papel → prototipo interactivo → implementación → user test` es el lifecycle de referencia. Los ejercicios históricos de HTML/CSS/JS prueban conceptos, pero no fijan herramientas actuales.

## 3. Modelo mental del navegador

**[SPEC]** ECMAScript define el lenguaje; HTML/DOM/Fetch/URL definen APIs y algoritmos de la plataforma; el browser host define event loops, rendering y políticas. “JavaScript es single-threaded” es insuficiente: documentos, workers, network/GPU/browser processes y task queues cooperan bajo contratos distintos.

### 3.1 Navegación y documento

```text
URL parse/resolve
→ navigation + redirects/policy
→ fetch/cache/service worker/network
→ response/MIME/encoding
→ HTML parse + script execution
→ DOM + styles
→ style/layout/paint/composite
→ pixels + accessibility tree + input
```

Una navegación reemplaza documentos y puede cruzar process/site boundary. History, back/forward cache, form resubmission, scroll restoration y unsaved state son parte del lifecycle, no detalles del router.

La plataforma identifica la navegación en curso y una nueva puede volver obsoleto/abortar trabajo anterior. Router, loader y componente deben compartir esa noción de ownership; nunca permitir que una completion vieja publique título, permisos o datos sobre el documento nuevo. Distinguir navegación same-document, cross-document y history traversal; `pushState` no carga ni valida un recurso remoto.

### 3.2 Event loop y responsividad

- tasks ejecutan callbacks/eventos; microtasks se drenan en checkpoints definidos;
- existen varias task sources/queues y el user agent elige entre queues; sólo se preserva el orden exigido dentro de una misma source, no una FIFO global inventada;
- una función larga bloquea input, style/layout y paint del main thread;
- `async/await` ordena promesas; no crea CPU paralela ni cancela por sí mismo;
- timers expresan mínimo/eligibility, no deadline exacto;
- Web Workers aíslan CPU work, pero serialization/transfer/lifecycle cuesta;
- `requestAnimationFrame` coordina visual update; no usarlo como cronómetro de negocio.

Toda operación asíncrona necesita ownership, cancelación/obsolescencia y regla contra late completion.

### 3.3 Rendering

DOM mutation no equivale a pixel inmediato. Distinguir:

- style recalculation;
- layout y dependencias geométricas;
- paint de display items;
- raster/compositing;
- presentation al display.

Evitar ciclos read-layout/write repetidos; agrupar cambios y medir. `transform/opacity` puede favorecer composición, pero layers consumen memoria y no eliminan main-thread work de handlers o layout previo.

**[PROD]** Chromium usa múltiples procesos y sandbox/site isolation. Renderer, browser, network y GPU poseen permisos/fallos distintos; iframe cross-site puede vivir fuera de proceso. No asumir memoria, event loop o DOM compartido entre frames/origins.

### 3.4 Input, foco y controladores

Un input físico produce eventos crudos; la plataforma los traduce, encola y despacha. Conservar el timestamp de ocurrencia cuando la latencia importa: el momento del handler no equivale al momento del gesto. Coalescing puede ser correcto para posición visual y fatal para un historial/auditoría; declarar cuál semántica se necesita.

```text
device/browser state
→ raw event
→ translated event + occurrence time
→ queue/coalescing
→ capture → target → bubble
→ controller state transition
→ model/view update + feedback
```

Gestos, drag, autocomplete, menus y composición de texto son máquinas de estados, no colecciones de callbacks. Definir cancelación, pointer capture, focus ownership, keyboard parity, composition/IME y cleanup al desaparecer el target. `preventDefault` y propagation stops requieren una razón de contrato: pueden romper navegación, scrolling, shortcuts o tecnología asistiva.

### 3.5 Matriz de capability web

Antes de adoptar una API/CSS feature:

```yaml
capability: ""
normative_source_and_status: "Living Standard | Recommendation | Draft"
target_browser_engine_versions: []
interoperability_evidence: "WPT/vendor tests/results"
secure_context_permission_policy: {}
accessibility_semantics: {}
failure_and_fallback: {}
polyfill_cost_and_trust: {}
real_browser_tests: []
removal_or_migration_trigger: ""
```

Un Working Draft no se presenta como Recommendation. “Soportado” no basta: comprobar semántica, flags, permisos, mobile/webview y combinación con assistive technology. Polyfill sólo puede reproducir lo observable desde JavaScript; no inventa integración nativa, aislamiento o accessibility API que el motor no expone.

## 4. HTML semántico primero

Elegir elementos por significado y comportamiento nativo:

- landmarks/header/nav/main/footer para estructura;
- heading levels describen jerarquía, no tamaño visual;
- `button` ejecuta acción; `a[href]` navega;
- `label` identifica control; `fieldset/legend` agrupa;
- lists/tables sólo para relaciones reales;
- `dialog`, `details`, inputs y validation se usan según soporte/contrato.

No reemplazar un control nativo por `div` + handlers salvo requerimiento demostrado. Recrearlo exige teclado, foco, nombre/rol/valor, disabled, high contrast, touch y assistive technology que el browser ya implementa.

### 4.1 Formularios y acciones

- labels persistentes; placeholder no sustituye nombre;
- formato, unidades y restricciones visibles antes del error;
- validación cliente mejora feedback, nunca reemplaza servidor;
- preservar input ante error recuperable;
- error se asocia al campo y existe también en resumen cuando procede;
- prevenir doble submit mediante estado e idempotencia backend;
- confirmación sólo donde riesgo/irreversibilidad la justifica; preferir undo cuando sea seguro.

`type`, `autocomplete` e `inputmode` son contratos distintos: control/validación, significado del dato y modalidad sugerida. Declarar los tres cuando proceda. Usar submit nativo y progressive enhancement conserva teclado, autofill y funcionamiento básico; interceptarlo exige preservar successful controls, encoding, validation, history y error recovery.

## 5. CSS, layout y sistema visual

El cascade tiene origen/importancia, layers, specificity, scope/order. Diseñar una política en lugar de ganar guerras con `!important`.

```text
tokens → reset/base → semantic primitives → components
→ compositions/layout → utilities → controlled overrides
```

- Grid para relaciones bidimensionales; Flexbox para distribución en un eje;
- intrinsic sizing y content wrapping antes de breakpoints arbitrarios;
- container/media queries según dependencia real;
- logical properties para writing modes;
- reservar dimensiones de media/contenido asíncrono para estabilidad;
- respetar zoom, font scaling, reduced motion, contrast y color schemes;
- nunca codificar significado sólo con color, posición o animación.

Un responsive layout conserva contenido, orden semántico y acciones críticas. CSS visual reordering no debe contradecir DOM/focus order.

### 5.1 Tokens y componentes

Tokens semánticos (`surface-danger`, `space-control-inline`) expresan intención; valores crudos son implementación. Cada componente documenta:

```yaml
anatomy: []
states: [default, hover, focus, active, disabled, loading, error, empty]
keyboard_and_focus: {}
content_constraints: {}
responsive_behavior: {}
accessibility_contract: {}
performance_cost: {}
tests: []
```

### 5.2 Jerarquía visual verificable

- **reducción:** cada elemento/variación necesita una función; probar si retirarlo destruye significado o acción;
- **regularidad:** poca variación incidental hace que la diferencia intencional sea visible;
- **contraste:** jerarquía y estado deben sobrevivir blur/squint, grayscale, zoom y temas;
- **grouping:** proximidad, whitespace y contención expresan relaciones antes que bordes decorativos;
- **alignment/grid:** minimizar ejes arbitrarios y alinear baselines; responsive reorganiza, no sólo escala;
- **typography:** medir legibilidad/readability con contenido real; limitar line length y variantes según contexto, no por cifra mágica;
- **color:** nunca canal único; combinar texto/forma/posición. Probar deficiencias de color y contrastes definidos por WCAG.

Mood y credibilidad son requisitos de producto, no excusa para degradar lectura o simular confianza. Ningún token visual convierte por sí mismo una jerarquía de información defectuosa en usable.

## 6. JavaScript y TypeScript como contratos parciales

TypeScript elimina clases de inconsistencias antes de ejecutar, pero no valida JSON, DOM, storage ni server responses. Toda frontera externa entra como `unknown` y se valida.

Reglas base:

- `strict`; evitar `any` y assertions sin prueba;
- discriminated unions para estados legales;
- branded/opaque types donde IDs/unidades se confunden;
- readonly/immutable data cuando simplifica ownership;
- exhaustive checks para state machines;
- separar domain value de wire/view model;
- fechas, moneda y decimal no son `number` sin contrato;
- generated API types no sustituyen runtime schema ni compatibility tests.

```ts
type Remote<T> =
  | { tag: "idle" }
  | { tag: "loading"; requestId: string }
  | { tag: "ready"; value: T; receivedAt: number }
  | { tag: "empty"; receivedAt: number }
  | { tag: "error"; error: UiError; retryable: boolean };
```

Esto impide combinaciones como `loading=true` y `data/error` incompatibles, siempre que el runtime input también se valide.

### 6.1 Toolchain, dependencias y release

- fijar runtime, package manager, lockfile, TypeScript, bundler y browser targets; CI usa instalación frozen/reproducible;
- separar source modules de chunks emitidos; revisar graph, duplicate dependencies, dynamic imports, CSS/assets y source maps;
- minificación/tree shaking no demuestran ausencia de side effects: verificar output y comportamiento de producción;
- dependency update incluye changelog/security/provenance, type/runtime tests, bundle diff y rollback; no mezclar renovación masiva con feature crítica;
- variables públicas se consideran datos del cliente aunque provengan del entorno de build; secretos sólo en frontera servidor;
- promover el mismo artifact/digest entre ambientes cuando sea práctico; config runtime se valida y versiona;
- sourcemaps privadas conservan correspondencia con release y acceso/retention acotados.

Budget de entrega separa `HTML/CSS/JS/fonts/media`, parse/compile/evaluate y código por interacción. Una reducción gzip que aumenta ejecución o memoria no es victoria.

### 6.2 Código industrial auditado: Chromium, React y TypeScript

**[PROD]** Segunda pasada fijada a Chromium `e5b0915cd04c456a7695009db757e65de7558692` (licencia BSD de tres cláusulas), React `eafeac097ba51e1eab809c07102126bd5f8e5425` (MIT) y TypeScript `d6c4afddb2c55f4a9dea7b59293a99a8fdea1799` (Apache-2.0). Son evidencia pública de implementación y tests en esos commits; no prueban cómo está configurado cada producto propietario ni obligan a copiar su arquitectura.

**Chromium — secuencia, thread y lifetime no son sinónimos.** Su modelo favorece secuencias de tareas y message passing, evita trabajo caro o I/O bloqueante en UI/IPC y asocia estado mutable a una secuencia comprobable. Un callback tardío debe tener dueño vivo; ownership único, weak references/cancelación y shutdown explícito son parte del contrato. Transferencia web: no compartir estado mutable entre main thread, worker, iframe o proceso como si fueran uno; definir mensaje, orden, backpressure, obsolescencia y cierre. Mover CPU a un worker no corrige por sí solo serialization cost, flooding ni publicación de resultados viejos.

**React — calcular puede interrumpirse; publicar no.** El scheduler mantiene por separado tareas disponibles y timers, asigna prioridad/expiración, cede al host y permite continuaciones. Fiber distingue render, suspensión/reintento y commit; por tanto un componente no puede depender de que cada render tentativo llegue a publicarse. Los tests de Strict Effects ejercitan `mount → cleanup → remount` para exponer efectos sin cleanup o con side effects no repetibles. Contrato transferible: render puro; identidad/key estable; efecto reversible e idempotente respecto de reintentos; async work cancelable/obsoleto; recurso externo adquirido y liberado simétricamente. Medir el commit y la experiencia visible, no contar renders como resultado de negocio.

**TypeScript — incremental es una optimización con obligación de equivalencia.** El compilador nativo auditado conserva snapshots, firmas, diagnostics y mapas de referencias para determinar archivos afectados. Sus tests aceptan fast path para cambios internos, exigen full rebuild cuando cambia un import/dependencia y fuerzan reconstrucción total si la cola del watcher desborda; también ejecutan ciclos concurrentes bajo detector de races. Transferencia: cache de build, HMR, generated types y project references nunca son autoridad. Renames, deletes, cambios de config/package/export y pérdida/coalescing de eventos deben invalidar de forma conservadora; CI compara periódicamente salida/diagnostics incrementales contra clean build. Aun con build correcto, input de runtime continúa necesitando validación.

Gate industrial para una ruta crítica:

1. Perfilar main thread, workers, long tasks, memoria, red y paint en hardware/OS representativos; un promedio local no cierra el gate.
2. Para cada tarea/callback: owner, secuencia, prioridad/deadline, cancelación, late-completion rule y shutdown test.
3. Para cada render/efecto: pureza, cleanup, repetición, interrupción, error y unmount durante I/O.
4. Para el toolchain: instalación frozen, clean build reproducible y prueba incremental tras edit/import/rename/delete/config change y watcher overflow simulado.
5. Una regresión se describe por benchmark, story, métrica, dispositivo y coste absoluto/relativo; reproducir, bisectar y trazar antes de justificarla. Si no existe ganancia demostrada o trade-off explícito de correctness/security, corregir o revertir.

## 7. Estado, datos y concurrencia visible

Clasificar estado:

| Tipo | Autoridad | Ejemplos |
|---|---|---|
| URL/navigation | browser/router | recurso, filtros compartibles, history |
| server state | API/stream | entidades, permisos, freshness |
| persisted client | storage | preferencia no sensible, draft |
| session UI | component/store | selection, disclosure, pending action |
| derived | función | totals, filtered view, validation |

No duplicar derived state. Evitar un global store por defecto: scope mínimo reduce invalidation, races y acoplamiento.

### 7.1 Fetch state machine

```text
IDLE → PENDING(request_id)
PENDING → SUCCESS | EMPTY | ERROR | CANCELED | SUPERSEDED
SUCCESS → REFRESHING(previous_value)
REFRESHING → SUCCESS | ERROR_WITH_STALE_VALUE
```

- AbortSignal cancela trabajo local/request según API, no deshace un commit remoto;
- response tardía sólo aplica si request/version sigue siendo current;
- optimistic UI necesita client operation ID, rollback/reconcile y conflicto visible;
- retry respeta método/idempotencia, deadline, connectivity y server hints;
- loading skeleton no falsifica estructura/avance ni mueve layout.

## 8. Arquitectura y estrategias de rendering

Elegir por ruta, no por dogma:

| Estrategia | Beneficio | Costo/riesgo |
|---|---|---|
| static/prerender | caché, poco runtime, primera vista | freshness/build invalidation |
| SSR | contenido inicial/metadata personalizados | server cost, hydration y dual execution |
| streaming SSR | progressive delivery | ordering, fallback y observability complejos |
| CSR | interacción rica/hosting simple | JS/bootstrap/data waterfall |
| islands/progressive enhancement | poco JS por zona | boundaries/tool constraints |

Hydration exige que server/client produzcan estructura compatible bajo locale, time, randomness, permissions y data version. Un mismatch no se “silencia”: localizar nondeterminism.

Separar:

```text
domain/use cases
↕ ports/contracts
data adapters + cache
↕ view models
semantic components + routes
```

El framework organiza render/lifecycle; no posee la verdad de API, accesibilidad, seguridad o producto.

MVC es una separación conceptual, no una obligación de clases. El modelo conserva datos e invariantes; la vista produce output; el controlador traduce input. En componentes modernos view/controller suelen estar acoplados, pero el dominio no debe absorber selección, hover, foco o disclosure efímeros. Dos vistas del mismo modelo pueden necesitar estado de interacción independiente.

### 8.1 Ruta, datos, caché y metadata

Cada ruta declara:

```yaml
identity_and_canonical_url: {}
authorization_and_data_dependencies: []
render_strategy_and_fallback: ""
cache_scope_and_revalidation: {}
metadata_and_structured_content: {}
loading_error_not_found_boundaries: {}
navigation_scroll_focus_behavior: {}
stream_hydration_client_boundaries: []
```

No cachear una respuesta personalizada bajo key pública ni compartir markup entre tenants. `Cache-Control`, CDN, server cache, framework data cache y client cache son capas distintas. Metadata/title/canonical/language/social cards y contenido indexable se generan desde la misma autoridad de ruta; SEO no justifica contenido diferente engañoso ni sacrificar semántica/accesibilidad.

## 9. Accesibilidad como corrección

**[SPEC]** WCAG 2.2 estructura criterios testables bajo perceivable, operable, understandable y robust. Objetivo general: AA completo por página/variación, salvo requisito distinto documentado. WAI-ARIA 1.2 es Recommendation; 1.3 sigue Working Draft al 2026-08-21.

Gate mínimo:

- HTML semántico y accessible name correcto;
- operación completa por teclado sin trap;
- foco visible, no oculto y restaurado tras modal/navigation;
- lectura/orden lógico y headings/landmarks coherentes;
- labels, instructions y errors programmatically associated;
- contrast, zoom/reflow, target size y motion preferences;
- status/error dinámico anunciado sin spam;
- media alternatives según contenido;
- auth no depende sólo de puzzle cognitivo;
- probar browser + screen reader representativos, keyboard-only, zoom y high contrast.

Diseñar universalmente desde el principio: capacidad visual, auditiva, motora y cognitiva varía entre personas, a lo largo del tiempo y por situación. Una alternativa equivalente no debe degradar o estigmatizar. Probar al menos salida alternativa, input alternativo y pérdida de un canal sensorial; captions, keyboard y estructura semántica benefician también contextos temporales/ruidosos.

ARIA cambia el accessibility tree, no añade conducta. `role="button"` no crea keyboard activation, focus, form behavior ni disabled semantics. First rule: usar host-language semantics cuando existen.

Automated scanners detectan una fracción. Conformance necesita revisión manual y, para productos materiales, usuarios con discapacidades.

WCAG conforma **páginas completas y procesos completos**, no una colección de componentes aislados. A nivel AA deben cumplirse todos los criterios A+AA aplicables, usar tecnologías accessibility-supported y respetar non-interference. En 2.2 probar explícitamente, entre otros, foco no totalmente oculto, alternativa a dragging, target mínimo de 24×24 CSS px con sus excepciones, entrada redundante y autenticación sin prueba cognitiva salvo alternativa/mecanismo admitido. Documentar criterio, técnica, evidencia, browser/AT y excepción; un badge o scanner verde no es una conformance claim.

## 10. Rendimiento centrado en experiencia

**[PROD]** Core Web Vitals vigentes: LCP ≤2,5 s, INP ≤200 ms y CLS ≤0,1 en p75, segmentado móvil/escritorio. Son gates generales, no definición completa de UX.

Presupuesto por ruta:

```yaml
field: {LCP_p75: "", INP_p75: "", CLS_p75: ""}
network: {html: "", css: "", js_initial: "", images_fonts: ""}
main_thread: {long_tasks: "", hydration: "", interaction_work: ""}
memory: {steady: "", navigation_growth: ""}
product: {time_to_task_success: "", error_recovery: ""}
```

### 10.1 Diagnóstico

```text
slow discovery/navigation?
→ network waterfall/cache/server
slow visual completion?
→ LCP resource, CSS, fonts, render blocking
slow interaction?
→ input delay + handler + render/presentation
unstable?
→ unsized/late content, fonts, insertion
degrades over time?
→ listeners, detached DOM, cache/store, media/GPU memory
```

- field/RUM explica usuarios reales; lab reproduce y atribuye;
- segmentar por ruta, release, device, network, geography y visibility;
- medir p75 y tails, no sólo Lighthouse score o desktop developer machine;
- code splitting puede reducir inicial y crear waterfalls/interaction delay; medir camino completo;
- virtualización reduce DOM, pero debe conservar teclado, focus, semantics y scroll anchoring.

## 11. Seguridad y privacidad de frontend

El browser aplica same-origin policy, CORS, cookie rules, secure contexts y isolation; no confundir una restricción de lectura con autorización servidor.

- XSS: contextual output encoding, safe DOM APIs, sanitización para HTML autorizado y CSP/Trusted Types cuando aplica;
- CSRF: cookie attributes + anti-CSRF/origin policy según arquitectura; CORS no es defensa general;
- auth: tokens/secrets fuera de bundles/URLs/logs; server verifica cada acción;
- clickjacking/embedding: CSP `frame-ancestors`/headers;
- third parties: inventario, integrity/provenance, permissions, egress y failure budget;
- `postMessage`: target origin exacto, validar `event.origin/source/data`;
- storage/cache/service worker: clasificación, expiry, logout/tenant boundary y migration;
- DOM/display no vuelve confiable contenido API/markdown/user;
- source maps, telemetry y error payloads no filtran PII/secrets.

Modelo DOM-XSS: `source → transformations → injection sink`. Preferir `textContent`, attributes tipados y creación DOM; centralizar/someter a revisión cualquier parsing de HTML, script o URL ejecutable. Trusted Types puede reducir sinks a valores creados por policies controladas, pero al 2026-08-21 su especificación publicada sigue Working Draft y necesita matriz de soporte/fallback. Una policy que retorna input sin sanitizar sólo tipa la vulnerabilidad.

CSP es defensa en profundidad, no sustituto de encoding/sanitización. Desplegar primero `Report-Only`, clasificar ruido y filtración de reportes, luego enforcement; evitar nonces reutilizados y aperturas `unsafe-*` sin riesgo documentado. `frame-ancestors`, `form-action`, `base-uri`, sources de script/worker/connect y reporting tienen propósitos distintos.

CORS decide si una respuesta puede compartirse con el caller cross-origin; no autentica ni autoriza la operación. Requests simples/form submissions pueden causar efectos sin preflight. Con credentials, origin/reflection/cache (`Vary: Origin`) y cookies se auditan juntos. Permissions Policy reduce capabilities disponibles para documento/iframes; no convierte contenido third-party en confiable.

Minimizar colección. Consentimiento no legitima cualquier tracking; declarar propósito, retention, acceso y delete/export. Analytics debe sobrevivir ad/script blocking sin romper el producto.

## 12. Internacionalización y contenido

- UTF-8 end-to-end y `lang` del documento/cambios;
- locale data mediante `Intl`, no listas caseras;
- almacenar instants/time zones/calendars y formatear en el borde;
- nombres, direcciones, plurales y orden no se fuerzan a un molde anglocéntrico;
- logical CSS y `dir`; aislar texto bidi dinámico;
- no concatenar fragmentos traducibles;
- layout resiste expansión, fuentes/fallback y mixed scripts;
- rutas, metadata, search y content negotiation declaran estrategia locale.

Collation no es ordenar code points; case conversion tampoco es universal. Mensajes dinámicos conservan frase completa y placeholders nombrados para permitir reordenamiento y plural rules. Formato de número, moneda, fecha y zona depende de locale/contexto; parsing de input debe declarar ambigüedad. Probar mixed-direction strings, nombres largos, grapheme clusters, IME, line breaking y truncation sin perder la acción.

Traducción lingüística, localización cultural y formato técnico son problemas relacionados, no idénticos.

## 13. Sistemas de diseño sin congelar el producto

Un design system es contratos + decisiones + governance, no sólo un catálogo visual:

```text
foundations/tokens
→ semantic primitives
→ accessible components
→ patterns/workflows
→ templates/examples
→ lint/tests/release/migration
```

Cada breaking visual/semantic change tiene changelog, codemod/migration y visual/accessibility regression. Evitar un componente universal con decenas de boolean props; composición y variantes tipadas mantienen estados legales.

## 14. Testing por capas

| Capa | Prueba | Detecta |
|---|---|---|
| pure/domain | unit/property | transformación, invariantes, bordes |
| component | DOM realista + user events | semantics, states, callbacks |
| contract | schema/API fixtures | drift, error/permission shapes |
| integration | route/data/cache | races, optimistic/reconcile |
| E2E real browser | critical journeys | wiring, navigation, browser behavior |
| accessibility | static + manual + AT | WCAG/interaction real |
| visual | controlled screenshots | layout/theme/regression |
| performance | lab + RUM | budget y release regression |
| usability | representative users/tasks | problema de diseño/producto |

Testing Library prioriza interacción observable sobre component internals. Playwright locators por role/label/text y auto-wait/actionability reducen races; `force`, sleeps y selectors estructurales suelen ocultar un contrato defectuoso.

Matriz E2E crítica:

```text
happy path
empty/partial/stale/error/offline
slow API + late response + cancel/navigation
double action + retry + ambiguous backend outcome
expired auth/permission change
keyboard + focus + screen-size/zoom
locale/RTL/long strings
back/forward/reload/deep link
```

Web Platform Tests validan interoperabilidad de la plataforma; no prueban el producto. Ejecutar browser matrix según usuarios/evidencia, no “sólo Chrome porque funciona”.

## 15. Observabilidad y experimentos

Eventos expresan outcomes/estado, no detalles accidentales del DOM:

```yaml
event: "checkout_completed"
schema_version: 1
operation_id: "pseudonymous/allowed"
route_release_experiment: {}
duration_and_result: {}
error_taxonomy: ""
privacy_classification: ""
```

- métricas de disponibilidad y JS errors por ruta/release;
- Web Vitals y long tasks con attribution controlada;
- task completion, abandonment y recovery, no vanity clicks;
- coverage/denominator y sample bias explícitos;
- session replay sólo con base legal, minimización y redaction probada;
- alertas ligadas a user harm y rollback.

A/B test necesita hipótesis previa, primary/guardrail metrics, unidad de randomización, sample-size/stopping policy y análisis de novelty/segmentos. Significancia sin tamaño de efecto, costo, accesibilidad o ética no decide producto.

Usability test y A/B test no son intercambiables: el primero explica problemas y modelos mentales en una muestra pequeña; el segundo estima diferencias agregadas bajo instrumentación y asignación correctas. Logs no muestran la intención ausente ni a usuarios que abandonaron antes de ser medidos.

## 16. Frameworks: adapter, no fundamento

Evaluar:

```text
platform/browser support
rendering/deployment needs
accessibility primitives
bundle/runtime/memory
data/cache/forms/routing semantics
testing/debugging/observability
upgrade/security cadence
team ecosystem and exit cost
```

React 19.2 es una referencia vigente al auditar, no default universal. Vue/Angular/Svelte/Web Components u otros pueden satisfacer contratos distintos. Fijar versión, renderer/meta-framework y server/client boundary; no mezclar patrones de documentación de releases incompatibles.

| Enfoque | Unidad/lifecycle | Fortaleza | Riesgo que exige prueba |
|---|---|---|---|
| plataforma + progressive enhancement | documento/form/custom behavior | poco runtime y semántica nativa | disciplina manual de composición/state |
| Web Components | custom element + lifecycle/attributes + opcional Shadow DOM | interoperabilidad de elemento/encapsulado | upgrade timing, retargeting, forms, SSR y a11y cross-shadow |
| React 19.2 + renderer/framework fijado | render/commit, state identity, effects, server/client boundaries | composición y ecosystem | redundant state/effects, hydration, bundle, framework-specific RSC contract |

Custom Elements pueden existir antes de su definición y luego ser upgraded; callbacks/reactions tienen ordering específico y no una garantía global entre elementos. Shadow DOM altera retargeting, style y traversal: encapsular CSS no justifica ocultar label, focus o estados a tests/AT.

Reglas transferibles en un sistema reactivo:

- render es cálculo; effects sincronizan con sistemas externos;
- state identity depende de posición/key/ownership definidos;
- derived values no requieren effect + state duplicado;
- cleanup debe tolerar mount/unmount/retry;
- concurrency/transition no corrige side effects impuros;
- server-only data/secrets nunca cruzan serialization boundary accidentalmente.

En React, effect es escape hatch para sincronizar un sistema externo, no un lugar genérico para derivar datos o responder a clicks. Las APIs de Server Components consumidas por bundlers/frameworks no comparten necesariamente la estabilidad semver del API de aplicación; pin exacto, integración oficial y upgrade tests son parte del contrato.

## 17. Integración API, realtime y offline

- UI modela HTTP/RPC error taxonomy, deadlines y ambiguous outcomes del backend;
- WebSocket reconnect no garantiza continuity: resubscribe, epoch/sequence/resync según protocolo;
- backpressure incluye render frequency; coalescing conserva estado final y audit requerido;
- offline queue sólo para operaciones con idempotencia/conflict policy demostrada;
- service worker versiona caches y activation/migration; no servir shell/data incompatible;
- connectivity signal no prueba reachability de la dependencia;
- datos sensibles y acciones peligrosas no se vuelven “offline-first” por preferencia estética.

Para dashboards de order book, `MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md` conserva sequence/state/freshness. La UI muestra `LIVE/STALE/GAPPED/RESYNCING`, no suaviza huecos ni presenta un mid viejo como actual.

## 18. Gates de producción

1. Tareas/outcomes y no-objetivos comprobados.
2. State machine cubre loading/empty/error/offline/permission/late response.
3. URL/history/deep-link/reload preservan contrato.
4. Semántica HTML, keyboard/focus y WCAG target probados.
5. Runtime inputs validados y estados ilegales no representables cuando sea práctico.
6. API auth/idempotency/reconcile y error text son correctos.
7. Responsive, zoom, locales, RTL y contenido extremo sobreviven.
8. Security/privacy gate y third-party inventory cerrados.
9. Browser matrix y critical E2E verdes sin sleeps arbitrarios.
10. Field/lab performance dentro de budget, incluidas tails.
11. Telemetry explica release/fallo sin secretos.
12. Rollout, feature flag, compatibility y rollback probados.
13. Usability test o evidencia equivalente cubre tareas críticas.
14. Hallazgos de error se convierten en prevención sistémica.

## 19. Contrato para Codex

```md
Usuarios, tareas y contexto:
Outcome y no-objetivos:
Rutas y estados:
API/schema/authority/freshness:
Browser/device/locale/AT matrix:
WCAG y keyboard/focus contract:
Performance budgets:
Security/privacy/trust boundaries:
Framework/tool versions:
Files/components in scope:
Tests and evidence required:
Rollout/rollback:
Open questions [OPEN]:
```

El agente debe inspeccionar design system, routes, API contracts, analytics y tests existentes; no inventar UI desde nombres de endpoints. Debe implementar todos los estados, preferir semántica nativa, validar entradas, cancelar/ignorar completions obsoletas, probar desde comportamiento de usuario y reportar evidencia/risks.

Entrega:

```text
task model and state machine:
information/interaction decisions:
files and boundaries:
accessibility:
data races/reconciliation:
security/privacy:
responsive/i18n:
tests and browser evidence:
performance measurement:
telemetry:
open risks and rollback:
```

## 20. Ruta práctica profesional

1. **Documento semántico:** contenido real, headings/landmarks/forms, keyboard y responsive sin framework.
2. **Interaction slice:** state machine, async cancellation, errors/recovery y component tests.
3. **Design system:** tokens + 5 primitives complejos accesibles; Story/examples y regression.
4. **Data product:** search/filter/pagination/cache/optimistic reconcile con contract tests.
5. **Realtime dashboard:** bounded updates, stale/gap states, virtualization accesible y RUM.
6. **SSR/offline:** hydration determinista, cache versioning, slow/offline/reload/back-forward.
7. **Capstone:** multi-locale WCAG 2.2 AA, security/privacy review, field budgets, experiment y rollback.

## 21. Suites ejecutables

| Suite | Gate |
|---|---|
| `semantic_a11y` | roles/names/labels/order/focus/keyboard + manual AT matrix |
| `ui_state_model` | todas las transiciones, late/duplicate/cancel y estados imposibles |
| `api_contract` | schema/version/errors/auth/ambiguous outcome |
| `navigation` | deep link, history, reload, bfcache-sensitive state |
| `responsive_i18n` | narrow/wide, zoom, long text, RTL, fonts/timezones |
| `browser_e2e` | journeys críticos en matrix fijada, sin sleep/force injustificado |
| `security_privacy` | XSS/CSRF/embed/message/storage/third parties/redaction |
| `performance` | Web Vitals field + lab attribution + bundle/memory budgets |
| `resilience` | offline/slow/partial/outage/update/service-worker rollback |
| `usability` | task success, time/errors/recovery + hallazgos cualitativos |

## 22. Auditoría transversal inicial

| Manual | Autoridad de frontera | Obligación frontend |
|---|---|---|
| `SOFTWARE_BACKEND_API_ENGINEERING.md` | negocio/API/auth/idempotencia | UI preserva lifecycle, error, deadline y reconcile |
| `NETWORKING_DISTRIBUTED_STREAMING.md` | transport/delivery/gap/backpressure | UI no inventa continuity y limita render/update |
| `DATABASE_STORAGE_INTERNALS.md` | authority/durability/query | cache cliente no se llama source of truth sin contrato |
| `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md` | CPU/memory/profiling/tails | medir main thread/memory/end-to-end en hardware real |
| `GPU_ACCELERATED_COMPUTING.md` | device correctness | WebGPU/canvas/media conservan lifecycle y fallback |
| `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md` | IAM/privacy/supply chain/SLO | browser no sustituye server auth; telemetry/third parties gobernados |
| manuales de IA | model/data/eval/agents | UI expresa incertidumbre, provenance, permisos y human control |
| `MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md` | book/order/risk/freshness | dashboard muestra state/sequence; trading UI no salta risk gate |
| `ALGORITHMS_DATA_STRUCTURES_PROBLEM_SOLVING.md` | estructura/corrección/complejidad | estado, búsqueda, ranking y virtualización preservan semántica y budgets |
| `NATIVE_MOBILE_DESKTOP_ENGINEERING.md` | lifecycle/ventanas/inputs/permisos/distribución del OS | interacción conserva state, accesibilidad y conventions por plataforma |

Cruces recíprocos cerrados el 2026-08-21: cada manual de la tabla contiene ahora su obligación frontend. Para algoritmos, DOM/state indexes, search, top-`k`, diff, virtualización y worker queues se eligen con invariantes, peor caso y memoria; medir main thread y experiencia real. Para `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md`, context/runtime views incluyen persona/browser/edge/API y estados offline/partial/stale/unknown; frontend decide interacción/accesibilidad y arquitectura decide boundaries/quality trade-offs. Para `DATA_ENGINEERING_ANALYTICS.md`, dashboard/visualización porta metric version, grain, filtros, timezone, freshness, completeness, uncertainty y lineage accesible; cache/export no cambia authority ni access policy. Para `TOOLCHAINS_BUILDS_PACKAGING_FFI.md`, Node/package manager/lock/browser targets/bundler/minifier/source maps forman artifact identity; testear contenido publicado, CSP/SRI y rollback del mismo digest. Para `NATIVE_MOBILE_DESKTOP_ENGINEERING.md`, este manual conserva task analysis, interaction, visual hierarchy y accessibility; nativo gobierna lifecycle, windows/scenes, hardware inputs, permissions, background y packaging. Una UI compartida debe respetar platform conventions y demostrar journeys con assistive technologies reales. Optimizar render no permite perder identidad estable, orden, foco, accesibilidad, freshness ni resultados. Regla de arbitraje: este manual decide interacción y estado visible; el manual de dominio decide verdad, autorización, delivery, durability, compute y riesgo. La UI no corrige silenciosamente una contradicción aguas abajo: la representa, bloquea la acción peligrosa y conserva evidencia.

## 23. Extensiones condicionadas

- ~~Auditar MIT 6.813/6.831 lectura por lectura~~: cerrado el 2026-08-21 sobre sus 22 lecturas públicas y secuencia práctica; ejemplos históricos se conservaron sólo como principios transferibles.
- ~~Completar inventario normativo WHATWG/W3C y browser-interoperability por capability~~: cerrado para el núcleo mediante matriz §3.5; capabilities especiales se fijan por proyecto.
- ~~Auditar React 19.2 y al menos un enfoque alternativo~~: cerrado para principios transferibles contra plataforma/progressive enhancement y Web Components; la segunda pasada añadió scheduler/Fiber/Strict Effects fijados por commit. Meta-framework específico se fija por proyecto.
- Profundizar WebGPU, Canvas, media, WebRTC y gráficos sólo si el producto los necesita.
- ~~Cerrar cruces recíprocos con las fronteras materiales~~: cerrado el 2026-08-21 e incorporadas también las capas posteriores aplicables.
- ~~Native mobile/desktop~~: núcleo general cerrado en `NATIVE_MOBILE_DESKTOP_ENGINEERING.md`; branding/visual craft avanzado, product management y research especializado requieren módulos si el puesto/proyecto los exige.

## 24. Fuentes primarias iniciales

### Plataforma y lenguaje

- WHATWG HTML: <https://html.spec.whatwg.org/>.
- WHATWG DOM: <https://dom.spec.whatwg.org/>.
- WHATWG Fetch: <https://fetch.spec.whatwg.org/>.
- WHATWG URL: <https://url.spec.whatwg.org/>.
- ECMAScript 2026: <https://tc39.es/ecma262/2026/multipage/>.
- TypeScript 6.0 release notes y handbook: <https://www.typescriptlang.org/docs/handbook/release-notes/typescript-6-0.html>, <https://www.typescriptlang.org/docs/handbook/intro.html>.
- CSS specifications: <https://www.w3.org/Style/CSS/specs.en.html>.
- Web Platform Tests: <https://web-platform-tests.org/>.
- Custom Elements y Shadow DOM/event dispatch: <https://html.spec.whatwg.org/multipage/custom-elements.html>, <https://dom.spec.whatwg.org/#shadow-trees>.

### HCI, accesibilidad e internacionalización

- MIT 6.813/6.831 Spring 2016, índice de las 22 lecturas, ejercicios y proyecto: <https://web.mit.edu/6.813/www/sp16/>.
- MIT OCW 6.831: <https://ocw.mit.edu/courses/6-831-user-interface-design-and-implementation-spring-2011/>.
- WCAG 2.2 Recommendation: <https://www.w3.org/TR/WCAG22/>.
- WAI-ARIA 1.2 Recommendation: <https://www.w3.org/TR/wai-aria-1.2/>.
- ARIA Authoring Practices Guide: <https://www.w3.org/WAI/ARIA/apg/>.
- W3C Internationalization: <https://www.w3.org/International/quicktips/>.

### Browser, rendimiento e implementación

- Chromium design documents: <https://www.chromium.org/developers/design-documents/>.
- Chromium fijado `e5b0915cd04c456a7695009db757e65de7558692`: [threading/tasks](https://github.com/chromium/chromium/blob/e5b0915cd04c456a7695009db757e65de7558692/docs/threading_and_tasks.md), [SiteInstance](https://github.com/chromium/chromium/blob/e5b0915cd04c456a7695009db757e65de7558692/content/browser/site_instance_impl.cc), [sandbox](https://github.com/chromium/chromium/blob/e5b0915cd04c456a7695009db757e65de7558692/docs/design/sandbox.md) y [tratamiento de regresiones](https://github.com/chromium/chromium/blob/e5b0915cd04c456a7695009db757e65de7558692/docs/speed/addressing_performance_regressions.md).
- Web Vitals: <https://web.dev/articles/vitals>.
- React versions/docs: <https://react.dev/versions>, <https://react.dev/learn>.
- React effects y Server Components: <https://react.dev/learn/you-might-not-need-an-effect>, <https://react.dev/reference/rsc/server-components>.
- React fijado `eafeac097ba51e1eab809c07102126bd5f8e5425`: [scheduler](https://github.com/facebook/react/blob/eafeac097ba51e1eab809c07102126bd5f8e5425/packages/scheduler/src/forks/Scheduler.js), [Fiber work loop](https://github.com/facebook/react/blob/eafeac097ba51e1eab809c07102126bd5f8e5425/packages/react-reconciler/src/ReactFiberWorkLoop.js) y [Strict Effects tests](https://github.com/facebook/react/blob/eafeac097ba51e1eab809c07102126bd5f8e5425/packages/react-reconciler/src/__tests__/StrictEffectsMode-test.js).
- TypeScript fijado `d6c4afddb2c55f4a9dea7b59293a99a8fdea1799`: [incremental program](https://github.com/microsoft/TypeScript/blob/d6c4afddb2c55f4a9dea7b59293a99a8fdea1799/tsc/internal/execute/incremental/program.go), [reference map](https://github.com/microsoft/TypeScript/blob/d6c4afddb2c55f4a9dea7b59293a99a8fdea1799/tsc/internal/execute/incremental/referencemap.go) y [watcher race/invalidation tests](https://github.com/microsoft/TypeScript/blob/d6c4afddb2c55f4a9dea7b59293a99a8fdea1799/tsc/internal/execute/tsctests/watcher_race_test.go).
- Testing Library: <https://testing-library.com/docs/>.
- Playwright: <https://playwright.dev/docs/actionability>, <https://playwright.dev/docs/best-practices>.

### Seguridad web de plataforma

- CSP Level 3 (Working Draft al 2026-08-21): <https://www.w3.org/TR/CSP3/>.
- Trusted Types (Working Draft al 2026-08-21): <https://www.w3.org/TR/trusted-types/>.
- Fetch/CORS: <https://fetch.spec.whatwg.org/#http-cors-protocol>.
- Permissions Policy: <https://www.w3.org/TR/permissions-policy-1/>.

## Regla final

```text
Make the user's task and system state visible.
Preserve platform semantics before adding abstraction.
Measure real outcomes, accessibility and tails.
Turn every material failure into a safer contract.
```
