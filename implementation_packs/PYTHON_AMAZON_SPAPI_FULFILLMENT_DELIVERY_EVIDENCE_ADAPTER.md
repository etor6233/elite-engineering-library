# Python Amazon SP-API Fulfillment Delivery Evidence Adapter

## 1. Metadata

```yaml
pack_id: "PYTHON-AMAZON-SPAPI-FULFILLMENT-DELIVERY-EVIDENCE-ADAPTER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa evidencia de entrega de Fulfillment Outbound v2020-07-01 con el SDK oficial Amazon amzn-sp-api 1.11.1; descarga documentos aprobados, excluye identidad/URLs firmadas y nunca inventa aceptación empresarial."
stacks: ["Python 3.10-3.14", "amzn-sp-api 1.11.1", "Amazon Fulfillment Outbound v2020-07-01"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "PYTHON-AMAZON-SPAPI-CATALOG-ADAPTER 0.1.x", "PYTHON-AMAZON-SPAPI-SHIPPING-TRACKING-ADAPTER 0.6.x", "PYTHON-AMAZON-SPAPI-EASYSHIP-HANDOVER-ADAPTER 0.1.x"]
incompatible_with: ["carrier universal", "pedido/marketplace o paquete no ligados a evidencia empresarial", "host o tipo documental no observado y aprobado", "status no terminal aprobado", "persistencia de URL firmada o identidad de destinatario", "redirect HTTP", "aceptación empresarial automática", "secreto en CLI/Markdown", "SDK distinto de 1.11.1"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0 AND MIT AND MPL-2.0 AND BSD-3-Clause"
upstream_sources: ["https://github.com/amzn/selling-partner-api-sdk/tree/8e792ae345a8d334ccdbdd03181f05f040e6a4fc", "https://github.com/amzn/selling-partner-api-models/tree/8e429486005c4ebdce5099e48cc48515a65359bb", "https://developer-docs.amazon.com/sp-api/docs/fulfillment-outbound-api-v2020-07-01-reference", "https://developer-docs.amazon.com/sp-api/docs/sp-api-release-notes?ld=SDESSOADirect", "https://developer-docs.amazon.com/sp-api/docs/get-recipient-details-for-attended-deliveries"]
verified_at: "2026-08-31"
```

## 2. Applicability

Use sólo para paquetes de órdenes Multi-Channel Fulfillment recuperados por `getFulfillmentOrder`, en marketplaces y programas demostrados. Amazon documenta que `deliveryDocumentList` puede contener foto o firma y que `dropOffLocation` informa dónde se entregó; este pack no asume un enum documental universal ni persiste sus atributos de identidad.

## 3. Architecture contract

Una aprobación liga el pedido, marketplace y cada paquete esperado al SHA-256 de un receipt empresarial. La consulta read-only usa la firma exacta del SDK fijado. Sólo se admite un documento requerido por tipo/paquete, status terminal aprobado y host HTTPS exacto. Descarga sin redirects, limita bytes, valida content type y magic bytes, escribe archivos con nombre SHA-256 mediante staging atómico y persiste únicamente una respuesta filtrada. URLs firmadas, tracking e identidad no se guardan ni se exponen en errores. El receipt conserva `automatic_business_delivery_acceptance=false`.

## 4. Exact file manifest

```text
CREATE amazon_spapi_fulfillment_delivery_evidence/requirements-direct.in
CREATE amazon_spapi_fulfillment_delivery_evidence/requirements.lock
CREATE amazon_spapi_fulfillment_delivery_evidence/sdk-artifact.lock.json
CREATE amazon_spapi_fulfillment_delivery_evidence/provider-profile.template.json
CREATE amazon_spapi_fulfillment_delivery_evidence/delivery-evidence-request.template.json
CREATE amazon_spapi_fulfillment_delivery_evidence/run_fulfillment_delivery_evidence.py
CREATE amazon_spapi_fulfillment_delivery_evidence/test_fulfillment_delivery_evidence.py
CREATE amazon_spapi_fulfillment_delivery_evidence/README.md
```

## 5. Materialization blocks

### FILE: `amazon_spapi_fulfillment_delivery_evidence/requirements-direct.in`
```yaml
block_id: "PY-AMAZON-SPAPI-FULFILLMENT-DELIVERY-EVIDENCE:REQUIREMENTS-DIRECT-IN:v1"
operation: CREATE
provenance: AUTHORED
source: "exact official PyPI wheel graph"
license: "LicenseRef-Workspace-Owner"
sha256: "a897426f3166913c0a3c5ad4d2bb2716d0c2245178bc47a18502a51700cd81ee"
variables: []
secrets_allowed: false
```
````text
amzn-sp-api @ https://files.pythonhosted.org/packages/36/3e/8744071121b96a7738fa485e584670cd4e6735b545cecc18516d10bd8881/amzn_sp_api-1.11.1-py3-none-any.whl#sha256=a88e4059c9fbfd454bfef67b82288a926f73cbdc608a8831f8383565736c70b0
````

