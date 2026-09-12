# Go App Wiring — V191

## Resultado estrecho

V191 materializa `GO-APP-WIRING 0.1.0`, el ensamblado del runtime de la franquicia: agente + LLM + dominio + costo + aprobación anti-estafa + FinOps en un único `App`, con validación de credenciales fail-closed que lista exactamente qué falta.

Es el punto de unión: `New(config)` devuelve todo cableado o un error que enumera los inputs faltantes. El agente no adivina ni se bloquea.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED`, ensambla los packs propios.

## Archivos materializados (3)

| Archivo | SHA-256 |
|---|---|
| internal/app/config.go | 1aa14bee28813b363f3cffe7a112f5f894cb86c17e7202da98733b1218b09884 |
| internal/app/app.go | 215bbd08d93a70264704140a41de006d77a7eac08c8a6d42bae6e04ad157ae43 |
| internal/app/app_test.go | 059e01726ebc58e5c6c55ac7915aa31eb0b2a08fb39db7e863f1ee8edbd8282e |

SHA-256 del pack: `000849265706e3dc307698ba0782a3bfb35785cf9f2f6b0c46b9aec95b942c55`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/app/ -count=1`: **3/3 PASS** (lista de inputs faltantes, ensamblado completo no-nil, fail-closed ante config incompleta).
- `go test ./... -count=1` (14 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 3/3 bloques reproducen byte a byte; `go test` sobre el árbol materializado (app + 9 dependencias) PASS.

## Cobertura de invariantes probada

1. Inputs faltantes → error con lista exacta (fail-closed).
2. Dinero nunca auto-aprueba (`AutoApproveMinorUnits: 0`).
3. Todos los componentes no-nil al ensamblar.

## Condiciones residuales

- `main.go` con servidor HTTP (rutas públicas + webhook) que usa este `App`: cablear al arrancar.
- Gate pnpm del frontend: entorno.

V191 une el runtime completo; la inyección de credenciales es runtime del proyecto.
