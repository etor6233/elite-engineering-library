"""AUTHORED deployment/evidence glue over existing library owners, not product code.

Offline planning is the default. No provider command is ever executed here.
The generated argv plans are consumed by the existing official-tool runner only
after an operator has approved the exact plan and target. SQLite is a single
controller's receipt journal, never the franchise's database or tenant budget.
"""
from __future__ import annotations
import argparse
import hashlib
import json
import re
import sqlite3
from datetime import datetime, timezone
from pathlib import Path

HEX = re.compile(r"[0-9a-f]{64}")
ID = re.compile(r"[a-z][a-z0-9-]{2,61}[a-z0-9]")
IMAGE = re.compile(r"[a-z0-9-]+-docker\.pkg\.dev/[a-z][a-z0-9-]+/[a-z0-9-]+/[a-z0-9/-]+@sha256:[0-9a-f]{64}")

def require(value, message):
    if not value:
        raise ValueError(message)

def raw(value):
    return (json.dumps(value, sort_keys=True, indent=2) + "\n").encode()

def digest(value):
    return hashlib.sha256(raw(value)).hexdigest()

def file_hash(path):
    return hashlib.sha256(Path(path).read_bytes()).hexdigest()

def read(path):
    return json.loads(Path(path).read_text(encoding="utf-8"))

def contained(root, name):
    require(isinstance(name, str) and name and "\\" not in name and ":" not in name, "relative path required")
    require(not name.startswith("/") and all(p not in ("", ".", "..") for p in name.split("/")), "unsafe relative path")
    root = Path(root).resolve()
    p = root / name
    for part in (p, *p.parents):
        require(not part.is_symlink() and not (hasattr(part, "is_junction") and part.is_junction()), "linked path rejected")
        if part == root:
            break
    p.resolve().relative_to(root)
    return p

def timestamp(value):
    require(isinstance(value, str) and value.endswith("Z"), "UTC Z timestamp required")
    return datetime.fromisoformat(value.replace("Z", "+00:00"))

def validate_profile(profile, library):
    require(profile["schema"] == "elite-cloud-composition/v403.1", "composition schema")
    require(profile["scope"] == "LIBRARY_INFRASTRUCTURE", "library scope only")
    require(profile["production_authorized"] is False, "production remains unauthorized")
    selection = profile["selection"]
    planpath = contained(library, selection["path"])
    require(file_hash(planpath) == selection["sha256"], "full composition changed")
    text = planpath.read_text(encoding="utf-8")
    plan = json.loads(re.search(r"```json\s*\n(.*?)\n```", text, re.S).group(1))
    if profile.get("extension_policy") == "BASE116_PLUS_EXACT_FOUR_V403":
        from seal_consumer import parse_plan,validate_selection
        base=profile["base_selection"];basepath=contained(library,base["path"])
        require(file_hash(basepath)==base["sha256"],"base selection changed")
        validate_selection(parse_plan(basepath.read_bytes()),plan)
    require(len(plan["packs"]) == selection["pack_count"], "pack count mismatch")
    require(len({p["packId"] for p in plan["packs"]}) == len(plan["packs"]), "duplicate pack")
    require(profile["pack_locks"] and len(profile["pack_locks"]) == len(plan["packs"]), "incomplete full selection")
    byid = {p["packId"]: p for p in plan["packs"]}
    require(len({p["packId"] for p in profile["pack_locks"]}) == len(profile["pack_locks"]) and {p["packId"] for p in profile["pack_locks"]} == set(byid), "pack lock identity set must exactly equal complete selection")
    for entry in profile["pack_locks"]:
        require(entry["packId"] in byid and entry["path"] == byid[entry["packId"]]["path"] and entry["version"] == byid[entry["packId"]]["version"], "selection identity changed")
        require(file_hash(contained(library, entry["path"])) == entry["sha256"], "selected pack changed: " + entry["packId"])
    for path, expected in profile["authority_locks"].items():
        require(file_hash(contained(library, path)) == expected, "authority changed: " + path)
    return {"result": "PASS", "claim": "complete selected source profile integrity", "pack_count": len(byid), "composition_sha256": digest(profile), "cloud_executed": False}

