from __future__ import annotations

import argparse
from datetime import datetime, timezone
import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import re
import subprocess
import sys
import tempfile
import time

from validate_production_admission import CONTROLS, DIGEST, HEX

SCHEMA = "elite-official-tool-run/v1"
ENV_NAME = re.compile(r"^[A-Z][A-Z0-9_]{0,63}$")
MAX_OUTPUT = 10 * 1024 * 1024


def safe_project_path(root: Path, value: object, label: str, *, file: bool) -> Path:
    if not isinstance(value, str) or "\\" in value:
        raise ValueError(f"{label} must be canonical project-relative")
    relative = PurePosixPath(value)
    if relative.is_absolute() or any(part in {"", ".."} for part in relative.parts):
        raise ValueError(f"{label} must be canonical project-relative")
    candidate = root.joinpath(*relative.parts)
    resolved = candidate.resolve(strict=True)
    resolved.relative_to(root)
    if candidate.is_symlink() or (file and not resolved.is_file()) or (not file and not resolved.is_dir()):
        raise ValueError(f"{label} has invalid type or symlink")
    return resolved


def atomic(path: Path, data: bytes) -> None:
    if path.exists():
        raise ValueError(f"output already exists: {path.name}")
    with tempfile.NamedTemporaryFile(dir=path.parent, prefix=f".{path.name}.", delete=False) as handle:
        temporary = Path(handle.name); handle.write(data); handle.flush(); os.fsync(handle.fileno())
    try:
        temporary.replace(path)
    except BaseException:
        temporary.unlink(missing_ok=True); raise


def bounded_process(arguments: list[str], cwd: Path, environment: dict[str, str], timeout: int) -> subprocess.CompletedProcess:
    """Bound each captured stream; supervise only the direct child, never restart.

    The deadline includes pipe EOF after child exit. No communicate(), reader
    threads, inherited stdin or output retry. OS process creation itself is not
    interruptible here. Descendants require a target-owned job/cgroup supervisor.
    """
    process = None
    streams = []
    buffers = [bytearray(), bytearray()]
    failure = None
    completed = None
    deadline = time.monotonic() + timeout
    try:
        flags = subprocess.CREATE_NO_WINDOW if os.name == "nt" else 0
        process = subprocess.Popen(arguments, cwd=cwd, env=environment, shell=False,
                                   stdin=subprocess.DEVNULL, stdout=subprocess.PIPE,
                                   stderr=subprocess.PIPE, bufsize=0, creationflags=flags)
        streams = [process.stdout, process.stderr]
        for stream in streams:
            os.set_blocking(stream.fileno(), False)
        pending = {0, 1}
        while True:
            if time.monotonic() >= deadline:
                failure = "official tool timed out"
                break
            progressed = False
            for index in tuple(pending):
                # At most one bounded read per stream per pass prevents starvation.
                size = min(65536, MAX_OUTPUT - len(buffers[index]) + 1)
                try:
                    chunk = os.read(streams[index].fileno(), size)
                except BlockingIOError:
                    continue
                if not chunk:
                    pending.remove(index)
                    continue
                progressed = True
                if len(buffers[index]) + len(chunk) > MAX_OUTPUT:
                    failure = "official tool output exceeded 10 MiB"
                    break
                buffers[index].extend(chunk)
            if failure:
                break
            code = process.poll()
            if code is not None and not pending:
                completed = subprocess.CompletedProcess(arguments, code, bytes(buffers[0]), bytes(buffers[1]))
                break
            if not progressed:
                time.sleep(min(.005, max(0, deadline - time.monotonic())))
    except (OSError, ValueError, subprocess.SubprocessError):
        failure = "official tool execution unavailable"
    finally:
        # Raw unbuffered handles: closing cannot wait on a buffered reader lock.
        # Close the pipes even when a descendant keeps its inherited write end.
        for stream in streams:
            try:
                stream.close()
            except OSError:
                failure = "official tool cleanup unavailable"
        if process is not None:
            try:
                if process.poll() is None:
                    process.kill()
                process.wait(timeout=5)
            except (OSError, subprocess.SubprocessError):
                failure = "official tool cleanup unavailable"
    if failure:
        # Do not chain TimeoutExpired/OSError: argv and tool output may be private.
        raise ValueError(failure) from None
    if completed is None:
        raise ValueError("official tool execution unavailable") from None
    return completed


