# Capability Gap Resolution Gate

## 1. Metadata

```yaml
pack_id: "CAPABILITY-GAP-RESOLUTION-GATE"
pack_version: "0.1.0"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un gate fail-closed que convierte una capability requerida sin cobertura en investigación oficial trazable, admisión honesta o bloqueo explícito, sin inventar código ni procedencia."
stacks: ["CPython 3.11+ stdlib"]
compatible_with: ["proyectos nuevos", "sistemas existentes", "cualquier stack de producto"]
incompatible_with: ["investigación sin evidencia", "branches móviles", "procedencia inferida", "autoaprobación comercial o legal"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://csrc.nist.gov/pubs/sp/800/218/final", "https://sre.google/sre-book/reliable-product-launches/", "https://www.nasa.gov/reference/5-0-product-realization/", "https://slsa.dev/spec/v1.0/levels"]
verified_at: "2026-09-05"
```

## 2. Applicability

Use cuando una capability `REQUIRED` no tenga un pack admitido compatible o cuando una fuente previa esté obsoleta, rechazada o incompleta. No use para eludir una decisión humana, convertir documentación en código oficial ni declarar que se buscó “todo Internet”. El alcance es registrar superficies oficiales pertinentes, candidatos observados, gates, agotamiento y trigger de reapertura.

## 3. Architecture contract

- El agente investiga; el gate valida el expediente y emite un receipt inmutable enlazado por hash.
- `RESEARCH_INCOMPLETE` falla para forzar continuidad. Sólo `USE_REUSABLE_PACK` produce `implementation_ready=true`.
- `USE_CONDITIONED_PACK` conserva blockers; `NO_ADMISSIBLE_SOURCE`, `ACCESS_BLOCKED` y `HUMAN_DECISION_REQUIRED` informan el límite sin inventar una solución.
- Código oficial exige URL oficial observada, revisión inmutable, SHA-256, licencia, claim estrecho y G0–G8. Código local conserva `AUTHORED`, autoridades rectoras y búsquedas que demostraron el hueco.
- La frescura se revalida al uso; un receipt no sustituye source identity, build, test, SCA, journey, despliegue o aceptación del target.
- El registro se devuelve a mapas, locks, catálogo, pack o ledger canónicos para evitar reinvestigación y divergencia.

## 4. Exact file manifest

```text
CREATE capability_gap_resolution/profile.example.json
CREATE capability_gap_resolution/resolve_capability_gap.py
CREATE capability_gap_resolution/verify_pack.py
CREATE capability_gap_resolution/README.md
```

## 5. Materialization blocks

### FILE: `capability_gap_resolution/profile.example.json`
```yaml
block_id: "CAPABILITY-GAP-RESOLUTION-GATE:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local, governed by NIST SSDF, Google SRE launch practice, NASA V&V and SLSA provenance"
license: "LicenseRef-Workspace-Owner"
sha256: "5ed2847511fd8c3b1864cba5f403bb46d9aa2822b7681a8c6756583206270114"
variables: []
secrets_allowed: false
```
````json
{
  "schema_version": 1,
  "capability_id": "EXAMPLE-CAPABILITY",
  "requirement_ref": "PROJECT_BLUEPRINT.md#example-capability",
  "observed_at": "2099-01-01T00:00:00Z",
  "max_age_days": 30,
  "authority_refs": [
    "PROJECT_AUTHORITY_MAP.md#example-capability"
  ],
  "source_searches": [
    {
      "search_id": "official-owner-docs",
      "query": "exact capability official implementation",
      "official_domains": ["example.com"],
      "executed_at": "2099-01-01T00:00:00Z",
      "evidence_refs": ["evidence/example-search.md"]
    }
  ],
  "candidates": [],
  "decision": {
    "state": "RESEARCH_INCOMPLETE",
    "selected_candidate_id": null,
    "reason": "Template only; perform and record current official research.",
    "canonical_updates": [],
    "blockers": ["No candidate has been audited."],
    "exhaustion": null
  }
}
````

### FILE: `capability_gap_resolution/resolve_capability_gap.py`
```yaml
block_id: "CAPABILITY-GAP-RESOLUTION-GATE:resolver:v1"
operation: CREATE
provenance: AUTHORED
source: "local, governed by NIST SSDF, Google SRE launch practice, NASA V&V and SLSA provenance"
license: "LicenseRef-Workspace-Owner"
sha256: "dbe60f23f6fba9e03ac1f27b5815dbbfafd79cb5360fb53bf702f4b754e58015"
variables: []
secrets_allowed: false
```
````python
#!/usr/bin/env python3
"""Validate a capability-gap investigation and emit an immutable decision receipt."""

