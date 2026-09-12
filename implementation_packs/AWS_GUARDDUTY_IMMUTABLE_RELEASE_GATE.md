# AWS GuardDuty Immutable Release Gate

## 1. Metadata

```yaml
pack_id: "AWS-GUARDDUTY-IMMUTABLE-RELEASE-GATE"
pack_version: "0.1.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa GuardDuty Malware Protection for S3 sobre la cuarentena del receptor, valida evento/tag/versión exactos y sólo libera contenido limpio a S3 KMS+Object Lock con recibo inmutable; ningún resultado autoriza persistencia empresarial."
stacks: ["Amazon GuardDuty Malware Protection for S3", "Amazon EventBridge", "Amazon SQS", "AWS Lambda Python 3.14", "Amazon S3 Object Lock", "AWS KMS", "Amazon DynamoDB", "Amazon SNS", "AWS SAM", "Boto3 1.43.83"]
compatible_with: ["AWS-SES-IMMUTABLE-EMAIL-RECEIVER 0.1.x", "SECURE-EMAIL-MIME-QUARANTINE-CORE 0.1.x", "OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x"]
incompatible_with: ["S3 sin versionado", "lectura sin tag limpio", "evento sin VersionId", "plan o cuenta no fijados", "persistencia automática", "deploy sin costo aprobado", "sample AWS desplegado sin hardening"]
license_expression: "LicenseRef-Workspace-Owner AND MIT-0 AND Apache-2.0"
upstream_sources: ["https://github.com/aws-samples/guardduty-malware-protection/tree/fd9d2bdf46eb52bb336a7e7a837bd38ef73002e2", "https://docs.aws.amazon.com/guardduty/latest/ug/monitor-with-eventbridge-s3-malware-protection.html", "https://docs.aws.amazon.com/guardduty/latest/ug/how-malware-protection-for-s3-gdu-works.html", "https://docs.aws.amazon.com/guardduty/latest/ug/tag-based-access-s3-malware-protection.html", "https://aws.amazon.com/blogs/security/using-amazon-guardduty-malware-protection-to-scan-uploads-to-amazon-s3/", "https://github.com/boto/boto3"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use después del receptor email cuando el proyecto elige AWS y acepta el costo de GuardDuty. Requiere cuenta/región, bucket de cuarentena versionado SSE-KMS, tenant/prefix, responsables de costo/seguridad/replay, aislamiento de lectura y canaries. Rechazar si se pretende tratar `UNSUPPORTED`, `ACCESS_DENIED`, `FAILED`, un tag ausente o una versión distinta como contenido limpio; si se desea desplegar el sample upstream sin las correcciones declaradas; o si el proyecto exige ejecución gratuita.

## 3. Architecture contract

El repositorio oficial AWS aporta el patrón GuardDuty plan→EventBridge→SQS/DLQ→Lambda y respuesta parcial. `handler.py` y `template.yaml` son adaptaciones declaradas: fijan plan/cuenta/región/bucket/key/version/eTag/tag, contenido hash-addressed, KMS/Object Lock, owner fencing, recibo durable, logs redacted, retain policies y permisos condicionados al tag limpio. El notice MIT-0 se preserva normalizado; los scripts/tests son `AUTHORED` y Boto3 proviene del lock hashado del receptor.

Sólo `COMPLETED/NO_THREATS_FOUND` más tag oficial coincidente en la versión exacta puede producir `ADMIT_TO_EXTRACTION`. Malware no se lee ni copia. Estados no escaneados se conservan y alertan. La copia exacta usa `VersionId`, `CopySourceIfMatch` e `IfNoneMatch`; los efectos parciales se reconcilian y el recibo se confirma último. Release y receipts son retenidos. El worker nunca cambia `automatic_business_persistence_authorized=false`.

La cuenta target debe demostrar que otros principals tampoco leen cuarentena. El stack no reemplaza la bucket policy del receptor porque hacerlo destruiría ownership entre stacks; se exige evidencia de aislamiento o merge explícito del TBAC oficial que preserve TLS y replay. Despliegue, tags y scan generan costo. Rollback deshabilita la regla/consumidor y conserva plan, buckets, tabla, colas, claves y evidencia; la remoción del plan es un cambio separado autorizado.

## 4. Exact file manifest

```text
CREATE aws_guardduty_immutable_release_gate/README.md
CREATE aws_guardduty_immutable_release_gate/handler.py
CREATE aws_guardduty_immutable_release_gate/test_handler.py
CREATE aws_guardduty_immutable_release_gate/template.yaml
CREATE aws_guardduty_immutable_release_gate/provider-profile.template.json
CREATE aws_guardduty_immutable_release_gate/build.ps1
CREATE aws_guardduty_immutable_release_gate/deploy.ps1
CREATE aws_guardduty_immutable_release_gate/verify_contract.ps1
CREATE aws_guardduty_immutable_release_gate/AWS_MIT_NO_ATTRIBUTION_LICENSE.txt
```

## 5. Materialization blocks

### FILE: `aws_guardduty_immutable_release_gate/README.md`
```yaml
block_id: "AWS-GUARDDUTY-IMMUTABLE-RELEASE-GATE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-guardduty-immutable-release-gate/README.md"
license: "LicenseRef-Workspace-Owner"
sha256: "ce4bb286853fb5783b58c1d273f40dbcefadffbe04ec48397ad7ba7e65b34f75"
variables: []
secrets_allowed: false
```
````text
# AWS GuardDuty immutable quarantine release gate

This component is a declared adaptation of the AWS `aws-samples/guardduty-malware-protection` CloudFormation/SQS/Lambda pattern at commit `fd9d2bdf46eb52bb336a7e7a837bd38ef73002e2`, combined with current public GuardDuty event, tag and at-least-once contracts. It preserves the upstream MIT-0 notice. The upstream archive remains independently pinned as `aws-guardduty-malware-protection-s3` in the official acquisition core.

It does not deploy the AWS sample unchanged. The 2026-08-28 audit found that the official sample copies by key without binding `versionId`, logs object keys, uses deletable buckets/keys, and its exact CloudFormation template returns cfn-lint warnings. Those limitations are recorded rather than hidden. `handler.py` and `template.yaml` are `ADAPTED`; tests and operational scripts are project-owned `AUTHORED` glue. Nothing is falsely presented as unchanged Amazon code.

The materialized path is:

`versioned KMS quarantine -> GuardDuty Malware Protection plan with managed tag -> EventBridge exact plan/bucket event -> encrypted SQS/DLQ -> account/region/plan/schema/version/eTag/tag validation -> DynamoDB owner-fenced lease -> exact clean version -> content-addressed KMS/Object-Lock release -> immutable security receipt -> next document gate`.

Only `COMPLETED + NO_THREATS_FOUND` with the matching `GuardDutyMalwareScanStatus` tag on the exact S3 version can be copied. The Lambda role has source read permission conditioned on that clean tag. `THREATS_FOUND` is never read or copied; threat names are hashed in the receipt. `UNSUPPORTED`, `ACCESS_DENIED` and `FAILED` remain unscanned/rejected and alert security. Every receipt uses the explicit `elite-aws-guardduty-release-receipt/v1` schema. No status authorizes business persistence: every receipt fixes `automatic_business_persistence_authorized=false`.

Release and receipt buckets are versioned, KMS-encrypted, Object-Locked in compliance mode, TLS-only, public-blocked and retained. Copies use source `VersionId`, `CopySourceIfMatch`, destination `IfNoneMatch: *`, SHA-256 addressing, exact metadata and retention verification. Receipts are committed last. Durable receipts, leases, owner fencing and partial batch responses make duplicate delivery and partial remote effects replayable without overwrite.

The GuardDuty plan role and resource are adapted from the official AWS template and narrowed to one receiver tenant prefix. AWS tagging is enabled before intake. The stack deliberately does not overwrite the quarantine bucket policy owned by the receiver stack. The project readiness gate must separately prove that no other workload role can read quarantine; this component's own role is already tag-conditioned. An organization-wide TBAC merge may use the official AWS policy, but must preserve the receiver's TLS policy and replay permissions.

Run `pwsh ./verify_contract.ps1`, then build into a new directory with `pwsh ./build.ps1 -ComposedRoot <profile-root> -OutputRoot <empty-output> -InstallDependencies`. The build reuses the receiver's exact hash-locked Boto3 1.43.83 graph so the profile carries one dependency authority.

Live deployment requires a completed provider profile, AWS CLI, SAM CLI, exact account/region, existing versioned SSE-KMS quarantine bucket, source key, security alert owner, cost owner, replay owner and both approval switches. GuardDuty Malware Protection, EventBridge, SQS, Lambda, DynamoDB, S3 tagging/storage, KMS, SNS and logs can incur AWS charges. Offline PASS is not a live malware verdict. Before enabling downstream extraction, send clean/EICAR/unsupported canaries, confirm the SNS subscription, inspect exact-version tags and receipts, exercise retry/DLQ/redrive, and demonstrate restore/rollback without deleting retained evidence.

Official references:

- https://github.com/aws-samples/guardduty-malware-protection/tree/fd9d2bdf46eb52bb336a7e7a837bd38ef73002e2
- https://docs.aws.amazon.com/guardduty/latest/ug/monitor-with-eventbridge-s3-malware-protection.html
- https://docs.aws.amazon.com/guardduty/latest/ug/how-malware-protection-for-s3-gdu-works.html
- https://docs.aws.amazon.com/guardduty/latest/ug/tag-based-access-s3-malware-protection.html
- https://aws.amazon.com/blogs/security/using-amazon-guardduty-malware-protection-to-scan-uploads-to-amazon-s3/
````

### FILE: `aws_guardduty_immutable_release_gate/handler.py`
```yaml
block_id: "AWS-GUARDDUTY-IMMUTABLE-RELEASE-GATE:handler:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/aws-samples/guardduty-malware-protection/blob/fd9d2bdf46eb52bb336a7e7a837bd38ef73002e2/cfn/template.yaml#L357-L399; narrowed and hardened for exact-version immutable release"
license: "MIT-0"
sha256: "7405f64023ca2861cb1066770b143964103f6f1cdb0d9836315706fabd25c29b"
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
import time
import uuid
from dataclasses import dataclass
from datetime import datetime, timedelta, timezone
from typing import Any
from urllib.parse import unquote_plus


class SecurityRejected(RuntimeError):
    def __init__(self, code: str):
        super().__init__(code)
        self.code = code


ACCOUNT_RE = re.compile(r"^[0-9]{12}$")
TENANT_RE = re.compile(r"^[a-z0-9][a-z0-9_-]{0,62}$")
HEX64_RE = re.compile(r"^[0-9a-f]{64}$")
ETAG_RE = re.compile(r"^[0-9A-Fa-f-]{1,128}$")
SAFE_KEY_RE = re.compile(r"^[A-Za-z0-9._/+=-]{1,1024}$")
STATUS_MAP = {
    "NO_THREATS_FOUND": "COMPLETED",
    "THREATS_FOUND": "COMPLETED",
    "UNSUPPORTED": "SKIPPED",
    "ACCESS_DENIED": "SKIPPED",
    "FAILED": "FAILED",
}


def _sha256(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def _canonical_json(value: Any) -> bytes:
    return (json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(",", ":")) + "\n").encode("utf-8")


def _checksum_b64(data: bytes) -> str:
    return base64.b64encode(hashlib.sha256(data).digest()).decode("ascii")


def _is_not_found(exc: Exception) -> bool:
    response = getattr(exc, "response", {})
    code = str(response.get("Error", {}).get("Code", ""))
    return code in {"404", "NoSuchKey", "NotFound"}


def _is_precondition(exc: Exception) -> bool:
    response = getattr(exc, "response", {})
    code = str(response.get("Error", {}).get("Code", ""))
    return code in {"PreconditionFailed", "412", "ConditionalCheckFailedException"}


