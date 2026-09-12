# Project Start Readiness Validator

## 1. Metadata

```yaml
pack_id: "PROJECT-START-READINESS-VALIDATOR"
pack_version: "0.7.0"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: REUSABLE_PACK
claim: "Materializa un asesor guiado y un gate ejecutable fail-closed: explica una ronda por vez desde un catálogo oficial-source-routed e impide READY_TO_BUILD sin ocho rondas, 48 capabilities y journeys release-bound que conecten persona/rol, interfaz, autorización, contratos, dominio, datos/efecto externo, auditoría, respuesta, ayuda, capacitación, soporte, actualización, pruebas y recuperación demostrados."
stacks: ["CPython 3.14.4 standard library", "portable JSON"]
compatible_with: ["PROJECT-START-READINESS-GATE 1.x", "TOTAL-SYSTEM-CAPABILITY-CONTRACT 1.x", "MARKDOWN-COMPOSITOR-CORE 0.2.x", "DEBEZIUM-POSTGRES-INBOX-CONSUMER 0.2.x"]
incompatible_with: ["credentials in readiness records", "absolute or traversal evidence paths", "moving source references", "project implementation before READY_TO_BUILD or explicit library preparation", "isolated frontend or backend capabilities", "unversioned help, training or support", "preconfigured or acknowledged consumer templates"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://github.com/github/spec-kit/tree/9118ed15a0ba65053469a94c560ea5d233f75884", "https://github.com/aws-samples/sample-well-architected-skills-and-steering/tree/4b58ba02670abcd67458eeafb44fb6118b6af2ef", "https://developers.openai.com/api/docs/guides/latest-model", "https://sre.google/sre-book/reliable-product-launches/", "https://www.nasa.gov/reference/system-engineering-handbook-appendix/", "https://learn.microsoft.com/en-us/dynamics365/business-central/hr-how-manage-absence", "https://learn.microsoft.com/en-us/dynamics365/business-central/ui-post-sales", "https://cloud.google.com/document-ai/docs/processors-list", "https://docs.aws.amazon.com/textract/latest/dg/limits-document.html", "https://www.arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG.pdf"]
verified_at: "2026-09-11"
```

## 2. Applicability

V283 / 0.6.1: los journeys seleccionados exclusivamente CLI/AUTOMATION ya no exigen una web/API ficticia. Deben enlazar RUNTIME, CONTRACTS, DOMAIN-MODULES y AUTHORIZATION requeridos; todos los eslabones de evidencia y versiones siguen siendo obligatorios. Cada superficie web/portal/mobile/desktop/API también exige su capability correspondiente enlazada. Regresión roja: dos tests; verde y reconstrucción: 70 tests, incluidos catorce subcasos de evidencia ausente. No habilita readiness local incompleto. Detalle: reconstruction_evidence/READINESS_CLI_CONTRACT_REPAIR_V283.md.

Use at the first project conversation and as a delta gate before extending an existing critical subsystem. It materializes only readiness tooling, not product code. The agent renders and explains rounds A–H conversationally; the validator prevents a completed-looking Markdown from opening implementation without the exact advisory prompt, machine-readable decisions and evidence.

This pack is entirely `AUTHORED` Elite control code. Official OpenAI, Google, NASA, Microsoft, AWS, ARCA, GitHub and AWS-samples sources govern the catalog or workflow; no source code or guarantee is copied or attributed here. Reject use if the project cannot retain non-secret evidence locally or if the owner refuses to classify all 48 surfaces.

## 3. Architecture contract

The advisory catalog contains exactly one source-routed topic per round. `render_project_advisory.py` validates duplicate keys, exact schema, official HTTPS hosts, topic/source identities, explicit non-assumed examples and A–H completeness, then atomically creates one deterministic project-root prompt. The agent presents it conversationally; the answer remains separate evidence. The project copies the disabled JSON template to `PROJECT_READINESS_GATE.json` and updates it with the human record. The validator rejects symlinks, unsafe evidence, secrets, capability drift, moving sources and any missing, modified or cross-round advisory prompt. After PASS, the consumer renderer transfers only the selected class/variant's ten exact mapping fields into a pristine closed template; it cannot manufacture operational evidence or enable deployment.

`READY_TO_BUILD` requires exact A–H completion with eight catalog-matching prompts, empty blocker arrays, a named P1 slice, essential/delivery capabilities, dependency implications, required integration probes or explicit later deferral, and a unique document class/variant inventory when documents are required. Every required capability must appear in at least one closed `connected_journeys` record; a P1 slice must select existing journey IDs. Each journey has one release/help/training/support version and evidence for persona/role, interface, authorization, API contract, domain and data/external effect, audit/observability, understandable response/errors, contextual help, training, support, change impact, rollback/recovery, E2E and negative tests. Unknown keys, incomplete links, isolated screens/backend effects, non-required capability references and version drift fail closed. Automatic storage additionally requires a class-specific approved decision and persistence mode. The mapping object uses the exact Debezium/CEL consumer field names; expression bytes must match their SHA, domain must equal the class, IDs/semantic version/output-key CSV/schema hashes must satisfy the same lexical contract, and corpus/evaluation/approval evidence must resolve. Transactional direct and outbox/CDC modes require their exact architecture capabilities. Exact sources/licenses, dependency/authority/operations proof and pack-plan collision/rollback evidence remain mandatory. Reports and prompts are atomic and never overwrite an existing artifact.

Semantic blocks and rejected CLI inputs exit 2; input/I/O diagnostics use PROJECT_READINESS_BLOCKED on stderr without fabricating a report. Argument parsing remains argparse exit 2. A PASS authorizes materializing the first planned slice, not production. Rollback removes only the newly generated report/tool directory; readiness history remains immutable project evidence.

## 4. Exact file manifest

```text
CREATE project_readiness_gate/project-readiness.template.json
CREATE project_readiness_gate/validate_project_readiness.py
CREATE project_readiness_gate/render_consumer_profile.py
CREATE project_readiness_gate/project-advisory-catalog.json
CREATE project_readiness_gate/render_project_advisory.py
CREATE project_readiness_gate/test_validate_project_readiness.py
CREATE project_readiness_gate/test_render_project_advisory.py
CREATE project_readiness_gate/README.md
CREATE project_readiness_gate/validate_library_readiness.py
CREATE project_readiness_gate/test_validate_library_readiness.py
CREATE project_readiness_gate/library-readiness.template.json
```

## 5. Materialization blocks

### FILE: `project_readiness_gate/project-readiness.template.json`
```yaml
block_id: "PROJECT-START-READINESS-VALIDATOR:template:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite project readiness and total capability contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "09cbbb9fdeefcbba7f8c09a1da9f558a0b81b834e3745bba2f1691aa160dd79c"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-project-readiness-gate/v2",
  "project": {
    "id": "",
    "owner": "",
    "status": "DISCOVERY",
    "platform_mode": "BLOCK",
    "risk_tier": "CRITICAL",
    "data_classification": "RESTRICTED",
    "updated_at": ""
  },
  "rounds": {
    "A": {"status": "PENDING", "evidence": [], "advisory_prompt_ref": ""},
    "B": {"status": "PENDING", "evidence": [], "advisory_prompt_ref": ""},
    "C": {"status": "PENDING", "evidence": [], "advisory_prompt_ref": ""},
    "D": {"status": "PENDING", "evidence": [], "advisory_prompt_ref": ""},
    "E": {"status": "PENDING", "evidence": [], "advisory_prompt_ref": ""},
    "F": {"status": "PENDING", "evidence": [], "advisory_prompt_ref": ""},
    "G": {"status": "PENDING", "evidence": [], "advisory_prompt_ref": ""},
    "H": {"status": "PENDING", "evidence": [], "advisory_prompt_ref": ""}
  },
  "blockers": {
    "critical_unknowns": ["PROJECT_OUTCOME_NOT_PROVEN"],
    "required_access_not_proven": [],
    "open_high_or_critical_failures": [],
    "open_dependency_blockers": [],
    "stale_or_conflicted_authorities": []
  },
  "first_vertical_slice": {
    "id": "",
    "journeys": [],
    "connected_journey_ids": [],
    "invariants": [],
    "acceptance_test_refs": [],
    "rollback_recovery_ref": ""
  },
  "capabilities": [
    {"id":"PRD-INTAKE","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"ARCH-DOMAIN","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"REPO-SCM","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"CONTRACTS","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"RUNTIME","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"WEB-PUBLIC","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"WEB-PORTALS","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"MOBILE","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"DESKTOP","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"EMBEDDED-IOT","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"API-BACKEND","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"DOMAIN-MODULES","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"IDENTITY","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"AUTHORIZATION","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"SECRETS-PKI","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"TX-DATABASE","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"CACHE","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"SEARCH","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"OBJECT-STORAGE","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"OUTBOX-INBOX","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"JOBS-WORKFLOWS","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"BROKER-STREAMING","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"INTEGRATIONS","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"PAYMENTS","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"MARKETPLACES","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"ADS-ATTRIBUTION","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"NOTIFICATIONS","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"DATA-INGEST","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"ANALYTICS-BI","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"ML-AI","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"RAG-AGENTS","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"GPU-ACCEL","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"NETWORK-EDGE","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"CONTAINERS","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"ORCHESTRATION","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"IAC-CLOUD","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"CI","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"CD-RELEASE","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"SUPPLY-CHAIN","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"OBSERVABILITY","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"SLO-INCIDENT","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"PERFORMANCE","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"SECURITY-APPSEC","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"PRIVACY-COMPLIANCE","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"BACKUP-DR","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"COST-FINOPS","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"TEST-PLATFORM","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]},
    {"id":"DOCS-OPS","classification":"BLOCKED","owner":"","reason_or_trigger":"","implementation_refs":[],"test_refs":[],"rollback_recovery_ref":"","blocking_conditions":[]}
  ],
  "connected_journeys": [],
  "documents": {
    "required": false,
    "decision_record": "",
    "corpus_ground_truth_status": "NOT_APPLICABLE",
    "automatic_storage_requested": false,
    "field_evaluation_evidence": [],
    "security_evidence": [],
    "classes": []
  },
  "integrations": [],
  "sources": {
    "lock_path": "PROJECT_EXTERNAL_SOURCE_LOCK.md",
    "status": "BLOCKED",
    "exact_revisions": false,
    "licenses_notices_resolved": false,
    "moving_references": [],
    "product_provenance_modes": []
  },
  "dependencies": {"status":"BLOCKED","evidence":[],"eol_policy":false,"locks":false,"rollback":false},
  "authorities": {"status":"BLOCKED","evidence":[],"conflicts":[]},
  "operations": {
    "security_privacy": {"status":"BLOCKED","evidence":[]},
    "slo_performance_capacity_cost": {"status":"BLOCKED","evidence":[]},
    "backup_restore_dr": {"status":"BLOCKED","evidence":[]},
    "deployment_canary_rollback": {"status":"BLOCKED","evidence":[]},
    "monitoring_incident_response": {"status":"BLOCKED","evidence":[]}
  },
  "pack_plan": {"path":"PROJECT_PACK_PLAN.md","status":"BLOCKED","collision_check":false,"rollback_defined":false,"evidence":[]}
}
````

### FILE: `project_readiness_gate/validate_project_readiness.py`
```yaml
block_id: "PROJECT-START-READINESS-VALIDATOR:validator:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite fail-closed semantic validator"
license: "LicenseRef-Workspace-Owner"
sha256: "d3b3c69b78cde3a3857c0713d1ee30c8ed217206aa3637fd38cf96efacde312b"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
from datetime import datetime
import hashlib
import json
import os
from pathlib import Path
import re
import sys
import tempfile
from typing import Any, Iterable

from render_project_advisory import AdvisoryError, validate_prompt_reference

SCHEMA = "elite-project-readiness-gate/v2"
CAPABILITIES = (
    "PRD-INTAKE", "ARCH-DOMAIN", "REPO-SCM", "CONTRACTS", "RUNTIME", "WEB-PUBLIC",
    "WEB-PORTALS", "MOBILE", "DESKTOP", "EMBEDDED-IOT", "API-BACKEND", "DOMAIN-MODULES",
    "IDENTITY", "AUTHORIZATION", "SECRETS-PKI", "TX-DATABASE", "CACHE", "SEARCH",
    "OBJECT-STORAGE", "OUTBOX-INBOX", "JOBS-WORKFLOWS", "BROKER-STREAMING", "INTEGRATIONS",
    "PAYMENTS", "MARKETPLACES", "ADS-ATTRIBUTION", "NOTIFICATIONS", "DATA-INGEST",
    "ANALYTICS-BI", "ML-AI", "RAG-AGENTS", "GPU-ACCEL", "NETWORK-EDGE", "CONTAINERS",
    "ORCHESTRATION", "IAC-CLOUD", "CI", "CD-RELEASE", "SUPPLY-CHAIN", "OBSERVABILITY",
    "SLO-INCIDENT", "PERFORMANCE", "SECURITY-APPSEC", "PRIVACY-COMPLIANCE", "BACKUP-DR",
    "COST-FINOPS", "TEST-PLATFORM", "DOCS-OPS",
)
REQUIRED_FILES = (
    ".specify/memory/constitution.md", "PROJECT_READINESS_RECORD.md", "PROJECT_FAILURE_LESSONS.md",
    "PROJECT_DEPENDENCY_UPDATE_RECORD.md", "PROJECT_AUTHORITY_FRESHNESS_RECORD.md",
    "PROJECT_VULNERABILITY_MONITORING_RECORD.md", "PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md",
    "PROJECT_BLUEPRINT.md", "PROJECT_AUTHORITY_MAP.md", "PROJECT_EXTERNAL_SOURCE_LOCK.md",
    "PROJECT_PACK_PLAN.md",
)
PROHIBITED_SECRET_KEYS = re.compile(r"(^|_)(password|token|secret|private_key|client_secret|api_key)($|_)", re.I)
CANONICAL_DOCUMENT_ID = re.compile(r"^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$")
SHA256_HEX = re.compile(r"^[0-9a-f]{64}$")
MAPPING_ID = re.compile(r"^[a-z][a-z0-9-]{2,63}$")
MAPPING_VERSION = re.compile(r"^[1-9][0-9]*\.[0-9]+\.[0-9]+$")
SCHEMA_ID = re.compile(r"^[a-z][a-z0-9._-]{2,127}$")
OUTPUT_KEY = re.compile(r"^[a-z][a-z0-9_]{0,62}$")
CONNECTED_JOURNEY_ID = re.compile(r"^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$")
SEMANTIC_VERSION = re.compile(r"^[1-9][0-9]*\.[0-9]+\.[0-9]+$")
CONNECTED_JOURNEY_SURFACES = {
    "WEB_PUBLIC", "WEB_PORTAL", "MOBILE", "DESKTOP", "API", "CLI",
    "AUTOMATION", "PHYSICAL_ASSISTED",
}
DELIVERY_CAPABILITY_BY_SURFACE = {
    "WEB_PUBLIC": "WEB-PUBLIC", "WEB_PORTAL": "WEB-PORTALS",
    "MOBILE": "MOBILE", "DESKTOP": "DESKTOP", "API": "API-BACKEND",
}
CONNECTED_JOURNEY_REFERENCE_FIELDS = (
    "frontend_or_interface_refs", "authorization_refs", "api_contract_refs",
    "domain_effect_refs", "data_or_external_effect_refs", "audit_observability_refs",
    "response_error_refs", "contextual_help_refs", "training_refs", "support_refs",
    "update_impact_refs", "rollback_recovery_refs", "e2e_test_refs", "negative_test_refs",
)
CONNECTED_JOURNEY_KEYS = {
    "id", "status", "release_version", "help_version", "training_version",
    "support_version", "personas", "roles", "surfaces", "capability_ids",
    *CONNECTED_JOURNEY_REFERENCE_FIELDS,
}
DOCUMENT_DECISIONS = {
    "CORPUS_BLOCKED", "ACCESS_BLOCKED", "POLICY_BLOCKED", "EVALUATION_REQUIRED",
    "REVIEW_ONLY", "LIMITED_AUTOMATION", "READY_FOR_AUTOMATIC_STORAGE",
}
MAPPING_FIELDS = {
    "mappingId", "mappingDomainType", "mappingVersion", "mappingExpression",
    "mappingExpressionSha256", "mappingExactOutputKeysCsv", "mappingInputSchemaId",
    "mappingInputSchemaSha256", "mappingOutputSchemaId", "mappingOutputSchemaSha256",
}


def nonempty(value: Any) -> bool:
    return isinstance(value, str) and bool(value.strip())


def string_list(value: Any) -> bool:
    return isinstance(value, list) and all(nonempty(item) for item in value)


def resolve_evidence(project_root: Path, value: str, label: str, errors: list[str]) -> None:
    if not nonempty(value):
        errors.append(f"{label}: empty evidence path")
        return
    path_text = value.split("#", 1)[0]
    candidate = Path(path_text)
    if candidate.is_absolute() or ".." in candidate.parts:
        errors.append(f"{label}: evidence path must be safe and project-relative: {value}")
        return
    try:
        resolved = (project_root / candidate).resolve(strict=True)
    except (FileNotFoundError, OSError):
        errors.append(f"{label}: evidence file missing: {value}")
        return
    try:
        resolved.relative_to(project_root)
    except ValueError:
        errors.append(f"{label}: evidence escaped project root: {value}")
        return
    if candidate.is_symlink() or not resolved.is_file() or resolved.stat().st_size == 0:
        errors.append(f"{label}: evidence must be a non-empty regular non-symlink file: {value}")


def evidence_list(project_root: Path, value: Any, label: str, errors: list[str], required: bool = True) -> None:
    if not string_list(value):
        errors.append(f"{label}: evidence must be an array of non-empty paths")
        return
    if required and not value:
        errors.append(f"{label}: at least one evidence path is required")
    for index, path in enumerate(value):
        resolve_evidence(project_root, path, f"{label}[{index}]", errors)


def validate_connected_journeys(value: Any, classification_by_id: dict[str, Any],
                                first_slice: dict[str, Any], project_root: Path,
                                errors: list[str]) -> None:
    if not isinstance(value, list) or not value:
        errors.append("connected_journeys must be a non-empty array")
        return
    journey_ids: list[str] = []
    covered_capabilities: set[str] = set()
    for index, item in enumerate(value):
        label = f"connected_journeys[{index}]"
        if not isinstance(item, dict):
            errors.append(f"{label} must be an object")
            continue
        if set(item) != CONNECTED_JOURNEY_KEYS:
            errors.append(f"{label} must contain exactly the connected journey fields")
        journey_id = item.get("id")
        if (not isinstance(journey_id, str) or not CONNECTED_JOURNEY_ID.fullmatch(journey_id)
                or not 3 <= len(journey_id) <= 128):
            errors.append(f"{label}.id must be a canonical lowercase hyphenated ID")
        else:
            journey_ids.append(journey_id)
        if item.get("status") != "PROVEN":
            errors.append(f"{label}.status must be PROVEN")
        versions = [item.get(field) for field in (
            "release_version", "help_version", "training_version", "support_version"
        )]
        if any(not isinstance(version, str) or not SEMANTIC_VERSION.fullmatch(version)
               for version in versions):
            errors.append(f"{label} versions must be explicit semantic versions")
        elif len(set(versions)) != 1:
            errors.append(f"{label} help, training and support versions must equal release_version")
        for field in ("personas", "roles"):
            entries = item.get(field)
            if not string_list(entries) or not entries or len(entries) != len(set(entries)):
                errors.append(f"{label}.{field} must be a non-empty unique string array")
        surfaces = item.get("surfaces")
        if (not string_list(surfaces) or not surfaces or len(surfaces) != len(set(surfaces))
                or not set(surfaces) <= CONNECTED_JOURNEY_SURFACES):
            errors.append(f"{label}.surfaces must be a non-empty unique supported surface array")
        capability_ids = item.get("capability_ids")
        if string_list(surfaces) and string_list(capability_ids):
            needed = {DELIVERY_CAPABILITY_BY_SURFACE[s] for s in surfaces
                      if s in DELIVERY_CAPABILITY_BY_SURFACE}
            if set(surfaces) & {"CLI", "AUTOMATION"}:
                needed.update({"RUNTIME", "CONTRACTS", "DOMAIN-MODULES", "AUTHORIZATION"})
            for capability_id in sorted(needed):
                if classification_by_id.get(capability_id) != "REQUIRED" or capability_id not in capability_ids:
                    errors.append(f"{label}.surfaces requires linked REQUIRED capability {capability_id}")
        if not string_list(capability_ids) or not capability_ids or len(capability_ids) != len(set(capability_ids)):
            errors.append(f"{label}.capability_ids must be a non-empty unique string array")
        else:
            for capability_id in capability_ids:
                if capability_id not in CAPABILITIES:
                    errors.append(f"{label}.capability_ids contains unknown capability {capability_id}")
                elif classification_by_id.get(capability_id) != "REQUIRED":
                    errors.append(f"{label}.capability_ids requires capabilities.{capability_id}=REQUIRED")
                else:
                    covered_capabilities.add(capability_id)
        for field in CONNECTED_JOURNEY_REFERENCE_FIELDS:
            evidence_list(project_root, item.get(field), f"{label}.{field}", errors)
    duplicates = sorted({item for item in journey_ids if journey_ids.count(item) > 1})
    if duplicates:
        errors.append(f"connected_journeys IDs must be unique; duplicates={duplicates}")
    required_capabilities = {
        capability_id for capability_id, classification in classification_by_id.items()
        if classification == "REQUIRED"
    }
    uncovered = sorted(required_capabilities - covered_capabilities)
    if uncovered:
        errors.append(f"connected_journeys must cover every REQUIRED capability; missing={uncovered}")
    selected = first_slice.get("connected_journey_ids") if isinstance(first_slice, dict) else None
    if not string_list(selected) or not selected or len(selected) != len(set(selected)):
        errors.append("first_vertical_slice.connected_journey_ids must be a non-empty unique string array")
    else:
        unknown = sorted(set(selected) - set(journey_ids))
        if unknown:
            errors.append(f"first_vertical_slice.connected_journey_ids contains unknown journeys {unknown}")


