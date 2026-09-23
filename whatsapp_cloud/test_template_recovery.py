"""AUTHORED regression: preserve exact template acceptance without a resend."""
import base64
import copy
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch

from test_whatsapp_cloud import profile, message
from whatsapp_cloud import recover_send_bridge, send_bridge


class TemplateRecoveryTests(unittest.TestCase):
    def test_exact_template_and_tampered_bindings(self):
        calls = []
        def transport(*args):
            calls.append(args)
            return 200, {}, b'{"messages":[{"id":"wamid.template.recovered"}]}'
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "receipt"
            sent = send_bridge(json.dumps({
                "schema": "elite-whatsapp-send-bridge/v1",
                "binding_sha256": "a"*64, "profile": profile(), "request": message(),
                "access_token": "synthetic-token", "output_directory": str(output)
            }).encode(), transport)
            frame = {
                "schema": "elite-whatsapp-recover-send/v1",
                "binding_sha256": "a"*64, "profile": profile(), "request": message(),
                "receipt": base64.b64encode((output/"SEND_RECEIPT.json").read_bytes()).decode(),
                "response": base64.b64encode((output/"provider-response.json").read_bytes()).decode()
            }
            with patch("whatsapp_cloud.stdlib_transport", side_effect=AssertionError("recovery attempted network")):
                self.assertEqual(recover_send_bridge(json.dumps(frame).encode()), sent)
                cases = []
                changed = copy.deepcopy(frame); changed["request"]["body_parameters"][0] = "another-order"; cases.append(changed)
                changed = copy.deepcopy(frame); changed["request"]["recipient"] = "5491112345679"; cases.append(changed)
                changed = copy.deepcopy(frame); changed["profile"]["phone_number_id"] = "123456788"; cases.append(changed)
                changed = copy.deepcopy(frame); changed["profile"]["graph_api_version"] = "v98.0"; cases.append(changed)
                changed = copy.deepcopy(frame); changed["request"]["kind"] = "unknown"; cases.append(changed)
                changed = copy.deepcopy(frame); changed["request"]["extra"] = True; cases.append(changed)
                changed = copy.deepcopy(frame); changed["request"] = []; cases.append(changed)
                changed = copy.deepcopy(frame); changed["response"] = base64.b64encode(b'{"messages":[{"id":"wamid.other"}]}').decode(); cases.append(changed)
                for index, changed in enumerate(cases):
                    with self.subTest(index=index), self.assertRaises(ValueError):
                        recover_send_bridge(json.dumps(changed).encode())
            self.assertEqual(len(calls), 1)


if __name__ == "__main__":
    unittest.main()
