# Secure Local File Ingestion Gate

## 1. Metadata

```yaml
pack_id: "SECURE-LOCAL-FILE-INGESTION-GATE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa una puerta local fail-closed sobre Google Magika CLI 1.1.0, Cisco Talos ClamAV 1.5.4 y VirusTotal/Google YARA-X 1.20.0, con binarios fijados, bases frescas verificadas, ruleset compilado y recibo atómico sin autorizar almacenamiento de negocio."
stacks: ["CPython 3.12 stdlib", "Magika CLI 1.1.0", "ClamAV 1.5.4", "YARA-X 1.20.0", "Windows x86-64"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "MICROSOFT-MARKITDOWN-LOCAL-RUNTIME 0.1.x", "official document runtime packs"]
incompatible_with: ["unapproved content types", "unreviewed YARA rules", "stale ClamAV databases", "archive extraction", "automatic business storage", "non-Windows binary hash reuse"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0 AND GPL-2.0-only AND BSD-3-Clause"
upstream_sources: ["https://github.com/google/magika/tree/5e2f437fb7b7452368c8c1fa9354858f5487a5c4", "https://github.com/Cisco-Talos/clamav/tree/fa59fca15872bb8a914ba4c68188bcc8a502cbdf", "https://github.com/VirusTotal/yara-x/tree/60ad06971467029e77967e59d580cbbe85a1474d"]
verified_at: "2026-08-26"
```

## 2. Applicability

Use como primera puerta ejecutable para un archivo local ya colocado en cuarentena. Cubre límite de bytes, antimalware local con bases oficiales y frescas, tipo real por contenido, reglas adicionales gobernadas y evidencia atómica. La política nace bloqueada y el agente debe obtener aprobación, tipos/MIME/scores, límites y ruleset exacto antes de usarla.

No es extractor documental, sandbox por sí solo, suite endpoint, garantía universal de seguridad ni autorización para persistir hechos de negocio. Los seis archivos materializados son coordinación/configuración/tests `AUTHORED`; los motores siguen siendo los binarios oficiales y conservan licencias Apache-2.0, GPL-2.0-only y BSD-3-Clause.

## 3. Architecture contract

Orden: validar política/rutas/hashes → ClamAV `official-db-only` con edad y límites → Magika JSON con label/MIME/score aprobados → YARA-X JSON con reglas compiladas y hash exacto → receipt nuevo y atómico. Cualquier detección, base vieja, timeout, error, JSON divergente, binario/ruleset alterado o condición no aprobada rechaza cerrado.

`freshclam` no se acepta sólo por exit code: el preparador obliga a que `clamscan` cargue las bases, aplique edad máxima y escanee un probe benigno antes de emitir inventario hash. El despliegue todavía debe ejecutar el worker aislado, sin egress durante el scan, con CPU/memoria/tiempo, storage de cuarentena/clean separado, métricas, recall e incident owner.

## 4. Exact file manifest

```text
CREATE secure_file_gate/security-policy.template.json
CREATE secure_file_gate/secure_local_file_gate.py
CREATE secure_file_gate/prepare_security_assets.py
CREATE secure_file_gate/test_secure_local_file_gate.py
CREATE secure_file_gate/test_prepare_security_assets.py
CREATE secure_file_gate/README.md
```

## 5. Materialization blocks

### FILE: `secure_file_gate/security-policy.template.json`
```yaml
block_id: "SECURE-LOCAL-FILE-INGESTION-GATE:policy:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed project policy template for exact official engines"
license: "LicenseRef-Workspace-Owner"
sha256: "9722a3c5a3780c407c3020f2720f44664c1a73b4412663cbb8025e352023653d"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-secure-local-file-gate/v1",
  "approval": {
    "status": "AWAITING_USER",
    "approval_id": "REQUIRED",
    "approved_by": "REQUIRED",
    "approved_at": "REQUIRED"
  },
  "allowed_content": [],
  "limits": {
    "max_input_bytes": 52428800,
    "process_timeout_seconds": 120,
    "clamav_database_max_age_days": 2,
    "clamav_max_scan_bytes": 104857600,
    "clamav_max_files": 1000,
    "clamav_max_recursion": 8,
    "clamav_bytecode_timeout_ms": 5000,
    "yara_timeout_seconds": 30,
    "yara_max_matches_per_pattern": 100
  },
  "yara": {
    "compiled_rules_sha256": "REQUIRED",
    "reject_on_any_match": true
  },
  "archives": {
    "expansion_allowed": false
  },
  "business_storage_authorized": false
}
````

### FILE: `secure_file_gate/secure_local_file_gate.py`
```yaml
block_id: "SECURE-LOCAL-FILE-INGESTION-GATE:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local orchestration against exact official CLI contracts; no upstream source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "5f9beb949b72662bb6ac0c9e04cd5c200371a6807c123eb23ca1752609122fa8"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
import hashlib
import json
import os
import re
import shutil
import subprocess
import tempfile
from dataclasses import dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Callable, Sequence


