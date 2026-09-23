"""AUTHORED offline setup for the admitted Windows CPython3.14 wheel graph.

Run from the composed project with an already acquired exact wheel directory.
The absent destination is created once; failures preserve logs, never delete it.
No account or credential is read, and no provider API is contacted.
"""
from __future__ import annotations

import argparse
import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import platform
import subprocess
import sys
from urllib.parse import urlsplit


def require(ok: bool, message: str) -> None:
    if not ok:
        raise ValueError(message)


def digest(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def setup(args: argparse.Namespace) -> dict:
    root = Path(__file__).resolve().parent
    require(sys.platform == "win32" and sys.version_info[:2] == (3, 14)
            and platform.machine().lower() in {"amd64", "x86_64"}, "fixed ABI required")
    require(digest(Path(sys.executable)) == args.python_sha256, "base Python pin mismatch")
    lock_path = root / "connected-runtime.lock.json"
    require(digest(lock_path) == args.lock_sha256, "runtime lock mismatch")
    lock = json.loads(lock_path.read_bytes())
    require(lock["schema"] == "elite-merchant-installed-source-lock/v1"
            and len(lock["wheels"]) == 20, "exact wheel graph required")
    require(args.destination.is_absolute() and not args.destination.exists(),
            "absolute absent runtime destination required")
    require(args.wheels.is_absolute() and args.wheels.is_dir(), "absolute wheel cache required")
    # Validate every artifact before any environment is created or pip is run.
    lines = []
    for row in lock["wheels"]:
        url = urlsplit(row["url"])
        require(url.scheme == "https" and url.netloc == "files.pythonhosted.org"
                and not url.query and not url.fragment, "official wheel origin required")
        name = PurePosixPath(url.path).name
        require(name.endswith(".whl") and "\\" not in name, "wheel filename required")
        wheel = args.wheels / name
        require(wheel.is_file() and wheel.stat().st_size == row["bytes"]
                and digest(wheel) == row["sha256"], "wheel pin mismatch: " + name)
        lines.append(row["distribution"] + " @ " + wheel.as_uri() + "#sha256=" + row["sha256"])
    args.destination.mkdir()
    env = {k: v for k, v in os.environ.items()
           if not k.startswith(("PIP_", "PYTHON", "GOOGLE_", "GCLOUD_", "CLOUDSDK_"))}
    env.update(PIP_CONFIG_FILE=os.devnull, PIP_DISABLE_PIP_VERSION_CHECK="1",
               PYTHONDONTWRITEBYTECODE="1")
    requirements = args.destination / "requirements-offline.lock"
    requirements.write_text("\n".join(lines) + "\n", encoding="utf-8", newline="\n")
    rows = []
    result = {"state": "FAIL", "live_account_access": False, "network_required": False,
              "base_python_sha256": args.python_sha256, "lock_sha256": args.lock_sha256}

    def run(name: str, command: list) -> None:
        path = args.destination / (name + ".log")
        with path.open("xb") as output:
            process = subprocess.run(list(map(str, command)), env=env, stdin=subprocess.DEVNULL,
                                     stdout=output, stderr=subprocess.STDOUT, timeout=180,
                                     creationflags=subprocess.CREATE_NO_WINDOW)
        rows.append({"step": name, "exit_code": process.returncode, "sha256": digest(path)})
        require(process.returncode == 0, "setup failed: " + name)

    try:
        runtime = args.destination / "venv"
        run("venv", [sys.executable, "-I", "-B", "-m", "venv", "--copies", runtime])
        python = runtime / "Scripts/python.exe"
        run("install", [python, "-I", "-B", "-m", "pip", "install", "--no-index", "--no-deps",
                        "--require-hashes", "-r", requirements])
        run("pip-check", [python, "-I", "-B", "-m", "pip", "check"])
        site = runtime / "Lib/site-packages"
        for name, sha in lock["files"].items():
            parts = PurePosixPath(name)
            require(not parts.is_absolute() and ".." not in parts.parts
                    and "\\" not in name and ":" not in name, "invalid installed path")
            require(digest(site / name) == sha, "installed source differs: " + name)
        script = root / "connected_worker.py"
        config = {"schema": "elite-merchant-runtime/v1", "python": str(python),
                  "python_sha256": digest(python), "script": str(script),
                  "script_sha256": digest(script), "owner_sha256": digest(root / "sync_product.py"),
                  "runtime_lock_sha256": args.lock_sha256,
                  "mode": "CREDENTIALS", "fixture_origin": ""}
        config_path = args.destination / "merchant-runtime.json"
        config_path.write_text(json.dumps(config, indent=2) + "\n", encoding="utf-8", newline="\n")
        result.update(state="PASS", installed_files=len(lock["files"]), wheels=len(lock["wheels"]),
                      runtime_file=str(config_path), runtime_sha256=digest(config_path),
                      condition="future user credential supplied through GOOGLE_APPLICATION_CREDENTIALS")
        return result
    finally:
        result["receipts"] = rows
        (args.destination / "setup-receipt.json").write_text(
            json.dumps(result, indent=2) + "\n", encoding="utf-8", newline="\n")


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--destination", type=Path, required=True)
    parser.add_argument("--wheels", type=Path, required=True)
    parser.add_argument("--python-sha256", required=True)
    parser.add_argument("--lock-sha256", required=True)
    args = parser.parse_args()
    print(json.dumps(setup(args), sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
