# Go PCI DSS Scope Core — V208

## Resultado estrecho

V208 materializa `GO-PCI-DSS-SCOPE-CORE 0.1.0`, el gobierno de alcance del CDE (PCI DSS 4.0): registro tenant-scoped de componentes, fail-closed ante datos de cuenta fuera de alcance, almacenamiento de SAD prohibido tras la autorización y SAQ obligatorio en alcance.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED`, gobernado por PCI DSS 4.0 (Req 3.2, 12.5.2; https://www.pcisecuritystandards.org). No es una autoevaluación PCI; es el contrato de alcance que la precede.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/pciscope/pciscope.go | 454151a17fb92854fac14cc2f5998fb048e667de537a32f69f41aedd8923cee9 |
| internal/pciscope/pciscope_test.go | cd2dd60efad4490ea9168a1c7d550eba97c0a7dd31fc3e7e9692119047c09c31 |

SHA-256 del pack: `c021d2c6fe677a5cfe385985f20928790cbf9a5bee7527f2bdf8b3ba62ec314f`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/pciscope/ -count=1`: **7/7 PASS**.
- `go test ./... -count=1` (30 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 2/2 bloques reproducen byte a byte.

## Cobertura de invariantes probada

1. Componente con datos de cuenta (PAN) fuera de alcance → rechazado.
2. En alcance con SAQ válido → aceptado.
3. Fuera de alcance sin datos de cuenta → aceptado.
4. Almacenamiento de SAD → prohibido.
5. SAQ ausente/desconocido → rechazado.
6. Duplicado e aislamiento de tenant.
7. `Audit` limpio sobre un conjunto válido.

## Condiciones residuales

- La autoevaluación SAQ completa, el ASV, el QSA y la certificación anual son procesos externos del operador, no código.