SCHEMA = "elite-secure-local-file-gate/v1"
MAGIKA_VERSION = "magika 1.1.0 standard_v3_3"
CLAMAV_VERSION_PREFIX = "ClamAV 1.5.4/"
YARA_VERSION = "yara-x-cli 1.20.0"
OFFICIAL_BINARY_SHA256 = {
    "magika": "3631dab2f57ec42b6646ce141397671132707cd7e7633f4511fe47171efe69eb",
    "clamscan": "f368912f61ddea8b302acf9114891ec57f8a45b0de77f076b1d19c960bc6f7cf",
    "yara_x": "d3e648651f3eaa5e833faeb394fc7dc7dfb583a31dd5fc572ab6cae7ac41f3fd",
}
HEX_256 = re.compile(r"^[0-9a-f]{64}$")


@dataclass(frozen=True)
class CommandResult:
    returncode: int
    stdout: str
    stderr: str


Runner = Callable[[Sequence[str], int], CommandResult]


class GateConfigurationError(RuntimeError):
    pass


class GateStageError(RuntimeError):
    def __init__(self, reason: str):
        super().__init__(reason)
        self.reason = reason


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def utc_now() -> str:
    return datetime.now(timezone.utc).isoformat().replace("+00:00", "Z")


def default_runner(args: Sequence[str], timeout: int) -> CommandResult:
    completed = subprocess.run(
        list(args),
        check=False,
        capture_output=True,
        text=True,
        encoding="utf-8",
        errors="replace",
        timeout=timeout,
        creationflags=getattr(subprocess, "CREATE_NO_WINDOW", 0),
    )
    return CommandResult(completed.returncode, completed.stdout, completed.stderr)


def require_regular_file(path: Path, label: str) -> None:
    if path.is_symlink() or not path.is_file():
        raise GateConfigurationError(f"{label} must be a regular non-symlink file")


def require_directory(path: Path, label: str) -> None:
    if path.is_symlink() or not path.is_dir():
        raise GateConfigurationError(f"{label} must be a non-symlink directory")


def load_policy(path: Path) -> dict[str, Any]:
    require_regular_file(path, "policy")
    policy = json.loads(path.read_text(encoding="utf-8"))
    if policy.get("schema") != SCHEMA:
        raise GateConfigurationError("unsupported policy schema")
    approval = policy.get("approval") or {}
    if approval.get("status") != "APPROVED":
        raise GateConfigurationError("policy approval is not APPROVED")
    for field in ("approval_id", "approved_by", "approved_at"):
        value = approval.get(field)
        if not isinstance(value, str) or not value.strip() or value == "REQUIRED":
            raise GateConfigurationError(f"approval.{field} is unresolved")
    try:
        approved_at = datetime.fromisoformat(approval["approved_at"].replace("Z", "+00:00"))
    except ValueError as error:
        raise GateConfigurationError("approval.approved_at is not ISO-8601") from error
    if approved_at.tzinfo is None:
        raise GateConfigurationError("approval.approved_at must include a timezone")
    allowed = policy.get("allowed_content")
    if not isinstance(allowed, list) or not allowed:
        raise GateConfigurationError("allowed_content must contain approved label/mime pairs")
    for item in allowed:
        if not isinstance(item, dict) or not isinstance(item.get("label"), str) or not isinstance(item.get("mime_type"), str):
            raise GateConfigurationError("each allowed_content item requires label and mime_type")
        score = item.get("min_score")
        if not isinstance(score, (int, float)) or not 0 <= score <= 1:
            raise GateConfigurationError("each allowed_content item requires min_score from 0 to 1")
    limits = policy.get("limits") or {}
    for field in (
        "max_input_bytes", "process_timeout_seconds", "clamav_database_max_age_days",
        "clamav_max_scan_bytes", "clamav_max_files", "clamav_max_recursion",
        "clamav_bytecode_timeout_ms", "yara_timeout_seconds", "yara_max_matches_per_pattern",
    ):
        if not isinstance(limits.get(field), int) or limits[field] <= 0:
            raise GateConfigurationError(f"limits.{field} must be a positive integer")
    yara = policy.get("yara") or {}
    rules_hash = yara.get("compiled_rules_sha256")
    if not isinstance(rules_hash, str) or not HEX_256.fullmatch(rules_hash):
        raise GateConfigurationError("yara.compiled_rules_sha256 must be an exact lowercase SHA-256")
    if yara.get("reject_on_any_match") is not True:
        raise GateConfigurationError("yara.reject_on_any_match must be true")
    if (policy.get("archives") or {}).get("expansion_allowed") is not False:
        raise GateConfigurationError("archive expansion must remain disabled")
    if policy.get("business_storage_authorized") is not False:
        raise GateConfigurationError("this gate cannot authorize business storage")
    return policy


