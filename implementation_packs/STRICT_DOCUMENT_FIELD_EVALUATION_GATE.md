# Strict Document Field Evaluation Gate

## 1. Metadata

```yaml
pack_id: "STRICT-DOCUMENT-FIELD-EVALUATION-GATE"
pack_version: "0.1.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa una puerta de evaluación fail-closed alrededor de AWS Labs Stickler 0.6.0: valida schema cerrado y hashes antes del motor oficial, exige comparadores exactos y cero FA/FD/FN/FP, y emite un recibo atómico que nunca autoriza almacenamiento automático."
stacks: ["CPython 3.12", "AWS Labs Stickler 0.6.0", "JSON Schema Draft-07"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x", "official document runtime packs"]
incompatible_with: ["open schemas", "moving Stickler versions", "unapproved corpus", "unreviewed ground truth", "automatic business storage"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://github.com/awslabs/stickler/tree/174ca9d3476c1ea2d0a36628c73d595604e0398d"]
verified_at: "2026-08-26"
```

## 2. Applicability

Use después de la puerta de seguridad y de un extractor seleccionado para medir payloads estructurados contra ground truth aprobado. Los cinco archivos son coordinación y pruebas locales `AUTHORED`; el algoritmo de comparación, validadores y modelos se importan de la release oficial AWS Labs Stickler 0.6.0 adquirida por commit/archivo/lock exactos.

No extrae PDF/imagen, no inventa schemas ni ground truth, no repara campos, no promete exactitud sobre documentos no incluidos y no autoriza persistencia. La release oficial tiene una brecha adversarial demostrada para extras dentro de `List[StructuredModel]`; esta puerta la contiene validando ambos JSON contra schema recursivamente cerrado antes de comparar. No se publica un parche local como código AWS.

## 3. Architecture contract

Flujo: aprobación y versión exactas → manifest/hash/path → Draft-07 válido → objetos cerrados/propiedades requeridas/referencias locales → validación de expected y actual mediante el validador oficial → modelo/comparación oficial → cero `fa/fd/fn/fp` y score 1.0 → receipt create-only sin valores originales. Cualquier divergencia rechaza cerrado.

El corpus, ground truth y schema pertenecen al proyecto y deben ser independientes del extractor evaluado. Los IDs se hashean en evidencia. El recibo siempre conserva `automatic_storage_authorized=false`. Producción requiere además clases cubiertas, drift, seguridad de ingreso, privacidad, revisión, carga, rollback y autoridad explícita de storage.

## 4. Exact file manifest

```text
CREATE strict_field_eval/evaluation-policy.template.json
CREATE strict_field_eval/evaluation-manifest.template.json
CREATE strict_field_eval/strict_field_evaluation.py
CREATE strict_field_eval/test_strict_field_evaluation.py
CREATE strict_field_eval/README.md
```

## 5. Materialization blocks

### FILE: `strict_field_eval/evaluation-policy.template.json`
```yaml
block_id: "STRICT-DOCUMENT-FIELD-EVALUATION-GATE:policy:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed policy for exact official AWS Labs Stickler engine"
license: "LicenseRef-Workspace-Owner"
sha256: "8919c0ee9ae2b3e4dfde05ceadb485675f0f0ae2911865de814dc8fe5a551d1c"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-strict-field-evaluation-policy/v1",
  "approval": {
    "status": "AWAITING_USER",
    "approval_id": "REQUIRED",
    "approved_by": "REQUIRED",
    "approved_at": "REQUIRED"
  },
  "engine": {
    "distribution": "stickler-eval",
    "version": "0.6.0",
    "source_commit": "174ca9d3476c1ea2d0a36628c73d595604e0398d",
    "source_archive_sha256": "f24ae4650b3ffc46d23e5b447a0e72a8f75a20e98df1b4ee0ca31490759707a3"
  },
  "limits": {
    "max_cases": 1000,
    "max_json_bytes": 10485760
  },
  "require_exact_comparators": true,
  "require_zero_confusion_errors": true,
  "automatic_storage_authorized": false
}
````

### FILE: `strict_field_eval/evaluation-manifest.template.json`
```yaml
block_id: "STRICT-DOCUMENT-FIELD-EVALUATION-GATE:manifest:v1"
operation: CREATE
provenance: AUTHORED
source: "local hash-linked corpus manifest; contains no upstream code"
license: "LicenseRef-Workspace-Owner"
sha256: "9b626832229405fb6d02a3248a1daba11b7f1a1e48e1cd3e64b2df4e0a49f50e"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-strict-field-evaluation-manifest/v1",
  "cases": [
    {
      "case_id": "REPLACE_WITH_OPAQUE_CASE_ID",
      "schema_path": "schemas/REPLACE.json",
      "schema_sha256": "REQUIRED",
      "expected_path": "ground_truth/REPLACE.json",
      "expected_sha256": "REQUIRED",
      "actual_path": "extractions/REPLACE.json",
      "actual_sha256": "REQUIRED"
    }
  ]
}
````