from __future__ import annotations

import argparse
import datetime as dt
import hashlib
import json
import os
import pathlib
import re
import sys
import tempfile
from urllib.parse import urlparse


class GapError(ValueError):
    pass


def fail(message: str) -> None:
    raise GapError(message)


def no_duplicates(pairs):
    result = {}
    for key, value in pairs:
        if key in result:
            fail(f"duplicate JSON key: {key}")
        result[key] = value
    return result


def read_json(path: pathlib.Path):
    try:
        return json.loads(path.read_text(encoding="utf-8"), object_pairs_hook=no_duplicates)
    except (OSError, UnicodeError, json.JSONDecodeError) as exc:
        fail(f"record is not strict UTF-8 JSON: {exc}")


def require_string(value, name: str) -> str:
    if not isinstance(value, str) or not value.strip():
        fail(f"{name} must be a non-empty string")
    return value.strip()


def require_list(value, name: str, *, nonempty: bool = False) -> list:
    if not isinstance(value, list) or (nonempty and not value):
        fail(f"{name} must be {'a non-empty ' if nonempty else ''}array")
    return value


def parse_time(value: str, name: str) -> dt.datetime:
    try:
        parsed = dt.datetime.fromisoformat(value.replace("Z", "+00:00"))
    except ValueError:
        fail(f"{name} must be ISO-8601 with timezone")
    if parsed.tzinfo is None:
        fail(f"{name} must include timezone")
    return parsed.astimezone(dt.timezone.utc)


def safe_ref(value, name: str) -> str:
    ref = require_string(value, name)
    if "\x00" in ref or ref.startswith(("file:", "javascript:")):
        fail(f"{name} contains an unsafe reference")
    return ref


def verify_ref(root: pathlib.Path, value, name: str) -> str:
    ref = safe_ref(value, name)
    if ref.startswith("https://"):
        fail(f"{name} must be preserved as a local evidence file; keep the official URL in source_url")
    raw_path = ref.split("#", 1)[0]
    target = (root / raw_path).resolve()
    try:
        target.relative_to(root)
    except ValueError:
        fail(f"{name} escapes project root")
    if not target.is_file():
        fail(f"{name} does not resolve to an evidence file: {raw_path}")
    return ref


def https_host(value, name: str) -> str:
    url = require_string(value, name)
    parsed = urlparse(url)
    if parsed.scheme != "https" or not parsed.hostname or parsed.username or parsed.password:
        fail(f"{name} must be an HTTPS URL without embedded credentials")
    return parsed.hostname.lower().rstrip(".")


def host_allowed(host: str, domains: set[str]) -> bool:
    return any(host == domain or host.endswith("." + domain) for domain in domains)


GATES = tuple(f"G{i}" for i in range(9))
GATE_VALUES = {"PASS", "FAIL", "NOT_APPLICABLE"}
TERMINAL_STATES = {
    "USE_REUSABLE_PACK",
    "USE_CONDITIONED_PACK",
    "NO_ADMISSIBLE_SOURCE",
    "ACCESS_BLOCKED",
    "HUMAN_DECISION_REQUIRED",
}


