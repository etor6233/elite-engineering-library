"""Local router failure and preservation tests; real admitted installs have separate receipts."""
import base64
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest
from unittest.mock import patch
import local_runtime as runtime


class PayloadTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(); self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name); (self.root / 'file.txt').write_bytes(b'locked')
        self.files = {'file.txt': {'bytes': 6, 'sha256': hashlib.sha256(b'locked').hexdigest()}}

    def test_exact(self):
        self.assertEqual(runtime.verify_payload(self.root, self.files), self.root)

    def test_missing(self):
        (self.root / 'file.txt').unlink()
        with self.assertRaises(runtime.SelectionError): runtime.verify_payload(self.root, self.files)

    def test_changed(self):
        (self.root / 'file.txt').write_bytes(b'drift!')
        with self.assertRaises(runtime.SelectionError): runtime.verify_payload(self.root, self.files)

    def test_unexpected_file(self):
        (self.root / 'addon.node').write_bytes(b'never executed')
        with self.assertRaises(runtime.SelectionError): runtime.verify_payload(self.root, self.files)

    def test_unexpected_empty_directory(self):
        (self.root / 'extra').mkdir()
        with self.assertRaises(runtime.SelectionError): runtime.verify_payload(self.root, self.files)


class NoticeTests(unittest.TestCase):
    def test_exact_supplied_bytes_without_final_newline(self):
        raw = b'original fixture\r\nwithout final newline'; h = hashlib.sha256(raw).hexdigest()
        value = {'entries': [], 'texts': {h: {'bytes': len(raw), 'content_base64': base64.b64encode(raw).decode()}}}
        self.assertIn(raw + b'\nEND ORIGINAL TEXT', runtime.notice_bytes(runtime.json_bytes(value)))
        value['texts'][h]['content_base64'] = base64.b64encode(raw + b'changed').decode()
        with self.assertRaises(runtime.SelectionError): runtime.notice_bytes(runtime.json_bytes(value))


@unittest.skipUnless(os.name == 'nt', 'Windows runtime orchestration')
class RecoveryTests(unittest.TestCase):
    def attempt(self, code, timeout):
        temp = tempfile.TemporaryDirectory(); self.addCleanup(temp.cleanup); root = Path(temp.name)
        # Only admission is replaced for this subprocess recovery fixture. No pnpm/runtime admission is claimed.
        value = {'scope': 'RECOVERY_FIXTURE_ONLY', 'argv': [sys.executable, '-I', '-B', '-c', code],
                 'cwd': str(root), 'environment': {'SystemRoot': os.environ['SystemRoot']}}
        with patch.object(runtime, 'verify_plan', return_value=value), patch.object(runtime, 'MAX_SECONDS', timeout):
            with self.assertRaises((runtime.SelectionError, subprocess.TimeoutExpired)): runtime.execute(root, 'fixture-only')
        self.assertTrue((root / runtime.STARTED).is_file()); self.assertTrue((root / runtime.LOG).is_file())
        result = json.loads((root / runtime.RESULT).read_text()); self.assertEqual(result['state'], 'FAIL')
        self.assertTrue(result['executed']); self.assertFalse(result['redistribution_admitted'])
        return result

    def test_nonzero_exit_preserves_result(self):
        result = self.attempt("raise SystemExit(23)", 10)
        self.assertEqual(result['exit_code'], 23)

    def test_timeout_preserves_result(self):
        result = self.attempt("import time; time.sleep(10)", 0.1)
        self.assertIn('timed out', result['error'])


if __name__ == '__main__':
    unittest.main()
