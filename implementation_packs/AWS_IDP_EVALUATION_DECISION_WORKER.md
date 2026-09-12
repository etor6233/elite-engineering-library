# AWS IDP Evaluation Decision Worker

## 1. Metadata

```yaml
pack_id: "AWS-IDP-EVALUATION-DECISION-WORKER"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un consumidor SQS idempotente que verifica snapshot/policy/review AWS IDP por VersionId/SHA/KMS/Object Lock, valida cada sección con JSON Schema cerrado y congela una decisión durable sin autorizar persistencia empresarial."
stacks: ["AWS Powertools for Lambda Python 3.34.0", "AWS Lambda Python 3.13", "Amazon SQS FIFO", "Amazon DynamoDB", "Amazon S3 Object Lock", "AWS KMS", "AWS SAM"]
compatible_with: ["AWS-IDP-IMMUTABLE-EVALUATION-HANDOFF 0.1.x", "AWS-POWERTOOLS-IDEMPOTENT-SQS-BATCH-COMPONENT 0.1.x", "STRICT-DOCUMENT-FIELD-EVALUATION-GATE 0.1.x", "OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x"]
incompatible_with: ["evaluationId no gobernado por productor", "policy móvil o no aprobada", "review sin binding autocontenido", "bucket sin Versioning/Object Lock/KMS", "schema abierto/remoto/combinatorio", "persistencia empresarial automática"]
license_expression: "LicenseRef-Workspace-Owner AND MIT-0"
upstream_sources: ["https://github.com/aws-powertools/powertools-lambda-python/tree/376757161b002f0c2f5d19d5cdf9def5b8c45704", "https://github.com/aws-powertools/powertools-lambda-python/blob/376757161b002f0c2f5d19d5cdf9def5b8c45704/examples/idempotency/src/integrate_idempotency_with_batch_processor.py", "https://github.com/aws-powertools/powertools-lambda-python/blob/376757161b002f0c2f5d19d5cdf9def5b8c45704/examples/idempotency/templates/sam.yaml", "https://github.com/aws-powertools/powertools-lambda-python/blob/376757161b002f0c2f5d19d5cdf9def5b8c45704/examples/batch_processing/sam/sqs_batch_processing.yaml", "https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html", "https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/transactional-outbox.html"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use después del handoff inmutable AWS IDP 0.1.2 y antes de cualquier persistencia. Requiere Powertools 3.34.0 oficial, policy aprobada e inmutable, schemas cerrados por clase, revisión humana completa, snapshots/review receipts Object-Locked, cuenta/IAM/KMS/SQS/Dynamo target y owners/evidencia. No sustituye evaluación offline contra ground truth ni transacción negocio+outbox.

## 3. Architecture contract

La composición `ADAPTED` conserva `BatchProcessor`, `process_partial_response`, `idempotent_function`, timeout context y validator oficiales. La clave obligatoria es `evaluationId`; snapshot/review forman payload validation y SHA-256. Cada referencia se limita al prefijo exacto del tenant y se lee por VersionId, checksum, metadata SHA, SSE-KMS y Object Lock COMPLIANCE. Policy exige aprobación, review para toda clase y schemas recursivamente cerrados sin refs/combinadores. `Completed` sólo pasa con cobertura exacta; `Skipped`/sin review y schema inválido producen decisión negativa durable. El receipt v2 determinista se escribe create-only/Object-Locked, conserva la referencia exacta del snapshot sin copiar sus valores y siempre fija business persistence false.

## 4. Exact file manifest

```text
CREATE aws_idp_evaluation_decision_worker/README.md
CREATE aws_idp_evaluation_decision_worker/handler.py
CREATE aws_idp_evaluation_decision_worker/test_handler.py
CREATE aws_idp_evaluation_decision_worker/evaluation-policy.template.json
CREATE aws_idp_evaluation_decision_worker/provider-profile.template.json
CREATE aws_idp_evaluation_decision_worker/template.yaml
CREATE aws_idp_evaluation_decision_worker/deploy.ps1
CREATE aws_idp_evaluation_decision_worker/verify_contract.ps1
```

## 5. Materialization blocks

### FILE: `aws_idp_evaluation_decision_worker/README.md`
```yaml
block_id: "AWS-IDP-EVALUATION-DECISION-WORKER:readme-md:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-evaluation-decision-worker/README.md"
license: "LicenseRef-Workspace-Owner"
sha256: "5477de232515242faef9968ec6a142fae55bcfaa9c8aab9d6a0593c2498a726d"
variables: []
secrets_allowed: false
```
````markdown
# AWS IDP evaluation decision worker

This `ADAPTED` component consumes the immutable handoff FIFO with the official AWS Powertools Batch Processor, DynamoDB idempotency and JSON Schema validator. AWS publishes those primitives and their composition examples; AWS does not publish or claim this project-specific decision handler.

The worker requires `evaluationId` as producer-owned idempotency key, validates payload tampering, reads snapshot/policy/review by exact S3 VersionId and SHA-256, requires SSE-KMS plus Object Lock COMPLIANCE, validates every section against an approved closed schema, and accepts human review only when the self-contained review evidence binds the same document/config/schema and covers every section. A skipped or absent review creates a durable negative decision instead of retrying forever.

Every v2 decision is deterministic, content-addressed, immutable, carries the exact immutable snapshot reference without extracted values, and fixes `automaticBusinessPersistenceAuthorized=false`. The Powertools record may expire after the configured long window; the decision receipt remains the durable replay authority. This worker does not write a business table. A separately approved transaction must consume only a positive decision and atomically write business data plus outbox, with idempotency and reconciliation.

The policy template intentionally starts unusable. Upload an approved policy to the existing Object Lock result bucket, then pin its key, VersionId and SHA-256. Live AWS, account/IAM, Powertools layer identity, KMS, Object Lock, SQS/DLQ, load, restore, policy/corpus and transactional persistence evidence remain required.
````

### FILE: `aws_idp_evaluation_decision_worker/handler.py`
```yaml
block_id: "AWS-IDP-EVALUATION-DECISION-WORKER:handler-py:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/aws-powertools/powertools-lambda-python/tree/376757161b002f0c2f5d19d5cdf9def5b8c45704"
license: "MIT-0 AND LicenseRef-Workspace-Owner"
sha256: "21b3256df1d977a84755c68498c7ad1740c219b86a67848183305dfced5cb1f5"
variables: []
secrets_allowed: false
```
````python
"""ADAPTED AWS IDP evaluation consumer using official AWS Powertools primitives."""
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
from aws_lambda_powertools.utilities.batch import BatchProcessor, EventType, process_partial_response
from aws_lambda_powertools.utilities.data_classes.sqs_event import SQSRecord
from aws_lambda_powertools.utilities.idempotency import DynamoDBPersistenceLayer, IdempotencyConfig, idempotent_function
from aws_lambda_powertools.utilities.typing import LambdaContext
from aws_lambda_powertools.utilities.validation import validate
from aws_lambda_powertools.utilities.validation.exceptions import SchemaValidationError


REQUEST_SCHEMA = "elite-aws-idp-evaluation-handoff/v1"
POLICY_SCHEMA = "elite-aws-idp-runtime-evaluation-policy/v1"
REVIEW_SCHEMA = "elite-aws-idp-review-evidence/v1"
DECISION_SCHEMA = "elite-aws-idp-evaluation-decision/v2"
HEX = re.compile(r"^[0-9a-f]{64}$")
TOKEN = re.compile(r"^[A-Za-z0-9][A-Za-z0-9._:/+=,@-]{0,511}$")
FORBIDDEN_SCHEMA = {"$ref", "$dynamicRef", "allOf", "anyOf", "oneOf", "not", "if", "then", "else", "patternProperties", "unevaluatedProperties", "unevaluatedItems"}


class Rejected(ValueError):
    pass


@dataclass(frozen=True)
class Config:
    tenant_id: str
    config_version: str
    schema_profile: str
    result_bucket: str
    result_kms_key_arn: str
    policy_key: str
    policy_version_id: str
    policy_sha256: str
    retain_days: int
    max_object_bytes: int

    @staticmethod
    def from_env() -> "Config":
        value = Config(
            tenant_id=os.environ.get("TENANT_ID", ""), config_version=os.environ.get("IDP_CONFIG_VERSION", ""),
            schema_profile=os.environ.get("SCHEMA_PROFILE", ""), result_bucket=os.environ.get("RESULT_BUCKET", ""),
            result_kms_key_arn=os.environ.get("RESULT_KMS_KEY_ARN", ""), policy_key=os.environ.get("POLICY_KEY", ""),
            policy_version_id=os.environ.get("POLICY_VERSION_ID", ""), policy_sha256=os.environ.get("POLICY_SHA256", ""),
            retain_days=int(os.environ.get("RETAIN_DAYS", "365")), max_object_bytes=int(os.environ.get("MAX_OBJECT_BYTES", "20971520")),
        )
        value.validate()
        return value

    def validate(self) -> None:
        if any(not TOKEN.fullmatch(value) for value in (self.tenant_id, self.config_version, self.schema_profile)):
            raise Rejected("identity configuration is invalid")
        if not self.result_bucket or not self.result_kms_key_arn or not self.policy_version_id:
            raise Rejected("storage configuration is missing")
        if not self.policy_key.startswith(f"policies/{self.tenant_id}/") or ".." in self.policy_key.split("/"):
            raise Rejected("policy key is outside the tenant prefix")
        if not HEX.fullmatch(self.policy_sha256):
            raise Rejected("policy hash is invalid")
        if not (1 <= self.retain_days <= 36500 and 1024 <= self.max_object_bytes <= 50 * 1024 * 1024):
            raise Rejected("limits are invalid")


@dataclass(frozen=True)
class Clients:
    s3: Any


def _clients() -> Clients:
    return Clients(s3=boto3.client("s3"))


def _sha(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def _canonical(value: Any) -> bytes:
    try:
        return json.dumps(value, sort_keys=True, separators=(",", ":"), ensure_ascii=False, allow_nan=False).encode("utf-8")
    except (TypeError, ValueError) as error:
        raise Rejected("value is not canonical JSON") from error


def _pairs(pairs: list[tuple[str, Any]]) -> dict[str, Any]:
    value: dict[str, Any] = {}
    for key, item in pairs:
        if key in value:
            raise Rejected("duplicate JSON key")
        value[key] = item
    return value


def _json(raw: bytes, label: str) -> Any:
    try:
        return json.loads(raw.decode("utf-8"), object_pairs_hook=_pairs)
    except (UnicodeDecodeError, json.JSONDecodeError) as error:
        raise Rejected(f"{label} is not strict UTF-8 JSON") from error


def _checksum(data: bytes) -> str:
    return base64.b64encode(hashlib.sha256(data).digest()).decode("ascii")


def _precondition(error: ClientError) -> bool:
    return str(error.response.get("Error", {}).get("Code", "")) in {"PreconditionFailed", "412", "ConditionalRequestConflict", "409"}


def _read_ref(clients: Clients, cfg: Config, ref: dict[str, Any], label: str) -> tuple[dict[str, Any], bytes]:
    if set(ref) != {"bucket", "key", "versionId", "sha256"} or ref.get("bucket") != cfg.result_bucket:
        raise Rejected(f"{label} reference is invalid")
    key, version, expected = ref.get("key"), ref.get("versionId"), ref.get("sha256")
    if not isinstance(key, str) or not key or ".." in key.split("/") or not isinstance(version, str) or not version or not isinstance(expected, str) or not HEX.fullmatch(expected):
        raise Rejected(f"{label} reference is incomplete")
    prefixes = {"policy": f"policies/{cfg.tenant_id}/", "snapshot": f"snapshots/{cfg.tenant_id}/", "review": f"reviews/{cfg.tenant_id}/"}
    if label not in prefixes or not key.startswith(prefixes[label]):
        raise Rejected(f"{label} reference is outside the tenant prefix")
    response = clients.s3.get_object(Bucket=cfg.result_bucket, Key=key, VersionId=version, ChecksumMode="ENABLED")
    if response.get("VersionId") != version or response.get("ServerSideEncryption") != "aws:kms" or response.get("SSEKMSKeyId") != cfg.result_kms_key_arn:
        raise Rejected(f"{label} storage authority mismatch")
    if response.get("ObjectLockMode") != "COMPLIANCE":
        raise Rejected(f"{label} is not compliance locked")
    body = response["Body"].read(cfg.max_object_bytes + 1)
    if len(body) > cfg.max_object_bytes or _sha(body) != expected or response.get("ChecksumSHA256") != _checksum(body):
        raise Rejected(f"{label} content integrity mismatch")
    if (response.get("Metadata") or {}).get("sha256") != expected:
        raise Rejected(f"{label} metadata integrity mismatch")
    value = _json(body, label)
    if not isinstance(value, dict):
        raise Rejected(f"{label} must be an object")
    return value, body


def _policy_ref(cfg: Config) -> dict[str, Any]:
    return {"bucket": cfg.result_bucket, "key": cfg.policy_key, "versionId": cfg.policy_version_id, "sha256": cfg.policy_sha256}


def _closed_schema(node: Any) -> None:
    if not isinstance(node, dict) or FORBIDDEN_SCHEMA.intersection(node):
        raise Rejected("policy contains an unsupported JSON Schema")
    kind = node.get("type")
    if kind == "object":
        props, required = node.get("properties"), node.get("required")
        if not isinstance(props, dict) or not props or node.get("additionalProperties") is not False or not isinstance(required, list) or set(required) != set(props):
            raise Rejected("every object schema must be closed and fully required")
        for child in props.values():
            _closed_schema(child)
    elif kind == "array":
        _closed_schema(node.get("items"))
    elif kind not in {"string", "integer", "number", "boolean"}:
        raise Rejected("schema type is not admitted")


def _validate_policy(policy: dict[str, Any], cfg: Config) -> dict[str, Any]:
    exact = {"schema", "approval", "configVersion", "schemaProfile", "requireHumanReviewForAll", "classSchemas", "automaticBusinessPersistenceAuthorized"}
    if set(policy) != exact or policy.get("schema") != POLICY_SCHEMA or policy.get("configVersion") != cfg.config_version or policy.get("schemaProfile") != cfg.schema_profile:
        raise Rejected("policy identity mismatch")
    approval = policy.get("approval")
    if not isinstance(approval, dict) or set(approval) != {"status", "approvalId", "approvedBy", "approvedAt"} or approval.get("status") != "APPROVED":
        raise Rejected("policy is not approved")
    if any(not isinstance(approval.get(key), str) or not approval[key].strip() or approval[key] == "REQUIRED" for key in ("approvalId", "approvedBy", "approvedAt")):
        raise Rejected("policy approval evidence is unresolved")
    try:
        approved = datetime.fromisoformat(approval["approvedAt"].replace("Z", "+00:00"))
    except ValueError as error:
        raise Rejected("policy approval time is invalid") from error
    if approved.tzinfo is None or policy.get("requireHumanReviewForAll") is not True or policy.get("automaticBusinessPersistenceAuthorized") is not False:
        raise Rejected("policy is not fail closed")
    schemas = policy.get("classSchemas")
    if not isinstance(schemas, dict) or not schemas or any(not TOKEN.fullmatch(key) for key in schemas):
        raise Rejected("class schemas are missing")
    for schema in schemas.values():
        _closed_schema(schema)
    return schemas


def _validate_request(request: dict[str, Any], cfg: Config) -> None:
    exact = {"schema", "evaluationId", "tenantId", "configVersion", "schemaProfile", "snapshot", "reviewAuthority", "reviewEvidence", "automaticBusinessPersistenceAuthorized"}
    if set(request) != exact or request.get("schema") != REQUEST_SCHEMA or request.get("tenantId") != cfg.tenant_id or request.get("configVersion") != cfg.config_version or request.get("schemaProfile") != cfg.schema_profile:
        raise Rejected("evaluation request identity mismatch")
    if request.get("automaticBusinessPersistenceAuthorized") is not False:
        raise Rejected("upstream request attempted to authorize persistence")
    snapshot = request.get("snapshot")
    if not isinstance(snapshot, dict) or not isinstance(request.get("evaluationId"), str) or not HEX.fullmatch(request["evaluationId"]):
        raise Rejected("evaluation request is incomplete")
    review = request.get("reviewEvidence")
    authority = request.get("reviewAuthority")
    suffix = review.get("sha256") if isinstance(review, dict) else authority
    expected = _sha((cfg.tenant_id + "\n" + cfg.config_version + "\n" + str(snapshot.get("sha256")) + "\n" + str(suffix)).encode("utf-8"))
    if expected != request["evaluationId"]:
        raise Rejected("evaluation id does not bind the request")


def _decision_key(cfg: Config, evaluation_id: str) -> str:
    return f"decisions/{cfg.tenant_id}/{_sha(evaluation_id.encode('utf-8'))}.json"


def _put_decision(clients: Clients, cfg: Config, key: str, body: bytes) -> str:
    digest = _sha(body)
    try:
        response = clients.s3.put_object(
            Bucket=cfg.result_bucket, Key=key, Body=body, ContentType="application/json", IfNoneMatch="*",
            ChecksumAlgorithm="SHA256", ChecksumSHA256=_checksum(body), ServerSideEncryption="aws:kms", SSEKMSKeyId=cfg.result_kms_key_arn,
            ObjectLockMode="COMPLIANCE", ObjectLockRetainUntilDate=datetime.now(timezone.utc) + timedelta(days=cfg.retain_days),
            Metadata={"sha256": digest, "schema": DECISION_SCHEMA},
        )
        if not response.get("VersionId"):
            raise Rejected("decision write returned no version")
        return response["VersionId"]
    except ClientError as error:
        if not _precondition(error):
            raise
        head = clients.s3.head_object(Bucket=cfg.result_bucket, Key=key, ChecksumMode="ENABLED")
        if (head.get("Metadata") or {}).get("sha256") != digest or head.get("ContentLength") != len(body) or head.get("ObjectLockMode") != "COMPLIANCE":
            raise Rejected("existing decision is divergent") from error
        if not head.get("VersionId"):
            raise Rejected("existing decision has no version") from error
        return head["VersionId"]


def evaluate_request_core(request: dict[str, Any], *, clients: Clients, cfg: Config) -> dict[str, Any]:
    _validate_request(request, cfg)
    policy, _ = _read_ref(clients, cfg, _policy_ref(cfg), "policy")
    schemas = _validate_policy(policy, cfg)
    snapshot, _ = _read_ref(clients, cfg, request["snapshot"], "snapshot")
    if snapshot.get("config_version") != cfg.config_version or snapshot.get("errors") or snapshot.get("processing_issues"):
        raise Rejected("snapshot is not admissible")
    sections = snapshot.get("sections")
    if not isinstance(sections, list) or not sections:
        raise Rejected("snapshot sections are missing")
    section_ids: set[str] = set()
    classes: list[str] = []
    schema_valid = True
    for section in sections:
        if not isinstance(section, dict) or not isinstance(section.get("section_id"), str) or section["section_id"] in section_ids:
            raise Rejected("snapshot section identity is invalid")
        classification, attributes = section.get("classification"), section.get("attributes")
        if classification not in schemas or not isinstance(attributes, dict):
            raise Rejected("snapshot class or attributes are not admitted")
        try:
            validate(event=attributes, schema=schemas[classification])
        except SchemaValidationError:
            schema_valid = False
        section_ids.add(section["section_id"])
        classes.append(classification)

    review_sha = None
    passed = False
    reason = "HUMAN_REVIEW_REQUIRED"
    if request["reviewAuthority"] == "AWS_IDP_HUMAN_REVIEW_COMPLETED":
        if not isinstance(request["reviewEvidence"], dict):
            raise Rejected("completed review has no evidence")
        review, _ = _read_ref(clients, cfg, request["reviewEvidence"], "review")
        review_sha = request["reviewEvidence"]["sha256"]
        if review.get("schema") != REVIEW_SCHEMA or review.get("tenantId") != cfg.tenant_id or review.get("reviewAuthority") != request["reviewAuthority"] or review.get("automaticBusinessPersistenceAuthorized") is not False:
            raise Rejected("review authority is invalid")
        if review.get("documentSha256") != request["snapshot"]["sha256"] or review.get("configVersion") != cfg.config_version or review.get("schemaProfile") != cfg.schema_profile:
            raise Rejected("review evidence is not bound to the snapshot")
        reviews = review.get("sections")
        if review.get("sectionCount") != len(section_ids) or not isinstance(reviews, list) or len(reviews) != len(section_ids) or not all(isinstance(item, dict) for item in reviews) or {item.get("sectionId") for item in reviews} != section_ids:
            raise Rejected("review evidence does not cover the snapshot")
        for item in reviews:
            if not isinstance(item.get("reviewerSubjectSha256"), str) or not HEX.fullmatch(item["reviewerSubjectSha256"]) or not isinstance(item.get("reviewedAt"), str):
                raise Rejected("review evidence entry is incomplete")
            try:
                reviewed_at = datetime.fromisoformat(item["reviewedAt"].replace("Z", "+00:00"))
            except ValueError as error:
                raise Rejected("review evidence time is invalid") from error
            if reviewed_at.tzinfo is None:
                raise Rejected("review evidence time has no timezone")
        passed = schema_valid
        reason = "SCHEMA_AND_HUMAN_REVIEW_VERIFIED" if schema_valid else "SCHEMA_VALIDATION_FAILED"
    elif request["reviewAuthority"] not in {"NOT_TRIGGERED", "SKIPPED_NOT_REVIEWED"} or request["reviewEvidence"] is not None:
        raise Rejected("review state is not admitted")
    elif not schema_valid:
        reason = "SCHEMA_VALIDATION_FAILED"

    decision = {
        "schema": DECISION_SCHEMA, "evaluationId": request["evaluationId"], "tenantId": cfg.tenant_id,
        "configVersion": cfg.config_version, "schemaProfile": cfg.schema_profile,
        "snapshotReference": request["snapshot"], "snapshotSha256": request["snapshot"]["sha256"], "policySha256": cfg.policy_sha256,
        "reviewEvidenceSha256": review_sha, "reviewAuthority": request["reviewAuthority"],
        "sectionCount": len(section_ids), "classifications": sorted(classes),
        "evaluationPassed": passed, "decisionReason": reason,
        "automaticBusinessPersistenceAuthorized": False,
    }
    body = _canonical(decision)
    version = _put_decision(clients, cfg, _decision_key(cfg, request["evaluationId"]), body)
    return {"evaluationId": request["evaluationId"], "decisionVersionId": version, "evaluationPassed": passed, "automaticBusinessPersistenceAuthorized": False}


processor = BatchProcessor(event_type=EventType.SQS)
persistence = DynamoDBPersistenceLayer(table_name=os.environ.get("IDEMPOTENCY_TABLE", ""))
idempotency = IdempotencyConfig(
    event_key_jmespath="evaluationId", payload_validation_jmespath="[snapshot.sha256, reviewAuthority, reviewEvidence.sha256]",
    raise_on_no_idempotency_key=True, expires_after_seconds=int(os.environ.get("IDEMPOTENCY_SECONDS", "315360000")), hash_function="sha256",
)


@idempotent_function(data_keyword_argument="request", config=idempotency, persistence_store=persistence)
def evaluate_request(request: dict[str, Any]) -> dict[str, Any]:
    return evaluate_request_core(request, clients=_clients(), cfg=Config.from_env())


def record_handler(record: SQSRecord) -> dict[str, Any]:
    request = _json(record.body.encode("utf-8"), "SQS message")
    if not isinstance(request, dict):
        raise Rejected("SQS message must be an object")
    return evaluate_request(request=request)


def lambda_handler(event: dict[str, Any], context: LambdaContext) -> dict[str, Any]:
    idempotency.register_lambda_context(context)
    return process_partial_response(event=event, context=context, processor=processor, record_handler=record_handler)
````

### FILE: `aws_idp_evaluation_decision_worker/test_handler.py`
```yaml
block_id: "AWS-IDP-EVALUATION-DECISION-WORKER:test-handler-py:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-evaluation-decision-worker/test_handler.py"
license: "LicenseRef-Workspace-Owner"
sha256: "bf76bec5d98ad2faa903ed81f8f61c793435b1b02ff781029bbc619b326d7b2a"
variables: []
secrets_allowed: false
```
````python
import base64
import hashlib
import json
import sys
import types
import unittest
from dataclasses import dataclass
from unittest.mock import patch


class FakeClientError(Exception):
    def __init__(self, code): self.response={"Error":{"Code":code}}


botocore=types.ModuleType("botocore"); exceptions=types.ModuleType("botocore.exceptions"); exceptions.ClientError=FakeClientError; botocore.exceptions=exceptions
sys.modules.setdefault("botocore",botocore);sys.modules.setdefault("botocore.exceptions",exceptions)
boto3=types.ModuleType("boto3");boto3.client=lambda _name:None;sys.modules.setdefault("boto3",boto3)

powertools=types.ModuleType("aws_lambda_powertools");utilities=types.ModuleType("aws_lambda_powertools.utilities")
batch=types.ModuleType("aws_lambda_powertools.utilities.batch")
class BatchProcessor:
    def __init__(self,**kwargs): self.kwargs=kwargs
class EventType: SQS="SQS"
def process_partial_response(**kwargs): return {"called":True,"event":kwargs["event"]}
batch.BatchProcessor=BatchProcessor;batch.EventType=EventType;batch.process_partial_response=process_partial_response
data_classes=types.ModuleType("aws_lambda_powertools.utilities.data_classes");sqs_event=types.ModuleType("aws_lambda_powertools.utilities.data_classes.sqs_event")
class SQSRecord: pass
sqs_event.SQSRecord=SQSRecord
idem=types.ModuleType("aws_lambda_powertools.utilities.idempotency")
class DynamoDBPersistenceLayer:
    def __init__(self,**kwargs): self.kwargs=kwargs
class IdempotencyConfig:
    def __init__(self,**kwargs): self.kwargs=kwargs;self.context=None
    def register_lambda_context(self,context): self.context=context
def idempotent_function(**_kwargs):
    def decorator(fn): return fn
    return decorator
idem.DynamoDBPersistenceLayer=DynamoDBPersistenceLayer;idem.IdempotencyConfig=IdempotencyConfig;idem.idempotent_function=idempotent_function
typing=types.ModuleType("aws_lambda_powertools.utilities.typing");typing.LambdaContext=object
validation=types.ModuleType("aws_lambda_powertools.utilities.validation");validation.validate=lambda **kwargs: None
validation_exceptions=types.ModuleType("aws_lambda_powertools.utilities.validation.exceptions")
class SchemaValidationError(Exception):pass
validation_exceptions.SchemaValidationError=SchemaValidationError
for name,module in [("aws_lambda_powertools",powertools),("aws_lambda_powertools.utilities",utilities),("aws_lambda_powertools.utilities.batch",batch),("aws_lambda_powertools.utilities.data_classes",data_classes),("aws_lambda_powertools.utilities.data_classes.sqs_event",sqs_event),("aws_lambda_powertools.utilities.idempotency",idem),("aws_lambda_powertools.utilities.typing",typing),("aws_lambda_powertools.utilities.validation",validation),("aws_lambda_powertools.utilities.validation.exceptions",validation_exceptions)]:sys.modules.setdefault(name,module)

import handler


class Body:
    def __init__(self,value):self.value=value
    def read(self,_limit):return self.value


@dataclass
class FakeS3:
    objects:dict
    def add(self,key,value,version="v1",kms="kms"):
        raw=handler._canonical(value);digest=handler._sha(raw)
        self.objects[(key,version)]={"Body":raw,"VersionId":version,"ServerSideEncryption":"aws:kms","SSEKMSKeyId":kms,"ObjectLockMode":"COMPLIANCE","ChecksumSHA256":base64.b64encode(hashlib.sha256(raw).digest()).decode(),"Metadata":{"sha256":digest},"ContentLength":len(raw)}
        return {"bucket":"result","key":key,"versionId":version,"sha256":digest}
    def get_object(self,**kwargs):
        value=self.objects[(kwargs["Key"],kwargs["VersionId"])].copy();value["Body"]=Body(value["Body"]);return value
    def put_object(self,**kwargs):
        key=(kwargs["Key"],"decision-v1")
        if any(k[0]==kwargs["Key"] for k in self.objects):raise FakeClientError("PreconditionFailed")
        self.objects[key]={"Body":kwargs["Body"],"VersionId":"decision-v1","ObjectLockMode":kwargs["ObjectLockMode"],"Metadata":kwargs["Metadata"],"ContentLength":len(kwargs["Body"])}
        return {"VersionId":"decision-v1"}
    def head_object(self,**kwargs):
        matches=[v for (key,_version),v in self.objects.items() if key==kwargs["Key"]]
        if not matches:raise FakeClientError("404")
        return matches[0]


def config():
    return handler.Config("tenant-a","cfg-7","invoice-v3","result","kms","policies/tenant-a/policy.json","p1","0"*64,365,1024*1024)


def policy():
    return {"schema":handler.POLICY_SCHEMA,"approval":{"status":"APPROVED","approvalId":"approval-7","approvedBy":"owner-7","approvedAt":"2026-08-28T10:00:00Z"},"configVersion":"cfg-7","schemaProfile":"invoice-v3","requireHumanReviewForAll":True,"classSchemas":{"Invoice":{"type":"object","properties":{"total":{"type":"string"}},"required":["total"],"additionalProperties":False}},"automaticBusinessPersistenceAuthorized":False}


def snapshot():
    return {"id":"secret-file.pdf","config_version":"cfg-7","sections":[{"section_id":"1","classification":"Invoice","attributes":{"total":"10.00"}}]}


def review(snapshot_sha):
    return {"schema":handler.REVIEW_SCHEMA,"tenantId":"tenant-a","documentIdSha256":"1"*64,"documentSha256":snapshot_sha,"configVersion":"cfg-7","schemaProfile":"invoice-v3","executionArnSha256":"2"*64,"trackingStateSha256":"3"*64,"sections":[{"sectionId":"1","reviewedAt":"2026-08-28T10:00:00+00:00","reviewerSubjectSha256":"4"*64}],"sectionCount":1,"reviewAuthority":"AWS_IDP_HUMAN_REVIEW_COMPLETED","automaticBusinessPersistenceAuthorized":False}


class WorkerTests(unittest.TestCase):
    def setUp(self):
        self.s3=FakeS3({});self.cfg=config()
        pref=self.s3.add(self.cfg.policy_key,policy(),self.cfg.policy_version_id)
        self.cfg=handler.Config(**(self.cfg.__dict__|{"policy_sha256":pref["sha256"]}))
        self.snapshot_ref=self.s3.add("snapshots/tenant-a/x.json",snapshot(),"s1")
        self.review_ref=self.s3.add("reviews/tenant-a/x.json",review(self.snapshot_ref["sha256"]),"r1")

    def request(self,authority="AWS_IDP_HUMAN_REVIEW_COMPLETED",review_ref=True):
        ref=self.review_ref if review_ref else None;suffix=ref["sha256"] if ref else authority
        eid=handler._sha(("tenant-a\ncfg-7\n"+self.snapshot_ref["sha256"]+"\n"+suffix).encode())
        return {"schema":handler.REQUEST_SCHEMA,"evaluationId":eid,"tenantId":"tenant-a","configVersion":"cfg-7","schemaProfile":"invoice-v3","snapshot":self.snapshot_ref,"reviewAuthority":authority,"reviewEvidence":ref,"automaticBusinessPersistenceAuthorized":False}

    def call(self,value):return handler.evaluate_request_core(value,clients=handler.Clients(self.s3),cfg=self.cfg)

    def test_completed_review_writes_positive_but_never_storage_authority(self):
        out=self.call(self.request());self.assertTrue(out["evaluationPassed"]);self.assertFalse(out["automaticBusinessPersistenceAuthorized"])
        decision=[v for (k,_),v in self.s3.objects.items() if k.startswith("decisions/")][0]["Body"]
        parsed=json.loads(decision);self.assertTrue(parsed["evaluationPassed"]);self.assertFalse(parsed["automaticBusinessPersistenceAuthorized"]);self.assertEqual(self.snapshot_ref,parsed["snapshotReference"])

    def test_skipped_and_not_triggered_are_durable_negative_decisions(self):
        for authority in ("SKIPPED_NOT_REVIEWED","NOT_TRIGGERED"):
            self.s3.objects={k:v for k,v in self.s3.objects.items() if not k[0].startswith("decisions/")}
            out=self.call(self.request(authority,False));self.assertFalse(out["evaluationPassed"])

    def test_evaluation_id_tamper_is_rejected(self):
        value=self.request();value["evaluationId"]="0"*64
        with self.assertRaises(handler.Rejected):self.call(value)

    def test_snapshot_hash_and_storage_authority_are_enforced(self):
        value=self.request();value["snapshot"]=dict(value["snapshot"]);value["snapshot"]["sha256"]="0"*64
        with self.assertRaises(handler.Rejected):self.call(value)
        self.s3.objects[("snapshots/tenant-a/x.json","s1")]["ObjectLockMode"]="GOVERNANCE"
        with self.assertRaises(handler.Rejected):self.call(self.request())

    def test_review_must_bind_document_config_schema_and_all_sections(self):
        for field,bad in (("documentSha256","0"*64),("configVersion","other"),("schemaProfile","other"),("sectionCount",2)):
            self.s3.objects.pop(("reviews/tenant-a/x.json","r1"),None)
            value=review(self.snapshot_ref["sha256"]);value[field]=bad;self.review_ref=self.s3.add("reviews/tenant-a/x.json",value,"r1")
            with self.subTest(field=field):
                with self.assertRaises(handler.Rejected):self.call(self.request())

    def test_policy_must_be_approved_closed_and_human_review_required(self):
        bads=[{"approval":policy()["approval"]|{"status":"AWAITING_USER"}},{"requireHumanReviewForAll":False},{"automaticBusinessPersistenceAuthorized":True}]
        for mutation in bads:
            self.s3.objects.pop((self.cfg.policy_key,self.cfg.policy_version_id),None);p=policy()|mutation;pref=self.s3.add(self.cfg.policy_key,p,self.cfg.policy_version_id);cfg2=handler.Config(**(self.cfg.__dict__|{"policy_sha256":pref["sha256"]}))
            with self.assertRaises(handler.Rejected):handler.evaluate_request_core(self.request(),clients=handler.Clients(self.s3),cfg=cfg2)

    def test_unknown_class_duplicate_section_and_schema_failure_rejected(self):
        original=self.s3.objects[("snapshots/tenant-a/x.json","s1")]
        for doc in (snapshot()|{"sections":[snapshot()["sections"][0]|{"classification":"Other"}]},snapshot()|{"sections":snapshot()["sections"]*2}):
            self.s3.objects.pop(("snapshots/tenant-a/x.json","s1"),None);self.snapshot_ref=self.s3.add("snapshots/tenant-a/x.json",doc,"s1");self.review_ref=self.s3.add("reviews/tenant-a/x2.json",review(self.snapshot_ref["sha256"]),"r2")
            with self.assertRaises(handler.Rejected):self.call(self.request())
        self.s3.objects[("snapshots/tenant-a/x.json","s1")]=original

    def test_schema_mismatch_is_durable_negative_not_retry(self):
        with patch.object(handler,"validate",side_effect=handler.SchemaValidationError("schema mismatch")):
            out=self.call(self.request())
        self.assertFalse(out["evaluationPassed"])
        decision=[json.loads(v["Body"]) for (k,_),v in self.s3.objects.items() if k.startswith("decisions/")][0]
        self.assertEqual("SCHEMA_VALIDATION_FAILED",decision["decisionReason"])

    def test_cross_tenant_snapshot_review_and_review_tenant_are_rejected(self):
        snapshot_ref=dict(self.snapshot_ref);snapshot_ref["key"]="snapshots/tenant-b/x.json"
        with self.assertRaisesRegex(handler.Rejected,"outside the tenant prefix"):handler._read_ref(handler.Clients(self.s3),self.cfg,snapshot_ref,"snapshot")
        review_ref=dict(self.review_ref);review_ref["key"]="reviews/tenant-b/x.json"
        with self.assertRaisesRegex(handler.Rejected,"outside the tenant prefix"):handler._read_ref(handler.Clients(self.s3),self.cfg,review_ref,"review")
        self.s3.objects.pop(("reviews/tenant-a/x.json","r1"));self.review_ref=self.s3.add("reviews/tenant-a/x.json",review(self.snapshot_ref["sha256"])|{"tenantId":"tenant-b"},"r1")
        with self.assertRaises(handler.Rejected):self.call(self.request())

    def test_replay_reconciles_same_decision_and_divergent_rejects(self):
        first=self.call(self.request());second=self.call(self.request());self.assertEqual(first["decisionVersionId"],second["decisionVersionId"])
        key=handler._decision_key(self.cfg,self.request()["evaluationId"]);stored=[v for (k,_),v in self.s3.objects.items() if k==key][0];stored["Metadata"]={"sha256":"0"*64}
        with self.assertRaises(handler.Rejected):self.call(self.request())

    def test_extra_request_key_and_upstream_storage_authority_rejected(self):
        value=self.request();value["extra"]=1
        with self.assertRaises(handler.Rejected):self.call(value)
        value=self.request();value["automaticBusinessPersistenceAuthorized"]=True
        with self.assertRaises(handler.Rejected):self.call(value)

    def test_decision_contains_no_document_values_or_reviewer_identity(self):
        self.call(self.request());raw=b"\n".join(v["Body"] for (k,_),v in self.s3.objects.items() if k.startswith("decisions/"))
        self.assertNotIn(b"secret-file",raw);self.assertNotIn(b"10.00",raw);self.assertNotIn(b"owner-7",raw)

    def test_lambda_registers_context_and_uses_partial_batch(self):
        context=object();result=handler.lambda_handler({"Records":[]},context);self.assertTrue(result["called"]);self.assertIs(handler.idempotency.context,context)


if __name__=="__main__":unittest.main()
````

### FILE: `aws_idp_evaluation_decision_worker/evaluation-policy.template.json`
```yaml
block_id: "AWS-IDP-EVALUATION-DECISION-WORKER:evaluation-policy-template-json:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-evaluation-decision-worker/evaluation-policy.template.json"
license: "LicenseRef-Workspace-Owner"
sha256: "6a1d92736fa3c49dfde1491c0df57699d4308f9dfa3c557d5503e3cdb5048fcd"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-aws-idp-runtime-evaluation-policy/v1",
  "approval": {"status": "AWAITING_USER", "approvalId": "REQUIRED", "approvedBy": "REQUIRED", "approvedAt": "REQUIRED"},
  "configVersion": "REPLACE_WITH_PINNED_CONFIG_VERSION",
  "schemaProfile": "REPLACE_WITH_PINNED_SCHEMA_PROFILE",
  "requireHumanReviewForAll": true,
  "classSchemas": {},
  "automaticBusinessPersistenceAuthorized": false
}
````

