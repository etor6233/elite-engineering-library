# V349 — native result publication and process-crash reconciliation

Library maintenance / T2809, continued from validated checkpoint140. AUTHORED
qualification apparatus only; CONDITIONED / NOT_ADMITTED as product. Existing
native launch/capture/shutdown code and product owner/profile bytes stay unchanged.

## Concrete delta and invariants

The local result store reserves a 32-hex run ID by exclusive directory creation
before invoking the governed launcher. A flushed claim binds that run ID, the
externally trusted profile digest and selected main-image hash/size. Failure to
publish the claim prevents launch. Any existing reservation, including an empty
one left by a crash, returns ALREADY_RESERVED and never automatically re-executes.
This is a local at-most-one-launch protocol under a trusted store, not exactly-once
business processing or a distributed scheduler/lease/restart policy.

After capture, a bounded final record binds the exact claim digest, launch status,
image identity and outcome status/exit/tree/shutdown plus stdout/stderr lengths
and hashes. No raw output, argv, environment, cwd or friendly profile ID persists.
Digests can still be sensitive metadata; retention/access/collector policy remains
target-owned. RECORDED means receipt publication, never successful domain work.
Nonzero refund exit1, cancellation and forced timeout remain explicit outcomes.

Directory chains are pinned using the prior native no-reparse/handle-path rules.
Claim/result are written to unique same-directory temporary files using exclusive
creation, a checked short-write loop and os.fsync before publication. MoveFileExW
uses WRITE_THROUGH without replace, copy-across-volume or reboot flags. Existing
final names are not overwritten. Failed part files remain as bounded-per-attempt
metadata evidence; no automatic cleanup, retry, deletion or retention policy.

Readers use fixed filenames, bounded reads through held non-reparse file handles,
duplicate-key rejection, exact schemas/types, run/profile/claim hash binding and
outcome checks. They distinguish ABSENT, UNRESOLVED, INVALID and RECORDED. A valid
claim without final result is unresolved; a valid result surviving a lost response
can be recovered by reading. Coordinated rewriting of all records by a trusted
filesystem owner can recompute hashes; this is not signature/WORM/anti-admin
tamper-proof storage. External expected receipt digests/attestation remain separate.

## Failure boundaries and evidence scope

Execution/publish/close failures remain explicit REJECTED, CLAIM_FAILED,
EXECUTION_UNCERTAIN or RESULT_UNCERTAIN. An uncertain return after a successful
rename does not erase the final record; inspect can reconcile it. Unexpected
execution exceptions do not export private exception text. The persisted claim
still forbids automatic retry. No business effect is inferred from absence of
a final record. A new ID is not a safe automatic substitute for a lost attempt.

The initial22tests pass. The expanded27tests also pass, including8independent
processes racing the same ID (one execution and one synthetic effect),24concurrent
thread callers, exclusive publication, short/zero writes, fsync/move failures,
lost publication acknowledgement, partial/oversized/tampered/cross-run records,
reparse rejection, held store path, strict image types and metadata-only output.
One further repeated-reader handle test is authored below and requires the fresh
canonical run before integration. No initial failure was hidden or waived.

Four real owner subprocesses use os._exit at before-claim, after-claim,
after-effect/before-result and after-result boundaries. Fresh parent reads and
retry attempts preserve the appropriate unresolved/recorded outcome, never create
a second synthetic effect and never invoke providers or business accounts.
These are process-crash experiments; the machine/storage controller was not
power-cycled. They do not prove power-loss durability, directory journal ordering,
disk firmware behavior, device loss, backup/restore or remote replication. Store
capacity/quota, permissions, retention and alert delivery remain target gates.

## Authorities and reproduction

Internal authorities: DATABASE_STORAGE_INTERNALS.md persistence/commit/recovery
distinctions; SECURITY_SRE_CLOUD_INFRASTRUCTURE.md privacy/operation boundaries;
ENGINEERING_EXECUTION_PLAYBOOK.md; PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md.

