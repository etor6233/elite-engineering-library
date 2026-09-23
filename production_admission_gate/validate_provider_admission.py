from __future__ import annotations
import argparse,json,re,sys
from pathlib import Path
from urllib.parse import urlparse
from validate_k6_load_admission import atomic,canonical_sha,proof,read_json,safe,utc,verify_proof

KINDS=("SANDBOX_CONTRACT","WEBHOOK_AUTH","IDEMPOTENCY_RECONCILIATION","RATE_LIMIT_RETRY")
FIELDS={"schema","project_id","environment","release_digest","evaluated_at","expires_at","executor_ref","target","providers","output"}
TARGET_FIELDS={"system_url","jurisdiction"}
PROVIDER_FIELDS={"provider_id","policy","runs"}
RUN_FIELDS={"kind","execution_receipt","observation","expected_arguments"}
POLICY_FIELDS={"schema","project_id","environment","release_digest","provider_id","approved_at","approvers","adapter_id","source_identity","operations","scopes","account_reference","sandbox_reference","webhook_mode","mutation_mode","terms_url","terms_version","terms_reviewed_at","data_regions","retention_days","quota_cost_approved","finance_owner","security_owner","reconciliation_owner"}
TOOLS={
 "stripe":{"name":"Stripe Go SDK","version":"86.3.0","source":"https://github.com/stripe/stripe-go/tree/a2df585a800a97fe8ec4ebf551b4449bdb3d90a1","source_identity":"stripe-go-86.3.0"},
 "mercadopago":{"name":"Mercado Pago Go SDK","version":"1.14.0","source":"https://github.com/mercadopago/sdk-go/tree/f910ee53fbb6819e435eaf3d0f800cb1fe74ae09","source_identity":"mercadopago-sdk-go-1.14.0"},
 "amazon-spapi":{"name":"Amazon Selling Partner API Python SDK","version":"1.11.1","source":"https://github.com/amzn/selling-partner-api-sdk/tree/8e792ae345a8d334ccdbdd03181f05f040e6a4fc","source_identity":"amazon-selling-partner-api-sdk-python-1.11.1"},
 "google-ads":{"name":"Google Ads Python SDK","version":"31.3.0","source":"https://github.com/googleads/google-ads-python/tree/f7bf312d26904ea50ad9ef02e395edccfaa1ebb7","source_identity":"google-ads-python-31.3.0"},
 "meta-ads":{"name":"Meta Business Python SDK","version":"26.0.1","source":"https://github.com/facebook/facebook-python-business-sdk/tree/788f363d15b1269ab5efb7cd00fb5e3b133cd99b","source_identity":"meta-business-sdk-python-26.0.1"},
 "tiktok-ads":{"name":"TikTok Business API Python SDK","version":"f809c396520df2d7b201a9ccc5378d822b728ed3","source":"https://github.com/tiktok/tiktok-business-api-sdk/tree/f809c396520df2d7b201a9ccc5378d822b728ed3","source_identity":"tiktok-business-api-sdk"},
 "meta-whatsapp":{"name":"Meta WhatsApp Cloud API examples","version":"de70ee908a67026e642aaee3703d20464e2a9466","source":"https://github.com/fbsamples/whatsapp-api-examples/tree/de70ee908a67026e642aaee3703d20464e2a9466","source_identity":"meta-whatsapp-api-examples"},
 "firebase-fcm":{"name":"Firebase Admin Go SDK","version":"4.21.0","source":"https://github.com/firebase/firebase-admin-go/tree/eebb06f2a643fbb59b1cb262874a943584475128","source_identity":"firebase-admin-go-4.21.0"},
 "google-merchant":{"name":"Google Merchant Products Python SDK","version":"1.8.0","source":"https://github.com/googleapis/google-cloud-python/tree/97d7b42cd74b41211f5ec8871cc0dd15debdb1a0/packages/google-shopping-merchant-products","source_identity":"google-shopping-merchant-products-python-1.8.0"},
 "mercadolibre":{"name":"Elite Mercado Libre official-HTTP-contract adapter","version":"0.2.0","source":"https://developers.mercadolibre.com.ar/en_us/categories-and-attributes/manage-questions-and-answers","source_identity":"GO-MERCADOLIBRE-MARKETPLACE-ADAPTER@0.2.0:AUTHORED"},
}

