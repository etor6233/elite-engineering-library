from __future__ import annotations
import argparse,json,re,sys
from pathlib import Path
from urllib.parse import urlparse
from validate_k6_load_admission import atomic,canonical_sha,proof,read_json,safe,utc,verify_proof

KINDS=("PREFLIGHT","APPLY_IMMUTABLE_DIGEST","CANARY_HEALTH","ROLLOUT_STATUS","ROLLBACK","POST_ROLLBACK_VERIFY","MIGRATION_COMPATIBILITY")
FIELDS={"schema","project_id","environment","release_digest","evaluated_at","expires_at","executor_ref","tool","target","policy","runs","output"}
RUN_FIELDS={"kind","execution_receipt","observation","expected_arguments"}
TARGET_FIELDS={"cluster_uid_sha256","context","namespace","workload_kind","workload_name","container_name","service_url"}
KUBECTL_SOURCE="https://github.com/kubernetes/kubernetes/releases/tag/v1.37.0"
KUBECTL_SHA="4721b614a67bb4932a0369e61f4a323d8c6ca00943d3a2ff14837c124da06f0e"

def nonempty(value:object,label:str)->str:
 if not isinstance(value,str) or not value.strip(): raise ValueError(f"{label} is required")
 return value

def https(value:object,label:str)->None:
 value=nonempty(value,label); parsed=urlparse(value)
 if parsed.scheme!="https" or not parsed.hostname or parsed.username or parsed.password or parsed.fragment: raise ValueError(f"{label} must be credential-free HTTPS")

def bind(value:dict[str,object],profile:dict[str,object],label:str)->None:
 for key in ("project_id","environment","release_digest"):
  if value.get(key)!=profile[key]: raise ValueError(f"{label} {key} mismatch")

def base_observation(value:object,profile:dict[str,object],kind:str,fields:set[str])->dict[str,object]:
 required={"schema","project_id","environment","release_digest","captured_at","kind"}|fields
 if not isinstance(value,dict) or set(value)!=required or value.get("schema")!="elite-deploy-rollback-observation/v1" or value.get("kind")!=kind: raise ValueError(f"{kind} observation fields/schema mismatch")
 bind(value,profile,kind); utc(value["captured_at"],f"{kind}.captured_at")
 return value

