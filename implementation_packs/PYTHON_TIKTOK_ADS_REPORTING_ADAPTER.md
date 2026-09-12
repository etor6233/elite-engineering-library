# Python TikTok Ads Reporting Adapter

## 1. Metadata

```yaml
pack_id: "PYTHON-TIKTOK-ADS-REPORTING-ADAPTER"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un adapter read-only sobre el wheel oficial TikTok Business API 1.1.3 fijado por URL y SHA-256 para ejecutar report_integrated_get aprobado, preservar páginas y emitir evidencia atómica sin exponer advertiser ID ni token."
stacks: ["CPython 3.14 Windows x86-64", "tiktok-business-api-sdk-official 1.1.3", "TikTok Business API v1.3"]
compatible_with: ["Python 3.14 hash-locked environments"]
incompatible_with: ["mutación de anuncios", "source móvil", "secreto en CLI/Markdown", "campos no aprobados", "write automático"]
license_expression: "LicenseRef-Workspace-Owner AND MIT AND third-party-wheel-licenses"
upstream_sources: ["https://pypi.org/project/tiktok-business-api-sdk-official/1.1.3/", "https://github.com/tiktok/tiktok-business-api-sdk"]
verified_at: "2026-09-04"
```

## 2. Applicability

Use sólo para reporting TikTok Ads cuando el usuario elija TikTok y demuestre Business Center/developer app, acceso al advertiser test, reporting scope, dimensions/metrics exactas, fechas, retención, cuota/costo y reconciliación. El adapter usa la distribución PyPI oficial recomendada por el README del proveedor, fijada por URL/SHA-256 y autor `TikTok Pte. Ltd.`. No usar para mutaciones, atribución autónoma ni escritura al dominio.

## 3. Architecture contract

La plantilla nace bloqueada. Exige `tiktok-business-api-sdk-official==1.1.3` antes de importar el SDK oficial. Sólo ejecuta `ReportingApi.report_integrated_get`, `GET /open_api/v1.3/report/integrated/get/`, con token desde una variable de entorno. Limita reporte BASIC, nivel admitido, listas aprobadas, rango máximo de 31 días, página 1..1000 y máximo 100 páginas. Conserva páginas y receipt hash-linked atómicamente, sin advertiser ID ni secreto. Los errores de versión/acceso/configuración/provider/paginación/output fallan cerrados y eliminan staging. Los fixtures no prueban acceso real, métricas, atribución, límites, costo o reconciliación.

## 4. Exact file manifest

```text
CREATE tiktok_ads_reporting/requirements-direct.in
CREATE tiktok_ads_reporting/requirements.lock
CREATE tiktok_ads_reporting/sdk-artifact.lock.json
CREATE tiktok_ads_reporting/provider-profile.template.json
CREATE tiktok_ads_reporting/query.template.json
CREATE tiktok_ads_reporting/run_reporting.py
CREATE tiktok_ads_reporting/test_reporting.py
CREATE tiktok_ads_reporting/README.md
```

## 5. Materialization blocks

### FILE: `tiktok_ads_reporting/requirements-direct.in`
```yaml
block_id: "PY-TIKTOK-ADS-REPORTING:requirements-direct:v1"
operation: CREATE
provenance: AUTHORED
source: "exact official SDK runtime dependencies"
license: "LicenseRef-Workspace-Owner"
sha256: "39b0e5b11026ce226f914edc7ee439ff52cdb19c314a15547635215c0ec3954b"
variables: []
secrets_allowed: false
```
````text
certifi==2026.7.22
python-dateutil==2.9.0.post0
six==1.17.0
tiktok-business-api-sdk-official==1.1.3
urllib3==2.7.0
````