@dataclass(frozen=True)
class Config:
    expected_account_id: str
    expected_region: str
    expected_plan_arn: str
    tenant_id: str
    quarantine_bucket: str
    quarantine_prefix: str
    release_bucket: str
    receipt_bucket: str
    source_kms_key_arn: str
    release_kms_key_arn: str
    receipt_kms_key_arn: str
    idempotency_table: str
    alert_topic_arn: str
    retain_days: int = 365
    lease_seconds: int = 240
    dedupe_days: int = 400
    max_object_bytes: int = 50 * 1024 * 1024
    retention_mode: str = "COMPLIANCE"

    def __post_init__(self) -> None:
        if not ACCOUNT_RE.fullmatch(self.expected_account_id):
            raise SecurityRejected("INVALID_EXPECTED_ACCOUNT_ID")
        if not TENANT_RE.fullmatch(self.tenant_id):
            raise SecurityRejected("INVALID_TENANT_ID")
        if not re.fullmatch(r"^[a-z]{2}(?:-gov)?-[a-z]+-[0-9]$", self.expected_region):
            raise SecurityRejected("INVALID_EXPECTED_REGION")
        if not self.expected_plan_arn.startswith("arn:") or ":guardduty:" not in self.expected_plan_arn:
            raise SecurityRejected("INVALID_EXPECTED_PLAN_ARN")
        expected_prefix = f"quarantine/{self.tenant_id}/"
        if self.quarantine_prefix != expected_prefix:
            raise SecurityRejected("INVALID_QUARANTINE_PREFIX")
        for value, code in (
            (self.quarantine_bucket, "QUARANTINE_BUCKET"),
            (self.release_bucket, "RELEASE_BUCKET"),
            (self.receipt_bucket, "RECEIPT_BUCKET"),
            (self.source_kms_key_arn, "SOURCE_KMS_KEY_ARN"),
            (self.release_kms_key_arn, "RELEASE_KMS_KEY_ARN"),
            (self.receipt_kms_key_arn, "RECEIPT_KMS_KEY_ARN"),
            (self.idempotency_table, "IDEMPOTENCY_TABLE"),
            (self.alert_topic_arn, "ALERT_TOPIC_ARN"),
        ):
            if not value or len(value) > 2048:
                raise SecurityRejected(f"INVALID_{code}")
        if not 1 <= self.retain_days <= 36500 or self.dedupe_days < self.retain_days:
            raise SecurityRejected("INVALID_RETENTION")
        if not 30 <= self.lease_seconds <= 900 or not 1 <= self.max_object_bytes <= 5 * 1024 * 1024 * 1024:
            raise SecurityRejected("INVALID_LIMIT")
        if self.retention_mode not in {"GOVERNANCE", "COMPLIANCE"}:
            raise SecurityRejected("INVALID_RETENTION_MODE")

    @classmethod
    def from_env(cls) -> "Config":
        return cls(
            expected_account_id=os.environ.get("EXPECTED_ACCOUNT_ID", ""),
            expected_region=os.environ.get("EXPECTED_REGION", ""),
            expected_plan_arn=os.environ.get("EXPECTED_PLAN_ARN", ""),
            tenant_id=os.environ.get("TENANT_ID", ""),
            quarantine_bucket=os.environ.get("QUARANTINE_BUCKET", ""),
            quarantine_prefix=os.environ.get("QUARANTINE_PREFIX", ""),
            release_bucket=os.environ.get("RELEASE_BUCKET", ""),
            receipt_bucket=os.environ.get("RECEIPT_BUCKET", ""),
            source_kms_key_arn=os.environ.get("SOURCE_KMS_KEY_ARN", ""),
            release_kms_key_arn=os.environ.get("RELEASE_KMS_KEY_ARN", ""),
            receipt_kms_key_arn=os.environ.get("RECEIPT_KMS_KEY_ARN", ""),
            idempotency_table=os.environ.get("IDEMPOTENCY_TABLE", ""),
            alert_topic_arn=os.environ.get("ALERT_TOPIC_ARN", ""),
            retain_days=int(os.environ.get("RETAIN_DAYS", "365")),
            lease_seconds=int(os.environ.get("LEASE_SECONDS", "240")),
            dedupe_days=int(os.environ.get("DEDUPE_DAYS", "400")),
            max_object_bytes=int(os.environ.get("MAX_OBJECT_BYTES", str(50 * 1024 * 1024))),
            retention_mode=os.environ.get("RETENTION_MODE", "COMPLIANCE"),
        )


@dataclass(frozen=True)
class ScanEvent:
    event_id: str
    account_id: str
    region: str
    plan_arn: str
    bucket: str
    key: str
    version_id: str
    etag: str
    status: str
    scan_status: str
    status_reasons: tuple[str, ...]
    threat_name_hashes: tuple[str, ...]

    @property
    def scan_id(self) -> str:
        authority = "\n".join((self.account_id, self.region, self.plan_arn, self.bucket, self.key, self.version_id, self.etag))
        return _sha256(authority.encode("utf-8"))


@dataclass(frozen=True)
class AwsClients:
    s3: Any
    dynamodb: Any
    sns: Any


def _required_string(value: Any, code: str, maximum: int = 2048) -> str:
    if not isinstance(value, str) or not value or len(value) > maximum:
        raise SecurityRejected(code)
    return value


def parse_guardduty_event(event: Any, config: Config) -> ScanEvent:
    if not isinstance(event, dict):
        raise SecurityRejected("INVALID_EVENT")
    if event.get("version") != "0" or event.get("source") != "aws.guardduty":
        raise SecurityRejected("INVALID_EVENT_AUTHORITY")
    if event.get("detail-type") != "GuardDuty Malware Protection Object Scan Result":
        raise SecurityRejected("INVALID_DETAIL_TYPE")
    account = _required_string(event.get("account"), "INVALID_EVENT_ACCOUNT", 12)
    region = _required_string(event.get("region"), "INVALID_EVENT_REGION", 64)
    if account != config.expected_account_id or region != config.expected_region:
        raise SecurityRejected("EVENT_ACCOUNT_REGION_MISMATCH")
    resources = event.get("resources")
    if resources != [config.expected_plan_arn]:
        raise SecurityRejected("EVENT_PLAN_MISMATCH")
    event_id = _required_string(event.get("id"), "INVALID_EVENT_ID", 128)
    try:
        uuid.UUID(event_id)
    except ValueError as exc:
        raise SecurityRejected("INVALID_EVENT_ID") from exc
    detail = event.get("detail")
    if not isinstance(detail, dict) or detail.get("schemaVersion") != "1.0" or detail.get("resourceType") != "S3_OBJECT":
        raise SecurityRejected("INVALID_SCAN_SCHEMA")
    scan_status = _required_string(detail.get("scanStatus"), "INVALID_SCAN_STATUS", 32)
    object_details = detail.get("s3ObjectDetails")
    result_details = detail.get("scanResultDetails")
    if not isinstance(object_details, dict) or not isinstance(result_details, dict):
        raise SecurityRejected("INVALID_SCAN_DETAILS")
    status = _required_string(result_details.get("scanResultStatus"), "INVALID_RESULT_STATUS", 32)
    if STATUS_MAP.get(status) != scan_status:
        raise SecurityRejected("RESULT_STATUS_MISMATCH")
    bucket = _required_string(object_details.get("bucketName"), "INVALID_SOURCE_BUCKET", 63)
    if bucket != config.quarantine_bucket:
        raise SecurityRejected("SOURCE_BUCKET_MISMATCH")
    key = unquote_plus(_required_string(object_details.get("objectKey"), "INVALID_SOURCE_KEY", 2048))
    if not SAFE_KEY_RE.fullmatch(key) or not key.startswith(config.quarantine_prefix) or ".." in key.split("/"):
        raise SecurityRejected("SOURCE_KEY_MISMATCH")
    relative = key[len(config.quarantine_prefix):]
    if not re.fullmatch(r"[0-9a-f]{64}/part-[0-9]{4}-[0-9a-f]{64}\.bin", relative):
        raise SecurityRejected("SOURCE_KEY_NOT_RECEIVER_ATTACHMENT")
    version_id = _required_string(object_details.get("versionId"), "MISSING_SOURCE_VERSION", 1024)
    etag = _required_string(object_details.get("eTag"), "MISSING_SOURCE_ETAG", 128).strip('"')
    if not ETAG_RE.fullmatch(etag):
        raise SecurityRejected("INVALID_SOURCE_ETAG")
    if not isinstance(object_details.get("s3Throttled"), bool):
        raise SecurityRejected("INVALID_THROTTLE_FLAG")
    raw_reasons = result_details.get("statusReasons")
    if raw_reasons is None:
        reasons: tuple[str, ...] = ()
    elif isinstance(raw_reasons, list) and all(isinstance(x, str) and 0 < len(x) <= 128 for x in raw_reasons):
        reasons = tuple(sorted(set(raw_reasons)))
    else:
        raise SecurityRejected("INVALID_STATUS_REASONS")
    raw_threats = result_details.get("threats")
    threat_hashes: list[str] = []
    if status == "THREATS_FOUND":
        if not isinstance(raw_threats, list) or not raw_threats:
            raise SecurityRejected("MISSING_THREAT_EVIDENCE")
        for threat in raw_threats:
            name = threat.get("name") if isinstance(threat, dict) else None
            threat_hashes.append(_sha256(_required_string(name, "INVALID_THREAT_EVIDENCE", 512).encode("utf-8")))
    elif raw_threats not in (None, []):
        raise SecurityRejected("UNEXPECTED_THREAT_EVIDENCE")
    return ScanEvent(
        event_id=event_id, account_id=account, region=region, plan_arn=config.expected_plan_arn,
        bucket=bucket, key=key, version_id=version_id, etag=etag, status=status,
        scan_status=scan_status, status_reasons=reasons, threat_name_hashes=tuple(threat_hashes),
    )


def _receipt_key(scan: ScanEvent, config: Config) -> str:
    return f"security-receipts/{config.tenant_id}/{scan.scan_id}.json"


def _existing_receipt(s3: Any, scan: ScanEvent, config: Config) -> dict[str, str] | None:
    key = _receipt_key(scan, config)
    try:
        head = s3.head_object(Bucket=config.receipt_bucket, Key=key, ExpectedBucketOwner=config.expected_account_id)
    except Exception as exc:
        if _is_not_found(exc):
            return None
        raise
    metadata = head.get("Metadata") or {}
    if metadata.get("scan-id") != scan.scan_id or metadata.get("scan-status") != scan.status:
        raise SecurityRejected("EXISTING_RECEIPT_AUTHORITY_MISMATCH")
    version_id = _required_string(head.get("VersionId"), "EXISTING_RECEIPT_VERSION_MISSING", 1024)
    return {"bucket": config.receipt_bucket, "key": key, "version_id": version_id}


def _acquire_lease(dynamodb: Any, scan: ScanEvent, config: Config, owner: str, now: int) -> None:
    expires = now + config.lease_seconds
    ttl = now + config.dedupe_days * 86400
    try:
        dynamodb.put_item(
            TableName=config.idempotency_table,
            Item={"pk": {"S": scan.scan_id}, "state": {"S": "IN_PROGRESS"}, "owner": {"S": owner},
                  "lease_expires": {"N": str(expires)}, "expires_at": {"N": str(ttl)}},
            ConditionExpression="attribute_not_exists(pk) OR #state = :failed OR lease_expires < :now",
            ExpressionAttributeNames={"#state": "state"},
            ExpressionAttributeValues={":failed": {"S": "FAILED"}, ":now": {"N": str(now)}},
        )
    except Exception as exc:
        if _is_precondition(exc):
            raise SecurityRejected("IDEMPOTENCY_LEASE_BUSY") from exc
        raise


