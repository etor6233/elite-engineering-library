# V396 — licencia original semver-utils y entrega completada

2026-09-11. Mantenimiento de biblioteca; checkpoint257 de entrada validado.
**Resultado:** cerrado el pendiente de localizar el texto original de licencia
semver-utils1.1.4 y entregarlo verificablemente en las tres recetas actuales.
El archivo original ofrece MIT OR Apache-2.0 y contiene el permiso MIT completo.
Se conserva entero, con la opción MIT seleccionada para el suplemento.
La biblioteca mantiene45/48: este cierre no cambia los oráculos ni la admisión
de TEST-02, TEST-03 o TEST-07.

## Evidencia de origen y discrepancia conservada

La consulta oficial https://registry.npmjs.org/semver-utils/1.1.4 respondió200;
sus2261bytes son idénticos a la metadata histórica321,
SHA256243bcee08073123a61e80d3955814c94ebe232b902c121f59dc67a3c0abe7bc7.
El manifiesto y registry declaran APACHEv2; nunca se sustituyó esa declaración
por una licencia inventada. El README trasladado y GitHub gitHead404 previos
quedan como historia; no se accedió por otra vía al host bloqueado.

La firma ECDSA verifica el mensaje exacto nombre@versión:integrity, y un mensaje
alterado se rechaza. La clave figura en el endpoint oficial
https://registry.npmjs.org/-/npm/v1/keys pero venció2025-01-29:
**HISTORICAL_SIGNATURE_EXPIRED_KEY**. Esto no establece firma vigente,
instante de firma confiable ni equivalencia Git/source/build.
Se preservan clave pública, firma, mensajes y hashes en los receipts.

HEAD no aportó Content-Length. Un GET Range bytes=0-0 obtuvo HTTP206,
Content-Range bytes0-0/4193, leyendo sólo1byte para fijar el tamaño de adquisición.
El tarball completo se adquirió después de reconstruir/probar el transporte,
materializar el perfil, responder todos sus inputs bajo la autorización de
mantenimiento existente y validar profile/approval/lock. No firma humana inventada.

## Adquisición de cuarentena y licencia exacta

URL: https://registry.npmjs.org/semver-utils/-/semver-utils-1.1.4.tgz
4193bytes; SHA256fa6980458971864f25f2afc7d7f6632b015d8b767296d48421e6ff13221bd2b3.
SHA5121239e82c4e4e1a60c0555ffc603a0de4a89a8cd69d8f3229f410073a161e407b7b8f4516c639a0b31e180f8f30a7851ed50a20ab7f05d4635424da85d6428ac4.

El transporte deja sólo el tarball opaco; ACQUIRED y PRESENT con rehash pasan.
El recibo histórico mantiene QUARANTINED_NOT_ADMITTED y todos los flags de
extracción, instalación, ejecución, admisión y firma verificada en false.
La inspección posterior, separada, descomprime en memoria con presupuesto1MiB:
19456bytes tar, seis archivos regulares,13868bytes de contenido. Se rechazan
rutas fuera de package, duplicados, enlaces, exceso de miembros/tamaño.
Sólo LICENSE y package.json se copian a evidencia; nada se instala/ejecuta.

package/LICENSE:1839bytes,
SHA25680b98c1b20edfc51abd2c802ee7a1d3c5561151d108a24f19c207b6935beaed8.
Contiene copyright2013AJONeal, elección explícita entre MIT y Apache-2.0,
los tres párrafos completos de permiso/condición/garantía MIT y el resumen
Apache original. No se presenta el resumen como texto Apache completo.
package/package.json:596bytes,
SHA256cf8e783b1772fbfb38d8994e1fc10879e9a2c6d7c170d792c768422aa8e69fa2.
package/semver-utils.js:3477bytes,
SHA2561f85bceb6f2e6cdbefadd5bf9969db8637c754c9d6553f3d7663dca6ade7c293.

El bundle pnpm fijado contiene un único marcador semver-utils/1.1.4.
Región byte1347913..1350984(exclusive), SHA256
aed20138b14658d55690a793c31b50d4f97e3c0d5260e323b7a74fb8279d9cf2.
Fuente original y región transformada son observaciones separadas;
no se afirma equivalencia de AST, bindings completos o build reproducible.

## Cambios canónicos y pruebas

OFFICIAL-UPSTREAM-ACQUISITION-CORE0.4.88:51files,196sources,25profiles.
Los195source records anteriores permanecen idénticos. Se extiende el owner
quarantine_wheel_transport para un único tarball npm revisado, SHA512 y4193bytes
exactos; no se amplía el roster opaco normal ni se debilita el requisito de
licencia de fuentes admitidas. Receipt npm separado y flags false.
0.4.87 continúa como candidato nativo aislado no promovido.

Cuatro suites pasan: upstream acquisition, source profiles, opaque transport y
quarantine. La última pasa140checks, conservando los74previos y agregando66:
identidad/URL/tamaño/digest/metadata/permisos/recibos/replay y no admisión.
51archivos reconstruidos byte-idénticos.

