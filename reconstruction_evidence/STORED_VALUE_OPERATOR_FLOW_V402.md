# V402 / checkpoint281 — gift y loyalty: operador HTTP, expiración y entrega

Sucesor de STORED_VALUE_TENDER_V402.md/checkpoint280. **Candidato local conectado, aún no publicado como pack**. Canónica82/1039 y45/48 intactos; T2802/TEST02 globales abiertos. Fuente aislada OdooCommunity19.0@99edb6dd82b7b560930c00b03b694ba700785370, LGPL3/ADAPTED; Go, SQL, transportes y UI son AUTHORED glue. No se adquiere ni cambia una dependencia. Este informe no sustituye source/notice/G0–G8 y reconstrucción exacta pendientes.

## Delta demostrado

- FAIL824 real: vencía la aprobación durante espera de outbox y se confirmaba una operación. RED retenido. Dos constraints SQL diferidos usan clock_timestamp al COMMIT. GREEN: cero operación/decisión/evento cuando vence; propuesta anterior continúa pending. Una inserción de propuesta vencida también revierte. Sólo el fixture usa leases cortos; perfil conserva mínimo60segundos. El core gift/loyalty conserva6operaciones/6eventos.
- HTTP conecta propuesta, revisión separada, recuperación, perfil, paginación acotada por pedido, allocation y financiación. Principal y organización están ligados al perfil del servidor. Versiones/importes int64 viajan como strings; payload/receipt JSON se conserva como texto y BFF comprueba SHA. Read valida el hash del recibo persistido. Moneda/pricing/points no se calculan en navegador.
- gift-http-connected-first: cinco tests superiores sin skips,65migraciones. Emisión/acumulación/canje y decisiones pasan ahora por HTTP real y SQL/IPC; fixture de bearer declarado en esa suite. Gift parcial y loyalty continúan a SDK/callback/GET/release. Se rechazan organización distinta aun si el actor tiene acceso a ambas, version JSON numérica y permiso de funding ausente.9007199254740993 permanece exacto.
- Portal /franchise/stored-value: consulta de pedido, programa fijado, propuesta, cuenta/importe/puntos de la revisión, motivo, separación de personas y comprobante. Guarda request antes de POST; resultado incierto mantiene referencia, GET recupera y reintento usa la misma solicitud. Funding total continúa al panel de entrega existente. No billetera pública/bearer ni ERP completo.
-49tests focales de BFF/handovers PASS; typecheck/build Next PASS. Incluye XOR estricto payment_attempt_id/funding_receipt_id y forwarding server-side. Build y módulos reutilizan la instalación contained existente; no se promociona pnpm por este PASS.
- Chromium real: JWE de rol y OIDC RS256/JWKS reales locales, BFF Next, PG y cálculo Python. Opera emisión→revisor distinto→canje→revisor distinto→funding→preparación→checklist→aceptación cliente→recibo comercial. Recupera respuestas perdidas de funding/preparación/release sin segundo POST. Resultado2operaciones,1funding,1handover aceptado,1recibo comercial,0payment_attempts. Roles/cross-org negados,0pageerrors y ancho390sin desborde. Screenshot inspeccionado. No login live ni clave de usuario.
- Contrato del host opcional: factory real, flags/perfil/hash/tenant/org/script/manifest, preflight del proceso; desactivado no monta rutas, activado exige auth. Pool lazy sin conexión en ese test, declarado; SQL se prueba por las suites conectadas.8subcasos PASS; go vet ./... y go build ./... offline PASS.
- Downgrade:0066/67/68 rechazan historia; snapshot completo antes/después byte-idéntico.0067 elimina sólo objetos propios sin CASCADE; una dependencia externa fixture impide el downgrade y queda intacta. Base vacía down/up de las tres migraciones PASS. PostgreSQL propio detenido.

## Fallos y límites preservados

FAIL825: logger inexistente en hook de main; corrección a slog observada, compilación posterior PASS. Extracción del helper de fixture movió op incorrectamente: logs de compilación antes/después preservados; no se contó como defecto del algoritmo. FAIL826: optional-property type bajo exactOptionalPropertyTypes; parámetro acepta undefined sin cambiar XOR. FAIL827: locator getByLabel exacto no resolvió select; snapshot accesible demostró combobox Operación; test corrige role/name y pasa5.2s. Build de producto no cambió por ese ajuste. REDs anteriores no se reescriben.