### FILE: `amazon_spapi_fulfillment_delivery_evidence/requirements.lock`
```yaml
block_id: "PY-AMAZON-SPAPI-FULFILLMENT-DELIVERY-EVIDENCE:REQUIREMENTS-LOCK:v1"
operation: CREATE
provenance: AUTHORED
source: "exact official PyPI wheel graph"
license: "LicenseRef-Workspace-Owner"
sha256: "d8a6c4dc588887acc11a157b1cf4a84f1e279be29eb1a1aed202121f463b2d1c"
variables: []
secrets_allowed: false
```
````text
amzn-sp-api @ https://files.pythonhosted.org/packages/36/3e/8744071121b96a7738fa485e584670cd4e6735b545cecc18516d10bd8881/amzn_sp_api-1.11.1-py3-none-any.whl#sha256=a88e4059c9fbfd454bfef67b82288a926f73cbdc608a8831f8383565736c70b0
certifi @ https://files.pythonhosted.org/packages/0b/a7/71ac2cff56fec219ed242bb11b8efb69fcc4bec75db06fb7bfe35de520e6/certifi-2026.7.22-py3-none-any.whl#sha256=62f22742b58a1a33014a2b6b706588a8d7e2a88ae7bd1a6ebe8c992928483775
charset-normalizer @ https://files.pythonhosted.org/packages/cc/61/d01fc49b8dea277640b55a9e15960dbca9fdc8c9fde18e572d39c59f4019/charset_normalizer-3.5.1-py3-none-any.whl#sha256=6df0ec430f9a831772c23ca5a224cba36517a58a84bb32c32bb59a9fa67c47f6
idna @ https://files.pythonhosted.org/packages/57/b0/0e52c878c53f245edd3a11020f20979b3f490f245af532c7cae3027754b5/idna-3.19-py3-none-any.whl#sha256=815e7be7a7806d54abb586dc943addc79e8b2ee16915059658cbeff4b1b43bf4
requests @ https://files.pythonhosted.org/packages/a0/f4/c67b0b3f1b9245e8d266f0f112c500d50e5b4e83cb6f3b71b6528104182a/requests-2.34.2-py3-none-any.whl#sha256=2a0d60c172f83ac6ab31e4554906c0f3b3588d37b5cb939b1c061f4907e278e0
six @ https://files.pythonhosted.org/packages/b7/ce/149a00dd41f10bc29e5921b496af8b574d8413afcd5e30dfa0ed46c2cc5e/six-1.17.0-py2.py3-none-any.whl#sha256=4721f391ed90541fddacab5acf947aa0d3dc7d27b2e1e8eda2be8970586c3274
urllib3 @ https://files.pythonhosted.org/packages/7f/3e/5db95bcf282c52709639744ca2a8b149baccf648e39c8cc87553df9eae0c/urllib3-2.7.0-py3-none-any.whl#sha256=9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897
````

### FILE: `amazon_spapi_fulfillment_delivery_evidence/sdk-artifact.lock.json`
```yaml
block_id: "PY-AMAZON-SPAPI-FULFILLMENT-DELIVERY-EVIDENCE:SDK-ARTIFACT-LOCK-JSON:v1"
operation: CREATE
provenance: AUTHORED
source: "local receipt for exact official Amazon source/model/wheel identities"
license: "LicenseRef-Workspace-Owner"
sha256: "c719ddd61b73d8481857e9ee6c9d223da9a837209addd8d7498b7a956d1bd148"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-python-sdk-artifact-lock/v1",
  "provider": "Amazon",
  "distribution": "amzn-sp-api",
  "version": "1.11.1",
  "wheel": {"filename":"amzn_sp_api-1.11.1-py3-none-any.whl","bytes":3557108,"sha256":"a88e4059c9fbfd454bfef67b82288a926f73cbdc608a8831f8383565736c70b0"},
  "source": {"repository":"amzn/selling-partner-api-sdk","release":"python-1.11.1","commit":"8e792ae345a8d334ccdbdd03181f05f040e6a4fc","archive_bytes":26447333,"archive_sha256":"5bf712836bda76619c796fdb2b3a0a13acb281ba8bcf1b39c8f4b92f8f4799a1"},
  "model_source": {"repository":"amzn/selling-partner-api-models","commit":"8e429486005c4ebdce5099e48cc48515a65359bb","archive_sha256":"22bd3d956e12fe887f2ad893979d676b449c8bac4f6be6fd8cb0d56ae149b2e3"},
  "license_expression": "Apache-2.0",
  "verified_at": "2026-08-31"
}
````

### FILE: `amazon_spapi_fulfillment_delivery_evidence/provider-profile.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-FULFILLMENT-DELIVERY-EVIDENCE:PROVIDER-PROFILE-TEMPLATE-JSON:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed wrapper/configuration governed by official Amazon API and models"
license: "LicenseRef-Workspace-Owner"
sha256: "8fc4a554f215ae0d9427ff90dd97b085d7f1d6ad89d22fd802f4c1c799912a42"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-fulfillment-delivery-evidence-profile/v1",
  "provider": "Amazon Fulfillment Outbound v2020-07-01",
  "sdk": "amzn-sp-api==1.11.1",
  "decision": "BLOCKED_ACCESS_AND_SCOPE_APPROVAL_REQUIRED",
  "region": "SANDBOX",
  "approved_marketplace_ids": [],
  "approved_delivery_document_types": [],
  "approved_document_hosts": [],
  "approved_content_types": ["image/jpeg", "image/png", "application/pdf"],
  "approved_terminal_shipment_statuses": [],
  "max_document_bytes": 0,
  "developer_registration": "NOT_PROVEN",
  "application_registration": "NOT_PROVEN",
  "amazon_fulfillment_role_and_scope": "NOT_PROVEN",
  "mcf_program_and_marketplace_support": "NOT_PROVEN",
  "sandbox_contract_test": "NOT_PROVEN",
  "quota_and_cost_approved": false,
  "delivery_evidence_data_handling_approved": false,
  "delivery_evidence_retention_approved": false,
  "delivery_evidence_owner": "",
  "client_id_env": "AMAZON_SP_API_CLIENT_ID",
  "client_secret_env": "AMAZON_SP_API_CLIENT_SECRET",
  "refresh_token_env": "AMAZON_SP_API_REFRESH_TOKEN",
  "automatic_business_delivery_acceptance": false
}
````

