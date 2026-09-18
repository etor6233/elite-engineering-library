# Windows Reference Telemetry Runtime

## 1. Metadata

```yaml
pack_id: "WINDOWS-REFERENCE-TELEMETRY-RUNTIME"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un runner opt-in y finito que verifica observabilidad integrada del main y CLI de referencia: identidad mTLS, logs/métricas/trazas, alertas Prometheus, retención por rotación, pérdida del colector, parada y recuperación local, con runtime/migraciones fijados y PostgreSQL nuevo. No instala un servicio productivo ni ejecuta pagos."
stacks: ["Windows x64 trusted local account", "Python 3.14 stdlib", "Go 1.26.8", "PostgreSQL 18.6", "minimal official-source OpenTelemetry Collector 0.160.0", "official-source Prometheus 3.14.0"]
compatible_with: ["GO-OFFICIAL-RETURN-REFUND-WORKER 0.1.5", "FRANCHISE_COMPLETE_PACK_PLAN selected migrations", "SECURE-OPS-DELIVERY-CORE 1.1.3"]
incompatible_with: ["existing databases", "provider credentials", "production service or secret custody", "untrusted same-user processes", "automatic retries", "unverified runtime lock"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://opentelemetry.io/docs/specs/otlp/", "https://opentelemetry.io/docs/collector/configuration/", "https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/", "https://learn.microsoft.com/en-us/windows/win32/procthread/job-objects", "https://docs.python.org/3/library/subprocess.html"]
verified_at: "2026-09-10"
```

## 2. Applicability

Library reference acceptance on the current trusted Windows x64 account only.
All supplied runtimes require exact source/build/license/security qualification.
The runner accepts a separately trusted digest of its complete runtime lock;
it does not turn arbitrary executable hashes into admission. All business rows
are explicitly synthetic and live in a fresh owned cluster whose data directory
is checked before mutation. Every real project retains its target gates.

## 3. Architecture contract

The native launcher pins main-image bytes/path, exact environment and budgets,
creates the process inside a Job Object, supplies only selected handles, captures
finite output and keeps a durable local reservation/result. Unknown delivery
never authorizes replay. The trusted manual-reset stop event is cooperative;
after finite grace, the owned job is terminated. This is not hostile-user isolation.

Actual Go main exports five closed outcomes through authenticated OTLP/HTTP.
The canonical worker owns its data/transport policy. Collector persists three
signals and exports counters to Prometheus. Actual rule evaluation, loss and
recovery are observed; no notification is sent to another person or service.
The runner uses a private reference directory, ephemeral keys, hash-bound config,
1MiB file rotation,2backups and1day retention. The2000-observation CLI probe
verifies rotation/expiry and retained bytes after process restart. No claim of
power-loss durability, PITR, full TSDB disk cap or production acceptance follows.

## 4. Exact file manifest

```text
CREATE reference_telemetry/governed_launch.py
CREATE reference_telemetry/job_probe.py
CREATE reference_telemetry/native_capture.py
CREATE reference_telemetry/README.md
CREATE reference_telemetry/result_store.py
CREATE reference_telemetry/run_reference_runtime.py
CREATE reference_telemetry/runtime-lock.example.json
CREATE reference_telemetry/seed.sql
CREATE reference_telemetry/test_capture.py
CREATE reference_telemetry/test_governed_launch.py
CREATE reference_telemetry/test_reference_profile.py
CREATE reference_telemetry/test_result_store.py
CREATE reference_telemetry/test_shutdown.py
```

## 5. Materialization blocks

### FILE: `reference_telemetry/governed_launch.py`
```yaml
block_id: "WINDOWS-REFERENCE-TELEMETRY:governed-launch:v1"
operation: CREATE
provenance: AUTHORED
source: "local V345-V351 native qualification and V372 integrated reference; official contracts govern behavior, no vendor implementation authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a1a5c67c2116c6e084859ef2017bcff7a75d20ae8279087e096851b2f29b2599"
variables: []
secrets_allowed: false
```
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
from urllib.parse import urlsplit, parse_qsl
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
    require(v.get('schema') in ('elite-native-launch-qualification/v1','elite-native-launch-qualification/v2','elite-native-launch-qualification/v3'))
    if v['schema'].endswith(('/v2','/v3')):
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
    system={'SYSTEMROOT','WINDIR','TEMP','TMP'}
    runtime={'DATABASE_URL','REFUND_WORKER_ID','REFUND_TELEMETRY_CONFIG','REFUND_TELEMETRY_CONFIG_SHA256'} if v['schema'].endswith('/v3') else set()
    require(isinstance(env,dict) and set(env)==system|runtime)
    for name in system:local_path(env[name])
    if runtime:
        # Explicit reference-only identity: trusted local account, synthetic
        # loopback PostgreSQL, no payment-provider credentials or remote DSN.
        require(all(isinstance(env[k],str) and 1<=len(env[k])<=4096 and not any(ord(c)<32 for c in env[k]) for k in runtime))
        try:u=urlsplit(env['DATABASE_URL']);port=u.port
        except ValueError:raise Rejected('LAUNCH_PROFILE_REJECTED') from None
        require(u.scheme=='postgres' and u.hostname=='127.0.0.1' and u.username=='postgres' and u.password is None and port is not None and 1024<=port<=65535 and not u.fragment)
        require(re.fullmatch(r'/elite_refund_test_[0-9a-f]{32}',u.path) and parse_qsl(u.query,strict_parsing=True)==[('sslmode','disable')])
        require(re.fullmatch(r'elite-reference-[a-z0-9-]{1,40}',env['REFUND_WORKER_ID']))
        local_path(env['REFUND_TELEMETRY_CONFIG']);require(HEX.fullmatch(env['REFUND_TELEMETRY_CONFIG_SHA256']))
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
        outcome=native_capture.capture([v['executable']['path'],*v['arguments']],v['cwd'],environment=v['environment'],timeout=b['timeout'],limit=b['output_bytes'],max_processes=b['processes'],max_memory=b['commit_bytes'],cancel=cancel,shutdown_grace=v.get('shutdown',{}).get('grace_seconds'),reference_environment=v['schema'].endswith('/v3'))
        status='CAPTURED'
    except (Rejected,OSError,ValueError):
        status='REJECTED'
    finally:
        if pin is not None and not pin.close():status='IDENTITY_CLEANUP_FAILED'
    return LaunchResult(status,digest,image,outcome)
````

### FILE: `reference_telemetry/job_probe.py`
```yaml
block_id: "WINDOWS-REFERENCE-TELEMETRY:job-probe:v1"
operation: CREATE
provenance: AUTHORED
source: "local V345-V351 native qualification and V372 integrated reference; official contracts govern behavior, no vendor implementation authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e33dbd14b0cd09262897c98263de0d572ac0b66673df6bce2f62c3bf9ade316c"
variables: []
secrets_allowed: false
```
````python
"""AUTHORED isolated Win32 qualification experiment; not a product supervisor.
No provider, service registration, restart, or global process enumeration.
"""
import ctypes as C
from ctypes import wintypes as W
from pathlib import Path
import os,sys,json,time,subprocess,hashlib
K=C.WinDLL('kernel32',use_last_error=True)
SIZE=C.c_size_t
class Basic(C.Structure):
 _fields_=[('process_time',C.c_longlong),('job_time',C.c_longlong),('flags',W.DWORD),('min_ws',SIZE),('max_ws',SIZE),('active',W.DWORD),('affinity',SIZE),('priority',W.DWORD),('scheduling',W.DWORD)]
class IO(C.Structure):
 _fields_=[(n,C.c_ulonglong) for n in ('read_ops','write_ops','other_ops','read_bytes','write_bytes','other_bytes')]
class Extended(C.Structure):
 _fields_=[('basic',Basic),('io',IO),('process_memory',SIZE),('job_memory',SIZE),('peak_process',SIZE),('peak_job',SIZE)]
class Startup(C.Structure):
 _fields_=[('cb',W.DWORD),('reserved',W.LPWSTR),('desktop',W.LPWSTR),('title',W.LPWSTR),*[(n,W.DWORD) for n in ('x','y','width','height','xchars','ychars','fill','flags')],('show',W.WORD),('reserved_bytes',W.WORD),('reserved_ptr',C.c_void_p),('stdin',W.HANDLE),('stdout',W.HANDLE),('stderr',W.HANDLE)]
class StartupEx(C.Structure):
 _fields_=[('startup',Startup),('attributes',C.c_void_p)]
class Process(C.Structure):
 _fields_=[('process',W.HANDLE),('thread',W.HANDLE),('pid',W.DWORD),('tid',W.DWORD)]
def api(name,args,result):
 f=getattr(K,name);f.argtypes=args;f.restype=result;return f
