# V347 — governed native executable launch qualification

Library maintenance / T2809, continued from validated checkpoint 136. AUTHORED
local qualification only, CONDITIONED / NOT_ADMITTED as product. No product pack,
profile, service registration, account, provider, dependency or runtime update.

## Delta and trust boundary

V346 supplies job-owned I/O and cleanup. This successor takes bounded profile
bytes and a separately supplied trusted SHA-256. Exact schema, duplicate-key
rejection and input limits bind executable path/hash/size, argv, cwd, explicit
four-variable environment and finite process/output/commit/time budgets. The
digest must come from an already trusted owner; computing it from untrusted input
does not authorize that input. A digest is not a signature or admission decision.

The image and each executable/cwd directory component are opened in order with
noninheritable Windows handles. Reparse components, alternate streams, device/UNC/
relative/ambiguous paths and a mismatched actual handle path are refused. The
selected regular disk file is hashed through its held handle and checked against
the profile size. Read-only sharing retains the file and directory handles until
the complete capture/cleanup returns. Existing writers, writes through hardlinks,
file replacement and directory rename are exercised locally. No copy-then-launch
fallback, PATH lookup or shell is used.

The existing capture API gains only an optional explicit environment argument;
its legacy default and all V346 regression cases are retained. The governed entry
always supplies the four profile values and admits no ambient PATH or custom
secret-bearing variable. This binds strings; it does not authenticate contents of
the OS, TEMP directories, DLLs, scripts, plugins or runtime dependency graph.

LaunchResult separates profile/identity rejection, identity cleanup failure and
captured process outcome. CAPTURED is not business success: refund exit 1 remains
exit 1, and a timed-out captured process retains TIMED_OUT plus tree cleanup.
Pre-cancel performs no native opens/launch; cancellation during capture uses V346.

Limits: trusted local drive mapping/OS/admin, same Windows x64 ABI and pinned
CPython 3.14.4. Not a sandbox, signature/ACL policy, anti-admin boundary, complete
runtime/dependency admission, service identity, durable result store, cooperative
shutdown or restart policy. Hashing is bounded by 1 GiB of local file bytes but
is not a strict wall-clock I/O deadline. Resource/time defaults remain synthetic
qualification values. No product incorporation until compatible G0–G8 admission.

## Failure retained

FAIL598 / LIB-FAIL-2314: initial suite ran 18 tests, 17 passed and one assertion
failed. Python CRT open(r+b) was refused with errno EACCES but winerror=None; the
test incorrectly required a Win32 code. Preserve launch-tests.log/json and
test-before598.py. The successor adds direct CreateFileW/GENERIC_WRITE refusal
with GetLastError 32 and independently checks CRT errno or native winerror for
each attempted write/delete/rename. Successful mutation still fails the test.
The corrected first suite passes 18 tests; two additional concurrency/leak tests
are included below and require fresh qualification, not an inferred PASS.

## Current official semantics (consulted, no upstream source copied)

- [CreateFileW](https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-createfilew):
  share modes remain in effect while the handle is open; omitted write/delete
  sharing restricts conflicting handles; OPEN_REPARSE_POINT permits inspection.
- [GetFinalPathNameByHandleW](https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-getfinalpathnamebyhandlew):
  obtain the opened handle's path and reject aliases outside this local contract.
- [CreateProcessW](https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-createprocessw):
  explicit executable, argv and environment; broader loader dependencies remain
  separate from the main image's selected hash.

Internal authorities: SECURITY_SRE_CLOUD_INFRASTRUCTURE.md least privilege,
supply chain and lifecycle boundaries; ENGINEERING_EXECUTION_PLAYBOOK.md;
PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md. Candidate gates concern only the
narrow observed boundary; full service G3–G8 remain incomplete. Provenance AUTHORED,
workspace-owner license, no external source acquisition or product license claim.

## Reproduction and progress accounting

Stage: C:/Users/NL/AppData/Local/Temp/elite-v347-254a34244df64eb6986ac22f885f1604. Use CPython 3.14.4. Reconstruct the three source blocks
below exactly. Reuse unchanged job_probe.py from the current V345 qualification
block and test_capture.py from the V346 qualification block. Supply the already
verified refund fixture JSON explicitly; it binds the real local binary/hash.
Run test_governed_launch.py <new-report.json> <refund-fixture.json>; then run the
23 V346 tests with their explicit fixture path/hash. Preserve independent logs.

