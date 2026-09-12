# AWS SES Immutable Email Receiver

## 1. Metadata

```yaml
pack_id: "AWS-SES-IMMUTABLE-EMAIL-RECEIVER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa recepción SES→S3 versionado/KMS→SQS/DLQ→Lambda, retiene el MIME exacto con S3 Object Lock, extrae por el core MIME probado y deja adjuntos hash-named y manifest versionado en cuarentena sin autorizar persistencia empresarial."
stacks: ["AWS SES", "Amazon S3", "S3 Object Lock", "AWS KMS", "Amazon SQS", "AWS Lambda Python 3.14", "Amazon DynamoDB", "AWS SAM", "Boto3 1.43.83"]
compatible_with: ["SECURE-EMAIL-MIME-QUARANTINE-CORE 0.1.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x", "OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x"]
incompatible_with: ["bucket sin versionado", "raw MIME no retenido", "evento sin VersionId", "regla SES sin scan", "persistencia automática", "dependencias móviles", "deploy sin aprobación"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://docs.aws.amazon.com/ses/latest/dg/receiving-email-action-s3.html", "https://docs.aws.amazon.com/ses/latest/dg/receiving-email-concepts.html", "https://docs.aws.amazon.com/AmazonS3/latest/userguide/notification-how-to-event-types-and-destinations.html", "https://docs.aws.amazon.com/AmazonS3/latest/userguide/notification-content-structure.html", "https://docs.aws.amazon.com/lambda/latest/dg/services-sqs-errorhandling.html", "https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-configure-lambda-function-trigger.html", "https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-lock-configure.html", "https://docs.aws.amazon.com/lambda/latest/dg/python-package.html", "https://github.com/boto/boto3"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use después de que el proyecto pruebe dominio/MX, identidad SES, región receptora, cuenta, IAM, costos, tenant, clases MIME y responsables. Crea la infraestructura y el worker para recibir un recipient/prefix ligado a un tenant, retener el original, extraer adjuntos y dejarlos en cuarentena. Rechazar si se pretende usar headers From/To como autoridad, saltar seguridad/revisión, desplegar sin live approval o tratar el PASS offline como prueba de la cuenta destino.

## 3. Architecture contract

AWS publica los servicios, SDK y contratos; los nueve archivos son integración local `AUTHORED`, no un sample atribuido a AWS. SES escribe primero en un landing versionado SSE-KMS. S3 entrega eventos al menos una vez y sin orden garantizado; el worker exige bucket/key/version/sequencer/etag/tamaño, owner de cuenta, lease DynamoDB con owner fencing y manifest durable. Antes de extraer retiene bytes content-addressed con checksum, `If-None-Match`, KMS, Object Lock y verificación de retención sobre VersionId. Adjuntos base64 pasan el core separado, se escriben hash-named en cuarentena, y el manifest se confirma al final. Los reintentos verifican efectos parciales sin overwrite.

SES spam/virus deben ser PASS exactamente una vez. Sender headers no rutean tenant ni autorizan negocio. Logs no incluyen PII. Queue/DLQ cifradas, batch failures parciales, visibility 6× timeout, PITR/TTL, KMS rotation, least privilege, alarms, retain policies y doble aprobación de deploy están materializados. `automatic_storage_authorized=false` es inmutable. Rollback desactiva/cambia el rule set y el stack retiene evidencia; Object Lock y deletion protection requieren procedimiento explícito, no borrado automático.

## 4. Exact file manifest

```text
CREATE aws_ses_immutable_email_receiver/README.md
CREATE aws_ses_immutable_email_receiver/handler.py
CREATE aws_ses_immutable_email_receiver/test_handler.py
CREATE aws_ses_immutable_email_receiver/template.yaml
CREATE aws_ses_immutable_email_receiver/requirements.txt
CREATE aws_ses_immutable_email_receiver/provider-profile.template.json
CREATE aws_ses_immutable_email_receiver/build.ps1
CREATE aws_ses_immutable_email_receiver/deploy.ps1
CREATE aws_ses_immutable_email_receiver/verify_contract.ps1
```

## 5. Materialization blocks

### FILE: `aws_ses_immutable_email_receiver/README.md`
```yaml
block_id: "AWS-SES-IMMUTABLE-EMAIL-RECEIVER:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-ses-immutable-email-receiver/README.md"
license: "LicenseRef-Workspace-Owner"
sha256: "58941a53bb7809b8787ca2b09e41d52e693dea5a4607df99a65b948f5248b349"
variables: []
secrets_allowed: false
```
````text
# AWS SES immutable email receiver

This component is project-owned `AUTHORED` integration code over public AWS services and the pinned official Boto3 SDK. It is not an AWS sample and it does not misattribute the receiver to Amazon. It composes with `SECURE-EMAIL-MIME-QUARANTINE-CORE`; the build copies that already verified extractor into the Lambda artifact.

The deployed path is:

`SES receipt rule -> versioned SSE-KMS landing S3 -> version-bearing S3 event -> encrypted SQS/DLQ -> owner-fenced DynamoDB lease -> exact landing version -> checksum-addressed S3 Object Lock retained original -> bounded MIME extraction -> hash-named versioned KMS quarantine -> manifest committed last`.

AWS documents that S3 event notifications are at-least-once and unordered. The handler therefore binds identity to bucket/key/version/sequencer, uses a durable manifest plus DynamoDB owner fencing, writes with `If-None-Match: *`, verifies Object Lock retention on the created version, and returns per-message SQS failures. A partial remote effect is never treated as success; replay verifies existing hashes/versions and continues.

SES scanning must be enabled. Exactly one `X-SES-Virus-Verdict: PASS` and one `X-SES-Spam-Verdict: PASS` are required by default. Sender headers never select a tenant or authorize business persistence. The receipt rule and landing prefix bind one configured tenant. Authentication results are recorded only as booleans/hash unless the project explicitly requires at least one pass.

The SES action deliberately omits `KmsKeyArn`: AWS documents that this option performs S3-client-side encryption and requires its encryption client for decryption. Instead, the landing bucket applies server-side KMS encryption with a customer-managed key. Retained raw MIME uses Object Lock and explicit per-version retention; quarantine remains versioned and KMS-encrypted. All buckets, queues, the table, log group and KMS key are retained on stack deletion.

`requirements.txt` pins the seven-wheel Boto3 1.43.83 graph with exact SHA-256 hashes. The 2026-08-28 OSV query returned no records for those exact versions. This dated result is evidence, not a promise about future advisories; the library dependency-update and vulnerability-monitoring contracts still apply.

Run `pwsh ./verify_contract.ps1 -ComposedRoot <materialized-profile-root>`. Build with `pwsh ./build.ps1 -ComposedRoot <root> -OutputRoot <new-dir> -InstallDependencies`. Deployment additionally requires a completed copy of `provider-profile.template.json`, AWS CLI, SAM CLI, a verified SES domain, the exact regional MX record, both approval flags, cost/security/replay owners and live AWS authority. `deploy.ps1` builds, deploys, activates the new rule set and verifies it is active.

This pack ends at security quarantine with `automatic_storage_authorized=false`. Magika/ClamAV/YARA, document classification/extraction, field evidence, validation and human review remain subsequent composed gates. Offline PASS does not claim live delivery, target accuracy, cost, quota, DNS propagation, disaster recovery or production authorization; those require a canary email, DLQ/replay drill, retained-version inspection and target evidence.
````

### FILE: `aws_ses_immutable_email_receiver/handler.py`
```yaml
block_id: "AWS-SES-IMMUTABLE-EMAIL-RECEIVER:handler:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-ses-immutable-email-receiver/handler.py over official Boto3 public API"
license: "LicenseRef-Workspace-Owner"
sha256: "e87642aacfb16d23bcb3b760f44df5b9395cb30cb46358fa3618d4c0f9e60859"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import base64
import hashlib
import json
import os
import re
import tempfile
import time
import uuid
from dataclasses import dataclass
from datetime import datetime, timedelta, timezone
from email import policy
from email.parser import BytesParser
from pathlib import Path
from typing import Any
from urllib.parse import unquote_plus

from extract_email_attachments import Limits, extract


class ReceiverRejected(RuntimeError):
    def __init__(self, code: str):
        super().__init__(code)
        self.code = code


TENANT_RE = re.compile(r"^[a-z0-9][a-z0-9_-]{0,62}$")
ACCOUNT_RE = re.compile(r"^[0-9]{12}$")
KEY_RE = re.compile(r"^[A-Za-z0-9._/+=-]{1,1024}$")
HEX_RE = re.compile(r"^[0-9A-Fa-f]+$")


def _sha256(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def _canonical_json(value: Any) -> bytes:
    return (json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(",", ":")) + "\n").encode("utf-8")


def _checksum_b64(data: bytes) -> str:
    return base64.b64encode(hashlib.sha256(data).digest()).decode("ascii")


def _env_bool(name: str, default: bool) -> bool:
    value = os.environ.get(name, str(default)).strip().lower()
    if value not in {"true", "false"}:
        raise ReceiverRejected(f"INVALID_{name}")
    return value == "true"