def validate_binary(path: Path, tool: str) -> str:
    require_regular_file(path, tool)
    actual = sha256_file(path)
    expected = OFFICIAL_BINARY_SHA256[tool]
    if actual != expected:
        raise GateConfigurationError(f"{tool} binary SHA-256 mismatch")
    return actual


def run_checked(runner: Runner, args: Sequence[str], timeout: int, reason: str) -> CommandResult:
    try:
        return runner(args, timeout)
    except (OSError, subprocess.TimeoutExpired) as error:
        raise GateStageError(reason) from error


def output_digest(result: CommandResult) -> str:
    return hashlib.sha256((result.stdout + "\n" + result.stderr).encode("utf-8")).hexdigest()


def commit_receipt(output: Path, receipt: dict[str, Any]) -> None:
    require_directory(output.parent, "output parent")
    if output.exists() or output.is_symlink():
        raise FileExistsError(f"output already exists: {output}")
    staging = Path(tempfile.mkdtemp(prefix=f".{output.name}.staging-", dir=output.parent))
    try:
        payload = json.dumps(receipt, indent=2, sort_keys=True, ensure_ascii=False) + "\n"
        (staging / "security-receipt.json").write_text(payload, encoding="utf-8", newline="\n")
        os.replace(staging, output)
    except Exception:
        shutil.rmtree(staging, ignore_errors=True)
        raise


def reject(receipt: dict[str, Any], reason: str) -> dict[str, Any]:
    receipt["decision"] = "REJECTED"
    receipt["reason"] = reason
    receipt["completed_at"] = utc_now()
    return receipt


