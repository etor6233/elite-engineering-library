# Python Official Meta Page Write Adapter

## 1. Metadata

```yaml
pack_id: "PYTHON-OFFICIAL-META-PAGE-WRITE-ADAPTER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Official SDK leaf for one approved Facebook Page text publication, exact GET observation/reconciliation and separately approved revocation; no automatic write retry."
stacks: ["CPython 3.14 Windows x86-64", "facebook-business 26.0.1", "Graph API v26.0"]
compatible_with: ["exact Meta 18-wheel lock", "caller-owned authorization and durable outbound fence"]
incompatible_with: ["claim of complete social scheduler", "automatic retry of uncertain POST/DELETE", "unapproved content or changed scope", "Instagram/TikTok/LinkedIn/media/ad mutations"]
license_expression: "LicenseRef-Workspace-Owner AND LicenseRef-Meta-Platform AND third-party-wheel-licenses"
upstream_sources: ["https://github.com/facebook/facebook-python-business-sdk/tree/788f363d15b1269ab5efb7cd00fb5e3b133cd99b"]
verified_at: "2026-09-11"
```

Unpublished candidate. Wrapper/tests/contracts are AUTHORED glue invoking the
DEPENDENCY_PIN official SDK; no local domain code is attributed to Meta. License
bytes are VERBATIM and only authorize the uses permitted by Meta Platform terms.
The adapter does not upgrade Ads reporting's read-only admission to posting.
Its independent narrow source/contract evidence is in seven-core-evidence.md.

## 2. Applicability and conditions

Selected Facebook Page, approved exact text and caller-held durable outbound fence.
The leaf checks binding; it does not perform the approval or claim a database
lease. Whole-core scheduler, quota/approval persistence, worker bridge, UI and
operator recovery remain LOCAL_IMPLEMENTATION_GAP, not missing credentials alone.
Current live permission/app/Page ownership contract must be verified before use.

## 3. Architecture contract

Scope plus approved-intent fields -> official SDK POST -> bound SDK GET -> receipt.
Any unconfirmed result is unknown and never automatically retried. Known-reference
reconciliation is GET only. Revoke requires a distinct intent and bound prior GET,
then raw official SDK DELETE success=True. The SDK object parser removes success,
so its returned object is deliberately not used to infer deletion. No scheduling,
pricing, moderation or customer-eligibility algorithm is invented.

## 4. Exact file manifest

```text
CREATE meta_page_write/adapter.py
CREATE meta_page_write/tests/test_adapter.py
CREATE meta_page_write/requirements-windows-py314.lock
CREATE meta_page_write/sdk-artifact.lock.json
CREATE meta_page_write/source-map.json
CREATE meta_page_write/LICENSE.Meta.txt
CREATE meta_page_write/README.md
```

## 5. Materialization blocks

### FILE: `meta_page_write/adapter.py`
```yaml
block_id: "PY-OFFICIAL-META-PAGE-WRITE:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "Local integration/verification/dependency declaration; official SDK code is dependency-pinned, not locally authored"
license: "LicenseRef-Workspace-Owner"
sha256: "a397d84a5f92fc8d9d3a44850ad568875da837e59800aeab324fcb5adc6a9e1e"
variables: []
secrets_allowed: false
```
````python
"""AUTHORED binding glue over the locked official Meta Business SDK.

This leaf does not authorize content, claim a durable delivery lease, schedule
posts or retry writes. Its caller must own those decisions before invoking it.
"""
from dataclasses import dataclass
import hashlib
import importlib.metadata
import json
import re

from facebook_business.adobjects.page import Page
from facebook_business.adobjects.pagepost import PagePost
from facebook_business.api import FacebookAdsApi
from facebook_business.session import FacebookSession


class InvalidBinding(ValueError):
    pass


class UnknownDelivery(RuntimeError):
    """No automatic retry is safe; retain the optional reference for GET."""

    def __init__(self, code, provider_reference=None):
        super().__init__(code)
        self.code = code
        self.provider_reference = provider_reference


def _hash(value):
    return isinstance(value, str) and re.fullmatch(r"[a-f0-9]{64}", value) is not None


def _identifier(value):
    return isinstance(value, str) and re.fullmatch(r"[A-Za-z0-9_-]{1,128}", value) is not None


