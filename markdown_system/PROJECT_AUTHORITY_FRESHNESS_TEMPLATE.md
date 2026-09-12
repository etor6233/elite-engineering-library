# Project Authority Freshness Record Template

> Copiar como `PROJECT_AUTHORITY_FRESHNESS_RECORD.md`. Aplicar `AUTHORITY_FRESHNESS_AND_SELF_CORRECTION_CONTRACT.md`. No almacenar credenciales, PII ni contratos privados completos; usar referencias autorizadas.

## 1. Control

```yaml
record_version: "1.0"
project_id: ""
owner: ""
updated_at: ""
review_due_items: []
stale_blockers: []
unresolved_conflicts: []
last_corpus_wide_consistency_gate: ""
```

## 2. Authorities

| Authority ID | Capability/claim | Fuente oficial o owner | Versión/fecha/hash | Applicability | Verificado | Próxima revisión/trigger | Estado | Evidence |
|---|---|---|---|---|---|---|---|---|
|  |  |  |  |  |  |  | `REVIEW_DUE` |  |

## 3. Eventos de contradicción/autocorrección

```yaml
- correction_id: CORRECTION-YYYYMMDD-001
  detected_at: ""
  affected_claim: ""
  previous_statement: ""
  contradictory_evidence: []
  classification: factual|contract|code|dependency|business-decision
  blast_radius: ""
  corrected_statement: ""
  supersedes: ""
  files_or_decisions_changed: []
  hashes_sbom_notices_regenerated: []
  gates_rerun: []
  related_failure_ids: []
  status: OPEN
  owner: ""
  next_review: ""
```

Estados: `OPEN|CURRENT_VERIFIED|CURRENT_CONDITIONED|REVIEW_DUE|STALE_BLOCKED|CONFLICT_UNRESOLVED|SUPERSEDED|NOT_APPLICABLE_WITH_REASON`.

## 4. Consistency matrix

| Claim | Blueprint/spec | Source lock | Pack/code | Tests/evidence | Readiness | Resultado |
|---|---|---|---|---|---|---|
|  |  |  |  |  |  | `PENDING` |

## 5. Checklist

- [ ] hechos temporales revalidados al usarlos;
- [ ] toda fuente externa es primaria y tiene fecha/versión;
- [ ] applicability/target diferenciados de “lo más nuevo”;
- [ ] memoria del agente no se usó como autoridad final;
- [ ] contradicciones congelaron el claim antes de mutar producto;
- [ ] before/after y blast radius preservados;
- [ ] source locks/packs/hashes/SBOM/notices actualizados cuando correspondía;
- [ ] fallos derivados registrados;
- [ ] gates dependientes repetidos;
- [ ] corrections comunicadas explícitamente;
- [ ] próxima revisión o trigger definido.