def run_gate(
    *, input_path: Path, output: Path, policy_path: Path, magika_exe: Path,
    clamscan_exe: Path, clamav_database: Path, yara_exe: Path,
    compiled_rules: Path, runner: Runner = default_runner,
) -> dict[str, Any]:
    policy = load_policy(policy_path)
    require_regular_file(input_path, "input")
    require_directory(clamav_database, "ClamAV database")
    require_regular_file(compiled_rules, "compiled YARA-X rules")
    size = input_path.stat().st_size
    limits = policy["limits"]
    if size > limits["max_input_bytes"]:
        raise GateConfigurationError("input exceeds the approved byte limit")
    if sha256_file(compiled_rules) != policy["yara"]["compiled_rules_sha256"]:
        raise GateConfigurationError("compiled YARA-X rules SHA-256 mismatch")
    tool_hashes = {
        "magika": validate_binary(magika_exe, "magika"),
        "clamscan": validate_binary(clamscan_exe, "clamscan"),
        "yara_x": validate_binary(yara_exe, "yara_x"),
    }
    receipt: dict[str, Any] = {
        "schema": "elite-secure-local-file-receipt/v1",
        "decision": "IN_PROGRESS",
        "reason": None,
        "started_at": utc_now(),
        "approval_id": policy["approval"]["approval_id"],
        "policy_sha256": sha256_file(policy_path),
        "input": {"sha256": sha256_file(input_path), "bytes": size},
        "tools": {name: {"sha256": value} for name, value in tool_hashes.items()},
        "business_storage_authorized": False,
        "security_claim": "No selected local detector matched within approved limits; this is not a universal safety guarantee.",
    }
    timeout = limits["process_timeout_seconds"]
    try:
        clam_version = run_checked(runner, [str(clamscan_exe), f"--database={clamav_database}", "--version"], timeout, "CLAMAV_VERSION_ERROR")
        if clam_version.returncode != 0 or not clam_version.stdout.strip().startswith(CLAMAV_VERSION_PREFIX):
            raise GateStageError("CLAMAV_VERSION_MISMATCH")
        receipt["tools"]["clamscan"]["version"] = clam_version.stdout.strip()
        clam = run_checked(runner, [
            str(clamscan_exe), f"--database={clamav_database}", "--official-db-only=yes",
            f"--fail-if-cvd-older-than={limits['clamav_database_max_age_days']}", "--stdout",
            "--no-summary", "--infected", "--follow-file-symlinks=0", "--follow-dir-symlinks=0",
            "--alert-encrypted=yes", "--alert-exceeds-max=yes",
            f"--max-filesize={limits['max_input_bytes']}", f"--max-scansize={limits['clamav_max_scan_bytes']}",
            f"--max-files={limits['clamav_max_files']}", f"--max-recursion={limits['clamav_max_recursion']}",
            f"--bytecode-timeout={limits['clamav_bytecode_timeout_ms']}", str(input_path),
        ], timeout, "CLAMAV_EXECUTION_ERROR")
        receipt["clamav"] = {"exit_code": clam.returncode, "output_sha256": output_digest(clam)}
        if clam.returncode == 1:
            raise GateStageError("CLAMAV_DETECTED")
        if clam.returncode != 0:
            raise GateStageError("CLAMAV_ERROR_OR_STALE_DATABASE")

        magika_version = run_checked(runner, [str(magika_exe), "--version"], timeout, "MAGIKA_VERSION_ERROR")
        if magika_version.returncode != 0 or magika_version.stdout.strip() != MAGIKA_VERSION:
            raise GateStageError("MAGIKA_VERSION_MISMATCH")
        receipt["tools"]["magika"]["version"] = magika_version.stdout.strip()
        magika = run_checked(runner, [str(magika_exe), "--json", "--no-colors", str(input_path)], timeout, "MAGIKA_EXECUTION_ERROR")
        if magika.returncode != 0:
            raise GateStageError("MAGIKA_ERROR")
        try:
            items = json.loads(magika.stdout)
            value = items[0]["result"]["value"]
            result = value["output"]
            label, mime_type, score = result["label"], result["mime_type"], float(value["score"])
        except (KeyError, IndexError, TypeError, ValueError, json.JSONDecodeError) as error:
            raise GateStageError("MAGIKA_INVALID_JSON") from error
        approved = any(
            item["label"] == label and item["mime_type"] == mime_type and score >= item["min_score"]
            for item in policy["allowed_content"]
        )
        receipt["content_type"] = {"label": label, "mime_type": mime_type, "score": score}
        if not approved:
            raise GateStageError("CONTENT_TYPE_NOT_APPROVED")

        yara_version = run_checked(runner, [str(yara_exe), "--version"], timeout, "YARA_VERSION_ERROR")
        if yara_version.returncode != 0 or yara_version.stdout.strip() != YARA_VERSION:
            raise GateStageError("YARA_VERSION_MISMATCH")
        receipt["tools"]["yara_x"]["version"] = yara_version.stdout.strip()
        yara = run_checked(runner, [
            str(yara_exe), "scan", "--compiled-rules", "--output-format=json", "--no-mmap", "--threads", "1",
            "--timeout", str(limits["yara_timeout_seconds"]),
            "--max-matches-per-pattern", str(limits["yara_max_matches_per_pattern"]),
            str(compiled_rules), str(input_path),
        ], timeout, "YARA_EXECUTION_ERROR")
        if yara.returncode != 0:
            raise GateStageError("YARA_ERROR")
        try:
            yara_payload = json.loads(yara.stdout)
            if yara_payload.get("version") != "1.20.0" or not isinstance(yara_payload.get("matches"), list):
                raise ValueError("invalid schema")
            matches = [str(item["rule"]) for item in yara_payload["matches"]]
        except (KeyError, TypeError, ValueError, json.JSONDecodeError) as error:
            raise GateStageError("YARA_INVALID_JSON") from error
        receipt["yara_x"] = {
            "compiled_rules_sha256": policy["yara"]["compiled_rules_sha256"],
            "matching_rules": sorted(matches),
        }
        if matches:
            raise GateStageError("YARA_RULE_MATCH")
        receipt["decision"] = "ADMITTED"
        receipt["reason"] = "ALL_SELECTED_GATES_PASSED"
        receipt["completed_at"] = utc_now()
    except GateStageError as error:
        reject(receipt, error.reason)
    commit_receipt(output, receipt)
    return receipt


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--input", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--policy", required=True, type=Path)
    parser.add_argument("--magika-exe", required=True, type=Path)
    parser.add_argument("--clamscan-exe", required=True, type=Path)
    parser.add_argument("--clamav-database", required=True, type=Path)
    parser.add_argument("--yara-exe", required=True, type=Path)
    parser.add_argument("--compiled-rules", required=True, type=Path)
    args = parser.parse_args()
    receipt = run_gate(
        input_path=args.input, output=args.output, policy_path=args.policy,
        magika_exe=args.magika_exe, clamscan_exe=args.clamscan_exe,
        clamav_database=args.clamav_database, yara_exe=args.yara_exe,
        compiled_rules=args.compiled_rules,
    )
    print(f"SECURE_LOCAL_FILE_GATE decision={receipt['decision']} reason={receipt['reason']}")
    return 0 if receipt["decision"] == "ADMITTED" else 2


