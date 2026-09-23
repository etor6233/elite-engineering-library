from __future__ import annotations

import argparse
from datetime import datetime, timezone
import hashlib
import importlib.metadata
import json
from pathlib import Path
import re
import shutil
import sys
import uuid
from typing import Any, Iterable, Protocol

from google.ads.googleads.client import GoogleAdsClient
from google.protobuf.json_format import MessageToDict


SDK_VERSION = "31.4.0"
API_VERSION = "v25"


class GoogleAdsService(Protocol):
    def search_stream(self, *, customer_id: str, query: str, timeout: float) -> Iterable[Any]: ...


def sha256(value: bytes) -> str:
    return hashlib.sha256(value).hexdigest()


def row_to_dict(row: Any) -> dict[str, Any]:
    if isinstance(row, dict):
        return row
    protobuf = getattr(row, "_pb", None)
    if protobuf is None:
        raise RuntimeError("Google Ads result row is not an official protobuf message")
    return MessageToDict(protobuf, preserving_proto_field_name=True)


def query_to_evidence(service: GoogleAdsService, customer_id: str, query: str, output_directory: Path, timeout_seconds: float = 60.0) -> dict[str, Any]:
    if importlib.metadata.version("google-ads") != SDK_VERSION:
        raise RuntimeError(f"google-ads {SDK_VERSION} is required")
    if not re.fullmatch(r"[0-9]{5,20}", customer_id):
        raise ValueError("customer_id must contain digits only")
    normalized = " ".join(query.split())
    if not normalized.upper().startswith("SELECT ") or len(normalized) > 10000 or ";" in normalized:
        raise ValueError("one bounded GAQL SELECT query is required")
    if timeout_seconds <= 0 or timeout_seconds > 300:
        raise ValueError("timeout must be within 1..300 seconds")
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".google-ads-report-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        rows: list[dict[str, Any]] = []
        for batch in service.search_stream(customer_id=customer_id, query=normalized, timeout=timeout_seconds):
            for row in getattr(batch, "results", []):
                rows.append(row_to_dict(row))
        response = {"api_version": API_VERSION, "rows": rows}
        response_bytes = (json.dumps(response, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "provider-response.json").write_bytes(response_bytes)
        receipt = {
            "schema": "elite-google-ads-reporting-receipt/v1",
            "created_at": datetime.now(timezone.utc).isoformat(),
            "provider": "Google Ads API",
            "sdk_version": SDK_VERSION,
            "api_version": API_VERSION,
            "customer_id_sha256": sha256(customer_id.encode("ascii")),
            "query_sha256": sha256(normalized.encode("utf-8")),
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
    parser.add_argument("--customer-id", required=True)
    parser.add_argument("--query-file", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--timeout-seconds", type=float, default=60.0)
    args = parser.parse_args()
    query = args.query_file.read_text(encoding="utf-8")
    client = GoogleAdsClient.load_from_env(version=API_VERSION)
    service = client.get_service("GoogleAdsService", version=API_VERSION)
    receipt = query_to_evidence(service, args.customer_id, query, args.output, args.timeout_seconds)
    print(json.dumps({"status": "GOOGLE_ADS_REPORTING_PASS", "receipt": receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
