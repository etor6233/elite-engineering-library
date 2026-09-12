# Python TikTok Lead Adapter

## 1. Metadata

```yaml
pack_id: "PYTHON-TIKTOK-LEAD-ADAPTER"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa una adaptación mínima sobre el transporte genérico del wheel oficial TikTok Business API 1.1.3 para construir exactamente Lead GET y LEAD subscription v1.3, validar su configuración, emitir tres artefactos hash-linked y analizar señales webhook sin autenticarlas ni persistirlas automáticamente."
stacks: ["CPython 3.12+", "tiktok-business-api-sdk-official 1.1.3", "TikTok Business API v1.3"]
compatible_with: ["durable idempotency owner", "authenticated retrieval reconciliation", "secret manager references"]
incompatible_with: ["webhook tratado como autenticado", "persistencia automática desde webhook", "credenciales embebidas", "SDK Lead atribuido falsamente a TikTok"]
license_expression: "LicenseRef-Workspace-Owner AND MIT AND third-party-wheel-licenses"
upstream_sources: ["https://github.com/tiktok/tiktok-business-api-sdk", "https://pypi.org/project/tiktok-business-api-sdk-official/1.1.3/", "https://business-api.tiktok.com/portal/docs/get-an-instant-form-lead-or-a-direct-message-lead/v1.3", "https://business-api.tiktok.com/portal/docs/create-a-subscription/v1.3"]
verified_at: "2026-09-04"
```

## 2. Applicability

Use para TikTok Instant Form o Direct Message Leads sólo después de demostrar cuenta/app/permisos, términos y consentimiento, cuenta de prueba, callback controlado, delivery del proveedor, deduplicación durable, reconciliación mediante retrieval autenticado, retención/borrado y cuota/costo. TikTok publica los contratos v1.3; el wheel oficial no trae métodos Lead generados. `adapter.py` es `ADAPTED` del patrón de métodos oficiales sobre `ApiClient.call_api`, no código Lead publicado por TikTok.

## 3. Architecture contract

El perfil nace bloqueado y prohíbe escritura de negocio automática. El cliente construye sólo `GET /open_api/v1.3/lead/get/` y `POST /open_api/v1.3/subscription/subscribe/`, conserva `Access-Token` en memoria, exige selectores admitidos y convierte cualquier excepción del transporte en un error sin secretos. El parser acepta únicamente el envelope Lead documentado, limita bytes/cardinalidad, conserva SHA-256 y una clave estable para deduplicación, pero etiqueta la señal `UNAUTHENTICATED_PROVIDER_SIGNAL`. TikTok documenta entrega al menos una vez; el owner de persistencia debe deduplicar duraderamente y reconciliar por la API autenticada.

## 4. Exact file manifest

```text
CREATE tiktok_lead_adapter/requirements.lock
CREATE tiktok_lead_adapter/sdk-artifact.lock.json
CREATE tiktok_lead_adapter/official-contract.lock.json
CREATE tiktok_lead_adapter/provider-profile.template.json
CREATE tiktok_lead_adapter/request.template.json
CREATE tiktok_lead_adapter/adapter.py
CREATE tiktok_lead_adapter/webhook.py
CREATE tiktok_lead_adapter/evidence.py
CREATE tiktok_lead_adapter/run_retrieval.py
CREATE tiktok_lead_adapter/test_tiktok_lead.py
CREATE tiktok_lead_adapter/README.md
```

## 5. Materialization blocks

### FILE: `tiktok_lead_adapter/requirements.lock`
```yaml
block_id: "PY-TIKTOK-LEAD:requirements-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "exact PyPI wheels for official TikTok SDK runtime"
license: "LicenseRef-Workspace-Owner"
sha256: "610bdd481a970f93b33b8d019428d8f30a4c78a4da4e057c3f096c91fe30ba34"
variables: []
secrets_allowed: false
```
````text
certifi @ https://files.pythonhosted.org/packages/0b/a7/71ac2cff56fec219ed242bb11b8efb69fcc4bec75db06fb7bfe35de520e6/certifi-2026.7.22-py3-none-any.whl#sha256=62f22742b58a1a33014a2b6b706588a8d7e2a88ae7bd1a6ebe8c992928483775
python-dateutil @ https://files.pythonhosted.org/packages/ec/57/56b9bcc3c9c6a792fcbaf139543cee77261f3651ca9da0c93f5c1221264b/python_dateutil-2.9.0.post0-py2.py3-none-any.whl#sha256=a8b2bc7bffae282281c8140a97d3aa9c14da0b136dfe83f850eea9a5f7470427
six @ https://files.pythonhosted.org/packages/b7/ce/149a00dd41f10bc29e5921b496af8b574d8413afcd5e30dfa0ed46c2cc5e/six-1.17.0-py2.py3-none-any.whl#sha256=4721f391ed90541fddacab5acf947aa0d3dc7d27b2e1e8eda2be8970586c3274
tiktok-business-api-sdk-official @ https://files.pythonhosted.org/packages/ec/bd/1cf2faaa9fb43a4e81dec0cbe1076c1ffe0d2536611650ff7891a68a932b/tiktok_business_api_sdk_official-1.1.3-py3-none-any.whl#sha256=663b4a1f2585c4b386144d968646f695befd733c464418bbaeeac41befff33b7
urllib3 @ https://files.pythonhosted.org/packages/7f/3e/5db95bcf282c52709639744ca2a8b149baccf648e39c8cc87553df9eae0c/urllib3-2.7.0-py3-none-any.whl#sha256=9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897
````

