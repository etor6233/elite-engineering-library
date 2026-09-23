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
