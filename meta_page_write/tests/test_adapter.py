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
