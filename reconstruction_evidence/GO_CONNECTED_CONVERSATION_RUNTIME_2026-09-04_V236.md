# Evidencia V236 — runtime conversacional conectado

## Resultado

```yaml
date: 2026-09-04
decision: REBUILD_VERIFIED_CONDITIONED
closed_findings: [LIB-FAIL-1743, LIB-FAIL-1758]
production_ready_project: false
```

V236 demuestra localmente un único recorrido ejecutable:

```text
channel envelope
→ durable claim/hash
→ authorized contact scope
→ safety + bounded history
→ tenant token budget
→ OpenAI Responses strict tool call
→ approval gate
→ authenticated/idempotent domain command
→ PostgreSQL completion
→ deterministic outbound delivery key
→ exact replay without repeating model or domain effect
```

Todos los archivos locales son `AUTHORED`. OpenAI Responses, PostgreSQL 18 locking y Google SRE gobiernan contratos estrechos; ninguna línea se atribuye a esas organizaciones.

## Artefactos

| Artefacto | Versión | SHA-256 del pack |
|---|---:|---|
| `GO_CONNECTED_CONVERSATION_RUNTIME.md` | 0.1.0 | `c9e6f2429b1dda0b8615e93d7506d39fcd14231cd1274d6847eb62e728f35415` |
| `GO_APP_WIRING.md` | 0.2.0 | `829a4ef7a3fa2290bd55d1538081d56887101721fe0cbb59b44b2e2dc087b9cc` |
| `GO_CHANNELS_CORE.md` | 0.4.0 | `a191a1e0d5649889b5f7f2f64eba1cf018ed7889feb0fc8c5e6b831af920974d` |
| `GO_AGENT_DOMAIN_BINDING.md` | 0.4.0 | `1edb646df616936bd7637abd2b5c450befc51361042f28b5ab1b612d6ea19e91` |
| `FRANCHISE_COMPLETE_PACK_PLAN.md` | 1.0 / 56 packs | `0d6b752f6d763f38b6c13ba72649d35ec635b707f8713fe33c6e506ff6f6862e` |

## Gates reproducidos

- Go exacto: `go1.26.7 windows/amd64`.
- PostgreSQL exacto: 18.6, puerto aislado 55445, base nueva `elite_conversation_v236`.
- Migraciones limpias: 46/46 PASS.
- Migración 0046: down → up → prueba focal PASS.
- PostgreSQL: replay exacto, divergencia, aislamiento de contacto/thread, reintentos máximos, un ganador concurrente e identidad terminal inmutable PASS.
- Runtime: cita, cotización, status, safety, presupuesto, tool no autorizada y errores inciertos retryable PASS.
- App E2E con Responses/domain HTTP controlados y PostgreSQL real: segundo dispatch produjo el mismo reply, cero inferencias adicionales y cero efectos de dominio adicionales PASS.
- Suite completa: 48 packages Go PASS con `TEST_DATABASE_URL` y `DATABASE_URL` reales; `go vet ./...` y `go build ./...` exit 0.
- Reconstrucción desde Markdown: 17/17 archivos de los cuatro packs, hashes exactos y `gofmt` no-op PASS.

## Condiciones retenidas

Este expediente no demuestra proveedor live, cuenta/permiso Meta o Google, schema/auth TikTok exactos, receipt outbound live, clasificador de PII/product safety productivo, budget distribuido, IdP, carga, seguridad ofensiva, backup/restore, canary/rollback ni aceptación empresarial. Cada proyecto debe cerrarlos en su target. El runtime tampoco ofrece tools de pagos/refunds irreversibles.

## Autoridades

- OpenAI Responses y function calling: <https://developers.openai.com/api/reference/resources/responses/methods/create>
- PostgreSQL explicit locking: <https://www.postgresql.org/docs/18/explicit-locking.html>
- Google SRE release/launch gates: <https://sre.google/sre-book/reliable-product-launches/> y <https://sre.google/sre-book/release-engineering/>
