# Dependency License Evidence Core

## 1. Metadata

```yaml
pack_id: "DEPENDENCY-LICENSE-EVIDENCE-CORE"
pack_version: "0.1.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un gate neutral que valida política de licencias y genera evidencia determinista con hashes de notices para Go y pnpm."
stacks: ["Python 3.14", "go-licenses 2.0.1", "pnpm 11/12"]
compatible_with: ["PORTABLE-CI-GATE-RUNNER 0.1.x", "SECURE-OPS-DELIVERY-CORE 0.1.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://github.com/google/go-licenses", "https://pnpm.io/cli/licenses", "https://cyclonedx.org/docs/1.7/json/"]
verified_at: "2026-08-25"
```

El gate no toma decisiones legales. Exige una política aprobada por el proyecto, rechaza expresiones denegadas o no revisadas y prueba que los artefactos de licencia existen y tienen hashes. Para Go se alimenta con `go-licenses report` y el directorio de `go-licenses save`; para pnpm con `pnpm licenses list --prod --json --long`.

## 2. Applicability

Use this pack for every distributable Go or pnpm artifact that needs machine-readable dependency evidence, policy decisions and notices. It does not provide legal advice or preapprove an unknown license; unresolved, denied or unclassified expressions fail closed.

## 3. Architecture contract

Evidence is generated from the exact build/install graph, normalized to explicit package/version/license records, checked against repository-owned policy and stored beside the immutable release. Policy acknowledgement is an audited exception, not a string bypass. No dependency scanner writes source code or receives product secrets.

## 4. Exact file manifest

```text
CREATE supply_chain/license_gate.py
CREATE supply_chain/test_license_gate.py
CREATE supply_chain/license-policy.example.json
```

## 5. Materialization blocks

