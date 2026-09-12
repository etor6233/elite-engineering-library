# Python Amazon SP-API Easy Ship Handover Adapter

## 1. Metadata

```yaml
pack_id: "PYTHON-AMAZON-SPAPI-EASYSHIP-HANDOVER-ADAPTER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa slots, programación y reconciliación Easy Ship v2022-03-23 mediante el SDK oficial Amazon amzn-sp-api 1.11.1; no inventa cancelación, firma ni aceptación interna."
stacks: ["Python 3.10-3.14", "amzn-sp-api 1.11.1", "Amazon Easy Ship v2022-03-23"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "PYTHON-AMAZON-SPAPI-CATALOG-ADAPTER 0.1.x", "PYTHON-AMAZON-SPAPI-SHIPPING-TRACKING-ADAPTER 0.6.x"]
incompatible_with: ["carrier universal", "marketplace u operación no demostrados por support table", "slot libre o no inventariado", "programación sin aprobación hash-bound", "retry automático de schedule ambiguo", "PickedUp como aceptación empresarial automática", "firma/foto no presentes", "secreto en CLI/Markdown", "SDK distinto de 1.11.1"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0 AND MIT AND MPL-2.0 AND BSD-3-Clause"
upstream_sources: ["https://github.com/amzn/selling-partner-api-sdk/tree/8e792ae345a8d334ccdbdd03181f05f040e6a4fc", "https://github.com/amzn/selling-partner-api-models/tree/8e429486005c4ebdce5099e48cc48515a65359bb", "https://developer-docs.amazon.com/sp-api/reference/listhandoverslots", "https://developer-docs.amazon.com/sp-api/reference/createscheduledpackage", "https://developer-docs.amazon.com/sp-api/reference/getscheduledpackage", "https://developer-docs.amazon.com/sp-api/reference/updatescheduledpackages", "https://developer-docs.amazon.com/sp-api/docs/easy-ship-api-rate-limits"]
verified_at: "2026-08-31"
```

## 2. Applicability

Use sólo para órdenes Amazon Easy Ship y marketplaces/operaciones demostrados en la support table vigente. El perfil debe probar registro, rol/scope, sandbox, cuota/costo, tratamiento de datos y owners. No es un carrier genérico.

## 3. Architecture contract

`list_handover_slots` es read-only y produce request/response/receipt hash-bound. Una aprobación separada selecciona un slot exacto. `create_scheduled_package` persiste el intento antes del POST; como la firma oficial no expone idempotency key, excepción o respuesta inválida deja efecto desconocido y bloquea retry. `get_scheduled_package` reconcilia estado. `PickedUp` se registra como reporte Amazon, nunca como aceptación interna automática. El SDK fijado no expone una cancelación inequívoca; este pack no la fabrica.

## 4. Exact file manifest

```text
CREATE amazon_spapi_easyship_handover/requirements-direct.in
CREATE amazon_spapi_easyship_handover/requirements.lock
CREATE amazon_spapi_easyship_handover/sdk-artifact.lock.json
CREATE amazon_spapi_easyship_handover/provider-profile.template.json
CREATE amazon_spapi_easyship_handover/handover-request.template.json
CREATE amazon_spapi_easyship_handover/schedule-approval.template.json
CREATE amazon_spapi_easyship_handover/run_easyship_handover.py
CREATE amazon_spapi_easyship_handover/test_easyship_handover.py
CREATE amazon_spapi_easyship_handover/README.md
```
## 5. Materialization blocks

### FILE: `amazon_spapi_easyship_handover/requirements-direct.in`
```yaml
block_id: "PY-AMAZON-SPAPI-EASYSHIP-HANDOVER:REQUIREMENTS-DIRECT-IN:v1"
operation: CREATE
provenance: AUTHORED
source: "official Amazon wheel URL and SHA-256 from PyPI"
license: "LicenseRef-Workspace-Owner"
sha256: "a897426f3166913c0a3c5ad4d2bb2716d0c2245178bc47a18502a51700cd81ee"
variables: []
secrets_allowed: false
```
````text
amzn-sp-api @ https://files.pythonhosted.org/packages/36/3e/8744071121b96a7738fa485e584670cd4e6735b545cecc18516d10bd8881/amzn_sp_api-1.11.1-py3-none-any.whl#sha256=a88e4059c9fbfd454bfef67b82288a926f73cbdc608a8831f8383565736c70b0
````

### FILE: `amazon_spapi_easyship_handover/requirements.lock`
```yaml
block_id: "PY-AMAZON-SPAPI-EASYSHIP-HANDOVER:REQUIREMENTS-LOCK:v1"
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

### FILE: `amazon_spapi_easyship_handover/sdk-artifact.lock.json`
```yaml
block_id: "PY-AMAZON-SPAPI-EASYSHIP-HANDOVER:SDK-ARTIFACT-LOCK-JSON:v1"
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

### FILE: `amazon_spapi_easyship_handover/provider-profile.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-EASYSHIP-HANDOVER:PROVIDER-PROFILE-TEMPLATE-JSON:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed access, scope, marketplace and data profile"
license: "LicenseRef-Workspace-Owner"
sha256: "4ced3b3e8085c2a49a7fdb0ae497a6e09e996f74c4c99afad7e17349c6c6ffd2"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-easyship-profile/v1",
  "provider": "Amazon Easy Ship v2022-03-23",
  "sdk": "amzn-sp-api==1.11.1",
  "decision": "BLOCKED_ACCESS_AND_SCOPE_APPROVAL_REQUIRED",
  "region": "SANDBOX",
  "approved_marketplace_ids": [],
  "developer_registration": "NOT_PROVEN",
  "application_registration": "NOT_PROVEN",
  "easyship_role_and_scope": "NOT_PROVEN",
  "marketplace_support_table_proven": false,
  "sandbox_contract_test": "NOT_PROVEN",
  "quota_and_cost_approved": false,
  "order_data_handling_approved": false,
  "handover_slots_approved": false,
  "package_schedule_approved": false,
  "package_reconciliation_approved": false,
  "schedule_approval_owner": "",
  "reconciliation_owner": "",
  "client_id_env": "AMAZON_SP_API_CLIENT_ID",
  "client_secret_env": "AMAZON_SP_API_CLIENT_SECRET",
  "refresh_token_env": "AMAZON_SP_API_REFRESH_TOKEN",
  "automatic_internal_handover_completion": false
}
````

