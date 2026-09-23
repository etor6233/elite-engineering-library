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
