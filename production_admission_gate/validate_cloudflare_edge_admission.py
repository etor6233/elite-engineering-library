from __future__ import annotations
import argparse,json,re,sys
from pathlib import Path
from urllib.parse import urlparse
from validate_k6_load_admission import atomic,canonical_sha,proof,read_json,safe,utc,verify_proof

KINDS=("CLOUDFLARE_CONFIG_SNAPSHOT","ORIGIN_BYPASS_NEGATIVE","PUBLIC_CACHE_BEHAVIOR","PRIVATE_CACHE_BYPASS")
FIELDS={"schema","project_id","environment","release_digest","evaluated_at","expires_at","executor_ref","tool","target","policy","runs","output"}
RUN_FIELDS={"kind","execution_receipt","observation","expected_arguments"}
TARGET_FIELDS={"zone_id","hostname","origin_address","public_url","private_url"}
SOURCE="https://github.com/cloudflare/cloudflare-go/releases/tag/v7.9.0"

def nonempty(value:object,label:str)->str:
 if not isinstance(value,str) or not value.strip(): raise ValueError(f"{label} is required")
 return value

def https(value:object,label:str)->None:
 value=nonempty(value,label); parsed=urlparse(value)
 if parsed.scheme!="https" or not parsed.hostname or parsed.username or parsed.password or parsed.fragment: raise ValueError(f"{label} must be credential-free HTTPS")

def bind(value:dict[str,object],profile:dict[str,object],label:str)->None:
 for key in ("project_id","environment","release_digest"):
  if value.get(key)!=profile[key]: raise ValueError(f"{label} {key} mismatch")

def validate_snapshot(value:dict[str,object],profile:dict[str,object],policy:dict[str,object])->None:
 required={"schema","project_id","environment","release_digest","captured_at","zone_id","hostname","dns_records","waf_rulesets","tls"}
 if set(value)!=required or value.get("schema")!="elite-cloudflare-edge-snapshot/v1": raise ValueError("Cloudflare snapshot fields/schema must be exact")
 bind(value,profile,"snapshot"); utc(value["captured_at"],"snapshot.captured_at")
 target=profile["target"]
 if value["zone_id"]!=target["zone_id"] or value["hostname"]!=target["hostname"]: raise ValueError("snapshot target mismatch")
 records=value["dns_records"]
 if not isinstance(records,list) or not records: raise ValueError("DNS records are required")
 address=[]
 for item in records:
  if not isinstance(item,dict) or set(item)!={"id","name","type","content","proxied"}: raise ValueError("DNS record fields must be exact")
  if item["name"]==target["hostname"] and item["type"] in ("A","AAAA","CNAME"): address.append(item)
 if not address or any(item["proxied"] is not True or not nonempty(item["id"],"DNS id") or not nonempty(item["content"],"DNS content") for item in address): raise ValueError("all hostname address records must be proxied")
 rules=value["waf_rulesets"]
 if not isinstance(rules,list) or not rules: raise ValueError("WAF rulesets are required")
 normalized={}
 for item in rules:
  if not isinstance(item,dict) or set(item)!={"id","name","kind","phase","version","status","enabled_actions"}: raise ValueError("WAF ruleset fields must be exact")
  if item["status"]!="active" or not isinstance(item["enabled_actions"],list) or not item["enabled_actions"]: raise ValueError("WAF rulesets must be active with enabled actions")
  normalized[item["id"]]=item
 managed=policy["managed_ruleset_ids"]; custom=policy["custom_ruleset_ids"]
 if not managed or not custom or any(i not in normalized for i in managed+custom): raise ValueError("approved managed and custom WAF rulesets are required")
 if any(normalized[i]["phase"]!="http_request_firewall_managed" or "execute" not in normalized[i]["enabled_actions"] for i in managed): raise ValueError("managed WAF execute ruleset is not active")
 defensive={"block","challenge","managed_challenge","js_challenge"}
 if any(normalized[i]["phase"]!="http_request_firewall_custom" or not defensive.intersection(normalized[i]["enabled_actions"]) for i in custom): raise ValueError("custom WAF defensive ruleset is not active")
 tls=value["tls"]
 if not isinstance(tls,dict) or set(tls)!={"mode","minimum_version","always_use_https","certificate_status"}: raise ValueError("TLS fields must be exact")
 if tls!={"mode":"full_strict","minimum_version":policy["minimum_tls_version"],"always_use_https":True,"certificate_status":"active"} or tls["minimum_version"] not in ("1.2","1.3"): raise ValueError("Full (strict) TLS policy is not proven")