### FILE: `strict_field_eval/strict_field_evaluation.py`
```yaml
block_id: "STRICT-DOCUMENT-FIELD-EVALUATION-GATE:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local orchestration against exact public AWS Labs APIs; no upstream source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "bd006c30a972b3498ce8246568fadfaeeb1e51e13506db2f38097c5b903c7515"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
import hashlib
import importlib.metadata as metadata
import json
import os
import re
import tempfile
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Iterable

from jsonschema.exceptions import SchemaError, ValidationError
from stickler.structured_object_evaluator import StructuredModel
from stickler.structured_object_evaluator.utils.json_schema_validator import (
    validate_instance_against_schema,
    validate_json_schema,
)


POLICY_SCHEMA = "elite-strict-field-evaluation-policy/v1"
MANIFEST_SCHEMA = "elite-strict-field-evaluation-manifest/v1"
ENGINE_DISTRIBUTION = "stickler-eval"
ENGINE_VERSION = "0.6.0"
ENGINE_COMMIT = "174ca9d3476c1ea2d0a36628c73d595604e0398d"
ENGINE_ARCHIVE_SHA256 = "f24ae4650b3ffc46d23e5b447a0e72a8f75a20e98df1b4ee0ca31490759707a3"
HEX256 = re.compile(r"^[0-9a-f]{64}$")
CASE_ID = re.compile(r"^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$")
UNSUPPORTED_SCHEMA_KEYWORDS = {
    "allOf", "anyOf", "oneOf", "not", "if", "then", "else",
    "patternProperties", "dependencies", "dependentSchemas",
    "unevaluatedProperties", "additionalItems",
}


class EvaluationConfigurationError(RuntimeError):
    pass


class EvaluationRejected(RuntimeError):
    pass


def sha256_bytes(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def require_regular_file(path: Path, label: str) -> None:
    if path.is_symlink() or not path.is_file():
        raise EvaluationConfigurationError(f"{label} must be a regular non-symlink file")


def load_json(path: Path, label: str, max_bytes: int) -> tuple[Any, str]:
    require_regular_file(path, label)
    size = path.stat().st_size
    if size <= 0 or size > max_bytes:
        raise EvaluationConfigurationError(f"{label} size is outside the approved limit")
    raw = path.read_bytes()
    try:
        value = json.loads(raw.decode("utf-8"))
    except (UnicodeDecodeError, json.JSONDecodeError) as error:
        raise EvaluationConfigurationError(f"{label} is not canonical UTF-8 JSON") from error
    return value, sha256_bytes(raw)


def require_exact_sha(value: Any, label: str) -> str:
    if not isinstance(value, str) or not HEX256.fullmatch(value):
        raise EvaluationConfigurationError(f"{label} must be an exact lowercase SHA-256")
    return value


def load_policy(path: Path) -> dict[str, Any]:
    policy, _ = load_json(path, "policy", 1024 * 1024)
    if not isinstance(policy, dict) or policy.get("schema") != POLICY_SCHEMA:
        raise EvaluationConfigurationError("unsupported policy schema")
    approval = policy.get("approval") or {}
    if approval.get("status") != "APPROVED":
        raise EvaluationConfigurationError("policy approval is not APPROVED")
    for key in ("approval_id", "approved_by", "approved_at"):
        value = approval.get(key)
        if not isinstance(value, str) or not value.strip() or value == "REQUIRED":
            raise EvaluationConfigurationError(f"approval.{key} is unresolved")
    try:
        approved_at = datetime.fromisoformat(approval["approved_at"].replace("Z", "+00:00"))
    except ValueError as error:
        raise EvaluationConfigurationError("approval.approved_at is not ISO-8601") from error
    if approved_at.tzinfo is None:
        raise EvaluationConfigurationError("approval.approved_at must include timezone")
    engine = policy.get("engine") or {}
    expected_engine = {
        "distribution": ENGINE_DISTRIBUTION,
        "version": ENGINE_VERSION,
        "source_commit": ENGINE_COMMIT,
        "source_archive_sha256": ENGINE_ARCHIVE_SHA256,
    }
    if engine != expected_engine:
        raise EvaluationConfigurationError("engine identity does not match the audited AWS Labs release")
    if metadata.version(ENGINE_DISTRIBUTION) != ENGINE_VERSION:
        raise EvaluationConfigurationError("installed stickler-eval version mismatch")
    limits = policy.get("limits") or {}
    for key in ("max_cases", "max_json_bytes"):
        if not isinstance(limits.get(key), int) or isinstance(limits.get(key), bool) or limits[key] <= 0:
            raise EvaluationConfigurationError(f"limits.{key} must be a positive integer")
    if policy.get("require_exact_comparators") is not True:
        raise EvaluationConfigurationError("exact comparators must remain required")
    if policy.get("require_zero_confusion_errors") is not True:
        raise EvaluationConfigurationError("zero confusion errors must remain required")
    if policy.get("automatic_storage_authorized") is not False:
        raise EvaluationConfigurationError("this evaluation gate cannot authorize automatic storage")
    return policy


def resolve_local_ref(root: dict[str, Any], ref: str) -> Any:
    if not ref.startswith("#/"):
        raise EvaluationConfigurationError("only local JSON Pointer references are allowed")
    node: Any = root
    for raw_part in ref[2:].split("/"):
        part = raw_part.replace("~1", "/").replace("~0", "~")
        if not isinstance(node, dict) or part not in node:
            raise EvaluationConfigurationError(f"unresolved local schema reference: {ref}")
        node = node[part]
    return node


def iter_schema_nodes(node: Any) -> Iterable[dict[str, Any]]:
    if isinstance(node, bool):
        raise EvaluationConfigurationError("boolean JSON Schemas are not admitted")
    if not isinstance(node, dict):
        return
    yield node
    for key, value in node.items():
        if key in ("properties", "definitions", "$defs") and isinstance(value, dict):
            for child in value.values():
                yield from iter_schema_nodes(child)
        elif key == "items":
            yield from iter_schema_nodes(value)


def validate_closed_exact_schema(schema: Any) -> dict[str, Any]:
    if not isinstance(schema, dict):
        raise EvaluationConfigurationError("schema document must be an object")
    try:
        validate_json_schema(schema)
    except (SchemaError, ValueError) as error:
        raise EvaluationConfigurationError("schema is not valid Draft-07") from error
    if schema.get("type") != "object":
        raise EvaluationConfigurationError("root schema type must be object")
    for node in iter_schema_nodes(schema):
        forbidden = UNSUPPORTED_SCHEMA_KEYWORDS.intersection(node)
        if forbidden:
            raise EvaluationConfigurationError(f"unsupported schema keyword: {sorted(forbidden)[0]}")
        ref = node.get("$ref")
        if ref is not None:
            if not isinstance(ref, str):
                raise EvaluationConfigurationError("$ref must be a string")
            target = resolve_local_ref(schema, ref)
            if not isinstance(target, dict):
                raise EvaluationConfigurationError("local $ref target must be a schema object")
        node_type = node.get("type")
        if node_type == "object":
            properties = node.get("properties")
            if not isinstance(properties, dict) or not properties:
                raise EvaluationConfigurationError("every object schema must declare non-empty properties")
            if node.get("additionalProperties") is not False:
                raise EvaluationConfigurationError("every object schema must set additionalProperties=false")
            required = node.get("required")
            if not isinstance(required, list) or set(required) != set(properties):
                raise EvaluationConfigurationError("every object schema must require every declared property")
        elif node_type == "array":
            if not isinstance(node.get("items"), dict):
                raise EvaluationConfigurationError("every array schema must declare one object-form items schema")
        elif node_type in ("string", "integer", "number", "boolean"):
            if node.get("x-aws-stickler-comparator") != "ExactComparator":
                raise EvaluationConfigurationError("every scalar must use AWS Stickler ExactComparator")
            if node.get("x-aws-stickler-threshold") != 1.0:
                raise EvaluationConfigurationError("every scalar must use Stickler threshold 1.0")
        elif node_type is not None:
            raise EvaluationConfigurationError(f"unsupported schema type: {node_type}")
    for node in iter_schema_nodes(schema):
        ref = node.get("$ref")
        if isinstance(ref, str):
            list(iter_schema_nodes(resolve_local_ref(schema, ref)))
    return schema


def safe_relative(root: Path, raw: Any, label: str) -> Path:
    if not isinstance(raw, str) or not raw or "\\" in raw:
        raise EvaluationConfigurationError(f"{label} must be a canonical relative POSIX path")
    relative = Path(raw)
    if relative.is_absolute() or any(part in ("", ".", "..") for part in relative.parts):
        raise EvaluationConfigurationError(f"{label} is unsafe")
    resolved = (root / relative).resolve()
    try:
        resolved.relative_to(root.resolve())
    except ValueError as error:
        raise EvaluationConfigurationError(f"{label} escapes data root") from error
    return resolved


def verify_hash(actual: str, expected: Any, label: str) -> None:
    if actual != require_exact_sha(expected, label):
        raise EvaluationConfigurationError(f"{label} mismatch")


def evaluate_case(case: dict[str, Any], data_root: Path, max_bytes: int) -> dict[str, Any]:
    case_id = case.get("case_id")
    if not isinstance(case_id, str) or not CASE_ID.fullmatch(case_id):
        raise EvaluationConfigurationError("case_id must be opaque and canonical")
    paths: dict[str, Path] = {}
    values: dict[str, Any] = {}
    hashes: dict[str, str] = {}
    for name in ("schema", "expected", "actual"):
        paths[name] = safe_relative(data_root, case.get(f"{name}_path"), f"{name}_path")
        values[name], hashes[name] = load_json(paths[name], name, max_bytes)
        verify_hash(hashes[name], case.get(f"{name}_sha256"), f"{name}_sha256")
    schema = validate_closed_exact_schema(values["schema"])
    if not isinstance(values["expected"], dict) or not isinstance(values["actual"], dict):
        raise EvaluationConfigurationError("expected and actual payloads must be JSON objects")
    for label in ("expected", "actual"):
        try:
            validate_instance_against_schema(values[label], schema)
        except ValidationError as error:
            raise EvaluationRejected(f"{label} payload does not conform to approved schema") from error
    model = StructuredModel.from_json_schema(schema)
    expected_model = model.from_json(values["expected"])
    actual_model = model.from_json(values["actual"])
    result = expected_model.compare_with(
        actual_model,
        include_confusion_matrix=True,
        document_non_matches=True,
        document_field_comparisons=True,
    )
    overall = (result.get("confusion_matrix") or {}).get("overall") or {}
    metrics: dict[str, int] = {}
    for key in ("fa", "fd", "fn", "fp"):
        value = overall.get(key)
        if not isinstance(value, int) or isinstance(value, bool) or value < 0:
            raise EvaluationRejected(f"Stickler did not return a valid {key} metric")
        metrics[key] = value
    score = result.get("overall_score")
    passed = all(value == 0 for value in metrics.values()) and score == 1.0
    return {
        "case_id_sha256": sha256_bytes(case_id.encode("utf-8")),
        "input_sha256": hashes,
        "confusion_errors": metrics,
        "overall_score": score,
        "passed": passed,
    }


def commit_receipt(output: Path, receipt: dict[str, Any]) -> None:
    parent = output.parent
    if parent.is_symlink() or not parent.is_dir():
        raise EvaluationConfigurationError("output parent must be a non-symlink directory")
    if output.exists() or output.is_symlink():
        raise FileExistsError(f"output already exists: {output}")
    stage = Path(tempfile.mkdtemp(prefix=f".{output.name}.staging-", dir=parent))
    try:
        candidate = stage / "evaluation-receipt.json"
        candidate.write_text(json.dumps(receipt, indent=2, sort_keys=True) + "\n", encoding="utf-8", newline="\n")
        os.replace(stage, output)
    except BaseException:
        if stage.exists():
            for child in stage.iterdir():
                child.unlink()
            stage.rmdir()
        raise


def run(policy_path: Path, manifest_path: Path, data_root: Path, output: Path) -> dict[str, Any]:
    policy = load_policy(policy_path)
    if data_root.is_symlink() or not data_root.is_dir():
        raise EvaluationConfigurationError("data root must be a non-symlink directory")
    manifest, manifest_hash = load_json(manifest_path, "manifest", 1024 * 1024)
    if not isinstance(manifest, dict) or manifest.get("schema") != MANIFEST_SCHEMA:
        raise EvaluationConfigurationError("unsupported manifest schema")
    cases = manifest.get("cases")
    if not isinstance(cases, list) or not cases or len(cases) > policy["limits"]["max_cases"]:
        raise EvaluationConfigurationError("manifest cases count is outside the approved limit")
    if not all(isinstance(item, dict) for item in cases):
        raise EvaluationConfigurationError("every case must be an object")
    ids = [item.get("case_id") for item in cases]
    if len(ids) != len(set(ids)):
        raise EvaluationConfigurationError("case_id values must be unique")
    results = [evaluate_case(item, data_root, policy["limits"]["max_json_bytes"]) for item in cases]
    totals = {key: sum(item["confusion_errors"][key] for item in results) for key in ("fa", "fd", "fn", "fp")}
    passed = all(item["passed"] for item in results) and all(value == 0 for value in totals.values())
    receipt = {
        "schema": "elite-strict-field-evaluation-receipt/v1",
        "evaluated_at": datetime.now(timezone.utc).isoformat().replace("+00:00", "Z"),
        "engine": {
            "distribution": ENGINE_DISTRIBUTION,
            "version": ENGINE_VERSION,
            "source_commit": ENGINE_COMMIT,
            "source_archive_sha256": ENGINE_ARCHIVE_SHA256,
        },
        "approval_id_sha256": sha256_bytes(policy["approval"]["approval_id"].encode("utf-8")),
        "manifest_sha256": manifest_hash,
        "case_count": len(results),
        "confusion_error_totals": totals,
        "cases": results,
        "evaluation_passed": passed,
        "automatic_storage_authorized": False,
    }
    commit_receipt(output, receipt)
    if not passed:
        raise EvaluationRejected("field evaluation rejected: one or more exact comparisons failed")
    return receipt


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--policy", type=Path, required=True)
    parser.add_argument("--manifest", type=Path, required=True)
    parser.add_argument("--data-root", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    try:
        receipt = run(args.policy, args.manifest, args.data_root, args.output)
    except (EvaluationConfigurationError, EvaluationRejected, FileExistsError, ValueError) as error:
        print(f"STRICT_FIELD_EVALUATION_REJECTED: {error}")
        return 2
    print(f"STRICT_FIELD_EVALUATION_PASS cases={receipt['case_count']}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `strict_field_eval/test_strict_field_evaluation.py`
```yaml
block_id: "STRICT-DOCUMENT-FIELD-EVALUATION-GATE:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local regression and adversarial tests against the exact official engine"
license: "LicenseRef-Workspace-Owner"
sha256: "2a361107be09d38f90dae01e6a79f999ce96aacd3f088fa4ed4c21f1ca12fbfb"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import hashlib
import json
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