@dataclass(frozen=True)
class Config:
    tenant_id: str
    expected_account_id: str
    landing_bucket: str
    landing_prefix: str
    retained_bucket: str
    quarantine_bucket: str
    kms_key_arn: str
    idempotency_table: str
    allowed_content_types: frozenset[str]
    receipt_prefix: str = "receipts"
    retain_days: int = 365
    dedupe_days: int = 400
    lease_seconds: int = 240
    max_raw_bytes: int = 30 * 1024 * 1024
    require_spam_pass: bool = True
    require_virus_pass: bool = True
    require_authentication_pass: bool = False
    retention_mode: str = "COMPLIANCE"

    def __post_init__(self) -> None:
        if not TENANT_RE.fullmatch(self.tenant_id):
            raise ReceiverRejected("INVALID_TENANT_ID")
        if not ACCOUNT_RE.fullmatch(self.expected_account_id):
            raise ReceiverRejected("INVALID_EXPECTED_ACCOUNT_ID")
        for value, code in ((self.landing_prefix, "LANDING_PREFIX"), (self.receipt_prefix, "RECEIPT_PREFIX")):
            if not KEY_RE.fullmatch(value) or value.startswith("/") or ".." in value.split("/"):
                raise ReceiverRejected(f"INVALID_{code}")
        if not self.landing_prefix.endswith("/"):
            raise ReceiverRejected("INVALID_LANDING_PREFIX")
        for value, code in (
            (self.landing_bucket, "LANDING_BUCKET"),
            (self.retained_bucket, "RETAINED_BUCKET"),
            (self.quarantine_bucket, "QUARANTINE_BUCKET"),
            (self.kms_key_arn, "KMS_KEY_ARN"),
            (self.idempotency_table, "IDEMPOTENCY_TABLE"),
        ):
            if not value or len(value) > 2048:
                raise ReceiverRejected(f"INVALID_{code}")
        if not self.allowed_content_types:
            raise ReceiverRejected("EMPTY_ALLOWED_CONTENT_TYPES")
        if self.retain_days < 1 or self.dedupe_days < self.retain_days or not 30 <= self.lease_seconds <= 900:
            raise ReceiverRejected("INVALID_RETENTION_OR_LEASE")
        if not 1 <= self.max_raw_bytes <= 30 * 1024 * 1024:
            raise ReceiverRejected("INVALID_MAX_RAW_BYTES")
        if self.retention_mode not in {"GOVERNANCE", "COMPLIANCE"}:
            raise ReceiverRejected("INVALID_RETENTION_MODE")

    @classmethod
    def from_env(cls) -> "Config":
        allowed = frozenset(x.strip().lower() for x in os.environ.get("ALLOWED_CONTENT_TYPES", "").split(",") if x.strip())
        return cls(
            tenant_id=os.environ.get("TENANT_ID", ""),
            expected_account_id=os.environ.get("EXPECTED_ACCOUNT_ID", ""),
            landing_bucket=os.environ.get("LANDING_BUCKET", ""),
            landing_prefix=os.environ.get("LANDING_PREFIX", ""),
            retained_bucket=os.environ.get("RETAINED_BUCKET", ""),
            quarantine_bucket=os.environ.get("QUARANTINE_BUCKET", ""),
            kms_key_arn=os.environ.get("KMS_KEY_ARN", ""),
            idempotency_table=os.environ.get("IDEMPOTENCY_TABLE", ""),
            allowed_content_types=allowed,
            receipt_prefix=os.environ.get("RECEIPT_PREFIX", "receipts"),
            retain_days=int(os.environ.get("RETAIN_DAYS", "365")),
            dedupe_days=int(os.environ.get("DEDUPE_DAYS", "400")),
            lease_seconds=int(os.environ.get("LEASE_SECONDS", "240")),
            max_raw_bytes=int(os.environ.get("MAX_RAW_BYTES", str(30 * 1024 * 1024))),
            require_spam_pass=_env_bool("REQUIRE_SPAM_PASS", True),
            require_virus_pass=_env_bool("REQUIRE_VIRUS_PASS", True),
            require_authentication_pass=_env_bool("REQUIRE_AUTHENTICATION_PASS", False),
            retention_mode=os.environ.get("RETENTION_MODE", "COMPLIANCE"),
        )


@dataclass(frozen=True)
class AwsClients:
    s3: Any
    dynamodb: Any


@dataclass(frozen=True)
class SourceObject:
    bucket: str
    key: str
    version_id: str
    sequencer: str
    etag: str
    size: int
    event_id: str


def _aws_code(exc: Exception) -> str:
    response = getattr(exc, "response", None)
    if isinstance(response, dict):
        error = response.get("Error")
        if isinstance(error, dict) and isinstance(error.get("Code"), str):
            return error["Code"]
    return exc.__class__.__name__


def _not_found(exc: Exception) -> bool:
    return _aws_code(exc) in {"404", "NoSuchKey", "NotFound", "ResourceNotFoundException"}


def _conditional(exc: Exception) -> bool:
    return _aws_code(exc) in {"412", "PreconditionFailed", "ConditionalCheckFailedException"}


def parse_s3_record(sqs_record: dict[str, Any], config: Config) -> SourceObject:
    message_id = sqs_record.get("messageId")
    body = sqs_record.get("body")
    if not isinstance(message_id, str) or not message_id or not isinstance(body, str):
        raise ReceiverRejected("INVALID_SQS_RECORD")
    try:
        payload = json.loads(body)
    except json.JSONDecodeError as exc:
        raise ReceiverRejected("INVALID_S3_EVENT_JSON") from exc
    records = payload.get("Records") if isinstance(payload, dict) else None
    if not isinstance(records, list) or len(records) != 1 or not isinstance(records[0], dict):
        raise ReceiverRejected("INVALID_S3_EVENT_CARDINALITY")
    event = records[0]
    if not str(event.get("eventVersion", "")).startswith("2.") or event.get("eventSource") != "aws:s3" or not str(event.get("eventName", "")).startswith("ObjectCreated:"):
        raise ReceiverRejected("INVALID_S3_EVENT_TYPE")
    s3 = event.get("s3")
    if not isinstance(s3, dict) or not isinstance(s3.get("bucket"), dict) or not isinstance(s3.get("object"), dict):
        raise ReceiverRejected("INVALID_S3_EVENT_SHAPE")
    bucket = s3["bucket"].get("name")
    obj = s3["object"]
    encoded_key, version_id, sequencer, etag, size = obj.get("key"), obj.get("versionId"), obj.get("sequencer"), obj.get("eTag"), obj.get("size")
    if bucket != config.landing_bucket or not all(isinstance(v, str) and v for v in (encoded_key, version_id, sequencer, etag)) or not isinstance(size, int):
        raise ReceiverRejected("S3_EVENT_AUTHORITY_MISMATCH")
    key = unquote_plus(encoded_key)
    if not KEY_RE.fullmatch(key) or not key.startswith(config.landing_prefix) or not HEX_RE.fullmatch(sequencer):
        raise ReceiverRejected("S3_EVENT_KEY_OR_SEQUENCE_REJECTED")
    if size < 1 or size > config.max_raw_bytes:
        raise ReceiverRejected("S3_EVENT_SIZE_REJECTED")
    event_material = f"{bucket}\n{key}\n{version_id}\n{sequencer}".encode("utf-8")
    return SourceObject(bucket, key, version_id, sequencer.lower(), etag.strip('"'), size, _sha256(event_material))


def _manifest_key(source: SourceObject, config: Config) -> str:
    return f"{config.receipt_prefix}/{config.tenant_id}/events/{source.event_id}/manifest.json"


def _already_complete(s3: Any, source: SourceObject, config: Config) -> bool:
    try:
        head = s3.head_object(Bucket=config.quarantine_bucket, Key=_manifest_key(source, config), ExpectedBucketOwner=config.expected_account_id)
    except Exception as exc:
        if _not_found(exc):
            return False
        raise
    metadata = head.get("Metadata") or {}
    if metadata.get("event-id") != source.event_id or head.get("ServerSideEncryption") != "aws:kms":
        raise ReceiverRejected("EXISTING_MANIFEST_AUTHORITY_MISMATCH")
    return True


def _acquire(dynamodb: Any, source: SourceObject, config: Config, owner: str, now: int) -> str:
    key = {"pk": {"S": f"EMAIL#{source.event_id}"}}
    item = {
        **key,
        "state": {"S": "IN_PROGRESS"},
        "owner": {"S": owner},
        "lease_expires": {"N": str(now + config.lease_seconds)},
        "expires_at": {"N": str(now + config.dedupe_days * 86400)},
        "attempts": {"N": "1"},
    }
    try:
        dynamodb.put_item(TableName=config.idempotency_table, Item=item, ConditionExpression="attribute_not_exists(pk)")
        return "ACQUIRED"
    except Exception as exc:
        if not _conditional(exc):
            raise
    existing = dynamodb.get_item(TableName=config.idempotency_table, Key=key, ConsistentRead=True).get("Item") or {}
    if existing.get("state", {}).get("S") == "COMPLETE":
        return "COMPLETE"
    try:
        dynamodb.update_item(
            TableName=config.idempotency_table,
            Key=key,
            UpdateExpression="SET #state=:progress, #owner=:owner, lease_expires=:lease, expires_at=:expires ADD attempts :one",
            ConditionExpression="(#state=:failed) OR (#state=:progress AND lease_expires < :now)",
            ExpressionAttributeNames={"#state": "state", "#owner": "owner"},
            ExpressionAttributeValues={
                ":progress": {"S": "IN_PROGRESS"}, ":failed": {"S": "FAILED"}, ":owner": {"S": owner},
                ":lease": {"N": str(now + config.lease_seconds)}, ":expires": {"N": str(now + config.dedupe_days * 86400)},
                ":now": {"N": str(now)}, ":one": {"N": "1"},
            },
        )
        return "ACQUIRED"
    except Exception as exc:
        if _conditional(exc):
            raise ReceiverRejected("IDEMPOTENCY_LEASE_BUSY") from exc
        raise


def _complete(dynamodb: Any, source: SourceObject, config: Config, owner: str, manifest_key: str, now: int) -> None:
    dynamodb.update_item(
        TableName=config.idempotency_table,
        Key={"pk": {"S": f"EMAIL#{source.event_id}"}},
        UpdateExpression="SET #state=:complete, manifest_key=:manifest, completed_at=:now REMOVE lease_expires",
        ConditionExpression="#state=:progress AND #owner=:owner",
        ExpressionAttributeNames={"#state": "state", "#owner": "owner"},
        ExpressionAttributeValues={
            ":complete": {"S": "COMPLETE"}, ":progress": {"S": "IN_PROGRESS"}, ":owner": {"S": owner},
            ":manifest": {"S": manifest_key}, ":now": {"N": str(now)},
        },
    )


def _fail(dynamodb: Any, source: SourceObject, config: Config, owner: str, code: str, now: int) -> None:
    try:
        dynamodb.update_item(
            TableName=config.idempotency_table,
            Key={"pk": {"S": f"EMAIL#{source.event_id}"}},
            UpdateExpression="SET #state=:failed, error_code=:code, failed_at=:now, lease_expires=:zero",
            ConditionExpression="#state=:progress AND #owner=:owner",
            ExpressionAttributeNames={"#state": "state", "#owner": "owner"},
            ExpressionAttributeValues={
                ":failed": {"S": "FAILED"}, ":progress": {"S": "IN_PROGRESS"}, ":owner": {"S": owner},
                ":code": {"S": code[:128]}, ":now": {"N": str(now)}, ":zero": {"N": "0"},
            },
        )
    except Exception:
        pass


def _read_raw(s3: Any, source: SourceObject, config: Config) -> bytes:
    response = s3.get_object(Bucket=source.bucket, Key=source.key, VersionId=source.version_id, ExpectedBucketOwner=config.expected_account_id)
    if response.get("VersionId") != source.version_id or str(response.get("ETag", "")).strip('"') != source.etag:
        raise ReceiverRejected("LANDING_VERSION_OR_ETAG_MISMATCH")
    body = response.get("Body")
    if body is None:
        raise ReceiverRejected("LANDING_BODY_MISSING")
    raw = body.read(config.max_raw_bytes + 1)
    if not isinstance(raw, bytes) or len(raw) != source.size or len(raw) > config.max_raw_bytes:
        raise ReceiverRejected("LANDING_BODY_SIZE_MISMATCH")
    return raw


