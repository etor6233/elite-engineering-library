"""AUTHORED whole-stack adapter over the same bounded runner/durable fence core.

One approved DAG step per invocation. Execution receipts never prove resource
health. Every dependency needs a separate hash-bound semantic reconciliation.
"""
from pathlib import Path
import argparse
from cloud_control import raw,read,digest,file_hash,contained,require,timestamp,authorize_record,release,bind_receipts
from execute_step import _run_fenced
from stack_manifest import build_stack
from semantic_contract import expected

def envelope(config,images,source_root,composition_path,artifact_root,fence,artifact_receipts=None):
    source=Path(source_root).resolve();cp=Path(composition_path).resolve();root=Path(artifact_root).resolve()
    plan=build_stack(config,images,source_root=source,composition=read(cp))
    e={"schema":"elite-cloud-stack-execution/v403.1","scope":"LIBRARY_INFRASTRUCTURE","environment":config["environment"],"production_authorized":False,"execution":"NOT_RUN","config":config,"images":images,"source_root":str(source),"composition_path":str(cp),"composition_sha256":file_hash(cp),"artifact_root":str(root),"execution_fence":fence,"stack":plan,"artifact_receipts":artifact_receipts or {}}
    validate_envelope(e);return e

def validate_envelope(e):
    require(set(e)=={"schema","scope","environment","production_authorized","execution","config","images","source_root","composition_path","composition_sha256","artifact_root","execution_fence","stack","artifact_receipts"},"exact stack execution envelope required")
    require(e["schema"]=="elite-cloud-stack-execution/v403.1" and e["scope"]=="LIBRARY_INFRASTRUCTURE" and e["environment"] in ("development","staging") and e["production_authorized"] is False and e["execution"]=="NOT_RUN","non-production stack envelope required")
    require(e["environment"]==e["config"]["environment"],"stack environment mismatch")
    for label in ("source_root","composition_path","artifact_root"):
        p=Path(e[label]);require(p.is_absolute(),"absolute pinned worker path required")
        for item in (p,*p.parents):require(not item.is_symlink() and not (hasattr(item,"is_junction") and item.is_junction()),"linked stack execution input rejected")
    require(file_hash(e["composition_path"])==e["composition_sha256"],"consumer composition bytes changed")
    regenerated=build_stack(e["config"],e["images"],source_root=e["source_root"],composition=read(e["composition_path"]))
    require(e["stack"]==regenerated,"stack commands/effects/dependencies/config differ from exact generator")
    for name,body in regenerated["manifests"].items():require(contained(Path(e["artifact_root"]),name).read_bytes()==raw(body),"runtime manifest bytes changed")
    if e["images"]["method"]=="EXECUTED_TARGET":
        require(set(e["artifact_receipts"])==set(e["images"]["artifacts"]),"all runtime families require exact artifact admission receipts")
        for family,ref in e["artifact_receipts"].items():
            p=contained(Path(e["artifact_root"]),ref["path"]);require(file_hash(p)==ref["sha256"],"artifact release receipt changed")
            proof=release(read(p),e["config"]["target"]);artifact=e["images"]["artifacts"][family]
            require(proof["image"]==artifact["image"] and proof["source_commit"]==artifact["source_commit"] and proof["composition_sha256"]==artifact["composition_sha256"],"artifact receipt source or image mismatch")
            bind_receipts(proof,Path(e["artifact_root"]))
    else:require(e["artifact_receipts"]=={},"fixtures must not fabricate real artifact admissions")
    fence=e["execution_fence"];require(fence["backend"] in ("LOCAL_FIXTURE","GCS_GENERATION_0"),"durable fence required")
    if fence["backend"]=="LOCAL_FIXTURE":
        require(set(fence)=={"backend","directory"} and Path(fence["directory"]).is_dir(),"persistent fixture fence directory required")
    else:
        import re
        require(set(fence)=={"backend","uri"} and re.fullmatch(r"gs://[a-z0-9][a-z0-9._-]+/[A-Za-z0-9/_-]+",fence["uri"]),"approved existing control-bucket namespace required")
    return regenerated

def observation(e,ref,now,method):
    require(isinstance(ref,dict) and set(ref)=={"path","sha256"},"hash-bound observation reference required")
    p=contained(Path(e["artifact_root"]),ref["path"]);require(file_hash(p)==ref["sha256"],"stack observation bytes changed")
    ob=read(p)
    require(ob["schema"]=="elite-cloud-stack-observation/v403.1" and ob["plan_sha256"]==digest(e) and ob["target_sha256"]==e["stack"]["target_sha256"] and ob["images_sha256"]==e["stack"]["images_sha256"] and ob["method"]==method and ob["result"]=="PASS","stack observation identity mismatch")
    require(timestamp(ob["observed_at"])<=timestamp(now)<=timestamp(ob["expires_at"]) and (timestamp(now)-timestamp(ob["observed_at"])).total_seconds()<=900,"stack observation stale")
    proof=ob["source_evidence"];proofpath=contained(Path(e["artifact_root"]),proof["path"]);require(file_hash(proofpath)==proof["sha256"],"source semantic evidence changed")
    evidence=read(proofpath)
    require(evidence["method"]==method and evidence["target_sha256"]==ob["target_sha256"] and evidence["images_sha256"]==ob["images_sha256"] and evidence["result"]=="PASS","semantic evidence target mismatch")
    require(timestamp(evidence["observed_at"])<=timestamp(now)<=timestamp(evidence["expires_at"]) and (timestamp(now)-timestamp(evidence["observed_at"])).total_seconds()<=900,"underlying semantic evidence stale")
    step,claim,name,checks=expected(e,ob)
    require(evidence["schema"]=="elite-cloud-target-semantic-report/v403.1" and evidence["claim"]==claim and evidence["gate_or_resource"]==name and evidence["for_step"]==step["id"] and evidence["desired_step_sha256"]==digest(step),"semantic owner report belongs to a different claim/resource")
    require(set(evidence["checks"])==checks and all(v is True for v in evidence["checks"].values()),"exact semantic checks incomplete")
    return ob

def validate_prerequisites(e,step,refs,now,method):
    require(set(refs)=={"gates","dependencies"},"exact prerequisite maps required")
    require(set(refs["gates"])==set(step["required_target_receipts"]),"required target gate omitted")
    require(set(refs["dependencies"])==set(step["depends_on"]),"DAG dependency reconciliation omitted")
    for name,ref in refs["gates"].items():
        ob=observation(e,ref,now,method);require(ob["kind"]=="TARGET_GATE_EVIDENCE" and ob["for_step"]==step["id"] and ob["name"]==name,"target gate binding mismatch")
    for name,ref in refs["dependencies"].items():
        ob=observation(e,ref,now,method);require(ob["kind"]=="STACK_RESOURCE_RECONCILIATION" and ob["for_step"]==name,"dependency observation step mismatch")
        receipt=ob["execution_receipt"];rp=contained(Path(e["artifact_root"]),receipt["path"]);require(file_hash(rp)==receipt["sha256"],"dependency execution receipt changed")
        dependency=next(s for s in e["stack"]["commands"] if s["id"]==name)
        run=read(rp)
        success=(run.get("result")=="BARRIER_EVIDENCE_CHECKED" if dependency["effect"]=="BARRIER" else type(run.get("exit_code")) is int and run["exit_code"]==0 and run.get("result")!="BARRIER_EVIDENCE_CHECKED")
        require(run.get("schema")=="elite-cloud-step-execution/v403.1" and run["plan_sha256"]==digest(e) and run["step"]==name and run["method"]==method and success,"dependency execution missing or failed")
        require(ob["desired_step_sha256"]==digest(dependency),"reconciliation desired resource changed")

def execute(e,step_id,approval,now,tool,runner,evidence,prerequisites,*,fixture=False):
    validate_envelope(e);authorize_record(e,approval,now)
    method="SIMULATED_PROVIDER" if fixture else "EXECUTED_TARGET"
    require(e["images"]["method"]==("SIMULATED_ARTIFACTS" if fixture else "EXECUTED_TARGET"),"fixture images cannot enter real stack execution")
    require(e["execution_fence"]["backend"]==("LOCAL_FIXTURE" if fixture else "GCS_GENERATION_0"),"fixture/real fence mismatch")
    matches=[s for s in e["stack"]["commands"] if s["id"]==step_id];require(len(matches)==1,"unknown stack step");step=matches[0]
    validate_prerequisites(e,step,prerequisites,now,method)
    out=Path(evidence).absolute();require(not out.exists(),"fresh stack receipt required")
    if step["effect"]=="BARRIER":
        require(not step["argv"],"barrier cannot execute a command")
        result={"schema":"elite-cloud-step-execution/v403.1","plan_sha256":digest(e),"step":step_id,"method":method,"result":"BARRIER_EVIDENCE_CHECKED","cloud_demonstrated":False,"production_authorized":False}
        out.parent.mkdir(parents=True,exist_ok=True)
        with out.open("xb") as f:f.write(raw(result))
        return result
    return _run_fenced(e,step,tool,runner,evidence,fixture=fixture,working_directory=e["artifact_root"])

if __name__=="__main__":
    p=argparse.ArgumentParser();sub=p.add_subparsers(dest="command",required=True)
    q=sub.add_parser("prepare");q.add_argument("--config",required=True);q.add_argument("--source",required=True);q.add_argument("--composition",required=True);q.add_argument("--artifact-root",required=True);q.add_argument("--fence",required=True);q.add_argument("--output",required=True);q.add_argument("--artifact-receipts")
    q=sub.add_parser("execute");q.add_argument("--plan",required=True);q.add_argument("--step",required=True);q.add_argument("--approval",required=True);q.add_argument("--tool",required=True);q.add_argument("--runner",required=True);q.add_argument("--prerequisites",required=True);q.add_argument("--receipt",required=True);q.add_argument("--authorize-execution",action="store_true")
    a=p.parse_args()
    if a.command=="prepare":
        b=read(a.config);result=envelope(b["config"],b["images"],a.source,a.composition,a.artifact_root,read(a.fence),read(a.artifact_receipts) if a.artifact_receipts else None);out=Path(a.output);require(not out.exists(),"fresh execution envelope required");out.write_bytes(raw(result))
    else:
        from datetime import datetime,timezone
        require(a.authorize_execution,"explicit cloud stack execution authorization required")
        result=execute(read(a.plan),a.step,read(a.approval),datetime.now(timezone.utc).isoformat().replace("+00:00","Z"),read(a.tool),read(a.runner),a.receipt,read(a.prerequisites))
    print(raw(result).decode(),end="")
