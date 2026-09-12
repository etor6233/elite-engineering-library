# V401 — contención de ZIP en pnpm y bloqueo de su release afectada

2026-09-11. Mantenimiento local, reparación AUTHORED; controles45/48.
Se cierra la implementación y comprobación de la contención en candidato aislado.
La admisión general, fuente/licencias/relinking y release final siguen abiertas.
Daybreak/V386 permanece diferido por indicación del usuario.

## Hallazgo actual y candidatos descartados

OSV2.5.1 fijado encontró GHSA-vwc7-r8mq-g2x9 / CVE-2026-76845 en adm-zip0.6.0
del pnpm11.25.0 original.455archivos originales siguen intactos. Se añadieron
al inventario tres inputs Noble ligados al source lock de OpenPGP:476identidades
conocidas. No es el cierre completo del grafo recursivo nativo. El473/0 histórico
se preserva como historia; no autoriza hoy ese mismo artefacto.

Autoridad: https://github.com/advisories/GHSA-vwc7-r8mq-g2x9 . Rango afectado
>=0.5.9 <=0.6.0, sin versión corregida publicada según el aviso consultado.
La PR https://github.com/cthackers/adm-zip/pull/575 está open/merged=false en el
receipt oficial; merge_commit_sha de una PR abierta no prueba integración.
pnpm11.26.0 conserva adm-zip0.6.0: rechazado como solución sin ejecutarlo.
pnpm12.4.1 cambia a ejecutables opcionales nativos:22identidades visibles no son
un SBOM del binario. Tarballs npm íntegros/SHA512 comprobados, firmas pendientes,
sin ejecución de esos dos candidatos ni afirmación de seguridad por ausencia de
un encabezado. No se ejecutaron exploits, PoC destructivas ni la parte diferida.

## Reparación canónica y límites

PNPM-ARTIFACT-SELECTION-GATE0.8.0 incorpora contain_zip.py,9archivos en total.
Acepta únicamente455archivos originales con bytes/hashes exactos. Publica un
candidato separado con442archivos de payload:441idénticos y un bundle cambiado.
Quita físicamente16módulos adm-zip y bloquea ZIP antes de crear su directorio o
descargar. Las dos entradas auxiliares también rechazan de inmediato. Mantiene
los notices originales, registra ADAPTATION.txt, diff y receipt determinista.
No cambia ni suplanta versión/firma upstream. Identidad local por SHA/receipt.

Original bundle: ddc64218bc85fb88d28b5def06eb01fafb39cb67f7a9465ce736456658cc7f11.
Candidato: aff07c57e598640266b0ca1c795879db1bbdd2eafe71e96be34026358ba1c42b.
No es una reparación oficial de adm-zip: elimina esa capacidad del candidato.
La descarga de runtimes binarios ZIP queda indisponible; usar un runtime admitido
por separado. Los planes anteriores siguen fijando la proyección original y
no pueden apuntarse al bundle cambiado. Ni source/relinking ni redistribución
global se consideran aprobados. El recibo del materializador no ejecuta nada;
las ejecuciones locales de calificación se documentan separadamente debajo.

VERIFY_EXECUTABLE_LIBRARY deja de ejecutar pnpm desde PATH para inferir admisión
por su versión. Reporta pnpm BLOCKED y la condición exacta; Full se detiene antes
de instalar con un runtime sin admisión. Un número de versión mayor no lo abre.
Reapertura: digest exacto con cierre de sus condiciones, o release oficial fija
admitida después de los gates correspondientes. No se acepta riesgo por defecto.

## Evidencia nueva y reutilización

- Fuente reconstruida y perfil1/9 exactos; siete archivos previos sin cambios,
  README actualizado y un archivo AUTHORED nuevo.90tests históricos de esas
  fuentes previas no se vuelven a ejecutar ni se presentan como pruebas nuevas.
- Nueve comprobaciones nuevas: input incorrecto/mutado, tres alteraciones de
  salida, archivo extra, destino ocupado conservado, sintaxis Node y frontera
  ejecutada desde las funciones del bundle real. Cuatro rechazosZIP con cero
  llamadas de descarga/tempDir; despacho TAR a fixture conservado.
- OSV2.5.1:475identidades conocidas restantes, cero hallazgos. El inventario es
  conservador, mantiene identidades de otras exclusiones y no prueba SBOM
  recursivo completo. adm-zip sólo se quita tras eliminación física comprobada,
  sin ignore de advisory ni versión ficticia. El hallazgo original permanece.
