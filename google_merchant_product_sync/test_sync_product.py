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
