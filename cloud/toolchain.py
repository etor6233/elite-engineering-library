"""AUTHORED exact archive acquisition, not automatic upstream admission.

Allows only the explicit official host + sha256 lock for each tool. Installs
portable artifacts below a worker directory; never modifies host/system PATH.
Every archive member is preflighted before extraction, with symlinks rejected.
"""
from __future__ import annotations
import argparse
import hashlib
import json
from pathlib import Path
import re
import os
import tarfile
import urllib.request
import urllib.parse
import zipfile
import tempfile
from cloud_control import contained, file_hash, raw, read, require

OFFICIAL={"go":{"go.dev","dl.google.com"},"node":{"nodejs.org"},"pwsh":{"github.com","release-assets.githubusercontent.com"},"dotnet":{"builds.dotnet.microsoft.com","download.visualstudio.microsoft.com"},"gcloud":{"dl.google.com"},"docker":{"download.docker.com"},"buildx":{"github.com","release-assets.githubusercontent.com"}}

def install(row, root, archive=None):
    require(row["tool"] in OFFICIAL and row["admission"] in ("ACQUISITION_ONLY","ADMITTED_EXACT_ARTIFACT"), "explicit pinned acquisition admission required")
    require(re.fullmatch(r"[0-9a-f]{64}",row["sha256"]) and row["version"] and row["platform"]=="linux/amd64", "exact version/digest/platform required")
    url=urllib.parse.urlparse(row["url"]); require(url.scheme=="https" and url.hostname in OFFICIAL[row["tool"]] and not url.username and not url.password,"official archive HTTPS required")
    root=Path(root).absolute();require(not root.exists(),"portable tool destination must be absent")
    if archive is None:
        request=urllib.request.Request(row["url"],headers={"User-Agent":"Elite-Library-V403-Qualification"})
        with urllib.request.urlopen(request,timeout=60) as r:
            final=urllib.parse.urlparse(r.url);require(final.scheme=="https" and final.hostname in OFFICIAL[row["tool"]],"redirect outside official artifact hosts")
            data=r.read(512*1024*1024+1);require(len(data)<=512*1024*1024,"tool download budget exceeded")
    else:data=Path(archive).read_bytes()
    require(hashlib.sha256(data).hexdigest()==row["sha256"],"tool archive digest mismatch")
    root.parent.mkdir(parents=True,exist_ok=True); cache=root.parent/(root.name+".archive")
    require(not cache.exists(),"fresh archive cache required");cache.write_bytes(data)
    if row["format"]=="binary":
        require(row["tool"]=="buildx","only fixed buildx raw binary supported")
        root.mkdir();binary=root/"docker-buildx";binary.write_bytes(data);binary.chmod(0o755)
        return {"result":"PASS","claim":"exact official artifact acquisition only","tool":row["tool"],"sha256":row["sha256"],"files":1,"runtime":"NOT_RUN","security_admission":"NOT_INFERRED"}
    with tempfile.TemporaryDirectory(prefix=".elite-case-probe-",dir=root.parent) as probe:
        p=Path(probe);(p/"CaseA").write_bytes(b"probe");case_sensitive=not (p/"casea").exists()
    seen=set();size=0
    def member(name, nbytes):
        nonlocal size
        # Root is absent and every member is preflighted before extraction.
        # Lexical checks avoid O(files * mounted-filesystem-stat) startup cost.
        name=name.rstrip("/")
        require(name and not name.startswith("/") and "\\" not in name and ":" not in name and all(p not in ("", ".", "..") for p in name.split("/")),"unsafe relative archive path")
        path=root/name; key=name if case_sensitive else name.casefold()
        require(key not in seen,"archive duplicate path");seen.add(key);size+=nbytes
        require(len(seen)<=50000 and size<=2*1024*1024*1024,"archive expansion limit")
        return path
    if row["format"]=="tar.gz":
        with tarfile.open(cache,"r:gz") as a:
            entries=[]
            for e in a.getmembers():
                if e.name in (".","./"):
                    require(e.isdir(),"root archive entry must be directory");continue
                while e.name.startswith("./"):e.name=e.name[2:]
                entries.append(e)
            regular=set();directories=set();links=[]
            for e in entries:
                require(e.isfile() or e.isdir() or e.issym(),"archive hardlinks/devices forbidden");member(e.name,e.size)
                if e.isfile():regular.add(e.name.rstrip("/"))
                elif e.isdir():directories.add(e.name.rstrip("/"))
                elif e.issym():links.append(e)
                directories.update(p.as_posix() for p in Path(e.name).parents if p.as_posix()!=".")
            for e in links:
                require(not e.linkname.startswith("/") and "\\" not in e.linkname and ":" not in e.linkname,"absolute archive link forbidden")
                rel=os.path.normpath(str(Path(e.name).parent/e.linkname)).replace("\\","/")
                require(not rel.startswith("../") and rel!=".." and not rel.startswith("/"),"archive link escapes root");require(rel in regular|directories,"link must target an in-archive regular file/directory, not a chain")
                require(not any(x.name.startswith(e.name.rstrip("/")+"/") for x in entries),"archive link parent traversal")
            root.mkdir();a.extractall(root,members=entries,filter="data")
    elif row["format"]=="zip":
        with zipfile.ZipFile(cache) as a:
            for e in a.infolist():
                require((e.external_attr>>16)&0o170000 != 0o120000,"archive symlink forbidden");member(e.filename,e.file_size)
            root.mkdir();a.extractall(root)
    else:raise ValueError("unsupported archive format")
    return {"result":"PASS","claim":"exact portable archive integrity and safe extraction","tool":row["tool"],"version":row["version"],"archive_sha256":row["sha256"],"files":len(seen),"runtime_execution":"NOT_RUN","security_admission":"NOT_INFERRED_FROM_ACQUISITION"}

def main():
    p=argparse.ArgumentParser();p.add_argument("--lock",required=True);p.add_argument("--output",required=True);p.add_argument("--offline-archives-root");a=p.parse_args();lock=read(a.lock)
    require(lock["schema"]=="elite-cloud-toolchain/v403.1","toolchain lock schema")
    require(lock["tools"] and not lock["pending"],"toolchain admission incomplete; do not substitute latest or host tools")
    receipts=[]
    for row in lock["tools"]:receipts.append(install(row,Path(a.output)/row["tool"],contained(Path(a.offline_archives_root),row["local_archive"]) if a.offline_archives_root else None))
    print(raw({"receipts":receipts}).decode(),end="")

if __name__=="__main__":main()
