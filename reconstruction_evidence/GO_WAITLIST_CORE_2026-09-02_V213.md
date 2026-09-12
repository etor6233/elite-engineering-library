# Go Waitlist Core — V213

## Resultado estrecho

V213 materializa `GO-WAITLIST-CORE 0.1.0`: lista de espera tenant-scoped con posición por orden de llegada y máquina de estados fail-closed.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/waitlist/waitlist.go | 14c4b10a48919b2c030e265316743d9596d27546b04ec16d8cde4fb91eb73e82 |
| internal/waitlist/waitlist_test.go | 4afccea308464aaad2b600b1ed4738a960205bc177dad24e4ff5c65fa5aeef96 |

SHA-256 del pack: `b8b4b5768f78181bfa9a4a91d6900b5d9b1e3ec9ce713f5a58d87d34c333d6a5`.

## Verificación

- `go test ./internal/waitlist/ -count=1`: 6/6 PASS.
- `go test ./... -count=1` (37 paquetes): PASS. `go vet`: exit 0.
- Round-trip: 2/2 bloques reproducen byte a byte.
