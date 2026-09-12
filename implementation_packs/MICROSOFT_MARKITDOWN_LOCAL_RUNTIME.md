# Microsoft MarkItDown Local Runtime

## 1. Metadata

```yaml
pack_id: "MICROSOFT-MARKITDOWN-LOCAL-RUNTIME"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un runtime local fail-closed sobre Microsoft MarkItDown 0.1.7: perfil V2 sin formatos habilitados por defecto, receipt ClamAV/Magika/YARA-X obligatorio, grafo Windows CPython 3.12 de 44 wheels fijados y evidencia atómica hash-linked."
stacks: ["CPython 3.12 Windows x86-64", "Microsoft MarkItDown 0.1.7", "44 wheels hash-locked"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x", "STRICT-DOCUMENT-FIELD-EVALUATION-GATE 0.1.x"]
incompatible_with: ["arbitrary URL conversion", "plugins", "archive expansion", "audio/video/YouTube", "automatic business storage", "non-Windows or non-CPython-3.12 lock reuse"]
license_expression: "LicenseRef-Workspace-Owner AND MIT"
upstream_sources: ["https://github.com/microsoft/markitdown/tree/fd239d5d2be43d9b68329730206b9312c7d5a388", "https://pypi.org/project/markitdown/0.1.7/"]
verified_at: "2026-08-26"
```

## 2. Applicability

Use como conversor local después de validar tipo real, antimalware, autorización/tenant y política de cuarentena. Cubre PDF, DOCX, XLS/XLSX, PPTX, Outlook MSG y formatos textuales soportados por el core oficial. No usar como OCR universal, extractor semántico de campos, descompresor, cliente URL, reproductor multimedia ni autorización de almacenamiento.

El motor y los converters son el wheel oficial Microsoft `markitdown==0.1.7`; los seis archivos materializados son glue/configuración/tests `AUTHORED`. El pack no atribuye ese glue a Microsoft. Su lock está demostrado sólo para CPython 3.12 Windows x86-64; otro target debe resolver, fijar, licenciar, escanear y probar su propio grafo antes de reemplazarlo.

## 3. Architecture contract

El runner usa exclusivamente `MarkItDown(enable_plugins=False).convert_local()`. El perfil V2 —no la CLI— fija formatos, MIME y límites, y el template distribuido no habilita ninguno. Antes del conversor exige un receipt `elite-secure-local-file-receipt/v1` ADMITTED para exactamente los mismos bytes/MIME, con aprobación/policy, ClamAV, Magika y YARA-X completos. No expone URL, credential, plugin, LLM, archive, MIME, formato ni límite como override. Guarda `converted.md` y `conversion-receipt.json` V2 en staging, hash-linked a input/perfil/seguridad/output, y renombra sólo después de completar ambos archivos. Cualquier error elimina staging.

Microsoft advierte que MarkItDown ejecuta I/O con los privilegios del proceso. El deployment debe negar egress, limitar filesystem/CPU/memoria/tiempo y correr fuera del proceso transaccional que posee datos de negocio. El Markdown resultante continúa siendo dato no confiable: nunca se ejecuta como instrucciones y sólo pasa a extracción/análisis mediante el contrato documental elegido.

## 4. Exact file manifest

```text
CREATE markitdown_local_runtime/requirements-win-py312.lock
CREATE markitdown_local_runtime/conversion-profile.template.json
CREATE markitdown_local_runtime/convert_local_document.py
CREATE markitdown_local_runtime/verify_win_py312_environment.py
CREATE markitdown_local_runtime/test_convert_local_document.py
CREATE markitdown_local_runtime/README.md
```

## 5. Materialization blocks