def validate(record: dict, as_of: dt.datetime, project_root: pathlib.Path) -> dict:
    if not isinstance(record, dict) or record.get("schema_version") != 1:
        fail("schema_version must equal 1")
    capability = require_string(record.get("capability_id"), "capability_id")
    if not re.fullmatch(r"[A-Z0-9][A-Z0-9._-]{2,79}", capability):
        fail("capability_id must be a stable uppercase identifier")
    requirement_ref = verify_ref(project_root, record.get("requirement_ref"), "requirement_ref")
    observed = parse_time(require_string(record.get("observed_at"), "observed_at"), "observed_at")
    max_age = record.get("max_age_days")
    if type(max_age) is not int or not 1 <= max_age <= 365:
        fail("max_age_days must be an integer from 1 to 365")
    if observed > as_of + dt.timedelta(minutes=5):
        fail("observed_at is in the future")
    if as_of - observed > dt.timedelta(days=max_age):
        fail("official research is stale; repeat it before deciding")
    authorities = [verify_ref(project_root, v, "authority_refs[]") for v in require_list(record.get("authority_refs"), "authority_refs", nonempty=True)]

    searches = require_list(record.get("source_searches"), "source_searches", nonempty=True)
    search_ids: set[str] = set()
    official_domains: set[str] = set()
    for index, search in enumerate(searches):
        if not isinstance(search, dict):
            fail(f"source_searches[{index}] must be an object")
        search_id = require_string(search.get("search_id"), f"source_searches[{index}].search_id")
        if search_id in search_ids:
            fail(f"duplicate search_id: {search_id}")
        search_ids.add(search_id)
        require_string(search.get("query"), f"source_searches[{index}].query")
        when = parse_time(require_string(search.get("executed_at"), f"source_searches[{index}].executed_at"), f"source_searches[{index}].executed_at")
        if when > observed + dt.timedelta(minutes=5) or observed - when > dt.timedelta(days=max_age):
            fail(f"source_searches[{index}] is outside the freshness window")
        for domain in require_list(search.get("official_domains"), f"source_searches[{index}].official_domains", nonempty=True):
            normalized = require_string(domain, "official_domains[]").lower().rstrip(".")
            if not re.fullmatch(r"[a-z0-9.-]+\.[a-z]{2,}", normalized):
                fail(f"invalid official domain: {normalized}")
            official_domains.add(normalized)
        for evidence in require_list(search.get("evidence_refs"), f"source_searches[{index}].evidence_refs", nonempty=True):
            verify_ref(project_root, evidence, "evidence_refs[]")

    candidates = require_list(record.get("candidates"), "candidates")
    by_id: dict[str, dict] = {}
    for index, candidate in enumerate(candidates):
        if not isinstance(candidate, dict):
            fail(f"candidates[{index}] must be an object")
        candidate_id = require_string(candidate.get("candidate_id"), f"candidates[{index}].candidate_id")
        if candidate_id in by_id:
            fail(f"duplicate candidate_id: {candidate_id}")
        origin = candidate.get("origin_kind")
        if origin not in {"OFFICIAL_CODE", "OFFICIAL_ARCHITECTURE", "LOCAL_AUTHORED"}:
            fail(f"candidate {candidate_id} has invalid origin_kind")
        provenance = candidate.get("provenance")
        if provenance not in {"VERBATIM", "ADAPTED", "DEPENDENCY_PIN", "AUTHORED"}:
            fail(f"candidate {candidate_id} has invalid provenance")
        if origin == "LOCAL_AUTHORED" and provenance != "AUTHORED":
            fail(f"local candidate {candidate_id} must remain AUTHORED")
        if origin != "LOCAL_AUTHORED" and provenance == "AUTHORED":
            fail(f"official candidate {candidate_id} cannot be relabeled AUTHORED")
        if origin == "LOCAL_AUTHORED":
            if candidate.get("source_url") is not None:
                fail(f"local candidate {candidate_id} cannot claim an official source_url")
            verify_ref(project_root, candidate.get("source_ref"), f"candidate {candidate_id}.source_ref")
        else:
            host = https_host(candidate.get("source_url"), f"candidate {candidate_id}.source_url")
            if not host_allowed(host, official_domains):
                fail(f"candidate {candidate_id} source host was not searched as official")
        require_string(candidate.get("claim"), f"candidate {candidate_id}.claim")
        require_list(candidate.get("non_claims"), f"candidate {candidate_id}.non_claims", nonempty=True)
        status = candidate.get("status")
        if status not in {"REJECTED", "CONDITIONED", "ELITE_REFERENCE", "REUSABLE_PACK"}:
            fail(f"candidate {candidate_id} has invalid status")
        revision = require_string(candidate.get("immutable_revision"), f"candidate {candidate_id}.immutable_revision")
        artifact_hash = require_string(candidate.get("artifact_sha256"), f"candidate {candidate_id}.artifact_sha256")
        artifact_ref = verify_ref(project_root, candidate.get("artifact_ref"), f"candidate {candidate_id}.artifact_ref")
        if origin == "OFFICIAL_CODE" and (revision.lower() in {"main", "master", "latest", "head"} or not re.fullmatch(r"[A-Za-z0-9._:+@/-]{7,160}", revision)):
            fail(f"official code candidate {candidate_id} lacks an immutable revision")
        if not re.fullmatch(r"[0-9a-f]{64}", artifact_hash):
            fail(f"candidate {candidate_id} lacks a lowercase SHA-256")
        artifact_path = (project_root / artifact_ref.split("#", 1)[0]).resolve()
        if hashlib.sha256(artifact_path.read_bytes()).hexdigest() != artifact_hash:
            fail(f"candidate {candidate_id} artifact SHA-256 does not match")
        require_string(candidate.get("license_expression"), f"candidate {candidate_id}.license_expression")
        verify_ref(project_root, candidate.get("license_evidence_ref"), f"candidate {candidate_id}.license_evidence_ref")
        gates = candidate.get("gates")
        if not isinstance(gates, dict) or set(gates) != set(GATES) or any(gates[g] not in GATE_VALUES for g in GATES):
            fail(f"candidate {candidate_id} must declare exactly G0..G8")
        if status == "REUSABLE_PACK" and any(gates[g] != "PASS" for g in GATES):
            fail(f"REUSABLE_PACK {candidate_id} requires G0..G8 PASS")
        if status == "ELITE_REFERENCE" and any(gates[g] != "PASS" for g in GATES[:8]):
            fail(f"ELITE_REFERENCE {candidate_id} requires G0..G7 PASS")
        if origin == "OFFICIAL_ARCHITECTURE" and status == "REUSABLE_PACK":
            fail(f"architecture-only candidate {candidate_id} cannot be a REUSABLE_PACK")
        if status == "REJECTED" and not require_list(candidate.get("rejection_reasons"), f"candidate {candidate_id}.rejection_reasons", nonempty=True):
            fail(f"rejected candidate {candidate_id} needs reasons")
        if origin == "LOCAL_AUTHORED":
            for authority in require_list(candidate.get("governing_authorities"), f"candidate {candidate_id}.governing_authorities", nonempty=True):
                verify_ref(project_root, authority, f"candidate {candidate_id}.governing_authorities[]")
            declared = set(require_list(candidate.get("no_source_search_ids"), f"candidate {candidate_id}.no_source_search_ids", nonempty=True))
            if not declared.issubset(search_ids):
                fail(f"local candidate {candidate_id} cites unknown source searches")
        by_id[candidate_id] = candidate

    decision = record.get("decision")
    if not isinstance(decision, dict):
        fail("decision must be an object")
    state = decision.get("state")
    if state == "RESEARCH_INCOMPLETE":
        fail("research remains incomplete; continue official discovery")
    if state not in TERMINAL_STATES:
        fail("decision.state is invalid")
    require_string(decision.get("reason"), "decision.reason")
    selected_id = decision.get("selected_candidate_id")
    selected = None
    if state.startswith("USE_"):
        selected_id = require_string(selected_id, "decision.selected_candidate_id")
        selected = by_id.get(selected_id)
        if selected is None:
            fail("selected_candidate_id does not exist")
        expected = "REUSABLE_PACK" if state == "USE_REUSABLE_PACK" else "CONDITIONED"
        if selected["status"] != expected:
            fail(f"{state} requires a {expected} candidate")
        for update in require_list(decision.get("canonical_updates"), "decision.canonical_updates", nonempty=True):
            verify_ref(project_root, update, "decision.canonical_updates[]")
        if state == "USE_CONDITIONED_PACK":
            require_list(decision.get("blockers"), "decision.blockers", nonempty=True)
    else:
        if selected_id is not None:
            fail(f"{state} cannot select a candidate")
        require_list(decision.get("blockers"), "decision.blockers", nonempty=True)

    if state == "NO_ADMISSIBLE_SOURCE":
        if any(c["status"] != "REJECTED" for c in candidates):
            fail("NO_ADMISSIBLE_SOURCE requires every discovered candidate to be REJECTED")
        exhaustion = decision.get("exhaustion")
        required_checks = {"searched_authority_registry", "searched_official_docs", "searched_official_repositories", "checked_release_security"}
        if not isinstance(exhaustion, dict) or any(exhaustion.get(key) is not True for key in required_checks):
            fail("NO_ADMISSIBLE_SOURCE requires all exhaustion checks")
        require_string(exhaustion.get("scope_reason"), "decision.exhaustion.scope_reason")
        require_string(exhaustion.get("next_review_trigger"), "decision.exhaustion.next_review_trigger")
        if len(official_domains) < 2 and not require_string(exhaustion.get("exclusive_authority_reason"), "decision.exhaustion.exclusive_authority_reason"):
            fail("single-domain exhaustion requires exclusive_authority_reason")

    implementation_ready = state == "USE_REUSABLE_PACK"
    return {
        "schema_version": 1,
        "gate": "CAPABILITY-GAP-RESOLUTION-GATE",
        "capability_id": capability,
        "requirement_ref": requirement_ref,
        "observed_at": observed.isoformat().replace("+00:00", "Z"),
        "authority_count": len(authorities),
        "search_count": len(searches),
        "official_domain_count": len(official_domains),
        "candidate_count": len(candidates),
        "decision_state": state,
        "selected_candidate_id": selected_id,
        "implementation_ready": implementation_ready,
        "outcome": "READY" if implementation_ready else "BLOCKED",
    }


