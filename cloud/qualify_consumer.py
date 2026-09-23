"""AUTHORED independent 120-pack Git fixture; never opens a real franchise.

The source commit must already exist in the local library. This proves separate
checkouts/locks/overlays on the chosen OS; it does not claim remote publication.
"""
from pathlib import Path
import argparse,subprocess,sys
from cloud_control import raw,file_hash,require,read
from seal_consumer import seal
from bootstrap import bootstrap

def qualify(library,commit,work,pwsh,tooling_invariant=False):
    work=Path(work).absolute();require(not work.exists(),"fresh disposable qualification root required")
    for p in (work,*work.parents):require(not p.is_symlink() and not (hasattr(p,"is_junction") and p.is_junction()),"linked qualification root rejected")
    work.mkdir(parents=True);project=work/"fictitious-project";project.mkdir();(project/"QUALIFICATION_ONLY.md").write_text("Synthetic consumer. No business startup, accounts, provider effects or production authority.\n")
    subprocess.run(["git","init","--quiet",str(project)],check=True,stdin=subprocess.DEVNULL,capture_output=True,timeout=30,shell=False)
    locks=project/".elite/library-pin";sealed=seal(library,commit,locks,str(Path(library).resolve()),fixture=True)
    result=bootstrap(project,locks/"library-checkout.lock.json",locks/"composition.json",work/"worker",pwsh,fixture=True,tooling_invariant=tooling_invariant)
    require(result["composition"]["pack_count"]==120 and result["extension_overlays"],"full extension materialization/overlays absent")
    receipt={"schema":"elite-full120-independent-bootstrap/v403.1","result":"PASS","method":"EXECUTED_"+sys.platform.upper()+"_LOCAL_GIT_REPOSITORIES","library_commit":commit,"seal":sealed,"bootstrap":result,"cloud_execution":"NOT_RUN","remote_commit_availability":"NOT_RUN","production_authorized":False,"source_tools":{n:file_hash(Path(__file__).with_name(n)) for n in ["qualify_consumer.py","seal_consumer.py","bootstrap.py","cloud_control.py"]}}
    (work/"qualification-receipt.json").write_bytes(raw(receipt));return receipt

if __name__=="__main__":
    p=argparse.ArgumentParser();p.add_argument("--library",required=True);p.add_argument("--commit",required=True);p.add_argument("--work",required=True);p.add_argument("--pwsh",required=True);p.add_argument("--tooling-invariant",action="store_true");a=p.parse_args()
    result=qualify(a.library,a.commit,a.work,a.pwsh,a.tooling_invariant);print(raw(result).decode(),end="")