def _page(value):
    return isinstance(value, str) and re.fullmatch(r"[1-9][0-9]{0,31}", value) is not None


@dataclass(frozen=True)
class Scope:
    tenant_id: str
    page_id: str
    profile_sha256: str

    def validate(self):
        if not _identifier(self.tenant_id) or not _page(self.page_id) or not _hash(self.profile_sha256):
            raise InvalidBinding("invalid configured page scope")


@dataclass(frozen=True)
class ApprovedIntent:
    tenant_id: str
    page_id: str
    profile_sha256: str
    approval_id: str
    delivery_key: str
    operation: str
    message: str
    content_sha256: str


@dataclass(frozen=True)
class Receipt:
    tenant_id: str
    page_id: str
    profile_sha256: str
    approval_id: str
    delivery_key: str
    provider_reference: str
    state: str
    content_sha256: str
    evidence_sha256: str


def create_api(app_id, app_secret, page_access_token):
    """Credentials stay in memory; no global default API or debug logger."""
    if importlib.metadata.version("facebook-business") != "26.0.1":
        raise InvalidBinding("SDK version differs from admitted lock")
    if not _page(app_id) or not isinstance(app_secret, str) or not app_secret or not isinstance(page_access_token, str) or not page_access_token:
        raise InvalidBinding("credentials are unavailable")
    session = FacebookSession(app_id, app_secret, page_access_token, timeout=10)
    session.requests.trust_env = False
    return FacebookAdsApi(session, api_version="v26.0", enable_debug_logger=False)


class FacebookPageWriteAdapter:
    def __init__(self, scope, api):
        if not isinstance(scope, Scope) or not isinstance(api, FacebookAdsApi):
            raise InvalidBinding("scope and official SDK API required")
        scope.validate()
        # This attribute is pinned by the exact official api.py source lock.
        if api._api_version != "v26.0" or api._enable_debug_logger:
            raise InvalidBinding("Graph API version differs from admitted lock")
        self.scope, self.api = scope, api

    def _validate(self, intent, operation):
        if not isinstance(intent, ApprovedIntent):
            raise InvalidBinding("bound intent required")
        if (intent.tenant_id, intent.page_id, intent.profile_sha256) != (self.scope.tenant_id, self.scope.page_id, self.scope.profile_sha256):
            raise InvalidBinding("intent does not match configured scope")
        if intent.operation != operation or not _identifier(intent.approval_id) or not _identifier(intent.delivery_key) or not 16 <= len(intent.delivery_key) <= 128:
            raise InvalidBinding("intent approval, action or delivery binding invalid")
        try:
            body = intent.message.encode("utf-8")
        except (AttributeError, UnicodeError):
            raise InvalidBinding("invalid UTF-8 content") from None
        # Local bounded-transport limit; not a claim about Meta's character limit.
        if not body.strip() or len(body) > 16384 or not _hash(intent.content_sha256) or hashlib.sha256(body).hexdigest() != intent.content_sha256:
            raise InvalidBinding("content differs from the approved bytes")

    def _reference(self, provider_reference):
        if not isinstance(provider_reference, str) or re.fullmatch(re.escape(self.scope.page_id) + r"_[1-9][0-9]{0,63}", provider_reference) is None:
            raise InvalidBinding("post reference does not belong to configured page")

    def _observe(self, intent, provider_reference):
        self._reference(provider_reference)
        try:
            post = PagePost(provider_reference, api=self.api).api_get(fields=["id", "from", "message", "is_published"])
            value = post.export_all_data()
            if value.get("id") != provider_reference or value.get("from", {}).get("id") != self.scope.page_id or value.get("message") != intent.message or value.get("is_published") is not True:
                raise ValueError("post binding differs")
        except Exception:
            # SDK exceptions can include request details; keep only a safe code.
            raise UnknownDelivery("POST_OBSERVATION_UNCONFIRMED", provider_reference) from None
        return {"id": provider_reference, "from": {"id": self.scope.page_id}, "message": intent.message, "is_published": True}

    def _receipt(self, intent, reference, state, observed):
        evidence = hashlib.sha256(json.dumps(observed, sort_keys=True, separators=(",", ":"), ensure_ascii=False).encode("utf-8")).hexdigest()
        return Receipt(intent.tenant_id, intent.page_id, intent.profile_sha256, intent.approval_id, intent.delivery_key, reference, state, intent.content_sha256, evidence)

    def publish_once(self, intent):
        """Caller must hold its durable outbound fence. One POST, then one GET."""
        self._validate(intent, "publish")
        reference = None
        try:
            result = Page(self.scope.page_id, api=self.api).create_feed(params={"message": intent.message, "published": True})
            reference = result.get("id")
            self._reference(reference)
        except Exception:
            raise UnknownDelivery("PUBLISH_OUTCOME_UNKNOWN") from None
        observed = self._observe(intent, reference)
        return self._receipt(intent, reference, "published", observed)

    def reconcile(self, intent, provider_reference):
        """GET only; requires a known post ID. Never republishes an unknown POST."""
        self._validate(intent, "publish")
        observed = self._observe(intent, provider_reference)
        return self._receipt(intent, provider_reference, "published", observed)

    def revoke_once(self, intent, provider_reference):
        """Separate caller-approved revoke; verifies page and original bytes first."""
        self._validate(intent, "revoke")
        observed = self._observe(intent, provider_reference)
        try:
            # PagePost.api_delete's ObjectParser discards the success field.
            # Keep the official SDK's raw response for this same node DELETE.
            result = self.api.call("DELETE", (provider_reference,)).json()
            if not isinstance(result, dict) or result.get("success") is not True:
                raise ValueError("delete not confirmed")
        except Exception:
            raise UnknownDelivery("REVOKE_OUTCOME_UNKNOWN", provider_reference) from None
        return self._receipt(intent, provider_reference, "revoked", {"observed_before_delete": observed, "delete_success": True})
