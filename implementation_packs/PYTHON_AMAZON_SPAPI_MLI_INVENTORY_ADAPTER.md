# Python Amazon SP-API Multi-Location Inventory Adapter

## 1. Metadata

```yaml
pack_id: "PYTHON-AMAZON-SPAPI-MLI-INVENTORY-ADAPTER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa snapshot y reemplazo reconciliado de inventario Amazon por Supply Source mediante Listings Items v2021-08-01 y el SDK oficial amzn-sp-api 1.11.1; no inventa soporte FBA por ubicación ni confirma inventario interno."
stacks: ["Python 3.10-3.14", "amzn-sp-api 1.11.1", "Amazon Listings Items v2021-08-01", "Amazon Multi-Location Inventory"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "PYTHON-AMAZON-SPAPI-SUPPLY-SOURCES-ADAPTER 0.1.x"]
incompatible_with: ["FBA location-level inventory", "marketplace o programa MLI no demostrado", "source Inactive/Archived", "PATCH sin approval y snapshots hash-bound", "retry automático de efecto ambiguo", "persistencia de seller/SKU/source IDs crudos", "commit automático de inventario interno", "secreto en CLI/Markdown", "SDK distinto de 1.11.1"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0 AND MIT AND MPL-2.0 AND BSD-3-Clause"
upstream_sources: ["https://github.com/amzn/selling-partner-api-sdk/tree/8e792ae345a8d334ccdbdd03181f05f040e6a4fc", "https://github.com/amzn/selling-partner-api-models/tree/8e429486005c4ebdce5099e48cc48515a65359bb", "https://developer-docs.amazon.com/sp-api/docs/mli-integration-guide", "https://developer-docs.amazon.com/sp-api/docs/listings-items-api-v2021-08-01-reference"]
verified_at: "2026-08-31"
```

## 2. Applicability

Use sólo para inventario Fulfilled by Merchant inscripto en Amazon Multi-Location Inventory. Amazon exige Supply Sources y Listings Items/Feeds; este pack implementa la lane individual Listings Items. Rechazar para FBA por ubicación, mercados fuera del programa, sources no activos o cualquier target sin cuenta, rol, sandbox, cuota, owner y política de datos demostrados.

## 3. Architecture contract

El snapshot consulta `get_listings_item(... includedData=fulfillmentAvailability)`, liga target/request/receipt y persiste únicamente hashes de seller, SKU y source IDs. El reemplazo exige approval sobre bytes exactos, relee para detectar drift, escribe attempt antes del efecto, construye `ListingsItemPatchRequest(product_type="PRODUCT")` con un único `PatchOperation(op="replace", path="/attributes/fulfillment_availability")`, rechaza issues y reconcilia con GET acotado. Excepción o falta de reconciliación produce `UNKNOWN_EFFECT`; jamás autoriza retry ni commit interno automático. El receipt Supply Sources debe contener sources únicos y `Active`.

## 4. Exact file manifest

```text
CREATE amazon_spapi_mli_inventory/requirements-direct.in
CREATE amazon_spapi_mli_inventory/requirements.lock
CREATE amazon_spapi_mli_inventory/sdk-artifact.lock.json
CREATE amazon_spapi_mli_inventory/provider-profile.template.json
CREATE amazon_spapi_mli_inventory/location-inventory-request.template.json
CREATE amazon_spapi_mli_inventory/mutation-approval.template.json
CREATE amazon_spapi_mli_inventory/run_mli_inventory.py
CREATE amazon_spapi_mli_inventory/test_mli_inventory.py
CREATE amazon_spapi_mli_inventory/README.md
```

## 5. Materialization blocks

### FILE: `amazon_spapi_mli_inventory/requirements-direct.in`
```yaml
block_id: "PY-AMAZON-SPAPI-MLI-INVENTORY:REQUIREMENTS-DIRECT-IN:v1"
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

### FILE: `amazon_spapi_mli_inventory/requirements.lock`
```yaml
block_id: "PY-AMAZON-SPAPI-MLI-INVENTORY:REQUIREMENTS-LOCK:v1"
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

### FILE: `amazon_spapi_mli_inventory/sdk-artifact.lock.json`
```yaml
block_id: "PY-AMAZON-SPAPI-MLI-INVENTORY:SDK-ARTIFACT-LOCK-JSON:v1"
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

### FILE: `amazon_spapi_mli_inventory/provider-profile.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-MLI-INVENTORY:PROVIDER-PROFILE-TEMPLATE-JSON:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration governed by official Amazon MLI, Listings Items and models"
license: "LicenseRef-Workspace-Owner"
sha256: "f1b4053d0854bf6eefcbef747419c583fc603064bdd0763df0516b08aad9c907"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-mli-inventory-profile/v1",
  "provider": "Amazon Listings Items v2021-08-01",
  "sdk": "amzn-sp-api==1.11.1",
  "decision": "BLOCKED_UNTIL_PROJECT_APPROVAL",
  "region": "SANDBOX",
  "developer_registration": "UNPROVEN",
  "application_registration": "UNPROVEN",
  "product_listing_role_and_scope": "UNPROVEN",
  "multi_location_inventory_program": "UNPROVEN",
  "supply_sources_inventory_receipt": "UNPROVEN",
  "dynamic_sandbox_contract_test": "UNPROVEN",
  "quota_and_cost_approved": false,
  "inventory_data_handling_approved": false,
  "inventory_authority_owner": "",
  "max_locations_per_request": 50,
  "max_quantity_per_location": 1000000,
  "reconciliation_attempts": 5,
  "reconciliation_delay_seconds": 2.0,
  "automatic_internal_inventory_commit": false,
  "client_id_env": "AMAZON_SPAPI_CLIENT_ID",
  "client_secret_env": "AMAZON_SPAPI_CLIENT_SECRET",
  "refresh_token_env": "AMAZON_SPAPI_REFRESH_TOKEN"
}
````

