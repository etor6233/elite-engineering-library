# Project Blueprint Contract

Antes de elegir tecnología o materializar código, crear `PROJECT_BLUEPRINT.md` con esta estructura. El YAML vive dentro del Markdown para que continúe siendo legible y portable entre agentes.

```yaml
blueprint_version: "1.1"
project:
  id: ""
  outcome: ""
  maturity_target: prototype|internal|production|regulated
business:
  model: ""
  buyers: []
  users: []
  operators: []
  jurisdictions: []
  failure_costs: []
journeys:
  must_have: []
  later: []
  out_of_scope: []
domain:
  entities: []
  invariants: []
  sources_of_truth: []
  organizations_and_roles: []
quality:
  availability_slo: "UNKNOWN"
  latency_budget: "UNKNOWN"
  capacity: "UNKNOWN"
  rpo: "UNKNOWN"
  rto: "UNKNOWN"
  privacy_security: []
  compliance_and_data_classes: []
  cost_budget_and_unit_economics: "UNKNOWN"
interfaces:
  public_web: required|optional|none
  admin: required|optional|none
  customer: required|optional|none
  partner_factory_supplier: required|optional|none
  mobile: required|optional|none
  desktop: required|optional|none
  machine_device_firmware: required|optional|none
platform:
  topology: modular_monolith|services|edge_hybrid|UNKNOWN
  runtime_candidates: []
  tenancy_and_isolation: "UNKNOWN"
  environments: [local, ci, staging, production]
  compute_and_packaging: "UNKNOWN"
  networking_dns_tls_edge: "UNKNOWN"
  secrets_keys_pki: "UNKNOWN"
  identity_sessions_authorization: "UNKNOWN"
data:
  transactional_source_of_truth: "UNKNOWN"
  cache: required|optional|none|UNKNOWN
  search: required|optional|none|UNKNOWN
  object_storage: required|optional|none|UNKNOWN
  messaging_streaming_workflows: required|optional|none|UNKNOWN
  analytics_bi: required|optional|none|UNKNOWN
  retention_backup_pitr_dr: "UNKNOWN"
intelligence:
  ml_ai: required|optional|none|UNKNOWN
  rag_agents: required|optional|none|UNKNOWN
  gpu_acceleration: required|optional|none|UNKNOWN
  model_data_eval_governance: "UNKNOWN"
integrations:
  providers: []
  payments_marketplaces_ads: []
  webhooks_reconciliation: "UNKNOWN"
  lead_capture_to_channel_identity_binding: "UNKNOWN"
  outbound_delivery_receipts_and_reconciliation: "UNKNOWN"
engineering:
  source_control_and_branching: "UNKNOWN"
  build_dependency_packaging: "UNKNOWN"
  ci_cd_release_rollback: "UNKNOWN"
  observability_slo_incident: "UNKNOWN"
  supply_chain_sbom_provenance_signing_licenses: "UNKNOWN"
  testing: [unit, property, contract, integration, e2e, accessibility, performance, security, recovery]
  documentation_adrs_runbooks: "UNKNOWN"
delivery:
  team_skills: []
  existing_stack: []
  deployment_constraints: []
  budget_deadline: "UNKNOWN"
decisions:
  configurable_variants: []
  code_extensions_expected: []
  irreversible_unknowns: []
capability_closure:
  required: []
  optional: []
  not_applicable_with_reason: []
  open_gaps: []
```

## Gate

Un valor `UNKNOWN` no autoriza al agente a inventarlo. Si afecta una decisión irreversible, queda `BLOCKING`; si es reversible, se registra supuesto, fecha de revisión y rollback.

El agente transforma cada journey `must_have` en:

```text
actor → entrada → autorización → invariantes → writes → side effects
→ resultado visible → telemetría → fallos → aceptación
```

## Gate de infraestructura total

Antes de programar, cada superficie de `platform`, `data`, `intelligence`, `integrations` y `engineering` debe quedar en una de cuatro clases:

```text
REQUIRED       → pack/implementación/gates deben entrar al plan
OPTIONAL       → trigger de activación y costo de incorporación explícitos
NONE           → exclusión consciente sin dependencias residuales
UNKNOWN        → supuesto reversible o blocker, nunca omisión silenciosa
```

`not_applicable_with_reason` exige una razón verificable. Por ejemplo, un sistema web puede declarar firmware `NONE`; no puede omitir backups, identidad o supply chain si persiste datos y se desplegará en producción.
