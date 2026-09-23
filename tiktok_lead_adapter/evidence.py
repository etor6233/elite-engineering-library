"""Hash-linked TikTok Lead retrieval artifacts without automatic persistence."""

from __future__ import annotations

import hashlib
import json
import os
import re
import shutil
import uuid
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Dict, Mapping, Tuple


SDK_VERSION = "1.1.3"
SDK_WHEEL_SHA256 = "663b4a1f2585c4b386144d968646f695befd733c464418bbaeeac41befff33b7"
API_VERSION = "v1.3"
CONTRACT_OBSERVED_AT = "2026-09-04"
PROVIDER = "tiktok_lead_generation"
SAFE_ID = re.compile(r"^[A-Za-z0-9_.:@/+ -]{1,256}$")
SAFE_FIELD = re.compile(r"^[^\x00-\x1f]{1,256}$")


class TikTokLeadEvidenceError(ValueError):
    pass


def _canonical(value: Mapping[str, Any]) -> bytes:
    return (json.dumps(value, ensure_ascii=False, separators=(",", ":"), sort_keys=True) + "\n").encode("utf-8")


def _sha(value: bytes) -> str:
    return hashlib.sha256(value).hexdigest()


def _identifier(value: object, name: str) -> str:
    if not isinstance(value, (str, int)):
        raise TikTokLeadEvidenceError(f"{name} is not a string-compatible identifier")
    result = str(value).strip()
    if not SAFE_ID.fullmatch(result):
        raise TikTokLeadEvidenceError(f"{name} is empty, too long or contains unsafe characters")
    return result


def validate_request(request: Mapping[str, Any]) -> Dict[str, Any]:
    if request.get("schema") != "elite-tiktok-lead-retrieval-request/v1":
        raise TikTokLeadEvidenceError("unsupported retrieval request schema")
    try:
        tenant = str(uuid.UUID(str(request.get("tenant_id", ""))))
    except ValueError:
        raise TikTokLeadEvidenceError("tenant_id must be a canonical UUID") from None
    organization = _identifier(request.get("organization_id"), "organization_id")
    mode = request.get("mode")
    if mode not in ("TEST", "LIVE"):
        raise TikTokLeadEvidenceError("mode must be TEST or LIVE")
    lead_source = request.get("lead_source")
    if lead_source not in ("INSTANT_FORM", "DIRECT_MESSAGE"):
        raise TikTokLeadEvidenceError("lead_source is not admitted")
    return {
        "schema": request["schema"],
        "tenant_id": tenant,
        "organization_id": organization,
        "mode": mode,
        "lead_source": lead_source,
        "advertiser_id": request.get("advertiser_id") or None,
        "library_id": request.get("library_id") or None,
        "page_id": request.get("page_id") or None,
    }


