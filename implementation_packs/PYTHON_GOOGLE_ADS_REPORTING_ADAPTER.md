# Python Google Ads Reporting Adapter

## 1. Metadata

```yaml
pack_id: "PYTHON-GOOGLE-ADS-REPORTING-ADAPTER"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un adapter read-only sobre el cliente oficial Google Ads 31.4.0/API v25 para ejecutar una GAQL SELECT aprobada, preservar filas y emitir evidencia atómica sin exponer customer ID."
stacks: ["Python 3.9-3.14", "google-ads 31.4.0", "Google Ads API v25"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x"]
incompatible_with: ["mutación de campañas", "conversion upload", "secreto en CLI/Markdown", "query no allowlisted", "write automático"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://pypi.org/project/google-ads/31.4.0/", "https://github.com/googleads/google-ads-python/tree/b8eb80ae277920d56fef467617795a7360c492fb", "https://developers.google.com/google-ads/api/docs/query/overview", "https://developers.google.com/google-ads/api/docs/release-notes"]
verified_at: "2026-09-04"
```

## 2. Applicability

Use sólo para reporting Google Ads cuando el usuario elija ese provider y cierre developer token, OAuth, customer/login hierarchy, cuenta test, query allowlist, datos, cuota/costo y retención. No usar para crear/modificar campañas, subir conversiones, decidir atribución ni escribir resultados directamente en el dominio.

El adapter es `AUTHORED` glue sobre `GoogleAdsClient`, `GoogleAdsService.search_stream` y protobufs oficiales. El SDK/modelos son Google; el código local no se atribuye a Google.

## 3. Architecture contract

El runner verifica distribución 31.4.0, customer ID sólo dígitos, una GAQL `SELECT` sin `;`, tamaño/timeout y output nuevo. Carga configuración únicamente con `GoogleAdsClient.load_from_env(version="v25")`, ejecuta stream, serializa todas las filas protobuf, hash del customer/query/response y commit atómico. No imprime credenciales ni customer ID en receipt.

Errores de config, query, API, row/output eliminan staging. Un resultado vacío es evidencia válida, no error. Rate limit/retry del SDK no prueba reconciliación. Performance, campos, privacy, quota/cost y freshness se validan en target; rollback desactiva query/job y conserva evidence.

## 4. Exact file manifest

```text
CREATE google_ads_reporting/requirements-direct.in
CREATE google_ads_reporting/requirements-windows-py314.lock
CREATE google_ads_reporting/sdk-artifact.lock.json
CREATE google_ads_reporting/provider-profile.template.json
CREATE google_ads_reporting/run_reporting_query.py
CREATE google_ads_reporting/test_reporting_query.py
CREATE google_ads_reporting/README.md
```

## 5. Materialization blocks

### FILE: `google_ads_reporting/requirements-direct.in`
```yaml
block_id: "PY-GOOGLE-ADS-REPORTING:requirements:v2"
operation: CREATE
provenance: AUTHORED
source: "exact official PyPI wheel URL/SHA-256"
license: "LicenseRef-Workspace-Owner"
sha256: "bac7c0ff2099893655d6f9c4e871910b9f5dc2b6e91cdba2c8e849e51f23afda"
variables: []
secrets_allowed: false
```
````text
google-ads @ https://files.pythonhosted.org/packages/5d/59/27f19973bee9d6ae64e052d20767e2f586067b8237dd935254927f9009b5/google_ads-31.4.0-py3-none-any.whl#sha256=210a7f6a7b5d40ef544090216dceeb90c9d2e7fba51c72329a690c9e5bb94474
````