if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `secure_file_gate/prepare_security_assets.py`
```yaml
block_id: "SECURE-LOCAL-FILE-INGESTION-GATE:assets:v1"
operation: CREATE
provenance: AUTHORED
source: "local preparation and proof wrapper around official freshclam, clamscan and YARA-X CLI"
license: "LicenseRef-Workspace-Owner"
sha256: "1fc86a14c7089989d2d264ca62c79faf522b734a7a19bab977f81b181fe0a28f"
variables: []
secrets_allowed: false
```
````python
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
````

### FILE: `secure_file_gate/test_secure_local_file_gate.py`
```yaml
block_id: "SECURE-LOCAL-FILE-INGESTION-GATE:runner-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local deterministic negative and atomicity regression suite"
license: "LicenseRef-Workspace-Owner"
sha256: "5bbb5ee670c6742a5a817c9fc20837c63fb7f27c7af9468d403f6d6154ccf156"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import json
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

import secure_local_file_gate as gate


class FakeRunner:
    def __init__(self) -> None:
        self.clam_exit = 0
        self.magika_label = "pdf"
        self.magika_mime = "application/pdf"
        self.magika_score = 0.99
        self.yara_matches: list[dict[str, str]] = []
        self.calls: list[list[str]] = []

    def __call__(self, args, timeout):
        call = [str(value) for value in args]
        self.calls.append(call)
        name = Path(call[0]).name
        if "--version" in call:
            if name == "clamscan.exe":
                return gate.CommandResult(0, "ClamAV 1.5.4/28104/Wed Aug 26 03:24:01 2026\n", "")
            if name == "magika.exe":
                return gate.CommandResult(0, gate.MAGIKA_VERSION + "\n", "")
            return gate.CommandResult(0, gate.YARA_VERSION + "\n", "")
        if name == "clamscan.exe":
            return gate.CommandResult(self.clam_exit, "detector output" if self.clam_exit else "", "")
        if name == "magika.exe":
            payload = [{"result": {"value": {"output": {
                "label": self.magika_label, "mime_type": self.magika_mime,
            }, "score": self.magika_score}}}]
            return gate.CommandResult(0, json.dumps(payload), "")
        return gate.CommandResult(0, json.dumps({"version": "1.20.0", "matches": self.yara_matches}), "")


class SecureLocalFileGateTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary = tempfile.TemporaryDirectory()
        self.root = Path(self.temporary.name)
        self.input = self.root / "invoice.pdf"
        self.input.write_bytes(b"%PDF-1.7\nfixture")
        self.policy = self.root / "policy.json"
        self.policy.write_text(json.dumps({
            "schema": gate.SCHEMA,
            "approval": {"status": "APPROVED", "approval_id": "SEC-1", "approved_by": "owner", "approved_at": "2026-08-26T00:00:00Z"},
            "allowed_content": [{"label": "pdf", "mime_type": "application/pdf", "min_score": 0.9}],
            "limits": {
                "max_input_bytes": 1024, "process_timeout_seconds": 10,
                "clamav_database_max_age_days": 2, "clamav_max_scan_bytes": 2048,
                "clamav_max_files": 10, "clamav_max_recursion": 2,
                "clamav_bytecode_timeout_ms": 1000, "yara_timeout_seconds": 5,
                "yara_max_matches_per_pattern": 10,
            },
            "yara": {"compiled_rules_sha256": "0" * 64, "reject_on_any_match": True},
            "archives": {"expansion_allowed": False}, "business_storage_authorized": False,
        }), encoding="utf-8")
        self.db = self.root / "db"
        self.db.mkdir()
        self.rules = self.root / "rules.yarc"
        self.rules.write_bytes(b"compiled")
        data = json.loads(self.policy.read_text(encoding="utf-8"))
        data["yara"]["compiled_rules_sha256"] = gate.sha256_file(self.rules)
        self.policy.write_text(json.dumps(data), encoding="utf-8")
        self.magika = self.root / "magika.exe"
        self.clam = self.root / "clamscan.exe"
        self.yara = self.root / "yr.exe"
        for path in (self.magika, self.clam, self.yara):
            path.write_bytes(path.name.encode("ascii"))
        self.hashes = {
            "magika": gate.sha256_file(self.magika),
            "clamscan": gate.sha256_file(self.clam),
            "yara_x": gate.sha256_file(self.yara),
        }

    def tearDown(self) -> None:
        self.temporary.cleanup()

    def execute(self, runner: FakeRunner, name: str = "evidence"):
        with patch.dict(gate.OFFICIAL_BINARY_SHA256, self.hashes, clear=True):
            return gate.run_gate(
                input_path=self.input, output=self.root / name, policy_path=self.policy,
                magika_exe=self.magika, clamscan_exe=self.clam, clamav_database=self.db,
                yara_exe=self.yara, compiled_rules=self.rules, runner=runner,
            )

    def test_all_selected_gates_admit_and_commit_hash_receipt(self):
        receipt = self.execute(FakeRunner())
        self.assertEqual(receipt["decision"], "ADMITTED")
        saved = json.loads((self.root / "evidence" / "security-receipt.json").read_text(encoding="utf-8"))
        self.assertEqual(saved["input"]["sha256"], gate.sha256_file(self.input))
        self.assertFalse(saved["business_storage_authorized"])

    def test_clamav_detection_rejects_before_other_content_engines(self):
        runner = FakeRunner()
        runner.clam_exit = 1
        receipt = self.execute(runner)
        self.assertEqual(receipt["reason"], "CLAMAV_DETECTED")
        self.assertFalse(any(Path(call[0]).name in {"magika.exe", "yr.exe"} for call in runner.calls))

    def test_stale_or_broken_clamav_database_rejects_closed(self):
        runner = FakeRunner()
        runner.clam_exit = 2
        receipt = self.execute(runner)
        self.assertEqual(receipt["reason"], "CLAMAV_ERROR_OR_STALE_DATABASE")

    def test_unapproved_content_type_rejects_before_yara(self):
        runner = FakeRunner()
        runner.magika_label = "zip"
        runner.magika_mime = "application/zip"
        receipt = self.execute(runner)
        self.assertEqual(receipt["reason"], "CONTENT_TYPE_NOT_APPROVED")
        self.assertFalse(any(Path(call[0]).name == "yr.exe" for call in runner.calls))

    def test_any_yara_match_rejects_and_records_only_rule_identifier(self):
        runner = FakeRunner()
        runner.yara_matches = [{"rule": "blocked_pattern", "file": "private-path"}]
        receipt = self.execute(runner)
        self.assertEqual(receipt["reason"], "YARA_RULE_MATCH")
        self.assertEqual(receipt["yara_x"]["matching_rules"], ["blocked_pattern"])
        self.assertNotIn("private-path", json.dumps(receipt))

    def test_awaiting_user_policy_blocks_before_any_tool(self):
        data = json.loads(self.policy.read_text(encoding="utf-8"))
        data["approval"]["status"] = "AWAITING_USER"
        self.policy.write_text(json.dumps(data), encoding="utf-8")
        runner = FakeRunner()
        with self.assertRaisesRegex(gate.GateConfigurationError, "not APPROVED"):
            self.execute(runner)
        self.assertEqual(runner.calls, [])

    def test_binary_hash_mismatch_blocks_before_execution(self):
        runner = FakeRunner()
        self.magika.write_bytes(b"tampered")
        with self.assertRaisesRegex(gate.GateConfigurationError, "SHA-256 mismatch"):
            self.execute(runner)
        self.assertEqual(runner.calls, [])


if __name__ == "__main__":
    unittest.main()
````

### FILE: `secure_file_gate/test_prepare_security_assets.py`
```yaml
block_id: "SECURE-LOCAL-FILE-INGESTION-GATE:asset-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local freshclam distrust, database proof, YARA compile and overwrite regression suite"
license: "LicenseRef-Workspace-Owner"
sha256: "b7be740a155b0c0acfe9471a3d4e83101f3c9b1bf9e7370149f11cab0bcd7984"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import json
import subprocess
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

