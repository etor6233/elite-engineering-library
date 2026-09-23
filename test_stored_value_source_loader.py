# AUTHORED source-closure regressions. Extra files contain only test exceptions;
# no provider calls, credentials, user data or harmful effects are involved.
import hashlib
import json
from pathlib import Path
import py_compile
import shutil
import subprocess
import sys
import tempfile
import unittest
from decimal import Decimal
from importlib.util import cache_from_source
from test_odoo_loyalty_derived import request

class SourceLoaderTests(unittest.TestCase):
    def run_copy(self, mutation=None):
        with tempfile.TemporaryDirectory(prefix='stored-value-source-fixture-') as temporary:
            root=Path(temporary)/'odoo_loyalty'
            shutil.copytree(Path(__file__).parent/'odoo_loyalty',root,ignore=shutil.ignore_patterns('__pycache__','*.pyc'))
            if mutation: mutation(root)
            lock=hashlib.sha256((root/'engine-lock.json').read_bytes()).hexdigest()
            return subprocess.run([sys.executable,'-I','-S','-B',str(root/'run.py'),lock],input=json.dumps(request()).encode(),capture_output=True,timeout=10)
    def test_ignores_unlocked_initializer(self):
        result=self.run_copy(lambda root:(root/'__init__.py').write_text("raise RuntimeError('UNLOCKED_INITIALIZER_FIXTURE')\n",encoding='utf-8'))
        self.assertEqual(result.returncode,0,result.stderr.decode())
        self.assertEqual([Decimal(p) for p in json.loads(result.stdout)['result']['points']],[Decimal(50),Decimal(50)])
    def test_ignores_unlocked_bytecode(self):
        def add(root):
            source=root/'fixture_cache_source.py'
            source.write_text("raise RuntimeError('UNLOCKED_BYTECODE_FIXTURE')\n",encoding='utf-8')
            cache=Path(cache_from_source(str(root/'protocol.py')));cache.parent.mkdir(exist_ok=True)
            py_compile.compile(str(source),cfile=str(cache),doraise=True,invalidation_mode=py_compile.PycInvalidationMode.UNCHECKED_HASH)
        result=self.run_copy(add)
        self.assertEqual(result.returncode,0,result.stderr.decode())
        self.assertEqual([Decimal(p) for p in json.loads(result.stdout)['result']['points']],[Decimal(50),Decimal(50)])
    def test_rejects_changed_locked_source(self):
        def add(root):
            p=root/'engine.py';p.write_bytes(p.read_bytes()+b'\n# changed fixture\n')
        result=self.run_copy(add)
        self.assertEqual(result.returncode,2)
        self.assertEqual(result.stdout,b'')


class SourceReplacementTests(unittest.TestCase):
    def test_exact_and_modified_source_rebuilds_preserve_licenses(self):
        from tools.rebuild_stored_value_source import rebuild
        original=Path(__file__).parent/'odoo_loyalty'
        baseline=hashlib.sha256((original/'engine-lock.json').read_bytes()).hexdigest()
        with tempfile.TemporaryDirectory(prefix='stored-value-replacement-fixture-') as temporary:
            root=Path(temporary)
            exact=root/'exact';receipt=rebuild(original,baseline,exact)
            self.assertEqual(receipt['classification'],'EXACT_REBUILD')
            self.assertEqual(receipt['replacement_lock_sha256'],baseline)
            (exact/'engine.py').write_bytes((exact/'engine.py').read_bytes()+b'\n# User modification fixture, no original attribution claimed.\n')
            modified=root/'modified';receipt=rebuild(exact,baseline,modified)
            self.assertEqual(receipt['classification'],'USER_MODIFIED_SOURCE_NOT_ADMITTED')
            self.assertEqual([r['path'] for r in receipt['changes']],['engine.py'])
            for name in ['LICENSE','COPYRIGHT']:
                self.assertEqual((original/name).read_bytes(),(modified/name).read_bytes())
            run=lambda digest:subprocess.run([sys.executable,'-I','-S','-B',str(modified/'run.py'),digest],input=json.dumps(request()).encode(),capture_output=True,timeout=10)
            self.assertEqual(run(baseline).returncode,2)
            self.assertEqual(run(receipt['replacement_lock_sha256']).returncode,0)
            with self.assertRaises(ValueError):rebuild(exact,baseline,modified)

if __name__=='__main__':unittest.main()
