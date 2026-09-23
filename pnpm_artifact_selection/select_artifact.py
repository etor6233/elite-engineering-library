#!/usr/bin/env python3
"""AUTHORED local projection of a pinned pnpm archive; never executes its bytes."""
from __future__ import annotations
import argparse
import hashlib
import io
import json
import os
from pathlib import Path, PurePosixPath
import re
import shutil
import stat
import sys
import tarfile
import tempfile

POLICY_SHA256 = "93c27e601ccc04bbadb755483186e96d3aca0bad6b37063c668e934e55bc43f6"
RECEIPT = "selection-receipt.json"
REPARSE = 0x400

class SelectionError(ValueError):
    pass

def require(condition, message):
    if not condition:
        raise SelectionError(message)

def digest(data):
    return hashlib.sha256(data).hexdigest()

def json_bytes(value):
    return (json.dumps(value, sort_keys=True, ensure_ascii=False, indent=2) + "\n").encode("utf-8")

def unique(pairs):
    result = {}
    for key, value in pairs:
        require(key not in result, "duplicate JSON key")
        result[key] = value
    return result

def decode(data):
    return json.loads(data.decode("utf-8"), object_pairs_hook=unique)

def no_reparse(path):
    path = Path(os.path.abspath(path))
    for item in [*reversed(path.parents), path]:
        try:
            info = item.lstat()
        except FileNotFoundError:
            continue
        require(not stat.S_ISLNK(info.st_mode) and not (getattr(info, "st_file_attributes", 0) & REPARSE), "reparse point rejected")
    return path

def read_bounded(path, limit):
    path = no_reparse(path)
    with path.open("rb") as stream:
        info = os.fstat(stream.fileno())
        require(stat.S_ISREG(info.st_mode) and info.st_size <= limit, "input size/type rejected")
        data = stream.read(limit + 1)
    require(len(data) <= limit, "input read budget exceeded")
    return data

def safe_name(name):
    require(isinstance(name, str) and name and "\\" not in name, "unsafe member name")
    parts = name.split("/")
    require(not name.startswith("/") and all(p not in ("", ".", "..") for p in parts), "unsafe member path")
    for part in parts:
        require(not re.search(r'[\x00-\x1f<>:"|?*]', part) and not part.endswith((".", " ")), "unsafe Windows member name")
        require(part.split(".")[0].upper() not in {"CON", "PRN", "AUX", "NUL", *[f"COM{i}" for i in range(1, 10)], *[f"LPT{i}" for i in range(1, 10)]}, "reserved Windows member name")
    require(str(PurePosixPath(name)) == name, "noncanonical member name")
    return name

def load_policy():
    raw = read_bounded(Path(__file__).with_name("selection-policy.json"), 262144)
    require(digest(raw) == POLICY_SHA256, "selection policy identity mismatch")
    policy = decode(raw)
    require(policy["schema"] == "elite-pnpm-selection-policy/v1", "policy schema rejected")
    entries = policy["files"]
    require(len(entries) == 455, "policy member count rejected")
    folded = set()
    for item in entries:
        name = safe_name(item["path"])
        require(name.casefold() not in folded, "policy case collision")
        folded.add(name.casefold())
        require(item["selection"] in ("unchanged", "excluded_native_family"), "policy selection rejected")
        require(type(item["bytes"]) is int and 0 <= item["bytes"] <= 22000000, "policy size rejected")
        require(re.fullmatch("[0-9a-f]{64}", item["sha256"]), "policy digest rejected")
    require(sum(x["selection"] == "unchanged" for x in entries) == 442, "selected count rejected")
    return policy

def acquisition_binding(raw, policy):
    value = decode(raw)
    require(isinstance(value, dict), "acquisition receipt must be object")
    for key, expected in policy["receipt_bindings"].items():
        require(type(value.get(key)) is type(expected) and value[key] == expected, "acquisition binding rejected: " + key)
    for key in ("profile_sha256", "approval_sha256", "lock_sha256"):
        require(isinstance(value.get(key), str) and re.fullmatch("[0-9a-f]{64}", value[key]), "acquisition reference missing: " + key)
    return digest(raw)

def make_receipt(policy, acquisition_sha):
    return {"schema": "elite-pnpm-selection-receipt/v1", "policy_sha256": POLICY_SHA256,
            "artifact_sha256": policy["archive_sha256"], "acquisition_receipt_sha256": acquisition_sha,
            "provenance": "ADAPTED_SELECTION_UNCHANGED_FILES",
            "selected_files": sum(x["selection"] == "unchanged" for x in policy["files"]),
            "excluded_files": sum(x["selection"] != "unchanged" for x in policy["files"]), "installed": False, "executed": False,
            "runtime_admitted": False, "redistribution_admitted": False}

