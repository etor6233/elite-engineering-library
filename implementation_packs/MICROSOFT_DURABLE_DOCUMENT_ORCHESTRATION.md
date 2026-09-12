# Microsoft Durable Document Orchestration

## 1. Metadata

```yaml
pack_id: "MICROSOFT-DURABLE-DOCUMENT-ORCHESTRATION"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa una orquestación documental determinista sobre Microsoft Durable Task Python 1.9.0 con seguridad, evidencia original, extracción, evaluación estricta, revisión humana durable, persistencia idempotente y evidencia final; liga los dos samples oficiales Human Interaction por path/SHA y bloquea su frontera anónima como autoridad de aprobación."
stacks: ["CPython 3.12.13", "Microsoft durabletask 1.9.0", "Durable Task Scheduler or official in-memory test backend"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x", "STRICT-DOCUMENT-FIELD-EVALUATION-GATE 0.1.x", "GO-AWS-ENTERPRISE-STORAGE-EMAIL-ADAPTERS 0.2.x", "GO-GOOGLE-CLOUD-STORAGE-ADAPTER 0.1.x", "PYTHON-AZURE-BLOB-IMMUTABLE-EVIDENCE-ADAPTER 0.1.x"]
incompatible_with: ["datos o bytes del documento en orchestration history", "persistencia automática sin evaluación o review", "retry ciego de efectos", "backend in-memory en producción", "imagen latest sin digest", "credenciales embebidas"]
license_expression: "LicenseRef-Workspace-Owner AND MIT"
upstream_sources: ["https://pypi.org/project/durabletask/1.9.0/", "https://github.com/microsoft/durabletask-python/tree/6dfdbac521d9d59d31d7ea975fb669d66f86f4c4", "https://github.com/Azure-Samples/Durable-Task-Scheduler/tree/5636c25ffbdaabaca1062c3b6c6e77072f2b2eb2", "https://learn.microsoft.com/azure/durable-task/sdks/quickstart-portable-durable-task-sdks"]
verified_at: "2026-09-04"
```

## 2. Applicability

Use cuando el proyecto necesite coordinar documentos desde cuarentena hasta persistencia de negocio sin perder estado entre reinicios. El usuario debe elegir backend, identidad, task hub, lane de extracción, storage inmutable, clases/schemas, política de revisión y autoridad de escritura. El template distribuido habilita cero efectos.

El backend in-memory oficial Microsoft sirve sólo para pruebas gratuitas y no es persistencia productiva ni WORM. Docker no estaba disponible en esta auditoría; el emulador DTS no se ejecutó. Producción exige endpoint aprobado, TLS, RBAC/managed identity, aislamiento, observabilidad, carga, recovery, versionado y evidencia live.

## 3. Architecture contract

La orquestación es determinista y sólo conserva referencias opacas, hashes y receipts. Los bytes, campos extraídos, identidad del revisor y secretos permanecen dentro de actividades/stores autorizados. Las fronteras fijas son security, retain-original, extract, evaluate, validate-review, persist y retain-final.

No hay retry automático de efectos. Un fallo ambiguo al retener el original produce reconciliación de storage; un fallo de persistencia produce UNKNOWN_BUSINESS_EFFECT; un fallo de evidencia final después del commit produce PARTIAL_EFFECT y conserva la transacción. Sólo AUTO_ACCEPT demostrado o review APPROVED hash-bound permite persistir. La clave idempotente liga tenant, hash documental, clase y schema.

El skill oficial Microsoft se conserva byte a byte como referencia. Sus ejemplos sin pin y la imagen con tag latest no son autoridad de ejecución; el lock local y el admission del proyecto prevalecen. LICENSE.md sólo añade newline final y se declara ADAPTED, nunca VERBATIM.

## 4. Exact file manifest

```text
CREATE durable_document_orchestration/requirements-direct.in
CREATE durable_document_orchestration/requirements-windows-py312.lock
CREATE durable_document_orchestration/official-artifact-lock.json
CREATE durable_document_orchestration/pipeline-profile.template.json
CREATE durable_document_orchestration/documentflow/__init__.py
CREATE durable_document_orchestration/documentflow/contracts.py
CREATE durable_document_orchestration/documentflow/orchestration.py
CREATE durable_document_orchestration/documentflow/runtime.py
CREATE durable_document_orchestration/test_contracts.py
CREATE durable_document_orchestration/test_pipeline.py
CREATE durable_document_orchestration/README.md
CREATE durable_document_orchestration/upstream/LICENSE.md
CREATE durable_document_orchestration/upstream/durable-task-python-SKILL.md
```

## 5. Materialization blocks

### FILE: `durable_document_orchestration/requirements-direct.in`
```yaml
block_id: "MICROSOFT-DURABLE-DOCUMENT:requirements-direct-in:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration around exact official Microsoft Durable Task APIs"
license: "LicenseRef-Workspace-Owner"
sha256: "43faf914614ecd498cd0613ed87f5e6eb309ccd38db223c7872472d7044e7aae"
variables: []
secrets_allowed: false
```
````text
durabletask==1.9.0
````

### FILE: `durable_document_orchestration/requirements-windows-py312.lock`
```yaml
block_id: "MICROSOFT-DURABLE-DOCUMENT:requirements-windows-py312-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration around exact official Microsoft Durable Task APIs"
license: "LicenseRef-Workspace-Owner"
sha256: "4d7cca8b66e29db5226811ef09e7532fdac1e0afbe2ab1e8fe247bcbc4e4ed7f"
variables: []
secrets_allowed: false
```
````text
asyncio==4.0.0 \
    --hash=sha256:c1eddb0659231837046809e68103969b2bef8b0400d59cfa6363f6b5ed8cc88b
durabletask==1.9.0 \
    --hash=sha256:a759af4ad8e6897922575886e93bf40c6b3af936eaed83798d492ee561607185
grpcio==1.83.0 \
    --hash=sha256:1c699bbb20f143c8f2bff219de578aa2dc1f919399d67dc702b038b986ee62df
packaging==26.3 \
    --hash=sha256:d7193f7c8e4e93f444fde0262bf90af30e16fa0ad0ad44cb553c87339b23cd1c
protobuf==7.36.0 \
    --hash=sha256:1781cc1de61249b750848029bca452c0a8b7e990080316b9bbc2518b2117b488
typing-extensions==4.16.0 \
    --hash=sha256:481caa481374e813c1b176ada14e97f1f67a4539ce9cfeb3f350d78d6370c2e8
````

