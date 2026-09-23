from __future__ import annotations
import argparse,json,re,sys
from itertools import permutations
from pathlib import Path
from urllib.parse import urlparse
from validate_k6_load_admission import atomic,canonical_sha,proof,read_json,safe,utc,verify_proof

KINDS=("OIDC_CONFORMANCE","ROLE_POSITIVE_JOURNEYS","UNAUTHORIZED_NEGATIVE_JOURNEYS","TENANT_ISOLATION","SESSION_REVOCATION_ROTATION")
FIELDS={"schema","project_id","environment","release_digest","evaluated_at","expires_at","executor_ref","tool","target","policy","runs","output"}
RUN_FIELDS={"kind","execution_receipt","observation","expected_arguments"}
TARGET_FIELDS={"issuer_url","client_id","audience","conformance_profile","application_url","tenant_ids","role_ids"}
OIDF_SOURCE="https://gitlab.com/openid/conformance-suite/-/releases/release-v5.2.4"
PLAYWRIGHT_SOURCE="https://github.com/microsoft/playwright/releases/tag/v1.62.1"

def nonempty(value:object,label:str)->str:
 if not isinstance(value,str) or not value.strip(): raise ValueError(f"{label} is required")
 return value

def https(value:object,label:str)->None:
 value=nonempty(value,label); parsed=urlparse(value)
 if parsed.scheme!="https" or not parsed.hostname or parsed.username or parsed.password or parsed.fragment: raise ValueError(f"{label} must be credential-free HTTPS")

def unique_strings(value:object,label:str,minimum:int=1)->list[str]:
 if not isinstance(value,list) or len(value)<minimum or len(set(value))!=len(value) or any(not isinstance(x,str) or not x.strip() for x in value): raise ValueError(f"{label} must contain at least {minimum} unique non-empty values")
 return value

def bind(value:dict[str,object],profile:dict[str,object],label:str)->None:
 for key in ("project_id","environment","release_digest"):
  if value.get(key)!=profile[key]: raise ValueError(f"{label} {key} mismatch")

def validate_oidc(value:dict[str,object],profile:dict[str,object])->None:
 fields={"schema","project_id","environment","release_digest","captured_at","kind","issuer_url","client_id","conformance_profile","suite_release","plan_id","summary","non_passed"}
 if set(value)!=fields or value.get("schema")!="elite-identity-observation/v1" or value.get("kind")!="OIDC_CONFORMANCE": raise ValueError("OIDC conformance observation fields/schema mismatch")
 bind(value,profile,"OIDC conformance"); utc(value["captured_at"],"OIDC captured_at"); target=profile["target"]
 if value["issuer_url"]!=target["issuer_url"] or value["client_id"]!=target["client_id"] or value["conformance_profile"]!=target["conformance_profile"] or value["suite_release"]!="release-v5.2.4" or not nonempty(value["plan_id"],"OIDC plan_id"): raise ValueError("OIDC exact target/profile/release mismatch")
 summary=value["summary"]
 if not isinstance(summary,dict) or set(summary)!={"passed","review","warning","skipped","failed","interrupted","total"} or any(not isinstance(v,int) or v<0 for v in summary.values()): raise ValueError("OIDC summary must contain exact non-negative counters")
 if summary["total"]<=0 or summary["total"]!=sum(summary[k] for k in ("passed","review","warning","skipped","failed","interrupted")) or summary["failed"] or summary["interrupted"]: raise ValueError("OIDC conformance contains failed/interrupted or inconsistent results")
 exceptions=value["non_passed"]
 if not isinstance(exceptions,list) or len(exceptions)!=summary["review"]+summary["warning"]+summary["skipped"]: raise ValueError("every non-passed OIDC result requires an approved disposition")
 seen=set()
 for item in exceptions:
  if not isinstance(item,dict) or set(item)!={"test_id","status","reason","reviewers"} or item["status"] not in ("REVIEW","WARNING","SKIPPED") or not nonempty(item["test_id"],"OIDC test_id") or item["test_id"] in seen or not nonempty(item["reason"],"OIDC disposition reason") or len(unique_strings(item["reviewers"],"OIDC disposition reviewers",2))<2: raise ValueError("OIDC non-passed disposition is invalid")
  seen.add(item["test_id"])

