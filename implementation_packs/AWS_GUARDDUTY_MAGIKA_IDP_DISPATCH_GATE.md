# AWS GuardDuty Magika IDP Dispatch Gate

## 1. Metadata

```yaml
pack_id: "AWS-GUARDDUTY-MAGIKA-IDP-DISPATCH-GATE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un límite exacto GuardDuty limpio → identificación Google Magika oficial → copia versionada e idempotente al input del AWS GenAI IDP Accelerator v0.6.5; registra recibo inmutable y nunca autoriza persistencia empresarial automática."
stacks: ["AWS GenAI IDP Accelerator v0.6.5", "Google Magika CLI 1.1.0", "AWS Lambda Python 3.14", "Amazon S3 Object Lock", "Amazon EventBridge", "Amazon SQS", "Amazon DynamoDB", "AWS KMS", "AWS SAM", "Boto3 1.43.83"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "AWS-GUARDDUTY-IMMUTABLE-RELEASE-GATE 0.1.x"]
incompatible_with: ["recibo GuardDuty sin versión/eTag", "S3 sin versionado", "Magika no fijado", "arquitectura distinta de x86_64", "configuración IDP móvil", "persistencia empresarial automática", "deploy sin pruebas o costo aprobado"]
license_expression: "LicenseRef-Workspace-Owner AND MIT-0 AND Apache-2.0"
upstream_sources: ["https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/tree/1b5fd74454e593de233342a02ee911af8ee38359", "https://github.com/google/magika/tree/5e2f437fb7b7452368c8c1fa9354858f5487a5c4", "https://docs.aws.amazon.com/AmazonS3/latest/API/API_CopyObject.html", "https://docs.aws.amazon.com/AmazonS3/latest/userguide/ev-events.html"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use sólo en AWS, después de un recibo `elite-aws-guardduty-release-receipt/v1` limpio y antes del input del AWS GenAI IDP Accelerator v0.6.5. Requiere fuente oficial adquirida, binario Linux x86-64 de Magika 1.1.0 fijado por SHA-256 y attestation, buckets versionados SSE-KMS, configuración IDP evaluada e inmutable, corpus representativo y owners de seguridad/documentos/costo/replay. Rechazar versiones móviles, tipos Magika no admitidos, evidencia incompleta o cualquier pretensión de persistir datos empresariales automáticamente.

## 3. Architecture contract

`handler.py` es glue `ADAPTED`, no una copia atribuida a AWS o Google. Implementa los contratos públicos de S3/EventBridge/SQS, el límite de entrada del acelerador oficial AWS y el ejecutable oficial Google Magika. Valida cuenta, región, recurso, versión mayor del evento, tenant, VersionId, eTag, tamaño, KMS, esquema/decisión de seguridad y metadatos exactos; descarga sólo esa versión para identificación y copia esa misma versión con `CopySourceIfMatch`, `IfNoneMatch: *`, checksum, clave content-addressed y `config-version` fijado.

DynamoDB aporta lease con owner fencing; el recibo de dispatch Object-Locked se confirma último y SQS usa respuesta parcial. El worker sólo registra hashes y tipos de error. El input oficial AWS continúa siendo at-least-once (`UP-FAIL-177`), por lo que el recibo fija `automatic_business_persistence_authorized=false`: evaluación de campos, revisión humana cuando corresponda e idempotencia empresarial permanecen como gates posteriores obligatorios.

## 4. Exact file manifest

```text
CREATE aws_guardduty_magika_idp_dispatch/README.md
CREATE aws_guardduty_magika_idp_dispatch/handler.py
CREATE aws_guardduty_magika_idp_dispatch/test_handler.py
CREATE aws_guardduty_magika_idp_dispatch/template.yaml
CREATE aws_guardduty_magika_idp_dispatch/provider-profile.template.json
CREATE aws_guardduty_magika_idp_dispatch/requirements-runtime.lock
CREATE aws_guardduty_magika_idp_dispatch/build.ps1
CREATE aws_guardduty_magika_idp_dispatch/deploy.ps1
CREATE aws_guardduty_magika_idp_dispatch/verify_contract.ps1
CREATE aws_guardduty_magika_idp_dispatch/GOOGLE_APACHE_LICENSE.txt
CREATE aws_guardduty_magika_idp_dispatch/AWS_MIT_NO_ATTRIBUTION_LICENSE.txt
```

## 5. Materialization blocks

### FILE: `aws_guardduty_magika_idp_dispatch/README.md`
```yaml
block_id: "AWS-GUARDDUTY-MAGIKA-IDP-DISPATCH-GATE:readme-md:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-guardduty-magika-idp-dispatch-gate/README.md"
license: "LicenseRef-Workspace-Owner"
sha256: "79fda91314ed08248a459dbac7489077cf5c3923f494ccc02cd756b5e184e92a"
variables: []
secrets_allowed: false
```
````markdown
# AWS GuardDuty release to Magika and GenAI IDP dispatch

This component connects an exact clean `elite-aws-guardduty-release-receipt/v1` to the official AWS GenAI IDP Accelerator v0.6.5 input boundary. It does not reimplement OCR, classification or extraction. It verifies the immutable security receipt and released S3 version, downloads those exact bytes only for Google Magika CLI v1.1.0 identification, then conditionally copies the same exact version to a content-addressed AWS IDP input key with a pinned `config-version`. A retained dispatch receipt is committed last.

`handler.py` is `ADAPTED` integration glue, not verbatim AWS or Google source. Its authorities are AWS S3/EventBridge/SQS conditional and version APIs, the official AWS IDP v0.6.5 input metadata contract, and the exact Google Magika Linux x86-64 release asset. The AWS IDP source remains acquired separately through `OFFICIAL-UPSTREAM-ACQUISITION-CORE`; this pack never relabels the full accelerator as local code.

The worker rejects wrong account, region, bucket, tenant, event major version, receipt version/eTag/KMS/schema, non-clean decision, changed release version/hash/size/KMS/metadata, unapproved Magika label/MIME/score/extension, unversioned or divergent IDP input and divergent receipts. A DynamoDB owner-fenced lease and deterministic S3 keys reconcile duplicate delivery and a copy completed before the receipt. SQS partial batch responses isolate failures. Logs contain only hashes and error types.

The official accelerator remains at-least-once at its S3/EventBridge intake and v0.6.5 has the specific limitations recorded as `UP-FAIL-177`. This gate prevents duplicate copies from its own replay but cannot truthfully promise that AWS will never redeliver the single Object Created event. Therefore its receipt fixes `automatic_business_persistence_authorized=false`; extracted output must still pass strict field evaluation, human review rules and an idempotent business persistence boundary.

Build requires the exact Google asset `magika-cli-x86_64-unknown-linux-gnu.tar.xz`, SHA-256 `6b4c1010...00485`, whose release asset, sidecar, GitHub API digest and GitHub attestation were verified. The build copies the official binary into the Lambda artifact and installs only the seven hash-locked Boto3 1.43.83 distributions. Live deployment requires an existing verified GuardDuty release lane, a separately acquired/deployed AWS IDP v0.6.5 stack, a named evaluated IDP configuration, all owners and explicit cost/effect approvals.
````

### FILE: `aws_guardduty_magika_idp_dispatch/handler.py`
```yaml
block_id: "AWS-GUARDDUTY-MAGIKA-IDP-DISPATCH-GATE:handler-py:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/tree/1b5fd74454e593de233342a02ee911af8ee38359; https://github.com/google/magika/tree/5e2f437fb7b7452368c8c1fa9354858f5487a5c4; integration glue hardened against official AWS S3/EventBridge/SQS and Google Magika contracts"
license: "LicenseRef-Workspace-Owner AND MIT-0 AND Apache-2.0"
sha256: "968225f37f07aeb7bb806927373d11a737b5ebcbd7605f48a831177cbace6220"
variables: []
secrets_allowed: false
```
````python
"""Version-bound GuardDuty release to AWS GenAI IDP input dispatch.

ADAPTED integration glue.  Authorities are AWS S3/EventBridge/SQS conditional
copy contracts, AWS GenAI IDP Accelerator v0.6.5 input metadata, and Google
Magika CLI v1.1.0 output.  This file is not verbatim AWS or Google source.
"""

from __future__ import annotations

import base64
import hashlib
import json
import os
import re
import shutil
import subprocess
import tempfile
import time
import uuid
from dataclasses import dataclass
from datetime import datetime, timedelta, timezone
from pathlib import Path
from typing import Any, Callable

import boto3
from botocore.exceptions import ClientError


RECEIPT_SCHEMA = "elite-aws-guardduty-release-receipt/v1"
DISPATCH_SCHEMA = "elite-aws-magika-idp-dispatch-receipt/v1"
MAGIKA_VERSION = "magika 1.1.0 standard_v3_3"
HEX64 = re.compile(r"^[0-9a-f]{64}$")
SAFE_VERSION = re.compile(r"^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$")
SAFE_EXTENSION = re.compile(r"^[a-z0-9]{1,10}$")


class Rejected(ValueError):
    """Permanent contract rejection."""


class Retryable(RuntimeError):
    """Transient or ambiguous effect; SQS must retry."""


@dataclass(frozen=True)
class Config:
    account_id: str
    region: str
    tenant_id: str
    receipt_bucket: str
    receipt_kms_key_arn: str
    release_bucket: str
    release_kms_key_arn: str
    idp_input_bucket: str
    idp_input_kms_key_arn: str
    dispatch_receipt_bucket: str
    dispatch_receipt_kms_key_arn: str
    idempotency_table: str
    config_version: str
    magika_binary: str
    magika_sha256: str
    allowed_types: tuple[dict[str, Any], ...]
    max_bytes: int
    retain_days: int
    lease_seconds: int

    @staticmethod
    def from_env() -> "Config":
        try:
            allowed = json.loads(os.environ.get("ALLOWED_MAGIKA_TYPES_JSON", "[]"))
        except json.JSONDecodeError as error:
            raise Rejected("ALLOWED_MAGIKA_TYPES_JSON is invalid") from error
        config = Config(
            account_id=os.environ.get("EXPECTED_ACCOUNT_ID", ""),
            region=os.environ.get("EXPECTED_REGION", ""),
            tenant_id=os.environ.get("TENANT_ID", ""),
            receipt_bucket=os.environ.get("SECURITY_RECEIPT_BUCKET", ""),
            receipt_kms_key_arn=os.environ.get("SECURITY_RECEIPT_KMS_KEY_ARN", ""),
            release_bucket=os.environ.get("RELEASE_BUCKET", ""),
            release_kms_key_arn=os.environ.get("RELEASE_KMS_KEY_ARN", ""),
            idp_input_bucket=os.environ.get("IDP_INPUT_BUCKET", ""),
            idp_input_kms_key_arn=os.environ.get("IDP_INPUT_KMS_KEY_ARN", ""),
            dispatch_receipt_bucket=os.environ.get("DISPATCH_RECEIPT_BUCKET", ""),
            dispatch_receipt_kms_key_arn=os.environ.get("DISPATCH_RECEIPT_KMS_KEY_ARN", ""),
            idempotency_table=os.environ.get("IDEMPOTENCY_TABLE", ""),
            config_version=os.environ.get("IDP_CONFIG_VERSION", ""),
            magika_binary=os.environ.get("MAGIKA_BINARY", "/var/task/bin/magika"),
            magika_sha256=os.environ.get("MAGIKA_SHA256", ""),
            allowed_types=tuple(allowed) if isinstance(allowed, list) else (),
            max_bytes=int(os.environ.get("MAX_DOCUMENT_BYTES", "20971520")),
            retain_days=int(os.environ.get("RETAIN_DAYS", "365")),
            lease_seconds=int(os.environ.get("LEASE_SECONDS", "900")),
        )
        config.validate()
        return config

    def validate(self) -> None:
        required = {
            "account_id": self.account_id, "region": self.region, "tenant_id": self.tenant_id,
            "receipt_bucket": self.receipt_bucket, "receipt_kms_key_arn": self.receipt_kms_key_arn,
            "release_bucket": self.release_bucket, "release_kms_key_arn": self.release_kms_key_arn,
            "idp_input_bucket": self.idp_input_bucket, "idp_input_kms_key_arn": self.idp_input_kms_key_arn,
            "dispatch_receipt_bucket": self.dispatch_receipt_bucket,
            "dispatch_receipt_kms_key_arn": self.dispatch_receipt_kms_key_arn,
            "idempotency_table": self.idempotency_table, "config_version": self.config_version,
            "magika_binary": self.magika_binary, "magika_sha256": self.magika_sha256,
        }
        if any(not value for value in required.values()):
            raise Rejected("required configuration is missing")
        if not re.fullmatch(r"[0-9]{12}", self.account_id):
            raise Rejected("EXPECTED_ACCOUNT_ID is invalid")
        if not re.fullmatch(r"[a-z0-9][a-z0-9-]{1,62}", self.tenant_id):
            raise Rejected("TENANT_ID is invalid")
        if not SAFE_VERSION.fullmatch(self.config_version):
            raise Rejected("IDP_CONFIG_VERSION is invalid")
        if not HEX64.fullmatch(self.magika_sha256):
            raise Rejected("MAGIKA_SHA256 is invalid")
        if not (1 <= self.max_bytes <= 20 * 1024 * 1024 and self.retain_days >= 1 and 30 <= self.lease_seconds <= 3600):
            raise Rejected("limits are invalid")
        if not self.allowed_types:
            raise Rejected("allowed Magika types must be non-empty")
        seen: set[tuple[str, str]] = set()
        for item in self.allowed_types:
            if not isinstance(item, dict) or set(item) != {"label", "mime_type", "extension", "min_score"}:
                raise Rejected("allowed Magika type has invalid fields")
            label, mime, extension = item["label"], item["mime_type"], item["extension"]
            score = item["min_score"]
            if not isinstance(label, str) or not re.fullmatch(r"[a-z0-9][a-z0-9_+-]{0,63}", label):
                raise Rejected("allowed Magika label is invalid")
            if not isinstance(mime, str) or not re.fullmatch(r"[a-z0-9.+-]+/[a-z0-9.+-]+", mime):
                raise Rejected("allowed MIME is invalid")
            if not isinstance(extension, str) or not SAFE_EXTENSION.fullmatch(extension):
                raise Rejected("allowed extension is invalid")
            if not isinstance(score, (int, float)) or isinstance(score, bool) or not 0.5 <= float(score) <= 1.0:
                raise Rejected("allowed Magika score is invalid")
            identity = (label, mime)
            if identity in seen:
                raise Rejected("allowed Magika identity is duplicated")
            seen.add(identity)


@dataclass(frozen=True)
class Clients:
    s3: Any
    dynamodb: Any


def _clients() -> Clients:
    return Clients(boto3.client("s3"), boto3.client("dynamodb"))


def _sha256(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def _strict_json(data: bytes, label: str) -> dict[str, Any]:
    def pairs(items: list[tuple[str, Any]]) -> dict[str, Any]:
        result: dict[str, Any] = {}
        for key, value in items:
            if key in result:
                raise Rejected(f"{label} contains duplicate JSON key")
            result[key] = value
        return result
    try:
        value = json.loads(data.decode("utf-8"), object_pairs_hook=pairs)
    except (UnicodeDecodeError, json.JSONDecodeError) as error:
        raise Rejected(f"{label} is invalid JSON") from error
    if not isinstance(value, dict):
        raise Rejected(f"{label} must be a JSON object")
    return value


def _canonical_json(value: Any) -> bytes:
    return (json.dumps(value, sort_keys=True, separators=(",", ":")) + "\n").encode("utf-8")


def _parse_event(record: dict[str, Any], config: Config) -> dict[str, str | int]:
    try:
        envelope = _strict_json(record["body"].encode("utf-8"), "SQS body")
        detail = envelope["detail"]
        obj = detail["object"]
        bucket = detail["bucket"]["name"]
    except (KeyError, AttributeError) as error:
        raise Rejected("event fields are missing") from error
    if envelope.get("version") != "0" or envelope.get("source") != "aws.s3" or envelope.get("detail-type") != "Object Created":
        raise Rejected("unexpected EventBridge envelope")
    if envelope.get("account") != config.account_id or envelope.get("region") != config.region:
        raise Rejected("event account or region mismatch")
    if envelope.get("resources") != [f"arn:aws:s3:::{config.receipt_bucket}"]:
        raise Rejected("event resource mismatch")
    event_version = detail.get("event-version")
    if not isinstance(event_version, str) or event_version.split(".", 1)[0] != "1":
        raise Rejected("unsupported S3 event major version")
    if bucket != config.receipt_bucket:
        raise Rejected("receipt bucket mismatch")
    key, version_id, etag, size = obj.get("key"), obj.get("version-id"), obj.get("etag"), obj.get("size")
    prefix = f"security-receipts/{config.tenant_id}/"
    if not isinstance(key, str) or not key.startswith(prefix) or not key.endswith(".json") or "/../" in key:
        raise Rejected("receipt key is outside tenant prefix")
    if not all(isinstance(value, str) and value for value in (version_id, etag)):
        raise Rejected("receipt version or etag is missing")
    if not isinstance(size, int) or isinstance(size, bool) or not 2 <= size <= 65536:
        raise Rejected("receipt event size is invalid")
    return {"key": key, "version_id": version_id, "etag": etag.strip('"'), "size": size}


def _read_receipt(s3: Any, event: dict[str, str | int], config: Config) -> tuple[dict[str, Any], str]:
    response = s3.get_object(
        Bucket=config.receipt_bucket, Key=event["key"], VersionId=event["version_id"],
        IfMatch=event["etag"], ExpectedBucketOwner=config.account_id,
    )
    body = response["Body"].read(65537)
    if len(body) != event["size"] or len(body) > 65536:
        raise Rejected("receipt bytes do not match event")
    if response.get("ServerSideEncryption") != "aws:kms" or response.get("SSEKMSKeyId") != config.receipt_kms_key_arn:
        raise Rejected("receipt encryption mismatch")
    receipt = _strict_json(body, "security receipt")
    required = {
        "schema", "receipt_version", "scan_id", "tenant_id", "provider", "provider_schema_version",
        "decision", "released_object", "automatic_business_persistence_authorized", "next_gate",
    }
    if not required.issubset(receipt):
        raise Rejected("security receipt contract is incomplete")
    if receipt["schema"] != RECEIPT_SCHEMA or receipt["receipt_version"] != 1:
        raise Rejected("security receipt schema mismatch")
    if receipt["tenant_id"] != config.tenant_id or receipt["provider"] != "AWS_GUARDDUTY_MALWARE_PROTECTION_S3":
        raise Rejected("security receipt authority mismatch")
    if receipt["decision"] != "ADMIT_TO_EXTRACTION" or receipt["next_gate"] != "DOCUMENT_CLASSIFICATION_AND_FIELD_EVIDENCE":
        raise Rejected("security receipt did not admit extraction")
    if receipt["automatic_business_persistence_authorized"] is not False:
        raise Rejected("security receipt must deny business persistence")
    if not isinstance(receipt["scan_id"], str) or not HEX64.fullmatch(receipt["scan_id"]):
        raise Rejected("scan identity is invalid")
    return receipt, _sha256(body)


def _verify_released(s3: Any, receipt: dict[str, Any], config: Config) -> dict[str, Any]:
    released = receipt.get("released_object")
    if not isinstance(released, dict):
        raise Rejected("released object is missing")
    required = {"bucket", "key", "version_id", "sha256", "bytes", "retention_verified"}
    if not required.issubset(released) or released["bucket"] != config.release_bucket or released["retention_verified"] is not True:
        raise Rejected("released object contract mismatch")
    sha, size = released["sha256"], released["bytes"]
    if not isinstance(sha, str) or not HEX64.fullmatch(sha) or not isinstance(size, int) or isinstance(size, bool) or not 1 <= size <= config.max_bytes:
        raise Rejected("released object hash or size is invalid")
    if released["key"] != f"clean/{config.tenant_id}/{sha}.bin" or not isinstance(released["version_id"], str) or not released["version_id"]:
        raise Rejected("released object location is invalid")
    head = s3.head_object(
        Bucket=config.release_bucket, Key=released["key"], VersionId=released["version_id"],
        ChecksumMode="ENABLED", ExpectedBucketOwner=config.account_id,
    )
    metadata = head.get("Metadata") or {}
    if head.get("ContentLength") != size or head.get("ServerSideEncryption") != "aws:kms" or head.get("SSEKMSKeyId") != config.release_kms_key_arn:
        raise Rejected("released object storage evidence mismatch")
    if metadata.get("content-sha256") != sha or metadata.get("tenant-id") != config.tenant_id:
        raise Rejected("released object metadata mismatch")
    etag = str(head.get("ETag", "")).strip('"')
    if not etag:
        raise Rejected("released object etag is missing")
    return {**released, "etag": etag, "raw_sha256": metadata.get("raw-sha256", "")}


def _download_exact(s3: Any, released: dict[str, Any], config: Config, target: Path) -> None:
    response = s3.get_object(
        Bucket=config.release_bucket, Key=released["key"], VersionId=released["version_id"],
        IfMatch=released["etag"], ExpectedBucketOwner=config.account_id,
    )
    digest = hashlib.sha256()
    total = 0
    with target.open("xb") as output:
        while True:
            chunk = response["Body"].read(1024 * 1024)
            if not chunk:
                break
            total += len(chunk)
            if total > config.max_bytes:
                raise Rejected("released object exceeds configured limit")
            digest.update(chunk)
            output.write(chunk)
    if total != released["bytes"] or digest.hexdigest() != released["sha256"]:
        raise Rejected("downloaded bytes do not match release receipt")


def _prepare_magika(config: Config) -> Path:
    source = Path(config.magika_binary)
    if not source.is_file() or source.is_symlink() or _sha256(source.read_bytes()) != config.magika_sha256:
        raise Rejected("Magika binary identity mismatch")
    target = Path(tempfile.gettempdir()) / f"magika-{config.magika_sha256[:16]}"
    if not target.exists():
        staging = target.with_name(target.name + f".{uuid.uuid4().hex}.tmp")
        shutil.copyfile(source, staging)
        os.chmod(staging, 0o700)
        try:
            os.replace(staging, target)
        except FileExistsError:
            staging.unlink(missing_ok=True)
    if _sha256(target.read_bytes()) != config.magika_sha256:
        raise Rejected("prepared Magika binary identity mismatch")
    os.chmod(target, 0o700)
    return target


def _identify(path: Path, config: Config, runner: Callable[..., Any] = subprocess.run) -> dict[str, Any]:
    binary = _prepare_magika(config)
    version = runner([str(binary), "--version"], capture_output=True, text=True, timeout=15, check=False)
    if version.returncode != 0 or version.stdout.strip() != MAGIKA_VERSION:
        raise Rejected("Magika version mismatch")
    result = runner([str(binary), "--json", "--no-colors", str(path)], capture_output=True, text=True, timeout=60, check=False)
    if result.returncode != 0:
        raise Rejected("Magika identification failed")
    try:
        items = json.loads(result.stdout)
        value = items[0]["result"]["value"]
        output = value["output"]
        label, mime, score = output["label"], output["mime_type"], float(value["score"])
        extensions = output["extensions"]
    except (KeyError, IndexError, TypeError, ValueError, json.JSONDecodeError) as error:
        raise Rejected("Magika output is invalid") from error
    for allowed in config.allowed_types:
        if allowed["label"] == label and allowed["mime_type"] == mime and score >= float(allowed["min_score"]):
            if not isinstance(extensions, list) or allowed["extension"] not in extensions:
                raise Rejected("configured extension is not authorized by Magika output")
            return {"label": label, "mime_type": mime, "score": score, "extension": allowed["extension"], "version": MAGIKA_VERSION}
    raise Rejected("content type is not admitted")


def _receipt_key(config: Config, sha: str) -> str:
    config_hash = _sha256(config.config_version.encode("utf-8"))[:16]
    return f"dispatch-receipts/{config.tenant_id}/{sha}/{config_hash}.json"


def _existing_dispatch(s3: Any, config: Config, sha: str) -> dict[str, str] | None:
    key = _receipt_key(config, sha)
    try:
        head = s3.head_object(Bucket=config.dispatch_receipt_bucket, Key=key, ExpectedBucketOwner=config.account_id)
    except ClientError as error:
        if error.response.get("Error", {}).get("Code") in {"404", "NoSuchKey", "NotFound"}:
            return None
        raise
    metadata = head.get("Metadata") or {}
    if metadata.get("content-sha256") != sha or metadata.get("config-version-sha256") != _sha256(config.config_version.encode("utf-8")):
        raise Rejected("existing dispatch receipt authority mismatch")
    version_id = head.get("VersionId")
    if not version_id:
        raise Rejected("existing dispatch receipt is not versioned")
    return {"bucket": config.dispatch_receipt_bucket, "key": key, "version_id": version_id}


def _acquire_lease(ddb: Any, config: Config, dispatch_id: str, owner: str, now: int) -> None:
    try:
        ddb.update_item(
            TableName=config.idempotency_table, Key={"DispatchId": {"S": dispatch_id}},
            UpdateExpression="SET #state=:working, lease_owner=:owner, lease_expires=:expires, updated_at=:now",
            ConditionExpression="attribute_not_exists(DispatchId) OR (#state=:working AND lease_expires < :now)",
            ExpressionAttributeNames={"#state": "state"},
            ExpressionAttributeValues={":working": {"S": "WORKING"}, ":owner": {"S": owner},
                                               ":expires": {"N": str(now + config.lease_seconds)}, ":now": {"N": str(now)}},
        )
    except ClientError as error:
        if error.response.get("Error", {}).get("Code") == "ConditionalCheckFailedException":
            raise Retryable("dispatch lease is owned by another worker") from error
        raise


def _finish_lease(ddb: Any, config: Config, dispatch_id: str, owner: str, receipt_key: str, now: int) -> None:
    ddb.update_item(
        TableName=config.idempotency_table, Key={"DispatchId": {"S": dispatch_id}},
        UpdateExpression="SET #state=:done, receipt_key=:receipt, completed_at=:now REMOVE lease_expires",
        ConditionExpression="#state=:working AND lease_owner=:owner",
        ExpressionAttributeNames={"#state": "state"},
        ExpressionAttributeValues={":done": {"S": "DONE"}, ":working": {"S": "WORKING"},
                                           ":owner": {"S": owner}, ":receipt": {"S": receipt_key}, ":now": {"N": str(now)}},
    )


def _is_precondition(error: ClientError) -> bool:
    return error.response.get("Error", {}).get("Code") in {"412", "PreconditionFailed", "ConditionalRequestConflict"}


def _copy_to_idp(s3: Any, released: dict[str, Any], detected: dict[str, Any], security_receipt_sha: str, config: Config) -> dict[str, Any]:
    key = f"security-released/{config.tenant_id}/{released['sha256']}.{detected['extension']}"
    metadata = {
        "config-version": config.config_version,
        "content-sha256": released["sha256"], "raw-sha256": released.get("raw_sha256", ""),
        "tenant-id": config.tenant_id, "security-receipt-sha256": security_receipt_sha,
        "magika-label": detected["label"], "magika-mime": detected["mime_type"],
    }
    try:
        response = s3.copy_object(
            CopySource={"Bucket": config.release_bucket, "Key": released["key"], "VersionId": released["version_id"]},
            CopySourceIfMatch=released["etag"], Bucket=config.idp_input_bucket, Key=key, IfNoneMatch="*",
            Metadata=metadata, MetadataDirective="REPLACE", ContentType=detected["mime_type"],
            ChecksumAlgorithm="SHA256", ServerSideEncryption="aws:kms", SSEKMSKeyId=config.idp_input_kms_key_arn,
            BucketKeyEnabled=True, ExpectedBucketOwner=config.account_id,
        )
        version_id = response.get("VersionId")
    except ClientError as error:
        if not _is_precondition(error):
            raise
        version_id = None
    head_args = {"Bucket": config.idp_input_bucket, "Key": key, "ChecksumMode": "ENABLED", "ExpectedBucketOwner": config.account_id}
    if version_id:
        head_args["VersionId"] = version_id
    head = s3.head_object(**head_args)
    observed = head.get("Metadata") or {}
    if head.get("ContentLength") != released["bytes"] or head.get("ServerSideEncryption") != "aws:kms" or head.get("SSEKMSKeyId") != config.idp_input_kms_key_arn:
        raise Rejected("IDP input storage evidence mismatch")
    if any(observed.get(key_name) != value for key_name, value in metadata.items()):
        raise Rejected("IDP input metadata mismatch")
    version_id = head.get("VersionId") or version_id
    if not version_id:
        raise Rejected("IDP input bucket is not versioned")
    return {"bucket": config.idp_input_bucket, "key": key, "version_id": version_id,
            "sha256": released["sha256"], "bytes": released["bytes"], "config_version": config.config_version}


def _put_dispatch_receipt(s3: Any, receipt: dict[str, Any], config: Config) -> dict[str, str]:
    sha = receipt["content_sha256"]
    key = _receipt_key(config, sha)
    data = _canonical_json(receipt)
    retain_until = datetime.now(timezone.utc) + timedelta(days=config.retain_days)
    metadata = {"content-sha256": sha, "config-version-sha256": _sha256(config.config_version.encode("utf-8")),
                "receipt-sha256": _sha256(data)}
    try:
        response = s3.put_object(
            Bucket=config.dispatch_receipt_bucket, Key=key, Body=data, ContentLength=len(data), ContentType="application/json",
            ChecksumAlgorithm="SHA256", ChecksumSHA256=base64.b64encode(hashlib.sha256(data).digest()).decode("ascii"),
            IfNoneMatch="*", Metadata=metadata, ServerSideEncryption="aws:kms", SSEKMSKeyId=config.dispatch_receipt_kms_key_arn,
            BucketKeyEnabled=True, ObjectLockMode="COMPLIANCE", ObjectLockRetainUntilDate=retain_until,
            ExpectedBucketOwner=config.account_id,
        )
        version_id = response.get("VersionId")
    except ClientError as error:
        if not _is_precondition(error):
            raise
        existing = _existing_dispatch(s3, config, sha)
        if existing is None:
            raise Retryable("dispatch receipt precondition was ambiguous") from error
        return existing
    head = s3.head_object(Bucket=config.dispatch_receipt_bucket, Key=key, VersionId=version_id,
                          ExpectedBucketOwner=config.account_id)
    if head.get("ServerSideEncryption") != "aws:kms" or head.get("SSEKMSKeyId") != config.dispatch_receipt_kms_key_arn:
        raise Rejected("dispatch receipt encryption mismatch")
    retention = s3.get_object_retention(Bucket=config.dispatch_receipt_bucket, Key=key, VersionId=version_id,
                                        ExpectedBucketOwner=config.account_id).get("Retention") or {}
    if retention.get("Mode") != "COMPLIANCE" or not retention.get("RetainUntilDate"):
        raise Rejected("dispatch receipt retention mismatch")
    return {"bucket": config.dispatch_receipt_bucket, "key": key, "version_id": version_id}


def process_record(record: dict[str, Any], clients: Clients, config: Config,
                   runner: Callable[..., Any] = subprocess.run) -> dict[str, Any]:
    event = _parse_event(record, config)
    security, security_sha = _read_receipt(clients.s3, event, config)
    released = _verify_released(clients.s3, security, config)
    existing = _existing_dispatch(clients.s3, config, released["sha256"])
    if existing:
        return {"state": "DUPLICATE", "receipt": existing, "content_sha256": released["sha256"]}
    dispatch_id = _sha256(f"{config.tenant_id}\0{released['sha256']}\0{config.config_version}".encode("utf-8"))
    owner = uuid.uuid4().hex
    now = int(time.time())
    _acquire_lease(clients.dynamodb, config, dispatch_id, owner, now)
    with tempfile.TemporaryDirectory(prefix="elite-idp-dispatch-") as directory:
        path = Path(directory) / f"{released['sha256']}.bin"
        _download_exact(clients.s3, released, config, path)
        detected = _identify(path, config, runner)
    target = _copy_to_idp(clients.s3, released, detected, security_sha, config)
    receipt = {
        "schema": DISPATCH_SCHEMA, "receipt_version": 1, "created_at": datetime.now(timezone.utc).isoformat(),
        "dispatch_id": dispatch_id, "tenant_id": config.tenant_id, "content_sha256": released["sha256"],
        "security_receipt_sha256": security_sha, "security_scan_id": security["scan_id"],
        "content_identification": detected, "idp_engine": {
            "owner": "AWS", "repository": "aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws",
            "release": "v0.6.5", "commit": "1b5fd74454e593de233342a02ee911af8ee38359",
            "config_version": config.config_version,
        },
        "idp_input_object": target, "next_gate": "AWS_GENAI_IDP_CLASSIFICATION_EXTRACTION_EVIDENCE",
        "automatic_business_persistence_authorized": False,
    }
    saved = _put_dispatch_receipt(clients.s3, receipt, config)
    _finish_lease(clients.dynamodb, config, dispatch_id, owner, saved["key"], int(time.time()))
    return {"state": "DISPATCHED", "receipt": saved, "content_sha256": released["sha256"]}


def handler(event: dict[str, Any], context: Any) -> dict[str, list[dict[str, str]]]:
    config = Config.from_env()
    clients = _clients()
    failures: list[dict[str, str]] = []
    for record in event.get("Records", []):
        message_id = str(record.get("messageId", ""))
        try:
            result = process_record(record, clients, config)
            print(json.dumps({"state": result["state"], "message_id_sha256": _sha256(message_id.encode("utf-8")),
                              "content_sha256": result["content_sha256"]}, sort_keys=True))
        except Exception as error:  # Lambda/SQS contract requires record-level failure.
            print(json.dumps({"state": "FAILED", "message_id_sha256": _sha256(message_id.encode("utf-8")),
                              "error_type": type(error).__name__}, sort_keys=True))
            failures.append({"itemIdentifier": message_id})
    return {"batchItemFailures": failures}
````

### FILE: `aws_guardduty_magika_idp_dispatch/test_handler.py`
```yaml
block_id: "AWS-GUARDDUTY-MAGIKA-IDP-DISPATCH-GATE:test-handler-py:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-guardduty-magika-idp-dispatch-gate/test_handler.py"
license: "LicenseRef-Workspace-Owner"
sha256: "3bb9ccb724937d4fa5b6a2dc04681d05be8a4b2de642e06447b51a62f8aa6029"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import contextlib
import hashlib
import io
import json
import sys
import tempfile
import types
import unittest
from dataclasses import replace
from pathlib import Path
from types import SimpleNamespace

