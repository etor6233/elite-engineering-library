# V312 — recepción y decisión de devoluciones recuperables

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

Mantenimiento T2802/T2803/T2804 desde checkpoint50. Fecha2026-09-08.
No cambia reglas comerciales, dependencias ni migraciones. Perfil67/745.
Journey0.10.17, Portal0.14.10, Browser0.1.22; correcciones AUTHORED,
REBUILD_VERIFIED/CONDITIONED. Sin ZIP, promoción integral ni efectos externos.

## Fallos observados y corrección

FAIL488: recepción200/replay409 y decisión200/replay409 eran durables,
pero ambas consultas exactas devolvían404. El sender genérico podía afirmar
no aplicado después de un commit y recargar sin conservar la referencia.
FAIL489: mover la organización del stock sintético dejaba un grafo incoherente;
la lista exponía dos casos y los comandos todavía creaban recepción, decisión,
cuatro solicitudes y dos eventos. No incidente ni dato real de producción.

La consulta GET /v1/franchise/returns/result selecciona por tenant, organización
y autorización con handover:manage y no-store; no usa la lista como oráculo.
Una única sentencia reconstruye autorización, recepción, decisión y solicitudes.
Valida bindings de entrega/pedido/stock/cliente/excepción y recepción. Service
rechaza proyecciones inconsistentes. ReturnCases usa snapshot repeatable-read,
límite1..100 y consulta de efectos sólo para las decisiones seleccionadas.
ReceiveReturn y DecideReturn bloquean el mismo grafo dentro de su transacción;
conservan inmutabilidad, auditoría atómica y conflictos409 ante duplicados.

El portal conserva sólo acción, autorización, recibo y SHA256 scoped en sesión.
No guarda serie, notas ni tokens en esa referencia. Fence antes de digest/POST,
storage fallido implica cero POST y consulta exacta funciona aun sin lista.
Comprueba huella y evidencia devuelta antes de habilitar Continuar operando.
El botón actualiza el caso y habilita una nueva acción sin repetir la anterior.
Conserva actor/fecha y muestra solicitudes como solicitadas, nunca ejecutadas.
Ayuda return-operations-view/1.0.0, regiones accesibles y revisión explícita.
Se retira el sender genérico; permanecen botones de consulta read-only existentes.
Un reembolso crea cuatro solicitudes; un cambio tres, sin fiscal. Regla existente.

## Evidencia conectada

| Gate | Resultado |
|---|---|
| Web reconstruida |115 PASS/1 SKIP y build PASS |
| Go reconstruido |baseline go test ./..., vet y build PASS; skips sin DB no equivalen a integración |
| PostgreSQL |11 tests por3 repeticiones PASS, incluyendo recuperación y grafo cruzado |
| Concurrencia |8 recepciones→1 éxito/7 conflictos;8 decisiones→1 éxito/7 conflictos; actor ganador conservado |
| Atomicidad |outbox duplicado revierte recepción; en decisión revierte decisión y todas las solicitudes |
| Inmutabilidad |6 intentos update/delete bloqueados; actor/notas/evidencia conservados |
| Alcance |tenant/org/autorización ajenos bloqueados; grafo cruzado rechaza lista/escritura sin efectos; locks bloquean cambios de entrega, pedido y stock |
| Identidad |101 autorizaciones posteriores desplazan el caso fuera de lista100; consulta exacta desde pool nuevo recupera evidencia |
| Navegadores |92 fases/4 proyectos: Chromium desktop/mobile, Firefox y WebKit desktop PASS |
| Nueva fase por proyecto |3 recepciones/3 decisiones/11 solicitudes/6 eventos/6 actores; pérdidas, JSON inválido, carreras200/409, storage0POST, GET503 y contenido alterado bloqueado |
| Lista caída |recuperación del caso conservado sin tarjeta ni lista; Continuar operando no envía POST |
| Permisos |consulta ajena/lector bloqueada, actor forjado400, replay409; backend privado no se muestra |
| Agenda existente |4 proyectos PASS |
| Reconstrucción |10/10 archivos byte-idénticos,67 packs/745 archivos; normalización LF explícita conservada |
| Reinicio sintético |autorizaciones/recepciones/decisiones/solicitudes/eventos idénticos; no restore/PITR del target |

