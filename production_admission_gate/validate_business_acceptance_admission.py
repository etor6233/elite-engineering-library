from __future__ import annotations
import argparse,json,re,sys
from datetime import datetime
from pathlib import Path
from urllib.parse import urlparse
from validate_k6_load_admission import atomic,canonical_sha,proof,read_json,safe,utc,verify_proof

FIELDS={"schema","project_id","environment","release_digest","evaluated_at","expires_at","executor_ref","tool","target","specification","playwright_run","approvals","defects","owner_signoff","output"}
TARGET_FIELDS={"application_url","business_scope_id","tenant_ids","role_ids"}
RUN_FIELDS={"execution_receipt","observation","expected_arguments"}
PLAYWRIGHT_SOURCE="https://github.com/microsoft/playwright/releases/tag/v1.62.1"
SPEC_KIT_SOURCE="https://github.com/github/spec-kit/releases/tag/v1.0.1"
SPEC_KIT_COMMIT="9118ed15a0ba65053469a94c560ea5d233f75884"

def nonempty(value:object,label:str)->str:
 if not isinstance(value,str) or not value.strip(): raise ValueError(f"{label} is required")
 return value

def unique(value:object,label:str,minimum:int=1)->list[str]:
 if not isinstance(value,list) or len(value)<minimum or len(set(value))!=len(value) or any(not isinstance(x,str) or not x.strip() for x in value): raise ValueError(f"{label} must contain at least {minimum} unique non-empty values")
 return value

def https(value:object,label:str)->None:
 value=nonempty(value,label); p=urlparse(value)
 if p.scheme!="https" or not p.hostname or p.username or p.password or p.fragment: raise ValueError(f"{label} must be credential-free HTTPS")

def bind(value:dict[str,object],profile:dict[str,object],label:str)->None:
 for key in ("project_id","environment","release_digest"):
  if value.get(key)!=profile[key]: raise ValueError(f"{label} {key} mismatch")