try:
    from botocore.exceptions import ClientError
except ModuleNotFoundError:
    class ClientError(Exception):
        def __init__(self, response, operation_name):
            self.response, self.operation_name = response, operation_name
            super().__init__(f"{operation_name}: {response['Error']['Code']}")
    boto3_stub = types.ModuleType("boto3")
    boto3_stub.client = lambda name: None
    botocore_stub = types.ModuleType("botocore")
    exceptions_stub = types.ModuleType("botocore.exceptions")
    exceptions_stub.ClientError = ClientError
    sys.modules.setdefault("boto3", boto3_stub)
    sys.modules.setdefault("botocore", botocore_stub)
    sys.modules.setdefault("botocore.exceptions", exceptions_stub)

import handler


def client_error(code: str, operation: str = "Test") -> ClientError:
    return ClientError({"Error": {"Code": code, "Message": code}}, operation)


class Body:
    def __init__(self, data: bytes): self.data, self.offset = data, 0
    def read(self, size: int = -1) -> bytes:
        if size < 0: size = len(self.data) - self.offset
        result = self.data[self.offset:self.offset + size]
        self.offset += len(result)
        return result


class FakeS3:
    def __init__(self):
        self.objects = {}
        self.latest = {}
        self.copy_calls = []
        self.fail_receipt_once = False

    def seed(self, bucket, key, version, data, *, kms, metadata=None, retained=False, etag=None):
        self.objects[(bucket, key, version)] = {
            "Body": data, "Metadata": metadata or {}, "SSEKMSKeyId": kms,
            "ETag": etag or hashlib.md5(data).hexdigest(), "retained": retained,
            "ContentType": "application/octet-stream",
        }
        self.latest[(bucket, key)] = version

    def _row(self, bucket, key, version=None):
        version = version or self.latest.get((bucket, key))
        if (bucket, key, version) not in self.objects: raise client_error("NoSuchKey", "HeadObject")
        return self.objects[(bucket, key, version)], version

    def get_object(self, Bucket, Key, VersionId=None, IfMatch=None, **kwargs):
        row, version = self._row(Bucket, Key, VersionId)
        if IfMatch and IfMatch.strip('"') != row["ETag"].strip('"'): raise client_error("PreconditionFailed", "GetObject")
        return {"Body": Body(row["Body"]), "ServerSideEncryption": "aws:kms", "SSEKMSKeyId": row["SSEKMSKeyId"],
                "VersionId": version, "ETag": row["ETag"]}

    def head_object(self, Bucket, Key, VersionId=None, **kwargs):
        row, version = self._row(Bucket, Key, VersionId)
        return {"ContentLength": len(row["Body"]), "Metadata": row["Metadata"], "ServerSideEncryption": "aws:kms",
                "SSEKMSKeyId": row["SSEKMSKeyId"], "VersionId": version, "ETag": row["ETag"]}

    def copy_object(self, **kwargs):
        self.copy_calls.append(kwargs)
        if kwargs.get("IfNoneMatch") == "*" and (kwargs["Bucket"], kwargs["Key"]) in self.latest:
            raise client_error("PreconditionFailed", "CopyObject")
        source = kwargs["CopySource"]
        row, _ = self._row(source["Bucket"], source["Key"], source["VersionId"])
        if kwargs["CopySourceIfMatch"].strip('"') != row["ETag"].strip('"'):
            raise client_error("PreconditionFailed", "CopyObject")
        version = "idp-v1"
        self.seed(kwargs["Bucket"], kwargs["Key"], version, row["Body"], kms=kwargs["SSEKMSKeyId"],
                  metadata=kwargs["Metadata"], etag=row["ETag"])
        return {"VersionId": version}

    def put_object(self, **kwargs):
        if self.fail_receipt_once and kwargs["Key"].startswith("dispatch-receipts/"):
            self.fail_receipt_once = False
            raise RuntimeError("injected receipt outage")
        if kwargs.get("IfNoneMatch") == "*" and (kwargs["Bucket"], kwargs["Key"]) in self.latest:
            raise client_error("PreconditionFailed", "PutObject")
        data = kwargs["Body"]
        version = "dispatch-v1"
        self.seed(kwargs["Bucket"], kwargs["Key"], version, data, kms=kwargs["SSEKMSKeyId"],
                  metadata=kwargs["Metadata"], retained=True)
        return {"VersionId": version}

    def get_object_retention(self, Bucket, Key, VersionId, **kwargs):
        row, _ = self._row(Bucket, Key, VersionId)
        return {"Retention": {"Mode": "COMPLIANCE", "RetainUntilDate": "future"}} if row["retained"] else {"Retention": {}}


