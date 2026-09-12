# Go Rate Limit Core — V218

## Resultado estrecho

V218 materializa `GO-RATE-LIMIT-CORE 0.1.0`: rate limiting / anti-abuse tenant-scoped con token bucket (deny fail-closed, refill a tasa fija, burst acotado).

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/ratelimit/ratelimit.go | 9d94c564e4f0d291de05289e6afbc3168bda31dedfa9f47c1cf16c72e971477f |
| internal/ratelimit/ratelimit_test.go | eaeffab126de36dc0ffebd94b49ddf654103b14b83d9553711275709f04a13ce |

SHA-256 del pack: `093dcf377902bb674475beb26b2c5915c92044047c326a7010766d96cdf623d5`.

## Verificación

- `go test ./internal/ratelimit/ -count=1`: 6/6 PASS.
- `go test ./... -count=1` (41 paquetes): PASS. `go vet`: exit 0.
- Round-trip: 2/2 bloques reproducen byte a byte.