def target(config):
    require(config["schema"] == "elite-cloud-target/v403.1", "target schema")
    require(config["environment"] in ("development", "staging"), "production requires separate existing production admission")
    require(config["platform"] == "linux/amd64", "target must be Linux amd64")
    for key in ("project", "service", "database_instance"):
        require(ID.fullmatch(config[key]), "invalid target " + key)
    require(re.fullmatch(r"[a-z]+-[a-z]+[0-9]", config["region"]), "invalid region")
    require(re.fullmatch(r"[a-z][a-z0-9-]{4,28}[a-z0-9]@" + re.escape(config["project"]) + r"\.iam\.gserviceaccount\.com", config["runtime_identity"]), "explicit service identity required")
    require(config["private_url"] is True and config["production_authorized"] is False, "private reference only")
    require(config["max_instances"] >= config["min_instances"] >= 1, "embedded background loops require minimum one API instance")
    require(config["max_instances"] <= 10 and config["cpu_always"] is True, "bounded instance-based API CPU required")
    require(config["database_policy"] == "RETAIN_BACKUP_PITR", "database must be retained")
    require(config["state_backend"] == "GCS_GENERATION_CAS", "ephemeral cloud controller journal forbidden")
    require(re.fullmatch(r"gs://[a-z0-9][a-z0-9._-]+/[a-zA-Z0-9/_-]+\.json", config["state_uri"]), "versioned state URI required")
    require(set(config["secrets"]) >= {"DATABASE_URL"}, "database secret reference required")
    for name, value in config["secrets"].items():
        require(re.fullmatch(r"[A-Z][A-Z0-9_]+", name) and re.fullmatch(r"[a-z][a-z0-9-]+:[1-9][0-9]*", value), "only explicit Secret Manager versions, no values/latest")
    require(set(config["environment_values"]) >= {"OIDC_ISSUER", "OIDC_AUDIENCE"}, "identity runtime values required")
    require(config["environment_values"]["OIDC_ISSUER"].startswith("https://"), "HTTPS issuer required")
    for name, value in config["environment_values"].items():
        require(re.fullmatch(r"[A-Z][A-Z0-9_]+", name), "invalid environment key")
        require(isinstance(value, str) and value and all(c not in value for c in "\r\n,\x00") and not any(s in name for s in ("SECRET", "TOKEN", "PASSWORD", "PRIVATE_KEY")), "secret or ambiguous env must use reference")
    return config

def release(value, config):
    require(value["schema"] == "elite-cloud-release/v403.1", "release schema")
    require(value["target"] == "linux/amd64" and value["environment_class"] == "TARGET_ARTIFACT", "local Windows artifact forbidden")
    require(IMAGE.fullmatch(value["image"]), "immutable Artifact Registry image required")
    require(value["image"].split("/")[1] == config["project"], "image project mismatch")
    require(re.fullmatch(r"[0-9a-f]{40}", value["source_commit"]), "exact source revision required")
    require(HEX.fullmatch(value["composition_sha256"]), "exact composition digest required")
    require(set(value["required_receipts"]) == {"build", "license", "sca", "secret_scan", "migrations", "restore", "runtime", "identity"}, "all required artifact/target receipts required")
    for key, r in value["required_receipts"].items():
        require(r["result"] == "PASS" and HEX.fullmatch(r["sha256"]) and r["image"] == value["image"] and r["target"] == value["target"], "receipt mismatch: " + key)
        require(r["method"] == "EXECUTED_TARGET", "fixture receipt cannot authorize target: " + key)
        require(isinstance(r["path"], str), "receipt path required")
    return value

def bind_receipts(value, root):
    for key, ref in value["required_receipts"].items():
        p = contained(root, ref["path"])
        require(file_hash(p) == ref["sha256"], "tampered receipt: " + key)
        r = read(p)
        require(r["image"] == value["image"] and r["target"] == value["target"] and r["result"] == "PASS" and r["method"] == "EXECUTED_TARGET", "unbound target receipt: " + key)

