# Document Pipeline Routing Gate

## 1. Metadata

```yaml
pack_id: "DOCUMENT-PIPELINE-ROUTING-GATE"
pack_version: "0.3.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un gate stdlib determinista que exige las 21 clases base, admite clases adicionales canónicas con el mismo contrato completo, enlaza cada lane activa a un perfil oficial exacto y emite un receipt V3 sin llamar providers ni autorizar storage."
stacks: ["CPython 3.12+ stdlib", "JSON", "SHA-256"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x", "STRICT-DOCUMENT-FIELD-EVALUATION-GATE 0.1.x", "Azure/Google/AWS/MarkItDown document profiles"]
incompatible_with: ["provider elegido por intuición", "clase omitida", "perfil o path arbitrario", "corpus/acceso no probado", "auto-storage", "clasificación o extracción implícita"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/tree/1b5fd74454e593de233342a02ee911af8ee38359", "https://github.com/microsoft/content-processing-solution-accelerator/tree/b47cec48475cfedd7109debf2f407ef0c5530311", "https://github.com/GoogleCloudPlatform/document-ai-samples/tree/001ba391ab4a2f40d001cc0387618cb3c3699523"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use antes de materializar una lane documental para convertir la decisión del usuario en selección verificable. AWS publica clasificación/segmentación por clases y schemas versionados; Microsoft publica un pipeline Extract/Map/Evaluate/Save; Google publica processors especializados y custom classifiers/extractors. El gate local no copia sus motores ni inventa un clasificador universal: sólo valida que el proyecto haya nombrado clases, owners, lanes, profiles, seguridad, corpus, acceso y evaluación.

No usar como classifier, OCR, extractor, benchmark, aprobación de precisión ni autorización de persistencia. Las tres fuentes upstream permanecen bajo sus clasificaciones y condiciones del source lock; este glue es íntegramente `AUTHORED`.

## 3. Architecture contract

El template enumera 21 clases base y es deliberadamente inválido mientras conserve owners vacíos o `AWAITING_USER`. Una configuración admitida conserva y clasifica cada fila base exactamente una vez; también puede añadir cualquier ID canónico kebab-case para una clase nueva. Base y adicionales usan `REQUIRED`, `OPTIONAL`, `NONE_WITH_REASON` o `BLOCKED` y exigen al menos una requerida. Toda clase activa, incluida una adicional, fija schema/fields/MIME/límites, una única lane `PRIMARY`, perfiles relativos seguros, receipt de seguridad, corpus autorizado con ground truth, acceso sandbox `PROVEN`, evaluación estricta y storage false.

El mapa provider→pack plan está cerrado a los perfiles Azure, Google, AWS y MarkItDown de esta biblioteca. El validator no ejecuta esos planes. Escribe un receipt V3 atómico con SHA-256 de la configuración, inventario base, IDs adicionales, planes seleccionados y cada ruta clase→provider→role→plan→provider profile; cualquier error o output ocupado deja cero commit.

## 4. Exact file manifest

```text
CREATE document_pipeline_routing/document-routing.template.json
CREATE document_pipeline_routing/validate_document_routing.py
CREATE document_pipeline_routing/test_validate_document_routing.py
CREATE document_pipeline_routing/README.md
```

## 5. Materialization blocks

### FILE: `document_pipeline_routing/document-routing.template.json`
```yaml
block_id: "DOCUMENT-PIPELINE-ROUTING-GATE:template:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed operationalization of the project document decision and official provider class/schema boundaries"
license: "LicenseRef-Workspace-Owner"
sha256: "618ece2942cf2e2296289b7ad15a16ba7b244f2d5860236092159e35b2740e33"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-document-pipeline-routing/v1",
  "project": "",
  "owner": "",
  "security_owner": "",
  "data_owner": "",
  "review_operations_owner": "",
  "automatic_storage_authorized": false,
  "classes": [
    {"id": "supplier-invoice", "decision": "AWAITING_USER"},
    {"id": "receipt", "decision": "AWAITING_USER"},
    {"id": "proforma-invoice", "decision": "AWAITING_USER"},
    {"id": "commercial-invoice", "decision": "AWAITING_USER"},
    {"id": "packing-list", "decision": "AWAITING_USER"},
    {"id": "purchase-order", "decision": "AWAITING_USER"},
    {"id": "purchase-order-confirmation", "decision": "AWAITING_USER"},
    {"id": "bill-of-lading-airway-bill", "decision": "AWAITING_USER"},
    {"id": "delivery-note", "decision": "AWAITING_USER"},
    {"id": "certificate-of-origin", "decision": "AWAITING_USER"},
    {"id": "quality-inspection-certificate", "decision": "AWAITING_USER"},
    {"id": "customs-declaration-clearance", "decision": "AWAITING_USER"},
    {"id": "insurance-freight-document", "decision": "AWAITING_USER"},
    {"id": "price-list-quotation", "decision": "AWAITING_USER"},
    {"id": "product-specification-catalog", "decision": "AWAITING_USER"},
    {"id": "contract-addendum", "decision": "AWAITING_USER"},
    {"id": "bank-payment-statement", "decision": "AWAITING_USER"},
    {"id": "warranty-service-claim", "decision": "AWAITING_USER"},
    {"id": "identity-regulatory-document", "decision": "AWAITING_USER"},
    {"id": "multi-document-package", "decision": "AWAITING_USER"},
    {"id": "other-named-class", "decision": "AWAITING_USER"}
  ]
}
````

