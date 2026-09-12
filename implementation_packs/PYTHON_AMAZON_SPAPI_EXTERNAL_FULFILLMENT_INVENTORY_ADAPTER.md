# Python Amazon SP-API External Fulfillment Inventory Adapter

## 1. Metadata

```yaml
pack_id: "PYTHON-AMAZON-SPAPI-EXTERNAL-FULFILLMENT-INVENTORY-ADAPTER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa FETCH y UPDATE batch de inventario absoluto por ubicación Amazon External Fulfillment v2024-09-11 mediante el SDK oficial amzn-sp-api 1.11.1; no convierte el acknowledgement provider en inventario interno confirmado."
stacks: ["Python 3.10-3.14", "amzn-sp-api 1.11.1", "Amazon External Fulfillment Inventory v2024-09-11"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "PYTHON-AMAZON-SPAPI-SUPPLY-SOURCES-ADAPTER 0.1.x", "PYTHON-AMAZON-SPAPI-MLI-INVENTORY-ADAPTER 0.1.x"]
incompatible_with: ["programa/canal no demostrado", "batch mayor a 10", "sequence no ligada por ubicación+SKU+marketplace+canal", "UPDATE sin FETCH y approval hash-bound", "retry automático de efecto ambiguo", "persistencia de IDs crudos", "commit automático ERP/POS/WMS", "secreto en CLI/Markdown", "SDK distinto de 1.11.1"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0 AND LicenseRef-Amazon-Software-License AND MIT AND MPL-2.0 AND BSD-3-Clause"
upstream_sources: ["https://github.com/amzn/selling-partner-api-sdk/tree/8e792ae345a8d334ccdbdd03181f05f040e6a4fc", "https://github.com/amzn/selling-partner-api-models/blob/8e429486005c4ebdce5099e48cc48515a65359bb/models/external-fulfillment/externalFulfillmentInventory_2024-09-11.json", "https://developer-docs.amazon.com/sp-api/reference/external-fulfillment-inventory-v2024-09-11", "https://developer-docs.amazon.com/sp-api/docs/external-fulfillment-apis-rate-limits"]
verified_at: "2026-08-31"
```

## 2. Applicability

Use sólo para programas y canales External Fulfillment Amazon aprobados que publiquen inventario absoluto location-level desde POS, ERP o WMS. La referencia oficial enumera Seller Flex/FBA Onsite, Multi Seller Flex, Easy ship, Self ship, MFN Self Delivery y Amazon Pharmacy; el modelo admite canales `FBA`, `MFN`, `DF`. Rechazar sin cuenta, programa, roles, mapping, sequence authority, sandbox, cuota, privacidad y owner demostrados.

## 3. Architecture contract

Cada batch contiene 1–10 identidades únicas. FETCH crea receipt sin IDs crudos. UPDATE exige sequence por identidad idéntica al receipt, approval sobre bytes exactos y un FETCH live sin drift. Antes del efecto persiste attempt; construye `BatchInventoryRequest`/`InventoryRequest`/`InventoryRequestParams`/`MarketplaceAttributes` oficiales y llama una vez a `batch_inventory`. Rechazo total es conocido; éxito parcial, excepción o sequence no reconciliada produce `UNKNOWN_EFFECT` sin retry. Un nuevo FETCH valida sequence, pero la cantidad solicitada no se promueve a hecho interno porque la respuesta oficial separa sellable/reserved.

## 4. Exact file manifest

```text
CREATE amazon_spapi_external_inventory/requirements-direct.in
CREATE amazon_spapi_external_inventory/requirements.lock
CREATE amazon_spapi_external_inventory/sdk-artifact.lock.json
CREATE amazon_spapi_external_inventory/provider-profile.template.json
CREATE amazon_spapi_external_inventory/inventory-batch-request.template.json
CREATE amazon_spapi_external_inventory/mutation-approval.template.json
CREATE amazon_spapi_external_inventory/run_external_inventory.py
CREATE amazon_spapi_external_inventory/test_external_inventory.py
CREATE amazon_spapi_external_inventory/README.md
```

## 5. Materialization blocks

### FILE: `amazon_spapi_external_inventory/requirements-direct.in`
```yaml
block_id: "PY-AMAZON-SPAPI-EXTERNAL-INVENTORY:REQUIREMENTS-DIRECT-IN:v1"
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

### FILE: `amazon_spapi_external_inventory/requirements.lock`
```yaml
block_id: "PY-AMAZON-SPAPI-EXTERNAL-INVENTORY:REQUIREMENTS-LOCK:v1"
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

### FILE: `amazon_spapi_external_inventory/sdk-artifact.lock.json`
```yaml
block_id: "PY-AMAZON-SPAPI-EXTERNAL-INVENTORY:SDK-ARTIFACT-LOCK-JSON:v1"
operation: CREATE
provenance: AUTHORED
source: "local receipt for exact official Amazon source/model/wheel identities and license conditions"
license: "LicenseRef-Workspace-Owner"
sha256: "949adc6a1d1531cd35cd6da7590fd3fa8f91e2f3f1627aada94c0df0c90ac175"
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
  "source": {"repository":"amzn/selling-partner-api-sdk","release":"python-1.11.1","commit":"8e792ae345a8d334ccdbdd03181f05f040e6a4fc","archive_bytes":26447333,"archive_sha256":"5bf712836bda76619c796fdb2b3a0a13acb281ba8bcf1b39c8f4b92f8f4799a1","license_expression":"Apache-2.0"},
  "model_source": {"repository":"amzn/selling-partner-api-models","commit":"8e429486005c4ebdce5099e48cc48515a65359bb","archive_sha256":"22bd3d956e12fe887f2ad893979d676b449c8bac4f6be6fd8cb0d56ae149b2e3","repository_license":"Apache-2.0","external_fulfillment_model_license":"LicenseRef-Amazon-Software-License"},
  "verified_at": "2026-08-31"
}
````

### FILE: `amazon_spapi_external_inventory/provider-profile.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-EXTERNAL-INVENTORY:PROVIDER-PROFILE-TEMPLATE-JSON:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed integration governed by official Amazon External Fulfillment Inventory"
license: "LicenseRef-Workspace-Owner"
sha256: "fe115064acf6d891b12132a859a7e9abf3b8fc3f6c56b37b97692faeed4fac4d"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-external-inventory-profile/v1",
  "provider": "Amazon External Fulfillment Inventory v2024-09-11",
  "sdk": "amzn-sp-api==1.11.1",
  "decision": "BLOCKED_UNTIL_PROJECT_APPROVAL",
  "region": "SANDBOX",
  "approved_channels": [],
  "developer_registration": "UNPROVEN",
  "application_registration": "UNPROVEN",
  "external_fulfillment_program": "UNPROVEN",
  "roles_and_scope": "UNPROVEN",
  "location_and_sku_mapping": "UNPROVEN",
  "sequence_authority": "UNPROVEN",
  "sandbox_contract_test": "UNPROVEN",
  "quota_and_cost_approved": false,
  "inventory_data_handling_approved": false,
  "inventory_authority_owner": "",
  "max_batch_requests": 10,
  "max_quantity": 2000000000,
  "automatic_internal_inventory_commit": false,
  "client_id_env": "AMAZON_SPAPI_CLIENT_ID",
  "client_secret_env": "AMAZON_SPAPI_CLIENT_SECRET",
  "refresh_token_env": "AMAZON_SPAPI_REFRESH_TOKEN"
}
````

