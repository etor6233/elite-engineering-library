# Go SMS Channel Adapter — V189

## Resultado estrecho

V189 materializa `GO-SMS-CHANNEL-ADAPTER 0.1.0`, el adapter SMS (Twilio-compatible) con basic auth y body form-encoded; credenciales inyectadas en runtime y tests con `httptest`.

Completa los canales de mensajería: WhatsApp/Instagram/Facebook Messenger (Meta Graph), email B2B (relay HTTP) y SMS (Twilio-compatible).

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED` sobre la API oficial Twilio Messaging.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/msgchannels/sms.go | f24a4c06c6c7c9c288506f9c1809e9d092ff689fd8623379bb5be25a1c0a1617 |
| internal/msgchannels/sms_test.go | 4b8e9093b8eae56b705faf7decf8bab3332435ee0aadc06c56edc5d88b4c23ba |

SHA-256 del pack: `da65117016b5d6054895a4c5a58abc55b724038a4e79499d2bfd17604b37c59b`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/msgchannels/ -count=1`: **8/8 PASS** (suma SMS config + envío con basic auth/form body).
- `go test ./... -count=1` (12 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 2/2 bloques reproducen byte a byte.

## Cobertura de invariantes probada

1. Basic auth con credenciales del entorno.
2. Status no-2xx → error (fail-closed).
3. URL absoluta http(s).

## Condiciones residuales

- Cuenta Twilio/relay real: runtime del proyecto.
- Gate pnpm del frontend: entorno.

V189 completa los canales de mensajería; no ejecuta contra cuentas reales sin credenciales.
