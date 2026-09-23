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