### FILE: `amazon_spapi_fulfillment_delivery_evidence/delivery-evidence-request.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-FULFILLMENT-DELIVERY-EVIDENCE:DELIVERY-EVIDENCE-REQUEST-TEMPLATE-JSON:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed wrapper/configuration governed by official Amazon API and models"
license: "LicenseRef-Workspace-Owner"
sha256: "0e348b8c67885fa99d2030fcbe87e53ba3f8e0b2684d5240db81df18633b33b1"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-fulfillment-delivery-evidence-request/v1",
  "approved_by": "",
  "approved_at": "",
  "business_order_receipt_sha256": "",
  "seller_fulfillment_order_id": "",
  "marketplace_id": "",
  "expected_packages": [
    {"amazon_shipment_id": "", "package_number": 0, "required_document_types": []}
  ]
}
````

### FILE: `amazon_spapi_fulfillment_delivery_evidence/run_fulfillment_delivery_evidence.py`
```yaml
block_id: "PY-AMAZON-SPAPI-FULFILLMENT-DELIVERY-EVIDENCE:RUN-FULFILLMENT-DELIVERY-EVIDENCE-PY:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed wrapper/configuration governed by official Amazon API and models"
license: "LicenseRef-Workspace-Owner"
sha256: "4b34bce527ee52f599a60cf2630f7a60bb8527a75e16022fee7d68681c19bf4f"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
from datetime import date, datetime, timezone
from decimal import Decimal
import hashlib
import importlib.metadata
import json
import os
from pathlib import Path
import re
import shutil
import sys
from typing import Any, Callable, Protocol
import urllib.error
import urllib.parse
import urllib.request
import uuid

from spapi import SPAPIClient, SPAPIConfig
from spapi.api.fulfillment_outbound_v2020_07_01.fba_outbound_api import FbaOutboundApi

SDK_VERSION = "1.11.1"
REGIONS = {"SANDBOX", "NA", "EU", "FE"}
REQUEST_KEYS = {"schema", "approved_by", "approved_at", "business_order_receipt_sha256", "seller_fulfillment_order_id", "marketplace_id", "expected_packages"}
PACKAGE_KEYS = {"amazon_shipment_id", "package_number", "required_document_types"}
CONTENT_EXTENSIONS = {"image/jpeg": "jpg", "image/png": "png", "application/pdf": "pdf"}


class FulfillmentService(Protocol):
    def get_fulfillment_order(self, seller_fulfillment_order_id: str, **kwargs: Any) -> Any: ...


Fetcher = Callable[[str, int, float], tuple[bytes, str]]


def sha256(value: bytes) -> str:
    return hashlib.sha256(value).hexdigest()


def json_value(value: Any) -> Any:
    if value is None or isinstance(value, (str, int, float, bool)):
        return value
    if isinstance(value, Decimal):
        return str(value)
    if isinstance(value, (date, datetime)):
        return value.isoformat()
    if isinstance(value, dict):
        return {str(key): json_value(item) for key, item in value.items()}
    if isinstance(value, (list, tuple)):
        return [json_value(item) for item in value]
    data = getattr(value, "data", None)
    if data is not None:
        return json_value(data)
    to_dict = getattr(value, "to_dict", None)
    if callable(to_dict):
        return json_value(to_dict())
    raise TypeError(f"unsupported provider response type: {type(value).__name__}")


def canonical_bytes(value: Any) -> bytes:
    return (json.dumps(json_value(value), ensure_ascii=False, sort_keys=True, separators=(",", ":")) + "\n").encode("utf-8")


def pretty_bytes(value: Any) -> bytes:
    return (json.dumps(json_value(value), ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")


def exact_keys(value: Any, keys: set[str], label: str) -> dict[str, Any]:
    if not isinstance(value, dict) or set(value) != keys:
        raise ValueError(f"{label} keys must be exactly {sorted(keys)}")
    return value


def nonempty(value: Any, label: str, maximum: int = 255) -> str:
    if not isinstance(value, str) or not value or value != value.strip() or len(value) > maximum:
        raise ValueError(f"{label} must be a canonical non-empty string")
    return value


def aware_datetime(value: Any, label: str) -> datetime:
    text = nonempty(value, label, 64)
    try:
        parsed = datetime.fromisoformat(text.replace("Z", "+00:00"))
    except ValueError as exc:
        raise ValueError(f"{label} must be ISO-8601") from exc
    if parsed.tzinfo is None or parsed.utcoffset() is None:
        raise ValueError(f"{label} must include an offset")
    return parsed.astimezone(timezone.utc)


def strict_object(path: Path, label: str) -> tuple[dict[str, Any], bytes]:
    raw = path.read_bytes()
    def hook(pairs: list[tuple[str, Any]]) -> dict[str, Any]:
        result: dict[str, Any] = {}
        for key, value in pairs:
            if key in result:
                raise ValueError(f"duplicate JSON key in {label}: {key}")
            result[key] = value
        return result
    value = json.loads(raw.decode("utf-8"), object_pairs_hook=hook)
    if not isinstance(value, dict):
        raise ValueError(f"{label} must be an object")
    return value, raw


def canonical_string_list(value: Any, label: str, maximum: int = 100) -> list[str]:
    if not isinstance(value, list) or not value or len(value) > maximum or len(value) != len(set(value)):
        raise ValueError(f"{label} must be a unique non-empty list")
    for item in value:
        nonempty(item, label)
    return value


def approved_profile(path: Path) -> dict[str, Any]:
    profile, _ = strict_object(path, "provider profile")
    required = {"schema", "provider", "sdk", "decision", "region", "approved_marketplace_ids", "approved_delivery_document_types", "approved_document_hosts", "approved_content_types", "approved_terminal_shipment_statuses", "max_document_bytes", "developer_registration", "application_registration", "amazon_fulfillment_role_and_scope", "mcf_program_and_marketplace_support", "sandbox_contract_test", "quota_and_cost_approved", "delivery_evidence_data_handling_approved", "delivery_evidence_retention_approved", "delivery_evidence_owner", "client_id_env", "client_secret_env", "refresh_token_env", "automatic_business_delivery_acceptance"}
    exact_keys(profile, required, "provider profile")
    if profile["schema"] != "elite-amazon-spapi-fulfillment-delivery-evidence-profile/v1" or profile["provider"] != "Amazon Fulfillment Outbound v2020-07-01" or profile["sdk"] != "amzn-sp-api==1.11.1" or profile["decision"] != "APPROVED":
        raise PermissionError("Fulfillment Outbound profile is not approved")
    if profile["region"] not in REGIONS or profile["automatic_business_delivery_acceptance"] is not False:
        raise PermissionError("region or business acceptance contract is invalid")
    for key in ("developer_registration", "application_registration", "amazon_fulfillment_role_and_scope", "mcf_program_and_marketplace_support", "sandbox_contract_test"):
        if profile[key] != "PROVEN":
            raise PermissionError(f"{key} must be PROVEN")
    for key in ("quota_and_cost_approved", "delivery_evidence_data_handling_approved", "delivery_evidence_retention_approved"):
        if profile[key] is not True:
            raise PermissionError(f"{key} must be true")
    canonical_string_list(profile["approved_marketplace_ids"], "approved marketplaces")
    canonical_string_list(profile["approved_delivery_document_types"], "approved document types")
    hosts = canonical_string_list(profile["approved_document_hosts"], "approved document hosts")
    if any(not re.fullmatch(r"[a-z0-9.-]+", host) for host in hosts):
        raise ValueError("approved document hosts must be canonical lower-case hosts")
    content_types = canonical_string_list(profile["approved_content_types"], "approved content types")
    if any(item not in CONTENT_EXTENSIONS for item in content_types):
        raise ValueError("approved content type is unsupported by the local verifier")
    canonical_string_list(profile["approved_terminal_shipment_statuses"], "approved terminal shipment statuses")
    if not isinstance(profile["max_document_bytes"], int) or isinstance(profile["max_document_bytes"], bool) or not 1 <= profile["max_document_bytes"] <= 20_000_000:
        raise ValueError("max_document_bytes must be within 1..20000000")
    nonempty(profile["delivery_evidence_owner"], "delivery_evidence_owner")
    for key in ("client_id_env", "client_secret_env", "refresh_token_env"):
        if not re.fullmatch(r"[A-Z][A-Z0-9_]{2,127}", nonempty(profile[key], key, 128)):
            raise ValueError(f"{key} must name an environment variable")
    return profile


def validate_request(request: dict[str, Any], receipt_bytes: bytes, profile: dict[str, Any], now: datetime) -> list[dict[str, Any]]:
    exact_keys(request, REQUEST_KEYS, "delivery evidence request")
    if request["schema"] != "elite-amazon-spapi-fulfillment-delivery-evidence-request/v1":
        raise ValueError("request schema is invalid")
    nonempty(request["approved_by"], "approved_by")
    if aware_datetime(request["approved_at"], "approved_at") > now:
        raise ValueError("approved_at must not be in the future")
    if request["business_order_receipt_sha256"] != sha256(receipt_bytes):
        raise PermissionError("request is not bound to the business order receipt")
    if not re.fullmatch(r"[0-9a-f]{64}", request["business_order_receipt_sha256"]):
        raise ValueError("business receipt SHA-256 is invalid")
    nonempty(request["seller_fulfillment_order_id"], "seller_fulfillment_order_id")
    if request["marketplace_id"] not in profile["approved_marketplace_ids"]:
        raise PermissionError("marketplace is not approved")
    packages = request["expected_packages"]
    if not isinstance(packages, list) or not packages or len(packages) > 100:
        raise ValueError("expected_packages must contain 1..100 entries")
    seen: set[tuple[str, int]] = set()
    for package in packages:
        exact_keys(package, PACKAGE_KEYS, "expected package")
        shipment_id = nonempty(package["amazon_shipment_id"], "amazon_shipment_id")
        number = package["package_number"]
        if not isinstance(number, int) or isinstance(number, bool) or number < 1 or (shipment_id, number) in seen:
            raise ValueError("expected package identity is invalid or duplicated")
        seen.add((shipment_id, number))
        required_types = canonical_string_list(package["required_document_types"], "required document types", 20)
        if any(item not in profile["approved_delivery_document_types"] for item in required_types):
            raise PermissionError("requested delivery document type is not approved")
    return packages


def validate_url(url: Any, approved_hosts: list[str]) -> str:
    text = nonempty(url, "delivery document URL", 4096)
    parsed = urllib.parse.urlsplit(text)
    if parsed.scheme != "https" or parsed.username or parsed.password or parsed.port not in (None, 443) or parsed.hostname not in approved_hosts or parsed.fragment:
        raise PermissionError("delivery document URL must use an approved exact HTTPS host without credentials or fragment")
    return text


def extract_documents(response: Any, request: dict[str, Any], expected: list[dict[str, Any]], profile: dict[str, Any]) -> tuple[list[dict[str, Any]], dict[str, Any], str]:
    provider = json_value(response)
    if not isinstance(provider, dict) or provider.get("errors") not in (None, []):
        raise ValueError("provider response contains errors or has invalid envelope")
    payload = provider.get("payload")
    if not isinstance(payload, dict):
        raise ValueError("provider response payload is missing")
    order = payload.get("fulfillment_order")
    if not isinstance(order, dict) or order.get("seller_fulfillment_order_id") != request["seller_fulfillment_order_id"] or order.get("marketplace_id") != request["marketplace_id"]:
        raise ValueError("provider order identity or marketplace does not match")
    shipments = payload.get("fulfillment_shipments")
    if not isinstance(shipments, list):
        raise ValueError("fulfillment_shipments must be a list")
    index: dict[tuple[str, int], tuple[str, dict[str, Any]]] = {}
    for shipment in shipments:
        if not isinstance(shipment, dict):
            raise ValueError("fulfillment shipment is invalid")
        shipment_id = nonempty(shipment.get("amazon_shipment_id"), "provider amazon_shipment_id")
        status = nonempty(shipment.get("fulfillment_shipment_status"), "fulfillment_shipment_status")
        packages = shipment.get("fulfillment_shipment_package") or []
        if not isinstance(packages, list):
            raise ValueError("fulfillment_shipment_package must be a list")
        for package in packages:
            if not isinstance(package, dict) or not isinstance(package.get("package_number"), int):
                raise ValueError("provider package is invalid")
            key = (shipment_id, package["package_number"])
            if key in index:
                raise ValueError("provider package identity is duplicated")
            index[key] = (status, package)
    selected: list[dict[str, Any]] = []
    filtered_packages: list[dict[str, Any]] = []
    for expected_package in expected:
        key = (expected_package["amazon_shipment_id"], expected_package["package_number"])
        if key not in index:
            raise ValueError("expected package is absent from provider response")
        status, package = index[key]
        if status not in profile["approved_terminal_shipment_statuses"]:
            raise PermissionError("shipment status is not approved as terminal delivery evidence")
        delivery = package.get("delivery_information")
        if not isinstance(delivery, dict):
            raise ValueError("delivery_information is missing")
        documents = delivery.get("delivery_document_list")
        if not isinstance(documents, list):
            raise ValueError("delivery_document_list is missing")
        by_type: dict[str, list[dict[str, Any]]] = {}
        for document in documents:
            if not isinstance(document, dict):
                raise ValueError("delivery document is invalid")
            document_type = nonempty(document.get("document_type"), "delivery document type")
            by_type.setdefault(document_type, []).append(document)
        filtered_docs: list[dict[str, Any]] = []
        for required_type in expected_package["required_document_types"]:
            matches = by_type.get(required_type, [])
            if len(matches) != 1:
                raise ValueError("required delivery document must occur exactly once")
            url = validate_url(matches[0].get("url"), profile["approved_document_hosts"])
            selected.append({"amazon_shipment_id": key[0], "package_number": key[1], "shipment_status": status, "document_type": required_type, "url": url})
            filtered_docs.append({"document_type": required_type, "url_sha256": sha256(url.encode())})
        dropoff = delivery.get("drop_off_location")
        dropoff_type = nonempty(dropoff.get("type"), "drop_off_location type") if isinstance(dropoff, dict) and dropoff.get("type") is not None else None
        filtered_packages.append({"amazon_shipment_id_sha256": sha256(key[0].encode()), "package_number": key[1], "shipment_status": status, "delivery_documents": filtered_docs, "drop_off_location_type": dropoff_type, "drop_off_attributes_present": bool(isinstance(dropoff, dict) and dropoff.get("attributes"))})
    filtered = {"schema": "elite-amazon-spapi-fulfillment-delivery-filtered-response/v1", "seller_fulfillment_order_id_sha256": sha256(request["seller_fulfillment_order_id"].encode()), "marketplace_id_sha256": sha256(request["marketplace_id"].encode()), "packages": filtered_packages, "recipient_identity_values_persisted": False, "delivery_document_urls_persisted": False}
    return selected, filtered, sha256(canonical_bytes(provider))


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        return None


def default_fetcher(url: str, maximum: int, timeout: float) -> tuple[bytes, str]:
    request = urllib.request.Request(url, headers={"Accept": "image/jpeg,image/png,application/pdf", "User-Agent": "elite-delivery-evidence/1"}, method="GET")
    try:
        with urllib.request.build_opener(NoRedirect).open(request, timeout=timeout) as response:
            if getattr(response, "status", 200) != 200:
                raise ValueError("delivery document HTTP status is not 200")
            content_type = (response.headers.get_content_type() or "").lower()
            length = response.headers.get("Content-Length")
            if length is not None and int(length) > maximum:
                raise ValueError("delivery document exceeds maximum size")
            data = response.read(maximum + 1)
    except (urllib.error.URLError, OSError, ValueError) as error:
        raise RuntimeError(f"delivery document fetch failed: {type(error).__name__}") from None
    if not data or len(data) > maximum:
        raise ValueError("delivery document is empty or exceeds maximum size")
    return data, content_type


def validate_document(data: bytes, content_type: str, approved: list[str], maximum: int) -> str:
    if not data or len(data) > maximum or content_type not in approved or content_type not in CONTENT_EXTENSIONS:
        raise ValueError("delivery document size or content type is invalid")
    valid = (content_type == "image/jpeg" and data.startswith(b"\xff\xd8\xff")) or (content_type == "image/png" and data.startswith(b"\x89PNG\r\n\x1a\n")) or (content_type == "application/pdf" and data.startswith(b"%PDF-"))
    if not valid:
        raise ValueError("delivery document magic bytes do not match content type")
    return CONTENT_EXTENSIONS[content_type]


def fetch_delivery_evidence(service: FulfillmentService, request: dict[str, Any], business_receipt_bytes: bytes, profile: dict[str, Any], output_directory: Path, fetcher: Fetcher = default_fetcher, timeout_seconds: float = 30.0, now: datetime | None = None) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION:
        raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if timeout_seconds <= 0 or timeout_seconds > 120:
        raise ValueError("timeout must be within 1..120 seconds")
    target = output_directory.resolve()
    if target.exists():
        raise FileExistsError("output_directory must not exist")
    current = (now or datetime.now(timezone.utc)).astimezone(timezone.utc)
    expected = validate_request(request, business_receipt_bytes, profile, current)
    response = service.get_fulfillment_order(request["seller_fulfillment_order_id"], _request_timeout=timeout_seconds)
    selected, filtered, provider_response_sha = extract_documents(response, request, expected, profile)
    downloaded: list[tuple[dict[str, Any], bytes, str, str]] = []
    seen_hashes: set[str] = set()
    for document in selected:
        data, content_type = fetcher(document["url"], profile["max_document_bytes"], timeout_seconds)
        extension = validate_document(data, content_type, profile["approved_content_types"], profile["max_document_bytes"])
        digest = sha256(data)
        if digest in seen_hashes:
            raise ValueError("delivery documents must have unique content hashes")
        seen_hashes.add(digest)
        downloaded.append((document, data, content_type, extension))
    if target.exists():
        raise FileExistsError("output_directory appeared during download")
    stage = target.parent / f".{target.name}.tmp-{uuid.uuid4().hex}"
    if stage.exists():
        raise FileExistsError("staging directory already exists")
    stage.mkdir(parents=True)
    try:
        (stage / "provider-filtered-response.json").write_bytes(pretty_bytes(filtered))
        records: list[dict[str, Any]] = []
        for document, data, content_type, extension in downloaded:
            digest = sha256(data)
            filename = f"{digest}.{extension}"
            (stage / filename).write_bytes(data)
            records.append({"amazon_shipment_id_sha256": sha256(document["amazon_shipment_id"].encode()), "package_number": document["package_number"], "shipment_status": document["shipment_status"], "document_type": document["document_type"], "content_type": content_type, "bytes": len(data), "sha256": digest, "filename": filename})
        receipt = {"schema": "elite-amazon-spapi-fulfillment-delivery-evidence-receipt/v1", "created_at": current.isoformat(), "provider": "Amazon Fulfillment Outbound v2020-07-01", "sdk_version": SDK_VERSION, "operation": "get_fulfillment_order", "business_order_receipt_sha256": sha256(business_receipt_bytes), "request_sha256": sha256(canonical_bytes(request)), "seller_fulfillment_order_id_sha256": sha256(request["seller_fulfillment_order_id"].encode()), "marketplace_id_sha256": sha256(request["marketplace_id"].encode()), "provider_response_sha256": provider_response_sha, "filtered_response_sha256": sha256((stage / "provider-filtered-response.json").read_bytes()), "documents": records, "provider_reports_delivery_evidence": True, "recipient_identity_values_persisted": False, "delivery_document_urls_persisted": False, "automatic_business_delivery_acceptance": False}
        (stage / "DELIVERY_EVIDENCE_RECEIPT.json").write_bytes(pretty_bytes(receipt))
        if target.exists():
            raise FileExistsError("output_directory appeared during staging")
        stage.replace(target)
        return receipt
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def service_from_environment(profile: dict[str, Any]) -> FbaOutboundApi:
    values: dict[str, str] = {}
    for key in ("client_id_env", "client_secret_env", "refresh_token_env"):
        value = os.environ.get(profile[key], "")
        if not value:
            raise PermissionError(f"required secret environment variable is missing: {profile[key]}")
        values[key] = value
    config = SPAPIConfig(client_id=values["client_id_env"], client_secret=values["client_secret_env"], refresh_token=values["refresh_token_env"], region=profile["region"])
    return FbaOutboundApi(SPAPIClient(config).api_client)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--request", required=True, type=Path)
    parser.add_argument("--business-order-receipt", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--timeout-seconds", type=float, default=30.0)
    args = parser.parse_args()
    profile = approved_profile(args.profile)
    request, _ = strict_object(args.request, "delivery evidence request")
    receipt = fetch_delivery_evidence(service_from_environment(profile), request, args.business_order_receipt.read_bytes(), profile, args.output, timeout_seconds=args.timeout_seconds)
    print(json.dumps({"status": "AMAZON_SPAPI_FULFILLMENT_DELIVERY_EVIDENCE_PASS", "receipt": receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `amazon_spapi_fulfillment_delivery_evidence/test_fulfillment_delivery_evidence.py`
```yaml
block_id: "PY-AMAZON-SPAPI-FULFILLMENT-DELIVERY-EVIDENCE:TEST-FULFILLMENT-DELIVERY-EVIDENCE-PY:v1"
operation: CREATE
provenance: AUTHORED
source: "local executable contract against exact official SDK models and fail-closed wrapper"
license: "LicenseRef-Workspace-Owner"
sha256: "042555e3779d6391274f2b843184763048bc9d1dea9441b30091b832dd3c708f"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

