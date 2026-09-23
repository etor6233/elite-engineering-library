from __future__ import annotations

import argparse
import hashlib
import json
import re
import subprocess
import sys
import time
from datetime import UTC, datetime
from pathlib import Path
from typing import Any

CATEGORIES = {"format", "unit", "property", "static", "build", "security", "supply-chain", "integration", "e2e", "performance", "recovery", "license"}
SECRET_KEY = re.compile(r"(^|_)(password|passwd|token|secret|private_key|api_key|client_secret)($|_)", re.I)
REDACTIONS = [
    re.compile(r"(?i)(authorization:\s*bearer\s+)[^\s]+"),
    re.compile(r"(?i)((?:password|passwd|token|api[_-]?key|client[_-]?secret)\s*[=:]\s*)[^\s]+"),
    re.compile(r"-----BEGIN (?:RSA )?PRIVATE KEY-----.*?-----END (?:RSA )?PRIVATE KEY-----", re.S),
]


def utc_now() -> str:
    return datetime.now(UTC).isoformat().replace("+00:00", "Z")


def redact(text: str, limit: int) -> str:
    value = text[:limit]
    for pattern in REDACTIONS:
        value = pattern.sub(r"\1[REDACTED]" if pattern.groups else "[REDACTED]", value)
    return value + ("\n[TRUNCATED]" if len(text) > limit else "")


def validate(plan: dict[str, Any], root: Path) -> list[str]:
    errors: list[str] = []
    if plan.get("schema_version") != "elite.ci-gates.v1": errors.append("schema_version: unsupported")
    required = plan.get("required_categories", [])
    if not isinstance(required, list) or not required or not set(required) <= CATEGORIES: errors.append("required_categories: invalid")
    limit = plan.get("max_output_bytes")
    if not isinstance(limit, int) or not 4096 <= limit <= 1048576: errors.append("max_output_bytes: must be 4096..1048576")
    gates = plan.get("gates", [])
    ids: set[str] = set()
    covered: set[str] = set()
    for index, gate in enumerate(gates if isinstance(gates, list) else []):
        prefix = f"gates[{index}]"
        gate_id = gate.get("id")
        if not isinstance(gate_id, str) or not re.fullmatch(r"[a-z][a-z0-9-]{1,63}", gate_id) or "replace" in gate_id: errors.append(f"{prefix}.id: invalid")
        elif gate_id in ids: errors.append(f"{prefix}.id: duplicate")
        else: ids.add(gate_id)
        category = gate.get("category")
        if category not in CATEGORIES: errors.append(f"{prefix}.category: invalid")
        elif gate.get("required") is True: covered.add(category)
        argv = gate.get("argv")
        if not isinstance(argv, list) or not argv or not all(isinstance(v, str) and v and "replace-me" not in v for v in argv): errors.append(f"{prefix}.argv: non-empty resolved string array required")
        if "env" in gate or any(SECRET_KEY.search(str(key)) for key in gate): errors.append(f"{prefix}: inline environment/secrets forbidden")
        timeout = gate.get("timeout_seconds")
        if not isinstance(timeout, int) or not 1 <= timeout <= 7200: errors.append(f"{prefix}.timeout_seconds: must be 1..7200")
        cwd = gate.get("cwd", ".")
        try:
            (root / cwd).resolve().relative_to(root.resolve())
        except (ValueError, TypeError):
            errors.append(f"{prefix}.cwd: escapes workspace")
    missing = set(required if isinstance(required, list) else []) - covered
    if missing: errors.append("required_categories: missing required gates: " + ", ".join(sorted(missing)))
    return sorted(errors)


def execute(plan: dict[str, Any], root: Path, plan_hash: str) -> dict[str, Any]:
    started = time.monotonic()
    report: dict[str, Any] = {"schema_version": "elite.ci-evidence.v1", "plan_sha256": plan_hash, "started_at_utc": utc_now(), "workspace": str(root), "gates": []}
    overall = "PASS"
    for gate in plan["gates"]:
        gate_started = time.monotonic()
        result: dict[str, Any] = {"id": gate["id"], "category": gate["category"], "required": gate["required"], "started_at_utc": utc_now()}
        try:
            completed = subprocess.run(gate["argv"], cwd=(root / gate.get("cwd", ".")).resolve(), shell=False, capture_output=True, text=True, encoding="utf-8", errors="replace", timeout=gate["timeout_seconds"], check=False)
            result.update(status="PASS" if completed.returncode == 0 else "FAIL", exit_code=completed.returncode, stdout=redact(completed.stdout, plan["max_output_bytes"]), stderr=redact(completed.stderr, plan["max_output_bytes"]))
        except subprocess.TimeoutExpired as error:
            result.update(status="TIMEOUT", exit_code=None, stdout=redact((error.stdout or "") if isinstance(error.stdout, str) else "", plan["max_output_bytes"]), stderr="gate timed out")
        except OSError as error:
            result.update(status="ERROR", exit_code=None, stdout="", stderr=redact(str(error), plan["max_output_bytes"]))
        result["duration_seconds"] = round(time.monotonic() - gate_started, 3)
        result["finished_at_utc"] = utc_now()
        report["gates"].append(result)
        if gate["required"] and result["status"] != "PASS":
            overall = "FAIL"
            if plan.get("fail_fast", True): break
    report.update(status=overall, duration_seconds=round(time.monotonic() - started, 3), finished_at_utc=utc_now())
    return report


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("plan", type=Path)
    parser.add_argument("--workspace", type=Path, default=Path.cwd())
    parser.add_argument("--evidence", type=Path, required=True)
    args = parser.parse_args()
    raw = args.plan.read_bytes()
    plan = json.loads(raw)
    root = args.workspace.resolve()
    errors = validate(plan, root)
    report = {"schema_version": "elite.ci-evidence.v1", "status": "INVALID", "plan_sha256": hashlib.sha256(raw).hexdigest(), "errors": errors} if errors else execute(plan, root, hashlib.sha256(raw).hexdigest())
    args.evidence.parent.mkdir(parents=True, exist_ok=True)
    args.evidence.write_text(json.dumps(report, indent=2) + "\n", encoding="utf-8")
    print(json.dumps({"status": report["status"], "evidence": str(args.evidence)}))
    return 0 if report["status"] == "PASS" else 2


if __name__ == "__main__":
    sys.exit(main())
