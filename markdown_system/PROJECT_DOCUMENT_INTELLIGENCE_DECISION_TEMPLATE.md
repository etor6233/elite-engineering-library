# Project Document Intelligence Decision — Template

Copiar al proyecto como `PROJECT_DOCUMENT_INTELLIGENCE_DECISION.md`. El agente completa una fila/registro por clase y variante; no resume varias clases con “similar”.

## 1. Estado global

```yaml
schema: elite-project-document-intelligence/v1
project: ""
owner: ""
security_owner: ""
data_owner: ""
review_operations_owner: ""
created_at: ""
updated_at: ""
status: DISCOVERY|CORPUS_BLOCKED|ACCESS_BLOCKED|EVALUATING|REVIEW_ONLY|LIMITED_AUTOMATION|READY_FOR_AUTOMATIC_STORAGE
critical_unknowns: []
required_access_not_proven: []
source_lock_path: PROJECT_EXTERNAL_SOURCE_LOCK.md
evaluation_evidence_root: ""
```

## 1.1 Contrato completo de ingesta de archivos

Este bloque se completa aunque todavía no exista OCR. Un `unknown`, vacío o sin evidencia conserva el proyecto en `DISCOVERY`/`POLICY_BLOCKED`; el agente no elige defaults silenciosos.

```yaml
file_ingestion:
  scope_owner: ""
  security_owner: ""
  operations_owner: ""
  incident_owner: ""
  ingress_channels:
    - channel: portal|api|email|sftp|object_event|scanner|mobile|marketplace|erp|other
      status: REQUIRED|OPTIONAL|NONE_WITH_REASON|BLOCKED
      trust_boundary: ""
      tenant_identity: ""
      actor_identity: ""
      authentication: ""
      transport_security: ""
      checksum_contract: ""
      partial_upload_cleanup: ""
      max_bytes: 0
      max_files_per_request: 0
      rate_limit: ""
      timeout: ""
      retry_and_idempotency: ""
      evidence_path: ""
  immutable_original:
    location: ""
    encryption_and_key_owner: ""
    tenant_isolation: ""
    object_lock_or_versioning: ""
    sha256_before_processing: false
    received_at_source: ""
    retention: ""
    legal_hold: ""
    deletion_and_tombstone: ""
  lifecycle_states:
    - RECEIVING
    - QUARANTINED
    - TYPE_REJECTED
    - MALWARE_REJECTED
    - CLEAN
    - EXTRACTING
    - EXTRACTION_REVIEW
    - VALIDATING
    - REVIEW_REQUIRED
    - APPROVED
    - PERSISTED
    - SUPERSEDED
    - DELETED
  content_type_gate:
    claimed_type_policy: ""
    magika_source_id: google-magika-cli-1.1.0
    accepted_labels: []
    review_labels: []
    rejected_labels: []
    minimum_confidence: ""
    extension_mismatch_action: REJECT|QUARANTINE|REVIEW
    unknown_or_polyglot_action: REJECT|QUARANTINE|REVIEW
    active_content_and_macros_policy: ""
    encrypted_document_policy: ""
  malware_gate:
    primary_source_id: cisco-clamav-1.5.4
    cloud_source_id_optional: aws-guardduty-malware-protection-s3
    signature_origin_and_authentication: ""
    signature_max_age: ""
    stale_or_update_failure_action: FAIL_CLOSED|QUARANTINE
    scanner_isolation: ""
    max_scan_bytes: 0
    max_scan_time: ""
    max_recursion: 0
    max_expanded_bytes: 0
    clean_result_ttl: ""
    quarantine_location: ""
    release_approvers: []
    rescan_and_recall_policy: ""
  yara_gate:
    source_id: virustotal-yara-x-1.20.0
    enabled: false
    ruleset_repository_or_archive: ""
    ruleset_version_and_sha256: ""
    rule_authors_and_reviewers: []
    rollout: shadow|canary|enforce
    false_positive_response: ""
    rollback: ""
  archive_gate:
    accepted_formats: []
    extraction_default: FORBIDDEN|LINUX_ONLY_APPROVED
    linux_source_id: google-safearchive-f7ce9d7
    encrypted_archive_action: REJECT|QUARANTINE|REVIEW
    nested_archive_action: REJECT|BOUNDED
    max_compressed_bytes: 0
    max_expanded_bytes: 0
    max_expansion_ratio: 0
    max_entries: 0
    max_depth: 0
    max_processing_time: ""
    path_traversal_action: REJECT_AND_INCIDENT
    symlink_hardlink_special_file_action: REJECT_AND_INCIDENT
    duplicate_or_unicode_collision_action: REJECT
    extraction_sandbox: ""
    cleanup_evidence: ""
  sensitive_data_gate:
    presidio_source_id: presidio-2.2.364
    recognizers_and_languages: []
    purpose_and_legal_basis: ""
    redact_before_logs_or_egress: false
    false_positive_and_negative_evaluation: ""
    reviewer_roles: []
  provenance_contract:
    required_fields:
      - file_id
      - tenant_id
      - source_channel
      - source_reference
      - actor_reference
      - received_at
      - original_name
      - original_bytes
      - original_sha256
      - claimed_content_type
      - detected_content_type
      - scanner_engine_version
      - scanner_signature_version
      - yara_ruleset_version
      - parser_or_model_version
      - schema_version
      - decision
      - evidence_reference
    append_only_audit_location: ""
  reliability:
    idempotency_key: ""
    duplicate_policy: ""
    queue_and_backpressure: ""
    retry_classes: []
    poison_and_dlq: ""
    replay_owner: ""
    reconciliation: ""
    disaster_recovery: ""
    observability_without_sensitive_payloads: ""
  admission_evidence:
    official_sources_acquired_with_receipts: false
    benign_and_malicious_fixture_gate: false
    stale_signature_failure_gate: false
    type_spoof_and_polyglot_gate: false
    hostile_archive_gate: false
    duplicate_and_partial_upload_gate: false
    tenant_isolation_gate: false
    load_and_backpressure_gate: false
    quarantine_release_and_recall_gate: false
    replay_reconciliation_and_restore_gate: false
    evidence_root: ""
```

