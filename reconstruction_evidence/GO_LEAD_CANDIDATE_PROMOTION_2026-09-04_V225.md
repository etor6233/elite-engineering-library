# Reconstruction evidence — GO Lead Candidate Promotion V225

## Resultado

```yaml
pack: GO-LEAD-CANDIDATE-PROMOTION
version: 0.1.0
status: REBUILD_VERIFIED_CONDITIONED
verified_at: 2026-09-04
pack_sha256: 99a4ea0e08b8a54ab875cbce3a37931db5d0a08dcde09894a469ac8451e8955f
materialized_files: 7
roundtrip_hashes: PASS_7_OF_7
go_test_full_with_postgresql: PASS
go_vet_full: PASS
go_build_full: PASS
postgres_migration: PASS_0045_UP_DOWN_UP
postgres_sql_tests: PASS_31
```

## Journey demostrado

```text
candidato 0044 durable
  -> decisión explícita con purpose/policy/evidence/mapping
  -> test lead bloqueado
  -> contacto derivado sólo de campos mapeados
  -> crm.lead + consentimiento + vínculo source + outbox + idempotencia
  -> commit atómico
  -> replay exacto sin duplicar / divergente fail-closed
```

El test PostgreSQL usa el runtime oficial 18.6 del entorno con data checksums. La suite completa pasó con `TEST_DATABASE_URL` y `DATABASE_URL` reales: 46 paquetes, incluidos `leadstream`, `leadpromotion`, CRM/journeys y adapters de plataforma. `go vet ./...` y `go build ./...` terminaron exit 0. La migración 0045 pasó `up`, `down`, `up`; las 31 pruebas SQL pasaron.

## Hashes reconstruidos

| Archivo | SHA-256 |
|---|---|
| `internal/leadpromotion/promotion.go` | `8a32c465cd46f9da0ad322080d0ce629f8c6cd38531129de4828ab9eb5f1d88d` |
| `internal/leadpromotion/promotion_test.go` | `8b9aadc43a95efccb4bdd705b9cfc6c8d8c35baf52e13eff58698a4de79c423b` |
| `internal/platform/postgres/lead_promotion.go` | `31001b1c9476e0d5bd2cd4d13cd3edfa5d0f4e832a03fff58dc9996e3624e6f2` |
| `internal/platform/postgres/lead_promotion_integration_test.go` | `fa60af0cf11f615e402657bbd4a8239971221fd3dac807eb325dd10f580e2cb4` |
| `db/migrations/0045_lead_candidate_promotion.up.sql` | `2d6b404be0e8611bbc7951b8759ac70ecd326390fb53b496047dae3782ea139d` |
| `db/migrations/0045_lead_candidate_promotion.down.sql` | `9988d288b62cc7e09b2e755790f3e4719e35dad4816167b840b00537e6b16c9f` |
| `db/tests/0045_lead_candidate_promotion.test.sql` | `02b6e88cab6ccf47431c9412b8950f834341229cddb0e4f3a56199245e56c498` |

## Procedencia y límites

Cinco archivos son `ADAPTED` y dos `AUTHORED`; ninguno es código copiado de Google, AWS o PostgreSQL. Gobiernan claims estrechos los contratos oficiales de Google Lead Form, entrega al menos una vez/idempotencia de AWS y transacciones/constraints de PostgreSQL. El pack no inventa consentimiento ni mapping: ambos son inputs obligatorios del proyecto.

No demuestra todavía Meta/TikTok, conversación durable, entrega por canal, LLM/tool calling, postback de conversiones, RLS/cifrado/retención del target ni aceptación productiva. Por eso permanece `CONDITIONED`.