def validate_document_classes(documents: dict[str, Any], capabilities: list[Any],
                              project_root: Path, errors: list[str]) -> None:
    classes = documents.get("classes")
    if not isinstance(classes, list) or not classes:
        errors.append("documents.classes must contain at least one class/variant when documents are required")
        return
    identities: list[tuple[str, str]] = []
    authorized = 0
    for index, item in enumerate(classes):
        label = f"documents.classes[{index}]"
        if not isinstance(item, dict):
            errors.append(f"{label} must be an object")
            continue
        class_id = item.get("class_id")
        variant_id = item.get("variant_id")
        for field, value in (("class_id", class_id), ("variant_id", variant_id)):
            if (not isinstance(value, str) or not CANONICAL_DOCUMENT_ID.fullmatch(value)
                    or not 3 <= len(value) <= 128):
                errors.append(f"{label}.{field} must be a canonical lowercase hyphenated ID")
        if isinstance(class_id, str) and isinstance(variant_id, str):
            identities.append((class_id, variant_id))
        if item.get("decision") not in DOCUMENT_DECISIONS:
            errors.append(f"{label}.decision is invalid")
        if not nonempty(item.get("owner")):
            errors.append(f"{label}.owner is required")
        evidence_list(project_root, item.get("inventory_evidence"), f"{label}.inventory_evidence", errors)
        automatic = item.get("automatic_storage_authorized")
        if not isinstance(automatic, bool):
            errors.append(f"{label}.automatic_storage_authorized must be boolean")
            continue
        if not automatic:
            if item.get("persistence_mode") not in (None, "") or item.get("mapping") not in (None, {}):
                errors.append(f"{label} cannot carry persistence configuration when automatic storage is not authorized")
            continue
        authorized += 1
        if item.get("decision") != "READY_FOR_AUTOMATIC_STORAGE":
            errors.append(f"{label}.decision must be READY_FOR_AUTOMATIC_STORAGE when automatic storage is authorized")
        persistence_mode = item.get("persistence_mode")
        if persistence_mode not in {"TRANSACTIONAL_DIRECT", "OUTBOX_CDC_CONSUMER"}:
            errors.append(f"{label}.persistence_mode is invalid for automatic storage")
        mapping = item.get("mapping")
        if not isinstance(mapping, dict):
            errors.append(f"{label}.mapping must be an object for automatic storage")
            continue
        if set(mapping) != MAPPING_FIELDS:
            errors.append(f"{label}.mapping must contain exactly the consumer mapping fields")
        for field in ("mappingId", "mappingDomainType", "mappingVersion", "mappingExpression",
                      "mappingExactOutputKeysCsv", "mappingInputSchemaId", "mappingOutputSchemaId"):
            if not nonempty(mapping.get(field)):
                errors.append(f"{label}.mapping.{field} is required for automatic storage")
        if not isinstance(mapping.get("mappingId"), str) or not MAPPING_ID.fullmatch(mapping["mappingId"]):
            errors.append(f"{label}.mapping.mappingId must match the consumer mapping ID contract")
        if not isinstance(mapping.get("mappingVersion"), str) or not MAPPING_VERSION.fullmatch(mapping["mappingVersion"]):
            errors.append(f"{label}.mapping.mappingVersion must match the consumer semantic version contract")
        for field in ("mappingInputSchemaId", "mappingOutputSchemaId"):
            if not isinstance(mapping.get(field), str) or not SCHEMA_ID.fullmatch(mapping[field]):
                errors.append(f"{label}.mapping.{field} must match the consumer schema ID contract")
        if mapping.get("mappingDomainType") != class_id:
            errors.append(f"{label}.mapping.mappingDomainType must equal class_id")
        expression = mapping.get("mappingExpression")
        if isinstance(expression, str):
            if expression != expression.strip() or "\n" in expression or "\r" in expression or len(expression) > 4096:
                errors.append(f"{label}.mapping.mappingExpression must be trimmed single-line and at most 4096 characters")
            expected_hash = hashlib.sha256(expression.encode("utf-8")).hexdigest()
            if mapping.get("mappingExpressionSha256") != expected_hash:
                errors.append(f"{label}.mapping.mappingExpressionSha256 must match exact UTF-8 expression bytes")
        for field in ("mappingExpressionSha256", "mappingInputSchemaSha256", "mappingOutputSchemaSha256"):
            value = mapping.get(field)
            if not isinstance(value, str) or not SHA256_HEX.fullmatch(value):
                errors.append(f"{label}.mapping.{field} must be lowercase SHA-256")
        keys_csv = mapping.get("mappingExactOutputKeysCsv")
        if isinstance(keys_csv, str):
            keys = keys_csv.split(",")
            if (not keys or any(not OUTPUT_KEY.fullmatch(key) for key in keys)
                    or len(keys) != len(set(keys)) or keys != sorted(keys)):
                errors.append(f"{label}.mapping.mappingExactOutputKeysCsv must be sorted unique canonical CSV")
        for field in ("corpus_evidence", "evaluation_evidence", "approval_evidence"):
            evidence_list(project_root, item.get(field), f"{label}.{field}", errors)
        persistence_dependencies = {"TX-DATABASE", "OBSERVABILITY", "BACKUP-DR", "PRIVACY-COMPLIANCE"}
        if persistence_mode == "OUTBOX_CDC_CONSUMER":
            persistence_dependencies.update({"OUTBOX-INBOX", "BROKER-STREAMING"})
        classifications = {
            candidate.get("id"): candidate.get("classification")
            for candidate in capabilities if isinstance(candidate, dict)
        }
        for dependency in sorted(persistence_dependencies):
            if classifications.get(dependency) != "REQUIRED":
                errors.append(f"{label}.persistence_mode requires capabilities.{dependency}=REQUIRED")
    if len(identities) != len(set(identities)):
        errors.append("documents.classes class_id/variant_id pairs must be unique")
    requested = documents.get("automatic_storage_requested")
    if requested is True and authorized == 0:
        errors.append("documents.automatic_storage_requested requires at least one authorized class")
    if requested is False and authorized:
        errors.append("documents cannot authorize a class when automatic_storage_requested is false")


def scan_secret_keys(value: Any, path: str, errors: list[str]) -> None:
    if isinstance(value, dict):
        for key, child in value.items():
            child_path = f"{path}.{key}"
            if PROHIBITED_SECRET_KEYS.search(str(key)) and child not in (None, "", [], {}):
                errors.append(f"{child_path}: secret material is prohibited; store only an identity reference")
            scan_secret_keys(child, child_path, errors)
    elif isinstance(value, list):
        for index, child in enumerate(value):
            scan_secret_keys(child, f"{path}[{index}]", errors)


def validate_required_files(project_root: Path, errors: list[str]) -> None:
    for relative in REQUIRED_FILES:
        resolve_evidence(project_root, relative, "required_artifact", errors)
    specs = [path for path in (project_root / "specs").glob("*/spec.md") if path.is_file() and not path.is_symlink() and path.stat().st_size]
    if not specs:
        errors.append("required_artifact: at least one non-empty specs/<feature>/spec.md is required")


def validate_record(record: Any, project_root: Path) -> list[str]:
    errors: list[str] = []
    if not isinstance(record, dict):
        return ["record must be a JSON object"]
    if record.get("schema") != SCHEMA:
        errors.append(f"schema must be {SCHEMA}")
    scan_secret_keys(record, "$", errors)
    validate_required_files(project_root, errors)

    project = record.get("project")
    if not isinstance(project, dict):
        errors.append("project must be an object")
        project = {}
    for field in ("id", "owner"):
        if not nonempty(project.get(field)):
            errors.append(f"project.{field} is required")
    if project.get("status") != "READY_TO_BUILD":
        errors.append("project.status must be READY_TO_BUILD")
    platform_mode = project.get("platform_mode")
    if platform_mode not in {"OFFICIAL_PLATFORM", "BUSINESS_CENTRAL_PLATFORM", "CUSTOM_PLATFORM"}:
        errors.append("project.platform_mode must be an implemented non-BLOCK choice")
    if project.get("risk_tier") not in {"LIGHT", "STANDARD", "HIGH", "CRITICAL"}:
        errors.append("project.risk_tier is invalid")
    if project.get("data_classification") not in {"PUBLIC", "INTERNAL", "CONFIDENTIAL", "RESTRICTED"}:
        errors.append("project.data_classification is invalid")
    timestamp = project.get("updated_at")
    try:
        if not nonempty(timestamp) or not timestamp.endswith("Z"):
            raise ValueError
        datetime.fromisoformat(timestamp[:-1] + "+00:00")
    except (TypeError, ValueError):
        errors.append("project.updated_at must be a real ISO-8601 UTC timestamp ending Z")

    rounds = record.get("rounds")
    if not isinstance(rounds, dict) or set(rounds) != set("ABCDEFGH"):
        errors.append("rounds must contain exactly A through H")
    else:
        for round_id in "ABCDEFGH":
            item = rounds[round_id]
            if not isinstance(item, dict) or item.get("status") not in {"ANSWERED", "PROVEN"}:
                errors.append(f"rounds.{round_id}.status must be ANSWERED or PROVEN")
                continue
            evidence_list(project_root, item.get("evidence"), f"rounds.{round_id}.evidence", errors)
            advisory_ref = item.get("advisory_prompt_ref")
            try:
                validate_prompt_reference(project_root, advisory_ref, round_id)
            except AdvisoryError as error:
                errors.append(f"rounds.{round_id}.advisory_prompt_ref: {error}")

    blockers = record.get("blockers")
    blocker_fields = ("critical_unknowns", "required_access_not_proven", "open_high_or_critical_failures",
                      "open_dependency_blockers", "stale_or_conflicted_authorities")
    if not isinstance(blockers, dict):
        errors.append("blockers must be an object")
    else:
        for field in blocker_fields:
            value = blockers.get(field)
            if not isinstance(value, list):
                errors.append(f"blockers.{field} must be an array")
            elif value:
                errors.append(f"blockers.{field} must be empty")

    slice_record = record.get("first_vertical_slice")
    if not isinstance(slice_record, dict):
        errors.append("first_vertical_slice must be an object")
    else:
        if not nonempty(slice_record.get("id")):
            errors.append("first_vertical_slice.id is required")
        for field in ("journeys", "invariants", "acceptance_test_refs"):
            value = slice_record.get(field)
            if not string_list(value) or not value:
                errors.append(f"first_vertical_slice.{field} must be a non-empty string array")
            elif field == "acceptance_test_refs":
                evidence_list(project_root, value, f"first_vertical_slice.{field}", errors)
        resolve_evidence(project_root, slice_record.get("rollback_recovery_ref", ""),
                         "first_vertical_slice.rollback_recovery_ref", errors)

    capabilities = record.get("capabilities")
    if not isinstance(capabilities, list):
        errors.append("capabilities must be an array")
        capabilities = []
    ids = [item.get("id") for item in capabilities if isinstance(item, dict)]
    duplicates = sorted({item for item in ids if ids.count(item) > 1})
    missing = sorted(set(CAPABILITIES) - set(ids))
    extra = sorted(set(ids) - set(CAPABILITIES), key=str)
    if len(capabilities) != len(CAPABILITIES) or duplicates or missing or extra:
        errors.append(f"capabilities must contain exactly 48 unique IDs; missing={missing} extra={extra} duplicates={duplicates}")
    for index, item in enumerate(capabilities):
        if not isinstance(item, dict) or item.get("id") not in CAPABILITIES:
            continue
        label = f"capabilities.{item['id']}"
        classification = item.get("classification")
        if classification not in {"REQUIRED", "OPTIONAL", "NONE_WITH_REASON"}:
            errors.append(f"{label}.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON")
            continue
        if not nonempty(item.get("reason_or_trigger")):
            errors.append(f"{label}.reason_or_trigger is required")
        if classification in {"REQUIRED", "OPTIONAL"} and not nonempty(item.get("owner")):
            errors.append(f"{label}.owner is required")
        conditions = item.get("blocking_conditions")
        if not isinstance(conditions, list) or conditions:
            errors.append(f"{label}.blocking_conditions must be an empty array")
        if classification == "REQUIRED":
            evidence_list(project_root, item.get("implementation_refs"), f"{label}.implementation_refs", errors)
            evidence_list(project_root, item.get("test_refs"), f"{label}.test_refs", errors)
            resolve_evidence(project_root, item.get("rollback_recovery_ref", ""), f"{label}.rollback_recovery_ref", errors)
    classification_by_id = {
        item["id"]: item.get("classification") for item in capabilities
        if isinstance(item, dict) and item.get("id") in CAPABILITIES
    }
    essential = {"PRD-INTAKE", "ARCH-DOMAIN", "REPO-SCM", "CONTRACTS", "RUNTIME",
                 "DOMAIN-MODULES", "CI", "SUPPLY-CHAIN", "SECURITY-APPSEC", "TEST-PLATFORM", "DOCS-OPS"}
    for capability_id in sorted(essential):
        if classification_by_id.get(capability_id) != "REQUIRED":
            errors.append(f"capabilities.{capability_id} must be REQUIRED for a buildable vertical slice")
    delivery_surfaces = {"WEB-PUBLIC", "WEB-PORTALS", "MOBILE", "DESKTOP", "EMBEDDED-IOT", "API-BACKEND"}
    selected_ids = slice_record.get("connected_journey_ids") if isinstance(slice_record, dict) else None
    journeys = record.get("connected_journeys")
    # A CLI/automation entry must belong to the selected slice, not an incidental
    # unselected journey. All its evidence, versions and permissions are validated below.
    operational_entry = (
        string_list(selected_ids) and bool(selected_ids) and isinstance(journeys, list)
        and all(any(isinstance(j, dict) and j.get("id") == selected_id
                    and j.get("status") == "PROVEN" and string_list(j.get("surfaces"))
                    and bool(j["surfaces"]) and set(j["surfaces"]) <= {"CLI", "AUTOMATION"}
                    for j in journeys) for selected_id in selected_ids)
    )
    if not any(classification_by_id.get(item) == "REQUIRED" for item in delivery_surfaces) and not operational_entry:
        errors.append("at least one delivery surface must be REQUIRED or the selected slice must have proven CLI/AUTOMATION journeys")
    dependency_rules = {
        "WEB-PORTALS": {"IDENTITY", "AUTHORIZATION"},
        "MOBILE": {"IDENTITY", "AUTHORIZATION"},
        "DESKTOP": {"IDENTITY", "AUTHORIZATION"},
        "API-BACKEND": {"IDENTITY", "AUTHORIZATION"},
        "TX-DATABASE": {"BACKUP-DR", "OBSERVABILITY", "SECURITY-APPSEC", "PRIVACY-COMPLIANCE"},
        "ML-AI": {"DATA-INGEST", "OBSERVABILITY"},
        "RAG-AGENTS": {"ML-AI", "AUTHORIZATION", "SECURITY-APPSEC"},
        "CONTAINERS": {"SUPPLY-CHAIN"},
        "IAC-CLOUD": {"SUPPLY-CHAIN"},
        "CD-RELEASE": {"SUPPLY-CHAIN"},
        "PAYMENTS": {"INTEGRATIONS"},
        "MARKETPLACES": {"INTEGRATIONS"},
        "ADS-ATTRIBUTION": {"INTEGRATIONS"},
        "EMBEDDED-IOT": {"IDENTITY", "AUTHORIZATION", "OBSERVABILITY"},
    }
    for trigger, dependencies_required in dependency_rules.items():
        if classification_by_id.get(trigger) == "REQUIRED":
            for dependency in sorted(dependencies_required):
                if classification_by_id.get(dependency) != "REQUIRED":
                    errors.append(f"capabilities.{trigger} requires capabilities.{dependency}=REQUIRED")

    validate_connected_journeys(record.get("connected_journeys"), classification_by_id,
                                slice_record if isinstance(slice_record, dict) else {},
                                project_root, errors)

    documents = record.get("documents")
    if not isinstance(documents, dict):
        errors.append("documents must be an object")
    elif documents.get("required") is True:
        resolve_evidence(project_root, documents.get("decision_record", ""), "documents.decision_record", errors)
        if documents.get("corpus_ground_truth_status") != "PROVEN":
            errors.append("documents.corpus_ground_truth_status must be PROVEN when documents are required")
        evidence_list(project_root, documents.get("security_evidence"), "documents.security_evidence", errors)
        if documents.get("automatic_storage_requested") is True:
            evidence_list(project_root, documents.get("field_evaluation_evidence"),
                          "documents.field_evaluation_evidence", errors)
        elif documents.get("automatic_storage_requested") is not False:
            errors.append("documents.automatic_storage_requested must be boolean")
        validate_document_classes(documents, capabilities, project_root, errors)
        for dependency in ("DATA-INGEST", "SECURITY-APPSEC"):
            if classification_by_id.get(dependency) != "REQUIRED":
                errors.append(f"required documents require capabilities.{dependency}=REQUIRED")
    else:
        if documents.get("corpus_ground_truth_status") != "NOT_APPLICABLE" or documents.get("automatic_storage_requested") is not False:
            errors.append("documents not required must be NOT_APPLICABLE with automatic storage false")
        if documents.get("classes") != []:
            errors.append("documents not required must have an empty classes array")

    integrations = record.get("integrations")
    if not isinstance(integrations, list):
        errors.append("integrations must be an array")
    else:
        integration_ids: list[str] = []
        for index, item in enumerate(integrations):
            label = f"integrations[{index}]"
            if not isinstance(item, dict) or not nonempty(item.get("id")):
                errors.append(f"{label}.id is required")
                continue
            integration_ids.append(item["id"])
            requirement = item.get("requirement")
            status = item.get("status")
            if requirement not in {"REQUIRED", "OPTIONAL", "NONE_WITH_REASON"}:
                errors.append(f"{label}.requirement is invalid")
            if requirement == "REQUIRED" and item.get("in_first_slice") is True:
                if status != "PROVEN":
                    errors.append(f"{label}.status must be PROVEN for a required first-slice integration")
                evidence_list(project_root, item.get("evidence"), f"{label}.evidence", errors)
            elif requirement == "REQUIRED":
                if status != "DEFERRED" or not nonempty(item.get("defer_reason")) or not nonempty(item.get("owner")):
                    errors.append(f"{label} required outside first slice must be DEFERRED with reason and owner")
            elif requirement == "OPTIONAL" and status not in {"PROVEN", "DEFERRED", "NOT_APPLICABLE"}:
                errors.append(f"{label}.status is invalid for OPTIONAL")
            elif requirement == "NONE_WITH_REASON" and (status != "NOT_APPLICABLE" or not nonempty(item.get("defer_reason"))):
                errors.append(f"{label} NONE_WITH_REASON requires NOT_APPLICABLE and a reason")
        if len(integration_ids) != len(set(integration_ids)):
            errors.append("integrations IDs must be unique")

    sources = record.get("sources")
    if not isinstance(sources, dict):
        errors.append("sources must be an object")
    else:
        resolve_evidence(project_root, sources.get("lock_path", ""), "sources.lock_path", errors)
        if sources.get("status") != "PROVEN" or sources.get("exact_revisions") is not True or sources.get("licenses_notices_resolved") is not True:
            errors.append("sources must be PROVEN with exact revisions and resolved licenses/notices")
        if sources.get("moving_references") != []:
            errors.append("sources.moving_references must be empty")
        modes = sources.get("product_provenance_modes")
        if not string_list(modes) or not modes:
            errors.append("sources.product_provenance_modes must be non-empty")
        elif platform_mode == "OFFICIAL_PLATFORM" and not set(modes) <= {"VERBATIM", "DEPENDENCY_PIN"}:
            errors.append("OFFICIAL_PLATFORM permits only VERBATIM or DEPENDENCY_PIN product provenance")

    dependencies = record.get("dependencies")
    if not isinstance(dependencies, dict) or dependencies.get("status") != "PROVEN":
        errors.append("dependencies.status must be PROVEN")
    else:
        evidence_list(project_root, dependencies.get("evidence"), "dependencies.evidence", errors)
        for field in ("eol_policy", "locks", "rollback"):
            if dependencies.get(field) is not True:
                errors.append(f"dependencies.{field} must be true")

    authorities = record.get("authorities")
    if not isinstance(authorities, dict) or authorities.get("status") != "PROVEN":
        errors.append("authorities.status must be PROVEN")
    else:
        evidence_list(project_root, authorities.get("evidence"), "authorities.evidence", errors)
        if authorities.get("conflicts") != []:
            errors.append("authorities.conflicts must be empty")

    operations = record.get("operations")
    expected_operations = ("security_privacy", "slo_performance_capacity_cost", "backup_restore_dr",
                           "deployment_canary_rollback", "monitoring_incident_response")
    if not isinstance(operations, dict) or set(operations) != set(expected_operations):
        errors.append("operations must contain exactly the five required control groups")
    else:
        for name in expected_operations:
            item = operations[name]
            if not isinstance(item, dict) or item.get("status") != "PROVEN":
                errors.append(f"operations.{name}.status must be PROVEN")
            else:
                evidence_list(project_root, item.get("evidence"), f"operations.{name}.evidence", errors)

    pack_plan = record.get("pack_plan")
    if not isinstance(pack_plan, dict):
        errors.append("pack_plan must be an object")
    else:
        resolve_evidence(project_root, pack_plan.get("path", ""), "pack_plan.path", errors)
        if pack_plan.get("status") != "PROVEN" or pack_plan.get("collision_check") is not True or pack_plan.get("rollback_defined") is not True:
            errors.append("pack_plan must be PROVEN with collision_check and rollback_defined true")
        evidence_list(project_root, pack_plan.get("evidence"), "pack_plan.evidence", errors)

    return sorted(set(errors))


