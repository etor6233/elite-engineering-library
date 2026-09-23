"""Real owned local files/processes; no provider, account or service changes."""
from pathlib import Path
import copy, ctypes as C, errno, hashlib, json, os, shutil, subprocess, sys, threading, unittest, uuid
from concurrent.futures import ThreadPoolExecutor
from unittest.mock import patch
import governed_launch as g

S=Path(__file__).resolve().parent
FIXTURE=json.loads(Path(sys.argv[2]).read_text())
ROOT=S/'launch-runs';ROOT.mkdir(exist_ok=True)
sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()

class Tests(unittest.TestCase):
    def setUp(self):
        self.work=ROOT/(self._testMethodName+'-'+uuid.uuid4().hex[:8]);self.work.mkdir()
        self.exe=self.work/'owned.exe';shutil.copyfile(FIXTURE['path'],self.exe)
        self.assertEqual(sha(self.exe),FIXTURE['sha256'])
        self.v={'schema':'elite-native-launch-qualification/v1','id':'owned-refund-probe','executable':{'path':str(self.exe),'sha256':sha(self.exe),'bytes':self.exe.stat().st_size},'arguments':[],'cwd':str(self.work),'environment':{'SYSTEMROOT':os.environ['SYSTEMROOT'],'WINDIR':os.environ['WINDIR'],'TEMP':str(self.work),'TMP':str(self.work)},'budgets':{'timeout':5,'output_bytes':4096,'processes':8,'commit_bytes':268435456}}
    def raw(self,v=None):
        raw=json.dumps(self.v if v is None else v,sort_keys=True).encode();return raw,hashlib.sha256(raw).hexdigest()
    def test_actual_hashed_refund_exit_and_output(self):
        r=g.launch(*self.raw());self.assertEqual(r.status,'CAPTURED');self.assertEqual(r.outcome.status,'COMPLETED');self.assertEqual(r.outcome.exit_code,1);self.assertIn(b'REFUND_HOST_FAILED',r.outcome.stderr);self.assertTrue(r.outcome.tree_empty);self.assertEqual(r.image['sha256'],self.v['executable']['sha256'])
    def test_profile_digest_mismatch_prevents_launch(self):
        with patch.object(g.native_capture,'capture') as capture:
            self.assertEqual(g.launch(self.raw()[0],'0'*64).status,'REJECTED');capture.assert_not_called()
    def test_duplicate_unknown_or_malformed_fields_reject(self):
        samples=[b'{"schema":1,"schema":2}',b'\xff',b'{',b'[]',b'{}',b'NaN',b' '*65537]
        for change in [lambda v:v.update(extra=True),lambda v:v['budgets'].update(processes=True),lambda v:v['environment'].update(PATH='C:\\bad'),lambda v:v['executable'].update(bytes=0),lambda v:v.update(arguments=['bad\0arg']),lambda v:v['budgets'].update(timeout=float('inf'))]:
            v=copy.deepcopy(self.v);change(v);samples.append(json.dumps(v).encode())
        for raw in samples:
            with self.subTest(raw_length=len(raw)),patch.object(g.native_capture,'capture') as c:
                self.assertEqual(g.launch(raw,hashlib.sha256(raw).hexdigest()).status,'REJECTED');c.assert_not_called()
    def test_ambiguous_paths_refused(self):
        for path in ['relative.exe','C:relative.exe','\\\\server\\share\\a.exe','\\\\?\\C:\\a.exe','C:\\x\\..\\a.exe','C:\\a.exe:stream','C:\\NUL.exe','C:\\x.\\a.exe','C:/a.exe']:
            with self.subTest(path=path):
                with self.assertRaises(g.Rejected):g.local_path(path)
    def test_wrong_binary_digest_and_size_prevent_launch(self):
        for key,value in [('sha256','0'*64),('bytes',1)]:
            v=copy.deepcopy(self.v);v['executable'][key]=value
            with self.subTest(key=key),patch.object(g.native_capture,'capture') as c:
                self.assertEqual(g.launch(*self.raw(v)).status,'REJECTED');c.assert_not_called()
        with self.exe.open('r+b') as f:f.write(b'X') # rejected hash path released every handle
    def test_existing_writer_refuses_identity_lock(self):
        with self.exe.open('r+b'):
            with patch.object(g.native_capture,'capture') as c:
                self.assertEqual(g.launch(*self.raw()).status,'REJECTED');c.assert_not_called()
    def test_write_delete_rename_blocked_until_close(self):
        pin=g.PinnedImage(self.v['executable'],str(self.work))
        try:
            h=g.create_file(str(self.exe),0x40000000,7,None,3,0,None)
            self.assertEqual(h,C.c_void_p(-1).value)
            self.assertEqual(C.get_last_error(),32)
            for action in [lambda:self.exe.open('r+b'),lambda:self.exe.unlink(),lambda:self.exe.rename(self.work/'changed.exe'),lambda:self.work.rename(self.work.with_name('renamed'))]:
                with self.assertRaises(OSError) as error:action()
                if error.exception.winerror is None:self.assertEqual(error.exception.errno,errno.EACCES)
                else:self.assertIn(error.exception.winerror,[5,32])
        finally:self.assertTrue(pin.close())
        renamed=self.exe.rename(self.work/'after-close.exe');self.assertTrue(renamed.is_file())
    def test_hardlink_writer_is_blocked(self):
        alias=self.work/'alias.exe';os.link(self.exe,alias)
        pin=g.PinnedImage(self.v['executable'],str(self.work))
        try:
            with self.assertRaises(OSError):alias.open('r+b')
        finally:pin.close()
    def test_real_launch_keeps_identity_locks_through_capture(self):
        original=g.native_capture.capture
        def guarded(*args,**kwargs):
            with self.assertRaises(OSError):self.exe.open('r+b')
            with self.assertRaises(OSError):self.work.rename(self.work.with_name('swap'))
            return original(*args,**kwargs)
        with patch.object(g.native_capture,'capture',side_effect=guarded):r=g.launch(*self.raw())
        self.assertEqual(r.outcome.exit_code,1);self.assertTrue(r.outcome.tree_empty)
    def test_explicit_environment_and_exact_argv(self):
        py=Path(sys.executable);self.v['executable']={'path':str(py),'sha256':sha(py),'bytes':py.stat().st_size}
        values=['space here','quote"here','trailing\\','&& ignored | > nope']
        self.v['arguments']=['-I','-S','-c','import os,sys,json; print(json.dumps([sys.argv[1:], os.getenv("TEMP"), os.getenv("V347_PRIVATE"), os.getenv("PATH")]))',*values]
        with patch.dict(os.environ,{'TEMP':'C:\\wrong-parent','V347_PRIVATE':'do-not-inherit','PATH':'do-not-inherit'}):r=g.launch(*self.raw())
        self.assertEqual(r.outcome.exit_code,0);self.assertEqual(json.loads(r.outcome.stdout),[values,str(self.work),None,None])
    def test_precancel_never_opens_identity_or_creates_process(self):
        event=threading.Event();event.set()
        with patch.object(g,'PinnedImage') as pin:
            self.assertEqual(g.launch(*self.raw(),cancel=event).status,'CANCELLED');pin.assert_not_called()
    def test_timeout_propagates_owned_tree_cleanup(self):
        py=Path(sys.executable);self.v['executable']={'path':str(py),'sha256':sha(py),'bytes':py.stat().st_size};self.v['arguments']=['-I','-S','-c','import time; time.sleep(30)'];self.v['budgets']['timeout']=.15
        r=g.launch(*self.raw());self.assertEqual(r.outcome.status,'TIMED_OUT');self.assertTrue(r.outcome.tree_empty)
    def test_missing_file_refused_without_process(self):
        self.exe.rename(self.work/'unavailable.exe')
        with patch.object(g.native_capture,'capture') as c:self.assertEqual(g.launch(*self.raw()).status,'REJECTED');c.assert_not_called()
    def test_directory_reparse_junction_refused(self):
        junction=self.work/'junction'
        pwsh='C:/Users/NL/.cache/codex-runtimes/codex-primary-runtime/dependencies/native/powershell/pwsh.exe'
        command="New-Item -ItemType Junction -Path '"+str(junction).replace("'","''")+"' -Target '"+str(self.work).replace("'","''")+"' | Out-Null"
        result=subprocess.run([pwsh,'-NoProfile','-Command',command],capture_output=True,timeout=15,creationflags=subprocess.CREATE_NO_WINDOW)
        self.assertEqual(result.returncode,0,result.stderr)
        self.v['executable']['path']=str(junction/'owned.exe')
        with patch.object(g.native_capture,'capture') as c:self.assertEqual(g.launch(*self.raw()).status,'REJECTED');c.assert_not_called()
    def test_identity_handles_are_noninheritable(self):
        get_info=g.win.api('GetHandleInformation',[g.W.HANDLE,C.POINTER(g.W.DWORD)],g.W.BOOL)
        pin=g.PinnedImage(self.v['executable'],str(self.work))
        try:
            for h in pin.handles:
                flags=g.W.DWORD();self.assertTrue(get_info(h,C.byref(flags)));self.assertEqual(flags.value&1,0)
        finally:pin.close()
    def test_capture_error_releases_identity(self):
        with patch.object(g.native_capture,'capture',side_effect=OSError('private-value')):
            result=g.launch(*self.raw());self.assertEqual(result.status,'REJECTED');self.assertNotIn('private-value',repr(result))
        with self.exe.open('r+b'):pass
    def test_cleanup_failure_never_becomes_captured(self):
        original=g.PinnedImage.close
        def failed(pin):original(pin);return False
        with patch.object(g.PinnedImage,'close',failed):r=g.launch(*self.raw())
        self.assertEqual(r.status,'IDENTITY_CLEANUP_FAILED')
    def test_invalid_environment_rejected_by_capture_api(self):
        with self.assertRaises(ValueError):g.native_capture.capture([str(self.exe)],self.work,environment={'PATH':'private'})
    def test_concurrent_read_pins_and_captures_do_not_conflict(self):
        raw,digest=self.raw()
        with ThreadPoolExecutor(max_workers=8) as pool:results=list(pool.map(lambda _:g.launch(raw,digest),range(16)))
        self.assertTrue(all(r.status=='CAPTURED' and r.outcome.status=='COMPLETED' and r.outcome.exit_code==1 and r.outcome.tree_empty for r in results))
        with self.exe.open('r+b'):pass
    def test_failed_hash_releases_all_native_handles(self):
        count=g.win.api('GetProcessHandleCount',[g.W.HANDLE,C.POINTER(g.W.DWORD)],g.W.BOOL)
        current=g.win.api('GetCurrentProcess',[],g.W.HANDLE)
        def observed():
            n=g.W.DWORD();self.assertTrue(count(current(),C.byref(n)));return n.value
        self.v['executable']['sha256']='0'*64
        for _ in range(3):self.assertEqual(g.launch(*self.raw()).status,'REJECTED')
        before=observed()
        for _ in range(24):self.assertEqual(g.launch(*self.raw()).status,'REJECTED')
        self.assertEqual(observed(),before)

class Result(unittest.TextTestResult):
    def __init__(self,*a,**kw):super().__init__(*a,**kw);self.observations=[]
    def addSuccess(self,test):super().addSuccess(test);self.observations.append({'id':test.id(),'status':'PASS'})
    def addFailure(self,test,err):super().addFailure(test,err);self.observations.append({'id':test.id(),'status':'FAIL','detail':self._exc_info_to_string(err,test)})
    def addError(self,test,err):super().addError(test,err);self.observations.append({'id':test.id(),'status':'ERROR','detail':self._exc_info_to_string(err,test)})
    def addSubTest(self,test,subtest,err):
        super().addSubTest(test,subtest,err)
        self.observations.append({'id':subtest.id(),'subcase':True,'status':'FAIL' if err else 'PASS'})
if __name__=='__main__':
    result=unittest.TextTestRunner(verbosity=2,resultclass=Result).run(unittest.defaultTestLoader.loadTestsFromTestCase(Tests))
    Path(sys.argv[1]).write_text(json.dumps({'tests_run':result.testsRun,'pass':result.wasSuccessful(),'skips':len(result.skipped),'outcomes':result.observations},indent=2)+'\n')
    raise SystemExit(0 if result.wasSuccessful() else 1)
