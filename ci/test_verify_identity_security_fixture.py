from __future__ import annotations

import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

from verify_identity_security_fixture import load_fixture, validate_fixture, verify_doors

ROOT = Path(__file__).resolve().parents[1]
FIXTURE = ROOT / "ci" / "fixtures" / "identity_authz_scenarios.json"


class VerifyIdentitySecurityFixtureTests(unittest.TestCase):
    def test_default_fixture_valid(self) -> None:
        data = load_fixture(FIXTURE)
        self.assertEqual([], validate_fixture(data))
        self.assertEqual([], verify_doors(ROOT, data["doors"]))

    def test_missing_door_detected(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            workspace = Path(tmp)
            data = json.loads(FIXTURE.read_text(encoding="utf-8"))
            data["doors"] = ["missing/door.md"]
            path = workspace / "fixture.json"
            path.write_text(json.dumps(data), encoding="utf-8")
            errors = verify_doors(workspace, data["doors"])
            self.assertTrue(any("missing door" in item for item in errors))

    def test_wildcard_scope_rejected(self) -> None:
        data = json.loads(FIXTURE.read_text(encoding="utf-8"))
        data["scope"] = "PRODUCTION"
        with self.assertRaises(ValueError):
            load_fixture_from_dict(data)

    def test_extra_fixture_field_rejected(self) -> None:
        data = json.loads(FIXTURE.read_text(encoding="utf-8"))
        data["description"] = "unsupported extra"
        with tempfile.NamedTemporaryFile("w", suffix=".json", delete=False) as handle:
            handle.write(json.dumps(data))
            path = Path(handle.name)
        try:
            with self.assertRaises(ValueError):
                load_fixture(path)
        finally:
            path.unlink(missing_ok=True)

    def test_cli_passes_on_workspace(self) -> None:
        completed = subprocess.run(
            [sys.executable, str(ROOT / "ci" / "verify_identity_security_fixture.py"), "--workspace", str(ROOT), "--skip-go"],
            capture_output=True,
            text=True,
            check=False,
        )
        self.assertEqual(0, completed.returncode, completed.stderr)
        self.assertIn("PASS", completed.stdout)


def load_fixture_from_dict(data: dict) -> dict:
    if data.get("schema") != "elite.identity-authz-fixture/v1":
        raise ValueError("fixture schema must be elite.identity-authz-fixture/v1")
    if data.get("scope") != "LOCAL_FIXTURES":
        raise ValueError("fixture scope must be LOCAL_FIXTURES")
    if data.get("live_effects") is not False:
        raise ValueError("fixture live_effects must be false")
    return data


if __name__ == "__main__":
    unittest.main()