### FILE: `google_ads_reporting/requirements-windows-py314.lock`
```yaml
block_id: "PY-GOOGLE-ADS-REPORTING:requirements-windows-py314:v1"
operation: CREATE
provenance: AUTHORED
source: "exact official PyPI wheel graph resolved for CPython 3.14 Windows x86-64"
license: "LicenseRef-Workspace-Owner"
sha256: "b48ac7347fe788d8c2aa60ab03c7e24ed22179d70faa232a52085f0d31947f73"
variables: []
secrets_allowed: false
```
````text
# Exact wheel graph resolved and verified for CPython 3.14 on Windows x86-64.
certifi @ https://files.pythonhosted.org/packages/0b/a7/71ac2cff56fec219ed242bb11b8efb69fcc4bec75db06fb7bfe35de520e6/certifi-2026.7.22-py3-none-any.whl#sha256=62f22742b58a1a33014a2b6b706588a8d7e2a88ae7bd1a6ebe8c992928483775
cffi @ https://files.pythonhosted.org/packages/a7/06/1c3e01e3ba14c39f6d10bfbac52753b7e22259e38088e5cfe1d704918690/cffi-2.1.1-cp314-cp314-win_amd64.whl#sha256=3222ba5d678f80a030e6afbcc33dc1ae5cb45facabb61cee2c7016b8432fde48
charset-normalizer @ https://files.pythonhosted.org/packages/7a/7c/4938c329b6a9d446f6a59aa2092ff7118f274209b5ed0e26893d1d30a63c/charset_normalizer-3.5.1-cp314-cp314-win_amd64.whl#sha256=c658c50ac0c98cd755a2dd50b7977d3bca7df401dcc47fbdfa87db53ef7d4e8b
cryptography @ https://files.pythonhosted.org/packages/42/8b/cb12b1b60c91b074ca6bf0fdd59aa8f10d8bc5f73af8faece86ef0421b37/cryptography-50.0.1-cp311-abi3-win_amd64.whl#sha256=aed8db4f6d71c51efb89530e12d9464e7bf2923d46c3205dc794a2a93f8c0648
google-ads @ https://files.pythonhosted.org/packages/5d/59/27f19973bee9d6ae64e052d20767e2f586067b8237dd935254927f9009b5/google_ads-31.4.0-py3-none-any.whl#sha256=210a7f6a7b5d40ef544090216dceeb90c9d2e7fba51c72329a690c9e5bb94474
google-api-core @ https://files.pythonhosted.org/packages/44/56/30c91c61b8f70d4c09285a005b94798729aeaf4ec8b90c1c360da8207728/google_api_core-2.36.0-py3-none-any.whl#sha256=e4d0b179260727ea5c42222426d9199285214dbef7bf48f8b16600c7f9a78944
google-auth @ https://files.pythonhosted.org/packages/3b/27/0f7247b8002a1404fdb5412aeb15fb0fe70288eb8f88a5873d1c0d8262cf/google_auth-2.57.1-py3-none-any.whl#sha256=ab439dee60a6856412bc058f13a68eb6a59c5b81526bb89cfdcf89ce9c0a48c9
google-auth-oauthlib @ https://files.pythonhosted.org/packages/6f/32/ebf82c6ae8c1caf975b3153d763f3b765dddafa01d2a6c2a142fde28bda7/google_auth_oauthlib-1.4.1-py3-none-any.whl#sha256=a1be43ec69fe563ac9b2c4d6fc4334b323b21cbdc59a638b5fa34dd4d5a2a348
googleapis-common-protos @ https://files.pythonhosted.org/packages/1a/7a/7d79170c6ce6f12e109df2b3879d6b934010cf4f99aea8de8b7e5408c174/googleapis_common_protos-1.75.3-py3-none-any.whl#sha256=a018d2bf098ca9fb6faa08d5bb780e2a2c2f73c566f069761331386c9596d3f2
grpcio @ https://files.pythonhosted.org/packages/65/22/fc9a622d885a7a37ff972a12faaef443d74e47407181da70d0ab62ab41f0/grpcio-1.83.1-cp314-cp314-win_amd64.whl#sha256=47e6934ad38779271e2e7cc5f78a63a407cf3d98114c65c1fdbcd3f5a716f29b
grpcio-status @ https://files.pythonhosted.org/packages/e7/34/b03ec688e5c8c6ce283ec8b214d9c171e97059eb55157465f6fd6db1562c/grpcio_status-1.83.1-py3-none-any.whl#sha256=2cd328ee62ef2b3eb957dd3b75db7dadcdbb76488fdc9ab3aba1ebfbbdc324a4
idna @ https://files.pythonhosted.org/packages/57/b0/0e52c878c53f245edd3a11020f20979b3f490f245af532c7cae3027754b5/idna-3.19-py3-none-any.whl#sha256=815e7be7a7806d54abb586dc943addc79e8b2ee16915059658cbeff4b1b43bf4
oauthlib @ https://files.pythonhosted.org/packages/be/9c/92789c596b8df838baa98fa71844d84283302f7604ed565dafe5a6b5041a/oauthlib-3.3.1-py3-none-any.whl#sha256=88119c938d2b8fb88561af5f6ee0eec8cc8d552b7bb1f712743136eb7523b7a1
opentelemetry-api @ https://files.pythonhosted.org/packages/ca/6f/a04e900f465ff3221ccc395522503e2d10e79fa21f2723c8e177aae1e0d1/opentelemetry_api-1.44.0-py3-none-any.whl#sha256=94b98c893a91b88657eaac1e3ba89618cdb85be6918196705354f34728b2cdef
proto-plus @ https://files.pythonhosted.org/packages/41/5d/0f04b85dafdc3250ced7f2592efc17dce7f40712e941e9632202481e600d/proto_plus-1.28.4-py3-none-any.whl#sha256=4b01341272f8a348db3f003b6143109f83ab43091019d5181b3fcdf500ab32aa
protobuf @ https://files.pythonhosted.org/packages/db/37/155788a0d8daded960375af604202805308169f9b859419ea0aa370946e2/protobuf-7.36.1-cp310-abi3-win_amd64.whl#sha256=51139351435d9b43d88a55eaa49fb6f737fbb478fb0cbf2cf694d1a04a9d3363
pyasn1 @ https://files.pythonhosted.org/packages/9a/3b/6163796d69c3977d1e4287bea4a6979161cbbdd170ebb430511e8e1999ce/pyasn1-0.6.4-py3-none-any.whl#sha256=deda9277cfd454080ec40b207fb6df82206a3a2688735233cdcd8d3d565f088b
pyasn1-modules @ https://files.pythonhosted.org/packages/47/8d/d529b5d697919ba8c11ad626e835d4039be708a35b0d22de83a269a6682c/pyasn1_modules-0.4.2-py3-none-any.whl#sha256=29253a9207ce32b64c3ac6600edc75368f98473906e8fd1043bd6b5b1de2c14a
pycparser @ https://files.pythonhosted.org/packages/0c/c3/44f3fbbfa403ea2a7c779186dc20772604442dde72947e7d01069cbe98e3/pycparser-3.0-py3-none-any.whl#sha256=b727414169a36b7d524c1c3e31839a521725078d7b2ff038656844266160a992
PyYAML @ https://files.pythonhosted.org/packages/23/20/bb6982b26a40bb43951265ba29d4c246ef0ff59c9fdcdf0ed04e0687de4d/pyyaml-6.0.3-cp314-cp314-win_amd64.whl#sha256=4a2e8cebe2ff6ab7d1050ecd59c25d4c8bd7e6f400f5f82b96557ac0abafd0ac
requests @ https://files.pythonhosted.org/packages/a0/f4/c67b0b3f1b9245e8d266f0f112c500d50e5b4e83cb6f3b71b6528104182a/requests-2.34.2-py3-none-any.whl#sha256=2a0d60c172f83ac6ab31e4554906c0f3b3588d37b5cb939b1c061f4907e278e0
requests-oauthlib @ https://files.pythonhosted.org/packages/3b/5d/63d4ae3b9daea098d5d6f5da83984853c1bbacd5dc826764b249fe119d24/requests_oauthlib-2.0.0-py2.py3-none-any.whl#sha256=7dd8a5c40426b779b0868c404bdef9768deccf22749cde15852df527e6269b36
typing-extensions @ https://files.pythonhosted.org/packages/49/d3/b8441a820a491ddfc024b0b0cf0393375b75ea13866d9c66727e54c2fc80/typing_extensions-4.16.0-py3-none-any.whl#sha256=481caa481374e813c1b176ada14e97f1f67a4539ce9cfeb3f350d78d6370c2e8
urllib3 @ https://files.pythonhosted.org/packages/7f/3e/5db95bcf282c52709639744ca2a8b149baccf648e39c8cc87553df9eae0c/urllib3-2.7.0-py3-none-any.whl#sha256=9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897
````

