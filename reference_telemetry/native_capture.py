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

def capture(argv,cwd,*,timeout=5.0,limit=10*1024*1024,cancel=None,max_processes=8,max_memory=256*1024*1024,environment=None,shutdown_grace=None,reference_environment=False):
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
    if type(reference_environment) is not bool or (reference_environment and environment is None):raise ValueError('CAPTURE_CONFIG_REJECTED')
    expected={'SYSTEMROOT','WINDIR','TEMP','TMP'}
    if reference_environment:expected |= {'DATABASE_URL','REFUND_WORKER_ID','REFUND_TELEMETRY_CONFIG','REFUND_TELEMETRY_CONFIG_SHA256'}
    if environment is not None:
        if not isinstance(environment,dict) or set(environment)!=expected or any(not isinstance(v,str) or not v or '\0' in v or len(v)>4096 for v in environment.values()):
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