def validate_observation(kind:str,value:dict[str,object],profile:dict[str,object],policy:dict[str,object])->None:
 fields={"schema","project_id","environment","release_digest","captured_at","kind","target","result"}
 if set(value)!=fields or value.get("schema")!="elite-cloudflare-edge-observation/v1" or value.get("kind")!=kind: raise ValueError(f"{kind} observation fields/schema mismatch")
 bind(value,profile,kind); utc(value["captured_at"],f"{kind}.captured_at")
 result=value["result"]
 if not isinstance(result,dict): raise ValueError(f"{kind} result must be an object")
 if kind=="ORIGIN_BYPASS_NEGATIVE":
  if value["target"]!=profile["target"]["origin_address"] or set(result)!={"connection_succeeded","http_status","blocked","protection_mode"}: raise ValueError("origin bypass observation mismatch")
  if result["blocked"] is not True or result["protection_mode"]!=policy["origin_protection_mode"] or (result["connection_succeeded"] is True and result["http_status"] not in (400,401,403,421,525,526)): raise ValueError("origin bypass was not blocked")
 elif kind=="PUBLIC_CACHE_BEHAVIOR":
  if value["target"]!=profile["target"]["public_url"] or set(result)!={"authorization_sent","first_status","second_status","first_cf_cache_status","second_cf_cache_status","cache_control","first_body_sha256","second_body_sha256"}: raise ValueError("public cache observation mismatch")
  if result["authorization_sent"] is not False or result["first_status"]!=200 or result["second_status"]!=200 or result["second_cf_cache_status"]!="HIT" or result["first_body_sha256"]!=result["second_body_sha256"] or not re.fullmatch(r"[0-9a-f]{64}",str(result["first_body_sha256"])) or re.search(r"(?:private|no-store)",str(result["cache_control"]),re.I): raise ValueError("public cache behavior is not proven")
 else:
  if value["target"]!=profile["target"]["private_url"] or set(result)!={"authorization_sent","status","cf_cache_status","cache_control","set_cookie","shared_response_reused"}: raise ValueError("private cache observation mismatch")
  if result["authorization_sent"] is not True or not isinstance(result["status"],int) or not 200<=result["status"]<400 or result["cf_cache_status"] not in ("BYPASS","DYNAMIC") or not re.search(r"(?:private|no-store)",str(result["cache_control"]),re.I) or result["shared_response_reused"] is not False: raise ValueError("private response cache isolation is not proven")

