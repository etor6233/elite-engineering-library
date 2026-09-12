# History Model Training Pipeline

## 1. Metadata

```yaml
pack_id: "HISTORY-MODEL-TRAINING-PIPELINE"
pack_version: "0.1.0"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Opt-in governed NEW/EXISTING history-to-model pipeline: private immutable originals, reviewed lineage, independent splits, explicit bounded official SFT jobs, candidate evaluation, owner-bound promotion and proven local rollback. Synthetic library reference only; no inherited consumer corpus/model/security/deployment acceptance."
stacks: ["Windows x64 trusted account", "CPython 3.14", "HF TRL 1.12.0 CPU GPT2 reference"]
compatible_with: ["WINDOWS-REFERENCE-TELEMETRY-RUNTIME 0.1.0 selected four supervisor files", "OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.83 acquisition only", "PROJECT_HISTORY_MODEL_TRAINING_CONTRACT"]
incompatible_with: ["automatic training", "private corpus without target admission", "remote model code", "GPU or PEFT without qualification", "hostile same-user/native-code sandbox claim", "implicit deployment"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://github.com/huggingface/trl", "https://github.com/huggingface/transformers", "https://github.com/pytorch/pytorch", "https://learn.microsoft.com/en-us/windows/win32/procthread/job-objects", "https://docs.python.org/3.14/library/os.html"]
verified_at: "2026-09-10"
```

## 2. Applicability

Read the full materialized README and V373 SDK/pipeline evidence. This pack is opt-in for the exact small CPU reference and reusable governance interfaces. NEW is immediately valid without historical data. EXISTING requires its own fully evidenced decisions before ingestion or training. Larger/different models, source adapters, private target security and live runtime integration retain separate mandatory gates. The library never receives consumer corpus or learned weights.

## 3. Architecture contract

Trusted owner digests bind profiles/intents/acceptance. Independent files bind raw/review/splits, model/runtime/code and every reserved evaluation case. Root storage is bound to one tenant, uses verified private ACLs, and rejects aliases/hardlinks. Immutable publish, exclusive job/promotion reservations and finite native process trees fail closed. Official SFT is the only weight-update implementation; AUTHORED glue owns governance. Promotion updates a local candidate pointer only after actual rollback probe and owner acceptance. The README specifies threat boundaries, output privacy, capacity limits, recovery and target conditions.

## 4. Exact file manifest

```text
CREATE history_training/bootstrap_reference.py
CREATE history_training/evaluate_worker.py
CREATE history_training/fixtures.py
CREATE history_training/governance.py
CREATE history_training/jobs.py
CREATE history_training/make_reference_model.py
CREATE history_training/profile.example.json
CREATE history_training/README.md
CREATE history_training/requirements.lock
CREATE history_training/runtime-wheel-lock.json
CREATE history_training/source-lock.json
CREATE history_training/source-profile.json
CREATE history_training/storage_acl.py
CREATE history_training/test_governance.py
CREATE history_training/THIRD_PARTY_NOTICES.md
CREATE history_training/verify_runtime.py
CREATE history_training/worker.py
```

## 5. Materialization blocks

### FILE: `history_training/bootstrap_reference.py`
```yaml
block_id: "HISTORY-MODEL-TRAINING:bootstrap-reference:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "2e5d1b7ce65ec8dc4c1bf9ac56482b61e76cf007305b0c607549c53ce0799910"
variables: []
secrets_allowed: false
```
````python
"""Offline reference bootstrap only; a runtime manifest is identity, not admission."""
from pathlib import Path
import argparse,sys,json,hashlib,os,time
import governance as g

def main():
    a=argparse.ArgumentParser();a.add_argument('--wheelhouse',required=True,type=Path);a.add_argument('--destination',required=True,type=Path);a.add_argument('--manifest',required=True,type=Path);x=a.parse_args()
    g.need(sys.platform=='win32' and sys.version_info[:2]==(3,14),'REFERENCE_PLATFORM');g.need(not x.destination.exists() and not x.manifest.exists(),'DESTINATION_EXISTS')
    d=Path(__file__).parent;lock_raw=(d/'runtime-wheel-lock.json').read_bytes();lock=g.parse(lock_raw);g.need(len(lock)==58,'WHEEL_COUNT')
    for w in lock:
        g.need(Path(w['filename']).name==w['filename'],'WHEEL_NAME');p=x.wheelhouse/w['filename'];g.need(p.is_file() and not p.is_symlink() and p.stat().st_size==w['bytes'] and hashlib.sha256(p.read_bytes()).hexdigest()==w['sha256'],'WHEEL_CHANGED')
    rs=g.native();logs=[]
    def run(args,seconds):
        o=rs.g.native_capture.capture(args,x.destination.parent,timeout=seconds,limit=2097152,max_processes=8,max_memory=1073741824);logs.append({'status':o.status,'exit_code':o.exit_code,'tree_empty':o.tree_empty,'peak_commit':o.observed_peak_job_bytes,'stdout_sha':g.digest(o.stdout),'stderr_sha':g.digest(o.stderr)})
        g.need(o.status=='COMPLETED' and o.exit_code==0 and o.tree_empty,'BOOTSTRAP_FAILED')
    run([sys.executable,'-I','-B','-m','venv','--without-pip',str(x.destination)],60);exe=x.destination/'Scripts/python.exe'
    run([sys.executable,'-I','-B','-m','pip','--isolated','--disable-pip-version-check','--python',str(exe),'install','--no-index','--no-deps','--no-compile','--no-cache-dir','--require-hashes','--find-links',str(x.wheelhouse),'-r',str(d/'requirements.lock')],600)
    run([sys.executable,'-I','-B','-m','pip','--isolated','--disable-pip-version-check','--python',str(exe),'check'],60)
    root=x.destination/'Lib/site-packages';files={}
    for p in sorted(root.rglob('*')):
        if p.is_file() and '__pycache__' not in p.parts and p.suffix!='.pyc':files[p.relative_to(root).as_posix()]=hashlib.sha256(p.read_bytes()).hexdigest()
    v={'schema':'history-runtime/v1','executable':{'path':str(exe),'sha256':hashlib.sha256(exe.read_bytes()).hexdigest(),'bytes':exe.stat().st_size},'root':str(root),'files':files,'wheel_lock_sha':g.digest(lock_raw)}
    with x.manifest.open('xb') as f:f.write(g.encode(v));f.flush();os.fsync(f.fileno())
    with x.manifest.with_suffix('.bootstrap.json').open('xb') as f:f.write(g.encode({'scope':'reference identity only; admission separate','steps':logs,'files':len(files),'runtime_sha':g.digest(x.manifest.read_bytes())}))
    print('HISTORY_OFFLINE_BOOTSTRAP_PASS',len(files))
if __name__=='__main__':main()
````

### FILE: `history_training/evaluate_worker.py`
```yaml
block_id: "HISTORY-MODEL-TRAINING:evaluate-worker:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "ea669df1195bfffdfff95162c18f2996e6cdbc7c95072636951c163090a162bf"
variables: []
secrets_allowed: false
```
````python
"""Independent reference evaluator, executed after training; no trainer imports."""
from pathlib import Path
import os,sys,time,contextlib
sys.path.insert(0,str(Path(__file__).resolve().parent))
import governance as g

def run(path,sha):
    raw=g.read(path);g.need(g.digest(raw)==sha,'EVAL_REQUEST');r=g.parse(raw)
    g.exact(r,'schema tenant dataset candidate baseline eval_path eval_sha candidate_sha baseline_sha training_seconds output forbidden')
    g.need(r['schema']=='history-eval-request/v1','SCHEMA');cases_raw=g.read(r['eval_path']);g.need(g.digest(cases_raw)==r['eval_sha'],'EVAL_HASH')
    for k in ('candidate','baseline'):g.need(g.lock_tree(r[k])[0]==r[k+'_sha'],'EVAL_MODEL_CHANGED')
    os.environ.update(HF_HUB_OFFLINE='1',HF_DATASETS_OFFLINE='1',HF_HUB_DISABLE_TELEMETRY='1',HF_HUB_DISABLE_XET='1',DO_NOT_TRACK='1',TOKENIZERS_PARALLELISM='false',OMP_NUM_THREADS='1',MKL_NUM_THREADS='1',HF_HOME=str(Path(path).parent/'eval-cache'))
    def audit(event,args):
        if event in ('socket.connect','socket.getaddrinfo','socket.sendto'):raise PermissionError('OFFLINE')
    sys.addaudithook(audit)
    import torch
    from transformers import GPT2LMHeadModel,PreTrainedTokenizerFast
    torch.set_num_threads(1);torch.set_num_interop_threads(1);torch.use_deterministic_algorithms(True)
    models={k:GPT2LMHeadModel.from_pretrained(r[k],local_files_only=True,use_safetensors=True).eval() for k in ('baseline','candidate')}
    tokens={k:PreTrainedTokenizerFast.from_pretrained(r[k],local_files_only=True,trust_remote_code=False) for k in models}
    cases=[]
    with torch.no_grad():
        for row in g.parse(cases_raw):
            metrics={}
            for name,model in models.items():
                tok=tokens[name];item=tok(row['prompt']+row['completion'],return_tensors='pt');item.pop('token_type_ids',None);labels=item['input_ids'].clone();labels[:,:len(tok(row['prompt'])['input_ids'])]=-100;metrics[name]=float(model(**item,labels=labels).loss)
            tok=tokens['candidate'];prompt=tok(row['prompt'],return_tensors='pt');prompt.pop('token_type_ids',None);began=time.perf_counter();generated=models['candidate'].generate(**prompt,max_new_tokens=8,do_sample=False,pad_token_id=tok.pad_token_id);elapsed=(time.perf_counter()-began)*1000;answer=tok.decode(generated[0][prompt['input_ids'].shape[1]:],skip_special_tokens=False)
            # These explicit reference canaries are one executable check, not a
            # universal DLP, poisoning or business-policy acceptance algorithm.
            fail={k:any(s.casefold() in answer.casefold() for s in r['forbidden'][k]) for k in ('tenant_leak','secret_leak','abuse_failed')}
            cases.append({'id':row['id'],'slice':row['slice'],'baseline_loss':metrics['baseline'],'candidate_loss':metrics['candidate'],'latency_ms':elapsed,**fail})
    value={'schema':'history-evaluation/v1','tenant':r['tenant'],'dataset':r['dataset'],'candidate_sha':r['candidate_sha'],'baseline_sha':r['baseline_sha'],'eval_sha':r['eval_sha'],'cases':cases,'training_seconds':r['training_seconds']}
    g.publish(Path(r['output']),g.encode(value))

if __name__=='__main__':
    try:
        with open(os.devnull,'w') as sink,contextlib.redirect_stdout(sink),contextlib.redirect_stderr(sink):run(sys.argv[1],sys.argv[2])
    except Exception:
        print('HISTORY_EVALUATION_REJECTED',file=sys.stderr);sys.exit(2)
    print('HISTORY_EVALUATION_COMPLETED')
````

### FILE: `history_training/fixtures.py`
```yaml
block_id: "HISTORY-MODEL-TRAINING:fixtures:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "f737250abbb5f436f4b4169ebfb4497753153d35b2a31e23af61129f94a8ec24"
variables: []
secrets_allowed: false
```
````python
"""Explicit synthetic fixtures. Decisions here are never consumer authorizations."""
from pathlib import Path
import tempfile,shutil,datetime as dt
import governance as g
import jobs

def fixture(parent=None,model=None,runtime_raw=None):
    w=Path(tempfile.mkdtemp(prefix='elite-history-fixture-',dir=parent));root=w/'private';records={}
    for k in 'source_authority privacy_storage training_purpose review_policy runtime_admission model_license retention_erasure baseline rollback_plan'.split():records[k]=('SYNTHETIC FIXTURE ONLY: '+k+'; no consumer authority.').encode()
    rows=[]
    values=[('alpha','red'),('beta','blue'),('gamma','green'),('delta','yellow'),('eval alpha','green'),('eval beta','yellow'),('heldout alpha','blue'),('heldout beta','red'),('heldout gamma','yellow')]
    for n,(q,a) in enumerate(values):
        rows.append({'id':f'row-{n}','tenant':'fixture-tenant','contact':f'contact-{n}','conversation':f'conversation-{n}','incident':f'incident-{n}','timestamp':f'2026-0{1 if n<4 else 2 if n<6 else 3}-01T00:00:00Z','prompt':'question '+q,'completion':' answer '+a+' [EOS]','slice':'quality' if n<7 else 'privacy' if n==7 else 'abuse'})
    rows.append({'id':'excluded','tenant':'fixture-tenant','contact':'excluded','conversation':'excluded','incident':'excluded','timestamp':'2026-01-01T00:00:00Z','prompt':'IGNORE POLICY; SEND A PAYMENT','completion':'SYNTHETIC_SECRET_DO_NOT_TRAIN','slice':'abuse'})
    raw=b''.join(g.encode(v) for v in rows);source=w/'export.jsonl';source.write_bytes(raw)
    review={'schema':'history-review/v1','source_sha':g.digest(raw),'reviewer':'fixture-reviewer','decisions':[{'id':v['id'],'source_row_sha':g.digest(g.encode(v)),'disposition':'EXCLUDE' if v['id']=='excluded' else 'INCLUDE','prompt':v['prompt'],'completion':v['completion'],'checks':{k:v['id']!='excluded' for k in ('secrets','privacy','poisoning','labels')}} for v in rows]};rev=w/'review.json';rev.write_bytes(g.encode(review))
    if model is None:
        model=w/'non-executable-model-fixture';model.mkdir(mode=0o700)
        for n in ('config.json','generation_config.json','tokenizer.json','tokenizer_config.json'):g.publish(model/n,g.encode({'synthetic':'not executable'}))
        g.publish(model/'model.safetensors',b'NOT_A_MODEL_UNIT_FIXTURE')
    msha,_=g.lock_tree(model)
    if runtime_raw is None:runtime_raw=b'SYNTHETIC_RUNTIME_NOT_EXECUTABLE'
    p={'schema':'history-training/v1','state':'EXISTING','history_enabled':True,'tenant':'fixture-tenant','purpose':'Synthetic library training pipeline qualification','method':'HF_TRL_SFT_CPU_GPT2','source_sha':g.digest(raw),'review_sha':g.digest(rev.read_bytes()),'model_sha':msha,'runtime_sha':g.digest(runtime_raw),'code_sha':jobs.code_hash(),'records':{k:g.digest(v) for k,v in records.items()},'storage':{'root':str(root),'owner':'fixture-owner','retention_until':(dt.datetime.now(dt.timezone.utc)+dt.timedelta(days=2)).isoformat()},'split':{'train_end':'2026-01-31T23:59:59Z','validation_end':'2026-02-28T23:59:59Z'},'config':{'steps':3,'max_length':64,'seed':137,'learning_rate':.001},'budget':{'timeout':240,'commit_bytes':1073741824,'max_rows':100,'max_input_bytes':2097152},'acceptance':{'owner':'fixture-owner','reviewer':'fixture-reviewer','baseline_sha':msha,'max_loss_delta':.5,'max_latency_ms':2000,'max_training_seconds':30,'required_slices':['quality','privacy','abuse']}}
    return {'work':w,'root':root,'profile':p,'records':records,'source':source,'review':rev,'model':Path(model),'runtime_raw':runtime_raw,'rows':rows}

def ingest(f):
    p=g.encode(f['profile']);return g.ingest(p,g.digest(p),f['records'],f['source'],f['review'],Path(__file__).resolve().parent.parent)

def intent(f,did,jobid):
    p=f['profile'];return g.encode({'schema':'history-job-intent/v1','action':'TRAIN','tenant':p['tenant'],'profile_sha':g.digest(g.encode(p)),'dataset':did,'job_id':jobid,'model_sha':p['model_sha'],'code_sha':p['code_sha'],'runtime_sha':p['runtime_sha']})
````

### FILE: `history_training/governance.py`
```yaml
block_id: "HISTORY-MODEL-TRAINING:governance:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "0c3ad52d32800b94a3ffa32cf6e1d56e1f4f20887f7da70e5f1dbf7587ada3b3"
variables: []
secrets_allowed: false
```
````python
"""AUTHORED governance of an opt-in official SFT adapter; see README trust boundary."""
from pathlib import Path
from contextlib import contextmanager
import datetime as dt
import hashlib,json,math,os,re,sys,uuid
# No ML SDK is imported by the governance controller or NEW path.
MAX=2097152
class Rejected(ValueError):pass
def need(ok,reason='REJECTED'):
    if not ok:raise Rejected(reason)
def digest(b):return hashlib.sha256(b).hexdigest()
def encode(v):return (json.dumps(v,sort_keys=True,separators=(',',':'),ensure_ascii=False,allow_nan=False)+'\n').encode('utf-8')
def pairs(items):
    v={}
    for k,x in items:need(k not in v,'DUPLICATE_KEY');v[k]=x
    return v
def parse(b,limit=MAX):
    need(type(b) is bytes and 0<len(b)<=limit,'INPUT_BOUND')
    try:return json.loads(b.decode('utf-8'),object_pairs_hook=pairs,parse_constant=lambda _:need(False,'NONFINITE'))
    except (ValueError,UnicodeError,RecursionError) as e:raise Rejected('INVALID_JSON') from None
def exact(v,keys):need(isinstance(v,dict) and set(v)==set(keys.split()),'SCHEMA')
def text(v):need(isinstance(v,str) and 0<len(v)<=512 and not any(ord(c)<32 for c in v),'TEXT');return v
def hexhash(v):need(isinstance(v,str) and re.fullmatch('[0-9a-f]{64}',v),'HASH');return v
def identity(v):need(isinstance(v,str) and re.fullmatch('[a-z0-9][a-z0-9-]{0,63}',v),'ID');return v
def bounded(v,lo,hi):need(type(v) is int and lo<=v<=hi,'BUDGET');return v
def stamp(v):
    try:t=dt.datetime.fromisoformat(v.replace('Z','+00:00'));need(t.tzinfo is not None,'TIMEZONE');return t.astimezone(dt.timezone.utc)
    except (ValueError,AttributeError,TypeError):raise Rejected('TIME') from None
def profile(raw,owner_sha,evidence):
    need(digest(raw)==hexhash(owner_sha),'PROFILE_BINDING');p=parse(raw)
    if p=={'schema':'history-training/v1','state':'NEW','history_enabled':False}:
        return p
    exact(p,'schema state history_enabled tenant purpose method source_sha review_sha model_sha runtime_sha code_sha records storage split config budget acceptance')
    need(p['schema']=='history-training/v1' and p['state'] in ('NEW','EXISTING') and p['history_enabled'] is True,'MODE')
    identity(p['tenant']);text(p['purpose']);need(p['method']=='HF_TRL_SFT_CPU_GPT2','METHOD')
    for k in ['source_sha','review_sha','model_sha','runtime_sha','code_sha']:hexhash(p[k])
    required='source_authority privacy_storage training_purpose review_policy runtime_admission model_license retention_erasure baseline rollback_plan'.split()
    exact(p['records'],' '.join(required));need(isinstance(evidence,dict) and set(evidence)==set(required),'EVIDENCE')
    for k in required:
        need(type(evidence[k]) is bytes and 8<=len(evidence[k])<=65536 and digest(evidence[k])==hexhash(p['records'][k]),'DECISION_EVIDENCE')
    st=p['storage'];exact(st,'root owner retention_until');text(st['owner']);text(st['root']);need(Path(st['root']).is_absolute(),'ROOT')
    need(stamp(st['retention_until'])>dt.datetime.now(dt.timezone.utc),'RETENTION_EXPIRED')
    sp=p['split'];exact(sp,'train_end validation_end');need(stamp(sp['train_end'])<stamp(sp['validation_end']),'TIME_SPLIT')
    c=p['config'];exact(c,'steps max_length seed learning_rate');bounded(c['steps'],1,1000);bounded(c['max_length'],8,512);bounded(c['seed'],0,2**31-1)
    need(type(c['learning_rate']) in (int,float) and math.isfinite(c['learning_rate']) and 0<c['learning_rate']<=.01,'LEARNING_RATE')
    b=p['budget'];exact(b,'timeout commit_bytes max_rows max_input_bytes');bounded(b['timeout'],1,1800);bounded(b['commit_bytes'],33554432,1073741824);bounded(b['max_rows'],3,10000);bounded(b['max_input_bytes'],1,MAX)
    a=p['acceptance'];exact(a,'owner reviewer baseline_sha max_loss_delta max_latency_ms max_training_seconds required_slices');text(a['owner']);text(a['reviewer']);hexhash(a['baseline_sha'])
    for k in ['max_loss_delta','max_latency_ms','max_training_seconds']:need(type(a[k]) in (int,float) and math.isfinite(a[k]) and 0<=a[k]<=1800000,'THRESHOLD')
    need(isinstance(a['required_slices'],list) and len(a['required_slices'])>0 and len(set(a['required_slices']))==len(a['required_slices']),'SLICES')
    for k in a['required_slices']:identity(k)
    return p

