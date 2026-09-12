# Go Applicant Onboarding Core — V223

## Resultado estrecho

V223 materializa `GO-APPLICANT-ONBOARDING-CORE 0.1.0`: alta segura anti-fake (franquiciado o cliente online) submitted→contact_verified→approved→activated|rejected, con deduplicación de contacto, velocidad anti-falsos, doble control y activación sólo desde aprobado+verificado.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/applicant/applicant.go | dd7a75914ea9dd00fd2aeecd9208cb185cf4294d11951f5f5a6bec45fec15a6f |
| internal/applicant/applicant_test.go | 5f06bda12ef32e80418b730afabb9c2d3461cb8da3f84df334120ea100cfe32c |

SHA-256 del pack: `3e99c4b1f11e165f9948953c7b933c3d4caa8c31a3efc5b0c4c47a56a2d07d1d`.

## Verificación

- `go test ./internal/applicant/ -count=1`: 8/8 PASS.
- `go test ./... -count=1` (43 paquetes): PASS. `go vet`: exit 0.
- Round-trip: 2/2 bloques reproducen byte a byte.

## Condición residual

El KYC documental/biométrico real (documento + liveness) es un adapter CONDITIONED a un proveedor externo; este pack gobierna el ciclo y deja el hook, no inventa verificación biométrica.