def _finish_lease(dynamodb: Any, scan: ScanEvent, config: Config, owner: str, state: str, receipt_key: str, now: int) -> None:
    try:
        dynamodb.update_item(
            TableName=config.idempotency_table, Key={"pk": {"S": scan.scan_id}},
            UpdateExpression="SET #state=:state, receipt_key=:receipt, completed_at=:now REMOVE lease_expires",
            ConditionExpression="#owner=:owner AND #state=:progress",
            ExpressionAttributeNames={"#owner": "owner", "#state": "state"},
            ExpressionAttributeValues={":owner": {"S": owner}, ":progress": {"S": "IN_PROGRESS"},
                                       ":state": {"S": state}, ":receipt": {"S": receipt_key}, ":now": {"N": str(now)}},
        )
    except Exception as exc:
        if _is_precondition(exc):
            raise SecurityRejected("IDEMPOTENCY_OWNER_FENCE") from exc
        raise


def _mark_failed(dynamodb: Any, scan: ScanEvent, config: Config, owner: str, code: str, now: int) -> None:
    try:
        dynamodb.update_item(
            TableName=config.idempotency_table, Key={"pk": {"S": scan.scan_id}},
            UpdateExpression="SET #state=:failed, error_code=:code, failed_at=:now REMOVE lease_expires",
            ConditionExpression="#owner=:owner AND #state=:progress",
            ExpressionAttributeNames={"#owner": "owner", "#state": "state"},
            ExpressionAttributeValues={":owner": {"S": owner}, ":progress": {"S": "IN_PROGRESS"},
                                       ":failed": {"S": "FAILED"}, ":code": {"S": code}, ":now": {"N": str(now)}},
        )
    except Exception:
        return


def _verify_guardduty_tag(s3: Any, scan: ScanEvent, config: Config) -> None:
    response = s3.get_object_tagging(Bucket=scan.bucket, Key=scan.key, VersionId=scan.version_id,
                                     ExpectedBucketOwner=config.expected_account_id)
    if response.get("VersionId") not in (None, scan.version_id):
        raise SecurityRejected("TAG_VERSION_MISMATCH")
    values = [row.get("Value") for row in response.get("TagSet", []) if row.get("Key") == "GuardDutyMalwareScanStatus"]
    if values != [scan.status]:
        raise SecurityRejected("GUARDDUTY_TAG_MISMATCH")


def _verify_clean_source(s3: Any, scan: ScanEvent, config: Config) -> dict[str, Any]:
    head = s3.head_object(Bucket=scan.bucket, Key=scan.key, VersionId=scan.version_id,
                          ChecksumMode="ENABLED", ExpectedBucketOwner=config.expected_account_id)
    if head.get("VersionId") != scan.version_id or str(head.get("ETag", "")).strip('"') != scan.etag:
        raise SecurityRejected("SOURCE_VERSION_ETAG_MISMATCH")
    if head.get("ServerSideEncryption") != "aws:kms" or head.get("SSEKMSKeyId") != config.source_kms_key_arn:
        raise SecurityRejected("SOURCE_ENCRYPTION_MISMATCH")
    size = head.get("ContentLength")
    metadata = head.get("Metadata") or {}
    content_sha = metadata.get("content-sha256", "")
    raw_sha = metadata.get("raw-sha256", "")
    if not isinstance(size, int) or not 1 <= size <= config.max_object_bytes:
        raise SecurityRejected("SOURCE_SIZE_LIMIT")
    if metadata.get("tenant-id") != config.tenant_id or not HEX64_RE.fullmatch(content_sha) or not HEX64_RE.fullmatch(raw_sha):
        raise SecurityRejected("SOURCE_METADATA_MISMATCH")
    if not scan.key.endswith(f"-{content_sha}.bin"):
        raise SecurityRejected("SOURCE_CONTENT_ADDRESS_MISMATCH")
    return {"bytes": size, "content_sha256": content_sha, "raw_sha256": raw_sha}


def _verify_retention(s3: Any, bucket: str, key: str, version_id: str, config: Config) -> None:
    retention = s3.get_object_retention(Bucket=bucket, Key=key, VersionId=version_id,
                                        ExpectedBucketOwner=config.expected_account_id).get("Retention") or {}
    if retention.get("Mode") != config.retention_mode:
        raise SecurityRejected("RETENTION_MODE_MISMATCH")
    until = retention.get("RetainUntilDate")
    if not isinstance(until, datetime) or until <= datetime.now(timezone.utc) + timedelta(days=config.retain_days - 1):
        raise SecurityRejected("RETENTION_NOT_PROVEN")


def _ensure_clean_copy(s3: Any, scan: ScanEvent, source: dict[str, Any], config: Config) -> dict[str, Any]:
    key = f"clean/{config.tenant_id}/{source['content_sha256']}.bin"
    expected_metadata = {
        "content-sha256": source["content_sha256"], "raw-sha256": source["raw_sha256"],
        "tenant-id": config.tenant_id, "scan-id": scan.scan_id,
        "source-version-sha256": _sha256(scan.version_id.encode("utf-8")),
    }
    try:
        head = s3.head_object(Bucket=config.release_bucket, Key=key, ChecksumMode="ENABLED",
                              ExpectedBucketOwner=config.expected_account_id)
    except Exception as exc:
        if not _is_not_found(exc):
            raise
        response = s3.copy_object(
            Bucket=config.release_bucket, Key=key,
            CopySource={"Bucket": scan.bucket, "Key": scan.key, "VersionId": scan.version_id},
            CopySourceIfMatch=f'"{scan.etag}"', IfNoneMatch="*", MetadataDirective="REPLACE",
            Metadata=expected_metadata, ContentType="application/octet-stream", ChecksumAlgorithm="SHA256",
            ServerSideEncryption="aws:kms", SSEKMSKeyId=config.release_kms_key_arn, BucketKeyEnabled=True,
            ObjectLockMode=config.retention_mode,
            ObjectLockRetainUntilDate=datetime.now(timezone.utc) + timedelta(days=config.retain_days),
            ExpectedBucketOwner=config.expected_account_id, ExpectedSourceBucketOwner=config.expected_account_id,
        )
        version_id = _required_string(response.get("VersionId"), "RELEASE_VERSION_MISSING", 1024)
        head = s3.head_object(Bucket=config.release_bucket, Key=key, VersionId=version_id, ChecksumMode="ENABLED",
                              ExpectedBucketOwner=config.expected_account_id)
    metadata = head.get("Metadata") or {}
    version_id = _required_string(head.get("VersionId"), "RELEASE_VERSION_MISSING", 1024)
    if metadata != expected_metadata or head.get("ContentLength") != source["bytes"]:
        raise SecurityRejected("RELEASE_OBJECT_MISMATCH")
    if head.get("ServerSideEncryption") != "aws:kms" or head.get("SSEKMSKeyId") != config.release_kms_key_arn:
        raise SecurityRejected("RELEASE_ENCRYPTION_MISMATCH")
    _verify_retention(s3, config.release_bucket, key, version_id, config)
    return {"bucket": config.release_bucket, "key": key, "version_id": version_id,
            "sha256": source["content_sha256"], "bytes": source["bytes"], "retention_verified": True}


def _put_receipt(s3: Any, scan: ScanEvent, config: Config, receipt: dict[str, Any]) -> dict[str, str]:
    key = _receipt_key(scan, config)
    data = _canonical_json(receipt)
    metadata = {"scan-id": scan.scan_id, "scan-status": scan.status, "receipt-sha256": _sha256(data),
                "tenant-id": config.tenant_id}
    try:
        response = s3.put_object(
            Bucket=config.receipt_bucket, Key=key, Body=data, ContentLength=len(data), ContentType="application/json",
            ChecksumAlgorithm="SHA256", ChecksumSHA256=_checksum_b64(data), IfNoneMatch="*", Metadata=metadata,
            ServerSideEncryption="aws:kms", SSEKMSKeyId=config.receipt_kms_key_arn, BucketKeyEnabled=True,
            ObjectLockMode=config.retention_mode,
            ObjectLockRetainUntilDate=datetime.now(timezone.utc) + timedelta(days=config.retain_days),
            ExpectedBucketOwner=config.expected_account_id,
        )
        version_id = _required_string(response.get("VersionId"), "RECEIPT_VERSION_MISSING", 1024)
    except Exception as exc:
        if not _is_precondition(exc):
            raise
        existing = _existing_receipt(s3, scan, config)
        if existing is None:
            raise SecurityRejected("RECEIPT_RACE_UNRESOLVED") from exc
        return existing
    head = s3.head_object(Bucket=config.receipt_bucket, Key=key, VersionId=version_id,
                          ChecksumMode="ENABLED", ExpectedBucketOwner=config.expected_account_id)
    if head.get("Metadata") != metadata or head.get("ContentLength") != len(data):
        raise SecurityRejected("RECEIPT_VERIFY_FAILED")
    _verify_retention(s3, config.receipt_bucket, key, version_id, config)
    return {"bucket": config.receipt_bucket, "key": key, "version_id": version_id}


def _alert(sns: Any, scan: ScanEvent, config: Config, decision: str) -> None:
    if decision == "ADMIT_TO_EXTRACTION":
        return
    body = _canonical_json({"scan_id": scan.scan_id, "status": scan.status, "decision": decision,
                            "status_reasons": list(scan.status_reasons)})
    sns.publish(TopicArn=config.alert_topic_arn, Message=body.decode("utf-8"),
                MessageAttributes={"scan_status": {"DataType": "String", "StringValue": scan.status}})


def process_event(event: Any, clients: AwsClients, config: Config, now: int | None = None) -> dict[str, Any]:
    now = int(time.time()) if now is None else now
    scan = parse_guardduty_event(event, config)
    existing = _existing_receipt(clients.s3, scan, config)
    if existing is not None:
        return {"state": "DUPLICATE", "scan_id": scan.scan_id, "receipt": existing}
    owner = uuid.uuid4().hex
    _acquire_lease(clients.dynamodb, scan, config, owner, now)
    try:
        _verify_guardduty_tag(clients.s3, scan, config)
        released = None
        if scan.status == "NO_THREATS_FOUND":
            source = _verify_clean_source(clients.s3, scan, config)
            released = _ensure_clean_copy(clients.s3, scan, source, config)
            decision = "ADMIT_TO_EXTRACTION"
        elif scan.status == "THREATS_FOUND":
            decision = "REJECT_MALWARE"
        else:
            decision = "REJECT_UNSCANNED"
        receipt = {
            "schema": "elite-aws-guardduty-release-receipt/v1",
            "receipt_version": 1, "scan_id": scan.scan_id, "tenant_id": config.tenant_id,
            "provider": "AWS_GUARDDUTY_MALWARE_PROTECTION_S3", "provider_schema_version": "1.0",
            "source": {"bucket": scan.bucket, "key_sha256": _sha256(scan.key.encode("utf-8")),
                       "version_id_sha256": _sha256(scan.version_id.encode("utf-8")), "etag": scan.etag},
            "scan": {"status": scan.status, "scan_status": scan.scan_status,
                     "status_reasons": list(scan.status_reasons), "threat_name_hashes": list(scan.threat_name_hashes)},
            "decision": decision, "released_object": released,
            "automatic_business_persistence_authorized": False,
            "next_gate": "DOCUMENT_CLASSIFICATION_AND_FIELD_EVIDENCE" if released else "SECURITY_REVIEW",
        }
        saved = _put_receipt(clients.s3, scan, config, receipt)
        _alert(clients.sns, scan, config, decision)
        _finish_lease(clients.dynamodb, scan, config, owner, decision, saved["key"], now)
        return {"state": decision, "scan_id": scan.scan_id, "receipt": saved, "released": released}
    except Exception as exc:
        code = exc.code if isinstance(exc, SecurityRejected) else "UNEXPECTED_PROVIDER_FAILURE"
        _mark_failed(clients.dynamodb, scan, config, owner, code, now)
        raise


