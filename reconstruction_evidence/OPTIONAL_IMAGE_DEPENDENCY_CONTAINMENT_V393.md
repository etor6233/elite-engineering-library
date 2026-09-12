# V393 — dependencia opcional de imágenes retirada de la referencia

2026-09-11. Mantenimiento de biblioteca; entrada checkpoint248 validada antes de cambios.

**Cierre de esta revisión:** cuatro packs incorporados y reconstruidos; instalación limpia y regresiones completas PASS. Preflight250:164pasos ejecutados PASS y56perfiles compuestos; inventario165packs/1561archivos/841Markdown. El estado global de disponibilidad permanece BLOCKED únicamente por Docker ausente. La biblioteca conserva45/48controles PASS: TEST02/03/07 continúan pendientes.

## Cambio incorporado

TS-GO-API-WEB-BRIDGE0.5.16 desactiva explícitamente el optimizador de imágenes de Next y omite Sharp mediante ignoredOptionalDependencies de pnpm11.25.0. La referencia no usa next/image ni transformaciones de imágenes en su código de aplicación. La adquisición/inspección nativa diferida V386 sigue intacta; esta selección no es un fix ni una admisión de libxml2, Sharp o sus DLL. Sus archivos, licencias, fuentes y fallos anteriores se conservan como evidencia histórica separada.

El lock generado por el pnpm fijado pasa de152 a122identidades npm:30eliminadas, cero añadidas. Las122entradas restantes mantienen exactamente sus stanzas e integridades. Se eliminan exclusivamente Sharp y su cierre opcional exclusivo; dependencias compartidas se conservan. Instalación nueva offline/frozen/ignore-scripts PASS, sin Sharp ni sus paquetes @img. Una segunda instalación limpia sobre la reconstrucción canónica LF también pasa; sus805archivos fuente permanecen byte-idénticos y los12571archivos instalados restantes vuelven a coincidir con el baseline autenticado. Ninguna versión nueva, cuenta, proveedor, regla comercial, migración o código Go.

MICROSOFT-PLAYWRIGHT-BROWSER-GATE0.1.37 agrega al recorrido existente una comprobación HTTP real: /_next/image debe responder404. El baseline de producción anterior respondió400 al mismo GET sin parámetros, con proceso propio terminado. El nuevo oracle pasa en las cuatro agendas conectadas. GO-HTTP-METRICS-REFERENCE0.1.13 fija el hash del plan padre actualizado; su ejecutable Go es byte-idéntico al anterior. Cuatro packs reconstruidos, cuatro planes actualizados;805archivos exactos en la composición de referencia,68packs/797archivos en la integral.

## Evidencia ejecutada

- 294pruebas web PASS y1skip heredado explícito; typecheck y build PASS. El skip no reemplaza las suites conectadas ejecutadas por separado.
- 92fases de cotización/pedido/pago/entrega/devolución y recuperación en cuatro proyectos de navegador; pruebas de pestañas y sesión activas. Cuatro agendas adicionales PASS, incluyendo el404 nuevo.
- Cuatro recorridos privados de lectura/paginación/fechas PASS con cinco negativos HTTP y cero escrituras; snapshot de siete tablas preservado. Configuración argentina; no se presenta la variante Nueva York histórica como rerun.
- Cuatro recorridos de cancelación/reconciliación PASS:16cancelaciones durables,16auditorías/16outbox,28POST incluidos los rechazos. Sin replay automático de resultado incierto.
- 12.571archivos reales de la instalación restante coinciden con la instalación anterior;8.534pertenecen a Next, env y SWC ya autenticados y también coinciden con el inventario de artefactos firmados. No sólo package.json. Los paths largos de Windows se comparan íntegramente mediante su representación extendida.
- Consulta actual a OSV oficial sobre las122versiones del lock: cero avisos devueltos. Esta consulta no es un SBOM completo de cada componente vendorizado/nativo ni cubre todo el tooling.

## Autoridades y condiciones de uso

