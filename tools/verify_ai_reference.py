"""AUTHORED bounded local fixture verifier; never opens a provider account."""
from pathlib import Path
import argparse,json,os,subprocess,sys,urllib.parse,hashlib

def main():
 p=argparse.ArgumentParser();p.add_argument('--go',type=Path,required=True);p.add_argument('--database-url',required=True);p.add_argument('--receipt',type=Path,required=True);a=p.parse_args()
 url=urllib.parse.urlparse(a.database_url)
 if url.scheme not in ('postgres','postgresql')or url.hostname!='127.0.0.1'or not url.path.startswith('/elite_payment_connected_')or a.receipt.exists():raise ValueError('owned loopback fixture database and absent receipt required')
 root=Path(__file__).resolve().parents[1];env=os.environ.copy()
 for k in list(env):
  if k.startswith(('ELITE_','PG','LLM_','OPENAI_'))or k in ('GOROOT','DATABASE_URL','TEST_DATABASE_URL'):env.pop(k)
 env.update(PAYMENT_CONNECTED_DB_URL=a.database_url,TEST_DATABASE_URL=a.database_url,GOWORK='off',GOPROXY='off',GOSUMDB='off',GOTOOLCHAIN='local',GOFLAGS='-mod=readonly',GOMAXPROCS='4')
 rows=[]
 for label,args,marker in [
  ('connected',['./internal/platform/postgres','-run','^TestAIConnectedReferenceEvaluation$'],'AI_CONNECTED_EVALUATION_PASS'),
  ('eval-contract',['./internal/aifoundation','-run','^TestEval'],''),
  ('provider-domain-contract',['./internal/domainbind','./internal/llmopenai','./internal/conversationruntime'],'')]:
  q=subprocess.run([str(a.go),'test','-mod=readonly',*args,'-count=1','-v','-timeout','120s'],cwd=root,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=150,creationflags=subprocess.CREATE_NO_WINDOW if os.name=='nt'else 0)
  log=a.receipt.with_name(a.receipt.name+'.'+label+'.log')
  with log.open('xb')as f:f.write(q.stdout)
  passed=q.returncode==0 and b'--- SKIP:'not in q.stdout and (not marker or marker.encode()in q.stdout)
  rows.append({'name':label,'exit_code':q.returncode,'passed':passed,'log':str(log),'sha256':hashlib.sha256(q.stdout).hexdigest()})
  if not passed:break
 result={'scope':'LOCAL_SYNTHETIC_MODEL_REAL_DOMAIN','state':'PASS'if len(rows)==3 and all(r['passed']for r in rows)else'FAIL','provider_calls_live':False,'runs':rows}
 with a.receipt.open('x',encoding='utf-8',newline='\n')as f:json.dump(result,f,indent=2);f.write('\n')
 print(result['state']);return 0 if result['state']=='PASS'else 2
if __name__=='__main__':sys.exit(main())