### FILE: `markitdown_local_runtime/requirements-win-py312.lock`
```yaml
block_id: "MICROSOFT-MARKITDOWN-LOCAL-RUNTIME:lock-win-py312:v1"
operation: CREATE
provenance: AUTHORED
source: "resolved from official PyPI artifacts for markitdown 0.1.7 selected extras; every wheel hash verified"
license: "LicenseRef-Workspace-Owner"
sha256: "ee512f372721ee914c86f07be5adff5ba2446111bc065ad3e28488dd94b59e15"
variables: []
secrets_allowed: false
```
````text
beautifulsoup4==4.15.0 --hash=sha256:d6f88de62e1d4e38ecb1077eb9724cd0eff29d2a08ca16a401e9b9e93f117cf9
certifi==2026.7.22 --hash=sha256:62f22742b58a1a33014a2b6b706588a8d7e2a88ae7bd1a6ebe8c992928483775
cffi==2.1.1 --hash=sha256:f53e442b08449d42821fa4a4fba000095af9f62742a500f978a9f557ec44339a
charset-normalizer==3.5.1 --hash=sha256:3617ac3cfd8b9888f145ad89dd6e692285834b0201c6074a5eeaad3fd4d668c2
click==8.5.0 --hash=sha256:255bc9599cf7748b4b1a446ccc735421bd08a2ae529a8b88597d3de5664ee360
cobble==0.1.4 --hash=sha256:36c91b1655e599fd428e2b95fdd5f0da1ca2e9f1abb0bc871dec21a0e78a2b44
coloredlogs==15.0.1 --hash=sha256:612ee75c546f53e92e70049c9dbfcc18c935a2b9a53b66085ce9ef6a6e5c0934
cryptography==50.0.1 --hash=sha256:aed8db4f6d71c51efb89530e12d9464e7bf2923d46c3205dc794a2a93f8c0648
defusedxml==0.7.1 --hash=sha256:a352e7e428770286cc899e2542b6cdaedb2b4953ff269a210103ec58f6198a61
et_xmlfile==2.0.0 --hash=sha256:7a91720bc756843502c3b7504c77b8fe44217c85c537d85037f0f536151b2caa
flatbuffers==25.12.19 --hash=sha256:7634f50c427838bb021c2d66a3d1168e9d199b0607e6329399f04846d42e20b4
humanfriendly==10.0 --hash=sha256:1697e1a8a8f550fd43c2865cd84542fc175a61dcb779b6fee18cf6b6ccba1477
idna==3.19 --hash=sha256:815e7be7a7806d54abb586dc943addc79e8b2ee16915059658cbeff4b1b43bf4
lxml==6.1.2 --hash=sha256:b97153ca609b434b712ddfb92cd6af101a7045a7724c542258bd4727a344472f
magika==0.6.3 --hash=sha256:e57f75674447b20cab4db928ae58ab264d7d8582b55183a0b876711c2b2787f3
mammoth==1.11.0 --hash=sha256:c077ab0d450bd7c0c6ecd529a23bf7e0fa8190c929e28998308ff4eada3f063b
markdownify==1.2.3 --hash=sha256:a189a0bedfd14009030fde5f85bb6f77c56897cb839b5c25315dd7d4e3e290ba
markitdown==0.1.7 --hash=sha256:4eca912c87c6aa6897284a7f4bf6769a23bccf8544530f5d8b175fbe3797c916
mpmath==1.3.0 --hash=sha256:a0b2b9fe80bbcd81a6647ff13108738cfb482d481d826cc0e02f5b35e5c88d2c
numpy==2.5.2 --hash=sha256:28ac63476ec7651484215ee7fa15a1f78b57c14621f01e392afe17b9a1390ce4
olefile==0.47 --hash=sha256:543c7da2a7adadf21214938bb79c83ea12b473a4b6ee4ad4bf854e7715e13d1f
onnxruntime==1.20.1 --hash=sha256:19c2d843eb074f385e8bbb753a40df780511061a63f9def1b216bf53860223fb
openpyxl==3.1.5 --hash=sha256:5282c12b107bffeef825f4617dc029afaf41d0ea60823bbb665ef3079dc79de2
packaging==26.3 --hash=sha256:d7193f7c8e4e93f444fde0262bf90af30e16fa0ad0ad44cb553c87339b23cd1c
pandas==3.0.5 --hash=sha256:80a611068e8a3ac23f7398c6c14eb46dc974e5cc9997f653e2dcfd1da74edd41
pdfminer.six==20260107 --hash=sha256:366585ba97e80dffa8f00cebe303d2f381884d8637af4ce422f1df3ef38111a9
pdfplumber==0.11.10 --hash=sha256:7741ea81bf165b474b153e6789d10d18e06b6ddcf3ec84289c3ef2fed6802580
pillow==12.3.0 --hash=sha256:a2b55dd6b2a4c4b7d87ffa56bdb33fdc5fdb9a462173861a7bc097f17d91cb09
protobuf==7.36.0 --hash=sha256:1781cc1de61249b750848029bca452c0a8b7e990080316b9bbc2518b2117b488
pycparser==3.0 --hash=sha256:b727414169a36b7d524c1c3e31839a521725078d7b2ff038656844266160a992
pypdfium2==5.13.0 --hash=sha256:47dcca2a8d507b5fd24f94c3c9d48fb379430f097bc20f01beff6c963ffbcedb
pyreadline3==3.5.6 --hash=sha256:8449b734232e42a5dcd74048e39b60db2839a4c38cf3ae2bf7707d58b5389c0d
python-dateutil==2.9.0.post0 --hash=sha256:a8b2bc7bffae282281c8140a97d3aa9c14da0b136dfe83f850eea9a5f7470427
python-dotenv==1.2.3 --hash=sha256:904552145e8bfed22162c09dab1c2b9b54fefa7b23ba780f4f26ca0316b0f0d9
python-pptx==1.0.2 --hash=sha256:160838e0b8565a8b1f67947675886e9fea18aa5e795db7ae531606d68e785cba
requests==2.34.2 --hash=sha256:2a0d60c172f83ac6ab31e4554906c0f3b3588d37b5cb939b1c061f4907e278e0
six==1.17.0 --hash=sha256:4721f391ed90541fddacab5acf947aa0d3dc7d27b2e1e8eda2be8970586c3274
soupsieve==2.9.2 --hash=sha256:8089a26fd974ca7a1f30276d3d8492ab266ab15af581642dfe8aa162e0c1c823
sympy==1.14.0 --hash=sha256:e091cc3e99d2141a0ba2847328f5479b05d94a6635cb96148ccb3f34671bd8f5
typing_extensions==4.16.0 --hash=sha256:481caa481374e813c1b176ada14e97f1f67a4539ce9cfeb3f350d78d6370c2e8
tzdata==2026.3 --hash=sha256:dc096730c87af6cab1b171c9d532be840741ff5d459015e7f6947bd7d7e54931
urllib3==2.7.0 --hash=sha256:9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897
xlrd==2.0.2 --hash=sha256:ea762c3d29f4cca48d82df517b6d89fbce4db3107f9d78713e48cd321d5c9aa9
xlsxwriter==3.2.9 --hash=sha256:9a5db42bc5dff014806c58a20b9eae7322a134abb6fce3c92c181bfb275ec5b3
````

