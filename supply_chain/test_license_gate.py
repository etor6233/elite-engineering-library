import json
import tempfile
import unittest
from argparse import Namespace
from pathlib import Path

from license_gate import GateError, run


class LicenseGateTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.root = Path(self.temp.name)
        self.policy = self.root / "policy.json"
        self.policy.write_text(json.dumps({"allowed_expressions": ["MIT", "Apache-2.0"], "acknowledged_expressions": {}, "denied_expressions": ["AGPL-3.0-only"]}), encoding="utf-8")

    def tearDown(self):
        self.temp.cleanup()

    def test_generates_hashed_pnpm_and_go_evidence(self):
        package = self.root / "package"
        package.mkdir()
        (package / "LICENSE").write_text("license text", encoding="utf-8")
        pnpm = self.root / "pnpm.json"
        pnpm.write_text(json.dumps({"MIT": [{"name": "example", "versions": ["1.0.0"], "paths": [str(package)]}]}), encoding="utf-8")
        notices = self.root / "notices"
        notices.mkdir()
        (notices / "LICENSE").write_text("go license", encoding="utf-8")
        go_report = self.root / "go.csv"
        go_report.write_text("example.org/module,https://example.org/LICENSE,Apache-2.0\n", encoding="utf-8")
        output = self.root / "evidence.json"
        result = run(Namespace(policy=str(self.policy), output=str(output), pnpm_report=str(pnpm), go_report=str(go_report), go_notices=str(notices)))
        self.assertEqual(result["status"], "PASS")
        self.assertTrue(output.exists())
        self.assertIn("sha256", output.read_text(encoding="utf-8"))

    def test_rejects_unreviewed_denied_and_missing_artifacts(self):
        package = self.root / "package"
        package.mkdir()
        report = self.root / "pnpm.json"
        for expression in ("LGPL-3.0-or-later", "AGPL-3.0-only"):
            report.write_text(json.dumps({expression: [{"name": "example", "versions": ["1"], "paths": [str(package)]}]}), encoding="utf-8")
            with self.assertRaises(GateError):
                run(Namespace(policy=str(self.policy), output=str(self.root / "out.json"), pnpm_report=str(report), go_report=None, go_notices=None))
        report.write_text(json.dumps({"MIT": [{"name": "example", "versions": ["1"], "paths": [str(package)]}]}), encoding="utf-8")
        with self.assertRaises(GateError):
            run(Namespace(policy=str(self.policy), output=str(self.root / "out.json"), pnpm_report=str(report), go_report=None, go_notices=None))


if __name__ == "__main__":
    unittest.main()