### FILE: `durable_document_orchestration/official-artifact-lock.json`
```yaml
block_id: "MICROSOFT-DURABLE-DOCUMENT:official-artifact-lock-json:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration around exact official Microsoft Durable Task APIs"
license: "LicenseRef-Workspace-Owner"
sha256: "b56457e301a656038e1dab5276ce6b12310121fcb216de686b8b43f4f01217ab"
variables: []
secrets_allowed: false
```
````text
{
  "schema_version": 1,
  "runtime_target": "CPython 3.12 / Windows x86-64",
  "sources": [
    {
      "id": "microsoft-durabletask-python-1.9.0",
      "repository": "https://github.com/microsoft/durabletask-python",
      "tag": "v1.9.0",
      "commit": "6dfdbac521d9d59d31d7ea975fb669d66f86f4c4",
      "commit_signature": "verified",
      "archive_bytes": 1418256,
      "archive_sha256": "c8c941b25804ae9abcc0470f5594cab4929216a981982616344eb2b84a7916e8",
      "wheel_bytes": 212413,
      "wheel_sha256": "a759af4ad8e6897922575886e93bf40c6b3af936eaed83798d492ee561607185",
      "sdist_bytes": 184873,
      "sdist_sha256": "55460fbfda8941e721096c4639e72111b03638708df1556af466160ee649478d",
      "human_interaction_sample_path": "examples/human_interaction.py",
      "human_interaction_sample_sha256": "7f97502962fac35775057651a2081f5e4b7e19bc116dba7145b70c8afb467669",
      "azure_functions_human_interaction_sample_path": "azure-functions-durable/samples/human-interaction/function_app.py",
      "azure_functions_human_interaction_sample_sha256": "42fdd853d67fcec3acc51ee12e6df1f4ac6c221d246f77a24585d023ec3f0e83",
      "license": "MIT",
      "license_sha256": "c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383"
    },
    {
      "id": "azure-samples-durable-task-scheduler-5636c25",
      "repository": "https://github.com/Azure-Samples/Durable-Task-Scheduler",
      "commit": "5636c25ffbdaabaca1062c3b6c6e77072f2b2eb2",
      "commit_signature": "verified",
      "archive_bytes": 8425148,
      "archive_sha256": "f9754c7d0f9d957aa01ec541d78f838fe097b0cdb0fcb57d7decb8b0c0e55347",
      "license": "MIT",
      "license_sha256": "d9a1b1e30d633d5732ea18e3cba9538d293ebc53e1a9e4e96ab739e0c5c4f1cb",
      "packaged_license_sha256": "c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383",
      "packaged_license_provenance": "ADAPTED_FINAL_NEWLINE_ONLY",
      "skill_path": ".github/skills/durable-task-python/SKILL.md",
      "skill_sha256": "553f0dc278c63bb129df5557a7d6f7d645bf827efba4cb077588688f88a2e9e2"
    }
  ],
  "runtime_packages": 6,
  "osv_findings": 0,
  "human_interaction_condition": "Official samples prove durable external-event plus timeout mechanics only. They do not authenticate an approver; the Azure Functions HTTP sample is ANONYMOUS. Production review must validate actor, tenant, resource authority, replay/idempotency and hash-bound decision before raising the event.",
  "verified_at": "2026-08-27"
}
````

### FILE: `durable_document_orchestration/pipeline-profile.template.json`
```yaml
block_id: "MICROSOFT-DURABLE-DOCUMENT:pipeline-profile-template-json:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration around exact official Microsoft Durable Task APIs"
license: "LicenseRef-Workspace-Owner"
sha256: "6df767c2e34ab66ecad033f41a9ba4b3ae2d184a6888e859fea7f4b0262e5426"
variables: []
secrets_allowed: false
```
````text
{
  "schema_version": 1,
  "enabled": false,
  "backend": "AWAITING_SELECTION",
  "host_address": "",
  "task_hub": "",
  "secure_channel": true,
  "managed_identity_approved": false,
  "emulator_digest": "",
  "document_classes": [],
  "storage_lane": "AWAITING_SELECTION",
  "review_timeout_seconds": 0,
  "max_document_bytes": 0,
  "automatic_business_persistence": false,
  "live_effects_approved": false,
  "owner": "",
  "notes": "The distributed template intentionally authorizes no backend, storage, provider, review timeout or business write."
}
````

### FILE: `durable_document_orchestration/documentflow/__init__.py`
```yaml
block_id: "MICROSOFT-DURABLE-DOCUMENT:documentflow-init-py:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration around exact official Microsoft Durable Task APIs"
license: "LicenseRef-Workspace-Owner"
sha256: "3157eb66954f95b5bcbb30adcab933dfa6be943eb17286ee15723c5844bdb317"
variables: []
secrets_allowed: false
```
````text
"""Fail-closed document orchestration on Microsoft Durable Task."""

from .contracts import expected_idempotency_key, validate_request
from .orchestration import document_pipeline_orchestrator
from .runtime import ActivityHandlers, ManagedTaskHubGrpcWorker, register_document_pipeline

__all__ = [
    "ActivityHandlers",
    "ManagedTaskHubGrpcWorker",
    "document_pipeline_orchestrator",
    "expected_idempotency_key",
    "register_document_pipeline",
    "validate_request",
]
````

### FILE: `durable_document_orchestration/documentflow/contracts.py`
```yaml
block_id: "MICROSOFT-DURABLE-DOCUMENT:documentflow-contracts-py:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration around exact official Microsoft Durable Task APIs"
license: "LicenseRef-Workspace-Owner"
sha256: "f591190226211095c9ebef22652737491232ceefa5d8c1b3be0ab8fb1559d6dd"
variables: []
secrets_allowed: false
```
````text
"""Closed, history-safe contracts for the durable document pipeline."""

from __future__ import annotations

import hashlib
import re
from collections.abc import Iterable, Mapping
from typing import Any

_SHA256 = re.compile(r"^[0-9a-f]{64}$")
_TOKEN = re.compile(r"^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$")
_REF = re.compile(r"^(urn|ref):[A-Za-z0-9][A-Za-z0-9._:/-]{0,255}$")

REQUEST_KEYS = frozenset(
    {
        "contract_version",
        "tenant_ref",
        "document_id",
        "document_sha256",
        "class_id",
        "schema_version",
        "profile_sha256",
        "routing_receipt_sha256",
        "input_ref",
        "idempotency_key",
        "review_timeout_seconds",
        "evidence_mode",
    }
)


def _mapping(value: Any, label: str) -> dict[str, Any]:
    if not isinstance(value, Mapping):
        raise ValueError(f"{label} must be an object")
    return dict(value)


def _exact(value: Mapping[str, Any], keys: Iterable[str], label: str) -> None:
    expected = frozenset(keys)
    actual = frozenset(value)
    if actual != expected:
        missing = sorted(expected - actual)
        extra = sorted(actual - expected)
        raise ValueError(f"{label} keys mismatch missing={missing} extra={extra}")


def _sha(value: Any, label: str) -> str:
    if not isinstance(value, str) or not _SHA256.fullmatch(value):
        raise ValueError(f"{label} must be lowercase SHA-256")
    return value


def _token(value: Any, label: str) -> str:
    if not isinstance(value, str) or not _TOKEN.fullmatch(value):
        raise ValueError(f"{label} must be an opaque safe token")
    return value


def _ref(value: Any, label: str, *, allow_empty: bool = False) -> str:
    if allow_empty and value == "":
        return ""
    if not isinstance(value, str) or not _REF.fullmatch(value):
        raise ValueError(f"{label} must be an opaque urn:/ref: reference")
    return value


def expected_idempotency_key(
    tenant_ref: str, document_sha256: str, class_id: str, schema_version: str
) -> str:
    material = "\n".join((tenant_ref, document_sha256, class_id, schema_version)).encode("utf-8")
    return hashlib.sha256(material).hexdigest()


def validate_request(value: Any) -> dict[str, Any]:
    request = _mapping(value, "request")
    _exact(request, REQUEST_KEYS, "request")
    if request["contract_version"] != 1:
        raise ValueError("contract_version must be 1")
    for field in ("tenant_ref", "document_id", "class_id", "schema_version"):
        request[field] = _token(request[field], field)
    for field in ("document_sha256", "profile_sha256", "routing_receipt_sha256"):
        request[field] = _sha(request[field], field)
    request["input_ref"] = _ref(request["input_ref"], "input_ref")
    if request["evidence_mode"] not in {
        "AWS_S3_OBJECT_LOCK",
        "GCS_OBJECT_RETENTION",
        "AZURE_BLOB_WORM",
        "TEST_ONLY_MEMORY",
    }:
        raise ValueError("evidence_mode is unsupported")
    timeout = request["review_timeout_seconds"]
    if not isinstance(timeout, int) or isinstance(timeout, bool) or not 1 <= timeout <= 2_592_000:
        raise ValueError("review_timeout_seconds must be 1..2592000")
    expected = expected_idempotency_key(
        request["tenant_ref"],
        request["document_sha256"],
        request["class_id"],
        request["schema_version"],
    )
    if request["idempotency_key"] != expected:
        raise ValueError("idempotency_key is not bound to tenant/document/class/schema")
    return request