### FILE: `aws_idp_evaluation_decision_worker/provider-profile.template.json`
```yaml
block_id: "AWS-IDP-EVALUATION-DECISION-WORKER:provider-profile-template-json:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-evaluation-decision-worker/provider-profile.template.json"
license: "LicenseRef-Workspace-Owner"
sha256: "be795f524c12533ad659549f8872953221bbcec90c2fb411dfa75c8c0b15e542"
variables: []
secrets_allowed: false
```
````json
{
  "acknowledgeConditionedState": false,
  "tenantId": "REPLACE_WITH_TENANT",
  "idpConfigVersion": "REPLACE_WITH_PINNED_CONFIG_VERSION",
  "schemaProfile": "REPLACE_WITH_PINNED_SCHEMA_PROFILE",
  "evaluationQueueArn": "REPLACE_WITH_QUEUE_ARN",
  "resultBucket": "REPLACE_WITH_RESULT_BUCKET",
  "resultBucketArn": "REPLACE_WITH_RESULT_BUCKET_ARN",
  "resultKmsKeyArn": "REPLACE_WITH_KMS_KEY_ARN",
  "logKmsKeyArn": "REPLACE_WITH_KMS_KEY_ARN",
  "powertoolsLayerArn": "REPLACE_WITH_OFFICIAL_POWERTOOLS_LAYER_ARN",
  "policyKey": "REPLACE_WITH_IMMUTABLE_POLICY_KEY",
  "policyVersionId": "REPLACE_WITH_POLICY_VERSION",
  "policySha256": "REQUIRED",
  "retainDays": 0,
  "maxObjectBytes": 0,
  "idempotencySeconds": 0,
  "owners": {"documents": "", "security": "", "privacy": "", "operations": "", "recovery": "", "cost": ""},
  "evidence": {"awsIdentity": "", "powertoolsRelease": "", "policyApproval": "", "corpusEvaluation": "", "hitl": "", "kmsObjectLock": "", "duplicateRace": "", "dlqReplay": "", "restore": "", "costApproval": ""}
}
````

