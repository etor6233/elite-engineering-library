"""AUTHORED fake-provider end-to-end fixtures. No network or provider binaries run."""
import base64
import copy
import hashlib
import json
from pathlib import Path
import sqlite3
import subprocess
import sys
import tempfile
import unittest
from unittest.mock import patch

import finops_cycle as f


def fixture_config():
    return {"schema": "elite-finops-cycle/v403.1", "production_authorized": False,
            "environment": "development", "project": "elite-fixture", "location": "US",
            "billing_account": "AAAAAA-BBBBBB-CCCCCC", "billing_table": "elite-fixture.billing.export_table",
            "currency": "USD", "period_policy": "UTC_CURRENT_MONTH", "maximum_bytes_billed": 10000000,
            "interval_seconds": 300, "max_state_bytes": 1048576, "max_age_seconds": 600,
            "limit_micros": 1000000, "journal_uri": "gs://elite-fixture-state/control/billing.sqlite",
            "fence_uri": "gs://elite-fixture-state/control/finops-intents",
            "notification": {"mode": "PUBSUB", "topic": "projects/elite-fixture/topics/billing-alerts", "dedupe_seconds": 3600}}


class FakeProvider:
    def __init__(self, seed, uri):
        self.uri = uri
        self.generation = 1
        self.objects = {uri: seed.read_bytes()}
        self.versions = {(uri, 1): seed.read_bytes()}
        self.rows = [{"project": "fixture-project", "service": "fixture-service", "currency": "USD", "amount_micros": "1500000"},
                     {"project": "_ACCOUNT_LEVEL", "service": "credit", "currency": "USD", "amount_micros": "-100000"}]
        self.calls = []
        self.queries = self.publishes = self.writes = 0
        self.conflict_at_write = None
        self.publish_timeout = False
        self.persisted_before_publish = False

    def __call__(self, argv, cwd, environment, timeout):
        args = argv[1:]
        self.calls.append(args)
        self.assert_timeout = timeout
        value = {}
        code = 0
        if "query" in args:
            self.queries += 1
            value = self.rows
        elif args[:3] == ["storage", "objects", "describe"]:
            content = self.objects[self.uri]
            value = {"generation": str(self.generation), "size": str(len(content)), "md5_hash": base64.b64encode(hashlib.md5(content).digest()).decode()}
        elif args[:2] == ["storage", "cp"]:
            source, destination = args[2:4]
            if source.startswith("gs://"):
                uri, generation = source.rsplit("#", 1)
                Path(destination).write_bytes(self.versions[uri, int(generation)])
            else:
                condition = int(next(x.split("=", 1)[1] for x in args if x.startswith("--if-generation-match=")))
                if destination == self.uri:
                    self.writes += 1
                    if self.conflict_at_write == self.writes:
                        self.generation += 1
                        self.versions[self.uri, self.generation] = self.objects[self.uri]
                    if condition != self.generation:
                        code = 1
                    else:
                        content = Path(source).read_bytes()
                        self.generation += 1
                        self.objects[self.uri] = content
                        self.versions[self.uri, self.generation] = content
                elif condition != 0 or destination in self.objects:
                    code = 1
                else:
                    self.objects[destination] = Path(source).read_bytes()
        elif args[:3] == ["pubsub", "topics", "publish"]:
            self.publishes += 1
            mirror = cwd / "fake-provider-check.sqlite"
            mirror.write_bytes(self.objects[self.uri])
            db = sqlite3.connect(mirror)
            try:
                self.persisted_before_publish = db.execute("SELECT status FROM finops_outbox").fetchone() == ("CLAIMED_RECONCILE_BEFORE_RETRY",)
            finally:
                db.close()
            if self.publish_timeout:
                raise ValueError("simulated acknowledgement timeout")
            value = {"messageIds": ["fake-message-1"]}
        else:
            raise AssertionError("unrecognized generated command: " + str(args))
        return subprocess.CompletedProcess(argv, code, json.dumps(value).encode(), b"")


class CycleTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.root = Path(self.temp.name)
        self.config = fixture_config()
        self.seed = self.root / "seed.sqlite"
        f.init_journal(self.seed)
        self.fake = FakeProvider(self.seed, self.config["journal_uri"])
        fixture_binary = str(Path(sys.executable).resolve())
        tool = {"binary": fixture_binary, "binary_sha256": f.file_hash(fixture_binary), "method": "SIMULATED_PROVIDER", "admission": "FIXTURE_ONLY"}
        self.tools = {"bq": dict(tool), "gcloud": dict(tool)}
        self.runner = {"path": __file__, "sha256": f.file_hash(__file__)}
        self.count = 0

    def tearDown(self):
        self.temp.cleanup()

    def approval(self):
        return {"schema": "elite-finops-approval/v403.1", "method": "SIMULATED_PROVIDER", "config_sha256": f.digest(self.config),
                "runtime_sha256": f.runtime_identity(self.tools, self.runner), "approved_by": "FIXTURE_OPERATOR",
                "valid_from": "2026-09-14T00:00:00Z", "expires_at": "2026-09-15T00:00:00Z",
                "billing_table_account_verified": True, "query_cost_authorized": True, "durable_state_authorized": True,
                "dispatch_authorized": self.config["notification"]["mode"] == "PUBSUB"}

    def run_cycle(self, now="2026-09-14T12:00:00Z", approval=None, fixture=True):
        self.count += 1
        with patch.object(f, "load_runner", return_value=self.fake):
            return f.cycle(self.config, approval or self.approval(), now, self.tools, self.runner, self.root / ("run-" + str(self.count)), fixture=fixture)

    def test_aggregate_durable_then_fake_publish_and_next_snapshot_replaces(self):
        receipt = self.run_cycle()
        report = receipt["report"]
        self.assertEqual(report["observed_micros"], 1400000)
        self.assertEqual(report["state"], "ALERT")
        self.assertTrue(report["alert_persisted"])
        self.assertFalse(report["hard_cap"])
        self.assertFalse(receipt["human_delivery_proven"])
        self.assertEqual(report["dispatch_status"], "PUBLISHED_NOT_DELIVERED")
        self.assertTrue(self.fake.persisted_before_publish)
        self.assertEqual(receipt["method"], "SIMULATED_PROVIDER")
        query = next(x for x in self.fake.calls if "query" in x)
        self.assertIn("--maximum_bytes_billed=10000000", query)
        self.assertIn("--max_rows=100001", query)
        self.assertIn("invoice.month = '202609'", query[-1])
        self.assertIn("UNNEST(credits)", query[-1])
        self.assertTrue(any(x.startswith("--bigqueryrc=") for x in query))
        self.fake.rows[0]["amount_micros"] = "1600000"
        second = self.run_cycle("2026-09-14T12:05:00Z")
        self.assertEqual(second["report"]["observed_micros"], 1500000)
        self.assertEqual(second["report"]["dispatch_status"], "DEDUPLICATED_PUBLISHED_NOT_DELIVERED")
        self.assertEqual(self.fake.publishes, 1)

    def test_same_slot_fence_prevents_second_billable_query(self):
        self.run_cycle()
        with self.assertRaisesRegex(ValueError, "reconcile"):
            self.run_cycle("2026-09-14T12:00:01Z")
        self.assertEqual(self.fake.queries, 1)
        self.assertEqual(self.fake.publishes, 1)

    def test_concurrent_journal_cas_rejects_without_publish(self):
        self.fake.conflict_at_write = 1
        with self.assertRaisesRegex(ValueError, "reconcile"):
            self.run_cycle()
        self.assertEqual(self.fake.publishes, 0)
        self.assertEqual(self.fake.objects[self.config["journal_uri"]], self.seed.read_bytes())

    def test_publish_timeout_keeps_durable_claim_and_never_automatically_replays(self):
        self.fake.publish_timeout = True
        receipt = self.run_cycle()
        self.assertEqual(receipt["report"]["dispatch_status"], "CLAIMED_RECONCILE_BEFORE_RETRY")
        receipt = self.run_cycle("2026-09-14T12:05:00Z")
        self.assertEqual(receipt["report"]["dispatch_status"], "DEDUPLICATED_CLAIMED_RECONCILE_BEFORE_RETRY")
        self.assertEqual(self.fake.publishes, 1)

    def test_post_publish_cas_conflict_preserves_prior_claim(self):
        self.fake.conflict_at_write = 2
        with self.assertRaisesRegex(ValueError, "reconcile"):
            self.run_cycle()
        self.fake.conflict_at_write = None
        receipt = self.run_cycle("2026-09-14T12:05:00Z")
        self.assertEqual(receipt["report"]["dispatch_status"], "DEDUPLICATED_CLAIMED_RECONCILE_BEFORE_RETRY")
        self.assertEqual(self.fake.publishes, 1)

    def test_none_mode_persists_without_publish(self):
        self.config["notification"] = {"mode": "NONE"}
        self.assertTrue(self.run_cycle()["report"]["alert_persisted"])
        self.assertEqual(self.fake.publishes, 0)

    def test_config_limits_types_prod_sql_and_scope_rejected_before_effect(self):
        cases = [("maximum_bytes_billed", True), ("maximum_bytes_billed", 1.5), ("maximum_bytes_billed", 0),
                 ("maximum_bytes_billed", "100"), ("limit_micros", True), ("production_authorized", True),
                 ("environment", "production"), ("billing_table", "project.dataset.table; DROP TABLE x"),
                 ("journal_uri", "gs://elite-state/../billing.sqlite"), ("period_policy", "ALL_TIME")]
        for key, value in cases:
            with self.subTest(key=key, value=value):
                config = copy.deepcopy(self.config)
                config[key] = value
                with self.assertRaises(ValueError):
                    f.validate_config(config)
        self.assertEqual(self.fake.calls, [])

    def test_export_cap_mixed_currency_duplicate_and_amount_types_rejected(self):
        invalid = [[dict(self.fake.rows[0], amount_micros=True)], [dict(self.fake.rows[0], amount_micros=1.5)],
                   [dict(self.fake.rows[0], amount_micros=str(2**63))], [dict(self.fake.rows[0], currency="EUR")],
                   [dict(self.fake.rows[0]), dict(self.fake.rows[0])], [self.fake.rows[0]] * 100000]
        for number, rows in enumerate(invalid):
            with self.subTest(number=number):
                self.fake.rows = rows
                with self.assertRaises(ValueError):
                    self.run_cycle("2026-09-14T13:" + str(number * 5).zfill(2) + ":00Z")
        self.assertEqual(self.fake.writes, 0)
        self.assertEqual(self.fake.publishes, 0)

    def test_approval_expiry_hash_method_and_dispatch_rejected_before_effect(self):
        for key, value in [("expires_at", "2026-09-14T00:00:00Z"), ("config_sha256", "0" * 64),
                           ("runtime_sha256", "0" * 64), ("method", "EXECUTED_TARGET"),
                           ("dispatch_authorized", False), ("query_cost_authorized", False),
                           ("billing_table_account_verified", False)]:
            with self.subTest(key=key):
                approval = self.approval()
                approval[key] = value
                with self.assertRaises(ValueError):
                    self.run_cycle(approval=approval)
        with self.assertRaises(ValueError):
            self.run_cycle(fixture=False)
        self.assertEqual(self.fake.calls, [])

    def test_seed_is_local_and_refuses_overwrite(self):
        with self.assertRaises(ValueError):
            f.init_journal(self.seed)
        self.assertEqual(self.fake.calls, [])


if __name__ == "__main__":
    unittest.main()