class ReadinessInputError(Exception):
    """A rejected CLI input; never an approval or a retry instruction."""


def atomic_write(path: Path, payload: bytes) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    temporary = None
    try:
        with tempfile.NamedTemporaryFile(dir=path.parent, prefix=f".{path.name}.", delete=False) as handle:
            temporary = Path(handle.name)
            handle.write(payload)
            handle.flush()
            os.fsync(handle.fileno())
        # Publish complete bytes exclusively. Unlike replace(), link() refuses
        # a destination created after the caller's preliminary exists check.
        # Unsupported filesystems fail closed; do not fall back to replacement.
        os.link(temporary, path)
    finally:
        if temporary is not None:
            temporary.unlink(missing_ok=True)


def main(argv: Iterable[str] | None = None) -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--project-root", required=True, type=Path)
    parser.add_argument("--record", default="PROJECT_READINESS_GATE.json")
    parser.add_argument("--report", type=Path)
    parser.add_argument("--scope", choices=("PROJECT", "LIBRARY_INFRASTRUCTURE"), default="PROJECT")
    parser.add_argument("--level", choices=("preparation", "release"))
    args = parser.parse_args(argv)
    try:
        root = args.project_root.resolve(strict=True)
        if not root.is_dir():
            raise ReadinessInputError("project root must be a directory")
        record_path = (root / args.record).resolve(strict=True)
        if record_path.parent != root or record_path.is_symlink() or not record_path.is_file():
            raise ReadinessInputError("record must be a regular file directly under project root")
        try:
            record = json.loads(record_path.read_bytes())
        except (UnicodeDecodeError, json.JSONDecodeError) as error:
            raise ReadinessInputError(f"invalid readiness JSON: {error}") from error
        if args.scope == "LIBRARY_INFRASTRUCTURE":
            from validate_library_readiness import validate_library_record
            if args.level is None:
                raise ReadinessInputError("library scope requires explicit --level preparation or release")
            errors = validate_library_record(record, root, args.level)
            ready = "READY_FOR_LIBRARY_WORK" if args.level == "preparation" else "READY_FOR_LIBRARY_USE"
            report = {"schema": "elite-library-readiness-report/v1", "record": record_path.name,
                      "record_sha256": hashlib.sha256(record_path.read_bytes()).hexdigest(),
                      "scope": args.scope, "level": args.level, "production_authorized": False,
                      "status": ready if not errors else "BLOCKED", "errors": errors}
        else:
            if args.level is not None:
                raise ReadinessInputError("--level is only valid for explicit library scope")
            errors = validate_record(record, root)
            report = {"schema": "elite-project-readiness-report/v1", "record": record_path.name,
                      "status": "READY_TO_BUILD" if not errors else "BLOCKED", "errors": errors,
                      "capabilities_expected": len(CAPABILITIES)}
        payload = (json.dumps(report, sort_keys=True, indent=2) + "\n").encode()
        if args.report:
            report_path = args.report if args.report.is_absolute() else root / args.report
            report_path = report_path.resolve(strict=False)
            try:
                report_path.relative_to(root)
            except ValueError as error:
                raise ReadinessInputError("report must stay within project root") from error
            if report_path.exists():
                raise ReadinessInputError("report path already exists")
            atomic_write(report_path, payload)
        sys.stdout.buffer.write(payload)
        return 0 if not errors else 2
    except (ReadinessInputError, OSError, UnicodeError) as error:
        print(f"PROJECT_READINESS_BLOCKED: {error}", file=sys.stderr)
        return 2


if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `project_readiness_gate/render_consumer_profile.py`
```yaml
block_id: "PROJECT-START-READINESS-VALIDATOR:consumer-profile-renderer:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite deterministic readiness-to-consumer orchestration"
license: "LicenseRef-Workspace-Owner"
sha256: "16978a72f46428f083062fa61062d37d806baefd4ae53d5fe8f8cd9694a2d53a"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
import json
from pathlib import Path, PurePosixPath
import sys
from typing import Iterable

from validate_project_readiness import MAPPING_FIELDS, atomic_write, validate_record


PROOF_FIELDS = {
    "migrationAppliedByOwner", "runtimeRoleIsNonOwner", "tlsAndAuthenticationProven",
    "duplicateRaceRollbackProven", "restartRedeliveryOrderingProven", "loadAndRecoveryProven",
    "observabilityRedactionProven", "mappingSchemasAndCorpusApproved",
    "mappingDeterminismAndLimitsProven", "rollbackProven",
}
STRING_FIELDS = {
    "topic", "consumerGroup", "bootstrapServers", "databaseJdbcUrlSecretRef",
    "databaseUsernameSecretRef", "databasePasswordSecretRef", "tlsTruststoreSecretRef",
    "tlsKeystoreSecretRef", "schemaRegistryUrl",
}
CONSUMER_FIELDS = {
    "schemaVersion", "acknowledgeConditionedState", *STRING_FIELDS, *MAPPING_FIELDS,
    "owners", "approvals", "proofs",
}


class RenderError(ValueError):
    pass


def reject_duplicate_keys(pairs: list[tuple[str, object]]) -> dict[str, object]:
    result: dict[str, object] = {}
    for key, value in pairs:
        if key in result:
            raise RenderError(f"duplicate JSON key: {key}")
        result[key] = value
    return result


def load_json(path: Path, label: str) -> dict[str, object]:
    try:
        value = json.loads(path.read_text(encoding="utf-8"), object_pairs_hook=reject_duplicate_keys)
    except (UnicodeDecodeError, json.JSONDecodeError) as error:
        raise RenderError(f"invalid {label} JSON: {error}") from error
    if not isinstance(value, dict):
        raise RenderError(f"{label} must be a JSON object")
    return value


def safe_path(root: Path, value: str, label: str, *, must_exist: bool) -> Path:
    if not value or "\\" in value:
        raise RenderError(f"{label} must be a canonical project-relative path")
    relative = PurePosixPath(value)
    if relative.is_absolute() or any(part in {"", ".", ".."} for part in relative.parts):
        raise RenderError(f"{label} must be a canonical project-relative path")
    candidate = root.joinpath(*relative.parts)
    current = root
    for part in relative.parts[:-1]:
        current = current / part
        if current.exists() and current.is_symlink():
            raise RenderError(f"{label} cannot traverse a symbolic link")
    resolved = candidate.resolve(strict=must_exist)
    try:
        resolved.relative_to(root)
    except ValueError as error:
        raise RenderError(f"{label} escapes project root") from error
    if must_exist and (candidate.is_symlink() or not resolved.is_file()):
        raise RenderError(f"{label} must be a regular non-symlink file")
    return resolved


def select_mapping(record: dict[str, object], class_id: str, variant_id: str) -> dict[str, object]:
    documents = record.get("documents")
    classes = documents.get("classes") if isinstance(documents, dict) else None
    matches = [
        item for item in classes or []
        if isinstance(item, dict) and item.get("class_id") == class_id and item.get("variant_id") == variant_id
    ]
    if len(matches) != 1:
        raise RenderError("class_id/variant_id must select exactly one readiness class")
    selected = matches[0]
    if selected.get("decision") != "READY_FOR_AUTOMATIC_STORAGE" or selected.get("automatic_storage_authorized") is not True:
        raise RenderError("selected class is not authorized for automatic storage")
    if selected.get("persistence_mode") != "OUTBOX_CDC_CONSUMER":
        raise RenderError("selected class is not configured for OUTBOX_CDC_CONSUMER")
    mapping = selected.get("mapping")
    if not isinstance(mapping, dict) or set(mapping) != MAPPING_FIELDS:
        raise RenderError("selected class does not contain the exact consumer mapping contract")
    return mapping


def render_profile(
    record: dict[str, object], project_root: Path, class_id: str, variant_id: str,
    template: dict[str, object],
) -> dict[str, object]:
    errors = validate_record(record, project_root)
    if errors:
        raise RenderError("readiness record is BLOCKED: " + " | ".join(errors))
    mapping = select_mapping(record, class_id, variant_id)
    if set(template) != CONSUMER_FIELDS or template.get("schemaVersion") != 1:
        raise RenderError("consumer template has an unexpected schema")
    if template.get("acknowledgeConditionedState") is not False:
        raise RenderError("consumer template must remain unacknowledged")
    if template.get("owners") != [] or template.get("approvals") != []:
        raise RenderError("consumer template owners and approvals must start empty")
    proofs = template.get("proofs")
    if not isinstance(proofs, dict) or set(proofs) != PROOF_FIELDS or any(value is not False for value in proofs.values()):
        raise RenderError("consumer template proofs must contain the exact closed proof set")
    if any(template.get(field) != "" for field in STRING_FIELDS | MAPPING_FIELDS):
        raise RenderError("consumer template configuration and mapping fields must start empty")
    output = dict(template)
    for field in MAPPING_FIELDS:
        output[field] = mapping[field]
    return output


def main(argv: Iterable[str] | None = None) -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--project-root", required=True, type=Path)
    parser.add_argument("--record", default="PROJECT_READINESS_GATE.json")
    parser.add_argument("--consumer-template", required=True)
    parser.add_argument("--class-id", required=True)
    parser.add_argument("--variant-id", required=True)
    parser.add_argument("--output", default="PROJECT_DOCUMENT_CONSUMER_PROFILE.json")
    args = parser.parse_args(argv)
    root = args.project_root.resolve(strict=True)
    if not root.is_dir():
        raise SystemExit("project root must be a directory")
    try:
        record_path = safe_path(root, args.record, "record", must_exist=True)
        template_path = safe_path(root, args.consumer_template, "consumer template", must_exist=True)
        output_path = safe_path(root, args.output, "output", must_exist=False)
        if output_path.exists():
            raise RenderError("output path already exists")
        record = load_json(record_path, "readiness record")
        template = load_json(template_path, "consumer template")
        profile = render_profile(record, root, args.class_id, args.variant_id, template)
        payload = (json.dumps(profile, sort_keys=True, indent=2) + "\n").encode("utf-8")
        atomic_write(output_path, payload)
    except (OSError, RenderError) as error:
        print(f"CONSUMER_PROFILE_RENDER_BLOCKED: {error}", file=sys.stderr)
        return 2
    print(f"CONSUMER_PROFILE_RENDER_PASS output={output_path.relative_to(root).as_posix()}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `project_readiness_gate/project-advisory-catalog.json`
```yaml
block_id: "PROJECT-START-READINESS-VALIDATOR:advisory-catalog:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite curated advisory catalog governed by official OpenAI, Google, NASA, Microsoft, AWS and ARCA authorities"
license: "LicenseRef-Workspace-Owner"
sha256: "d5c336bf2c493e334ecc253aac5ed3d9e7c90d0ebd09fb78120da0c28bff6015"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-project-advisory-catalog/v1",
  "catalog_version": "1.0.0",
  "sources": [
    {
      "id": "OPENAI-LEAN-INSTRUCTIONS",
      "owner": "OpenAI",
      "url": "https://developers.openai.com/api/docs/guides/latest-model",
      "claim": "Keep agent instructions lean, state each instruction once, expose relevant tools and validate changes on representative work."
    },
    {
      "id": "GOOGLE-RELIABLE-LAUNCHES",
      "owner": "Google",
      "url": "https://sre.google/sre-book/reliable-product-launches/",
      "claim": "Use curated, practical and adaptable launch questions, shared infrastructure, explicit dependencies, failure modes and staged rollout."
    },
    {
      "id": "NASA-SE-HANDBOOK",
      "owner": "NASA",
      "url": "https://www.nasa.gov/reference/system-engineering-handbook-appendix/",
      "claim": "Separate requirements, verification, validation, integration and stakeholder acceptance evidence."
    },
    {
      "id": "MICROSOFT-ABSENCE",
      "owner": "Microsoft",
      "url": "https://learn.microsoft.com/en-us/dynamics365/business-central/hr-how-manage-absence",
      "claim": "Employee absence is explicit operational data with causes, periods and units; it is distinct from software licensing."
    },
    {
      "id": "MICROSOFT-POST-SALES",
      "owner": "Microsoft",
      "url": "https://learn.microsoft.com/en-us/dynamics365/business-central/ui-post-sales",
      "claim": "Posting sales documents creates durable operational and accounting consequences and requires explicit business decisions."
    },
    {
      "id": "GOOGLE-DOCUMENT-AI",
      "owner": "Google",
      "url": "https://cloud.google.com/document-ai/docs/processors-list",
      "claim": "Document processing capabilities are processor- and document-type-specific; project data, schemas and evaluation remain project inputs."
    },
    {
      "id": "AWS-TEXTRACT",
      "owner": "Amazon Web Services",
      "url": "https://docs.aws.amazon.com/textract/latest/dg/limits-document.html",
      "claim": "Document services have explicit format, size, page, operation and service limits that must be selected and tested."
    },
    {
      "id": "ARCA-WSFE",
      "owner": "ARCA",
      "url": "https://www.arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG.pdf",
      "claim": "Argentine electronic invoicing requires the applicable official service contract, authentication and homologation/production authority."
    }
  ],
  "topics": [
    {
      "id": "A-OUTCOME-AUTHORITY",
      "round": "A",
      "title": "Resultado, responsables y límites reales",
      "plain_language": "Define qué negocio se construye, quién puede decidir, dónde operará, qué fecha y presupuesto existen y qué queda fuera. El owner es la persona que puede aceptar o rechazar una decisión; no es un nombre decorativo.",
      "why_required": "Sin resultado medible, jurisdicción y autoridad, el agente no puede distinguir una necesidad real de una suposición ni elegir prioridades, proveedores o evidencia de aceptación.",
      "questions": [
        "¿Cuál es el resultado comercial u operativo medible del primer lanzamiento?",
        "¿Quién decide producto, operación, seguridad y gasto, y qué decisiones requieren tu aprobación?",
        "¿Qué países, idiomas, monedas, fecha, presupuesto y exclusiones gobiernan el primer alcance?"
      ],
      "required_inputs": ["nombre del proyecto", "owners y responsabilidades", "resultado medible", "jurisdicciones", "fecha y presupuesto", "fuera de alcance"],
      "example_answer": "EJEMPLO, NO RESPUESTA ASUMIDA: Argentina; owner comercial Ana, owner técnico Luis; primer objetivo: publicar catálogo y convertir leads medidos; presupuesto mensual con tope aprobado; Brasil y pagos quedan fuera de P1.",
      "verification": ["confirmación del owner", "constitution y spec coherentes", "restricciones y métricas escritas"],
      "not_applicable_when": "Una subpregunta puede marcarse NO_APLICA sólo si no afecta el primer slice y se registra el motivo y el trigger que obligaría a reabrirla.",
      "source_ids": ["OPENAI-LEAN-INSTRUCTIONS", "GOOGLE-RELIABLE-LAUNCHES", "NASA-SE-HANDBOOK"]
    },
    {
      "id": "B-ACTORS-JOURNEYS",
      "round": "B",
      "title": "Personas, permisos y recorridos completos",
      "plain_language": "Describe qué puede hacer cada visitante, cliente, empleado, franquicia y proveedor de principio a fin. Licencia laboral o ausencia significa indisponibilidad de una persona; cancelación revierte una reserva; no-show registra que alguien no asistió; entrega completa incluye preparación, identidad, serie física, aceptación y evidencia.",
      "why_required": "Los roles y estados evitan exposición entre organizaciones, doble reserva, acciones sin autoridad y experiencias donde el usuario queda atrapado a mitad del recorrido.",
      "questions": [
        "¿Qué actores existen y qué organización, objeto y acción puede operar cada uno?",
        "¿Cuáles son los journeys P1, incluidos error, cancelación, no-show, reversa y recuperación?",
        "¿Cómo se administran horarios, ausencias/licencias laborales, turnos, recursos y entrega final?"
      ],
      "required_inputs": ["actores", "matriz de permisos", "journeys Given/When/Then", "reglas de cancelación y no-show", "reglas de ausencia", "evidencia de entrega"],
      "example_answer": "EJEMPLO, NO RESPUESTA ASUMIDA: el visitante consulta modelos sin login; el cliente reserva test-drive; la franquicia confirma recurso; una ausencia bloquea ese recurso; el cliente cancela antes del límite; no-show requiere actor y motivo; la entrega compara cliente, organización y número de serie antes de aceptar.",
      "verification": ["journeys positivos y negativos", "pruebas de autorización y aislamiento", "transiciones e identidad auditables"],
      "not_applicable_when": "Sólo NO_APLICA si el actor o journey no existe en el alcance aprobado; no se permite omitir cancelación o recuperación de una operación que sí existe.",
      "source_ids": ["MICROSOFT-ABSENCE", "GOOGLE-RELIABLE-LAUNCHES", "NASA-SE-HANDBOOK"]
    },
    {
      "id": "C-ENTERPRISE-FISCAL",
      "round": "C",
      "title": "Modelo empresarial, stock, dinero y fiscalidad",
      "plain_language": "Define productos, variantes, series, stock, compras, ventas, pagos, franquicias y contabilidad con sus estados y dueños. Regla fiscal aprobada es una decisión validada por responsable contable; certificado es la identidad criptográfica de la empresa; POS es el punto/canal de venta; homologación ARCA es el ambiente de prueba oficial antes de producción.",
      "why_required": "Estos conceptos generan compromisos de inventario, dinero, impuestos y auditoría. Inventarlos puede producir stock imposible, documentos fiscales inválidos o conciliaciones irreparables.",
      "questions": [
        "¿Cuáles son las entidades, estados, invariantes, fuentes de verdad y reglas configurables del negocio?",
        "¿Cómo fluyen compra, producción, importación, stock, precio, venta, cobro, devolución, regalía y postventa?",
        "¿Qué autoridad contable aprueba reglas fiscales, puntos de venta, certificados, homologación y paso a producción ARCA?"
      ],
      "required_inputs": ["glosario y estados", "owners por agregado", "reglas de stock y dinero", "política de franquicia", "contador/autoridad fiscal", "referencias seguras a certificado y ambientes"],
      "example_answer": "EJEMPLO, NO RESPUESTA ASUMIDA: serie única por vehículo; reserva no descuenta stock hasta confirmación; devolución genera reversa; contador externo aprueba IVA y comprobantes; certificado se entrega por secret store; primero homologación ARCA y luego producción con aprobación separada.",
      "verification": ["invariantes y concurrencia", "conciliación y reversa", "evidencia del responsable contable", "probe separado de homologación y producción"],
      "not_applicable_when": "Fiscalidad o una etapa comercial sólo puede ser NO_APLICA si el primer slice no emite, cobra, registra ni promete esa operación y existe trigger explícito para reabrirla.",
      "source_ids": ["MICROSOFT-POST-SALES", "ARCA-WSFE", "NASA-SE-HANDBOOK"]
    },
    {
      "id": "D-DOCUMENT-EVIDENCE",
      "round": "D",
      "title": "Archivos, extracción y almacenamiento preciso",
      "plain_language": "Cada clase y variante identifica un layout real. Corpus es el conjunto representativo de archivos; ground truth son valores revisados; schema define tipos y campos; mapping transforma esos campos a los nombres exactos del sistema consumidor. Ningún proveedor convierte un archivo desconocido en dato correcto sin estas evidencias.",
      "why_required": "Persistir una extracción incorrecta contamina compras, stock, pagos y fiscalidad. La precisión debe demostrarse por clase, variante y campo antes de automatizar.",
      "questions": [
        "¿Qué canales, clases, variantes, emisores, idiomas, layouts, volúmenes y límites de archivo existen?",
        "¿Dónde están originales, checksums, corpus autorizado, ground truth, schemas y reglas cruzadas por clase?",
        "¿Qué errores son tolerables por campo y quién revisa, aprueba, rechaza o corrige ambigüedades?",
        "Si habrá persistencia automática, ¿cuál es el modo y el mapping exacto al consumer con hashes y evidencia?"
      ],
      "required_inputs": ["inventario clase/variante", "archivos autorizados", "ground truth", "schemas", "mapping", "umbrales por campo", "owner de aprobación", "retención y seguridad"],
      "example_answer": "EJEMPLO, NO RESPUESTA ASUMIDA: invoice/supplier-a-v1 y packing-list/factory-b-es-v2; 200 documentos anonimizados por variante; campos y unidades en schemas versionados; número, moneda y total requieren exactitud aprobada; discrepancias van a revisión; persistencia sólo tras evaluación y mapping hash-locked.",
      "verification": ["bytes y SHA-256 del original", "evaluación separada por clase/variante/campo", "schema y mapping reproducibles", "revisión humana donde la política lo exige"],
      "not_applicable_when": "Persistencia automática puede ser NO_APLICA y conservar extracción para revisión. Corpus o schema no son NO_APLICA si el sistema almacenará datos extraídos de esa clase.",
      "source_ids": ["GOOGLE-DOCUMENT-AI", "AWS-TEXTRACT", "NASA-SE-HANDBOOK"]
    },
    {
      "id": "E-INTEGRATION-ACCESS",
      "round": "E",
      "title": "Proveedores, cuentas y accesos demostrables",
      "plain_language": "Una integración no está lista por conocer su nombre. Requiere producto y versión exactos, cuenta owner, sandbox, identidad, scopes, cuotas, costos, webhooks, reconciliación y evidencia de un probe seguro. Los secretos nunca se pegan en chat ni Markdown.",
      "why_required": "El código puede compilar y aun fallar por cuenta, país, permisos, contrato, límite o ambiente incorrectos. El probe evita simular acceso inexistente.",
      "questions": [
        "¿Qué integraciones son REQUIRED, OPTIONAL o NO_APLICA y para qué journey?",
        "¿Quién posee cada cuenta y cómo entregará una referencia segura de identidad sin revelar secretos?",
        "¿Qué sandbox, scopes, cuotas, costos, webhooks, fixtures, reconciliation y soporte se demostrarán?"
      ],
      "required_inputs": ["proveedor/producto/país", "owner de cuenta", "ambiente", "identity reference", "scopes", "cuotas y costos", "probe y evidencia", "fallback/exit plan"],
      "example_answer": "EJEMPLO, NO RESPUESTA ASUMIDA: Mercado Pago sandbox, owner Finanzas, identidad mp-sandbox-franchise en secret store, scopes mínimos, webhook firmado, tope de costo aprobado, fixture de pago y conciliación diaria; producción queda bloqueada hasta probe separado.",
      "verification": ["probe read-only o sandbox autorizado", "cuenta y ambiente observados", "request ID/schema redactados", "fallo cerrado ante ambiente equivocado"],
      "not_applicable_when": "Sólo si ningún journey aprobado depende del proveedor; debe registrarse el motivo y la condición futura que reabre la integración.",
      "source_ids": ["GOOGLE-RELIABLE-LAUNCHES", "OPENAI-LEAN-INSTRUCTIONS"]
    },
    {
      "id": "F-UX-CHANNELS",
      "round": "F",
      "title": "Experiencia completa para público, clientes y empresa",
      "plain_language": "UX completa significa que cada actor entiende qué ocurre, puede corregir errores y finaliza su objetivo en web pública, cliente, franquicia, administración, fábrica o proveedor. Incluye responsive, accesibilidad, rendimiento, contenido, contacto, estados vacíos/error/espera, ayuda contextual, capacitación y soporte vinculados a la misma versión operativa.",
      "why_required": "Una API correcta no garantiza que visitantes conviertan, clientes confíen ni operadores trabajen sin atajos manuales. Una pantalla, capacitación o respuesta de soporte desconectada de permisos, efectos durables y versión tampoco cierra una capacidad. Los canales y journeys deben validarse de extremo a extremo con usuarios y dispositivos reales.",
      "questions": [
        "¿Qué pantallas y canales necesita cada actor para completar cada journey P1?",
        "¿Qué branding, contenido, SEO, analytics, consentimiento, accesibilidad, dispositivos y navegadores se aceptan?",
        "¿Qué presupuestos de rendimiento y pruebas de error, espera, recuperación y entrega final gobiernan la aceptación?",
        "¿Qué ayuda contextual, práctica, capacitación, evaluación, soporte y reentrenamiento necesita cada rol, y cómo se ligan a permisos y versión?"
      ],
      "required_inputs": ["mapa de pantallas", "contenido y assets autorizados", "matriz de dispositivos/navegadores", "criterios de accesibilidad", "performance budgets", "owners editoriales", "matriz rol→journey→ayuda→capacitación→soporte", "versiones de contenido operativo"],
      "example_answer": "EJEMPLO, NO RESPUESTA ASUMIDA: visitante móvil filtra modelos y contacta; cliente ve turno y entrega; franquicia gestiona disponibilidad con ayuda contextual y práctica versionada; admin audita; soporte recibe contexto autorizado y runbook de la misma release; Chrome/Edge/Safari actuales, teclado y lector de pantalla; p75 web vitals objetivo y textos de error aprobados.",
      "verification": ["E2E por actor y efecto durable", "navegadores y responsive", "accesibilidad", "métricas de rendimiento", "ayuda/capacitación/soporte iguales a la release", "reentrenamiento ante cambio de journey", "aceptación del owner de negocio"],
      "not_applicable_when": "Un canal puede ser NO_APLICA si ningún journey P1 lo usa. Estados de error y recuperación no son opcionales para una pantalla incluida; ayuda, capacitación o soporte sólo pueden omitirse con razón aprobada cuando ninguna persona opera ni recibe el efecto.",
      "source_ids": ["GOOGLE-RELIABLE-LAUNCHES", "NASA-SE-HANDBOOK"]
    },
    {
      "id": "G-SECURITY-PRIVACY",
      "round": "G",
      "title": "Datos, aislamiento, seguridad y cumplimiento",
      "plain_language": "Clasifica cada dato, propósito, acceso y retención. Aislamiento significa que una franquicia no puede leer ni modificar otra. IdP autentica identidades; autorización decide acciones; KMS/PKI protege claves y certificados; threat model enumera abusos y defensas.",
      "why_required": "Sin límites de datos e identidad, cualquier módulo puede filtrar información, mezclar franquicias o usar privilegios excesivos aunque sus funciones de negocio pasen tests.",
      "questions": [
        "¿Qué datos son personales, financieros, secretos o restringidos y con qué propósito, base, retención y borrado?",
        "¿Qué IdP, MFA, sesiones, roles, políticas por organización y break-glass se usarán?",
        "¿Qué amenazas, fraude, abuso, logs, auditoría, pruebas ofensivas y acceptance owner deben demostrarse?"
      ],
      "required_inputs": ["clasificación por campo", "política de tenant", "IdP y roles", "secret/KMS references", "retención y derechos", "threat model", "plan de pruebas"],
      "example_answer": "EJEMPLO, NO RESPUESTA ASUMIDA: email CONFIDENTIAL por atención comercial; retención aprobada; OIDC con MFA para administración; policy por organization_id; certificados sólo por referencia KMS; pruebas negativas entre dos franquicias y revisión ofensiva antes de piloto.",
      "verification": ["tests de aislamiento", "tokens y scopes reales redactados", "secret scanning", "audit logs sin PII innecesaria", "security acceptance"],
      "not_applicable_when": "Una categoría de dato puede ser NO_APLICA si se demuestra que no se recoge ni deriva. Autorización y aislamiento no son NO_APLICA cuando existen usuarios u organizaciones múltiples.",
      "source_ids": ["GOOGLE-RELIABLE-LAUNCHES", "NASA-SE-HANDBOOK"]
    },
    {
      "id": "H-OPERATIONS-PRODUCTION",
      "round": "H",
      "title": "Plataforma, operación, recuperación y producción",
      "plain_language": "Distingue local, CI, DEV, STAGE y PROD. SLI/SLO mide el servicio; RPO/RTO define pérdida y tiempo de recuperación; PITR restaura a un instante; canary expone gradualmente; rollback revierte. CDN/WAF, IdP, PostgreSQL, proveedores y observabilidad deben existir realmente en el target elegido.",
      "why_required": "Compilar y probar localmente no demuestra capacidad, seguridad, recuperación ni operación productiva. La aceptación necesita el artefacto y entorno reales con owners y evidencia.",
      "questions": [
        "¿Qué plataforma, regiones, red, DNS/TLS, CDN/WAF, identidades, IaC y entornos se usarán?",
        "¿Cuáles son SLO, carga, capacidad, costo, alertas, on-call, runbooks y respuesta a incidentes?",
        "¿Cómo se probarán backup/PITR/restore, RPO/RTO, canary, rollback, seguridad y piloto?",
        "¿Qué evidencia y qué autoridad separada habilitan DEV, piloto y producción?"
      ],
      "required_inputs": ["target y accesos", "topología e IaC", "SLO/carga/costo", "on-call/runbooks", "backup/restore", "rollout/rollback", "acceptance owners"],
      "example_answer": "EJEMPLO, NO RESPUESTA ASUMIDA: DEV en región aprobada con IaC; PostgreSQL gestionado; CDN/WAF e IdP definidos; SLO y carga medidos; restore en entorno aislado; canary al 5%; rollback automático; owner técnico habilita DEV y comité negocio-seguridad autoriza piloto/PROD con evidencia separada.",
      "verification": ["deploy del artefacto exacto", "prueba de carga y seguridad", "restore/PITR observado", "canary y rollback", "aceptación empresarial del piloto"],
      "not_applicable_when": "Una tecnología puede ser NO_APLICA por arquitectura demostrada; producción nunca se considera cubierta por evidencia sólo local o DEV.",
      "source_ids": ["GOOGLE-RELIABLE-LAUNCHES", "NASA-SE-HANDBOOK", "OPENAI-LEAN-INSTRUCTIONS"]
    }
  ]
}
````

