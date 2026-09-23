"""Execute one approved TikTok Lead retrieval and emit three atomic artifacts."""

from __future__ import annotations

import argparse
import json
import os
from pathlib import Path

from business_api_client.api_client import ApiClient

from adapter import TikTokLeadClient
from evidence import build_retrieval_artifacts, validate_request, write_artifacts_atomically


def _load(path: Path):
    return json.loads(path.read_text(encoding="utf-8"))


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--request", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()
    profile = _load(args.profile)
    request = validate_request(_load(args.request))
    token_name = profile.get("access_token_environment_variable")
    if not isinstance(token_name, str) or not token_name:
        raise RuntimeError("profile lacks access token environment variable name")
    token = os.environ.get(token_name, "")
    client = TikTokLeadClient(ApiClient(), profile)
    result = client.get_lead(
        token,
        request["lead_source"],
        advertiser_id=request["advertiser_id"],
        library_id=request["library_id"],
        page_id=request["page_id"],
    )
    artifacts = build_retrieval_artifacts(result, request)
    write_artifacts_atomically(args.output, artifacts)
    receipt = json.loads(artifacts["retrieval-receipt.json"])
    print(json.dumps(receipt, separators=(",", ":"), sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
