from __future__ import annotations

import argparse
from datetime import datetime, timezone
import hashlib
import json
from pathlib import Path, PurePosixPath
import re
import sys

SCHEMA = "elite-production-admission/v1"
CONTROLS = {
    "EDGE_CDN_WAF", "IDENTITY_AUTHORIZATION", "PROVIDERS", "POSTGRES_RECOVERY",
    "LOAD_RESILIENCE", "OFFENSIVE_SECURITY", "DEPLOY_ROLLBACK", "BUSINESS_ACCEPTANCE",
}
REQUIRED_ASSERTIONS = {
    "EDGE_CDN_WAF": {"dns_proxy_active", "waf_rulesets_active", "tls_policy_pass", "origin_bypass_blocked", "cache_behavior_pass"},
    "IDENTITY_AUTHORIZATION": {"oidc_conformance_pass", "role_positive_journeys_pass", "unauthorized_negative_journeys_pass", "tenant_isolation_pass", "session_revocation_rotation_pass"},
    "PROVIDERS": {"all_required_sandboxes_pass", "webhook_auth_pass", "idempotency_reconciliation_pass", "rate_limit_retry_pass", "terms_cost_scope_approved"},
    "POSTGRES_RECOVERY": {"pg_verifybackup_pass", "restore_boot_pass", "data_integrity_pass", "pitr_pass", "rpo_met", "rto_met"},
    "LOAD_RESILIENCE": {"k6_thresholds_pass", "expected_workload_represented", "soak_pass", "saturation_recovery_pass", "no_error_budget_breach"},
    "OFFENSIVE_SECURITY": {"zap_automation_pass", "authenticated_scope_pass", "api_scope_pass", "zero_unaccepted_high_critical", "manual_review_pass"},
    "DEPLOY_ROLLBACK": {"immutable_digest_deployed", "canary_health_pass", "rollout_status_pass", "rollback_executed", "previous_digest_restored", "migration_compatibility_pass"},
    "BUSINESS_ACCEPTANCE": {"acceptance_scenarios_pass", "finance_operations_security_approved", "no_open_sev1_sev2", "approvers_distinct", "owner_signoff"},
}
DIGEST = re.compile(r"^sha256:[0-9a-f]{64}$")
HEX = re.compile(r"^[0-9a-f]{64}$")


def utc(value: object, label: str, errors: list[str]) -> datetime | None:
    try:
        if not isinstance(value, str) or not value.endswith("Z"):
            raise ValueError
        return datetime.fromisoformat(value[:-1] + "+00:00")
    except ValueError:
        errors.append(f"{label} must be ISO-8601 UTC ending Z")
        return None


def evidence(root: Path, item: object, label: str, errors: list[str]) -> None:
    if not isinstance(item, dict) or set(item) != {"path", "sha256", "bytes"}:
        errors.append(f"{label} must contain exactly path, sha256 and bytes")
        return
    raw = item.get("path")
    if not isinstance(raw, str) or "\\" in raw:
        errors.append(f"{label}.path must be canonical project-relative")
        return
    rel = PurePosixPath(raw)
    if rel.is_absolute() or any(part in {"", ".", ".."} for part in rel.parts):
        errors.append(f"{label}.path must be canonical project-relative")
        return
    candidate = root.joinpath(*rel.parts)
    try:
        resolved = candidate.resolve(strict=True)
        resolved.relative_to(root)
    except (FileNotFoundError, ValueError):
        errors.append(f"{label}.path is missing or escapes project root")
        return
    if candidate.is_symlink() or not resolved.is_file():
        errors.append(f"{label}.path must be a regular non-symlink file")
        return
    data = resolved.read_bytes()
    if item.get("bytes") != len(data) or not HEX.fullmatch(str(item.get("sha256", ""))) or item["sha256"] != hashlib.sha256(data).hexdigest():
        errors.append(f"{label} bytes/SHA-256 mismatch")


