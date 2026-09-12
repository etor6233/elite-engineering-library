# Python Amazon SP-API Supply Sources Adapter

## 1. Metadata

```yaml
pack_id: "PYTHON-AMAZON-SPAPI-SUPPLY-SOURCES-ADAPTER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa inventario y ciclo create/update/status/archive de Amazon Supply Sources v2020-07-01 mediante el SDK oficial amzn-sp-api 1.11.1; no inventa multi-location universal ni activa ubicaciones internas."
stacks: ["Python 3.10-3.14", "amzn-sp-api 1.11.1", "Amazon Supply Sources v2020-07-01"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "PYTHON-AMAZON-SPAPI-CATALOG-ADAPTER 0.1.x", "PYTHON-AMAZON-SPAPI-SHIPPING-TRACKING-ADAPTER 0.6.x", "PYTHON-AMAZON-SPAPI-FULFILLMENT-DELIVERY-EVIDENCE-ADAPTER 0.1.x"]
incompatible_with: ["red multi-location universal", "marketplace/programa no demostrado", "mutación sin inventario y snapshot hash-bound", "retry automático de efecto ambiguo", "archive de source Active", "persistencia de dirección/contacto en receipts", "activación interna automática", "secreto en CLI/Markdown", "SDK distinto de 1.11.1"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0 AND MIT AND MPL-2.0 AND BSD-3-Clause"
upstream_sources: ["https://github.com/amzn/selling-partner-api-sdk/tree/8e792ae345a8d334ccdbdd03181f05f040e6a4fc", "https://github.com/amzn/selling-partner-api-models/tree/8e429486005c4ebdce5099e48cc48515a65359bb", "https://developer-docs.amazon.com/sp-api/docs/supply-sources-api-v2020-07-01-reference", "https://developer-docs.amazon.com/sp-api/docs/sp-api-release-notes?ld=SDESSOADirect"]
verified_at: "2026-08-31"
```

## 2. Applicability

Use sólo para supply sources Amazon —tiendas o depósitos registrados para capacidades de fulfillment— en cuenta, marketplace y programa demostrados. Amazon publica esta API para configurar información y capacidades por ubicación; este pack no sustituye el modelo interno de sucursales, inventario, promesas, transporte ni conciliación.

## 3. Architecture contract

Inventario recorre páginas con tokens únicos, obtiene cada source completo y rechaza IDs/códigos duplicados o divergentes. Persiste ID/código/alias/status y hashes, nunca dirección/contacto. Cada CREATE/UPDATE/STATUS/ARCHIVE exige approval ligado a request, receipt e instantánea exacta, vuelve a leer inventario para detectar drift, escribe attempt antes del efecto, usa modelos/métodos oficiales y reconcilia con `get_supply_source`. Toda excepción posterior al inicio crea `UNKNOWN_EFFECT` y bloquea retry. Archive exige status `Inactive`; ningún status Amazon activa automáticamente una ubicación interna.

## 4. Exact file manifest

```text
CREATE amazon_spapi_supply_sources/requirements-direct.in
CREATE amazon_spapi_supply_sources/requirements.lock
CREATE amazon_spapi_supply_sources/sdk-artifact.lock.json
CREATE amazon_spapi_supply_sources/provider-profile.template.json
CREATE amazon_spapi_supply_sources/supply-source-request.template.json
CREATE amazon_spapi_supply_sources/mutation-approval.template.json
CREATE amazon_spapi_supply_sources/run_supply_sources.py
CREATE amazon_spapi_supply_sources/test_supply_sources.py
CREATE amazon_spapi_supply_sources/README.md
```

## 5. Materialization blocks

### FILE: `amazon_spapi_supply_sources/requirements-direct.in`
```yaml
block_id: "PY-AMAZON-SPAPI-SUPPLY-SOURCES:REQUIREMENTS-DIRECT-IN:v1"
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

### FILE: `amazon_spapi_supply_sources/requirements.lock`
```yaml
block_id: "PY-AMAZON-SPAPI-SUPPLY-SOURCES:REQUIREMENTS-LOCK:v1"
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

### FILE: `amazon_spapi_supply_sources/sdk-artifact.lock.json`
```yaml
block_id: "PY-AMAZON-SPAPI-SUPPLY-SOURCES:SDK-ARTIFACT-LOCK-JSON:v1"
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

### FILE: `amazon_spapi_supply_sources/provider-profile.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-SUPPLY-SOURCES:PROVIDER-PROFILE-TEMPLATE-JSON:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration governed by official Amazon Supply Sources API and models"
license: "LicenseRef-Workspace-Owner"
sha256: "28770811a55173a44d422d171fa337745603c55538c0331b1a2d7eab75e3fc73"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-supply-sources-profile/v1",
  "provider": "Amazon Supply Sources v2020-07-01",
  "sdk": "amzn-sp-api==1.11.1",
  "decision": "BLOCKED_ACCESS_SCOPE_AND_LOCATION_POLICY_REQUIRED",
  "region": "SANDBOX",
  "approved_operations": [],
  "page_size": 50,
  "max_pages": 20,
  "max_sources": 500,
  "developer_registration": "NOT_PROVEN",
  "application_registration": "NOT_PROVEN",
  "supply_sources_role_and_scope": "NOT_PROVEN",
  "marketplace_and_program_support": "NOT_PROVEN",
  "dynamic_sandbox_contract_test": "NOT_PROVEN",
  "quota_and_cost_approved": false,
  "business_location_data_handling_approved": false,
  "contact_data_processing_approved": false,
  "location_lifecycle_owner": "",
  "archive_requires_inactive": true,
  "automatic_internal_location_activation": false,
  "client_id_env": "AMAZON_SP_API_CLIENT_ID",
  "client_secret_env": "AMAZON_SP_API_CLIENT_SECRET",
  "refresh_token_env": "AMAZON_SP_API_REFRESH_TOKEN"
}
````

### FILE: `amazon_spapi_supply_sources/supply-source-request.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-SUPPLY-SOURCES:SUPPLY-SOURCE-REQUEST-TEMPLATE-JSON:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration governed by official Amazon Supply Sources API and models"
license: "LicenseRef-Workspace-Owner"
sha256: "f81531c45b9d8e0ca70155a1aab7d68224bea05f94e003a4e6989a91c9e775ca"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-supply-source-mutation/v1",
  "operation": "CREATE|UPDATE|STATUS|ARCHIVE",
  "supply_source_id": null,
  "payload": {}
}
````