class FakeDDB:
    def __init__(self): self.rows = {}; self.calls = []
    def update_item(self, **kwargs):
        self.calls.append(kwargs)
        key = kwargs["Key"]["DispatchId"]["S"]
        values = kwargs["ExpressionAttributeValues"]
        if ":expires" in values:
            if key in self.rows and self.rows[key].get("state") == "WORKING": raise client_error("ConditionalCheckFailedException", "UpdateItem")
            self.rows[key] = {"state": "WORKING", "owner": values[":owner"]["S"]}
        else:
            if key not in self.rows or self.rows[key]["owner"] != values[":owner"]["S"]: raise client_error("ConditionalCheckFailedException", "UpdateItem")
            self.rows[key]["state"] = "DONE"
        return {}


class Runner:
    def __init__(self, *, label="pdf", mime="application/pdf", score=0.99, extensions=None, version=handler.MAGIKA_VERSION):
        self.label, self.mime, self.score = label, mime, score
        self.extensions, self.version = extensions or ["pdf"], version
    def __call__(self, args, **kwargs):
        if "--version" in args: return SimpleNamespace(returncode=0, stdout=self.version, stderr="")
        output = [{"result": {"status": "ok", "value": {"score": self.score, "output": {
            "label": self.label, "mime_type": self.mime, "extensions": self.extensions}}}}]
        return SimpleNamespace(returncode=0, stdout=json.dumps(output), stderr="")


def config(binary: Path) -> handler.Config:
    return handler.Config(
        account_id="111122223333", region="us-east-1", tenant_id="tenant-a",
        receipt_bucket="security-receipts", receipt_kms_key_arn="arn:kms:receipt",
        release_bucket="release", release_kms_key_arn="arn:kms:release",
        idp_input_bucket="idp-input", idp_input_kms_key_arn="arn:kms:idp",
        dispatch_receipt_bucket="dispatch-receipts", dispatch_receipt_kms_key_arn="arn:kms:dispatch",
        idempotency_table="dispatch-table", config_version="invoice-v1", magika_binary=str(binary),
        magika_sha256=hashlib.sha256(binary.read_bytes()).hexdigest(),
        allowed_types=({"label": "pdf", "mime_type": "application/pdf", "extension": "pdf", "min_score": 0.95},),
        max_bytes=1024, retain_days=365, lease_seconds=900,
    )


def receipt(content: bytes) -> bytes:
    sha = hashlib.sha256(content).hexdigest()
    value = {
        "schema": handler.RECEIPT_SCHEMA, "receipt_version": 1, "scan_id": "a" * 64,
        "tenant_id": "tenant-a", "provider": "AWS_GUARDDUTY_MALWARE_PROTECTION_S3",
        "provider_schema_version": "1.0", "decision": "ADMIT_TO_EXTRACTION",
        "released_object": {"bucket": "release", "key": f"clean/tenant-a/{sha}.bin", "version_id": "release-v1",
                            "sha256": sha, "bytes": len(content), "retention_verified": True},
        "automatic_business_persistence_authorized": False,
        "next_gate": "DOCUMENT_CLASSIFICATION_AND_FIELD_EVIDENCE",
    }
    return handler._canonical_json(value)


def event(receipt_bytes: bytes, *, account="111122223333", region="us-east-1", version="receipt-v1", etag=None):
    key = "security-receipts/tenant-a/" + "a" * 64 + ".json"
    envelope = {"version": "0", "source": "aws.s3", "detail-type": "Object Created", "account": account, "region": region,
                "resources": ["arn:aws:s3:::security-receipts"], "detail": {"event-version": "1.1",
                "bucket": {"name": "security-receipts"}, "object": {"key": key, "version-id": version,
                "etag": etag or hashlib.md5(receipt_bytes).hexdigest(), "size": len(receipt_bytes)}}}
    return {"messageId": "message-1", "body": json.dumps(envelope)}


class DispatchTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.binary = Path(self.temp.name) / "magika"
        self.binary.write_bytes(b"official-magika-binary")
        self.config = config(self.binary)
        self.s3, self.ddb = FakeS3(), FakeDDB()
        self.clients = handler.Clients(self.s3, self.ddb)
        self.content = b"%PDF-1.7\n"
        self.receipt = receipt(self.content)
        self.event = event(self.receipt)
        key = json.loads(self.event["body"])["detail"]["object"]["key"]
        self.s3.seed("security-receipts", key, "receipt-v1", self.receipt, kms="arn:kms:receipt")
        sha = hashlib.sha256(self.content).hexdigest()
        self.s3.seed("release", f"clean/tenant-a/{sha}.bin", "release-v1", self.content, kms="arn:kms:release",
                     metadata={"content-sha256": sha, "raw-sha256": "b" * 64, "tenant-id": "tenant-a"}, etag="release-etag")

    def tearDown(self): self.temp.cleanup()

    def process(self, runner=None):
        return handler.process_record(self.event, self.clients, self.config, runner or Runner())

    def test_clean_exact_version_is_identified_copied_and_receipted(self):
        result = self.process()
        self.assertEqual("DISPATCHED", result["state"])
        call = self.s3.copy_calls[0]
        self.assertEqual("release-v1", call["CopySource"]["VersionId"])
        self.assertEqual("release-etag", call["CopySourceIfMatch"])
        self.assertEqual("*", call["IfNoneMatch"])
        self.assertTrue(call["Key"].endswith(".pdf"))
        self.assertEqual("invoice-v1", call["Metadata"]["config-version"])
        saved = self.s3.objects[("dispatch-receipts", result["receipt"]["key"], "dispatch-v1")]["Body"]
        dispatch = json.loads(saved)
        self.assertFalse(dispatch["automatic_business_persistence_authorized"])
        self.assertEqual("v0.6.5", dispatch["idp_engine"]["release"])

    def test_duplicate_short_circuits_before_magika_and_copy(self):
        first = self.process()
        second = self.process(Runner(version="must-not-run"))
        self.assertEqual("DUPLICATE", second["state"])
        self.assertEqual(1, len(self.s3.copy_calls))
        self.assertEqual(first["receipt"]["key"], second["receipt"]["key"])

    def test_rejected_scan_never_downloads_or_copies(self):
        value = json.loads(self.receipt); value["decision"] = "REJECTED_THREATS_FOUND"
        bad = handler._canonical_json(value)
        key = json.loads(self.event["body"])["detail"]["object"]["key"]
        self.s3.seed("security-receipts", key, "receipt-v1", bad, kms="arn:kms:receipt")
        self.event = event(bad)
        with self.assertRaisesRegex(handler.Rejected, "did not admit"): self.process()
        self.assertEqual([], self.s3.copy_calls)

    def test_account_region_resource_and_major_version_are_closed(self):
        for mutate in (
            lambda body: body.__setitem__("account", "999900001111"),
            lambda body: body.__setitem__("region", "eu-west-1"),
            lambda body: body.__setitem__("resources", ["arn:aws:s3:::other"]),
            lambda body: body["detail"].__setitem__("event-version", "2.0"),
        ):
            body = json.loads(self.event["body"]); mutate(body)
            with self.assertRaises(handler.Rejected): handler._parse_event({"body": json.dumps(body)}, self.config)

    def test_receipt_requires_exact_version_etag_kms_and_schema(self):
        bad_event = event(self.receipt, version="wrong")
        with self.assertRaises(ClientError): handler._read_receipt(self.s3, handler._parse_event(bad_event, self.config), self.config)
        wrong = json.loads(self.receipt); wrong.pop("schema")
        data = handler._canonical_json(wrong); key = json.loads(self.event["body"])["detail"]["object"]["key"]
        self.s3.seed("security-receipts", key, "receipt-v1", data, kms="arn:kms:receipt")
        with self.assertRaises(handler.Rejected): handler._read_receipt(self.s3, handler._parse_event(event(data), self.config), self.config)

    def test_release_location_metadata_encryption_and_hash_are_bound(self):
        sha = hashlib.sha256(self.content).hexdigest(); key = f"clean/tenant-a/{sha}.bin"
        for field, value in (("SSEKMSKeyId", "wrong"), ("Metadata", {"content-sha256": "0" * 64, "tenant-id": "tenant-a"}), ("Body", b"changed")):
            row = self.s3.objects[("release", key, "release-v1")]; old = row[field]; row[field] = value
            with self.assertRaises(handler.Rejected): self.process()
            row[field] = old

    def test_magika_version_type_score_and_extension_fail_closed(self):
        cases = [Runner(version="wrong"), Runner(label="zip", mime="application/zip"), Runner(score=0.5), Runner(extensions=["bin"])]
        for candidate in cases:
            self.ddb = FakeDDB()
            self.clients = handler.Clients(self.s3, self.ddb)
            with self.assertRaises(handler.Rejected): self.process(candidate)

    def test_preexisting_correct_idp_copy_reconciles(self):
        sha = hashlib.sha256(self.content).hexdigest(); key = f"security-released/tenant-a/{sha}.pdf"
        metadata = {"config-version": "invoice-v1", "content-sha256": sha, "raw-sha256": "b" * 64,
                    "tenant-id": "tenant-a", "security-receipt-sha256": hashlib.sha256(self.receipt).hexdigest(),
                    "magika-label": "pdf", "magika-mime": "application/pdf"}
        self.s3.seed("idp-input", key, "existing-v", self.content, kms="arn:kms:idp", metadata=metadata)
        self.assertEqual("DISPATCHED", self.process()["state"])

    def test_preexisting_divergent_idp_copy_fails(self):
        sha = hashlib.sha256(self.content).hexdigest(); key = f"security-released/tenant-a/{sha}.pdf"
        self.s3.seed("idp-input", key, "existing-v", b"wrong", kms="arn:kms:idp", metadata={})
        with self.assertRaises(handler.Rejected): self.process()

    def test_receipt_outage_leaves_copy_for_explicit_reconciliation(self):
        self.s3.fail_receipt_once = True
        with self.assertRaisesRegex(RuntimeError, "injected"): self.process()
        self.assertEqual(1, len(self.s3.copy_calls))
        self.assertTrue(any(bucket == "idp-input" for bucket, _ in self.s3.latest))

    def test_active_lease_is_retryable(self):
        dispatch_id = hashlib.sha256(f"tenant-a\0{hashlib.sha256(self.content).hexdigest()}\0invoice-v1".encode()).hexdigest()
        self.ddb.rows[dispatch_id] = {"state": "WORKING", "owner": "other"}
        with self.assertRaises(handler.Retryable): self.process()

    def test_configuration_rejects_empty_types_and_duplicate_identity(self):
        with self.assertRaises(handler.Rejected): replace(self.config, allowed_types=()).validate()
        duplicate = self.config.allowed_types + self.config.allowed_types
        with self.assertRaises(handler.Rejected): replace(self.config, allowed_types=duplicate).validate()

    def test_sqs_partial_batch_and_redacted_log(self):
        original_config, original_clients, original_identify = handler.Config.from_env, handler._clients, handler._identify
        try:
            handler.Config.from_env = staticmethod(lambda: self.config)
            handler._clients = lambda: self.clients
            handler._identify = lambda path, config, runner=handler.subprocess.run: {
                "label": "pdf", "mime_type": "application/pdf", "score": 0.99,
                "extension": "pdf", "version": handler.MAGIKA_VERSION,
            }
            failed = {"messageId": "secret-message", "body": "{}"}
            output = io.StringIO()
            with contextlib.redirect_stdout(output): result = handler.handler({"Records": [self.event, failed]}, None)
            self.assertEqual([{"itemIdentifier": "secret-message"}], result["batchItemFailures"])
            self.assertNotIn("secret-message", output.getvalue())
            self.assertNotIn("security-receipts/tenant-a", output.getvalue())
        finally:
            handler.Config.from_env, handler._clients, handler._identify = original_config, original_clients, original_identify

    def test_event_duplicate_json_key_rejected(self):
        with self.assertRaisesRegex(handler.Rejected, "duplicate JSON key"):
            handler._parse_event({"body": '{"version":"0","version":"0"}'}, self.config)


if __name__ == "__main__":
    unittest.main()
````

### FILE: `aws_guardduty_magika_idp_dispatch/template.yaml`
```yaml
block_id: "AWS-GUARDDUTY-MAGIKA-IDP-DISPATCH-GATE:template-yaml:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-guardduty-magika-idp-dispatch-gate/template.yaml"
license: "LicenseRef-Workspace-Owner"
sha256: "54d494d1f59e238129e00eced566f910733bd9fed82b1c8112087f6b86147997"
variables: []
secrets_allowed: false
```
````yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Description: Version-bound GuardDuty release through Google Magika to AWS GenAI IDP v0.6.5 input

Parameters:
  ExpectedAccountId:
    Type: String
    AllowedPattern: '^[0-9]{12}$'
  TenantId:
    Type: String
    AllowedPattern: '^[a-z0-9][a-z0-9-]{1,62}$'
  SecurityReceiptBucketName:
    Type: String
  SecurityReceiptKmsKeyArn:
    Type: String
  ReleaseBucketName:
    Type: String
  ReleaseKmsKeyArn:
    Type: String
  IdpInputBucketName:
    Type: String
  IdpInputKmsKeyArn:
    Type: String
  IdpConfigVersion:
    Type: String
    AllowedPattern: '^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$'
  AllowedMagikaTypesJson:
    Type: String
    Default: '[]'
  MaxDocumentBytes:
    Type: Number
    Default: 20971520
    MinValue: 1
    MaxValue: 20971520
  RetainDays:
    Type: Number
    Default: 365
    MinValue: 1
  SecurityAlertEmail:
    Type: String
    AllowedPattern: '^[^@\s]+@[^@\s]+\.[^@\s]+$'

