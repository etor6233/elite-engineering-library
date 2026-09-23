from __future__ import annotations

import copy
from datetime import datetime, timedelta, timezone
import hashlib
import json
from pathlib import Path
import tempfile
import unittest

from validate_k6_load_admission import canonical_sha, validate


class K6LoadAdmissionTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temp = tempfile.TemporaryDirectory(); self.root = Path(self.temp.name).resolve(); (self.root / "evidence").mkdir()
        now = datetime.now(timezone.utc); stamp = lambda value: value.isoformat().replace("+00:00", "Z")
        self.profile = {"schema":"elite-k6-load-admission/v1", "project_id":"revestex", "environment":"production", "release_digest":"sha256:" + "a" * 64, "evaluated_at":stamp(now), "expires_at":stamp(now + timedelta(days=7)), "executor_ref":"ci:run/42", "tool":{"name":"Grafana k6", "version":"2.2.0", "source":"https://github.com/grafana/k6/releases/tag/v2.2.0", "sha256":"87dfa91bc3e47bc4bd77911d59d7ff79f25cd76fb8322c97072f08f08a0da5ed"}, "target":{"service":"https://production.invalid", "region":"primary"}, "error_budget":"evidence/budget.json", "runs":[], "output":"evidence/load-control.json"}
        budget_runs = {}
        for kind in ("BASELINE", "SOAK", "SATURATION_RECOVERY"):
            budget_runs[kind] = {"minimum_iterations":2, "minimum_duration_seconds":60, "maximum_http_req_failed_rate":0.01, "maximum_http_req_duration_p95_ms":500}
        self.write("evidence/budget.json", {"schema":"elite-load-error-budget/v1", "project_id":"revestex", "environment":"production", "release_digest":self.profile["release_digest"], "approved_at":stamp(now), "approvers":["performance-owner", "service-owner"], "runs":budget_runs})
        for kind in ("BASELINE", "SOAK", "SATURATION_RECOVERY"):
            stem = kind.lower(); script = self.root / "evidence" / f"{stem}.js"; script.write_text("export default function() {}\n", encoding="utf-8")
            args = ["run", "--summary-export", f"evidence/{stem}-summary.json", f"evidence/{stem}.js"]
            target = dict(self.profile["target"]); target.update({"run_kind":kind, "workload_script_sha256":hashlib.sha256(script.read_bytes()).hexdigest(), "minimum_duration_seconds":60.0})
            stdout = self.root / "evidence" / f"{stem}.stdout.bin"; stderr = self.root / "evidence" / f"{stem}.stderr.bin"; stdout.write_bytes(b"pass\n"); stderr.write_bytes(b"")
            pf = lambda path: {"path":path.relative_to(self.root).as_posix(), "bytes":path.stat().st_size, "sha256":hashlib.sha256(path.read_bytes()).hexdigest()}
            receipt = {"schema":"elite-official-tool-execution/v1", "project_id":"revestex", "environment":"production", "release_digest":self.profile["release_digest"], "control_id":"LOAD_RESILIENCE", "executed_at":stamp(now), "exit_code":0, "tool":dict(self.profile["tool"]), "arguments_sha256":canonical_sha(args), "working_directory":".", "target":target, "environment_variable_names":[], "stdout":pf(stdout), "stderr":pf(stderr)}
            self.write(f"evidence/{stem}-receipt.json", receipt)
            metrics = [{"name":"iterations", "values":{"count":2}}, {"name":"http_req_failed", "values":{"rate":0.0}}, {"name":"http_req_duration", "values":{"p(95)":120.0}}]
            self.write(f"evidence/{stem}-summary.json", {"version":"1.0.0", "metadata":{"k6_version":"2.2.0"}, "results":{"metrics":metrics}})
            self.profile["runs"].append({"kind":kind, "execution_receipt":f"evidence/{stem}-receipt.json", "summary":f"evidence/{stem}-summary.json", "script":f"evidence/{stem}.js", "expected_arguments":args})

    def tearDown(self) -> None: self.temp.cleanup()
    def write(self, relative: str, value: object) -> None: (self.root / relative).write_text(json.dumps(value), encoding="utf-8")

    def test_three_distinct_runs_emit_semantic_load_receipt(self) -> None:
        result = validate(self.profile, self.root); self.assertEqual("PASS", result["result"]); self.assertEqual(16, len(result["evidence"])); self.assertTrue(all(result["assertions"].values()))

    def test_nonzero_k6_exit_is_rejected(self) -> None:
        value = copy.deepcopy(self.profile); path = self.root / value["runs"][0]["execution_receipt"]; receipt = json.loads(path.read_text()); receipt["exit_code"] = 99; self.write(value["runs"][0]["execution_receipt"], receipt)
        with self.assertRaisesRegex(ValueError, "identity/exit"): validate(value, self.root)

    def test_script_tamper_is_rejected(self) -> None:
        (self.root / self.profile["runs"][1]["script"]).write_text("tampered\n", encoding="utf-8")
        with self.assertRaisesRegex(ValueError, "binding mismatch"): validate(self.profile, self.root)

    def test_argument_drift_is_rejected(self) -> None:
        value = copy.deepcopy(self.profile); value["runs"][0]["expected_arguments"].append("--insecure-skip-tls-verify")
        with self.assertRaisesRegex(ValueError, "hash-bound"): validate(value, self.root)

    def test_budget_breach_is_rejected(self) -> None:
        path = self.root / self.profile["runs"][2]["summary"]; summary = json.loads(path.read_text()); summary["results"]["metrics"][1]["values"]["rate"] = 0.02; self.write(self.profile["runs"][2]["summary"], summary)
        with self.assertRaisesRegex(ValueError, "failure rate"): validate(self.profile, self.root)

    def test_missing_soak_is_rejected(self) -> None:
        value = copy.deepcopy(self.profile); value["runs"].pop(1)
        with self.assertRaisesRegex(ValueError, "three exact"): validate(value, self.root)

    def test_single_approver_is_rejected(self) -> None:
        path = self.root / self.profile["error_budget"]; budget = json.loads(path.read_text()); budget["approvers"] = ["one"]; self.write(self.profile["error_budget"], budget)
        with self.assertRaisesRegex(ValueError, "two distinct"): validate(self.profile, self.root)

    def test_legacy_summary_is_rejected(self) -> None:
        path = self.root / self.profile["runs"][0]["summary"]; summary = json.loads(path.read_text()); summary["version"] = "legacy"; self.write(self.profile["runs"][0]["summary"], summary)
        with self.assertRaisesRegex(ValueError, "machine-readable"): validate(self.profile, self.root)


if __name__ == "__main__": unittest.main()