### FILE: `amazon_spapi_supply_sources/mutation-approval.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-SUPPLY-SOURCES:MUTATION-APPROVAL-TEMPLATE-JSON:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration governed by official Amazon Supply Sources API and models"
license: "LicenseRef-Workspace-Owner"
sha256: "8febfd7f27d50beaca0c45a391b6eaef690a4bf017d609ea9d43ff6000ee4db3"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-supply-source-approval/v1",
  "approved_by": "",
  "approved_at": "",
  "request_sha256": "",
  "inventory_receipt_sha256": "",
  "expected_source_snapshot_sha256": null,
  "allow_unknown_effect_retry": false
}
````

### FILE: `amazon_spapi_supply_sources/run_supply_sources.py`
```yaml
block_id: "PY-AMAZON-SPAPI-SUPPLY-SOURCES:RUN-SUPPLY-SOURCES-PY:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration governed by official Amazon Supply Sources API and models"
license: "LicenseRef-Workspace-Owner"
sha256: "02a45d64d6b4164a0d055d3c79e6e69209256adae6b5a53f2e1782c043760f17"
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
from typing import Any, Protocol
import uuid

from spapi import SPAPIClient, SPAPIConfig
from spapi.api.supply_sources_v2020_07_01.supply_sources_api import SupplySourcesApi
import spapi.models.supply_sources_v2020_07_01 as official_models

SDK_VERSION = "1.11.1"
OPERATIONS = {"CREATE", "UPDATE", "STATUS", "ARCHIVE"}
REGIONS = {"SANDBOX", "NA", "EU", "FE"}
REQUEST_KEYS = {"schema", "operation", "supply_source_id", "payload"}
APPROVAL_KEYS = {"schema", "approved_by", "approved_at", "request_sha256", "inventory_receipt_sha256", "expected_source_snapshot_sha256", "allow_unknown_effect_retry"}
REQUIRED_MODEL_FIELDS = {
    "CreateSupplySourceRequest": {"supply_source_code", "alias", "address"},
    "Address": {"name", "address_line1", "state_or_region", "country_code"},
    "ThroughputConfig": {"throughput_unit"},
}


class SupplySourcesService(Protocol):
    def get_supply_sources(self, **kwargs: Any) -> Any: ...
    def get_supply_source(self, supply_source_id: str, **kwargs: Any) -> Any: ...
    def create_supply_source(self, payload: Any, **kwargs: Any) -> Any: ...
    def update_supply_source(self, supply_source_id: str, **kwargs: Any) -> Any: ...
    def update_supply_source_status(self, supply_source_id: str, **kwargs: Any) -> Any: ...
    def archive_supply_source(self, supply_source_id: str, **kwargs: Any) -> Any: ...


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


def prune_none(value: Any) -> Any:
    if isinstance(value, dict):
        return {key: prune_none(item) for key, item in value.items() if item is not None}
    if isinstance(value, list):
        return [prune_none(item) for item in value]
    return value


def canonical_bytes(value: Any) -> bytes:
    return (json.dumps(json_value(value), ensure_ascii=False, sort_keys=True, separators=(",", ":")) + "\n").encode("utf-8")


def pretty_bytes(value: Any) -> bytes:
    return (json.dumps(json_value(value), ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")


def nonempty(value: Any, label: str, maximum: int = 255) -> str:
    if not isinstance(value, str) or not value or value != value.strip() or len(value) > maximum:
        raise ValueError(f"{label} must be a trimmed non-empty string up to {maximum} characters")
    return value


def strict_object(path: Path, label: str) -> tuple[dict[str, Any], bytes]:
    raw = path.read_bytes()
    value = json.loads(raw, object_pairs_hook=lambda pairs: _unique_object(pairs, label))
    if not isinstance(value, dict):
        raise ValueError(f"{label} must be an object")
    return value, raw


def _unique_object(pairs: list[tuple[str, Any]], label: str) -> dict[str, Any]:
    result: dict[str, Any] = {}
    for key, value in pairs:
        if key in result:
            raise ValueError(f"duplicate key in {label}: {key}")
        result[key] = value
    return result


def exact_keys(value: Any, expected: set[str], label: str) -> dict[str, Any]:
    if not isinstance(value, dict) or set(value) != expected:
        raise ValueError(f"{label} keys must be exactly {sorted(expected)}")
    return value


def aware_datetime(value: Any, label: str) -> datetime:
    text = nonempty(value, label, 64)
    try:
        parsed = datetime.fromisoformat(text.replace("Z", "+00:00"))
    except ValueError as error:
        raise ValueError(f"{label} must be ISO-8601") from error
    if parsed.tzinfo is None:
        raise ValueError(f"{label} must include an offset")
    return parsed.astimezone(timezone.utc)


def approved_profile(path: Path) -> dict[str, Any]:
    profile, _ = strict_object(path, "provider profile")
    required = {"schema", "provider", "sdk", "decision", "region", "approved_operations", "page_size", "max_pages", "max_sources", "developer_registration", "application_registration", "supply_sources_role_and_scope", "marketplace_and_program_support", "dynamic_sandbox_contract_test", "quota_and_cost_approved", "business_location_data_handling_approved", "contact_data_processing_approved", "location_lifecycle_owner", "archive_requires_inactive", "automatic_internal_location_activation", "client_id_env", "client_secret_env", "refresh_token_env"}
    exact_keys(profile, required, "provider profile")
    if profile["schema"] != "elite-amazon-spapi-supply-sources-profile/v1" or profile["provider"] != "Amazon Supply Sources v2020-07-01" or profile["sdk"] != "amzn-sp-api==1.11.1" or profile["decision"] != "APPROVED":
        raise PermissionError("Supply Sources profile is not approved")
    if profile["region"] not in REGIONS:
        raise PermissionError("region is not approved")
    operations = profile["approved_operations"]
    if not isinstance(operations, list) or not operations or len(operations) != len(set(operations)) or any(item not in OPERATIONS for item in operations):
        raise ValueError("approved_operations must be a unique non-empty official operation subset")
    for key in ("page_size", "max_pages", "max_sources"):
        if not isinstance(profile[key], int) or isinstance(profile[key], bool) or profile[key] < 1:
            raise ValueError(f"{key} must be a positive integer")
    if profile["page_size"] > 100 or profile["max_pages"] > 100 or profile["max_sources"] > 10_000:
        raise ValueError("pagination bounds exceed local safety limits")
    for key in ("developer_registration", "application_registration", "supply_sources_role_and_scope", "marketplace_and_program_support", "dynamic_sandbox_contract_test"):
        if profile[key] != "PROVEN":
            raise PermissionError(f"{key} must be PROVEN")
    for key in ("quota_and_cost_approved", "business_location_data_handling_approved", "contact_data_processing_approved", "archive_requires_inactive"):
        if profile[key] is not True:
            raise PermissionError(f"{key} must be true")
    if profile["automatic_internal_location_activation"] is not False:
        raise PermissionError("automatic internal location activation must remain false")
    nonempty(profile["location_lifecycle_owner"], "location_lifecycle_owner")
    for key in ("client_id_env", "client_secret_env", "refresh_token_env"):
        if not re.fullmatch(r"[A-Z][A-Z0-9_]{2,127}", nonempty(profile[key], key, 128)):
            raise ValueError(f"{key} must name an environment variable")
    return profile


def enum_values(class_name: str) -> set[str]:
    cls = getattr(official_models, class_name)
    return {value for key, value in vars(cls).items() if key.isupper() and isinstance(value, str)}


def build_official_value(type_name: str, value: Any, label: str) -> Any:
    if type_name == "str":
        return nonempty(value, label, 4096)
    if type_name == "bool":
        if not isinstance(value, bool):
            raise ValueError(f"{label} must be a boolean")
        return value
    if type_name == "int":
        if not isinstance(value, int) or isinstance(value, bool) or value < 0:
            raise ValueError(f"{label} must be a non-negative integer")
        return value
    list_match = re.fullmatch(r"List\[(?P<inner>[A-Za-z0-9_]+)\]", type_name)
    if list_match:
        if not isinstance(value, list) or not value or len(value) > 100:
            raise ValueError(f"{label} must be a non-empty list up to 100 items")
        return [build_official_value(list_match.group("inner"), item, f"{label}[{index}]") for index, item in enumerate(value)]
    values = enum_values(type_name) if hasattr(official_models, type_name) else set()
    cls = getattr(official_models, type_name, None)
    if values and getattr(cls, "swagger_types", {}) == {}:
        if value not in values:
            raise ValueError(f"{label} must be one of {sorted(values)}")
        return value
    if not isinstance(cls, type) or not isinstance(value, dict) or not value:
        raise ValueError(f"{label} must be a non-empty official {type_name} object")
    allowed = set(cls.swagger_types)
    unknown = set(value) - allowed
    missing = REQUIRED_MODEL_FIELDS.get(type_name, set()) - set(value)
    if unknown or missing:
        raise ValueError(f"{label} has unknown={sorted(unknown)} missing={sorted(missing)}")
    kwargs = {key: build_official_value(cls.swagger_types[key], item, f"{label}.{key}") for key, item in value.items()}
    return cls(**kwargs)


def provider_dict(value: Any, label: str) -> dict[str, Any]:
    data = json_value(value)
    if not isinstance(data, dict) or data.get("errors") not in (None, []):
        raise ValueError(f"{label} is invalid or contains provider errors")
    return data


def filtered_source(value: Any) -> dict[str, Any]:
    source = provider_dict(value, "supply source")
    source_id = nonempty(source.get("supply_source_id"), "provider supply_source_id")
    code = nonempty(source.get("supply_source_code"), "provider supply_source_code")
    alias = nonempty(source.get("alias"), "provider alias")
    status = source.get("status")
    if status not in enum_values("SupplySourceStatusReadOnly"):
        raise ValueError("provider supply source status is invalid")
    address = source.get("address")
    if not isinstance(address, dict):
        raise ValueError("provider supply source address is missing")
    configuration = source.get("configuration")
    capabilities = source.get("capabilities")
    return {
        "supply_source_id": source_id,
        "supply_source_code": code,
        "alias": alias,
        "status": status,
        "address_sha256": sha256(canonical_bytes(prune_none(address))),
        "configuration_sha256": sha256(canonical_bytes(prune_none(configuration))),
        "capabilities_sha256": sha256(canonical_bytes(prune_none(capabilities))),
        "created_at": source.get("created_at"),
        "updated_at": source.get("updated_at"),
        "address_or_contact_values_persisted": False,
    }


def collect_inventory(service: SupplySourcesService, profile: dict[str, Any], timeout_seconds: float = 30.0) -> tuple[list[dict[str, Any]], list[str]]:
    if timeout_seconds <= 0 or timeout_seconds > 120:
        raise ValueError("timeout must be within 1..120 seconds")
    token: str | None = None
    seen_tokens: set[str] = set()
    summaries: list[dict[str, Any]] = []
    page_hashes: list[str] = []
    for _ in range(profile["max_pages"]):
        kwargs: dict[str, Any] = {"page_size": profile["page_size"], "_request_timeout": timeout_seconds}
        if token is not None:
            kwargs["next_page_token"] = token
        page = provider_dict(service.get_supply_sources(**kwargs), "supply source page")
        page_hashes.append(sha256(canonical_bytes(page)))
        entries = page.get("supply_sources")
        if not isinstance(entries, list):
            raise ValueError("provider supply_sources must be a list")
        summaries.extend(entries)
        if len(summaries) > profile["max_sources"]:
            raise ValueError("provider inventory exceeds max_sources")
        next_token = page.get("next_page_token")
        if next_token in (None, ""):
            break
        token = nonempty(next_token, "next_page_token", 4096)
        if token in seen_tokens:
            raise ValueError("provider pagination token cycle detected")
        seen_tokens.add(token)
    else:
        raise ValueError("provider inventory exceeds max_pages")
    sources: list[dict[str, Any]] = []
    seen_ids: set[str] = set()
    seen_codes: set[str] = set()
    for summary in summaries:
        if not isinstance(summary, dict):
            raise ValueError("provider supply source summary is invalid")
        source_id = nonempty(summary.get("supply_source_id"), "summary supply_source_id")
        code = nonempty(summary.get("supply_source_code"), "summary supply_source_code")
        if source_id in seen_ids or code in seen_codes:
            raise ValueError("provider inventory contains duplicate source identity or code")
        seen_ids.add(source_id)
        seen_codes.add(code)
        source = filtered_source(service.get_supply_source(source_id, _request_timeout=timeout_seconds))
        if source["supply_source_id"] != source_id or source["supply_source_code"] != code:
            raise ValueError("provider summary and full source identity differ")
        sources.append(source)
    return sorted(sources, key=lambda item: item["supply_source_id"]), page_hashes


def inventory_receipt(service: SupplySourcesService, profile: dict[str, Any], output_directory: Path, timeout_seconds: float = 30.0, now: datetime | None = None) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION:
        raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    target = output_directory.resolve()
    if target.exists():
        raise FileExistsError("output_directory must not exist")
    sources, page_hashes = collect_inventory(service, profile, timeout_seconds)
    current = (now or datetime.now(timezone.utc)).astimezone(timezone.utc)
    receipt = {"schema": "elite-amazon-spapi-supply-sources-inventory/v1", "created_at": current.isoformat(), "provider": "Amazon Supply Sources v2020-07-01", "sdk_version": SDK_VERSION, "page_response_sha256": page_hashes, "source_count": len(sources), "sources": sources, "inventory_sha256": sha256(canonical_bytes(sources)), "address_or_contact_values_persisted": False, "automatic_internal_location_activation": False}
    stage = target.parent / f".{target.name}.tmp-{uuid.uuid4().hex}"
    if target.exists() or stage.exists():
        raise FileExistsError("inventory output appeared during collection")
    stage.mkdir(parents=True)
    try:
        (stage / "SUPPLY_SOURCES_INVENTORY.json").write_bytes(pretty_bytes(receipt))
        if target.exists():
            raise FileExistsError("inventory output appeared during staging")
        stage.replace(target)
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise
    return receipt


def validate_mutation(request: dict[str, Any], request_bytes: bytes, approval: dict[str, Any], inventory: dict[str, Any], inventory_bytes: bytes, profile: dict[str, Any], now: datetime) -> tuple[str, Any, dict[str, Any] | None]:
    exact_keys(request, REQUEST_KEYS, "mutation request")
    exact_keys(approval, APPROVAL_KEYS, "mutation approval")
    if request["schema"] != "elite-amazon-spapi-supply-source-mutation/v1" or approval["schema"] != "elite-amazon-spapi-supply-source-approval/v1":
        raise ValueError("request or approval schema is invalid")
    operation = request["operation"]
    if operation not in profile["approved_operations"]:
        raise PermissionError("operation is not approved")
    nonempty(approval["approved_by"], "approved_by")
    if aware_datetime(approval["approved_at"], "approved_at") > now:
        raise ValueError("approved_at must not be in the future")
    if approval["request_sha256"] != sha256(request_bytes) or approval["inventory_receipt_sha256"] != sha256(inventory_bytes):
        raise PermissionError("approval is not bound to request and inventory receipt")
    if approval["allow_unknown_effect_retry"] is not False:
        raise PermissionError("unknown effects must never be retried automatically")
    if inventory.get("schema") != "elite-amazon-spapi-supply-sources-inventory/v1" or inventory.get("inventory_sha256") != sha256(canonical_bytes(inventory.get("sources"))):
        raise ValueError("inventory receipt is invalid")
    sources = inventory["sources"]
    if not isinstance(sources, list):
        raise ValueError("inventory sources are invalid")
    source_id = request["supply_source_id"]
    current_source: dict[str, Any] | None = None
    if operation == "CREATE":
        if source_id is not None or approval["expected_source_snapshot_sha256"] is not None:
            raise ValueError("CREATE must not provide existing source identity")
        model = build_official_value("CreateSupplySourceRequest", request["payload"], "payload")
        code = model.supply_source_code
        if any(item.get("supply_source_code") == code for item in sources):
            raise ValueError("supply source code already exists")
        return operation, model, None
    source_id = nonempty(source_id, "supply_source_id")
    matches = [item for item in sources if item.get("supply_source_id") == source_id]
    if len(matches) != 1:
        raise ValueError("existing supply source is absent or duplicated")
    current_source = matches[0]
    if approval["expected_source_snapshot_sha256"] != sha256(canonical_bytes(current_source)):
        raise PermissionError("approval is not bound to the existing source snapshot")
    if current_source.get("status") == "Archived":
        raise PermissionError("archived supply source cannot be mutated")
    if operation == "UPDATE":
        model = build_official_value("UpdateSupplySourceRequest", request["payload"], "payload")
    elif operation == "STATUS":
        exact_keys(request["payload"], {"status"}, "status payload")
        status = build_official_value("SupplySourceStatus", request["payload"]["status"], "payload.status")
        model = official_models.UpdateSupplySourceStatusRequest(status=status)
    elif operation == "ARCHIVE":
        exact_keys(request["payload"], set(), "archive payload")
        if profile["archive_requires_inactive"] and current_source.get("status") != "Inactive":
            raise PermissionError("archive requires an Inactive source")
        model = None
    else:
        raise ValueError("unsupported operation")
    return operation, model, current_source


def deep_contains(actual: Any, expected: Any) -> bool:
    if isinstance(expected, dict):
        return isinstance(actual, dict) and all(key in actual and deep_contains(actual[key], item) for key, item in expected.items())
    if isinstance(expected, list):
        return isinstance(actual, list) and len(actual) == len(expected) and all(deep_contains(a, e) for a, e in zip(actual, expected))
    return actual == expected


def exclusive_write(path: Path, value: Any) -> None:
    with path.open("xb") as stream:
        stream.write(pretty_bytes(value))


def execute_mutation(service: SupplySourcesService, request: dict[str, Any], request_bytes: bytes, approval: dict[str, Any], inventory: dict[str, Any], inventory_bytes: bytes, profile: dict[str, Any], output_directory: Path, timeout_seconds: float = 30.0, now: datetime | None = None) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION:
        raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if timeout_seconds <= 0 or timeout_seconds > 120:
        raise ValueError("timeout must be within 1..120 seconds")
    target = output_directory.resolve()
    if target.exists():
        raise FileExistsError("output_directory must not exist")
    current_time = (now or datetime.now(timezone.utc)).astimezone(timezone.utc)
    operation, model, previous = validate_mutation(request, request_bytes, approval, inventory, inventory_bytes, profile, current_time)
    live_sources, _ = collect_inventory(service, profile, timeout_seconds)
    if sha256(canonical_bytes(live_sources)) != inventory["inventory_sha256"]:
        raise PermissionError("live supply source inventory drifted from approved receipt")
    target.mkdir(parents=True)
    attempt = {"schema": "elite-amazon-spapi-supply-source-attempt/v1", "created_at": current_time.isoformat(), "operation": operation, "request_sha256": sha256(request_bytes), "inventory_receipt_sha256": sha256(inventory_bytes), "previous_source_snapshot_sha256": sha256(canonical_bytes(previous)) if previous else None, "status": "PENDING", "automatic_retry_authorized": False}
    exclusive_write(target / "MUTATION_ATTEMPT.json", attempt)
    effect_started = False
    try:
        source_id = request["supply_source_id"]
        effect_started = True
        if operation == "CREATE":
            response = provider_dict(service.create_supply_source(model, _request_timeout=timeout_seconds), "create response")
            source_id = nonempty(response.get("supply_source_id"), "created supply_source_id")
            if response.get("supply_source_code") != model.supply_source_code:
                raise ValueError("create response code does not match request")
        elif operation == "UPDATE":
            service.update_supply_source(source_id, payload=model, _request_timeout=timeout_seconds)
        elif operation == "STATUS":
            service.update_supply_source_status(source_id, payload=model, _request_timeout=timeout_seconds)
        else:
            service.archive_supply_source(source_id, _request_timeout=timeout_seconds)
        provider_source = provider_dict(service.get_supply_source(source_id, _request_timeout=timeout_seconds), "reconciled supply source")
        reconciled = filtered_source(provider_source)
        if reconciled["supply_source_id"] != source_id:
            raise ValueError("reconciled source identity differs")
        if operation == "CREATE" and reconciled["supply_source_code"] != model.supply_source_code:
            raise ValueError("created source code was not reconciled")
        if operation == "UPDATE" and not deep_contains(provider_source, prune_none(model.to_dict())):
            raise ValueError("updated source does not contain requested official model fields")
        if operation == "STATUS" and reconciled["status"] != model.status:
            raise ValueError("updated source status was not reconciled")
        if operation == "ARCHIVE" and reconciled["status"] != "Archived":
            raise ValueError("archived source status was not reconciled")
        exclusive_write(target / "RECONCILED_SOURCE.json", reconciled)
        receipt = {"schema": "elite-amazon-spapi-supply-source-mutation-receipt/v1", "completed_at": datetime.now(timezone.utc).isoformat(), "operation": operation, "supply_source_id": source_id, "request_sha256": sha256(request_bytes), "inventory_receipt_sha256": sha256(inventory_bytes), "reconciled_source_sha256": sha256(canonical_bytes(reconciled)), "provider_effect_reconciled": True, "address_or_contact_values_persisted": False, "automatic_internal_location_activation": False, "automatic_retry_authorized": False}
        exclusive_write(target / "MUTATION_RECEIPT.json", receipt)
        return receipt
    except BaseException as error:
        if effect_started:
            unknown = {"schema": "elite-amazon-spapi-supply-source-unknown-effect/v1", "recorded_at": datetime.now(timezone.utc).isoformat(), "operation": operation, "request_sha256": sha256(request_bytes), "error_type": type(error).__name__, "status": "UNKNOWN_EFFECT", "automatic_retry_authorized": False}
            try:
                exclusive_write(target / "UNKNOWN_EFFECT.json", unknown)
            except FileExistsError:
                pass
            raise RuntimeError(f"supply source mutation effect is unknown: {type(error).__name__}") from None
        raise


def service_from_environment(profile: dict[str, Any]) -> SupplySourcesApi:
    values: dict[str, str] = {}
    for key in ("client_id_env", "client_secret_env", "refresh_token_env"):
        value = os.environ.get(profile[key], "")
        if not value:
            raise PermissionError(f"required secret environment variable is missing: {profile[key]}")
        values[key] = value
    config = SPAPIConfig(client_id=values["client_id_env"], client_secret=values["client_secret_env"], refresh_token=values["refresh_token_env"], region=profile["region"])
    return SupplySourcesApi(SPAPIClient(config).api_client)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--mode", required=True, choices=("inventory", "mutate"))
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--request", type=Path)
    parser.add_argument("--approval", type=Path)
    parser.add_argument("--inventory-receipt", type=Path)
    parser.add_argument("--timeout-seconds", type=float, default=30.0)
    args = parser.parse_args()
    profile = approved_profile(args.profile)
    service = service_from_environment(profile)
    if args.mode == "inventory":
        result = inventory_receipt(service, profile, args.output, args.timeout_seconds)
    else:
        if not args.request or not args.approval or not args.inventory_receipt:
            raise ValueError("mutate mode requires request, approval and inventory receipt")
        request, request_bytes = strict_object(args.request, "mutation request")
        approval, _ = strict_object(args.approval, "mutation approval")
        inventory, inventory_bytes = strict_object(args.inventory_receipt, "inventory receipt")
        result = execute_mutation(service, request, request_bytes, approval, inventory, inventory_bytes, profile, args.output, args.timeout_seconds)
    print(json.dumps({"status": "AMAZON_SPAPI_SUPPLY_SOURCES_PASS", "result": result}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `amazon_spapi_supply_sources/test_supply_sources.py`
```yaml
block_id: "PY-AMAZON-SPAPI-SUPPLY-SOURCES:TEST-SUPPLY-SOURCES-PY:v1"
operation: CREATE
provenance: AUTHORED
source: "local executable contract against exact official Supply Sources API/models"
license: "LicenseRef-Workspace-Owner"
sha256: "509640c6e8a553fab4bc2800257693edcf220212adede107baa39fdbbe2893ec"
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

