from __future__ import annotations
import copy,hashlib,json
from datetime import datetime,timedelta,timezone
from pathlib import Path
import tempfile,unittest
from validate_k6_load_admission import canonical_sha
from validate_zap_offensive_admission import KINDS,SOURCE,validate
class ZapAdmissionTests(unittest.TestCase):
 def setUp(self):
  self.temp=tempfile.TemporaryDirectory(); self.root=Path(self.temp.name).resolve(); (self.root/"evidence").mkdir(); now=datetime.now(timezone.utc); z=lambda d:d.isoformat().replace("+00:00","Z"); digest="sha256:"+"e"*64
  self.profile={"schema":"elite-zap-offensive-admission/v1","project_id":"revestex","environment":"production","release_digest":digest,"evaluated_at":z(now),"expires_at":z(now+timedelta(days=7)),"executor_ref":"security:engagement/42","target":{"base_url":"https://app.example.test","api_base_url":"https://api.example.test","authorization_ticket":"SEC-42"},"policy":"evidence/policy.json","manual_review":"evidence/review.json","runs":[],"output":"evidence/offensive-control.json"}
  self.write("evidence/policy.json",{"schema":"elite-offensive-security-policy/v1","project_id":"revestex","environment":"production","release_digest":digest,"approved_at":z(now),"approvers":["security-owner","system-owner"],"accepted_alerts":[]})
  self.write("evidence/review.json",{"schema":"elite-offensive-manual-review/v1","project_id":"revestex","environment":"production","release_digest":digest,"reviewed_at":z(now),"reviewers":["security-reviewer","business-reviewer"],"authorization_confirmed":True,"authenticated_web_reviewed":True,"authenticated_api_reviewed":True,"business_logic_reviewed":True,"notes":"Authorized scope and residual risk reviewed."})
  for kind in KINDS:
   stem=kind.lower(); plan=self.root/"evidence"/f"{stem}.yaml"; plan.write_text("env:\n  contexts: []\n",encoding="utf-8"); args=["-cmd","-autorun",f"evidence/{stem}.yaml"]
   out=self.root/"evidence"/f"{stem}.stdout.bin"; err=self.root/"evidence"/f"{stem}.stderr.bin"; out.write_bytes(b"Automation plan succeeded\n"); err.write_bytes(b""); pf=lambda p:{"path":p.relative_to(self.root).as_posix(),"bytes":p.stat().st_size,"sha256":hashlib.sha256(p.read_bytes()).hexdigest()}; target=dict(self.profile["target"]); target.update({"run_kind":kind,"automation_plan_sha256":pf(plan)["sha256"]})
   receipt={"schema":"elite-official-tool-execution/v1","project_id":"revestex","environment":"production","release_digest":digest,"control_id":"OFFENSIVE_SECURITY","executed_at":z(now),"exit_code":0,"tool":{"name":"OWASP ZAP","version":"2.17.0","source":SOURCE,"sha256":"f"*64},"arguments_sha256":canonical_sha(args),"working_directory":".","target":target,"environment_variable_names":[],"stdout":pf(out),"stderr":pf(err)}; self.write(f"evidence/{stem}-receipt.json",receipt); self.write(f"evidence/{stem}-report.json",{"@version":"2.17.0","site":[{"@name":self.profile["target"]["base_url" if kind=="WEB_AUTHENTICATED" else "api_base_url"],"alerts":[{"pluginid":"10021","riskcode":"1","name":"Low fixture","instances":[]}]}]}); self.profile["runs"].append({"kind":kind,"execution_receipt":f"evidence/{stem}-receipt.json","plan":f"evidence/{stem}.yaml","report":f"evidence/{stem}-report.json","expected_arguments":args})
 def tearDown(self): self.temp.cleanup()
 def write(self,p,v): (self.root/p).write_text(json.dumps(v),encoding="utf-8")
 def test_web_api_and_manual_review_emit_semantic_receipt(self): self.assertTrue(all(validate(self.profile,self.root)["assertions"].values()))
 def test_unauthorized_target_is_rejected(self):
  v=copy.deepcopy(self.profile); v["target"]["authorization_ticket"]=""
  with self.assertRaisesRegex(ValueError,"authorization ticket"): validate(v,self.root)
 def test_high_alert_is_rejected(self):
  p=self.root/self.profile["runs"][0]["report"]; v=json.loads(p.read_text()); v["site"][0]["alerts"][0]["riskcode"]="3"; self.write(self.profile["runs"][0]["report"],v)
  with self.assertRaisesRegex(ValueError,"unaccepted high"): validate(self.profile,self.root)
 def test_expired_exception_is_rejected(self):
  p=self.root/self.profile["policy"]; v=json.loads(p.read_text()); v["accepted_alerts"]=[{"pluginid":"42","riskcode":3,"reason":"temporary","owner":"security","expires_at":"2020-01-01T00:00:00Z"}]; self.write(self.profile["policy"],v)
  with self.assertRaisesRegex(ValueError,"expires before"): validate(self.profile,self.root)
 def test_plan_tamper_is_rejected(self):
  (self.root/self.profile["runs"][0]["plan"]).write_text("tampered")
  with self.assertRaisesRegex(ValueError,"binding mismatch"): validate(self.profile,self.root)
 def test_missing_api_run_is_rejected(self):
  v=copy.deepcopy(self.profile); v["runs"]=v["runs"][:1]
  with self.assertRaisesRegex(ValueError,"web and API"): validate(v,self.root)
 def test_manual_review_false_is_rejected(self):
  p=self.root/self.profile["manual_review"]; v=json.loads(p.read_text()); v["business_logic_reviewed"]=False; self.write(self.profile["manual_review"],v)
  with self.assertRaisesRegex(ValueError,"manual review assertions"): validate(self.profile,self.root)
 def test_nonzero_zap_exit_is_rejected(self):
  p=self.root/self.profile["runs"][1]["execution_receipt"]; v=json.loads(p.read_text()); v["exit_code"]=2; self.write(self.profile["runs"][1]["execution_receipt"],v)
  with self.assertRaisesRegex(ValueError,"identity/exit"): validate(self.profile,self.root)
if __name__=="__main__": unittest.main()
