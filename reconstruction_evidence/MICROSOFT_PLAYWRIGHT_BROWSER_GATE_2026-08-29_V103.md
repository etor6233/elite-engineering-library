# Microsoft Playwright Browser Gate — V103

Fecha gobernante: 2026-08-29. Alcance: admisión, materialización y ejecución local reproducible del gate de navegador para el perfil web Elite. Esta evidencia no promueve el perfil ni el proyecto a producción.

## Autoridad oficial fijada

- Repositorio: `microsoft/playwright`; release `v1.62.1`; commit firmado/verificado `26a9e470a7b3c7822084b09fb7f13902c5f37b51`.
- Archive GitHub: 42.927.948 bytes, SHA-256 `29377110a042e60e8c8192eae145e8e8554b58b31073bdcdf9e25b47d019b186`.
- npm `@playwright/test@1.62.1`: 8.750 bytes, SHA-256 `009534220efd98c0361d8c4ee7e3db1ed510ab88a23e98b5081ef1c8fed64965`, integrity `sha512-DTcUc8qii+cpHvtOwggMtBRMjKZHXYWdw8syRYu2vtzuq4Wxphqq4NfCs5Zt44L6mA8rfDfj+PHnxFc/FeK6mQ==`.
- Manifest oficial `packages/playwright-test/package.json`: 754 bytes, SHA-256 `be51f784a5742a1aa075e0173f5f04f189a791cc30c40e3bb18be584f696764e`.
- `packages/playwright-core/browsers.json`: 1.780 bytes, SHA-256 `f306eed529599b1eaf2f8a85db9de2b23e1a3fe36c2b66434b7c9434fb627a99`; Chromium revision `1234`, versión `151.0.7922.34`.
- LICENSE oficial original: 11.601 bytes/SHA-256 `45873d00a0dd243596deb4aa23b2493b3d1f0671921bf2538ea431d7380220eb`. Copia empaquetada declarada `ADAPTED` sólo por CRLF→LF: 11.399 bytes/SHA-256 `7fab1461b41970ff376f1c9303a637076bfaaeb71cd12dd3a1c44aaf59a1a2b9`.
- Grafo npm fijado: cuatro paquetes. `pnpm audit` y consulta batch oficial OSV fechada: cero findings conocidos. Esto no garantiza ausencia futura de vulnerabilidades.

Fuentes oficiales: `https://github.com/microsoft/playwright/releases/tag/v1.62.1`, `https://github.com/microsoft/playwright/tree/26a9e470a7b3c7822084b09fb7f13902c5f37b51`, `https://playwright.dev/docs/best-practices`, `https://playwright.dev/docs/ci`.

## Código materializable y proveniencia

`MICROSOFT_PLAYWRIGHT_BROWSER_GATE.md` 0.1.0 materializa nueve archivos. Ocho son `AUTHORED` por Elite: manifest/lock de instalación, source lock, configuración, dos tests, verifier y README. Uno es la licencia Microsoft `ADAPTED` exclusivamente por normalización de fin de línea declarada. El ejecutable Playwright no se copia ni se atribuye falsamente: `pnpm install --ignore-workspace --frozen-lockfile` obtiene el paquete oficial exacto definido por el lock.

El pack canónico final tiene 33.733 bytes/SHA-256 `6ce8edecdf1c47b6f168496979a6e931f089c05840bdf4debf7f851ed23c85cd`. Roundtrip Markdown: 9/9 archivos y hashes PASS. Frozen offline install: tres paquetes activos en Windows más la entrada condicional de plataforma del lock; cero descargas desde el store auditado.

## Pruebas tangibles

Entorno: Windows x86-64, Node 24.14.1, pnpm 11.19.0, Chromium oficial revision 1234.

- runtime real Microsoft Playwright: desktop Chromium PASS + mobile emulation Chromium PASS = 2/2;
- perfil `ENTERPRISE_WEB_PACK_PLAN.md`: 4 packs/53 archivos, sin colisiones;
- web: frozen offline install PASS, TypeScript/typegen PASS, 14/14 tests PASS, Next 16.3.2 production build PASS;
- target browser: home público desktop PASS + móvil PASS = 2/2;
- assertions target: `main`, H1 único, enlace catálogo, `lang`, CSP, `X-Content-Type-Options`, `Referrer-Policy`, `Permissions-Policy`, ausencia de overflow horizontal, errores de página y console errors;
- `VERIFY_LIBRARY_PASS`: 72 packs/692 archivos, perfil web 53;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 72 packs, 112 fuentes y `microsoft_playwright_browser_gates=1`.

El plan web final tiene SHA-256 `c120e0c344655ba8aa765bf6653c9117245b5697e61b2bfd4c5f4d552d9164be`; su record registra exactamente Playwright 0.1.0 junto a BFF, OIDC y license gate.

## Fallos convertidos en memoria

`LIB-FAIL-1106`–`LIB-FAIL-1127` conservan paths inferidos, SCA opcional encadenada, layout pnpm, cwd, webServer, Next CLI, conteo de receipt, proveniencia de licencia, arrays/update, contextos largos, path del compositor, aislamiento de workspace, reincidencia `foreach` y cleanup read-only observados. Todos están cerrados por regresión; no se ocultó ningún fallo para obtener el PASS.

## Límites que permanecen condicionados

La evidencia no demuestra roles autenticados ni IdP real, efectos backend/PostgreSQL, Firefox/WebKit, interoperabilidad con tecnología asistiva, métricas de rendimiento, seguridad ofensiva, despliegue edge/cloud, canary/rollback ni producción. Esos gates necesitan las decisiones, accesos y entornos reales del proyecto. El smoke local no se puede convertir honestamente en una afirmación universal.
