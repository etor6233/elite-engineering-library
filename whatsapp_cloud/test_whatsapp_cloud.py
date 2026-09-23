from __future__ import annotations

import hashlib
import hmac
import json
from pathlib import Path
import tempfile
import unittest
import subprocess
import sys

from whatsapp_cloud import build_template_payload, send_template_to_evidence, verify_packaged_official_source, verify_subscription, verify_webhook_signature, webhook_to_evidence, send_bridge


def profile():
    return {"provider": "Meta WhatsApp Cloud API", "source_id": "meta-whatsapp-api-examples", "source_commit": "de70ee908a67026e642aaee3703d20464e2a9466", "decision": "PROVEN", "official_source_license_accepted": True, "platform_terms_accepted": True, "business_account_proven": True, "app_registration_proven": True, "phone_number_id_proven": True, "message_template_approval_proven": True, "webhook_subscription_proven": True, "test_recipient_consent_proven": True, "quota_and_cost_approved": True, "reconciliation_approved": True, "business_account_id": "987654321", "graph_api_version": "v99.0", "phone_number_id": "123456789", "access_token_environment_variable": "WHATSAPP_ACCESS_TOKEN", "app_secret_environment_variable": "META_APP_SECRET", "verify_token_environment_variable": "WHATSAPP_VERIFY_TOKEN", "approved_templates": [{"name": "order_update", "language_code": "es_AR", "body_parameter_count": 2}], "data_retention": "30 days", "automatic_business_write": False}


def message():
    return {"recipient": "5491112345678", "template_name": "order_update", "language_code": "es_AR", "body_parameters": ["A-1", "despachado"]}


class WhatsAppCloudTests(unittest.TestCase):
    def test_signed_foreign_scope_is_rejected(self):
        payload = {"object": "whatsapp_business_account", "entry": [{"id": "999999999", "changes": [{"field": "messages", "value": {"messaging_product": "whatsapp", "metadata": {"phone_number_id": "999999999"}, "statuses": [{"id": "wamid.1", "status": "delivered", "timestamp": "1603086313", "recipient_id": "5491112345678"}]}}]}]}
        raw = json.dumps(payload).encode()
        signature = "sha256=" + hmac.new(b"app-secret", raw, hashlib.sha256).hexdigest()
        with tempfile.TemporaryDirectory() as temp:
            with self.assertRaises(PermissionError):
                webhook_to_evidence(raw, signature, "app-secret", Path(temp) / "evidence", profile=profile())

    def test_packaged_official_source_hashes(self):
        root = Path(__file__).parent
        verify_packaged_official_source(root, root / "official-source.lock.json")

    def test_profile_and_template_fail_closed(self):
        blocked = profile(); blocked["decision"] = "BLOCKED"
        with self.assertRaises(PermissionError): build_template_payload(blocked, message())
        bad = message(); bad["body_parameters"] = ["only-one"]
        with self.assertRaises(PermissionError): build_template_payload(profile(), bad)

    def test_send_preserves_response_and_redacts_identifiers(self):
        calls = []
        def transport(method, url, headers, body, timeout):
            calls.append((method, url, headers, body, timeout))
            return 200, {"content-type": "application/json"}, b'{"messages":[{"id":"wamid.secret-id"}]}'
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = send_template_to_evidence(profile(), message(), "secret-token", output, transport)
            serialized = json.dumps(receipt)
            self.assertNotIn("5491112345678", serialized); self.assertNotIn("123456789", serialized); self.assertNotIn("secret", serialized)
            self.assertEqual(calls[0][0], "POST"); self.assertEqual(calls[0][1], "https://graph.facebook.com/v99.0/123456789/messages")
            self.assertEqual(calls[0][2]["Authorization"], "Bearer secret-token")
            self.assertTrue((output / "provider-response.json").is_file())

    def test_send_provider_failure_is_atomic(self):
        def transport(*_): return 403, {}, b'{"error":{"message":"denied"}}'
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            with self.assertRaisesRegex(RuntimeError, "HTTP 403"):
                send_template_to_evidence(profile(), message(), "token", output, transport)
            self.assertFalse(output.exists()); self.assertEqual(list(Path(temp).iterdir()), [])

    def test_subscription_is_exact(self):
        self.assertEqual(verify_subscription({"hub.mode": "subscribe", "hub.verify_token": "verify", "hub.challenge": "42"}, "verify"), "42")
        with self.assertRaises(PermissionError): verify_subscription({"hub.mode": "subscribe", "hub.verify_token": "wrong", "hub.challenge": "42"}, "verify")

    def test_signature_uses_raw_body_and_rejects_tamper(self):
        body, secret = b'{"object":"whatsapp_business_account"}', "app-secret"
        signature = "sha256=" + hmac.new(secret.encode(), body, hashlib.sha256).hexdigest()
        verify_webhook_signature(body, signature, secret)
        with self.assertRaises(PermissionError): verify_webhook_signature(body + b" ", signature, secret)

    def test_webhook_normalizes_and_redacts(self):
        payload = {"object": "whatsapp_business_account", "entry": [{"id": "987654321", "changes": [{"field": "messages", "value": {"messaging_product": "whatsapp", "metadata": {"phone_number_id": "123456789"}, "messages": [{"id": "wamid.1", "from": "5491112345678", "type": "text", "timestamp": "1", "text": {"body": "private"}}], "statuses": [{"id": "wamid.2", "status": "delivered", "timestamp": "2", "recipient_id": "5491112345678"}]}}]}]}
        raw, secret = json.dumps(payload, separators=(",", ":")).encode(), "app-secret"
        signature = "sha256=" + hmac.new(secret.encode(), raw, hashlib.sha256).hexdigest()
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = webhook_to_evidence(raw, signature, secret, output, profile=profile())
            normalized = (output / "normalized-events.json").read_text(encoding="utf-8")
            self.assertEqual(receipt["event_count"], 2); self.assertFalse(receipt["raw_payload_persisted"])
            self.assertNotIn("5491112345678", normalized); self.assertNotIn("private", normalized)

    def test_webhook_invalid_envelope_is_atomic(self):
        raw, secret = b'{"object":"wrong","entry":[]}', "app-secret"
        signature = "sha256=" + hmac.new(secret.encode(), raw, hashlib.sha256).hexdigest()
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            with self.assertRaises(ValueError): webhook_to_evidence(raw, signature, secret, output, profile=profile())
            self.assertFalse(output.exists())


