from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path
from typing import Any

SECRET_KEYS = re.compile(r"(^|_)(password|passwd|token|secret|private_key|api_key|client_secret)($|_)", re.I)
SECRET_MARKERS = ("-----BEGIN PRIVATE KEY-----", "-----BEGIN RSA PRIVATE KEY-----", "AKIA")
SECRET_REF = re.compile(r"^(secret-manager|vault|workload-identity|kms)://[^\s]+$")
DIGEST = re.compile(r"^sha256:[0-9a-f]{64}$")


def require(condition: bool, message: str, errors: list[str]) -> None:
    if not condition:
        errors.append(message)


def reject_embedded_secrets(value: Any, path: str, errors: list[str]) -> None:
    if isinstance(value, dict):
        for key, child in value.items():
            child_path = f"{path}.{key}"
            if SECRET_KEYS.search(key) and key not in {"secret_scan_evidence", "secrets"}:
                errors.append(f"{child_path}: secret-shaped key is forbidden")
            reject_embedded_secrets(child, child_path, errors)
    elif isinstance(value, list):
        for index, child in enumerate(value):
            reject_embedded_secrets(child, f"{path}[{index}]", errors)
    elif isinstance(value, str) and any(marker in value for marker in SECRET_MARKERS):
        errors.append(f"{path}: embedded credential marker is forbidden")


def nonempty_string(value: Any) -> bool:
    return isinstance(value, str) and bool(value.strip()) and "replace-me" not in value


