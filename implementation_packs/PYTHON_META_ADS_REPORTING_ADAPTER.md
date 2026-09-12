# Python Meta Ads Reporting Adapter

## 1. Metadata

```yaml
pack_id: "PYTHON-META-ADS-REPORTING-ADAPTER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un adapter read-only sobre el SDK oficial Meta Business 26.0.1/Graph API v26.0 para ejecutar Ads Insights aprobados, preservar filas y emitir evidencia atómica sin exponer el ad account ID."
stacks: ["CPython 3.14 Windows x86-64", "facebook-business 26.0.1", "Meta Graph API v26.0"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x"]
incompatible_with: ["mutación de campañas", "Conversions API", "secreto en CLI/Markdown", "fields no aprobados", "write automático"]
license_expression: "LicenseRef-Workspace-Owner AND LicenseRef-Meta-Platform AND third-party-wheel-licenses"
upstream_sources: ["https://github.com/facebook/facebook-python-business-sdk/tree/788f363d15b1269ab5efb7cd00fb5e3b133cd99b", "https://github.com/facebook/capi-param-builder/tree/66a7eb84e29a65aa999089366c1f7c3c0e49bdb7", "https://pypi.org/project/facebook-business/26.0.1/"]
verified_at: "2026-08-26"
```

## 2. Applicability

Use sólo para reporting Meta Ads cuando el usuario elija Meta y cierre registro de app, verificación de negocio, permiso `ads_read`, cuenta publicitaria test, ownership, fields, fechas, retención, cuota/costo y reconciliación. No usar para crear/modificar campañas, Conversions API, atribución autónoma ni escritura directa al dominio.

El adapter es glue `AUTHORED` sobre `FacebookAdsApi`, `AdAccount` y `AdAccount.get_insights` oficiales. El SDK/modelos son Meta; el código local no se atribuye a Meta. La licencia oficial limita ambos SDKs a servicios/APIs Facebook y debe conservarse. El wheel `capi-param-builder-python` 1.3.0 omite `LICENSE`; el source oficial fijado sí contiene la misma licencia y gobierna su redistribución.

## 3. Architecture contract

La plantilla nace bloqueada. Sólo acepta Graph API `v26.0`, SDK exacto, cuenta numérica, fields localmente permitidos y explícitamente aprobados, nivel conocido, rango ISO ordenado de hasta 31 días, límite 1..1000, secrets desde variables de entorno y output nuevo. Ejecuta únicamente `GET /insights` por el SDK oficial. Conserva filas y recibo hash-linked de forma atómica sin cuenta ni secretos en el receipt.

Error de profile, versión, input, SDK/provider, serialización u output elimina staging. Resultado vacío es evidencia válida. Fixtures no prueban permisos, disponibilidad de métricas, atribución, cuota, costo, reconciliación ni equivalencia con UI/export; esos gates requieren la cuenta real seleccionada.

## 4. Exact file manifest

```text
CREATE meta_ads_reporting/requirements-direct.in
CREATE meta_ads_reporting/requirements-windows-py314.lock
CREATE meta_ads_reporting/sdk-artifact.lock.json
CREATE meta_ads_reporting/provider-profile.template.json
CREATE meta_ads_reporting/query.template.json
CREATE meta_ads_reporting/run_insights.py
CREATE meta_ads_reporting/test_insights.py
CREATE meta_ads_reporting/README.md
```

## 5. Materialization blocks

### FILE: `meta_ads_reporting/requirements-direct.in`
```yaml
block_id: "PY-META-ADS-REPORTING:requirements-direct:v1"
operation: CREATE
provenance: AUTHORED
source: "exact official PyPI wheel URL and SHA-256"
license: "LicenseRef-Workspace-Owner"
sha256: "b0fe2f0d83c899fd9986a9990b9102b613774bd912aaffd5219736b014051c88"
variables: []
secrets_allowed: false
```
````text
facebook-business @ https://files.pythonhosted.org/packages/bd/fb/038e02129ff35dbe1eac38bd7ab40234a08576359ef8806ca6d132e2d558/facebook_business-26.0.1-py3-none-any.whl#sha256=41555ce87ac105617ed678a489a630108ac162fbdccfbbf65fa20ae1de7cd0ca
````

