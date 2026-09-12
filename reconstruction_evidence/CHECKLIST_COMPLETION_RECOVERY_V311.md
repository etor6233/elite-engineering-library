# V311 — recuperación de checklist completada y entrega presentada

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

Mantenimiento T2802/T2803/T2804, FAIL457, desde checkpoint48 validado.
Fecha2026-09-08 UTC. Misma composición67/745; sin nuevas dependencias,
migraciones, reglas comerciales, proveedor ni efecto externo. No promoción/ZIP.

## Fallo y corrección

El rojo confirma POST200, replay409,2 respuestas y1 evento auditado durables,
pero GET404. El sender genérico podía afirmar no aplicado después del commit.
Journey0.10.16 conserva el contrato de completado: CAS prepared→presented,
respuestas y actor atómicos, binding y respuestas inmutables. Añade GET
/v1/franchise/handovers/{id}/checklist-result con handover:manage y no-store.
Una consulta SQL exacta valida tenant/org/handover/cliente/pedido/stock y la
checklist publicada; descarta respuestas de otro binding o actor. Devuelve
respuestas ordenadas, actor/fecha originales y estado/versión actuales sin
reabrir ni volver a presentar. Service rechaza proyecciones incoherentes.
No devuelve datos del cliente ni reemplaza los owners de aceptación/rechazo.

Portal0.14.9 conserva sólo IDs, versiones y SHA256 scoped; nunca respuestas,
series, evidencias ni tokens en esa referencia. Fence antes de digest/POST;
fallo de storage0POST, bloqueo tras incertidumbre, GET sin POST y comparación
de checklist/respuestas normalizadas. Consulta permite estado posterior
accepted/rejected y muestra actor real; no atribuye automáticamente el efecto
al último intento. Preparar otra presentación exige consulta positiva y
ausencia de POST/GET activo. Ayuda checklist-completion-view/1.0.0, región
accesible y mensajes que conservan la incertidumbre.
Browser0.1.21 conecta el nuevo recorrido a las21 fases previas. Código AUTHORED,
REBUILD_VERIFIED/CONDITIONED; pruebas locales no autorizan promoción integral.

## Verificación

| Gate | Resultado |
|---|---|
| HTTP rojo→verde | GET404 tras commit pasa a200 con respuestas y actor originales; replay conserva409 |
| Web |114 PASS/1 SKIP explícito, build PASS |
| Go recompuesto |test/vet/build ./... PASS; skips sin DB no contados como conexión |
| Dominio |14 proyecciones incoherentes rechazadas; presented/accepted/rejected válidos |
| PostgreSQL |9 tests,3 repeticiones;8 completados concurrentes→1 éxito/7 conflictos; respuestas/evento/actor únicos |
| Atomicidad |outbox duplicado revierte respuestas y estado/version a prepared/1; no resultado parcial |
| Inmutabilidad/scope |update/delete de respuestas y cambio de actor rechazados; tenant/org/id ajenos, prepared y stock cruzado bloqueados |
| Recuperación |aceptación y rechazo por owners reales conservan respuestas/actor/fecha; fresh pool recupera estado actual |
| Browser recompuesto |88 fases/4 proyectos: Chromium desktop/mobile, Firefox y WebKit desktop PASS |
| Nueva fase |3 entregas/6 respuestas/3 eventos/3 actores por proyecto; pérdida, JSON inválido, carrera200/409, storage0POST, recarga, GET503, contenido distinto bloqueado, actor forjado y permisos |
| Avance posterior |HTTP del cliente acepta1/rechaza1; UI recupera1 accepted/1 rejected/1 presented sin cambiar estado |
| Agenda existente |4 proyectos PASS |
| Reconstrucción |10/10 archivos byte-idénticos;67 packs/745 archivos |
| Reinicio |snapshot de entregas completadas/respuestas/eventos idéntico tras restart sintético; no restore/PITR del target |

FAIL485 documenta recuperación ausente. FAIL486 es recurrencia de FAIL438:
segunda unidad sintética con VIN/batería NULL viola UNIQUE NULLS NOT DISTINCT.
Se declaran identificadores sintéticos distintos; no se debilita DDL ni se
demuestra compatibilidad universal de unidades sin VIN/batería. FAIL487 conserva
ParserError de edición inline, sin ejecución/modificación; reemplazo literal
por apply_patch y gates finales sobre esa fuente. No se ocultan los rojos.

FAIL457 permanece OPEN para recepción y decisión de retornos. La creación de
entregas prepared en pruebas es un fixture explícito, no una regla de liberación
comercial. Esa decisión, ronda D/canal-formato, V295 y ARCA diferida permanecen
pendientes. No garantía multi-tab, aceptación de negocio, provider, SCA/seguridad
integral, carga ni operación del target por este cierre.