create_job=api('CreateJobObjectW',[C.c_void_p,W.LPCWSTR],W.HANDLE)
set_job=api('SetInformationJobObject',[W.HANDLE,C.c_int,C.c_void_p,W.DWORD],W.BOOL)
assign=api('AssignProcessToJobObject',[W.HANDLE,W.HANDLE],W.BOOL)
close=api('CloseHandle',[W.HANDLE],W.BOOL)
create=api('CreateProcessW',[W.LPCWSTR,W.LPWSTR,C.c_void_p,C.c_void_p,W.BOOL,W.DWORD,C.c_void_p,W.LPCWSTR,C.POINTER(Startup),C.POINTER(Process)],W.BOOL)
resume=api('ResumeThread',[W.HANDLE],W.DWORD)
terminate=api('TerminateProcess',[W.HANDLE,W.UINT],W.BOOL)
wait=api('WaitForSingleObject',[W.HANDLE,W.DWORD],W.DWORD)
open_process=api('OpenProcess',[W.DWORD,W.BOOL,W.DWORD],W.HANDLE)
in_job=api('IsProcessInJob',[W.HANDLE,W.HANDLE,C.POINTER(W.BOOL)],W.BOOL)
initialize=api('InitializeProcThreadAttributeList',[C.c_void_p,W.DWORD,W.DWORD,C.POINTER(SIZE)],W.BOOL)
update=api('UpdateProcThreadAttribute',[C.c_void_p,W.DWORD,SIZE,C.c_void_p,SIZE,C.c_void_p,C.c_void_p],W.BOOL)
delete=api('DeleteProcThreadAttributeList',[C.c_void_p],None)
def require(ok,name):
 if not ok:raise RuntimeError(name+' failed code='+str(C.get_last_error()))
def write(p,obj):
 data=json.dumps(obj).encode()
 with p.open('xb') as f:f.write(data);f.flush();os.fsync(f.fileno())
def await_file(p,timeout=5):
 until=time.monotonic()+timeout
 while time.monotonic()<until:
  if p.is_file():
   try:return json.loads(p.read_text())
   except json.JSONDecodeError:pass
  time.sleep(.01)
 raise RuntimeError('fixture handshake timed out')
def job():
 h=create_job(None,None);require(h,'CreateJobObject')
 info=Extended();info.basic.flags=0x2000 # KILL_ON_JOB_CLOSE, no breakaway.
 if not set_job(h,9,C.byref(info),C.sizeof(info)):
  close(h);raise RuntimeError('SetInformationJobObject failed')
 return h
def suspended(root,code,h):
 # Membership is supplied to CreateProcess itself. There is no running or
 # suspended unowned-child window between CreateProcess and AssignProcess.
 si=StartupEx();si.startup.cb=C.sizeof(si);pi=Process()
 size=SIZE()
 require(not initialize(None,1,0,C.byref(size)) and C.get_last_error()==122 and size.value>0,'AttributeSize')
 storage=C.create_string_buffer(size.value)
 require(initialize(storage,1,0,C.byref(size)),'InitializeAttributes')
 handles=(W.HANDLE*1)(h)
 try:
  require(update(storage,0,0x0002000d,handles,C.sizeof(handles),None,None),'JobListAttribute')
  si.attributes=C.cast(storage,C.c_void_p)
  cmd=C.create_unicode_buffer(subprocess.list2cmdline([str(Path(sys.executable).resolve()),'-I','-c',code]))
  env={k:v for k,v in os.environ.items() if k.upper() in {'SYSTEMROOT','WINDIR','TEMP','TMP'}}
  block=C.create_unicode_buffer('\0'.join(k+'='+v for k,v in sorted(env.items()))+'\0\0')
  require(create(str(Path(sys.executable).resolve()),cmd,None,None,False,0x4|0x08000000|0x400|0x80000,block,str(root),C.cast(C.byref(si),C.POINTER(Startup)),C.byref(pi)),'CreateProcessWithJob')
  return pi
 finally:delete(storage)
def tree_code():
 nested="import os,time,json,pathlib;pathlib.Path('nested.json').write_text(json.dumps({'pid':os.getpid()}));time.sleep(30)"
 return f"import subprocess,sys,time;subprocess.Popen([sys.executable,'-I','-c',{nested!r}],creationflags=0x08000000);time.sleep(30)"
def owner(root):
 h=job();pi=suspended(root,tree_code(),h)
 try:
  require(resume(pi.thread)!=0xffffffff,'ResumeThread')
  require(close(pi.thread),'CloseThread');pi.thread=None
  nested=await_file(root/'nested.json');write(root/'owner-ready.json',{'root':pi.pid,'nested':nested['pid']})
  await_file(root/'owner-exit.json')
  os._exit(0) # Deliberately bypass Python finally to exercise kernel cleanup.
 finally:
  close(h)
  if wait(pi.process,0)==258:terminate(pi.process,99);wait(pi.process,5000)
  close(pi.process)
  if pi.thread:close(pi.thread)
def check_live(handle,h=None):
 require(wait(handle,0)==258,'owned fixture alive')
 if h:
  member=W.BOOL();require(in_job(handle,h,C.byref(member)) and member.value,'job membership')
def normal(root):
 h=job();pi=None;nested_handle=None
 try:
  pi=suspended(root,tree_code(),h)
  require(not (root/'nested.json').exists(),'not started while suspended')
  require(resume(pi.thread)!=0xffffffff,'ResumeThread');close(pi.thread);pi.thread=None
  info=await_file(root/'nested.json');nested_handle=open_process(0x100000|0x1000,False,info['pid']);require(nested_handle,'OpenOwnedNested')
  check_live(pi.process,h);check_live(nested_handle,h)
  require(close(h),'CloseJob');h=None
  require(wait(pi.process,5000)==0,'root terminated');require(wait(nested_handle,5000)==0,'nested terminated')
  return {'case':'last_handle_close','root_signaled':True,'nested_signaled':True,'assigned_before_resume':True}
 finally:
  if h:close(h)
  if pi:
   if wait(pi.process,0)==258:terminate(pi.process,99);wait(pi.process,5000)
   close(pi.process)
   if pi.thread:close(pi.thread)
  if nested_handle:close(nested_handle)
def assignment_refusal(root):
 try:
  pi=suspended(root,"from pathlib import Path;Path('must-not-run').write_text('bad')",None)
 except RuntimeError as error:
  require(str(error).startswith('CreateProcessWithJob failed'),'creation refused by OS')
  require(not (root/'must-not-run').exists(),'no child effect')
  return {'case':'assignment_refusal','no_child_effect':True,'creation_failed':True}
 else:
  terminate(pi.process,99);wait(pi.process,5000);close(pi.thread);close(pi.process)
  raise RuntimeError('invalid job unexpectedly accepted')
def breakaway_refusal(root):
 h=job();pi=None
 escaped="from pathlib import Path;Path('escaped').write_text('bad')"
 body="import subprocess,sys,json,pathlib\ntry:\n subprocess.Popen([sys.executable,'-I','-c',"+repr(escaped)+"],creationflags=0x01000000|0x08000000)\nexcept OSError as error:\n pathlib.Path('refused.json').write_text(json.dumps({'winerror':error.winerror}))"
 try:
  pi=suspended(root,body,h);require(resume(pi.thread)!=0xffffffff,'ResumeThread')
  require(wait(pi.process,5000)==0,'breakaway fixture finished')
  info=await_file(root/'refused.json');require(info['winerror']==5,'breakaway denied with access error')
  require(not (root/'escaped').exists(),'no escaped child effect')
  return {'case':'breakaway_refusal','winerror':5,'no_escaped_effect':True}
 finally:
  close(h)
  if pi:
   if wait(pi.process,0)==258:terminate(pi.process,99);wait(pi.process,5000)
   close(pi.thread);close(pi.process)
def parent_death(root):
 flags=subprocess.CREATE_NO_WINDOW
 env={k:v for k,v in os.environ.items() if k.upper() in {'SYSTEMROOT','WINDIR','TEMP','TMP'}}
 wrapper=subprocess.Popen([sys.executable,'-I',str(Path(__file__).resolve()),'owner',str(root)],stdin=subprocess.DEVNULL,stdout=subprocess.PIPE,stderr=subprocess.PIPE,creationflags=flags,env=env)
 handles=[]
 try:
  info=await_file(root/'owner-ready.json')
  for key in ('root','nested'):
   h=open_process(0x100000|0x1000,False,info[key]);require(h,'OpenOwnedChild');handles.append(h);check_live(h)
  write(root/'owner-exit.json',{'exit':True})
  wrapper.wait(timeout=5);require(wrapper.returncode==0,'owner deliberate exit')
  for h in handles:require(wait(h,5000)==0,'kernel cleanup after owner death')
  return {'case':'abrupt_owner_exit','root_signaled':True,'nested_signaled':True,'python_finally_bypassed':True}
 finally:
  if wrapper.poll() is None:wrapper.kill();wrapper.wait(timeout=5)
  wrapper.stdout.close();wrapper.stderr.close()
  for h in handles:close(h)