### FILE: `meta_ads_reporting/requirements-windows-py314.lock`
```yaml
block_id: "PY-META-ADS-REPORTING:requirements-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "exact PyPI wheel graph resolved for CPython 3.14 Windows x86-64"
license: "LicenseRef-Workspace-Owner"
sha256: "a4b4b65d2a19be4bb51892bcb25a4e8127f7b2c16b8f49a2f40b319797314e19"
variables: []
secrets_allowed: false
```
````text
# Exact wheel graph resolved and verified for CPython 3.14 on Windows x86-64.
aiohappyeyeballs @ https://files.pythonhosted.org/packages/71/43/1947f06babed6b3f1d7f38b0c767f52df66bfb2bc10b468c4a7de9eceff2/aiohappyeyeballs-2.7.1-py3-none-any.whl#sha256=9243213661e29250eb41368e5daa826fc017156c3b8a11440826b2e3ed376472
aiohttp @ https://files.pythonhosted.org/packages/f5/8b/c7baa1ba1eda4db6989baefe5de6d99834921b84ebd7918624febcb9f290/aiohttp-3.14.3-cp314-cp314-win_amd64.whl#sha256=8b3b60de05f3dcb6f6a00f818bb2ec781cee4de0645f59ccaf99b1d1823b6100
aiosignal @ https://files.pythonhosted.org/packages/fb/76/641ae371508676492379f16e2fa48f4e2c11741bd63c48be4b12a6b09cba/aiosignal-1.4.0-py3-none-any.whl#sha256=053243f8b92b990551949e63930a839ff0cf0b0ebbe0597b0f3fb19e1a0fe82e
attrs @ https://files.pythonhosted.org/packages/64/b4/17d4b0b2a2dc85a6df63d1157e028ed19f90d4cd97c36717afef2bc2f395/attrs-26.1.0-py3-none-any.whl#sha256=c647aa4a12dfbad9333ca4e71fe62ddc36f4e63b2d260a37a8b83d2f043ac309
capi-param-builder-python @ https://files.pythonhosted.org/packages/4f/1d/463144a77228b0137b498f70eb4ac60808fddba7bc890308729b76197391/capi_param_builder_python-1.3.0-py3-none-any.whl#sha256=641fd8832378033670783a8841dc3168fd75f5ab9629779d11a6ce0a6687f86d
certifi @ https://files.pythonhosted.org/packages/0b/a7/71ac2cff56fec219ed242bb11b8efb69fcc4bec75db06fb7bfe35de520e6/certifi-2026.7.22-py3-none-any.whl#sha256=62f22742b58a1a33014a2b6b706588a8d7e2a88ae7bd1a6ebe8c992928483775
charset-normalizer @ https://files.pythonhosted.org/packages/7a/7c/4938c329b6a9d446f6a59aa2092ff7118f274209b5ed0e26893d1d30a63c/charset_normalizer-3.5.1-cp314-cp314-win_amd64.whl#sha256=c658c50ac0c98cd755a2dd50b7977d3bca7df401dcc47fbdfa87db53ef7d4e8b
curlify @ https://files.pythonhosted.org/packages/9e/f8/912ebddbff8ea603d4c90fa31557096f927b17efd30a166ce7ac1242910a/curlify-3.0.0-py3-none-any.whl#sha256=52060c0eb7a656b7bde6b668c32f337bed4d736ce230755767e3ada56a09c338
facebook-business @ https://files.pythonhosted.org/packages/bd/fb/038e02129ff35dbe1eac38bd7ab40234a08576359ef8806ca6d132e2d558/facebook_business-26.0.1-py3-none-any.whl#sha256=41555ce87ac105617ed678a489a630108ac162fbdccfbbf65fa20ae1de7cd0ca
frozenlist @ https://files.pythonhosted.org/packages/59/ad/9caa9b9c836d9ad6f067157a531ac48b7d36499f5036d4141ce78c230b1b/frozenlist-1.8.0-cp314-cp314-win_amd64.whl#sha256=3e0761f4d1a44f1d1a47996511752cf3dcec5bbdd9cc2b4fe595caf97754b7a0
idna @ https://files.pythonhosted.org/packages/57/b0/0e52c878c53f245edd3a11020f20979b3f490f245af532c7cae3027754b5/idna-3.19-py3-none-any.whl#sha256=815e7be7a7806d54abb586dc943addc79e8b2ee16915059658cbeff4b1b43bf4
multidict @ https://files.pythonhosted.org/packages/e0/bf/52f25716bbe93745595800f36fb17b73711f14da59ed0bb2eba141bc9f0f/multidict-6.7.1-cp314-cp314-win_amd64.whl#sha256=5e01429a929600e7dab7b166062d9bb54a5eed752384c7384c968c2afab8f50f
propcache @ https://files.pythonhosted.org/packages/61/d2/45c9defbaa1ea297035d9d4cce9e8f80daafbf19319c6007f157c6256ea9/propcache-0.5.2-cp314-cp314-win_amd64.whl#sha256=81e3a30b0bb60caa22033dd0f8a3618d1d67356212514f62c57db75cb0ef410c
pycountry @ https://files.pythonhosted.org/packages/9c/42/7703bd45b62fecd44cd7d3495423097e2f7d28bc2e99e7c1af68892ab157/pycountry-26.2.16-py3-none-any.whl#sha256=115c4baf7cceaa30f59a4694d79483c9167dbce7a9de4d3d571c5f3ea77c305a
requests @ https://files.pythonhosted.org/packages/a0/f4/c67b0b3f1b9245e8d266f0f112c500d50e5b4e83cb6f3b71b6528104182a/requests-2.34.2-py3-none-any.whl#sha256=2a0d60c172f83ac6ab31e4554906c0f3b3588d37b5cb939b1c061f4907e278e0
six @ https://files.pythonhosted.org/packages/b7/ce/149a00dd41f10bc29e5921b496af8b574d8413afcd5e30dfa0ed46c2cc5e/six-1.17.0-py2.py3-none-any.whl#sha256=4721f391ed90541fddacab5acf947aa0d3dc7d27b2e1e8eda2be8970586c3274
urllib3 @ https://files.pythonhosted.org/packages/7f/3e/5db95bcf282c52709639744ca2a8b149baccf648e39c8cc87553df9eae0c/urllib3-2.7.0-py3-none-any.whl#sha256=9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897
yarl @ https://files.pythonhosted.org/packages/cf/52/6daa2ee9d95e5c98b8128f8df91eb692eb423ab274b8cf08db52152fad26/yarl-1.24.5-cp314-cp314-win_amd64.whl#sha256=5ba4f78df2bcc19f764a4b26a8a4f5049c110090ad5825993aacb052bf8003ad
````

