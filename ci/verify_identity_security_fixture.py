#!/usr/bin/env python3
"""Fail-closed LOCAL_FIXTURES verifier for identity/security/authz doors."""

from __future__ import annotations

import argparse
import json
import shutil
import subprocess
import sys
from itertools import permutations
from pathlib import Path
from typing import Any

SCHEMA = "elite.identity-authz-fixture/v1"
REQUIRED_UNAUTHORIZED = {
    "unauthenticated",
    "expired_token",
    "wrong_issuer",
    "wrong_audience",
    "insufficient_role",
    "object_ownership",
}
REQUIRED_SESSION = {
    "logout_revocation",
    "token_revocation",
    "signing_key_rotation",
    "session_expiry",
    "replay_rejection",
}
ALLOWED_FIXTURE_FIELDS = {
    "schema",
    "scope",
    "live_effects",
    "tenants",
    "unauthorized_scenarios",
    "tenant_actions",
    "session_cases",
    "invariants",
    "doors",
}


def load_fixture(path: Path) -> dict[str, Any]:
    data = json.loads(path.read_text(encoding="utf-8"))
    if data.get("schema") != SCHEMA:
        raise ValueError(f"fixture schema must be {SCHEMA}")
    if data.get("scope") != "LOCAL_FIXTURES":
        raise ValueError("fixture scope must be LOCAL_FIXTURES")
    if data.get("live_effects") is not False:
        raise ValueError("fixture live_effects must be false")
    extra = set(data) - ALLOWED_FIXTURE_FIELDS
    if extra:
        raise ValueError(f"fixture has unsupported fields: {', '.join(sorted(extra))}")
    return data


def validate_fixture(data: dict[str, Any]) -> list[str]:
    errors: list[str] = []
    missing = ALLOWED_FIXTURE_FIELDS - set(data)
    if missing:
        errors.append("fixture missing required fields: " + ", ".join(sorted(missing)))
    tenants = data.get("tenants")
    if not isinstance(tenants, list) or len(tenants) < 2 or len(set(tenants)) != len(tenants):
        errors.append("tenants: need at least two unique tenant ids")
    unauthorized = data.get("unauthorized_scenarios")
    if not isinstance(unauthorized, list) or set(unauthorized) != REQUIRED_UNAUTHORIZED:
        errors.append("unauthorized_scenarios: must contain exactly the six mandatory cases")
    actions = data.get("tenant_actions")
    if not isinstance(actions, list) or set(actions) != {"read", "write", "list"}:
        errors.append("tenant_actions: must include read, write, list")
    session = data.get("session_cases")
    if not isinstance(session, list) or set(session) != REQUIRED_SESSION:
        errors.append("session_cases: must include all five session revocation cases")
    invariants = data.get("invariants")
    if not isinstance(invariants, dict):
        errors.append("invariants: object required")
    else:
        for key in ("data_disclosed", "mutation_observed", "old_credential_accepted"):
            if invariants.get(key) is not False:
                errors.append(f"invariants.{key}: must be false for LOCAL_FIXTURES")
    doors = data.get("doors")
    if not isinstance(doors, list) or not doors:
        errors.append("doors: non-empty path list required")
    if isinstance(tenants, list) and isinstance(actions, list) and len(tenants) >= 2:
        expected = len(list(permutations(tenants, 2))) * len(actions)
        if expected < 6:
            errors.append("tenant matrix: cross-tenant combinations too small")
    return errors


def verify_doors(workspace: Path, doors: list[str]) -> list[str]:
    errors: list[str] = []
    for rel in doors:
        if not (workspace / rel).is_file():
            errors.append(f"missing door: {rel}")
    return errors


def run_go_identity_tests(workspace: Path) -> tuple[bool, str]:
    go = shutil.which("go")
    if go is None:
        return True, "go not on PATH; skipped identity unit tests"
    completed = subprocess.run(
        [go, "test", "./internal/platform/identity/...", "-count=1"],
        cwd=workspace,
        capture_output=True,
        text=True,
        check=False,
    )
    if completed.returncode != 0:
        detail = (completed.stderr or completed.stdout or "go test failed").strip()
        return False, detail[:4000]
    return True, "go test ./internal/platform/identity/... PASS"


def main() -> int:
    parser = argparse.ArgumentParser(description="Verify identity/security/authz LOCAL_FIXTURES doors")
    parser.add_argument("--workspace", type=Path, default=Path.cwd())
    parser.add_argument(
        "--fixture",
        type=Path,
        default=Path("ci/fixtures/identity_authz_scenarios.json"),
    )
    parser.add_argument("--skip-go", action="store_true")
    args = parser.parse_args()
    workspace = args.workspace.resolve()
    fixture_path = (workspace / args.fixture).resolve() if not args.fixture.is_absolute() else args.fixture

    errors: list[str] = []
    try:
        fixture = load_fixture(fixture_path)
    except (OSError, json.JSONDecodeError, ValueError) as error:
        print(f"IDENTITY_SECURITY_FIXTURE_REJECTED: {error}", file=sys.stderr)
        return 1

    errors.extend(validate_fixture(fixture))
    errors.extend(verify_doors(workspace, fixture.get("doors", [])))

    if errors:
        for item in errors:
            print(f"IDENTITY_SECURITY_DOOR_FAIL: {item}", file=sys.stderr)
        return 1

    if not args.skip_go:
        ok, message = run_go_identity_tests(workspace)
        print(message)
        if not ok:
            print(f"IDENTITY_SECURITY_GO_TEST_FAIL: {message}", file=sys.stderr)
            return 1

    print(
        json.dumps(
            {
                "schema": "elite.identity-security-verify/v1",
                "scope": "LOCAL_FIXTURES",
                "status": "PASS",
                "fixture": str(fixture_path.relative_to(workspace)),
                "doors_checked": len(fixture.get("doors", [])),
            }
        )
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