````

### FILE: `meta_page_write/tests/test_adapter.py`
```yaml
block_id: "PY-OFFICIAL-META-PAGE-WRITE:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "Local integration/verification/dependency declaration; official SDK code is dependency-pinned, not locally authored"
license: "LicenseRef-Workspace-Owner"
sha256: "36f452d3842e71252b130e352b85619ff342138b43eb58cb34d8651d3930c0b8"
variables: []
secrets_allowed: false
```
````python
"""Real locked SDK + requests transport fixtures; every HTTP send is intercepted."""
from dataclasses import replace, asdict
import hashlib
import json
import unittest
from urllib.parse import urlsplit, parse_qs

import requests
from requests.adapters import BaseAdapter

from meta_page_write.adapter import ApprovedIntent, Scope, FacebookPageWriteAdapter, InvalidBinding, UnknownDelivery, create_api


class FixtureTransport(BaseAdapter):
    def __init__(self, steps):
        self.steps = list(steps)
        self.calls = []

    def send(self, request, **kwargs):
        url = urlsplit(request.url)
        method, path, result = self.steps.pop(0)
        # No networking exists in this transport; unexpected calls fail the test.
        assert url.scheme == "https" and url.hostname == "graph.facebook.com"
        assert (request.method, url.path.rstrip("/")) == (method, path)
        assert kwargs["timeout"] == 10
        body = request.body.decode() if isinstance(request.body, bytes) else request.body or ""
        self.calls.append({"method": request.method, "path": url.path.rstrip("/"), "body": parse_qs(body), "query": parse_qs(url.query)})
        if isinstance(result, Exception):
            raise result
        response = requests.Response()
        response.status_code = 200
        response._content = json.dumps(result).encode()
        response.headers["Content-Type"] = "application/json"
        response.encoding = "utf-8"
        response.request = request
        return response

    def close(self):
        pass


