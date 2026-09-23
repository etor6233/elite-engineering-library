from __future__ import annotations

import copy
import hashlib
import json
import os
from pathlib import Path
import sys
import tempfile
import unittest
import time
import traceback
import threading
from unittest.mock import patch
import run_official_tool as runner

from run_official_tool import run


class OfficialToolRunnerTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temp = tempfile.TemporaryDirectory(); self.root = Path(self.temp.name).resolve(); (self.root / "work").mkdir()
        binary = Path(sys.executable).resolve()
        self.profile = {"schema":"elite-official-tool-run/v1", "project_id":"revestex", "environment":"production", "release_digest":"sha256:" + "a" * 64, "control_id":"LOAD_RESILIENCE", "tool":{"path":str(binary), "sha256":hashlib.sha256(binary.read_bytes()).hexdigest(), "name":"CPython fixture", "version":sys.version.split()[0], "source":"https://python.org"}, "arguments":["-I", "-c", "print('official-output')"], "working_directory":"work", "environment_variable_names":[], "timeout_seconds":10, "target":{"id":"production:test"}}

    def tearDown(self) -> None: self.temp.cleanup()

    def test_exact_binary_runs_without_shell_and_hashes_outputs(self) -> None:
        receipt = run(self.profile, self.root, self.root / "receipt.json", True)
        self.assertEqual(0, receipt["exit_code"]); self.assertEqual(b"official-output\r\n" if os.name == "nt" else b"official-output\n", (self.root / "receipt.stdout.bin").read_bytes())
        self.assertEqual(hashlib.sha256((self.root / "receipt.stdout.bin").read_bytes()).hexdigest(), receipt["stdout"]["sha256"])
        expected = hashlib.sha256(json.dumps(self.profile["arguments"], ensure_ascii=False, separators=(",", ":")).encode()).hexdigest()
        self.assertEqual(expected, receipt["arguments_sha256"]); self.assertEqual("work", receipt["working_directory"])

    def test_explicit_authorization_is_required(self) -> None:
        with self.assertRaisesRegex(ValueError, "authorize"): run(self.profile, self.root, self.root / "receipt.json", False)

    def test_binary_hash_mismatch_fails_before_execution(self) -> None:
        value = copy.deepcopy(self.profile); value["tool"]["sha256"] = "0" * 64
        with self.assertRaisesRegex(ValueError, "SHA-256"): run(value, self.root, self.root / "receipt.json", True)

    def test_missing_environment_reference_fails_without_value_leak(self) -> None:
        value = copy.deepcopy(self.profile); value["environment_variable_names"] = ["ELITE_MISSING_TEST_SECRET"]
        with self.assertRaisesRegex(ValueError, "ELITE_MISSING_TEST_SECRET"): run(value, self.root, self.root / "receipt.json", True)

    def test_nonzero_exit_is_preserved_not_promoted(self) -> None:
        value = copy.deepcopy(self.profile); value["arguments"] = ["-I", "-c", "raise SystemExit(7)"]
        receipt = run(value, self.root, self.root / "receipt.json", True); self.assertEqual(7, receipt["exit_code"])


    def assert_no_receipt(self):
        self.assertEqual([], list(self.root.glob("receipt*")))

    def test_excess_is_stopped_before_later_effect(self):
        self.profile["arguments"] = ["-I", "-c", "import os,time,pathlib;os.write(1,b'x'*65536);time.sleep(.3);pathlib.Path('continued').write_text('bad')"]
        with patch.object(runner, "MAX_OUTPUT", 1024), self.assertRaisesRegex(ValueError, "output exceeded"):
            run(self.profile, self.root, self.root / "receipt.json", True)
        self.assertFalse((self.root / "work/continued").exists())
        self.assert_no_receipt()

    def test_timeout_has_no_private_argument_and_child_is_reaped(self):
        self.profile["timeout_seconds"] = 1
        self.profile["arguments"] = ["-I", "-c", "import time;time.sleep(4)", "PRIVATE_MARKER_345"]
        children = []
        popen = runner.subprocess.Popen
        def launch(*args, **kwargs):
            child = popen(*args, **kwargs); children.append(child); return child
        before = {t.ident for t in threading.enumerate()}
        start = time.monotonic()
        with patch.object(runner.subprocess, "Popen", side_effect=launch):
            try: run(self.profile, self.root, self.root / "receipt.json", True)
            except ValueError as error:
                self.assertEqual("official tool timed out", str(error))
                self.assertNotIn("PRIVATE_MARKER_345", traceback.format_exc())
            else: self.fail("timeout accepted")
        self.assertLess(time.monotonic() - start, 3)
        self.assertEqual(1, len(children))
        self.assertIsNotNone(children[0].poll())
        self.assertTrue(children[0].stdout.closed and children[0].stderr.closed)
        self.assertEqual(before, {t.ident for t in threading.enumerate()})
        self.assert_no_receipt()

    def test_both_streams_at_exact_limit_preserve_binary_content(self):
        limit = 1024 * 1024
        self.profile["arguments"] = ["-I", "-c", "import os;[(os.write(1,bytes(range(256))*16),os.write(2,bytes(reversed(range(256)))*16)) for _ in range(256)]"]
        with patch.object(runner, "MAX_OUTPUT", limit):
            receipt = run(self.profile, self.root, self.root / "receipt.json", True)
        self.assertEqual(0, receipt["exit_code"])
        for kind, data in [("stdout", bytes(range(256))*4096), ("stderr", bytes(reversed(range(256)))*4096)]:
            self.assertEqual(data, (self.root / receipt[kind]["path"]).read_bytes())
            self.assertEqual(limit, receipt[kind]["bytes"])
            self.assertEqual(hashlib.sha256(data).hexdigest(), receipt[kind]["sha256"])

    def test_each_stream_one_byte_over_limit_rejected(self):
        for fd in (1, 2):
            with self.subTest(fd=fd):
                self.profile["arguments"] = ["-I", "-c", f"import os;os.write({fd},b'x'*1025)"]
                with patch.object(runner, "MAX_OUTPUT", 1024), self.assertRaisesRegex(ValueError, "output exceeded"):
                    run(self.profile, self.root, self.root / "receipt.json", True)
                self.assert_no_receipt()

    def test_closed_streams_do_not_make_live_child_successful(self):
        self.profile["timeout_seconds"] = 1
        self.profile["arguments"] = ["-I", "-c", "import os,time;os.close(1);os.close(2);time.sleep(4)"]
        with self.assertRaisesRegex(ValueError, "timed out"):
            run(self.profile, self.root, self.root / "receipt.json", True)
        self.assert_no_receipt()

    def test_inherited_pipe_has_deadline_without_waiting_for_descendant(self):
        # The descendant is finite and exits itself. This explicitly does NOT
        # claim tree termination; a service manager must supply that boundary.
        self.profile["timeout_seconds"] = 1
        child_code = "import time,pathlib;time.sleep(2.5);pathlib.Path('descendant-finished').write_text('done')"
        self.profile["arguments"] = ["-I", "-c", f"import subprocess,sys;subprocess.Popen([sys.executable,'-I','-c',{child_code!r}])"]
        start = time.monotonic()
        try:
            with self.assertRaisesRegex(ValueError, "timed out"):
                run(self.profile, self.root, self.root / "receipt.json", True)
            self.assertLess(time.monotonic() - start, 2)
            self.assert_no_receipt()
        finally:
            until = time.monotonic() + 4
            while not (self.root / "work/descendant-finished").exists() and time.monotonic() < until:
                time.sleep(.02)
            self.assertTrue((self.root / "work/descendant-finished").is_file())

    def test_stdin_is_closed_and_unselected_environment_is_absent(self):
        self.profile["arguments"] = ["-I", "-c", "import sys,os;assert sys.stdin.buffer.read()==b'';assert 'ELITE_UNSELECTED_SECRET' not in os.environ;print('closed')"]
        with patch.dict(os.environ, {"ELITE_UNSELECTED_SECRET": "PRIVATE_MARKER_345"}):
            receipt = run(self.profile, self.root, self.root / "receipt.json", True)
        self.assertEqual(0, receipt["exit_code"])

    def test_boolean_timeout_rejected_before_process_creation(self):
        self.profile["timeout_seconds"] = True
        with patch.object(runner.subprocess, "Popen") as launch, self.assertRaisesRegex(ValueError, "timeout_seconds"):
            run(self.profile, self.root, self.root / "receipt.json", True)
        launch.assert_not_called()

    def test_read_error_is_static_and_cleans_up_direct_child(self):
        self.profile["arguments"] = ["-I", "-c", "import time;time.sleep(4)"]
        children = []; popen = runner.subprocess.Popen; read = runner.os.read
        def launch(*args, **kwargs):
            child = popen(*args, **kwargs); children.append(child); return child
        def fail_read(fd, count):
            if children and fd in (children[0].stdout.fileno(), children[0].stderr.fileno()):
                raise OSError("PRIVATE_MARKER_345")
            return read(fd, count)
        with patch.object(runner.subprocess, "Popen", side_effect=launch), patch.object(runner.os, "read", side_effect=fail_read):
            try: run(self.profile, self.root, self.root / "receipt.json", True)
            except ValueError as error:
                self.assertEqual("official tool execution unavailable", str(error))
                self.assertNotIn("PRIVATE_MARKER_345", traceback.format_exc())
            else: self.fail("read failure accepted")
        self.assertIsNotNone(children[0].poll())
        self.assert_no_receipt()



if __name__ == "__main__": unittest.main()
