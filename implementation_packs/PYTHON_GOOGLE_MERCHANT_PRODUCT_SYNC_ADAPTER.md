# Python Google Merchant Product Sync Adapter

## 1. Metadata

```yaml
pack_id: "PYTHON-GOOGLE-MERCHANT-PRODUCT-SYNC-ADAPTER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un adapter sobre el cliente oficial GA Google Merchant Products 1.8.0/v1 para insertar un ProductInput aprobado, preservar evidencia durable y consultar Product status sin confundir inserción con aprobación."
stacks: ["Python 3.14 Windows x86-64 verified", "google-shopping-merchant-products 1.8.0", "Google Merchant Products v1"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x"]
incompatible_with: ["credencial en source/CLI", "data source no API primary", "offer no aprobado", "aprobación inferida", "write local automático", "lock Windows usado en otro ABI"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0 AND MIT AND MPL-2.0 AND BSD-3-Clause AND PSF-2.0"
upstream_sources: ["https://github.com/googleapis/google-cloud-python/tree/97d7b42cd74b41211f5ec8871cc0dd15debdb1a0/packages/google-shopping-merchant-products", "https://pypi.org/project/google-shopping-merchant-products/1.8.0/", "https://developers.google.com/merchant/api/guides/products/add-manage", "https://developers.google.com/merchant/api/client-libraries/python"]
verified_at: "2026-08-26"
```

## 2. Applicability

Use cuando el usuario elija Google Merchant API y cierre proyecto/billing, API, ADC, Merchant account, data source API primary, test account, offer IDs, escritura externa, cuota/costo, retención, refresh y reconciliación. Materializa un `productInputs.insert` seguido de `products.get`; no crea accounts/data sources, no elimina productos, no gestiona promotions/inventory/orders y no escribe el dominio local.

El SDK, modelos y servicios v1 son código oficial Google bajo Apache-2.0. El wrapper, perfiles, validaciones y receipts son glue `AUTHORED` local, no código atribuido a Google.

## 3. Architecture contract

La lane verificada instala veinte wheels exactos por URL/SHA-256 para CPython 3.14 Windows x86-64. Otros targets parten del wheel superior fijado y deben producir un lock candidato propio antes de ejecutar. El runner usa Application Default Credentials, nunca acepta tokens, construye mensajes oficiales `ProductInput`, `ProductAttributes`, `Price` e `InsertProductInputRequest`, deshabilita retry implícito del write y exige output nuevo.

Tras el insert, persiste request/response/receipt antes de consultar estado. `NotFound` acotado queda `INSERT_ACCEPTED_PRODUCT_NOT_YET_VISIBLE`; otro fallo de status conserva evidencia post-write y se propaga. Un `get` exitoso preserva `product_status`, pero `approval_inferred=false`: Google declara que insert sólo inicia procesamiento y no prueba aprobación. Ningún response entra al dominio automáticamente.

## 4. Exact file manifest

```text
CREATE google_merchant_product_sync/requirements-direct.in
CREATE google_merchant_product_sync/requirements-windows-py314.lock
CREATE google_merchant_product_sync/sdk-artifact.lock.json
CREATE google_merchant_product_sync/provider-profile.template.json
CREATE google_merchant_product_sync/product.template.json
CREATE google_merchant_product_sync/sync_product.py
CREATE google_merchant_product_sync/test_sync_product.py
CREATE google_merchant_product_sync/README.md
```

## 5. Materialization blocks

### FILE: `google_merchant_product_sync/requirements-direct.in`
```yaml
block_id: "PY-GOOGLE-MERCHANT-SYNC:requirements-direct:v1"
operation: CREATE
provenance: AUTHORED
source: "official Google wheel URL and SHA-256 from PyPI"
license: "LicenseRef-Workspace-Owner"
sha256: "b0d037fe32e656ee47a9a7f8ca36d90e5595008fe2e0abf55e6d4237083f9f5b"
variables: []
secrets_allowed: false
```
````text
google-shopping-merchant-products @ https://files.pythonhosted.org/packages/ea/7f/4de247c1900a4bb9a2e6bec51128c907545980133f1f39293c21f7c17527/google_shopping_merchant_products-1.8.0-py3-none-any.whl#sha256=722ef095eca35126255c129964a61818f0d7bc9eeef7f1da12732d81c1dcd505
````