## Fuzzing finito y límites

GO_NATIVE_FUZZ_GATE0.1.0 materializado como kit de prueba separado (4 archivos),
runner verificado con1 positivo/2 negativos. Perfil habilitado:10s,2 workers,
FuzzReturnCaseResult,8 semillas del dominio (return/exchange y etapas reales,
incluido binding corrupto). Invariantes: lectura sin mint/mutación, fallo retorna
DTO vacío, éxito no reescribe evidencia y respeta scope/padres/owners/estado.
Baseline, semillas y fuzz PASS:489402 ejecuciones; no input fallido observado.
Receipt conserva perfil, toolchain y locks. No reemplaza SAST, DAST, carga,
seguridad ofensiva, autenticación live ni prueba exhaustiva de PostgreSQL.

FAIL490: retirar el sender retiró por error el tipo Command usado por agenda;
restaurado sólo el tipo, build rojo y verde conservados. La fecha inicialmente
mal rotulada se corrigió con before/after registrado. FAIL491: Python produjo
CRLF, el compositor LF; el gate hash detectó diferencia. Se conservaron bytes
previos, se normalizaron sólo terminadores y todos los gates finales corrieron
sobre rebuilt. FAIL492: wrapper de agenda emitió PASS ante cuatro SKIP por
flag ausente; log preservado, gate explícito repetido y cuatro receipts reales
exigidos antes del cierre. FAIL493: el generador contó81 PASS con subtests
en vez de33 tests raíz; se ancló el conteo sin alterar logs ni repetir runtime.
FAIL473 conserva recurrencias de diagnóstico/rutas truncadas.

FAIL457 se cierra para los comandos ambiguos inventariados tras V299–V312;
no equivale a cerrar todo T2804. Entregas prepared y autorizaciones son fixtures
sintéticos por owners, no regla comercial de liberación aprobada. Ronda D/corpus,
canal histórico, target/IdP/proveedores V295 y aceptación material siguen abiertos;
ARCA diferida. Faltan los cierres integrales T2801–T2810, no se infiere porcentaje.
No se prueba coordinación multi-tab, pérdida completa del storage, acciones
downstream ejecutadas ni aceptación humana de accesibilidad/usabilidad.

## Receipts reproducibles

Staging: `%LOCALAPPDATA%/Temp/elite-v312-5d957894149648ef938b06860d608ad3`.
Los logs, scripts, snapshots sintéticos y receipt están retenidos allí.

