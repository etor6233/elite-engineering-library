from __future__ import annotations

import argparse
import base64
import hashlib
import json
import os
import re
import shutil
import tempfile
from dataclasses import dataclass
from email import policy
from email.message import Message
from email.parser import BytesParser
from pathlib import Path
from typing import Any


class IntakeRejected(RuntimeError):
    def __init__(self, code: str):
        super().__init__(code)
        self.code = code


@dataclass(frozen=True)
class Limits:
    max_raw_bytes: int = 36 * 1024 * 1024
    max_parts: int = 100
    max_depth: int = 12
    max_headers_per_part: int = 100
    max_header_chars: int = 8192
    max_attachments: int = 20
    max_attachment_bytes: int = 20 * 1024 * 1024
    max_total_attachment_bytes: int = 30 * 1024 * 1024


TENANT_RE = re.compile(r"^[a-z0-9][a-z0-9_-]{0,62}$")
HEX64_RE = re.compile(r"^[0-9a-f]{64}$")
SOURCE_TOKEN_RE = re.compile(r"^[A-Za-z0-9._:/+=-]{1,1024}$")
BASE64_WS_RE = re.compile(rb"[\t\r\n ]+")


def canonical_json(value: Any) -> bytes:
    return (json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(",", ":")) + "\n").encode("utf-8")