### FILE: `amazon_spapi_mli_inventory/location-inventory-request.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-MLI-INVENTORY:LOCATION-INVENTORY-REQUEST-TEMPLATE-JSON:v1"
operation: CREATE
provenance: AUTHORED
source: "local strict request governed by official fulfillment_availability schema"
license: "LicenseRef-Workspace-Owner"
sha256: "ec23f5ef7f13c858d05d10491751367f09775230608b5f7c1820a461f5cda766"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-mli-inventory-replace/v1",
  "seller_id": "REPLACE_ME",
  "sku": "REPLACE_ME",
  "marketplace_id": "REPLACE_ME",
  "locations": [
    {"supply_source_id": "REPLACE_ME", "quantity": 0}
  ]
}
````

### FILE: `amazon_spapi_mli_inventory/mutation-approval.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-MLI-INVENTORY:MUTATION-APPROVAL-TEMPLATE-JSON:v1"
operation: CREATE
provenance: AUTHORED
source: "local exact-byte approval envelope"
license: "LicenseRef-Workspace-Owner"
sha256: "62aaaf5b1b6408668f288edd4f0cdeaa0d2febfc71925aba58ec003aa0d46f87"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-mli-inventory-approval/v1",
  "approved_by": "",
  "approved_at": "1970-01-01T00:00:00Z",
  "request_sha256": "",
  "supply_sources_inventory_receipt_sha256": "",
  "expected_listing_snapshot_sha256": "",
  "allow_unknown_effect_retry": false
}
````