### FILE: `project_readiness_gate/render_project_advisory.py`
```yaml
block_id: "PROJECT-START-READINESS-VALIDATOR:advisory-renderer:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite deterministic project advisory rendering and receipt validation"
license: "LicenseRef-Workspace-Owner"
sha256: "3b287f648be48baca56677de831d933d36410c14d3e2200a955fa96ca6df2105"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import re
import sys
from typing import Any, Iterable
from urllib.parse import urlparse


SCHEMA = "elite-project-advisory-catalog/v1"
CATALOG_FILE = "project-advisory-catalog.json"
ROUNDS = tuple("ABCDEFGH")
TOP_KEYS = {"schema", "catalog_version", "sources", "topics"}
SOURCE_KEYS = {"id", "owner", "url", "claim"}
TOPIC_KEYS = {
    "id", "round", "title", "plain_language", "why_required", "questions",
    "required_inputs", "example_answer", "verification", "not_applicable_when", "source_ids",
}
ALLOWED_SOURCE_HOSTS = {
    "developers.openai.com", "sre.google", "www.nasa.gov", "learn.microsoft.com",
    "cloud.google.com", "docs.aws.amazon.com", "www.arca.gob.ar",
}
IDENTIFIER = re.compile(r"^[A-Z][A-Z0-9-]{2,63}$")
SEMVER = re.compile(r"^[0-9]+\.[0-9]+\.[0-9]+$")


class AdvisoryError(ValueError):
    pass


def nonempty(value: Any) -> bool:
    return isinstance(value, str) and bool(value.strip()) and value == value.strip()


def strict_object(pairs: list[tuple[str, Any]]) -> dict[str, Any]:
    result: dict[str, Any] = {}
    for key, value in pairs:
        if key in result:
            raise AdvisoryError(f"duplicate JSON key: {key}")
        result[key] = value
    return result


def string_array(value: Any) -> bool:
    return isinstance(value, list) and bool(value) and all(nonempty(item) for item in value)


def load_catalog(path: Path | None = None) -> tuple[dict[str, Any], str]:
    catalog_path = path or Path(__file__).with_name(CATALOG_FILE)
    try:
        payload = catalog_path.read_bytes()
    except OSError as error:
        raise AdvisoryError(f"cannot read advisory catalog: {error}") from error
    try:
        catalog = json.loads(payload.decode("utf-8"), object_pairs_hook=strict_object)
    except (UnicodeDecodeError, json.JSONDecodeError, AdvisoryError) as error:
        raise AdvisoryError(f"invalid advisory catalog JSON: {error}") from error
    validate_catalog(catalog)
    return catalog, hashlib.sha256(payload).hexdigest()


def validate_catalog(catalog: Any) -> None:
    if not isinstance(catalog, dict) or set(catalog) != TOP_KEYS:
        raise AdvisoryError(f"catalog keys must be exactly {sorted(TOP_KEYS)}")
    if catalog.get("schema") != SCHEMA:
        raise AdvisoryError(f"catalog schema must be {SCHEMA}")
    if not isinstance(catalog.get("catalog_version"), str) or not SEMVER.fullmatch(catalog["catalog_version"]):
        raise AdvisoryError("catalog_version must be semantic x.y.z")

    sources = catalog.get("sources")
    if not isinstance(sources, list) or not sources:
        raise AdvisoryError("sources must be a non-empty array")
    source_ids: set[str] = set()
    for index, source in enumerate(sources):
        label = f"sources[{index}]"
        if not isinstance(source, dict) or set(source) != SOURCE_KEYS:
            raise AdvisoryError(f"{label} keys must be exactly {sorted(SOURCE_KEYS)}")
        source_id = source.get("id")
        if not isinstance(source_id, str) or not IDENTIFIER.fullmatch(source_id):
            raise AdvisoryError(f"{label}.id is invalid")
        if source_id in source_ids:
            raise AdvisoryError(f"duplicate source id: {source_id}")
        source_ids.add(source_id)
        for field in ("owner", "claim"):
            if not nonempty(source.get(field)):
                raise AdvisoryError(f"{label}.{field} is required")
        url = source.get("url")
        parsed = urlparse(url if isinstance(url, str) else "")
        if parsed.scheme != "https" or parsed.hostname not in ALLOWED_SOURCE_HOSTS or parsed.username or parsed.password:
            raise AdvisoryError(f"{label}.url is not an admitted official HTTPS authority")

    topics = catalog.get("topics")
    if not isinstance(topics, list) or not topics:
        raise AdvisoryError("topics must be a non-empty array")
    topic_ids: set[str] = set()
    round_counts = {round_id: 0 for round_id in ROUNDS}
    for index, topic in enumerate(topics):
        label = f"topics[{index}]"
        if not isinstance(topic, dict) or set(topic) != TOPIC_KEYS:
            raise AdvisoryError(f"{label} keys must be exactly {sorted(TOPIC_KEYS)}")
        topic_id = topic.get("id")
        if not isinstance(topic_id, str) or not IDENTIFIER.fullmatch(topic_id):
            raise AdvisoryError(f"{label}.id is invalid")
        if topic_id in topic_ids:
            raise AdvisoryError(f"duplicate topic id: {topic_id}")
        topic_ids.add(topic_id)
        round_id = topic.get("round")
        if round_id not in ROUNDS:
            raise AdvisoryError(f"{label}.round must be A through H")
        round_counts[round_id] += 1
        for field in ("title", "plain_language", "why_required", "example_answer", "not_applicable_when"):
            if not nonempty(topic.get(field)):
                raise AdvisoryError(f"{label}.{field} is required")
        if not topic["example_answer"].startswith("EJEMPLO, NO RESPUESTA ASUMIDA:"):
            raise AdvisoryError(f"{label}.example_answer must declare that it is not an assumed answer")
        for field in ("questions", "required_inputs", "verification", "source_ids"):
            if not string_array(topic.get(field)):
                raise AdvisoryError(f"{label}.{field} must be a non-empty string array")
        unknown_sources = sorted(set(topic["source_ids"]) - source_ids)
        if unknown_sources:
            raise AdvisoryError(f"{label}.source_ids unknown={unknown_sources}")
        if len(topic["source_ids"]) != len(set(topic["source_ids"])):
            raise AdvisoryError(f"{label}.source_ids contains duplicates")
    if any(count != 1 for count in round_counts.values()):
        raise AdvisoryError(f"catalog must contain exactly one topic per round: {round_counts}")


def bullet_lines(values: list[str]) -> list[str]:
    return [f"- {value}" for value in values]


def render_round(catalog: dict[str, Any], digest: str, round_id: str) -> bytes:
    if round_id not in ROUNDS:
        raise AdvisoryError("round must be A through H")
    topic = next((item for item in catalog["topics"] if item["round"] == round_id), None)
    if topic is None:
        raise AdvisoryError(f"catalog has no topic for round {round_id}")
    sources = {item["id"]: item for item in catalog["sources"]}
    lines = [
        f"# Asesoramiento de proyecto — ronda {round_id}", "",
        f"Catalog-Version: {catalog['catalog_version']}",
        f"Catalog-SHA256: {digest}",
        f"Topic-ID: {topic['id']}", "",
        "Este prompt fue generado desde el catálogo ejecutable. El agente debe explicarlo en lenguaje simple, presentar un grupo manejable de preguntas y registrar la respuesta en los artefactos de readiness. No debe completar respuestas por inferencia ni pedir secretos en chat.", "",
        f"## {topic['title']}", "", topic["plain_language"], "",
        "## Por qué se necesita", "", topic["why_required"], "",
        "## Preguntas que el agente debe ayudarte a responder", "",
        *bullet_lines(topic["questions"]), "",
        "## Información que puedes proporcionar", "",
        *bullet_lines(topic["required_inputs"]), "",
        "## Ejemplo ilustrativo", "", topic["example_answer"], "",
        "## Cómo se comprobará", "", *bullet_lines(topic["verification"]), "",
        "## Cuándo puede responderse NO_APLICA", "", topic["not_applicable_when"], "",
        "## Autoridades oficiales", "",
    ]
    for source_id in topic["source_ids"]:
        source = sources[source_id]
        lines.append(f"- {source['owner']} — {source['claim']} — {source['url']}")
    lines.extend(["", "## Regla de avance", "", "La ronda sólo se cierra con respuesta o evidencia real. Si falta una decisión o acceso material, el agente explica exactamente qué falta, conserva BLOCKED/AWAITING_USER y continúa con trabajo read-only que no dependa de inventarlo.", ""])
    return "\n".join(lines).encode("utf-8")


def canonical_prompt_name(round_id: str) -> str:
    if round_id not in ROUNDS:
        raise AdvisoryError("round must be A through H")
    return f"PROJECT_ADVISORY_{round_id}.md"


def direct_project_path(project_root: Path, reference: str, expected_name: str) -> Path:
    if not nonempty(reference) or "\\" in reference:
        raise AdvisoryError("advisory prompt reference must be a canonical relative path")
    pure = PurePosixPath(reference)
    if pure.is_absolute() or len(pure.parts) != 1 or pure.parts[0] in {".", ".."}:
        raise AdvisoryError("advisory prompt reference must be a direct project-root file")
    if pure.name != expected_name:
        raise AdvisoryError(f"advisory prompt reference must be {expected_name}")
    root = project_root.resolve(strict=True)
    candidate = root / pure.name
    if candidate.is_symlink():
        raise AdvisoryError("advisory prompt must not be a symlink")
    return candidate


def validate_prompt_reference(project_root: Path, reference: str, round_id: str) -> None:
    catalog, digest = load_catalog()
    expected_name = canonical_prompt_name(round_id)
    candidate = direct_project_path(project_root, reference, expected_name)
    if not candidate.is_file():
        raise AdvisoryError(f"advisory prompt is missing: {reference}")
    try:
        observed = candidate.read_bytes()
    except OSError as error:
        raise AdvisoryError(f"cannot read advisory prompt: {error}") from error
    expected = render_round(catalog, digest, round_id)
    if observed != expected:
        raise AdvisoryError(f"advisory prompt bytes do not match catalog for round {round_id}")


def atomic_create(path: Path, payload: bytes) -> None:
    flags = os.O_WRONLY | os.O_CREAT | os.O_EXCL
    descriptor = os.open(path, flags, 0o600)
    try:
        with os.fdopen(descriptor, "wb") as stream:
            stream.write(payload)
            stream.flush()
            os.fsync(stream.fileno())
    except Exception:
        try:
            path.unlink(missing_ok=True)
        finally:
            raise


def main(argv: Iterable[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description="Render one fail-closed Elite project advisory round")
    parser.add_argument("--project-root", default=".")
    parser.add_argument("--round", required=True, choices=ROUNDS)
    parser.add_argument("--output", required=True)
    args = parser.parse_args(argv)
    try:
        project_root = Path(args.project_root).resolve(strict=True)
        expected_name = canonical_prompt_name(args.round)
        output = direct_project_path(project_root, args.output, expected_name)
        if output.exists() or output.is_symlink():
            raise AdvisoryError(f"refusing to overwrite advisory prompt: {args.output}")
        catalog, digest = load_catalog()
        payload = render_round(catalog, digest, args.round)
        atomic_create(output, payload)
        validate_prompt_reference(project_root, args.output, args.round)
    except (AdvisoryError, OSError) as error:
        print(f"PROJECT_ADVISORY_FAILED: {error}", file=sys.stderr)
        return 2
    print(json.dumps({"status": "PROJECT_ADVISORY_READY", "round": args.round,
                      "output": args.output, "sha256": hashlib.sha256(payload).hexdigest()},
                     sort_keys=True, separators=(",", ":")))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `project_readiness_gate/test_render_project_advisory.py`
```yaml
block_id: "PROJECT-START-READINESS-VALIDATOR:advisory-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite advisory catalog, deterministic rendering and fail-closed regression suite"
license: "LicenseRef-Workspace-Owner"
sha256: "25cb15e2d5ee4b41f52f57280db4b30eb599cd86367502c8a696608a6e4fbe93"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

