# V348 — bounded cooperative shutdown qualification

Library maintenance / T2809, continued from validated checkpoint138. AUTHORED
qualification apparatus only; CONDITIONED / NOT_ADMITTED as product. No product
pack/profile, runtime/dependency, user identity or installed service changed.

## Concrete delta

The V347 candidate gains an opt-in shutdown protocol. A strict v2 launch profile
binds protocol win32-inherited-event/v1 and a finite grace interval (0,60] seconds.
V1 profiles keep their previous immediate cancellation/timeout behavior. Invalid
protocols, fields and grace values reject before creation. A new unnamed manual-
reset event is initially nonsignaled; only its duplicate with SYNCHRONIZE access
joins the existing explicit HANDLE_LIST. The parent retains the original handle
and closes its inheritable duplicate before resuming the contained process.
ELITE_STOP_EVENT_HANDLE is injected internally, never supplied by profile input.

Cancellation or the execution deadline requests the event once. The first reason
is retained, capture keeps draining bounded stdout/stderr, and the additional
grace deadline is not renewed. Root exit, both EOFs and zero active job processes
are required for EXITED_DURING_GRACE. Otherwise the existing owned-job termination
and finite cleanup run. Output/read errors still stop immediately; SetEvent failure
returns STOP_SIGNAL_FAILED. Cleanup failure never upgrades the receipt to success.

Outcome.shutdown_state is NOT_REQUESTED, EXITED_DURING_GRACE, FORCED or
SIGNAL_FAILED (REQUESTED is internal only). Exit codes retain their original
meaning, including nonzero drained exits. EXITED_DURING_GRACE describes process
lifecycle only: a worker may finish naturally without observing the event. It is
not application acknowledgement, domain commit, business success or durable
supervisor result delivery. The synthetic drain fixture separately observes its
event, finishes exactly one already-started unit, fsyncs its local marker and
emits final output. That fixture is not a mounted business worker.

The execution timeout plus optional grace plus up to five seconds of existing
cleanup form separate budgets; hashing/setup/OS scheduling are not strict wall-
clock SLO guarantees. Abrupt supervisor death still kills its owned job; it cannot
offer graceful drain after the owner is gone. No automatic restart or retry.

## Failures, scope correction and preserved evidence

FAIL599 / LIB-FAIL-2315: initial 16 tests passed, but the additional native access
probe disproved a broad wait-only security claim. Direct SetEvent on the inherited
handle fails5; a same-user child can duplicate the event with MODIFY_STATE and
signal that duplicate. event-access-probe.json preserves the exact observations;
last-error values after successful calls are stale, not errors. The protocol
requires a trusted cooperative worker. It does not isolate hostile same-user
processes or replace service identity/DACL/privilege/sandbox admission. Original
code/tests and the initial COOPERATIVE label are retained; the current label is
the narrower EXITED_DURING_GRACE. The suite characterizes amplification explicitly
so a future agent cannot reinstate the rejected security claim from direct denial.

FAIL600 / LIB-FAIL-2316: the retained capture suite's silent-descendant timeout and
tree-empty assertions passed, but immediate TemporaryDirectory removal failed32.
Original capture-regression.log/json and test-capture-before600.py are preserved.
The successor test opens that specific descendant while alive, confirms membership
in the owned job before cleanup, requires its process handle to signal, then
checks owned-directory rename/release under a separate five-second budget. No
arbitrary PID is killed, no live descendant accepted, and no deadline is removed.
Job accounting/root completion is not claimed to imply immediate filesystem
quiescence for every descendant. Other process-lifecycle cases remain intact.

## Executed candidate qualification and remaining conditions

Current candidate:19 shutdown tests,20 governed-launch tests and23 capture/resource
tests PASS, zero skips. Includes cancellation/drain, deadline/grace, hung child,
nonzero exit, output overflow during drain, signal failure, event setup failure,
pre-cancel, strict v2 profile, legacy v1 behavior,16 concurrent event instances,
20 repeated lifetimes without handle growth, direct-handle denial and same-user
amplification, natural exit without acknowledgement, and abrupt owner death with
live root/descendant waiting while the event exists. Exact canonical rebuild is
required below before integration; these are local synthetic finite experiments.