def atomic_write(path: pathlib.Path, payload: dict) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    if path.exists():
        fail("receipt already exists")
    handle, temp_name = tempfile.mkstemp(prefix=path.name + ".tmp-", dir=path.parent)
    try:
        with os.fdopen(handle, "w", encoding="utf-8", newline="\n") as stream:
            json.dump(payload, stream, indent=2, sort_keys=True)
            stream.write("\n")
        os.replace(temp_name, path)
    except Exception:
        try:
            os.unlink(temp_name)
        except FileNotFoundError:
            pass
        raise


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--record", required=True, type=pathlib.Path)
    parser.add_argument("--receipt", required=True, type=pathlib.Path)
    parser.add_argument("--project-root", type=pathlib.Path, help="root used to resolve evidence references; defaults to record directory")
    parser.add_argument("--as-of", help="ISO-8601 time for deterministic verification")
    args = parser.parse_args()
    try:
        as_of = parse_time(args.as_of, "--as-of") if args.as_of else dt.datetime.now(dt.timezone.utc)
        raw = args.record.read_bytes()
        record = read_json(args.record)
        project_root = (args.project_root or args.record.parent).resolve()
        if not project_root.is_dir():
            fail("project root does not exist")
        receipt = validate(record, as_of, project_root)
        receipt["record_sha256"] = hashlib.sha256(raw).hexdigest()
        atomic_write(args.receipt, receipt)
        print(f"CAPABILITY_GAP_RESOLUTION_PASS capability={receipt['capability_id']} state={receipt['decision_state']} ready={str(receipt['implementation_ready']).lower()}")
        return 0
    except (GapError, OSError) as exc:
        print(f"CAPABILITY_GAP_RESOLUTION_FAIL: {exc}", file=sys.stderr)
        return 2


