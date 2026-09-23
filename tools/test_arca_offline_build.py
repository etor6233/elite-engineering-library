"""AUTHORED negative tests for the executable source/dependency trust boundary."""
import argparse
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys

from assemble_arca_svcutil import assemble


def run(target: Path, inputs: Path, runtime: Path) -> None:
    target.mkdir()
    cfg = json.loads(inputs.read_text())
    script = Path(__file__).with_name("generate-arca-wsfe-client.ps1")
    results = []
    for name in ["tool-drift", "forged-receipt", "extra-tool-file", "missing-tool-file",
                 "feed-drift", "wsdl-drift", "offline-required"]:
        area = target / name
        area.mkdir()
        tool = runtime
        feed = Path(cfg["runtime_packages"])
        cache = Path(cfg["wsdl_cache"])
        if name in {"tool-drift", "forged-receipt", "extra-tool-file", "missing-tool-file"}:
            tool = area / "tool"
            shutil.copytree(runtime, tool)
            if name in {"tool-drift", "forged-receipt"}:
                (tool / "NuGet.Common.dll").write_bytes(b"deliberately invalid input")
            if name == "forged-receipt":
                import hashlib
                receipt = json.loads((tool / "ASSEMBLY_RECEIPT.json").read_text())
                for row in receipt["files"]:
                    if row["path"] == "NuGet.Common.dll":
                        row["sha256"] = hashlib.sha256((tool / row["path"]).read_bytes()).hexdigest()
                (tool / "ASSEMBLY_RECEIPT.json").write_text(json.dumps(receipt))
            if name == "extra-tool-file":
                (tool / "extra.dll").write_bytes(b"untrusted")
            if name == "missing-tool-file":
                # An absent inventory entry is simulated by renaming our own fixture.
                (tool / "NuGet.Common.dll").rename(tool / "different-name.dll")
        if name == "feed-drift":
            feed = area / "feed"
            shutil.copytree(cfg["runtime_packages"], feed)
            next(feed.glob("*.nupkg")).write_bytes(b"untrusted archive")
        if name == "wsdl-drift":
            cache = area / "cache"
            shutil.copytree(cfg["wsdl_cache"], cache)
            (cache / "wsaa-homologation.wsdl").write_text("<different/>")
        env = os.environ.copy()
        env.update(DOTNET_ROOT=cfg["dotnet_root"], DOTNET_ROOT_X64=cfg["dotnet_root"],
                   DOTNET_CLI_TELEMETRY_OPTOUT="1", DOTNET_SVCUTIL_TELEMETRY_OPTOUT="1")
        args = [cfg["pwsh"]["path"], "-NoLogo", "-NoProfile", "-File", str(script),
                "-DotnetExecutable", str(Path(cfg["dotnet_root"]) / "dotnet.exe"),
                "-CacheDirectory", str(cache), "-OutputDirectory", str(area / "output"),
                "-SvcutilRuntimeDirectory", str(tool), "-OfflinePackageDirectory", str(feed)]
        if name != "offline-required":
            args.append("-Offline")
        with (area / "rejection.log").open("xb") as log:
            result = subprocess.run(args, env=env, stdout=log, stderr=subprocess.STDOUT,
                                    timeout=30, creationflags=subprocess.CREATE_NO_WINDOW)
        text = (area / "rejection.log").read_text(errors="replace")
        expected = {"tool-drift": "tool member hash mismatch", "forged-receipt": "tool receipt file drift",
                    "extra-tool-file": "tool contains undeclared files", "missing-tool-file": "tool member hash mismatch",
                    "feed-drift": "offline package pin mismatch", "wsdl-drift": "WSDL identity mismatch",
                    "offline-required": "supports only pinned offline homologation"}[name]
        assert result.returncode != 0 and expected in text, (name, text)
        assert not (area / "output/GENERATION_RECEIPT.json").exists()
        assert not (area / "output/Generated").exists()
        results.append({"case": name, "state": "PASS_REJECTED"})
    invalid = target / "invalid-original.nupkg"
    invalid.write_bytes(b"invalid official archive")
    output = target / "unpublished"
    try:
        assemble(invalid, Path(cfg["svcutil_packages"]), output)
        raise AssertionError("Invalid original archive accepted")
    except ValueError as error:
        assert "Input hash mismatch" in str(error)
    assert not output.exists()
    results.append({"case": "original-archive-drift", "state": "PASS_REJECTED"})
    result = {"state": "PASS_OFFLINE_TRUST_CONTRACTS", "cases": results,
              "actual_generation_claim": "Separate two independent fresh-cache generation receipts required"}
    (target / "result.json").write_text(json.dumps(result, indent=2) + "\n")
    print(json.dumps(result))


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--target", type=Path, required=True)
    parser.add_argument("--inputs", type=Path, required=True)
    parser.add_argument("--runtime", type=Path, required=True)
    args = parser.parse_args()
    run(args.target, args.inputs, args.runtime)
