"""AUTHORED build orchestration for the complete selected Linux source tree.

Reuses the selected application's cmd packages, standalone Next build, ARCA
generated client/worker and existing quality runner. No business code is copied
or rewritten. The source hash manifest and every invoked tool binary are inputs.
This does not turn a successful compile into deploy/readiness authorization.
"""
from pathlib import Path
import argparse
import os
import subprocess
import sys
from cloud_control import contained,file_hash,raw,read,require
import hashlib
import importlib.util
import shutil

def plan(source, output, tools, inputs):
    source=Path(source).resolve();output=Path(output).absolute()
    require(not output.exists(),"fresh target build output required")
    require(inputs["target"]=="linux/amd64" and inputs["scope"]=="LIBRARY_INFRASTRUCTURE","target build identity")
    for path,h in inputs["source_files"].items():require(file_hash(contained(source,path))==h,"source input drift: "+path)
    require(inputs["source_files"],"source inventory required")
    for name in ("go","node","pnpm","dotnet"):
        require(file_hash(tools[name]["path"])==tools[name]["sha256"],"tool identity mismatch: "+name)
        require(tools[name]["platform"]=="linux/amd64","Windows tooling cannot build this target")
    command_dirs=sorted(p.parent.relative_to(source).as_posix() for p in (source/"cmd").glob("*/main.go"))
    require("cmd/electromobility-api" in command_dirs and "cmd/arca-fiscal-worker" in command_dirs,"full API/fiscal selected commands absent")
    all_sources={p.relative_to(source).as_posix() for p in source.rglob("*") if p.is_file() and ".git" not in p.parts and "node_modules" not in p.parts and ".next" not in p.parts}
    require(set(inputs["source_files"])==all_sources,"complete source inventory required; omissions rejected")
    commands=[]
    # Acquire only versions named by selected lockfiles. No update/latest path.
    for mod in sorted(source.rglob("go.mod")):
        if "node_modules" in mod.parts:continue
        commands.append({"id":"go-deps-"+str(len(commands)),"category":"supply-chain","argv":[tools["go"]["path"],"mod","download"],"cwd":str(mod.parent)})
    require(Path(inputs["offline_pnpm_store"]).is_dir(),"pre-acquired verified offline pnpm store required")
    commands.append({"id":"web-install-locked","category":"supply-chain","argv":[tools["node"]["path"],"--jitless","--no-addons",tools["pnpm"]["path"],"install","--frozen-lockfile","--offline","--ignore-scripts","--ignore-pnpmfile","--store-dir",inputs["offline_pnpm_store"],"--package-import-method=copy","--verify-store-integrity"],"cwd":str(source)})
    for step in inputs.get("generation_steps",[]):
        require(step["owner_path"] in ("tools/assemble_arca_svcutil.py","tools/generate-arca-wsfe-client.ps1"),"only existing admitted generation owners")
        require(file_hash(contained(source,step["owner_path"]))==step["owner_sha256"],"generation owner changed")
        argv=step["argv"];owner=str(contained(source,step["owner_path"]))
        require(isinstance(argv,list) and all(isinstance(x,str) and x for x in argv),"generation argv required")
        require(file_hash(argv[0])==step["executable_sha256"],"generation executable identity changed")
        prefix=[argv[0],"-B",owner] if step["owner_path"].endswith(".py") else [argv[0],"-NoProfile","-File",owner]
        require(argv[:len(prefix)]==prefix,"exact original owner invocation required")
        require(step["input_files"] and all(file_hash(contained(source,rel))==h for rel,h in step["input_files"].items()),"generation source inputs missing or changed")
        commands.append({"id":"arca-generation-"+str(len(commands)),"category":"build","argv":step["argv"],"cwd":str(source)})
    require((source/"arca/fiscal/generated/Elite.Arca.Wsfe.Generated.csproj").is_file() or inputs.get("generation_steps"),"existing ARCA generated source or executable exact-owner generation inputs required")
    commands.append({"id":"arca-restore-locked","category":"supply-chain","argv":[tools["dotnet"]["path"],"restore","cloud/cert_bridge/Elite.Cloud.ArcaLauncher.csproj","--locked-mode","--runtime","linux-x64"],"cwd":str(source)})
    for package in command_dirs:
        commands.append({"id":"go-"+package.split("/")[-1],"category":"build","argv":[tools["go"]["path"],"build","-mod=readonly","-trimpath","-buildvcs=false","-ldflags=-buildid=","-o",str(output/"bin"/package.split("/")[-1]),"./"+package],"cwd":str(source)})
    commands.extend([
        {"id":"web-build","category":"build","argv":[tools["node"]["path"],str(source/"node_modules/next/dist/bin/next"),"build"],"cwd":str(source)},
        {"id":"arca-linux-publish","category":"build","argv":[tools["dotnet"]["path"],"publish","cloud/cert_bridge/Elite.Cloud.ArcaLauncher.csproj","--no-restore","--configuration","Release","--runtime","linux-x64","--self-contained","false","--output",str(output/"arca")],"cwd":str(source)}])
    return {"schema":"elite-linux-build-plan/v403.1","target":"linux/amd64","execution":"NOT_RUN","commands":commands,"required_inputs":["exact adapted pnpm projection from current selected owner; upstream raw archive is not a replacement","all selected Go module dependencies from go.sum","ARCA generation source inputs and Linux NuGet restore from admitted official source/tool manifests","server-action encryption key reference unique to this build, never original V402 fixture key"],"environment":{"GOOS":"linux","GOARCH":"amd64","CGO_ENABLED":"0","GOTOOLCHAIN":"local","GOPROXY":"https://proxy.golang.org","GOSUMDB":"sum.golang.org","NEXT_TELEMETRY_DISABLED":"1","CI":"true","TZ":"UTC"}}