def validate(record: object, root: Path, now: datetime | None = None) -> list[str]:
    errors: list[str] = []
    if not isinstance(record, dict):
        return ["record must be an object"]
    if set(record) != {"schema", "project_id", "environment", "release_digest", "evaluated_at", "controls"}:
        errors.append("record fields must be exact")
    if record.get("schema") != SCHEMA:
        errors.append(f"schema must be {SCHEMA}")
    if not isinstance(record.get("project_id"), str) or not record["project_id"].strip():
        errors.append("project_id is required")
    if record.get("environment") != "production":
        errors.append("environment must be production")
    if not DIGEST.fullmatch(str(record.get("release_digest", ""))) or set(str(record.get("release_digest", "")).split(":")[-1]) == {"0"}:
        errors.append("release_digest must be a real sha256 digest")
    evaluated = utc(record.get("evaluated_at"), "evaluated_at", errors)
    clock = now or datetime.now(timezone.utc)
    if evaluated and (evaluated > clock or (clock - evaluated).total_seconds() > 3600):
        errors.append("evaluated_at must be within the last hour")
    controls = record.get("controls")
    if not isinstance(controls, list):
        return sorted(set(errors + ["controls must be an array"]))
    ids = [item.get("id") for item in controls if isinstance(item, dict)]
    if len(controls) != 8 or set(ids) != CONTROLS or len(ids) != len(set(ids)):
        errors.append("controls must contain the eight unique production controls")
    for index, item in enumerate(controls):
        label = f"controls[{index}]"
        if not isinstance(item, dict):
            errors.append(f"{label} must be an object")
            continue
        required = {"id", "project_id", "environment", "release_digest", "result", "executed_at", "expires_at", "executor_ref", "tool", "target", "assertions", "evidence"}
        if set(item) != required:
            errors.append(f"{label} fields must be exact")
        for field in ("project_id", "environment", "release_digest"):
            if item.get(field) != record.get(field):
                errors.append(f"{label}.{field} must match admission record")
        if item.get("result") != "PASS":
            errors.append(f"{label}.result must be PASS")
        start = utc(item.get("executed_at"), f"{label}.executed_at", errors)
        end = utc(item.get("expires_at"), f"{label}.expires_at", errors)
        if start and end and not (start <= clock <= end and (end - start).total_seconds() <= 30 * 86400):
            errors.append(f"{label} evidence is stale, future or valid for more than 30 days")
        if not isinstance(item.get("executor_ref"), str) or not item["executor_ref"].strip():
            errors.append(f"{label}.executor_ref is required")
        tool = item.get("tool")
        if not isinstance(tool, dict) or set(tool) != {"name", "version", "source", "digest"} or not all(isinstance(tool.get(k), str) and tool[k].strip() for k in tool):
            errors.append(f"{label}.tool identity is incomplete")
        elif not DIGEST.fullmatch(tool["digest"]):
            errors.append(f"{label}.tool.digest must be sha256")
        if not isinstance(item.get("target"), dict) or not item["target"]:
            errors.append(f"{label}.target must identify the real target")
        assertions = item.get("assertions")
        expected_assertions = REQUIRED_ASSERTIONS.get(str(item.get("id")), set())
        if not isinstance(assertions, dict) or set(assertions) != expected_assertions or any(value is not True for value in assertions.values()):
            errors.append(f"{label}.assertions must contain exactly the required true semantic assertions")
        files = item.get("evidence")
        if not isinstance(files, list) or not files:
            errors.append(f"{label}.evidence must be non-empty")
        else:
            for number, proof in enumerate(files):
                evidence(root, proof, f"{label}.evidence[{number}]", errors)
    return sorted(set(errors))


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--project-root", required=True, type=Path)
    parser.add_argument("--record", default="PRODUCTION_ADMISSION.json")
    parser.add_argument("--report", default="PRODUCTION_ADMISSION_REPORT.json")
    args = parser.parse_args()
    root = args.project_root.resolve(strict=True)
    record_path = (root / args.record).resolve(strict=True)
    if record_path.parent != root or record_path.is_symlink():
        raise SystemExit("record must be a regular file directly under project root")
    record = json.loads(record_path.read_text(encoding="utf-8"))
    errors = validate(record, root)
    report = {"schema": "elite-production-admission-report/v1", "status": "PRODUCTION_ADMITTED" if not errors else "BLOCKED", "errors": errors}
    output = root / args.report
    if output.exists():
        raise SystemExit("report already exists")
    output.write_text(json.dumps(report, indent=2, sort_keys=True) + "\n", encoding="utf-8", newline="\n")
    print(json.dumps(report, indent=2, sort_keys=True))
    return 0 if not errors else 2


if __name__ == "__main__":
    sys.exit(main())
