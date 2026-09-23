from __future__ import annotations
import copy, hashlib, json
from datetime import datetime, timedelta, timezone
from pathlib import Path
import tempfile, unittest
from validate_k6_load_admission import canonical_sha
from validate_postgres_recovery_admission import KINDS, SOURCE, TOOLS, validate

class PostgresRecoveryAdmissionTests(unittest.TestCase):
    def setUp(self):
        self.temp=tempfile.TemporaryDirectory(); self.root=Path(self.temp.name).resolve(); (self.root/"evidence").mkdir(); now=datetime.now(timezone.utc); z=lambda d:d.isoformat().replace("+00:00","Z"); digest="sha256:"+"b"*64
        self.profile={"schema":"elite-postgres-recovery-admission/v1","project_id":"revestex","environment":"production","release_digest":digest,"evaluated_at":z(now),"expires_at":z(now+timedelta(days=7)),"executor_ref":"ci:recovery/42","target":{"cluster":"primary","backup_id":"basebackup-42","recovery_target_time":z(now-timedelta(seconds=10)),"restore_cluster":"isolated-drill-42"},"policy":"evidence/policy.json","drill":"evidence/drill.json","runs":[],"output":"evidence/postgres-control.json"}
        self.write("evidence/policy.json",{"schema":"elite-postgres-recovery-policy/v1","project_id":"revestex","environment":"production","release_digest":digest,"approved_at":z(now-timedelta(days=1)),"approvers":["database-owner","business-owner"],"maximum_rpo_seconds":60,"maximum_rto_seconds":300})
        self.write("evidence/drill.json",{"schema":"elite-postgres-recovery-drill/v1","project_id":"revestex","environment":"production","release_digest":digest,"backup_id":"basebackup-42","recovery_target_time":self.profile["target"]["recovery_target_time"],"recovered_through":z(now-timedelta(seconds=20)),"restore_started_at":z(now-timedelta(seconds=200)),"restore_ready_at":z(now-timedelta(seconds=20)),"wal_archive":{"timeline":1,"start_lsn":"0/1000000","end_lsn":"0/2000000","location_digest":"sha256:"+"c"*64},"invariants":[{"name":"orders-total","pass":True},{"name":"ledger-balance","pass":True}]})
        for kind in KINDS:
            stem=kind.lower(); inp=self.root/"evidence"/f"{stem}.input"; inp.write_text(kind+"\n"); args=["--elite-fixture",kind]; out=self.root/"evidence"/f"{stem}.stdout.bin"; err=self.root/"evidence"/f"{stem}.stderr.bin"; out.write_bytes(b"PASS\n"); err.write_bytes(b"")
            pf=lambda p:{"path":p.relative_to(self.root).as_posix(),"bytes":p.stat().st_size,"sha256":hashlib.sha256(p.read_bytes()).hexdigest()}; target=dict(self.profile["target"]); target.update({"run_kind":kind,"input_sha256s":[pf(inp)["sha256"]]})
            receipt={"schema":"elite-official-tool-execution/v1","project_id":"revestex","environment":"production","release_digest":digest,"control_id":"POSTGRES_RECOVERY","executed_at":z(now),"exit_code":0,"tool":{"name":TOOLS[kind],"version":"18.6","source":SOURCE,"sha256":"d"*64},"arguments_sha256":canonical_sha(args),"working_directory":".","target":target,"environment_variable_names":[],"stdout":pf(out),"stderr":pf(err)}; self.write(f"evidence/{stem}-receipt.json",receipt); self.profile["runs"].append({"kind":kind,"execution_receipt":f"evidence/{stem}-receipt.json","expected_arguments":args,"input_files":[f"evidence/{stem}.input"]})
    def tearDown(self): self.temp.cleanup()
    def write(self,path,value): (self.root/path).write_text(json.dumps(value),encoding="utf-8")
    def test_full_restore_pitr_and_objectives_emit_semantic_receipt(self): self.assertTrue(all(validate(self.profile,self.root)["assertions"].values()))
    def test_verifybackup_nonzero_is_rejected(self):
        p=self.root/self.profile["runs"][0]["execution_receipt"]; v=json.loads(p.read_text()); v["exit_code"]=1; self.write(self.profile["runs"][0]["execution_receipt"],v)
        with self.assertRaisesRegex(ValueError,"identity/exit"): validate(self.profile,self.root)
    def test_logical_or_verify_only_without_five_runs_is_rejected(self):
        v=copy.deepcopy(self.profile); v["runs"]=v["runs"][:1]
        with self.assertRaisesRegex(ValueError,"five ordered"): validate(v,self.root)
    def test_rpo_breach_is_rejected(self):
        p=self.root/self.profile["drill"]; v=json.loads(p.read_text()); v["recovered_through"]="2020-01-01T00:00:00Z"; self.write(self.profile["drill"],v)
        with self.assertRaisesRegex(ValueError,"RPO exceeds"): validate(self.profile,self.root)
    def test_rto_breach_is_rejected(self):
        p=self.root/self.profile["drill"]; v=json.loads(p.read_text()); v["restore_started_at"]="2020-01-01T00:00:00Z"; self.write(self.profile["drill"],v)
        with self.assertRaisesRegex(ValueError,"RTO exceeds"): validate(self.profile,self.root)
    def test_failed_invariant_is_rejected(self):
        p=self.root/self.profile["drill"]; v=json.loads(p.read_text()); v["invariants"][0]["pass"]=False; self.write(self.profile["drill"],v)
        with self.assertRaisesRegex(ValueError,"invariants"): validate(self.profile,self.root)
    def test_input_tamper_is_rejected(self):
        (self.root/self.profile["runs"][3]["input_files"][0]).write_text("tampered")
        with self.assertRaisesRegex(ValueError,"binding mismatch"): validate(self.profile,self.root)
    def test_tool_version_drift_is_rejected(self):
        p=self.root/self.profile["runs"][4]["execution_receipt"]; v=json.loads(p.read_text()); v["tool"]["version"]="18.5"; self.write(self.profile["runs"][4]["execution_receipt"],v)
        with self.assertRaisesRegex(ValueError,"tool identity"): validate(self.profile,self.root)

if __name__=="__main__": unittest.main()