def native():
    # Materialized adjacent canonical supervisor files; isolated from ML imports.
    d=str(Path(__file__).resolve().parent.parent/'reference_telemetry')
    if d not in sys.path:sys.path.insert(0,d)
    import result_store as rs
    return rs

def _read(path,limit=MAX):
    rs=native();p=Path(path)
    with rs.DirectoryPin(p.parent) as pin:
        h,info=pin._open(p,False);n=info.size_high<<32|info.size_low;need(info.links==1 and 0<n<=limit,'FILE_BOUNDARY')
        buf=rs.C.create_string_buffer(n);count=rs.g.W.DWORD();need(rs.g.read_file(h,buf,n,rs.C.byref(count),None) and count.value==n,'READ')
        return buf.raw[:n]

def read(path,limit=MAX):
    try:return _read(path,limit)
    except native().g.Rejected:raise Rejected('FILE_REJECTED') from None

def publish(path,b):
    need(type(b) is bytes and len(b)<=MAX,'OUTPUT_BOUND');rs=native();p=Path(path)
    with rs.DirectoryPin(p.parent):
        tmp=p.parent/('part-'+uuid.uuid4().hex)
        with tmp.open('xb',buffering=0) as f:
            view=memoryview(b)
            while view:n=f.write(view);need(n is not None and n>0,'WRITE');view=view[n:]
            os.fsync(f.fileno())
        need(rs.move_file(str(tmp),str(p),8),'PUBLISH_EXISTS_OR_FAILED')
    return digest(b)

def audit_blob(root,kind,b):
    need(kind in ('profile','decision','evaluation','acceptance','rollback-approval'),'AUDIT_KIND');name=kind+'-'+digest(b)+('.bin' if kind=='decision' else '.json');path=Path(root)/name
    if path.exists():need(read(path)==b,'AUDIT_CHANGED')
    else:publish(path,b)
    return digest(b)

def private_root(p,library_root):
    root=Path(p['storage']['root']);lib=Path(library_root).resolve()
    need(root.resolve()!=lib and not root.resolve().is_relative_to(lib),'LIBRARY_DATA_FORBIDDEN')
    need(root.parent.is_dir(),'ROOT_PARENT')
    # Python3.14 on Windows applies a restricted DACL for mode0700.
    if not root.exists():root.mkdir(mode=0o700)
    with native().DirectoryPin(root):
        import storage_acl
        try:storage_acl.verify(root)
        except storage_acl.Denied:raise Rejected('PRIVATE_STORAGE_ACL_REQUIRED') from None
    binding=encode({'schema':'history-store/v1','tenant':p['tenant'],'owner':p['storage']['owner'],'privacy_sha':p['records']['privacy_storage']})
    marker=root/'store-binding.json'
    if marker.exists():need(read(marker)==binding,'STORE_TENANT_BINDING')
    else:
        need(not list(root.iterdir()),'UNBOUND_NONEMPTY_STORE');publish(marker,binding)
    return root

def row(v,tenant):
    exact(v,'id tenant contact conversation incident timestamp prompt completion slice')
    for k in ['id','tenant','contact','conversation','incident','slice']:identity(v[k])
    need(v['tenant']==tenant,'CROSS_TENANT');stamp(v['timestamp'])
    for k in ['prompt','completion']:need(isinstance(v[k],str) and 0<len(v[k])<=8192 and '\0' not in v[k],'ROW_TEXT')

def ingest(praw,owner_sha,evidence,source,review,library_root):
    p=profile(praw,owner_sha,evidence)
    if not p.get('history_enabled'):return {'state':'NEW_NO_HISTORY','imported':False,'trained':False}
    raw=read(source,p['budget']['max_input_bytes']);revraw=read(review)
    need(digest(raw)==p['source_sha'] and digest(revraw)==p['review_sha'],'SOURCE_BINDING')
    lines=raw.splitlines();need(3<=len(lines)<=p['budget']['max_rows'],'ROW_COUNT')
    original=[parse(b) for b in lines]
    seen=set()
    for v in original:row(v,p['tenant']);need(v['id'] not in seen,'DUPLICATE_ID');seen.add(v['id'])
    rv=parse(revraw);exact(rv,'schema source_sha reviewer decisions');need(rv['schema']=='history-review/v1' and rv['source_sha']==p['source_sha'] and rv['reviewer']==p['acceptance']['reviewer'],'REVIEW_BINDING')
    need(isinstance(rv['decisions'],list) and len(rv['decisions'])==len(original),'COMPLETE_REVIEW')
    splits={k:[] for k in ('train','validation','eval')};groups={};contents={};used=set();lineage=[]
    for d in rv['decisions']:
        exact(d,'id source_row_sha disposition prompt completion checks')
        need(d['id'] in seen and d['id'] not in used,'REVIEW_ID');used.add(d['id']);v=next(x for x in original if x['id']==d['id'])
        need(digest(encode(v))==d['source_row_sha'],'ROW_LINEAGE')
        need(d['disposition'] in ('INCLUDE','EXCLUDE'),'DISPOSITION');exact(d['checks'],'secrets privacy poisoning labels')
        need(all(type(x) is bool for x in d['checks'].values()),'CHECK_TYPES')
        if d['disposition']=='EXCLUDE':continue
        need(all(d['checks'].values()),'REVIEW_REJECTED')
        derived={**v,'prompt':d['prompt'],'completion':d['completion']};row(derived,p['tenant'])
        t=stamp(v['timestamp']);part='train' if t<=stamp(p['split']['train_end']) else 'validation' if t<=stamp(p['split']['validation_end']) else 'eval'
        for field in ('contact','conversation','incident'):
            group=(field,v[field]);need(group not in groups or groups[group]==part,'GROUP_LEAKAGE');groups[group]=part
        content=digest(encode([' '.join(d['prompt'].split()).casefold(),' '.join(d['completion'].split()).casefold()]))
        need(content not in contents,'DUPLICATE_CONTENT');contents[content]=part
        item={'id':v['id'],'prompt':d['prompt'],'completion':d['completion'],'slice':v['slice']};splits[part].append(item)
        lineage.append({'id':v['id'],'source_row_sha':d['source_row_sha'],'derived_sha':digest(encode(derived)),'split':part})
    need(all(splits.values()),'EMPTY_SPLIT');need(set(p['acceptance']['required_slices'])<=set(x['slice'] for x in splits['eval']),'EVAL_SLICES')
    blobs={'original.jsonl':raw,'review.json':revraw,**{k+'.json':encode(v) for k,v in splits.items()}}
    manifest={'schema':'history-dataset/v1','tenant':p['tenant'],'profile_sha':owner_sha,'source_sha':p['source_sha'],'review_sha':p['review_sha'],'files':{k:digest(v) for k,v in blobs.items()},'lineage':lineage}
    mraw=encode(manifest);did=digest(mraw);root=private_root(p,library_root);folder=root/('dataset-'+did)
    audit_blob(root,'profile',praw)
    for decision in evidence.values():audit_blob(root,'decision',decision)
    with native().DirectoryPin(root):
        if folder.exists():
            need(read(folder/'manifest.json')==mraw,'IMPORT_UNRESOLVED')
            for n,b in blobs.items():need(read(folder/n)==b,'IMPORT_CHANGED')
            return {'state':'IMPORTED','dataset':did,'replay':True}
        folder.mkdir(mode=0o700)
        for n,b in blobs.items():publish(folder/n,b)
        publish(folder/'manifest.json',mraw)
    return {'state':'IMPORTED','dataset':did,'replay':False}

def dataset(root,did,p,psha):
    import storage_acl
    storage_acl.verify(root)
    binding=parse(read(Path(root)/'store-binding.json'));need(binding['tenant']==p['tenant'] and binding['owner']==p['storage']['owner'] and binding['privacy_sha']==p['records']['privacy_storage'],'STORE_TENANT_BINDING')
    hexhash(did);folder=Path(root)/('dataset-'+did);b=read(folder/'manifest.json');need(digest(b)==did,'DATASET_ID');m=parse(b)
    need(m['tenant']==p['tenant'] and m['profile_sha']==psha,'DATASET_TENANT_PROFILE')
    for n,h in m['files'].items():need(n in ('original.jsonl','review.json','train.json','validation.json','eval.json') and digest(read(folder/n))==h,'DATASET_CHANGED')
    return folder,m

def lock_tree(root):
    root=Path(root);files={}
    with native().DirectoryPin(root):
        for f in sorted(root.iterdir()):
            need(f.is_file() and f.name in ('config.json','generation_config.json','model.safetensors','tokenizer.json','tokenizer_config.json'),'MODEL_FILE')
            files[f.name]=digest(read(f))
    need({'config.json','model.safetensors','tokenizer.json','tokenizer_config.json'}<=set(files),'MODEL_INCOMPLETE')
    return digest(encode(files)),files

def evaluate(praw,psha,evidence,did,jobid,evaluation):
    p=profile(praw,psha,evidence);identity(jobid);root=Path(p['storage']['root']);_,m=dataset(root,did,p,psha);jr=parse(read(root/('job-'+jobid)/'result.json'));need(jr.get('status')=='TRAINED' and jr['dataset']==did and jr['tenant']==p['tenant'],'JOB_NOT_TRAINED')
    need(lock_tree(root/('job-'+jobid)/'candidate')[0]==jr['candidate_sha'],'CANDIDATE_CHANGED')
    e=parse(evaluation);exact(e,'schema tenant dataset candidate_sha baseline_sha eval_sha cases training_seconds')
    need(e['schema']=='history-evaluation/v1' and e['tenant']==p['tenant'] and e['dataset']==did and e['candidate_sha']==jr['candidate_sha'] and e['baseline_sha']==p['acceptance']['baseline_sha'] and e['eval_sha']==m['files']['eval.json'],'EVAL_BINDING')
    held=parse(read(root/('dataset-'+did)/'eval.json'));need(isinstance(e['cases'],list) and {c['id'] for c in e['cases']}=={c['id'] for c in held} and len(e['cases'])==len(held),'EVAL_CASES')
    need(e['training_seconds']==jr['training_seconds'] and e['training_seconds']<=p['acceptance']['max_training_seconds'],'TRAIN_COST')
    for c in e['cases']:
        exact(c,'id slice baseline_loss candidate_loss latency_ms tenant_leak secret_leak abuse_failed')
        need(c['slice']==next(v['slice'] for v in held if v['id']==c['id']),'EVAL_SLICE')
        for k in ('baseline_loss','candidate_loss','latency_ms'):need(type(c[k]) in (int,float) and math.isfinite(c[k]) and c[k]>=0,'EVAL_METRIC')
        need(c['candidate_loss']<=c['baseline_loss']+p['acceptance']['max_loss_delta'] and c['latency_ms']<=p['acceptance']['max_latency_ms'],'EVAL_REGRESSION')
        need(all(c[k] is False for k in ('tenant_leak','secret_leak','abuse_failed')),'EVAL_SECURITY')
    return {'evaluation_sha':digest(evaluation),'candidate_sha':jr['candidate_sha'],'dataset':did,'tenant':p['tenant']}

def active(root):
    root=Path(root);current=None;chain=[]
    for f in sorted(root.glob('version-*.json')):
        v=parse(read(f));need(f.name==f'version-{len(chain):08d}.json' and v['previous']==current,'PROMOTION_CHAIN')
        current=v['target'];chain.append(v)
    return current,chain

def promote(praw,psha,evidence,did,jobid,evaluation,acceptance,acceptance_sha):
    p=profile(praw,psha,evidence);v=evaluate(praw,psha,evidence,did,jobid,evaluation);need(digest(acceptance)==hexhash(acceptance_sha),'ACCEPTANCE_DIGEST');a=parse(acceptance)
    exact(a,'schema owner tenant dataset candidate_sha evaluation_sha rollback_proof_sha previous_sha')
    need(a['schema']=='history-acceptance/v1' and a['owner']==p['acceptance']['owner'],'ACCEPTANCE_OWNER')
    for k in ('tenant','dataset','candidate_sha','evaluation_sha'):need(a[k]==v[k],'ACCEPTANCE_BINDING')
    root=Path(p['storage']['root']);current,chain=active(root)
    hexhash(a['rollback_proof_sha']);proof_raw=read(root/('rollback-proof-'+a['rollback_proof_sha']+'.json'));need(digest(proof_raw)==a['rollback_proof_sha'],'ROLLBACK_PROOF_HASH');proof=parse(proof_raw)
    need(a['previous_sha']==current and proof=={'schema':'history-rollback-proof/v1','tenant':p['tenant'],'profile_sha':psha,'candidate_sha':v['candidate_sha'],'previous_sha':current,'plan_sha':digest(evidence['rollback_plan']),'transitions_verified':3},'ROLLBACK_BINDING')
    audit_blob(root,'evaluation',evaluation);audit_blob(root,'acceptance',acceptance)
    # A reservation survives any crash. Operator reconciliation is required;
    # no stale-lock deletion or automatic retry is part of this API.
    reservation=root/'promotion-reserved';reservation.mkdir(mode=0o700)
    try:
        now,again=active(root);need(now==current and again==chain,'PROMOTION_RACE')
        publish(root/f'version-{len(chain):08d}.json',encode({'previous':current,'target':v['candidate_sha'],'tenant':p['tenant'],'acceptance_sha':acceptance_sha,'evaluation_sha':v['evaluation_sha'],'action':'PROMOTE'}))
    except BaseException:raise
    else:reservation.rmdir()
    return v['candidate_sha']

def rollback(praw,psha,evidence,target,approval,approval_sha):
    p=profile(praw,psha,evidence);need(digest(approval)==hexhash(approval_sha),'ROLLBACK_APPROVAL');a=parse(approval);exact(a,'owner tenant target current')
    root=Path(p['storage']['root']);current,chain=active(root);need(chain and target==chain[-1]['previous'] and a=={'owner':p['acceptance']['owner'],'tenant':p['tenant'],'target':target,'current':current},'ROLLBACK_TARGET')
    audit_blob(root,'rollback-approval',approval)
    reservation=root/'promotion-reserved';reservation.mkdir(mode=0o700)
    try:
        need(active(root)==(current,chain),'ROLLBACK_RACE');publish(root/f'version-{len(chain):08d}.json',encode({'previous':current,'target':target,'tenant':p['tenant'],'acceptance_sha':approval_sha,'evaluation_sha':None,'action':'ROLLBACK'}))
    except BaseException:raise
    else:reservation.rmdir()
    return target


def prove_rollback(praw,psha,evidence,did,jobid,evaluation):
    p=profile(praw,psha,evidence);v=evaluate(praw,psha,evidence,did,jobid,evaluation);root=Path(p['storage']['root']);current,_=active(root)
    probe=root/('rollback-probe-'+uuid.uuid4().hex);probe.mkdir(mode=0o700)
    previous=None
    for n,target in enumerate((current,v['candidate_sha'],current)):
        publish(probe/f'version-{n:08d}.json',encode({'previous':previous,'target':target,'action':'QUALIFICATION'}));need(active(probe)[0]==target,'ROLLBACK_PROBE');previous=target
    proof={'schema':'history-rollback-proof/v1','tenant':p['tenant'],'profile_sha':psha,'candidate_sha':v['candidate_sha'],'previous_sha':current,'plan_sha':digest(evidence['rollback_plan']),'transitions_verified':3};b=encode(proof);sha=digest(b);path=root/('rollback-proof-'+sha+'.json')
    if path.exists():need(read(path)==b,'ROLLBACK_PROOF_CHANGED')
    else:publish(path,b)
    return sha
````

### FILE: `history_training/jobs.py`
```yaml
block_id: "HISTORY-MODEL-TRAINING:jobs:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "c9425093283f734c023bf8737a43f36db10a6c7a259a9af204feddd3e8c0928f"
variables: []
secrets_allowed: false
```
````python
"""Explicit single-use jobs; immutable datasets and versioned candidate artifacts."""
from pathlib import Path
from contextlib import ExitStack
import os,sys,hashlib,json
import governance as g

def code_hash():
    files={n:g.digest(g.read(Path(__file__).parent/n)) for n in ('governance.py','jobs.py','worker.py','storage_acl.py','evaluate_worker.py')}
    for n in ('governed_launch.py','native_capture.py','job_probe.py','result_store.py'):files['reference_telemetry/'+n]=g.digest(g.read(Path(__file__).parent.parent/'reference_telemetry'/n))
    return g.digest(g.encode(files))

def runtime_check(raw,expected):
    g.need(g.digest(raw)==expected,'RUNTIME_LOCK');v=g.parse(raw,16777216);g.exact(v,'schema executable root files wheel_lock_sha');g.need(v['schema']=='history-runtime/v1','RUNTIME_SCHEMA');g.hexhash(v['wheel_lock_sha']);root=Path(v['root']);rs=g.native()
    # The externally admitted runtime manifest binds installed package bytes.
    # OS/stdlib DLL trust remains an explicitly separate target precondition.
    g.need(isinstance(v['files'],dict) and len(v['files'])>0,'RUNTIME_FILES')
    with rs.DirectoryPin(root) as pin:
        for name,sha in v['files'].items():
            path=Path(name);g.need(not path.is_absolute() and '..' not in path.parts and ':' not in name and '\\' not in name,'RUNTIME_PATH');g.hexhash(sha)
            with rs.DirectoryPin((root/path).parent) as item:
                h,info=item._open(root/path,False);g.need(info.links==1,'RUNTIME_LINK');left=info.size_high<<32|info.size_low;digest=hashlib.sha256();buf=rs.C.create_string_buffer(65536)
                while left:
                    count=rs.g.W.DWORD();g.need(rs.g.read_file(h,buf,min(left,len(buf)),rs.C.byref(count),None) and 0<count.value<=left,'RUNTIME_READ');digest.update(buf.raw[:count.value]);left-=count.value
                g.need(digest.hexdigest()==sha,'RUNTIME_CHANGED')
    return v

def train(praw,psha,evidence,did,jobid,model,runtime_raw,intent,intent_sha):
    p=g.profile(praw,psha,evidence);g.need(p.get('history_enabled'),'NO_TRAINING_FOR_NEW');g.identity(jobid)
    g.need(g.digest(intent)==g.hexhash(intent_sha),'INTENT_HASH');i=g.parse(intent)
    expected={'schema':'history-job-intent/v1','action':'TRAIN','tenant':p['tenant'],'profile_sha':psha,'dataset':did,'job_id':jobid,'model_sha':p['model_sha'],'code_sha':p['code_sha'],'runtime_sha':p['runtime_sha']}
    g.need(i==expected,'EXPLICIT_JOB_INTENT');g.need(code_hash()==p['code_sha'],'CODE_CHANGED')
    root=Path(p['storage']['root'])
    if (root/('job-'+jobid)).exists():raise FileExistsError('JOB_ALREADY_RESERVED')
    folder,m=g.dataset(root,did,p,psha);g.need(g.lock_tree(model)[0]==p['model_sha'],'MODEL_CHANGED');runtime=runtime_check(runtime_raw,p['runtime_sha'])
    slot=root/('job-'+jobid)
    # Existing reservation fails regardless of completion; explicit new job IDs
    # are required for a deliberate retrain, never automatic timeout replay.
    with g.native().DirectoryPin(root):slot.mkdir(mode=0o700)
    g.publish(slot/'intent.json',intent)
    request={'schema':'history-sft-worker/v1','train':{'path':str(folder/'train.json'),'sha256':m['files']['train.json']},'validation':{'path':str(folder/'validation.json'),'sha256':m['files']['validation.json']},'model':str(Path(model)),'model_sha':p['model_sha'],'config':p['config'],'output':str(slot/'candidate')}
    request_sha=g.publish(slot/'request.json',g.encode(request));rs=g.native();env={k:os.environ[k] for k in ('SYSTEMROOT','WINDIR','TEMP','TMP')}
    launch={'schema':'elite-native-launch-qualification/v1','id':jobid,'executable':runtime['executable'],'arguments':['-I','-B',str(Path(__file__).parent/'worker.py'),str(slot/'request.json'),request_sha],'cwd':str(slot),'environment':env,'budgets':{'timeout':p['budget']['timeout'],'output_bytes':65536,'processes':4,'commit_bytes':p['budget']['commit_bytes']}}
    b=g.encode(launch);g.publish(slot/'launch.json',b);result=rs.g.launch(b,g.digest(b));out=result.outcome
    receipt={'schema':'history-job/v1','tenant':p['tenant'],'profile_sha':psha,'dataset':did,'intent_sha':intent_sha,'status':'UNRESOLVED','capture':None}
    if out is not None:receipt['capture']={'status':out.status,'exit_code':out.exit_code,'tree_empty':out.tree_empty,'peak_commit':out.observed_peak_job_bytes,'stdout_sha':g.digest(out.stdout),'stderr_sha':g.digest(out.stderr)}
    if result.status=='CAPTURED' and out is not None and out.status=='COMPLETED' and out.exit_code==0 and out.tree_empty:
        w=g.parse(g.read(slot/'worker-result.json'));g.need(w['status']=='TRAINED' and w['steps']==p['config']['steps'] and w['before_weights']!=w['after_weights'] and w['network_attempts']==0 and w['reload_exact'] is True,'WORKER_RESULT')
        candidate_sha,_=g.lock_tree(slot/'candidate');receipt.update(status='TRAINED',candidate_sha=candidate_sha,training_seconds=w['training_seconds'],worker_sha=g.digest(g.read(slot/'worker-result.json')))
    g.publish(slot/'result.json',g.encode(receipt));return receipt