### FILE: `amazon_spapi_mli_inventory/run_mli_inventory.py`
```yaml
block_id: "PY-AMAZON-SPAPI-MLI-INVENTORY:RUN-MLI-INVENTORY-PY:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed adapter over exact official Amazon SDK methods/models"
license: "LicenseRef-Workspace-Owner"
sha256: "dd25364471ee41e017c673b6506b20e0b9643b4012884dd6ff01d54a7ea5f523"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
from datetime import datetime, timezone
import hashlib
import importlib.metadata
import json
import os
from pathlib import Path
import re
import shutil
import sys
import time
from typing import Any, Protocol
import uuid

from spapi import SPAPIClient, SPAPIConfig
from spapi.api.listings_items_v2021_08_01.listings_api import ListingsApi
from spapi.models.listings_items_v2021_08_01 import ListingsItemPatchRequest, PatchOperation

SDK_VERSION = "1.11.1"
REGIONS = {"SANDBOX", "NA", "EU", "FE"}
REQUEST_KEYS = {"schema", "seller_id", "sku", "marketplace_id", "locations"}
LOCATION_KEYS = {"supply_source_id", "quantity"}
APPROVAL_KEYS = {
    "schema",
    "approved_by",
    "approved_at",
    "request_sha256",
    "supply_sources_inventory_receipt_sha256",
    "expected_listing_snapshot_sha256",
    "allow_unknown_effect_retry",
}


class ListingsService(Protocol):
    def get_listings_item(self, seller_id: str, sku: str, marketplace_ids: list[str], **kwargs: Any) -> Any: ...
    def patch_listings_item(self, seller_id: str, sku: str, marketplace_ids: list[str], body: Any, **kwargs: Any) -> Any: ...


class ProviderRejectedError(Exception):
    pass


def sha256(value: bytes) -> str:
    return hashlib.sha256(value).hexdigest()


def json_value(value: Any) -> Any:
    if value is None or isinstance(value, (str, int, float, bool)):
        return value
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


def _unique_object(pairs: list[tuple[str, Any]], label: str) -> dict[str, Any]:
    result: dict[str, Any] = {}
    for key, value in pairs:
        if key in result:
            raise ValueError(f"duplicate key in {label}: {key}")
        result[key] = value
    return result


def strict_object(path: Path, label: str) -> tuple[dict[str, Any], bytes]:
    raw = path.read_bytes()
    value = json.loads(raw, object_pairs_hook=lambda pairs: _unique_object(pairs, label))
    if not isinstance(value, dict):
        raise ValueError(f"{label} must be an object")
    return value, raw


def exact_keys(value: Any, expected: set[str], label: str) -> dict[str, Any]:
    if not isinstance(value, dict) or set(value) != expected:
        raise ValueError(f"{label} keys must be exactly {sorted(expected)}")
    return value


def nonempty(value: Any, label: str, maximum: int = 255) -> str:
    if not isinstance(value, str) or not value or value != value.strip() or len(value) > maximum:
        raise ValueError(f"{label} must be a trimmed non-empty string up to {maximum} characters")
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
    required = {
        "schema", "provider", "sdk", "decision", "region", "developer_registration",
        "application_registration", "product_listing_role_and_scope", "multi_location_inventory_program",
        "supply_sources_inventory_receipt", "dynamic_sandbox_contract_test", "quota_and_cost_approved",
        "inventory_data_handling_approved", "inventory_authority_owner", "max_locations_per_request",
        "max_quantity_per_location", "reconciliation_attempts", "reconciliation_delay_seconds",
        "automatic_internal_inventory_commit", "client_id_env", "client_secret_env", "refresh_token_env",
    }
    exact_keys(profile, required, "provider profile")
    if (
        profile["schema"] != "elite-amazon-spapi-mli-inventory-profile/v1"
        or profile["provider"] != "Amazon Listings Items v2021-08-01"
        or profile["sdk"] != "amzn-sp-api==1.11.1"
        or profile["decision"] != "APPROVED"
    ):
        raise PermissionError("MLI inventory profile is not approved")
    if profile["region"] not in REGIONS:
        raise PermissionError("region is not approved")
    for key in (
        "developer_registration", "application_registration", "product_listing_role_and_scope",
        "multi_location_inventory_program", "supply_sources_inventory_receipt", "dynamic_sandbox_contract_test",
    ):
        if profile[key] != "PROVEN":
            raise PermissionError(f"{key} must be PROVEN")
    for key in ("quota_and_cost_approved", "inventory_data_handling_approved"):
        if profile[key] is not True:
            raise PermissionError(f"{key} must be true")
    if profile["automatic_internal_inventory_commit"] is not False:
        raise PermissionError("automatic internal inventory commit must remain false")
    nonempty(profile["inventory_authority_owner"], "inventory_authority_owner")
    integer_limits = {
        "max_locations_per_request": (1, 100),
        "max_quantity_per_location": (1, 2_000_000_000),
        "reconciliation_attempts": (1, 20),
    }
    for key, (minimum, maximum) in integer_limits.items():
        value = profile[key]
        if not isinstance(value, int) or isinstance(value, bool) or not minimum <= value <= maximum:
            raise ValueError(f"{key} must be within {minimum}..{maximum}")
    delay = profile["reconciliation_delay_seconds"]
    if not isinstance(delay, (int, float)) or isinstance(delay, bool) or not 0 <= delay <= 60:
        raise ValueError("reconciliation_delay_seconds must be within 0..60")
    for key in ("client_id_env", "client_secret_env", "refresh_token_env"):
        if not re.fullmatch(r"[A-Z][A-Z0-9_]{2,127}", nonempty(profile[key], key, 128)):
            raise ValueError(f"{key} must name an environment variable")
    return profile


def validate_supply_sources_receipt(value: dict[str, Any]) -> set[str]:
    if value.get("schema") != "elite-amazon-spapi-supply-sources-inventory/v1":
        raise ValueError("Supply Sources inventory receipt schema is invalid")
    sources = value.get("sources")
    if not isinstance(sources, list) or value.get("inventory_sha256") != sha256(canonical_bytes(sources)):
        raise ValueError("Supply Sources inventory receipt hash is invalid")
    identifiers: set[str] = set()
    for source in sources:
        if not isinstance(source, dict):
            raise ValueError("Supply Sources inventory entry is invalid")
        source_id = nonempty(source.get("supply_source_id"), "supply_source_id")
        if source_id in identifiers or source.get("status") != "Active":
            raise ValueError("Supply Sources receipt must contain unique Active sources")
        identifiers.add(source_id)
    if not identifiers:
        raise ValueError("Supply Sources receipt must contain at least one Active source")
    return identifiers


def validate_request(value: dict[str, Any], profile: dict[str, Any], allowed_sources: set[str]) -> list[dict[str, Any]]:
    exact_keys(value, REQUEST_KEYS, "inventory request")
    if value["schema"] != "elite-amazon-spapi-mli-inventory-replace/v1":
        raise ValueError("inventory request schema is invalid")
    nonempty(value["seller_id"], "seller_id", 256)
    nonempty(value["sku"], "sku", 256)
    nonempty(value["marketplace_id"], "marketplace_id", 64)
    locations = value["locations"]
    if not isinstance(locations, list) or not locations or len(locations) > profile["max_locations_per_request"]:
        raise ValueError("locations must be a non-empty bounded list")
    seen: set[str] = set()
    result: list[dict[str, Any]] = []
    for index, location in enumerate(locations):
        exact_keys(location, LOCATION_KEYS, f"locations[{index}]")
        source_id = nonempty(location["supply_source_id"], f"locations[{index}].supply_source_id")
        quantity = location["quantity"]
        if source_id in seen or source_id not in allowed_sources:
            raise PermissionError("location is duplicated or absent from approved Supply Sources receipt")
        if not isinstance(quantity, int) or isinstance(quantity, bool) or not 0 <= quantity <= profile["max_quantity_per_location"]:
            raise ValueError("location quantity is outside the approved range")
        seen.add(source_id)
        result.append({"supply_source_id": source_id, "quantity": quantity})
    return sorted(result, key=lambda item: item["supply_source_id"])


def provider_dict(value: Any, label: str) -> dict[str, Any]:
    data = json_value(value)
    if not isinstance(data, dict):
        raise ValueError(f"{label} is invalid")
    return data


def target_hash(request: dict[str, Any]) -> str:
    return sha256(canonical_bytes({key: request[key] for key in ("seller_id", "sku", "marketplace_id")}))


def filtered_availability(item: Any, request: dict[str, Any]) -> list[dict[str, Any]]:
    data = provider_dict(item, "listing item")
    if data.get("sku") != request["sku"]:
        raise ValueError("provider listing SKU differs from requested SKU")
    values = data.get("fulfillment_availability")
    if not isinstance(values, list):
        raise ValueError("provider fulfillmentAvailability is missing")
    result: list[dict[str, Any]] = []
    seen: set[str] = set()
    for value in values:
        if not isinstance(value, dict):
            raise ValueError("provider fulfillment availability entry is invalid")
        source_id = nonempty(value.get("fulfillment_channel_code"), "provider fulfillment_channel_code")
        quantity = value.get("quantity")
        if source_id in seen or not isinstance(quantity, int) or isinstance(quantity, bool) or quantity < 0:
            raise ValueError("provider fulfillment availability identity or quantity is invalid")
        seen.add(source_id)
        result.append({"supply_source_id_sha256": sha256(source_id.encode("utf-8")), "quantity": quantity})
    return sorted(result, key=lambda entry: entry["supply_source_id_sha256"])


def fetch_availability(service: ListingsService, request: dict[str, Any], timeout_seconds: float) -> list[dict[str, Any]]:
    item = service.get_listings_item(
        request["seller_id"],
        request["sku"],
        [request["marketplace_id"]],
        included_data=["fulfillmentAvailability"],
        _request_timeout=timeout_seconds,
    )
    return filtered_availability(item, request)


def expected_availability(locations: list[dict[str, Any]]) -> list[dict[str, Any]]:
    return sorted(
        [
            {"supply_source_id_sha256": sha256(item["supply_source_id"].encode("utf-8")), "quantity": item["quantity"]}
            for item in locations
        ],
        key=lambda entry: entry["supply_source_id_sha256"],
    )


def exclusive_write(path: Path, value: Any) -> None:
    with path.open("xb") as stream:
        stream.write(pretty_bytes(value))


def write_directory(target: Path, filename: str, value: Any) -> None:
    resolved = target.resolve()
    if resolved.exists():
        raise FileExistsError("output_directory must not exist")
    stage = resolved.parent / f".{resolved.name}.tmp-{uuid.uuid4().hex}"
    if stage.exists():
        raise FileExistsError("temporary output already exists")
    stage.mkdir(parents=True)
    try:
        exclusive_write(stage / filename, value)
        if resolved.exists():
            raise FileExistsError("output appeared during staging")
        stage.replace(resolved)
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def snapshot_listing(
    service: ListingsService,
    request: dict[str, Any],
    request_bytes: bytes,
    supply_receipt: dict[str, Any],
    supply_receipt_bytes: bytes,
    profile: dict[str, Any],
    output_directory: Path,
    timeout_seconds: float = 30.0,
    now: datetime | None = None,
) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION:
        raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if not 0 < timeout_seconds <= 120:
        raise ValueError("timeout must be within 0..120 seconds")
    if output_directory.resolve().exists():
        raise FileExistsError("output_directory must not exist")
    locations = validate_request(request, profile, validate_supply_sources_receipt(supply_receipt))
    availability = fetch_availability(service, request, timeout_seconds)
    receipt = {
        "schema": "elite-amazon-spapi-mli-listing-snapshot/v1",
        "created_at": (now or datetime.now(timezone.utc)).astimezone(timezone.utc).isoformat(),
        "provider": "Amazon Listings Items v2021-08-01",
        "sdk_version": SDK_VERSION,
        "target_identity_sha256": target_hash(request),
        "request_sha256": sha256(request_bytes),
        "supply_sources_inventory_receipt_sha256": sha256(supply_receipt_bytes),
        "requested_location_count": len(locations),
        "availability": availability,
        "availability_sha256": sha256(canonical_bytes(availability)),
        "raw_seller_sku_or_source_ids_persisted": False,
        "automatic_internal_inventory_commit": False,
    }
    write_directory(output_directory, "LISTING_INVENTORY.json", receipt)
    return receipt


def sanitized_issues(response: dict[str, Any]) -> list[dict[str, Any]]:
    issues = response.get("issues") or []
    if not isinstance(issues, list):
        raise ValueError("provider submission issues are invalid")
    result: list[dict[str, Any]] = []
    for issue in issues:
        if not isinstance(issue, dict):
            raise ValueError("provider submission issue is invalid")
        result.append(
            {
                "code": nonempty(issue.get("code"), "provider issue code", 256),
                "severity": nonempty(issue.get("severity"), "provider issue severity", 64),
                "attribute_names": issue.get("attribute_names") or [],
                "marketplace_ids": issue.get("marketplace_ids") or [],
            }
        )
    return result


def replace_inventory(
    service: ListingsService,
    request: dict[str, Any],
    request_bytes: bytes,
    approval: dict[str, Any],
    supply_receipt: dict[str, Any],
    supply_receipt_bytes: bytes,
    listing_snapshot: dict[str, Any],
    listing_snapshot_bytes: bytes,
    profile: dict[str, Any],
    output_directory: Path,
    timeout_seconds: float = 30.0,
    now: datetime | None = None,
) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION:
        raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if not 0 < timeout_seconds <= 120:
        raise ValueError("timeout must be within 0..120 seconds")
    target = output_directory.resolve()
    if target.exists():
        raise FileExistsError("output_directory must not exist")
    locations = validate_request(request, profile, validate_supply_sources_receipt(supply_receipt))
    exact_keys(approval, APPROVAL_KEYS, "mutation approval")
    if approval["schema"] != "elite-amazon-spapi-mli-inventory-approval/v1":
        raise ValueError("mutation approval schema is invalid")
    current_time = (now or datetime.now(timezone.utc)).astimezone(timezone.utc)
    nonempty(approval["approved_by"], "approved_by")
    if aware_datetime(approval["approved_at"], "approved_at") > current_time:
        raise ValueError("approved_at must not be in the future")
    if approval["allow_unknown_effect_retry"] is not False:
        raise PermissionError("unknown effects must never be retried automatically")
    bindings = {
        "request_sha256": sha256(request_bytes),
        "supply_sources_inventory_receipt_sha256": sha256(supply_receipt_bytes),
        "expected_listing_snapshot_sha256": sha256(listing_snapshot_bytes),
    }
    if any(approval[key] != value for key, value in bindings.items()):
        raise PermissionError("approval is not bound to exact request and receipts")
    if (
        listing_snapshot.get("schema") != "elite-amazon-spapi-mli-listing-snapshot/v1"
        or listing_snapshot.get("target_identity_sha256") != target_hash(request)
        or listing_snapshot.get("request_sha256") != sha256(request_bytes)
        or listing_snapshot.get("supply_sources_inventory_receipt_sha256") != sha256(supply_receipt_bytes)
        or listing_snapshot.get("availability_sha256") != sha256(canonical_bytes(listing_snapshot.get("availability")))
    ):
        raise ValueError("listing snapshot is invalid or belongs to another target")
    live_before = fetch_availability(service, request, timeout_seconds)
    if sha256(canonical_bytes(live_before)) != listing_snapshot["availability_sha256"]:
        raise PermissionError("live listing inventory drifted from approved snapshot")
    target.mkdir(parents=True)
    attempt = {
        "schema": "elite-amazon-spapi-mli-inventory-attempt/v1",
        "created_at": current_time.isoformat(),
        "target_identity_sha256": target_hash(request),
        **bindings,
        "desired_availability_sha256": sha256(canonical_bytes(expected_availability(locations))),
        "status": "PENDING",
        "automatic_retry_authorized": False,
    }
    exclusive_write(target / "MUTATION_ATTEMPT.json", attempt)
    official_value = [
        {"fulfillment_channel_code": item["supply_source_id"], "quantity": item["quantity"]}
        for item in locations
    ]
    body = ListingsItemPatchRequest(
        product_type="PRODUCT",
        patches=[PatchOperation(op="replace", path="/attributes/fulfillment_availability", value=official_value)],
    )
    effect_started = False
    try:
        effect_started = True
        raw_response = service.patch_listings_item(
            request["seller_id"],
            request["sku"],
            [request["marketplace_id"]],
            body,
            included_data=["issues"],
            _request_timeout=timeout_seconds,
        )
        response = provider_dict(raw_response, "patch submission response")
        issues = sanitized_issues(response)
        status = nonempty(response.get("status"), "provider submission status", 64)
        if status != "ACCEPTED" or any(issue["severity"].upper() == "ERROR" for issue in issues):
            rejection = {
                "schema": "elite-amazon-spapi-mli-provider-rejection/v1",
                "recorded_at": datetime.now(timezone.utc).isoformat(),
                "status": status,
                "issues": issues,
                "request_sha256": bindings["request_sha256"],
                "provider_effect_reconciled": False,
                "automatic_retry_authorized": False,
            }
            exclusive_write(target / "PROVIDER_REJECTED.json", rejection)
            raise ProviderRejectedError("provider rejected inventory submission")
        desired = expected_availability(locations)
        observed: list[dict[str, Any]] | None = None
        for attempt_index in range(profile["reconciliation_attempts"]):
            observed = fetch_availability(service, request, timeout_seconds)
            if observed == desired:
                break
            if attempt_index + 1 < profile["reconciliation_attempts"] and profile["reconciliation_delay_seconds"]:
                time.sleep(profile["reconciliation_delay_seconds"])
        if observed != desired:
            raise RuntimeError("accepted inventory submission was not reconciled")
        receipt = {
            "schema": "elite-amazon-spapi-mli-inventory-receipt/v1",
            "completed_at": datetime.now(timezone.utc).isoformat(),
            "target_identity_sha256": target_hash(request),
            "request_sha256": bindings["request_sha256"],
            "supply_sources_inventory_receipt_sha256": bindings["supply_sources_inventory_receipt_sha256"],
            "previous_listing_snapshot_sha256": bindings["expected_listing_snapshot_sha256"],
            "reconciled_availability_sha256": sha256(canonical_bytes(observed)),
            "provider_submission_status": status,
            "provider_issue_count": len(issues),
            "provider_effect_reconciled": True,
            "raw_seller_sku_or_source_ids_persisted": False,
            "automatic_internal_inventory_commit": False,
            "automatic_retry_authorized": False,
        }
        exclusive_write(target / "MUTATION_RECEIPT.json", receipt)
        return receipt
    except ProviderRejectedError as error:
        raise PermissionError(str(error)) from None
    except BaseException as error:
        if effect_started:
            unknown = {
                "schema": "elite-amazon-spapi-mli-inventory-unknown-effect/v1",
                "recorded_at": datetime.now(timezone.utc).isoformat(),
                "target_identity_sha256": target_hash(request),
                "request_sha256": bindings["request_sha256"],
                "error_type": type(error).__name__,
                "status": "UNKNOWN_EFFECT",
                "automatic_retry_authorized": False,
            }
            try:
                exclusive_write(target / "UNKNOWN_EFFECT.json", unknown)
            except FileExistsError:
                pass
            raise RuntimeError(f"inventory replacement effect is unknown: {type(error).__name__}") from None
        raise


def service_from_environment(profile: dict[str, Any]) -> ListingsApi:
    values: dict[str, str] = {}
    for key in ("client_id_env", "client_secret_env", "refresh_token_env"):
        value = os.environ.get(profile[key], "")
        if not value:
            raise PermissionError(f"required secret environment variable is missing: {profile[key]}")
        values[key] = value
    config = SPAPIConfig(
        client_id=values["client_id_env"],
        client_secret=values["client_secret_env"],
        refresh_token=values["refresh_token_env"],
        region=profile["region"],
    )
    return ListingsApi(SPAPIClient(config).api_client)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--mode", required=True, choices=("snapshot", "replace"))
    parser.add_argument("--request", required=True, type=Path)
    parser.add_argument("--supply-sources-inventory", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--approval", type=Path)
    parser.add_argument("--listing-snapshot", type=Path)
    parser.add_argument("--timeout-seconds", type=float, default=30.0)
    args = parser.parse_args()
    profile = approved_profile(args.profile)
    request, request_bytes = strict_object(args.request, "inventory request")
    supply_receipt, supply_receipt_bytes = strict_object(args.supply_sources_inventory, "Supply Sources inventory receipt")
    service = service_from_environment(profile)
    if args.mode == "snapshot":
        result = snapshot_listing(service, request, request_bytes, supply_receipt, supply_receipt_bytes, profile, args.output, args.timeout_seconds)
    else:
        if not args.approval or not args.listing_snapshot:
            raise ValueError("replace mode requires approval and listing snapshot")
        approval, _ = strict_object(args.approval, "mutation approval")
        listing_snapshot, listing_snapshot_bytes = strict_object(args.listing_snapshot, "listing snapshot")
        result = replace_inventory(
            service, request, request_bytes, approval, supply_receipt, supply_receipt_bytes,
            listing_snapshot, listing_snapshot_bytes, profile, args.output, args.timeout_seconds,
        )
    print(json.dumps({"status": "AMAZON_SPAPI_MLI_INVENTORY_PASS", "result": result}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `amazon_spapi_mli_inventory/test_mli_inventory.py`
```yaml
block_id: "PY-AMAZON-SPAPI-MLI-INVENTORY:TEST-MLI-INVENTORY-PY:v1"
operation: CREATE
provenance: AUTHORED
source: "local official-contract and negative regression suite"
license: "LicenseRef-Workspace-Owner"
sha256: "405cc770c41987b8ae836278e20c9d75a9bcbc58e46d607dfe141699c6cab5bf"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import copy
from datetime import datetime, timezone
import inspect
import json
from pathlib import Path
import tempfile
import unittest