### FILE: `document_pipeline_routing/validate_document_routing.py`
```yaml
block_id: "DOCUMENT-PIPELINE-ROUTING-GATE:validator:v1"
operation: CREATE
provenance: AUTHORED
source: "deterministic local validator; no provider algorithm or vendor code copied"
license: "LicenseRef-Workspace-Owner"
sha256: "c2ceb1a63887dfd188cd4f8abf387c3439e3bb06451f9e23068715802406d2f5"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
from datetime import datetime, timezone
import hashlib
import json
from pathlib import Path
import re
import shutil
import uuid
from typing import Any


SCHEMA = "elite-document-pipeline-routing/v1"
RECEIPT_SCHEMA = "elite-document-pipeline-routing-receipt/v3"
SECURITY_RECEIPT_SCHEMA = "elite-secure-local-file-receipt/v1"
REQUIRED_CLASS_IDS = {
    "supplier-invoice", "receipt", "proforma-invoice", "commercial-invoice", "packing-list",
    "purchase-order", "purchase-order-confirmation", "bill-of-lading-airway-bill", "delivery-note",
    "certificate-of-origin", "quality-inspection-certificate", "customs-declaration-clearance",
    "insurance-freight-document", "price-list-quotation", "product-specification-catalog",
    "contract-addendum", "bank-payment-statement", "warranty-service-claim",
    "identity-regulatory-document", "multi-document-package", "other-named-class",
}
PROVIDER_PLANS = {
    "AZURE": "markdown_system/AZURE_DOCUMENT_RUNTIME_PACK_PLAN.md",
    "GOOGLE": "markdown_system/GOOGLE_DOCUMENT_RUNTIME_PACK_PLAN.md",
    "AWS": "markdown_system/AWS_TEXTRACT_DOCUMENT_RUNTIME_PACK_PLAN.md",
    "MARKITDOWN_LOCAL": "markdown_system/MARKITDOWN_LOCAL_RUNTIME_PACK_PLAN.md",
}
MIME_TYPES = {
    "application/pdf", "image/png", "image/jpeg", "image/tiff", "text/plain", "text/csv",
    "text/html", "text/markdown", "application/json", "application/xml",
    "application/vnd.ms-excel", "application/vnd.ms-outlook",
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    "application/vnd.openxmlformats-officedocument.presentationml.presentation",
}
ID_PATTERN = re.compile(r"^[a-z0-9]+(?:-[a-z0-9]+)*$")
VERSION_PATTERN = re.compile(r"^[A-Za-z0-9][A-Za-z0-9._/-]*$")


def sha256_bytes(payload: bytes) -> str:
    return hashlib.sha256(payload).hexdigest()


def nonempty(value: object, label: str) -> str:
    if not isinstance(value, str) or not value.strip():
        raise ValueError(f"{label} is required")
    return value.strip()


def safe_relative(value: object, label: str) -> str:
    text = nonempty(value, label).replace("\\", "/")
    path = Path(text)
    if path.is_absolute() or ".." in path.parts or "." in path.parts or any(not part for part in path.parts):
        raise ValueError(f"{label} must be a canonical relative path")
    return text


def unique_strings(value: object, label: str, allowed: set[str] | None = None) -> list[str]:
    if not isinstance(value, list) or not value:
        raise ValueError(f"{label} must be a non-empty array")
    result: list[str] = []
    for item in value:
        text = nonempty(item, label)
        if allowed is not None and text not in allowed:
            raise ValueError(f"{label} contains an unsupported value")
        if text in result:
            raise ValueError(f"{label} contains duplicates")
        result.append(text)
    return result


def positive_int(value: object, label: str) -> int:
    if not isinstance(value, int) or isinstance(value, bool) or value <= 0:
        raise ValueError(f"{label} must be a positive integer")
    return value


def validate_lane(lane: object, label: str) -> dict[str, str]:
    if not isinstance(lane, dict):
        raise ValueError(f"{label} must be an object")
    provider = nonempty(lane.get("provider"), f"{label}.provider")
    if provider not in PROVIDER_PLANS:
        raise ValueError(f"{label}.provider is unsupported")
    role = lane.get("role")
    if role not in {"PRIMARY", "EVALUATION"}:
        raise ValueError(f"{label}.role must be PRIMARY or EVALUATION")
    plan = safe_relative(lane.get("pack_plan"), f"{label}.pack_plan")
    if plan != PROVIDER_PLANS[provider]:
        raise ValueError(f"{label}.pack_plan does not match provider")
    profile_path = safe_relative(lane.get("provider_profile_path"), f"{label}.provider_profile_path")
    return {"provider": provider, "role": role, "pack_plan": plan, "provider_profile_path": profile_path}


def validate_active_class(item: dict[str, Any], class_id: str) -> dict[str, Any]:
    schema_version = nonempty(item.get("schema_version"), f"{class_id}.schema_version")
    if VERSION_PATTERN.fullmatch(schema_version) is None:
        raise ValueError(f"{class_id}.schema_version is invalid")
    mime_types = unique_strings(item.get("mime_types"), f"{class_id}.mime_types", MIME_TYPES)
    max_bytes = positive_int(item.get("max_bytes"), f"{class_id}.max_bytes")
    max_pages = positive_int(item.get("max_pages"), f"{class_id}.max_pages")
    fields = unique_strings(item.get("required_fields"), f"{class_id}.required_fields")
    lanes_value = item.get("candidate_lanes")
    if not isinstance(lanes_value, list) or not lanes_value:
        raise ValueError(f"{class_id}.candidate_lanes must be a non-empty array")
    lanes = [validate_lane(lane, f"{class_id}.candidate_lanes[{index}]") for index, lane in enumerate(lanes_value)]
    providers = [lane["provider"] for lane in lanes]
    if len(set(providers)) != len(providers):
        raise ValueError(f"{class_id}.candidate_lanes contains duplicate providers")
    if sum(lane["role"] == "PRIMARY" for lane in lanes) != 1:
        raise ValueError(f"{class_id}.candidate_lanes requires exactly one PRIMARY")
    security = item.get("security")
    if not isinstance(security, dict):
        raise ValueError(f"{class_id}.security is required")
    security_policy_path = safe_relative(security.get("policy_path"), f"{class_id}.security.policy_path")
    if security.get("required_receipt_schema") != SECURITY_RECEIPT_SCHEMA:
        raise ValueError(f"{class_id}.security.required_receipt_schema mismatch")
    evaluation = item.get("evaluation")
    if not isinstance(evaluation, dict):
        raise ValueError(f"{class_id}.evaluation is required")
    strict_profile_path = safe_relative(evaluation.get("strict_profile_path"), f"{class_id}.evaluation.strict_profile_path")
    if evaluation.get("status") not in {"EVALUATION_REQUIRED", "REVIEW_ONLY"}:
        raise ValueError(f"{class_id}.evaluation.status cannot authorize storage")
    corpus = item.get("corpus")
    if not isinstance(corpus, dict):
        raise ValueError(f"{class_id}.corpus is required")
    if corpus.get("authorized") is not True:
        raise ValueError(f"{class_id}.corpus must be explicitly authorized")
    ground_truth_documents = positive_int(corpus.get("ground_truth_documents"), f"{class_id}.corpus.ground_truth_documents")
    ground_truth_owner = nonempty(corpus.get("ground_truth_owner"), f"{class_id}.corpus.ground_truth_owner")
    corpus_evidence_path = safe_relative(corpus.get("evidence_path"), f"{class_id}.corpus.evidence_path")
    access = item.get("access")
    if not isinstance(access, dict):
        raise ValueError(f"{class_id}.access is required")
    if access.get("probe_status") != "PROVEN":
        raise ValueError(f"{class_id}.access.probe_status must be PROVEN")
    identity_reference = nonempty(access.get("identity_reference"), f"{class_id}.access.identity_reference")
    region = nonempty(access.get("region"), f"{class_id}.access.region")
    access_evidence_path = safe_relative(access.get("evidence_path"), f"{class_id}.access.evidence_path")
    if item.get("automatic_storage") is not False:
        raise ValueError(f"{class_id}.automatic_storage must be false")
    return {
        "id": class_id,
        "decision": item["decision"],
        "schema_version": schema_version,
        "mime_types": mime_types,
        "max_bytes": max_bytes,
        "max_pages": max_pages,
        "required_fields": fields,
        "candidate_lanes": lanes,
        "security_policy_path": security_policy_path,
        "strict_profile_path": strict_profile_path,
        "ground_truth_documents": ground_truth_documents,
        "ground_truth_owner": ground_truth_owner,
        "corpus_evidence_path": corpus_evidence_path,
        "identity_reference": identity_reference,
        "region": region,
        "access_evidence_path": access_evidence_path,
        "automatic_storage": False,
    }


def validate_configuration(path: Path) -> tuple[dict[str, Any], bytes]:
    resolved = path.resolve(strict=True)
    if path.is_symlink() or not resolved.is_file():
        raise ValueError("routing configuration must be a regular non-symlink file")
    payload = resolved.read_bytes()
    try:
        value = json.loads(payload)
    except (UnicodeDecodeError, json.JSONDecodeError) as error:
        raise ValueError("routing configuration is not valid UTF-8 JSON") from error
    if not isinstance(value, dict) or value.get("schema") != SCHEMA:
        raise ValueError("routing configuration schema mismatch")
    owners = {name: nonempty(value.get(name), name) for name in ("project", "owner", "security_owner", "data_owner", "review_operations_owner")}
    if value.get("automatic_storage_authorized") is not False:
        raise ValueError("routing configuration must deny automatic storage")
    classes = value.get("classes")
    if not isinstance(classes, list):
        raise ValueError("classes must be an array")
    observed: dict[str, dict[str, Any]] = {}
    for item in classes:
        if not isinstance(item, dict):
            raise ValueError("each class must be an object")
        class_id = nonempty(item.get("id"), "class.id")
        if ID_PATTERN.fullmatch(class_id) is None:
            raise ValueError(f"non-canonical class id: {class_id}")
        if class_id in observed:
            raise ValueError(f"duplicate class id: {class_id}")
        decision = item.get("decision")
        if decision in {"REQUIRED", "OPTIONAL"}:
            observed[class_id] = validate_active_class(item, class_id)
        elif decision in {"NONE_WITH_REASON", "BLOCKED"}:
            observed[class_id] = {"id": class_id, "decision": decision, "reason": nonempty(item.get("reason"), f"{class_id}.reason")}
        else:
            raise ValueError(f"{class_id}.decision is not complete")
    missing = sorted(REQUIRED_CLASS_IDS - set(observed))
    if missing:
        raise ValueError("required class inventory is incomplete: " + ",".join(missing))
    if not any(item["decision"] == "REQUIRED" for item in observed.values()):
        raise ValueError("at least one document class must be REQUIRED")
    return {**owners, "classes": [observed[class_id] for class_id in sorted(observed)]}, payload


def write_receipt(config_path: Path, output_directory: Path) -> dict[str, Any]:
    normalized, payload = validate_configuration(config_path)
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".document-routing-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        active_classes = [item for item in normalized["classes"] if item["decision"] in {"REQUIRED", "OPTIONAL"}]
        selected_plans = sorted({lane["pack_plan"] for item in active_classes for lane in item["candidate_lanes"]})
        selected_routes = [
            {
                "class_id": item["id"],
                "class_decision": item["decision"],
                "provider": lane["provider"],
                "role": lane["role"],
                "pack_plan": lane["pack_plan"],
                "provider_profile_path": lane["provider_profile_path"],
            }
            for item in active_classes
            for lane in item["candidate_lanes"]
        ]
        selected_routes.sort(key=lambda item: (item["class_id"], item["role"], item["provider"]))
        receipt = {
            "schema": RECEIPT_SCHEMA,
            "created_at": datetime.now(timezone.utc).isoformat(),
            "configuration_sha256": sha256_bytes(payload),
            "project": normalized["project"],
            "baseline_class_ids": sorted(REQUIRED_CLASS_IDS),
            "custom_class_ids": sorted(item["id"] for item in normalized["classes"] if item["id"] not in REQUIRED_CLASS_IDS),
            "required_classes": sorted(item["id"] for item in normalized["classes"] if item["decision"] == "REQUIRED"),
            "optional_classes": sorted(item["id"] for item in normalized["classes"] if item["decision"] == "OPTIONAL"),
            "selected_pack_plans": selected_plans,
            "selected_routes": selected_routes,
            "automatic_storage_authorized": False,
            "next_gate": "MATERIALIZE_SELECTED_PLANS_THEN_RUN_REAL_SECURITY_PROVIDER_AND_STRICT_EVALUATION",
        }
        (stage / "document-routing-receipt.json").write_text(json.dumps(receipt, indent=2, sort_keys=True) + "\n", encoding="utf-8", newline="\n")
        stage.rename(output)
        return receipt
    except Exception:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate an explicit document-class routing decision without calling providers")
    parser.add_argument("--configuration", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()
    receipt = write_receipt(args.configuration, args.output)
    print(json.dumps(receipt, sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `document_pipeline_routing/test_validate_document_routing.py`
```yaml
block_id: "DOCUMENT-PIPELINE-ROUTING-GATE:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "stdlib positive and fail-closed routing regressions"
license: "LicenseRef-Workspace-Owner"
sha256: "16f67a71da2a6cadd517b60120e48691548fff915415eb51c4e1b1fca21f95ee"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import hashlib
import json
from pathlib import Path
import tempfile
import unittest