from spapi.api.supply_sources_v2020_07_01.supply_sources_api import SupplySourcesApi
from spapi.models.supply_sources_v2020_07_01.address import Address
from spapi.models.supply_sources_v2020_07_01.create_supply_source_request import CreateSupplySourceRequest
from spapi.models.supply_sources_v2020_07_01.supply_source_status import SupplySourceStatus
from spapi.models.supply_sources_v2020_07_01.update_supply_source_request import UpdateSupplySourceRequest
from spapi.models.supply_sources_v2020_07_01.update_supply_source_status_request import UpdateSupplySourceStatusRequest

from run_supply_sources import approved_profile, canonical_bytes, execute_mutation, inventory_receipt, json_value, pretty_bytes, prune_none, sha256

NOW = datetime(2026, 8, 31, 20, 0, tzinfo=timezone.utc)


def profile() -> dict:
    return {
        "approved_operations": ["CREATE", "UPDATE", "STATUS", "ARCHIVE"],
        "page_size": 50,
        "max_pages": 5,
        "max_sources": 20,
        "archive_requires_inactive": True,
    }


def source(source_id: str = "source-1", code: str = "BRANCH-001", alias: str = "Branch One", status: str = "Active") -> dict:
    return {
        "supply_source_id": source_id,
        "supply_source_code": code,
        "alias": alias,
        "status": status,
        "address": {"name": "PRIVATE BUSINESS", "address_line1": "PRIVATE STREET", "state_or_region": "BA", "country_code": "AR", "phone": "PRIVATE PHONE"},
        "configuration": {"timezone": "America/Argentina/Buenos_Aires"},
        "capabilities": {"outbound": {"is_supported": True}},
        "created_at": "2026-08-30T10:00:00Z",
        "updated_at": "2026-08-30T10:00:00Z",
    }


