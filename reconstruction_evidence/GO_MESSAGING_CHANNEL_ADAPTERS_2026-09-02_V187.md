# Go Messaging Channel Adapters — V187

## Resultado estrecho

V187 materializa `GO-MESSAGING-CHANNEL-ADAPTERS 0.1.0`, los adapters de canales de mensajería: Meta Graph (WhatsApp/Instagram/Facebook Messenger) con verificación HMAC-SHA256 de webhook y envío de texto, y email B2B vía API HTTP (SendGrid-compatible). Credenciales inyectadas en runtime; tests con `httptest` y vectores HMAC conocidos.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED` sobre el contrato `channels.Channel` y las APIs oficiales (Meta Graph, SendGrid).

## Archivos materializados (4)

| Archivo | SHA-256 |
|---|---|
| internal/msgchannels/metagraph.go | 5dd17bdc4afb336921c86e2e2acf3b29075435bd4bb93f618fb46d579574115d |
| internal/msgchannels/email.go | c9cb0931a2ef3bb68c61a13542b05cb01d253e7fe0bf9dde2a0ee9d50f98559a |
| internal/msgchannels/metagraph_test.go | 0a612b41352d56bffa4ec4497333e3fd68fcd95c735083099b59422f2e2dfaf7 |
| internal/msgchannels/email_test.go | b941460fe1e6df733e0661fb310666e549b8bdfa42bc519d5e19548a25a66c07 |

SHA-256 del pack: `abddef9489caafcb97de21b1083ce947ac73ec4474c833d507ad6f17f04d101c`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/msgchannels/ -count=1`: **6/6 PASS** (firma HMAC válida/inválida/no-sha256, parseo de webhook, envío de texto con servidor falso, config Graph con negativos, envío de email con relay falso, config email con negativos).
- `go test ./... -count=1` (11 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 4/4 bloques reproducen byte a byte; `go test` sobre el árbol materializado PASS.

## Cobertura de invariantes probada

1. Firma HMAC-SHA256 del webhook con compare constante.
2. Status no-2xx → error (fail-closed).
3. Credenciales sólo del entorno (Bearer).
4. URL/endpoint absolutos http(s).

## Condiciones residuales

- Cuenta Meta/relay email reales: runtime del proyecto (tokens inyectados ahí).
- Canal SMS (Twilio-compatible): mismo patrón HTTP, pendiente de agregar.
- Binding del `Domain` a los endpoints Go: pendiente.

V187 conecta los canales de mensajería; no ejecuta contra cuentas reales sin credenciales.
