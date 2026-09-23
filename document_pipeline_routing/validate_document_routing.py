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
    "PADDLEOCR_LOCAL": "markdown_system/PADDLEOCR_LOCAL_RUNTIME_PACK_PLAN.md",
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
