# Go Referrals Core — V196

## Resultado estrecho

V196 materializa `GO-REFERRALS-CORE 0.1.0`, el programa de referidos tenant-scoped: códigos únicos, sin auto-referido, un referido por persona y máquina de estados pending→converted→rewarded fail-closed.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED`.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/referrals/referrals.go | dcf15673467caf0f57746125c7e4e6285ab5105ad4bb2b9320e5137acf6866dc |
| internal/referrals/referrals_test.go | 3aa7541a1981127f7ffe3acd01b1ef6b724824eaf28bda497bfda1653aff23fa |

SHA-256 del pack: `4397e3aeb0236f4ab36f83dd910e016cdda24fef905580c14999d39d7049f28e`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/referrals/ -count=1`: **5/5 PASS** (máquina de estados, reward requiere converted, auto-referido rechazado, código/referido duplicados, aislamiento tenant).
- `go test ./... -count=1` (18 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 2/2 bloques reproducen byte a byte.

## Cobertura de invariantes probada

1. referrer ≠ referee (sin auto-referido).
2. Código y referido únicos.
3. Estados ordenados pending→converted→rewarded (reward sin converted falla).
4. tenant-scoped.

## Condiciones residuales

- Recompensa automática (loyalty/cupón) al reward: composición del dominio.

V196 cierra el hueco de referidos; la recompensa es composición del proyecto.
