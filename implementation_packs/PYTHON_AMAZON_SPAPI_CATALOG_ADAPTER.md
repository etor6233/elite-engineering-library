# Python Amazon SP-API Catalog Adapter

## 1. Metadata

```yaml
pack_id: "PYTHON-AMAZON-SPAPI-CATALOG-ADAPTER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un adapter read-only sobre el SDK oficial Amazon amzn-sp-api 1.11.1 y Catalog Items 2022-04-01 para consultar un ASIN aprobado, preservar la respuesta y emitir evidencia atómica sin escribir en el negocio."
stacks: ["Python 3.10-3.14", "amzn-sp-api 1.11.1", "Amazon Catalog Items 2022-04-01"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x"]
incompatible_with: ["Orders/PII/RDT", "Listings/Feeds writes", "secreto en CLI/Markdown", "ASIN no aprobado", "write automático", "rate limiter inferido"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0 AND MIT AND MPL-2.0 AND BSD-3-Clause"
upstream_sources: ["https://github.com/amzn/selling-partner-api-sdk/tree/8e792ae345a8d334ccdbdd03181f05f040e6a4fc", "https://pypi.org/project/amzn-sp-api/1.11.1/", "https://developer-docs.amazon.com/sp-api/docs/onboarding-overview", "https://developer-docs.amazon.com/sp-api/reference/getcatalogitem"]
verified_at: "2026-08-26"
```

## 2. Applicability

Use sólo cuando el usuario elija Amazon Selling Partner API, complete registro de developer y aplicación, roles, región/marketplace, ASINs, datasets, sandbox, cuota/costo y reconciliación. Esta lane cubre una consulta `get_catalog_item`; no cubre Orders, Listings, Feeds, Reports, restricted data tokens, paginación, rate limiting, retries ni sincronización productiva.

El SDK/wheel/modelos son código oficial de Amazon bajo Apache-2.0. El runner, los perfiles y los receipts son glue `AUTHORED` local y nunca se atribuyen a Amazon.

## 3. Architecture contract

El runner instala un grafo de siete wheels por URL y SHA-256 exactos, valida versión de distribución, carga secretos sólo desde variables de entorno canónicas y exige un profile `APPROVED` con accesos probados. Acepta un ASIN y marketplace canónicos, un subset de datasets oficiales y una sola llamada read-only. Serializa el modelo oficial, redacta identificadores del receipt y realiza commit atómico a un destino nuevo.

La plantilla nace bloqueada. Error de profile, secreto, versión, input, SDK/provider, serialización u output elimina staging. No existe fallback a SDK comunitario. El SDK Python oficial declara soporte básico pero no RDT ni rate limiter; este pack no inventa esas capacidades. Un sandbox real, headers/rate policy, retry, reconciliation, security y carga siguen siendo evidencia del target.

## 4. Exact file manifest

```text
CREATE amazon_spapi_catalog/requirements-direct.in
CREATE amazon_spapi_catalog/requirements.lock
CREATE amazon_spapi_catalog/sdk-artifact.lock.json
CREATE amazon_spapi_catalog/provider-profile.template.json
CREATE amazon_spapi_catalog/run_catalog_query.py
CREATE amazon_spapi_catalog/test_catalog_query.py
CREATE amazon_spapi_catalog/README.md
```

## 5. Materialization blocks

### FILE: `amazon_spapi_catalog/requirements-direct.in`
```yaml
block_id: "PY-AMAZON-SPAPI-CATALOG:requirements-direct:v1"
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

### FILE: `amazon_spapi_catalog/requirements.lock`
```yaml
block_id: "PY-AMAZON-SPAPI-CATALOG:requirements-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "exact official PyPI wheels selected from package metadata"
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