from validate_document_routing import PROVIDER_PLANS, REQUIRED_CLASS_IDS, validate_configuration, write_receipt


def complete_configuration() -> dict:
    classes = []
    for class_id in sorted(REQUIRED_CLASS_IDS):
        if class_id != "supplier-invoice":
            classes.append({"id": class_id, "decision": "NONE_WITH_REASON", "reason": "not required by test project"})
            continue
        classes.append({
            "id": class_id,
            "decision": "REQUIRED",
            "schema_version": "supplier-invoice/v1",
            "mime_types": ["application/pdf"],
            "max_bytes": 10485760,
            "max_pages": 1,
            "required_fields": ["invoice_number", "total_amount"],
            "candidate_lanes": [{
                "provider": "AWS",
                "role": "PRIMARY",
                "pack_plan": PROVIDER_PLANS["AWS"],
                "provider_profile_path": "config/aws-document-profile.json",
            }],
            "security": {
                "policy_path": "config/secure-file-policy.json",
                "required_receipt_schema": "elite-secure-local-file-receipt/v1",
            },
            "evaluation": {
                "strict_profile_path": "config/strict-field-evaluation.json",
                "status": "EVALUATION_REQUIRED",
            },
            "corpus": {
                "authorized": True,
                "ground_truth_documents": 10,
                "ground_truth_owner": "data-owner",
                "evidence_path": "evidence/corpus/supplier-invoice.json",
            },
            "access": {
                "probe_status": "PROVEN",
                "identity_reference": "secret-manager://aws-role-reference",
                "region": "us-east-1",
                "evidence_path": "evidence/access/aws-sandbox.json",
            },
            "automatic_storage": False,
        })
    return {
        "schema": "elite-document-pipeline-routing/v1",
        "project": "test-project",
        "owner": "business-owner",
        "security_owner": "security-owner",
        "data_owner": "data-owner",
        "review_operations_owner": "review-owner",
        "automatic_storage_authorized": False,
        "classes": classes,
    }


class DocumentRoutingTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary = tempfile.TemporaryDirectory()
        self.root = Path(self.temporary.name)
        self.configuration = self.root / "routing.json"

    def tearDown(self) -> None:
        self.temporary.cleanup()

    def write(self, value: dict) -> None:
        self.configuration.write_text(json.dumps(value, sort_keys=True) + "\n", encoding="utf-8", newline="\n")

    def test_complete_explicit_inventory_writes_hash_linked_receipt(self) -> None:
        self.write(complete_configuration())
        output = self.root / "evidence"
        receipt = write_receipt(self.configuration, output)
        persisted = json.loads((output / "document-routing-receipt.json").read_text(encoding="utf-8"))
        self.assertEqual(receipt, persisted)
        self.assertEqual(receipt["configuration_sha256"], hashlib.sha256(self.configuration.read_bytes()).hexdigest())
        self.assertEqual(receipt["required_classes"], ["supplier-invoice"])
        self.assertEqual(receipt["selected_pack_plans"], [PROVIDER_PLANS["AWS"]])
        self.assertEqual(receipt["schema"], "elite-document-pipeline-routing-receipt/v3")
        self.assertEqual(receipt["baseline_class_ids"], sorted(REQUIRED_CLASS_IDS))
        self.assertEqual(receipt["custom_class_ids"], [])
        self.assertEqual(receipt["selected_routes"], [{
            "class_id": "supplier-invoice",
            "class_decision": "REQUIRED",
            "provider": "AWS",
            "role": "PRIMARY",
            "pack_plan": PROVIDER_PLANS["AWS"],
            "provider_profile_path": "config/aws-document-profile.json",
        }])
        self.assertFalse(receipt["automatic_storage_authorized"])

    def test_receipt_preserves_every_class_lane_and_profile_path(self) -> None:
        value = complete_configuration()
        active = next(item for item in value["classes"] if item["decision"] == "REQUIRED")
        active["candidate_lanes"].append({
            "provider": "GOOGLE",
            "role": "EVALUATION",
            "pack_plan": PROVIDER_PLANS["GOOGLE"],
            "provider_profile_path": "config/google-document-profile.json",
        })
        self.write(value)
        receipt = write_receipt(self.configuration, self.root / "route-evidence")
        self.assertEqual(receipt["selected_pack_plans"], sorted([PROVIDER_PLANS["AWS"], PROVIDER_PLANS["GOOGLE"]]))
        self.assertEqual(
            [(route["provider"], route["role"], route["provider_profile_path"]) for route in receipt["selected_routes"]],
            [("GOOGLE", "EVALUATION", "config/google-document-profile.json"), ("AWS", "PRIMARY", "config/aws-document-profile.json")],
        )

    def test_additional_canonical_class_is_routed_with_full_contract(self) -> None:
        value = complete_configuration()
        source = next(item for item in value["classes"] if item["decision"] == "REQUIRED")
        custom = json.loads(json.dumps(source))
        custom["id"] = "letter-of-credit-amendment"
        custom["schema_version"] = "letter-of-credit-amendment/v1"
        custom["required_fields"] = ["credit_number", "amendment_number", "effective_date"]
        custom["candidate_lanes"] = [{
            "provider": "GOOGLE",
            "role": "PRIMARY",
            "pack_plan": PROVIDER_PLANS["GOOGLE"],
            "provider_profile_path": "config/google-letter-of-credit-profile.json",
        }]
        custom["corpus"]["evidence_path"] = "evidence/corpus/letter-of-credit-amendment.json"
        custom["access"]["identity_reference"] = "secret-manager://google-document-ai-reference"
        custom["access"]["region"] = "us"
        custom["access"]["evidence_path"] = "evidence/access/google-sandbox.json"
        value["classes"].append(custom)
        self.write(value)
        receipt = write_receipt(self.configuration, self.root / "custom-route-evidence")
        self.assertEqual(receipt["custom_class_ids"], ["letter-of-credit-amendment"])
        self.assertEqual(receipt["required_classes"], ["letter-of-credit-amendment", "supplier-invoice"])
        self.assertIn(PROVIDER_PLANS["GOOGLE"], receipt["selected_pack_plans"])
        self.assertEqual(
            [route["class_id"] for route in receipt["selected_routes"]],
            ["letter-of-credit-amendment", "supplier-invoice"],
        )

    def test_additional_class_id_must_be_canonical(self) -> None:
        value = complete_configuration()
        value["classes"].append({"id": "../../new class", "decision": "NONE_WITH_REASON", "reason": "invalid identifier"})
        self.write(value)
        with self.assertRaisesRegex(ValueError, "non-canonical class id"):
            validate_configuration(self.configuration)

    def test_distributed_template_is_intentionally_not_routable(self) -> None:
        template = Path(__file__).with_name("document-routing.template.json")
        with self.assertRaisesRegex(ValueError, "project is required"):
            validate_configuration(template)

    def test_missing_or_duplicate_class_fails(self) -> None:
        value = complete_configuration()
        value["classes"].pop()
        self.write(value)
        with self.assertRaisesRegex(ValueError, "inventory is incomplete"):
            validate_configuration(self.configuration)
        value = complete_configuration()
        value["classes"].append(value["classes"][0].copy())
        self.write(value)
        with self.assertRaisesRegex(ValueError, "duplicate class id"):
            validate_configuration(self.configuration)

    def test_unclassified_and_missing_reason_fail(self) -> None:
        value = complete_configuration()
        value["classes"][0] = {"id": value["classes"][0]["id"], "decision": "AWAITING_USER"}
        self.write(value)
        with self.assertRaisesRegex(ValueError, "decision is not complete"):
            validate_configuration(self.configuration)
        value = complete_configuration()
        inactive = next(item for item in value["classes"] if item["decision"] == "NONE_WITH_REASON")
        inactive["reason"] = ""
        self.write(value)
        with self.assertRaisesRegex(ValueError, "reason is required"):
            validate_configuration(self.configuration)

    def test_provider_plan_mapping_and_safe_paths_are_closed(self) -> None:
        value = complete_configuration()
        active = next(item for item in value["classes"] if item["decision"] == "REQUIRED")
        active["candidate_lanes"][0]["pack_plan"] = PROVIDER_PLANS["GOOGLE"]
        self.write(value)
        with self.assertRaisesRegex(ValueError, "does not match provider"):
            validate_configuration(self.configuration)
        active["candidate_lanes"][0]["pack_plan"] = PROVIDER_PLANS["AWS"]
        active["candidate_lanes"][0]["provider_profile_path"] = "../secret.json"
        self.write(value)
        with self.assertRaisesRegex(ValueError, "canonical relative path"):
            validate_configuration(self.configuration)

    def test_required_class_needs_security_corpus_access_and_evaluation(self) -> None:
        cases = (
            ("security", None, "security is required"),
            ("corpus", {"authorized": False}, "corpus must be explicitly authorized"),
            ("access", {"probe_status": "NOT_PROVIDED"}, "probe_status must be PROVEN"),
            ("evaluation", {"strict_profile_path": "config/eval.json", "status": "READY_FOR_AUTOMATIC_STORAGE"}, "cannot authorize storage"),
        )
        for field, replacement, message in cases:
            value = complete_configuration()
            active = next(item for item in value["classes"] if item["decision"] == "REQUIRED")
            active[field] = replacement
            self.write(value)
            with self.assertRaisesRegex(ValueError, message):
                validate_configuration(self.configuration)

    def test_storage_authorization_is_rejected_at_both_levels(self) -> None:
        value = complete_configuration()
        value["automatic_storage_authorized"] = True
        self.write(value)
        with self.assertRaisesRegex(ValueError, "deny automatic storage"):
            validate_configuration(self.configuration)
        value = complete_configuration()
        active = next(item for item in value["classes"] if item["decision"] == "REQUIRED")
        active["automatic_storage"] = True
        self.write(value)
        with self.assertRaisesRegex(ValueError, "automatic_storage must be false"):
            validate_configuration(self.configuration)

    def test_existing_output_is_not_overwritten(self) -> None:
        self.write(complete_configuration())
        output = self.root / "evidence"
        output.mkdir()
        with self.assertRaises(FileExistsError):
            write_receipt(self.configuration, output)