### FILE: `markitdown_local_runtime/conversion-profile.template.json`
```yaml
block_id: "MICROSOFT-MARKITDOWN-LOCAL-RUNTIME:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "fail-closed local conversion profile derived from Microsoft security guidance"
license: "LicenseRef-Workspace-Owner"
sha256: "77575ed1d0f0123a1836989b515b62bf256d07b97221bd8db1517b8e3c4639dd"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-markitdown-local-conversion-profile/v2",
  "provider": "Microsoft MarkItDown",
  "runtime": "CPython 3.12 Windows x86-64",
  "package": "markitdown==0.1.7",
  "max_input_bytes": 52428800,
  "max_output_bytes": 104857600,
  "remote_urls": false,
  "plugins": false,
  "llm": false,
  "archive_expansion": false,
  "automatic_business_storage": false,
  "formats": [
    {"extension": ".pdf", "mime_type": "application/pdf", "decision": "BLOCKED_CONFIGURATION_REQUIRED"},
    {"extension": ".docx", "mime_type": "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "decision": "BLOCKED_CONFIGURATION_REQUIRED"},
    {"extension": ".xlsx", "mime_type": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "decision": "BLOCKED_CONFIGURATION_REQUIRED"},
    {"extension": ".pptx", "mime_type": "application/vnd.openxmlformats-officedocument.presentationml.presentation", "decision": "BLOCKED_CONFIGURATION_REQUIRED"},
    {"extension": ".txt", "mime_type": "text/plain", "decision": "BLOCKED_CONFIGURATION_REQUIRED"}
  ]
}
````

### FILE: `markitdown_local_runtime/convert_local_document.py`
```yaml
block_id: "MICROSOFT-MARKITDOWN-LOCAL-RUNTIME:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "minimal wrapper invoking Microsoft MarkItDown convert_local public API"
license: "LicenseRef-Workspace-Owner"
sha256: "7475c4d1d033b2fc16adaa23485f7056ebc3e2f5bef5d030d8edd8dc8976be09"
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
import uuid
from typing import Any, Protocol


MARKITDOWN_VERSION = "0.1.7"
PROFILE_SCHEMA = "elite-markitdown-local-conversion-profile/v2"
SECURITY_RECEIPT_SCHEMA = "elite-secure-local-file-receipt/v1"
FORMAT_MIME_TYPES = {
    ".csv": "text/csv",
    ".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    ".htm": "text/html",
    ".html": "text/html",
    ".ipynb": "application/x-ipynb+json",
    ".json": "application/json",
    ".md": "text/markdown",
    ".msg": "application/vnd.ms-outlook",
    ".pdf": "application/pdf",
    ".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
    ".txt": "text/plain",
    ".xls": "application/vnd.ms-excel",
    ".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    ".xml": "application/xml",
}


class ConversionResult(Protocol):
    text_content: str


class LocalConverter(Protocol):
    def convert_local(self, path: Path) -> ConversionResult: ...


@dataclass(frozen=True)
class ConversionRequest:
    input_path: Path
    output_directory: Path
    profile_path: Path
    security_receipt_path: Path


@dataclass(frozen=True)
class FormatSelection:
    extension: str
    mime_type: str
    max_input_bytes: int
    max_output_bytes: int
    profile_sha256: str


def sha256_bytes(payload: bytes) -> str:
    return hashlib.sha256(payload).hexdigest()


def is_sha256(value: object) -> bool:
    return isinstance(value, str) and len(value) == 64 and all(character in "0123456789abcdef" for character in value)


def read_json_file(path: Path, label: str) -> tuple[dict[str, Any], bytes]:
    resolved = path.resolve(strict=True)
    if path.is_symlink() or not resolved.is_file():
        raise ValueError(f"{label} must be a regular non-symlink file")
    payload = resolved.read_bytes()
    try:
        value = json.loads(payload)
    except (UnicodeDecodeError, json.JSONDecodeError) as error:
        raise ValueError(f"{label} is not valid UTF-8 JSON") from error
    if not isinstance(value, dict):
        raise ValueError(f"{label} must be a JSON object")
    return value, payload