def owner_before_resume(root):
 h=job();pi=suspended(root,"from pathlib import Path;Path('must-not-run').write_text('bad')",h)
 # Deliberately no explicit AssignProcessToJobObject and no ResumeThread.
 write(root/'ready.json',{'pid':pi.pid})
 await_file(root/'exit.json');os._exit(0)
def death_before_resume(root):
 env={k:v for k,v in os.environ.items() if k.upper() in {'SYSTEMROOT','WINDIR','TEMP','TMP'}}
 wrapper=subprocess.Popen([sys.executable,'-I',str(Path(__file__).resolve()),'owner-before-resume',str(root)],stdin=subprocess.DEVNULL,stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL,creationflags=subprocess.CREATE_NO_WINDOW,env=env)
 h=None
 try:
  info=await_file(root/'ready.json');h=open_process(0x100000|0x1000|1,False,info['pid']);require(h,'OpenOwnedSuspended');check_live(h)
  write(root/'exit.json',{'exit':True});wrapper.wait(timeout=5)
  require(wrapper.returncode==0,'owner exited');require(wait(h,5000)==0,'creation-bound child terminated')
  require(not (root/'must-not-run').exists(),'no child effect before resume')
  return {'case':'owner_death_before_resume','child_signaled':True,'no_user_effect':True,'no_post_creation_assignment':True}
 finally:
  if wrapper.poll() is None:wrapper.kill();wrapper.wait(timeout=5)
  if h:
   if wait(h,0)==258:terminate(h,99);wait(h,5000)
   close(h)
def main():
 if len(sys.argv)>1 and sys.argv[1]=='owner-before-resume':owner_before_resume(Path(sys.argv[2]));return
 if len(sys.argv)>1 and sys.argv[1]=='owner':owner(Path(sys.argv[2]));return
 root=Path(__file__).resolve().parent/'runs';root.mkdir()
 outcomes=[]
 for iteration in range(3):
  for method in (normal,assignment_refusal,parent_death,breakaway_refusal,death_before_resume):
   dest=root/f'{iteration}-{method.__name__}';dest.mkdir();start=time.monotonic();result=method(dest);result['seconds']=time.monotonic()-start;result['iteration']=iteration;outcomes.append(result)
 result={'schema':'elite-native-job-qualification/v1','status':'PASS','product_admitted':False,'python':sys.version,'os':sys.getwindowsversion().build,'sizes':{'extended':C.sizeof(Extended),'startup':C.sizeof(Startup),'process':C.sizeof(Process),'startup_ex':C.sizeof(StartupEx)},'outcomes':outcomes,'script_sha256':hashlib.sha256(Path(__file__).read_bytes()).hexdigest()}
 write(Path(__file__).with_suffix('.receipt.json'),result);print(json.dumps(result,indent=2))
if __name__=='__main__':main()
````

### FILE: `reference_telemetry/native_capture.py`
```yaml
block_id: "WINDOWS-REFERENCE-TELEMETRY:native-capture:v1"
operation: CREATE
provenance: AUTHORED
source: "local V345-V351 native qualification and V372 integrated reference; official contracts govern behavior, no vendor implementation authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f1e00ba6afe414d494617e3c69ab3c8733729ef47b83154567737cc023341269"
variables: []
secrets_allowed: false
```
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
````

### FILE: `reference_telemetry/README.md`
```yaml
block_id: "WINDOWS-REFERENCE-TELEMETRY:README:v1"
operation: CREATE
provenance: AUTHORED
source: "local V345-V351 native qualification and V372 integrated reference; official contracts govern behavior, no vendor implementation authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8cd8af13123d8e00a1f1a05fc3a5b2efeeb3fa0a835e5062c8893173209e291c"
variables: []
secrets_allowed: false
```
````markdown
# Windows reference telemetry runtime

This is an opt-in, finite reference acceptance runner for the library. It mounts
the actual refund main and CLI, official-source minimal Collector and Prometheus,
and a fresh isolated PostgreSQL cluster. It is not a production service installer,
secret manager, payment environment, notification recipient or training pipeline.
Use only on a trusted Windows x64 host under the current local account.

Materialize GO-OFFICIAL-RETURN-REFUND-WORKER0.1.5 as part of the compatible full
consumer composition. Build its two cmd binaries with the admitted Go1.26.8.
Acquire and qualify the exact Collector0.160.0 and Prometheus3.14.0 sources through
the maintenance-telemetry-source-inspection profile. The V372 evidence preserves
the exact OCB selection, generated source, module locks, build commands, security
patches, licenses, tests and artifact hashes; official release binaries are not
interchangeable with this minimal locally built distribution. Verify the local
PostgreSQL18.6 runtime independently. Do not fill the runtime lock from reputation.

Copy runtime-lock.example.json outside the source tree. Record each verified
binary path, byte size and SHA-256, every selected *.up.sql filename/hash from the
consumer, and the exact seed.sql hash. Bind that entire JSON by an independently
trusted SHA-256 at invocation. The example is intentionally rejected unchanged.
No runtime downloads, installations or external accounts occur in this runner.

Run:

```text
python -B reference_telemetry/run_reference_runtime.py --runtime-lock LOCK.json --runtime-lock-sha256 TRUSTED_SHA256 --consumer MATERIALIZED_CONSUMER --work-root OWNED_EVIDENCE_DIRECTORY
```

The evidence root must already exist and be short enough for the native path
protocol. A new UUID directory receives a current-user/SYSTEM ACL, ephemeral
certificates, exact configs, a fresh cluster and bounded process receipts. The
CA private key is never saved. Client/server keys expire after two hours and
must never be published. The cluster identity is checked before any mutation.
All owned processes stop in finally; files and failed evidence are retained.

Five host attempts verify bad-config rejection before claims, actual missing
provider blocking, DB-error signals, graceful stop, collector loss and recovery.
A handled step does not imply a successful refund. No provider key is inherited.
OTLP, metrics and the Prometheus API require mutual TLS; rogue/missing clients
fail. The separate fixed health response is bound to local HTTP.
Actual OTLP output proves log/span correlation and cumulative metric invariants.
Prometheus executes the error and collector-health rules and their recovery.
A2000-observation/8-instance CLI probe drives1MiB rotation,2backup retention and
expiry; retained files remain byte-identical after collector process restart.
This is local file/process recovery, not power-loss durability, PITR or DR.
Prometheus64MB/1h retention is a configured TSDB policy, not a hard total-disk cap.

Regression suites:

```text
python -B reference_telemetry/test_reference_profile.py policy.json
python -B reference_telemetry/test_capture.py capture.json REFUND_EXE REFUND_SHA256
python -B reference_telemetry/test_governed_launch.py launch.json REFUND_FIXTURE_JSON
python -B reference_telemetry/test_shutdown.py shutdown.json
python -B reference_telemetry/test_result_store.py result.json REFUND_FIXTURE_JSON
```

REFUND_FIXTURE_JSON contains path and sha256 for that exact executable. The
supervisor is a trusted-launcher primitive with finite capture, job limits,
cooperative stop and durable local metadata. It never retries an uncertain run.
It does not sandbox administrators or hostile same-user processes. Runtime,
service identity, secret custody, retention obligations, alert routing, workload,
security and acceptance must be requalified in every real consumer target.
````

### FILE: `reference_telemetry/result_store.py`
```yaml
block_id: "WINDOWS-REFERENCE-TELEMETRY:result-store:v1"
operation: CREATE
provenance: AUTHORED
source: "local V345-V351 native qualification and V372 integrated reference; official contracts govern behavior, no vendor implementation authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "464d52c47273c9309948edbda371c0587a38de43da41f4a3640f09b7d2f3c669"
variables: []
secrets_allowed: false
```
````python
"""AUTHORED Windows local crash-recovery qualification, NOT_ADMITTED product.

Trusted local store/OS/worker only. Atomic reservation forbids automatic replay
of an uncertain launch. No power-loss, WORM, anti-admin or business-success claim.
Only fixed metadata and digests persist; raw stdout/argv/environment do not.
"""
from pathlib import Path
from dataclasses import dataclass
import ctypes as C
import hashlib,json,os,re,uuid
import governed_launch as g

move_file=g.win.api('MoveFileExW',[g.W.LPCWSTR,g.W.LPCWSTR,g.W.DWORD],g.W.BOOL)
MAX_RECORD=16384
RUN=re.compile(r'[0-9a-f]{32}')
CAPTURE={'COMPLETED','TIMED_OUT','OUTPUT_LIMIT','CANCELLED','IO_FAILED','LAUNCH_FAILED','CLEANUP_FAILED','STOP_SIGNAL_FAILED'}
STOP={'NOT_REQUESTED','EXITED_DURING_GRACE','FORCED','SIGNAL_FAILED'}
def require(ok):
    if not ok:raise g.Rejected('RESULT_RECORD_REJECTED') from None
def hex64(v):return isinstance(v,str) and bool(g.HEX.fullmatch(v))
def digest(raw):return hashlib.sha256(raw).hexdigest()
def encode(value):
    raw=(json.dumps(value,sort_keys=True,separators=(',',':'),allow_nan=False)+'\n').encode()
    require(len(raw)<=MAX_RECORD);return raw

