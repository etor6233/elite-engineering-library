from __future__ import annotations

import argparse
from datetime import datetime, timezone
import hashlib
import importlib.metadata
import json
import os
from pathlib import Path
import re
import shutil
import sys
from typing import Any, Callable, Iterable
import uuid

from facebook_business.adobjects.leadgenform import LeadgenForm
from facebook_business.api import FacebookAdsApi


SDK_VERSION = "26.0.1"
API_VERSION = "v26.0"
SOURCE_COMMIT = "788f363d15b1269ab5efb7cd00fb5e3b133cd99b"
ALLOWED_PROVIDER_FIELDS = {
    "id", "created_time", "form_id", "campaign_id", "campaign_name",
    "adset_id", "adset_name", "ad_id", "ad_name", "field_data",
    "is_organic", "platform", "custom_disclaimer_responses", "home_listing",
}


def sha256(value: bytes) -> str:
    return hashlib.sha256(value).hexdigest()


def _read_json(path: Path) -> dict[str, Any]:
    value = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(value, dict):
        raise ValueError(f"{path.name} must contain one JSON object")
    return value


def validate_configuration(profile: dict[str, Any], retrieval: dict[str, Any]) -> tuple[str, str, str, list[str], set[str], int]:
    proof_fields = (
        "app_registration_proven", "business_verification_proven",
        "lead_access_permission_proven", "page_and_form_ownership_proven",
        "test_lead_contract_proven", "pii_storage_controls_proven",
        "quota_and_cost_approved", "reconciliation_approved",
    )
    if (
        profile.get("provider") != "Meta Marketing API"
        or profile.get("sdk") != f"facebook-business=={SDK_VERSION}"
        or profile.get("graph_api_version") != API_VERSION
    ):
        raise ValueError("profile provider/SDK/API identity mismatch")
    if profile.get("decision") != "PROVEN" or any(profile.get(key) is not True for key in proof_fields):
        raise PermissionError("Meta lead access/policy profile remains blocked")
    if profile.get("automatic_business_write") is not False or profile.get("automatic_contact_eligibility") is not False:
        raise ValueError("automatic writes and contact eligibility must remain false")
    if not str(profile.get("data_retention", "")).strip():
        raise ValueError("data_retention is required")
    tenant_id = str(profile.get("tenant_id", "")).strip()
    organization_id = str(profile.get("organization_id", "")).strip()
    if not tenant_id or not organization_id:
        raise ValueError("tenant_id and organization_id are required")
    form_id = str(profile.get("form_id", ""))
    if not re.fullmatch(r"[0-9]{5,30}", form_id):
        raise ValueError("form_id must contain digits only")
    fields = retrieval.get("fields")
    approved_fields = profile.get("approved_provider_fields")
    if not isinstance(fields, list) or not fields or len(fields) > 20 or len(fields) != len(set(fields)):
        raise ValueError("fields must be a unique non-empty list with at most 20 entries")
    if not {"id", "created_time", "field_data"}.issubset(set(fields)):
        raise ValueError("id, created_time and field_data are mandatory")
    if not isinstance(approved_fields, list) or not set(fields).issubset(set(approved_fields)) or not set(fields).issubset(ALLOWED_PROVIDER_FIELDS):
        raise PermissionError("all provider fields must be locally allowed and explicitly approved")
    approved_form_fields = profile.get("approved_form_field_names")
    if not isinstance(approved_form_fields, list) or not approved_form_fields:
        raise ValueError("approved_form_field_names must be a non-empty list")
    approved_names: set[str] = set()
    for value in approved_form_fields:
        if not isinstance(value, str) or not re.fullmatch(r"[A-Za-z0-9_.-]{1,128}", value) or value in approved_names:
            raise ValueError("approved_form_field_names contains an invalid or duplicate name")
        approved_names.add(value)
    mode = retrieval.get("retrieval_mode")
    if mode not in {"TEST", "LIVE"}:
        raise ValueError("retrieval_mode must be TEST or LIVE")
    max_rows = retrieval.get("max_rows")
    if not isinstance(max_rows, int) or isinstance(max_rows, bool) or max_rows < 1 or max_rows > 1000:
        raise ValueError("max_rows must be within 1..1000")
    return tenant_id, organization_id, form_id, fields, approved_names, max_rows


def _row_to_dict(row: Any) -> dict[str, Any]:
    if isinstance(row, dict):
        return row
    exporter = getattr(row, "export_all_data", None)
    if not callable(exporter):
        raise RuntimeError("Meta lead row is not an official SDK object")
    value = exporter()
    if not isinstance(value, dict):
        raise RuntimeError("Meta SDK lead export must return an object")
    return value