### FILE: `aws_idp_evaluation_decision_worker/template.yaml`
```yaml
block_id: "AWS-IDP-EVALUATION-DECISION-WORKER:template-yaml:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/aws-powertools/powertools-lambda-python/tree/376757161b002f0c2f5d19d5cdf9def5b8c45704/examples"
license: "MIT-0 AND LicenseRef-Workspace-Owner"
sha256: "d9fdcfcd57ceabae0db5b8a271b67505c62a66253b08c668da24514a2d7f46d8"
variables: []
secrets_allowed: false
```
````yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Description: Idempotent AWS IDP evaluation decision worker using official AWS Powertools primitives

Parameters:
  EvaluationQueueArn: {Type: String, AllowedPattern: '^arn:aws[a-zA-Z-]*:sqs:[a-z0-9-]+:[0-9]{12}:[A-Za-z0-9_-]+\.fifo$'}
  ResultBucketName: {Type: String, MinLength: 3}
  ResultBucketArn: {Type: String, AllowedPattern: '^arn:aws[a-zA-Z-]*:s3:::[A-Za-z0-9.-]+$'}
  ResultKmsKeyArn: {Type: String, AllowedPattern: '^arn:aws[a-zA-Z-]*:kms:[a-z0-9-]+:[0-9]{12}:key/[0-9a-f-]+$'}
  LogKmsKeyArn: {Type: String, AllowedPattern: '^arn:aws[a-zA-Z-]*:kms:[a-z0-9-]+:[0-9]{12}:key/[0-9a-f-]+$'}
  PowertoolsLayerArn: {Type: String, AllowedPattern: '^arn:aws[a-zA-Z-]*:lambda:[a-z0-9-]+:[0-9]{12}:layer:[A-Za-z0-9-_]+:[0-9]+$'}
  TenantId: {Type: String, AllowedPattern: '^[A-Za-z0-9][A-Za-z0-9._-]{0,62}$'}
  IdpConfigVersion: {Type: String, AllowedPattern: '^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$'}
  SchemaProfile: {Type: String, AllowedPattern: '^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$'}
  PolicyKey: {Type: String, MinLength: 1, MaxLength: 1024}
  PolicyVersionId: {Type: String, MinLength: 1, MaxLength: 1024}
  PolicySha256: {Type: String, AllowedPattern: '^[0-9a-f]{64}$'}
  RetainDays: {Type: Number, Default: 365, MinValue: 1, MaxValue: 36500}
  MaxObjectBytes: {Type: Number, Default: 20971520, MinValue: 1024, MaxValue: 52428800}
  IdempotencySeconds: {Type: Number, Default: 315360000, MinValue: 86400, MaxValue: 3153600000}

