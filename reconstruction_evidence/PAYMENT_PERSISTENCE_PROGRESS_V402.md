# V402 — pago durable: correcciones locales verificadas

Alcance: infraestructura de biblioteca. El bridge nuevo sigue en staging fuera de la carpeta canónica; este expediente no cierra TEST02 ni declara producción. La publicación y reconstrucción del perfil conectado siguen pendientes.

La suite PaymentCheckoutDurableProjection pasó con PostgreSQL real y un transporte de proveedor sintético: ocho workers concurrentes produjeron una sola llamada; observación exacta habilitó captured; callbacks durables invalidaron generaciones anteriores; replay no duplicó jobs; importe, moneda, pedido, intento y modo divergentes mantuvieron hold. Una devolución parcial conserva el estado captured bajo hold, una total puede proyectar refunded, y una observación anterior no reduce el reembolso observado ni revierte estados financieros terminales. Una conexión de otra organización y una cancelación entre Bind y Claim son rechazadas; la cancelación no deja fence ni evento de envío.

La prueba usa Commerce.RecordOrderPayment para crear la solicitud; INSERT de pedido es fixture y no prueba aquí cotización→pedido. Su transporte no es un cobro real ni el E2E de SDK. El E2E conectado del SDK, callback, JobProcessor y handover tiene un expediente separado cuando su revisión quede congelada. Vet y build pasaron en las revisiones identificadas por sus receipts, no en cualquier cambio posterior.

FAIL809 conservó el RED con dos fallos reales: parcial confundido con total y terminal regresando a captured. semantic-green2 y admission corrigieron esas regresiones. El intento semantic-green se detuvo por inventario55 frente a56 antes de iniciar PostgreSQL; no ejecutó producto. first descargó dependencias Go; los dos runs siguientes desactivaron GOPROXY/GOSUMDB y usan runtime exacto, loopback y bases descartables con parada de proceso verificada.

Tras esos runs se añadieron unique de provider payment reference y compatibilidad con URL nula de Stripe complete/expired. Son cambios pendientes del E2E siguiente, no cubiertos retroactivamente por esos PASS. [Contrato oficial de Checkout Session](https://docs.stripe.com/api/checkout/sessions/object), consultado2026-09-11: la URL es nullable y sólo está disponible mientras la sesión está activa; complete no significa por sí mismo que el pago esté cobrado. El SDK ya toleraba URL ausente; el guard incorrecto estaba en el glue de persistencia.

## Evidencia local inmutable

Stage: <LOCALAPPDATA>/Temp/elite-v402-library-infra. Los hashes de abajo ligan la revisión concreta y preservan la historia; no son una promesa de portabilidad fuera de esta máquina.

| Archivo | SHA-256 |
|---|---|
| `payment-pg-first/result.json` | `8be8f988445ef3e9f7946d7d1e04d11758126112ed7b5d9e85b6e8e466a6127d` |
| `payment-pg-semantic-red/result.json` | `9b1f8ce82786c3af2d6ec1c5ba841bbcb6006521f582df409af95c4b8bb8c8bf` |
| `payment-pg-semantic-red/connected.log` | `af7f8aa5ba9464a8b8c9c143e0aa3b004cfcf762bd9f02a71dd712d4a40538f1` |
| `payment-pg-semantic-green2/result.json` | `3d6ff283e5a57dcb950f21e0d66162ded79f3ecbd03c0f34210c2234a0dbd290` |
| `payment-pg-semantic-green2/connected.log` | `4b3dd3d582187f96c8141186b58f2ebd6759f066dd4c7c818a35ef95c9a205a9` |
| `payment-pg-admission/result.json` | `3cf7b733ce3d879bd2533d6b591df792efe668fa81a83703f37d4debbeea319d` |
| `payment-pg-admission/connected.log` | `159df17ea01ca085028319084b9c4e28edfacd369581b67aaf5bfe9e4dfd375f` |
| `payment-pg-admission/vet.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `payment-pg-admission/build.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `bc-payment-reference-integration.json` | `f5c5ac4f10bd8a020b7bf8cc40eb2ed1f0cf92dd1a13ea945c574d7dc19fb319` |
| `published-deltas-integration.json` | `fa46f76686e79f72bef1c7c41db88c355c3be48e1566ff0b1427f280b47e2cf1` |