class ScopedWebhookTests(unittest.TestCase):
    def payload(self):
        return {"object": "whatsapp_business_account", "entry": [{"id": "987654321", "changes": [{"field": "messages", "value": {"messaging_product": "whatsapp", "metadata": {"phone_number_id": "123456789"}, "statuses": [{"id": "wamid.example", "status": "delivered", "timestamp": "1603086313", "recipient_id": "5491112345678"}]}}]}]}

    def run_payload(self, payload, output, configured=None):
        raw = payload if isinstance(payload, bytes) else json.dumps(payload).encode()
        signature = "sha256=" + hmac.new(b"app-secret", raw, hashlib.sha256).hexdigest()
        return webhook_to_evidence(raw, signature, "app-secret", output, profile=configured or profile())

    def test_independent_account_and_phone_scope_rejection(self):
        for target in ("account", "phone", "missing_account", "missing_phone", "mixed_batch"):
            with self.subTest(target=target), tempfile.TemporaryDirectory() as temp:
                payload = self.payload()
                entry = payload["entry"][0]
                metadata = entry["changes"][0]["value"]["metadata"]
                if target == "account": entry["id"] = "999999999"
                if target == "phone": metadata["phone_number_id"] = "999999999"
                if target == "missing_account": del entry["id"]
                if target == "missing_phone": metadata.clear()
                if target == "mixed_batch":
                    foreign = self.payload()["entry"][0]
                    foreign["id"] = "999999999"
                    payload["entry"].append(foreign)
                with self.assertRaises(PermissionError):
                    self.run_payload(payload, Path(temp) / "evidence")
                self.assertEqual(list(Path(temp).iterdir()), [])

    def test_malformed_status_is_atomic(self):
        variants = [("id", ""), ("id", "a" * 513), ("id", "private\ntext"), ("recipient_id", None), ("recipient_id", ""), ("timestamp", 1603086313), ("timestamp", "-1"), ("timestamp", "01"), ("timestamp", ""), ("status", "accepted"), ("status", "unknown"), ("status", {})]
        for key, value in variants:
            with self.subTest(key=key, value=value), tempfile.TemporaryDirectory() as temp:
                payload = self.payload()
                payload["entry"][0]["changes"][0]["value"]["statuses"][0][key] = value
                with self.assertRaises(ValueError): self.run_payload(payload, Path(temp) / "evidence")
                self.assertEqual(list(Path(temp).iterdir()), [])

    def test_malformed_containers_rejected(self):
        for target in ("entry", "changes", "metadata", "product", "statuses", "messages", "field"):
            with self.subTest(target=target), tempfile.TemporaryDirectory() as temp:
                payload = self.payload()
                entry = payload["entry"][0]
                change = entry["changes"][0]
                value = change["value"]
                if target == "entry": payload["entry"] = {}
                elif target == "changes": entry["changes"] = None
                elif target == "field": change["field"] = "other"
                elif target == "product": value["messaging_product"] = "other"
                else: value[target] = None
                with self.assertRaises(ValueError): self.run_payload(payload, Path(temp) / "evidence")

    def test_scope_profile_cannot_be_omitted_or_unproven(self):
        with tempfile.TemporaryDirectory() as temp:
            for key, value in (("business_account_id", ""), ("phone_number_id", 123456789), ("decision", "BLOCKED")):
                configured = profile(); configured[key] = value
                with self.subTest(key=key), self.assertRaises((ValueError, PermissionError)):
                    self.run_payload(self.payload(), Path(temp) / "evidence", configured)
            with self.assertRaises(TypeError):
                webhook_to_evidence(b"{}", "signature", "app-secret", Path(temp) / "evidence")
            self.assertEqual(list(Path(temp).iterdir()), [])

    def test_duplicate_json_keys_rejected(self):
        raw = json.dumps(self.payload()).replace('"status": "delivered"', '"status": "sent", "status": "delivered"').encode()
        with tempfile.TemporaryDirectory() as temp, self.assertRaisesRegex(ValueError, "duplicate JSON"):
            self.run_payload(raw, Path(temp) / "evidence")

    def test_local_body_and_event_budgets(self):
        with tempfile.TemporaryDirectory() as temp:
            with self.assertRaisesRegex(ValueError, "body budget"):
                self.run_payload(b" " * (1024 * 1024 + 1), Path(temp) / "body")
            payload = self.payload()
            value = payload["entry"][0]["changes"][0]["value"]
            value["statuses"] *= 1001
            with self.assertRaisesRegex(ValueError, "bounded lists"):
                self.run_payload(payload, Path(temp) / "events")
            value["statuses"] = value["statuses"][:600]
            payload["entry"][0]["changes"] *= 2
            with self.assertRaisesRegex(ValueError, "event budget"):
                self.run_payload(payload, Path(temp) / "aggregate")
            self.assertEqual(list(Path(temp).iterdir()), [])

    def test_status_identity_preserved_without_ordering_or_delivery_inference(self):
        with tempfile.TemporaryDirectory() as temp:
            payload = self.payload()
            statuses = payload["entry"][0]["changes"][0]["value"]["statuses"]
            statuses[:] = [{**statuses[0], "status": state} for state in ("read", "sent", "failed", "delivered", "deleted")]
            output = Path(temp) / "first"
            receipt = self.run_payload(payload, output)
            events = json.loads((output / "normalized-events.json").read_bytes())["events"]
            self.assertEqual([e["status"] for e in events], ["read", "sent", "failed", "delivered", "deleted"])
            self.assertEqual(len({e["event_key_sha256"] for e in events}), 5)
            self.assertEqual(receipt["schema"], "elite-whatsapp-cloud-webhook-receipt/v2")
            self.assertFalse(receipt["automatic_business_write"])
            self.assertEqual(receipt["normalized_sha256"], hashlib.sha256((output / "normalized-events.json").read_bytes()).hexdigest())
            # Format/order changes must not invent new provider-event identities.
            statuses.reverse()
            second = Path(temp) / "second"
            self.run_payload(json.dumps(payload, indent=2).encode(), second)
            later = json.loads((second / "normalized-events.json").read_bytes())["events"]
            self.assertEqual({e["event_key_sha256"] for e in events}, {e["event_key_sha256"] for e in later})
            for value in ("987654321", "123456789", "5491112345678", "wamid.example"):
                self.assertNotIn(value, json.dumps(events) + json.dumps(receipt))
            self.assertEqual(events[0]["recipient_sha256"], hashlib.sha256(b"5491112345678").hexdigest())

    def test_existing_output_is_not_overwritten(self):
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            self.run_payload(self.payload(), output)
            before = (output / "normalized-events.json").read_bytes()
            with self.assertRaises(FileExistsError): self.run_payload(self.payload(), output)
            self.assertEqual(before, (output / "normalized-events.json").read_bytes())