import strict_field_evaluation as gate


def digest(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def write_json(path: Path, value: object) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(value, sort_keys=True) + "\n", encoding="utf-8", newline="\n")


def exact_scalar(kind: str) -> dict[str, object]:
    return {
        "type": kind,
        "x-aws-stickler-comparator": "ExactComparator",
        "x-aws-stickler-threshold": 1.0,
    }


SCHEMA = {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "x-aws-stickler-model-name": "InvoiceEvaluation",
    "x-aws-stickler-match-threshold": 1.0,
    "additionalProperties": False,
    "required": ["invoice_number", "total", "lines"],
    "properties": {
        "invoice_number": exact_scalar("string"),
        "total": exact_scalar("string"),
        "lines": {"type": "array", "items": {"$ref": "#/definitions/Line"}},
    },
    "definitions": {
        "Line": {
            "type": "object",
            "x-aws-stickler-model-name": "InvoiceLineEvaluation",
            "x-aws-stickler-match-threshold": 1.0,
            "additionalProperties": False,
            "required": ["sku", "quantity"],
            "properties": {
                "sku": exact_scalar("string"),
                "quantity": exact_scalar("integer"),
            },
        }
    },
}

BASE = {
    "invoice_number": "INV-001",
    "total": "1250.00",
    "lines": [{"sku": "SKU-1", "quantity": 2}],
}


