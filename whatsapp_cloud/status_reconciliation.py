"""AUTHORED correlation of an anchored send receipt with signed Meta statuses.

Reuses the existing Meta-adapted signature/scope boundary. Never sends messages
or changes the PostgreSQL fence. Hash anchors must come from trusted storage,
not from the webhook caller or a newly computed hash of an untrusted file.
"""
from __future__ import annotations

from datetime import datetime, timezone
import hmac
import json
from pathlib import Path
import re
import shutil
from typing import Any
import uuid
import base64
import tempfile

from whatsapp_cloud import SOURCE_COMMIT, _sha256, _unique_object, normalize_verified_webhook, validate_profile


def status_bridge(raw: bytes) -> dict[str, Any]:
    """Internal isolated-process contract. Secrets and raw bodies stay in memory."""
    if not 0 < len(raw) <= 12 * 1024 * 1024:
        raise ValueError("status frame exceeds local budget")
    frame = json.loads(raw.decode("utf-8"), object_pairs_hook=_unique_object)
    fields = {"schema", "profile", "send_receipt", "expected_send_receipt_sha256", "webhooks", "app_secret", "evidence_directory"}
    if not isinstance(frame, dict) or set(frame) != fields or frame["schema"] != "elite-whatsapp-status-bridge/v1":
        raise ValueError("status frame contract mismatch")
    if not isinstance(frame["app_secret"], str) or not 0 < len(frame["app_secret"]) <= 16384:
        raise ValueError("app secret unavailable")
    if not isinstance(frame["evidence_directory"], str) or not Path(frame["evidence_directory"]).is_absolute():
        raise ValueError("absolute evidence directory required")
    batches = frame["webhooks"]
    if not isinstance(batches, list) or not 1 <= len(batches) <= 8:
        raise ValueError("bounded webhook batch required")
    signed = []
    for item in batches:
        if not isinstance(item, dict) or set(item) != {"body", "signature"} or not isinstance(item["body"], str):
            raise ValueError("status batch contract mismatch")
        signed.append((base64.b64decode(item["body"], validate=True), item["signature"]))
    from whatsapp_cloud import verify_packaged_official_source
    root = Path(__file__).resolve().parent
    verify_packaged_official_source(root, root / "official-source.lock.json")
    with tempfile.TemporaryDirectory(prefix="whatsapp-status-", dir=frame["evidence_directory"]) as temp:
        output = Path(temp) / "verified"
        receipt = reconcile_status_webhooks(profile=frame["profile"],
            send_receipt_bytes=base64.b64decode(frame["send_receipt"], validate=True),
            expected_send_receipt_sha256=frame["expected_send_receipt_sha256"],
            signed_webhooks=signed, app_secret=frame["app_secret"], output_directory=output)
        observations = json.loads((output / "status-observations.json").read_bytes())
        return {"schema": "elite-whatsapp-status-result/v1", "binding_sha256": _sha256(raw),
                "receipt": receipt, "events": observations["events"]}