class SendBridgeTests(unittest.TestCase):
    def frame(self, output):
        return {"schema": "elite-whatsapp-send-bridge/v1", "binding_sha256": "b" * 64, "profile": profile(), "request": message(), "access_token": "synthetic-secret", "output_directory": str(output)}

    def test_bridge_uses_existing_adapter_and_binds_receipt(self):
        calls = []
        def transport(method, url, headers, body, timeout):
            calls.append((method, url, json.loads(body)))
            return 200, {}, b'{"messages":[{"id":"wamid.bridge"}]}'
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            result = send_bridge(json.dumps(self.frame(output)).encode(), transport)
            self.assertEqual(len(calls), 1)
            self.assertEqual(calls[0][0], "POST")
            self.assertEqual(calls[0][2]["to"], message()["recipient"])
            self.assertEqual(result["binding_sha256"], "b" * 64)
            self.assertEqual(result["provider_message_id"], "wamid.bridge")
            self.assertEqual(result["evidence_sha256"], hashlib.sha256((output / "SEND_RECEIPT.json").read_bytes()).hexdigest())
            self.assertNotIn("synthetic-secret", json.dumps(result))
            with self.assertRaises(FileExistsError): send_bridge(json.dumps(self.frame(output)).encode(), transport)
            self.assertEqual(len(calls), 1)

    def test_bridge_rejects_frame_without_provider(self):
        def forbidden(*args): self.fail("provider must not be invoked")
        with tempfile.TemporaryDirectory() as temp:
            for key, value in (("schema", "wrong"), ("binding_sha256", "bad"), ("profile", []), ("request", []), ("access_token", ""), ("output_directory", "relative")):
                frame = self.frame(Path(temp) / "evidence"); frame[key] = value
                with self.subTest(key=key), self.assertRaises(ValueError): send_bridge(json.dumps(frame).encode(), forbidden)
            frame = self.frame(Path(temp) / "evidence"); frame["extra"] = True
            with self.assertRaises(ValueError): send_bridge(json.dumps(frame).encode(), forbidden)
            with self.assertRaises(ValueError): send_bridge(b" " * 131073, forbidden)

    def test_uncertain_provider_never_retried_inside_bridge(self):
        calls = []
        def lost(*args): calls.append(1); raise TimeoutError("synthetic private transport")
        with tempfile.TemporaryDirectory() as temp, self.assertRaises(TimeoutError):
            send_bridge(json.dumps(self.frame(Path(temp) / "evidence")).encode(), lost)
        self.assertEqual(calls, [1])

    def test_ambiguous_provider_identity_is_not_accepted(self):
        with tempfile.TemporaryDirectory() as temp:
            for index, body in enumerate((b'{"messages":[{"id":"one"},{"id":"two"}]}', b'{"messages":[{"id":"bad\\nidentity"}]}')):
                with self.subTest(index=index), self.assertRaises(ValueError):
                    send_bridge(json.dumps(self.frame(Path(temp) / str(index))).encode(), lambda *args: (200, {}, body))

    def test_real_cli_fails_closed_without_logging_secret(self):
        frame = self.frame(Path(tempfile.gettempdir()) / "must-not-be-created-whatsapp")
        frame["profile"]["decision"] = "BLOCKED"
        result = subprocess.run([sys.executable, "-I", "-B", str(Path(__file__).with_name("whatsapp_cloud.py")), "--send-bridge"], input=json.dumps(frame).encode(), capture_output=True, timeout=10)
        self.assertEqual(result.returncode, 2)
        self.assertEqual(result.stdout, b"")
        self.assertEqual(result.stderr, b"WHATSAPP_SEND_UNCERTAIN_OR_REJECTED\n")


if __name__ == "__main__": unittest.main()
