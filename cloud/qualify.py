"""AUTHORED bounded local qualification; records methods and current source hashes."""
from pathlib import Path
import argparse,json,os,subprocess,sys,tempfile,time,shutil
from cloud_control import file_hash,raw,require

def run(root,argv,env=None,timeout=180):
    start=time.monotonic()
    try:q=subprocess.run(argv,cwd=root,env=env,stdin=subprocess.DEVNULL,capture_output=True,timeout=timeout,shell=False)
    except subprocess.TimeoutExpired as e:q=subprocess.CompletedProcess(argv,124,e.stdout or b"",b"Bounded local qualification timeout; no target PASS.")
    return {"argv":argv,"exit_code":q.returncode,"seconds":round(time.monotonic()-start,3),"stdout_sha256":__import__('hashlib').sha256(q.stdout).hexdigest(),"stderr_sha256":__import__('hashlib').sha256(q.stderr).hexdigest(),"stdout":q.stdout.decode(errors="replace")[-4096:],"stderr":q.stderr.decode(errors="replace")[-4096:]}

def qualify(root,node,receipt,dotnet=None):
    root=Path(root).resolve();receipt=Path(receipt);require(not receipt.exists(),"fresh qualification receipt required")
    results=[run(root,[sys.executable,"-B","-m","unittest","test_cloud_control","test_consumer_build","test_tool_admission","test_stack_manifest","test_execute_stack","test_finops_cycle","-q"]),run(root,[node,"--test","test-cloud-run-transport.mjs"])]
    if dotnet:
        require(sys.platform=="linux","certificate test is Linux only")
        with tempfile.TemporaryDirectory(prefix="elite-cert-home-") as home:
            os.chmod(home,0o700);env=dict(os.environ);env.update(HOME=home,DOTNET_CLI_TELEMETRY_OPTOUT="1",DOTNET_GENERATE_ASPNET_CERTIFICATE="false",DOTNET_SYSTEM_GLOBALIZATION_INVARIANT="1")
            results.append(run(root,[dotnet,"run","--project","cert_bridge/Elite.Cloud.CertificateFixture.csproj","--configuration","Release"],env))
    candidates=[p for p in root.iterdir() if p.is_file()]
    candidates += [p for p in (root/"cert_bridge").iterdir() if p.is_file()]
    candidates += [p for p in (root/"cloud_overlay").rglob("*") if p.is_file()]
    sources={p.relative_to(root).as_posix():file_hash(p) for p in candidates}
    runtime={"python_version":sys.version,"python_sha256":file_hash(Path(sys.executable).resolve()),"node_sha256":file_hash(Path(shutil.which(node) or node).resolve())}
    if dotnet:runtime["dotnet_launcher_sha256"]=file_hash(Path(dotnet).resolve())
    report={"schema":"elite-cloud-local-qualification/v403.1","result":"PASS" if all(x["exit_code"]==0 for x in results) else "FAIL","method":"EXECUTED_"+sys.platform.upper()+"_OFFLINE_FIXTURES","runtime":runtime,"sources":sources,"checks":results,"cloud_execution":"NOT_RUN","production_authorized":False,"certificate_fixture":"EXECUTED_LINUX_SYNTHETIC" if dotnet else "NOT_RUN_IN_THIS_RECEIPT"}
    receipt.parent.mkdir(parents=True,exist_ok=True);receipt.write_bytes(raw(report));return report

if __name__=="__main__":
    p=argparse.ArgumentParser();p.add_argument("--node",required=True);p.add_argument("--dotnet");p.add_argument("--receipt",required=True);a=p.parse_args();r=qualify(Path(__file__).parent,a.node,a.receipt,a.dotnet);print(json.dumps({"result":r["result"],"receipt":a.receipt}));sys.exit(0 if r["result"]=="PASS" else 2)