from datetime import datetime, timezone
import inspect
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch
import urllib.error

from spapi.api.fulfillment_outbound_v2020_07_01.fba_outbound_api import FbaOutboundApi
from spapi.models.fulfillment_outbound_v2020_07_01.delivery_document import DeliveryDocument
from spapi.models.fulfillment_outbound_v2020_07_01.delivery_information import DeliveryInformation
from spapi.models.fulfillment_outbound_v2020_07_01.fulfillment_shipment_package import FulfillmentShipmentPackage

from run_fulfillment_delivery_evidence import approved_profile, default_fetcher, fetch_delivery_evidence, sha256

NOW = datetime(2026, 8, 31, 19, 0, tzinfo=timezone.utc)
BUSINESS_RECEIPT = b'{"schema":"business-order-receipt/v1"}\n'
PHOTO_URL = "https://proofs.amazon.example/photo?signature=secret"
SIGNATURE_URL = "https://proofs.amazon.example/signature?signature=secret"


def profile() -> dict:
    return {"approved_marketplace_ids": ["A1MARKET"], "approved_delivery_document_types": ["PHOTO", "SIGNATURE"], "approved_document_hosts": ["proofs.amazon.example"], "approved_content_types": ["image/jpeg", "image/png", "application/pdf"], "approved_terminal_shipment_statuses": ["SHIPPED"], "max_document_bytes": 1024}


