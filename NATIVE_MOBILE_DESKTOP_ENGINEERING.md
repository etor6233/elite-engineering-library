# Ingeniería de aplicaciones nativas móviles y de escritorio — manual operativo

> **Estado:** núcleo general v1 auditado el 2026-08-21; fuentes, integridad documental y cruces materiales cerrados. Capstone offline-first y contrato ejecutable disponibles en `ENGINEERING_EXECUTION_PLAYBOOK.md`; packages/dispositivos reales requieren evidencia propia.
> **Autoridad:** Android Developers; Apple Developer; Microsoft Learn/Windows App SDK; freedesktop/Flatpak; Flutter; React Native; Electron y Tauri para fronteras multiplataforma.
> **Propósito:** construir clientes instalables resilientes, seguros, accesibles y medibles que respeten lifecycle, recursos, permisos, distribución y convenciones de cada plataforma.

## Cómo usar este manual

```text
product journeys + datos + devices + plataformas
→ lifecycle/state/offline contract
→ arquitectura UI/domain/data/platform
→ permisos + seguridad + accesibilidad
→ performance/battery/memory/network budgets
→ tests en simulador y hardware real
→ build/sign/package/store/update
→ field telemetry + crash recovery + rollback
```

Una aplicación nativa no es un frontend web dentro de otro formato. Vive bajo un sistema operativo que puede recrear UI, suspender o terminar procesos, negar recursos, cambiar conectividad, limitar background work y exigir identidad criptográfica para distribuir actualizaciones.

## Procedencia

- `[SPEC]`: contrato o documentación oficial de plataforma.
- `[ACADEMIC]`: material académico identificado.
- `[CODE]`: implementación o sample oficial.
- `[PROD]`: práctica operativa atribuida.
- `[MEASURED]`: medición reproducible en dispositivo y escenario declarados.
- `[IMPL]`: diseño, código, prueba, presupuesto o runbook derivado.
- `[OPEN]`: decisión aún dependiente del producto/plataforma.

## 1. Contrato del producto nativo

Antes de elegir framework:

```text
usuarios, tareas críticas y riesgo:
plataformas/OS mínimos/form factors:
dispositivos baseline y premium:
inputs: touch/mouse/keyboard/stylus/gamepad/voz:
sensores/hardware/archivos/IPC requeridos:
online, degradado y offline:
estado efímero, durable y remoto:
background/push/deep links:
privacidad, permisos y datos sensibles:
accesibilidad/localización:
SLO de startup, frames, interacción y sync:
budgets de memoria, batería, red, disco y tamaño:
distribución, firma, canales y updates:
crash/telemetry/support/rollback:
```

El dispositivo lento, la red mala, el storage casi lleno y el proceso recién recreado son condiciones normales, no casos excepcionales.

## 2. Decidir native, shared core o multiplataforma

### 2.1 Criterios

| Opción | Elegir cuando | Costo principal |
|---|---|---|
| UI nativa por plataforma | máxima integración, hardware, accesibilidad, latencia o convenciones específicas | equipos/código duplicado |
| core compartido + UI nativa | dominio/protocolo complejo compartible y UX específica | FFI, toolchains y debugging cruzado |
| Flutter | producto visual consistente y amplia reutilización justifican engine/render propio | integración, binary/runtime y comportamiento de plataforma |
| React Native | equipo React y reuse de lógica aceleran entrega sin renunciar a módulos nativos | doble runtime, threading/native boundaries y ecosystem |
| Electron | desktop con fuerte base web y costo de Chromium aceptable | memoria, tamaño, updates y gran superficie privilegiada |
| Tauri | desktop web UI con core Rust y WebView del sistema | variación del WebView/OS, IPC y capability design |

### 2.2 Matriz de decisión obligatoria

Puntuar con evidencia:

- journeys críticos y fidelity nativa;
- hardware/APIs no cubiertas;
- latency/frame/startup/memory/battery;
- accessibility y platform conventions;
- offline/background;
- seguridad del boundary UI ↔ privilegiado;
- habilidades y capacidad de dos equipos;
- dependencia de plugins y velocidad de upgrades;
- test matrix y release cadence;
- vida esperada y costo de migración.

No decidir por demo, popularidad o porcentaje teórico de líneas compartidas. Medir una vertical slice en los dispositivos más débiles soportados.

### 2.3 Compartir semántica, no accidentes

Compartir preferentemente:

- modelos de dominio y reglas puras;
- validación, codecs y protocolos;
- sync/conflict policy;
- clientes API generados;
- casos y fixtures de aceptación.

Mantener específicos cuando corresponda:

- navegación, ventanas y back behavior;
- permisos, background y notifications;
- hardware, media y filesystem;
- accessibility semantics e inputs;
- UI que debe sentirse propia del sistema.

## 3. Lifecycle como máquina de estados

### 3.1 Regla central

Nunca asumir:

- que el proceso seguirá vivo;
- que una pantalla conserva memoria;
- que background equivale a ejecución libre;
- que shutdown/crash entregará un callback final;
- que una activación comienza desde la pantalla principal;
- que existe una sola ventana o escena.

Modelar explícitamente:

```text
not_running
→ launching/activating
→ foreground_active
↔ foreground_inactive
↔ background_limited
→ suspended/terminated
```

Los estados y transiciones reales varían por plataforma. El diseño portable se basa en obligaciones: restaurar, persistir, cancelar, liberar y reanudar de forma idempotente.

### 3.2 Startup y activación

Una app puede activarse por:

- icono;
- URL/deep link/universal or app link;
- archivo/documento;
- push/notification action;
- share extension;
- protocol handler;
- background event;
- restore/restart/update.

Parsear la activación como input tipado, autenticarla, deduplicarla y resolver navegación sólo después de que las dependencias mínimas estén listas. No ejecutar comandos privilegiados directamente desde una URL.

### 3.3 Persistencia por intención

| Estado | Ejemplo | Destino |
|---|---|---|
| frame/transitorio | hover, animación | memoria UI |
| UI restorable | scroll, tab, draft no crítico | saved state/local |
| dominio durable | orden, documento, configuración | base/archivo transaccional |
| secreto | token, key | secure storage/keychain |
| regenerable | thumbnail/cache | cache evictable |
| remoto | truth del servidor | local mirror + sync metadata |