## Reproducción e integridad

Staging `%LOCALAPPDATA%/Temp/elite-v311-d148716b286c44e9b0c7f4a257bc0263`.
Go1.26.7, Node24.14.1, pnpm11.19.0, Python3.14.4 y PostgreSQL18.6 fijados previamente.
gate.ps1 -Mode Go|Postgres|Browser -Root rebuilt; Browser -Project '' ejecuta4.
Web candidate compilada con archivos idénticos. Logs/traces/snapshots y scripts
sintéticos fuera de distribución. No datos privados ni llamadas a providers.

| Archivo | SHA256 |
|---|---|
| `src/components/franchise-command-panel.tsx` | `2a0c33893db1f2f46591cbb218728a1ff43ce06cd7be0d4cd65d51795aedd63b` |
| `src/app/api/enterprise/franchise/commands/route.ts` | `5b6721ec1b8cecf1be9e432aeaea8c7295e81479fd6dbad4c2225e90228c7523` |
| `src/app/api/enterprise/franchise/commands/route.test.ts` | `212474ba58fa2310963e85506c90cd5f71c61f21f78f79092e00d8c6d242d379` |
| `internal/franchisejourney/service.go` | `43ba7980cd8ba8c8e821efd25326d50aac3e21ef9d0468ea782dd85637fd9e30` |
| `internal/franchisejourney/service_test.go` | `1bf264f77d50bbed763c6072232fd6b7114022f9f94f41534de330dc4ae899bd` |
| `internal/platform/httpapi/franchisejourney.go` | `f9fd12a187515a7204672d0346c9eb3da2222ba955d0d24372e0bda636ac5af8` |
| `internal/platform/httpapi/franchisejourney_test.go` | `12704a41b45f412ddd14a9cb2e481b940537b4568362c39d5f73e53557fc732c` |
| `internal/platform/postgres/franchisejourney.go` | `e9650e3250ef9b10f1a4d7624e926f3b493f4ae2cc503015715e3972562481c8` |
| `internal/platform/postgres/franchisejourney_integration_test.go` | `8fac6b4b06dda49152f3126c3d4b4f985d8e51936617603a6ae82da4b3476085` |
| `microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs` | `7f41d23422623113f790af0e3f85014cf6f0c571a7750a4b6f8011d88c7447ad` |

| Log | SHA256 |
|---|---|
| `completion-red.log` | `87fc2f172d96970ccbe7bedec654f777f6650c56f0b7bb7737a77c162cd3eb7d` |
| `completion-r1.log` | `f7a2ab7f3d2bb690316c38e21860eba0379fd36ea85a209d9362f995a06cbadd` |
| `postgres-r1.log` | `a7a0f671e46c7bde4ca8092ec0bc7b3aa85e0b90278d69704888bbaadc82a73f` |
| `postgres-r2.log` | `bcab38e64a2cca243f8291ce3939f1b93725176c4db49f0b2116d156540514df` |
| `browser-r1.log` | `c3ba08e814ecfebd8fda11d1ff31174629b826731fd928b56b242781abe9a783` |
| `browser-final.log` | `8f4d4f0e7dd9593f884f9c65c256191992e4479f944b8bdde058bf243dba99e7` |
| `postgres-final.log` | `323e79573992402c0338b114ea281f4d62b56e8cf55633d3443bdc2190024e7c` |
| `go-test-final.log` | `73cacfdd85e00687cc4b4cb962ab649cc427dd72a9841113aa07abb338de06f8` |
| `go-vet-final.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `go-build-final.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `web-test-r2.log` | `dad326581fe062729159529d124b09fc7846ab7e6ef77ee7ec73789812bcf6ea` |
| `web-build-r2.log` | `d97fdb3641c3633d296eb6776816e8168b15150741ea856f2649870f4c9b5667` |
| `agenda-final.log` | `eaa18f194eb8d4c9c1f39fee4327e414e3028a7522699cb2805a30de7d12120c` |
| `compose-rebuilt.log` | `73e89a95b2841e0d811a640d59f266a8c0276c503f37f54b34dfbe1ce1f4d12d` |
| `restart.log` | `42de6d936c023711d990301d79700ce05ffeae998501649f41c1949e7dc44a73` |

## Cierre general

VERIFY_LIBRARY_PASS con160 packs/1436 archivos/744 Markdown y51 perfiles;
composición de franquicia67/745. Checkpoint49 validado antes del verificador.
Assurance TEST30/EVID21 enlazada: plan PASS, evidence integral BLOCKED por
pendientes previos preservados. El cierre local no equivale al100% del roadmap.
Log library-final.log SHA256: eb017f6fe134facaa8eddaa45f35861c9ddc5f36343c6cebfd1abcf5f7c17640.