def request() -> dict:
    return {"schema": "elite-amazon-spapi-fulfillment-delivery-evidence-request/v1", "approved_by": "owner@example.invalid", "approved_at": "2026-08-31T18:59:00+00:00", "business_order_receipt_sha256": sha256(BUSINESS_RECEIPT), "seller_fulfillment_order_id": "seller-order-1", "marketplace_id": "A1MARKET", "expected_packages": [{"amazon_shipment_id": "shipment-1", "package_number": 1, "required_document_types": ["PHOTO", "SIGNATURE"]}]}


def response() -> dict:
    return {"payload": {"fulfillment_order": {"seller_fulfillment_order_id": "seller-order-1", "marketplace_id": "A1MARKET", "destination_address": {"name": "PRIVATE PERSON", "address_line1": "PRIVATE STREET"}}, "fulfillment_shipments": [{"amazon_shipment_id": "shipment-1", "fulfillment_shipment_status": "SHIPPED", "fulfillment_shipment_package": [{"package_number": 1, "tracking_number": "PRIVATE-TRACKING", "delivery_information": {"delivery_document_list": [{"document_type": "PHOTO", "url": PHOTO_URL}, {"document_type": "SIGNATURE", "url": SIGNATURE_URL}], "drop_off_location": {"type": "RECEPTIONIST", "attributes": {"recipientName": "PRIVATE NAME"}}}}]}]}, "errors": None}


