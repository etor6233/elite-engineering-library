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