### FILE: `meta_ads_reporting/sdk-artifact.lock.json`
```yaml
block_id: "PY-META-ADS-REPORTING:sdk-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "local receipt derived from official GitHub/PyPI metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "03eb5d7bd845cd02d42045e1d74ee307976cc31eadc5b8c4925b222f86c5a887"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-python-sdk-artifact-lock/v1",
  "provider": "Meta",
  "distribution": "facebook-business",
  "version": "26.0.1",
  "verified_platform_lock": "CPython 3.14 / Windows x86-64 / 18 wheels",
  "wheel": {
    "filename": "facebook_business-26.0.1-py3-none-any.whl",
    "bytes": 1554832,
    "sha256": "41555ce87ac105617ed678a489a630108ac162fbdccfbbf65fa20ae1de7cd0ca"
  },
  "source": {
    "repository": "facebook/facebook-python-business-sdk",
    "release": "26.0.1",
    "commit": "788f363d15b1269ab5efb7cd00fb5e3b133cd99b",
    "archive_bytes": 2070568,
    "archive_sha256": "bd3bc1f14072662eecb05fb6f43888eb652c9a1f039894f99a7e7609e4452db7",
    "commit_signature_verified": false
  },
  "capi_parameter_builder_source": {
    "repository": "facebook/capi-param-builder",
    "release": "v1.3.0-python",
    "commit": "66a7eb84e29a65aa999089366c1f7c3c0e49bdb7",
    "archive_bytes": 674665,
    "archive_sha256": "be8d21efab5a1ea16b1f6e84d4219d9e343436b8789ede407d0b94d0a54dabf5"
  },
  "license_expression": "LicenseRef-Meta-Platform",
  "license_sha256": "48d97b3c936203a750a3288c2b924327769d62327383a1060244b3d3c05ad01f",
  "distribution_notice": "capi-param-builder-python 1.3.0 wheel omits LICENSE; preserve the exact official source LICENSE with redistributions",
  "verified_at": "2026-08-26"
}
````

### FILE: `meta_ads_reporting/provider-profile.template.json`
```yaml
block_id: "PY-META-ADS-REPORTING:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed Meta access/query profile"
license: "LicenseRef-Workspace-Owner"
sha256: "36eb179d5afcc0c20ab1608eeaa04b328402b1c9b7275b83a5e1ba2b0799bfed"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-meta-ads-reporting-profile/v1",
  "provider": "Meta Marketing API",
  "sdk": "facebook-business==26.0.1",
  "graph_api_version": "v26.0",
  "decision": "BLOCKED_ACCESS_AND_QUERY_APPROVAL_REQUIRED",
  "app_registration_proven": false,
  "business_verification_proven": false,
  "ads_read_permission_proven": false,
  "test_ad_account_contract_proven": false,
  "app_id_environment_variable": "META_APP_ID",
  "app_secret_environment_variable": "META_APP_SECRET",
  "access_token_environment_variable": "META_ACCESS_TOKEN",
  "ad_account_id": "",
  "approved_fields": [],
  "data_retention": "",
  "quota_and_cost_approved": false,
  "reconciliation_approved": false,
  "automatic_business_write": false
}
````

