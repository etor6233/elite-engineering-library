"""Finite, synthetic HTTP -> official metrics -> canonical alert qualification."""
from pathlib import Path
import argparse,subprocess,os,json,hashlib,socket,urllib.request,urllib.error,urllib.parse,ssl,time,uuid,csv,re
p=argparse.ArgumentParser()
for n in ['runtime-lock','runtime-lock-sha256','consumer','host','host-sha256','rules','rules-sha256','work-root']:p.add_argument('--'+n,required=True)
a=p.parse_args();sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
lock_path=Path(a.runtime_lock);assert sha(lock_path)==a.runtime_lock_sha256
lock=json.loads(lock_path.read_text());bins={}
for n in ['postgres','initdb','pg_ctl','psql','prometheus','promtool','telemetry-reference']:
 pin=lock['binaries'][n];q=Path(pin['path']);assert q.is_file() and q.stat().st_size==pin['bytes'] and sha(q)==pin['sha256'];bins[n]=q
host=Path(a.host);rules=Path(a.rules);assert sha(host)==a.host_sha256 and sha(rules)==a.rules_sha256
consumer=Path(a.consumer);migrations=sorted((consumer/'db/migrations').glob('*.up.sql'))
assert {x.name:sha(x) for x in migrations}==lock['migrations']
root=Path(a.work_root).resolve(strict=True);W=root/('http-'+uuid.uuid4().hex);W.mkdir()
system=Path(os.environ['SYSTEMROOT'])/'System32'
sid=list(csv.reader(subprocess.check_output([str(system/'whoami.exe'),'/user','/fo','csv','/nh'],text=True).splitlines()))[0][1]
assert re.fullmatch(r'S-[0-9-]+',sid)
subprocess.run([str(system/'icacls.exe'),str(W),'/inheritance:r','/grant:r','*'+sid+':(OI)(CI)F','*S-1-5-18:(OI)(CI)F'],check=True,stdout=subprocess.DEVNULL,creationflags=subprocess.CREATE_NO_WINDOW)
print('HTTP_REFERENCE_STAGE '+str(W),flush=True)
env={k:os.environ[k] for k in ['SYSTEMROOT','WINDIR']};env.update(TEMP=str(W),TMP=str(W),PGCONNECT_TIMEOUT='3',COMSPEC=str(system/'cmd.exe'))
procs={};handles=[];events=[];result={'state':'RUNNING','rules_sha256':a.rules_sha256,'host_sha256':a.host_sha256,'requests':{},'synthetic_only':True}
start=time.monotonic();deadline=start+1140
def log(kind,**values):
 events.append({'seconds':round(time.monotonic()-start,3),'kind':kind,**values})
 (W/'progress.json').write_text(json.dumps(events,indent=2)+'\n',encoding='utf-8',newline='\n')
 print(json.dumps(events[-1]),flush=True)
def command(name,args,seconds=30):
 q=subprocess.run([str(v) for v in args],env=env,stdin=subprocess.DEVNULL,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=seconds,creationflags=subprocess.CREATE_NO_WINDOW)
 (W/(name+'.log')).write_bytes(q.stdout)
 (W/(name+'.process.json')).write_text(json.dumps({'exit_code':q.returncode,'output_bytes':len(q.stdout)}),encoding='utf-8')
 if q.returncode:raise RuntimeError('reference command failed: '+name+' exit='+str(q.returncode))
 return q.stdout.decode().strip()
def port():
 with socket.socket() as s:s.bind(('127.0.0.1',0));return s.getsockname()[1]
pgport=port();promport=port();DATA=W/'pgdata';DB='elite_http_test_'+uuid.uuid4().hex
def launch(name,args,extra=None):
 f=(W/(name+'.log')).open('xb');handles.append(f)
 q=subprocess.Popen([str(v) for v in args],cwd=W,env=env|dict(extra or {}),stdin=subprocess.DEVNULL,stdout=f,stderr=subprocess.STDOUT,creationflags=subprocess.CREATE_NO_WINDOW);procs[name]=q;return q
def wait(fn,label,seconds=30):
 end=min(time.monotonic()+seconds,deadline)
 while time.monotonic()<end:
  try:
   value=fn()
   if value:return value
  except (OSError,ValueError,urllib.error.URLError):pass
  time.sleep(.25)
 raise TimeoutError(label)
def sql(query,db=None):
 return command('sql-'+uuid.uuid4().hex,[bins['psql'],'-X','-w','-qAt','-h','127.0.0.1','-p',pgport,'-U','postgres','-d',db or DB,'-v','ON_ERROR_STOP=1','-c',query],20)
