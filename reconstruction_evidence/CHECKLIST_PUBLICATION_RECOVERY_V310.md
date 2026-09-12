# V310 — recuperación de publicación de checklist

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

Mantenimiento T2802/T2803/T2804, FAIL457, desde checkpoint46 validado.
Fecha2026-09-08 UTC. Misma composición67/745; sin nuevas dependencias,
migraciones, reglas de entrega, proveedor ni efecto externo. No promoción/ZIP.

## Fallo y corrección

La publicación natural org/id/versión ya impide duplicación y mutaciones.
El rojo mostró actor vacío en el evento y404 al consultar la versión publicada.
El sender genérico podía afirmar no aplicado después de perder una respuesta.
Journey0.10.15 conserva el POST201 y el conflicto409 para identidad existente;
añade actor autenticado al evento atómico y GET /v1/franchise/delivery-checklists/result
con handover:manage, tenant/org/id/versión exactos y no-store. Una consulta SQL
devuelve sólo published y todos sus ítems ordenados en el mismo snapshot.
Service rechaza una proyección incoherente; no corrige datos al leer.

Portal0.14.8 usa referencia mínima scoped por sesión/identidad/organización:
ID, versión y SHA256 del contenido normalizado igual que el BFF. No títulos,
respuestas, prompts ni tokens en esa referencia. Fence antes del digest/POST,
fallo de almacenamiento0POST, bloqueo tras incertidumbre, consulta sin POST y
comparación exacta del contenido. Preparar otra versión requiere consulta
positiva y ausencia de POST/GET activo; no republica ni modifica la anterior.
Ayuda versionada checklist-publication-view/1.0.0 y resultados accesibles.
Browser0.1.20 añade el recorrido a las20 fases previas. Cambios AUTHORED;
REBUILD_VERIFIED/CONDITIONED no equivale a admisión integral de producción.

## Verificación

| Gate | Resultado |
|---|---|
| HTTP rojo→verde | actor vacío/GET404 pasan a actor checklist-publisher y GET200 con2 ítems ordenados |
| Web |113 PASS/1 SKIP explícito; build PASS |
| Go recompuesto |test/vet/build ./... PASS; skips sin DB no contados como integración |
| Dominio |12 proyecciones incoherentes rechazadas y positivo válido |
| PostgreSQL |8 tests,3 repeticiones; publicación8 concurrentes→1 alta/7 conflictos,1 evento con actor ganador |
| Atomicidad |outbox duplicado revierte template/ítems/publicación; lectura con nueva pool; tenant/org/id ajenos y draft rechazados |
| Inmutabilidad |versiones1/2 separadas; cinco mutaciones directas rechazadas, versión1 preservada |
| Browser recompuesto |84 fases/4 proyectos: Chromium desktop/mobile, Firefox y WebKit desktop PASS |
| Nueva fase |3 versiones/6 ítems/3 eventos/3 actores por proyecto; pérdida, JSON inválido, concurrencia201/409, storage0POST, recarga, GET503, contenido distinto bloqueado, permisos, actor forjado y versión anterior preservada |
| Agenda existente |4 proyectos PASS |
| Reconstrucción |10/10 archivos byte-idénticos;67 packs/745 archivos |
| Reinicio |snapshot templates/ítems/eventos idéntico tras reiniciar cluster sintético; no restore/PITR del target |

FAIL481 cubre actor/consulta ausentes. FAIL482 conserva compilación fallida por
import strconv omitido y corrección. FAIL483 detuvo el gate de reconstrucción:
el test JS generado con Python tenía CRLF; comparación demostró sólo907 CRLF→LF.
Normalizado antes de gates finales, fuente canónica intacta y hashes10/10.
No se retiraron guards ni se atribuyó PASS al intento fallido.

FAIL457 permanece OPEN para completar checklist y recibir/decidir retornos.
Ronda D canal/formato histórico sin respuesta; liberación comercial, V295 y
ARCA diferida intactos. Este cierre no demuestra garantía multi-tab, aceptación
de negocio, providers, SCA/seguridad integral, carga ni operación del target.

## Reproducción e integridad

Staging `%LOCALAPPDATA%/Temp/elite-v310-c131ab7f56f64bb9b124a3bec9b58e7d`.
Go1.26.7, Node24.14.1, pnpm11.19.0, Python3.14.4 y PostgreSQL18.6 fijados previamente.
gate.ps1 -Mode Go|Postgres|Browser -Root rebuilt; Browser -Project '' ejecuta4.
Web candidate compilada con archivos idénticos. Logs/traces/snapshots y scripts
sintéticos fuera de distribución. No datos privados ni llamadas a providers.

