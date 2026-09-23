"""AUTHORED assembly glue; exact official binaries, declared dependency adaptation.

This local WSDL-only projection excludes the original named-pipe import helpers.
It is not a Microsoft release or a general-purpose replacement for svcutil.
"""
from __future__ import annotations

import argparse
import hashlib
import json
from pathlib import Path
import stat
import zipfile


def digest(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def plain(path: Path) -> Path:
    path = path.absolute()
    if str(path).startswith(("\\\\", "//")):
        raise ValueError("Local paths only")
    for item in (path, *path.parents):
        if item.exists() and item.lstat().st_file_attributes & stat.FILE_ATTRIBUTE_REPARSE_POINT:
            raise ValueError("Reparse input/output rejected")
    return path


def read_pinned(path: Path, sha256: str) -> bytes:
    data = plain(path).read_bytes()
    if digest(data) != sha256:
        raise ValueError("Input hash mismatch: " + path.name)
    return data


def encode(value: object) -> bytes:
    return (json.dumps(value, indent=2, ensure_ascii=False) + "\n").encode("utf-8")


def assemble(source_archive: Path, package_directory: Path, output: Path) -> dict:
    assets = Path(__file__).resolve().parent.parent / "arca/fiscal"
    lock = json.loads((assets / "svcutil-adaptation.lock.json").read_text(encoding="utf-8"))
    if lock["schema"] != "elite-arca-svcutil-adaptation/v1" or lock["scope"] != "PINNED_LOCAL_WSDL_ONLY":
        raise ValueError("Unsupported adaptation")
    output = plain(output)
    package_directory = plain(package_directory)
    source_archive = plain(source_archive)
    for source in (assets, package_directory, source_archive):
        if output == source or output in source.parents or source in output.parents:
            raise ValueError("Input and output must be disjoint")
    if output.exists():
        raise ValueError("Output must be absent")
    read_pinned(source_archive, lock["original"]["sha256"])
    payload: dict[str, bytes] = {}
    with zipfile.ZipFile(source_archive) as archive:
        if len(archive.namelist()) != len(set(archive.namelist())):
            raise ValueError("Duplicate archive names")
        for row in lock["original_files"]:
            data = archive.read(row["member"])
            if digest(data) != row["sha256"]:
                raise ValueError("Original member mismatch")
            payload[row["output"]] = data
    deps = json.loads(payload["dotnet-svcutil.deps.json"])
    for row in lock["patches"]:
        archive_path = package_directory / row["archive"]
        read_pinned(archive_path, row["sha256"])
        with zipfile.ZipFile(archive_path) as archive:
            data = archive.read(row["member"])
        if digest(data) != row["dll_sha256"]:
            raise ValueError("Patch member mismatch")
        payload[row["name"] + ".dll"] = data
        old = row["name"] + "/6.12.1"
        new = row["name"] + "/6.12.5"
        for target in deps["targets"].values():
            node = target.pop(old)
            if list(node["runtime"]) != [row["member"]]:
                raise ValueError("Runtime TFM drift")
            node["runtime"][row["member"]] = row["assembly_identity"]
            target[new] = node
        library = deps["libraries"].pop(old)
        library.update(sha512="sha512-" + row["sha512"], path=row["name"].lower() + "/6.12.5",
                       hashPath=row["name"].lower() + ".6.12.5.nupkg.sha512")
        deps["libraries"][new] = library
    names = {row["name"] for row in lock["patches"]}
    for target in deps["targets"].values():
        for node in target.values():
            for name, version in node.get("dependencies", {}).items():
                if name in names:
                    if version != "6.12.1":
                        raise ValueError("Unexpected original NuGet edge")
                    node["dependencies"][name] = "6.12.5"
    payload["dotnet-svcutil.deps.json"] = encode(deps)
    for row in lock["legal_files"]:
        payload[row["output"]] = read_pinned(assets / row["input"], row["sha256"])
    payload["ADAPTATION.json"] = encode({k: lock[k] for k in ("schema", "scope", "original", "patches", "excluded_original_files", "notice")})
    inventory = [{"path": name, "sha256": digest(data), "bytes": len(data)} for name, data in sorted(payload.items())]
    receipt = {"schema": "elite-arca-svcutil-assembly/v1", "scope": lock["scope"],
               "provenance": "ADAPTED_RUNTIME_COMPOSITION", "files": inventory,
               "inventory_sha256": digest(encode(inventory)), "production_admitted": False}
    if receipt["inventory_sha256"] != lock["expected_inventory_sha256"]:
        raise ValueError("Adapted output differs from the admitted inventory")
    # Every input and output is checked before publishing a new directory.
    output.mkdir(parents=True, exist_ok=False)
    for name, data in payload.items():
        path = output / name
        if path.is_absolute() and not path.is_relative_to(output) or ".." in Path(name).parts:
            raise ValueError("Output escape")
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_bytes(data)
    (output / "ASSEMBLY_RECEIPT.json").write_bytes(encode(receipt))
    return receipt


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--source-archive", type=Path, required=True)
    parser.add_argument("--package-directory", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    receipt = assemble(args.source_archive, args.package_directory, args.output)
    print(json.dumps({"state": "EXACT_ADAPTED_RUNTIME", "files": len(receipt["files"]), "inventory_sha256": receipt["inventory_sha256"]}))


if __name__ == "__main__":
    main()