def nonempty(value:object,label:str)->str:
 if not isinstance(value,str) or not value.strip() or "replace-me" in value.lower(): raise ValueError(f"{label} is required")
 return value

def https(value:object,label:str)->None:
 value=nonempty(value,label); parsed=urlparse(value)
 if parsed.scheme!="https" or not parsed.hostname or parsed.username or parsed.password or parsed.fragment: raise ValueError(f"{label} must be credential-free HTTPS")

def strings(value:object,label:str,minimum:int=1)->list[str]:
 if not isinstance(value,list) or len(value)<minimum or len(set(value))!=len(value) or any(not isinstance(x,str) or not x.strip() for x in value): raise ValueError(f"{label} must contain at least {minimum} unique values")
 return value

def bind(value:dict[str,object],profile:dict[str,object],label:str)->None:
 for key in ("project_id","environment","release_digest"):
  if value.get(key)!=profile[key]: raise ValueError(f"{label} {key} mismatch")

def common_observation(value:object,profile:dict[str,object],policy:dict[str,object],kind:str)->dict[str,object]:
 if not isinstance(value,dict): raise ValueError(f"{kind} observation must be an object")
 required={"schema","project_id","environment","release_digest","captured_at","provider_id","kind","account_reference","operations"}
 if not required.issubset(value) or value.get("schema")!="elite-provider-observation/v1" or value.get("kind")!=kind: raise ValueError(f"{kind} observation identity mismatch")
 bind(value,profile,kind); utc(value["captured_at"],f"{kind}.captured_at")
 if value.get("provider_id")!=policy["provider_id"] or value.get("account_reference")!=policy["account_reference"] or value.get("operations")!=policy["operations"]: raise ValueError(f"{kind} provider/account/operations mismatch")
 return value

def sandbox(value:object,profile:dict[str,object],policy:dict[str,object])->None:
 value=common_observation(value,profile,policy,"SANDBOX_CONTRACT"); fields={"schema","project_id","environment","release_digest","captured_at","provider_id","kind","account_reference","operations","sandbox_reference","cases"}
 if set(value)!=fields or value["sandbox_reference"]!=policy["sandbox_reference"]: raise ValueError("sandbox observation fields/reference mismatch")
 cases=value["cases"]
 if not isinstance(cases,list) or {x.get("operation") for x in cases if isinstance(x,dict)}!=set(policy["operations"]): raise ValueError("every approved provider operation requires a sandbox case")
 for item in cases:
  if set(item)!={"operation","contract_pass","provider_response_observed","no_unapproved_production_effect","response_schema_sha256"} or item["contract_pass"] is not True or item["provider_response_observed"] is not True or item["no_unapproved_production_effect"] is not True or not re.fullmatch(r"[0-9a-f]{64}",str(item["response_schema_sha256"])): raise ValueError("sandbox contract case failed")