def lambda_handler(event: Any, context: Any) -> dict[str, list[dict[str, str]]]:
    config = Config.from_env()
    import boto3
    clients = AwsClients(s3=boto3.client("s3"), dynamodb=boto3.client("dynamodb"), sns=boto3.client("sns"))
    failures: list[dict[str, str]] = []
    records = event.get("Records") if isinstance(event, dict) else None
    if not isinstance(records, list):
        raise SecurityRejected("INVALID_SQS_EVENT")
    for record in records:
        message_id = str(record.get("messageId", "")) if isinstance(record, dict) else ""
        try:
            if not message_id or len(message_id) > 256:
                raise SecurityRejected("INVALID_SQS_MESSAGE_ID")
            body = json.loads(record.get("body", ""))
            result = process_event(body, clients, config)
            print(json.dumps({"level": "INFO", "scan_id": result["scan_id"], "state": result["state"]}, sort_keys=True))
        except Exception as exc:
            code = exc.code if isinstance(exc, SecurityRejected) else "UNEXPECTED_PROVIDER_FAILURE"
            print(json.dumps({"level": "ERROR", "message_id_sha256": _sha256(message_id.encode("utf-8")),
                              "error_code": code, "state": "RETRY"}, sort_keys=True))
            if message_id:
                failures.append({"itemIdentifier": message_id})
    return {"batchItemFailures": failures}
````

