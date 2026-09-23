"""AUTHORED local activation journal around the existing Windows job launcher.

Single trusted local operator. No SCM registration, cloud rollout, hostile-user
isolation or secret storage. Failed candidate health restores the prior process
and pointer; no automatic database downgrade is ever issued.
"""
from __future__ import annotations
import ctypes,hashlib,json,msvcrt,os,re,subprocess,sys,threading,time,uuid
from pathlib import Path
from urllib.parse import urlsplit
from urllib.request import build_opener,ProxyHandler
from build_local_reference import plain,require,read_json,relative,sha,raw_json,tree

def loopback(value,scheme):
    u=urlsplit(value)
    require(u.scheme==scheme and u.hostname=='127.0.0.1'and u.port and 1024<=u.port<=65535 and not u.password and not u.fragment,'owned loopback endpoint required')
    return u
def runtime_config(value):
    require(set(value) in ({'database_url','issuer','api_port','web_port'},{'database_url','issuer','api_port','web_port','metrics_port'}),'runtime config fields')
    u=loopback(value['database_url'],'postgres')
    require(u.username=='postgres'and re.fullmatch(r'/elite_payment_connected_[0-9a-f]{32}',u.path)and u.query=='sslmode=disable','synthetic reference database required')
    u=loopback(value['issuer'],'http');require(not u.username and u.path in ('','/'),'fixture issuer required')
    ports=[value[x]for x in ['api_port','web_port','metrics_port']if x in value]
    require(all(type(x)is int and 1024<=x<=65535 for x in ports)and len(set(ports))==len(ports),'distinct loopback ports required')
    return value
def verify_release(root,inventory_path,digest):
    root=plain(root);files=read_json(inventory_path,digest)
    require(files and tree(root)==files,'artifact inventory mismatch')
    profile=json.loads((root/'LOCAL_REFERENCE_ONLY.json').read_text())
    require(profile['scope']=='LOCAL_FIXTURES'and profile['production_admitted']is False,'local fixture artifact required')
    require(sha(root/'node.exe')==profile['node_sha256'],'node binding changed')
    return {'root':str(root),'inventory':str(plain(inventory_path)),'inventory_sha256':digest,'source_sha256':profile['source_sha256']}
def atomic(path,value):
    temp=path.with_name(path.name+'.'+uuid.uuid4().hex+'.tmp')
    with temp.open('xb')as f:f.write(raw_json(value));f.flush();os.fsync(f.fileno())
    os.replace(temp,path)
def probe(config,seconds=60,alive=lambda:True):
    client=build_opener(ProxyHandler({}));end=time.monotonic()+seconds;checks={}
    while time.monotonic()<end:
        require(alive(),'candidate process exited before health')
        try:
            for name,url in [('api',f"http://127.0.0.1:{config['api_port']}/health/live"),('web',f"http://127.0.0.1:{config['web_port']}/")]:
                with client.open(url,timeout=2)as r:
                    body=r.read(1048576);require(r.status==200,'health status')
                    require(json.loads(body).get('status')=='ok'if name=='api'else b'<html'in body.lower(),'health body')
                    checks[name]={'status':r.status,'body_sha256':hashlib.sha256(body).hexdigest()}
            return checks
        except (OSError,ValueError):time.sleep(.1)
    raise ValueError('candidate health timeout')