### FILE: `google_merchant_product_sync/requirements-windows-py314.lock`
```yaml
block_id: "PY-GOOGLE-MERCHANT-SYNC:requirements-windows:v1"
operation: CREATE
provenance: AUTHORED
source: "exact official PyPI wheel graph resolved for CPython 3.14 Windows x86-64"
license: "LicenseRef-Workspace-Owner"
sha256: "4b9f2dec1970d8a3a3ee1caf6309407d70cbaf70468c7f350fc911b012f0cd25"
variables: []
secrets_allowed: false
```
````text
# Exact wheel graph resolved and verified for CPython 3.14 on Windows x86-64.
certifi @ https://files.pythonhosted.org/packages/0b/a7/71ac2cff56fec219ed242bb11b8efb69fcc4bec75db06fb7bfe35de520e6/certifi-2026.7.22-py3-none-any.whl#sha256=62f22742b58a1a33014a2b6b706588a8d7e2a88ae7bd1a6ebe8c992928483775
cffi @ https://files.pythonhosted.org/packages/a7/06/1c3e01e3ba14c39f6d10bfbac52753b7e22259e38088e5cfe1d704918690/cffi-2.1.1-cp314-cp314-win_amd64.whl#sha256=3222ba5d678f80a030e6afbcc33dc1ae5cb45facabb61cee2c7016b8432fde48
charset-normalizer @ https://files.pythonhosted.org/packages/7a/7c/4938c329b6a9d446f6a59aa2092ff7118f274209b5ed0e26893d1d30a63c/charset_normalizer-3.5.1-cp314-cp314-win_amd64.whl#sha256=c658c50ac0c98cd755a2dd50b7977d3bca7df401dcc47fbdfa87db53ef7d4e8b
cryptography @ https://files.pythonhosted.org/packages/42/8b/cb12b1b60c91b074ca6bf0fdd59aa8f10d8bc5f73af8faece86ef0421b37/cryptography-50.0.1-cp311-abi3-win_amd64.whl#sha256=aed8db4f6d71c51efb89530e12d9464e7bf2923d46c3205dc794a2a93f8c0648
google-api-core @ https://files.pythonhosted.org/packages/bc/c1/a8a92ae1bc4b1a8f804c776d7d3f0c771b78a62c3ad4df1be41b3fd8c767/google_api_core-2.34.0-py3-none-any.whl#sha256=cdf9c67e7ca2402d86ccbfde5f2503fc83e3cc3f58cc78456ae96cad24a6d2de
google-auth @ https://files.pythonhosted.org/packages/00/f3/8508a702c094af5f6e89773f4dfdeee74913df0f41a02c21b5e7dc3d75cd/google_auth-2.57.0-py3-none-any.whl#sha256=180dafe015cfb62193bea26b677500fab5b9fd51a1e825ebf3ad9b182047ae59
google-shopping-merchant-products @ https://files.pythonhosted.org/packages/ea/7f/4de247c1900a4bb9a2e6bec51128c907545980133f1f39293c21f7c17527/google_shopping_merchant_products-1.8.0-py3-none-any.whl#sha256=722ef095eca35126255c129964a61818f0d7bc9eeef7f1da12732d81c1dcd505
google-shopping-type @ https://files.pythonhosted.org/packages/37/b0/2b6287659916fd70564951e5d5797ca987da3cc67897fd7a1d1724075401/google_shopping_type-1.5.0-py3-none-any.whl#sha256=afc1180f2d068713bc1bee4b1c6882876023be6ef42e1af020991dfab2595072
googleapis-common-protos @ https://files.pythonhosted.org/packages/47/5b/1c9e55363c3b1890a98cae813de5b4ea327845756cd8fb7ee690140c7eac/googleapis_common_protos-1.75.2-py3-none-any.whl#sha256=6b83302f554ea93a0f48409c7fc2050f954bcbcddb7e3a9c76d4a823cb22920e
grpcio @ https://files.pythonhosted.org/packages/a1/00/b1b26431c9d54eee11724fd6e5585473a2ed47fbc1fb95e5204906a642ce/grpcio-1.83.0-cp314-cp314-win_amd64.whl#sha256=2bb48cb5e6dd005ca12b89ce4b6ac0b48ff3112c747542ee7986ef611a8ca6d9
grpcio-status @ https://files.pythonhosted.org/packages/d6/00/73204406228cf989bea6b0fd9fe4702fab49a8a152a0c6f90856dadb6ac7/grpcio_status-1.83.0-py3-none-any.whl#sha256=f6a838a7c5fb84ae98833ec0ef81ed438c26e11e54b2ddb8e92ad328c861de69
idna @ https://files.pythonhosted.org/packages/57/b0/0e52c878c53f245edd3a11020f20979b3f490f245af532c7cae3027754b5/idna-3.19-py3-none-any.whl#sha256=815e7be7a7806d54abb586dc943addc79e8b2ee16915059658cbeff4b1b43bf4
proto-plus @ https://files.pythonhosted.org/packages/41/5d/0f04b85dafdc3250ced7f2592efc17dce7f40712e941e9632202481e600d/proto_plus-1.28.4-py3-none-any.whl#sha256=4b01341272f8a348db3f003b6143109f83ab43091019d5181b3fcdf500ab32aa
protobuf @ https://files.pythonhosted.org/packages/0e/4e/12cb93270967a2affff5b3f720694700d4d87712a67afd05c8cb3f6fa52c/protobuf-7.36.0-cp310-abi3-win_amd64.whl#sha256=1781cc1de61249b750848029bca452c0a8b7e990080316b9bbc2518b2117b488
pyasn1 @ https://files.pythonhosted.org/packages/9a/3b/6163796d69c3977d1e4287bea4a6979161cbbdd170ebb430511e8e1999ce/pyasn1-0.6.4-py3-none-any.whl#sha256=deda9277cfd454080ec40b207fb6df82206a3a2688735233cdcd8d3d565f088b
pyasn1-modules @ https://files.pythonhosted.org/packages/47/8d/d529b5d697919ba8c11ad626e835d4039be708a35b0d22de83a269a6682c/pyasn1_modules-0.4.2-py3-none-any.whl#sha256=29253a9207ce32b64c3ac6600edc75368f98473906e8fd1043bd6b5b1de2c14a
pycparser @ https://files.pythonhosted.org/packages/0c/c3/44f3fbbfa403ea2a7c779186dc20772604442dde72947e7d01069cbe98e3/pycparser-3.0-py3-none-any.whl#sha256=b727414169a36b7d524c1c3e31839a521725078d7b2ff038656844266160a992
requests @ https://files.pythonhosted.org/packages/a0/f4/c67b0b3f1b9245e8d266f0f112c500d50e5b4e83cb6f3b71b6528104182a/requests-2.34.2-py3-none-any.whl#sha256=2a0d60c172f83ac6ab31e4554906c0f3b3588d37b5cb939b1c061f4907e278e0
typing-extensions @ https://files.pythonhosted.org/packages/49/d3/b8441a820a491ddfc024b0b0cf0393375b75ea13866d9c66727e54c2fc80/typing_extensions-4.16.0-py3-none-any.whl#sha256=481caa481374e813c1b176ada14e97f1f67a4539ce9cfeb3f350d78d6370c2e8
urllib3 @ https://files.pythonhosted.org/packages/7f/3e/5db95bcf282c52709639744ca2a8b149baccf648e39c8cc87553df9eae0c/urllib3-2.7.0-py3-none-any.whl#sha256=9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897
````

