# Go Data Residency Core — V211

## Resultado estrecho

V211 materializa `GO-DATA-RESIDENCY-CORE 0.1.0`, la aplicación de política de residencia de datos: tenant-scoped, deny-by-default y fail-closed.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED`, gobernado por GDPR Capítulo V (transferencias) y requisitos de soberanía de datos.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/residency/residency.go | 1e59c31a34c8ea8af06d0b27fa86488c0ebf30981b140b33dfcc489a25b572ae |
| internal/residency/residency_test.go | 4f962bba4b9a4560262442547b4e47034a1292936899c14ede9590c72e8a0e89 |

SHA-256 del pack: `0f229489bff15f573ea815efcc2f526f4c708db750674b50d07fc0afd7f9f8ac`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/residency/ -count=1`: **6/6 PASS**.
- `go test ./... -count=1` (32 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 2/2 bloques reproducen byte a byte.

## Cobertura de invariantes probada

1. Región permitida aceptada.
2. Región fuera del conjunto → rechazada.
3. Categoría sin política → deny-by-default.
4. Política inválida (tenant vacío, regiones vacías, región mal formada) rechazada.
5. Reemplazo de política y dedup de regiones.
6. Aislamiento de tenant.

## Condiciones residuales

- La colocación física real (replicación entre regiones, respaldo) es composición del proyecto.
