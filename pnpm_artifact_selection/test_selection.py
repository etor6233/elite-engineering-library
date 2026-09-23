"""AUTHORED safety regressions; tiny synthetic inputs never represent pnpm admission."""
import concurrent.futures
import contextlib
import io
import json
import os
from pathlib import Path
import subprocess
import tarfile
import tempfile
import unittest
from unittest import mock
import select_artifact as m

def tar_bytes(items):
    out = io.BytesIO()
    with tarfile.open(fileobj=out, mode="w:gz") as tar:
        for name, data, kind in items:
            info = tarfile.TarInfo(name); info.type = kind; info.size = len(data)
            if kind != tarfile.REGTYPE:
                info.linkname = "../../outside"
            tar.addfile(info, io.BytesIO(data))
    return out.getvalue()

class SelectionTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory(prefix="elite-selection-test-")
        self.root = Path(self.tmp.name)
        self.items = [("package/bin/tool.js", b"fixture only", tarfile.REGTYPE),
                      ("package/LICENSE", b"synthetic notice", tarfile.REGTYPE),
                      ("package/native.exe", b"not executable", tarfile.REGTYPE)]
        self.payload = tar_bytes(self.items)
        self.entries = [{"path": n[8:], "bytes": len(data), "sha256": m.digest(data),
                         "selection": "excluded_native_family" if n.endswith(".exe") else "unchanged"}
                        for n, data, _ in self.items]
        self.policy = {"archive_sha256": m.digest(self.payload), "archive_bytes": len(self.payload),
                       "files": self.entries, "receipt_bindings": {"id": "synthetic", "executed": False}}
        self.acquisition = {"id": "synthetic", "executed": False,
                            "profile_sha256": "a" * 64, "approval_sha256": "b" * 64, "lock_sha256": "c" * 64}
        self.artifact = self.root / "fixture.tgz"; self.artifact.write_bytes(self.payload)
        self.receipt = self.root / "acquisition.json"; self.receipt.write_bytes(m.json_bytes(self.acquisition))
        self.target = self.root / "selected"

    def tearDown(self):
        self.tmp.cleanup()

    def build(self, target=None):
        return m.materialize(self.artifact, self.receipt, m.digest(self.receipt.read_bytes()), target or self.target, self.policy)

    def test_real_locked_policy(self):
        policy = m.load_policy()
        self.assertEqual(len(policy["files"]), 455)
        self.assertEqual(sum(x["selection"] == "unchanged" for x in policy["files"]), 442)

    def test_projection_deterministic_and_exclusion(self):
        self.build(); other = self.root / "other"; self.build(other)
        left = {p.relative_to(self.target).as_posix(): p.read_bytes() for p in self.target.rglob("*") if p.is_file()}
        right = {p.relative_to(other).as_posix(): p.read_bytes() for p in other.rglob("*") if p.is_file()}
        self.assertEqual(left, right)
        self.assertNotIn("payload/native.exe", left)
        self.assertFalse(m.verify(self.target, self.policy)["runtime_admitted"])

    def test_occupied_targets_untouched(self):
        for kind in ("file", "directory"):
            with self.subTest(kind=kind):
                target = self.root / kind
                if kind == "file": target.write_bytes(b"owner")
                else: target.mkdir()
                with self.assertRaises(m.SelectionError): self.build(target)
                self.assertTrue(target.exists())
                if kind == "file": self.assertEqual(target.read_bytes(), b"owner")

    def test_archive_and_receipt_tamper_before_output(self):
        for file in (self.artifact, self.receipt):
            before = file.read_bytes()
            file.write_bytes(before + b"changed")
            with self.assertRaises((m.SelectionError, ValueError)): self.build()
            self.assertFalse(self.target.exists())
            file.write_bytes(before)

    def test_receipt_hash_is_required(self):
        with self.assertRaises(m.SelectionError):
            m.materialize(self.artifact, self.receipt, "0" * 64, self.target, self.policy)
        self.assertFalse(self.target.exists())

    def test_acquisition_bindings_and_no_admission(self):
        for key, value in (("id", "wrong"), ("executed", True), ("executed", 0), ("approval_sha256", "")):
            receipt = dict(self.acquisition); receipt[key] = value
            self.receipt.write_bytes(m.json_bytes(receipt))
            with self.subTest(key=key, value=value), self.assertRaises(m.SelectionError): self.build()
            self.assertFalse(self.target.exists())

    def test_duplicate_json_rejected(self):
        with self.assertRaises(m.SelectionError): m.decode(b'{"id":1,"id":2}')

    def test_tar_duplicate_unexpected_missing(self):
        for index, items in enumerate((self.items + [self.items[0]], self.items + [("package/extra", b"x", tarfile.REGTYPE)], self.items[:-1])):
            with self.subTest(count=len(items)), self.assertRaises(m.SelectionError):
                m.project_payload(tar_bytes(items), self.entries, self.root / ("parser-" + str(index)))

    def test_tar_links_and_digest_mismatch(self):
        for kind in (tarfile.SYMTYPE, tarfile.LNKTYPE, tarfile.FIFOTYPE):
            items = [(self.items[0][0], self.items[0][1], kind), *self.items[1:]]
            with self.subTest(kind=kind), self.assertRaises(m.SelectionError):
                m.project_payload(tar_bytes(items), self.entries, self.root / "parser")
        entries = [dict(e) for e in self.entries]; entries[2]["sha256"] = "0" * 64
        with self.assertRaises(m.SelectionError): m.project_payload(self.payload, entries, self.root / "digest")

    def test_unsafe_names(self):
        for name in ("../outside", "/absolute", "x\\y", "a//b", "a/./b", "a/b.", "a/b ", "a:stream", "NUL.txt", "com1", "x\x00y"):
            with self.subTest(name=name), self.assertRaises(m.SelectionError): m.safe_name(name)

    def test_extraction_fault_cleans_only_staging(self):
        before = self.artifact.read_bytes()
        def fault(_payload, _entries, target):
            (target / "partial").write_bytes(b"partial")
            raise OSError("synthetic disk failure")
        with mock.patch.object(m, "project_payload", side_effect=fault), self.assertRaises(OSError): self.build()
        self.assertFalse(self.target.exists())
        self.assertEqual(list(self.root.glob(".pnpm-selection-*")), [])
        self.assertEqual(self.artifact.read_bytes(), before)

    def test_late_competitor_survives(self):
        rename = m.os.rename
        def compete(source, target):
            target.mkdir(); (target / "owner").write_bytes(b"other publisher")
            return rename(source, target)
        with mock.patch.object(m.os, "rename", side_effect=compete), self.assertRaises(OSError): self.build()
        self.assertEqual((self.target / "owner").read_bytes(), b"other publisher")
        self.assertEqual(list(self.root.glob(".pnpm-selection-*")), [])

    def test_concurrent_publication_one_winner(self):
        def attempt():
            try: self.build(); return "ok"
            except (OSError, m.SelectionError): return "blocked"
        with concurrent.futures.ThreadPoolExecutor(max_workers=2) as pool:
            results = list(pool.map(lambda _: attempt(), range(2)))
        self.assertEqual(sorted(results), ["blocked", "ok"])
        m.verify(self.target, self.policy)

    def test_verify_detects_missing_extra_and_changed(self):
        self.build()
        file = self.target / "payload/LICENSE"; original = file.read_bytes()
        for mutate in (lambda: file.write_bytes(b"changed"), lambda: file.unlink()):
            mutate()
            with self.assertRaises(m.SelectionError): m.verify(self.target, self.policy)
            file.write_bytes(original)
        (self.target / "extra").write_bytes(b"unexpected")
        with self.assertRaises(m.SelectionError): m.verify(self.target, self.policy)

    def test_false_admission_receipt_rejected(self):
        self.build(); file = self.target / m.RECEIPT; data = m.decode(file.read_bytes()); data["runtime_admitted"] = True
        file.write_bytes(m.json_bytes(data))
        with self.assertRaises(m.SelectionError): m.verify(self.target, self.policy)

    def test_junction_input_and_target_parent_rejected(self):
        actual = self.root / "real"; actual.mkdir(); link = self.root / "junction"
        result = subprocess.run(["cmd.exe", "/d", "/c", "mklink", "/J", str(link), str(actual)], capture_output=True, creationflags=0x08000000)
        self.assertEqual(result.returncode, 0, result.stderr.decode(errors="replace"))
        try:
            with self.assertRaises(m.SelectionError): m.no_reparse(link / "missing")
            with self.assertRaises(m.SelectionError): self.build(link / "target")
            self.assertEqual(list(actual.iterdir()), [])
        finally:
            link.rmdir()

    def test_cli_policy_tamper_and_missing_paths_no_traceback(self):
        with mock.patch.object(m, "POLICY_SHA256", "0" * 64), contextlib.redirect_stderr(io.StringIO()) as err:
            self.assertEqual(m.main(["verify", "--target", str(self.target)]), 2)
        self.assertNotIn("Traceback", err.getvalue())
        with contextlib.redirect_stderr(io.StringIO()) as err:
            self.assertEqual(m.main(["verify", "--target", str(self.target)]), 2)
        self.assertNotIn("Traceback", err.getvalue())

if __name__ == "__main__":
    unittest.main()
