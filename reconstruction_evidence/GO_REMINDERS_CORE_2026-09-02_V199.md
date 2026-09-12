# Go Reminders Core — V199

## Resultado estrecho

V199 materializa `GO-REMINDERS-CORE 0.1.0`, el scheduler de recordatorios de turnos: due (vencido y no enviado) sobre SMS/email/push, marcado enviado una sola vez, tenant-scoped y fail-closed.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED`.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/reminders/reminders.go | ba68269128726f271850276be60e888328d4b63aa9e02f8ac821b3dafed6e900 |
| internal/reminders/reminders_test.go | 2b4dd8d59d514658b8a6859b85ecf843efa89aff5c2433548190383fbe319a79 |

SHA-256 del pack: `efe16e3d13cfa7d109d2783d665f33b9e61d40c861a13d186beaef07e70dc0cb`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/reminders/ -count=1`: **4/4 PASS** (due y mark-sent, rechazo pasado, inválido, duplicado).
- `go test ./... -count=1` (21 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 2/2 bloques reproducen byte a byte.

## Cobertura de invariantes probada

1. At en el futuro.
2. due = at ≤ now y no enviado.
3. Enviado una sola vez.
4. Canal admitido; tenant-scoped.

## Condiciones residuales

- Envío real por canal (msgchannels): composición del proyecto.

V199 cierra el hueco de recordatorios; el envío real es composición.