### FILE: `tiktok_lead_adapter/sdk-artifact.lock.json`
```yaml
block_id: "PY-TIKTOK-LEAD:sdk-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "official TikTok wheel and repository identity receipt"
license: "LicenseRef-Workspace-Owner"
sha256: "f07209b62988ed9a69ff8bf0555d09ccef45679506f3ace3c0dc866ddb37f684"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-python-sdk-artifact-lock/v1",
  "provider": "TikTok Pte. Ltd.",
  "distribution": "tiktok-business-api-sdk-official",
  "version": "1.1.3",
  "module_declared_version": "1.2.1",
  "version_note": "The exact official wheel metadata declares distribution 1.1.3 while business_api_client.__version__ declares 1.2.1; both values are pinned and tested without asserting equivalence.",
  "wheel": {
    "filename": "tiktok_business_api_sdk_official-1.1.3-py3-none-any.whl",
    "bytes": 678182,
    "sha256": "663b4a1f2585c4b386144d968646f695befd733c464418bbaeeac41befff33b7",
    "url": "https://files.pythonhosted.org/packages/ec/bd/1cf2faaa9fb43a4e81dec0cbe1076c1ffe0d2536611650ff7891a68a932b/tiktok_business_api_sdk_official-1.1.3-py3-none-any.whl"
  },
  "publisher_evidence": {
    "pypi_author": "TikTok Pte. Ltd.",
    "repository": "https://github.com/tiktok/tiktok-business-api-sdk",
    "repository_commit": "f809c396520df2d7b201a9ccc5378d822b728ed3",
    "provider_readme_recommends_distribution": true,
    "lead_client_generated_by_provider": false
  },
  "api_version": "v1.3",
  "license_expression": "MIT",
  "automatic_business_write": false
}
````

### FILE: `tiktok_lead_adapter/official-contract.lock.json`
```yaml
block_id: "PY-TIKTOK-LEAD:official-contract:v1"
operation: CREATE
provenance: AUTHORED
source: "official TikTok Business API v1.3 Lead documentation"
license: "LicenseRef-Workspace-Owner"
sha256: "d70d004982d4e3b2fe79737ebe24157ffe935f93a10d13d7d3cb300a93a8312e"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-tiktok-lead-official-contract/v1",
  "provider": "TikTok Business API",
  "observed_at": "2026-09-04",
  "documentation_revision_publicly_exposed": false,
  "lead_get": {
    "documentation": "https://business-api.tiktok.com/portal/docs/get-an-instant-form-lead-or-a-direct-message-lead/v1.3",
    "method": "GET",
    "path": "/open_api/v1.3/lead/get/",
    "header": "Access-Token",
    "lead_sources": ["INSTANT_FORM", "DIRECT_MESSAGE"]
  },
  "subscription_create": {
    "documentation": "https://business-api.tiktok.com/portal/docs/create-a-subscription/v1.3",
    "method": "POST",
    "path": "/open_api/v1.3/subscription/subscribe/",
    "subscribe_entity": "LEAD",
    "response_field": "subscription_id"
  },
  "webhook": {
    "delivery_semantics": "AT_LEAST_ONCE",
    "duplicate_delivery_possible": true,
    "public_signature_mechanism_demonstrated": false,
    "automatic_persistence_allowed": false
  }
}
````

### FILE: `tiktok_lead_adapter/provider-profile.template.json`
```yaml
block_id: "PY-TIKTOK-LEAD:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed TikTok Lead access and operation profile"
license: "LicenseRef-Workspace-Owner"
sha256: "d27333640d1ce2c94fc808fd40ca398a167d2b06421d265e7f8269001dc1ec3e"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-tiktok-lead-profile/v1",
  "decision": "BLOCKED_ACCESS_AND_LIVE_PROOF_REQUIRED",
  "business_account_proven": false,
  "developer_app_proven": false,
  "lead_access_proven": false,
  "terms_and_consent_approved": false,
  "test_account_proven": false,
  "callback_control_and_tls_proven": false,
  "provider_delivery_test_proven": false,
  "durable_idempotency_proven": false,
  "retrieval_reconciliation_proven": false,
  "retention_and_deletion_approved": false,
  "quota_and_cost_approved": false,
  "access_token_environment_variable": "TIKTOK_ACCESS_TOKEN",
  "app_id_environment_variable": "TIKTOK_APP_ID",
  "app_secret_environment_variable": "TIKTOK_APP_SECRET",
  "automatic_business_write": false
}
````

### FILE: `tiktok_lead_adapter/adapter.py`
```yaml
block_id: "PY-TIKTOK-LEAD:adapter:v1"
operation: CREATE
provenance: ADAPTED
source: "TikTok official generated API method pattern plus official Lead v1.3 endpoint contracts"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "e813f4cd6a3fae56bc3e4f2fb979229b465cc8706c7de4d15e68f996283c920a"
variables: []
secrets_allowed: false
```
````python
"""Strict Lead API wrapper over TikTok's official generic ApiClient transport."""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Dict, List, Mapping, Optional, Sequence, Tuple
from urllib.parse import urlparse

from business_api_client.api_client import ApiClient


LEAD_GET_PATH = "/open_api/v1.3/lead/get/"
SUBSCRIBE_PATH = "/open_api/v1.3/subscription/subscribe/"
LEAD_SOURCES = frozenset({"INSTANT_FORM", "DIRECT_MESSAGE"})


class TikTokLeadConfigurationError(ValueError):
    pass


class TikTokLeadProviderError(RuntimeError):
    pass


def _require_nonempty(value: object, name: str, maximum: int = 512) -> str:
    if not isinstance(value, str) or not value.strip() or len(value) > maximum:
        raise TikTokLeadConfigurationError(f"{name} must be a non-empty bounded string")
    return value.strip()


def _require_identifier(value: object, name: str) -> str:
    result = _require_nonempty(value, name, 128)
    if not result.isdigit():
        raise TikTokLeadConfigurationError(f"{name} must contain decimal digits only")
    return result


def validate_profile(profile: Mapping[str, Any]) -> None:
    required_true = (
        "business_account_proven",
        "developer_app_proven",
        "lead_access_proven",
        "terms_and_consent_approved",
        "test_account_proven",
        "callback_control_and_tls_proven",
        "provider_delivery_test_proven",
        "durable_idempotency_proven",
        "retrieval_reconciliation_proven",
        "retention_and_deletion_approved",
        "quota_and_cost_approved",
    )
    if profile.get("decision") != "PROVEN":
        raise PermissionError("TikTok Lead profile decision is not PROVEN")
    missing = [name for name in required_true if profile.get(name) is not True]
    if missing:
        raise PermissionError("TikTok Lead profile lacks required proof")
    if profile.get("automatic_business_write") is not False:
        raise PermissionError("This adapter does not authorize automatic business writes")


def _selectors(
    lead_source: str,
    advertiser_id: Optional[str],
    library_id: Optional[str],
    page_id: Optional[str],
) -> List[Tuple[str, str]]:
    if lead_source not in LEAD_SOURCES:
        raise TikTokLeadConfigurationError("lead_source is not admitted")
    if bool(advertiser_id) == bool(library_id):
        raise TikTokLeadConfigurationError("exactly one of advertiser_id or library_id is required")
    pairs: List[Tuple[str, str]] = [("lead_source", lead_source)]
    if advertiser_id:
        pairs.append(("advertiser_id", _require_identifier(advertiser_id, "advertiser_id")))
    if library_id:
        pairs.append(("library_id", _require_identifier(library_id, "library_id")))
    if lead_source == "INSTANT_FORM":
        pairs.append(("page_id", _require_identifier(page_id, "page_id")))
    elif page_id is not None:
        raise TikTokLeadConfigurationError("page_id is unsupported for DIRECT_MESSAGE")
    return pairs


def _call(api_client: ApiClient, path: str, method: str, query: Sequence[Tuple[str, str]], headers: Dict[str, str], body: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
    try:
        result = api_client.call_api(
            path,
            method,
            {},
            list(query),
            headers,
            body=body,
            post_params=[],
            files={},
            response_type="InlineResponse200",
            auth_settings=[],
            async_req=False,
            _return_http_data_only=True,
            _preload_content=True,
            _request_timeout=(5, 30),
            collection_formats={},
        )
    except Exception:
        raise TikTokLeadProviderError("TikTok Lead API request failed") from None
    if not isinstance(result, dict) or not isinstance(result.get("data"), dict):
        raise TikTokLeadProviderError("TikTok Lead API returned an invalid response")
    request_id = result.get("request_id")
    if not isinstance(request_id, str) or not request_id:
        raise TikTokLeadProviderError("TikTok Lead API response lacks request_id")
    return result


@dataclass(frozen=True)
class TikTokLeadClient:
    api_client: ApiClient
    profile: Mapping[str, Any]

    def __post_init__(self) -> None:
        validate_profile(self.profile)

    def get_lead(self, access_token: str, lead_source: str, *, advertiser_id: Optional[str] = None, library_id: Optional[str] = None, page_id: Optional[str] = None) -> Dict[str, Any]:
        token = _require_nonempty(access_token, "access_token", 4096)
        query = _selectors(lead_source, advertiser_id, library_id, page_id)
        result = _call(self.api_client, LEAD_GET_PATH, "GET", query, {"Accept": "application/json", "Access-Token": token})
        data = result["data"]
        if not isinstance(data.get("lead_data"), dict) or not isinstance(data.get("meta_data"), dict):
            raise TikTokLeadProviderError("TikTok Lead response lacks lead_data or meta_data")
        return result

    def create_subscription(self, access_token: str, app_id: str, app_secret: str, callback_url: str, lead_source: str, *, advertiser_id: Optional[str] = None, library_id: Optional[str] = None, page_id: Optional[str] = None) -> Dict[str, Any]:
        token = _require_nonempty(access_token, "access_token", 4096)
        app = _require_nonempty(app_id, "app_id", 256)
        secret = _require_nonempty(app_secret, "app_secret", 4096)
        callback = _require_nonempty(callback_url, "callback_url", 2048)
        parsed = urlparse(callback)
        if parsed.scheme != "https" or not parsed.hostname or parsed.username or parsed.password or parsed.fragment:
            raise TikTokLeadConfigurationError("callback_url must be a credential-free HTTPS URL")
        selectors = dict(_selectors(lead_source, advertiser_id, library_id, page_id))
        detail: Dict[str, Any] = {"access_token": token, "lead_source": lead_source}
        detail.update({key: value for key, value in selectors.items() if key != "lead_source"})
        body = {
            "app_id": app,
            "secret": secret,
            "subscribe_entity": "LEAD",
            "callback_url": callback,
            "subscription_detail": detail,
        }
        result = _call(self.api_client, SUBSCRIBE_PATH, "POST", (), {"Accept": "application/json", "Content-Type": "application/json", "Access-Token": token}, body)
        subscription_id = result["data"].get("subscription_id")
        if not isinstance(subscription_id, str) or not subscription_id:
            raise TikTokLeadProviderError("TikTok subscription response lacks subscription_id")
        return result
````

### FILE: `tiktok_lead_adapter/request.template.json`
```yaml
block_id: "PY-TIKTOK-LEAD:request:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed request bound to official TikTok Lead selectors"
license: "LicenseRef-Workspace-Owner"
sha256: "004084436c6138fb3297c8b81216bf7cbc8531475b869520c4a2eea3e9af6950"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-tiktok-lead-retrieval-request/v1",
  "tenant_id": "00000000-0000-4000-8000-000000000000",
  "organization_id": "replace-with-approved-organization",
  "mode": "TEST",
  "lead_source": "INSTANT_FORM",
  "advertiser_id": "",
  "library_id": "",
  "page_id": ""
}
````

