import concurrent.futures
import base64
import hashlib
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest
from unittest.mock import patch
import plan_install as gate

@unittest.skipUnless(os.name == 'nt', 'Windows exclusive-publication contract')
class RoutingTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory(prefix='elite-route-test-')
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name)
        self.project, self.projection, self.store, self.cache = [self.root / x for x in ('project', 'projection', 'store', 'cache')]
        for p in (self.project, self.projection, self.store, self.cache): p.mkdir()
        self.node = self.root / 'node.exe'; self.node.write_bytes(b'fixture never executed')
        self.acq = self.root / 'acquisition.json'; self.acq.write_bytes(b'{}\n')
        self.acq_sha = gate.digest(self.acq.read_bytes())
        (self.projection / 'selection-receipt.json').write_bytes(b'fixture receipt')
        self.filehash = {}
        for name in ('package.json', 'pnpm-lock.yaml'):
            (self.project / name).write_bytes((name + '\n').encode())
            self.filehash[name] = gate.digest((self.project / name).read_bytes())
        self.output = self.root / 'plan'
        self.catalog = json.loads(Path(gate.__file__).with_name('retained-notice-catalog.json').read_bytes())
        for entry in self.catalog['entries']:
            if entry['origin'] == 'ORIGINAL_SELECTED_PAYLOAD':
                dest = self.projection / 'payload' / entry['payload_path']
                dest.parent.mkdir(parents=True, exist_ok=True)
                dest.write_bytes(base64.b64decode(entry['content_base64']))
        self.notice_bindings = {}
        for name, original in gate.BLUEOAK.items():
            binding = dict(original)
            raw = gate.json_bytes({'name': name, 'version': binding['version'], 'license': 'BlueOak-1.0.0'})
            dest = self.projection / 'payload' / binding['path']; dest.parent.mkdir(parents=True, exist_ok=True); dest.write_bytes(raw)
            binding['sha256'] = gate.digest(raw)
            self.notice_bindings[name] = binding
        notice_patch = patch.dict(gate.BLUEOAK, self.notice_bindings, clear=True)
        notice_patch.start(); self.addCleanup(notice_patch.stop)
        bundle = b''; modules = []
        for index, original in enumerate(gate.QRCODE_MODULES):
            raw = ('synthetic module ' + str(index) + '\n').encode()
            item = dict(original); item.update(start=len(bundle), end=len(bundle)+len(raw), region_sha256=gate.digest(raw))
            modules.append(item); bundle += raw
        self.qrcode_bundle = self.projection / 'payload' / gate.QRCODE_BUNDLE
        self.qrcode_bundle.parent.mkdir(parents=True, exist_ok=True); self.qrcode_bundle.write_bytes(bundle)
        for patched in [patch.object(gate, 'QRCODE_BUNDLE_SHA256', gate.digest(bundle)),
                        patch.object(gate, 'QRCODE_MODULES', modules),
                        patch.object(gate, 'NEXT_PATH_REGION', {'start': 0, 'end': len(bundle), 'sha256': gate.digest(bundle)}),
                        patch.object(gate, 'SEMVER_REGION', dict(gate.SEMVER_REGION, start=0, end=len(bundle), region_sha256=gate.digest(bundle)))]:
            patched.start(); self.addCleanup(patched.stop)
        for p in [patch.dict(gate.PROFILES, {'fixture': self.filehash}, clear=True),
                  patch.object(gate, 'NODE_SHA256', gate.digest(self.node.read_bytes())),
                  patch.object(gate, 'load_policy', return_value={}),
                  patch.object(gate, 'acquisition_binding', return_value=self.acq_sha),
                  patch.object(gate, 'verify', return_value={'acquisition_receipt_sha256': self.acq_sha})]:
            p.start(); self.addCleanup(p.stop)

    def create(self, **changes):
        args = dict(consumer='fixture', project=self.project, projection=self.projection, node=self.node,
                    store=self.store, cache=self.cache, target=self.output, acquisition_receipt=self.acq, acquisition_sha256=self.acq_sha)
        args.update(changes)
        return gate.create_plan(**args)

    def rejected(self, **changes):
        with self.assertRaises((gate.SelectionError, OSError, ValueError)):
            self.create(**changes)
        self.assertFalse(self.output.exists())
        self.assertEqual([], list(self.root.glob('.pnpm-plan-*')))

    def test_recipe_is_bound_and_never_admitted(self):
        sha = self.create(); value = gate.verify_plan(self.output, sha)
        self.assertFalse(value['executed']); self.assertFalse(value['runtime_admitted'])
        self.assertFalse(value['redistribution_admitted']); self.assertFalse(value['shell'])
        for arg in ['--offline','--frozen-lockfile','--ignore-scripts','--ignore-pnpmfile','--ignore-workspace','--package-import-method=copy']:
            self.assertIn(arg, value['argv'])
        self.assertEqual(str(self.project), value['cwd'])
        self.assertTrue(any(x.startswith('--config.userconfig=') for x in value['argv']))
        self.assertTrue(any(x.startswith('--config.globalconfig=') for x in value['argv']))
        self.assertFalse(any(x.startswith(('--userconfig=', '--globalconfig=')) for x in value['argv']))

    def test_host_secrets_and_options_are_not_inherited(self):
        with patch.dict(os.environ, {'NODE_OPTIONS':'--require=evil.cjs','NPM_CONFIG_REGISTRY':'secret-value','PNPM_HOME':'secret-home','AWS_SECRET_ACCESS_KEY':'secret-key'}):
            sha=self.create(); value=gate.verify_plan(self.output,sha)
        self.assertEqual('REPLACE_NOT_MERGE',value['environment_mode'])
        self.assertNotIn('NODE_OPTIONS',value['environment']); self.assertNotIn('NPM_CONFIG_REGISTRY',value['environment'])
        self.assertNotIn('secret-',json.dumps(value))

    def test_unknown_consumer(self): self.rejected(consumer='arbitrary')
    def test_manifest_drift(self):
        (self.project/'package.json').write_text('changed'); self.rejected()
    def test_lock_drift(self):
        (self.project/'pnpm-lock.yaml').write_text('changed'); self.rejected()
    def test_local_npmrc(self):
        (self.project/'.npmrc').write_text('registry=changed'); self.rejected()
    def test_ancestor_hook(self):
        (self.root/'.pnpmfile.cjs').write_text('throw new Error()'); self.rejected()
    def test_workspace_not_in_profile(self):
        (self.project/'pnpm-workspace.yaml').write_text('packages: [evil]'); self.rejected()
    def test_wrong_node(self):
        self.node.write_bytes(b'changed'); self.rejected()
    def test_wrong_acquisition_sha(self): self.rejected(acquisition_sha256='0'*64)
    def test_projection_from_other_receipt(self):
        with patch.object(gate,'verify',return_value={'acquisition_receipt_sha256':'1'*64}): self.rejected()
    def test_existing_install_is_rejected(self):
        (self.project/'node_modules').mkdir(); self.rejected()
    def test_known_workspace_policy_is_preserved(self):
        p=self.project/'pnpm-workspace.yaml';p.write_bytes(b'known-workspace')
        self.filehash[p.name]=gate.digest(p.read_bytes())
        sha=self.create(); value=gate.verify_plan(self.output,sha)
        self.assertNotIn('--ignore-workspace',value['argv'])
        self.assertIn(p.name,value['consumer_files'])
    def test_missing_metadata_cache_is_rejected(self):
        self.rejected(cache=self.root/'absent-cache')
    def test_overlap(self): self.rejected(target=self.project/'plan')
    def test_existing_destination_is_preserved(self):
        self.output.mkdir(); (self.output/'marker').write_bytes(b'old')
        with self.assertRaises(gate.SelectionError): self.create()
        self.assertEqual(b'old',(self.output/'marker').read_bytes())
    def test_late_competitor_is_preserved(self):
        original=gate.os.rename
        def compete(a,b):
            Path(b).mkdir(); (Path(b)/'marker').write_bytes(b'winner'); original(a,b)
        with patch.object(gate.os,'rename',side_effect=compete):
            with self.assertRaises(OSError): self.create()
        self.assertEqual(b'winner',(self.output/'marker').read_bytes()); self.assertEqual([],list(self.root.glob('.pnpm-plan-*')))
    def test_normal_publication_failure_cleans_owned_stage(self):
        with patch.object(gate.os,'rename',side_effect=PermissionError('injected')): self.rejected()
    def test_two_writers_have_one_winner(self):
        def attempt(_):
            try: self.create(); return 'win'
            except (gate.SelectionError,OSError): return 'rejected'
        with concurrent.futures.ThreadPoolExecutor(max_workers=2) as pool: results=list(pool.map(attempt,range(2)))
        self.assertEqual(['rejected','win'],sorted(results)); self.assertEqual([],list(self.root.glob('.pnpm-plan-*')))
    def test_tampered_plan_hash(self):
        sha=self.create(); p=self.output/gate.RECIPE; value=json.loads(p.read_bytes()); value['runtime_admitted']=True; p.write_bytes(gate.json_bytes(value))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output,sha)
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output,gate.digest(p.read_bytes()))
    def test_config_change_is_rejected(self):
        sha=self.create(); (self.output/gate.EMPTY_FILES[0]).write_bytes(b'unsafe=true')
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output,sha)
    def test_extra_file_is_rejected(self):
        sha=self.create(); (self.output/'home/evil.cjs').write_bytes(b'unsafe')
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output,sha)
    def test_missing_directory_is_rejected(self):
        sha=self.create(); (self.output/'home/data').rmdir()
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output,sha)
    def test_drift_after_plan_is_rejected(self):
        sha=self.create(); (self.project/'pnpm-lock.yaml').write_bytes(b'drift')
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output,sha)
    def test_junction_input_rejected(self):
        link=self.root/'junction'; r=subprocess.run(['cmd','/d','/c','mklink','/J',str(link),str(self.project)],capture_output=True)
        self.assertEqual(0,r.returncode,r.stderr.decode(errors='replace'))
        try:self.rejected(project=link)
        finally:link.rmdir()
    def test_invalid_json_cli_is_actionable(self):
        sha=self.create(); p=self.output/gate.RECIPE;p.write_bytes(b'{')
        self.assertEqual(2,gate.main(['verify','--target',str(self.output),'--sha256',gate.digest(p.read_bytes())]))


    def test_notice_delivered_for_all_five_exact_declarations(self):
        sha = self.create(); value = gate.verify_plan(self.output, sha)
        raw = (self.output / gate.NOTICE).read_bytes()
        self.assertIn(b'https://blueoakcouncil.org/license/1.0.0', raw)
        for name, row in self.notice_bindings.items():
            self.assertIn((name + '@' + row['version']).encode(), raw)
            self.assertIn(row['sha256'].encode(), raw)
        self.assertEqual(hashlib.sha256(raw).hexdigest(), value['supplemental_notice']['sha256'])
        self.assertEqual('elite-pnpm-install-plan/v6', value['schema'])

    def test_missing_notice_rejected(self):
        sha = self.create(); (self.output / gate.NOTICE).unlink()
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)

    def test_changed_notice_rejected(self):
        sha = self.create(); p = self.output / gate.NOTICE
        p.write_bytes(p.read_bytes().replace(b'https://blueoakcouncil.org/license/1.0.0', b'https://invalid.example/license'))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)

    def test_forged_notice_and_rehashed_receipt_rejected(self):
        self.create(); p = self.output / gate.NOTICE; p.write_bytes(b'false license')
        receipt = self.output / gate.RECIPE; value = json.loads(receipt.read_bytes())
        value['supplemental_notice']['sha256'] = gate.digest(p.read_bytes())
        receipt.write_bytes(gate.json_bytes(value))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, gate.digest(receipt.read_bytes()))

    def test_notice_manifest_drift_rejected(self):
        p = self.projection / 'payload' / self.notice_bindings['chownr']['path']
        p.write_bytes(p.read_bytes() + b' ')
        self.rejected()

    def test_notice_declaration_mismatch_rejected_even_with_matching_hash(self):
        # Independent semantic guard; projection verifier is mocked only for the fixture.
        for field, changed in [('name', 'another-work'), ('version', '9.0.0'), ('license', 'MIT')]:
            with self.subTest(field=field):
                p = self.projection / 'payload' / self.notice_bindings['chownr']['path']
                original = p.read_bytes(); raw = json.loads(original); raw[field] = changed
                p.write_bytes(gate.json_bytes(raw))
                before = self.notice_bindings['chownr']['sha256']
                self.notice_bindings['chownr']['sha256'] = gate.digest(p.read_bytes())
                try: self.rejected()
                finally: p.write_bytes(original); self.notice_bindings['chownr']['sha256'] = before

    def test_missing_notice_manifest_rejected(self):
        (self.projection / 'payload' / self.notice_bindings['chownr']['path']).unlink()
        self.rejected()

    def test_notice_write_failure_cleans_unpublished_stage(self):
        original = Path.write_bytes
        def fail_notice(path, data):
            if path.name == gate.NOTICE: raise PermissionError('injected notice write fault')
            return original(path, data)
        with patch.object(Path, 'write_bytes', fail_notice): self.rejected()


    def test_qrcode_notice_contains_author_permission_and_all_modules(self):
        sha = self.create(); value = gate.verify_plan(self.output, sha)
        raw = (self.output / gate.QRCODE_NOTICE).read_bytes()
        self.assertIn(b'Copyright (c) 2009 Kazuhiko Arase', raw)
        self.assertIn(b'Permission is hereby granted', raw)
        self.assertIn(b'The above copyright notice and this permission notice', raw)
        self.assertIn(b'THE SOFTWARE IS PROVIDED', raw)
        self.assertIn(gate.QRCODE_HEADER.encode('utf-8'), raw)
        for item in gate.QRCODE_MODULES: self.assertIn(item['path'].encode(), raw)
        self.assertEqual(gate.digest(raw), value['qrcode_notice']['sha256'])
        self.assertTrue((self.output / gate.NOTICE).is_file())

    def test_missing_qrcode_notice_rejected(self):
        sha = self.create(); (self.output / gate.QRCODE_NOTICE).unlink()
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)

    def test_qrcode_author_removed_rejected(self):
        sha = self.create(); p = self.output / gate.QRCODE_NOTICE
        p.write_bytes(p.read_bytes().replace(b'Kazuhiko Arase', b'unknown author'))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)

    def test_qrcode_permission_removed_rejected(self):
        sha = self.create(); p = self.output / gate.QRCODE_NOTICE
        p.write_bytes(p.read_bytes().replace(gate.QRCODE_LICENSE.encode('utf-8'), b'MIT'))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)

    def test_qrcode_notice_and_receipt_rehashed_together_rejected(self):
        self.create(); p = self.output / gate.QRCODE_NOTICE; p.write_bytes(b'Copyright omitted')
        receipt = self.output / gate.RECIPE; value = json.loads(receipt.read_bytes())
        value['qrcode_notice']['sha256'] = gate.digest(p.read_bytes()); receipt.write_bytes(gate.json_bytes(value))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, gate.digest(receipt.read_bytes()))

    def test_qrcode_bundle_missing_rejected(self):
        self.qrcode_bundle.unlink(); self.rejected()

    def test_qrcode_bundle_drift_rejected(self):
        self.qrcode_bundle.write_bytes(self.qrcode_bundle.read_bytes() + b'changed'); self.rejected()

    def test_qrcode_notice_write_failure_cleans_unpublished_stage(self):
        original = Path.write_bytes
        def fail_notice(path, data):
            if path.name == gate.QRCODE_NOTICE: raise PermissionError('injected QRCode notice write fault')
            return original(path, data)
        with patch.object(Path, 'write_bytes', fail_notice): self.rejected()


    def test_semver_original_license_delivered_and_scoped(self):
        sha = self.create(); value = gate.verify_plan(self.output, sha)
        raw = (self.output / gate.SEMVER_NOTICE).read_bytes()
        self.assertIn(gate.SEMVER_LICENSE.encode('utf-8'), raw)
        for text in [b'Copyright 2013 AJ ONeal', b'Permission is hereby granted', b'copies or substantial portions', b'THE SOFTWARE IS PROVIDED', b'APACHEv2', b'expired on 2025-01-29']:
            self.assertIn(text, raw)
        self.assertEqual(gate.digest(raw), value['semver_notice']['sha256'])
        self.assertFalse(value['runtime_admitted']); self.assertFalse(value['redistribution_admitted'])
    def test_missing_semver_notice_rejected(self):
        sha = self.create(); (self.output / gate.SEMVER_NOTICE).unlink()
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)
    def test_semver_copyright_removed_rejected(self):
        sha = self.create(); p = self.output / gate.SEMVER_NOTICE
        p.write_bytes(p.read_bytes().replace(b'Copyright 2013 AJ ONeal', b''))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)
    def test_semver_permission_removed_rejected(self):
        sha = self.create(); p = self.output / gate.SEMVER_NOTICE
        p.write_bytes(p.read_bytes().replace(gate.SEMVER_LICENSE.encode('utf-8'), b'MIT'))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)
    def test_semver_notice_and_receipt_rehashed_together_rejected(self):
        self.create(); p = self.output / gate.SEMVER_NOTICE; p.write_bytes(b'MIT')
        receipt = self.output / gate.RECIPE; value = json.loads(receipt.read_bytes())
        value['semver_notice']['sha256'] = gate.digest(p.read_bytes()); receipt.write_bytes(gate.json_bytes(value))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, gate.digest(receipt.read_bytes()))
    def test_semver_region_drift_rejected(self):
        with patch.dict(gate.SEMVER_REGION, {'region_sha256': '0' * 64}): self.rejected()
    def test_semver_original_license_drift_rejected(self):
        with patch.object(gate, 'SEMVER_LICENSE', 'MIT'): self.rejected()
    def test_semver_notice_write_failure_cleans_unpublished_stage(self):
        original = Path.write_bytes
        def fail_notice(path, data):
            if path.name == gate.SEMVER_NOTICE: raise PermissionError('injected semver notice write fault')
            return original(path, data)
        with patch.object(Path, 'write_bytes', fail_notice): self.rejected()


    def test_retained_collection_delivers_all_exact_texts(self):
        sha = self.create(); value = gate.verify_plan(self.output, sha)
        raw = (self.output / gate.RETAINED_NOTICE).read_bytes()
        for entry in self.catalog['entries']:
            self.assertIn(base64.b64decode(entry['content_base64']), raw)
            self.assertIn(entry['sha256'].encode(), raw)
        self.assertEqual(47, value['retained_notices']['copies'])
        self.assertEqual('RETAINED_TEXTS_ONLY', value['retained_notices']['scope'])
        self.assertFalse(value['redistribution_admitted'])
    def test_retained_delivery_missing_rejected(self):
        sha = self.create(); (self.output / gate.RETAINED_NOTICE).unlink()
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)
    def test_retained_permission_removed_rejected(self):
        sha = self.create(); p = self.output / gate.RETAINED_NOTICE
        p.write_bytes(p.read_bytes().replace(b'Permission', b'Omitted'))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)
    def test_retained_delivery_and_receipt_rehashed_rejected(self):
        self.create(); p = self.output / gate.RETAINED_NOTICE; p.write_bytes(b'All licensed')
        receipt = self.output / gate.RECIPE; value = json.loads(receipt.read_bytes())
        value['retained_notices']['sha256'] = gate.digest(p.read_bytes()); receipt.write_bytes(gate.json_bytes(value))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, gate.digest(receipt.read_bytes()))
    def test_original_retained_notice_drift_rejected(self):
        entry = next(e for e in self.catalog['entries'] if e['origin'] == 'ORIGINAL_SELECTED_PAYLOAD')
        (self.projection / 'payload' / entry['payload_path']).write_bytes(b'changed')
        self.rejected()
    def test_retained_write_failure_cleans_unpublished_stage(self):
        original = Path.write_bytes
        def fail(path, data):
            if path.name == gate.RETAINED_NOTICE: raise PermissionError('injected retained notice write fault')
            return original(path, data)
        with patch.object(Path, 'write_bytes', fail): self.rejected()
    def catalog_rejected(self, mutate):
        value = json.loads(json.dumps(self.catalog)); mutate(value); raw = gate.json_bytes(value)
        original = gate.read_bounded
        def read(path, limit):
            if Path(path).name == 'retained-notice-catalog.json': return raw
            return original(path, limit)
        with patch.object(gate, 'read_bounded', read), patch.object(gate, 'CATALOG_SHA256', gate.digest(raw)):
            self.rejected()
    def test_catalog_duplicate_path_rejected(self):
        self.catalog_rejected(lambda v: v['entries'].__setitem__(1, v['entries'][0]))
    def test_catalog_traversal_rejected(self):
        self.catalog_rejected(lambda v: v['entries'][0].update(path='../outside'))
    def test_catalog_inflated_status_rejected(self):
        self.catalog_rejected(lambda v: v.update(release_ready=True))
    def test_catalog_content_tampered_rejected(self):
        self.catalog_rejected(lambda v: v['entries'][0].update(content_base64='YQ=='))
    def test_catalog_invalid_encoding_rejected(self):
        self.catalog_rejected(lambda v: v['entries'][0].update(content_base64='!!!!'))
    def test_catalog_origin_inflated_rejected(self):
        self.catalog_rejected(lambda v: v['entries'][22].update(origin='ORIGINAL_SELECTED_PAYLOAD'))
    def test_catalog_count_drift_rejected(self):
        self.catalog_rejected(lambda v: v['entries'].pop())
    def test_catalog_extra_field_rejected(self):
        self.catalog_rejected(lambda v: v.update(licensed=True))
    def test_catalog_wrong_hash_rejected(self):
        with patch.object(gate, 'CATALOG_SHA256', '0'*64): self.rejected()


    def test_next_path_original_source_manifest_license_delivered(self):
        sha = self.create(); value = gate.verify_plan(self.output, sha)
        raw = (self.output / gate.NEXT_PATH_NOTICE).read_bytes()
        data = json.loads(Path(gate.__file__).with_name('next-path-source-evidence.json').read_bytes())
        for entry in data['files']: self.assertIn(base64.b64decode(entry['content_base64']), raw)
        self.assertEqual('SOURCE_AND_LICENSE_DELIVERY_ONLY', value['next_path_source']['scope'])
        self.assertIn(b'conditional loader adapter', raw)
        self.assertFalse(value['redistribution_admitted'])
    def test_next_path_source_notice_missing_rejected(self):
        sha = self.create(); (self.output / gate.NEXT_PATH_NOTICE).unlink()
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)
    def test_next_path_source_body_mutation_rejected(self):
        sha = self.create(); p = self.output / gate.NEXT_PATH_NOTICE
        p.write_bytes(p.read_bytes().replace(b'path.relative(from, to)', b'path.relative(to, from)'))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)
    def test_next_path_full_license_omission_rejected(self):
        sha = self.create(); p = self.output / gate.NEXT_PATH_NOTICE
        data = json.loads(Path(gate.__file__).with_name('next-path-source-evidence.json').read_bytes())
        original = base64.b64decode(data['files'][2]['content_base64'])
        p.write_bytes(p.read_bytes().replace(original, b'MPL-2.0'))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)
    def test_next_path_notice_and_receipt_rehashed_rejected(self):
        self.create(); p = self.output / gate.NEXT_PATH_NOTICE; p.write_bytes(b'Source omitted')
        receipt = self.output / gate.RECIPE; value = json.loads(receipt.read_bytes())
        value['next_path_source']['sha256'] = gate.digest(p.read_bytes()); receipt.write_bytes(gate.json_bytes(value))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, gate.digest(receipt.read_bytes()))
    def test_next_path_bundle_region_changed_rejected(self):
        with patch.dict(gate.NEXT_PATH_REGION, {'sha256': '0'*64}): self.rejected()
    def test_next_path_evidence_changed_rejected(self):
        with patch.object(gate, 'NEXT_PATH_DATA_SHA256', '0'*64): self.rejected()
    def test_next_path_write_failure_cleans_stage(self):
        original = Path.write_bytes
        def fail(path, raw):
            if path.name == gate.NEXT_PATH_NOTICE: raise PermissionError('injected next-path source write fault')
            return original(path, raw)
        with patch.object(Path, 'write_bytes', fail): self.rejected()

if __name__ == '__main__': unittest.main()
