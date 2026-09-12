# V346 — native supervisor I/O and resource qualification

Library maintenance / T2809, continued from validated checkpoint134.
AUTHORED qualification code, CONDITIONED / NOT_ADMITTED as product. The existing
SECURE-OPS1.1.3 product owner and all selected profile payloads remain unchanged.
This advances the missing supervisor candidate; it does not install a Windows
service, change its identity, restart real work or configure a telemetry backend.

## Concrete delta

V345 proved creation-bound JOB_LIST membership but did not redirect operational
I/O. V346 composes that exact API apparatus with CREATE_NO_WINDOW and a second
process-creation attribute, HANDLE_LIST: only NUL stdin and the two pipe write
handles are inherited. Read ends and the job handle stay parent-only. Parent
copies of inherited handles close before ResumeThread. There is no fallback to
uncontained creation or broad inheritance if setup fails.

The owning instance captures two binary streams independently and fairly with
nonblocking reads up to64KiB per stream/pass, a10MiB maximum retained per stream,
a caller deadline and cancellation. Completion requires the root exit status,
both EOFs and zero job-active processes. A live silent descendant cannot turn
root exit0 into success. Termination addresses the owned job and exact root
handles; it never enumerates/kills unrelated process IDs or replays work.

Fixture-only default budgets are8active processes and256MiB job commit, with
explicit test overrides. These are not production sizing decisions. The limits
are set and read back before creation. Cleanup has a shared five-second budget;
its failure remains CLEANUP_FAILED with no empty-tree proof. Capture statuses are
COMPLETED, TIMED_OUT, OUTPUT_LIMIT, CANCELLED, IO_FAILED, LAUNCH_FAILED or
CLEANUP_FAILED. COMPLETED preserves nonzero and259application exit values; it is
not semantic operation success. Raw output is excluded from Outcome repr but
remains available to its caller and is not automatically redacted or stored.

## Executed candidate evidence

- 23ordinary tests PASS, zero skips. Retained subcases counted separately
  in JSON; no test/fixture count is converted to a roadmap completion percentage.
- Separate stdout/stderr each1MiB at the exact cap, byte-for-byte/hash semantics;
  each stream limit+1 rejects before a later marker. Late descendant output after
  root exit is captured. Silent/open-pipe descendants time out and the tree empties.
- Pre-cancel avoids launch; active cancel empties the job. Injected read error
  produces a static failure and cleanup. Failed second pipe setup releases prior
  handles. Explicit cleanup failure never becomes a successful receipt.
- An intentionally inheritable sentinel file handle is excluded from the child;
  its numeric handle is checked by resolved file identity to avoid treating a
  coincidentally reused numeric value as a leaked handle. Sixteen concurrent
  launches retain separate streams and EOF. Twenty-four repeated runs preserve
  the observed parent process handle count after warm-up.
- A pipe-owning supervisor abruptly exits while root and descendant are alive;
  both captured owned process handles signal termination. The fixture establishes
  descendant readiness before saturating a shared output pipe. This verifies
  owner death with actual redirected I/O, not a cleanup callback mock.
- Active-process budget1refuses a descendant with Windows error1816 and no child
  marker. Python128MiB allocation under64MiB raises MemoryError. A native paired
  control denies128MiB VirtualAlloc under64MiB (error1455), grants the same request
  under256MiB and permits1MiB under both. Successful allocations are released.
- The previously compiled local refund0.1.4 binary is supplied explicitly by path
  and verified SHA-256; without database/provider environment it exits1 with the
  fixed REFUND_HOST_FAILED code. No provider, database or financial effect.
- The unchanged V345 native lifecycle apparatus reruns its five scenarios x3
  (15observations). Current files still require canonical rebuild below.

## Failures retained and corrected

