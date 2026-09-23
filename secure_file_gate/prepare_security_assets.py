from __future__ import annotations

import argparse
import hashlib
import json
import os
import shutil
import subprocess
import tempfile
from datetime import datetime, timezone
from pathlib import Path
from typing import Sequence


FRESHCLAM_SHA256 = "4032fdd45184333d7c963eec87664a8339c6f31dc5ae94bb5008fe5e88feaf1f"
CLAMSCAN_SHA256 = "f368912f61ddea8b302acf9114891ec57f8a45b0de77f076b1d19c960bc6f7cf"
YARA_SHA256 = "d3e648651f3eaa5e833faeb394fc7dc7dfb583a31dd5fc572ab6cae7ac41f3fd"


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def require_exact(path: Path, expected: str, label: str) -> None:
    if path.is_symlink() or not path.is_file() or sha256_file(path) != expected:
        raise RuntimeError(f"{label} is not the admitted official binary")


def run(args: Sequence[str], timeout: int) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        list(args), check=False, capture_output=True, text=True, encoding="utf-8", errors="replace",
        timeout=timeout, creationflags=getattr(subprocess, "CREATE_NO_WINDOW", 0),
    )


def write_atomic(path: Path, payload: dict[str, object]) -> None:
    temporary = path.with_name(path.name + ".tmp")
    temporary.write_text(json.dumps(payload, indent=2, sort_keys=True) + "\n", encoding="utf-8", newline="\n")
    os.replace(temporary, path)


def update_clamav(freshclam: Path, clamscan: Path, database: Path, max_age_days: int, timeout: int) -> None:
    require_exact(freshclam, FRESHCLAM_SHA256, "freshclam")
    require_exact(clamscan, CLAMSCAN_SHA256, "clamscan")
    if database.is_symlink():
        raise RuntimeError("database directory cannot be a symlink")
    database.mkdir(parents=True, exist_ok=True)
    config_root = Path(tempfile.mkdtemp(prefix="elite-freshclam-"))
    config = config_root / "freshclam.conf"
    try:
        escaped = str(database.resolve()).replace("\\", "\\\\")
        config.write_text(
            "DatabaseMirror database.clamav.net\n"
            f"DatabaseDirectory {escaped}\n"
            "TestDatabases yes\nCompressLocalDatabase no\nBytecode yes\n",
            encoding="utf-8", newline="\n",
        )
        updated = run([str(freshclam), f"--config-file={config}", "--quiet", "--stdout"], timeout)
        version = run([str(clamscan), f"--database={database}", "--version"], timeout)
        probe = config_root / "benign-probe.txt"
        probe.write_text("Elite ClamAV database verification probe.\n", encoding="utf-8", newline="\n")
        scan = run([
            str(clamscan), f"--database={database}", "--official-db-only=yes",
            f"--fail-if-cvd-older-than={max_age_days}", "--stdout", "--no-summary", "--infected",
            "--follow-file-symlinks=0", "--follow-dir-symlinks=0", str(probe),
        ], timeout)
        if version.returncode != 0 or not version.stdout.strip().startswith("ClamAV 1.5.4/") or scan.returncode != 0:
            raise RuntimeError("freshclam result was not proven by a fresh official database scan")
        inventory = []
        for path in sorted(database.iterdir(), key=lambda item: item.name):
            if path.is_file() and path.suffix.lower() in {".cvd", ".cld", ".sign"}:
                inventory.append({"name": path.name, "bytes": path.stat().st_size, "sha256": sha256_file(path)})
        names = {item["name"] for item in inventory}
        if not all(any(f"{stem}.{suffix}" in names for suffix in ("cvd", "cld")) for stem in ("main", "daily", "bytecode")):
            raise RuntimeError("required official ClamAV databases are missing")
        write_atomic(database / "clamav-database-receipt.json", {
            "schema": "elite-clamav-database-receipt/v1",
            "verified_at": datetime.now(timezone.utc).isoformat().replace("+00:00", "Z"),
            "freshclam_exit_code": updated.returncode,
            "validation": "clamscan version plus official-db-only freshness probe",
            "clamav_version": version.stdout.strip(),
            "max_age_days": max_age_days,
            "files": inventory,
        })
    finally:
        shutil.rmtree(config_root, ignore_errors=True)


