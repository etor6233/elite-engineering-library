# Go OTP Verification Core — V222

## Resultado estrecho

V222 materializa `GO-OTP-VERIFICATION-CORE 0.1.0`: verificación de email/teléfono por código de un solo uso (hash SHA-256, TTL, intentos acotados, fail-closed).

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/otp/otp.go | c72b8920acc63eb109cfcbb6d59b4d81102abd6f5f0f893b75cb18596d1ea0ae |
| internal/otp/otp_test.go | d8a26827750d130a65e8b52bc8790896fd266efec5e32b2f0016c493d518d49c |

SHA-256 del pack: `946c28a03dc179bef203d1264f6e238d884cbae30a1e1a0551c6682a7580b33c`.

## Verificación

- `go test ./internal/otp/ -count=1`: 5/5 PASS.
- `go test ./... -count=1` (43 paquetes): PASS. `go vet`: exit 0.
- Round-trip: 2/2 bloques reproducen byte a byte.