from contextlib import redirect_stderr, redirect_stdout
from copy import deepcopy
import io
import json
from pathlib import Path
import tempfile
import unittest

import render_project_advisory as target


class AdvisoryTests(unittest.TestCase):
    def setUp(self):
        self.catalog, self.digest = target.load_catalog()
        self.temp = tempfile.TemporaryDirectory()
        self.root = Path(self.temp.name)

    def tearDown(self):
        self.temp.cleanup()

    def test_catalog_is_complete_and_officially_routed(self):
        target.validate_catalog(self.catalog)
        self.assertEqual({item["round"] for item in self.catalog["topics"]}, set(target.ROUNDS))
        self.assertEqual(len(self.catalog["topics"]), 8)
        self.assertTrue(all(item["source_ids"] for item in self.catalog["topics"]))

    def test_catalog_unknown_field_rejected(self):
        catalog = deepcopy(self.catalog)
        catalog["invented"] = True
        with self.assertRaisesRegex(target.AdvisoryError, "catalog keys must be exactly"):
            target.validate_catalog(catalog)

    def test_catalog_duplicate_json_key_rejected(self):
        path = self.root / target.CATALOG_FILE
        path.write_text('{"schema":"a","schema":"b"}', encoding="utf-8")
        with self.assertRaisesRegex(target.AdvisoryError, "duplicate JSON key"):
            target.load_catalog(path)

    def test_catalog_requires_exactly_one_topic_per_round(self):
        catalog = deepcopy(self.catalog)
        catalog["topics"][-1]["round"] = "A"
        with self.assertRaisesRegex(target.AdvisoryError, "exactly one topic per round"):
            target.validate_catalog(catalog)

    def test_catalog_rejects_unknown_source(self):
        catalog = deepcopy(self.catalog)
        catalog["topics"][0]["source_ids"] = ["INVENTED-SOURCE"]
        with self.assertRaisesRegex(target.AdvisoryError, "unknown"):
            target.validate_catalog(catalog)

    def test_examples_cannot_be_presented_as_assumed_answers(self):
        catalog = deepcopy(self.catalog)
        catalog["topics"][0]["example_answer"] = "Argentina is the answer"
        with self.assertRaisesRegex(target.AdvisoryError, "not an assumed answer"):
            target.validate_catalog(catalog)

    def test_render_is_deterministic_and_explains_material_terms(self):
        first = target.render_round(self.catalog, self.digest, "D")
        second = target.render_round(self.catalog, self.digest, "D")
        self.assertEqual(first, second)
        text = first.decode("utf-8")
        for required in ("Corpus", "ground truth", "schema", "mapping", "NO_APLICA", self.digest):
            self.assertIn(required, text)

    def test_cli_creates_exact_prompt_and_validator_accepts_it(self):
        stdout = io.StringIO()
        with redirect_stdout(stdout):
            code = target.main(["--project-root", str(self.root), "--round", "C",
                                "--output", "PROJECT_ADVISORY_C.md"])
        self.assertEqual(code, 0)
        receipt = json.loads(stdout.getvalue())
        self.assertEqual(receipt["status"], "PROJECT_ADVISORY_READY")
        target.validate_prompt_reference(self.root, "PROJECT_ADVISORY_C.md", "C")

    def test_cli_refuses_overwrite(self):
        path = self.root / "PROJECT_ADVISORY_A.md"
        path.write_text("existing\n", encoding="utf-8")
        stderr = io.StringIO()
        with redirect_stderr(stderr):
            code = target.main(["--project-root", str(self.root), "--round", "A",
                                "--output", path.name])
        self.assertEqual(code, 2)
        self.assertIn("refusing to overwrite", stderr.getvalue())
        self.assertEqual(path.read_text(encoding="utf-8"), "existing\n")

    def test_tampered_prompt_rejected(self):
        path = self.root / "PROJECT_ADVISORY_H.md"
        path.write_bytes(target.render_round(self.catalog, self.digest, "H") + b"tampered\n")
        with self.assertRaisesRegex(target.AdvisoryError, "bytes do not match"):
            target.validate_prompt_reference(self.root, path.name, "H")

    def test_traversal_and_cross_round_reference_rejected(self):
        with self.assertRaisesRegex(target.AdvisoryError, "direct project-root"):
            target.validate_prompt_reference(self.root, "../PROJECT_ADVISORY_A.md", "A")
        (self.root / "PROJECT_ADVISORY_A.md").write_bytes(
            target.render_round(self.catalog, self.digest, "B")
        )
        with self.assertRaisesRegex(target.AdvisoryError, "bytes do not match"):
            target.validate_prompt_reference(self.root, "PROJECT_ADVISORY_A.md", "A")


if __name__ == "__main__":
    unittest.main()
````