### FILE: `aws_guardduty_immutable_release_gate/test_handler.py`
```yaml
block_id: "AWS-GUARDDUTY-IMMUTABLE-RELEASE-GATE:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-guardduty-immutable-release-gate/test_handler.py"
license: "LicenseRef-Workspace-Owner"
sha256: "5005cadad3055dcfe4e8491772165d36fa8ebc09b79cf97759b8500f5f676e2e"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import io
import json
import os
import unittest
from contextlib import redirect_stdout
from datetime import datetime, timedelta, timezone
from unittest.mock import patch

import handler


class FakeAwsError(RuntimeError):
    def __init__(self, code: str):
        super().__init__(code)
        self.response = {"Error": {"Code": code}}


class FakeS3:
    def __init__(self) -> None:
        self.objects: dict[tuple[str, str, str], dict] = {}
        self.latest: dict[tuple[str, str], str] = {}
        self.version_counter = 0
        self.calls: list[tuple[str, dict]] = []
        self.fail_receipt_once = False

    def seed(self, bucket: str, key: str, version: str, body: bytes, metadata: dict[str, str], kms: str,
             etag: str = "abc123", tags: dict[str, str] | None = None, retained: bool = False) -> None:
        self.objects[(bucket, key, version)] = {
            "Body": body, "ContentLength": len(body), "Metadata": dict(metadata), "SSEKMSKeyId": kms,
            "ServerSideEncryption": "aws:kms", "ETag": etag, "VersionId": version,
            "TagSet": [{"Key": k, "Value": v} for k, v in (tags or {}).items()],
            "Retention": {"Mode": "COMPLIANCE", "RetainUntilDate": datetime.now(timezone.utc) + timedelta(days=400)} if retained else None,
        }
        self.latest[(bucket, key)] = version

    def _find(self, bucket: str, key: str, version: str | None) -> dict:
        version = version or self.latest.get((bucket, key))
        value = self.objects.get((bucket, key, version or ""))
        if value is None:
            raise FakeAwsError("NoSuchKey")
        return value

    def head_object(self, Bucket, Key, VersionId=None, **kwargs):
        self.calls.append(("head", {"Bucket": Bucket, "Key": Key, "VersionId": VersionId, **kwargs}))
        return dict(self._find(Bucket, Key, VersionId))

    def get_object_tagging(self, Bucket, Key, VersionId, **kwargs):
        self.calls.append(("tags", {"Bucket": Bucket, "Key": Key, "VersionId": VersionId, **kwargs}))
        value = self._find(Bucket, Key, VersionId)
        return {"VersionId": VersionId, "TagSet": list(value["TagSet"])}

    def get_object_retention(self, Bucket, Key, VersionId, **kwargs):
        value = self._find(Bucket, Key, VersionId)
        return {"Retention": value["Retention"] or {}}

    def _new_version(self) -> str:
        self.version_counter += 1
        return f"generated-{self.version_counter}"

    def copy_object(self, Bucket, Key, CopySource, **kwargs):
        self.calls.append(("copy", {"Bucket": Bucket, "Key": Key, "CopySource": dict(CopySource), **kwargs}))
        if kwargs.get("IfNoneMatch") == "*" and (Bucket, Key) in self.latest:
            raise FakeAwsError("PreconditionFailed")
        source = self._find(CopySource["Bucket"], CopySource["Key"], CopySource["VersionId"])
        if kwargs.get("CopySourceIfMatch", "").strip('"') != source["ETag"]:
            raise FakeAwsError("PreconditionFailed")
        version = self._new_version()
        self.seed(Bucket, Key, version, source["Body"], kwargs["Metadata"], kwargs["SSEKMSKeyId"], source["ETag"], retained=True)
        return {"VersionId": version, "CopyObjectResult": {"ETag": source["ETag"]}}

    def put_object(self, Bucket, Key, Body, **kwargs):
        self.calls.append(("put", {"Bucket": Bucket, "Key": Key, **kwargs}))
        if self.fail_receipt_once and Key.startswith("security-receipts/"):
            self.fail_receipt_once = False
            raise RuntimeError("injected receipt outage")
        if kwargs.get("IfNoneMatch") == "*" and (Bucket, Key) in self.latest:
            raise FakeAwsError("PreconditionFailed")
        data = bytes(Body)
        version = self._new_version()
        self.seed(Bucket, Key, version, data, kwargs["Metadata"], kwargs["SSEKMSKeyId"], handler._sha256(data)[:32], retained=True)
        return {"VersionId": version, "ChecksumSHA256": kwargs.get("ChecksumSHA256")}


class FakeDynamoDB:
    def __init__(self) -> None:
        self.items: dict[str, dict[str, str | int]] = {}

    def put_item(self, TableName, Item, ExpressionAttributeValues, **kwargs):
        pk = Item["pk"]["S"]
        current = self.items.get(pk)
        now = int(ExpressionAttributeValues[":now"]["N"])
        if current and current["state"] != "FAILED" and int(current.get("lease_expires", 0)) >= now:
            raise FakeAwsError("ConditionalCheckFailedException")
        self.items[pk] = {"state": "IN_PROGRESS", "owner": Item["owner"]["S"],
                          "lease_expires": int(Item["lease_expires"]["N"]), "expires_at": int(Item["expires_at"]["N"])}

    def update_item(self, Key, ExpressionAttributeValues, **kwargs):
        pk = Key["pk"]["S"]
        item = self.items[pk]
        if item["owner"] != ExpressionAttributeValues[":owner"]["S"] or item["state"] != "IN_PROGRESS":
            raise FakeAwsError("ConditionalCheckFailedException")
        if ":failed" in ExpressionAttributeValues:
            item["state"] = "FAILED"
            item["error_code"] = ExpressionAttributeValues[":code"]["S"]
        else:
            item["state"] = ExpressionAttributeValues[":state"]["S"]
            item["receipt_key"] = ExpressionAttributeValues[":receipt"]["S"]
        item.pop("lease_expires", None)


class FakeSNS:
    def __init__(self) -> None:
        self.messages: list[dict] = []

    def publish(self, **kwargs):
        self.messages.append(kwargs)
        return {"MessageId": "alert-1"}


def make_config() -> handler.Config:
    return handler.Config(
        expected_account_id="111122223333", expected_region="us-east-1",
        expected_plan_arn="arn:aws:guardduty:us-east-1:111122223333:malware-protection-plan/plan-1",
        tenant_id="tenant-a", quarantine_bucket="quarantine", quarantine_prefix="quarantine/tenant-a/",
        release_bucket="release", receipt_bucket="receipts",
        source_kms_key_arn="arn:aws:kms:us-east-1:111122223333:key/source",
        release_kms_key_arn="arn:aws:kms:us-east-1:111122223333:key/release",
        receipt_kms_key_arn="arn:aws:kms:us-east-1:111122223333:key/receipt",
        idempotency_table="security-decisions", alert_topic_arn="arn:aws:sns:us-east-1:111122223333:security",
    )


def make_event(status: str = "NO_THREATS_FOUND") -> dict:
    scan_status = handler.STATUS_MAP[status]
    result = {"scanResultStatus": status, "threats": None, "statusReasons": None}
    if status == "THREATS_FOUND":
        result["threats"] = [{"name": "EICAR-Test-File", "source": "AMAZON"}]
    elif status in {"UNSUPPORTED", "ACCESS_DENIED"}:
        result["statusReasons"] = ["PASSWORD_PROTECTED" if status == "UNSUPPORTED" else "UNAUTHORIZED_TO_GET_OBJECT"]
    return {
        "version": "0", "id": "72c7d362-737a-6dce-fc78-9e27a0171419",
        "detail-type": "GuardDuty Malware Protection Object Scan Result", "source": "aws.guardduty",
        "account": "111122223333", "region": "us-east-1",
        "resources": ["arn:aws:guardduty:us-east-1:111122223333:malware-protection-plan/plan-1"],
        "detail": {"schemaVersion": "1.0", "scanStatus": scan_status, "resourceType": "S3_OBJECT",
                   "s3ObjectDetails": {"bucketName": "quarantine",
                                       "objectKey": "quarantine/tenant-a/" + "a" * 64 + "/part-0001-" + "b" * 64 + ".bin",
                                       "eTag": "abc123", "versionId": "v-source", "s3Throttled": False},
                   "scanResultDetails": result},
    }


class GuardDutyReleaseTests(unittest.TestCase):
    def setUp(self) -> None:
        self.config = make_config()
        self.s3 = FakeS3()
        self.ddb = FakeDynamoDB()
        self.sns = FakeSNS()
        self.clients = handler.AwsClients(self.s3, self.ddb, self.sns)
        self.source_key = make_event()["detail"]["s3ObjectDetails"]["objectKey"]
        self.s3.seed("quarantine", self.source_key, "v-source", b"PDF", {
            "content-sha256": "b" * 64, "raw-sha256": "a" * 64, "tenant-id": "tenant-a", "part-index": "1",
        }, self.config.source_kms_key_arn, tags={"GuardDutyMalwareScanStatus": "NO_THREATS_FOUND"})

    def _receipt_body(self) -> dict:
        receipt_keys = [row for row in self.s3.latest if row[0] == "receipts"]
        self.assertEqual(1, len(receipt_keys))
        version = self.s3.latest[receipt_keys[0]]
        return json.loads(self.s3.objects[(receipt_keys[0][0], receipt_keys[0][1], version)]["Body"])

    def test_clean_version_is_released_and_receipt_blocks_business_persistence(self):
        result = handler.process_event(make_event(), self.clients, self.config, now=100)
        self.assertEqual("ADMIT_TO_EXTRACTION", result["state"])
        self.assertEqual(1, len([x for x in self.s3.calls if x[0] == "copy"]))
        copy = [x[1] for x in self.s3.calls if x[0] == "copy"][0]
        self.assertEqual("v-source", copy["CopySource"]["VersionId"])
        self.assertEqual("*", copy["IfNoneMatch"])
        self.assertEqual('"abc123"', copy["CopySourceIfMatch"])
        receipt = self._receipt_body()
        self.assertEqual("elite-aws-guardduty-release-receipt/v1", receipt["schema"])
        self.assertFalse(receipt["automatic_business_persistence_authorized"])
        self.assertEqual("DOCUMENT_CLASSIFICATION_AND_FIELD_EVIDENCE", receipt["next_gate"])
        self.assertNotIn(self.source_key, json.dumps(receipt))
        self.assertEqual([], self.sns.messages)

    def test_durable_receipt_short_circuits_duplicate(self):
        first = handler.process_event(make_event(), self.clients, self.config, now=100)
        before = len(self.s3.calls)
        second = handler.process_event(make_event(), self.clients, self.config, now=101)
        self.assertEqual("DUPLICATE", second["state"])
        self.assertEqual(before + 1, len(self.s3.calls))
        self.assertEqual(first["scan_id"], second["scan_id"])

    def test_threat_is_not_read_or_copied_and_alert_redacts_name(self):
        event = make_event("THREATS_FOUND")
        self.s3.objects[("quarantine", self.source_key, "v-source")]["TagSet"] = [{"Key": "GuardDutyMalwareScanStatus", "Value": "THREATS_FOUND"}]
        result = handler.process_event(event, self.clients, self.config, now=100)
        self.assertEqual("REJECT_MALWARE", result["state"])
        self.assertEqual([], [x for x in self.s3.calls if x[0] == "copy"])
        self.assertFalse(any(x[0] == "head" and x[1]["Bucket"] == "quarantine" for x in self.s3.calls))
        receipt_text = json.dumps(self._receipt_body())
        self.assertNotIn("EICAR-Test-File", receipt_text)
        self.assertEqual(1, len(self.sns.messages))

    def test_unsupported_is_durable_rejection_not_false_clean(self):
        self.s3.objects[("quarantine", self.source_key, "v-source")]["TagSet"] = [{"Key": "GuardDutyMalwareScanStatus", "Value": "UNSUPPORTED"}]
        result = handler.process_event(make_event("UNSUPPORTED"), self.clients, self.config, now=100)
        self.assertEqual("REJECT_UNSCANNED", result["state"])
        self.assertEqual("SECURITY_REVIEW", self._receipt_body()["next_gate"])
        self.assertEqual(1, len(self.sns.messages))

    def test_authority_and_status_mismatches_fail_closed(self):
        cases = []
        wrong_account = make_event(); wrong_account["account"] = "999900001111"; cases.append(wrong_account)
        wrong_region = make_event(); wrong_region["region"] = "eu-west-1"; cases.append(wrong_region)
        wrong_plan = make_event(); wrong_plan["resources"] = ["arn:aws:guardduty:us-east-1:111122223333:malware-protection-plan/other"]; cases.append(wrong_plan)
        wrong_status = make_event(); wrong_status["detail"]["scanStatus"] = "SKIPPED"; cases.append(wrong_status)
        for value in cases:
            with self.subTest(value=value):
                with self.assertRaises(handler.SecurityRejected):
                    handler.parse_guardduty_event(value, self.config)

    def test_only_receiver_attachment_keys_and_exact_versions_are_admitted(self):
        cases = []
        no_version = make_event(); no_version["detail"]["s3ObjectDetails"].pop("versionId"); cases.append(no_version)
        traversal = make_event(); traversal["detail"]["s3ObjectDetails"]["objectKey"] = "quarantine/tenant-a/../bad"; cases.append(traversal)
        manifest = make_event(); manifest["detail"]["s3ObjectDetails"]["objectKey"] = "receipts/tenant-a/manifest.json"; cases.append(manifest)
        for value in cases:
            with self.assertRaises(handler.SecurityRejected):
                handler.parse_guardduty_event(value, self.config)

    def test_guardduty_tag_must_match_exact_version_and_event(self):
        self.s3.objects[("quarantine", self.source_key, "v-source")]["TagSet"] = [{"Key": "GuardDutyMalwareScanStatus", "Value": "THREATS_FOUND"}]
        with self.assertRaisesRegex(handler.SecurityRejected, "GUARDDUTY_TAG_MISMATCH"):
            handler.process_event(make_event(), self.clients, self.config, now=100)
        tags = [x[1] for x in self.s3.calls if x[0] == "tags"]
        self.assertEqual("v-source", tags[0]["VersionId"])

    def test_source_metadata_encryption_size_and_etag_are_enforced(self):
        mutations = [
            ("SSEKMSKeyId", "wrong"), ("ETag", "wrong"), ("ContentLength", 0),
            ("Metadata", {"content-sha256": "b" * 64, "raw-sha256": "a" * 64, "tenant-id": "other"}),
        ]
        for field, value in mutations:
            with self.subTest(field=field):
                self.setUp()
                self.s3.objects[("quarantine", self.source_key, "v-source")][field] = value
                with self.assertRaises(Exception):
                    handler.process_event(make_event(), self.clients, self.config, now=100)

    def test_partial_copy_is_reconciled_without_overwrite(self):
        self.s3.fail_receipt_once = True
        with self.assertRaisesRegex(RuntimeError, "injected receipt outage"):
            handler.process_event(make_event(), self.clients, self.config, now=100)
        self.assertEqual(1, len([x for x in self.s3.calls if x[0] == "copy"]))
        result = handler.process_event(make_event(), self.clients, self.config, now=101)
        self.assertEqual("ADMIT_TO_EXTRACTION", result["state"])
        self.assertEqual(1, len([x for x in self.s3.calls if x[0] == "copy"]))

    def test_preexisting_wrong_release_object_fails_closed(self):
        key = "clean/tenant-a/" + "b" * 64 + ".bin"
        self.s3.seed("release", key, "attacker", b"bad", {"content-sha256": "0" * 64},
                     self.config.release_kms_key_arn, retained=True)
        with self.assertRaisesRegex(handler.SecurityRejected, "RELEASE_OBJECT_MISMATCH"):
            handler.process_event(make_event(), self.clients, self.config, now=100)

    def test_existing_receipt_requires_matching_scan_authority(self):
        scan = handler.parse_guardduty_event(make_event(), self.config)
        key = handler._receipt_key(scan, self.config)
        self.s3.seed("receipts", key, "bad", b"{}", {"scan-id": "wrong", "scan-status": "NO_THREATS_FOUND"},
                     self.config.receipt_kms_key_arn, retained=True)
        with self.assertRaisesRegex(handler.SecurityRejected, "EXISTING_RECEIPT_AUTHORITY_MISMATCH"):
            handler.process_event(make_event(), self.clients, self.config, now=100)

    def test_busy_lease_retries_and_expired_lease_is_reclaimed(self):
        scan = handler.parse_guardduty_event(make_event(), self.config)
        self.ddb.items[scan.scan_id] = {"state": "IN_PROGRESS", "owner": "other", "lease_expires": 200}
        with self.assertRaisesRegex(handler.SecurityRejected, "IDEMPOTENCY_LEASE_BUSY"):
            handler.process_event(make_event(), self.clients, self.config, now=100)
        result = handler.process_event(make_event(), self.clients, self.config, now=201)
        self.assertEqual("ADMIT_TO_EXTRACTION", result["state"])

    def test_owner_fence_rejects_stale_completion(self):
        scan = handler.parse_guardduty_event(make_event(), self.config)
        self.ddb.items[scan.scan_id] = {"state": "IN_PROGRESS", "owner": "owner-new", "lease_expires": 200}
        with self.assertRaisesRegex(handler.SecurityRejected, "IDEMPOTENCY_OWNER_FENCE"):
            handler._finish_lease(self.ddb, scan, self.config, "owner-old", "DONE", "receipt", 100)

    def test_missing_object_lock_proof_rejects_release(self):
        original = self.s3.get_object_retention
        self.s3.get_object_retention = lambda **kwargs: {"Retention": {}}
        try:
            with self.assertRaisesRegex(handler.SecurityRejected, "RETENTION_MODE_MISMATCH"):
                handler.process_event(make_event(), self.clients, self.config, now=100)
        finally:
            self.s3.get_object_retention = original

    def test_batch_response_retries_only_failed_message(self):
        env = {
            "EXPECTED_ACCOUNT_ID": self.config.expected_account_id, "EXPECTED_REGION": self.config.expected_region,
            "EXPECTED_PLAN_ARN": self.config.expected_plan_arn, "TENANT_ID": self.config.tenant_id,
            "QUARANTINE_BUCKET": self.config.quarantine_bucket, "QUARANTINE_PREFIX": self.config.quarantine_prefix,
            "RELEASE_BUCKET": self.config.release_bucket, "RECEIPT_BUCKET": self.config.receipt_bucket,
            "SOURCE_KMS_KEY_ARN": self.config.source_kms_key_arn, "RELEASE_KMS_KEY_ARN": self.config.release_kms_key_arn,
            "RECEIPT_KMS_KEY_ARN": self.config.receipt_kms_key_arn, "IDEMPOTENCY_TABLE": self.config.idempotency_table,
            "ALERT_TOPIC_ARN": self.config.alert_topic_arn,
        }
        good = {"messageId": "good", "body": json.dumps(make_event())}
        bad = {"messageId": "bad", "body": "not-json"}
        fake_boto3 = type("Boto3", (), {"client": lambda _, name: {"s3": self.s3, "dynamodb": self.ddb, "sns": self.sns}[name]})()
        with patch.dict(os.environ, env, clear=True), patch.dict("sys.modules", {"boto3": fake_boto3}):
            with redirect_stdout(io.StringIO()) as output:
                result = handler.lambda_handler({"Records": [good, bad]}, None)
        self.assertEqual({"batchItemFailures": [{"itemIdentifier": "bad"}]}, result)
        self.assertNotIn(self.source_key, output.getvalue())

    def test_config_rejects_unsafe_prefix_and_empty_authority(self):
        values = dict(make_config().__dict__)
        values["quarantine_prefix"] = "quarantine/other/"
        with self.assertRaisesRegex(handler.SecurityRejected, "INVALID_QUARANTINE_PREFIX"):
            handler.Config(**values)
        values = dict(make_config().__dict__)
        values["alert_topic_arn"] = ""
        with self.assertRaisesRegex(handler.SecurityRejected, "INVALID_ALERT_TOPIC_ARN"):
            handler.Config(**values)


if __name__ == "__main__":
    unittest.main(verbosity=2)
````

### FILE: `aws_guardduty_immutable_release_gate/template.yaml`
```yaml
block_id: "AWS-GUARDDUTY-IMMUTABLE-RELEASE-GATE:template:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/aws-samples/guardduty-malware-protection/blob/fd9d2bdf46eb52bb336a7e7a837bd38ef73002e2/cfn/template.yaml; narrowed/hardened, deletable resources removed"
license: "MIT-0"
sha256: "ec58aa530de0b4a6d0d04c8d5d644100f25288b00ccf4a5f0b586badb53399ea"
variables: []
secrets_allowed: false
```
````yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Description: Conditioned GuardDuty Malware Protection for S3 release gate adapted from aws-samples/guardduty-malware-protection

Parameters:
  TenantId:
    Type: String
    AllowedPattern: '^[a-z0-9][a-z0-9_-]{0,62}$'
  QuarantineBucketName:
    Type: String
    AllowedPattern: '^(?![.-])(?!.*\.\.)(?!.*-\.)(?!.*\.-)[a-z0-9.-]{3,63}(?<![.-])$'
  SourceKmsKeyArn:
    Type: String
    AllowedPattern: '^arn:(aws|aws-us-gov|aws-cn):kms:[a-z0-9-]+:[0-9]{12}:key/[A-Za-z0-9-]+$'
  ReleaseBucketName:
    Type: String
    AllowedPattern: '^(?![.-])(?!.*\.\.)(?!.*-\.)(?!.*\.-)[a-z0-9.-]{3,63}(?<![.-])$'
  ReceiptBucketName:
    Type: String
    AllowedPattern: '^(?![.-])(?!.*\.\.)(?!.*-\.)(?!.*\.-)[a-z0-9.-]{3,63}(?<![.-])$'
  SecurityAlertEmail:
    Type: String
    AllowedPattern: '^[^\s@]+@[^\s@]+\.[^\s@]+$'
  RetainDays:
    Type: Number
    Default: 365
    MinValue: 1
    MaxValue: 36500
  DedupeDays:
    Type: Number
    Default: 400
    MinValue: 1
    MaxValue: 36500
  MaxObjectBytes:
    Type: Number
    Default: 52428800
    MinValue: 1
    MaxValue: 5368709120
  ReservedConcurrency:
    Type: Number
    Default: 5
    MinValue: 2
    MaxValue: 100

Rules:
  DedupeCoversRetention:
    Assertions:
      - Assert: !Not [!Equals [!Ref DedupeDays, 0]]
        AssertDescription: DedupeDays must be non-zero; deploy.ps1 additionally requires it to be at least RetainDays.

Resources:
  ReleaseKey:
    Type: AWS::KMS::Key
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      Description: KMS key for clean objects released after exact GuardDuty decision
      EnableKeyRotation: true
      PendingWindowInDays: 30
      KeyPolicy:
        Version: '2012-10-17'
        Statement:
          - Sid: AccountAdministration
            Effect: Allow
            Principal: {AWS: !Sub 'arn:${AWS::Partition}:iam::${AWS::AccountId}:root'}
            Action: 'kms:*'
            Resource: '*'

  ReceiptKey:
    Type: AWS::KMS::Key
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      Description: KMS key for immutable GuardDuty security receipts
      EnableKeyRotation: true
      PendingWindowInDays: 30
      KeyPolicy:
        Version: '2012-10-17'
        Statement:
          - Sid: AccountAdministration
            Effect: Allow
            Principal: {AWS: !Sub 'arn:${AWS::Partition}:iam::${AWS::AccountId}:root'}
            Action: 'kms:*'
            Resource: '*'

  ReleaseBucket:
    Type: AWS::S3::Bucket
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      BucketName: !Ref ReleaseBucketName
      BucketEncryption:
        ServerSideEncryptionConfiguration:
          - BucketKeyEnabled: true
            ServerSideEncryptionByDefault:
              SSEAlgorithm: aws:kms
              KMSMasterKeyID: !GetAtt ReleaseKey.Arn
      VersioningConfiguration: {Status: Enabled}
      ObjectLockEnabled: true
      ObjectLockConfiguration:
        ObjectLockEnabled: Enabled
        Rule:
          DefaultRetention:
            Mode: COMPLIANCE
            Days: !Ref RetainDays
      OwnershipControls:
        Rules: [{ObjectOwnership: BucketOwnerEnforced}]
      PublicAccessBlockConfiguration:
        BlockPublicAcls: true
        BlockPublicPolicy: true
        IgnorePublicAcls: true
        RestrictPublicBuckets: true

  ReceiptBucket:
    Type: AWS::S3::Bucket
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      BucketName: !Ref ReceiptBucketName
      BucketEncryption:
        ServerSideEncryptionConfiguration:
          - BucketKeyEnabled: true
            ServerSideEncryptionByDefault:
              SSEAlgorithm: aws:kms
              KMSMasterKeyID: !GetAtt ReceiptKey.Arn
      VersioningConfiguration: {Status: Enabled}
      ObjectLockEnabled: true
      ObjectLockConfiguration:
        ObjectLockEnabled: Enabled
        Rule:
          DefaultRetention:
            Mode: COMPLIANCE
            Days: !Ref RetainDays
      OwnershipControls:
        Rules: [{ObjectOwnership: BucketOwnerEnforced}]
      PublicAccessBlockConfiguration:
        BlockPublicAcls: true
        BlockPublicPolicy: true
        IgnorePublicAcls: true
        RestrictPublicBuckets: true

  ReleaseBucketPolicy:
    Type: AWS::S3::BucketPolicy
    Properties:
      Bucket: !Ref ReleaseBucket
      PolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Sid: DenyInsecureTransport
            Effect: Deny
            Principal: '*'
            Action: 's3:*'
            Resource: [!GetAtt ReleaseBucket.Arn, !Sub '${ReleaseBucket.Arn}/*']
            Condition: {Bool: {'aws:SecureTransport': 'false'}}
          - Sid: DenyUnencryptedWrites
            Effect: Deny
            Principal: '*'
            Action: s3:PutObject
            Resource: !Sub '${ReleaseBucket.Arn}/*'
            Condition: {StringNotEquals: {'s3:x-amz-server-side-encryption': 'aws:kms'}}

  ReceiptBucketPolicy:
    Type: AWS::S3::BucketPolicy
    Properties:
      Bucket: !Ref ReceiptBucket
      PolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Sid: DenyInsecureTransport
            Effect: Deny
            Principal: '*'
            Action: 's3:*'
            Resource: [!GetAtt ReceiptBucket.Arn, !Sub '${ReceiptBucket.Arn}/*']
            Condition: {Bool: {'aws:SecureTransport': 'false'}}
          - Sid: DenyUnencryptedWrites
            Effect: Deny
            Principal: '*'
            Action: s3:PutObject
            Resource: !Sub '${ReceiptBucket.Arn}/*'
            Condition: {StringNotEquals: {'s3:x-amz-server-side-encryption': 'aws:kms'}}

  GuardDutyPlanRole:
    Type: AWS::IAM::Role
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      Description: Narrow role used only by GuardDuty Malware Protection for the configured quarantine prefix
      AssumeRolePolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Effect: Allow
            Principal: {Service: malware-protection-plan.guardduty.amazonaws.com}
            Action: sts:AssumeRole
      Policies:
        - PolicyName: GuardDutyMalwareProtectionExactBucket
          PolicyDocument:
            Version: '2012-10-17'
            Statement:
              - Sid: AllowManagedRuleToSendS3EventsToGuardDuty
                Effect: Allow
                Action: [events:PutRule, events:DeleteRule, events:PutTargets, events:RemoveTargets]
                Resource: !Sub 'arn:${AWS::Partition}:events:${AWS::Region}:${AWS::AccountId}:rule/DO-NOT-DELETE-AmazonGuardDutyMalwareProtectionS3*'
                Condition: {StringLike: {events:ManagedBy: malware-protection-plan.guardduty.amazonaws.com}}
              - Sid: AllowGuardDutyToMonitorEventBridgeManagedRule
                Effect: Allow
                Action: [events:DescribeRule, events:ListTargetsByRule]
                Resource: !Sub 'arn:${AWS::Partition}:events:${AWS::Region}:${AWS::AccountId}:rule/DO-NOT-DELETE-AmazonGuardDutyMalwareProtectionS3*'
              - Sid: AllowPostScanTag
                Effect: Allow
                Action: [s3:PutObjectTagging, s3:GetObjectTagging, s3:PutObjectVersionTagging, s3:GetObjectVersionTagging]
                Resource: !Sub 'arn:${AWS::Partition}:s3:::${QuarantineBucketName}/quarantine/${TenantId}/*'
              - Sid: AllowEnableS3EventBridgeEvents
                Effect: Allow
                Action: [s3:PutBucketNotification, s3:GetBucketNotification]
                Resource: !Sub 'arn:${AWS::Partition}:s3:::${QuarantineBucketName}'
              - Sid: AllowPutValidationObject
                Effect: Allow
                Action: s3:PutObject
                Resource: !Sub 'arn:${AWS::Partition}:s3:::${QuarantineBucketName}/malware-protection-resource-validation-object'
              - Sid: AllowCheckBucketOwnership
                Effect: Allow
                Action: [s3:ListBucket, s3:GetBucketLocation]
                Resource: !Sub 'arn:${AWS::Partition}:s3:::${QuarantineBucketName}'
              - Sid: AllowMalwareScan
                Effect: Allow
                Action: [s3:GetObject, s3:GetObjectVersion]
                Resource: !Sub 'arn:${AWS::Partition}:s3:::${QuarantineBucketName}/quarantine/${TenantId}/*'
              - Sid: AllowDecryptForMalwareScan
                Effect: Allow
                Action: [kms:GenerateDataKey, kms:Decrypt]
                Resource: !Ref SourceKmsKeyArn
                Condition: {StringLike: {kms:ViaService: 's3.*.amazonaws.com'}}

  GuardDutyMalwareProtectionPlan:
    Type: AWS::GuardDuty::MalwareProtectionPlan
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      Actions:
        Tagging: {Status: ENABLED}
      ProtectedResource:
        S3Bucket:
          BucketName: !Ref QuarantineBucketName
          ObjectPrefixes: [!Sub 'quarantine/${TenantId}/']
      Role: !GetAtt GuardDutyPlanRole.Arn

  SecurityAlertTopic:
    Type: AWS::SNS::Topic
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      KmsMasterKeyId: alias/aws/sns
      Subscription:
        - Protocol: email
          Endpoint: !Ref SecurityAlertEmail

  ScanResultDLQ:
    Type: AWS::SQS::Queue
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      MessageRetentionPeriod: 1209600
      SqsManagedSseEnabled: true

  ScanResultQueue:
    Type: AWS::SQS::Queue
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      VisibilityTimeout: 1800
      MessageRetentionPeriod: 1209600
      SqsManagedSseEnabled: true
      RedrivePolicy:
        deadLetterTargetArn: !GetAtt ScanResultDLQ.Arn
        maxReceiveCount: 5

  GuardDutyScanResultRule:
    Type: AWS::Events::Rule
    Properties:
      State: ENABLED
      EventPattern:
        source: [aws.guardduty]
        detail-type: [GuardDuty Malware Protection Object Scan Result]
        resources: [!GetAtt GuardDutyMalwareProtectionPlan.Arn]
        detail:
          resourceType: [S3_OBJECT]
          s3ObjectDetails:
            bucketName: [!Ref QuarantineBucketName]
          scanResultDetails:
            scanResultStatus: [NO_THREATS_FOUND, THREATS_FOUND, UNSUPPORTED, ACCESS_DENIED, FAILED]
      Targets:
        - Id: ScanResultQueue
          Arn: !GetAtt ScanResultQueue.Arn

  ScanResultQueuePolicy:
    Type: AWS::SQS::QueuePolicy
    Properties:
      Queues: [!Ref ScanResultQueue]
      PolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Sid: EventBridgeExactRule
            Effect: Allow
            Principal: {Service: events.amazonaws.com}
            Action: sqs:SendMessage
            Resource: !GetAtt ScanResultQueue.Arn
            Condition:
              ArnEquals: {aws:SourceArn: !GetAtt GuardDutyScanResultRule.Arn}
              StringEquals: {aws:SourceAccount: !Ref AWS::AccountId}

  DecisionTable:
    Type: AWS::DynamoDB::Table
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      BillingMode: PAY_PER_REQUEST
      AttributeDefinitions: [{AttributeName: pk, AttributeType: S}]
      KeySchema: [{AttributeName: pk, KeyType: HASH}]
      PointInTimeRecoverySpecification: {PointInTimeRecoveryEnabled: true}
      SSESpecification: {SSEEnabled: true}
      TimeToLiveSpecification: {AttributeName: expires_at, Enabled: true}
      DeletionProtectionEnabled: true

  ReleaseFunctionRole:
    Type: AWS::IAM::Role
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      AssumeRolePolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Effect: Allow
            Principal: {Service: lambda.amazonaws.com}
            Action: sts:AssumeRole
      Policies:
        - PolicyName: ExactVersionGuardDutyRelease
          PolicyDocument:
            Version: '2012-10-17'
            Statement:
              - Sid: ReadScanTagForExactVersion
                Effect: Allow
                Action: [s3:GetObjectTagging, s3:GetObjectVersionTagging]
                Resource: !Sub 'arn:${AWS::Partition}:s3:::${QuarantineBucketName}/quarantine/${TenantId}/*'
              - Sid: ReadOnlyCleanExactVersion
                Effect: Allow
                Action: [s3:GetObject, s3:GetObjectVersion]
                Resource: !Sub 'arn:${AWS::Partition}:s3:::${QuarantineBucketName}/quarantine/${TenantId}/*'
                Condition:
                  StringEquals:
                    s3:ExistingObjectTag/GuardDutyMalwareScanStatus: NO_THREATS_FOUND
              - Sid: ReleaseAndVerify
                Effect: Allow
                Action: [s3:PutObject, s3:GetObject, s3:GetObjectVersion, s3:GetObjectRetention]
                Resource: !Sub '${ReleaseBucket.Arn}/clean/${TenantId}/*'
              - Sid: ReceiptAndVerify
                Effect: Allow
                Action: [s3:PutObject, s3:GetObject, s3:GetObjectVersion, s3:GetObjectRetention]
                Resource: !Sub '${ReceiptBucket.Arn}/security-receipts/${TenantId}/*'
              - Sid: UseSourceKey
                Effect: Allow
                Action: [kms:Decrypt]
                Resource: !Ref SourceKmsKeyArn
                Condition: {StringLike: {kms:ViaService: 's3.*.amazonaws.com'}}
              - Sid: UseDestinationKeys
                Effect: Allow
                Action: [kms:Decrypt, kms:Encrypt, kms:GenerateDataKey]
                Resource: [!GetAtt ReleaseKey.Arn, !GetAtt ReceiptKey.Arn]
                Condition: {StringLike: {kms:ViaService: 's3.*.amazonaws.com'}}
              - Sid: DecisionLease
                Effect: Allow
                Action: [dynamodb:PutItem, dynamodb:UpdateItem]
                Resource: !GetAtt DecisionTable.Arn
              - Sid: RedactedSecurityAlert
                Effect: Allow
                Action: sns:Publish
                Resource: !Ref SecurityAlertTopic
              - Sid: QueueConsumption
                Effect: Allow
                Action: [sqs:ReceiveMessage, sqs:DeleteMessage, sqs:GetQueueAttributes]
                Resource: !GetAtt ScanResultQueue.Arn
              - Sid: Logs
                Effect: Allow
                Action: [logs:CreateLogStream, logs:PutLogEvents]
                Resource: !GetAtt ReleaseLogGroup.Arn

  ReleaseLogGroup:
    Type: AWS::Logs::LogGroup
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      RetentionInDays: 365

  ReleaseFunction:
    Type: AWS::Serverless::Function
    Properties:
      CodeUri: build/function
      Handler: handler.lambda_handler
      Runtime: python3.14
      Architectures: [arm64]
      Role: !GetAtt ReleaseFunctionRole.Arn
      Timeout: 300
      MemorySize: 512
      ReservedConcurrentExecutions: !Ref ReservedConcurrency
      LoggingConfig:
        LogFormat: JSON
        ApplicationLogLevel: INFO
        SystemLogLevel: WARN
        LogGroup: !Ref ReleaseLogGroup
      Environment:
        Variables:
          EXPECTED_ACCOUNT_ID: !Ref AWS::AccountId
          EXPECTED_REGION: !Ref AWS::Region
          EXPECTED_PLAN_ARN: !GetAtt GuardDutyMalwareProtectionPlan.Arn
          TENANT_ID: !Ref TenantId
          QUARANTINE_BUCKET: !Ref QuarantineBucketName
          QUARANTINE_PREFIX: !Sub 'quarantine/${TenantId}/'
          RELEASE_BUCKET: !Ref ReleaseBucket
          RECEIPT_BUCKET: !Ref ReceiptBucket
          SOURCE_KMS_KEY_ARN: !Ref SourceKmsKeyArn
          RELEASE_KMS_KEY_ARN: !GetAtt ReleaseKey.Arn
          RECEIPT_KMS_KEY_ARN: !GetAtt ReceiptKey.Arn
          IDEMPOTENCY_TABLE: !Ref DecisionTable
          ALERT_TOPIC_ARN: !Ref SecurityAlertTopic
          RETAIN_DAYS: !Ref RetainDays
          DEDUPE_DAYS: !Ref DedupeDays
          MAX_OBJECT_BYTES: !Ref MaxObjectBytes
          RETENTION_MODE: COMPLIANCE

  ScanResultMapping:
    Type: AWS::Lambda::EventSourceMapping
    Properties:
      EventSourceArn: !GetAtt ScanResultQueue.Arn
      FunctionName: !Ref ReleaseFunction
      Enabled: true
      BatchSize: 10
      MaximumBatchingWindowInSeconds: 5
      FunctionResponseTypes: [ReportBatchItemFailures]
      ScalingConfig:
        MaximumConcurrency: !Ref ReservedConcurrency

  DLQVisibleAlarm:
    Type: AWS::CloudWatch::Alarm
    Properties:
      AlarmDescription: GuardDuty scan result messages reached the DLQ
      Namespace: AWS/SQS
      MetricName: ApproximateNumberOfMessagesVisible
      Dimensions: [{Name: QueueName, Value: !GetAtt ScanResultDLQ.QueueName}]
      Statistic: Maximum
      Period: 60
      EvaluationPeriods: 1
      Threshold: 0
      ComparisonOperator: GreaterThanThreshold
      TreatMissingData: notBreaching
      AlarmActions: [!Ref SecurityAlertTopic]

  QueueAgeAlarm:
    Type: AWS::CloudWatch::Alarm
    Properties:
      AlarmDescription: GuardDuty scan result processing is delayed
      Namespace: AWS/SQS
      MetricName: ApproximateAgeOfOldestMessage
      Dimensions: [{Name: QueueName, Value: !GetAtt ScanResultQueue.QueueName}]
      Statistic: Maximum
      Period: 60
      EvaluationPeriods: 2
      Threshold: 600
      ComparisonOperator: GreaterThanThreshold
      TreatMissingData: notBreaching
      AlarmActions: [!Ref SecurityAlertTopic]

  FunctionErrorsAlarm:
    Type: AWS::CloudWatch::Alarm
    Properties:
      AlarmDescription: GuardDuty release function errors
      Namespace: AWS/Lambda
      MetricName: Errors
      Dimensions: [{Name: FunctionName, Value: !Ref ReleaseFunction}]
      Statistic: Sum
      Period: 60
      EvaluationPeriods: 1
      Threshold: 0
      ComparisonOperator: GreaterThanThreshold
      TreatMissingData: notBreaching
      AlarmActions: [!Ref SecurityAlertTopic]

Outputs:
  MalwareProtectionPlanArn:
    Value: !GetAtt GuardDutyMalwareProtectionPlan.Arn
  MalwareProtectionPlanId:
    Value: !GetAtt GuardDutyMalwareProtectionPlan.MalwareProtectionPlanId
  ReleaseBucketName:
    Value: !Ref ReleaseBucket
  ReceiptBucketName:
    Value: !Ref ReceiptBucket
  ScanResultQueueUrl:
    Value: !Ref ScanResultQueue
  ScanResultDLQUrl:
    Value: !Ref ScanResultDLQ
  DecisionTableName:
    Value: !Ref DecisionTable
  SecurityAlertTopicArn:
    Value: !Ref SecurityAlertTopic
````

### FILE: `aws_guardduty_immutable_release_gate/provider-profile.template.json`
```yaml
block_id: "AWS-GUARDDUTY-IMMUTABLE-RELEASE-GATE:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-guardduty-immutable-release-gate/provider-profile.template.json"
license: "LicenseRef-Workspace-Owner"
sha256: "b723a187807051daca364d337e518be78cf305f3f24727d0ee13b3843fe1a2ce"
variables: []
secrets_allowed: false
```
````json
{
  "profile_version": 1,
  "approve_live_effects": false,
  "approve_guardduty_costs": false,
  "aws_account_id": "",
  "aws_region": "",
  "stack_name": "",
  "tenant_id": "",
  "quarantine_bucket_name": "",
  "source_kms_key_arn": "",
  "release_bucket_name": "",
  "receipt_bucket_name": "",
  "security_alert_email": "",
  "retain_days": 365,
  "dedupe_days": 400,
  "max_object_bytes": 52428800,
  "reserved_concurrency": 5,
  "cost_owner": "",
  "security_owner": "",
  "replay_owner": "",
  "quarantine_read_isolation_evidence": "",
  "clean_canary_evidence": "",
  "malware_canary_evidence": "",
  "unsupported_canary_evidence": "",
  "sns_confirmation_evidence": "",
  "dlq_redrive_evidence": "",
  "rollback_evidence": ""
}
````

### FILE: `aws_guardduty_immutable_release_gate/build.ps1`
```yaml
block_id: "AWS-GUARDDUTY-IMMUTABLE-RELEASE-GATE:build:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-guardduty-immutable-release-gate/build.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "65c866fba913939f9da5c9a0e9f4cbf36659ca7d6e77026c2aa3e5803bb9fa62"
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
$component = Join-Path $root 'aws_guardduty_immutable_release_gate'
$receiver = Join-Path $root 'aws_ses_immutable_email_receiver'
foreach ($path in @($component, $receiver)) {
  if (-not (Test-Path -LiteralPath $path -PathType Container)) { throw "required composed directory missing: $path" }
}
$output = [IO.Path]::GetFullPath($OutputRoot)
if (Test-Path -LiteralPath $output) { throw "output must not exist: $output" }
$parent = Split-Path -Parent $output
if (-not (Test-Path -LiteralPath $parent -PathType Container)) { throw "output parent missing: $parent" }
$temp = Join-Path $parent ('.' + [IO.Path]::GetFileName($output) + '.' + [Guid]::NewGuid().ToString('N'))
$functionDir = Join-Path $temp 'build/function'