def validate(profile:object,root:Path)->dict[str,object]:
 if not isinstance(profile,dict) or set(profile)!=FIELDS: raise ValueError("profile fields must be exact")
 if profile["schema"]!="elite-cloudflare-edge-admission/v1" or profile["environment"]!="production" or not nonempty(profile["project_id"],"project_id") or not re.fullmatch(r"sha256:[0-9a-f]{64}",str(profile["release_digest"])): raise ValueError("profile identity is invalid")
 start=utc(profile["evaluated_at"],"evaluated_at"); end=utc(profile["expires_at"],"expires_at")
 if not start<end or (end-start).total_seconds()>30*86400: raise ValueError("admission validity must be positive and at most 30 days")
 nonempty(profile["executor_ref"],"executor_ref")
 if profile["tool"]!={"name":"Cloudflare cloudflare-go","version":"7.9.0","source":SOURCE}: raise ValueError("exact Cloudflare cloudflare-go v7.9.0 identity is required")
 target=profile["target"]
 if not isinstance(target,dict) or set(target)!=TARGET_FIELDS: raise ValueError("target fields must be exact")
 for key in ("zone_id","hostname","origin_address"): nonempty(target[key],f"target.{key}")
 if target["origin_address"]==target["hostname"]: raise ValueError("origin must be distinct from public hostname")
 https(target["public_url"],"target.public_url"); https(target["private_url"],"target.private_url")
 if urlparse(target["public_url"]).hostname!=target["hostname"] or urlparse(target["private_url"]).hostname!=target["hostname"]: raise ValueError("edge URLs must use the admitted hostname")
 policy_path=safe(root,profile["policy"],"policy"); policy=read_json(policy_path,"policy")
 policy_fields={"schema","project_id","environment","release_digest","approved_at","approvers","zone_id","hostname","managed_ruleset_ids","custom_ruleset_ids","minimum_tls_version","origin_protection_mode","public_cache_path","private_cache_path"}
 if set(policy)!=policy_fields or policy.get("schema")!="elite-cloudflare-edge-policy/v1": raise ValueError("edge policy fields/schema must be exact")
 bind(policy,profile,"policy"); utc(policy["approved_at"],"policy.approved_at")
 approvers=policy["approvers"]
 if not isinstance(approvers,list) or len(approvers)<2 or len(set(approvers))!=len(approvers) or any(not isinstance(x,str) or not x.strip() for x in approvers): raise ValueError("edge policy requires two distinct approvers")
 if policy["zone_id"]!=target["zone_id"] or policy["hostname"]!=target["hostname"] or policy["origin_protection_mode"] not in ("authenticated_origin_pull","mtls","cloudflare_tunnel","private_network"): raise ValueError("edge policy target/origin protection mismatch")
 for key in ("managed_ruleset_ids","custom_ruleset_ids"):
  if not isinstance(policy[key],list) or len(set(policy[key]))!=len(policy[key]) or any(not isinstance(x,str) or not x for x in policy[key]): raise ValueError(f"{key} must contain unique IDs")
 if urlparse(target["public_url"]).path!=policy["public_cache_path"] or urlparse(target["private_url"]).path!=policy["private_cache_path"]: raise ValueError("approved cache paths mismatch")
 runs=profile["runs"]
 if not isinstance(runs,list) or len(runs)!=4 or any(not isinstance(x,dict) or set(x)!=RUN_FIELDS for x in runs) or [x["kind"] for x in runs]!=list(KINDS): raise ValueError("four ordered edge runs are required")
 evidence=[proof(root,policy_path)]; digest=""
 for run in runs:
  kind=run["kind"]; observation_path=safe(root,run["observation"],f"{kind}.observation"); observation=read_json(observation_path,f"{kind}.observation")
  if kind=="CLOUDFLARE_CONFIG_SNAPSHOT": validate_snapshot(observation,profile,policy)
  else: validate_observation(kind,observation,profile,policy)
  observation_proof=proof(root,observation_path); receipt_path=safe(root,run["execution_receipt"],f"{kind}.receipt"); receipt=read_json(receipt_path,f"{kind}.receipt")
  if receipt.get("schema")!="elite-official-tool-execution/v1" or receipt.get("project_id")!=profile["project_id"] or receipt.get("environment")!="production" or receipt.get("release_digest")!=profile["release_digest"] or receipt.get("control_id")!="EDGE_CDN_WAF" or receipt.get("exit_code")!=0: raise ValueError(f"{kind} execution identity/exit mismatch")
  tool=receipt.get("tool")
  if not isinstance(tool,dict) or set(tool)!={"name","version","source","sha256"} or tool.get("name")!="Project Cloudflare edge probe" or tool.get("version")!="cloudflare-go/7.9.0" or tool.get("source")!=SOURCE or not re.fullmatch(r"[0-9a-f]{64}",str(tool.get("sha256"))): raise ValueError(f"{kind} exact probe identity mismatch")
  digest=tool["sha256"] if not digest else digest
  if tool["sha256"]!=digest: raise ValueError("all edge runs must use the same probe binary")
  args=run["expected_arguments"]
  if not isinstance(args,list) or any(not isinstance(x,str) or not x for x in args) or receipt.get("arguments_sha256")!=canonical_sha(args): raise ValueError(f"{kind} arguments are not hash-bound")
  expected_target=dict(target); expected_target.update({"run_kind":kind,"observation_sha256":observation_proof["sha256"]})
  if receipt.get("target")!=expected_target: raise ValueError(f"{kind} target/observation binding mismatch")
  for channel in ("stdout","stderr"): verify_proof(root,receipt.get(channel),f"{kind}.{channel}")
  evidence.extend([proof(root,receipt_path),observation_proof,receipt["stdout"],receipt["stderr"]])
 return {"id":"EDGE_CDN_WAF","project_id":profile["project_id"],"environment":"production","release_digest":profile["release_digest"],"result":"PASS","executed_at":profile["evaluated_at"],"expires_at":profile["expires_at"],"executor_ref":profile["executor_ref"],"tool":{"name":"Cloudflare cloudflare-go + project edge probe","version":"7.9.0","source":SOURCE,"digest":"sha256:"+digest},"target":target,"assertions":{"dns_proxy_active":True,"waf_rulesets_active":True,"tls_policy_pass":True,"origin_bypass_blocked":True,"cache_behavior_pass":True},"evidence":evidence}

def main()->int:
 p=argparse.ArgumentParser(); p.add_argument("--project-root",required=True,type=Path); p.add_argument("--profile",required=True); a=p.parse_args()
 try:
  root=a.project_root.resolve(strict=True); profile=read_json(safe(root,a.profile,"profile"),"profile"); output=safe(root,profile.get("output"),"output",exists=False); output.parent.resolve(strict=True).relative_to(root); receipt=validate(profile,root); atomic(output,(json.dumps(receipt,indent=2,sort_keys=True)+"\n").encode()); print(f"CLOUDFLARE_EDGE_ADMISSION_PASS project={receipt['project_id']} release={receipt['release_digest']}"); return 0
 except (ValueError,OSError) as error: print(f"CLOUDFLARE_EDGE_ADMISSION_FAILED: {error}",file=sys.stderr); return 1
if __name__=="__main__": sys.exit(main())