class FakeService:
    def __init__(self, payload=None):
        self.payload = payload or response()
        self.calls = []

    def get_fulfillment_order(self, seller_fulfillment_order_id, **kwargs):
        self.calls.append((seller_fulfillment_order_id, kwargs))
        return self.payload


def fetcher(url: str, maximum: int, timeout: float):
    if url == PHOTO_URL:
        return b"\xff\xd8\xffphoto", "image/jpeg"
    if url == SIGNATURE_URL:
        return b"\x89PNG\r\n\x1a\nsignature", "image/png"
    raise AssertionError("unexpected URL")


class FulfillmentDeliveryEvidenceTests(unittest.TestCase):
    def test_distributed_profile_is_blocked_by_default(self):
        with self.assertRaises(PermissionError):
            approved_profile(Path(__file__).with_name("provider-profile.template.json"))

    def test_official_sdk_signature_and_delivery_models(self):
        self.assertEqual(list(inspect.signature(FbaOutboundApi.get_fulfillment_order).parameters)[:2], ["self", "seller_fulfillment_order_id"])
        document = DeliveryDocument(document_type="PHOTO", url=PHOTO_URL)
        information = DeliveryInformation(delivery_document_list=[document])
        package = FulfillmentShipmentPackage(package_number=1, carrier_code="ILLUSTRATIVE_CARRIER", delivery_information=information)
        self.assertEqual(package.delivery_information.delivery_document_list[0].document_type, "PHOTO")

    def test_default_fetcher_redacts_signed_url_on_network_failure(self):
        class BrokenOpener:
            def open(self, request, timeout):
                raise urllib.error.URLError(f"failed URL {request.full_url}")

        with patch("urllib.request.build_opener", return_value=BrokenOpener()):
            with self.assertRaises(RuntimeError) as captured:
                default_fetcher(PHOTO_URL, 1024, 1.0)
        self.assertNotIn("signature=secret", str(captured.exception))
        self.assertIsNone(captured.exception.__cause__)

    def test_downloads_hash_named_documents_without_urls_or_identity(self):
        with tempfile.TemporaryDirectory() as root:
            target = Path(root) / "evidence"
            receipt = fetch_delivery_evidence(FakeService(), request(), BUSINESS_RECEIPT, profile(), target, fetcher, now=NOW)
            self.assertEqual(len(receipt["documents"]), 2)
            self.assertFalse(receipt["automatic_business_delivery_acceptance"])
            self.assertFalse(receipt["recipient_identity_values_persisted"])
            for document in receipt["documents"]:
                self.assertTrue((target / document["filename"]).is_file())
            persisted = "".join(path.read_text(errors="ignore") for path in target.glob("*.json"))
            self.assertNotIn("signature=secret", persisted)
            self.assertNotIn("PRIVATE PERSON", persisted)
            self.assertNotIn("PRIVATE NAME", persisted)
            self.assertNotIn("PRIVATE-TRACKING", persisted)

    def test_business_receipt_tamper_blocks_before_provider(self):
        with tempfile.TemporaryDirectory() as root:
            service = FakeService()
            target = Path(root) / "evidence"
            with self.assertRaises(PermissionError):
                fetch_delivery_evidence(service, request(), b"tampered", profile(), target, fetcher, now=NOW)
            self.assertEqual(service.calls, [])
            self.assertFalse(target.exists())

    def test_wrong_order_or_marketplace_blocks_before_download(self):
        with tempfile.TemporaryDirectory() as root:
            bad = response()
            bad["payload"]["fulfillment_order"]["seller_fulfillment_order_id"] = "other"
            calls = []
            with self.assertRaises(ValueError):
                fetch_delivery_evidence(FakeService(bad), request(), BUSINESS_RECEIPT, profile(), Path(root) / "evidence", lambda *args: calls.append(args), now=NOW)
            self.assertEqual(calls, [])

    def test_missing_package_or_required_document_fails_closed(self):
        with tempfile.TemporaryDirectory() as root:
            bad = response()
            bad["payload"]["fulfillment_shipments"][0]["fulfillment_shipment_package"][0]["delivery_information"]["delivery_document_list"] = []
            with self.assertRaises(ValueError):
                fetch_delivery_evidence(FakeService(bad), request(), BUSINESS_RECEIPT, profile(), Path(root) / "evidence", fetcher, now=NOW)

    def test_unapproved_host_or_nonterminal_status_fails_before_download(self):
        with tempfile.TemporaryDirectory() as root:
            bad = response()
            bad["payload"]["fulfillment_shipments"][0]["fulfillment_shipment_package"][0]["delivery_information"]["delivery_document_list"][0]["url"] = "https://evil.example/proof"
            with self.assertRaises(PermissionError):
                fetch_delivery_evidence(FakeService(bad), request(), BUSINESS_RECEIPT, profile(), Path(root) / "evidence", fetcher, now=NOW)

    def test_bad_magic_or_oversize_leaves_no_output(self):
        with tempfile.TemporaryDirectory() as root:
            target = Path(root) / "evidence"
            def bad_fetch(url, maximum, timeout):
                return b"not-an-image", "image/jpeg"
            with self.assertRaises(ValueError):
                fetch_delivery_evidence(FakeService(), request(), BUSINESS_RECEIPT, profile(), target, bad_fetch, now=NOW)
            self.assertFalse(target.exists())

    def test_duplicate_document_content_is_rejected(self):
        with tempfile.TemporaryDirectory() as root:
            target = Path(root) / "evidence"
            def duplicate(url, maximum, timeout):
                return b"\xff\xd8\xffsame", "image/jpeg"
            with self.assertRaises(ValueError):
                fetch_delivery_evidence(FakeService(), request(), BUSINESS_RECEIPT, profile(), target, duplicate, now=NOW)
            self.assertFalse(target.exists())

    def test_existing_output_is_never_overwritten_or_called(self):
        with tempfile.TemporaryDirectory() as root:
            target = Path(root) / "evidence"
            target.mkdir()
            marker = target / "keep"
            marker.write_text("keep")
            service = FakeService()
            with self.assertRaises(FileExistsError):
                fetch_delivery_evidence(service, request(), BUSINESS_RECEIPT, profile(), target, fetcher, now=NOW)
            self.assertEqual(service.calls, [])
            self.assertEqual(marker.read_text(), "keep")

    def test_output_created_during_download_is_not_overwritten(self):
        with tempfile.TemporaryDirectory() as root:
            target = Path(root) / "evidence"

            def racing_fetcher(url, maximum, timeout):
                if not target.exists():
                    target.mkdir()
                    (target / "intruder").write_text("keep")
                return fetcher(url, maximum, timeout)

            with self.assertRaises(FileExistsError):
                fetch_delivery_evidence(FakeService(), request(), BUSINESS_RECEIPT, profile(), target, racing_fetcher, now=NOW)
            self.assertEqual((target / "intruder").read_text(), "keep")
            self.assertEqual(list(Path(root).glob(".evidence.tmp-*")), [])