def select_format(profile_path: Path, extension: str) -> FormatSelection:
    profile, profile_bytes = read_json_file(profile_path, "conversion profile")
    if profile.get("schema") != PROFILE_SCHEMA:
        raise ValueError("unsupported conversion profile schema")
    if profile.get("provider") != "Microsoft MarkItDown":
        raise ValueError("conversion profile provider mismatch")
    if profile.get("runtime") != "CPython 3.12 Windows x86-64":
        raise ValueError("conversion profile runtime mismatch")
    if profile.get("package") != f"markitdown=={MARKITDOWN_VERSION}":
        raise ValueError("conversion profile package mismatch")
    for control in ("remote_urls", "plugins", "llm", "archive_expansion", "automatic_business_storage"):
        if profile.get(control) is not False:
            raise ValueError(f"conversion profile must explicitly disable {control}")
    max_input = profile.get("max_input_bytes")
    max_output = profile.get("max_output_bytes")
    if not isinstance(max_input, int) or isinstance(max_input, bool) or max_input <= 0:
        raise ValueError("conversion profile max_input_bytes must be positive")
    if not isinstance(max_output, int) or isinstance(max_output, bool) or max_output <= 0:
        raise ValueError("conversion profile max_output_bytes must be positive")
    formats = profile.get("formats")
    if not isinstance(formats, list) or not formats:
        raise ValueError("conversion profile formats must be a non-empty array")
    selected: dict[str, Any] | None = None
    seen: set[str] = set()
    for item in formats:
        if not isinstance(item, dict):
            raise ValueError("conversion profile contains an invalid format")
        item_extension = item.get("extension")
        if not isinstance(item_extension, str) or item_extension not in FORMAT_MIME_TYPES:
            raise ValueError("conversion profile contains an unsupported extension")
        if item_extension in seen:
            raise ValueError("conversion profile contains duplicate extensions")
        seen.add(item_extension)
        if item.get("mime_type") != FORMAT_MIME_TYPES[item_extension]:
            raise ValueError("conversion profile MIME type mismatch")
        if item_extension == extension:
            selected = item
    if selected is None:
        raise ValueError("input extension is not declared by the conversion profile")
    if selected.get("decision") != "REQUIRED":
        raise ValueError("input format is blocked or not configured by the conversion profile")
    return FormatSelection(
        extension=extension,
        mime_type=FORMAT_MIME_TYPES[extension],
        max_input_bytes=max_input,
        max_output_bytes=max_output,
        profile_sha256=sha256_bytes(profile_bytes),
    )


def verify_security_receipt(receipt_path: Path, *, input_sha256: str, input_bytes: int, mime_type: str) -> str:
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
    if observed_input.get("sha256") != input_sha256 or observed_input.get("bytes") != input_bytes:
        raise ValueError("security receipt does not match input bytes")
    observed_type = receipt.get("content_type")
    if not isinstance(observed_type, dict) or observed_type.get("mime_type") != mime_type:
        raise ValueError("security receipt content type does not match input")
    tools = receipt.get("tools")
    if not isinstance(tools, dict):
        raise ValueError("security receipt tool evidence is missing")
    for name in ("clamscan", "magika", "yara_x"):
        item = tools.get(name)
        if not isinstance(item, dict) or not is_sha256(item.get("sha256")) or not isinstance(item.get("version"), str) or not item["version"]:
            raise ValueError(f"security receipt {name} evidence is invalid")
    clamav = receipt.get("clamav")
    if not isinstance(clamav, dict) or clamav.get("exit_code") != 0 or not is_sha256(clamav.get("output_sha256")):
        raise ValueError("security receipt ClamAV result is invalid")
    yara = receipt.get("yara_x")
    if not isinstance(yara, dict) or not is_sha256(yara.get("compiled_rules_sha256")) or yara.get("matching_rules") != []:
        raise ValueError("security receipt YARA-X result is invalid")
    return sha256_bytes(receipt_bytes)