### FILE: `amazon_spapi_external_inventory/inventory-batch-request.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-EXTERNAL-INVENTORY:INVENTORY-BATCH-REQUEST-TEMPLATE-JSON:v1"
operation: CREATE
provenance: AUTHORED
source: "local strict envelope over official InventoryRequest fields"
license: "LicenseRef-Workspace-Owner"
sha256: "645494fe69b78a03f94b7a4682dfab2873372e920edc0ac24f0ca169cf79f27b"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-external-inventory-batch/v1",
  "operation": "FETCH",
  "items": [
    {
      "location_id": "REPLACE_ME",
      "sku_id": "REPLACE_ME",
      "marketplace_id": "REPLACE_ME",
      "channel_name": "MFN",
      "quantity": null,
      "client_sequence_number": null
    }
  ]
}
````

### FILE: `amazon_spapi_external_inventory/mutation-approval.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-EXTERNAL-INVENTORY:MUTATION-APPROVAL-TEMPLATE-JSON:v1"
operation: CREATE
provenance: AUTHORED
source: "local exact-byte approval envelope"
license: "LicenseRef-Workspace-Owner"
sha256: "263d8d5837542e61a64c81421e7c25e512bfa4b88f6bc0f1b43b3bef9e9c7fc3"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-external-inventory-approval/v1",
  "approved_by": "",
  "approved_at": "1970-01-01T00:00:00Z",
  "request_sha256": "",
  "expected_fetch_receipt_sha256": "",
  "allow_unknown_effect_retry": false
}
````

