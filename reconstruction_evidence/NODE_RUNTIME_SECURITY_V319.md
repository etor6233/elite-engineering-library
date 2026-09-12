# Identidad, advisories y selección del runtime Node — V319

Mantenimiento T2803/T2809, 2026-09-08. Candidato local Node24.20.0 Windows x64
para el BFF/portales/OIDC y gates Node existentes. No código de producto nuevo,
instalación global, cambio de engine/lock de aplicación, provider live, cuenta,
gasto, scheduler, ZIP de release ni promoción. El conjunto sigue BLOCKED.

## Hallazgo y fuentes oficiales

El Node instalado24.14.1 tiene SHA256
58e74bf02fc5bbacc41dcb8bef089961cd5bddd37830b87784e4fc624d145d1f.
La revisión comenzó por los avisos oficiales de junio y julio2026, posteriores
a esa versión. Reachability del target permanece UNKNOWN. V317 escaneó Go;
su scan de módulos/npm no certificaba este runtime ni todos sus componentes.

Fuentes preservadas: [seguridad junio](https://nodejs.org/en/blog/vulnerability/june-2026-security-releases),
[seguridad julio](https://nodejs.org/en/blog/vulnerability/july-2026-security-releases),
[release24.20.0](https://nodejs.org/en/blog/release/v24.20.0),
[checker oficial](https://github.com/nodejs/is-my-node-vulnerable/tree/c37a56bad56e34fe5223ddd3cb223cc4158136ae),
[datos security-wg](https://github.com/nodejs/security-wg/tree/21bf6214b10d3a70250828c05b5e7026f4e27f50/vuln/core).

Node24.20.0: commit71b8b174857e25106d39b61a9e6f30d927da8b01;
ZIP SHA2566cac9ffbca8f6a47091e4b5c772e0606049c3871cb67d900c0cedde630e545ba.
1994 archivos extraídos con rutas validadas; licencia y hashes conservados.
SHASUMS autenticado con gpgv2.4.9 y fingerprint
5BE8A3F6C8A5C01D106C0AD820B1A390B168D356, presente en README oficial fijado.
Keyring oficial nodejs/release-keys5b7f55f4a7e35d1176d27a6b81b0c3c3b794216b.
El primer gpgv de Git falló por ruta Windows: su plaintext queda no autenticado.
El segundo, cwd/rutas relativas y output distinto, devuelve exit0/VALIDSIG y
plaintext byte-idéntico al checksum oficial. No se omitió la firma.

## Tres alcances de seguridad distintos

1. **Node core:** OSV npm/node devuelve0 incluso para el control vulnerable;
   se rechaza la equivalencia. El motor oficial Node1.6.1, con sólo2 URLs fijadas
   y semver7.8.5, evalúa194 avisos y el schedule oficial. Cada fragmento coincide
   con su índice, Git blob SHA y SHA256. Consulta acotada30s y hashes comprobados:
   Node24.14.1/win32 →22CVEs;24.20.0/win32 →0. Esto no significa SCA nativo de OSV
   ni ausencia de vulnerabilidades desconocidas/componentes no cubiertos.
2. **ZIP completo:** OSV2.5.1 con inventario completo145 paquetes npm da4 afectados
   y9 advisories. brace-expansion5.0.7, ip-address10.2.0, tar7.5.19 y undici6.27.0.
   Se conservan todos los IDs, manifests, JSON y9 avisos oficiales. El ZIP firmado
   queda REJECTED para distribución completa; no se edita ni se silencian matches.
3. **Selección mínima:** se descarga el node.exe independiente oficial y se prueba
   byte-idéntico al del ZIP,93381448 bytes,SHA256
   5c976096e04e5c2c1f091938926234cc9fbebfe9787ddd149351b3b0ecc707b5.
   Se acompaña con LICENSE íntegra SHA256
   ed34dd8e3f0a78dbaf00d0444ce8e285b015b765379c2e17880455f70370f8e9.
   Sólo estos2 archivos, sin npm. Es selección AUTHORED de bytes oficiales
   VERBATIM, no una nueva release upstream ni prueba de que npm quedó corregido.

Los avisos npm remiten a los owners oficiales de
[brace-expansion](https://github.com/juliangruber/brace-expansion/security/advisories),
[ip-address](https://github.com/beaugunderson/ip-address/security/advisories),
[node-tar](https://github.com/isaacs/node-tar/security/advisories/GHSA-r292-9mhp-454m)
y [undici](https://github.com/nodejs/undici/security/advisories).
Los datos del probe Node y su grafo CLI2 paquetes/0matches no cubren el entrypoint
Actions que se omitió, ni convierten la firma del runtime en un scan de malware.

## Funcionalidad y reconstrucción

Dos composiciones canónicas67/746, primero con Node del ZIP y luego con el
ejecutable independiente: frozen/offline install,115 tests web PASS/1SKIP
preexistente, build Next16.3.2,92 fases Playwright1.62.1 y4 agenda PostgreSQL18.6
PASS por composición. Las92 fases son23 por Chromium desktop/mobile, Firefox y
WebKit, verificadas por bloque de proyecto, sin FAIL/SKIP del gate conectado.
Go1.26.7 se seleccionó desde official-toolchain, no current-toolchain1.26.8.
TLS negativo, autorización por rol/tenant, firma/idempotencia y recuperación
conservan sus oráculos; no se relajó TLS ni se agregó retry de POST.

Los746 archivos del producto son idénticos entre ambos árboles probados.
745 conservan el SHA materializado; next-env.d.ts es regenerado por Next e
idéntico entre builds. No hay cambio canónico de implementación que sincronizar.
El proyecto fija pnpm11.19.0: se observó el entrypoint real bajo la selección
automática y luego se llamó directamente con Node exacto, evitando otro salto.
Los891 archivos del pnpm observado coinciden con su tarball npm11.19.0, sin extras,
verificado por integrity SHA512; esto prueba identidad, no SCA de su bundle compilado.
El stdout de pnpm incluye progreso; la identidad del hijo se lee de JSON y hash,
no de una línea supuesta. Los logs iniciales fallidos se preservan.

## Rendimiento y retorno local

Ocho arranques alternados baseline/candidato del mismo build Next,256 GET locales
a icon.svg con cuerpo SHA idéntico y CSP nonce fresco. Cuatro arranques/128 GET
por versión. El último par prueba retorno local al ejecutable anterior y vuelta
al candidato; no prueba migraciones, backup/restore ni rollback productivo.
Los límites se fijaron antes: startup mediano y p95 ≤3×, RSS muestreada ≤2×;
todos pasan. La versión antigua sólo se retiene para evidencia, nunca como
rollback productivo seguro. No carga externa, datos de negocio ni provider.

- baseline: startup mediano0.552s; p50/p95/p99 15.16/26.80/28.36ms; RSS muestreada máxima129347584 bytes.
- candidate: startup mediano0.555s; p50/p95/p99 14.78/25.86/29.02ms; RSS muestreada máxima114253824 bytes.

Muestra pequeña y ruta local limitada: números descriptivos, no SLO de negocio,
throughput de órdenes ni garantía de costo/producción.

## Gap del checker y continuación

Seis probes con fixtures locales y el mismo motor oficial: baseline1/candidato0,
plataforma inválida1 y JSON malformado1; corpus vacío{} devuelve0 incorrecto para
baseline vulnerable; fuente inmóvil requiere kill externo3s. Fuente y logs before
se conservan. No se escribió un matcher nuevo ni se alteró el corpus original.
El gate CAPABILITY_GAP_RESOLUTION valida USE_CONDITIONED_PACK/ready=false:
G0–G3 PASS; G4–G8 no admitidos para CLI autónomo, lifecycle operativo o pack portable.
La consulta puntual con hashes exactos es evidencia acotada; no autoriza reutilizar
el CLI sin controles. Reabrir con adapter mínimo que preserve motor oficial,
integridad/freshness, deadlines/límites, negativos y reconstrucción, o fix upstream.

Autoridades, lock, mapa de gaps, registros de dependencias/freshness/monitoring,
roadmap, assurance y ledger reciben la resolución. FAIL512/514/515 corregidos;
511/513/516/517 conservan sus límites abiertos y trigger. El global24.14.1 no
se altera; próximos gates Node deben seleccionar explícitamente el candidato
local verificado. No se declara limpio pnpm compilado, todos los runtimes,
Windows, target, monitor ni el sistema completo. TEST38/EVID29 describe sólo
esta auditoría; readiness D–H y aceptación integral siguen BLOCKED.

Staging `%LOCALAPPDATA%/Temp/elite-v319-952f1100edd1424f842acd958e3e1f75`.
Los scripts de preparación/modificación son de una sola ejecución. Retomar desde
checkpoint, no repetir mkdir/copies/mutaciones sobre destinos usados. Reproducir
en destino ausente; conservar JSON/logs negativos y no modificar PostgreSQL real.

## Receipts del staging

| Archivo | SHA256 |
|---|---|
| metadata-identities.json | a785b5fd641a90fde0c3fe3f91843dd80b78f2fba8f0fb62f821825acb4eaef4 |
| metadata-receipts.json | 133410b2a4dbea08122a31a1be6c651072779f34565a88c43170357174024070 |
| core-advisory-receipts.json | db017dc9a52df5506c5220e97bd42b38d010d05d3db5f0ad22261bba2947994d |
| gpg-verification.log | 1776a2ca80e2361808cf34f9bfe39d3f57d3bae1c18abb68230424cffea2cfb6 |
| gpg-verified.log | ff49ca1ace566aec88ee59673be6ffbc5c515e5699f73d45fe74d7ce8411abdd |
| metadata/SHASUMS256.authenticated.txt | c614a913a302d35d1c3a1348d07d84c52e28dceb0b97fd91e4ade05906ca2bd7 |
| candidate-runtime.json | 599863571693e82fbcb0306897395c9355140a9e72422d7f307b9e54c2f68d69 |
| candidate-runtime-files.json | 5a637c17c3ab05cd0bb1088d3711ee43a4fab9b80da7c574bb1bb3309afbb261 |
| baseline-runtime.json | fa0970df29b5c6d96d9a9f19d69d37acb22e3beacb72817a9c76d8af2be5c889 |
| node-alias-probe.json | 8abe536182d5e4ec42750176f2bc8e24acc03c8f21c4b63c3c2bbb0342235992 |
| node-alias-24.14.1-osv.json | b0104fca6360c3dd17924ac7ad2510fd45aa572ef646a6e89541f58039817dcb |
| node-alias-24.20.0-osv.json | 459036866f3b63b407ee19967188d041a0ad59cf44a5566942235f9aff601672 |
| official-probe-adaptation.json | 5f5360118d09253547e7fa96b6264b3b68c9c5731551fd76268cb16151b3c223 |
| official-probe-osv.json | 4b49d1d2fbedc442d567ded5da918287b9b161b4fc3416b82354fc5032043076 |
| official-node-advisory-results.json | 11e6df2127411e32902ea05bc52f2216d1acb74856d343948103aa5fcf211cda |
| official-probe-negative-results.json | 8c17b5432a8ad553ee99b5d8dcb4aef840fde6d9612a607ca2e1e9cfa454cce0 |
| node-runtime-gap.json | ab83b3a6391e4fe5b14426a47cee748567e3292997f1b151e974ec0d2ad506c0 |
| node-runtime-gap-receipt.json | 0bd191394fc64b3d4b1a6b947128a656e08a483b6c760fc3cd3fc6cad8371208 |
| gap-audit.md | 0602e19277946a61ab2270ae406984e8801c93157fef0130e0d3d8d1a28a3649 |
| bundled-npm-inventory.json | 4ca20a27361e1994e9152040ccbda6d4358cedc12caf3a298df7037ab67eb733 |
| bundled-npm.cdx.json | c7ab2995a71d2e2907ab4b20eeec381cfb6a352031b4db0fea144b5825ebe5d8 |
| bundled-npm-osv.json | 157887adf084b29b4e40f093361b9d52767f277229005f24e32c90695e5f9c67 |
| bundled-npm-all-osv.json | 7245150924e77949c60931f5d339bc41ac69091f9247307571984124d3b01f16 |
| bundled-npm-result.json | ff9b6e093c02b9eebad30c95baaec94aacbe3ccb3ffb18e7b0bcd9fca8f4c46f |
| npm-advisories/receipts.json | 517695b299c90122543a0e30fb829e64efe5954ad285ac05de95c85159b45fb5 |
| node-only-selection.json | 9275b6a3038fc353921d2cab35b8590defb0a13f60d854373c686542f9c50c32 |
| pnpm-artifact-result.json | a2e3f9d1f7c1c84564a4ed7331e3463117ea8398de931e57b934a32982a59649 |
| pnpm-artifact-audit/files.json | 1a22da836259658d8e7230b3a72b65cefb7d9c1cb8872933e4610f6d99658d1e |
| pnpm-launches.jsonl | ff6b425e2f7cb3c14dfabb967d4908ec54a250b5b8baa8c796343d9b947d5ff5 |
| pnpm-node-identity.json | f055ed9d7d9166002ed2c17ed222fa9b3ba40855ecb4aed0e5a2acbba499c41e |
| node-only-pnpm-node-identity.json | 5e23f4ec84303f30b4a91aed3b7d4ef9b96e52cef713a6eddbf326dd50698f73 |
| web-tests.log | ef13761351e9c800d1ede13e35f5f1d82887dcb89c16906116afb082a7b7487c |
| web-build.log | 857e880ad7ea4e76ab319768161c4e076b84704a427c8519bd91d43392701403 |
| candidate-browser.log | 04297db98c626232c66dc4440f0af6108625295ec58825de637adee5ad227fd0 |
| candidate-agenda.log | 02a5a7cc68382d2f252e89962a68764a3bb0c5c34d32b564605de738ac0cbfed |
| node-only-web-tests.log | 29d3dc5a76344817aa08b7b895e90e92a0d544ad686ad7b99452ac14fd7ee7de |
| node-only-web-build.log | 7d0f1fab85e23277a5ad55e7232ba7a9e0ba2495fa9b252feb0a787954c37c44 |
| node-only-browser.log | adfcf177048a66fe551bd498200c4256586520192abf50954f11588bfef0d33a |
| node-only-agenda.log | be208d4a94cbc333381520be9640929ae5ea4eb95cb974ada1da863182ae6762 |
| rebuild-parity.json | ba576956914e7c5389e1a78d4e36a265f228614afc3a414d1ff4cca473ca7d83 |
| runtime-smoke-results.json | af9074152bce71f053d27eba77c7e19904a4c2369325ba15caadfdd973f76cf6 |
| node-only-candidate/MATERIALIZATION_RECORD.md | 602f96a2a5674b0805900ea856a19e438c7e8ac6b1d52926030eae9718037280 |

## Corrección de clasificación del ledger

Checkpoint65 fue válido; el primer verificador rechazó4 filas locales OPEN en el ledger reusable. Log library-final.log SHA256:2010457a91da471b07c03d13148c24dec4aec7d8b1922a2009332bda5ee4ffe2. Conforme al contrato,2226/2228/2231 conservan REJECTED_COMPONENT y2232 UPSTREAM_OPEN, en su sección correspondiente. No se altera el verificador ni se cierran511/513/516/517 del proyecto. FAIL518 registra el error propio y su corrección;2233 conserva la lección.

## Verificación final de biblioteca

VERIFY_LIBRARY_PASS:160 packs/1437 archivos/752 Markdown y51 perfiles. Franquicia67/746. El verificador usa el checkpoint66 válido; log library-revalidated.log SHA256:8c34a37c2f0c57367bf31d09f2be79a4d776ac6df5db1466ff8e8f5be323a398. TEST38/EVID29 plan PASS; evidencia integral BLOCKED. Checkpoint67 conserva93refs/12must_read y continúa con el gap Node oficial y SCA de herramientas/runtimes; no promoción ni cierre de los incidentes upstream.
