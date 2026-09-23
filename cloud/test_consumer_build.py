"""Consumer sealing and build-plan regression fixtures; no target runtime claim."""
from pathlib import Path
import copy,json,subprocess,tempfile,unittest
from unittest.mock import patch
import sys
from cloud_control import file_hash,raw,validate_profile
from seal_consumer import seal,validate_selection,BASE,EXT,AUTHORITIES,EXTRAS
from bootstrap import checkout
from build_target import plan

class ConsumerTests(unittest.TestCase):
    def setUp(self): self.tmp=tempfile.TemporaryDirectory();self.root=Path(self.tmp.name)
    def tearDown(self): self.tmp.cleanup()
    def fixture(self):
        repo=self.root/"library";repo.mkdir()
        base={"packs":[{"packId":"BASE-"+str(i),"path":"packs/"+str(i)+".md","version":"1","files":["*"],"variables":{}} for i in range(116)]}
        ext=copy.deepcopy(base);ext["packs"] += [{"packId":x,"path":"packs/"+x+".md","version":"0.2.0"} for x in sorted(EXTRAS)]
        for x in ext["packs"]:
            p=repo/x["path"];p.parent.mkdir(exist_ok=True);p.write_text(x["packId"])
        for x in AUTHORITIES:
            p=repo/x;p.parent.mkdir(parents=True,exist_ok=True);p.write_text("fixture authority "+x)
        for path,obj in [(BASE,base),(EXT,ext)]: (repo/path).write_text("```json\n"+json.dumps(obj)+"\n```\n")
        for argv in [["git","init","-q"],["git","config","core.autocrlf","false"],["git","add","."],["git","-c","user.name=Fixture","-c","user.email=fixture@example.test","commit","-qm","fixture"]]:subprocess.run(argv,cwd=repo,check=True,capture_output=True)
        commit=subprocess.check_output(["git","rev-parse","HEAD"],cwd=repo).decode().strip();return repo,commit,base,ext
    def test_full_120_git_object_seal_distinct_checkout(self):
        repo,commit,_,_=self.fixture();out=self.root/"consumer";receipt=seal(repo,commit,out,str(repo),fixture=True)
        lock=json.loads((out/"library-checkout.lock.json").read_text());checkout(lock,self.root/"worker",fixture=True)
        self.assertEqual(validate_profile(json.loads((out/"composition.json").read_text()),self.root/"worker")["pack_count"],120)
        self.assertFalse(receipt["production_authorized"])
    def test_seal_reads_commit_not_dirty_live_tree(self):
        repo,commit,_,_=self.fixture();(repo/EXT).write_text("dirty uncommitted")
        self.assertEqual(seal(repo,commit,self.root/"consumer",str(repo),fixture=True)["result"],"PASS")
    def test_base_mutation_and_duplicate_extra_rejected(self):
        _,_,base,ext=self.fixture();ext["packs"][1]["version"]="2"
        with self.assertRaises(ValueError):validate_selection(base,ext)
        ext=copy.deepcopy(base);ext["packs"] += [{"packId":next(iter(EXTRAS))}]*4
        with self.assertRaises(ValueError):validate_selection(base,ext)
    def test_seal_cannot_write_inside_library(self):
        repo,commit,_,_=self.fixture()
        with self.assertRaises(ValueError):seal(repo,commit,repo/"consumer",str(repo),fixture=True)
    def test_build_plan_covers_sources_and_locked_offline_install(self):
        source=self.root/"source";source.mkdir();store=self.root/"offline-store";store.mkdir()
        for rel in ["cmd/electromobility-api/main.go","cmd/arca-fiscal-worker/main.go","go.mod","package.json","arca/fiscal/generated/Elite.Arca.Wsfe.Generated.csproj"]:
            p=source/rel;p.parent.mkdir(parents=True,exist_ok=True);p.write_text("fixture")
        tools={}
        for name in ["go","node","pnpm","dotnet"]:
            p=self.root/name;p.write_text("fixturetool");tools[name]={"path":str(p),"sha256":file_hash(p),"platform":"linux/amd64"}
        inputs={"target":"linux/amd64","scope":"LIBRARY_INFRASTRUCTURE","offline_pnpm_store":str(store),"source_files":{p.relative_to(source).as_posix():file_hash(p) for p in source.rglob("*") if p.is_file()}}
        result=plan(source,self.root/"build",tools,inputs);install=next(x for x in result["commands"] if x["id"]=="web-install-locked")
        self.assertIn("--offline",install["argv"]);self.assertIn("--ignore-pnpmfile",install["argv"]);self.assertIn("cloud/cert_bridge/Elite.Cloud.ArcaLauncher.csproj",str(result));self.assertEqual(result["execution"],"NOT_RUN")
        del inputs["source_files"]["go.mod"]
        with self.assertRaises(ValueError):plan(source,self.root/"build",tools,inputs)
    def test_oci_packaging_binds_exact_outputs_and_all_five_images(self):
        from oci_target import execute,CONTROLLER_SOURCES
        source=self.root/"src";build=self.root/"built";(source/"cloud").mkdir(parents=True);(build/"bin").mkdir(parents=True)
        (build/"bin/api").write_bytes(b"fixture-binary");(source/".next/standalone").mkdir(parents=True);(source/".next/standalone/server.js").write_text("fixture-js")
        for name in ["Dockerfile.web","Dockerfile.workers","Dockerfile.arca","Dockerfile.migrations","Dockerfile.controller"]:(source/"cloud"/name).write_text("fixture Dockerfile")
        for name in CONTROLLER_SOURCES:
            p=source/name;p.parent.mkdir(parents=True,exist_ok=True);p.write_text("fixture controller or unchanged runner source")
        (source/"deploy").mkdir();(source/"deploy/postgres-migrate.sh").write_text("fixture migration owner")
        (source/"cloud/image-lock.json").write_bytes(raw({"images":[{"id":id,"reference":"official.example/"+id+"@sha256:"+"a"*64} for id in ["go-runtime","node-runtime","dotnet-runtime","postgres-client","gcloud-runtime"]]}))
        inventory={"bin/api":file_hash(build/"bin/api"),"web/.next/standalone/server.js":file_hash(source/".next/standalone/server.js")}
        (build/"artifact-manifest.json").write_bytes(raw({"schema":"elite-linux-built-artifact/v403.1","result":"PASS","production_authorized":False,"source_files":{p.relative_to(source).as_posix():file_hash(p) for p in source.rglob('*') if p.is_file()},"files":inventory}))
        called=[]
        def fake(argv,cwd,env,timeout):
            called.append(argv);Path(argv[argv.index("--metadata-file")+1]).write_bytes(raw({"containerimage.digest":"sha256:"+"b"*64}));Path(argv[argv.index("--output")+1].split("dest=",1)[1]).write_bytes(b"fixture-oci");return subprocess.CompletedProcess(argv,0)
        tool={"binary":sys.executable,"binary_sha256":file_hash(sys.executable),"admission":"ADMITTED_EXACT_ARTIFACT"}
        with patch("oci_target.load_runner",return_value=fake),patch("tool_admission.validate",return_value={}):result=execute(source,build,self.root/"oci",tool,{"path":"fixture","sha256":"fixture"},authorize=True)
        self.assertEqual(set(result["components"]),{"workers","web","arca","migrations","controller"});self.assertEqual(len(called),5);self.assertTrue(all("--network=none" in x and "--push" not in x for x in called));self.assertEqual(result["cloud"],"NOT_RUN")
        for name in CONTROLLER_SOURCES:self.assertEqual(file_hash(self.root/"oci/context"/name),file_hash(source/name))
        (build/"bin/api").write_bytes(b"changed")
        with patch("oci_target.load_runner") as runner,patch("tool_admission.validate",return_value={}):
            with self.assertRaises(ValueError):execute(source,build,self.root/"oci2",tool,{"path":"fixture","sha256":"fixture"},authorize=True)
            runner.assert_not_called()

if __name__=="__main__":unittest.main()
