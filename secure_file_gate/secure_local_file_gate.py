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