Official APIs were consulted for semantics, not copied as source:
[MoveFileExW](https://learn.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-movefileexw)
describes move flags and overwrite behavior;
[FlushFileBuffers](https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-flushfilebuffers)
describes flushing buffered file data;
[Python os.fsync](https://docs.python.org/3.14/library/os.html#os.fsync)
describes file flushing and its Windows implementation. No rename API test is
promoted into a complete storage durability guarantee. Python docs display3.14.7;
execution stays pinned to CPython3.14.4. No dependency/runtime/source acquisition.

Stage: <LOCALAPPDATA>/Temp/elite-v349-18671f2c25d04ab1b6dd04cb6a282262. Rebuild the two current source blocks below and reuse
the six native qualification files from their V345/V347/V348 canonical sources.
Run test_result_store.py <new-report.json> <refund-fixture.json>, followed by the
retained shutdown/launch/capture suites with their explicit binary/hash inputs.
Same local Windows x64 host, owned synthetic files/processes only. Workspace-owner
license; AUTHORED local adapter, not a Microsoft implementation or certification.

Next: actual worker lifecycle/result-store mounting with domain acknowledgement,
then full service/source/runtime/security/identity/target admission. The local
candidate remains CONDITIONED; no product pack or selected profile is promoted.
Other roadmap fronts remain open; test counts do not define a global percentage.

## Qualification source: result_store.py (AUTHORED / NOT_ADMITTED)

SHA-256: `464d52c47273c9309948edbda371c0587a38de43da41f4a3640f09b7d2f3c669`

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

## Qualification source: test_result_store.py (AUTHORED / NOT_ADMITTED)

SHA-256: `2b3cd608d606413228778e400ca017f4bb671662e795e0e9d95c093376afd314`

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

## Canonical rebuild and integration141

Observed 2026-09-09T02:16:40.754242Z. All eight files reconstruct exactly. The rebuilt result store
passes28tests, including the additional24valid/24invalid reads without native
handle growth. The unchanged prior candidate suites pass19shutdown+20launch+
23capture tests. Total90ordinary tests, zero skips; no result failure is treated
as a successful business operation. Four process-crash cases and one eight-process
race are real child executions. No machine power-cycle or remote effect occurred.

Execution validator1.3.1 materialized15files and capability-gap0.1.0 composed4files;
gap self-check3positives/6negatives PASS. Product owner/profiles and all six reused
native qualification source files remain unchanged. Preflight141 is pending here.

| Evidence | SHA-256 |
|---|---|
| result-tests.json | `0b70414ac48822dd49904f34bca5e880be1010d085e7a5ae7fd8dd6532d692cb` |
| result-extended.json | `0a862436dee8225a2d35adf9382c829d9b421fc5b5524b2f2656deadc92ed49d` |
| rebuilt-result.json | `a0d6088680bcaca86b4f5fe5901e8d704ca69d3caef80249bdafd68861e7537e` |
| rebuilt-result.log | `ed0bee8768bef2cc3923ee1f1e1b2714de50e48e84b5dd8d111a8227bdcad3d0` |
| rebuilt-shutdown.json | `6b0586d112235212d1a8b31b663dd99b84706630eb4de97ea6e80463f6e4c549` |
| rebuilt-launch.json | `9dc3c8efed68851cf6170b0706956b637ab147fd4a1987fce53099d9da149d36` |
| rebuilt-capture.json | `c4b869ed05b0f2082f785da240aab9346736b414b47724192311db00336fa798` |
| canonical-parity.json | `99934f56ddccf28fc93d25866449bb0e04ec378e8a820d99374d61ad8422904e` |
| baseline140.json | `0a66e656c858d6294e7a1effae354ba544baf7f63815825b4cbb2d79efafd003` |

## Closure142 — local result publication qualification complete

Preflight141 completed with 154 executed checks PASS. Structural inventory:
162 packs/1461 materializable files/787 Markdown/53 profiles.
Docker remains unavailable. Network/provider/target opt-in skips remain unexecuted;
this is not an overall production PASS. The pending statement above preserves the
integration141 snapshot and is superseded by this closure.

The eight qualification files did not change after exact rebuild and28new+
62retained tests PASS, zero skips. Four process-crash boundaries and the
eight-process single-ID race remain actual subprocess evidence. No new failure
entry was required in this wave; existing ledger 2525 unique IDs stays preserved.
Product owner/profiles and all six prior native files remain unchanged.

The next work is actual worker lifecycle/result-store mounting with domain
acknowledgement, followed by full trusted profile/identity/runtime/security and
storage/target admission. RECORDED is metadata publication only; UNRESOLVED never
authorizes an automatic retry or a fresh ID. Process-crash evidence does not admit
power-loss durability, WORM, exactly-once business work or a real service rollout.
Other roadmap fronts remain open; no global remaining percentage is computed.

Final successor decision/receipt: gap-record142.json and gap-receipt142.json in
this stage. Integration141 artifacts remain immutable historical snapshots.

| Closure evidence | SHA-256 |
|---|---|
| preflight141.json | `c6fb9f00aa5603a6068ebb84af96b0fdd68d771ed30bc06dfa306a8e2a3c5f68` |
| preflight141.log | `01b049870c929b16f54396ef652fde3f60ee3ba12d1a6d2d72f211b24bc57684` |
| structural-closure141.json | `4d56e977d12584bb4c9cfa2b8ff0d271a4c5662d4bd835402e690e16b431be91` |
| gap-record141.json | `a58e048061fdc9ca7ec1b1de458de02b12034b22348402bb96a22c0b88a98d77` |
| gap-receipt141.json | `aec24517f7996c4d99231e4bc85901814730f412f21ef4e74876359963c12ef1` |
