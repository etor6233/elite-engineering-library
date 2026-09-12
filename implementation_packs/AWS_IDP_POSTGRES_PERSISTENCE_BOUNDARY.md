# AWS IDP PostgreSQL Persistence Boundary

## 1. Metadata

```yaml
pack_id: "AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el límite S3 EventBridge/SQS/Lambda que verifica una decisión AWS IDP v2 y su snapshot inmutable, exige policy de persistencia aprobada y escribe document intake más outbox en una única transacción Aurora PostgreSQL Data API."
stacks: ["AWS Lambda Python 3.13", "AWS Powertools 3.34.0", "Amazon S3 EventBridge", "Amazon SQS", "Aurora PostgreSQL RDS Data API", "PostgreSQL", "AWS SAM"]
compatible_with: ["AWS-IDP-EVALUATION-DECISION-WORKER 0.2.x", "PG-TX-FOUNDATION 0.1.x", "AWS-POWERTOOLS-IDEMPOTENT-SQS-BATCH-COMPONENT 0.1.x", "DEBEZIUM-POSTGRES-OUTBOX-RUNTIME 0.1.x"]
incompatible_with: ["decision sin snapshotReference v2", "policy no aprobada o móvil", "RDS no Aurora Data API", "secret DB table-owner para runtime", "outbox exactly-once", "persistencia sin RLS/duplicate/race/restore evidence"]
license_expression: "LicenseRef-Workspace-Owner AND MIT-0"
upstream_sources: ["https://docs.aws.amazon.com/AmazonS3/latest/userguide/ev-events.html", "https://docs.aws.amazon.com/prescriptive-guidance/latest/lambda-event-filtering-partial-batch-responses-for-sqs/best-practices-partial-batch-responses.html", "https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/data-api-operations.html", "https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/transactional-outbox.html", "https://github.com/aws-powertools/powertools-lambda-python/tree/376757161b002f0c2f5d19d5cdf9def5b8c45704", "https://github.com/dapr/dapr/tree/2dcf2e3548faf3c740b3d29ee300e2132db65cff"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use después del decision worker 0.2.0 y antes de cualquier mapping ERP/CRM. Requiere foundation PostgreSQL, Aurora Data API, policy Object-Locked aprobada, S3 EventBridge habilitado y evidencias de RLS, duplicados, outbox recovery, restore, rollback y costo. Persiste staging revisado; no convierte atributos genéricos en asientos o inventario final.

## 3. Architecture contract

Amazon S3 entrega la versión exacta del decision object a EventBridge; SQS aporta buffer/DLQ y Lambda responde partial batch. El handler ignora de forma exitosa decisiones negativas ya durables. Para positivas verifica account/region/bucket/prefix/version, KMS/Object Lock/checksum/hash, decision v2, review authority, policy aprobada, profile/classes y snapshot exacto. Una sola transacción Data API fija tenant RLS, inserta o reconcilia el intake, compara hashes, inserta/reconcilia una outbox append-only con columnas oficiales Debezium y commit; toda excepción hace rollback. RLS separa writer tenant y rol NOLOGIN CDC sólo lectura. La entrega sigue siendo at-least-once.

## 4. Exact file manifest

```text
CREATE aws_idp_postgres_persistence_boundary/README.md
CREATE aws_idp_postgres_persistence_boundary/handler.py
CREATE aws_idp_postgres_persistence_boundary/test_handler.py
CREATE aws_idp_postgres_persistence_boundary/001_document_intake_outbox.sql
CREATE aws_idp_postgres_persistence_boundary/persistence-policy.template.json
CREATE aws_idp_postgres_persistence_boundary/project-profile.template.json
CREATE aws_idp_postgres_persistence_boundary/official-contract-lock.json
CREATE aws_idp_postgres_persistence_boundary/template.yaml
CREATE aws_idp_postgres_persistence_boundary/apply_migration.ps1
CREATE aws_idp_postgres_persistence_boundary/deploy.ps1
CREATE aws_idp_postgres_persistence_boundary/verify_contract.ps1
```

## 5. Materialization blocks

### FILE: `aws_idp_postgres_persistence_boundary/README.md`
```yaml
block_id: "AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-postgres-persistence-boundary/README.md"
license: "LicenseRef-Workspace-Owner"
sha256: "55631f6a99f8bd6e3481be5b65c91c09f548f91800fd9e843e3c6deb23adc1b7"
variables: []
secrets_allowed: false
```
````markdown
# AWS IDP to PostgreSQL persistence boundary

This conditioned component consumes versioned S3 `Object Created` events for immutable AWS IDP evaluation decisions. It follows the published Amazon S3 EventBridge envelope, AWS Lambda SQS partial-batch guidance, Aurora RDS Data API transaction APIs, AWS transactional-outbox guidance, and the separately preserved Dapr 1.18.3 outbox implementation/tests. The integration code is locally authored/adapted and is not falsely attributed to AWS, PostgreSQL, Dapr, or Red Hat.

A v2 positive decision is accepted only when it binds the configured tenant/config, completed human-review authority, immutable snapshot VersionId/SHA/KMS/Object Lock, approved persistence policy and admitted schema profile/classifications. Negative decisions are intentionally acknowledged without a database call because their immutable decision receipt is already the durable rejection authority.

For a positive decision, one RDS Data API transaction sets the tenant context, inserts the validated document intake row idempotently, locks and compares every authoritative hash, inserts the append-only `document_intelligence.outbox_event` shaped to Debezium's official `id/aggregatetype/aggregateid/type/payload` contract, verifies the outbox payload, and commits. Any divergence or infrastructure error rolls back. Duplicate delivery reconciles against the existing row and outbox. Debezium delivery remains at-least-once; downstream consumers still require inbox/idempotency.

Before deployment:

1. materialize and apply `POSTGRES-TRANSACTIONAL-FOUNDATION` before `001_document_intake_outbox.sql`;
2. create a least-privilege Aurora PostgreSQL Data API secret/role that is not the table owner and prove RLS isolation; separately grant the Debezium login membership in the migration-created NOLOGIN role `elite_outbox_cdc`, which has only schema usage and outbox select;
3. complete and independently approve `persistence-policy.template.json`, upload it create-only to the Object Lock bucket and pin key/VersionId/SHA-256;
4. enable S3→EventBridge on the existing result bucket, complete every project profile field and retain migration/duplicate/race/outbox recovery/restore/rollback/cost evidence;
5. run verifier, cfn-lint, SAM build and DEV duplicate/tamper/failover tests before any production authorization.

This component stores the reviewed snapshot as JSONB staging evidence, not arbitrary final ERP mutations. Project-specific mappings must consume `document.intake.accepted` through an idempotent inbox and their own domain transaction. No extracted data, reviewer identity or secret is logged by this handler.
````

### FILE: `aws_idp_postgres_persistence_boundary/handler.py`
```yaml
block_id: "AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY:handler:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/aws-powertools/powertools-lambda-python/blob/376757161b002f0c2f5d19d5cdf9def5b8c45704/examples/batch_processing/sam/sqs_batch_processing.py"
license: "MIT-0 AND LicenseRef-Workspace-Owner"
sha256: "c89b656a10c87249a8e2b7194e633f6cc530a9c30054a0ae8c3b6c5bce6a8b38"
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
from dataclasses import dataclass
from datetime import datetime
from typing import Any

import boto3
from aws_lambda_powertools.utilities.batch import BatchProcessor, EventType, process_partial_response


DECISION_SCHEMA = "elite-aws-idp-evaluation-decision/v2"
POLICY_SCHEMA = "elite-document-persistence-policy/v1"
SNAPSHOT_SCHEMA = "elite-aws-idp-document-snapshot/v1"
HEX = re.compile(r"^[0-9a-f]{64}$")
TOKEN = re.compile(r"^[A-Za-z0-9][A-Za-z0-9._:/+=,@-]{0,511}$")
UUID = re.compile(r"^[0-9a-f]{8}-[0-9a-f]{4}-[1-8][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$")
processor = BatchProcessor(event_type=EventType.SQS)


class Rejected(ValueError):
    pass