def compile_yara(yara: Path, source: Path, output: Path, timeout: int) -> None:
    require_exact(yara, YARA_SHA256, "YARA-X")
    if source.is_symlink() or not (source.is_file() or source.is_dir()):
        raise RuntimeError("YARA source must be a regular file or directory")
    if output.exists() or output.is_symlink():
        raise FileExistsError(f"output already exists: {output}")
    output.parent.mkdir(parents=True, exist_ok=True)
    if output.parent.is_symlink() or not output.parent.is_dir():
        raise RuntimeError("YARA output parent must be a non-symlink directory")
    staging = output.with_name(output.name + ".tmp")
    receipt = output.with_suffix(output.suffix + ".receipt.json")
    try:
        compiled = run([str(yara), "compile", "--output", str(staging), str(source)], timeout)
        if compiled.returncode != 0 or not staging.is_file():
            raise RuntimeError("YARA-X compilation failed")
        if source.is_file():
            source_files = [source]
            probe_target = source
            source_root = source.parent
        else:
            source_files = sorted(
                [path for path in source.rglob("*") if path.is_file() and path.suffix.lower() in {".yar", ".yara"}],
                key=lambda path: path.as_posix(),
            )
            if not source_files:
                raise RuntimeError("YARA source directory has no .yar or .yara files")
            probe_target = source_files[0]
            source_root = source
        probe = run([str(yara), "scan", "--compiled-rules", "--output-format=json", "--timeout", "5", str(staging), str(probe_target)], timeout)
        if probe.returncode != 0:
            raise RuntimeError("compiled YARA-X rules did not pass a load/scan probe")
        payload = json.loads(probe.stdout)
        if payload.get("version") != "1.20.0" or not isinstance(payload.get("matches"), list):
            raise RuntimeError("YARA-X probe returned an unexpected schema")
        os.replace(staging, output)
        source_inventory = [
            {
                "path": path.relative_to(source_root).as_posix(),
                "bytes": path.stat().st_size,
                "sha256": sha256_file(path),
            }
            for path in source_files
        ]
        write_atomic(receipt, {
            "schema": "elite-yara-x-rules-receipt/v1",
            "compiled_at": datetime.now(timezone.utc).isoformat().replace("+00:00", "Z"),
            "yara_x_version": "1.20.0",
            "compiled_rules_sha256": sha256_file(output),
            "source_files": source_inventory,
        })
    finally:
        staging.unlink(missing_ok=True)


def main() -> int:
    parser = argparse.ArgumentParser()
    sub = parser.add_subparsers(dest="command", required=True)
    clam = sub.add_parser("update-clamav")
    clam.add_argument("--freshclam", required=True, type=Path)
    clam.add_argument("--clamscan", required=True, type=Path)
    clam.add_argument("--database", required=True, type=Path)
    clam.add_argument("--max-age-days", type=int, default=2)
    clam.add_argument("--timeout", type=int, default=900)
    yara = sub.add_parser("compile-yara")
    yara.add_argument("--yara", required=True, type=Path)
    yara.add_argument("--source", required=True, type=Path)
    yara.add_argument("--output", required=True, type=Path)
    yara.add_argument("--timeout", type=int, default=120)
    args = parser.parse_args()
    if args.command == "update-clamav":
        update_clamav(args.freshclam, args.clamscan, args.database, args.max_age_days, args.timeout)
        print("CLAMAV_DATABASE_PREPARED")
    else:
        compile_yara(args.yara, args.source, args.output, args.timeout)
        print("YARA_X_RULES_PREPARED")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