Resources:
  DispatchKey:
    Type: AWS::KMS::Key
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      Description: KMS key for secure document dispatch receipts, queues and table
      EnableKeyRotation: true
      PendingWindowInDays: 30
      KeyPolicy:
        Version: '2012-10-17'
        Statement:
          - Sid: RootAccountAdministration
            Effect: Allow
            Principal: {AWS: !Sub 'arn:${AWS::Partition}:iam::${AWS::AccountId}:root'}
            Action: 'kms:*'
            Resource: '*'
          - Sid: AllowDispatchFunctionUse
            Effect: Allow
            Principal: {AWS: !GetAtt DispatchFunctionRole.Arn}
            Action: [kms:Encrypt, kms:Decrypt, kms:GenerateDataKey, kms:DescribeKey]
            Resource: '*'
          - Sid: AllowSqsUse
            Effect: Allow
            Principal: {Service: sqs.amazonaws.com}
            Action: [kms:Decrypt, kms:GenerateDataKey]
            Resource: '*'

  DispatchReceiptBucket:
    Type: AWS::S3::Bucket
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      BucketEncryption:
        ServerSideEncryptionConfiguration:
          - BucketKeyEnabled: true
            ServerSideEncryptionByDefault: {SSEAlgorithm: aws:kms, KMSMasterKeyID: !GetAtt DispatchKey.Arn}
      ObjectLockEnabled: true
      ObjectLockConfiguration:
        ObjectLockEnabled: Enabled
        Rule:
          DefaultRetention: {Mode: COMPLIANCE, Days: !Ref RetainDays}
      VersioningConfiguration: {Status: Enabled}
      PublicAccessBlockConfiguration:
        BlockPublicAcls: true
        BlockPublicPolicy: true
        IgnorePublicAcls: true
        RestrictPublicBuckets: true
      OwnershipControls:
        Rules: [{ObjectOwnership: BucketOwnerEnforced}]

  DispatchReceiptBucketPolicy:
    Type: AWS::S3::BucketPolicy
    Properties:
      Bucket: !Ref DispatchReceiptBucket
      PolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Sid: DenyInsecureTransport
            Effect: Deny
            Principal: '*'
            Action: 's3:*'
            Resource: [!GetAtt DispatchReceiptBucket.Arn, !Sub '${DispatchReceiptBucket.Arn}/*']
            Condition: {Bool: {'aws:SecureTransport': 'false'}}

  DispatchTable:
    Type: AWS::DynamoDB::Table
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      BillingMode: PAY_PER_REQUEST
      AttributeDefinitions: [{AttributeName: DispatchId, AttributeType: S}]
      KeySchema: [{AttributeName: DispatchId, KeyType: HASH}]
      DeletionProtectionEnabled: true
      PointInTimeRecoverySpecification: {PointInTimeRecoveryEnabled: true}
      SSESpecification: {SSEEnabled: true, SSEType: KMS, KMSMasterKeyId: !GetAtt DispatchKey.Arn}

  DispatchDLQ:
    Type: AWS::SQS::Queue
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      KmsMasterKeyId: !GetAtt DispatchKey.Arn
      MessageRetentionPeriod: 1209600
      SqsManagedSseEnabled: false

  DispatchQueue:
    Type: AWS::SQS::Queue
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      KmsMasterKeyId: !GetAtt DispatchKey.Arn
      MessageRetentionPeriod: 1209600
      VisibilityTimeout: 1800
      RedrivePolicy:
        deadLetterTargetArn: !GetAtt DispatchDLQ.Arn
        maxReceiveCount: 5

  DispatchQueuePolicy:
    Type: AWS::SQS::QueuePolicy
    Properties:
      Queues: [!Ref DispatchQueue]
      PolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Sid: EventBridgeSendOnly
            Effect: Allow
            Principal: {Service: events.amazonaws.com}
            Action: sqs:SendMessage
            Resource: !GetAtt DispatchQueue.Arn
            Condition: {ArnEquals: {'aws:SourceArn': !GetAtt SecurityReceiptRule.Arn}}

  SecurityReceiptRule:
    Type: AWS::Events::Rule
    Properties:
      State: ENABLED
      EventPattern:
        source: [aws.s3]
        detail-type: [Object Created]
        account: [!Ref ExpectedAccountId]
        region: [!Ref 'AWS::Region']
        resources: [!Sub 'arn:${AWS::Partition}:s3:::${SecurityReceiptBucketName}']
        detail:
          bucket: {name: [!Ref SecurityReceiptBucketName]}
          object:
            key:
              - prefix: !Sub 'security-receipts/${TenantId}/'
      Targets:
        - Arn: !GetAtt DispatchQueue.Arn
          Id: dispatch-queue
          DeadLetterConfig: {Arn: !GetAtt DispatchDLQ.Arn}
          RetryPolicy: {MaximumEventAgeInSeconds: 86400, MaximumRetryAttempts: 10}

  DispatchFunction:
    Type: AWS::Serverless::Function
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      CodeUri: build/function
      Handler: handler.handler
      Runtime: python3.14
      Architectures: [x86_64]
      Timeout: 300
      MemorySize: 1024
      EphemeralStorage: {Size: 1024}
      ReservedConcurrentExecutions: 10
      Tracing: Active
      AutoPublishAlias: live
      DeploymentPreference: {Type: Canary10Percent5Minutes}
      Environment:
        Variables:
          EXPECTED_ACCOUNT_ID: !Ref ExpectedAccountId
          EXPECTED_REGION: !Ref 'AWS::Region'
          TENANT_ID: !Ref TenantId
          SECURITY_RECEIPT_BUCKET: !Ref SecurityReceiptBucketName
          SECURITY_RECEIPT_KMS_KEY_ARN: !Ref SecurityReceiptKmsKeyArn
          RELEASE_BUCKET: !Ref ReleaseBucketName
          RELEASE_KMS_KEY_ARN: !Ref ReleaseKmsKeyArn
          IDP_INPUT_BUCKET: !Ref IdpInputBucketName
          IDP_INPUT_KMS_KEY_ARN: !Ref IdpInputKmsKeyArn
          DISPATCH_RECEIPT_BUCKET: !Ref DispatchReceiptBucket
          DISPATCH_RECEIPT_KMS_KEY_ARN: !GetAtt DispatchKey.Arn
          IDEMPOTENCY_TABLE: !Ref DispatchTable
          IDP_CONFIG_VERSION: !Ref IdpConfigVersion
          MAGIKA_BINARY: /var/task/bin/magika
          MAGIKA_SHA256: 03a56971f9e121666d702db611fd08f1c8bae5a3c31538f1b45b3ef99e6df84d
          ALLOWED_MAGIKA_TYPES_JSON: !Ref AllowedMagikaTypesJson
          MAX_DOCUMENT_BYTES: !Ref MaxDocumentBytes
          RETAIN_DAYS: !Ref RetainDays
          LEASE_SECONDS: 900
      Policies:
        - Statement:
            - Sid: ReadExactSecurityReceipts
              Effect: Allow
              Action: [s3:GetObject, s3:GetObjectVersion]
              Resource: !Sub 'arn:${AWS::Partition}:s3:::${SecurityReceiptBucketName}/security-receipts/${TenantId}/*'
            - Sid: ReadExactReleasedDocuments
              Effect: Allow
              Action: [s3:GetObject, s3:GetObjectVersion]
              Resource: !Sub 'arn:${AWS::Partition}:s3:::${ReleaseBucketName}/clean/${TenantId}/*'
            - Sid: WriteOnlyContentAddressedIdpInput
              Effect: Allow
              Action: [s3:PutObject]
              Resource: !Sub 'arn:${AWS::Partition}:s3:::${IdpInputBucketName}/security-released/${TenantId}/*'
            - Sid: VerifyIdpInputVersions
              Effect: Allow
              Action: [s3:GetObject, s3:GetObjectVersion]
              Resource: !Sub 'arn:${AWS::Partition}:s3:::${IdpInputBucketName}/security-released/${TenantId}/*'
            - Sid: DispatchReceiptEvidence
              Effect: Allow
              Action: [s3:GetObject, s3:GetObjectVersion, s3:PutObject, s3:GetObjectRetention]
              Resource: !Sub '${DispatchReceiptBucket.Arn}/dispatch-receipts/${TenantId}/*'
            - Sid: DispatchLease
              Effect: Allow
              Action: [dynamodb:UpdateItem]
              Resource: !GetAtt DispatchTable.Arn
            - Sid: KmsUse
              Effect: Allow
              Action: [kms:Decrypt, kms:Encrypt, kms:GenerateDataKey, kms:DescribeKey]
              Resource: [!Ref SecurityReceiptKmsKeyArn, !Ref ReleaseKmsKeyArn, !Ref IdpInputKmsKeyArn, !GetAtt DispatchKey.Arn]
      Events:
        DispatchBatch:
          Type: SQS
          Properties:
            Queue: !GetAtt DispatchQueue.Arn
            BatchSize: 10
            FunctionResponseTypes: [ReportBatchItemFailures]
            ScalingConfig: {MaximumConcurrency: 10}

  AlertTopic:
    Type: AWS::SNS::Topic
    Properties:
      KmsMasterKeyId: !GetAtt DispatchKey.Arn
      Subscription: [{Endpoint: !Ref SecurityAlertEmail, Protocol: email}]

  DispatchDlqAlarm:
    Type: AWS::CloudWatch::Alarm
    Properties:
      AlarmActions: [!Ref AlertTopic]
      ComparisonOperator: GreaterThanThreshold
      EvaluationPeriods: 1
      MetricName: ApproximateNumberOfMessagesVisible
      Namespace: AWS/SQS
      Period: 60
      Statistic: Maximum
      Threshold: 0
      TreatMissingData: notBreaching
      Dimensions: [{Name: QueueName, Value: !GetAtt DispatchDLQ.QueueName}]

  DispatchFunctionErrorsAlarm:
    Type: AWS::CloudWatch::Alarm
    Properties:
      AlarmActions: [!Ref AlertTopic]
      ComparisonOperator: GreaterThanThreshold
      EvaluationPeriods: 2
      MetricName: Errors
      Namespace: AWS/Lambda
      Period: 60
      Statistic: Sum
      Threshold: 0
      TreatMissingData: notBreaching
      Dimensions: [{Name: FunctionName, Value: !Ref DispatchFunction}]

Outputs:
  DispatchReceiptBucketName: {Value: !Ref DispatchReceiptBucket}
  DispatchQueueUrl: {Value: !Ref DispatchQueue}
  DispatchTableName: {Value: !Ref DispatchTable}
````

### FILE: `aws_guardduty_magika_idp_dispatch/provider-profile.template.json`
```yaml
block_id: "AWS-GUARDDUTY-MAGIKA-IDP-DISPATCH-GATE:provider-profile-template-json:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-guardduty-magika-idp-dispatch-gate/provider-profile.template.json"
license: "LicenseRef-Workspace-Owner"
sha256: "9d0d8fc656f47db039e166a65600bf95d97eb0bd053ac4413dfe529992e842ba"
variables: []
secrets_allowed: false
```
````json
{
  "profile_version": 1,
  "aws_account_id": "",
  "aws_region": "",
  "stack_name": "",
  "tenant_id": "",
  "security_receipt_bucket_name": "",
  "security_receipt_kms_key_arn": "",
  "release_bucket_name": "",
  "release_kms_key_arn": "",
  "idp_input_bucket_name": "",
  "idp_input_kms_key_arn": "",
  "aws_idp_stack_name": "",
  "idp_config_version": "",
  "allowed_magika_types": [],
  "max_document_bytes": 20971520,
  "retain_days": 365,
  "security_owner": "",
  "document_operations_owner": "",
  "cost_owner": "",
  "replay_owner": "",
  "aws_idp_release": "v0.6.5",
  "aws_idp_commit": "1b5fd74454e593de233342a02ee911af8ee38359",
  "magika_release": "cli/v1.1.0",
  "magika_commit": "5e2f437fb7b7452368c8c1fa9354858f5487a5c4",
  "magika_linux_asset_sha256": "6b4c1010c84d1f4f06205ccef4597f1690bcd7744f46d841eee26426bc100485",
  "aws_idp_source_acquired_and_verified": false,
  "aws_idp_stack_outputs_verified": false,
  "guardduty_canaries_proven": false,
  "idp_config_corpus_evaluated": false,
  "downstream_business_idempotency_proven": false,
  "approve_live_effects": false,
  "approve_provider_costs": false
}
````

### FILE: `aws_guardduty_magika_idp_dispatch/requirements-runtime.lock`
```yaml
block_id: "AWS-GUARDDUTY-MAGIKA-IDP-DISPATCH-GATE:requirements-runtime-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-guardduty-magika-idp-dispatch-gate/requirements-runtime.lock"
license: "LicenseRef-Workspace-Owner"
sha256: "95c37a7ffe3d5e3e3dc9c188d4370a934e7d853c2c3494c99cc7884e54df24cc"
variables: []
secrets_allowed: false
```
````text
boto3==1.43.83 \
    --hash=sha256:73a3564f737d4516625964eee709a498fa98ccee6aca929febad2b0b5fbeae1e
botocore==1.43.83 \
    --hash=sha256:bf75a6cf587c22d968e43e79fe122c39f82deafbe9c3422bc5d3e80b6210fc98
jmespath==1.1.0 \
    --hash=sha256:a5663118de4908c91729bea0acadca56526eb2698e83de10cd116ae0f4e97c64
python-dateutil==2.9.0.post0 \
    --hash=sha256:a8b2bc7bffae282281c8140a97d3aa9c14da0b136dfe83f850eea9a5f7470427
s3transfer==0.19.2 \
    --hash=sha256:d8168eccca828cbb2cd573675333f3bddd254313a9c42494b84c76b539e8ba25
six==1.17.0 \
    --hash=sha256:4721f391ed90541fddacab5acf947aa0d3dc7d27b2e1e8eda2be8970586c3274
urllib3==2.7.0 \
    --hash=sha256:9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897
````