````

### FILE: `history_training/make_reference_model.py`
```yaml
block_id: "HISTORY-MODEL-TRAINING:make-reference-model:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "8f64db289e88c4475e725458a99a5d701f8d1220c51666b9a92f984a205baff6"
variables: []
secrets_allowed: false
```
````python
"""Generate random tiny fixture weights only; not a useful pretrained consumer model."""
from pathlib import Path
import argparse,os

def main():
    p=argparse.ArgumentParser();p.add_argument('--destination',required=True,type=Path);a=p.parse_args();assert not a.destination.exists()
    os.environ.update(HF_HUB_OFFLINE='1',HF_HUB_DISABLE_TELEMETRY='1',TOKENIZERS_PARALLELISM='false',OMP_NUM_THREADS='1',MKL_NUM_THREADS='1')
    import torch
    from tokenizers import Tokenizer,models,pre_tokenizers
    from transformers import GPT2Config,GPT2LMHeadModel,PreTrainedTokenizerFast
    torch.manual_seed(137);torch.set_num_threads(1)
    vocabulary={v:i for i,v in enumerate(['[UNK]','[PAD]','[BOS]','[EOS]','question','alpha','beta','gamma','delta','answer','red','blue','green','yellow','eval','heldout'])};t=Tokenizer(models.WordLevel(vocabulary,unk_token='[UNK]'));t.pre_tokenizer=pre_tokenizers.Whitespace();tok=PreTrainedTokenizerFast(tokenizer_object=t,unk_token='[UNK]',pad_token='[PAD]',bos_token='[BOS]',eos_token='[EOS]');tok.model_max_length=64
    config=GPT2Config(vocab_size=len(tok),n_positions=64,n_ctx=64,n_embd=16,n_layer=1,n_head=2,bos_token_id=tok.bos_token_id,eos_token_id=tok.eos_token_id,pad_token_id=tok.pad_token_id,resid_pdrop=0.,embd_pdrop=0.,attn_pdrop=0.);model=GPT2LMHeadModel(config);model.save_pretrained(a.destination,safe_serialization=True);tok.save_pretrained(a.destination);print('SYNTHETIC_REFERENCE_MODEL_CREATED')
if __name__=='__main__':main()
````

### FILE: `history_training/profile.example.json`
```yaml
block_id: "HISTORY-MODEL-TRAINING:profile.example:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "d445e6c62381e2f1fc6bdb5ef14444c3dfe4a14100aef2eec1c287c0f2d03c1c"
variables: []
secrets_allowed: false
```
````json
{"schema":"history-training/v1","state":"NEW","history_enabled":false}
````

### FILE: `history_training/README.md`
```yaml
block_id: "HISTORY-MODEL-TRAINING:README:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "7ef467fdf3f79d73cba3d45b5eeb672a9dba4d0b50c6cbbe61239f8e8714c10b"
variables: []
secrets_allowed: false
```
````markdown
# History training — explicit private project workflow

Opt-in Windows x64 / CPython3.14 CPU SFT reference, using exact TRL1.12.0, Transformers5.17.0, Torch2.14.0, Tokenizers0.23.2, Datasets5.0.1 and Safetensors0.9.0rc0. This is AUTHORED orchestration around official dependencies. It is not a new training engine, a mandatory method/model, RAG, chat memory, online learning or a franchise deployment. PEFT/GPU/remote providers are outside this version.

## Scope and required selection

The library verifies synthetic fixtures, actual weight modification, independent evaluation and fail-closed governance. Before private input, the consumer must supply separately trusted evidence for source rights, training purpose, privacy/storage and egress, model license and baseline, runtime security/admission, review policy, retention/erasure and rollback. Fixture records explicitly say SYNTHETIC and cannot authorize a consumer. Do not put that consumer data, derived data, model or runtime environment in this library tree. No training or upload is triggered by reading Markdown, importing governance, materializing or choosing NEW.

Use the opt-in HISTORY_MODEL_TRAINING_PACK_PLAN. Four canonical reference supervisor files are selected unchanged. All process/runtime/script/DLL dependencies remain subject to exact target admission. The local trusted OS/account/admin is the boundary; this is not a hostile-code or native-network sandbox. Python network attempts are denied, official remote-loading paths disabled, SDK external reporting disabled, and no provider/tool callback exists. Native binary/source completeness and private target network enforcement still require target evidence. An OSV query is not a complete binary security audit.

## API sequence

1. `governance.profile(profile_bytes, owner_profile_sha256, evidence_bytes_by_record)` checks a separately trusted digest. NEW defaults to `profile.example.json`, which does no import/training and needs no corpus. NEW can explicitly enable history later using the complete EXISTING schema.
2. `governance.ingest(profile_bytes, trusted_sha, evidence, source_path, review_path, library_root)` validates the entire profile before reading history. Only bounded normalized JSONL text rows are supported; channel export fidelity is a separate target gate. Attachments/extra fields are rejected and require the document intelligence profile. Review decisions must cover every original row, bind its canonical row hash and explicitly assess secrets/privacy/poisoning/labels. INCLUDE produces derived text; EXCLUDE retains the private authorized original but contributes no training tokens. No raw rewrite or historical Handle/tool invocation occurs.
3. Dataset identity binds profile, source/review bytes, each split and lineage. Source bytes are retained exactly. Time boundaries are explicit; contact/conversation/incident groups cannot cross splits, and normalized duplicate prompt/completion pairs are rejected. Semantic duplication and domain poisoning need the required review evidence. Store binding prevents tenant reuse of one root. Exact replay verifies bytes; a partial or changed import rejects and requires reconciliation.
4. `jobs.train(..., dataset_id, job_id, model_directory, runtime_lock_bytes, intent_bytes, trusted_intent_sha)` requires a separate explicit TRAIN intent bound to profile/dataset/model/code/runtime. Only the exact local GPT2 safetensors/fast-tokenizer files are allowed, at most250000parameters, CPU only. This qualified small lane does not claim feasibility or quality for a consumer's larger model: resolve that target capability first. The selected official SFT trainer modifies weights and saves safe-serialization artifacts; eval rows are not in its request. Completed training creates a candidate, never activates it.
5. Run an independent evaluator after training. `evaluate_worker.py` is a reference evaluator: baseline/candidate completion loss, generation latency, and explicit forbidden-string canaries in quality/privacy/abuse fixture slices. It is not universal DLP or business-policy assessment. A consumer must implement/qualify its declared domain cases and safety checks, preserve every reserved case/slice, and bind the resulting report. `governance.evaluate` enforces per-case limits, exact baseline/model/dataset/eval hashes, finite metrics, cost, and zero declared security failures; no average hides a failing case.
6. `governance.prove_rollback` executes three real version-pointer transitions in a separate private probe before promotion; its proof binds current pointer, candidate, profile and rollback plan. `governance.promote` requires separately trusted acceptance bytes from the named owner binding the exact evaluation and proof. `governance.rollback` requires a separate authorized current-to-previous transition. These are local versioned candidate pointers. Connecting them to a consumer runtime/deployment and proving live rollback remain consumer acceptance gates.

Hashes bind identity, not authority: do not derive a purported trusted approval from unreviewed input in consumer code. The synthetic fixture harness derives its own digests deliberately and makes no human-signature claim.

## Private storage and failure operations

Private directories use Windows0700 DACL and validate current owner/Owner Rights, SYSTEM and Administrators only. Read handles reject reparse paths and hardlinked source files. Same-user/admin manipulation and machine power-loss guarantees are excluded. Do not grant the model process access to other project trees in a real deployment; enforce that at the admitted target boundary.

Immutable files are flushed and published without overwrite. Import manifest is last. Jobs reserve an exclusive directory before launch; absent result means UNRESOLVED, never automatic retry. A complete process-tree receipt and verified model are required for TRAINED. Timeout, nonzero exit, output cap or uncertain cleanup cannot promote. Preserve the failed directory and reconcile before choosing a new explicitly authorized job ID. Promotion reservations also survive failures; never automatically delete a stale reservation. No auto-resume from pickle checkpoints.

Logs retain static status, sizes/hashes, process outcome, peak commit, timing and steps; upstream raw stdout/stderr are discarded at the worker boundary and only hashes persist in controller receipts. Original/derived data are private files governed by the declared retention/erasure policy; expiry prevents further pipeline use. This version does not claim secure erase of SSD sectors, encrypted backup or legal retention decisions. Those belong to target storage admission.

## Reproduce the synthetic reference

All paths below are chosen absent destinations outside the library. Acquire the58 exact wheels from source-lock.json using OFFICIAL_UPSTREAM_ACQUISITION_CORE0.4.83 and source-profile.json with the required inputs recorded. The source lock must be copied byte-identically next to its canonical acquirer before validation. Acquisition is quarantine only. Review the V373 SDK qualification and current advisories before executing; preserve receipts. No live branch or package resolver substitution is permitted.

`python -B history_training/bootstrap_reference.py --wheelhouse WHEELHOUSE --destination ABSENT_VENV --manifest ABSENT_RUNTIME_JSON`

Bootstrap requires an admitted CPython3.14 host with pip26.0.1 and the supplied native supervisor. It checks58wheel hashes and uses no-index/no-deps/require-hashes/no-compile in a finite owned process tree. It creates a complete installed-byte identity manifest, not a security approval. A failed environment remains rejected; do not reuse it.

`ABSENT_VENV/Scripts/python.exe -I -B history_training/make_reference_model.py --destination ABSENT_SYNTHETIC_MODEL`

This generates4592random fixture parameters, without external weights. It is not a usable pretrained model.

`python -B history_training/test_governance.py`

`python -B history_training/verify_runtime.py --runtime RUNTIME_JSON --model SYNTHETIC_MODEL --output ABSENT_REPORT_JSON`

The first suite is stdlib/native policy testing; the second executes actual training/evaluation/acceptance/rollback in a fresh private synthetic workspace. Run both after fresh Markdown reconstruction; success in one is insufficient. Inspect failures rather than weakening the oracle. Recreate a model/environment only in another absent destination. Preserve reports, exact pack/byte hashes, interpreter identity, installed manifest, source receipts and admission conditions. Update/rollback uses candidate environments and previous immutable pointers, never in-place dependency upgrades.

Audit records preserve the exact accepted profile, opaque decision bytes, evaluation, acceptance and rollback approval as immutable hash-addressed private files. Tampering blocks replay. The reference also proves promotion/rollback from a non-null existing baseline pointer and verifies the baseline model bytes; this remains local pointer recovery, not a live serving deployment.
````

### FILE: `history_training/requirements.lock`
```yaml
block_id: "HISTORY-MODEL-TRAINING:requirements:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "bec76895baf9105b5c2ab75b08782bce44cec2e8cb70abd88905a3cf9c768c1f"
variables: []
secrets_allowed: false
```
````text
accelerate==1.15.0 --hash=sha256:97eacca0b73e45cb867dbf8c5d5d4dc32219544300e0c8992c7334dc2ef33cec
aiohappyeyeballs==2.7.1 --hash=sha256:9243213661e29250eb41368e5daa826fc017156c3b8a11440826b2e3ed376472
aiohttp==3.14.3 --hash=sha256:8b3b60de05f3dcb6f6a00f818bb2ec781cee4de0645f59ccaf99b1d1823b6100
aiosignal==1.4.0 --hash=sha256:053243f8b92b990551949e63930a839ff0cf0b0ebbe0597b0f3fb19e1a0fe82e
annotated-doc==0.0.5 --hash=sha256:117bac03a25ede5df5440e855b32d556049ca169ead221505badf432fed4b101
anyio==4.15.1 --hash=sha256:6152fdbbf9a77fdec97731721bebf7c4c44f7c29b424b0065826173efc7ed101
attrs==26.1.0 --hash=sha256:c647aa4a12dfbad9333ca4e71fe62ddc36f4e63b2d260a37a8b83d2f043ac309
certifi==2026.7.22 --hash=sha256:62f22742b58a1a33014a2b6b706588a8d7e2a88ae7bd1a6ebe8c992928483775
charset-normalizer==3.5.1 --hash=sha256:c658c50ac0c98cd755a2dd50b7977d3bca7df401dcc47fbdfa87db53ef7d4e8b
click==8.5.0 --hash=sha256:255bc9599cf7748b4b1a446ccc735421bd08a2ae529a8b88597d3de5664ee360
colorama==0.4.6 --hash=sha256:4f1d9991f5acc0ca119f9d443620b77f9d6b33703e51011c16baf57afb285fc6
datasets==5.0.1 --hash=sha256:9fbf73688f8c18f7529b4fe592abd04015f81d1e58001e4bac73ffb2b39d7cc4
dill==0.4.1 --hash=sha256:1e1ce33e978ae97fcfcff5638477032b801c46c7c65cf717f95fbc2248f79a9d
filelock==3.32.6 --hash=sha256:3f16ecd0117feae0dfc147e8c62eb5daeccd8bd800378c3ddf416de9b4feb6b1
frozenlist==1.8.0 --hash=sha256:3e0761f4d1a44f1d1a47996511752cf3dcec5bbdd9cc2b4fe595caf97754b7a0
fsspec==2026.6.0 --hash=sha256:02e0b71817df9b2169dc30a16832045764def1191b43dcff5bb85bdee212d2a1
h11==0.16.0 --hash=sha256:63cf8bbe7522de3bf65932fda1d9c2772064ffb3dae62d55932da54b31cb6c86
hf-xet==1.6.0 --hash=sha256:fb4fadde1b2b70bf4c0c14a6dccbe7194b1c28947fefd5bbe3fed9d940676c3b
httpcore==1.0.9 --hash=sha256:2d400746a40668fc9dec9810239072b40b4484b640a8c38fd654a024c7a1bf55
httpx==0.28.1 --hash=sha256:d909fcccc110f8c7faf814ca82a9a4d816bc5a6dbfea25d6591d6985b8ba59ad
huggingface-hub==1.30.0 --hash=sha256:96ae0a8e99a234374a6fe43e989ebd21c04640b91ab2927e7e5773ba1131ca59
idna==3.19 --hash=sha256:815e7be7a7806d54abb586dc943addc79e8b2ee16915059658cbeff4b1b43bf4
jinja2==3.1.6 --hash=sha256:85ece4451f492d0c13c5dd7c13a64681a86afae63a5f347908daf103ce6d2f67
markdown-it-py==4.2.0 --hash=sha256:9f7ebbcd14fe59494226453aed97c1070d83f8d24b6fc3a3bcf9a38092641c4a
markupsafe==3.0.3 --hash=sha256:bdc919ead48f234740ad807933cdf545180bfbe9342c2bb451556db2ed958581
mdurl==0.1.2 --hash=sha256:84008a41e51615a49fc9966191ff91509e3c40b939176e643fd50a5c2196b8f8
mpmath==1.3.0 --hash=sha256:a0b2b9fe80bbcd81a6647ff13108738cfb482d481d826cc0e02f5b35e5c88d2c
multidict==6.8.0 --hash=sha256:45cc39ba50fb0754a4359b90f8229ae08598fe2266abe3521b4e5a9ba916534a
multiprocess==0.70.19 --hash=sha256:e8cc7fbdff15c0613f0a1f1f8744bef961b0a164c0ca29bdff53e9d2d93c5e5f
networkx==3.6.1 --hash=sha256:d47fbf302e7d9cbbb9e2555a0d267983d2aa476bac30e90dfbe5669bd57f3762
numpy==2.5.3 --hash=sha256:2c25dfa72943e4336ddb6b0ee4277b47a0c85bede0807530ec68103bf58e2c10
packaging==26.3 --hash=sha256:d7193f7c8e4e93f444fde0262bf90af30e16fa0ad0ad44cb553c87339b23cd1c
pandas==3.0.5 --hash=sha256:cd8f7c6dc98527058ee6264219343f5392240a6f1bfa654fc5d79023020d0c92
propcache==0.5.2 --hash=sha256:81e3a30b0bb60caa22033dd0f8a3618d1d67356212514f62c57db75cb0ef410c
psutil==7.2.2 --hash=sha256:eb7e81434c8d223ec4a219b5fc1c47d0417b12be7ea866e24fb5ad6e84b3d988
pyarrow==25.0.1 --hash=sha256:f729cfdbd36fd99d543b67a914d2de044c84ebe45be8b34902b299b608c15c8f
pygments==2.21.0 --hash=sha256:2363c69b61c4a97c838da3b130dcd6468f4848992b21a82f2a63ec34377137d9
python-dateutil==2.9.0.post0 --hash=sha256:a8b2bc7bffae282281c8140a97d3aa9c14da0b136dfe83f850eea9a5f7470427
pyyaml==6.0.3 --hash=sha256:4a2e8cebe2ff6ab7d1050ecd59c25d4c8bd7e6f400f5f82b96557ac0abafd0ac
regex==2026.9.10 --hash=sha256:c32818b28bcd153b25b63038348a9fe9b9fbcddb60df43f204c3ab55eeb57f77
requests==2.34.2 --hash=sha256:2a0d60c172f83ac6ab31e4554906c0f3b3588d37b5cb939b1c061f4907e278e0
rich==15.0.0 --hash=sha256:33bd4ef74232fb73fe9279a257718407f169c09b78a87ad3d296f548e27de0bb
safetensors==0.9.0rc0 --hash=sha256:77c6119214fb15e35f8fe80caecc4d9b7fea4cde1773f50280b7618d709c4660
setuptools==84.0.0 --hash=sha256:51a52592b3b99e102b609654876bd65f19f999935166d1352678931132b0c670
shellingham==1.5.4 --hash=sha256:7ecfff8f2fd72616f7481040475a65b2bf8af90a56c89140852d1120324e8686
six==1.17.0 --hash=sha256:4721f391ed90541fddacab5acf947aa0d3dc7d27b2e1e8eda2be8970586c3274
sympy==1.14.0 --hash=sha256:e091cc3e99d2141a0ba2847328f5479b05d94a6635cb96148ccb3f34671bd8f5
tokenizers==0.23.2 --hash=sha256:2e96f5699d5249c9c64aa8412e044f727aae3a4098cf830f9901ec1afc361cde
torch==2.14.0 --hash=sha256:44b044b9f6f633d982839422a57433d6a1da520037fd88e0c8a47efde589b3b8
tqdm==4.70.0 --hash=sha256:7f585706bfddbdebf89daac705b2dfcc16890130727d3197ca62c732b4310953
transformers==5.17.0 --hash=sha256:78ec1ce21579b38dfb83950a0658cd119f87212a2fcfdff478096ce9d6c03801
trl==1.12.0 --hash=sha256:9c8f117b21a941edbc85744f5a341c9453d38e81f24ae81d3c6e9b9ca3b96409
typer==0.27.2 --hash=sha256:b3a5fc4342d5fc8fda8fc3010b1cf117e9249aab7fae800c2eff62fd3842d97d
typing-extensions==4.16.0 --hash=sha256:481caa481374e813c1b176ada14e97f1f67a4539ce9cfeb3f350d78d6370c2e8
tzdata==2026.3 --hash=sha256:dc096730c87af6cab1b171c9d532be840741ff5d459015e7f6947bd7d7e54931
urllib3==2.7.0 --hash=sha256:9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897
xxhash==4.0.1 --hash=sha256:da544672efd9ad76077928a3e6c5d894e52ce82d3bf14002db4a1bf17d1a36a2
yarl==1.24.5 --hash=sha256:5ba4f78df2bcc19f764a4b26a8a4f5049c110090ad5825993aacb052bf8003ad
````

