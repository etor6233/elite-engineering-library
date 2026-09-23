"""AUTHORED consumer lock; generated after choosing a library commit, not inside it.

Every selected pack is read from the chosen Git object, not a mutable working
tree. The lock belongs to the consumer; no self-referential pack hash is needed.
"""
from pathlib import Path
import argparse
import hashlib
import json
import re
import subprocess
from cloud_control import contained,raw,require

BASE="markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md"
EXT="markdown_system/FRANCHISE_COMPLETE_EXTENSION_PACK_PLAN_V403.md"
EXTRAS={"BUSINESS-FUNCTION-OPERATING-V403","TS-DESIGN-SYSTEM-V403","TS-FRANCHISE-EXPERIENCE-V403","FRANCHISE-CLOUD-EXECUTION-V403"}
AUTHORITIES=["AGENT_SYSTEM_START.md","INSTALL_AGENT_BRIDGE.ps1","materialize_markdown_pack.ps1","implementation_packs/MARKDOWN_COMPOSITOR_CORE.md","implementation_packs/PROJECT_OPERATING_CONNECTION.md","markdown_system/FRANCHISE_PROJECT_OPERATING_PROTOCOL.md"]

def git_bytes(repo,commit,path):
    contained(repo,path)
    q=subprocess.run(["git","show",commit+":"+path],cwd=repo,stdin=subprocess.DEVNULL,capture_output=True,timeout=30,shell=False)
    require(q.returncode==0,"selected Git object lacks required file: "+path)
    return q.stdout

def parse_plan(data):
    match=re.search(rb"```json\s*\n(.*?)\n```",data,re.S)
    require(match is not None,"selected plan JSON absent")
    return json.loads(match.group(1))

def validate_selection(base,extended):
    require(len(base["packs"])==116,"baseline must retain all 116 selected packs")
    require(extended["packs"][:116]==base["packs"],"base selection changed or reordered")
    tail=extended["packs"][116:]
    require(len(tail)==4 and {p["packId"] for p in tail}==EXTRAS,"exact four V403 extensions required")
    require(len({p["packId"] for p in extended["packs"]})==120,"duplicate selected pack")
    return extended

def seal(repo,commit,output,repository="https://github.com/etor6233/elite-engineering-library.git",fixture=False):
    repo=Path(repo).resolve();output=Path(output).absolute()
    require(re.fullmatch(r"[0-9a-f]{40}",commit),"full exact commit required")
    require(not output.exists(),"consumer lock destination must be absent")
    for p in (output,*output.parents):
        require(not p.is_symlink() and not (hasattr(p,"is_junction") and p.is_junction()),"linked output rejected")
    require(output!=repo and repo not in output.parents,"consumer lock must be outside library")
    require(repository=="https://github.com/etor6233/elite-engineering-library.git" or (fixture and Path(repository).is_dir()),"approved read-only origin required")
    base_bytes=git_bytes(repo,commit,BASE);plan_bytes=git_bytes(repo,commit,EXT)
    plan=validate_selection(parse_plan(base_bytes),parse_plan(plan_bytes))
    files={};packs=[]
    for item in plan["packs"]:
        value=hashlib.sha256(git_bytes(repo,commit,item["path"])).hexdigest();files[item["path"]]=value
        packs.append({k:item[k] for k in ("packId","path","version")} | {"sha256":value})
    authorities={p:hashlib.sha256(git_bytes(repo,commit,p)).hexdigest() for p in AUTHORITIES}
    files.update(authorities);files[BASE]=hashlib.sha256(base_bytes).hexdigest();files[EXT]=hashlib.sha256(plan_bytes).hexdigest()
    profile={"schema":"elite-cloud-composition/v403.1","scope":"LIBRARY_INFRASTRUCTURE","production_authorized":False,"selection":{"path":EXT,"sha256":files[EXT],"pack_count":120},"base_selection":{"path":BASE,"sha256":files[BASE],"pack_count":116},"pack_locks":packs,"authority_locks":authorities,"extension_policy":"BASE116_PLUS_EXACT_FOUR_V403"}
    lock={"schema":"elite-library-checkout/v403.1","repository":repository,"commit":commit,"files":files,"publication_status":"NOT_CHECKED; exact commit must be available on chosen origin before remote bootstrap"}
    output.mkdir(parents=True);(output/"library-checkout.lock.json").write_bytes(raw(lock));(output/"composition.json").write_bytes(raw(profile))
    receipt={"schema":"elite-consumer-seal/v403.1","result":"PASS","claim":"120 selected Git-object file hashes","library_commit":commit,"production_authorized":False,"remote_execution":"NOT_RUN","files":{p.name:hashlib.sha256(p.read_bytes()).hexdigest() for p in output.iterdir()}}
    (output/"seal-receipt.json").write_bytes(raw(receipt));return receipt

if __name__=="__main__":
    p=argparse.ArgumentParser();p.add_argument("--library",required=True);p.add_argument("--commit",required=True);p.add_argument("--output",required=True);a=p.parse_args()
    print(raw(seal(a.library,a.commit,a.output)).decode(),end="")