import prepare_security_assets as assets


class PrepareSecurityAssetsTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary = tempfile.TemporaryDirectory()
        self.root = Path(self.temporary.name)
        self.fresh = self.root / "freshclam.exe"
        self.clam = self.root / "clamscan.exe"
        self.yara = self.root / "yr.exe"
        for path in (self.fresh, self.clam, self.yara):
            path.write_bytes(path.name.encode("ascii"))

    def tearDown(self) -> None:
        self.temporary.cleanup()

    def completed(self, code=0, stdout="", stderr=""):
        return subprocess.CompletedProcess([], code, stdout, stderr)

    def test_update_does_not_trust_freshclam_exit_without_scan_proof(self):
        db = self.root / "db"
        hashes = {
            "FRESHCLAM_SHA256": assets.sha256_file(self.fresh),
            "CLAMSCAN_SHA256": assets.sha256_file(self.clam),
        }
        def fake_run(args, timeout):
            if Path(args[0]).name == "freshclam.exe":
                db.mkdir(exist_ok=True)
                for name in ("main.cvd", "daily.cvd", "bytecode.cvd"):
                    (db / name).write_bytes(name.encode("ascii"))
                return self.completed(0)
            if "--version" in args:
                return self.completed(0, "ClamAV 1.5.4/28104/Wed Aug 26 03:24:01 2026\n")
            return self.completed(2, "", "stale")
        with patch.multiple(assets, **hashes), patch.object(assets, "run", side_effect=fake_run):
            with self.assertRaisesRegex(RuntimeError, "not proven"):
                assets.update_clamav(self.fresh, self.clam, db, 2, 30)
        self.assertFalse((db / "clamav-database-receipt.json").exists())

    def test_update_records_inventory_only_after_fresh_scan(self):
        db = self.root / "db"
        hashes = {
            "FRESHCLAM_SHA256": assets.sha256_file(self.fresh),
            "CLAMSCAN_SHA256": assets.sha256_file(self.clam),
        }
        def fake_run(args, timeout):
            if Path(args[0]).name == "freshclam.exe":
                db.mkdir(exist_ok=True)
                for name in ("main.cvd", "daily.cvd", "bytecode.cvd"):
                    (db / name).write_bytes(name.encode("ascii"))
                return self.completed(0)
            if "--version" in args:
                return self.completed(0, "ClamAV 1.5.4/28104/Wed Aug 26 03:24:01 2026\n")
            return self.completed(0)
        with patch.multiple(assets, **hashes), patch.object(assets, "run", side_effect=fake_run):
            assets.update_clamav(self.fresh, self.clam, db, 2, 30)
        receipt = json.loads((db / "clamav-database-receipt.json").read_text(encoding="utf-8"))
        self.assertEqual(len(receipt["files"]), 3)
        self.assertIn("official-db-only", receipt["validation"])

    def test_yara_compile_is_atomic_and_hash_receipted(self):
        source = self.root / "approved.yar"
        source.write_text("rule approved { condition: false }\n", encoding="utf-8")
        output = self.root / "approved.yarc"
        def fake_run(args, timeout):
            if "compile" in args:
                Path(args[args.index("--output") + 1]).write_bytes(b"compiled")
                return self.completed(0)
            return self.completed(0, '{"version":"1.20.0","matches":[]}')
        with patch.object(assets, "YARA_SHA256", assets.sha256_file(self.yara)), patch.object(assets, "run", side_effect=fake_run):
            assets.compile_yara(self.yara, source, output, 30)
        receipt = json.loads(output.with_suffix(".yarc.receipt.json").read_text(encoding="utf-8"))
        self.assertEqual(receipt["compiled_rules_sha256"], assets.sha256_file(output))

    def test_existing_yara_output_fails_without_overwrite(self):
        source = self.root / "approved.yar"
        source.write_text("rule approved { condition: false }\n", encoding="utf-8")
        output = self.root / "approved.yarc"
        output.write_bytes(b"existing")
        with patch.object(assets, "YARA_SHA256", assets.sha256_file(self.yara)):
            with self.assertRaises(FileExistsError):
                assets.compile_yara(self.yara, source, output, 30)
        self.assertEqual(output.read_bytes(), b"existing")


