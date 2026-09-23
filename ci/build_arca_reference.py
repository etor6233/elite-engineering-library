"""AUTHORED offline build glue for the existing ARCA Go/UDS/.NET owners."""
from __future__ import annotations

import hashlib
import json
import os
from pathlib import Path
import shutil
import sys
import zipfile

from build_local_reference import command, plain, raw_json, read_json, relative, require, sha


def build_arca(src: Path, dist: Path, work: Path, logs: Path, go: Path,
               environment: dict, input_path: str, input_sha256: str) -> dict:
    cfg = read_json(input_path, input_sha256)
    require(set(cfg) == {"schema_version", "scope", "dotnet_root", "pwsh", "svcutil_archive",
                        "svcutil_packages", "wsdl_cache", "runtime_packages"}, "Exact ARCA build inputs required")
    require(cfg["schema_version"] == 1 and cfg["scope"] == "LOCAL_FIXTURES", "ARCA scope mismatch")
    assets = src / "arca/fiscal"
    sdk = plain(cfg["dotnet_root"])
    sdk_lock = json.loads((assets / "dotnet-sdk.lock.json").read_text())
    require(sdk_lock["schema"] == "elite-arca-dotnet-sdk/v1", "Unknown .NET SDK inventory")
    for row in sdk_lock["files"]:
        path = relative(sdk, row["path"])
        require(path.stat().st_size == row["bytes"] and sha(path) == row["sha256"], "SDK input changed")
    runtime_lock = json.loads((assets / "dotnet-runtime.lock.json").read_text())
    require(runtime_lock["schema"] == "elite-arca-dotnet-runtime/v1", "Unknown .NET runtime projection")
    for row in runtime_lock["files"]:
        path = relative(sdk, row["path"])
        require(path.stat().st_size == row["bytes"] and sha(path) == row["sha256"], "Runtime pin changed")
    require(set(cfg["pwsh"]) == {"path", "sha256"}, "Exact PowerShell pin required")
    pwsh = plain(cfg["pwsh"]["path"])
    require(sha(pwsh) == cfg["pwsh"]["sha256"], "PowerShell pin changed")
    for name in ["svcutil_archive", "svcutil_packages", "wsdl_cache", "runtime_packages"]:
        cfg[name] = str(plain(cfg[name]))
    area = work / "arca"
    require(not area.exists(), "Fresh ARCA build area required")
    area.mkdir()
    generated = src / "arca/fiscal/generated"
    require(not generated.exists(), "Generated output must be absent")
    dotnet = sdk / "dotnet.exe"
    env = dict(environment)
    env.update(DOTNET_ROOT=str(sdk), DOTNET_ROOT_X64=str(sdk), DOTNET_CLI_TELEMETRY_OPTOUT="1",
               DOTNET_SVCUTIL_TELEMETRY_OPTOUT="1", DOTNET_NOLOGO="1", DOTNET_ROLL_FORWARD="Major",
               DOTNET_CLI_HOME=str(area / "home"), NUGET_PACKAGES=str(area / "packages"),
               APPDATA=str(area / "roaming"), LOCALAPPDATA=str(area / "local"),
               TEMP=str(area), TMP=str(area), PATH=str(sdk) + os.pathsep + env["PATH"])
    # NuGet's Windows configuration discovery requires the machine directory
    # locators even when all restore sources are explicitly local and cleared.
    for name in ["ProgramFiles", "ProgramFiles(x86)", "ProgramData", "SystemDrive"]:
        require(bool(os.environ.get(name)), "Windows machine directory missing: " + name)
        env[name] = os.environ[name]
    steps = []
    steps.append(command("arca-tool", [sys.executable, "-X", "utf8", "-B", src / "tools/assemble_arca_svcutil.py",
        "--source-archive", cfg["svcutil_archive"], "--package-directory", cfg["svcutil_packages"],
        "--output", area / "svcutil"], src, env, logs))
    steps.append(command("arca-generate", [pwsh, "-NoLogo", "-NoProfile", "-File", src / "tools/generate-arca-wsfe-client.ps1",
        "-DotnetExecutable", dotnet, "-CacheDirectory", cfg["wsdl_cache"], "-OutputDirectory", generated,
        "-SvcutilRuntimeDirectory", area / "svcutil", "-OfflinePackageDirectory", cfg["runtime_packages"], "-Offline"], src, env, logs))
    artifact = dist / "arca"
    require(not artifact.exists(), "ARCA artifact must be absent")
    artifact.mkdir()
    projects = {"worker": "worker/Elite.Arca.Wsfe.Worker/Elite.Arca.Wsfe.Worker.csproj",
                "fixture": "worker/fixtures/Elite.Arca.Wsfe.Fixture/Elite.Arca.Wsfe.Fixture.csproj"}
    for name, rel in projects.items():
        project = assets / rel
        steps.append(command("arca-" + name + "-restore", [dotnet, "restore", project, "--locked-mode", "--no-cache",
            "--warnaserror", "--configfile", generated / "NuGet.Config", "-p:NuGetAudit=false"], src, env, logs))
        steps.append(command("arca-" + name + "-publish", [dotnet, "publish", project, "--configuration", "Release",
            "--no-restore", "--no-self-contained", "--output", artifact / name, "--warnaserror", "-p:UseAppHost=false",
            "-p:Deterministic=true", "-p:ContinuousIntegrationBuild=true", "-p:DebugType=None",
            "-p:PathMap=" + str(src) + "=/_/source", "-p:NuGetAudit=false"], src, env, logs))
    (artifact / "go").mkdir()
    for name in ["arca-fiscal-worker", "arca-parameter-worker"]:
        steps.append(command(name, [go, "build", "-trimpath", "-buildvcs=false", "-ldflags=-s -w -buildid=",
            "-o", artifact / "go" / (name + ".exe"), "./cmd/" + name], src, env, logs))
    for row in runtime_lock["files"]:
        destination = artifact / "dotnet" / row["path"]
        destination.parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(relative(sdk, row["path"]), destination)
        require(sha(destination) == row["sha256"], "Runtime copy drift")
    steps.append(command("arca-runtime", [artifact / "dotnet/dotnet.exe", "--list-runtimes"], artifact, env, logs))
    runtime_output = (logs / "arca-runtime.log").read_text()
    require("Microsoft.NETCore.App 10.0.11" in runtime_output and "Microsoft.AspNetCore.App 10.0.11" in runtime_output,
            "Portable .NET runtime incomplete")
    require((artifact / "dotnet/shared").as_posix() in runtime_output.replace("\\", "/"),
            "Runtime executable did not resolve its own shipped frameworks")
    notices = []
    for row in json.loads((assets / "nuget-runtime.lock.json").read_text())["packages"]:
        archive = plain(Path(cfg["runtime_packages"]) / row["archive"])
        require(sha(archive) == row["sha256"], "NuGet archive drift")
        with zipfile.ZipFile(archive) as package:
            for legal in row["legal"]:
                data = package.read(legal["member"])
                require(hashlib.sha256(data).hexdigest() == legal["sha256"], "NuGet notice drift")
                dest = artifact / "dependency-notices" / row["name"] / legal["member"]
                dest.parent.mkdir(parents=True, exist_ok=True)
                dest.write_bytes(data)
        if "supplemental_legal" in row:
            legal = row["supplemental_legal"]
            source = relative(assets, legal["input"])
            require(sha(source) == legal["sha256"], "Supplemental MIT notice drift")
            dest = artifact / "dependency-notices" / row["name"] / legal["file"]
            dest.parent.mkdir(parents=True, exist_ok=True)
            shutil.copyfile(source, dest)
        notices.append(row)
    (artifact / "nuget-dependency-notices.json").write_bytes(raw_json({"packages": notices}))
    source = artifact / "generated-source"
    source.mkdir()
    for rel in ["Generated/WsaaReference.cs", "Generated/WsfeV1Reference.cs", "GENERATION_RECEIPT.json",
                "packages.lock.json", "Elite.Arca.Wsfe.Generated.csproj"]:
        path = source / rel
        path.parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(generated / rel, path)
    shutil.copyfile(src / "docs/ARCA_LOCAL_REFERENCE.md", artifact / "START_ARCA.md")
    receipt = {"schema": "elite-arca-reference-artifact/v1", "scope": "LOCAL_FIXTURES",
               "runtime_files": len(runtime_lock["files"]), "runtime_lock_sha256": sha(assets / "dotnet-runtime.lock.json"),
               "generation_receipt_sha256": sha(generated / "GENERATION_RECEIPT.json"),
               "provider_state": "CONDITIONED_USER_CREDENTIALS", "production_admitted": False}
    (artifact / "ARCA_REFERENCE.json").write_bytes(raw_json(receipt))
    return {"state": "PASS", "steps": steps, "receipt": receipt, "fresh_nuget_cache": True}
