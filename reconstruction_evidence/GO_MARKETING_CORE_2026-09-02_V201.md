# Go Marketing Core — V201

## Resultado estrecho

V201 materializa `GO-MARKETING-CORE 0.1.0`, las campañas de marketing tenant-scoped: canal (SMS/email/push), mensaje y destinatarios deduplicados, programadas con envío idempotente por destinatario y completado automático.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED`.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/marketing/marketing.go | 69b63c97377db228c5f548a6fa3dcdb4a0c1f72d544cd6888d807004a2ab910f |
| internal/marketing/marketing_test.go | 9954bf39eef824da02172bfc44e398379b7bc1980d5f35bf0b3def802d4969d4 |

SHA-256 del pack: `a9122f7daeed08463f6ab5035dab94150dede1ab50798c000f682bfc29c350d8`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/marketing/ -count=1`: **6/6 PASS** (dedup destinatarios, envío y completado, dedup por destinatario, inválido, duplicado, destinatario desconocido).
- `go test ./... -count=1` (23 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 2/2 bloques reproducen byte a byte.

## Cobertura de invariantes probada

1. Canal admitido, mensaje y destinatarios no vacíos.
2. Destinatarios deduplicados.
3. Envío idempotente por destinatario.
4. Completado automático cuando todos enviados.
5. tenant-scoped.

## Condiciones residuales

- Dispatch real por canal (msgchannels): composición del proyecto.

V201 cierra el hueco de campañas; el dispatch real es composición.