class PageWriteTests(unittest.TestCase):
    def setUp(self):
        self.scope = Scope("tenant-one", "123", "a" * 64)
        message = "Texto aprobado\nSin cambios: ñ"
        self.intent = ApprovedIntent("tenant-one", "123", "a" * 64, "approval-one", "delivery-key-0001", "publish", message, hashlib.sha256(message.encode()).hexdigest())
        self.observed = {"id": "123_456", "from": {"id": "123"}, "message": message, "is_published": True}

    def adapter(self, steps):
        api = create_api("789", "fixture-app-secret", "fixture-page-token")
        transport = FixtureTransport(steps)
        api._session.requests.mount("https://", transport)
        api._session.requests.mount("http://", transport)
        self.addCleanup(api._session.requests.close)
        return FacebookPageWriteAdapter(self.scope, api), transport

    def test_publish_official_sdk_then_bound_get(self):
        adapter, transport = self.adapter([("POST", "/v26.0/123/feed", {"id": "123_456"}), ("GET", "/v26.0/123_456", self.observed)])
        receipt = adapter.publish_once(self.intent)
        self.assertEqual(receipt.state, "published")
        self.assertEqual(receipt.provider_reference, "123_456")
        self.assertEqual(receipt.profile_sha256, self.scope.profile_sha256)
        self.assertEqual(receipt.approval_id, "approval-one")
        self.assertEqual(receipt.delivery_key, "delivery-key-0001")
        self.assertEqual(receipt.content_sha256, self.intent.content_sha256)
        self.assertEqual(transport.calls[0]["body"], {"message": [self.intent.message], "published": ["true"]})
        self.assertEqual(transport.calls[1]["query"]["fields"], ["id,from,message,is_published"])
        self.assertEqual(transport.calls[0]["query"]["access_token"], ["fixture-page-token"])
        self.assertEqual(len(transport.calls[0]["query"]["appsecret_proof"][0]), 64)
        self.assertNotIn("fixture-page-token", json.dumps(asdict(receipt)))
        self.assertEqual(transport.steps, [])

    def test_timeout_does_not_retry_or_claim_publication(self):
        adapter, transport = self.adapter([("POST", "/v26.0/123/feed", requests.Timeout("fixture-page-token"))])
        with self.assertRaises(UnknownDelivery) as caught:
            adapter.publish_once(self.intent)
        self.assertEqual(caught.exception.code, "PUBLISH_OUTCOME_UNKNOWN")
        self.assertIsNone(caught.exception.provider_reference)
        self.assertNotIn("fixture-page-token", str(caught.exception))
        self.assertEqual(len(transport.calls), 1)

    def test_mismatched_observation_retains_reference_for_get_only_reconciliation(self):
        for change in [{"from": {"id": "999"}}, {"message": "changed"}, {"is_published": False}, {"id": "123_999"}]:
            with self.subTest(change=change):
                wrong = {**self.observed, **change}
                adapter, transport = self.adapter([("POST", "/v26.0/123/feed", {"id": "123_456"}), ("GET", "/v26.0/123_456", wrong), ("GET", "/v26.0/123_456", self.observed)])
                with self.assertRaises(UnknownDelivery) as caught:
                    adapter.publish_once(self.intent)
                self.assertEqual(caught.exception.provider_reference, "123_456")
                self.assertEqual(adapter.reconcile(self.intent, "123_456").state, "published")
                self.assertEqual([x["method"] for x in transport.calls], ["POST", "GET", "GET"])

    def test_wrong_page_response_cannot_be_followed_or_reconciled(self):
        adapter, transport = self.adapter([("POST", "/v26.0/123/feed", {"id": "999_456"})])
        with self.assertRaises(UnknownDelivery):
            adapter.publish_once(self.intent)
        with self.assertRaises(InvalidBinding):
            adapter.reconcile(self.intent, "999_456")
        self.assertEqual(len(transport.calls), 1)

    def test_separate_revoke_approval_checks_original_content_before_delete(self):
        adapter, transport = self.adapter([("GET", "/v26.0/123_456", self.observed), ("DELETE", "/v26.0/123_456", {"success": True})])
        with self.assertRaises(InvalidBinding):
            adapter.revoke_once(self.intent, "123_456")
        intent = replace(self.intent, operation="revoke", approval_id="revoke-approval", delivery_key="revoke-delivery-0001")
        receipt = adapter.revoke_once(intent, "123_456")
        self.assertEqual(receipt.state, "revoked")
        self.assertEqual(receipt.approval_id, "revoke-approval")
        self.assertEqual([x["method"] for x in transport.calls], ["GET", "DELETE"])

    def test_invalid_binding_never_sends(self):
        adapter, transport = self.adapter([])
        changes = [dict(tenant_id="other"), dict(page_id="999"), dict(profile_sha256="b"*64), dict(content_sha256="b"*64), dict(approval_id=""), dict(delivery_key="short"), dict(operation="revoke"), dict(message="\ud800"), dict(message="x"*16385)]
        for change in changes:
            with self.subTest(change=list(change)):
                with self.assertRaises(InvalidBinding):
                    adapter.publish_once(replace(self.intent, **change))
        self.assertEqual(transport.calls, [])
        adapter.api._enable_debug_logger = True
        with self.assertRaises(InvalidBinding):
            FacebookPageWriteAdapter(self.scope, adapter.api)

    def test_revoke_requires_explicit_success_and_never_retries(self):
        intent = replace(self.intent, operation="revoke", approval_id="revoke-approval", delivery_key="revoke-delivery-0001")
        for result in [{"success": False}, {}, requests.Timeout("fixture-page-token")]:
            with self.subTest(result=type(result).__name__):
                adapter, transport = self.adapter([("GET", "/v26.0/123_456", self.observed), ("DELETE", "/v26.0/123_456", result)])
                with self.assertRaises(UnknownDelivery) as caught:
                    adapter.revoke_once(intent, "123_456")
                self.assertEqual(caught.exception.provider_reference, "123_456")
                self.assertEqual(caught.exception.code, "REVOKE_OUTCOME_UNKNOWN")
                self.assertNotIn("fixture-page-token", str(caught.exception))
                self.assertEqual([x["method"] for x in transport.calls], ["GET", "DELETE"])


