from __future__ import annotations

import argparse
from datetime import datetime, timezone
import hashlib
import json
from pathlib import Path, PurePosixPath
import re
import sys
import tempfile

DIGEST = re.compile(r"sha256:[0-9a-f]{64}")
SHA = re.compile(r"[0-9a-f]{64}")
KINDS = ("BASELINE", "SOAK", "SATURATION_RECOVERY")
PROFILE_FIELDS = {"schema", "project_id", "environment", "release_digest", "evaluated_at", "expires_at", "executor_ref", "tool", "target", "error_budget", "runs", "output"}
RUN_FIELDS = {"kind", "execution_receipt", "summary", "script", "expected_arguments"}
K6_SHA256 = "87dfa91bc3e47bc4bd77911d59d7ff79f25cd76fb8322c97072f08f08a0da5ed"


def canonical_sha(value: object) -> str:
    return hashlib.sha256(json.dumps(value, ensure_ascii=False, separators=(",", ":")).encode()).hexdigest()


def safe(root: Path, raw: object, label: str, *, exists: bool = True) -> Path:
    if not isinstance(raw, str) or not raw or "\\" in raw:
        raise ValueError(f"{label} must be a non-empty POSIX relative path")
    pure = PurePosixPath(raw)
    if pure.is_absolute() or any(part in ("", ".", "..") for part in pure.parts):
        raise ValueError(f"{label} is unsafe")
    path = root.joinpath(*pure.parts).resolve(strict=exists)
    path.relative_to(root)
    if exists and not path.is_file():
        raise ValueError(f"{label} must be a file")
    return path


def read_json(path: Path, label: str) -> dict[str, object]:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, UnicodeError, json.JSONDecodeError) as error:
        raise ValueError(f"{label} is not valid UTF-8 JSON") from error
    if not isinstance(value, dict):
        raise ValueError(f"{label} must be an object")
    return value


def proof(root: Path, path: Path) -> dict[str, object]:
    data = path.read_bytes()
    return {"path": path.relative_to(root).as_posix(), "bytes": len(data), "sha256": hashlib.sha256(data).hexdigest()}


def verify_proof(root: Path, item: object, label: str) -> Path:
    if not isinstance(item, dict) or set(item) != {"path", "bytes", "sha256"}:
        raise ValueError(f"{label} proof fields must be exact")
    path = safe(root, item["path"], label)
    actual = proof(root, path)
    if actual != item:
        raise ValueError(f"{label} proof mismatch")
    return path


def utc(raw: object, label: str) -> datetime:
    if not isinstance(raw, str) or not raw.endswith("Z"):
        raise ValueError(f"{label} must be ISO-8601 UTC ending Z")
    try:
        return datetime.fromisoformat(raw[:-1] + "+00:00").astimezone(timezone.utc)
    except ValueError as error:
        raise ValueError(f"{label} must be ISO-8601 UTC ending Z") from error


def metric(summary: dict[str, object], name: str) -> dict[str, object]:
    results = summary.get("results")
    metrics = results.get("metrics") if isinstance(results, dict) else None
    if not isinstance(metrics, list):
        raise ValueError("k6 summary results.metrics must be a list")
    matches = [item for item in metrics if isinstance(item, dict) and item.get("name") == name]
    if len(matches) != 1 or not isinstance(matches[0].get("values"), dict):
        raise ValueError(f"k6 summary requires exactly one {name} metric")
    return matches[0]["values"]


def number(values: dict[str, object], key: str, label: str) -> float:
    value = values.get(key)
    if isinstance(value, bool) or not isinstance(value, (int, float)) or value < 0:
        raise ValueError(f"{label} must be a non-negative number")
    return float(value)


def atomic(path: Path, data: bytes) -> None:
    if path.exists():
        raise ValueError("output already exists")
    path.parent.mkdir(parents=False, exist_ok=True)
    with tempfile.NamedTemporaryFile(dir=path.parent, prefix=f".{path.name}.", delete=False) as handle:
        temporary = Path(handle.name); handle.write(data); handle.flush()
    try:
        temporary.replace(path)
    except Exception:
        temporary.unlink(missing_ok=True); raise


