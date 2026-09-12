# Azure Content Understanding Document Runtime

## 1. Metadata

```yaml
pack_id: "AZURE-CONTENT-UNDERSTANDING-DOCUMENT-RUNTIME"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un runtime Python y lock productivo de nueve wheels para CPython 3.12/Windows x86-64 que usan el SDK GA oficial Microsoft Content Understanding 1.1.0 para analizar únicamente archivos admitidos por el gate local y clases REQUIRED del perfil, preservar resultado completo y encadenar hashes de input, seguridad, configuración y provider de forma atómica."
stacks: ["Python 3.10+", "azure-ai-contentunderstanding 1.1.0", "Azure Content Understanding API 2025-11-01"]
compatible_with: ["OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE 0.1.x", "OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x", "STRICT-DOCUMENT-FIELD-EVALUATION-GATE 0.1.x"]
incompatible_with: ["SAS URL obligatorio", "secreto por CLI", "auto-storage sin corpus", "archivo remoto no controlado"]
license_expression: "LicenseRef-Workspace-Owner AND MIT"
upstream_sources: ["https://pypi.org/project/azure-ai-contentunderstanding/1.1.0/", "https://github.com/Azure/azure-sdk-for-python/tree/129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb/sdk/contentunderstanding/azure-ai-contentunderstanding", "https://github.com/Azure-Samples/azure-ai-content-understanding-python/tree/fbc18880179397729abfa9506bda865d009dc711"]
verified_at: "2026-08-27"
```

## 2. Applicability

Use después de seleccionar Azure Content Understanding, adquirir/instalar el wheel oficial exacto y producir un receipt `ADMITTED` con el gate local. Acepta documentos locales PDF/PNG/JPEG/TIFF, una clase declarada `REQUIRED` en el perfil y credenciales entregadas por environment reference. Rechazar cuando el provider elegido sea otro, falte cuenta/endpoint/región/analyzer, el input sea remoto/symlink/otro formato, el receipt no coincida byte a byte o se pretenda almacenar campos automáticamente sólo por recibir un resultado.

El runtime es `AUTHORED` glue mínimo sobre `ContentUnderstandingClient.begin_analyze_binary`; el SDK/model/result son Microsoft. El template sólo activa los prebuilt IDs documentados para invoice y purchase order como candidatos; proforma, packing list, BOL, customs y clases custom permanecen bloqueadas hasta analyzer/corpus/evaluación propios.

## 3. Architecture contract

El runner valida versión de distribución, profile exacto, clase `REQUIRED`, archivo regular no symlink, extensión y 1..20 MiB antes de leer. El analyzer nunca es libre por CLI: se deriva del profile. Exige además un receipt del gate local con schema, decisión `ADMITTED`, prohibición explícita de storage y coincidencia exacta de SHA/bytes/MIME. Envía bytes locales al método oficial con MIME explícito; no crea SAS ni recibe secretos en CLI. Después del LRO exige contenidos, verifica analyzer y API GA cuando el provider los devuelve, guarda el resultado completo y un receipt con hashes/bytes/versiones en staging, y renombra sólo al finalizar.

Un error de input, SDK, servicio, analyzer o output elimina staging y no crea destino final. El resultado conserva fields/confidence/source oficiales sin traducirlos a hechos de negocio. Performance/costo/latencia dependen del servicio, región, documento y analyzer y se miden en sandbox. Rollback conserva original/result/receipt, deja de enrutar al analyzer nuevo y reevalúa schema/model antes de reprocess.

## 4. Exact file manifest

```text
CREATE azure_document_runtime/requirements-direct.in
CREATE azure_document_runtime/requirements-windows-py312.lock
CREATE azure_document_runtime/document-profile.template.json
CREATE azure_document_runtime/run_content_understanding.py
CREATE azure_document_runtime/test_content_understanding_runtime.py
CREATE azure_document_runtime/README.md
```

## 5. Materialization blocks

### FILE: `azure_document_runtime/requirements-direct.in`
```yaml
block_id: "AZURE-CU-DOCUMENT-RUNTIME:requirement:v1"
operation: CREATE
provenance: AUTHORED
source: "exact official PyPI wheel URL and SHA-256"
license: "LicenseRef-Workspace-Owner"
sha256: "0f701ee7b6efdc3638a391769fae4bd8ba32103ae6422fff96e5a448329a1eb2"
variables: []
secrets_allowed: false
```
````text
azure-ai-contentunderstanding @ https://files.pythonhosted.org/packages/67/d6/6f60ac39b73a0b8a2bc86cd6895e26eccc20d9b7bde86814fe6a36964dff/azure_ai_contentunderstanding-1.1.0-py3-none-any.whl#sha256=d1d6bdeffe02f5c8cc5309f0e120a1338f51af04638605f9617b7b501e2d1dd6
````

