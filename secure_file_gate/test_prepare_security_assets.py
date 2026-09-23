from __future__ import annotations

import json
import subprocess
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

import prepare_security_assets as assets


class PrepareSecurityAssetsTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary = tempfile.TemporaryDirectory()
        self.root = Path(self.temporary.name)
        self.fresh = self.root / "freshclam.exe"
        self.clam = self.root / "clamscan.exe"
        self.yara = self.root / "yr.exe"
        for path in (self.fresh, self.clam, self.yara):
            path.write_bytes(path.name.encode("ascii"))

    def tearDown(self) -> None:
        self.temporary.cleanup()

    def completed(self, code=0, stdout="", stderr=""):
        return subprocess.CompletedProcess([], code, stdout, stderr)

    def test_update_does_not_trust_freshclam_exit_without_scan_proof(self):
        db = self.root / "db"
        hashes = {
            "FRESHCLAM_SHA256": assets.sha256_file(self.fresh),
            "CLAMSCAN_SHA256": assets.sha256_file(self.clam),
        }
        def fake_run(args, timeout):
            if Path(args[0]).name == "freshclam.exe":
                db.mkdir(exist_ok=True)
                for name in ("main.cvd", "daily.cvd", "bytecode.cvd"):
                    (db / name).write_bytes(name.encode("ascii"))
                return self.completed(0)
            if "--version" in args:
                return self.completed(0, "ClamAV 1.5.4/28104/Wed Aug 26 03:24:01 2026\n")
            return self.completed(2, "", "stale")
        with patch.multiple(assets, **hashes), patch.object(assets, "run", side_effect=fake_run):
            with self.assertRaisesRegex(RuntimeError, "not proven"):
                assets.update_clamav(self.fresh, self.clam, db, 2, 30)
        self.assertFalse((db / "clamav-database-receipt.json").exists())

    def test_update_records_inventory_only_after_fresh_scan(self):
        db = self.root / "db"
        hashes = {
            "FRESHCLAM_SHA256": assets.sha256_file(self.fresh),
            "CLAMSCAN_SHA256": assets.sha256_file(self.clam),
        }
        def fake_run(args, timeout):
            if Path(args[0]).name == "freshclam.exe":
                db.mkdir(exist_ok=True)
                for name in ("main.cvd", "daily.cvd", "bytecode.cvd"):
                    (db / name).write_bytes(name.encode("ascii"))
                return self.completed(0)
            if "--version" in args:
                return self.completed(0, "ClamAV 1.5.4/28104/Wed Aug 26 03:24:01 2026\n")
            return self.completed(0)
        with patch.multiple(assets, **hashes), patch.object(assets, "run", side_effect=fake_run):
            assets.update_clamav(self.fresh, self.clam, db, 2, 30)
        receipt = json.loads((db / "clamav-database-receipt.json").read_text(encoding="utf-8"))
        self.assertEqual(len(receipt["files"]), 3)
        self.assertIn("official-db-only", receipt["validation"])

    def test_yara_compile_is_atomic_and_hash_receipted(self):
        source = self.root / "approved.yar"
        source.write_text("rule approved { condition: false }\n", encoding="utf-8")
        output = self.root / "approved.yarc"
        def fake_run(args, timeout):
            if "compile" in args:
                Path(args[args.index("--output") + 1]).write_bytes(b"compiled")
                return self.completed(0)
            return self.completed(0, '{"version":"1.20.0","matches":[]}')
        with patch.object(assets, "YARA_SHA256", assets.sha256_file(self.yara)), patch.object(assets, "run", side_effect=fake_run):
            assets.compile_yara(self.yara, source, output, 30)
        receipt = json.loads(output.with_suffix(".yarc.receipt.json").read_text(encoding="utf-8"))
        self.assertEqual(receipt["compiled_rules_sha256"], assets.sha256_file(output))

    def test_existing_yara_output_fails_without_overwrite(self):
        source = self.root / "approved.yar"
        source.write_text("rule approved { condition: false }\n", encoding="utf-8")
        output = self.root / "approved.yarc"
        output.write_bytes(b"existing")
        with patch.object(assets, "YARA_SHA256", assets.sha256_file(self.yara)):
            with self.assertRaises(FileExistsError):
                assets.compile_yara(self.yara, source, output, 30)
        self.assertEqual(output.read_bytes(), b"existing")


if __name__ == "__main__":
    unittest.main()
