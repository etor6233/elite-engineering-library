# Gate Node oficial con datos verificados y ejecución acotada — V320

Mantenimiento T2803/T2809,2026-09-08. NODE-OFFICIAL-RUNTIME-ADVISORY-GATE0.1.0:
10 archivos,7 AUTHORED/1 ADAPTED/2 VERBATIM. Perfil focal1/10 fuera de franquicia
67/746. No modifica producto, instalación global, motor de OSV/Dependabot,
scheduler, provider, cuenta, presupuesto ni decisión comercial.

## Resolución del gap de V319

El CLI oficial Node1.6.1 aceptaba corpus{} como limpio y carecía de deadline
interno; los rojos V319 se preservan. El nuevo adapter cambia exclusivamente
getJson a lectura local: matching, filtrado de plataforma y EOL oficiales se
conservan. Comparación before/after demuestra igualdad fuera de ese bloque
y dos comentarios explicativos. El código no implementa un matcher alternativo.

Fuente [nodejs/is-my-node-vulnerable](https://github.com/nodejs/is-my-node-vulnerable/tree/c37a56bad56e34fe5223ddd3cb223cc4158136ae),
commitc37a56bad56e34fe5223ddd3cb223cc4158136ae,MIT. ascii.js y LICENSE conservan
bytes originales. Los7 wrappers/lock/perfil/ayuda/tests son del workspace y no
se atribuyen a Node. semver7.8.5ISC es dependencia externa, no53 bloques embebidos;
[release oficial](https://github.com/npm/node-semver/releases/tag/v7.8.5),
commit6e05b7637396ac66522cff8731f07cfe0ef49a29. Node24.20.0 independiente y LICENSE
son inputs externos fijados, adquiridos/verificados en V319; no se incluye npm.

## Controles antes de ejecutar

Perfil con esquema estricto, campos únicos, enabled=false seguro y edad1..7días.
El lock queda vinculado por SHA al runner; su fecha no se puede refrescar como
input del perfil. Corpus/schedule deben coincidir con los hashes de los snapshots
oficiales194 registros fijados; vacío, truncado, alterado, vencido o fecha futura
falla antes de ejecutar el motor. Actualizar datos requiere nueva investigación,
before/after, lock/binding, gates y revisión de pack.

Semver verifica su conjunto exacto53 archivos/licencia y cada SHA; rechaza extras,
faltantes, vínculos y más de256 entradas. Límite1MiB por archivo y4MiB de bundle.
Motor, bridge y dependencias se copian desde bytes ya verificados a un directorio
temporal aislado. Sin consultas HTTP del motor; bridge bloquea http/https.request.
El Node externo se comprueba por SHA, versión/platform y LICENSE antes de usar.
El entorno no hereda NODE_OPTIONS/NODE_PATH/DEBUG ni argumentos privados.

Child20s, identity probe5s, heap V8 de64MiB y salida combinada64KiB. Timeout,
salida excesiva y exit no cero se convierten en error estático; no se publica su
causa privada. Estos límites no son sandbox del OS ni límite de RSS completo.
Se confía en el código del pack instalado y el host local, no en un admin hostil.
El receipt nuevo vincula runtime, perfil, lock, snapshots y fechas y declara scope
y production_admitted=false; nunca pisa evidencia existente. El proyecto conserva
triage, retención y renovación. No se crea ciclo persistente ni destinos de alerta.

## Fallos encontrados durante la implementación

FAIL519: el primer runner verificaba motor pero copiaba bridge sin digest. Un
bridge sustituido por JSON fijo produjo PASS sin consultar el motor. El rojo
se conserva; agregar bridge_sha256 al lock y exigirlo antes de copiar corrige la
frontera. La suite incluye sustitución del bridge y del motor.

FAIL520: semver del ZIP Node no es byte-idéntico al tarball npm del mismo número.
Comparación de53 archivos:52 difieren sólo por CRLF/LF y README por whitespace
adicional. No se normalizó para fingir igualdad: se seleccionaron los53 archivos
con sus bytes exactos del tarball oficial fijado por integrity SHA512 y SHA256
d85045d4300d7d57c891336b95df532e73f34c22ffcd222452b6d08b9d127d5d.
Los bytes y hashes anteriores siguen disponibles; licenciaISC original preservada.
Tarball/metadata exactos y las30 pruebas se verificaron sobre esa selección.

Recurrencia FAIL491: el writer Python produjo CRLF en archivos locales y el
compositor rechazó el hash del primer bloque. Se preservaron before/log,
normalizaron sólo terminadores del candidato, recalcularon lock/binding y se
sincronizaron10 bloques mediante update_pack_from_tree.ps1. Los dos archivos
VERBATIM ya tenían LF y no se alteraron. Se recompuso en destino nuevo, sin
repetir el writer parcial ni cambiar el verificador.

## Evidencia ejecutada

30 pruebas PASS sobre la selección final y sobre el árbol reconstruido, sin skips.
Incluyen control vulnerable24.14.1 del motor, Node actual24.20.0, plataforma inválida,
vacío/truncado/oversize/stale/future, perfiles inválidos, bridge/motor/semver/lock
alterados, licencias faltantes, runtime ajeno, receipt existente y env hostil.
Children reales demuestran timeout, exceso de salida y causa privada excluida.
La frontera junction usa inyección unitaria; no se afirma sandbox Windows.

Cuatro CLI reales desde rebuilt: actual exit0/receipt, disabled exit2, corpus{}
exit2 sin receipt y destino existente exit2 con bytes intactos. El control de
versión vulnerable se prueba en el API del motor, no como ejecución del viejo
runtime dentro del gate. No se atribuye ese caso al CLI final.
Rebuild10/10 idéntico al candidato; modificación upstream acotada comprobada.
OSV2.5.1 exacto, all-packages, grafo CLI is-my-node-vulnerable1.6.1+semver7.8.5:
2 paquetes/0matches conocidos. No incluye Actions entrypoint, npm/pnpm compilado,
CPython, todos los componentes nativos ni security testing del deployment.

G0–G8: PASS para la instancia compatible actual; resolver emite
USE_REUSABLE_PACK/ready=true. El pack global REBUILD_VERIFIED/CONDITIONED mantiene
inputs exactos, licencia, freshness y ejecución de gates como condiciones.
El receipt de resolución no convierte en admitido al CLI upstream sin guards,
no valida el alias OSV npm/node y no certifica monitoring o producción.
README/source-lock/plan permiten materializar y operar el gate sin nueva aprobación
ceremonial cuando el proyecto satisface esas condiciones; lo demás sigue bloqueado.

## Continuidad

TEST39/EVID30 registra esta implementación; assurance pasa de AUTHORED a MIXED
con source lock explícito. Se corrigen además los blockers históricos que aún
describían observabilidad0.2/histograma pendiente:0.3 y reparación V301 existen,
pero integración/política/operación target siguen pendientes. Conservar historial.
La siguiente auditoría aborda pnpm11.19.0 compilado y otros runtimes, seguida de
hosts/supervisión/alertas/retención T2809. D/canal histórico, liberación comercial,
target/IdP/providers y ARCA diferida mantienen sus decisiones pendientes.
No ZIP/release/deploy ni porcentaje de completitud inventado.

Staging `%LOCALAPPDATA%/Temp/elite-v320-513682952df14bd2b01959d4ff04e3d8`.
Scripts de preparación/corrección son de una sola ejecución; los tests/gates usan
recibos y destinos nuevos. V319 y los rojos/locks before no se reescriben.

| Receipt | SHA256 |
|---|---|
| bridge-red.log | 6ef68c939d3c29cae9508c3458f243699aa6cd438c2e3edc46817992b7db4e1c |
| bridge-green.log | 0b6f194bf98cffe2da458ae70d186c7c0bb8e4ff8aaea3e43591558a39c6f57e |
| candidate-tests.log | c23f4d67293d8b4f75bed725edcc3ba3b93f35be1528a484f20198ffcca5c17f |
| official-semver-tests.log | 35097fd8c349a889ab263941be13fc985f9c51228745dea08cee38bf85def934 |
| candidate-final-tests.log | 33cbb9e6ceb5af4d6301a85267f9be31a3725348c1ace58900b6cd39582a5190 |
| rebuild-compose.log | 436715eca87e4050290ae2957127ea3a964d574486b1d1f9f9b3d42ac44ee686 |
| rebuild-lf.log | 5b48116417e9e7d3fdfa1baabc422b3557dfef02c47c12d6345345573df09b82 |
| rebuilt-lf-tests.log | b61076b873a3d5929aa11b4ab181d37b553c492719e4a9d204a40cb59914d611 |
| semver-comparison.json | 11b73b5ef0da73a30bc1ea6a411d11d281938cd71074b805b4310c5be6ac5aa0 |
| semver-readme.diff | 53565668f542001ea9b4aef17f7bac5b28a5191de41748d6dabca153c4d943bc |
| semver-official-files.json | 1a84006001e9bb02eda76e81bb4dd0375b3cb081ecb51abc3f5f0ad88884e726 |
| semver-artifact/registry.json | 595600ec8060055e4582121441e8227bc6ab12ff63b42030686e18ec46226ff2 |
| sca-results.json | 56f9a249745c1f1a2773019c0ef6f616c45793ea6152091e04e373adfd4a10b9 |
| sca.log | 70b2e88e288fb97b22f521ed5f4d63419165fa3ce3741310ffeaa0c8c0b342bb |
| normalization.json | 8c132f6a118e71691e84bd1728252059a04c21d0505e28c6cc41e5d7125fa131 |
| rebuild-parity.json | de7f10f488008c347dfcb3c99ef63e3e2fa27859f624e462b362016fe5f32076 |
| cli-results.json | 8d9e629d5709b6d8f572bdeb0f0c6ccf649b9979cf238a2b4a3849cb20117d01 |
| valid-cli-receipt.json | 4867e2b90bbee5ffb317678744da542233eef4fc4b6763989ebc88595a2e66ed |
| qualification/qualification.md | 752174fba3e4e09ab9fb44f52094f9aa492d26d919f4a418c33556146819893e |
| node-runtime-gap.json | 3d3ed9c22e3c3d9010f090b43162a95b3064e8e1a715c8c3c6477adaa79577fd |
| node-runtime-gap-receipt.json | 88a0ec365ed2161be8344c665efa5800499aefe270e774d0f76277671f64220f |
| rebuilt-lf/node_runtime_advisory_gate/profile.example.json | 29cb59f281cd9fa40d8d994a951d43935eb39ff31cffea373edcba11e019d881 |
| rebuilt-lf/node_runtime_advisory_gate/source-lock.json | 98ee8520ec08f0d6dd988a5d477402dfdb99f9b13ae0b1c2a27c6767a70db247 |
| rebuilt-lf/node_runtime_advisory_gate/run_gate.py | 510cfe27a052ba803d8748091772f2727212fd0ae20a55e6a7fb27db4b587123 |
| rebuilt-lf/node_runtime_advisory_gate/bridge.cjs | 2e6a524bf2ae602a02f5fdf2ef59310e174490540d37afc56b559f90e91b0fb6 |
| rebuilt-lf/node_runtime_advisory_gate/verify_pack.py | 6368e937709d430fc3a7993de68b57d97146645a8af602d10ea63c9227f5debe |
| rebuilt-lf/node_runtime_advisory_gate/README.md | b5c83e4e571f364ad8f64412dd62da69d33c533f1abfbb1a5b04f505883e1d99 |
| rebuilt-lf/node_runtime_advisory_gate/THIRD_PARTY_NOTICES.md | ecf555b3fb96298493d253287132691844ce592335909e5d2ccd2e6f0092d586 |
| rebuilt-lf/node_runtime_advisory_gate/official/is-vulnerable.js | 769e22cd956780430210535f562d1d31cb6fda5cdfe22790980b0e9595678c46 |
| rebuilt-lf/node_runtime_advisory_gate/official/ascii.js | a22ca24d5ca9dd023219f32ed668359a15493ac6e351d3640538067ceb4761bc |
| rebuilt-lf/node_runtime_advisory_gate/official/LICENSE | bbf7df376bad39f8be96f9044ad58b77f0a48cbd1b482ecc5893916f7129b91c |

## Integración en el verificador raíz

El primer VERIFY_LIBRARY validó checkpoint68 y todos los packs, pero rechazó52 planes descubiertos/51 compuestos. FAIL522 corrigió el registro explícito del owner: composición Node1/10 y ID exacto, sin bajar el gate de cobertura. Log rojo library-final.log SHA256:de0f9f19dd144e4a440cda7b1e1c66ae8a49bcfbcdc38a248fe346bb985a082d. Before verifier-before.ps1 preservado.

## Verificación final y continuidad

VERIFY_LIBRARY_PASS:161 packs/1447 archivos/755 markdowns;52 perfiles descubiertos y compuestos, Node1/10 y franquicia67/746. Receipt library-revalidated.log SHA256:5300dec9f292a7058728314b75fe950261b096d224ba62954616ea91a85bec24. Checkpoint69 validado por el gate; cierre documental70 conserva DISCOVERY/BLOCKED. Continúa V321: inventario del bundle pnpm11.19.0 revela471 paquetes y brace-expansion5.0.8 afectado por GHSA-rgw5-rvv9-x895; admisión reabierta, sin promoción ni cambio global.
