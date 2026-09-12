# AWS IDP Immutable Evaluation Handoff

## 1. Metadata

```yaml
pack_id: "AWS-IDP-IMMUTABLE-EVALUATION-HANDOFF"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un hook postprocessing síncrono para AWS GenAI IDP v0.6.5 que carga el Document mediante el helper oficial, valida consistentemente la cobertura HITL Completed contra tracking, fija evidencia de review/salida por hash/versión/Object Lock y entrega una solicitud FIFO a evaluación estricta sin autorizar persistencia empresarial."
stacks: ["AWS GenAI IDP Accelerator v0.6.5", "idp_common 0.6.5", "AWS Lambda Python 3.13", "AWS SAM", "Amazon DynamoDB", "Amazon S3 Object Lock", "Amazon SQS FIFO", "AWS KMS"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "AWS-GUARDDUTY-MAGIKA-IDP-DISPATCH-GATE 0.1.x", "STRICT-DOCUMENT-FIELD-EVALUATION-GATE 0.1.x"]
incompatible_with: ["AWS IDP distinto del commit fijado", "hook onError continue", "allowDocumentUpdate true", "configuración activa móvil", "clases sin allow-list", "bucket sin Versioning/Object Lock/KMS", "persistencia empresarial automática"]
license_expression: "LicenseRef-Workspace-Owner AND MIT-0"
upstream_sources: ["https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/tree/1b5fd74454e593de233342a02ee911af8ee38359", "https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/blob/1b5fd74454e593de233342a02ee911af8ee38359/docs/feature-platform.md", "https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/blob/1b5fd74454e593de233342a02ee911af8ee38359/lib/idp_common_pkg/idp_common/hooks/__init__.py", "https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/blob/1b5fd74454e593de233342a02ee911af8ee38359/src/lambda/complete_section_review/index.py", "https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-lock.html", "https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/fifo-queues-exactly-once-processing.html"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use sólo como `postprocessing` del AWS GenAI IDP Accelerator v0.6.5 exacto, después de su evaluación y antes de cualquier consumidor empresarial. Requiere el `IDPCommonBaseLayerArn` publicado por el stack oficial, configuración IDP fijada y evaluada, clases documentales admitidas, HITL definido, cuenta/región/state machine exactos, KMS/Object Lock/SQS FIFO y owners/evidencia completos. No convierte un PASS offline en autorización AWS ni promete exactitud universal o exactly-once empresarial.

## 3. Architecture contract

El handler es glue `ADAPTED`: reutiliza `idp_common.hooks.load_hook_document` y el contrato de tracking HITL oficiales, pero AWS no es autor del hardening local. Valida punto/feature/args/config/ejecución, rechaza errores e issues y sólo admite clases configuradas. `PendingReview`/`InProgress` produce receipt diferido. `Completed` exige lectura DynamoDB consistente, cobertura exacta sin pending/skipped/duplicados, reviewer conocido y timestamp zonificado; la identidad se reduce a SHA-256 y se guarda evidencia Object-Locked separada y autocontenida, ligada a document SHA/config/schema/ejecución. `Skipped` queda explícitamente `SKIPPED_NOT_REVIEWED`. Luego produce snapshot canónico content-addressed, mensaje FIFO y receipt comprometido al final. Un replay con receipt final no reenvía. El consumidor todavía debe deduplicar por `evaluationId`, y toda salida fija `automaticBusinessPersistenceAuthorized=false`.

## 4. Exact file manifest

```text
CREATE aws_idp_immutable_evaluation_handoff/README.md
CREATE aws_idp_immutable_evaluation_handoff/handler.py
CREATE aws_idp_immutable_evaluation_handoff/test_handler.py
CREATE aws_idp_immutable_evaluation_handoff/hook-config.template.yaml
CREATE aws_idp_immutable_evaluation_handoff/provider-profile.template.json
CREATE aws_idp_immutable_evaluation_handoff/template.yaml
CREATE aws_idp_immutable_evaluation_handoff/deploy.ps1
CREATE aws_idp_immutable_evaluation_handoff/verify_contract.ps1
```

## 5. Materialization blocks
### FILE: `aws_idp_immutable_evaluation_handoff/README.md`
```yaml
block_id: "AWS-IDP-IMMUTABLE-EVALUATION-HANDOFF:readme-md:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-immutable-evaluation-handoff/README.md"
license: "LicenseRef-Workspace-Owner"
sha256: "c20429a6a02cafeabc4bc945b2532ea5909e24a7f6e227ca876ebbf30e78cb80"
variables: []
secrets_allowed: false
```
````markdown
# AWS IDP immutable evaluation handoff

This component is a synchronous `postprocessing` pipeline hook for the official AWS GenAI IDP Accelerator v0.6.5. It loads the exact inline or compressed `Document` through the official `idp_common.hooks.load_hook_document` helper supplied by the accelerator's `IDPCommonBaseLayerArn` output. It does not reimplement OCR, classification, extraction, confidence, rule validation, evaluation or HITL.

The hook fails closed unless the hook point, feature id, Step Functions execution ARN, configuration version, tenant argument, terminal document state and classification allow-list all match the project profile. Pending or in-progress HITL produces only an immutable deferred receipt. For `Completed`, the hook consistently reads the official tracking row and requires exact, non-skipped, non-duplicated section coverage, known reviewer subjects and timezone-aware timestamps. It stores only SHA-256 reviewer identifiers in a separate Object Lock review-evidence object; raw subjects and emails are never copied into the handoff artifacts. `Skipped` is explicitly labelled `SKIPPED_NOT_REVIEWED` and cannot masquerade as human review.

A terminal document is canonicalized, SHA-256 addressed and stored in a versioned SSE-KMS Object Lock bucket. Review evidence, when required, is bound into the evaluation id. The component then emits a FIFO evaluation request and commits an immutable dispatch receipt last. Replays reconcile the receipt and do not resend after a committed dispatch.

The handler is `ADAPTED` integration code built against AWS's published pipeline-hook and review-tracking contracts. AWS does not claim this local hardening. The hook configuration must use `onError: fail` and `allowDocumentUpdate: false`; the latter prevents this delivery gate from mutating the official document. A returned handoff result always fixes `automaticBusinessPersistenceAuthorized` to `false`. This evidence detects a coherent tracking snapshot; it is not a cryptographic signature by AWS and does not repair concurrent mutations inside the upstream review handler. Downstream schema evaluation, representative-corpus acceptance, human approval policy and transactional business persistence remain separate gates.

The FIFO queue provides concurrent duplicate suppression, while the immutable dispatch receipt provides durable replay reconciliation. A consumer must still key idempotency on `evaluationId`: FIFO deduplication alone is not a permanent exactly-once guarantee. Live AWS deployment, KMS policies, Object Lock retention, alarms, DLQ redrive, concurrency, costs and restoration require target evidence and owner approval.
````

### FILE: `aws_idp_immutable_evaluation_handoff/handler.py`
```yaml
block_id: "AWS-IDP-IMMUTABLE-EVALUATION-HANDOFF:handler-py:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/tree/1b5fd74454e593de233342a02ee911af8ee38359; integration hardening against the official v0.6.5 pipeline-hook and idp_common contracts"
license: "LicenseRef-Workspace-Owner AND MIT-0"
sha256: "f3e27bd5f60b5665e8f5e5ba48051d346408830ebb3e813b7a31fa4d42e6b7c4"
variables: []
secrets_allowed: false
```
````python
"""Fail-closed AWS GenAI IDP v0.6.5 postprocessing handoff.

ADAPTED integration code. The input contract and ``load_hook_document`` helper
are published by AWS under MIT-0. This file is not represented as verbatim AWS
source and never authorizes automatic business persistence.
"""

from __future__ import annotations

import base64
import hashlib
import json
import os
import re
from dataclasses import dataclass
from datetime import datetime, timedelta, timezone
from typing import Any

import boto3
from botocore.exceptions import ClientError
from idp_common.hooks import load_hook_document


HANDOFF_SCHEMA = "elite-aws-idp-evaluation-handoff/v1"
DEFERRED_SCHEMA = "elite-aws-idp-hitl-deferred/v1"
REVIEW_SCHEMA = "elite-aws-idp-review-evidence/v1"
SAFE_TOKEN = re.compile(r"^[A-Za-z0-9][A-Za-z0-9._:/+=,@-]{0,511}$")
TERMINAL_HITL = {None, "Completed", "Skipped"}
DEFERRED_HITL = {"PendingReview", "InProgress"}


class Rejected(ValueError):
    """Permanent contract rejection; the configured pipeline hook must fail."""


@dataclass(frozen=True)
class Config:
    account_id: str
    region: str
    feature_id: str
    execution_name_prefix: str
    tenant_id: str
    config_version: str
    schema_profile: str
    result_bucket: str
    result_kms_key_arn: str
    evaluation_queue_url: str
    working_bucket: str
    tracking_table_name: str
    allowed_classes: frozenset[str]
    retain_days: int
    max_canonical_bytes: int

    @staticmethod
    def from_env() -> "Config":
        try:
            classes = json.loads(os.environ.get("ALLOWED_CLASSES_JSON", "[]"))
        except json.JSONDecodeError as error:
            raise Rejected("ALLOWED_CLASSES_JSON is invalid") from error
        config = Config(
            account_id=os.environ.get("EXPECTED_ACCOUNT_ID", ""),
            region=os.environ.get("EXPECTED_REGION", ""),
            feature_id=os.environ.get("EXPECTED_FEATURE_ID", ""),
            execution_name_prefix=os.environ.get("EXECUTION_NAME_PREFIX", ""),
            tenant_id=os.environ.get("TENANT_ID", ""),
            config_version=os.environ.get("IDP_CONFIG_VERSION", ""),
            schema_profile=os.environ.get("SCHEMA_PROFILE", ""),
            result_bucket=os.environ.get("RESULT_BUCKET", ""),
            result_kms_key_arn=os.environ.get("RESULT_KMS_KEY_ARN", ""),
            evaluation_queue_url=os.environ.get("EVALUATION_QUEUE_URL", ""),
            working_bucket=os.environ.get("WORKING_BUCKET", ""),
            tracking_table_name=os.environ.get("TRACKING_TABLE_NAME", ""),
            allowed_classes=frozenset(classes) if isinstance(classes, list) else frozenset(),
            retain_days=int(os.environ.get("RETAIN_DAYS", "365")),
            max_canonical_bytes=int(os.environ.get("MAX_CANONICAL_BYTES", "20971520")),
        )
        config.validate()
        return config

    def validate(self) -> None:
        strings = (
            self.feature_id, self.execution_name_prefix, self.tenant_id,
            self.config_version, self.schema_profile, self.result_bucket,
            self.result_kms_key_arn, self.evaluation_queue_url, self.working_bucket,
            self.tracking_table_name,
        )
        if not re.fullmatch(r"[0-9]{12}", self.account_id):
            raise Rejected("EXPECTED_ACCOUNT_ID is invalid")
        if not re.fullmatch(r"[a-z]{2}(?:-gov)?-[a-z]+-[0-9]", self.region):
            raise Rejected("EXPECTED_REGION is invalid")
        if any(not value or not SAFE_TOKEN.fullmatch(value) for value in strings[:5]):
            raise Rejected("identity configuration is invalid")
        if any(not value for value in strings[5:]):
            raise Rejected("storage configuration is missing")
        if not self.evaluation_queue_url.endswith(".fifo"):
            raise Rejected("evaluation queue must be FIFO")
        if not re.fullmatch(r"[A-Za-z0-9_.-]{3,255}", self.tracking_table_name):
            raise Rejected("tracking table name is invalid")
        if not self.allowed_classes or any(not isinstance(x, str) or not SAFE_TOKEN.fullmatch(x) for x in self.allowed_classes):
            raise Rejected("allowed classes are invalid")
        if not (1 <= self.retain_days <= 36500 and 1024 <= self.max_canonical_bytes <= 50 * 1024 * 1024):
            raise Rejected("limits are invalid")


@dataclass(frozen=True)
class Clients:
    s3: Any
    sqs: Any
    tracking: Any


def _clients(cfg: Config) -> Clients:
    return Clients(
        boto3.client("s3"), boto3.client("sqs"),
        boto3.resource("dynamodb").Table(cfg.tracking_table_name),
    )


def _canonical(value: Any) -> bytes:
    try:
        return json.dumps(value, sort_keys=True, separators=(",", ":"), ensure_ascii=False, allow_nan=False).encode("utf-8")
    except (TypeError, ValueError) as error:
        raise Rejected("document is not canonical JSON") from error


def _sha(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def _checksum(data: bytes) -> str:
    return base64.b64encode(hashlib.sha256(data).digest()).decode("ascii")


def _is_precondition(error: ClientError) -> bool:
    code = str(error.response.get("Error", {}).get("Code", ""))
    return code in {"PreconditionFailed", "412", "ConditionalRequestConflict", "409"}


def _execution_arn(event: dict[str, Any], cfg: Config) -> str:
    value = event.get("executionArn")
    prefix = f"arn:aws:states:{cfg.region}:{cfg.account_id}:execution:{cfg.execution_name_prefix}"
    if not isinstance(value, str) or not value.startswith(prefix) or not SAFE_TOKEN.fullmatch(value):
        raise Rejected("execution ARN is not admitted")
    return value


def _validate_event(event: dict[str, Any], cfg: Config) -> str:
    if event.get("hookPoint") != "postprocessing":
        raise Rejected("only postprocessing is admitted")
    if event.get("featureId") != cfg.feature_id:
        raise Rejected("feature id mismatch")
    args = event.get("argsMap")
    expected = {"tenantId": cfg.tenant_id, "configVersion": cfg.config_version, "schemaProfile": cfg.schema_profile}
    if not isinstance(args, dict) or args != expected:
        raise Rejected("hook args do not match the project profile")
    return _execution_arn(event, cfg)


def _validate_document(document: dict[str, Any], cfg: Config) -> tuple[str, Any, tuple[str, ...]]:
    doc_id = document.get("id") or document.get("document_id")
    if not isinstance(doc_id, str) or not doc_id or len(doc_id.encode("utf-8")) > 1024:
        raise Rejected("document id is invalid")
    if document.get("config_version") != cfg.config_version:
        raise Rejected("document config version mismatch")
    if document.get("evaluation_status") in {"RUNNING", "FAILED"}:
        raise Rejected("official evaluation is not in an admissible state")
    if document.get("status") == "REDACTED_SUPERSEDED":
        raise Rejected("superseded unredacted document is forbidden")
    if document.get("errors") or document.get("processing_issues"):
        raise Rejected("document contains processing errors or issues")
    hitl = document.get("hitl_status")
    if hitl not in TERMINAL_HITL | DEFERRED_HITL:
        raise Rejected("unknown HITL state")
    sections = document.get("sections")
    if not isinstance(sections, list) or not sections:
        raise Rejected("document sections are missing")
    identities: set[str] = set()
    classes: list[str] = []
    for section in sections:
        if not isinstance(section, dict):
            raise Rejected("section is invalid")
        sid = section.get("section_id")
        classification = section.get("classification")
        if not isinstance(sid, str) or not sid or sid in identities:
            raise Rejected("section identity is invalid or duplicated")
        if classification not in cfg.allowed_classes:
            raise Rejected("section classification is not admitted")
        if section.get("processing_issues"):
            raise Rejected("section contains processing issues")
        identities.add(sid)
        classes.append(str(classification))
    return doc_id, hitl, tuple(classes)


def _head(clients: Clients, cfg: Config, key: str) -> dict[str, Any] | None:
    try:
        return clients.s3.head_object(Bucket=cfg.result_bucket, Key=key, ChecksumMode="ENABLED")
    except ClientError as error:
        if str(error.response.get("Error", {}).get("Code", "")) in {"404", "NoSuchKey", "NotFound"}:
            return None
        raise


def _put_immutable(clients: Clients, cfg: Config, key: str, body: bytes, sha256: str, schema: str = HANDOFF_SCHEMA) -> dict[str, Any]:
    retain_until = datetime.now(timezone.utc) + timedelta(days=cfg.retain_days)
    try:
        result = clients.s3.put_object(
            Bucket=cfg.result_bucket, Key=key, Body=body, ContentType="application/json",
            ChecksumAlgorithm="SHA256", ChecksumSHA256=_checksum(body), IfNoneMatch="*",
            ServerSideEncryption="aws:kms", SSEKMSKeyId=cfg.result_kms_key_arn,
            ObjectLockMode="COMPLIANCE", ObjectLockRetainUntilDate=retain_until,
            Metadata={"sha256": sha256, "schema": schema},
        )
        if not result.get("VersionId"):
            raise Rejected("immutable put returned no VersionId")
        return result
    except ClientError as error:
        if not _is_precondition(error):
            raise
        existing = _head(clients, cfg, key)
        if not existing:
            raise Rejected("conditional collision cannot be reconciled") from error
        metadata = existing.get("Metadata") or {}
        if metadata.get("sha256") != sha256 or existing.get("ContentLength") != len(body):
            raise Rejected("immutable object is divergent") from error
        return existing


def _review_evidence(clients: Clients, cfg: Config, doc_id: str, document: dict[str, Any], document_sha: str, execution_arn: str) -> dict[str, Any]:
    response = clients.tracking.get_item(
        Key={"PK": f"doc#{doc_id}", "SK": "none"}, ConsistentRead=True,
    )
    item = response.get("Item")
    if not isinstance(item, dict) or item.get("HITLCompleted") is not True:
        raise Rejected("completed HITL has no completed tracking record")
    if item.get("HITLStatus") == "Review Skipped":
        raise Rejected("skipped review cannot produce completed review evidence")
    if item.get("HITLSectionsPending") not in (None, []):
        raise Rejected("completed HITL still has pending sections")
    if item.get("HITLSectionsSkipped") not in (None, []):
        raise Rejected("completed HITL contains skipped sections")

    section_ids = {section["section_id"] for section in document["sections"]}
    document_completed = document.get("hitl_sections_completed")
    tracking_completed = item.get("HITLSectionsCompleted")
    if not isinstance(document_completed, list) or set(document_completed) != section_ids or len(document_completed) != len(section_ids):
        raise Rejected("document review coverage is incomplete or duplicated")
    if not isinstance(tracking_completed, list) or set(tracking_completed) != section_ids or len(tracking_completed) != len(section_ids):
        raise Rejected("tracking review coverage is incomplete or duplicated")

    history = item.get("HITLReviewHistory")
    if not isinstance(history, list) or len(history) != len(section_ids):
        raise Rejected("review history is missing or ambiguous")
    normalized: list[dict[str, str]] = []
    seen: set[str] = set()
    for record in history:
        if not isinstance(record, dict) or record.get("action") == "skip_all":
            raise Rejected("review history contains a non-review action")
        section_id = record.get("sectionId")
        subject = record.get("reviewedBy")
        reviewed_at = record.get("reviewedAt")
        email = record.get("reviewedByEmail", "")
        if section_id not in section_ids or section_id in seen:
            raise Rejected("review history section is unknown or duplicated")
        if not isinstance(subject, str) or not subject.strip() or subject.strip().lower() == "unknown":
            raise Rejected("reviewer subject is missing")
        if not isinstance(email, str) or not isinstance(reviewed_at, str):
            raise Rejected("review identity or timestamp is invalid")
        try:
            instant = datetime.fromisoformat(reviewed_at.replace("Z", "+00:00"))
        except ValueError as error:
            raise Rejected("review timestamp is invalid") from error
        if instant.tzinfo is None:
            raise Rejected("review timestamp has no timezone")
        seen.add(section_id)
        proof = {
            "sectionId": section_id,
            "reviewedAt": instant.astimezone(timezone.utc).isoformat(),
            "reviewerSubjectSha256": _sha(subject.strip().encode("utf-8")),
        }
        if email.strip():
            proof["reviewerEmailSha256"] = _sha(email.strip().lower().encode("utf-8"))
        normalized.append(proof)
    if seen != section_ids:
        raise Rejected("review history does not cover every section")
    normalized.sort(key=lambda value: value["sectionId"])
    tracking_state = {
        "HITLCompleted": True,
        "HITLStatus": item.get("HITLStatus"),
        "HITLSectionsCompleted": sorted(tracking_completed),
        "HITLSectionsPending": [],
        "HITLSectionsSkipped": [],
        "reviews": normalized,
    }
    return {
        "schema": REVIEW_SCHEMA,
        "tenantId": cfg.tenant_id,
        "documentIdSha256": _sha(doc_id.encode("utf-8")),
        "documentSha256": document_sha,
        "configVersion": cfg.config_version,
        "schemaProfile": cfg.schema_profile,
        "executionArnSha256": _sha(execution_arn.encode("utf-8")),
        "trackingStateSha256": _sha(_canonical(tracking_state)),
        "sections": normalized,
        "sectionCount": len(section_ids),
        "reviewAuthority": "AWS_IDP_HUMAN_REVIEW_COMPLETED",
        "automaticBusinessPersistenceAuthorized": False,
    }


def _receipt_document(schema: str, cfg: Config, doc_id: str, hitl: Any, classes: tuple[str, ...], document_sha: str, execution_arn: str, review_authority: str, review_sha: str | None) -> dict[str, Any]:
    return {
        "schema": schema,
        "tenantId": cfg.tenant_id,
        "documentIdSha256": _sha(doc_id.encode("utf-8")),
        "documentSha256": document_sha,
        "configVersion": cfg.config_version,
        "schemaProfile": cfg.schema_profile,
        "hitlStatus": hitl,
        "classifications": list(classes),
        "executionArnSha256": _sha(execution_arn.encode("utf-8")),
        "reviewAuthority": review_authority,
        "reviewEvidenceSha256": review_sha,
        "automaticBusinessPersistenceAuthorized": False,
    }


def lambda_handler(event: dict[str, Any], _context: Any, *, clients: Clients | None = None, cfg: Config | None = None) -> dict[str, Any]:
    cfg = cfg or Config.from_env()
    clients = clients or _clients(cfg)
    execution_arn = _validate_event(event, cfg)
    model = load_hook_document(event, working_bucket=cfg.working_bucket)
    document = model.to_dict()
    if document.get("workflow_execution_arn") not in {None, execution_arn}:
        raise Rejected("document workflow execution mismatch")
    doc_id, hitl, classes = _validate_document(document, cfg)
    canonical = _canonical(document)
    if len(canonical) > cfg.max_canonical_bytes:
        raise Rejected("canonical document exceeds configured bound")
    document_sha = _sha(canonical)
    identity = _sha((cfg.tenant_id + "\n" + doc_id).encode("utf-8"))

    if hitl in DEFERRED_HITL:
        deferred = _receipt_document(DEFERRED_SCHEMA, cfg, doc_id, hitl, classes, document_sha, execution_arn, "PENDING", None)
        deferred_body = _canonical(deferred)
        deferred_key = f"deferred/{cfg.tenant_id}/{identity}/{document_sha}.json"
        _put_immutable(clients, cfg, deferred_key, deferred_body, _sha(deferred_body), DEFERRED_SCHEMA)
        return {"schema": DEFERRED_SCHEMA, "documentSha256": document_sha, "evaluationDispatched": False, "automaticBusinessPersistenceAuthorized": False}

    review_authority = "NOT_TRIGGERED"
    review_sha = None
    review_ref = None
    if hitl == "Completed":
        review = _review_evidence(clients, cfg, doc_id, document, document_sha, execution_arn)
        review_body = _canonical(review)
        review_sha = _sha(review_body)
        review_key = f"reviews/{cfg.tenant_id}/{identity}/{review_sha}.json"
        stored_review = _put_immutable(clients, cfg, review_key, review_body, review_sha, REVIEW_SCHEMA)
        review_authority = review["reviewAuthority"]
        review_ref = {"bucket": cfg.result_bucket, "key": review_key, "versionId": stored_review.get("VersionId"), "sha256": review_sha}
    elif hitl == "Skipped":
        review_authority = "SKIPPED_NOT_REVIEWED"

    snapshot_key = f"snapshots/{cfg.tenant_id}/{identity}/{document_sha}.json"
    snapshot = _put_immutable(clients, cfg, snapshot_key, canonical, document_sha)
    evaluation_id = _sha((cfg.tenant_id + "\n" + cfg.config_version + "\n" + document_sha + "\n" + (review_sha or review_authority)).encode("utf-8"))
    dispatch_key = f"dispatch/{cfg.tenant_id}/{identity}/{evaluation_id}.json"
    if _head(clients, cfg, dispatch_key):
        return {"schema": HANDOFF_SCHEMA, "evaluationId": evaluation_id, "documentSha256": document_sha, "evaluationDispatched": True, "reconciled": True, "automaticBusinessPersistenceAuthorized": False}

    request = {
        "schema": HANDOFF_SCHEMA, "evaluationId": evaluation_id, "tenantId": cfg.tenant_id,
        "configVersion": cfg.config_version, "schemaProfile": cfg.schema_profile,
        "snapshot": {"bucket": cfg.result_bucket, "key": snapshot_key, "versionId": snapshot.get("VersionId"), "sha256": document_sha},
        "reviewAuthority": review_authority, "reviewEvidence": review_ref,
        "automaticBusinessPersistenceAuthorized": False,
    }
    clients.sqs.send_message(
        QueueUrl=cfg.evaluation_queue_url, MessageBody=_canonical(request).decode("utf-8"),
        MessageGroupId=identity, MessageDeduplicationId=evaluation_id,
    )
    dispatch = _receipt_document(HANDOFF_SCHEMA, cfg, doc_id, hitl, classes, document_sha, execution_arn, review_authority, review_sha)
    dispatch.update({"evaluationId": evaluation_id, "snapshotKey": snapshot_key, "snapshotVersionId": snapshot.get("VersionId"), "reviewEvidence": review_ref})
    dispatch_body = _canonical(dispatch)
    _put_immutable(clients, cfg, dispatch_key, dispatch_body, _sha(dispatch_body))
    return {"schema": HANDOFF_SCHEMA, "evaluationId": evaluation_id, "documentSha256": document_sha, "evaluationDispatched": True, "reconciled": False, "automaticBusinessPersistenceAuthorized": False}
````

### FILE: `aws_idp_immutable_evaluation_handoff/test_handler.py`
```yaml
block_id: "AWS-IDP-IMMUTABLE-EVALUATION-HANDOFF:test-handler-py:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/blob/1b5fd74454e593de233342a02ee911af8ee38359/lib/idp_common_pkg/tests/unit/test_hooks_helpers.py; local negative and replay tests"
license: "LicenseRef-Workspace-Owner AND MIT-0"
sha256: "f92aacbab93c3dcec3c30d188e7e941b55c1d76f644c01450db88d7f61acaf4f"
variables: []
secrets_allowed: false
```
````python
import json
import sys
import types
import unittest
from dataclasses import dataclass
from unittest.mock import patch


class FakeClientError(Exception):
    def __init__(self, code):
        self.response = {"Error": {"Code": code}}


botocore = types.ModuleType("botocore")
botocore_exceptions = types.ModuleType("botocore.exceptions")
botocore_exceptions.ClientError = FakeClientError
botocore.exceptions = botocore_exceptions
sys.modules.setdefault("botocore", botocore)
sys.modules.setdefault("botocore.exceptions", botocore_exceptions)
boto3 = types.ModuleType("boto3")
boto3.client = lambda _name: None
sys.modules.setdefault("boto3", boto3)


class FakeDocument:
    def __init__(self, value):
        self.value = value

    def to_dict(self):
        return self.value


hooks = types.ModuleType("idp_common.hooks")
hooks.load_hook_document = lambda event, working_bucket=None: FakeDocument(event["document"])
idp_common = types.ModuleType("idp_common")
idp_common.hooks = hooks
sys.modules.setdefault("idp_common", idp_common)
sys.modules.setdefault("idp_common.hooks", hooks)

import handler


@dataclass
class FakeS3:
    objects: dict

    def put_object(self, **kwargs):
        key = kwargs["Key"]
        if key in self.objects:
            raise FakeClientError("PreconditionFailed")
        body = kwargs["Body"]
        self.objects[key] = {
            "Body": body,
            "Metadata": kwargs["Metadata"],
            "ContentLength": len(body),
            "VersionId": "v-" + str(len(self.objects) + 1),
        }
        return {"VersionId": self.objects[key]["VersionId"], "ETag": "etag"}

    def head_object(self, **kwargs):
        key = kwargs["Key"]
        if key not in self.objects:
            raise FakeClientError("404")
        return self.objects[key]


@dataclass
class FakeSQS:
    messages: list
    fail: bool = False

    def send_message(self, **kwargs):
        if self.fail:
            raise RuntimeError("sqs unavailable")
        self.messages.append(kwargs)
        return {"MessageId": "m1"}


@dataclass
class FakeTracking:
    items: dict

    def get_item(self, **kwargs):
        self.last_request = kwargs
        item = self.items.get((kwargs["Key"]["PK"], kwargs["Key"]["SK"]))
        return {"Item": item} if item is not None else {}


def _messages(queue):
    return "\n".join(message["MessageBody"] for message in queue.messages).encode("utf-8")


def review_item():
    return {
        "HITLCompleted": True, "HITLStatus": "Review Completed",
        "HITLSectionsPending": [], "HITLSectionsSkipped": [],
        "HITLSectionsCompleted": ["1"],
        "HITLReviewHistory": [{
            "sectionId": "1", "reviewedBy": "reviewer-subject-7",
            "reviewedByEmail": "Reviewer@Example.com",
            "reviewedAt": "2026-08-28T12:30:00+00:00",
        }],
    }


def cfg():
    return handler.Config(
        account_id="123456789012", region="us-east-1",
        feature_id="elite-immutable-evaluation-handoff",
        execution_name_prefix="EliteIdp-",
        tenant_id="tenant-a", config_version="config-v7",
        schema_profile="invoice-v3", result_bucket="result",
        result_kms_key_arn="arn:aws:kms:us-east-1:123456789012:key/00000000-0000-0000-0000-000000000000",
        evaluation_queue_url="https://sqs.us-east-1.amazonaws.com/123456789012/evaluation.fifo",
        working_bucket="working", tracking_table_name="idp-tracking",
        allowed_classes=frozenset({"Invoice", "PackingList"}),
        retain_days=365, max_canonical_bytes=1024 * 1024,
    )


def event(hitl=None):
    doc = {
        "id": "imports/invoice-7.pdf", "config_version": "config-v7",
        "status": "EVALUATING", "evaluation_status": "COMPLETED",
        "workflow_execution_arn": "arn:aws:states:us-east-1:123456789012:execution:EliteIdp-main:exec-7",
        "sections": [{"section_id": "1", "classification": "Invoice", "attributes": {"total": "10.00"}}],
        "metering": {"OCR/textract": {"pages": 1}},
    }
    if hitl is not None:
        doc["hitl_status"] = hitl
    if hitl == "Completed":
        doc["hitl_sections_completed"] = ["1"]
    return {
        "hookPoint": "postprocessing", "featureId": "elite-immutable-evaluation-handoff",
        "executionArn": "arn:aws:states:us-east-1:123456789012:execution:EliteIdp-main:exec-7",
        "args": [], "argsMap": {"tenantId": "tenant-a", "configVersion": "config-v7", "schemaProfile": "invoice-v3"},
        "document": doc,
    }


class HandoffTests(unittest.TestCase):
    def setUp(self):
        self.s3 = FakeS3({})
        self.sqs = FakeSQS([])
        self.tracking = FakeTracking({("doc#imports/invoice-7.pdf", "none"): review_item()})
        self.clients = handler.Clients(self.s3, self.sqs, self.tracking)

    def call(self, value):
        return handler.lambda_handler(value, None, clients=self.clients, cfg=cfg())

    def test_terminal_document_snapshots_dispatches_and_commits(self):
        out = self.call(event())
        self.assertTrue(out["evaluationDispatched"])
        self.assertFalse(out["automaticBusinessPersistenceAuthorized"])
        self.assertEqual(1, len(self.sqs.messages))
        request = json.loads(self.sqs.messages[0]["MessageBody"])
        self.assertEqual(out["evaluationId"], request["evaluationId"])
        self.assertFalse(request["automaticBusinessPersistenceAuthorized"])
        self.assertTrue(any(key.startswith("snapshots/") for key in self.s3.objects))
        self.assertTrue(any(key.startswith("dispatch/") for key in self.s3.objects))

    def test_replay_reconciles_without_resending(self):
        first = self.call(event())
        second = self.call(event())
        self.assertEqual(first["evaluationId"], second["evaluationId"])
        self.assertTrue(second["reconciled"])
        self.assertEqual(1, len(self.sqs.messages))

    def test_pending_hitl_only_writes_deferred_receipt(self):
        out = self.call(event("PendingReview"))
        self.assertFalse(out["evaluationDispatched"])
        self.assertEqual([], self.sqs.messages)
        self.assertTrue(all(key.startswith("deferred/") for key in self.s3.objects))

    def test_in_progress_hitl_is_also_deferred(self):
        self.assertFalse(self.call(event("InProgress"))["evaluationDispatched"])

    def test_completed_review_is_hashed_immutable_and_dispatches(self):
        completed = self.call(event("Completed"))
        self.assertTrue(completed["evaluationDispatched"])
        request = json.loads(self.sqs.messages[0]["MessageBody"])
        self.assertEqual("AWS_IDP_HUMAN_REVIEW_COMPLETED", request["reviewAuthority"])
        self.assertIsNotNone(request["reviewEvidence"])
        self.assertTrue(any(key.startswith("reviews/") for key in self.s3.objects))
        review = json.loads(self.s3.objects[request["reviewEvidence"]["key"]]["Body"])
        self.assertEqual(request["snapshot"]["sha256"], review["documentSha256"])
        self.assertEqual(request["configVersion"], review["configVersion"])
        self.assertEqual(request["schemaProfile"], review["schemaProfile"])
        self.assertEqual(handler._sha(event("Completed")["executionArn"].encode("utf-8")), review["executionArnSha256"])
        self.assertIs(self.tracking.last_request["ConsistentRead"], True)

    def test_skipped_dispatches_but_is_never_human_reviewed(self):
        skipped = self.call(event("Skipped"))
        self.assertTrue(skipped["evaluationDispatched"])
        request = json.loads(self.sqs.messages[0]["MessageBody"])
        self.assertEqual("SKIPPED_NOT_REVIEWED", request["reviewAuthority"])
        self.assertIsNone(request["reviewEvidence"])

    def test_completed_review_missing_or_skip_all_is_rejected(self):
        self.tracking.items.clear()
        with self.assertRaises(handler.Rejected): self.call(event("Completed"))
        self.tracking.items[("doc#imports/invoice-7.pdf", "none")] = review_item() | {"HITLStatus": "Review Skipped"}
        with self.assertRaises(handler.Rejected): self.call(event("Completed"))

    def test_completed_review_requires_exact_section_coverage(self):
        for mutation in (
            {"HITLSectionsPending": ["1"]},
            {"HITLSectionsSkipped": ["1"]},
            {"HITLSectionsCompleted": []},
        ):
            self.tracking.items[("doc#imports/invoice-7.pdf", "none")] = review_item() | mutation
            with self.subTest(mutation=mutation):
                with self.assertRaises(handler.Rejected): self.call(event("Completed"))
        value = event("Completed"); value["document"]["hitl_sections_completed"] = []
        self.tracking.items[("doc#imports/invoice-7.pdf", "none")] = review_item()
        with self.assertRaises(handler.Rejected): self.call(value)

    def test_completed_review_rejects_ambiguous_identity_and_time(self):
        bad_records = [
            review_item()["HITLReviewHistory"] * 2,
            [review_item()["HITLReviewHistory"][0] | {"reviewedBy": "unknown"}],
            [review_item()["HITLReviewHistory"][0] | {"reviewedAt": "2026-08-28T12:30:00"}],
            [review_item()["HITLReviewHistory"][0] | {"sectionId": "other"}],
        ]
        for records in bad_records:
            self.tracking.items[("doc#imports/invoice-7.pdf", "none")] = review_item() | {"HITLReviewHistory": records}
            with self.subTest(records=records):
                with self.assertRaises(handler.Rejected): self.call(event("Completed"))

    def test_review_receipts_and_requests_contain_no_raw_identity(self):
        self.call(event("Completed"))
        payload = b"\n".join(value["Body"] for value in self.s3.objects.values()) + _messages(self.sqs)
        self.assertNotIn(b"reviewer-subject-7", payload)
        self.assertNotIn(b"Reviewer@Example.com", payload)
        self.assertNotIn(b"reviewer@example.com", payload)

    def test_completed_and_skipped_hitl_dispatch(self):
        self.assertTrue(self.call(event("Completed"))["evaluationDispatched"])
        self.s3, self.sqs = FakeS3({}), FakeSQS([])
        self.clients = handler.Clients(self.s3, self.sqs, self.tracking)
        self.assertTrue(self.call(event("Skipped"))["evaluationDispatched"])

    def test_wrong_hook_point_feature_args_and_execution_rejected(self):
        cases = []
        value = event(); value["hookPoint"] = "postExtraction"; cases.append(value)
        value = event(); value["featureId"] = "other"; cases.append(value)
        value = event(); value["argsMap"]["tenantId"] = "other"; cases.append(value)
        value = event(); value["executionArn"] = "arn:aws:states:us-west-2:123456789012:execution:EliteIdp-main:x"; cases.append(value)
        for candidate in cases:
            with self.subTest(candidate=candidate.get("hookPoint")):
                with self.assertRaises(handler.Rejected):
                    self.call(candidate)

    def test_config_version_and_class_rejected(self):
        value = event(); value["document"]["config_version"] = "latest"
        with self.assertRaises(handler.Rejected): self.call(value)
        value = event(); value["document"]["sections"][0]["classification"] = "Unknown"
        with self.assertRaises(handler.Rejected): self.call(value)

    def test_running_failed_evaluation_and_workflow_mismatch_rejected(self):
        for state in ("RUNNING", "FAILED"):
            value = event(); value["document"]["evaluation_status"] = state
            with self.assertRaises(handler.Rejected): self.call(value)
        value = event(); value["document"]["workflow_execution_arn"] = "arn:aws:states:us-east-1:123456789012:execution:EliteIdp-main:other"
        with self.assertRaises(handler.Rejected): self.call(value)

    def test_duplicate_sections_processing_errors_and_redacted_rejected(self):
        value = event(); value["document"]["sections"].append(dict(value["document"]["sections"][0]))
        with self.assertRaises(handler.Rejected): self.call(value)
        value = event(); value["document"]["processing_issues"] = [{"message": "bad"}]
        with self.assertRaises(handler.Rejected): self.call(value)
        value = event(); value["document"]["sections"][0]["processing_issues"] = [{"message": "bad"}]
        with self.assertRaises(handler.Rejected): self.call(value)
        value = event(); value["document"]["status"] = "REDACTED_SUPERSEDED"
        with self.assertRaises(handler.Rejected): self.call(value)

    def test_unknown_hitl_and_oversize_rejected(self):
        with self.assertRaises(handler.Rejected): self.call(event("WaitingForever"))
        small = cfg().__dict__ | {"max_canonical_bytes": 1024}
        value = event(); value["document"]["sections"][0]["attributes"]["blob"] = "x" * 2000
        with self.assertRaises(handler.Rejected):
            handler.lambda_handler(value, None, clients=self.clients, cfg=handler.Config(**small))

    def test_divergent_conditional_collision_rejected(self):
        value = event()
        canonical = handler._canonical(value["document"])
        sha = handler._sha(canonical)
        identity = handler._sha(b"tenant-a\nimports/invoice-7.pdf")
        key = f"snapshots/tenant-a/{identity}/{sha}.json"
        self.s3.objects[key] = {"Body": b"evil", "Metadata": {"sha256": "0" * 64}, "ContentLength": 4, "VersionId": "v-old"}
        with self.assertRaises(handler.Rejected): self.call(value)

    def test_sqs_failure_never_commits_dispatch_receipt(self):
        self.sqs.fail = True
        with self.assertRaises(RuntimeError): self.call(event())
        self.assertTrue(any(key.startswith("snapshots/") for key in self.s3.objects))
        self.assertFalse(any(key.startswith("dispatch/") for key in self.s3.objects))

    def test_official_loader_receives_working_bucket(self):
        with patch.object(handler, "load_hook_document", wraps=hooks.load_hook_document) as loader:
            self.call(event())
            loader.assert_called_once()
            self.assertEqual("working", loader.call_args.kwargs["working_bucket"])

    def test_no_response_can_authorize_business_persistence(self):
        for state in (None, "PendingReview", "InProgress", "Completed", "Skipped"):
            self.s3, self.sqs = FakeS3({}), FakeSQS([])
            self.clients = handler.Clients(self.s3, self.sqs, self.tracking)
            self.assertIs(False, self.call(event(state))["automaticBusinessPersistenceAuthorized"])


if __name__ == "__main__":
    unittest.main()
````

### FILE: `aws_idp_immutable_evaluation_handoff/hook-config.template.yaml`
```yaml
block_id: "AWS-IDP-IMMUTABLE-EVALUATION-HANDOFF:hook-config-template-yaml:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/blob/1b5fd74454e593de233342a02ee911af8ee38359/docs/feature-platform.md"
license: "LicenseRef-Workspace-Owner AND MIT-0"
sha256: "2239a38c0454171109bc6efe417c299d6c28e9c2399563165d5fab589519359a"
variables: []
secrets_allowed: false
```
````yaml
postprocessing:
  enabled: true
  featureId: elite-immutable-evaluation-handoff
  arn: REPLACE_WITH_EliteIdpEvaluationHandoffFunctionArn
  onError: fail
  args:
    - { key: tenantId, value: REPLACE_WITH_TENANT_ID }
    - { key: configVersion, value: REPLACE_WITH_PINNED_IDP_CONFIG_VERSION }
    - { key: schemaProfile, value: REPLACE_WITH_PINNED_SCHEMA_PROFILE }
  allowDocumentUpdate: false
````

### FILE: `aws_idp_immutable_evaluation_handoff/provider-profile.template.json`
```yaml
block_id: "AWS-IDP-IMMUTABLE-EVALUATION-HANDOFF:provider-profile-template-json:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-immutable-evaluation-handoff/provider-profile.template.json"
license: "LicenseRef-Workspace-Owner"
sha256: "946302965c2c247c3988e320c6638917ae89fd05a6429147e08baceb7d3f5f71"
variables: []
secrets_allowed: false
```
````json
{
  "acknowledgeConditionedState": false,
  "accountId": "REPLACE_WITH_12_DIGIT_ACCOUNT",
  "region": "REPLACE_WITH_REGION",
  "tenantId": "REPLACE_WITH_TENANT",
  "idpCommit": "1b5fd74454e593de233342a02ee911af8ee38359",
  "idpTag": "v0.6.5",
  "idpConfigVersion": "REPLACE_WITH_PINNED_CONFIG_VERSION",
  "schemaProfile": "REPLACE_WITH_PINNED_SCHEMA_PROFILE",
  "allowedClasses": [],
  "executionNamePrefix": "REPLACE_WITH_STATE_MACHINE_NAME_PREFIX",
  "workingBucket": "REPLACE_WITH_IDP_WORKING_BUCKET",
  "workingKmsKeyArn": "REPLACE_WITH_IDP_WORKING_KMS_KEY_ARN",
  "trackingTableName": "REPLACE_WITH_IDP_TRACKING_TABLE_NAME",
  "trackingTableArn": "REPLACE_WITH_IDP_TRACKING_TABLE_ARN",
  "idpCommonBaseLayerArn": "REPLACE_WITH_OFFICIAL_STACK_OUTPUT",
  "resultKmsKeyArn": "REPLACE_WITH_KMS_KEY_ARN",
  "logKmsKeyArn": "REPLACE_WITH_KMS_KEY_ARN",
  "retainDays": 0,
  "maxCanonicalBytes": 0,
  "owners": {"documents": "", "security": "", "privacy": "", "operations": "", "cost": "", "recovery": ""},
  "evidence": {"awsIdentity": "", "idpDeployment": "", "configEvaluation": "", "representativeCorpus": "", "hitlPolicy": "", "reviewAuthority": "", "kmsAndObjectLock": "", "dlqAndReplay": "", "restore": "", "costApproval": ""}
}
````

### FILE: `aws_idp_immutable_evaluation_handoff/template.yaml`
```yaml
block_id: "AWS-IDP-IMMUTABLE-EVALUATION-HANDOFF:template-yaml:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/tree/1b5fd74454e593de233342a02ee911af8ee38359; AWS SAM/S3/SQS/KMS contracts with local fail-closed hardening"
license: "LicenseRef-Workspace-Owner AND MIT-0"
sha256: "fdd7ec0e6576d2e403ae6f8ac91d08a1c39a3b159e21435fcbf4345c1438d7f7"
variables: []
secrets_allowed: false
```
````yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Description: Fail-closed immutable evaluation handoff for AWS GenAI IDP v0.6.5 postprocessing

Parameters:
  IdpCommonBaseLayerArn: {Type: String, AllowedPattern: '^arn:aws[a-zA-Z-]*:lambda:[a-z0-9-]+:[0-9]{12}:layer:[A-Za-z0-9-_]+:[0-9]+$'}
  WorkingBucketName: {Type: String, MinLength: 3}
  WorkingBucketArn: {Type: String, AllowedPattern: '^arn:aws[a-zA-Z-]*:s3:::[A-Za-z0-9.-]+$'}
  WorkingKmsKeyArn: {Type: String, AllowedPattern: '^arn:aws[a-zA-Z-]*:kms:[a-z0-9-]+:[0-9]{12}:key/[0-9a-f-]+$'}
  TrackingTableName: {Type: String, AllowedPattern: '^[A-Za-z0-9_.-]{3,255}$'}
  TrackingTableArn: {Type: String, AllowedPattern: '^arn:aws[a-zA-Z-]*:dynamodb:[a-z0-9-]+:[0-9]{12}:table/[A-Za-z0-9_.-]{3,255}$'}
  ResultKmsKeyArn: {Type: String, AllowedPattern: '^arn:aws[a-zA-Z-]*:kms:[a-z0-9-]+:[0-9]{12}:key/[0-9a-f-]+$'}
  LogKmsKeyArn: {Type: String, AllowedPattern: '^arn:aws[a-zA-Z-]*:kms:[a-z0-9-]+:[0-9]{12}:key/[0-9a-f-]+$'}
  TenantId: {Type: String, AllowedPattern: '^[a-z0-9][a-z0-9-]{1,62}$'}
  IdpConfigVersion: {Type: String, AllowedPattern: '^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$'}
  SchemaProfile: {Type: String, AllowedPattern: '^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$'}
  ExecutionNamePrefix: {Type: String, MinLength: 1, MaxLength: 128}
  AllowedClassesJson: {Type: String, Default: '[]'}
  RetainDays: {Type: Number, Default: 365, MinValue: 1, MaxValue: 36500}
  MaxCanonicalBytes: {Type: Number, Default: 20971520, MinValue: 1024, MaxValue: 52428800}

Resources:
  ResultBucket:
    Type: AWS::S3::Bucket
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      ObjectLockEnabled: true
      VersioningConfiguration: {Status: Enabled}
      BucketEncryption:
        ServerSideEncryptionConfiguration:
          - ServerSideEncryptionByDefault: {SSEAlgorithm: aws:kms, KMSMasterKeyID: !Ref ResultKmsKeyArn}
            BucketKeyEnabled: true
      PublicAccessBlockConfiguration: {BlockPublicAcls: true, BlockPublicPolicy: true, IgnorePublicAcls: true, RestrictPublicBuckets: true}
      OwnershipControls: {Rules: [{ObjectOwnership: BucketOwnerEnforced}]}
      LifecycleConfiguration:
        Rules:
          - Id: AbortIncompleteMultipartUploads
            Status: Enabled
            AbortIncompleteMultipartUpload: {DaysAfterInitiation: 1}
  ResultBucketPolicy:
    Type: AWS::S3::BucketPolicy
    Properties:
      Bucket: !Ref ResultBucket
      PolicyDocument:
        Version: '2012-10-17'
        Statement:
          - {Sid: DenyInsecureTransport, Effect: Deny, Principal: '*', Action: 's3:*', Resource: [!GetAtt ResultBucket.Arn, !Sub '${ResultBucket.Arn}/*'], Condition: {Bool: {'aws:SecureTransport': 'false'}}}
          - {Sid: DenyUnencryptedObjectWrites, Effect: Deny, Principal: '*', Action: 's3:PutObject', Resource: !Sub '${ResultBucket.Arn}/*', Condition: {StringNotEquals: {'s3:x-amz-server-side-encryption': 'aws:kms'}}}
          - {Sid: DenyWrongKmsKey, Effect: Deny, Principal: '*', Action: 's3:PutObject', Resource: !Sub '${ResultBucket.Arn}/*', Condition: {StringNotEquals: {'s3:x-amz-server-side-encryption-aws-kms-key-id': !Ref ResultKmsKeyArn}}}
  EvaluationDLQ:
    Type: AWS::SQS::Queue
    Properties: {FifoQueue: true, QueueName: !Sub '${AWS::StackName}-evaluation-dlq.fifo', KmsMasterKeyId: !Ref ResultKmsKeyArn, MessageRetentionPeriod: 1209600}
  EvaluationQueue:
    Type: AWS::SQS::Queue
    Properties:
      FifoQueue: true
      QueueName: !Sub '${AWS::StackName}-evaluation.fifo'
      KmsMasterKeyId: !Ref ResultKmsKeyArn
      VisibilityTimeout: 960
      MessageRetentionPeriod: 1209600
      RedrivePolicy: {deadLetterTargetArn: !GetAtt EvaluationDLQ.Arn, maxReceiveCount: 5}
  HandoffRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument: {Version: '2012-10-17', Statement: [{Effect: Allow, Principal: {Service: lambda.amazonaws.com}, Action: 'sts:AssumeRole'}]}
      Policies:
        - PolicyName: ExactHandoffEffects
          PolicyDocument:
            Version: '2012-10-17'
            Statement:
              - {Effect: Allow, Action: ['logs:CreateLogStream', 'logs:PutLogEvents'], Resource: !Sub '${HandoffLogGroup.Arn}:*'}
              - {Effect: Allow, Action: ['s3:GetObject', 's3:GetObjectVersion'], Resource: !Sub '${WorkingBucketArn}/compressed_documents/*'}
              - {Effect: Allow, Action: ['s3:PutObject', 's3:PutObjectRetention', 's3:GetObject', 's3:GetObjectVersion', 's3:GetObjectRetention'], Resource: !Sub '${ResultBucket.Arn}/*'}
              - {Effect: Allow, Action: ['sqs:SendMessage'], Resource: !GetAtt EvaluationQueue.Arn}
              - {Effect: Allow, Action: ['dynamodb:GetItem'], Resource: !Ref TrackingTableArn}
              - {Effect: Allow, Action: ['kms:Decrypt', 'kms:Encrypt', 'kms:GenerateDataKey', 'kms:DescribeKey'], Resource: !Ref ResultKmsKeyArn}
              - {Effect: Allow, Action: ['kms:Decrypt', 'kms:DescribeKey'], Resource: !Ref WorkingKmsKeyArn}
  HandoffFunction:
    Type: AWS::Serverless::Function
    Properties:
      CodeUri: .
      Handler: handler.lambda_handler
      FunctionName: !Sub '${AWS::StackName}-immutable-evaluation-handoff'
      Runtime: python3.13
      Architectures: [x86_64]
      MemorySize: 1024
      Timeout: 900
      ReservedConcurrentExecutions: 10
      Role: !GetAtt HandoffRole.Arn
      Layers: [!Ref IdpCommonBaseLayerArn]
      Environment:
        Variables:
          EXPECTED_ACCOUNT_ID: !Ref AWS::AccountId
          EXPECTED_REGION: !Ref AWS::Region
          EXPECTED_FEATURE_ID: elite-immutable-evaluation-handoff
          EXECUTION_NAME_PREFIX: !Ref ExecutionNamePrefix
          TENANT_ID: !Ref TenantId
          IDP_CONFIG_VERSION: !Ref IdpConfigVersion
          SCHEMA_PROFILE: !Ref SchemaProfile
          RESULT_BUCKET: !Ref ResultBucket
          RESULT_KMS_KEY_ARN: !Ref ResultKmsKeyArn
          EVALUATION_QUEUE_URL: !Ref EvaluationQueue
          WORKING_BUCKET: !Ref WorkingBucketName
          TRACKING_TABLE_NAME: !Ref TrackingTableName
          ALLOWED_CLASSES_JSON: !Ref AllowedClassesJson
          RETAIN_DAYS: !Ref RetainDays
          MAX_CANONICAL_BYTES: !Ref MaxCanonicalBytes
  HandoffLogGroup:
    Type: AWS::Logs::LogGroup
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties: {LogGroupName: !Sub '/aws/lambda/${AWS::StackName}-immutable-evaluation-handoff', RetentionInDays: 30, KmsKeyId: !Ref LogKmsKeyArn}

Outputs:
  EliteIdpEvaluationHandoffFunctionArn: {Value: !GetAtt HandoffFunction.Arn}
  EvaluationQueueArn: {Value: !GetAtt EvaluationQueue.Arn}
  EvaluationDLQArn: {Value: !GetAtt EvaluationDLQ.Arn}
  ResultBucketName: {Value: !Ref ResultBucket}
````

### FILE: `aws_idp_immutable_evaluation_handoff/deploy.ps1`
```yaml
block_id: "AWS-IDP-IMMUTABLE-EVALUATION-HANDOFF:deploy-ps1:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-immutable-evaluation-handoff/deploy.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "5c13dc48bcf10c74df1ab6c02c300eae5e90b1af33c4faca383726b9f6515dc4"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding(SupportsShouldProcess=$true, ConfirmImpact='High')]
param(
  [Parameter(Mandatory)][string]$ProfilePath,
  [Parameter(Mandatory)][string]$StackName,
  [switch]$Execute
)
$ErrorActionPreference = 'Stop'
$profile = Get-Content -LiteralPath $ProfilePath -Raw | ConvertFrom-Json -Depth 32
if ($profile.acknowledgeConditionedState -ne $true) { throw 'Conditioned state is not acknowledged.' }
foreach ($owner in @('documents','security','privacy','operations','cost','recovery')) { if ([string]::IsNullOrWhiteSpace($profile.owners.$owner)) { throw "Missing owner: $owner" } }
foreach ($evidence in @('awsIdentity','idpDeployment','configEvaluation','representativeCorpus','hitlPolicy','reviewAuthority','kmsAndObjectLock','dlqAndReplay','restore','costApproval')) { if ([string]::IsNullOrWhiteSpace($profile.evidence.$evidence)) { throw "Missing evidence: $evidence" } }
if ($profile.idpCommit -ne '1b5fd74454e593de233342a02ee911af8ee38359' -or $profile.idpTag -ne 'v0.6.5') { throw 'AWS IDP identity mismatch.' }
if (-not $profile.allowedClasses -or $profile.allowedClasses.Count -eq 0) { throw 'allowedClasses must be non-empty.' }
if (-not $Execute) { throw 'Dry gate complete. Re-run with -Execute and approve the change explicitly.' }
$parameters = @(
  "IdpCommonBaseLayerArn=$($profile.idpCommonBaseLayerArn)", "WorkingBucketName=$($profile.workingBucket)",
  "WorkingBucketArn=arn:aws:s3:::$($profile.workingBucket)", "WorkingKmsKeyArn=$($profile.workingKmsKeyArn)", "ResultKmsKeyArn=$($profile.resultKmsKeyArn)",
  "TrackingTableName=$($profile.trackingTableName)", "TrackingTableArn=$($profile.trackingTableArn)",
  "LogKmsKeyArn=$($profile.logKmsKeyArn)", "TenantId=$($profile.tenantId)",
  "IdpConfigVersion=$($profile.idpConfigVersion)", "SchemaProfile=$($profile.schemaProfile)",
  "ExecutionNamePrefix=$($profile.executionNamePrefix)", "AllowedClassesJson=$($profile.allowedClasses | ConvertTo-Json -Compress)",
  "RetainDays=$($profile.retainDays)", "MaxCanonicalBytes=$($profile.maxCanonicalBytes)"
)
if (-not $PSCmdlet.ShouldProcess($StackName, 'deploy immutable AWS IDP evaluation handoff')) { return }
& sam deploy --template-file (Join-Path $PSScriptRoot 'template.yaml') --stack-name $StackName --capabilities CAPABILITY_IAM --parameter-overrides @parameters
if ($LASTEXITCODE -ne 0) { throw "sam deploy failed: $LASTEXITCODE" }
````