### FILE: `tiktok_ads_reporting/requirements.lock`
```yaml
block_id: "PY-TIKTOK-ADS-REPORTING:requirements-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "exact PyPI wheels for official SDK runtime"
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

### FILE: `tiktok_ads_reporting/sdk-artifact.lock.json`
```yaml
block_id: "PY-TIKTOK-ADS-REPORTING:sdk-artifact-lock:v2"
operation: CREATE
provenance: AUTHORED
source: "official TikTok PyPI wheel identity and publisher evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "7244ba16c0748bf090e49dccf8a2c3f5eb01bda78dc883a03a1091ce1c2af3bd"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-python-sdk-artifact-lock/v1",
  "provider": "TikTok Pte. Ltd.",
  "distribution": "tiktok-business-api-sdk-official",
  "version": "1.1.3",
  "requires_python": ">=3.4",
  "wheel": {
    "filename": "tiktok_business_api_sdk_official-1.1.3-py3-none-any.whl",
    "bytes": 678182,
    "sha256": "663b4a1f2585c4b386144d968646f695befd733c464418bbaeeac41befff33b7",
    "url": "https://files.pythonhosted.org/packages/ec/bd/1cf2faaa9fb43a4e81dec0cbe1076c1ffe0d2536611650ff7891a68a932b/tiktok_business_api_sdk_official-1.1.3-py3-none-any.whl",
    "uploaded_at": "2026-02-27T22:42:23.997612Z"
  },
  "publisher_evidence": {
    "pypi_author": "TikTok Pte. Ltd.",
    "repository": "https://github.com/tiktok/tiktok-business-api-sdk",
    "provider_readme_recommends_distribution": true
  },
  "resolved_environment": {
    "python": "3.14.4",
    "platform": "Windows x86-64",
    "pip": "26.0.1",
    "dependency_count": 5,
    "resolved_at": "2026-09-04"
  },
  "api_version": "v1.3",
  "license_expression": "MIT",
  "automatic_business_write": false
}
````

### FILE: `tiktok_ads_reporting/provider-profile.template.json`
```yaml
block_id: "PY-TIKTOK-ADS-REPORTING:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed TikTok access/query profile"
license: "LicenseRef-Workspace-Owner"
sha256: "26b941545bc041eb4a41025ae286ec7145fd55d364e286732798b86c77cc5d5e"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-tiktok-ads-reporting-profile/v1",
  "provider": "TikTok Business API",
  "sdk_distribution": "tiktok-business-api-sdk-official",
  "sdk_version": "1.1.3",
  "decision": "BLOCKED_ACCESS_AND_QUERY_APPROVAL_REQUIRED",
  "business_account_proven": false,
  "developer_app_proven": false,
  "reporting_scope_proven": false,
  "test_advertiser_contract_proven": false,
  "access_token_environment_variable": "TIKTOK_ACCESS_TOKEN",
  "advertiser_id": "",
  "approved_dimensions": [],
  "approved_metrics": [],
  "data_retention": "",
  "quota_and_cost_approved": false,
  "reconciliation_approved": false,
  "automatic_business_write": false
}
````

### FILE: `tiktok_ads_reporting/query.template.json`
```yaml
block_id: "PY-TIKTOK-ADS-REPORTING:query:v1"
operation: CREATE
provenance: AUTHORED
source: "local bounded official reporting query template"
license: "LicenseRef-Workspace-Owner"
sha256: "b7d712a121f1603612b7ab1741030f45e2feb48fa92f86e594b57fbb376ea65d"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-tiktok-integrated-report-query/v1",
  "report_type": "BASIC",
  "data_level": "AUCTION_CAMPAIGN",
  "dimensions": ["campaign_id", "stat_time_day"],
  "metrics": ["spend", "impressions", "clicks"],
  "start_date": "2026-08-01",
  "end_date": "2026-08-02",
  "page_size": 1000
}
````

