# Project Official Source Profile Record — Template

> Copiar como `PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md`. Este registro une las respuestas del usuario con un perfil materializado por `OFFICIAL-UPSTREAM-ACQUISITION-CORE`; no contiene secretos.

## 1. Control

```yaml
schema: elite-project-official-source-profile-record/v1
project_id: ""
owner: ""
status: DISCOVERY|AWAITING_USER|VALIDATED|ACQUIRED|BLOCKED
updated_at: ""
selected_profiles: []
profile_receipts: []
external_source_lock: PROJECT_EXTERNAL_SOURCE_LOCK.md
critical_unknowns: []
required_access_not_proven: []
production_blockers: []
```

## 2. Selección por capability

| Capability | Perfil | Requerido | Motivo | IDs seleccionados | IDs excluidos y razón | Owner | Estado |
|---|---|---|---|---|---|---|---|
| secure file ingestion local | `secure-file-ingestion-leaders` |  |  |  |  |  | `AWAITING_USER` |
| secure archive Linux | `secure-archive-linux` |  |  |  |  |  | `AWAITING_USER` |
| AWS document pipeline | `aws-secure-document-pipeline` |  | elegir como máximo un pipeline cloud |  |  |  | `AWAITING_USER` |
| Google document pipeline | `google-secure-document-pipeline` |  | elegir como máximo un pipeline cloud |  |  |  | `AWAITING_USER` |
| Microsoft document pipeline | `microsoft-secure-document-pipeline` |  | elegir como máximo un pipeline cloud |  |  |  | `AWAITING_USER` |
| document intelligence | `document-intelligence-leaders` |  |  |  |  |  | `AWAITING_USER` |
| commerce/communications | `commerce-communications-leaders` |  |  |  |  |  | `AWAITING_USER` |
| enterprise platform | `enterprise-platform-leaders` |  |  |  |  |  | `AWAITING_USER` |

## 3. Respuestas obligatorias

Por cada cadena de `required_user_inputs` del JSON materializado conservar:

| Perfil | Input requerido exacto | Respuesta o `NONE_WITH_REASON` | Fuente | Owner | Evidencia/probe | Estado |
|---|---|---|---|---|---|---|
|  |  |  | `user|contract|probe` |  |  | `UNKNOWN` |

Ninguna fila puede marcarse `PROVEN` por inferencia. Tokens, certificados, cookies, claves y documentos privados se entregan por el mecanismo seguro del entorno y aquí sólo se registra una referencia lógica.

Los tres perfiles cloud son alternativas de adquisición, no capas acumulativas. Seleccionar más de uno exige una decisión explícita de evaluación multi-provider, aislamiento de costos/datos y un plan de comparación; nunca se hace para “tener todo” sin caso aprobado.

## 4. Adquisición exacta

```yaml
profile_path: ""
profile_sha256: ""
lock_path: ""
lock_sha256: ""
validate_command: ""
validate_observed: ""
acquisition_command_redacted: ""
acquired_at: ""
aggregate_receipt: ""
aggregate_receipt_sha256: ""
source_receipts: []
license_notice_evidence: []
sbom_evidence: []
```

El comando permitido es `elite_sources/apply_source_profile.ps1`. Alterar manualmente el JSON materializado abre una revisión de source lock; no se permite añadir un ID móvil o desconocido.

## 5. Condiciones que permanecen abiertas

Copiar literalmente cada `production_blockers` del perfil y enlazar su cierre:

| Perfil | Blocker exacto | Cierre requerido | Evidencia | Owner | Estado |
|---|---|---|---|---|---|
|  |  |  |  |  | `OPEN` |

La adquisición completa sólo cambia `status` a `ACQUIRED`; nunca convierte automáticamente un blocker en `CLOSED` ni autoriza producción.

## 6. Gate

```yaml
all_selected_profiles_from_materialized_pack: false
all_required_inputs_answered: false
all_source_ids_in_exact_lock: false
all_selected_sources_have_receipts: false
all_licenses_notices_preserved: false
all_required_access_proven: false
all_production_blockers_closed_or_explicitly_blocking: false
ready_for_project_pack_plan: false
ready_for_production: false
```

`ready_for_project_pack_plan` puede ser `true` con blockers productivos abiertos sólo si esos blockers no afectan el primer plan autorizado y permanecen visibles. `ready_for_production` exige todos los blockers aplicables cerrados con evidencia real.
