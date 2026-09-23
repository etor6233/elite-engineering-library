from __future__ import annotations

import hashlib
import hmac
import itertools
import json
from pathlib import Path
import tempfile
import unittest
import base64
import os
import subprocess
import sys
from unittest.mock import patch

from status_reconciliation import reconcile_status_webhooks
from test_whatsapp_cloud import message, profile
from whatsapp_cloud import send_bridge


def event(status="delivered", timestamp="1603086314", **extra):
    return {"id": "wamid.synthetic", "recipient_id": "5491112345678", "timestamp": timestamp, "status": status, **extra}


def webhook(events, *, account="987654321", phone="123456789"):
    raw = json.dumps({"object": "whatsapp_business_account", "entry": [{"id": account, "changes": [{"field": "messages", "value": {"messaging_product": "whatsapp", "metadata": {"phone_number_id": phone}, "statuses": events}}]}]}, separators=(",", ":")).encode()
    return raw, "sha256=" + hmac.new(b"test-app-secret", raw, hashlib.sha256).hexdigest()


class StatusReconciliationTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.calls = 0
        def transport(*args):
            self.calls += 1
            return 200, {}, b'{"messages":[{"id":"wamid.synthetic"}]}'
        frame = {"schema": "elite-whatsapp-send-bridge/v1", "binding_sha256": "a" * 64, "profile": profile(), "request": message(), "access_token": "test-access-token", "output_directory": str(self.root / "send")}
        self.send = send_bridge(json.dumps(frame).encode(), transport)
        self.receipt = (self.root / "send" / "SEND_RECEIPT.json").read_bytes()

    def reconcile(self, batches=None, **overrides):
        args = dict(profile=profile(), send_receipt_bytes=self.receipt,
                    expected_send_receipt_sha256=self.send["evidence_sha256"],
                    signed_webhooks=batches if batches is not None else [webhook([event()])],
                    app_secret="test-app-secret", output_directory=self.root / "result")
        args.update(overrides)
        return reconcile_status_webhooks(**args)

    def test_send_bridge_to_signed_statuses_uses_durable_receipt_anchor_without_resend(self):
        with patch("whatsapp_cloud.urlopen", side_effect=AssertionError("no provider calls allowed")):
            result = self.reconcile()
        self.assertEqual(result["last_observed_status"], "delivered")
        self.assertEqual(result["send_receipt_sha256"], self.send["evidence_sha256"])
        self.assertEqual(self.calls, 1)
        self.assertFalse(result["resend_authorized"])
        self.assertFalse(result["automatic_business_write"])
        observed = (self.root / "result" / "status-observations.json").read_bytes()
        self.assertEqual(hashlib.sha256(observed).hexdigest(), result["observations_sha256"])
        for secret in ("test-app-secret", "test-access-token", "5491112345678", "wamid.synthetic"):
            self.assertNotIn(secret, observed.decode() + json.dumps(result))

    def test_all_arrival_permutations_have_identical_observation_hash(self):
        events = [event("sent", "1603086313"), event("delivered", "1603086314"), event("read", "1603086315")]
        hashes = set()
        for index, order in enumerate(itertools.permutations(events)):
            result = self.reconcile([webhook([item]) for item in order], output_directory=self.root / str(index))
            self.assertEqual(result["last_observed_status"], "read")
            hashes.add(result["observations_sha256"])
        self.assertEqual(len(hashes), 1)

    def test_repeated_batch_and_event_are_deduplicated(self):
        batch = webhook([event(), event()])
        result = self.reconcile([batch, batch])
        self.assertEqual(result["matching_unique_events"], 1)
        self.assertEqual(result["duplicate_events"], 3)

    def test_same_timestamp_does_not_invent_precedence(self):
        result = self.reconcile([webhook([event("sent"), event("delivered")])])
        self.assertEqual(result["decision"], "AMBIGUOUS_LATEST_TIMESTAMP")
        self.assertIsNone(result["last_observed_status"])
        self.assertEqual(result["matching_unique_events"], 2)

    def test_other_message_is_not_evidence_of_our_delivery(self):
        result = self.reconcile([webhook([event(id="wamid.other")])])
        self.assertEqual(result["decision"], "NOT_OBSERVED_IN_INPUT")
        self.assertIsNone(result["last_observed_status"])
        self.assertEqual(result["excluded_events"], 1)

    def test_mixed_batch_does_not_correlate_other_message(self):
        result = self.reconcile([webhook([event(), event("read", id="wamid.other")])])
        self.assertEqual(result["last_observed_status"], "delivered")
        self.assertEqual(result["excluded_events"], 1)

    def test_foreign_account_phone_or_recipient_rejects_atomically(self):
        for batch in [webhook([event()], account="999999999"), webhook([event()], phone="999999999"), webhook([event(recipient_id="5491112345679")])]:
            with self.subTest(batch_digest=hashlib.sha256(batch[0]).hexdigest()):
                with self.assertRaises(PermissionError): self.reconcile([batch])
                self.assertFalse((self.root / "result").exists())

    def test_signature_tamper_even_in_unrelated_message_rejects_entire_batch(self):
        raw, sig = webhook([event(id="wamid.other")])
        with self.assertRaises(PermissionError): self.reconcile([(raw + b" ", sig)])
        self.assertFalse((self.root / "result").exists())

    def test_send_anchor_tamper_rejects_before_output(self):
        with self.assertRaises(PermissionError): self.reconcile(send_receipt_bytes=self.receipt + b" ")
        self.assertFalse((self.root / "result").exists())

    def test_invalid_anchor_type_and_shape(self):
        for value in (None, "", "z" * 64):
            with self.subTest(value=value), self.assertRaises(ValueError):
                self.reconcile(expected_send_receipt_sha256=value)

    def test_same_event_identity_with_different_payload_is_not_silently_deduplicated(self):
        with self.assertRaisesRegex(ValueError, "divergent"):
            self.reconcile([webhook([event("failed", errors=[{"code": 1}])]), webhook([event("failed", errors=[{"code": 2}])])])
        self.assertFalse((self.root / "result").exists())

    def test_failed_and_deleted_are_retained_without_authorizing_retry(self):
        result = self.reconcile([webhook([event("failed", "1603086313"), event("deleted", "1603086314")])])
        self.assertEqual(result["observed_statuses"], ["deleted", "failed"])
        self.assertFalse(result["resend_authorized"])

    def test_profile_mismatch(self):
        changed = profile(); changed["graph_api_version"] = "v98.0"
        with self.assertRaises(PermissionError): self.reconcile(profile=changed)

    def test_empty_excessive_and_malformed_input(self):
        for value in ([], [webhook([event()])] * 9, [(b"{}",)], [(b"{}", None)]):
            with self.subTest(size=len(value)), self.assertRaises(ValueError): self.reconcile(value)
        self.assertFalse((self.root / "result").exists())

    def test_existing_output_is_never_overwritten(self):
        self.reconcile()
        previous = (self.root / "result" / "STATUS_RECONCILIATION.json").read_bytes()
        with self.assertRaises(FileExistsError): self.reconcile()
        self.assertEqual((self.root / "result" / "STATUS_RECONCILIATION.json").read_bytes(), previous)

    def test_write_failure_does_not_publish_partial_result(self):
        with patch.object(Path, "write_text", side_effect=OSError("fixture disk unavailable")):
            with self.assertRaises(OSError): self.reconcile()
        self.assertFalse((self.root / "result").exists())
        self.assertEqual(list(self.root.glob(".whatsapp-status-stage-*")), [])

    def test_late_bad_batch_leaves_no_partial_observations(self):
        raw, sig = webhook([event("read")])
        with self.assertRaises(PermissionError):
            self.reconcile([webhook([event()]), (raw + b" ", sig)])
        self.assertFalse((self.root / "result").exists())

    def test_bounded_raw_input_before_signature_parsing(self):
        with self.assertRaises(ValueError):
            self.reconcile([(b" " * (1024 * 1024 + 1), "sha256=" + "0" * 64)])
        with self.assertRaises(ValueError):
            self.reconcile(send_receipt_bytes=b" " * 65537)

    def test_receipt_contract_is_checked_even_when_hash_matches(self):
        for changes in ({"schema": "foreign"}, {"recipient_sha256": ""}, {"automatic_business_write": True}):
            altered = json.loads(self.receipt); altered.update(changes)
            raw = json.dumps(altered).encode()
            with self.subTest(changes=changes), self.assertRaises(ValueError):
                self.reconcile(send_receipt_bytes=raw, expected_send_receipt_sha256=hashlib.sha256(raw).hexdigest())

    def test_status_bridge_isolated_process_and_safe_failure(self):
        raw, signature = webhook([event()])
        frame = {"schema": "elite-whatsapp-status-bridge/v1", "profile": profile(),
                 "send_receipt": base64.b64encode(self.receipt).decode(),
                 "expected_send_receipt_sha256": self.send["evidence_sha256"],
                 "webhooks": [{"body": base64.b64encode(raw).decode(), "signature": signature}],
                 "app_secret": "test-app-secret", "evidence_directory": str(self.root)}
        environment = {key: os.environ[key] for key in ("SystemRoot",) if key in os.environ}
        args = [sys.executable, "-I", "-B", str(Path(__file__).resolve().parent / "whatsapp_cloud.py"), "--status-bridge"]
        encoded = json.dumps(frame).encode()
        result = subprocess.run(args, input=encoded, capture_output=True, env=environment, timeout=10)
        self.assertEqual(result.returncode, 0, result.stderr)
        reply = json.loads(result.stdout)
        self.assertEqual(reply["binding_sha256"], hashlib.sha256(encoded).hexdigest())
        self.assertEqual(reply["receipt"]["last_observed_status"], "delivered")
        self.assertEqual(list(self.root.glob("whatsapp-status-*")), [])
        frame["webhooks"][0]["signature"] = "sha256=" + "0" * 64
        result = subprocess.run(args, input=json.dumps(frame).encode(), capture_output=True, env=environment, timeout=10)
        self.assertEqual(result.returncode, 2)
        self.assertEqual(result.stdout, b"")
        self.assertEqual(result.stderr, b"WHATSAPP_STATUS_UNVERIFIED\n")


if __name__ == "__main__":
    unittest.main()