def validate_observation(value:object,profile:dict[str,object],policy:dict[str,object],kind:str)->None:
 target=profile["target"]
 if kind=="PREFLIGHT":
  v=base_observation(value,profile,kind,{"client_version","server_version","cluster_uid_sha256","context","namespace","authorized_verbs","workload_found"})
  if v["client_version"]!="v1.37.0" or not re.fullmatch(r"v1\.(3[4-7])\.[0-9]+",str(v["server_version"])): raise ValueError("kubectl/client-server version or supported skew mismatch")
  if v["cluster_uid_sha256"]!=target["cluster_uid_sha256"] or v["context"]!=target["context"] or v["namespace"]!=target["namespace"]: raise ValueError("preflight target mismatch")
  if v["authorized_verbs"]!=["get","list","patch","watch"] or v["workload_found"] is not True: raise ValueError("preflight authorization/workload check failed")
 elif kind=="APPLY_IMMUTABLE_DIGEST":
  v=base_observation(value,profile,kind,{"workload_kind","workload_name","container_name","requested_image_digest","observed_image_digest","revision_before","revision_after","generation","observed_generation"})
  if (v["workload_kind"],v["workload_name"],v["container_name"])!=(target["workload_kind"],target["workload_name"],target["container_name"]): raise ValueError("apply target mismatch")
  if v["requested_image_digest"]!=policy["desired_image_digest"] or v["observed_image_digest"]!=policy["desired_image_digest"]: raise ValueError("immutable desired image digest was not observed")
  if not all(isinstance(v[x],int) and v[x]>=1 for x in ("revision_before","revision_after","generation","observed_generation")) or v["revision_after"]<=v["revision_before"] or v["generation"]!=v["observed_generation"]: raise ValueError("apply revision/generation did not converge")
 elif kind=="CANARY_HEALTH":
  v=base_observation(value,profile,kind,{"service_url","samples","healthy_samples","error_rate","p95_latency_ms","probes_pass","alerts_firing","error_budget_breached"})
  if v["service_url"]!=target["service_url"] or not isinstance(v["samples"],int) or v["samples"]<policy["canary"]["minimum_samples"] or v["healthy_samples"]!=v["samples"]: raise ValueError("canary sample policy failed")
  if not isinstance(v["error_rate"],(int,float)) or v["error_rate"]>policy["canary"]["max_error_rate"] or not isinstance(v["p95_latency_ms"],(int,float)) or v["p95_latency_ms"]>policy["canary"]["max_p95_latency_ms"] or v["probes_pass"] is not True or v["alerts_firing"]!=[] or v["error_budget_breached"] is not False: raise ValueError("canary health/SLO policy failed")
 elif kind=="ROLLOUT_STATUS":
  v=base_observation(value,profile,kind,{"workload_kind","workload_name","progressing","progressing_reason","available","desired_replicas","updated_replicas","available_replicas","old_replicas","application_probes_pass"})
  if (v["workload_kind"],v["workload_name"])!=(target["workload_kind"],target["workload_name"]): raise ValueError("rollout target mismatch")
  if v["progressing"] is not True or v["progressing_reason"]!="NewReplicaSetAvailable" or v["available"] is not True or not isinstance(v["desired_replicas"],int) or v["desired_replicas"]<1 or v["updated_replicas"]!=v["desired_replicas"] or v["available_replicas"]!=v["desired_replicas"] or v["old_replicas"]!=0 or v["application_probes_pass"] is not True: raise ValueError("rollout status/application probes failed")
 elif kind=="ROLLBACK":
  v=base_observation(value,profile,kind,{"workload_kind","workload_name","from_revision","to_revision","rollback_command_issued","rollout_complete"})
  if (v["workload_kind"],v["workload_name"])!=(target["workload_kind"],target["workload_name"]) or v["from_revision"]!=policy["desired_revision"] or v["to_revision"]!=policy["previous_revision"] or v["rollback_command_issued"] is not True or v["rollout_complete"] is not True: raise ValueError("actual rollback execution was not proven")
 elif kind=="POST_ROLLBACK_VERIFY":
  v=base_observation(value,profile,kind,{"observed_image_digest","desired_replicas","available_replicas","service_health_pass","data_invariants_pass","alerts_firing"})
  if v["observed_image_digest"]!=policy["previous_image_digest"] or not isinstance(v["desired_replicas"],int) or v["desired_replicas"]<1 or v["available_replicas"]!=v["desired_replicas"] or v["service_health_pass"] is not True or v["data_invariants_pass"] is not True or v["alerts_firing"]!=[]: raise ValueError("previous digest/health was not restored")
 else:
  v=base_observation(value,profile,kind,{"migration_id","expand_contract_used","forward_compatible","backward_compatible","old_new_versions_coexisted","rollback_data_safe","integrity_checks_pass","destructive_change_before_convergence"})
  if v["migration_id"]!=policy["migration"]["migration_id"] or any(v[x] is not True for x in ("expand_contract_used","forward_compatible","backward_compatible","old_new_versions_coexisted","rollback_data_safe","integrity_checks_pass")) or v["destructive_change_before_convergence"] is not False: raise ValueError("migration compatibility/rollback safety failed")