### FILE: `google_merchant_product_sync/sdk-artifact.lock.json`
```yaml
block_id: "PY-GOOGLE-MERCHANT-SYNC:sdk-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "local receipt derived from official GitHub/PyPI metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "b070f37394487a1fd674c9e492543275b232255379e03095ea52008ca536f0ee"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-python-sdk-artifact-lock/v1",
  "provider": "Google",
  "distribution": "google-shopping-merchant-products",
  "version": "1.8.0",
  "requires_python": ">=3.10",
  "release_level": "Production/Stable",
  "wheel": {
    "filename": "google_shopping_merchant_products-1.8.0-py3-none-any.whl",
    "bytes": 244077,
    "sha256": "722ef095eca35126255c129964a61818f0d7bc9eeef7f1da12732d81c1dcd505"
  },
  "sdist": {
    "filename": "google_shopping_merchant_products-1.8.0.tar.gz",
    "bytes": 241375,
    "sha256": "08f0ab5ccbdf4df65cdf8c0d763ba161acc02f61897f6bf023a20223acbbb7f9"
  },
  "source": {
    "repository": "googleapis/google-cloud-python",
    "release": "google-shopping-merchant-products-v1.8.0",
    "commit": "97d7b42cd74b41211f5ec8871cc0dd15debdb1a0",
    "archive_bytes": 206581382,
    "archive_sha256": "cde3981299462b8346f725d6fd4bf7973972b6b4263cfc7399081e0514ad6b1c"
  },
  "verified_platform_lock": "CPython 3.14 / Windows x86-64 / 20 wheels",
  "license_expression": "Apache-2.0",
  "verified_at": "2026-08-26"
}
````

### FILE: `google_merchant_product_sync/provider-profile.template.json`
```yaml
block_id: "PY-GOOGLE-MERCHANT-SYNC:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed Google Merchant account/write profile"
license: "LicenseRef-Workspace-Owner"
sha256: "09184bf9e9e37e1dbb5603232e66f828ed66c65ca33e7b151bd434b8b76d7374"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-google-merchant-product-sync-profile/v1",
  "provider": "Google Merchant API",
  "sdk": "google-shopping-merchant-products==1.8.0",
  "api": "Products v1",
  "decision": "BLOCKED_ACCOUNT_AND_EXTERNAL_WRITE_APPROVAL_REQUIRED",
  "cloud_project_reference": "",
  "merchant_account_id": "",
  "data_source_name": "",
  "data_source_type": "API_PRIMARY_NOT_PROVEN",
  "merchant_api_enabled": false,
  "billing_approved": false,
  "adc_identity_reference": "",
  "account_access_proven": false,
  "test_account_contract_proven": false,
  "external_product_write_approved": false,
  "approved_offer_ids": [],
  "quota_and_cost_approved": false,
  "refresh_cadence_days": 30,
  "reconciliation_owner": "",
  "evidence_retention": "",
  "automatic_local_business_write": false
}
````

### FILE: `google_merchant_product_sync/product.template.json`
```yaml
block_id: "PY-GOOGLE-MERCHANT-SYNC:product-template:v1"
operation: CREATE
provenance: AUTHORED
source: "local bounded mapping to official ProductInput/ProductAttributes fields"
license: "LicenseRef-Workspace-Owner"
sha256: "8a78df49b1cc2bb3e1172dd2f994eb7239d5c1291e76b56ac09d1c5afbf1d737"
variables: []
secrets_allowed: false
```
````json
{
  "offer_id": "SKU12345",
  "content_language": "en",
  "feed_label": "AR",
  "title": "REPLACE_WITH_APPROVED_TITLE",
  "description": "REPLACE_WITH_APPROVED_DESCRIPTION",
  "link": "https://example.invalid/products/SKU12345",
  "image_link": "https://example.invalid/images/SKU12345.jpg",
  "availability": "IN_STOCK",
  "condition": "NEW",
  "price": {
    "amount_micros": 1000000,
    "currency_code": "ARS"
  },
  "brand": "",
  "gtins": []
}
````