if __name__ == "__main__":
    unittest.main()
````

### FILE: `secure_file_gate/README.md`
```yaml
block_id: "SECURE-LOCAL-FILE-INGESTION-GATE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local operating instructions tied to exact official upstream identities and limits"
license: "LicenseRef-Workspace-Owner"
sha256: "62b4f28e84eef4b6ffcbc46672dd18a3c6694461d5ea4eef3fda370cc3e2d59e"
variables: []
secrets_allowed: false
```
````markdown
# Secure Local File Gate

This pack is authored orchestration around three exact official Windows x86-64 release binaries: Google Magika CLI 1.1.0, Cisco Talos ClamAV 1.5.4 and VirusTotal/Google YARA-X 1.20.0. It does not copy or relabel their source. Acquire the pinned assets through `OFFICIAL-UPSTREAM-ACQUISITION-CORE` and retain their Apache-2.0, GPL-2.0-only and BSD-3-Clause terms.

The policy template intentionally starts at `AWAITING_USER`, has no allowed content type and contains an unresolved YARA rules hash. The agent must obtain a real owner, approval ID/date, exact label/MIME/score pairs, limits and an approved ruleset before execution. The runner refuses unresolved configuration.

Prepare mutable ClamAV signatures and compiled project rules:

```powershell
python prepare_security_assets.py update-clamav --freshclam <freshclam.exe> --clamscan <clamscan.exe> --database <database-directory>
python prepare_security_assets.py compile-yara --yara <yr.exe> --source <approved-rules> --output <new-rules.yarc>
```

Do not trust `freshclam` exit status by itself. The preparer also makes `clamscan` load the official databases, enforce maximum age and scan a benign probe before writing a hash inventory receipt. YARA source provenance, review, false-positive testing, canary and rollback remain project approvals; this library never supplies invented detection rules.