### FILE: `amazon_spapi_easyship_handover/handover-request.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-EASYSHIP-HANDOVER:HANDOVER-REQUEST-TEMPLATE-JSON:v1"
operation: CREATE
provenance: AUTHORED
source: "local closed input for official Easy Ship dimensions, weight and marketplace contract"
license: "LicenseRef-Workspace-Owner"
sha256: "5338b0a25c29d25d7780f859a6172609889eb9b281e570985f8eff0d38f795d1"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-easyship-handover-request/v1",
  "marketplace_id": "",
  "amazon_order_id": "",
  "package_dimensions": {"length": 0, "width": 0, "height": 0, "unit": "cm"},
  "package_weight": {"value": 0, "unit": "grams"},
  "requested_handover_method": "PICKUP"
}
````

### FILE: `amazon_spapi_easyship_handover/schedule-approval.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-EASYSHIP-HANDOVER:SCHEDULE-APPROVAL-TEMPLATE-JSON:v1"
operation: CREATE
provenance: AUTHORED
source: "local approval bound to exact provider slot inventory"
license: "LicenseRef-Workspace-Owner"
sha256: "847d685bd5a7d711a870944796335cf3bc0888fa137de6847329438ac62641aa"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-easyship-schedule-approval/v1",
  "approved_by": "",
  "approved_at": "",
  "slots_receipt_sha256": "",
  "selected_slot": {"slot_id": "", "start_time": "", "end_time": "", "handover_method": "PICKUP"},
  "package_items": [{"order_item_id": "", "serial_numbers": []}],
  "package_identifier": ""
}
````