### FILE: `amazon_spapi_catalog/sdk-artifact.lock.json`
```yaml
block_id: "PY-AMAZON-SPAPI-CATALOG:sdk-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "local artifact receipt derived from official GitHub and PyPI metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "4cdddb62491f5ae44593e6025b8a31eefc5f8c17a093ad35c5d7461f28717901"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-python-sdk-artifact-lock/v1",
  "provider": "Amazon",
  "distribution": "amzn-sp-api",
  "version": "1.11.1",
  "sdk_requires_python": ">=3.9",
  "locked_graph_requires_python": ">=3.10",
  "wheel": {
    "filename": "amzn_sp_api-1.11.1-py3-none-any.whl",
    "bytes": 3557108,
    "sha256": "a88e4059c9fbfd454bfef67b82288a926f73cbdc608a8831f8383565736c70b0",
    "url": "https://files.pythonhosted.org/packages/36/3e/8744071121b96a7738fa485e584670cd4e6735b545cecc18516d10bd8881/amzn_sp_api-1.11.1-py3-none-any.whl"
  },
  "dependency_wheels": [
    {"distribution": "certifi", "version": "2026.7.22", "bytes": 136983, "sha256": "62f22742b58a1a33014a2b6b706588a8d7e2a88ae7bd1a6ebe8c992928483775"},
    {"distribution": "charset-normalizer", "version": "3.5.1", "bytes": 68658, "sha256": "6df0ec430f9a831772c23ca5a224cba36517a58a84bb32c32bb59a9fa67c47f6"},
    {"distribution": "idna", "version": "3.19", "bytes": 68550, "sha256": "815e7be7a7806d54abb586dc943addc79e8b2ee16915059658cbeff4b1b43bf4"},
    {"distribution": "requests", "version": "2.34.2", "bytes": 73075, "sha256": "2a0d60c172f83ac6ab31e4554906c0f3b3588d37b5cb939b1c061f4907e278e0"},
    {"distribution": "six", "version": "1.17.0", "bytes": 11050, "sha256": "4721f391ed90541fddacab5acf947aa0d3dc7d27b2e1e8eda2be8970586c3274"},
    {"distribution": "urllib3", "version": "2.7.0", "bytes": 131087, "sha256": "9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897"}
  ],
  "source": {
    "repository": "amzn/selling-partner-api-sdk",
    "release": "python-1.11.1",
    "commit": "8e792ae345a8d334ccdbdd03181f05f040e6a4fc",
    "archive_bytes": 26447333,
    "archive_sha256": "5bf712836bda76619c796fdb2b3a0a13acb281ba8bcf1b39c8f4b92f8f4799a1"
  },
  "license_expression": "Apache-2.0",
  "verified_at": "2026-08-26"
}
````

### FILE: `amazon_spapi_catalog/provider-profile.template.json`
```yaml
block_id: "PY-AMAZON-SPAPI-CATALOG:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed Amazon SP-API access/scope profile"
license: "LicenseRef-Workspace-Owner"
sha256: "5617acc2affba2d0a7559bc592f6cb82283515c141881609c7dfa166867efb90"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-amazon-spapi-catalog-profile/v1",
  "provider": "Amazon Selling Partner API",
  "sdk": "amzn-sp-api==1.11.1",
  "api": "Catalog Items 2022-04-01",
  "decision": "BLOCKED_ACCESS_AND_SCOPE_APPROVAL_REQUIRED",
  "region": "SANDBOX",
  "marketplace_id": "",
  "approved_asins": [],
  "included_data_allowlist": ["summaries"],
  "developer_registration": "NOT_PROVEN",
  "application_registration": "NOT_PROVEN",
  "roles_and_data_access": "NOT_PROVEN",
  "client_id_env": "AMAZON_SP_API_CLIENT_ID",
  "client_secret_env": "AMAZON_SP_API_CLIENT_SECRET",
  "refresh_token_env": "AMAZON_SP_API_REFRESH_TOKEN",
  "sandbox_contract_test": "NOT_PROVEN",
  "quota_and_cost_approved": false,
  "reconciliation_owner": "",
  "automatic_business_write": false
}
````

