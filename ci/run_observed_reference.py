"""AUTHORED finite local API/Next -> official Prometheus -> fault/recovery probe.

Only a caller-selected, owned synthetic PostgreSQL reference is accepted. Every
binary, artifact, rule and configuration is bound to an external SHA. No live
provider, notification recipient, production identity or SLO is inferred.
"""
from pathlib import Path
import argparse,csv,hashlib,json,os,queue,re,socket,ssl,subprocess,threading,time,uuid
from concurrent.futures import ThreadPoolExecutor
from urllib.request import Request,build_opener,ProxyHandler,HTTPSHandler
from urllib.error import HTTPError
from urllib.parse import quote,urlsplit
from local_release import Deployment,verify_release,runtime_config
from build_local_reference import plain,read_json,require,sha,raw_json,environment

def main():
    p=argparse.ArgumentParser()
    for n in ['artifact','inventory','inventory-sha256','runtime-lock','runtime-lock-sha256','database-config','database-config-sha256','rules','rules-sha256','work']:
        p.add_argument('--'+n,required=True)
    a=p.parse_args();W=plain(a.work);require(not W.exists(),'work directory must be absent');W.mkdir()
    release=verify_release(a.artifact,a.inventory,a.inventory_sha256);A=Path(release['root']);source=Path(__file__).resolve().parent
    lock=read_json(a.runtime_lock,a.runtime_lock_sha256);bins={}
    for name in ['postgres','pg_ctl','psql','prometheus','promtool','telemetry-reference']:
        pin=lock['binaries'][name];v=plain(pin['path']);require(v.is_file()and v.stat().st_size==pin['bytes']and sha(v)==pin['sha256'],'runtime changed: '+name);bins[name]=v
    conf=read_json(a.database_config,a.database_config_sha256);require(set(conf)=={'database_url','data_path'},'database configuration fields')
    data=plain(conf['data_path']);u=urlsplit(conf['database_url']);db=u.path.removeprefix('/')
    runtime_config({'database_url':conf['database_url'],'issuer':'http://127.0.0.1:6544','api_port':6545,'web_port':6546})
    require(data.is_dir()and not(data/'postmaster.pid').exists(),'owned database must be stopped')
    rules=plain(a.rules);require(sha(rules)==a.rules_sha256,'rules changed')
    system=Path(os.environ['SYSTEMROOT'])/'System32'
    sid=list(csv.reader(subprocess.check_output([str(system/'whoami.exe'),'/user','/fo','csv','/nh'],text=True).splitlines()))[0][1]
    require(re.fullmatch(r'S-[0-9-]+',sid),'local identity unresolved')
    subprocess.run([str(system/'icacls.exe'),str(W),'/inheritance:r','/grant:r','*'+sid+':(OI)(CI)F','*S-1-5-18:(OI)(CI)F'],check=True,stdout=subprocess.DEVNULL,creationflags=subprocess.CREATE_NO_WINDOW)
    env=environment(A/'node.exe',W);env.update(PGCONNECT_TIMEOUT='3',COMSPEC=str(system/'cmd.exe'))
    procs={};handles=[];dep=None;renamed=False;events=[];start=time.monotonic();deadline=start+1170
    result={'state':'FAIL','scope':'LOCAL_FIXTURES','accounts_used':False,'live_effects':False,'requests':{},'source_sha256':release['source_sha256'],'host_sha256':sha(A/'api/electromobility-api.exe'),'rules_sha256':a.rules_sha256,'runtime_lock_sha256':a.runtime_lock_sha256}
    def log(kind,**values):
        events.append({'seconds':round(time.monotonic()-start,3),'kind':kind,**values});(W/'progress.json').write_bytes(raw_json(events));print(json.dumps(events[-1]),flush=True)
    def command(name,args,seconds=30):
        q=subprocess.run(list(map(str,args)),env=env,stdin=subprocess.DEVNULL,capture_output=True,timeout=seconds,creationflags=subprocess.CREATE_NO_WINDOW)
        (W/(name+'.log')).write_bytes(q.stdout+q.stderr);require(q.returncode==0,'reference command failed: '+name);return q.stdout.decode().strip()
    def sql(query,database=db):
        return command('sql-'+uuid.uuid4().hex,[bins['psql'],'-X','-w','-qAt','-h','127.0.0.1','-p',u.port,'-U','postgres','-d',database,'-v','ON_ERROR_STOP=1','-c',query],15)
    def launch(name,args):
        f=(W/(name+'-'+uuid.uuid4().hex+'.log')).open('xb');handles.append(f)
        q=subprocess.Popen(list(map(str,args)),cwd=W,env=env,stdin=subprocess.DEVNULL,stdout=f,stderr=subprocess.STDOUT,creationflags=subprocess.CREATE_NO_WINDOW);procs[name]=q;return q
    def wait(fn,label,seconds=30):
        end=min(time.monotonic()+seconds,deadline)
        while time.monotonic()<end:
            try:
                value=fn()
                if value:return value
            except (OSError,ValueError):pass
            time.sleep(.25)
        raise TimeoutError(label)
    def port():
        with socket.socket()as s:s.bind(('127.0.0.1',0));return s.getsockname()[1]
    plain_client=build_opener(ProxyHandler({}))
    def get(url,client=plain_client):
        with client.open(url,timeout=5)as r:
            raw=r.read(2_000_001);require(len(raw)<=2_000_000,'response budget exceeded');return raw
    try:
        launch('postgres',[bins['postgres'],'-D',data])
        def pg_ready():
            require(procs['postgres'].poll()is None,'owned PostgreSQL exited')
            return Path(sql('show data_directory','postgres')).resolve()==data.resolve()
        wait(pg_ready,'owned PostgreSQL')
        require(sql('select count(*) from library_delivery.migration')=='84','current migration ledger required')
        sql("insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values('b0000000-0000-4000-8000-000000000001','delivery-local','Delivery Fixture','Delivery Fixture')on conflict do nothing; insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values('b0000000-0000-4000-8000-000000000001','b0000000-0000-4000-8000-000000000002','reference','Reference','store')on conflict do nothing")
        f=(W/'fixture-stderr.log').open('xb');handles.append(f)
        fixture=subprocess.Popen([A/'node.exe',source/'local_identity_fixture.cjs','--ttl-seconds','1200'],cwd=W,env=env,stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=f,creationflags=subprocess.CREATE_NO_WINDOW);procs['fixture']=fixture
        ready=queue.Queue();threading.Thread(target=lambda:ready.put(fixture.stdout.readline()),daemon=True).start();issuer=json.loads(ready.get(timeout=10))['issuer']
        cfg={'database_url':conf['database_url'],'issuer':issuer,'api_port':port(),'web_port':port(),'metrics_port':port()};runtime_config(cfg)
        dep=Deployment(W/'host',cfg);result['activation']=dep.activate(A,a.inventory,a.inventory_sha256)
        key='local-observability-'+uuid.uuid4().hex;token=''
        body=raw_json({'OrganizationID':'b0000000-0000-4000-8000-000000000002','Currency':'ARS','TotalMinorUnits':1900})
        def refresh():
            nonlocal token
            token=json.loads(get(issuer+'/fixture-token'))['token']
        def request(k=key,credential=None,payload=body):
            req=Request(f"http://127.0.0.1:{cfg['api_port']}/v1/orders?contact=PRIVATE_QUERY@example.invalid",data=payload,
                headers={'Authorization':'Bearer '+(token if credential is None else credential),'Content-Type':'application/json','Idempotency-Key':k,'Cookie':'PRIVATE_COOKIE','X-Private':'PRIVATE_HEADER'},method='POST')
            then=time.perf_counter()
            try:
                with build_opener(ProxyHandler({})).open(req,timeout=8)as r:status=r.status;raw=r.read(16385)
            except HTTPError as e:status=e.code;raw=e.read(16385)
            require(len(raw)<=16384 and b'PRIVATE_'not in raw,'response privacy/budget');return status,raw,time.perf_counter()-then
        def record(r):
            status,raw,seconds=r;result['requests'][str(status)]=result['requests'].get(str(status),0)+1;return r
        refresh();first=record(request());require(first[0]==201,'first order');order=json.loads(first[1])
        again=record(request());require(again[0]==200 and json.loads(again[1])==order,'order replay')
        require(record(request(key+'-auth','PRIVATE_TOKEN'))[0]==401,'identity fail-closed')
        require(record(request(key+'-body',payload=b'{"PRIVATE_BODY":true}'))[0]==400,'invalid body')
        def metrics(label):
            raw=get(f"http://127.0.0.1:{cfg['metrics_port']}/metrics")
            for forbidden in [b'PRIVATE_',b'b0000000',b'local-delivery-customer',b'customer_order',b'http_route',b'server_address',b'url_',b'target_info',b'otel_scope',token.encode()]:require(forbidden not in raw,'metric attribute escaped closed profile')
            (W/(label+'.prom')).write_bytes(raw);return raw
        metrics('healthy');command('certificates',[bins['telemetry-reference'],'--certificates',W])
        tls=ssl.create_default_context(cafile=str(W/'ca.pem'));tls.minimum_version=ssl.TLSVersion.TLSv1_3;tls.load_cert_chain(str(W/'client.pem'),str(W/'client.key'))
        client=build_opener(ProxyHandler({}),HTTPSHandler(context=tls));promport=port()
        def api(path):return json.loads(get('https://127.0.0.1:'+str(promport)+path,client))
        def query(expr,at=None):return api('/api/v1/query?query='+quote(expr)+(''if at is None else '&time='+str(at)))
        def alerts():return [v for v in api('/api/v1/alerts')['data']['alerts']if v['labels']['alertname']=='PlatformHighErrorRate']
        (W/'prometheus.json').write_bytes(raw_json({'global':{'scrape_interval':'1s','evaluation_interval':'30s'},'rule_files':[str(rules)],'scrape_configs':[{'job_name':'reference-http','static_configs':[{'targets':[f"127.0.0.1:{cfg['metrics_port']}"]}]}]}))
        (W/'web.json').write_bytes(raw_json({'tls_server_config':{'cert_file':str(W/'server.pem'),'key_file':str(W/'server.key'),'client_ca_file':str(W/'ca.pem'),'client_auth_type':'RequireAndVerifyClientCert','min_version':'TLS13'}}))
        command('rules-check',[bins['promtool'],'check','config',W/'prometheus.json'])
        promargs=[bins['prometheus'],'--config.file='+str(W/'prometheus.json'),'--web.listen-address=127.0.0.1:'+str(promport),'--web.config.file='+str(W/'web.json'),'--storage.tsdb.path='+str(W/'tsdb'),'--storage.tsdb.retention.time=1h','--storage.tsdb.retention.size=64MB','--log.level=error']
        launch('prometheus',promargs);wait(lambda:query('up{job="reference-http"}')['data']['result'],'first scrape');time.sleep(3)
        sql('alter table sales.customer_order rename to customer_order_reference_fault');renamed=True
        pending=fired=None;last=None;count=0;next_error=time.monotonic();log('fault_injected',rule_hold_seconds=600)
        while time.monotonic()<deadline-130:
            require(dep.instance.thread.is_alive()and procs['prometheus'].poll()is None,'observed host stopped')
            if time.monotonic()>=next_error:
                refresh();require(record(request(key+'-failure-'+str(count)))[0]==500,'actual database fault did not fail');count+=1;next_error+=240;metrics('fault-'+str(count));log('actual_http_500',count=count)
            active=alerts();state=active[0]['state']if active else 'inactive'
            if state!=last:log('alert_state',state=state);last=state
            if state=='pending'and pending is None:pending=time.monotonic()
            if state=='firing':
                fired=time.monotonic();require(pending is not None and fired-pending>=590,'hold shortened');result['firing_alert']=active;break
            time.sleep(1)
        require(fired is not None,'alert did not fire within finite budget')
        sql('alter table sales.customer_order_reference_fault rename to customer_order');renamed=False;refresh()
        started=time.monotonic()
        with ThreadPoolExecutor(max_workers=8)as pool:load=list(pool.map(lambda _:request(),range(1000)))
        for r in load:record(r);require(r[0]==200 and json.loads(r[1])==order,'concurrent replay diverged')
        timings=sorted(r[2]for r in load);elapsed=time.monotonic()-started;require(elapsed<90,'load exceeded finite local budget')
        result['load']={'requests':1000,'concurrency':8,'seconds':round(elapsed,3),'p95_ms':round(timings[949]*1000,3),'max_ms':round(timings[-1]*1000,3),'all_same_order':True,'production_slo_claim':False}
        wait(lambda:not any(x['state']=='firing'for x in alerts()),'alert recovery',70);metrics('recovered')
        fixed=time.time()-5;expr='sum(http_server_request_duration_seconds_count)';before=query(expr,fixed)
        procs['prometheus'].terminate();procs['prometheus'].wait(10);launch('prometheus',promargs)
        after=wait(lambda:query(expr,fixed),'Prometheus WAL restart');require(before['data']['result']==after['data']['result']and before['data']['result'],'retained series changed')
        result['retention']={'restart_fixed_timestamp_query_identical':True,'time':'1h','block_size':'64MB','disk_hard_cap_claim':False,'total_tsdb_bytes':sum(v.stat().st_size for v in(W/'tsdb').rglob('*')if v.is_file())}
        result['host_recovery']=dep.recover();refresh();replay=record(request());require(replay[0]==200 and json.loads(replay[1])==order,'host recovery replay')
        count=sql("select count(*) from sales.customer_order where tenant_id='b0000000-0000-4000-8000-000000000001'and order_id='"+order['id']+"'");require(count=='1','duplicate durable order')
        result['stop']=dep.close();dep=None;require(result['stop']['tree_empty']and result['stop']['exit_code']==0,'supervised stop failed')
        result.update(state='PASS',durable_order_count=1,pending_to_firing_seconds=round(fired-pending,3));log('recovered',load=result['load'],order_count=1)
    except BaseException as e:
        result.update(error_type=type(e).__name__,error=str(e));raise
    finally:
        if renamed:
            try:sql('alter table sales.customer_order_reference_fault rename to customer_order');result['fault_restored_after_failure']=True
            except BaseException as e:result['fault_restore_failure']=str(e)
        if dep:
            try:result['failure_stop']=dep.close()
            except BaseException as e:result['failure_stop_error']=str(e)
        for name in ['fixture','prometheus','postgres']:
            q=procs.get(name)
            if q is None:continue
            if q.poll()is None:
                try:
                    if name=='fixture':q.stdin.write(b'STOP\n');q.stdin.flush();q.wait(5)
                    elif name=='postgres':command('postgres-stop',[bins['pg_ctl'],'-D',data,'-m','fast','-w','stop']);q.wait(10)
                    else:q.terminate();q.wait(10)
                except BaseException:q.terminate();q.wait(5)
            if name=='fixture':q.stdin.close();q.stdout.close()
        for f in handles:f.close()
        result['children_exited']=all(q.poll()is not None for q in procs.values());result['seconds']=round(time.monotonic()-start,3)
        (W/'result.json').write_bytes(raw_json(result));print('OBSERVED_REFERENCE '+json.dumps(result),flush=True)
    return 0
if __name__=='__main__':raise SystemExit(main())