@dataclass(frozen=True)
class Config:
    source_tenant_id: str
    business_tenant_id: str
    config_version: str
    account_id: str
    region: str
    result_bucket: str
    result_kms_key_arn: str
    policy_key: str
    policy_version_id: str
    policy_sha256: str
    resource_arn: str
    secret_arn: str
    database: str
    max_object_bytes: int

    @staticmethod
    def from_env() -> "Config":
        value = Config(
            source_tenant_id=os.environ.get("SOURCE_TENANT_ID", ""),
            business_tenant_id=os.environ.get("BUSINESS_TENANT_ID", ""),
            config_version=os.environ.get("IDP_CONFIG_VERSION", ""),
            account_id=os.environ.get("SOURCE_ACCOUNT_ID", ""),
            region=os.environ.get("SOURCE_REGION", ""),
            result_bucket=os.environ.get("RESULT_BUCKET", ""),
            result_kms_key_arn=os.environ.get("RESULT_KMS_KEY_ARN", ""),
            policy_key=os.environ.get("PERSISTENCE_POLICY_KEY", ""),
            policy_version_id=os.environ.get("PERSISTENCE_POLICY_VERSION_ID", ""),
            policy_sha256=os.environ.get("PERSISTENCE_POLICY_SHA256", ""),
            resource_arn=os.environ.get("RDS_RESOURCE_ARN", ""),
            secret_arn=os.environ.get("RDS_SECRET_ARN", ""),
            database=os.environ.get("RDS_DATABASE", ""),
            max_object_bytes=int(os.environ.get("MAX_OBJECT_BYTES", "20971520")),
        )
        value.validate()
        return value

    def validate(self) -> None:
        if not TOKEN.fullmatch(self.source_tenant_id) or not TOKEN.fullmatch(self.config_version):
            raise Rejected("source identity configuration is invalid")
        if not UUID.fullmatch(self.business_tenant_id):
            raise Rejected("business tenant must be a canonical UUID")
        if not re.fullmatch(r"[0-9]{12}", self.account_id) or not TOKEN.fullmatch(self.region):
            raise Rejected("AWS source authority is invalid")
        if not self.result_bucket or not self.result_kms_key_arn or not self.policy_version_id:
            raise Rejected("storage authority is incomplete")
        if not self.policy_key.startswith(f"persistence-policies/{self.source_tenant_id}/") or ".." in self.policy_key.split("/"):
            raise Rejected("policy key is outside the tenant prefix")
        if not HEX.fullmatch(self.policy_sha256):
            raise Rejected("policy hash is invalid")
        if not self.resource_arn or not self.secret_arn or not TOKEN.fullmatch(self.database):
            raise Rejected("RDS Data API authority is incomplete")
        if not 1024 <= self.max_object_bytes <= 50 * 1024 * 1024:
            raise Rejected("object limit is invalid")


@dataclass(frozen=True)
class Clients:
    s3: Any
    rds: Any


def _clients() -> Clients:
    return Clients(s3=boto3.client("s3"), rds=boto3.client("rds-data"))