### FILE: `supply_chain/license_gate.py`
```yaml
block_id: "DEPENDENCY-LICENSE-EVIDENCE:gate:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "9f9390357f03ae0da1f7f803c165113bdc7194012a84e9f1c9b6d91b88ae3e53"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
import csv
import hashlib
import json
from pathlib import Path
from typing import Any


class GateError(Exception):
    pass


def sha256(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as stream:
        for chunk in iter(lambda: stream.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def admitted(expression: str, policy: dict[str, Any]) -> str:
    if expression in policy.get("denied_expressions", []):
        raise GateError(f"denied license expression: {expression}")
    if expression in policy.get("allowed_expressions", []):
        return "allowed"
    acknowledgement = policy.get("acknowledged_expressions", {}).get(expression)
    if not isinstance(acknowledgement, str) or len(acknowledgement.strip()) < 8:
        raise GateError(f"unreviewed license expression: {expression}")
    return "acknowledged"


def license_files(package: Path) -> list[Path]:
    names = ("LICENSE*", "LICENCE*", "COPYING*", "NOTICE*")
    found = {candidate for pattern in names for candidate in package.glob(pattern) if candidate.is_file()}
    return sorted(found, key=lambda value: value.name.casefold())


def pnpm_components(report_path: Path, policy: dict[str, Any]) -> list[dict[str, Any]]:
    report = json.loads(report_path.read_text(encoding="utf-8"))
    if not isinstance(report, dict):
        raise GateError("pnpm report must be an object")
    components: list[dict[str, Any]] = []
    for expression, packages in report.items():
        status = admitted(expression, policy)
        if not isinstance(packages, list):
            raise GateError("pnpm license group must be an array")
        for package in packages:
            name, versions, paths = package.get("name"), package.get("versions"), package.get("paths")
            if not isinstance(name, str) or not isinstance(versions, list) or not versions or not isinstance(paths, list) or not paths:
                raise GateError("invalid pnpm package record")
            artifacts = license_files(Path(paths[0]))
            if not artifacts:
                raise GateError(f"license artifact missing for pnpm package: {name}")
            components.append({"ecosystem": "npm", "name": name, "versions": sorted(str(value) for value in versions), "license": expression, "admission": status, "artifacts": [{"name": item.name, "sha256": sha256(item)} for item in artifacts]})
    return components


def go_components(report_path: Path, notices: Path, policy: dict[str, Any]) -> list[dict[str, Any]]:
    artifacts = [item for item in notices.rglob("*") if item.is_file()]
    if not artifacts:
        raise GateError("go-licenses save produced no notice artifacts")
    rows: list[dict[str, Any]] = []
    with report_path.open("r", encoding="utf-8", newline="") as stream:
        for row in csv.reader(stream):
            if len(row) != 3 or not all(row):
                raise GateError("invalid go-licenses CSV row")
            package, url, expression = row
            rows.append({"ecosystem": "go", "name": package, "license": expression, "license_url": url, "admission": admitted(expression, policy)})
    if not rows:
        raise GateError("go-licenses report is empty")
    notice_hashes = [{"relative_path": item.relative_to(notices).as_posix(), "sha256": sha256(item)} for item in sorted(artifacts)]
    return rows + [{"ecosystem": "go-notice-bundle", "artifacts": notice_hashes}]


def run(args: argparse.Namespace) -> dict[str, Any]:
    policy = json.loads(Path(args.policy).read_text(encoding="utf-8"))
    components: list[dict[str, Any]] = []
    if args.pnpm_report:
        components.extend(pnpm_components(Path(args.pnpm_report), policy))
    if args.go_report or args.go_notices:
        if not args.go_report or not args.go_notices:
            raise GateError("go report and saved notices must be supplied together")
        components.extend(go_components(Path(args.go_report), Path(args.go_notices), policy))
    if not components:
        raise GateError("at least one ecosystem report is required")
    result = {"schema": "elite-license-evidence/v1", "status": "PASS", "components": sorted(components, key=lambda item: (item["ecosystem"], item.get("name", "")))}
    output = Path(args.output)
    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text(json.dumps(result, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    return result


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--policy", required=True)
    parser.add_argument("--output", required=True)
    parser.add_argument("--pnpm-report")
    parser.add_argument("--go-report")
    parser.add_argument("--go-notices")
    args = parser.parse_args()
    try:
        result = run(args)
    except (GateError, OSError, ValueError, json.JSONDecodeError) as error:
        print(f"LICENSE_GATE_FAILED: {error}")
        return 2
    print(f"LICENSE_GATE_PASS components={len(result['components'])}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `supply_chain/test_license_gate.py`
```yaml
block_id: "DEPENDENCY-LICENSE-EVIDENCE:test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "3f282227fd33f592434c2433ced21adfb1e706ef27f92eaf681b85c2b465aa90"
variables: []
secrets_allowed: false
```
````python
import json
import tempfile
import unittest
from argparse import Namespace
from pathlib import Path

from license_gate import GateError, run


class LicenseGateTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.root = Path(self.temp.name)
        self.policy = self.root / "policy.json"
        self.policy.write_text(json.dumps({"allowed_expressions": ["MIT", "Apache-2.0"], "acknowledged_expressions": {}, "denied_expressions": ["AGPL-3.0-only"]}), encoding="utf-8")

    def tearDown(self):
        self.temp.cleanup()

    def test_generates_hashed_pnpm_and_go_evidence(self):
        package = self.root / "package"
        package.mkdir()
        (package / "LICENSE").write_text("license text", encoding="utf-8")
        pnpm = self.root / "pnpm.json"
        pnpm.write_text(json.dumps({"MIT": [{"name": "example", "versions": ["1.0.0"], "paths": [str(package)]}]}), encoding="utf-8")
        notices = self.root / "notices"
        notices.mkdir()
        (notices / "LICENSE").write_text("go license", encoding="utf-8")
        go_report = self.root / "go.csv"
        go_report.write_text("example.org/module,https://example.org/LICENSE,Apache-2.0\n", encoding="utf-8")
        output = self.root / "evidence.json"
        result = run(Namespace(policy=str(self.policy), output=str(output), pnpm_report=str(pnpm), go_report=str(go_report), go_notices=str(notices)))
        self.assertEqual(result["status"], "PASS")
        self.assertTrue(output.exists())
        self.assertIn("sha256", output.read_text(encoding="utf-8"))

    def test_rejects_unreviewed_denied_and_missing_artifacts(self):
        package = self.root / "package"
        package.mkdir()
        report = self.root / "pnpm.json"
        for expression in ("LGPL-3.0-or-later", "AGPL-3.0-only"):
            report.write_text(json.dumps({expression: [{"name": "example", "versions": ["1"], "paths": [str(package)]}]}), encoding="utf-8")
            with self.assertRaises(GateError):
                run(Namespace(policy=str(self.policy), output=str(self.root / "out.json"), pnpm_report=str(report), go_report=None, go_notices=None))
        report.write_text(json.dumps({"MIT": [{"name": "example", "versions": ["1"], "paths": [str(package)]}]}), encoding="utf-8")
        with self.assertRaises(GateError):
            run(Namespace(policy=str(self.policy), output=str(self.root / "out.json"), pnpm_report=str(report), go_report=None, go_notices=None))


if __name__ == "__main__":
    unittest.main()
````

### FILE: `supply_chain/license-policy.example.json`
```yaml
block_id: "DEPENDENCY-LICENSE-EVIDENCE:policy:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "06c632fc1a8e5f2ef224a195783289c78fff9f5c4172d6576cac5f484b4a9943"
variables: []
secrets_allowed: false
```
````json
{
  "allowed_expressions": [
    "0BSD",
    "Apache-2.0",
    "BSD-3-Clause",
    "CC-BY-4.0",
    "ISC",
    "MIT"
  ],
  "acknowledged_expressions": {},
  "denied_expressions": [
    "AGPL-3.0-only",
    "AGPL-3.0-or-later"
  ]
}
````

## 6. Configuration surface

| Input | Type | Safe default | Secret | Validation/effect |
|---|---|---|---|---|
| allow/deny policy | versioned JSON | restrictive example | no | unknown expressions require explicit acknowledgement |
| artifact/dependency graph | generated evidence | none | no | must match the shipped target and lock/module files |
| output directory | relative path | project evidence directory | no | must remain inside the workspace |

## 7. Dependency bill

| Tool | Pin policy | Use | License | Scope | Official source |
|---|---|---|---|---|---|
| `go-licenses` | exact CI pin selected by project | Go report/check/save | Apache-2.0 | build | GitHub/google |
| pnpm | exact toolchain pin selected by project | JavaScript license graph | MIT | build | pnpm.io |
| CycloneDX schema | `1.7` | SBOM contract | specification terms | evidence | cyclonedx.org |
| Python | `3.14+` stdlib | normalization/policy tests | PSF-2.0 | verification | python.org |

## 8. Apply order

Materialize after product lockfiles and build targets exist. Pin the scanners, generate evidence from the exact distributable, normalize, apply policy, preserve license texts/notices and bind hashes to the release record. For an existing workspace, merge policy deliberately. Rollback restores the prior policy and evidence only together with the corresponding prior artifact.

## 9. Verification

1. Fijar versiones de `go-licenses`, Go y pnpm en CI; no usar `latest`.
2. Ejecutar report/check/save sobre el ejecutable realmente distribuido e incluir tests cuando sus dependencias se redistribuyan.
3. Ejecutar pnpm sobre la instalación productiva exacta; no confundir lockfile con artefacto desplegado.
4. Una expresión fuera de `allowed_expressions` requiere una entrada trazable en `acknowledged_expressions`; este ejemplo no preaprueba copyleft.
5. Conservar el JSON de evidencia, notices y sus hashes junto al SBOM/artefacto.
6. Repetir el gate si cambia un lockfile, módulo, target OS/arch o modo de packaging.

## 10. Reconstruction evidence

Clean reconstruction and positive/negative policy cases are recorded in `reconstruction_evidence/DEPENDENCY_LICENSE_EVIDENCE_2026-08-24_V1.md`; the final full-library audit verifies current embedded hashes and manifest parity.
