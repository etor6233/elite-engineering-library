"""AUTHORED exact build-output packaging over Docker buildx and bounded owner.

Exports OCI archives without pushing or deploying. Newly pinned Linux images
remain candidates until their own license/SCA/runtime evidence is supplied.
"""
from pathlib import Path
import argparse,shutil,os
from cloud_control import contained,file_hash,raw,read,require
from execute_step import load_runner

CONTROLLER_SOURCES=["cloud/finops_cycle.py","cloud/cloud_control.py","cloud/execute_step.py","cloud/tool_admission.py","production_admission_gate/run_official_tool.py","production_admission_gate/validate_production_admission.py"]

def prepare(source,build,destination):
    source=Path(source).resolve();build=Path(build).resolve();destination=Path(destination).absolute()
    require(not destination.exists(),"fresh OCI context required")
    report=read(build/"artifact-manifest.json")
    require(report["schema"]=="elite-linux-built-artifact/v403.1" and report["result"]=="PASS" and report["production_authorized"] is False,"Linux build receipt required")
    for name in ["cloud/Dockerfile.web","cloud/Dockerfile.workers","cloud/Dockerfile.arca","cloud/Dockerfile.migrations","cloud/Dockerfile.controller","cloud/image-lock.json","deploy/postgres-migrate.sh",*CONTROLLER_SOURCES]:
        require(file_hash(contained(source,name))==report["source_files"].get(name),"packaging source changed after target build")
    for name,expected in report["files"].items():
        actual=contained(source,name[4:]) if name.startswith("web/") else contained(build,name)
        require(file_hash(actual)==expected,"built artifact changed before packaging")
    destination.mkdir(parents=True)
    for name in report["files"]:
        actual=contained(source,name[4:]) if name.startswith("web/") else contained(build,name)
        rel=name[4:] if name.startswith("web/") else "build/"+name
        out=contained(destination,rel);out.parent.mkdir(parents=True,exist_ok=True);shutil.copyfile(actual,out)
    for name in ["Dockerfile.web","Dockerfile.workers","Dockerfile.arca","Dockerfile.migrations","Dockerfile.controller"]:shutil.copyfile(source/"cloud"/name,destination/name)
    for name in CONTROLLER_SOURCES:
        out=contained(destination,name);out.parent.mkdir(parents=True,exist_ok=True);shutil.copyfile(contained(source,name),out)
    for name,h in report["source_files"].items():
        if name.startswith("db/migrations/") or name=="deploy/postgres-migrate.sh":
            require(file_hash(contained(source,name))==h,"migration source changed after build")
            out=contained(destination,name);out.parent.mkdir(parents=True,exist_ok=True);shutil.copyfile(source/name,out)
    (destination/"public").mkdir(exist_ok=True)
    return report

def execute(source,build,output,tool,runner,*,authorize=False):
    require(authorize,"explicit OCI candidate build authorization required")
    require(tool["admission"]=="ADMITTED_EXACT_ARTIFACT" and file_hash(tool["binary"])==tool["binary_sha256"],"OCI engine admission/identity required")
    from tool_admission import validate
    validate(tool)
    output=Path(output).absolute();context=output/"context";source=Path(source).resolve()
    require(not output.exists(),"fresh OCI output required")
    lock=read(source/"cloud/image-lock.json");images={x["id"]:x for x in lock["images"]}
    # Candidate execution may build locally; production admission still separate.
    output.mkdir(parents=True);built=prepare(source,build,context)
    bounded=load_runner(runner["path"],runner["sha256"]);receipts={}
    for component,dockerfile,arg,imageid in [("workers","Dockerfile.workers","GO_RUNTIME","go-runtime"),("web","Dockerfile.web","NODE_RUNTIME","node-runtime"),("arca","Dockerfile.arca","DOTNET_RUNTIME","dotnet-runtime"),("migrations","Dockerfile.migrations","POSTGRES_CLIENT","postgres-client"),("controller","Dockerfile.controller","GCLOUD_RUNTIME","gcloud-runtime")]:
        image=images[imageid];require("@sha256:" in image["reference"],"immutable runtime image required")
        metadata=output/(component+"-metadata.json");archive=output/(component+".oci.tar")
        argv=[tool["binary"],"buildx","build","--platform","linux/amd64","--network=none","--provenance=mode=max","--build-arg",arg+"="+image["reference"],"--file",str(context/dockerfile),"--metadata-file",str(metadata),"--output","type=oci,dest="+str(archive),str(context)]
        result=bounded(argv,context,dict(os.environ),1200)
        require(result.returncode==0 and archive.is_file() and metadata.is_file(),"OCI build failed; no promotion")
        meta=read(metadata);digest=meta.get("containerimage.digest","");require(digest.startswith("sha256:") and len(digest)==71,"OCI image digest absent")
        receipts[component]={"image_digest":digest,"archive_sha256":file_hash(archive),"metadata_sha256":file_hash(metadata),"base_image":image["reference"],"target_runtime":"NOT_RUN"}
    report={"schema":"elite-oci-archives/v403.1","result":"PASS","claim":"exact Linux build output packaged into OCI archives","build_manifest_sha256":file_hash(Path(build)/"artifact-manifest.json"),"components":receipts,"production_authorized":False,"push":"NOT_RUN","cloud":"NOT_RUN"}
    (output/"oci-receipt.json").write_bytes(raw(report));return report

if __name__=="__main__":
    p=argparse.ArgumentParser();p.add_argument("--source",required=True);p.add_argument("--build",required=True);p.add_argument("--output",required=True);p.add_argument("--tool",required=True);p.add_argument("--runner",required=True);p.add_argument("--authorize-candidate-build",action="store_true");a=p.parse_args()
    print(raw(execute(a.source,a.build,a.output,read(a.tool),read(a.runner),authorize=a.authorize_candidate_build)).decode(),end="")
