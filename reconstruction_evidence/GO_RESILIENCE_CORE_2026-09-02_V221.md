# Go Resilience Core — V221

## Resultado estrecho

V221 materializa `GO-RESILIENCE-CORE 0.1.0`: circuit breaker (closed→open→half-open con sonda única), timeout y retry budget acotado con backoff inyectable.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/resilience/resilience.go | 21f34f08a8d04a23ded5544a3b81dbcd4b26ca273021e647db7302a761c4812e |
| internal/resilience/resilience_test.go | 7dc894ebffecceb734967643f99e95bdf93223bc21402153bca190d92266849d |

SHA-256 del pack: `aa8bd4932dd4b6df9d65f02119aba914a79ad58220083e61ea1a1a7397f611d1`.

## Verificación

- `go test ./internal/resilience/ -count=1`: 6/6 PASS.
- `go test ./... -count=1` (41 paquetes): PASS. `go vet`: exit 0.
- Round-trip: 2/2 bloques reproducen byte a byte.