class GateTest(unittest.TestCase):
    def setUp(self) -> None:
        self.temp = tempfile.TemporaryDirectory()
        self.root = Path(self.temp.name)
        self.policy = self.root / "policy.json"
        policy = {
            "schema": gate.POLICY_SCHEMA,
            "approval": {
                "status": "APPROVED",
                "approval_id": "APR-001",
                "approved_by": "owner",
                "approved_at": "2026-08-26T12:00:00Z",
            },
            "engine": {
                "distribution": gate.ENGINE_DISTRIBUTION,
                "version": gate.ENGINE_VERSION,
                "source_commit": gate.ENGINE_COMMIT,
                "source_archive_sha256": gate.ENGINE_ARCHIVE_SHA256,
            },
            "limits": {"max_cases": 10, "max_json_bytes": 100000},
            "require_exact_comparators": True,
            "require_zero_confusion_errors": True,
            "automatic_storage_authorized": False,
        }
        write_json(self.policy, policy)

    def tearDown(self) -> None:
        self.temp.cleanup()

    def prepare(self, expected: object = BASE, actual: object = BASE, schema: object = SCHEMA) -> Path:
        paths = {
            "schema": self.root / "schemas" / "invoice.json",
            "expected": self.root / "ground_truth" / "case.json",
            "actual": self.root / "extractions" / "case.json",
        }
        write_json(paths["schema"], schema)
        write_json(paths["expected"], expected)
        write_json(paths["actual"], actual)
        manifest = {
            "schema": gate.MANIFEST_SCHEMA,
            "cases": [{
                "case_id": "CASE-001",
                **{f"{name}_path": path.relative_to(self.root).as_posix() for name, path in paths.items()},
                **{f"{name}_sha256": digest(path) for name, path in paths.items()},
            }],
        }
        manifest_path = self.root / "manifest.json"
        write_json(manifest_path, manifest)
        return manifest_path

    def test_exact_case_passes_and_receipt_contains_no_raw_values(self) -> None:
        manifest = self.prepare()
        output = self.root / "receipt"
        receipt = gate.run(self.policy, manifest, self.root, output)
        self.assertTrue(receipt["evaluation_passed"])
        self.assertFalse(receipt["automatic_storage_authorized"])
        serialized = (output / "evaluation-receipt.json").read_text(encoding="utf-8")
        self.assertNotIn("INV-001", serialized)
        self.assertNotIn("SKU-1", serialized)
        self.assertNotIn("APR-001", serialized)

    def test_root_extra_is_rejected_by_closed_schema(self) -> None:
        actual = {**BASE, "invented_supplier": "hallucinated"}
        with self.assertRaises(gate.EvaluationRejected):
            gate.run(self.policy, self.prepare(actual=actual), self.root, self.root / "out")

    def test_nested_list_extra_is_rejected_before_known_stickler_gap(self) -> None:
        actual = {**BASE, "lines": [{"sku": "SKU-1", "quantity": 2, "invented_discount": "50%"}]}
        with self.assertRaises(gate.EvaluationRejected):
            gate.run(self.policy, self.prepare(actual=actual), self.root, self.root / "out")

    def test_missing_required_field_is_rejected(self) -> None:
        actual = {"invoice_number": "INV-001", "lines": BASE["lines"]}
        with self.assertRaises(gate.EvaluationRejected):
            gate.run(self.policy, self.prepare(actual=actual), self.root, self.root / "out")

    def test_schema_rejection_does_not_echo_raw_value(self) -> None:
        secret_value = "CUSTOMER-TAX-ID-SECRET"
        actual = {**BASE, "invented_supplier": secret_value}
        with self.assertRaises(gate.EvaluationRejected) as raised:
            gate.run(self.policy, self.prepare(actual=actual), self.root, self.root / "out")
        self.assertNotIn(secret_value, str(raised.exception))

    def test_duplicate_ground_truth_row_causes_false_negative(self) -> None:
        expected = {**BASE, "lines": [BASE["lines"][0], BASE["lines"][0]]}
        manifest = self.prepare(expected=expected)
        with self.assertRaises(gate.EvaluationRejected):
            gate.run(self.policy, manifest, self.root, self.root / "out")
        receipt = json.loads((self.root / "out" / "evaluation-receipt.json").read_text(encoding="utf-8"))
        self.assertGreaterEqual(receipt["confusion_error_totals"]["fn"], 1)

    def test_wrong_exact_value_causes_false_discovery(self) -> None:
        actual = {**BASE, "total": "9999.99"}
        with self.assertRaises(gate.EvaluationRejected):
            gate.run(self.policy, self.prepare(actual=actual), self.root, self.root / "out")
        receipt = json.loads((self.root / "out" / "evaluation-receipt.json").read_text(encoding="utf-8"))
        self.assertGreaterEqual(receipt["confusion_error_totals"]["fd"], 1)

    def test_open_object_schema_is_rejected(self) -> None:
        schema = json.loads(json.dumps(SCHEMA))
        del schema["definitions"]["Line"]["additionalProperties"]
        with self.assertRaises(gate.EvaluationConfigurationError):
            gate.run(self.policy, self.prepare(schema=schema), self.root, self.root / "out")

    def test_external_reference_is_rejected(self) -> None:
        schema = json.loads(json.dumps(SCHEMA))
        schema["properties"]["lines"]["items"] = {"$ref": "https://example.invalid/line.json"}
        with self.assertRaises(gate.EvaluationConfigurationError):
            gate.run(self.policy, self.prepare(schema=schema), self.root, self.root / "out")

    def test_tampered_input_hash_is_rejected(self) -> None:
        manifest = self.prepare()
        value = json.loads(manifest.read_text(encoding="utf-8"))
        value["cases"][0]["actual_sha256"] = "0" * 64
        write_json(manifest, value)
        with self.assertRaises(gate.EvaluationConfigurationError):
            gate.run(self.policy, manifest, self.root, self.root / "out")

    def test_installed_engine_version_mismatch_is_rejected(self) -> None:
        with patch.object(gate.metadata, "version", return_value="0.7.0"):
            with self.assertRaises(gate.EvaluationConfigurationError):
                gate.load_policy(self.policy)

    def test_output_is_create_only(self) -> None:
        manifest = self.prepare()
        output = self.root / "out"
        gate.run(self.policy, manifest, self.root, output)
        with self.assertRaises(FileExistsError):
            gate.run(self.policy, manifest, self.root, output)