Resources:
  IdempotencyTable:
    Type: AWS::DynamoDB::Table
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties:
      BillingMode: PAY_PER_REQUEST
      AttributeDefinitions: [{AttributeName: id, AttributeType: S}]
      KeySchema: [{AttributeName: id, KeyType: HASH}]
      TimeToLiveSpecification: {AttributeName: expiration, Enabled: true}
      PointInTimeRecoverySpecification: {PointInTimeRecoveryEnabled: true}
      SSESpecification: {SSEEnabled: true, SSEType: KMS, KMSMasterKeyId: !Ref ResultKmsKeyArn}
  WorkerRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument: {Version: '2012-10-17', Statement: [{Effect: Allow, Principal: {Service: lambda.amazonaws.com}, Action: 'sts:AssumeRole'}]}
      Policies:
        - PolicyName: ExactDecisionWorkerEffects
          PolicyDocument:
            Version: '2012-10-17'
            Statement:
              - {Effect: Allow, Action: ['logs:CreateLogStream', 'logs:PutLogEvents'], Resource: !Sub 'arn:${AWS::Partition}:logs:${AWS::Region}:${AWS::AccountId}:log-group:/aws/lambda/${AWS::StackName}-idp-evaluation-decision:*'}
              - {Effect: Allow, Action: ['sqs:ReceiveMessage', 'sqs:DeleteMessage', 'sqs:GetQueueAttributes'], Resource: !Ref EvaluationQueueArn}
              - {Effect: Allow, Action: ['dynamodb:GetItem', 'dynamodb:PutItem', 'dynamodb:UpdateItem', 'dynamodb:DeleteItem'], Resource: !GetAtt IdempotencyTable.Arn}
              - {Effect: Allow, Action: ['s3:GetObject', 's3:GetObjectVersion', 's3:GetObjectRetention', 's3:PutObject', 's3:PutObjectRetention'], Resource: !Sub '${ResultBucketArn}/*'}
              - {Effect: Allow, Action: ['kms:Decrypt', 'kms:Encrypt', 'kms:GenerateDataKey', 'kms:DescribeKey'], Resource: !Ref ResultKmsKeyArn}
  WorkerFunction:
    Type: AWS::Serverless::Function
    Properties:
      CodeUri: .
      Handler: handler.lambda_handler
      FunctionName: !Sub '${AWS::StackName}-idp-evaluation-decision'
      Runtime: python3.13
      MemorySize: 1024
      Timeout: 300
      ReservedConcurrentExecutions: 20
      Role: !GetAtt WorkerRole.Arn
      Layers: [!Ref PowertoolsLayerArn]
      Environment:
        Variables:
          POWERTOOLS_SERVICE_NAME: elite-idp-evaluation-decision
          IDEMPOTENCY_TABLE: !Ref IdempotencyTable
          TENANT_ID: !Ref TenantId
          IDP_CONFIG_VERSION: !Ref IdpConfigVersion
          SCHEMA_PROFILE: !Ref SchemaProfile
          RESULT_BUCKET: !Ref ResultBucketName
          RESULT_KMS_KEY_ARN: !Ref ResultKmsKeyArn
          POLICY_KEY: !Ref PolicyKey
          POLICY_VERSION_ID: !Ref PolicyVersionId
          POLICY_SHA256: !Ref PolicySha256
          RETAIN_DAYS: !Ref RetainDays
          MAX_OBJECT_BYTES: !Ref MaxObjectBytes
          IDEMPOTENCY_SECONDS: !Ref IdempotencySeconds
      Events:
        EvaluationBatch:
          Type: SQS
          Properties:
            Queue: !Ref EvaluationQueueArn
            BatchSize: 10
            FunctionResponseTypes: [ReportBatchItemFailures]
            ScalingConfig: {MaximumConcurrency: 20}
  WorkerLogGroup:
    Type: AWS::Logs::LogGroup
    DeletionPolicy: Retain
    UpdateReplacePolicy: Retain
    Properties: {LogGroupName: !Sub '/aws/lambda/${AWS::StackName}-idp-evaluation-decision', RetentionInDays: 30, KmsKeyId: !Ref LogKmsKeyArn}