def _receipt(
    value: Any,
    *,
    label: str,
    stage: str,
    document_sha256: str,
    statuses: set[str],
    extra_keys: Iterable[str] = (),
) -> dict[str, Any]:
    receipt = _mapping(value, label)
    keys = {
        "stage",
        "status",
        "document_sha256",
        "receipt_ref",
        "receipt_sha256",
        *extra_keys,
    }
    _exact(receipt, keys, label)
    if receipt["stage"] != stage or receipt["status"] not in statuses:
        raise ValueError(f"{label} stage/status is invalid")
    if receipt["document_sha256"] != document_sha256:
        raise ValueError(f"{label} document hash mismatch")
    _ref(receipt["receipt_ref"], f"{label}.receipt_ref")
    _sha(receipt["receipt_sha256"], f"{label}.receipt_sha256")
    return receipt


def validate_security(value: Any, document_sha256: str) -> dict[str, Any]:
    return _receipt(
        value,
        label="security receipt",
        stage="security",
        document_sha256=document_sha256,
        statuses={"ADMITTED", "REJECTED"},
    )


def validate_retention(value: Any, document_sha256: str, stage: str) -> dict[str, Any]:
    receipt = _receipt(
        value,
        label=f"{stage} receipt",
        stage=stage,
        document_sha256=document_sha256,
        statuses={"VERIFIED", "UNKNOWN_EFFECT", "PARTIAL_VERIFICATION"},
        extra_keys={"object_ref", "version_ref"},
    )
    _ref(receipt["object_ref"], f"{stage}.object_ref", allow_empty=receipt["status"] != "VERIFIED")
    _ref(receipt["version_ref"], f"{stage}.version_ref", allow_empty=receipt["status"] != "VERIFIED")
    return receipt


def validate_extraction(value: Any, request: Mapping[str, Any]) -> dict[str, Any]:
    receipt = _receipt(
        value,
        label="extraction receipt",
        stage="extraction",
        document_sha256=request["document_sha256"],
        statuses={"EXTRACTED", "AMBIGUOUS", "REJECTED"},
        extra_keys={"result_ref", "result_sha256", "schema_version"},
    )
    _ref(receipt["result_ref"], "extraction.result_ref")
    _sha(receipt["result_sha256"], "extraction.result_sha256")
    if receipt["schema_version"] != request["schema_version"]:
        raise ValueError("extraction schema version mismatch")
    return receipt


def validate_evaluation(value: Any, document_sha256: str, result_sha256: str) -> dict[str, Any]:
    receipt = _receipt(
        value,
        label="evaluation receipt",
        stage="evaluation",
        document_sha256=document_sha256,
        statuses={"AUTO_ACCEPT", "REVIEW_REQUIRED", "REJECTED"},
        extra_keys={"result_sha256"},
    )
    if _sha(receipt["result_sha256"], "evaluation.result_sha256") != result_sha256:
        raise ValueError("evaluation result hash mismatch")
    return receipt


def validate_review_event(value: Any, result_sha256: str) -> dict[str, Any]:
    event = _mapping(value, "review event")
    _exact(
        event,
        {"decision", "subject_result_sha256", "review_receipt_ref", "review_receipt_sha256"},
        "review event",
    )
    if event["decision"] not in {"APPROVED", "REJECTED"}:
        raise ValueError("review decision is invalid")
    if _sha(event["subject_result_sha256"], "review.subject_result_sha256") != result_sha256:
        raise ValueError("review result hash mismatch")
    _ref(event["review_receipt_ref"], "review.review_receipt_ref")
    _sha(event["review_receipt_sha256"], "review.review_receipt_sha256")
    return event


def validate_review(value: Any, document_sha256: str, result_sha256: str) -> dict[str, Any]:
    receipt = _receipt(
        value,
        label="review receipt",
        stage="review",
        document_sha256=document_sha256,
        statuses={"APPROVED", "REJECTED"},
        extra_keys={"result_sha256"},
    )
    if _sha(receipt["result_sha256"], "review.result_sha256") != result_sha256:
        raise ValueError("validated review result hash mismatch")
    return receipt


def validate_persistence(value: Any, request: Mapping[str, Any], result_sha256: str) -> dict[str, Any]:
    receipt = _receipt(
        value,
        label="persistence receipt",
        stage="business_persistence",
        document_sha256=request["document_sha256"],
        statuses={"COMMITTED", "ALREADY_COMMITTED"},
        extra_keys={"result_sha256", "idempotency_key", "transaction_ref"},
    )
    if _sha(receipt["result_sha256"], "persistence.result_sha256") != result_sha256:
        raise ValueError("persistence result hash mismatch")
    if receipt["idempotency_key"] != request["idempotency_key"]:
        raise ValueError("persistence idempotency key mismatch")
    _ref(receipt["transaction_ref"], "persistence.transaction_ref")
    return receipt
````

### FILE: `durable_document_orchestration/documentflow/orchestration.py`
```yaml
block_id: "MICROSOFT-DURABLE-DOCUMENT:documentflow-orchestration-py:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration around exact official Microsoft Durable Task APIs"
license: "LicenseRef-Workspace-Owner"
sha256: "af72a431a24bbd6403d94681d4a075542987dd5b7ef2318f180d181f7c68caee"
variables: []
secrets_allowed: false
```
````text
"""Deterministic document workflow implemented on Microsoft Durable Task."""

from __future__ import annotations

from datetime import timedelta
from typing import Any

from durabletask import task

from .contracts import (
    validate_evaluation,
    validate_extraction,
    validate_persistence,
    validate_request,
    validate_retention,
    validate_review,
    validate_review_event,
    validate_security,
)


def _outcome(
    request: dict[str, Any],
    status: str,
    *,
    receipt_ref: str = "",
    transaction_ref: str = "",
    action_required: str = "NONE",
) -> dict[str, str]:
    return {
        "status": status,
        "document_sha256": request["document_sha256"],
        "idempotency_key": request["idempotency_key"],
        "receipt_ref": receipt_ref,
        "transaction_ref": transaction_ref,
        "action_required": action_required,
    }