class FakeService:
    def __init__(self, sources=None):
        items = sources if sources is not None else [source()]
        self.sources = {item["supply_source_id"]: json.loads(json.dumps(item)) for item in items}
        self.calls = []
        self.fail_create = False
        self.ignore_status = False

    def get_supply_sources(self, **kwargs):
        self.calls.append(("list", kwargs))
        return {"supply_sources": [{"supply_source_id": item["supply_source_id"], "supply_source_code": item["supply_source_code"], "alias": item["alias"], "address": item["address"]} for item in self.sources.values()], "next_page_token": None}

    def get_supply_source(self, supply_source_id, **kwargs):
        self.calls.append(("get", supply_source_id, kwargs))
        return json.loads(json.dumps(self.sources[supply_source_id]))

    def create_supply_source(self, payload, **kwargs):
        self.calls.append(("create", payload, kwargs))
        if self.fail_create:
            raise TimeoutError("provider timeout with sensitive request")
        data = prune_none(json_value(payload))
        source_id = "source-created"
        self.sources[source_id] = {"supply_source_id": source_id, "supply_source_code": data["supply_source_code"], "alias": data["alias"], "status": "Inactive", "address": data["address"], "configuration": None, "capabilities": None, "created_at": "2026-08-31T20:00:00Z", "updated_at": "2026-08-31T20:00:00Z"}
        return {"supply_source_id": source_id, "supply_source_code": data["supply_source_code"]}

    def update_supply_source(self, supply_source_id, **kwargs):
        self.calls.append(("update", supply_source_id, kwargs))
        data = prune_none(json_value(kwargs["payload"]))
        self.sources[supply_source_id].update(data)

    def update_supply_source_status(self, supply_source_id, **kwargs):
        self.calls.append(("status", supply_source_id, kwargs))
        if not self.ignore_status:
            self.sources[supply_source_id]["status"] = kwargs["payload"].status

    def archive_supply_source(self, supply_source_id, **kwargs):
        self.calls.append(("archive", supply_source_id, kwargs))
        self.sources[supply_source_id]["status"] = "Archived"


