# Go Meta Lead Evidence Import — evidencia V227

## Resultado

```yaml
pack: GO-META-LEAD-EVIDENCE-IMPORT
version: 0.1.0
authority: SUPPORTED_REFERENCE
implementation: REBUILD_VERIFIED
admission: CONDITIONED
verified_at: 2026-09-04
go: 1.26.7
postgresql: 18.6
materialization: 7/7 PASS
gofmt_idempotence: 4/4 PASS
python_to_go_contract: PASS
postgres_integration: PASS
full_go_packages: 47 PASS
go_vet: PASS
go_build: PASS
```

## Claim demostrado

Los tres artefactos producidos por V226 se validan completos antes de la primera escritura: identidad SDK/API/commit, schemas, response→batch→receipt hashes, cardinalidades, scope y equivalencia raw/candidato. El importer entrega cada resultado al owner `postgres.LeadIngress` de V224. No escribe CRM ni concede contacto.

La integración PostgreSQL real demostró en una misma ejecución:

- un candidato Meta durable con `contact_eligibility=pending_policy`;
- un rechazo Meta durable con tenant/organization aun cuando `candidate=null`;
- raw append-only y un outbox por resultado;
- replay del mismo lote como dos duplicados y sin nuevas filas;
- suite completa de 47 paquetes Go, `go vet ./...` y `go build ./...`.

## Procedencia

- Contratos de `Lead` y `LeadgenForm` Meta: commit exacto `788f363d15b1269ab5efb7cd00fb5e3b133cd99b`.
- Contrato productor: `PYTHON-META-LEAD-RECONCILIATION-ADAPTER 0.1.0`.
- Owner durable: `GO-OMNICHANNEL-LEAD-INGRESS 0.1.0`.
- Promoción posterior: `GO-LEAD-CANDIDATE-PROMOTION 0.1.0`.
- Los cuatro archivos Go son `ADAPTED`; las tres fixtures son salidas locales `AUTHORED`. No se presenta código local como Meta-authored.

## Hashes materializados

| Archivo | SHA-256 |
|---|---|
| `internal/leadstream/meta_import.go` | `1a51c555ceb30c42e339dfcabd216e76ae1781b5bcd218b1e9077b680dda2f2f` |
| `internal/leadstream/meta_import_test.go` | `29e3cefea572e1c16b2ec0a5d4f09dc258481f0ed58198e6bcf6e98a8e08c936` |
| `internal/leadstream/testdata/meta-v226-valid/provider-response.json` | `fec8986d87cb71e9dfdffb40494ac8aa579c7a56a613605ec6978fe920a2d798` |
| `internal/leadstream/testdata/meta-v226-valid/lead-candidates.json` | `57a5fc4e5a4fce062b244f5c0937fb8e4cb1b43caed8dd00e4efbfe22ebab0c9` |
| `internal/leadstream/testdata/meta-v226-valid/RETRIEVAL_RECEIPT.json` | `dd37223aaf77a733b474496ed2f446ba92d7c91c9343360ebe44602e876c491e` |
| `internal/platform/postgres/meta_import_integration_test.go` | `062b40c83027aa6ead8687174f12294c14db604c9daf8209c35988b3acbc12c8` |
| `cmd/meta-lead-import/main.go` | `446cf69e93fb0817d7523152f69d5e689f477dc4b0fe67bcde7c4b5799e2ad15` |

Pack Markdown SHA-256: `6a760bcf50855ffddb3eb6b4b7503edb1c84203416823b0f5e9d15063bfbffaf`.

## Límites retenidos

- La prueba usó fixtures y PostgreSQL local, no una cuenta Meta live.
- Si el Store falla durante una lista ya validada, puede quedar un prefijo durable; el receipt de progreso y el replay idempotente permiten converger sin duplicar.
- Un lote V226 truncado se importa pero queda marcado como fuente incompleta; requiere continuación/reconciliación operacional.
- Falta scheduling, storage cifrado/retención target, webhook live, política V225 real, conversación, outcome y conversion postback.