FAIL595: job accounting reached zero before the root process handle was signaled.
The first cleanup tried TerminateProcess again, received error5 and falsely
reported CLEANUP_FAILED. Native call trace preserves active1→0, wait258,
TerminateProcess0/error5, then wait0. The correction waits for already-requested
termination within the cleanup budget. If the deadline is actually missed it
keeps failure and only makes a best-effort exact-root termination. Real failing
cases and an injected transitional-wait regression both pass after correction.

FAIL596: the first test incorrectly assumed PeakJobMemoryUsed was an enforced
committed-memory ceiling. Observed counters exceeded64MiB even when the128MiB
allocation failed. Native paired controls isolate budget-dependent allocation
refusal while retaining the counter (145391616and145125376 in the recorded probe).
The output is renamed observed_peak_job_bytes. No counter is clamped; the failed
assertion/source is preserved. Exact internal counter accounting is not claimed
resolved; this host observation is not a guarantee of RSS or total system memory.

FAIL597: the owner-death fixture flooded a shared pipe before the descendant
ready handshake. Diagnostic owner stderr was empty; a direct observation showed
three job-active processes, no ready file and no stderr. The revised fixture
starts its descendant hidden and waits for readiness before flooding output.
The original failure and source remain; no assertion is removed. This is a
fixture-order correction, not a claimed explanation of undocumented kernel
startup internals. Actual root/descendant handles now signal after owner death.

## Authority and scope of admission