class CyclingService(FakeService):
    def get_supply_sources(self, **kwargs):
        self.calls.append(("list", kwargs))
        return {"supply_sources": [], "next_page_token": "cycle"}


class TwoPageService(FakeService):
    def get_supply_sources(self, **kwargs):
        self.calls.append(("list", kwargs))
        ordered = list(self.sources.values())
        if kwargs.get("next_page_token") is None:
            items, token = ordered[:1], "page-2"
        elif kwargs.get("next_page_token") == "page-2":
            items, token = ordered[1:], None
        else:
            raise AssertionError("unexpected pagination token")
        return {"supply_sources": [{"supply_source_id": item["supply_source_id"], "supply_source_code": item["supply_source_code"], "alias": item["alias"], "address": item["address"]} for item in items], "next_page_token": token}


def create_request() -> dict:
    return {"schema": "elite-amazon-spapi-supply-source-mutation/v1", "operation": "CREATE", "supply_source_id": None, "payload": {"supply_source_code": "BRANCH-NEW", "alias": "New Branch", "address": {"name": "New Branch", "address_line1": "Street 1", "city": "Buenos Aires", "state_or_region": "BA", "postal_code": "1000", "country_code": "AR"}}}


def mutation_request(operation: str, payload: dict, source_id: str = "source-1") -> dict:
    return {"schema": "elite-amazon-spapi-supply-source-mutation/v1", "operation": operation, "supply_source_id": source_id, "payload": payload}