## 2. Inventario documental

Registrar al menos estas clases como `REQUIRED`, `OPTIONAL`, `NONE_WITH_REASON` o `BLOCKED`:

| Clase | Decisión | Variantes/emisores | Idioma/país | Volumen | Criticidad | Owner |
|---|---|---|---|---:|---|---|
| invoice / supplier invoice | | | | | | |
| receipt / ticket | | | | | | |
| proforma invoice | | | | | | |
| commercial invoice | | | | | | |
| packing list | | | | | | |
| purchase order | | | | | | |
| purchase order confirmation | | | | | | |
| bill of lading / airway bill | | | | | | |
| delivery note / remito | | | | | | |
| certificate of origin | | | | | | |
| quality/inspection certificate | | | | | | |
| customs declaration/clearance | | | | | | |
| insurance/freight document | | | | | | |
| price list/quotation | | | | | | |
| product specification/catalog | | | | | | |
| contract/addendum | | | | | | |
| bank/payment statement | | | | | | |
| warranty/service/claim document | | | | | | |
| identity or regulatory document | | | | | | |
| multi-document package | | | | | | |
| other named class | | | | | | |

## 3. Registro obligatorio por clase/variante

```yaml
document_class: ""
variant: ""
decision: REQUIRED|OPTIONAL|NONE_WITH_REASON|BLOCKED
business_journey: ""
emitters: []
countries: []
languages: []
input:
  mime_types: []
  digital_pdf: true
  scans: true
  photos: false
  handwriting: false
  multipage: true
  multi_document_bundle: false
  dpi_observed: []
  max_pages: 0
  max_bytes: 0
  expected_daily_volume: 0
  peak_concurrency: 0
schema_version: ""
fields:
  - name: ""
    type: ""
    required: true
    cardinality: one|many
    criticality: critical|high|normal
    normalization: ""
    cross_checks: []
    storage_policy: BLOCK|REVIEW_ONLY|AUTO_IF_PROVEN
relationships:
  - type: "PO↔invoice|invoice↔packing-list|packing-list↔BOL|other"
    identifiers: []
    invariants: []
data_policy:
  pii: []
  financial: []
  secrets: []
  retention: ""
  residency: ""
  encryption: ""
  reviewer_roles: []
corpus:
  authorized: false
  total_documents: 0
  ground_truth_documents: 0
  train_documents: 0
  validation_documents: 0
  test_documents: 0
  hard_cases: []
  ground_truth_owner: ""
  evidence_path: ""
official_lane:
  provider: AZURE|GOOGLE|AWS|IBM|NVIDIA|MICROSOFT|SAP|ORACLE|MULTI_PROVIDER
  source_ids: []
  prebuilt_processor: ""
  custom_model_required: false
  classifier_required: false
  splitter_required: false
  review_platform: ""
access:
  environment: sandbox|test|production
  identity_reference: ""
  required_scopes: []
  region: ""
  probe_status: NOT_PROVIDED|FAILED|PROVEN
  probe_evidence: ""
model:
  id: ""
  version: ""
  schema_version: ""
  trained_at: ""
  training_evidence: ""
evaluation:
  exact_match_by_field: {}
  normalized_match_by_field: {}
  missing_by_field: {}
  false_values_by_field: {}
  out_of_distribution_tests: []
  duplicate_tests: []
  corrupt_partial_tests: []
  arithmetic_reference_tests: []
  latency_p50_ms: 0
  latency_p95_ms: 0
  cost_per_document: 0
  evidence_path: ""
review:
  queue: ""
  assignment_policy: ""
  correction_audit: ""
  escalation_sla: ""
  replay_after_correction: ""
storage:
  original_immutable_location: ""
  original_hash_algorithm: SHA-256
  provenance_fields: []
  idempotency_key: ""
  transaction_boundary: ""
  extracted_state: PROVISIONAL|REVIEWED|APPROVED
rollback:
  model: ""
  schema: ""
  data_correction: ""
  evidence: ""
decision_record:
  status: BLOCK|REVIEW_ONLY|LIMITED_AUTOMATION|READY_FOR_AUTOMATIC_STORAGE
  business_owner_approval: ""
  data_owner_approval: ""
  security_owner_approval: ""
  decided_at: ""
```