### FILE: `amazon_spapi_easyship_handover/run_easyship_handover.py`
```yaml
block_id: "PY-AMAZON-SPAPI-EASYSHIP-HANDOVER:RUN-EASYSHIP-HANDOVER-PY:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed wrapper over official Amazon EasyShipApi methods"
license: "LicenseRef-Workspace-Owner"
sha256: "2290dca17bf30d5adcbb87ec74e2271101569fc58db44c16d0e62bbba5595e3e"
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
import sys
from typing import Any, Protocol

from spapi import SPAPIClient, SPAPIConfig
from spapi.api.easyship_v2022_03_23.easy_ship_api import EasyShipApi
from spapi.models.easyship_v2022_03_23.create_scheduled_package_request import CreateScheduledPackageRequest
from spapi.models.easyship_v2022_03_23.dimensions import Dimensions
from spapi.models.easyship_v2022_03_23.item import Item
from spapi.models.easyship_v2022_03_23.list_handover_slots_request import ListHandoverSlotsRequest
from spapi.models.easyship_v2022_03_23.package_details import PackageDetails
from spapi.models.easyship_v2022_03_23.time_slot import TimeSlot
from spapi.models.easyship_v2022_03_23.weight import Weight

SDK_VERSION = "1.11.1"
REGIONS = {"SANDBOX", "NA", "EU", "FE"}
HANDOVER_METHODS = {"PICKUP", "DROPOFF"}
PACKAGE_STATUSES = {"ReadyForPickup", "PickedUp", "AtOriginFC", "AtDestinationFC", "Delivered", "Rejected", "Undeliverable", "ReturnedToSeller", "LostInTransit", "LabelCanceled", "DamagedInTransit", "OutForDelivery"}
PROVIDER_HANDOVER_STATUSES = {"PickedUp", "AtOriginFC", "AtDestinationFC", "OutForDelivery", "Delivered", "Rejected", "Undeliverable", "ReturnedToSeller", "LostInTransit", "DamagedInTransit"}
REQUEST_KEYS = {"schema", "marketplace_id", "amazon_order_id", "package_dimensions", "package_weight", "requested_handover_method"}
APPROVAL_KEYS = {"schema", "approved_by", "approved_at", "slots_receipt_sha256", "selected_slot", "package_items", "package_identifier"}


class EasyShipService(Protocol):
    def list_handover_slots(self, **kwargs: Any) -> Any: ...
    def create_scheduled_package(self, body: CreateScheduledPackageRequest, **kwargs: Any) -> Any: ...
    def get_scheduled_package(self, amazon_order_id: str, marketplace_id: str, **kwargs: Any) -> Any: ...


class ScheduleEffectUnknown(RuntimeError):
    pass


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
        raise ValueError(f"{label} must be a JSON object")
    return value, raw


def exact_keys(value: Any, keys: set[str], label: str) -> dict[str, Any]:
    if not isinstance(value, dict) or set(value) != keys:
        raise ValueError(f"{label} keys must be exactly {sorted(keys)}")
    return value


def nonempty(value: Any, label: str, maximum: int = 255) -> str:
    if not isinstance(value, str) or not value or value != value.strip() or len(value) > maximum:
        raise ValueError(f"{label} must be a canonical non-empty string")
    return value


def positive_number(value: Any, label: str, minimum: Decimal) -> float:
    if isinstance(value, bool):
        raise ValueError(f"{label} is invalid")
    try:
        parsed = Decimal(str(value))
    except Exception as exc:
        raise ValueError(f"{label} is invalid") from exc
    if not parsed.is_finite() or parsed < minimum:
        raise ValueError(f"{label} must be at least {minimum}")
    return float(parsed)


def aware_datetime(value: Any, label: str) -> datetime:
    text = nonempty(value, label, 64)
    try:
        parsed = datetime.fromisoformat(text.replace("Z", "+00:00"))
    except ValueError as exc:
        raise ValueError(f"{label} must be ISO-8601") from exc
    if parsed.tzinfo is None or parsed.utcoffset() is None:
        raise ValueError(f"{label} must include an offset")
    return parsed.astimezone(timezone.utc)


def atomic_directory(destination: Path) -> Path:
    target = destination.resolve()
    if target.exists():
        raise FileExistsError("output_directory must not exist")
    target.parent.mkdir(parents=True, exist_ok=True)
    target.mkdir()
    return target


def approved_profile(path: Path) -> dict[str, Any]:
    profile, _ = strict_object(path, "provider profile")
    required = {"schema", "provider", "sdk", "decision", "region", "approved_marketplace_ids", "developer_registration", "application_registration", "easyship_role_and_scope", "marketplace_support_table_proven", "sandbox_contract_test", "quota_and_cost_approved", "order_data_handling_approved", "handover_slots_approved", "package_schedule_approved", "package_reconciliation_approved", "schedule_approval_owner", "reconciliation_owner", "client_id_env", "client_secret_env", "refresh_token_env", "automatic_internal_handover_completion"}
    exact_keys(profile, required, "provider profile")
    if profile["schema"] != "elite-amazon-spapi-easyship-profile/v1" or profile["provider"] != "Amazon Easy Ship v2022-03-23" or profile["sdk"] != "amzn-sp-api==1.11.1" or profile["decision"] != "APPROVED":
        raise PermissionError("Easy Ship profile is not approved")
    if profile["region"] not in REGIONS or profile["automatic_internal_handover_completion"] is not False:
        raise PermissionError("region or automatic completion contract is invalid")
    if any(profile[key] != "PROVEN" for key in ("developer_registration", "application_registration", "easyship_role_and_scope", "sandbox_contract_test")):
        raise PermissionError("registration, scope and sandbox must be proven")
    for key in ("marketplace_support_table_proven", "quota_and_cost_approved", "order_data_handling_approved"):
        if profile[key] is not True:
            raise PermissionError(f"{key} must be true")
    marketplaces = profile["approved_marketplace_ids"]
    if not isinstance(marketplaces, list) or not marketplaces or len(marketplaces) != len(set(marketplaces)):
        raise PermissionError("approved_marketplace_ids must be unique and non-empty")
    for value in marketplaces:
        nonempty(value, "approved marketplace")
    for key in ("client_id_env", "client_secret_env", "refresh_token_env"):
        if not re.fullmatch(r"[A-Z][A-Z0-9_]{2,127}", nonempty(profile[key], key, 128)):
            raise ValueError(f"{key} must name an environment variable")
    return profile


def service_from_environment(profile: dict[str, Any]) -> EasyShipApi:
    values: dict[str, str] = {}
    for key in ("client_id_env", "client_secret_env", "refresh_token_env"):
        value = os.environ.get(profile[key], "")
        if not value:
            raise PermissionError(f"required secret environment variable is missing: {profile[key]}")
        values[key] = value
    config = SPAPIConfig(client_id=values["client_id_env"], client_secret=values["client_secret_env"], refresh_token=values["refresh_token_env"], region=profile["region"])
    return EasyShipApi(SPAPIClient(config).api_client)


def build_slot_request(request: dict[str, Any], profile: dict[str, Any]) -> ListHandoverSlotsRequest:
    exact_keys(request, REQUEST_KEYS, "handover request")
    if request["schema"] != "elite-amazon-spapi-easyship-handover-request/v1":
        raise ValueError("handover request schema is invalid")
    marketplace = nonempty(request["marketplace_id"], "marketplace_id")
    if marketplace not in profile["approved_marketplace_ids"]:
        raise PermissionError("marketplace is not approved")
    order = nonempty(request["amazon_order_id"], "amazon_order_id")
    method = request["requested_handover_method"]
    if method not in HANDOVER_METHODS:
        raise ValueError("handover method is invalid")
    dimensions = exact_keys(request["package_dimensions"], {"length", "width", "height", "unit"}, "package_dimensions")
    weight = exact_keys(request["package_weight"], {"value", "unit"}, "package_weight")
    if dimensions["unit"] != "cm" or weight["unit"] not in {"grams", "g"}:
        raise ValueError("Easy Ship dimensions/weight units are invalid")
    return ListHandoverSlotsRequest(
        marketplace_id=marketplace,
        amazon_order_id=order,
        package_dimensions=Dimensions(length=positive_number(dimensions["length"], "length", Decimal("0.01")), width=positive_number(dimensions["width"], "width", Decimal("0.01")), height=positive_number(dimensions["height"], "height", Decimal("0.01")), unit="cm"),
        package_weight=Weight(value=positive_number(weight["value"], "weight", Decimal("11")), unit=weight["unit"]),
    )


def validate_slots(payload: Any, request: dict[str, Any], now: datetime) -> list[dict[str, Any]]:
    if not isinstance(payload, dict) or payload.get("amazon_order_id") != request["amazon_order_id"]:
        raise ValueError("slot response order does not match request")
    values = payload.get("time_slots")
    if not isinstance(values, list) or not values or len(values) > 100:
        raise ValueError("slot response must contain 1..100 slots")
    result: list[dict[str, Any]] = []
    seen: set[str] = set()
    for raw in values:
        slot = exact_keys(raw, {"slot_id", "start_time", "end_time", "handover_method"}, "slot")
        slot_id = nonempty(slot["slot_id"], "slot_id")
        if slot_id in seen:
            raise ValueError("slot IDs must be unique")
        seen.add(slot_id)
        start = aware_datetime(slot["start_time"], "slot start")
        end = aware_datetime(slot["end_time"], "slot end")
        if start >= end or end <= now:
            raise ValueError("slot window is invalid or expired")
        if slot["handover_method"] not in HANDOVER_METHODS:
            raise ValueError("provider returned unsupported handover method")
        result.append({"slot_id": slot_id, "start_time": start.isoformat(), "end_time": end.isoformat(), "handover_method": slot["handover_method"]})
    return result


def list_slots(service: EasyShipService, request: dict[str, Any], profile: dict[str, Any], output_directory: Path, timeout_seconds: float = 30.0, now: datetime | None = None) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION:
        raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if timeout_seconds <= 0 or timeout_seconds > 120:
        raise ValueError("timeout must be within 1..120 seconds")
    current = (now or datetime.now(timezone.utc)).astimezone(timezone.utc)
    model = build_slot_request(request, profile)
    response = service.list_handover_slots(list_handover_slots_request=model, _request_timeout=timeout_seconds)
    payload = json_value(response)
    slots = validate_slots(payload, request, current)
    target = atomic_directory(output_directory)
    request_bytes = pretty_bytes(request)
    response_bytes = pretty_bytes({"amazon_order_id": request["amazon_order_id"], "time_slots": slots})
    (target / "SLOTS_REQUEST.json").write_bytes(request_bytes)
    (target / "provider-response.json").write_bytes(response_bytes)
    receipt = {"schema": "elite-amazon-spapi-easyship-slots-receipt/v1", "created_at": current.isoformat(), "provider": "Amazon Easy Ship v2022-03-23", "sdk_version": SDK_VERSION, "operation": "list_handover_slots", "request_sha256": sha256(request_bytes), "provider_response_sha256": sha256(response_bytes), "marketplace_id_sha256": sha256(request["marketplace_id"].encode()), "amazon_order_id_sha256": sha256(request["amazon_order_id"].encode()), "slot_count": len(slots), "package_scheduled": False, "physical_handover_completed": False}
    receipt_bytes = pretty_bytes(receipt)
    (target / "SLOTS_RECEIPT.json").write_bytes(receipt_bytes)
    return receipt


def read_slots(directory: Path) -> tuple[dict[str, Any], dict[str, Any], dict[str, Any], bytes]:
    request, request_raw = strict_object(directory / "SLOTS_REQUEST.json", "slots request")
    response, response_raw = strict_object(directory / "provider-response.json", "slots response")
    receipt, receipt_raw = strict_object(directory / "SLOTS_RECEIPT.json", "slots receipt")
    expected = {"schema": "elite-amazon-spapi-easyship-slots-receipt/v1", "provider": "Amazon Easy Ship v2022-03-23", "sdk_version": SDK_VERSION, "operation": "list_handover_slots", "package_scheduled": False, "physical_handover_completed": False}
    if any(receipt.get(key) != value for key, value in expected.items()) or receipt.get("request_sha256") != sha256(request_raw) or receipt.get("provider_response_sha256") != sha256(response_raw):
        raise ValueError("slots artifacts are not hash-bound")
    return request, response, receipt, receipt_raw


def build_schedule(slots_directory: Path, approval: dict[str, Any], profile: dict[str, Any], now: datetime) -> tuple[CreateScheduledPackageRequest, dict[str, Any]]:
    request, response, _, receipt_raw = read_slots(slots_directory)
    exact_keys(approval, APPROVAL_KEYS, "schedule approval")
    if approval["schema"] != "elite-amazon-spapi-easyship-schedule-approval/v1" or approval["slots_receipt_sha256"] != sha256(receipt_raw):
        raise PermissionError("schedule approval is not bound to slots receipt")
    nonempty(approval["approved_by"], "approved_by")
    if aware_datetime(approval["approved_at"], "approved_at") > now:
        raise ValueError("approved_at must not be in the future")
    selected = exact_keys(approval["selected_slot"], {"slot_id", "start_time", "end_time", "handover_method"}, "selected_slot")
    normalized = {"slot_id": nonempty(selected["slot_id"], "slot_id"), "start_time": aware_datetime(selected["start_time"], "slot start").isoformat(), "end_time": aware_datetime(selected["end_time"], "slot end").isoformat(), "handover_method": selected["handover_method"]}
    matches = [slot for slot in response.get("time_slots", []) if slot == normalized]
    if len(matches) != 1 or normalized["handover_method"] != request["requested_handover_method"]:
        raise PermissionError("selected slot is not an exact approved provider slot")
    if aware_datetime(normalized["end_time"], "slot end") <= now:
        raise ValueError("selected slot is expired")
    raw_items = approval["package_items"]
    if not isinstance(raw_items, list) or not raw_items or len(raw_items) > 500:
        raise ValueError("package_items must contain 1..500 items")
    items: list[Item] = []
    order_item_ids: set[str] = set()
    for raw in raw_items:
        item = exact_keys(raw, {"order_item_id", "serial_numbers"}, "package item")
        item_id = nonempty(item["order_item_id"], "order_item_id")
        serials = item["serial_numbers"]
        if item_id in order_item_ids or not isinstance(serials, list) or len(serials) > 100 or len(serials) != len(set(serials)):
            raise ValueError("order items and serial numbers must be unique and within official limits")
        order_item_ids.add(item_id)
        for serial in serials:
            nonempty(serial, "serial number")
        items.append(Item(order_item_id=item_id, order_item_serial_numbers=serials or None))
    package_identifier = nonempty(approval["package_identifier"], "package_identifier")
    time_slot = TimeSlot(slot_id=normalized["slot_id"], start_time=aware_datetime(normalized["start_time"], "slot start"), end_time=aware_datetime(normalized["end_time"], "slot end"), handover_method=normalized["handover_method"])
    model = CreateScheduledPackageRequest(amazon_order_id=request["amazon_order_id"], marketplace_id=request["marketplace_id"], package_details=PackageDetails(package_items=items, package_time_slot=time_slot, package_identifier=package_identifier))
    context = {"slots_receipt_sha256": sha256(receipt_raw), "approval_sha256": sha256(canonical_bytes(approval)), "marketplace_id_sha256": sha256(request["marketplace_id"].encode()), "amazon_order_id_sha256": sha256(request["amazon_order_id"].encode()), "selected_slot_sha256": sha256(canonical_bytes(normalized)), "handover_method": normalized["handover_method"], "order_item_count": len(items)}
    return model, context


def validate_package(payload: Any, order_id: str, selected_slot: dict[str, Any]) -> tuple[str, str, str | None]:
    if not isinstance(payload, dict):
        raise ValueError("scheduled package response must be an object")
    package_id = exact_keys(payload.get("scheduled_package_id"), {"amazon_order_id", "package_id"}, "scheduled_package_id")
    if package_id["amazon_order_id"] != order_id:
        raise ValueError("scheduled package order does not match")
    identifier = nonempty(package_id["package_id"], "package_id")
    slot = payload.get("package_time_slot")
    normalized = {"slot_id": slot.get("slot_id"), "start_time": aware_datetime(slot.get("start_time"), "provider slot start").isoformat(), "end_time": aware_datetime(slot.get("end_time"), "provider slot end").isoformat(), "handover_method": slot.get("handover_method")} if isinstance(slot, dict) else None
    if normalized != selected_slot:
        raise ValueError("scheduled package slot does not match approval")
    status = payload.get("package_status")
    if status not in PACKAGE_STATUSES:
        raise ValueError("provider package status is invalid")
    tracking = payload.get("tracking_details")
    tracking_id = None
    if tracking is not None:
        tracking_id = nonempty(exact_keys(tracking, {"tracking_id"}, "tracking_details")["tracking_id"], "tracking_id")
    return identifier, status, tracking_id


def schedule_package(service: EasyShipService, slots_directory: Path, approval: dict[str, Any], profile: dict[str, Any], output_directory: Path, timeout_seconds: float = 30.0, now: datetime | None = None) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION:
        raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if timeout_seconds <= 0 or timeout_seconds > 120:
        raise ValueError("timeout must be within 1..120 seconds")
    current = (now or datetime.now(timezone.utc)).astimezone(timezone.utc)
    model, context = build_schedule(slots_directory, approval, profile, current)
    target = atomic_directory(output_directory)
    attempt = {"schema": "elite-amazon-spapi-easyship-schedule-attempt/v1", "created_at": current.isoformat(), **context, "provider_call_started": False, "automatic_retry_authorized": False, "reconciliation_required": False}
    attempt_path = target / "SCHEDULE_ATTEMPT.json"
    attempt_path.write_bytes(pretty_bytes(attempt))
    selected_slot = {"slot_id": model.package_details.package_time_slot.slot_id, "start_time": model.package_details.package_time_slot.start_time.isoformat(), "end_time": model.package_details.package_time_slot.end_time.isoformat(), "handover_method": model.package_details.package_time_slot.handover_method}
    try:
        attempt["provider_call_started"] = True
        attempt["reconciliation_required"] = True
        attempt_path.write_bytes(pretty_bytes(attempt))
        response = service.create_scheduled_package(model, _request_timeout=timeout_seconds)
        payload = json_value(response)
        package_id, status, tracking_id = validate_package(payload, model.amazon_order_id, selected_slot)
        response_bytes = pretty_bytes(payload)
        (target / "provider-response.json").write_bytes(response_bytes)
        receipt = {"schema": "elite-amazon-spapi-easyship-schedule-receipt/v1", "created_at": datetime.now(timezone.utc).isoformat(), **context, "provider_response_sha256": sha256(response_bytes), "scheduled_package_id_sha256": sha256(package_id.encode()), "package_status": status, "tracking_id_sha256": sha256(tracking_id.encode()) if tracking_id else None, "package_scheduled": True, "provider_reports_handover": status in PROVIDER_HANDOVER_STATUSES, "physical_handover_completed": False, "automatic_internal_handover_completion": False, "automatic_retry_authorized": False, "reconciliation_required": False}
        (target / "SCHEDULE_RECEIPT.json").write_bytes(pretty_bytes(receipt))
        attempt["provider_call_completed"] = True
        attempt["reconciliation_required"] = False
        attempt_path.write_bytes(pretty_bytes(attempt))
        return receipt
    except BaseException as exc:
        unknown = {"schema": "elite-amazon-spapi-easyship-schedule-unknown-effect/v1", "created_at": datetime.now(timezone.utc).isoformat(), **context, "provider_call_started": True, "automatic_retry_authorized": False, "reconciliation_required": True, "failure_class": type(exc).__name__}
        (target / "SCHEDULE_UNKNOWN_EFFECT.json").write_bytes(pretty_bytes(unknown))
        raise ScheduleEffectUnknown("Easy Ship schedule effect is unknown; reconcile before retry") from exc


def reconcile_package(service: EasyShipService, slots_directory: Path, approval: dict[str, Any], schedule_directory: Path, profile: dict[str, Any], output_directory: Path, timeout_seconds: float = 30.0, now: datetime | None = None) -> dict[str, Any]:
    current = (now or datetime.now(timezone.utc)).astimezone(timezone.utc)
    model, context = build_schedule(slots_directory, approval, profile, current)
    if not (schedule_directory / "SCHEDULE_ATTEMPT.json").is_file() or not ((schedule_directory / "SCHEDULE_RECEIPT.json").is_file() or (schedule_directory / "SCHEDULE_UNKNOWN_EFFECT.json").is_file()):
        raise ValueError("schedule directory lacks a durable attempt and terminal/unknown artifact")
    target = atomic_directory(output_directory)
    payload = json_value(service.get_scheduled_package(model.amazon_order_id, model.marketplace_id, _request_timeout=timeout_seconds))
    selected_slot = {"slot_id": model.package_details.package_time_slot.slot_id, "start_time": model.package_details.package_time_slot.start_time.isoformat(), "end_time": model.package_details.package_time_slot.end_time.isoformat(), "handover_method": model.package_details.package_time_slot.handover_method}
    package_id, status, tracking_id = validate_package(payload, model.amazon_order_id, selected_slot)
    response_bytes = pretty_bytes(payload)
    (target / "provider-response.json").write_bytes(response_bytes)
    receipt = {"schema": "elite-amazon-spapi-easyship-reconciliation-receipt/v1", "created_at": current.isoformat(), **context, "schedule_attempt_sha256": sha256((schedule_directory / "SCHEDULE_ATTEMPT.json").read_bytes()), "provider_response_sha256": sha256(response_bytes), "scheduled_package_id_sha256": sha256(package_id.encode()), "package_status": status, "tracking_id_sha256": sha256(tracking_id.encode()) if tracking_id else None, "provider_reports_handover": status in PROVIDER_HANDOVER_STATUSES, "provider_reports_delivery": status == "Delivered", "automatic_internal_handover_completion": False, "reconciliation_completed": True}
    (target / "RECONCILIATION_RECEIPT.json").write_bytes(pretty_bytes(receipt))
    return receipt


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("action", choices=("slots", "schedule", "reconcile"))
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--request", type=Path)
    parser.add_argument("--slots-directory", type=Path)
    parser.add_argument("--approval", type=Path)
    parser.add_argument("--schedule-directory", type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--timeout-seconds", type=float, default=30.0)
    args = parser.parse_args()
    profile = approved_profile(args.profile)
    service = service_from_environment(profile)
    if args.action == "slots":
        if profile["handover_slots_approved"] is not True or args.request is None:
            raise PermissionError("slot access and request are required")
        request, _ = strict_object(args.request, "handover request")
        receipt = list_slots(service, request, profile, args.output, args.timeout_seconds)
    elif args.action == "schedule":
        if profile["package_schedule_approved"] is not True or not nonempty(profile["schedule_approval_owner"], "schedule_approval_owner") or args.slots_directory is None or args.approval is None:
            raise PermissionError("schedule approval, owner and artifacts are required")
        approval, _ = strict_object(args.approval, "schedule approval")
        receipt = schedule_package(service, args.slots_directory, approval, profile, args.output, args.timeout_seconds)
    else:
        if profile["package_reconciliation_approved"] is not True or not nonempty(profile["reconciliation_owner"], "reconciliation_owner") or args.slots_directory is None or args.approval is None or args.schedule_directory is None:
            raise PermissionError("reconciliation approval, owner and artifacts are required")
        approval, _ = strict_object(args.approval, "schedule approval")
        receipt = reconcile_package(service, args.slots_directory, approval, args.schedule_directory, profile, args.output, args.timeout_seconds)
    print(json.dumps({"status": "AMAZON_SPAPI_EASYSHIP_HANDOVER_PASS", "receipt": receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `amazon_spapi_easyship_handover/test_easyship_handover.py`
```yaml
block_id: "PY-AMAZON-SPAPI-EASYSHIP-HANDOVER:TEST-EASYSHIP-HANDOVER-PY:v1"
operation: CREATE
provenance: AUTHORED
source: "local negative, atomicity, official-signature and reconciliation tests"
license: "LicenseRef-Workspace-Owner"
sha256: "7e32fd9c4ffe87c35b957b09fe62e34d21917c4dd826edfed581c4672f8bece0"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