if __name__ == "__main__":
    unittest.main()
````

### FILE: `document_pipeline_routing/README.md`
```yaml
block_id: "DOCUMENT-PIPELINE-ROUTING-GATE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local usage and non-attribution boundary"
license: "LicenseRef-Workspace-Owner"
sha256: "6e03916b4d1c7a92d7ef138a36d107de82926bcf32dc3fcb92e59b84fc4f1edf"
variables: []
secrets_allowed: false
```
````markdown
# Document pipeline routing gate

This directory is deterministic local `AUTHORED` control code. It is not AWS, Microsoft or Google code and it performs no classification, extraction, provider call or storage.

Its purpose is to stop an agent from guessing a document class, provider or implementation plan. The 21 baseline rows must always remain explicit so common enterprise document families cannot be omitted. A project may append any additional canonical kebab-case class ID; every additional class must pass the same schema, field, MIME, size, security, corpus, access, evaluation, provider-profile and no-storage contract as a baseline class. The distributed template is intentionally invalid until the user completes it.

The exact provider-to-plan map selects only library profiles that already wrap official public SDKs and source locks:

- Azure → `AZURE_DOCUMENT_RUNTIME_PACK_PLAN.md`;
- Google → `GOOGLE_DOCUMENT_RUNTIME_PACK_PLAN.md`;
- AWS → `AWS_TEXTRACT_DOCUMENT_RUNTIME_PACK_PLAN.md`;
- Microsoft MarkItDown local → `MARKITDOWN_LOCAL_RUNTIME_PACK_PLAN.md`.