### FILE: `google_merchant_product_sync/sync_product.py`
```yaml
block_id: "PY-GOOGLE-MERCHANT-SYNC:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local wrapper invoking official Google Merchant Products v1 clients/messages"
license: "LicenseRef-Workspace-Owner"
sha256: "7c926b92e124d2a38f253009a6c238e436d52ba2bc6aa9675381f8e30ced9d75"
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
import time
import uuid
from typing import Any, Callable, Protocol
from urllib.parse import urlparse

from google.api_core.exceptions import NotFound
from google.protobuf.json_format import MessageToDict
from google.shopping import merchant_products_v1 as merchant
from google.shopping.type import Price


SDK_VERSION = "1.8.0"
PROFILE_SCHEMA = "elite-google-merchant-product-sync-profile/v1"
AVAILABILITY = {
    "IN_STOCK": merchant.Availability.IN_STOCK,
    "OUT_OF_STOCK": merchant.Availability.OUT_OF_STOCK,
    "PREORDER": merchant.Availability.PREORDER,
    "LIMITED_AVAILABILITY": merchant.Availability.LIMITED_AVAILABILITY,
    "BACKORDER": merchant.Availability.BACKORDER,
}
CONDITION = {
    "NEW": merchant.Condition.NEW,
    "USED": merchant.Condition.USED,
    "REFURBISHED": merchant.Condition.REFURBISHED,
}
PRODUCT_KEYS = {
    "offer_id", "content_language", "feed_label", "title", "description",
    "link", "image_link", "availability", "condition", "price", "brand", "gtins",
}


class ProductInputsClient(Protocol):
    def insert_product_input(self, request: Any, *, retry: Any, timeout: float) -> Any: ...


class ProductsClient(Protocol):
    def get_product(self, request: Any, *, retry: Any, timeout: float) -> Any: ...


def sha256(value: bytes) -> str:
    return hashlib.sha256(value).hexdigest()


def json_bytes(value: Any) -> bytes:
    return (json.dumps(value, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")


def message_dict(value: Any) -> dict[str, Any]:
    if isinstance(value, dict):
        return value
    protobuf = getattr(value, "_pb", None)
    if protobuf is not None:
        return MessageToDict(protobuf, preserving_proto_field_name=True)
    to_dict = getattr(value, "to_dict", None)
    if callable(to_dict):
        result = to_dict()
        if isinstance(result, dict):
            return result
    raise TypeError(f"unsupported provider message type: {type(value).__name__}")


def write_json_atomic(path: Path, value: Any) -> bytes:
    data = json_bytes(value)
    temporary = path.with_name(f".{path.name}.{uuid.uuid4().hex}.tmp")
    temporary.write_bytes(data)
    temporary.replace(path)
    return data


def valid_https_url(value: Any) -> bool:
    if not isinstance(value, str) or len(value) > 2048:
        return False
    parsed = urlparse(value)
    return parsed.scheme == "https" and bool(parsed.hostname) and parsed.username is None and parsed.password is None


def read_approved_profile(path: Path) -> dict[str, Any]:
    profile = json.loads(path.read_text(encoding="utf-8"))
    if profile.get("schema") != PROFILE_SCHEMA:
        raise ValueError("unsupported profile schema")
    if profile.get("sdk") != f"google-shopping-merchant-products=={SDK_VERSION}":
        raise ValueError("profile SDK version mismatch")
    if profile.get("decision") not in {"APPROVED_FOR_TEST_ACCOUNT", "APPROVED_FOR_PRODUCTION"}:
        raise PermissionError("Google Merchant external write decision is not approved")
    account = str(profile.get("merchant_account_id", ""))
    if not re.fullmatch(r"[0-9]{3,20}", account):
        raise ValueError("merchant_account_id must contain 3..20 digits")
    expected_source = rf"accounts/{account}/dataSources/[0-9]{{1,20}}"
    if not re.fullmatch(expected_source, str(profile.get("data_source_name", ""))):
        raise ValueError("data_source_name must belong to the approved Merchant account")
    if profile.get("data_source_type") != "API_PRIMARY_PROVEN":
        raise PermissionError("an API primary data source must be proven")
    for key in (
        "merchant_api_enabled", "billing_approved", "account_access_proven",
        "test_account_contract_proven", "external_product_write_approved", "quota_and_cost_approved",
    ):
        if profile.get(key) is not True:
            raise PermissionError(f"{key} must be proven")
    if not str(profile.get("cloud_project_reference", "")).strip():
        raise ValueError("cloud_project_reference is required")
    if not str(profile.get("adc_identity_reference", "")).strip():
        raise ValueError("adc_identity_reference is required")
    if not str(profile.get("reconciliation_owner", "")).strip() or not str(profile.get("evidence_retention", "")).strip():
        raise ValueError("reconciliation_owner and evidence_retention are required")
    if profile.get("automatic_local_business_write") is not False:
        raise ValueError("automatic local business writes are forbidden")
    cadence = profile.get("refresh_cadence_days")
    if not isinstance(cadence, int) or isinstance(cadence, bool) or cadence < 1 or cadence > 30:
        raise ValueError("refresh_cadence_days must be within 1..30")
    approved = profile.get("approved_offer_ids")
    if not isinstance(approved, list) or not approved or any(not isinstance(value, str) for value in approved):
        raise ValueError("approved_offer_ids must be a non-empty string array")
    return profile


def build_insert_request(profile: dict[str, Any], product: dict[str, Any]) -> merchant.InsertProductInputRequest:
    if set(product) != PRODUCT_KEYS:
        raise ValueError("product keys must match the exact supported schema")
    offer_id = product["offer_id"]
    if not isinstance(offer_id, str) or not re.fullmatch(r"[A-Za-z0-9._-]{1,50}", offer_id):
        raise ValueError("offer_id must be a safe 1..50 character identifier")
    if offer_id not in profile["approved_offer_ids"]:
        raise PermissionError("offer_id is not approved by the profile")
    language = product["content_language"]
    label = product["feed_label"]
    if not isinstance(language, str) or not re.fullmatch(r"[a-z]{2}", language):
        raise ValueError("content_language must be a two-letter lowercase code")
    if not isinstance(label, str) or not re.fullmatch(r"[A-Z0-9-]{2,20}", label):
        raise ValueError("feed_label must be 2..20 uppercase alphanumeric/hyphen characters")
    title = product["title"]
    description = product["description"]
    if not isinstance(title, str) or not title.strip() or len(title) > 150:
        raise ValueError("title must contain 1..150 characters")
    if not isinstance(description, str) or not description.strip() or len(description) > 5000:
        raise ValueError("description must contain 1..5000 characters")
    if not valid_https_url(product["link"]) or not valid_https_url(product["image_link"]):
        raise ValueError("link and image_link must be credential-free HTTPS URLs")
    if product["availability"] not in AVAILABILITY or product["condition"] not in CONDITION:
        raise ValueError("availability or condition is unsupported")
    price = product["price"]
    if not isinstance(price, dict) or set(price) != {"amount_micros", "currency_code"}:
        raise ValueError("price must contain amount_micros and currency_code")
    amount = price["amount_micros"]
    currency = price["currency_code"]
    if not isinstance(amount, int) or isinstance(amount, bool) or amount <= 0 or amount > 10**18:
        raise ValueError("amount_micros must be a positive bounded integer")
    if not isinstance(currency, str) or not re.fullmatch(r"[A-Z]{3}", currency):
        raise ValueError("currency_code must contain three uppercase letters")
    brand = product["brand"]
    if not isinstance(brand, str) or len(brand) > 70:
        raise ValueError("brand must contain at most 70 characters")
    gtins = product["gtins"]
    if not isinstance(gtins, list) or len(gtins) > 10 or any(not isinstance(value, str) or not re.fullmatch(r"(?:[0-9]{8}|[0-9]{12,14})", value) for value in gtins):
        raise ValueError("gtins must contain at most 10 numeric GTIN-8/12/13/14 values")

    attributes = merchant.ProductAttributes(
        title=title.strip(),
        description=description.strip(),
        link=product["link"],
        image_link=product["image_link"],
        availability=AVAILABILITY[product["availability"]],
        condition=CONDITION[product["condition"]],
        price=Price(amount_micros=amount, currency_code=currency),
        brand=brand or None,
        gtins=gtins,
    )
    product_input = merchant.ProductInput(
        offer_id=offer_id,
        content_language=language,
        feed_label=label,
        product_attributes=attributes,
    )
    return merchant.InsertProductInputRequest(
        parent=f"accounts/{profile['merchant_account_id']}",
        data_source=profile["data_source_name"],
        product_input=product_input,
    )


def sync_product(
    inputs_client: ProductInputsClient,
    products_client: ProductsClient,
    request: merchant.InsertProductInputRequest,
    output_directory: Path,
    *,
    timeout_seconds: float = 30.0,
    status_attempts: int = 3,
    status_interval_seconds: float = 1.0,
    sleep_fn: Callable[[float], None] = time.sleep,
) -> dict[str, Any]:
    if importlib.metadata.version("google-shopping-merchant-products") != SDK_VERSION:
        raise RuntimeError(f"google-shopping-merchant-products {SDK_VERSION} is required")
    if timeout_seconds <= 0 or timeout_seconds > 120:
        raise ValueError("timeout must be within 1..120 seconds")
    if status_attempts < 1 or status_attempts > 12 or status_interval_seconds < 0 or status_interval_seconds > 60:
        raise ValueError("status polling bounds are invalid")
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".google-merchant-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    request_payload = message_dict(request)
    request_data = json_bytes(request_payload)
    (stage / "insert-request.json").write_bytes(request_data)
    try:
        inserted = inputs_client.insert_product_input(request, retry=None, timeout=timeout_seconds)
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise
    inserted_payload = message_dict(inserted)
    product_name = getattr(inserted, "product", None)
    if not isinstance(product_name, str) or not re.fullmatch(r"accounts/[0-9]{3,20}/products/.+", product_name):
        shutil.rmtree(stage, ignore_errors=True)
        raise RuntimeError("provider insert response omitted the canonical product resource name")
    insert_data = json_bytes(inserted_payload)
    (stage / "insert-response.json").write_bytes(insert_data)
    offer_id = request.product_input.offer_id
    receipt: dict[str, Any] = {
        "schema": "elite-google-merchant-product-sync-receipt/v1",
        "created_at": datetime.now(timezone.utc).isoformat(),
        "provider": "Google Merchant API",
        "sdk_version": SDK_VERSION,
        "api_version": "products/v1",
        "operation": "productInputs.insert + products.get",
        "processing_state": "INSERT_ACCEPTED_STATUS_UNCHECKED",
        "offer_id_sha256": sha256(offer_id.encode("utf-8")),
        "product_name_sha256": sha256(product_name.encode("utf-8")),
        "request_sha256": sha256(request_data),
        "insert_response_sha256": sha256(insert_data),
        "external_product_write": True,
        "automatic_local_business_write": False,
        "approval_inferred": False,
    }
    (stage / "SYNC_RECEIPT.json").write_bytes(json_bytes(receipt))
    stage.replace(output)

    request_get = merchant.GetProductRequest(name=product_name)
    for attempt in range(1, status_attempts + 1):
        try:
            product = products_client.get_product(request_get, retry=None, timeout=timeout_seconds)
            product_payload = message_dict(product)
            product_data = write_json_atomic(output / "product-response.json", product_payload)
            receipt["processing_state"] = "PROCESSED_PRODUCT_STATUS_OBSERVED"
            receipt["status_attempts"] = attempt
            receipt["product_response_sha256"] = sha256(product_data)
            write_json_atomic(output / "SYNC_RECEIPT.json", receipt)
            return receipt
        except NotFound:
            if attempt < status_attempts:
                sleep_fn(status_interval_seconds)
                continue
            receipt["processing_state"] = "INSERT_ACCEPTED_PRODUCT_NOT_YET_VISIBLE"
            receipt["status_attempts"] = attempt
            write_json_atomic(output / "SYNC_RECEIPT.json", receipt)
            return receipt
        except BaseException as error:
            receipt["processing_state"] = "INSERT_ACCEPTED_STATUS_CHECK_FAILED"
            receipt["status_attempts"] = attempt
            receipt["status_error_type"] = type(error).__name__
            write_json_atomic(output / "SYNC_RECEIPT.json", receipt)
            raise
    raise AssertionError("unreachable status loop")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--product", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--timeout-seconds", type=float, default=30.0)
    args = parser.parse_args()
    profile = read_approved_profile(args.profile)
    product = json.loads(args.product.read_text(encoding="utf-8"))
    request = build_insert_request(profile, product)
    inputs_client = merchant.ProductInputsServiceClient()
    products_client = merchant.ProductsServiceClient()
    receipt = sync_product(inputs_client, products_client, request, args.output, timeout_seconds=args.timeout_seconds)
    print(json.dumps({"status": "GOOGLE_MERCHANT_PRODUCT_SYNC_PASS", "receipt": receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `google_merchant_product_sync/test_sync_product.py`
```yaml
block_id: "PY-GOOGLE-MERCHANT-SYNC:test:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract/negative/durability tests plus official SDK signature probe"
license: "LicenseRef-Workspace-Owner"
sha256: "8e4a97a4d45f6a4e39509894784e5c36f141d6d49b6e75d1c49b0e68468d4075"
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