### FILE: `tiktok_ads_reporting/run_reporting.py`
```yaml
block_id: "PY-TIKTOK-ADS-REPORTING:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local wrapper over official ReportingApi.report_integrated_get"
license: "LicenseRef-Workspace-Owner"
sha256: "cb4a0b60b56b6de8339f74d656cb8c060fc1521e4e4cd76571c4f27041eb852b"
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


SDK_DISTRIBUTION = "tiktok-business-api-sdk-official"
SDK_VERSION = "1.1.3"
ALLOWED_LEVELS = {"AUCTION_ADVERTISER", "AUCTION_CAMPAIGN", "AUCTION_ADGROUP", "AUCTION_AD"}


def sha256(value: bytes) -> str:
    return hashlib.sha256(value).hexdigest()


def _read_object(path: Path) -> dict[str, Any]:
    value = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(value, dict):
        raise ValueError(f"{path.name} must contain one JSON object")
    return value


def validate_configuration(profile: dict[str, Any], query: dict[str, Any]) -> tuple[str, dict[str, Any]]:
    required = ("business_account_proven", "developer_app_proven", "reporting_scope_proven", "test_advertiser_contract_proven", "quota_and_cost_approved", "reconciliation_approved")
    if profile.get("provider") != "TikTok Business API" or profile.get("sdk_distribution") != SDK_DISTRIBUTION or profile.get("sdk_version") != SDK_VERSION:
        raise ValueError("TikTok profile authority mismatch")
    if profile.get("decision") != "PROVEN" or any(profile.get(key) is not True for key in required):
        raise PermissionError("TikTok access/query profile remains blocked")
    if profile.get("automatic_business_write") is not False or not str(profile.get("data_retention", "")).strip():
        raise ValueError("retention is required and automatic_business_write must remain false")
    advertiser_id = str(profile.get("advertiser_id", ""))
    if not re.fullmatch(r"[0-9]{5,30}", advertiser_id):
        raise ValueError("advertiser_id must contain digits only")
    report_type, level = query.get("report_type"), query.get("data_level")
    if report_type != "BASIC" or level not in ALLOWED_LEVELS:
        raise ValueError("only BASIC reporting at an admitted auction level is supported")
    identifier = re.compile(r"[a-z][a-z0-9_]{0,79}")
    dimensions, metrics = query.get("dimensions"), query.get("metrics")
    for name, values, approved, maximum in (("dimensions", dimensions, profile.get("approved_dimensions"), 10), ("metrics", metrics, profile.get("approved_metrics"), 30)):
        if not isinstance(values, list) or not values or len(values) > maximum or len(values) != len(set(values)) or any(not isinstance(value, str) or not identifier.fullmatch(value) for value in values):
            raise ValueError(f"{name} must be a bounded unique identifier list")
        if not isinstance(approved, list) or not set(values).issubset(set(approved)):
            raise PermissionError(f"all {name} must be explicitly approved")
    start, end = date.fromisoformat(str(query.get("start_date", ""))), date.fromisoformat(str(query.get("end_date", "")))
    if end < start or (end - start).days > 31:
        raise ValueError("report range must be ordered and at most 31 days")
    page_size = query.get("page_size")
    if not isinstance(page_size, int) or isinstance(page_size, bool) or not 1 <= page_size <= 1000:
        raise ValueError("page_size must be within 1..1000")
    return advertiser_id, {"report_type": report_type, "data_level": level, "dimensions": dimensions, "metrics": metrics, "start_date": start.isoformat(), "end_date": end.isoformat(), "page_size": page_size}


def _response_to_dict(response: Any) -> dict[str, Any]:
    if isinstance(response, dict):
        value = response
    else:
        exporter = getattr(response, "to_dict", None)
        if not callable(exporter):
            raise RuntimeError("TikTok response is not an official SDK response")
        value = exporter()
    if not isinstance(value, dict):
        raise RuntimeError("TikTok response export must be an object")
    return value


def query_to_evidence(service: Any, access_token: str, advertiser_id: str, query: dict[str, Any], output_directory: Path) -> dict[str, Any]:
    if not access_token:
        raise PermissionError("TikTok access token is unavailable")
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".tiktok-report-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        pages, row_count, page = [], 0, 1
        while True:
            response = service.report_integrated_get(query["report_type"], access_token, advertiser_id=advertiser_id, data_level=query["data_level"], dimensions=query["dimensions"], metrics=query["metrics"], start_date=query["start_date"], end_date=query["end_date"], page=page, page_size=query["page_size"])
            payload = _response_to_dict(response)
            if payload.get("code") not in (0, None):
                raise RuntimeError(f"TikTok provider rejected report with code {payload.get('code')}")
            pages.append(payload)
            data = payload.get("data") if isinstance(payload.get("data"), dict) else {}
            rows = data.get("list") if isinstance(data.get("list"), list) else []
            row_count += len(rows)
            page_info = data.get("page_info") if isinstance(data.get("page_info"), dict) else {}
            total_page = page_info.get("total_page", page)
            if not isinstance(total_page, int) or total_page < page or total_page > 100:
                raise RuntimeError("TikTok page_info is invalid or exceeds 100-page safety bound")
            if page >= total_page:
                break
            page += 1
        response_bytes = (json.dumps({"sdk_distribution": SDK_DISTRIBUTION, "sdk_version": SDK_VERSION, "pages": pages}, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "provider-response.json").write_bytes(response_bytes)
        query_bytes = json.dumps(query, sort_keys=True, separators=(",", ":")).encode("utf-8")
        receipt = {"schema": "elite-tiktok-ads-reporting-receipt/v1", "created_at": datetime.now(timezone.utc).isoformat(), "provider": "TikTok Business API", "sdk_distribution": SDK_DISTRIBUTION, "sdk_version": SDK_VERSION, "advertiser_id_sha256": sha256(advertiser_id.encode("ascii")), "query_sha256": sha256(query_bytes), "page_count": len(pages), "row_count": row_count, "provider_response_sha256": sha256(response_bytes), "automatic_business_write": False}
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
    if importlib.metadata.version(SDK_DISTRIBUTION) != SDK_VERSION:
        raise RuntimeError(f"{SDK_DISTRIBUTION} {SDK_VERSION} is required")
    from business_api_client.api.reporting_api import ReportingApi
    profile, query_document = _read_object(args.profile), _read_object(args.query)
    advertiser_id, query = validate_configuration(profile, query_document)
    secret_name = str(profile.get("access_token_environment_variable", ""))
    if not re.fullmatch(r"[A-Z][A-Z0-9_]{2,80}", secret_name):
        raise ValueError("access token reference must be a canonical environment variable name")
    receipt = query_to_evidence(ReportingApi(), os.environ.get(secret_name, ""), advertiser_id, query, args.output)
    print(json.dumps({"status": "TIKTOK_ADS_REPORTING_PASS", "receipt": receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `tiktok_ads_reporting/test_reporting.py`
```yaml
block_id: "PY-TIKTOK-ADS-REPORTING:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed and official contract tests"
license: "LicenseRef-Workspace-Owner"
sha256: "39423c3d26ade4788b7d32b6478d268bae7570a61ea2a9122a8dc01bba75f9ed"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import json
from pathlib import Path
import tempfile
import unittest

