"""AUTHORED worker-relative input binding, no credential or business decisions.

The existing official acquisition owners supply the admitted adapted pnpm
projection, verified offline dependency store and ARCA generation record. This
adapter never silently substitutes the vulnerable raw upstream pnpm payload.
"""
from pathlib import Path
import argparse
from cloud_control import file_hash,raw,read,require,contained

def prepare(source,tool_root,projection,store,generation,output):
    source=Path(source).resolve();tool_root=Path(tool_root).resolve();projection=Path(projection).resolve();store=Path(store).resolve();output=Path(output).absolute()
    require(not output.exists() and store.is_dir(),"fresh inputs destination and verified offline store required")
    policy=read(source/"pnpm_artifact_selection/local-runtime-policy.json")
    for rel,row in policy["files"].items():require(file_hash(contained(projection,rel))==row["sha256"],"adapted pnpm projection mismatch")
    executables={"go":tool_root/"go/go/bin/go","node":tool_root/"node/node-v24.20.0-linux-x64/bin/node","dotnet":tool_root/"dotnet/dotnet","pnpm":projection/"payload/bin/pnpm.cjs"}
    # Find only the exact launcher declared by the current projection policy.
    if not executables["pnpm"].is_file():
        candidates=[projection/x for x in policy["files"] if x.endswith("bin/pnpm.cjs")]
        require(len(candidates)==1,"unique exact adapted pnpm launcher required");executables["pnpm"]=candidates[0]
    tools={name:{"path":str(path),"sha256":file_hash(path),"platform":"linux/amd64"} for name,path in executables.items()};tools["pnpm"]["projection_root"]=str(projection)
    inputs={"scope":"LIBRARY_INFRASTRUCTURE","target":"linux/amd64","source_files":{f.relative_to(source).as_posix():file_hash(f) for f in source.rglob("*") if f.is_file() and not {".git","node_modules",".next"}.intersection(f.parts)},"quality_runner_sha256":file_hash(source/"ci/run_quality_gates.py"),"offline_pnpm_store":str(store),"generation_steps":generation}
    output.mkdir(parents=True);(output/"tools.json").write_bytes(raw(tools));(output/"inputs.json").write_bytes(raw(inputs));return {"result":"PASS","claim":"exact build input binding","runtime_admission":"NOT_INFERRED","files":len(inputs["source_files"])}

if __name__=="__main__":
    p=argparse.ArgumentParser();p.add_argument("--source",required=True);p.add_argument("--tool-root",required=True);p.add_argument("--pnpm-projection",required=True);p.add_argument("--offline-store",required=True);p.add_argument("--generation-record",required=True);p.add_argument("--output",required=True);a=p.parse_args()
    print(raw(prepare(a.source,a.tool_root,a.pnpm_projection,a.offline_store,read(a.generation_record)["generation_steps"],a.output)).decode(),end="")