if __name__ == "__main__":
    unittest.main()
````

### FILE: `amazon_spapi_fulfillment_delivery_evidence/README.md`
```yaml
block_id: "PY-AMAZON-SPAPI-FULFILLMENT-DELIVERY-EVIDENCE:README-MD:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed wrapper/configuration governed by official Amazon API and models"
license: "LicenseRef-Workspace-Owner"
sha256: "b260c8edd79b3baf148f366e244cb6b55a2276f1ec857631b31f78f0e1272996"
variables: []
secrets_allowed: false
```
````markdown
# Amazon SP-API Fulfillment delivery evidence

Wrapper local fail-closed sobre el SDK oficial Amazon `amzn-sp-api==1.11.1` y Fulfillment Outbound v2020-07-01.

`getFulfillmentOrder` liga la consulta a un receipt empresarial y paquetes esperados. Sólo descarga tipos/hosts aprobados; las URLs firmadas y valores de identidad del destinatario nunca se persisten. Cada foto, firma o PDF se valida por tamaño, content type y magic bytes, y se guarda con nombre SHA-256. El receipt siempre conserva `automatic_business_delivery_acceptance=false`.

El perfil distribuido está bloqueado. Antes de una llamada real deben probarse registro, rol Amazon Fulfillment, programa MCF/marketplace, sandbox, cuota/costo, tipos/hosts observados, privacidad, retención y owner. No es carrier universal y la presencia de un documento no sustituye aceptación contractual ni inspección del bien.