Persistir al confirmar acciones importantes, no esperar una señal de terminación.

### 3.4 Recursos por lifecycle

Al perder visibilidad/foco:

- pausar animación/render/media no necesaria;
- reducir timers/polling;
- cancelar trabajo ligado a la UI;
- guardar datos críticos;
- liberar cámara, micrófono, GPU y handles compartidos;
- ocultar datos sensibles antes del snapshot del app switcher;
- mantener sólo trabajo permitido y explicable al usuario.

## 4. Arquitectura de referencia

```text
platform shell / lifecycle / activation
              ↓ events
UI view ← immutable UI state ← state holder
              ↑ intents
        use cases/domain
              ↓
repositories = policy + single source of truth
       ↓ local                 ↓ remote/device
database/files/cache       API/sensors/platform adapters
```

### 4.1 Reglas

- separar presentación, dominio, datos y adapters de plataforma;
- UI declarativa como función de state cuando el framework lo permita;
- flujo unidireccional: state baja, events suben;
- un único dueño mutable por dato;
- exponer state inmutable y operaciones explícitas;
- aislar SDK/OS detrás de adapters estrechos;
- inyectar reloj, red, filesystem, scheduler y identifiers;
- dependencias apuntan hacia contratos estables;
- no añadir domain layer ceremonial si no reduce complejidad o reuse.

**[SPEC]** Android recomienda separation of concerns, UI/data layers, una domain layer opcional, single source of truth y unidirectional data flow. Flutter documenta el mismo principio general de UI/data con responsabilidades y límites definidos.

### 4.2 State holder

Debe:

- combinar flows/sources;
- convertir dominio a UI state;
- aceptar intents tipados;
- representar idle/loading/content/empty/error/partial;
- ignorar resultados obsoletos;
- sobrevivir sólo el scope correcto;
- ser testeable sin renderer.

No debe poseer ventanas, contexts globales ni recursos más longevos que su scope.

### 4.3 Errores como estado

Distinguir:

- validación local;
- auth/permission;
- offline/timeout;
- conflict;
- server rejection;
- corruption/invariant;
- unsupported platform/version;
- cancellation.

Una cancelación esperada por navegación no es un error operativo. Una pantalla con cache stale puede ser usable mientras el refresh falla.

## 5. Concurrencia y aislamiento

### 5.1 Main/UI thread

El hilo principal procesa input, layout, commits y a menudo render coordination. No bloquearlo con:

- red o disk I/O;
- parsing/serialization grande;
- crypto/compression;
- locks contended;
- inference pesada;
- consultas sin límites;
- callbacks de FFI impredecibles.

Mover trabajo no basta: limitar concurrencia, cancelar lo obsoleto y volver a UI con state coherente.

### 5.2 Structured concurrency

Preferir tareas ligadas al scope que las necesita:

```text
screen scope closes → cancel requests/streams
new query supersedes old → cancel or discard old result
app backgrounds → cancel UI work; persist allowed durable work
logout → revoke tasks, secrets, caches and subscriptions
```

Definir ownership, cancellation, timeout, retry, priority, isolation y delivery thread. No crear fire-and-forget salvo que se transfiera a un scheduler durable.

### 5.3 Prioridades e inversión

- prioridad de tarea debe reflejar valor al usuario;
- no hacer sync masivo mientras una interacción espera el mismo recurso;
- evitar locks atravesando IPC/await;
- batch/coalesce de updates sin sumar latencia perceptible;
- aplicar backpressure a sensores, sockets y streams;
- medir queue delay además de execution time.

### 5.4 Swift, Kotlin y boundaries

**[SPEC]** Swift concurrency ofrece tasks, task groups y aislamiento; el modo estricto ayuda a detectar cruces inseguros. En Kotlin, coroutines sirven para trabajo asíncrono en proceso; trabajo que debe sobrevivir salida/reinicio requiere una facility durable como WorkManager en Android.

Para FFI/native modules aplicar `TOOLCHAINS_BUILDS_PACKAGING_FFI.md`: ownership, allocator, lifetime, errors, callback thread y cancellation deben formar un contrato verificable.

## 6. Datos locales y offline-first

### 6.1 Local source of truth

Para journeys offline:

```text
network → validate → transaction local → observable state → UI
UI intent → optimistic/local transaction → outbox → network → reconcile
```

La UI lee del store local canónico; el sync actualiza el store. Evitar que pantalla y red compitan como dos truths.

**[SPEC]** La guía Android offline-first exige al menos lectura funcional sin red, propone una fuente local canónica y trata writes, queues y conflictos como decisiones explícitas.

### 6.2 Sync contract

Definir por entidad/operación:

- stable local ID y server ID;
- version/vector/timestamp semantics;
- idempotency key;
- dirty/tombstone/pending state;
- upload/download ordering;
- conflict detector y resolver;
- retry/backoff/jitter;
- auth expiry;
- partial progress/checkpoint;
- retention/compaction;
- user-visible sync state.

### 6.3 Estrategias de escritura

| Estrategia | Uso | Riesgo |
|---|---|---|
| online-only | irreversible o requiere validación autoritativa inmediata | no funciona offline |
| queued | telemetry/logs y operaciones tolerantes | atraso/duplicado |
| lazy/optimistic | creación/edición crítica para UX | conflict/reconciliation |

No mostrar “guardado” si sólo está en memoria. Distinguir guardado local, pendiente de sync y confirmado remoto.

### 6.4 Conflictos

Elegir por semántica:

- server/client wins sólo si pérdida aceptable;
- last-write-wins requiere reloj/meaning adecuados;
- field merge si campos son independientes;
- append/CRDT si dominio lo soporta;
- intervención humana para decisiones irreversibles.

Registrar base/local/remote, resolver determinísticamente y probar reintentos. Un tombstone evita que una descarga resucite una eliminación.

### 6.5 Migrations y corrupción

- schema migrations forwards probadas sobre snapshots reales;
- downgrade policy explícita;
- backup/restore/import versionados;
- atomic replace y fsync donde durability importe;
- recovery ante disco lleno, kill y archivo corrupto;
- separar cache borrable de datos del usuario;
- encryption keys fuera del archivo cifrado.

## 7. Red móvil y resiliencia

Diseñar para latency alta, roaming, captive portal, cambio Wi-Fi/celular, pérdida parcial y reorder.

### 7.1 Cliente

