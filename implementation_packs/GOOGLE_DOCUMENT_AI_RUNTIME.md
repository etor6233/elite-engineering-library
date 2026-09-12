# Google Cloud Document AI Runtime

## 1. Metadata

```yaml
pack_id: "GOOGLE-DOCUMENT-AI-RUNTIME"
pack_version: "0.1.3"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un runtime Python AUTHORED sobre los SDK oficiales Google Cloud Document AI 3.15.0 y Cloud Storage 3.13.1: processor version gobernado sólo por perfil, receipt obligatoria del gate local, evidencia hash-linked, commit atómico y lock exacto de 23 wheels para CPython 3.12/Windows x86-64."
stacks: ["CPython 3.12 / Windows x86-64", "google-cloud-documentai 3.15.0", "google-cloud-storage 3.13.1", "Google Cloud Document AI v1"]
compatible_with: ["OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE 0.3.x", "OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x", "STRICT-DOCUMENT-FIELD-EVALUATION-GATE 0.1.x"]
incompatible_with: ["processor default mutable", "processor por CLI", "receipt de seguridad ausente o no vinculada", "auto-storage sin corpus", "archivo remoto no controlado"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://github.com/googleapis/google-cloud-python/tree/google-cloud-documentai-v3.15.0/packages/google-cloud-documentai", "https://pypi.org/project/google-cloud-documentai/3.15.0/", "https://pypi.org/project/google-cloud-storage/3.13.1/", "https://docs.cloud.google.com/python/docs/reference/documentai/3.15.0/google.cloud.documentai_v1.services.document_processor_service.DocumentProcessorServiceClient"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use después de seleccionar Google Cloud Document AI, crear un processor en un proyecto/región autorizados y adquirir el wheel oficial exacto. Acepta PDF/PNG/JPEG/TIFF sólo después de una receipt `ADMITTED` del gate local, exige el recurso completo de una `processorVersion` dentro del perfil aprobado y usa Application Default Credentials. Rechazar cuando falten cuenta, facturación aprobada, IAM mínimo, región, versión desplegada, corpus, evaluación o cuando se pretenda almacenar automáticamente la respuesta como hechos de negocio.

El runner es glue `AUTHORED` mínimo alrededor de `RawDocument`, `ProcessRequest` y `DocumentProcessorServiceClient.process_document`; el SDK, request, response y modelos son Google. El template no activa ningún processor: `supplier-invoice` queda `CONFIGURATION_REQUIRED` hasta que el proyecto aporte su recurso/schema/evidencia, y todas las demás clases continúan bloqueadas en rutas oficiales de Custom Extractor/Classifier/Splitter.

## 3. Architecture contract

El runner verifica la distribución 3.15.0, clase `REQUIRED` única, provider/SDK/schema del perfil, recurso versionado exacto, archivo regular no symlink, MIME y 1..20 MiB. Exige que SHA-256/bytes/MIME coincidan con una receipt `ADMITTED` completa de ClamAV/Magika/YARA-X; el CLI no admite un processor libre. Envía bytes inline al endpoint regional derivado, conserva todo `ProcessResponse` mediante la serialización oficial proto-plus y produce receipt v2 enlazando perfil, seguridad, input, output, processor y schema. No acepta archivos de credenciales ni tokens por argumentos; el cliente usa ADC.

Cualquier fallo de validación, SDK, credencial, servicio, respuesta o output elimina staging y no crea destino final. La respuesta mantiene texto, layout, entidades, confidence, anchors, normalization y estado de revisión tal como Google los devuelve; este pack no transforma predicciones en hechos de negocio. Latencia, precisión, costo y cuotas se miden en sandbox con corpus representativo. Rollback fija la versión anterior y reprocesa sólo desde originales/evidencia autorizados.

## 4. Exact file manifest

```text
CREATE google_document_runtime/requirements-direct.in
CREATE google_document_runtime/requirements-windows-py312.lock
CREATE google_document_runtime/document-profile.template.json
CREATE google_document_runtime/run_document_ai.py
CREATE google_document_runtime/test_document_ai_runtime.py
CREATE google_document_runtime/README.md
```

## 5. Materialization blocks

### FILE: `google_document_runtime/requirements-direct.in`
```yaml
block_id: "GOOGLE-DOCUMENT-AI-RUNTIME:requirement:v1"
operation: CREATE
provenance: AUTHORED
source: "exact official PyPI wheel URL and SHA-256"
license: "LicenseRef-Workspace-Owner"
sha256: "92391f41cea2e572e91a95857b32d310bf00024d4277e1eb7401ae7e7de12e59"
variables: []
secrets_allowed: false
```
````text
google-cloud-documentai @ https://files.pythonhosted.org/packages/55/d1/2a873f97cb08bb592f5b944e39f040161d5f7d2d4edbfae4c7515869b12a/google_cloud_documentai-3.15.0-py3-none-any.whl#sha256=f040f4f9db43411184197a808b11fde52b580723a9fca336f43ea7c6a885bfd8
google-cloud-storage @ https://files.pythonhosted.org/packages/06/6f/d69f0e185e08ddb58c323a0a935af2b492907b5de362bc08933b0a3b5644/google_cloud_storage-3.13.1-py3-none-any.whl#sha256=98208de6c21e85cecd3eb44551894efff33d98365500e178867d4305854a770a
````

### FILE: `google_document_runtime/requirements-windows-py312.lock`
```yaml
block_id: "GOOGLE-DOCUMENT-AI-RUNTIME:windows-py312-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "dated target-specific resolution of official PyPI artifacts; every wheel hash was recomputed locally"
license: "LicenseRef-Workspace-Owner"
sha256: "c4d503af18cb0fd5bcccd4cde58d7baf79db61f6af12ce9d87a4c2ec7c50aee0"
variables: []
secrets_allowed: false
```
````text
# Target: CPython 3.12 / Windows x86-64. Resolved 2026-08-28; install with --require-hashes.
certifi==2026.7.22 --hash=sha256:62f22742b58a1a33014a2b6b706588a8d7e2a88ae7bd1a6ebe8c992928483775
cffi==2.1.1 --hash=sha256:f53e442b08449d42821fa4a4fba000095af9f62742a500f978a9f557ec44339a
charset-normalizer==3.5.1 --hash=sha256:3617ac3cfd8b9888f145ad89dd6e692285834b0201c6074a5eeaad3fd4d668c2
cryptography==50.0.1 --hash=sha256:aed8db4f6d71c51efb89530e12d9464e7bf2923d46c3205dc794a2a93f8c0648
google-api-core==2.34.0 --hash=sha256:cdf9c67e7ca2402d86ccbfde5f2503fc83e3cc3f58cc78456ae96cad24a6d2de
google-auth==2.57.0 --hash=sha256:180dafe015cfb62193bea26b677500fab5b9fd51a1e825ebf3ad9b182047ae59
google-cloud-core==2.7.0 --hash=sha256:c18a250904cfdda021eb3ae8b8238c9f9ca272a4cbbfb5cba946b3fe3022eed1
google-cloud-documentai==3.15.0 --hash=sha256:f040f4f9db43411184197a808b11fde52b580723a9fca336f43ea7c6a885bfd8
google-cloud-storage==3.13.1 --hash=sha256:98208de6c21e85cecd3eb44551894efff33d98365500e178867d4305854a770a
google-crc32c==1.8.0 --hash=sha256:3cc0c8912038065eafa603b238abf252e204accab2a704c63b9e14837a854962
google-resumable-media==2.10.2 --hash=sha256:e3cedc827a4ea41e216582d74346f1fb9fceb625a8c3c53912f2ca1d663334d7
googleapis-common-protos==1.75.2 --hash=sha256:6b83302f554ea93a0f48409c7fc2050f954bcbcddb7e3a9c76d4a823cb22920e
grpcio==1.83.0 --hash=sha256:1c699bbb20f143c8f2bff219de578aa2dc1f919399d67dc702b038b986ee62df
grpcio-status==1.83.0 --hash=sha256:f6a838a7c5fb84ae98833ec0ef81ed438c26e11e54b2ddb8e92ad328c861de69
idna==3.19 --hash=sha256:815e7be7a7806d54abb586dc943addc79e8b2ee16915059658cbeff4b1b43bf4
proto-plus==1.28.4 --hash=sha256:4b01341272f8a348db3f003b6143109f83ab43091019d5181b3fcdf500ab32aa
protobuf==7.36.0 --hash=sha256:1781cc1de61249b750848029bca452c0a8b7e990080316b9bbc2518b2117b488
pyasn1==0.6.4 --hash=sha256:deda9277cfd454080ec40b207fb6df82206a3a2688735233cdcd8d3d565f088b
pyasn1_modules==0.4.2 --hash=sha256:29253a9207ce32b64c3ac6600edc75368f98473906e8fd1043bd6b5b1de2c14a
pycparser==3.0 --hash=sha256:b727414169a36b7d524c1c3e31839a521725078d7b2ff038656844266160a992
requests==2.34.2 --hash=sha256:2a0d60c172f83ac6ab31e4554906c0f3b3588d37b5cb939b1c061f4907e278e0
typing_extensions==4.16.0 --hash=sha256:481caa481374e813c1b176ada14e97f1f67a4539ce9cfeb3f350d78d6370c2e8
urllib3==2.7.0 --hash=sha256:9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897
````

### FILE: `google_document_runtime/document-profile.template.json`
```yaml
block_id: "GOOGLE-DOCUMENT-AI-RUNTIME:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed profile using Google documented processor types"
license: "LicenseRef-Workspace-Owner"
sha256: "7e2322bcc52779e6840369f8d3638f0df4a66471f36434be6327b7a9974f4d5c"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-google-document-ai-profile/v2",
  "provider": "Google Cloud Document AI",
  "sdk": "google-cloud-documentai==3.15.0",
  "classes": [
    {"id": "supplier-invoice", "decision": "CONFIGURATION_REQUIRED", "processor_type": "INVOICE_PROCESSOR", "processor_version_resource": "", "schema_version": "", "automatic_storage": false},
    {"id": "commercial-invoice", "decision": "BLOCKED_EVALUATION_REQUIRED", "processor_type": "INVOICE_PROCESSOR", "processor_version_resource": "", "schema_version": "", "automatic_storage": false},
    {"id": "proforma-invoice", "decision": "BLOCKED_CUSTOM_EXTRACTOR_REQUIRED", "processor_type": "CUSTOM_EXTRACTION_PROCESSOR", "processor_version_resource": "", "schema_version": "", "automatic_storage": false},
    {"id": "purchase-order", "decision": "BLOCKED_CUSTOM_EXTRACTOR_REQUIRED", "processor_type": "CUSTOM_EXTRACTION_PROCESSOR", "processor_version_resource": "", "schema_version": "", "automatic_storage": false},
    {"id": "packing-list", "decision": "BLOCKED_CUSTOM_EXTRACTOR_REQUIRED", "processor_type": "CUSTOM_EXTRACTION_PROCESSOR", "processor_version_resource": "", "schema_version": "", "automatic_storage": false},
    {"id": "bill-of-lading", "decision": "BLOCKED_CUSTOM_EXTRACTOR_REQUIRED", "processor_type": "CUSTOM_EXTRACTION_PROCESSOR", "processor_version_resource": "", "schema_version": "", "automatic_storage": false},
    {"id": "delivery-note", "decision": "BLOCKED_CUSTOM_EXTRACTOR_REQUIRED", "processor_type": "CUSTOM_EXTRACTION_PROCESSOR", "processor_version_resource": "", "schema_version": "", "automatic_storage": false},
    {"id": "certificate-of-origin", "decision": "BLOCKED_CUSTOM_EXTRACTOR_REQUIRED", "processor_type": "CUSTOM_EXTRACTION_PROCESSOR", "processor_version_resource": "", "schema_version": "", "automatic_storage": false},
    {"id": "quality-inspection-certificate", "decision": "BLOCKED_CUSTOM_EXTRACTOR_REQUIRED", "processor_type": "CUSTOM_EXTRACTION_PROCESSOR", "processor_version_resource": "", "schema_version": "", "automatic_storage": false},
    {"id": "customs-declaration", "decision": "BLOCKED_CUSTOM_EXTRACTOR_REQUIRED", "processor_type": "CUSTOM_EXTRACTION_PROCESSOR", "processor_version_resource": "", "schema_version": "", "automatic_storage": false},
    {"id": "mixed-document-batch", "decision": "BLOCKED_CUSTOM_CLASSIFIER_OR_SPLITTER_REQUIRED", "processor_type": "CUSTOM_CLASSIFICATION_PROCESSOR", "processor_version_resource": "", "schema_version": "", "automatic_storage": false}
  ]
}
````