def convert_to_evidence(converter: LocalConverter, request: ConversionRequest) -> dict[str, object]:
    if importlib.metadata.version("markitdown") != MARKITDOWN_VERSION:
        raise RuntimeError(f"markitdown {MARKITDOWN_VERSION} is required")
    source = request.input_path.resolve(strict=True)
    if request.input_path.is_symlink() or not source.is_file():
        raise ValueError("input must be a regular non-symlink file")
    extension = source.suffix.lower()
    selection = select_format(request.profile_path, extension)
    input_payload = source.read_bytes()
    if not input_payload or len(input_payload) > selection.max_input_bytes:
        raise ValueError("document byte size is outside the profile limit")
    input_sha256 = sha256_bytes(input_payload)
    security_receipt_sha256 = verify_security_receipt(
        request.security_receipt_path,
        input_sha256=input_sha256,
        input_bytes=len(input_payload),
        mime_type=selection.mime_type,
    )
    output = request.output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".markitdown-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        result = converter.convert_local(source)
        text = result.text_content
        if not isinstance(text, str) or not text.strip():
            raise RuntimeError("converter returned empty text")
        converted = text.encode("utf-8")
        if len(converted) > selection.max_output_bytes:
            raise RuntimeError("converted output exceeds profile limit")
        (stage / "converted.md").write_bytes(converted)
        receipt: dict[str, object] = {
            "schema": "elite-markitdown-local-conversion-receipt/v2",
            "created_at": datetime.now(timezone.utc).isoformat(),
            "converter": "Microsoft MarkItDown",
            "markitdown_version": MARKITDOWN_VERSION,
            "method": "convert_local",
            "plugins_enabled": False,
            "remote_url_used": False,
            "input_extension": extension,
            "input_mime_type": selection.mime_type,
            "input_bytes": len(input_payload),
            "input_sha256": input_sha256,
            "output_bytes": len(converted),
            "output_sha256": sha256_bytes(converted),
            "conversion_profile_sha256": selection.profile_sha256,
            "security_receipt_sha256": security_receipt_sha256,
            "business_storage_authorized": False,
        }
        receipt_bytes = (json.dumps(receipt, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "conversion-receipt.json").write_bytes(receipt_bytes)
        stage.rename(output)
        return receipt
    except Exception:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description="Convert one security-admitted local document with Microsoft MarkItDown")
    parser.add_argument("--input", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    parser.add_argument("--profile", type=Path, required=True)
    parser.add_argument("--security-receipt", type=Path, required=True)
    return parser


def main() -> int:
    args = build_parser().parse_args()
    from markitdown import MarkItDown

    converter = MarkItDown(enable_plugins=False)
    receipt = convert_to_evidence(
        converter,
        ConversionRequest(args.input, args.output, args.profile, args.security_receipt),
    )
    print(json.dumps(receipt, sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `markitdown_local_runtime/verify_win_py312_environment.py`
```yaml
block_id: "MICROSOFT-MARKITDOWN-LOCAL-RUNTIME:environment:v1"
operation: CREATE
provenance: AUTHORED
source: "exact platform and installed distribution verifier for the resolved wheel graph"
license: "LicenseRef-Workspace-Owner"
sha256: "2ff9321dabe77ad5e8f4ba133baf1fb9aaa33cab168c2010d8aa8961ffd13019"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import importlib.metadata
from pathlib import Path
import platform
import re
import sys


LOCK_LINE = re.compile(r"^([A-Za-z0-9_.-]+)==([^ ]+) --hash=sha256:([0-9a-f]{64})$")


def main() -> int:
    if sys.implementation.name != "cpython" or sys.version_info[:2] != (3, 12):
        raise RuntimeError("this verified lock requires CPython 3.12")
    if sys.platform != "win32" or platform.machine().lower() not in {"amd64", "x86_64"}:
        raise RuntimeError("this verified lock requires Windows x86-64")
    lock = Path(__file__).with_name("requirements-win-py312.lock")
    expected: dict[str, str] = {}
    for number, raw in enumerate(lock.read_text(encoding="utf-8").splitlines(), 1):
        match = LOCK_LINE.fullmatch(raw)
        if match is None:
            raise RuntimeError(f"invalid lock line {number}")
        name, version, _ = match.groups()
        normalized = re.sub(r"[-_.]+", "-", name).lower()
        if normalized in expected:
            raise RuntimeError(f"duplicate lock distribution: {name}")
        expected[normalized] = version
    if len(expected) != 44:
        raise RuntimeError(f"expected 44 locked distributions, found {len(expected)}")
    for name, version in expected.items():
        observed = importlib.metadata.version(name)
        if observed != version:
            raise RuntimeError(f"distribution mismatch: {name} expected={version} observed={observed}")
    print("MARKITDOWN_WIN_PY312_ENVIRONMENT_PASS distributions=44")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `markitdown_local_runtime/test_convert_local_document.py`
```yaml
block_id: "MICROSOFT-MARKITDOWN-LOCAL-RUNTIME:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "unit regressions for local-only invocation, hashes and atomic failure"
license: "LicenseRef-Workspace-Owner"
sha256: "ef2e86ddad63e4c027a5ecacac00514eb0a1eee83a61a6bfd4492d1533ebb4ae"
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
from unittest.mock import patch

from convert_local_document import ConversionRequest, build_parser, convert_to_evidence


class FakeResult:
    def __init__(self, text: str) -> None:
        self.text_content = text


class FakeConverter:
    def __init__(self, text: str = "# Converted\n\n|a|b|\n|-|-|\n|1|2|\n") -> None:
        self.text = text
        self.calls: list[Path] = []

    def convert_local(self, path: Path) -> FakeResult:
        self.calls.append(path)
        return FakeResult(self.text)


class LocalConversionTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary = tempfile.TemporaryDirectory()
        self.root = Path(self.temporary.name)
        self.source = self.root / "invoice.pdf"
        self.payload = b"%PDF-1.7 official-fixture-placeholder"
        self.source.write_bytes(self.payload)
        self.profile = self.root / "conversion-profile.json"
        self.security = self.root / "security-receipt.json"
        self.write_profile()
        self.write_security_receipt()

    def tearDown(self) -> None:
        self.temporary.cleanup()

    def write_profile(
        self,
        *,
        decision: str = "REQUIRED",
        package: str = "markitdown==0.1.7",
        mime_type: str = "application/pdf",
        duplicate: bool = False,
        storage: bool = False,
        max_output: int = 1024,
    ) -> None:
        formats = [{"extension": ".pdf", "mime_type": mime_type, "decision": decision}]
        if duplicate:
            formats.append(formats[0].copy())
        value = {
            "schema": "elite-markitdown-local-conversion-profile/v2",
            "provider": "Microsoft MarkItDown",
            "runtime": "CPython 3.12 Windows x86-64",
            "package": package,
            "max_input_bytes": 1024,
            "max_output_bytes": max_output,
            "remote_urls": False,
            "plugins": False,
            "llm": False,
            "archive_expansion": False,
            "automatic_business_storage": storage,
            "formats": formats,
        }
        self.profile.write_text(json.dumps(value) + "\n", encoding="utf-8", newline="\n")

    def write_security_receipt(
        self,
        *,
        decision: str = "ADMITTED",
        sha256: str | None = None,
        mime_type: str = "application/pdf",
        include_tools: bool = True,
        storage: bool = False,
    ) -> None:
        value = {
            "schema": "elite-secure-local-file-receipt/v1",
            "decision": decision,
            "reason": "ALL_SELECTED_GATES_PASSED" if decision == "ADMITTED" else "TEST_REJECTION",
            "approval_id": "security-approval-test",
            "policy_sha256": "1" * 64,
            "input": {"sha256": sha256 or hashlib.sha256(self.payload).hexdigest(), "bytes": len(self.payload)},
            "content_type": {"mime_type": mime_type, "label": "pdf", "score": 1.0},
            "clamav": {"exit_code": 0, "output_sha256": "2" * 64},
            "yara_x": {"compiled_rules_sha256": "3" * 64, "matching_rules": []},
            "business_storage_authorized": storage,
        }
        if include_tools:
            value["tools"] = {
                "clamscan": {"sha256": "4" * 64, "version": "ClamAV 1.5.1"},
                "magika": {"sha256": "5" * 64, "version": "magika 1.0.1"},
                "yara_x": {"sha256": "6" * 64, "version": "1.20.0"},
            }
        self.security.write_text(json.dumps(value) + "\n", encoding="utf-8", newline="\n")

    def request(self, output: Path, *, profile: Path | None = None, security: Path | None = None) -> ConversionRequest:
        return ConversionRequest(self.source, output, profile or self.profile, security or self.security)

    def run_conversion(self, converter: FakeConverter, request: ConversionRequest) -> dict[str, object]:
        with patch.object(importlib.metadata, "version", return_value="0.1.7"):
            return convert_to_evidence(converter, request)

    def test_positive_is_atomic_and_hash_links_profile_security_input_and_output(self) -> None:
        output = self.root / "evidence"
        fake = FakeConverter()
        receipt = self.run_conversion(fake, self.request(output))
        converted = (output / "converted.md").read_bytes()
        persisted = json.loads((output / "conversion-receipt.json").read_text(encoding="utf-8"))
        self.assertEqual(receipt, persisted)
        self.assertEqual(receipt["schema"], "elite-markitdown-local-conversion-receipt/v2")
        self.assertEqual(fake.calls, [self.source.resolve()])
        self.assertEqual(receipt["input_sha256"], hashlib.sha256(self.payload).hexdigest())
        self.assertEqual(receipt["output_sha256"], hashlib.sha256(converted).hexdigest())
        self.assertEqual(receipt["conversion_profile_sha256"], hashlib.sha256(self.profile.read_bytes()).hexdigest())
        self.assertEqual(receipt["security_receipt_sha256"], hashlib.sha256(self.security.read_bytes()).hexdigest())
        self.assertFalse(receipt["plugins_enabled"])
        self.assertFalse(receipt["remote_url_used"])
        self.assertFalse(receipt["business_storage_authorized"])

    def test_profile_is_execution_authority_and_template_enables_nothing(self) -> None:
        fake = FakeConverter()
        template = Path(__file__).with_name("conversion-profile.template.json")
        for profile in (template, self.profile):
            if profile == self.profile:
                self.write_profile(decision="BLOCKED_CONFIGURATION_REQUIRED")
            with self.assertRaisesRegex(ValueError, "blocked or not configured"):
                self.run_conversion(fake, self.request(self.root / f"out-{profile.name}", profile=profile))
        for mutation, message in (
            ({"package": "markitdown==9.9.9"}, "package mismatch"),
            ({"mime_type": "application/zip"}, "MIME type mismatch"),
            ({"duplicate": True}, "duplicate extensions"),
            ({"storage": True}, "disable automatic_business_storage"),
        ):
            self.write_profile(**mutation)
            with self.assertRaisesRegex(ValueError, message):
                self.run_conversion(fake, self.request(self.root / f"never-{message[:3]}"))
        self.assertEqual(fake.calls, [])

    def test_security_receipt_must_be_complete_match_and_admit(self) -> None:
        fake = FakeConverter()
        mutations = (
            ({"decision": "REJECTED"}, "did not admit"),
            ({"sha256": "0" * 64}, "does not match input bytes"),
            ({"mime_type": "application/zip"}, "content type does not match"),
            ({"include_tools": False}, "tool evidence is missing"),
            ({"storage": True}, "deny business storage"),
        )
        for index, (mutation, message) in enumerate(mutations):
            self.write_security_receipt(**mutation)
            with self.assertRaisesRegex(ValueError, message):
                self.run_conversion(fake, self.request(self.root / f"never-{index}"))
        missing = self.root / "missing-security.json"
        with self.assertRaises(FileNotFoundError):
            self.run_conversion(fake, self.request(self.root / "never-missing", security=missing))
        self.assertEqual(fake.calls, [])

    def test_disguised_or_unsupported_content_is_rejected_before_converter(self) -> None:
        fake = FakeConverter()
        self.write_security_receipt(mime_type="application/zip")
        with self.assertRaisesRegex(ValueError, "content type does not match"):
            self.run_conversion(fake, self.request(self.root / "never-disguised"))
        self.source = self.root / "archive.zip"
        self.payload = b"PK\x03\x04"
        self.source.write_bytes(self.payload)
        self.write_security_receipt(mime_type="application/zip")
        with self.assertRaisesRegex(ValueError, "not declared"):
            self.run_conversion(fake, self.request(self.root / "never-archive"))
        self.assertEqual(fake.calls, [])

    def test_empty_or_oversized_provider_result_leaves_no_output(self) -> None:
        for index, text in enumerate(("", "x" * 2048)):
            output = self.root / f"evidence-{index}"
            with self.assertRaisesRegex(RuntimeError, "empty|exceeds profile"):
                self.run_conversion(FakeConverter(text), self.request(output))
            self.assertFalse(output.exists())
        self.assertEqual(list(self.root.glob(".markitdown-stage-*")), [])

    def test_existing_output_and_wrong_distribution_fail_closed(self) -> None:
        output = self.root / "evidence"
        output.mkdir()
        with self.assertRaises(FileExistsError):
            self.run_conversion(FakeConverter(), self.request(output))
        with patch.object(importlib.metadata, "version", return_value="0.1.6"):
            with self.assertRaisesRegex(RuntimeError, "0.1.7 is required"):
                convert_to_evidence(FakeConverter(), self.request(self.root / "never"))

    def test_cli_requires_profile_and_security_without_policy_overrides(self) -> None:
        destinations = {action.dest for action in build_parser()._actions}
        self.assertEqual(destinations, {"help", "input", "output", "profile", "security_receipt"})


if __name__ == "__main__":
    unittest.main()
````

### FILE: `markitdown_local_runtime/README.md`
```yaml
block_id: "MICROSOFT-MARKITDOWN-LOCAL-RUNTIME:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local operating and admission instructions tied to Microsoft upstream security guidance"
license: "LicenseRef-Workspace-Owner"
sha256: "0a9e6ff4fcb227108f563c7ab960884c8922b837104a5782f34d8b7db8e4cdc1"
variables: []
secrets_allowed: false
```
````markdown
# MarkItDown local conversion runtime

This directory is local `AUTHORED` integration code around the official Microsoft `markitdown==0.1.7` wheel. It is not Microsoft-authored code and does not imply Microsoft sponsorship.

The runtime invokes only `MarkItDown(enable_plugins=False).convert_local()`. It accepts no URL, plugin, archive-expansion, LLM, storage, MIME or size override. The conversion profile is the execution authority and the distributed template enables zero formats. A project owner must configure the exact formats, real MIME types and limits after discovery.

Every conversion requires the hash-linked `elite-secure-local-file-receipt/v1` produced for the same bytes by the secure local-file gate. The runtime checks ADMITTED/reason, policy and approval identity, MIME, ClamAV, Magika and YARA-X evidence before calling MarkItDown. It writes the Markdown and V2 receipt atomically and always records `business_storage_authorized=false`.

Install only in an empty CPython 3.12 Windows x86-64 virtual environment with the hash lock:

```text
python -m pip install --require-hashes --only-binary=:all: -r requirements-win-py312.lock
python verify_win_py312_environment.py
python -m unittest -v test_convert_local_document.py
python convert_local_document.py --input <admitted-file> --output <new-directory> --profile <approved-profile.json> --security-receipt <security-receipt.json>
```

Microsoft warns that MarkItDown performs I/O with the current process privileges. Run this lane with egress denied, minimal filesystem access and bounded CPU/memory/time. Converted Markdown remains untrusted data. Conversion does not prove semantic field accuracy and never authorizes business storage; route it into the strict field-evaluation gate with an approved closed schema and representative ground truth.
````

## 6. Configuration surface

| Variable/field | Type | Default | Regla | Secreto | Cambio |
|---|---|---|---|---|---|
| input | local path | none | archivo regular, no symlink, extensión allowlisted | no | nueva evidencia |
| output | directory | none | no debe existir | no | nunca overwrite |
| profile | JSON path | none | schema V2, provider/runtime/package, controles false, formatos únicos | no | aprobación del proyecto |
| security receipt | JSON path | none | ADMITTED, mismo SHA/bytes/MIME y evidencia completa de tres herramientas | no | por archivo |
| max input bytes | integer | profile | positivo y antes de conversión; no CLI override | no | nueva aprobación |
| max output bytes | integer | profile | positivo y antes de commit; no CLI override | no | nueva aprobación |
| formats | array | cero habilitados en template | extensión/MIME exactos; sólo `REQUIRED` ejecuta | no | nueva aprobación |
| plugins/URL/LLM | boolean | false | no configurables por CLI | no | requiere otro pack y threat review |

## 7. Dependency bill

| Dependencia | Identidad | Uso | Licencia/procedencia | Gate |
|---|---|---|---|---|
| Microsoft MarkItDown | 0.1.7 wheel SHA-256 `4eca912c…` | converters y `convert_local` | Microsoft, MIT, PyPI/source fijados | versión/import/fixtures |
| grafo Python | 43 wheels transitivos exactos | PDF/Office/MSG/core/Magika | metadata y licenses dentro de wheels; set permisivo observado | `--require-hashes`, environment verifier, legal report |
| CPython | 3.12 Windows x86-64 | runtime demostrado | PSF/runtime del target | preflight exacto |

La consulta OSV por 44 paquetes devolvió cero findings conocidos el 2026-08-27T16:33:25Z. No es garantía futura: cada instalación/release reejecuta OSV, licencias y el proceso de actualización gobernado.

## 8. Apply order

1. Clasificar esta capability como requerida y aprobar Windows CPython 3.12 o crear un lock nuevo del target.
2. Configurar el perfil V2 con formatos/MIME/límites exactos; el template bloqueado no es ejecutable.
3. Ejecutar `SECURE-LOCAL-FILE-INGESTION-GATE` y conservar su receipt ADMITTED para los mismos bytes.
4. Crear venv vacío e instalar con `--require-hashes --only-binary`.
5. Ejecutar `pip check`, environment verifier y unit tests.
6. Ejecutar fixtures oficiales y luego corpus autorizado del proyecto.
7. Conectar el Markdown no confiable al gate estricto/analyzer seleccionado; no a tablas de negocio directas.
8. Medir límites, fallos, costo operativo y rollback antes de tráfico.

## 9. Verification

```powershell
python -m pip install --require-hashes --only-binary=:all: -r requirements-win-py312.lock
python -m pip check
python verify_win_py312_environment.py
python -m unittest -v test_convert_local_document.py
python convert_local_document.py --input <official-or-authorized-fixture> --output <new-directory> --profile <approved-profile.json> --security-receipt <security-receipt.json>
```

Esperado: 44 distribuciones exactas, `pip check` verde, environment PASS, siete tests PASS y receipt V2 hash-linked. Los fixtures autorizados deben producir output no vacío. Negativos prueban template bloqueado, perfil/package/MIME/duplicado/storage inválidos, receipt ausente/rechazado/mismatched/incompleto, archivo disfrazado, archive desconocido, output vacío/grande/ocupado y versión divergente sin invocar al conversor ni hacer commit parcial.

## 10. Reconstruction evidence

Evidencia gobernante: `reconstruction_evidence/MICROSOFT_MARKITDOWN_2026-08-27_V2.md`. Evidencia histórica del lock/fixtures: `reconstruction_evidence/MICROSOFT_MARKITDOWN_LOCAL_RUNTIME_2026-08-26_V1.md`.

- source Microsoft v0.1.7/commit/archive/licencia fijados;
- wheel 71.093 bytes y SHA-256 oficial verificado;
- 44 wheels/82.672.112 bytes resueltos para CPython 3.12 Windows x86-64;
- 44 consultas OSV, cero findings conocidos;
- instalación offline, `pip check`, imports y environment verifier PASS;
- siete unit tests adversariales PASS y conversión real con wheel oficial 0.1.7 PASS;
- seis fixtures upstream oficiales convertidos con hashes conservados;
- llamadas cloud, URLs, plugins, archives, audio/video/YouTube y almacenamiento de negocio no ejecutados ni autorizados.