from spapi.api.listings_items_v2021_08_01.listings_api import ListingsApi
from spapi.models.listings_items_v2021_08_01 import (
    FulfillmentAvailability,
    Issue,
    Item,
    ListingsItemPatchRequest,
    ListingsItemSubmissionResponse,
    PatchOperation,
)

import run_mli_inventory as subject


NOW = datetime(2026, 8, 31, 12, 0, tzinfo=timezone.utc)


def profile() -> dict:
    return {
        "schema": "elite-amazon-spapi-mli-inventory-profile/v1",
        "provider": "Amazon Listings Items v2021-08-01",
        "sdk": "amzn-sp-api==1.11.1",
        "decision": "APPROVED",
        "region": "SANDBOX",
        "developer_registration": "PROVEN",
        "application_registration": "PROVEN",
        "product_listing_role_and_scope": "PROVEN",
        "multi_location_inventory_program": "PROVEN",
        "supply_sources_inventory_receipt": "PROVEN",
        "dynamic_sandbox_contract_test": "PROVEN",
        "quota_and_cost_approved": True,
        "inventory_data_handling_approved": True,
        "inventory_authority_owner": "inventory-owner",
        "max_locations_per_request": 50,
        "max_quantity_per_location": 1000000,
        "reconciliation_attempts": 3,
        "reconciliation_delay_seconds": 0,
        "automatic_internal_inventory_commit": False,
        "client_id_env": "AMAZON_SPAPI_CLIENT_ID",
        "client_secret_env": "AMAZON_SPAPI_CLIENT_SECRET",
        "refresh_token_env": "AMAZON_SPAPI_REFRESH_TOKEN",
    }


