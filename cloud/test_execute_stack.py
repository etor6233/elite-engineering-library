"""Whole-DAG behavioral fixtures through unchanged bounded_process; no cloud."""
from pathlib import Path
import copy,json,os,subprocess,sys,tempfile,unittest
from unittest.mock import patch
from cloud_control import raw,read,digest,file_hash
from execute_stack import envelope,execute
from stack_manifest import build_stack
from test_stack_manifest import fixture_bundle,SOURCE,ROOT
from semantic_contract import GATE_CHECKS,resource_checks

NOW="2026-09-14T12:00:00Z";EXP="2026-09-14T12:10:00Z"
class ExecuteStackTests(unittest.TestCase):
    def setUp(self):
        self.tmp=tempfile.TemporaryDirectory();self.root=Path(self.tmp.name);self.art=self.root/"artifacts";self.art.mkdir();self.journal=self.root/"intent";self.journal.mkdir()
        b=fixture_bundle();generated=build_stack(b["config"],b["images"],source_root=SOURCE)
        for name,value in generated["manifests"].items():p=self.art/name;p.parent.mkdir(parents=True,exist_ok=True);p.write_bytes(raw(value))
        self.plan=envelope(b["config"],b["images"],SOURCE,ROOT/"cloud/composition.json",self.art,{"backend":"LOCAL_FIXTURE","directory":str(self.journal)})
        self.tool={"binary":sys.executable,"binary_sha256":file_hash(sys.executable),"method":"SIMULATED_PROVIDER","admission":"FIXTURE_ONLY","prefix_arguments":[str(ROOT/"cloud/fake_stack_provider.py"),str(self.art/"synthetic-state.json")]}
        owner=Path(os.environ.get("ELITE_CLOUD_OWNER_ROOT",ROOT/"cloud/owners/secure-ops/production_admission_gate"))/"run_official_tool.py"
        if not owner.exists():owner=ROOT/"production_admission_gate/run_official_tool.py"
        self.runner={"path":str(owner),"sha256":file_hash(owner)};self.runs={};self.counter=0
    def tearDown(self):self.tmp.cleanup()
    def approval(self):return {"plan_sha256":digest(self.plan),"scope":"CLOUD_REFERENCE_EXECUTION","environment":self.plan["environment"],"approved_by":"synthetic-fixture-operator","approved_at":NOW,"expires_at":EXP,"provider_access_proven":True,"budget_authorized":True}
    def write(self,body):
        self.counter+=1;p=self.art/("proof-"+str(self.counter)+".json");p.write_bytes(raw(body));return {"path":p.name,"sha256":file_hash(p)}
    def ob(self,kind,step,name="",run=None):
        common={"method":"SIMULATED_PROVIDER","target_sha256":self.plan["stack"]["target_sha256"],"images_sha256":self.plan["stack"]["images_sha256"],"result":"PASS","observed_at":NOW,"expires_at":EXP}
        claim="TARGET_GATE_VERIFIED" if kind=="TARGET_GATE_EVIDENCE" else "STACK_RESOURCE_STATE_RECONCILED"
        checks=GATE_CHECKS[name] if kind=="TARGET_GATE_EVIDENCE" else resource_checks(step)
        proof=self.write(common|{"schema":"elite-cloud-target-semantic-report/v403.1","claim":claim,"gate_or_resource":name if kind=="TARGET_GATE_EVIDENCE" else step["id"],"for_step":step["id"],"desired_step_sha256":digest(step),"checks":{k:True for k in checks},"fixture_notice":"SYNTHETIC_STATE_COMPARISON_ONLY"})
        body=common|{"schema":"elite-cloud-stack-observation/v403.1","plan_sha256":digest(self.plan),"kind":kind,"for_step":step["id"],"name":name,"source_evidence":proof}
        if run:body.update(execution_receipt=run,desired_step_sha256=digest(step))
        return self.write(body)
    def refs(self,step):
        gates={x:self.ob("TARGET_GATE_EVIDENCE",step,x) for x in step["required_target_receipts"]}
        deps={x:self.ob("STACK_RESOURCE_RECONCILIATION",next(s for s in self.plan["stack"]["commands"] if s["id"]==x),run=self.runs[x]) for x in step["depends_on"]}
        return {"gates":gates,"dependencies":deps}
    def first(self):return next(s for s in self.plan["stack"]["commands"] if not s["depends_on"] and s["effect"]=="WRITE")
    def test_complete_dag_executes_each_effect_once_with_separate_observations(self):
        pending=list(self.plan["stack"]["commands"])
        while pending:
            ready=next(s for s in pending if set(s["depends_on"])<=set(self.runs));receipt=self.art/(ready["id"]+".receipt.json")
            result=execute(self.plan,ready["id"],self.approval(),NOW,self.tool,self.runner,receipt,self.refs(ready),fixture=True)
            self.assertFalse(result["cloud_demonstrated"]);self.runs[ready["id"]]={"path":receipt.name,"sha256":file_hash(receipt)};pending.remove(ready)
        state=read(self.art/"synthetic-state.json");self.assertEqual(len(state),sum(s["effect"]=="WRITE" for s in self.plan["stack"]["commands"]));self.assertEqual(len(self.runs),len(self.plan["stack"]["commands"]))
    def test_changed_argv_and_manifests_reject_before_effect(self):
        step=self.first();refs=self.refs(step);original=copy.deepcopy(self.plan);self.plan["stack"]["commands"][0]["argv"].append("--allow-unauthenticated")
        with self.assertRaises(ValueError):execute(self.plan,step["id"],self.approval(),NOW,self.tool,self.runner,self.art/"bad.json",refs,fixture=True)
        self.plan=original;p=next(iter(self.plan["stack"]["manifests"]));(self.art/p).write_bytes(b"changed")
        with self.assertRaises(ValueError):execute(self.plan,step["id"],self.approval(),NOW,self.tool,self.runner,self.art/"bad2.json",refs,fixture=True)
        self.assertFalse((self.art/"synthetic-state.json").exists())
    def test_missing_gate_and_expired_underlying_evidence_reject(self):
        step=self.first();refs=self.refs(step);key=next(iter(refs["gates"]));missing=copy.deepcopy(refs);missing["gates"].pop(key)
        with self.assertRaises(ValueError):execute(self.plan,step["id"],self.approval(),NOW,self.tool,self.runner,self.art/"bad.json",missing,fixture=True)
        ob=read(self.art/refs["gates"][key]["path"]);p=self.art/ob["source_evidence"]["path"];v=read(p);v.update(observed_at="2020-01-01T00:00:00Z",expires_at="2020-01-01T00:15:00Z");p.write_bytes(raw(v));ob["source_evidence"]["sha256"]=file_hash(p);refs["gates"][key]=self.write(ob)
        with self.assertRaises(ValueError):execute(self.plan,step["id"],self.approval(),NOW,self.tool,self.runner,self.art/"bad2.json",refs,fixture=True)
    def test_ambiguous_effect_never_replays_with_fresh_receipt_name(self):
        step=self.first();refs=self.refs(step);calls=[]
        def uncertain(*args):calls.append(True);raise ValueError("synthetic crash after effect")
        with patch("execute_step.load_runner",return_value=uncertain):
            with self.assertRaises(ValueError):execute(self.plan,step["id"],self.approval(),NOW,self.tool,self.runner,self.art/"uncertain.json",refs,fixture=True)
            with self.assertRaises(FileExistsError):execute(self.plan,step["id"],self.approval(),NOW,self.tool,self.runner,self.art/"retry.json",refs,fixture=True)
        self.assertEqual(len(calls),1);self.assertEqual(read(self.art/"uncertain.json")["result"],"UNKNOWN_RECONCILE_BEFORE_RETRY")
    def test_real_artifact_label_without_receipt_tree_is_rejected(self):
        from execute_stack import validate_envelope
        self.plan["images"]["method"]="EXECUTED_TARGET"
        self.plan["stack"]=build_stack(self.plan["config"],self.plan["images"],source_root=SOURCE)
        with self.assertRaisesRegex(ValueError,"artifact admission receipts"):validate_envelope(self.plan)
    def test_other_claim_and_arbitrary_true_checks_are_rejected(self):
        step=self.first();refs=self.refs(step);key=next(iter(refs["gates"]));ob=read(self.art/refs["gates"][key]["path"]);p=self.art/ob["source_evidence"]["path"];proof=read(p)
        proof["checks"]={"foo":True};p.write_bytes(raw(proof));ob["source_evidence"]["sha256"]=file_hash(p);refs["gates"][key]=self.write(ob)
        with self.assertRaisesRegex(ValueError,"semantic checks"):execute(self.plan,step["id"],self.approval(),NOW,self.tool,self.runner,self.art/"bad.json",refs,fixture=True)
        refs=self.refs(step);ob=read(self.art/refs["gates"][key]["path"]);p=self.art/ob["source_evidence"]["path"];proof=read(p);proof["for_step"]="another-resource";p.write_bytes(raw(proof));ob["source_evidence"]["sha256"]=file_hash(p);refs["gates"][key]=self.write(ob)
        with self.assertRaisesRegex(ValueError,"different claim/resource"):execute(self.plan,step["id"],self.approval(),NOW,self.tool,self.runner,self.art/"bad2.json",refs,fixture=True)
    def test_boolean_exit_code_cannot_prove_dependency_success(self):
        step=next(s for s in self.plan["stack"]["commands"] if s["depends_on"]);dep=step["depends_on"][0]
        self.runs={x:self.write({"schema":"elite-cloud-step-execution/v403.1","plan_sha256":digest(self.plan),"step":x,"method":"SIMULATED_PROVIDER","exit_code":False}) for x in step["depends_on"]}
        with self.assertRaisesRegex(ValueError,"dependency execution"):execute(self.plan,step["id"],self.approval(),NOW,self.tool,self.runner,self.art/"bad.json",self.refs(step),fixture=True)
    def test_barrier_receipt_cannot_replace_a_write_receipt(self):
        commands={s["id"]:s for s in self.plan["stack"]["commands"]}
        step=next(s for s in commands.values() if s["depends_on"] and all(commands[d]["effect"]=="WRITE" for d in s["depends_on"]))
        self.runs={x:self.write({"schema":"elite-cloud-step-execution/v403.1","plan_sha256":digest(self.plan),"step":x,"method":"SIMULATED_PROVIDER","result":"BARRIER_EVIDENCE_CHECKED"}) for x in step["depends_on"]}
        with self.assertRaisesRegex(ValueError,"dependency execution"):execute(self.plan,step["id"],self.approval(),NOW,self.tool,self.runner,self.art/"bad.json",self.refs(step),fixture=True)

if __name__=="__main__":unittest.main()