def _verify_ses_headers(raw: bytes, config: Config) -> dict[str, Any]:
    message = BytesParser(policy=policy.default).parsebytes(raw, headersonly=True)
    virus = message.get_all("X-SES-Virus-Verdict", [])
    spam = message.get_all("X-SES-Spam-Verdict", [])
    auth = [str(value) for value in message.get_all("Authentication-Results", [])]
    if config.require_virus_pass and ([str(v).strip().upper() for v in virus] != ["PASS"]):
        raise ReceiverRejected("SES_VIRUS_VERDICT_REJECTED")
    if config.require_spam_pass and ([str(v).strip().upper() for v in spam] != ["PASS"]):
        raise ReceiverRejected("SES_SPAM_VERDICT_REJECTED")
    auth_joined = ";".join(auth).lower()
    auth_pass = bool(re.search(r"(?:^|[;\s])(?:spf|dkim|dmarc)=pass(?:[;\s]|$)", auth_joined))
    if config.require_authentication_pass and not auth_pass:
        raise ReceiverRejected("SES_AUTHENTICATION_REJECTED")
    return {
        "virus_pass": [str(v).strip().upper() for v in virus] == ["PASS"],
        "spam_pass": [str(v).strip().upper() for v in spam] == ["PASS"],
        "authentication_pass_observed": auth_pass,
        "authentication_headers_sha256": _sha256("\n".join(auth).encode("utf-8")),
    }


def _prove_object_lock(s3: Any, config: Config) -> None:
    response = s3.get_object_lock_configuration(Bucket=config.retained_bucket, ExpectedBucketOwner=config.expected_account_id)
    if (response.get("ObjectLockConfiguration") or {}).get("ObjectLockEnabled") != "Enabled":
        raise ReceiverRejected("RETAINED_BUCKET_OBJECT_LOCK_NOT_PROVEN")


def _verify_retention(s3: Any, config: Config, key: str, version_id: str, expected_until: datetime | None) -> datetime:
    response = s3.get_object_retention(Bucket=config.retained_bucket, Key=key, VersionId=version_id, ExpectedBucketOwner=config.expected_account_id)
    retention = response.get("Retention") or {}
    observed = retention.get("RetainUntilDate")
    if retention.get("Mode") != config.retention_mode or not isinstance(observed, datetime):
        raise ReceiverRejected("RETAINED_OBJECT_RETENTION_MISMATCH")
    observed = observed.astimezone(timezone.utc).replace(microsecond=0)
    if expected_until is not None and observed != expected_until:
        raise ReceiverRejected("RETAINED_OBJECT_RETENTION_MISMATCH")
    if observed <= datetime.now(timezone.utc):
        raise ReceiverRejected("RETAINED_OBJECT_RETENTION_EXPIRED")
    return observed


def _put_retained(s3: Any, raw: bytes, config: Config, now: int) -> dict[str, Any]:
    digest = _sha256(raw)
    key = f"raw/{config.tenant_id}/{digest}.eml"
    retain_until = datetime.fromtimestamp(now, timezone.utc).replace(microsecond=0) + timedelta(days=config.retain_days)
    _prove_object_lock(s3, config)
    try:
        response = s3.put_object(
            Bucket=config.retained_bucket, Key=key, Body=raw, ContentLength=len(raw), ContentType="message/rfc822",
            ChecksumAlgorithm="SHA256", ChecksumSHA256=_checksum_b64(raw), IfNoneMatch="*",
            Metadata={"content-sha256": digest, "tenant-id": config.tenant_id}, ExpectedBucketOwner=config.expected_account_id,
            ServerSideEncryption="aws:kms", SSEKMSKeyId=config.kms_key_arn,
            ObjectLockMode=config.retention_mode, ObjectLockRetainUntilDate=retain_until,
        )
        version_id = response.get("VersionId")
        if not isinstance(version_id, str) or not version_id:
            raise ReceiverRejected("RETAINED_OBJECT_VERSION_MISSING")
        observed_until = _verify_retention(s3, config, key, version_id, retain_until)
    except Exception as exc:
        if not _conditional(exc):
            raise
        head = s3.head_object(Bucket=config.retained_bucket, Key=key, ExpectedBucketOwner=config.expected_account_id)
        metadata = head.get("Metadata") or {}
        if metadata.get("content-sha256") != digest or head.get("ContentLength") != len(raw) or head.get("ServerSideEncryption") != "aws:kms":
            raise ReceiverRejected("EXISTING_RETAINED_OBJECT_MISMATCH") from exc
        version_id = head.get("VersionId")
        if not isinstance(version_id, str) or not version_id:
            raise ReceiverRejected("EXISTING_RETAINED_VERSION_MISSING") from exc
        observed_until = _verify_retention(s3, config, key, version_id, None)
    return {"uri": f"s3://{config.retained_bucket}/{key}", "bucket": config.retained_bucket, "key": key, "version_id": version_id, "sha256": digest, "retain_until": observed_until.isoformat().replace("+00:00", "Z"), "retention_verified": True}


def _put_quarantine_object(s3: Any, config: Config, key: str, data: bytes, metadata: dict[str, str]) -> dict[str, Any]:
    digest = _sha256(data)
    full_metadata = {**metadata, "content-sha256": digest, "tenant-id": config.tenant_id, "security-decision": "PENDING"}
    try:
        response = s3.put_object(
            Bucket=config.quarantine_bucket, Key=key, Body=data, ContentLength=len(data), ContentType="application/octet-stream",
            ChecksumAlgorithm="SHA256", ChecksumSHA256=_checksum_b64(data), IfNoneMatch="*", Metadata=full_metadata,
            ExpectedBucketOwner=config.expected_account_id, ServerSideEncryption="aws:kms", SSEKMSKeyId=config.kms_key_arn,
            Tagging="security-decision=PENDING&automatic-storage-authorized=false",
        )
        version_id = response.get("VersionId")
    except Exception as exc:
        if not _conditional(exc):
            raise
        head = s3.head_object(Bucket=config.quarantine_bucket, Key=key, ExpectedBucketOwner=config.expected_account_id)
        observed = head.get("Metadata") or {}
        if observed.get("content-sha256") != digest or head.get("ContentLength") != len(data) or head.get("ServerSideEncryption") != "aws:kms":
            raise ReceiverRejected("EXISTING_QUARANTINE_OBJECT_MISMATCH") from exc
        version_id = head.get("VersionId")
    if not isinstance(version_id, str) or not version_id:
        raise ReceiverRejected("QUARANTINE_VERSION_MISSING")
    return {"bucket": config.quarantine_bucket, "key": key, "version_id": version_id, "sha256": digest, "bytes": len(data), "security_decision": "PENDING"}


def _put_manifest(s3: Any, source: SourceObject, config: Config, manifest: dict[str, Any]) -> dict[str, str]:
    key = _manifest_key(source, config)
    data = _canonical_json(manifest)
    try:
        response = s3.put_object(
            Bucket=config.quarantine_bucket, Key=key, Body=data, ContentLength=len(data), ContentType="application/json",
            ChecksumAlgorithm="SHA256", ChecksumSHA256=_checksum_b64(data), IfNoneMatch="*",
            Metadata={"event-id": source.event_id, "manifest-sha256": _sha256(data), "tenant-id": config.tenant_id},
            ExpectedBucketOwner=config.expected_account_id, ServerSideEncryption="aws:kms", SSEKMSKeyId=config.kms_key_arn,
        )
        version_id = response.get("VersionId")
    except Exception as exc:
        if not _conditional(exc):
            raise
        head = s3.head_object(Bucket=config.quarantine_bucket, Key=key, ExpectedBucketOwner=config.expected_account_id)
        if (head.get("Metadata") or {}).get("manifest-sha256") != _sha256(data):
            raise ReceiverRejected("EXISTING_MANIFEST_CONTENT_MISMATCH") from exc
        version_id = head.get("VersionId")
    if not isinstance(version_id, str) or not version_id:
        raise ReceiverRejected("MANIFEST_VERSION_MISSING")
    return {"bucket": config.quarantine_bucket, "key": key, "version_id": version_id, "sha256": _sha256(data)}


def process_record(sqs_record: dict[str, Any], clients: AwsClients, config: Config, now: int | None = None) -> dict[str, Any]:
    now = int(time.time()) if now is None else now
    source = parse_s3_record(sqs_record, config)
    if _already_complete(clients.s3, source, config):
        return {"state": "ALREADY_COMPLETE", "event_id": source.event_id}
    owner = str(uuid.uuid4())
    acquired = _acquire(clients.dynamodb, source, config, owner, now)
    if acquired == "COMPLETE":
        return {"state": "ALREADY_COMPLETE", "event_id": source.event_id}
    try:
        raw = _read_raw(clients.s3, source, config)
        ses_evidence = _verify_ses_headers(raw, config)
        retained = _put_retained(clients.s3, raw, config, now)
        source_receipt = {
            "receipt_version": 1, "decision": "ADMIT", "tenant_id": config.tenant_id, "raw_sha256": retained["sha256"],
            "source_object": {"bucket": source.bucket, "key": source.key, "version_id": source.version_id, "sequencer": source.sequencer},
            "retained_original": {"uri": retained["uri"], "version_id": retained["version_id"]},
        }
        with tempfile.TemporaryDirectory(prefix="email-quarantine-") as temp:
            output = Path(temp) / "attachments"
            core_manifest = extract(raw, source_receipt, output, set(config.allowed_content_types), Limits(max_raw_bytes=config.max_raw_bytes))
            uploaded: list[dict[str, Any]] = []
            for attachment in core_manifest["attachments"]:
                data = (output / attachment["quarantine_name"]).read_bytes()
                key = f"quarantine/{config.tenant_id}/{retained['sha256']}/{attachment['quarantine_name']}"
                uploaded.append(_put_quarantine_object(clients.s3, config, key, data, {"raw-sha256": retained["sha256"], "part-index": str(attachment["part_index"])}))
        final_manifest = {
            "receipt_version": 1, "state": "QUARANTINED", "automatic_storage_authorized": False,
            "event_id": source.event_id, "tenant_id": config.tenant_id,
            "source": {"bucket": source.bucket, "key": source.key, "version_id": source.version_id, "sequencer": source.sequencer, "etag": source.etag},
            "ses_evidence": ses_evidence, "retained_original": retained, "core_manifest": core_manifest,
            "quarantine_objects": uploaded,
        }
        manifest_receipt = _put_manifest(clients.s3, source, config, final_manifest)
        _complete(clients.dynamodb, source, config, owner, manifest_receipt["key"], now)
        return {"state": "QUARANTINED", "event_id": source.event_id, "attachments": len(uploaded), "manifest": manifest_receipt}
    except Exception as exc:
        _fail(clients.dynamodb, source, config, owner, getattr(exc, "code", _aws_code(exc)), now)
        raise


