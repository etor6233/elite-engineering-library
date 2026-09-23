"""AUTHORED non-production FinOps orchestration over existing billing/runner owners.

Only generated read-only billing SQL, generation-fenced state writes and an
explicitly approved Pub/Sub publication are executable. No restriction API.
"""
from __future__ import annotations

import argparse
import base64
import hashlib
import json
import os
from pathlib import Path
import re
import sqlite3
from datetime import datetime, timezone

from cloud_control import (BillingJournal, billing_query, normalize_billing_rows,
                           digest, file_hash, raw, read, require, timestamp)
from execute_step import load_runner


CONFIG_KEYS = {"schema", "production_authorized", "environment", "project", "location",
               "billing_account", "billing_table", "currency", "period_policy",
               "maximum_bytes_billed", "interval_seconds", "max_state_bytes",
               "max_age_seconds", "limit_micros", "journal_uri", "fence_uri", "notification"}
URI = re.compile(r"gs://[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]/[A-Za-z0-9_./-]+")


def integer(value, minimum, maximum, message):
    require(type(value) is int and minimum <= value <= maximum, message)


def validate_config(config):
    require(type(config) is dict and set(config) == CONFIG_KEYS, "exact cycle configuration required")
    require(config["schema"] == "elite-finops-cycle/v403.1", "cycle schema")
    require(config["production_authorized"] is False and config["environment"] in ("development", "staging"), "non-production only")
    require(re.fullmatch(r"[a-z][a-z0-9-]{4,61}[a-z0-9]", config["project"]), "project identity")
    require(re.fullmatch(r"[A-Za-z][A-Za-z0-9-]{1,31}", config["location"]), "billing location")
    require(re.fullmatch(r"[A-Z0-9]{6}-[A-Z0-9]{6}-[A-Z0-9]{6}", config["billing_account"]), "billing account identity")
    require(re.fullmatch(r"[A-Z]{3}", config["currency"]), "currency identity")
    require(config["period_policy"] == "UTC_CURRENT_MONTH", "explicit current invoice month policy")
    billing_query(config["billing_table"], "202601")  # Existing owner's SQL identifier validation.
    integer(config["maximum_bytes_billed"], 1, 10**15, "explicit integer query byte limit required")
    integer(config["interval_seconds"], 60, 86400, "cycle interval bounds")
    integer(config["max_age_seconds"], 1, 86400, "observation age bounds")
    integer(config["max_state_bytes"], 65536, 64 * 1024 * 1024, "state size bounds")
    integer(config["limit_micros"], 1, 2**63 - 1, "integer alert threshold required")
    for key in ("journal_uri", "fence_uri"):
        value = config[key]
        require(isinstance(value, str) and URI.fullmatch(value) and not any(p in ("", ".", "..") for p in value[5:].split("/")), "exact safe GCS URI required")
    require(config["journal_uri"].endswith(".sqlite") and not config["fence_uri"].endswith(".sqlite"), "separate journal and fence prefixes")
    n = config["notification"]
    require(type(n) is dict and n.get("mode") in ("NONE", "PUBSUB"), "notification mode")
    require(set(n) == ({"mode"} if n["mode"] == "NONE" else {"mode", "topic", "dedupe_seconds"}), "exact notification fields")
    if n["mode"] == "PUBSUB":
        require(re.fullmatch(r"projects/" + re.escape(config["project"]) + r"/topics/[A-Za-z][A-Za-z0-9._~-]{2,254}", n["topic"]), "approved project topic required")
        integer(n["dedupe_seconds"], config["interval_seconds"], 31 * 86400, "explicit alert deduplication window required")


def safe_path(path):
    path = Path(path).absolute()
    for part in (path, *path.parents):
        require(not part.is_symlink() and not (hasattr(part, "is_junction") and part.is_junction()), "linked controller path rejected")
    return path


def initialize_tables(journal):
    journal.db.executescript("""
        CREATE TABLE IF NOT EXISTS finops_cycles (
          id TEXT PRIMARY KEY, config_sha TEXT NOT NULL, report TEXT NOT NULL);
        CREATE TABLE IF NOT EXISTS finops_outbox (
          id TEXT PRIMARY KEY, status TEXT NOT NULL, payload TEXT NOT NULL);
    """)
    journal.db.commit()


