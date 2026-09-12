# Project Dependency Update Record Template

> Copiar como `PROJECT_DEPENDENCY_UPDATE_RECORD.md`. No incluir tokens, credenciales, URLs privadas ni payloads sensibles. Aplicar `DEPENDENCY_UPDATE_CONTRACT.md` y enlazar fallos en `PROJECT_FAILURE_LESSONS.md`.

## 1. Control

```yaml
record_version: "1.0"
project_id: ""
owner: ""
updated_at: ""
baseline_artifact_digest: ""
baseline_sbom: ""
baseline_license_report: ""
open_security_updates: []
open_eol_updates: []
blocked_updates: []
active_candidate: ""
last_full_review: ""
next_review: ""
```

## 2. Inventario

| Component ID | Tipo/scope | Versión/digest/commit actual | Fuente | Licencia/notices | Manifest/lock | Consumers | EOL | Policy | Owner |
|---|---|---|---|---|---|---|---|---|---|
|  |  |  |  |  |  |  |  |  |  |

## 3. Evaluación de candidato

```yaml
update_id: UPDATE-YYYYMMDD-001
component_id: ""
trigger: security|release|eol|bug|compatibility|periodic|provider-contract
status: DISCOVERED
current_identity: ""
candidate_version: ""
candidate_commit_or_digest: ""
official_release: ""
official_changelog: ""
official_migration_guide: ""
published_or_observed_hashes: []
signature_or_provenance: ""
license_delta: ""
notice_delta: ""
dependency_graph_delta: ""
breaking_changes: []
new_defaults_or_privileges: []
consumers: []
blast_radius: ""
rollback_artifact: ""
rollback_data_config_protocol: ""
related_failures: []
advisory_ids: []
affected_ranges: []
observed_artifact_identity: ""
runtime_exposure: "UNKNOWN"
reachability: UNKNOWN
reachability_evidence: []
exploit_maturity: unknown
kev_status: unknown
upstream_fix_status: unknown
official_fixed_release: ""
official_patch_commit: ""
response_lane: ""
containment: []
local_patch_provenance: NONE
exception_owner: ""
exception_expiry: ""
incident_id: ""
owner: ""
```

### Respuesta a vulnerabilidad

| Comprobación | Resultado | Fuente/evidencia | Estado/acción |
|---|---|---|---|
| artefacto y consumers afectados |  |  | `UNKNOWN` |
| advisory/rangos oficiales |  |  | `PENDING` |
| reachability/exploitability |  |  | `UNKNOWN` |
| release o patch commit oficial |  |  | `UNKNOWN` |
| contención y blast radius |  |  | `PENDING` |
| lane seleccionada y provenance |  |  | `PENDING` |
| secretos/datos/tenants expuestos |  |  | `UNKNOWN` |
| incidente/forensics/rotación |  |  | `NONE_WITH_REASON` |

## 4. Gates baseline versus candidate

| Gate | Baseline | Candidate | Evidencia | Estado | Condición/acción |
|---|---|---|---|---|---|
| frozen install/clean build |  |  |  | `PENDING` |  |
| unit/property/negative |  |  |  | `PENDING` |  |
| contract/integration |  |  |  | `PENDING` |  |
| schema/mixed-version/rollback |  |  |  | `PENDING` |  |
| security/advisories/malware |  |  |  | `PENDING` |  |
| license/notices/SBOM/provenance |  |  |  | `PENDING` |  |
| API/ABI/protocol/config |  |  |  | `PENDING` |  |
| performance/cost |  |  |  | `PENDING` |  |
| deploy/canary/recovery |  |  |  | `PENDING` |  |

## 5. Decisión

```yaml
decision: BLOCKED
decided_by: ""
decided_at: ""
reason: ""
promoted_artifact_digest: ""
promoted_lock_hashes: []
promoted_sbom: ""
promoted_notices: []
canary_evidence: ""
rollback_evidence: ""
observation_window: ""
residual_risks: []
supersedes: ""
```

## 6. Checklist

- [ ] fuentes oficiales y revisión inmutable;
- [ ] licencia/notices/terms comparados;
- [ ] breaking changes, defaults, scripts y privilegios revisados;
- [ ] graph/lock/SBOM reproducibles;
- [ ] baseline y candidate ejecutados sobre gates equivalentes;
- [ ] fallos registrados, sin warnings/skips ocultos;
- [ ] generated code revisado y contract-tested;
- [ ] migración/mixed version/rollback o forward recovery probado;
- [ ] mismo digest promovido, sin rebuild ambiental;
- [ ] versión anterior y evidencia retenidas;
- [ ] readiness/authority maps actualizados si cambió una capability.
- [ ] una vulnerabilidad reabrió la admisión y bloqueó promotion hasta resolver su lane;
- [ ] un parche local se identifica como `ADAPTED_PATCH`/`AUTHORED_PATCH`, nunca como release upstream;
- [ ] excepción por falso positivo/no alcanzable tiene evidencia, owner y expiración;
- [ ] incidente, SBOM, consumers, locks, notices, runbooks y lecciones quedaron actualizados.