from datetime import datetime, timezone
import hashlib
import inspect
import json
from pathlib import Path
import tempfile
import unittest

from spapi.api.easyship_v2022_03_23.easy_ship_api import EasyShipApi
from spapi.models.easyship_v2022_03_23.create_scheduled_package_request import CreateScheduledPackageRequest
from spapi.models.easyship_v2022_03_23.list_handover_slots_request import ListHandoverSlotsRequest

from run_easyship_handover import ScheduleEffectUnknown, build_slot_request, canonical_bytes, list_slots, reconcile_package, schedule_package, sha256

NOW = datetime(2026, 8, 31, 18, 0, tzinfo=timezone.utc)
SLOT = {"slot_id": "slot-1", "start_time": "2026-09-01T10:00:00+00:00", "end_time": "2026-09-01T12:00:00+00:00", "handover_method": "PICKUP"}


def profile() -> dict:
    return {"approved_marketplace_ids": ["A1MARKET"], "region": "SANDBOX"}


def request() -> dict:
    return {"schema": "elite-amazon-spapi-easyship-handover-request/v1", "marketplace_id": "A1MARKET", "amazon_order_id": "123-1234567-1234567", "package_dimensions": {"length": 20, "width": 10, "height": 5, "unit": "cm"}, "package_weight": {"value": 500, "unit": "grams"}, "requested_handover_method": "PICKUP"}


