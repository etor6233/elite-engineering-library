# Go Warranty & Claims Core — V214

## Resultado estrecho

V214 materializa `GO-WARRANTY-CLAIMS-CORE 0.1.0`: ciclo de garantía/reclamos tenant-scoped open→in_review→approved/rejected→resolved/closed, fail-closed con nota obligatoria al rechazar/resolver.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/warranty/warranty.go | f84fd2a588045f2f10e155ac12b660aac13352504e1366948b5716c8ef8e4b6a |
| internal/warranty/warranty_test.go | 381e338341277d642f84c0e7489f421c7656b308b0786372d3afdee817f1aced |

SHA-256 del pack: `5633258019a207a5a6317e90123b77c20f1970e52331d8f6a0a7bfa368d2b723`.

## Verificación

- `go test ./internal/warranty/ -count=1`: 4/4 PASS.
- `go test ./... -count=1` (37 paquetes): PASS. `go vet`: exit 0.
- Round-trip: 2/2 bloques reproducen byte a byte.