| Archivo | SHA256 |
|---|---|
| `internal/franchisejourney/service.go` | `29a83f4c8d95ddbdd3cf9f2c9b2bc61cb2e77a1d533e49a3bfac02961b9f8da4` |
| `internal/franchisejourney/service_test.go` | `d567fcd56a576fd0441d37823435ca1e7d368ae1e8506e29d9ef2829b2fad77d` |
| `internal/platform/httpapi/franchisejourney.go` | `7a5417c20fa322996b9a7b150f6d06442e151486f24e3dd7d3723439684845ad` |
| `internal/platform/httpapi/franchisejourney_test.go` | `dd894060994f89276738592513995ca243f1935288c12b1f8ada179e812cabb8` |
| `internal/platform/postgres/franchisejourney.go` | `5d725766a073f81894c78a9a16d68637b57f6a9f40f23801ffcc577d8bfe151d` |
| `internal/platform/postgres/franchisejourney_integration_test.go` | `fdf124f7e9b38f75aee380c7ebaba0ce29a839029f964bfad64acb61d776a3e4` |
| `microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs` | `60926c45f29d4d0b5e7482416682ee4d228a1c4b88742da77ad1fd01f6ec967e` |
| `src/components/franchise-command-panel.tsx` | `286c2a16c6ac04213db853cde2c734b292131d2e6f351384aa8efcfadfa4e035` |
| `src/app/api/enterprise/franchise/commands/route.ts` | `194523661df0c87aecdc1d82321991fa9b422282f59b82914d0fe6ad7c9a94e3` |
| `src/app/api/enterprise/franchise/commands/route.test.ts` | `0598f201a278782122a8875845cb130a14a96f8aa0115d2e4430177d574f540c` |

| Log | SHA256 |
|---|---|
| `checklist-red.log` | `aff5cef0869fe2076debe2edb11e67563240512caeebbec1a7cf3df973541c3b` |
| `checklist-r1.log` | `fd51db41b2b669511b3d7bc516f45ffdccea887a7f7bbab8c380d3214b0d6b76` |
| `checklist-r2.log` | `f75bca78f53995f1013c1fd36ab365bd1bfad4a4318097fd680d28f227821168` |
| `browser-r1.log` | `31c2e5eee7052f98cde07ea040308225fde3baefbca5fe5120b7675209005f53` |
| `browser-final.log` | `808dc1d5494772190d288e1f7590b82b2057268f284ae548f2c79df7edf29e37` |
| `postgres-final.log` | `16957f3c583e98989706b5c259e7b887c9a4c388e34d8f67b14b2c5da9517f8b` |
| `go-test-final.log` | `f8baf2bbeb46b57c5403648136ccce0d39512981f707e34a0249ea6d47788299` |
| `go-vet-final.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `go-build-final.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `web-test-final.log` | `1411c612a1aa4c3c5794a2f446deea565e5d05473e51709894fafa7c4e03d323` |
| `web-build-final.log` | `95988dd5736f101b003b3e30338cb0f00fcd3e9be92813a1b6e8859ebbea4474` |
| `agenda-final.log` | `ccfbace9e1e246b224f881c7af423dcbfd5590aa56e171e455e1c8f1a3c29197` |
| `compose-rebuilt.log` | `7b694f09704cdb500bc5aaa53287b884d4c53338b03a5a059076bd1c79b68a87` |
| `normalization.log` | `6325123dc508c9a92d5d3307c990eb9ce6e08ad35c1cef856e06fb7bfbef44ff` |
| `restart.log` | `15696e3119a85c1ce64260dfc9b2415a066ee372b84a4d9c94ea1a87bee5e41e` |

## Cierre general

VERIFY_LIBRARY_PASS con160 packs/1436 archivos/743 Markdown y51 perfiles;
composición de franquicia67/745. Checkpoint47 validado antes del verificador.
Assurance TEST29/EVID20 enlazada: plan PASS, evidence integral BLOCKED por
pendientes previos preservados. No se transforma el cierre local en100%.
Log library-final.log SHA256: 70391d430f8c521a3d7dc0b2ad5de576ee4a282aaf197bab95eccb6b26e80caa.