- timeouts por fase y total;
- cancellation por lifecycle;
- retries sólo idempotentes o con idempotency key;
- exponential backoff con jitter y retry budget;
- paginación/streaming con límites;
- cache validators/versioning;
- request coalescing;
- bounded concurrency;
- TLS y hostname validation estándar;
- error envelope versionado;
- observabilidad sin PII/secrets.

### 7.2 Connectivity no es reachability

Que exista interfaz de red no prueba que el servicio responda. Usar señales de conectividad para optimizar, no para decidir definitivamente éxito; ejecutar la operación y manejar su resultado.

### 7.3 Tiempo real

Para WebSocket/stream:

- reconnect con jitter;
- session/sequence/resume token;
- gap detection y snapshot recovery;
- bounded buffer y drop/coalesce policy;
- foreground/background behavior;
- auth refresh;
- heartbeat basado en protocolo;
- no renderizar cada evento si excede frame budget.

La ingestión de order books de alta frecuencia pertenece al backend/systems; el cliente representa snapshots/deltas ya normalizados y limita visualización. Ver `MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md` y `NETWORKING_DISTRIBUTED_STREAMING.md`.

## 8. Background work, push y jobs

### 8.1 Clasificar primero

```text
¿debe sobrevivir a la pantalla? → app/process scope
¿debe sobrevivir al proceso/reboot? → durable scheduler
¿es visible y continuo para el usuario? → foreground/user-visible facility
¿es exacto por semántica? → API específica de alarm/event
¿puede esperar condiciones óptimas? → constrained deferrable job
```

### 8.2 Obligaciones

- input serializable y pequeño;
- idempotencia/dedupe;
- checkpoint/restart;
- deadlines y expiration handler;
- constraints de red/batería/carga;
- bounded retries;
- progress/cancel visible cuando corresponda;
- no depender de process globals;
- preservar privacidad en notifications.

**[SPEC]** Android recomienda WorkManager para trabajo persistente y permite constraints, chaining, retries y persistencia tras reboot; no es sustituto universal de trabajo inmediato. Apple restringe background y pide hacer lo mínimo posible, usando capabilities/facilities específicas.

### 8.3 Push

Push es señal, no base de datos ni garantía de entrega:

- payload mínimo y no sensible;
- fetch/reconcile desde truth autoritativa;
- collapse/dedupe/version;
- expiration;
- user controls y quiet hours;
- token rotation/revocation;
- métricas sent/delivered/opened sin confundirlas.

## 9. Deep links, intents, archivos e IPC

Todo input externo es no confiable.

- usar enlaces verificados cuando sea posible;
- canonicalizar URL/path;
- allowlist de scheme/host/action/type;
- validar tamaño, encoding y schema;
- autenticar sesión y autorizar acción;
- pedir confirmación contextual para irreversible;
- no interpolar en shell/SQL/path;
- limitar exported components/handlers;
- validar sender/origin en IPC;
- capability tokens y least privilege;
- prevenir replay con nonce/idempotency cuando aplique.

Probar cold launch, app activa, app autenticándose, recurso inexistente, link repetido y versión vieja.

## 10. Permisos, sensores y hardware

### 10.1 Permiso como journey

1. ofrecer valor sin permiso si es posible;
2. explicar necesidad en contexto;
3. pedir el mínimo alcance y duración;
4. manejar deny, restricted y revoked;
5. no pedir de nuevo en loop;
6. dirigir a settings sólo con explicación;
7. apagar recurso al dejar de necesitarlo;
8. auditar SDKs transitivos.

**[SPEC]** Android recomienda minimizar permisos, auditar accesos y usar storage/identifiers de menor alcance. Apple y Flatpak modelan acceso a recursos mediante entitlements/sandbox/portals.

### 10.2 Resource broker

Centralizar cámara, micrófono, ubicación, Bluetooth, USB, clipboard y notifications:

- ownership exclusivo/compartido;
- lifecycle y interruption;
- permission state;
- timeout/cancellation;
- mock/fake interface;
- audit trail sin capturar contenido sensible;
- fallback ante hardware ausente.

### 10.3 Frecuencia y privacidad

Elegir menor precisión/frecuencia suficiente. Sensor más preciso aumenta batería y sensibilidad de datos. No recolectar “por si acaso”; definir retention, propósito y eliminación.

## 11. Seguridad y privacidad local

### 11.1 Modelo de amenazas del cliente

Actores/superficies:

- dispositivo perdido o comprometido;
- app maliciosa vecina;
- deep link/file/document hostil;
- red hostil/MITM;
- dependencia/plugin comprometido;
- WebView/render content;
- IPC privilegiado;
- logs/backups/screenshots/clipboard;
- update/installer adulterado;
- reverse engineering/tampering.

El cliente no puede guardar un secreto compartido universal contra su propio usuario. Las decisiones autoritativas permanecen en servidor.

### 11.2 Storage

- tokens/keys pequeños en Keychain/Keystore/credential facility;
- datos privados en app-private storage;
- excluir caches/secrets de backups según política;
- no guardar secrets en preferences, source, logs o analytics;
- bloquear UI/snapshot sensible;
- borrar session data al logout según contrato;
- rotar/revocar server-side;
- usar crypto de plataforma, nunca diseños propios.

**[SPEC]** Apple Keychain permite condicionar acceso al estado del dispositivo. Android proporciona sandbox y storage privado. Esto no elimina riesgo en un dispositivo comprometido.

### 11.3 WebViews y desktop web shells

Separar contenido no confiable del proceso privilegiado:

- no habilitar Node/system APIs a remote content;
- context isolation/sandbox;
- CSP restrictiva;
- navigation/window allowlist;
- IPC schemas y sender/origin validation;
- capabilities por ventana/feature;
- no abrir URLs/shell con input sin validar;
- actualizar framework/WebView strategy;
- minimizar plugins y permisos.

**[SPEC]** Electron advierte que combina Chromium, Node, dependencies y código con privilegios mayores que un browser; su checklist exige aislamiento, sandbox e IPC validado. Tauri separa frontend WebView y core Rust mediante capabilities/IPC; una capability amplia fusiona fronteras de confianza.

### 11.4 Supply chain y firma

Aplicar `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md` y `TOOLCHAINS_BUILDS_PACKAGING_FFI.md`:

- dependencies locked/scanned;
- build reproducible donde sea viable;
- identidad de publisher protegida;
- CI signing con acceso mínimo;
- artifacts/SBOM/provenance inmutables;
- installer, binaries, libraries y updater firmados;
- rollback/compromise response.

## 12. Performance como experiencia y presupuesto

### 12.1 Journeys y métricas

| Dimensión | Medir |
|---|---|
| startup | cold/warm/hot, first frame y fully usable |
| input | event-to-response y main-thread stalls |
| render | frame time, missed frames, hitches/jank |
| memory | steady/peak, leaks, churn, pressure kills |
| CPU/GPU | utilization, work/frame y contention |
| energy/thermal | wakeups, radios, sensors, throttling |
| network | bytes, requests, reconnects, background use |
| storage | size, writes, I/O latency, disk-full behavior |
| package | download/install/update size y time |

Reportar distribución por device/OS/app version/journey; no sólo promedio.

### 12.2 Startup

- dibujar shell útil rápido;
- diferir inicialización no crítica;
- lazy-load features y SDKs;
- evitar I/O/network síncronos;
- medir cada activation path;
- no hacer trabajo dos veces entre splash y primera pantalla;
- cargar data local antes de refresh remoto;
- preservar símbolos/traces para atribución.

**[SPEC]** Android distingue time to initial display y time to full display; su guía observada recomienda medir startup, jank y field/lab data. Tratar sus números orientativos por versión como baseline inicial, no como SLO universal.

### 12.3 Frames y respuesta

El presupuesto por frame depende del refresh rate. El trabajo continuo debe terminar dentro del intervalo con margen para OS/compositor.

- perfilar release build en hardware real;
- virtualizar listas;
- reducir invalidation/recomposition;
- precalcular/batch fuera de main;
- no decodificar imágenes grandes al tamaño original;
- evitar layout passes y overdraw innecesarios;
- rate-limit telemetry y state updates;
- separar UI-thread time de render/GPU time.

**[SPEC]** Apple señala que 100 ms puede volver perceptible una demora discreta y que trabajo durante interacción continua requiere presupuestos mucho menores; medir en hardware soportado, incluidos dispositivos antiguos.

### 12.4 Memoria

- heap snapshots y retain graphs;
- ownership de images/buffers/caches;
- bounded caches con pressure callbacks;
- streaming/chunking de archivos;
- cancelar subscriptions/tasks al cerrar scope;
- evitar copias en FFI/serialization cuando estén medidas, sin romper seguridad;
- probar repeated navigation y long session;
- manejar warning/trim/low-memory.

### 12.5 Batería y térmica

Costo no es sólo CPU: radios, GPS, camera, display, wakeups y GPU dominan según caso.

- coalescer red/trabajo;
- constraints charging/unmetered;
- bajar sensor accuracy/frequency;
- pausar invisible;
- evitar polling si hay event/push;
- reducir animations bajo preferencias/situación;
- medir sesiones largas y thermal throttling.

### 12.6 Método científico

```text
field symptom → reproducible journey → trace/profile
→ bottleneck hypothesis → one controlled change
→ before/after distribution → regression gate
```

Apple documenta este ciclo y recomienda combinar datos de campo con Xcode/Instruments y performance tests. No optimizar por intuición ni microbenchmark aislado.

## 13. Android

### 13.1 Arquitectura y lifecycle

- components/activities son entrypoints efímeros, no data stores;
- single-activity/Compose es una opción moderna, no excusa para acoplar todo;
- preservar state ante configuration change y process death;
- repositories coordinan local/network;
- ViewModel/state holder tiene scope explícito;
- navigation/deep links se prueban desde cold start;
- layouts adaptan window sizes, foldables y rotation.

### 13.2 Background

- coroutine/thread para trabajo en proceso cancelable;
- WorkManager para trabajo durable diferible/reintentable;
- foreground service sólo para caso permitido y user-visible;
- APIs específicas vencen a un servicio genérico;
- constraints y exponential backoff protegen batería/red;
- workers idempotentes y con checkpoint.

### 13.3 Calidad

Medir con herramientas oficiales según pregunta:

- Macrobenchmark para journeys/startup/scroll;
- Microbenchmark para código aislado;
- Perfetto/system trace para scheduling;
- profiler de memoria/CPU/network;
- Android vitals/Play data para campo;
- tests en API levels, form factors y vendors relevantes.

### 13.4 Seguridad y release

- sandbox y app-private storage;
- permisos mínimos y runtime denial;
- exported components explícitos;
- Network Security Configuration/secure transport;
- Credential Manager/biometric sólo con threat model;
- Play Integrity como señal de riesgo, no autorización única;
- app signing, bundles y staged rollout;
- data safety/privacy declarations alineadas con código y SDKs.

### 13.5 Código AndroidX/AOSP — proceso efímero, intención durable

**[PROD]** AndroidX fijado en `3c4787e430e6df3f2d0699486de6fa8702f53f30` y AOSP `frameworks/base` en `1cdfff555f4a21f71ccc978290e2e212e2f8b168` (Apache-2.0). Son snapshots de implementación pública; cada app todavía debe fijar API level, versiones Jetpack y comportamiento de vendors/dispositivos reales.

`WorkerWrapper` comprueba elegibilidad persistida, usa ID generacional, cambia `ENQUEUED → RUNNING → SUCCEEDED/FAILED` o vuelve a `ENQUEUED` dentro de transacciones y sólo desbloquea dependencias cuando todos sus prerequisitos terminaron. `ForceStopRunnable` reconcilia jobs del sistema con la base local y reencola trabajo que quedó `RUNNING` tras crash/force-stop. La consecuencia no es que WorkManager vuelva exactly-once a una API remota: el worker puede repetir su cuerpo y necesita idempotency key/outbox/checkpoint del dominio.

`SavedStateHandle` requiere registro temprano y conserva estado pequeño serializable por key. Sirve para recreación de UI, no para datasets, secretos, jobs ni verdad de negocio. `ProcessLifecycleOwner` agrega counters de activities y demora pause/stop para absorber configuration changes; no entrega precisión milisegundo ni `ON_DESTROY` del proceso. Persistir intención antes de depender de un callback de background.