class Instance:
    def __init__(self,release,config,root):
        self.cancel=threading.Event();self.result=None;self.config=config;self.root=root
        source=Path(__file__).resolve().parent.parent;sys.path.insert(0,str(source/'reference_telemetry'));import governed_launch
        script=source/'ci/serve_local_reference.py'
        job={'release':release,'runtime':config,'worker_sha256':sha(script)}
        config_path=root/'worker.json';atomic(config_path,job)
        exe=plain(sys.executable);runtime_env={k:str(Path(os.environ['SYSTEMROOT']))for k in ['SYSTEMROOT','WINDIR']}
        runtime_env.update(TEMP=str(root),TMP=str(root))
        profile={'schema':'elite-native-launch-qualification/v2','id':'local-franchise',
          'executable':{'path':str(exe),'sha256':sha(exe),'bytes':exe.stat().st_size},
          'arguments':['-X','utf8','-B',str(script),'--profile',str(config_path),'--sha256',sha(config_path)],
          'cwd':str(root),'environment':runtime_env,
          'budgets':{'timeout':1200 if 'metrics_port'in config else 600,'output_bytes':2097152,'processes':12,'commit_bytes':1073741824},
          'shutdown':{'protocol':'win32-inherited-event/v1','grace_seconds':25}}
        raw=raw_json(profile);atomic(root/'launch.json',profile)
        def run():self.result=governed_launch.launch(raw,hashlib.sha256(raw).hexdigest(),cancel=self.cancel)
        self.thread=threading.Thread(target=run);self.thread.start()
    def stop(self):
        self.cancel.set();self.thread.join(35);require(not self.thread.is_alive(),'launcher did not finish')
        r=self.result;require(r and r.status=='CAPTURED'and r.outcome and r.outcome.tree_empty,'process tree not empty')
        out=r.outcome
        (self.root/'stdout.log').write_bytes(out.stdout);(self.root/'stderr.log').write_bytes(out.stderr)
        receipt={'status':out.status,'exit_code':out.exit_code,'tree_empty':out.tree_empty,
          'shutdown_state':out.shutdown_state,'observed_peak_job_bytes':out.observed_peak_job_bytes,
          'stdout_sha256':sha(self.root/'stdout.log'),'stderr_sha256':sha(self.root/'stderr.log')}
        atomic(self.root/'stop.json',receipt);return receipt
class Deployment:
    def __init__(self,root,config):
        self.root=plain(root);self.root.mkdir(exist_ok=True);self.config=runtime_config(config);self.instance=None;self.serial=0
        self.lock=(self.root/'operator.lock').open('a+b');self.lock.seek(0)
        try:msvcrt.locking(self.lock.fileno(),msvcrt.LK_NBLCK,1)
        except OSError:self.lock.close();raise ValueError('another local operator owns deployment')
        if os.fstat(self.lock.fileno()).st_size==0:self.lock.write(b'0');self.lock.flush()
    def _start(self,release):
        verify_release(release['root'],release['inventory'],release['inventory_sha256'])
        self.serial+=1;work=self.root/('run-'+uuid.uuid4().hex);work.mkdir()
        self.instance=Instance(release,self.config,work)
        return probe(self.config,alive=self.instance.thread.is_alive)
    def _stop(self):
        if self.instance:
            instance=self.instance;self.instance=None;return instance.stop()
    def current(self):
        p=self.root/'current.json';return json.loads(p.read_text())if p.exists()else None
    def activate(self,artifact,inventory_path,digest):
        release=verify_release(artifact,inventory_path,digest);prior=self.current()
        if prior:verify_release(prior['root'],prior['inventory'],prior['inventory_sha256'])
        receipt={'schema':'elite-local-activation/v1','state':'STARTING','previous':prior,'candidate':release,'migration_action':'NONE'}
        journal=self.root/('activation-'+uuid.uuid4().hex+'.json');atomic(journal,receipt)
        try:
            receipt['stopped_previous']=self._stop();receipt['health']=self._start(release)
            atomic(self.root/'current.json',release);receipt['state']='ACTIVE'
        except BaseException:
            receipt['stopped_candidate']=self._stop()
            if prior:
                receipt['rollback_health']=self._start(prior);receipt['state']='ROLLED_BACK'
                require(self.current()==prior,'stable pointer changed on failed candidate')
            else:receipt['state']='FAILED_NO_PREVIOUS'
            raise
        finally:atomic(journal,receipt)
        return receipt
    def recover(self):
        prior=self.current();require(prior is not None,'no committed release to recover')
        self._stop();return self._start(prior)
    def close(self):
        try:return self._stop()
        finally:
            self.lock.seek(0);msvcrt.locking(self.lock.fileno(),msvcrt.LK_UNLCK,1);self.lock.close()
