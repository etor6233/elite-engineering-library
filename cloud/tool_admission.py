"""AUTHORED binding adapter; exact acquisition is never runtime admission.

The consumer's existing admission workflow supplies G0-G8 owner reports. This
module verifies the reports and runtime tree instead of trusting an ADMITTED
string or just the small gcloud/Docker launcher.
"""
from pathlib import Path
from cloud_control import read,file_hash,contained,require,digest

def validate(tool):
    require(tool["admission"]=="ADMITTED_EXACT_ARTIFACT","tool runtime admission required")
    ref=tool["admission_receipt"];require(set(ref)=={"path","sha256"},"hash-bound owner admission receipt required")
    path=Path(ref["path"]).resolve();require(file_hash(path)==ref["sha256"],"tool admission receipt changed")
    report=read(path)
    require(report["schema"]=="elite-cloud-tool-admission/v403.1" and report["result"]=="PASS" and report["scope"]=="CLOUD_REFERENCE_TOOL_RUNTIME" and report["production_authorized"] is False,"tool admission scope/result required")
    require(report["binary_sha256"]==tool["binary_sha256"]==file_hash(tool["binary"]),"admitted tool binary changed")
    require(report["runtime_files"] and report["runtime_manifest_sha256"]==digest(report["runtime_files"]),"whole tool runtime inventory required")
    root=Path(report["runtime_root"]).resolve()
    observed={p.relative_to(root).as_posix() for p in root.rglob("*") if p.is_file() and not p.is_symlink()}
    require(observed==set(report["runtime_files"]),"runtime inventory omission/addition rejected")
    links={p.relative_to(root).as_posix():p.readlink().as_posix() for p in root.rglob("*") if p.is_symlink()}
    require(links==report.get("runtime_links",{}),"runtime link inventory changed")
    for rel,target in links.items():require((root/rel).resolve().is_relative_to(root),"runtime link escapes admitted tree")
    for rel,h in report["runtime_files"].items():require(file_hash(contained(root,rel))==h,"admitted runtime dependency changed")
    require(not Path(tool["binary"]).is_symlink() and Path(tool["binary"]).resolve().is_relative_to(root),"launcher outside admitted runtime root or alias")
    require(set(report["gates"])=={"G"+str(i) for i in range(9)},"G0-G8 owner reports required")
    for gate,proof in report["gates"].items():
        p=contained(path.parent,proof["path"]);require(file_hash(p)==proof["sha256"],"owner gate evidence changed")
        g=read(p);require(g["gate"]==gate and g["result"]=="PASS" and g["method"] in ("EXECUTED_LOCAL","EXECUTED_TARGET") and g["runtime_manifest_sha256"]==report["runtime_manifest_sha256"],"owner gate not executed for exact runtime")
    return report
