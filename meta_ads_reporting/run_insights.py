from __future__ import annotations

import argparse
from datetime import date, datetime, timezone
import hashlib
import importlib.metadata
import json
import os
from pathlib import Path
import re
import shutil
import sys
from typing import Any, Callable
import uuid

from facebook_business.adobjects.adaccount import AdAccount
from facebook_business.api import FacebookAdsApi


SDK_VERSION = "26.0.1"
API_VERSION = "v26.0"
ALLOWED_FIELDS = {
    "account_id", "campaign_id", "campaign_name", "adset_id", "adset_name",
    "ad_id", "ad_name", "date_start", "date_stop", "impressions", "reach",
    "clicks", "spend",
}
ALLOWED_LEVELS = {"account", "campaign", "adset", "ad"}


def sha256(value: bytes) -> str:
    return hashlib.sha256(value).hexdigest()


def _read_json(path: Path) -> dict[str, Any]:
    value = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(value, dict):
        raise ValueError(f"{path.name} must contain one JSON object")
    return value


def validate_configuration(profile: dict[str, Any], query: dict[str, Any]) -> tuple[str, list[str], dict[str, Any]]:
    required = (
        "app_registration_proven", "business_verification_proven", "ads_read_permission_proven",
        "test_ad_account_contract_proven", "quota_and_cost_approved", "reconciliation_approved",
    )
    if profile.get("provider") != "Meta Marketing API" or profile.get("sdk") != f"facebook-business=={SDK_VERSION}" or profile.get("graph_api_version") != API_VERSION:
        raise ValueError("profile provider/SDK/API identity mismatch")
    if profile.get("decision") != "PROVEN" or any(profile.get(key) is not True for key in required):
        raise PermissionError("Meta access/query profile remains blocked")
    if profile.get("automatic_business_write") is not False:
        raise ValueError("automatic_business_write must remain false")
    if not str(profile.get("data_retention", "")).strip():
        raise ValueError("data_retention is required")
    account_id = str(profile.get("ad_account_id", ""))
    if not re.fullmatch(r"[0-9]{5,30}", account_id):
        raise ValueError("ad_account_id must contain digits only")
    fields = query.get("fields")
    approved = profile.get("approved_fields")
    if not isinstance(fields, list) or not fields or len(fields) > 20 or len(fields) != len(set(fields)):
        raise ValueError("fields must be a unique non-empty list with at most 20 entries")
    if not isinstance(approved, list) or not set(fields).issubset(set(approved)) or not set(fields).issubset(ALLOWED_FIELDS):
        raise PermissionError("all fields must be locally allowed and explicitly approved")
    level = query.get("level")
    if level not in ALLOWED_LEVELS:
        raise ValueError("unsupported reporting level")
    since = date.fromisoformat(str(query.get("since", "")))
    until = date.fromisoformat(str(query.get("until", "")))
    if until < since or (until - since).days > 31:
        raise ValueError("time range must be ordered and at most 31 days")
    limit = query.get("limit")
    if not isinstance(limit, int) or isinstance(limit, bool) or limit < 1 or limit > 1000:
        raise ValueError("limit must be within 1..1000")
    params = {"level": level, "time_range": {"since": since.isoformat(), "until": until.isoformat()}, "limit": limit}
    return account_id, fields, params


def _row_to_dict(row: Any) -> dict[str, Any]:
    if isinstance(row, dict):
        return row
    exporter = getattr(row, "export_all_data", None)
    if not callable(exporter):
        raise RuntimeError("Meta Ads result row is not an official SDK object")
    value = exporter()
    if not isinstance(value, dict):
        raise RuntimeError("Meta SDK row export must return an object")
    return value


def query_to_evidence(account_factory: Callable[[str], Any], account_id: str, fields: list[str], params: dict[str, Any], output_directory: Path) -> dict[str, Any]:
    if importlib.metadata.version("facebook-business") != SDK_VERSION:
        raise RuntimeError(f"facebook-business {SDK_VERSION} is required")
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".meta-ads-report-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        rows = [_row_to_dict(row) for row in account_factory(f"act_{account_id}").get_insights(fields=fields, params=params)]
        response = {"graph_api_version": API_VERSION, "rows": rows}
        response_bytes = (json.dumps(response, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "provider-response.json").write_bytes(response_bytes)
        query_bytes = json.dumps({"fields": fields, "params": params}, sort_keys=True, separators=(",", ":")).encode("utf-8")
        receipt = {
            "schema": "elite-meta-ads-reporting-receipt/v1",
            "created_at": datetime.now(timezone.utc).isoformat(),
            "provider": "Meta Marketing API",
            "sdk_version": SDK_VERSION,
            "graph_api_version": API_VERSION,
            "ad_account_id_sha256": sha256(account_id.encode("ascii")),
            "query_sha256": sha256(query_bytes),
            "row_count": len(rows),
            "provider_response_sha256": sha256(response_bytes),
            "automatic_business_write": False,
        }
        (stage / "QUERY_RECEIPT.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        stage.replace(output)
        return receipt
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--query", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()
    profile, query = _read_json(args.profile), _read_json(args.query)
    account_id, fields, params = validate_configuration(profile, query)
    names = ("app_id_environment_variable", "app_secret_environment_variable", "access_token_environment_variable")
    environment_names = [str(profile.get(name, "")) for name in names]
    if any(not re.fullmatch(r"[A-Z][A-Z0-9_]{2,80}", name) for name in environment_names):
        raise ValueError("secret references must be canonical environment variable names")
    values = [os.environ.get(name, "") for name in environment_names]
    if any(not value for value in values):
        raise PermissionError("required Meta credential environment variables are unavailable")
    api = FacebookAdsApi.init(app_id=values[0], app_secret=values[1], access_token=values[2], api_version=API_VERSION)
    receipt = query_to_evidence(lambda identifier: AdAccount(identifier, api=api), account_id, fields, params, args.output)
    print(json.dumps({"status": "META_ADS_REPORTING_PASS", "receipt": receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