def document_pipeline_orchestrator(ctx: task.OrchestrationContext, raw_request: dict[str, Any]):
    """Run a hash-linked document pipeline without placing document fields in history."""
    request = validate_request(raw_request)
    document_sha256 = request["document_sha256"]
    ctx.set_custom_status({"stage": "security", "document_sha256": document_sha256})

    security = yield ctx.call_activity("documentflow_security", input={"request": request})
    security = validate_security(security, document_sha256)
    if security["status"] == "REJECTED":
        return _outcome(request, "REJECTED_SECURITY", receipt_ref=security["receipt_ref"])

    ctx.set_custom_status({"stage": "retain_original", "document_sha256": document_sha256})
    try:
        original = yield ctx.call_activity(
            "documentflow_retain_original", input={"request": request, "security": security}
        )
    except task.TaskFailedError:
        return _outcome(request, "UNKNOWN_ORIGINAL_STORAGE_EFFECT", action_required="RECONCILE_STORAGE")
    original = validate_retention(original, document_sha256, "original_evidence")
    if original["status"] != "VERIFIED":
        return _outcome(
            request,
            "AWAITING_ORIGINAL_STORAGE_RECONCILIATION",
            receipt_ref=original["receipt_ref"],
            action_required="RECONCILE_STORAGE",
        )

    ctx.set_custom_status({"stage": "extract", "document_sha256": document_sha256})
    extracted = yield ctx.call_activity(
        "documentflow_extract", input={"request": request, "original": original}
    )
    extracted = validate_extraction(extracted, request)
    if extracted["status"] == "REJECTED":
        return _outcome(request, "REJECTED_EXTRACTION", receipt_ref=extracted["receipt_ref"])

    evaluated = yield ctx.call_activity(
        "documentflow_evaluate", input={"request": request, "extraction": extracted}
    )
    evaluated = validate_evaluation(evaluated, document_sha256, extracted["result_sha256"])
    if evaluated["status"] == "REJECTED":
        return _outcome(request, "REJECTED_VALIDATION", receipt_ref=evaluated["receipt_ref"])

    review = None
    if extracted["status"] == "AMBIGUOUS" or evaluated["status"] == "REVIEW_REQUIRED":
        ctx.set_custom_status({"stage": "human_review", "document_sha256": document_sha256})
        review_event = ctx.wait_for_external_event("review_completed")
        timeout = ctx.create_timer(timedelta(seconds=request["review_timeout_seconds"]))
        winner = yield task.when_any([review_event, timeout])
        if winner == timeout:
            return _outcome(
                request,
                "REVIEW_TIMEOUT",
                receipt_ref=evaluated["receipt_ref"],
                action_required="HUMAN_REVIEW",
            )
        event = validate_review_event(review_event.get_result(), extracted["result_sha256"])
        review = yield ctx.call_activity(
            "documentflow_validate_review",
            input={"request": request, "extraction": extracted, "review_event": event},
        )
        review = validate_review(review, document_sha256, extracted["result_sha256"])
        if review["status"] == "REJECTED":
            return _outcome(request, "REVIEW_REJECTED", receipt_ref=review["receipt_ref"])

    ctx.set_custom_status({"stage": "business_persistence", "document_sha256": document_sha256})
    try:
        persisted = yield ctx.call_activity(
            "documentflow_persist",
            input={
                "request": request,
                "extraction": extracted,
                "evaluation": evaluated,
                "review": review,
            },
        )
    except task.TaskFailedError:
        return _outcome(request, "UNKNOWN_BUSINESS_EFFECT", action_required="RECONCILE_BUSINESS")
    persisted = validate_persistence(persisted, request, extracted["result_sha256"])

    ctx.set_custom_status({"stage": "retain_final", "document_sha256": document_sha256})
    try:
        final = yield ctx.call_activity(
            "documentflow_retain_final",
            input={
                "request": request,
                "original": original,
                "extraction": extracted,
                "evaluation": evaluated,
                "review": review,
                "persistence": persisted,
            },
        )
    except task.TaskFailedError:
        return _outcome(
            request,
            "PARTIAL_EFFECT",
            receipt_ref=persisted["receipt_ref"],
            transaction_ref=persisted["transaction_ref"],
            action_required="RECONCILE_FINAL_EVIDENCE",
        )
    final = validate_retention(final, document_sha256, "final_evidence")
    if final["status"] != "VERIFIED":
        return _outcome(
            request,
            "PARTIAL_EFFECT",
            receipt_ref=final["receipt_ref"],
            transaction_ref=persisted["transaction_ref"],
            action_required="RECONCILE_FINAL_EVIDENCE",
        )

    ctx.set_custom_status({"stage": "completed", "document_sha256": document_sha256})
    return _outcome(
        request,
        "COMPLETED",
        receipt_ref=final["receipt_ref"],
        transaction_ref=persisted["transaction_ref"],
    )
````

### FILE: `durable_document_orchestration/documentflow/runtime.py`
```yaml
block_id: "MICROSOFT-DURABLE-DOCUMENT:documentflow-runtime-py:v1"
operation: CREATE
provenance: ADAPTED
source: "adapted from Microsoft durabletask 1.9.0 TaskHubGrpcWorker.start at commit 6dfdbac521d9d59d31d7ea975fb669d66f86f4c4; local event-loop cleanup and document activity registration"
license: "LicenseRef-Workspace-Owner"
sha256: "ea6796d1d3f015e1d532f32e8d478997452ed3c169741b69797fc6e99d692b87"
variables: []
secrets_allowed: false
```
````text
"""Worker registration for project-supplied, explicitly named activities."""

from __future__ import annotations

import asyncio
from collections.abc import Callable
from dataclasses import dataclass
from threading import Thread
from typing import Any

from durabletask import task
from durabletask.worker import TaskHubGrpcWorker, WorkItemFilters

from .orchestration import document_pipeline_orchestrator

Activity = Callable[[task.ActivityContext, dict[str, Any]], dict[str, Any]]


@dataclass(frozen=True)
class ActivityHandlers:
    security: Activity
    retain_original: Activity
    extract: Activity
    evaluate: Activity
    validate_review: Activity
    persist: Activity
    retain_final: Activity


class ManagedTaskHubGrpcWorker(TaskHubGrpcWorker):
    """Pinned 1.9.0 lifecycle with deterministic event-loop cleanup.

    Microsoft durabletask 1.9.0 creates a loop in the worker thread but does not
    close it after ``_async_run_loop`` returns. This narrow subclass preserves
    the official start sequence and closes that owned loop after shutdown.
    """

    def start(self) -> None:
        if self._is_running:
            raise RuntimeError("The worker is already running.")
        if self._auto_generate_work_item_filters:
            self._work_item_filters = WorkItemFilters._from_registry(self._registry)
        self._shutdown.clear()

        def run_loop() -> None:
            loop = asyncio.new_event_loop()
            try:
                asyncio.set_event_loop(loop)
                loop.run_until_complete(self._async_run_loop())
                loop.run_until_complete(loop.shutdown_asyncgens())
                loop.run_until_complete(loop.shutdown_default_executor())
            finally:
                asyncio.set_event_loop(None)
                loop.close()

        self._logger.info(f"Starting gRPC worker that connects to {self._host_address}")
        self._runLoop = Thread(target=run_loop)
        self._runLoop.start()
        self._is_running = True


def _named_activity(name: str, handler: Activity) -> Activity:
    def activity(ctx: task.ActivityContext, payload: dict[str, Any]) -> dict[str, Any]:
        return handler(ctx, payload)

    activity.__name__ = name
    activity.__qualname__ = name
    return activity


def register_document_pipeline(worker: TaskHubGrpcWorker, handlers: ActivityHandlers) -> tuple[str, ...]:
    """Register the deterministic orchestrator and all explicit activity borders."""
    names = [worker.add_orchestrator(document_pipeline_orchestrator)]
    for name, handler in (
        ("documentflow_security", handlers.security),
        ("documentflow_retain_original", handlers.retain_original),
        ("documentflow_extract", handlers.extract),
        ("documentflow_evaluate", handlers.evaluate),
        ("documentflow_validate_review", handlers.validate_review),
        ("documentflow_persist", handlers.persist),
        ("documentflow_retain_final", handlers.retain_final),
    ):
        names.append(worker.add_activity(_named_activity(name, handler)))
    return tuple(names)
````

### FILE: `durable_document_orchestration/test_contracts.py`
```yaml
block_id: "MICROSOFT-DURABLE-DOCUMENT:test-contracts-py:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration around exact official Microsoft Durable Task APIs"
license: "LicenseRef-Workspace-Owner"
sha256: "79e9156721ed9dd051a930854bb906fb4b5a1a37a076b8de2d8f19aaff4dd801"
variables: []
secrets_allowed: false
```
````text
import unittest

from documentflow.contracts import (
    expected_idempotency_key,
    validate_request,
    validate_review_event,
)


def request() -> dict:
    document_sha = "a" * 64
    value = {
        "contract_version": 1,
        "tenant_ref": "tenant-01",
        "document_id": "document-01",
        "document_sha256": document_sha,
        "class_id": "commercial_invoice",
        "schema_version": "invoice-v1",
        "profile_sha256": "b" * 64,
        "routing_receipt_sha256": "c" * 64,
        "input_ref": "urn:quarantine:object-01",
        "review_timeout_seconds": 30,
        "evidence_mode": "TEST_ONLY_MEMORY",
    }
    value["idempotency_key"] = expected_idempotency_key(
        value["tenant_ref"], document_sha, value["class_id"], value["schema_version"]
    )
    return value