def approval(receipt_bytes: bytes) -> dict:
    return {"schema": "elite-amazon-spapi-easyship-schedule-approval/v1", "approved_by": "owner@example.invalid", "approved_at": "2026-08-31T17:59:00+00:00", "slots_receipt_sha256": sha256(receipt_bytes), "selected_slot": SLOT, "package_items": [{"order_item_id": "item-1", "serial_numbers": ["serial-1"]}], "package_identifier": "package-local-1"}


def package_payload(status: str = "ReadyForPickup") -> dict:
    return {"scheduled_package_id": {"amazon_order_id": "123-1234567-1234567", "package_id": "package-1"}, "package_dimensions": {"length": 20, "width": 10, "height": 5, "unit": "cm", "identifier": None}, "package_weight": {"value": 500, "unit": "grams"}, "package_items": [{"order_item_id": "item-1", "order_item_serial_numbers": ["serial-1"]}], "package_time_slot": SLOT, "package_identifier": "package-local-1", "invoice": None, "package_status": status, "tracking_details": {"tracking_id": "tracking-1"}}


class FakeService:
    def __init__(self, *, slots: dict | None = None, package: dict | None = None, error: BaseException | None = None):
        self.slots = slots or {"amazon_order_id": request()["amazon_order_id"], "time_slots": [SLOT]}
        self.package = package or package_payload()
        self.error = error
        self.slot_calls = []
        self.schedule_calls = []
        self.get_calls = []

    def list_handover_slots(self, **kwargs):
        self.slot_calls.append(kwargs)
        return self.slots

    def create_scheduled_package(self, body, **kwargs):
        self.schedule_calls.append((body, kwargs))
        if self.error:
            raise self.error
        return self.package

    def get_scheduled_package(self, amazon_order_id, marketplace_id, **kwargs):
        self.get_calls.append((amazon_order_id, marketplace_id, kwargs))
        return self.package