### FILE: `amazon_spapi_external_inventory/run_external_inventory.py`
```yaml
block_id: "PY-AMAZON-SPAPI-EXTERNAL-INVENTORY:RUN-EXTERNAL-INVENTORY-PY:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed adapter over exact official Amazon SDK methods/models"
license: "LicenseRef-Workspace-Owner"
sha256: "4425aac93eee30e5a88eeb5cbb8e22e1c5fbde4127ab62df5d6383d0f77d022b"
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
from typing import Any, Protocol
from urllib.parse import quote
import uuid

from spapi import SPAPIClient, SPAPIConfig
from spapi.api.external_fulfillment_inventory_v2024_09_11.batch_inventory_api import BatchInventoryApi
from spapi.models.external_fulfillment_inventory_v2024_09_11 import (
    BatchInventoryRequest,
    HttpMethod,
    InventoryRequest,
    InventoryRequestParams,
    MarketplaceAttributes,
)

SDK_VERSION = "1.11.1"
CHANNELS = {"FBA", "MFN", "DF"}
REGIONS = {"SANDBOX", "NA", "EU", "FE"}
REQUEST_KEYS = {"schema", "operation", "items"}
ITEM_KEYS = {"location_id", "sku_id", "marketplace_id", "channel_name", "quantity", "client_sequence_number"}
APPROVAL_KEYS = {"schema", "approved_by", "approved_at", "request_sha256", "expected_fetch_receipt_sha256", "allow_unknown_effect_retry"}


class InventoryService(Protocol):
    def batch_inventory(self, body: Any, **kwargs: Any) -> Any: ...


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
        "schema", "provider", "sdk", "decision", "region", "approved_channels",
        "developer_registration", "application_registration", "external_fulfillment_program",
        "roles_and_scope", "location_and_sku_mapping", "sequence_authority", "sandbox_contract_test",
        "quota_and_cost_approved", "inventory_data_handling_approved", "inventory_authority_owner",
        "max_batch_requests", "max_quantity", "automatic_internal_inventory_commit",
        "client_id_env", "client_secret_env", "refresh_token_env",
    }
    exact_keys(profile, required, "provider profile")
    if (
        profile["schema"] != "elite-amazon-spapi-external-inventory-profile/v1"
        or profile["provider"] != "Amazon External Fulfillment Inventory v2024-09-11"
        or profile["sdk"] != "amzn-sp-api==1.11.1"
        or profile["decision"] != "APPROVED"
    ):
        raise PermissionError("External Fulfillment Inventory profile is not approved")
    if profile["region"] not in REGIONS:
        raise PermissionError("region is not approved")
    channels = profile["approved_channels"]
    if not isinstance(channels, list) or not channels or len(channels) != len(set(channels)) or any(channel not in CHANNELS for channel in channels):
        raise PermissionError("approved_channels must be a unique official channel subset")
    for key in (
        "developer_registration", "application_registration", "external_fulfillment_program",
        "roles_and_scope", "location_and_sku_mapping", "sequence_authority", "sandbox_contract_test",
    ):
        if profile[key] != "PROVEN":
            raise PermissionError(f"{key} must be PROVEN")
    for key in ("quota_and_cost_approved", "inventory_data_handling_approved"):
        if profile[key] is not True:
            raise PermissionError(f"{key} must be true")
    if profile["automatic_internal_inventory_commit"] is not False:
        raise PermissionError("automatic internal inventory commit must remain false")
    nonempty(profile["inventory_authority_owner"], "inventory_authority_owner")
    if not isinstance(profile["max_batch_requests"], int) or isinstance(profile["max_batch_requests"], bool) or not 1 <= profile["max_batch_requests"] <= 10:
        raise ValueError("max_batch_requests must be within 1..10")
    if not isinstance(profile["max_quantity"], int) or isinstance(profile["max_quantity"], bool) or not 1 <= profile["max_quantity"] <= 2_000_000_000:
        raise ValueError("max_quantity is invalid")
    for key in ("client_id_env", "client_secret_env", "refresh_token_env"):
        if not re.fullmatch(r"[A-Z][A-Z0-9_]{2,127}", nonempty(profile[key], key, 128)):
            raise ValueError(f"{key} must name an environment variable")
    return profile


def validate_request(value: dict[str, Any], profile: dict[str, Any]) -> list[dict[str, Any]]:
    exact_keys(value, REQUEST_KEYS, "inventory request")
    if value["schema"] != "elite-amazon-spapi-external-inventory-batch/v1" or value["operation"] not in {"FETCH", "UPDATE"}:
        raise ValueError("inventory request schema or operation is invalid")
    items = value["items"]
    if not isinstance(items, list) or not items or len(items) > profile["max_batch_requests"]:
        raise ValueError("items must contain 1..max_batch_requests entries")
    seen: set[tuple[str, str, str, str]] = set()
    result: list[dict[str, Any]] = []
    for index, item in enumerate(items):
        exact_keys(item, ITEM_KEYS, f"items[{index}]")
        normalized = {
            "location_id": nonempty(item["location_id"], f"items[{index}].location_id", 256),
            "sku_id": nonempty(item["sku_id"], f"items[{index}].sku_id", 256),
            "marketplace_id": nonempty(item["marketplace_id"], f"items[{index}].marketplace_id", 64),
            "channel_name": item["channel_name"],
            "quantity": item["quantity"],
            "client_sequence_number": item["client_sequence_number"],
        }
        if normalized["channel_name"] not in profile["approved_channels"]:
            raise PermissionError("item channel is not approved")
        identity = tuple(normalized[key] for key in ("location_id", "sku_id", "marketplace_id", "channel_name"))
        if identity in seen:
            raise ValueError("inventory request contains duplicate identity")
        seen.add(identity)
        if value["operation"] == "FETCH":
            if normalized["quantity"] is not None or normalized["client_sequence_number"] is not None:
                raise ValueError("FETCH quantity and client_sequence_number must be null")
        else:
            for key in ("quantity", "client_sequence_number"):
                field = normalized[key]
                maximum = profile["max_quantity"] if key == "quantity" else 9_223_372_036_854_775_807
                if not isinstance(field, int) or isinstance(field, bool) or not 0 <= field <= maximum:
                    raise ValueError(f"UPDATE {key} is invalid")
        result.append(normalized)
    return sorted(result, key=lambda item: tuple(item[key] for key in ("location_id", "sku_id", "marketplace_id", "channel_name")))


def identity_view(item: dict[str, Any]) -> dict[str, Any]:
    return {key: item[key] for key in ("location_id", "sku_id", "marketplace_id", "channel_name")}


def target_set_sha256(items: list[dict[str, Any]]) -> str:
    return sha256(canonical_bytes([identity_view(item) for item in items]))


def uri_for(operation: str, item: dict[str, Any]) -> str:
    action = "fetch" if operation == "FETCH" else "update"
    return f"/inventory/{action}?locationId={quote(item['location_id'], safe='')}&skuId={quote(item['sku_id'], safe='')}"


def official_batch(operation: str, items: list[dict[str, Any]]) -> BatchInventoryRequest:
    requests: list[InventoryRequest] = []
    for item in items:
        attributes = MarketplaceAttributes(channel_name=item["channel_name"], marketplace_id=item["marketplace_id"])
        params = InventoryRequestParams(marketplace_attributes=attributes)
        if operation == "UPDATE":
            params.quantity = item["quantity"]
            params.client_sequence_number = item["client_sequence_number"]
        requests.append(InventoryRequest(uri=uri_for(operation, item), method=HttpMethod.POST, body=params))
    return BatchInventoryRequest(requests=requests)


def response_rows(value: Any, expected: list[dict[str, Any]]) -> tuple[list[dict[str, Any]], list[dict[str, Any]]]:
    data = json_value(value)
    responses = data.get("responses") if isinstance(data, dict) else None
    if not isinstance(responses, list) or len(responses) != len(expected):
        raise ValueError("provider batch response cardinality is invalid")
    successes: list[dict[str, Any]] = []
    failures: list[dict[str, Any]] = []
    expected_by_identity = {tuple(identity_view(item).values()): item for item in expected}
    seen: set[tuple[str, str, str, str]] = set()
    for response in responses:
        if not isinstance(response, dict) or not isinstance(response.get("status"), dict) or not isinstance(response.get("body"), dict):
            raise ValueError("provider inventory response is malformed")
        status = response["status"]
        body = response["body"]
        attributes = body.get("marketplace_attributes")
        if not isinstance(attributes, dict):
            raise ValueError("provider marketplace attributes are missing")
        identity = (
            nonempty(body.get("location_id"), "provider location_id", 256),
            nonempty(body.get("sku_id"), "provider sku_id", 256),
            nonempty(attributes.get("marketplace_id"), "provider marketplace_id", 64),
            attributes.get("channel_name"),
        )
        if identity not in expected_by_identity or identity in seen:
            raise ValueError("provider response identity is unexpected or duplicated")
        seen.add(identity)
        errors = body.get("actionable_errors") or []
        if not isinstance(errors, list):
            raise ValueError("provider actionable errors are invalid")
        sanitized_errors = []
        for error in errors:
            if not isinstance(error, dict):
                raise ValueError("provider actionable error is invalid")
            sanitized_errors.append({"error_type": error.get("error_type"), "error_sub_type": error.get("error_sub_type")})
        status_code = status.get("status_code")
        row = {
            "identity": identity,
            "status_code": status_code,
            "client_sequence_number": body.get("client_sequence_number"),
            "sellable_quantity": body.get("sellable_quantity"),
            "reserved_quantity": body.get("reserved_quantity"),
            "actionable_errors": sanitized_errors,
        }
        if status_code == 200 and not errors:
            for key in ("client_sequence_number", "sellable_quantity", "reserved_quantity"):
                field = row[key]
                if not isinstance(field, int) or isinstance(field, bool) or field < 0:
                    raise ValueError(f"provider {key} is invalid")
            successes.append(row)
        else:
            failures.append(row)
    if len(seen) != len(expected):
        raise ValueError("provider response omitted an identity")
    return successes, failures


def filtered_rows(rows: list[dict[str, Any]]) -> list[dict[str, Any]]:
    result = []
    for row in rows:
        location_id, sku_id, marketplace_id, channel_name = row["identity"]
        result.append({
            "location_id_sha256": sha256(location_id.encode("utf-8")),
            "sku_id_sha256": sha256(sku_id.encode("utf-8")),
            "marketplace_id_sha256": sha256(marketplace_id.encode("utf-8")),
            "channel_name": channel_name,
            "client_sequence_number": row["client_sequence_number"],
            "sellable_quantity": row["sellable_quantity"],
            "reserved_quantity": row["reserved_quantity"],
        })
    return sorted(result, key=lambda item: (item["location_id_sha256"], item["sku_id_sha256"], item["marketplace_id_sha256"], item["channel_name"]))


def filtered_identity(item: dict[str, Any]) -> tuple[str, str, str, str]:
    if "location_id_sha256" in item:
        return (item["location_id_sha256"], item["sku_id_sha256"], item["marketplace_id_sha256"], item["channel_name"])
    return (
        sha256(item["location_id"].encode("utf-8")),
        sha256(item["sku_id"].encode("utf-8")),
        sha256(item["marketplace_id"].encode("utf-8")),
        item["channel_name"],
    )


def fetch_provider(service: InventoryService, items: list[dict[str, Any]], timeout_seconds: float) -> list[dict[str, Any]]:
    response = service.batch_inventory(official_batch("FETCH", items), _request_timeout=timeout_seconds)
    successes, failures = response_rows(response, items)
    if failures:
        raise PermissionError("provider rejected one or more FETCH entries")
    return filtered_rows(successes)


def exclusive_write(path: Path, value: Any) -> None:
    with path.open("xb") as stream:
        stream.write(pretty_bytes(value))


def write_directory(target: Path, filename: str, value: Any) -> None:
    resolved = target.resolve()
    if resolved.exists():
        raise FileExistsError("output_directory must not exist")
    stage = resolved.parent / f".{resolved.name}.tmp-{uuid.uuid4().hex}"
    stage.mkdir(parents=True)
    try:
        exclusive_write(stage / filename, value)
        if resolved.exists():
            raise FileExistsError("output appeared during staging")
        stage.replace(resolved)
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def fetch_inventory(service: InventoryService, request: dict[str, Any], request_bytes: bytes, profile: dict[str, Any], output_directory: Path, timeout_seconds: float = 30.0, now: datetime | None = None) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION:
        raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if not 0 < timeout_seconds <= 120:
        raise ValueError("timeout must be within 0..120 seconds")
    if output_directory.resolve().exists():
        raise FileExistsError("output_directory must not exist")
    items = validate_request(request, profile)
    if request["operation"] != "FETCH":
        raise ValueError("fetch mode requires FETCH request")
    rows = fetch_provider(service, items, timeout_seconds)
    receipt = {
        "schema": "elite-amazon-spapi-external-inventory-fetch-receipt/v1",
        "created_at": (now or datetime.now(timezone.utc)).astimezone(timezone.utc).isoformat(),
        "provider": "Amazon External Fulfillment Inventory v2024-09-11",
        "sdk_version": SDK_VERSION,
        "request_sha256": sha256(request_bytes),
        "target_set_sha256": target_set_sha256(items),
        "items": rows,
        "items_sha256": sha256(canonical_bytes(rows)),
        "raw_location_sku_or_marketplace_ids_persisted": False,
        "automatic_internal_inventory_commit": False,
    }
    write_directory(output_directory, "FETCH_RECEIPT.json", receipt)
    return receipt


def update_inventory(service: InventoryService, request: dict[str, Any], request_bytes: bytes, approval: dict[str, Any], fetch_receipt: dict[str, Any], fetch_receipt_bytes: bytes, profile: dict[str, Any], output_directory: Path, timeout_seconds: float = 30.0, now: datetime | None = None) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION:
        raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if not 0 < timeout_seconds <= 120:
        raise ValueError("timeout must be within 0..120 seconds")
    target = output_directory.resolve()
    if target.exists():
        raise FileExistsError("output_directory must not exist")
    items = validate_request(request, profile)
    if request["operation"] != "UPDATE":
        raise ValueError("update mode requires UPDATE request")
    exact_keys(approval, APPROVAL_KEYS, "mutation approval")
    current_time = (now or datetime.now(timezone.utc)).astimezone(timezone.utc)
    if approval["schema"] != "elite-amazon-spapi-external-inventory-approval/v1" or not nonempty(approval["approved_by"], "approved_by"):
        raise ValueError("mutation approval is invalid")
    if aware_datetime(approval["approved_at"], "approved_at") > current_time:
        raise ValueError("approved_at must not be in the future")
    if approval["allow_unknown_effect_retry"] is not False:
        raise PermissionError("unknown effects must never be retried automatically")
    if approval["request_sha256"] != sha256(request_bytes) or approval["expected_fetch_receipt_sha256"] != sha256(fetch_receipt_bytes):
        raise PermissionError("approval is not bound to exact request and fetch receipt")
    if (
        fetch_receipt.get("schema") != "elite-amazon-spapi-external-inventory-fetch-receipt/v1"
        or fetch_receipt.get("target_set_sha256") != target_set_sha256(items)
        or fetch_receipt.get("items_sha256") != sha256(canonical_bytes(fetch_receipt.get("items")))
    ):
        raise ValueError("fetch receipt is invalid or belongs to another target set")
    baseline = fetch_receipt["items"]
    update_sequences = {filtered_identity(item): item["client_sequence_number"] for item in items}
    baseline_sequences = {filtered_identity(item): item["client_sequence_number"] for item in baseline}
    if update_sequences != baseline_sequences:
        raise PermissionError("UPDATE sequence numbers do not match approved FETCH receipt")
    live_before = fetch_provider(service, items, timeout_seconds)
    if live_before != baseline:
        raise PermissionError("live inventory drifted from approved FETCH receipt")
    target.mkdir(parents=True)
    attempt = {
        "schema": "elite-amazon-spapi-external-inventory-attempt/v1",
        "created_at": current_time.isoformat(),
        "request_sha256": sha256(request_bytes),
        "fetch_receipt_sha256": sha256(fetch_receipt_bytes),
        "target_set_sha256": target_set_sha256(items),
        "desired_quantities_sha256": sha256(canonical_bytes([item["quantity"] for item in items])),
        "status": "PENDING",
        "automatic_retry_authorized": False,
    }
    exclusive_write(target / "MUTATION_ATTEMPT.json", attempt)
    effect_started = False
    try:
        effect_started = True
        raw_response = service.batch_inventory(official_batch("UPDATE", items), _request_timeout=timeout_seconds)
        successes, failures = response_rows(raw_response, items)
        if failures:
            rejection = {
                "schema": "elite-amazon-spapi-external-inventory-provider-response/v1",
                "recorded_at": datetime.now(timezone.utc).isoformat(),
                "success_count": len(successes),
                "failure_count": len(failures),
                "failures": [{"status_code": row["status_code"], "actionable_errors": row["actionable_errors"]} for row in failures],
                "automatic_retry_authorized": False,
            }
            exclusive_write(target / ("PROVIDER_REJECTED.json" if not successes else "PARTIAL_EFFECT.json"), rejection)
            if successes:
                raise RuntimeError("provider returned a partial batch effect")
            raise ProviderRejectedError("provider rejected all inventory updates")
        response_sequences = {filtered_identity({"location_id": row["identity"][0], "sku_id": row["identity"][1], "marketplace_id": row["identity"][2], "channel_name": row["identity"][3]}): row["client_sequence_number"] for row in successes}
        if response_sequences != update_sequences:
            raise RuntimeError("provider update response sequence mismatch")
        reconciled = fetch_provider(service, items, timeout_seconds)
        reconciled_sequences = {filtered_identity(item): item["client_sequence_number"] for item in reconciled}
        if reconciled_sequences != update_sequences:
            raise RuntimeError("inventory update sequence was not reconciled")
        receipt = {
            "schema": "elite-amazon-spapi-external-inventory-update-receipt/v1",
            "completed_at": datetime.now(timezone.utc).isoformat(),
            "request_sha256": sha256(request_bytes),
            "previous_fetch_receipt_sha256": sha256(fetch_receipt_bytes),
            "target_set_sha256": target_set_sha256(items),
            "provider_response_sha256": sha256(canonical_bytes(filtered_rows(successes))),
            "reconciled_items_sha256": sha256(canonical_bytes(reconciled)),
            "provider_effect_reconciled_by_sequence": True,
            "requested_absolute_quantity_not_promoted_to_internal_fact": True,
            "raw_location_sku_or_marketplace_ids_persisted": False,
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
                "schema": "elite-amazon-spapi-external-inventory-unknown-effect/v1",
                "recorded_at": datetime.now(timezone.utc).isoformat(),
                "request_sha256": sha256(request_bytes),
                "target_set_sha256": target_set_sha256(items),
                "error_type": type(error).__name__,
                "status": "UNKNOWN_EFFECT",
                "automatic_retry_authorized": False,
            }
            try:
                exclusive_write(target / "UNKNOWN_EFFECT.json", unknown)
            except FileExistsError:
                pass
            raise RuntimeError(f"external inventory effect is unknown: {type(error).__name__}") from None
        raise


def service_from_environment(profile: dict[str, Any]) -> BatchInventoryApi:
    values = {}
    for key in ("client_id_env", "client_secret_env", "refresh_token_env"):
        value = os.environ.get(profile[key], "")
        if not value:
            raise PermissionError(f"required secret environment variable is missing: {profile[key]}")
        values[key] = value
    config = SPAPIConfig(client_id=values["client_id_env"], client_secret=values["client_secret_env"], refresh_token=values["refresh_token_env"], region=profile["region"])
    return BatchInventoryApi(SPAPIClient(config).api_client)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--mode", required=True, choices=("fetch", "update"))
    parser.add_argument("--request", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--approval", type=Path)
    parser.add_argument("--fetch-receipt", type=Path)
    parser.add_argument("--timeout-seconds", type=float, default=30.0)
    args = parser.parse_args()
    profile = approved_profile(args.profile)
    request, request_bytes = strict_object(args.request, "inventory request")
    service = service_from_environment(profile)
    if args.mode == "fetch":
        result = fetch_inventory(service, request, request_bytes, profile, args.output, args.timeout_seconds)
    else:
        if not args.approval or not args.fetch_receipt:
            raise ValueError("update mode requires approval and fetch receipt")
        approval, _ = strict_object(args.approval, "mutation approval")
        receipt, receipt_bytes = strict_object(args.fetch_receipt, "fetch receipt")
        result = update_inventory(service, request, request_bytes, approval, receipt, receipt_bytes, profile, args.output, args.timeout_seconds)
    print(json.dumps({"status": "AMAZON_SPAPI_EXTERNAL_INVENTORY_PASS", "result": result}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `amazon_spapi_external_inventory/test_external_inventory.py`
```yaml
block_id: "PY-AMAZON-SPAPI-EXTERNAL-INVENTORY:TEST-EXTERNAL-INVENTORY-PY:v1"
operation: CREATE
provenance: AUTHORED
source: "local official-contract and negative regression suite"
license: "LicenseRef-Workspace-Owner"
sha256: "ca10c8c2e15548bc3d204220b5da0af3342508a8e67cba5833b3877c83dda4e9"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import copy
from datetime import datetime, timezone
import inspect
from pathlib import Path
import tempfile
import unittest
from urllib.parse import parse_qs, urlparse