### FILE: `amazon_spapi_catalog/run_catalog_query.py`
```yaml
block_id: "PY-AMAZON-SPAPI-CATALOG:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local wrapper invoking official Amazon SP-API CatalogApi"
license: "LicenseRef-Workspace-Owner"
sha256: "53761b913d78ec5b5298bf1b6eda1633f799ab146cf3ae22724957df0bf57eba"
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

from spapi import CatalogApi, SPAPIClient, SPAPIConfig


SDK_VERSION = "1.11.1"
CATALOG_API_VERSION = "2022-04-01"
REGIONS = {"SANDBOX", "NA", "EU", "FE"}
INCLUDED_DATA = {
    "attributes",
    "classifications",
    "dimensions",
    "identifiers",
    "images",
    "productTypes",
    "relationships",
    "salesRanks",
    "summaries",
    "vendorDetails",
}


class CatalogService(Protocol):
    def get_catalog_item(self, asin: str, marketplace_ids: list[str], **kwargs: Any) -> Any: ...


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
    if profile.get("schema") != "elite-amazon-spapi-catalog-profile/v1":
        raise ValueError("unsupported profile schema")
    if profile.get("sdk") != f"amzn-sp-api=={SDK_VERSION}":
        raise ValueError("profile SDK version mismatch")
    if profile.get("decision") not in {"APPROVED_FOR_SANDBOX", "APPROVED_FOR_PRODUCTION"}:
        raise PermissionError("Amazon SP-API access decision is not approved")
    if profile.get("region") not in REGIONS:
        raise ValueError("unsupported Amazon SP-API region")
    if profile.get("developer_registration") != "PROVEN" or profile.get("application_registration") != "PROVEN":
        raise PermissionError("developer and application registration must be proven")
    if profile.get("roles_and_data_access") != "PROVEN":
        raise PermissionError("roles and data access must be proven")
    if profile.get("sandbox_contract_test") != "PROVEN":
        raise PermissionError("sandbox contract test must be proven")
    if profile.get("quota_and_cost_approved") is not True:
        raise PermissionError("quota and cost approval is required")
    if profile.get("automatic_business_write") is not False:
        raise ValueError("automatic business writes are forbidden")
    for key, expected in {
        "client_id_env": "AMAZON_SP_API_CLIENT_ID",
        "client_secret_env": "AMAZON_SP_API_CLIENT_SECRET",
        "refresh_token_env": "AMAZON_SP_API_REFRESH_TOKEN",
    }.items():
        if profile.get(key) != expected:
            raise ValueError(f"{key} must use the canonical secret environment reference")
    return profile


def query_catalog_item(
    service: CatalogService,
    asin: str,
    marketplace_id: str,
    included_data: list[str],
    output_directory: Path,
    timeout_seconds: float = 30.0,
) -> dict[str, Any]:
    if importlib.metadata.version("amzn-sp-api") != SDK_VERSION:
        raise RuntimeError(f"amzn-sp-api {SDK_VERSION} is required")
    if not re.fullmatch(r"[A-Z0-9]{10}", asin):
        raise ValueError("ASIN must be exactly 10 uppercase alphanumeric characters")
    if not re.fullmatch(r"[A-Z0-9]{8,20}", marketplace_id):
        raise ValueError("marketplace_id must be 8..20 uppercase alphanumeric characters")
    selected = list(dict.fromkeys(included_data))
    if not selected or len(selected) > len(INCLUDED_DATA) or any(value not in INCLUDED_DATA for value in selected):
        raise ValueError("included_data must be a non-empty subset of the official Catalog Items datasets")
    if timeout_seconds <= 0 or timeout_seconds > 120:
        raise ValueError("timeout must be within 1..120 seconds")
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".amazon-spapi-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        response = service.get_catalog_item(
            asin,
            [marketplace_id],
            included_data=selected,
            _request_timeout=timeout_seconds,
        )
        provider_payload = json_value(response)
        response_bytes = (json.dumps(provider_payload, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "provider-response.json").write_bytes(response_bytes)
        receipt = {
            "schema": "elite-amazon-spapi-catalog-receipt/v1",
            "created_at": datetime.now(timezone.utc).isoformat(),
            "provider": "Amazon Selling Partner API",
            "sdk_version": SDK_VERSION,
            "api_version": CATALOG_API_VERSION,
            "operation": "get_catalog_item",
            "asin_sha256": sha256(asin.encode("ascii")),
            "marketplace_id_sha256": sha256(marketplace_id.encode("ascii")),
            "included_data": selected,
            "provider_response_sha256": sha256(response_bytes),
            "automatic_business_write": False,
        }
        (stage / "QUERY_RECEIPT.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        stage.replace(output)
        return receipt
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def service_from_environment(profile: dict[str, Any]) -> CatalogApi:
    values: dict[str, str] = {}
    for key in ("client_id_env", "client_secret_env", "refresh_token_env"):
        env_name = profile[key]
        value = os.environ.get(env_name, "")
        if not value:
            raise PermissionError(f"required secret environment variable is missing: {env_name}")
        values[key] = value
    config = SPAPIConfig(
        client_id=values["client_id_env"],
        client_secret=values["client_secret_env"],
        refresh_token=values["refresh_token_env"],
        region=profile["region"],
    )
    client = SPAPIClient(config)
    return CatalogApi(client.api_client)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--asin", required=True)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--timeout-seconds", type=float, default=30.0)
    args = parser.parse_args()
    profile = read_approved_profile(args.profile)
    if args.asin not in profile.get("approved_asins", []):
        raise PermissionError("ASIN is not approved by the profile")
    service = service_from_environment(profile)
    receipt = query_catalog_item(
        service,
        args.asin,
        profile["marketplace_id"],
        profile["included_data_allowlist"],
        args.output,
        args.timeout_seconds,
    )
    print(json.dumps({"status": "AMAZON_SPAPI_CATALOG_PASS", "receipt": receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `amazon_spapi_catalog/test_catalog_query.py`
```yaml
block_id: "PY-AMAZON-SPAPI-CATALOG:test:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract/negative tests plus official SDK signature probe"
license: "LicenseRef-Workspace-Owner"
sha256: "a42fb98df9a113dbdcba5959e6b1967137311343bc6e2c34857359eda89533e4"
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

