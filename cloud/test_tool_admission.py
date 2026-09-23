"""AUTHORED admission-binding fixtures; do not export synthetic gates as real."""
from pathlib import Path
import tempfile,unittest
from cloud_control import raw,digest,file_hash,read
from tool_admission import validate
class ToolAdmissionTests(unittest.TestCase):
    def setUp(self):
        self.tmp=tempfile.TemporaryDirectory();self.root=Path(self.tmp.name);self.runtime=self.root/"runtime";self.runtime.mkdir();(self.runtime/"tool").write_text("fixture")
        files={"tool":file_hash(self.runtime/"tool")};gates={}
        for i in range(9):
            g="G"+str(i);p=self.root/(g+".json");p.write_bytes(raw({"gate":g,"result":"PASS","method":"EXECUTED_LOCAL","runtime_manifest_sha256":digest(files)}));gates[g]={"path":p.name,"sha256":file_hash(p)}
        self.report={"schema":"elite-cloud-tool-admission/v403.1","result":"PASS","scope":"CLOUD_REFERENCE_TOOL_RUNTIME","production_authorized":False,"binary_sha256":files["tool"],"runtime_root":str(self.runtime),"runtime_files":files,"runtime_manifest_sha256":digest(files),"gates":gates}
        self.path=self.root/"admission.json";self.save()
    def save(self):
        self.path.write_bytes(raw(self.report));self.tool={"binary":str(self.runtime/"tool"),"binary_sha256":self.report["binary_sha256"],"admission":"ADMITTED_EXACT_ARTIFACT","admission_receipt":{"path":str(self.path),"sha256":file_hash(self.path)}}
    def tearDown(self):self.tmp.cleanup()
    def test_bound_runtime_and_executed_owner_reports(self):self.assertEqual(validate(self.tool)["result"],"PASS")
    def test_missing_owner_gate_rejected(self):
        del self.report["gates"]["G7"];self.save()
        with self.assertRaises(ValueError):validate(self.tool)
    def test_extra_runtime_dependency_rejected(self):
        (self.runtime/"unexpected-plugin").write_text("new")
        with self.assertRaises(ValueError):validate(self.tool)
    def test_not_run_gate_and_tamper_rejected(self):
        p=self.root/"G5.json";v=read(p);v["method"]="NOT_RUN";p.write_bytes(raw(v));self.report["gates"]["G5"]["sha256"]=file_hash(p);self.save()
        with self.assertRaises(ValueError):validate(self.tool)
        self.path.write_text("{}")
        with self.assertRaises(ValueError):validate(self.tool)
if __name__=="__main__":unittest.main()