def deployment_commands(config, image, prior_revision):
    artifact={"image":image}
    suffix=image.split(":")[-1][:12]
    revision=config["service"]+"-"+suffix
    common = ["--project", config["project"], "--region", config["region"], "--quiet"]
    envs = dict(config["environment_values"], HTTP_ADDRESS=":8080", ARCA_ENABLED="false")
    commands = [
      {"id": "describe-before", "effect": "READ", "argv": ["gcloud", "run", "services", "describe", config["service"], *common, "--format=json"]},
      {"id": "deploy-no-traffic", "effect": "WRITE", "argv": ["gcloud", "run", "deploy", config["service"], *common, "--image", artifact["image"], "--revision-suffix", suffix, "--no-traffic", "--no-allow-unauthenticated", "--ingress", "internal-and-cloud-load-balancing", "--service-account", config["runtime_identity"], "--min-instances", str(config["min_instances"]), "--max-instances", str(config["max_instances"]), "--no-cpu-throttling", "--port", "8080", "--timeout", "60s", "--set-secrets", ",".join(k+"="+v for k,v in sorted(config["secrets"].items())), "--set-env-vars", ",".join(k+"="+v for k,v in sorted(envs.items()))]},
      {"id": "describe-candidate", "effect": "READ", "argv": ["gcloud", "run", "revisions", "describe", revision, *common, "--format=json"]},
      {"id": "promote-after-health", "effect": "WRITE", "requires": "authenticated candidate health and migration receipt", "argv": ["gcloud", "run", "services", "update-traffic", config["service"], *common, "--to-revisions", revision+"=100"]},
      {"id": "rollback", "effect": "WRITE", "requires": "rollback approval and previous digest health", "argv": ["gcloud", "run", "services", "update-traffic", config["service"], *common, "--to-revisions", prior_revision+"=100"]},
      {"id": "describe-after", "effect": "READ", "argv": ["gcloud", "run", "services", "describe", config["service"], *common, "--format=json"]},
    ]
    return commands

def deploy_plan(config, artifact, prior_revision, evidence_root, before_snapshot="before.json"):
    target(config); release(artifact, config); bind_receipts(artifact, evidence_root)
    require(re.fullmatch(re.escape(config["service"]) + r"-[a-z0-9-]+", prior_revision), "known previous revision required")
    before_path=contained(evidence_root,before_snapshot);before=read(before_path)
    require(before["project"]==config["project"] and before["region"]==config["region"] and before["service"]==config["service"],"before snapshot target mismatch")
    require(before["revision"]==prior_revision and before["traffic"]=={prior_revision:100} and IMAGE.fullmatch(before["image"]),"before snapshot exact prior image and traffic required")
    require(before["method"]=="EXECUTED_TARGET", "before snapshot must be observed target; synthetic data only in isolated test")
    suffix = artifact["image"].split(":")[-1][:12]
    revision = config["service"] + "-" + suffix
    require(revision != prior_revision, "new release equals prior revision")
    commands=deployment_commands(config,artifact["image"],prior_revision)
    return {"schema": "elite-cloud-command-plan/v403.1", "scope":"LIBRARY_INFRASTRUCTURE", "environment":config["environment"], "execution": "NOT_RUN", "production_authorized": False, "target_config":config, "target_sha256": digest(config), "release_sha256": digest(artifact), "image": artifact["image"], "prior_revision": prior_revision, "prior_image":before["image"], "before_snapshot_sha256":file_hash(before_path), "candidate_revision": revision, "execution_fence":{"backend":"GCS_GENERATION_0","uri":config["state_uri"][:-5]+"-intents"}, "evidence_root":str(Path(evidence_root).resolve()), "commands": commands, "mandatory_manual_gates": ["target IAM and protected URL acceptance", "previous revision digest captured", "Cloud Run 10s termination versus current 15s host shutdown: forced stop/retry recovery demonstrated", "no schema rollback without expand-contract compatibility", "separate frontend image and authenticated API transport"]}

def execution_plan(plan):
    require(plan.get("schema")=="elite-cloud-command-plan/v403.1" and plan.get("scope")=="LIBRARY_INFRASTRUCTURE" and plan.get("production_authorized") is False and plan.get("environment") in ("development","staging") and plan.get("execution")=="NOT_RUN","invalid non-production execution plan")
    require(HEX.fullmatch(plan["target_sha256"]) and HEX.fullmatch(plan["release_sha256"]) and IMAGE.fullmatch(plan["image"]) and IMAGE.fullmatch(plan["prior_image"]),"execution plan identities")
    require([c["id"] for c in plan["commands"]]==["describe-before","deploy-no-traffic","describe-candidate","promote-after-health","rollback","describe-after"],"exact ordered execution steps required")
    require(plan["execution_fence"]["backend"] in ("GCS_GENERATION_0","LOCAL_FIXTURE"),"execution fence missing")
    target(plan["target_config"]);require(digest(plan["target_config"])==plan["target_sha256"],"target configuration drift")
    require(plan["commands"]==deployment_commands(plan["target_config"],plan["image"],plan["prior_revision"]),"execution commands/effects differ from exact generated contract")
    return plan