### FILE: `project_readiness_gate/test_validate_project_readiness.py`
```yaml
block_id: "PROJECT-START-READINESS-VALIDATOR:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite readiness positive, negative, dependency and CLI regression suite"
license: "LicenseRef-Workspace-Owner"
sha256: "56b1bcdc2ed3fee71093c409bf2540519f6ebf9ce46422d9f5e0d9c26a2c299c"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

from copy import deepcopy
import json
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

import render_consumer_profile as renderer
import render_project_advisory as advisory
import validate_project_readiness as target


ESSENTIAL = {
    "PRD-INTAKE", "ARCH-DOMAIN", "REPO-SCM", "CONTRACTS", "RUNTIME", "DOMAIN-MODULES",
    "CI", "SUPPLY-CHAIN", "SECURITY-APPSEC", "TEST-PLATFORM", "DOCS-OPS", "WEB-PUBLIC",
}


class GateTests(unittest.TestCase):
    def operational_record(self, surface="CLI"):
        record = deepcopy(self.record)
        for item in record["capabilities"]:
            if item["id"] == "WEB-PUBLIC":
                item.update(classification="NONE_WITH_REASON", reason_or_trigger="No web in this CLI fixture")
            if item["id"] == "AUTHORIZATION":
                item.update(classification="REQUIRED", owner="operator", reason_or_trigger="Local action authority",
                            implementation_refs=["evidence.md"], test_refs=["evidence.md"],
                            rollback_recovery_ref="evidence.md")
        journey = record["connected_journeys"][0]
        journey["surfaces"] = [surface]
        journey["capability_ids"] = sorted((ESSENTIAL - {"WEB-PUBLIC"}) | {"AUTHORIZATION"})
        return record

    def test_cli_only_complete_journey(self):
        self.assertEqual(target.validate_record(self.operational_record(), self.root), [])

    def test_automation_only_complete_journey(self):
        self.assertEqual(target.validate_record(self.operational_record("AUTOMATION"), self.root), [])

    def test_operational_entry_cannot_remove_authorization(self):
        record = self.operational_record()
        for item in record["capabilities"]:
            if item["id"] == "AUTHORIZATION":
                item["classification"] = "NONE_WITH_REASON"
        record["connected_journeys"][0]["capability_ids"].remove("AUTHORIZATION")
        self.assertTrue(target.validate_record(record, self.root))

    def test_operational_entry_requires_selected_journey(self):
        record = self.operational_record()
        record["first_vertical_slice"]["connected_journey_ids"] = ["not-selected"]
        self.assertTrue(target.validate_record(record, self.root))

    def test_operational_entry_preserves_every_evidence_boundary(self):
        for field in target.CONNECTED_JOURNEY_REFERENCE_FIELDS:
            with self.subTest(field=field):
                record = self.operational_record()
                record["connected_journeys"][0][field] = []
                self.assertTrue(target.validate_record(record, self.root))

    def test_operational_entry_cannot_hide_web_surface(self):
        record = self.operational_record()
        record["connected_journeys"][0]["surfaces"].append("WEB_PUBLIC")
        self.assertTrue(target.validate_record(record, self.root))

    def test_operational_entry_requires_version_consistency(self):
        record = self.operational_record()
        record["connected_journeys"][0]["support_version"] = "2.0.0"
        self.assertTrue(target.validate_record(record, self.root))

    def test_operational_entry_requires_proven_status(self):
        record = self.operational_record()
        record["connected_journeys"][0]["status"] = "PLANNED"
        self.assertTrue(target.validate_record(record, self.root))

    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.root = Path(self.temp.name)
        template_path = Path(__file__).with_name("project-readiness.template.json")
        self.record = json.loads(template_path.read_text(encoding="utf-8"))
        for relative in target.REQUIRED_FILES:
            path = self.root / relative
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_text(f"evidence for {relative}\n", encoding="utf-8")
        spec = self.root / "specs" / "p1" / "spec.md"
        spec.parent.mkdir(parents=True)
        spec.write_text("P1 acceptance specification\n", encoding="utf-8")
        (self.root / "evidence.md").write_text("verified owner evidence\n", encoding="utf-8")
        project = self.record["project"]
        project.update({"id": "demo", "owner": "owner", "status": "READY_TO_BUILD",
                        "platform_mode": "CUSTOM_PLATFORM", "risk_tier": "STANDARD",
                        "data_classification": "INTERNAL", "updated_at": "2026-08-28T12:00:00Z"})
        advisory_catalog, advisory_digest = advisory.load_catalog()
        for round_id, item in self.record["rounds"].items():
            advisory_name = advisory.canonical_prompt_name(round_id)
            (self.root / advisory_name).write_bytes(
                advisory.render_round(advisory_catalog, advisory_digest, round_id)
            )
            item.update({"status": "ANSWERED", "evidence": ["PROJECT_READINESS_RECORD.md"],
                         "advisory_prompt_ref": advisory_name})
        for key in self.record["blockers"]:
            self.record["blockers"][key] = []
        self.record["first_vertical_slice"] = {
            "id": "P1", "journeys": ["visitor sees catalog"], "invariants": ["published only"],
            "connected_journey_ids": ["visitor-catalog"],
            "acceptance_test_refs": ["specs/p1/spec.md"], "rollback_recovery_ref": "evidence.md",
        }
        for item in self.record["capabilities"]:
            item.update({"classification": "NONE_WITH_REASON", "reason_or_trigger": "not required for P1",
                         "blocking_conditions": []})
            if item["id"] in ESSENTIAL:
                item.update({"classification": "REQUIRED", "owner": "team",
                             "reason_or_trigger": "required by P1",
                             "implementation_refs": ["PROJECT_PACK_PLAN.md"],
                             "test_refs": ["specs/p1/spec.md"], "rollback_recovery_ref": "evidence.md"})
        self.record["connected_journeys"] = [{
            "id": "visitor-catalog", "status": "PROVEN", "release_version": "1.0.0",
            "help_version": "1.0.0", "training_version": "1.0.0", "support_version": "1.0.0",
            "personas": ["visitor"], "roles": ["anonymous"], "surfaces": ["WEB_PUBLIC"],
            "capability_ids": sorted(ESSENTIAL),
            **{field: ["evidence.md"] for field in target.CONNECTED_JOURNEY_REFERENCE_FIELDS},
        }]
        self.record["documents"] = {"required": False, "decision_record": "",
                                      "corpus_ground_truth_status": "NOT_APPLICABLE",
                                      "automatic_storage_requested": False,
                                      "field_evaluation_evidence": [], "security_evidence": [],
                                      "classes": []}
        self.record["integrations"] = []
        self.record["sources"] = {"lock_path": "PROJECT_EXTERNAL_SOURCE_LOCK.md", "status": "PROVEN",
                                    "exact_revisions": True, "licenses_notices_resolved": True,
                                    "moving_references": [], "product_provenance_modes": ["AUTHORED"]}
        self.record["dependencies"] = {"status": "PROVEN", "evidence": ["PROJECT_DEPENDENCY_UPDATE_RECORD.md"],
                                         "eol_policy": True, "locks": True, "rollback": True}
        self.record["authorities"] = {"status": "PROVEN", "evidence": ["PROJECT_AUTHORITY_FRESHNESS_RECORD.md"],
                                        "conflicts": []}
        self.record["operations"] = {
            name: {"status": "PROVEN", "evidence": ["evidence.md"]} for name in (
                "security_privacy", "slo_performance_capacity_cost", "backup_restore_dr",
                "deployment_canary_rollback", "monitoring_incident_response")
        }
        self.record["pack_plan"] = {"path": "PROJECT_PACK_PLAN.md", "status": "PROVEN",
                                      "collision_check": True, "rollback_defined": True,
                                      "evidence": ["PROJECT_PACK_PLAN.md"]}

    def tearDown(self):
        self.temp.cleanup()

    def errors(self, record=None):
        return target.validate_record(record or self.record, self.root)

    def assert_error(self, fragment, record=None):
        errors = self.errors(record)
        self.assertTrue(any(fragment in item for item in errors), errors)

    def capability(self, capability_id):
        return next(item for item in self.record["capabilities"] if item["id"] == capability_id)

    def require(self, *ids):
        for capability_id in ids:
            item = self.capability(capability_id)
            item.update({"classification": "REQUIRED", "owner": "team", "reason_or_trigger": "required",
                         "implementation_refs": ["PROJECT_PACK_PLAN.md"], "test_refs": ["specs/p1/spec.md"],
                         "rollback_recovery_ref": "evidence.md", "blocking_conditions": []})
            capability_ids = self.record["connected_journeys"][0]["capability_ids"]
            if capability_id not in capability_ids:
                capability_ids.append(capability_id)

    def closed_consumer_template(self):
        template = {field: "" for field in renderer.STRING_FIELDS | target.MAPPING_FIELDS}
        template.update({"schemaVersion": 1, "acknowledgeConditionedState": False,
                         "owners": [], "approvals": [],
                         "proofs": {field: False for field in renderer.PROOF_FIELDS}})
        return template

    def test_ready_record_passes(self):
        self.assertEqual(self.errors(), [])

    def test_template_fails_closed(self):
        template = json.loads(Path(__file__).with_name("project-readiness.template.json").read_text(encoding="utf-8"))
        self.assertGreater(len(self.errors(template)), 20)

    def test_missing_required_artifact_rejected(self):
        (self.root / "PROJECT_BLUEPRINT.md").unlink()
        self.assert_error("PROJECT_BLUEPRINT.md")

    def test_open_blocker_rejected(self):
        self.record["blockers"]["critical_unknowns"] = ["tax jurisdiction"]
        self.assert_error("critical_unknowns must be empty")

    def test_pending_round_rejected(self):
        self.record["rounds"]["D"]["status"] = "PENDING"
        self.assert_error("rounds.D.status")

    def test_answered_round_requires_exact_advisory_prompt(self):
        self.record["rounds"]["D"]["advisory_prompt_ref"] = ""
        self.assert_error("rounds.D.advisory_prompt_ref")

    def test_tampered_advisory_prompt_rejected(self):
        (self.root / "PROJECT_ADVISORY_F.md").write_text("invented prompt\n", encoding="utf-8")
        self.assert_error("rounds.F.advisory_prompt_ref: advisory prompt bytes do not match")

    def test_cross_round_advisory_prompt_rejected(self):
        self.record["rounds"]["B"]["advisory_prompt_ref"] = "PROJECT_ADVISORY_A.md"
        self.assert_error("rounds.B.advisory_prompt_ref: advisory prompt reference must be PROJECT_ADVISORY_B.md")

    def test_missing_and_duplicate_capability_rejected(self):
        self.record["capabilities"][-1] = deepcopy(self.record["capabilities"][0])
        self.assert_error("exactly 48 unique IDs")

    def test_essential_capability_cannot_be_none(self):
        self.capability("CI")["classification"] = "NONE_WITH_REASON"
        self.assert_error("capabilities.CI must be REQUIRED")

    def test_required_capability_needs_evidence(self):
        self.capability("RUNTIME")["implementation_refs"] = []
        self.assert_error("capabilities.RUNTIME.implementation_refs")

    def test_connected_journeys_are_mandatory(self):
        self.record["connected_journeys"] = []
        self.assert_error("connected_journeys must be a non-empty array")

    def test_every_required_capability_must_be_connected(self):
        self.capability("TX-DATABASE").update({
            "classification": "REQUIRED", "owner": "team", "reason_or_trigger": "required",
            "implementation_refs": ["PROJECT_PACK_PLAN.md"], "test_refs": ["specs/p1/spec.md"],
            "rollback_recovery_ref": "evidence.md", "blocking_conditions": [],
        })
        self.assert_error("connected_journeys must cover every REQUIRED capability")

    def test_connected_journey_cannot_reference_nonrequired_capability(self):
        self.record["connected_journeys"][0]["capability_ids"].append("MOBILE")
        self.assert_error("capabilities.MOBILE=REQUIRED")

    def test_connected_journey_versions_are_release_bound(self):
        self.record["connected_journeys"][0]["training_version"] = "1.0.1"
        self.assert_error("training and support versions must equal release_version")

    def test_connected_journey_requires_every_evidence_link(self):
        self.record["connected_journeys"][0]["contextual_help_refs"] = []
        self.assert_error("contextual_help_refs: at least one evidence path is required")

    def test_connected_journey_schema_is_closed(self):
        self.record["connected_journeys"][0]["unreviewed_shortcut"] = True
        self.assert_error("must contain exactly the connected journey fields")

    def test_connected_journey_ids_are_unique(self):
        self.record["connected_journeys"].append(deepcopy(self.record["connected_journeys"][0]))
        self.assert_error("connected_journeys IDs must be unique")

    def test_first_slice_requires_known_connected_journey(self):
        self.record["first_vertical_slice"]["connected_journey_ids"] = ["missing-journey"]
        self.assert_error("contains unknown journeys")

    def test_dependency_rule_rejected(self):
        self.require("TX-DATABASE")
        self.assert_error("capabilities.TX-DATABASE requires capabilities.BACKUP-DR")

    def test_dependency_rule_passes_when_closed(self):
        self.require("TX-DATABASE", "BACKUP-DR", "OBSERVABILITY", "PRIVACY-COMPLIANCE")
        self.assertEqual(self.errors(), [])

    def test_documents_require_corpus_and_security(self):
        self.record["documents"].update({"required": True, "decision_record": "evidence.md",
                                          "corpus_ground_truth_status": "BLOCKED"})
        self.assert_error("corpus_ground_truth_status")
        self.assert_error("security_evidence")

    def test_documents_storage_requires_field_evaluation(self):
        self.require("DATA-INGEST")
        self.record["documents"].update({"required": True, "decision_record": "evidence.md",
                                          "corpus_ground_truth_status": "PROVEN",
                                          "security_evidence": ["evidence.md"],
                                          "automatic_storage_requested": True})
        self.assert_error("field_evaluation_evidence")

    def test_required_documents_require_exact_class_inventory(self):
        self.require("DATA-INGEST")
        self.record["documents"].update({"required": True, "decision_record": "evidence.md",
                                          "corpus_ground_truth_status": "PROVEN",
                                          "security_evidence": ["evidence.md"]})
        self.assert_error("documents.classes must contain at least one")

    def test_automatic_storage_requires_exact_mapping_authority(self):
        self.require("DATA-INGEST", "TX-DATABASE", "OBSERVABILITY", "BACKUP-DR", "PRIVACY-COMPLIANCE",
                     "OUTBOX-INBOX", "BROKER-STREAMING")
        expression = "{'invoice_number': doc.invoice_number, 'total': doc.total}"
        digest = __import__("hashlib").sha256(expression.encode()).hexdigest()
        schema_digest = "a" * 64
        self.record["documents"].update({
            "required": True, "decision_record": "evidence.md", "corpus_ground_truth_status": "PROVEN",
            "automatic_storage_requested": True, "field_evaluation_evidence": ["evidence.md"],
            "security_evidence": ["evidence.md"], "classes": [{
                "class_id": "invoice", "variant_id": "supplier-v1", "decision": "READY_FOR_AUTOMATIC_STORAGE",
                "owner": "finance", "inventory_evidence": ["evidence.md"],
                "automatic_storage_authorized": True, "persistence_mode": "OUTBOX_CDC_CONSUMER",
                "mapping": {"mappingId": "invoice-to-erp", "mappingDomainType": "invoice",
                            "mappingVersion": "1.0.0", "mappingExpression": expression,
                            "mappingExpressionSha256": digest,
                            "mappingExactOutputKeysCsv": "invoice_number,total",
                            "mappingInputSchemaId": "invoice-review-v1",
                            "mappingInputSchemaSha256": schema_digest,
                            "mappingOutputSchemaId": "erp-invoice-v1",
                            "mappingOutputSchemaSha256": schema_digest},
                "corpus_evidence": ["evidence.md"], "evaluation_evidence": ["evidence.md"],
                "approval_evidence": ["evidence.md"],
            }],
        })
        self.assertEqual(self.errors(), [])

    def test_automatic_storage_rejects_cross_class_or_unhashed_mapping(self):
        self.require("DATA-INGEST")
        self.record["documents"].update({
            "required": True, "decision_record": "evidence.md", "corpus_ground_truth_status": "PROVEN",
            "automatic_storage_requested": True, "field_evaluation_evidence": ["evidence.md"],
            "security_evidence": ["evidence.md"], "classes": [{
                "class_id": "Packing_List", "variant_id": "default", "decision": "LIMITED_AUTOMATION",
                "owner": "supply", "inventory_evidence": ["evidence.md"],
                "automatic_storage_authorized": True, "persistence_mode": "OUTBOX_CDC_CONSUMER",
                "mapping": {"mappingId": "packing-list-to-erp", "mappingDomainType": "invoice",
                            "mappingVersion": "1.0.0", "mappingExpression": "{'items': doc.items}",
                            "mappingExpressionSha256": "NOT-A-SHA", "mappingExactOutputKeysCsv": "items",
                            "mappingInputSchemaId": "packing-list-review-v1",
                            "mappingInputSchemaSha256": "B" * 64,
                            "mappingOutputSchemaId": "erp-packing-list-v1",
                            "mappingOutputSchemaSha256": "c" * 64},
                "corpus_evidence": ["evidence.md"], "evaluation_evidence": ["evidence.md"],
                "approval_evidence": ["evidence.md"],
            }],
        })
        self.assert_error("class_id must be a canonical")
        self.assert_error("decision must be READY_FOR_AUTOMATIC_STORAGE")
        self.assert_error("mappingDomainType must equal class_id")
        self.assert_error("mappingExpressionSha256 must match exact UTF-8 expression bytes")
        self.assert_error("mappingExpressionSha256 must be lowercase SHA-256")
        self.assert_error("mappingInputSchemaSha256 must be lowercase SHA-256")

    def test_document_class_variant_pairs_are_unique(self):
        self.require("DATA-INGEST")
        item = {"class_id": "invoice", "variant_id": "supplier-v1", "decision": "REVIEW_ONLY",
                "owner": "finance", "inventory_evidence": ["evidence.md"],
                "automatic_storage_authorized": False}
        self.record["documents"].update({"required": True, "decision_record": "evidence.md",
                                          "corpus_ground_truth_status": "PROVEN",
                                          "security_evidence": ["evidence.md"],
                                          "classes": [deepcopy(item), deepcopy(item)]})
        self.assert_error("class_id/variant_id pairs must be unique")

    def test_mapping_expression_must_be_single_line_and_hash_exact(self):
        self.test_automatic_storage_requires_exact_mapping_authority()
        item = self.record["documents"]["classes"][0]
        item["mapping"]["mappingExpression"] += "\n"
        self.assert_error("mappingExpression must be trimmed single-line")
        self.assert_error("mappingExpressionSha256 must match exact UTF-8 expression bytes")

    def test_mapping_output_keys_must_be_sorted_unique_canonical_csv(self):
        self.test_automatic_storage_requires_exact_mapping_authority()
        self.record["documents"]["classes"][0]["mapping"]["mappingExactOutputKeysCsv"] = "total,invoice_number,total"
        self.assert_error("mappingExactOutputKeysCsv must be sorted unique canonical CSV")

    def test_mapping_identity_must_match_consumer_contract(self):
        self.test_automatic_storage_requires_exact_mapping_authority()
        mapping = self.record["documents"]["classes"][0]["mapping"]
        mapping.update({"mappingId": "x", "mappingVersion": "0.1.0",
                        "mappingInputSchemaId": "x", "mappingOutputSchemaId": "Y"})
        self.assert_error("mappingId must match the consumer mapping ID contract")
        self.assert_error("mappingVersion must match the consumer semantic version contract")
        self.assert_error("mappingInputSchemaId must match the consumer schema ID contract")
        self.assert_error("mappingOutputSchemaId must match the consumer schema ID contract")

    def test_outbox_mapping_requires_runtime_capabilities(self):
        self.test_automatic_storage_requires_exact_mapping_authority()
        self.capability("BROKER-STREAMING")["classification"] = "NONE_WITH_REASON"
        self.assert_error("persistence_mode requires capabilities.BROKER-STREAMING=REQUIRED")

    def test_transactional_direct_does_not_require_broker_or_outbox(self):
        self.test_automatic_storage_requires_exact_mapping_authority()
        self.record["documents"]["classes"][0]["persistence_mode"] = "TRANSACTIONAL_DIRECT"
        for capability_id in ("BROKER-STREAMING", "OUTBOX-INBOX"):
            self.capability(capability_id).update({"classification": "NONE_WITH_REASON",
                                                   "owner": "", "reason_or_trigger": "not used by direct mode",
                                                   "implementation_refs": [], "test_refs": [],
                                                   "rollback_recovery_ref": ""})
            self.record["connected_journeys"][0]["capability_ids"].remove(capability_id)
        self.assertEqual(self.errors(), [])

    def test_mapping_contract_rejects_unknown_fields(self):
        self.test_automatic_storage_requires_exact_mapping_authority()
        self.record["documents"]["classes"][0]["mapping"]["untrustedOverride"] = True
        self.assert_error("must contain exactly the consumer mapping fields")

    def test_review_only_rejects_dormant_persistence_configuration(self):
        self.require("DATA-INGEST")
        self.record["documents"].update({"required": True, "decision_record": "evidence.md",
                                          "corpus_ground_truth_status": "PROVEN",
                                          "security_evidence": ["evidence.md"], "classes": [{
            "class_id": "invoice", "variant_id": "supplier-v1", "decision": "REVIEW_ONLY",
            "owner": "finance", "inventory_evidence": ["evidence.md"],
            "automatic_storage_authorized": False, "persistence_mode": "OUTBOX_CDC_CONSUMER",
            "mapping": {"mappingId": "dormant"},
        }]})
        self.assert_error("cannot carry persistence configuration")

    def test_required_first_slice_integration_needs_probe(self):
        self.record["integrations"] = [{"id": "payments", "requirement": "REQUIRED",
                                         "in_first_slice": True, "status": "NOT_PROVIDED",
                                         "evidence": [], "owner": "finance", "defer_reason": ""}]
        self.assert_error("must be PROVEN")

    def test_required_later_integration_may_be_explicitly_deferred(self):
        self.record["integrations"] = [{"id": "ads", "requirement": "REQUIRED",
                                         "in_first_slice": False, "status": "DEFERRED", "evidence": [],
                                         "owner": "marketing", "defer_reason": "outside P1"}]
        self.assertEqual(self.errors(), [])

    def test_official_platform_rejects_authored_product(self):
        self.record["project"]["platform_mode"] = "OFFICIAL_PLATFORM"
        self.assert_error("OFFICIAL_PLATFORM permits only")

    def test_official_platform_accepts_dependency_pin(self):
        self.record["project"]["platform_mode"] = "OFFICIAL_PLATFORM"
        self.record["sources"]["product_provenance_modes"] = ["VERBATIM", "DEPENDENCY_PIN"]
        self.assertEqual(self.errors(), [])

    def test_moving_source_rejected(self):
        self.record["sources"]["moving_references"] = ["main"]
        self.assert_error("moving_references must be empty")

    def test_secret_literal_key_rejected(self):
        self.record["integrations"] = [{"id": "bad", "requirement": "OPTIONAL", "in_first_slice": False,
                                         "status": "DEFERRED", "owner": "team", "defer_reason": "later",
                                         "api_token": "literal-secret", "evidence": []}]
        self.assert_error("secret material is prohibited")

    def test_evidence_traversal_rejected(self):
        self.record["rounds"]["A"]["evidence"] = ["../outside.md"]
        self.assert_error("project-relative")

    def test_invalid_timestamp_rejected(self):
        self.record["project"]["updated_at"] = "2026-02-30T12:00:00Z"
        self.assert_error("real ISO-8601")

    def test_cli_ready_writes_atomic_report(self):
        record_path = self.root / "PROJECT_READINESS_GATE.json"
        record_path.write_text(json.dumps(self.record), encoding="utf-8")
        result = subprocess.run(
            [sys.executable, str(Path(target.__file__)), "--project-root", str(self.root),
             "--report", "PROJECT_READINESS_REPORT.json"], capture_output=True, text=True, check=False)
        self.assertEqual(result.returncode, 0, result.stderr + result.stdout)
        report = json.loads((self.root / "PROJECT_READINESS_REPORT.json").read_text(encoding="utf-8"))
        self.assertEqual(report["status"], "READY_TO_BUILD")
        self.assertEqual(report["errors"], [])
        self.assertEqual(list(self.root.glob(".PROJECT_READINESS_REPORT.json.*")), [])

    def test_cli_blocked_returns_two(self):
        self.record["blockers"]["critical_unknowns"] = ["owner decision"]
        (self.root / "PROJECT_READINESS_GATE.json").write_text(json.dumps(self.record), encoding="utf-8")
        result = subprocess.run(
            [sys.executable, str(Path(target.__file__)), "--project-root", str(self.root)],
            capture_output=True, text=True, check=False)
        self.assertEqual(result.returncode, 2, result.stderr + result.stdout)
        self.assertEqual(json.loads(result.stdout)["status"], "BLOCKED")

    def test_renderer_copies_exact_mapping_and_keeps_profile_closed(self):
        self.test_automatic_storage_requires_exact_mapping_authority()
        profile = renderer.render_profile(self.record, self.root, "invoice", "supplier-v1",
                                          self.closed_consumer_template())
        mapping = self.record["documents"]["classes"][0]["mapping"]
        self.assertEqual({field: profile[field] for field in target.MAPPING_FIELDS}, mapping)
        self.assertIs(profile["acknowledgeConditionedState"], False)
        self.assertTrue(all(value is False for value in profile["proofs"].values()))
        self.assertEqual(profile["databasePasswordSecretRef"], "")

    def test_renderer_rejects_blocked_readiness(self):
        self.test_automatic_storage_requires_exact_mapping_authority()
        self.record["blockers"]["critical_unknowns"] = ["unresolved"]
        with self.assertRaisesRegex(renderer.RenderError, "readiness record is BLOCKED"):
            renderer.render_profile(self.record, self.root, "invoice", "supplier-v1",
                                    self.closed_consumer_template())

    def test_renderer_requires_exact_class_and_variant(self):
        self.test_automatic_storage_requires_exact_mapping_authority()
        with self.assertRaisesRegex(renderer.RenderError, "select exactly one"):
            renderer.render_profile(self.record, self.root, "invoice", "customer-v1",
                                    self.closed_consumer_template())

    def test_renderer_rejects_direct_persistence_mode(self):
        self.test_automatic_storage_requires_exact_mapping_authority()
        self.record["documents"]["classes"][0]["persistence_mode"] = "TRANSACTIONAL_DIRECT"
        with self.assertRaisesRegex(renderer.RenderError, "OUTBOX_CDC_CONSUMER"):
            renderer.render_profile(self.record, self.root, "invoice", "supplier-v1",
                                    self.closed_consumer_template())

    def test_renderer_rejects_tampered_or_preconfigured_template(self):
        self.test_automatic_storage_requires_exact_mapping_authority()
        template = self.closed_consumer_template()
        template["mappingExpression"] = "unreviewed"
        with self.assertRaisesRegex(renderer.RenderError, "must start empty"):
            renderer.render_profile(self.record, self.root, "invoice", "supplier-v1", template)
        template = self.closed_consumer_template()
        template["proofs"]["rollbackProven"] = True
        with self.assertRaisesRegex(renderer.RenderError, "exact closed proof set"):
            renderer.render_profile(self.record, self.root, "invoice", "supplier-v1", template)

    def test_renderer_rejects_duplicate_json_keys(self):
        duplicate = self.root / "duplicate.json"
        duplicate.write_text('{"schemaVersion":1,"schemaVersion":2}', encoding="utf-8")
        with self.assertRaisesRegex(renderer.RenderError, "duplicate JSON key"):
            renderer.load_json(duplicate, "consumer template")

    def test_renderer_cli_is_atomic_and_refuses_overwrite(self):
        self.test_automatic_storage_requires_exact_mapping_authority()
        (self.root / "PROJECT_READINESS_GATE.json").write_text(json.dumps(self.record), encoding="utf-8")
        (self.root / "consumer.template.json").write_text(
            json.dumps(self.closed_consumer_template()), encoding="utf-8")
        command = [sys.executable, str(Path(renderer.__file__)), "--project-root", str(self.root),
                   "--consumer-template", "consumer.template.json", "--class-id", "invoice",
                   "--variant-id", "supplier-v1", "--output", "consumer.profile.json"]
        first = subprocess.run(command, capture_output=True, text=True, check=False)
        self.assertEqual(first.returncode, 0, first.stderr + first.stdout)
        self.assertIn("CONSUMER_PROFILE_RENDER_PASS", first.stdout)
        self.assertEqual(list(self.root.glob(".consumer.profile.json.*")), [])
        second = subprocess.run(command, capture_output=True, text=True, check=False)
        self.assertEqual(second.returncode, 2, second.stderr + second.stdout)
        self.assertIn("output path already exists", second.stderr)


    def assert_cli_blocked_unchanged(self, args):
        before = {p.relative_to(self.root).as_posix(): p.read_bytes() for p in self.root.rglob('*') if p.is_file()}
        result = subprocess.run([sys.executable, '-X', 'utf8', '-B', str(Path(target.__file__)), *map(str, args)], capture_output=True, timeout=15)
        self.assertEqual(result.returncode, 2, result.stderr.decode('utf-8', errors='replace'))
        self.assertIn(b'PROJECT_READINESS_BLOCKED:', result.stderr)
        self.assertNotIn(b'Traceback', result.stderr)
        self.assertEqual(result.stdout, b'')
        self.assertEqual(before, {p.relative_to(self.root).as_posix(): p.read_bytes() for p in self.root.rglob('*') if p.is_file()})

    def test_cli_invalid_inputs_are_controlled(self):
        record = self.root / 'PROJECT_READINESS_GATE.json'
        record.write_text(json.dumps(self.record), encoding='utf-8')
        malformed = self.root / 'malformed.json'
        malformed.write_bytes(b'{')
        invalid_utf8 = self.root / 'invalid-utf8.json'
        invalid_utf8.write_bytes(b'\xff')
        directory = self.root / 'record-directory'
        directory.mkdir()
        for label, project, source in [
            ('missing-root', self.root/'absent', record.name),
            ('file-root', record, record.name),
            ('missing-record', self.root, 'absent.json'),
            ('directory-record', self.root, directory.name),
            ('malformed-json', self.root, malformed.name),
            ('invalid-utf8', self.root, invalid_utf8.name),
        ]:
            with self.subTest(case=label):
                self.assert_cli_blocked_unchanged(['--project-root', project, '--record', source, '--report', 'new-report.json'])

    def test_cli_report_errors_preserve_files(self):
        (self.root/'PROJECT_READINESS_GATE.json').write_text(json.dumps(self.record), encoding='utf-8')
        owned = self.root/'owned-report.json'
        owned.write_bytes(b'\xef\xbb\xbfUser owned report\r\n')
        parent_file = self.root/'not-a-directory'
        parent_file.write_bytes(b'original parent file')
        with tempfile.TemporaryDirectory() as separate:
            escaped = Path(separate)/'report.json'
            for label, destination in [('existing', owned), ('outside-root', escaped), ('file-parent', parent_file/'report.json')]:
                with self.subTest(case=label):
                    self.assert_cli_blocked_unchanged(['--project-root', self.root, '--report', destination])
                    self.assertEqual(list(Path(separate).iterdir()), [])

    def test_atomic_report_refuses_a_late_destination(self):
        destination = self.root/'late-report.json'
        original = b'\xef\xbb\xbfEvidence written by another caller\r\n'
        destination.write_bytes(original)
        with self.assertRaises(FileExistsError):
            target.atomic_write(destination, b'new report')
        self.assertEqual(destination.read_bytes(), original)
        self.assertEqual(list(self.root.glob('.late-report.json.*')), [])

    def test_atomic_report_concurrent_publish_has_one_winner(self):
        program = '''import pathlib,sys,time
import validate_project_readiness as target
gate=pathlib.Path(sys.argv[1]); destination=pathlib.Path(sys.argv[2]); payload=sys.argv[3].encode()
deadline=time.monotonic()+10
while not gate.exists():
    if time.monotonic()>deadline: raise SystemExit(3)
    time.sleep(0.01)
try: target.atomic_write(destination,payload)
except FileExistsError: raise SystemExit(2)
'''
        for round_id in range(3):
            gate=self.root/f'go-{round_id}'
            destination=self.root/f'concurrent-{round_id}.json'
            payloads=['first-complete-payload','second-complete-payload']
            children=[]
            try:
                for payload in payloads:
                    children.append(subprocess.Popen([sys.executable,'-X','utf8','-B','-c',program,str(gate),str(destination),payload],cwd=Path(target.__file__).parent,stdout=subprocess.PIPE,stderr=subprocess.PIPE))
                gate.write_bytes(b'go')
                outcomes=[]
                for child in children:
                    out,err=child.communicate(timeout=15)
                    outcomes.append(child.returncode)
                    self.assertEqual(out,b'')
                    self.assertEqual(err,b'')
                self.assertEqual(sorted(outcomes),[0,2])
                self.assertEqual(destination.read_bytes(),payloads[outcomes.index(0)].encode())
                self.assertEqual(list(self.root.glob(f'.{destination.name}.*')),[])
            finally:
                for child in children:
                    if child.poll() is None:
                        child.kill()
                        child.communicate(timeout=5)

    def test_cli_report_io_errors_do_not_publish_readiness(self):
        import contextlib
        import io
        from unittest import mock
        (self.root/'PROJECT_READINESS_GATE.json').write_text(json.dumps(self.record), encoding='utf-8')
        before = {p.relative_to(self.root).as_posix(): p.read_bytes() for p in self.root.rglob('*') if p.is_file()}
        for error in (PermissionError('synthetic denied'), OSError('synthetic I/O error')):
            with self.subTest(error=type(error).__name__):
                output, diagnostic = io.StringIO(), io.StringIO()
                with mock.patch.object(target, 'atomic_write', side_effect=error), contextlib.redirect_stdout(output), contextlib.redirect_stderr(diagnostic):
                    result = target.main(['--project-root', str(self.root), '--report', 'new-report.json'])
                self.assertEqual(result, 2)
                self.assertEqual(output.getvalue(), '')
                self.assertIn('PROJECT_READINESS_BLOCKED:', diagnostic.getvalue())
                self.assertEqual(before, {p.relative_to(self.root).as_posix(): p.read_bytes() for p in self.root.rglob('*') if p.is_file()})

if __name__ == "__main__":
    unittest.main()
````

### FILE: `project_readiness_gate/README.md`
```yaml
block_id: "PROJECT-START-READINESS-VALIDATOR:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite operating instructions and attribution boundary"
license: "LicenseRef-Workspace-Owner"
sha256: "90c69ee86499b450b1ca1c81f9c81cc278b06982d157eca8ac8e164b607b2b3a"
variables: []
secrets_allowed: false
```
````markdown
# Executable project readiness gate

This gate is Elite-authored orchestration. It operationalizes the project intake and capability contracts; it is not product code and is not attributed to GitHub, AWS, Microsoft or another upstream.

1. Copy `project-readiness.template.json` to the project root as `PROJECT_READINESS_GATE.json`.
2. Before each round, render its exact advisory prompt. Present it conversationally, record the real answer/evidence and set that round's `advisory_prompt_ref` to the generated file:

```powershell
python project_readiness_gate/render_project_advisory.py --project-root . --round A --output PROJECT_ADVISORY_A.md
```

3. Complete rounds A–H while maintaining the human-readable `PROJECT_READINESS_RECORD.md` and all required evidence artifacts. The prompt explains unfamiliar terms, required inputs, an explicitly non-assumed example, verification and the valid `NO_APLICA` boundary.
4. Never place credentials or private key material in prompts or records. Store only logical identity references in the appropriate project evidence.
5. Run the validator before composing any product pack:

```powershell
python project_readiness_gate/validate_project_readiness.py --project-root . --report PROJECT_READINESS_REPORT.json
```

Exit `0` and `status=READY_TO_BUILD` are required. Exit `2` is an expected fail-closed readiness decision. The validator checks all 48 capabilities, dependency implications, eight intake rounds, the exact catalog-generated advisory prompt for every answered round, required project files, local evidence paths, blockers, first vertical slice, documents/corpus, integrations/access, exact sources/licenses, dependency lifecycle, authority freshness, operations and pack-plan rollback. Every `REQUIRED` capability must also be covered by a `connected_journeys` record. Each record is release-bound and proves personas, roles, delivery surface, interface, authorization, API contract, domain and data/external effect, audit/observability, understandable errors, contextual help, training, support, update impact, rollback/recovery, E2E and negative tests. The first slice selects existing connected journey IDs; isolated screens, backend-only effects or unversioned operating content keep readiness blocked. Required documents need a unique `{class_id, variant_id}` inventory. Automatic storage additionally requires an explicitly authorized class, approved decision and persistence mode. Its mapping object contains exactly the Debezium/CEL consumer mapping fields—no unknown override—plus expression bytes, matching SHA-256, canonical output-key CSV, input/output schema IDs+SHA and corpus/evaluation/approval evidence. `TRANSACTIONAL_DIRECT` requires database/observability/backup/privacy; `OUTBOX_CDC_CONSUMER` additionally requires outbox and broker. Review-only classes cannot carry dormant persistence configuration.

After the complete readiness record passes, an authorized `OUTBOX_CDC_CONSUMER` mapping can be transferred without manual transcription into a freshly materialized, still-closed Debezium consumer template:

```powershell
python project_readiness_gate/render_consumer_profile.py --project-root . --consumer-template debezium_postgres_inbox_consumer/project-profile.template.json --class-id invoice --variant-id supplier-v1 --output PROJECT_DOCUMENT_CONSUMER_PROFILE.json
```

The renderer validates readiness again, selects one exact class/variant, rejects direct-mode mappings, duplicate JSON keys, altered/preconfigured templates, symlink/path escape and overwrite, then copies only the ten mapping fields. It never supplies endpoints, secret references, owners, approvals, acknowledgements or proofs. The generated profile therefore remains conditioned and cannot itself authorize deployment.

The advisory renderer does not invent or answer. It deterministically turns the curated catalog into one manageable prompt per round. The agent explains that prompt, lets the user answer in ordinary language, asks only the material follow-up, updates both records and reruns the validator. Modified, missing, cross-round or traversal prompt references fail closed; a syntactically complete record without evidence remains blocked.


## Explicit library infrastructure scope (0.7.0 candidate)

This mode implements the user's library acceptance contract. It is AUTHORED
validator glue; it does not originate from, or certify, any external company.
The existing default PROJECT validator and READY_TO_BUILD conditions are unchanged.

Copy `library-readiness.template.json` into a local maintenance record. All refs
are `{"path": "project-relative/file", "sha256": "exact 64 hex characters"}`;
unknown fields, missing/modified evidence, symlinks/reparse points and unsafe
paths fail closed. The template is deliberately blocked.

```powershell
python project_readiness_gate/validate_project_readiness.py --project-root . --record PROJECT_READINESS_GATE.json --scope LIBRARY_INFRASTRUCTURE --level preparation
python project_readiness_gate/validate_project_readiness.py --project-root . --record PROJECT_READINESS_GATE.json --scope LIBRARY_INFRASTRUCTURE --level release
```

Preparation requires the exact user scope decision, resolution plan, authority
decisions, implementation assurance, admitted tooling, source lock and notices;
all are bound to their bytes. T2801 must be PREPARATION_PROVEN with a real local
receipt. Remaining controls may be PENDING. READY_FOR_LIBRARY_WORK opens the
authorized local work; it does not close T2801, TEST02 or any other pending task.
Local blockers must have an exact category, affected controls and a hash-bound
resolution plan. They remain OPEN during preparation; this authorizes correcting
them. Every unresolved local blocker rejects release, and none can be silently
converted into missing credentials.

Release requires PROVEN_LOCAL for every T2801–T2810 and ARCA_INFRA control with
successful same-revision execution receipts bound to one final inventory. Each
receipt includes the command, exit code, verified assertions, input/output hashes
and a declaration that live effects were not performed. The full inventory must
equal two independent product trees outside the canonical root, including extra
or missing files. Every file has exactly one provenance assignment and a matching
G0–G8 admission record with hash-verified evidence. AUTHORED BUSINESS_LOGIC is
rejected. AUTHORED requires role GLUE, an unavoidable-glue justification and the
same admission evidence; a company label is not a provenance mode.

Only USER_CREDENTIALS and the exact DAYBREAK_LIBXML2 boundary may be deferred.
Credential deferral requires complete implementation plus successful local
contract evidence; it cannot replace a pending control at release. The Daybreak
entry is exactly id daybreak-libxml2 / control T2803 / code_complete false /
local_receipt null with an exclusion record and reopening trigger. It conveys
no approval of libxml2, other native dependencies or security findings.

The final success is READY_FOR_LIBRARY_USE with production_authorized false.
It never emits READY_TO_BUILD and cannot be fed to the PROJECT gate as approval.
Production credentials, policy/country/corpus choices, live validation and target
acceptance retain their own project workflow. No secrets belong in these records.

Evidence validation proves file identity and declared relationships, not the
truth of arbitrary prose or the semantic quality of a falsely authored receipt.
Actual runners must generate the receipts; the admitting agent must inspect the
source/provenance review and reject business logic mislabeled as glue. Hashing a
claim does not prove that claim. Unit test fixtures are explicitly synthetic and
never constitute actual library readiness evidence.
````

### FILE: `project_readiness_gate/validate_library_readiness.py`
```yaml
block_id: "PROJECT-START-READINESS-VALIDATOR:library-validator:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite local validator glue for the explicit V402 user library scope; no external code attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "80a77f41aea74fa43e03f83a11a211795312d21cf9439d8b36e14c5801b0ad8e"
variables: []
secrets_allowed: false
```
````python
"""AUTHORED validator glue. Local library evidence is never production approval."""
from __future__ import annotations

import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import re
from typing import Any

SCHEMA = "elite-library-infrastructure-readiness/v1"
CONTROLS = tuple(f"T28{i:02d}" for i in range(1, 11)) + ("ARCA_INFRA",)
SHA = re.compile(r"^[0-9a-f]{64}$")
FIELDS = {"schema", "scope", "revision", "production_authorized", "request",
          "plan", "authority_decisions", "assurance", "tooling", "source_lock",
          "notices", "blockers", "deferred", "controls", "reference"}
BASE_REFS = ("request", "plan", "authority_decisions", "assurance", "tooling",
             "source_lock", "notices")


def digest(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def exact(value: Any, keys: set[str], label: str, errors: list[str]) -> bool:
    if not isinstance(value, dict) or set(value) != keys:
        errors.append(f"{label}: exact fields required: {','.join(sorted(keys))}")
        return False
    return True


def text(value: Any) -> bool:
    return isinstance(value, str) and bool(value.strip())


def safe_file(root: Path, name: Any, *, allow_empty=False) -> Path:
    if (not isinstance(name, str) or not name or "\\" in name or ":" in name
            or "#" in name or "\x00" in name):
        raise ValueError("canonical relative file path required")
    path = PurePosixPath(name)
    if path.is_absolute() or any(p in {"..", ".", ""} for p in name.split("/")):
        raise ValueError("unsafe relative file path")
    current = root
    for part in path.parts:
        current = current / part
        info = current.lstat()
        if current.is_symlink() or getattr(info, "st_file_attributes", 0) & 0x400:
            raise ValueError("symlink/reparse evidence rejected")
    current.resolve(strict=True).relative_to(root.resolve(strict=True))
    if not current.is_file() or (not allow_empty and current.stat().st_size == 0):
        raise ValueError("non-empty regular file required")
    return current


def reference(root: Path, value: Any, label: str, errors: list[str], *, json_data=False):
    if not exact(value, {"path", "sha256"}, label, errors):
        return None
    try:
        if not isinstance(value["sha256"], str) or not SHA.fullmatch(value["sha256"]):
            raise ValueError("invalid SHA-256")
        payload = safe_file(root, value["path"]).read_bytes()
        if digest(payload) != value["sha256"]:
            raise ValueError("SHA-256 mismatch")
        return json.loads(payload) if json_data else payload
    except (OSError, ValueError, UnicodeError) as exc:
        errors.append(f"{label}: {exc}")
        return None


def receipt(root: Path, ref: Any, control: str, revision: str,
            inventory_sha: str | None, errors: list[str]) -> None:
    value = reference(root, ref, f"{control}.receipt", errors, json_data=True)
    keys = {"schema", "scope", "revision", "control", "outcome", "command",
            "exit_code", "assertions", "inputs", "outputs", "inventory_sha256",
            "live_effects_performed"}
    if value is None or not exact(value, keys, f"{control}.receipt", errors):
        return
    if (value["schema"] != "elite-library-control-receipt/v1"
            or value["scope"] != "LOCAL_FIXTURES" or value["revision"] != revision
            or value["control"] != control or value["outcome"] != "PASS"
            or type(value["exit_code"]) is not int or value["exit_code"] != 0
            or value["live_effects_performed"] is not False):
        errors.append(f"{control}: receipt must be a successful exact-revision local run")
    if not text(value["command"]):
        errors.append(f"{control}: reproducible command required")
    if not isinstance(value["assertions"], list) or not value["assertions"] or not all(map(text, value["assertions"])):
        errors.append(f"{control}: concrete verified assertions required")
    for name in ("inputs", "outputs"):
        entries = value[name]
        if not isinstance(entries, list) or not entries:
            errors.append(f"{control}: {name} must be a non-empty hash-bound array")
        else:
            for index, item in enumerate(entries):
                reference(root, item, f"{control}.{name}[{index}]", errors)
    if inventory_sha is not None and value["inventory_sha256"] != inventory_sha:
        errors.append(f"{control}: receipt is not bound to the final reference inventory")


def reference_trees(root: Path, value: Any, revision: str, errors: list[str]) -> str | None:
    if not exact(value, {"first_root", "second_root", "inventory"}, "reference", errors):
        return None
    manifest = reference(root, value["inventory"], "reference.inventory", errors, json_data=True)
    if manifest is None or not exact(manifest, {"schema", "revision", "files", "provenance"}, "inventory", errors):
        return None
    if manifest["schema"] != "elite-library-reference-inventory/v1" or manifest["revision"] != revision:
        errors.append("inventory: exact schema and revision required")
    files = manifest["files"]
    expected = {}
    if not isinstance(files, list) or not files:
        errors.append("inventory: a non-empty complete file inventory is required")
    else:
        for item in files:
            if not exact(item, {"path", "sha256"}, "inventory.file", errors):
                continue
            name, sha = item["path"], item["sha256"]
            if not text(name) or not isinstance(sha, str) or not SHA.fullmatch(sha):
                errors.append("inventory: invalid file identity")
            elif name.casefold() in {p.casefold() for p in expected}:
                errors.append("inventory: duplicate/case-colliding file")
            else:
                expected[name] = sha
    roots = []
    for field in ("first_root", "second_root"):
        try:
            path = Path(value[field])
            if not path.is_absolute():
                raise ValueError("reference root must be absolute")
            # Inspect ancestors before resolving; a link may otherwise disappear.
            for part in (path, *path.parents):
                if part.is_symlink() or getattr(part.lstat(), "st_file_attributes", 0) & 0x400:
                    raise ValueError("reference root traverses symlink/reparse point")
            path = path.resolve(strict=True)
            if not path.is_dir() or path == root or root in path.parents or path in root.parents:
                raise ValueError("reference product must be outside canonical root")
            roots.append(path)
            observed = {}
            for directory, directories, names in os.walk(path, followlinks=False):
                for name in directories:
                    child = Path(directory) / name
                    if child.is_symlink() or getattr(child.lstat(), "st_file_attributes", 0) & 0x400:
                        raise ValueError("reference contains symlink/reparse directory")
                for name in names:
                    child = Path(directory) / name
                    relative = child.relative_to(path).as_posix()
                    observed[relative] = digest(safe_file(path, relative, allow_empty=True).read_bytes())
            if observed != expected:
                errors.append(f"reference.{field}: tree differs from the exact inventory")
        except (OSError, ValueError, TypeError) as exc:
            errors.append(f"reference.{field}: {exc}")
    if len(roots) == 2 and (roots[0] == roots[1] or roots[0] in roots[1].parents or roots[1] in roots[0].parents):
        errors.append("reference: independent non-overlapping materialization destinations required")
    groups, covered = manifest["provenance"], set()
    if not isinstance(groups, list) or not groups:
        errors.append("inventory: provenance groups required for every file")
    else:
        for index, group in enumerate(groups):
            label = f"provenance[{index}]"
            if not exact(group, {"mode", "role", "files", "admission", "glue_justification"}, label, errors):
                continue
            mode = group["mode"]
            if not isinstance(mode, str) or mode not in {"VERBATIM", "DEPENDENCY_PIN", "ADAPTED", "AUTHORED"}:
                errors.append(f"{label}: declared provenance mode required")
            role = group["role"]
            if not isinstance(role, str) or role not in {"GLUE", "BUSINESS_LOGIC", "DEPENDENCY"}:
                errors.append(f"{label}: explicit implementation role required")
            if mode == "AUTHORED" and role != "GLUE":
                errors.append(f"{label}: AUTHORED business logic is outside this library acceptance scope")
            if mode == "AUTHORED" and not text(group["glue_justification"]):
                errors.append(f"{label}: AUTHORED requires unavoidable glue justification")
            selected = group["files"]
            if not isinstance(selected, list) or not selected or not all(isinstance(p, str) for p in selected):
                errors.append(f"{label}: exact files required")
                continue
            if len(set(selected)) != len(selected) or covered.intersection(selected):
                errors.append(f"{label}: duplicate provenance assignment")
            covered.update(selected)
            admission = reference(root, group["admission"], f"{label}.admission", errors, json_data=True)
            if admission is None or not exact(admission, {"schema", "revision", "mode", "role", "files", "source_lock", "gates"}, label, errors):
                continue
            if (admission["schema"] != "elite-library-admission/v1" or admission["revision"] != revision
                    or admission["mode"] != mode or admission["role"] != role or admission["files"] != selected):
                errors.append(f"{label}: exact admission scope required")
            reference(root, admission["source_lock"], f"{label}.source_lock", errors)
            if exact(admission["gates"], {f"G{i}" for i in range(9)}, f"{label}.gates", errors):
                for gate, result in admission["gates"].items():
                    if not exact(result, {"status", "evidence"}, f"{label}.{gate}", errors):
                        continue
                    if result["status"] != "PASS":
                        errors.append(f"{label}.{gate}: PASS required")
                    reference(root, result["evidence"], f"{label}.{gate}.evidence", errors)
    if covered != set(expected):
        errors.append("inventory: provenance coverage must equal the complete file set")
    return value["inventory"].get("sha256")


def validate_library_record(record: Any, root: Path, level: str) -> list[str]:
    errors: list[str] = []
    root = root.resolve(strict=True)
    if level not in {"preparation", "release"}:
        return ["library readiness level must be preparation or release"]
    if not exact(record, FIELDS, "library record", errors):
        return errors
    if record["schema"] != SCHEMA or record["scope"] != "LIBRARY_INFRASTRUCTURE":
        errors.append("library schema/scope must be explicit; project readiness cannot be inherited")
    if record["production_authorized"] is not False:
        errors.append("production_authorized must remain false")
    revision = record["revision"]
    if not text(revision):
        errors.append("exact non-empty revision required")
    for field in BASE_REFS:
        reference(root, record[field], field, errors)
    blockers = record["blockers"]
    if not isinstance(blockers, list):
        errors.append("blockers must be an explicit array")
    else:
        blocker_ids = set()
        for blocker in blockers:
            if not exact(blocker, {"id", "kind", "status", "affected_controls", "resolution_plan"}, "blocker", errors):
                continue
            if not text(blocker["id"]) or blocker["id"] in blocker_ids:
                errors.append("blocker: unique non-empty id required")
            else:
                blocker_ids.add(blocker["id"])
            if blocker["kind"] not in ("SOURCE_PROVENANCE", "DEPENDENCY_ADMISSION", "MISSING_IMPLEMENTATION", "MISSING_EVIDENCE") or blocker["status"] != "OPEN":
                errors.append("blocker: exact open resolution category required")
            affected = blocker["affected_controls"]
            if not isinstance(affected, list) or not affected or not all(c in CONTROLS for c in affected):
                errors.append("blocker: affected controls must be explicit")
            reference(root, blocker["resolution_plan"], "blocker.resolution_plan", errors)
        if level == "release" and blockers:
            errors.append("unresolved local blockers prevent release readiness")
    deferred = record["deferred"]
    ids = set()
    if not isinstance(deferred, list):
        errors.append("deferred must be an array")
    else:
        for item in deferred:
            if not exact(item, {"id", "kind", "control", "reason", "trigger", "evidence", "code_complete", "local_receipt"}, "deferred", errors):
                continue
            if not text(item["id"]):
                errors.append("deferred: unique non-empty id required")
                continue
            if item["id"] in ids:
                errors.append("deferred: unique non-empty id required")
            ids.add(item["id"])
            if not text(item["reason"]) or not text(item["trigger"]):
                errors.append("deferred: explicit reason and reopening trigger required")
            reference(root, item["evidence"], "deferred.evidence", errors)
            if item["kind"] == "DAYBREAK_LIBXML2":
                if item["id"] != "daybreak-libxml2" or item["control"] != "T2803" or item["code_complete"] is not False or item["local_receipt"] is not None:
                    errors.append("Daybreak: only exact excluded libxml2 boundary may be deferred")
            elif item["kind"] == "USER_CREDENTIALS":
                if item["control"] not in CONTROLS or item["code_complete"] is not True:
                    errors.append("credentials cannot defer missing implementation")
                receipt(root, item["local_receipt"], item["control"], revision, None, errors)
            else:
                errors.append("deferred: only USER_CREDENTIALS or DAYBREAK_LIBXML2 are allowed")
    inventory_sha = None
    if level == "release":
        inventory_sha = reference_trees(root, record["reference"], revision, errors)
    elif record["reference"] is not None:
        errors.append("preparation reference must be null; use release to verify materialized product")
    controls = record["controls"]
    if not exact(controls, set(CONTROLS), "controls", errors):
        return sorted(set(errors))
    for control, value in controls.items():
        if not exact(value, {"status", "receipt"}, control, errors):
            continue
        required = level == "release"
        if level == "preparation" and control == "T2801":
            if value["status"] != "PREPARATION_PROVEN":
                errors.append("T2801: PREPARATION_PROVEN required; this is not task completion")
            receipt(root, value["receipt"], control, revision, None, errors)
        elif value["status"] == "PROVEN_LOCAL":
            receipt(root, value["receipt"], control, revision, inventory_sha, errors)
        elif value["status"] != "PENDING" or required or value["receipt"] is not None:
            errors.append(f"{control}: {'PROVEN_LOCAL required' if required else 'invalid pending control'}")
    return sorted(set(errors))
````

### FILE: `project_readiness_gate/test_validate_library_readiness.py`
```yaml
block_id: "PROJECT-START-READINESS-VALIDATOR:library-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite local validator glue for the explicit V402 user library scope; no external code attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "a104fe2b7c0db52166ef1c2022b96531f7153f5f9007afe0d0438f87bc343a90"
variables: []
secrets_allowed: false
```
````python
"""Synthetic validator cases; a PASS here never accepts actual library controls."""
from __future__ import annotations
import copy
import json
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

import validate_library_readiness as gate


class LibraryReadiness(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.base = Path(self.temp.name)
        self.root = self.base / "library"
        self.root.mkdir()
        self.anchor = self.write("evidence/method.txt", b"Synthetic fixture only.\n")
        self.record = {"schema": gate.SCHEMA, "scope": "LIBRARY_INFRASTRUCTURE",
                       "revision": "fixture-v1", "production_authorized": False,
                       **{name: self.anchor for name in gate.BASE_REFS},
                       "blockers": [], "deferred": [], "reference": None,
                       "controls": {name: {"status": "PENDING", "receipt": None} for name in gate.CONTROLS}}
        self.record["controls"]["T2801"] = {"status": "PREPARATION_PROVEN", "receipt": self.run_receipt("T2801")}

    def tearDown(self):
        self.temp.cleanup()

    def write(self, name, value):
        if not isinstance(value, bytes):
            value = (json.dumps(value, sort_keys=True, indent=2) + "\n").encode()
        path = self.root / name
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_bytes(value)
        return {"path": name, "sha256": gate.digest(value)}

    def run_receipt(self, control, inventory=None):
        return self.write(f"evidence/{control}.json", {
            "schema": "elite-library-control-receipt/v1", "scope": "LOCAL_FIXTURES",
            "revision": "fixture-v1", "control": control, "outcome": "PASS",
            "command": "python -m unittest fixture", "exit_code": 0,
            "assertions": ["synthetic input bound to synthetic output"],
            "inputs": [self.anchor], "outputs": [self.anchor],
            "inventory_sha256": inventory, "live_effects_performed": False})

    def release(self):
        data = b"print('synthetic fixture')\n"
        for name in ("build-one", "build-two"):
            (self.base / name).mkdir()
            (self.base / name / "fixture.py").write_bytes(data)
        admission = self.write("evidence/admission.json", {
            "schema": "elite-library-admission/v1", "revision": "fixture-v1",
            "mode": "AUTHORED", "role": "GLUE", "files": ["fixture.py"], "source_lock": self.anchor,
            "gates": {f"G{i}": {"status": "PASS", "evidence": self.anchor} for i in range(9)}})
        inventory = self.write("evidence/inventory.json", {
            "schema": "elite-library-reference-inventory/v1", "revision": "fixture-v1",
            "files": [{"path": "fixture.py", "sha256": gate.digest(data)}],
            "provenance": [{"mode": "AUTHORED", "role": "GLUE", "files": ["fixture.py"],
                            "admission": admission, "glue_justification": "Test-only validator fixture."}]})
        self.record["reference"] = {"first_root": str(self.base / "build-one"),
                                    "second_root": str(self.base / "build-two"), "inventory": inventory}
        for name in gate.CONTROLS:
            self.record["controls"][name] = {"status": "PROVEN_LOCAL", "receipt": self.run_receipt(name, inventory["sha256"])}

    def check(self, contains, level="preparation"):
        errors = gate.validate_library_record(self.record, self.root, level)
        self.assertTrue(any(contains in message for message in errors), errors)

    def test_preparation_has_distinct_success_without_product_claim(self):
        self.assertEqual(gate.validate_library_record(self.record, self.root, "preparation"), [])

    def test_release_requires_remaining_work(self):
        self.check("reference", "release")
        self.check("T2802: PROVEN_LOCAL required", "release")

    def test_release_accepts_only_complete_exact_synthetic_evidence(self):
        self.release()
        self.assertEqual(gate.validate_library_record(self.record, self.root, "release"), [])

    def test_production_claim_is_rejected(self):
        self.record["production_authorized"] = True
        self.check("production_authorized")

    def test_project_schema_cannot_be_inherited(self):
        self.record["schema"] = "elite-project-readiness-gate/v2"
        self.check("schema/scope")

    def test_unknown_fields_are_rejected(self):
        self.record["status"] = "READY_TO_BUILD"
        self.check("exact fields")

    def test_missing_and_tampered_evidence_are_rejected(self):
        for ref in ({"path": "missing", "sha256": "0" * 64}, {"path": self.anchor["path"], "sha256": "0" * 64}):
            with self.subTest(ref=ref):
                self.record["assurance"] = ref
                self.check("assurance:")

    def test_path_traversal_and_ambiguous_paths_rejected(self):
        for name in ("../outside", "evidence/../method.txt", "C:/outside", "evidence\\method.txt", "evidence//method.txt", "evidence/method.txt#fake"):
            with self.subTest(path=name):
                self.record["tooling"] = {"path": name, "sha256": self.anchor["sha256"]}
                self.check("tooling:")

    def test_unresolved_local_bug_blocks_preparation(self):
        self.record["blockers"] = ["dependency SCA finding"]
        self.check("blocker: exact fields")

    def test_explicit_resolution_plan_allows_work_but_never_release(self):
        self.record["blockers"] = [{"id": "source-review", "kind": "SOURCE_PROVENANCE",
                                    "status": "OPEN", "affected_controls": ["T2801", "T2802"],
                                    "resolution_plan": self.anchor}]
        self.assertEqual(gate.validate_library_record(self.record, self.root, "preparation"), [])
        self.check("unresolved local blockers prevent release readiness", "release")

    def test_receipt_is_bound_to_control_and_revision(self):
        self.record["controls"]["T2801"]["receipt"] = self.run_receipt("T2802")
        self.check("exact-revision local run")

    def test_target_receipt_cannot_pass_as_local(self):
        path = self.root / "evidence/T2801.json"
        receipt = json.loads(path.read_bytes())
        receipt["live_effects_performed"] = True
        self.record["controls"]["T2801"]["receipt"] = self.write("evidence/T2801.json", receipt)
        self.check("exact-revision local run")

    def test_receipt_without_actual_outputs_rejected(self):
        path = self.root / "evidence/T2801.json"
        receipt = json.loads(path.read_bytes())
        receipt["outputs"] = []
        self.record["controls"]["T2801"]["receipt"] = self.write("evidence/T2801.json", receipt)
        self.check("outputs must be")

    def deferred(self, kind="USER_CREDENTIALS"):
        return {"id": "payment-access", "kind": kind, "control": "T2805",
                "reason": "No user account provided", "trigger": "User configures account later",
                "evidence": self.anchor, "code_complete": True, "local_receipt": self.run_receipt("T2805")}

    def test_credential_deferral_requires_complete_code_and_tests(self):
        self.record["deferred"] = [self.deferred()]
        self.assertEqual(gate.validate_library_record(self.record, self.root, "preparation"), [])
        self.record["deferred"][0]["code_complete"] = False
        self.check("cannot defer missing implementation")

    def test_generic_sca_and_code_deferral_rejected(self):
        for kind in ("SCA", "MISSING_CODE", "POLICY", "DAYBREAK_OTHER"):
            self.record["deferred"] = [self.deferred(kind)]
            self.check("only USER_CREDENTIALS")

    def test_daybreak_exclusion_is_exact(self):
        item = self.deferred("DAYBREAK_LIBXML2")
        item.update(id="daybreak-libxml2", control="T2803", code_complete=False, local_receipt=None)
        self.record["deferred"] = [item]
        self.assertEqual(gate.validate_library_record(self.record, self.root, "preparation"), [])
        item["control"] = "T2807"
        self.check("only exact excluded libxml2")

    def test_rebuild_detects_changed_missing_and_extra_files(self):
        self.release()
        target = self.base / "build-two" / "fixture.py"
        target.write_bytes(b"modified\n")
        self.check("tree differs", "release")
        target.unlink()
        self.check("tree differs", "release")
        target.write_bytes((self.base / "build-one" / "fixture.py").read_bytes())
        (target.parent / "extra.py").write_bytes(b"extra\n")
        self.check("tree differs", "release")

    def test_rebuild_rejects_same_destination_and_canonical_product(self):
        self.release()
        self.record["reference"]["second_root"] = self.record["reference"]["first_root"]
        self.check("independent non-overlapping", "release")
        self.record["reference"]["second_root"] = str(self.root)
        self.check("outside canonical root", "release")

    def mutate_inventory(self, mutate):
        value = json.loads((self.root / "evidence/inventory.json").read_bytes())
        mutate(value)
        ref = self.write("evidence/inventory.json", value)
        self.record["reference"]["inventory"] = ref
        for name in gate.CONTROLS:
            self.record["controls"][name]["receipt"] = self.run_receipt(name, ref["sha256"])

    def test_authored_cannot_masquerade_as_vendor_code(self):
        self.release()
        self.mutate_inventory(lambda value: value["provenance"][0].update(mode="GOOGLE"))
        self.check("declared provenance mode", "release")

    def test_authored_requires_glue_justification(self):
        self.release()
        self.mutate_inventory(lambda value: value["provenance"][0].update(glue_justification=""))
        self.check("unavoidable glue", "release")

    def test_authored_business_logic_rejected_even_with_glue_text(self):
        self.release()
        self.mutate_inventory(lambda value: value["provenance"][0].update(role="BUSINESS_LOGIC"))
        self.check("AUTHORED business logic", "release")

    def test_preparation_does_not_close_t2801(self):
        self.record["controls"]["T2801"]["status"] = "PROVEN_LOCAL"
        self.check("not task completion")

    def test_malformed_deferred_identity_fails_closed(self):
        item = self.deferred()
        item["id"] = []
        self.record["deferred"] = [item]
        self.check("unique non-empty id")

    def test_uncovered_source_file_blocks_release(self):
        self.release()
        self.mutate_inventory(lambda value: value.update(provenance=[]))
        self.check("coverage must equal", "release")

    def test_unpassed_g0_g8_blocks_release(self):
        self.release()
        value = json.loads((self.root / "evidence/admission.json").read_bytes())
        value["gates"]["G7"]["status"] = "CONDITIONED"
        ref = self.write("evidence/admission.json", value)
        self.mutate_inventory(lambda value: value["provenance"][0].update(admission=ref))
        self.check("G7: PASS required", "release")

    def test_old_receipts_cannot_prove_new_inventory(self):
        self.release()
        self.record["controls"]["T2802"]["receipt"] = self.run_receipt("T2802", "0" * 64)
        self.check("not bound to the final", "release")

    def test_cli_requires_explicit_scope_and_level(self):
        path = self.root / "PROJECT_READINESS_GATE.json"
        path.write_text(json.dumps(self.record), encoding="utf-8")
        command = [sys.executable, "-X", "utf8", "-B", str(Path(__file__).with_name("validate_project_readiness.py")), "--project-root", str(self.root)]
        default = subprocess.run(command, capture_output=True, timeout=10)
        self.assertEqual(default.returncode, 2)
        self.assertNotIn(b'"status": "READY_TO_BUILD"', default.stdout)
        absent_level = subprocess.run(command + ["--scope", "LIBRARY_INFRASTRUCTURE"], capture_output=True, timeout=10)
        self.assertEqual(absent_level.returncode, 2)
        result = subprocess.run(command + ["--scope", "LIBRARY_INFRASTRUCTURE", "--level", "preparation"], capture_output=True, timeout=10)
        self.assertEqual(result.returncode, 0, result.stderr)
        report = json.loads(result.stdout)
        self.assertEqual(report["status"], "READY_FOR_LIBRARY_WORK")
        self.assertFalse(report["production_authorized"])
        self.assertEqual(report["record_sha256"], gate.digest(path.read_bytes()))


if __name__ == "__main__":
    unittest.main()
````

### FILE: `project_readiness_gate/library-readiness.template.json`
```yaml
block_id: "PROJECT-START-READINESS-VALIDATOR:library-template:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite local validator glue for the explicit V402 user library scope; no external code attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "ccfc2c4ec1f51db802ec3ef65c0050682527fd4865a415aac38eb652edc40508"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-library-infrastructure-readiness/v1",
  "scope": "LIBRARY_INFRASTRUCTURE",
  "revision": "",
  "production_authorized": false,
  "request": {
    "path": "",
    "sha256": ""
  },
  "plan": {
    "path": "",
    "sha256": ""
  },
  "authority_decisions": {
    "path": "",
    "sha256": ""
  },
  "assurance": {
    "path": "",
    "sha256": ""
  },
  "tooling": {
    "path": "",
    "sha256": ""
  },
  "source_lock": {
    "path": "",
    "sha256": ""
  },
  "notices": {
    "path": "",
    "sha256": ""
  },
  "blockers": [
    {
      "id": "preparation-incomplete",
      "kind": "MISSING_EVIDENCE",
      "status": "OPEN",
      "affected_controls": [
        "T2801"
      ],
      "resolution_plan": {
        "path": "",
        "sha256": ""
      }
    }
  ],
  "deferred": [],
  "reference": null,
  "controls": {
    "T2801": {
      "status": "PENDING",
      "receipt": null
    },
    "T2802": {
      "status": "PENDING",
      "receipt": null
    },
    "T2803": {
      "status": "PENDING",
      "receipt": null
    },
    "T2804": {
      "status": "PENDING",
      "receipt": null
    },
    "T2805": {
      "status": "PENDING",
      "receipt": null
    },
    "T2806": {
      "status": "PENDING",
      "receipt": null
    },
    "T2807": {
      "status": "PENDING",
      "receipt": null
    },
    "T2808": {
      "status": "PENDING",
      "receipt": null
    },
    "T2809": {
      "status": "PENDING",
      "receipt": null
    },
    "T2810": {
      "status": "PENDING",
      "receipt": null
    },
    "ARCA_INFRA": {
      "status": "PENDING",
      "receipt": null
    }
  }
}
````

