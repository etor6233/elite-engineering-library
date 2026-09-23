from __future__ import annotations
import copy,hashlib,json,tempfile,unittest
from datetime import datetime,timedelta,timezone
from pathlib import Path
from validate_k6_load_admission import canonical_sha
from validate_provider_admission import KINDS,TOOLS,validate

class ProviderAdmissionTests(unittest.TestCase):
 def setUp(self):
  self.temp=tempfile.TemporaryDirectory(); self.root=Path(self.temp.name); (self.root/"evidence").mkdir(); now=datetime.now(timezone.utc); self.now=now; z=lambda d:d.isoformat().replace("+00:00","Z"); digest="sha256:"+"c"*64
  self.profile={"schema":"elite-provider-admission/v1","project_id":"revestex","environment":"production","release_digest":digest,"evaluated_at":z(now),"expires_at":z(now+timedelta(days=7)),"executor_ref":"providers:change/77","target":{"system_url":"https://app.example.test","jurisdiction":"AR"},"providers":[],"output":"evidence/providers-control.json"}
  self.add_provider("stripe","REQUIRED","MUTATING",["create_payment_intent","receive_payment_event"],"stripe-account-hash")
  self.add_provider("google-merchant","NOT_APPLICABLE","MUTATING",["insert_product","get_product_status"],"merchant-account-hash")
 def tearDown(self): self.temp.cleanup()
 def file(self,p,data): q=self.root/p; q.parent.mkdir(parents=True,exist_ok=True); q.write_bytes(data); return q
 def write(self,p,v): q=self.root/p; q.parent.mkdir(parents=True,exist_ok=True); q.write_text(json.dumps(v),encoding="utf-8"); return q
 def pf(self,p): data=p.read_bytes(); return {"path":p.relative_to(self.root).as_posix(),"bytes":len(data),"sha256":hashlib.sha256(data).hexdigest()}
 def add_provider(self,pid,webhook_mode,mutation_mode,operations,account):
  z=lambda d:d.isoformat().replace("+00:00","Z"); exact=TOOLS[pid]; stem=f"evidence/{pid}"; policy={"schema":"elite-provider-policy/v1","project_id":"revestex","environment":"production","release_digest":self.profile["release_digest"],"provider_id":pid,"approved_at":z(self.now),"approvers":["finance-owner","security-owner"],"adapter_id":f"elite-{pid}-adapter","source_identity":exact["source_identity"],"operations":operations,"scopes":["minimum-required"],"account_reference":account,"sandbox_reference":f"{pid}-sandbox-hash","webhook_mode":webhook_mode,"mutation_mode":mutation_mode,"terms_url":"https://example.test/provider-terms","terms_version":"2026-08","terms_reviewed_at":z(self.now-timedelta(days=1)),"data_regions":["approved-region"],"retention_days":365,"quota_cost_approved":True,"finance_owner":"finance-owner","security_owner":"security-owner","reconciliation_owner":"operations-owner"}; self.write(f"{stem}-policy.json",policy); runs=[]
  base={"schema":"elite-provider-observation/v1","project_id":"revestex","environment":"production","release_digest":self.profile["release_digest"],"captured_at":z(self.now),"provider_id":pid,"account_reference":account,"operations":operations}
  observations={
   "SANDBOX_CONTRACT":dict(base,kind="SANDBOX_CONTRACT",sandbox_reference=policy["sandbox_reference"],cases=[{"operation":op,"contract_pass":True,"provider_response_observed":True,"no_unapproved_production_effect":True,"response_schema_sha256":"d"*64} for op in operations]),
   "WEBHOOK_AUTH":dict(base,kind="WEBHOOK_AUTH",mode=webhook_mode,positive_signature_verified=webhook_mode=="REQUIRED",altered_payload_rejected=webhook_mode=="REQUIRED",stale_or_replayed_rejected=webhook_mode=="REQUIRED",duplicate_converged=webhook_mode=="REQUIRED",durable_inbox=webhook_mode=="REQUIRED",missed_event_reconciled=webhook_mode=="REQUIRED",not_applicable_reason="provider operations expose no callback" if webhook_mode=="NOT_APPLICABLE" else "",reviewers=["security-owner","operations-owner"] if webhook_mode=="NOT_APPLICABLE" else []),
   "IDEMPOTENCY_RECONCILIATION":dict(base,kind="IDEMPOTENCY_RECONCILIATION",mutation_mode=mutation_mode,idempotency_key_or_dedup_used=True,duplicate_same_remote_identity=True,remote_effect_count=1,uncertain_outcome_reconciled=True,local_state_matches_provider=True,recovery_or_compensation_pass=True),
   "RATE_LIMIT_RETRY":dict(base,kind="RATE_LIMIT_RETRY",throttle_observed=True,retry_after_or_backoff_honored=True,jitter_or_provider_sdk_retry=True,max_attempts=4,attempts_observed=3,permanent_error_not_retried=True,exhaustion_durable=True,within_approved_quota=True,retry_storm_absent=True)}
  for kind in KINDS:
   key=kind.lower(); obs=self.write(f"{stem}-{key}.json",observations[kind]); args=["provider-check",pid,key]; out=self.file(f"{stem}-{key}.stdout",b"pass\n"); err=self.file(f"{stem}-{key}.stderr",b""); target=dict(self.profile["target"],provider_id=pid,account_reference=account,run_kind=kind,observation_sha256=self.pf(obs)["sha256"]); receipt={"schema":"elite-official-tool-execution/v1","project_id":"revestex","environment":"production","release_digest":self.profile["release_digest"],"control_id":"PROVIDERS","executed_at":z(self.now),"exit_code":0,"tool":{"name":exact["name"],"version":exact["version"],"source":exact["source"],"sha256":("a" if pid=="stripe" else "b")*64},"arguments_sha256":canonical_sha(args),"working_directory":".","target":target,"environment_variable_names":["PROVIDER_CREDENTIAL_REF"],"stdout":self.pf(out),"stderr":self.pf(err)}; self.write(f"{stem}-{key}-receipt.json",receipt); runs.append({"kind":kind,"execution_receipt":f"{stem}-{key}-receipt.json","observation":f"{stem}-{key}.json","expected_arguments":args})
  self.profile["providers"].append({"provider_id":pid,"policy":f"{stem}-policy.json","runs":runs})
 def observation(self,pindex,rindex): return self.root/self.profile["providers"][pindex]["runs"][rindex]["observation"]
 def receipt(self,pindex,rindex): return self.root/self.profile["providers"][pindex]["runs"][rindex]["execution_receipt"]
 def mutate(self,path,change): v=json.loads(path.read_text()); change(v); self.write(path.relative_to(self.root).as_posix(),v)
 def test_complete_selected_provider_matrix_emits_semantic_receipt(self): self.assertTrue(all(validate(self.profile,self.root)["assertions"].values()))
 def test_missing_selected_provider_is_rejected(self):
  self.profile["providers"]=[]
  with self.assertRaisesRegex(ValueError,"at least one"): validate(self.profile,self.root)
 def test_exact_official_tool_identity_is_required(self):
  self.mutate(self.receipt(0,0),lambda v:v["tool"].update(version="latest"))
  with self.assertRaisesRegex(ValueError,"exact official tool"): validate(self.profile,self.root)
 def test_every_operation_needs_sandbox_case(self):
  self.mutate(self.observation(0,0),lambda v:v.update(cases=v["cases"][:1]))
  with self.assertRaisesRegex(ValueError,"every approved"): validate(self.profile,self.root)
 def test_sandbox_cannot_hide_production_effect(self):
  self.mutate(self.observation(0,0),lambda v:v["cases"][0].update(no_unapproved_production_effect=False))
  with self.assertRaisesRegex(ValueError,"sandbox contract case failed"): validate(self.profile,self.root)
 def test_webhook_tamper_acceptance_is_rejected(self):
  self.mutate(self.observation(0,1),lambda v:v.update(altered_payload_rejected=False))
  with self.assertRaisesRegex(ValueError,"webhook authenticity"): validate(self.profile,self.root)
 def test_webhook_not_applicable_needs_two_reviewers(self):
  self.mutate(self.observation(1,1),lambda v:v.update(reviewers=["security-owner"]))
  with self.assertRaisesRegex(ValueError,"at least 2"): validate(self.profile,self.root)
 def test_duplicate_remote_effect_is_rejected(self):
  self.mutate(self.observation(0,2),lambda v:v.update(remote_effect_count=2))
  with self.assertRaisesRegex(ValueError,"idempotency/reconciliation"): validate(self.profile,self.root)
 def test_uncertain_outcome_must_reconcile(self):
  self.mutate(self.observation(0,2),lambda v:v.update(uncertain_outcome_reconciled=False))
  with self.assertRaisesRegex(ValueError,"idempotency/reconciliation"): validate(self.profile,self.root)
 def test_unbounded_retry_is_rejected(self):
  self.mutate(self.observation(0,3),lambda v:v.update(max_attempts=99,attempts_observed=99))
  with self.assertRaisesRegex(ValueError,"bounded provider"): validate(self.profile,self.root)
 def test_permanent_error_retry_is_rejected(self):
  self.mutate(self.observation(0,3),lambda v:v.update(permanent_error_not_retried=False))
  with self.assertRaisesRegex(ValueError,"bounded provider"): validate(self.profile,self.root)
 def test_stale_terms_review_is_rejected(self):
  p=self.root/self.profile["providers"][0]["policy"]; self.mutate(p,lambda v:v.update(terms_reviewed_at=(self.now-timedelta(days=91)).isoformat().replace("+00:00","Z")))
  with self.assertRaisesRegex(ValueError,"no older than 90"): validate(self.profile,self.root)
 def test_nonzero_execution_is_rejected(self):
  self.mutate(self.receipt(0,0),lambda v:v.update(exit_code=3))
  with self.assertRaisesRegex(ValueError,"identity/exit"): validate(self.profile,self.root)
 def test_secret_bearing_environment_name_is_rejected(self):
  self.mutate(self.receipt(0,0),lambda v:v.update(environment_variable_names=["PROVIDER_ACCESS_TOKEN"]))
  with self.assertRaisesRegex(ValueError,"references, not secret"): validate(self.profile,self.root)

if __name__=="__main__": unittest.main()