AOSP Binder mantiene identidad del caller por transacción/thread, permite limpiar/restaurar esa identidad en un bloque acotado y advierte sobre llamadas bloqueantes a código remoto. `oneway` cambia respuesta/backpressure/identidad disponible; no vuelve segura una cola ilimitada. Gate IPC: autorización con UID/package/user correctos antes de `clearCallingIdentity`, restore en `finally`, payload acotado/versionado, timeout/death recipient, no bloqueo de main/Binder pool, callbacks duplicados/tardíos y proceso remoto muerto durante cada paso.

Gate Android de recuperación:

1. matar el proceso antes/después de persistir intención, iniciar side effect, actualizar checkpoint y publicar UI;
2. duplicar scheduler/job/callback y actualizar/cancelar una generación mientras corre la anterior;
3. reboot, force-stop, clock change, constraint flapping, storage full/locked/corrupt y OS/vendor que pierde un job;
4. reconstruir desde DB/server, no desde singleton/ViewModel; comprobar que estado salvado restaura sólo navegación/input pequeño;
5. verificar resultado externo una vez por clave o una duplicación explícitamente tolerada, además de batería, retries, oldest work y causa de stop.

## 14. Plataformas Apple

### 14.1 Arquitectura y escenas

- SwiftUI `App`/scene o UIKit lifecycle según producto;
- cada escena puede tener state/lifecycle independiente;
- modelos observables no deben ocultar side effects;
- mantener dominio/storage fuera de views;
- soportar memory pressure, protected-data state y restoration;
- quietar trabajo y proteger snapshot al background;
- no asumir background execution ilimitada.

### 14.2 Concurrencia y performance

- aislamiento de UI/main actor explícito;
- structured tasks y cancellation;
- evitar priority inversion y synchronous wait;
- Instruments para Time Profiler, allocations/leaks, energy, I/O y network;
- MetricKit/Xcode Organizer para campo;
- XCTest/Swift Testing para logic, integration, UI y performance;
- dispositivo real y oldest supported hardware.

### 14.3 Seguridad y distribución

- Keychain con accessibility class correcta;
- App Transport Security;
- sandbox/entitlements mínimos;
- code signing de todos los nested executables;
- Hardened Runtime en macOS donde corresponda;
- App Store/TestFlight o Developer ID;
- notarization para distribución macOS externa aplicable;
- probar archive/export final, no sólo debug build.

**[SPEC]** Apple distingue notarization de App Review: el servicio inspecciona software firmado y emite un ticket verificable por Gatekeeper. Automatizar upload/log/stapling y probar el deliverable final.

## 15. Windows desktop

### 15.1 Modelo

- Windows App SDK/WinUI para nuevas apps cuando encaje; también puede integrarse con WPF, WinForms o Win32;
- separar windows/views de dominio y services;
- manejar launch, file/protocol activation, restart y multi-instance;
- guardar state periódicamente: un crash o cierre externo no promete cleanup;
- reducir GPU/timers/polling al minimizar/perder foco;
- probar DPI, scaling, multi-monitor, themes y múltiples inputs.

**[SPEC]** Windows App SDK desktop no usa la suspensión automática UWP tradicional, pero sigue bajo políticas de energía/recursos y cierres inesperados; no asumir CPU plena en background.

### 15.2 Privilegios e IPC

- ejecutar sin admin;
- aislar feature elevada en proceso mínimo;
- validar named pipes/COM/protocol/file inputs;
- ACLs y identity correctas;
- memory-safe language para parsing no confiable cuando sea viable;
- firmar EXE, DLL, installer y updater;
- crash dumps/symbols protegidos y correlacionables.

### 15.3 Packaging

MSIX aporta identity, install/uninstall y update confiables cuando sus constraints sirven. Si se distribuye unpackaged, declarar bootstrap/runtime/update/uninstall explícitamente. Probar clean install, upgrade, repair, per-user/per-machine, roaming profile y enterprise restrictions.

## 16. Linux desktop

### 16.1 Integración

- XDG paths para config/data/cache/state;
- desktop entry, icons, MIME y URL handlers;
- Wayland-first y X11 sólo según support matrix;
- portals para file chooser, URI, notifications, screenshots y recursos;
- D-Bus interfaces mínimas y versionadas;
- themes/fonts/scaling/input methods/accessibility;
- packaging por distribución definido, no “funciona en mi distro”.

### 16.2 Flatpak

**[SPEC]** Flatpak ejecuta apps en sandbox con acceso host mínimo por defecto; portals conceden recursos seleccionados sin permisos estáticos amplios.

- preferir portal a `--filesystem=host/home`;
- filesystem read-only y XDG scopes cuando alcance;
- evitar session/system bus completos;
- limitar talk/own names;
- testear permisos con la app realmente empaquetada;
- auditar D-Bus/portal access;
- declarar runtime/extensions y update channel.

AppImage, deb/rpm, Snap u otro formato cambian update, sandbox, dependencies y support. Seleccionar y probar por canal; no asumir equivalencia.

## 17. Flutter y React Native

### 17.1 Flutter

**[SPEC]** Flutter dibuja/composita su UI y usa un platform embedder para lifecycle, input, ventanas, threading y platform messages.

- separar UI/data y view models;
- no ocultar plataforma bajo abstractions falsas;
- medir raster/UI threads, shaders y startup;
- probar plugins en cada OS/arch/version;
- tipar method channels y errores;
- minimizar cruces/chatty IPC;
- validar binary size y engine/runtime updates.

### 17.2 React Native

- compartir lógica/UI sólo donde semántica coincida;
- archivos/branches específicos para divergencia real;
- medir JS y UI/native frame rates;
- evitar trabajo JS largo que demore interaction;
- validar Turbo Modules/Fabric/codegen compatibility;
- módulos nativos con ownership/thread contract;
- probar release build sin dev tooling;
- upgrade framework/plugins como proyecto planificado.

**[SPEC]** React Native documenta código platform-specific y native modules/components; su arquitectura moderna reduce el bridge serializado, pero no elimina boundaries, threads ni necesidad de perfilado.

### 17.3 Gate multiplataforma

- [ ] vertical slice usa hardware/API crítica;
- [ ] behavior y accessibility por plataforma;
- [ ] startup/frame/memory/battery medidos;
- [ ] offline/background/deep-link probados;
- [ ] native boundary tipado y testeado;
- [ ] plugin maintenance/security evaluados;
- [ ] store/package/signing automatizados;
- [ ] escape hatch nativo y costo conocidos.

