from __future__ import annotations

import copy
from datetime import datetime, timedelta, timezone
import hashlib
import tempfile
from pathlib import Path
import unittest

from validate_production_admission import CONTROLS, REQUIRED_ASSERTIONS, validate


class ProductionAdmissionTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temp = tempfile.TemporaryDirectory()
        self.root = Path(self.temp.name).resolve()
        proof = self.root / "evidence.json"
        proof.write_bytes(b'{"pass":true}\n')
        now = datetime(2026, 8, 29, 18, tzinfo=timezone.utc)
        self.now = now
        stamp = lambda value: value.isoformat().replace("+00:00", "Z")
        evidence = {"path": "evidence.json", "bytes": proof.stat().st_size, "sha256": hashlib.sha256(proof.read_bytes()).hexdigest()}
        digest = "sha256:" + "a" * 64
        self.record = {
            "schema": "elite-production-admission/v1", "project_id": "revestex", "environment": "production",
            "release_digest": digest, "evaluated_at": stamp(now), "controls": [
                {"id": control, "project_id": "revestex", "environment": "production", "release_digest": digest,
                 "result": "PASS", "executed_at": stamp(now - timedelta(minutes=5)), "expires_at": stamp(now + timedelta(days=7)),
                 "executor_ref": "workload-identity://ci/production", "tool": {"name": "official-tool", "version": "1.0.0", "source": "https://example.invalid/official", "digest": "sha256:" + "b" * 64},
                 "target": {"id": f"production:{control}"}, "assertions": {key: True for key in sorted(REQUIRED_ASSERTIONS[control])}, "evidence": [evidence]}
                for control in sorted(CONTROLS)
            ]}

    def tearDown(self) -> None:
        self.temp.cleanup()

    def test_complete_fresh_hash_bound_record_passes(self) -> None:
        self.assertEqual([], validate(self.record, self.root, self.now))

    def test_missing_control_fails(self) -> None:
        value = copy.deepcopy(self.record); value["controls"].pop()
        self.assertTrue(any("eight unique" in error for error in validate(value, self.root, self.now)))

    def test_release_mismatch_fails(self) -> None:
        value = copy.deepcopy(self.record); value["controls"][0]["release_digest"] = "sha256:" + "c" * 64
        self.assertTrue(any("release_digest must match" in error for error in validate(value, self.root, self.now)))

    def test_stale_receipt_fails(self) -> None:
        value = copy.deepcopy(self.record); value["controls"][0]["expires_at"] = "2026-08-29T17:00:00Z"
        self.assertTrue(any("stale" in error for error in validate(value, self.root, self.now)))

    def test_tampered_evidence_fails(self) -> None:
        value = copy.deepcopy(self.record); (self.root / "evidence.json").write_text("tampered", encoding="utf-8")
        self.assertTrue(any("SHA-256 mismatch" in error for error in validate(value, self.root, self.now)))

    def test_false_assertion_fails(self) -> None:
        value = copy.deepcopy(self.record); key = next(iter(value["controls"][0]["assertions"])); value["controls"][0]["assertions"][key] = False
        self.assertTrue(any("semantic assertions" in error for error in validate(value, self.root, self.now)))

    def test_generic_exit_zero_assertion_cannot_replace_semantics(self) -> None:
        value = copy.deepcopy(self.record); value["controls"][0]["assertions"] = {"official_process_exit_zero": True}
        self.assertTrue(any("semantic assertions" in error for error in validate(value, self.root, self.now)))

    def test_duplicate_control_fails(self) -> None:
        value = copy.deepcopy(self.record); value["controls"][1]["id"] = value["controls"][0]["id"]
        self.assertTrue(any("eight unique" in error for error in validate(value, self.root, self.now)))

    def test_symlink_evidence_fails_when_supported(self) -> None:
        link = self.root / "link.json"
        try:
            link.symlink_to(self.root / "evidence.json")
        except OSError:
            self.skipTest("symlinks unavailable")
        value = copy.deepcopy(self.record); value["controls"][0]["evidence"][0]["path"] = "link.json"
        self.assertTrue(any("non-symlink" in error for error in validate(value, self.root, self.now)))


if __name__ == "__main__":
    unittest.main()