def pg_ready():
 if procs['postgres'].poll() is not None:raise RuntimeError('owned PostgreSQL exited')
 try:observed=sql('show data_directory','postgres')
 except RuntimeError:return False
 assert Path(observed).resolve()==DATA.resolve()
 return True
def get(url,context=None):
 with urllib.request.urlopen(url,context=context,timeout=5) as r:return r.read(2_000_001)
def api(path):return json.loads(get('https://127.0.0.1:'+str(promport)+path,tls))
def alert():
 return [x for x in api('/api/v1/alerts')['data']['alerts'] if x['labels']['alertname']=='PlatformHighErrorRate']
def firing():return any(x['state']=='firing' for x in alert())
def request(key,token='synthetic-local-reference-only',body=None):
 data=body or b'{"OrganizationID":"integration-org","Currency":"USD","TotalMinorUnits":100}'
 req=urllib.request.Request('http://'+address['api']+'/v1/orders?contact=PRIVATE_QUERY@example.invalid',data=data,headers={'Authorization':'Bearer '+token,'Content-Type':'application/json','Idempotency-Key':key,'Cookie':'PRIVATE_COOKIE','X-Private':'PRIVATE_HEADER'},method='POST')
 try:
  with urllib.request.urlopen(req,timeout=15) as r:status=r.status;payload=r.read(16384)
 except urllib.error.HTTPError as e:status=e.code;payload=e.read(16384)
 result['requests'][str(status)]=result['requests'].get(str(status),0)+1
 assert b'PRIVATE_' not in payload
 return status,payload
def verify_metrics(label):
 data=get('http://'+address['metrics']+'/metrics')
 assert len(data)<=2_000_000
 for token in [b'PRIVATE_',b'integration-org',b'synthetic-local-reference-only',b'customer_order',b'Authorization',b'http_route',b'server_address',b'url_',b'target_info',b'otel_scope']:
  assert token not in data,token
 (W/(label+'.prom')).write_bytes(data)
 return data