[Next: opción unoptimized](https://nextjs.org/docs/app/api-reference/components/image) y [self-hosting](https://nextjs.org/docs/app/guides/self-hosting) documentan desactivar la optimización. [pnpm: ignoredOptionalDependencies](https://pnpm.io/settings/dependency-resolution) documenta la omisión selectiva; además se verificó su presencia y ejecución en el pnpm11.25.0 exacto, sin asumir toda la documentación de la versión actual para esa versión fijada.

La observación de releases oficiales del11de septiembre sigue encontrando [sharp0.35.4](https://github.com/lovell/sharp/releases/tag/v0.35.4), [sharp-libvips1.3.3](https://github.com/lovell/sharp-libvips/releases/tag/v1.3.3) y [MXE8.18.6](https://github.com/libvips/build-win64-mxe/releases/tag/v8.18.6). Sólo se consultaron metadatos públicos; no se reprodujo ninguna vulnerabilidad ni se adquirió un nuevo binario.

AUTHORED/CONDITIONED: este perfil no ofrece transformación de imágenes en runtime. Antes de habilitarla se debe admitir una canalización compatible, sus versiones/procedencia nativa/licencias/SCA y gates de rendimiento web. No restaurar automáticamente el grafo excluido como rollback apto para producción. Para revertir técnicamente en investigación existen fuente/lock anteriores exactos; conservan sus bloqueos. Las imágenes estáticas futuras requieren su validación normal. La omisión de paquetes no acredita por sí sola seguridad integral.

Implementation assurance: contratos oficiales y package.json firmado; comparación del grafo y artefactos; pruebas funcionales/HTTP/browser/datos/recovery; procesos propios con tiempo finito y parada verificada; reconstrucción canónica y paridad. No se reejecutó Lighthouse ni se declara mejora de rendimiento, accesibilidad o admisión del tooling por este cambio. No existe una ruta de transformación que requiera dimensionar carga de codecs en este perfil.

## Estado y próximos gates

45/48controles siguen PASS. TEST02 conserva su integración funcional pendiente por owners; TEST03 conserva las condiciones del tooling pnpm, componentes vendorizados/nativos restantes y trust boundaries; TEST07 depende de esos gates y de la firma/verificación del artefacto final exacto. No se altera ningún oracle, status o clasificación REQUIRED para fabricar el cierre. V392 verifica un ZIP anterior; no se reutiliza su aceptación como prueba de este payload nuevo.

El aislamiento Sharp elimina su instalación de este perfil concreto. La restricción V386 permanece para su workstream de investigación, pero no debe describirse como si esa DLL siguiera instalada en el perfil nuevo. El conteo45/48 es de controles, no un porcentaje que haga equivalentes los tres pendientes a tres archivos pequeños.

La reconstrucción y compilación aquí registradas se completaron. Preflight249 rechazó el inventario descriptivo840Markdown: el informe nuevo lleva el conteo real a841. No produjo un resultado integral ni se omite su fallo. El owner se corrige a165/1561/841/56; El sucesor Preflight250 completo ya pasó sus164pasos ejecutados; sólo Docker permanece indisponible.

## Fallos de verificación preservados

FAIL719: rutas supuestas ausentes, recuperadas con inventario. FAIL770: receipt de composición contado como fuente, corregido a805másreceipt. FAIL768: Path.is_file rechazó un path existente de270caracteres; prefix extendido en ambos lados conserva las12571comparaciones. El primer literal generado fue rechazado por sintaxis y se conserva. El primer parser de comparación de lock incluyó ignoredOptionalDependencies dentro del último paquete; delimitación por sección corrige el parser sin modificar el lock ni sus oráculos. FAIL780: la primera reconstrucción rechazó los finales CRLF del candidato antes de modificar los packs raíz. Se conservaron esos bytes, se normalizaron cinco archivos a LF con texto decodificado idéntico y se volvieron a ejecutar las294pruebas, typecheck/build y todas las92fases/4agendas/4lecturas/4cancelaciones sobre los bytes finales. El lock no cambió de hash; OSV se reutiliza sólo por esa identidad exacta. FAIL781: la enumeración con paths ordinarios omitía16archivos en el candidato y27en la instalación limpia, aunque los8534archivos firmados se verificaban por path extendido. Enumerar desde la raíz extendida recupera exactamente12571archivos en ambos árboles, con mapas completos idénticos y comparación de cada hash contra el baseline. Los conteos previos12555/12544 fueron incompletos y se conservan, no se interpretan como archivos faltantes. Scripts y logs originales permanecen en el stage; no son fallos de los tests de producto.

## Receipts

Raw evidence se conserva bajo Temp, stage elite-v393-d68478838b4746f8b0f32686a9205535; resolver local elite-v393-current.txt. Ningún node_modules, secreto, binario o dato de consumidor se incluye en los Markdown portables.

| Evidencia relativa al stage | SHA-256 |
|---|---|
| official-release-observation.json | 9850f5f85cc5ae8c11f8c29766da5c3c1af1b1e52a8375fa38b4fa6fe22a8a7a |
| candidate-lock.json | b02d494128f70fe76031ea7086b856652bd33e2b98a74323fbc13df464c27484 |
| candidate-install.json | 27725cfc9615b8e82617df664e02b78921363e7b2ba585c3b70ae810c241e7db |
| candidate-tests.json | d805f7ca42265794bc222775a0bad00dfc586a5630dfdd2fe1237c8d15c80540 |
| candidate-typecheck.json | 03f7a035098a338800ebc616e73af8fe27a6754f2e23804319ad487ef5ea91f6 |
| candidate-build.json | a80c57f7cbe46e4ede6b447690a5ce77fcabc871320a859946809171f01b3e06 |
| exact-package-delta.json | 623db0e9e153688c3795a3d8cc1e9f7444825819f276c0fcc6ae1a5174d009fb |
| candidate-lock-sca.json | 8ef60fb11dd1ee3fc50030608559b2372bbb807237c45bb7c49a12a5cdb6e912 |
| candidate-osv-request.json | a1900f97c831512b070320c9fd4523964d896436cbc10cd7e174cda7f9841b9e |
| candidate-osv-response.json | a7ac2d1b1928b9dbd8070b7938b9420bf877dd39e923767a9b0b2498d9a2b396 |
| artifact-parity.json | 728e4eac8375e272cc5f3d2c73007b54fcbf685a4f977bcf01e2248ceb889cf0 |
| remaining-install-file-parity.json | 967b5e46ce270622aded04c200da6682667945008e0eb44ed34a5490d654d706 |
| baseline-image-endpoint.json | d6b8ea89d9e2d4120a2b28bd4972ce81ae899253a6190ade7a8afefb4c574511 |
| newline-normalization.json | d89dec98913763e82d99439fe6dcb05d5b0b9412226beaa521b8ee9330dc0d36 |
| final-gates.json | f22abaf77d3c5675490c1efc1205c666652e43eb9bda66ff2b546ad75bc45b6b |
| canonical-pending-changes.json | 2389fab63163646e66ba68afb31ffcf44362a4981246b3fc7305913efc23d289 |
| final-parity.json | 5514aca3fe339c4356538e62ef805ef171dd1df44100b1c4ce07d22b488aefd9 |
| complete-inventory-correspondence.json | 5ede8cce8d04409bcdd48a4c34c54e476b3547019202ef0dd88903cf0c97cf61 |
| final-consumer-install.json | ce48afae77d46dd3260530cf1736148504f09c510c7c2054e9675bf7c55c911a |
| final-artifact-parity.json | 728e4eac8375e272cc5f3d2c73007b54fcbf685a4f977bcf01e2248ceb889cf0 |
| final-remaining-install-file-parity.json | 967b5e46ce270622aded04c200da6682667945008e0eb44ed34a5490d654d706 |
| final-source-format.json | d193b4f7b1753e7a439706591017f6a7b2d3856b53b57140f3abce0ceb92aadd |
| connected-9dbacbdeb87243e7aca4d7f928e9053e/result.json | 5a2c0c12cd87072af9649c197dfd942a59a04c0994a4fdd0a132353461e18ca8 |
| connected-5997616ed8ca4d31943fe2da5edfcbb1/result.json | 61021a6452a64774217e4656fefc86572f9c922da37f0b1cb72f394d53c9b81a |
| connected-45c3e35654e4430392c401e2ff9028ba/result.json | d44eebca256c8dcb68f0b856cc07ba8434923c01f311c7b80d7a83d738113295 |

Preflight249 fallido por FAIL770: log SHA-256 a50f432541f3e5e0e507d9c6d9614bfc4a434e062667b4943cf7a37e686c7d78; proceso SHA-256 c819603ffe17bc7c511a88a1627a9b9b3868f56231cbd6d35dd09bc04ada8d9a. La fuente de producto, los48oráculos y la selección portable no cambiaron por esta corrección del resumen.

## Cierre verificado del sucesor

Se reconstruyeron TS-GO-API-WEB-BRIDGE0.5.16/58archivos, MICROSOFT-PLAYWRIGHT-BROWSER-GATE0.1.37/14, TS-MULTIROLE-ONBOARDING0.1.8/7 y GO-HTTP-METRICS-REFERENCE0.1.13/8. Multirole sólo alinea su guía con la versión BFF seleccionada; código y permisos intactos. Los siete archivos materializados cambiados y los805archivos de la composición se compararon según el contrato canónico, incluidas las diferencias generadas declaradas.

Preflight250 conserva allow_network=false, availability_is_admission=false y164steps PASS. Su status general es BLOCKED por Docker ausente; no se traduce a PASS global. Las suites reales de esta revisión usaron sus binarios fijados Go/PostgreSQL y bases sintéticas propias, sin Docker. Se compusieron56perfiles y confirmó165packs/1561archivos/841Markdown. El checkpoint posterior sólo agrega estos resultados; código, lock y oráculos no cambian.

Comparación estructurada contra el contrato de entrada:48tests,45passed y exactamente TEST-02/03/07blocked; acciones, oráculos y estados intactos. EVID-99 aporta evidencia acotada de dependencia/artefactos. FAIL780/781 se reparan con reconstrucción, instalación limpia y paridad completas; FAIL770 administrativo se repara por el gate íntegro250, sin cambiar selectores.

| Recibo final | SHA-256 |
|---|---|
| preflight250.json | 39d8c350d6b71c2dac72f51cd894cb5992fd0b481e1912ba579d0af78e2666b1 |
| preflight250-process.json | 00bd467630d4fe6cffe4f4ef0dc2d4324f2758d24b99a14611a81dabf366a4c0 |
| preflight250.log | 25a1e1fa87749711e3addf462232b79adb62b72123a5343a6dc7e013aef232c0 |
| acceptance-preservation.json | 3b945205812e92bbfff3f752595c01612a7be8cfebc8721df8734fc74a5d2193 |