from spapi.api.external_fulfillment_inventory_v2024_09_11.batch_inventory_api import BatchInventoryApi
from spapi.models.external_fulfillment_inventory_v2024_09_11 import (
    ActionableError,
    BatchInventoryRequest,
    BatchInventoryResponse,
    HttpMethod,
    HttpStatusLine,
    InventoryRequest,
    InventoryRequestParams,
    InventoryResponse,
    InventoryResponseBody,
    MarketplaceAttributes,
)

import run_external_inventory as subject


NOW = datetime(2026, 8, 31, 12, 0, tzinfo=timezone.utc)


def profile() -> dict:
    return {
        "schema": "elite-amazon-spapi-external-inventory-profile/v1",
        "provider": "Amazon External Fulfillment Inventory v2024-09-11",
        "sdk": "amzn-sp-api==1.11.1",
        "decision": "APPROVED",
        "region": "SANDBOX",
        "approved_channels": ["MFN", "FBA"],
        "developer_registration": "PROVEN",
        "application_registration": "PROVEN",
        "external_fulfillment_program": "PROVEN",
        "roles_and_scope": "PROVEN",
        "location_and_sku_mapping": "PROVEN",
        "sequence_authority": "PROVEN",
        "sandbox_contract_test": "PROVEN",
        "quota_and_cost_approved": True,
        "inventory_data_handling_approved": True,
        "inventory_authority_owner": "inventory-owner",
        "max_batch_requests": 10,
        "max_quantity": 2000000000,
        "automatic_internal_inventory_commit": False,
        "client_id_env": "AMAZON_SPAPI_CLIENT_ID",
        "client_secret_env": "AMAZON_SPAPI_CLIENT_SECRET",
        "refresh_token_env": "AMAZON_SPAPI_REFRESH_TOKEN",
    }