try {
  New-Item -ItemType Directory -Path $functionDir | Out-Null
  Copy-Item -LiteralPath (Join-Path $component 'handler.py') -Destination $functionDir
  Copy-Item -LiteralPath (Join-Path $component 'AWS_MIT_NO_ATTRIBUTION_LICENSE.txt') -Destination $functionDir
  Copy-Item -LiteralPath (Join-Path $receiver 'requirements.txt') -Destination $functionDir
  Copy-Item -LiteralPath (Join-Path $component 'template.yaml') -Destination $temp

  $env:PYTHONDONTWRITEBYTECODE = '1'
  $previousPythonPath = $env:PYTHONPATH
  try {
    $env:PYTHONPATH = $component
    & $PythonExecutable (Join-Path $component 'test_handler.py') -v
    if ($LASTEXITCODE -ne 0) { throw 'GuardDuty release tests failed' }
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
    dependency_authority = 'aws_ses_immutable_email_receiver/requirements.txt'
    upstream_commit = 'fd9d2bdf46eb52bb336a7e7a837bd38ef73002e2'
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
  Write-Output "AWS_GUARDDUTY_RELEASE_BUILD_PASS output=$output files=$($receipt.files.Count) dependencies=$([bool]$InstallDependencies)"
} catch {
  if (Test-Path -LiteralPath $temp) { Remove-Item -LiteralPath $temp -Recurse -Force }
  throw
}
````

### FILE: `aws_guardduty_immutable_release_gate/deploy.ps1`
```yaml
block_id: "AWS-GUARDDUTY-IMMUTABLE-RELEASE-GATE:deploy:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-guardduty-immutable-release-gate/deploy.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "dbcc40302f6db2927406bab8b2ade069f597dd0d8156dd7117b5fc63d24f771e"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)][string] $ProfilePath,
  [Parameter(Mandatory = $true)][string] $BuildRoot,
  [switch] $ApproveLiveEffects,
  [switch] $ApproveGuardDutyCosts
)

$ErrorActionPreference = 'Stop'
$profile = Get-Content -LiteralPath (Resolve-Path -LiteralPath $ProfilePath).Path -Raw | ConvertFrom-Json
if (-not $ApproveLiveEffects -or $profile.approve_live_effects -ne $true) { throw 'live AWS effects require both -ApproveLiveEffects and approve_live_effects=true' }
if (-not $ApproveGuardDutyCosts -or $profile.approve_guardduty_costs -ne $true) { throw 'GuardDuty costs require both -ApproveGuardDutyCosts and approve_guardduty_costs=true' }
foreach ($field in @('aws_account_id','aws_region','stack_name','tenant_id','quarantine_bucket_name','source_kms_key_arn','release_bucket_name','receipt_bucket_name','security_alert_email','cost_owner','security_owner','replay_owner','quarantine_read_isolation_evidence')) {
  if ([string]::IsNullOrWhiteSpace([string]$profile.$field)) { throw "profile field is required: $field" }
}
if ([string]$profile.aws_account_id -notmatch '^[0-9]{12}$') { throw 'aws_account_id must be 12 digits' }
if ([int]$profile.dedupe_days -lt [int]$profile.retain_days) { throw 'dedupe_days must be at least retain_days' }
if ([int]$profile.reserved_concurrency -lt 2) { throw 'reserved_concurrency must be at least 2 for SQS scaling' }
$build = (Resolve-Path -LiteralPath $BuildRoot).Path
$template = Join-Path $build 'template.yaml'
if (-not (Test-Path -LiteralPath $template -PathType Leaf)) { throw 'built template missing' }
foreach ($tool in @('aws','sam')) {
  if (-not (Get-Command $tool -ErrorAction SilentlyContinue)) { throw "$tool CLI is required" }
}

$identity = (& aws sts get-caller-identity --region $profile.aws_region --output json | ConvertFrom-Json)
if ($LASTEXITCODE -ne 0 -or [string]$identity.Account -ne [string]$profile.aws_account_id) { throw 'AWS caller account mismatch' }
$versioning = (& aws s3api get-bucket-versioning --bucket $profile.quarantine_bucket_name --expected-bucket-owner $profile.aws_account_id --region $profile.aws_region --output json | ConvertFrom-Json)
if ($LASTEXITCODE -ne 0 -or $versioning.Status -ne 'Enabled') { throw 'quarantine bucket versioning is not enabled' }
$encryption = (& aws s3api get-bucket-encryption --bucket $profile.quarantine_bucket_name --expected-bucket-owner $profile.aws_account_id --region $profile.aws_region --output json | ConvertFrom-Json)
if ($LASTEXITCODE -ne 0) { throw 'quarantine bucket encryption could not be verified' }
$kmsRules = @($encryption.ServerSideEncryptionConfiguration.Rules | Where-Object { $_.ApplyServerSideEncryptionByDefault.SSEAlgorithm -eq 'aws:kms' -and $_.ApplyServerSideEncryptionByDefault.KMSMasterKeyID -eq $profile.source_kms_key_arn })
if ($kmsRules.Count -ne 1) { throw 'quarantine bucket exact KMS authority mismatch' }

$samBuild = Join-Path ([IO.Path]::GetTempPath()) ('elite-guardduty-sam-' + [Guid]::NewGuid().ToString('N'))
try {
  & sam build --template-file $template --build-dir $samBuild --cached:$false
  if ($LASTEXITCODE -ne 0) { throw 'sam build failed' }
  $parameters = @(
    "TenantId=$($profile.tenant_id)", "QuarantineBucketName=$($profile.quarantine_bucket_name)",
    "SourceKmsKeyArn=$($profile.source_kms_key_arn)", "ReleaseBucketName=$($profile.release_bucket_name)",
    "ReceiptBucketName=$($profile.receipt_bucket_name)", "SecurityAlertEmail=$($profile.security_alert_email)",
    "RetainDays=$($profile.retain_days)", "DedupeDays=$($profile.dedupe_days)",
    "MaxObjectBytes=$($profile.max_object_bytes)", "ReservedConcurrency=$($profile.reserved_concurrency)"
  )
  & sam deploy --template-file (Join-Path $samBuild 'template.yaml') --stack-name $profile.stack_name --region $profile.aws_region --capabilities CAPABILITY_IAM --resolve-s3 --no-confirm-changeset --no-fail-on-empty-changeset --parameter-overrides $parameters
  if ($LASTEXITCODE -ne 0) { throw 'sam deploy failed' }
} finally {
  if (Test-Path -LiteralPath $samBuild) { Remove-Item -LiteralPath $samBuild -Recurse -Force }
}

$outputs = @(& aws cloudformation describe-stacks --stack-name $profile.stack_name --region $profile.aws_region --query 'Stacks[0].Outputs' --output json | ConvertFrom-Json)
if ($LASTEXITCODE -ne 0) { throw 'stack outputs unavailable' }
$planId = [string](($outputs | Where-Object OutputKey -eq 'MalwareProtectionPlanId').OutputValue)
if ([string]::IsNullOrWhiteSpace($planId)) { throw 'malware protection plan id output missing' }
$plan = (& aws guardduty get-malware-protection-plan --malware-protection-plan-id $planId --region $profile.aws_region --output json | ConvertFrom-Json)
if ($LASTEXITCODE -ne 0 -or $plan.Status -ne 'ACTIVE') { throw 'GuardDuty malware protection plan is not ACTIVE' }
Write-Output "AWS_GUARDDUTY_RELEASE_DEPLOY_PASS stack=$($profile.stack_name) plan_id=$planId status=$($plan.Status)"
Write-Warning 'Deployment is not admission. Confirm SNS and execute clean, malware, unsupported, duplicate, DLQ/redrive and rollback canaries; record evidence in the provider profile before downstream enablement.'
````

### FILE: `aws_guardduty_immutable_release_gate/verify_contract.ps1`
```yaml
block_id: "AWS-GUARDDUTY-IMMUTABLE-RELEASE-GATE:verify:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-guardduty-immutable-release-gate/verify_contract.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "6ec90e68fdfe84e4138ce63ac95c5f1acf2b17b0d229a02db70db2df45fa611a"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [string] $PythonExecutable = 'python',
  [string] $CfnLintSitePackages = ''
)