### FILE: `azure_document_runtime/requirements-windows-py312.lock`
```yaml
block_id: "AZURE-CU-DOCUMENT-RUNTIME:windows-py312-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "dated resolution of exact official PyPI wheels for CPython 3.12 Windows x86-64"
license: "LicenseRef-Workspace-Owner"
sha256: "ee950561fd78e273ffab6384cd4c303f0ef785c7236c59804d4c847f15381d5f"
variables: []
secrets_allowed: false
```
````text
azure-ai-contentunderstanding==1.1.0 --hash=sha256:d1d6bdeffe02f5c8cc5309f0e120a1338f51af04638605f9617b7b501e2d1dd6
azure-core==1.41.0 --hash=sha256:522b4011e8180b1a3dcd2024396a4e7fe9ac37fb8597db47163d230b5efe892d
certifi==2026.7.22 --hash=sha256:62f22742b58a1a33014a2b6b706588a8d7e2a88ae7bd1a6ebe8c992928483775
charset-normalizer==3.5.1 --hash=sha256:3617ac3cfd8b9888f145ad89dd6e692285834b0201c6074a5eeaad3fd4d668c2
idna==3.19 --hash=sha256:815e7be7a7806d54abb586dc943addc79e8b2ee16915059658cbeff4b1b43bf4
isodate==0.7.2 --hash=sha256:28009937d8031054830160fce6d409ed342816b543597cece116d966c6d99e15
requests==2.34.2 --hash=sha256:2a0d60c172f83ac6ab31e4554906c0f3b3588d37b5cb939b1c061f4907e278e0
typing-extensions==4.16.0 --hash=sha256:481caa481374e813c1b176ada14e97f1f67a4539ce9cfeb3f350d78d6370c2e8
urllib3==2.7.0 --hash=sha256:9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897
````

### FILE: `azure_document_runtime/document-profile.template.json`
```yaml
block_id: "AZURE-CU-DOCUMENT-RUNTIME:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed profile using Microsoft documented prebuilt analyzer IDs"
license: "LicenseRef-Workspace-Owner"
sha256: "c022818b366477cb992d30e391457d6219e14ea53179245880641703c3c0e438"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-azure-content-understanding-document-profile/v1",
  "provider": "Microsoft Azure Content Understanding",
  "sdk": "azure-ai-contentunderstanding==1.1.0",
  "api_version": "2025-11-01",
  "classes": [
    {"id": "supplier-invoice", "decision": "REQUIRED", "analyzer_id": "prebuilt-invoice", "schema_version": "PROVIDER_PREBUILT", "automatic_storage": false},
    {"id": "purchase-order", "decision": "REQUIRED", "analyzer_id": "prebuilt-purchaseOrder", "schema_version": "PROVIDER_PREBUILT", "automatic_storage": false},
    {"id": "proforma-invoice", "decision": "BLOCKED_CUSTOM_ANALYZER_REQUIRED", "analyzer_id": "", "schema_version": "", "automatic_storage": false},
    {"id": "packing-list", "decision": "BLOCKED_CUSTOM_ANALYZER_REQUIRED", "analyzer_id": "", "schema_version": "", "automatic_storage": false},
    {"id": "commercial-invoice", "decision": "BLOCKED_EVALUATION_REQUIRED", "analyzer_id": "prebuilt-invoice", "schema_version": "PROVIDER_PREBUILT_BASELINE_ONLY", "automatic_storage": false},
    {"id": "bill-of-lading", "decision": "BLOCKED_CUSTOM_ANALYZER_REQUIRED", "analyzer_id": "", "schema_version": "", "automatic_storage": false},
    {"id": "delivery-note", "decision": "BLOCKED_CUSTOM_ANALYZER_REQUIRED", "analyzer_id": "", "schema_version": "", "automatic_storage": false},
    {"id": "certificate-of-origin", "decision": "BLOCKED_CUSTOM_ANALYZER_REQUIRED", "analyzer_id": "", "schema_version": "", "automatic_storage": false},
    {"id": "quality-inspection-certificate", "decision": "BLOCKED_CUSTOM_ANALYZER_REQUIRED", "analyzer_id": "", "schema_version": "", "automatic_storage": false},
    {"id": "customs-declaration", "decision": "BLOCKED_CUSTOM_ANALYZER_REQUIRED", "analyzer_id": "", "schema_version": "", "automatic_storage": false}
  ]
}
````

### FILE: `azure_document_runtime/run_content_understanding.py`
```yaml
block_id: "AZURE-CU-DOCUMENT-RUNTIME:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local wrapper invoking Microsoft GA SDK public API"
license: "LicenseRef-Workspace-Owner"
sha256: "0e1fc1021eb5defba64a1cba57135fc371c34acbd4c8e08ff90dc500a4e2884e"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
from dataclasses import dataclass
from datetime import datetime, timezone
import hashlib
import importlib.metadata
import json
import os
from pathlib import Path
import shutil
import sys
import uuid
from typing import Any, Protocol