if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `capability_gap_resolution/verify_pack.py`
```yaml
block_id: "CAPABILITY-GAP-RESOLUTION-GATE:verifier:v1"
operation: CREATE
provenance: AUTHORED
source: "local, governed by NIST SSDF, Google SRE launch practice, NASA V&V and SLSA provenance"
license: "LicenseRef-Workspace-Owner"
sha256: "5a906f6ebae58108265ec0afa07b38ca5f255749e7d6db32bd6d4bcd3a3ebbb4"
variables: []
secrets_allowed: false
```
````python
#!/usr/bin/env python3
from __future__ import annotations

import copy
import hashlib
import json
import pathlib
import subprocess
import sys
import tempfile


AS_OF = "2026-09-05T12:00:00Z"
BASE = {
    "schema_version": 1,
    "capability_id": "SECURITY-SAST",
    "requirement_ref": "PROJECT_BLUEPRINT.md#security-sast",
    "observed_at": "2026-09-05T11:00:00Z",
    "max_age_days": 30,
    "authority_refs": ["PROJECT_AUTHORITY_MAP.md#security-sast"],
    "source_searches": [{
        "search_id": "microsoft-official",
        "query": "official Microsoft SAST source release",
        "official_domains": ["github.com", "microsoft.com"],
        "executed_at": "2026-09-05T10:00:00Z",
        "evidence_refs": ["evidence/microsoft-sast.md"]
    }],
    "candidates": [{
        "candidate_id": "official-tool-v1",
        "origin_kind": "OFFICIAL_CODE",
        "provenance": "DEPENDENCY_PIN",
        "source_url": "https://github.com/microsoft/example",
        "immutable_revision": "0123456789abcdef0123456789abcdef01234567",
        "artifact_ref": "evidence/tool.zip",
        "artifact_sha256": "REPLACED_BY_TEST",
        "license_expression": "MIT",
        "license_evidence_ref": "evidence/license.md",
        "claim": "Runs the narrowly audited SAST gate.",
        "non_claims": ["Does not prove absence of vulnerabilities."],
        "status": "REUSABLE_PACK",
        "gates": {f"G{i}": "PASS" for i in range(9)},
        "rejection_reasons": []
    }],
    "decision": {
        "state": "USE_REUSABLE_PACK",
        "selected_candidate_id": "official-tool-v1",
        "reason": "Exact source and all admission gates pass.",
        "canonical_updates": ["PROJECT_PACK_PLAN.md#security-sast"],
        "blockers": [],
        "exhaustion": None
    }
}