def sha256(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def _required_string(obj: dict[str, Any], field: str, pattern: re.Pattern[str] | None = None) -> str:
    value = obj.get(field)
    if not isinstance(value, str) or not value or (pattern and not pattern.fullmatch(value)):
        raise IntakeRejected(f"INVALID_RECEIPT_{field.upper()}")
    return value


def validate_receipt(receipt: dict[str, Any], raw: bytes) -> dict[str, str]:
    if receipt.get("receipt_version") != 1 or receipt.get("decision") != "ADMIT":
        raise IntakeRejected("RECEIPT_NOT_ADMITTED")
    tenant_id = _required_string(receipt, "tenant_id", TENANT_RE)
    raw_sha = _required_string(receipt, "raw_sha256", HEX64_RE)
    if raw_sha != sha256(raw):
        raise IntakeRejected("RAW_SHA256_MISMATCH")
    source = receipt.get("source_object")
    if not isinstance(source, dict):
        raise IntakeRejected("INVALID_RECEIPT_SOURCE_OBJECT")
    result = {"tenant_id": tenant_id, "raw_sha256": raw_sha}
    for field in ("bucket", "key", "version_id", "sequencer"):
        result[field] = _required_string(source, field, SOURCE_TOKEN_RE)
    retained = receipt.get("retained_original")
    if not isinstance(retained, dict):
        raise IntakeRejected("INVALID_RECEIPT_RETAINED_ORIGINAL")
    result["retained_uri"] = _required_string(retained, "uri", SOURCE_TOKEN_RE)
    result["retained_version"] = _required_string(retained, "version_id", SOURCE_TOKEN_RE)
    return result


def _strict_base64(part: Message) -> bytes:
    if (part.get("Content-Transfer-Encoding") or "").strip().lower() != "base64":
        raise IntakeRejected("ATTACHMENT_TRANSFER_ENCODING_NOT_BASE64")
    payload = part.get_payload(decode=False)
    if not isinstance(payload, str):
        raise IntakeRejected("ATTACHMENT_PAYLOAD_NOT_TEXT")
    compact = BASE64_WS_RE.sub(b"", payload.encode("ascii", "strict"))
    try:
        return base64.b64decode(compact, validate=True)
    except Exception as exc:
        raise IntakeRejected("ATTACHMENT_BASE64_INVALID") from exc


def _walk(message: Message, limits: Limits) -> list[tuple[int, Message]]:
    rows: list[tuple[int, Message]] = []
    stack: list[tuple[int, Message]] = [(0, message)]
    while stack:
        depth, part = stack.pop()
        if depth > limits.max_depth:
            raise IntakeRejected("MIME_DEPTH_LIMIT")
        rows.append((depth, part))
        if len(rows) > limits.max_parts:
            raise IntakeRejected("MIME_PART_LIMIT")
        if part.defects:
            raise IntakeRejected("MIME_DEFECT")
        if len(part.items()) > limits.max_headers_per_part:
            raise IntakeRejected("MIME_HEADER_COUNT_LIMIT")
        if any(len(str(k)) + len(str(v)) > limits.max_header_chars for k, v in part.items()):
            raise IntakeRejected("MIME_HEADER_SIZE_LIMIT")
        if part.is_multipart():
            children = part.get_payload()
            if not isinstance(children, list):
                raise IntakeRejected("MIME_MULTIPART_INVALID")
            stack.extend((depth + 1, child) for child in reversed(children))
    return rows


def extract(raw: bytes, receipt: dict[str, Any], output_dir: Path, allowed_types: set[str], limits: Limits = Limits()) -> dict[str, Any]:
    if not raw or len(raw) > limits.max_raw_bytes:
        raise IntakeRejected("RAW_SIZE_LIMIT")
    authority = validate_receipt(receipt, raw)
    if output_dir.exists():
        raise IntakeRejected("OUTPUT_ALREADY_EXISTS")
    message = BytesParser(policy=policy.default).parsebytes(raw)
    parts = _walk(message, limits)
    attachments: list[dict[str, Any]] = []
    payloads: list[bytes] = []
    total = 0
    for _, part in parts:
        if part.get_content_disposition() != "attachment":
            continue
        if part.is_multipart() or part.get_content_maintype() == "message":
            raise IntakeRejected("NESTED_MESSAGE_ATTACHMENT_REJECTED")
        if len(attachments) >= limits.max_attachments:
            raise IntakeRejected("ATTACHMENT_COUNT_LIMIT")
        content_type = part.get_content_type().lower()
        if content_type not in allowed_types:
            raise IntakeRejected("ATTACHMENT_CONTENT_TYPE_REJECTED")
        data = _strict_base64(part)
        if not data or len(data) > limits.max_attachment_bytes:
            raise IntakeRejected("ATTACHMENT_SIZE_LIMIT")
        total += len(data)
        if total > limits.max_total_attachment_bytes:
            raise IntakeRejected("ATTACHMENT_TOTAL_SIZE_LIMIT")
        digest = sha256(data)
        filename = part.get_filename() or ""
        index = len(attachments) + 1
        stored_name = f"part-{index:04d}-{digest}.bin"
        attachments.append({
            "part_index": index,
            "content_type_claimed": content_type,
            "filename_sha256": sha256(filename.encode("utf-8")),
            "bytes": len(data),
            "sha256": digest,
            "quarantine_name": stored_name,
            "security_decision": "PENDING",
        })
        payloads.append(data)
    if not attachments:
        raise IntakeRejected("NO_ATTACHMENTS")
    manifest = {
        "manifest_version": 1,
        "state": "QUARANTINE_EXTRACTED",
        "automatic_storage_authorized": False,
        "tenant_id": authority["tenant_id"],
        "source": {k: authority[k] for k in ("bucket", "key", "version_id", "sequencer", "raw_sha256")},
        "retained_original": {"uri": authority["retained_uri"], "version_id": authority["retained_version"]},
        "attachments": attachments,
    }
    output_dir.parent.mkdir(parents=True, exist_ok=True)
    temp_dir = Path(tempfile.mkdtemp(prefix=f".{output_dir.name}.", dir=output_dir.parent))
    try:
        for meta, data in zip(attachments, payloads, strict=True):
            target = temp_dir / meta["quarantine_name"]
            with target.open("xb") as handle:
                handle.write(data)
                handle.flush()
                os.fsync(handle.fileno())
        manifest_bytes = canonical_json(manifest)
        with (temp_dir / "manifest.json").open("xb") as handle:
            handle.write(manifest_bytes)
            handle.flush()
            os.fsync(handle.fileno())
        os.replace(temp_dir, output_dir)
        return manifest
    except Exception:
        shutil.rmtree(temp_dir, ignore_errors=True)
        raise


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--input-eml", type=Path, required=True)
    parser.add_argument("--source-receipt", type=Path, required=True)
    parser.add_argument("--profile", type=Path, required=True)
    parser.add_argument("--output-dir", type=Path, required=True)
    args = parser.parse_args()
    profile = json.loads(args.profile.read_text(encoding="utf-8"))
    allowed = profile.get("allowed_content_types")
    if not isinstance(allowed, list) or not allowed or not all(isinstance(x, str) for x in allowed):
        raise IntakeRejected("INVALID_PROFILE_ALLOWED_TYPES")
    receipt = json.loads(args.source_receipt.read_text(encoding="utf-8"))
    manifest = extract(args.input_eml.read_bytes(), receipt, args.output_dir, set(allowed))
    print(json.dumps({"status": manifest["state"], "attachments": len(manifest["attachments"])}, sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
