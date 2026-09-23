from __future__ import annotations

import copy
import json
import tempfile
import unittest
from unittest.mock import patch
from pathlib import Path

from validate_operational_readiness import validate


class ReadinessTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.template = json.loads((Path(__file__).parent / "operational-readiness.example.json").read_text(encoding="utf-8"))

    def valid(self) -> dict:
        value = copy.deepcopy(self.template)
        value["service"] = {"name": "orders", "owner": "team:commerce", "tier": 1}
        value["secrets"]["provider"] = "production-vault"
        value["secrets"]["references"] = ["vault://production/orders/database"]
        value["secrets"]["break_glass_runbook"] = "ops/runbooks/break-glass.md"
        value["telemetry"]["backend"] = "regional-otel-gateway"
        value["slos"][0]["owner"] = "team:commerce"
        value["supply_chain"]["lockfiles"] = ["go.sum"]
        value["release"]["artifact_digest"] = "sha256:" + "a" * 64
        value["release"]["rollback_command"] = ["deployer", "rollback", "--digest", "sha256:" + "b" * 64]
        return value

    def test_valid_contract_passes(self) -> None:
        self.assertEqual([], validate(self.valid()))

    def test_inline_secret_key_fails(self) -> None:
        value = self.valid()
        value["database_password"] = "do-not-store-this"
        self.assertTrue(any("secret-shaped" in error for error in validate(value)))

    def test_missing_signal_fails(self) -> None:
        value = self.valid()
        value["telemetry"]["signals"] = ["logs", "metrics"]
        self.assertTrue(any("telemetry.signals" in error for error in validate(value)))

    def test_placeholder_digest_fails(self) -> None:
        value = self.valid()
        value["release"]["artifact_digest"] = "sha256:" + "0" * 64
        self.assertTrue(any("artifact_digest" in error for error in validate(value)))

    def test_slo_without_runbook_fails(self) -> None:
        value = self.valid()
        value["slos"][0]["runbook"] = "replace-me"
        self.assertTrue(any("runbook" in error for error in validate(value)))

    def test_missing_restore_evidence_fails_closed(self) -> None:
        value = self.valid()
        value["recovery"]["restore_evidence"] = ""
        self.assertTrue(any("restore_evidence" in error for error in validate(value)))

    def write_evidence(self, root: Path, value: dict) -> None:
        # These files test presence only; they are not operational proofs.
        refs = [value["telemetry"]["redaction_test"], value["supply_chain"]["sbom"]["artifact"],
                value["supply_chain"]["provenance"]["artifact"], value["supply_chain"]["secret_scan_evidence"],
                value["release"]["rollback_test"], value["release"]["approval_record"], value["recovery"]["restore_evidence"]]
        for ref in refs:
            path = root / ref
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_text("{}\n", encoding="utf-8")

    def test_complete_evidence_file_set_passes_presence_gate(self) -> None:
        with tempfile.TemporaryDirectory() as folder:
            value = self.valid(); root = Path(folder)
            self.write_evidence(root, value)
            self.assertEqual([], validate(value, True, root))

    def test_artifact_references_cannot_disable_evidence_check(self) -> None:
        missing = object()
        for artifact in ("sbom", "provenance"):
            for mode in (False, True):
                for bad in (missing, None, "", "replace-me", 42, [], {}):
                    with self.subTest(artifact=artifact, require_files=mode, bad_type=type(bad).__name__):
                        with tempfile.TemporaryDirectory() as folder:
                            root = Path(folder); value = self.valid()
                            self.write_evidence(root, value)
                            if bad is missing:
                                del value["supply_chain"][artifact]["artifact"]
                            else:
                                value["supply_chain"][artifact]["artifact"] = bad
                            errors = validate(value, mode, root)
                            self.assertTrue(any(f"supply_chain.{artifact}.artifact" in e for e in errors), errors)

    def test_missing_or_directory_evidence_is_rejected(self) -> None:
        with tempfile.TemporaryDirectory() as folder:
            root = Path(folder); value = self.valid(); self.write_evidence(root, value)
            for artifact in ("sbom", "provenance"):
                path = root / value["supply_chain"][artifact]["artifact"]
                contents = path.read_bytes(); path.unlink()
                self.assertTrue(validate(value, True, root))
                path.mkdir()
                self.assertTrue(validate(value, True, root))
                path.rmdir(); path.write_bytes(contents)

    def test_evidence_escape_is_rejected_even_when_file_exists(self) -> None:
        with tempfile.TemporaryDirectory() as folder:
            root = Path(folder) / "project"; root.mkdir()
            value = self.valid(); self.write_evidence(root, value)
            outside = root.parent / "outside.json"; outside.write_text("{}", encoding="utf-8")
            value["supply_chain"]["sbom"]["artifact"] = "../outside.json"
            self.assertTrue(any("escapes root" in e for e in validate(value, True, root)))

    def test_invalid_filesystem_evidence_is_reported(self) -> None:
        with tempfile.TemporaryDirectory() as folder:
            root = Path(folder); value = self.valid(); self.write_evidence(root, value)
            value["supply_chain"]["sbom"]["artifact"] = "evidence/invalid\x00.json"
            errors = validate(value, True, root)
            self.assertTrue(errors)

    def test_unreadable_evidence_returns_validation_failure(self) -> None:
        with tempfile.TemporaryDirectory() as folder:
            root = Path(folder); value = self.valid(); self.write_evidence(root, value)
            with patch.object(Path, "is_file", side_effect=PermissionError("synthetic access failure")):
                errors = validate(value, True, root)
            self.assertTrue(any("cannot inspect file" in e for e in errors), errors)


if __name__ == "__main__":
    unittest.main()