def webhook(value:object,profile:dict[str,object],policy:dict[str,object])->None:
 value=common_observation(value,profile,policy,"WEBHOOK_AUTH"); fields={"schema","project_id","environment","release_digest","captured_at","provider_id","kind","account_reference","operations","mode","positive_signature_verified","altered_payload_rejected","stale_or_replayed_rejected","duplicate_converged","durable_inbox","missed_event_reconciled","not_applicable_reason","reviewers"}
 if set(value)!=fields or value["mode"]!=policy["webhook_mode"]: raise ValueError("webhook observation fields/mode mismatch")
 if value["mode"]=="REQUIRED":
  if any(value[k] is not True for k in ("positive_signature_verified","altered_payload_rejected","stale_or_replayed_rejected","duplicate_converged","durable_inbox","missed_event_reconciled")) or value["not_applicable_reason"]!="" or value["reviewers"]!=[]: raise ValueError("required webhook authenticity/durability evidence failed")
 elif value["mode"]=="NOT_APPLICABLE":
  if any(value[k] is not False for k in ("positive_signature_verified","altered_payload_rejected","stale_or_replayed_rejected","duplicate_converged","durable_inbox","missed_event_reconciled")) or not nonempty(value["not_applicable_reason"],"webhook not_applicable_reason") or len(strings(value["reviewers"],"webhook reviewers",2))<2: raise ValueError("webhook N/A requires reason and two reviewers")
 else: raise ValueError("webhook mode must be REQUIRED or NOT_APPLICABLE")

def idempotency(value:object,profile:dict[str,object],policy:dict[str,object])->None:
 value=common_observation(value,profile,policy,"IDEMPOTENCY_RECONCILIATION"); fields={"schema","project_id","environment","release_digest","captured_at","provider_id","kind","account_reference","operations","mutation_mode","idempotency_key_or_dedup_used","duplicate_same_remote_identity","remote_effect_count","uncertain_outcome_reconciled","local_state_matches_provider","recovery_or_compensation_pass"}
 if set(value)!=fields or value["mutation_mode"]!=policy["mutation_mode"]: raise ValueError("idempotency observation fields/mode mismatch")
 if value["idempotency_key_or_dedup_used"] is not True or value["duplicate_same_remote_identity"] is not True or value["remote_effect_count"]!=1 or value["uncertain_outcome_reconciled"] is not True or value["local_state_matches_provider"] is not True or value["recovery_or_compensation_pass"] is not True: raise ValueError("idempotency/reconciliation evidence failed")

def rate_limit(value:object,profile:dict[str,object],policy:dict[str,object])->None:
 value=common_observation(value,profile,policy,"RATE_LIMIT_RETRY"); fields={"schema","project_id","environment","release_digest","captured_at","provider_id","kind","account_reference","operations","throttle_observed","retry_after_or_backoff_honored","jitter_or_provider_sdk_retry","max_attempts","attempts_observed","permanent_error_not_retried","exhaustion_durable","within_approved_quota","retry_storm_absent"}
 if set(value)!=fields or value["throttle_observed"] is not True or value["retry_after_or_backoff_honored"] is not True or value["jitter_or_provider_sdk_retry"] is not True or not isinstance(value["max_attempts"],int) or not 2<=value["max_attempts"]<=10 or not isinstance(value["attempts_observed"],int) or not 2<=value["attempts_observed"]<=value["max_attempts"] or any(value[k] is not True for k in ("permanent_error_not_retried","exhaustion_durable","within_approved_quota","retry_storm_absent")): raise ValueError("bounded provider rate-limit/retry evidence failed")

