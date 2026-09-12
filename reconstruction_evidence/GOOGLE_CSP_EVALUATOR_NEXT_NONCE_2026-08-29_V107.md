# Google CSP Evaluator + Next.js nonce — V107

Fecha de corte: 2026-08-29. Esta evidencia gobierna únicamente el refuerzo CSP del perfil web V107; no convierte la biblioteca ni un proyecto derivado en producción aprobada.

## Autoridades públicas y decisión

- Google CSP Evaluator: <https://github.com/google/csp-evaluator>, licencia Apache-2.0, release/tag `v1.1.8`, commit exacto `2e00e419d0be0c204a91ce08f7cc4cb838a2b75d`, tree `1231fb555f2ea7063a9c5fbdf333a5bc9dd8514e`.
- Guía CSP oficial de Next.js/Vercel: <https://nextjs.org/docs/app/guides/content-security-policy>. El patrón seleccionado genera nonce por request en `proxy.ts`, fija CSP tanto en request como response y fuerza render dinámico mediante `connection()`.
- Advisory Next.js nonce XSS: <https://github.com/vercel/next.js/security/advisories/GHSA-ffhc-5mcf-pf4q>. El perfil usa Next.js `16.3.2`, posterior a la corrección `16.2.5` indicada por el advisory.
- Microsoft `@microsoft/eslint-plugin-sdl` 1.1.0 fue evaluado y rechazado: su única rama peer es ESLint `^9`, ya EOL en la fecha de corte. No se copió, distribuyó ni presentó como gate vigente.

Google CSP Evaluator se admite sólo como dependencia de desarrollo condicionada. Su propio README dice que no es un producto oficial de Google, no ofrece garantía y su lista de bypasses no es exhaustiva. Un resultado limpio no demuestra ausencia universal de XSS ni reemplaza pruebas ofensivas.

## Identidad exacta adquirida

| artefacto | identidad |
|---|---|
| source ZIP por commit | 61.437 bytes; SHA-256 `60065e01be2332ffc902bf71821baeb08d8d29c10cd8d2598fe9fb4d3e32ab4d` |
| source `LICENSE` | SHA-256 `58d1e17ffe5109a7ae296caafcadfdbe6a7d176f0bc4ab01e12a689b0499d8bd` |
| source `package.json` | 761 bytes; SHA-256 `76ede5870a6a38313483a134ea388d07d9caa50831ec0305cb0697dcc9601874` |
| npm `csp_evaluator@1.1.8` | 79.066 bytes; SHA-256 `2b4f2a30ed67a39908708c05ba559accbad0b5a6c0fa29f1761a641a342fd904` |
| npm integrity | `sha512-EwOnfYuNbTytvbMKsLixTrRgnjOa0WZCxGy8A9nnSYAicrdwn+T/epU/yjgymmOxlgKnvH+8wXt+7p/8ak5Feg==` |
| npm shasum | `7bd461ec69a66cb067ff31ff6aacea101e88a46b` |

El tag es lightweight y el commit figura unsigned. El paquete npm declara cero dependencias runtime, contiene 120 archivos, una signature de registry y cero attestations. La licencia del source y la del tarball coinciden. Estos límites permanecen en el source lock y no se ocultan por la marca del repositorio.

## Código materializable incorporado

`TYPESCRIPT_GO_API_WEB_BRIDGE` 0.3.0 materializa 36 archivos. V107 agrega:

- `src/platform/security/csp-policy.ts`, `AUTHORED`: política estricta con `default-src 'self'`, script nonce + `strict-dynamic`, `unsafe-eval` sólo en desarrollo, `object-src 'none'`, `base-uri`/`form-action` self y `frame-ancestors 'none'`;
- `src/platform/security/csp-policy.test.ts`, `AUTHORED`: cuatro regresiones ejecutadas con el paquete Google exacto;
- `src/proxy.ts`, `ADAPTED`: adaptación declarada de la guía oficial Next.js; genera UUID aleatorio, lo codifica base64, sobreescribe headers entrantes no confiables y fija CSP en request/response;
- `layout.tsx` fuerza render dinámico con `await connection()`; la CSP estática incompleta anterior fue retirada de `next.config.ts`.

`MICROSOFT_PLAYWRIGHT_BROWSER_GATE` 0.1.2 materializa 10 archivos y añade `tests/nonce-csp.spec.mjs`, `AUTHORED`. La prueba exige nonce base64 de al menos 16 bytes, ausencia de `unsafe-inline`/`unsafe-eval` en producción, `strict-dynamic`, igualdad entre header y propiedad DOM `.nonce` de todos los scripts, y rotación después de reload.

Hashes canónicos al cierre previo a esta evidencia:

- bridge: `3bf501eea3d30107b392a44df160e60a84c177a5406350bab3673149ae5c4d87`;
- Playwright gate: `7c9a6aada592f382596db713720f91b7c7f24bdffe689bcdec897cd5a7855222`;
- acquisition core 0.4.69: `fd97b12d68e6a04d170a5dfe2411b29b14d5f196d38123a9bfb3b59e83324d0a`;
- perfil web: `39485883ae72e332d37523a257fc9b8353ef83f360d872fe0e1d6bd59ec08ffd`.

## Resultados ejecutados desde Markdown

La política anterior sin `script-src` produjo un finding HIGH esperado. La política nonce V107 produjo sólo avisos informativos de entradas ignoradas y cero findings accionables.

El perfil limpio `ENTERPRISE_WEB_PACK_PLAN.md` materializó 5 packs/69 archivos y pasó:

- `pnpm install --frozen-lockfile` offline;
- TypeScript/Next type generation;
- 7 archivos Vitest / 21 tests;
- build productivo Next.js 16.3.2 con todas las rutas dinámicas;
- Microsoft Playwright runtime 4/4;
- Microsoft Playwright home+nonce/CSP 8/8 en Chromium desktop, Chromium mobile, Firefox y WebKit;
- Google Lighthouse 13.4.1, cinco runs: performance mediana `0.99`, accessibility `1.00`, best-practices `1.00`, SEO `1.00`;
- auditoría ejecutable global: `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=73 upstream_sources=115`;
- Foundation con Go oficial 1.26.7 verificado: `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Foundation backend=go-test web=typecheck-test-build browser=playwright-chromium-firefox-webkit quality=lighthouse-five-run`.

## Condiciones que permanecen abiertas

El nonce fuerza render dinámico y reduce el uso de caché estática; la guía oficial lo advierte. Los tests locales no prueban comportamiento del CDN, WAF, reverse proxy, hosting standalone ni headers mutados por el edge. Antes de producción el proyecto debe ejecutar sobre su entorno real: nonce distinto por response, igualdad header/scripts, cache isolation, report-only/rollout, rollback y observabilidad sin datos sensibles.

Tampoco están demostrados por este pack Trusted Types, ausencia de todos los sinks directos/raw HTML, seguridad ofensiva, IdP/roles/backend reales, accesibilidad con tecnología asistiva, RUM/carga ni aceptación de negocio. Esas condiciones no se convierten en PASS por inferencia.