class ContractTests(unittest.TestCase):
    def test_valid_request_is_hash_bound(self):
        value = request()
        self.assertEqual(validate_request(value), value)

    def test_unknown_key_fails_closed(self):
        value = request()
        value["raw_invoice_number"] = "must-not-enter-history"
        with self.assertRaisesRegex(ValueError, "keys mismatch"):
            validate_request(value)

    def test_idempotency_key_mismatch_fails(self):
        value = request()
        value["idempotency_key"] = "0" * 64
        with self.assertRaisesRegex(ValueError, "not bound"):
            validate_request(value)

    def test_filesystem_path_is_not_an_opaque_input_reference(self):
        value = request()
        value["input_ref"] = "C:\\secret\\invoice.pdf"
        with self.assertRaisesRegex(ValueError, "opaque"):
            validate_request(value)

    def test_review_event_is_exact_and_hash_bound(self):
        result_sha = "d" * 64
        event = {
            "decision": "APPROVED",
            "subject_result_sha256": result_sha,
            "review_receipt_ref": "ref:review:01",
            "review_receipt_sha256": "e" * 64,
        }
        self.assertEqual(validate_review_event(event, result_sha), event)
        event["subject_result_sha256"] = "f" * 64
        with self.assertRaisesRegex(ValueError, "hash mismatch"):
            validate_review_event(event, result_sha)


if __name__ == "__main__":
    unittest.main()
````

### FILE: `durable_document_orchestration/test_pipeline.py`
```yaml
block_id: "MICROSOFT-DURABLE-DOCUMENT:test-pipeline-py:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration around exact official Microsoft Durable Task APIs"
license: "LicenseRef-Workspace-Owner"
sha256: "852a230d139881feb4a62601f3e6a8ee13522b3f0351379f7b3b3c688a6d9e64"
variables: []
secrets_allowed: false
```
````text
import json
import socket
import unittest

from durabletask import client
from durabletask.testing import create_test_backend

from documentflow import (
    ActivityHandlers,
    ManagedTaskHubGrpcWorker,
    expected_idempotency_key,
    register_document_pipeline,
)


def _free_port() -> int:
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        sock.bind(("127.0.0.1", 0))
        return int(sock.getsockname()[1])


def _request(timeout: int = 5) -> dict:
    document_sha = "a" * 64
    value = {
        "contract_version": 1,
        "tenant_ref": "tenant-01",
        "document_id": "document-01",
        "document_sha256": document_sha,
        "class_id": "commercial_invoice",
        "schema_version": "invoice-v1",
        "profile_sha256": "b" * 64,
        "routing_receipt_sha256": "c" * 64,
        "input_ref": "urn:quarantine:object-01",
        "review_timeout_seconds": timeout,
        "evidence_mode": "TEST_ONLY_MEMORY",
    }
    value["idempotency_key"] = expected_idempotency_key(
        value["tenant_ref"], document_sha, value["class_id"], value["schema_version"]
    )
    return value


class FakeHandlers:
    def __init__(self, *, security="ADMITTED", extraction="EXTRACTED", evaluation="AUTO_ACCEPT",
                 review="APPROVED", original="VERIFIED", final="VERIFIED",
                 fail_persist=False, fail_final=False):
        self.security_status = security
        self.extraction_status = extraction
        self.evaluation_status = evaluation
        self.review_status = review
        self.original_status = original
        self.final_status = final
        self.fail_persist = fail_persist
        self.fail_final = fail_final
        self.calls: list[str] = []

    @staticmethod
    def _base(stage: str, status: str, document_sha: str) -> dict:
        return {
            "stage": stage,
            "status": status,
            "document_sha256": document_sha,
            "receipt_ref": f"ref:receipt:{stage}",
            "receipt_sha256": "1" * 64,
        }

    def security(self, _, payload):
        self.calls.append("security")
        return self._base("security", self.security_status, payload["request"]["document_sha256"])

    def retain_original(self, _, payload):
        self.calls.append("retain_original")
        value = self._base("original_evidence", self.original_status, payload["request"]["document_sha256"])
        value.update({"object_ref": "ref:object:original", "version_ref": "ref:version:original"})
        return value

    def extract(self, _, payload):
        self.calls.append("extract")
        value = self._base("extraction", self.extraction_status, payload["request"]["document_sha256"])
        value.update({
            "result_ref": "ref:result:normalized",
            "result_sha256": "d" * 64,
            "schema_version": payload["request"]["schema_version"],
        })
        return value

    def evaluate(self, _, payload):
        self.calls.append("evaluate")
        value = self._base("evaluation", self.evaluation_status, payload["request"]["document_sha256"])
        value["result_sha256"] = payload["extraction"]["result_sha256"]
        return value

    def validate_review(self, _, payload):
        self.calls.append("validate_review")
        value = self._base("review", self.review_status, payload["request"]["document_sha256"])
        value["result_sha256"] = payload["extraction"]["result_sha256"]
        return value

    def persist(self, _, payload):
        self.calls.append("persist")
        if self.fail_persist:
            raise RuntimeError("simulated ambiguous persistence effect")
        value = self._base("business_persistence", "COMMITTED", payload["request"]["document_sha256"])
        value.update({
            "result_sha256": payload["extraction"]["result_sha256"],
            "idempotency_key": payload["request"]["idempotency_key"],
            "transaction_ref": "ref:transaction:01",
        })
        return value

    def retain_final(self, _, payload):
        self.calls.append("retain_final")
        if self.fail_final:
            raise RuntimeError("simulated final evidence failure")
        value = self._base("final_evidence", self.final_status, payload["request"]["document_sha256"])
        value.update({"object_ref": "ref:object:final", "version_ref": "ref:version:final"})
        return value

    def as_handlers(self) -> ActivityHandlers:
        return ActivityHandlers(
            security=self.security,
            retain_original=self.retain_original,
            extract=self.extract,
            evaluate=self.evaluate,
            validate_review=self.validate_review,
            persist=self.persist,
            retain_final=self.retain_final,
        )


def _run(handlers: FakeHandlers, request: dict, review_event: dict | None = None) -> dict:
    port = _free_port()
    host = f"localhost:{port}"
    backend = create_test_backend(port=port)
    try:
        with ManagedTaskHubGrpcWorker(host_address=host) as task_worker:
            names = register_document_pipeline(task_worker, handlers.as_handlers())
            if len(names) != 8 or len(set(names)) != 8:
                raise AssertionError("orchestrator/activity names are not unique")
            task_worker.start()
            with client.TaskHubGrpcClient(host_address=host) as task_client:
                instance_id = task_client.schedule_new_orchestration(
                    "document_pipeline_orchestrator", input=request
                )
                if review_event is not None:
                    task_client.raise_orchestration_event(instance_id, "review_completed", data=review_event)
                state = task_client.wait_for_orchestration_completion(instance_id, timeout=15)
                if state is None:
                    raise AssertionError("orchestration state missing")
                state.raise_if_failed()
                return json.loads(state.serialized_output)
    finally:
        backend.stop()
        backend.reset()


