# Project Failure Lessons Template

> Copiar como `PROJECT_FAILURE_LESSONS.md` al iniciar o continuar un proyecto. No incluir secretos, PII, documentos privados ni payloads completos. Este artefacto aplica `FAILURE_LEARNING_CONTRACT.md`.

## 1. Control

```yaml
ledger_version: "1.0"
project_id: ""
owner: ""
created_at: ""
updated_at: ""
open_high_or_critical: []
recurring_failures: []
last_clean_rebuild_evidence: ""
```

## 2. Entradas

Copiar una entrada por causa material:

```yaml
- failure_id: FAIL-YYYYMMDD-001
  fingerprint: ""
  related_failures: []
  first_seen: ""
  last_seen: ""
  recurrence_count: 1
  status: OPEN
  severity: HIGH
  classification: CODE
  phase: test
  capability: ""
  scope: ""
  command_redacted: ""
  toolchain: []
  error_signature_redacted: ""
  violated_invariant: ""
  blast_radius: ""
  reproduction: ""
  hypotheses_rejected: []
  root_cause: ""
  canonical_correction: ""
  regression_test: ""
  clean_rebuild_evidence: ""
  dependent_gates_rerun: []
  upstream_reference: ""
  residual_risk: ""
  owner: ""
  next_action: ""
```

Estados: `OPEN|DIAGNOSED|FIXED|REGRESSION_PROVEN|BLOCKED_EXTERNAL|UPSTREAM_OPEN|REJECTED_COMPONENT|ACCEPTED_RISK`.

## 3. Resumen operativo

| ID | Huella | Severidad | Estado | Recurrencias | Capability | Invariante | Próximo gate |
|---|---|---|---|---:|---|---|---|
|  |  |  |  |  |  |  |  |

## 4. Lecciones aplicadas al plan actual

| Lección previa | Decisión preventiva | Archivo/contrato/gate | Evidencia |
|---|---|---|---|
|  |  |  |  |

## 5. Condiciones upstream seleccionadas

| Fuente/revisión | Fallo o condición retenida | Impacto en el proyecto | Aislamiento/gate | Estado |
|---|---|---|---|---|
|  |  |  |  |  |

## 6. Cierre de ciclo

- [ ] cada fallo observado tiene ID y huella;
- [ ] cada retry conserva su resultado y presupuesto;
- [ ] ningún warning material, skip o test deshabilitado quedó sin clasificación;
- [ ] cada corrección volvió a la fuente canónica;
- [ ] cada fix tiene regresión o permanece honestamente abierto/bloqueado;
- [ ] se reconstruyó desde un destino vacío;
- [ ] se repitieron los gates dependientes previamente verdes;
- [ ] las lecciones reutilizables se propusieron al ledger de la biblioteca sin datos privados;
- [ ] los riesgos aceptados tienen owner, vencimiento y control compensatorio;
- [ ] el estado de readiness refleja los fallos abiertos reales.
