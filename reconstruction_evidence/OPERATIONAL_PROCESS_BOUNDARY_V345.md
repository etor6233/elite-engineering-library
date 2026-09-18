# V345 — finite operational process boundary

Library maintenance / T2809. SECURE-OPS-DELIVERY-CORE1.1.2 →1.1.3;
37materializable files, three existing AUTHORED blocks changed. No new runtime,
SDK, service manager, provider activation or business behavior. Conditioned
admission remains; the overall roadmap and production are not complete.

## Defect and demonstrated correction

The exact previous runner used subprocess.run(capture_output=True), then checked
output length. A real child exceeded a1024byte test cap and performed a later
synthetic marker write before rejection. A second real child timed out and the
chained TimeoutExpired traceback exposed a synthetic private argv marker.
Both unchanged-behavior baseline probes fail and then pass after the repair.

The existing runner now reads both raw pipes fairly and nonblocking: one read
per stream/pass, at most64KiB/read and10MiB retained per stream. It checks a
monotonic deadline until process exit AND both EOFs. On failure it closes raw
pipes, kills/reaps only the direct child with a five-second cleanup budget and
returns a static error without private exception causes. No reader thread,
communicate retry, inherited stdin or restart is introduced. Windows launches
are hidden. Output at the exact cap remains valid; a single extra byte rejects.
The buffer bound is not a hard process-RSS limit and cannot prevent a child effect
that occurred before detection. OS creation/scheduling cannot be bounded here.

A finite descendant deliberately keeps the pipe open after its parent exits.
The runner rejects on its own deadline; the descendant exits naturally after
2.5seconds. This explicitly does not prove tree termination. Successful raw
stdout/stderr remain potentially sensitive evidence: ACL, review, retention and
storage quota must be supplied by the target. Receipt schema stays unchanged.

## Executed evidence

- Canonical production-admission suite: 110tests,
  109PASS/1SKIP. The skip is the existing OS symlink-creation
  permission case; it is not a new omission or proof of symlink handling.
  The candidate also passed the same suite before canonical integration.
- 14runner tests pass in candidate, fresh37file rebuild and full67pack/746file
  consumer. Includes timeout/private argv, both binary streams at1MiB, each stream
  limit+1, post-limit effect prevention, closed streams/live child, inherited pipe,
  DEVNULL stdin, environment exclusion, boolean timeout and injected read error.
  The direct child is reaped; captured pipes close; no reader threads remain.
- Unchanged operational-readiness validator: 12tests PASS, zero skips.
- All37owner files match candidate, fresh rebuild and consumer. Exactly three
  change;34remain exact. Full consumer comparison derives746paths from the
  materialization record: three changed,743unchanged, including all Go/DB/provider
  and dependency locks. Generated composition receipts remain separate evidence.
- Actual previously compiled refund0.1.4 executable is run by the repaired runner
  without DATABASE_URL or provider environment: exit1 and REFUND_HOST_FAILED
  retained in hash-bound evidence. No synthetic production label in the existing
  receipt schema is interpreted as a real production target or acceptance.
- Capability-gap kit freshly composed: self-check3positives/6negatives PASS.
  Full library Preflight133 is still pending at this integration checkpoint.

## Current primary authority and service-supervisor gap

Observed 2026-09-08T23:42:47.602171Z. Internal authorities: SECURITY_SRE_CLOUD_INFRASTRUCTURE.md
§12.5/process-health distinction, ENGINEERING_EXECUTION_PLAYBOOK.md and
PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md G0–G8.