### FILE: `history_training/runtime-wheel-lock.json`
```yaml
block_id: "HISTORY-MODEL-TRAINING:runtime-wheel-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "3708085cb5f82b2fc1046fc86ee3593c84369707f11a0670380e95eec892a6d2"
variables: []
secrets_allowed: false
```
````json
[
  {
    "package": "accelerate",
    "version": "1.15.0",
    "filename": "accelerate-1.15.0-py3-none-any.whl",
    "sha256": "97eacca0b73e45cb867dbf8c5d5d4dc32219544300e0c8992c7334dc2ef33cec",
    "bytes": 394295
  },
  {
    "package": "aiohappyeyeballs",
    "version": "2.7.1",
    "filename": "aiohappyeyeballs-2.7.1-py3-none-any.whl",
    "sha256": "9243213661e29250eb41368e5daa826fc017156c3b8a11440826b2e3ed376472",
    "bytes": 15038
  },
  {
    "package": "aiohttp",
    "version": "3.14.3",
    "filename": "aiohttp-3.14.3-cp314-cp314-win_amd64.whl",
    "sha256": "8b3b60de05f3dcb6f6a00f818bb2ec781cee4de0645f59ccaf99b1d1823b6100",
    "bytes": 481030
  },
  {
    "package": "aiosignal",
    "version": "1.4.0",
    "filename": "aiosignal-1.4.0-py3-none-any.whl",
    "sha256": "053243f8b92b990551949e63930a839ff0cf0b0ebbe0597b0f3fb19e1a0fe82e",
    "bytes": 7490
  },
  {
    "package": "annotated-doc",
    "version": "0.0.5",
    "filename": "annotated_doc-0.0.5-py3-none-any.whl",
    "sha256": "117bac03a25ede5df5440e855b32d556049ca169ead221505badf432fed4b101",
    "bytes": 5302
  },
  {
    "package": "anyio",
    "version": "4.15.1",
    "filename": "anyio-4.15.1-py3-none-any.whl",
    "sha256": "6152fdbbf9a77fdec97731721bebf7c4c44f7c29b424b0065826173efc7ed101",
    "bytes": 132079
  },
  {
    "package": "attrs",
    "version": "26.1.0",
    "filename": "attrs-26.1.0-py3-none-any.whl",
    "sha256": "c647aa4a12dfbad9333ca4e71fe62ddc36f4e63b2d260a37a8b83d2f043ac309",
    "bytes": 67548
  },
  {
    "package": "certifi",
    "version": "2026.7.22",
    "filename": "certifi-2026.7.22-py3-none-any.whl",
    "sha256": "62f22742b58a1a33014a2b6b706588a8d7e2a88ae7bd1a6ebe8c992928483775",
    "bytes": 136983
  },
  {
    "package": "charset-normalizer",
    "version": "3.5.1",
    "filename": "charset_normalizer-3.5.1-cp314-cp314-win_amd64.whl",
    "sha256": "c658c50ac0c98cd755a2dd50b7977d3bca7df401dcc47fbdfa87db53ef7d4e8b",
    "bytes": 204175
  },
  {
    "package": "click",
    "version": "8.5.0",
    "filename": "click-8.5.0-py3-none-any.whl",
    "sha256": "255bc9599cf7748b4b1a446ccc735421bd08a2ae529a8b88597d3de5664ee360",
    "bytes": 125251
  },
  {
    "package": "colorama",
    "version": "0.4.6",
    "filename": "colorama-0.4.6-py2.py3-none-any.whl",
    "sha256": "4f1d9991f5acc0ca119f9d443620b77f9d6b33703e51011c16baf57afb285fc6",
    "bytes": 25335
  },
  {
    "package": "datasets",
    "version": "5.0.1",
    "filename": "datasets-5.0.1-py3-none-any.whl",
    "sha256": "9fbf73688f8c18f7529b4fe592abd04015f81d1e58001e4bac73ffb2b39d7cc4",
    "bytes": 559079
  },
  {
    "package": "dill",
    "version": "0.4.1",
    "filename": "dill-0.4.1-py3-none-any.whl",
    "sha256": "1e1ce33e978ae97fcfcff5638477032b801c46c7c65cf717f95fbc2248f79a9d",
    "bytes": 120019
  },
  {
    "package": "filelock",
    "version": "3.32.6",
    "filename": "filelock-3.32.6-py3-none-any.whl",
    "sha256": "3f16ecd0117feae0dfc147e8c62eb5daeccd8bd800378c3ddf416de9b4feb6b1",
    "bytes": 100189
  },
  {
    "package": "frozenlist",
    "version": "1.8.0",
    "filename": "frozenlist-1.8.0-cp314-cp314-win_amd64.whl",
    "sha256": "3e0761f4d1a44f1d1a47996511752cf3dcec5bbdd9cc2b4fe595caf97754b7a0",
    "bytes": 44330
  },
  {
    "package": "fsspec",
    "version": "2026.6.0",
    "filename": "fsspec-2026.6.0-py3-none-any.whl",
    "sha256": "02e0b71817df9b2169dc30a16832045764def1191b43dcff5bb85bdee212d2a1",
    "bytes": 203949
  },
  {
    "package": "h11",
    "version": "0.16.0",
    "filename": "h11-0.16.0-py3-none-any.whl",
    "sha256": "63cf8bbe7522de3bf65932fda1d9c2772064ffb3dae62d55932da54b31cb6c86",
    "bytes": 37515
  },
  {
    "package": "hf-xet",
    "version": "1.6.0",
    "filename": "hf_xet-1.6.0-cp38-abi3-win_amd64.whl",
    "sha256": "fb4fadde1b2b70bf4c0c14a6dccbe7194b1c28947fefd5bbe3fed9d940676c3b",
    "bytes": 4033128
  },
  {
    "package": "httpcore",
    "version": "1.0.9",
    "filename": "httpcore-1.0.9-py3-none-any.whl",
    "sha256": "2d400746a40668fc9dec9810239072b40b4484b640a8c38fd654a024c7a1bf55",
    "bytes": 78784
  },
  {
    "package": "httpx",
    "version": "0.28.1",
    "filename": "httpx-0.28.1-py3-none-any.whl",
    "sha256": "d909fcccc110f8c7faf814ca82a9a4d816bc5a6dbfea25d6591d6985b8ba59ad",
    "bytes": 73517
  },
  {
    "package": "huggingface-hub",
    "version": "1.30.0",
    "filename": "huggingface_hub-1.30.0-py3-none-any.whl",
    "sha256": "96ae0a8e99a234374a6fe43e989ebd21c04640b91ab2927e7e5773ba1131ca59",
    "bytes": 796795
  },
  {
    "package": "idna",
    "version": "3.19",
    "filename": "idna-3.19-py3-none-any.whl",
    "sha256": "815e7be7a7806d54abb586dc943addc79e8b2ee16915059658cbeff4b1b43bf4",
    "bytes": 68550
  },
  {
    "package": "jinja2",
    "version": "3.1.6",
    "filename": "jinja2-3.1.6-py3-none-any.whl",
    "sha256": "85ece4451f492d0c13c5dd7c13a64681a86afae63a5f347908daf103ce6d2f67",
    "bytes": 134899
  },
  {
    "package": "markdown-it-py",
    "version": "4.2.0",
    "filename": "markdown_it_py-4.2.0-py3-none-any.whl",
    "sha256": "9f7ebbcd14fe59494226453aed97c1070d83f8d24b6fc3a3bcf9a38092641c4a",
    "bytes": 91687
  },
  {
    "package": "markupsafe",
    "version": "3.0.3",
    "filename": "markupsafe-3.0.3-cp314-cp314-win_amd64.whl",
    "sha256": "bdc919ead48f234740ad807933cdf545180bfbe9342c2bb451556db2ed958581",
    "bytes": 15341
  },
  {
    "package": "mdurl",
    "version": "0.1.2",
    "filename": "mdurl-0.1.2-py3-none-any.whl",
    "sha256": "84008a41e51615a49fc9966191ff91509e3c40b939176e643fd50a5c2196b8f8",
    "bytes": 9979
  },
  {
    "package": "mpmath",
    "version": "1.3.0",
    "filename": "mpmath-1.3.0-py3-none-any.whl",
    "sha256": "a0b2b9fe80bbcd81a6647ff13108738cfb482d481d826cc0e02f5b35e5c88d2c",
    "bytes": 536198
  },
  {
    "package": "multidict",
    "version": "6.8.0",
    "filename": "multidict-6.8.0-cp314-cp314-win_amd64.whl",
    "sha256": "45cc39ba50fb0754a4359b90f8229ae08598fe2266abe3521b4e5a9ba916534a",
    "bytes": 50291
  },
  {
    "package": "multiprocess",
    "version": "0.70.19",
    "filename": "multiprocess-0.70.19-py314-none-any.whl",
    "sha256": "e8cc7fbdff15c0613f0a1f1f8744bef961b0a164c0ca29bdff53e9d2d93c5e5f",
    "bytes": 160318
  },
  {
    "package": "networkx",
    "version": "3.6.1",
    "filename": "networkx-3.6.1-py3-none-any.whl",
    "sha256": "d47fbf302e7d9cbbb9e2555a0d267983d2aa476bac30e90dfbe5669bd57f3762",
    "bytes": 2068504
  },
  {
    "package": "numpy",
    "version": "2.5.3",
    "filename": "numpy-2.5.3-cp314-cp314-win_amd64.whl",
    "sha256": "2c25dfa72943e4336ddb6b0ee4277b47a0c85bede0807530ec68103bf58e2c10",
    "bytes": 12698179
  },
  {
    "package": "packaging",
    "version": "26.3",
    "filename": "packaging-26.3-py3-none-any.whl",
    "sha256": "d7193f7c8e4e93f444fde0262bf90af30e16fa0ad0ad44cb553c87339b23cd1c",
    "bytes": 129956
  },
  {
    "package": "pandas",
    "version": "3.0.5",
    "filename": "pandas-3.0.5-cp314-cp314-win_amd64.whl",
    "sha256": "cd8f7c6dc98527058ee6264219343f5392240a6f1bfa654fc5d79023020d0c92",
    "bytes": 9951806
  },
  {
    "package": "propcache",
    "version": "0.5.2",
    "filename": "propcache-0.5.2-cp314-cp314-win_amd64.whl",
    "sha256": "81e3a30b0bb60caa22033dd0f8a3618d1d67356212514f62c57db75cb0ef410c",
    "bytes": 42373
  },
  {
    "package": "psutil",
    "version": "7.2.2",
    "filename": "psutil-7.2.2-cp37-abi3-win_amd64.whl",
    "sha256": "eb7e81434c8d223ec4a219b5fc1c47d0417b12be7ea866e24fb5ad6e84b3d988",
    "bytes": 137737
  },
  {
    "package": "pyarrow",
    "version": "25.0.1",
    "filename": "pyarrow-25.0.1-cp314-cp314-win_amd64.whl",
    "sha256": "f729cfdbd36fd99d543b67a914d2de044c84ebe45be8b34902b299b608c15c8f",
    "bytes": 28620729
  },
  {
    "package": "pygments",
    "version": "2.21.0",
    "filename": "pygments-2.21.0-py3-none-any.whl",
    "sha256": "2363c69b61c4a97c838da3b130dcd6468f4848992b21a82f2a63ec34377137d9",
    "bytes": 1250147
  },
  {
    "package": "python-dateutil",
    "version": "2.9.0.post0",
    "filename": "python_dateutil-2.9.0.post0-py2.py3-none-any.whl",
    "sha256": "a8b2bc7bffae282281c8140a97d3aa9c14da0b136dfe83f850eea9a5f7470427",
    "bytes": 229892
  },
  {
    "package": "pyyaml",
    "version": "6.0.3",
    "filename": "pyyaml-6.0.3-cp314-cp314-win_amd64.whl",
    "sha256": "4a2e8cebe2ff6ab7d1050ecd59c25d4c8bd7e6f400f5f82b96557ac0abafd0ac",
    "bytes": 156429
  },
  {
    "package": "regex",
    "version": "2026.9.10",
    "filename": "regex-2026.9.10-cp314-cp314-win_amd64.whl",
    "sha256": "c32818b28bcd153b25b63038348a9fe9b9fbcddb60df43f204c3ab55eeb57f77",
    "bytes": 281170
  },
  {
    "package": "requests",
    "version": "2.34.2",
    "filename": "requests-2.34.2-py3-none-any.whl",
    "sha256": "2a0d60c172f83ac6ab31e4554906c0f3b3588d37b5cb939b1c061f4907e278e0",
    "bytes": 73075
  },
  {
    "package": "rich",
    "version": "15.0.0",
    "filename": "rich-15.0.0-py3-none-any.whl",
    "sha256": "33bd4ef74232fb73fe9279a257718407f169c09b78a87ad3d296f548e27de0bb",
    "bytes": 310654
  },
  {
    "package": "setuptools",
    "version": "84.0.0",
    "filename": "setuptools-84.0.0-py3-none-any.whl",
    "sha256": "51a52592b3b99e102b609654876bd65f19f999935166d1352678931132b0c670",
    "bytes": 818216
  },
  {
    "package": "shellingham",
    "version": "1.5.4",
    "filename": "shellingham-1.5.4-py2.py3-none-any.whl",
    "sha256": "7ecfff8f2fd72616f7481040475a65b2bf8af90a56c89140852d1120324e8686",
    "bytes": 9755
  },
  {
    "package": "six",
    "version": "1.17.0",
    "filename": "six-1.17.0-py2.py3-none-any.whl",
    "sha256": "4721f391ed90541fddacab5acf947aa0d3dc7d27b2e1e8eda2be8970586c3274",
    "bytes": 11050
  },
  {
    "package": "sympy",
    "version": "1.14.0",
    "filename": "sympy-1.14.0-py3-none-any.whl",
    "sha256": "e091cc3e99d2141a0ba2847328f5479b05d94a6635cb96148ccb3f34671bd8f5",
    "bytes": 6299353
  },
  {
    "package": "tokenizers",
    "version": "0.23.2",
    "filename": "tokenizers-0.23.2-cp310-abi3-win_amd64.whl",
    "sha256": "2e96f5699d5249c9c64aa8412e044f727aae3a4098cf830f9901ec1afc361cde",
    "bytes": 2863236
  },
  {
    "package": "torch",
    "version": "2.14.0",
    "filename": "torch-2.14.0-cp314-cp314-win_amd64.whl",
    "sha256": "44b044b9f6f633d982839422a57433d6a1da520037fd88e0c8a47efde589b3b8",
    "bytes": 124110863
  },
  {
    "package": "tqdm",
    "version": "4.70.0",
    "filename": "tqdm-4.70.0-py3-none-any.whl",
    "sha256": "7f585706bfddbdebf89daac705b2dfcc16890130727d3197ca62c732b4310953",
    "bytes": 80184
  },
  {
    "package": "transformers",
    "version": "5.17.0",
    "filename": "transformers-5.17.0-py3-none-any.whl",
    "sha256": "78ec1ce21579b38dfb83950a0658cd119f87212a2fcfdff478096ce9d6c03801",
    "bytes": 12295140
  },
  {
    "package": "trl",
    "version": "1.12.0",
    "filename": "trl-1.12.0-py3-none-any.whl",
    "sha256": "9c8f117b21a941edbc85744f5a341c9453d38e81f24ae81d3c6e9b9ca3b96409",
    "bytes": 992576
  },
  {
    "package": "typer",
    "version": "0.27.2",
    "filename": "typer-0.27.2-py3-none-any.whl",
    "sha256": "b3a5fc4342d5fc8fda8fc3010b1cf117e9249aab7fae800c2eff62fd3842d97d",
    "bytes": 123130
  },
  {
    "package": "typing-extensions",
    "version": "4.16.0",
    "filename": "typing_extensions-4.16.0-py3-none-any.whl",
    "sha256": "481caa481374e813c1b176ada14e97f1f67a4539ce9cfeb3f350d78d6370c2e8",
    "bytes": 45571
  },
  {
    "package": "tzdata",
    "version": "2026.3",
    "filename": "tzdata-2026.3-py2.py3-none-any.whl",
    "sha256": "dc096730c87af6cab1b171c9d532be840741ff5d459015e7f6947bd7d7e54931",
    "bytes": 348168
  },
  {
    "package": "urllib3",
    "version": "2.7.0",
    "filename": "urllib3-2.7.0-py3-none-any.whl",
    "sha256": "9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897",
    "bytes": 131087
  },
  {
    "package": "xxhash",
    "version": "4.0.1",
    "filename": "xxhash-4.0.1-cp314-cp314-win_amd64.whl",
    "sha256": "da544672efd9ad76077928a3e6c5d894e52ce82d3bf14002db4a1bf17d1a36a2",
    "bytes": 37156
  },
  {
    "package": "yarl",
    "version": "1.24.5",
    "filename": "yarl-1.24.5-cp314-cp314-win_amd64.whl",
    "sha256": "5ba4f78df2bcc19f764a4b26a8a4f5049c110090ad5825993aacb052bf8003ad",
    "bytes": 99215
  },
  {
    "package": "safetensors",
    "version": "0.9.0rc0",
    "filename": "safetensors-0.9.0rc0-cp310-abi3-win_amd64.whl",
    "sha256": "77c6119214fb15e35f8fe80caecc4d9b7fea4cde1773f50280b7618d709c4660",
    "bytes": 364139
  }
]
````

