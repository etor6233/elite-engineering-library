"""AUTHORED child of the existing finite Windows Job Object launcher."""
from pathlib import Path
import argparse,ctypes,os,shutil,subprocess,sys,time
from local_release import read_json,runtime_config,verify_release,require,sha
from build_local_reference import environment,copy_tree

def main():
    p=argparse.ArgumentParser();p.add_argument('--profile',required=True);p.add_argument('--sha256',required=True);a=p.parse_args()
    job=read_json(a.profile,a.sha256);require(set(job)=={'release','runtime','worker_sha256'},'worker profile fields')
    require(sha(Path(__file__))==job['worker_sha256'],'worker source changed')
    r=job['release'];verify_release(r['root'],r['inventory'],r['inventory_sha256']);root=Path(r['root']);cfg=runtime_config(job['runtime'])
    raw=os.environ.get('ELITE_STOP_EVENT_HANDLE','');require(raw.isdecimal()and int(raw)>0,'trusted inherited stop event required');handle=int(raw)
    wait=ctypes.WinDLL('kernel32',use_last_error=True).WaitForSingleObject
    wait.argtypes=[ctypes.c_void_p,ctypes.c_uint32];wait.restype=ctypes.c_uint32
    env=environment(root/'node.exe',Path.cwd());env.update(DATABASE_URL=cfg['database_url'],OIDC_ISSUER=cfg['issuer'],OIDC_AUDIENCE='elite-local-reference',
       HTTP_ADDRESS=f"127.0.0.1:{cfg['api_port']}",HOSTNAME='127.0.0.1',PORT=str(cfg['web_port']),NODE_ENV='production',
       APP_BASE_URL=f"http://127.0.0.1:{cfg['web_port']}",ENTERPRISE_API_BASE_URL=f"http://127.0.0.1:{cfg['api_port']}")
    require(wait(handle,0)==258,'stop requested before launch')
    if 'metrics_port'in cfg:
        env.update(HTTP_METRICS_ENABLED='true',HTTP_METRICS_ADDRESS=f"127.0.0.1:{cfg['metrics_port']}")
    procs=[]
    try:
        # Restrict extra inheritance to the selected stop handle for the Go host.
        os.set_handle_inheritable(handle,True);startup=subprocess.STARTUPINFO();startup.lpAttributeList={'handle_list':[handle]}
        api_env={**env,'ELITE_STOP_EVENT_HANDLE':raw}
        api=subprocess.Popen([root/'api/electromobility-api.exe'],cwd=root,env=api_env,stdin=subprocess.DEVNULL,
            stdout=sys.stdout,stderr=sys.stderr,startupinfo=startup,creationflags=subprocess.CREATE_NO_WINDOW)
        procs.append(api);os.set_handle_inheritable(handle,False)
        webroot=Path.cwd()/'web';copy_tree(root/'web',webroot)
        script=Path(__file__).with_name('web_local_stop.cjs')
        web=subprocess.Popen([root/'node.exe',script,webroot/'server.js'],cwd=webroot,env=env,stdin=subprocess.PIPE,
            stdout=sys.stdout,stderr=sys.stderr,creationflags=subprocess.CREATE_NO_WINDOW)
        procs.append(web)
        while wait(handle,50)==258:
            require(all(x.poll()is None for x in procs),'application exited unexpectedly')
        if web.poll()is None:web.stdin.write(b'STOP\n');web.stdin.flush();web.stdin.close()
        # Go uses its inherited event; Next handles its own SIGTERM listeners.
        deadline=time.monotonic()+20
        for proc,expected in [(api,0),(web,143)]:
            proc.wait(max(.1,deadline-time.monotonic()));require(proc.returncode==expected,'application shutdown failed')
        print('LOCAL_REFERENCE_STOPPED',flush=True);return 0
    finally:
        for proc in procs:
            if proc.poll()is None:proc.terminate()
        for proc in procs:proc.wait(5)
if __name__=='__main__':raise SystemExit(main())