class PipelineTests(unittest.TestCase):
    def test_auto_accept_runs_every_effect_once(self):
        handlers = FakeHandlers()
        result = _run(handlers, _request())
        self.assertEqual(result["status"], "COMPLETED")
        self.assertEqual(
            handlers.calls,
            ["security", "retain_original", "extract", "evaluate", "persist", "retain_final"],
        )

    def test_security_rejection_stops_before_storage(self):
        handlers = FakeHandlers(security="REJECTED")
        result = _run(handlers, _request())
        self.assertEqual(result["status"], "REJECTED_SECURITY")
        self.assertEqual(handlers.calls, ["security"])

    def test_original_partial_verification_stops_for_reconciliation(self):
        handlers = FakeHandlers(original="PARTIAL_VERIFICATION")
        result = _run(handlers, _request())
        self.assertEqual(result["status"], "AWAITING_ORIGINAL_STORAGE_RECONCILIATION")
        self.assertEqual(result["action_required"], "RECONCILE_STORAGE")
        self.assertEqual(handlers.calls, ["security", "retain_original"])

    def test_human_approval_is_hash_bound_then_persists(self):
        handlers = FakeHandlers(evaluation="REVIEW_REQUIRED")
        event = {
            "decision": "APPROVED",
            "subject_result_sha256": "d" * 64,
            "review_receipt_ref": "ref:review:01",
            "review_receipt_sha256": "e" * 64,
        }
        result = _run(handlers, _request(), event)
        self.assertEqual(result["status"], "COMPLETED")
        self.assertIn("validate_review", handlers.calls)
        self.assertEqual(handlers.calls.count("persist"), 1)

    def test_human_rejection_never_persists(self):
        handlers = FakeHandlers(evaluation="REVIEW_REQUIRED", review="REJECTED")
        event = {
            "decision": "REJECTED",
            "subject_result_sha256": "d" * 64,
            "review_receipt_ref": "ref:review:02",
            "review_receipt_sha256": "f" * 64,
        }
        result = _run(handlers, _request(), event)
        self.assertEqual(result["status"], "REVIEW_REJECTED")
        self.assertNotIn("persist", handlers.calls)

    def test_review_timeout_never_persists(self):
        handlers = FakeHandlers(evaluation="REVIEW_REQUIRED")
        result = _run(handlers, _request(timeout=1))
        self.assertEqual(result["status"], "REVIEW_TIMEOUT")
        self.assertNotIn("persist", handlers.calls)

    def test_ambiguous_persistence_effect_is_not_retried(self):
        handlers = FakeHandlers(fail_persist=True)
        result = _run(handlers, _request())
        self.assertEqual(result["status"], "UNKNOWN_BUSINESS_EFFECT")
        self.assertEqual(result["action_required"], "RECONCILE_BUSINESS")
        self.assertEqual(handlers.calls.count("persist"), 1)
        self.assertNotIn("retain_final", handlers.calls)

    def test_final_evidence_failure_preserves_committed_transaction(self):
        handlers = FakeHandlers(fail_final=True)
        result = _run(handlers, _request())
        self.assertEqual(result["status"], "PARTIAL_EFFECT")
        self.assertEqual(result["transaction_ref"], "ref:transaction:01")
        self.assertEqual(result["action_required"], "RECONCILE_FINAL_EVIDENCE")
        self.assertEqual(handlers.calls.count("persist"), 1)
        self.assertEqual(handlers.calls.count("retain_final"), 1)


if __name__ == "__main__":
    unittest.main()
````

### FILE: `durable_document_orchestration/README.md`
```yaml
block_id: "MICROSOFT-DURABLE-DOCUMENT:readme-md:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration around exact official Microsoft Durable Task APIs"
license: "LicenseRef-Workspace-Owner"
sha256: "3761424d1b43dba53139f34b998f4cfb8aa79d808e0d5f7ad5243dcfa3e04cb4"
variables: []
secrets_allowed: false
```
````text
# Durable document orchestration

This directory materializes a fail-closed document workflow on Microsoft `durabletask==1.9.0`. The SDK, in-memory test backend and upstream agent skill are Microsoft MIT code; `documentflow` is local integration glue and is never attributed to Microsoft. `ManagedTaskHubGrpcWorker` is explicitly `ADAPTED` from the pinned official `TaskHubGrpcWorker.start`: it preserves that sequence and closes the worker-owned event loop after shutdown because the exact 1.9.0 implementation does not.

The orchestration carries only opaque references, hashes and receipts. Document bytes, extracted field values, reviewer identity and secrets must remain in activities and approved stores. The sequence is security admission, retained original, extraction, strict evaluation, optional durable human review, idempotent business persistence and retained final evidence.

Install the frozen Windows/CPython 3.12 graph:

```powershell
python -m pip install --require-hashes -r requirements-windows-py312.lock
python -m pip check
python -m unittest -v
```

`durabletask.testing.create_test_backend` is official Microsoft code and runs the integration suite without Docker or Azure. It is test-only and in-memory; it is not production persistence and not WORM. Production requires an approved Durable Task Scheduler endpoint or another supported backend, managed identity/RBAC/TLS, a task hub, observability, load/recovery/versioning evidence and one of the separately verified S3/GCS/Azure immutable evidence lanes.

`upstream/durable-task-python-SKILL.md` is retained verbatim as official agent guidance. `official-artifact-lock.json` also binds Microsoft's exact standalone and Azure Functions human-interaction samples: they prove external-event plus durable-timeout mechanics, not approver identity. The Azure HTTP sample is explicitly anonymous. Its unpinned `pip install` and `mcr...:latest` examples are reference-only: never execute them automatically. Resolve exact package hashes and an immutable container digest through the project admission process first. The distributed profile authorizes nothing.
````

### FILE: `durable_document_orchestration/upstream/LICENSE.md`
```yaml
block_id: "MICROSOFT-DURABLE-DOCUMENT:upstream-license-md:v1"
operation: CREATE
provenance: ADAPTED
source: "Azure-Samples/Durable-Task-Scheduler 5636c25ffbdaabaca1062c3b6c6e77072f2b2eb2 LICENSE.md; final newline only"
license: "MIT"
sha256: "c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383"
variables: []
secrets_allowed: false
```
````text
    MIT License

    Copyright (c) Microsoft Corporation.

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE
````

### FILE: `durable_document_orchestration/upstream/durable-task-python-SKILL.md`
```yaml
block_id: "MICROSOFT-DURABLE-DOCUMENT:upstream-durable-task-python-skill-md:v1"
operation: CREATE
provenance: VERBATIM
source: "Azure-Samples/Durable-Task-Scheduler 5636c25ffbdaabaca1062c3b6c6e77072f2b2eb2 exact .github skill"
license: "MIT"
sha256: "553f0dc278c63bb129df5557a7d6f7d645bf827efba4cb077588688f88a2e9e2"
variables: []
secrets_allowed: false
```
````text
---
name: durable-task-python
description: Build durable, fault-tolerant workflows in Python using the Durable Task SDK with Azure Durable Task Scheduler. Use when creating orchestrations, activities, entities, or implementing patterns like function chaining, fan-out/fan-in, human interaction, or stateful agents. Applies to any Python application requiring durable execution, state persistence, or distributed transactions without Azure Functions dependency.
---

# Durable Task Python SDK with Durable Task Scheduler

Build fault-tolerant, stateful workflows in Python applications using the Durable Task SDK connected to Azure Durable Task Scheduler.

## Quick Start

### Required Packages

```bash
pip install durabletask durabletask-azuremanaged azure-identity
```

Or add to `requirements.txt`:

```text
durabletask
durabletask-azuremanaged
azure-identity
```

### Minimal Worker + Client Setup

```python
import os
from azure.identity import DefaultAzureCredential
from durabletask import task
from durabletask.client import OrchestrationStatus
from durabletask.azuremanaged.client import DurableTaskSchedulerClient
from durabletask.azuremanaged.worker import DurableTaskSchedulerWorker


# Activity function
def hello(ctx: task.ActivityContext, name: str) -> str:
    return f"Hello {name}!"


# Orchestrator function
def my_orchestration(ctx: task.OrchestrationContext, input: str):
    result = yield ctx.call_activity(hello, input=input)
    return result


# Configuration - defaults to local emulator
taskhub = os.getenv("TASKHUB", "default")
endpoint = os.getenv("ENDPOINT", "http://localhost:8080")
secure_channel = endpoint != "http://localhost:8080"
credential = None if endpoint == "http://localhost:8080" else DefaultAzureCredential()

# Start worker and run orchestration
with DurableTaskSchedulerWorker(
    host_address=endpoint,
    secure_channel=secure_channel,
    taskhub=taskhub,
    token_credential=credential
) as worker:
    worker.add_orchestrator(my_orchestration)
    worker.add_activity(hello)
    worker.start()

    # Create client and schedule orchestration
    dts_client = DurableTaskSchedulerClient(
        host_address=endpoint,
        secure_channel=secure_channel,
        taskhub=taskhub,
        token_credential=credential
    )
    
    instance_id = dts_client.schedule_new_orchestration(my_orchestration, input="World")
    state = dts_client.wait_for_orchestration_completion(instance_id, timeout=60)
    
    if state and state.runtime_status == OrchestrationStatus.COMPLETED:
        print(f"Result: {state.serialized_output}")
```

## Pattern Selection Guide