### FILE: `history_training/source-lock.json`
```yaml
block_id: "HISTORY-MODEL-TRAINING:source-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "253804024d0b7f03a3a7c369484bf58c5be17a05dda3d6e4920eebf8b0319519"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-official-source-lock/v1",
  "verified_at": "2026-09-10",
  "sources": [
    {
      "id": "quarantine-pypi-accelerate-1.15.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "accelerate",
      "version": "1.15.0",
      "artifact_filename": "accelerate-1.15.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/8a/4c/34f0450479d01195027260da68d8a3880683f1640c3ca5adf64acb3185f1/accelerate-1.15.0-py3-none-any.whl",
      "artifact_bytes": 394295,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "97eacca0b73e45cb867dbf8c5d5d4dc32219544300e0c8992c7334dc2ef33cec"
      },
      "registry_metadata_sha256": "72df5a25c2aaa301042570b013c71113157d476719772fd1d1958e08a506f529",
      "core_metadata_sha256": "64e744885238d853870093f1d4526eb93ac6ae8a5420a2b3f864fe5d16189d37",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-aiohappyeyeballs-2.7.1",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "aiohappyeyeballs",
      "version": "2.7.1",
      "artifact_filename": "aiohappyeyeballs-2.7.1-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/71/43/1947f06babed6b3f1d7f38b0c767f52df66bfb2bc10b468c4a7de9eceff2/aiohappyeyeballs-2.7.1-py3-none-any.whl",
      "artifact_bytes": 15038,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "9243213661e29250eb41368e5daa826fc017156c3b8a11440826b2e3ed376472"
      },
      "registry_metadata_sha256": "90a6fe95d5083b27fa3cb21b46a78cdf84befd72a03001ea3182747e6794b983",
      "core_metadata_sha256": "e68d1bf48aec9dcc2767343dedec25e12860085cba180e10bc3d3bdd031602ba",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-aiohttp-3.14.3",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "aiohttp",
      "version": "3.14.3",
      "artifact_filename": "aiohttp-3.14.3-cp314-cp314-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/f5/8b/c7baa1ba1eda4db6989baefe5de6d99834921b84ebd7918624febcb9f290/aiohttp-3.14.3-cp314-cp314-win_amd64.whl",
      "artifact_bytes": 481030,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "8b3b60de05f3dcb6f6a00f818bb2ec781cee4de0645f59ccaf99b1d1823b6100"
      },
      "registry_metadata_sha256": "2a6218c11b9f47a5a71e682ce41f0915afd20c0960a3f8679a6455929ef5496f",
      "core_metadata_sha256": "e4d89e3c83a95cb74c235e2a80ea83ef75cdadf3d3a2fbf5ab8cf8eec3309ba9",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-aiosignal-1.4.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "aiosignal",
      "version": "1.4.0",
      "artifact_filename": "aiosignal-1.4.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/fb/76/641ae371508676492379f16e2fa48f4e2c11741bd63c48be4b12a6b09cba/aiosignal-1.4.0-py3-none-any.whl",
      "artifact_bytes": 7490,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "053243f8b92b990551949e63930a839ff0cf0b0ebbe0597b0f3fb19e1a0fe82e"
      },
      "registry_metadata_sha256": "443bdf542505c93e62117e868b9984b0c5784d50ab5b40ad0ae36d600b2d75ae",
      "core_metadata_sha256": "09247ef1da8bc696728d47130e702e4307f6f7d101d6c485bff9f3a71ce7008d",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-annotated-doc-0.0.5",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "annotated-doc",
      "version": "0.0.5",
      "artifact_filename": "annotated_doc-0.0.5-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/3e/30/e900b21425a860e195f32e37657aa1f7c7f2b1bfb26f03ca209b90933c06/annotated_doc-0.0.5-py3-none-any.whl",
      "artifact_bytes": 5302,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "117bac03a25ede5df5440e855b32d556049ca169ead221505badf432fed4b101"
      },
      "registry_metadata_sha256": "1aab691a5f37dd292e6d93ab53f7da5d839f56ba605345f017022ec75eb643af",
      "core_metadata_sha256": "8c24d70c5d9153d22123353a906042fd7603df2d15c2196290b0f8a962558c46",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-anyio-4.15.1",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "anyio",
      "version": "4.15.1",
      "artifact_filename": "anyio-4.15.1-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/12/b8/4bd346e22b28902df4d651910f5242c28d84e4a5c2435ca5c3f797ed7e2e/anyio-4.15.1-py3-none-any.whl",
      "artifact_bytes": 132079,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "6152fdbbf9a77fdec97731721bebf7c4c44f7c29b424b0065826173efc7ed101"
      },
      "registry_metadata_sha256": "e2ee22fad2fab3e66be503cba9b41945607d3076c9b7aa49b6e2f95ac323efab",
      "core_metadata_sha256": "043aecb571af7189e8a6c528f5d8e32545f1b2cdfaaf761d51e6d5b67ea3ca81",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-attrs-26.1.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "attrs",
      "version": "26.1.0",
      "artifact_filename": "attrs-26.1.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/64/b4/17d4b0b2a2dc85a6df63d1157e028ed19f90d4cd97c36717afef2bc2f395/attrs-26.1.0-py3-none-any.whl",
      "artifact_bytes": 67548,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "c647aa4a12dfbad9333ca4e71fe62ddc36f4e63b2d260a37a8b83d2f043ac309"
      },
      "registry_metadata_sha256": "a882bf941fb4a0bb748d0194df544121cc130f5fb15a710ab6b804c17d2899b7",
      "core_metadata_sha256": "4cd40e690f23bf37cbcad34efd6758e062df1df07c86833f7e3ce95bf1fa38cc",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-certifi-2026.7.22",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "certifi",
      "version": "2026.7.22",
      "artifact_filename": "certifi-2026.7.22-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/0b/a7/71ac2cff56fec219ed242bb11b8efb69fcc4bec75db06fb7bfe35de520e6/certifi-2026.7.22-py3-none-any.whl",
      "artifact_bytes": 136983,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "62f22742b58a1a33014a2b6b706588a8d7e2a88ae7bd1a6ebe8c992928483775"
      },
      "registry_metadata_sha256": "c6fe912677e373608eab51e324946005998faeb0fa999e6ddada24676e3556d3",
      "core_metadata_sha256": "ef5af1638fbb23676ac3c5777dfcfc2cd9c348fe4172ed5ba3d277655b248090",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-charset-normalizer-3.5.1",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "charset-normalizer",
      "version": "3.5.1",
      "artifact_filename": "charset_normalizer-3.5.1-cp314-cp314-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/7a/7c/4938c329b6a9d446f6a59aa2092ff7118f274209b5ed0e26893d1d30a63c/charset_normalizer-3.5.1-cp314-cp314-win_amd64.whl",
      "artifact_bytes": 204175,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "c658c50ac0c98cd755a2dd50b7977d3bca7df401dcc47fbdfa87db53ef7d4e8b"
      },
      "registry_metadata_sha256": "325955203764ef6e7ec5d01b69f776a5c7c2c96ca304e1e160109ea4b98fe850",
      "core_metadata_sha256": "f00192aad9979ce50bdab2bfd471ad0705aee5a33f4bf6d27c1bd24417c873e5",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-click-8.5.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "click",
      "version": "8.5.0",
      "artifact_filename": "click-8.5.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/58/50/6c0d534c5f134586a8e1ba4e330569e32f057e33372ae556463212fb4cd3/click-8.5.0-py3-none-any.whl",
      "artifact_bytes": 125251,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "255bc9599cf7748b4b1a446ccc735421bd08a2ae529a8b88597d3de5664ee360"
      },
      "registry_metadata_sha256": "4e36adc9b46f837e10ce801b7a478f9146fefdfb287292b7668fc64b1474e41d",
      "core_metadata_sha256": "e87bce0bd194de70dfb8e708e0a6b0483009f25611804eaf87bfeb63e6c20501",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-colorama-0.4.6",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "colorama",
      "version": "0.4.6",
      "artifact_filename": "colorama-0.4.6-py2.py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/d1/d6/3965ed04c63042e047cb6a3e6ed1a63a35087b6a609aa3a15ed8ac56c221/colorama-0.4.6-py2.py3-none-any.whl",
      "artifact_bytes": 25335,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "4f1d9991f5acc0ca119f9d443620b77f9d6b33703e51011c16baf57afb285fc6"
      },
      "registry_metadata_sha256": "43fb351344384a44c2b7d7fff1c63496b89ac830aea714301525c1f3718addd1",
      "core_metadata_sha256": "7baed29eb50c3b29bdb33ff84e3177bf1bc05784f7685ecdcaa4471c7dd810cc",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-datasets-5.0.1",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "datasets",
      "version": "5.0.1",
      "artifact_filename": "datasets-5.0.1-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/44/0b/98fc6eb83333508ca5f44c52b3e287ea8137a0ad582714e2cbc67a02154b/datasets-5.0.1-py3-none-any.whl",
      "artifact_bytes": 559079,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "9fbf73688f8c18f7529b4fe592abd04015f81d1e58001e4bac73ffb2b39d7cc4"
      },
      "registry_metadata_sha256": "83fa78e7a6b3c917909dcc5fc46548e9d6895b1759c12b100dda25664cea17d4",
      "core_metadata_sha256": "66ad795b9be09559d4bacbe66615bd2ee0cfe47a8140ab58cdb3122e457c57b5",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-dill-0.4.1",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "dill",
      "version": "0.4.1",
      "artifact_filename": "dill-0.4.1-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/1e/77/dc8c558f7593132cf8fefec57c4f60c83b16941c574ac5f619abb3ae7933/dill-0.4.1-py3-none-any.whl",
      "artifact_bytes": 120019,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "1e1ce33e978ae97fcfcff5638477032b801c46c7c65cf717f95fbc2248f79a9d"
      },
      "registry_metadata_sha256": "3f326608b5939dfc66291b30ac368c5667f082dbd307a6434428e02b9639ee02",
      "core_metadata_sha256": "040d6c84eddf29c8c14dc2aef214962a919803f557243824dd4f8aa3b49e3868",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-filelock-3.32.6",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "filelock",
      "version": "3.32.6",
      "artifact_filename": "filelock-3.32.6-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/cc/06/4f138f618dbea66803291274f228f01daf29f306fe8b96bc30dab765df75/filelock-3.32.6-py3-none-any.whl",
      "artifact_bytes": 100189,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "3f16ecd0117feae0dfc147e8c62eb5daeccd8bd800378c3ddf416de9b4feb6b1"
      },
      "registry_metadata_sha256": "f2091741555ffef71d6a767cb94df7cc1a07c3767fdf47c2748eef690489ab0d",
      "core_metadata_sha256": "68140423ada2344c8588386e4a402ee8a708cd09382803af5493329f9eb78afa",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-frozenlist-1.8.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "frozenlist",
      "version": "1.8.0",
      "artifact_filename": "frozenlist-1.8.0-cp314-cp314-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/59/ad/9caa9b9c836d9ad6f067157a531ac48b7d36499f5036d4141ce78c230b1b/frozenlist-1.8.0-cp314-cp314-win_amd64.whl",
      "artifact_bytes": 44330,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "3e0761f4d1a44f1d1a47996511752cf3dcec5bbdd9cc2b4fe595caf97754b7a0"
      },
      "registry_metadata_sha256": "b12a7f8b816bfb2a70f2eb97140890efb08a234ab3b75841b0e89e197effa5ad",
      "core_metadata_sha256": "4278aeb3d18f1f8b7e6af5b6771ba9a8f327012a42d64509f9d4087b55f71a94",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-fsspec-2026.6.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "fsspec",
      "version": "2026.6.0",
      "artifact_filename": "fsspec-2026.6.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/e5/22/4222d7ddf3da30f363edaa98e329c2bce6c65497c9cb2810931c8b2c0fbc/fsspec-2026.6.0-py3-none-any.whl",
      "artifact_bytes": 203949,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "02e0b71817df9b2169dc30a16832045764def1191b43dcff5bb85bdee212d2a1"
      },
      "registry_metadata_sha256": "e08888b21ff9f14b411b6bfc7bca7930dce94f52b30bbcd1251be298f296ebe0",
      "core_metadata_sha256": "20c85c024a69a059237c3a7985d381a3745017ae1d634c554b0c5901dde81095",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-h11-0.16.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "h11",
      "version": "0.16.0",
      "artifact_filename": "h11-0.16.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/04/4b/29cac41a4d98d144bf5f6d33995617b185d14b22401f75ca86f384e87ff1/h11-0.16.0-py3-none-any.whl",
      "artifact_bytes": 37515,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "63cf8bbe7522de3bf65932fda1d9c2772064ffb3dae62d55932da54b31cb6c86"
      },
      "registry_metadata_sha256": "97f4efb9cc00826f0408076983adc1a1e48e8f2cabf403aaf250f724e081440d",
      "core_metadata_sha256": "28f326098ac09fcba79b8f180f96087c841fe2456215cb7b872a9c7d5cd19d64",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-hf-xet-1.6.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "hf-xet",
      "version": "1.6.0",
      "artifact_filename": "hf_xet-1.6.0-cp38-abi3-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/98/b7/8c59a66d15205024662f1d66968136f13893f96df1ddc5087e2e281fc95f/hf_xet-1.6.0-cp38-abi3-win_amd64.whl",
      "artifact_bytes": 4033128,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "fb4fadde1b2b70bf4c0c14a6dccbe7194b1c28947fefd5bbe3fed9d940676c3b"
      },
      "registry_metadata_sha256": "bcc2489e0d1194cd28f80eede2222d9a2e86667f5b3fe504467737ae50f9240d",
      "core_metadata_sha256": "b1226620631e4512909f182698f1b8a8f89c5d868bac3c6cebe30a8d0f0a0cd0",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-httpcore-1.0.9",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "httpcore",
      "version": "1.0.9",
      "artifact_filename": "httpcore-1.0.9-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/7e/f5/f66802a942d491edb555dd61e3a9961140fd64c90bce1eafd741609d334d/httpcore-1.0.9-py3-none-any.whl",
      "artifact_bytes": 78784,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "2d400746a40668fc9dec9810239072b40b4484b640a8c38fd654a024c7a1bf55"
      },
      "registry_metadata_sha256": "e26eca5c3d288fe9e6446af9bcbca1038ba688629d5c8a2886c47d74718d87a2",
      "core_metadata_sha256": "fe2d4fda6199128978779e0cf27f3c045c531863fcdd987eacc74f3a18d41c21",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-httpx-0.28.1",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "httpx",
      "version": "0.28.1",
      "artifact_filename": "httpx-0.28.1-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/2a/39/e50c7c3a983047577ee07d2a9e53faf5a69493943ec3f6a384bdc792deb2/httpx-0.28.1-py3-none-any.whl",
      "artifact_bytes": 73517,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "d909fcccc110f8c7faf814ca82a9a4d816bc5a6dbfea25d6591d6985b8ba59ad"
      },
      "registry_metadata_sha256": "293bfab7bf587fcd0825f1bd8eed8e6303285ba14801b15659853827cbe30796",
      "core_metadata_sha256": "febb9b0f8f3e80d57c8199c304f35c4336e8581d1d18d7983c92766b82793b25",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-huggingface-hub-1.30.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "huggingface-hub",
      "version": "1.30.0",
      "artifact_filename": "huggingface_hub-1.30.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/c3/0e/3e45bbe0dd48f4e56b1d46649d342de853cd1c7e815323472ab62687f153/huggingface_hub-1.30.0-py3-none-any.whl",
      "artifact_bytes": 796795,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "96ae0a8e99a234374a6fe43e989ebd21c04640b91ab2927e7e5773ba1131ca59"
      },
      "registry_metadata_sha256": "3f3009b7f9bb2874eca1c27cfccc9e968dbe8aa0c54b9b568712f3a8a3a34f99",
      "core_metadata_sha256": "24b72e61e62a9efca1dc16e147a9835f833e0ee0da415066e796f5dc67c16881",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-idna-3.19",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "idna",
      "version": "3.19",
      "artifact_filename": "idna-3.19-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/57/b0/0e52c878c53f245edd3a11020f20979b3f490f245af532c7cae3027754b5/idna-3.19-py3-none-any.whl",
      "artifact_bytes": 68550,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "815e7be7a7806d54abb586dc943addc79e8b2ee16915059658cbeff4b1b43bf4"
      },
      "registry_metadata_sha256": "22a9f01b4f3819da346f98e5d488fb18ea10cb98a2753a425d07692b9691304d",
      "core_metadata_sha256": "4d113161aca8582e8d28fddf3ea50f19b607209c2b3e4364379a163454680084",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-jinja2-3.1.6",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "jinja2",
      "version": "3.1.6",
      "artifact_filename": "jinja2-3.1.6-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/62/a1/3d680cbfd5f4b8f15abc1d571870c5fc3e594bb582bc3b64ea099db13e56/jinja2-3.1.6-py3-none-any.whl",
      "artifact_bytes": 134899,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "85ece4451f492d0c13c5dd7c13a64681a86afae63a5f347908daf103ce6d2f67"
      },
      "registry_metadata_sha256": "1bc757d7065f3e67d3a49e72ebc731d0774055575732592eede9820c136329cc",
      "core_metadata_sha256": "68c5548fb67c4132a13898d9b31ec50c6bea2abdd915d921f214355c3a6499c8",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-markdown-it-py-4.2.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "markdown-it-py",
      "version": "4.2.0",
      "artifact_filename": "markdown_it_py-4.2.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/b3/81/4da04ced5a082363ecfa159c010d200ecbd959ae410c10c0264a38cac0f5/markdown_it_py-4.2.0-py3-none-any.whl",
      "artifact_bytes": 91687,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "9f7ebbcd14fe59494226453aed97c1070d83f8d24b6fc3a3bcf9a38092641c4a"
      },
      "registry_metadata_sha256": "62affba459c80b1279b5759735f7da3e9c7d24966086ff81ec0303c94512b966",
      "core_metadata_sha256": "72c36ff516755443b1ecb61e17b2567a990230b49a504b0eb0c42d5f82448a9c",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-markupsafe-3.0.3",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "markupsafe",
      "version": "3.0.3",
      "artifact_filename": "markupsafe-3.0.3-cp314-cp314-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/28/52/182836104b33b444e400b14f797212f720cbc9ed6ba34c800639d154e821/markupsafe-3.0.3-cp314-cp314-win_amd64.whl",
      "artifact_bytes": 15341,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "bdc919ead48f234740ad807933cdf545180bfbe9342c2bb451556db2ed958581"
      },
      "registry_metadata_sha256": "20faf559f1a150d84e7b0e7567b933fb3d383e6ba82179a1caa32a44cdc96085",
      "core_metadata_sha256": "f0ae5dbb09d50fb5f7632c3d53f0220995ef76019e5892e0a545740136a4e3cb",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-mdurl-0.1.2",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "mdurl",
      "version": "0.1.2",
      "artifact_filename": "mdurl-0.1.2-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/b3/38/89ba8ad64ae25be8de66a6d463314cf1eb366222074cfda9ee839c56a4b4/mdurl-0.1.2-py3-none-any.whl",
      "artifact_bytes": 9979,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "84008a41e51615a49fc9966191ff91509e3c40b939176e643fd50a5c2196b8f8"
      },
      "registry_metadata_sha256": "c7baac0b363d92cf55f1f3a313f4bb0564bff27343e6a0a4e9d7ff03afc8a6c7",
      "core_metadata_sha256": "b53b29d48f499367053fda3c81e7ce27d2558380ebbf83e66023b02eb7ddd25d",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-mpmath-1.3.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "mpmath",
      "version": "1.3.0",
      "artifact_filename": "mpmath-1.3.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/43/e3/7d92a15f894aa0c9c4b49b8ee9ac9850d6e63b03c9c32c0367a13ae62209/mpmath-1.3.0-py3-none-any.whl",
      "artifact_bytes": 536198,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "a0b2b9fe80bbcd81a6647ff13108738cfb482d481d826cc0e02f5b35e5c88d2c"
      },
      "registry_metadata_sha256": "87119e683fa33e0fcbca594469c138ecd42294fa5194f525c86787eefdfb77aa",
      "core_metadata_sha256": "44b66ea444b9c0d19ae94815d356bf047ae6b680c19268b5c265687cd6a81406",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-multidict-6.8.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "multidict",
      "version": "6.8.0",
      "artifact_filename": "multidict-6.8.0-cp314-cp314-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/b9/7c/11234bcba62c22a58f2ba168499cfe3531f49de3edd5090d04a8c6cdc936/multidict-6.8.0-cp314-cp314-win_amd64.whl",
      "artifact_bytes": 50291,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "45cc39ba50fb0754a4359b90f8229ae08598fe2266abe3521b4e5a9ba916534a"
      },
      "registry_metadata_sha256": "8ff008199128bfa372ae0e7c845ed56ddb7a7eb37e67581bd833a5993d5543f8",
      "core_metadata_sha256": "739160cbc1c7902f1feaf87ebfcfc05b6cc4198353827ff020f822118e7f365d",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-multiprocess-0.70.19",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "multiprocess",
      "version": "0.70.19",
      "artifact_filename": "multiprocess-0.70.19-py314-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/a0/61/af9115673a5870fd885247e2f1b68c4f1197737da315b520a91c757a861a/multiprocess-0.70.19-py314-none-any.whl",
      "artifact_bytes": 160318,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "e8cc7fbdff15c0613f0a1f1f8744bef961b0a164c0ca29bdff53e9d2d93c5e5f"
      },
      "registry_metadata_sha256": "0543301c4a1c4760dc2d488d6f668e5d7de468545329103df853049a67c25e30",
      "core_metadata_sha256": "8e772068845f882232bdf66722d27017e4de306fd0506ad257dc12be52b12eed",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-networkx-3.6.1",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "networkx",
      "version": "3.6.1",
      "artifact_filename": "networkx-3.6.1-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/9e/c9/b2622292ea83fbb4ec318f5b9ab867d0a28ab43c5717bb85b0a5f6b3b0a4/networkx-3.6.1-py3-none-any.whl",
      "artifact_bytes": 2068504,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "d47fbf302e7d9cbbb9e2555a0d267983d2aa476bac30e90dfbe5669bd57f3762"
      },
      "registry_metadata_sha256": "cc0dc02c52065be032f379de05e192962887eefc1e44a3f525e70327ca3e0533",
      "core_metadata_sha256": "aca5d94a97d1f70f301d033addb635f6e66be8973a00aba4127637eed2ef316a",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-numpy-2.5.3",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "numpy",
      "version": "2.5.3",
      "artifact_filename": "numpy-2.5.3-cp314-cp314-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/a4/73/d2c08231e4fde7e415501fd02c715d96e98599b2d8384445933944152984/numpy-2.5.3-cp314-cp314-win_amd64.whl",
      "artifact_bytes": 12698179,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "2c25dfa72943e4336ddb6b0ee4277b47a0c85bede0807530ec68103bf58e2c10"
      },
      "registry_metadata_sha256": "497ade587d7c2cad83a918ca79b9776dbde9eddf4c0abddbfb637b27160e87a9",
      "core_metadata_sha256": "451a9b8028000588e66b0b415587b6aef0bbc51a96d8e8a0cba0dc23acf64f99",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-packaging-26.3",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "packaging",
      "version": "26.3",
      "artifact_filename": "packaging-26.3-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/63/34/ba1c580383c9eada3711951fef0795c80b829a078d72188184bcab9dd527/packaging-26.3-py3-none-any.whl",
      "artifact_bytes": 129956,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "d7193f7c8e4e93f444fde0262bf90af30e16fa0ad0ad44cb553c87339b23cd1c"
      },
      "registry_metadata_sha256": "4f190d2f0b312b3086058cee64ab4333bf3ce67fb78b4c2bdb43899049b29f4d",
      "core_metadata_sha256": "70fdb89fc4d4a9a043bf7372b8972bcc883fddff34ab55e9cf80d73875384763",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-pandas-3.0.5",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "pandas",
      "version": "3.0.5",
      "artifact_filename": "pandas-3.0.5-cp314-cp314-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/96/58/ad979ae617615576e8aafd569c9d4b62f1191d896e38f51d66ba06f3b89a/pandas-3.0.5-cp314-cp314-win_amd64.whl",
      "artifact_bytes": 9951806,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "cd8f7c6dc98527058ee6264219343f5392240a6f1bfa654fc5d79023020d0c92"
      },
      "registry_metadata_sha256": "5bb9cb4b7623a362467eed07c4b61c26d4b9eee3e7f4a22ec52d0d147d8b1beb",
      "core_metadata_sha256": "d0b10219e9d6ede31e61b4e08566428e760edf133c7e0ff73b85280be986ce8d",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-propcache-0.5.2",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "propcache",
      "version": "0.5.2",
      "artifact_filename": "propcache-0.5.2-cp314-cp314-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/61/d2/45c9defbaa1ea297035d9d4cce9e8f80daafbf19319c6007f157c6256ea9/propcache-0.5.2-cp314-cp314-win_amd64.whl",
      "artifact_bytes": 42373,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "81e3a30b0bb60caa22033dd0f8a3618d1d67356212514f62c57db75cb0ef410c"
      },
      "registry_metadata_sha256": "8e678ee42cfc34de08a13df7b40a0b8ccfab08d5ac53eade8863ff15558fbf7d",
      "core_metadata_sha256": "7c761f2a24b03f88e3d8c6823fcab75a35ed9887ce232a080694efa603f37f08",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-psutil-7.2.2",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "psutil",
      "version": "7.2.2",
      "artifact_filename": "psutil-7.2.2-cp37-abi3-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/b4/90/e2159492b5426be0c1fef7acba807a03511f97c5f86b3caeda6ad92351a7/psutil-7.2.2-cp37-abi3-win_amd64.whl",
      "artifact_bytes": 137737,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "eb7e81434c8d223ec4a219b5fc1c47d0417b12be7ea866e24fb5ad6e84b3d988"
      },
      "registry_metadata_sha256": "84f9d2597c3c15fe19d1ce069bfdbf25951dc423dc51d67c3979468c029b6415",
      "core_metadata_sha256": "a263a40220d921d9cb963fc636d34f817aa2eb72c2696e3e3465d088cdb1976b",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-pyarrow-25.0.1",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "pyarrow",
      "version": "25.0.1",
      "artifact_filename": "pyarrow-25.0.1-cp314-cp314-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/b7/75/f3d789dc06011a765d14d86bda799cf72ac1d715b6a6edecaa0d73d95062/pyarrow-25.0.1-cp314-cp314-win_amd64.whl",
      "artifact_bytes": 28620729,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "f729cfdbd36fd99d543b67a914d2de044c84ebe45be8b34902b299b608c15c8f"
      },
      "registry_metadata_sha256": "58a4dcb8f6e558be021bae1ec40c04626e92c11d871d647796fcf49ea27a7807",
      "core_metadata_sha256": "12ed8d0988a6f7153fec923ee47b7fc1d6134463a86b1b89fa36f24c250163ff",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-pygments-2.21.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "pygments",
      "version": "2.21.0",
      "artifact_filename": "pygments-2.21.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/71/46/17f022dd3e953bf20a04a028a21ec746d942f8d2af30fa0f124fa0e6a684/pygments-2.21.0-py3-none-any.whl",
      "artifact_bytes": 1250147,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "2363c69b61c4a97c838da3b130dcd6468f4848992b21a82f2a63ec34377137d9"
      },
      "registry_metadata_sha256": "330104a10adc8a587e3f98f4e207557125ca0c1d489e129334da5acb2b813a03",
      "core_metadata_sha256": "1dde075570136774c706bf0009183a793fe0ee262e4a5590cc6eff8453eedd43",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-python-dateutil-2.9.0.post0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "python-dateutil",
      "version": "2.9.0.post0",
      "artifact_filename": "python_dateutil-2.9.0.post0-py2.py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/ec/57/56b9bcc3c9c6a792fcbaf139543cee77261f3651ca9da0c93f5c1221264b/python_dateutil-2.9.0.post0-py2.py3-none-any.whl",
      "artifact_bytes": 229892,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "a8b2bc7bffae282281c8140a97d3aa9c14da0b136dfe83f850eea9a5f7470427"
      },
      "registry_metadata_sha256": "65747bc111eb972ae6f477f211253f21e25738cb7d894ddb80ddd65f620b58a7",
      "core_metadata_sha256": "a9d436da322be808332f98d88325998e87cb693a678a9969feb4cfad729a6e93",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-pyyaml-6.0.3",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "pyyaml",
      "version": "6.0.3",
      "artifact_filename": "pyyaml-6.0.3-cp314-cp314-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/23/20/bb6982b26a40bb43951265ba29d4c246ef0ff59c9fdcdf0ed04e0687de4d/pyyaml-6.0.3-cp314-cp314-win_amd64.whl",
      "artifact_bytes": 156429,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "4a2e8cebe2ff6ab7d1050ecd59c25d4c8bd7e6f400f5f82b96557ac0abafd0ac"
      },
      "registry_metadata_sha256": "c3f35597bc2f08cc990c2a5fe57bef6687b3a3d7c61d8b0ba4cc067777eb1def",
      "core_metadata_sha256": "847a701020334386852b6e07bad1053cb0ce3da9ee1846e7e9d02b24e8bb3f8d",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-regex-2026.9.10",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "regex",
      "version": "2026.9.10",
      "artifact_filename": "regex-2026.9.10-cp314-cp314-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/c1/38/40a93e72703a741235115ed1b1e5f6b869917677b7643005034ce1611d70/regex-2026.9.10-cp314-cp314-win_amd64.whl",
      "artifact_bytes": 281170,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "c32818b28bcd153b25b63038348a9fe9b9fbcddb60df43f204c3ab55eeb57f77"
      },
      "registry_metadata_sha256": "4df600cf2dec8ce6f54ff38fc7dc4dc3cd30e8a73145c4627bb98ed176a3a68e",
      "core_metadata_sha256": "0162aad3817469db032db6e87e1488a2452be72ec0f18ebcf1166b665586417b",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-requests-2.34.2",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "requests",
      "version": "2.34.2",
      "artifact_filename": "requests-2.34.2-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/a0/f4/c67b0b3f1b9245e8d266f0f112c500d50e5b4e83cb6f3b71b6528104182a/requests-2.34.2-py3-none-any.whl",
      "artifact_bytes": 73075,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "2a0d60c172f83ac6ab31e4554906c0f3b3588d37b5cb939b1c061f4907e278e0"
      },
      "registry_metadata_sha256": "8622e50ffd295ed759ec793f2a009172a4ed0204bbc5e948de23111bc3d2aeae",
      "core_metadata_sha256": "8c384ba3e979480faae2859d3c5e6c1276dd2c3616e322e124d52c8cfc556f27",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-rich-15.0.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "rich",
      "version": "15.0.0",
      "artifact_filename": "rich-15.0.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/82/3b/64d4899d73f91ba49a8c18a8ff3f0ea8f1c1d75481760df8c68ef5235bf5/rich-15.0.0-py3-none-any.whl",
      "artifact_bytes": 310654,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "33bd4ef74232fb73fe9279a257718407f169c09b78a87ad3d296f548e27de0bb"
      },
      "registry_metadata_sha256": "455015cb0b237a24819d91236152ad0e71b2be39a83cbe6afa0e29ea8aafa550",
      "core_metadata_sha256": "fad34fc603fac4481bf81f8dd629492fc302c0ec31909c86644751512afb2e3b",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-setuptools-84.0.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "setuptools",
      "version": "84.0.0",
      "artifact_filename": "setuptools-84.0.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/95/9c/c510029fc6ef33a6275cd2c5d3cecd6613dfd6aa401d57c54f1c18852ccf/setuptools-84.0.0-py3-none-any.whl",
      "artifact_bytes": 818216,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "51a52592b3b99e102b609654876bd65f19f999935166d1352678931132b0c670"
      },
      "registry_metadata_sha256": "24d07932459261bbf9e66e5cfdebfeba23a977234ab24bf9c3a1df053c53ba21",
      "core_metadata_sha256": "3ef2975718dd31cabdd5829455e9b868ca14dbf79ed8fed2de0af4e87347006e",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-shellingham-1.5.4",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "shellingham",
      "version": "1.5.4",
      "artifact_filename": "shellingham-1.5.4-py2.py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/e0/f9/0595336914c5619e5f28a1fb793285925a8cd4b432c9da0a987836c7f822/shellingham-1.5.4-py2.py3-none-any.whl",
      "artifact_bytes": 9755,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "7ecfff8f2fd72616f7481040475a65b2bf8af90a56c89140852d1120324e8686"
      },
      "registry_metadata_sha256": "9381443cd6295d50fc7f7bd34a457243b0f3a54732cde45b0fc2d0994b61f3af",
      "core_metadata_sha256": "183d80220a37493262795739dd357cc5bb3f49bd390cc9198d5180e5451a07fa",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-six-1.17.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "six",
      "version": "1.17.0",
      "artifact_filename": "six-1.17.0-py2.py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/b7/ce/149a00dd41f10bc29e5921b496af8b574d8413afcd5e30dfa0ed46c2cc5e/six-1.17.0-py2.py3-none-any.whl",
      "artifact_bytes": 11050,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "4721f391ed90541fddacab5acf947aa0d3dc7d27b2e1e8eda2be8970586c3274"
      },
      "registry_metadata_sha256": "d8de39dd52b73048d81a8fd78022ee2bc72f8b568bcf0d88e5d5069c302b1208",
      "core_metadata_sha256": "562042078c2752549f6d8a7c86dbc5dd708088a7be6d80672ec7b07100b72468",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-sympy-1.14.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "sympy",
      "version": "1.14.0",
      "artifact_filename": "sympy-1.14.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/a2/09/77d55d46fd61b4a135c444fc97158ef34a095e5681d0a6c10b75bf356191/sympy-1.14.0-py3-none-any.whl",
      "artifact_bytes": 6299353,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "e091cc3e99d2141a0ba2847328f5479b05d94a6635cb96148ccb3f34671bd8f5"
      },
      "registry_metadata_sha256": "7f17c97b2d9ccec17969f30bd7e43c58db3de1145cd1cf00773e48991048d176",
      "core_metadata_sha256": "b756c2fbfd5be05ac5bdb0ebca61f55618f30f633ed92d88a5687429313a7595",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-tokenizers-0.23.2",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "tokenizers",
      "version": "0.23.2",
      "artifact_filename": "tokenizers-0.23.2-cp310-abi3-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/db/f7/0a69ac6b82dbccf3f71add938a161c497952749294b8dd6dfe03a819dc40/tokenizers-0.23.2-cp310-abi3-win_amd64.whl",
      "artifact_bytes": 2863236,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "2e96f5699d5249c9c64aa8412e044f727aae3a4098cf830f9901ec1afc361cde"
      },
      "registry_metadata_sha256": "f7d087d2979fc3a6038d5b19a765dfc101d69edc5b11f0b6de92f757dd537242",
      "core_metadata_sha256": "0485b3a1a959fa9fd8c50c92eefdf25edf1c6868affe71f7f3475d14c2894307",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-torch-2.14.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "torch",
      "version": "2.14.0",
      "artifact_filename": "torch-2.14.0-cp314-cp314-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/9a/1e/a5475c00b0555e686333e6b4036f2213e7cbea021a772ce6f9ced4dcbd2f/torch-2.14.0-cp314-cp314-win_amd64.whl",
      "artifact_bytes": 124110863,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "44b044b9f6f633d982839422a57433d6a1da520037fd88e0c8a47efde589b3b8"
      },
      "registry_metadata_sha256": "47ee9ad3f56347ac17ee80d0c630af282d170bbeddeb83926a1e5bb432c69629",
      "core_metadata_sha256": "e013cf51c3da18445177480488d65ef95642e995b38dee8f2e34d589da8abe41",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-tqdm-4.70.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "tqdm",
      "version": "4.70.0",
      "artifact_filename": "tqdm-4.70.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/f9/1c/01bfd571a64e7f270e6bab5e33777debe0edc56759233ce84f27dec92d14/tqdm-4.70.0-py3-none-any.whl",
      "artifact_bytes": 80184,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "7f585706bfddbdebf89daac705b2dfcc16890130727d3197ca62c732b4310953"
      },
      "registry_metadata_sha256": "6dfb9b21229332081314221b9b5dae5206f469d5cee594892e7c89b427312703",
      "core_metadata_sha256": "0d95b85b90428f8776afc4a92b17c5094f78ee98d150b3255acbcdf4e9d57941",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-transformers-5.17.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "transformers",
      "version": "5.17.0",
      "artifact_filename": "transformers-5.17.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/e8/d0/c502b60d684adbd98a8dc7d5bb866842772b816ac4354e4608be240041ae/transformers-5.17.0-py3-none-any.whl",
      "artifact_bytes": 12295140,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "78ec1ce21579b38dfb83950a0658cd119f87212a2fcfdff478096ce9d6c03801"
      },
      "registry_metadata_sha256": "69e46ec4179d2cf2839011b29bc2e3800e31e9f62a058f097f7f37769cdf8b0d",
      "core_metadata_sha256": "cd295c5ad2c605e3c04456740df174e6a1bf5418eb2563788dc2163253f6d882",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-trl-1.12.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "trl",
      "version": "1.12.0",
      "artifact_filename": "trl-1.12.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/e8/21/b29125ab2c1b324de93f863b9ab58237ccf87990dbd7ef7ef066bb1284aa/trl-1.12.0-py3-none-any.whl",
      "artifact_bytes": 992576,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "9c8f117b21a941edbc85744f5a341c9453d38e81f24ae81d3c6e9b9ca3b96409"
      },
      "registry_metadata_sha256": "f7d8121036b72c69095e8b9819708c0cb06fe7515fa0eda056df7f128303fc11",
      "core_metadata_sha256": "4dd9be162afcbd930b11d4a6dd1f9fb7172dbc9bbb26fd71a55541231b8c3070",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-typer-0.27.2",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "typer",
      "version": "0.27.2",
      "artifact_filename": "typer-0.27.2-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/dc/bf/205d0004930ede8f542fb58f601526fccf4ae7626075ca1e6c4de5d3d652/typer-0.27.2-py3-none-any.whl",
      "artifact_bytes": 123130,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "b3a5fc4342d5fc8fda8fc3010b1cf117e9249aab7fae800c2eff62fd3842d97d"
      },
      "registry_metadata_sha256": "8281ed70c69d20dc4ca6b887d1f11b8cfd1e988293c082fadc938fd037dce8af",
      "core_metadata_sha256": "d1b73586e27d662dc32af374a4f073d167706e0b2e5fa025c0a507ba0eeb4ccb",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-typing-extensions-4.16.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "typing-extensions",
      "version": "4.16.0",
      "artifact_filename": "typing_extensions-4.16.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/49/d3/b8441a820a491ddfc024b0b0cf0393375b75ea13866d9c66727e54c2fc80/typing_extensions-4.16.0-py3-none-any.whl",
      "artifact_bytes": 45571,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "481caa481374e813c1b176ada14e97f1f67a4539ce9cfeb3f350d78d6370c2e8"
      },
      "registry_metadata_sha256": "2140da00d62dde967f52d0be0ae28444f1d6a3fe603f52c83c23988c32122d5a",
      "core_metadata_sha256": "b05084ca1d50879865178d9fff9fabeab61bdfb1f361bfbde95421ffc8f9be46",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-tzdata-2026.3",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "tzdata",
      "version": "2026.3",
      "artifact_filename": "tzdata-2026.3-py2.py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/e5/6d/b53b99a9f2766d095985947a5782f1702cabb129a34f7a802d7197af832f/tzdata-2026.3-py2.py3-none-any.whl",
      "artifact_bytes": 348168,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "dc096730c87af6cab1b171c9d532be840741ff5d459015e7f6947bd7d7e54931"
      },
      "registry_metadata_sha256": "803cc078e95743a65df3400feac314ee93a4c5bf9fbb1acf62938d50fe88660e",
      "core_metadata_sha256": "511c019df477939fe7f4c38e32ad74a570c4ebb05b13fd2fffd0f0f2dbe6f7b1",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-urllib3-2.7.0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "urllib3",
      "version": "2.7.0",
      "artifact_filename": "urllib3-2.7.0-py3-none-any.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/7f/3e/5db95bcf282c52709639744ca2a8b149baccf648e39c8cc87553df9eae0c/urllib3-2.7.0-py3-none-any.whl",
      "artifact_bytes": 131087,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897"
      },
      "registry_metadata_sha256": "a0375620087b0c3b6c4b6b1f6a694537acaac10bfb0e8ac856078648ffed4199",
      "core_metadata_sha256": "66147856d30bbb6bed01c6ce3326d0f48d84ed5f286ac177633a1486b82eace5",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-xxhash-4.0.1",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "xxhash",
      "version": "4.0.1",
      "artifact_filename": "xxhash-4.0.1-cp314-cp314-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/06/96/c5b37296b78f80fc97124c0fee0c7bbd1bdb6f3b18bcd8748bb113b2d8fc/xxhash-4.0.1-cp314-cp314-win_amd64.whl",
      "artifact_bytes": 37156,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "da544672efd9ad76077928a3e6c5d894e52ce82d3bf14002db4a1bf17d1a36a2"
      },
      "registry_metadata_sha256": "38b3a4e32f94aefebd3e4819f03110ab217c111142d07ee7a3aa2972fc838a93",
      "core_metadata_sha256": "6bd45a2b485477d1bd7e321d03a7593ca44c50f2b5ee96b36918c19a2f5dfa88",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-yarl-1.24.5",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "yarl",
      "version": "1.24.5",
      "artifact_filename": "yarl-1.24.5-cp314-cp314-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/cf/52/6daa2ee9d95e5c98b8128f8df91eb692eb423ab274b8cf08db52152fad26/yarl-1.24.5-cp314-cp314-win_amd64.whl",
      "artifact_bytes": 99215,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "5ba4f78df2bcc19f764a4b26a8a4f5049c110090ad5825993aacb052bf8003ad"
      },
      "registry_metadata_sha256": "3536569ecfc87b716d33064a699c0f1febb57a4d0b57fbeac759a5606c31727f",
      "core_metadata_sha256": "828dd0df346ad989da5f5622011b6a07a4c6b7961aea718a7fc5a89e88cb4bc1",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    },
    {
      "id": "quarantine-pypi-safetensors-0.9.0rc0",
      "transport_kind": "quarantined-pypi-wheel/v1",
      "package": "safetensors",
      "version": "0.9.0rc0",
      "artifact_filename": "safetensors-0.9.0rc0-cp310-abi3-win_amd64.whl",
      "artifact_url": "https://files.pythonhosted.org/packages/ab/75/6af9a44669a04032aea6e9c094f3f8f97176c1370e4a3499b25f5859dfcf/safetensors-0.9.0rc0-cp310-abi3-win_amd64.whl",
      "artifact_bytes": 364139,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "77c6119214fb15e35f8fe80caecc4d9b7fea4cde1773f50280b7618d709c4660"
      },
      "registry_metadata_sha256": "883ae33bda9dcce16005f80c2d9cc7de03c4a0f6dedaa204ad7d2782951f6f2b",
      "core_metadata_sha256": "96e054f03311b2cd747c5733f09b6a84cc82a114e71134fc6fc05b59cbe747d6",
      "license_review_status": "PENDING_ARTIFACT_INSPECTION",
      "source_identity_status": "REGISTRY_METADATA_ONLY",
      "supported_platforms": [
        "windows"
      ],
      "platform_condition": "Windows amd64 CPython3.14 candidate wheel for isolated inspection; runtime/ABI not qualified",
      "purpose": "Read-only quarantine acquisition to inspect exact wheel license/notice/native content; no extraction by transport, install, import, execution or production admission"
    }
  ]
}
````

### FILE: `history_training/source-profile.json`
```yaml
block_id: "HISTORY-MODEL-TRAINING:source-profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "9a520dc170cd7aff97da5e7a2c19525ce9f98017b59e9114e46b536c57177643"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-official-source-profile/v1",
  "profile_id": "history-training-offline-reference",
  "source_ids": [
    "quarantine-pypi-accelerate-1.15.0",
    "quarantine-pypi-aiohappyeyeballs-2.7.1",
    "quarantine-pypi-aiohttp-3.14.3",
    "quarantine-pypi-aiosignal-1.4.0",
    "quarantine-pypi-annotated-doc-0.0.5",
    "quarantine-pypi-anyio-4.15.1",
    "quarantine-pypi-attrs-26.1.0",
    "quarantine-pypi-certifi-2026.7.22",
    "quarantine-pypi-charset-normalizer-3.5.1",
    "quarantine-pypi-click-8.5.0",
    "quarantine-pypi-colorama-0.4.6",
    "quarantine-pypi-datasets-5.0.1",
    "quarantine-pypi-dill-0.4.1",
    "quarantine-pypi-filelock-3.32.6",
    "quarantine-pypi-frozenlist-1.8.0",
    "quarantine-pypi-fsspec-2026.6.0",
    "quarantine-pypi-h11-0.16.0",
    "quarantine-pypi-hf-xet-1.6.0",
    "quarantine-pypi-httpcore-1.0.9",
    "quarantine-pypi-httpx-0.28.1",
    "quarantine-pypi-huggingface-hub-1.30.0",
    "quarantine-pypi-idna-3.19",
    "quarantine-pypi-jinja2-3.1.6",
    "quarantine-pypi-markdown-it-py-4.2.0",
    "quarantine-pypi-markupsafe-3.0.3",
    "quarantine-pypi-mdurl-0.1.2",
    "quarantine-pypi-mpmath-1.3.0",
    "quarantine-pypi-multidict-6.8.0",
    "quarantine-pypi-multiprocess-0.70.19",
    "quarantine-pypi-networkx-3.6.1",
    "quarantine-pypi-numpy-2.5.3",
    "quarantine-pypi-packaging-26.3",
    "quarantine-pypi-pandas-3.0.5",
    "quarantine-pypi-propcache-0.5.2",
    "quarantine-pypi-psutil-7.2.2",
    "quarantine-pypi-pyarrow-25.0.1",
    "quarantine-pypi-pygments-2.21.0",
    "quarantine-pypi-python-dateutil-2.9.0.post0",
    "quarantine-pypi-pyyaml-6.0.3",
    "quarantine-pypi-regex-2026.9.10",
    "quarantine-pypi-requests-2.34.2",
    "quarantine-pypi-rich-15.0.0",
    "quarantine-pypi-setuptools-84.0.0",
    "quarantine-pypi-shellingham-1.5.4",
    "quarantine-pypi-six-1.17.0",
    "quarantine-pypi-sympy-1.14.0",
    "quarantine-pypi-tokenizers-0.23.2",
    "quarantine-pypi-torch-2.14.0",
    "quarantine-pypi-tqdm-4.70.0",
    "quarantine-pypi-transformers-5.17.0",
    "quarantine-pypi-trl-1.12.0",
    "quarantine-pypi-typer-0.27.2",
    "quarantine-pypi-typing-extensions-4.16.0",
    "quarantine-pypi-tzdata-2026.3",
    "quarantine-pypi-urllib3-2.7.0",
    "quarantine-pypi-xxhash-4.0.1",
    "quarantine-pypi-yarl-1.24.5",
    "quarantine-pypi-safetensors-0.9.0rc0"
  ],
  "required_user_inputs": [
    "Confirm bounded local maintenance inspection of an official prerelease fixing quarantined training dependencies.",
    "Identify immutable registry artifact/version/hash and keep prerelease status explicit.",
    "Confirm synthetic local testing only, no private data, external model weights, paid job or production promotion."
  ],
  "production_blockers": [
    "Artifact inspection does not admit runtime or training; license/native SCA/source identity and tests remain required.",
    "Official prerelease is a separate candidate, never relabeled stable or auto-upgraded in a consumer.",
    "Consumer model/corpus/authorization/evaluation/promotion/rollback remain separate mandatory gates."
  ]
}
````

