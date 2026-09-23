"""Closed v3 reference environment and actual input-rejection regressions."""
from pathlib import Path
from unittest.mock import patch
import copy,hashlib,json,os,subprocess,sys,tempfile,unittest
import governed_launch as g
import native_capture as n
class Tests(unittest.TestCase):
 def setUp(self):
  self.tmp=tempfile.TemporaryDirectory(prefix='elite-reference-policy-');self.addCleanup(self.tmp.cleanup);self.root=Path(self.tmp.name)
  self.v={'schema':'elite-native-launch-qualification/v3','id':'reference','executable':{'path':str(self.root/'host.exe'),'bytes':1,'sha256':'a'*64},'arguments':[],'cwd':str(self.root),'environment':{'SYSTEMROOT':os.environ['SYSTEMROOT'],'WINDIR':os.environ['WINDIR'],'TEMP':str(self.root),'TMP':str(self.root),'DATABASE_URL':'postgres://postgres@127.0.0.1:15432/elite_refund_test_'+'a'*32+'?sslmode=disable','REFUND_WORKER_ID':'elite-reference-test','REFUND_TELEMETRY_CONFIG':str(self.root/'config.json'),'REFUND_TELEMETRY_CONFIG_SHA256':'b'*64},'budgets':{'timeout':3,'output_bytes':1024,'processes':1,'commit_bytes':67108864},'shutdown':{'protocol':'win32-inherited-event/v1','grace_seconds':1}}
 def raw(self,v):
  b=json.dumps(v).encode();return b,hashlib.sha256(b).hexdigest()
 def test_reference_profile_accepts_only_explicit_shape(self):self.assertEqual(g.profile(*self.raw(self.v)),self.v)
 def test_remote_credentials_and_query_ambiguity_rejected(self):
  base=self.v['environment']['DATABASE_URL']
  for dsn in [base.replace('127.0.0.1','localhost'),base.replace('127.0.0.1','example.test'),base.replace('postgres@','postgres:secret@'),base.replace('15432','80'),base+'&sslmode=verify-full',base+'#fragment',base.replace('postgres://','https://'),base.replace('elite_refund_test_','existing_')]:
   with self.subTest(dsn=dsn):
    v=copy.deepcopy(self.v);v['environment']['DATABASE_URL']=dsn
    with self.assertRaises((g.Rejected,ValueError)):g.profile(*self.raw(v))
 def test_environment_cannot_inject_provider_path_or_stop_handle(self):
  for key in ['PATH','STRIPE_SECRET_KEY','MERCADO_PAGO_ACCESS_TOKEN','ELITE_STOP_EVENT_HANDLE']:
   with self.subTest(key=key):
    v=copy.deepcopy(self.v);v['environment'][key]='synthetic-secret'
    with patch.object(g.native_capture,'capture') as c:self.assertEqual(g.launch(*self.raw(v)).status,'REJECTED');c.assert_not_called()
 def test_old_profile_cannot_enable_runtime_environment(self):
  for version in ['v1','v2']:
   v=copy.deepcopy(self.v);v['schema']='elite-native-launch-qualification/'+version
   if version=='v1':v.pop('shutdown')
   with self.assertRaises(g.Rejected):g.profile(*self.raw(v))
 def test_capture_default_and_reference_shape_reject_before_child(self):
  bad=copy.deepcopy(self.v['environment']);bad['PATH']='unsafe'
  for env,flag in [(self.v['environment'],False),(bad,True),(None,True),(self.v['environment'],'yes')]:
   with self.subTest(flag=flag),patch.object(n,'OwnedCapture') as launch:
    with self.assertRaises(ValueError):n.capture([str(self.root/'host.exe')],self.root,environment=env,reference_environment=flag)
    launch.assert_not_called()
 def test_runtime_lock_digest_rejected_before_work_directory(self):
  p=self.root/'lock.json';p.write_text('{"synthetic":"private-marker"}')
  r=subprocess.run([sys.executable,'-B',str(Path(__file__).with_name('run_reference_runtime.py')),'--runtime-lock',str(p),'--runtime-lock-sha256','0'*64,'--consumer',str(self.root),'--work-root',str(self.root)],capture_output=True,timeout=10)
  self.assertEqual(r.returncode,2);self.assertEqual(r.stderr.strip(),b'REFERENCE_INPUT_REJECTED');self.assertNotIn(b'private-marker',r.stdout+r.stderr);self.assertEqual(list(self.root.iterdir()),[p])
if __name__=='__main__':
 result=unittest.TextTestRunner(verbosity=2).run(unittest.defaultTestLoader.loadTestsFromTestCase(Tests))
 if len(sys.argv)>1:Path(sys.argv[1]).write_text(json.dumps({'pass':result.wasSuccessful(),'tests':result.testsRun,'skips':len(result.skipped)}))
 raise SystemExit(0 if result.wasSuccessful() else 1)