### FILE: `google_ads_reporting/sdk-artifact.lock.json`
```yaml
block_id: "PY-GOOGLE-ADS-REPORTING:sdk-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "local receipt derived from official GitHub/PyPI metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "5024ccf18d86ea6ab2766779140eba9e17ba5eb08f14b067811687d0c40ca889"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-python-sdk-artifact-lock/v1",
  "provider": "Google",
  "distribution": "google-ads",
  "version": "31.4.0",
  "requires_python": "<3.15,>=3.9",
  "release_level": "Production/Stable",
  "wheel": {
    "filename": "google_ads-31.4.0-py3-none-any.whl",
    "bytes": 23712034,
    "sha256": "210a7f6a7b5d40ef544090216dceeb90c9d2e7fba51c72329a690c9e5bb94474",
    "url": "https://files.pythonhosted.org/packages/5d/59/27f19973bee9d6ae64e052d20767e2f586067b8237dd935254927f9009b5/google_ads-31.4.0-py3-none-any.whl"
  },
  "source": {
    "repository": "https://github.com/googleads/google-ads-python",
    "tag": "31.4.0",
    "commit": "b8eb80ae277920d56fef467617795a7360c492fb",
    "commit_signature_verified": true
  },
  "resolved_environment": {
    "python": "3.14.4",
    "platform": "Windows x86-64",
    "pip": "26.0.1",
    "dependency_count": 24,
    "resolved_at": "2026-09-04"
  },
  "api_version": "v25",
  "automatic_business_write": false
}
````

### FILE: `google_ads_reporting/provider-profile.template.json`
```yaml
block_id: "PY-GOOGLE-ADS-REPORTING:profile:v2"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed Google Ads access/query profile"
license: "LicenseRef-Workspace-Owner"
sha256: "9c637c170d661703ae354ff86c77997ec237c2244997b1f77eda2f4fef5a5e00"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-google-ads-reporting-profile/v1",
  "provider": "Google Ads API",
  "sdk": "google-ads==31.4.0",
  "api_version": "v25",
  "decision": "BLOCKED_ACCESS_AND_QUERY_APPROVAL_REQUIRED",
  "authentication": "GoogleAdsClient.load_from_env",
  "developer_token_access": "NOT_PROVIDED",
  "oauth_identity_reference": "",
  "customer_id": "",
  "login_customer_id": "",
  "query_allowlist": [],
  "data_retention": "",
  "cost_and_quota_approved": false,
  "automatic_business_write": false
}
````