### FILE: `meta_ads_reporting/query.template.json`
```yaml
block_id: "PY-META-ADS-REPORTING:query:v1"
operation: CREATE
provenance: AUTHORED
source: "local bounded Ads Insights query template"
license: "LicenseRef-Workspace-Owner"
sha256: "5b328270f544d80a9a8f8288e02d15a4b8cad6529c7f240db2124c5b3e4b890d"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-meta-ads-insights-query/v1",
  "fields": ["campaign_id", "campaign_name", "date_start", "date_stop", "impressions", "clicks", "spend"],
  "level": "campaign",
  "since": "2026-08-01",
  "until": "2026-08-02",
  "limit": 100
}
````

### FILE: `meta_ads_reporting/run_insights.py`
```yaml
block_id: "PY-META-ADS-REPORTING:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local wrapper invoking official Meta Business SDK APIs"
license: "LicenseRef-Workspace-Owner"
sha256: "777a5754bc9e7774baa5dfb819967225149115b6bb0f9b2c780be349497d1171"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
from datetime import date, datetime, timezone
import hashlib
import importlib.metadata
import json
import os
from pathlib import Path
import re
import shutil
import sys
from typing import Any, Callable
import uuid

from facebook_business.adobjects.adaccount import AdAccount
from facebook_business.api import FacebookAdsApi


SDK_VERSION = "26.0.1"
API_VERSION = "v26.0"
ALLOWED_FIELDS = {
    "account_id", "campaign_id", "campaign_name", "adset_id", "adset_name",
    "ad_id", "ad_name", "date_start", "date_stop", "impressions", "reach",
    "clicks", "spend",
}
ALLOWED_LEVELS = {"account", "campaign", "adset", "ad"}


def sha256(value: bytes) -> str:
    return hashlib.sha256(value).hexdigest()


def _read_json(path: Path) -> dict[str, Any]:
    value = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(value, dict):
        raise ValueError(f"{path.name} must contain one JSON object")
    return value


