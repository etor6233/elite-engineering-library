from __future__ import annotations
import hashlib,json,tempfile,unittest
from datetime import datetime,timedelta,timezone
from pathlib import Path
from validate_business_acceptance_admission import PLAYWRIGHT_SOURCE,SPEC_KIT_COMMIT,SPEC_KIT_SOURCE,validate
from validate_k6_load_admission import canonical_sha

class BusinessAcceptanceAdmissionTests(unittest.TestCase):
 def setUp(self):
  self.temp=tempfile.TemporaryDirectory(); self.root=Path(self.temp.name); (self.root/"evidence").mkdir(); now=datetime.now(timezone.utc); self.now=now; z=lambda d:d.isoformat().replace("+00:00","Z"); digest="sha256:"+"e"*64; self.target={"application_url":"https://app.example.test","business_scope_id":"launch-112","tenant_ids":["north","south"],"role_ids":["admin","seller"]}
  self.profile={"schema":"elite-business-acceptance-admission/v1","project_id":"revestex","environment":"production","release_digest":digest,"evaluated_at":z(now),"expires_at":z(now+timedelta(days=7)),"executor_ref":"acceptance:change/113","tool":{"name":"Microsoft Playwright","version":"1.62.1","source":PLAYWRIGHT_SOURCE},"target":self.target,"specification":"evidence/specification.json","playwright_run":{},"approvals":"evidence/approvals.json","defects":"evidence/defects.json","owner_signoff":"evidence/owner.json","output":"evidence/business-control.json"}
  artifacts=[]
  for kind in ("constitution","spec","plan","tasks","convergence"):
   p=self.file(f"evidence/{kind}.md",f"# {kind}\n".encode()); artifacts.append({"kind":kind,"evidence":self.pf(p)})
  scenarios=[{"scenario_id":"sale-order","capability_id":"commerce","actor_role":"seller","tenant_id":"north","preconditions":["catalog active"],"expected_outcomes":["order committed","receipt visible"]},{"scenario_id":"admin-audit","capability_id":"administration","actor_role":"admin","tenant_id":"south","preconditions":["admin authenticated"],"expected_outcomes":["audit visible"]}]
  spec={"schema":"elite-business-acceptance-specification/v1","project_id":"revestex","environment":"production","release_digest":digest,"approved_at":z(now),"spec_kit":{"name":"GitHub Spec Kit","version":"1.0.1","source":SPEC_KIT_SOURCE,"commit":SPEC_KIT_COMMIT},"artifacts":artifacts,"capabilities":["commerce","administration"],"scenarios":scenarios}; sp=self.write(self.profile["specification"],spec); spec_hash=self.pf(sp)["sha256"]
  results=[]
  for item in scenarios:
   ep=self.file(f"evidence/{item['scenario_id']}.json",b'{"verified":true}\n'); results.append({"scenario_id":item["scenario_id"],"status":"PASS","expected_outcomes_observed":True,"backend_effect_verified":True,"financial_effect":"VERIFIED" if item["capability_id"]=="commerce" else "NOT_APPLICABLE","evidence":self.pf(ep)})
  obs={"schema":"elite-business-acceptance-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"application_url":self.target["application_url"],"business_scope_id":"launch-112","specification_sha256":spec_hash,"results":results}; op=self.write("evidence/observation.json",obs); obs_hash=self.pf(op)["sha256"]
  args=["test","acceptance","--project=chromium-desktop"]; out=self.file("evidence/playwright.stdout",b"2 passed\n"); err=self.file("evidence/playwright.stderr",b""); receipt={"schema":"elite-official-tool-execution/v1","project_id":"revestex","environment":"production","release_digest":digest,"control_id":"BUSINESS_ACCEPTANCE","executed_at":z(now),"exit_code":0,"tool":{"name":"Microsoft Playwright","version":"1.62.1","source":PLAYWRIGHT_SOURCE,"sha256":"a"*64},"arguments_sha256":canonical_sha(args),"working_directory":".","target":dict(self.target,specification_sha256=spec_hash,observation_sha256=obs_hash),"environment_variable_names":["TEST_USER_TOKEN_REF"],"stdout":self.pf(out),"stderr":self.pf(err)}; self.write("evidence/playwright-receipt.json",receipt); self.profile["playwright_run"]={"execution_receipt":"evidence/playwright-receipt.json","observation":"evidence/observation.json","expected_arguments":args}
  approvals=[]
  for idx,role in enumerate(("finance","operations","security"),1):
   ep=self.file(f"evidence/{role}-approval.txt",f"approved {role}\n".encode()); approvals.append({"role":role,"subject_sha256":"sha256:"+str(idx)*64,"decision":"APPROVED","approved_at":z(now),"evidence":self.pf(ep)})
  ap=self.write(self.profile["approvals"],{"schema":"elite-business-functional-approvals/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"specification_sha256":spec_hash,"observation_sha256":obs_hash,"approvals":approvals})
  dp=self.write(self.profile["defects"],{"schema":"elite-business-defect-register/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"specification_sha256":spec_hash,"defects":[{"defect_id":"UI-LOW-1","severity":4,"status":"ACCEPTED","reason":"cosmetic only","reviewers":["product-owner","ux-owner"],"expires_at":z(now+timedelta(days=5))}]})
  ep=self.file("evidence/owner-approval.txt",b"approved owner\n"); self.write(self.profile["owner_signoff"],{"schema":"elite-business-owner-signoff/v1","project_id":"revestex","environment":"production","release_digest":digest,"approved_at":z(now),"role":"business_owner","subject_sha256":"sha256:"+"4"*64,"decision":"APPROVED","specification_sha256":spec_hash,"observation_sha256":obs_hash,"approvals_sha256":self.pf(ap)["sha256"],"defects_sha256":self.pf(dp)["sha256"],"evidence":self.pf(ep)})
 def tearDown(self): self.temp.cleanup()
 def file(self,p,data): q=self.root/p; q.write_bytes(data); return q
 def write(self,p,v): q=self.root/p; q.write_text(json.dumps(v),encoding="utf-8"); return q
 def pf(self,p): data=p.read_bytes(); return {"path":p.relative_to(self.root).as_posix(),"bytes":len(data),"sha256":hashlib.sha256(data).hexdigest()}
 def mutate(self,path,change): p=self.root/path; v=json.loads(p.read_text()); change(v); self.write(path,v)
 def test_complete_business_evidence_emits_semantic_receipt(self): self.assertTrue(all(validate(self.profile,self.root)["assertions"].values()))
 def test_missing_capability_scenario_is_rejected(self):
  self.mutate(self.profile["specification"],lambda v:v.update(scenarios=v["scenarios"][:1]))
  with self.assertRaisesRegex(ValueError,"every capability"): validate(self.profile,self.root)
 def test_failed_scenario_is_rejected(self):
  self.mutate("evidence/observation.json",lambda v:v["results"][0].update(status="FAIL"))
  with self.assertRaisesRegex(ValueError,"did not prove"): validate(self.profile,self.root)
 def test_backend_effect_is_required(self):
  self.mutate("evidence/observation.json",lambda v:v["results"][0].update(backend_effect_verified=False))
  with self.assertRaisesRegex(ValueError,"did not prove"): validate(self.profile,self.root)
 def test_observation_tamper_breaks_receipt_binding(self):
  self.mutate("evidence/observation.json",lambda v:v["results"][0].update(financial_effect="NOT_APPLICABLE"))
  with self.assertRaisesRegex(ValueError,"target/observation binding"): validate(self.profile,self.root)
 def test_nonzero_playwright_exit_is_rejected(self):
  self.mutate("evidence/playwright-receipt.json",lambda v:v.update(exit_code=1))
  with self.assertRaisesRegex(ValueError,"identity/exit"): validate(self.profile,self.root)
 def test_wrong_playwright_version_is_rejected(self):
  self.mutate("evidence/playwright-receipt.json",lambda v:v["tool"].update(version="latest"))
  with self.assertRaisesRegex(ValueError,"exact Microsoft Playwright execution"): validate(self.profile,self.root)
 def test_secret_environment_name_is_rejected(self):
  self.mutate("evidence/playwright-receipt.json",lambda v:v.update(environment_variable_names=["TEST_USER_TOKEN"]))
  with self.assertRaisesRegex(ValueError,"references, not secrets"): validate(self.profile,self.root)
 def test_missing_security_approval_is_rejected(self):
  self.mutate(self.profile["approvals"],lambda v:v.update(approvals=v["approvals"][:2]))
  with self.assertRaisesRegex(ValueError,"finance, operations and security"): validate(self.profile,self.root)
 def test_same_functional_approver_is_rejected(self):
  self.mutate(self.profile["approvals"],lambda v:v["approvals"][1].update(subject_sha256=v["approvals"][0]["subject_sha256"]))
  with self.assertRaisesRegex(ValueError,"must be distinct"): validate(self.profile,self.root)
 def test_open_sev1_is_rejected(self):
  self.mutate(self.profile["defects"],lambda v:v["defects"].append({"defect_id":"BLOCKER","severity":1,"status":"ACCEPTED","reason":"later","reviewers":["a","b"],"expires_at":(self.now+timedelta(days=1)).isoformat().replace("+00:00","Z")}))
  with self.assertRaisesRegex(ValueError,"Sev1/Sev2"): validate(self.profile,self.root)
 def test_unreviewed_accepted_defect_is_rejected(self):
  self.mutate(self.profile["defects"],lambda v:v["defects"][0].update(reviewers=["one"]))
  with self.assertRaisesRegex(ValueError,"at least 2"): validate(self.profile,self.root)
 def test_owner_must_be_independent(self):
  self.mutate(self.profile["owner_signoff"],lambda v:v.update(subject_sha256="sha256:"+"1"*64))
  with self.assertRaisesRegex(ValueError,"independently owned"): validate(self.profile,self.root)
 def test_owner_signoff_tamper_is_rejected(self):
  self.mutate(self.profile["owner_signoff"],lambda v:v.update(defects_sha256="0"*64))
  with self.assertRaisesRegex(ValueError,"signoff evidence binding"): validate(self.profile,self.root)

if __name__=="__main__": unittest.main()