## 4. Comparación de lanes oficiales

Cuando el documento/campo sea crítico o no exista processor prebuilt específico, conservar una fila por candidato:

| Clase/campo | Source ID/modelo oficial | Corpus test | Exact | Falso | Missing | p95 | Costo | Condición | Decisión |
|---|---|---:|---:|---:|---:|---:|---:|---|---|
| | | | | | | | | | |

No elegir proveedor por marca, benchmark genérico o confianza reportada sin ground truth. La decisión usa el mismo corpus holdout y el costo/riesgo del proyecto.

## 5. Gate final

```yaml
all_required_classes_classified: false
all_ingress_channels_classified: false
all_file_lifecycle_states_and_owners_proven: false
all_original_hash_quarantine_and_retention_controls_proven: false
all_type_malware_rule_and_archive_gates_proven: false
all_signature_and_ruleset_staleness_fail_closed: false
all_duplicate_partial_poison_replay_and_reconciliation_paths_proven: false
all_file_security_admission_evidence_linked: false
all_required_fields_have_policy: false
all_required_sources_exactly_locked: false
all_required_access_proven: false
all_required_corpus_authorized: false
all_ground_truth_approved: false
all_models_and_schemas_versioned: false
all_critical_fields_zero_false_values_under_declared_policy: false
all_ambiguous_and_ood_inputs_fail_to_review: false
all_storage_idempotency_proven: false
all_review_and_correction_paths_proven: false
all_model_schema_data_rollbacks_proven: false
final_status: NOT_READY|REVIEW_ONLY_ACCEPTED|READY_FOR_AUTOMATIC_STORAGE
```

El agente no cambia `final_status` por haber terminado de llenar texto. Cada `true` enlaza evidencia reproducible.