class DirectoryPin(g.PinnedImage):
    def __init__(self,path):
        self.handles=[]
        try:
            p=g.local_path(str(path))
            for directory in [*reversed(p.parents),p]:self._open(directory,True)
            self.path=Path(p)
        except BaseException:self.close();raise
    def __enter__(self):return self
    def __exit__(self,*_):require(self.close())
    def read(self,name):
        require(name in ('claim.json','result.json'))
        # Read through an opened, read-shared, non-reparse handle. Use a distinct
        # temporary pin so every early parse/validation path closes the file.
        with DirectoryPin(self.path) as reader:
            h,info=reader._open(self.path/name,False)
            size=info.size_high<<32|info.size_low;require(0<size<=MAX_RECORD)
            data=bytearray();buf=C.create_string_buffer(MAX_RECORD)
            while len(data)<size:
                count=g.W.DWORD();require(g.read_file(h,buf,size-len(data),C.byref(count),None) and 0<count.value<=size-len(data))
                data.extend(buf.raw[:count.value])
            return bytes(data)

def publish_file(directory,name,raw):
    require(name in ('claim.json','result.json') and type(raw) is bytes and 0<len(raw)<=MAX_RECORD)
    directory=Path(directory);g.local_path(str(directory))
    # Caller holds the directory chain. Same-directory path, no replace/copy/
    # delayed-reboot flags; failed temp files are retained for reconciliation.
    temp=directory/('part-'+uuid.uuid4().hex)
    with temp.open('xb',buffering=0) as f:
        remaining=memoryview(raw)
        while remaining:
            n=f.write(remaining);require(type(n) is int and 0<n<=len(remaining));remaining=remaining[n:]
        os.fsync(f.fileno())
    require(move_file(str(temp),str(directory/name),8))
    return digest(raw)

def parse(raw):
    require(type(raw) is bytes and 0<len(raw)<=MAX_RECORD)
    try:v=json.loads(raw.decode('utf-8'),object_pairs_hook=g.unique_pairs,parse_constant=lambda _:require(False))
    except (ValueError,UnicodeError,RecursionError):raise g.Rejected('RESULT_RECORD_REJECTED') from None
    require(isinstance(v,dict));return v
def validate_claim(v,run_id,profile_sha):
    require(set(v)=={'schema','run_id','profile_sha256','executable'})
    require(v['schema']=='elite-launch-claim/v1' and v['run_id']==run_id and v['profile_sha256']==profile_sha)
    e=v['executable'];require(isinstance(e,dict) and set(e)=={'sha256','bytes'} and hex64(e['sha256']) and type(e['bytes']) is int and 1<=e['bytes']<=1073741824)
def validate_result(v,claim,claim_sha):
    require(set(v)=={'schema','run_id','profile_sha256','claim_sha256','launch_status','image','outcome'})
    require(v['schema']=='elite-launch-result/v1' and v['run_id']==claim['run_id'] and v['profile_sha256']==claim['profile_sha256'] and v['claim_sha256']==claim_sha)
    require(v['launch_status'] in {'REJECTED','CANCELLED','CAPTURED','IDENTITY_CLEANUP_FAILED'})
    image=v['image']
    if image is not None:
        require(isinstance(image,dict) and set(image)=={'sha256','bytes'} and hex64(image['sha256']) and type(image['bytes']) is int and image==claim['executable'])
    o=v['outcome']
    if o is None:
        require(v['launch_status'] in {'REJECTED','CANCELLED'});return
    require(v['launch_status'] in {'CAPTURED','IDENTITY_CLEANUP_FAILED'} and image is not None)
    require(isinstance(o,dict) and set(o)=={'status','exit_code','tree_empty','shutdown_state','stdout','stderr'})
    require(o['status'] in CAPTURE and o['shutdown_state'] in STOP and type(o['tree_empty']) is bool)
    require(o['exit_code'] is None or type(o['exit_code']) is int and 0<=o['exit_code']<=0xffffffff)
    if o['status']=='COMPLETED':require(type(o['exit_code']) is int and o['tree_empty'] and o['shutdown_state']=='NOT_REQUESTED')
    if o['status']=='CLEANUP_FAILED':require(not o['tree_empty'])
    if o['shutdown_state']=='EXITED_DURING_GRACE':require(o['status'] in {'CANCELLED','TIMED_OUT','CLEANUP_FAILED'} and type(o['exit_code']) is int)
    for field in ['stdout','stderr']:
        stream=o[field];require(isinstance(stream,dict) and set(stream)=={'bytes','sha256'} and type(stream['bytes']) is int and 0<=stream['bytes']<=10485760 and hex64(stream['sha256']))
        if stream['bytes']==0:require(stream['sha256']==digest(b''))

@dataclass(frozen=True)
class Stored:
    state:str
    run_id:str
    receipt_sha256:str|None=None

def inspect(root,run_id,profile_sha):
    try:
        require(isinstance(run_id,str) and RUN.fullmatch(run_id) and hex64(profile_sha))
        with DirectoryPin(root) as base:
            slot=base.path/run_id
            if not slot.exists():return Stored('ABSENT',run_id)
            with DirectoryPin(slot) as pin:
                if not (slot/'claim.json').exists():
                    return Stored('INVALID' if (slot/'result.json').exists() else 'UNRESOLVED',run_id)
                claim_raw=pin.read('claim.json');claim=parse(claim_raw);validate_claim(claim,run_id,profile_sha)
                if not (slot/'result.json').exists():return Stored('UNRESOLVED',run_id)
                raw=pin.read('result.json');v=parse(raw);validate_result(v,claim,digest(claim_raw))
                return Stored('RECORDED',run_id,digest(raw))
    except (g.Rejected,OSError,ValueError,TypeError):return Stored('INVALID',run_id if isinstance(run_id,str) and RUN.fullmatch(run_id) else '')

def execute(root,run_id,profile_raw,trusted_profile_sha,*,cancel=None):
    # RECORDED means local receipt publication, not successful business work.
    # Existing slots, including empty/partial ones, NEVER authorize another launch.
    state='REJECTED';receipt_sha=None
    safe_id=run_id if isinstance(run_id,str) and RUN.fullmatch(run_id) else ''
    try:
        require(safe_id)
        profile=g.profile(profile_raw,trusted_profile_sha)
        with DirectoryPin(root) as base:
            slot=base.path/run_id
            try:slot.mkdir()
            except FileExistsError:return Stored('ALREADY_RESERVED',run_id)
            state='CLAIM_FAILED'
            with DirectoryPin(slot):
                claim={'schema':'elite-launch-claim/v1','run_id':run_id,'profile_sha256':trusted_profile_sha,'executable':{k:profile['executable'][k] for k in ['sha256','bytes']}}
                claim_raw=encode(claim);claim_sha=publish_file(slot,'claim.json',claim_raw)
                state='EXECUTION_UNCERTAIN'
                result=g.launch(profile_raw,trusted_profile_sha,cancel=cancel)
                state='RESULT_UNCERTAIN'
                out=result.outcome
                stream=lambda b:{'bytes':len(b),'sha256':digest(b)}
                value={'schema':'elite-launch-result/v1','run_id':run_id,'profile_sha256':trusted_profile_sha,'claim_sha256':claim_sha,'launch_status':result.status,'image':None if result.image is None else {k:result.image[k] for k in ['sha256','bytes']},'outcome':None if out is None else {'status':out.status,'exit_code':out.exit_code,'tree_empty':out.tree_empty,'shutdown_state':out.shutdown_state,'stdout':stream(out.stdout),'stderr':stream(out.stderr)}}
                validate_result(value,claim,claim_sha)
                receipt_sha=publish_file(slot,'result.json',encode(value));state='RECORDED'
        return Stored(state,run_id,receipt_sha)
    except Exception:
        # A publication/close error after execution cannot trigger a replay.
        return Stored('RESULT_UNCERTAIN' if state=='RECORDED' else state,safe_id)
````

### FILE: `reference_telemetry/run_reference_runtime.py`
```yaml
block_id: "WINDOWS-REFERENCE-TELEMETRY:run-reference-runtime:v1"
operation: CREATE
provenance: AUTHORED
source: "local V345-V351 native qualification and V372 integrated reference; official contracts govern behavior, no vendor implementation authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "997b21d7c4a9f4982c52785402bc14afd38359bb85ec37694b3f33669f965a87"
variables: []
secrets_allowed: false
```
````python
"""Reference-only complete executable mounting; never a real payment target."""
from pathlib import Path
import os,json,subprocess,threading,time,socket,uuid,hashlib,ssl,urllib.request,urllib.parse,sys,traceback,datetime,ipaddress
import argparse,re,csv
S=Path(__file__).resolve().parent
sys.path.insert(0,str(S));import result_store as rs
parser=argparse.ArgumentParser(description="Finite Windows reference telemetry acceptance. No existing database or provider access.")
parser.add_argument('--runtime-lock',required=True);parser.add_argument('--runtime-lock-sha256',required=True);parser.add_argument('--consumer',required=True);parser.add_argument('--work-root',required=True)
args=parser.parse_args()
def sha(p):
 h=hashlib.sha256()
 with p.open('rb') as f:
  for chunk in iter(lambda:f.read(1048576),b''):h.update(chunk)
 return h.hexdigest()