$ErrorActionPreference = 'Stop'
$root = $PSScriptRoot
$handler = Join-Path $root 'handler.py'
$tests = Join-Path $root 'test_handler.py'
$template = Join-Path $root 'template.yaml'
$license = Join-Path $root 'AWS_MIT_NO_ATTRIBUTION_LICENSE.txt'
foreach ($path in @($handler,$tests,$template,$license,'README.md','provider-profile.template.json','build.ps1','deploy.ps1')) {
  $full = if ([IO.Path]::IsPathRooted($path)) { $path } else { Join-Path $root $path }
  if (-not (Test-Path -LiteralPath $full -PathType Leaf)) { throw "required file missing: $full" }
}
$licenseHash = (Get-FileHash -LiteralPath $license -Algorithm SHA256).Hash.ToLowerInvariant()
if ($licenseHash -ne '5025c7bbedbc8da868b6dce2ab689225b3f33c43b3f03368b0a331809b6ada26') { throw 'normalized AWS MIT-0 notice hash mismatch' }

$env:PYTHONDONTWRITEBYTECODE = '1'
$previousPythonPath = $env:PYTHONPATH
try {
  $env:PYTHONPATH = $root
  & $PythonExecutable $tests -v
  if ($LASTEXITCODE -ne 0) { throw 'unit/contract tests failed' }
  & $PythonExecutable -m py_compile $handler $tests
  if ($LASTEXITCODE -ne 0) { throw 'Python compile failed' }
} finally {
  $env:PYTHONPATH = $previousPythonPath
}