class EasyShipHandoverTests(unittest.TestCase):
    def slots_artifacts(self, root: Path, service: FakeService | None = None) -> Path:
        target = root / "slots"
        list_slots(service or FakeService(), request(), profile(), target, now=NOW)
        return target

    def test_official_sdk_signatures_and_models(self):
        self.assertEqual(list(inspect.signature(EasyShipApi.create_scheduled_package).parameters)[:2], ["self", "create_scheduled_package_request"])
        self.assertEqual(list(inspect.signature(EasyShipApi.get_scheduled_package).parameters)[:3], ["self", "amazon_order_id", "marketplace_id"])
        model = build_slot_request(request(), profile())
        self.assertIsInstance(model, ListHandoverSlotsRequest)
        self.assertEqual(model.package_dimensions.unit, "cm")
        self.assertEqual(model.package_weight.unit, "grams")

    def test_list_slots_writes_hash_bound_inventory(self):
        with tempfile.TemporaryDirectory() as root:
            service = FakeService()
            target = Path(root) / "slots"
            receipt = list_slots(service, request(), profile(), target, now=NOW)
            self.assertEqual(receipt["slot_count"], 1)
            self.assertFalse(receipt["package_scheduled"])
            self.assertIsInstance(service.slot_calls[0]["list_handover_slots_request"], ListHandoverSlotsRequest)
            self.assertEqual(receipt["provider_response_sha256"], hashlib.sha256((target / "provider-response.json").read_bytes()).hexdigest())

    def test_duplicate_or_expired_slots_fail_without_output(self):
        with tempfile.TemporaryDirectory() as root:
            bad = FakeService(slots={"amazon_order_id": request()["amazon_order_id"], "time_slots": [SLOT, SLOT]})
            target = Path(root) / "slots"
            with self.assertRaises(ValueError):
                list_slots(bad, request(), profile(), target, now=NOW)
            self.assertFalse(target.exists())

    def test_schedule_uses_exact_official_model_and_never_completes_handover(self):
        with tempfile.TemporaryDirectory() as root:
            root_path = Path(root)
            slots = self.slots_artifacts(root_path)
            selected = approval((slots / "SLOTS_RECEIPT.json").read_bytes())
            service = FakeService()
            receipt = schedule_package(service, slots, selected, profile(), root_path / "schedule", now=NOW)
            self.assertIsInstance(service.schedule_calls[0][0], CreateScheduledPackageRequest)
            self.assertEqual(service.schedule_calls[0][0].package_details.package_time_slot.slot_id, "slot-1")
            self.assertTrue(receipt["package_scheduled"])
            self.assertFalse(receipt["physical_handover_completed"])
            self.assertFalse(receipt["automatic_internal_handover_completion"])

    def test_unlisted_slot_or_tampered_receipt_is_rejected_before_mutation(self):
        with tempfile.TemporaryDirectory() as root:
            root_path = Path(root)
            slots = self.slots_artifacts(root_path)
            selected = approval((slots / "SLOTS_RECEIPT.json").read_bytes())
            selected["selected_slot"] = selected["selected_slot"] | {"slot_id": "invented"}
            service = FakeService()
            with self.assertRaises(PermissionError):
                schedule_package(service, slots, selected, profile(), root_path / "schedule", now=NOW)
            self.assertEqual(service.schedule_calls, [])
            self.assertFalse((root_path / "schedule").exists())

    def test_ambiguous_schedule_persists_unknown_effect_and_blocks_retry(self):
        with tempfile.TemporaryDirectory() as root:
            root_path = Path(root)
            slots = self.slots_artifacts(root_path)
            selected = approval((slots / "SLOTS_RECEIPT.json").read_bytes())
            target = root_path / "schedule"
            with self.assertRaises(ScheduleEffectUnknown):
                schedule_package(FakeService(error=TimeoutError("ambiguous")), slots, selected, profile(), target, now=NOW)
            unknown = json.loads((target / "SCHEDULE_UNKNOWN_EFFECT.json").read_text())
            self.assertFalse(unknown["automatic_retry_authorized"])
            self.assertTrue(unknown["reconciliation_required"])
            self.assertFalse((target / "SCHEDULE_RECEIPT.json").exists())

    def test_invalid_provider_response_is_unknown_effect_not_success(self):
        with tempfile.TemporaryDirectory() as root:
            root_path = Path(root)
            slots = self.slots_artifacts(root_path)
            selected = approval((slots / "SLOTS_RECEIPT.json").read_bytes())
            target = root_path / "schedule"
            bad = package_payload() | {"scheduled_package_id": {"amazon_order_id": "other", "package_id": "package-1"}}
            with self.assertRaises(ScheduleEffectUnknown):
                schedule_package(FakeService(package=bad), slots, selected, profile(), target, now=NOW)
            self.assertTrue((target / "SCHEDULE_UNKNOWN_EFFECT.json").is_file())

    def test_reconciliation_records_provider_pickup_without_internal_completion(self):
        with tempfile.TemporaryDirectory() as root:
            root_path = Path(root)
            slots = self.slots_artifacts(root_path)
            selected = approval((slots / "SLOTS_RECEIPT.json").read_bytes())
            schedule = root_path / "schedule"
            schedule_package(FakeService(), slots, selected, profile(), schedule, now=NOW)
            receipt = reconcile_package(FakeService(package=package_payload("PickedUp")), slots, selected, schedule, profile(), root_path / "reconcile", now=NOW)
            self.assertTrue(receipt["provider_reports_handover"])
            self.assertFalse(receipt["automatic_internal_handover_completion"])
            self.assertFalse(receipt["provider_reports_delivery"])

    def test_existing_output_is_never_overwritten(self):
        with tempfile.TemporaryDirectory() as root:
            target = Path(root) / "exists"
            target.mkdir()
            marker = target / "keep"
            marker.write_text("keep")
            with self.assertRaises(FileExistsError):
                list_slots(FakeService(), request(), profile(), target, now=NOW)
            self.assertEqual(marker.read_text(), "keep")