def validate(profile:object,root:Path)->dict[str,object]:
 if not isinstance(profile,dict) or set(profile)!=FIELDS: raise ValueError("profile fields must be exact")
 if profile["schema"]!="elite-business-acceptance-admission/v1" or profile["environment"]!="production" or not nonempty(profile["project_id"],"project_id") or not re.fullmatch(r"sha256:[0-9a-f]{64}",str(profile["release_digest"])) or set(str(profile["release_digest"]).split(":")[-1])=={"0"}: raise ValueError("profile identity is invalid")
 start=utc(profile["evaluated_at"],"evaluated_at"); end=utc(profile["expires_at"],"expires_at")
 if not start<end or (end-start).total_seconds()>30*86400: raise ValueError("admission validity must be positive and at most 30 days")
 nonempty(profile["executor_ref"],"executor_ref")
 if profile["tool"]!={"name":"Microsoft Playwright","version":"1.62.1","source":PLAYWRIGHT_SOURCE}: raise ValueError("exact Microsoft Playwright 1.62.1 identity is required")
 target=profile["target"]
 if not isinstance(target,dict) or set(target)!=TARGET_FIELDS: raise ValueError("target fields must be exact")
 https(target["application_url"],"target.application_url"); nonempty(target["business_scope_id"],"target.business_scope_id"); unique(target["tenant_ids"],"target.tenant_ids",1); unique(target["role_ids"],"target.role_ids",1)

 spec_path=safe(root,profile["specification"],"specification"); spec=read_json(spec_path,"specification"); spec_proof=proof(root,spec_path)
 sfields={"schema","project_id","environment","release_digest","approved_at","spec_kit","artifacts","capabilities","scenarios"}
 if not isinstance(spec,dict) or set(spec)!=sfields or spec.get("schema")!="elite-business-acceptance-specification/v1": raise ValueError("specification fields/schema must be exact")
 bind(spec,profile,"specification"); utc(spec["approved_at"],"specification.approved_at")
 if spec["spec_kit"]!={"name":"GitHub Spec Kit","version":"1.0.1","source":SPEC_KIT_SOURCE,"commit":SPEC_KIT_COMMIT}: raise ValueError("exact GitHub Spec Kit 1.0.1 identity is required")
 artifacts=spec["artifacts"]
 required_artifacts={"constitution","spec","plan","tasks","convergence"}
 if not isinstance(artifacts,list) or {x.get("kind") for x in artifacts if isinstance(x,dict)}!=required_artifacts or any(not isinstance(x,dict) or set(x)!={"kind","evidence"} for x in artifacts): raise ValueError("five exact Spec Kit artifacts are required")
 evidence=[spec_proof]
 for item in artifacts:
  verify_proof(root,item["evidence"],f"specification artifact {item['kind']}"); evidence.append(item["evidence"])
 capabilities=unique(spec["capabilities"],"specification.capabilities",1); scenarios=spec["scenarios"]
 sf={"scenario_id","capability_id","actor_role","tenant_id","preconditions","expected_outcomes"}
 if not isinstance(scenarios,list) or not scenarios or any(not isinstance(x,dict) or set(x)!=sf for x in scenarios): raise ValueError("acceptance scenarios must use the exact schema")
 scenario_ids=[]
 for item in scenarios:
  scenario_ids.append(nonempty(item["scenario_id"],"scenario_id"))
  if item["capability_id"] not in capabilities or item["actor_role"] not in target["role_ids"] or item["tenant_id"] not in target["tenant_ids"] or not unique(item["preconditions"],"scenario.preconditions") or not unique(item["expected_outcomes"],"scenario.expected_outcomes"): raise ValueError("scenario is not bound to declared capability/role/tenant/outcomes")
 if len(set(scenario_ids))!=len(scenario_ids) or {x["capability_id"] for x in scenarios}!=set(capabilities): raise ValueError("every capability requires scenarios with unique IDs")

 run=profile["playwright_run"]
 if not isinstance(run,dict) or set(run)!=RUN_FIELDS: raise ValueError("playwright_run fields must be exact")
 obs_path=safe(root,run["observation"],"playwright observation"); obs=read_json(obs_path,"playwright observation"); obs_proof=proof(root,obs_path)
 ofields={"schema","project_id","environment","release_digest","captured_at","application_url","business_scope_id","specification_sha256","results"}
 if not isinstance(obs,dict) or set(obs)!=ofields or obs.get("schema")!="elite-business-acceptance-observation/v1": raise ValueError("business observation fields/schema mismatch")
 bind(obs,profile,"observation"); utc(obs["captured_at"],"observation.captured_at")
 if obs["application_url"]!=target["application_url"] or obs["business_scope_id"]!=target["business_scope_id"] or obs["specification_sha256"]!=spec_proof["sha256"]: raise ValueError("business observation target/specification mismatch")
 results=obs["results"]; rf={"scenario_id","status","expected_outcomes_observed","backend_effect_verified","financial_effect","evidence"}
 if not isinstance(results,list) or {x.get("scenario_id") for x in results if isinstance(x,dict)}!=set(scenario_ids) or len(results)!=len(scenario_ids): raise ValueError("every approved scenario requires exactly one result")
 for item in results:
  if not isinstance(item,dict) or set(item)!=rf or item["status"]!="PASS" or item["expected_outcomes_observed"] is not True or item["backend_effect_verified"] is not True or item["financial_effect"] not in ("VERIFIED","NOT_APPLICABLE"): raise ValueError("acceptance scenario result did not prove its outcomes/effects")
  verify_proof(root,item["evidence"],f"scenario {item.get('scenario_id')} evidence"); evidence.append(item["evidence"])
 receipt_path=safe(root,run["execution_receipt"],"playwright receipt"); receipt=read_json(receipt_path,"playwright receipt")
 if receipt.get("schema")!="elite-official-tool-execution/v1" or receipt.get("project_id")!=profile["project_id"] or receipt.get("environment")!="production" or receipt.get("release_digest")!=profile["release_digest"] or receipt.get("control_id")!="BUSINESS_ACCEPTANCE" or receipt.get("exit_code")!=0: raise ValueError("Playwright execution identity/exit mismatch")
 tool=receipt.get("tool")
 if not isinstance(tool,dict) or set(tool)!={"name","version","source","sha256"} or (tool.get("name"),tool.get("version"),tool.get("source"))!=("Microsoft Playwright","1.62.1",PLAYWRIGHT_SOURCE) or not re.fullmatch(r"[0-9a-f]{64}",str(tool.get("sha256"))) or set(str(tool.get("sha256")))=={"0"}: raise ValueError("exact Microsoft Playwright execution identity is required")
 args=run["expected_arguments"]
 if not isinstance(args,list) or not args or any(not isinstance(x,str) or not x for x in args) or receipt.get("arguments_sha256")!=canonical_sha(args): raise ValueError("Playwright arguments are not hash-bound")
 expected_target=dict(target); expected_target.update({"specification_sha256":spec_proof["sha256"],"observation_sha256":obs_proof["sha256"]})
 if receipt.get("target")!=expected_target: raise ValueError("Playwright target/observation binding mismatch")
 env=receipt.get("environment_variable_names")
 if not isinstance(env,list) or any(re.search(r"(?:secret|password|token|key|credential)",str(x),re.I) and not str(x).endswith("_REF") for x in env): raise ValueError("Playwright environment must expose references, not secrets")
 for channel in ("stdout","stderr"): verify_proof(root,receipt.get(channel),f"playwright.{channel}")
 evidence.extend([obs_proof,proof(root,receipt_path),receipt["stdout"],receipt["stderr"]])

 approvals_path=safe(root,profile["approvals"],"approvals"); approvals=read_json(approvals_path,"approvals"); approvals_proof=proof(root,approvals_path)
 afields={"schema","project_id","environment","release_digest","captured_at","specification_sha256","observation_sha256","approvals"}
 if not isinstance(approvals,dict) or set(approvals)!=afields or approvals.get("schema")!="elite-business-functional-approvals/v1": raise ValueError("functional approvals fields/schema mismatch")
 bind(approvals,profile,"approvals"); utc(approvals["captured_at"],"approvals.captured_at")
 if approvals["specification_sha256"]!=spec_proof["sha256"] or approvals["observation_sha256"]!=obs_proof["sha256"]: raise ValueError("functional approvals evidence binding mismatch")
 required_roles={"finance","operations","security"}; decisions=approvals["approvals"]; df={"role","subject_sha256","decision","approved_at","evidence"}
 if not isinstance(decisions,list) or {x.get("role") for x in decisions if isinstance(x,dict)}!=required_roles or len(decisions)!=3: raise ValueError("finance, operations and security approvals are required")
 subjects=[]
 for item in decisions:
  if not isinstance(item,dict) or set(item)!=df or not re.fullmatch(r"sha256:[0-9a-f]{64}",str(item["subject_sha256"])) or item["decision"]!="APPROVED": raise ValueError("functional approval is invalid")
  utc(item["approved_at"],f"{item.get('role')} approved_at"); verify_proof(root,item["evidence"],f"{item.get('role')} approval evidence"); subjects.append(item["subject_sha256"]); evidence.append(item["evidence"])
 if len(set(subjects))!=3: raise ValueError("finance, operations and security approvers must be distinct")
 evidence.append(approvals_proof)

 defects_path=safe(root,profile["defects"],"defects"); defects=read_json(defects_path,"defects"); defects_proof=proof(root,defects_path)
 dfields={"schema","project_id","environment","release_digest","captured_at","specification_sha256","defects"}
 if not isinstance(defects,dict) or set(defects)!=dfields or defects.get("schema")!="elite-business-defect-register/v1": raise ValueError("defect register fields/schema mismatch")
 bind(defects,profile,"defects"); utc(defects["captured_at"],"defects.captured_at")
 if defects["specification_sha256"]!=spec_proof["sha256"] or not isinstance(defects["defects"],list): raise ValueError("defect register binding/list mismatch")
 seen=set()
 for item in defects["defects"]:
  fields={"defect_id","severity","status","reason","reviewers","expires_at"}
  if not isinstance(item,dict) or set(item)!=fields or not nonempty(item["defect_id"],"defect_id") or item["defect_id"] in seen or item["severity"] not in (1,2,3,4) or item["status"] not in ("CLOSED","ACCEPTED"): raise ValueError("open, duplicate or malformed defect is forbidden")
  seen.add(item["defect_id"])
  if item["severity"] in (1,2) and item["status"]!="CLOSED": raise ValueError("Sev1/Sev2 defects cannot be accepted open")
  if item["status"]=="ACCEPTED":
   nonempty(item["reason"],"accepted defect reason"); unique(item["reviewers"],"accepted defect reviewers",2); expiry=utc(item["expires_at"],"accepted defect expires_at")
   if expiry<=start: raise ValueError("accepted defect expiry must be after evaluation")
 evidence.append(defects_proof)

 signoff_path=safe(root,profile["owner_signoff"],"owner_signoff"); signoff=read_json(signoff_path,"owner_signoff"); signoff_proof=proof(root,signoff_path)
 xfields={"schema","project_id","environment","release_digest","approved_at","role","subject_sha256","decision","specification_sha256","observation_sha256","approvals_sha256","defects_sha256","evidence"}
 if not isinstance(signoff,dict) or set(signoff)!=xfields or signoff.get("schema")!="elite-business-owner-signoff/v1": raise ValueError("owner signoff fields/schema mismatch")
 bind(signoff,profile,"owner signoff"); utc(signoff["approved_at"],"owner signoff approved_at")
 if signoff["role"]!="business_owner" or signoff["decision"]!="APPROVED" or not re.fullmatch(r"sha256:[0-9a-f]{64}",str(signoff["subject_sha256"])) or signoff["subject_sha256"] in subjects: raise ValueError("business owner signoff must be approved and independently owned")
 if (signoff["specification_sha256"],signoff["observation_sha256"],signoff["approvals_sha256"],signoff["defects_sha256"])!=(spec_proof["sha256"],obs_proof["sha256"],approvals_proof["sha256"],defects_proof["sha256"]): raise ValueError("business owner signoff evidence binding mismatch")
 verify_proof(root,signoff["evidence"],"business owner signoff evidence"); evidence.extend([signoff_proof,signoff["evidence"]])
 return {"id":"BUSINESS_ACCEPTANCE","project_id":profile["project_id"],"environment":"production","release_digest":profile["release_digest"],"result":"PASS","executed_at":profile["evaluated_at"],"expires_at":profile["expires_at"],"executor_ref":profile["executor_ref"],"tool":{"name":"GitHub Spec Kit + Microsoft Playwright","version":"1.0.1 + 1.62.1","source":SPEC_KIT_SOURCE,"digest":"sha256:"+canonical_sha([SPEC_KIT_COMMIT,tool["sha256"]])},"target":target,"assertions":{"acceptance_scenarios_pass":True,"finance_operations_security_approved":True,"no_open_sev1_sev2":True,"approvers_distinct":True,"owner_signoff":True},"evidence":evidence}

def main()->int:
 p=argparse.ArgumentParser(); p.add_argument("--project-root",required=True,type=Path); p.add_argument("--profile",required=True); a=p.parse_args()
 try:
  root=a.project_root.resolve(strict=True); profile=read_json(safe(root,a.profile,"profile"),"profile"); output=safe(root,profile.get("output"),"output",exists=False); output.parent.resolve(strict=True).relative_to(root); receipt=validate(profile,root); atomic(output,(json.dumps(receipt,indent=2,sort_keys=True)+"\n").encode()); print(f"BUSINESS_ACCEPTANCE_ADMISSION_PASS project={receipt['project_id']} release={receipt['release_digest']}"); return 0
 except (ValueError,OSError) as error: print(f"BUSINESS_ACCEPTANCE_ADMISSION_FAILED: {error}",file=sys.stderr); return 1
if __name__=="__main__": sys.exit(main())