### FILE: `aws_idp_immutable_evaluation_handoff/verify_contract.ps1`
```yaml
block_id: "AWS-IDP-IMMUTABLE-EVALUATION-HANDOFF:verify-contract-ps1:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-immutable-evaluation-handoff/verify_contract.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "0f9e48d16201cfa48b0732e821bda98a1833080fb4e41655c80c232f8b304ab9"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param([string]$ComponentRoot=$PSScriptRoot)
$ErrorActionPreference = 'Stop'
$required = @('README.md','handler.py','test_handler.py','template.yaml','hook-config.template.yaml','provider-profile.template.json','deploy.ps1','verify_contract.ps1')
foreach ($name in $required) { if (-not (Test-Path -LiteralPath (Join-Path $ComponentRoot $name) -PathType Leaf)) { throw "Missing $name" } }
$profile = Get-Content -LiteralPath (Join-Path $ComponentRoot 'provider-profile.template.json') -Raw | ConvertFrom-Json -Depth 32
if ($profile.acknowledgeConditionedState -ne $false -or $profile.allowedClasses.Count -ne 0 -or $profile.retainDays -ne 0 -or $profile.maxCanonicalBytes -ne 0) { throw 'Provider profile must begin fail-closed.' }
$hook = Get-Content -LiteralPath (Join-Path $ComponentRoot 'hook-config.template.yaml') -Raw
foreach ($needle in @('onError: fail','allowDocumentUpdate: false','featureId: elite-immutable-evaluation-handoff')) { if (-not $hook.Contains($needle)) { throw "Hook configuration missing $needle" } }
$handler = Get-Content -LiteralPath (Join-Path $ComponentRoot 'handler.py') -Raw
foreach ($needle in @('load_hook_document','ConsistentRead=True','HITLReviewHistory','SKIPPED_NOT_REVIEWED','reviewerSubjectSha256','IfNoneMatch="*"','ObjectLockMode="COMPLIANCE"','MessageDeduplicationId=evaluation_id','automaticBusinessPersistenceAuthorized')) { if (-not $handler.Contains($needle)) { throw "Handler contract missing $needle" } }
if ($handler -match 'print\(' -or $handler -match 'logger\.') { throw 'Handler must not log document identifiers or contents.' }
Push-Location $ComponentRoot
try {
  & python -B -c "from pathlib import Path; [compile(Path(p).read_text(encoding='utf-8'), p, 'exec') for p in ('handler.py','test_handler.py')]"
  if ($LASTEXITCODE -ne 0) { throw 'source compilation failed.' }
  & python -B -m unittest -v test_handler.py
  if ($LASTEXITCODE -ne 0) { throw 'unit tests failed.' }
} finally { Pop-Location }
Write-Output 'AWS_IDP_IMMUTABLE_EVALUATION_HANDOFF_PASS'
````

## 6. Configuration surface

`provider-profile.template.json` empieza cerrado: acknowledgement false, allow-list vacía, retención/tamaño cero, owners y evidencias vacíos. Los parámetros exigidos fijan cuenta, región, tenant, commit/tag IDP, config version, schema profile, state-machine prefix, working bucket/KMS, tracking table name/ARN, `IDPCommonBaseLayerArn`, KMS de resultados/logs, clases, retención y límites. `hook-config.template.yaml` exige `onError: fail` y `allowDocumentUpdate: false`.

## 7. Dependency bill

- AWS GenAI IDP Accelerator tag `v0.6.5`, commit `1b5fd74454e593de233342a02ee911af8ee38359`, licencia MIT-0.
- `idp_common` 0.6.5 y su base layer producida por el build oficial del mismo commit; Python `>=3.12,<3.14`.
- AWS Lambda Python 3.13, Boto3 provisto por el layer/runtime, AWS SAM, DynamoDB consistent read, S3 Versioning/Object Lock/KMS y SQS FIFO/DLQ.
- Ninguna dependencia propia se presenta como código verbatim AWS; 4 archivos son `ADAPTED` y 4 `AUTHORED`.

## 8. Apply order

1. Adquirir y verificar el commit/tag AWS mediante `OFFICIAL-UPSTREAM-ACQUISITION-CORE`.
2. Desplegar/evaluar AWS IDP y fijar sus outputs `IDPCommonBaseLayerArn`, working bucket/KMS y tracking table name/ARN.
3. Completar profile/owners/evidencia, ejecutar verifier y cfn-lint; aprobar costo/efectos.
4. Desplegar el stack del handoff; registrar su ARN en la config IDP exacta con el template provisto.
5. Conectar la FIFO únicamente a un evaluador que valide `evaluationId`, snapshot VersionId/SHA y schema profile.
6. Probar inline/comprimido, HITL doble, replay, race, SQS/DLQ, restore y rollback en DEV; sólo después considerar un gate empresarial separado.

## 9. Verification

`verify_contract.ps1` exige los ocho archivos, profile cerrado, hook fail-closed, consistent review tracking, `Skipped` no-review, identidades hashadas, primitives de immutability/dedup/no-storage, compilación y 20 pruebas. La admisión también exige cfn-lint, round-trip manifest/bloques, el helper/model real `idp_common` 0.6.5 desde el commit fijado y los verificadores globales. Live AWS, carreras durante review, grandes documentos comprimidos, KMS/Object Lock, DLQ/redrive, alarmas, costo, restore y downstream transactional siguen condicionados hasta evidencia target.

## 10. Reconstruction evidence

Reconstruido el 2026-08-28 desde AWS IDP tag `v0.6.5`/commit `1b5fd74454e593de233342a02ee911af8ee38359`. El componente pasó 20 pruebas, cfn-lint 1.55.1 y conserva el round-trip real de `Document`/`load_hook_document` con `idp_common` 0.6.5 en Python 3.12.14. La evidencia V93 registra review tracking, hashes, comandos, fallos y límites; V92 queda como snapshot histórico 0.1.0. No se afirma deploy productivo, firma AWS del receipt ni persistencia automática.