from run_reporting import query_to_evidence, validate_configuration


def proven_profile():
    return {"provider": "TikTok Business API", "sdk_distribution": "tiktok-business-api-sdk-official", "sdk_version": "1.1.3", "decision": "PROVEN", "business_account_proven": True, "developer_app_proven": True, "reporting_scope_proven": True, "test_advertiser_contract_proven": True, "quota_and_cost_approved": True, "reconciliation_approved": True, "automatic_business_write": False, "data_retention": "30 days", "advertiser_id": "123456789", "approved_dimensions": ["campaign_id", "stat_time_day"], "approved_metrics": ["spend", "impressions", "clicks"]}


def valid_query():
    return {"report_type": "BASIC", "data_level": "AUCTION_CAMPAIGN", "dimensions": ["campaign_id", "stat_time_day"], "metrics": ["spend", "impressions"], "start_date": "2026-08-01", "end_date": "2026-08-02", "page_size": 1000}


class FakeService:
    def __init__(self, responses=None, error=None): self.responses = list(responses or []); self.error = error; self.calls = []
    def report_integrated_get(self, report_type, access_token, **kwargs):
        self.calls.append((report_type, access_token, kwargs))
        if self.error: raise self.error
        return self.responses.pop(0)


class TikTokReportingTests(unittest.TestCase):
    def test_configuration_fails_closed(self):
        profile = proven_profile(); profile["decision"] = "BLOCKED_ACCESS_AND_QUERY_APPROVAL_REQUIRED"
        with self.assertRaises(PermissionError): validate_configuration(profile, valid_query())
        profile = proven_profile(); query = valid_query(); query["metrics"] = ["unapproved"]
        with self.assertRaises(PermissionError): validate_configuration(profile, query)

    def test_configuration_rejects_bad_range_and_identifier(self):
        profile = proven_profile(); profile["advertiser_id"] = "adv_123"
        with self.assertRaises(ValueError): validate_configuration(profile, valid_query())
        profile = proven_profile(); query = valid_query(); query["end_date"] = "2026-10-01"
        with self.assertRaises(ValueError): validate_configuration(profile, query)

    def test_paginates_preserves_rows_and_redacts_advertiser(self):
        service = FakeService([{"code": 0, "data": {"list": [{"metrics": {"spend": "1.0"}}], "page_info": {"total_page": 2}}}, {"code": 0, "data": {"list": [{"metrics": {"spend": "2.0"}}], "page_info": {"total_page": 2}}}])
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = query_to_evidence(service, "secret-token", "123456789", valid_query(), output)
            self.assertEqual(receipt["page_count"], 2); self.assertEqual(receipt["row_count"], 2)
            self.assertNotIn("123456789", json.dumps(receipt)); self.assertNotIn("secret-token", json.dumps(receipt))
            self.assertEqual([call[2]["page"] for call in service.calls], [1, 2])

    def test_provider_error_is_atomic(self):
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            with self.assertRaisesRegex(RuntimeError, "provider failed"):
                query_to_evidence(FakeService(error=RuntimeError("provider failed")), "token", "123456789", valid_query(), output)
            self.assertFalse(output.exists()); self.assertEqual(list(Path(temp).iterdir()), [])

    def test_provider_code_and_page_bound_fail_closed(self):
        with tempfile.TemporaryDirectory() as temp:
            with self.assertRaisesRegex(RuntimeError, "code 40001"):
                query_to_evidence(FakeService([{"code": 40001, "message": "denied"}]), "token", "123456789", valid_query(), Path(temp) / "one")
            with self.assertRaisesRegex(RuntimeError, "100-page"):
                query_to_evidence(FakeService([{"code": 0, "data": {"page_info": {"total_page": 101}}}]), "token", "123456789", valid_query(), Path(temp) / "two")

    def test_official_sdk_contract(self):
        from business_api_client.api.reporting_api import ReportingApi
        self.assertTrue(callable(getattr(ReportingApi, "report_integrated_get", None)))
        source = Path(__import__("business_api_client.api.reporting_api", fromlist=["x"]).__file__).read_text(encoding="utf-8")
        self.assertIn("'/open_api/v1.3/report/integrated/get/', 'GET'", source)
        self.assertIn("header_params['Access-Token']", source)