```powershell
python -m pip install --require-hashes -r requirements.lock
python -m unittest -v test_fulfillment_delivery_evidence.py
```
````

## 6. Configuration surface

El perfil distribuido queda bloqueado. El usuario debe demostrar registro de developer y aplicación, rol/scope Amazon Fulfillment, soporte MCF/marketplace, sandbox, cuota/costo, tratamiento y retención de evidencia, owner, marketplaces, status terminales, tipos documentales, hosts y content types observados. Los secretos se referencian exclusivamente por nombres de variables de entorno.

## 7. Dependency bill

`requirements.lock` fija `amzn-sp-api==1.11.1` y todo su grafo por URL oficial PyPI y SHA-256. `sdk-artifact.lock.json` fija wheel, commit de SDK, commit de modelos, hashes y Apache-2.0. Los ocho archivos del pack son `AUTHORED`; no se presentan como código copiado de Amazon.

## 8. Apply order

1. Materializar junto a `OFFICIAL-UPSTREAM-ACQUISITION-CORE`. 2. Instalar con `python -m pip install --require-hashes -r requirements.lock`. 3. Probar y aprobar el perfil sin incluir secretos. 4. Generar request desde receipt empresarial real y paquetes esperados. 5. Ejecutar contra sandbox. 6. Revisar receipt y binarios hash-named. 7. Integrar aceptación empresarial como decisión separada.

## 9. Verification

Debe pasar materialización 8/8, igualdad SHA por archivo, `pip check` y las 12 pruebas: firma/modelos oficiales, perfil bloqueado, binding del receipt, identidad pedido/marketplace, paquete/tipo exactos, host/status, bytes/magic, duplicados, no-overwrite inicial y concurrente, redacción de URL firmada y ausencia de identidad/URL/tracking persistidos. Un proveedor vivo y aceptación empresarial permanecen gates del proyecto.

## 10. Reconstruction evidence

`reconstruction_evidence/AMAZON_SPAPI_FULFILLMENT_DELIVERY_EVIDENCE_EXECUTION_INVENTORY_2026-08-31_V154.md` registra hashes, comandos y resultados del árbol autor y del árbol reconstruido. La evidencia demuestra el adapter local y el SDK fijado; no demuestra credenciales, disponibilidad provider, entrega física ni producción.