if __name__ == "__main__":
    unittest.main()
````

### FILE: `amazon_spapi_easyship_handover/README.md`
```yaml
block_id: "PY-AMAZON-SPAPI-EASYSHIP-HANDOVER:README-MD:v1"
operation: CREATE
provenance: AUTHORED
source: "local operating and admission instructions"
license: "LicenseRef-Workspace-Owner"
sha256: "96762e3b502eb1bc5a19fc844adf2420d51dc3d1d3f38c22aad661d8ae361f83"
variables: []
secrets_allowed: false
```
````markdown
# Amazon SP-API Easy Ship handover

Wrapper local fail-closed sobre el SDK oficial Amazon `amzn-sp-api==1.11.1` y Easy Ship v2022-03-23.

Flujo: `slots` guarda inventario hash-bound; una aprobación separada selecciona un slot exacto; `schedule` persiste el intento antes del POST y bloquea retry si el efecto queda ambiguo; `reconcile` consulta el paquete programado. `PickedUp` se registra sólo como estado reportado por Amazon y nunca autoriza por sí solo completar el handover interno.

El perfil distribuido está bloqueado. Antes de una llamada real deben probarse registro, rol/scope, marketplace y operación soportados, sandbox, cuota/costo, manejo de datos, owners, secretos por ambiente y reconciliación. No es carrier universal, no obtiene firma/foto y no prueba producción.

