from __future__ import annotations
import copy,hashlib,json,tempfile,unittest
from datetime import datetime,timedelta,timezone
from pathlib import Path
from validate_deploy_rollback_admission import KINDS,KUBECTL_SHA,KUBECTL_SOURCE,validate
from validate_k6_load_admission import canonical_sha

class DeployRollbackAdmissionTests(unittest.TestCase):
 def setUp(self):
  self.temp=tempfile.TemporaryDirectory(); self.root=Path(self.temp.name); (self.root/"evidence").mkdir(); now=datetime.now(timezone.utc); z=lambda d:d.isoformat().replace("+00:00","Z"); digest="sha256:"+"e"*64; desired="registry.example.test/revestex@sha256:"+"a"*64; previous="registry.example.test/revestex@sha256:"+"b"*64
  self.target={"cluster_uid_sha256":"sha256:"+"c"*64,"context":"prod-ar","namespace":"revestex","workload_kind":"Deployment","workload_name":"api","container_name":"api","service_url":"https://api.example.test/health"}
  self.profile={"schema":"elite-deploy-rollback-admission/v1","project_id":"revestex","environment":"production","release_digest":digest,"evaluated_at":z(now),"expires_at":z(now+timedelta(days=7)),"executor_ref":"deploy:change/112","tool":{"name":"Kubernetes kubectl","version":"1.37.0","source":KUBECTL_SOURCE,"sha256":KUBECTL_SHA},"target":self.target,"policy":"evidence/policy.json","runs":[],"output":"evidence/deploy-control.json"}
  self.policy={"schema":"elite-deploy-rollback-policy/v1","project_id":"revestex","environment":"production","release_digest":digest,"approved_at":z(now),"approvers":["platform-owner","business-owner"],"target":self.target,"desired_image_digest":desired,"previous_image_digest":previous,"desired_revision":42,"previous_revision":41,"timeout_seconds":600,"canary":{"minimum_samples":3,"max_error_rate":0.01,"max_p95_latency_ms":500,"automatic_rollback":True},"migration":{"migration_id":"expand-112","expand_contract_required":True,"forward_backward_required":True,"rollback_data_recovery_ref":"recovery:112"}}
  self.write("evidence/policy.json",self.policy)
  common={"schema":"elite-deploy-rollback-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now)}
  observations={
   "PREFLIGHT":dict(common,kind="PREFLIGHT",client_version="v1.37.0",server_version="v1.37.0",cluster_uid_sha256=self.target["cluster_uid_sha256"],context="prod-ar",namespace="revestex",authorized_verbs=["get","list","patch","watch"],workload_found=True),
   "APPLY_IMMUTABLE_DIGEST":dict(common,kind="APPLY_IMMUTABLE_DIGEST",workload_kind="Deployment",workload_name="api",container_name="api",requested_image_digest=desired,observed_image_digest=desired,revision_before=41,revision_after=42,generation=12,observed_generation=12),
   "CANARY_HEALTH":dict(common,kind="CANARY_HEALTH",service_url=self.target["service_url"],samples=5,healthy_samples=5,error_rate=0.0,p95_latency_ms=120,probes_pass=True,alerts_firing=[],error_budget_breached=False),
   "ROLLOUT_STATUS":dict(common,kind="ROLLOUT_STATUS",workload_kind="Deployment",workload_name="api",progressing=True,progressing_reason="NewReplicaSetAvailable",available=True,desired_replicas=3,updated_replicas=3,available_replicas=3,old_replicas=0,application_probes_pass=True),
   "ROLLBACK":dict(common,kind="ROLLBACK",workload_kind="Deployment",workload_name="api",from_revision=42,to_revision=41,rollback_command_issued=True,rollout_complete=True),
   "POST_ROLLBACK_VERIFY":dict(common,kind="POST_ROLLBACK_VERIFY",observed_image_digest=previous,desired_replicas=3,available_replicas=3,service_health_pass=True,data_invariants_pass=True,alerts_firing=[]),
   "MIGRATION_COMPATIBILITY":dict(common,kind="MIGRATION_COMPATIBILITY",migration_id="expand-112",expand_contract_used=True,forward_compatible=True,backward_compatible=True,old_new_versions_coexisted=True,rollback_data_safe=True,integrity_checks_pass=True,destructive_change_before_convergence=False)}
  for kind in KINDS:
   stem=kind.lower(); obs=self.write(f"evidence/{stem}.json",observations[kind]); args=["--context","prod-ar","--namespace","revestex",kind.lower()]; out=self.file(f"evidence/{stem}.stdout",b"pass\n"); err=self.file(f"evidence/{stem}.stderr",b""); target=dict(self.target); target.update({"run_kind":kind,"observation_sha256":self.pf(obs)["sha256"]}); receipt={"schema":"elite-official-tool-execution/v1","project_id":"revestex","environment":"production","release_digest":digest,"control_id":"DEPLOY_ROLLBACK","executed_at":z(now),"exit_code":0,"tool":{"name":"Kubernetes kubectl","version":"1.37.0","source":KUBECTL_SOURCE,"sha256":KUBECTL_SHA},"arguments_sha256":canonical_sha(args),"working_directory":".","target":target,"environment_variable_names":["KUBECONFIG_REF"],"stdout":self.pf(out),"stderr":self.pf(err)}; self.write(f"evidence/{stem}-receipt.json",receipt); self.profile["runs"].append({"kind":kind,"execution_receipt":f"evidence/{stem}-receipt.json","observation":f"evidence/{stem}.json","expected_arguments":args})
 def tearDown(self): self.temp.cleanup()
 def file(self,p,data): q=self.root/p; q.write_bytes(data); return q
 def write(self,p,v): q=self.root/p; q.write_text(json.dumps(v),encoding="utf-8"); return q
 def pf(self,p): data=p.read_bytes(); return {"path":p.relative_to(self.root).as_posix(),"bytes":len(data),"sha256":hashlib.sha256(data).hexdigest()}
 def mutate(self,index,change): p=self.root/self.profile["runs"][index]["observation"]; v=json.loads(p.read_text()); change(v); self.write(self.profile["runs"][index]["observation"],v)
 def test_complete_realistic_evidence_emits_semantic_receipt(self): self.assertTrue(all(validate(self.profile,self.root)["assertions"].values()))
 def test_mutable_image_is_rejected(self):
  p=self.root/self.profile["policy"]; v=json.loads(p.read_text()); v["desired_image_digest"]="registry.example.test/revestex:latest"; self.write(self.profile["policy"],v)
  with self.assertRaisesRegex(ValueError,"immutable registry digest"): validate(self.profile,self.root)
 def test_apply_digest_drift_is_rejected(self):
  self.mutate(1,lambda v:v.update(observed_image_digest="registry.example.test/revestex@sha256:"+"d"*64))
  with self.assertRaisesRegex(ValueError,"desired image digest"): validate(self.profile,self.root)
 def test_canary_budget_breach_is_rejected(self):
  self.mutate(2,lambda v:v.update(error_budget_breached=True))
  with self.assertRaisesRegex(ValueError,"canary health"): validate(self.profile,self.root)
 def test_rollout_old_replicas_are_rejected(self):
  self.mutate(3,lambda v:v.update(old_replicas=1))
  with self.assertRaisesRegex(ValueError,"rollout status"): validate(self.profile,self.root)
 def test_rollback_not_executed_is_rejected(self):
  self.mutate(4,lambda v:v.update(rollback_command_issued=False))
  with self.assertRaisesRegex(ValueError,"actual rollback"): validate(self.profile,self.root)
 def test_previous_digest_not_restored_is_rejected(self):
  self.mutate(5,lambda v:v.update(observed_image_digest=self.policy["desired_image_digest"]))
  with self.assertRaisesRegex(ValueError,"previous digest"): validate(self.profile,self.root)
 def test_unsafe_migration_is_rejected(self):
  self.mutate(6,lambda v:v.update(backward_compatible=False))
  with self.assertRaisesRegex(ValueError,"migration compatibility"): validate(self.profile,self.root)
 def test_wrong_cluster_is_rejected(self):
  self.mutate(0,lambda v:v.update(cluster_uid_sha256="sha256:"+"f"*64))
  with self.assertRaisesRegex(ValueError,"preflight target"): validate(self.profile,self.root)
 def test_nonzero_kubectl_exit_is_rejected(self):
  p=self.root/self.profile["runs"][3]["execution_receipt"]; v=json.loads(p.read_text()); v["exit_code"]=1; self.write(self.profile["runs"][3]["execution_receipt"],v)
  with self.assertRaisesRegex(ValueError,"identity/exit"): validate(self.profile,self.root)
 def test_observation_tamper_breaks_receipt_binding(self):
  self.mutate(3,lambda v:v.update(desired_replicas=4,updated_replicas=4,available_replicas=4))
  with self.assertRaisesRegex(ValueError,"target/observation binding"): validate(self.profile,self.root)
 def test_wrong_kubectl_hash_is_rejected(self):
  p=self.root/self.profile["runs"][0]["execution_receipt"]; v=json.loads(p.read_text()); v["tool"]["sha256"]="0"*64; self.write(self.profile["runs"][0]["execution_receipt"],v)
  with self.assertRaisesRegex(ValueError,"exact kubectl identity"): validate(self.profile,self.root)
 def test_secret_kubeconfig_name_is_rejected(self):
  p=self.root/self.profile["runs"][0]["execution_receipt"]; v=json.loads(p.read_text()); v["environment_variable_names"]=["KUBECONFIG"]; self.write(self.profile["runs"][0]["execution_receipt"],v)
  with self.assertRaisesRegex(ValueError,"references, not secrets"): validate(self.profile,self.root)
 def test_missing_ordered_run_is_rejected(self):
  self.profile["runs"]=self.profile["runs"][:-1]
  with self.assertRaisesRegex(ValueError,"seven ordered"): validate(self.profile,self.root)

if __name__=="__main__": unittest.main()