Current roadmap has 11 owner tasks: T2800 complete, T2801/T2809 blocked and eight
pending. Ten open tasks are not 10 equally sized units of work and cannot provide
a global remaining percentage or completion date. This wave advances T2809 but
does not close it. Other owners remain unchanged; no 48-surface completion claim.

## Source files — qualification apparatus only

## Qualification source: native_capture.py (AUTHORED / NOT_ADMITTED)

SHA-256: `f54c6676529d29cb3395b0a97ee95b1dec1ba4414480410660711bc1ab103939`

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

def capture(argv,cwd,*,timeout=5.0,limit=10*1024*1024,cancel=None,max_processes=8,max_memory=256*1024*1024,environment=None):
    # Qualification API accepts only explicit absolute executables; no PATH/shell.
    if not isinstance(argv,list) or not argv or len(argv)>129 or any(not isinstance(x,str) or '\0' in x or len(x)>4096 for x in argv):
        raise ValueError('INVALID_ARGUMENTS') from None
    if not Path(argv[0]).is_absolute() or not Path(cwd).is_absolute():raise ValueError('ABSOLUTE_PATH_REQUIRED') from None
    if type(timeout) not in (int,float) or not math.isfinite(timeout) or not 0<timeout<=7200:raise ValueError('INVALID_TIMEOUT') from None
    if type(limit) is not int or not 1<=limit<=10*1024*1024:raise ValueError('INVALID_OUTPUT_LIMIT') from None
    if type(max_processes) is not int or not 1<=max_processes<=64:raise ValueError('INVALID_PROCESS_LIMIT') from None
    if type(max_memory) is not int or not 32*1024*1024<=max_memory<=1024*1024*1024:raise ValueError('INVALID_MEMORY_LIMIT') from None
    if cancel is not None and not isinstance(cancel,threading.Event):raise ValueError('INVALID_CANCELLATION') from None
    if environment is not None:
        if not isinstance(environment,dict) or set(environment)!={'SYSTEMROOT','WINDIR','TEMP','TMP'} or any(not isinstance(v,str) or not v or '\0' in v or len(v)>4096 for v in environment.values()):
            raise ValueError('INVALID_ENVIRONMENT') from None
        environment=dict(environment)
    if cancel is not None and cancel.is_set():return Outcome('CANCELLED',None,b'',b'',True)
    deadline=time.monotonic()+timeout;owner=None;buffers=[bytearray(),bytearray()];status='LAUNCH_FAILED';code=None;empty=False
    try:
        owner=OwnedCapture(argv,Path(cwd),clean_environment() if environment is None else environment,max_processes=max_processes,max_memory=max_memory)
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

## Qualification source: governed_launch.py (AUTHORED / NOT_ADMITTED)

SHA-256: `e8ce23a923fccfef2417aad652989520cb3972d5a05bd70ed3b3e99d3873408f`

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
    require(isinstance(v,dict) and set(v)=={'schema','id','executable','arguments','cwd','environment','budgets'})
    require(v['schema']=='elite-native-launch-qualification/v1' and isinstance(v['id'],str) and re.fullmatch(r'[a-z0-9][a-z0-9-]{0,63}',v['id']))
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
        outcome=native_capture.capture([v['executable']['path'],*v['arguments']],v['cwd'],environment=v['environment'],timeout=b['timeout'],limit=b['output_bytes'],max_processes=b['processes'],max_memory=b['commit_bytes'],cancel=cancel)
        status='CAPTURED'
    except (Rejected,OSError,ValueError):
        status='REJECTED'
    finally:
        if pin is not None and not pin.close():status='IDENTITY_CLEANUP_FAILED'
    return LaunchResult(status,digest,image,outcome)
````

## Qualification source: test_governed_launch.py (AUTHORED / NOT_ADMITTED)

SHA-256: `7000903fd4e56fc07dbdc3ebad6436f432902b32f28009be0b55c7c6e49ab3de`

````python
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
````

## Canonical rebuild and integration 137