def execute(source,output,tools,inputs,*,authorize_acquisition=False):
    require(authorize_acquisition,"explicit locked-dependency acquisition/build authorization required")
    require(sys.platform=="linux","execute target builds on Linux, not Windows cross-build inference")
    source=Path(source).resolve();output=Path(output).absolute();build=plan(source,output,tools,inputs)
    # Entire adapted pnpm projection must match the existing owner, not just its
    # tiny launcher. Linux runtime compatibility remains separately evidenced.
    projection=Path(tools["pnpm"]["projection_root"]).resolve()
    for rel,entry in read(source/"pnpm_artifact_selection/local-runtime-policy.json")["files"].items():
        require(file_hash(contained(projection,rel))==entry["sha256"],"adapted pnpm projection changed: "+rel)
    runner=source/"ci/run_quality_gates.py"
    require(file_hash(runner)==inputs["quality_runner_sha256"],"existing quality runner changed")
    gateplan={"schema_version":"elite.ci-gates.v1","required_categories":["supply-chain","build"],"max_output_bytes":65536,"fail_fast":True,"gates":[{"id":c["id"],"category":c["category"],"argv":c["argv"],"cwd":str(Path(c["cwd"]).relative_to(source)) or ".","timeout_seconds":1200,"required":True} for c in build["commands"]]}
    output.mkdir(parents=True);gatefile=output/"build-gates.json";gatefile.write_bytes(raw(gateplan))
    env=dict(os.environ);env.update(build["environment"]);env["PATH"]=str(Path(tools["node"]["path"]).parent)+os.pathsep+str(Path(tools["go"]["path"]).parent)+os.pathsep+str(Path(tools["dotnet"]["path"]).parent)+os.pathsep+env.get("PATH","")
    q=subprocess.run([sys.executable,"-B",str(runner),str(gatefile),"--workspace",str(source),"--evidence",str(output/"build-evidence.json")],env=env,timeout=7200,shell=False)
    require(q.returncode==0,"target build failed; exact existing runner evidence retained")
    for path,h in inputs["source_files"].items():require(file_hash(contained(source,path))==h,"build changed locked source: "+path)
    require((source/".next/standalone/server.js").is_file(),"standalone web output absent")
    # Keep source and build metadata outside product images; Dockerfile.web uses
    # standalone + static + public explicitly, .NET/Go use output subtrees.
    artifact={}
    for tree in (output/"bin",output/"arca",source/".next/standalone",source/".next/static",source/"public"):
        if tree.name=="public" and not tree.exists():tree.mkdir()
        require(tree.is_dir(),"required target artifact tree missing: "+tree.name)
        for f in sorted(tree.rglob("*")):
            if f.is_file():artifact[str(f.relative_to(output)) if f.is_relative_to(output) else "web/"+str(f.relative_to(source))]=file_hash(f)
    report={"schema":"elite-linux-built-artifact/v403.1","result":"PASS","claim":"Linux target build only","production_authorized":False,"source_manifest_sha256":hashlib.sha256(raw(inputs["source_files"])).hexdigest(),"source_files":inputs["source_files"],"files":artifact,"artifact_manifest_sha256":hashlib.sha256(raw(artifact)).hexdigest(),"oci_build":"NOT_RUN","target_runtime":"NOT_RUN"}
    (output/"artifact-manifest.json").write_bytes(raw(report));return report

def main():
    p=argparse.ArgumentParser();p.add_argument("--source",required=True);p.add_argument("--output",required=True);p.add_argument("--tools",required=True);p.add_argument("--inputs",required=True);p.add_argument("--execute",action="store_true");p.add_argument("--authorize-locked-acquisition",action="store_true");a=p.parse_args()
    result=execute(a.source,a.output,read(a.tools),read(a.inputs),authorize_acquisition=a.authorize_locked_acquisition) if a.execute else plan(a.source,a.output,read(a.tools),read(a.inputs))
    print(raw(result).decode(),end="")

if __name__=="__main__":main()