def validate(profile:object,root:Path)->dict[str,object]:
 if not isinstance(profile,dict) or set(profile)!=FIELDS: raise ValueError("profile fields must be exact")
 if profile["schema"]!="elite-deploy-rollback-admission/v1" or profile["environment"]!="production" or not nonempty(profile["project_id"],"project_id") or not re.fullmatch(r"sha256:[0-9a-f]{64}",str(profile["release_digest"])) or set(str(profile["release_digest"]).split(":")[-1])=={"0"}: raise ValueError("profile identity is invalid")
 start=utc(profile["evaluated_at"],"evaluated_at"); end=utc(profile["expires_at"],"expires_at")
 if not start<end or (end-start).total_seconds()>30*86400: raise ValueError("admission validity must be positive and at most 30 days")
 nonempty(profile["executor_ref"],"executor_ref")
 if profile["tool"]!={"name":"Kubernetes kubectl","version":"1.37.0","source":KUBECTL_SOURCE,"sha256":KUBECTL_SHA}: raise ValueError("exact official kubectl 1.37.0 identity is required")
 target=profile["target"]
 if not isinstance(target,dict) or set(target)!=TARGET_FIELDS: raise ValueError("target fields must be exact")
 if not re.fullmatch(r"sha256:[0-9a-f]{64}",str(target["cluster_uid_sha256"])): raise ValueError("cluster_uid_sha256 must be sha256")
 for key in ("context","namespace","workload_name","container_name"): nonempty(target[key],f"target.{key}")
 if target["workload_kind"]!="Deployment": raise ValueError("only Kubernetes Deployment rollout semantics are supported")
 https(target["service_url"],"target.service_url")
 policy_path=safe(root,profile["policy"],"policy"); policy=read_json(policy_path,"policy")
 pfields={"schema","project_id","environment","release_digest","approved_at","approvers","target","desired_image_digest","previous_image_digest","desired_revision","previous_revision","timeout_seconds","canary","migration"}
 if not isinstance(policy,dict) or set(policy)!=pfields or policy.get("schema")!="elite-deploy-rollback-policy/v1": raise ValueError("deploy policy fields/schema must be exact")
 bind(policy,profile,"policy"); utc(policy["approved_at"],"policy.approved_at")
 if not isinstance(policy["approvers"],list) or len(policy["approvers"])<2 or len(set(policy["approvers"]))!=len(policy["approvers"]) or any(not nonempty(x,"policy approver") for x in policy["approvers"]): raise ValueError("deploy policy requires at least two distinct approvers")
 if policy["target"]!=target: raise ValueError("deploy policy target mismatch")
 for key in ("desired_image_digest","previous_image_digest"):
  if not re.fullmatch(r"[^\s@]+@sha256:[0-9a-f]{64}",str(policy[key])): raise ValueError(f"{key} must be an immutable registry digest")
 if policy["desired_image_digest"]==policy["previous_image_digest"]: raise ValueError("desired and previous image digests must differ")
 if not all(isinstance(policy[x],int) and policy[x]>=1 for x in ("desired_revision","previous_revision","timeout_seconds")) or policy["desired_revision"]<=policy["previous_revision"] or policy["timeout_seconds"]>3600: raise ValueError("revision/timeout policy is invalid")
 canary=policy["canary"]
 if not isinstance(canary,dict) or set(canary)!={"minimum_samples","max_error_rate","max_p95_latency_ms","automatic_rollback"} or not isinstance(canary["minimum_samples"],int) or canary["minimum_samples"]<3 or not isinstance(canary["max_error_rate"],(int,float)) or not 0<=canary["max_error_rate"]<=0.05 or not isinstance(canary["max_p95_latency_ms"],int) or canary["max_p95_latency_ms"]<1 or canary["automatic_rollback"] is not True: raise ValueError("canary policy is invalid")
 migration=policy["migration"]
 if not isinstance(migration,dict) or set(migration)!={"migration_id","expand_contract_required","forward_backward_required","rollback_data_recovery_ref"} or not nonempty(migration["migration_id"],"migration_id") or migration["expand_contract_required"] is not True or migration["forward_backward_required"] is not True or not nonempty(migration["rollback_data_recovery_ref"],"rollback_data_recovery_ref"): raise ValueError("migration policy is incomplete")
 runs=profile["runs"]
 if not isinstance(runs,list) or len(runs)!=len(KINDS) or any(not isinstance(x,dict) or set(x)!=RUN_FIELDS for x in runs) or [x["kind"] for x in runs]!=list(KINDS): raise ValueError("seven ordered deploy/rollback runs are required")
 evidence=[proof(root,policy_path)]; tool_hash=None
 for run in runs:
  kind=run["kind"]; observation_path=safe(root,run["observation"],f"{kind}.observation"); observation=read_json(observation_path,f"{kind}.observation"); validate_observation(observation,profile,policy,kind); observation_proof=proof(root,observation_path)
  receipt_path=safe(root,run["execution_receipt"],f"{kind}.receipt"); receipt=read_json(receipt_path,f"{kind}.receipt")
  if receipt.get("schema")!="elite-official-tool-execution/v1" or receipt.get("project_id")!=profile["project_id"] or receipt.get("environment")!="production" or receipt.get("release_digest")!=profile["release_digest"] or receipt.get("control_id")!="DEPLOY_ROLLBACK" or receipt.get("exit_code")!=0: raise ValueError(f"{kind} execution identity/exit mismatch")
  tool=receipt.get("tool")
  if tool!={"name":"Kubernetes kubectl","version":"1.37.0","source":KUBECTL_SOURCE,"sha256":KUBECTL_SHA}: raise ValueError(f"{kind} exact kubectl identity mismatch")
  tool_hash=tool_hash or tool["sha256"]
  if tool["sha256"]!=tool_hash: raise ValueError("all deploy/rollback runs must use the same kubectl hash")
  args=run["expected_arguments"]
  if not isinstance(args,list) or not args or any(not isinstance(x,str) or not x for x in args) or receipt.get("arguments_sha256")!=canonical_sha(args): raise ValueError(f"{kind} arguments are not hash-bound")
  expected_target=dict(target); expected_target.update({"run_kind":kind,"observation_sha256":observation_proof["sha256"]})
  if receipt.get("target")!=expected_target: raise ValueError(f"{kind} target/observation binding mismatch")
  env=receipt.get("environment_variable_names")
  if not isinstance(env,list) or any(re.search(r"(?:secret|password|token|key|credential|kubeconfig)",str(x),re.I) and not str(x).endswith("_REF") for x in env): raise ValueError(f"{kind} environment must expose references, not secrets")
  for channel in ("stdout","stderr"): verify_proof(root,receipt.get(channel),f"{kind}.{channel}")
  evidence.extend([proof(root,receipt_path),observation_proof,receipt["stdout"],receipt["stderr"]])
 return {"id":"DEPLOY_ROLLBACK","project_id":profile["project_id"],"environment":"production","release_digest":profile["release_digest"],"result":"PASS","executed_at":profile["evaluated_at"],"expires_at":profile["expires_at"],"executor_ref":profile["executor_ref"],"tool":{"name":"Kubernetes kubectl","version":"1.37.0","source":KUBECTL_SOURCE,"digest":"sha256:"+KUBECTL_SHA},"target":target,"assertions":{"immutable_digest_deployed":True,"canary_health_pass":True,"rollout_status_pass":True,"rollback_executed":True,"previous_digest_restored":True,"migration_compatibility_pass":True},"evidence":evidence}

def main()->int:
 p=argparse.ArgumentParser(); p.add_argument("--project-root",required=True,type=Path); p.add_argument("--profile",required=True); a=p.parse_args()
 try:
  root=a.project_root.resolve(strict=True); profile=read_json(safe(root,a.profile,"profile"),"profile"); output=safe(root,profile.get("output"),"output",exists=False); output.parent.resolve(strict=True).relative_to(root); receipt=validate(profile,root); atomic(output,(json.dumps(receipt,indent=2,sort_keys=True)+"\n").encode()); print(f"DEPLOY_ROLLBACK_ADMISSION_PASS project={receipt['project_id']} release={receipt['release_digest']}"); return 0
 except (ValueError,OSError) as error: print(f"DEPLOY_ROLLBACK_ADMISSION_FAILED: {error}",file=sys.stderr); return 1
if __name__=="__main__": sys.exit(main())