def validate_roles(value:dict[str,object],profile:dict[str,object])->None:
 fields={"schema","project_id","environment","release_digest","captured_at","kind","application_url","journeys"}
 if set(value)!=fields or value.get("schema")!="elite-identity-observation/v1" or value.get("kind")!="ROLE_POSITIVE_JOURNEYS": raise ValueError("role observation fields/schema mismatch")
 bind(value,profile,"roles"); utc(value["captured_at"],"roles captured_at")
 if value["application_url"]!=profile["target"]["application_url"]: raise ValueError("role target mismatch")
 journeys=value["journeys"]
 if not isinstance(journeys,list) or {x.get("role_id") for x in journeys if isinstance(x,dict)}!=set(profile["target"]["role_ids"]): raise ValueError("every declared role requires a positive journey")
 for item in journeys:
  if set(item)!={"role_id","tenant_id","journey_id","authenticated","status","expected_backend_effect","effect_observed"} or item["tenant_id"] not in profile["target"]["tenant_ids"] or not nonempty(item["journey_id"],"role journey_id") or item["authenticated"] is not True or not isinstance(item["status"],int) or not 200<=item["status"]<300 or item["expected_backend_effect"] is not True or item["effect_observed"] is not True: raise ValueError("positive role journey did not prove its backend effect")

def validate_unauthorized(value:dict[str,object],profile:dict[str,object])->None:
 required={"unauthenticated","expired_token","wrong_issuer","wrong_audience","insufficient_role","object_ownership"}
 fields={"schema","project_id","environment","release_digest","captured_at","kind","application_url","scenarios"}
 if set(value)!=fields or value.get("schema")!="elite-identity-observation/v1" or value.get("kind")!="UNAUTHORIZED_NEGATIVE_JOURNEYS": raise ValueError("unauthorized observation fields/schema mismatch")
 bind(value,profile,"unauthorized"); utc(value["captured_at"],"unauthorized captured_at")
 if value["application_url"]!=profile["target"]["application_url"]: raise ValueError("unauthorized target mismatch")
 scenarios=value["scenarios"]
 if not isinstance(scenarios,list) or {x.get("scenario") for x in scenarios if isinstance(x,dict)}!=required: raise ValueError("all mandatory unauthorized scenarios are required")
 for item in scenarios:
  if set(item)!={"scenario","status","data_disclosed","mutation_observed"} or item["status"] not in (401,403,404) or item["data_disclosed"] is not False or item["mutation_observed"] is not False: raise ValueError("unauthorized journey leaked data or changed state")

def validate_tenants(value:dict[str,object],profile:dict[str,object])->None:
 fields={"schema","project_id","environment","release_digest","captured_at","kind","tenant_ids","scenarios"}
 if set(value)!=fields or value.get("schema")!="elite-identity-observation/v1" or value.get("kind")!="TENANT_ISOLATION": raise ValueError("tenant observation fields/schema mismatch")
 bind(value,profile,"tenants"); utc(value["captured_at"],"tenants captured_at"); tenants=profile["target"]["tenant_ids"]
 if value["tenant_ids"]!=tenants: raise ValueError("tenant set/order mismatch")
 scenarios=value["scenarios"]; expected={(a,b,action) for a,b in permutations(tenants,2) for action in ("read","write","list")}
 if not isinstance(scenarios,list) or len(scenarios)!=len(expected): raise ValueError("complete ordered cross-tenant matrix is required")
 actual=set()
 for item in scenarios:
  if not isinstance(item,dict) or set(item)!={"actor_tenant","target_tenant","action","resource_key","status","data_disclosed","mutation_observed"} or not nonempty(item["resource_key"],"tenant resource_key"): raise ValueError("tenant scenario fields are invalid")
  key=(item["actor_tenant"],item["target_tenant"],item["action"]); actual.add(key)
  if item["actor_tenant"]==item["target_tenant"] or item["status"] not in (403,404) or item["data_disclosed"] is not False or item["mutation_observed"] is not False: raise ValueError("cross-tenant isolation failed")
 if actual!=expected: raise ValueError("complete ordered cross-tenant matrix is required")

