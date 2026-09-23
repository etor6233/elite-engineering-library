from pathlib import Path
import ctypes as C
from ctypes import wintypes as W
import hashlib,json,os,sys,time,threading,unittest,tempfile,subprocess
from unittest.mock import patch
import native_capture as n
from native_capture import win

class CaptureTests(unittest.TestCase):
    def setUp(self):
        self.temp=tempfile.TemporaryDirectory();self.root=Path(self.temp.name).resolve()
    def tearDown(self):self.temp.cleanup()
    def call(self,code,**kw):return n.capture([sys.executable,'-I','-c',code],self.root,**kw)
    def assert_complete(self,result,code=0):
        self.assertEqual('COMPLETED',result.status);self.assertEqual(code,result.exit_code);self.assertTrue(result.tree_empty)
    def test_separate_binary_streams_at_exact_limit(self):
        result=self.call("import os;[(os.write(1,bytes(range(256))*16),os.write(2,bytes(reversed(range(256)))*16)) for _ in range(256)]",limit=1024*1024)
        self.assert_complete(result)
        self.assertEqual(bytes(range(256))*4096,result.stdout);self.assertEqual(bytes(reversed(range(256)))*4096,result.stderr)
    def test_each_stream_over_limit_stops_tree_before_later_marker(self):
        for fd in (1,2):
            with self.subTest(fd=fd):
                result=self.call(f"import os,time,pathlib;os.write({fd},b'x'*65536);time.sleep(.4);pathlib.Path('later').write_text('bad')",limit=1024)
                self.assertEqual('OUTPUT_LIMIT',result.status);self.assertTrue(result.tree_empty);self.assertLessEqual(len(result.stdout),1024);self.assertLessEqual(len(result.stderr),1024)
                self.assertFalse((self.root/'later').exists())
    def test_stdin_eof_and_unselected_environment_absent(self):
        with patch.dict(os.environ,{'ELITE_CAPTURE_PRIVATE':'DO_NOT_INHERIT'}):
            result=self.call("import sys,os;assert sys.stdin.buffer.read()==b'';assert 'ELITE_CAPTURE_PRIVATE' not in os.environ;print('ok')")
        self.assert_complete(result);self.assertIn(b'ok',result.stdout)
    def test_nonzero_and_still_active_numeric_exit_preserved(self):
        for code in (7,259):
            with self.subTest(code=code):self.assert_complete(self.call(f'raise SystemExit({code})'),code)
    def test_timeout_cleans_child_and_descendant_with_open_pipe(self):
        nested="import os,time,pathlib;pathlib.Path('nested-ready').write_text(str(os.getpid()));time.sleep(20)"
        original=n.OwnedCapture.close;observed=[]
        def close_with_owned_descendant(owner):
            handle=win.open_process(0x100000|0x1000,False,int((self.root/'nested-ready').read_text()))
            self.assertTrue(handle)
            try:
                membership=W.BOOL();self.assertTrue(win.in_job(handle,owner.job,C.byref(membership)));self.assertTrue(membership.value)
                self.assertEqual(win.wait(handle,0),258)
                result=original(owner)
                self.assertEqual(win.wait(handle,5000),0);observed.append('owned-descendant-signaled')
                return result
            finally:win.close(handle)
        with patch.object(n.OwnedCapture,'close',close_with_owned_descendant):
            result=self.call(f"import subprocess,sys,time;subprocess.Popen([sys.executable,'-I','-c',{nested!r}]);time.sleep(20)",timeout=.6)
        self.assertEqual(observed,['owned-descendant-signaled'])
        self.assertTrue((self.root/'nested-ready').exists());self.assertEqual('TIMED_OUT',result.status);self.assertTrue(result.tree_empty)
        # Native tree accounting is distinct from post-termination file release.
        # Observe only this owned fixture, under a finite deadline, after the
        # exact descendant handle has signaled. Never waive expiry or a live child.
        until=time.monotonic()+5;moved=self.root.with_name(self.root.name+'-released')
        while True:
            try:self.root.rename(moved);break
            except OSError as error:
                if error.winerror not in (5,32) or time.monotonic()>=until:raise
                time.sleep(.005)
        moved.rename(self.root)
    def test_silent_descendant_cannot_be_success_when_root_exits(self):
        nested="import time,pathlib,os;pathlib.Path('silent-ready').write_text(str(os.getpid()));time.sleep(20)"
        original=n.OwnedCapture.close
        observed=[]
        def close_with_owned_descendant(owner):
            pid=int((self.root/'silent-ready').read_text())
            handle=win.open_process(0x100000|0x1000,False,pid)
            self.assertTrue(handle)
            try:
                membership=W.BOOL();self.assertTrue(win.in_job(handle,owner.job,C.byref(membership)));self.assertTrue(membership.value)
                self.assertEqual(win.wait(handle,0),258)
                result=original(owner)
                self.assertEqual(win.wait(handle,5000),0)
                observed.append('owned-descendant-signaled')
                return result
            finally:win.close(handle)
        with patch.object(n.OwnedCapture,'close',close_with_owned_descendant):
            result=self.call(f"import subprocess,sys;subprocess.Popen([sys.executable,'-I','-c',{nested!r}],stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL)",timeout=.6)
        self.assertEqual(observed,['owned-descendant-signaled'])
        self.assertTrue((self.root/'silent-ready').exists());self.assertEqual('TIMED_OUT',result.status);self.assertTrue(result.tree_empty)
        # Fixture teardown is a separate resource-release observation, not an
        # inference from job accounting. Retry only this exact owned rename,
        # after proving descendant termination, under a finite deadline.
        until=time.monotonic()+5;moved=self.root.with_name(self.root.name+'-released')
        while True:
            try:self.root.rename(moved);break
            except OSError as error:
                if error.winerror not in (5,32) or time.monotonic()>=until:raise
                time.sleep(.005)
        moved.rename(self.root)
    def test_finite_descendant_output_after_root_exit_is_collected(self):
        nested="import time,os;time.sleep(.1);os.write(1,b'late-child');os.write(2,b'late-error')"
        result=self.call(f"import subprocess,sys;subprocess.Popen([sys.executable,'-I','-c',{nested!r}])")
        self.assert_complete(result);self.assertEqual(b'late-child',result.stdout);self.assertEqual(b'late-error',result.stderr)
    def test_pre_cancel_does_not_launch(self):
        cancel=threading.Event();cancel.set()
        with patch.object(n,'OwnedCapture') as launch:result=self.call('raise SystemExit(0)',cancel=cancel)
        launch.assert_not_called();self.assertEqual('CANCELLED',result.status);self.assertTrue(result.tree_empty)
    def test_active_cancel_cleans_tree(self):
        cancel=threading.Event()
        def trigger():
            until=time.monotonic()+5
            while not (self.root/'ready').exists() and time.monotonic()<until:time.sleep(.005)
            cancel.set()
        thread=threading.Thread(target=trigger);thread.start()
        try:result=self.call("from pathlib import Path;import time;Path('ready').write_text('yes');time.sleep(20)",cancel=cancel)
        finally:thread.join(timeout=6)
        self.assertFalse(thread.is_alive());self.assertEqual('CANCELLED',result.status);self.assertTrue(result.tree_empty)
    def test_read_failure_is_static_and_cleans(self):
        read=os.read
        def injected(fd,count):raise OSError('PRIVATE_PATH_AND_TOKEN')
        with patch.object(n.os,'read',side_effect=injected):result=self.call('import time;time.sleep(20)')
        self.assertEqual('IO_FAILED',result.status);self.assertTrue(result.tree_empty);self.assertNotIn('PRIVATE_PATH_AND_TOKEN',repr(result))
    def test_unlisted_inheritable_handle_is_not_in_child(self):
        # Give the child the numeric handle of an intentionally inheritable file.
        # GetFinalPathName distinguishes accidental handle-value reuse from leak.
        sentinel=self.root/'not-inherited.txt';sentinel.write_text('PRIVATE_SENTINEL')
        fd=os.open(sentinel,os.O_RDONLY|os.O_BINARY);os.set_inheritable(fd,True);h=msvcrt_handle(fd)
        code=f"""import ctypes as c
from ctypes import wintypes as w
k=c.WinDLL('kernel32',use_last_error=True)
f=k.GetFinalPathNameByHandleW;f.argtypes=[w.HANDLE,w.LPWSTR,w.DWORD,w.DWORD];f.restype=w.DWORD
b=c.create_unicode_buffer(4096);r=f({h},b,len(b),0)
assert r==0 or 'not-inherited.txt' not in b.value
print('excluded')
"""
        try:result=self.call(code)
        finally:os.close(fd)
        self.assert_complete(result);self.assertIn(b'excluded',result.stdout)
    def test_concurrent_capture_has_isolated_streams_and_eof(self):
        results=[];errors=[]
        def child():
            try:results.append(self.call("import os;os.write(1,b'own');os.write(2,b'err')"))
            except BaseException as e:errors.append(type(e).__name__)
        threads=[threading.Thread(target=child) for _ in range(16)]
        for thread in threads:thread.start()
        for thread in threads:thread.join(timeout=8)
        self.assertFalse(errors);self.assertTrue(all(not t.is_alive() for t in threads));self.assertEqual(16,len(results))
        for result in results:self.assert_complete(result);self.assertEqual(b'own',result.stdout);self.assertEqual(b'err',result.stderr)
    def test_bad_limits_fail_before_launch(self):
        for key,values in [('timeout',[True,0,float('inf'),float('nan'),7201]),('limit',[True,0,10485761]),('max_processes',[True,0,65]),('max_memory',[True,0,1073741825])]:
            for value in values:
                with self.subTest(key=key,value=value),patch.object(n,'OwnedCapture') as launch,self.assertRaises(ValueError):self.call('pass',**{key:value})
                launch.assert_not_called()
    def test_launch_failure_does_not_create_child_effect(self):
        result=n.capture([str(self.root/'nonexistent.exe')],self.root)
        self.assertEqual('LAUNCH_FAILED',result.status);self.assertEqual(b'',result.stdout);self.assertEqual(b'',result.stderr)
    def test_cleanup_waits_after_accounting_empty_without_retermination(self):
        owner=n.OwnedCapture([sys.executable,'-I','-c','import time;time.sleep(20)'],self.root,n.clean_environment());owner.start()
        actual_wait=win.wait
        def delayed(handle,milliseconds):
            # A nonblocking wait could see the transitional still-active state;
            # bounded wait must wait for the already-requested kernel cleanup.
            if handle==owner.pi.process and milliseconds==0:return 258
            return actual_wait(handle,milliseconds)
        try:
            with patch.object(win,'wait',side_effect=delayed),patch.object(win,'terminate',wraps=win.terminate) as retry:
                self.assertTrue(owner.close())
                retry.assert_not_called()
            self.assertTrue(owner.empty)
        finally:owner.close()

    def test_real_refund_binary_retains_exit_and_fixed_code(self):
        exe=Path(sys.argv[2]).resolve(strict=True)
        expected=sys.argv[3]
        self.assertEqual(expected,hashlib.sha256(exe.read_bytes()).hexdigest())
        # Explicit hashed local fixture; no DB/provider environment is inherited.
        result=n.capture([str(exe)],self.root)
        self.assert_complete(result,1);self.assertIn(b'REFUND_HOST_FAILED',result.stderr)
    def test_job_commit_budget_refuses_large_allocation(self):
        budget=64*1024*1024
        code="try:\n value=bytearray(128*1024*1024)\nexcept MemoryError:\n print('MEMORY_DENIED')\nelse:\n raise SystemExit(17)"
        result=self.call(code,max_memory=budget)
        self.assert_complete(result);self.assertIn(b'MEMORY_DENIED',result.stdout)
        self.assertIsNotNone(result.observed_peak_job_bytes);self.assertGreater(result.observed_peak_job_bytes,0)

    def test_native_commit_denied_and_granted_by_selected_budget(self):
        code="""import ctypes as c,json
from ctypes import wintypes as w
k=c.WinDLL('kernel32',use_last_error=True)
a=k.VirtualAlloc;a.argtypes=[c.c_void_p,c.c_size_t,w.DWORD,w.DWORD];a.restype=c.c_void_p
f=k.VirtualFree;f.argtypes=[c.c_void_p,c.c_size_t,w.DWORD];f.restype=w.BOOL
large=a(None,128*1024*1024,0x3000,4);error=c.get_last_error()
if large:assert f(large,0,0x8000)
small=a(None,1024*1024,0x3000,4);assert small;assert f(small,0,0x8000)
print(json.dumps({'large_granted':bool(large),'small_granted':bool(small),'large_error':error if not large else None}))
"""
        for mib,allowed in [(64,False),(256,True)]:
            with self.subTest(budget_mib=mib):
                result=self.call(code,max_memory=mib*1024*1024);self.assert_complete(result)
                self.assertEqual({'large_granted':allowed,'small_granted':True,'large_error':None if allowed else 1455},json.loads(result.stdout))
                self.assertGreater(result.observed_peak_job_bytes,0)

    def test_failed_pipe_setup_releases_created_handles(self):
        current=win.api('GetCurrentProcess',[],W.HANDLE)
        count=win.api('GetProcessHandleCount',[W.HANDLE,C.POINTER(W.DWORD)],W.BOOL)
        def observed():
            value=W.DWORD();self.assertTrue(count(current(),C.byref(value)));return value.value
        before=observed();pipe=n.create_pipe;calls=0
        def fail_second(*args):
            nonlocal calls
            calls+=1
            return 0 if calls==2 else pipe(*args)
        with patch.object(n,'create_pipe',side_effect=fail_second),patch.object(win,'create') as launch:
            result=self.call('pass')
        self.assertEqual('LAUNCH_FAILED',result.status);launch.assert_not_called();self.assertEqual(before,observed())

    def test_cleanup_failure_is_not_a_success_receipt(self):
        with patch.object(n,'terminate_job',return_value=0):
            result=self.call('import time;time.sleep(20)',timeout=.1)
        # Closing the owned job still requests termination, but failed explicit
        # cleanup must stay visible. No false proof of an empty tree is emitted.
        self.assertEqual('CLEANUP_FAILED',result.status);self.assertFalse(result.tree_empty)

    def test_active_process_budget_refuses_descendant_creation(self):
        child="from pathlib import Path;Path('quota-escape').write_text('bad')"
        code="import subprocess,sys,json\ntry:\n subprocess.Popen([sys.executable,'-I','-c',"+repr(child)+"])\nexcept OSError as error:\n print(json.dumps({'denied':error.winerror}))\nelse:\n raise SystemExit(17)"
        result=self.call(code,max_processes=1)
        self.assert_complete(result);self.assertEqual({'denied':1816},json.loads(result.stdout));self.assertFalse((self.root/'quota-escape').exists())

    def test_abrupt_owner_death_with_redirected_pipes_kills_tree(self):
        module_dir=Path(__file__).resolve().parent
        script=self.root/'owner.py'
        nested="import os,time,pathlib;pathlib.Path('nested.json').write_text(str(os.getpid()));time.sleep(20)"
        child=f"import os,subprocess,sys,time,pathlib;subprocess.Popen([sys.executable,'-I','-c',{nested!r}],creationflags=0x08000000)\nuntil=time.monotonic()+5\nwhile not pathlib.Path('nested.json').exists() and time.monotonic()<until:time.sleep(.005)\nos.write(1,b'x'*65536);time.sleep(20)"
        script.write_text("import sys,os,time,json\nfrom pathlib import Path\nsys.path.insert(0,sys.argv[1]);import native_capture as n\nroot=Path(sys.argv[2])\nowner=n.OwnedCapture([sys.executable,'-I','-c',"+repr(child)+"],root,n.clean_environment());owner.start()\nuntil=time.monotonic()+5\nwhile not (root/'nested.json').is_file() and time.monotonic()<until:time.sleep(.005)\nn.win.write(root/'owner-ready.json',{'root':owner.pi.pid,'nested':int((root/'nested.json').read_text())})\nn.win.await_file(root/'exit.json');os._exit(0)\n",encoding='utf-8')
        wrapper=subprocess.Popen([sys.executable,'-I',str(script),str(module_dir),str(self.root)],stdin=subprocess.DEVNULL,stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL,creationflags=subprocess.CREATE_NO_WINDOW,env=n.clean_environment())
        handles=[]
        try:
            info=win.await_file(self.root/'owner-ready.json')
            for key in ('root','nested'):
                handle=win.open_process(0x100000|0x1000,False,info[key]);self.assertTrue(handle);handles.append(handle);self.assertEqual(258,win.wait(handle,0))
            win.write(self.root/'exit.json',{'exit':True});wrapper.wait(timeout=5)
            self.assertEqual(0,wrapper.returncode)
            for handle in handles:self.assertEqual(0,win.wait(handle,5000))
        finally:
            if wrapper.poll() is None:wrapper.kill();wrapper.wait(timeout=5)
            for handle in handles:win.close(handle)

    def test_owned_handles_stable_after_repeated_runs(self):
        current=win.api('GetCurrentProcess',[],W.HANDLE)
        count=win.api('GetProcessHandleCount',[W.HANDLE,C.POINTER(W.DWORD)],W.BOOL)
        def observed():
            v=W.DWORD();self.assertTrue(count(current(),C.byref(v)));return v.value
        self.assert_complete(self.call('pass'));before=observed()
        for _ in range(24):self.assert_complete(self.call('pass'))
        self.assertEqual(before,observed())

def msvcrt_handle(fd):
    import msvcrt
    return msvcrt.get_osfhandle(fd)

if __name__=='__main__':
    outcomes=[]
    class Result(unittest.TextTestResult):
        def addSubTest(self,test,subtest,error):
            outcomes.append({'id':subtest.id(),'status':'PASS' if error is None else 'FAIL','subtest':True});super().addSubTest(test,subtest,error)
        def addSuccess(self,test):outcomes.append({'id':test.id(),'status':'PASS'});super().addSuccess(test)
        def addError(self,test,error):outcomes.append({'id':test.id(),'status':'ERROR'});super().addError(test,error)
        def addFailure(self,test,error):outcomes.append({'id':test.id(),'status':'FAIL'});super().addFailure(test,error)
    suite=unittest.defaultTestLoader.loadTestsFromTestCase(CaptureTests)
    result=unittest.TextTestRunner(verbosity=2,resultclass=Result).run(suite)
    Path(sys.argv[1]).write_text(json.dumps({'tests_run':result.testsRun,'pass':result.wasSuccessful(),'skips':len(result.skipped),'outcomes':outcomes},indent=2)+'\n')
    sys.exit(not result.wasSuccessful())