## 18. Electron y Tauri

### 18.1 Proceso y privilegios

```text
untrusted/render/UI process
       ↓ validated narrow IPC
privileged main/core process
       ↓ least privilege adapter
filesystem/network/OS
```

Renderer comprometido no debe poder ejecutar capacidades generales. APIs privilegiadas se exponen por operaciones de dominio, no por `exec(path, args)` genérico.

### 18.2 Electron

- renderer sandbox/context isolation;
- Node integration deshabilitada para contenido remoto;
- CSP y navigation/window controls;
- preload API mínima e inmutable;
- validar sender y schema en cada IPC;
- utility processes para aislamiento/carga;
- actualización frecuente por Chromium/Node/Electron;
- medir memory/startup/installer/update delta.

### 18.3 Tauri

- capabilities por window/webview y feature;
- commands estrechos con input validation;
- remote URLs sin core access salvo necesidad demostrada;
- CSP/isolation pattern según threat model;
- auditar plugins Rust/frontend;
- probar WebView real de cada OS;
- updater/signing y migrations integrados.

No declarar automáticamente que Tauri es “más seguro” o Electron “más lento”: comparar architecture, privileges, bundle y mediciones del producto concreto.

## 19. UI adaptativa, accesibilidad e inputs

El diseño visual detallado pertenece a `FRONTEND_PRODUCT_ENGINEERING_UX.md`; aquí mandan integración y garantías del OS.

### 19.1 Adaptación

- window size, no device label;
- rotation/fold/resizing/multi-window;
- safe areas/system bars/notches;
- font scaling/dynamic type;
- localization, RTL y texto expansivo;
- high contrast/dark mode/reduced motion;
- mouse/keyboard/touch/stylus/gamepad;
- focus, shortcuts, context menus y drag/drop.

### 19.2 Accessibility contract

- semantic role/name/value/state/action;
- logical traversal/focus order;
- keyboard-only completion;
- screen reader announcements no redundantes;
- touch target y contrast;
- zoom/text scaling sin pérdida;
- captions/transcripts/haptics alternatives;
- errores asociados al control y recuperables;
- automated scan + manual assistive-tech test.

Native look sin semantic accessibility no es nativo de calidad.

## 20. Testing

### 20.1 Pirámide y matrices

| Nivel | Verifica |
|---|---|
| pure/unit | dominio, reducers, validation, sync/conflicts |
| component/view | state→UI, semantics, layouts |
| integration | DB, filesystem, API, scheduler, native module |
| UI journey | task crítica con OS real/simulado |
| package/install | firma, permisos, activation, update |
| field/canary | devices/OS/conditions reales |

La pirámide no elimina UI tests: reduce su cantidad y reserva alta fidelidad para journeys críticos.

### 20.2 Fault matrix

Inyectar:

- process death en cada paso;
- rotate/resize/background/foreground;
- offline, slow, loss, server 5xx y auth expiry;
- duplicate/out-of-order response;
- disk full/corruption/migration interruption;
- permission deny/revoke/restricted;
- low memory/thermal/battery saver;
- clock/timezone/locale change;
- update con pending work y old schema;
- inaccessible hardware/sensor interruption.

### 20.3 Device matrix

Elegir por riesgo, no combinatoria ciega:

- oldest/newest supported OS;
- low/high memory y CPU/GPU;
- architectures;
- screen sizes/densities/form factors;
- vendors/drivers relevantes;
- clean/upgrade installs;
- locale/accessibility modes;
- real device para performance, sensors y power.

### 20.4 Determinismo de test

- fake clock/network/location/scheduler;
- seeded data;
- idling/synchronization por estado, no sleeps;
- reset explícito de permissions/storage/account;
- artifacts: logs, video, screenshot, trace, device metadata;
- flaky test como defecto con owner y deadline, no retry infinito.

## 21. Observabilidad y operación de campo

### 21.1 Eventos mínimos

```text
app/session install identity pseudonymous
app/build/OS/device class
journey + phase + duration
result/error taxonomy
network/cache/offline state
background job attempt/result
crash/hang/ANR/OOM
release channel/experiment
```

No enviar payloads, tokens, URLs sensibles, texto del usuario o exact location salvo contrato de privacidad explícito.

### 21.2 Métricas

- crash-free users/sessions;
- hangs/ANR;
- startup and interaction p50/p95/p99;
- jank/hitches;
- memory/termination;
- battery/network/storage;
- sync freshness/conflicts/backlog;
- install/update success;
- permission funnel;
- error recovery success.

Separar lab regression de field SLO. Correlacionar por build digest y symbol mapping.

### 21.3 Kill switch y compatibilidad

Backend/feature flags deben permitir:

- apagar feature defectuosa;
- mantener clientes viejos durante ventana declarada;
- degradar payload/capability;
- revocar app/token comprometido con cuidado;
- migrar por etapas;
- preservar offline pending operations.

No usar remote config para descargar código no revisado o evadir el canal de distribución.

## 22. Build, firma, publicación y updates

### 22.1 Pipeline

```text
pinned source/deps/SDK
→ build release matrix
→ unit/integration/UI/performance/security
→ package final
→ sign/notarize where applicable
→ verify/install/smoke on clean device
→ staged/canary publish
→ field gates
→ expand or halt/rollback
```

Promover el mismo artifact probado. Guardar symbols, mappings, manifests, privacy declarations, SBOM y provenance ligados al digest.

### 22.2 Signing

- claves fuera del repo y runner general;
- hardware/managed key storage cuando corresponda;
- roles y approval separados;
- rotation/revocation/recovery documentados;
- timestamps/notarization logs;
- nested components firmados;
- updater verifica firma antes de reemplazar;
- reproducir quién firmó qué digest y cuándo.

### 22.3 Rollout

- internal → beta/canary → small cohort → gradual → full;
- gates por crash/hang/startup/core journey;
- pause automático/humano;
- database/config/protocol compatibility;
- server feature flags independientes de store latency;
- rollback cuando canal lo permita o forward fix seguro;
- support notes y known issues.

Una mobile store release puede no poder revertirse instantáneamente; diseñar backend y migrations para convivencia y forward recovery.

### 22.4 Update correctness

Probar:

- N-1/N-2 → N;
- update interrumpido;
- insufficient disk;
- pending background jobs/outbox;
- session/key migration;
- deep link durante update;
- old client con new server;
- new client con staged old server;
- uninstall/reinstall/restore backup.

## 23. Anti-patrones

- guardar dominio sólo en Activity/View/window;
- asumir callback de termination;
- request de red sin cancellation;
- main-thread disk/JSON/crypto;
- polling permanente;
- background service para cualquier necesidad;
- permiso al launch sin contexto;
- token en preferences/log;
- WebView remoto con privilegios;
- IPC genérico sin allowlist/schema;
- “offline” que sólo cachea una pantalla;
- optimistic UI sin pending/error/reconcile;
- last-write-wins accidental;
- medir sólo emulator/debug/premium device;
- promedio sin percentiles/device class;
- test con sleeps/retries ilimitados;
- package/store artifact distinto del probado;
- migration no reversible ni forward-recoverable;
- cross-platform sin native escape hatch;
- duplicar toda UI por dogma cuando core compartido reduce riesgo.

## 24. Gate profesional

- [ ] Product/platform/device/support contract.
- [ ] Native/shared/cross-platform decision medida.
- [ ] Lifecycle/activation/process-death model.
- [ ] UI/domain/data/platform boundaries.
- [ ] State ownership/UDF/cancellation.
- [ ] Offline truth/sync/conflict/migrations.
- [ ] Background/push/deep links/IPC.
- [ ] Permissions/hardware/privacy/threat model.
- [ ] Startup/frame/memory/battery/network budgets.
- [ ] Accessibility/input/adaptive UI.
- [ ] Fault/device/package test matrices.
- [ ] Field telemetry con minimización de datos.
- [ ] Signing/store/update/rollback/compatibility.
- [ ] Runbook, symbols y support window.

## 25. Contrato para Codex

```text
Producto y journeys:
Plataformas/OS/devices:
Arquitectura/framework y razones:
Lifecycle/activation/background:
Datos/offline/sync:
Permisos/hardware/IPC:
Seguridad/privacidad:
Performance/accessibility budgets:
Distribución/update/support:

Entregar:
1) platform/product contract y ADR;
2) state machine de lifecycle y activaciones;
3) módulos UI/domain/data/platform;
4) ownership, concurrency y cancellation;
5) local truth, outbox, sync y conflict policy;
6) permission/threat/privacy design;
7) performance budgets y measurement plan;
8) unit/integration/UI/fault/device suites;
9) signed package, install/update/rollback pipeline;
10) field SLO, dashboards, alertas y runbook.

No asumir proceso, red, permiso ni hardware disponible.
No bloquear main/UI thread.
No guardar secrets en storage general ni logs.
No exponer IPC o WebView privilegiado genérico.
No declarar performance sin release build y dispositivo real.
No declarar offline sin process-death/conflict tests.
No publicar un artifact distinto del validado.
```

## 26. Suites de aceptación

| Suite | Inyección | Aceptación |
|---|---|---|
| `lifecycle_matrix` | rotate/resize/background/kill/relaunch/multi-window | state durable intacto, tareas canceladas/reanudadas correctamente |
| `activation_matrix` | icon/link/file/push cold y warm | routing autorizado, idempotente y recuperable |
| `offline_sync` | loss/reorder/duplicate/auth expiry/conflict | local usable, outbox bounded, reconcile determinístico |
| `main_thread` | critical journeys bajo carga | budgets y cero I/O/blocking prohibido según trace |
| `resource_pressure` | low memory/disk/battery/thermal | degrada sin corrupción ni pérdida silenciosa |
| `permission` | deny/revoke/restricted/partial | UX funcional/degradada y sin acceso indebido |
| `native_boundary` | malformed IPC/FFI/callback after close | rechazo seguro, sin leak/UAF/deadlock |
| `accessibility` | keyboard/screen reader/text scale/contrast | journeys críticos completables |
| `package_upgrade` | clean/N-2/interrupted/offline update | firma, migrations, pending work y recovery correctos |
| `field_canary` | cohort/device/OS reales | gates de crash/hang/startup/journey antes de expansión |

## 27. Capstones

### A. Mobile offline-first

Android+iOS o framework justificado: local truth, outbox, auth refresh, conflicts, push/deep links, background sync, accessibility, staged release y field telemetry. Kill en cada transición.

### B. Desktop privilegiado mínimo

Windows+macOS+Linux o support matrix justificada: filesystem por selector/portal, IPC capability-based, signed package/updater, sandbox, multi-window, crash recovery y package-real tests.

### C. Cliente de streaming observable

Snapshot+deltas con sequence/gap recovery, backpressure/coalescing, offline state, UI frame budget y reconnect. El cliente no reconstruye una verdad incompleta sin detectar gaps.

## 28. Inventario de autoridades

| Fuente | Cobertura |
|---|---|
| Android Developers | architecture/UDF/SSOT, lifecycle, offline-first, background, privacy/security, performance |
| Apple Developer | SwiftUI/UIKit lifecycle, concurrency, Keychain/sandbox, performance/testing, signing/notarization |
| Microsoft Learn | Windows App SDK, lifecycle, security/accessibility/performance, MSIX |
| Flatpak docs | sandbox, portals, permissions, packaging/debugging |
| Flutter docs | engine/embedder y arquitectura UI/data |
| React Native docs | new architecture, threading/rendering, platform code y native modules |
| Electron docs | multi-process, sandbox e hardening |
| Tauri docs | WebView/core trust boundary, IPC y capabilities |

Las APIs, políticas de stores, SDKs y versiones cambian. Pin en cada proyecto y verificar la documentación vigente antes de release.

## 29. Fuentes oficiales verificadas