def fetch_request() -> dict:
    return {
        "schema": "elite-amazon-spapi-external-inventory-batch/v1",
        "operation": "FETCH",
        "items": [
            {"location_id": "LOC A/1", "sku_id": "SKU A?", "marketplace_id": "MARKET-A", "channel_name": "MFN", "quantity": None, "client_sequence_number": None},
            {"location_id": "LOC-B", "sku_id": "SKU-B", "marketplace_id": "MARKET-B", "channel_name": "FBA", "quantity": None, "client_sequence_number": None},
        ],
    }


def update_request(sequence_a: int = 10, sequence_b: int = 20) -> dict:
    value = fetch_request()
    value["operation"] = "UPDATE"
    value["items"][0].update(quantity=17, client_sequence_number=sequence_a)
    value["items"][1].update(quantity=23, client_sequence_number=sequence_b)
    return value


def encoded(value: dict) -> bytes:
    return subject.pretty_bytes(value)


class FakeInventory:
    def __init__(self):
        self.state = {
            ("LOC A/1", "SKU A?", "MARKET-A", "MFN"): {"sequence": 10, "sellable": 5, "reserved": 2},
            ("LOC-B", "SKU-B", "MARKET-B", "FBA"): {"sequence": 20, "sellable": 8, "reserved": 3},
        }
        self.calls: list[BatchInventoryRequest] = []
        self.raise_error: Exception | None = None
        self.raise_call_number: int | None = None
        self.reject_indexes: set[int] = set()
        self.apply_updates = True

    def batch_inventory(self, body, **kwargs):
        self.calls.append(body)
        if self.raise_error and (self.raise_call_number is None or len(self.calls) == self.raise_call_number):
            raise self.raise_error
        responses = []
        for index, request in enumerate(body.requests):
            parsed = urlparse(request.uri)
            query = parse_qs(parsed.query)
            attributes = request.body.marketplace_attributes
            identity = (query["locationId"][0], query["skuId"][0], attributes.marketplace_id, attributes.channel_name)
            state = self.state[identity]
            is_update = parsed.path.endswith("/update")
            if is_update and index in self.reject_indexes:
                responses.append(
                    InventoryResponse(
                        status=HttpStatusLine(status_code=400, reason_phrase="SECRET REASON"),
                        body=InventoryResponseBody(
                            location_id=identity[0], sku_id=identity[1], marketplace_attributes=attributes,
                            actionable_errors=[ActionableError(error_type="INVALID_INPUT", error_sub_type="SEQUENCE")],
                        ),
                    )
                )
                continue
            if is_update and self.apply_updates:
                state["sequence"] = request.body.client_sequence_number
                state["sellable"] = request.body.quantity
            responses.append(
                InventoryResponse(
                    status=HttpStatusLine(status_code=200, reason_phrase="Success"),
                    body=InventoryResponseBody(
                        client_sequence_number=state["sequence"], location_id=identity[0], sku_id=identity[1],
                        sellable_quantity=state["sellable"], reserved_quantity=state["reserved"],
                        marketplace_attributes=attributes, actionable_errors=[],
                    ),
                )
            )
        return BatchInventoryResponse(responses=responses)