- Tres instalaciones reales offline/frozen/ignore-scripts/ignore-pnpmfile,
  entorno reemplazado y config vacía, copy/verify-store-integrity: web, Playwright
  y Lighthouse PASS; cero descargas, seis manifests/locks idénticos. Sus tres
  comandos --version funcionan. La primera invocación exec de Next omitió
  configuración store/offline: pnpm reinstaló y descargó71paquetes. Su log se
  conserva como FAIL807; no fue offline. La invocación corregida con config
  explícita termina Already up to date/Next16.3.4 sin descargas. No se repiten
  las suites de negocio/browsers.
- La comparación directa terminó con12571/12571archivos idénticos a la
  referencia, cero diferencias/faltantes. El primer lector falló por Windows260;
  el corregido usa rutas extendidas. Terminó antes de la solicitud de detener la
  relectura opcional: no hubo proceso detenido. El resumen preliminar incompleto
  queda preservado en final-correction-before; aquí rige el resultado completado.
- Función real Get-ToolchainReport comprobada con stubs: cero ejecuciones de
  pnpm rechazado, siete probes independientes y Full bloqueado antes de instalar.
- Inventario con parser canónico:167packs/1588archivos,1331AUTHORED,
  150ADAPTED/107VERBATIM. El perfil de producto70/820 no cambia.
  No se repiten los164pasos generales de Preflight por este cambio acotado.

## Continuidad

H tiene su prompt generado y referencia exacta; BLOCKED conserva la frontera
entre pruebas locales y aceptación de release/producción. No se inventan SLO,
RPO/RTO, regiones ni respuestas del usuario. TEST02 sigue con integración y
claims de los owners restantes; TEST03 con admisión completa del tooling/grafo;
TEST07 depende del artefacto y gates finales. Nada permite declarar48/48.
No repetir esta contención, las tres instalaciones, ni entregas de notices sin
delta. Conservar el hallazgo original y la condición del candidato local.

## Artefactos locales verificables

Stage: C:/Users/NL/AppData/Local/Temp/elite-v401-20260911-closure.
La implementación queda en el pack canónico; los recibos locales no equivalen
a una release portable firmada. Ningún provider, dato real o despliegue afectado.

| Archivo del stage | SHA256 |
|---|---|
| nested-scan-result.json | a0fab5b543c529421243e548bbbb651b87d939456d07667c61e4e1a8e2944422 |
| security-authorities.json | 9591cfd27e622ba301fef0fb1c5d249eca70c59bb3d4c7035cbfbb13b82e1ccb |
| candidate-acquisition-receipts.json | b36931d0d3284d48bafd5189215334a7406898c1214d43e977761c434c33b572 |
| containment-tests.json | 546471e7fbe9a20df3dd3d22347b69eabeaa5e15e36b9573b4c043cb569ec2b0 |
| contained-scan-result.json | 9b6cf002685bde579f254f9b2110d80dc380ceeea7c799142d6e8d8256d90deb |
| contained-consumers-result.json | 500ed8c4b2460b4e73777b720c0c636f0b7927ce9afc1a9ca664ef41b9ecadca |
| canonical-delta.json | c849273ea1f3c7040862d1321d03be16ba6593ec64f3322930d3ed949e5df734 |
| preflight-admission-result.json | 0d3dcc6acffc5ef9b1cdcf6fe7926f3b25b080dd50a031b022464b5e785dec67 |
| zip-disabled/containment-receipt.json | df2f0a626448bace5d306978bcd669eb8307de1d35bf9c8ac71f3322ab638fec |
| zip-disabled/adaptation.diff | 24940c7b258ebf8640939a04151023cb215b484698cb06e540322447574a271e |
| installed-payload-comparison.json | 35ef7f987bfccd7dc3d19c29c635d707935c666479b67eed82a300ab0f1d5a4e |
| dispatch-correction.json | 6717134d980a8d74c84369c9b67191d3354ff0d6781da3b4ed56da32ee261801 |

## Cierre272

Assurance enlaza EVID-107; readiness recalculado conserva39observaciones/BLOCKED.
Regresión afectada test_toolchain_resolution.ps1:46checks PASS, sin auditoría general.
Checkpoint272 y contrato validados. No controles ni criterios promovidos.