- Android app architecture: <https://developer.android.com/topic/architecture>
- Android offline-first: <https://developer.android.com/topic/architecture/data-layer/offline-first>
- Android background tasks: <https://developer.android.com/develop/background-work/background-tasks>
- Android persistent work: <https://developer.android.com/develop/background-work/background-tasks/persistent>
- Android performance measurement: <https://developer.android.com/topic/performance/measuring-performance>
- Android security checklist: <https://developer.android.com/privacy-and-security/security-tips>
- Android privacy checklist: <https://developer.android.com/privacy-and-security/about>
- AndroidX fijado `3c4787e430e6df3f2d0699486de6fa8702f53f30`: [WorkerWrapper](https://github.com/androidx/androidx/blob/3c4787e430e6df3f2d0699486de6fa8702f53f30/work/work-runtime/src/main/java/androidx/work/impl/WorkerWrapper.kt), [force-stop reconciliation](https://github.com/androidx/androidx/blob/3c4787e430e6df3f2d0699486de6fa8702f53f30/work/work-runtime/src/main/java/androidx/work/impl/utils/ForceStopRunnable.java), [SavedStateHandle support](https://github.com/androidx/androidx/blob/3c4787e430e6df3f2d0699486de6fa8702f53f30/lifecycle/lifecycle-viewmodel-savedstate/src/commonMain/kotlin/androidx/lifecycle/SavedStateHandleSupport.kt) y [process lifecycle](https://github.com/androidx/androidx/blob/3c4787e430e6df3f2d0699486de6fa8702f53f30/lifecycle/lifecycle-process/src/main/java/androidx/lifecycle/ProcessLifecycleOwner.kt).
- AOSP `frameworks/base` fijado `1cdfff555f4a21f71ccc978290e2e212e2f8b168`: [Binder](https://github.com/aosp-mirror/platform_frameworks_base/blob/1cdfff555f4a21f71ccc978290e2e212e2f8b168/core/java/android/os/Binder.java), [Parcel](https://github.com/aosp-mirror/platform_frameworks_base/blob/1cdfff555f4a21f71ccc978290e2e212e2f8b168/core/java/android/os/Parcel.java) y [ActivityThread](https://github.com/aosp-mirror/platform_frameworks_base/blob/1cdfff555f4a21f71ccc978290e2e212e2f8b168/core/java/android/app/ActivityThread.java).
- Apple SwiftUI apps: <https://developer.apple.com/documentation/technologyoverviews/swiftui>
- Apple UIKit lifecycle: <https://developer.apple.com/documentation/uikit/managing-your-app-s-life-cycle>
- Apple Swift concurrency: <https://developer.apple.com/documentation/swift/concurrency>
- Apple security: <https://developer.apple.com/documentation/security/>
- Apple app performance: <https://developer.apple.com/documentation/xcode/improving-your-app-s-performance>
- Apple testing: <https://developer.apple.com/documentation/xcode/testing>
- Apple notarization: <https://developer.apple.com/documentation/security/notarizing-macos-software-before-distribution>
- Apple Human Interface Guidelines: <https://developer.apple.com/design/human-interface-guidelines/>
- Windows App SDK: <https://learn.microsoft.com/en-us/windows/apps/windows-app-sdk/>
- Windows app lifecycle: <https://learn.microsoft.com/en-us/windows/apps/develop/launch/app-lifecycle>
- Windows best practices: <https://learn.microsoft.com/en-us/windows/apps/get-started/best-practices>
- MSIX: <https://learn.microsoft.com/en-us/windows/msix/>
- Flatpak concepts: <https://docs.flatpak.org/en/latest/basic-concepts.html>
- Flatpak sandbox permissions: <https://docs.flatpak.org/en/latest/sandbox-permissions.html>
- Flutter architecture: <https://docs.flutter.dev/resources/architectural-overview>
- Flutter app architecture: <https://docs.flutter.dev/app-architecture/guide>
- React Native architecture: <https://reactnative.dev/architecture/overview>
- React Native native platform: <https://reactnative.dev/docs/native-platform>
- React Native performance: <https://reactnative.dev/docs/performance>
- Electron process model: <https://www.electronjs.org/docs/latest/tutorial/process-model>
- Electron security: <https://www.electronjs.org/docs/latest/tutorial/security>
- Tauri core concepts: <https://v2.tauri.app/concept/>
- Tauri security/capabilities: <https://v2.tauri.app/security/> y <https://v2.tauri.app/reference/acl/capability/>

## 30. Cruces iniciales

| Manual | Frontera |
|---|---|
| `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md` | drivers, boundaries, state y evolución |
| `TOOLCHAINS_BUILDS_PACKAGING_FFI.md` | SDK/target, native modules, firma/package/release |
| `FRONTEND_PRODUCT_ENGINEERING_UX.md` | interaction, accessibility, state UI y design systems |
| `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md` | threads, memory, profiling y latency |
| `GPU_ACCELERATED_COMPUTING.md` | render/inferencia device, async completion, memory y thermal |
| `SOFTWARE_BACKEND_API_ENGINEERING.md` | auth/API/idempotency/compatibility/feature flags |
| `NETWORKING_DISTRIBUTED_STREAMING.md` | reconnect, sequence, streaming y backpressure |
| `DATABASE_STORAGE_INTERNALS.md` | local durability, migrations y sync metadata |
| `DATA_ENGINEERING_ANALYTICS.md` | telemetry, schemas, quality y privacy |
| `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md` | threat model, secrets, signing, incident y field SLO |
| `AI_ENGINEERING_MASTER_MAP.md` | on-device inference, model packaging y privacy |

Cruces materiales cerrados el 2026-08-21: arquitectura, toolchains, frontend, sistemas, GPU, backend, redes, storage, datos, seguridad e IA incorporan la obligación recíproca que cambia una decisión. Este manual gobierna lifecycle del OS, app state, offline/background, permisos, integración nativa y distribución instalada; los receptores conservan dominio, seguridad, física, datos y experiencia. No se añadieron cruces ceremoniales.

## 31. Puerta de cierre

- [x] Contrato, decisión native/shared/cross-platform.
- [x] Lifecycle/state/activation/concurrency.
- [x] Offline/sync/background/network.
- [x] Permisos/hardware/security/privacy.
- [x] Performance/accessibility/testing/observability.
- [x] Android/Apple/Windows/Linux.
- [x] Flutter/React Native/Electron/Tauri.
- [x] Signing/package/update/rollout.
- [x] Codex/suites/capstones/fuentes.
- [x] Auditoría Markdown/referencias: fences, encabezados, caracteres y referencias internas.
- [x] Cruces recíprocos materiales y arbitraje.
- [x] Contrato ejecutable y capstone native offline-first disponibles en `ENGINEERING_EXECUTION_PLAYBOOK.md`; la app y device evidence permanecen correctamente `planned`.

El conocimiento general v1 queda auditado. No declarar suites ejecutadas sin artifacts, dispositivos y resultados reales.