PNPM-ARTIFACT-SELECTION-GATE0.5.0 mantiene6files,67tests:
17selection+50routing. Recipe v4 genera SEMVER-UTILS-NOTICE.md y conserva
BlueOak/QRCode byte-idénticos. Recalcula el aviso desde constantes/fuentes
fijadas y rechaza ausencia, copyright/permiso alterado o falsificación conjunta
del notice y receipt, aun rehasheados. El bloque sigue ADAPTED:
lógica AUTHORED más texto original VERBATIM bajo opción MIT declarada.
No se incluye el algoritmo semver como nueva implementación.

Tres consumidores compuestos frescos desde el perfil canónico:
enterprise-web, Playwright y Lighthouse. Cada receta pasa preparación,
verificación y cuatro negativos reales;12rechazos comprobados.
Los442archivos de payload y su receipt quedan byte-idénticos a V395.
Las recetas no ejecutaron pnpm ni admiten runtime/redistribución.
El Preflight general sí tiene sus instalaciones offline congeladas previas;
no confundir ese alcance con la no ejecución de estas tres recetas.

## Fallos, límites y continuidad

FAIL787: helper asumió fences sin líneas vacías y su primera corrección insertó
signos plus; ambos errores se conservaron y el compositor los rechazó antes
de promover. Corrección limitada al helper, sin tocar gates. FAIL719:
consultas a paths supuestos fallaron; inventario identificó los paths reales.
FAIL788: licencia exacta desconocida y no entregada; cerrado con fuente original,
artifact utilizable y pruebas, no con una etiqueta APACHEv2 normalizada.

No se retoma V386/libxml2 ni el bloqueo Daybreak: queda diferido según la
instrucción del usuario. Quedan abiertas la integración integral, las otras
obligaciones/licencias nativas y de pnpm, seguridad restante y release.
No48/48 ni declaración100%lista. Rollback: snapshots de ambos packs y planes,
recetas anteriores intactas; nuevo material se crea en destinos ausentes.

## Evidencia reproducible

Stage: Temp/elite-v396-e5a797060309431d8e37863e7f5eeea3; resolver elite-v396-current.txt
Snapshots previos, helpers, metadata, recibos, artifact y notices se conservan
localmente allí. Los dos packs canónicos reconstruyen los cambios sin depender
de una instalación de semver-utils. Hashes de evidencia de este cierre:

| Archivo relativo al stage | SHA-256 |
|---|---|
| registry-observations.json | 55506eb3e24a8bd11e0bb8c5a7880a85321785d87b4d6a3148c6fbe230f9adeb |
| tarball-range-metadata.json | 33f5ef894760dcce2668bc8510f726539b1e9fef733278591fd180a8f192ab44 |
| npm-keys.json | faf23d8753d5bb79df250f10391ac89b63ecf7743e48487a544a99c847f9c8df |
| registry-signature-evidence.json | f5c7712a964d963af32ab77493e1495e39c433b81b11308e6d66de6152c26232 |
| license-inspection.json | b4de1791082b0b54c69a79aca4e1b363185ef1418af8c43cd7afa0b8ffc54e9a |
| semver-bundle-binding.json | a4b518c8f1c9e9454f4944eb19f57c8f24641d554ec9a4bd382fc89c5df791bd |
| core-rebuild-parity.json | cf28d18e580513bb9d7bb440491eab999d8121ceec195a468aea64d7f05cd54a |
| test_upstream_acquisition-run2.log | 65fc96b0158e5caf8168af58001d0cd7f02579f0b97076865792942e68fa19b0 |
| test_source_profiles-run2.log | 0ea7b0958f936263d117e4d7c8787aba4985683e0ba52d189c31daae76414cab |
| test_opaque_artifact_transport-run2.log | 79a97b5ada26e75b2d182ba58d615608364a4b90b9cb5829f6c526ba5853dbc8 |
| test_quarantine_wheel_transport-run2.log | 690f230c4f9e5a92d87e69d69d14909e55d58429a258da0764a86b5fa5bc1d2e |
| pnpm-rebuild-parity.json | efbe58a656b3f06db931b8643be88c5f625191fae0cfde2dd34ba8c41abdc645 |
| pnpm-tests.log | 9d973bba884b84cc03f478709a0a38e197693fbe99bf8585e3bc9c82b73f6685 |
| profile-approval.json | c45d05378960de4be7424266e4a3b539ac60fb273f851d811e0f5dd2e4368ab8 |
| profile-validate.log | c7f9d2fb0831c5417148313a3f3e8cff3ce9c09ced6bb20210ec7b0ec3ac6753 |
| profile-acquire.log | 4273f12db2cb25143618c7405c28a52a5ccb99bf13779981edd0a7c60e8d1427 |
| profile-replay.log | a0e1cd6ce9333b7d9bdf37cf40d3932fc431b93f04d0357199cdc6098b964a82 |
| real-delivery-results.json | d1d6d770326cad12f09f881f611bc73ab4e5be2a690965870c7c62d9d69bdb57 |
| projection-before-after.json | 4b56bb3ee5a5bfa16b7ac0832b67866ff5d5348e7357be81e1e3d200eca9f27b |