from google.api_core.exceptions import NotFound
from google.shopping import merchant_products_v1 as merchant

from sync_product import build_insert_request, read_approved_profile, sync_product


def approved_profile():
    return {
        "schema": "elite-google-merchant-product-sync-profile/v1",
        "sdk": "google-shopping-merchant-products==1.8.0",
        "decision": "APPROVED_FOR_TEST_ACCOUNT",
        "cloud_project_reference": "projects/test",
        "merchant_account_id": "12345",
        "data_source_name": "accounts/12345/dataSources/67890",
        "data_source_type": "API_PRIMARY_PROVEN",
        "merchant_api_enabled": True,
        "billing_approved": True,
        "adc_identity_reference": "workload:test",
        "account_access_proven": True,
        "test_account_contract_proven": True,
        "external_product_write_approved": True,
        "approved_offer_ids": ["SKU12345"],
        "quota_and_cost_approved": True,
        "refresh_cadence_days": 30,
        "reconciliation_owner": "commerce",
        "evidence_retention": "30d",
        "automatic_local_business_write": False,
    }


def product_spec():
    return {
        "offer_id": "SKU12345",
        "content_language": "en",
        "feed_label": "AR",
        "title": "Electric bicycle",
        "description": "Approved product description",
        "link": "https://shop.example.com/products/SKU12345",
        "image_link": "https://shop.example.com/images/SKU12345.jpg",
        "availability": "IN_STOCK",
        "condition": "NEW",
        "price": {"amount_micros": 1500000000, "currency_code": "ARS"},
        "brand": "Example",
        "gtins": ["12345678"],
    }


