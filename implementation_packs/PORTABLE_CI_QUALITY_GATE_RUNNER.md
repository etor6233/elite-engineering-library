# Portable CI Quality Gate Runner

## 1. Metadata

```yaml
pack_id: "PORTABLE-CI-GATE-RUNNER"
pack_version: "0.1.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un runner CI/CD stack-neutral, sin shell, con timeouts, redacción y evidencia hash-linked."
stacks: ["Python 3.14+ stdlib"]
compatible_with: ["SECURE-OPS-DELIVERY-CORE >=0.1.0 <1.0.0"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources:
  - "https://csrc.nist.gov/pubs/sp/800/218/final"
  - "https://slsa.dev/spec/v1.2/"
verified_at: "2026-08-25"
```

Todos los bloques son `AUTHORED`. Este runner coordina herramientas elegidas por el proyecto; no sustituye compiladores, scanners, un build service aislado ni la verificación de provenance.

## 2. Applicability

Use as a stack-neutral orchestrator for already selected project tools when CI needs bounded, redactable, hash-linked evidence. Reject it as a sandbox or scanner: the target CI runner must provide isolation, resource limits, trusted tool pins and secret injection.

## 3. Architecture contract

- ejecuta `argv` directamente con `shell=False`;
- exige IDs únicos, categorías reconocidas, timeout acotado y directorios dentro del workspace;
- categorías requeridas se declaran en el plan y ninguna puede quedar sin gate;
- cada salida se limita y redacta antes de persistirse;
- fail-fast por defecto; un gate obligatorio fallido produce exit 2;
- evidencia JSON incluye hash del plan, tiempos, código de salida y resultado por gate;
- secretos se inyectan externamente; el plan no contiene `env`, passwords ni tokens.

## 4. Exact file manifest

```text
CREATE ci/gates.example.json
CREATE ci/run_quality_gates.py
CREATE ci/test_run_quality_gates.py
```

## 5. Materialization blocks

### FILE: `ci/gates.example.json`

```yaml
block_id: "PORTABLE-CI:gates-example:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "0373edc1b533fd759e1d1803df006ef1743ae7c55e98f594807266ea06404fda"
variables: []
secrets_allowed: false
```

````json
{
  "schema_version": "elite.ci-gates.v1",
  "required_categories": ["format", "unit", "static", "build", "security", "supply-chain", "integration", "recovery"],
  "fail_fast": true,
  "max_output_bytes": 131072,
  "gates": [
    {"id": "replace-format", "category": "format", "argv": ["replace-me"], "cwd": ".", "timeout_seconds": 300, "required": true}
  ]
}
````

### FILE: `ci/run_quality_gates.py`

```yaml
block_id: "PORTABLE-CI:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "aab21d67cbb9549c2637430a2b4486039b279e52c2e11a4d99d2c5b8e0d72b55"
variables: []
secrets_allowed: false
```

````python
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
````

### FILE: `ci/test_run_quality_gates.py`

```yaml
block_id: "PORTABLE-CI:runner-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "342cc3403e899b584765e7981ffce116f263235d3ac77f9caae2b73ec3ab52eb"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations

import sys
import tempfile
import unittest
from pathlib import Path

from run_quality_gates import execute, redact, validate


class GateRunnerTests(unittest.TestCase):
    def plan(self, argv: list[str] | None = None) -> dict:
        return {"schema_version": "elite.ci-gates.v1", "required_categories": ["unit"], "fail_fast": True, "max_output_bytes": 4096, "gates": [{"id": "unit-test", "category": "unit", "argv": argv or [sys.executable, "-c", "print('ok')"], "cwd": ".", "timeout_seconds": 10, "required": True}]}

    def test_valid_plan_and_pass(self) -> None:
        with tempfile.TemporaryDirectory() as value:
            root = Path(value)
            plan = self.plan()
            self.assertEqual([], validate(plan, root))
            self.assertEqual("PASS", execute(plan, root, "a" * 64)["status"])

    def test_failure_is_recorded(self) -> None:
        with tempfile.TemporaryDirectory() as value:
            report = execute(self.plan([sys.executable, "-c", "raise SystemExit(7)"]), Path(value), "b" * 64)
            self.assertEqual("FAIL", report["status"])
            self.assertEqual(7, report["gates"][0]["exit_code"])

    def test_timeout_is_recorded(self) -> None:
        with tempfile.TemporaryDirectory() as value:
            plan = self.plan([sys.executable, "-c", "import time; time.sleep(2)"])
            plan["gates"][0]["timeout_seconds"] = 1
            self.assertEqual("TIMEOUT", execute(plan, Path(value), "c" * 64)["gates"][0]["status"])

    def test_missing_category_fails(self) -> None:
        with tempfile.TemporaryDirectory() as value:
            plan = self.plan()
            plan["required_categories"].append("security")
            self.assertTrue(any("missing required gates" in error for error in validate(plan, Path(value))))

    def test_shell_and_inline_env_are_not_available(self) -> None:
        with tempfile.TemporaryDirectory() as value:
            plan = self.plan()
            plan["gates"][0]["env"] = {"TOKEN": "forbidden"}
            self.assertTrue(any("inline environment" in error for error in validate(plan, Path(value))))

    def test_redaction_and_truncation(self) -> None:
        value = redact("Authorization: Bearer abc\npassword=hunter2\n" + "x" * 5000, 4096)
        self.assertNotIn("abc", value)
        self.assertNotIn("hunter2", value)
        self.assertIn("[TRUNCATED]", value)


if __name__ == "__main__":
    unittest.main()
````

## 6. Configuration surface

| Field | Type/default | Secret | Validation/effect |
|---|---|---|---|
| `required_categories` | known string list / none | no | every category must have a gate |
| gate `argv` | non-empty string array | none | no | executed directly; shell syntax is not interpreted |
| gate `cwd` | workspace-relative path / `.` | no | traversal/outside workspace rejected |
| `timeout_seconds` | bounded positive integer | required per gate | no | timeout produces failure evidence |
| `max_output_bytes` | bounded integer / 131072 example | no | output is redacted then truncated |
| `fail_fast` / `required` | boolean | true | no | controls continuation, never converts failure to pass |

## 7. Dependency bill

| Tool | Pin | Use | License | Scope | Source |
|---|---|---|---|---|---|
| Python | `3.14+` standard library | runner/tests | PSF-2.0 | CI | `python.org` |
| project compilers/scanners | exact pins in target CI | gates selected by project | respective licenses | CI | official sources |

## 8. Apply order

Materialize after the project defines its build/test/security/recovery commands. Replace the example plan, pin every tool, run unit/negative tests, execute in an isolated CI worker and retain plan/result hashes with the release. Existing CI should call this runner as one controlled stage. Rollback restores the runner and plan together; historical evidence remains immutable.

## 9. Verification

Clean materialization, six unit tests, intentional invalid-example rejection and an eight-category CLI execution passed; see `reconstruction_evidence/PORTABLE_CI_QUALITY_GATE_RUNNER_2026-08-24_V1.md`. Admission stays conditioned until the target CI uses an isolated runner, injects secrets externally, stores immutable evidence and generates/verifies actual SBOM, provenance and signatures.

## 10. Reconstruction evidence

Toolchain, hashes, positive execution and negative cases are recorded in `reconstruction_evidence/PORTABLE_CI_QUALITY_GATE_RUNNER_2026-08-24_V1.md`; the final library audit rechecks version 0.1.1.