from azure.ai.contentunderstanding import ContentUnderstandingClient
from azure.core.credentials import AzureKeyCredential


SDK_VERSION = "1.1.0"
API_VERSION = "2025-11-01"
PROFILE_SCHEMA = "elite-azure-content-understanding-document-profile/v1"
SECURITY_RECEIPT_SCHEMA = "elite-secure-local-file-receipt/v1"
ALLOWED_MIME_TYPES = {
    ".pdf": "application/pdf",
    ".png": "image/png",
    ".jpg": "image/jpeg",
    ".jpeg": "image/jpeg",
    ".tif": "image/tiff",
    ".tiff": "image/tiff",
}


class Poller(Protocol):
    def result(self) -> Any: ...


class AnalyzerClient(Protocol):
    def begin_analyze_binary(self, analyzer_id: str, binary_input: bytes, *, content_type: str) -> Poller: ...


@dataclass(frozen=True)
class AnalysisRequest:
    document_class: str
    input_path: Path
    output_directory: Path
    profile_path: Path
    security_receipt_path: Path
    max_bytes: int = 20 * 1024 * 1024


def sha256_bytes(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def read_json_file(path: Path, label: str) -> tuple[dict[str, Any], bytes]:
    resolved = path.resolve(strict=True)
    if path.is_symlink() or not resolved.is_file():
        raise ValueError(f"{label} must be a regular non-symlink file")
    raw = resolved.read_bytes()
    try:
        value = json.loads(raw)
    except (UnicodeDecodeError, json.JSONDecodeError) as error:
        raise ValueError(f"{label} must be valid UTF-8 JSON") from error
    if not isinstance(value, dict):
        raise ValueError(f"{label} root must be an object")
    return value, raw


def select_analyzer(profile_path: Path, document_class: str) -> tuple[str, str, str]:
    profile, profile_bytes = read_json_file(profile_path, "document profile")
    if profile.get("schema") != PROFILE_SCHEMA:
        raise ValueError("unsupported document profile schema")
    if profile.get("provider") != "Microsoft Azure Content Understanding":
        raise ValueError("document profile provider mismatch")
    if profile.get("sdk") != f"azure-ai-contentunderstanding=={SDK_VERSION}":
        raise ValueError("document profile SDK mismatch")
    if profile.get("api_version") != API_VERSION:
        raise ValueError("document profile API version mismatch")
    classes = profile.get("classes")
    if not isinstance(classes, list) or not classes:
        raise ValueError("document profile classes must be a non-empty array")
    by_id: dict[str, dict[str, Any]] = {}
    for item in classes:
        if not isinstance(item, dict) or not isinstance(item.get("id"), str) or not item["id"]:
            raise ValueError("document profile contains an invalid class")
        if item["id"] in by_id:
            raise ValueError("document profile contains duplicate class ids")
        by_id[item["id"]] = item
    selected = by_id.get(document_class)
    if selected is None:
        raise ValueError("document class is not declared by the profile")
    if selected.get("decision") != "REQUIRED":
        raise ValueError("document class is blocked by the profile")
    analyzer_id = selected.get("analyzer_id")
    schema_version = selected.get("schema_version")
    if not isinstance(analyzer_id, str) or not analyzer_id:
        raise ValueError("required document class has no analyzer")
    if not isinstance(schema_version, str) or not schema_version:
        raise ValueError("required document class has no schema version")
    if selected.get("automatic_storage") is not False:
        raise ValueError("document profile must explicitly deny automatic storage")
    return analyzer_id, schema_version, sha256_bytes(profile_bytes)


def verify_security_receipt(
    receipt_path: Path,
    *,
    payload_sha256: str,
    payload_bytes: int,
    mime_type: str,
) -> str:
    receipt, receipt_bytes = read_json_file(receipt_path, "security receipt")
    if receipt.get("schema") != SECURITY_RECEIPT_SCHEMA:
        raise ValueError("unsupported security receipt schema")
    if receipt.get("decision") != "ADMITTED":
        raise ValueError("security receipt did not admit the input")
    if receipt.get("business_storage_authorized") is not False:
        raise ValueError("security receipt must explicitly deny business storage")
    observed_input = receipt.get("input")
    if not isinstance(observed_input, dict):
        raise ValueError("security receipt input evidence is missing")
    if observed_input.get("sha256") != payload_sha256 or observed_input.get("bytes") != payload_bytes:
        raise ValueError("security receipt does not match input bytes")
    observed_type = receipt.get("content_type")
    if not isinstance(observed_type, dict) or observed_type.get("mime_type") != mime_type:
        raise ValueError("security receipt content type does not match input")
    return sha256_bytes(receipt_bytes)


def analyze_to_evidence(client: AnalyzerClient, request: AnalysisRequest) -> dict[str, Any]:
    if importlib.metadata.version("azure-ai-contentunderstanding") != SDK_VERSION:
        raise RuntimeError(f"azure-ai-contentunderstanding {SDK_VERSION} is required")
    if not request.document_class or request.max_bytes <= 0:
        raise ValueError("document_class and positive max_bytes are required")
    source = request.input_path.resolve(strict=True)
    if request.input_path.is_symlink() or not source.is_file():
        raise ValueError("input must be a regular non-symlink file")
    mime_type = ALLOWED_MIME_TYPES.get(source.suffix.lower())
    if mime_type is None:
        raise ValueError("unsupported document extension")
    payload = source.read_bytes()
    if not payload or len(payload) > request.max_bytes:
        raise ValueError("document byte size is outside the allowed range")
    payload_sha256 = sha256_bytes(payload)
    analyzer_id, schema_version, profile_sha256 = select_analyzer(request.profile_path, request.document_class)
    security_receipt_sha256 = verify_security_receipt(
        request.security_receipt_path,
        payload_sha256=payload_sha256,
        payload_bytes=len(payload),
        mime_type=mime_type,
    )
    output = request.output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".document-analysis-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        result = client.begin_analyze_binary(
            analyzer_id=analyzer_id,
            binary_input=payload,
            content_type=mime_type,
        ).result()
        raw = result.as_dict()
        observed_analyzer = raw.get("analyzerId") or raw.get("analyzer_id")
        if observed_analyzer and observed_analyzer != analyzer_id:
            raise RuntimeError("provider result analyzer does not match request")
        observed_api_version = raw.get("apiVersion") or raw.get("api_version")
        if observed_api_version and observed_api_version != API_VERSION:
            raise RuntimeError("provider result API version does not match GA profile")
        contents = raw.get("contents")
        if not isinstance(contents, list) or not contents:
            raise RuntimeError("provider result contains no analyzed content")
        result_bytes = (json.dumps(raw, ensure_ascii=False, sort_keys=True, indent=2, default=str) + "\n").encode("utf-8")
        (stage / "provider-result.json").write_bytes(result_bytes)
        receipt = {
            "schema": "elite-azure-content-understanding-analysis-receipt/v2",
            "created_at": datetime.now(timezone.utc).isoformat(),
            "provider": "Microsoft Azure Content Understanding",
            "sdk_version": SDK_VERSION,
            "sdk_api_version": API_VERSION,
            "provider_result_api_version": observed_api_version,
            "document_class": request.document_class,
            "analyzer_id": analyzer_id,
            "schema_version": schema_version,
            "document_profile_sha256": profile_sha256,
            "security_receipt_sha256": security_receipt_sha256,
            "input_filename": source.name,
            "input_bytes": len(payload),
            "input_sha256": payload_sha256,
            "content_type": mime_type,
            "content_count": len(contents),
            "provider_result_sha256": sha256_bytes(result_bytes),
            "automatic_storage_authorized": False,
        }
        (stage / "ANALYSIS_RECEIPT.json").write_text(
            json.dumps(receipt, sort_keys=True, indent=2) + "\n",
            encoding="utf-8",
            newline="\n",
        )
        stage.replace(output)
        return receipt
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--document-class", required=True)
    parser.add_argument("--input", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--security-receipt", required=True, type=Path)
    parser.add_argument("--endpoint-env", default="CONTENTUNDERSTANDING_ENDPOINT")
    parser.add_argument("--key-env", default="CONTENTUNDERSTANDING_KEY")
    parser.add_argument("--max-bytes", type=int, default=20 * 1024 * 1024)
    args = parser.parse_args()
    endpoint = os.environ.get(args.endpoint_env, "").strip()
    key = os.environ.get(args.key_env, "").strip()
    if not endpoint or not key:
        raise RuntimeError("endpoint/key environment references are required; secrets are never accepted as CLI values")
    client = ContentUnderstandingClient(endpoint=endpoint, credential=AzureKeyCredential(key))
    receipt = analyze_to_evidence(
        client,
        AnalysisRequest(
            args.document_class,
            args.input,
            args.output,
            args.profile,
            args.security_receipt,
            args.max_bytes,
        ),
    )
    print(json.dumps({"status": "AZURE_CONTENT_UNDERSTANDING_ANALYSIS_PASS", "receipt": receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `azure_document_runtime/test_content_understanding_runtime.py`
```yaml
block_id: "AZURE-CU-DOCUMENT-RUNTIME:test:v1"
operation: CREATE
provenance: AUTHORED
source: "local stdlib regression plus official SDK import/signature probe"
license: "LicenseRef-Workspace-Owner"
sha256: "45c00425a893a7b0bd132e5858a26e60cfb95f4a5890a956b5e071b55173ac73"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import hashlib
import json
from pathlib import Path
import tempfile
import unittest

from run_content_understanding import AnalysisRequest, analyze_to_evidence


class FakeResult:
    def __init__(self, body): self.body = body
    def as_dict(self): return self.body


class FakePoller:
    def __init__(self, body): self.body = body
    def result(self): return FakeResult(self.body)


class FakeClient:
    def __init__(self, body): self.body = body; self.calls = []
    def begin_analyze_binary(self, analyzer_id, binary_input, *, content_type):
        self.calls.append((analyzer_id, binary_input, content_type))
        return FakePoller(self.body)


class RuntimeTests(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory()
        self.root = Path(self.temporary.name)
        self.source = self.root / "invoice.pdf"
        self.source.write_bytes(b"%PDF-1.7 fixture")
        self.profile = Path(__file__).parent / "document-profile.template.json"
        self.security = self.root / "security-receipt.json"
        self.write_security_receipt()

    def tearDown(self):
        self.temporary.cleanup()

    def write_security_receipt(self, *, decision="ADMITTED", sha256=None, mime_type="application/pdf"):
        payload = self.source.read_bytes()
        receipt = {
            "schema": "elite-secure-local-file-receipt/v1",
            "decision": decision,
            "reason": "ALL_SELECTED_GATES_PASSED" if decision == "ADMITTED" else "TEST_REJECTION",
            "business_storage_authorized": False,
            "input": {
                "sha256": sha256 or hashlib.sha256(payload).hexdigest(),
                "bytes": len(payload),
            },
            "content_type": {"mime_type": mime_type, "label": "pdf", "score": 1.0},
        }
        self.security.write_text(json.dumps(receipt) + "\n", encoding="utf-8", newline="\n")

    def request(self, output, *, document_class="supplier-invoice", profile=None, security=None):
        return AnalysisRequest(
            document_class,
            self.source,
            output,
            profile or self.profile,
            security or self.security,
        )

    def test_official_sdk_contract_and_evidence_chain(self):
        client = FakeClient({
            "analyzerId": "prebuilt-invoice",
            "apiVersion": "2025-11-01",
            "contents": [{"kind": "document", "fields": {"InvoiceId": {
                "type": "string", "valueString": "INV-1", "confidence": 0.99, "source": "D(1,1)"
            }}}],
        })
        output = self.root / "evidence"
        receipt = analyze_to_evidence(client, self.request(output))
        self.assertFalse(receipt["automatic_storage_authorized"])
        self.assertEqual(receipt["document_class"], "supplier-invoice")
        self.assertEqual(receipt["analyzer_id"], "prebuilt-invoice")
        self.assertEqual(receipt["sdk_api_version"], "2025-11-01")
        self.assertRegex(receipt["security_receipt_sha256"], r"^[0-9a-f]{64}$")
        self.assertRegex(receipt["document_profile_sha256"], r"^[0-9a-f]{64}$")
        self.assertEqual(client.calls[0][2], "application/pdf")
        self.assertTrue((output / "provider-result.json").is_file())
        self.assertTrue((output / "ANALYSIS_RECEIPT.json").is_file())

    def test_provider_failures_are_atomic(self):
        for index, (body, needle) in enumerate([
            ({}, "no analyzed content"),
            ({"analyzerId": "other", "contents": [{}]}, "does not match request"),
            ({"apiVersion": "2025-05-01-preview", "contents": [{}]}, "does not match GA profile"),
        ]):
            output = self.root / f"failed-{index}"
            with self.assertRaisesRegex(RuntimeError, needle):
                analyze_to_evidence(FakeClient(body), self.request(output))
            self.assertFalse(output.exists())
        occupied = self.root / "occupied"
        occupied.mkdir()
        with self.assertRaises(FileExistsError):
            analyze_to_evidence(FakeClient({"contents": [{}]}), self.request(occupied))

    def test_profile_is_execution_authority(self):
        client = FakeClient({"contents": [{}]})
        with self.assertRaisesRegex(ValueError, "blocked by the profile"):
            analyze_to_evidence(client, self.request(self.root / "blocked", document_class="packing-list"))
        with self.assertRaisesRegex(ValueError, "not declared"):
            analyze_to_evidence(client, self.request(self.root / "unknown", document_class="invented-class"))
        self.assertEqual(client.calls, [])

    def test_security_receipt_must_match_and_admit(self):
        client = FakeClient({"contents": [{}]})
        self.write_security_receipt(decision="REJECTED")
        with self.assertRaisesRegex(ValueError, "did not admit"):
            analyze_to_evidence(client, self.request(self.root / "rejected"))
        self.write_security_receipt(sha256="0" * 64)
        with self.assertRaisesRegex(ValueError, "does not match input bytes"):
            analyze_to_evidence(client, self.request(self.root / "mismatch"))
        self.write_security_receipt(mime_type="image/png")
        with self.assertRaisesRegex(ValueError, "content type does not match"):
            analyze_to_evidence(client, self.request(self.root / "wrong-type"))
        self.assertEqual(client.calls, [])

    def test_extension_and_output_fail_closed(self):
        unsupported = self.root / "invoice.exe"
        unsupported.write_bytes(b"MZ")
        request = AnalysisRequest(
            "supplier-invoice",
            unsupported,
            self.root / "never",
            self.profile,
            self.security,
        )
        with self.assertRaisesRegex(ValueError, "unsupported"):
            analyze_to_evidence(FakeClient({"contents": [{}]}), request)


if __name__ == "__main__":
    unittest.main()
````

### FILE: `azure_document_runtime/README.md`
```yaml
block_id: "AZURE-CU-DOCUMENT-RUNTIME:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "e566bbdff1930823d13adf4076557d46d774f7566e4c8d53721fc3ef81e0cc50"
variables: []
secrets_allowed: false
```
````markdown
# Azure Content Understanding document runtime

This runtime calls Microsoft Azure Content Understanding through the exact GA Python SDK 1.1.0. It submits local PDF/image bytes with `begin_analyze_binary`, preserves the complete official result, hashes the original and result, records analyzer/API versions and writes atomically to a new evidence directory.

Endpoint and key are read only from named environment variables; secret values are never CLI parameters or receipts. The caller selects a document class, never an arbitrary analyzer ID. The runtime derives the analyzer only from a `REQUIRED` class in the exact profile, records the profile hash and rejects every `BLOCKED_*` class. The template maps official prebuilt invoice and purchase-order analyzers and leaves proformas, packing lists, bills of lading, customs and other custom classes blocked until a versioned custom analyzer and project corpus are supplied.

Every call also requires `--security-receipt` from `SECURE-LOCAL-FILE-INGESTION-GATE`. Its schema and `ADMITTED` decision must match the input bytes, SHA-256 and MIME type exactly; the analysis receipt preserves the security receipt hash. A filename extension alone is never the trust decision.

The output is evidence, not an authorization to persist business fields. The project must evaluate each class/field against approved ground truth, configure review/rejection and explicitly promote its storage policy.
````

## 6. Configuration surface

| Parámetro | Tipo/default | Validación | Secreto | Efecto |
|---|---|---|---|---|
| document class | string | debe existir y estar `REQUIRED` | no | selecciona analyzer sólo desde profile |
| document profile | JSON path | schema/provider/SDK/API/clases únicas; auto-storage false | no | fija analyzer/schema y queda hash-linked |
| security receipt | JSON path | `ADMITTED`; SHA/bytes/MIME coinciden; storage false | contiene evidencia | une cuarentena/scan con análisis |
| input | path local | regular, no symlink, MIME allowlist, size | contiene datos sensibles | bytes enviados al servicio |
| output | path | debe no existir | no | evidence-only result/receipt |
| endpoint env name | string / `CONTENTUNDERSTANDING_ENDPOINT` | value no vacío | endpoint no | resource seleccionado |
| key env name | string / `CONTENTUNDERSTANDING_KEY` | value no vacío | sí; value nunca CLI/receipt | autentica SDK |
| max bytes | int / 20 MiB | positivo | no | preflight de capacidad |

## 7. Dependency bill

| Package/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente |
|---|---|---|---|---|---|
| Python | >=3.9 upstream; 3.12.13 en evidencia | runtime/test | PSF-2.0 | runtime | python.org |
| `azure-ai-contentunderstanding` | 1.1.0 wheel SHA fijado | cliente/model/LRO | MIT | runtime | Microsoft PyPI |
| Azure Content Understanding API | 2025-11-01 observado | analysis | términos Azure | external | Microsoft Learn |

El archivo direct fija el top-level y `requirements-windows-py312.lock` fija los nueve wheels resueltos para CPython 3.12/Windows x86-64 con un hash admitido por distribución. Otro OS/arquitectura/Python debe generar, licenciar, escanear y admitir su propio lock; nunca reutiliza éste silenciosamente.

## 8. Apply order

1. Componer adquisición, file-security, artifact core, runtime y evaluación estricta; seleccionar Azure en el expediente.
2. Adquirir wheel con approval hash-linked; generar venv/lock transitivo.
3. Completar profile por todas las clases requeridas; crear/versionar custom analyzers faltantes.
4. Ejecutar el gate local y entregar su receipt exacto al runtime; ninguna extensión sustituye esta evidencia.
5. Probar imports/tests locales y endpoint sandbox por secret references.
6. Ejecutar corpus holdout; conservar original/security receipt/profile/result/analysis receipt/ground truth.
7. Sólo después decidir field policies/review/storage en el expediente del proyecto.

En proyecto existente, ejecutar sobre copia/corpus aislado y no sobrescribir parsers/resultados. Rollback conserva evidencia y revierte routing/model/schema.

## 9. Verification

```powershell
python -m unittest -v test_content_understanding_runtime.py
python -c "import importlib.metadata as m; assert m.version('azure-ai-contentunderstanding') == '1.1.0'"
```

Esperado: cinco tests PASS + contract probe. Positivo verifica método oficial, MIME, contents y cadena hash input/security/profile/result/receipt. Negativos cubren clase bloqueada/desconocida, security receipt rechazado o divergente, analyzer/API provider divergentes, contenidos vacíos, output ocupado y formato no permitido sin dejar destino. La llamada sandbox real queda bloqueada hasta endpoint/key/analyzer/corpus autorizados.

## 10. Reconstruction evidence

- materialización limpia: seis archivos con SHA/manifest verificados;
- Python 3.12.13 + wheel oficial 1.1.0;
- tests oficiales aislables: 183 PASS; suite integral desde archive condicionada por `.git`/test-proxy;
- `unittest` del wrapper: 5 PASS;
- contract probe: firma `begin_analyze_binary(analyzer_id,binary_input,content_type)` PASS;
- ningún secreto/red/servicio utilizado; output fixture no autoriza storage;
- fallo `LIB-FAIL-039` conserva la selección de venv incompleto;
- fecha/revisor: 2026-08-27 / Codex; expedientes `AZURE_CONTENT_UNDERSTANDING_DOCUMENT_RUNTIME_2026-08-26_V1.md` y `MICROSOFT_AZURE_CONTENT_UNDERSTANDING_PYTHON_2026-08-27_V1.md`.