### FILE: `google_ads_reporting/run_reporting_query.py`
```yaml
block_id: "PY-GOOGLE-ADS-REPORTING:runner:v2"
operation: CREATE
provenance: AUTHORED
source: "local wrapper invoking official Google Ads client API"
license: "LicenseRef-Workspace-Owner"
sha256: "80249919a57d8ceff3e8d97f705e94d9f8b4861ff9ba52ec59dcfb184f0d364a"
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
from pathlib import Path
import re
import shutil
import sys
import uuid
from typing import Any, Iterable, Protocol

from google.ads.googleads.client import GoogleAdsClient
from google.protobuf.json_format import MessageToDict


SDK_VERSION = "31.4.0"
API_VERSION = "v25"


class GoogleAdsService(Protocol):
    def search_stream(self, *, customer_id: str, query: str, timeout: float) -> Iterable[Any]: ...


def sha256(value: bytes) -> str:
    return hashlib.sha256(value).hexdigest()


def row_to_dict(row: Any) -> dict[str, Any]:
    if isinstance(row, dict):
        return row
    protobuf = getattr(row, "_pb", None)
    if protobuf is None:
        raise RuntimeError("Google Ads result row is not an official protobuf message")
    return MessageToDict(protobuf, preserving_proto_field_name=True)


def query_to_evidence(service: GoogleAdsService, customer_id: str, query: str, output_directory: Path, timeout_seconds: float = 60.0) -> dict[str, Any]:
    if importlib.metadata.version("google-ads") != SDK_VERSION:
        raise RuntimeError(f"google-ads {SDK_VERSION} is required")
    if not re.fullmatch(r"[0-9]{5,20}", customer_id):
        raise ValueError("customer_id must contain digits only")
    normalized = " ".join(query.split())
    if not normalized.upper().startswith("SELECT ") or len(normalized) > 10000 or ";" in normalized:
        raise ValueError("one bounded GAQL SELECT query is required")
    if timeout_seconds <= 0 or timeout_seconds > 300:
        raise ValueError("timeout must be within 1..300 seconds")
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".google-ads-report-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        rows: list[dict[str, Any]] = []
        for batch in service.search_stream(customer_id=customer_id, query=normalized, timeout=timeout_seconds):
            for row in getattr(batch, "results", []):
                rows.append(row_to_dict(row))
        response = {"api_version": API_VERSION, "rows": rows}
        response_bytes = (json.dumps(response, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "provider-response.json").write_bytes(response_bytes)
        receipt = {
            "schema": "elite-google-ads-reporting-receipt/v1",
            "created_at": datetime.now(timezone.utc).isoformat(),
            "provider": "Google Ads API",
            "sdk_version": SDK_VERSION,
            "api_version": API_VERSION,
            "customer_id_sha256": sha256(customer_id.encode("ascii")),
            "query_sha256": sha256(normalized.encode("utf-8")),
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
    parser.add_argument("--customer-id", required=True)
    parser.add_argument("--query-file", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--timeout-seconds", type=float, default=60.0)
    args = parser.parse_args()
    query = args.query_file.read_text(encoding="utf-8")
    client = GoogleAdsClient.load_from_env(version=API_VERSION)
    service = client.get_service("GoogleAdsService", version=API_VERSION)
    receipt = query_to_evidence(service, args.customer_id, query, args.output, args.timeout_seconds)
    print(json.dumps({"status": "GOOGLE_ADS_REPORTING_PASS", "receipt": receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `google_ads_reporting/test_reporting_query.py`
```yaml
block_id: "PY-GOOGLE-ADS-REPORTING:test:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract/negative tests plus official SDK signature probe"
license: "LicenseRef-Workspace-Owner"
sha256: "fe2b0dd9c72238990d67e7202d1dc6e7446468a8e5c2b2f77ba49339372c1671"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import json
from pathlib import Path
from types import SimpleNamespace
import tempfile
import unittest