### FILE: `tiktok_lead_adapter/webhook.py`
```yaml
block_id: "PY-TIKTOK-LEAD:webhook:v1"
operation: CREATE
provenance: AUTHORED
source: "local strict parser governed by TikTok official Lead webhook field table and at-least-once semantics"
license: "LicenseRef-Workspace-Owner"
sha256: "3084633d45626a79d14d486a042fe835e83ff39c11e142f38c0abffcd574f1bb"
variables: []
secrets_allowed: false
```
````python
"""Parser for TikTok Lead webhook signals; it never authenticates or persists them."""

from __future__ import annotations

import hashlib
import json
from dataclasses import dataclass
from typing import Any, Dict, List, Mapping, Optional, Tuple


MAX_WEBHOOK_BYTES = 1_048_576
MAX_ENTRIES = 100
MAX_CHANGES = 200


class TikTokWebhookError(ValueError):
    pass


def _bounded_string(value: object, name: str, maximum: int = 1024) -> str:
    if not isinstance(value, (str, int)):
        raise TikTokWebhookError(f"{name} must be a string-compatible identifier")
    result = str(value).strip()
    if not result or len(result) > maximum:
        raise TikTokWebhookError(f"{name} is empty or too long")
    return result


@dataclass(frozen=True)
class LeadWebhookEntry:
    lead_id: str
    page_id: Optional[str]
    create_time: int
    campaign_id: Optional[str]
    adgroup_id: Optional[str]
    ad_id: Optional[str]
    changes: Tuple[Tuple[str, Any], ...]


@dataclass(frozen=True)
class LeadWebhookSignal:
    request_id: Optional[str]
    notification_time: int
    raw_sha256: str
    delivery_key: str
    trust_state: str
    entries: Tuple[LeadWebhookEntry, ...]

    def redacted_receipt(self) -> Dict[str, Any]:
        return {
            "schema": "elite-tiktok-lead-webhook-receipt/v1",
            "trust_state": self.trust_state,
            "request_id_present": self.request_id is not None,
            "notification_time": self.notification_time,
            "raw_sha256": self.raw_sha256,
            "delivery_key": self.delivery_key,
            "entry_count": len(self.entries),
            "automatic_persistence_allowed": False,
        }


def _optional_identifier(entry: Mapping[str, Any], name: str) -> Optional[str]:
    value = entry.get(name)
    return None if value is None else _bounded_string(value, name, 256)