def build_retrieval_artifacts(result: Mapping[str, Any], request_input: Mapping[str, Any]) -> Dict[str, bytes]:
    request = validate_request(request_input)
    if not isinstance(result, dict) or not isinstance(result.get("data"), dict):
        raise TikTokLeadEvidenceError("invalid official SDK response")
    data = result["data"]
    lead_data = data.get("lead_data")
    meta = data.get("meta_data")
    request_id = result.get("request_id")
    if not isinstance(lead_data, dict) or not isinstance(meta, dict) or not isinstance(request_id, str) or not request_id:
        raise TikTokLeadEvidenceError("response lacks lead_data, meta_data or request_id")
    lead_id = _identifier(meta.get("lead_id"), "meta_data.lead_id")
    lead_source = meta.get("lead_source")
    if lead_source != request["lead_source"]:
        raise TikTokLeadEvidenceError("response lead_source does not match request")
    page_id = meta.get("page_id")
    if lead_source == "INSTANT_FORM":
        if _identifier(page_id, "meta_data.page_id") != _identifier(request["page_id"], "request.page_id"):
            raise TikTokLeadEvidenceError("response page_id does not match request")
    elif page_id not in (None, ""):
        raise TikTokLeadEvidenceError("DIRECT_MESSAGE response unexpectedly contains page_id")
    created = meta.get("create_time")
    if not isinstance(created, str):
        raise TikTokLeadEvidenceError("meta_data.create_time is required")
    try:
        submitted = datetime.strptime(created, "%Y-%m-%d %H:%M:%S").replace(tzinfo=timezone.utc)
    except ValueError:
        raise TikTokLeadEvidenceError("meta_data.create_time is not official UTC format") from None

    normalized_fields = []
    rejected = []
    for key in sorted(lead_data):
        if not isinstance(key, str) or not SAFE_FIELD.fullmatch(key):
            rejected.append("UNSAFE_DYNAMIC_FIELD_NAME")
            continue
        value = lead_data[key]
        if not isinstance(value, str) or not value.strip():
            rejected.append("NON_STRING_DYNAMIC_FIELD_REQUIRES_MAPPING")
            continue
        normalized_fields.append({"id": key, "value": value})
    if not lead_data:
        rejected.append("EMPTY_LEAD_DATA")

    provider_response = {
        "schema": "elite-tiktok-lead-sdk-response/v1",
        "sdk_version": SDK_VERSION,
        "api_version": API_VERSION,
        "request_identity": {
            "mode": request["mode"],
            "lead_source": request["lead_source"],
            "advertiser_id_sha256": "" if request["advertiser_id"] is None else _sha(str(request["advertiser_id"]).encode("utf-8")),
            "library_id_sha256": "" if request["library_id"] is None else _sha(str(request["library_id"]).encode("utf-8")),
            "page_id_sha256": "" if request["page_id"] is None else _sha(str(request["page_id"]).encode("utf-8")),
        },
        "data": data,
        "request_id": request_id,
    }
    provider_bytes = _canonical(provider_response)
    status = "REJECTED" if rejected else "CANDIDATE"
    candidate = None
    if status == "CANDIDATE":
        candidate = {
            "tenant_id": request["tenant_id"],
            "organization_id": request["organization_id"],
            "provider": PROVIDER,
            "provider_lead_id": lead_id,
            "form_id": "" if page_id is None else str(page_id),
            "campaign_id": "" if meta.get("campaign_id") is None else _identifier(meta.get("campaign_id"), "campaign_id"),
            "ad_group_id": "" if meta.get("adgroup_id") is None else _identifier(meta.get("adgroup_id"), "adgroup_id"),
            "creative_id": "" if meta.get("ad_id") is None else _identifier(meta.get("ad_id"), "ad_id"),
            "source_kind": lead_source,
            "submitted_at": submitted.isoformat().replace("+00:00", "Z"),
            "is_test": request["mode"] == "TEST",
            "fields": normalized_fields,
            "contact_eligibility": "pending_policy",
        }
    outcome = {
        "status": status,
        "normalization_codes": sorted(set(rejected)),
        "provider_lead_id_sha256": _sha(lead_id.encode("utf-8")),
        "candidate": candidate,
    }
    batch = {
        "schema": "elite-tiktok-lead-candidate-batch/v1",
        "provider": PROVIDER,
        "sdk_version": SDK_VERSION,
        "sdk_wheel_sha256": SDK_WHEEL_SHA256,
        "api_version": API_VERSION,
        "contract_observed_at": CONTRACT_OBSERVED_AT,
        "tenant_id": request["tenant_id"],
        "organization_id": request["organization_id"],
        "provider_response_sha256": _sha(provider_bytes),
        "outcomes": [outcome],
        "automatic_business_write": False,
        "automatic_contact_eligibility": False,
    }
    batch_bytes = _canonical(batch)
    receipt = {
        "schema": "elite-tiktok-lead-retrieval-receipt/v1",
        "sdk_version": SDK_VERSION,
        "sdk_wheel_sha256": SDK_WHEEL_SHA256,
        "api_version": API_VERSION,
        "contract_observed_at": CONTRACT_OBSERVED_AT,
        "retrieval_mode": request["mode"],
        "request_id_sha256": _sha(request_id.encode("utf-8")),
        "provider_response_sha256": _sha(provider_bytes),
        "candidate_batch_sha256": _sha(batch_bytes),
        "row_count": 1,
        "candidate_count": 1 if status == "CANDIDATE" else 0,
        "rejected_count": 1 if status == "REJECTED" else 0,
        "automatic_business_write": False,
        "automatic_contact_eligibility": False,
    }
    return {
        "provider-response.json": provider_bytes,
        "candidate-batch.json": batch_bytes,
        "retrieval-receipt.json": _canonical(receipt),
    }


def write_artifacts_atomically(output: Path, artifacts: Mapping[str, bytes]) -> None:
    output = output.resolve()
    if output.exists():
        raise TikTokLeadEvidenceError("output path must not exist")
    staging = output.parent / ("." + output.name + ".tmp-" + uuid.uuid4().hex)
    if staging.exists():
        raise TikTokLeadEvidenceError("staging path collision")
    staging.mkdir(parents=False)
    try:
        for name in ("provider-response.json", "candidate-batch.json", "retrieval-receipt.json"):
            value = artifacts.get(name)
            if not isinstance(value, bytes) or not value:
                raise TikTokLeadEvidenceError(f"artifact is missing: {name}")
            path = staging / name
            with path.open("xb") as handle:
                handle.write(value)
                handle.flush()
                os.fsync(handle.fileno())
        staging.replace(output)
    except Exception:
        shutil.rmtree(staging, ignore_errors=True)
        raise