Run the gate only inside a sandboxed worker with denied egress, bounded CPU/memory/time and separate quarantine/clean storage:

```powershell
python secure_local_file_gate.py --input <quarantined-file> --output <new-evidence-directory> --policy <approved-policy.json> --magika-exe <magika.exe> --clamscan-exe <clamscan.exe> --clamav-database <database-directory> --yara-exe <yr.exe> --compiled-rules <approved-rules.yarc>
```

`ADMITTED` means only that the selected official engines passed the approved limits and rules at that moment. It is not a universal claim that the file is safe or semantically correct, and the receipt always keeps `business_storage_authorized=false`. Conversion and field extraction occur later, with an approved analyzer/schema/corpus and field-level evidence.
````

## 6. Configuration surface

| Campo | Default | Regla | Cambio |
|---|---|---|---|
| approval | `AWAITING_USER` | owner/id/fecha con timezone y `APPROVED` | nuevo approval |
| allowed_content | vacío | label + MIME + score mínimo exactos | corpus/gate de aceptación |
| límites | conservadores | bytes/timeout/edad DB/recursión/matches positivos | carga y threat test |
| compiled_rules_sha256 | `REQUIRED` | 64 hex lowercase; cualquier match rechaza | revisión/canary/rollback |
| archive expansion | false | no habilitable en este pack | pack separado aprobado |
| business storage | false | inmutable | analyzer/schema/corpus posterior |

## 7. Dependency bill

| Dependencia | Identidad exacta | Uso | Licencia | Gate |
|---|---|---|---|---|
| Google Magika | CLI 1.1.0, commit `5e2f437…`, exe SHA `3631dab2…` | tipo por contenido | Apache-2.0 | versión/hash/JSON/allowlist |
| Cisco Talos ClamAV | 1.5.4, commit `fa59fca…`, clamscan SHA `f368912f…`, freshclam SHA `4032fdd4…` | antimalware y firmas oficiales | GPL-2.0-only | firmas/edad/probe/límites/exit 0-1-2 |
| VirusTotal/Google YARA-X | 1.20.0, commit `60ad069…`, exe SHA `d3e64865…` | reglas adicionales | BSD-3-Clause | compile/hash/JSON/timeout/reject-any |
| CPython | 3.12 stdlib | coordinación y receipts | runtime target | compileall + 11 tests |

## 8. Apply order

1. Adquirir los tres assets exactos mediante el source profile y aceptar licencias.
2. Completar política y aprobación; no usar el template bloqueado.
3. Actualizar bases con `prepare_security_assets.py update-clamav` y comprobar receipt.
4. Revisar reglas, falsos positivos, canary/rollback; compilar con `compile-yara` y copiar el SHA a política.
5. Ejecutar once tests y un fixture benigno autorizado con binarios/bases reales.
6. Desplegar worker aislado, sin secretos ni acceso a tablas de negocio.
7. Sólo después conectar conversión/extracción; el receipt nunca sustituye evidencia semántica por campo.

## 9. Verification

```powershell
python -m compileall -q secure_file_gate
python -m unittest discover -s secure_file_gate -p 'test_*.py' -v
python secure_file_gate/prepare_security_assets.py update-clamav --freshclam <freshclam.exe> --clamscan <clamscan.exe> --database <db>
python secure_file_gate/prepare_security_assets.py compile-yara --yara <yr.exe> --source <approved-rules> --output <new.yarc>
python secure_file_gate/secure_local_file_gate.py --input <quarantined-file> --output <new-evidence> --policy <approved-policy> --magika-exe <magika.exe> --clamscan-exe <clamscan.exe> --clamav-database <db> --yara-exe <yr.exe> --compiled-rules <new.yarc>
```

Esperado: 11 tests PASS; preparación sólo emite receipt tras scan fresco real; fixture benigno permitido produce `ADMITTED`; binario/ruleset alterado, base vieja/error, tipo no aprobado, match YARA y política incompleta producen rechazo sin almacenamiento de negocio.

## 10. Reconstruction evidence

Evidencia gobernante: `reconstruction_evidence/SECURE_LOCAL_FILE_INGESTION_GATE_2026-08-26_V1.md`.

- tres releases/commits/assets/licencias oficiales exactos;
- contratos CLI y JSON verificados con binarios release;
- bases main 63, daily 28104 y bytecode 339 actualizadas y probadas;
- YARA-X source/compiled scan reales;
- 11 unit tests y un recorrido real de los tres motores;
- no EICAR atribuido falsamente a ClamAV: Defender lo interceptó en la auditoría previa.