class FakeMessage:
    def __init__(self, payload, product=None): self.payload = payload; self.product = product
    def to_dict(self): return self.payload


class FakeInputs:
    def __init__(self, response=None, error=None): self.response=response; self.error=error; self.calls=[]
    def insert_product_input(self, request, *, retry, timeout):
        self.calls.append((request, retry, timeout))
        if self.error: raise self.error
        return self.response


class FakeProducts:
    def __init__(self, outcomes): self.outcomes=list(outcomes); self.calls=[]
    def get_product(self, request, *, retry, timeout):
        self.calls.append((request, retry, timeout))
        outcome=self.outcomes.pop(0)
        if isinstance(outcome, BaseException): raise outcome
        return outcome


class MerchantSyncTests(unittest.TestCase):
    def test_builds_official_v1_request(self):
        request = build_insert_request(approved_profile(), product_spec())
        self.assertIsInstance(request, merchant.InsertProductInputRequest)
        self.assertEqual(request.parent, "accounts/12345")
        self.assertEqual(request.data_source, "accounts/12345/dataSources/67890")
        self.assertEqual(request.product_input.offer_id, "SKU12345")
        self.assertEqual(request.product_input.product_attributes.price.currency_code, "ARS")
        self.assertEqual(request.product_input.product_attributes.availability, merchant.Availability.IN_STOCK)

    def test_insert_and_processed_status_are_preserved(self):
        request = build_insert_request(approved_profile(), product_spec())
        name = "accounts/12345/products/en~AR~SKU12345"
        inputs = FakeInputs(FakeMessage({"product": name, "offer_id": "SKU12345"}, name))
        products = FakeProducts([FakeMessage({"name": name, "product_status": {"destination_statuses": []}})])
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = sync_product(inputs, products, request, output, status_interval_seconds=0)
            self.assertEqual(receipt["processing_state"], "PROCESSED_PRODUCT_STATUS_OBSERVED")
            self.assertTrue((output / "insert-request.json").is_file())
            self.assertTrue((output / "insert-response.json").is_file())
            self.assertTrue((output / "product-response.json").is_file())
            self.assertNotIn("SKU12345", json.dumps(receipt))
            self.assertFalse(receipt["approval_inferred"])
            self.assertIsNone(inputs.calls[0][1])

    def test_not_found_is_retained_as_pending_not_approval(self):
        request = build_insert_request(approved_profile(), product_spec())
        name = "accounts/12345/products/en~AR~SKU12345"
        inputs = FakeInputs(FakeMessage({"product": name}, name))
        products = FakeProducts([NotFound("pending"), NotFound("pending")])
        with tempfile.TemporaryDirectory() as temp:
            receipt = sync_product(inputs, products, request, Path(temp) / "out", status_attempts=2, status_interval_seconds=0)
            self.assertEqual(receipt["processing_state"], "INSERT_ACCEPTED_PRODUCT_NOT_YET_VISIBLE")
            self.assertFalse(receipt["approval_inferred"])

    def test_insert_failure_is_atomic(self):
        request = build_insert_request(approved_profile(), product_spec())
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "out"
            with self.assertRaisesRegex(RuntimeError, "insert failed"):
                sync_product(FakeInputs(error=RuntimeError("insert failed")), FakeProducts([]), request, output)
            self.assertFalse(output.exists())
            self.assertEqual(list(Path(temp).iterdir()), [])

    def test_status_failure_preserves_post_write_evidence(self):
        request = build_insert_request(approved_profile(), product_spec())
        name = "accounts/12345/products/en~AR~SKU12345"
        inputs = FakeInputs(FakeMessage({"product": name}, name))
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "out"
            with self.assertRaisesRegex(RuntimeError, "status failed"):
                sync_product(inputs, FakeProducts([RuntimeError("status failed")]), request, output)
            receipt = json.loads((output / "SYNC_RECEIPT.json").read_text(encoding="utf-8"))
            self.assertEqual(receipt["processing_state"], "INSERT_ACCEPTED_STATUS_CHECK_FAILED")
            self.assertEqual(receipt["status_error_type"], "RuntimeError")

    def test_profile_and_product_fail_closed(self):
        blocked = approved_profile(); blocked["decision"] = "BLOCKED_ACCOUNT_AND_EXTERNAL_WRITE_APPROVAL_REQUIRED"
        with tempfile.TemporaryDirectory() as temp:
            path = Path(temp) / "profile.json"
            path.write_text(json.dumps(blocked), encoding="utf-8")
            with self.assertRaises(PermissionError): read_approved_profile(path)
        invalid = product_spec(); invalid["offer_id"] = "NOT-APPROVED"
        with self.assertRaises(PermissionError): build_insert_request(approved_profile(), invalid)
        invalid = product_spec(); invalid["link"] = "http://insecure.example.com"
        with self.assertRaises(ValueError): build_insert_request(approved_profile(), invalid)

    def test_official_sdk_contract(self):
        insert = inspect.signature(merchant.ProductInputsServiceClient.insert_product_input)
        get = inspect.signature(merchant.ProductsServiceClient.get_product)
        self.assertIn("request", insert.parameters)
        self.assertIn("retry", insert.parameters)
        self.assertIn("request", get.parameters)
        self.assertIn("product_status", merchant.Product.meta.fields)