if __name__ == "__main__":
    unittest.main(verbosity=2)
````

### FILE: `strict_field_eval/README.md`
```yaml
block_id: "STRICT-DOCUMENT-FIELD-EVALUATION-GATE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local operational boundary and provenance statement"
license: "LicenseRef-Workspace-Owner"
sha256: "4f597141f14f413e859f48dfb0a3bc2ddc78d6bda98c102f411ee7ccfc5df70a"
variables: []
secrets_allowed: false
```
````markdown
# Strict document field evaluation gate

This directory is local `AUTHORED` orchestration around the official AWS Labs
Stickler 0.6.0 engine. It does not copy or modify AWS source and must never be
described as AWS-authored application code.

Use the exact source profile and approval generated by
`elite_sources/apply_source_profile.ps1`, acquire the pinned source, then run
`uv sync --frozen` in that exact checkout. The admitted release is commit
`174ca9d3476c1ea2d0a36628c73d595604e0398d`; moving tags or `main` are rejected.

Before a run, the project owner must replace the policy and manifest templates,
approve opaque cases, provide independently reviewed ground truth, and provide a
closed Draft-07 JSON Schema. Every object must set `additionalProperties: false`
and require every property. Every scalar must select Stickler's
`ExactComparator` with threshold `1.0`.