def parse_lead_webhook(raw_body: bytes) -> LeadWebhookSignal:
    if not isinstance(raw_body, bytes) or not raw_body or len(raw_body) > MAX_WEBHOOK_BYTES:
        raise TikTokWebhookError("webhook body is empty, non-bytes or too large")
    raw_sha = hashlib.sha256(raw_body).hexdigest()
    try:
        payload = json.loads(raw_body.decode("utf-8"))
    except (UnicodeDecodeError, json.JSONDecodeError):
        raise TikTokWebhookError("webhook body is not valid UTF-8 JSON") from None
    if not isinstance(payload, dict) or payload.get("object") != 1:
        raise TikTokWebhookError("webhook object is not TikTok Lead")
    notification_time = payload.get("time")
    if not isinstance(notification_time, int) or notification_time < 0:
        raise TikTokWebhookError("webhook time must be a non-negative integer")
    raw_entries = payload.get("entry")
    if not isinstance(raw_entries, list) or not raw_entries or len(raw_entries) > MAX_ENTRIES:
        raise TikTokWebhookError("webhook entry list is empty or exceeds the safety bound")
    entries: List[LeadWebhookEntry] = []
    for raw_entry in raw_entries:
        if not isinstance(raw_entry, dict):
            raise TikTokWebhookError("webhook entry must be an object")
        lead_id = _bounded_string(raw_entry.get("id"), "id", 256)
        create_time = raw_entry.get("create_time")
        if not isinstance(create_time, int) or create_time < 0:
            raise TikTokWebhookError("entry create_time must be a non-negative integer")
        raw_changes = raw_entry.get("changes")
        if not isinstance(raw_changes, list) or not raw_changes or len(raw_changes) > MAX_CHANGES:
            raise TikTokWebhookError("entry changes are empty or exceed the safety bound")
        changes: List[Tuple[str, Any]] = []
        seen = set()
        for change in raw_changes:
            if not isinstance(change, dict):
                raise TikTokWebhookError("change must be an object")
            field = _bounded_string(change.get("field"), "field", 256)
            if field in seen:
                raise TikTokWebhookError("duplicate change field")
            seen.add(field)
            changes.append((field, change.get("value")))
        entries.append(LeadWebhookEntry(lead_id, _optional_identifier(raw_entry, "page_id"), create_time, _optional_identifier(raw_entry, "campaign_id"), _optional_identifier(raw_entry, "adgroup_id"), _optional_identifier(raw_entry, "ad_id"), tuple(changes)))
    request_value = payload.get("request_id")
    request_id = None if request_value is None else _bounded_string(request_value, "request_id", 256)
    stable = json.dumps({"request_id": request_id, "time": notification_time, "entry_ids": [[entry.lead_id, entry.page_id, entry.create_time] for entry in entries]}, ensure_ascii=False, separators=(",", ":"), sort_keys=True).encode("utf-8")
    delivery_key = hashlib.sha256(stable).hexdigest()
    return LeadWebhookSignal(request_id, notification_time, raw_sha, delivery_key, "UNAUTHENTICATED_PROVIDER_SIGNAL", tuple(entries))
````

### FILE: `tiktok_lead_adapter/evidence.py`
```yaml
block_id: "PY-TIKTOK-LEAD:evidence:v1"
operation: CREATE
provenance: AUTHORED
source: "local hash-linked evidence and conservative dynamic-field normalization"
license: "LicenseRef-Workspace-Owner"
sha256: "5acfcbb4c421926d04eb11c652c5a6d0952c076afad03cb89711c868c68e4f11"
variables: []
secrets_allowed: false
```
````python
"""Hash-linked TikTok Lead retrieval artifacts without automatic persistence."""

from __future__ import annotations

import hashlib
import json
import os
import re
import shutil
import uuid
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Dict, Mapping, Tuple


SDK_VERSION = "1.1.3"
SDK_WHEEL_SHA256 = "663b4a1f2585c4b386144d968646f695befd733c464418bbaeeac41befff33b7"
API_VERSION = "v1.3"
CONTRACT_OBSERVED_AT = "2026-09-04"
PROVIDER = "tiktok_lead_generation"
SAFE_ID = re.compile(r"^[A-Za-z0-9_.:@/+ -]{1,256}$")
SAFE_FIELD = re.compile(r"^[^\x00-\x1f]{1,256}$")


class TikTokLeadEvidenceError(ValueError):
    pass


def _canonical(value: Mapping[str, Any]) -> bytes:
    return (json.dumps(value, ensure_ascii=False, separators=(",", ":"), sort_keys=True) + "\n").encode("utf-8")


def _sha(value: bytes) -> str:
    return hashlib.sha256(value).hexdigest()


def _identifier(value: object, name: str) -> str:
    if not isinstance(value, (str, int)):
        raise TikTokLeadEvidenceError(f"{name} is not a string-compatible identifier")
    result = str(value).strip()
    if not SAFE_ID.fullmatch(result):
        raise TikTokLeadEvidenceError(f"{name} is empty, too long or contains unsafe characters")
    return result


def validate_request(request: Mapping[str, Any]) -> Dict[str, Any]:
    if request.get("schema") != "elite-tiktok-lead-retrieval-request/v1":
        raise TikTokLeadEvidenceError("unsupported retrieval request schema")
    try:
        tenant = str(uuid.UUID(str(request.get("tenant_id", ""))))
    except ValueError:
        raise TikTokLeadEvidenceError("tenant_id must be a canonical UUID") from None
    organization = _identifier(request.get("organization_id"), "organization_id")
    mode = request.get("mode")
    if mode not in ("TEST", "LIVE"):
        raise TikTokLeadEvidenceError("mode must be TEST or LIVE")
    lead_source = request.get("lead_source")
    if lead_source not in ("INSTANT_FORM", "DIRECT_MESSAGE"):
        raise TikTokLeadEvidenceError("lead_source is not admitted")
    return {
        "schema": request["schema"],
        "tenant_id": tenant,
        "organization_id": organization,
        "mode": mode,
        "lead_source": lead_source,
        "advertiser_id": request.get("advertiser_id") or None,
        "library_id": request.get("library_id") or None,
        "page_id": request.get("page_id") or None,
    }


