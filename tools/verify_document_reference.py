"""AUTHORED finite local fixture runner; does not provision or address live services."""
from pathlib import Path
import argparse,hashlib,json,os,subprocess,time
from urllib.parse import urlsplit
p=argparse.ArgumentParser();p.add_argument('--go',type=Path,required=True);p.add_argument('--fixture',type=Path,required=True);p.add_argument('--receipt',type=Path,required=True);args=p.parse_args()
root=Path(__file__).resolve().parents[1]
sha=lambda path:hashlib.sha256(path.read_bytes()).hexdigest()
if sha(args.go)!='21761eceb9302062c9623fb699f332c8c7fe000f15f70efe8da01a2cfbbc16b9':raise SystemExit('exact admitted Go1.26.8 executable required')
if sha(args.fixture)!='489f0c63b6a05e0ae0fcfbf29121299fc2ccacab35d5cd40a5d63d42764078cb':raise SystemExit('exact public fixture required')
url=urlsplit(os.environ.get('DOCUMENT_CONNECTED_DB_URL',''))
if url.scheme not in ('postgres','postgresql') or url.hostname not in ('127.0.0.1','localhost') or not url.path.startswith(('/elite_document_reference_','/elite_payment_connected_')):raise SystemExit('owned loopback synthetic document database required; apply selected migrations first')
if args.receipt.exists() or not args.receipt.parent.is_dir():raise SystemExit('receipt must be new with existing parent')
env=os.environ.copy();env.pop('GOROOT',None);env.update(GOTOOLCHAIN='local',GOPROXY='off',GOSUMDB='off',GOWORK='off',GOFLAGS='-mod=readonly',GOMAXPROCS='4',DOCUMENT_TEST_ROOT=str(root),ELITE_DOCUMENT_REFERENCE_FIXTURE=str(args.fixture.resolve()))
steps=[('connected',['test','./internal/platform/postgres','-run','^TestDocument','-count=1','-v','-timeout','120s']),('host',['test','./cmd/electromobility-api','-run','^TestDocumentHostActivation$','-count=1','-v','-timeout','60s']),('bounds',['test','./internal/documentbridge','./internal/approval','-run','^(TestDocument|FuzzDocument)','-count=1','-v']),('fuzz',['test','./internal/documentbridge','-run','^$','-fuzz','^FuzzDocumentOriginalBoundary$','-fuzztime','3s','-parallel','2']),('vet',['vet','./internal/documentbridge','./internal/approval','./internal/platform/postgres','./internal/platform/httpapi','./cmd/electromobility-api'])]
result={'schema':'elite-document-reference-verification/v1','state':'RUNNING','scope':'LOCAL_FIXTURES','live_effects':False,'detectors':'SIMULATED; no scanner effectiveness proof','go_sha256':sha(args.go),'fixture_sha256':sha(args.fixture),'steps':[]};start=time.monotonic()
try:
 for label,argv in steps:
  log=args.receipt.with_name(args.receipt.stem+'-'+label+'.log')
  with log.open('xb')as f:q=subprocess.run([str(args.go),*argv],cwd=root,env=env,stdout=f,stderr=subprocess.STDOUT,timeout=150,creationflags=getattr(subprocess,'CREATE_NO_WINDOW',0))
  raw=log.read_bytes();result['steps'].append({'name':label,'exit_code':q.returncode,'log_sha256':sha(log),'log':str(log),'no_skip':b'--- SKIP:' not in raw})
  if q.returncode or b'--- SKIP:'in raw:raise RuntimeError(label+' failed; consult preserved log')
 result['state']='PASS'
except BaseException as e:
 result.update(state='FAIL',failure_type=type(e).__name__);raise
finally:
 result['seconds']=round(time.monotonic()-start,3);args.receipt.write_text(json.dumps(result,indent=2)+'\n',encoding='utf-8',newline='\n')