def _sha(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def _checksum(data: bytes) -> str:
    return base64.b64encode(hashlib.sha256(data).digest()).decode("ascii")


def _pairs(pairs: list[tuple[str, Any]]) -> dict[str, Any]:
    value: dict[str, Any] = {}
    for key, item in pairs:
        if key in value:
            raise Rejected("duplicate JSON key")
        value[key] = item
    return value


def _json(raw: bytes, label: str) -> dict[str, Any]:
    try:
        value = json.loads(raw.decode("utf-8"), object_pairs_hook=_pairs)
    except (UnicodeDecodeError, json.JSONDecodeError) as error:
        raise Rejected(f"{label} is not strict UTF-8 JSON") from error
    if not isinstance(value, dict):
        raise Rejected(f"{label} must be an object")
    return value


def _canonical(value: Any) -> bytes:
    try:
        return json.dumps(value, sort_keys=True, separators=(",", ":"), ensure_ascii=False, allow_nan=False).encode("utf-8")
    except (TypeError, ValueError) as error:
        raise Rejected("value is not canonical JSON") from error


def _read_ref(clients: Clients, cfg: Config, ref: dict[str, Any], label: str) -> tuple[dict[str, Any], bytes]:
    if set(ref) != {"bucket", "key", "versionId", "sha256"} or ref.get("bucket") != cfg.result_bucket:
        raise Rejected(f"{label} reference is invalid")
    prefixes = {
        "snapshot": f"snapshots/{cfg.source_tenant_id}/",
        "policy": f"persistence-policies/{cfg.source_tenant_id}/",
    }
    key, version, expected = ref.get("key"), ref.get("versionId"), ref.get("sha256")
    if label not in prefixes or not isinstance(key, str) or not key.startswith(prefixes[label]) or ".." in key.split("/"):
        raise Rejected(f"{label} reference is outside the tenant prefix")
    if not isinstance(version, str) or not version or not isinstance(expected, str) or not HEX.fullmatch(expected):
        raise Rejected(f"{label} reference is incomplete")
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
    return _json(body, label), body


def _read_decision_event(clients: Clients, cfg: Config, event: dict[str, Any]) -> tuple[dict[str, Any], bytes, str]:
    required = {"version", "id", "detail-type", "source", "account", "time", "region", "resources", "detail"}
    if not required.issubset(event):
        raise Rejected("EventBridge envelope is incomplete")
    if event.get("version") != "0" or event.get("source") != "aws.s3" or event.get("detail-type") != "Object Created":
        raise Rejected("event type is not admitted")
    if event.get("account") != cfg.account_id or event.get("region") != cfg.region:
        raise Rejected("event AWS authority mismatch")
    detail = event.get("detail")
    if not isinstance(detail, dict) or str(detail.get("event-version", "")).split(".")[0] != "1" or detail.get("reason") != "PutObject":
        raise Rejected("S3 event detail is not admitted")
    bucket, obj = detail.get("bucket"), detail.get("object")
    if not isinstance(bucket, dict) or bucket.get("name") != cfg.result_bucket or not isinstance(obj, dict):
        raise Rejected("decision event bucket is invalid")
    key, version = obj.get("key"), obj.get("version-id")
    if not isinstance(key, str) or not key.startswith(f"decisions/{cfg.source_tenant_id}/") or ".." in key.split("/"):
        raise Rejected("decision event key is outside the tenant prefix")
    if not isinstance(version, str) or not version:
        raise Rejected("decision event has no version")
    response = clients.s3.get_object(Bucket=cfg.result_bucket, Key=key, VersionId=version, ChecksumMode="ENABLED")
    if response.get("VersionId") != version or response.get("ServerSideEncryption") != "aws:kms" or response.get("SSEKMSKeyId") != cfg.result_kms_key_arn or response.get("ObjectLockMode") != "COMPLIANCE":
        raise Rejected("decision storage authority mismatch")
    body = response["Body"].read(cfg.max_object_bytes + 1)
    digest = _sha(body)
    if len(body) > cfg.max_object_bytes or response.get("ChecksumSHA256") != _checksum(body) or (response.get("Metadata") or {}).get("sha256") != digest:
        raise Rejected("decision content integrity mismatch")
    return _json(body, "decision"), body, version


def _policy(clients: Clients, cfg: Config) -> dict[str, Any]:
    ref = {"bucket": cfg.result_bucket, "key": cfg.policy_key, "versionId": cfg.policy_version_id, "sha256": cfg.policy_sha256}
    value, _ = _read_ref(clients, cfg, ref, "policy")
    exact = {"schema", "tenantId", "businessTenantId", "configVersion", "approval", "allowedProfiles", "automaticPersistenceAuthorized"}
    if set(value) != exact or value.get("schema") != POLICY_SCHEMA or value.get("tenantId") != cfg.source_tenant_id or value.get("businessTenantId") != cfg.business_tenant_id or value.get("configVersion") != cfg.config_version:
        raise Rejected("persistence policy identity mismatch")
    approval = value.get("approval")
    if not isinstance(approval, dict) or set(approval) != {"status", "approvalId", "approvedBy", "approvedAt"} or approval.get("status") != "APPROVED":
        raise Rejected("persistence policy is not approved")
    if any(not isinstance(approval.get(key), str) or not approval[key].strip() or approval[key] == "REQUIRED" for key in ("approvalId", "approvedBy", "approvedAt")):
        raise Rejected("persistence policy approval evidence is unresolved")
    try:
        approved_at = datetime.fromisoformat(approval["approvedAt"].replace("Z", "+00:00"))
    except ValueError as error:
        raise Rejected("persistence policy approval time is invalid") from error
    if approved_at.tzinfo is None:
        raise Rejected("persistence policy approval time has no timezone")
    profiles = value.get("allowedProfiles")
    if value.get("automaticPersistenceAuthorized") is not True or not isinstance(profiles, dict) or not profiles:
        raise Rejected("persistence policy remains fail closed")
    for name, profile in profiles.items():
        if not TOKEN.fullmatch(name) or not isinstance(profile, dict) or set(profile) != {"domainType", "allowedClassifications"}:
            raise Rejected("persistence profile is invalid")
        if not TOKEN.fullmatch(str(profile.get("domainType", ""))) or not isinstance(profile.get("allowedClassifications"), list) or not profile["allowedClassifications"]:
            raise Rejected("persistence profile authority is incomplete")
        allowed = profile["allowedClassifications"]
        if any(not isinstance(item, str) or not TOKEN.fullmatch(item) for item in allowed) or len(allowed) != len(set(allowed)):
            raise Rejected("persistence profile classification is invalid")
    return profiles


def _validate_decision(decision: dict[str, Any], cfg: Config) -> bool:
    exact = {"schema", "evaluationId", "tenantId", "configVersion", "schemaProfile", "snapshotReference", "snapshotSha256", "policySha256", "reviewEvidenceSha256", "reviewAuthority", "sectionCount", "classifications", "evaluationPassed", "decisionReason", "automaticBusinessPersistenceAuthorized"}
    if set(decision) != exact or decision.get("schema") != DECISION_SCHEMA or decision.get("tenantId") != cfg.source_tenant_id or decision.get("configVersion") != cfg.config_version:
        raise Rejected("decision identity mismatch")
    if not isinstance(decision.get("evaluationId"), str) or not HEX.fullmatch(decision["evaluationId"]):
        raise Rejected("decision evaluation id is invalid")
    if decision.get("automaticBusinessPersistenceAuthorized") is not False:
        raise Rejected("decision attempted to grant persistence authority")
    if decision.get("evaluationPassed") is not True:
        return False
    if decision.get("decisionReason") != "SCHEMA_AND_HUMAN_REVIEW_VERIFIED" or decision.get("reviewAuthority") != "AWS_IDP_HUMAN_REVIEW_COMPLETED":
        raise Rejected("positive decision lacks completed review authority")
    ref = decision.get("snapshotReference")
    if not isinstance(ref, dict) or decision.get("snapshotSha256") != ref.get("sha256"):
        raise Rejected("decision snapshot binding is invalid")
    classes = decision.get("classifications")
    if not isinstance(classes, list) or not classes or classes != sorted(classes) or any(not isinstance(item, str) for item in classes):
        raise Rejected("decision classifications are invalid")
    if not isinstance(decision.get("sectionCount"), int) or decision["sectionCount"] <= 0:
        raise Rejected("decision section count is invalid")
    return True


def _validate_snapshot(snapshot: dict[str, Any], decision: dict[str, Any], cfg: Config) -> None:
    if snapshot.get("config_version") != cfg.config_version or snapshot.get("errors") or snapshot.get("processing_issues"):
        raise Rejected("snapshot is not admissible")
    sections = snapshot.get("sections")
    if not isinstance(sections, list) or len(sections) != decision["sectionCount"]:
        raise Rejected("snapshot section count mismatch")
    classes = sorted(section.get("classification") for section in sections if isinstance(section, dict))
    if len(classes) != len(sections) or classes != decision["classifications"]:
        raise Rejected("snapshot classification binding mismatch")


def _param(name: str, value: Any) -> dict[str, Any]:
    if isinstance(value, int):
        field = {"longValue": value}
    else:
        field = {"stringValue": str(value)}
    return {"name": name, "value": field}


def _field(value: dict[str, Any]) -> Any:
    for name in ("stringValue", "longValue", "booleanValue", "doubleValue"):
        if name in value:
            return value[name]
    if value.get("isNull") is True:
        return None
    raise Rejected("RDS Data API returned an unsupported field")


def _execute(clients: Clients, cfg: Config, transaction_id: str, sql: str, parameters: list[dict[str, Any]]) -> dict[str, Any]:
    return clients.rds.execute_statement(resourceArn=cfg.resource_arn, secretArn=cfg.secret_arn, database=cfg.database, transactionId=transaction_id, sql=sql, parameters=parameters)


def _persist(clients: Clients, cfg: Config, decision: dict[str, Any], decision_body: bytes, decision_version: str, snapshot: dict[str, Any], snapshot_body: bytes, domain_type: str) -> bool:
    decision_sha, snapshot_sha, payload_sha = _sha(decision_body), _sha(snapshot_body), _sha(_canonical(snapshot))
    common = [
        _param("tenant", cfg.business_tenant_id), _param("source_tenant", cfg.source_tenant_id),
        _param("evaluation", decision["evaluationId"]), _param("decision_sha", decision_sha),
        _param("decision_version", decision_version), _param("snapshot_sha", snapshot_sha),
        _param("policy_sha", cfg.policy_sha256), _param("config_version", cfg.config_version),
        _param("schema_profile", decision["schemaProfile"]), _param("domain_type", domain_type),
        _param("section_count", decision["sectionCount"]), _param("payload_sha", payload_sha),
        _param("payload", _canonical(snapshot).decode("utf-8")),
    ]
    begin = clients.rds.begin_transaction(resourceArn=cfg.resource_arn, secretArn=cfg.secret_arn, database=cfg.database)
    transaction_id = begin.get("transactionId")
    if not isinstance(transaction_id, str) or not transaction_id:
        raise Rejected("RDS Data API returned no transaction id")
    try:
        _execute(clients, cfg, transaction_id, "/* elite:set-tenant */ select set_config('app.tenant_id',:tenant,true)", [_param("tenant", cfg.business_tenant_id)])
        inserted = _execute(clients, cfg, transaction_id, """/* elite:insert-document */
insert into document_intelligence.intake_document(tenant_id,source_tenant_id,evaluation_id,decision_sha256_hex,decision_version_id,snapshot_sha256_hex,persistence_policy_sha256_hex,config_version,schema_profile,domain_type,section_count,document_payload_sha256_hex,document_payload)
values(cast(:tenant as uuid),:source_tenant,:evaluation,:decision_sha,:decision_version,:snapshot_sha,:policy_sha,:config_version,:schema_profile,:domain_type,:section_count,:payload_sha,cast(:payload as jsonb))
on conflict (tenant_id,evaluation_id) do nothing returning evaluation_id""", common)
        created = bool(inserted.get("records"))
        selected = _execute(clients, cfg, transaction_id, """/* elite:select-document */
select decision_sha256_hex,decision_version_id,snapshot_sha256_hex,persistence_policy_sha256_hex,config_version,schema_profile,domain_type,section_count,document_payload_sha256_hex
from document_intelligence.intake_document where tenant_id=cast(:tenant as uuid) and evaluation_id=:evaluation for update""", [_param("tenant", cfg.business_tenant_id), _param("evaluation", decision["evaluationId"])])
        records = selected.get("records") or []
        if len(records) != 1:
            raise Rejected("document idempotency authority is missing")
        actual = [_field(field) for field in records[0]]
        expected = [decision_sha, decision_version, snapshot_sha, cfg.policy_sha256, cfg.config_version, decision["schemaProfile"], domain_type, decision["sectionCount"], payload_sha]
        if actual != expected:
            raise Rejected("existing document intake is divergent")
        outbox_payload = _canonical({"evaluationId": decision["evaluationId"], "snapshotSha256": snapshot_sha, "domainType": domain_type}).decode("utf-8")
        _execute(clients, cfg, transaction_id, """/* elite:insert-outbox */
insert into document_intelligence.outbox_event(id,aggregatetype,aggregateid,type,payload,tenant_id,schema_version,occurred_at)
values(cast(substr(:evaluation,1,8)||'-'||substr(:evaluation,9,4)||'-7'||substr(:evaluation,14,3)||'-a'||substr(:evaluation,18,3)||'-'||substr(:evaluation,21,12) as uuid),'document_intake',:evaluation,'document.intake.accepted',cast(:outbox_payload as jsonb),cast(:tenant as uuid),1,clock_timestamp())
on conflict (tenant_id,aggregatetype,aggregateid,type) do nothing""", [_param("tenant", cfg.business_tenant_id), _param("evaluation", decision["evaluationId"]), _param("outbox_payload", outbox_payload)])
        outbox = _execute(clients, cfg, transaction_id, """/* elite:select-outbox */
select payload::text from document_intelligence.outbox_event where tenant_id=cast(:tenant as uuid) and aggregatetype='document_intake' and aggregateid=:evaluation and type='document.intake.accepted' for update""", [_param("tenant", cfg.business_tenant_id), _param("evaluation", decision["evaluationId"])])
        if len(outbox.get("records") or []) != 1 or json.loads(str(_field(outbox["records"][0][0]))) != json.loads(outbox_payload):
            raise Rejected("outbox authority is missing or divergent")
        clients.rds.commit_transaction(resourceArn=cfg.resource_arn, secretArn=cfg.secret_arn, transactionId=transaction_id)
        return created
    except Exception:
        clients.rds.rollback_transaction(resourceArn=cfg.resource_arn, secretArn=cfg.secret_arn, transactionId=transaction_id)
        raise


def process_event_core(event: dict[str, Any], *, clients: Clients, cfg: Config) -> dict[str, Any]:
    decision, decision_body, decision_version = _read_decision_event(clients, cfg, event)
    if not _validate_decision(decision, cfg):
        return {"evaluationId": decision["evaluationId"], "state": "IGNORED_NEGATIVE", "persisted": False}
    profiles = _policy(clients, cfg)
    schema_profile = decision.get("schemaProfile")
    if schema_profile not in profiles:
        raise Rejected("schema profile is not approved for persistence")
    profile = profiles[schema_profile]
    if not set(decision["classifications"]).issubset(set(profile["allowedClassifications"])):
        raise Rejected("decision classification is not approved for persistence")
    snapshot, snapshot_body = _read_ref(clients, cfg, decision["snapshotReference"], "snapshot")
    _validate_snapshot(snapshot, decision, cfg)
    created = _persist(clients, cfg, decision, decision_body, decision_version, snapshot, snapshot_body, profile["domainType"])
    return {"evaluationId": decision["evaluationId"], "state": "PERSISTED" if created else "RECONCILED", "persisted": True}


def _record_handler(record: dict[str, Any]) -> dict[str, Any]:
    body = record.get("body")
    if not isinstance(body, str):
        raise Rejected("SQS body is missing")
    return process_event_core(_json(body.encode("utf-8"), "SQS body"), clients=_clients(), cfg=Config.from_env())


def lambda_handler(event: dict[str, Any], context: Any) -> dict[str, Any]:
    return process_partial_response(event=event, record_handler=_record_handler, processor=processor, context=context)
````

### FILE: `aws_idp_postgres_persistence_boundary/test_handler.py`
```yaml
block_id: "AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY:test-handler:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-postgres-persistence-boundary/test_handler.py"
license: "LicenseRef-Workspace-Owner"
sha256: "f5c475c50b45cd5a617b0c157371b9a0b95b1f75bad36dc1a7621d3d63e63c67"
variables: []
secrets_allowed: false
```
````python
import base64,hashlib,io,json,sys,types,unittest

boto3=types.ModuleType("boto3");boto3.client=lambda name:None;sys.modules["boto3"]=boto3
batch=types.ModuleType("aws_lambda_powertools.utilities.batch")
class BatchProcessor:
    def __init__(self,**kwargs):self.kwargs=kwargs
class EventType:SQS="SQS"
def process_partial_response(**kwargs):return {"called":True,"context":kwargs.get("context")}
batch.BatchProcessor=BatchProcessor;batch.EventType=EventType;batch.process_partial_response=process_partial_response
for name in ("aws_lambda_powertools","aws_lambda_powertools.utilities"):
    sys.modules[name]=types.ModuleType(name)
sys.modules["aws_lambda_powertools.utilities.batch"]=batch

import handler


def canonical(value):return json.dumps(value,sort_keys=True,separators=(",",":"),ensure_ascii=False).encode()
def sha(value):return hashlib.sha256(value).hexdigest()
def checksum(value):return base64.b64encode(hashlib.sha256(value).digest()).decode()


class FakeS3:
    def __init__(self,kms):self.kms=kms;self.objects={}
    def add(self,key,value,version):
        body=canonical(value);digest=sha(body);self.objects[(key,version)]={"Body":body,"VersionId":version,"ServerSideEncryption":"aws:kms","SSEKMSKeyId":self.kms,"ObjectLockMode":"COMPLIANCE","ChecksumSHA256":checksum(body),"Metadata":{"sha256":digest}}
        return {"bucket":"result","key":key,"versionId":version,"sha256":digest}
    def get_object(self,**kwargs):
        value=dict(self.objects[(kwargs["Key"],kwargs["VersionId"])]);value["Body"]=io.BytesIO(value["Body"]);return value


def params(items):
    result={}
    for item in items:
        value=item["value"]
        result[item["name"]]=next(iter(value.values()))
    return result


class FakeRDS:
    def __init__(self):self.documents={};self.outbox={};self.commits=0;self.rollbacks=0;self.fail_outbox=False
    def begin_transaction(self,**kwargs):return {"transactionId":"tx-1"}
    def execute_statement(self,**kwargs):
        sql=kwargs["sql"];p=params(kwargs.get("parameters",[]));key=(p.get("tenant"),p.get("evaluation"))
        if "elite:set-tenant" in sql:return {"records":[[{"stringValue":p["tenant"]}]]}
        if "elite:insert-document" in sql:
            if key in self.documents:return {"records":[]}
            self.documents[key]=[p[n] for n in ("decision_sha","decision_version","snapshot_sha","policy_sha","config_version","schema_profile","domain_type","section_count","payload_sha")]
            return {"records":[[{"stringValue":p["evaluation"]}]]}
        if "elite:select-document" in sql:
            values=self.documents.get(key);return {"records":[] if values is None else [[({"longValue":v} if isinstance(v,int) else {"stringValue":v}) for v in values]]}
        if "elite:insert-outbox" in sql:
            if self.fail_outbox:raise RuntimeError("outbox unavailable")
            self.outbox.setdefault(key,p["outbox_payload"]);return {"numberOfRecordsUpdated":1}
        if "elite:select-outbox" in sql:
            value=self.outbox.get(key);return {"records":[] if value is None else [[{"stringValue":value}]]}
        raise AssertionError(sql)
    def commit_transaction(self,**kwargs):self.commits+=1;return {"transactionStatus":"Transaction Committed"}
    def rollback_transaction(self,**kwargs):self.rollbacks+=1;return {"transactionStatus":"Rollback Complete"}


def snapshot():
    return {"config_version":"cfg-7","errors":[],"processing_issues":[],"sections":[{"section_id":"1","classification":"Invoice","attributes":{"invoiceId":"A-7","total":"10.00"}}]}


def policy():
    return {"schema":handler.POLICY_SCHEMA,"tenantId":"tenant-a","businessTenantId":"018f4d4a-7b36-7a21-8d10-2f4c54c28b10","configVersion":"cfg-7","approval":{"status":"APPROVED","approvalId":"approval-7","approvedBy":"owner-7","approvedAt":"2026-08-28T10:00:00Z"},"allowedProfiles":{"invoice-v3":{"domainType":"supplier_invoice","allowedClassifications":["Invoice"]}},"automaticPersistenceAuthorized":True}


class BoundaryTests(unittest.TestCase):
    def setUp(self):
        self.cfg=handler.Config("tenant-a","018f4d4a-7b36-7a21-8d10-2f4c54c28b10","cfg-7","111122223333","us-east-1","result","kms","persistence-policies/tenant-a/policy.json","p1","0"*64,"cluster","secret","elite",1024*1024)
        self.s3=FakeS3("kms");self.rds=FakeRDS();self.snapshot_ref=self.s3.add("snapshots/tenant-a/x.json",snapshot(),"s1")
        self.policy_ref=self.s3.add(self.cfg.policy_key,policy(),"p1");self.cfg=handler.Config(**(self.cfg.__dict__|{"policy_sha256":self.policy_ref["sha256"]}))
        self.decision=self.make_decision();self.decision_ref=self.s3.add("decisions/tenant-a/d.json",self.decision,"d1")

    def make_decision(self,passed=True):
        return {"schema":handler.DECISION_SCHEMA,"evaluationId":"a"*64,"tenantId":"tenant-a","configVersion":"cfg-7","schemaProfile":"invoice-v3","snapshotReference":self.snapshot_ref,"snapshotSha256":self.snapshot_ref["sha256"],"policySha256":"b"*64,"reviewEvidenceSha256":"c"*64 if passed else None,"reviewAuthority":"AWS_IDP_HUMAN_REVIEW_COMPLETED" if passed else "SKIPPED_NOT_REVIEWED","sectionCount":1,"classifications":["Invoice"],"evaluationPassed":passed,"decisionReason":"SCHEMA_AND_HUMAN_REVIEW_VERIFIED" if passed else "HUMAN_REVIEW_REQUIRED","automaticBusinessPersistenceAuthorized":False}

    def event(self):
        return {"version":"0","id":"e1","detail-type":"Object Created","source":"aws.s3","account":"111122223333","time":"2026-08-28T10:00:00Z","region":"us-east-1","resources":["arn:aws:s3:::result"],"detail":{"version":"0","event-version":"1.1","bucket":{"name":"result"},"object":{"key":"decisions/tenant-a/d.json","size":1,"etag":"x","version-id":"d1","sequencer":"1"},"request-id":"r","requester":"111122223333","source-ip-address":"1.2.3.4","reason":"PutObject"}}

    def call(self):return handler.process_event_core(self.event(),clients=handler.Clients(self.s3,self.rds),cfg=self.cfg)

    def test_positive_persists_document_and_outbox_in_one_transaction(self):
        out=self.call();self.assertEqual("PERSISTED",out["state"]);self.assertEqual(1,self.rds.commits);self.assertEqual(0,self.rds.rollbacks);self.assertEqual(1,len(self.rds.documents));self.assertEqual(1,len(self.rds.outbox))

    def test_duplicate_reconciles_without_second_rows(self):
        self.call();out=self.call();self.assertEqual("RECONCILED",out["state"]);self.assertEqual(1,len(self.rds.documents));self.assertEqual(1,len(self.rds.outbox));self.assertEqual(2,self.rds.commits)

    def test_divergent_existing_document_rolls_back(self):
        self.call();next(iter(self.rds.documents.values()))[0]="0"*64
        with self.assertRaises(handler.Rejected):self.call()
        self.assertEqual(1,self.rds.rollbacks)

    def test_outbox_failure_rolls_back(self):
        self.rds.fail_outbox=True
        with self.assertRaises(RuntimeError):self.call()
        self.assertEqual(1,self.rds.rollbacks);self.assertEqual(0,self.rds.commits)

    def test_negative_decision_is_ignored_without_database(self):
        self.s3.objects.pop(("decisions/tenant-a/d.json","d1"));self.decision=self.make_decision(False);self.decision_ref=self.s3.add("decisions/tenant-a/d.json",self.decision,"d1")
        out=self.call();self.assertEqual("IGNORED_NEGATIVE",out["state"]);self.assertEqual(0,self.rds.commits)

    def test_event_account_region_bucket_prefix_and_version_are_authoritative(self):
        for path,value in (("account","000000000000"),("region","eu-west-1")):
            event=self.event();event[path]=value
            with self.assertRaises(handler.Rejected):handler.process_event_core(event,clients=handler.Clients(self.s3,self.rds),cfg=self.cfg)
        event=self.event();event["detail"]["object"]["key"]="decisions/tenant-b/d.json"
        with self.assertRaises(handler.Rejected):handler.process_event_core(event,clients=handler.Clients(self.s3,self.rds),cfg=self.cfg)

    def test_event_minor_extension_is_accepted(self):
        event=self.event();event["future-minor-field"]={"safe":True};out=handler.process_event_core(event,clients=handler.Clients(self.s3,self.rds),cfg=self.cfg);self.assertTrue(out["persisted"])

    def test_duplicate_classifications_preserve_section_cardinality(self):
        value=snapshot();value["sections"].append({"section_id":"2","classification":"Invoice","attributes":{"invoiceId":"A-8","total":"20.00"}})
        self.s3.objects.pop(("snapshots/tenant-a/x.json","s1"));self.snapshot_ref=self.s3.add("snapshots/tenant-a/x.json",value,"s1");self.decision=self.make_decision();self.decision["snapshotReference"]=self.snapshot_ref;self.decision["snapshotSha256"]=self.snapshot_ref["sha256"];self.decision["sectionCount"]=2;self.decision["classifications"]=["Invoice","Invoice"]
        self.s3.objects.pop(("decisions/tenant-a/d.json","d1"));self.s3.add("decisions/tenant-a/d.json",self.decision,"d1");out=self.call();self.assertTrue(out["persisted"])

    def test_decision_snapshot_and_policy_storage_authority_are_enforced(self):
        for key in (("decisions/tenant-a/d.json","d1"),("snapshots/tenant-a/x.json","s1"),(self.cfg.policy_key,"p1")):
            old=self.s3.objects[key]["ObjectLockMode"];self.s3.objects[key]["ObjectLockMode"]="GOVERNANCE"
            with self.assertRaises(handler.Rejected):self.call()
            self.s3.objects[key]["ObjectLockMode"]=old

    def test_policy_must_approve_profile_classification_and_automatic_persistence(self):
        for mutation in ({"automaticPersistenceAuthorized":False},{"allowedProfiles":{}},{"businessTenantId":"018f4d4a-7b36-7a21-8d10-2f4c54c28b11"}):
            self.s3.objects.pop((self.cfg.policy_key,"p1"));ref=self.s3.add(self.cfg.policy_key,policy()|mutation,"p1");cfg=handler.Config(**(self.cfg.__dict__|{"policy_sha256":ref["sha256"]}))
            with self.assertRaises(handler.Rejected):handler.process_event_core(self.event(),clients=handler.Clients(self.s3,self.rds),cfg=cfg)

    def test_snapshot_and_decision_must_remain_bound(self):
        self.decision["snapshotSha256"]="0"*64;self.s3.objects.pop(("decisions/tenant-a/d.json","d1"));self.s3.add("decisions/tenant-a/d.json",self.decision,"d1")
        with self.assertRaises(handler.Rejected):self.call()

    def test_lambda_uses_official_partial_batch_processor(self):
        context=object();out=handler.lambda_handler({"Records":[]},context);self.assertTrue(out["called"]);self.assertIs(context,out["context"])


if __name__=="__main__":unittest.main()
````

### FILE: `aws_idp_postgres_persistence_boundary/001_document_intake_outbox.sql`
```yaml
block_id: "AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY:migration:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-postgres-persistence-boundary/001_document_intake_outbox.sql"
license: "LicenseRef-Workspace-Owner"
sha256: "cb46073c7f1ad9fcbd969ab38e30ef48b41e6d83ea359e5966cc4747fffaba51"
variables: []
secrets_allowed: false
```
````sql
begin;

create schema if not exists document_intelligence;

create table if not exists document_intelligence.intake_document (
  tenant_id uuid not null,
  source_tenant_id text not null,
  evaluation_id text not null,
  decision_sha256_hex text not null,
  decision_version_id text not null,
  snapshot_sha256_hex text not null,
  persistence_policy_sha256_hex text not null,
  config_version text not null,
  schema_profile text not null,
  domain_type text not null,
  section_count integer not null,
  document_payload_sha256_hex text not null,
  document_payload jsonb not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, evaluation_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (evaluation_id ~ '^[0-9a-f]{64}$'),
  check (decision_sha256_hex ~ '^[0-9a-f]{64}$'),
  check (snapshot_sha256_hex ~ '^[0-9a-f]{64}$'),
  check (persistence_policy_sha256_hex ~ '^[0-9a-f]{64}$'),
  check (document_payload_sha256_hex ~ '^[0-9a-f]{64}$'),
  check (section_count > 0),
  check (jsonb_typeof(document_payload) = 'object')
);

alter table document_intelligence.intake_document enable row level security;

drop policy if exists intake_document_tenant_isolation on document_intelligence.intake_document;
create policy intake_document_tenant_isolation on document_intelligence.intake_document
  using (tenant_id = nullif(current_setting('app.tenant_id', true), '')::uuid)
  with check (tenant_id = nullif(current_setting('app.tenant_id', true), '')::uuid);

create index if not exists intake_document_created_idx
  on document_intelligence.intake_document (tenant_id, created_at, evaluation_id);

create table if not exists document_intelligence.outbox_event (
  id uuid not null primary key,
  aggregatetype varchar(255) not null,
  aggregateid varchar(255) not null,
  type varchar(255) not null,
  payload jsonb not null,
  tenant_id uuid not null references platform.tenant (tenant_id),
  schema_version integer not null check (schema_version > 0),
  occurred_at timestamptz not null default clock_timestamp(),
  unique (tenant_id, aggregatetype, aggregateid, type),
  check (jsonb_typeof(payload) = 'object')
);

alter table document_intelligence.outbox_event enable row level security;

do $$
begin
  if not exists (select 1 from pg_roles where rolname = 'elite_outbox_cdc') then
    create role elite_outbox_cdc nologin;
  end if;
end
$$;

grant usage on schema document_intelligence to elite_outbox_cdc;
grant select on document_intelligence.outbox_event to elite_outbox_cdc;

drop policy if exists outbox_event_tenant_isolation on document_intelligence.outbox_event;
create policy outbox_event_tenant_isolation on document_intelligence.outbox_event
  using (
    tenant_id = nullif(current_setting('app.tenant_id', true), '')::uuid
    or pg_has_role(current_user, 'elite_outbox_cdc', 'member')
  )
  with check (tenant_id = nullif(current_setting('app.tenant_id', true), '')::uuid);

create index if not exists document_outbox_occurred_idx
  on document_intelligence.outbox_event (occurred_at, id);

commit;
````

### FILE: `aws_idp_postgres_persistence_boundary/persistence-policy.template.json`
```yaml
block_id: "AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY:persistence-policy:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-postgres-persistence-boundary/persistence-policy.template.json"
license: "LicenseRef-Workspace-Owner"
sha256: "f8f19fc52f122729c506b42b3af0cc2934b8f6d1e9ab43b9d35b5bbbd1cbbf3d"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-document-persistence-policy/v1",
  "tenantId": "REPLACE_WITH_SOURCE_TENANT",
  "businessTenantId": "00000000-0000-0000-0000-000000000000",
  "configVersion": "REPLACE_WITH_CONFIG_VERSION",
  "approval": {
    "status": "AWAITING_USER",
    "approvalId": "REQUIRED",
    "approvedBy": "REQUIRED",
    "approvedAt": "REQUIRED"
  },
  "allowedProfiles": {},
  "automaticPersistenceAuthorized": false
}
````

### FILE: `aws_idp_postgres_persistence_boundary/project-profile.template.json`
```yaml
block_id: "AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY:project-profile:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-postgres-persistence-boundary/project-profile.template.json"
license: "LicenseRef-Workspace-Owner"
sha256: "ce7c94063d58b0b470af65971f33f09e99c00e87efdb25cdaa74b1c7f3c5c87e"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-aws-idp-postgres-persistence-profile/v1",
  "acknowledgeConditionedState": false,
  "approvedBy": "REQUIRED",
  "approvedAt": "REQUIRED",
  "sourceTenantId": "REQUIRED",
  "businessTenantId": "00000000-0000-0000-0000-000000000000",
  "idpConfigVersion": "REQUIRED",
  "sourceAccountId": "000000000000",
  "sourceRegion": "REQUIRED",
  "resultBucket": "REQUIRED",
  "resultBucketArn": "REQUIRED",
  "resultKmsKeyArn": "REQUIRED",
  "logKmsKeyArn": "REQUIRED",
  "powertoolsLayerArn": "REQUIRED",
  "persistencePolicyKey": "REQUIRED",
  "persistencePolicyFile": "REQUIRED",
  "persistencePolicyVersionId": "REQUIRED",
  "persistencePolicySha256": "REQUIRED",
  "rdsResourceArn": "REQUIRED",
  "rdsSecretArn": "REQUIRED",
  "rdsDatabase": "REQUIRED",
  "maxObjectBytes": 0,
  "migrationAppliedEvidence": "REQUIRED",
  "rlsIsolationEvidence": "REQUIRED",
  "duplicateRaceEvidence": "REQUIRED",
  "outboxRecoveryEvidence": "REQUIRED",
  "restoreEvidence": "REQUIRED",
  "rollbackEvidence": "REQUIRED",
  "costApproval": "REQUIRED"
}
````