def run(root: pathlib.Path, name: str, record: dict, expected: int) -> None:
    record_path = root / f"{name}.json"
    receipt = root / f"{name}.receipt.json"
    record_path.write_text(json.dumps(record), encoding="utf-8")
    result = subprocess.run(
        [sys.executable, str(pathlib.Path(__file__).with_name("resolve_capability_gap.py")), "--record", str(record_path), "--receipt", str(receipt), "--project-root", str(root), "--as-of", AS_OF],
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        check=False,
    )
    if result.returncode != expected:
        raise AssertionError(f"{name}: expected {expected}, got {result.returncode}: {result.stdout} {result.stderr}")
    if expected == 0 and not receipt.exists():
        raise AssertionError(f"{name}: receipt missing")
    if expected != 0 and receipt.exists():
        raise AssertionError(f"{name}: failed case wrote a receipt")


with tempfile.TemporaryDirectory(prefix="elite-capability-gap-") as temp:
    root = pathlib.Path(temp)
    (root / "evidence").mkdir()
    for relative in ["PROJECT_BLUEPRINT.md", "PROJECT_AUTHORITY_MAP.md", "PROJECT_PACK_PLAN.md", "evidence/microsoft-sast.md", "evidence/license.md"]:
        (root / relative).write_text(f"evidence for {relative}\n", encoding="utf-8")
    artifact = root / "evidence/tool.zip"
    artifact.write_bytes(b"exact audited artifact")
    BASE["candidates"][0]["artifact_sha256"] = hashlib.sha256(artifact.read_bytes()).hexdigest()
    run(root, "reusable", copy.deepcopy(BASE), 0)

    conditioned = copy.deepcopy(BASE)
    conditioned["candidates"][0]["status"] = "CONDITIONED"
    conditioned["candidates"][0]["gates"]["G8"] = "NOT_APPLICABLE"
    conditioned["decision"].update(state="USE_CONDITIONED_PACK", blockers=["Target load evidence required."])
    run(root, "conditioned", conditioned, 0)

    no_source = copy.deepcopy(BASE)
    no_source["candidates"][0]["status"] = "REJECTED"
    no_source["candidates"][0]["gates"]["G4"] = "FAIL"
    no_source["candidates"][0]["rejection_reasons"] = ["Current dependency graph is vulnerable."]
    no_source["decision"] = {
        "state": "NO_ADMISSIBLE_SOURCE",
        "selected_candidate_id": None,
        "reason": "Every discovered official candidate failed a mandatory gate.",
        "canonical_updates": ["PROJECT_AUTHORITY_MAP.md#security-sast"],
        "blockers": ["No reusable official implementation."],
        "exhaustion": {
            "searched_authority_registry": True,
            "searched_official_docs": True,
            "searched_official_repositories": True,
            "checked_release_security": True,
            "scope_reason": "Microsoft and GitHub are the selected official source surfaces.",
            "next_review_trigger": "New signed release or security advisory resolution.",
            "exclusive_authority_reason": ""
        }
    }
    run(root, "no-source", no_source, 0)

    mobile = copy.deepcopy(BASE)
    mobile["candidates"][0]["immutable_revision"] = "main"
    run(root, "mobile-revision", mobile, 2)

    false_origin = copy.deepcopy(BASE)
    false_origin["candidates"][0]["provenance"] = "AUTHORED"
    run(root, "false-origin", false_origin, 2)

    stale = copy.deepcopy(BASE)
    stale["observed_at"] = "2025-01-01T00:00:00Z"
    run(root, "stale", stale, 2)

    incomplete = copy.deepcopy(BASE)
    incomplete["decision"]["state"] = "RESEARCH_INCOMPLETE"
    incomplete["decision"]["selected_candidate_id"] = None
    run(root, "incomplete", incomplete, 2)

    bad_exhaustion = copy.deepcopy(no_source)
    bad_exhaustion["decision"]["exhaustion"]["checked_release_security"] = False
    run(root, "bad-exhaustion", bad_exhaustion, 2)

    remote_only = copy.deepcopy(BASE)
    remote_only["candidates"][0]["artifact_ref"] = "https://github.com/microsoft/example/archive.zip"
    run(root, "remote-evidence-not-preserved", remote_only, 2)