Run:

```text
python -m unittest -v test_validate_document_routing.py
python validate_document_routing.py --configuration <completed-routing.json> --output <new-evidence-directory>
```

The V3 receipt preserves the baseline inventory, explicit custom-class IDs and the complete class → provider → role → exact pack plan → configured provider-profile path map, so the agent never has to reconstruct or guess a routing decision after validation. It only proves that the configuration is complete and hash-linked. It never proves provider accuracy or authorizes persistence. Every selected profile still has to run real security, provider, strict evaluation, corpus, review, load and recovery gates.
````

## 6. Configuration surface

| Campo | Tipo | Regla | Secreto | Efecto |
|---|---|---|---|---|
| owners | strings | cinco identidades no vacías | referencias, no credenciales | responsabilidad explícita |
| classes | 21 objetos base + 0..N adicionales | inventario base exacto/único; IDs adicionales kebab-case únicos | no | ninguna clase base se omite y una clase nueva no requiere cambiar código |
| decision | enum | sin `AWAITING_USER` | no | requerida/opcional/no aplica/bloqueada |
| candidate_lanes | array | provider único, exactamente un PRIMARY, plan exacto | no | selección materializable |
| provider_profile_path | relative path | canónico, sin traversal | no | perfil posterior |
| corpus/access | objects | autorizado + ground truth >0; probe PROVEN | referencias | bloquea suposición |
| security/evaluation | relative paths/status | receipt de seguridad V1 + evaluación no-storage | no | encadena gates |
| automatic storage | boolean | siempre false | no | jamás se promueve aquí |