try:
 command('certificates',[bins['telemetry-reference'],'--certificates',W])
 tls=ssl.create_default_context(cafile=str(W/'ca.pem'));tls.minimum_version=ssl.TLSVersion.TLSv1_3;tls.load_cert_chain(str(W/'client.pem'),str(W/'client.key'))
 command('initdb',[bins['initdb'],'-D',DATA,'-U','postgres','-A','trust','--no-locale','-E','UTF8'])
 with (DATA/'postgresql.conf').open('a') as f:f.write(f"\nlisten_addresses='127.0.0.1'\nport={pgport}\nfsync=on\nsynchronous_commit=on\nfull_page_writes=on\n")
 launch('postgres',[bins['postgres'],'-D',DATA]);wait(pg_ready,'owned postgres')
 sql('create database '+DB,'postgres')
 for m in migrations:command('migration-'+m.stem,[bins['psql'],'-X','-w','-qAt','-h','127.0.0.1','-p',pgport,'-U','postgres','-d',DB,'-v','ON_ERROR_STOP=1','-f',m])
 sql("insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values('018f4d4a-7b36-7a21-8d10-2f4c54c28b01','http-reference','Synthetic Local Only','Synthetic Local Only'); insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values('018f4d4a-7b36-7a21-8d10-2f4c54c28b01','integration-org','reference','Synthetic Local Only','enterprise')")
 launch('host',[host,'--ttl','20m'],{'REFERENCE_DATABASE_URL':f'postgres://postgres@127.0.0.1:{pgport}/{DB}?sslmode=disable'})
 def host_ready():
  if procs['host'].poll() is not None:raise RuntimeError('reference host exited')
  text=(W/'host.log').read_text()
  if '\n' not in text:return None
  return json.loads(text.splitlines()[0])
 address=wait(host_ready,'host startup')
 assert set(address)=={'api','metrics'} and all(re.fullmatch(r'127\.0\.0\.1:[0-9]{1,5}',v) for v in address.values())
 first=request('reference-http-order-0001');assert first[0]==201,first
 again=request('reference-http-order-0001');assert again[0]==200 and json.loads(first[1])==json.loads(again[1])
 assert request('reference-http-unauth-0001','PRIVATE_TOKEN')[0]==401
 assert request('reference-http-invalid-0001',body=b'{"PRIVATE_BODY":true}')[0]==400
 assert sql("select count(*) from sales.customer_order")=='1'
 verify_metrics('healthy')
 (W/'prometheus.json').write_text(json.dumps({'global':{'scrape_interval':'1s','evaluation_interval':'30s'},'rule_files':[str(rules)],'scrape_configs':[{'job_name':'reference-http','static_configs':[{'targets':[address['metrics']]}]}]},indent=2))
 (W/'web.json').write_text(json.dumps({'tls_server_config':{'cert_file':str(W/'server.pem'),'key_file':str(W/'server.key'),'client_ca_file':str(W/'ca.pem'),'client_auth_type':'RequireAndVerifyClientCert','min_version':'TLS13'}}))
 command('rules-check',[bins['promtool'],'check','config',W/'prometheus.json'])
 launch('prometheus',[bins['prometheus'],'--config.file='+str(W/'prometheus.json'),'--web.listen-address=127.0.0.1:'+str(promport),'--web.config.file='+str(W/'web.json'),'--storage.tsdb.path='+str(W/'tsdb'),'--storage.tsdb.retention.time=1h','--storage.tsdb.retention.size=64MB','--log.level=error'])
 wait(lambda:api('/api/v1/status/config')['status']=='success','prometheus startup')
 wait(lambda:api('/api/v1/query?query='+urllib.parse.quote('up{job="reference-http"}'))['data']['result'],'first scrape')
 time.sleep(3)
 sql('alter table sales.customer_order rename to customer_order_reference_fault')
 fault=time.monotonic();next_error=fault;pending=None;fired=None;last_state=None;error_count=0
 log('fault_injected',method='rename owned synthetic order table',interval_seconds=240,rule_hold_seconds=600)
 while time.monotonic()<deadline-100:
  if time.monotonic()>=next_error:
   assert request('reference-http-failure-'+str(error_count).zfill(5))[0]==500
   error_count+=1;next_error+=240;verify_metrics('fault-'+str(error_count))
   log('actual_http_500',count=error_count)
  active=alert();state=active[0]['state'] if active else 'inactive'
  if state!=last_state:log('alert_state',state=state);last_state=state
  if state=='pending' and pending is None:pending=time.monotonic()
  if state=='firing':
   assert pending is not None and time.monotonic()-pending>=590
   fired=time.monotonic();result['firing_alert']=active;break
  assert procs['host'].poll() is None and procs['prometheus'].poll() is None
  time.sleep(1)
 assert fired is not None,'alert did not fire within finite budget'
 expr='sum(rate(http_server_request_duration_seconds_count{http_response_status_code=~"5.."}[5m])) / sum(rate(http_server_request_duration_seconds_count[5m]))'
 current=api('/api/v1/query?query='+urllib.parse.quote(expr));legacy=api('/api/v1/query?query='+urllib.parse.quote(expr.replace('/ sum(rate(http_server_request_duration_seconds_count[5m]))','/ clamp_min(sum(rate(http_server_request_duration_seconds_count[5m])), 1)')))
 result['observed_ratios']={'corrected':current,'legacy':legacy}
 assert float(current['data']['result'][0]['value'][1])>0.99
 assert float(legacy['data']['result'][0]['value'][1])<0.01, 'legacy window ratio was not below threshold'
 sql('alter table sales.customer_order_reference_fault rename to customer_order')
 for i in range(500):assert request('reference-http-order-0001')[0]==200
 wait(lambda:not firing(),'alert recovery',70)
 assert sql('select count(*) from sales.customer_order')=='1'
 assert sql('select count(*) from platform.idempotency_record')=='1'
 assert sql("select count(*) from platform.outbox_event where event_type='order.created'")=='1'
 verify_metrics('recovered')
 result.update(state='PASS',actual_sparse_failures=error_count,pending_to_firing_seconds=round(fired-pending,3),fault_to_firing_seconds=round(fired-fault,3),same_order_after_recovery=True,order_count=1,outbox_count=1,idempotency_count=1)
 log('recovered',order_count=1,outbox_count=1,idempotency_count=1)
except BaseException as e:
 result.update(state='FAIL',failure_type=type(e).__name__,failure=str(e))
 raise
finally:
 if 'host' in procs and procs['host'].poll() is None:procs['host'].terminate();procs['host'].wait(10)
 if 'prometheus' in procs and procs['prometheus'].poll() is None:procs['prometheus'].terminate();procs['prometheus'].wait(10)
 if 'postgres' in procs and procs['postgres'].poll() is None:command('postgres-stop',[bins['pg_ctl'],'-D',DATA,'-m','fast','-w','stop']);procs['postgres'].wait(10)
 for f in handles:f.close()
 result['children_exited']=all(q.poll() is not None for q in procs.values());result['elapsed_seconds']=round(time.monotonic()-start,3)
 result['event_log_sha256']=sha(W/'progress.json') if (W/'progress.json').exists() else None
 (W/'result.json').write_text(json.dumps(result,indent=2)+'\n',encoding='utf-8',newline='\n')
 print('HTTP_REFERENCE_RESULT '+json.dumps(result),flush=True)
