# Secure Email MIME Quarantine Core

## 1. Metadata

```yaml
pack_id: "SECURE-EMAIL-MIME-QUARANTINE-CORE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un extractor MIME fail-closed que exige raw MIME retenido y receipt ligado a versión/hash, aplica límites, no usa filenames como paths, escribe adjuntos hash-named y manifest atómico con persistencia automática bloqueada."
stacks: ["CPython 3.12+ stdlib", "PowerShell 7 verification"]
compatible_with: ["AWS-POWERTOOLS-IDEMPOTENT-SQS-BATCH-COMPONENT 0.1.x", "AWS-LAMBDA-DURABLE-EXECUTION-COMPONENT 0.1.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x"]
incompatible_with: ["raw MIME no retenido", "receipts sin version/hash", "persistencia automática", "filenames como object keys", "adjuntos no-base64"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://github.com/aws-samples/sample-amazon-ses-mail-manager-attachment-pipeline/tree/79314a93bda431bd03a4c1236bc4b22e75e2f41f", "https://github.com/aws-samples/serverless-mail/tree/70ac9310a313cc95f1092083810659961afb88ba", "https://github.com/aws-powertools/powertools-lambda-python"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use este core después de que Mail Manager/SES u otro provider haya retenido el MIME exacto y emitido un receipt ADMIT con tenant, bucket, key, version ID, sequencer, SHA-256 y retained-original. No recibe SMTP, no llama AWS y no autoriza persistencia: entrega bytes en cuarentena para el gate Magika/ClamAV/YARA.

## 3. Architecture contract

El código es glue local `AUTHORED`, no código AWS. Corrige de manera demostrable las brechas observadas en samples oficiales: filename nunca es path, cada payload se liga a hash/part index, el raw se liga a receipt/version, existen límites de raw/parts/depth/headers/count/bytes, sólo base64 estricto, escritura atómica y `automatic_storage_authorized=false`. Debe componerse con Powertools idempotency/partial batch, Durable Execution, almacenamiento retenido y seguridad ya materializados por packs separados.

## 4. Exact file manifest

```text
CREATE secure_email_mime_quarantine_core/README.md
CREATE secure_email_mime_quarantine_core/extract_email_attachments.py
CREATE secure_email_mime_quarantine_core/test_extract_email_attachments.py
CREATE secure_email_mime_quarantine_core/provider-profile.template.json
CREATE secure_email_mime_quarantine_core/verify_contract.ps1
```

## 5. Materialization blocks

### FILE: `secure_email_mime_quarantine_core/README.md`

```yaml
block_id: "SECURE-EMAIL-MIME:01:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://secure-email-mime-quarantine-core/README.md"
license: "LicenseRef-Workspace-Owner"
sha256: "6bb9e77c142770f45a7eaa3f8c9ec552d278dfe4ba5c2d72b27febb1d327786e"
variables: []
secrets_allowed: false
```
````text
# Secure email MIME quarantine core

This component begins only after a provider has retained the exact raw MIME object and produced a version-bound admission receipt. It parses attachments into a new quarantine directory, names every byte payload from its part index and SHA-256 rather than the untrusted filename, and writes a deterministic manifest with `automatic_storage_authorized=false`.

The implementation is local `AUTHORED` integration code. It is not presented as AWS source. Its contract is derived from failures observed in official AWS samples and is intended to compose with the separately materialized AWS Powertools idempotent partial-batch component, Lambda Durable Execution component, retained S3 Object Lock lane, and Magika/ClamAV/YARA security gate.

Only base64 attachments and explicitly allowlisted claimed MIME types are accepted. This is deliberately fail-closed. A later security gate must identify the real content from bytes before any extraction or business persistence.

Run `pwsh ./verify_contract.ps1`. Live AWS remains blocked until the project proves Mail Manager/SES, DNS/MX, tenant routing, KMS, retained raw MIME, SQS/DLQ/redrive, DynamoDB idempotency, Lambda packaging, security engines, costs and outage/replay reconciliation.
````

### FILE: `secure_email_mime_quarantine_core/extract_email_attachments.py`

```yaml
block_id: "SECURE-EMAIL-MIME:02:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://secure-email-mime-quarantine-core/extract_email_attachments.py"
license: "LicenseRef-Workspace-Owner"
sha256: "38d428ee13a4c179a361753f075fe59850e99606e39eed715c6f64fdb4652835"
variables: []
secrets_allowed: false
```
````text
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
````

### FILE: `secure_email_mime_quarantine_core/test_extract_email_attachments.py`

```yaml
block_id: "SECURE-EMAIL-MIME:03:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://secure-email-mime-quarantine-core/test_extract_email_attachments.py"
license: "LicenseRef-Workspace-Owner"
sha256: "0bd11c0b83596e6a47325eba76ae00e9cac5c7f2667e68b1f0d44f30b445f3a5"
variables: []
secrets_allowed: false
```
````text
import base64
import json
import tempfile
import unittest
from email.message import EmailMessage
from pathlib import Path

from extract_email_attachments import IntakeRejected, Limits, extract, sha256


ALLOWED = {"application/pdf", "text/csv", "image/jpeg"}


def email_bytes(items=(("invoice.pdf", "application", "pdf", b"PDF"),)):
    msg = EmailMessage()
    msg["From"] = "sender@example.com"
    msg["To"] = "intake@example.com"
    msg["Subject"] = "documents"
    msg.set_content("attached")
    for name, main, sub, data in items:
        msg.add_attachment(data, maintype=main, subtype=sub, filename=name, cte="base64")
    return msg.as_bytes()


def receipt(raw):
    return {
        "receipt_version": 1,
        "decision": "ADMIT",
        "tenant_id": "tenant-a",
        "raw_sha256": sha256(raw),
        "source_object": {"bucket": "raw-bucket", "key": "mail/1", "version_id": "v1", "sequencer": "001"},
        "retained_original": {"uri": "s3://retained/mail/1", "version_id": "locked-v1"},
    }


class ExtractTests(unittest.TestCase):
    def run_extract(self, raw=None, rec=None, limits=Limits(), allowed=ALLOWED):
        raw = raw or email_bytes()
        rec = rec or receipt(raw)
        root = tempfile.TemporaryDirectory()
        self.addCleanup(root.cleanup)
        out = Path(root.name) / "out"
        result = extract(raw, rec, out, allowed, limits)
        return result, out

    def rejected(self, code, raw=None, rec=None, limits=Limits(), allowed=ALLOWED):
        with self.assertRaises(IntakeRejected) as ctx:
            self.run_extract(raw, rec, limits, allowed)
        self.assertEqual(code, ctx.exception.code)

    def test_positive_is_hash_named_atomic_and_blocks_storage(self):
        result, out = self.run_extract()
        item = result["attachments"][0]
        self.assertEqual(b"PDF", (out / item["quarantine_name"]).read_bytes())
        self.assertFalse(result["automatic_storage_authorized"])
        self.assertEqual("PENDING", item["security_decision"])
        self.assertTrue((out / "manifest.json").is_file())

    def test_filename_traversal_is_never_a_path(self):
        raw = email_bytes((("../../secret.pdf", "application", "pdf", b"PDF"),))
        result, out = self.run_extract(raw, receipt(raw))
        self.assertNotIn("secret", result["attachments"][0]["quarantine_name"])
        self.assertEqual(2, len(list(out.iterdir())))

    def test_multiple_same_names_do_not_collide(self):
        raw = email_bytes((("same.pdf", "application", "pdf", b"A"), ("same.pdf", "application", "pdf", b"B")))
        result, _ = self.run_extract(raw, receipt(raw))
        self.assertEqual(2, len({x["quarantine_name"] for x in result["attachments"]}))

    def test_manifest_is_deterministic(self):
        raw = email_bytes()
        a, _ = self.run_extract(raw, receipt(raw))
        b, _ = self.run_extract(raw, receipt(raw))
        self.assertEqual(a, b)

    def test_receipt_must_admit(self):
        raw = email_bytes(); rec = receipt(raw); rec["decision"] = "REJECT"
        self.rejected("RECEIPT_NOT_ADMITTED", raw, rec)

    def test_raw_hash_must_match(self):
        raw = email_bytes(); rec = receipt(raw); rec["raw_sha256"] = "0" * 64
        self.rejected("RAW_SHA256_MISMATCH", raw, rec)

    def test_source_version_is_required(self):
        raw = email_bytes(); rec = receipt(raw); rec["source_object"]["version_id"] = ""
        self.rejected("INVALID_RECEIPT_VERSION_ID", raw, rec)

    def test_retained_original_is_required(self):
        raw = email_bytes(); rec = receipt(raw); rec.pop("retained_original")
        self.rejected("INVALID_RECEIPT_RETAINED_ORIGINAL", raw, rec)

    def test_raw_size_limit(self):
        raw = email_bytes()
        self.rejected("RAW_SIZE_LIMIT", raw, receipt(raw), Limits(max_raw_bytes=10))

    def test_attachment_count_limit(self):
        raw = email_bytes((("a.pdf", "application", "pdf", b"A"), ("b.pdf", "application", "pdf", b"B")))
        self.rejected("ATTACHMENT_COUNT_LIMIT", raw, receipt(raw), Limits(max_attachments=1))

    def test_attachment_size_limit(self):
        raw = email_bytes((("a.pdf", "application", "pdf", b"AB"),))
        self.rejected("ATTACHMENT_SIZE_LIMIT", raw, receipt(raw), Limits(max_attachment_bytes=1))

    def test_total_size_limit(self):
        raw = email_bytes((("a.pdf", "application", "pdf", b"AA"), ("b.pdf", "application", "pdf", b"BB")))
        self.rejected("ATTACHMENT_TOTAL_SIZE_LIMIT", raw, receipt(raw), Limits(max_total_attachment_bytes=3))

    def test_content_type_allowlist(self):
        raw = email_bytes((("a.exe", "application", "octet-stream", b"MZ"),))
        self.rejected("ATTACHMENT_CONTENT_TYPE_REJECTED", raw, receipt(raw))

    def test_non_base64_rejected(self):
        msg = EmailMessage(); msg.set_content("body"); msg.add_attachment("plain", subtype="csv", filename="a.csv", cte="quoted-printable")
        raw = msg.as_bytes()
        self.rejected("ATTACHMENT_TRANSFER_ENCODING_NOT_BASE64", raw, receipt(raw))

    def test_no_attachments_rejected(self):
        msg = EmailMessage(); msg.set_content("body"); raw = msg.as_bytes()
        self.rejected("NO_ATTACHMENTS", raw, receipt(raw))

    def test_existing_output_rejected(self):
        raw = email_bytes(); rec = receipt(raw)
        with tempfile.TemporaryDirectory() as td:
            out = Path(td) / "out"; out.mkdir()
            with self.assertRaises(IntakeRejected) as ctx: extract(raw, rec, out, ALLOWED)
            self.assertEqual("OUTPUT_ALREADY_EXISTS", ctx.exception.code)

    def test_invalid_tenant_rejected(self):
        raw = email_bytes(); rec = receipt(raw); rec["tenant_id"] = "../tenant"
        self.rejected("INVALID_RECEIPT_TENANT_ID", raw, rec)

    def test_part_limit(self):
        raw = email_bytes((("a.pdf", "application", "pdf", b"A"),))
        self.rejected("MIME_PART_LIMIT", raw, receipt(raw), Limits(max_parts=1))


if __name__ == "__main__":
    unittest.main(verbosity=2)
````

### FILE: `secure_email_mime_quarantine_core/provider-profile.template.json`

```yaml
block_id: "SECURE-EMAIL-MIME:04:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://secure-email-mime-quarantine-core/provider-profile.template.json"
license: "LicenseRef-Workspace-Owner"
sha256: "ba64a662170f8bba249cf8e194a2d20eb20001c2e7aed6ea5f4024879baac34b"
variables: []
secrets_allowed: false
```
````text
{"allowed_content_types":[],"automatic_storage_authorized":false,"profile_version":1}
````

### FILE: `secure_email_mime_quarantine_core/verify_contract.ps1`

```yaml
block_id: "SECURE-EMAIL-MIME:05:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://secure-email-mime-quarantine-core/verify_contract.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "be019f9acddb8e165582b7010b9ae1e5659176456f49daf0cd86801e58880072"
variables: []
secrets_allowed: false
```
````text
#requires -Version 7.0
[CmdletBinding()]
param([string]$PythonExecutable = 'python')
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest
$source = Get-Content -Raw -LiteralPath (Join-Path $PSScriptRoot 'extract_email_attachments.py')
foreach ($token in @('automatic_storage_authorized', 'RAW_SHA256_MISMATCH', 'MIME_PART_LIMIT', 'ATTACHMENT_BASE64_INVALID', 'os.replace', 'security_decision')) {
  if (-not $source.Contains($token, [StringComparison]::Ordinal)) { throw "Missing contract token: $token" }
}
& $PythonExecutable -B -m py_compile (Join-Path $PSScriptRoot 'extract_email_attachments.py') (Join-Path $PSScriptRoot 'test_extract_email_attachments.py')
if ($LASTEXITCODE -ne 0) { throw 'Python compile failed' }
Push-Location $PSScriptRoot
try { & $PythonExecutable -B -m unittest -v test_extract_email_attachments.py; if ($LASTEXITCODE -ne 0) { throw 'Tests failed' } }
finally { Pop-Location }
Write-Output 'SECURE_EMAIL_MIME_QUARANTINE_CORE_PASS tests=18 storage_authorized=0 filename_as_path=0 receipt_bound=1 atomic=1'
````

## 6. Configuration surface

El profile empieza con allowlist vacía. El proyecto debe fijar MIME permitidos y conservar `automatic_storage_authorized=false`. Los límites tienen defaults cerrados y sólo pueden ampliarse mediante decisión/evidencia del target.

## 7. Dependency bill

CPython 3.12+ stdlib y PowerShell 7. No descarga paquetes ni usa red. La integración AWS posterior usa únicamente packs oficiales separados.

## 8. Apply order

Readiness y blueprint; cuenta/costo/provider; raw MIME retenido; receipt version/hash; este core; gate de seguridad sobre cada payload; routing/extracción/evaluación; revisión y persistencia sólo con autoridad.

## 9. Verification

Materializar en carpeta vacía y ejecutar `pwsh ./secure_email_mime_quarantine_core/verify_contract.ps1`. Debe compilar y pasar 18 pruebas. Luego mutar receipt/hash/filename/encoding/límites y confirmar fallo cerrado.

## 10. Reconstruction evidence

2026-08-28: suite CPython ejecutó 18/18; verifier PASS. Un fixture negativo fue corregido porque MIME no permitido interceptaba antes que encoding; el fallo y su regresión deben registrarse. Live SMTP/AWS, malware y storage continúan condicionados.