_CLIENTS: AwsClients | None = None


def _runtime_clients() -> AwsClients:
    global _CLIENTS
    if _CLIENTS is None:
        import boto3
        _CLIENTS = AwsClients(s3=boto3.client("s3"), dynamodb=boto3.client("dynamodb"))
    return _CLIENTS


def lambda_handler(event: dict[str, Any], context: Any) -> dict[str, list[dict[str, str]]]:
    config = Config.from_env()
    records = event.get("Records") if isinstance(event, dict) else None
    if not isinstance(records, list):
        raise ReceiverRejected("INVALID_SQS_BATCH")
    failures: list[dict[str, str]] = []
    for record in records:
        message_id = record.get("messageId") if isinstance(record, dict) else None
        if not isinstance(message_id, str) or not message_id:
            raise ReceiverRejected("INVALID_SQS_MESSAGE_ID")
        try:
            result = process_record(record, _runtime_clients(), config)
            print(json.dumps({"level": "INFO", "state": result["state"], "event_id": result["event_id"]}, sort_keys=True))
        except Exception as exc:
            code = getattr(exc, "code", _aws_code(exc))
            print(json.dumps({"level": "ERROR", "state": "RETRY", "message_id_sha256": _sha256(message_id.encode()), "error_code": code}, sort_keys=True))
            failures.append({"itemIdentifier": message_id})
    return {"batchItemFailures": failures}
````

### FILE: `aws_ses_immutable_email_receiver/test_handler.py`
```yaml
block_id: "AWS-SES-IMMUTABLE-EMAIL-RECEIVER:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-ses-immutable-email-receiver/test_handler.py"
license: "LicenseRef-Workspace-Owner"
sha256: "1c42d9766421e0ad543442afe660812dfa58b67ffd3df1dd534b9407b805113f"
variables: []
secrets_allowed: false
```
````python
import base64
import io
import json
import os
import tempfile
import unittest
from datetime import timezone
from email.message import EmailMessage
from unittest.mock import patch
from urllib.parse import quote_plus

import handler


class AwsError(RuntimeError):
    def __init__(self, code):
        super().__init__(code)
        self.response = {"Error": {"Code": code}}


class FakeS3:
    def __init__(self):
        self.objects = {}
        self.retentions = {}
        self.put_log = []
        self.version = 0
        self.object_lock_enabled = True
        self.fail_once_contains = None

    def _next_version(self):
        self.version += 1
        return f"version-{self.version}"

    def add_landing(self, bucket, key, data):
        version = self._next_version()
        etag = handler._sha256(data)[:32]
        self.objects[(bucket, key)] = {
            "Body": data, "VersionId": version, "ETag": etag, "ContentLength": len(data),
            "Metadata": {}, "ServerSideEncryption": "aws:kms", "SSEKMSKeyId": "kms",
        }
        return version, etag

    def head_object(self, Bucket, Key, **kwargs):
        value = self.objects.get((Bucket, Key))
        if value is None:
            raise AwsError("404")
        return {k: v for k, v in value.items() if k != "Body"}

    def get_object(self, Bucket, Key, VersionId, **kwargs):
        value = self.objects.get((Bucket, Key))
        if value is None:
            raise AwsError("NoSuchKey")
        if value["VersionId"] != VersionId:
            raise AwsError("NoSuchVersion")
        result = {k: v for k, v in value.items() if k != "Body"}
        result["Body"] = io.BytesIO(value["Body"])
        return result

    def get_object_lock_configuration(self, **kwargs):
        return {"ObjectLockConfiguration": {"ObjectLockEnabled": "Enabled" if self.object_lock_enabled else "Disabled"}}

    def get_object_retention(self, Bucket, Key, VersionId, **kwargs):
        value = self.retentions.get((Bucket, Key, VersionId))
        if value is None:
            return {"Retention": {}}
        return {"Retention": value}

    def put_object(self, Bucket, Key, Body, **kwargs):
        if self.fail_once_contains and self.fail_once_contains in Key:
            self.fail_once_contains = None
            raise AwsError("ServiceUnavailable")
        if kwargs.get("IfNoneMatch") == "*" and (Bucket, Key) in self.objects:
            raise AwsError("PreconditionFailed")
        data = Body if isinstance(Body, bytes) else Body.read()
        version = self._next_version()
        value = {
            "Body": data, "VersionId": version, "ETag": handler._sha256(data)[:32], "ContentLength": len(data),
            "Metadata": kwargs.get("Metadata", {}), "ServerSideEncryption": kwargs.get("ServerSideEncryption"),
            "SSEKMSKeyId": kwargs.get("SSEKMSKeyId"),
        }
        self.objects[(Bucket, Key)] = value
        if "ObjectLockMode" in kwargs:
            self.retentions[(Bucket, Key, version)] = {
                "Mode": kwargs["ObjectLockMode"], "RetainUntilDate": kwargs["ObjectLockRetainUntilDate"].astimezone(timezone.utc)
            }
        self.put_log.append((Bucket, Key, kwargs))
        return {"VersionId": version, "ETag": value["ETag"], "ChecksumSHA256": kwargs.get("ChecksumSHA256")}


class FakeDynamoDB:
    def __init__(self):
        self.items = {}

    @staticmethod
    def _pk(mapping):
        return mapping["pk"]["S"]

    def put_item(self, Item, **kwargs):
        pk = self._pk(Item)
        if pk in self.items:
            raise AwsError("ConditionalCheckFailedException")
        self.items[pk] = json.loads(json.dumps(Item))
        return {}

    def get_item(self, Key, **kwargs):
        return {"Item": json.loads(json.dumps(self.items.get(self._pk(Key), {})))}

    def update_item(self, Key, UpdateExpression, ExpressionAttributeValues, **kwargs):
        pk = self._pk(Key)
        item = self.items.get(pk)
        if item is None:
            raise AwsError("ConditionalCheckFailedException")
        owner = ExpressionAttributeValues.get(":owner", {}).get("S")
        if ":complete" in ExpressionAttributeValues:
            if item["state"]["S"] != "IN_PROGRESS" or item["owner"]["S"] != owner:
                raise AwsError("ConditionalCheckFailedException")
            item["state"] = {"S": "COMPLETE"}
            item["manifest_key"] = ExpressionAttributeValues[":manifest"]
            item["completed_at"] = ExpressionAttributeValues[":now"]
            item.pop("lease_expires", None)
        elif ":failed" in ExpressionAttributeValues and ":progress" in ExpressionAttributeValues and "ADD attempts" in UpdateExpression:
            state = item["state"]["S"]
            now = int(ExpressionAttributeValues[":now"]["N"])
            lease = int(item.get("lease_expires", {"N": "0"})["N"])
            if not (state == "FAILED" or (state == "IN_PROGRESS" and lease < now)):
                raise AwsError("ConditionalCheckFailedException")
            item["state"] = {"S": "IN_PROGRESS"}
            item["owner"] = {"S": owner}
            item["lease_expires"] = ExpressionAttributeValues[":lease"]
            item["expires_at"] = ExpressionAttributeValues[":expires"]
            item["attempts"] = {"N": str(int(item["attempts"]["N"]) + 1)}
        elif ":failed" in ExpressionAttributeValues:
            if item["state"]["S"] != "IN_PROGRESS" or item["owner"]["S"] != owner:
                raise AwsError("ConditionalCheckFailedException")
            item["state"] = {"S": "FAILED"}
            item["error_code"] = ExpressionAttributeValues[":code"]
            item["failed_at"] = ExpressionAttributeValues[":now"]
            item["lease_expires"] = {"N": "0"}
        else:
            raise AssertionError(UpdateExpression)
        return {}


def mime_bytes(virus="PASS", spam="PASS", payload=b"invoice-data"):
    message = EmailMessage()
    message["From"] = "untrusted@example.net"
    message["To"] = "intake@example.com"
    message["Subject"] = "document"
    message["X-SES-Virus-Verdict"] = virus
    message["X-SES-Spam-Verdict"] = spam
    message["Authentication-Results"] = "example.com; dkim=pass header.i=example.net"
    message.set_content("attached")
    message.add_attachment(payload, maintype="application", subtype="pdf", filename="../../invoice.pdf", cte="base64")
    return message.as_bytes()


