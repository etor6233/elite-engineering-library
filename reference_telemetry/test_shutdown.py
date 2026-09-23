"""Owned synthetic work; no business/provider effects or real service restart."""
from pathlib import Path
import ctypes as C, hashlib,json,os,subprocess,sys,threading,time,unittest,uuid
from concurrent.futures import ThreadPoolExecutor
from unittest.mock import patch
import native_capture as n
import governed_launch as g
S=Path(__file__).resolve().parent
ROOT=S/'shutdown-runs';ROOT.mkdir(exist_ok=True)
PRELUDE='''import ctypes as c, os, sys, time
from pathlib import Path
k=c.WinDLL('kernel32',use_last_error=True)
k.WaitForSingleObject.argtypes=[c.c_void_p,c.c_ulong];k.WaitForSingleObject.restype=c.c_ulong
k.SetEvent.argtypes=[c.c_void_p];k.SetEvent.restype=c.c_int
k.ResetEvent.argtypes=[c.c_void_p];k.ResetEvent.restype=c.c_int
h=int(os.environ['ELITE_STOP_EVENT_HANDLE'])
'''
DRAIN=PRELUDE+'''
Path('ready').write_text('in-flight')
assert k.WaitForSingleObject(h,10000)==0
time.sleep(.06)
with open('committed','xb') as f:f.write(b'one-finished-unit');f.flush();os.fsync(f.fileno())
sys.stdout.buffer.write(b'drained-output');sys.stderr.buffer.write(b'drained-error')
assert k.WaitForSingleObject(h,0)==0
'''

