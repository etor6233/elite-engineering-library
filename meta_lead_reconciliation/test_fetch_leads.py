from __future__ import annotations

import inspect
import json
from pathlib import Path
import tempfile
import unittest

from facebook_business.adobjects.leadgenform import LeadgenForm
from fetch_leads import fetch_to_evidence, validate_configuration


class FakeRow:
    def __init__(self, value): self.value = value
    def export_all_data(self): return self.value


class FakeForm:
    def __init__(self, identifier, rows=None, error=None):
        self.identifier = identifier
        self.rows = rows or []
        self.error = error
        self.calls = []

    def _run(self, mode, fields, params):
        self.calls.append((mode, fields, params))
        if self.error: raise self.error
        return self.rows

    def get_leads(self, *, fields, params): return self._run("LIVE", fields, params)
    def get_test_leads(self, *, fields, params): return self._run("TEST", fields, params)


def proven_profile():
    return {
        "provider": "Meta Marketing API", "sdk": "facebook-business==26.0.1", "graph_api_version": "v26.0",
        "decision": "PROVEN", "app_registration_proven": True, "business_verification_proven": True,
        "lead_access_permission_proven": True, "page_and_form_ownership_proven": True,
        "test_lead_contract_proven": True, "pii_storage_controls_proven": True,
        "quota_and_cost_approved": True, "reconciliation_approved": True,
        "automatic_business_write": False, "automatic_contact_eligibility": False,
        "data_retention": "30 days", "tenant_id": "tenant-1", "organization_id": "org-1",
        "form_id": "123456789", "approved_provider_fields": ["id", "created_time", "form_id", "field_data", "campaign_id", "adset_id", "ad_id", "platform"],
        "approved_form_field_names": ["email", "full_name"],
    }


def valid_retrieval(mode="LIVE"):
    return {"retrieval_mode": mode, "fields": ["id", "created_time", "form_id", "field_data", "campaign_id"], "max_rows": 100}


def valid_row():
    return {
        "id": "lead-1", "created_time": "2026-09-04T12:00:00+0000", "form_id": "123456789",
        "campaign_id": "campaign-1", "field_data": [
            {"name": "email", "values": ["person@example.test"]},
            {"name": "full_name", "values": ["Test Person"]},
        ],
    }