class ReceiverTests(unittest.TestCase):
    def setUp(self):
        self.s3 = FakeS3()
        self.ddb = FakeDynamoDB()
        self.clients = handler.AwsClients(self.s3, self.ddb)
        self.config = handler.Config(
            tenant_id="tenant-a", expected_account_id="123456789012", landing_bucket="landing",
            landing_prefix="incoming/tenant-a/", retained_bucket="retained", quarantine_bucket="quarantine",
            kms_key_arn="arn:aws:kms:us-east-1:123456789012:key/test", idempotency_table="idempotency",
            allowed_content_types=frozenset({"application/pdf"}), retain_days=365, dedupe_days=400,
        )

    def make_record(self, data=None, message_id="message-1", key="incoming/tenant-a/message-1"):
        data = mime_bytes() if data is None else data
        version, etag = self.s3.add_landing("landing", key, data)
        event = {"Records": [{
            "eventVersion": "2.1", "eventSource": "aws:s3", "eventName": "ObjectCreated:Put",
            "s3": {"bucket": {"name": "landing"}, "object": {"key": quote_plus(key), "size": len(data), "eTag": etag, "versionId": version, "sequencer": "00AF"}},
        }]}
        return {"messageId": message_id, "body": json.dumps(event)}

    def test_success_retains_extracts_and_commits_manifest_last(self):
        record = self.make_record()
        result = handler.process_record(record, self.clients, self.config, now=1_800_000_000)
        self.assertEqual("QUARANTINED", result["state"])
        self.assertEqual(1, result["attachments"])
        keys = [key for _, key, _ in self.s3.put_log]
        self.assertTrue(keys[0].startswith("raw/tenant-a/"))
        self.assertTrue(keys[1].startswith("quarantine/tenant-a/"))
        self.assertTrue(keys[-1].endswith("/manifest.json"))
        manifest = json.loads(self.s3.objects[("quarantine", keys[-1])]["Body"])
        self.assertFalse(manifest["automatic_storage_authorized"])
        self.assertTrue(manifest["retained_original"]["retention_verified"])
        self.assertEqual("PENDING", manifest["quarantine_objects"][0]["security_decision"])
        self.assertNotIn("invoice.pdf", json.dumps(manifest))

    def test_duplicate_uses_durable_manifest_without_remote_writes(self):
        record = self.make_record()
        first = handler.process_record(record, self.clients, self.config, now=1_800_000_000)
        writes = len(self.s3.put_log)
        second = handler.process_record(record, self.clients, self.config, now=1_800_000_001)
        self.assertEqual("ALREADY_COMPLETE", second["state"])
        self.assertEqual(writes, len(self.s3.put_log))
        self.assertEqual(first["event_id"], second["event_id"])

    def test_partial_effect_replays_without_overwriting_retained_original(self):
        record = self.make_record()
        self.s3.fail_once_contains = "quarantine/tenant-a/"
        with self.assertRaises(AwsError):
            handler.process_record(record, self.clients, self.config, now=1_800_000_000)
        raw_keys = [key for bucket, key, _ in self.s3.put_log if bucket == "retained"]
        self.assertEqual(1, len(raw_keys))
        result = handler.process_record(record, self.clients, self.config, now=1_800_000_001)
        self.assertEqual("QUARANTINED", result["state"])
        self.assertEqual(1, len([key for bucket, key, _ in self.s3.put_log if bucket == "retained"]))

    def test_virus_rejection_preserves_landing_and_blocks_all_derived_writes(self):
        record = self.make_record(mime_bytes(virus="FAIL"))
        with self.assertRaisesRegex(handler.ReceiverRejected, "SES_VIRUS_VERDICT_REJECTED"):
            handler.process_record(record, self.clients, self.config, now=1_800_000_000)
        self.assertEqual([], self.s3.put_log)
        self.assertEqual("FAILED", next(iter(self.ddb.items.values()))["state"]["S"])

    def test_wrong_bucket_prefix_missing_version_and_oversize_are_rejected(self):
        record = self.make_record()
        event = json.loads(record["body"])
        for mutation in (
            lambda e: e["Records"][0]["s3"]["bucket"].update(name="wrong"),
            lambda e: e["Records"][0]["s3"]["object"].update(key="other%2Fkey"),
            lambda e: e["Records"][0]["s3"]["object"].update(versionId=""),
            lambda e: e["Records"][0]["s3"]["object"].update(size=self.config.max_raw_bytes + 1),
        ):
            candidate = json.loads(json.dumps(event))
            mutation(candidate)
            with self.assertRaises(handler.ReceiverRejected):
                handler.parse_s3_record({"messageId": "x", "body": json.dumps(candidate)}, self.config)

    def test_unproven_object_lock_blocks_retained_write(self):
        record = self.make_record()
        self.s3.object_lock_enabled = False
        with self.assertRaisesRegex(handler.ReceiverRejected, "OBJECT_LOCK_NOT_PROVEN"):
            handler.process_record(record, self.clients, self.config, now=1_800_000_000)
        self.assertEqual([], self.s3.put_log)

    def test_owner_fencing_rejects_stale_completion(self):
        record = self.make_record()
        source = handler.parse_s3_record(record, self.config)
        self.assertEqual("ACQUIRED", handler._acquire(self.ddb, source, self.config, "owner-1", 100))
        with self.assertRaises(AwsError):
            handler._complete(self.ddb, source, self.config, "owner-2", "manifest", 101)

    def test_busy_lease_retries_and_expired_lease_can_be_reclaimed(self):
        record = self.make_record()
        source = handler.parse_s3_record(record, self.config)
        handler._acquire(self.ddb, source, self.config, "owner-1", 100)
        with self.assertRaisesRegex(handler.ReceiverRejected, "IDEMPOTENCY_LEASE_BUSY"):
            handler._acquire(self.ddb, source, self.config, "owner-2", 101)
        self.assertEqual("ACQUIRED", handler._acquire(self.ddb, source, self.config, "owner-2", 100 + self.config.lease_seconds + 1))

    def test_existing_manifest_with_wrong_authority_fails_closed(self):
        record = self.make_record()
        source = handler.parse_s3_record(record, self.config)
        key = handler._manifest_key(source, self.config)
        self.s3.put_object(Bucket="quarantine", Key=key, Body=b"{}", Metadata={"event-id": "wrong"}, ServerSideEncryption="aws:kms")
        with self.assertRaisesRegex(handler.ReceiverRejected, "EXISTING_MANIFEST_AUTHORITY_MISMATCH"):
            handler.process_record(record, self.clients, self.config, now=1_800_000_000)

    def test_batch_response_retries_only_failed_message(self):
        good = self.make_record(message_id="good", key="incoming/tenant-a/good")
        bad = {"messageId": "bad", "body": "not-json"}
        env = {
            "TENANT_ID": self.config.tenant_id, "EXPECTED_ACCOUNT_ID": self.config.expected_account_id,
            "LANDING_BUCKET": self.config.landing_bucket, "LANDING_PREFIX": self.config.landing_prefix,
            "RETAINED_BUCKET": self.config.retained_bucket, "QUARANTINE_BUCKET": self.config.quarantine_bucket,
            "KMS_KEY_ARN": self.config.kms_key_arn, "IDEMPOTENCY_TABLE": self.config.idempotency_table,
            "ALLOWED_CONTENT_TYPES": "application/pdf", "RETAIN_DAYS": "365", "DEDUPE_DAYS": "400",
        }
        with patch.dict(os.environ, env, clear=True), patch.object(handler, "_CLIENTS", self.clients):
            response = handler.lambda_handler({"Records": [good, bad]}, None)
        self.assertEqual([{"itemIdentifier": "bad"}], response["batchItemFailures"])

    def test_config_rejects_empty_allowlist_and_unsafe_prefix(self):
        values = dict(self.config.__dict__)
        values["allowed_content_types"] = frozenset()
        with self.assertRaises(handler.ReceiverRejected):
            handler.Config(**values)
        values = dict(self.config.__dict__)
        values["landing_prefix"] = "../incoming/"
        with self.assertRaises(handler.ReceiverRejected):
            handler.Config(**values)


if __name__ == "__main__":
    unittest.main()
````

### FILE: `aws_ses_immutable_email_receiver/template.yaml`
```yaml
block_id: "AWS-SES-IMMUTABLE-EMAIL-RECEIVER:iac:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-ses-immutable-email-receiver/template.yaml constrained by official AWS CloudFormation contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "9a74bf3abefa56e58ad67adc752181e15f95e7e0ce82dedaf8502d2ca6593381"
variables: []
secrets_allowed: false
```
````yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Description: Conditioned SES to immutable raw email and attachment quarantine receiver

Parameters:
  TenantId:
    Type: String
    AllowedPattern: '^[a-z0-9][a-z0-9_-]{0,62}$'
  ReceiveRecipient:
    Type: String
    MinLength: 3
  LandingBucketName:
    Type: String
    MinLength: 3
  RetainedBucketName:
    Type: String
    MinLength: 3
  QuarantineBucketName:
    Type: String
    MinLength: 3
  AllowedContentTypes:
    Type: String
    MinLength: 3
  RetainDays:
    Type: Number
    Default: 365
    MinValue: 1
  DedupeDays:
    Type: Number
    Default: 400
    MinValue: 1
  RetentionMode:
    Type: String
    Default: COMPLIANCE
    AllowedValues: [GOVERNANCE, COMPLIANCE]
  ReservedConcurrency:
    Type: Number
    Default: 5
    MinValue: 2
    MaxValue: 100