$handlerText = Get-Content -LiteralPath $handler -Raw
foreach ($token in @('VersionId=scan.version_id','CopySourceIfMatch','IfNoneMatch="*"','GuardDutyMalwareScanStatus','EXPECTED_PLAN_ARN','IDEMPOTENCY_OWNER_FENCE','ObjectLockRetainUntilDate','elite-aws-guardduty-release-receipt/v1','automatic_business_persistence_authorized','batchItemFailures','message_id_sha256')) {
  if (-not $handlerText.Contains($token, [StringComparison]::Ordinal)) { throw "handler contract token missing: $token" }
}
$templateText = Get-Content -LiteralPath $template -Raw
foreach ($token in @('AWS::GuardDuty::MalwareProtectionPlan','Tagging: {Status: ENABLED}','ObjectPrefixes','ReportBatchItemFailures','MaximumConcurrency','SqsManagedSseEnabled: true','PointInTimeRecoveryEnabled: true','DeletionProtectionEnabled: true','ObjectLockEnabled: true','s3:ExistingObjectTag/GuardDutyMalwareScanStatus','NO_THREATS_FOUND','DeletionPolicy: Retain')) {
  if (-not $templateText.Contains($token, [StringComparison]::Ordinal)) { throw "template contract token missing: $token" }
}
if ($templateText.Contains('DeletionPolicy: Delete', [StringComparison]::Ordinal) -or $templateText.Contains('RemovalPolicy.DESTROY', [StringComparison]::Ordinal)) { throw 'destructive infrastructure policy detected' }

if (-not [string]::IsNullOrWhiteSpace($CfnLintSitePackages)) {
  $site = (Resolve-Path -LiteralPath $CfnLintSitePackages).Path
  $previousPythonPath = $env:PYTHONPATH
  try {
    $env:PYTHONPATH = $site
    & $PythonExecutable -c 'from cfnlint.runner import main; main()' $template
    if ($LASTEXITCODE -ne 0) { throw 'cfn-lint failed' }
  } finally {
    $env:PYTHONPATH = $previousPythonPath
  }
}

$profile = Get-Content -LiteralPath (Join-Path $root 'provider-profile.template.json') -Raw | ConvertFrom-Json
if ($profile.approve_live_effects -ne $false -or $profile.approve_guardduty_costs -ne $false) { throw 'provider template must fail closed' }
Write-Output 'AWS_GUARDDUTY_IMMUTABLE_RELEASE_PASS tests=16 official_source=aws-samples/guardduty-malware-protection@fd9d2bdf46eb52bb336a7e7a837bd38ef73002e2 storage_authorized=0'
````

### FILE: `aws_guardduty_immutable_release_gate/AWS_MIT_NO_ATTRIBUTION_LICENSE.txt`
```yaml
block_id: "AWS-GUARDDUTY-IMMUTABLE-RELEASE-GATE:aws-notice:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/aws-samples/guardduty-malware-protection/blob/fd9d2bdf46eb52bb336a7e7a837bd38ef73002e2/LICENSE; text preserved, newline/terminal blank normalized by Markdown materialization"
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

| Variable | Tipo/default seguro | Validación | Secreto | Efecto |
|---|---|---|---|---|
| account/region/plan | exactos/sin default | 12 dígitos, región, ARN output | no | autoridad del evento |
| tenant/quarantine prefix | exactos/sin default | prefijo del receptor | no | aislamiento |
| buckets/KMS | exactos/sin default | owner, versioning, SSE-KMS | no | fuentes/destinos |
| retention/dedupe/limit | 365/400/50 MiB | rangos y dedupe≥retención | no | evidencia/recursos |
| approvals/owners | false/vacío | doble costo+efecto y owners | no | habilita LIVE |
| canary/isolation records | vacío | requerido por admission, no por build | no | habilita downstream |

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| AWS sample | commit `fd9d2bdf46eb52bb336a7e7a837bd38ef73002e2` | base declarada de plan/SQS/Lambda | MIT-0 | reference/adaptation | aws-samples |
| Boto3 graph | `1.43.83`, 7 wheels hashados por receiver | APIs S3/DDB/SNS | Apache-2.0 y transitivas | runtime | PyPI/AWS SDK |
| Python | `3.14` | Lambda | PSF | runtime | AWS Lambda |
| cfn-lint | `1.55.1` | validación fechada | MIT-0 | audit | AWS CloudFormation |

## 8. Apply order

Componer acquisition+receiver+MIME+este pack; completar readiness/profile; verificar; construir en destino nuevo con dependencia hashada; validar cuenta/bucket/KMS; desplegar sólo con doble aprobación; confirmar plan/SNS; ejecutar canaries clean/malware/unsupported/duplicate/DLQ/redrive/rollback; registrar evidence; habilitar extracción únicamente para receipts limpios. Workspace existente: compositor rechaza colisión. Rollback: deshabilitar consumidor/rule y preservar evidencia; no borrar automáticamente.

## 9. Verification

`verify_contract.ps1` exige 16 tests, compile, tokens de seguridad, notice hashado y profile fail-closed; con site-packages opcional ejecuta cfn-lint 1.55.1. `build.ps1 -InstallDependencies` instala 7 wheels con `--require-hashes`, importa Boto3/Botocore exactos y emite 2.191 hashes. Live requiere plan ACTIVE y canaries target; el PASS offline no los reemplaza.

## 10. Reconstruction evidence

Reconstruido el 2026-08-28 en Windows/PowerShell 7/Python con 9 archivos. Tests 16/16, py_compile, cfn-lint 1.55.1 exit 0 y build hash-locked 2.191 archivos PASS. Source AWS auditado en commit fijado; su template exacto terminó cfn-lint exit 4 por warnings y su lane directa permanece rechazada. Evidencia detallada se registra en `reconstruction_evidence/` después del gate global.
