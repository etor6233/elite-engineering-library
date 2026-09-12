# Composición limpia de franquicia y freshness de providers — 2026-09-04 V241

## Decisión auditada

```yaml
evidence_scope: library_and_clean_local_composition
library_structure: PASS
library_online_audit: PASS
franchise_clean_materialization: PASS
franchise_local_build_and_tests: PASS
provider_sdk_freshness: PASS_FOR_LOCKED_CLAIMS
production_project: NOT_CLAIMED
open_local_failures: 0
```

Esta evidencia no declara una franquicia concreta lista para producción. Demuestra que la biblioteca puede reconstruir su composición local vigente y que los claims enumerados abajo fueron ejecutados con revisiones fijadas. Las cuentas, credenciales, reglas empresariales/fiscales aprobadas, corpus documentales propios, infraestructura, CDN/WAF, IdP, carga, seguridad ofensiva, despliegue, rollback, recovery del proveedor y aceptación siguen siendo gates del proyecto real.

## Identidad de la composición

- plan canónico: `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`;
- selección: 61 packs únicos;
- producto materializado desde destino vacío: 659 archivos;
- runtime principal: Go 1.26.7; `GO-CONVERSATIONAL-AGENT` aporta el núcleo determinista importado por ese runtime y no crea un segundo backend;
- frontend: Node.js 24.14.1, pnpm 11.19.0, Next.js 16.3.2;
- datos: PostgreSQL 18.6, 49 migraciones `up`;
- fiscal Argentina: .NET SDK 10.0.400, `dotnet-svcutil` 8.0.0 y contratos ARCA WSAA/WSFE fijados por URL, tamaño y SHA-256.

## Freshness y procedencia de providers

### Google Ads

- distribución oficial `google-ads` 31.4.0;
- wheel oficial fijado por SHA-256 `210a7f6a7b5d40ef544090216dceeb90c9d2e7fba51c72329a690c9e5bb94474`;
- release/commit oficial fijado: `b8eb80ae277920d56fef467617795a7360c492fb`;
- lock transitivo: 24 wheels exactos para Windows CPython 3.14;
- prueba limpia: instalación `--require-hashes --no-deps`, `pip check`, 3/3 tests y probe del SDK/API v25.

Autoridades: <https://pypi.org/project/google-ads/>, <https://github.com/googleads/google-ads-python/releases/tag/31.4.0> y <https://developers.google.com/google-ads/api/docs/release-notes>.

### TikTok Business API

El snapshot fuente oficial `f809c396...` fue rechazado como runtime porque su propio grafo de imports referencia un modelo ausente. No fue parcheado ni presentado falsamente como producción. La fuente oficial recomienda la distribución PyPI de TikTok, por lo que el pack vigente usa:

- distribución `tiktok-business-api-sdk-official` 1.1.3, autor publicada `TikTok Pte. Ltd.`;
- wheel oficial fijado por SHA-256 `663b4a1f2585c4b386144d968646f695befd733c464418bbaeeac41befff33b7`;
- lock exacto de cinco wheels;
- dos entornos limpios: import, endpoint de reporting y 6/6 tests PASS.

Autoridades: <https://github.com/tiktok/tiktok-business-api-sdk> y <https://pypi.org/project/tiktok-business-api-sdk-official/>.

El SDK de reporting no se presenta como adapter de lead ingestion. Sus cuentas, autenticación, subscriptions, payload live y reconciliación de leads siguen siendo gates separados del proyecto.

### Meta Lead y demás adapters

El verificador global ahora materializa y prueba explícitamente `PYTHON_META_LEAD_RECONCILIATION_ADAPTER` bajo su lock de 18 wheels. La auditoría online final terminó con `provider_adapters=15`. Las suites focales de producto pasaron Google Ads 3/3, Google Merchant 7/7, Meta Ads 6/6, Meta Lead 10/10 y TikTok 6/6. Una consulta oficial OSV batch sobre los locks combinados de providers informó 41 paquetes y cero vulnerabilidades conocidas en esa consulta fechada; no sustituye monitoring continuo ni revalidación al actualizar.

## Reconstrucción y pruebas sobre los bytes finales

| Superficie | Evidencia ejecutada | Resultado |
|---|---|---|
| biblioteca raíz | `VERIFY_LIBRARY.ps1` | 151 packs, 1.340 archivos materializables, 44 perfiles; PASS previo a este expediente |
| auditoría online | `VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit -AllowNetwork` con CPython 3.12.14 y Go 1.26.7 | PASS; 151 packs, 121 fuentes, 15 adapters |
| MarkItDown | lock Windows CPython 3.12 | 44 distribuciones, `pip check` y tests PASS |
| backend | `go test ./...`, `go vet ./...`, `go build ./...` | PASS sobre composición limpia |
| frontend | install frozen/offline, typecheck, Vitest, build | 9 archivos/38 tests; 19 páginas Next; PASS |
| PostgreSQL | cluster 18.6 nuevo, 49 migraciones y suites | 35 SQL + 45 integración Go PostgreSQL; cero skips/fallos |
| backup/restore | distribución PostgreSQL 18.6 completa | dump custom + manifest SHA-256 `1bee6e9a...304e8`, restore separado, invariantes, cleanup y shutdown PASS |
| ARCA generator | WSDL homologación y toolchain fijados | dos clientes regenerados, build y NuGet vulnerability gate PASS |
| ARCA bridge | cliente recién generado + credenciales/bridge materializados | invoice 15/15 + parámetros 15/15; PASS |
| ARCA UDS worker | producto final + cliente recién generado | Go completo, vet/build, .NET build y 12/12 worker; PASS |
| Python dependency-free | 15 archivos de test | PASS; el caso symlink que Windows no puede crear se ejecutó 9/9 en WSL Linux |
| navegador | Microsoft Playwright 1.62.1, Chromium desktop/mobile, Firefox y WebKit | runtime 4/4 + home público final 8/8 PASS |
| calidad web | Google Lighthouse 13.4.1, Chromium rev 1234 | cinco corridas, 155 audits/run, cero warnings; mínimos performance 0.99, accessibility 1, best-practices 1, SEO 1 |

## Fallos aprendidos durante V241

Los fallos `LIB-FAIL-1899` a `LIB-FAIL-1910` quedaron registrados y cerrados mediante regresión: omisión Meta Lead del verifier, contador stale, intérprete 3.14 incompatible con lock 3.12, toolchain PostgreSQL incompleto, dependencias sobre temporales efímeros, `$LASTEXITCODE` heredado, pnpm no exacto, variable `$HOME` case-insensitive, cwd de pack incorrecto y nombre de summary Lighthouse supuesto. Ningún gate se relajó para producir PASS.

Inventario después del cierre: 1.910 IDs únicos de fallos/lecciones de biblioteca más 202 condiciones upstream, sin estados locales `OBSERVED`, `FIX_PROPOSED` o `BLOCKED` abiertos.

## Límite honesto y siguiente gate

La biblioteca queda demostrada para composición y validación local del alcance anterior. Un proyecto sólo cambia a `READY_TO_BUILD` después del intake, las decisiones A–H y el blueprint; sólo cambia a producción cuando sus proveedores y entorno reales prueban journeys completos, identidad/autorización, datos, seguridad, carga, backup/PITR/restore, despliegue/canary/rollback, soporte, capacitación y aceptación. Ese límite es parte de la solución: impide que el agente confunda código reconstruible con una empresa ya operativa.
