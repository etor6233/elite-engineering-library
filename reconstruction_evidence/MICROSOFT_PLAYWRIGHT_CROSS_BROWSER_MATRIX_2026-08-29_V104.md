# Microsoft Playwright Cross-Browser Matrix — V104

Fecha gobernante: 2026-08-29. Esta evidencia sucede únicamente la matriz Chromium 2+2 de V103; conserva sin cambios la identidad Microsoft Playwright 1.62.1, el grafo npm, la licencia/proveniencia y todos los límites de producción allí declarados.

## Revisiones oficiales fijadas

El `browsers.json` exacto del commit Microsoft `26a9e470a7b3c7822084b09fb7f13902c5f37b51`, SHA-256 `f306eed529599b1eaf2f8a85db9de2b23e1a3fe36c2b66434b7c9434fb627a99`, declara:

- Chromium `151.0.7922.34`, revision `1234`;
- Firefox `153.0`, revision `1538`;
- WebKit `26.5`, revision `2336`.

Los binarios fueron obtenidos por el instalador oficial Playwright; no se copian dentro de la biblioteca. Firefox descargó 119,9 MiB y WebKit 59,6 MiB desde `cdn.playwright.dev`; la disponibilidad de red/disco sigue siendo condición de la máquina que no tenga esas revisiones cacheadas.

## Ejecución

- runtime aislado: Chromium desktop PASS, Chromium mobile PASS, Firefox desktop PASS, WebKit desktop PASS = 4/4;
- perfil web: 4 packs/53 archivos; frozen offline install y Next 16.3.2 production build PASS;
- home real: las cuatro variantes anteriores PASS = 4/4;
- assertions preservadas: status 200, `main`, H1 único/no vacío, catálogo, locale, headers de seguridad, sin overflow, page errors ni console errors.

WebKit reveló un fallo real: la directiva CSP productiva `upgrade-insecure-requests` promovió siete chunks `http://127.0.0.1:4173/_next/...` a HTTPS y produjo `SSL connect error`. No se retiró la directiva, no se habilitó `bypassCSP` y no se ignoró la consola. El test WebKit intercepta exclusivamente `https://127.0.0.1:4173/`, solicita el mismo path al mismo loopback HTTP y entrega esa respuesta al request promovido; el response principal conserva el CSP y el test lo comprueba. Este shim sólo autoriza el smoke local; producción sigue exigiendo HTTPS/edge reales.

## Código y límites

El pack `MICROSOFT_PLAYWRIGHT_BROWSER_GATE.md` 0.1.1 conserva nueve archivos y la misma proveniencia V103: ocho `AUTHORED`, una licencia `ADAPTED` sólo CRLF→LF, cero source Microsoft falsamente atribuido. Los cambios V104 son configuración/test/lock/verifier locales que llaman las APIs oficiales Microsoft y fijan las tres revisiones.

`LIB-FAIL-1128`–`LIB-FAIL-1135` conservan el probe booleano, dos patches rechazados, el fallo WebKit, el cwd incorrecto, el bypass inefectivo y los dos fallos de cleanup. Ninguno se ocultó. La matriz no prueba roles autenticados, issuer/edge/backend real, tecnología asistiva, performance/load, seguridad ofensiva, deploy, canary, rollback ni producción.
