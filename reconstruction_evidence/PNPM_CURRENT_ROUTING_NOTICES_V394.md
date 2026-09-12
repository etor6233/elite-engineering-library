# V394 — pnpm: rutas actuales y aviso BlueOak verificable

2026-09-11. Mantenimiento de biblioteca. Entrada: checkpoint251 validado;45/48.
Resultado: resuelto FAIL782 en el preparador de recetas;3perfiles actuales
preparan y verifican correctamente. Las recetas entregan localmente un aviso
BlueOak verificable para5paquetes.51tests PASS y10mutaciones reales rechazadas.
El bloque pnpm completo y TEST02/03/07 permanecen abiertos.

## Problema reproducido y cambio

El código0.2.0 conservaba hashes anteriores: enterprise-web y Playwright
rechazaban sus package.json actuales. Lighthouse pasaba. El baseline conserva
los hashes anteriores y vigentes; no se debilitó la comparación.0.3.0 fija
BFF0.5.16 y Playwright0.1.37; lock Playwright y ambos inputs Lighthouse intactos.
La nueva comprobación Preflight compone ENTERPRISE_WEB_PACK_PLAN vigente y
compara sus tres perfiles con el preparador. Evita que futuros cambios vuelvan
a quedar ocultos por fixtures sintéticos. El pack conserva6archivos; cambian
README, plan_install.py y test_routing.py. Selección del artefacto y política
de sus455miembros intactas. El preparador no incorpora dependencias ni ejecuta pnpm.

## Aviso de licencia entregado

Los cinco package.json seleccionados son byte-idénticos a las declaraciones
del autor en los commits fijados en V339. Nombres/versiones:
chownr3.0.0, isexe4.0.0, minipass7.1.3, tar7.5.22 y yallist5.0.0.
Todos declaran BlueOak-1.0.0. La autoridad oficial permite entregar el texto
o el enlace de la licencia; se elige su enlace oficial:
https://blueoakcouncil.org/license/1.0.0 (consultado2026-09-11).

Cada receta v2 incluye BLUEOAK-NOTICE.md con las cinco identidades, hashes de
manifiestos, declaraciones del autor fijadas por commit y enlace de licencia.
El recibo liga su SHA-256; verify recalcula el contenido desde la política
propia y rechaza su ausencia, edición o rehash conjunto del aviso/recibo.
El aviso es un suplemento local declarado; no se inventa un LICENSE original
para chownr ni se elimina la salvedad yallist sobre paquetes anidados.
El payload conserva sus442archivos y22notices byte-idénticos. La receta no
redistribuye ese payload: su futura entrega debe acompañar el aviso. El hueco
de generación/entrega local del aviso queda resuelto; entrega downstream,
otras licencias, publicación y fuente completa siguen pendientes.

## Evidencia ejecutada

-51tests en candidato y nuevamente desde Markdown reconstruido:17selección y34routing.
-Seis fuentes reconstruidas idénticas al candidato, con UTF-8/LF.
-Composición canónica actual:6packs/126archivos; sus3perfiles preparan y verifican con
  proyección, recibo de adquisición y Node exactos, sin mocks ni ejecución pnpm.
-Siete mutaciones reales de manifiesto/lock/workspace y tres del aviso rechazadas.
  Después de restaurar cada input, verify vuelve a pasar.
-Pruebas sintéticas adicionales: borrado de aviso/manifiesto, falsificación conjunta
  del recibo, declaración incorrecta y fallo de escritura antes de publicar.
-442archivos de payload y receipt original intactos;5declaraciones de autor exactas.

Las pruebas de races/configuración/recibos/Node/junction previas permanecen.
El cambio no toca frontend, backend, lock de producto ni los48oráculos.
Reusar evidencia runtime previa sólo por identidad del producto; no se afirma
un nuevo build de producto ni scan OSV por las recetas. El Preflight sí realiza las instalaciones offline habituales de sus gates de Playwright/Lighthouse.

## Límites y continuidad