def load_inputs():
 with Path(args.runtime_lock).open('rb') as f:raw=f.read(65537)
 if len(raw)>65536 or not re.fullmatch('[0-9a-f]{64}',args.runtime_lock_sha256) or hashlib.sha256(raw).hexdigest()!=args.runtime_lock_sha256:raise ValueError('REFERENCE_INPUT_REJECTED')
 v=json.loads(raw,object_pairs_hook=rs.g.unique_pairs)
 names={'elite-otelcol','prometheus','promtool','refund-host','telemetry-reference','postgres','initdb','pg_ctl','psql'}
 if set(v)!={'schema','binaries','migrations','seed_sha256'} or v['schema']!='elite-reference-telemetry-runtime/v1' or set(v['binaries'])!=names:raise ValueError('REFERENCE_INPUT_REJECTED')
 bins={}
 for name,item in v['binaries'].items():
  if set(item)!={'path','sha256','bytes'}:raise ValueError('REFERENCE_INPUT_REJECTED')
  rs.g.local_path(item['path']);p=Path(item['path'])
  if p.suffix.lower()!='.exe' or not p.is_file() or p.is_symlink() or type(item['bytes']) is not int or not 1<=item['bytes']<=1073741824 or p.stat().st_size!=item['bytes'] or not re.fullmatch('[0-9a-f]{64}',item['sha256']) or sha(p)!=item['sha256']:raise ValueError('REFERENCE_RUNTIME_IDENTITY_REJECTED')
  bins[name]=p
 consumer=Path(args.consumer).resolve(strict=True);migrations=sorted((consumer/'db/migrations').glob('*.up.sql'))
 if not migrations or {p.name:sha(p) for p in migrations}!=v['migrations'] or sha(S/'seed.sql')!=v['seed_sha256']:raise ValueError('REFERENCE_FIXTURE_IDENTITY_REJECTED')
 root=Path(args.work_root).resolve(strict=True)
 if not root.is_dir() or len(str(root))>145:raise ValueError('REFERENCE_WORK_ROOT_REJECTED')
 return bins,migrations,root
try:BINS,MIGRATIONS,ROOT=load_inputs()
except (ValueError,OSError,KeyError,TypeError,rs.g.Rejected):
 print('REFERENCE_INPUT_REJECTED',file=sys.stderr);raise SystemExit(2)
W=ROOT/('integrated-'+uuid.uuid4().hex);W.mkdir()
# Restrict only the newly created reference directory. The ephemeral CA key is
# never written; client/server test keys inherit current-user and SYSTEM access.
system=Path(os.environ['SYSTEMROOT'])/'System32'
identity=subprocess.run([str(system/'whoami.exe'),'/user','/fo','csv','/nh'],capture_output=True,text=True,timeout=5,creationflags=subprocess.CREATE_NO_WINDOW,check=True)
sid=list(csv.reader(identity.stdout.splitlines()))[0][1]
if not re.fullmatch(r'S-[0-9-]+',sid):raise SystemExit('REFERENCE_IDENTITY_REJECTED')
acl=subprocess.run([str(system/'icacls.exe'),str(W),'/inheritance:r','/grant:r','*'+sid+':(OI)(CI)F','*S-1-5-18:(OI)(CI)F'],capture_output=True,timeout=5,creationflags=subprocess.CREATE_NO_WINDOW)
if acl.returncode:raise SystemExit('REFERENCE_DIRECTORY_ACL_REJECTED')
print('REFERENCE_STAGE '+str(W),flush=True)
def port():
 with socket.socket() as s:s.bind(('127.0.0.1',0));return s.getsockname()[1]
ports={n:port() for n in ['pg','otlp','metrics','health','prom']};procs={};logs=[];receipts={};host_threads=[]
env={k:os.environ[k] for k in ['SYSTEMROOT','WINDIR']};env.update(TEMP=str(W),TMP=str(W),PGCONNECT_TIMEOUT='3',COMSPEC=str(system/'cmd.exe'))
def command(args,name,timeout=60):
 p=subprocess.run([str(x) for x in args],env=env,stdin=subprocess.DEVNULL,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,creationflags=subprocess.CREATE_NO_WINDOW,timeout=timeout);(W/(name+'.log')).write_bytes(p.stdout)
 if p.returncode:raise RuntimeError(name+' failed: '+p.stdout.decode(errors='replace')[-2000:])
 return p.stdout.decode().strip()
def start(name,args):
 f=(W/(name+'-'+str(len(logs))+'.log')).open('xb');logs.append(f);p=subprocess.Popen([str(a) for a in args],cwd=W,env=env,stdin=subprocess.DEVNULL,stdout=f,stderr=subprocess.STDOUT,creationflags=subprocess.CREATE_NO_WINDOW);procs[name]=p;return p

def wait(pred,label,seconds=20):
 end=time.monotonic()+seconds
 while time.monotonic()<end:
  try:
   value=pred()
   if value:return value
  except (OSError,ValueError,urllib.error.URLError):pass
  time.sleep(.1)
 raise TimeoutError(label)
def sql(query,db=None):
 if db is not None and db!='postgres' and not re.fullmatch(r'elite_refund_test_[0-9a-f]{32}',db):raise ValueError('REFERENCE_DB_REJECTED')
 return command([BINS['psql'],'-X','-w','-qAt','-h','127.0.0.1','-p',ports['pg'],'-U','postgres','-d',db or DB,'-v','ON_ERROR_STOP=1','-c',query],'sql-'+uuid.uuid4().hex,15)
def certs():
 command([BINS['telemetry-reference'],'--certificates',W],'reference-certificates')

def context(role='client'):
 c=ssl.create_default_context(cafile=str(W/'ca.pem'));c.minimum_version=ssl.TLSVersion.TLSv1_3
 if role:c.load_cert_chain(str(W/(role+'.pem')),str(W/(role+'.key')))
 return c

def get(endpoint,path,role='client'):
 with urllib.request.urlopen('https://127.0.0.1:'+str(ports[endpoint])+path,context=context(role),timeout=2) as r:return r.read()

def signals():
 result={'logs':[],'metrics':[],'traces':[]}
 for p in W.glob('signals*.json'):
  for line in p.read_bytes().splitlines():
   j=json.loads(line)
   for k,t in [('resourceLogs','logs'),('resourceMetrics','metrics'),('resourceSpans','traces')]:
    if k in j:result[t].extend(j[k])
 return result

def attributes(a):return {x['key']:next(iter(x['value'].values())) for x in a}
def log_records():return [(attributes(r['resource']['attributes']),x) for r in signals()['logs'] for s in r['scopeLogs'] for x in s['logRecords']]
def outcomes():return [attributes(x['attributes'])['outcome'] for _,x in log_records()]
def api(path):return json.loads(get('prom',path))
def alerts():return api('/api/v1/alerts')['data']['alerts']
def firing(name):return [a for a in alerts() if a['labels']['alertname']==name and a['state']=='firing']
def launch_host(label,db,extra=None):
 root=W/('host-'+label);root.mkdir();store=root/'store';store.mkdir();exe=BINS['refund-host'];run=uuid.uuid4().hex
 e={'SYSTEMROOT':os.environ['SYSTEMROOT'],'WINDIR':os.environ['WINDIR'],'TEMP':str(root),'TMP':str(root),'DATABASE_URL':f'postgres://postgres@127.0.0.1:{ports["pg"]}/{db}?sslmode=disable','REFUND_WORKER_ID':'elite-reference-'+label,'REFUND_TELEMETRY_CONFIG':str(W/'reporter.json'),'REFUND_TELEMETRY_CONFIG_SHA256':sha(W/'reporter.json')};e.update(extra or {})
 profile={'schema':'elite-native-launch-qualification/v3','id':'reference-refund-host','executable':{'path':str(exe),'bytes':exe.stat().st_size,'sha256':sha(exe)},'arguments':[],'cwd':str(root),'environment':e,'budgets':{'timeout':120,'output_bytes':8192,'processes':8,'commit_bytes':536870912},'shutdown':{'protocol':'win32-inherited-event/v1','grace_seconds':5}}
 raw=rs.encode(profile);(root/'profile.json').write_bytes(raw);digest=rs.digest(raw);cancel=threading.Event();result={}
 def worker():
  try:result['receipt']=rs.execute(store,run,raw,digest,cancel=cancel)
  except BaseException as exc:result['error']=repr(exc)
 th=threading.Thread(target=worker);th.start();host_threads.append((th,cancel));return root,store,run,digest,cancel,th,result