| Pattern | Use When |
|---------|----------|
| **Function Chaining** | Sequential steps where each depends on the previous |
| **Fan-Out/Fan-In** | Parallel processing with aggregated results |
| **Human Interaction** | Workflow pauses for external input/approval |
| **Durable Entities** | Stateful objects with operations (counters, accounts) |
| **Sub-Orchestrations** | Reusable workflow components or version isolation |
| **Eternal Orchestrations** | Long-running background processes with `continue_as_new` |
| **Monitoring** | Periodic polling with configurable timeouts |

See [references/patterns.md](references/patterns.md) for detailed implementations.

## Orchestration Structure

### Basic Orchestrator

```python
def my_orchestration(ctx: task.OrchestrationContext, input: str):
    """Orchestrator function - MUST be deterministic"""
    # Call activities sequentially
    step1 = yield ctx.call_activity(step1_activity, input=input)
    step2 = yield ctx.call_activity(step2_activity, input=step1)
    return step2
```

### Basic Activity

```python
def my_activity(ctx: task.ActivityContext, input: str) -> str:
    """Activity function - can have side effects, I/O, non-determinism"""
    # Perform actual work here
    print(f"Processing: {input}")
    return f"Processed: {input}"
```

### Registering with Worker

```python
with DurableTaskSchedulerWorker(...) as worker:
    worker.add_orchestrator(my_orchestration)
    worker.add_activity(step1_activity)
    worker.add_activity(step2_activity)
    worker.start()
```

## Critical Rules

### Orchestration Determinism

Orchestrations replay from history - all code MUST be deterministic. When an orchestration resumes, it replays all previous code to rebuild state. Non-deterministic code produces different results on replay, causing failures.

**NEVER do inside orchestrations:**
- `datetime.now()`, `datetime.utcnow()` → Use `ctx.current_utc_datetime`
- `uuid.uuid4()` → Use `ctx.new_uuid()`
- `random.random()` → Pass random values from activities
- Direct I/O, HTTP calls, database access → Move to activities
- `time.sleep()`, `asyncio.sleep()` → Use `ctx.create_timer()`
- Environment variables that may change → Pass as input or use activities
- Global mutable state → Pass state through activity results

**ALWAYS use:**
- `yield ctx.call_activity()` - Call activities
- `yield ctx.call_sub_orchestrator()` - Call sub-orchestrations
- `yield ctx.create_timer()` - Durable delays
- `yield ctx.wait_for_external_event()` - Wait for events
- `ctx.current_utc_datetime` - Current time
- `ctx.new_uuid()` - Generate GUIDs
- `ctx.set_custom_status()` - Set status

### Non-Determinism Patterns (WRONG vs CORRECT)

#### Getting Current Time

```python
# WRONG - datetime.now() returns different value on replay
def bad_orchestration(ctx: task.OrchestrationContext, _):
    current_time = datetime.now()  # Non-deterministic!
    if current_time.hour < 12:
        yield ctx.call_activity(morning_activity)

# CORRECT - ctx.current_utc_datetime is replayed consistently
def good_orchestration(ctx: task.OrchestrationContext, _):
    current_time = ctx.current_utc_datetime  # Deterministic
    if current_time.hour < 12:
        yield ctx.call_activity(morning_activity)
```

#### Generating UUIDs/Random Values

```python
# WRONG - uuid4() generates different value on replay
def bad_orchestration(ctx: task.OrchestrationContext, _):
    order_id = str(uuid.uuid4())  # Non-deterministic!
    yield ctx.call_activity(create_order, input=order_id)

# CORRECT - ctx.new_uuid() replays the same value
def good_orchestration(ctx: task.OrchestrationContext, _):
    order_id = str(ctx.new_uuid())  # Deterministic
    yield ctx.call_activity(create_order, input=order_id)
```

#### Random Numbers

```python
# WRONG - random produces different values on replay
def bad_orchestration(ctx: task.OrchestrationContext, _):
    delay = random.randint(1, 10)  # Non-deterministic!
    yield ctx.create_timer(timedelta(seconds=delay))

# CORRECT - generate random in activity, pass to orchestrator
def get_random_delay(ctx: task.ActivityContext, _) -> int:
    return random.randint(1, 10)  # OK in activity

def good_orchestration(ctx: task.OrchestrationContext, _):
    delay = yield ctx.call_activity(get_random_delay)  # Deterministic
    yield ctx.create_timer(timedelta(seconds=delay))
```

#### Sleeping/Delays

```python
# WRONG - time.sleep blocks and doesn't persist
def bad_orchestration(ctx: task.OrchestrationContext, _):
    yield ctx.call_activity(step1)
    time.sleep(60)  # Non-durable! Lost on restart
    yield ctx.call_activity(step2)

# CORRECT - ctx.create_timer is durable
def good_orchestration(ctx: task.OrchestrationContext, _):
    yield ctx.call_activity(step1)
    yield ctx.create_timer(timedelta(seconds=60))  # Durable timer
    yield ctx.call_activity(step2)
```

#### HTTP Calls and I/O

```python
# WRONG - HTTP call in orchestrator is non-deterministic
def bad_orchestration(ctx: task.OrchestrationContext, url: str):
    import requests
    response = requests.get(url)  # Non-deterministic!
    return response.json()

# CORRECT - move I/O to activity
def fetch_data(ctx: task.ActivityContext, url: str) -> dict:
    import requests
    response = requests.get(url)  # OK in activity
    return response.json()

def good_orchestration(ctx: task.OrchestrationContext, url: str):
    data = yield ctx.call_activity(fetch_data, input=url)  # Deterministic
    return data
```

#### Database Access

```python
# WRONG - database query in orchestrator
def bad_orchestration(ctx: task.OrchestrationContext, user_id: str):
    import sqlite3
    conn = sqlite3.connect('db.sqlite')  # Non-deterministic!
    cursor = conn.execute("SELECT * FROM users WHERE id=?", (user_id,))
    user = cursor.fetchone()
    # ...

# CORRECT - database access in activity
def get_user(ctx: task.ActivityContext, user_id: str) -> dict:
    import sqlite3
    conn = sqlite3.connect('db.sqlite')  # OK in activity
    cursor = conn.execute("SELECT * FROM users WHERE id=?", (user_id,))
    return dict(cursor.fetchone())

def good_orchestration(ctx: task.OrchestrationContext, user_id: str):
    user = yield ctx.call_activity(get_user, input=user_id)
    # ...
```

#### Environment Variables

```python
# WRONG - env var might change between replays
def bad_orchestration(ctx: task.OrchestrationContext, _):
    api_endpoint = os.getenv("API_ENDPOINT")  # Could change!
    yield ctx.call_activity(call_api, input=api_endpoint)

# CORRECT - pass config as input or read in activity
def good_orchestration(ctx: task.OrchestrationContext, config: dict):
    api_endpoint = config["api_endpoint"]  # From input, deterministic
    yield ctx.call_activity(call_api, input=api_endpoint)

# ALSO CORRECT - read env var in activity
def call_api(ctx: task.ActivityContext, _) -> str:
    api_endpoint = os.getenv("API_ENDPOINT")  # OK in activity
    # make the call...
```

#### Conditional Logic Based on External State

```python
# WRONG - file existence can change between replays
def bad_orchestration(ctx: task.OrchestrationContext, path: str):
    if os.path.exists(path):  # Non-deterministic!
        yield ctx.call_activity(process_file, input=path)

# CORRECT - check in activity
def check_file_exists(ctx: task.ActivityContext, path: str) -> bool:
    return os.path.exists(path)  # OK in activity

def good_orchestration(ctx: task.OrchestrationContext, path: str):
    exists = yield ctx.call_activity(check_file_exists, input=path)
    if exists:  # Deterministic - based on activity result
        yield ctx.call_activity(process_file, input=path)
```

#### Dictionary/Set Iteration Order

