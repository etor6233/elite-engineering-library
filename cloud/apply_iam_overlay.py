"""AUTHORED IAM overlay adapter; reuses the V403 UI transactional overlay owner."""
from pathlib import Path
import argparse
import importlib.util
from cloud_control import file_hash,raw,read,require

def apply(target,verify=False):
    target=Path(target).resolve();bundle=target/"cloud_overlay"
    manifest=read(bundle/"overlay-manifest.json")
    engine=target/"experience_overlay/apply_overlay.py"
    require(file_hash(engine)==manifest["overlay_engine_sha256"],"existing shared overlay engine changed")
    spec=importlib.util.spec_from_file_location("v403_shared_overlay",engine);module=importlib.util.module_from_spec(spec);spec.loader.exec_module(module)
    result=module.apply(bundle,target,verify);result["claim"]="HASH_LOCKED_CLOUD_IAM_OVERLAY"
    return result

if __name__=="__main__":
    p=argparse.ArgumentParser();p.add_argument("--target",required=True);p.add_argument("--report",required=True);p.add_argument("--verify",action="store_true");a=p.parse_args()
    report=Path(a.report);require(not report.exists(),"fresh IAM receipt required");result=apply(a.target,a.verify)
    report.parent.mkdir(parents=True,exist_ok=True);report.write_bytes(raw(result));print(raw(result).decode(),end="")
