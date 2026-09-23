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
