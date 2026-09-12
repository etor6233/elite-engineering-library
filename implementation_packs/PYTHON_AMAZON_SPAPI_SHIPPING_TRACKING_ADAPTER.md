# Python Amazon SP-API Shipping Tracking Adapter

## 1. Metadata

```yaml
pack_id: "PYTHON-AMAZON-SPAPI-SHIPPING-TRACKING-ADAPTER"
pack_version: "0.6.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa ciclo Shipping V2 incluyendo claim aprobado mediante SDK oficial Amazon amzn-sp-api 1.11.1; efectos ligados a purchase/quote/insured-value/proofs y sin afirmar handover o entrega."
stacks: ["Python 3.10-3.14", "amzn-sp-api 1.11.1", "Amazon Shipping V2"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "PYTHON-AMAZON-SPAPI-CATALOG-ADAPTER 0.1.x"]
incompatible_with: ["carrier universal", "tracking/formulario como prueba automática de handover/entrega", "compra sin quote íntegro/aprobación/cargo máximo", "cancelación desde shipment ID libre", "collection form sin inventario/selección aprobada", "claim desde tracking libre, valor no asegurado o proof host no aprobado", "retry automático de efecto mutante ambiguo", "cuenta/región/carrier no aprobados", "secreto en CLI/Markdown", "SDK distinto de 1.11.1"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0 AND MIT AND MPL-2.0 AND BSD-3-Clause"
upstream_sources: ["https://github.com/amzn/selling-partner-api-sdk/tree/8e792ae345a8d334ccdbdd03181f05f040e6a4fc", "https://github.com/amzn/selling-partner-api-models/tree/8e429486005c4ebdce5099e48cc48515a65359bb", "https://developer-docs.amazon.com/sp-api/reference/gettracking", "https://developer-docs.amazon.com/sp-api/reference/getrates", "https://developer-docs.amazon.com/sp-api/reference/purchaseshipment", "https://developer-docs.amazon.com/sp-api/reference/cancelshipment", "https://developer-docs.amazon.com/sp-api/reference/getshipmentdocuments", "https://developer-docs.amazon.com/sp-api/reference/getunmanifestedshipments", "https://developer-docs.amazon.com/sp-api/reference/generatecollectionform", "https://developer-docs.amazon.com/sp-api/reference/createclaim", "https://github.com/amzn/selling-partner-api-models/issues/5270"]
verified_at: "2026-08-31"
```

## 2. Applicability

Use sólo con acceso/business/carrier/sandbox/costo/datos/reconciliación probados. Claim requiere purchase+quote+rate-request+selection íntegros, approval, insured value y proof hosts exactos. Collection y ciclo anterior conservan límites. No es carrier genérico; form/claim no prueban pickup/entrega.

## 3. Architecture contract

Los wrappers llaman métodos oficiales del ciclo, incluido `create_claim(CreateClaimRequest)`. Claim deriva tracking/package desde purchase, verifica quote/request/selection, limita moneda/valor a insured value, allowlistea proof hosts y persiste attempt antes del POST. Amazon no expone idempotency key para claim: efecto ambiguo bloquea retry. Collection/ciclo conservan contratos; nada completa handover/entrega.

## 4. Exact file manifest

```text
CREATE amazon_spapi_shipping_tracking/requirements-direct.in
CREATE amazon_spapi_shipping_tracking/requirements.lock
CREATE amazon_spapi_shipping_tracking/sdk-artifact.lock.json
CREATE amazon_spapi_shipping_tracking/provider-profile.template.json
CREATE amazon_spapi_shipping_tracking/run_shipping_tracking.py
CREATE amazon_spapi_shipping_tracking/test_shipping_tracking.py
CREATE amazon_spapi_shipping_tracking/rate-request.template.json
CREATE amazon_spapi_shipping_tracking/run_shipping_rates.py
CREATE amazon_spapi_shipping_tracking/test_shipping_rates.py
CREATE amazon_spapi_shipping_tracking/purchase-selection.template.json
CREATE amazon_spapi_shipping_tracking/run_shipping_purchase.py
CREATE amazon_spapi_shipping_tracking/test_shipping_purchase.py
CREATE amazon_spapi_shipping_tracking/cancellation-approval.template.json
CREATE amazon_spapi_shipping_tracking/run_shipping_postpurchase.py
CREATE amazon_spapi_shipping_tracking/test_shipping_postpurchase.py
CREATE amazon_spapi_shipping_tracking/collection-inventory-request.template.json
CREATE amazon_spapi_shipping_tracking/collection-form-approval.template.json
CREATE amazon_spapi_shipping_tracking/run_shipping_collection.py
CREATE amazon_spapi_shipping_tracking/test_shipping_collection.py
CREATE amazon_spapi_shipping_tracking/claim-approval.template.json
CREATE amazon_spapi_shipping_tracking/run_shipping_claim.py
CREATE amazon_spapi_shipping_tracking/test_shipping_claim.py
CREATE amazon_spapi_shipping_tracking/README.md
```

## 5. Materialization blocks

### FILE: `amazon_spapi_shipping_tracking/requirements-direct.in`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:requirements-direct:v1"
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

### FILE: `amazon_spapi_shipping_tracking/requirements.lock`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:requirements-lock:v1"
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

### FILE: `amazon_spapi_shipping_tracking/sdk-artifact.lock.json`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:sdk-lock:v1"
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

### FILE: `amazon_spapi_shipping_tracking/provider-profile.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed access/scope/data profile"
license: "LicenseRef-Workspace-Owner"
sha256: "ad536ec3775038f30f373726bf7c24261faa14eb8faf18c330f39b0d29bde928"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-shipping-tracking-profile/v1",
  "provider": "Amazon Shipping V2",
  "sdk": "amzn-sp-api==1.11.1",
  "decision": "BLOCKED_ACCESS_AND_SCOPE_APPROVAL_REQUIRED",
  "region": "SANDBOX",
  "approved_carrier_ids": [],
  "developer_registration": "NOT_PROVEN",
  "application_registration": "NOT_PROVEN",
  "shipping_role_and_scope": "NOT_PROVEN",
  "sandbox_contract_test": "NOT_PROVEN",
  "quota_and_cost_approved": false,
  "tracking_data_handling_approved": false,
  "rate_quotes_approved": false,
  "shipment_purchase_approved": false,
  "shipment_cancellation_approved": false,
  "shipment_documents_approved": false,
  "postpurchase_data_handling_approved": false,
  "collection_inventory_approved": false,
  "collection_form_generation_approved": false,
  "collection_form_data_handling_approved": false,
  "claims_approved": false,
  "claim_data_handling_approved": false,
  "approved_claim_proof_hosts": [],
  "shipping_business_id": "",
  "purchase_approval_owner": "",
  "cancellation_approval_owner": "",
  "collection_form_approval_owner": "",
  "claim_approval_owner": "",
  "reconciliation_owner": "",
  "client_id_env": "AMAZON_SP_API_CLIENT_ID",
  "client_secret_env": "AMAZON_SP_API_CLIENT_SECRET",
  "refresh_token_env": "AMAZON_SP_API_REFRESH_TOKEN",
  "automatic_delivery_completion": false
}
````

### FILE: `amazon_spapi_shipping_tracking/run_shipping_tracking.py`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local wrapper invoking official Amazon ShippingApi"
license: "LicenseRef-Workspace-Owner"
sha256: "7b8195e3e929110b2fa7e4e4e9382f3d6778f793c806119a0b09789b6773d2f8"
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
import uuid
from typing import Any, Protocol

from spapi import ShippingApi, SPAPIClient, SPAPIConfig

SDK_VERSION = "1.11.1"
REGIONS = {"SANDBOX", "NA", "EU", "FE"}


class TrackingService(Protocol):
    def get_tracking(self, tracking_id: str, carrier_id: str, **kwargs: Any) -> Any: ...


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
    to_dict = getattr(value, "to_dict", None)
    if callable(to_dict):
        return json_value(to_dict())
    raise TypeError(f"unsupported provider response type: {type(value).__name__}")


def read_approved_profile(path: Path) -> dict[str, Any]:
    profile = json.loads(path.read_text(encoding="utf-8"))
    if profile.get("schema") != "elite-amazon-spapi-shipping-tracking-profile/v1":
        raise ValueError("unsupported profile schema")
    if profile.get("sdk") != f"amzn-sp-api=={SDK_VERSION}":
        raise ValueError("profile SDK version mismatch")
    if profile.get("decision") not in {"APPROVED_FOR_SANDBOX", "APPROVED_FOR_PRODUCTION"}:
        raise PermissionError("Amazon Shipping access decision is not approved")
    if profile.get("region") not in REGIONS:
        raise ValueError("unsupported Amazon SP-API region")
    for field in ("developer_registration", "application_registration", "shipping_role_and_scope", "sandbox_contract_test"):
        if profile.get(field) != "PROVEN":
            raise PermissionError(f"{field} must be proven")
    if profile.get("quota_and_cost_approved") is not True or profile.get("tracking_data_handling_approved") is not True:
        raise PermissionError("quota/cost and tracking-data handling approvals are required")
    carriers = profile.get("approved_carrier_ids")
    if not isinstance(carriers, list) or not carriers:
        raise ValueError("approved_carrier_ids must be a non-empty list")
    if any(not isinstance(value, str) or not re.fullmatch(r"[A-Za-z0-9][A-Za-z0-9_.-]{1,63}", value) for value in carriers):
        raise ValueError("approved_carrier_ids contains an invalid value")
    if not isinstance(profile.get("reconciliation_owner"), str) or not profile["reconciliation_owner"].strip():
        raise PermissionError("reconciliation_owner is required")
    if profile.get("automatic_delivery_completion") is not False:
        raise ValueError("tracking must not complete a business delivery automatically")
    for key, expected in {"client_id_env":"AMAZON_SP_API_CLIENT_ID","client_secret_env":"AMAZON_SP_API_CLIENT_SECRET","refresh_token_env":"AMAZON_SP_API_REFRESH_TOKEN"}.items():
        if profile.get(key) != expected:
            raise ValueError(f"{key} must use the canonical secret environment reference")
    return profile