def authorize(plan, approval, now):
    execution_plan(plan)
    return authorize_record(plan,approval,now)

def authorize_record(plan,approval,now):
    """Shared approval fields; callers validate their exact plan schema first."""
    require(approval["plan_sha256"] == digest(plan), "approval must name exact plan")
    require(approval["scope"] == "CLOUD_REFERENCE_EXECUTION" and approval["environment"] in ("development", "staging"), "approval scope")
    require(approval["environment"]==plan["environment"],"approval target environment mismatch")
    require(isinstance(approval["approved_by"],str) and approval["approved_by"].strip() and approval["approved_by"] != "fixture", "human approval identity required")
    require(timestamp(approval["approved_at"]) <= timestamp(now) <= timestamp(approval["expires_at"]), "approval expired or not yet valid")
    require(approval["provider_access_proven"] is True and approval["budget_authorized"] is True, "account and budget approval required")
    return {"result": "PASS", "claim": "authorization record consistency only", "execution": "NOT_RUN", "plan_sha256": digest(plan)}

def reconcile_deployment(plan, observation):
    require(observation["schema"] == "elite-cloud-observation/v403.1", "observation schema")
    require(observation["plan_sha256"] == digest(plan), "observation plan mismatch")
    require(observation["method"] in ("SIMULATED_PROVIDER", "EXECUTED_TARGET"), "observation method")
    wanted = plan["candidate_revision"] if observation["operation"] == "promote" else plan["prior_revision"]
    require(observation["operation"] in ("promote", "rollback"), "operation")
    require(observation["traffic"] == {wanted: 100}, "traffic does not match requested operation")
    require(all(observation[k] is True for k in ("ready","identity_checked","data_health","pending_work_reconciled")), "incomplete observed health")
    if observation["operation"] == "promote":
        require(observation["image"] == plan["image"], "deployed digest drift")
    else:
        require(IMAGE.fullmatch(observation["image"]) and observation["image"] == plan["prior_image"], "prior immutable image not restored")
    return {"result": "PASS", "claim": "deployment observation reconciliation", "method": observation["method"], "cloud_demonstrated": observation["method"] == "EXECUTED_TARGET", "operation": observation["operation"], "observation_sha256": digest(observation)}

