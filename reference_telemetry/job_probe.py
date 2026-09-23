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