Observed 2026-09-09T01:32:54.285964Z. All five qualification files reconstruct byte-for-byte: three
current blocks plus two unchanged V345/V346 apparatus files. From this rebuilt
tree, 20 launch tests and 23 capture/resource tests PASS with zero
skips. Sixteen concurrent captures pass; 24 failed-hash acquisitions retain the
same warmed-up native process handle count. The real pinned refund binary retains
exit 1 and REFUND_HOST_FAILED. Python argv quoting and explicit environment are
tested in a real subprocess; no PATH/private parent variables leak into it.

No product owner or selected profile changed. Execution validator 1.3.1 was
materialized again (15 files), and capability-gap 0.1.0 self-checks pass three
positive/six negative cases. FAIL598 is RECOVERED as a test-oracle error with an
additional direct native denial probe. Product service qualification remains open.
Preflight 137 is pending at this integration checkpoint.

| Evidence | SHA-256 |
|---|---|
| launch-tests.json | `91f2165a0ff3a33561926f043bf1a30e663c8f475896a2b7a8910add51384cfa` |
| launch-tests.log | `45cc6dcf0319e471616efdc873507cca210f6b6bfa5a64463844e717995f982f` |
| test-before598.py | `051ea31fafaa75ab6890f058d57a768e43e340b7900b26104a7d77f2fe332410` |
| launch-green.json | `24aaa7d92726b23d574d55b767813894d3831e0a9310df9eb86d1ac9a32c333c` |
| rebuilt-launch-tests.json | `9dc3c8efed68851cf6170b0706956b637ab147fd4a1987fce53099d9da149d36` |
| rebuilt-launch-tests.log | `f04e0057ae76b8024e7c667557cbc020a0c4a2c96f142e49119baeec50fe0bc6` |
| rebuilt-capture-tests.json | `c4b869ed05b0f2082f785da240aab9346736b414b47724192311db00336fa798` |
| rebuilt-capture-tests.log | `af3bc46b5afec94a0f15e6fadfa7c4819f9711f9b2260a20c870d421b5136b57` |
| canonical-parity.json | `90cc1021c4892114906f46af8708153f9bf172b09956ce0a9d063dc6bf5c6683` |
| baseline136.json | `242aa70c7c0f5a62a40548241105ef1c287479ee766a6b82305431462d75f122` |

## Closure 138 — governed launch qualification complete, service gates open

Preflight 137 completed with 154 executed checks PASS, including structural
inventory 162 packs/1461 materializable files/785 Markdown and
53 profiles. Docker remains unavailable. Network/provider/target opt-in
gates remain unexecuted where marked; this is not an overall production PASS.
The previous pending statement records integration 137 and is superseded here.

The qualified code did not change after its clean reconstruction and tests.
Twenty launch-profile/file-identity tests and 23 retained capture/resource tests
PASS with zero skips; the V345 base's unchanged 15 lifecycle observations remain
reused evidence, not newly counted executions. Product SECURE-OPS 1.1.3 and all
selected profiles stay unchanged. No new runtime, service or provider activation.

FAIL598 remains recovered only as the test-oracle error, with native refusal32
and CRT EACCES controls. Main-image pinning and immutable profile bytes do not
establish trusted profile selection/signing, DLL/runtime graph integrity, SCA,
privilege sandbox or service identity. No candidate promotion on test counts.

Next independent work is cooperative shutdown/lifecycle and durable operational
results, then complete compatible service admission and target gates. Ten roadmap
owner tasks remain open with partial progress; no remaining-time percentage can
be derived from these differently sized fronts or the 48 capability surfaces.

The final successor record/receipt are gap-record138.json and gap-receipt138.json
in this stage. The integration-137 artifacts remain unchanged historical snapshots.

| Closure evidence | SHA-256 |
|---|---|
| preflight137.json | `7849859d43e8b53fd8dfe2b94528efac1512598d3ef660acde724bd13b838aab` |
| preflight137.log | `3f49ddfa47e0d0cb1fecccd48a9add903b9b30e02eedfec3321cc29c785da1cb` |
| structural-closure137.json | `d38c07a34a36ece5fcc944809bded2b515804b90bcc2edefe33d75c57b1f2a3a` |
| gap-record137.json | `15166bfaae846070b9cd70d49a1f2686f97e74e4aae7304d31513e80410cc151` |
| gap-receipt137.json | `e650dbda0e4aa358247d4924646c45c540301367510301050778b9823dd4e8eb` |