def finish_host(h,stop=True):
 root,store,run,digest,cancel,th,result=h
 if stop:cancel.set()
 th.join(12);assert not th.is_alive(),result
 assert 'error' not in result,result
 receipt=result['receipt'];assert receipt.state=='RECORDED',receipt
 assert rs.inspect(store,run,digest)==receipt
 assert rs.execute(store,run,(root/'profile.json').read_bytes(),digest).state=='ALREADY_RESERVED'
 data=json.loads((store/run/'result.json').read_text());receipts[root.name]=data;return data['outcome']

def main():
 global DB
 certs()
 tls={'cert_file':str(W/'server.pem'),'key_file':str(W/'server.key'),'client_ca_file':str(W/'ca.pem'),'min_version':'1.3'}
 config={'receivers':{'otlp':{'protocols':{'http':{'endpoint':'127.0.0.1:'+str(ports['otlp']),'tls':tls,'max_request_body_size':32768,'read_timeout':'2s','write_timeout':'2s','idle_timeout':'2s'}}}},'exporters':{'file':{'path':str(W/'signals.json'),'format':'json','rotation':{'max_megabytes':1,'max_days':1,'max_backups':2}},'prometheus':{'endpoint':'127.0.0.1:'+str(ports['metrics']),'tls':tls,'metric_expiration':'5s','resource_constant_labels':{'included':['service.name','deployment.environment.name','service.instance.id']}}},'processors':{'memory_limiter':{'check_interval':'1s','limit_mib':256,'spike_limit_mib':32}},'extensions':{'health_check':{'endpoint':'127.0.0.1:'+str(ports['health'])}},'service':{'extensions':['health_check'],'telemetry':{'logs':{'level':'error'},'metrics':{'level':'none'}},'pipelines':{'logs':{'receivers':['otlp'],'processors':['memory_limiter'],'exporters':['file']},'traces':{'receivers':['otlp'],'processors':['memory_limiter'],'exporters':['file']},'metrics':{'receivers':['otlp'],'processors':['memory_limiter'],'exporters':['file','prometheus']}}}}
 (W/'collector.json').write_text(json.dumps(config,indent=2))
 pc={'global':{'scrape_interval':'1s','evaluation_interval':'1s'},'rule_files':[str(W/'rules.json')],'scrape_configs':[{'job_name':'reference-refund','scheme':'https','tls_config':{'ca_file':str(W/'ca.pem'),'cert_file':str(W/'client.pem'),'key_file':str(W/'client.key'),'min_version':'TLS13'},'static_configs':[{'targets':['127.0.0.1:'+str(ports['metrics'])]}]}]}
 rules={'groups':[{'name':'reference-refund','interval':'1s','rules':[{'alert':'ReferenceCollectorUnavailable','expr':'up{job="reference-refund"} == 0','for':'2s','annotations':{'summary':'El colector no responde; revisar el proceso y la identidad TLS.'}},{'alert':'ReferenceRefundStepError','expr':'sum(increase(elite_refund_observations_total{outcome="error"}[10s])) > 0','for':'1s','annotations':{'summary':'El worker informó un fallo; consultar el resultado durable antes de reintentar.'}},{'alert':'ReferenceRefundHeartbeatMissing','expr':'absent(elite_refund_observations_total{outcome="idle"})','for':'8s','annotations':{'summary':'No hay señales recientes del worker.'}}]}]}
 for n,v in [('prometheus',pc),('rules',rules),('web',{'tls_server_config':{'cert_file':str(W/'server.pem'),'key_file':str(W/'server.key'),'client_ca_file':str(W/'ca.pem'),'client_auth_type':'RequireAndVerifyClientCert','min_version':'TLS13'}})]: (W/(n+'.json')).write_text(json.dumps(v,indent=2))
 reporter={'endpoint':'https://127.0.0.1:'+str(ports['otlp']),'timeout_ms':1500}
 for k,n in [('ca','ca.pem'),('certificate','client.pem'),('private_key','client.key')]:reporter[k+'_file']=str(W/n);reporter[k+'_sha256']=sha(W/n)
 (W/'reporter.json').write_text(json.dumps(reporter))
 command([BINS['elite-otelcol'],'validate','--config',W/'collector.json'],'collector-validate')
 command([BINS['promtool'],'check','config',W/'prometheus.json'],'promtool-config')
 start('collector',[BINS['elite-otelcol'],'--config',W/'collector.json']);wait(lambda:(get('metrics','/metrics'),True),'collector metrics')
 for role in [None,'rogue']:
  try:get('metrics','/metrics',role);raise AssertionError('unauthenticated metrics accepted')
  except urllib.error.HTTPError:raise AssertionError('unauthenticated TLS reached HTTP')
  except (ssl.SSLError,urllib.error.URLError):pass
 start('prometheus',[BINS['prometheus'],'--config.file='+str(W/'prometheus.json'),'--web.listen-address=127.0.0.1:'+str(ports['prom']),'--web.config.file='+str(W/'web.json'),'--storage.tsdb.path='+str(W/'tsdb'),'--storage.tsdb.retention.time=1h','--storage.tsdb.retention.size=64MB','--log.level=error'])
 wait(lambda:api('/api/v1/status/config')['status']=='success','prometheus API')
 for endpoint,path in [('otlp','/v1/logs'),('prom','/api/v1/status/config')]:
  for role in [None,'rogue']:
   try:get(endpoint,path,role);raise AssertionError('unauthenticated TLS accepted')
   except urllib.error.HTTPError:raise AssertionError('unauthenticated TLS reached HTTP')
   except (ssl.SSLError,urllib.error.URLError):pass
 DATA=W/'pgdata';command([BINS['initdb'],'-D',DATA,'-U','postgres','-A','trust','--no-locale','-E','UTF8'],'initdb')
 with (DATA/'postgresql.conf').open('a') as f:f.write(f"\nlisten_addresses='127.0.0.1'\nport={ports['pg']}\nfsync=on\nsynchronous_commit=on\nfull_page_writes=on\n")
 start('postgres',[BINS['postgres'],'-D',DATA])
 def owned_database():
  if procs['postgres'].poll() is not None:raise RuntimeError('OWNED_POSTGRES_START_FAILED')
  observed=sql('show data_directory','postgres')
  if Path(observed).resolve()!=DATA.resolve():raise RuntimeError('OWNED_POSTGRES_IDENTITY_MISMATCH')
  return True
 wait(owned_database,'owned postgres startup')
 DB='elite_refund_test_'+uuid.uuid4().hex;sql('create database '+DB,'postgres')
 migrations=MIGRATIONS
 for p in migrations:command([BINS['psql'],'-X','-w','-qAt','-h','127.0.0.1','-p',ports['pg'],'-U','postgres','-d',DB,'-v','ON_ERROR_STOP=1','-f',p],'migration-'+p.stem)
 command([BINS['psql'],'-X','-w','-qAt','-h','127.0.0.1','-p',ports['pg'],'-U','postgres','-d',DB,'-v','ON_ERROR_STOP=1','-f',S/'seed.sql'],'seed')
 bad=launch_host('config-tamper',DB,{'REFUND_TELEMETRY_CONFIG_SHA256':'0'*64});rejected=finish_host(bad,False)
 assert rejected['exit_code']==1 and rejected['tree_empty'] and not log_records()
 assert sql("select status||':'||attempt_count from sales.return_effect_execution where request_id='refund'")=='requested:0'
 h=launch_host('healthy',DB);wait(lambda:'handled' in outcomes() and outcomes().count('idle')>=4,'actual main signals')
 assert sql("select status from sales.return_effect_execution where request_id='refund'")=='blocked'
 assert sql("select last_error_code from sales.return_effect_execution where request_id='refund'")=='PROVIDER_CONFIG_MISSING'
 assert sql('select count(*) from payment.return_refund_observation')=='0'
 metric=get('metrics','/metrics');(W/'metrics-healthy.txt').write_bytes(metric);assert b'elite_refund_observations_total' in metric
 wait(lambda:api('/api/v1/query?'+urllib.parse.urlencode({'query':'elite_refund_observations_total'}))['data']['result'],'scraped host counters')
 o=finish_host(h);assert o['exit_code']==0 and o['shutdown_state']=='EXITED_DURING_GRACE',o
 wait(lambda:'stopped' in outcomes(),'graceful stop exported')
 print('ACTUAL_MAIN_HEALTHY_PASS',flush=True)
 # Actual SQL permission failure, without exposing its sensitive diagnostic.
 sql("alter table sales.return_effect_execution rename to private_customer_diagnostic")
 h=launch_host('db-error',DB);wait(lambda:outcomes().count('error')>=4,'database error observed');wait(lambda:firing('ReferenceRefundStepError'),'real error alert')
 (W/'alerts-error.json').write_text(json.dumps(alerts(),indent=2));finish_host(h);sql("alter table sales.private_customer_diagnostic rename to return_effect_execution")
 print('REAL_ERROR_ALERT_PASS',flush=True)
 h=launch_host('sink-loss',DB);count=len(log_records());wait(lambda:len(log_records())>=count+5,'sink-loss baseline')
 procs['collector'].terminate();procs['collector'].wait(10)
 o=finish_host(h,False);assert o['exit_code']==1 and o['tree_empty'],o
 wait(lambda:firing('ReferenceCollectorUnavailable'),'real sink health alert');(W/'alerts-sink-down.json').write_text(json.dumps(alerts(),indent=2))
 before={p.name:sha(p) for p in W.glob('signals*.json')};assert before
 start('collector',[BINS['elite-otelcol'],'--config',W/'collector.json']);wait(lambda:(get('metrics','/metrics'),True),'collector restart')
 assert {p.name:sha(p) for p in W.glob('signals*.json')}==before,'restart truncated retained signals'
 h=launch_host('recovered',DB);count=len(log_records());wait(lambda:len(log_records())>=count+6,'new instance after recovery');wait(lambda:not firing('ReferenceCollectorUnavailable'),'sink alert recovery');finish_host(h)
 print('SINK_FAILURE_RECOVERY_PASS',flush=True)
 # Correlation is checked from actual Collector OTLP output, not sender intent.
 sig=signals();lr=log_records();tr=[(attributes(r['resource']['attributes']),x) for r in sig['traces'] for s in r['scopeSpans'] for x in s['spans']]
 traces={(r['service.instance.id'],x['traceId'],x['spanId']) for r,x in tr}
 assert all((r['service.instance.id'],x['traceId'],x['spanId']) in traces for r,x in lr),('missing correlated spans',len(lr),len(tr))
 instances={r['service.instance.id'] for r,x in lr};assert len(instances)==4
 for instance in instances:
  seq=[int(attributes(x['attributes'])['sequence']) for r,x in lr if r['service.instance.id']==instance];assert seq==list(range(1,len(seq)+1))
 raw=b'\n'.join(p.read_bytes() for p in W.glob('signals*.json'))+metric
 for secret in [b'private_customer_diagnostic',b'postgres://',DB.encode(),b'PRIVATE KEY',b'PROVIDER_CONFIG_MISSING']:assert secret not in raw,secret
 assert all(r['deployment.environment.name']=='reference' for r,x in lr)
 # Verify actual persisted cumulative counters/histograms per instance/outcome.
 previous={};points=0
 for resource in sig['metrics']:
  identity=attributes(resource['resource']['attributes'])['service.instance.id']
  for scope in resource['scopeMetrics']:
   for m in scope['metrics']:
    kind='sum' if 'sum' in m else 'histogram';body=m[kind];assert body['aggregationTemporality']==2
    if kind=='sum':assert body['isMonotonic'] is True
    for point in body['dataPoints']:
     key=(identity,m['name'],attributes(point['attributes'])['outcome']);count=int(point['asInt'] if kind=='sum' else point.get('count','0'))
     series_started=point['startTimeUnixNano'];before=previous.get(key,(0,series_started));assert count>=before[0] and series_started==before[1]
     if kind=='histogram':assert sum(map(int,point['bucketCounts']))==count and point['sum']>=0 and point['explicitBounds']==[.001,.01,.1,1,10,60]
     previous[key]=(count,series_started);points+=1
 (W/'operational-signals.jsonl').write_bytes(b'\n'.join(p.read_bytes() for p in W.glob('signals*.json')))
 receipts['operational_evidence']={'sha256':sha(W/'operational-signals.jsonl'),'points_checked':points,'series_checked':len(previous)}
 # CLI executes the same closed-outcome reporter, bounded to2000observations.
 # The synthetic expired backup tests the configured age policy on real rotation.
 expired=W/'signals-2000-01-01T00-00-00.000-size.json';expired.write_text('{"resourceLogs":[]}\n')
 began=time.monotonic();load=json.loads(command([BINS['telemetry-reference'],W/'reporter.json',sha(W/'reporter.json'),'2000','8'],'reference-load',130));assert load['observations']==2000 and load['instances']==8
 wait(lambda:len(list(W.glob('signals*.json')))<=3 and not expired.exists(),'file rotation/age retention')
 retained=list(W.glob('signals*.json'));assert len(retained)==3 and all(p.stat().st_size<=1048576 for p in retained)
 for p in retained:
  for line in p.read_bytes().splitlines():json.loads(line)
 final_metrics=get('metrics','/metrics');(W/'metrics-after-load.txt').write_bytes(final_metrics)
 values=[line for line in final_metrics.decode().splitlines() if line.startswith('elite_refund_observations_total{') and 'outcome="idle"' in line and line.endswith(' 250')];assert len(values)==8,values
 file_snapshot={p.name:sha(p) for p in retained}
 procs['collector'].terminate();procs['collector'].wait(10);start('collector',[BINS['elite-otelcol'],'--config',W/'collector.json']);wait(lambda:(get('metrics','/metrics'),True),'rotated collector restart')
 assert {p.name:sha(p) for p in W.glob('signals*.json')}==file_snapshot
 receipts['load_retention']={'observations':2000,'instances':8,'seconds':time.monotonic()-began,'files':{p.name:{'bytes':p.stat().st_size,'sha256':sha(p)} for p in retained},'age_fixture_removed':not expired.exists(),'restart_identical':True}
 print('CLI_LOAD_ROTATION_RETENTION_PASS',flush=True)
 return {'pass':True,'actual_main':True,'instances':len(instances),'log_records':len(lr),'trace_records':len(tr),'metric_records':len(sig['metrics']),'metric_points_checked':points,'cases':['mutual-TLS-negative','main-PostgreSQL-blocked-provider','three-signals','counter-scrape','real-alert','graceful-stop','sink-loss-exit1','no-replay','collector-restart-retention','alert-recovery','correlation','privacy'],'ports':ports,'runtime_lock_sha256':args.runtime_lock_sha256,'binary_sha256':{n:sha(BINS[n]) for n in ['elite-otelcol','prometheus','promtool','refund-host']},'process_receipts':receipts}