class BillingJournal:
    """Single-writer durable observations, distinct from tenant GO-FINOPS owner.

Snapshots replace a complete account/month atomically; never add a repeatedly
observed aggregate. Export delay and provider credit adjustments remain visible.
The backing file must be persisted by the GCS CAS adapter on cloud workers.
"""
    def __init__(self, path):
        self.db = sqlite3.connect(path, timeout=10)
        self.db.execute("PRAGMA journal_mode=WAL")
        self.db.execute("PRAGMA foreign_keys=ON")
        self.db.executescript("CREATE TABLE IF NOT EXISTS snapshots(provider TEXT, account TEXT, period TEXT, generation INTEGER, observed_at TEXT, sha TEXT, currency TEXT, payload TEXT, PRIMARY KEY(provider,account,period)); CREATE TABLE IF NOT EXISTS internal_costs(id TEXT PRIMARY KEY, currency TEXT, amount_micros INTEGER, payload_sha TEXT);")

    def close(self):
        self.db.close()

    def import_snapshot(self, item):
        require(item["schema"] == "elite-billing-snapshot/v403.1", "billing schema")
        require(item["kind"] == "PROVIDER_BILLING_OBSERVATION", "internal estimate is not provider billing")
        require(item["complete_scope"] is True and item["export_delay_possible"] is True, "complete scope and export latency required")
        proof=item["source_proof"]
        require(proof["truncated"] is False and proof["row_count"]==len(item["rows"]) and 0<=proof["row_count"]<100000 and proof["method"] in ("SIMULATED_PROVIDER","EXECUTED_PROVIDER_EXPORT"),"billing export completeness proof required; cap/truncation rejected")
        require(re.fullmatch(r"[A-Z]{3}", item["currency"]), "currency")
        require(re.fullmatch(r"20[0-9]{2}-[01][0-9]", item["period"]) and 1 <= int(item["period"][-2:]) <= 12, "period")
        require(isinstance(item["generation"], int) and item["generation"] > 0, "source generation")
        timestamp(item["observed_at"])
        seen = set()
        for row in item["rows"]:
            require(set(row) == {"project", "service", "amount_micros"}, "exact normalized billing fields")
            require(isinstance(row["amount_micros"], int) and not isinstance(row["amount_micros"], bool), "integer micros required; signed credits permitted")
            key = row["project"], row["service"]
            require(all(isinstance(s, str) and s for s in key) and key not in seen, "duplicate billing aggregation key")
            seen.add(key)
        k = item["provider"], item["account"], item["period"]
        require(all(isinstance(s, str) and s for s in k), "provider scope")
        with self.db:
            self.db.execute("BEGIN IMMEDIATE")
            old = self.db.execute("SELECT generation,sha,observed_at FROM snapshots WHERE provider=? AND account=? AND period=?", k).fetchone()
            if old and old[1] == digest(item):
                return "DUPLICATE_NO_EFFECT"
            require(not old or item["generation"] > old[0] and timestamp(item["observed_at"]) >= timestamp(old[2]), "stale or conflicting billing snapshot")
            self.db.execute("INSERT OR REPLACE INTO snapshots VALUES(?,?,?,?,?,?,?,?)", (*k, item["generation"], item["observed_at"], digest(item), item["currency"], raw(item).decode()))
        return "RECONCILED"

    def record_internal(self, item):
        require(item["kind"] == "INTERNAL_ESTIMATE" and type(item["amount_micros"]) is int and item["amount_micros"] >= 0, "explicit internal nonnegative estimate")
        require(re.fullmatch(r"[A-Z]{3}", item["currency"]), "currency")
        h = digest(item)
        with self.db:
            old = self.db.execute("SELECT payload_sha FROM internal_costs WHERE id=?", (item["id"],)).fetchone()
            require(not old or old[0] == h, "idempotency key payload conflict")
            self.db.execute("INSERT OR IGNORE INTO internal_costs VALUES(?,?,?,?)", (item["id"], item["currency"], item["amount_micros"], h))

    def budget(self, policy, now):
        require(policy["kind"] == "AGGREGATE_ALERT_POLICY", "budget is an alert policy, not a global hard cap")
        require(type(policy["limit_micros"]) is int and policy["limit_micros"] > 0, "positive limit")
        require(policy["scope"] and len({(s["provider"],s["account"]) for s in policy["scope"]}) == len(policy["scope"]), "scope must be nonempty and unique")
        total = 0; missing = []; stale = []
        for scope in policy["scope"]:
            r = self.db.execute("SELECT payload FROM snapshots WHERE provider=? AND account=? AND period=?", (scope["provider"], scope["account"], policy["period"])).fetchone()
            if not r:
                missing.append(scope); continue
            item = json.loads(r[0]); require(item["currency"] == policy["currency"], "cross currency requires separately admitted FX; no implicit conversion")
            age = (timestamp(now)-timestamp(item["observed_at"])).total_seconds()
            if age < 0 or age > policy["max_age_seconds"]:
                stale.append(scope)
            total += sum(row["amount_micros"] for row in item["rows"])
        return {"kind": "AGGREGATE_BILLING_ALERT", "currency": policy["currency"], "observed_micros": total, "limit_micros": policy["limit_micros"], "state": "UNKNOWN" if missing or stale else "ALERT" if total >= policy["limit_micros"] else "WITHIN_OBSERVED_BUDGET", "missing": missing, "stale": stale, "hard_cap": False, "internal_estimates_added_to_billing": False, "residual_costs": ["delayed billing exports", "storage and backups", "networking", "external APIs", "credit adjustments"]}

def billing_query(table, period):
    require(re.fullmatch(r"[a-z][a-z0-9-]+\.[a-zA-Z_][a-zA-Z0-9_]*\.[a-zA-Z_][a-zA-Z0-9_]*", table), "exact billing table identifier")
    require(re.fullmatch(r"20[0-9]{2}(0[1-9]|1[0-2])", period), "invoice month YYYYMM")
    # Standard export: all projects/services, signed credits, single currency
    # aggregation; account-level unallocated cost remains a separate row.
    sql = "SELECT COALESCE(project.id, '_ACCOUNT_LEVEL') AS project, service.id AS service, currency, CAST(ROUND(SUM(cost + IFNULL((SELECT SUM(c.amount) FROM UNNEST(credits) c), 0))*1000000) AS INT64) AS amount_micros FROM `" + table + "` WHERE invoice.month = '" + period + "' GROUP BY project, service, currency"
    return ["bq", "query", "--use_legacy_sql=false", "--format=json", "--max_rows=100001", sql]

