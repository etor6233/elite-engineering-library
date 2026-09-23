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