def build_retrieval_artifacts(result: Mapping[str, Any], request_input: Mapping[str, Any]) -> Dict[str, bytes]:
    request = validate_request(request_input)
    if not isinstance(result, dict) or not isinstance(result.get("data"), dict):
        raise TikTokLeadEvidenceError("invalid official SDK response")
    data = result["data"]
    lead_data = data.get("lead_data")
    meta = data.get("meta_data")
    request_id = result.get("request_id")
    if not isinstance(lead_data, dict) or not isinstance(meta, dict) or not isinstance(request_id, str) or not request_id:
        raise TikTokLeadEvidenceError("response lacks lead_data, meta_data or request_id")
    lead_id = _identifier(meta.get("lead_id"), "meta_data.lead_id")
    lead_source = meta.get("lead_source")
    if lead_source != request["lead_source"]:
        raise TikTokLeadEvidenceError("response lead_source does not match request")
    page_id = meta.get("page_id")
    if lead_source == "INSTANT_FORM":
        if _identifier(page_id, "meta_data.page_id") != _identifier(request["page_id"], "request.page_id"):
            raise TikTokLeadEvidenceError("response page_id does not match request")
    elif page_id not in (None, ""):
        raise TikTokLeadEvidenceError("DIRECT_MESSAGE response unexpectedly contains page_id")
    created = meta.get("create_time")
    if not isinstance(created, str):
        raise TikTokLeadEvidenceError("meta_data.create_time is required")
    try:
        submitted = datetime.strptime(created, "%Y-%m-%d %H:%M:%S").replace(tzinfo=timezone.utc)
    except ValueError:
        raise TikTokLeadEvidenceError("meta_data.create_time is not official UTC format") from None

    normalized_fields = []
    rejected = []
    for key in sorted(lead_data):
        if not isinstance(key, str) or not SAFE_FIELD.fullmatch(key):
            rejected.append("UNSAFE_DYNAMIC_FIELD_NAME")
            continue
        value = lead_data[key]
        if not isinstance(value, str) or not value.strip():
            rejected.append("NON_STRING_DYNAMIC_FIELD_REQUIRES_MAPPING")
            continue
        normalized_fields.append({"id": key, "value": value})
    if not lead_data:
        rejected.append("EMPTY_LEAD_DATA")

    provider_response = {
        "schema": "elite-tiktok-lead-sdk-response/v1",
        "sdk_version": SDK_VERSION,
        "api_version": API_VERSION,
        "request_identity": {
            "mode": request["mode"],
            "lead_source": request["lead_source"],
            "advertiser_id_sha256": "" if request["advertiser_id"] is None else _sha(str(request["advertiser_id"]).encode("utf-8")),
            "library_id_sha256": "" if request["library_id"] is None else _sha(str(request["library_id"]).encode("utf-8")),
            "page_id_sha256": "" if request["page_id"] is None else _sha(str(request["page_id"]).encode("utf-8")),
        },
        "data": data,
        "request_id": request_id,
    }
    provider_bytes = _canonical(provider_response)
    status = "REJECTED" if rejected else "CANDIDATE"
    candidate = None
    if status == "CANDIDATE":
        candidate = {
            "tenant_id": request["tenant_id"],
            "organization_id": request["organization_id"],
            "provider": PROVIDER,
            "provider_lead_id": lead_id,
            "form_id": "" if page_id is None else str(page_id),
            "campaign_id": "" if meta.get("campaign_id") is None else _identifier(meta.get("campaign_id"), "campaign_id"),
            "ad_group_id": "" if meta.get("adgroup_id") is None else _identifier(meta.get("adgroup_id"), "adgroup_id"),
            "creative_id": "" if meta.get("ad_id") is None else _identifier(meta.get("ad_id"), "ad_id"),
            "source_kind": lead_source,
            "submitted_at": submitted.isoformat().replace("+00:00", "Z"),
            "is_test": request["mode"] == "TEST",
            "fields": normalized_fields,
            "contact_eligibility": "pending_policy",
        }
    outcome = {
        "status": status,
        "normalization_codes": sorted(set(rejected)),
        "provider_lead_id_sha256": _sha(lead_id.encode("utf-8")),
        "candidate": candidate,
    }
    batch = {
        "schema": "elite-tiktok-lead-candidate-batch/v1",
        "provider": PROVIDER,
        "sdk_version": SDK_VERSION,
        "sdk_wheel_sha256": SDK_WHEEL_SHA256,
        "api_version": API_VERSION,
        "contract_observed_at": CONTRACT_OBSERVED_AT,
        "tenant_id": request["tenant_id"],
        "organization_id": request["organization_id"],
        "provider_response_sha256": _sha(provider_bytes),
        "outcomes": [outcome],
        "automatic_business_write": False,
        "automatic_contact_eligibility": False,
    }
    batch_bytes = _canonical(batch)
    receipt = {
        "schema": "elite-tiktok-lead-retrieval-receipt/v1",
        "sdk_version": SDK_VERSION,
        "sdk_wheel_sha256": SDK_WHEEL_SHA256,
        "api_version": API_VERSION,
        "contract_observed_at": CONTRACT_OBSERVED_AT,
        "retrieval_mode": request["mode"],
        "request_id_sha256": _sha(request_id.encode("utf-8")),
        "provider_response_sha256": _sha(provider_bytes),
        "candidate_batch_sha256": _sha(batch_bytes),
        "row_count": 1,
        "candidate_count": 1 if status == "CANDIDATE" else 0,
        "rejected_count": 1 if status == "REJECTED" else 0,
        "automatic_business_write": False,
        "automatic_contact_eligibility": False,
    }
    return {
        "provider-response.json": provider_bytes,
        "candidate-batch.json": batch_bytes,
        "retrieval-receipt.json": _canonical(receipt),
    }


def write_artifacts_atomically(output: Path, artifacts: Mapping[str, bytes]) -> None:
    output = output.resolve()
    if output.exists():
        raise TikTokLeadEvidenceError("output path must not exist")
    staging = output.parent / ("." + output.name + ".tmp-" + uuid.uuid4().hex)
    if staging.exists():
        raise TikTokLeadEvidenceError("staging path collision")
    staging.mkdir(parents=False)
    try:
        for name in ("provider-response.json", "candidate-batch.json", "retrieval-receipt.json"):
            value = artifacts.get(name)
            if not isinstance(value, bytes) or not value:
                raise TikTokLeadEvidenceError(f"artifact is missing: {name}")
            path = staging / name
            with path.open("xb") as handle:
                handle.write(value)
                handle.flush()
                os.fsync(handle.fileno())
        staging.replace(output)
    except Exception:
        shutil.rmtree(staging, ignore_errors=True)
        raise
````

### FILE: `tiktok_lead_adapter/run_retrieval.py`
```yaml
block_id: "PY-TIKTOK-LEAD:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local approved retrieval runner over the official TikTok transport"
license: "LicenseRef-Workspace-Owner"
sha256: "2ca7547891c17fd6986a1fd733af607ebc1270f54bee52bc7f80c08a7553d656"
variables: []
secrets_allowed: false
```
````python
"""Execute one approved TikTok Lead retrieval and emit three atomic artifacts."""

