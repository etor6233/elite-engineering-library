from __future__ import annotations
import copy,hashlib,json,tempfile,unittest
from datetime import datetime,timedelta,timezone
from pathlib import Path
from validate_cloudflare_edge_admission import KINDS,SOURCE,validate
from validate_k6_load_admission import canonical_sha

class CloudflareEdgeAdmissionTests(unittest.TestCase):
 def setUp(self):
  self.temp=tempfile.TemporaryDirectory(); self.root=Path(self.temp.name); (self.root/"evidence").mkdir(); now=datetime.now(timezone.utc); z=lambda d:d.isoformat().replace("+00:00","Z"); digest="sha256:"+"e"*64
  self.target={"zone_id":"zone-123","hostname":"app.example.test","origin_address":"203.0.113.10","public_url":"https://app.example.test/assets/app.js","private_url":"https://app.example.test/api/me"}
  self.profile={"schema":"elite-cloudflare-edge-admission/v1","project_id":"revestex","environment":"production","release_digest":digest,"evaluated_at":z(now),"expires_at":z(now+timedelta(days=7)),"executor_ref":"edge:change/42","tool":{"name":"Cloudflare cloudflare-go","version":"7.9.0","source":SOURCE},"target":self.target,"policy":"evidence/policy.json","runs":[],"output":"evidence/edge-control.json"}
  self.policy={"schema":"elite-cloudflare-edge-policy/v1","project_id":"revestex","environment":"production","release_digest":digest,"approved_at":z(now),"approvers":["security-owner","platform-owner"],"zone_id":"zone-123","hostname":"app.example.test","managed_ruleset_ids":["managed-1"],"custom_ruleset_ids":["custom-1"],"minimum_tls_version":"1.2","origin_protection_mode":"authenticated_origin_pull","public_cache_path":"/assets/app.js","private_cache_path":"/api/me"}; self.write("evidence/policy.json",self.policy)
  values={
   "CLOUDFLARE_CONFIG_SNAPSHOT":{"schema":"elite-cloudflare-edge-snapshot/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"zone_id":"zone-123","hostname":"app.example.test","dns_records":[{"id":"dns-1","name":"app.example.test","type":"A","content":"198.51.100.9","proxied":True}],"waf_rulesets":[{"id":"managed-1","name":"Cloudflare Managed","kind":"zone","phase":"http_request_firewall_managed","version":"9","status":"active","enabled_actions":["execute"]},{"id":"custom-1","name":"Project Custom","kind":"zone","phase":"http_request_firewall_custom","version":"3","status":"active","enabled_actions":["block"]}],"tls":{"mode":"full_strict","minimum_version":"1.2","always_use_https":True,"certificate_status":"active"}},
   "ORIGIN_BYPASS_NEGATIVE":{"schema":"elite-cloudflare-edge-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"kind":"ORIGIN_BYPASS_NEGATIVE","target":"203.0.113.10","result":{"connection_succeeded":True,"http_status":403,"blocked":True,"protection_mode":"authenticated_origin_pull"}},
   "PUBLIC_CACHE_BEHAVIOR":{"schema":"elite-cloudflare-edge-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"kind":"PUBLIC_CACHE_BEHAVIOR","target":"https://app.example.test/assets/app.js","result":{"authorization_sent":False,"first_status":200,"second_status":200,"first_cf_cache_status":"MISS","second_cf_cache_status":"HIT","cache_control":"public, max-age=3600","first_body_sha256":"a"*64,"second_body_sha256":"a"*64}},
   "PRIVATE_CACHE_BYPASS":{"schema":"elite-cloudflare-edge-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"kind":"PRIVATE_CACHE_BYPASS","target":"https://app.example.test/api/me","result":{"authorization_sent":True,"status":200,"cf_cache_status":"BYPASS","cache_control":"private, no-store","set_cookie":True,"shared_response_reused":False}}}
  for kind in KINDS:
   stem=kind.lower(); obs=self.write(f"evidence/{stem}.json",values[kind]); args=["--mode",stem,"--zone-id","zone-123","--hostname","app.example.test"]
   out=self.file(f"evidence/{stem}.stdout",b"edge probe pass\n"); err=self.file(f"evidence/{stem}.stderr",b""); target=dict(self.target); target.update({"run_kind":kind,"observation_sha256":self.pf(obs)["sha256"]})
   receipt={"schema":"elite-official-tool-execution/v1","project_id":"revestex","environment":"production","release_digest":digest,"control_id":"EDGE_CDN_WAF","executed_at":z(now),"exit_code":0,"tool":{"name":"Project Cloudflare edge probe","version":"cloudflare-go/7.9.0","source":SOURCE,"sha256":"f"*64},"arguments_sha256":canonical_sha(args),"working_directory":".","target":target,"environment_variable_names":["CLOUDFLARE_API_TOKEN"],"stdout":self.pf(out),"stderr":self.pf(err)}; self.write(f"evidence/{stem}-receipt.json",receipt); self.profile["runs"].append({"kind":kind,"execution_receipt":f"evidence/{stem}-receipt.json","observation":f"evidence/{stem}.json","expected_arguments":args})
 def tearDown(self): self.temp.cleanup()
 def file(self,p,data): q=self.root/p; q.write_bytes(data); return q
 def write(self,p,v): q=self.root/p; q.write_text(json.dumps(v),encoding="utf-8"); return q
 def pf(self,p): data=p.read_bytes(); return {"path":p.relative_to(self.root).as_posix(),"bytes":len(data),"sha256":hashlib.sha256(data).hexdigest()}
 def mutate_observation(self,index,change): p=self.root/self.profile["runs"][index]["observation"]; v=json.loads(p.read_text()); change(v); self.write(self.profile["runs"][index]["observation"],v)
 def test_complete_edge_evidence_emits_semantic_receipt(self): self.assertTrue(all(validate(self.profile,self.root)["assertions"].values()))
 def test_unproxied_dns_is_rejected(self):
  self.mutate_observation(0,lambda v:v["dns_records"][0].update(proxied=False))
  with self.assertRaisesRegex(ValueError,"proxied"): validate(self.profile,self.root)
 def test_missing_managed_waf_is_rejected(self):
  self.mutate_observation(0,lambda v:v.update(waf_rulesets=v["waf_rulesets"][1:]))
  with self.assertRaisesRegex(ValueError,"managed and custom"): validate(self.profile,self.root)
 def test_weak_tls_is_rejected(self):
  self.mutate_observation(0,lambda v:v["tls"].update(mode="flexible"))
  with self.assertRaisesRegex(ValueError,"strict"): validate(self.profile,self.root)
 def test_origin_bypass_success_is_rejected(self):
  self.mutate_observation(1,lambda v:v["result"].update(blocked=False,http_status=200))
  with self.assertRaisesRegex(ValueError,"bypass was not blocked"): validate(self.profile,self.root)
 def test_public_cache_miss_twice_is_rejected(self):
  self.mutate_observation(2,lambda v:v["result"].update(second_cf_cache_status="MISS"))
  with self.assertRaisesRegex(ValueError,"public cache behavior"): validate(self.profile,self.root)
 def test_private_cache_hit_is_rejected(self):
  self.mutate_observation(3,lambda v:v["result"].update(cf_cache_status="HIT",shared_response_reused=True))
  with self.assertRaisesRegex(ValueError,"cache isolation"): validate(self.profile,self.root)
 def test_nonzero_probe_exit_is_rejected(self):
  p=self.root/self.profile["runs"][0]["execution_receipt"]; v=json.loads(p.read_text()); v["exit_code"]=2; self.write(self.profile["runs"][0]["execution_receipt"],v)
  with self.assertRaisesRegex(ValueError,"identity/exit"): validate(self.profile,self.root)
 def test_observation_tamper_is_rejected_by_receipt_binding(self):
  self.mutate_observation(2,lambda v:v["result"].update(cache_control="public, max-age=7200"))
  with self.assertRaisesRegex(ValueError,"target/observation binding"): validate(self.profile,self.root)

if __name__=="__main__": unittest.main()