def make_baseline(root: Path, service: FakeInventory) -> tuple[dict, bytes]:
    req = fetch_request()
    target = root / "fetch"
    receipt = subject.fetch_inventory(service, req, encoded(req), profile(), target, now=NOW)
    return receipt, (target / "FETCH_RECEIPT.json").read_bytes()


def approval(request_bytes: bytes, receipt_bytes: bytes) -> dict:
    return {
        "schema": "elite-amazon-spapi-external-inventory-approval/v1",
        "approved_by": "inventory-owner",
        "approved_at": "2026-08-31T11:00:00Z",
        "request_sha256": subject.sha256(request_bytes),
        "expected_fetch_receipt_sha256": subject.sha256(receipt_bytes),
        "allow_unknown_effect_retry": False,
    }


class ExternalInventoryTests(unittest.TestCase):
    def test_01_exact_official_signatures_models_and_enums(self):
        self.assertEqual(list(inspect.signature(BatchInventoryApi.batch_inventory).parameters), ["self", "body", "kwargs"])
        self.assertEqual(BatchInventoryRequest.swagger_types, {"requests": "List[InventoryRequest]"})
        self.assertEqual(InventoryRequest.swagger_types, {"uri": "str", "method": "HttpMethod", "body": "InventoryRequestParams"})
        self.assertEqual(InventoryRequestParams.swagger_types, {"client_sequence_number": "int", "quantity": "int", "marketplace_attributes": "MarketplaceAttributes"})
        self.assertEqual(HttpMethod.POST, "POST")
        with self.assertRaises(ValueError):
            MarketplaceAttributes(channel_name="INVENTED", marketplace_id="M")

    def test_02_distributed_profile_is_blocked(self):
        with self.assertRaises(PermissionError):
            subject.approved_profile(Path(__file__).with_name("provider-profile.template.json"))

    def test_03_fetch_uses_official_models_encodes_uri_and_filters_ids(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            service = FakeInventory()
            req = fetch_request()
            receipt = subject.fetch_inventory(service, req, encoded(req), profile(), root / "out", now=NOW)
            model = service.calls[0]
            self.assertIsInstance(model, BatchInventoryRequest)
            self.assertTrue(all(isinstance(item, InventoryRequest) for item in model.requests))
            self.assertEqual(model.requests[0].method, HttpMethod.POST)
            self.assertIn("locationId=LOC%20A%2F1", model.requests[0].uri)
            self.assertIn("skuId=SKU%20A%3F", model.requests[0].uri)
            persisted = (root / "out" / "FETCH_RECEIPT.json").read_bytes()
            for secret in (b"LOC A/1", b"LOC-B", b"SKU A?", b"SKU-B", b"MARKET-A", b"MARKET-B"):
                self.assertNotIn(secret, persisted)
            self.assertFalse(receipt["automatic_internal_inventory_commit"])

    def test_04_request_rejects_unknown_duplicate_bounds_and_wrong_mode_fields(self):
        prof = profile()
        candidates = []
        extra = fetch_request(); extra["extra"] = True; candidates.append(extra)
        duplicate = fetch_request(); duplicate["items"].append(copy.deepcopy(duplicate["items"][0])); candidates.append(duplicate)
        channel = fetch_request(); channel["items"][0]["channel_name"] = "DF"; candidates.append(channel)
        fetch_quantity = fetch_request(); fetch_quantity["items"][0]["quantity"] = 1; candidates.append(fetch_quantity)
        boolean = update_request(); boolean["items"][0]["quantity"] = True; candidates.append(boolean)
        for candidate in candidates:
            with self.assertRaises((ValueError, PermissionError)):
                subject.validate_request(candidate, prof)

    def test_05_positive_update_binds_sequence_uses_official_models_and_reconciles(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            service = FakeInventory()
            baseline, baseline_bytes = make_baseline(root, service)
            req = update_request()
            req_bytes = encoded(req)
            output = root / "update"
            receipt = subject.update_inventory(service, req, req_bytes, approval(req_bytes, baseline_bytes), baseline, baseline_bytes, profile(), output, now=NOW)
            update_model = service.calls[2]
            self.assertIsInstance(update_model, BatchInventoryRequest)
            self.assertTrue(all(isinstance(item.body, InventoryRequestParams) for item in update_model.requests))
            self.assertEqual(update_model.requests[0].body.quantity, 17)
            self.assertEqual(update_model.requests[0].body.client_sequence_number, 10)
            self.assertTrue(receipt["provider_effect_reconciled_by_sequence"])
            self.assertTrue(receipt["requested_absolute_quantity_not_promoted_to_internal_fact"])
            self.assertTrue((output / "MUTATION_ATTEMPT.json").exists())
            persisted = b"".join(path.read_bytes() for path in output.iterdir())
            for secret in (b"LOC A/1", b"LOC-B", b"SKU A?", b"SKU-B", b"MARKET-A", b"MARKET-B"):
                self.assertNotIn(secret, persisted)

    def test_06_live_drift_blocks_before_update_effect(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            service = FakeInventory()
            baseline, baseline_bytes = make_baseline(root, service)
            service.state[("LOC A/1", "SKU A?", "MARKET-A", "MFN")]["sellable"] = 99
            req = update_request(); req_bytes = encoded(req)
            with self.assertRaises(PermissionError):
                subject.update_inventory(service, req, req_bytes, approval(req_bytes, baseline_bytes), baseline, baseline_bytes, profile(), root / "update", now=NOW)
            self.assertEqual(len(service.calls), 2)

    def test_07_swapped_sequences_block_per_identity_before_effect(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            service = FakeInventory()
            baseline, baseline_bytes = make_baseline(root, service)
            req = update_request(sequence_a=20, sequence_b=10); req_bytes = encoded(req)
            with self.assertRaises(PermissionError):
                subject.update_inventory(service, req, req_bytes, approval(req_bytes, baseline_bytes), baseline, baseline_bytes, profile(), root / "update", now=NOW)
            self.assertEqual(len(service.calls), 1)

    def test_08_tampered_approval_blocks_before_provider(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            service = FakeInventory()
            baseline, baseline_bytes = make_baseline(root, service)
            req = update_request(); req_bytes = encoded(req)
            approve = approval(req_bytes, baseline_bytes); approve["request_sha256"] = "0" * 64
            with self.assertRaises(PermissionError):
                subject.update_inventory(service, req, req_bytes, approve, baseline, baseline_bytes, profile(), root / "update", now=NOW)
            self.assertEqual(len(service.calls), 1)

    def test_09_all_provider_rejections_are_sanitized_and_known(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            service = FakeInventory()
            baseline, baseline_bytes = make_baseline(root, service)
            service.reject_indexes = {0, 1}
            req = update_request(); req_bytes = encoded(req); output = root / "update"
            with self.assertRaises(PermissionError):
                subject.update_inventory(service, req, req_bytes, approval(req_bytes, baseline_bytes), baseline, baseline_bytes, profile(), output, now=NOW)
            self.assertTrue((output / "PROVIDER_REJECTED.json").exists())
            self.assertFalse((output / "UNKNOWN_EFFECT.json").exists())
            self.assertNotIn("SECRET REASON", (output / "PROVIDER_REJECTED.json").read_text("utf-8"))

    def test_10_partial_batch_records_partial_and_unknown_effect(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            service = FakeInventory()
            baseline, baseline_bytes = make_baseline(root, service)
            service.reject_indexes = {1}
            req = update_request(); req_bytes = encoded(req); output = root / "update"
            with self.assertRaisesRegex(RuntimeError, "effect is unknown"):
                subject.update_inventory(service, req, req_bytes, approval(req_bytes, baseline_bytes), baseline, baseline_bytes, profile(), output, now=NOW)
            self.assertTrue((output / "PARTIAL_EFFECT.json").exists())
            self.assertTrue((output / "UNKNOWN_EFFECT.json").exists())

    def test_11_provider_exception_records_unknown_without_message(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            service = FakeInventory()
            baseline, baseline_bytes = make_baseline(root, service)
            service.raise_error = RuntimeError("SECRET PROVIDER ERROR")
            service.raise_call_number = 3
            req = update_request(); req_bytes = encoded(req); output = root / "update"
            with self.assertRaisesRegex(RuntimeError, "effect is unknown: RuntimeError"):
                subject.update_inventory(service, req, req_bytes, approval(req_bytes, baseline_bytes), baseline, baseline_bytes, profile(), output, now=NOW)
            persisted = (output / "UNKNOWN_EFFECT.json").read_text("utf-8")
            self.assertNotIn("SECRET PROVIDER ERROR", persisted)
            self.assertIn('"automatic_retry_authorized": false', persisted)

    def test_12_existing_output_blocks_before_provider(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory); output = root / "exists"; output.mkdir()
            service = FakeInventory(); req = fetch_request()
            with self.assertRaises(FileExistsError):
                subject.fetch_inventory(service, req, encoded(req), profile(), output, now=NOW)
            self.assertEqual(service.calls, [])


if __name__ == "__main__":
    unittest.main(verbosity=2)
````

### FILE: `amazon_spapi_external_inventory/README.md`
```yaml
block_id: "PY-AMAZON-SPAPI-EXTERNAL-INVENTORY:README-MD:v1"
operation: CREATE
provenance: AUTHORED
source: "local operator contract with explicit upstream attribution and limits"
license: "LicenseRef-Workspace-Owner"
sha256: "39cc039ac4efe5894d6783bcf9ba879be4ee243a93f9e316a8ce99df0b983c52"
variables: []
secrets_allowed: false
```
````markdown
# Amazon External Fulfillment Inventory adapter

Adapter local sobre `amzn-sp-api==1.11.1` para `BatchInventoryApi.batch_inventory` v2024-09-11. Amazon publica este API para inventario absoluto por ubicación desde POS, ERP o WMS sobre programas external fulfillment admitidos.

## Procedencia

- SDK oficial Amazon fijado por wheel/source SHA-256; source Apache-2.0.
- Modelo oficial External Fulfillment Inventory fijado por commit. Su OpenAPI declara Amazon Software License; este pack la conserva como condición y no copia el modelo.
- Scripts, configuración y pruebas son `AUTHORED`; no se presentan como código escrito por Amazon.

## Flujo

Instale `requirements.lock` con `--require-hashes`, complete el perfil y ejecute primero `FETCH`. El receipt filtra location/SKU/marketplace y conserva sólo hashes, cantidades y sequence. Para `UPDATE`, use las sequence exactas del fetch, ligue aprobación a los bytes de request y receipt y ejecute una sola vez. El adapter vuelve a hacer FETCH antes del efecto, escribe attempt, envía el modelo oficial y reconcilia sequence mediante otro FETCH.

No reintente `UNKNOWN_EFFECT`. Un response aceptado no hace commit automático en ERP/POS/WMS. Cuenta, programa, roles, mapping, sequence authority, sandbox, cuota, privacidad y pruebas live continúan obligatorios.
````

## 6. Configuration surface

| Campo | Default seguro | Regla |
|---|---|---|
| `decision` | blocked | Sólo `APPROVED`. |
| programas/roles/mapping/sequence/sandbox | `UNPROVEN` | Todos deben ser `PROVEN`. |
| canales | vacío | Subset único de `FBA`,`MFN`,`DF`. |
| cuota/datos | false | Ambos true. |
| batch/cantidad | 10 / 2.000.000.000 | Límites locales no superables. |
| commit interno | false | Cualquier true falla. |
| secretos | nombres env | Valores sólo en environment. |

## 7. Dependency bill

| Package/tool | Pin | Uso | Licencia | Fuente |
|---|---|---|---|---|
| `amzn-sp-api` | 1.11.1 wheel SHA `a88e…70b0` | API/modelos | Apache-2.0 + condición ASL del modelo | Amazon/PyPI |
| dependencias Python | URLs+SHA exactos | HTTP/runtime | MIT/MPL/BSD | PyPI |
| CPython | 3.10–3.14 | runtime/test | PSF-2.0 | python.org |

## 8. Apply order

1. Materializar y verificar hashes; instalar lock con `--require-hashes`.
2. Completar y probar perfil/cuenta/programa/rol/mapping/sequence/sandbox/cuota/política.
3. Ejecutar FETCH y revisar receipt.
4. Construir UPDATE con sequence por identidad y approval externo ligado a bytes exactos.
5. Ejecutar una vez; resolver manualmente `UNKNOWN_EFFECT` antes de reintentar.
6. El ERP/POS/WMS sólo consume receipts reconciliados según su propia transacción y política.

## 9. Verification

```text
python -m pip install --require-hashes -r requirements.lock
python -m pip check
python -m unittest -v test_external_inventory.py
```

Éxito focal: `pip check` y 12/12 pruebas en árbol autor y reconstruido. Live requiere cuenta, programa, roles, mapping, sequence behavior, sandbox, cuota, privacidad, carga y reconciliación real.

## 10. Reconstruction evidence

`reconstruction_evidence/AMAZON_SPAPI_EXTERNAL_FULFILLMENT_INVENTORY_EXECUTION_INVENTORY_2026-08-31_V157.md` fija hashes/resultados. Los nueve archivos son `AUTHORED`; el SDK/modelo oficial no se atribuye al glue local. Estado `CONDITIONED` hasta gates live y del proyecto.
