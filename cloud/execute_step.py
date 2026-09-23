"""AUTHORED narrow non-production execution adapter, one approved step at a time.

Reuses bounded_process from the immutable SECURE-OPS-DELIVERY-CORE owner.
Human approval, exact tool/runner identity and fresh step evidence are required.
No automatic retries after an ambiguous provider result. Fixtures use a separate
method and can never create an EXECUTED_TARGET receipt.
"""
from __future__ import annotations
import importlib.util
import json
import os
from pathlib import Path
import sys
import re
from cloud_control import authorize, digest, file_hash, raw, read, require, contained, timestamp

PREREQUISITES={"deploy-no-traffic":{"target_access","budget","migrations","before_state"},"promote-after-health":{"candidate_health","identity","migrations","termination_recovery"},"rollback":{"rollback_approval","prior_health","migrations","pending_work_reconciliation"}}

def validate_prerequisite(plan, step, reference, method, now):
    require(isinstance(reference,dict) and set(reference)=={"path","sha256"},"hash-bound prerequisite reference required")
    path=contained(Path(plan["evidence_root"]),reference["path"])
    require(file_hash(path)==reference["sha256"],"prerequisite bytes changed")
    p=read(path)
    require(p["schema"]=="elite-cloud-step-prerequisite/v403.1" and p["plan_sha256"]==digest(plan) and p["result"]=="PASS" and p["for_step"]==step and p["method"]==method,"prerequisite identity or method mismatch")
    require(p["target_sha256"]==plan["target_sha256"] and p["image"]==(plan["prior_image"] if step=="rollback" else plan["image"]),"prerequisite target/image mismatch")
    require(timestamp(p["observed_at"])<=timestamp(now)<=timestamp(p["expires_at"]) and (timestamp(now)-timestamp(p["observed_at"])).total_seconds()<=900,"prerequisite stale or expired")
    require(set(p["checks"])==PREREQUISITES[step] and all(v is True for v in p["checks"].values()),"step-specific prerequisites incomplete")
    require(p["evidence"] and {e["check"] for e in p["evidence"]}==PREREQUISITES[step],"step evidence coverage incomplete")
    for e in p["evidence"]:
        ep=contained(Path(plan["evidence_root"]),e["path"])
        require(file_hash(ep)==e["sha256"],"prerequisite observation tampered")
        ob=read(ep)
        require(ob["check"]==e["check"] and ob["result"]=="PASS" and ob["method"]==method and ob["target_sha256"]==p["target_sha256"] and ob["image"]==p["image"],"prerequisite semantic observation mismatch")
        require(timestamp(ob["observed_at"])<=timestamp(now)<=timestamp(ob["expires_at"]) and (timestamp(now)-timestamp(ob["observed_at"])).total_seconds()<=900,"individual prerequisite observation stale or expired")

def load_runner(path, expected):
    path=Path(path).resolve();require(file_hash(path)==expected,"bounded runner hash mismatch")
    sys.path.insert(0,str(path.parent))
    spec=importlib.util.spec_from_file_location("elite_existing_official_runner",path)
    module=importlib.util.module_from_spec(spec);spec.loader.exec_module(module)
    return module.bounded_process

def execute(plan, step_id, approval, now, tool, runner, evidence, *, fixture=False, prerequisite=None):
    authorize(plan,approval,now)
    require(tool["method"] == ("SIMULATED_PROVIDER" if fixture else "EXECUTED_TARGET"),"explicit execution method required")
    binary=Path(tool["binary"]).resolve();require(file_hash(binary)==tool["binary_sha256"],"tool binary identity mismatch")
    require(tool["admission"]=="FIXTURE_ONLY" if fixture else tool["admission"]=="ADMITTED_EXACT_ARTIFACT", "tool admission incomplete")
    matches=[x for x in plan["commands"] if x["id"]==step_id];require(len(matches)==1,"unknown or duplicate step")
    step=matches[0]
    if step["effect"]=="WRITE":
        validate_prerequisite(plan,step_id,prerequisite,tool["method"],now)
    return _run_fenced(plan,step,tool,runner,evidence,fixture=fixture)

def _run_fenced(plan,step,tool,runner,evidence,*,fixture=False,working_directory=None):
    """Private shared effect core. Public adapters validate schema/approval/gates first."""
    step_id=step["id"];binary=Path(tool["binary"]).resolve()
    require(tool["method"] == ("SIMULATED_PROVIDER" if fixture else "EXECUTED_TARGET"),"explicit execution method required")
    require(file_hash(binary)==tool["binary_sha256"],"tool binary identity mismatch")
    require(tool["admission"]==("FIXTURE_ONLY" if fixture else "ADMITTED_EXACT_ARTIFACT"),"tool admission incomplete")
    if not fixture:
        from tool_admission import validate
        validate(tool)
    out=Path(evidence).absolute();require(not out.exists(),"step receipt already exists; reconcile instead of replay")
    out.parent.mkdir(parents=True,exist_ok=True)
    bounded=load_runner(runner["path"],runner["sha256"])
    args=[str(binary),*tool.get("prefix_arguments",[]),*step["argv"]]
    if not fixture:
        require(step["argv"][0]=="gcloud" and not tool.get("prefix_arguments"),"exact gcloud binary only")
        args=[str(binary),*step["argv"][1:]]
    # Only explicit externally supplied environment names, never inline values.
    allowed={"PATH","SYSTEMROOT","WINDIR","TEMP","TMP","LANG","LC_ALL","SSL_CERT_FILE","SSL_CERT_DIR"}
    allowed.update(tool.get("environment_variable_names",[]))
    require(not any(x in allowed for x in ("GITHUB_TOKEN","GH_TOKEN")),"repository token must not enter deployer")
    env={k:v for k,v in os.environ.items() if k in allowed}
    intent={"schema":"elite-cloud-intent/v403.1","plan_sha256":digest(plan),"step":step_id,"arguments_sha256":digest(args),"tool_binary_sha256":tool["binary_sha256"],"state":"CLAIMED_RECONCILE_BEFORE_RETRY"}
    intent_name=digest(plan)+"-"+step_id+".json"
    if fixture:
        # A stable controller journal independent of the chosen receipt filename.
        require(plan["execution_fence"]["backend"]=="LOCAL_FIXTURE","fixture fence must be in approved plan")
        journal=Path(plan["execution_fence"]["directory"]).resolve();require(journal.is_dir(),"persistent fixture journal required")
        require("journal_directory" not in tool or Path(tool["journal_directory"]).resolve()==journal,"unapproved fence namespace override")
        intent_path=journal/intent_name
        with intent_path.open("xb") as f:
            f.write(raw(intent));f.flush();os.fsync(f.fileno())
    else:
        # An ephemeral worker file is NOT a cloud fence. Create a durable object
        # before any provider effect, with generation-match=0 and no delete path.
        require(plan["execution_fence"]["backend"]=="GCS_GENERATION_0","durable GCS fence must be in approved plan")
        uri=plan["execution_fence"]["uri"]
        require("fence_uri" not in tool or tool["fence_uri"]==uri,"unapproved fence namespace override")
        require(re.fullmatch(r"gs://[a-z0-9][a-z0-9._-]+/[a-zA-Z0-9/_-]+",uri),"durable GCS intent namespace required")
        intent_path=out.parent/(intent_name+".request")
        with intent_path.open("xb") as f:
            f.write(raw(intent));f.flush();os.fsync(f.fileno())
        claim=bounded([str(binary),"storage","cp",str(intent_path),uri+"/"+intent_name,"--if-generation-match=0","--quiet"],out.parent,env,60)
        require(claim.returncode==0,"durable intent conflict or uncertain claim; reconcile before any effect")
    try:
        q=bounded(args,Path(working_directory) if working_directory else out.parent,env,60)
    except BaseException:
        uncertain={"schema":"elite-cloud-step-execution/v403.1","plan_sha256":digest(plan),"step":step_id,"method":tool["method"],"cloud_demonstrated":False,"result":"UNKNOWN_RECONCILE_BEFORE_RETRY","intent":intent_name}
        with out.open("xb") as f:f.write(raw(uncertain));f.flush();os.fsync(f.fileno())
        raise
    # Persist only exit code and output digests. Raw CLI output may contain private
    # configuration. A separate semantic probe provides the minimal observation.
    import hashlib
    result={"schema":"elite-cloud-step-execution/v403.1","plan_sha256":digest(plan),"step":step_id,"method":tool["method"],"cloud_demonstrated":False,"tool_binary_sha256":tool["binary_sha256"],"runner_sha256":runner["sha256"],"arguments_sha256":digest(args),"exit_code":q.returncode,"result":"COMMAND_SUCCEEDED_RECONCILIATION_REQUIRED" if q.returncode==0 else "UNKNOWN_RECONCILE_BEFORE_RETRY","stdout_sha256":hashlib.sha256(q.stdout).hexdigest(),"stderr_sha256":hashlib.sha256(q.stderr).hexdigest()}
    with out.open("xb") as f:f.write(raw(result));f.flush();os.fsync(f.fileno())
    return result

def main():
    import argparse
    p=argparse.ArgumentParser();p.add_argument("--plan",required=True);p.add_argument("--step",required=True);p.add_argument("--approval",required=True);p.add_argument("--tool",required=True);p.add_argument("--runner",required=True);p.add_argument("--receipt",required=True);p.add_argument("--prerequisite");p.add_argument("--authorize-execution",action="store_true");a=p.parse_args()
    require(a.authorize_execution,"explicit execution authorization required")
    from datetime import datetime,timezone
    result=execute(read(a.plan),a.step,read(a.approval),datetime.now(timezone.utc).isoformat().replace("+00:00","Z"),read(a.tool),read(a.runner),a.receipt,prerequisite=read(a.prerequisite) if a.prerequisite else None)
    print(raw(result).decode(),end="");return 0 if result["exit_code"]==0 else 2

if __name__=="__main__":sys.exit(main())