def run(profile: object, root: Path, output: Path, authorized: bool) -> dict[str, object]:
    if not authorized:
        raise ValueError("explicit --authorize-execution is required")
    if not isinstance(profile, dict) or set(profile) != {"schema", "project_id", "environment", "release_digest", "control_id", "tool", "arguments", "working_directory", "environment_variable_names", "timeout_seconds", "target"}:
        raise ValueError("profile fields must be exact")
    if profile["schema"] != SCHEMA or profile["environment"] != "production" or profile["control_id"] not in CONTROLS:
        raise ValueError("profile identity is invalid")
    if not isinstance(profile["project_id"], str) or not profile["project_id"].strip() or not DIGEST.fullmatch(str(profile["release_digest"])):
        raise ValueError("project/release identity is invalid")
    tool = profile["tool"]
    if not isinstance(tool, dict) or set(tool) != {"path", "sha256", "name", "version", "source"} or not HEX.fullmatch(str(tool.get("sha256", ""))):
        raise ValueError("tool identity is invalid")
    binary = Path(str(tool["path"])).resolve(strict=True)
    if Path(str(tool["path"])).is_symlink() or not binary.is_file() or hashlib.sha256(binary.read_bytes()).hexdigest() != tool["sha256"]:
        raise ValueError("tool binary SHA-256 mismatch or invalid file")
    arguments = profile["arguments"]
    if not isinstance(arguments, list) or len(arguments) > 128 or any(not isinstance(v, str) or "\x00" in v or len(v) > 4096 for v in arguments):
        raise ValueError("arguments must be a bounded argv array")
    cwd = safe_project_path(root, profile["working_directory"], "working_directory", file=False)
    names = profile["environment_variable_names"]
    if not isinstance(names, list) or len(names) != len(set(names)) or any(not isinstance(v, str) or not ENV_NAME.fullmatch(v) for v in names):
        raise ValueError("environment_variable_names are invalid")
    missing = [name for name in names if name not in os.environ]
    if missing:
        raise ValueError(f"required environment variables are absent: {missing}")
    timeout = profile["timeout_seconds"]
    if type(timeout) is not int or not 1 <= timeout <= 7200:
        raise ValueError("timeout_seconds must be 1..7200")
    base_names = {"PATH", "SYSTEMROOT", "WINDIR", "TEMP", "TMP", "LANG", "LC_ALL", "SSL_CERT_FILE", "SSL_CERT_DIR"}
    environment = {name: value for name, value in os.environ.items() if name in base_names or name in names}
    started = datetime.now(timezone.utc)
    completed = bounded_process([str(binary), *arguments], cwd, environment, timeout)
    stdout = output.with_suffix(".stdout.bin"); stderr = output.with_suffix(".stderr.bin")
    atomic(stdout, completed.stdout); atomic(stderr, completed.stderr)
    proof = lambda path, data: {"path": path.relative_to(root).as_posix(), "bytes": len(data), "sha256": hashlib.sha256(data).hexdigest()}
    arguments_sha256 = hashlib.sha256(json.dumps(arguments, ensure_ascii=False, separators=(",", ":")).encode()).hexdigest()
    receipt = {"schema": "elite-official-tool-execution/v1", "project_id": profile["project_id"], "environment": "production", "release_digest": profile["release_digest"], "control_id": profile["control_id"], "executed_at": started.isoformat().replace("+00:00", "Z"), "exit_code": completed.returncode, "tool": {key: tool[key] for key in ("name", "version", "source", "sha256")}, "arguments_sha256": arguments_sha256, "working_directory": Path(profile["working_directory"]).as_posix(), "target": profile["target"], "environment_variable_names": names, "stdout": proof(stdout, completed.stdout), "stderr": proof(stderr, completed.stderr)}
    atomic(output, (json.dumps(receipt, indent=2, sort_keys=True) + "\n").encode())
    return receipt


def main() -> int:
    parser = argparse.ArgumentParser(); parser.add_argument("--project-root", required=True, type=Path); parser.add_argument("--profile", required=True); parser.add_argument("--output", required=True); parser.add_argument("--authorize-execution", action="store_true"); args = parser.parse_args()
    root = args.project_root.resolve(strict=True); profile_path = safe_project_path(root, args.profile, "profile", file=True)
    output = root.joinpath(*PurePosixPath(args.output).parts).resolve(strict=False); output.parent.resolve(strict=True).relative_to(root)
    receipt = run(json.loads(profile_path.read_text(encoding="utf-8")), root, output, args.authorize_execution)
    print(json.dumps(receipt, indent=2, sort_keys=True)); return 0 if receipt["exit_code"] == 0 else 2


if __name__ == "__main__": sys.exit(main())