Outputs:
  WorkerFunctionArn: {Value: !GetAtt WorkerFunction.Arn}
  IdempotencyTableName: {Value: !Ref IdempotencyTable}
````

### FILE: `aws_idp_evaluation_decision_worker/deploy.ps1`
```yaml
block_id: "AWS-IDP-EVALUATION-DECISION-WORKER:deploy-ps1:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-evaluation-decision-worker/deploy.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "a3823849ba2a902cc58c5f57b84110e70ac20214cdbacb9dc620dfce5b1c4c6d"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding(SupportsShouldProcess=$true, ConfirmImpact='High')]
param([Parameter(Mandatory)][string]$ProfilePath,[Parameter(Mandatory)][string]$StackName,[switch]$Execute)
$ErrorActionPreference='Stop'
$p=Get-Content -LiteralPath $ProfilePath -Raw|ConvertFrom-Json -Depth 64
if($p.acknowledgeConditionedState-ne$true){throw 'Conditioned state is not acknowledged.'}
foreach($n in @('documents','security','privacy','operations','recovery','cost')){if([string]::IsNullOrWhiteSpace($p.owners.$n)){throw "Missing owner: $n"}}
foreach($n in @('awsIdentity','powertoolsRelease','policyApproval','corpusEvaluation','hitl','kmsObjectLock','duplicateRace','dlqReplay','restore','costApproval')){if([string]::IsNullOrWhiteSpace($p.evidence.$n)){throw "Missing evidence: $n"}}
if($p.policySha256-notmatch'^[0-9a-f]{64}$'){throw 'Policy SHA-256 is invalid.'}
if($p.retainDays-lt1-or$p.maxObjectBytes-lt1024-or$p.idempotencySeconds-lt86400){throw 'Limits are unresolved.'}
if(-not$Execute){throw 'Dry gate complete. Re-run with -Execute and approve explicitly.'}
$args=@("EvaluationQueueArn=$($p.evaluationQueueArn)","ResultBucketName=$($p.resultBucket)","ResultBucketArn=$($p.resultBucketArn)","ResultKmsKeyArn=$($p.resultKmsKeyArn)","LogKmsKeyArn=$($p.logKmsKeyArn)","PowertoolsLayerArn=$($p.powertoolsLayerArn)","TenantId=$($p.tenantId)","IdpConfigVersion=$($p.idpConfigVersion)","SchemaProfile=$($p.schemaProfile)","PolicyKey=$($p.policyKey)","PolicyVersionId=$($p.policyVersionId)","PolicySha256=$($p.policySha256)","RetainDays=$($p.retainDays)","MaxObjectBytes=$($p.maxObjectBytes)","IdempotencySeconds=$($p.idempotencySeconds)")
if(-not$PSCmdlet.ShouldProcess($StackName,'deploy AWS IDP evaluation decision worker')){return}
& sam deploy --template-file (Join-Path $PSScriptRoot 'template.yaml') --stack-name $StackName --capabilities CAPABILITY_IAM --parameter-overrides @args
if($LASTEXITCODE-ne0){throw "sam deploy failed: $LASTEXITCODE"}
````

### FILE: `aws_idp_evaluation_decision_worker/verify_contract.ps1`
```yaml
block_id: "AWS-IDP-EVALUATION-DECISION-WORKER:verify-contract-ps1:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-evaluation-decision-worker/verify_contract.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "1563020506885f7518a6f90a83ea7cbbe6a7a2e425c15b7543f92b20ce0f56da"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]param([string]$ComponentRoot=$PSScriptRoot)
$ErrorActionPreference='Stop'
$files=@('README.md','handler.py','test_handler.py','evaluation-policy.template.json','provider-profile.template.json','template.yaml','deploy.ps1','verify_contract.ps1')
foreach($f in $files){if(-not(Test-Path -LiteralPath (Join-Path $ComponentRoot $f)-PathType Leaf)){throw "Missing $f"}}
$p=Get-Content -LiteralPath (Join-Path $ComponentRoot 'provider-profile.template.json')-Raw|ConvertFrom-Json -Depth 64
if($p.acknowledgeConditionedState-ne$false-or$p.retainDays-ne0-or$p.maxObjectBytes-ne0-or$p.idempotencySeconds-ne0){throw 'Profile must begin fail closed.'}
$policy=Get-Content -LiteralPath (Join-Path $ComponentRoot 'evaluation-policy.template.json')-Raw|ConvertFrom-Json -Depth 64
if($policy.approval.status-ne'AWAITING_USER'-or$policy.requireHumanReviewForAll-ne$true-or$policy.automaticBusinessPersistenceAuthorized-ne$false-or@($policy.classSchemas.PSObject.Properties).Count-ne0){throw 'Policy template must begin fail closed.'}
$h=Get-Content -LiteralPath (Join-Path $ComponentRoot 'handler.py')-Raw
foreach($n in @('BatchProcessor','process_partial_response','idempotent_function','raise_on_no_idempotency_key=True','payload_validation_jmespath','hash_function="sha256"','VersionId=version','ObjectLockMode="COMPLIANCE"','automaticBusinessPersistenceAuthorized')){if(-not$h.Contains($n)){throw "Handler contract missing $n"}}
if($h-match'print\('-or$h-match'logger\.'){throw 'Handler may not log documents or identifiers.'}
Push-Location $ComponentRoot
try{python -B -c "from pathlib import Path;[compile(Path(p).read_text(encoding='utf-8'),p,'exec') for p in ('handler.py','test_handler.py')]";if($LASTEXITCODE-ne0){throw 'compile failed'};python -B -m unittest -v test_handler.py;if($LASTEXITCODE-ne0){throw 'tests failed'}}finally{Pop-Location}
'AWS_IDP_EVALUATION_DECISION_WORKER_PASS'
````

