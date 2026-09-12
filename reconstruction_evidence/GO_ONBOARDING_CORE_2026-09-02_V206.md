# Go Onboarding Core — V206

## Resultado estrecho

V206 materializa `GO-ONBOARDING-CORE 0.1.0`, el onboarding del franquiciado: pasos ordenados, checklist con progreso y completado idempotente, tenant-scoped.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED`.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/onboarding/onboarding.go | d8473e8a77d3e440fee994c75cd28173f0e74b405a20a23e79df24af3d395013 |
| internal/onboarding/onboarding_test.go | 8b7a2c9aea28f9c76975c397ce3a2e6d552e4dde5feb161a944eb2d8b6d48fc2 |

SHA-256 del pack: `0348f51991a57109678e608451fb5e2cd39c12b19c2eef6fbdb30a9b5315aca5`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/onboarding/ -count=1`: **5/5 PASS** (progreso/done, idempotente, paso desconocido, dedup/inválido, aislamiento tenant).
- `go test ./... -count=1` (28 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 2/2 bloques reproducen byte a byte.

## Cobertura de invariantes probada

1. Pasos no vacíos y deduplicados.
2. Paso debe existir.
3. Completado idempotente.
4. tenant-scoped.

## Condiciones residuales

- Activación del franquiciado (roles/permisos) al completar: composición del dominio.

V206 cierra el hueco de onboarding; la activación es composición.
