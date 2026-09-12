# V313 — actualización de seguridad y cobertura de dependencias de la composición

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Los cuatro binarios retenidos V313 prueban Go1.26.8 mediante buildinfo. La igualdad before/after permanece demostrada; la atribución1.26.7 era incorrecta. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

2026-09-08. Mantenimiento T2803/T2809 desde checkpoint52 validado.
Perfil67 packs/745 archivos, sin nuevos paquetes de producto, migraciones,
reglas comerciales, cuentas, scheduled tasks, mensajes, deploy ni gasto de servicio.
Los cambios son locks de dependencias y un guard de test descartable. No ZIP
ni promoción integral. Blueprint y readiness BLOCKED siguen gobernando.

## Hallazgos y sucesión de candidatos

OSV oficial2.5.1, commit c84fa4568f2526d0333e9a914ea8a0a5f74ad68b,
SHA Windows25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6,
detectó GO-2026-5970/CVE-2026-56852 en x/text0.29.0 de los dos consumers Go.
El advisory oficial describe falta de progreso ante UTF-8 inválido y fija0.39.0.
Se consultaron advisory y patch antes de cambiar locks. Input oficial f3cc80:
el probe finito falla en0.29.0 por falta de progreso y pasa en0.39.0 conservando
bytes. No se afirma alcance explotable del producto sin prueba adicional.

La actualización mínima x/text0.39.0 eleva x/sync0.17.0→0.21.0 por MVS.
Primer candidato0.4.1/0.1.1 y sus receipts se conservan. El scan directo quedó
limpio, pero ampliar a go list -m all detectó x/mod0.37.0 seleccionado por
herramientas upstream, ausente de los paquetes compilados de ambos consumidores.
GO-2026-6179/6180 exigen x/mod0.40.0. Consultar x/text0.40.0/0.41.0 mostró
que todavía seleccionan x/mod afectado; no se actualizó por novedad.

El sucesor canónico es GO_ENTERPRISE_BACKEND_CORE0.4.2 y
GO_OFFICIAL_RETURN_REFUND_WORKER0.1.2: x/text0.39.0, x/sync0.21.0 y piso
explícito x/mod0.40.0. Comments de go.mod explican conservarlo al mantener
locks; tidy puede retirar un requisito no compilado y reabrir el gate completo.
No hay backport local ni código atribuido falsamente a upstream.

Los tres módulos conservan receipts de proxy.golang.org/sum.golang.org,
commit, sum y SHA de ZIP/LICENSE. Licencia BSD-3-Clause verificada íntegra:
conservar copyright/condiciones/disclaimer y no endorsement. No se distribuyó
un binario de producto. Go1.26.7 supera el mínimo1.26.6 corregido de cmd/go
para esos dos advisories; no se infiere compromiso de cache ni se ejecuta la
limpieza destructiva sugerida por el advisory.