def validate_configuration(profile: dict[str, Any], query: dict[str, Any]) -> tuple[str, list[str], dict[str, Any]]:
    required = (
        "app_registration_proven", "business_verification_proven", "ads_read_permission_proven",
        "test_ad_account_contract_proven", "quota_and_cost_approved", "reconciliation_approved",
    )
    if profile.get("provider") != "Meta Marketing API" or profile.get("sdk") != f"facebook-business=={SDK_VERSION}" or profile.get("graph_api_version") != API_VERSION:
        raise ValueError("profile provider/SDK/API identity mismatch")
    if profile.get("decision") != "PROVEN" or any(profile.get(key) is not True for key in required):
        raise PermissionError("Meta access/query profile remains blocked")
    if profile.get("automatic_business_write") is not False:
        raise ValueError("automatic_business_write must remain false")
    if not str(profile.get("data_retention", "")).strip():
        raise ValueError("data_retention is required")
    account_id = str(profile.get("ad_account_id", ""))
    if not re.fullmatch(r"[0-9]{5,30}", account_id):
        raise ValueError("ad_account_id must contain digits only")
    fields = query.get("fields")
    approved = profile.get("approved_fields")
    if not isinstance(fields, list) or not fields or len(fields) > 20 or len(fields) != len(set(fields)):
        raise ValueError("fields must be a unique non-empty list with at most 20 entries")
    if not isinstance(approved, list) or not set(fields).issubset(set(approved)) or not set(fields).issubset(ALLOWED_FIELDS):
        raise PermissionError("all fields must be locally allowed and explicitly approved")
    level = query.get("level")
    if level not in ALLOWED_LEVELS:
        raise ValueError("unsupported reporting level")
    since = date.fromisoformat(str(query.get("since", "")))
    until = date.fromisoformat(str(query.get("until", "")))
    if until < since or (until - since).days > 31:
        raise ValueError("time range must be ordered and at most 31 days")
    limit = query.get("limit")
    if not isinstance(limit, int) or isinstance(limit, bool) or limit < 1 or limit > 1000:
        raise ValueError("limit must be within 1..1000")
    params = {"level": level, "time_range": {"since": since.isoformat(), "until": until.isoformat()}, "limit": limit}
    return account_id, fields, params


def _row_to_dict(row: Any) -> dict[str, Any]:
    if isinstance(row, dict):
        return row
    exporter = getattr(row, "export_all_data", None)
    if not callable(exporter):
        raise RuntimeError("Meta Ads result row is not an official SDK object")
    value = exporter()
    if not isinstance(value, dict):
        raise RuntimeError("Meta SDK row export must return an object")
    return value