### FILE: `aws_idp_postgres_persistence_boundary/official-contract-lock.json`
```yaml
block_id: "AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY:official-contract-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-postgres-persistence-boundary/official-contract-lock.json"
license: "LicenseRef-Workspace-Owner"
sha256: "238fa1efac74b0f720e2674bd26d26f8f5f5368237db7d108910e9002e1ce3ed"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-official-contract-lock/v1",
  "verifiedAt": "2026-08-28",
  "contracts": [
    {
      "authority": "Amazon S3",
      "url": "https://docs.aws.amazon.com/AmazonS3/latest/userguide/ev-events.html",
      "required": ["Object Created", "event-version major 1", "bucket.name", "object.key", "object.version-id", "reason PutObject"]
    },
    {
      "authority": "AWS Lambda",
      "url": "https://docs.aws.amazon.com/prescriptive-guidance/latest/lambda-event-filtering-partial-batch-responses-for-sqs/best-practices-partial-batch-responses.html",
      "required": ["SQS DLQ", "ReportBatchItemFailures", "idempotent duplicate handling"]
    },
    {
      "authority": "Amazon Aurora RDS Data API",
      "url": "https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/data-api-operations.html",
      "required": ["BeginTransaction", "ExecuteStatement parameters", "CommitTransaction", "RollbackTransaction"]
    },
    {
      "authority": "AWS Prescriptive Guidance",
      "url": "https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/transactional-outbox.html",
      "required": ["business write and outbox in one transaction", "consumer idempotency", "duplicate delivery"]
    },
    {
      "authority": "Dapr Runtime 1.18.3",
      "url": "https://github.com/dapr/dapr/tree/2dcf2e3548faf3c740b3d29ee300e2132db65cff",
      "required": ["official transactional outbox implementation and tests remain separately materializable"]
    }
  ]
}
````

### FILE: `aws_idp_postgres_persistence_boundary/template.yaml`
```yaml
block_id: "AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY:template:v1"
operation: CREATE
provenance: ADAPTED
source: "https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-function.html"
license: "MIT-0 AND LicenseRef-Workspace-Owner"
sha256: "4171df26580af38bdc84daa73797585e6fd80572400853cb8ffd374885df8ac5"
variables: []
secrets_allowed: false
```
````yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Description: Immutable AWS IDP decision to Aurora PostgreSQL document intake and transactional outbox

Parameters:
  SourceTenantId: {Type: String, AllowedPattern: '^[A-Za-z0-9][A-Za-z0-9._:/+=,@-]{0,511}$'}
  BusinessTenantId: {Type: String, AllowedPattern: '^[0-9a-f]{8}-[0-9a-f]{4}-[1-8][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'}
  IdpConfigVersion: {Type: String, AllowedPattern: '^[A-Za-z0-9][A-Za-z0-9._:/+=,@-]{0,511}$'}
  SourceAccountId: {Type: String, AllowedPattern: '^[0-9]{12}$'}
  SourceRegion: {Type: String}
  ResultBucketName: {Type: String}
  ResultBucketArn: {Type: String}
  ResultKmsKeyArn: {Type: String}
  LogKmsKeyArn: {Type: String}
  PowertoolsLayerArn: {Type: String}
  PersistencePolicyKey: {Type: String}
  PersistencePolicyVersionId: {Type: String, NoEcho: true}
  PersistencePolicySha256: {Type: String, AllowedPattern: '^[0-9a-f]{64}$'}
  RdsResourceArn: {Type: String}
  RdsSecretArn: {Type: String}
  RdsDatabase: {Type: String}
  MaxObjectBytes: {Type: Number, MinValue: 1024, MaxValue: 52428800}

Resources:
  PersistenceDeadLetterQueue:
    Type: AWS::SQS::Queue
    Properties:
      QueueName: !Sub '${AWS::StackName}-persistence-dlq'
      KmsMasterKeyId: !Ref ResultKmsKeyArn
      MessageRetentionPeriod: 1209600
      SqsManagedSseEnabled: false

  PersistenceQueue:
    Type: AWS::SQS::Queue
    Properties:
      QueueName: !Sub '${AWS::StackName}-persistence'
      KmsMasterKeyId: !Ref ResultKmsKeyArn
      SqsManagedSseEnabled: false
      VisibilityTimeout: 180
      RedrivePolicy:
        deadLetterTargetArn: !GetAtt PersistenceDeadLetterQueue.Arn
        maxReceiveCount: 8

  PersistenceEventRule:
    Type: AWS::Events::Rule
    Properties:
      EventPattern:
        source: [aws.s3]
        detail-type: [Object Created]
        account: [!Ref SourceAccountId]
        region: [!Ref SourceRegion]
        detail:
          reason: [PutObject]
          bucket:
            name: [!Ref ResultBucketName]
          object:
            key:
              - prefix: !Sub 'decisions/${SourceTenantId}/'
      State: ENABLED
      Targets:
        - Arn: !GetAtt PersistenceQueue.Arn
          Id: PersistenceQueue

  PersistenceQueuePolicy:
    Type: AWS::SQS::QueuePolicy
    Properties:
      Queues: [!Ref PersistenceQueue]
      PolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Effect: Allow
            Principal: {Service: events.amazonaws.com}
            Action: sqs:SendMessage
            Resource: !GetAtt PersistenceQueue.Arn
            Condition:
              ArnEquals:
                aws:SourceArn: !GetAtt PersistenceEventRule.Arn

  PersistenceFunctionRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Effect: Allow
            Principal: {Service: lambda.amazonaws.com}
            Action: sts:AssumeRole
      Policies:
        - PolicyName: exact-runtime-authority
          PolicyDocument:
            Version: '2012-10-17'
            Statement:
              - Effect: Allow
                Action: [sqs:ReceiveMessage, sqs:DeleteMessage, sqs:GetQueueAttributes]
                Resource: !GetAtt PersistenceQueue.Arn
              - Effect: Allow
                Action: [s3:GetObjectVersion, s3:GetObjectVersionAttributes]
                Resource:
                  - !Sub '${ResultBucketArn}/decisions/${SourceTenantId}/*'
                  - !Sub '${ResultBucketArn}/snapshots/${SourceTenantId}/*'
                  - !Sub '${ResultBucketArn}/persistence-policies/${SourceTenantId}/*'
              - Effect: Allow
                Action: [kms:Decrypt, kms:GenerateDataKey]
                Resource: [!Ref ResultKmsKeyArn, !Ref LogKmsKeyArn]
              - Effect: Allow
                Action: [rds-data:BeginTransaction, rds-data:ExecuteStatement, rds-data:CommitTransaction, rds-data:RollbackTransaction]
                Resource: !Ref RdsResourceArn
              - Effect: Allow
                Action: secretsmanager:GetSecretValue
                Resource: !Ref RdsSecretArn
              - Effect: Allow
                Action: [logs:CreateLogStream, logs:PutLogEvents]
                Resource: !Sub 'arn:${AWS::Partition}:logs:${AWS::Region}:${AWS::AccountId}:log-group:/aws/lambda/${AWS::StackName}-persist:*'

  PersistenceFunction:
    Type: AWS::Serverless::Function
    Properties:
      FunctionName: !Sub '${AWS::StackName}-persist'
      Runtime: python3.13
      Handler: handler.lambda_handler
      CodeUri: .
      Role: !GetAtt PersistenceFunctionRole.Arn
      Layers: [!Ref PowertoolsLayerArn]
      Architectures: [arm64]
      MemorySize: 512
      Timeout: 60
      ReservedConcurrentExecutions: 20
      Environment:
        Variables:
          SOURCE_TENANT_ID: !Ref SourceTenantId
          BUSINESS_TENANT_ID: !Ref BusinessTenantId
          IDP_CONFIG_VERSION: !Ref IdpConfigVersion
          SOURCE_ACCOUNT_ID: !Ref SourceAccountId
          SOURCE_REGION: !Ref SourceRegion
          RESULT_BUCKET: !Ref ResultBucketName
          RESULT_KMS_KEY_ARN: !Ref ResultKmsKeyArn
          PERSISTENCE_POLICY_KEY: !Ref PersistencePolicyKey
          PERSISTENCE_POLICY_VERSION_ID: !Ref PersistencePolicyVersionId
          PERSISTENCE_POLICY_SHA256: !Ref PersistencePolicySha256
          RDS_RESOURCE_ARN: !Ref RdsResourceArn
          RDS_SECRET_ARN: !Ref RdsSecretArn
          RDS_DATABASE: !Ref RdsDatabase
          MAX_OBJECT_BYTES: !Ref MaxObjectBytes
          POWERTOOLS_SERVICE_NAME: idp-persistence-boundary
          POWERTOOLS_LOG_LEVEL: INFO
      LoggingConfig:
        LogFormat: JSON
        ApplicationLogLevel: INFO
        SystemLogLevel: WARN

  PersistenceLogGroup:
    Type: AWS::Logs::LogGroup
    Properties:
      LogGroupName: !Sub '/aws/lambda/${AWS::StackName}-persist'
      KmsKeyId: !Ref LogKmsKeyArn
      RetentionInDays: 365

  PersistenceEventSource:
    Type: AWS::Lambda::EventSourceMapping
    Properties:
      EventSourceArn: !GetAtt PersistenceQueue.Arn
      FunctionName: !Ref PersistenceFunction
      BatchSize: 10
      MaximumBatchingWindowInSeconds: 5
      FunctionResponseTypes: [ReportBatchItemFailures]
      Enabled: true

Outputs:
  PersistenceQueueArn: {Value: !GetAtt PersistenceQueue.Arn}
  PersistenceDeadLetterQueueArn: {Value: !GetAtt PersistenceDeadLetterQueue.Arn}
  PersistenceFunctionArn: {Value: !GetAtt PersistenceFunction.Arn}
````

### FILE: `aws_idp_postgres_persistence_boundary/apply_migration.ps1`
```yaml
block_id: "AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY:apply-migration:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-postgres-persistence-boundary/apply_migration.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "4ca5f56eb0d0d41aa840dc1e48ad91af069c80987ddc65ebec5d746acebd027e"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding(SupportsShouldProcess)]param(
  [Parameter(Mandatory)][string]$PgServiceName,
  [Parameter(Mandatory)][string]$FoundationSql,
  [string]$ComponentRoot=$PSScriptRoot
)
$ErrorActionPreference='Stop'
if($PgServiceName -notmatch '^[A-Za-z0-9_.-]{1,128}$'){throw 'PG service name is invalid.'}
$foundation=(Resolve-Path -LiteralPath $FoundationSql).Path
$migration=(Resolve-Path -LiteralPath (Join-Path $ComponentRoot '001_document_intake_outbox.sql')).Path
if(-not(Get-Command psql -ErrorAction SilentlyContinue)){throw 'psql is required.'}
if(-not$PSCmdlet.ShouldProcess("PostgreSQL service $PgServiceName",'apply foundation and document intake migrations')){return}
& psql "service=$PgServiceName" -X -v ON_ERROR_STOP=1 -f $foundation
if($LASTEXITCODE-ne0){throw 'PostgreSQL foundation migration failed.'}
& psql "service=$PgServiceName" -X -v ON_ERROR_STOP=1 -f $migration
if($LASTEXITCODE-ne0){throw 'Document intake migration failed.'}
'AWS_IDP_POSTGRES_PERSISTENCE_MIGRATION_PASS'
````