Observed 2026-09-09T01:10:17.250920+00:00.
[Microsoft creation attributes](https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-updateprocthreadattribute)
require an explicit inheritable HANDLE_LIST and keep attribute values alive through
creation. JOB_LIST governs creation-bound containment from V345.
[Redirected I/O documentation](https://learn.microsoft.com/en-us/windows/win32/procthread/creating-a-child-process-with-redirected-input-and-output)
explains pipe-end ownership; its sample was consulted for semantics, not copied or
admitted. [CPython msvcrt](https://docs.python.org/3.14/library/msvcrt.html)
defines OS-handle to CRT-descriptor ownership. Execution remains CPython3.14.4;
current documentation displays3.14.7 and does not update the pinned runtime.

[Job limits](https://learn.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-jobobject_basic_limit_information),
[extended limits](https://learn.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-jobobject_extended_limit_information),
[job accounting](https://learn.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-jobobject_basic_accounting_information)
and [termination](https://learn.microsoft.com/en-us/windows/win32/api/jobapi2/nf-jobapi2-terminatejobobject)
govern the limits/observations above. [VirtualAlloc](https://learn.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc)
distinguishes committing virtual memory from physically touching pages. The
allocator controls prove the selected local rejection; they do not certify
production sizing or complete resource exhaustion resistance.

Internal authorities: SECURITY_SRE_CLOUD_INFRASTRUCTURE.md §12.5 and lifecycle/
least privilege boundaries; ENGINEERING_EXECUTION_PLAYBOOK.md; PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md.
Provenance and local workspace license are preserved. No Microsoft/Python source
is redistributed by this authored adapter; no runtime/dependency is added.

For the full service-supervisor capability, G0–G2 remain supported; G3–G8 remain
conditioned by real service identity, artifact authentication, installation and
lifecycle policy, cooperative shutdown/drain, restart/heartbeat budget, durable
receipt/collector/retention/alerts, security review and target acceptance. A job
is not a filesystem/network/privilege sandbox. Same-user hostile processes,
handle duplication attacks, arbitrary untrusted code and alternate OS/ABI are
outside this qualified claim. Do not incorporate the candidate in product until
the capability-gap gate admits the compatible implementation and target gates.

Next concrete work: govern executable identity and launch profile, then qualify
supervisor lifecycle/shutdown and operational result delivery using synthetic
work, while preserving the new handle/resource and prior launch-window tests.
No user account, service registration, restart of real jobs or spending is needed
for that reference qualification. Retention/alert destination/policy remain target
conditions. No broad absence-of-source or complete-SCA claim is made.

## Reproduction and assurance

The base job_probe.py is reconstructed from the Current qualification source
block in OPERATIONAL_PROCESS_BOUNDARY_V345.md; expected SHA-256:
`e33dbd14b0cd09262897c98263de0d572ac0b66673df6bce2f62c3bf9ade316c`. It is reused, not edited or copied into a new
product owner. Reconstruct the two following current qualification blocks into
a fresh external directory alongside it. Test arguments are output JSON path,
explicit local refund fixture executable and its expected SHA-256; all children
are synthetic. No generated files belong in the portable library tree.

Correctness: binary streams, exit codes, negative controls and real process trees.
Integration: job creation+handles+capture+cleanup+local refund executable.
Security: explicit handle/environment boundaries and static error codes, with
privilege sandbox/identity/source-admission limitations retained. Reliability:
EOF/tree completion, cancellation/deadline and parent-death failure injection.
Performance: finite stream/process/commit budgets and16concurrent captures;
no throughput/SLO/RSS claim. Operability: explicit outcomes and cleanup evidence;
service/collector/alerts remain incomplete. Maintainability: exact reused base,
two authored files, no dependency. Compatibility: Windows x64/build26200 and
CPython3.14.4 only. Delivery: canonical reconstruction plus test receipts, no
release/product profile promotion. These are narrow implementation-assurance
observations, not completion of the project assurance contract.

Stage: C:/Users/NL/AppData/Local/Temp/elite-v346-6b17f3aa330f426aadc7f8fe46a1bab9

| Evidence | SHA-256 |
|---|---|
| baseline.json | `decfb3f367774e074fb23a97828d75addb02606561907a8b745d7499ed43ccb4` |
| candidate-tests.log | `19eb43b443d622d02201a888350f50412bdaef9e2242a0d156aa8d8ad5f5d956` |
| candidate-tests.json | `8b3cd6eb1ac12a9e3c1f680e81d77244aec54a6895b861b73c3a67498ab5d0f3` |
| native_capture-before595.py | `6d79db63038f10bb129baec9ff9768ca699ed76dabea71d9600f375096f6f447` |
| cleanup-diagnosis.json | `8dd22458f38137ce2e8b2d0d979b97fc6030810e7b7de1dab1f200d9ba6f9d5d` |
| green-tests.json | `c996b56520a34d6cd1a3bb93d06cb8da6d05731debfc7993d6a6307a064d380a` |
| resource-tests.log | `fefd2eef308298a76bf0e3daf57eaff066b4602525c77a75484443a36c2ed7f1` |
| resource-tests.json | `264c960dbae38daac7e482788fd8a01b1d10b90ce6dd7b0515522c49a1d16777` |
| test_capture-before596.py | `46ef1833549d35a6ab087f6b0e6fd66d3525311f43ff3b1f50b6e65f039885ec` |
| owner-diagnosis-stderr.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| owner-diagnosis-error.log | `11064d45e889e306c17c87a5be9b29a353fd85f1fde866f451de1c5e24edb664` |
| owner-diagnosis-source.py | `c403e2705cf1856dc5982c3986939ad93f28763a66dad96b2756f3c6fd6028ea` |
| owner-child-diagnosis.json | `d95515d1464d0741d0fd85e1a74d262f8447d80cbbf201534d964a4eaf64d21d` |
| memory-native-receipt.json | `abe02ab8b891f9cbb218540bf279b3995d1d92696271adfa8d26b2572260b4fe` |
| final-candidate-tests.json | `a70f75e6e7a5b319943d848d6a7da7493f302a5f3e24120454d602899c5bb68d` |
| portable-candidate-tests.json | `c4b869ed05b0f2082f785da240aab9346736b414b47724192311db00336fa798` |
| refund-fixture.json | `7ca8f783cad926cae2135920208b07b56bc407f91600d2c7950e11f2385e4ce4` |
| candidate/job_probe.receipt.json | `964c49b0794a1b24d24d8e9e06cf452a0aeaa92d5dd3162ab5a7fbc8b8a20f60` |

## Qualification source: native_capture.py (AUTHORED / NOT_ADMITTED)

````python
"""AUTHORED Windows job + bounded I/O qualification; NOT_ADMITTED product.
No service registration, restart, secrets, telemetry export or privilege sandbox.
"""
from dataclasses import dataclass,field
from pathlib import Path
import ctypes as C
from ctypes import wintypes as W
import os,sys,time,subprocess,msvcrt,math,threading
import job_probe as win

class Security(C.Structure):
    _fields_=[('length',W.DWORD),('descriptor',C.c_void_p),('inherit',W.BOOL)]
class Accounting(C.Structure):
    _fields_=[('user',C.c_longlong),('kernel',C.c_longlong),('period_user',C.c_longlong),('period_kernel',C.c_longlong),('faults',W.DWORD),('total',W.DWORD),('active',W.DWORD),('terminated',W.DWORD)]
create_pipe=win.api('CreatePipe',[C.POINTER(W.HANDLE),C.POINTER(W.HANDLE),C.POINTER(Security),W.DWORD],W.BOOL)
set_handle=win.api('SetHandleInformation',[W.HANDLE,W.DWORD,W.DWORD],W.BOOL)
query=win.api('QueryInformationJobObject',[W.HANDLE,C.c_int,C.c_void_p,W.DWORD,C.c_void_p],W.BOOL)
terminate_job=win.api('TerminateJobObject',[W.HANDLE,W.UINT],W.BOOL)
exit_code=win.api('GetExitCodeProcess',[W.HANDLE,C.POINTER(W.DWORD)],W.BOOL)

class BoundaryError(Exception):
    pass

def check(ok):
    if not ok:raise BoundaryError('NATIVE_OPERATION_FAILED') from None

@dataclass(frozen=True)
class Outcome:
    status: str
    exit_code: int|None
    stdout: bytes=field(repr=False)
    stderr: bytes=field(repr=False)
    tree_empty: bool
    # Windows peak accounting can include denied commit attempts on this host.
    # This observation is neither an RSS ceiling nor proof of granted allocation.
    observed_peak_job_bytes: int|None=None

class OwnedCapture:
    """Own one fresh job, root process, thread, pipes and CRT descriptors.

    Call close on every path, including capture failure. Read ends are raw,
    nonblocking and parent-only; HANDLE_LIST grants only the three std handles.
    This reference assumes exclusive local use of each instance.
    """
    def __init__(self,argv,cwd,env,*,max_processes=8,max_memory=256*1024*1024):
        self.job=None;self.pi=win.Process();self.read_fds=[];self.write_handles=[]
        self.nul_fd=None;self.raw_handles=[];self.closed=False;self.empty=False;self.peak_commit=None
        try:
            self.job=win.job()
            limits=win.Extended();limits.basic.flags=0x2000|0x8|0x200
            limits.basic.active=max_processes;limits.job_memory=max_memory
            check(win.set_job(self.job,9,C.byref(limits),C.sizeof(limits)))
            observed=win.Extended();check(query(self.job,9,C.byref(observed),C.sizeof(observed),None))
            check(observed.basic.active==max_processes and observed.job_memory==max_memory and observed.basic.flags & 0x2208==0x2208)
            sa=Security(C.sizeof(Security),None,True)
            for _ in range(2):
                read=W.HANDLE();write=W.HANDLE()
                check(create_pipe(C.byref(read),C.byref(write),C.byref(sa),0))
                self.raw_handles.extend([read.value,write.value]);self.write_handles.append(write.value)
                check(set_handle(read,1,0))
                fd=msvcrt.open_osfhandle(read.value,os.O_RDONLY|os.O_BINARY)
                self.raw_handles.remove(read.value);self.read_fds.append(fd)
                os.set_blocking(fd,False)
            self.nul_fd=os.open(os.devnull,os.O_RDONLY|os.O_BINARY)
            os.set_inheritable(self.nul_fd,True)
            self._create(argv,cwd,env)
            # Parent copies must close before resume, otherwise EOF cannot prove
            # the child side closed. Already-inherited child handles remain valid.
            self._close_writes()
        except BaseException:
            self.close()
            raise

    def _create(self,argv,cwd,env):
        size=win.SIZE()
        check(not win.initialize(None,2,0,C.byref(size)) and C.get_last_error()==122 and size.value>0)
        storage=C.create_string_buffer(size.value)
        check(win.initialize(storage,2,0,C.byref(size)))
        jobs=(W.HANDLE*1)(self.job)
        permitted=(W.HANDLE*3)(msvcrt.get_osfhandle(self.nul_fd),*self.write_handles)
        si=win.StartupEx();si.startup.cb=C.sizeof(si);si.startup.flags=0x100
        si.startup.stdin,si.startup.stdout,si.startup.stderr=permitted
        try:
            check(win.update(storage,0,0x2000d,jobs,C.sizeof(jobs),None,None))
            check(win.update(storage,0,0x20002,permitted,C.sizeof(permitted),None,None))
            si.attributes=C.cast(storage,C.c_void_p)
            command=C.create_unicode_buffer(subprocess.list2cmdline(argv))
            block=C.create_unicode_buffer('\0'.join(k+'='+v for k,v in sorted(env.items()))+'\0\0')
            check(win.create(argv[0],command,None,None,True,0x4|0x08000000|0x400|0x80000,block,str(cwd),C.cast(C.byref(si),C.POINTER(win.Startup)),C.byref(self.pi)))
        finally:
            win.delete(storage)

    def _close_writes(self):
        for h in list(self.write_handles):
            check(win.close(h));self.raw_handles.remove(h);self.write_handles.remove(h)
        if self.nul_fd is not None:
            os.close(self.nul_fd);self.nul_fd=None

    def start(self):
        check(win.resume(self.pi.thread)!=0xffffffff)
        check(win.close(self.pi.thread));self.pi.thread=None

    def active(self):
        accounting=Accounting()
        check(query(self.job,1,C.byref(accounting),C.sizeof(accounting),None))
        return accounting.active

    def code(self):
        # Exit code259 is a valid application exit; distinguish via process wait.
        status=win.wait(self.pi.process,0)
        if status==258:return None
        check(status==0);value=W.DWORD();check(exit_code(self.pi.process,C.byref(value)))
        return value.value

    def close(self):
        if self.closed:return self.empty
        ok=True;cleanup_deadline=time.monotonic()+5
        # No arbitrary PID lookup; operate solely on handles created here.
        if self.job:
            try:
                limits=win.Extended();check(query(self.job,9,C.byref(limits),C.sizeof(limits),None))
                self.peak_commit=limits.peak_job
                if self.active():check(terminate_job(self.job,99))
                while self.active() and time.monotonic()<cleanup_deadline:time.sleep(.005)
                self.empty=self.active()==0
                ok=self.empty and ok
            except (BoundaryError,OSError):ok=False
            ok=bool(win.close(self.job)) and ok;self.job=None
        # Job accounting may become empty before the root handle is signaled.
        # Wait for termination already requested by TerminateJob/CloseJob;
        # do not re-terminate a process in that transition (AccessDenied race).
        if self.pi.process:
            milliseconds=max(0,int((cleanup_deadline-time.monotonic())*1000))
            reaped=win.wait(self.pi.process,milliseconds)==0
            if not reaped:
                # Best effort containment only; a missed cleanup deadline can
                # never become success, even if termination is then accepted.
                win.terminate(self.pi.process,99)
                ok=False
            ok=bool(win.close(self.pi.process)) and ok;self.pi.process=None
        if self.pi.thread:ok=bool(win.close(self.pi.thread)) and ok;self.pi.thread=None
        for fd in self.read_fds:
            try:os.close(fd)
            except OSError:ok=False
        self.read_fds.clear()
        if self.nul_fd is not None:
            try:os.close(self.nul_fd)
            except OSError:ok=False
            self.nul_fd=None
        for h in self.raw_handles:ok=bool(win.close(h)) and ok
        self.raw_handles.clear();self.write_handles.clear();self.closed=True
        self.empty=self.empty and ok
        return ok

def clean_environment():
    return {k:v for k,v in os.environ.items() if k.upper() in {'SYSTEMROOT','WINDIR','TEMP','TMP'}}

def capture(argv,cwd,*,timeout=5.0,limit=10*1024*1024,cancel=None,max_processes=8,max_memory=256*1024*1024):
    # Qualification API accepts only explicit absolute executables; no PATH/shell.
    if not isinstance(argv,list) or not argv or len(argv)>129 or any(not isinstance(x,str) or '\0' in x or len(x)>4096 for x in argv):
        raise ValueError('INVALID_ARGUMENTS') from None
    if not Path(argv[0]).is_absolute() or not Path(cwd).is_absolute():raise ValueError('ABSOLUTE_PATH_REQUIRED') from None
    if type(timeout) not in (int,float) or not math.isfinite(timeout) or not 0<timeout<=7200:raise ValueError('INVALID_TIMEOUT') from None
    if type(limit) is not int or not 1<=limit<=10*1024*1024:raise ValueError('INVALID_OUTPUT_LIMIT') from None
    if type(max_processes) is not int or not 1<=max_processes<=64:raise ValueError('INVALID_PROCESS_LIMIT') from None
    if type(max_memory) is not int or not 32*1024*1024<=max_memory<=1024*1024*1024:raise ValueError('INVALID_MEMORY_LIMIT') from None
    if cancel is not None and not isinstance(cancel,threading.Event):raise ValueError('INVALID_CANCELLATION') from None
    if cancel is not None and cancel.is_set():return Outcome('CANCELLED',None,b'',b'',True)
    deadline=time.monotonic()+timeout;owner=None;buffers=[bytearray(),bytearray()];status='LAUNCH_FAILED';code=None;empty=False
    try:
        owner=OwnedCapture(argv,Path(cwd),clean_environment(),max_processes=max_processes,max_memory=max_memory)
        if cancel is not None and cancel.is_set():status='CANCELLED'
        elif time.monotonic()>=deadline:status='TIMED_OUT'
        else:
            owner.start();pending={0,1};status='IO_FAILED'
            while True:
                if cancel is not None and cancel.is_set():status='CANCELLED';break
                if time.monotonic()>=deadline:status='TIMED_OUT';break
                progress=False
                for index in tuple(pending):
                    try:chunk=os.read(owner.read_fds[index],min(65536,limit-len(buffers[index])+1))
                    except BlockingIOError:continue
                    if not chunk:pending.remove(index);continue
                    progress=True
                    if len(buffers[index])+len(chunk)>limit:status='OUTPUT_LIMIT';break
                    buffers[index].extend(chunk)
                if status=='OUTPUT_LIMIT':break
                code=owner.code()
                if code is not None and not pending and owner.active()==0:status='COMPLETED';break
                if not progress:time.sleep(min(.005,max(0,deadline-time.monotonic())))
    except (OSError,ValueError,BoundaryError):
        # Exceptions may carry command/path/private data; return a fixed code.
        pass
    finally:
        if owner is not None:
            if not owner.close():status='CLEANUP_FAILED'
            empty=owner.empty
    return Outcome(status,code,bytes(buffers[0]),bytes(buffers[1]),empty,owner.peak_commit if owner is not None else None)
````

## Qualification source: test_capture.py (AUTHORED / NOT_ADMITTED)

````python
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
        result=self.call(f"import subprocess,sys,time;subprocess.Popen([sys.executable,'-I','-c',{nested!r}]);time.sleep(20)",timeout=.6)
        self.assertTrue((self.root/'nested-ready').exists());self.assertEqual('TIMED_OUT',result.status);self.assertTrue(result.tree_empty)
    def test_silent_descendant_cannot_be_success_when_root_exits(self):
        nested="import time,pathlib;pathlib.Path('silent-ready').write_text('ready');time.sleep(20)"
        result=self.call(f"import subprocess,sys;subprocess.Popen([sys.executable,'-I','-c',{nested!r}],stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL)",timeout=.6)
        self.assertTrue((self.root/'silent-ready').exists());self.assertEqual('TIMED_OUT',result.status);self.assertTrue(result.tree_empty)
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
````

## Canonical rebuild and integration135

Fresh reconstruction of both current source blocks and the unchanged V345 base
is byte-identical. Rebuilt suite:23ordinary tests PASS, zero skips;
base lifecycle:15observations PASS. Current candidate and rebuilt
runs are distinct; no assertion/case dropped after failure diagnosis. All original
production owner/profile bytes stay unchanged. Capability-gap kit0.1.0 composed
again and self-check3positives/6negatives PASS. Preflight135 is pending at this
integration checkpoint, not passed by inference from V345.

- rebuilt-tests.json: `c4b869ed05b0f2082f785da240aab9346736b414b47724192311db00336fa798`
- rebuilt/job_probe.receipt.json: `cce8531a086bf6a5d972619ae4516aff0f94a8fa94c3bfbf5fd6f741bfa804b3`
- canonical-parity.json: `a7f363f78dc4f512a67a2aa697f99a4c96381950e598bd111dabbffaf24c78ab`

## Closure 136 — V346 qualification complete, service admission still open

Preflight 135 completed with 154 executed checks PASS. Its structural gate
verified 162 packs, 1461 materializable files, 784 Markdown files
and 53 profiles. Docker is the only unavailable toolchain; this availability
blocker and all network/provider/target opt-in skips remain explicit. A passing
local library Preflight is not a production or full service-supervisor admission.

The earlier pending statements describe the integration-135 snapshot. This
closure supersedes them: the two current qualification files and unchanged V345
base reconstructed exactly; 23 ordinary tests with zero skips and 15 prior
lifecycle observations passed. Candidate/rebuilt SHA-256 equality is preserved.
No qualification source, product pack or profile changed during Preflight.

FAIL595 is regression-proven for termination accounting/signaling; FAIL596 is
recovered with the peak-counter limitation and paired native allocation controls;
FAIL597 is recovered for the ordered owner-death fixture. Original failures and
sources remain retained. No broader security, memory, lifecycle or release claim
is inferred from those fixes.

Next bounded work: executable identity and governed launch profile in this local
candidate, then cooperative shutdown/lifecycle and durable result delivery. Keep
the current 23 tests and 15 lifecycle observations as regression evidence. Service
identity, target policy, compatibility/security/SCA and G0–G8 admission remain
required before product incorporation. Other roadmap fronts remain open.

The final successor capability decision is gap-record136.json with validation
receipt gap-receipt136.json in the V346 stage. The original integration-135 record
and receipt remain historical. Their structural PASS does not close service gates.

| Closure evidence | SHA-256 |
|---|---|
| preflight135.json | `31e13dd1a870799e56f74645c65de005cfc177999f74448179b51399aa230ecf` |
| preflight135.log | `58dc541dbd77a4b2f37dc8f67750e6b35c84d1532fcbe35efb557a82a6e82c89` |
| structural-closure135.json | `534649af8ffe9949a08b6565803ccec16f1cc5e0d7b1ed2acc7182e40bc45d93` |
| rebuilt-tests.json | `c4b869ed05b0f2082f785da240aab9346736b414b47724192311db00336fa798` |
| rebuilt/job_probe.receipt.json | `cce8531a086bf6a5d972619ae4516aff0f94a8fa94c3bfbf5fd6f741bfa804b3` |
| canonical-parity.json | `a7f363f78dc4f512a67a2aa697f99a4c96381950e598bd111dabbffaf24c78ab` |
| gap-record135.json | `524c67908ea3fa58b63462e173babeefc2241006bd3b7a689edf29d867d6431d` |
| gap-receipt135.json | `e0b53aac795ed339bbf5500ccceeb21639bd18a4b9995c9daa347d989322a783` |
