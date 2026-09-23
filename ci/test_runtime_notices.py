"""Local packaging regressions; all fixture data is synthetic."""
from pathlib import Path
import tempfile,json,hashlib,unittest
import collect_runtime_notices as n

class NoticeBoundaryTests(unittest.TestCase):
    def setUp(self):
        self.temp=tempfile.TemporaryDirectory();self.addCleanup(self.temp.cleanup);self.root=Path(self.temp.name)
        self.source=self.root/'source';self.runtime=self.root/'runtime';self.out=self.root/'out'
        for p in [self.source,self.runtime,self.out]:p.mkdir()
        self.rel='node_modules/.pnpm/example@1.0.0/node_modules/example/package.json'
        data=b'{"name":"example","version":"1.0.0","license":"MIT"}\n'
        for root in [self.source,self.runtime]:
            p=root/self.rel;p.parent.mkdir(parents=True);p.write_bytes(data)
        self.notice='licenses/example.txt';p=self.source/self.notice;p.parent.mkdir();p.write_text('Synthetic license fixture only. '*8)
        self.lock={'schema':'elite-reference-runtime-notices/v1','source_scope':'SYNTHETIC_FIXTURE',
          'manifests':{self.rel:{'sha256':hashlib.sha256(data).hexdigest(),'name':'example','version':'1.0.0','license':'MIT','notice_refs':[self.notice],'source_basis':'synthetic'}},
          'texts':{self.notice:n.sha(p)},'build_only_exclusions':[]}
        self.path=self.source/'lock.json'

    def run_collect(self):
        self.path.write_text(json.dumps(self.lock));return n.collect(self.source,self.runtime,self.out,self.path)

    def rejected(self):
        with self.assertRaises(ValueError):self.run_collect()
        self.assertFalse((self.out/'runtime-dependency-notices.json').exists())

    def test_unknown_runtime_manifest_blocks_release(self):
        p=self.runtime/'node_modules/injected/package.json';p.parent.mkdir();p.write_text('{}');self.rejected()

    def test_missing_runtime_dependency_blocks_release(self):
        (self.runtime/self.rel).unlink();self.rejected()

    def test_modified_runtime_manifest_blocks_release(self):
        (self.runtime/self.rel).write_text('{"name":"different"}');self.rejected()

    def test_changed_installed_source_blocks_release(self):
        (self.source/self.rel).write_text('{}');self.rejected()

    def test_changed_or_missing_notice_blocks_release(self):
        (self.source/self.notice).write_text('changed');self.rejected()
        (self.source/self.notice).unlink();self.rejected()

    def test_empty_notice_mapping_blocks_release(self):
        self.lock['manifests'][self.rel]['notice_refs']=[];self.rejected()

    def test_empty_notice_file_is_not_legal_evidence(self):
        p=self.source/self.notice;p.write_bytes(b'');self.lock['texts'][self.notice]=n.sha(p);self.rejected()

    def test_path_escape_blocks_before_publication(self):
        self.lock['texts']['../outside.txt']='0'*64;self.rejected()

    def test_unadmitted_native_payload_blocks_release(self):
        p=self.runtime/'node_modules/payload.node';p.write_bytes(b'synthetic-not-a-binary');self.rejected()

    def test_excluded_build_tool_cannot_be_shipped(self):
        self.lock['build_only_exclusions']=['example@1.0.0'];self.rejected()

    def test_existing_catalogue_is_never_overwritten(self):
        self.run_collect();before=(self.out/'runtime-dependency-notices.json').read_bytes()
        with self.assertRaises(ValueError):self.run_collect()
        self.assertEqual(before,(self.out/'runtime-dependency-notices.json').read_bytes())

    def test_duplicate_json_fields_rejected(self):
        self.path.write_text('{"schema":1,"schema":2}')
        with self.assertRaisesRegex(ValueError,'duplicate'):n.collect(self.source,self.runtime,self.out,self.path)

if __name__=='__main__':unittest.main()