def query_to_evidence(account_factory: Callable[[str], Any], account_id: str, fields: list[str], params: dict[str, Any], output_directory: Path) -> dict[str, Any]:
    if importlib.metadata.version("facebook-business") != SDK_VERSION:
        raise RuntimeError(f"facebook-business {SDK_VERSION} is required")
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".meta-ads-report-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        rows = [_row_to_dict(row) for row in account_factory(f"act_{account_id}").get_insights(fields=fields, params=params)]
        response = {"graph_api_version": API_VERSION, "rows": rows}
        response_bytes = (json.dumps(response, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "provider-response.json").write_bytes(response_bytes)
        query_bytes = json.dumps({"fields": fields, "params": params}, sort_keys=True, separators=(",", ":")).encode("utf-8")
        receipt = {
            "schema": "elite-meta-ads-reporting-receipt/v1",
            "created_at": datetime.now(timezone.utc).isoformat(),
            "provider": "Meta Marketing API",
            "sdk_version": SDK_VERSION,
            "graph_api_version": API_VERSION,
            "ad_account_id_sha256": sha256(account_id.encode("ascii")),
            "query_sha256": sha256(query_bytes),
            "row_count": len(rows),
            "provider_response_sha256": sha256(response_bytes),
            "automatic_business_write": False,
        }
        (stage / "QUERY_RECEIPT.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        stage.replace(output)
        return receipt
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--query", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()
    profile, query = _read_json(args.profile), _read_json(args.query)
    account_id, fields, params = validate_configuration(profile, query)
    names = ("app_id_environment_variable", "app_secret_environment_variable", "access_token_environment_variable")
    environment_names = [str(profile.get(name, "")) for name in names]
    if any(not re.fullmatch(r"[A-Z][A-Z0-9_]{2,80}", name) for name in environment_names):
        raise ValueError("secret references must be canonical environment variable names")
    values = [os.environ.get(name, "") for name in environment_names]
    if any(not value for value in values):
        raise PermissionError("required Meta credential environment variables are unavailable")
    api = FacebookAdsApi.init(app_id=values[0], app_secret=values[1], access_token=values[2], api_version=API_VERSION)
    receipt = query_to_evidence(lambda identifier: AdAccount(identifier, api=api), account_id, fields, params, args.output)
    print(json.dumps({"status": "META_ADS_REPORTING_PASS", "receipt": receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `meta_ads_reporting/test_insights.py`
```yaml
block_id: "PY-META-ADS-REPORTING:test:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract/negative tests plus official SDK signature probe"
license: "LicenseRef-Workspace-Owner"
sha256: "10ca41eba4749e8afc9090ef413f01f5b0cca3c0b6d36d8774bae1152ab65681"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import inspect
import json
from pathlib import Path
import tempfile
import unittest

from facebook_business.adobjects.adaccount import AdAccount
from run_insights import query_to_evidence, validate_configuration


class FakeRow:
    def __init__(self, value): self.value = value
    def export_all_data(self): return self.value


class FakeAccount:
    def __init__(self, identifier, rows=None, error=None): self.identifier = identifier; self.rows = rows or []; self.error = error; self.calls = []
    def get_insights(self, *, fields, params):
        self.calls.append((fields, params))
        if self.error: raise self.error
        return self.rows


def proven_profile():
    return {
        "provider": "Meta Marketing API", "sdk": "facebook-business==26.0.1", "graph_api_version": "v26.0",
        "decision": "PROVEN", "app_registration_proven": True, "business_verification_proven": True,
        "ads_read_permission_proven": True, "test_ad_account_contract_proven": True,
        "quota_and_cost_approved": True, "reconciliation_approved": True,
        "automatic_business_write": False, "data_retention": "30 days", "ad_account_id": "123456789",
        "approved_fields": ["campaign_id", "impressions", "clicks"],
    }


def valid_query():
    return {"fields": ["campaign_id", "impressions"], "level": "campaign", "since": "2026-08-01", "until": "2026-08-02", "limit": 100}


class MetaInsightsTests(unittest.TestCase):
    def test_configuration_fails_closed(self):
        profile = proven_profile(); profile["decision"] = "BLOCKED_ACCESS_AND_QUERY_APPROVAL_REQUIRED"
        with self.assertRaises(PermissionError): validate_configuration(profile, valid_query())
        profile = proven_profile(); query = valid_query(); query["fields"] = ["actions"]
        with self.assertRaises(PermissionError): validate_configuration(profile, query)

    def test_configuration_rejects_range_and_account(self):
        profile = proven_profile(); profile["ad_account_id"] = "act_123"
        with self.assertRaises(ValueError): validate_configuration(profile, valid_query())
        profile = proven_profile(); query = valid_query(); query["until"] = "2026-10-01"
        with self.assertRaises(ValueError): validate_configuration(profile, query)

    def test_preserves_rows_and_redacts_account(self):
        account = FakeAccount("", [FakeRow({"campaign_id": "1", "impressions": "12"})])
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = query_to_evidence(lambda identifier: self._bind(account, identifier), "123456789", ["campaign_id", "impressions"], {"level": "campaign", "time_range": {"since": "2026-08-01", "until": "2026-08-02"}, "limit": 100}, output)
            self.assertEqual(account.identifier, "act_123456789")
            self.assertEqual(receipt["row_count"], 1)
            self.assertNotIn("123456789", json.dumps(receipt))
            self.assertEqual(json.loads((output / "provider-response.json").read_text(encoding="utf-8"))["rows"][0]["impressions"], "12")

    @staticmethod
    def _bind(account, identifier): account.identifier = identifier; return account

    def test_provider_failure_is_atomic(self):
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            with self.assertRaisesRegex(RuntimeError, "provider failed"):
                query_to_evidence(lambda identifier: FakeAccount(identifier, error=RuntimeError("provider failed")), "123456789", ["campaign_id"], {"level": "campaign", "time_range": {"since": "2026-08-01", "until": "2026-08-02"}, "limit": 1}, output)
            self.assertFalse(output.exists())
            self.assertEqual(list(Path(temp).iterdir()), [])

    def test_existing_output_is_rejected(self):
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"; output.mkdir()
            with self.assertRaises(FileExistsError):
                query_to_evidence(lambda identifier: FakeAccount(identifier), "123456789", ["campaign_id"], {"level": "campaign", "time_range": {"since": "2026-08-01", "until": "2026-08-02"}, "limit": 1}, output)

    def test_official_sdk_contract(self):
        signature = inspect.signature(AdAccount.get_insights)
        self.assertIn("fields", signature.parameters)
        self.assertIn("params", signature.parameters)
        self.assertIn("is_async", signature.parameters)


if __name__ == "__main__": unittest.main()
````

### FILE: `meta_ads_reporting/README.md`
```yaml
block_id: "PY-META-ADS-REPORTING:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local operator contract"
license: "LicenseRef-Workspace-Owner"
sha256: "6d12c50bc572ea91ece3c111b3ec04e8e16e8d75ae85a0b7f25be74162d458f8"
variables: []
secrets_allowed: false
```
````markdown
# Meta Ads Reporting Adapter

Read-only wrapper over the official Meta Business SDK for Python 26.0.1 and Graph API v26.0. It executes only `AdAccount.get_insights`, preserves the provider rows and creates a hash-linked receipt. The local wrapper is `AUTHORED`; it is not Meta code.

The template is blocked. Before use, prove app/business registration, `ads_read`, a test ad account, exact account ownership, field allowlist, retention, quota/cost and reconciliation. Put app ID, app secret and access token only in the named environment variables. Never place credentials in JSON, Markdown, arguments or output.

The frozen lock is specific to CPython 3.14 on Windows x86-64. For another target, resolve and admit a separate hash-complete wheel graph. The official Meta license restricts the SDK and `capi-param-builder` to Facebook web services/APIs. The `capi-param-builder-python` 1.3.0 wheel omits its license file, so preserve the exact root LICENSE from official source commit `66a7eb84e29a65aa999089366c1f7c3c0e49bdb7`.

```powershell
python -m venv .venv
.venv\Scripts\python -m pip install --require-hashes -r .\meta_ads_reporting\requirements-windows-py314.lock
Push-Location .\meta_ads_reporting
..\.venv\Scripts\python -m unittest -v test_insights.py
Pop-Location
```

No fixture proves real permissions, field availability, provider totals, attribution, rate limits or billing. Admit production only after sandbox/test-account comparison and reconciliation against Meta UI/export for the selected dates and fields.
````

## 6. Configuration surface

No contiene secretos. El proyecto debe completar el profile con decisión `PROVEN`, evidencia de app/business/`ads_read`/cuenta test, cuenta numérica, fields permitidos, retención, cuota/costo y reconciliación. Los tres secretos se resuelven exclusivamente desde los nombres de variables de entorno del profile. La query separada fija fields, nivel, fechas y límite.

## 7. Dependency bill

| Dependencia | Pin verificado | Licencia | Uso |
|---|---|---|---|
| `facebook-business` | 26.0.1, wheel SHA `41555ce8…0ca` | Meta Platform License | SDK oficial Ads Insights |
| `capi-param-builder-python` | 1.3.0, wheel SHA `641fd883…86d` | Meta Platform License desde source fijado; wheel sin LICENSE | dependencia SDK |
| `pycountry` | 26.2.16 | LGPL-2.1-only | dependencia SDK |
| requests/aiohttp y transitivas | 15 wheels exactos | Apache-2.0/MIT/BSD/PSF/MPL según distribución | transporte/runtime |

El lock de 18 wheels es target-specific para CPython 3.14 Windows x86-64. Otro target requiere lock separado, hash-completo y revisión de licencias. La licencia Meta no permite reutilizar este SDK fuera de sus web services/APIs.

## 8. Apply order

1. Materializar en destino vacío y verificar hashes.
2. Completar readiness/acceso/terms/costo sin secretos embebidos.
3. Crear venv CPython 3.14 Windows x86-64 e instalar `--require-hashes`.
4. Ejecutar tests offline/locales.
5. Comparar una query mínima contra test account y export/UI; aprobar reconciliación.
6. Habilitar job read-only acotado; rollback lo desactiva y conserva evidence.

## 9. Verification

```powershell
python -m venv .venv
.venv\Scripts\python -m pip install --disable-pip-version-check --require-hashes -r .\meta_ads_reporting\requirements-windows-py314.lock
.venv\Scripts\python -m pip check
Push-Location .\meta_ads_reporting
..\.venv\Scripts\python -m unittest -v test_insights.py
Pop-Location
```

Resultado verificado: 18-wheel hash install, `pip check` y seis tests PASS. Sandbox/test account, totals, attribution, limits/cost y reconciliación siguen `PROJECT_CONDITIONED`.

## 10. Reconstruction evidence

- `reconstruction_evidence/META_ADS_REPORTING_ADAPTER_2026-08-26_V1.md`
- SDK release/tag/commit/archive/wheel/licencia fijados; ambos commits oficiales observados estaban sin firma Git verificada.
- El wheel CAPI sin licencia embebida queda como condición explícita; licencia oficial preservada por source commit/hash.