Run from the frozen Stickler environment:

```text
python strict_field_eval/strict_field_evaluation.py --policy project/evaluation-policy.json --manifest project/evaluation-manifest.json --data-root project/evaluation-data --output project/evaluation-receipt
python -m unittest -v strict_field_eval/test_strict_field_evaluation.py
```

The gate rejects schema/payload/hash/version divergence and requires zero
`fa`, `fd`, `fn`, and `fp` plus score 1.0 for every case. This independently
closes the demonstrated Stickler 0.6.0 nested-list-extra gap by validating both
payloads against a recursively closed schema before comparison. Receipts contain
only hashes, counts, metrics, engine identity, and a hashed approval ID.

Even a PASS has `automatic_storage_authorized: false`. It proves only the
approved evaluation corpus. Production storage still requires document-class
coverage, secure ingestion, provider/runtime evidence, privacy/retention,
human-review policy, drift/load/rollback tests, and a separate project decision.
````

## 6. Configuration surface

| Variable | Tipo | Default | Validación | Secreto | Mutabilidad | Efecto |
|---|---|---|---|---|---|---|
| policy | JSON path | none | aprobación completa, engine exacto, límites positivos, storage false | no | por expediente | autoridad y límites |
| manifest | JSON path | none | casos únicos, paths relativos y SHA-256 exactos | puede revelar paths | por corrida | corpus evaluado |
| data-root | directory | none | directorio real no symlink; paths confinados | contiene datos | por corrida | schemas/ground truth/extracciones |
| output | absent path | none | parent real, target inexistente | no | create-only | receipt atómico |

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| AWS Labs Stickler | `stickler-eval 0.6.0`, commit `174ca9d3476c1ea2d0a36628c73d595604e0398d`, archive SHA `f24ae4650b3ffc46d23e5b447a0e72a8f75a20e98df1b4ee0ca31490759707a3`, `uv.lock` SHA `325f27921c235b3a3620d3523e615fa3f63525a09deb0ed2fad96da901884dd6` | schema/model/comparación/Hungarian/confusion matrix | Apache-2.0 + NOTICE | runtime/test | https://github.com/awslabs/stickler/releases/tag/v0.6.0 |
| Python | 3.12.13 en evidencia | runner/tests | PSF | runtime | https://www.python.org/ |
| Orquestación local | archivos completos de este pack | fail-closed policy/hash/receipt | LicenseRef-Workspace-Owner | runtime/test | local, no atribuida a AWS |