### FILE: `google_document_runtime/run_document_ai.py`
```yaml
block_id: "GOOGLE-DOCUMENT-AI-RUNTIME:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local wrapper invoking Google official SDK public API"
license: "LicenseRef-Workspace-Owner"
sha256: "2ec44f12873ba5fc0e16bd6361a3e23ccfc204806b6a0b4441e3f7e90038b106"
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
from pathlib import Path
import shutil
import sys
import uuid
from typing import Any, Protocol

from google.cloud import documentai_v1 as documentai


SDK_VERSION = "3.15.0"
PROFILE_SCHEMA = "elite-google-document-ai-profile/v2"
SECURITY_RECEIPT_SCHEMA = "elite-secure-local-file-receipt/v1"
ALLOWED_MIME_TYPES = {
    ".pdf": "application/pdf",
    ".png": "image/png",
    ".jpg": "image/jpeg",
    ".jpeg": "image/jpeg",
    ".tif": "image/tiff",
    ".tiff": "image/tiff",
}


class ProcessorClient(Protocol):
    def process_document(
        self, request: documentai.ProcessRequest
    ) -> documentai.ProcessResponse: ...


@dataclass(frozen=True)
class AnalysisRequest:
    document_class: str
    input_path: Path
    output_directory: Path
    profile_path: Path
    security_receipt_path: Path
    max_bytes: int = 20 * 1024 * 1024


@dataclass(frozen=True)
class ProcessorSelection:
    document_class: str
    processor_type: str
    processor_version_resource: str
    schema_version: str
    project: str
    location: str
    processor: str
    processor_version: str
    profile_sha256: str


def sha256_bytes(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def is_sha256(value: object) -> bool:
    return isinstance(value, str) and len(value) == 64 and all(
        character in "0123456789abcdef" for character in value
    )


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


def parse_processor_version_resource(name: str) -> dict[str, str]:
    segments = name.split("/")
    if (
        len(segments) != 8
        or segments[0] != "projects"
        or segments[2] != "locations"
        or segments[4] != "processors"
        or segments[6] != "processorVersions"
    ):
        raise ValueError("an exact processor version resource is required")
    values = {
        "project": segments[1],
        "location": segments[3],
        "processor": segments[5],
        "processor_version": segments[7],
    }
    if any(
        not value
        or value.strip() != value
        or value in {".", ".."}
        or any(character.isspace() for character in value)
        for value in values.values()
    ):
        raise ValueError("processor version resource contains an empty or unsafe segment")
    return values


def select_processor(profile_path: Path, document_class: str) -> ProcessorSelection:
    profile, profile_bytes = read_json_file(profile_path, "document profile")
    if profile.get("schema") != PROFILE_SCHEMA:
        raise ValueError("unsupported document profile schema")
    if profile.get("provider") != "Google Cloud Document AI":
        raise ValueError("document profile provider mismatch")
    if profile.get("sdk") != f"google-cloud-documentai=={SDK_VERSION}":
        raise ValueError("document profile SDK mismatch")
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
        raise ValueError("document class is blocked or not configured by the profile")
    if selected.get("automatic_storage") is not False:
        raise ValueError("document profile must explicitly deny automatic storage")
    processor_type = selected.get("processor_type")
    resource_name = selected.get("processor_version_resource")
    schema_version = selected.get("schema_version")
    if not isinstance(processor_type, str) or not processor_type:
        raise ValueError("required document class has no processor type")
    if not isinstance(resource_name, str) or not resource_name:
        raise ValueError("required document class has no processor version resource")
    if not isinstance(schema_version, str) or not schema_version:
        raise ValueError("required document class has no schema version")
    resource = parse_processor_version_resource(resource_name)
    return ProcessorSelection(
        document_class=document_class,
        processor_type=processor_type,
        processor_version_resource=resource_name,
        schema_version=schema_version,
        project=resource["project"],
        location=resource["location"],
        processor=resource["processor"],
        processor_version=resource["processor_version"],
        profile_sha256=sha256_bytes(profile_bytes),
    )


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
    if receipt.get("decision") != "ADMITTED" or receipt.get("reason") != "ALL_SELECTED_GATES_PASSED":
        raise ValueError("security receipt did not admit the input")
    if receipt.get("business_storage_authorized") is not False:
        raise ValueError("security receipt must explicitly deny business storage")
    if not isinstance(receipt.get("approval_id"), str) or not receipt["approval_id"]:
        raise ValueError("security receipt approval id is missing")
    if not is_sha256(receipt.get("policy_sha256")):
        raise ValueError("security receipt policy hash is invalid")
    observed_input = receipt.get("input")
    if not isinstance(observed_input, dict):
        raise ValueError("security receipt input evidence is missing")
    if observed_input.get("sha256") != payload_sha256 or observed_input.get("bytes") != payload_bytes:
        raise ValueError("security receipt does not match input bytes")
    observed_type = receipt.get("content_type")
    if not isinstance(observed_type, dict) or observed_type.get("mime_type") != mime_type:
        raise ValueError("security receipt content type does not match input")
    tools = receipt.get("tools")
    if not isinstance(tools, dict):
        raise ValueError("security receipt tool evidence is missing")
    for name in ("clamscan", "magika", "yara_x"):
        tool = tools.get(name)
        if (
            not isinstance(tool, dict)
            or not is_sha256(tool.get("sha256"))
            or not isinstance(tool.get("version"), str)
            or not tool["version"]
        ):
            raise ValueError(f"security receipt {name} evidence is invalid")
    clamav = receipt.get("clamav")
    if not isinstance(clamav, dict) or clamav.get("exit_code") != 0 or not is_sha256(clamav.get("output_sha256")):
        raise ValueError("security receipt ClamAV result is invalid")
    yara = receipt.get("yara_x")
    if (
        not isinstance(yara, dict)
        or not is_sha256(yara.get("compiled_rules_sha256"))
        or yara.get("matching_rules") != []
    ):
        raise ValueError("security receipt YARA-X result is invalid")
    return sha256_bytes(receipt_bytes)


def analyze_to_evidence(client: ProcessorClient, request: AnalysisRequest) -> dict[str, Any]:
    if importlib.metadata.version("google-cloud-documentai") != SDK_VERSION:
        raise RuntimeError(f"google-cloud-documentai {SDK_VERSION} is required")
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
    selection = select_processor(request.profile_path, request.document_class)
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
        api_request = documentai.ProcessRequest(
            name=selection.processor_version_resource,
            raw_document=documentai.RawDocument(content=payload, mime_type=mime_type),
        )
        response = client.process_document(request=api_request)
        if not isinstance(response, documentai.ProcessResponse):
            raise RuntimeError("official SDK ProcessResponse is required")
        document = response.document
        if not document or not (document.text or document.pages or document.entities):
            raise RuntimeError("provider response contains no document content")
        raw = documentai.ProcessResponse.to_dict(
            response, preserving_proto_field_name=True
        )
        result_bytes = (
            json.dumps(raw, ensure_ascii=False, sort_keys=True, indent=2, default=str) + "\n"
        ).encode("utf-8")
        (stage / "provider-response.json").write_bytes(result_bytes)
        receipt = {
            "schema": "elite-google-document-ai-analysis-receipt/v2",
            "created_at": datetime.now(timezone.utc).isoformat(),
            "provider": "Google Cloud Document AI",
            "sdk_version": SDK_VERSION,
            "document_class": selection.document_class,
            "processor_type": selection.processor_type,
            "processor_version_resource": selection.processor_version_resource,
            "processor_schema_version": selection.schema_version,
            "project": selection.project,
            "location": selection.location,
            "processor": selection.processor,
            "processor_version": selection.processor_version,
            "api_endpoint": f"{selection.location}-documentai.googleapis.com",
            "document_profile_sha256": selection.profile_sha256,
            "security_receipt_sha256": security_receipt_sha256,
            "input_filename": source.name,
            "input_bytes": len(payload),
            "input_sha256": payload_sha256,
            "content_type": mime_type,
            "page_count": len(document.pages),
            "entity_count": len(document.entities),
            "provider_response_sha256": sha256_bytes(result_bytes),
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


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser()
    parser.add_argument("--document-class", required=True)
    parser.add_argument("--input", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--security-receipt", required=True, type=Path)
    parser.add_argument("--max-bytes", type=int, default=20 * 1024 * 1024)
    return parser


def main() -> int:
    args = build_parser().parse_args()
    selection = select_processor(args.profile, args.document_class)
    client = documentai.DocumentProcessorServiceClient(
        client_options={
            "api_endpoint": f"{selection.location}-documentai.googleapis.com"
        }
    )
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
    print(
        json.dumps(
            {"status": "GOOGLE_DOCUMENT_AI_ANALYSIS_PASS", "receipt": receipt},
            sort_keys=True,
        )
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `google_document_runtime/test_document_ai_runtime.py`
```yaml
block_id: "GOOGLE-DOCUMENT-AI-RUNTIME:test:v1"
operation: CREATE
provenance: AUTHORED
source: "local stdlib regression plus official SDK import/signature probe"
license: "LicenseRef-Workspace-Owner"
sha256: "b378ddb00600bf7adf0388bdb3ed6e4443853744fbd77367f9b257f415ef3dd3"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import hashlib
import importlib.metadata
import json
from pathlib import Path
import tempfile
import unittest