[Python subprocess](https://docs.python.org/3.14/library/subprocess.html) documents
memory collection by communicate and the OS-creation limit on timeouts.
[Python os](https://docs.python.org/3.14/library/os.html#os.set_blocking) documents
nonblocking pipe support on Windows from3.12. Current3.14 docs display3.14.7;
executed runtime remains the separately hashed3.14.4. No automatic upgrade,
new runtime admission, SCA or runtime-wide security approval is inferred.

[Microsoft Job Objects](https://learn.microsoft.com/en-us/windows/win32/procthread/job-objects)
defines owned process groups and their lifecycle limits; it is a relevant Windows
primitive, not an installed service manager or an admitted library wrapper.
[Official systemd service source](https://github.com/systemd/systemd/blob/main/man/systemd.service.xml)
describes watchdog/restart and rate limiting. This moving documentation branch is
research only; no code/license/release pin or Linux execution is admitted.
The freedesktop manual fetch failed in the browser; its official repository
provided the documentation fallback. No source absence is inferred from the fetch.

Existing-owner search across implementation_packs found contracts/reporter
interfaces, finite runners and target supervisor requirements, but no compatible
admitted durable service supervisor. Retain the operations owner CONDITIONED;
materialized gap decision must keep implementation_ready=false for the missing
service layer. Finite-runner defect repair is maintenance of existing code, not
implementation of that missing capability and not a workaround for its gate.

Candidate G0 provenance PASS (AUTHORED, exact canonical hash); G1 workspace-owner
license declaration PASS for local maintenance; G2 authority alignment PASS.
For the broad service-supervisor requirement G3–G8 FAIL: service deployment,
process-tree/startup races, bounded restart/heartbeat, privilege identity, collector,
retention/alerting, resource limits and target recovery are not demonstrated.
The narrow finite-runner correction has the code/tests described above but does
not upgrade those service gates. No claim of exhaustive Internet search is made.

Next action: qualify a native supervisor candidate in isolated reference fixtures,
starting with Windows Job Objects (owned launch/tree cleanup and parent-death),
then audit exact source/runtime/rights and G0–G8 before product incorporation.
Preserve POSIX/systemd as a separately conditioned target option, not a silently
selected deployment. No human question is required for these local experiments.
No automatic restart policy or command authorization for real effects is invented.

## Narrow implementation assurance

Correctness/integration: baseline red→green and actual child/binary execution.
Security/privacy: fixed error messages, selected environment, bounded reads;
not a sandbox or evidence redactor. Reliability: timeout through EOF, explicit
cleanup failure, no retry of uncertain effects. Performance: bounded retained
streams, fair reads; no load/SLO or RSS claim. Operability: exact exit/receipt
on completion and explicit failure otherwise; service/alerts are still missing.
Maintainability/provenance: three AUTHORED blocks in the existing owner, no new
dependency. Compatibility: existing schema/semantic validators plus743identical
consumer files. Verification: fresh reconstruction and full affected Python
suites. Delivery: three exact profile pins1.1.3; no release/deploy promotion.

Stage: <LOCALAPPDATA>/Temp/elite-v345-bb97f0a21710475989c6e37712e10a28

| Evidence | SHA-256 |
|---|---|
| baseline-pack.md | `ca141dad8b42cf29cca05936feae3d973cfa6a8dfd33d00f13f29a05102dc629` |
| red_probe.py | `5ac40d950998b52122b6ad41a0410205c6be0f18640bc55dcd3a282722f729c1` |
| red.log | `179db02cb7c487cc4c9d31e023ec81d2145dd0a0ac627d59014350b927f60d38` |
| green.log | `8d7c4c1d4bbf168dd7ae3810dc3bc4744b750f0643feaf43b5a7fed62f4061da` |
| candidate-tests.json | `368cfb4becd5eee94964502ded7a5ee8dace18597749a8b9c17e093b3f1d4e3d` |
| rebuilt-tests.json | `655e0b37f2b5873dfab0874439395899ace5cf93a8b64988323bc5cf8bb86318` |
| consumer-runner-tests.json | `f9578ad0a8ff75a5fc1a01e81db32cf1ba61678bbfe5acab098cd13ba92e1bf8` |
| readiness-tests.json | `8858f691d6eeb89f88f9c8d20dd0169634e40ce48cf0081105f05856160dcd18` |
| canonical-parity.json | `300b5c2f0e70b87a627246827003a9a611a04a2f0b308aa024a0b73dc094c281` |
| consumer-delta.json | `2c0345d43a577d41044ad95a7e129643c62e4d1194582331908db9a795cf3d39` |
| worker-reference.json | `a8ea51d1f39ec8976869aaf47a9a7937870edd834632b5cc1869c6c559867828` |
| runtime-receipt.json | `9b69ad0df365fe675c4786b99b56c72463a09e9793cca301587e119ec2832f3d` |
| changed-blocks.json | `69adafc292e133cc4c267d4cf6c064df90695df5db741121fdfee2d13badf41d` |
| changed-profiles.json | `eb33a90e05b0b702cf88dfb94ffdb6acb9b7f9bb35414f4571c541acc874baf6` |
| verify_parity-before593.py | `49199997dced33644f773b7138f44739973f73d3dd7b931be4f0adcd9601a983` |

## Native qualification continuation — preserved experiment, not admitted product

After the finite-runner rebuild, an AUTHORED Windows Job Object experiment
executed three scenarios three times: last-handle close terminates root and nested
child; invalid assignment leaves child suspended and terminates it without user
code; abrupt owner os._exit bypasses Python finally yet kernel closes the job and
terminates both children. It records observed process handles and membership,
never enumerates or kills unrelated processes. Windows build26200, CPython3.14.4;
ABI sizes144/104/24 for Extended/Startup/Process were observed. No real service,
provider or global installation. First nine-run receipt and source are retained.

The source below adds fail-path cleanup for the fixture itself and a fourth case:
a descendant attempting CREATE_BREAKAWAY_FROM_JOB must be refused with error5.
It must be reconstructed to a new isolated directory and executed before claiming
that fourth scenario PASS. This is a qualification harness, not a selected
materialization block; no pack/profile gains product code from it.

Microsoft sources observed2026-09-08:
[CreateJobObjectW](https://learn.microsoft.com/en-us/windows/win32/api/jobapi2/nf-jobapi2-createjobobjectw),
[AssignProcessToJobObject](https://learn.microsoft.com/en-us/windows/win32/api/jobapi2/nf-jobapi2-assignprocesstojobobject),
[limit structure](https://learn.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-jobobject_extended_limit_information),
[CreateProcessW](https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-createprocessw).
The adapter is authored against documented APIs; no Microsoft sample copied,
no claim that Microsoft authored the Python code. Job handles are unnamed and
non-inheritable, children begin suspended, assignment precedes resume, and the
job does not permit breakaway. This controls process membership and lifetime,
not filesystem/network privilege or all Windows service-manager behavior.

### Superseded qualification source (AUTHORED / REJECTED_FOR_LAUNCH_WINDOW)

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
def suspended(root,code):
 si=Startup();si.cb=C.sizeof(si);pi=Process()
 cmd=C.create_unicode_buffer(subprocess.list2cmdline([str(Path(sys.executable).resolve()),'-I','-c',code]))
 env={k:v for k,v in os.environ.items() if k.upper() in {'SYSTEMROOT','WINDIR','TEMP','TMP'}}
 block=C.create_unicode_buffer('\0'.join(k+'='+v for k,v in sorted(env.items()))+'\0\0')
 require(create(str(Path(sys.executable).resolve()),cmd,None,None,False,0x4|0x08000000|0x400,block,str(root),C.byref(si),C.byref(pi)),'CreateProcess')
 return pi
def tree_code():
 nested="import os,time,json,pathlib;pathlib.Path('nested.json').write_text(json.dumps({'pid':os.getpid()}));time.sleep(30)"
 return f"import subprocess,sys,time;subprocess.Popen([sys.executable,'-I','-c',{nested!r}],creationflags=0x08000000);time.sleep(30)"
def owner(root):
 h=job();pi=suspended(root,tree_code())
 try:
  require(assign(h,pi.process),'AssignProcess');require(resume(pi.thread)!=0xffffffff,'ResumeThread')
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
  pi=suspended(root,tree_code());require(assign(h,pi.process),'AssignProcess')
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
 pi=suspended(root,"from pathlib import Path;Path('must-not-run').write_text('bad')")
 try:
  require(not assign(None,pi.process),'invalid-job refusal')
  require(terminate(pi.process,99),'TerminateOwnedSuspended');require(wait(pi.process,5000)==0,'reap suspended')
  require(not (root/'must-not-run').exists(),'no effect before admitted assignment')
  return {'case':'assignment_refusal','no_child_effect':True,'reaped':True}
 finally:
  if wait(pi.process,0)==258:terminate(pi.process,99);wait(pi.process,5000)
  close(pi.thread);close(pi.process)
def breakaway_refusal(root):
 h=job();pi=None
 escaped="from pathlib import Path;Path('escaped').write_text('bad')"
 body="import subprocess,sys,json,pathlib\ntry:\n subprocess.Popen([sys.executable,'-I','-c',"+repr(escaped)+"],creationflags=0x01000000|0x08000000)\nexcept OSError as error:\n pathlib.Path('refused.json').write_text(json.dumps({'winerror':error.winerror}))"
 try:
  pi=suspended(root,body);require(assign(h,pi.process),'AssignProcess');require(resume(pi.thread)!=0xffffffff,'ResumeThread')
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
def main():
 if len(sys.argv)>1 and sys.argv[1]=='owner':owner(Path(sys.argv[2]));return
 root=Path(__file__).resolve().parent/'runs';root.mkdir()
 outcomes=[]
 for iteration in range(3):
  for method in (normal,assignment_refusal,parent_death,breakaway_refusal):
   dest=root/f'{iteration}-{method.__name__}';dest.mkdir();start=time.monotonic();result=method(dest);result['seconds']=time.monotonic()-start;result['iteration']=iteration;outcomes.append(result)
 result={'schema':'elite-native-job-qualification/v1','status':'PASS','product_admitted':False,'python':sys.version,'os':sys.getwindowsversion().build,'sizes':{'extended':C.sizeof(Extended),'startup':C.sizeof(Startup),'process':C.sizeof(Process)},'outcomes':outcomes,'script_sha256':hashlib.sha256(Path(__file__).read_bytes()).hexdigest()}
 write(Path(__file__).with_suffix('.receipt.json'),result);print(json.dumps(result,indent=2))
if __name__=='__main__':main()
````

## Native startup correction — FAIL594

The initial twelve observations above remain true but are insufficient for
admission. A deliberately interrupted owner between CreateProcess(suspended) and
AssignProcessToJobObject left a suspended child alive. The controller held the
exact owned process handle, observed it still alive, then terminated/reaped it;
no user code ran. This rejects the initial launch design for parent-death safety.
The prior source is preserved above as REJECTED_FOR_LAUNCH_WINDOW.

The successor uses STARTUPINFOEXW and PROC_THREAD_ATTRIBUTE_JOB_LIST so membership
is supplied inside CreateProcessW itself, still suspended and hidden. The mutable
attribute list and job-handle array live through creation and are released with
DeleteProcThreadAttributeList. There is no subsequent assignment step. Invalid job
handle refuses creation; no fallback to uncontained launch is permitted.

[Microsoft UpdateProcThreadAttribute](https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-updateprocthreadattribute)
documents JOB_LIST from Windows10/Server2016. Its value0x2000d is derived from
Microsoft's [WinBase.h definitions](https://raw.githubusercontent.com/microsoft/win32metadata/main/generation/WinSDK/RecompiledIdlHeaders/um/WinBase.h)
(attribute number13 plus INPUT0x20000), inspected as research only; the moving
header is not a source pin or redistributed upstream. A search also returned an
unofficial SDK mirror, which was not used as authority. The raw official header
supplied the constants. Extended startup ABI112bytes observed on this x64 host.

The successor passes five scenarios three times (15observations): whole-tree
last-handle cleanup, invalid-job create refusal, abrupt owner exit with descendants,
breakaway refusal, and owner death immediately after creation before resume.
The last scenario now observes the suspended child terminated by the kernel with
no user effect. This closes only the demonstrated launch-window defect in the
qualification apparatus. Standard service identity, stdout pipe inheritance,
heartbeat/restart/rate policy, multi-owner lifecycle, handle exhaustion, ACL/SCA,
collector/retention/alerts and target service acceptance are still unproven.
No supervisor code is admitted into product profiles by these experiments.

### Current qualification source (AUTHORED / NOT_ADMITTED)

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

## Closure134

Preflight133 completed: all154executed checks PASS, including structural
inventory162packs/1461materializable files/783Markdown
and53profiles. Availability remains BLOCKED only by Docker. Network,
provider, deployment and opt-in runtime skips are not promoted to executed gates.
No product code changed after this Preflight; the subsequent changes preserve
qualification evidence and the corrected experimental apparatus only.

The current native apparatus reconstructs byte-for-byte from this canonical
report and again passes all15observations. The initial
pre-assignment design is rejected and retained; FAIL594 is regression-proven
only for the successor launch-window correction. No target supervision claim.
The prior twelve passing cases remain historical, not a replacement for the
additional failure injection or a service gate PASS.

FAIL591 recovered;592regression-proven in the existing finite runner;
593recovered by exact output inventory;594regression-proven in the isolated
native experiment. All failed hypotheses and original evidence remain retained.
Next work: native candidate operational I/O/handle inheritance, shutdown and
resource/privilege limits, followed by complete source/runtime/compatibility and
G0–G8 qualification before product admission. No restart policy or provider effect
may be inferred from these process-lifecycle experiments. Other roadmap fronts
remain open and no global completion percentage is computed.

The successor gap record/receipt is preserved as supervisor-gap-record-final.json
and supervisor-gap-receipt-final.json in the same stage. The selected owner stays
CONDITIONED; the native experiment is an additional CONDITIONED candidate. Receipt
validation is structural and never upgrades incomplete service G3–G8 evidence.

| Closure evidence | SHA-256 |
|---|---|
| preflight133.log | `1d0760e0efae5e3a375ded80a403f439e8cf0a461434a6d32f61187d21f4246f` |
| preflight133.json | `7224d7586276cb9f370ef5afd08ccf29545c56fac2f1db59a899c46c9b9fcb60` |
| structural-closure.json | `9e046971db8d2fa4f4f5ee7c87da919917d6836672c72ead769ebbad38c8b326` |
| supervisor-gap-record.json | `1865e2c7801b17ca986ec5dbff8da0fca6e6f81426f2a031373939374e7a8d99` |
| supervisor-gap-receipt.json | `a9a19fd5a75ac9afef6ad226cd6f876d70fd7eb48c683dde0a80f3fb4a4b56e7` |
| native-qualification/job_probe.py | `a944eeeb0d9a3be83e411b2b89e9a5b46d0cc87aa7603bcf5d7c80ba5ab3caca` |
| native-qualification/job_probe.receipt.json | `dd88c7620274ed0c9a8d119749c8e4b18bfcb1467d4013e1db771331274ef7f3` |
| native-rebuilt/job_probe.py | `fe235b081ad2bd44e053b365f268150b4c1c7e3e49ba25ac1b089d5f91eade3c` |
| native-rebuilt/job_probe.receipt.json | `27dd55b2bcaf5aaeff08f193b925d903c6452e2e26b97601cfbb321eb71148bf` |
| launch_gap_red.py | `80cd6751a8bcc2b4c2be0034250eac2517385f7a2d00298cb1cd95593d643bcd` |
| launch-gap-red.json | `0342168d95d964eb8a1b264c28e80b1572e6826183f74925e6e6af674aea4927` |
| native-atomic/job_probe.receipt.json | `0313bea0d77f4e5f50bd251b57ada3a5ffd69e16158c812b1f99e810ccf3aee7` |
| native-atomic-rebuilt/job_probe.py | `e33dbd14b0cd09262897c98263de0d572ac0b66673df6bce2f62c3bf9ade316c` |
| native-atomic-rebuilt/job_probe.receipt.json | `15297f27c4bb4325845d28f0e29a9ddab63b6341994f5e162bb64ab5a0ad98f8` |
| native-parity.json | `12e00173550140b19b89c3f65e0f9b5fe2e73e200d495c3a00888cacb4474546` |