Next: durable operational result delivery and mounting lifecycle into the actual
selected worker. Trusted profile selection/signature, runtime/DLL/SCA/compatibility,
service identity, restart policy, real target collection/retention/alerts/recovery
and G0–G8 admission remain separate. No business/provider/ARCA action, account,
service registration, global install or spending. Other roadmap fronts stay open.

## Authority and reproduction

Current Microsoft official documentation was consulted, not copied as source:
[CreateEventW](https://learn.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-createeventw)
defines unnamed/manual-reset event behavior;
[DuplicateHandle](https://learn.microsoft.com/en-us/windows/win32/api/handleapi/nf-handleapi-duplicatehandle)
allows selected access and warns some object types can gain greater access;
[synchronization rights](https://learn.microsoft.com/en-us/windows/win32/sync/synchronization-object-security-and-access-rights)
and [SetEvent](https://learn.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-setevent)
distinguish waiting from direct modification. No named IPC or external source
artifact is acquired. Workspace-owner license and AUTHORED provenance apply.
Internal authorities: SECURITY_SRE_CLOUD_INFRASTRUCTURE.md lifecycle/supply-chain
boundaries, ENGINEERING_EXECUTION_PLAYBOOK.md and public architecture admission.

Stage: C:/Users/NL/AppData/Local/Temp/elite-v348-93f1c51f1d86411fb3b5c5e94aaa54e8. Pinned runtime:CPython3.14.4, local Windows x64. Rebuild
the four current source blocks below, reuse job_probe.py from V345 and unchanged
test_governed_launch.py from V347. Execute test_shutdown.py <new-report.json>,
test_governed_launch.py <new-report.json> <refund-fixture.json>, and test_capture.py
<new-report.json> <refund-binary-path> <expected-hash>. The previous refund fixture
remains explicitly hash-checked and never connects to DB/provider in these tests.
No test totals are translated into a roadmap completion percentage.

## Qualification source: native_capture.py (AUTHORED / NOT_ADMITTED)

SHA-256: `7e19a3ce7ef16ed10c7bddaaee97ecaba65e7eed07ed16fdcb5a20be4f56e61f`

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
create_event=win.api('CreateEventW',[C.c_void_p,W.BOOL,W.BOOL,W.LPCWSTR],W.HANDLE)
set_event=win.api('SetEvent',[W.HANDLE],W.BOOL)
duplicate=win.api('DuplicateHandle',[W.HANDLE,W.HANDLE,W.HANDLE,C.POINTER(W.HANDLE),W.DWORD,W.BOOL,W.DWORD],W.BOOL)
current_process=win.api('GetCurrentProcess',[],W.HANDLE)

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
    # Lifecycle observation only; never an acknowledgement of domain commit/drain.
    shutdown_state: str='NOT_REQUESTED'

class OwnedCapture:
    """Own one fresh job, root process, thread, pipes and CRT descriptors.

    Call close on every path, including capture failure. Read ends are raw,
    nonblocking and parent-only; HANDLE_LIST grants only the three std handles.
    This reference assumes exclusive local use of each instance.
    """
    def __init__(self,argv,cwd,env,*,max_processes=8,max_memory=256*1024*1024,cooperative=False):
        self.job=None;self.pi=win.Process();self.read_fds=[];self.write_handles=[]
        self.nul_fd=None;self.raw_handles=[];self.closed=False;self.empty=False;self.peak_commit=None
        self.stop_event=None;self.child_stop=None
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
            if cooperative:
                # Only trusted workers: SYNCHRONIZE limits direct use of this
                # handle, but same-user duplication can amplify event access.
                self.stop_event=create_event(None,True,False,None);check(self.stop_event)
                child=W.HANDLE()
                check(duplicate(current_process(),self.stop_event,current_process(),C.byref(child),0x100000,True,0))
                self.child_stop=child.value
                env=dict(env);env['ELITE_STOP_EVENT_HANDLE']=str(self.child_stop)
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
        inherited=[msvcrt.get_osfhandle(self.nul_fd),*self.write_handles]
        if self.child_stop is not None:inherited.append(self.child_stop)
        permitted=(W.HANDLE*len(inherited))(*inherited)
        si=win.StartupEx();si.startup.cb=C.sizeof(si);si.startup.flags=0x100
        si.startup.stdin,si.startup.stdout,si.startup.stderr=inherited[:3]
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
        if self.child_stop is not None:
            check(win.close(self.child_stop));self.child_stop=None

    def request_stop(self):
        check(self.stop_event is not None and set_event(self.stop_event))

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
        for name in ['child_stop','stop_event']:
            h=getattr(self,name)
            if h is not None:ok=bool(win.close(h)) and ok;setattr(self,name,None)
        self.raw_handles.clear();self.write_handles.clear();self.closed=True
        self.empty=self.empty and ok
        return ok

def clean_environment():
    return {k:v for k,v in os.environ.items() if k.upper() in {'SYSTEMROOT','WINDIR','TEMP','TMP'}}

def capture(argv,cwd,*,timeout=5.0,limit=10*1024*1024,cancel=None,max_processes=8,max_memory=256*1024*1024,environment=None,shutdown_grace=None):
    # Qualification API accepts only explicit absolute executables; no PATH/shell.
    if not isinstance(argv,list) or not argv or len(argv)>129 or any(not isinstance(x,str) or '\0' in x or len(x)>4096 for x in argv):
        raise ValueError('INVALID_ARGUMENTS') from None
    if not Path(argv[0]).is_absolute() or not Path(cwd).is_absolute():raise ValueError('ABSOLUTE_PATH_REQUIRED') from None
    if type(timeout) not in (int,float) or not math.isfinite(timeout) or not 0<timeout<=7200:raise ValueError('INVALID_TIMEOUT') from None
    if type(limit) is not int or not 1<=limit<=10*1024*1024:raise ValueError('INVALID_OUTPUT_LIMIT') from None
    if type(max_processes) is not int or not 1<=max_processes<=64:raise ValueError('INVALID_PROCESS_LIMIT') from None
    if type(max_memory) is not int or not 32*1024*1024<=max_memory<=1024*1024*1024:raise ValueError('INVALID_MEMORY_LIMIT') from None
    if cancel is not None and not isinstance(cancel,threading.Event):raise ValueError('INVALID_CANCELLATION') from None
    if shutdown_grace is not None and (type(shutdown_grace) not in (int,float) or not math.isfinite(shutdown_grace) or not 0<shutdown_grace<=60):raise ValueError('INVALID_SHUTDOWN_GRACE') from None
    if environment is not None:
        if not isinstance(environment,dict) or set(environment)!={'SYSTEMROOT','WINDIR','TEMP','TMP'} or any(not isinstance(v,str) or not v or '\0' in v or len(v)>4096 for v in environment.values()):
            raise ValueError('INVALID_ENVIRONMENT') from None
        environment=dict(environment)
    if cancel is not None and cancel.is_set():return Outcome('CANCELLED',None,b'',b'',True)
    deadline=time.monotonic()+timeout;owner=None;buffers=[bytearray(),bytearray()];status='LAUNCH_FAILED';code=None;empty=False
    shutdown='NOT_REQUESTED';stop_reason=None;grace_deadline=None
    try:
        owner=OwnedCapture(argv,Path(cwd),clean_environment() if environment is None else environment,max_processes=max_processes,max_memory=max_memory,cooperative=shutdown_grace is not None)
        if cancel is not None and cancel.is_set():status='CANCELLED'
        elif time.monotonic()>=deadline:status='TIMED_OUT'
        else:
            owner.start();pending={0,1};status='IO_FAILED'
            while True:
                now=time.monotonic()
                if stop_reason is None:
                    reason='CANCELLED' if cancel is not None and cancel.is_set() else ('TIMED_OUT' if now>=deadline else None)
                    if reason is not None:
                        if shutdown_grace is None:status=reason;break
                        stop_reason=reason;grace_deadline=now+shutdown_grace;shutdown='REQUESTED'
                        try:owner.request_stop()
                        except BoundaryError:status='STOP_SIGNAL_FAILED';shutdown='SIGNAL_FAILED';break
                elif now>=grace_deadline:
                    status=stop_reason;shutdown='FORCED';break
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
                if code is not None and not pending and owner.active()==0:
                    status=stop_reason or 'COMPLETED'
                    if stop_reason is not None:shutdown='EXITED_DURING_GRACE'
                    break
                if not progress:time.sleep(min(.005,max(0,(grace_deadline if stop_reason is not None else deadline)-time.monotonic())))
    except (OSError,ValueError,BoundaryError):
        # Exceptions may carry command/path/private data; return a fixed code.
        pass
    finally:
        if shutdown=='REQUESTED':shutdown='FORCED'
        if owner is not None:
            if not owner.close():status='CLEANUP_FAILED'
            empty=owner.empty
    return Outcome(status,code,bytes(buffers[0]),bytes(buffers[1]),empty,owner.peak_commit if owner is not None else None,shutdown)
````

## Qualification source: governed_launch.py (AUTHORED / NOT_ADMITTED)

SHA-256: `a66a2936bcaea449c473744a8384d4140e713c1132008e8617eb10a7be44555a`

````python
"""AUTHORED local qualification only; not an admitted service supervisor.

An externally trusted profile digest binds executable, argv, environment and
budgets. Native handles pin the selected path and file while hashing/launching.
This is neither a signature authority nor DLL/script/runtime security admission.
"""
from pathlib import PureWindowsPath
from dataclasses import dataclass, field
import ctypes as C
from ctypes import wintypes as W
import hashlib, json, math, re, threading
import job_probe as win
import native_capture

create_file=win.api('CreateFileW',[W.LPCWSTR,W.DWORD,W.DWORD,C.c_void_p,W.DWORD,W.DWORD,W.HANDLE],W.HANDLE)
final_path=win.api('GetFinalPathNameByHandleW',[W.HANDLE,W.LPWSTR,W.DWORD,W.DWORD],W.DWORD)
file_info=win.api('GetFileInformationByHandle',[W.HANDLE,C.c_void_p],W.BOOL)
read_file=win.api('ReadFile',[W.HANDLE,C.c_void_p,W.DWORD,C.POINTER(W.DWORD),C.c_void_p],W.BOOL)
get_type=win.api('GetFileType',[W.HANDLE],W.DWORD)
class Info(C.Structure):
    _fields_=[('attributes',W.DWORD),('creation',W.FILETIME),('access',W.FILETIME),('write',W.FILETIME),('volume',W.DWORD),('size_high',W.DWORD),('size_low',W.DWORD),('links',W.DWORD),('index_high',W.DWORD),('index_low',W.DWORD)]
class Rejected(Exception):
    pass
def require(ok):
    if not ok: raise Rejected('LAUNCH_PROFILE_REJECTED') from None
HEX=re.compile(r'[0-9a-f]{64}')
RESERVED=re.compile(r'(?:CON|PRN|AUX|NUL|COM[0-9¹²³]|LPT[0-9¹²³])(?:\..*)?',re.I)

def local_path(value):
    require(isinstance(value,str) and 3<=len(value)<=240 and re.match(r'^[A-Za-z]:\\',value))
    require(not any(ord(c)<32 or c in '/<>"|?*' for c in value) and ':' not in value[2:])
    parts=value[3:].split('\\') if len(value)>3 else []
    require(all(p and p not in ('.','..') and not p.endswith((' ','.')) and not RESERVED.fullmatch(p) for p in parts))
    return PureWindowsPath(value)

def unique_pairs(pairs):
    result={}
    for k,v in pairs:
        require(k not in result);result[k]=v
    return result

def profile(raw, trusted_sha256):
    require(type(raw) is bytes and 1<=len(raw)<=65536 and isinstance(trusted_sha256,str) and HEX.fullmatch(trusted_sha256))
    require(hashlib.sha256(raw).hexdigest()==trusted_sha256)
    try:v=json.loads(raw.decode('utf-8'),object_pairs_hook=unique_pairs,parse_constant=lambda _:require(False))
    except (ValueError,UnicodeError,RecursionError):raise Rejected('LAUNCH_PROFILE_REJECTED') from None
    require(isinstance(v,dict))
    fields={'schema','id','executable','arguments','cwd','environment','budgets'}
    require(v.get('schema') in ('elite-native-launch-qualification/v1','elite-native-launch-qualification/v2'))
    if v['schema'].endswith('/v2'):
        fields.add('shutdown');stop=v.get('shutdown')
        require(isinstance(stop,dict) and set(stop)=={'protocol','grace_seconds'} and stop['protocol']=='win32-inherited-event/v1')
        require(type(stop['grace_seconds']) in (int,float) and math.isfinite(stop['grace_seconds']) and 0<stop['grace_seconds']<=60)
    require(set(v)==fields and isinstance(v['id'],str) and re.fullmatch(r'[a-z0-9][a-z0-9-]{0,63}',v['id']))
    e=v['executable'];require(isinstance(e,dict) and set(e)=={'path','sha256','bytes'})
    local_path(e['path']);require(e['path'].lower().endswith('.exe') and isinstance(e['sha256'],str) and HEX.fullmatch(e['sha256']))
    require(type(e['bytes']) is int and 1<=e['bytes']<=1024*1024*1024)
    a=v['arguments'];require(isinstance(a,list) and len(a)<=128 and all(isinstance(s,str) and '\0' not in s and len(s)<=4096 for s in a))
    require(sum(len(s)+3 for s in a)+len(e['path'])<30000)
    local_path(v['cwd']);env=v['environment']
    require(isinstance(env,dict) and set(env)=={'SYSTEMROOT','WINDIR','TEMP','TMP'})
    for p in env.values():local_path(p)
    b=v['budgets'];require(isinstance(b,dict) and set(b)=={'timeout','output_bytes','processes','commit_bytes'})
    require(type(b['timeout']) in (int,float) and math.isfinite(b['timeout']) and 0<b['timeout']<=7200)
    for key,lo,hi in [('output_bytes',1,10485760),('processes',1,64),('commit_bytes',33554432,1073741824)]:require(type(b[key]) is int and lo<=b[key]<=hi)
    return v

class PinnedImage:
    """Noninheritable, read-shared handles; close only after capture cleanup.

    Reject local path aliases/reparse components. Pin directories before opening
    descendants and verify each actual handle path. File hashing reads the held
    handle; read-only sharing refuses competing writers/deleters. Assumes a
    trusted local drive mapping/OS/admin and no hostile kernel/privilege changes.
    """
    def __init__(self, executable, cwd):
        self.handles=[];self.directories=set();self.identity=None
        try:
            exe=local_path(executable['path']);working=local_path(cwd)
            for directory in [*reversed(exe.parent.parents),exe.parent,*reversed(working.parents),working]:
                key=str(directory).casefold()
                if key not in self.directories:
                    self._open(directory,True);self.directories.add(key)
            handle,info=self._open(exe,False)
            require((info.size_high<<32|info.size_low)==executable['bytes'])
            digest=hashlib.sha256();remaining=executable['bytes'];buf=C.create_string_buffer(65536)
            while remaining:
                count=W.DWORD();require(read_file(handle,buf,min(remaining,len(buf)),C.byref(count),None) and 0<count.value<=remaining)
                digest.update(buf.raw[:count.value]);remaining-=count.value
            require(digest.hexdigest()==executable['sha256'])
            self.identity={'volume':info.volume,'file_index':info.index_high<<32|info.index_low,'sha256':digest.hexdigest(),'bytes':executable['bytes']}
        except BaseException:
            self.close();raise
    def _open(self,path,directory):
        h=create_file(str(path),0x80 if directory else 0x80000000,1,None,3,0x200000|(0x2000000 if directory else 0),None)
        require(h not in (None,C.c_void_p(-1).value))
        self.handles.append(h);info=Info();require(file_info(h,C.byref(info)))
        require(get_type(h)==1 and not info.attributes&0x400 and bool(info.attributes&0x10)==directory)
        buf=C.create_unicode_buffer(32768);n=final_path(h,buf,len(buf),0)
        require(0<n<len(buf) and buf.value.casefold()==('\\\\?\\'+str(path)).casefold())
        return h,info
    def close(self):
        ok=True
        for h in reversed(self.handles):ok=bool(win.close(h)) and ok
        self.handles.clear()
        return ok

@dataclass(frozen=True)
class LaunchResult:
    status:str
    profile_sha256:str
    image:dict|None
    outcome:native_capture.Outcome|None=field(repr=False)

def launch(raw, trusted_sha256, *, cancel=None):
    # The caller must obtain the digest from a trusted owner, not from raw itself.
    # Parsing precedes any native handle acquisition or process creation.
    pin=None;outcome=None;image=None;status='REJECTED'
    digest=trusted_sha256 if isinstance(trusted_sha256,str) and HEX.fullmatch(trusted_sha256) else ''
    try:
        v=profile(raw,trusted_sha256)
        require(cancel is None or isinstance(cancel,threading.Event))
        if cancel is not None and cancel.is_set():return LaunchResult('CANCELLED',digest,None,None)
        pin=PinnedImage(v['executable'],v['cwd']);image=pin.identity
        b=v['budgets']
        outcome=native_capture.capture([v['executable']['path'],*v['arguments']],v['cwd'],environment=v['environment'],timeout=b['timeout'],limit=b['output_bytes'],max_processes=b['processes'],max_memory=b['commit_bytes'],cancel=cancel,shutdown_grace=v.get('shutdown',{}).get('grace_seconds'))
        status='CAPTURED'
    except (Rejected,OSError,ValueError):
        status='REJECTED'
    finally:
        if pin is not None and not pin.close():status='IDENTITY_CLEANUP_FAILED'
    return LaunchResult(status,digest,image,outcome)
````

## Qualification source: test_shutdown.py (AUTHORED / NOT_ADMITTED)

SHA-256: `6abe6e0507aea4d6e40ad2aa4228fc0cedb622cbdfefb074bd8e672186aedb79`

````python
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
````

## Qualification source: test_capture.py (AUTHORED / NOT_ADMITTED)

SHA-256: `6f7c84136252a979a0a2ce6049059f37c5f28a0ca5498d1fa7c9a33aeed0c8c7`

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
````

## Canonical rebuild and integration139

Observed 2026-09-09T01:54:26.586723Z. All six files reconstruct exactly. The rebuilt tree passes
19 shutdown +20 governed-launch +23 capture/resource tests, zero skips. The
silent-descendant case retains its timeout/tree assertions and adds owned-handle
termination and bounded directory-release proof. Same-user event amplification
is retained as an explicit limitation test, never called a security admission.
Execution validator1.3.1 materialized15files; capability-gap0.1.0 composed4files
and self-check3positive/6negative PASS. Product owner/profiles unchanged.
Preflight139 is pending at this integration checkpoint.

| Evidence | SHA-256 |
|---|---|
| shutdown-tests.json | `dbcbb9b318aa5407496d79a13ab1411481a0a0e90e1f5c45d893c21b981509f4` |
| shutdown-green.json | `6b0586d112235212d1a8b31b663dd99b84706630eb4de97ea6e80463f6e4c549` |
| native-before599.py | `5571da0dc505b9fd20d07e5909605e232104af335db2df77c152812c9e28b510` |
| tests-before599.py | `de4f9d44bfa4328940adff10edb3e9bfd554937768930e4b1368b682a7b6a590` |
| event_access_probe.py | `1725b2cbfa77a906ef270b4d84ba95817a4847a7f0857a357d94266ca0d04ac8` |
| event-access-probe.json | `757f92b693046b51ad30136618a085a1170b3d881e46697bf3fec3702fb9f522` |
| capture-regression.json | `0e667587d5badc5f067a44bf27a440ed586a1837b838b6409532ab3ff4e06518` |
| capture-regression.log | `85c7f96913e8e24c57c5b29bbbd983f0adfd95e52ea09ff0ff422797c7e8c0b9` |
| test-capture-before600.py | `1597f76517f0fe9e95a77f31f979f35e5ecc11aed8f83550d99d9a7f1b650e08` |
| capture-green.json | `c4b869ed05b0f2082f785da240aab9346736b414b47724192311db00336fa798` |
| rebuilt-shutdown.json | `6b0586d112235212d1a8b31b663dd99b84706630eb4de97ea6e80463f6e4c549` |
| rebuilt-launch.json | `9dc3c8efed68851cf6170b0706956b637ab147fd4a1987fce53099d9da149d36` |
| rebuilt-capture.json | `c4b869ed05b0f2082f785da240aab9346736b414b47724192311db00336fa798` |
| canonical-parity.json | `ca989008e00dafda9c02b0db644dde987672bdc51b2fd037dd87aafcd0763e04` |

## Closure140 — shutdown qualification complete, full service gates open

Preflight139 completed: 154 executed checks PASS, with structural inventory
162 packs/1461 materializable files/786 Markdown and 53 profiles.
Docker remains unavailable. Every network/provider/target opt-in skip remains
unexecuted; local structural/test PASS is not production readiness. This closure
supersedes the pending statement at integration139 without rewriting its history.

The six reconstructed qualification files did not change after 19 shutdown,
20 launch and23 capture/resource tests passed with zero skips. No product pack,
selected profile or runtime changed. V345 base lifecycle evidence remains reused,
not newly counted. FAIL599 is recovered with the trusted-worker/same-user-access
limitation and lifecycle-only result label; FAIL600 is recovered by exact owned
descendant termination and finite fixture resource-release evidence.

Next: durable operational result delivery and actual selected worker lifecycle
mounting. Main-image identity, v2 profile and stop/grace experiments do not admit
trusted profile selection, DLL/runtime/SCA, service identity, restart policy,
hostile-worker sandbox or target collector/retention/alerts/recovery/acceptance.
Full G0–G8 service admission remains CONDITIONED. No global completion percentage.

Final successor decision/receipt: gap-record140.json and gap-receipt140.json in
this stage. Integration139 artifacts remain unchanged historical evidence.

| Closure evidence | SHA-256 |
|---|---|
| preflight139.json | `865b5c929fa40af9839eb3e2dfdcac71ca3e209515724b9a0734ef7fea523297` |
| preflight139.log | `d5e2f5326c22e686bbfe0beb8bb46768f4c5f65deb33db2ea814615fe2fc35ad` |
| structural-closure139.json | `908b15961324c789428e90a22153db1e807cddb2d0d0a573433d8bc91cb79c21` |
| gap-record139.json | `d739fa50aa2edbc6609b8c8d6a7373c7a708cd5cb6ed6163c078bf7ab7edfc10` |
| gap-receipt139.json | `66091a57f19d22b9a8dec9da9b6b09a5d480a5e82a5b9d8a88ac5c37a40aad5e` |