from spapi import CatalogApi, SPAPIConfig

from run_catalog_query import query_catalog_item, read_approved_profile


class FakeItem:
    def to_dict(self):
        return {"asin": "B012345678", "summaries": [{"marketplaceId": "ATVPDKIKX0DER", "itemName": "Test"}]}


class FakeCatalog:
    def __init__(self, response=None, error=None):
        self.response = response
        self.error = error
        self.calls = []

    def get_catalog_item(self, asin, marketplace_ids, **kwargs):
        self.calls.append((asin, marketplace_ids, kwargs))
        if self.error:
            raise self.error
        return self.response


class CatalogAdapterTests(unittest.TestCase):
    def test_preserves_official_response_and_redacts_identifiers(self):
        service = FakeCatalog(FakeItem())
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = query_catalog_item(service, "B012345678", "ATVPDKIKX0DER", ["summaries"], output)
            self.assertTrue((output / "provider-response.json").is_file())
            payload = json.loads((output / "provider-response.json").read_text(encoding="utf-8"))
            self.assertEqual(payload["asin"], "B012345678")
            self.assertNotIn("B012345678", json.dumps(receipt))
            self.assertNotIn("ATVPDKIKX0DER", json.dumps(receipt))
            self.assertFalse(receipt["automatic_business_write"])
            self.assertEqual(service.calls[0][2]["_request_timeout"], 30.0)

    def test_rejects_invalid_requests_and_is_atomic(self):
        with tempfile.TemporaryDirectory() as temp:
            root = Path(temp)
            cases = [
                ("bad", "ATVPDKIKX0DER", ["summaries"]),
                ("B012345678", "bad-id", ["summaries"]),
                ("B012345678", "ATVPDKIKX0DER", []),
                ("B012345678", "ATVPDKIKX0DER", ["orders"]),
            ]
            for index, (asin, marketplace, datasets) in enumerate(cases):
                output = root / f"out-{index}"
                with self.assertRaises(ValueError):
                    query_catalog_item(FakeCatalog(FakeItem()), asin, marketplace, datasets, output)
                self.assertFalse(output.exists())
            occupied = root / "occupied"
            occupied.mkdir()
            with self.assertRaises(FileExistsError):
                query_catalog_item(FakeCatalog(FakeItem()), "B012345678", "ATVPDKIKX0DER", ["summaries"], occupied)

    def test_provider_error_leaves_no_partial_output(self):
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "out"
            with self.assertRaisesRegex(RuntimeError, "provider failed"):
                query_catalog_item(FakeCatalog(error=RuntimeError("provider failed")), "B012345678", "ATVPDKIKX0DER", ["summaries"], output)
            self.assertFalse(output.exists())
            self.assertEqual(list(Path(temp).iterdir()), [])

    def test_profile_blocks_until_access_is_proven(self):
        with tempfile.TemporaryDirectory() as temp:
            path = Path(temp) / "profile.json"
            path.write_text(json.dumps({
                "schema": "elite-amazon-spapi-catalog-profile/v1",
                "sdk": "amzn-sp-api==1.11.1",
                "decision": "BLOCKED_ACCESS_AND_SCOPE_APPROVAL_REQUIRED",
            }), encoding="utf-8")
            with self.assertRaises(PermissionError):
                read_approved_profile(path)

    def test_official_sdk_contract(self):
        catalog_signature = inspect.signature(CatalogApi.get_catalog_item)
        config_signature = inspect.signature(SPAPIConfig)
        self.assertIn("asin", catalog_signature.parameters)
        self.assertIn("marketplace_ids", catalog_signature.parameters)
        self.assertIn("refresh_token", config_signature.parameters)
        self.assertEqual(CatalogApi.api_models_module, "spapi.models.catalogitems_v2022_04_01")