Preflight258 pendiente tras checkpoint. Aún no se atribuye PASS global a esta revisión.

Preflight258 no inició: FAIL790 por parámetro inválido en helper de invocación, sin gates ejecutados; log/process conservados. Reintento259 con argumentos exactos del runner anterior y sólo etiqueta de evidencia actualizada. FAIL789: selección real de30planes del core más1planner comprobada antes/después; no sólo inicialización. Los31pins se actualizaron y se preservan originales en before; composición global pendiente.

Preflight259 rechazó la ruta home absoluta en este reporte portable antes de ejecutar suites; FAIL791 corregido a referencia Temp/stage y pointer. El reporte anterior se conserva en stage/report-before-portability-fix.md. Preflight260 pendiente.

Preflight260 detectó FAIL792: seis oráculos fijos de composición necesitan+1 por el nuevo perfil AUTHORED. Se derivaron y auditaron los conjuntos exactos de archivos de los30planes que seleccionan core; se conservó igualdad estricta. No se cambia ningún oráculo de aceptación de los48controles. Ver composition-count-delta.json; Preflight261 pendiente. La metadata de fuentes del planner incorpora también la URL exacta npm y su plan declara67tests; fuentes ejecutables probadas intactas.

## Cierre comprobado — checkpoint262

Preflight261 completo:164pasos ejecutados PASS y56perfiles de composición.
Inventario165packs/1562archivos materializables/844Markdown/56perfiles.
Procedencia1307AUTHORED/148ADAPTED/107VERBATIM. Franquicia integral68packs/
797files y referencia HTTP69packs/805files permanecen iguales.
Disponibilidad global BLOCKED solamente por Docker ausente; el PASS de los
pasos ejecutados no constituye admisión productiva ni cierre de los48controles.

La licencia original semver-utils está identificada y entregada en tres recetas,
con67tests y12negativos reales. El transporte aporta140checks y57fuentes
canónicas idénticas a las reconstrucciones probadas. Los195registros fuente
previos,43archivos del core no modificados y todo el payload442 permanecen
intactos. El verificador sólo cambió seis conteos exactos afectados por el
nuevo archivo; sus comparaciones siguen siendo estrictas.

FAIL787/788/789/790/791/792 reparados en sus alcances acotados. Se conservan
helpers, reportes y ejecuciones fallidos. Ningún fallo se convirtió en PASS
mediante exclusión ni se cambió un oráculo de aceptación del contrato.
TEST-02 mantiene auditoría contractual integral de cores; TEST-03 mantiene
seguridad/grafo exacto y tooling/nativos restantes; TEST-07 mantiene build
reproducible, SCA, notices/SBOM y firma del artefacto final.
45/48 permanece explícito. V386/Daybreak sigue registrado como diferido,
sin reintento de esa parte.

Próximo alcance: consolidar la entrega de los avisos pnpm restantes ya
investigados y su cobertura exacta, manteniendo separados permisos, procedencia
de publicación y obligaciones de fuente/relinking. No reabrir semver-utils como
licencia desconocida ni repetir los endpoints bloqueados.

| Evidencia final relativa al stage | SHA-256 |
|---|---|
| preflight258.log | 62eba06f9c61e01f45d06365d35eb7a8b72738a24adcc777a2c968ae2c3443c2 |
| preflight258-process.json | 647ef17e66be9cc9c1a399810e1905c853c247adb6432db3c59f9c8737e9c784 |
| preflight259.log | 9a2d2aad0e3b8b2594206a857ce7556c7243087d2333661c0324b98c333b5e93 |
| preflight259-process.json | bb0c62dda078b89fe483cceb5ed92c5d494ae52be1eda871331d6c924332e03f |
| preflight260.log | 1e1c6e1a7a0643c723bbc1a8459c623261969b8a7c6743e78d3f18384aacfaa9 |
| preflight260-process.json | 13853afad4ac400eb74da51e6e1ea0b16cbaba681b19a7a0a5c553c1dea1af46 |
| preflight261.json | 5275e5d69dcca69b8da8e05dc76860436730556194e945a6b143c6299e26af86 |
| preflight261-process.json | 16d8917525171cbbadb24bd041b8bb0a44c5b590c0de3b247a4f0af4e8ea5296 |
| preflight261.log | c3aac4ede331a4dc193e4e010bc40c6845d14a3f44d05da2994295de7f06028d |
| plan-fanout.json | 29d18e96a750766976dde3a0d73fd99c490614347347d6449e831b9072309b26 |
| composition-count-delta.json | 19be9cc8551db1c79396412f6add78e67778c8e4b949bb95f6593c4b202adcb8 |
| scope-delta-audit.json | d91c5e96db8ed73ccd0ee0839f11961651736a184ffda1702e59bc739cfcba1b |
| final-canonical-source-parity.json | c38fb1443b5de2fed723a1b7b4128baa660519c58d5060866840a1e320091fbb |