if __name__ == "__main__": unittest.main()
````

### FILE: `tiktok_ads_reporting/README.md`
```yaml
block_id: "PY-TIKTOK-ADS-REPORTING:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local operator contract"
license: "LicenseRef-Workspace-Owner"
sha256: "469f2a2ec7be92db0d152baec326ddcc10649b7958edab3ec18e29a414aab59a"
variables: []
secrets_allowed: false
```
````markdown
# TikTok Ads Reporting Adapter

Read-only wrapper over TikTok's official PyPI distribution `tiktok-business-api-sdk-official==1.1.3`. TikTok's provider-owned GitHub README recommends this distribution and PyPI identifies its author as TikTok Pte. Ltd. The exact universal wheel and four runtime dependencies are fixed by URL and SHA-256; no community substitute or locally invented SDK code is used.

The profile starts blocked. Prove Business/developer app access, reporting permission, test advertiser, exact dimensions/metrics/date semantics, retention, quota/cost and reconciliation. The access token belongs only in its named environment variable. This pack calls only official `ReportingApi.report_integrated_get`, which issues `GET /open_api/v1.3/report/integrated/get/`, preserves every page and emits a hash-linked receipt without advertiser ID or token.

```powershell
python -m venv .venv
.venv\Scripts\python -m pip install --require-hashes --no-deps -r .\tiktok_ads_reporting\requirements.lock
.venv\Scripts\python -m pip check
Push-Location .\tiktok_ads_reporting
..\.venv\Scripts\python -m unittest -v test_reporting.py
Pop-Location
```

