from pathlib import Path
import hashlib,json,os,subprocess,sys,threading,time,unittest,uuid
from concurrent.futures import ThreadPoolExecutor
from unittest.mock import patch
import result_store as s
S=Path(__file__).resolve().parent;ROOT=S/'result-runs';ROOT.mkdir(exist_ok=True)
FIX=json.loads(Path(sys.argv[2]).read_text());sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()

class Tests(unittest.TestCase):
    def setUp(self):
        self.base=ROOT/uuid.uuid4().hex;self.base.mkdir();self.store=self.base/'store';self.store.mkdir();self.run=uuid.uuid4().hex
        exe=Path(FIX['path']);self.assertEqual(sha(exe),FIX['sha256'])
        self.p={'schema':'elite-native-launch-qualification/v1','id':'local-result-probe','executable':{'path':str(exe),'bytes':exe.stat().st_size,'sha256':sha(exe)},'arguments':[],'cwd':str(self.base),'environment':{'SYSTEMROOT':os.environ['SYSTEMROOT'],'WINDIR':os.environ['WINDIR'],'TEMP':str(self.base),'TMP':str(self.base)},'budgets':{'timeout':5,'output_bytes':4096,'processes':8,'commit_bytes':268435456}}
    def raw(self):raw=s.encode(self.p);return raw,s.digest(raw)
    def run_once(self):return s.execute(self.store,self.run,*self.raw())
    def read(self):return s.inspect(self.store,self.run,self.raw()[1])
    def python(self,code):
        exe=Path(sys.executable);self.p['executable']={'path':str(exe),'sha256':sha(exe),'bytes':exe.stat().st_size};self.p['arguments']=['-I','-S','-c',code]
    def test_real_refund_record_reopens_and_preserves_nonzero_exit(self):
        r=self.run_once();self.assertEqual(r.state,'RECORDED');self.assertEqual(self.read(),r)
        v=json.loads((self.store/self.run/'result.json').read_text());self.assertEqual(v['outcome']['exit_code'],1);self.assertTrue(v['outcome']['tree_empty']);self.assertEqual(v['launch_status'],'CAPTURED')
    def test_no_raw_private_output_arguments_or_environment_persist(self):
        marker='PRIVATE-349-ARG-OUTPUT';self.python('import sys;print(sys.argv[1]);sys.stderr.write(sys.argv[1])');self.p['arguments'].append(marker)
        self.assertEqual(self.run_once().state,'RECORDED')
        for p in (self.store/self.run).iterdir():self.assertNotIn(marker.encode(),p.read_bytes())
        v=json.loads((self.store/self.run/'result.json').read_text());self.assertEqual(v['outcome']['stderr'],{'bytes':len(marker),'sha256':s.digest(marker.encode())})
    def test_existing_result_never_reexecutes_or_overwrites(self):
        self.assertEqual(self.run_once().state,'RECORDED');before=sha(self.store/self.run/'result.json')
        with patch.object(s.g,'launch') as launch:self.assertEqual(self.run_once().state,'ALREADY_RESERVED');launch.assert_not_called()
        self.assertEqual(sha(self.store/self.run/'result.json'),before)
    def test_empty_reservation_is_unresolved_not_replayable(self):
        (self.store/self.run).mkdir();self.assertEqual(self.read().state,'UNRESOLVED')
        with patch.object(s.g,'launch') as launch:self.assertEqual(self.run_once().state,'ALREADY_RESERVED');launch.assert_not_called()
    def test_concurrent_same_id_has_one_actual_launch(self):
        original=s.g.launch;calls=[];lock=threading.Lock()
        def counted(*a,**kw):
            with lock:calls.append(1)
            return original(*a,**kw)
        with patch.object(s.g,'launch',side_effect=counted),ThreadPoolExecutor(max_workers=12) as pool:results=list(pool.map(lambda _:self.run_once(),range(24)))
        self.assertEqual(len(calls),1);self.assertEqual(sum(r.state=='RECORDED' for r in results),1);self.assertEqual(sum(r.state=='ALREADY_RESERVED' for r in results),23)
    def test_invalid_profile_and_path_do_not_reserve(self):
        for run in ['../outside','A'*32,'','0'*33]:self.assertEqual(s.execute(self.store,run,*self.raw()).state,'REJECTED')
        self.assertEqual(s.execute(self.store,self.run,self.raw()[0],'0'*64).state,'REJECTED');self.assertEqual(list(self.store.iterdir()),[])
    def test_claim_flush_failure_prevents_launch(self):
        with patch.object(s.os,'fsync',side_effect=OSError('private-disk-error')),patch.object(s.g,'launch') as launch:
            r=self.run_once();self.assertEqual(r.state,'CLAIM_FAILED');launch.assert_not_called();self.assertNotIn('private',repr(r))
        self.assertEqual(self.read().state,'UNRESOLVED');self.assertEqual(self.run_once().state,'ALREADY_RESERVED')
    def test_short_writes_are_completed_before_publication(self):
        original=Path.open
        class Short:
            def __init__(self,f):self.f=f
            def __enter__(self):return self
            def __exit__(self,*a):self.f.close()
            def fileno(self):return self.f.fileno()
            def write(self,data):return self.f.write(data[:7])
        def opened(path,*a,**kw):
            f=original(path,*a,**kw);return Short(f) if a and a[0]=='xb' and path.name.startswith('part-') else f
        with patch.object(Path,'open',opened):self.assertEqual(self.run_once().state,'RECORDED')
        self.assertEqual(self.read().state,'RECORDED')
    def test_zero_write_fails_before_launch(self):
        original=Path.open
        class Zero:
            def __init__(self,f):self.f=f
            def __enter__(self):return self
            def __exit__(self,*a):self.f.close()
            def write(self,_):return 0
        def opened(path,*a,**kw):
            f=original(path,*a,**kw);return Zero(f) if a and a[0]=='xb' and path.name.startswith('part-') else f
        with patch.object(Path,'open',opened),patch.object(s.g,'launch') as launch:self.assertEqual(self.run_once().state,'CLAIM_FAILED');launch.assert_not_called()
    def test_result_publish_failure_preserves_uncertainty(self):
        original=s.publish_file
        def fail(directory,name,raw):
            if name=='result.json':raise OSError('private-result-error')
            return original(directory,name,raw)
        with patch.object(s,'publish_file',side_effect=fail):self.assertEqual(self.run_once().state,'RESULT_UNCERTAIN')
        self.assertEqual(self.read().state,'UNRESOLVED');self.assertEqual(self.run_once().state,'ALREADY_RESERVED')
    def test_lost_ack_after_publish_is_reconciled_without_reexecution(self):
        original=s.publish_file
        def fail(directory,name,raw):
            value=original(directory,name,raw)
            if name=='result.json':raise OSError('ack-lost')
            return value
        with patch.object(s,'publish_file',side_effect=fail):self.assertEqual(self.run_once().state,'RESULT_UNCERTAIN')
        self.assertEqual(self.read().state,'RECORDED')
        with patch.object(s.g,'launch') as launch:self.assertEqual(self.run_once().state,'ALREADY_RESERVED');launch.assert_not_called()
    def test_unexpected_execution_exception_has_no_private_trace(self):
        with patch.object(s.g,'launch',side_effect=RuntimeError('PRIVATE-EXCEPTION')):r=self.run_once()
        self.assertEqual(r.state,'EXECUTION_UNCERTAIN');self.assertNotIn('PRIVATE',repr(r));self.assertEqual(self.read().state,'UNRESOLVED')
    def test_atomic_publish_refuses_existing_destination(self):
        directory=self.store/self.run;directory.mkdir();(directory/'result.json').write_bytes(b'original')
        with s.DirectoryPin(directory):
            with self.assertRaises(s.g.Rejected):s.publish_file(directory,'result.json',b'new')
        self.assertEqual((directory/'result.json').read_bytes(),b'original')
    def test_partial_and_oversize_final_records_are_invalid(self):
        self.assertEqual(self.run_once().state,'RECORDED');path=self.store/self.run/'result.json';original=path.read_bytes()
        for raw in [b'{',b' '*16385,b'{"schema":1,"schema":2}',b'[]']:
            path.write_bytes(raw);self.assertEqual(self.read().state,'INVALID')
        path.write_bytes(original);self.assertEqual(self.read().state,'RECORDED')
    def test_cross_run_and_claim_hash_tampering_are_invalid(self):
        self.assertEqual(self.run_once().state,'RECORDED');path=self.store/self.run/'result.json';value=json.loads(path.read_text())
        for key,changed in [('run_id','0'*32),('profile_sha256','0'*64),('claim_sha256','0'*64)]:
            altered=dict(value);altered[key]=changed;path.write_bytes(s.encode(altered));self.assertEqual(self.read().state,'INVALID')
    def test_semantic_exit_and_output_metadata_are_checked(self):
        self.assertEqual(self.run_once().state,'RECORDED');path=self.store/self.run/'result.json';value=json.loads(path.read_text())
        for change in [lambda o:o.update(exit_code=True),lambda o:o.update(status='invented'),lambda o:o.update(tree_empty=False),lambda o:o['stdout'].update(sha256='not-a-digest'),lambda o:o.update(extra=1)]:
            altered=json.loads(json.dumps(value));change(altered['outcome']);path.write_bytes(s.encode(altered));self.assertEqual(self.read().state,'INVALID')
    def test_pre_cancel_is_recorded_without_child(self):
        event=threading.Event();event.set();r=s.execute(self.store,self.run,*self.raw(),cancel=event);self.assertEqual(r.state,'RECORDED')
        value=json.loads((self.store/self.run/'result.json').read_text());self.assertEqual(value['launch_status'],'CANCELLED');self.assertIsNone(value['outcome'])
    def test_shutdown_result_is_recorded_without_business_success(self):
        self.python('import time;time.sleep(30)');self.p['schema']='elite-native-launch-qualification/v2';self.p['shutdown']={'protocol':'win32-inherited-event/v1','grace_seconds':.1};self.p['budgets']['timeout']=.1
        self.assertEqual(self.run_once().state,'RECORDED');value=json.loads((self.store/self.run/'result.json').read_text());self.assertEqual(value['outcome']['status'],'TIMED_OUT');self.assertEqual(value['outcome']['shutdown_state'],'FORCED')
    def crash(self,where):
        self.python("from pathlib import Path;Path('effect').open('ab').write(b'one')")
        raw,digest=self.raw();(self.base/'profile.json').write_bytes(raw)
        wrapper=self.base/'crash.py';wrapper.write_text('''import sys,os
from pathlib import Path
sys.path.insert(0,sys.argv[1]);import result_store as s
original=s.publish_file
def publish(directory,name,raw):
    if sys.argv[5]=='before-claim' and name=='claim.json':os._exit(41)
    if sys.argv[5]=='before-result' and name=='result.json':os._exit(42)
    result=original(directory,name,raw)
    if sys.argv[5]=='after-claim' and name=='claim.json':os._exit(43)
    if sys.argv[5]=='after-result' and name=='result.json':os._exit(44)
    return result
s.publish_file=publish
s.execute(Path(sys.argv[2]),sys.argv[3],Path(sys.argv[2]).parent.joinpath('profile.json').read_bytes(),sys.argv[4])
''')
        p=subprocess.run([sys.executable,'-I',str(wrapper),str(S),str(self.store),self.run,digest,where],capture_output=True,timeout=15,creationflags=subprocess.CREATE_NO_WINDOW)
        self.assertIn(p.returncode,[41,42,43,44],p.stderr)
    def test_process_crash_before_claim_leaves_reserved_id(self):
        self.crash('before-claim');self.assertEqual(self.read().state,'UNRESOLVED');self.assertFalse((self.base/'effect').exists());self.assertEqual(self.run_once().state,'ALREADY_RESERVED')
    def test_process_crash_after_claim_never_replays(self):
        self.crash('after-claim');self.assertEqual(self.read().state,'UNRESOLVED');self.assertFalse((self.base/'effect').exists());self.assertEqual(self.run_once().state,'ALREADY_RESERVED')
    def test_process_crash_after_effect_before_result_stays_uncertain(self):
        self.crash('before-result');self.assertEqual(self.read().state,'UNRESOLVED');self.assertEqual((self.base/'effect').read_bytes(),b'one');self.assertEqual(self.run_once().state,'ALREADY_RESERVED');self.assertEqual((self.base/'effect').read_bytes(),b'one')
    def test_process_crash_after_result_recovers_complete_receipt(self):
        self.crash('after-result');self.assertEqual(self.read().state,'RECORDED');self.assertEqual((self.base/'effect').read_bytes(),b'one');self.assertEqual(self.run_once().state,'ALREADY_RESERVED')
    def test_independent_processes_reserve_only_one_execution(self):
        self.python("from pathlib import Path;Path('effect').open('ab').write(b'one')")
        raw,digest=self.raw();(self.base/'profile.json').write_bytes(raw)
        script=self.base/'race.py';script.write_text("import sys\nfrom pathlib import Path\nsys.path.insert(0,sys.argv[1]);import result_store as s\nprint(s.execute(Path(sys.argv[2]),sys.argv[3],Path(sys.argv[2]).parent.joinpath('profile.json').read_bytes(),sys.argv[4]).state)\n")
        def run(_):
            p=subprocess.run([sys.executable,'-I',str(script),str(S),str(self.store),self.run,digest],capture_output=True,timeout=15,creationflags=subprocess.CREATE_NO_WINDOW)
            self.assertEqual(p.returncode,0,p.stderr);return p.stdout.strip()
        with ThreadPoolExecutor(max_workers=8) as pool:results=list(pool.map(run,range(8)))
        self.assertEqual(results.count(b'RECORDED'),1);self.assertEqual(results.count(b'ALREADY_RESERVED'),7);self.assertEqual((self.base/'effect').read_bytes(),b'one')
    def test_renamed_store_and_reparse_slot_are_refused(self):
        # Native directory pins hold the selected parent while execution runs.
        original=s.g.launch
        def guarded(*a,**kw):
            with self.assertRaises(OSError):self.store.rename(self.base/'swapped')
            return original(*a,**kw)
        with patch.object(s.g,'launch',side_effect=guarded):self.assertEqual(self.run_once().state,'RECORDED')
        target=self.base/'elsewhere';target.mkdir();self.run=uuid.uuid4().hex
        pwsh='C:/Users/NL/.cache/codex-runtimes/codex-primary-runtime/dependencies/native/powershell/pwsh.exe'
        command="New-Item -ItemType Junction -Path '"+str(self.store/self.run).replace("'","''")+"' -Target '"+str(target).replace("'","''")+"' | Out-Null"
        p=subprocess.run([pwsh,'-NoProfile','-Command',command],capture_output=True,timeout=15,creationflags=subprocess.CREATE_NO_WINDOW);self.assertEqual(p.returncode,0,p.stderr)
        self.assertEqual(self.read().state,'INVALID')
        with patch.object(s.g,'launch') as launch:self.assertEqual(self.run_once().state,'ALREADY_RESERVED');launch.assert_not_called()
        self.assertEqual(list(target.iterdir()),[])
    def test_move_error_after_success_is_uncertain_but_readable(self):
        original=s.move_file
        def uncertain(source,target,flags):
            result=original(source,target,flags)
            return False if str(target).endswith('result.json') else result
        with patch.object(s,'move_file',side_effect=uncertain):self.assertEqual(self.run_once().state,'RESULT_UNCERTAIN')
        self.assertEqual(self.read().state,'RECORDED')
    def test_writer_readback_never_exposes_partial_public_record(self):
        original=s.publish_file;observations=[]
        def checked(directory,name,raw):
            observations.append(s.inspect(self.store,self.run,self.raw()[1]).state)
            result=original(directory,name,raw)
            observations.append(s.inspect(self.store,self.run,self.raw()[1]).state)
            return result
        with patch.object(s,'publish_file',side_effect=checked):self.assertEqual(self.run_once().state,'RECORDED')
        self.assertEqual(observations,['UNRESOLVED','UNRESOLVED','UNRESOLVED','RECORDED'])
    def test_reader_rejects_image_byte_boolean(self):
        self.assertEqual(self.run_once().state,'RECORDED');p=self.store/self.run/'result.json';value=json.loads(p.read_text());value['image']['bytes']=True;p.write_bytes(s.encode(value));self.assertEqual(self.read().state,'INVALID')
    def test_success_and_invalid_reads_release_native_handles(self):
        self.assertEqual(self.run_once().state,'RECORDED')
        count=s.g.win.api('GetProcessHandleCount',[s.g.W.HANDLE,s.g.C.POINTER(s.g.W.DWORD)],s.g.W.BOOL)
        current=s.g.win.api('GetCurrentProcess',[],s.g.W.HANDLE)
        def observed():v=s.g.W.DWORD();self.assertTrue(count(current(),s.g.C.byref(v)));return v.value
        self.read();before=observed()
        for _ in range(24):self.assertEqual(self.read().state,'RECORDED')
        p=self.store/self.run/'result.json';p.write_bytes(b'{')
        for _ in range(24):self.assertEqual(self.read().state,'INVALID')
        self.assertEqual(observed(),before)

class Result(unittest.TextTestResult):
    def __init__(self,*a,**kw):super().__init__(*a,**kw);self.outcomes=[]
    def addSuccess(self,t):super().addSuccess(t);self.outcomes.append({'id':t.id(),'status':'PASS'})
    def addFailure(self,t,e):super().addFailure(t,e);self.outcomes.append({'id':t.id(),'status':'FAIL','detail':self._exc_info_to_string(e,t)})
    def addError(self,t,e):super().addError(t,e);self.outcomes.append({'id':t.id(),'status':'ERROR','detail':self._exc_info_to_string(e,t)})
if __name__=='__main__':
    result=unittest.TextTestRunner(verbosity=2,resultclass=Result).run(unittest.defaultTestLoader.loadTestsFromTestCase(Tests))
    Path(sys.argv[1]).write_text(json.dumps({'tests_run':result.testsRun,'pass':result.wasSuccessful(),'skips':len(result.skipped),'outcomes':result.outcomes},indent=2)+'\n')
    raise SystemExit(0 if result.wasSuccessful() else 1)