def query_tracking(service: TrackingService, tracking_id: str, carrier_id: str, output_directory: Path, timeout_seconds: float = 30.0) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION:
        raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if not re.fullmatch(r"[A-Za-z0-9][A-Za-z0-9_.:/-]{4,127}", tracking_id):
        raise ValueError("tracking_id is invalid")
    if not re.fullmatch(r"[A-Za-z0-9][A-Za-z0-9_.-]{1,63}", carrier_id):
        raise ValueError("carrier_id is invalid")
    if timeout_seconds <= 0 or timeout_seconds > 120:
        raise ValueError("timeout must be within 1..120 seconds")
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".amazon-shipping-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        response = service.get_tracking(tracking_id, carrier_id, _request_timeout=timeout_seconds)
        payload = json_value(response)
        response_bytes = (json.dumps(payload, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "provider-response.json").write_bytes(response_bytes)
        receipt = {
            "schema": "elite-amazon-spapi-shipping-tracking-receipt/v1",
            "created_at": datetime.now(timezone.utc).isoformat(),
            "provider": "Amazon Shipping V2",
            "sdk_version": SDK_VERSION,
            "operation": "get_tracking",
            "tracking_id_sha256": sha256(tracking_id.encode("utf-8")),
            "carrier_id": carrier_id,
            "provider_response_sha256": sha256(response_bytes),
            "automatic_delivery_completion": False,
        }
        (stage / "TRACKING_RECEIPT.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        stage.replace(output)
        return receipt
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def service_from_environment(profile: dict[str, Any]) -> ShippingApi:
    values: dict[str, str] = {}
    for key in ("client_id_env", "client_secret_env", "refresh_token_env"):
        env_name = profile[key]
        value = os.environ.get(env_name, "")
        if not value:
            raise PermissionError(f"required secret environment variable is missing: {env_name}")
        values[key] = value
    config = SPAPIConfig(client_id=values["client_id_env"], client_secret=values["client_secret_env"], refresh_token=values["refresh_token_env"], region=profile["region"])
    return ShippingApi(SPAPIClient(config).api_client)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--tracking-id", required=True)
    parser.add_argument("--carrier-id", required=True)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--timeout-seconds", type=float, default=30.0)
    args = parser.parse_args()
    profile = read_approved_profile(args.profile)
    if args.carrier_id not in profile["approved_carrier_ids"]:
        raise PermissionError("carrier_id is not approved by the profile")
    receipt = query_tracking(service_from_environment(profile), args.tracking_id, args.carrier_id, args.output, args.timeout_seconds)
    print(json.dumps({"status":"AMAZON_SPAPI_SHIPPING_TRACKING_PASS","receipt":receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `amazon_spapi_shipping_tracking/test_shipping_tracking.py`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local negative/atomicity tests and official SDK signature probe"
license: "LicenseRef-Workspace-Owner"
sha256: "179665ccd8ff45e35245585e3761a84c2185a027409d39313c1d16ada5c6bb88"
variables: []
secrets_allowed: false
```
````python
import inspect
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch

from spapi import ShippingApi, SPAPIConfig
from run_shipping_tracking import query_tracking, read_approved_profile


class Result:
    def to_dict(self):
        return {"payload": {"trackingId": "SECRET-TRACK", "eventHistory": [{"eventCode": "Delivered"}]}}


class FakeService:
    def __init__(self, error=None):
        self.error = error
        self.calls = []

    def get_tracking(self, tracking_id, carrier_id, **kwargs):
        self.calls.append((tracking_id, carrier_id, kwargs))
        if self.error:
            raise self.error
        return Result()


def approved_profile():
    return {
        "schema":"elite-amazon-spapi-shipping-tracking-profile/v1",
        "provider":"Amazon Shipping V2",
        "sdk":"amzn-sp-api==1.11.1",
        "decision":"APPROVED_FOR_SANDBOX",
        "region":"SANDBOX",
        "approved_carrier_ids":["AMZN_US"],
        "developer_registration":"PROVEN",
        "application_registration":"PROVEN",
        "shipping_role_and_scope":"PROVEN",
        "sandbox_contract_test":"PROVEN",
        "quota_and_cost_approved":True,
        "tracking_data_handling_approved":True,
        "reconciliation_owner":"logistics-team",
        "client_id_env":"AMAZON_SP_API_CLIENT_ID",
        "client_secret_env":"AMAZON_SP_API_CLIENT_SECRET",
        "refresh_token_env":"AMAZON_SP_API_REFRESH_TOKEN",
        "automatic_delivery_completion":False,
    }


class ShippingTrackingTests(unittest.TestCase):
    def test_preserves_official_response_and_redacts_receipt(self):
        with tempfile.TemporaryDirectory() as root:
            output = Path(root) / "result"
            service = FakeService()
            receipt = query_tracking(service, "TRACK-12345", "AMZN_US", output)
            self.assertEqual(service.calls, [("TRACK-12345", "AMZN_US", {"_request_timeout":30.0})])
            self.assertIn("SECRET-TRACK", (output / "provider-response.json").read_text())
            self.assertNotIn("TRACK-12345", json.dumps(receipt))
            self.assertFalse(receipt["automatic_delivery_completion"])

    def test_invalid_inputs_and_existing_output_fail_before_provider(self):
        with tempfile.TemporaryDirectory() as root:
            service = FakeService()
            for tracking, carrier in [("bad", "AMZN_US"), ("TRACK-12345", "bad carrier")]:
                with self.assertRaises(ValueError):
                    query_tracking(service, tracking, carrier, Path(root) / (tracking + carrier).replace(" ", "_"))
            existing = Path(root) / "existing"
            existing.mkdir()
            with self.assertRaises(FileExistsError):
                query_tracking(service, "TRACK-12345", "AMZN_US", existing)
            self.assertEqual(service.calls, [])

    def test_provider_failure_is_atomic(self):
        with tempfile.TemporaryDirectory() as root:
            output = Path(root) / "result"
            with self.assertRaises(RuntimeError):
                query_tracking(FakeService(RuntimeError("provider unavailable")), "TRACK-12345", "AMZN_US", output)
            self.assertFalse(output.exists())
            self.assertEqual(list(Path(root).iterdir()), [])

    def test_profile_is_closed_and_exact(self):
        with tempfile.TemporaryDirectory() as root:
            path = Path(root) / "profile.json"
            blocked = approved_profile()
            blocked["decision"] = "BLOCKED_ACCESS_AND_SCOPE_APPROVAL_REQUIRED"
            path.write_text(json.dumps(blocked), encoding="utf-8")
            with self.assertRaises(PermissionError):
                read_approved_profile(path)
            valid = approved_profile()
            path.write_text(json.dumps(valid), encoding="utf-8")
            self.assertEqual(read_approved_profile(path)["approved_carrier_ids"], ["AMZN_US"])
            valid["automatic_delivery_completion"] = True
            path.write_text(json.dumps(valid), encoding="utf-8")
            with self.assertRaises(ValueError):
                read_approved_profile(path)

    def test_official_sdk_contract(self):
        self.assertEqual(ShippingApi.api_models_module, "spapi.models.shipping_v2")
        signature = inspect.signature(ShippingApi.get_tracking)
        self.assertEqual(list(signature.parameters)[:3], ["self", "tracking_id", "carrier_id"])
        self.assertIn("refresh_token", inspect.signature(SPAPIConfig).parameters)

    def test_version_mismatch_fails_before_provider(self):
        with tempfile.TemporaryDirectory() as root, patch("importlib.metadata.version", return_value="0.0.0"):
            service = FakeService()
            with self.assertRaises(RuntimeError):
                query_tracking(service, "TRACK-12345", "AMZN_US", Path(root) / "result")
            self.assertEqual(service.calls, [])


if __name__ == "__main__":
    unittest.main(verbosity=2)
````

### FILE: `amazon_spapi_shipping_tracking/rate-request.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:rate-request:v1"
operation: CREATE
provenance: AUTHORED
source: "local closed template for official Amazon Shipping V2 models"
license: "LicenseRef-Workspace-Owner"
sha256: "2323c854b38fcec0de1ee4c0971924fc8de36b83c11eb5ebb47e4f0c2baf565f"
variables: []
secrets_allowed: false
```
````json
{
  "ship_to": {"name":"","address_line1":"","city":"","state_or_region":"","country_code":"","postal_code":"","email":"","phone_number":""},
  "ship_from": {"name":"","address_line1":"","city":"","state_or_region":"","country_code":"","postal_code":"","email":"","phone_number":""},
  "return_to": {"name":"","address_line1":"","city":"","state_or_region":"","country_code":"","postal_code":"","email":"","phone_number":""},
  "ship_date": "",
  "package": {"dimensions":{"length":0,"width":0,"height":0,"unit":"CENTIMETER"},"weight":{"value":0,"unit":"KILOGRAM"},"insured_value":{"value":0,"unit":""},"is_hazmat":false,"seller_display_name":"","package_client_reference_id":"","items":[{"item_value":{"value":0,"unit":""},"description":"","item_identifier":"","quantity":0,"weight":{"value":0,"unit":"KILOGRAM"},"is_hazmat":false,"product_type":""}]}
}
````

### FILE: `amazon_spapi_shipping_tracking/run_shipping_rates.py`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:rates-runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local wrapper constructing and invoking official Amazon Shipping V2 models"
license: "LicenseRef-Workspace-Owner"
sha256: "65a66157e7537518480822e4f84e7647ccdc47f3650244a7398f3ac0e2ec1d8f"
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
from pathlib import Path
import re
import shutil
import sys
import uuid
from typing import Any, Protocol

from spapi import ShippingApi
from spapi.models.shipping_v2.address import Address
from spapi.models.shipping_v2.channel_details import ChannelDetails
from spapi.models.shipping_v2.channel_type import ChannelType
from spapi.models.shipping_v2.currency import Currency
from spapi.models.shipping_v2.dimensions import Dimensions
from spapi.models.shipping_v2.get_rates_request import GetRatesRequest
from spapi.models.shipping_v2.item import Item
from spapi.models.shipping_v2.package import Package
from spapi.models.shipping_v2.weight import Weight

from run_shipping_tracking import SDK_VERSION, json_value, read_approved_profile, service_from_environment, sha256

ADDRESS_KEYS = {"name","address_line1","address_line2","address_line3","company_name","state_or_region","city","country_code","postal_code","email","phone_number"}
ADDRESS_REQUIRED = {"name","address_line1","city","country_code","postal_code","email","phone_number"}
ROOT_KEYS = {"ship_to","ship_from","return_to","ship_date","package"}
PACKAGE_KEYS = {"dimensions","weight","insured_value","is_hazmat","seller_display_name","package_client_reference_id","items"}
DIMENSION_UNITS = {"INCH", "CENTIMETER"}
WEIGHT_UNITS = {"GRAM", "KILOGRAM", "OUNCE", "POUND"}


class RateService(Protocol):
    def get_rates(self, body: GetRatesRequest, **kwargs: Any) -> Any: ...


def exact_keys(value: dict[str, Any], allowed: set[str], required: set[str], label: str) -> None:
    unknown = set(value) - allowed
    missing = required - set(value)
    if unknown or missing:
        raise ValueError(f"{label} keys invalid; missing={sorted(missing)} unknown={sorted(unknown)}")


def nonempty(value: Any, label: str, maximum: int = 256) -> str:
    if not isinstance(value, str) or not value.strip() or len(value) > maximum:
        raise ValueError(f"{label} must be a non-empty bounded string")
    return value


def build_address(value: Any, label: str) -> Address:
    if not isinstance(value, dict):
        raise ValueError(f"{label} must be an object")
    exact_keys(value, ADDRESS_KEYS, ADDRESS_REQUIRED, label)
    country = nonempty(value["country_code"], f"{label}.country_code", 2)
    if not re.fullmatch(r"[A-Z]{2}", country):
        raise ValueError(f"{label}.country_code must be ISO alpha-2 uppercase")
    kwargs = {key: nonempty(item, f"{label}.{key}") for key, item in value.items() if item not in (None, "")}
    return Address(**kwargs)


def positive(value: Any, label: str) -> float:
    if isinstance(value, bool) or not isinstance(value, (int, float)) or value <= 0:
        raise ValueError(f"{label} must be positive")
    return float(value)


def build_currency(value: Any, label: str) -> Currency:
    if not isinstance(value, dict) or set(value) != {"value", "unit"}:
        raise ValueError(f"{label} keys are invalid")
    unit = nonempty(value["unit"], f"{label}.unit", 3)
    if not re.fullmatch(r"[A-Z]{3}", unit):
        raise ValueError(f"{label}.unit must be ISO-4217 uppercase")
    return Currency(value=positive(value["value"], f"{label}.value"), unit=unit)


def build_items(value: Any) -> list[Item]:
    keys = {"item_value","description","item_identifier","quantity","weight","is_hazmat","product_type"}
    if not isinstance(value, list) or not value or len(value) > 100:
        raise ValueError("items must contain 1..100 entries")
    result = []
    for index, item in enumerate(value):
        if not isinstance(item, dict):
            raise ValueError(f"items[{index}] must be an object")
        exact_keys(item, keys, keys, f"items[{index}]")
        weight = item["weight"]
        if not isinstance(weight, dict) or set(weight) != {"value","unit"} or weight["unit"] not in WEIGHT_UNITS:
            raise ValueError(f"items[{index}].weight is invalid")
        if isinstance(item["quantity"], bool) or not isinstance(item["quantity"], int) or item["quantity"] <= 0:
            raise ValueError(f"items[{index}].quantity must be positive integer")
        if not isinstance(item["is_hazmat"], bool):
            raise ValueError(f"items[{index}].is_hazmat must be boolean")
        result.append(Item(item_value=build_currency(item["item_value"], f"items[{index}].item_value"), description=nonempty(item["description"], f"items[{index}].description"), item_identifier=nonempty(item["item_identifier"], f"items[{index}].item_identifier", 128), quantity=item["quantity"], weight=Weight(value=positive(weight["value"], f"items[{index}].weight.value"), unit=weight["unit"]), is_hazmat=item["is_hazmat"], product_type=nonempty(item["product_type"], f"items[{index}].product_type", 64)))
    return result


def build_rate_request(payload: Any) -> GetRatesRequest:
    if not isinstance(payload, dict):
        raise ValueError("rate request must be an object")
    exact_keys(payload, ROOT_KEYS, ROOT_KEYS, "request")
    package = payload["package"]
    if not isinstance(package, dict):
        raise ValueError("package must be an object")
    exact_keys(package, PACKAGE_KEYS, PACKAGE_KEYS, "package")
    dimensions = package["dimensions"]
    weight = package["weight"]
    if not isinstance(dimensions, dict) or set(dimensions) != {"length","width","height","unit"}:
        raise ValueError("dimensions keys are invalid")
    if not isinstance(weight, dict) or set(weight) != {"value","unit"}:
        raise ValueError("weight keys are invalid")
    if dimensions["unit"] not in DIMENSION_UNITS or weight["unit"] not in WEIGHT_UNITS:
        raise ValueError("measurement unit is not admitted by the official model")
    if not isinstance(package["is_hazmat"], bool):
        raise ValueError("is_hazmat must be boolean")
    ship_date = nonempty(payload["ship_date"], "ship_date", 40)
    try:
        parsed = datetime.fromisoformat(ship_date.replace("Z", "+00:00"))
    except ValueError as exc:
        raise ValueError("ship_date must be ISO-8601 with timezone") from exc
    if parsed.tzinfo is None:
        raise ValueError("ship_date must include timezone")
    model_package = Package(
        dimensions=Dimensions(length=positive(dimensions["length"], "length"), width=positive(dimensions["width"], "width"), height=positive(dimensions["height"], "height"), unit=dimensions["unit"]),
        weight=Weight(value=positive(weight["value"], "weight"), unit=weight["unit"]),
        insured_value=build_currency(package["insured_value"], "insured_value"),
        is_hazmat=package["is_hazmat"],
        seller_display_name=nonempty(package["seller_display_name"], "seller_display_name", 128),
        package_client_reference_id=nonempty(package["package_client_reference_id"], "package_client_reference_id", 64),
        items=build_items(package["items"]),
    )
    return GetRatesRequest(ship_to=build_address(payload["ship_to"], "ship_to"), ship_from=build_address(payload["ship_from"], "ship_from"), return_to=build_address(payload["return_to"], "return_to"), ship_date=ship_date, packages=[model_package], channel_details=ChannelDetails(channel_type=ChannelType.EXTERNAL))


def quote_rates(service: RateService, payload: dict[str, Any], output_directory: Path, timeout_seconds: float = 30.0) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION:
        raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if timeout_seconds <= 0 or timeout_seconds > 120:
        raise ValueError("timeout must be within 1..120 seconds")
    model = build_rate_request(payload)
    canonical_request = (json.dumps(payload, ensure_ascii=False, sort_keys=True, separators=(",",":")) + "\n").encode("utf-8")
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".amazon-rates-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        response = service.get_rates(model, _request_timeout=timeout_seconds)
        response_bytes = (json.dumps(json_value(response), ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "provider-response.json").write_bytes(response_bytes)
        receipt = {"schema":"elite-amazon-spapi-shipping-rates-receipt/v1","created_at":datetime.now(timezone.utc).isoformat(),"provider":"Amazon Shipping V2","sdk_version":SDK_VERSION,"operation":"get_rates","request_sha256":sha256(canonical_request),"package_count":1,"provider_response_sha256":sha256(response_bytes),"shipment_purchased":False,"automatic_delivery_completion":False}
        (stage / "RATES_RECEIPT.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        stage.replace(output)
        return receipt
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--request", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--timeout-seconds", type=float, default=30.0)
    args = parser.parse_args()
    profile = read_approved_profile(args.profile)
    if profile.get("rate_quotes_approved") is not True:
        raise PermissionError("rate_quotes_approved is required")
    payload = json.loads(args.request.read_text(encoding="utf-8"))
    receipt = quote_rates(service_from_environment(profile), payload, args.output, args.timeout_seconds)
    print(json.dumps({"status":"AMAZON_SPAPI_SHIPPING_RATES_PASS","receipt":receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `amazon_spapi_shipping_tracking/test_shipping_rates.py`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:rates-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local official-model, closure and atomicity regressions"
license: "LicenseRef-Workspace-Owner"
sha256: "9b8551597edc934a9ccd4db6cf85f935587c2176346ca12b6adc775bc5e13b6f"
variables: []
secrets_allowed: false
```
````python
import json
from pathlib import Path
import tempfile
import unittest

from spapi.models.shipping_v2.get_rates_request import GetRatesRequest
from run_shipping_rates import build_rate_request, quote_rates


def request_payload():
    address = {"name":"Warehouse","address_line1":"1 Main St","city":"Seattle","state_or_region":"WA","country_code":"US","postal_code":"98101","email":"ops@example.invalid","phone_number":"+12025550100"}
    item = {"item_value":{"value":1000,"unit":"USD"},"description":"Electric mobility part","item_identifier":"SKU-0001","quantity":1,"weight":{"value":12.5,"unit":"KILOGRAM"},"is_hazmat":False,"product_type":"OTHER"}
    return {"ship_to":address | {"name":"Customer","address_line1":"2 Main St"},"ship_from":address,"return_to":address,"ship_date":"2026-09-01T10:00:00+00:00","package":{"dimensions":{"length":40,"width":30,"height":20,"unit":"CENTIMETER"},"weight":{"value":12.5,"unit":"KILOGRAM"},"insured_value":{"value":1000,"unit":"USD"},"is_hazmat":False,"seller_display_name":"Elite Mobility","package_client_reference_id":"shipment-0001","items":[item]}}


class Result:
    def to_dict(self):
        return {"payload":{"request_token":"request-token","rates":[{"rate_id":"rate-1","carrier_id":"AMZN_US"}]}}


class FakeService:
    def __init__(self, error=None):
        self.error = error
        self.calls = []
    def get_rates(self, body, **kwargs):
        self.calls.append((body, kwargs))
        if self.error:
            raise self.error
        return Result()


class ShippingRateTests(unittest.TestCase):
    def test_builds_official_models_and_one_package(self):
        model = build_rate_request(request_payload())
        self.assertIsInstance(model, GetRatesRequest)
        self.assertEqual(len(model.packages), 1)
        self.assertEqual(model.packages[0].dimensions.unit, "CENTIMETER")
        self.assertEqual(model.packages[0].weight.unit, "KILOGRAM")

    def test_quote_is_atomic_and_never_purchases(self):
        with tempfile.TemporaryDirectory() as root:
            service = FakeService()
            output = Path(root) / "result"
            receipt = quote_rates(service, request_payload(), output)
            self.assertEqual(len(service.calls), 1)
            self.assertIsInstance(service.calls[0][0], GetRatesRequest)
            self.assertFalse(receipt["shipment_purchased"])
            self.assertFalse(receipt["automatic_delivery_completion"])
            self.assertTrue((output / "provider-response.json").is_file())
            self.assertNotIn("Customer", json.dumps(receipt))

    def test_unknown_keys_invalid_units_values_and_dates_fail(self):
        cases = []
        value = request_payload(); value["unexpected"] = True; cases.append(value)
        value = request_payload(); value["package"]["dimensions"]["unit"] = "METER"; cases.append(value)
        value = request_payload(); value["package"]["weight"]["value"] = 0; cases.append(value)
        value = request_payload(); value["ship_date"] = "2026-09-01"; cases.append(value)
        value = request_payload(); value["ship_to"]["country_code"] = "usa"; cases.append(value)
        for value in cases:
            with self.subTest(value=value), self.assertRaises(ValueError):
                build_rate_request(value)

    def test_provider_failure_removes_stage(self):
        with tempfile.TemporaryDirectory() as root:
            output = Path(root) / "result"
            with self.assertRaises(RuntimeError):
                quote_rates(FakeService(RuntimeError("provider failed")), request_payload(), output)
            self.assertFalse(output.exists())
            self.assertEqual(list(Path(root).iterdir()), [])

    def test_existing_output_and_version_are_closed(self):
        with tempfile.TemporaryDirectory() as root:
            output = Path(root) / "result"; output.mkdir()
            service = FakeService()
            with self.assertRaises(FileExistsError):
                quote_rates(service, request_payload(), output)
            self.assertEqual(service.calls, [])


if __name__ == "__main__":
    unittest.main(verbosity=2)
````

### FILE: `amazon_spapi_shipping_tracking/README.md`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local bounded runbook with official authority boundaries"
license: "LicenseRef-Workspace-Owner"
sha256: "279e82ef2a05a74af37d6d7c738211cc22e4f36a09afd899e8ef1ccef9fdc969"
variables: []
secrets_allowed: false
```
````markdown
# Amazon SP-API Shipping V2 tracking

This module invokes Amazon's official `amzn-sp-api` Python SDK 1.11.1 and its generated `ShippingApi.get_tracking(tracking_id, carrier_id)`, `ShippingApi.get_rates(body)` and `ShippingApi.purchase_shipment(body, x_amzn_idempotency_key=..., x_amzn_shipping_business_id=...)` contracts. The SDK is Amazon code under Apache-2.0; the wrapper/profile/tests are local `AUTHORED` integration and are not represented as Amazon-authored.

The distributed profile is blocked. Before use, prove developer/application registration, Shipping role/scope, supported region/carrier, exact sandbox behavior, quotas/cost, tracking-data handling, reconciliation ownership and secret injection. Install only `requirements.lock` with hashes honored by the package installer.

The output preserves the official provider response in a new atomic directory and hashes the tracking ID in the receipt. Tracking evidence never completes a business delivery automatically. Provider response can contain customer/location data and must remain under project retention/access/redaction controls.

Rate quotes construct Amazon's official `Address`, `Dimensions`, `Weight`, `Package` and `GetRatesRequest` models. The request schema admits exactly one package because an open Amazon model issue dated 2026-06-12 reports that Shipping V2 rejects multiple packages although the OpenAPI array has no `maxItems`; this conservative boundary must be re-audited against an official fix. Quote output does not purchase a shipment.

Purchase is a separate state-changing command. It accepts only a V148 quote whose response hash and receipt hash still match, only inside Amazon's official ten-minute rate window, and only when the selected carrier, charge and document options are exactly approved. It constructs `RequestedDocumentSpecification` and `PurchaseShipmentRequest` from the offered rate. Before the POST it persists a deterministic `x-amzn-IdempotencyKey`. Any exception after that boundary leaves `PURCHASE_UNKNOWN_EFFECT.json`, blocks automatic retry and requires reconciliation with the same key.

Post-purchase uses only the completed purchase artifact. Cancellation obtains `shipmentId` from its hash-bound response, persists an attempt before Amazon's official `cancel_shipment` PUT and blocks automatic retry because that operation exposes no idempotency key. Document retrieval obtains shipment/package IDs from purchase, re-verifies the original purchase selection, calls official `get_shipment_documents`, validates Base64/type/format/cardinality/size and writes hash-named files plus a receipt atomically.

Collection handover is two-step. `get_unmanifested_shipments` first inventories carrier/from-address locations with no mutation. A separate approval binds that response and one exact location; only then `generate_collection_form` runs with Amazon's idempotency key. The returned Base64 PDF is validated and hash-named. Its receipt always states `physical_handover_completed=false`: generating a form is not carrier acceptance.

Claims are derived from purchase + quote + original rate request + purchase selection, never a free tracking ID. Declared value cannot exceed the request's insured value/currency and proof URLs must use exact approved HTTPS hosts. Amazon's `create_claim` exposes no idempotency key, so an attempt is persisted before POST and every ambiguous response blocks retry pending reconciliation.

Local gate:

```text
python -m unittest -v test_shipping_tracking.py
python -m unittest -v test_shipping_rates.py
python -m unittest -v test_shipping_purchase.py
python -m unittest -v test_shipping_postpurchase.py
python -m unittest -v test_shipping_collection.py
python -m unittest -v test_shipping_claim.py
```

Production additionally requires a real account/sandbox call, rate-limit/retry policy from observed headers, reconciliation against known shipment/claim, duplicate/out-of-order tests, carrier/geography support, security/privacy/legal review, load, alarms, rollback and business acceptance. Purchase/cancel/docs/collection/claim are implemented but blocked; pickup booking and physical carrier acceptance remain outside this version. It does not support arbitrary carriers or prove delivery.
````

### FILE: `amazon_spapi_shipping_tracking/purchase-selection.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:purchase-selection:v1"
operation: CREATE
provenance: AUTHORED
source: "local exact quote/rate/charge/document approval contract over official Amazon getRates output"
license: "LicenseRef-Workspace-Owner"
sha256: "e62b3b1f7b83a6663069cfd5e7f34dbf57a2c30da5de79ec696b0ac400872dac"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-shipping-purchase-selection/v1",
  "approved_by": "",
  "approved_at": "",
  "rates_receipt_sha256": "",
  "rate_id": "",
  "maximum_charge": {
    "value": "0.00",
    "unit": "USD"
  },
  "document": {
    "format": "PDF",
    "size": {
      "width": 4,
      "length": 6,
      "unit": "INCH"
    },
    "dpi": 300,
    "page_layout": "DEFAULT",
    "need_file_joining": false,
    "requested_document_types": ["LABEL"]
  }
}
````

### FILE: `amazon_spapi_shipping_tracking/run_shipping_purchase.py`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:purchase-runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed wrapper over official Amazon ShippingApi.purchase_shipment and generated request models"
license: "LicenseRef-Workspace-Owner"
sha256: "e348f13a4090d78fac295522f838867f8abcf68838ef758d319b1ecb1a98b177"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
from datetime import datetime, timezone
from decimal import Decimal, InvalidOperation
import importlib.metadata
import json
from pathlib import Path
import re
import sys
from typing import Any, Callable, Protocol
import uuid

from spapi.models.shipping_v2.document_size import DocumentSize
from spapi.models.shipping_v2.purchase_shipment_request import PurchaseShipmentRequest
from spapi.models.shipping_v2.requested_document_specification import RequestedDocumentSpecification

from run_shipping_tracking import SDK_VERSION, json_value, read_approved_profile, service_from_environment, sha256

SELECTION_KEYS = {"schema", "approved_by", "approved_at", "rates_receipt_sha256", "rate_id", "maximum_charge", "document"}
DOCUMENT_KEYS = {"format", "size", "dpi", "page_layout", "need_file_joining", "requested_document_types"}
SIZE_KEYS = {"width", "length", "unit"}
HEX64 = re.compile(r"[0-9a-f]{64}")


class PurchaseService(Protocol):
    def purchase_shipment(self, body: PurchaseShipmentRequest, **kwargs: Any) -> Any: ...


class PurchaseEffectUnknown(RuntimeError):
    pass


def exact_keys(value: Any, required: set[str], label: str) -> dict[str, Any]:
    if not isinstance(value, dict) or set(value) != required:
        raise ValueError(f"{label} keys must be exactly {sorted(required)}")
    return value


def nonempty(value: Any, label: str, maximum: int = 256) -> str:
    if not isinstance(value, str) or not value.strip() or len(value) > maximum:
        raise ValueError(f"{label} must be a non-empty bounded string")
    return value


def aware_datetime(value: Any, label: str) -> datetime:
    text = nonempty(value, label, 40)
    try:
        result = datetime.fromisoformat(text.replace("Z", "+00:00"))
    except ValueError as exc:
        raise ValueError(f"{label} must be ISO-8601") from exc
    if result.tzinfo is None:
        raise ValueError(f"{label} must include timezone")
    return result.astimezone(timezone.utc)


def decimal_value(value: Any, label: str, *, positive: bool = False) -> Decimal:
    if isinstance(value, bool) or not isinstance(value, (str, int, float)):
        raise ValueError(f"{label} must be decimal-compatible")
    try:
        result = Decimal(str(value))
    except InvalidOperation as exc:
        raise ValueError(f"{label} is invalid") from exc
    if not result.is_finite() or (positive and result <= 0):
        raise ValueError(f"{label} must be finite and positive")
    return result


def read_json(path: Path, label: str) -> tuple[dict[str, Any], bytes]:
    raw = path.read_bytes()
    try:
        value = json.loads(raw.decode("utf-8"))
    except (UnicodeDecodeError, json.JSONDecodeError) as exc:
        raise ValueError(f"{label} must be canonical UTF-8 JSON") from exc
    if not isinstance(value, dict):
        raise ValueError(f"{label} must be an object")
    return value, raw


def verify_quote(quote_directory: Path, selection: dict[str, Any], now: datetime) -> tuple[dict[str, Any], bytes, bytes]:
    receipt, receipt_bytes = read_json(quote_directory / "RATES_RECEIPT.json", "rates receipt")
    response, response_bytes = read_json(quote_directory / "provider-response.json", "rates response")
    required = {"schema", "created_at", "provider", "sdk_version", "operation", "request_sha256", "package_count", "provider_response_sha256", "shipment_purchased", "automatic_delivery_completion"}
    exact_keys(receipt, required, "rates receipt")
    if receipt != (receipt | {"schema":"elite-amazon-spapi-shipping-rates-receipt/v1", "provider":"Amazon Shipping V2", "sdk_version":SDK_VERSION, "operation":"get_rates", "package_count":1, "shipment_purchased":False, "automatic_delivery_completion":False}):
        raise ValueError("rates receipt contract mismatch")
    for field in ("request_sha256", "provider_response_sha256"):
        if not isinstance(receipt[field], str) or not HEX64.fullmatch(receipt[field]):
            raise ValueError(f"rates receipt {field} is invalid")
    if sha256(response_bytes) != receipt["provider_response_sha256"]:
        raise ValueError("rates response hash mismatch")
    if sha256(receipt_bytes) != selection["rates_receipt_sha256"]:
        raise PermissionError("selection is not bound to the exact rates receipt")
    created = aware_datetime(receipt["created_at"], "rates receipt created_at")
    if now < created or (now - created).total_seconds() >= 600:
        raise PermissionError("Amazon rate token is outside the official ten-minute purchase window")
    return response, response_bytes, receipt_bytes


def document_request(rate: dict[str, Any], selected: dict[str, Any]) -> RequestedDocumentSpecification:
    exact_keys(selected, DOCUMENT_KEYS, "document")
    size = exact_keys(selected["size"], SIZE_KEYS, "document.size")
    candidates = rate.get("supported_document_specifications")
    if not isinstance(candidates, list) or not candidates:
        raise ValueError("selected rate has no supported document specifications")
    matches = [item for item in candidates if isinstance(item, dict) and item.get("format") == selected["format"] and item.get("size") == size]
    if len(matches) != 1:
        raise ValueError("document format and size must select exactly one offered specification")
    options = matches[0].get("print_options")
    if not isinstance(options, list) or not options:
        raise ValueError("selected document specification has no print options")
    supported_dpis: set[int] = set()
    page_layouts: set[str] = set()
    joining: set[bool] = set()
    supported_types: set[str] = set()
    mandatory_types: set[str] = set()
    for option in options:
        if not isinstance(option, dict):
            raise ValueError("provider print option is invalid")
        supported_dpis.update(item for item in option.get("supported_dpis", []) if isinstance(item, int) and not isinstance(item, bool))
        page_layouts.update(item for item in option.get("supported_page_layouts", []) if isinstance(item, str))
        joining.update(item for item in option.get("supported_file_joining_options", []) if isinstance(item, bool))
        for detail in option.get("supported_document_details", []):
            if not isinstance(detail, dict) or not isinstance(detail.get("name"), str) or not isinstance(detail.get("is_mandatory"), bool):
                raise ValueError("provider supported document detail is invalid")
            supported_types.add(detail["name"])
            if detail["is_mandatory"]:
                mandatory_types.add(detail["name"])
    dpi = selected["dpi"]
    if dpi not in supported_dpis or selected["page_layout"] not in page_layouts or selected["need_file_joining"] not in joining:
        raise ValueError("selected print option was not offered by getRates")
    types = selected["requested_document_types"]
    if not isinstance(types, list) or not types or len(types) != len(set(types)) or any(not isinstance(item, str) for item in types):
        raise ValueError("requested_document_types must be a non-empty unique string list")
    if not set(types).issubset(supported_types) or not mandatory_types.issubset(set(types)):
        raise ValueError("requested document types do not satisfy the offered specification")
    return RequestedDocumentSpecification(format=selected["format"], size=DocumentSize(width=size["width"], length=size["length"], unit=size["unit"]), dpi=dpi, page_layout=selected["page_layout"], need_file_joining=selected["need_file_joining"], requested_document_types=types)


def prepare_purchase(quote_directory: Path, selection: dict[str, Any], profile: dict[str, Any], now: datetime) -> tuple[PurchaseShipmentRequest, dict[str, Any]]:
    exact_keys(selection, SELECTION_KEYS, "selection")
    if selection["schema"] != "elite-amazon-spapi-shipping-purchase-selection/v1":
        raise ValueError("unsupported selection schema")
    nonempty(selection["approved_by"], "approved_by")
    approved_at = aware_datetime(selection["approved_at"], "approved_at")
    if approved_at > now:
        raise ValueError("approved_at must not be in the future")
    if not isinstance(selection["rates_receipt_sha256"], str) or not HEX64.fullmatch(selection["rates_receipt_sha256"]):
        raise ValueError("rates_receipt_sha256 is invalid")
    rate_id = nonempty(selection["rate_id"], "rate_id", 256)
    response, response_bytes, receipt_bytes = verify_quote(quote_directory, selection, now)
    payload = response.get("payload")
    if not isinstance(payload, dict):
        raise ValueError("rates response payload is invalid")
    request_token = nonempty(payload.get("request_token"), "request_token", 512)
    rates = payload.get("rates")
    if not isinstance(rates, list):
        raise ValueError("rates response rates is invalid")
    matches = [rate for rate in rates if isinstance(rate, dict) and rate.get("rate_id") == rate_id]
    if len(matches) != 1:
        raise ValueError("rate_id must identify exactly one offered rate")
    rate = matches[0]
    if rate.get("carrier_id") not in profile["approved_carrier_ids"]:
        raise PermissionError("selected rate carrier is not approved")
    if rate.get("requires_additional_inputs") is not False:
        raise PermissionError("selected rate requires getAdditionalInputs and is not admitted by this lane")
    maximum = exact_keys(selection["maximum_charge"], {"value", "unit"}, "maximum_charge")
    unit = nonempty(maximum["unit"], "maximum_charge.unit", 3)
    if not re.fullmatch(r"[A-Z]{3}", unit):
        raise ValueError("maximum_charge.unit must be ISO-4217 uppercase")
    total = exact_keys(rate.get("total_charge"), {"value", "unit"}, "selected rate total_charge")
    if total["unit"] != unit or decimal_value(total["value"], "selected rate charge", positive=True) > decimal_value(maximum["value"], "maximum charge", positive=True):
        raise PermissionError("selected rate exceeds the exact approved currency/maximum")
    document = document_request(rate, selection["document"])
    model = PurchaseShipmentRequest(request_token=request_token, rate_id=rate_id, requested_document_specification=document)
    fingerprint = sha256(response_bytes + b"\x00" + receipt_bytes + b"\x00" + (json.dumps(selection, sort_keys=True, separators=(",", ":")) + "\n").encode("utf-8"))
    idempotency_key = str(uuid.uuid5(uuid.NAMESPACE_URL, "elite-amazon-shipping-purchase:" + fingerprint))
    context = {"quote_response_sha256":sha256(response_bytes), "rates_receipt_sha256":sha256(receipt_bytes), "selection_sha256":sha256((json.dumps(selection, sort_keys=True, separators=(",", ":")) + "\n").encode("utf-8")), "rate_id_sha256":sha256(rate_id.encode("utf-8")), "idempotency_key":idempotency_key}
    return model, context


def purchase(service: PurchaseService, quote_directory: Path, selection: dict[str, Any], profile: dict[str, Any], output_directory: Path, timeout_seconds: float = 30.0, clock: Callable[[], datetime] | None = None) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION:
        raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if timeout_seconds <= 0 or timeout_seconds > 120:
        raise ValueError("timeout must be within 1..120 seconds")
    now = (clock or (lambda: datetime.now(timezone.utc)))().astimezone(timezone.utc)
    model, context = prepare_purchase(quote_directory, selection, profile, now)
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist; reconcile an existing attempt instead of creating another")
    output.parent.mkdir(parents=True, exist_ok=True)
    output.mkdir()
    attempt = {"schema":"elite-amazon-spapi-shipping-purchase-attempt/v1", "created_at":now.isoformat(), "provider":"Amazon Shipping V2", "sdk_version":SDK_VERSION, "operation":"purchase_shipment", **context, "provider_call_started":False, "automatic_retry_authorized":False, "reconciliation_required":False}
    attempt_path = output / "PURCHASE_ATTEMPT.json"
    attempt_path.write_text(json.dumps(attempt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
    try:
        attempt["provider_call_started"] = True
        attempt["reconciliation_required"] = True
        attempt_path.write_text(json.dumps(attempt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        response = service.purchase_shipment(model, x_amzn_idempotency_key=context["idempotency_key"], x_amzn_shipping_business_id=profile["shipping_business_id"], _request_timeout=timeout_seconds)
        payload = json_value(response)
        response_bytes = (json.dumps(payload, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        result = payload.get("payload") if isinstance(payload, dict) else None
        if not isinstance(result, dict) or not nonempty(result.get("shipment_id"), "purchase response shipment_id", 256) or not isinstance(result.get("package_document_details"), list) or not result["package_document_details"]:
            raise ValueError("purchase response lacks official required shipment/package document details")
        (output / "provider-response.json").write_bytes(response_bytes)
        receipt = {"schema":"elite-amazon-spapi-shipping-purchase-receipt/v1", "created_at":datetime.now(timezone.utc).isoformat(), "provider":"Amazon Shipping V2", "sdk_version":SDK_VERSION, "operation":"purchase_shipment", **context, "provider_response_sha256":sha256(response_bytes), "shipment_id_sha256":sha256(result["shipment_id"].encode("utf-8")), "shipment_purchased":True, "automatic_delivery_completion":False, "reconciliation_required":False}
        (output / "PURCHASE_RECEIPT.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        attempt["reconciliation_required"] = False
        attempt["provider_call_completed"] = True
        attempt_path.write_text(json.dumps(attempt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        return receipt
    except BaseException as exc:
        unknown = {"schema":"elite-amazon-spapi-shipping-purchase-unknown-effect/v1", "created_at":datetime.now(timezone.utc).isoformat(), **context, "provider_call_started":True, "automatic_retry_authorized":False, "reconciliation_required":True, "failure_class":type(exc).__name__}
        (output / "PURCHASE_UNKNOWN_EFFECT.json").write_text(json.dumps(unknown, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        raise PurchaseEffectUnknown("purchase effect is unknown; reconcile with the same idempotency key before any retry") from exc


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--quote-directory", required=True, type=Path)
    parser.add_argument("--selection", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--timeout-seconds", type=float, default=30.0)
    args = parser.parse_args()
    profile = read_approved_profile(args.profile)
    if profile.get("shipment_purchase_approved") is not True:
        raise PermissionError("shipment_purchase_approved is required")
    for field in ("shipping_business_id", "purchase_approval_owner"):
        nonempty(profile.get(field), field)
    selection = json.loads(args.selection.read_text(encoding="utf-8"))
    receipt = purchase(service_from_environment(profile), args.quote_directory, selection, profile, args.output, args.timeout_seconds)
    print(json.dumps({"status":"AMAZON_SPAPI_SHIPPING_PURCHASE_PASS", "receipt":receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `amazon_spapi_shipping_tracking/test_shipping_purchase.py`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:purchase-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local offline contract tests against exact official Amazon model classes and method kwargs"
license: "LicenseRef-Workspace-Owner"
sha256: "a4ed2d36792b8b7c7a28f2d200d2ce147e10a2810039db8a03cd559723722ac1"
variables: []
secrets_allowed: false
```
````python
import json
from datetime import datetime, timezone
from pathlib import Path
import tempfile
import unittest

from spapi.models.shipping_v2.purchase_shipment_request import PurchaseShipmentRequest
from run_shipping_purchase import PurchaseEffectUnknown, prepare_purchase, purchase, sha256

NOW = datetime(2026, 8, 31, 15, 5, tzinfo=timezone.utc)


def profile():
    return {"approved_carrier_ids":["AMZN_US"], "shipping_business_id":"AmazonShipping_US"}


def selection(receipt_bytes):
    return {"schema":"elite-amazon-spapi-shipping-purchase-selection/v1", "approved_by":"shipping-owner@example.invalid", "approved_at":"2026-08-31T15:04:00+00:00", "rates_receipt_sha256":sha256(receipt_bytes), "rate_id":"rate-1", "maximum_charge":{"value":"12.50","unit":"USD"}, "document":{"format":"PDF", "size":{"width":4,"length":6,"unit":"INCH"}, "dpi":300, "page_layout":"DEFAULT", "need_file_joining":False, "requested_document_types":["LABEL"]}}


def quote(root: Path):
    directory = root / "quote"; directory.mkdir()
    response = {"payload":{"request_token":"request-token", "rates":[{"rate_id":"rate-1", "carrier_id":"AMZN_US", "total_charge":{"value":"12.00","unit":"USD"}, "requires_additional_inputs":False, "supported_document_specifications":[{"format":"PDF", "size":{"width":4,"length":6,"unit":"INCH"}, "print_options":[{"supported_dpis":[203,300], "supported_page_layouts":["DEFAULT"], "supported_file_joining_options":[False], "supported_document_details":[{"name":"LABEL","is_mandatory":True},{"name":"PACKSLIP","is_mandatory":False}]}]}]}]}}
    response_bytes = (json.dumps(response, sort_keys=True, indent=2) + "\n").encode()
    (directory / "provider-response.json").write_bytes(response_bytes)
    receipt = {"schema":"elite-amazon-spapi-shipping-rates-receipt/v1", "created_at":"2026-08-31T15:00:00+00:00", "provider":"Amazon Shipping V2", "sdk_version":"1.11.1", "operation":"get_rates", "request_sha256":"a"*64, "package_count":1, "provider_response_sha256":sha256(response_bytes), "shipment_purchased":False, "automatic_delivery_completion":False}
    receipt_bytes = (json.dumps(receipt, sort_keys=True, indent=2) + "\n").encode()
    (directory / "RATES_RECEIPT.json").write_bytes(receipt_bytes)
    return directory, receipt_bytes


class Result:
    def to_dict(self):
        return {"payload":{"shipment_id":"shipment-1", "package_document_details":[{"package_client_reference_id":"package-1", "package_documents":[{"type":"LABEL"}], "tracking_id":"tracking-1"}], "promise":{}}}


class FakeService:
    def __init__(self, error=None): self.error, self.calls = error, []
    def purchase_shipment(self, body, **kwargs):
        self.calls.append((body, kwargs))
        if self.error: raise self.error
        return Result()


class PurchaseTests(unittest.TestCase):
    def test_builds_exact_official_request_from_bound_quote(self):
        with tempfile.TemporaryDirectory() as root:
            directory, receipt_bytes = quote(Path(root))
            model, context = prepare_purchase(directory, selection(receipt_bytes), profile(), NOW)
            self.assertIsInstance(model, PurchaseShipmentRequest)
            self.assertEqual(model.request_token, "request-token")
            self.assertEqual(model.rate_id, "rate-1")
            self.assertEqual(model.requested_document_specification.requested_document_types, ["LABEL"])
            self.assertRegex(context["idempotency_key"], r"^[0-9a-f-]{36}$")

    def test_success_persists_response_and_exact_idempotency_header(self):
        with tempfile.TemporaryDirectory() as root:
            directory, receipt_bytes = quote(Path(root)); service = FakeService(); output = Path(root) / "purchase"
            receipt = purchase(service, directory, selection(receipt_bytes), profile(), output, clock=lambda: NOW)
            self.assertTrue(receipt["shipment_purchased"])
            self.assertFalse(receipt["automatic_delivery_completion"])
            self.assertEqual(service.calls[0][1]["x_amzn_idempotency_key"], receipt["idempotency_key"])
            self.assertEqual(service.calls[0][1]["x_amzn_shipping_business_id"], "AmazonShipping_US")
            self.assertTrue((output / "PURCHASE_ATTEMPT.json").is_file())
            self.assertTrue((output / "provider-response.json").is_file())

    def test_provider_failure_is_durable_unknown_effect_not_retryable(self):
        with tempfile.TemporaryDirectory() as root:
            directory, receipt_bytes = quote(Path(root)); output = Path(root) / "purchase"
            with self.assertRaises(PurchaseEffectUnknown):
                purchase(FakeService(TimeoutError("ambiguous")), directory, selection(receipt_bytes), profile(), output, clock=lambda: NOW)
            unknown = json.loads((output / "PURCHASE_UNKNOWN_EFFECT.json").read_text())
            self.assertFalse(unknown["automatic_retry_authorized"])
            self.assertTrue(unknown["reconciliation_required"])
            self.assertFalse((output / "PURCHASE_RECEIPT.json").exists())

    def test_tamper_expiry_charge_additional_inputs_and_document_fail_before_call(self):
        mutators = [
            lambda s, q: s.update(rates_receipt_sha256="0"*64),
            lambda s, q: s["maximum_charge"].update(value="1.00"),
            lambda s, q: q["payload"]["rates"][0].update(requires_additional_inputs=True),
            lambda s, q: s["document"].update(dpi=600),
        ]
        for mutator in mutators:
            with self.subTest(mutator=mutator), tempfile.TemporaryDirectory() as root:
                directory, receipt_bytes = quote(Path(root)); selected = selection(receipt_bytes)
                response_path = directory / "provider-response.json"; response = json.loads(response_path.read_text())
                mutator(selected, response)
                if response != json.loads(response_path.read_text()):
                    response_bytes = (json.dumps(response, sort_keys=True, indent=2) + "\n").encode(); response_path.write_bytes(response_bytes)
                    receipt = json.loads((directory / "RATES_RECEIPT.json").read_text()); receipt["provider_response_sha256"] = sha256(response_bytes)
                    receipt_bytes = (json.dumps(receipt, sort_keys=True, indent=2) + "\n").encode(); (directory / "RATES_RECEIPT.json").write_bytes(receipt_bytes); selected["rates_receipt_sha256"] = sha256(receipt_bytes)
                service = FakeService()
                with self.assertRaises((ValueError, PermissionError)):
                    purchase(service, directory, selected, profile(), Path(root) / "purchase", clock=lambda: NOW)
                self.assertEqual(service.calls, [])
        with tempfile.TemporaryDirectory() as root:
            directory, receipt_bytes = quote(Path(root)); service = FakeService()
            with self.assertRaises(PermissionError):
                purchase(service, directory, selection(receipt_bytes), profile(), Path(root) / "purchase", clock=lambda: datetime(2026,8,31,15,10,tzinfo=timezone.utc))
            self.assertEqual(service.calls, [])

    def test_existing_output_blocks_second_purchase(self):
        with tempfile.TemporaryDirectory() as root:
            directory, receipt_bytes = quote(Path(root)); output = Path(root) / "purchase"; output.mkdir(); service = FakeService()
            with self.assertRaises(FileExistsError):
                purchase(service, directory, selection(receipt_bytes), profile(), output, clock=lambda: NOW)
            self.assertEqual(service.calls, [])


if __name__ == "__main__":
    unittest.main(verbosity=2)
````

### FILE: `amazon_spapi_shipping_tracking/cancellation-approval.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:cancellation-approval:v1"
operation: CREATE
provenance: AUTHORED
source: "local approval binding cancellation to the exact completed Amazon purchase receipt"
license: "LicenseRef-Workspace-Owner"
sha256: "dc87f254ec001397fcf07b7383a5c34595b9cefa69f26bffbd1e0e0c1fc2a214"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-shipping-cancellation-approval/v1",
  "approved_by": "",
  "approved_at": "",
  "purchase_receipt_sha256": "",
  "reason_code": "CUSTOMER_REQUEST"
}
````

### FILE: `amazon_spapi_shipping_tracking/run_shipping_postpurchase.py`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:postpurchase-runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed wrapper over official Amazon cancel_shipment and get_shipment_documents methods"
license: "LicenseRef-Workspace-Owner"
sha256: "79c2aafae920fb362619e19918496e13c65885060df672ab731d52d04e11fa5d"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
import base64
import binascii
from datetime import datetime, timezone
import importlib.metadata
import json
from pathlib import Path
import re
import shutil
import sys
from typing import Any, Protocol
import uuid

from run_shipping_purchase import HEX64, aware_datetime, exact_keys, nonempty
from run_shipping_tracking import SDK_VERSION, json_value, read_approved_profile, service_from_environment, sha256

PURCHASE_RECEIPT_KEYS = {"schema", "created_at", "provider", "sdk_version", "operation", "quote_response_sha256", "rates_receipt_sha256", "selection_sha256", "rate_id_sha256", "idempotency_key", "provider_response_sha256", "shipment_id_sha256", "shipment_purchased", "automatic_delivery_completion", "reconciliation_required"}
CANCEL_APPROVAL_KEYS = {"schema", "approved_by", "approved_at", "purchase_receipt_sha256", "reason_code"}
CANCEL_REASONS = {"CUSTOMER_REQUEST", "OPERATIONAL_ERROR", "RATE_REJECTED"}
FORMAT_EXTENSIONS = {"PDF":"pdf", "PNG":"png", "ZPL":"zpl"}
MAX_DOCUMENTS = 10
MAX_DOCUMENT_BYTES = 20 * 1024 * 1024
MAX_TOTAL_BYTES = 50 * 1024 * 1024


class PostPurchaseService(Protocol):
    def cancel_shipment(self, shipment_id: str, **kwargs: Any) -> Any: ...
    def get_shipment_documents(self, shipment_id: str, package_client_reference_id: str, **kwargs: Any) -> Any: ...


class CancellationEffectUnknown(RuntimeError):
    pass


def canonical_bytes(value: Any) -> bytes:
    return (json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(",", ":")) + "\n").encode("utf-8")


def read_object(path: Path, label: str) -> tuple[dict[str, Any], bytes]:
    raw = path.read_bytes()
    try:
        value = json.loads(raw.decode("utf-8"))
    except (UnicodeDecodeError, json.JSONDecodeError) as exc:
        raise ValueError(f"{label} must be UTF-8 JSON") from exc
    if not isinstance(value, dict):
        raise ValueError(f"{label} must be an object")
    return value, raw


def verify_purchase(purchase_directory: Path) -> tuple[str, list[dict[str, Any]], dict[str, Any], bytes]:
    receipt, receipt_bytes = read_object(purchase_directory / "PURCHASE_RECEIPT.json", "purchase receipt")
    exact_keys(receipt, PURCHASE_RECEIPT_KEYS, "purchase receipt")
    expected = {"schema":"elite-amazon-spapi-shipping-purchase-receipt/v1", "provider":"Amazon Shipping V2", "sdk_version":SDK_VERSION, "operation":"purchase_shipment", "shipment_purchased":True, "automatic_delivery_completion":False, "reconciliation_required":False}
    if any(receipt.get(key) != value for key, value in expected.items()):
        raise ValueError("purchase receipt is not a completed Amazon purchase")
    for field in ("quote_response_sha256", "rates_receipt_sha256", "selection_sha256", "rate_id_sha256", "provider_response_sha256", "shipment_id_sha256"):
        if not isinstance(receipt[field], str) or not HEX64.fullmatch(receipt[field]):
            raise ValueError(f"purchase receipt {field} is invalid")
    response, response_bytes = read_object(purchase_directory / "provider-response.json", "purchase response")
    if sha256(response_bytes) != receipt["provider_response_sha256"]:
        raise ValueError("purchase response hash mismatch")
    payload = response.get("payload")
    if not isinstance(payload, dict):
        raise ValueError("purchase response payload is invalid")
    shipment_id = nonempty(payload.get("shipment_id"), "shipment_id", 256)
    if sha256(shipment_id.encode("utf-8")) != receipt["shipment_id_sha256"]:
        raise ValueError("shipment ID is not bound to the purchase receipt")
    details = payload.get("package_document_details")
    if not isinstance(details, list) or not details or any(not isinstance(item, dict) for item in details):
        raise ValueError("purchase response package details are invalid")
    return shipment_id, details, receipt, receipt_bytes


def cancel_shipment(service: PostPurchaseService, purchase_directory: Path, approval: dict[str, Any], profile: dict[str, Any], output_directory: Path, timeout_seconds: float = 30.0, now: datetime | None = None) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION:
        raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if timeout_seconds <= 0 or timeout_seconds > 120:
        raise ValueError("timeout must be within 1..120 seconds")
    exact_keys(approval, CANCEL_APPROVAL_KEYS, "cancellation approval")
    if approval["schema"] != "elite-amazon-spapi-shipping-cancellation-approval/v1" or approval["reason_code"] not in CANCEL_REASONS:
        raise ValueError("cancellation approval schema/reason is invalid")
    nonempty(approval["approved_by"], "approved_by")
    current = (now or datetime.now(timezone.utc)).astimezone(timezone.utc)
    if aware_datetime(approval["approved_at"], "approved_at") > current:
        raise ValueError("approved_at must not be in the future")
    shipment_id, _, _, receipt_bytes = verify_purchase(purchase_directory)
    if not isinstance(approval["purchase_receipt_sha256"], str) or sha256(receipt_bytes) != approval["purchase_receipt_sha256"]:
        raise PermissionError("cancellation is not bound to the exact purchase receipt")
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist; reconcile the existing cancellation attempt")
    output.parent.mkdir(parents=True, exist_ok=True); output.mkdir()
    context = {"purchase_receipt_sha256":sha256(receipt_bytes), "shipment_id_sha256":sha256(shipment_id.encode("utf-8")), "approval_sha256":sha256(canonical_bytes(approval)), "reason_code":approval["reason_code"]}
    attempt = {"schema":"elite-amazon-spapi-shipping-cancellation-attempt/v1", "created_at":current.isoformat(), **context, "provider_call_started":False, "automatic_retry_authorized":False, "reconciliation_required":False}
    attempt_path = output / "CANCELLATION_ATTEMPT.json"
    attempt_path.write_text(json.dumps(attempt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
    try:
        attempt["provider_call_started"] = True; attempt["reconciliation_required"] = True
        attempt_path.write_text(json.dumps(attempt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        response = service.cancel_shipment(shipment_id, x_amzn_shipping_business_id=profile["shipping_business_id"], _request_timeout=timeout_seconds)
        payload = json_value(response)
        if not isinstance(payload, dict) or set(payload) != {"payload"} or payload["payload"] not in ({}, None):
            raise ValueError("cancelShipment did not return the documented empty payload")
        response_bytes = (json.dumps(payload, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (output / "provider-response.json").write_bytes(response_bytes)
        receipt = {"schema":"elite-amazon-spapi-shipping-cancellation-receipt/v1", "created_at":datetime.now(timezone.utc).isoformat(), **context, "provider_response_sha256":sha256(response_bytes), "shipment_cancelled":True, "automatic_retry_authorized":False, "reconciliation_required":False}
        (output / "CANCELLATION_RECEIPT.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        attempt["provider_call_completed"] = True; attempt["reconciliation_required"] = False
        attempt_path.write_text(json.dumps(attempt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        return receipt
    except BaseException as exc:
        unknown = {"schema":"elite-amazon-spapi-shipping-cancellation-unknown-effect/v1", "created_at":datetime.now(timezone.utc).isoformat(), **context, "provider_call_started":True, "automatic_retry_authorized":False, "reconciliation_required":True, "failure_class":type(exc).__name__}
        (output / "CANCELLATION_UNKNOWN_EFFECT.json").write_text(json.dumps(unknown, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        raise CancellationEffectUnknown("cancellation effect is unknown; reconcile before retry") from exc


def retrieve_documents(service: PostPurchaseService, purchase_directory: Path, purchase_selection: dict[str, Any], package_reference: str, profile: dict[str, Any], output_directory: Path, timeout_seconds: float = 30.0) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION:
        raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if timeout_seconds <= 0 or timeout_seconds > 120:
        raise ValueError("timeout must be within 1..120 seconds")
    shipment_id, details, receipt, receipt_bytes = verify_purchase(purchase_directory)
    if sha256(canonical_bytes(purchase_selection)) != receipt["selection_sha256"]:
        raise PermissionError("purchase selection does not match the completed purchase")
    document = purchase_selection.get("document")
    if not isinstance(document, dict) or document.get("format") not in FORMAT_EXTENSIONS or not isinstance(document.get("dpi"), int):
        raise ValueError("purchase selection document format/dpi is invalid")
    package_reference = nonempty(package_reference, "package_client_reference_id", 256)
    matches = [item for item in details if item.get("package_client_reference_id") == package_reference]
    if len(matches) != 1:
        raise ValueError("package reference must identify exactly one purchased package")
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".amazon-documents-stage-{uuid.uuid4().hex}"; stage.mkdir()
    try:
        response = service.get_shipment_documents(shipment_id, package_reference, format=document["format"], dpi=document["dpi"], x_amzn_shipping_business_id=profile["shipping_business_id"], _request_timeout=timeout_seconds)
        payload = json_value(response)
        result = payload.get("payload") if isinstance(payload, dict) else None
        detail = result.get("package_document_detail") if isinstance(result, dict) else None
        if result.get("shipment_id") != shipment_id or not isinstance(detail, dict) or detail.get("package_client_reference_id") != package_reference:
            raise ValueError("document response is not bound to requested shipment/package")
        documents = detail.get("package_documents")
        if not isinstance(documents, list) or not documents or len(documents) > MAX_DOCUMENTS:
            raise ValueError("document response cardinality is invalid")
        allowed_types = set(document.get("requested_document_types", [])); manifest = []; total = 0
        for index, item in enumerate(documents):
            if not isinstance(item, dict) or item.get("type") not in allowed_types or item.get("format") != document["format"] or not isinstance(item.get("contents"), str):
                raise ValueError("provider document type/format/content is invalid")
            try:
                decoded = base64.b64decode(item["contents"], validate=True)
            except (ValueError, binascii.Error) as exc:
                raise ValueError("provider document is not canonical base64") from exc
            total += len(decoded)
            if not decoded or len(decoded) > MAX_DOCUMENT_BYTES or total > MAX_TOTAL_BYTES:
                raise ValueError("decoded provider document exceeds admitted size")
            digest = sha256(decoded); name = f"{index + 1:02d}-{item['type'].lower()}-{digest}.{FORMAT_EXTENSIONS[item['format']]}"
            (stage / name).write_bytes(decoded)
            manifest.append({"name":name, "type":item["type"], "format":item["format"], "bytes":len(decoded), "sha256":digest})
        response_bytes = (json.dumps(payload, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "provider-response.json").write_bytes(response_bytes)
        receipt_out = {"schema":"elite-amazon-spapi-shipping-documents-receipt/v1", "created_at":datetime.now(timezone.utc).isoformat(), "purchase_receipt_sha256":sha256(receipt_bytes), "shipment_id_sha256":sha256(shipment_id.encode("utf-8")), "package_reference_sha256":sha256(package_reference.encode("utf-8")), "provider_response_sha256":sha256(response_bytes), "documents":manifest, "automatic_delivery_completion":False}
        (stage / "DOCUMENTS_RECEIPT.json").write_text(json.dumps(receipt_out, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        stage.replace(output); return receipt_out
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True); raise


def main() -> int:
    parser = argparse.ArgumentParser(); sub = parser.add_subparsers(dest="command", required=True)
    cancel = sub.add_parser("cancel"); cancel.add_argument("--profile", required=True, type=Path); cancel.add_argument("--purchase-directory", required=True, type=Path); cancel.add_argument("--approval", required=True, type=Path); cancel.add_argument("--output", required=True, type=Path); cancel.add_argument("--timeout-seconds", type=float, default=30.0)
    docs = sub.add_parser("documents"); docs.add_argument("--profile", required=True, type=Path); docs.add_argument("--purchase-directory", required=True, type=Path); docs.add_argument("--purchase-selection", required=True, type=Path); docs.add_argument("--package-client-reference-id", required=True); docs.add_argument("--output", required=True, type=Path); docs.add_argument("--timeout-seconds", type=float, default=30.0)
    args = parser.parse_args(); profile = read_approved_profile(args.profile)
    for field in ("shipping_business_id", "reconciliation_owner"):
        nonempty(profile.get(field), field)
    if args.command == "cancel":
        if profile.get("shipment_cancellation_approved") is not True or not nonempty(profile.get("cancellation_approval_owner"), "cancellation_approval_owner"):
            raise PermissionError("shipment cancellation approval and owner are required")
        approval = json.loads(args.approval.read_text(encoding="utf-8")); receipt = cancel_shipment(service_from_environment(profile), args.purchase_directory, approval, profile, args.output, args.timeout_seconds)
        status = "AMAZON_SPAPI_SHIPPING_CANCELLATION_PASS"
    else:
        if profile.get("shipment_documents_approved") is not True or profile.get("postpurchase_data_handling_approved") is not True:
            raise PermissionError("shipment documents and postpurchase data handling approvals are required")
        selection = json.loads(args.purchase_selection.read_text(encoding="utf-8")); receipt = retrieve_documents(service_from_environment(profile), args.purchase_directory, selection, args.package_client_reference_id, profile, args.output, args.timeout_seconds)
        status = "AMAZON_SPAPI_SHIPPING_DOCUMENTS_PASS"
    print(json.dumps({"status":status, "receipt":receipt}, sort_keys=True)); return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `amazon_spapi_shipping_tracking/test_shipping_postpurchase.py`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:postpurchase-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local offline cancel/document contracts against exact official Amazon SDK method kwargs and response shapes"
license: "LicenseRef-Workspace-Owner"
sha256: "9c743d6eb239b36b0031f2d48e5c095649fa397db068487cdcd89f4194c102d5"
variables: []
secrets_allowed: false
```
````python
import base64
import json
from datetime import datetime, timezone
from pathlib import Path
import tempfile
import unittest

from run_shipping_postpurchase import CancellationEffectUnknown, cancel_shipment, canonical_bytes, retrieve_documents, sha256

NOW = datetime(2026, 8, 31, 16, 0, tzinfo=timezone.utc)


def profile(): return {"shipping_business_id":"AmazonShipping_US"}


def purchase(root: Path):
    directory = root / "purchase"; directory.mkdir()
    response = {"payload":{"shipment_id":"shipment-1", "package_document_details":[{"package_client_reference_id":"package-1", "package_documents":[{"type":"LABEL"}], "tracking_id":"tracking-1"}], "promise":{}}}
    response_bytes = (json.dumps(response, sort_keys=True, indent=2) + "\n").encode(); (directory / "provider-response.json").write_bytes(response_bytes)
    selection = {"schema":"elite-amazon-spapi-shipping-purchase-selection/v1", "approved_by":"owner@example.invalid", "approved_at":"2026-08-31T15:00:00+00:00", "rates_receipt_sha256":"b"*64, "rate_id":"rate-1", "maximum_charge":{"value":"12.50","unit":"USD"}, "document":{"format":"PDF", "size":{"width":4,"length":6,"unit":"INCH"}, "dpi":300, "page_layout":"DEFAULT", "need_file_joining":False, "requested_document_types":["LABEL"]}}
    receipt = {"schema":"elite-amazon-spapi-shipping-purchase-receipt/v1", "created_at":"2026-08-31T15:01:00+00:00", "provider":"Amazon Shipping V2", "sdk_version":"1.11.1", "operation":"purchase_shipment", "quote_response_sha256":"a"*64, "rates_receipt_sha256":"b"*64, "selection_sha256":sha256(canonical_bytes(selection)), "rate_id_sha256":"c"*64, "idempotency_key":"00000000-0000-5000-8000-000000000000", "provider_response_sha256":sha256(response_bytes), "shipment_id_sha256":sha256(b"shipment-1"), "shipment_purchased":True, "automatic_delivery_completion":False, "reconciliation_required":False}
    receipt_bytes = (json.dumps(receipt, sort_keys=True, indent=2) + "\n").encode(); (directory / "PURCHASE_RECEIPT.json").write_bytes(receipt_bytes)
    return directory, selection, receipt_bytes


def approval(receipt_bytes): return {"schema":"elite-amazon-spapi-shipping-cancellation-approval/v1", "approved_by":"cancel-owner@example.invalid", "approved_at":"2026-08-31T15:59:00+00:00", "purchase_receipt_sha256":sha256(receipt_bytes), "reason_code":"CUSTOMER_REQUEST"}


class Result:
    def __init__(self, value): self.value = value
    def to_dict(self): return self.value


class FakeService:
    def __init__(self, cancel_error=None): self.cancel_error, self.cancel_calls, self.document_calls = cancel_error, [], []
    def cancel_shipment(self, shipment_id, **kwargs):
        self.cancel_calls.append((shipment_id, kwargs))
        if self.cancel_error: raise self.cancel_error
        return Result({"payload":{}})
    def get_shipment_documents(self, shipment_id, package_reference, **kwargs):
        self.document_calls.append((shipment_id, package_reference, kwargs))
        encoded = base64.b64encode(b"%PDF-1.7\nlabel").decode()
        return Result({"payload":{"shipment_id":shipment_id, "package_document_detail":{"package_client_reference_id":package_reference, "package_documents":[{"type":"LABEL", "format":"PDF", "contents":encoded}], "tracking_id":"tracking-1"}}})


class PostPurchaseTests(unittest.TestCase):
    def test_cancel_uses_only_bound_shipment_and_business(self):
        with tempfile.TemporaryDirectory() as root:
            directory, _, receipt_bytes = purchase(Path(root)); service = FakeService(); output = Path(root) / "cancel"
            receipt = cancel_shipment(service, directory, approval(receipt_bytes), profile(), output, now=NOW)
            self.assertTrue(receipt["shipment_cancelled"]); self.assertEqual(service.cancel_calls[0][0], "shipment-1")
            self.assertEqual(service.cancel_calls[0][1]["x_amzn_shipping_business_id"], "AmazonShipping_US")

    def test_cancel_timeout_is_unknown_and_never_auto_retry(self):
        with tempfile.TemporaryDirectory() as root:
            directory, _, receipt_bytes = purchase(Path(root)); output = Path(root) / "cancel"
            with self.assertRaises(CancellationEffectUnknown): cancel_shipment(FakeService(TimeoutError("ambiguous")), directory, approval(receipt_bytes), profile(), output, now=NOW)
            unknown = json.loads((output / "CANCELLATION_UNKNOWN_EFFECT.json").read_text()); self.assertTrue(unknown["reconciliation_required"]); self.assertFalse(unknown["automatic_retry_authorized"])

    def test_cancel_tamper_and_existing_output_fail_before_provider(self):
        with tempfile.TemporaryDirectory() as root:
            directory, _, receipt_bytes = purchase(Path(root)); service = FakeService(); value = approval(receipt_bytes); value["purchase_receipt_sha256"] = "0"*64
            with self.assertRaises(PermissionError): cancel_shipment(service, directory, value, profile(), Path(root) / "cancel", now=NOW)
            self.assertEqual(service.cancel_calls, [])
            output = Path(root) / "exists"; output.mkdir()
            with self.assertRaises(FileExistsError): cancel_shipment(service, directory, approval(receipt_bytes), profile(), output, now=NOW)
            self.assertEqual(service.cancel_calls, [])

    def test_documents_are_bound_decoded_hash_named_and_receipted(self):
        with tempfile.TemporaryDirectory() as root:
            directory, selection, _ = purchase(Path(root)); service = FakeService(); output = Path(root) / "docs"
            receipt = retrieve_documents(service, directory, selection, "package-1", profile(), output)
            self.assertEqual(receipt["documents"][0]["sha256"], sha256(b"%PDF-1.7\nlabel")); self.assertEqual(len(service.document_calls), 1)
            self.assertEqual(service.document_calls[0][2]["format"], "PDF"); self.assertEqual(service.document_calls[0][2]["dpi"], 300)
            self.assertTrue((output / receipt["documents"][0]["name"]).is_file()); self.assertFalse(receipt["automatic_delivery_completion"])

    def test_documents_reject_selection_package_and_provider_mismatch_atomically(self):
        with tempfile.TemporaryDirectory() as root:
            directory, selection, _ = purchase(Path(root)); service = FakeService(); bad = json.loads(json.dumps(selection)); bad["document"]["dpi"] = 203
            with self.assertRaises(PermissionError): retrieve_documents(service, directory, bad, "package-1", profile(), Path(root) / "bad-selection")
            with self.assertRaises(ValueError): retrieve_documents(service, directory, selection, "other", profile(), Path(root) / "bad-package")
            self.assertEqual(service.document_calls, [])

    def test_documents_invalid_base64_removes_stage(self):
        class Bad(FakeService):
            def get_shipment_documents(self, shipment_id, package_reference, **kwargs):
                return Result({"payload":{"shipment_id":shipment_id, "package_document_detail":{"package_client_reference_id":package_reference, "package_documents":[{"type":"LABEL", "format":"PDF", "contents":"not-base64!"}]}}})
        with tempfile.TemporaryDirectory() as root:
            directory, selection, _ = purchase(Path(root)); output = Path(root) / "docs"
            with self.assertRaises(ValueError): retrieve_documents(Bad(), directory, selection, "package-1", profile(), output)
            self.assertFalse(output.exists()); self.assertEqual([item for item in Path(root).iterdir() if item.name.startswith(".amazon-documents-stage-")], [])


if __name__ == "__main__": unittest.main(verbosity=2)
````

### FILE: `amazon_spapi_shipping_tracking/collection-inventory-request.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:collection-inventory-request:v1"
operation: CREATE
provenance: AUTHORED
source: "local closed request over official Amazon ClientReferenceDetail/GetUnmanifestedShipmentsRequest models"
license: "LicenseRef-Workspace-Owner"
sha256: "b587b2d58711a18794ba18ff49abd4946f8809279a734530b953f67c915be2b8"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-shipping-collection-inventory-request/v1",
  "client_reference_details": [
    {
      "client_reference_type": "IntegratorShipperId",
      "client_reference_id": ""
    }
  ]
}
````

### FILE: `amazon_spapi_shipping_tracking/collection-form-approval.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:collection-form-approval:v1"
operation: CREATE
provenance: AUTHORED
source: "local approval bound to exact Amazon unmanifested inventory carrier and address"
license: "LicenseRef-Workspace-Owner"
sha256: "d330f11b825e02f5e797fd814d93971749322ee9b45060b679116fdad67bfd0e"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-shipping-collection-form-approval/v1",
  "approved_by": "",
  "approved_at": "",
  "inventory_receipt_sha256": "",
  "carrier_id": "",
  "ship_from_address_sha256": ""
}
````

### FILE: `amazon_spapi_shipping_tracking/run_shipping_collection.py`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:collection-runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed wrapper over official Amazon get_unmanifested_shipments/generate_collection_form"
license: "LicenseRef-Workspace-Owner"
sha256: "c17ad986bcf7e78eee7e59d918c78c756d2435c94b8b1799f24d4bed339043c5"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
import base64
import binascii
from datetime import datetime, timezone
import importlib.metadata
import json
from pathlib import Path
import re
import shutil
import sys
from typing import Any, Protocol
import uuid

from spapi.models.shipping_v2.client_reference_detail import ClientReferenceDetail
from spapi.models.shipping_v2.generate_collection_form_request import GenerateCollectionFormRequest
from spapi.models.shipping_v2.get_unmanifested_shipments_request import GetUnmanifestedShipmentsRequest

from run_shipping_postpurchase import canonical_bytes, read_object
from run_shipping_purchase import HEX64, aware_datetime, exact_keys, nonempty
from run_shipping_rates import build_address
from run_shipping_tracking import SDK_VERSION, json_value, read_approved_profile, service_from_environment, sha256

REFERENCE_TYPES = {"IntegratorShipperId", "IntegratorMerchantId"}
REQUEST_KEYS = {"schema", "client_reference_details"}
APPROVAL_KEYS = {"schema", "approved_by", "approved_at", "inventory_receipt_sha256", "carrier_id", "ship_from_address_sha256"}
INVENTORY_RECEIPT_KEYS = {"schema", "created_at", "provider", "sdk_version", "operation", "request_sha256", "provider_response_sha256", "collection_form_generated"}
MAX_FORM_BYTES = 20 * 1024 * 1024


class CollectionService(Protocol):
    def get_unmanifested_shipments(self, body: GetUnmanifestedShipmentsRequest, **kwargs: Any) -> Any: ...
    def generate_collection_form(self, body: GenerateCollectionFormRequest, **kwargs: Any) -> Any: ...


class CollectionFormEffectUnknown(RuntimeError):
    pass


def build_references(request: dict[str, Any]) -> list[ClientReferenceDetail]:
    exact_keys(request, REQUEST_KEYS, "collection inventory request")
    if request["schema"] != "elite-amazon-spapi-shipping-collection-inventory-request/v1":
        raise ValueError("unsupported collection inventory schema")
    values = request["client_reference_details"]
    if not isinstance(values, list) or not values or len(values) > 2:
        raise ValueError("client_reference_details must contain 1..2 entries")
    result = []; seen = set()
    for index, item in enumerate(values):
        exact_keys(item, {"client_reference_type", "client_reference_id"}, f"client_reference_details[{index}]")
        kind = item["client_reference_type"]
        if kind not in REFERENCE_TYPES or kind in seen:
            raise ValueError("client reference type is invalid or duplicated")
        seen.add(kind); result.append(ClientReferenceDetail(client_reference_type=kind, client_reference_id=nonempty(item["client_reference_id"], "client_reference_id", 128)))
    return result


def inventory(service: CollectionService, request: dict[str, Any], profile: dict[str, Any], output_directory: Path, timeout_seconds: float = 30.0) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION:
        raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if timeout_seconds <= 0 or timeout_seconds > 120:
        raise ValueError("timeout must be within 1..120 seconds")
    model = GetUnmanifestedShipmentsRequest(client_reference_details=build_references(request))
    request_bytes = canonical_bytes(request); output = output_directory.resolve()
    if output.exists(): raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True); stage = output.parent / f".amazon-collection-inventory-{uuid.uuid4().hex}"; stage.mkdir()
    try:
        response = service.get_unmanifested_shipments(model, x_amzn_shipping_business_id=profile["shipping_business_id"], _request_timeout=timeout_seconds)
        payload = json_value(response); carriers = payload.get("unmanifested_carrier_information_list") if isinstance(payload, dict) else None
        if not isinstance(carriers, list): raise ValueError("unmanifested shipment response is invalid")
        for carrier in carriers:
            if not isinstance(carrier, dict) or not isinstance(carrier.get("carrier_id"), str) or not isinstance(carrier.get("unmanifested_shipment_location_list"), list):
                raise ValueError("unmanifested carrier information is invalid")
            for location in carrier["unmanifested_shipment_location_list"]:
                if not isinstance(location, dict) or not isinstance(location.get("address"), dict): raise ValueError("unmanifested location is invalid")
        response_bytes = (json.dumps(payload, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8"); (stage / "provider-response.json").write_bytes(response_bytes)
        receipt = {"schema":"elite-amazon-spapi-shipping-collection-inventory-receipt/v1", "created_at":datetime.now(timezone.utc).isoformat(), "provider":"Amazon Shipping V2", "sdk_version":SDK_VERSION, "operation":"get_unmanifested_shipments", "request_sha256":sha256(request_bytes), "provider_response_sha256":sha256(response_bytes), "collection_form_generated":False}
        (stage / "COLLECTION_INVENTORY_RECEIPT.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8"); stage.replace(output); return receipt
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True); raise


def verify_inventory(directory: Path, request: dict[str, Any], approval: dict[str, Any], now: datetime) -> tuple[list[ClientReferenceDetail], dict[str, Any], bytes, dict[str, Any]]:
    exact_keys(approval, APPROVAL_KEYS, "collection form approval")
    if approval["schema"] != "elite-amazon-spapi-shipping-collection-form-approval/v1": raise ValueError("unsupported collection approval schema")
    nonempty(approval["approved_by"], "approved_by")
    if aware_datetime(approval["approved_at"], "approved_at") > now: raise ValueError("approved_at must not be in the future")
    receipt, receipt_bytes = read_object(directory / "COLLECTION_INVENTORY_RECEIPT.json", "collection inventory receipt")
    exact_keys(receipt, INVENTORY_RECEIPT_KEYS, "collection inventory receipt")
    expected = {"schema":"elite-amazon-spapi-shipping-collection-inventory-receipt/v1", "provider":"Amazon Shipping V2", "sdk_version":SDK_VERSION, "operation":"get_unmanifested_shipments", "collection_form_generated":False}
    if any(receipt.get(k) != v for k,v in expected.items()): raise ValueError("collection inventory receipt mismatch")
    if not HEX64.fullmatch(str(approval["inventory_receipt_sha256"])) or sha256(receipt_bytes) != approval["inventory_receipt_sha256"]: raise PermissionError("approval is not bound to inventory receipt")
    if receipt["request_sha256"] != sha256(canonical_bytes(request)): raise ValueError("inventory request does not match receipt")
    response, response_bytes = read_object(directory / "provider-response.json", "collection inventory response")
    if sha256(response_bytes) != receipt["provider_response_sha256"]: raise ValueError("collection inventory response hash mismatch")
    carrier_id = nonempty(approval["carrier_id"], "carrier_id", 64); locations = []
    for carrier in response.get("unmanifested_carrier_information_list", []):
        if isinstance(carrier, dict) and carrier.get("carrier_id") == carrier_id:
            locations.extend(item.get("address") for item in carrier.get("unmanifested_shipment_location_list", []) if isinstance(item, dict))
    matches = [item for item in locations if isinstance(item, dict) and sha256(canonical_bytes(item)) == approval["ship_from_address_sha256"]]
    if len(matches) != 1: raise ValueError("approval must identify exactly one unmanifested carrier/address")
    return build_references(request), matches[0], receipt_bytes, response


def generate(service: CollectionService, request: dict[str, Any], approval: dict[str, Any], inventory_directory: Path, profile: dict[str, Any], output_directory: Path, timeout_seconds: float = 30.0, now: datetime | None = None) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION: raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if timeout_seconds <= 0 or timeout_seconds > 120: raise ValueError("timeout must be within 1..120 seconds")
    current = (now or datetime.now(timezone.utc)).astimezone(timezone.utc)
    references, address, inventory_receipt_bytes, inventory_response = verify_inventory(inventory_directory, request, approval, current)
    model = GenerateCollectionFormRequest(client_reference_details=references, carrier_id=approval["carrier_id"], ship_from_address=build_address(address, "ship_from_address"))
    fingerprint = sha256(canonical_bytes(request) + canonical_bytes(approval) + canonical_bytes(inventory_response)); key = str(uuid.uuid5(uuid.NAMESPACE_URL, "elite-amazon-collection-form:" + fingerprint))
    output = output_directory.resolve()
    if output.exists(): raise FileExistsError("output_directory must not exist; reconcile existing generation attempt")
    output.parent.mkdir(parents=True, exist_ok=True); output.mkdir()
    context = {"inventory_receipt_sha256":sha256(inventory_receipt_bytes), "approval_sha256":sha256(canonical_bytes(approval)), "carrier_id":approval["carrier_id"], "ship_from_address_sha256":approval["ship_from_address_sha256"], "idempotency_key":key}
    attempt = {"schema":"elite-amazon-spapi-shipping-collection-form-attempt/v1", "created_at":current.isoformat(), **context, "provider_call_started":False, "automatic_retry_authorized":False, "reconciliation_required":False}
    attempt_path = output / "COLLECTION_FORM_ATTEMPT.json"; attempt_path.write_text(json.dumps(attempt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
    try:
        attempt["provider_call_started"] = True; attempt["reconciliation_required"] = True; attempt_path.write_text(json.dumps(attempt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        response = service.generate_collection_form(model, x_amzn_idempotency_key=key, x_amzn_shipping_business_id=profile["shipping_business_id"], _request_timeout=timeout_seconds)
        payload = json_value(response); document = payload.get("collections_form_document") if isinstance(payload, dict) else None
        if not isinstance(document, dict) or document.get("document_format") != "PDF" or not isinstance(document.get("base64_encoded_content"), str): raise ValueError("collection form response is invalid")
        try: decoded = base64.b64decode(document["base64_encoded_content"], validate=True)
        except (ValueError, binascii.Error) as exc: raise ValueError("collection form is not canonical Base64") from exc
        if not decoded.startswith(b"%PDF-") or len(decoded) > MAX_FORM_BYTES: raise ValueError("collection form is not an admitted PDF")
        digest = sha256(decoded); filename = f"collection-form-{digest}.pdf"; (output / filename).write_bytes(decoded)
        response_bytes = (json.dumps(payload, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8"); (output / "provider-response.json").write_bytes(response_bytes)
        receipt = {"schema":"elite-amazon-spapi-shipping-collection-form-receipt/v1", "created_at":datetime.now(timezone.utc).isoformat(), **context, "provider_response_sha256":sha256(response_bytes), "document":{"name":filename, "bytes":len(decoded), "sha256":digest, "format":"PDF"}, "collection_form_generated":True, "physical_handover_completed":False, "reconciliation_required":False}
        (output / "COLLECTION_FORM_RECEIPT.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8"); attempt["provider_call_completed"] = True; attempt["reconciliation_required"] = False; attempt_path.write_text(json.dumps(attempt, sort_keys=True, indent=2) + "\n", encoding="utf-8"); return receipt
    except BaseException as exc:
        unknown = {"schema":"elite-amazon-spapi-shipping-collection-form-unknown-effect/v1", "created_at":datetime.now(timezone.utc).isoformat(), **context, "provider_call_started":True, "automatic_retry_authorized":False, "reconciliation_required":True, "failure_class":type(exc).__name__}; (output / "COLLECTION_FORM_UNKNOWN_EFFECT.json").write_text(json.dumps(unknown, sort_keys=True, indent=2) + "\n", encoding="utf-8"); raise CollectionFormEffectUnknown("collection form effect is unknown; reconcile with the same idempotency key") from exc


def main() -> int:
    parser = argparse.ArgumentParser(); sub = parser.add_subparsers(dest="command", required=True)
    inv = sub.add_parser("inventory"); inv.add_argument("--profile", required=True, type=Path); inv.add_argument("--request", required=True, type=Path); inv.add_argument("--output", required=True, type=Path); inv.add_argument("--timeout-seconds", type=float, default=30.0)
    gen = sub.add_parser("generate"); gen.add_argument("--profile", required=True, type=Path); gen.add_argument("--request", required=True, type=Path); gen.add_argument("--inventory-directory", required=True, type=Path); gen.add_argument("--approval", required=True, type=Path); gen.add_argument("--output", required=True, type=Path); gen.add_argument("--timeout-seconds", type=float, default=30.0)
    args = parser.parse_args(); profile = read_approved_profile(args.profile); nonempty(profile.get("shipping_business_id"), "shipping_business_id"); service = service_from_environment(profile); request = json.loads(args.request.read_text(encoding="utf-8"))
    if args.command == "inventory":
        if profile.get("collection_inventory_approved") is not True: raise PermissionError("collection inventory approval is required")
        receipt = inventory(service, request, profile, args.output, args.timeout_seconds); status = "AMAZON_SPAPI_SHIPPING_COLLECTION_INVENTORY_PASS"
    else:
        if profile.get("collection_form_generation_approved") is not True or profile.get("collection_form_data_handling_approved") is not True or not nonempty(profile.get("collection_form_approval_owner"), "collection_form_approval_owner"): raise PermissionError("collection generation/data/owner approvals are required")
        approval = json.loads(args.approval.read_text(encoding="utf-8")); receipt = generate(service, request, approval, args.inventory_directory, profile, args.output, args.timeout_seconds); status = "AMAZON_SPAPI_SHIPPING_COLLECTION_FORM_PASS"
    print(json.dumps({"status":status, "receipt":receipt}, sort_keys=True)); return 0


if __name__ == "__main__": sys.exit(main())
````

### FILE: `amazon_spapi_shipping_tracking/test_shipping_collection.py`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:collection-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local offline inventory/generate contracts against exact official Amazon SDK models and kwargs"
license: "LicenseRef-Workspace-Owner"
sha256: "5669972af6cd3efd7606949b437cb28b4447041e185ff42d43b11c580e3267d3"
variables: []
secrets_allowed: false
```
````python
import base64
import json
from datetime import datetime, timezone
from pathlib import Path
import tempfile
import unittest

from spapi.models.shipping_v2.generate_collection_form_request import GenerateCollectionFormRequest
from spapi.models.shipping_v2.get_unmanifested_shipments_request import GetUnmanifestedShipmentsRequest
from run_shipping_collection import CollectionFormEffectUnknown, canonical_bytes, generate, inventory, sha256

NOW = datetime(2026,8,31,17,0,tzinfo=timezone.utc)

def request(): return {"schema":"elite-amazon-spapi-shipping-collection-inventory-request/v1", "client_reference_details":[{"client_reference_type":"IntegratorShipperId", "client_reference_id":"shipper-1"}]}
def profile(): return {"shipping_business_id":"AmazonShipping_US"}
ADDRESS = {"name":"Warehouse", "address_line1":"1 Main St", "city":"Seattle", "state_or_region":"WA", "country_code":"US", "postal_code":"98101", "email":"ops@example.invalid", "phone_number":"+12025550100"}

class Result:
    def __init__(self, value): self.value=value
    def to_dict(self): return self.value

class Fake:
    def __init__(self, generate_error=None, bad_document=False): self.inventory_calls=[]; self.generate_calls=[]; self.generate_error=generate_error; self.bad_document=bad_document
    def get_unmanifested_shipments(self, body, **kwargs):
        self.inventory_calls.append((body,kwargs)); return Result({"unmanifested_carrier_information_list":[{"carrier_id":"AMZN_US", "carrier_name":"Amazon", "unmanifested_shipment_location_list":[{"address":ADDRESS, "last_manifest_date":"2026-08-30"}]}]})
    def generate_collection_form(self, body, **kwargs):
        self.generate_calls.append((body,kwargs))
        if self.generate_error: raise self.generate_error
        content = "bad!" if self.bad_document else base64.b64encode(b"%PDF-1.7\ncollection").decode()
        return Result({"collections_form_document":{"base64_encoded_content":content, "document_format":"PDF"}})

def build_inventory(root: Path, service=None):
    output=root/"inventory"; inventory(service or Fake(), request(), profile(), output); receipt_bytes=(output/"COLLECTION_INVENTORY_RECEIPT.json").read_bytes()
    approval={"schema":"elite-amazon-spapi-shipping-collection-form-approval/v1", "approved_by":"handover@example.invalid", "approved_at":"2026-08-31T16:59:00+00:00", "inventory_receipt_sha256":sha256(receipt_bytes), "carrier_id":"AMZN_US", "ship_from_address_sha256":sha256(canonical_bytes(ADDRESS))}
    return output,approval

class CollectionTests(unittest.TestCase):
    def test_inventory_uses_official_model_and_is_read_only(self):
        with tempfile.TemporaryDirectory() as root:
            service=Fake(); receipt=inventory(service,request(),profile(),Path(root)/"inventory")
            self.assertIsInstance(service.inventory_calls[0][0],GetUnmanifestedShipmentsRequest); self.assertFalse(receipt["collection_form_generated"])
    def test_generate_uses_selected_offering_idempotency_and_hash_pdf(self):
        with tempfile.TemporaryDirectory() as root:
            directory,approval=build_inventory(Path(root)); service=Fake(); output=Path(root)/"form"; receipt=generate(service,request(),approval,directory,profile(),output,now=NOW)
            self.assertIsInstance(service.generate_calls[0][0],GenerateCollectionFormRequest); self.assertEqual(service.generate_calls[0][0].carrier_id,"AMZN_US")
            self.assertEqual(service.generate_calls[0][1]["x_amzn_idempotency_key"],receipt["idempotency_key"]); self.assertFalse(receipt["physical_handover_completed"]); self.assertTrue((output/receipt["document"]["name"]).is_file())
    def test_generate_timeout_is_durable_unknown_effect(self):
        with tempfile.TemporaryDirectory() as root:
            directory,approval=build_inventory(Path(root)); output=Path(root)/"form"
            with self.assertRaises(CollectionFormEffectUnknown): generate(Fake(TimeoutError("ambiguous")),request(),approval,directory,profile(),output,now=NOW)
            unknown=json.loads((output/"COLLECTION_FORM_UNKNOWN_EFFECT.json").read_text()); self.assertFalse(unknown["automatic_retry_authorized"]); self.assertTrue(unknown["reconciliation_required"])
    def test_tampered_receipt_address_and_request_fail_before_generate(self):
        with tempfile.TemporaryDirectory() as root:
            directory,approval=build_inventory(Path(root)); service=Fake(); bad=dict(approval); bad["inventory_receipt_sha256"]="0"*64
            with self.assertRaises(PermissionError): generate(service,request(),bad,directory,profile(),Path(root)/"one",now=NOW)
            bad=dict(approval); bad["ship_from_address_sha256"]="0"*64
            with self.assertRaises(ValueError): generate(service,request(),bad,directory,profile(),Path(root)/"two",now=NOW)
            altered=request(); altered["client_reference_details"][0]["client_reference_id"]="other"
            with self.assertRaises(ValueError): generate(service,altered,approval,directory,profile(),Path(root)/"three",now=NOW)
            self.assertEqual(service.generate_calls,[])
    def test_invalid_reference_and_existing_output_fail_before_provider(self):
        bad=request(); bad["client_reference_details"][0]["client_reference_type"]="UNKNOWN"
        with tempfile.TemporaryDirectory() as root:
            service=Fake()
            with self.assertRaises(ValueError): inventory(service,bad,profile(),Path(root)/"inventory")
            output=Path(root)/"exists"; output.mkdir()
            with self.assertRaises(FileExistsError): inventory(service,request(),profile(),output)
            self.assertEqual(service.inventory_calls,[])
    def test_invalid_form_payload_becomes_unknown_without_false_success(self):
        with tempfile.TemporaryDirectory() as root:
            directory,approval=build_inventory(Path(root)); output=Path(root)/"form"
            with self.assertRaises(CollectionFormEffectUnknown): generate(Fake(bad_document=True),request(),approval,directory,profile(),output,now=NOW)
            self.assertFalse((output/"COLLECTION_FORM_RECEIPT.json").exists())

if __name__=="__main__": unittest.main(verbosity=2)
````

### FILE: `amazon_spapi_shipping_tracking/claim-approval.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:claim-approval:v1"
operation: CREATE
provenance: AUTHORED
source: "local claim approval bound to exact purchase/quote/package/insured value and proof URLs"
license: "LicenseRef-Workspace-Owner"
sha256: "7dca1a6a4b04bda9714939b3e5a4273e7f262f7eb9cf292516ae76ad9a163260"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-shipping-claim-approval/v1",
  "approved_by": "",
  "approved_at": "",
  "purchase_receipt_sha256": "",
  "rates_receipt_sha256": "",
  "package_client_reference_id": "",
  "declared_value": {
    "value": "0.00",
    "unit": "USD"
  },
  "claim_reason": "DAMAGED_IN_TRANSIT",
  "proofs": [],
  "settlement_type": "REFUND"
}
````

### FILE: `amazon_spapi_shipping_tracking/run_shipping_claim.py`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:claim-runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed wrapper over official Amazon CreateClaimRequest/create_claim"
license: "LicenseRef-Workspace-Owner"
sha256: "0f2211c6084dfab5da0bf1fd13a92031aaef8f7c4e6814512b7d377840e968e3"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
from datetime import datetime, timezone
import importlib.metadata
import json
from pathlib import Path
import re
import sys
from typing import Any, Protocol
from urllib.parse import urlsplit

from spapi.models.shipping_v2.create_claim_request import CreateClaimRequest
from spapi.models.shipping_v2.currency import Currency

from run_shipping_postpurchase import canonical_bytes, read_object, verify_purchase
from run_shipping_purchase import aware_datetime, decimal_value, exact_keys, nonempty
from run_shipping_tracking import SDK_VERSION, json_value, read_approved_profile, service_from_environment, sha256

APPROVAL_KEYS = {"schema", "approved_by", "approved_at", "purchase_receipt_sha256", "rates_receipt_sha256", "package_client_reference_id", "declared_value", "claim_reason", "proofs", "settlement_type"}
CLAIM_REASONS = {"LOST_IN_TRANSIT", "DAMAGED_IN_TRANSIT", "DELIVERED_NOT_RECEIVED", "ITEM_MISSING_SWITCHEROO", "COD_ABUSE"}
SETTLEMENT_TYPES = {"REFUND", "CERTIFICATE_OF_FACT"}


class ClaimService(Protocol):
    def create_claim(self, body: CreateClaimRequest, **kwargs: Any) -> Any: ...


class ClaimEffectUnknown(RuntimeError):
    pass


def validate_proofs(values: Any, approved_hosts: Any) -> list[str]:
    if not isinstance(approved_hosts, list) or not approved_hosts or any(not isinstance(item, str) or not re.fullmatch(r"[a-z0-9.-]+", item) for item in approved_hosts):
        raise PermissionError("approved_claim_proof_hosts must be a non-empty canonical host list")
    if not isinstance(values, list) or not values or len(values) > 10 or len(values) != len(set(values)):
        raise ValueError("proofs must contain 1..10 unique HTTPS URLs")
    for value in values:
        if not isinstance(value, str) or len(value) > 2048: raise ValueError("proof URL is invalid")
        parsed = urlsplit(value)
        if parsed.scheme != "https" or parsed.username or parsed.password or parsed.port not in (None,443) or parsed.hostname not in approved_hosts or parsed.query or parsed.fragment:
            raise PermissionError("proof URL must use an approved exact HTTPS host without credentials/query/fragment")
    return values


def build_claim(purchase_directory: Path, quote_directory: Path, rate_request: dict[str, Any], purchase_selection: dict[str, Any], approval: dict[str, Any], profile: dict[str, Any], now: datetime) -> tuple[CreateClaimRequest, dict[str, Any]]:
    exact_keys(approval, APPROVAL_KEYS, "claim approval")
    if approval["schema"] != "elite-amazon-spapi-shipping-claim-approval/v1" or approval["claim_reason"] not in CLAIM_REASONS or approval["settlement_type"] not in SETTLEMENT_TYPES:
        raise ValueError("claim approval schema/reason/settlement is invalid")
    nonempty(approval["approved_by"], "approved_by")
    if aware_datetime(approval["approved_at"], "approved_at") > now: raise ValueError("approved_at must not be in the future")
    shipment_id, details, purchase_receipt, purchase_receipt_bytes = verify_purchase(purchase_directory)
    if approval["purchase_receipt_sha256"] != sha256(purchase_receipt_bytes): raise PermissionError("claim is not bound to purchase receipt")
    rates_receipt, rates_receipt_bytes = read_object(quote_directory / "RATES_RECEIPT.json", "rates receipt")
    quote_response_bytes = (quote_directory / "provider-response.json").read_bytes()
    if approval["rates_receipt_sha256"] != sha256(rates_receipt_bytes) or purchase_receipt["rates_receipt_sha256"] != sha256(rates_receipt_bytes): raise PermissionError("claim/purchase are not bound to rates receipt")
    if rates_receipt.get("provider_response_sha256") != sha256(quote_response_bytes) or purchase_receipt["quote_response_sha256"] != sha256(quote_response_bytes): raise ValueError("quote response is not bound to purchase")
    if rates_receipt.get("request_sha256") != sha256(canonical_bytes(rate_request)): raise ValueError("rate request is not bound to quote")
    if purchase_receipt["selection_sha256"] != sha256(canonical_bytes(purchase_selection)): raise ValueError("purchase selection is not bound to purchase")
    package = rate_request.get("package") if isinstance(rate_request, dict) else None
    if not isinstance(package, dict): raise ValueError("rate request package is invalid")
    package_reference = nonempty(approval["package_client_reference_id"], "package_client_reference_id", 256)
    if package.get("package_client_reference_id") != package_reference: raise ValueError("claim package does not match rate request")
    matches = [item for item in details if item.get("package_client_reference_id") == package_reference]
    if len(matches) != 1: raise ValueError("claim package must identify exactly one purchased package")
    tracking_id = nonempty(matches[0].get("tracking_id"), "tracking_id", 256)
    insured = exact_keys(package.get("insured_value"), {"value","unit"}, "insured_value")
    declared = exact_keys(approval["declared_value"], {"value","unit"}, "declared_value")
    if declared["unit"] != insured["unit"] or decimal_value(declared["value"], "declared value", positive=True) > decimal_value(insured["value"], "insured value", positive=True): raise PermissionError("claim value exceeds insured value or currency")
    proofs = validate_proofs(approval["proofs"], profile.get("approved_claim_proof_hosts"))
    model = CreateClaimRequest(tracking_id=tracking_id, declared_value=Currency(value=float(decimal_value(declared["value"], "declared value", positive=True)), unit=declared["unit"]), claim_reason=approval["claim_reason"], proofs=proofs, settlement_type=approval["settlement_type"])
    context = {"purchase_receipt_sha256":sha256(purchase_receipt_bytes), "rates_receipt_sha256":sha256(rates_receipt_bytes), "shipment_id_sha256":sha256(shipment_id.encode()), "tracking_id_sha256":sha256(tracking_id.encode()), "package_reference_sha256":sha256(package_reference.encode()), "approval_sha256":sha256(canonical_bytes(approval))}
    return model, context


def create_claim(service: ClaimService, purchase_directory: Path, quote_directory: Path, rate_request: dict[str, Any], purchase_selection: dict[str, Any], approval: dict[str, Any], profile: dict[str, Any], output_directory: Path, timeout_seconds: float = 30.0, now: datetime | None = None) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION: raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if timeout_seconds <= 0 or timeout_seconds > 120: raise ValueError("timeout must be within 1..120 seconds")
    current=(now or datetime.now(timezone.utc)).astimezone(timezone.utc); model,context=build_claim(purchase_directory,quote_directory,rate_request,purchase_selection,approval,profile,current)
    output=output_directory.resolve()
    if output.exists(): raise FileExistsError("output_directory must not exist; reconcile existing claim attempt")
    output.parent.mkdir(parents=True,exist_ok=True); output.mkdir(); attempt={"schema":"elite-amazon-spapi-shipping-claim-attempt/v1","created_at":current.isoformat(),**context,"provider_call_started":False,"automatic_retry_authorized":False,"reconciliation_required":False}; path=output/"CLAIM_ATTEMPT.json"; path.write_text(json.dumps(attempt,sort_keys=True,indent=2)+"\n",encoding="utf-8")
    try:
        attempt["provider_call_started"]=True; attempt["reconciliation_required"]=True; path.write_text(json.dumps(attempt,sort_keys=True,indent=2)+"\n",encoding="utf-8")
        response=service.create_claim(model,x_amzn_shipping_business_id=profile["shipping_business_id"],_request_timeout=timeout_seconds); payload=json_value(response); claim_id=nonempty(payload.get("claim_id") if isinstance(payload,dict) else None,"claim_id",256)
        response_bytes=(json.dumps(payload,ensure_ascii=False,sort_keys=True,indent=2)+"\n").encode(); (output/"provider-response.json").write_bytes(response_bytes)
        receipt={"schema":"elite-amazon-spapi-shipping-claim-receipt/v1","created_at":datetime.now(timezone.utc).isoformat(),**context,"provider_response_sha256":sha256(response_bytes),"claim_id_sha256":sha256(claim_id.encode()),"claim_created":True,"automatic_retry_authorized":False,"reconciliation_required":False}; (output/"CLAIM_RECEIPT.json").write_text(json.dumps(receipt,sort_keys=True,indent=2)+"\n",encoding="utf-8"); attempt["provider_call_completed"]=True; attempt["reconciliation_required"]=False; path.write_text(json.dumps(attempt,sort_keys=True,indent=2)+"\n",encoding="utf-8"); return receipt
    except BaseException as exc:
        unknown={"schema":"elite-amazon-spapi-shipping-claim-unknown-effect/v1","created_at":datetime.now(timezone.utc).isoformat(),**context,"provider_call_started":True,"automatic_retry_authorized":False,"reconciliation_required":True,"failure_class":type(exc).__name__}; (output/"CLAIM_UNKNOWN_EFFECT.json").write_text(json.dumps(unknown,sort_keys=True,indent=2)+"\n",encoding="utf-8"); raise ClaimEffectUnknown("claim effect is unknown; reconcile before retry") from exc


def main()->int:
    parser=argparse.ArgumentParser(); parser.add_argument("--profile",required=True,type=Path); parser.add_argument("--purchase-directory",required=True,type=Path); parser.add_argument("--quote-directory",required=True,type=Path); parser.add_argument("--rate-request",required=True,type=Path); parser.add_argument("--purchase-selection",required=True,type=Path); parser.add_argument("--approval",required=True,type=Path); parser.add_argument("--output",required=True,type=Path); parser.add_argument("--timeout-seconds",type=float,default=30.0); args=parser.parse_args(); profile=read_approved_profile(args.profile)
    if profile.get("claims_approved") is not True or profile.get("claim_data_handling_approved") is not True or not nonempty(profile.get("claim_approval_owner"),"claim_approval_owner"): raise PermissionError("claim approvals/owner are required")
    receipt=create_claim(service_from_environment(profile),args.purchase_directory,args.quote_directory,json.loads(args.rate_request.read_text(encoding="utf-8")),json.loads(args.purchase_selection.read_text(encoding="utf-8")),json.loads(args.approval.read_text(encoding="utf-8")),profile,args.output,args.timeout_seconds); print(json.dumps({"status":"AMAZON_SPAPI_SHIPPING_CLAIM_PASS","receipt":receipt},sort_keys=True)); return 0


if __name__=="__main__": sys.exit(main())
````

### FILE: `amazon_spapi_shipping_tracking/test_shipping_claim.py`
```yaml
block_id: "PY-AMAZON-SPAPI-SHIPPING-TRACKING:claim-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local offline claim contracts against exact official Amazon SDK model and method kwargs"
license: "LicenseRef-Workspace-Owner"
sha256: "0737b5f327363a5fa068b6834298ea122dcbd7148916ed4ede7e7824d2b3b6a8"
variables: []
secrets_allowed: false
```
````python
import json
from datetime import datetime, timezone
from pathlib import Path
import tempfile
import unittest

from spapi.models.shipping_v2.create_claim_request import CreateClaimRequest
from run_shipping_claim import ClaimEffectUnknown, build_claim, canonical_bytes, create_claim, sha256

NOW=datetime(2026,8,31,18,0,tzinfo=timezone.utc)
def profile(): return {"shipping_business_id":"AmazonShipping_US","approved_claim_proof_hosts":["evidence.example.com"]}
def artifacts(root:Path):
    rate={"package":{"insured_value":{"value":1000,"unit":"USD"},"package_client_reference_id":"package-1"}}
    selection={"schema":"elite-amazon-spapi-shipping-purchase-selection/v1","document":{"format":"PDF"}}
    quote=root/"quote"; quote.mkdir(); quote_response=b'{"payload":{}}\n'; (quote/"provider-response.json").write_bytes(quote_response)
    rates={"request_sha256":sha256(canonical_bytes(rate)),"provider_response_sha256":sha256(quote_response)}; rates_bytes=(json.dumps(rates,sort_keys=True,indent=2)+"\n").encode(); (quote/"RATES_RECEIPT.json").write_bytes(rates_bytes)
    purchase=root/"purchase"; purchase.mkdir(); response={"payload":{"shipment_id":"shipment-1","package_document_details":[{"package_client_reference_id":"package-1","package_documents":[],"tracking_id":"tracking-1"}]}}; response_bytes=(json.dumps(response,sort_keys=True,indent=2)+"\n").encode(); (purchase/"provider-response.json").write_bytes(response_bytes)
    receipt={"schema":"elite-amazon-spapi-shipping-purchase-receipt/v1","created_at":"2026-08-31T17:00:00+00:00","provider":"Amazon Shipping V2","sdk_version":"1.11.1","operation":"purchase_shipment","quote_response_sha256":sha256(quote_response),"rates_receipt_sha256":sha256(rates_bytes),"selection_sha256":sha256(canonical_bytes(selection)),"rate_id_sha256":"c"*64,"idempotency_key":"key","provider_response_sha256":sha256(response_bytes),"shipment_id_sha256":sha256(b"shipment-1"),"shipment_purchased":True,"automatic_delivery_completion":False,"reconciliation_required":False}; receipt_bytes=(json.dumps(receipt,sort_keys=True,indent=2)+"\n").encode(); (purchase/"PURCHASE_RECEIPT.json").write_bytes(receipt_bytes)
    approval={"schema":"elite-amazon-spapi-shipping-claim-approval/v1","approved_by":"claims@example.invalid","approved_at":"2026-08-31T17:59:00+00:00","purchase_receipt_sha256":sha256(receipt_bytes),"rates_receipt_sha256":sha256(rates_bytes),"package_client_reference_id":"package-1","declared_value":{"value":"900","unit":"USD"},"claim_reason":"DAMAGED_IN_TRANSIT","proofs":["https://evidence.example.com/proof/1"],"settlement_type":"REFUND"}
    return purchase,quote,rate,selection,approval
class Result:
    def to_dict(self): return {"claim_id":"claim-1"}
class Fake:
    def __init__(self,error=None): self.error=error; self.calls=[]
    def create_claim(self,body,**kwargs): self.calls.append((body,kwargs)); (_ for _ in ()).throw(self.error) if self.error else None; return Result()
class ClaimTests(unittest.TestCase):
    def test_builds_official_claim_from_bound_artifacts(self):
        with tempfile.TemporaryDirectory() as root:
            p,q,r,s,a=artifacts(Path(root)); model,_=build_claim(p,q,r,s,a,profile(),NOW); self.assertIsInstance(model,CreateClaimRequest); self.assertEqual(model.tracking_id,"tracking-1"); self.assertEqual(model.claim_reason,"DAMAGED_IN_TRANSIT"); self.assertIsNone(model.is_replacement_package_sent)
    def test_success_persists_hashed_claim_id(self):
        with tempfile.TemporaryDirectory() as root:
            p,q,r,s,a=artifacts(Path(root)); service=Fake(); receipt=create_claim(service,p,q,r,s,a,profile(),Path(root)/"claim",now=NOW); self.assertTrue(receipt["claim_created"]); self.assertEqual(receipt["claim_id_sha256"],sha256(b"claim-1")); self.assertEqual(service.calls[0][1]["x_amzn_shipping_business_id"],"AmazonShipping_US")
    def test_timeout_is_unknown_and_not_retryable(self):
        with tempfile.TemporaryDirectory() as root:
            p,q,r,s,a=artifacts(Path(root)); output=Path(root)/"claim"
            with self.assertRaises(ClaimEffectUnknown): create_claim(Fake(TimeoutError("ambiguous")),p,q,r,s,a,profile(),output,now=NOW)
            unknown=json.loads((output/"CLAIM_UNKNOWN_EFFECT.json").read_text()); self.assertFalse(unknown["automatic_retry_authorized"]); self.assertTrue(unknown["reconciliation_required"])
    def test_excess_value_bad_host_and_tamper_fail_before_provider(self):
        mutators=[lambda a:a["declared_value"].update(value="1001"),lambda a:a.update(proofs=["https://evil.example/proof"]),lambda a:a.update(purchase_receipt_sha256="0"*64)]
        for mutate in mutators:
            with self.subTest(mutate=mutate),tempfile.TemporaryDirectory() as root:
                p,q,r,s,a=artifacts(Path(root)); mutate(a); service=Fake()
                with self.assertRaises((ValueError,PermissionError)): create_claim(service,p,q,r,s,a,profile(),Path(root)/"claim",now=NOW)
                self.assertEqual(service.calls,[])
    def test_existing_output_blocks_provider(self):
        with tempfile.TemporaryDirectory() as root:
            p,q,r,s,a=artifacts(Path(root)); output=Path(root)/"claim"; output.mkdir(); service=Fake()
            with self.assertRaises(FileExistsError): create_claim(service,p,q,r,s,a,profile(),output,now=NOW)
            self.assertEqual(service.calls,[])
if __name__=="__main__": unittest.main(verbosity=2)
````

## 6. Configuration surface

Perfil: decisión, región/business/carriers/acceso/costo/datos; approvals/owners de ciclo, collection y claims; proof hosts exactos. Claim runtime recibe purchase/quote/rate-request/selection/approval, output nuevo y timeout.

## 7. Dependency bill

`amzn-sp-api` 1.11.1 y seis wheels transitivos por URL/SHA exactos. Código SDK/modelos Amazon Apache-2.0; wrapper y controles locales `AUTHORED`. Ejecutar revisión de licencias efectiva en el proyecto.

## 8. Apply order

Materializar; instalar lock; aprobar perfil; ejecutar treinta y tres pruebas; probar sandbox ciclo+collection+claim; reconciliar `UNKNOWN_EFFECT`; no usar claim/form como prueba física.

## 9. Verification

Exigir 23/23 bloques, hashes, frozen install, 33 pruebas, métodos/modelos oficiales, sandbox real, claim binding/insured-value/proof hosts/UNKNOWN_EFFECT, collection/ciclo previos, privacidad/legal, reconciliación, carga, alarmas y rollback. Form/claim no prueban pickup/entrega.

## 10. Reconstruction evidence

V147 tracking; V148 rates; V149 purchase; V150 cancel/docs; V151 collection; V152 claim 0.6.0. Ningún glue local se atribuye a Amazon.