if __name__ == "__main__":
    unittest.main()
````

### FILE: `amazon_spapi_catalog/README.md`
```yaml
block_id: "PY-AMAZON-SPAPI-CATALOG:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local bounded runbook with official authority boundaries"
license: "LicenseRef-Workspace-Owner"
sha256: "cd7ba25c5766989f07f4b1faf061a42289950e81f1efa663d6a6f9bd0cdecd5f"
variables: []
secrets_allowed: false
```
````markdown
# Amazon SP-API Catalog evidence adapter

This materialized module invokes Amazon's official `amzn-sp-api` Python SDK `1.11.1` and the official Catalog Items `2022-04-01` client. The SDK is Amazon code under Apache-2.0. `run_catalog_query.py` is local integration glue and is never represented as Amazon-authored code.

The shipped profile is deliberately blocked. Before execution, prove developer and application registration, least-privilege roles, marketplace and ASIN scope, sandbox contract behavior, quota/cost approval, reconciliation ownership, and secret injection through the three canonical environment variables. Do not put credentials in Markdown, JSON, command arguments, logs, or receipts.

Install the complete hash-pinned wheel graph in an isolated Python 3.10+ environment, copy and complete the profile outside source control, and run one approved ASIN lookup:

```text
python -m pip install --require-hashes -r requirements.lock
python run_catalog_query.py --profile C:\secure\amazon-profile.json --asin B012345678 --output C:\evidence\amazon-catalog-B012345678
```

The adapter preserves the provider payload plus an atomic receipt and performs no business write. It deliberately does not claim Orders, Listings, Feeds, Reports, restricted data tokens, pagination, rate-limit control, retries, reconciliation, or production readiness. Amazon's Python SDK reports basic API support but no built-in RDT or rate limiter; add neither by assumption.
````

## 6. Configuration surface