def reconcile_status_webhooks(*, profile: dict[str, Any], send_receipt_bytes: bytes,
                             expected_send_receipt_sha256: str,
                             signed_webhooks: list[tuple[bytes, str]], app_secret: str,
                             output_directory: Path) -> dict[str, Any]:
    """Publish a bounded observation snapshot, not a global delivery truth.

    Supply SEND_RECEIPT.json bytes and the EvidenceSHA256 already stored by the
    Go outbound fence. Each webhook is verified on its original bytes, including
    unrelated events in a shared batch. The snapshot covers only these inputs.
    """
    validate_profile(profile)
    if not isinstance(send_receipt_bytes, bytes) or not 0 < len(send_receipt_bytes) <= 65536:
        raise ValueError("send receipt exceeds local budget")
    if not isinstance(expected_send_receipt_sha256, str) or not re.fullmatch(r"[0-9a-f]{64}", expected_send_receipt_sha256):
        raise ValueError("trusted send receipt anchor is required")
    if not hmac.compare_digest(_sha256(send_receipt_bytes), expected_send_receipt_sha256):
        raise PermissionError("send receipt anchor mismatch")
    try:
        send = json.loads(send_receipt_bytes.decode("utf-8"), object_pairs_hook=_unique_object)
    except (UnicodeDecodeError, json.JSONDecodeError, RecursionError) as error:
        raise ValueError("send receipt must be UTF-8 JSON") from error
    if not isinstance(send, dict) or send.get("schema") != "elite-whatsapp-cloud-send-receipt/v1" or send.get("source_commit") != SOURCE_COMMIT or send.get("automatic_business_write") is not False:
        raise ValueError("send receipt contract mismatch")
    for key in ("phone_number_id_sha256", "recipient_sha256", "request_sha256", "response_sha256", "message_id_sha256"):
        if not isinstance(send.get(key), str) or not re.fullmatch(r"[0-9a-f]{64}", send[key]):
            raise ValueError("send receipt digest is invalid")
    if send.get("graph_api_version") != profile["graph_api_version"] or send["phone_number_id_sha256"] != _sha256(profile["phone_number_id"].encode("ascii")):
        raise PermissionError("send receipt profile mismatch")
    if not isinstance(signed_webhooks, list) or not 1 <= len(signed_webhooks) <= 8:
        raise ValueError("one to eight signed webhook bodies are required")

    matched: dict[str, dict[str, Any]] = {}
    raw_hashes: set[str] = set()
    excluded = 0
    duplicate = 0
    for item in signed_webhooks:
        if not isinstance(item, tuple) or len(item) != 2 or not isinstance(item[1], str):
            raise ValueError("webhook input must be raw bytes and signature")
        raw, signature = item
        scope, events = normalize_verified_webhook(raw, signature, app_secret, profile=profile)
        raw_hashes.add(_sha256(raw))
        for event in events:
            if event["kind"] != "status" or event["id_sha256"] != send["message_id_sha256"]:
                excluded += 1
                continue
            if event["recipient_sha256"] != send["recipient_sha256"]:
                raise PermissionError("matched provider message has another recipient")
            key = event["event_key_sha256"]
            if key in matched:
                if matched[key]["event_sha256"] != event["event_sha256"]:
                    raise ValueError("same status identity has divergent provider evidence")
                duplicate += 1
                continue
            matched[key] = event

    # Meta documents out-of-order delivery. Sort by provider time, never arrival.
    # Equal-time differing states remain explicit; do not invent precedence.
    timeline = sorted(matched.values(), key=lambda event: (int(event["timestamp"]), event["event_key_sha256"]))
    last_time = timeline[-1]["timestamp"] if timeline else None
    last_states = sorted({event["status"] for event in timeline if event["timestamp"] == last_time})
    decision = "NOT_OBSERVED_IN_INPUT" if not timeline else "AMBIGUOUS_LATEST_TIMESTAMP" if len(last_states) != 1 else "MATCHED_OBSERVATIONS"
    observations = {"schema": "elite-whatsapp-status-observations/v1", "events": timeline}
    normalized = (json.dumps(observations, sort_keys=True, indent=2) + "\n").encode("utf-8")
    receipt = {
        "schema": "elite-whatsapp-status-reconciliation/v1", "source_commit": SOURCE_COMMIT,
        "created_at": datetime.now(timezone.utc).isoformat(), **scope,
        "send_receipt_sha256": expected_send_receipt_sha256,
        "message_id_sha256": send["message_id_sha256"], "recipient_sha256": send["recipient_sha256"],
        "profile_sha256": _sha256(json.dumps(profile, sort_keys=True, separators=(",", ":")).encode("utf-8")),
        "webhook_body_sha256": sorted(raw_hashes), "observations_sha256": _sha256(normalized),
        "decision": decision, "last_observed_timestamp": last_time,
        "last_observed_status": last_states[0] if len(last_states) == 1 else None,
        "observed_statuses": sorted({event["status"] for event in timeline}),
        "matching_unique_events": len(timeline), "duplicate_events": duplicate,
        "excluded_events": excluded, "raw_payload_persisted": False,
        "automatic_business_write": False, "resend_authorized": False,
    }
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".whatsapp-status-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        (stage / "status-observations.json").write_bytes(normalized)
        (stage / "STATUS_RECONCILIATION.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        stage.replace(output)
        return receipt
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise
