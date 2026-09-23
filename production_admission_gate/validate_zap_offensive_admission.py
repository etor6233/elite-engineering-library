from __future__ import annotations
import argparse, json, re, sys
from pathlib import Path
from urllib.parse import urlparse
from validate_k6_load_admission import atomic, canonical_sha, proof, read_json, safe, utc, verify_proof

KINDS=("WEB_AUTHENTICATED","API_AUTHENTICATED")
FIELDS={"schema","project_id","environment","release_digest","evaluated_at","expires_at","executor_ref","target","policy","manual_review","runs","output"}
RUN_FIELDS={"kind","execution_receipt","plan","report","expected_arguments"}
SOURCE="https://github.com/zaproxy/zaproxy/releases/tag/v2.17.0"

def https_url(value: object, label: str) -> None:
    if not isinstance(value,str): raise ValueError(f"{label} must be HTTPS")
    parsed=urlparse(value)
    if parsed.scheme!="https" or not parsed.hostname or parsed.username or parsed.password or parsed.fragment: raise ValueError(f"{label} must be credential-free HTTPS")

def alerts(report: dict[str,object]) -> list[dict[str,object]]:
    if not isinstance(report.get("site"),list) or not report["site"]: raise ValueError("ZAP report requires non-empty site list")
    found=[]
    for site in report["site"]:
        if not isinstance(site,dict) or not isinstance(site.get("@name"),str) or not isinstance(site.get("alerts"),list): raise ValueError("ZAP site/alerts shape is invalid")
        for alert in site["alerts"]:
            if not isinstance(alert,dict) or not {"pluginid","riskcode","name","instances"}.issubset(alert) or not isinstance(alert["instances"],list): raise ValueError("ZAP alert shape is invalid")
            try: risk=int(alert["riskcode"])
            except (TypeError,ValueError) as error: raise ValueError("ZAP riskcode is invalid") from error
            if risk not in (0,1,2,3,4): raise ValueError("ZAP riskcode is outside admitted range")
            found.append({"pluginid":str(alert["pluginid"]),"riskcode":risk,"name":str(alert["name"])})
    return found

def validate(profile: object, root: Path) -> dict[str,object]:
    if not isinstance(profile,dict) or set(profile)!=FIELDS: raise ValueError("profile fields must be exact")
    if profile["schema"]!="elite-zap-offensive-admission/v1" or profile["environment"]!="production" or not isinstance(profile["project_id"],str) or not profile["project_id"].strip() or not re.fullmatch(r"sha256:[0-9a-f]{64}",str(profile["release_digest"])): raise ValueError("profile identity is invalid")
    start=utc(profile["evaluated_at"],"evaluated_at"); end=utc(profile["expires_at"],"expires_at")
    if not start<end or (end-start).total_seconds()>30*86400: raise ValueError("admission validity must be positive and at most 30 days")
    if not isinstance(profile["executor_ref"],str) or not profile["executor_ref"].strip() or not isinstance(profile["target"],dict) or set(profile["target"])!={"base_url","api_base_url","authorization_ticket"}: raise ValueError("exact authorized target is required")
    target=profile["target"]; https_url(target["base_url"],"base_url"); https_url(target["api_base_url"],"api_base_url")
    if not isinstance(target["authorization_ticket"],str) or not target["authorization_ticket"].strip(): raise ValueError("authorization ticket is required")

    policy_path=safe(root,profile["policy"],"policy"); policy=read_json(policy_path,"policy")
    if set(policy)!={"schema","project_id","environment","release_digest","approved_at","approvers","accepted_alerts"} or policy.get("schema")!="elite-offensive-security-policy/v1": raise ValueError("policy fields/schema must be exact")
    for key in ("project_id","environment","release_digest"):
        if policy.get(key)!=profile[key]: raise ValueError(f"policy {key} mismatch")
    utc(policy["approved_at"],"policy.approved_at"); approvers=policy["approvers"]
    if not isinstance(approvers,list) or len(approvers)<2 or len(set(approvers))!=len(approvers) or any(not isinstance(x,str) or not x.strip() for x in approvers): raise ValueError("policy requires two distinct approvers")
    accepted=policy["accepted_alerts"]
    if not isinstance(accepted,list): raise ValueError("accepted_alerts must be a list")
    accepted_ids=set()
    for item in accepted:
        if not isinstance(item,dict) or set(item)!={"pluginid","riskcode","reason","owner","expires_at"} or int(item.get("riskcode",-1))<3 or not all(isinstance(item.get(k),str) and item[k].strip() for k in ("pluginid","reason","owner")): raise ValueError("accepted alert record is invalid")
        if utc(item["expires_at"],"accepted_alert.expires_at")<end: raise ValueError("accepted alert expires before admission")
        key=(item["pluginid"],int(item["riskcode"]));
        if key in accepted_ids: raise ValueError("duplicate accepted alert")
        accepted_ids.add(key)

    review_path=safe(root,profile["manual_review"],"manual_review"); review=read_json(review_path,"manual_review")
    if set(review)!={"schema","project_id","environment","release_digest","reviewed_at","reviewers","authorization_confirmed","authenticated_web_reviewed","authenticated_api_reviewed","business_logic_reviewed","notes"} or review.get("schema")!="elite-offensive-manual-review/v1": raise ValueError("manual review fields/schema must be exact")
    for key in ("project_id","environment","release_digest"):
        if review.get(key)!=profile[key]: raise ValueError(f"manual review {key} mismatch")
    utc(review["reviewed_at"],"manual_review.reviewed_at"); reviewers=review["reviewers"]
    if not isinstance(reviewers,list) or len(reviewers)<2 or len(set(reviewers))!=len(reviewers) or any(not isinstance(x,str) or not x.strip() for x in reviewers): raise ValueError("manual review requires two distinct reviewers")
    if any(review.get(k) is not True for k in ("authorization_confirmed","authenticated_web_reviewed","authenticated_api_reviewed","business_logic_reviewed")): raise ValueError("manual review assertions must all pass")
    if not isinstance(review["notes"],str) or not review["notes"].strip(): raise ValueError("manual review notes are required")

    runs=profile["runs"]
    if not isinstance(runs,list) or len(runs)!=2 or any(not isinstance(x,dict) or set(x)!=RUN_FIELDS for x in runs) or [x["kind"] for x in runs]!=list(KINDS): raise ValueError("authenticated web and API runs are required")
    evidence=[proof(root,policy_path),proof(root,review_path)]; unaccepted=[]
    for run in runs:
        kind=run["kind"]; plan_path=safe(root,run["plan"],f"{kind}.plan"); plan_proof=proof(root,plan_path)
        receipt_path=safe(root,run["execution_receipt"],f"{kind}.receipt"); receipt=read_json(receipt_path,f"{kind}.receipt")
        if receipt.get("schema")!="elite-official-tool-execution/v1" or receipt.get("project_id")!=profile["project_id"] or receipt.get("environment")!="production" or receipt.get("release_digest")!=profile["release_digest"] or receipt.get("control_id")!="OFFENSIVE_SECURITY" or receipt.get("exit_code")!=0: raise ValueError(f"{kind} execution identity/exit mismatch")
        tool=receipt.get("tool")
        if not isinstance(tool,dict) or set(tool)!={"name","version","source","sha256"} or tool.get("name")!="OWASP ZAP" or tool.get("version")!="2.17.0" or tool.get("source")!=SOURCE or not re.fullmatch(r"[0-9a-f]{64}",str(tool.get("sha256"))): raise ValueError(f"{kind} exact ZAP identity mismatch")
        args=run["expected_arguments"]
        if not isinstance(args,list) or any(not isinstance(x,str) or not x for x in args) or receipt.get("arguments_sha256")!=canonical_sha(args): raise ValueError(f"{kind} arguments are not hash-bound")
        expected_target=dict(target); expected_target.update({"run_kind":kind,"automation_plan_sha256":plan_proof["sha256"]})
        if receipt.get("target")!=expected_target: raise ValueError(f"{kind} target/plan binding mismatch")
        for channel in ("stdout","stderr"): verify_proof(root,receipt.get(channel),f"{kind}.{channel}")
        report_path=safe(root,run["report"],f"{kind}.report"); report=read_json(report_path,f"{kind}.report")
        if str(report.get("@version"))!="2.17.0": raise ValueError(f"{kind} ZAP report version mismatch")
        for alert in alerts(report):
            if alert["riskcode"]>=3 and (alert["pluginid"],alert["riskcode"]) not in accepted_ids: unaccepted.append((kind,alert["pluginid"],alert["riskcode"]))
        evidence.extend([proof(root,receipt_path),plan_proof,proof(root,report_path),receipt["stdout"],receipt["stderr"]])
    if unaccepted: raise ValueError(f"unaccepted high/critical ZAP alerts: {unaccepted}")
    return {"id":"OFFENSIVE_SECURITY","project_id":profile["project_id"],"environment":"production","release_digest":profile["release_digest"],"result":"PASS","executed_at":profile["evaluated_at"],"expires_at":profile["expires_at"],"executor_ref":profile["executor_ref"],"tool":{"name":"OWASP ZAP","version":"2.17.0","source":SOURCE,"digest":"sha256:"+canonical_sha([run["kind"] for run in runs])},"target":target,"assertions":{"zap_automation_pass":True,"authenticated_scope_pass":True,"api_scope_pass":True,"zero_unaccepted_high_critical":True,"manual_review_pass":True},"evidence":evidence}

def main()->int:
    p=argparse.ArgumentParser(); p.add_argument("--project-root",required=True,type=Path); p.add_argument("--profile",required=True); a=p.parse_args()
    try:
        root=a.project_root.resolve(strict=True); profile=read_json(safe(root,a.profile,"profile"),"profile"); output=safe(root,profile.get("output"),"output",exists=False); output.parent.resolve(strict=True).relative_to(root); receipt=validate(profile,root); atomic(output,(json.dumps(receipt,indent=2,sort_keys=True)+"\n").encode()); print(f"ZAP_OFFENSIVE_ADMISSION_PASS project={receipt['project_id']} release={receipt['release_digest']}"); return 0
    except (ValueError,OSError) as error: print(f"ZAP_OFFENSIVE_ADMISSION_FAILED: {error}",file=sys.stderr); return 1
if __name__=="__main__": sys.exit(main())