## 6. Configuration surface

Profile y policy empiezan cerrados: acknowledgement false, owners/evidencia vacíos, policy `AWAITING_USER`, schemas vacíos y límites cero. Deben fijarse tenant/config/schema, queue ARN, result bucket/ARN/KMS, Powertools layer ARN, policy key/VersionId/SHA, retención, tamaño e idempotency window.

## 7. Dependency bill

- AWS Powertools for Lambda Python 3.34.0, release/commit/archive/wheel ya fijados por el pack oficial componente, MIT-0.
- Boto3/Botocore del runtime/layer AWS; Lambda Python 3.13; SQS FIFO; DynamoDB; S3 Object Lock/KMS; AWS SAM.
- Seis archivos `AUTHORED`, dos `ADAPTED`, cero `VERBATIM`; ninguna composición local se atribuye a AWS.

## 8. Apply order

1. Adquirir Powertools 3.34.0 y validar su pack oficial.
2. Evaluar schemas/corpus con strict-field gate; aprobar policy y subirla create-only al bucket Object Lock.
3. Completar profile/owners/evidencia; verificar tests/cfn-lint y aprobar efectos/costo.
4. Desplegar después del handoff 0.1.2, usando su queue/bucket/KMS exactos.
5. Probar duplicates/races/tamper/skipped/DLQ/redrive/load/restore/rollback en DEV.
6. Conectar una persistencia sólo mediante transacción negocio+outbox separada que exija decisión positiva.

## 9. Verification

El verifier exige 8 archivos, templates fail-closed, primitives Powertools/VersionId/ObjectLock/no-storage, compile y 13 pruebas. Admisión exige cfn-lint, round-trip, perfil combinado y gates globales. Live AWS, layer/account/IAM/KMS, corpus/policy/HITL, carga/DLQ/restore y persistencia transaccional siguen condicionados.

## 10. Reconstruction evidence

Reconstruido el 2026-08-28 contra AWS Powertools 3.34.0 commit `376757161b002f0c2f5d19d5cdf9def5b8c45704` y AWS IDP handoff 0.1.2. V95 registra el receipt v2 con snapshot reference; V94 conserva hashes y 13 pruebas de aislamiento tenant, invalidación durable de schema, cfn-lint, findings upstream y límites. No se afirma deploy live ni autorización empresarial.