class Tests(unittest.TestCase):
    def setUp(self):self.work=ROOT/(self._testMethodName+'-'+uuid.uuid4().hex[:8]);self.work.mkdir()
    def run_child(self,source,**kwargs):return n.capture([sys.executable,'-I','-S','-c',source],self.work,**kwargs)
    def cancel_ready(self,event):
        deadline=time.monotonic()+8
        while not (self.work/'ready').exists() and time.monotonic()<deadline:time.sleep(.005)
        event.set()
    def cancellation(self,source=DRAIN,grace=1):
        event=threading.Event();thread=threading.Thread(target=self.cancel_ready,args=(event,));thread.start()
        try:return self.run_child(source,cancel=event,timeout=10,shutdown_grace=grace)
        finally:thread.join(9);self.assertFalse(thread.is_alive());self.assertTrue((self.work/'ready').exists())
    def test_cancel_drains_inflight_and_captures_final_output(self):
        r=self.cancellation();self.assertEqual((r.status,r.shutdown_state,r.exit_code),('CANCELLED','EXITED_DURING_GRACE',0));self.assertTrue(r.tree_empty);self.assertEqual(r.stdout,b'drained-output');self.assertEqual(r.stderr,b'drained-error');self.assertEqual((self.work/'committed').read_bytes(),b'one-finished-unit')
    def test_timeout_requests_drain_and_retains_timeout(self):
        r=self.run_child(DRAIN,timeout=.3,shutdown_grace=1)
        self.assertEqual((r.status,r.shutdown_state,r.exit_code),('TIMED_OUT','EXITED_DURING_GRACE',0));self.assertTrue(r.tree_empty);self.assertTrue((self.work/'committed').exists())
    def test_noncooperative_child_is_forced_after_grace(self):
        start=time.monotonic();r=self.run_child('import time; time.sleep(30)',timeout=.1,shutdown_grace=.12)
        self.assertEqual((r.status,r.shutdown_state),('TIMED_OUT','FORCED'));self.assertTrue(r.tree_empty);self.assertLess(time.monotonic()-start,6)
    def test_inherited_handle_direct_signal_and_reset_are_denied(self):
        source=PRELUDE+'''
assert not k.SetEvent(h) and c.get_last_error()==5
assert not k.ResetEvent(h) and c.get_last_error()==5
assert k.WaitForSingleObject(h,0)==258
print('wait-only')
'''
        r=self.run_child(source,shutdown_grace=1);self.assertEqual(r.exit_code,0);self.assertEqual(r.stdout,b'wait-only\r\n');self.assertEqual(r.shutdown_state,'NOT_REQUESTED')
    def test_stop_signal_is_issued_once(self):
        original=n.OwnedCapture.request_stop;calls=[]
        def tracked(owner):calls.append(1);return original(owner)
        with patch.object(n.OwnedCapture,'request_stop',tracked):r=self.cancellation()
        self.assertEqual(calls,[1]);self.assertEqual(r.shutdown_state,'EXITED_DURING_GRACE')
    def test_signal_failure_forces_cleanup_with_fixed_code(self):
        with patch.object(n,'set_event',return_value=False):r=self.run_child('import time; time.sleep(30)',timeout=.1,shutdown_grace=1)
        self.assertEqual((r.status,r.shutdown_state),('STOP_SIGNAL_FAILED','SIGNAL_FAILED'));self.assertTrue(r.tree_empty)
    def test_output_limit_during_drain_is_still_enforced(self):
        source=PRELUDE+"\nassert k.WaitForSingleObject(h,10000)==0\nsys.stdout.buffer.write(b'x'*4096);sys.stdout.buffer.flush();time.sleep(30)"
        r=self.run_child(source,timeout=.2,shutdown_grace=1,limit=256)
        self.assertEqual((r.status,r.shutdown_state),('OUTPUT_LIMIT','FORCED'));self.assertTrue(r.tree_empty);self.assertLessEqual(len(r.stdout),256)
    def test_descendant_keeps_grace_open_until_force(self):
        source=PRELUDE+"\nimport subprocess\nsubprocess.Popen([sys.executable,'-I','-S','-c','import time;time.sleep(30)'],creationflags=0x08000000)\nassert k.WaitForSingleObject(h,10000)==0\n"
        r=self.run_child(source,timeout=.25,shutdown_grace=.15)
        self.assertEqual((r.status,r.shutdown_state),('TIMED_OUT','FORCED'));self.assertTrue(r.tree_empty)
    def test_nonzero_drained_exit_is_preserved(self):
        r=self.run_child(PRELUDE+'\nassert k.WaitForSingleObject(h,10000)==0\nsys.exit(7)',timeout=.2,shutdown_grace=1)
        self.assertEqual((r.status,r.shutdown_state,r.exit_code),('TIMED_OUT','EXITED_DURING_GRACE',7))
    def test_precancel_has_no_effect_or_event(self):
        event=threading.Event();event.set()
        with patch.object(n,'create_event') as create:
            r=self.run_child('raise Exception()',cancel=event,shutdown_grace=1)
            create.assert_not_called();self.assertEqual(r.shutdown_state,'NOT_REQUESTED');self.assertTrue(r.tree_empty)
    def test_invalid_grace_rejected_before_creation(self):
        for value in [True,0,-1,61,float('nan'),float('inf'),'1']:
            with self.subTest(value=str(value)),patch.object(n,'create_event') as create:
                with self.assertRaises(ValueError):self.run_child('',shutdown_grace=value)
                create.assert_not_called()
    def test_event_duplicate_failure_releases_handles(self):
        with patch.object(n,'duplicate',return_value=False):r=self.run_child('',shutdown_grace=1)
        self.assertEqual(r.status,'LAUNCH_FAILED');self.assertFalse(r.tree_empty)
    def test_profile_v2_governs_protocol_and_grace(self):
        exe=Path(sys.executable)
        p={'schema':'elite-native-launch-qualification/v2','id':'shutdown-fixture','executable':{'path':str(exe),'bytes':exe.stat().st_size,'sha256':hashlib.sha256(exe.read_bytes()).hexdigest()},'arguments':['-I','-S','-c',DRAIN],'cwd':str(self.work),'environment':{'SYSTEMROOT':os.environ['SYSTEMROOT'],'WINDIR':os.environ['WINDIR'],'TEMP':str(self.work),'TMP':str(self.work)},'budgets':{'timeout':.3,'output_bytes':4096,'processes':8,'commit_bytes':268435456},'shutdown':{'protocol':'win32-inherited-event/v1','grace_seconds':1}}
        raw=json.dumps(p).encode();r=g.launch(raw,hashlib.sha256(raw).hexdigest());self.assertEqual(r.status,'CAPTURED');self.assertEqual(r.outcome.shutdown_state,'EXITED_DURING_GRACE')
        for change in [{'protocol':'unknown','grace_seconds':1},{'protocol':'win32-inherited-event/v1','grace_seconds':0}]:
            p['shutdown']=change;raw=json.dumps(p).encode()
            with patch.object(g.native_capture,'capture') as cap:self.assertEqual(g.launch(raw,hashlib.sha256(raw).hexdigest()).status,'REJECTED');cap.assert_not_called()
    def test_concurrent_events_are_isolated(self):
        source=PRELUDE+'\nassert k.WaitForSingleObject(h,10000)==0\nprint("stopped")'
        with ThreadPoolExecutor(max_workers=8) as pool:results=list(pool.map(lambda _:self.run_child(source,timeout=.2,shutdown_grace=1),range(16)))
        self.assertTrue(all(r.shutdown_state=='EXITED_DURING_GRACE' and r.tree_empty and r.exit_code==0 for r in results))
    def test_repeated_event_lifetimes_do_not_leak_handles(self):
        count=n.win.api('GetProcessHandleCount',[n.W.HANDLE,C.POINTER(n.W.DWORD)],n.W.BOOL)
        def observed():v=n.W.DWORD();self.assertTrue(count(n.current_process(),C.byref(v)));return v.value
        self.run_child('',shutdown_grace=1);before=observed()
        for _ in range(20):self.assertEqual(self.run_child('',shutdown_grace=1).exit_code,0)
        self.assertEqual(observed(),before)
    def test_legacy_protocol_has_no_event_environment(self):
        r=self.run_child('import os; assert "ELITE_STOP_EVENT_HANDLE" not in os.environ')
        self.assertEqual(r.exit_code,0);self.assertEqual(r.shutdown_state,'NOT_REQUESTED')
    def test_grace_exit_does_not_imply_worker_acknowledgement(self):
        r=self.run_child('import time;time.sleep(.2)',timeout=.05,shutdown_grace=1)
        self.assertEqual(r.shutdown_state,'EXITED_DURING_GRACE');self.assertTrue(r.tree_empty)
        self.assertFalse((self.work/'committed').exists())
    def test_same_user_duplication_is_explicitly_outside_security_claim(self):
        source=PRELUDE+'''
k.GetCurrentProcess.restype=c.c_void_p
k.DuplicateHandle.argtypes=[c.c_void_p,c.c_void_p,c.c_void_p,c.POINTER(c.c_void_p),c.c_ulong,c.c_int,c.c_ulong];k.DuplicateHandle.restype=c.c_int
k.CloseHandle.argtypes=[c.c_void_p];k.CloseHandle.restype=c.c_int
out=c.c_void_p()
assert k.DuplicateHandle(k.GetCurrentProcess(),h,k.GetCurrentProcess(),c.byref(out),2,False,0)
assert k.SetEvent(out)
assert k.WaitForSingleObject(h,0)==0
assert k.CloseHandle(out)
print('not-a-hostile-worker-sandbox')
'''
        r=self.run_child(source,shutdown_grace=1);self.assertEqual(r.exit_code,0);self.assertIn(b'not-a-hostile-worker-sandbox',r.stdout)
    def test_owner_death_still_kills_event_waiting_root_and_descendant(self):
        nested="import os,time,pathlib;pathlib.Path('nested.json').write_text(str(os.getpid()));time.sleep(30)"
        child=PRELUDE+f"\nimport subprocess\nsubprocess.Popen([sys.executable,'-I','-S','-c',{nested!r}],creationflags=0x08000000)\nk.WaitForSingleObject(h,30000)\n"
        wrapper=self.work/'owner.py'
        wrapper.write_text("import sys,os,time\nfrom pathlib import Path\nsys.path.insert(0,sys.argv[1]);import native_capture as n\nroot=Path(sys.argv[2])\nowner=n.OwnedCapture([sys.executable,'-I','-S','-c',"+repr(child)+"],root,n.clean_environment(),cooperative=True);owner.start()\nuntil=time.monotonic()+5\nwhile not (root/'nested.json').exists() and time.monotonic()<until:time.sleep(.005)\nn.win.write(root/'ready.json',{'root':owner.pi.pid,'nested':int((root/'nested.json').read_text())})\nn.win.await_file(root/'exit.json');os._exit(0)\n")
        handles=[]
        with (self.work/'owner.stderr').open('wb') as errors:
            process=subprocess.Popen([sys.executable,'-I',str(wrapper),str(S),str(self.work)],stdin=subprocess.DEVNULL,stdout=subprocess.DEVNULL,stderr=errors,creationflags=subprocess.CREATE_NO_WINDOW,env=n.clean_environment())
            try:
                data=n.win.await_file(self.work/'ready.json')
                for key in ['root','nested']:
                    h=n.win.open_process(0x100000|0x1000,False,data[key]);self.assertTrue(h);handles.append(h);self.assertEqual(n.win.wait(h,0),258)
                n.win.write(self.work/'exit.json',{'exit':True});process.wait(timeout=5);self.assertEqual(process.returncode,0)
                for h in handles:self.assertEqual(n.win.wait(h,5000),0)
            finally:
                if process.poll() is None:process.kill();process.wait(timeout=5)
                for h in handles:n.win.close(h)

class Result(unittest.TextTestResult):
    def __init__(self,*a,**kw):super().__init__(*a,**kw);self.outcomes=[]
    def addSuccess(self,t):super().addSuccess(t);self.outcomes.append({'id':t.id(),'status':'PASS'})
    def addFailure(self,t,e):super().addFailure(t,e);self.outcomes.append({'id':t.id(),'status':'FAIL','detail':self._exc_info_to_string(e,t)})
    def addError(self,t,e):super().addError(t,e);self.outcomes.append({'id':t.id(),'status':'ERROR','detail':self._exc_info_to_string(e,t)})
    def addSubTest(self,t,s,e):super().addSubTest(t,s,e);self.outcomes.append({'id':s.id(),'subcase':True,'status':'FAIL' if e else 'PASS'})
if __name__=='__main__':
    result=unittest.TextTestRunner(verbosity=2,resultclass=Result).run(unittest.defaultTestLoader.loadTestsFromTestCase(Tests))
    Path(sys.argv[1]).write_text(json.dumps({'tests_run':result.testsRun,'pass':result.wasSuccessful(),'skips':len(result.skipped),'outcomes':result.outcomes},indent=2)+'\n')
    raise SystemExit(0 if result.wasSuccessful() else 1)