| Key | Required | Secret | Gate |
|---|---:|---:|---|
| `decision` | yes | no | `APPROVED_FOR_SANDBOX` o `APPROVED_FOR_PRODUCTION` |
| `region`, `marketplace_id`, `approved_asins` | yes | no | scope explícito y allowlisted |
| `included_data_allowlist` | yes | no | subset de datasets Catalog Items oficiales |
| developer/application/roles/sandbox | yes | no | todos `PROVEN` |
| `quota_and_cost_approved`, `reconciliation_owner` | yes | no | cerrados antes de llamar |
| `AMAZON_SP_API_CLIENT_ID` | yes | yes | sólo environment/runtime secret store |
| `AMAZON_SP_API_CLIENT_SECRET` | yes | yes | sólo environment/runtime secret store |
| `AMAZON_SP_API_REFRESH_TOKEN` | yes | yes | sólo environment/runtime secret store |

## 7. Dependency bill

| Dependency | Exact version/artifact | License | Purpose |
|---|---|---|---|
| Amazon `amzn-sp-api` | 1.11.1 wheel SHA `a88e4059…70b0` | Apache-2.0 | cliente/modelos oficiales |
| requests | 2.34.2 wheel SHA `2a0d60c1…78e0` | Apache-2.0 | transporte requerido por SDK |
| urllib3 | 2.7.0 wheel SHA `9fb4c81e…6897` | MIT | HTTP requerido |
| certifi | 2026.7.22 wheel SHA `62f22742…3775` | MPL-2.0 | CA bundle |
| charset-normalizer | 3.5.1 wheel SHA `6df0ec43…47f6` | MIT | dependency de requests |
| idna | 3.19 wheel SHA `815e7be7…bf4` | BSD-3-Clause | dependency de requests |
| six | 1.17.0 wheel SHA `4721f391…3274` | MIT | compatibility requerida por SDK |

## 8. Apply order

1. Materializar `OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x` y este pack.
2. Verificar lock/source/wheel/licencias y crear un entorno Python 3.10+ aislado.
3. Instalar `requirements.lock` con `--require-hashes`; ejecutar `pip check` y tests.
4. Completar el profile fuera del source tree y aportar secretos por referencias runtime.
5. Probar sandbox con un ASIN allowlisted; preservar payload/receipt y reconciliar manualmente.
6. Sólo después de evidencia real evaluar retry/rate/reconciliation o nuevas operaciones como packs separados.

Rollback: deshabilitar la integración/job, revocar/rotar credenciales, conservar evidencia y retirar el output nuevo; nunca borrar datos de negocio para compensar una lectura.

## 9. Verification

```text
python -m venv .venv
.venv/Scripts/python -m pip install --require-hashes -r amazon_spapi_catalog/requirements.lock
.venv/Scripts/python -m pip check
cd amazon_spapi_catalog
../.venv/Scripts/python -m unittest -v test_catalog_query.py
```

Esperado local verificado: siete wheels instalados sin dependencias rotas; `amzn-sp-api 1.11.1`; cinco tests PASS; firma `CatalogApi.get_catalog_item(asin, marketplace_ids, ...)` y `SPAPIConfig(refresh_token=...)` presentes. No se ejecutó una llamada Amazon ni se demostró una cuenta/sandbox.

Gates productivos: developer/app/roles, OAuth refresh-token lifecycle, sandbox/marketplace, contract fixture, cuota/rate/retry, reconciliación, observabilidad sin secretos, seguridad, carga y rollback deben estar `PROVEN`. Hasta entonces `admission: CONDITIONED`.

## 10. Reconstruction evidence

`reconstruction_evidence/AMAZON_SP_API_CATALOG_ADAPTER_2026-08-26_V1.md` registra tag/commit verificado, archive/wheel/dependency hashes, discrepancia upstream de metadata, instalación exacta, procesos detenidos, cinco tests y gates globales. Evidencia local no sustituye acceso Amazon ni producción.