def request() -> dict:
    return {
        "schema": "elite-amazon-spapi-mli-inventory-replace/v1",
        "seller_id": "SELLER-SECRET",
        "sku": "SKU-SECRET",
        "marketplace_id": "MARKET-SECRET",
        "locations": [
            {"supply_source_id": "SRC-A", "quantity": 7},
            {"supply_source_id": "SRC-B", "quantity": 11},
        ],
    }


def supply_receipt(status: str = "Active") -> dict:
    sources = [
        {"supply_source_id": "SRC-A", "status": status},
        {"supply_source_id": "SRC-B", "status": "Active"},
    ]
    return {
        "schema": "elite-amazon-spapi-supply-sources-inventory/v1",
        "sources": sources,
        "inventory_sha256": subject.sha256(subject.canonical_bytes(sources)),
    }


def encoded(value: dict) -> bytes:
    return subject.pretty_bytes(value)


class FakeListings:
    def __init__(self, availability: list[tuple[str, int]] | None = None):
        self.availability = availability or [("SRC-A", 1), ("SRC-B", 2)]
        self.get_calls: list[tuple] = []
        self.patch_calls: list[tuple] = []
        self.status = "ACCEPTED"
        self.issues: list[Issue] = []
        self.patch_error: Exception | None = None
        self.apply_patch = True

    def get_listings_item(self, seller_id, sku, marketplace_ids, **kwargs):
        self.get_calls.append((seller_id, sku, marketplace_ids, kwargs))
        return Item(
            sku=sku,
            fulfillment_availability=[
                FulfillmentAvailability(fulfillment_channel_code=source, quantity=quantity)
                for source, quantity in self.availability
            ],
        )

    def patch_listings_item(self, seller_id, sku, marketplace_ids, body, **kwargs):
        self.patch_calls.append((seller_id, sku, marketplace_ids, body, kwargs))
        if self.patch_error:
            raise self.patch_error
        if self.apply_patch:
            value = body.patches[0].value
            self.availability = [(item["fulfillment_channel_code"], item["quantity"]) for item in value]
        return ListingsItemSubmissionResponse(
            sku=sku,
            status=self.status,
            submission_id="submission-secret",
            issues=self.issues,
        )