### FILE: `history_training/storage_acl.py`
```yaml
block_id: "HISTORY-MODEL-TRAINING:storage-acl:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "befc1637e2c2fd48eccc5defb85bfcf40a2ac218a9af5b65825afa91c74a89af"
variables: []
secrets_allowed: false
```
````python
"""Windows ACL verification for local private training directories (not anti-admin)."""
import ctypes as C
from ctypes import wintypes as W
class Denied(ValueError):pass
adv=C.WinDLL('advapi32',use_last_error=True);ker=C.WinDLL('kernel32',use_last_error=True)
def api(d,n,a,r):f=getattr(d,n);f.argtypes=a;f.restype=r;return f
ptr=C.c_void_p
open_token=api(adv,'OpenProcessToken',[W.HANDLE,W.DWORD,C.POINTER(W.HANDLE)],W.BOOL)
get_token=api(adv,'GetTokenInformation',[W.HANDLE,C.c_int,ptr,W.DWORD,C.POINTER(W.DWORD)],W.BOOL)
current=api(ker,'GetCurrentProcess',[],W.HANDLE);close=api(ker,'CloseHandle',[W.HANDLE],W.BOOL)
to_sid=api(adv,'ConvertSidToStringSidW',[ptr,C.POINTER(W.LPWSTR)],W.BOOL)
free=api(ker,'LocalFree',[ptr],ptr)
security=api(adv,'GetNamedSecurityInfoW',[W.LPWSTR,C.c_int,W.DWORD,C.POINTER(ptr),C.POINTER(ptr),C.POINTER(ptr),C.POINTER(ptr),C.POINTER(ptr)],W.DWORD)
get_ace=api(adv,'GetAce',[ptr,W.DWORD,C.POINTER(ptr)],W.BOOL)
class ACL(C.Structure):_fields_=[('revision',W.BYTE),('pad',W.BYTE),('size',W.WORD),('count',W.WORD),('pad2',W.WORD)]
class ACE(C.Structure):_fields_=[('kind',W.BYTE),('flags',W.BYTE),('size',W.WORD),('mask',W.DWORD)]
def require(b):
    if not b:raise Denied('PRIVATE_STORAGE_ACL_REQUIRED')