Resources:
  DataKey:
    Type: AWS::KMS::Key
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      Description: Email intake landing, retained raw, quarantine, SQS, DynamoDB and logs
      EnableKeyRotation: true
      PendingWindowInDays: 30
      KeyPolicy:
        Version: '2012-10-17'
        Statement:
          - Sid: RootAccountAuthority
            Effect: Allow
            Principal:
              AWS: !Sub arn:${AWS::Partition}:iam::${AWS::AccountId}:root
            Action: kms:*
            Resource: '*'

  DataKeyAlias:
    Type: AWS::KMS::Alias
    Properties:
      AliasName: !Sub alias/${AWS::StackName}-email-intake
      TargetKeyId: !Ref DataKey

  DeadLetterQueue:
    Type: AWS::SQS::Queue
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      KmsMasterKeyId: !GetAtt DataKey.Arn
      MessageRetentionPeriod: 1209600
      SqsManagedSseEnabled: false

  IntakeQueue:
    Type: AWS::SQS::Queue
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      KmsMasterKeyId: !GetAtt DataKey.Arn
      MaximumMessageSize: 262144
      MessageRetentionPeriod: 1209600
      VisibilityTimeout: 1800
      RedrivePolicy:
        deadLetterTargetArn: !GetAtt DeadLetterQueue.Arn
        maxReceiveCount: 5
      SqsManagedSseEnabled: false

  IntakeQueuePolicy:
    Type: AWS::SQS::QueuePolicy
    Properties:
      Queues: [!Ref IntakeQueue]
      PolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Sid: AllowOnlyLandingBucketEvents
            Effect: Allow
            Principal:
              Service: s3.amazonaws.com
            Action: sqs:SendMessage
            Resource: !GetAtt IntakeQueue.Arn
            Condition:
              ArnEquals:
                aws:SourceArn: !Sub arn:${AWS::Partition}:s3:::${LandingBucketName}
              StringEquals:
                aws:SourceAccount: !Ref AWS::AccountId
          - Sid: DenyNonTLS
            Effect: Deny
            Principal: '*'
            Action: sqs:*
            Resource: !GetAtt IntakeQueue.Arn
            Condition:
              Bool:
                aws:SecureTransport: false

  LandingBucket:
    Type: AWS::S3::Bucket
    DependsOn: IntakeQueuePolicy
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      BucketName: !Ref LandingBucketName
      BucketEncryption:
        ServerSideEncryptionConfiguration:
          - BucketKeyEnabled: true
            ServerSideEncryptionByDefault:
              SSEAlgorithm: aws:kms
              KMSMasterKeyID: !GetAtt DataKey.Arn
      NotificationConfiguration:
        QueueConfigurations:
          - Event: s3:ObjectCreated:*
            Queue: !GetAtt IntakeQueue.Arn
            Filter:
              S3Key:
                Rules:
                  - Name: prefix
                    Value: !Sub incoming/${TenantId}/
      OwnershipControls:
        Rules:
          - ObjectOwnership: BucketOwnerEnforced
      PublicAccessBlockConfiguration:
        BlockPublicAcls: true
        BlockPublicPolicy: true
        IgnorePublicAcls: true
        RestrictPublicBuckets: true
      VersioningConfiguration:
        Status: Enabled

  LandingBucketPolicy:
    Type: AWS::S3::BucketPolicy
    Properties:
      Bucket: !Ref LandingBucket
      PolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Sid: DenyNonTLS
            Effect: Deny
            Principal: '*'
            Action: s3:*
            Resource:
              - !GetAtt LandingBucket.Arn
              - !Sub ${LandingBucket.Arn}/*
            Condition:
              Bool:
                aws:SecureTransport: false

  RetainedBucket:
    Type: AWS::S3::Bucket
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      BucketName: !Ref RetainedBucketName
      BucketEncryption:
        ServerSideEncryptionConfiguration:
          - BucketKeyEnabled: true
            ServerSideEncryptionByDefault:
              SSEAlgorithm: aws:kms
              KMSMasterKeyID: !GetAtt DataKey.Arn
      ObjectLockEnabled: true
      ObjectLockConfiguration:
        ObjectLockEnabled: Enabled
        Rule:
          DefaultRetention:
            Days: !Ref RetainDays
            Mode: !Ref RetentionMode
      OwnershipControls:
        Rules:
          - ObjectOwnership: BucketOwnerEnforced
      PublicAccessBlockConfiguration:
        BlockPublicAcls: true
        BlockPublicPolicy: true
        IgnorePublicAcls: true
        RestrictPublicBuckets: true
      VersioningConfiguration:
        Status: Enabled

  RetainedBucketPolicy:
    Type: AWS::S3::BucketPolicy
    Properties:
      Bucket: !Ref RetainedBucket
      PolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Sid: DenyNonTLS
            Effect: Deny
            Principal: '*'
            Action: s3:*
            Resource:
              - !GetAtt RetainedBucket.Arn
              - !Sub ${RetainedBucket.Arn}/*
            Condition:
              Bool:
                aws:SecureTransport: false

  QuarantineBucket:
    Type: AWS::S3::Bucket
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      BucketName: !Ref QuarantineBucketName
      BucketEncryption:
        ServerSideEncryptionConfiguration:
          - BucketKeyEnabled: true
            ServerSideEncryptionByDefault:
              SSEAlgorithm: aws:kms
              KMSMasterKeyID: !GetAtt DataKey.Arn
      OwnershipControls:
        Rules:
          - ObjectOwnership: BucketOwnerEnforced
      PublicAccessBlockConfiguration:
        BlockPublicAcls: true
        BlockPublicPolicy: true
        IgnorePublicAcls: true
        RestrictPublicBuckets: true
      VersioningConfiguration:
        Status: Enabled

  QuarantineBucketPolicy:
    Type: AWS::S3::BucketPolicy
    Properties:
      Bucket: !Ref QuarantineBucket
      PolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Sid: DenyNonTLS
            Effect: Deny
            Principal: '*'
            Action: s3:*
            Resource:
              - !GetAtt QuarantineBucket.Arn
              - !Sub ${QuarantineBucket.Arn}/*
            Condition:
              Bool:
                aws:SecureTransport: false

  IdempotencyTable:
    Type: AWS::DynamoDB::Table
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      AttributeDefinitions:
        - AttributeName: pk
          AttributeType: S
      BillingMode: PAY_PER_REQUEST
      DeletionProtectionEnabled: true
      KeySchema:
        - AttributeName: pk
          KeyType: HASH
      PointInTimeRecoverySpecification:
        PointInTimeRecoveryEnabled: true
      SSESpecification:
        KMSMasterKeyId: !GetAtt DataKey.Arn
        SSEEnabled: true
        SSEType: KMS
      TimeToLiveSpecification:
        AttributeName: expires_at
        Enabled: true

  WorkerRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Effect: Allow
            Principal:
              Service: lambda.amazonaws.com
            Action: sts:AssumeRole
      Policies:
        - PolicyName: EmailIntakeWorkerLeastPrivilege
          PolicyDocument:
            Version: '2012-10-17'
            Statement:
              - Sid: ReadExactLandingVersions
                Effect: Allow
                Action: [s3:GetObject, s3:GetObjectVersion]
                Resource: !Sub ${LandingBucket.Arn}/incoming/${TenantId}/*
              - Sid: RetainedOriginals
                Effect: Allow
                Action: [s3:GetObject, s3:GetObjectRetention, s3:GetObjectLockConfiguration, s3:PutObject, s3:PutObjectRetention]
                Resource:
                  - !GetAtt RetainedBucket.Arn
                  - !Sub ${RetainedBucket.Arn}/raw/${TenantId}/*
              - Sid: QuarantineObjectsAndReceipts
                Effect: Allow
                Action: [s3:GetObject, s3:PutObject, s3:PutObjectTagging]
                Resource:
                  - !Sub ${QuarantineBucket.Arn}/quarantine/${TenantId}/*
                  - !Sub ${QuarantineBucket.Arn}/receipts/${TenantId}/*
              - Sid: IdempotencyOwnerFencing
                Effect: Allow
                Action: [dynamodb:GetItem, dynamodb:PutItem, dynamodb:UpdateItem]
                Resource: !GetAtt IdempotencyTable.Arn
              - Sid: ConsumeIntakeQueue
                Effect: Allow
                Action: [sqs:ChangeMessageVisibility, sqs:DeleteMessage, sqs:GetQueueAttributes, sqs:ReceiveMessage]
                Resource: !GetAtt IntakeQueue.Arn
              - Sid: UseDataKey
                Effect: Allow
                Action: [kms:Decrypt, kms:DescribeKey, kms:Encrypt, kms:GenerateDataKey]
                Resource: !GetAtt DataKey.Arn
              - Sid: WriteFunctionLogs
                Effect: Allow
                Action: [logs:CreateLogStream, logs:PutLogEvents]
                Resource: !Sub arn:${AWS::Partition}:logs:${AWS::Region}:${AWS::AccountId}:log-group:/aws/lambda/${AWS::StackName}-email-receiver:*

  ReceiverFunction:
    Type: AWS::Serverless::Function
    Properties:
      FunctionName: !Sub ${AWS::StackName}-email-receiver
      CodeUri: function/
      Handler: handler.lambda_handler
      Runtime: python3.14
      Architectures: [x86_64]
      Role: !GetAtt WorkerRole.Arn
      Timeout: 300
      MemorySize: 1024
      EphemeralStorage:
        Size: 2048
      ReservedConcurrentExecutions: !Ref ReservedConcurrency
      Tracing: Active
      Environment:
        Variables:
          TENANT_ID: !Ref TenantId
          EXPECTED_ACCOUNT_ID: !Ref AWS::AccountId
          LANDING_BUCKET: !Ref LandingBucket
          LANDING_PREFIX: !Sub incoming/${TenantId}/
          RETAINED_BUCKET: !Ref RetainedBucket
          QUARANTINE_BUCKET: !Ref QuarantineBucket
          KMS_KEY_ARN: !GetAtt DataKey.Arn
          IDEMPOTENCY_TABLE: !Ref IdempotencyTable
          ALLOWED_CONTENT_TYPES: !Ref AllowedContentTypes
          RETAIN_DAYS: !Ref RetainDays
          DEDUPE_DAYS: !Ref DedupeDays
          RETENTION_MODE: !Ref RetentionMode
          REQUIRE_SPAM_PASS: 'true'
          REQUIRE_VIRUS_PASS: 'true'
          REQUIRE_AUTHENTICATION_PASS: 'false'

  ReceiverLogGroup:
    Type: AWS::Logs::LogGroup
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      LogGroupName: !Sub /aws/lambda/${AWS::StackName}-email-receiver
      KmsKeyId: !GetAtt DataKey.Arn
      RetentionInDays: 30

  QueueEventSource:
    Type: AWS::Lambda::EventSourceMapping
    Properties:
      BatchSize: 10
      Enabled: true
      EventSourceArn: !GetAtt IntakeQueue.Arn
      FunctionName: !Ref ReceiverFunction
      FunctionResponseTypes: [ReportBatchItemFailures]
      MaximumBatchingWindowInSeconds: 5
      ScalingConfig:
        MaximumConcurrency: !Ref ReservedConcurrency

  SesDeliveryRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Effect: Allow
            Principal:
              Service: ses.amazonaws.com
            Action: sts:AssumeRole
            Condition:
              StringEquals:
                AWS:SourceAccount: !Ref AWS::AccountId
      Policies:
        - PolicyName: DeliverOnlyToLandingPrefix
          PolicyDocument:
            Version: '2012-10-17'
            Statement:
              - Effect: Allow
                Action: s3:PutObject
                Resource: !Sub ${LandingBucket.Arn}/incoming/${TenantId}/*
              - Effect: Allow
                Action: [kms:Encrypt, kms:GenerateDataKey]
                Resource: !GetAtt DataKey.Arn

  SesRuleSet:
    Type: AWS::SES::ReceiptRuleSet

  SesRule:
    Type: AWS::SES::ReceiptRule
    DependsOn: LandingBucketPolicy
    Properties:
      RuleSetName: !Ref SesRuleSet
      Rule:
        Enabled: true
        Name: !Sub ${AWS::StackName}-store-before-process
        Recipients: [!Ref ReceiveRecipient]
        ScanEnabled: true
        TlsPolicy: Require
        Actions:
          - S3Action:
              BucketName: !Ref LandingBucket
              IamRoleArn: !GetAtt SesDeliveryRole.Arn
              ObjectKeyPrefix: !Sub incoming/${TenantId}/

  DeadLetterAlarm:
    Type: AWS::CloudWatch::Alarm
    Properties:
      AlarmDescription: Email intake has records in its dead-letter queue
      ComparisonOperator: GreaterThanOrEqualToThreshold
      Dimensions:
        - Name: QueueName
          Value: !GetAtt DeadLetterQueue.QueueName
      EvaluationPeriods: 1
      MetricName: ApproximateNumberOfMessagesVisible
      Namespace: AWS/SQS
      Period: 60
      Statistic: Maximum
      Threshold: 1
      TreatMissingData: notBreaching

  OldestMessageAlarm:
    Type: AWS::CloudWatch::Alarm
    Properties:
      AlarmDescription: Email intake processing is delayed
      ComparisonOperator: GreaterThanThreshold
      Dimensions:
        - Name: QueueName
          Value: !GetAtt IntakeQueue.QueueName
      EvaluationPeriods: 2
      MetricName: ApproximateAgeOfOldestMessage
      Namespace: AWS/SQS
      Period: 300
      Statistic: Maximum
      Threshold: 900
      TreatMissingData: notBreaching

Outputs:
  RuleSetName:
    Value: !Ref SesRuleSet
  LandingBucket:
    Value: !Ref LandingBucket
  RetainedBucket:
    Value: !Ref RetainedBucket
  QuarantineBucket:
    Value: !Ref QuarantineBucket
  IntakeQueueUrl:
    Value: !Ref IntakeQueue
  DeadLetterQueueUrl:
    Value: !Ref DeadLetterQueue
  IdempotencyTable:
    Value: !Ref IdempotencyTable
````

### FILE: `aws_ses_immutable_email_receiver/requirements.txt`
```yaml
block_id: "AWS-SES-IMMUTABLE-EMAIL-RECEIVER:requirements:v1"
operation: CREATE
provenance: AUTHORED
source: "exact official PyPI Boto3 1.43.83 wheel graph verified 2026-08-28"
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

### FILE: `aws_ses_immutable_email_receiver/provider-profile.template.json`
```yaml
block_id: "AWS-SES-IMMUTABLE-EMAIL-RECEIVER:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-ses-immutable-email-receiver/provider-profile.template.json"
license: "LicenseRef-Workspace-Owner"
sha256: "9e1bd96e9a29e760826e5c806a92596184a72ede2c54d070ddf4fac2f9d40091"
variables: []
secrets_allowed: false
```
````json
{
  "profile_version": 1,
  "approve_live_effects": false,
  "aws_account_id": "",
  "aws_region": "",
  "stack_name": "",
  "tenant_id": "",
  "receive_recipient": "",
  "landing_bucket_name": "",
  "retained_bucket_name": "",
  "quarantine_bucket_name": "",
  "allowed_content_types": [],
  "retain_days": 365,
  "dedupe_days": 400,
  "retention_mode": "COMPLIANCE",
  "reserved_concurrency": 5,
  "dns_mx_proven": false,
  "ses_identity_proven": false,
  "cost_owner": "",
  "security_owner": "",
  "replay_owner": ""
}
````

### FILE: `aws_ses_immutable_email_receiver/build.ps1`
```yaml
block_id: "AWS-SES-IMMUTABLE-EMAIL-RECEIVER:build:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-ses-immutable-email-receiver/build.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "339692a57a15f8ebe864e1d4f5a6f17abbca9c0aaf017e79414ade4328650a97"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)][string] $ComposedRoot,
  [Parameter(Mandatory = $true)][string] $OutputRoot,
  [string] $PythonExecutable = 'python',
  [switch] $InstallDependencies
)

$ErrorActionPreference = 'Stop'
$root = (Resolve-Path -LiteralPath $ComposedRoot).Path
$receiver = Join-Path $root 'aws_ses_immutable_email_receiver'
$mime = Join-Path $root 'secure_email_mime_quarantine_core'
foreach ($path in @($receiver, $mime)) {
  if (-not (Test-Path -LiteralPath $path -PathType Container)) { throw "required composed directory missing: $path" }
}
$output = [IO.Path]::GetFullPath($OutputRoot)
if (Test-Path -LiteralPath $output) { throw "output must not exist: $output" }
$parent = Split-Path -Parent $output
if (-not (Test-Path -LiteralPath $parent -PathType Container)) { throw "output parent missing: $parent" }
$temp = Join-Path $parent ('.' + [IO.Path]::GetFileName($output) + '.' + [Guid]::NewGuid().ToString('N'))
$functionDir = Join-Path $temp 'function'

try {
  New-Item -ItemType Directory -Path $functionDir | Out-Null
  Copy-Item -LiteralPath (Join-Path $receiver 'handler.py') -Destination $functionDir
  Copy-Item -LiteralPath (Join-Path $receiver 'requirements.txt') -Destination $functionDir
  Copy-Item -LiteralPath (Join-Path $mime 'extract_email_attachments.py') -Destination $functionDir
  Copy-Item -LiteralPath (Join-Path $receiver 'template.yaml') -Destination $temp

  $env:PYTHONDONTWRITEBYTECODE = '1'
  $previousPythonPath = $env:PYTHONPATH
  try {
    $env:PYTHONPATH = "$receiver$([IO.Path]::PathSeparator)$mime"
    & $PythonExecutable (Join-Path $receiver 'test_handler.py') -v
    if ($LASTEXITCODE -ne 0) { throw 'receiver unit tests failed' }
  } finally {
    $env:PYTHONPATH = $previousPythonPath
  }

  if ($InstallDependencies) {
    & $PythonExecutable -m pip install --disable-pip-version-check --require-hashes --only-binary=:all: --no-compile --target $functionDir -r (Join-Path $functionDir 'requirements.txt')
    if ($LASTEXITCODE -ne 0) { throw 'hash-locked dependency installation failed' }
    & $PythonExecutable -c "import sys; sys.path.insert(0, r'$functionDir'); import boto3, botocore; assert boto3.__version__ == '1.43.83'; assert botocore.__version__ == '1.43.83'"
    if ($LASTEXITCODE -ne 0) { throw 'installed dependency import/version check failed' }
  }

  $receipt = @{
    receipt_version = 1
    dependencies_installed = [bool]$InstallDependencies
    files = @(Get-ChildItem -LiteralPath $temp -File -Recurse | ForEach-Object {
      @{
        path = [IO.Path]::GetRelativePath($temp, $_.FullName).Replace('\','/')
        bytes = $_.Length
        sha256 = (Get-FileHash -LiteralPath $_.FullName -Algorithm SHA256).Hash.ToLowerInvariant()
      }
    } | Sort-Object path)
  }
  $receipt | ConvertTo-Json -Depth 5 | Set-Content -LiteralPath (Join-Path $temp 'BUILD_RECEIPT.json') -Encoding utf8NoBOM
  Move-Item -LiteralPath $temp -Destination $output
  Write-Output "AWS_SES_RECEIVER_BUILD_PASS output=$output files=$($receipt.files.Count) dependencies=$([bool]$InstallDependencies)"
} catch {
  if (Test-Path -LiteralPath $temp) { Remove-Item -LiteralPath $temp -Recurse -Force }
  throw
}
````

### FILE: `aws_ses_immutable_email_receiver/deploy.ps1`
```yaml
block_id: "AWS-SES-IMMUTABLE-EMAIL-RECEIVER:deploy:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-ses-immutable-email-receiver/deploy.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "2e1893e326c719518007e2ad619fbdb4074794b735206b21c5e8490c0816ad95"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)][string] $ComposedRoot,
  [Parameter(Mandatory = $true)][string] $ProfilePath,
  [Parameter(Mandatory = $true)][string] $BuildRoot,
  [switch] $ApproveLiveEffects
)

$ErrorActionPreference = 'Stop'
$profileFile = (Resolve-Path -LiteralPath $ProfilePath).Path
$profile = Get-Content -LiteralPath $profileFile -Raw | ConvertFrom-Json
if (-not $ApproveLiveEffects -or $profile.approve_live_effects -ne $true) { throw 'live AWS effects require both -ApproveLiveEffects and approve_live_effects=true' }
foreach ($field in @('aws_account_id','aws_region','stack_name','tenant_id','receive_recipient','landing_bucket_name','retained_bucket_name','quarantine_bucket_name','cost_owner','security_owner','replay_owner')) {
  if ([string]::IsNullOrWhiteSpace([string]$profile.$field)) { throw "profile field is required: $field" }
}
if ([string]$profile.aws_account_id -notmatch '^[0-9]{12}$') { throw 'aws_account_id must be 12 digits' }
if ($profile.dns_mx_proven -ne $true -or $profile.ses_identity_proven -ne $true) { throw 'DNS MX and SES identity evidence must be acknowledged before deployment' }
$allowed = @($profile.allowed_content_types)
if ($allowed.Count -eq 0 -or @($allowed | Where-Object { $_ -notmatch '^[a-z0-9.+-]+/[a-z0-9.+-]+$' }).Count -gt 0) { throw 'allowed_content_types must be a non-empty MIME allowlist' }
if ((Get-Command aws -ErrorAction SilentlyContinue) -eq $null -or (Get-Command sam -ErrorAction SilentlyContinue) -eq $null) { throw 'AWS CLI and AWS SAM CLI are required' }

$observedAccount = (& aws sts get-caller-identity --query Account --output text --region $profile.aws_region).Trim()
if ($LASTEXITCODE -ne 0 -or $observedAccount -cne [string]$profile.aws_account_id) { throw "AWS account mismatch: expected=$($profile.aws_account_id) observed=$observedAccount" }
$domain = ([string]$profile.receive_recipient -split '@')[-1].TrimEnd('.').ToLowerInvariant()
if ($domain -notmatch '^[a-z0-9.-]+\.[a-z]{2,}$') { throw 'receive_recipient must contain a valid owned domain' }
$identityStatus = (& aws ses get-identity-verification-attributes --identities $domain --query "VerificationAttributes.'$domain'.VerificationStatus" --output text --region $profile.aws_region).Trim()
if ($LASTEXITCODE -ne 0 -or $identityStatus -cne 'Success') { throw "SES domain identity is not verified in target region: $domain ($identityStatus)" }
$mx = @(Resolve-DnsName -Name $domain -Type MX -ErrorAction Stop | Where-Object Type -eq 'MX')
$expectedMx = "inbound-smtp.$($profile.aws_region).amazonaws.com"
if (@($mx | Where-Object { ([string]$_.NameExchange).TrimEnd('.') -ceq $expectedMx }).Count -eq 0) { throw "DNS MX does not prove SES receiving endpoint: $expectedMx" }

$build = [IO.Path]::GetFullPath($BuildRoot)
if (Test-Path -LiteralPath $build) { throw "BuildRoot must not exist: $build" }
& (Join-Path $ComposedRoot 'aws_ses_immutable_email_receiver/build.ps1') -ComposedRoot $ComposedRoot -OutputRoot $build -InstallDependencies
if ($LASTEXITCODE -ne 0) { throw 'receiver build failed' }

$samBuild = Join-Path $build '.sam-build'
& sam build --template-file (Join-Path $build 'template.yaml') --build-dir $samBuild --no-cached
if ($LASTEXITCODE -ne 0) { throw 'sam build failed' }
$parameterOverrides = @(
  "TenantId=$($profile.tenant_id)", "ReceiveRecipient=$($profile.receive_recipient)",
  "LandingBucketName=$($profile.landing_bucket_name)", "RetainedBucketName=$($profile.retained_bucket_name)",
  "QuarantineBucketName=$($profile.quarantine_bucket_name)", "AllowedContentTypes=$($allowed -join ',')",
  "RetainDays=$($profile.retain_days)", "DedupeDays=$($profile.dedupe_days)",
  "RetentionMode=$($profile.retention_mode)", "ReservedConcurrency=$($profile.reserved_concurrency)"
)
& sam deploy --template-file (Join-Path $samBuild 'template.yaml') --stack-name $profile.stack_name --region $profile.aws_region --capabilities CAPABILITY_IAM --resolve-s3 --no-confirm-changeset --no-fail-on-empty-changeset --parameter-overrides $parameterOverrides
if ($LASTEXITCODE -ne 0) { throw 'sam deploy failed' }
$ruleSet = (& aws cloudformation describe-stacks --stack-name $profile.stack_name --region $profile.aws_region --query "Stacks[0].Outputs[?OutputKey=='RuleSetName'].OutputValue | [0]" --output text).Trim()
if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrWhiteSpace($ruleSet) -or $ruleSet -ceq 'None') { throw 'deployed receipt rule set output is missing' }
& aws ses set-active-receipt-rule-set --rule-set-name $ruleSet --region $profile.aws_region
if ($LASTEXITCODE -ne 0) { throw 'receipt rule set activation failed' }
$active = (& aws ses describe-active-receipt-rule-set --region $profile.aws_region --query Metadata.Name --output text).Trim()
if ($LASTEXITCODE -ne 0 -or $active -cne $ruleSet) { throw "active receipt rule set mismatch: expected=$ruleSet observed=$active" }
Write-Output "AWS_SES_RECEIVER_DEPLOY_PASS stack=$($profile.stack_name) account=$observedAccount region=$($profile.aws_region) active_rule_set=$active"
````

### FILE: `aws_ses_immutable_email_receiver/verify_contract.ps1`
```yaml
block_id: "AWS-SES-IMMUTABLE-EMAIL-RECEIVER:verify:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-ses-immutable-email-receiver/verify_contract.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "d802546d46b0b2ed4121c6b507e8fca61b380e1e9ad9a34e73f2dd3e09351a96"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)][string] $ComposedRoot,
  [string] $PythonExecutable = 'python'
)

$ErrorActionPreference = 'Stop'
$root = (Resolve-Path -LiteralPath $ComposedRoot).Path
$receiver = Join-Path $root 'aws_ses_immutable_email_receiver'
$mime = Join-Path $root 'secure_email_mime_quarantine_core'
$expected = @('README.md','handler.py','test_handler.py','template.yaml','requirements.txt','provider-profile.template.json','build.ps1','deploy.ps1','verify_contract.ps1')
foreach ($name in $expected) {
  if (-not (Test-Path -LiteralPath (Join-Path $receiver $name) -PathType Leaf)) { throw "missing receiver file: $name" }
}
if (-not (Test-Path -LiteralPath (Join-Path $mime 'extract_email_attachments.py') -PathType Leaf)) { throw 'composed MIME core is missing' }

$template = Get-Content -LiteralPath (Join-Path $receiver 'template.yaml') -Raw
foreach ($token in @(
  'AWS::SES::ReceiptRule','ScanEnabled: true','TlsPolicy: Require','AWS::S3::Bucket','ObjectLockEnabled: true',
  'Mode: !Ref RetentionMode','VersioningConfiguration:','SSEAlgorithm: aws:kms','AWS::SQS::Queue','maxReceiveCount: 5',
  'VisibilityTimeout: 1800','AWS::Lambda::EventSourceMapping','ReportBatchItemFailures','AWS::DynamoDB::Table',
  'PointInTimeRecoveryEnabled: true','DeletionProtectionEnabled: true','ReservedConcurrentExecutions:',
  'AWS::CloudWatch::Alarm','DeletionPolicy: Retain','UpdateReplacePolicy: Retain','aws:SecureTransport: false'
)) {
  if (-not $template.Contains($token, [StringComparison]::Ordinal)) { throw "template contract token missing: $token" }
}
if ($template -match '(?m)^\s+KmsKeyArn:' -or $template -match '(?m)^\s+DeletionPolicy:\s+Delete\s*$' -or $template -match '(?m)^\s+Action:\s+\x27?\*\x27?\s*$') {
  throw 'template contains forbidden SES client encryption, destructive deletion, or wildcard action'
}
if ([regex]::Matches($template, 'DeletionPolicy: Retain').Count -lt 8) { throw 'durable resources are not retained consistently' }
if (-not $template.Contains('Timeout: 300') -or -not $template.Contains('VisibilityTimeout: 1800')) { throw 'SQS visibility must remain at least six times Lambda timeout' }

$requirements = Get-Content -LiteralPath (Join-Path $receiver 'requirements.txt') -Raw
$locks = @{
  'boto3==1.43.83'='73a3564f737d4516625964eee709a498fa98ccee6aca929febad2b0b5fbeae1e'
  'botocore==1.43.83'='bf75a6cf587c22d968e43e79fe122c39f82deafbe9c3422bc5d3e80b6210fc98'
  'jmespath==1.1.0'='a5663118de4908c91729bea0acadca56526eb2698e83de10cd116ae0f4e97c64'
  'python-dateutil==2.9.0.post0'='a8b2bc7bffae282281c8140a97d3aa9c14da0b136dfe83f850eea9a5f7470427'
  's3transfer==0.19.2'='d8168eccca828cbb2cd573675333f3bddd254313a9c42494b84c76b539e8ba25'
  'six==1.17.0'='4721f391ed90541fddacab5acf947aa0d3dc7d27b2e1e8eda2be8970586c3274'
  'urllib3==2.7.0'='9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897'
}
foreach ($entry in $locks.GetEnumerator()) {
  if (-not $requirements.Contains($entry.Key, [StringComparison]::Ordinal) -or -not $requirements.Contains($entry.Value, [StringComparison]::Ordinal)) { throw "dependency lock drifted: $($entry.Key)" }
}
if ([regex]::Matches($requirements, '(?m)^[a-z]').Count -ne 7 -or [regex]::Matches($requirements, '--hash=sha256:').Count -ne 7) { throw 'dependency graph must contain exactly seven hash-locked packages' }

$profile = Get-Content -LiteralPath (Join-Path $receiver 'provider-profile.template.json') -Raw | ConvertFrom-Json
if ($profile.approve_live_effects -ne $false -or $profile.dns_mx_proven -ne $false -or $profile.ses_identity_proven -ne $false -or @($profile.allowed_content_types).Count -ne 0) { throw 'provider profile must start fail-closed' }
$deploy = Get-Content -LiteralPath (Join-Path $receiver 'deploy.ps1') -Raw
foreach ($token in @('-ApproveLiveEffects','approve_live_effects','get-caller-identity','get-identity-verification-attributes','Resolve-DnsName','set-active-receipt-rule-set','describe-active-receipt-rule-set')) {
  if (-not $deploy.Contains($token, [StringComparison]::Ordinal)) { throw "deployment gate missing: $token" }
}
$handler = Get-Content -LiteralPath (Join-Path $receiver 'handler.py') -Raw
foreach ($token in @('VersionId=source.version_id','ExpectedBucketOwner=config.expected_account_id','IfNoneMatch="*"','ChecksumSHA256','get_object_lock_configuration','get_object_retention','IDEMPOTENCY_LEASE_BUSY','automatic_storage_authorized','batchItemFailures','message_id_sha256')) {
  if (-not $handler.Contains($token, [StringComparison]::Ordinal)) { throw "handler contract token missing: $token" }
}

$env:PYTHONDONTWRITEBYTECODE = '1'
$previousPythonPath = $env:PYTHONPATH
try {
  $env:PYTHONPATH = "$receiver$([IO.Path]::PathSeparator)$mime"
  & $PythonExecutable -c "import ast,pathlib; [ast.parse(pathlib.Path(p).read_text(encoding='utf-8')) for p in [r'$(Join-Path $receiver 'handler.py')',r'$(Join-Path $receiver 'test_handler.py')']]"
  if ($LASTEXITCODE -ne 0) { throw 'Python parse failed' }
  & $PythonExecutable (Join-Path $receiver 'test_handler.py') -v
  if ($LASTEXITCODE -ne 0) { throw 'receiver tests failed' }
} finally {
  $env:PYTHONPATH = $previousPythonPath
}
Write-Output 'AWS_SES_IMMUTABLE_EMAIL_RECEIVER_PASS tests=11 dependencies=7'
````

## 6. Configuration surface

El profile empieza con efectos live, MX, SES identity y allowlist sin autoridad. Cuenta/región/stack/tenant/recipient/buckets/allowlist/retención/concurrencia y owners son obligatorios. `deploy.ps1` revalida cuenta STS, SES identity y MX regional, exige aprobación en profile más switch y recién entonces crea efectos. Secrets no se aceptan en el profile; AWS usa credential chain/SSO/roles externos.

## 7. Dependency bill

| Package/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---:|---|---|---|---|
| boto3/botocore | 1.43.83/1.43.83 + hashes | AWS APIs | Apache-2.0 | runtime | PyPI/GitHub boto |
| jmespath | 1.1.0 + hash | SDK | MIT | runtime | PyPI |
| python-dateutil | 2.9.0.post0 + hash | SDK | Apache-2.0/BSD | runtime | PyPI |
| s3transfer | 0.19.2 + hash | SDK | Apache-2.0 | runtime | PyPI/GitHub boto |
| six | 1.17.0 + hash | SDK | MIT | runtime | PyPI |
| urllib3 | 2.7.0 + hash | SDK HTTP | MIT | runtime | PyPI |
| CPython | 3.14 Lambda | worker | PSF-2.0 | runtime | AWS managed runtime |
| PowerShell | 7+ | gates/build/deploy | MIT | build | Microsoft |
| AWS SAM/CLI | target-installed | build/deploy live | Apache-2.0 | build | AWS |

Los siete wheels exactos fueron descargados e instalados con `--require-hashes --only-binary`; importaron Boto3/Botocore 1.43.83. Query OSV datada 2026-08-28: cero records. Actualización futura reabre hashes, SCA, tests, lint y target canary.

## 8. Apply order

Materializar perfil en carpeta vacía; completar readiness/proveedor; copiar y completar el profile sin secrets; ejecutar verifier; ejecutar build en destino nuevo; validar target; desplegar sólo con doble aprobación; activar/verificar rule set; enviar canary; inspeccionar landing/version/Object Lock/quarantine/manifest; inducir retry/DLQ/replay; habilitar recién después gates de seguridad/extracción. Workspace existente: componer en paths sin colisión. Rollback: cambiar/desactivar rule set y preservar todos los recursos/evidencia; nunca borrar Object Lock automáticamente.

## 9. Verification

`pwsh ./aws_ses_immutable_email_receiver/verify_contract.ps1 -ComposedRoot <root>` debe dar 11/11 y PASS. `build.ps1 -InstallDependencies` debe instalar sólo hashes y emitir `BUILD_RECEIPT.json`; el audit produjo 2.191 archivos de artifact y Boto3/Botocore import PASS. `cfn-lint 1.55.1 template.yaml` debe salir 0. Live: `deploy.ps1` debe probar cuenta/identity/MX, stack y rule set activo; canary/retry/DLQ/replay/restore/costos continúan siendo evidencia del target, no se simulan.

## 10. Reconstruction evidence

2026-08-28, Windows/PowerShell 7/CPython 3.14.4: 9/9 archivos staged; Python AST PASS; 11 tests PASS; core MIME 18 tests ya probado; requirements 7 wheels/16.324.178 bytes descargados con hashes, instalación `--require-hashes` e imports 1.43.83 PASS; build receipt 2.191 files; cfn-lint 1.55.1 exit 0. No AWS credentials/SAM CLI ni live effects se usaron. Fallos de runner Windows, entrypoint cfn-lint y mínimo de ScalingConfig se corrigieron y deben quedar en el ledger.