def make_snapshot(root: Path, service: FakeListings, req: dict, supply: dict, prof: dict) -> tuple[dict, bytes, bytes, bytes]:
    req_bytes = encoded(req)
    supply_bytes = encoded(supply)
    target = root / "snapshot"
    value = subject.snapshot_listing(service, req, req_bytes, supply, supply_bytes, prof, target, now=NOW)
    raw = (target / "LISTING_INVENTORY.json").read_bytes()
    return value, raw, req_bytes, supply_bytes


def approval(req_bytes: bytes, supply_bytes: bytes, snapshot_bytes: bytes) -> dict:
    return {
        "schema": "elite-amazon-spapi-mli-inventory-approval/v1",
        "approved_by": "inventory-owner",
        "approved_at": "2026-08-31T11:00:00Z",
        "request_sha256": subject.sha256(req_bytes),
        "supply_sources_inventory_receipt_sha256": subject.sha256(supply_bytes),
        "expected_listing_snapshot_sha256": subject.sha256(snapshot_bytes),
        "allow_unknown_effect_retry": False,
    }


class MliInventoryTests(unittest.TestCase):
    def test_01_exact_official_api_signatures_and_models(self):
        self.assertEqual(
            list(inspect.signature(ListingsApi.get_listings_item).parameters),
            ["self", "seller_id", "sku", "marketplace_ids", "kwargs"],
        )
        self.assertEqual(
            list(inspect.signature(ListingsApi.patch_listings_item).parameters),
            ["self", "seller_id", "sku", "marketplace_ids", "body", "kwargs"],
        )
        self.assertEqual(ListingsItemPatchRequest.swagger_types, {"product_type": "str", "patches": "List[PatchOperation]"})
        self.assertEqual(PatchOperation.swagger_types, {"op": "str", "path": "str", "value": "List[Dict[str, object]]"})
        model = PatchOperation(op="replace", path="/attributes/fulfillment_availability", value=[])
        self.assertEqual(model.op, "replace")
        with self.assertRaises(ValueError):
            PatchOperation(op="invented", path="/attributes/fulfillment_availability", value=[])

    def test_02_distributed_profile_is_blocked(self):
        template = Path(__file__).with_name("provider-profile.template.json")
        with self.assertRaises(PermissionError):
            subject.approved_profile(template)

    def test_03_snapshot_filters_raw_identifiers_and_binds_receipts(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            service = FakeListings()
            req = request()
            supply = supply_receipt()
            value, raw, req_bytes, supply_bytes = make_snapshot(root, service, req, supply, profile())
            text = raw.decode("utf-8")
            for secret in ("SELLER-SECRET", "SKU-SECRET", "MARKET-SECRET", "SRC-A", "SRC-B"):
                self.assertNotIn(secret, text)
            self.assertEqual(value["request_sha256"], subject.sha256(req_bytes))
            self.assertEqual(value["supply_sources_inventory_receipt_sha256"], subject.sha256(supply_bytes))
            self.assertEqual(service.get_calls[0][3]["included_data"], ["fulfillmentAvailability"])
            self.assertFalse(value["automatic_internal_inventory_commit"])

    def test_04_request_rejects_unknown_duplicate_and_invalid_quantity(self):
        prof = profile()
        allowed = subject.validate_supply_sources_receipt(supply_receipt())
        for mutate in (
            lambda value: value.update(extra=True),
            lambda value: value["locations"].append(copy.deepcopy(value["locations"][0])),
            lambda value: value["locations"][0].update(supply_source_id="UNKNOWN"),
            lambda value: value["locations"][0].update(quantity=-1),
            lambda value: value["locations"][0].update(quantity=True),
        ):
            candidate = request()
            mutate(candidate)
            with self.assertRaises((ValueError, PermissionError)):
                subject.validate_request(candidate, prof, allowed)

    def test_05_supply_receipt_hash_and_active_status_are_mandatory(self):
        invalid_hash = supply_receipt()
        invalid_hash["inventory_sha256"] = "0" * 64
        with self.assertRaises(ValueError):
            subject.validate_supply_sources_receipt(invalid_hash)
        with self.assertRaises(ValueError):
            subject.validate_supply_sources_receipt(supply_receipt("Inactive"))

    def test_06_approval_tampering_blocks_before_patch(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            service = FakeListings()
            req = request()
            supply = supply_receipt()
            snapshot, snapshot_bytes, req_bytes, supply_bytes = make_snapshot(root, service, req, supply, profile())
            approve = approval(req_bytes, supply_bytes, snapshot_bytes)
            approve["request_sha256"] = "0" * 64
            with self.assertRaises(PermissionError):
                subject.replace_inventory(service, req, req_bytes, approve, supply, supply_bytes, snapshot, snapshot_bytes, profile(), root / "replace", now=NOW)
            self.assertEqual(service.patch_calls, [])

    def test_07_live_drift_blocks_before_patch(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            service = FakeListings()
            req = request()
            supply = supply_receipt()
            snapshot, snapshot_bytes, req_bytes, supply_bytes = make_snapshot(root, service, req, supply, profile())
            service.availability = [("SRC-A", 99), ("SRC-B", 2)]
            with self.assertRaises(PermissionError):
                subject.replace_inventory(service, req, req_bytes, approval(req_bytes, supply_bytes, snapshot_bytes), supply, supply_bytes, snapshot, snapshot_bytes, profile(), root / "replace", now=NOW)
            self.assertEqual(service.patch_calls, [])

    def test_08_replace_uses_exact_official_contract_and_reconciles(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            service = FakeListings()
            req = request()
            supply = supply_receipt()
            snapshot, snapshot_bytes, req_bytes, supply_bytes = make_snapshot(root, service, req, supply, profile())
            output = root / "replace"
            receipt = subject.replace_inventory(service, req, req_bytes, approval(req_bytes, supply_bytes, snapshot_bytes), supply, supply_bytes, snapshot, snapshot_bytes, profile(), output, now=NOW)
            body = service.patch_calls[0][3]
            self.assertIsInstance(body, ListingsItemPatchRequest)
            self.assertEqual(body.product_type, "PRODUCT")
            self.assertEqual(len(body.patches), 1)
            self.assertIsInstance(body.patches[0], PatchOperation)
            self.assertEqual(body.patches[0].op, "replace")
            self.assertEqual(body.patches[0].path, "/attributes/fulfillment_availability")
            self.assertEqual(service.patch_calls[0][4]["included_data"], ["issues"])
            self.assertTrue(receipt["provider_effect_reconciled"])
            self.assertFalse(receipt["automatic_internal_inventory_commit"])
            persisted = b"".join(path.read_bytes() for path in output.iterdir())
            for secret in (b"SELLER-SECRET", b"SKU-SECRET", b"MARKET-SECRET", b"SRC-A", b"SRC-B", b"submission-secret"):
                self.assertNotIn(secret, persisted)

    def test_09_provider_rejection_is_sanitized_and_not_unknown(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            service = FakeListings()
            service.status = "INVALID"
            service.issues = [Issue(code="BAD_ATTRIBUTE", message="SECRET-PROVIDER-MESSAGE", severity="ERROR", attribute_names=["fulfillment_availability"], categories=["INVALID_ATTRIBUTE"], marketplace_ids=["MARKET-SECRET"])]
            req = request()
            supply = supply_receipt()
            snapshot, snapshot_bytes, req_bytes, supply_bytes = make_snapshot(root, service, req, supply, profile())
            output = root / "replace"
            with self.assertRaises(PermissionError):
                subject.replace_inventory(service, req, req_bytes, approval(req_bytes, supply_bytes, snapshot_bytes), supply, supply_bytes, snapshot, snapshot_bytes, profile(), output, now=NOW)
            self.assertTrue((output / "PROVIDER_REJECTED.json").exists())
            self.assertFalse((output / "UNKNOWN_EFFECT.json").exists())
            persisted = (output / "PROVIDER_REJECTED.json").read_text("utf-8")
            self.assertNotIn("SECRET-PROVIDER-MESSAGE", persisted)

    def test_10_provider_exception_records_unknown_effect_without_message(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            service = FakeListings()
            service.patch_error = RuntimeError("SECRET-PROVIDER-ERROR")
            req = request()
            supply = supply_receipt()
            snapshot, snapshot_bytes, req_bytes, supply_bytes = make_snapshot(root, service, req, supply, profile())
            output = root / "replace"
            with self.assertRaisesRegex(RuntimeError, "effect is unknown: RuntimeError"):
                subject.replace_inventory(service, req, req_bytes, approval(req_bytes, supply_bytes, snapshot_bytes), supply, supply_bytes, snapshot, snapshot_bytes, profile(), output, now=NOW)
            persisted = (output / "UNKNOWN_EFFECT.json").read_text("utf-8")
            self.assertNotIn("SECRET-PROVIDER-ERROR", persisted)
            self.assertIn('"automatic_retry_authorized": false', persisted)

    def test_11_reconciliation_mismatch_is_unknown_effect(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            service = FakeListings()
            service.apply_patch = False
            req = request()
            supply = supply_receipt()
            snapshot, snapshot_bytes, req_bytes, supply_bytes = make_snapshot(root, service, req, supply, profile())
            output = root / "replace"
            with self.assertRaisesRegex(RuntimeError, "effect is unknown"):
                subject.replace_inventory(service, req, req_bytes, approval(req_bytes, supply_bytes, snapshot_bytes), supply, supply_bytes, snapshot, snapshot_bytes, profile(), output, now=NOW)
            self.assertTrue((output / "MUTATION_ATTEMPT.json").exists())
            self.assertTrue((output / "UNKNOWN_EFFECT.json").exists())
            self.assertFalse((output / "MUTATION_RECEIPT.json").exists())

    def test_12_existing_output_blocks_all_provider_calls(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            output = root / "exists"
            output.mkdir()
            service = FakeListings()
            req = request()
            supply = supply_receipt()
            with self.assertRaises(FileExistsError):
                subject.snapshot_listing(service, req, encoded(req), supply, encoded(supply), profile(), output, now=NOW)
            self.assertEqual(service.get_calls, [])

    def test_13_provider_permission_exception_is_unknown_not_rejection(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            service = FakeListings()
            service.patch_error = PermissionError("provider credential detail")
            req = request()
            supply = supply_receipt()
            snapshot, snapshot_bytes, req_bytes, supply_bytes = make_snapshot(root, service, req, supply, profile())
            output = root / "replace"
            with self.assertRaisesRegex(RuntimeError, "effect is unknown: PermissionError"):
                subject.replace_inventory(service, req, req_bytes, approval(req_bytes, supply_bytes, snapshot_bytes), supply, supply_bytes, snapshot, snapshot_bytes, profile(), output, now=NOW)
            self.assertTrue((output / "UNKNOWN_EFFECT.json").exists())
            self.assertFalse((output / "PROVIDER_REJECTED.json").exists())


if __name__ == "__main__":
    unittest.main(verbosity=2)
````

### FILE: `amazon_spapi_mli_inventory/README.md`
```yaml
block_id: "PY-AMAZON-SPAPI-MLI-INVENTORY:README-MD:v1"
operation: CREATE
provenance: AUTHORED
source: "local operator contract with explicit upstream attribution and limits"
license: "LicenseRef-Workspace-Owner"
sha256: "1219842f2557248eb0855b79bf74b1545836e4bd0c3e9dda5e371a5847bf9b20"
variables: []
secrets_allowed: false
```
````markdown
# Amazon SP-API Multi-Location Inventory adapter

Adapter local sobre `amzn-sp-api==1.11.1` y los modelos oficiales Listings Items v2021-08-01. Implementa lectura y reemplazo explícito de `/attributes/fulfillment_availability` para `productType=PRODUCT`, como exige la guía oficial Multi-Location Inventory (MLI).

## Procedencia

- SDK y modelos: Amazon, Apache-2.0, hashes fijados en `sdk-artifact.lock.json`.
- Los scripts, configuración y tests de este directorio son `AUTHORED`; no se presentan como código escrito por Amazon.

## Uso condicionado

Instale únicamente el lock exacto con hashes. Complete y apruebe el perfil. Genere primero un inventario con el adapter Supply Sources y conserve su receipt. Ejecute `snapshot` para obtener el estado actual; calcule los tres hashes de aprobación sobre bytes exactos; recién después ejecute `replace`.

```text
python -m pip install --require-hashes -r requirements.lock
python run_mli_inventory.py --profile provider-profile.json --mode snapshot --request location-inventory-request.json --supply-sources-inventory SUPPLY_SOURCES_INVENTORY.json --output listing-before
python run_mli_inventory.py --profile provider-profile.json --mode replace --request location-inventory-request.json --approval mutation-approval.json --supply-sources-inventory SUPPLY_SOURCES_INVENTORY.json --listing-snapshot listing-before/LISTING_INVENTORY.json --output replacement
```

El adapter relee el listing antes de mutar, escribe `MUTATION_ATTEMPT.json` antes del PATCH oficial, rechaza issues provider, y reconcilia mediante GET. Excepción o reconciliación no demostrada produce `UNKNOWN_EFFECT.json`; nunca se autoriza retry automático. Los outputs guardan hashes de seller/SKU/source IDs, no esos valores crudos.

No demuestra cuenta, sandbox, MLI real, consistencia temporal de Amazon, sincronización ERP/POS/WMS ni producción. Tampoco activa inventario interno automáticamente: `automatic_internal_inventory_commit=false`.
````

## 6. Configuration surface

| Campo | Tipo/default seguro | Validación y efecto |
|---|---|---|
| `decision` | enum / blocked | Sólo `APPROVED` habilita ejecución. |
| registros, rol, MLI, receipt y sandbox | enum / `UNPROVEN` | Los seis deben ser `PROVEN`. |
| cuota y tratamiento de inventario | bool / false | Ambos deben ser true. |
| owner | string / vacío | Obligatorio y no secreto. |
| límites | enteros acotados | Máximo 100 locations, 2.000.000.000 unidades y 20 lecturas. |
| demora de reconciliación | 0..60 / 2s | Espera sólo entre GET posteriores al PATCH. |
| `automatic_internal_inventory_commit` | false | Cualquier true falla cerrado. |
| nombres de variables secretas | identificador env | El secreto vive sólo en environment. |

## 7. Dependency bill

| Package/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| `amzn-sp-api` | 1.11.1 wheel SHA `a88e…70b0` | Listings API y modelos | Apache-2.0 | runtime/test | PyPI + `amzn/selling-partner-api-sdk` |
| dependencias Python | URLs y SHA en `requirements.lock` | HTTP/runtime | MIT/MPL/BSD | runtime | PyPI |
| CPython | 3.10–3.14 | ejecución | PSF-2.0 | runtime/test | python.org |

## 8. Apply order

1. Materializar en workspace vacío y comprobar hashes.
2. Instalar `requirements.lock` con `--require-hashes` en entorno aislado.
3. Completar perfil y probar cuenta/rol/MLI/dynamic sandbox/cuota/política.
4. Generar receipt Supply Sources con el pack compatible y conservar bytes exactos.
5. Ejecutar snapshot, revisar estado y emitir approval externo ligado a request/receipts.
6. Ejecutar replace; no reintentar si aparece `UNKNOWN_EFFECT`; reconciliar manualmente.
7. Mantener el sistema interno sin commit hasta que su propia transacción consuma un receipt reconciliado.

## 9. Verification

```text
python -m pip install --require-hashes -r requirements.lock
python -m pip check
python -m unittest -v test_mli_inventory.py
```

Éxito focal: instalación exacta, `pip check` PASS y 13/13 pruebas PASS en árbol autor y reconstruido. Live exige cuenta, rol Product Listing, inscripción MLI, sources reales, dynamic sandbox, cuota, privacidad, observabilidad y prueba de reconciliación; no se infiere desde mocks.

## 10. Reconstruction evidence

`reconstruction_evidence/AMAZON_SPAPI_MLI_INVENTORY_EXECUTION_INVENTORY_2026-08-31_V156.md` registra hashes y resultados. Nueve bloques `AUTHORED`: usan literalmente métodos/modelos del SDK Amazon, pero el glue no se atribuye a Amazon. El estado queda `CONDITIONED` hasta gates live y del proyecto.