def inventory(service: FakeService, root: Path) -> tuple[dict, bytes]:
    directory = root / "inventory"
    receipt = inventory_receipt(service, profile(), directory, now=NOW)
    return receipt, (directory / "SUPPLY_SOURCES_INVENTORY.json").read_bytes()


def approval(request: dict, request_bytes: bytes, receipt: dict, receipt_bytes: bytes, source_id: str | None) -> dict:
    snapshot = None
    if source_id is not None:
        matches = [item for item in receipt["sources"] if item["supply_source_id"] == source_id]
        snapshot = sha256(canonical_bytes(matches[0]))
    return {"schema": "elite-amazon-spapi-supply-source-approval/v1", "approved_by": "owner@example.invalid", "approved_at": "2026-08-31T19:59:00+00:00", "request_sha256": sha256(request_bytes), "inventory_receipt_sha256": sha256(receipt_bytes), "expected_source_snapshot_sha256": snapshot, "allow_unknown_effect_retry": False}


class SupplySourcesTests(unittest.TestCase):
    def test_official_sdk_signatures_models_and_enums(self):
        signatures = {name: list(inspect.signature(getattr(SupplySourcesApi, name)).parameters)[:2] for name in ("create_supply_source", "get_supply_source", "update_supply_source", "update_supply_source_status", "archive_supply_source")}
        self.assertEqual(signatures["create_supply_source"], ["self", "payload"])
        self.assertEqual(signatures["get_supply_source"], ["self", "supply_source_id"])
        self.assertEqual(signatures["update_supply_source"], ["self", "supply_source_id"])
        self.assertEqual(SupplySourceStatus.ACTIVE, "Active")
        address = Address(name="Branch", address_line1="Street", state_or_region="BA", country_code="AR")
        create = CreateSupplySourceRequest(supply_source_code="B-1", alias="Branch", address=address)
        update = UpdateSupplySourceRequest(alias="Branch 2")
        status = UpdateSupplySourceStatusRequest(status=SupplySourceStatus.INACTIVE)
        self.assertEqual(create.address.country_code, "AR")
        self.assertEqual(update.alias, "Branch 2")
        self.assertEqual(status.status, "Inactive")

    def test_distributed_profile_is_blocked_by_default(self):
        with self.assertRaises(PermissionError):
            approved_profile(Path(__file__).with_name("provider-profile.template.json"))

    def test_inventory_filters_address_contact_and_is_hash_bound(self):
        with tempfile.TemporaryDirectory() as root:
            receipt, raw = inventory(FakeService(), Path(root))
            self.assertEqual(receipt["source_count"], 1)
            self.assertEqual(receipt["inventory_sha256"], sha256(canonical_bytes(receipt["sources"])))
            persisted = raw.decode()
            self.assertNotIn("PRIVATE STREET", persisted)
            self.assertNotIn("PRIVATE PHONE", persisted)
            self.assertIn("address_sha256", persisted)

    def test_inventory_follows_each_unique_page_and_gets_full_sources(self):
        with tempfile.TemporaryDirectory() as root:
            service = TwoPageService([source("source-2", "BRANCH-002", "Branch Two"), source()])
            receipt, _ = inventory(service, Path(root))
            self.assertEqual(receipt["source_count"], 2)
            self.assertEqual(len(receipt["page_response_sha256"]), 2)
            self.assertEqual([item["supply_source_id"] for item in receipt["sources"]], ["source-1", "source-2"])
            self.assertEqual(len([call for call in service.calls if call[0] == "get"]), 2)

    def test_pagination_cycle_and_duplicate_identity_fail_closed(self):
        with tempfile.TemporaryDirectory() as root:
            with self.assertRaises(ValueError):
                inventory_receipt(CyclingService([]), profile(), Path(root) / "cycle", now=NOW)
            duplicates = [source("source-1", "DUP"), source("source-2", "DUP")]
            with self.assertRaises(ValueError):
                inventory_receipt(FakeService(duplicates), profile(), Path(root) / "duplicate", now=NOW)

    def test_create_uses_official_model_attempt_and_reconciliation(self):
        with tempfile.TemporaryDirectory() as root:
            base = Path(root)
            service = FakeService([])
            receipt, receipt_bytes = inventory(service, base)
            request = create_request()
            request_bytes = pretty_bytes(request)
            result = execute_mutation(service, request, request_bytes, approval(request, request_bytes, receipt, receipt_bytes, None), receipt, receipt_bytes, profile(), base / "mutation", now=NOW)
            create_call = [call for call in service.calls if call[0] == "create"]
            self.assertEqual(len(create_call), 1)
            self.assertIsInstance(create_call[0][1], CreateSupplySourceRequest)
            self.assertTrue(result["provider_effect_reconciled"])
            self.assertTrue((base / "mutation/MUTATION_ATTEMPT.json").is_file())
            self.assertTrue((base / "mutation/MUTATION_RECEIPT.json").is_file())
            persisted = "".join(path.read_text() for path in (base / "mutation").glob("*.json"))
            self.assertNotIn("Street 1", persisted)

    def test_update_maps_full_official_nested_model(self):
        with tempfile.TemporaryDirectory() as root:
            base = Path(root)
            service = FakeService()
            receipt, receipt_bytes = inventory(service, base)
            request = mutation_request("UPDATE", {"alias": "Branch Updated", "configuration": {"timezone": "America/Argentina/Buenos_Aires"}, "capabilities": {"outbound": {"is_supported": False}}})
            request_bytes = pretty_bytes(request)
            result = execute_mutation(service, request, request_bytes, approval(request, request_bytes, receipt, receipt_bytes, "source-1"), receipt, receipt_bytes, profile(), base / "mutation", now=NOW)
            call = [item for item in service.calls if item[0] == "update"][0]
            self.assertIsInstance(call[2]["payload"], UpdateSupplySourceRequest)
            self.assertEqual(service.sources["source-1"]["alias"], "Branch Updated")
            self.assertTrue(result["provider_effect_reconciled"])

    def test_status_uses_official_enum_and_reconciles(self):
        with tempfile.TemporaryDirectory() as root:
            base = Path(root)
            service = FakeService()
            receipt, receipt_bytes = inventory(service, base)
            request = mutation_request("STATUS", {"status": "Inactive"})
            request_bytes = pretty_bytes(request)
            execute_mutation(service, request, request_bytes, approval(request, request_bytes, receipt, receipt_bytes, "source-1"), receipt, receipt_bytes, profile(), base / "mutation", now=NOW)
            call = [item for item in service.calls if item[0] == "status"][0]
            self.assertIsInstance(call[2]["payload"], UpdateSupplySourceStatusRequest)
            self.assertEqual(call[2]["payload"].status, "Inactive")

    def test_archive_requires_inactive_then_reconciles_archived(self):
        with tempfile.TemporaryDirectory() as root:
            base = Path(root)
            active = FakeService()
            receipt, receipt_bytes = inventory(active, base)
            request = mutation_request("ARCHIVE", {})
            request_bytes = pretty_bytes(request)
            with self.assertRaises(PermissionError):
                execute_mutation(active, request, request_bytes, approval(request, request_bytes, receipt, receipt_bytes, "source-1"), receipt, receipt_bytes, profile(), base / "blocked", now=NOW)
            self.assertFalse((base / "blocked").exists())
        with tempfile.TemporaryDirectory() as root:
            base = Path(root)
            inactive = FakeService([source(status="Inactive")])
            receipt, receipt_bytes = inventory(inactive, base)
            request = mutation_request("ARCHIVE", {})
            request_bytes = pretty_bytes(request)
            result = execute_mutation(inactive, request, request_bytes, approval(request, request_bytes, receipt, receipt_bytes, "source-1"), receipt, receipt_bytes, profile(), base / "mutation", now=NOW)
            self.assertEqual(inactive.sources["source-1"]["status"], "Archived")
            self.assertTrue(result["provider_effect_reconciled"])

    def test_stale_inventory_blocks_before_provider_mutation(self):
        with tempfile.TemporaryDirectory() as root:
            base = Path(root)
            service = FakeService()
            receipt, receipt_bytes = inventory(service, base)
            service.sources["source-1"]["alias"] = "Drifted"
            request = mutation_request("STATUS", {"status": "Inactive"})
            request_bytes = pretty_bytes(request)
            with self.assertRaises(PermissionError):
                execute_mutation(service, request, request_bytes, approval(request, request_bytes, receipt, receipt_bytes, "source-1"), receipt, receipt_bytes, profile(), base / "mutation", now=NOW)
            self.assertEqual([call for call in service.calls if call[0] == "status"], [])
            self.assertFalse((base / "mutation").exists())

    def test_tampered_approval_and_unknown_model_field_block_before_effect(self):
        with tempfile.TemporaryDirectory() as root:
            base = Path(root)
            service = FakeService([])
            receipt, receipt_bytes = inventory(service, base)
            request = create_request()
            request_bytes = pretty_bytes(request)
            bad = approval(request, request_bytes, receipt, receipt_bytes, None)
            bad["request_sha256"] = "0" * 64
            with self.assertRaises(PermissionError):
                execute_mutation(service, request, request_bytes, bad, receipt, receipt_bytes, profile(), base / "tampered", now=NOW)
            unknown = create_request()
            unknown["payload"]["invented_field"] = True
            unknown_bytes = pretty_bytes(unknown)
            with self.assertRaises(ValueError):
                execute_mutation(service, unknown, unknown_bytes, approval(unknown, unknown_bytes, receipt, receipt_bytes, None), receipt, receipt_bytes, profile(), base / "unknown", now=NOW)
            self.assertEqual([call for call in service.calls if call[0] == "create"], [])

    def test_provider_exception_records_unknown_effect_without_sensitive_error(self):
        with tempfile.TemporaryDirectory() as root:
            base = Path(root)
            service = FakeService([])
            service.fail_create = True
            receipt, receipt_bytes = inventory(service, base)
            request = create_request()
            request_bytes = pretty_bytes(request)
            with self.assertRaises(RuntimeError) as captured:
                execute_mutation(service, request, request_bytes, approval(request, request_bytes, receipt, receipt_bytes, None), receipt, receipt_bytes, profile(), base / "mutation", now=NOW)
            self.assertNotIn("sensitive", str(captured.exception))
            unknown = json.loads((base / "mutation/UNKNOWN_EFFECT.json").read_text())
            self.assertEqual(unknown["status"], "UNKNOWN_EFFECT")
            self.assertFalse(unknown["automatic_retry_authorized"])
            self.assertFalse((base / "mutation/MUTATION_RECEIPT.json").exists())

    def test_reconciliation_mismatch_is_unknown_effect(self):
        with tempfile.TemporaryDirectory() as root:
            base = Path(root)
            service = FakeService()
            service.ignore_status = True
            receipt, receipt_bytes = inventory(service, base)
            request = mutation_request("STATUS", {"status": "Inactive"})
            request_bytes = pretty_bytes(request)
            with self.assertRaises(RuntimeError):
                execute_mutation(service, request, request_bytes, approval(request, request_bytes, receipt, receipt_bytes, "source-1"), receipt, receipt_bytes, profile(), base / "mutation", now=NOW)
            self.assertTrue((base / "mutation/UNKNOWN_EFFECT.json").is_file())
            self.assertFalse((base / "mutation/MUTATION_RECEIPT.json").exists())

    def test_existing_output_never_calls_provider(self):
        with tempfile.TemporaryDirectory() as root:
            base = Path(root)
            service = FakeService([])
            receipt, receipt_bytes = inventory(service, base)
            request = create_request()
            request_bytes = pretty_bytes(request)
            target = base / "mutation"
            target.mkdir()
            (target / "keep").write_text("keep")
            calls_before = len(service.calls)
            with self.assertRaises(FileExistsError):
                execute_mutation(service, request, request_bytes, approval(request, request_bytes, receipt, receipt_bytes, None), receipt, receipt_bytes, profile(), target, now=NOW)
            self.assertEqual(len(service.calls), calls_before)
            self.assertEqual((target / "keep").read_text(), "keep")


