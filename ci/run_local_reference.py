"""Finite local operator command. No provider activation or production effects."""
from pathlib import Path
import argparse,json,time
from local_release import Deployment,read_json,require,sha,verify_release
def main():
    p=argparse.ArgumentParser(description=__doc__)
    for name in ['artifact','inventory','inventory-sha256','runtime','runtime-sha256','migration-receipt','migration-receipt-sha256','state-directory']:p.add_argument('--'+name,required=True)
    p.add_argument('--seconds',type=int,default=300);a=p.parse_args();require(1<=a.seconds<=480,'finite duration 1..480 required')
    cfg=read_json(a.runtime,a.runtime_sha256);migration=read_json(a.migration_receipt,a.migration_receipt_sha256)
    require(migration['state']=='PASS'and migration['scope']=='LOCAL_FIXTURES'and migration['database']==cfg['database_url'],'migration receipt mismatch')
    files={p.name:sha(p)for p in(Path(a.artifact)/'migrations').glob('*.up.sql')};require(files==migration['files'],'migration receipt does not cover release')
    dep=Deployment(a.state_directory,cfg)
    try:
        dep.activate(a.artifact,a.inventory,a.inventory_sha256)
        print(f"LOCAL_REFERENCE_ACTIVE http://127.0.0.1:{cfg['web_port']} duration={a.seconds}",flush=True)
        deadline=time.monotonic()+a.seconds
        while time.monotonic()<deadline:
            require(dep.instance.thread.is_alive(),'application exited');time.sleep(.2)
    finally:
        receipt=dep.close();require(receipt and receipt['exit_code']==0,'local shutdown failed')
if __name__=='__main__':main()