### FILE: `aws_idp_postgres_persistence_boundary/deploy.ps1`
```yaml
block_id: "AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY:deploy:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-postgres-persistence-boundary/deploy.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "881d32595ad2bb22ac9959f540ac8b79cf51ce855f9b47a355d16fbc52319a10"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding(SupportsShouldProcess)]param(
  [Parameter(Mandatory)][string]$Profile,
  [Parameter(Mandatory)][string]$StackName,
  [string]$Region,
  [string]$ComponentRoot=$PSScriptRoot
)
$ErrorActionPreference='Stop'
$profilePath=(Resolve-Path -LiteralPath $Profile).Path
$p=Get-Content -LiteralPath $profilePath -Raw|ConvertFrom-Json -Depth 64
if($p.schema-ne'elite-aws-idp-postgres-persistence-profile/v1'-or$p.acknowledgeConditionedState-ne$true){throw 'Profile is not approved.'}
$required=@('approvedBy','approvedAt','sourceTenantId','businessTenantId','idpConfigVersion','sourceRegion','resultBucket','resultBucketArn','resultKmsKeyArn','logKmsKeyArn','powertoolsLayerArn','persistencePolicyKey','persistencePolicyFile','persistencePolicyVersionId','persistencePolicySha256','rdsResourceArn','rdsSecretArn','rdsDatabase','migrationAppliedEvidence','rlsIsolationEvidence','duplicateRaceEvidence','outboxRecoveryEvidence','restoreEvidence','rollbackEvidence','costApproval')
foreach($name in $required){$value=[string]$p.$name;if([string]::IsNullOrWhiteSpace($value)-or$value-eq'REQUIRED'-or$value-eq'REPLACE'){throw "Unresolved profile field: $name"}}
if([string]$p.sourceAccountId-notmatch'^[0-9]{12}$'-or[string]$p.businessTenantId-notmatch'^[0-9a-f]{8}-[0-9a-f]{4}-[1-8][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'){throw 'Profile identity is invalid.'}
if([int64]$p.maxObjectBytes-lt1024-or[int64]$p.maxObjectBytes-gt52428800){throw 'maxObjectBytes is invalid.'}
$policyPath=(Resolve-Path -LiteralPath ([string]$p.persistencePolicyFile)).Path
$policy=Get-Content -LiteralPath $policyPath -Raw|ConvertFrom-Json -Depth 64
if($policy.schema-ne'elite-document-persistence-policy/v1'-or$policy.approval.status-ne'APPROVED'-or$policy.automaticPersistenceAuthorized-ne$true){throw 'Persistence policy is not approved.'}
$policyHash=(Get-FileHash -LiteralPath $policyPath -Algorithm SHA256).Hash.ToLowerInvariant()
if($policyHash-ne[string]$p.persistencePolicySha256){throw 'Persistence policy SHA-256 mismatch.'}
foreach($tool in @('sam','aws')){if(-not(Get-Command $tool -ErrorAction SilentlyContinue)){throw "$tool is required."}}
if(-not$Region){$Region=[string]$p.sourceRegion}
& sam validate --lint --template-file (Join-Path $ComponentRoot 'template.yaml') --region $Region
if($LASTEXITCODE-ne0){throw 'sam validate failed.'}
& sam build --template-file (Join-Path $ComponentRoot 'template.yaml') --build-dir (Join-Path $ComponentRoot '.aws-sam/build') --cached:$false
if($LASTEXITCODE-ne0){throw 'sam build failed.'}
if(-not$PSCmdlet.ShouldProcess($StackName,'deploy AWS IDP PostgreSQL persistence boundary')){return}
$args=@("SourceTenantId=$($p.sourceTenantId)","BusinessTenantId=$($p.businessTenantId)","IdpConfigVersion=$($p.idpConfigVersion)","SourceAccountId=$($p.sourceAccountId)","SourceRegion=$($p.sourceRegion)","ResultBucketName=$($p.resultBucket)","ResultBucketArn=$($p.resultBucketArn)","ResultKmsKeyArn=$($p.resultKmsKeyArn)","LogKmsKeyArn=$($p.logKmsKeyArn)","PowertoolsLayerArn=$($p.powertoolsLayerArn)","PersistencePolicyKey=$($p.persistencePolicyKey)","PersistencePolicyVersionId=$($p.persistencePolicyVersionId)","PersistencePolicySha256=$($p.persistencePolicySha256)","RdsResourceArn=$($p.rdsResourceArn)","RdsSecretArn=$($p.rdsSecretArn)","RdsDatabase=$($p.rdsDatabase)","MaxObjectBytes=$($p.maxObjectBytes)")
& sam deploy --stack-name $StackName --region $Region --capabilities CAPABILITY_IAM --resolve-s3 --confirm-changeset --parameter-overrides $args
if($LASTEXITCODE-ne0){throw 'sam deploy failed.'}
````

### FILE: `aws_idp_postgres_persistence_boundary/verify_contract.ps1`
```yaml
block_id: "AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY:verify-contract:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://aws-idp-postgres-persistence-boundary/verify_contract.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "d8e396b365dd78f4e3c9fea246328336a7189e9eb2dd6e9203e079af2cb9e446"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]param([string]$ComponentRoot=$PSScriptRoot)
$ErrorActionPreference='Stop'
$files=@('README.md','handler.py','test_handler.py','001_document_intake_outbox.sql','persistence-policy.template.json','project-profile.template.json','official-contract-lock.json','template.yaml','apply_migration.ps1','deploy.ps1','verify_contract.ps1')
foreach($file in $files){if(-not(Test-Path -LiteralPath (Join-Path $ComponentRoot $file)-PathType Leaf)){throw "Missing $file"}}
$profile=Get-Content -LiteralPath (Join-Path $ComponentRoot 'project-profile.template.json')-Raw|ConvertFrom-Json -Depth 64
if($profile.acknowledgeConditionedState-ne$false-or$profile.maxObjectBytes-ne0){throw 'Project profile must begin fail closed.'}
$policy=Get-Content -LiteralPath (Join-Path $ComponentRoot 'persistence-policy.template.json')-Raw|ConvertFrom-Json -Depth 64
if($policy.approval.status-ne'AWAITING_USER'-or$policy.automaticPersistenceAuthorized-ne$false-or@($policy.allowedProfiles.PSObject.Properties).Count-ne0){throw 'Persistence policy must begin fail closed.'}
$handler=Get-Content -LiteralPath (Join-Path $ComponentRoot 'handler.py')-Raw
foreach($marker in @('BatchProcessor','process_partial_response','ObjectLockMode','version-id','begin_transaction','commit_transaction','rollback_transaction','elite:insert-document','elite:insert-outbox','automaticBusinessPersistenceAuthorized')){if(-not$handler.Contains($marker)){throw "Handler contract missing $marker"}}
if($handler-match'print\('-or$handler-match'logger\.'){throw 'Handler may not log document values or identifiers.'}
$sql=Get-Content -LiteralPath (Join-Path $ComponentRoot '001_document_intake_outbox.sql')-Raw
foreach($marker in @('foreign key (tenant_id) references platform.tenant','enable row level security','current_setting(''app.tenant_id''','document_payload_sha256_hex','aggregatetype varchar(255)','aggregateid varchar(255)','document_intelligence.outbox_event','create role elite_outbox_cdc nologin','pg_has_role(current_user, ''elite_outbox_cdc'', ''member'')','grant select on document_intelligence.outbox_event')){if(-not$sql.Contains($marker)){throw "Migration contract missing $marker"}}
Push-Location $ComponentRoot
try{
  python -B -c "from pathlib import Path;[compile(Path(p).read_text(encoding='utf-8'),p,'exec') for p in ('handler.py','test_handler.py')]"
  if($LASTEXITCODE-ne0){throw 'compile failed'}
  python -B -m unittest -v test_handler.py
  if($LASTEXITCODE-ne0){throw 'tests failed'}
}finally{Pop-Location}
'AWS_IDP_POSTGRES_PERSISTENCE_BOUNDARY_PASS'
````