| Archivo | SHA256 |
|---|---|
| returns-red.log | `ab73b52153ff1d56782fced49be9982dae68ac827370dcea1f9b7c10d70db466` |
| returns-scope-red.log | `f57591f9ae2e6a74ecdfcda1e96020c622e779fadf18c453217a9b717b3a1488` |
| returns-r1.log | `74862a148624a5efe79ce560e05c58ab62ec94a0344dfe9700f407cab867b4ec` |
| identity-r1.log | `1b90e46a2a35b2874a0f22c39e59e1adb2c60da149ce8070e5d4216485eaf470` |
| postgres-rebuilt.log | `fe68fcf95fb545bcedd7f5b5f07fa28ffeb705cfdcf5723a33997bc948e359b5` |
| report-guard-red.log | `df79195c60224507c99ecf4a4501302e8ac0493c5451bdc412527ec0cfd3fa9b` |
| web-test-rebuilt.log | `6fd9b290358773d8de6daffb1043128432e04321057de56513cc2a9e9eee15b5` |
| web-build-r1.log | `c19756e30240208841b77e22a64975a9c125e9c28f3d4413e0959aff329b9ada` |
| web-build-rebuilt.log | `260cd078dcf00d2f17668968efce757cb932682f2fc5b132ce387ae25dee8b3e` |
| browser-final.log | `dd59cfa9cba4715562619f4f25820c282c8b2bead90eb165a2de2debe4c4189a` |
| agenda-final.log | `6314a598b559289c1d8ea964312bb8fc64a2c707c4b09145e45a4878bcb4217c` |
| agenda-skipped.log | `677f5156e2f349a5b22d64e6de5f99f352bed4671d2abb298a27da9b21095fc8` |
| fuzz-runner-verify.log | `cf4189cef7b85e35b36ecdde6c1bee8ed93ff78e1c5497afc88fbbdbb01acf95` |
| fuzz-profile.json | `41672e1a46f8cd26a7edb24766b48ebdb9d12ca900ffced8fd1225b8e93b0c74` |
| fuzz-receipt.json | `ff8a4588c26b7023bb7f3fa67093202d9f88a0c038facddb21c44185116ba118` |
| fuzz-rebuilt.log | `92dd9cf678d6175a584f6e6948bfcada85c697d1c1ab53952788882b592d340d` |
| go-vet-rebuilt.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| go-build-rebuilt.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| canonical.log | `0b25e9a5e67f8d5be6ac0fbe06a487e30e46045f3c0829fcd2def9f46c806a1e` |
| normalization.json | `0ee038f2dbe1593e260fe9de97be6468e4cbb22e094f0ddbf79a834bbffa60c2` |
| identity-rebuilt.log | `4235de86b9a2060d1b6b62e9dbe51d267dbffc06c7933d31c9a286e13ca8d48e` |
| returns-before-restart.json | `8665b1631a5fb9feccde9194d5d5fc3af1ac1ee62e0f32070b7defba42c1a441` |
| returns-after-restart.json | `8665b1631a5fb9feccde9194d5d5fc3af1ac1ee62e0f32070b7defba42c1a441` |
| restart.log | `dea29ee59dbb83ce32372a2a99c843cb676c241e965cb0b7dd35854a85da639d` |

| Archivo canónico reconstruido | SHA256 |
|---|---|
| internal/franchisejourney/service.go | `eca21da528bed93046aa36a435ae4d91157c4efe229cbc55729d608371c3216b` |
| internal/franchisejourney/service_test.go | `6fc88efaef5ae994bb1b786a94e8d86325cd348ff900602dc3e7d835eaee159d` |
| internal/platform/httpapi/franchisejourney.go | `2439a78ae334f0a62376963c8a627e52a4627c39e7926b3927fda361230964fe` |
| internal/platform/httpapi/franchisejourney_test.go | `c14cdd1cf871ca472f328b9fab2d3482888e9e563c232defac777dba9fbd90d9` |
| internal/platform/postgres/franchisejourney.go | `7fdee6385a6fd7ed9be4991251d4cd03f97af61ca0092b89706618be80794e8e` |
| internal/platform/postgres/franchisejourney_integration_test.go | `af93b84ec8f18f2b7afad402f399f23fde3865d8ab22d874f42a8da74767afe4` |
| src/components/franchise-command-panel.tsx | `5424f0ef08974699bc8b55cdeba2a3fde2fc0dceb44b43e432c3d5c3644eb9fb` |
| src/app/api/enterprise/franchise/commands/route.ts | `c48ac2f8c73b027b93e24e1e26ff0701309ba53993c8156a20d62519cebf9f56` |
| src/app/api/enterprise/franchise/commands/route.test.ts | `5e190d894b1e396b6d99e86e321b2216cc2cf64b24c21ea7b1f160b5fe0f086e` |
| microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs | `361705624a8eebffffaa4b3dc5ccd902bda950ad863111a0246d38c2e30d86e6` |

## Cierre general

VERIFY_LIBRARY_PASS con160 packs/1436 archivos/745 Markdown y51 perfiles;
composición67/745. Checkpoint51 validado antes del verificador. Assurance
TEST31/EVID22 enlazada: plan PASS, evidence integral BLOCKED por pendientes
previos preservados. No equivale al100% del roadmap. Log library-final.log
SHA256: 31ee5b76db5c426996e3f278c46eff05ad28c800489246794458f7480aac92a1.