from run_reporting_query import query_to_evidence


class FakeService:
    def __init__(self, batches): self.batches = batches; self.calls = []
    def search_stream(self, *, customer_id, query, timeout): self.calls.append((customer_id, query, timeout)); return self.batches


class ReportingTests(unittest.TestCase):
    def test_preserves_rows_and_redacts_customer(self):
        service = FakeService([SimpleNamespace(results=[{"campaign": {"id": "1"}, "metrics": {"clicks": "2"}}])])
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = query_to_evidence(service, "1234567890", "SELECT campaign.id, metrics.clicks FROM campaign", output)
            self.assertEqual(receipt["row_count"], 1)
            self.assertNotIn("1234567890", json.dumps(receipt))
            raw = json.loads((output / "provider-response.json").read_text(encoding="utf-8"))
            self.assertEqual(raw["rows"][0]["metrics"]["clicks"], "2")
            self.assertFalse(receipt["automatic_business_write"])

    def test_empty_report_is_valid_evidence(self):
        with tempfile.TemporaryDirectory() as temp:
            receipt = query_to_evidence(FakeService([]), "12345", "SELECT customer.id FROM customer", Path(temp) / "out")
            self.assertEqual(receipt["row_count"], 0)

    def test_fails_closed_and_atomic(self):
        invalid = ["DELETE campaign", "SELECT x FROM y; SELECT z", " ", "MUTATE x"]
        with tempfile.TemporaryDirectory() as temp:
            root = Path(temp)
            for index, query in enumerate(invalid):
                output = root / f"out-{index}"
                with self.assertRaises(ValueError): query_to_evidence(FakeService([]), "12345", query, output)
                self.assertFalse(output.exists())
            with self.assertRaises(ValueError): query_to_evidence(FakeService([]), "123-45", "SELECT x FROM y", root / "customer")
            occupied = root / "occupied"; occupied.mkdir()
            with self.assertRaises(FileExistsError): query_to_evidence(FakeService([]), "12345", "SELECT x FROM y", occupied)


if __name__ == "__main__":
    unittest.main()
````