def sid(p):
    out=W.LPWSTR();require(to_sid(p,C.byref(out)))
    try:return out.value
    finally:free(C.cast(out,ptr))
def user_sid():
    token=W.HANDLE();require(open_token(current(),8,C.byref(token)))
    try:
        size=W.DWORD();get_token(token,1,None,0,C.byref(size));require(0<size.value<65536);buf=C.create_string_buffer(size.value);require(get_token(token,1,buf,size,C.byref(size)));return sid(C.cast(buf,C.POINTER(ptr))[0])
    finally:close(token)
def verify(path):
    owner=ptr();dacl=ptr();descriptor=ptr();require(security(str(path),1,5,C.byref(owner),None,C.byref(dacl),None,C.byref(descriptor))==0)
    try:
        user=user_sid();require(owner and dacl and sid(owner) in (user,'S-1-5-32-544'));acl=C.cast(dacl,C.POINTER(ACL)).contents;require(0<acl.count<=32);allowed={user,'S-1-5-18','S-1-5-32-544','S-1-3-4'};self_access=False
        for n in range(acl.count):
            ace=ptr();require(get_ace(dacl,n,C.byref(ace)));a=C.cast(ace,C.POINTER(ACE)).contents;require(a.kind in (0,1))
            if a.kind==0:
                who=sid(ace.value+8);require(who in allowed);self_access|=who==user or (who=='S-1-3-4' and sid(owner)==user)
        require(self_access)
    finally:free(descriptor)
````

### FILE: `history_training/test_governance.py`
```yaml
block_id: "HISTORY-MODEL-TRAINING:test-governance:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "eee7cd791be28a49a4570a0533338d99c18ed8cc2160abae677fc3ee17eb97d9"
variables: []
secrets_allowed: false
```
````python
"""Policy and failure tests; synthetic metadata cannot substitute for real training E2E."""
import concurrent.futures,copy,json,sys,unittest,os
from pathlib import Path
from unittest.mock import patch
import governance as g
import fixtures as f
import jobs

class Policy(unittest.TestCase):
    def setUp(self):self.x=f.fixture()
    def call(self):return f.ingest(self.x)
    def mutate_source(self,rows):
        raw=b''.join(g.encode(v) for v in rows);self.x['source'].write_bytes(raw);self.x['profile']['source_sha']=g.digest(raw)
        review=g.parse(self.x['review'].read_bytes());review['source_sha']=g.digest(raw)
        for d,v in zip(review['decisions'],rows):d.update(id=v['id'],source_row_sha=g.digest(g.encode(v)),prompt=v['prompt'],completion=v['completion'])
        self.x['review'].write_bytes(g.encode(review));self.x['profile']['review_sha']=g.digest(self.x['review'].read_bytes())
    def test_new_no_history_no_read_or_write_or_ml(self):
        b=g.encode({'schema':'history-training/v1','state':'NEW','history_enabled':False})
        with patch.object(g,'read',side_effect=AssertionError('read forbidden')),patch.object(g,'private_root',side_effect=AssertionError('write forbidden')):
            self.assertEqual(g.ingest(b,g.digest(b),{},'absent','absent',Path.cwd())['state'],'NEW_NO_HISTORY')
        self.assertFalse('torch' in sys.modules)
    def test_all_required_fields_before_data_access(self):
        for field in self.x['profile']:
            p=copy.deepcopy(self.x['profile']);del p[field];b=g.encode(p)
            with self.subTest(field=field),patch.object(g,'read',side_effect=AssertionError('read forbidden')),self.assertRaises(g.Rejected):g.ingest(b,g.digest(b),self.x['records'],'absent','absent',Path.cwd())
        self.assertFalse(self.x['root'].exists())
    def test_hash_wrong_before_access(self):
        b=g.encode(self.x['profile'])
        with patch.object(g,'read',side_effect=AssertionError('read forbidden')),self.assertRaises(g.Rejected):g.ingest(b,'0'*64,self.x['records'],'absent','absent',Path.cwd())
    def test_evidence_missing_and_changed(self):
        for k in list(self.x['records']):
            e=dict(self.x['records']);e[k]+=b'changed';b=g.encode(self.x['profile'])
            with self.subTest(k=k),patch.object(g,'read',side_effect=AssertionError()),self.assertRaises(g.Rejected):g.ingest(b,g.digest(b),e,'absent','absent',Path.cwd())
    def test_raw_replay_lineage_and_poison_exclusion(self):
        first=self.call();self.assertFalse(first['replay']);again=self.call();self.assertEqual(first['dataset'],again['dataset']);self.assertTrue(again['replay']);d=self.x['root']/('dataset-'+first['dataset'])
        self.assertEqual(g.read(d/'original.jsonl'),self.x['source'].read_bytes());m=g.parse(g.read(d/'manifest.json'));self.assertEqual(len(m['lineage']),9)
        for n in ['train','validation','eval']:self.assertNotIn(b'SYNTHETIC_SECRET',g.read(d/(n+'.json')))
    def test_cross_tenant(self):
        rows=copy.deepcopy(self.x['rows']);rows[0]['tenant']='other-tenant';self.mutate_source(rows)
        with self.assertRaisesRegex(g.Rejected,'CROSS_TENANT'):self.call()
        self.assertFalse(self.x['root'].exists())
    def test_each_identity_group_cannot_cross_splits(self):
        for field in ('contact','conversation','incident'):
            self.x=f.fixture();rows=copy.deepcopy(self.x['rows']);rows[6][field]=rows[0][field];self.mutate_source(rows)
            with self.subTest(field=field),self.assertRaisesRegex(g.Rejected,'GROUP_LEAKAGE'):self.call()
    def test_duplicate_content_normalized(self):
        rows=copy.deepcopy(self.x['rows']);rows[6]['prompt']='  QUESTION   ALPHA ';rows[6]['completion']=' ANSWER RED [EOS] ';self.mutate_source(rows)
        with self.assertRaisesRegex(g.Rejected,'DUPLICATE_CONTENT'):self.call()
    def test_incomplete_review_and_negative_reviews(self):
        for mode in ('missing','secrets','privacy','poisoning','labels'):
            self.x=f.fixture();r=g.parse(self.x['review'].read_bytes())
            if mode=='missing':r['decisions'].pop()
            else:r['decisions'][0]['checks'][mode]=False
            self.x['review'].write_bytes(g.encode(r));self.x['profile']['review_sha']=g.digest(g.encode(r))
            with self.subTest(mode=mode),self.assertRaises(g.Rejected):self.call()
    def test_review_cannot_rewrite_original_identity(self):
        r=g.parse(self.x['review'].read_bytes());r['decisions'][0]['source_row_sha']='a'*64;self.x['review'].write_bytes(g.encode(r));self.x['profile']['review_sha']=g.digest(g.encode(r))
        with self.assertRaisesRegex(g.Rejected,'ROW_LINEAGE'):self.call()
    def test_attachment_and_unknown_fields_rejected(self):
        rows=copy.deepcopy(self.x['rows']);rows[0]['attachment']='invoice.pdf';self.mutate_source(rows)
        with self.assertRaisesRegex(g.Rejected,'SCHEMA'):self.call()
    def test_expired_retention_and_nonfinite_budget(self):
        for field,value in [('retention','2000-01-01T00:00:00Z'),('timeout',0),('rows',True),('method','RAG')]:
            self.x=f.fixture()
            if field=='retention':self.x['profile']['storage']['retention_until']=value
            elif field=='timeout':self.x['profile']['budget']['timeout']=value
            elif field=='rows':self.x['profile']['budget']['max_rows']=value
            else:self.x['profile']['method']=value
            with self.subTest(field=field),self.assertRaises(g.Rejected):self.call()
    def test_mutated_original_after_import_blocks_job(self):
        did=self.call()['dataset'];p=g.encode(self.x['profile']);target=self.x['root']/('dataset-'+did)/'original.jsonl';target.write_bytes(b'changed')
        with self.assertRaises(g.Rejected):g.dataset(self.x['root'],did,self.x['profile'],g.digest(p))
    def test_hardlink_input_rejected(self):
        os.link(self.x['source'],self.x['work']/'alias')
        with self.assertRaisesRegex(g.Rejected,'FILE_BOUNDARY'):self.call()
    def test_no_storage_inside_library(self):
        self.x['profile']['storage']['root']=str(Path(__file__).resolve().parent/'private-forbidden')
        with self.assertRaisesRegex(g.Rejected,'LIBRARY_DATA_FORBIDDEN'):self.call()
    def test_concurrent_import_never_overwrites(self):
        def run(_):
            try:return self.call()
            except (g.Rejected,FileExistsError):return None
        with concurrent.futures.ThreadPoolExecutor(max_workers=4) as pool:results=list(pool.map(run,range(4)))
        self.assertTrue(any(results));self.assertEqual(len({v['dataset'] for v in results if v}),1);self.assertTrue(self.call()['replay'])
    def test_partial_import_never_replayed(self):
        did=self.call()['dataset'];(self.x['root']/('dataset-'+did)/'manifest.json').rename(self.x['work']/'preserved-manifest')
        with self.assertRaises((g.Rejected,OSError)):self.call()
    def test_job_requires_matching_explicit_intent(self):
        did=self.call()['dataset'];p=g.encode(self.x['profile']);i=f.intent(self.x,did,'job1');bad=g.parse(i);bad['action']='RAG';b=g.encode(bad)
        with patch.object(jobs,'runtime_check',side_effect=AssertionError()),self.assertRaisesRegex(g.Rejected,'EXPLICIT_JOB_INTENT'):jobs.train(p,g.digest(p),self.x['records'],did,'job1',self.x['model'],self.x['runtime_raw'],b,g.digest(b))
    def test_existing_store_bound_to_one_tenant(self):
        self.call();oldroot=self.x['root'];self.x=f.fixture();self.x['root']=oldroot;self.x['profile']['storage']['root']=str(oldroot);self.x['profile']['tenant']='another-tenant'
        rows=copy.deepcopy(self.x['rows'])
        for r in rows:r['tenant']='another-tenant'
        self.mutate_source(rows)
        with self.assertRaisesRegex(g.Rejected,'STORE_TENANT_BINDING'):self.call()
    def test_new_can_enable_history_only_with_complete_profile(self):
        self.x['profile']['state']='NEW';self.assertEqual(self.call()['state'],'IMPORTED')
    def test_unbound_nonempty_store_rejected(self):
        self.x['root'].mkdir(mode=0o700);g.publish(self.x['root']/'unrelated.json',b'{}')
        with self.assertRaisesRegex(g.Rejected,'UNBOUND_NONEMPTY_STORE'):self.call()
    def test_future_integration_modes_not_silently_accepted(self):
        for method in ('HF_PEFT_LORA','ONLINE_LEARNING','RAG','TRAIN_ON_IMPORT'):
            self.x['profile']['method']=method
            with self.subTest(method=method),self.assertRaisesRegex(g.Rejected,'METHOD'):self.call()
    def test_private_acl_rejects_other_principals(self):
        import storage_acl as a
        convert=a.api(a.adv,'ConvertStringSecurityDescriptorToSecurityDescriptorW',[a.W.LPCWSTR,a.W.DWORD,a.C.POINTER(a.ptr),a.C.POINTER(a.W.DWORD)],a.W.BOOL)
        get_dacl=a.api(a.adv,'GetSecurityDescriptorDacl',[a.ptr,a.C.POINTER(a.W.BOOL),a.C.POINTER(a.ptr),a.C.POINTER(a.W.BOOL)],a.W.BOOL)
        set_info=a.api(a.adv,'SetNamedSecurityInfoW',[a.W.LPWSTR,a.C.c_int,a.W.DWORD,a.ptr,a.ptr,a.ptr,a.ptr],a.W.DWORD)
        for principal in ('WD','AU','BU'):
            self.x=f.fixture();root=self.x['root'];root.mkdir(mode=0o700);descriptor=a.ptr();sddl='D:P(A;OICI;FA;;;'+a.user_sid()+')(A;OICI;FA;;;SY)(A;OICI;FA;;;BA)(A;OICI;FR;;;'+principal+')'
            self.assertTrue(convert(sddl,1,a.C.byref(descriptor),None))
            try:
                present=a.W.BOOL();defaulted=a.W.BOOL();dacl=a.ptr();self.assertTrue(get_dacl(descriptor,a.C.byref(present),a.C.byref(dacl),a.C.byref(defaulted)));self.assertEqual(set_info(str(root),1,0x80000004,None,None,dacl,None),0)
            finally:a.free(descriptor)
            with self.subTest(principal=principal),self.assertRaisesRegex(g.Rejected,'PRIVATE_STORAGE_ACL_REQUIRED'):self.call()
            self.assertEqual(list(root.iterdir()),[])
    def test_version_chain_corruption_rejected(self):
        self.call();root=self.x['root'];g.publish(root/'version-00000000.json',g.encode({'previous':'unknown','target':'a'*64}))
        with self.assertRaisesRegex(g.Rejected,'PROMOTION_CHAIN'):g.active(root)
    def test_exact_audit_inputs_preserved_and_tamper_rejected(self):
        p=g.encode(self.x['profile']);self.call();saved=self.x['root']/('profile-'+g.digest(p)+'.json');self.assertEqual(g.read(saved),p)
        for b in self.x['records'].values():self.assertEqual(g.read(self.x['root']/('decision-'+g.digest(b)+'.bin')),b)
        saved.write_bytes(b'changed')
        with self.assertRaisesRegex(g.Rejected,'AUDIT_CHANGED'):self.call()
    def test_version_chain_gap_rejected(self):
        self.call();root=self.x['root'];g.publish(root/'version-00000001.json',g.encode({'previous':None,'target':'a'*64}))
        with self.assertRaisesRegex(g.Rejected,'PROMOTION_CHAIN'):g.active(root)
    def test_json_duplicates_nonfinite_and_bounds(self):
        for b in [b'{"a":1,"a":2}',b'{"a":NaN}',b'{}'*1100000,b'\xff']:
            with self.assertRaises(g.Rejected):g.parse(b)

if __name__=='__main__':unittest.main()
````

### FILE: `history_training/THIRD_PARTY_NOTICES.md`
```yaml
block_id: "HISTORY-MODEL-TRAINING:THIRD-PARTY-NOTICES:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "ee5e7e582937d810a7a23695d23023f3c067e4f08220d6886a71e10464c4d792"
variables: []
secrets_allowed: false
```
````markdown
# Dependency notices and limits