### FILE: `aws_guardduty_magika_idp_dispatch/build.ps1`
```yaml
block_id: "AWS-GUARDDUTY-MAGIKA-IDP-DISPATCH-GATE:build-ps1:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-guardduty-magika-idp-dispatch-gate/build.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "372530af352671b1697704886df8afe464faf405f9fe4aa646f7e8d24351a56d"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)][string] $ComposedRoot,
  [Parameter(Mandatory = $true)][string] $MagikaArchive,
  [Parameter(Mandatory = $true)][string] $MagikaChecksumFile,
  [Parameter(Mandatory = $true)][string] $OutputRoot,
  [string] $PythonExecutable = 'python'
)

$ErrorActionPreference = 'Stop'
$root = (Resolve-Path -LiteralPath $ComposedRoot).Path
$component = Join-Path $root 'aws_guardduty_magika_idp_dispatch'
if (-not (Test-Path -LiteralPath $component -PathType Container)) { throw 'dispatch component is missing' }
$archive = (Resolve-Path -LiteralPath $MagikaArchive).Path
$sidecar = (Resolve-Path -LiteralPath $MagikaChecksumFile).Path
if ((Get-Item -LiteralPath $archive).Length -ne 8625112) { throw 'Magika archive size mismatch' }
if ((Get-FileHash -Algorithm SHA256 -LiteralPath $archive).Hash.ToLowerInvariant() -cne '6b4c1010c84d1f4f06205ccef4597f1690bcd7744f46d841eee26426bc100485') { throw 'Magika archive SHA-256 mismatch' }
if ((Get-Item -LiteralPath $sidecar).Length -ne 110 -or (Get-FileHash -Algorithm SHA256 -LiteralPath $sidecar).Hash.ToLowerInvariant() -cne '4028a7c1fe789dea5f84e879a6dec00733cf042eff3ee7089bee3962fd923e19') { throw 'Magika sidecar identity mismatch' }
$sidecarText = [IO.File]::ReadAllText($sidecar).Trim()
if ($sidecarText -cne '6b4c1010c84d1f4f06205ccef4597f1690bcd7744f46d841eee26426bc100485 *magika-cli-x86_64-unknown-linux-gnu.tar.xz') { throw 'Magika sidecar content mismatch' }
$tar = Get-Command tar -ErrorAction Stop

$output = [IO.Path]::GetFullPath($OutputRoot)
if (Test-Path -LiteralPath $output) { throw 'output must not exist' }
$parent = Split-Path -Parent $output
if (-not (Test-Path -LiteralPath $parent -PathType Container)) { throw 'output parent is missing' }
$temp = Join-Path $parent ('.' + [IO.Path]::GetFileName($output) + '.' + [Guid]::NewGuid().ToString('N'))
$extract = Join-Path $temp 'magika-extract'
$function = Join-Path $temp 'build/function'

try {
  New-Item -ItemType Directory -Path $extract,$function | Out-Null
  & $tar.Source -xJf $archive -C $extract
  if ($LASTEXITCODE -ne 0) { throw 'Magika archive extraction failed' }
  $binaries = @(Get-ChildItem -LiteralPath $extract -Recurse -File | Where-Object Name -CEQ 'magika')
  if ($binaries.Count -ne 1) { throw 'Magika archive must contain exactly one magika binary' }
  $binary = $binaries[0]
  if ($binary.Length -ne 29225280 -or (Get-FileHash -Algorithm SHA256 -LiteralPath $binary.FullName).Hash.ToLowerInvariant() -cne '03a56971f9e121666d702db611fd08f1c8bae5a3c31538f1b45b3ef99e6df84d') { throw 'Magika binary identity mismatch' }

  Copy-Item -LiteralPath (Join-Path $component 'handler.py') -Destination $function
  Copy-Item -LiteralPath (Join-Path $component 'requirements-runtime.lock') -Destination (Join-Path $function 'requirements.txt')
  Copy-Item -LiteralPath (Join-Path $component 'GOOGLE_APACHE_LICENSE.txt') -Destination $function
  Copy-Item -LiteralPath (Join-Path $component 'AWS_MIT_NO_ATTRIBUTION_LICENSE.txt') -Destination $function
  Copy-Item -LiteralPath (Join-Path $component 'template.yaml') -Destination $temp

  & $PythonExecutable -m pip install --disable-pip-version-check --require-hashes --only-binary=:all: --no-compile --target $function -r (Join-Path $function 'requirements.txt')
  if ($LASTEXITCODE -ne 0) { throw 'hash-locked dependency installation failed' }
  & $PythonExecutable -B -c "import sys; sys.path.insert(0, r'$function'); import boto3, botocore; assert boto3.__version__ == '1.43.83'; assert botocore.__version__ == '1.43.83'"
  if ($LASTEXITCODE -ne 0) { throw 'dependency version check failed' }
  $binDirectory = Join-Path $function 'bin'
  if (-not (Test-Path -LiteralPath $binDirectory -PathType Container)) { New-Item -ItemType Directory -Path $binDirectory | Out-Null }
  $magikaDestination = Join-Path $binDirectory 'magika'
  if (Test-Path -LiteralPath $magikaDestination) { throw 'dependency install unexpectedly created bin/magika' }
  Copy-Item -LiteralPath $binary.FullName -Destination $magikaDestination

  $priorPythonPath = $env:PYTHONPATH
  $priorBytecode = $env:PYTHONDONTWRITEBYTECODE
  try {
    $env:PYTHONPATH = "$component$([IO.Path]::PathSeparator)$function"
    $env:PYTHONDONTWRITEBYTECODE = '1'
    & $PythonExecutable -B -m py_compile (Join-Path $component 'handler.py') (Join-Path $component 'test_handler.py')
    if ($LASTEXITCODE -ne 0) { throw 'Python compile failed' }
    Push-Location $component
    try { & $PythonExecutable -B -m unittest -v test_handler.py; if ($LASTEXITCODE -ne 0) { throw 'dispatch tests failed' } }
    finally { Pop-Location }
  } finally {
    $env:PYTHONPATH = $priorPythonPath
    $env:PYTHONDONTWRITEBYTECODE = $priorBytecode
  }

  $files = @(Get-ChildItem -LiteralPath $temp -Recurse -File | Where-Object { $_.FullName -notlike "$extract*" } | ForEach-Object {
    @{path=[IO.Path]::GetRelativePath($temp,$_.FullName).Replace('\','/');bytes=$_.Length;sha256=(Get-FileHash -Algorithm SHA256 -LiteralPath $_.FullName).Hash.ToLowerInvariant()}
  } | Sort-Object path)
  $receipt = [ordered]@{
    schema='elite-aws-magika-idp-dispatch-build-receipt/v1';
    dependencies=@{boto3='1.43.83';botocore='1.43.83'};
    magika=@{release='cli/v1.1.0';commit='5e2f437fb7b7452368c8c1fa9354858f5487a5c4';archive_sha256='6b4c1010c84d1f4f06205ccef4597f1690bcd7744f46d841eee26426bc100485';binary_sha256='03a56971f9e121666d702db611fd08f1c8bae5a3c31538f1b45b3ef99e6df84d'};
    aws_idp=@{release='v0.6.5';commit='1b5fd74454e593de233342a02ee911af8ee38359'};
    tests=14;files=$files
  }
  $receipt | ConvertTo-Json -Depth 8 | Set-Content -LiteralPath (Join-Path $temp 'BUILD_RECEIPT.json') -Encoding utf8NoBOM
  Remove-Item -LiteralPath $extract -Recurse -Force
  Move-Item -LiteralPath $temp -Destination $output
  Write-Output "AWS_GUARDDUTY_MAGIKA_IDP_DISPATCH_BUILD_PASS output=$output files=$($files.Count) tests=14"
} catch {
  if (Test-Path -LiteralPath $temp) { Remove-Item -LiteralPath $temp -Recurse -Force }
  throw
}
````

### FILE: `aws_guardduty_magika_idp_dispatch/deploy.ps1`
```yaml
block_id: "AWS-GUARDDUTY-MAGIKA-IDP-DISPATCH-GATE:deploy-ps1:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-guardduty-magika-idp-dispatch-gate/deploy.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "d1de067f5638832582f0706217a5457f72c5946c31bd5117e8a79a4b52e9e581"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)][string] $BuildRoot,
  [Parameter(Mandatory = $true)][string] $ProviderProfile,
  [switch] $Execute
)

$ErrorActionPreference = 'Stop'
$build = (Resolve-Path -LiteralPath $BuildRoot).Path
$profile = Get-Content -Raw -LiteralPath (Resolve-Path -LiteralPath $ProviderProfile).Path | ConvertFrom-Json -Depth 30
if ($profile.profile_version -ne 1) { throw 'unsupported provider profile' }
$required = @('aws_account_id','aws_region','stack_name','tenant_id','security_receipt_bucket_name','security_receipt_kms_key_arn','release_bucket_name','release_kms_key_arn','idp_input_bucket_name','idp_input_kms_key_arn','aws_idp_stack_name','idp_config_version','security_owner','document_operations_owner','cost_owner','replay_owner')
foreach ($field in $required) { if ([string]::IsNullOrWhiteSpace([string]$profile.$field)) { throw "profile field is required: $field" } }
if ($profile.aws_account_id -notmatch '^[0-9]{12}$' -or $profile.tenant_id -notmatch '^[a-z0-9][a-z0-9-]{1,62}$' -or $profile.idp_config_version -notmatch '^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$') { throw 'profile identity is invalid' }
if ($profile.aws_idp_release -cne 'v0.6.5' -or $profile.aws_idp_commit -cne '1b5fd74454e593de233342a02ee911af8ee38359') { throw 'AWS IDP identity drifted' }
if ($profile.magika_release -cne 'cli/v1.1.0' -or $profile.magika_commit -cne '5e2f437fb7b7452368c8c1fa9354858f5487a5c4' -or $profile.magika_linux_asset_sha256 -cne '6b4c1010c84d1f4f06205ccef4597f1690bcd7744f46d841eee26426bc100485') { throw 'Magika identity drifted' }
$types = @($profile.allowed_magika_types)
if ($types.Count -eq 0) { throw 'allowed_magika_types must be non-empty' }
foreach ($type in $types) {
  if (@($type.PSObject.Properties.Name | Sort-Object) -join ',' -cne 'extension,label,mime_type,min_score') { throw 'allowed Magika type fields are invalid' }
  if ($type.label -notmatch '^[a-z0-9][a-z0-9_+-]{0,63}$' -or $type.mime_type -notmatch '^[a-z0-9.+-]+/[a-z0-9.+-]+$' -or $type.extension -notmatch '^[a-z0-9]{1,10}$' -or [double]$type.min_score -lt 0.5 -or [double]$type.min_score -gt 1) { throw 'allowed Magika type value is invalid' }
}
foreach ($proof in @('aws_idp_source_acquired_and_verified','aws_idp_stack_outputs_verified','guardduty_canaries_proven','idp_config_corpus_evaluated','downstream_business_idempotency_proven','approve_live_effects','approve_provider_costs')) {
  if ($profile.$proof -ne $true) { throw "live gate is not proven: $proof" }
}
if ([int64]$profile.max_document_bytes -lt 1 -or [int64]$profile.max_document_bytes -gt 20971520 -or [int]$profile.retain_days -lt 1) { throw 'profile limits are invalid' }

$aws = (Get-Command aws -ErrorAction Stop).Source
$sam = (Get-Command sam -ErrorAction Stop).Source
$identity = & $aws sts get-caller-identity --output json | ConvertFrom-Json
if ($LASTEXITCODE -ne 0 -or $identity.Account -cne $profile.aws_account_id) { throw 'AWS caller account mismatch' }
foreach ($bucketSpec in @(
  @{name=$profile.security_receipt_bucket_name;kms=$profile.security_receipt_kms_key_arn},
  @{name=$profile.release_bucket_name;kms=$profile.release_kms_key_arn},
  @{name=$profile.idp_input_bucket_name;kms=$profile.idp_input_kms_key_arn}
)) {
  $versioning = & $aws s3api get-bucket-versioning --bucket $bucketSpec.name --region $profile.aws_region --output json | ConvertFrom-Json
  if ($LASTEXITCODE -ne 0 -or $versioning.Status -cne 'Enabled') { throw "bucket is not versioned: $($bucketSpec.name)" }
  $encryption = & $aws s3api get-bucket-encryption --bucket $bucketSpec.name --region $profile.aws_region --output json | ConvertFrom-Json -Depth 20
  $rules = @($encryption.ServerSideEncryptionConfiguration.Rules | Where-Object { $_.ApplyServerSideEncryptionByDefault.SSEAlgorithm -eq 'aws:kms' -and $_.ApplyServerSideEncryptionByDefault.KMSMasterKeyID -eq $bucketSpec.kms })
  if ($LASTEXITCODE -ne 0 -or $rules.Count -ne 1) { throw "bucket KMS mismatch: $($bucketSpec.name)" }
}
$idpStack = & $aws cloudformation describe-stacks --stack-name $profile.aws_idp_stack_name --region $profile.aws_region --output json | ConvertFrom-Json -Depth 30
if ($LASTEXITCODE -ne 0 -or @($idpStack.Stacks).Count -ne 1) { throw 'AWS IDP stack was not resolved exactly once' }
$idpInputOutput = @($idpStack.Stacks[0].Outputs | Where-Object OutputKey -in @('S3InputBucketName','InputBucketName'))
if ($idpInputOutput.Count -ne 1 -or $idpInputOutput[0].OutputValue -cne $profile.idp_input_bucket_name) { throw 'AWS IDP input bucket output mismatch' }

