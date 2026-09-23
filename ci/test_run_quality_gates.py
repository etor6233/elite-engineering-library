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