def normalize_billing_rows(rows, identity):
    require(isinstance(rows,list) and len(rows)<100000,"billing row cap reached; paginate exact job before import")
    currencies={x["currency"] for x in rows};require(len(currencies)<=1,"mixed currency export requires separate scopes; no implicit FX")
    require(not currencies or identity["currency"] in currencies,"billing currency mismatch")
    out=dict(identity)
    out["rows"]=[]
    for x in rows:
        require(set(x)=={"project","service","currency","amount_micros"} and re.fullmatch(r"-?[0-9]+",str(x["amount_micros"])),"unexpected billing export shape")
        out["rows"].append({"project":x["project"],"service":x["service"],"amount_micros":int(x["amount_micros"])})
    out["source_proof"]={"method":"EXECUTED_PROVIDER_EXPORT","truncated":False,"row_count":len(rows)}
    return out

def state_cas_plan(config, expected_generation, state_file):
    target(config)
    require(type(expected_generation) is int and expected_generation >= 0, "exact GCS generation required")
    require(Path(state_file).is_file(), "new durable state file required")
    return {"expected_generation": expected_generation, "sha256": file_hash(state_file), "argv": ["gcloud", "storage", "cp", str(state_file), config["state_uri"], "--if-generation-match=" + str(expected_generation), "--quiet"], "conflict_action": "RELOAD_AND_RECONCILE; never overwrite unconditionally", "executed": False}

def suspend_plan(config, proof):
    target(config)
    required = ("ingress_closed", "schedules_paused", "webhooks_buffered_or_provider_replay_proven", "new_jobs_fenced", "active_jobs_zero", "outbox_reconciled", "backup_verified", "restore_proven", "resume_recipe_verified")
    require(all(proof.get(k) is True for k in required), "cannot suspend with incomplete drain/recovery proof")
    require(proof["target_sha256"] == digest(config) and proof["method"] in ("SIMULATED_PROVIDER", "EXECUTED_TARGET"), "suspend proof binding")
    return {"schema": "elite-cloud-suspend/v403.1", "method": proof["method"], "execution": "NOT_RUN", "database_deleted": False, "zero_cost_claim": False, "commands": [
      ["gcloud", "run", "services", "update", config["service"], "--project", config["project"], "--region", config["region"], "--scaling=0", "--quiet"],
      ["gcloud", "sql", "instances", "patch", config["database_instance"], "--project", config["project"], "--activation-policy=NEVER", "--quiet"]], "resume_order": ["Cloud SQL activation ALWAYS and observed health", "migrations compatibility and secrets references", "API instances 1 and health", "resume fenced workers and reconcile outbox", "restore scheduler and ingress last"], "retained_billable": ["database storage", "backups", "IP addresses", "object storage", "logs", "external API charges"]}

def main():
    parser=argparse.ArgumentParser(); sub=parser.add_subparsers(dest="cmd", required=True)
    p=sub.add_parser("verify-composition"); p.add_argument("--profile",required=True); p.add_argument("--library",required=True)
    p=sub.add_parser("deploy-plan"); p.add_argument("--target",required=True);p.add_argument("--release",required=True);p.add_argument("--prior-revision",required=True);p.add_argument("--evidence-root",required=True)
    p=sub.add_parser("billing-query");p.add_argument("--table",required=True);p.add_argument("--period",required=True)
    args=parser.parse_args()
    try:
        if args.cmd=="verify-composition": result=validate_profile(read(args.profile),args.library)
        elif args.cmd=="deploy-plan": result=deploy_plan(read(args.target),read(args.release),args.prior_revision,args.evidence_root)
        else: result={"argv":billing_query(args.table,args.period),"execution":"NOT_RUN"}
        print(raw(result).decode(),end="")
    except (ValueError,KeyError,OSError,TypeError) as exc:
        print(raw({"result":"FAIL","reason":str(exc),"execution":"NOT_RUN"}).decode(),end="");raise SystemExit(2)

if __name__ == "__main__": main()