def validate(profile:object,root:Path)->dict[str,object]:
 if not isinstance(profile,dict) or set(profile)!=FIELDS: raise ValueError("profile fields must be exact")
 if profile["schema"]!="elite-provider-admission/v1" or profile["environment"]!="production" or not nonempty(profile["project_id"],"project_id") or not re.fullmatch(r"sha256:[0-9a-f]{64}",str(profile["release_digest"])): raise ValueError("profile identity is invalid")
 start=utc(profile["evaluated_at"],"evaluated_at"); end=utc(profile["expires_at"],"expires_at")
 if not start<end or (end-start).total_seconds()>30*86400: raise ValueError("admission validity must be positive and at most 30 days")
 nonempty(profile["executor_ref"],"executor_ref"); target=profile["target"]
 if not isinstance(target,dict) or set(target)!=TARGET_FIELDS: raise ValueError("target fields must be exact")
 https(target["system_url"],"target.system_url"); nonempty(target["jurisdiction"],"target.jurisdiction")
 providers=profile["providers"]
 if not isinstance(providers,list) or not providers or any(not isinstance(x,dict) or set(x)!=PROVIDER_FIELDS for x in providers): raise ValueError("at least one exact provider entry is required")
 ids=[x["provider_id"] for x in providers]
 if len(ids)!=len(set(ids)) or any(x not in TOOLS for x in ids): raise ValueError("provider IDs must be unique and admitted")
 evidence=[]; tool_hashes=[]; summaries=[]
 validators={"SANDBOX_CONTRACT":sandbox,"WEBHOOK_AUTH":webhook,"IDEMPOTENCY_RECONCILIATION":idempotency,"RATE_LIMIT_RETRY":rate_limit}
 for entry in providers:
  provider_id=entry["provider_id"]; exact=TOOLS[provider_id]; policy_path=safe(root,entry["policy"],f"{provider_id}.policy"); policy=read_json(policy_path,f"{provider_id}.policy")
  if not isinstance(policy,dict) or set(policy)!=POLICY_FIELDS or policy.get("schema")!="elite-provider-policy/v1" or policy.get("provider_id")!=provider_id: raise ValueError(f"{provider_id} policy fields/schema mismatch")
  bind(policy,profile,f"{provider_id}.policy"); approved=utc(policy["approved_at"],f"{provider_id}.approved_at"); reviewed=utc(policy["terms_reviewed_at"],f"{provider_id}.terms_reviewed_at")
  if reviewed>approved or (approved-reviewed).total_seconds()>90*86400: raise ValueError(f"{provider_id} terms review must be no older than 90 days at approval")
  approvers=strings(policy["approvers"],f"{provider_id}.approvers",2)
  for key in ("adapter_id","account_reference","sandbox_reference","terms_version","finance_owner","security_owner","reconciliation_owner"): nonempty(policy[key],f"{provider_id}.{key}")
  if len({policy["finance_owner"],policy["security_owner"],policy["reconciliation_owner"]})<3 or not set((policy["finance_owner"],policy["security_owner"])).issubset(set(approvers)): raise ValueError(f"{provider_id} finance/security/reconciliation ownership is not independent")
  strings(policy["operations"],f"{provider_id}.operations"); strings(policy["scopes"],f"{provider_id}.scopes"); strings(policy["data_regions"],f"{provider_id}.data_regions")
  https(policy["terms_url"],f"{provider_id}.terms_url")
  if policy["source_identity"]!=exact["source_identity"] or policy["webhook_mode"] not in ("REQUIRED","NOT_APPLICABLE") or policy["mutation_mode"] not in ("READ_ONLY","MUTATING") or policy["quota_cost_approved"] is not True or not isinstance(policy["retention_days"],int) or not 1<=policy["retention_days"]<=3650: raise ValueError(f"{provider_id} source/modes/terms/cost policy invalid")
  runs=entry["runs"]
  if not isinstance(runs,list) or len(runs)!=4 or any(not isinstance(x,dict) or set(x)!=RUN_FIELDS for x in runs) or [x["kind"] for x in runs]!=list(KINDS): raise ValueError(f"{provider_id} requires four ordered provider runs")
  evidence.append(proof(root,policy_path)); family_hash=None
  for run in runs:
   kind=run["kind"]; observation_path=safe(root,run["observation"],f"{provider_id}.{kind}.observation"); observation=read_json(observation_path,f"{provider_id}.{kind}.observation"); validators[kind](observation,profile,policy); observation_proof=proof(root,observation_path)
   receipt_path=safe(root,run["execution_receipt"],f"{provider_id}.{kind}.receipt"); receipt=read_json(receipt_path,f"{provider_id}.{kind}.receipt")
   if receipt.get("schema")!="elite-official-tool-execution/v1" or receipt.get("project_id")!=profile["project_id"] or receipt.get("environment")!="production" or receipt.get("release_digest")!=profile["release_digest"] or receipt.get("control_id")!="PROVIDERS" or receipt.get("exit_code")!=0: raise ValueError(f"{provider_id} {kind} execution identity/exit mismatch")
   tool=receipt.get("tool")
   if not isinstance(tool,dict) or set(tool)!={"name","version","source","sha256"} or any(tool.get(k)!=exact[k] for k in ("name","version","source")) or not re.fullmatch(r"[0-9a-f]{64}",str(tool.get("sha256"))): raise ValueError(f"{provider_id} exact official tool identity mismatch")
   if family_hash is None: family_hash=tool["sha256"]
   if family_hash!=tool["sha256"]: raise ValueError(f"{provider_id} runs must use one exact tool/source hash")
   args=run["expected_arguments"]
   if not isinstance(args,list) or not args or any(not isinstance(x,str) or not x for x in args) or receipt.get("arguments_sha256")!=canonical_sha(args): raise ValueError(f"{provider_id} {kind} arguments are not hash-bound")
   expected_target={"system_url":target["system_url"],"jurisdiction":target["jurisdiction"],"provider_id":provider_id,"account_reference":policy["account_reference"],"run_kind":kind,"observation_sha256":observation_proof["sha256"]}
   if receipt.get("target")!=expected_target: raise ValueError(f"{provider_id} {kind} target/observation binding mismatch")
   if not isinstance(receipt.get("environment_variable_names"),list) or any(re.search(r"(?:secret|password|token|key|credential)",str(x),re.I) and not str(x).endswith("_REF") for x in receipt["environment_variable_names"]): raise ValueError(f"{provider_id} environment must expose references, not secret-bearing names")
   for channel in ("stdout","stderr"): verify_proof(root,receipt.get(channel),f"{provider_id}.{kind}.{channel}")
   evidence.extend([proof(root,receipt_path),observation_proof,receipt["stdout"],receipt["stderr"]])
  tool_hashes.append(family_hash); summaries.append({"provider_id":provider_id,"adapter_id":policy["adapter_id"],"source_identity":policy["source_identity"],"account_reference":policy["account_reference"],"operations":policy["operations"],"webhook_mode":policy["webhook_mode"],"mutation_mode":policy["mutation_mode"]})
 return {"id":"PROVIDERS","project_id":profile["project_id"],"environment":"production","release_digest":profile["release_digest"],"result":"PASS","executed_at":profile["evaluated_at"],"expires_at":profile["expires_at"],"executor_ref":profile["executor_ref"],"tool":{"name":"Selected exact provider SDKs and official-contract adapters","version":"policy-locked","source":"OFFICIAL-UPSTREAM-ACQUISITION-CORE:commerce-communications-leaders","digest":"sha256:"+canonical_sha(tool_hashes)},"target":{"system_url":target["system_url"],"jurisdiction":target["jurisdiction"],"providers":summaries},"assertions":{"all_required_sandboxes_pass":True,"webhook_auth_pass":True,"idempotency_reconciliation_pass":True,"rate_limit_retry_pass":True,"terms_cost_scope_approved":True},"evidence":evidence}

def main()->int:
 p=argparse.ArgumentParser(); p.add_argument("--project-root",required=True,type=Path); p.add_argument("--profile",required=True); a=p.parse_args()
 try:
  root=a.project_root.resolve(strict=True); profile=read_json(safe(root,a.profile,"profile"),"profile"); output=safe(root,profile.get("output"),"output",exists=False); output.parent.resolve(strict=True).relative_to(root); receipt=validate(profile,root); atomic(output,(json.dumps(receipt,indent=2,sort_keys=True)+"\n").encode()); print(f"PROVIDER_ADMISSION_PASS project={receipt['project_id']} release={receipt['release_digest']} providers={len(receipt['target']['providers'])}"); return 0
 except (ValueError,OSError) as error: print(f"PROVIDER_ADMISSION_FAILED: {error}",file=sys.stderr); return 1
if __name__=="__main__": sys.exit(main())