## 6. Configuration surface

El project profile y la policy empiezan cerrados. El usuario debe fijar tenant fuente y UUID de negocio, config, cuenta/región, bucket/KMS, policy file/key/VersionId/SHA, Powertools layer, Aurora/secret/database, límite de bytes y evidencias de migración/RLS/duplicates/outbox/restore/rollback/costo. Ningún secret se guarda en Markdown; RDS usa Secrets Manager por ARN.

## 7. Dependency bill

- AWS Powertools Lambda Python 3.34.0/commit `376757…`, patrón batch MIT-0; handler local `ADAPTED`.
- Boto3 RDS Data API/S3 del runtime AWS; Lambda Python 3.13; S3 EventBridge; SQS/DLQ; KMS; Secrets Manager; Aurora PostgreSQL.
- `PG-TX-FOUNDATION` 0.1.x aporta tenant/inbox; Debezium Event Router 3.6.1.Final y AWS outbox guidance aportan contratos publicados, no bytes copiados aquí.
- Nueve archivos `AUTHORED`, dos `ADAPTED`, cero `VERBATIM`.

## 8. Apply order

1. Materializar foundation PostgreSQL y este pack; ejecutar verifiers/cfn-lint.
2. Aplicar ambas migraciones con `apply_migration.ps1` mediante PGSERVICE sin secrets en argumentos.
3. Crear role/secret runtime no-owner y demostrar RLS cross-tenant.
4. Completar/aprobar policy, subirla create-only/Object-Locked y fijar key/VersionId/SHA.
5. Habilitar S3 EventBridge, completar profile y desplegar DEV con change set confirmado.
6. Probar positive/negative, tamper, duplicates/races, partial failures, DLQ/redrive, outbox outage/recovery, restore/rollback/costo.
7. Conectar mappings de dominio mediante consumer inbox idempotente; no escribir ERP directamente desde este staging.

## 9. Verification

El verifier exige 11 archivos, templates fail-closed, marcadores oficiales S3/SQS/Data API/outbox Debezium, RLS writer/CDC, compile y 12 pruebas. cfn-lint y round-trip son obligatorios. Live AWS/Aurora, policy real, migration/RLS/race/load/DLQ/CDC/restore/rollback/costo siguen condicionados.

## 10. Reconstruction evidence

V96 registra boundary 0.2 Debezium-shaped y runtime CDC; V95 conserva el primer intake+outbox transaccional. No se afirma AWS/Aurora/Kafka live ni exactly-once.