if __name__ == "__main__":
    unittest.main()
````

### FILE: `google_merchant_product_sync/README.md`
```yaml
block_id: "PY-GOOGLE-MERCHANT-SYNC:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local bounded runbook with official authority boundaries"
license: "LicenseRef-Workspace-Owner"
sha256: "8d51fbed45a12a6f094298b5e64ecb6b9e18c4a9cab7537328fb5ecc295a2998"
variables: []
secrets_allowed: false
```
````markdown
# Google Merchant Product Sync evidence adapter

This module invokes Google's official GA `google-shopping-merchant-products` Python client `1.8.0`, `ProductInputsServiceClient.insert_product_input`, and `ProductsServiceClient.get_product`. The SDK and generated v1 models are Google code under Apache-2.0. `sync_product.py`, profiles, validation, and receipts are local integration glue and are not represented as Google-authored code.

The provided profile is blocked. Before use, prove the Google Cloud project and billing decision, Merchant API enablement, ADC identity, Merchant account access, an API primary data source, a test-account contract, allowed offer IDs, quota/cost, retention, reconciliation ownership, and explicit authorization for the external product write. Credentials remain in Application Default Credentials or an authorized runtime secret store; never place tokens or keys in Markdown, JSON, CLI arguments, logs, or receipts.

For the verified CPython 3.14 / Windows x86-64 lane:

```text
python -m venv .venv
.venv\Scripts\python -m pip install --require-hashes -r requirements-windows-py314.lock
.venv\Scripts\python -m pip check
.venv\Scripts\python -m unittest -v test_sync_product.py
.venv\Scripts\python sync_product.py --profile C:\secure\merchant-profile.json --product C:\approved\product.json --output C:\evidence\merchant-SKU12345
```

For another OS/Python ABI, treat `requirements-direct.in` as the exact top-level input and produce a new fully hash-pinned candidate lock under the dependency update contract before execution. Do not reuse the Windows lock on another target.

The adapter stores the exact insert request/response and then reads the processed Product status. A successful insert is never labelled approved: Google states that insertion starts processing and approval must be read from `product_status`. If the product is not yet visible, the receipt remains pending; if the post-write status check fails, durable insert evidence is retained before the error is raised. No local business record is written automatically.
````

## 6. Configuration surface

| Key | Required | Secret | Gate |
|---|---:|---:|---|
| Cloud project, API, billing | yes | no | habilitación y costo aprobados |
| `merchant_account_id`, `data_source_name` | yes | no | misma cuenta; `API_PRIMARY_PROVEN` |
| ADC identity | yes | reference only | identidad/scopes de menor privilegio probados |
| test-account contract | yes | no | `true` incluso antes de producción |
| external product write | yes | no | aprobación explícita; offer allowlist |
| refresh cadence | yes | no | 1..30 días, conforme expiración indicada por Google |
| quota/cost, retention, reconciliation owner | yes | no | cerrados antes del write |
| product JSON | yes | business data | shape exacto; IDs/URLs/enums/precio/GTIN validados |

## 7. Dependency bill

| Dependency | Exact version/artifact | License | Purpose |
|---|---|---|---|
| Google Merchant Products | 1.8.0 wheel SHA `722ef095…d505` | Apache-2.0 | clientes/mensajes oficiales v1 |
| Google Shopping Type | 1.5.0 SHA `afc1180f…5072` | Apache-2.0 | `Price` oficial |
| google-api-core/auth/common-protos | 2.34.0 / 2.57.0 / 1.75.2 | Apache-2.0 | transporte/auth/protos |
| grpcio/grpcio-status/protobuf/proto-plus | 1.83.0 / 1.83.0 / 7.36.0 / 1.28.4 | Apache/BSD | runtime generado |
| requests/urllib3/certifi/charset/idna | 2.34.2 / 2.7.0 / 2026.7.22 / 3.5.1 / 3.19 | Apache/MIT/MPL/BSD | HTTP y trust store |
| cryptography/cffi/pycparser/pyasn1 modules | 50.0.1 / 2.1.1 / 3.0 / 0.6.4+0.4.2 | Apache/BSD/MIT | autenticación |
| typing-extensions | 4.16.0 | PSF-2.0 | typing runtime |

Los veinte URLs y hashes completos están en `requirements-windows-py314.lock`; la tabla no sustituye ese lock.

## 8. Apply order

1. Materializar `OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x` y este pack.
2. En Windows/CPython 3.14 instalar el lock exacto; en otro target generar y aprobar un lock candidato desde `requirements-direct.in`.
3. Ejecutar `pip check`, siete tests y probe de firmas/modelos oficiales.
4. Completar profile y product fuera del source tree; configurar ADC por mecanismo runtime autorizado.
5. Ejecutar test account, conservar request/insert/status/receipt y reconciliar con Merchant Center.
6. Promover sólo con product status/issue handling, refresh, quota, observabilidad, security y rollback probados.

Rollback: detener sync/refresh, preservar receipt y, si corresponde, ejecutar un pack de delete compensatorio explícitamente aprobado; nunca borrar localmente ni ocultar que el write externo ya ocurrió.

## 9. Verification

```text
python -m venv .venv
.venv/Scripts/python -m pip install --require-hashes -r google_merchant_product_sync/requirements-windows-py314.lock
.venv/Scripts/python -m pip check
cd google_merchant_product_sync
../.venv/Scripts/python -m unittest -v test_sync_product.py
```

Esperado verificado: veinte wheels exactos instalados; `pip check` PASS; SDK 1.8.0; siete tests PASS; firmas `insert_product_input`/`get_product` y campo `Product.product_status` presentes. No hubo llamada Merchant ni cuenta real.

Gates productivos: cuenta/data source/ADC/test account, write/offer, product spec/GTIN ownership, status/issues, refresh ≤30 días, quota/cost, concurrencia/idempotencia, observabilidad sin datos sensibles, seguridad, carga, compensación y reconciliación deben estar `PROVEN`. Hasta entonces `admission: CONDITIONED`.

## 10. Reconstruction evidence

`reconstruction_evidence/GOOGLE_MERCHANT_PRODUCT_SYNC_ADAPTER_2026-08-26_V1.md` registra release/commit/archive/wheel/sdist, grafo de veinte wheels, instalación limpia, siete tests, probe de firmas, fallo de expansión retenido y gates globales. Evidencia local no sustituye cuenta Google ni producción.