45/48: TEST02 integración integral, TEST03 tooling/nativos/seguridad restantes,
TEST07release. pnpm aún requiere QRCode MIT, nested/source obligations, semver-utils
y correspondencia/publicación completa según V338–341/V364. FAIL575 sigue abierto.
V386/libxml2 permanece diferido por la restricción ya registrada del servicio;
no se repitió ni rodeó aquella operación. runtime_admitted=false,
redistribution_admitted=false y executed=false en las3recetas.
Rollback: restaurar pack/plan/verifier previos desde before; preservar nuevas recetas
como evidencia, sin sustituir directorios existentes ni convertir un v1 en v2.

## Receipts

Raw evidencia local: Temp/elite-v394-a9f68fc2b82c4dc79d0ee9b843936b78,
resolver elite-v394-current.txt. Los datos locales no se incluyen en la distribución.

| Archivo relativo al stage | SHA-256 |
|---|---|
| baseline-routing.json | d9cc0664f1745cdb784326806a2587ec928df38500c56d13c9cfb93f3a691f5c |
| current-consumer-inputs.json | 5afa48787a42b4da864f22538950db59e0288d9db176fadcbed068d99f330c20 |
| blueoak-bindings.json | d9518f7cf6939fdbe969baf07748e0b9821e8aefc9b432f3760c7cf5de9f5ee3 |
| source-parity.json | ed0e5516a01195aace95d33f0d2f3eeb7e02fb364c02c8c06d886977197c7441 |
| rebuilt-tests.log | 1a836187823fdcff6c3fd0462899ca6539686c79d8ab9a5c4cde64fac0da1aae |
| current-web-composition.log | 50502449356ba218c750f5b3d2fb93379d236419b12f837ac865617c899a51c4 |
| real-routing-results.json | ec7b0929c1da7d82fb4eeacdaa5c7f591281cdcecc7d062d2a0fb700d27f338f |
| projection-before-after.json | 4b56bb3ee5a5bfa16b7ac0832b67866ff5d5348e7357be81e1e3d200eca9f27b |
| author-declaration-parity.json | a2cb930fc94b9abea56de946af5d1de61a91be09114f286e35799375bda0e1a6 |

Preflight252 completo:164pasos ejecutados PASS,56perfiles; inventario165packs/1561archivos/842Markdown. Incluye las34pruebas routing y la nueva comprobación de3perfiles canónicos actuales. Disponibilidad general BLOCKED únicamente por Docker ausente; no es una aprobación productiva.

## Cierre de verificación general

Preflight252 conserva allow_network=false; ejecuta pnpm11.25.0 en los gates offline habituales de Playwright y Lighthouse. Las3recetas nuevas sólo se prepararon/verificaron, sin ejecutarlas. La comparación canónica de los3perfiles pasa. El código incorporado permanece idéntico al rebuild probado; los48tests de aceptación conservan acciones, oráculos y estados. FAIL782 reparado; obligaciones pnpm restantes y TEST02/03/07 siguen abiertos. EVID-100 enlaza este alcance acotado.

| Recibo final | SHA-256 |
|---|---|
| preflight252.json | 01fa61d3b9e90d0cc2298841761e85678c4ed66df988f62121effabc8f653866 |
| preflight252-process.json | 29a60dd4656abf0d8c8e10628a217bce5cc803a90dd2a7681a01bef3925ceea8 |
| preflight252.log | d1763c72ee6d115ac9c15759012490f62094002454561a22621f6c3df79ed89d |
| acceptance-preservation.json | 60f07e2cd5fa52a985158ca705afa9af0cddd79855536e6a2e9d8aaa2fe2c06d |

Rectificación V395/FAIL786: el texto original extendía incorrectamente la no-ejecución del preparador al Preflight general. Los logs252 conservados muestran3paquetes Playwright y119Lighthouse instalados offline con pnpm11.25.0. Original antes de rectificación y prueba de alcance preservados en stage395; las3recetas no se ejecutaron, ni cambian sus flags o la admisión. Ver QRCODE_VENDOR_NOTICE_DELIVERY_V395.md.