## 7. Dependency bill

| Dependencia | Pin | Uso | Licencia | Gate |
|---|---|---|---|---|
| CPython stdlib | 3.12+ | JSON, paths, SHA, atomic staging, tests | PSF | compile + unittest |
| perfiles documentales Elite | versiones fijadas por sus planes | destino de selección, no runtime de este pack | por pack/upstream | mapping exacto |
| AWS/Microsoft/Google official source | commits del source lock | autoridad de patrones/clases/pipeline | MIT-0/MIT/Apache-2.0 | no se copia ni ejecuta aquí |

## 8. Apply order

1. Materializar en target vacío.
2. Copiar el template a una configuración de proyecto y preguntar al usuario cada owner/clase/variante; no editar el template canónico.
3. Completar las 21 clases base y añadir cada clase/variante adicional nombrada por el proyecto; para todas completar corpus, acceso, seguridad, evaluación y lanes.
4. Ejecutar tests y validator; conservar receipt.
5. Materializar sólo los `selected_pack_plans` del receipt y ejecutar sus approvals/gates reales.
6. Mantener storage false hasta evidencia de provider/corpus/campos/review/recovery separada.

## 9. Verification

```text
python -m py_compile validate_document_routing.py
python -m unittest -v test_validate_document_routing.py
python validate_document_routing.py --configuration <completed.json> --output <new-evidence-directory>
```

Esperado: once tests PASS. Los positivos clasifican 21/21, seleccionan AWS por mapping exacto, preservan cada lane AWS/Google y añaden `letter-of-credit-amendment` con schema/corpus/access completos sin cambiar el validator. Los negativos rechazan ID adicional no canónico, template incompleto, clase base faltante/duplicada, `AWAITING_USER`, razón vacía, provider/plan/path inválido, security/corpus/access/evaluation incompletos, storage true y overwrite. Receipt V3 enumera por separado baseline/custom IDs.

## 10. Reconstruction evidence

Evidencia gobernante: `reconstruction_evidence/DOCUMENT_PIPELINE_ROUTING_GATE_2026-08-28_V3.md`; V1/V2 permanecen snapshots históricos. El gate prueba configuración explícita, no clasificación, extracción, precisión ni producción. Las condiciones upstream están en `LIBRARY_FAILURE_LEARNING_LEDGER.md` y se trasladan al proyecto elegido.
