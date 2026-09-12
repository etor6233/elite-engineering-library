# Go Channels Core — V184

## Resultado estrecho

V184 materializa `GO-CHANNELS-CORE 0.1.0`, la abstracción de canal del agente: contrato `Channel` (receive/send), registro por código y dispatcher que rutea mensajes entrantes al responder y devuelve la respuesta por el mismo canal, sin ampliar tenant ni canal.

Los adapters de SDK reales (Meta Graph para WhatsApp/Instagram/Facebook Messenger, Twilio para SMS/Messages, SMTP para email B2B) son `CONDITIONED` (requieren cuentas y red). El adapter WhatsApp ya existe (`PYTHON_META_WHATSAPP_CLOUD_ADAPTER`).

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — core `AUTHORED`.

## Archivos materializados (4)

| Archivo | SHA-256 |
|---|---|
| internal/channels/channel.go | cc160e33117b4e990d1c097e67bd0ddcd5996a68619b386f8f6bce9988289b5b |
| internal/channels/registry.go | 38dec01943776ea47e30b3e82a8aa3892a5a00a831abf97d8381188b8178edce |
| internal/channels/dispatcher.go | ad39fed0db2dea9ba54a3049db8d50153b00fab29330e0ae935b45ccb4d58d63 |
| internal/channels/channels_test.go | b8006ac4477f766279067baf3e1f1298ce1e6ca5f50f72595471d0dc41f7ff27 |

SHA-256 del pack: `bacd1007b906cb8997c10e5f7dff96de1a837518df2591dcb0c2b0f7a63faa07`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/channels/ -count=1`: **4/4 PASS** (validación de mensaje, registro un-canal-por-código, dispatcher rutea y responde por el mismo canal, canal desconocido).
- `go test ./... -count=1` (9 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 4/4 bloques reproducen byte a byte; `go test` sobre el árbol materializado PASS.

## Cobertura de invariantes probada

1. Un canal por código (duplicados rechazados).
2. La respuesta sale por el mismo canal y al mismo `ExternalID`.
3. Canal desconocido → `ErrUnknownChannel`.
4. Tenant pre-configurado, nunca inferido del contenido.

## Condiciones residuales

- Adapters de SDK reales (Meta Graph/Twilio/SMTP) con sus cuentas: CONDITIONED.
- Resolución de tenant por contacto (canal→tenant): composición del proyecto.
- Frontend multi-rol + onboarding: pendiente.

V184 añade la abstracción de canales; no declara los SDK reales conectados.
