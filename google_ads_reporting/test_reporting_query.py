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