def validate(profile: object, root: Path) -> dict[str, object]:
    if not isinstance(profile, dict) or set(profile) != PROFILE_FIELDS:
        raise ValueError("profile fields must be exact")
    if profile["schema"] != "elite-k6-load-admission/v1" or profile["environment"] != "production":
        raise ValueError("profile identity is invalid")
    if not isinstance(profile["project_id"], str) or not profile["project_id"].strip() or not DIGEST.fullmatch(str(profile["release_digest"])):
        raise ValueError("project/release identity is invalid")
    start = utc(profile["evaluated_at"], "evaluated_at"); end = utc(profile["expires_at"], "expires_at")
    if not start < end or (end - start).total_seconds() > 30 * 86400:
        raise ValueError("admission validity must be positive and at most 30 days")
    if not isinstance(profile["executor_ref"], str) or not profile["executor_ref"].strip() or not isinstance(profile["target"], dict) or not profile["target"]:
        raise ValueError("executor_ref and real target are required")
    tool = profile["tool"]
    expected_tool = {"name":"Grafana k6", "version":"2.2.0", "source":"https://github.com/grafana/k6/releases/tag/v2.2.0", "sha256":K6_SHA256}
    if tool != expected_tool:
        raise ValueError("exact Grafana k6 2.2.0 tool identity is required")

    budget_path = safe(root, profile["error_budget"], "error_budget")
    budget = read_json(budget_path, "error_budget")
    budget_fields = {"schema", "project_id", "environment", "release_digest", "approved_at", "approvers", "runs"}
    if set(budget) != budget_fields or budget.get("schema") != "elite-load-error-budget/v1":
        raise ValueError("error budget fields/schema must be exact")
    for key in ("project_id", "environment", "release_digest"):
        if budget.get(key) != profile[key]: raise ValueError(f"error budget {key} mismatch")
    utc(budget["approved_at"], "error_budget.approved_at")
    approvers = budget["approvers"]
    if not isinstance(approvers, list) or len(approvers) < 2 or any(not isinstance(x, str) or not x.strip() for x in approvers) or len(set(approvers)) != len(approvers):
        raise ValueError("error budget requires at least two distinct approvers")
    limits = budget["runs"]
    if not isinstance(limits, dict) or set(limits) != set(KINDS):
        raise ValueError("error budget must define exactly three run kinds")

    runs = profile["runs"]
    if not isinstance(runs, list) or len(runs) != 3 or any(not isinstance(run, dict) or set(run) != RUN_FIELDS for run in runs):
        raise ValueError("runs must contain three exact records")
    if [run["kind"] for run in runs] != list(KINDS):
        raise ValueError("runs must be ordered BASELINE, SOAK, SATURATION_RECOVERY")
    evidence = [proof(root, budget_path)]
    for run in runs:
        kind = run["kind"]; limit = limits[kind]
        required_limits = {"minimum_iterations", "minimum_duration_seconds", "maximum_http_req_failed_rate", "maximum_http_req_duration_p95_ms"}
        if not isinstance(limit, dict) or set(limit) != required_limits:
            raise ValueError(f"{kind} budget fields must be exact")
        minimum_iterations = number(limit, "minimum_iterations", f"{kind}.minimum_iterations")
        minimum_duration = number(limit, "minimum_duration_seconds", f"{kind}.minimum_duration_seconds")
        maximum_failed = number(limit, "maximum_http_req_failed_rate", f"{kind}.maximum_http_req_failed_rate")
        maximum_p95 = number(limit, "maximum_http_req_duration_p95_ms", f"{kind}.maximum_http_req_duration_p95_ms")
        if minimum_iterations < 1 or minimum_duration < 1 or maximum_failed > 1 or maximum_p95 <= 0:
            raise ValueError(f"{kind} budget values are invalid")
        script_path = safe(root, run["script"], f"{kind}.script"); script_proof = proof(root, script_path)
        receipt_path = safe(root, run["execution_receipt"], f"{kind}.execution_receipt"); receipt = read_json(receipt_path, f"{kind}.execution_receipt")
        if receipt.get("schema") != "elite-official-tool-execution/v1" or receipt.get("project_id") != profile["project_id"] or receipt.get("environment") != "production" or receipt.get("release_digest") != profile["release_digest"] or receipt.get("control_id") != "LOAD_RESILIENCE" or receipt.get("exit_code") != 0:
            raise ValueError(f"{kind} execution identity/exit mismatch")
        receipt_tool = receipt.get("tool")
        if receipt_tool != {"name":tool["name"], "version":tool["version"], "source":tool["source"], "sha256":tool["sha256"]}:
            raise ValueError(f"{kind} tool identity mismatch")
        arguments = run["expected_arguments"]
        if not isinstance(arguments, list) or any(not isinstance(value, str) or not value for value in arguments) or receipt.get("arguments_sha256") != canonical_sha(arguments):
            raise ValueError(f"{kind} arguments are not hash-bound")
        expected_target = dict(profile["target"]); expected_target.update({"run_kind":kind, "workload_script_sha256":script_proof["sha256"], "minimum_duration_seconds":minimum_duration})
        if receipt.get("target") != expected_target:
            raise ValueError(f"{kind} target/script/duration binding mismatch")
        for channel in ("stdout", "stderr"):
            verify_proof(root, receipt.get(channel), f"{kind}.{channel}")
        summary_path = safe(root, run["summary"], f"{kind}.summary"); summary = read_json(summary_path, f"{kind}.summary")
        metadata = summary.get("metadata")
        if summary.get("version") != "1.0.0" or not isinstance(metadata, dict) or metadata.get("k6_version") != "2.2.0":
            raise ValueError(f"{kind} requires k6 v2 machine-readable summary")
        if number(metric(summary, "iterations"), "count", f"{kind}.iterations") < minimum_iterations:
            raise ValueError(f"{kind} iterations below approved minimum")
        if number(metric(summary, "http_req_failed"), "rate", f"{kind}.http_req_failed") > maximum_failed:
            raise ValueError(f"{kind} failure rate exceeds approved budget")
        if number(metric(summary, "http_req_duration"), "p(95)", f"{kind}.http_req_duration") > maximum_p95:
            raise ValueError(f"{kind} p95 exceeds approved budget")
        evidence.extend([proof(root, receipt_path), proof(root, summary_path), script_proof, receipt["stdout"], receipt["stderr"]])

    return {"id":"LOAD_RESILIENCE", "project_id":profile["project_id"], "environment":"production", "release_digest":profile["release_digest"], "result":"PASS", "executed_at":profile["evaluated_at"], "expires_at":profile["expires_at"], "executor_ref":profile["executor_ref"], "tool":{"name":tool["name"], "version":tool["version"], "source":tool["source"], "digest":"sha256:" + tool["sha256"]}, "target":profile["target"], "assertions":{"k6_thresholds_pass":True, "expected_workload_represented":True, "soak_pass":True, "saturation_recovery_pass":True, "no_error_budget_breach":True}, "evidence":evidence}


def main() -> int:
    parser = argparse.ArgumentParser(); parser.add_argument("--project-root", required=True, type=Path); parser.add_argument("--profile", required=True); args = parser.parse_args()
    try:
        root = args.project_root.resolve(strict=True); profile_path = safe(root, args.profile, "profile"); profile = read_json(profile_path, "profile")
        output = safe(root, profile.get("output"), "output", exists=False); output.parent.resolve(strict=True).relative_to(root)
        receipt = validate(profile, root); atomic(output, (json.dumps(receipt, indent=2, sort_keys=True) + "\n").encode())
        print(f"K6_LOAD_ADMISSION_PASS project={receipt['project_id']} release={receipt['release_digest']}"); return 0
    except (ValueError, OSError) as error:
        print(f"K6_LOAD_ADMISSION_FAILED: {error}", file=sys.stderr); return 1


if __name__ == "__main__": sys.exit(main())