El grafo oficial completo se instala sólo con `uv sync --frozen` desde el source adquirido; no se reemplaza por `pip install latest` ni se copian rangos del `pyproject.toml`.

## 8. Apply order

1. Componer `OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.19` y este pack en target vacío.
2. Completar/aprobar `strict-field-evaluation.json`; adquirir sólo `awslabs-stickler-eval-0.6.0` mediante el runner materializado.
3. En el checkout exacto ejecutar `uv sync --frozen`; no editar upstream.
4. Copiar los templates al expediente del proyecto, completar clases/schema/corpus/ground truth y hashes.
5. Ejecutar tests y luego runner. Un rechazo conserva receipt si la comparación alcanzó métricas; errores de configuración no crean receipt.
6. Rollback: retirar únicamente el directorio nuevo si no contiene evidencia posterior; nunca modificar un receipt existente.

## 9. Verification

```text
python test_strict_field_evaluation.py
expected: 12 tests PASS against stickler-eval 0.6.0 frozen environment

python strict_field_evaluation.py --policy APPROVED.json --manifest MANIFEST.json --data-root DATA --output RECEIPT
expected positive: STRICT_FIELD_EVALUATION_PASS cases=N; receipt has zero FA/FD/FN/FP and storage=false
expected negative: exit 2 or exception for open/nested-extra/missing/duplicate/wrong/hash/version/path/output divergence
```

Repetir además los 1.543 tests oficiales en Linux: el audit Windows conserva 1.535 PASS/6 `SIGALRM` FAIL/2 skip. Promotion exige corpus del proyecto, security/privacy, drift, carga, licencias/SBOM y la corrección o contención demostrada de cada fallo upstream.

## 10. Reconstruction evidence

Entorno limpio corto Windows, Python 3.12.13 y `uv sync --frozen` sobre el commit exacto: 12/12 tests locales PASS usando el motor oficial; schema abierto, extra raíz, extra anidado, campo requerido ausente, fila duplicada, valor incorrecto, ref externa, hash alterado, versión alterada, filtración de valor crudo y overwrite fueron rechazados. El receipt positivo no contiene valores de documento/aprobación y mantiene storage false.

Evidencia: `reconstruction_evidence/OFFICIAL_FIELD_EVIDENCE_CODE_REAUDIT_2026-08-26_V1.md` y `reconstruction_evidence/STRICT_DOCUMENT_FIELD_EVALUATION_GATE_2026-08-26_V1.md`. Estado condicionado: los resultados no autorizan extrapolar a clases no representadas ni usar Stickler solo como storage gate.