from __future__ import annotations

import argparse
import json
import os
from pathlib import Path

from business_api_client.api_client import ApiClient

from adapter import TikTokLeadClient
from evidence import build_retrieval_artifacts, validate_request, write_artifacts_atomically


def _load(path: Path):
    return json.loads(path.read_text(encoding="utf-8"))


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--request", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()
    profile = _load(args.profile)
    request = validate_request(_load(args.request))
    token_name = profile.get("access_token_environment_variable")
    if not isinstance(token_name, str) or not token_name:
        raise RuntimeError("profile lacks access token environment variable name")
    token = os.environ.get(token_name, "")
    client = TikTokLeadClient(ApiClient(), profile)
    result = client.get_lead(
        token,
        request["lead_source"],
        advertiser_id=request["advertiser_id"],
        library_id=request["library_id"],
        page_id=request["page_id"],
    )
    artifacts = build_retrieval_artifacts(result, request)
    write_artifacts_atomically(args.output, artifacts)
    receipt = json.loads(artifacts["retrieval-receipt.json"])
    print(json.dumps(receipt, separators=(",", ":"), sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `tiktok_lead_adapter/test_tiktok_lead.py`
```yaml
block_id: "PY-TIKTOK-LEAD:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract, safety and official SDK identity regression tests"
license: "LicenseRef-Workspace-Owner"
sha256: "af9b48fa589f0bf5155fe83e138bd86cc44f44caafb7fbc80119be464b0542d9"
variables: []
secrets_allowed: false
```
````python
import inspect
import importlib.metadata
import json
import tempfile
import unittest
from pathlib import Path

import business_api_client
from business_api_client.api_client import ApiClient
from business_api_client.configuration import Configuration

from adapter import (
    LEAD_GET_PATH,
    SUBSCRIBE_PATH,
    TikTokLeadClient,
    TikTokLeadConfigurationError,
    TikTokLeadProviderError,
    validate_profile,
)
from webhook import TikTokWebhookError, parse_lead_webhook
from evidence import TikTokLeadEvidenceError, build_retrieval_artifacts, write_artifacts_atomically


class FakeApiClient:
    def __init__(self, result=None, error=None):
        self.result = result
        self.error = error
        self.calls = []

    def call_api(self, *args, **kwargs):
        self.calls.append((args, kwargs))
        if self.error:
            raise self.error
        return self.result


def proven_profile():
    return {
        "decision": "PROVEN",
        "business_account_proven": True,
        "developer_app_proven": True,
        "lead_access_proven": True,
        "terms_and_consent_approved": True,
        "test_account_proven": True,
        "callback_control_and_tls_proven": True,
        "provider_delivery_test_proven": True,
        "durable_idempotency_proven": True,
        "retrieval_reconciliation_proven": True,
        "retention_and_deletion_approved": True,
        "quota_and_cost_approved": True,
        "automatic_business_write": False,
    }


class TikTokLeadAdapterTests(unittest.TestCase):
    def test_profile_fails_closed(self):
        profile = proven_profile()
        profile["provider_delivery_test_proven"] = False
        with self.assertRaises(PermissionError):
            validate_profile(profile)
        validate_profile(proven_profile())

    def test_client_construction_enforces_profile(self):
        profile = proven_profile()
        profile["provider_delivery_test_proven"] = False
        with self.assertRaises(PermissionError):
            TikTokLeadClient(FakeApiClient(), profile)
        self.assertIsInstance(TikTokLeadClient(FakeApiClient(), proven_profile()), TikTokLeadClient)

    def test_instant_form_get_uses_exact_official_contract(self):
        fake = FakeApiClient({"data": {"lead_data": {"Email": "a@example.com"}, "meta_data": {"lead_source": "INSTANT_FORM"}}, "request_id": "req-1"})
        result = TikTokLeadClient(fake, proven_profile()).get_lead("token", "INSTANT_FORM", advertiser_id="123", page_id="456")
        self.assertEqual(result["request_id"], "req-1")
        args, kwargs = fake.calls[0]
        self.assertEqual(args[:2], (LEAD_GET_PATH, "GET"))
        self.assertEqual(args[3], [("lead_source", "INSTANT_FORM"), ("advertiser_id", "123"), ("page_id", "456")])
        self.assertEqual(args[4]["Access-Token"], "token")
        self.assertEqual(kwargs["response_type"], "InlineResponse200")

    def test_direct_message_rejects_page_and_accepts_library(self):
        client = TikTokLeadClient(FakeApiClient({"data": {"lead_data": {}, "meta_data": {}}, "request_id": "req-2"}), proven_profile())
        with self.assertRaises(TikTokLeadConfigurationError):
            client.get_lead("token", "DIRECT_MESSAGE", advertiser_id="123", page_id="456")
        client.get_lead("token", "DIRECT_MESSAGE", library_id="789")

    def test_selector_contract_rejects_ambiguous_owner(self):
        client = TikTokLeadClient(FakeApiClient(), proven_profile())
        with self.assertRaises(TikTokLeadConfigurationError):
            client.get_lead("token", "INSTANT_FORM", advertiser_id="1", library_id="2", page_id="3")

    def test_subscription_uses_exact_official_contract(self):
        fake = FakeApiClient({"data": {"subscription_id": "sub-1"}, "request_id": "req-3"})
        result = TikTokLeadClient(fake, proven_profile()).create_subscription("token", "app", "secret", "https://leads.example.test/tiktok", "INSTANT_FORM", advertiser_id="123", page_id="456")
        self.assertEqual(result["data"]["subscription_id"], "sub-1")
        args, kwargs = fake.calls[0]
        self.assertEqual(args[:2], (SUBSCRIBE_PATH, "POST"))
        self.assertEqual(args[4]["Access-Token"], "token")
        self.assertEqual(kwargs["body"]["subscribe_entity"], "LEAD")
        self.assertEqual(kwargs["body"]["subscription_detail"]["lead_source"], "INSTANT_FORM")

    def test_provider_error_does_not_expose_secrets(self):
        fake = FakeApiClient(error=RuntimeError("secret token leaked by fake"))
        with self.assertRaises(TikTokLeadProviderError) as raised:
            TikTokLeadClient(fake, proven_profile()).create_subscription("token", "app", "secret", "https://example.test/callback", "DIRECT_MESSAGE", advertiser_id="123")
        self.assertEqual(str(raised.exception), "TikTok Lead API request failed")
        self.assertNotIn("secret", str(raised.exception))
        self.assertNotIn("token", str(raised.exception))

    def test_webhook_is_parsed_but_remains_unauthenticated(self):
        raw = json.dumps({"request_id": "r-1", "object": 1, "time": 1770000000, "entry": [{"id": "lead-1", "page_id": "page-1", "campaign_id": "campaign-1", "create_time": 1769999999, "changes": [{"field": "email", "value": "person@example.com"}, {"field": "name", "value": "Person"}]}]}, separators=(",", ":")).encode()
        signal = parse_lead_webhook(raw)
        receipt = signal.redacted_receipt()
        self.assertEqual(signal.trust_state, "UNAUTHENTICATED_PROVIDER_SIGNAL")
        self.assertFalse(receipt["automatic_persistence_allowed"])
        self.assertNotIn("person@example.com", json.dumps(receipt))
        self.assertEqual(parse_lead_webhook(raw).delivery_key, signal.delivery_key)

    def test_divergent_payload_has_distinct_raw_hash(self):
        one = b'{"object":1,"time":1,"entry":[{"id":"1","create_time":1,"changes":[{"field":"email","value":"a"}]}]}'
        two = b'{"object":1,"time":1,"entry":[{"id":"1","create_time":1,"changes":[{"field":"email","value":"b"}]}]}'
        first = parse_lead_webhook(one)
        second = parse_lead_webhook(two)
        self.assertEqual(first.delivery_key, second.delivery_key)
        self.assertNotEqual(first.raw_sha256, second.raw_sha256)

    def test_webhook_rejects_duplicate_fields_and_wrong_object(self):
        duplicate = b'{"object":1,"time":1,"entry":[{"id":"1","create_time":1,"changes":[{"field":"email","value":"a"},{"field":"email","value":"b"}]}]}'
        with self.assertRaises(TikTokWebhookError):
            parse_lead_webhook(duplicate)
        with self.assertRaises(TikTokWebhookError):
            parse_lead_webhook(b'{"object":2,"time":1,"entry":[]}')

    def test_official_sdk_transport_identity(self):
        self.assertEqual(importlib.metadata.version("tiktok-business-api-sdk-official"), "1.1.3")
        self.assertEqual(business_api_client.__version__, "1.2.1")
        self.assertEqual(Configuration().host, "https://business-api.tiktok.com")
        self.assertTrue(callable(ApiClient.call_api))
        source = inspect.getsource(ApiClient.call_api)
        self.assertIn("TikTokSDKResponse", source)
        self.assertIn("TiktokSDKError", source)

    def test_retrieval_artifacts_are_hash_linked_and_lossless_for_strings(self):
        request = {"schema": "elite-tiktok-lead-retrieval-request/v1", "tenant_id": "11111111-1111-4111-8111-111111111111", "organization_id": "store-1", "mode": "TEST", "lead_source": "INSTANT_FORM", "advertiser_id": "123", "library_id": "", "page_id": "456"}
        response = {"request_id": "req-1", "data": {"lead_data": {"Email": "person@example.test", "Name": "Person"}, "meta_data": {"lead_source": "INSTANT_FORM", "lead_id": "lead-1", "page_id": "456", "campaign_id": "11", "adgroup_id": "22", "ad_id": "33", "create_time": "2026-09-04 12:00:00"}}}
        artifacts = build_retrieval_artifacts(response, request)
        provider = json.loads(artifacts["provider-response.json"])
        batch = json.loads(artifacts["candidate-batch.json"])
        receipt = json.loads(artifacts["retrieval-receipt.json"])
        self.assertEqual(provider["data"]["lead_data"]["Email"], "person@example.test")
        self.assertEqual(provider["request_identity"]["mode"], "TEST")
        self.assertNotEqual(provider["request_identity"]["advertiser_id_sha256"], "123")
        self.assertEqual(batch["outcomes"][0]["status"], "CANDIDATE")
        self.assertEqual(receipt["candidate_count"], 1)
        self.assertEqual(receipt["retrieval_mode"], "TEST")
        self.assertFalse(receipt["automatic_business_write"])
        self.assertNotIn("person@example.test", json.dumps(receipt))

    def test_non_string_dynamic_field_is_rejected_without_data_loss(self):
        request = {"schema": "elite-tiktok-lead-retrieval-request/v1", "tenant_id": "11111111-1111-4111-8111-111111111111", "organization_id": "store-1", "mode": "LIVE", "lead_source": "DIRECT_MESSAGE", "advertiser_id": "123", "library_id": "", "page_id": ""}
        response = {"request_id": "req-2", "data": {"lead_data": {"multi": ["a", "b"]}, "meta_data": {"lead_source": "DIRECT_MESSAGE", "lead_id": "lead-2", "create_time": "2026-09-04 12:00:00"}}}
        artifacts = build_retrieval_artifacts(response, request)
        self.assertEqual(json.loads(artifacts["candidate-batch.json"])["outcomes"][0]["status"], "REJECTED")
        self.assertEqual(json.loads(artifacts["provider-response.json"])["data"]["lead_data"]["multi"], ["a", "b"])

    def test_retrieval_response_must_match_source_and_page(self):
        request = {"schema": "elite-tiktok-lead-retrieval-request/v1", "tenant_id": "11111111-1111-4111-8111-111111111111", "organization_id": "store-1", "mode": "TEST", "lead_source": "INSTANT_FORM", "advertiser_id": "123", "library_id": "", "page_id": "456"}
        response = {"request_id": "req", "data": {"lead_data": {"Email": "a"}, "meta_data": {"lead_source": "INSTANT_FORM", "lead_id": "lead", "page_id": "999", "create_time": "2026-09-04 12:00:00"}}}
        with self.assertRaises(TikTokLeadEvidenceError):
            build_retrieval_artifacts(response, request)
        response["data"]["meta_data"]["page_id"] = "456"
        response["data"]["meta_data"]["lead_source"] = "DIRECT_MESSAGE"
        with self.assertRaises(TikTokLeadEvidenceError):
            build_retrieval_artifacts(response, request)

    def test_atomic_artifact_write_rejects_overwrite(self):
        request = {"schema": "elite-tiktok-lead-retrieval-request/v1", "tenant_id": "11111111-1111-4111-8111-111111111111", "organization_id": "store-1", "mode": "LIVE", "lead_source": "DIRECT_MESSAGE", "advertiser_id": "123", "library_id": "", "page_id": ""}
        response = {"request_id": "req", "data": {"lead_data": {"Email": "a"}, "meta_data": {"lead_source": "DIRECT_MESSAGE", "lead_id": "lead", "create_time": "2026-09-04 12:00:00"}}}
        artifacts = build_retrieval_artifacts(response, request)
        with tempfile.TemporaryDirectory() as temporary:
            output = Path(temporary) / "evidence"
            write_artifacts_atomically(output, artifacts)
            self.assertEqual(sorted(path.name for path in output.iterdir()), ["candidate-batch.json", "provider-response.json", "retrieval-receipt.json"])
            with self.assertRaises(TikTokLeadEvidenceError):
                write_artifacts_atomically(output, artifacts)

    def test_retrieval_receipt_rejects_missing_identity_and_bad_time(self):
        request = {"schema": "elite-tiktok-lead-retrieval-request/v1", "tenant_id": "11111111-1111-4111-8111-111111111111", "organization_id": "store-1", "mode": "LIVE", "lead_source": "DIRECT_MESSAGE", "advertiser_id": "123", "library_id": "", "page_id": ""}
        response = {"request_id": "req", "data": {"lead_data": {"Email": "a"}, "meta_data": {"lead_source": "DIRECT_MESSAGE", "lead_id": "", "create_time": "bad"}}}
        with self.assertRaises(TikTokLeadEvidenceError):
            build_retrieval_artifacts(response, request)


if __name__ == "__main__":
    unittest.main()
````

### FILE: `tiktok_lead_adapter/README.md`
```yaml
block_id: "PY-TIKTOK-LEAD:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local operator and limitations contract"
license: "LicenseRef-Workspace-Owner"
sha256: "9e6bbffdfed367b5e37cca778e025f381852d5ab092167a5e5bdb2fdc8e008cf"
variables: []
secrets_allowed: false
```
````markdown
# TikTok Lead Adapter

This pack uses the exact official TikTok Business API wheel `tiktok-business-api-sdk-official==1.1.3` and its generic `ApiClient.call_api` transport. The wheel's distribution metadata says `1.1.3`, while its own `business_api_client.__version__` constant says `1.2.1`; the artifact lock and tests preserve both upstream facts and use filename plus SHA-256 as executable identity. TikTok publishes the Lead retrieval and subscription contracts in its v1.3 portal, but the pinned SDK does not contain provider-generated Lead methods. `adapter.py` is therefore explicitly `ADAPTED`: it follows the official generated-method request shape without claiming TikTok authored this wrapper.

The provider profile starts blocked and is a mandatory constructor argument; client construction fails unless every required proof is `PROVEN`. Before any live call, prove account/app/Lead access, terms and consent, a test account, callback control and TLS, provider delivery, durable idempotency, authenticated retrieval reconciliation, retention/deletion and quota/cost. Secrets belong only in the named environment variables.

TikTok documents webhook delivery as at least once, so duplicates are normal. `webhook.py` parses the documented Lead envelope and emits stable hashes, but marks every inbound signal `UNAUTHENTICATED_PROVIDER_SIGNAL`: the reviewed public documentation did not demonstrate a signature mechanism. The pack never persists a webhook automatically. A project must use durable deduplication, retrieve/reconcile through the authenticated API and keep an authorized audit trail.

```powershell
python -m venv .venv
.venv\Scripts\python -m pip install --require-hashes --no-deps -r .\tiktok_lead_adapter\requirements.lock
.venv\Scripts\python -m pip check
Push-Location .\tiktok_lead_adapter
..\.venv\Scripts\python -m unittest -v test_tiktok_lead.py
Pop-Location
```

`run_retrieval.py` performs one approved authenticated GET and writes a provider SDK response, candidate batch and retrieval receipt atomically. The receipt links both data artifacts by SHA-256 and contains no lead fields. Dynamic values are admitted as candidates only when they are non-empty strings; arrays, objects, nulls or empty values remain intact in provider evidence and produce a rejected outcome requiring an approved mapping.

The sixteen tests prove mandatory profile enforcement, exact endpoint/method/header/body construction, input gates, secret-safe errors, official wheel transport identity, hash-linked retrieval artifacts, atomic no-overwrite output and deterministic parsing/deduplication signals. They do not prove live permissions, provider delivery, webhook authenticity, field semantics, consent, reconciliation, quotas, cost or production persistence.
````

## 6. Configuration surface

No incluye secretos. El proyecto debe completar el profile sin cambiar `automatic_business_write=false`; app ID, secret y access token se resuelven desde las variables de entorno nombradas. `PROVEN` exige todos los receipts live del perfil. Advertiser/library/page, callback, consentimiento, retención y reconciliación pertenecen al blueprint del proyecto.

## 7. Dependency bill

| Dependencia | Pin verificado | Licencia | Uso |
|---|---|---|---|
| `tiktok-business-api-sdk-official` | 1.1.3 wheel SHA `663b4a…f33b7` | MIT | transporte oficial `ApiClient.call_api` |
| `urllib3` | 2.7.0 | MIT | transporte del SDK |
| `python-dateutil` | 2.9.0.post0 | Apache-2.0/BSD | dependencia SDK |
| `certifi` / `six` | 2026.7.22 / 1.17.0 | MPL-2.0 / MIT | runtime transitivo |

## 8. Apply order

1. Materializar y verificar los once archivos.
2. Instalar los cinco wheels exactos en un entorno aislado con hash enforcement y ejecutar `pip check`.
3. Completar el perfil y mantener secretos sólo en variables del entorno/secret manager.
4. Ejecutar los dieciséis tests antes de cualquier acceso.
5. Probar suscripción y delivery en cuenta de prueba; conservar request/subscription IDs redacted.
6. Demostrar deduplicación durable y retrieval/reconciliación antes de persistir un lead.

## 9. Verification

```powershell
python -m venv .venv
.venv\Scripts\python -m pip install --require-hashes --no-deps -r .\tiktok_lead_adapter\requirements.lock
.venv\Scripts\python -m pip check
Push-Location .\tiktok_lead_adapter
..\.venv\Scripts\python -m unittest -v test_tiktok_lead.py
Pop-Location
```

Resultado local esperado y demostrado sólo después del round-trip: wheel/dependencias exactos, `pip check`, dieciséis tests, contrato del transporte oficial y tres artefactos hash-linked/no-overwrite. No demuestra cuenta real, webhook auténtico, delivery, deduplicación durable, consentimiento, semántica de campos, reconciliación, cuota/costo ni persistencia productiva.

## 10. Reconstruction evidence

- Wheel oficial TikTok 1.1.3: URL y SHA-256 exactos, autor TikTok Pte. Ltd., licencia MIT.
- Repositorio oficial fijado: commit `f809c396520df2d7b201a9ccc5378d822b728ed3`; el árbol no contiene cliente Lead/subscription generado.
- Portal oficial v1.3: endpoints, headers, parámetros, response fields y garantía webhook at-least-once observados el 2026-09-04.
- La procedencia por bloque impide presentar el wrapper local como código Lead publicado por TikTok.