```python
# POTENTIALLY WRONG - dict iteration order may vary (Python < 3.7)
def risky_orchestration(ctx: task.OrchestrationContext, items: dict):
    for key in items:  # Order might not be guaranteed
        yield ctx.call_activity(process, input=key)

# CORRECT - use sorted keys for deterministic order
def good_orchestration(ctx: task.OrchestrationContext, items: dict):
    for key in sorted(items.keys()):  # Guaranteed order
        yield ctx.call_activity(process, input=key)
```

#### Thread-Local or Global State

```python
# WRONG - global state can change
counter = 0

def bad_orchestration(ctx: task.OrchestrationContext, _):
    global counter
    counter += 1  # Non-deterministic across replays!
    yield ctx.call_activity(process, input=counter)

# CORRECT - pass state through orchestration input/output
def good_orchestration(ctx: task.OrchestrationContext, counter: int):
    counter += 1  # Local variable, deterministic
    yield ctx.call_activity(process, input=counter)
    # If continuing, pass counter forward
    ctx.continue_as_new(counter)
```

### Using yield

In Python, orchestrator functions use `yield` to await durable operations:

```python
# CORRECT - use yield
result = yield ctx.call_activity(my_activity, input="data")

# WRONG - will not work
result = ctx.call_activity(my_activity, input="data")  # Missing yield!
```

### Error Handling

```python
def orchestrator_with_error_handling(ctx: task.OrchestrationContext, input: str):
    try:
        result = yield ctx.call_activity(risky_activity, input=input)
        return result
    except task.TaskFailedError as e:
        # Activity failed - implement compensation
        ctx.set_custom_status({"error": str(e)})
        yield ctx.call_activity(compensation_activity, input=input)
        return "Compensated"
```

### Retry Policies

```python
from durabletask.task import RetryPolicy

retry_policy = RetryPolicy(
    first_retry_interval=5,  # seconds
    max_number_of_attempts=3,
    backoff_coefficient=2.0,
    max_retry_interval=60,  # seconds
    retry_timeout=300  # seconds
)

def orchestrator(ctx: task.OrchestrationContext, _):
    result = yield ctx.call_activity(
        unreliable_activity, 
        input="data",
        retry_policy=retry_policy
    )
    return result
```

## Working with Custom Types

The SDK supports dataclasses, namedtuples, and custom classes:

```python
from dataclasses import dataclass

@dataclass
class Order:
    product: str
    quantity: int
    cost: float

def process_order(ctx: task.ActivityContext, order: Order) -> str:
    return f"Processed {order.quantity}x {order.product}"

def order_workflow(ctx: task.OrchestrationContext, order: Order):
    result = yield ctx.call_activity(process_order, input=order)
    return result
```

## Connection & Authentication

### Local Emulator (Default)

```python
# No authentication required
taskhub = "default"
endpoint = "http://localhost:8080"
credential = None
secure_channel = False
```

### Azure with DefaultAzureCredential

```python
from azure.identity import DefaultAzureCredential

taskhub = "my-taskhub"
endpoint = "https://my-scheduler.region.durabletask.io"
credential = DefaultAzureCredential()
secure_channel = True
```

### Authentication Helper

```python
def get_connection_config():
    endpoint = os.getenv("ENDPOINT", "http://localhost:8080")
    taskhub = os.getenv("TASKHUB", "default")
    
    is_local = endpoint == "http://localhost:8080"
    
    return {
        "host_address": endpoint,
        "taskhub": taskhub,
        "secure_channel": not is_local,
        "token_credential": None if is_local else DefaultAzureCredential()
    }

config = get_connection_config()
worker = DurableTaskSchedulerWorker(**config)
client = DurableTaskSchedulerClient(**config)
```

## Local Development with Emulator

```bash
# Pull and run the emulator
docker pull mcr.microsoft.com/dts/dts-emulator:latest
docker run -d -p 8080:8080 -p 8082:8082 --name dts-emulator mcr.microsoft.com/dts/dts-emulator:latest

# Dashboard available at http://localhost:8082
```

## Client Operations

```python
# Schedule new orchestration
instance_id = client.schedule_new_orchestration(my_orchestration, input="data")

# Schedule with custom instance ID
instance_id = client.schedule_new_orchestration(
    my_orchestration, 
    input="data",
    instance_id="my-custom-id"
)

# Wait for completion
state = client.wait_for_orchestration_completion(instance_id, timeout=60)

# Get current status
state = client.get_orchestration_state(instance_id)

# Raise external event
client.raise_orchestration_event(instance_id, "approval_received", data=approval_data)

# Terminate orchestration
client.terminate_orchestration(instance_id, output="User cancelled")

# Suspend/Resume
client.suspend_orchestration(instance_id)
client.resume_orchestration(instance_id)
```

## References

- **[patterns.md](references/patterns.md)** - Detailed pattern implementations (Fan-Out/Fan-In, Human Interaction, Entities, Sub-Orchestrations)
- **[setup.md](references/setup.md)** - Azure Durable Task Scheduler provisioning and deployment
````

## 6. Configuration surface

- `pipeline-profile.template.json` exige backend, host, task hub, TLS/identity, clases, storage lane, timeout, tamaño, owner y approvals; nace deshabilitado.
- El request runtime admite sólo tokens/referencias opacas y tres hashes de autoridad; rechaza paths, claves extra y claves idempotentes no derivadas.
- Backends productivos y adapters de actividad se inyectan desde el proyecto; ningún secreto vive en el pack.

## 7. Dependency bill

- Microsoft `durabletask==1.9.0`, commit firmado `6dfdbac521d9d59d31d7ea975fb669d66f86f4c4`, MIT; samples oficiales Human Interaction standalone/Azure Functions fijados por path y SHA-256, sólo como autoridad del patrón durable, nunca como auth.
- Lock CPython 3.12/Windows x86-64 de 6 wheels con hashes; 0 hallazgos OSV al 2026-08-27.
- Skill oficial Azure Samples commit firmado `5636c25ffbdaabaca1062c3b6c6e77072f2b2eb2`, MIT, conservado VERBATIM.
- El glue `documentflow` y sus pruebas son AUTHORED y no se atribuyen a Microsoft.

## 8. Apply order

1. Ejecutar preflight, routing documental y adquisición oficial.
2. Materializar el pack y verificar hashes/licencia/skill.
3. Instalar el lock con `--require-hashes`; ejecutar compileall, pip-check y unittest.
4. Completar el perfil; seleccionar handlers ya admitidos para seguridad, extracción, evaluación, persistencia y storage.
5. Probar con backend in-memory; luego ejecutar emulator/live sólo con backend, digest, identidad, costo y efectos aprobados.

## 9. Verification

- materialización 13/13 y SHA-256 por archivo;
- wheel↔sdist↔commit 60/60 Python byte-identical;
- suite core Microsoft 1.034 PASS, 1 skip, 14 deselected; 12 warnings datetime retenidos;
- suite local 13 PASS sobre el backend in-memory oficial; pipeline 8/8 adicional con `ResourceWarning` elevado a error y cero event loops/sockets propios sin cerrar;
- install frozen `--require-hashes`, pip-check, compileall y OSV 6/0;
- negativos para raw history, hash/idempotencia, security reject, storage parcial, review reject/timeout, persistencia ambigua y evidencia final parcial;
- emulator DTS, Azure live y carga/soak/recovery siguen condicionados; el lifecycle local usa `ManagedTaskHubGrpcWorker` `ADAPTED` y cierra el loop propiedad del thread en la revisión 0.1.2.

## 10. Reconstruction evidence

Expedientes canónicos: `reconstruction_evidence/MICROSOFT_DURABLE_DOCUMENT_ORCHESTRATION_2026-08-27_V1.md` y `reconstruction_evidence/MICROSOFT_DURABLE_HUMAN_INTERACTION_2026-08-27_V1.md`. El segundo liga los samples oficiales exactos, demuestra 13 pruebas sobre el backend oficial de test y conserva su ausencia de auth. Hasta que exista backend productivo y adapters reales aprobados, el estado es `REBUILD_VERIFIED / CONDITIONED`, nunca garantía de producción.