Falta publicar notices/licencias/source/replacement con G0–G8 del claim estrecho, pack(s), owners/dependency closure y rebuild exacto de perfiles afectados. El resto del gap map T2802 se mantiene abierto; no se convierte en CONDITIONED por credenciales lo que aún requiere código/gates. ARCA penúltimo; Daybreak/libxml2 último sólo expediente diferido.

## Recibos

Todos los paths siguientes son relativos al staging elite-v402-library-infra; JSON compañero conserva el inventario y hashes.

| Recibo | SHA-256 |
|---|---|
| `gift-expiry-red/result.json` | `24b5e72dccb7a91aa53a1edf322de1dbcf8a6715616a3a42fa2111085858ec9f` |
| `gift-expiry-red/connected.log` | `6a88bc2d7286c84f910ffc7e6abeaffbc8bc3af64801c40934bea3c4f71c5ffa` |
| `gift-expiry-green/result.json` | `a6bd7be82b12afc46c3c7684a49fb2ce7e93bf5fe3d2971a4118978f430201cc` |
| `gift-expiry-green/connected.log` | `f8582e2724a01e8288fbedffc1688ab93e4daed53d59304e2077aae62846115f` |
| `gift-http-connected-first/result.json` | `939dc1297eb885ea1a32d062ee61b3953cc5041ecced6b8a3a465aba5b4d0428` |
| `gift-http-connected-first/connected.log` | `9651180b749dbc0085746beac945a9a5d3a75264b502129d4c1dbcf8ca639278` |
| `gift-web-first/focused-tests.log` | `0a8bfa5c2fba80e49cbaae09b9b87a601ab458c2209bd9d4b26134bee8492f19` |
| `gift-web-first/build.log` | `3bba6eba726d3759e36ee23896eeedc737f8967fd4472d6757975c636e049d03` |
| `gift-web-first/typecheck.log` | `c214f20b44236d01495dd0b6fdf220c012761748a174af7da97665f0fbfc0af8` |
| `gift-web-first/typecheck-corrected.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `gift-browser-pg-first/result.json` | `0be1614ae3ef745ea0587452c69281c5d1e907fa13a3d8da2b7707d3924a3708` |
| `gift-browser-pg-role-locator/result.json` | `95a940416cd5768fea0b3d45f3de6d558e6d53904efa60f9c67d8cca63352d04` |
| `gift-browser-pg-role-locator/connected.log` | `0613ac0dc90bfa520f7f7f80ebc165c336c866f384693b37f0b59d97b934a92d` |
| `gift-tender-downgrade/result.json` | `8b82f5da324c17d87550b2fd5d37aea4427f88d81ee4c727090d629fc0ff3b8d` |
| `gift-final-go-gates.json` | `4f0e640c5ef9c30dde0f2d989aad300d027b0f9c590ddb0adbdbf7c397560019` |
| `gift-final-host-contract.log` | `5f98b77ff6d6b89c6a68704479e34658e376ee26ccbc4435f44cd74ab4973613` |
| `gift-final-vet.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `gift-final-build.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `gift-tender-operator-manifest.json` | `3e1a0d3b0d8a92550d39210ba8911ebf0526c8297bb3d6d94ae916f3333425e4` |
| `gift-tender-candidate/stored-value-browser-artifacts-2351150312/browser.log` | `ce682fae8e3d9029032d2865c2fb6123db209e0e71b54beb1ed7648c5fa512cd` |
| `gift-tender-candidate/stored-value-browser-artifacts-2351150312/api-post-counts.json` | `b95a525ea53bb4de2ff9420e65527d5d18b25458e56acfc4dd8a332161d0e2f6` |
| `gift-tender-candidate/stored-value-browser-artifacts-2351150312/browsers/stored-value-connected-sou-ce493-delivery-without-a-provider-chromium-desktop/stored-value-recovered-mobile.png` | `5b916d50b2f68992c3a6d288941bd92746d1afead7278bb18c431d6667ccc616` |