## 6. Configuration surface

| Variable | Type | Safe default | Validation | Secret | Mutability | Effect |
|---|---|---|---|---|---|---|
| `--project-root` | existing directory | none | resolved exact root | no | invocation | bounds all reads/report |
| `--record` | root filename | `PROJECT_READINESS_GATE.json` | regular direct child, no symlink | no | invocation | readiness authority |
| `--report` | relative path | stdout only | must remain under root and not exist | no | invocation | atomic evidence output |
| `--round` | enum | none | exactly `A` through `H` | no | invocation | selects one advisory topic |
| advisory `--output` | root filename | none | exact `PROJECT_ADVISORY_<ROUND>.md`, new regular file | no | invocation | immutable prompt receipt |
| template decisions | structured JSON | all blocked | semantic rules + local evidence + connected journeys | references only | conversation | opens or blocks build |

The record prohibits literal credentials/private keys. Secret access remains outside Markdown/JSON; only logical identity references and redacted probe evidence are allowed.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| CPython | 3.14.4 verified | JSON/path/CLI/tests | PSF-2.0 | runtime/test | python.org distribution installed on audit host |
| Python standard library | 3.14.4 | all behavior | PSF-2.0 | runtime | CPython |

No PyPI package, network service, Git repository or paid action is required.

## 8. Apply order