if __name__ == "__main__":
    unittest.main()
````

### FILE: `meta_page_write/requirements-windows-py314.lock`
```yaml
block_id: "PY-OFFICIAL-META-PAGE-WRITE:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "Local integration/verification/dependency declaration; official SDK code is dependency-pinned, not locally authored"
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

### FILE: `meta_page_write/sdk-artifact.lock.json`
```yaml
block_id: "PY-OFFICIAL-META-PAGE-WRITE:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "Local integration/verification/dependency declaration; official SDK code is dependency-pinned, not locally authored"
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

### FILE: `meta_page_write/source-map.json`
```yaml
block_id: "PY-OFFICIAL-META-PAGE-WRITE:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "Local integration/verification/dependency declaration; official SDK code is dependency-pinned, not locally authored"
license: "LicenseRef-Workspace-Owner"
sha256: "ee4211dd1e4a8e6c22caeca743eddea24de7d119b42157c753a127cbebdd3141"
variables: []
secrets_allowed: false
```
````json
{
  "provenance": "AUTHORED glue / DEPENDENCY_PIN official SDK; no adapted domain algorithm",
  "sdk_commit": "788f363d15b1269ab5efb7cd00fb5e3b133cd99b",
  "sources": [
    {
      "path": "LICENSE",
      "url": "https://raw.githubusercontent.com/facebook/facebook-python-business-sdk/788f363d15b1269ab5efb7cd00fb5e3b133cd99b/LICENSE",
      "bytes": 1035,
      "sha256": "48d97b3c936203a750a3288c2b924327769d62327383a1060244b3d3c05ad01f"
    },
    {
      "path": "facebook_business/adobjects/page.py",
      "url": "https://raw.githubusercontent.com/facebook/facebook-python-business-sdk/788f363d15b1269ab5efb7cd00fb5e3b133cd99b/facebook_business/adobjects/page.py",
      "bytes": 222734,
      "sha256": "0d08e317fe078de5ad2805ab5501423fe51fea39ef90ae4a0160122a5914ee51"
    },
    {
      "path": "facebook_business/adobjects/pagepost.py",
      "url": "https://raw.githubusercontent.com/facebook/facebook-python-business-sdk/788f363d15b1269ab5efb7cd00fb5e3b133cd99b/facebook_business/adobjects/pagepost.py",
      "bytes": 26704,
      "sha256": "b0b2009a0667ad57914423284877b104068dca1af4bdfb7f6543b0bec41f0b16"
    },
    {
      "path": "facebook_business/api.py",
      "url": "https://raw.githubusercontent.com/facebook/facebook-python-business-sdk/788f363d15b1269ab5efb7cd00fb5e3b133cd99b/facebook_business/api.py",
      "bytes": 30015,
      "sha256": "2f2e537354859758faba0d89a910b2212b478e0b2aaee771bb1c239abf4c2c72"
    },
    {
      "path": "facebook_business/session.py",
      "url": "https://raw.githubusercontent.com/facebook/facebook-python-business-sdk/788f363d15b1269ab5efb7cd00fb5e3b133cd99b/facebook_business/session.py",
      "bytes": 2319,
      "sha256": "72dd2e43efacc27728222d6c176202f6ff8919a85ba19811ac71fa357bc01587"
    },
    {
      "path": "facebook_business/adobjects/objectparser.py",
      "url": "https://raw.githubusercontent.com/facebook/facebook-python-business-sdk/788f363d15b1269ab5efb7cd00fb5e3b133cd99b/facebook_business/adobjects/objectparser.py",
      "bytes": 3661,
      "sha256": "9634dee1363b3930469b74942a5b5f58a53bfef5f21bd91519a2e9a25f1de9be"
    }
  ],
  "claims": [
    {
      "path": "facebook_business/adobjects/page.py",
      "symbol": "Page.create_feed",
      "line": 2222,
      "end_line": 2363
    },
    {
      "path": "facebook_business/adobjects/pagepost.py",
      "symbol": "PagePost.api_get",
      "line": 140,
      "end_line": 170
    },
    {
      "path": "facebook_business/adobjects/pagepost.py",
      "symbol": "PagePost.api_delete",
      "line": 110,
      "end_line": 138
    },
    {
      "path": "facebook_business/api.py",
      "symbol": "FacebookAdsApi.call",
      "line": 236,
      "end_line": 339
    },
    {
      "path": "facebook_business/adobjects/objectparser.py",
      "symbol": "ObjectParser.parse_single",
      "line": 43,
      "end_line": 88
    }
  ],
  "declared_delta": "Revoke uses FacebookAdsApi.call DELETE node envelope because ObjectParser.parse_single intentionally discards success; success is True is mandatory. No retry, queue, moderation, scheduling or permission-approval implementation is claimed."
}
````