def init_journal(path):
    """Offline seed only; uploading with generation-match=0 is a separate approved step."""
    path = safe_path(path)
    require(not path.exists(), "journal seed destination already exists")
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("xb"):
        pass
    journal = BillingJournal(path)
    try:
        initialize_tables(journal)
    finally:
        journal.close()
    os.chmod(path, 0o600)
    return {"schema": "elite-finops-seed/v403.1", "sha256": file_hash(path), "provider_executed": False}


def runtime_identity(tools, runner):
    return digest({"tools": tools, "runner": runner})


def cycle(config, approval, now, tools, runner, workdir, *, fixture=False):
    """Run once, returning a local receipt; remote intent survives every ambiguous failure.

    Real runs require exact tool runtime admission and a hash-bound, expiring
    approval. Tests inject the existing runner loader, never change method claims.
    """
    validate_config(config)
    instant = timestamp(now)
    method = "SIMULATED_PROVIDER" if fixture else "EXECUTED_TARGET"
    require(approval.get("schema") == "elite-finops-approval/v403.1" and approval.get("method") == method, "approval method/schema")
    require(approval.get("config_sha256") == digest(config) and approval.get("runtime_sha256") == runtime_identity(tools, runner), "approval config/runtime identity")
    require(timestamp(approval["valid_from"]) <= instant <= timestamp(approval["expires_at"]), "approval outside authorized window")
    require(isinstance(approval.get("approved_by"), str) and approval["approved_by"].strip(), "explicit approver required")
    for key in ("billing_table_account_verified", "query_cost_authorized", "durable_state_authorized"):
        require(approval.get(key) is True, "approval incomplete: " + key)
    require(approval.get("dispatch_authorized") is (config["notification"]["mode"] == "PUBSUB"), "dispatch approval mismatch")
    require(set(tools) == {"bq", "gcloud"}, "exact bq/gcloud tool set")
    for tool in tools.values():
        binary = safe_path(tool["binary"])
        require(tool.get("method") == method and file_hash(binary) == tool["binary_sha256"], "tool method or bytes mismatch")
        require(tool.get("admission") == ("FIXTURE_ONLY" if fixture else "ADMITTED_EXACT_ARTIFACT"), "tool admission required")
        require(not tool.get("prefix_arguments"), "tool prefix arguments rejected")
        if not fixture:
            from tool_admission import validate
            validate(tool)
    require(file_hash(safe_path(runner["path"])) == runner["sha256"], "bounded runner hash mismatch")
    bounded = load_runner(runner["path"], runner["sha256"])
    root = safe_path(workdir)
    require(not root.exists(), "fresh private working directory required")
    root.mkdir(mode=0o700, parents=True)
    os.chmod(root, 0o700)
    (root / "empty.bigqueryrc").write_bytes(b"")
    argv_hashes = []

    def run(name, arguments, *, json_output=False):
        tool = tools[name]
        allowed = {"PATH", "SYSTEMROOT", "WINDIR", "TEMP", "TMP", "LANG", "LC_ALL", "SSL_CERT_FILE", "SSL_CERT_DIR"}
        names = tool.get("environment_variable_names", [])
        require(type(names) is list and all(isinstance(n, str) and re.fullmatch(r"[A-Z_][A-Z0-9_]*", n) for n in names), "environment names only")
        require(not {"GH_TOKEN", "GITHUB_TOKEN"}.intersection(names), "repository credentials forbidden")
        allowed.update(names)
        env = {k: v for k, v in os.environ.items() if k in allowed}
        arguments = [str(safe_path(tool["binary"])), *arguments]
        argv_hashes.append(digest(arguments))
        result = bounded(arguments, root, env, 180)
        require(result.returncode == 0, "provider command failed; reconcile durable intent before any replay")
        if json_output:
            return json.loads(result.stdout)
        return None

    gflags = ["--project=" + config["project"], "--quiet"]
    slot = int(instant.timestamp()) // config["interval_seconds"]
    cycle_id = digest({"config_sha256": digest(config), "slot": slot})
    intent = {"schema": "elite-finops-intent/v403.1", "cycle_id": cycle_id, "method": method,
              "config_sha256": digest(config), "runtime_sha256": runtime_identity(tools, runner),
              "state": "CLAIMED_RECONCILE_BEFORE_RETRY", "observed_at": now}
    intent_path = root / "intent.json"
    intent_path.write_bytes(raw(intent))
    run("gcloud", ["storage", "cp", str(intent_path), config["fence_uri"] + "/" + cycle_id + ".json", "--if-generation-match=0", *gflags])

    def metadata():
        value = run("gcloud", ["storage", "objects", "describe", config["journal_uri"], "--format=json", *gflags], json_output=True)
        require(type(value) is dict and type(value.get("generation")) in (str, int) and re.fullmatch(r"[1-9][0-9]*", str(value["generation"])), "journal generation missing")
        require(type(value.get("size")) in (str, int) and re.fullmatch(r"[0-9]+", str(value["size"])), "journal size missing")
        integer(int(value["size"]), 1, config["max_state_bytes"], "journal size limit")
        require(isinstance(value.get("md5_hash"), str), "journal transport checksum required")
        return value

    before = metadata()
    local = root / "billing.sqlite"
    run("gcloud", ["storage", "cp", config["journal_uri"] + "#" + str(before["generation"]), str(local), *gflags])
    require(local.is_file() and not local.is_symlink() and local.stat().st_size == int(before["size"]), "journal download size mismatch")
    require(base64.b64encode(hashlib.md5(local.read_bytes()).digest()).decode() == before["md5_hash"], "journal download checksum mismatch")
    os.chmod(local, 0o600)
    journal = BillingJournal(local)
    try:
        journal.db.execute("PRAGMA trusted_schema=OFF")
        require(journal.db.execute("PRAGMA quick_check").fetchone() == ("ok",), "journal integrity failure")
        initialize_tables(journal)
        require(journal.db.execute("SELECT 1 FROM finops_cycles WHERE id=?", (cycle_id,)).fetchone() is None, "cycle already persisted")
        period = instant.strftime("%Y-%m")
        previous = journal.db.execute("SELECT generation FROM snapshots WHERE provider=? AND account=? AND period=?", ("gcp", config["billing_account"], period)).fetchone()
        query = billing_query(config["billing_table"], instant.strftime("%Y%m"))
        args = ["--bigqueryrc=" + str(root / "empty.bigqueryrc"), "--project_id=" + config["project"], "--location=" + config["location"],
                *query[1:-1], "--maximum_bytes_billed=" + str(config["maximum_bytes_billed"]), query[-1]]
        rows = run("bq", args, json_output=True)
        # BQ INT64 JSON strings are allowed; bool/float coercion and aggregate overflow are not.
        require(type(rows) is list and len(rows) < 100000, "complete bounded export required")
        for row in rows:
            require(type(row) is dict and type(row.get("amount_micros")) in (str, int) and re.fullmatch(r"-?[0-9]+", str(row["amount_micros"])), "export amount must be exact INT64")
            require(-(2**63) <= int(row["amount_micros"]) < 2**63, "export amount overflow")
        identity = {"schema": "elite-billing-snapshot/v403.1", "kind": "PROVIDER_BILLING_OBSERVATION", "provider": "gcp", "account": config["billing_account"],
                    "period": period, "generation": previous[0] + 1 if previous else 1, "observed_at": now,
                    "currency": config["currency"], "complete_scope": True, "export_delay_possible": True}
        snapshot = normalize_billing_rows(rows, identity)
        snapshot["source_proof"]["method"] = "SIMULATED_PROVIDER" if fixture else "EXECUTED_PROVIDER_EXPORT"
        journal.import_snapshot(snapshot)
        policy = {"kind": "AGGREGATE_ALERT_POLICY", "period": period, "currency": config["currency"], "limit_micros": config["limit_micros"],
                  "max_age_seconds": config["max_age_seconds"], "scope": [{"provider": "gcp", "account": config["billing_account"]}]}
        report = journal.budget(policy, now)
        report.update({"schema": "elite-finops-report/v403.1", "method": method, "cycle_id": cycle_id, "period": period,
                       "observed_at": now, "export_freshness": "NOT_PROVEN_EXPORT_MAY_BE_DELAYED", "alert_persisted": False,
                       "dispatch_status": "NOT_REQUESTED", "automatic_restriction": False})
        outbox_id = None
        n = config["notification"]
        if n["mode"] == "PUBSUB" and report["state"] == "ALERT":
            outbox_id = digest({"config": digest(config), "period": period, "window": int(instant.timestamp()) // n["dedupe_seconds"]})
            exists = journal.db.execute("SELECT status FROM finops_outbox WHERE id=?", (outbox_id,)).fetchone()
            if exists:
                report["dispatch_status"] = "DEDUPLICATED_" + exists[0]
                outbox_id = None
            else:
                report["dispatch_status"] = "CLAIMED_RECONCILE_BEFORE_RETRY"
                journal.db.execute("INSERT INTO finops_outbox VALUES (?,?,?)", (outbox_id, "CLAIMED_RECONCILE_BEFORE_RETRY", raw(report).decode()))
        report["alert_persisted"] = True  # Becomes observable only after CAS succeeds.
        journal.db.execute("INSERT INTO finops_cycles VALUES (?,?,?)", (cycle_id, digest(config), raw(report).decode()))
        journal.db.commit()

        def persist(expected, number):
            path = root / ("state-" + str(number) + ".sqlite")
            backup = sqlite3.connect(path)
            try:
                journal.db.backup(backup)
                backup.execute("PRAGMA journal_mode=DELETE")
            finally:
                backup.close()
            integer(path.stat().st_size, 1, config["max_state_bytes"], "journal growth limit")
            os.chmod(path, 0o600)
            with path.open("r+b") as handle:
                os.fsync(handle.fileno())
            checksum = base64.b64encode(hashlib.md5(path.read_bytes()).digest()).decode()
            run("gcloud", ["storage", "cp", str(path), config["journal_uri"], "--if-generation-match=" + str(expected), "--content-md5=" + checksum, *gflags])
            after = metadata()
            require(int(after["generation"]) > int(expected) and int(after["size"]) == path.stat().st_size and after["md5_hash"] == checksum,
                    "CAS result changed before observation; reconcile without publication")
            return after

        persisted = persist(before["generation"], 1)
        if outbox_id:
            # The durable claim prevents automatic duplicate publish after a timeout/crash.
            # No adapter retries. Pub/Sub/transport may duplicate; subscribers must deduplicate payload.id.
            payload = {"schema": "elite-finops-alert/v403.1", "id": outbox_id, "period": period, "currency": report["currency"],
                       "observed_micros": report["observed_micros"], "limit_micros": report["limit_micros"], "hard_cap": False,
                       "method": method, "export_delay_possible": True}
            status = "CLAIMED_RECONCILE_BEFORE_RETRY"
            try:
                published = run("gcloud", ["pubsub", "topics", "publish", n["topic"], "--message=" + json.dumps(payload, sort_keys=True),
                                           "--format=json", *gflags], json_output=True)
                require(type(published) is dict and type(published.get("messageIds")) is list and published["messageIds"] and
                        all(isinstance(x, str) and x for x in published["messageIds"]), "publication acknowledgement missing")
                status = "PUBLISHED_NOT_DELIVERED"
            except (ValueError, OSError):
                pass  # Persisted CLAIMED remains; no automatic retry on ambiguous publish.
            report["dispatch_status"] = status
            journal.db.execute("UPDATE finops_outbox SET status=? WHERE id=?", (status, outbox_id))
            journal.db.execute("UPDATE finops_cycles SET report=? WHERE id=?", (raw(report).decode(), cycle_id))
            journal.db.commit()
            persisted = persist(persisted["generation"], 2)
        receipt = {"schema": "elite-finops-cycle-receipt/v403.1", "method": method, "config_sha256": digest(config),
                   "runtime_sha256": runtime_identity(tools, runner), "report": report, "query_sha256": digest(args),
                   "snapshot_sha256": digest(snapshot), "journal_generation": str(persisted["generation"]), "argv_sha256": argv_hashes,
                   "production_authorized": False, "human_delivery_proven": False}
        (root / "receipt.json").write_bytes(raw(receipt))
        return receipt
    finally:
        journal.close()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest="command", required=True)
    init = sub.add_parser("init")
    init.add_argument("--output", required=True)
    run = sub.add_parser("run")
    for name in ("config", "approval", "tools", "runner", "workdir"):
        run.add_argument("--" + name, required=True)
    run.add_argument("--authorize-execution", action="store_true")
    args = parser.parse_args()
    if args.command == "init":
        result = init_journal(args.output)
    else:
        require(args.authorize_execution, "explicit execution authorization required")
        now = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
        result = cycle(read(args.config), read(args.approval), now, read(args.tools), read(args.runner), args.workdir)
    print(json.dumps(result, sort_keys=True))


if __name__ == "__main__":
    main()