1. Materialize this pack alone into an empty project/tooling target.
2. Copy the template as `PROJECT_READINESS_GATE.json`; keep the original unchanged.
3. For each pending round, render `PROJECT_ADVISORY_<ROUND>.md`, explain it in manageable groups and record its path, answer and evidence; never put secrets in chat/files.
4. Create the remaining human/readiness artifacts, classify the 48 capabilities and map every required capability into release-bound connected journeys with all interface-to-recovery evidence links.
5. Run both test modules, then the readiness CLI without a report while blocked.
6. Only after prompts and evidence are complete, run with a new report path and require exit 0/status `READY_TO_BUILD`.
7. For an admitted outbox consumer, render the selected class/variant mapping from the passing record into a pristine consumer template; require the renderer PASS while all operational proofs remain false.
8. Compose and configure product packs only after those receipts; rerun readiness on every material delta.

Existing targets abort on collisions. Rollback deletes only the newly materialized tooling/report after preserving readiness/failure history.

## 9. Verification

- `python -m py_compile validate_project_readiness.py render_consumer_profile.py render_project_advisory.py test_validate_project_readiness.py test_render_project_advisory.py` → exit 0.
- `python -m unittest -v test_render_project_advisory.py test_validate_project_readiness.py` → 75/75 PASS.
- Distributed template → exit 2/BLOCKED with more than twenty reasons.
- Complete isolated fixture → exit 0/READY_TO_BUILD and atomic report.
- Negative coverage: advisory duplicate/unknown fields, missing round/source, assumed example, tamper, cross-round, traversal and overwrite; readiness missing artifact/blocker/round/capability/dependency/document mapping/access/provenance/secret/path/timestamp; connected journey missing/duplicate/unknown/uncovered capability/version drift/schema extension/evidence gap/first-slice mismatch; consumer renderer blocked/direct/duplicate/tampered/preconfigured/overwrite.
- Fresh materialization → eight file hashes identical; global library verifier composes the profile and compares the mapping field set against the consumer.

Project-specific proof remains mandatory; these tests prove the gate, not the truth of a future project's evidence.

## 10. Reconstruction evidence

Staging V170 on CPython 3.14.4: syntax PASS and 62 unit/CLI regressions PASS. Eight official-source-routed topics render deterministically; round F now collects the role→journey→help→training→support matrix and versioned operational content. Catalog drift, assumed examples, tamper, cross-round, traversal and overwrite fail closed. A valid readiness fixture requires all eight exact prompt receipts, every mandatory artifact, 48 capability rows, connected journey closure, evidence and exact direct/outbox mapping contracts. Consumer rendering remains closed and the distributed template remains blocked. Final evidence is issued only after fresh Markdown reconstruction, profile composition and the global gate.

## V328 — diagnóstico y publicación exclusiva verificados

0.6.2 conserva schema, catálogo, prompts y validate_record. Maneja explícitamente
inputs/OS/Unicode del CLI y publica el reporte completo mediante hard link
exclusivo después de flush/fsync; no reemplaza un destino tardío. El filesystem
debe soportar esta operación; de lo contrario falla cerrado, sin fallback que
sobrescriba. No demuestra durabilidad del directorio ante corte eléctrico ni
defensa ante reemplazo adversarial concurrente de directorios. 75tests PASS y8archivos reconstruidos;6archivos y validate_record idénticos al
baseline. Los negativos anteriores permanecen y el expediente real conserva
42bloqueos. Ver reconstruction_evidence/READINESS_REPORT_PUBLICATION_V328.md.


## V402 — explicit library infrastructure scope, 0.7.0

Adds three AUTHORED validator/glue files (11 total), no product/domain code.
The user explicitly authorized preparing library infrastructure without secrets.
PROJECT remains the default and retains its prior READY_TO_BUILD checks.
LIBRARY_INFRASTRUCTURE requires explicit preparation/release level. Preparation
emits READY_FOR_LIBRARY_WORK from hash-bound scope/assurance/plan/tool evidence;
T2801 is PREPARATION_PROVEN, not closed. Open local code/source/dependency gaps
need affected controls and an exact resolution plan; they always reject release.
Release requires all eleven controls, exact external independent materializations,
inventory-bound local receipts and declared file provenance/admission G0–G8.
AUTHORED BUSINESS_LOGIC is rejected; AUTHORED glue needs explicit justification.
Credential and exact Daybreak exclusions cannot hide missing implementation.
The only final status is READY_FOR_LIBRARY_USE with production_authorized false.

Verification: 102 tests PASS (75 existing + 27 new) on CPython 3.14.4; source and
reconstructed kit bytes are compared separately. This validates tooling, not the
truth of an arbitrary admission statement or the actual library release.
