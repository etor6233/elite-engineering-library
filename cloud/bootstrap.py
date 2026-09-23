"""AUTHORED two-repository bootstrap; delegates composition and bridge unchanged.

The caller runs from the project repository. All paths are worker-relative.
Library is a distinct pinned checkout without write credentials. No source
upgrade, cloud authentication, provisioning, or live application startup.
"""
from __future__ import annotations
import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess
import sys

from cloud_control import contained, digest, file_hash, raw, read, require, validate_profile

def run(argv, cwd, tooling_invariant=False):
    env=dict(os.environ)
    if tooling_invariant:env["DOTNET_SYSTEM_GLOBALIZATION_INVARIANT"]="1"
    q = subprocess.run(list(map(str, argv)), cwd=cwd, env=env, stdin=subprocess.DEVNULL, stdout=subprocess.PIPE, stderr=subprocess.PIPE, timeout=300, shell=False)
    require(q.returncode == 0, "bootstrap command failed: " + Path(str(argv[0])).name + " (exit " + str(q.returncode) + ")")
    return q.stdout.decode("utf-8", errors="replace").strip()

def checkout(lock, destination, fixture=False):
    require(lock["schema"] == "elite-library-checkout/v403.1", "checkout schema")
    require(re.fullmatch(r"[0-9a-f]{40}", lock["commit"]), "exact Git commit required")
    if not fixture:
        require(lock["repository"] == "https://github.com/etor6233/elite-engineering-library.git", "approved read-only library remote required")
    else:
        require(Path(lock["repository"]).is_dir(), "fixture must be a local repository")
    p=Path(destination).absolute()
    require(not p.exists(), "library destination must be absent")
    p.parent.mkdir(parents=True, exist_ok=True)
    run(["git", "-c", "core.autocrlf=false", "clone", "--no-checkout", "--no-hardlinks", lock["repository"], str(p)], p.parent)
    run(["git", "config", "core.autocrlf", "false"], p)
    run(["git", "checkout", "--detach", lock["commit"]], p)
    require(run(["git", "rev-parse", "HEAD"], p) == lock["commit"], "checkout revision mismatch")
    # Prevent an accidental ordinary push. This is not a hostile-agent sandbox.
    run(["git", "remote", "set-url", "--push", "origin", "DISABLED-LIBRARY-READONLY"], p)
    for name, expected in lock["files"].items():
        require(file_hash(contained(p, name)) == expected, "checkout locked file mismatch: " + name)
    return {"result":"PASS", "claim":"distinct pinned repository checkout", "commit":lock["commit"], "files_checked":len(lock["files"]), "method":"LOCAL_GIT_FIXTURE" if fixture else "REMOTE_GIT_CHECKOUT", "production_authorized":False}

def bootstrap(project, lockfile, profilefile, worker, pwsh="pwsh", fixture=False, checkout_only=False, tooling_invariant=False):
    project=Path(project).resolve(); worker=Path(worker).absolute()
    for part in (worker,*worker.parents):
        require(not part.is_symlink() and not (hasattr(part,"is_junction") and part.is_junction()),"linked worker path rejected")
    worker=worker.resolve()
    require((project / ".git").exists(), "project must have its own Git repository")
    require(project != worker and project not in worker.parents, "worker must be outside project root")
    require(not worker.exists(), "worker destination must be absent")
    worker.mkdir(parents=True)
    library=worker/"library"; result=checkout(read(lockfile), library, fixture)
    if checkout_only:
        return result
    result["composition"]=validate_profile(read(profilefile),library)
    kit=worker/"compositor"
    run([pwsh,"-NoLogo","-NoProfile","-File",library/"materialize_markdown_pack.ps1","-PackFile",library/"implementation_packs/MARKDOWN_COMPOSITOR_CORE.md","-Destination",kit],worker,tooling_invariant)
    materialized=worker/"composition"
    run([pwsh,"-NoLogo","-NoProfile","-File",kit/"tools/compose-markdown-project.ps1","-PlanFile",library/read(profilefile)["selection"]["path"],"-LibraryRoot",library,"-Destination",materialized],worker,tooling_invariant)
    if read(profilefile).get("extension_policy")=="BASE116_PLUS_EXACT_FOUR_V403":
        run([sys.executable,"-B",materialized/"experience_overlay/apply_overlay.py","--target",materialized,"--report",worker/"ui-overlay.json"],worker)
        run([sys.executable,"-B",materialized/"cloud/apply_iam_overlay.py","--target",materialized,"--report",worker/"iam-overlay.json"],worker)
        result["extension_overlays"]={"ui_sha256":file_hash(worker/"ui-overlay.json"),"iam_sha256":file_hash(worker/"iam-overlay.json"),"scope":"worker composition; project source untouched"}
    # Current project already owns its code: never overwrite it with a new scaffold.
    # Qualification composes into the separate worker tree and retains its receipt.
    connection=worker/"connection-tools"
    run([pwsh,"-NoLogo","-NoProfile","-File",library/"materialize_markdown_pack.ps1","-PackFile",library/"implementation_packs/PROJECT_OPERATING_CONNECTION.md","-Destination",connection],worker,tooling_invariant)
    installer=connection/"operating_connection/INSTALL_PROJECT_OPERATING_BRIDGE.ps1"
    result["next_command"]=[str(pwsh),"-NoProfile","-File",str(installer),"-LibraryRoot",str(library),"-ProjectRoot",str(project)]
    result["bridge_execution"]="NOT_RUN; existing connection 18/18 reused; native Linux bridge qualification remains separate"
    result["worker_relative_layout"]={"library":"library","compositor":"compositor","composition":"composition","connection_tools":"connection-tools"}
    result["composition_receipt_sha256"]=file_hash(materialized/"MATERIALIZATION_RECORD.md")
    result["source_integrity_after"]=validate_profile(read(profilefile),library)
    result["tooling_globalization"]="INVARIANT_EXPLICIT_COMPOSITOR_ONLY" if tooling_invariant else "HOST_DEFAULT"
    return result

def main():
    p=argparse.ArgumentParser();p.add_argument("--project",required=True);p.add_argument("--lock",required=True);p.add_argument("--profile",required=True);p.add_argument("--worker",required=True);p.add_argument("--pwsh",default="pwsh");p.add_argument("--receipt",required=True);p.add_argument("--checkout-only",action="store_true");p.add_argument("--tooling-invariant",action="store_true")
    a=p.parse_args(); out=Path(a.receipt);require(not out.exists(),"fresh receipt required")
    result=bootstrap(a.project,a.lock,a.profile,a.worker,a.pwsh,checkout_only=a.checkout_only,tooling_invariant=a.tooling_invariant)
    out.parent.mkdir(parents=True,exist_ok=True);out.write_bytes(raw(result));print(raw(result).decode(),end="")

if __name__=="__main__":main()