def _normalize_row(row: dict[str, Any], tenant_id: str, organization_id: str, approved_names: set[str], is_test: bool) -> dict[str, Any]:
    candidate: dict[str, Any] = {
        "tenant_id": tenant_id,
        "organization_id": organization_id,
        "provider": "meta_lead_ads",
        "provider_lead_id": str(row.get("id", "")),
        "form_id": str(row.get("form_id", "")),
        "campaign_id": str(row.get("campaign_id", "")),
        "ad_group_id": str(row.get("adset_id", "")),
        "creative_id": str(row.get("ad_id", "")),
        "source_kind": str(row.get("platform", "")),
        "submitted_at": str(row.get("created_time", "")),
        "is_test": is_test,
        "fields": [],
        "contact_eligibility": "pending_policy",
    }
    errors: list[str] = []
    if not candidate["provider_lead_id"]:
        errors.append("MISSING_LEAD_ID")
    if not candidate["submitted_at"]:
        errors.append("MISSING_CREATED_TIME")
    field_data = row.get("field_data")
    if not isinstance(field_data, list):
        errors.append("INVALID_FIELD_DATA")
        field_data = []
    seen: set[str] = set()
    for item in field_data:
        if not isinstance(item, dict) or not isinstance(item.get("name"), str) or not isinstance(item.get("values"), list):
            errors.append("INVALID_FIELD_ENTRY")
            continue
        name = item["name"]
        values = item["values"]
        if name not in approved_names:
            errors.append("UNAPPROVED_FORM_FIELD")
            continue
        if name in seen:
            errors.append("DUPLICATE_FORM_FIELD")
            continue
        if len(values) != 1 or not isinstance(values[0], str) or not values[0].strip():
            errors.append("INVALID_FIELD_VALUE")
            continue
        seen.add(name)
        candidate["fields"].append({"id": name, "value": values[0]})
    return {
        "status": "CANDIDATE" if not errors else "REJECTED",
        "normalization_codes": sorted(set(errors)),
        "candidate": candidate if not errors else None,
        "provider_lead_id_sha256": sha256(candidate["provider_lead_id"].encode("utf-8")),
    }


def fetch_to_evidence(
    form_factory: Callable[[str], Any], tenant_id: str, organization_id: str,
    form_id: str, fields: list[str], approved_names: set[str], max_rows: int,
    retrieval_mode: str, output_directory: Path,
) -> dict[str, Any]:
    if importlib.metadata.version("facebook-business") != SDK_VERSION:
        raise RuntimeError(f"facebook-business {SDK_VERSION} is required")
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".meta-lead-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        form = form_factory(form_id)
        method = form.get_test_leads if retrieval_mode == "TEST" else form.get_leads
        rows: list[dict[str, Any]] = []
        truncated = False
        for raw in method(fields=fields, params={}):
            if len(rows) == max_rows:
                truncated = True
                break
            rows.append(_row_to_dict(raw))
        response = {
            "graph_api_version": API_VERSION,
            "retrieval_mode": retrieval_mode,
            "rows": rows,
            "truncated_by_local_bound": truncated,
        }
        response_bytes = (json.dumps(response, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "provider-response.json").write_bytes(response_bytes)
        outcomes = [_normalize_row(row, tenant_id, organization_id, approved_names, retrieval_mode == "TEST") for row in rows]
        candidate_document = {
            "schema": "elite-meta-lead-candidate-batch/v1",
            "provider": "meta_lead_ads",
            "source_commit": SOURCE_COMMIT,
            "tenant_id": tenant_id,
            "organization_id": organization_id,
            "provider_response_sha256": sha256(response_bytes),
            "outcomes": outcomes,
            "automatic_business_write": False,
            "automatic_contact_eligibility": False,
        }
        candidate_bytes = (json.dumps(candidate_document, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "lead-candidates.json").write_bytes(candidate_bytes)
        query_bytes = json.dumps({"fields": fields, "max_rows": max_rows, "retrieval_mode": retrieval_mode}, sort_keys=True, separators=(",", ":")).encode("utf-8")
        receipt = {
            "schema": "elite-meta-lead-reconciliation-receipt/v1",
            "created_at": datetime.now(timezone.utc).isoformat(),
            "provider": "Meta Marketing API",
            "sdk_version": SDK_VERSION,
            "graph_api_version": API_VERSION,
            "source_commit": SOURCE_COMMIT,
            "form_id_sha256": sha256(form_id.encode("ascii")),
            "query_sha256": sha256(query_bytes),
            "provider_response_sha256": sha256(response_bytes),
            "candidate_batch_sha256": sha256(candidate_bytes),
            "row_count": len(rows),
            "candidate_count": sum(item["status"] == "CANDIDATE" for item in outcomes),
            "rejected_count": sum(item["status"] == "REJECTED" for item in outcomes),
            "truncated_by_local_bound": truncated,
            "automatic_business_write": False,
            "automatic_contact_eligibility": False,
        }
        (stage / "RETRIEVAL_RECEIPT.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        stage.replace(output)
        return receipt
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--retrieval", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()
    profile, retrieval = _read_json(args.profile), _read_json(args.retrieval)
    tenant_id, organization_id, form_id, fields, approved_names, max_rows = validate_configuration(profile, retrieval)
    names = ("app_id_environment_variable", "app_secret_environment_variable", "access_token_environment_variable")
    environment_names = [str(profile.get(name, "")) for name in names]
    if any(not re.fullmatch(r"[A-Z][A-Z0-9_]{2,80}", name) for name in environment_names):
        raise ValueError("secret references must be canonical environment variable names")
    values = [os.environ.get(name, "") for name in environment_names]
    if any(not value for value in values):
        raise PermissionError("required Meta credential environment variables are unavailable")
    api = FacebookAdsApi.init(app_id=values[0], app_secret=values[1], access_token=values[2], api_version=API_VERSION)
    receipt = fetch_to_evidence(
        lambda identifier: LeadgenForm(identifier, api=api), tenant_id, organization_id,
        form_id, fields, approved_names, max_rows, retrieval["retrieval_mode"], args.output,
    )
    print(json.dumps({"status": "META_LEAD_RECONCILIATION_PASS", "receipt": receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