if __name__ == "__main__":
    unittest.main()
````

### FILE: `amazon_spapi_supply_sources/README.md`
```yaml
block_id: "PY-AMAZON-SPAPI-SUPPLY-SOURCES:README-MD:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration governed by official Amazon Supply Sources API and models"
license: "LicenseRef-Workspace-Owner"
sha256: "602298c4925f8b695e56588df60ef9934ca7d15eb1e0e786ecef4e3ceba2466e"
variables: []
secrets_allowed: false
```
````markdown
# Amazon SP-API Supply Sources

Wrapper local fail-closed sobre `amzn-sp-api==1.11.1` y Amazon Supply Sources v2020-07-01 para inventariar y gobernar tiendas/depósitos usados por fulfillment multi-location.

La lane lista todas las páginas, obtiene cada source completo, rechaza ciclos/duplicados/drift y persiste sólo identificadores empresariales, status y hashes de dirección/configuración/capacidades. `CREATE`, `UPDATE`, `STATUS` y `ARCHIVE` usan modelos/métodos oficiales; cada mutación requiere approval ligado al request, inventario y snapshot exactos, escribe el intento antes del efecto y reconcilia por `get_supply_source`. Una excepción posterior al inicio queda `UNKNOWN_EFFECT` y nunca autoriza retry automático.