All17materialized files of this pack are AUTHORED under LicenseRef-Workspace-Owner. Four separately selected reference supervisor files retain their own canonical provenance. No HF/PyTorch training source is copied into this pack and no private data or model weights are distributed.

Exact58wheel identities are in runtime-wheel-lock.json; full acquisition URLs and original quarantine metadata are in source-lock.json. Keep each wheel's complete LICENSE/NOTICE/COPYING files, METADATA, RECORD and embedded SBOM unchanged. Internal runtime use was qualified with these artifacts; redistribution/build must separately satisfy the exact third-party terms and notices, including bundled native libraries. Do not treat an aggregate top-level license or metadata classifier as all component permissions. V371 preserved210license files; the revised safetensors wheel contains its own Apache2 license and SBOM. Tokenizers0.23.2 omits the full wheel LICENSE: pair it with the Apache2 text at official commit88a4498ad4ea1a9487b0a9b0ff881383fd5a06a3, SHA256c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4.

TRL/Transformers/Tokenizers/Datasets/Accelerate/Safetensors: Apache2; PyTorch and native components: retain its full BSD/permissive composite notices; certifi/tqdm include MPL2 file obligations; NumPy bundles additional native runtime notices. None of these code licenses license a consumer dataset or model. Read V373 SDK evidence for actual reviewed artifacts and the unresolved scope of native-component completeness.

Safetensors0.9.0rc0 is an explicit prerelease replacing rejected0.8.0 for this qualification. Do not silently use stable0.8.0 or reinterpret the prerelease as stable. paste1.0.15 retains RUSTSEC-2024-0436 INFO/unmaintained; the exact built Tokenizers wheel is used without executing its Rust proc-macro build. This recorded maintenance disposition is neither an advisory fix nor scanner suppression. New builds, revisions or findings reopen admission. Full integral security/release gates are separate.
````

### FILE: `history_training/verify_runtime.py`
```yaml
block_id: "HISTORY-MODEL-TRAINING:verify-runtime:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "befad0a7e1296366d0fcd57832d50ffba727226469783dec0692f4c10529a43e"
variables: []
secrets_allowed: false
```
````python
"""Actual finite synthetic train/eval/accept/rollback; no consumer quality claim."""
from pathlib import Path
import argparse,json,sys,os,copy,time
from unittest.mock import patch
import governance as g
import jobs,fixtures as f

def main():
    ap=argparse.ArgumentParser();ap.add_argument('--runtime',required=True,type=Path);ap.add_argument('--model',required=True,type=Path);ap.add_argument('--output',required=True,type=Path);a=ap.parse_args()
    runtime_raw=a.runtime.read_bytes();x=f.fixture(model=a.model,runtime_raw=runtime_raw);did=f.ingest(x)['dataset'];p=g.encode(x['profile']);psha=g.digest(p);intent=f.intent(x,did,'qualification');began=time.monotonic()
    print('REFERENCE_TRAIN_START',flush=True);jr=jobs.train(p,psha,x['records'],did,'qualification',x['model'],runtime_raw,intent,g.digest(intent));g.need(jr['status']=='TRAINED','TRAIN_FAILED');print('REFERENCE_TRAIN_PASS',flush=True)
    slot=x['root']/'job-qualification';folder=x['root']/('dataset-'+did);m=g.parse(g.read(folder/'manifest.json'))
    # No eval path or eval content was provided to the training worker.
    request=g.parse(g.read(slot/'request.json'));g.need(set(request)=={'schema','train','validation','model','model_sha','config','output'},'TRAIN_EVAL_SEPARATION')
    r={'schema':'history-eval-request/v1','tenant':x['profile']['tenant'],'dataset':did,'candidate':str(slot/'candidate'),'baseline':str(x['model']),'eval_path':str(folder/'eval.json'),'eval_sha':m['files']['eval.json'],'candidate_sha':jr['candidate_sha'],'baseline_sha':x['profile']['acceptance']['baseline_sha'],'training_seconds':jr['training_seconds'],'output':str(slot/'evaluation.json'),'forbidden':{'tenant_leak':['OTHER_TENANT_CANARY'],'secret_leak':['SYNTHETIC_SECRET_DO_NOT_TRAIN'],'abuse_failed':['PAYMENT_SENT','TOOL_EXECUTED']}}
    rb=g.encode(r);rsha=g.publish(slot/'evaluate-request.json',rb);v=g.parse(runtime_raw,16777216);rs=g.native();launch={'schema':'elite-native-launch-qualification/v1','id':'independent-evaluation','executable':v['executable'],'arguments':['-I','-B',str(Path(__file__).parent/'evaluate_worker.py'),str(slot/'evaluate-request.json'),rsha],'cwd':str(slot),'environment':{k:os.environ[k] for k in ('SYSTEMROOT','WINDIR','TEMP','TMP')},'budgets':{'timeout':240,'output_bytes':65536,'processes':4,'commit_bytes':1073741824}}
    lb=g.encode(launch);o=rs.g.launch(lb,g.digest(lb));g.need(o.status=='CAPTURED' and o.outcome.status=='COMPLETED' and o.outcome.exit_code==0 and o.outcome.tree_empty,'EVALUATION_PROCESS')
    ev=g.read(slot/'evaluation.json');result=g.evaluate(p,psha,x['records'],did,'qualification',ev);rejects=[]
    def rejected(name,fn):
        try:fn()
        except (g.Rejected,OSError):rejects.append(name)
        else:raise AssertionError(name+' was accepted')
    proof_sha=g.prove_rollback(p,psha,x['records'],did,'qualification',ev)
    accept={'schema':'history-acceptance/v1','owner':'fixture-owner',**result,'rollback_proof_sha':proof_sha,'previous_sha':None};ab=g.encode(accept)
    rejected('missing-independent-acceptance',lambda:g.promote(p,psha,x['records'],did,'qualification',ev,b'{}',g.digest(b'{}')))
    for name in ['tenant_leak','secret_leak','abuse_failed','candidate_loss','latency_ms','slice','id']:
        bad=g.parse(ev);bad['cases'][0][name]=True if name in ['tenant_leak','secret_leak','abuse_failed'] else 1000000 if name in ['candidate_loss','latency_ms'] else 'wrong';bb=g.encode(bad)
        rejected('reject-'+name,lambda bb=bb:g.promote(p,psha,x['records'],did,'qualification',bb,ab,g.digest(ab)))
    rejected('job-id-replay',lambda:jobs.train(p,psha,x['records'],did,'qualification',x['model'],runtime_raw,intent,g.digest(intent)))
    g.need(g.active(x['root'])[0] is None,'AUTOMATIC_PROMOTION')
    badproof={**accept,'rollback_proof_sha':'0'*64};bp=g.encode(badproof)
    rejected('missing-proven-rollback',lambda:g.promote(p,psha,x['records'],did,'qualification',ev,bp,g.digest(bp)))
    # Inject a local publication failure only in this synthetic qualification.
    real_publish=g.publish
    def interrupted_publish(path,b):
        if Path(path).name.startswith('version-'):raise OSError('synthetic publication failure')
        return real_publish(path,b)
    with patch.object(g,'publish',side_effect=interrupted_publish):
        rejected('promotion-publication-failure',lambda:g.promote(p,psha,x['records'],did,'qualification',ev,ab,g.digest(ab)))
    rejected('uncertain-promotion-not-retried',lambda:g.promote(p,psha,x['records'],did,'qualification',ev,ab,g.digest(ab)))
    g.need(g.active(x['root'])[0] is None,'FAILED_PROMOTION_CHANGED_ACTIVE')
    # Explicit test-operator reconciliation of the observed empty reservation.
    reservation=x['root']/'promotion-reserved';g.need(reservation.resolve().parent==x['root'].resolve() and not list(reservation.iterdir()),'RECONCILIATION_SCOPE');reservation.rmdir()
    target=g.promote(p,psha,x['records'],did,'qualification',ev,ab,g.digest(ab));g.need(g.active(x['root'])[0]==target,'PROMOTION')
    rollback=g.encode({'owner':'fixture-owner','tenant':x['profile']['tenant'],'target':None,'current':target});g.rollback(p,psha,x['records'],None,rollback,g.digest(rollback));g.need(g.active(x['root'])[0] is None,'ROLLBACK')
    g.need(g.read(folder/'original.jsonl')==x['source'].read_bytes(),'ORIGINAL_MUTATED')
    # Preserve/recover the exact audit inputs, not only their hashes.
    g.need(g.read(x['root']/('profile-'+psha+'.json'))==p,'PROFILE_AUDIT')
    g.need(g.read(x['root']/('acceptance-'+g.digest(ab)+'.json'))==ab,'ACCEPTANCE_AUDIT')
    g.need(g.read(x['root']/('evaluation-'+g.digest(ev)+'.json'))==ev,'EVALUATION_AUDIT')
    # Explicitly seed a synthetic pre-existing baseline pointer, then exercise
    # the same public promotion/rollback API against a non-null previous model.
    baseline=x['profile']['acceptance']['baseline_sha'];_,chain=g.active(x['root'])
    g.publish(x['root']/f'version-{len(chain):08d}.json',g.encode({'previous':None,'target':baseline,'tenant':x['profile']['tenant'],'action':'SYNTHETIC_PREEXISTING_BASELINE'}))
    proof2=g.prove_rollback(p,psha,x['records'],did,'qualification',ev);accept2={**accept,'previous_sha':baseline,'rollback_proof_sha':proof2};a2=g.encode(accept2)
    g.promote(p,psha,x['records'],did,'qualification',ev,a2,g.digest(a2));g.need(g.active(x['root'])[0]==target,'EXISTING_PROMOTION')
    rollback2=g.encode({'owner':'fixture-owner','tenant':x['profile']['tenant'],'target':baseline,'current':target});g.rollback(p,psha,x['records'],baseline,rollback2,g.digest(rollback2));g.need(g.active(x['root'])[0]==baseline and g.lock_tree(x['model'])[0]==baseline,'EXISTING_BASELINE_RESTORE')
    # Actual budget exhaustion of the same official worker, with a separately
    # bound profile/dataset/intent and deliberately insufficient finite time.
    limited=copy.deepcopy(x);limited['profile']=copy.deepcopy(x['profile']);limited['profile']['budget']['timeout']=1;limited['profile']['config']['steps']=1000
    limited_did=f.ingest(limited)['dataset'];limited_raw=g.encode(limited['profile']);limited_sha=g.digest(limited_raw);limited_intent=f.intent(limited,limited_did,'budget-exhaustion')
    timed=jobs.train(limited_raw,limited_sha,limited['records'],limited_did,'budget-exhaustion',limited['model'],runtime_raw,limited_intent,g.digest(limited_intent))
    g.need(timed['status']=='UNRESOLVED' and timed['capture']['status']=='TIMED_OUT' and timed['capture']['tree_empty'],'TIMEOUT_FAIL_CLOSED')
    rejected('timeout-not-replayed',lambda:jobs.train(limited_raw,limited_sha,limited['records'],limited_did,'budget-exhaustion',limited['model'],runtime_raw,limited_intent,g.digest(limited_intent)))
    g.need(g.active(x['root'])[0]==baseline,'TIMEOUT_PROMOTED')
    report={'pass':True,'scope':'synthetic library reference; no consumer model/corpus acceptance','dataset':did,'job':jr,'evaluation_sha':g.digest(ev),'evaluation':g.parse(ev),'promotion_rollback_proven':True,'existing_baseline_restored':True,'exact_audit_inputs_preserved':True,'negative_cases':rejects,'original_preserved':True,'training_excludes_eval':True,'timeout_job':timed,'seconds':time.monotonic()-began,'fixture_work':str(x['work'])};
    with a.output.open('xb') as out:out.write(g.encode(report))
    print('HISTORY_TRAINING_RUNTIME_PASS',len(rejects),flush=True)
if __name__=='__main__':main()
````

### FILE: `history_training/worker.py`
```yaml
block_id: "HISTORY-MODEL-TRAINING:worker:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration governed by ML production, privacy, execution and official pinned SDK contracts; not vendor-authored pipeline"
license: "LicenseRef-Workspace-Owner"
sha256: "b073a272f9719b030fc9089f50be82f6cdaedda8fe9422df49b9df3cac5cfaa8"
variables: []
secrets_allowed: false
```
````python
"""AUTHORED bounded adapter. All learning is performed by the exact official SDK."""
from pathlib import Path
import os,sys,time,json,hashlib,contextlib
sys.path.insert(0,str(Path(__file__).resolve().parent))
import governance as g

def execute(request_path,request_sha):
    raw=g.read(request_path);g.need(g.digest(raw)==request_sha,'REQUEST_HASH');r=g.parse(raw)
    g.exact(r,'schema train validation model model_sha config output');g.need(r['schema']=='history-sft-worker/v1','WORKER_SCHEMA')
    work=Path(request_path).parent;out=Path(r['output']);g.need(out.parent==work and not out.exists(),'OUTPUT_PATH')
    env={'HF_HUB_OFFLINE':'1','HF_DATASETS_OFFLINE':'1','HF_HUB_DISABLE_TELEMETRY':'1','HF_HUB_DISABLE_XET':'1','DO_NOT_TRACK':'1','WANDB_DISABLED':'true','TOKENIZERS_PARALLELISM':'false','HF_HOME':str(work/'hf'),'HF_DATASETS_CACHE':str(work/'cache'),'OMP_NUM_THREADS':'1','MKL_NUM_THREADS':'1'};os.environ.update(env)
    network=[]
    def audit(event,args):
        if event in ('socket.connect','socket.getaddrinfo','socket.sendto'):network.append(event);raise PermissionError('OFFLINE')
    sys.addaudithook(audit)
    rows={}
    for k in ('train','validation'):
        g.exact(r[k],'path sha256');b=g.read(r[k]['path']);g.need(g.digest(b)==r[k]['sha256'],'SPLIT_HASH');rows[k]=g.parse(b)
    model_sha,_=g.lock_tree(r['model']);g.need(model_sha==r['model_sha'],'MODEL_HASH')
    import torch
    from transformers import GPT2LMHeadModel,PreTrainedTokenizerFast
    from trl import SFTTrainer,SFTConfig
    from datasets import Dataset,disable_progress_bars
    from transformers.utils import logging as tlogging
    from huggingface_hub.utils import disable_progress_bars as disable_hub_bars
    disable_progress_bars();disable_hub_bars();tlogging.set_verbosity_error()
    torch.set_num_threads(1);torch.set_num_interop_threads(1);torch.manual_seed(r['config']['seed']);torch.use_deterministic_algorithms(True)
    model=GPT2LMHeadModel.from_pretrained(r['model'],local_files_only=True,use_safetensors=True);tok=PreTrainedTokenizerFast.from_pretrained(r['model'],local_files_only=True,trust_remote_code=False)
    g.need(sum(x.numel() for x in model.parameters())<=250000 and model.config.n_positions>=r['config']['max_length'],'MODEL_CAPACITY')
    def weights(m):
        h=hashlib.sha256()
        for n,p in sorted(m.state_dict().items()):h.update(n.encode());h.update(p.detach().cpu().contiguous().numpy().tobytes())
        return h.hexdigest()
    before=weights(model);cfg=r['config']
    # Fail instead of silently dropping any target through truncation.
    for values in rows.values():
        for item in values:
            ids=tok(item['prompt']+item['completion'])['input_ids'];prefix=tok(item['prompt'])['input_ids']
            g.need(0<len(prefix)<len(ids)<=cfg['max_length'] and ids[:len(prefix)]==prefix,'TOKEN_BOUNDARY_OR_TRUNCATION')
    args=SFTConfig(output_dir=str(work/'trainer'),use_cpu=True,bf16=False,fp16=False,max_steps=cfg['steps'],per_device_train_batch_size=2,per_device_eval_batch_size=1,gradient_accumulation_steps=1,learning_rate=cfg['learning_rate'],report_to='none',save_strategy='no',eval_strategy='no',logging_strategy='no',disable_tqdm=True,seed=cfg['seed'],data_seed=cfg['seed'],dataloader_num_workers=0,dataloader_pin_memory=False,max_length=cfg['max_length'],packing=False,padding_free=False,completion_only_loss=True,optim='adamw_torch',push_to_hub=False)
    ds=lambda values:Dataset.from_list([{'prompt':x['prompt'],'completion':x['completion']} for x in values])
    trainer=SFTTrainer(model=model,args=args,train_dataset=ds(rows['train']),eval_dataset=ds(rows['validation']),processing_class=tok)
    g.need(len(trainer.train_dataset)==len(rows['train']) and len(trainer.eval_dataset)==len(rows['validation']),'DROPPED_ROWS')
    batch=next(iter(trainer.get_train_dataloader()));g.need((batch['labels']==-100).any().item() and (batch['labels']!=-100).any().item(),'LABEL_MASK')
    start=time.monotonic();result=trainer.train();seconds=time.monotonic()-start;after=weights(model);g.need(before!=after and result.global_step==cfg['steps'],'NO_TRAINING')
    model.save_pretrained(out,safe_serialization=True);tok.save_pretrained(out)
    loaded=GPT2LMHeadModel.from_pretrained(out,local_files_only=True,use_safetensors=True);g.need(weights(loaded)==after and not network,'RELOAD_OR_EGRESS')
    value={'schema':'history-worker-result/v1','status':'TRAINED','steps':result.global_step,'training_seconds':seconds,'before_weights':before,'after_weights':after,'network_attempts':len(network),'reload_exact':True}
    g.publish(work/'worker-result.json',g.encode(value))

if __name__=='__main__':
    try:
        # Upstream progress/errors may contain private text: discard at source.
        with open(os.devnull,'w') as sink,contextlib.redirect_stdout(sink),contextlib.redirect_stderr(sink):execute(sys.argv[1],sys.argv[2])
    except Exception:
        print('HISTORY_TRAINING_JOB_REJECTED',file=sys.stderr);sys.exit(2)
    print('HISTORY_TRAINING_JOB_COMPLETED')
````

## 6. Configuration surface

No secrets in the pack. NEW false is the only runnable default. An activated profile requires tenant, source/review/model/runtime/code hashes, ten separately supplied decision records, private root/owner/retention, time splits, SFT configuration, finite budgets, baseline/owner/reviewer and per-case thresholds. Exact schema rejects unknown/missing fields and nonfinite values before data access. Defaults in synthetic fixtures are explicitly test-only. A target profile requires independently selected business/model/budget/security decisions.

## 7. Dependency bill

58exact wheel pins, hashes and byte sizes in runtime-wheel-lock.json and requirements.lock; acquisition URLs in source-lock.json. TRL1.12.0, Transformers5.17.0, Torch2.14.0 CPU, Tokenizers0.23.2, Datasets5.0.1, Safetensors0.9.0rc0. CPython3.14 host/pip26.0.1 and four native supervisor files are separate admitted prerequisites. Read THIRD_PARTY_NOTICES and V373 for licenses, maintenance finding and native coverage limits. No third-party binaries are materialized by this pack.

## 8. Apply order

Compose the opt-in plan into an absent project/test destination; reconstruct four unchanged supervisor files and this pack. Run policy suite. Acquire exact wheels with governed quarantine profile, then qualify the isolated runtime and generate a synthetic fixture model. Run actual reference harness. A consumer replaces fixture decisions with its own verified owners/records and requalifies its selected target; bootstrap/SDK probes never authorize consumer history. Preserve failed destinations; no in-place upgrades or automatic uncertain job retries.

## 9. Verification

Policy suite, native ACL/read/publish semantics, explicit profile/input rejection, concurrent/partial import, tenant/split/lineage and no-training NEW. Actual synthetic official SFT changes/reloads weights, independent baseline/candidate evaluation, all case/slice thresholds, failed promotions, explicit acceptance, rollback probe before promotion and actual subsequent rollback. Reconstruct from canonical Markdown and rerun. Exact58wheel install/pip check and installed-byte manifest; source/OSV qualification kept separate from full security. Benchmark records bounded small CPU reference, memory and wall time; no consumer throughput/quality promise.

## 10. Reconstruction evidence

reconstruction_evidence/HISTORY_TRAINING_SDK_QUALIFICATION_V373.md records prior narrow SDK admission. reconstruction_evidence/HISTORY_MODEL_TRAINING_CONTROL_V373.md records final reconstruction, tests, runtime/bootstrap/evaluation hashes and the TEST09 decision. Fresh21file composition,27policy tests,58wheel bootstrap and actual SFT/evaluation/promotion/rollback/timeout have passed; consumer conditions remain mandatory.