from google.cloud import documentai_v1 as documentai

from run_document_ai import (
    AnalysisRequest,
    analyze_to_evidence,
    build_parser,
    parse_processor_version_resource,
)


PROCESSOR_VERSION = (
    "projects/project-1/locations/us/processors/invoice-1/"
    "processorVersions/pretrained-invoice-v2.1-2023-12-15"
)
HEX = "a" * 64


class FakeClient:
    def __init__(self, response):
        self.response = response
        self.calls = []

    def process_document(self, request):
        self.calls.append(request)
        return self.response


class RuntimeTests(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory()
        self.root = Path(self.temporary.name)
        self.source = self.root / "invoice.pdf"
        self.source.write_bytes(b"%PDF-1.7\nfixture")
        self.profile = self.root / "document-profile.json"
        self.security = self.root / "security-receipt.json"
        self.write_profile()
        self.write_security_receipt()

    def tearDown(self):
        self.temporary.cleanup()

    def write_profile(self, *, decision="REQUIRED", provider="Google Cloud Document AI", duplicate=False):
        selected = {
            "id": "supplier-invoice",
            "decision": decision,
            "processor_type": "INVOICE_PROCESSOR",
            "processor_version_resource": PROCESSOR_VERSION,
            "schema_version": "supplier-invoice/v1",
            "automatic_storage": False,
        }
        classes = [selected]
        if duplicate:
            classes.append(dict(selected))
        profile = {
            "schema": "elite-google-document-ai-profile/v2",
            "provider": provider,
            "sdk": "google-cloud-documentai==3.15.0",
            "classes": classes,
        }
        self.profile.write_text(json.dumps(profile) + "\n", encoding="utf-8", newline="\n")

    def write_security_receipt(
        self,
        *,
        decision="ADMITTED",
        sha256=None,
        mime_type="application/pdf",
        include_tools=True,
    ):
        payload = self.source.read_bytes()
        receipt = {
            "schema": "elite-secure-local-file-receipt/v1",
            "decision": decision,
            "reason": "ALL_SELECTED_GATES_PASSED" if decision == "ADMITTED" else "TEST_REJECTION",
            "approval_id": "security-approval-test",
            "policy_sha256": HEX,
            "input": {
                "sha256": sha256 or hashlib.sha256(payload).hexdigest(),
                "bytes": len(payload),
            },
            "content_type": {"mime_type": mime_type, "label": "pdf", "score": 1.0},
            "business_storage_authorized": False,
            "clamav": {"exit_code": 0, "output_sha256": HEX},
            "yara_x": {"compiled_rules_sha256": HEX, "matching_rules": []},
        }
        if include_tools:
            receipt["tools"] = {
                "clamscan": {"sha256": HEX, "version": "ClamAV 1.5.1"},
                "magika": {"sha256": HEX, "version": "magika 1.0.1"},
                "yara_x": {"sha256": HEX, "version": "1.20.0"},
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

    def response(self):
        return documentai.ProcessResponse(
            document=documentai.Document(
                text="Invoice INV-1",
                pages=[documentai.Document.Page(page_number=1)],
                entities=[
                    documentai.Document.Entity(
                        type_="invoice_id", mention_text="INV-1", confidence=0.99
                    )
                ],
            )
        )

    def test_official_sdk_contract_and_hash_linked_evidence(self):
        self.assertEqual(importlib.metadata.version("google-cloud-documentai"), "3.15.0")
        client = FakeClient(self.response())
        output = self.root / "evidence"
        receipt = analyze_to_evidence(client, self.request(output))
        self.assertEqual(receipt["schema"], "elite-google-document-ai-analysis-receipt/v2")
        self.assertEqual(receipt["document_class"], "supplier-invoice")
        self.assertEqual(receipt["processor_type"], "INVOICE_PROCESSOR")
        self.assertEqual(receipt["processor_version_resource"], PROCESSOR_VERSION)
        self.assertEqual(receipt["processor_schema_version"], "supplier-invoice/v1")
        self.assertFalse(receipt["automatic_storage_authorized"])
        self.assertRegex(receipt["document_profile_sha256"], r"^[0-9a-f]{64}$")
        self.assertRegex(receipt["security_receipt_sha256"], r"^[0-9a-f]{64}$")
        self.assertEqual(receipt["page_count"], 1)
        self.assertEqual(receipt["entity_count"], 1)
        self.assertEqual(len(client.calls), 1)
        api_request = client.calls[0]
        self.assertIsInstance(api_request, documentai.ProcessRequest)
        self.assertEqual(api_request.name, PROCESSOR_VERSION)
        self.assertEqual(api_request.raw_document.content, self.source.read_bytes())
        self.assertEqual(api_request.raw_document.mime_type, "application/pdf")
        self.assertTrue((output / "provider-response.json").is_file())
        self.assertTrue((output / "ANALYSIS_RECEIPT.json").is_file())

    def test_profile_is_execution_authority(self):
        client = FakeClient(self.response())
        template = Path(__file__).parent / "document-profile.template.json"
        with self.assertRaisesRegex(ValueError, "blocked or not configured"):
            analyze_to_evidence(client, self.request(self.root / "template", profile=template))
        self.write_profile(decision="BLOCKED_EVALUATION_REQUIRED")
        with self.assertRaisesRegex(ValueError, "blocked or not configured"):
            analyze_to_evidence(client, self.request(self.root / "blocked"))
        self.write_profile()
        with self.assertRaisesRegex(ValueError, "not declared"):
            analyze_to_evidence(client, self.request(self.root / "unknown", document_class="invented"))
        self.write_profile(duplicate=True)
        with self.assertRaisesRegex(ValueError, "duplicate"):
            analyze_to_evidence(client, self.request(self.root / "duplicate"))
        self.write_profile(provider="Other")
        with self.assertRaisesRegex(ValueError, "provider mismatch"):
            analyze_to_evidence(client, self.request(self.root / "provider"))
        self.assertEqual(client.calls, [])

    def test_security_receipt_must_be_complete_match_and_admit(self):
        client = FakeClient(self.response())
        self.write_security_receipt(decision="REJECTED")
        with self.assertRaisesRegex(ValueError, "did not admit"):
            analyze_to_evidence(client, self.request(self.root / "rejected"))
        self.write_security_receipt(sha256="0" * 64)
        with self.assertRaisesRegex(ValueError, "does not match input bytes"):
            analyze_to_evidence(client, self.request(self.root / "mismatch"))
        self.write_security_receipt(mime_type="image/png")
        with self.assertRaisesRegex(ValueError, "content type does not match"):
            analyze_to_evidence(client, self.request(self.root / "wrong-type"))
        self.write_security_receipt(include_tools=False)
        with self.assertRaisesRegex(ValueError, "tool evidence is missing"):
            analyze_to_evidence(client, self.request(self.root / "missing-tools"))
        self.assertEqual(client.calls, [])

    def test_provider_and_output_failures_are_atomic(self):
        for index, response in enumerate(
            [object(), documentai.ProcessResponse(document=documentai.Document())]
        ):
            output = self.root / f"failed-{index}"
            with self.assertRaises(RuntimeError):
                analyze_to_evidence(FakeClient(response), self.request(output))
            self.assertFalse(output.exists())
        occupied = self.root / "occupied"
        occupied.mkdir()
        client = FakeClient(self.response())
        with self.assertRaises(FileExistsError):
            analyze_to_evidence(client, self.request(occupied))
        self.assertEqual(client.calls, [])

    def test_extension_and_resource_fail_closed(self):
        unsupported = self.root / "invoice.exe"
        unsupported.write_bytes(b"MZ")
        with self.assertRaisesRegex(ValueError, "unsupported"):
            analyze_to_evidence(
                FakeClient(self.response()),
                AnalysisRequest(
                    "supplier-invoice", unsupported, self.root / "never", self.profile, self.security
                ),
            )
        with self.assertRaisesRegex(ValueError, "exact processor version"):
            parse_processor_version_resource("projects/p/locations/us/processors/x")
        with self.assertRaisesRegex(ValueError, "unsafe segment"):
            parse_processor_version_resource(
                "projects/p/locations/us/processors/x/processorVersions/.."
            )

    def test_cli_cannot_override_profile_processor(self):
        destinations = {action.dest for action in build_parser()._actions}
        self.assertIn("document_class", destinations)
        self.assertIn("profile", destinations)
        self.assertIn("security_receipt", destinations)
        self.assertNotIn("processor_version_resource", destinations)


if __name__ == "__main__":
    unittest.main()
````

### FILE: `google_document_runtime/README.md`
```yaml
block_id: "GOOGLE-DOCUMENT-AI-RUNTIME:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "d35a36ec1612cb2ab7c041b93ff5ad05bad0dcd4dfcf8b8ee4fe158ee23e823b"
variables: []
secrets_allowed: false
```
````markdown
# Google Cloud Document AI runtime

This component calls Google Cloud Document AI through the exact official Python SDK 3.15.0. Google owns the SDK, generated clients, request/response models and service. The profile, security chaining, atomic evidence writer and tests in this directory are `AUTHORED` orchestration and are not presented as Google code.

The command no longer accepts a processor resource chosen by the caller. It requires `--document-class`, an approved profile and an `ADMITTED` receipt from the secure local-file gate. The selected class must be uniquely declared as `REQUIRED`, must contain an exact `projects/{project}/locations/{location}/processors/{processor}/processorVersions/{version}` resource, and must explicitly deny automatic storage. Input SHA-256, bytes and detected MIME must match the security receipt, including ClamAV, Magika and YARA-X evidence.

The distributed template intentionally has `CONFIGURATION_REQUIRED` for supplier invoices and blocks every other class. A project agent must ask the owner for the Google project/location/processor/version, approved schema version, account/IAM, region/residency, corpus and review route; it may change a class to `REQUIRED` only after those values and their project evidence exist. It must never invent a processor or silently use the mutable default version.

For the audited target, install the frozen 23-wheel graph with hashes. It includes the exact Google Cloud Storage 3.13.1 client required by the official batch-processing sample:

```text
python -m pip install --require-hashes -r requirements-windows-py312.lock
python -m pip check
python -m unittest -v test_document_ai_runtime.py
python run_document_ai.py --document-class supplier-invoice --profile <approved-profile.json> --security-receipt <security-receipt.json> --input <local.pdf> --output <new-evidence-directory>
```

Authentication uses Google Application Default Credentials. Credential JSON and tokens are never CLI values or evidence fields. The endpoint is derived from the profile-owned processor version. The complete official `ProcessResponse` is preserved and hash-linked, but no extracted field is written to a business database. The strict field-evaluation gate must validate field evidence and the representative ground-truth corpus before any later storage decision.

The exact tag, commit, source archive, wheel and sdist were audited separately. The current SDK is `REBUILD_VERIFIED / CONDITIONED`, not a reusable product claim: its official unit suite resolves an open dependency graph and, on the dated CPython 3.12 Windows audit, both protobuf variants produced 1,879 PASS and six mTLS failures. No live GCP call, processor accuracy, billing, latency, quota or production readiness is claimed.
````

## 6. Configuration surface

| Variable/argument | Tipo | Default | Validación | Secreto | Mutabilidad/efecto |
|---|---|---|---|---|---|
| `document_class` | id de clase | none | existe una vez y decisión `REQUIRED` | no | por ejecución; no elige processor libre |
| `profile` | JSON local aprobado | none | schema/provider/SDK/clases/recurso/schema version/storage deny | no | cambio exige evaluación y rollback |
| `security_receipt` | JSON local | none | `ADMITTED`, policy/tools/resultados y SHA/bytes/MIME coincidentes | no | nueva por input |
| `input` | path local | none | archivo regular no symlink, MIME permitido, tamaño | no | por ejecución |
| `output` | path nuevo | none | no debe existir; commit atómico | no | por ejecución |
| ADC | reference externa | Google chain | identidad/IAM del entorno | sí | nunca se serializa |
| `max_bytes` | entero | 20 MiB | mayor que cero | no | reinicio de ejecución |

No hay variable que habilite almacenamiento automático. Una configuración inválida falla antes de llamar al provider.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| `google-cloud-documentai` | 3.15.0 wheel + SHA-256 | cliente v1, protos, ADC | Apache-2.0 | runtime | PyPI/Google Cloud |
| `google-cloud-storage` | 3.13.1 wheel + SHA-256 | cliente oficial requerido por batch para resultados JSON en GCS | Apache-2.0 | runtime | PyPI/Google Cloud |
| grafo productivo | 23 wheels exactos + SHA-256 | SDKs, gRPC, auth, TLS, protobuf, HTTP y CRC32C/media resumable | licencias de cada distribución | runtime | PyPI |
| CPython | 3.12 / Windows x86-64 target auditado | runner/tests | PSF-2.0 | runtime/test | python.org |

`requirements-direct.in` conserva la identidad del artefacto Google directo. `requirements-windows-py312.lock` congela el target auditado con `--require-hashes`; otro OS/arquitectura/Python debe resolver, hashear, escanear, probar y aprobar su propio lock sin reutilizar éste.

## 8. Apply order

En workspace vacío: componer adquisición + gate local + artefacto SDK + runtime + evaluación estricta; adquirir/verificar source/wheels; instalar el lock con hashes; ejecutar seguridad sobre el input; pedir y completar proyecto/región/processor version/schema/corpus/review; promover sólo esa clase a `REQUIRED`; configurar ADC/IAM y ejecutar tests antes de una llamada sandbox. En workspace existente: detectar colisiones, no sobrescribir, integrar mediante compositor y mantener runtime/evidencia fuera del path de almacenamiento empresarial.

Para rollback, retirar tráfico, restaurar `processor_version_resource` anterior, conservar evidencia, revertir el adapter/mapping y repetir evaluación. Nunca borrar originales ni reusar un output parcial.

## 9. Verification

```text
python -m pip install --require-hashes -r requirements-windows-py312.lock
python -m pip check
python -m unittest -v test_document_ai_runtime.py
python -c "import inspect, importlib.metadata as m; from google.cloud import documentai_v1 as d, storage; assert m.version('google-cloud-documentai') == '3.15.0'; assert m.version('google-cloud-storage') == '3.13.1'; assert storage.Client; assert 'ProcessRequest' in str(inspect.signature(d.DocumentProcessorServiceClient.process_document))"
```

Éxito local esperado: seis tests, lock/hash install, `pip check`, request oficial, perfil como autoridad, receipt completa/vinculada, negativos atómicos y `GOOGLE_DOCUMENT_AI_SDK_CONTRACT_PASS`. La suite Google no está completamente verde con su grafo móvil actual: ambas variantes retienen seis fallos mTLS. Antes de uso empresarial faltan llamada sandbox real, corpus/ground truth por clase, métricas por campo y ambigüedad, revisión, privacidad/residencia, costo/cuota, carga/fallos, IAM y rollback.

## 10. Reconstruction evidence

- identidad: tag oficial `google-cloud-documentai-v3.15.0`, commit verificado `5accbb42b3f8372c4d03ef59a6b0000b165f2200`, source archive, sdist y wheel fijados; 75 archivos sdist y 67 wheel comparados byte a byte con el tag, cero divergencias;
- entorno limpio: CPython 3.12.14 Windows x86-64; 23 wheels productivos exactos, `pip --require-hashes`, `pip check`, imports Document AI/Storage y consulta oficial Google OSV fechada con cero matches conocidos;
- materialización: 6/6 archivos con SHA-256 verificado;
- tests AUTHORED: 6/6 PASS; probe SDK PASS; prueba adversarial previa conservada como `LIB-FAIL-231` y cerrada por regresión;
- tests Google: protobuf `upb` 1879 PASS/6 FAIL y `python` 1879 PASS/6 FAIL; `UP-FAIL-075`, sin parche local ni afirmación verde;
- llamadas cloud: ninguna; no existían credenciales, processor ni corpus autorizados;
- divergencias: glue y tests son `AUTHORED`; no se redistribuye código fuente del SDK, sólo referencia al wheel oficial Apache-2.0 exacto;
- fecha/revisor: 2026-08-27 / Codex; expediente `GOOGLE_DOCUMENT_AI_PYTHON_2026-08-27_V2.md`.