El perfil distribuido está bloqueado. Deben probarse registro, rol/scope, marketplace/programa, dynamic sandbox, cuota/costo, tratamiento de ubicación/contacto y owner. Archivar exige primero `Inactive`. Ningún estado Amazon activa automáticamente una ubicación interna, y este adapter no es una red multi-location universal ni reemplaza inventario, disponibilidad, promesas, transporte o conciliación del negocio.

```powershell
python -m pip install --require-hashes -r requirements.lock
python -m unittest -v test_supply_sources.py
```
````

## 6. Configuration surface

El perfil distribuido queda bloqueado. El usuario debe demostrar developer/app, rol/scope, marketplace/programa, dynamic sandbox, cuota/costo, tratamiento de ubicación/contacto y owner. Debe elegir operaciones explícitas, límites de paginación y secretos sólo por variables de entorno. `archive_requires_inactive=true` y `automatic_internal_location_activation=false` no pueden relajarse.

## 7. Dependency bill

`requirements.lock` fija `amzn-sp-api==1.11.1` y su grafo por URL PyPI y SHA-256. `sdk-artifact.lock.json` fija wheel, commit del SDK, commit de modelos, hashes y Apache-2.0. Los nueve archivos son `AUTHORED`; no se presentan como producto Amazon copiado.

## 8. Apply order

1. Materializar junto a `OFFICIAL-UPSTREAM-ACQUISITION-CORE`. 2. Instalar con `--require-hashes`. 3. Completar y aprobar perfil. 4. Generar inventario contra dynamic sandbox. 5. Crear request y approval ligados al receipt exacto. 6. Ejecutar una sola mutación. 7. Revisar reconciliación o `UNKNOWN_EFFECT`. 8. Reconciliar con el modelo interno sin activar automáticamente.

## 9. Verification

Debe pasar materialización 9/9, hashes exactos, `pip check` y 14 pruebas: firmas/modelos/enums oficiales, perfil bloqueado, paginación positiva, ciclo/duplicados, filtrado de datos, create/update/status/archive, archive sólo Inactive, drift, tamper/campo inventado, efecto ambiguo, reconciliación divergente y no-overwrite. Cuenta/dynamic sandbox/live, disponibilidad, inventario location-level, promesas, carga y producción siguen siendo gates del proyecto.

## 10. Reconstruction evidence

`reconstruction_evidence/AMAZON_SPAPI_SUPPLY_SOURCES_EXECUTION_INVENTORY_2026-08-31_V155.md` registra archivos, hashes y resultados autor/reconstruido. Demuestra el adapter local y SDK fijado; no demuestra acceso, soporte del marketplace, efecto Amazon live ni multi-location productivo.