Fuentes primarias: [GO-2026-5970](https://pkg.go.dev/vuln/GO-2026-5970),
[GO-2026-6179](https://pkg.go.dev/vuln/GO-2026-6179),
[GO-2026-6180](https://pkg.go.dev/vuln/GO-2026-6180).
Patches oficiales retenidos: CL794100, CL814960 y CL815000. Los redirects del
browser fallaron; la API pública Gerrit devolvió patches verificables, no se
atribuyó contenido a las páginas inaccesibles.

## Cobertura comprobada

El scan recursivo no detectó seis locks Python de wheel URL+SHA; forzar el
parser requirements.txt produjo exit128/0paquetes. FAIL495 no se ocultó como
scan limpio. Se adquirieron47 wheels únicas (53843594 bytes), se verificaron
sus SHA y METADATA sin instalar/importar y se proyectaron90 pins por consumer.
Cada fila conserva lockSHA/línea/URL/wheelSHA/METADATASHA. Se verificaron tags
compatibles y134 requisitos activos, incluidos extras, para CPython3.14.4
Windows AMD64: cero incompatibilidades o ausencias en ese entorno fijado.

Dos proyecciones Go conservan19/18 módulos seleccionados;9 módulos en cada
consumer se usan en paquetes/tests. No equiparar grafo de selección con código
compilado ni eliminar findings por esa diferencia. Se retienen ambos grafos.
La pauta vuelve a ZERO_COST_VULNERABILITY_MONITORING_PROFILE.md: contrastar
extractores por consumer, proyecciones trazables y pisos transitivos explícitos.
[Uso oficial de OSV para archivos/proyecciones](https://google.github.io/osv-scanner/usage/scan-source).

El scan final directo de rebuilt-final + proyecciones verificadas devuelve
exit0:21 fuentes,465 ocurrencias,347 paquetes distintos y cero hallazgos
conocidos en la consulta. No se suprimieron resultados ni se ejecutó osv fix.
Los JSON originales, fechas, fuentes y códigos de salida se conservan.

## Gates y equivalencia ejecutable

| Verificación | Evidencia |
|---|---|
| Artefactos |3 módulos Go oficiales, commit/sum/ZIP/LICENSE;47 wheels por SHA y METADATA |
| Regresión norm |input oficial provoca no progreso en0.29.0;0.39.0 PASS con bytes conservados |
| Upstream |norm, errgroup, semaphore y singleflight PASS; suites oficiales x/mod sumdb/note/dirhash/tlog PASS |
| Go ambos consumidores |test/vet/build -mod=readonly y mod verify PASS en el sucesor |
| Backend PostgreSQL |11 tests por3 repeticiones,33 PASS raíz; concurrencia, scope, atomicidad y recuperación |
| Refund PostgreSQL |2 tests por3 repeticiones,6 PASS raíz,53 migraciones en base descartable exclusiva |
| Web |115 PASS/1 SKIP explícito y build nuevo frozen/offline |
| Recorridos conectados |92 fases/4 proyectos y agenda4 PASS; sin proveedor real ni efecto externo |
| Equivalencia del sucesor |317 archivos Go idénticos,388/260 variantes de paquetes/imports compilados idénticas, API/refund binarios byte-idénticos con trimpath |
| Reconstrucción |5/5 archivos cambiados desde fuentes canónicas;67/745 |
| Fuzzing |GO_NATIVE_FUZZ_GATE0.1.0, target de dominio ReturnCaseResult/8 semillas, presupuesto10s,2 workers, PASS; receipt del sucesor |
| SCA final |21 fuentes/465 ocurrencias/347 paquetes distintos, cero findings conocidos; cobertura separada de runtime |

Las92 fases corrieron sobre el primer rebuild con x/text0.39.0/x/sync0.21.0.
El piso adicional x/mod no cambia fuente Go ni grafo compilado y produce los
dos binarios idénticos citados; esa evidencia permite transferir esos recorridos
al sucesor sin repetirlos por una dependencia no compilada. Test/vet/build,
mod verify y fuzz del sucesor sí se ejecutaron; agenda/refund finales usaron
su fuente. No se transfirió un PASS a código funcional distinto.

El test histórico de reembolso omite triggers para seed/cleanup de un tenant
sintético. Se agregó guard de URL literal loopback, base elite_refund_test con
identidad única, rechazo de fallbacks remotos y timeout30s. Ocho targets negativos.
Sólo corrió en base nueva aislada, nunca en la base auditada de V296. Su PASS
prueba compatibilidad del worker, no integridad del journey que origina requests.
Mejorar ese fixture sigue T2802; no habilita pagos ni reglas comerciales.

## Fallos conservados

FAIL494: dependencia vulnerable, corregida por artefacto oficial. FAIL495:
cobertura Python cero, corregida con inventario verificado. FAIL496: guard de
licencia esperaba Google Inc. y la fuente dice Google LLC; lectura íntegra y
receipt sucesor, sin cambio de licencia. FAIL497: wrapper esperaba timeout pero
el guard del probe detectó antes la falta de progreso; rojo finito conservado.
FAIL498: test de refund aceptaba URL arbitraria para cleanup; guard y base aislada.
FAIL499: identidad web comparó template next-env con imports generados; se
conservó diferencia y se construyó de nuevo, sin editar derivado para fingir hash.
FAIL500: grafo seleccionado reveló x/mod afectado; piso explícito y rescan completo.
FAIL501: -match de PowerShell confundió escenarios failure exitosos con FAIL;
parser anclado y case-sensitive, cinco probes, resume sólo de gates pendientes.

## Límites y continuidad

Este cierre es de dependencias fijadas y sources identificadas, no certificación
de binarios externos/vendorizados, biblioteca estándar, herramientas descargadas
por otro gate, imágenes, runtime desplegado, call reachability, SAST, DAST o carga.
OSV no prueba licencias ni vulnerabilidades desconocidas. Dependabot/periodicidad,
retención/alerta operativa y monitor runtime no están demostrados. No se registra
Scheduled Task, no Actions programadas y no se declara ZERO_COST_BASELINE_READY.
Readiness y assurance integral permanecen BLOCKED; Ronda D/canal histórico,
liberación comercial, target/IdP/proveedores V295 y ARCA diferida siguen intactos.
Siguiente owner: T2809 montaje host/supervisor y señales; V301 histograma ya
reparado, no repetirlo ni incorporar un core candidato como cobertura integral.
No porcentaje de cierre inferido; T2801–T2810 siguen con sus criterios propios.

## Receipts y reconstrucción

Staging: `%LOCALAPPDATA%/Temp/elite-v313-36edb2a519a84574bda36e43c6ad46b0`.
Baseline de cuatro locks, candidate1-locks, baseline-packs y candidate1-packs
preservados. El rollback vulnerable sirve como evidencia, no autorización de uso.
Scripts y artefactos quedan retenidos; no se borran registros para ocultar rojos.

| Archivo en reports/ | SHA256 |
|---|---|
| intake.json | `375efc1608e35bb52d7148b15870e2dcfda160842bc22186aa4d9722a34fa51b` |
| source-scan.json | `77b129664c0b2d901970855dec0921c8ccbcf3985168addecb8f6bd36c886dd9` |
| source-scan-exit.json | `0f45bd44e84addd2f1349a02861f2ec8c5489063539a5cf98a99544f3fc56fde` |
| python-scan.log | `8de7a648e09d58a2114822e6f82f7dc3fc1e8bd572071152dbad1aec17a627b3` |
| python-scan-exit.json | `986a99af0c0c1b4b096e87132d4db13cd05bd0c12ea5ca80ca9c8c057954183c` |
| python-artifact-inventory.json | `3a2f9824f963cabca87e781f09a67d0586eaef9ee3016a5e8d4bcd20189a9596` |
| python-closure.json | `0e831cf207e9e7d6de4c56fd1cc002a6c87af11985697bc978f86ddad97eb77d` |
| python-projection-scan.json | `5c2cc5e296e059ff526c603e8cbf7551725c6a4debc3896a9eafa570beedf575` |
| norm-before.log | `61ed6661e1387f0e0acb6f32140b2edb8a203e8ed6b1784cc5c13e73766ddc04` |
| norm-after.log | `6a6796dba8984ee4a7a899df4de32c844787b4ae07e64b5172a531f6ec7383df` |
| norm-regression.json | `a8b151dd4cce901886cb086654673ca77689da75e240b0c13edca3a5d2303ecd` |
| xtext-artifact-review.json | `b0ce1cb27bc8db8d161a9b46545f2329a1340bba4d430d137039e6d3acd5f913` |
| xsync-artifact-review.json | `f1858c449e27407ca6d0ee5b4041cc3cdd6d527033ba70d3963e4d9048c82b10` |
| xmod-artifact-review.json | `ad43baf260e719261bcc69b58ff784356541b118dfae797949152d78e48adab2` |
| official-794100.patch | `b790922e2d2b5e9d3f33cd3e490532b5302ae4dae69a8409d9b3ee404996ff14` |
| official-814960.patch | `d011aa094f758f759ace8e7974a1c1bc79c3fbdb037623fb309c69f35701cca1` |
| official-815000.patch | `4bf95ae68028bde6a3b658872ca6e7a57649a3537f42aca6dc30a37f98f74d1a` |
| upstream-package-tests.log | `6280c6e0b6ed8f651c46f89ce1cc9df8cd99ac2cd17b3c3256a8e380fe03d824` |
| xmod-upstream-tests.log | `ff2e335a47f7ec3c9ece53b387779c0e245dc44a64a4ecd46e5b484754af3d1c` |
| final-combined-scan.json | `301459e854ecc32d0436e8801e946dad9dd98537e418e308a54d8d582f9b5b02` |
| successor-combined-scan.json | `5f94a43f98abcb497fbaaa5010f8b9005791e7ac56cfa3effc547423fabc87ad` |
| successor-combined-scan-exit.json | `9d74e43c50ca51695d8d466438fcdca64f1c45c5b8963713a7bb0912406692c0` |
| scan-coverage-review.json | `b0ee75e192d3e10adcae0ff15e5484d96ea272a559e7e5b248389a3960e1cc82` |
| go-module-closure-successor.json | `bcaa93e3ec8a4ede8dc7b81b567e9c8eca71918636352212d2b127650121f19b` |
| compiled-equivalence.json | `e48f7bba8edb450699e4e36d773ea8d966588f0a2b9115ea124a6e7a261fe049` |
| backend-test-successor.log | `15833068351eb9aa3a755eaca0e18fd40e708ff12edb0d1ac8fbac26587cbdc5` |
| backend-vet-successor.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| backend-build-successor.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| backend-mod-verify.log | `acdb6a6a98dc2c297af31bb3538319778be8c8ef263cd0a8c9799de9c4998533` |
| refund-test-successor.log | `e431b25d4199952cc1a6e60862214465862e92b36d574270097070adf155e221` |
| refund-vet-successor.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| refund-build-successor.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| refund-mod-verify.log | `acdb6a6a98dc2c297af31bb3538319778be8c8ef263cd0a8c9799de9c4998533` |
| postgres-final.log | `32635a409db982055a81be8100d11fe35e01df7f8488e10902cd301d3f2c4e7e` |
| refund-postgres-final.log | `c64cfa9e075bad2ade76d6d5cc797f73d4471586394aec85ad134c92fa13d7f8` |
| refund-database.json | `e72523be0fc16d8f4bfbcdb5fe498ea0dbfb6662437c6846489b9126f7100a34` |
| web-test-final.log | `c1e378d53eb59225b3e34e4fb236e75a2f962f848591595aa38f240b0a6f0222` |
| web-build-final.log | `c8e3b0afe2369a5177cdea60e76c2ff802b34895b410db30723611b585c5b07c` |
| browser-final.log | `5875ed96b5a0a2471ec1c0dbbd01970f9197bc5a52ff5493193b0a30759caebb` |
| agenda-final.log | `6f2dbf654280c663575438b615ca612d8415a755e0dff4dd6df7bcff58921047` |
| canonical-successor.log | `cb90a872233531f59b290739cc3b9a65d3d3c05c7a06033aaa2344df151a7611` |
| fuzz-receipt.json | `b841af88c6138d8a25b21c0ad076401212b99df30521e359d0c0adbd50ebb250` |
| fuzz-successor.log | `151ffae07332849a01150b75b10dbaf0ab4cc067fde19019f1f1745d340377bb` |
| receipt-parser-regression.json | `6de471128b5d6e8ee5580478b96cc67ed819b912bf3fcc15e215d0d33e89928f` |

| Fuente reconstruida | SHA256 |
|---|---|
| return_refund_worker/go.mod | `32f96cb7d5255316b23ab3f597d5fa6ef508a64fd444c0ece5cd9e3ca9803aba` |
| return_refund_worker/go.sum | `ea80bd93ebcdd4c1cd14c9b759c419f1c903ccd007b57014e16827fbfb14a13a` |
| return_refund_worker/internal/refundworker/postgres_integration_test.go | `1156f3214404c84d646d87033886869b814b339a0588ee327067e170726d0b36` |
| go.mod | `c2011001bc9840c3e54a5e162601fbaea798800ac5f5af45e4bd7a7c997ff17b` |
| go.sum | `ae7cd054cab0711086c2d8d98e7171744d11b9a5e3bbcf91193f67e18fe8ee1c` |

## Cierre general

VERIFY_LIBRARY_PASS con160 packs/1436 archivos/746 Markdown y51 perfiles;
franquicia67/745. Checkpoint53 validado antes del verificador. TEST32/EVID23
enlazados, plan PASS y evidence integral BLOCKED. No equivale al100% del roadmap.
Log library-final.log SHA256: 578b857caae19924e56edef39292f1469c04ae558d61420cd600364879e79e0b.