def validate(document: dict[str, Any], require_evidence_files: bool = False, root: Path | None = None) -> list[str]:
    errors: list[str] = []
    require(document.get("schema_version") == "elite.operational-readiness.v1", "schema_version: unsupported", errors)
    reject_embedded_secrets(document, "$", errors)

    service = document.get("service", {})
    require(nonempty_string(service.get("name")), "service.name: unresolved", errors)
    require(nonempty_string(service.get("owner")), "service.owner: unresolved", errors)
    require(service.get("tier") in {1, 2, 3}, "service.tier: must be 1, 2 or 3", errors)
    require(set(document.get("environments", [])) >= {"staging", "production"}, "environments: staging and production required", errors)

    secrets = document.get("secrets", {})
    require(nonempty_string(secrets.get("provider")), "secrets.provider: unresolved", errors)
    refs = secrets.get("references", [])
    require(bool(refs) and all(isinstance(ref, str) and SECRET_REF.fullmatch(ref) for ref in refs), "secrets.references: invalid secret reference", errors)
    require(secrets.get("workload_identity") is True, "secrets.workload_identity: must be true", errors)
    require(isinstance(secrets.get("rotation_days"), int) and 1 <= secrets["rotation_days"] <= 365, "secrets.rotation_days: must be 1..365", errors)
    require(nonempty_string(secrets.get("break_glass_runbook")), "secrets.break_glass_runbook: unresolved", errors)

    telemetry = document.get("telemetry", {})
    require(set(telemetry.get("signals", [])) == {"logs", "metrics", "traces"}, "telemetry.signals: logs, metrics and traces required", errors)
    require(telemetry.get("protocol") == "otlp", "telemetry.protocol: must be otlp", errors)
    for key in ("backend", "redaction_test", "retention_policy"):
        require(nonempty_string(telemetry.get(key)), f"telemetry.{key}: unresolved", errors)
    require({"service.name", "deployment.environment.name", "trace_id"} <= set(telemetry.get("correlation_fields", [])), "telemetry.correlation_fields: required fields missing", errors)

    slos = document.get("slos", [])
    require(isinstance(slos, list) and len(slos) > 0, "slos: at least one required", errors)
    for index, slo in enumerate(slos if isinstance(slos, list) else []):
        prefix = f"slos[{index}]"
        for key in ("name", "sli", "owner", "runbook"):
            require(nonempty_string(slo.get(key)), f"{prefix}.{key}: unresolved", errors)
        objective = slo.get("objective")
        require(isinstance(objective, (int, float)) and 0.9 <= objective < 1, f"{prefix}.objective: must be >=0.9 and <1", errors)
        require(slo.get("window_days") in {7, 28, 30, 90}, f"{prefix}.window_days: unsupported", errors)
        require(len(slo.get("page_burn_rates", [])) >= 2, f"{prefix}.page_burn_rates: two windows required", errors)
        require(len(slo.get("ticket_burn_rates", [])) >= 2, f"{prefix}.ticket_burn_rates: two windows required", errors)

    chain = document.get("supply_chain", {})
    require(bool(chain.get("lockfiles")) and all(nonempty_string(v) for v in chain.get("lockfiles", [])), "supply_chain.lockfiles: unresolved", errors)
    require(chain.get("sbom", {}).get("format") == "CycloneDX" and chain.get("sbom", {}).get("spec_version") == "1.7", "supply_chain.sbom: CycloneDX 1.7 required", errors)
    require(chain.get("provenance", {}).get("format") == "SLSA" and chain.get("provenance", {}).get("spec_version") == "1.2", "supply_chain.provenance: SLSA 1.2 required", errors)
    for artifact in ("sbom", "provenance"):
        require(nonempty_string(chain.get(artifact, {}).get("artifact")), f"supply_chain.{artifact}.artifact: evidence reference required", errors)
    policy = chain.get("vulnerability_policy", {})
    require(policy.get("max_critical") == 0, "supply_chain.vulnerability_policy.max_critical: must be zero", errors)
    for key in ("license_policy", "secret_scan_evidence"):
        require(nonempty_string(chain.get(key)), f"supply_chain.{key}: unresolved", errors)

    release = document.get("release", {})
    require(bool(DIGEST.fullmatch(str(release.get("artifact_digest", "")))) and set(str(release.get("artifact_digest", "")).split(":", 1)[1]) != {"0"}, "release.artifact_digest: real immutable digest required", errors)
    require(release.get("deployment_strategy") in {"canary", "blue-green", "rolling"}, "release.deployment_strategy: unsupported", errors)
    require(release.get("migration_strategy") == "expand-migrate-contract", "release.migration_strategy: expand-migrate-contract required", errors)
    require(len(release.get("health_gates", [])) >= 3, "release.health_gates: at least three required", errors)
    rollback = release.get("rollback_command", [])
    require(isinstance(rollback, list) and len(rollback) >= 2 and all(nonempty_string(v) for v in rollback), "release.rollback_command: argv array required", errors)
    for key in ("rollback_test", "approval_record"):
        require(nonempty_string(release.get(key)), f"release.{key}: unresolved", errors)

    recovery = document.get("recovery", {})
    require(isinstance(recovery.get("rpo_minutes"), int) and recovery["rpo_minutes"] > 0, "recovery.rpo_minutes: positive integer required", errors)
    require(isinstance(recovery.get("rto_minutes"), int) and recovery["rto_minutes"] > 0, "recovery.rto_minutes: positive integer required", errors)
    for key in ("backup_policy", "restore_evidence", "last_drill_utc"):
        require(nonempty_string(recovery.get(key)), f"recovery.{key}: unresolved", errors)
    require(len(recovery.get("dependency_order", [])) >= 3, "recovery.dependency_order: incomplete", errors)

    if require_evidence_files:
        base = root or Path.cwd()
        evidence_paths = {
            "telemetry.redaction_test": telemetry.get("redaction_test"),
            "supply_chain.sbom.artifact": chain.get("sbom", {}).get("artifact"),
            "supply_chain.provenance.artifact": chain.get("provenance", {}).get("artifact"),
            "supply_chain.secret_scan_evidence": chain.get("secret_scan_evidence"),
            "release.rollback_test": release.get("rollback_test"),
            "release.approval_record": release.get("approval_record"),
            "recovery.restore_evidence": recovery.get("restore_evidence"),
        }
        for field, item in evidence_paths.items():
            if not nonempty_string(item):
                errors.append(f"{field}: evidence reference required")
                continue
            try:
                candidate = (base / item).resolve()
                try:
                    candidate.relative_to(base.resolve())
                except ValueError:
                    errors.append(f"{field}: evidence path escapes root")
                    continue
                require(candidate.is_file(), f"{field}: evidence file not found", errors)
            except (OSError, ValueError, RuntimeError):
                errors.append(f"{field}: cannot inspect file")
    return sorted(set(errors))


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("document", type=Path)
    parser.add_argument("--require-evidence-files", action="store_true")
    parser.add_argument("--root", type=Path, default=Path.cwd())
    args = parser.parse_args()
    document = json.loads(args.document.read_text(encoding="utf-8"))
    errors = validate(document, args.require_evidence_files, args.root)
    print(json.dumps({"status": "PASS" if not errors else "FAIL", "errors": errors}, indent=2))
    return 0 if not errors else 2


if __name__ == "__main__":
    sys.exit(main())
