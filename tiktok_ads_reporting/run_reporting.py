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


SDK_DISTRIBUTION = "tiktok-business-api-sdk-official"
SDK_VERSION = "1.1.3"
ALLOWED_LEVELS = {"AUCTION_ADVERTISER", "AUCTION_CAMPAIGN", "AUCTION_ADGROUP", "AUCTION_AD"}


def sha256(value: bytes) -> str:
    return hashlib.sha256(value).hexdigest()


def _read_object(path: Path) -> dict[str, Any]:
    value = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(value, dict):
        raise ValueError(f"{path.name} must contain one JSON object")
    return value


def validate_configuration(profile: dict[str, Any], query: dict[str, Any]) -> tuple[str, dict[str, Any]]:
    required = ("business_account_proven", "developer_app_proven", "reporting_scope_proven", "test_advertiser_contract_proven", "quota_and_cost_approved", "reconciliation_approved")
    if profile.get("provider") != "TikTok Business API" or profile.get("sdk_distribution") != SDK_DISTRIBUTION or profile.get("sdk_version") != SDK_VERSION:
        raise ValueError("TikTok profile authority mismatch")
    if profile.get("decision") != "PROVEN" or any(profile.get(key) is not True for key in required):
        raise PermissionError("TikTok access/query profile remains blocked")
    if profile.get("automatic_business_write") is not False or not str(profile.get("data_retention", "")).strip():
        raise ValueError("retention is required and automatic_business_write must remain false")
    advertiser_id = str(profile.get("advertiser_id", ""))
    if not re.fullmatch(r"[0-9]{5,30}", advertiser_id):
        raise ValueError("advertiser_id must contain digits only")
    report_type, level = query.get("report_type"), query.get("data_level")
    if report_type != "BASIC" or level not in ALLOWED_LEVELS:
        raise ValueError("only BASIC reporting at an admitted auction level is supported")
    identifier = re.compile(r"[a-z][a-z0-9_]{0,79}")
    dimensions, metrics = query.get("dimensions"), query.get("metrics")
    for name, values, approved, maximum in (("dimensions", dimensions, profile.get("approved_dimensions"), 10), ("metrics", metrics, profile.get("approved_metrics"), 30)):
        if not isinstance(values, list) or not values or len(values) > maximum or len(values) != len(set(values)) or any(not isinstance(value, str) or not identifier.fullmatch(value) for value in values):
            raise ValueError(f"{name} must be a bounded unique identifier list")
        if not isinstance(approved, list) or not set(values).issubset(set(approved)):
            raise PermissionError(f"all {name} must be explicitly approved")
    start, end = date.fromisoformat(str(query.get("start_date", ""))), date.fromisoformat(str(query.get("end_date", "")))
    if end < start or (end - start).days > 31:
        raise ValueError("report range must be ordered and at most 31 days")
    page_size = query.get("page_size")
    if not isinstance(page_size, int) or isinstance(page_size, bool) or not 1 <= page_size <= 1000:
        raise ValueError("page_size must be within 1..1000")
    return advertiser_id, {"report_type": report_type, "data_level": level, "dimensions": dimensions, "metrics": metrics, "start_date": start.isoformat(), "end_date": end.isoformat(), "page_size": page_size}


def _response_to_dict(response: Any) -> dict[str, Any]:
    if isinstance(response, dict):
        value = response
    else:
        exporter = getattr(response, "to_dict", None)
        if not callable(exporter):
            raise RuntimeError("TikTok response is not an official SDK response")
        value = exporter()
    if not isinstance(value, dict):
        raise RuntimeError("TikTok response export must be an object")
    return value


def query_to_evidence(service: Any, access_token: str, advertiser_id: str, query: dict[str, Any], output_directory: Path) -> dict[str, Any]:
    if not access_token:
        raise PermissionError("TikTok access token is unavailable")
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".tiktok-report-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        pages, row_count, page = [], 0, 1
        while True:
            response = service.report_integrated_get(query["report_type"], access_token, advertiser_id=advertiser_id, data_level=query["data_level"], dimensions=query["dimensions"], metrics=query["metrics"], start_date=query["start_date"], end_date=query["end_date"], page=page, page_size=query["page_size"])
            payload = _response_to_dict(response)
            if payload.get("code") not in (0, None):
                raise RuntimeError(f"TikTok provider rejected report with code {payload.get('code')}")
            pages.append(payload)
            data = payload.get("data") if isinstance(payload.get("data"), dict) else {}
            rows = data.get("list") if isinstance(data.get("list"), list) else []
            row_count += len(rows)
            page_info = data.get("page_info") if isinstance(data.get("page_info"), dict) else {}
            total_page = page_info.get("total_page", page)
            if not isinstance(total_page, int) or total_page < page or total_page > 100:
                raise RuntimeError("TikTok page_info is invalid or exceeds 100-page safety bound")
            if page >= total_page:
                break
            page += 1
        response_bytes = (json.dumps({"sdk_distribution": SDK_DISTRIBUTION, "sdk_version": SDK_VERSION, "pages": pages}, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "provider-response.json").write_bytes(response_bytes)
        query_bytes = json.dumps(query, sort_keys=True, separators=(",", ":")).encode("utf-8")
        receipt = {"schema": "elite-tiktok-ads-reporting-receipt/v1", "created_at": datetime.now(timezone.utc).isoformat(), "provider": "TikTok Business API", "sdk_distribution": SDK_DISTRIBUTION, "sdk_version": SDK_VERSION, "advertiser_id_sha256": sha256(advertiser_id.encode("ascii")), "query_sha256": sha256(query_bytes), "page_count": len(pages), "row_count": row_count, "provider_response_sha256": sha256(response_bytes), "automatic_business_write": False}
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
    if importlib.metadata.version(SDK_DISTRIBUTION) != SDK_VERSION:
        raise RuntimeError(f"{SDK_DISTRIBUTION} {SDK_VERSION} is required")
    from business_api_client.api.reporting_api import ReportingApi
    profile, query_document = _read_object(args.profile), _read_object(args.query)
    advertiser_id, query = validate_configuration(profile, query_document)
    secret_name = str(profile.get("access_token_environment_variable", ""))
    if not re.fullmatch(r"[A-Z][A-Z0-9_]{2,80}", secret_name):
        raise ValueError("access token reference must be a canonical environment variable name")
    receipt = query_to_evidence(ReportingApi(), os.environ.get(secret_name, ""), advertiser_id, query, args.output)
    print(json.dumps({"status": "TIKTOK_ADS_REPORTING_PASS", "receipt": receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
