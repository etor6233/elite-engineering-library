from __future__ import annotations

import json
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

import secure_local_file_gate as gate


class FakeRunner:
    def __init__(self) -> None:
        self.clam_exit = 0
        self.magika_label = "pdf"
        self.magika_mime = "application/pdf"
        self.magika_score = 0.99
        self.yara_matches: list[dict[str, str]] = []
        self.calls: list[list[str]] = []

    def __call__(self, args, timeout):
        call = [str(value) for value in args]
        self.calls.append(call)
        name = Path(call[0]).name
        if "--version" in call:
            if name == "clamscan.exe":
                return gate.CommandResult(0, "ClamAV 1.5.4/28104/Wed Aug 26 03:24:01 2026\n", "")
            if name == "magika.exe":
                return gate.CommandResult(0, gate.MAGIKA_VERSION + "\n", "")
            return gate.CommandResult(0, gate.YARA_VERSION + "\n", "")
        if name == "clamscan.exe":
            return gate.CommandResult(self.clam_exit, "detector output" if self.clam_exit else "", "")
        if name == "magika.exe":
            payload = [{"result": {"value": {"output": {
                "label": self.magika_label, "mime_type": self.magika_mime,
            }, "score": self.magika_score}}}]
            return gate.CommandResult(0, json.dumps(payload), "")
        return gate.CommandResult(0, json.dumps({"version": "1.20.0", "matches": self.yara_matches}), "")


class SecureLocalFileGateTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary = tempfile.TemporaryDirectory()
        self.root = Path(self.temporary.name)
        self.input = self.root / "invoice.pdf"
        self.input.write_bytes(b"%PDF-1.7\nfixture")
        self.policy = self.root / "policy.json"
        self.policy.write_text(json.dumps({
            "schema": gate.SCHEMA,
            "approval": {"status": "APPROVED", "approval_id": "SEC-1", "approved_by": "owner", "approved_at": "2026-08-26T00:00:00Z"},
            "allowed_content": [{"label": "pdf", "mime_type": "application/pdf", "min_score": 0.9}],
            "limits": {
                "max_input_bytes": 1024, "process_timeout_seconds": 10,
                "clamav_database_max_age_days": 2, "clamav_max_scan_bytes": 2048,
                "clamav_max_files": 10, "clamav_max_recursion": 2,
                "clamav_bytecode_timeout_ms": 1000, "yara_timeout_seconds": 5,
                "yara_max_matches_per_pattern": 10,
            },
            "yara": {"compiled_rules_sha256": "0" * 64, "reject_on_any_match": True},
            "archives": {"expansion_allowed": False}, "business_storage_authorized": False,
        }), encoding="utf-8")
        self.db = self.root / "db"
        self.db.mkdir()
        self.rules = self.root / "rules.yarc"
        self.rules.write_bytes(b"compiled")
        data = json.loads(self.policy.read_text(encoding="utf-8"))
        data["yara"]["compiled_rules_sha256"] = gate.sha256_file(self.rules)
        self.policy.write_text(json.dumps(data), encoding="utf-8")
        self.magika = self.root / "magika.exe"
        self.clam = self.root / "clamscan.exe"
        self.yara = self.root / "yr.exe"
        for path in (self.magika, self.clam, self.yara):
            path.write_bytes(path.name.encode("ascii"))
        self.hashes = {
            "magika": gate.sha256_file(self.magika),
            "clamscan": gate.sha256_file(self.clam),
            "yara_x": gate.sha256_file(self.yara),
        }

    def tearDown(self) -> None:
        self.temporary.cleanup()

    def execute(self, runner: FakeRunner, name: str = "evidence"):
        with patch.dict(gate.OFFICIAL_BINARY_SHA256, self.hashes, clear=True):
            return gate.run_gate(
                input_path=self.input, output=self.root / name, policy_path=self.policy,
                magika_exe=self.magika, clamscan_exe=self.clam, clamav_database=self.db,
                yara_exe=self.yara, compiled_rules=self.rules, runner=runner,
            )

    def test_all_selected_gates_admit_and_commit_hash_receipt(self):
        receipt = self.execute(FakeRunner())
        self.assertEqual(receipt["decision"], "ADMITTED")
        saved = json.loads((self.root / "evidence" / "security-receipt.json").read_text(encoding="utf-8"))
        self.assertEqual(saved["input"]["sha256"], gate.sha256_file(self.input))
        self.assertFalse(saved["business_storage_authorized"])

    def test_clamav_detection_rejects_before_other_content_engines(self):
        runner = FakeRunner()
        runner.clam_exit = 1
        receipt = self.execute(runner)
        self.assertEqual(receipt["reason"], "CLAMAV_DETECTED")
        self.assertFalse(any(Path(call[0]).name in {"magika.exe", "yr.exe"} for call in runner.calls))

    def test_stale_or_broken_clamav_database_rejects_closed(self):
        runner = FakeRunner()
        runner.clam_exit = 2
        receipt = self.execute(runner)
        self.assertEqual(receipt["reason"], "CLAMAV_ERROR_OR_STALE_DATABASE")

    def test_unapproved_content_type_rejects_before_yara(self):
        runner = FakeRunner()
        runner.magika_label = "zip"
        runner.magika_mime = "application/zip"
        receipt = self.execute(runner)
        self.assertEqual(receipt["reason"], "CONTENT_TYPE_NOT_APPROVED")
        self.assertFalse(any(Path(call[0]).name == "yr.exe" for call in runner.calls))

    def test_any_yara_match_rejects_and_records_only_rule_identifier(self):
        runner = FakeRunner()
        runner.yara_matches = [{"rule": "blocked_pattern", "file": "private-path"}]
        receipt = self.execute(runner)
        self.assertEqual(receipt["reason"], "YARA_RULE_MATCH")
        self.assertEqual(receipt["yara_x"]["matching_rules"], ["blocked_pattern"])
        self.assertNotIn("private-path", json.dumps(receipt))

    def test_awaiting_user_policy_blocks_before_any_tool(self):
        data = json.loads(self.policy.read_text(encoding="utf-8"))
        data["approval"]["status"] = "AWAITING_USER"
        self.policy.write_text(json.dumps(data), encoding="utf-8")
        runner = FakeRunner()
        with self.assertRaisesRegex(gate.GateConfigurationError, "not APPROVED"):
            self.execute(runner)
        self.assertEqual(runner.calls, [])

    def test_binary_hash_mismatch_blocks_before_execution(self):
        runner = FakeRunner()
        self.magika.write_bytes(b"tampered")
        with self.assertRaisesRegex(gate.GateConfigurationError, "SHA-256 mismatch"):
            self.execute(runner)
        self.assertEqual(runner.calls, [])


if __name__ == "__main__":
    unittest.main()