$template = Join-Path $build 'template.yaml'
if (-not (Test-Path -LiteralPath (Join-Path $build 'build/function/handler.py') -PathType Leaf)) { throw 'built function is missing' }
& $sam validate --lint --template-file $template --region $profile.aws_region
if ($LASTEXITCODE -ne 0) { throw 'SAM validation failed' }
if (-not $Execute) { Write-Output 'AWS_GUARDDUTY_MAGIKA_IDP_DISPATCH_PREFLIGHT_PASS execute=0'; return }
$typesJson = ConvertTo-Json -Compress -Depth 8 -InputObject $types
$deployArgs = @(
  'deploy','--template-file',$template,'--stack-name',$profile.stack_name,'--region',$profile.aws_region,
  '--capabilities','CAPABILITY_IAM','--no-fail-on-empty-changeset','--parameter-overrides',
  "ExpectedAccountId=$($profile.aws_account_id)", "TenantId=$($profile.tenant_id)",
  "SecurityReceiptBucketName=$($profile.security_receipt_bucket_name)", "SecurityReceiptKmsKeyArn=$($profile.security_receipt_kms_key_arn)",
  "ReleaseBucketName=$($profile.release_bucket_name)", "ReleaseKmsKeyArn=$($profile.release_kms_key_arn)",
  "IdpInputBucketName=$($profile.idp_input_bucket_name)", "IdpInputKmsKeyArn=$($profile.idp_input_kms_key_arn)",
  "IdpConfigVersion=$($profile.idp_config_version)", "AllowedMagikaTypesJson=$typesJson",
  "MaxDocumentBytes=$($profile.max_document_bytes)", "RetainDays=$($profile.retain_days)", "SecurityAlertEmail=$($profile.security_owner)"
)
& $sam @deployArgs
if ($LASTEXITCODE -ne 0) { throw 'SAM deploy failed' }
Write-Output 'AWS_GUARDDUTY_MAGIKA_IDP_DISPATCH_DEPLOY_PASS'
````

### FILE: `aws_guardduty_magika_idp_dispatch/verify_contract.ps1`
```yaml
block_id: "AWS-GUARDDUTY-MAGIKA-IDP-DISPATCH-GATE:verify-contract-ps1:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-guardduty-magika-idp-dispatch-gate/verify_contract.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "b289df34d051bdeca070667f6bb9dcf734ee259b23232d15e6a917b9185dcd65"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param([string] $PythonExecutable = 'python')

$ErrorActionPreference = 'Stop'
$root = $PSScriptRoot
$expected = @(
  'AWS_MIT_NO_ATTRIBUTION_LICENSE.txt','GOOGLE_APACHE_LICENSE.txt','README.md','build.ps1','deploy.ps1',
  'handler.py','provider-profile.template.json','requirements-runtime.lock','template.yaml','test_handler.py','verify_contract.ps1'
)
$observed = @(Get-ChildItem -LiteralPath $root -File | Select-Object -ExpandProperty Name | Sort-Object)
if (@(Compare-Object $expected $observed).Count -ne 0) { throw 'component file allowlist mismatch' }

$handler = Get-Content -Raw -LiteralPath (Join-Path $root 'handler.py')
foreach ($token in @(
  'elite-aws-guardduty-release-receipt/v1','elite-aws-magika-idp-dispatch-receipt/v1','VersionId=event["version_id"]',
  'VersionId=released["version_id"]','CopySourceIfMatch=released["etag"]','IfNoneMatch="*"','ChecksumAlgorithm="SHA256"',
  'MAGIKA_VERSION = "magika 1.1.0 standard_v3_3"','automatic_business_persistence_authorized": False',
  'IDEMPOTENCY_TABLE','batchItemFailures','message_id_sha256'
)) { if (-not $handler.Contains($token)) { throw "handler contract token missing: $token" } }
foreach ($forbidden in @('logger.info(f"Processing event','print(event','automatic_business_persistence_authorized": True')) {
  if ($handler.Contains($forbidden)) { throw "forbidden handler token found: $forbidden" }
}

$template = Get-Content -Raw -LiteralPath (Join-Path $root 'template.yaml')
foreach ($token in @(
  'ObjectLockEnabled: true','Mode: COMPLIANCE','DeletionProtectionEnabled: true','PointInTimeRecoveryEnabled: true',
  'VersioningConfiguration: {Status: Enabled}','ReportBatchItemFailures','VisibilityTimeout: 1800','maxReceiveCount: 5',
  'Runtime: python3.14','Architectures: [x86_64]','Canary10Percent5Minutes','03a56971f9e121666d702db611fd08f1c8bae5a3c31538f1b45b3ef99e6df84d'
)) { if (-not $template.Contains($token)) { throw "template contract token missing: $token" } }

$lock = Get-Content -Raw -LiteralPath (Join-Path $root 'requirements-runtime.lock')
foreach ($token in @('boto3==1.43.83','botocore==1.43.83','--hash=sha256:')) { if (-not $lock.Contains($token)) { throw "dependency lock token missing: $token" } }
if (@($lock -split "`n" | Where-Object { $_ -match '^[a-z]' }).Count -ne 7) { throw 'dependency lock must contain seven distributions' }

$profile = Get-Content -Raw -LiteralPath (Join-Path $root 'provider-profile.template.json') | ConvertFrom-Json -Depth 20
if ($profile.profile_version -ne 1 -or @($profile.allowed_magika_types).Count -ne 0) { throw 'profile must begin fail-closed' }
foreach ($proof in @('aws_idp_source_acquired_and_verified','aws_idp_stack_outputs_verified','guardduty_canaries_proven','idp_config_corpus_evaluated','downstream_business_idempotency_proven','approve_live_effects','approve_provider_costs')) {
  if ($profile.$proof -ne $false) { throw "profile proof must begin false: $proof" }
}

$parseErrors = $null
[Management.Automation.Language.Parser]::ParseFile((Join-Path $root 'build.ps1'), [ref]$null, [ref]$parseErrors) | Out-Null
if ($parseErrors.Count -ne 0) { throw "build.ps1 parse failed: $($parseErrors[0].Message)" }
[Management.Automation.Language.Parser]::ParseFile((Join-Path $root 'deploy.ps1'), [ref]$null, [ref]$parseErrors) | Out-Null
if ($parseErrors.Count -ne 0) { throw "deploy.ps1 parse failed: $($parseErrors[0].Message)" }

$priorPythonPath = $env:PYTHONPATH
$priorBytecode = $env:PYTHONDONTWRITEBYTECODE
try {
  $env:PYTHONPATH = $root
  $env:PYTHONDONTWRITEBYTECODE = '1'
  & $PythonExecutable -B -m py_compile (Join-Path $root 'handler.py') (Join-Path $root 'test_handler.py')
  if ($LASTEXITCODE -ne 0) { throw 'Python compile failed' }
  Push-Location $root
  try { & $PythonExecutable -B -m unittest -v test_handler.py; if ($LASTEXITCODE -ne 0) { throw 'dispatch tests failed' } }
  finally { Pop-Location }
} finally {
  $env:PYTHONPATH = $priorPythonPath
  $env:PYTHONDONTWRITEBYTECODE = $priorBytecode
}

Write-Output 'AWS_GUARDDUTY_MAGIKA_IDP_DISPATCH_VERIFY_PASS tests=14 storage_authorized=0 exact_version=1 magika=1'
````

### FILE: `aws_guardduty_magika_idp_dispatch/GOOGLE_APACHE_LICENSE.txt`
```yaml
block_id: "AWS-GUARDDUTY-MAGIKA-IDP-DISPATCH-GATE:google-apache-license-txt:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/google/magika/blob/5e2f437fb7b7452368c8c1fa9354858f5487a5c4/LICENSE; adapted only by adding the final LF required by the Markdown materialization contract"
license: "Apache-2.0"
sha256: "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30"
variables: []
secrets_allowed: false
```
````text

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### FILE: `aws_guardduty_magika_idp_dispatch/AWS_MIT_NO_ATTRIBUTION_LICENSE.txt`
```yaml
block_id: "AWS-GUARDDUTY-MAGIKA-IDP-DISPATCH-GATE:aws-mit-no-attribution-license-txt:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/blob/1b5fd74454e593de233342a02ee911af8ee38359/LICENSE"
license: "MIT-0"
sha256: "5025c7bbedbc8da868b6dce2ab689225b3f33c43b3f03368b0a331809b6ada26"
variables: []
secrets_allowed: false
```
````text
MIT No Attribution

Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
````

## 6. Configuration surface

Completar una copia de `provider-profile.template.json`: cuenta/región/stack/tenant, buckets y claves KMS de recibo/release/input, stack y `config-version` exactos de AWS IDP, allowlist Magika por label/MIME/extensión/score, límites, owners y siete pruebas/aprobaciones inicialmente `false`. No hay secretos ni sustituciones dentro de los bloques; credenciales se obtienen por el provider chain de AWS.

## 7. Dependency bill

Runtime: Python 3.14 x86-64; siete wheels Boto3/Botocore 1.43.83 fijados por hash; binario oficial Google Magika CLI 1.1.0 Linux x86-64, asset SHA-256 `6b4c1010c84d1f4f06205ccef4597f1690bcd7744f46d841eee26426bc100485`, binario SHA-256 `03a56971f9e121666d702db611fd08f1c8bae5a3c31538f1b45b3ef99e6df84d`; AWS SAM/CLI para validar/desplegar y cfn-lint 1.55.1 usado en reconstrucción. AWS IDP v0.6.5 y Magika se adquieren por el core oficial, no se duplican como fuente local.

## 8. Apply order

1. Materializar adquisición oficial y verificar lock/asset/sidecar/attestation. 2. Materializar y probar GuardDuty release 0.1.x. 3. Materializar este pack y ejecutar `verify_contract.ps1`. 4. Construir en directorio vacío con el asset oficial Magika. 5. Completar provider profile y evidencia target. 6. Validar SAM/cfn-lint y desplegar sólo con aprobaciones explícitas. 7. Probar canaries, duplicados, DLQ/redrive, alertas, rollback y retención. 8. Conectar la salida a evaluación/revisión/persistencia idempotente, nunca directamente al modelo empresarial.

## 9. Verification

PASS requiere manifiesto↔bloques/hash exactos, compilación Python, 14 tests de contrato, allowlist cerrada, tokens de exact-version/conditional-copy/Object-Lock/lease/partial-batch, lock de siete dependencias, build offline del asset Magika verificado y template cfn-lint limpio. Live readiness además exige identidades AWS, aislamiento, configuración/corpus IDP, canaries GuardDuty, prueba downstream de idempotencia, owners y aprobaciones de efectos/costo. El PASS offline no es un veredicto de malware ni precisión documental productiva.

## 10. Reconstruction evidence

Reconstruido el 2026-08-28 desde AWS IDP commit `1b5fd74454e593de233342a02ee911af8ee38359`/tag `v0.6.5`, Google Magika commit `5e2f437fb7b7452368c8c1fa9354858f5487a5c4`/release `cli/v1.1.0`, asset/sidecar/API digest/attestation coincidentes. El árbol staging pasó 14 tests, build hash-locked y cfn-lint sin errores; la evidencia V91 registra comandos, hashes, fallos aprendidos y límites no resueltos. Este pack no afirma deploy productivo ni exactitud universal.
