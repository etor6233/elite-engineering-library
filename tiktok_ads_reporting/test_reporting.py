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
