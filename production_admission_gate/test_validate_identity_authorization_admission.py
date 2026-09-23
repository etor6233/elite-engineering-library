from __future__ import annotations
import copy,hashlib,json,tempfile,unittest
from datetime import datetime,timedelta,timezone
from itertools import permutations
from pathlib import Path
from validate_identity_authorization_admission import KINDS,OIDF_SOURCE,PLAYWRIGHT_SOURCE,validate
from validate_k6_load_admission import canonical_sha

class IdentityAuthorizationAdmissionTests(unittest.TestCase):
 def setUp(self):
  self.temp=tempfile.TemporaryDirectory(); self.root=Path(self.temp.name); (self.root/"evidence").mkdir(); now=datetime.now(timezone.utc); z=lambda d:d.isoformat().replace("+00:00","Z"); digest="sha256:"+"e"*64
  self.target={"issuer_url":"https://id.example.test/tenant","client_id":"web-client","audience":"revestex-api","conformance_profile":"oidcc-basic-certification-test-plan","application_url":"https://app.example.test","tenant_ids":["north","south"],"role_ids":["admin","seller"]}
  self.profile={"schema":"elite-identity-authorization-admission/v1","project_id":"revestex","environment":"production","release_digest":digest,"evaluated_at":z(now),"expires_at":z(now+timedelta(days=7)),"executor_ref":"identity:change/42","tool":{"name":"OpenID Foundation Conformance Suite","version":"5.2.4","source":OIDF_SOURCE},"target":self.target,"policy":"evidence/policy.json","runs":[],"output":"evidence/identity-control.json"}
  self.write("evidence/policy.json",{"schema":"elite-identity-authorization-policy/v1","project_id":"revestex","environment":"production","release_digest":digest,"approved_at":z(now),"approvers":["security-owner","business-owner"],"issuer_url":self.target["issuer_url"],"client_id":"web-client","audience":"revestex-api","conformance_profile":self.target["conformance_profile"],"tenant_ids":["north","south"],"role_ids":["admin","seller"],"break_glass":{"configured":True,"max_minutes":30,"two_person_approval":True,"audit_required":True,"post_use_review_required":True}})
  observations={
   "OIDC_CONFORMANCE":{"schema":"elite-identity-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"kind":"OIDC_CONFORMANCE","issuer_url":self.target["issuer_url"],"client_id":"web-client","conformance_profile":self.target["conformance_profile"],"suite_release":"release-v5.2.4","plan_id":"plan-42","summary":{"passed":3,"review":0,"warning":0,"skipped":1,"failed":0,"interrupted":0,"total":4},"non_passed":[{"test_id":"optional-encryption","status":"SKIPPED","reason":"feature excluded from approved profile","reviewers":["security-owner","identity-owner"]}]},
   "ROLE_POSITIVE_JOURNEYS":{"schema":"elite-identity-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"kind":"ROLE_POSITIVE_JOURNEYS","application_url":self.target["application_url"],"journeys":[{"role_id":role,"tenant_id":"north","journey_id":f"{role}-write","authenticated":True,"status":200,"expected_backend_effect":True,"effect_observed":True} for role in self.target["role_ids"]]},
   "UNAUTHORIZED_NEGATIVE_JOURNEYS":{"schema":"elite-identity-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"kind":"UNAUTHORIZED_NEGATIVE_JOURNEYS","application_url":self.target["application_url"],"scenarios":[{"scenario":case,"status":401 if case in ("unauthenticated","expired_token","wrong_issuer","wrong_audience") else 403,"data_disclosed":False,"mutation_observed":False} for case in ("unauthenticated","expired_token","wrong_issuer","wrong_audience","insufficient_role","object_ownership")]},
   "TENANT_ISOLATION":{"schema":"elite-identity-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"kind":"TENANT_ISOLATION","tenant_ids":["north","south"],"scenarios":[{"actor_tenant":a,"target_tenant":b,"action":action,"resource_key":"shared-id-42","status":403,"data_disclosed":False,"mutation_observed":False} for a,b in permutations(self.target["tenant_ids"],2) for action in ("read","write","list")]},
   "SESSION_REVOCATION_ROTATION":{"schema":"elite-identity-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"kind":"SESSION_REVOCATION_ROTATION","issuer_url":self.target["issuer_url"],"cases":[{"case":case,"status":401,"old_credential_accepted":False,"rejected":True,"audit_event_observed":True} for case in ("logout_revocation","token_revocation","signing_key_rotation","session_expiry","replay_rejection")]}}
  for kind in KINDS:
   stem=kind.lower(); obs=self.write(f"evidence/{stem}.json",observations[kind]); args=["--plan","plan-42"] if kind=="OIDC_CONFORMANCE" else ["test",stem]; out=self.file(f"evidence/{stem}.stdout",b"pass\n"); err=self.file(f"evidence/{stem}.stderr",b""); tool={"name":"OpenID Foundation Conformance Runner","version":"5.2.4","source":OIDF_SOURCE,"sha256":"a"*64} if kind=="OIDC_CONFORMANCE" else {"name":"Microsoft Playwright","version":"1.62.1","source":PLAYWRIGHT_SOURCE,"sha256":"b"*64}; target=dict(self.target); target.update({"run_kind":kind,"observation_sha256":self.pf(obs)["sha256"]})
   receipt={"schema":"elite-official-tool-execution/v1","project_id":"revestex","environment":"production","release_digest":digest,"control_id":"IDENTITY_AUTHORIZATION","executed_at":z(now),"exit_code":0,"tool":tool,"arguments_sha256":canonical_sha(args),"working_directory":".","target":target,"environment_variable_names":["OIDC_CLIENT_SECRET_REF"],"stdout":self.pf(out),"stderr":self.pf(err)}; self.write(f"evidence/{stem}-receipt.json",receipt); self.profile["runs"].append({"kind":kind,"execution_receipt":f"evidence/{stem}-receipt.json","observation":f"evidence/{stem}.json","expected_arguments":args})
 def tearDown(self): self.temp.cleanup()
 def file(self,p,data): q=self.root/p; q.write_bytes(data); return q
 def write(self,p,v): q=self.root/p; q.write_text(json.dumps(v),encoding="utf-8"); return q
 def pf(self,p): data=p.read_bytes(); return {"path":p.relative_to(self.root).as_posix(),"bytes":len(data),"sha256":hashlib.sha256(data).hexdigest()}
 def mutate(self,index,change): p=self.root/self.profile["runs"][index]["observation"]; v=json.loads(p.read_text()); change(v); self.write(self.profile["runs"][index]["observation"],v)
 def test_complete_identity_evidence_emits_semantic_receipt(self): self.assertTrue(all(validate(self.profile,self.root)["assertions"].values()))
 def test_failed_oidc_case_is_rejected(self):
  self.mutate(0,lambda v:v["summary"].update(passed=2,failed=1))
  with self.assertRaisesRegex(ValueError,"failed/interrupted"): validate(self.profile,self.root)
 def test_unreviewed_oidc_warning_is_rejected(self):
  self.mutate(0,lambda v:(v["summary"].update(passed=2,warning=1),v.update(non_passed=v["non_passed"])))
  with self.assertRaisesRegex(ValueError,"every non-passed"): validate(self.profile,self.root)
 def test_missing_role_journey_is_rejected(self):
  self.mutate(1,lambda v:v.update(journeys=v["journeys"][:1]))
  with self.assertRaisesRegex(ValueError,"every declared role"): validate(self.profile,self.root)
 def test_unauthorized_mutation_is_rejected(self):
  self.mutate(2,lambda v:v["scenarios"][0].update(mutation_observed=True))
  with self.assertRaisesRegex(ValueError,"changed state"): validate(self.profile,self.root)
 def test_incomplete_tenant_matrix_is_rejected(self):
  self.mutate(3,lambda v:v.update(scenarios=v["scenarios"][:-1]))
  with self.assertRaisesRegex(ValueError,"complete ordered"): validate(self.profile,self.root)
 def test_cross_tenant_disclosure_is_rejected(self):
  self.mutate(3,lambda v:v["scenarios"][0].update(data_disclosed=True))
  with self.assertRaisesRegex(ValueError,"isolation failed"): validate(self.profile,self.root)
 def test_old_session_acceptance_is_rejected(self):
  self.mutate(4,lambda v:v["cases"][0].update(old_credential_accepted=True,rejected=False,status=200))
  with self.assertRaisesRegex(ValueError,"assertion failed"): validate(self.profile,self.root)
 def test_nonzero_official_tool_exit_is_rejected(self):
  p=self.root/self.profile["runs"][0]["execution_receipt"]; v=json.loads(p.read_text()); v["exit_code"]=2; self.write(self.profile["runs"][0]["execution_receipt"],v)
  with self.assertRaisesRegex(ValueError,"identity/exit"): validate(self.profile,self.root)
 def test_observation_tamper_is_rejected_by_receipt_binding(self):
  self.mutate(1,lambda v:v["journeys"][0].update(journey_id="renamed"))
  with self.assertRaisesRegex(ValueError,"target/observation binding"): validate(self.profile,self.root)
 def test_secret_bearing_environment_name_is_rejected(self):
  p=self.root/self.profile["runs"][1]["execution_receipt"]; v=json.loads(p.read_text()); v["environment_variable_names"]=["OIDC_CLIENT_SECRET"]; self.write(self.profile["runs"][1]["execution_receipt"],v)
  with self.assertRaisesRegex(ValueError,"references, not secret"): validate(self.profile,self.root)
 def test_incomplete_break_glass_is_rejected(self):
  p=self.root/self.profile["policy"]; v=json.loads(p.read_text()); v["break_glass"]["two_person_approval"]=False; self.write(self.profile["policy"],v)
  with self.assertRaisesRegex(ValueError,"break-glass"): validate(self.profile,self.root)

if __name__=="__main__": unittest.main()