### FILE: `google_ads_reporting/README.md`
```yaml
block_id: "PY-GOOGLE-ADS-REPORTING:readme:v2"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "2e90344dc27ebce3d0124ef669e43c491fd8ad5644c624acbfe76f7cb9ba898c"
variables: []
secrets_allowed: false
```
````markdown
# Google Ads reporting adapter

This read-only adapter uses the exact official `google-ads==31.4.0` client and API `v25`. It loads credentials and developer-token configuration only through `GoogleAdsClient.load_from_env`, executes one bounded GAQL `SELECT` with `GoogleAdsService.search_stream`, preserves every modeled result row, and writes an atomic response/receipt. The receipt hashes the customer ID and query instead of exposing them and never authorizes a business write.

It does not create or mutate campaigns, upload conversions, choose attribution policy, or claim that report data is reconciled with orders. Before a real query the project must prove developer token/API access, OAuth identity, manager/login customer hierarchy, test account, query allowlist, field data policy, quota/cost, retention, logging redaction, reconciliation and rollback.

```text
python -m unittest -v test_reporting_query.py
python run_reporting_query.py --customer-id <digits> --query-file <approved.gaql> --output <new-directory>
```

The direct wheel and the complete CPython 3.14/Windows x86-64 dependency graph are hash-pinned. A different interpreter or platform must resolve, lock, scan and prove its own graph before deployment.
````

## 6. Configuration surface

| Config | Default | Validación | Secreto | Efecto |
|---|---|---|---|---|
| Google Ads env config | none | SDK official env names, access proven | developer/OAuth sí | identity/account |
| customer ID | none | 5..20 digits | confidential; receipt hashed | dataset |
| API version | v25 fixed | SDK 31.4.0 must support | no | schema |
| query file | none | one SELECT, no semicolon, allowlist external | may expose schema | report |
| timeout | 60s | 1..300 | no | deadline |
| output | none | must not exist | may contain business data | evidence |

El profile inicia bloqueado y `automatic_business_write=false`.

## 7. Dependency bill

| Package/tool | Pin | Uso | Licencia | Fase | Fuente |
|---|---|---|---|---|---|
| `google-ads` | 31.4.0 wheel SHA `210a7f...4474` | client/services/protobuf | Apache-2.0 | runtime | PyPI/Google |
| transitive wheel graph | 24 exact URLs/SHA-256 | reproducible py314/Windows runtime | mixed permissive; scan required | runtime | PyPI |
| Python | 3.9..<3.15 | runner/tests | PSF-2.0 | runtime | python.org |

El wheel directo y el grafo py314/Windows están fijados; otra plataforma debe resolver, fijar y scanear su propio grafo antes de deploy.

## 8. Apply order

Materializar, adquirir wheel exacto, resolver lock transitivo, tests, completar profile/access, aprobar GAQL y ejecutar en cuenta test read-only. En existente, integrar como job de reporting aislado. Rollback deshabilita schedule/query, conserva response/receipt y revoca scopes si aplica.

## 9. Verification

```text
python -m unittest -v test_reporting_query.py
python -m pip install --require-hashes --no-deps -r requirements-windows-py314.lock
python -m pip check
python -c "import importlib.metadata as m; import google.ads.googleads.v25; assert m.version('google-ads') == '31.4.0'"
```

Éxito: 24 wheels hash-verified + pip check + 3 tests + SDK/API contract PASS; rows/empty/redaction/atomic/negativos. Pendiente: SCA vigente, developer token/OAuth/account test/query real, quota/cost/privacy/reconciliation/load.

## 10. Reconstruction evidence

- Python 3.14.4; release/tag oficial 31.4.0 en commit firmado `b8eb80ae277920d56fef467617795a7360c492fb`; wheel 23,712,034 bytes/SHA-256 `210a7f6a...4474`;
- 7/7 archivos; grafo Windows py314 de 24 wheels con URL/SHA-256 y `pip check` PASS;
- 3 tests + import/version/API v25 probe PASS;
- sin llamada Google ni credenciales;
- primer install combinado se interrumpió tras quedar sin salida; probes separados import/version pasaron (`LIB-FAIL-043`);
- 2026-09-04 / Codex; freshness y lock ejecutable V2; el live gate sigue condicionado.
