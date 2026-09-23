"""Mutation tests for local activation integrity; native connected test is separate."""
from pathlib import Path
import hashlib,json,tempfile,unittest
from local_release import Deployment,runtime_config,verify_release,raw_json,sha,tree
class LocalReleaseTests(unittest.TestCase):
    def test_external_digest_required(self):
        with tempfile.TemporaryDirectory()as d:
            p=Path(d);f=p/'manifest.json';f.write_text('{}')
            with self.assertRaisesRegex(ValueError,'digest'):verify_release(p,f,'0'*64)
    def test_remote_database_and_unknown_env_rejected(self):
        base={'database_url':'postgres://postgres@127.0.0.1:6543/elite_payment_connected_'+'a'*32+'?sslmode=disable','issuer':'http://127.0.0.1:6544','api_port':6545,'web_port':6546}
        self.assertEqual(runtime_config(base),base)
        self.assertEqual(runtime_config(base|{'metrics_port':6547})['metrics_port'],6547)
        for port in [True,0,6545,6546,65536,'6547']:
            with self.subTest(metrics_port=port),self.assertRaises(ValueError):runtime_config(base|{'metrics_port':port})
        for patch in [{'database_url':base['database_url'].replace('127.0.0.1','example.com')},{'issuer':'http://example.com:6544'},{'api_port':6546},{'provider_key':'unexpected'}]:
            with self.subTest(patch=patch),self.assertRaises(ValueError):runtime_config(base|patch)
    def test_exclusive_operator(self):
        with tempfile.TemporaryDirectory()as d:
            cfg={'database_url':'postgres://postgres@127.0.0.1:6543/elite_payment_connected_'+'a'*32+'?sslmode=disable','issuer':'http://127.0.0.1:6544','api_port':6545,'web_port':6546}
            first=Deployment(Path(d),cfg)
            try:
                with self.assertRaisesRegex(ValueError,'another local operator'):Deployment(Path(d),cfg)
            finally:first.close()
if __name__=='__main__':unittest.main()