def validate_sessions(value:dict[str,object],profile:dict[str,object])->None:
 required={"logout_revocation","token_revocation","signing_key_rotation","session_expiry","replay_rejection"}
 fields={"schema","project_id","environment","release_digest","captured_at","kind","issuer_url","cases"}
 if set(value)!=fields or value.get("schema")!="elite-identity-observation/v1" or value.get("kind")!="SESSION_REVOCATION_ROTATION": raise ValueError("session observation fields/schema mismatch")
 bind(value,profile,"sessions"); utc(value["captured_at"],"sessions captured_at")
 if value["issuer_url"]!=profile["target"]["issuer_url"]: raise ValueError("session issuer mismatch")
 cases=value["cases"]
 if not isinstance(cases,list) or {x.get("case") for x in cases if isinstance(x,dict)}!=required: raise ValueError("all revocation/rotation/session cases are required")
 for item in cases:
  if set(item)!={"case","status","old_credential_accepted","rejected","audit_event_observed"} or item["status"] not in (401,403) or item["old_credential_accepted"] is not False or item["rejected"] is not True or item["audit_event_observed"] is not True: raise ValueError("session revocation/rotation assertion failed")

def validate(profile:object,root:Path)->dict[str,object]:
 if not isinstance(profile,dict) or set(profile)!=FIELDS: raise ValueError("profile fields must be exact")
 if profile["schema"]!="elite-identity-authorization-admission/v1" or profile["environment"]!="production" or not nonempty(profile["project_id"],"project_id") or not re.fullmatch(r"sha256:[0-9a-f]{64}",str(profile["release_digest"])): raise ValueError("profile identity is invalid")
 start=utc(profile["evaluated_at"],"evaluated_at"); end=utc(profile["expires_at"],"expires_at")
 if not start<end or (end-start).total_seconds()>30*86400: raise ValueError("admission validity must be positive and at most 30 days")
 nonempty(profile["executor_ref"],"executor_ref")
 if profile["tool"]!={"name":"OpenID Foundation Conformance Suite","version":"5.2.4","source":OIDF_SOURCE}: raise ValueError("exact OpenID Foundation Conformance Suite 5.2.4 identity is required")
 target=profile["target"]
 if not isinstance(target,dict) or set(target)!=TARGET_FIELDS: raise ValueError("target fields must be exact")
 https(target["issuer_url"],"target.issuer_url"); https(target["application_url"],"target.application_url")
 for key in ("client_id","audience","conformance_profile"): nonempty(target[key],f"target.{key}")
 unique_strings(target["tenant_ids"],"target.tenant_ids",2); unique_strings(target["role_ids"],"target.role_ids",1)
 policy_path=safe(root,profile["policy"],"policy"); policy=read_json(policy_path,"policy")
 policy_fields={"schema","project_id","environment","release_digest","approved_at","approvers","issuer_url","client_id","audience","conformance_profile","tenant_ids","role_ids","break_glass"}
 if not isinstance(policy,dict) or set(policy)!=policy_fields or policy.get("schema")!="elite-identity-authorization-policy/v1": raise ValueError("identity policy fields/schema must be exact")
 bind(policy,profile,"policy"); utc(policy["approved_at"],"policy.approved_at"); unique_strings(policy["approvers"],"policy.approvers",2)
 for key in ("issuer_url","client_id","audience","conformance_profile","tenant_ids","role_ids"):
  if policy[key]!=target[key]: raise ValueError(f"identity policy {key} mismatch")
 bg=policy["break_glass"]
 if not isinstance(bg,dict) or set(bg)!={"configured","max_minutes","two_person_approval","audit_required","post_use_review_required"} or bg["configured"] is not True or not isinstance(bg["max_minutes"],int) or not 1<=bg["max_minutes"]<=60 or any(bg[k] is not True for k in ("two_person_approval","audit_required","post_use_review_required")): raise ValueError("break-glass governance is incomplete")
 runs=profile["runs"]
 if not isinstance(runs,list) or len(runs)!=5 or any(not isinstance(x,dict) or set(x)!=RUN_FIELDS for x in runs) or [x["kind"] for x in runs]!=list(KINDS): raise ValueError("five ordered identity runs are required")
 evidence=[proof(root,policy_path)]; tool_hashes={}
 validators={"OIDC_CONFORMANCE":validate_oidc,"ROLE_POSITIVE_JOURNEYS":validate_roles,"UNAUTHORIZED_NEGATIVE_JOURNEYS":validate_unauthorized,"TENANT_ISOLATION":validate_tenants,"SESSION_REVOCATION_ROTATION":validate_sessions}
 for run in runs:
  kind=run["kind"]; observation_path=safe(root,run["observation"],f"{kind}.observation"); observation=read_json(observation_path,f"{kind}.observation"); validators[kind](observation,profile); observation_proof=proof(root,observation_path)
  receipt_path=safe(root,run["execution_receipt"],f"{kind}.receipt"); receipt=read_json(receipt_path,f"{kind}.receipt")
  if receipt.get("schema")!="elite-official-tool-execution/v1" or receipt.get("project_id")!=profile["project_id"] or receipt.get("environment")!="production" or receipt.get("release_digest")!=profile["release_digest"] or receipt.get("control_id")!="IDENTITY_AUTHORIZATION" or receipt.get("exit_code")!=0: raise ValueError(f"{kind} execution identity/exit mismatch")
  tool=receipt.get("tool"); expected=("OpenID Foundation Conformance Runner","5.2.4",OIDF_SOURCE) if kind=="OIDC_CONFORMANCE" else ("Microsoft Playwright","1.62.1",PLAYWRIGHT_SOURCE)
  if not isinstance(tool,dict) or set(tool)!={"name","version","source","sha256"} or (tool.get("name"),tool.get("version"),tool.get("source"))!=expected or not re.fullmatch(r"[0-9a-f]{64}",str(tool.get("sha256"))): raise ValueError(f"{kind} exact official tool identity mismatch")
  family="oidf" if kind=="OIDC_CONFORMANCE" else "playwright"; tool_hashes.setdefault(family,tool["sha256"])
  if tool_hashes[family]!=tool["sha256"]: raise ValueError(f"all {family} runs must use the same binary/source hash")
  args=run["expected_arguments"]
  if not isinstance(args,list) or not args or any(not isinstance(x,str) or not x for x in args) or receipt.get("arguments_sha256")!=canonical_sha(args): raise ValueError(f"{kind} arguments are not hash-bound")
  expected_target=dict(target); expected_target.update({"run_kind":kind,"observation_sha256":observation_proof["sha256"]})
  if receipt.get("target")!=expected_target: raise ValueError(f"{kind} target/observation binding mismatch")
  if not isinstance(receipt.get("environment_variable_names"),list) or any(re.search(r"(?:secret|password|token|key|credential)",str(x),re.I) and not str(x).endswith("_REF") for x in receipt["environment_variable_names"]): raise ValueError(f"{kind} environment must expose references, not secret-bearing names")
  for channel in ("stdout","stderr"): verify_proof(root,receipt.get(channel),f"{kind}.{channel}")
  evidence.extend([proof(root,receipt_path),observation_proof,receipt["stdout"],receipt["stderr"]])
 return {"id":"IDENTITY_AUTHORIZATION","project_id":profile["project_id"],"environment":"production","release_digest":profile["release_digest"],"result":"PASS","executed_at":profile["evaluated_at"],"expires_at":profile["expires_at"],"executor_ref":profile["executor_ref"],"tool":{"name":"OpenID Foundation Conformance Suite + Microsoft Playwright","version":"5.2.4 + 1.62.1","source":OIDF_SOURCE,"digest":"sha256:"+canonical_sha([tool_hashes["oidf"],tool_hashes["playwright"]])},"target":target,"assertions":{"oidc_conformance_pass":True,"role_positive_journeys_pass":True,"unauthorized_negative_journeys_pass":True,"tenant_isolation_pass":True,"session_revocation_rotation_pass":True},"evidence":evidence}

def main()->int:
 p=argparse.ArgumentParser(); p.add_argument("--project-root",required=True,type=Path); p.add_argument("--profile",required=True); a=p.parse_args()
 try:
  root=a.project_root.resolve(strict=True); profile=read_json(safe(root,a.profile,"profile"),"profile"); output=safe(root,profile.get("output"),"output",exists=False); output.parent.resolve(strict=True).relative_to(root); receipt=validate(profile,root); atomic(output,(json.dumps(receipt,indent=2,sort_keys=True)+"\n").encode()); print(f"IDENTITY_AUTHORIZATION_ADMISSION_PASS project={receipt['project_id']} release={receipt['release_digest']}"); return 0
 except (ValueError,OSError) as error: print(f"IDENTITY_AUTHORIZATION_ADMISSION_FAILED: {error}",file=sys.stderr); return 1
if __name__=="__main__": sys.exit(main())