### FILE: `meta_page_write/LICENSE.Meta.txt`
```yaml
block_id: "PY-OFFICIAL-META-PAGE-WRITE:file6:v1"
operation: CREATE
provenance: VERBATIM
source: "Meta official SDK LICENSE at locked commit788f363d15b1269ab5efb7cd00fb5e3b133cd99b"
license: "LicenseRef-Meta-Platform"
sha256: "48d97b3c936203a750a3288c2b924327769d62327383a1060244b3d3c05ad01f"
variables: []
secrets_allowed: false
```
````text
Copyright (c) Meta Platforms, Inc. and affiliates.
All rights reserved.

You are hereby granted a non-exclusive, worldwide, royalty-free license to use,
copy, modify, and distribute this software in source code or binary form for use
in connection with the web services and APIs provided by Facebook.

As with any software that integrates with the Facebook platform, your use of
this software is subject to the Facebook Platform Policy
[http://developers.facebook.com/policy/]. This copyright notice shall be
included in all copies or substantial portions of the software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
````

### FILE: `meta_page_write/README.md`
```yaml
block_id: "PY-OFFICIAL-META-PAGE-WRITE:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "Local integration/verification/dependency declaration; official SDK code is dependency-pinned, not locally authored"
license: "LicenseRef-Workspace-Owner"
sha256: "916b79a011cc78a532f5072a5581f5b4cd29b51599152c2ae60b4abb42666db1"
variables: []
secrets_allowed: false
```
````markdown
# Official Meta Page write leaf

This AUTHORED adapter invokes the official Meta Business SDK 26.0.1 pinned by
exact wheel graph and commit. It provides one POST of approved text to the selected
Facebook Page, bound GET observation/reconciliation, and separately approved
revocation. It adds no Instagram, TikTok, LinkedIn, ad mutation, media or paid boost.

Install the exact existing CPython 3.14 Windows lock into an isolated environment:

    python -m pip install --require-hashes --no-deps -r meta_page_write/requirements-windows-py314.lock
    python -m pip check
    python -m unittest discover -s meta_page_write/tests -v

The dependency lock is unchanged from the admitted Meta Ads reporting runtime;
its read-only adapter is not reused or relabeled as a write adapter. Meta licenses
permit use with Meta APIs; LICENSE.Meta.txt is mandatory, including for the CAPI
wheel whose license file is absent. This is not an OSI/general-purpose license.

Create an API with create_api(app_id, app_secret, page_access_token), then construct
FacebookPageWriteAdapter(Scope(...), api). Credentials resolve in the caller's
secret manager/environment; this module has no CLI or embedded credentials. Graph
v26.0 is fixed; debug logging is rejected. Calls use the official SDK transport.
Tests mount an in-memory requests transport, so no network publication occurs.

The caller supplies ApprovedIntent only after its existing identity, approval,
revocation, scheduling, quota and durable outbound-fence owners authorize it.
The approval ID is a correlation binding, not proof that this leaf has performed
authorization. Tenant/page/profile/content hashes must match the configured
scope; content is preserved byte-for-byte. The 16384-byte ceiling is a local
transport/memory bound, not a claim about provider character limits.

publish_once performs one POST then a GET. A receipt requires the exact post ID,
Page owner, message and published flag from the GET. A timeout or mismatched result
is UnknownDelivery; do not resend. With a known ID, reconcile performs GET only.
Without an ID, keep the durable fence unresolved pending owner reconciliation;
this adapter never guesses a post from similar text. revoke_once requires a
separate revoke intent, checks the original bytes/Page first, then requires raw
SDK DELETE success=True. It does not infer deletion from an empty object or a
permission/error response. Receipt hashes cover the normalized verified fields,
not an unmodified raw HTTP body. Requests and secrets are not included in receipts.

Scope still open for the complete social core: durable schedule/approval binding,
lease/quotas, dispatch worker wiring to communication.outbound_delivery, receipt
persistence, notification/UI, operator reconciliation and target token lifecycle.
These are local code/integration gaps and are not only missing credentials.
The existing outbound fence may be reused; this leaf does not claim that wiring
has been performed. The six other core gaps are listed in seven-core-decisions.md.

Official Page SDK methods were checked at their immutable commit. Current Meta
Pages API/permissions documentation endpoints returned tool retrieval errors;
no third-party mirror was used as authority. Live permission/app/Page ownership
contract must be verified before enabling a target. No live-readiness claim.

Verification: seven focused tests with real locked SDK objects/requests transport,
including wrong tenant/page/profile/content, unconfirmed observation, GET-only
reconciliation, POST timeout and false/missing/timeout DELETE outcomes. Source
files installed from the wheel match the fixed official source. First revoke
failure and correction are preserved in PROJECT_FAILURE_LESSONS.md. No core,
Commerce, Accounting, PG or previous fuzz suite was rerun.
````

## 6. Configuration and runtime

No embedded credential, CLI send command or global default API. The caller supplies
app/page credentials from its own secret handling when it chooses to enable the
adapter; fixtures use synthetic values and intercept all requests. Graph v26.0,
SDK26.0.1, per-call timeout10s, local content bound16384bytes; debug disabled.

## 7. Dependencies and license

The previously admitted Meta reporting 18-wheel lock is reused exactly, with no
new version or package. Clean isolated hash-required install, pip check and source
comparisons passed. Source/license mapping is materialized. Preserve all wheel
licenses and LICENSE.Meta.txt; the CAPI wheel omits its source license.

## 8. Operation and rollback

At most2 SDK requests per publish/revoke operation and1 per reconciliation; no
queue/retry engine. The caller stores receipt or unknown state using its existing
fence before further work, owns token lifecycle and per-provider quotas. Rollback
stops dispatch before removing the adapter and retains unknown receipts for
operator recovery. No full-core deployment or live permission acceptance claimed.

## 9. Verification

Seven focused real-SDK transport fixture tests. Prior seven cores and14 source
blocks match V374 exactly; their tests are reused without rerun. The initial
DELETE parser assumption failure is retained and corrected, with false/missing/
timeout regressions. No live API request, PG or previous domain/fuzz suite runs.

## 10. Reconstruction

Materialize into an absent directory; compare all7 outputs with these hashes.
Use the locked runtime and unittest discovery command in README. This candidate
requires independent root admission/wiring before claiming full social delivery.
