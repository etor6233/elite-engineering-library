# Go LLM Economy Core — V181

## Resultado estrecho

V181 materializa `GO-LLM-ECONOMY-CORE 0.1.0`, la capa de gobierno de costo que hace el gasto de LLM por respuesta casi nulo: cache semántica (`$0` en hit), compresión de prompt con presupuesto de tokens, ledger de gasto por tenant fail-closed y un gobernador que decide el camino más barato (cache → LLM → bloqueo).

No afirma el costo exacto de ningún proveedor (lo fija el proveedor y el volumen). Garantiza la arquitectura que minimiza las llamadas al modelo; el tokenizer exacto es del proveedor (CONDITIONED).

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — núcleo `AUTHORED`.
- `internal/search 0.1.0` (similitud coseno) — de este repositorio.

## Archivos materializados (6)

| Archivo | SHA-256 |
|---|---|
| internal/economy/economy.go | 0b319c126bf3a5ebbec8872864006d65ce7da7eb34eb418cd223d198df9d5f91 |
| internal/economy/semantic_cache.go | bf0925828c9122983016e1d910685b6ee778fcad00316e81762e06740694e9a3 |
| internal/economy/compressor.go | e6f00f9b0c910a953ffcb747c912559f287241d6262baaaba6ee6f8f16598069 |
| internal/economy/ledger.go | 0c3d3e89eb0ae356842cad62b7d8b0070b27ab47e37f4fd19a18694d37fbfe91 |
| internal/economy/governor.go | 03ea503ac1bd78b383334f6dfd5700d80c82670ce26e5a3220510ce817a40a42 |
| internal/economy/economy_test.go | 93432123b9933300a2ac5238660343833ffd8086fbb1384b3dd5028a98e6a5d4 |

SHA-256 del pack: `cde3ebad7202f4830a066a2106e17a4a4fcdb48f776173029bc6ac58ec27de80`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/economy/ -count=1`: **5/5 PASS** (estimación de tokens, cache semántica hit/miss/dim-mismatch, compresión con presupuesto, ledger fail-closed, gobernador cache→llm→blocked).
- `go test ./... -count=1` (6 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 6/6 bloques reproducen byte a byte; `go test` sobre el árbol materializado (economy + search) PASS.

## Cobertura de invariantes probada

1. Hit de cache semántica no llama al modelo (`$0`).
2. Prompt comprimido nunca excede su presupuesto.
3. Tenant sobre-presupuesto no puede gastar (`ErrBudgetExceeded`).
4. Vectores de dimensión distinta o cero no producen score inventado.

## Condiciones residuales

- Tokenizer exacto y adapter de LLM real: CONDITIONED (proveedor).
- Índice ANN para cache semántica a escala: expediente separado si el workload lo exige.
- Tiering de proveedor por región/precio: decisión de proyecto.

V181 añade el gobierno de costo del LLM; no declara el costo real de ningún proveedor.