class MetaLeadReconciliationTests(unittest.TestCase):
    def test_configuration_fails_closed(self):
        profile = proven_profile(); profile["decision"] = "BLOCKED_ACCESS_AND_POLICY_APPROVAL_REQUIRED"
        with self.assertRaises(PermissionError): validate_configuration(profile, valid_retrieval())
        profile = proven_profile(); profile["automatic_contact_eligibility"] = True
        with self.assertRaises(ValueError): validate_configuration(profile, valid_retrieval())

    def test_configuration_rejects_unapproved_fields_and_invalid_form(self):
        profile = proven_profile(); query = valid_retrieval(); query["fields"].append("ad_name")
        with self.assertRaises(PermissionError): validate_configuration(profile, query)
        profile = proven_profile(); profile["form_id"] = "form_123"
        with self.assertRaises(ValueError): validate_configuration(profile, valid_retrieval())

    def test_live_fetch_preserves_provider_and_emits_pending_candidate(self):
        form = FakeForm("", [FakeRow(valid_row())])
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = fetch_to_evidence(lambda identifier: self._bind(form, identifier), "tenant-1", "org-1", "123456789", valid_retrieval()["fields"], {"email", "full_name"}, 100, "LIVE", output)
            self.assertEqual(form.calls[0][0], "LIVE")
            self.assertEqual(receipt["candidate_count"], 1)
            self.assertNotIn("123456789", json.dumps(receipt))
            self.assertNotIn("person@example.test", json.dumps(receipt))
            provider = json.loads((output / "provider-response.json").read_text(encoding="utf-8"))
            candidate = json.loads((output / "lead-candidates.json").read_text(encoding="utf-8"))["outcomes"][0]["candidate"]
            batch = json.loads((output / "lead-candidates.json").read_text(encoding="utf-8"))
            self.assertEqual(provider["rows"][0]["field_data"][0]["values"][0], "person@example.test")
            self.assertEqual(batch["tenant_id"], "tenant-1")
            self.assertEqual(batch["organization_id"], "org-1")
            self.assertEqual(candidate["contact_eligibility"], "pending_policy")
            self.assertFalse(candidate["is_test"])
            self.assertEqual(candidate["fields"][0]["value"], "person@example.test")

    def test_test_fetch_marks_candidates_test(self):
        form = FakeForm("", [valid_row()])
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            fetch_to_evidence(lambda identifier: self._bind(form, identifier), "tenant-1", "org-1", "123456789", valid_retrieval("TEST")["fields"], {"email", "full_name"}, 100, "TEST", output)
            candidate = json.loads((output / "lead-candidates.json").read_text(encoding="utf-8"))["outcomes"][0]["candidate"]
            self.assertEqual(form.calls[0][0], "TEST")
            self.assertTrue(candidate["is_test"])

    def test_unknown_form_field_is_rejected_without_losing_raw(self):
        row = valid_row(); row["field_data"].append({"name": "unapproved", "values": ["secret"]})
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = fetch_to_evidence(lambda identifier: FakeForm(identifier, [row]), "tenant-1", "org-1", "123456789", valid_retrieval()["fields"], {"email", "full_name"}, 100, "LIVE", output)
            outcome = json.loads((output / "lead-candidates.json").read_text(encoding="utf-8"))["outcomes"][0]
            self.assertEqual(receipt["rejected_count"], 1)
            self.assertEqual(outcome["normalization_codes"], ["UNAPPROVED_FORM_FIELD"])
            self.assertIsNone(outcome["candidate"])
            self.assertIn("unapproved", (output / "provider-response.json").read_text(encoding="utf-8"))

    def test_multiple_values_are_retained_raw_and_rejected_as_ambiguous(self):
        row = valid_row(); row["field_data"][0]["values"] = ["first@example.test", "second@example.test"]
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = fetch_to_evidence(lambda identifier: FakeForm(identifier, [row]), "tenant-1", "org-1", "123456789", valid_retrieval()["fields"], {"email", "full_name"}, 100, "LIVE", output)
            outcome = json.loads((output / "lead-candidates.json").read_text(encoding="utf-8"))["outcomes"][0]
            self.assertEqual(receipt["rejected_count"], 1)
            self.assertEqual(outcome["normalization_codes"], ["INVALID_FIELD_VALUE"])
            self.assertIsNone(outcome["candidate"])
            self.assertIn("second@example.test", (output / "provider-response.json").read_text(encoding="utf-8"))

    def test_local_bound_is_explicit(self):
        rows = [valid_row(), {**valid_row(), "id": "lead-2"}]
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = fetch_to_evidence(lambda identifier: FakeForm(identifier, rows), "tenant-1", "org-1", "123456789", valid_retrieval()["fields"], {"email", "full_name"}, 1, "LIVE", output)
            self.assertEqual(receipt["row_count"], 1)
            self.assertTrue(receipt["truncated_by_local_bound"])

    def test_provider_failure_is_atomic(self):
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            with self.assertRaisesRegex(RuntimeError, "provider failed"):
                fetch_to_evidence(lambda identifier: FakeForm(identifier, error=RuntimeError("provider failed")), "tenant-1", "org-1", "123456789", valid_retrieval()["fields"], {"email"}, 1, "LIVE", output)
            self.assertFalse(output.exists())
            self.assertEqual(list(Path(temp).iterdir()), [])

    def test_existing_output_is_rejected(self):
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"; output.mkdir()
            with self.assertRaises(FileExistsError):
                fetch_to_evidence(lambda identifier: FakeForm(identifier), "tenant-1", "org-1", "123456789", valid_retrieval()["fields"], {"email"}, 1, "LIVE", output)

    def test_official_sdk_contract(self):
        for method_name in ("get_leads", "get_test_leads"):
            method = getattr(LeadgenForm, method_name)
            signature = inspect.signature(method)
            self.assertIn("fields", signature.parameters)
            self.assertIn("params", signature.parameters)
        self.assertIn("/leads", inspect.getsource(LeadgenForm.get_leads))
        self.assertIn("/test_leads", inspect.getsource(LeadgenForm.get_test_leads))

    @staticmethod
    def _bind(form, identifier): form.identifier = identifier; return form


if __name__ == "__main__": unittest.main()