Fixtures prove package/API import, pagination, provider errors, atomic evidence and official SDK endpoint/header contract. They do not prove account access, field semantics, provider totals, attribution, rate limits, cost or production reconciliation.
````

## 6. Configuration surface

No contiene secretos. El proyecto debe completar el profile con decisión `PROVEN`, evidencia de business/developer app/reporting/test advertiser, advertiser numérico, dimensions/metrics permitidas, retención, cuota/costo y reconciliación. El token se resuelve sólo desde la variable de entorno indicada. El SDK se instala únicamente desde el lock exacto materializado.

## 7. Dependency bill

| Dependencia | Pin verificado | Licencia | Uso |
|---|---|---|---|
| `tiktok-business-api-sdk-official` | 1.1.3 wheel SHA `663b4a…f33b7` | MIT | SDK oficial TikTok |
| `urllib3` | 2.7.0 | MIT | transporte oficial SDK |
| `python-dateutil` | 2.9.0.post0 | Apache-2.0/BSD | fechas oficial SDK |
| `certifi` / `six` | 2026.7.22 / 1.17.0 | MPL-2.0 / MIT | runtime transitivo |

El lock de cinco wheels universales se verificó en CPython 3.14 Windows x86-64 con `--require-hashes --no-deps` y `pip check`.

## 8. Apply order

1. Materializar este pack en destino vacío y verificar hashes.
2. Crear un entorno aislado e instalar los cinco wheels exactos con hash enforcement.
3. Verificar versión, import y método/endpoint oficial.
4. Completar acceso/terms/costo/query sin secretos embebidos.
5. Ejecutar los seis tests contra el wheel oficial instalado.
6. Comparar una query mínima contra export/UI y aprobar reconciliación antes de operar read-only.

## 9. Verification

```powershell
python -m venv .venv
.venv\Scripts\python -m pip install --require-hashes --no-deps -r .\tiktok_ads_reporting\requirements.lock
.venv\Scripts\python -m pip check
Push-Location .\tiktok_ads_reporting
..\.venv\Scripts\python -m unittest -v test_reporting.py
Pop-Location
```

Resultado local verificado: instalación hash-frozen de cinco wheels, `pip check`, package/ReportingApi import, seis tests y contrato oficial endpoint/header PASS. Cuenta real, métricas, atribución, cuota/costo y reconciliación siguen `PROJECT_CONDITIONED`.

## 10. Reconstruction evidence

- `reconstruction_evidence/TIKTOK_ADS_REPORTING_ADAPTER_2026-08-26_V1.md` documenta el candidato source posteriormente rechazado por import graph roto.
- PyPI 1.1.3: wheel/autor/licencia/fecha/URL/SHA-256 fijados; cinco dependencias y seis tests verificados.