def project_payload(payload, entries, destination):
    """Internal parser, tested with synthetic archives; public CLI uses only locked bytes."""
    expected = {"package/" + item["path"]: item for item in entries}
    seen = set()
    with tarfile.open(fileobj=io.BytesIO(payload), mode="r:gz") as archive:
        for member in archive:
            require(member.name in expected and member.name not in seen, "unexpected/duplicate archive member")
            seen.add(member.name)
            item = expected[member.name]
            require(member.isreg() and not member.sparse and member.size == item["bytes"], "member type/size rejected")
            name = safe_name(item["path"])
            stream = archive.extractfile(member)
            require(stream is not None, "member content absent")
            with stream:
                data = stream.read(item["bytes"] + 1)
            require(len(data) == item["bytes"] and digest(data) == item["sha256"], "member digest rejected")
            if item["selection"] == "unchanged":
                out = destination.joinpath(*name.split("/"))
                out.parent.mkdir(parents=True, exist_ok=True)
                with out.open("xb") as writer:
                    writer.write(data)
    require(seen == set(expected), "archive is missing locked members")

def verify(target, policy):
    target = no_reparse(target)
    expected = {"payload/" + x["path"]: x for x in policy["files"] if x["selection"] == "unchanged"}
    allowed_dirs = {"payload"}
    for rel in expected:
        allowed_dirs.update(str(x) for x in PurePosixPath(rel).parents if str(x) != ".")
    found = set()
    def walk(directory):
        with os.scandir(directory) as items:
            for item in items:
                path = Path(item.path)
                rel = path.relative_to(target).as_posix()
                no_reparse(path)
                if item.is_dir(follow_symlinks=False):
                    require(rel in allowed_dirs, "unexpected output directory")
                    walk(path)
                else:
                    require(rel in expected or rel == RECEIPT, "unexpected output file")
                    found.add(rel)
                    if rel in expected:
                        entry = expected[rel]
                        data = read_bounded(path, entry["bytes"])
                        require(len(data) == entry["bytes"] and digest(data) == entry["sha256"], "output member changed")
    walk(target)
    require(found == set(expected) | {RECEIPT}, "output files missing")
    raw = read_bounded(target / RECEIPT, 4096)
    receipt = decode(raw)
    require(isinstance(receipt, dict), "selection receipt must be object")
    source_sha = receipt.get("acquisition_receipt_sha256")
    require(isinstance(source_sha, str) and re.fullmatch("[0-9a-f]{64}", source_sha), "source receipt hash rejected")
    require(raw == json_bytes(make_receipt(policy, source_sha)), "selection receipt changed")
    return receipt

def materialize(artifact, acquisition_receipt, acquisition_sha256, target, policy):
    require(os.name == "nt", "publication supported only on Windows")
    target = no_reparse(target)
    require(target.parent.is_dir() and not os.path.lexists(target), "target must be absent with existing parent")
    raw_receipt = read_bounded(acquisition_receipt, 1048576)
    require(digest(raw_receipt) == acquisition_sha256, "acquisition receipt identity mismatch")
    acquisition_sha = acquisition_binding(raw_receipt, policy)
    payload = read_bounded(artifact, policy["archive_bytes"])
    require(len(payload) == policy["archive_bytes"] and digest(payload) == policy["archive_sha256"], "archive identity mismatch")
    staging = Path(tempfile.mkdtemp(prefix=".pnpm-selection-", dir=target.parent))
    try:
        (staging / "payload").mkdir()
        project_payload(payload, policy["files"], staging / "payload")
        (staging / RECEIPT).write_bytes(json_bytes(make_receipt(policy, acquisition_sha)))
        verify(staging, policy)
        no_reparse(target.parent)
        # Windows rename refuses any occupied destination, including an empty directory.
        os.rename(staging, target)
    finally:
        if staging.exists():
            # Only this invocation's exact newly-created sibling may be removed.
            require(staging.parent == target.parent and staging.name.startswith(".pnpm-selection-"), "cleanup scope rejected")
            shutil.rmtree(staging)
    return verify(target, policy)

def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest="action", required=True)
    create = sub.add_parser("materialize", help="Create a candidate selection, never install or execute")
    for arg in ("artifact", "acquisition-receipt", "acquisition-receipt-sha256", "target"):
        create.add_argument("--" + arg, required=True)
    check = sub.add_parser("verify", help="Verify exact candidate bytes, not runtime admission")
    check.add_argument("--target", required=True)
    args = parser.parse_args(argv)
    try:
        policy = load_policy()
        if args.action == "materialize":
            materialize(args.artifact, args.acquisition_receipt, args.acquisition_receipt_sha256, args.target, policy)
        else:
            verify(args.target, policy)
        print("PNPM_SELECTION_BYTES_PASS runtime_admitted=false redistribution_admitted=false")
        return 0
    except (SelectionError, OSError, UnicodeError, ValueError, KeyError, TypeError, RecursionError, tarfile.TarError) as exc:
        print("PNPM_SELECTION_BLOCKED: " + str(exc), file=sys.stderr)
        return 2

if __name__ == "__main__":
    raise SystemExit(main())