result={}
try:result=main()
except BaseException as e:traceback.print_exc();result={'pass':False,'error_type':type(e).__name__,'detail':str(e),'receipts':receipts}
finally:
 for th,c in host_threads:c.set();th.join(12)
 for n,p in procs.items():
  if p.poll() is None:
   if n=='postgres':
    try:command([BINS['pg_ctl'],'-D',W/'pgdata','-m','fast','-w','-t','40','stop'],'postgres-stop',45)
    except Exception as e:result['cleanup_error']=str(e);result['pass']=False
   else:p.terminate()
   try:p.wait(12)
   except subprocess.TimeoutExpired:p.kill();p.wait(5);result['forced_cleanup']=n;result['pass']=False
 for f in logs:f.close()
 result['all_owned_processes_stopped']=all(p.poll() is not None for p in procs.values());result['stage']=str(W)
 (W/'result.json').write_text(json.dumps(result,indent=2)+'\n');print(json.dumps({k:v for k,v in result.items() if k not in ['process_receipts','receipts']}),flush=True)
raise SystemExit(0 if result.get('pass') else 1)
````

### FILE: `reference_telemetry/runtime-lock.example.json`
```yaml
block_id: "WINDOWS-REFERENCE-TELEMETRY:runtime-lock.example:v1"
operation: CREATE
provenance: AUTHORED
source: "local V345-V351 native qualification and V372 integrated reference; official contracts govern behavior, no vendor implementation authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "455a9992d6a7fe36ba844973ad99511895800776d80341fca21736203ed0b322"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-reference-telemetry-runtime/v1",
  "binaries": {
    "elite-otelcol": {
      "path": "REPLACE_WITH_VERIFIED_ABSOLUTE_EXE_PATH",
      "bytes": 0,
      "sha256": "REPLACE_WITH_VERIFIED_SHA256"
    },
    "prometheus": {
      "path": "REPLACE_WITH_VERIFIED_ABSOLUTE_EXE_PATH",
      "bytes": 0,
      "sha256": "REPLACE_WITH_VERIFIED_SHA256"
    },
    "promtool": {
      "path": "REPLACE_WITH_VERIFIED_ABSOLUTE_EXE_PATH",
      "bytes": 0,
      "sha256": "REPLACE_WITH_VERIFIED_SHA256"
    },
    "refund-host": {
      "path": "REPLACE_WITH_VERIFIED_ABSOLUTE_EXE_PATH",
      "bytes": 0,
      "sha256": "REPLACE_WITH_VERIFIED_SHA256"
    },
    "telemetry-reference": {
      "path": "REPLACE_WITH_VERIFIED_ABSOLUTE_EXE_PATH",
      "bytes": 0,
      "sha256": "REPLACE_WITH_VERIFIED_SHA256"
    },
    "postgres": {
      "path": "REPLACE_WITH_VERIFIED_ABSOLUTE_EXE_PATH",
      "bytes": 0,
      "sha256": "REPLACE_WITH_VERIFIED_SHA256"
    },
    "initdb": {
      "path": "REPLACE_WITH_VERIFIED_ABSOLUTE_EXE_PATH",
      "bytes": 0,
      "sha256": "REPLACE_WITH_VERIFIED_SHA256"
    },
    "pg_ctl": {
      "path": "REPLACE_WITH_VERIFIED_ABSOLUTE_EXE_PATH",
      "bytes": 0,
      "sha256": "REPLACE_WITH_VERIFIED_SHA256"
    },
    "psql": {
      "path": "REPLACE_WITH_VERIFIED_ABSOLUTE_EXE_PATH",
      "bytes": 0,
      "sha256": "REPLACE_WITH_VERIFIED_SHA256"
    }
  },
  "migrations": {},
  "seed_sha256": "6fa405b703985ccbde52f1d3f37f130826d10d0056afbcccd667b0a557792b9a"
}
````

### FILE: `reference_telemetry/seed.sql`
```yaml
block_id: "WINDOWS-REFERENCE-TELEMETRY:seed:v1"
operation: CREATE
provenance: AUTHORED
source: "local V345-V351 native qualification and V372 integrated reference; official contracts govern behavior, no vendor implementation authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6fa405b703985ccbde52f1d3f37f130826d10d0056afbcccd667b0a557792b9a"
variables: []
secrets_allowed: false
```
````sql
-- AUTHORED local fixture. Pre-existing captured order / authorized return are
-- synthetic starting conditions, NOT an originating sales/approval workflow.
-- All actual constraints/triggers remain enabled. Fresh database only.
begin;
insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','fixture351','Synthetic','Synthetic');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','store','store','Synthetic','store');
insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','model','model','Synthetic','bicycle','active');
insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','variant','model','variant','Synthetic','{}','active');
insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','customer','Synthetic','synthetic@example.test');
insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','stock','store','variant','SYNTHETIC-SERIAL','sold',1,clock_timestamp());
insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','order','store','customer','delivered','ARS',1000,1);
insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units,allocated_stock_unit_id)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','order','line','variant',1,1000,'stock');
insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','payment','order','stripe','synthetic-pi','fixture-payment-00000001','captured','ARS',1000,3);
insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','handover','store','order','customer','stock','prepared',1);
insert into sales.delivery_exception(tenant_id,exception_id,organization_id,handover_id,customer_principal_id,reason_code,details,rejection_evidence_sha256_hex,state,version)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','exception','store','handover','customer','synthetic-return','Synthetic',repeat('a',64),'open',1);
insert into sales.return_authorization(tenant_id,authorization_id,organization_id,exception_id,handover_id,order_id,stock_unit_id,customer_principal_id,disposition,exact_cost_source_order_id,state,authorized_by_subject)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','authorization','store','exception','handover','order','stock','customer','return','order','authorized','synthetic-operator');
insert into sales.return_receipt(tenant_id,receipt_id,authorization_id,organization_id,order_id,stock_unit_id,customer_principal_id,received_serial_number,condition_code,notes,evidence_sha256_hex,received_by_subject)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','receipt','authorization','store','order','stock','customer','SYNTHETIC-SERIAL','opened','Synthetic',repeat('a',64),'synthetic-operator');
insert into sales.return_disposition(tenant_id,disposition_id,receipt_id,inventory_action,customer_remedy,notes,decided_by_subject)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','disposition','receipt','restock','refund','Synthetic','synthetic-operator');
insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','refund','disposition','refund','payment','requested','synthetic-refund-0000000001');
create schema qualification;
create table qualification.effect(run_id text primary key,idempotency_key text not null unique);
commit;
````

### FILE: `reference_telemetry/test_capture.py`
```yaml
block_id: "WINDOWS-REFERENCE-TELEMETRY:test-capture:v1"
operation: CREATE
provenance: AUTHORED
source: "local V345-V351 native qualification and V372 integrated reference; official contracts govern behavior, no vendor implementation authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8a78f1c54d788edb18e5d97a1cb2f4ca8c490af5adb0db0aaae2bad28873bf38"
variables: []
secrets_allowed: false
```
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
````

### FILE: `reference_telemetry/test_governed_launch.py`
```yaml
block_id: "WINDOWS-REFERENCE-TELEMETRY:test-governed-launch:v1"
operation: CREATE
provenance: AUTHORED
source: "local V345-V351 native qualification and V372 integrated reference; official contracts govern behavior, no vendor implementation authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9aa326b779065c254e64691935800fb665b53c92dacc7f6e624cd242f82015c5"
variables: []
secrets_allowed: false
```
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
        pwsh='<USERPROFILE>/.cache/codex-runtimes/codex-primary-runtime/dependencies/native/powershell/pwsh.exe'
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

### FILE: `reference_telemetry/test_reference_profile.py`
```yaml
block_id: "WINDOWS-REFERENCE-TELEMETRY:test-reference-profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local V345-V351 native qualification and V372 integrated reference; official contracts govern behavior, no vendor implementation authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "fc6501c0260bda4c8fce417c8f7c556408603d86849d72e313921781bb689148"
variables: []
secrets_allowed: false
```
````python
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
````

### FILE: `reference_telemetry/test_result_store.py`
```yaml
block_id: "WINDOWS-REFERENCE-TELEMETRY:test-result-store:v1"
operation: CREATE
provenance: AUTHORED
source: "local V345-V351 native qualification and V372 integrated reference; official contracts govern behavior, no vendor implementation authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "654a463ab0a0d39c8573b4e37379487ad5834e744b70239281c921f1266e3f63"
variables: []
secrets_allowed: false
```
````python
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
        pwsh='<USERPROFILE>/.cache/codex-runtimes/codex-primary-runtime/dependencies/native/powershell/pwsh.exe'
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
````

### FILE: `reference_telemetry/test_shutdown.py`
```yaml
block_id: "WINDOWS-REFERENCE-TELEMETRY:test-shutdown:v1"
operation: CREATE
provenance: AUTHORED
source: "local V345-V351 native qualification and V372 integrated reference; official contracts govern behavior, no vendor implementation authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6abe6e0507aea4d6e40ad2aa4228fc0cedb622cbdfefb074bd8e672186aedb79"
variables: []
secrets_allowed: false
```
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

## 6. Configuration surface

No templated secrets. Runtime lock and separately trusted SHA-256, reconstructed
consumer and owned work directory are explicit arguments. The example lock is
intentionally rejected. Private ephemeral keys never enter evidence exports.
TLS certificates and static reference limits are generated only after all locked
binaries, migrations and synthetic seed bytes pass verification.

## 7. Dependency bill

Python3.14 standard library on trusted Windows x64; no pip dependency. Native
Win32 Job Object/handle APIs use ctypes. External executables are separately
qualified Go1.26.8 builds of the compatible refund0.1.5 commands, minimal official
Collector0.160.0 and Prometheus/promtool3.14.0, plus PostgreSQL18.6 utilities.
Their exact binaries/versions/source licenses/Go module locks and SCA disposition
are in V372; they are not bundled in this AUTHORED pack. Any changed dependency
requires requalification and a new independently hashed runtime lock.

## 8. Apply order

Reconstruct into an absent directory; compose the compatible consumer separately.
Acquire/qualify/build exact external runtime candidates and both Go commands.
Run policy/native tests, create the exact local runtime lock and independently
verify its digest, then execute the finite runner with a new owned work root.
The runner creates only its own synthetic cluster; never point it at an existing
database. Preserve result/logs and inspect owned-process cleanup before promotion.

## 9. Verification

Materialize into an absent destination and verify all13file hashes. Execute the
six reference-policy tests plus23capture,20launch,19shutdown and28result-store
regressions. Build both Go entrypoints from the compatible worker pack, validate
Collector/Prometheus config, and run the finite integrated acceptance. A passed
local reference is evidence for TEST-05, not every project or every SRE capability.
Use reconstruction_evidence/INTEGRATED_TELEMETRY_CONTROL_V372.md for exact build,
source, SCA/package-exclusion reasoning, original failures and observed receipts.

### Failure modes and recovery

Input/hash mismatch fails before a work directory or child process. Runtime,
TLS, config, database identity, signal, counter, correlation, alert, retention
or process-cleanup failure leaves a failed receipt and retained source/logs.
No loop silently retries a failed host attempt. Use the existing reservation and
database observations to reconcile uncertainty; do not infer business success.

## 10. Reconstruction evidence

See reconstruction_evidence/INTEGRATED_TELEMETRY_CONTROL_V372.md: exact13file
reconstruction,96native/policy tests, actual integrated runtime and preserved
source/build/security/license/failure receipts. Full library structural gate
checks all mandatory contract sections and hashes; runtime PASS alone is not it.

CONDITIONED on the exact trusted Windows reference scope, verified input/runtime
identities and all observed gates. Product service identity, secret management,
regulatory retention, real alert routing, sustained workload, monitoring runtime,
security/offensive testing, deploy/rollback/DR and business acceptance remain
consumer responsibilities. No customer corpus or training action is included.