print("CAPABILITY_GAP_RESOLUTION_PACK_PASS positives=3 negatives=6")
````

### FILE: `capability_gap_resolution/README.md`
```yaml
block_id: "CAPABILITY-GAP-RESOLUTION-GATE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local, governed by NIST SSDF, Google SRE launch practice, NASA V&V and SLSA provenance"
license: "LicenseRef-Workspace-Owner"
sha256: "2d996d23d3e11bf2ad11d2e27d0846e4d67e4acae0a652484bffdb75e4a4f3ff"
variables: []
secrets_allowed: false
```
````markdown
# Capability gap resolution gate

Use this gate only after the project classifies a capability as required and the current catalog has no compatible admitted pack. Research current primary authorities and official code surfaces, preserve exact evidence files, audit each candidate through G0–G8, then validate the record with `--project-root`. Every evidence reference must resolve to a local file inside that root and candidate artifact hashes must match its preserved bytes; keep the official remote location separately in `source_url`.

`USE_REUSABLE_PACK` is the only result that sets `implementation_ready=true`. A conditioned candidate remains blocked on its declared target evidence. If official code does not exist or all candidates fail, report `NO_ADMISSIBLE_SOURCE` with the completed search scope and next review trigger; never rename authored glue as vendor code. `RESEARCH_INCOMPLETE` intentionally fails so the agent continues discovery instead of silently improvising.

The gate validates supplied evidence; it does not browse, prove that a source is official, approve licenses, or replace builds, tests, SCA, project journeys and production gates. Those observations must be produced by the agent with current official sources and preserved in the referenced evidence.
````

## 6. Configuration surface

| Campo | Tipo | Default seguro | Validación | Secreto | Mutabilidad | Efecto |
|---|---|---|---|---|---|---|
| `capability_id` | ID | ninguno | estable, mayúsculas | no | por capability | enlaza requirement y catálogo |
| `observed_at` / `max_age_days` | tiempo | ninguno/30 | timezone, 1–365, no vencido | no | cada investigación | impide decisiones obsoletas |
| `source_searches` | lista | vacía/rechazada | dominios oficiales y evidence refs | no | append-only | demuestra investigación |
| `candidates` | lista | vacía | identidad/licencia/hash/claim/G0–G8 | no | append-only | conserva aceptación/rechazo |
| `decision` | objeto | `RESEARCH_INCOMPLETE` | estado, candidato, blockers, agotamiento | no | superseding record | controla implementación |

## 7. Dependency bill

| Tool | Pin | Uso | Licencia | Scope | Fuente |
|---|---|---|---|---|---|
| CPython | 3.11+ del proyecto, fijado por su lock | JSON estricto, validación y receipt | PSF-2.0 | tooling | <https://www.python.org/> |

Las autoridades listadas gobiernan el método; no se redistribuye código suyo.

## 8. Apply order

Materializar en tooling del proyecto. Copiar el perfil con un nombre por capability. Registrar investigación actual y ejecutar el resolver. Si falla por incompleto o stale, continuar investigación. Si queda bloqueado, informar causa/evidencia/trigger y mantener la capability sin implementar. Si queda lista, actualizar mapa/lock/pack plan canónicos antes de componer. En proyecto existente, limitar cambios al delta inventariado. Rollback: retirar tooling no referenciado; nunca borrar receipts usados por una release.

## 9. Verification

Ejecutar `python -m py_compile` y `python capability_gap_resolution/verify_pack.py`. Éxito: tres positivos —reusable, conditioned, no-source honesto— y seis negativos. En proyecto, verificar además links/bytes/licencia/build/tests/SCA/compatibilidad y todos los journeys/gates aplicables.

## 10. Reconstruction evidence

Windows x64, CPython 3.12/3.14 stdlib: cuatro archivos materializados, ambos scripts compilados y `positives=3 negatives=6` PASS. Evidencia: `reconstruction_evidence/CAPABILITY_GAP_RESOLUTION_GATE_2026-09-05_V251.md`.