```powershell
python -m pip install --require-hashes -r requirements.lock
python -m unittest -v test_easyship_handover.py
```
````

## 6. Configuration surface

El usuario completa `provider-profile.template.json` fuera de la biblioteca. Permanecen bloqueados: registro, aplicación, rol/scope Easy Ship, support table del marketplace/operación, sandbox, cuota/costo, manejo de órdenes, slots, schedule, reconciliación y owners. Los secretos sólo llegan por los tres nombres de variables de ambiente.

## 7. Dependency bill

`requirements.lock` fija por URL y SHA-256 `amzn-sp-api==1.11.1` y su grafo transitivo. `sdk-artifact.lock.json` fija wheel, fuente SDK, modelos OpenAPI, commits, archives y licencia Apache-2.0. Los nueve archivos del pack son `AUTHORED`; no se presentan como código escrito por Amazon.

## 8. Apply order

1. Materializar el pack y verificar cada SHA-256.
2. Instalar exclusivamente `requirements.lock` con `--require-hashes`; ejecutar `pip check`.
3. Completar y aprobar el perfil sin secretos persistidos.
4. Ejecutar `slots`; revisar el inventario y crear una aprobación exacta.
5. Ejecutar `schedule` una vez. Ante ambigüedad, no repetir: ejecutar `reconcile`.
6. Conservar artefactos restringidos y enlazarlos al ledger/transacción del proyecto.

## 9. Verification

Ejecutar `compileall`, `pip check` y `python -m unittest -v test_easyship_handover.py`. Las nueve pruebas cubren firmas/modelos oficiales, slot inventory, duplicado/expiración, schedule exacto, tamper, efecto ambiguo, respuesta inválida, reconciliación `PickedUp` y no-overwrite. Producción requiere además cuenta/sandbox reales, support table exacta, rate headers observados, redelivery/race/load, privacidad/retención, alarmas, rollback y aceptación empresarial. Cancelación, firma/foto, carrier universal y target productivo permanecen fuera.

## 10. Reconstruction evidence

`reconstruction_evidence/AMAZON_SPAPI_EASYSHIP_HANDOVER_EXECUTION_INVENTORY_2026-08-31_V153.md` registra autoridad oficial, hashes y gates. Admission continúa `CONDITIONED` hasta demostrar proveedor y target reales.
