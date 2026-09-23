"""AUTHORED offline composition adapter over the existing 116-pack runtime.

Emits Google Cloud Run v1 manifests (JSON is YAML-compatible) and argv only.
No provider client, subprocess, credentials, executor, or production admission.
Cloud Run Services, WorkerPools and finite Jobs are deliberately distinct.
"""
from __future__ import annotations
import argparse
import copy
import os
import re
from pathlib import Path
from cloud_control import contained, digest, file_hash, read, raw, require, HEX, IMAGE

ENTRYPOINTS = {
    "electromobility-api": "SERVICE_WITH_LOOPS", "api": "REFERENCE_ONLY",
    "return-effect-worker": "WORKER_POOL", "return-exchange-worker": "WORKER_POOL",
    "return-accounting-worker": "WORKER_POOL", "return-fiscal-worker": "ARCA_WORKER_POOL",
    "arca-fiscal-worker": "ARCA_SIDECAR", "arca-parameter-worker": "ARCA_SIDECAR",
    "social-publishing": "SERVICE_WITH_LOOPS", "meta-lead-webhook": "SERVICE",
    "customer-survey-retention": "FINITE_JOB", "meta-lead-import": "FILE_IMPORT_CLI",
    "tiktok-lead-import": "FILE_IMPORT_CLI", "warranty-profile": "OPERATOR_CLI",
}
WORKER_ENV = {"return-effect-worker": "RETURN_EFFECT_WORKER_ID", "return-exchange-worker": "RETURN_EXCHANGE_WORKER_ID", "return-accounting-worker": "RETURN_ACCOUNTING_WORKER_ID", "return-fiscal-worker": "RETURN_FISCAL_WORKER_ID", "arca-fiscal-worker": "FISCAL_WORKER_ID", "arca-parameter-worker": "FISCAL_PARAMETER_WORKER_ID"}
REQUIRED_TOP = {"schema", "scope", "environment", "production_authorized", "target", "source", "workloads", "web", "arca", "database", "storage", "identity", "scheduler", "documents", "cost_controller", "ci"}
COMPONENTS = {"api", "web", "go-workers", "arca", "documents", "identity", "database", "migrations", "storage", "scheduler", "cost-controller", "ci"}
NAME = re.compile(r"[a-z][a-z0-9-]{2,48}[a-z0-9]")
SECRET = re.compile(r"([a-z][a-z0-9-]{2,62}):([1-9][0-9]*)")

def exact(value, keys, label):
    require(isinstance(value, dict) and set(value) == set(keys), "exact " + label + " fields required")

def source_inventory(root):
    root = Path(root)
    require(root.is_dir(), "source composition directory required")
    for p in (root, *root.parents):
        require(not p.is_symlink() and not (hasattr(p, "is_junction") and p.is_junction()), "linked source root rejected")
    observed = {p.parent.name for p in (root / "cmd").glob("*/main.go")}
    require(observed == set(ENTRYPOINTS), "entrypoint inventory changed; classify every executable before generating")
    refs = {"cmd/" + name + "/main.go": file_hash(contained(root, "cmd/" + name + "/main.go")) for name in ENTRYPOINTS}
    migrations = sorted((root / "db/migrations").glob("*.up.sql"))
    require(migrations and all(re.fullmatch(r"[0-9]{4}_.+\.up\.sql", p.name) for p in migrations), "ordered selected migration filenames required")
    ordinals = [int(p.name[:4]) for p in migrations]
    require(ordinals[0] == 1 and len(set(ordinals)) == len(ordinals), "unique selected migration ordinals required")
    # The selected baseline deliberately omits historical 0063/0064. Selection
    # position is the cursor; the greatest filename ordinal is never a count.
    for p in migrations:
        name = p.relative_to(root).as_posix(); refs[name] = file_hash(contained(root, name))
    for name in ("package.json", "go.mod", "arca/fiscal/worker/Elite.Arca.Wsfe.Worker/Program.cs", "arca/fiscal/worker/Elite.Arca.Wsfe.Worker/WorkerHost.cs"):
        refs[name] = file_hash(contained(root, name))
    return {"files": refs, "migration_paths": [p.relative_to(root).as_posix() for p in migrations], "entrypoint_kinds": ENTRYPOINTS}

def environment(values, secrets):
    require(isinstance(values, dict) and isinstance(secrets, dict) and not set(values).intersection(secrets), "environment/secrets must be disjoint maps")
    result = []
    for key, value in sorted(values.items()):
        require(re.fullmatch(r"[A-Z][A-Z0-9_]*", key) and isinstance(value, str) and value and len(value) <= 1024 and not any(x in value for x in ("\x00", "\r", "\n", "-----BEGIN")) and not re.search(r"://[^/\s]*@", value), "unsafe environment value")
        require(key not in ("DATABASE_URL", "PGPASSWORD") and not any(x in key for x in ("TOKEN", "PASSWORD", "PRIVATE_KEY", "SECRET", "HMAC_KEY")), "secret values must use references")
        result.append({"name": key, "value": value})
    for key, value in sorted(secrets.items()):
        require(re.fullmatch(r"[A-Z][A-Z0-9_]*", key) and isinstance(value, str) and SECRET.fullmatch(value), "exact Secret Manager name:version required")
        name, version = value.split(":")
        result.append({"name": key, "valueFrom": {"secretKeyRef": {"name": name, "key": version}}})
    return result

def validate_config(config, images, source_root, composition):
    exact(config, REQUIRED_TOP, "stack configuration")
    require(config["schema"] == "elite-cloud-stack-config/v403.1" and config["scope"] == "LIBRARY_INFRASTRUCTURE" and config["environment"] in ("development", "staging") and config["production_authorized"] is False, "non-production library configuration required")
    t = config["target"]
    exact(t, {"project", "project_number", "region", "platform", "network", "subnetwork"}, "target")
    require(NAME.fullmatch(t["project"]) and re.fullmatch(r"[0-9]{6,20}", t["project_number"]) and re.fullmatch(r"[a-z]+-[a-z]+[0-9]", t["region"]) and t["platform"] == "linux/amd64", "exact Linux target required")
    require(NAME.fullmatch(t["network"]) and NAME.fullmatch(t["subnetwork"]), "explicit existing private network references required")
    inv = source_inventory(source_root)
    exact(config["source"], {"composition_sha256", "inventory_sha256", "source_commit", "pack_count"}, "source")
    count = composition["selection"]["pack_count"]
    require(count in (116, 120) and config["source"]["pack_count"] == count == len(composition["pack_locks"]), "complete base116 composition required")
    if count == 120:
        from seal_consumer import EXTRAS
        require(composition.get("extension_policy") == "BASE116_PLUS_EXACT_FOUR_V403" and composition.get("base_selection", {}).get("pack_count") == 116 and {p["packId"] for p in composition["pack_locks"][116:]} == EXTRAS, "exact four extensions over base116 required")
    require(len({p["packId"] for p in composition["pack_locks"]}) == count and config["source"]["composition_sha256"] == digest(composition) and config["source"]["inventory_sha256"] == digest(inv), "source composition or executable inventory drift")
    require(re.fullmatch(r"[a-f0-9]{40}", config["source"]["source_commit"]), "pinned source commit required")
    exact(images, {"schema", "method", "artifacts"}, "images")
    require(images["schema"] == "elite-cloud-stack-images/v403.1" and images["method"] in ("SIMULATED_ARTIFACTS", "EXECUTED_TARGET"), "explicit image evidence method required")
    require(set(images["artifacts"]) == {"api", "web", "workers", "migrations", "arca", "controller"}, "all runtime image families required")
    for artifact in images["artifacts"].values():
        exact(artifact, {"image", "promoted_image", "target_sha256", "source_commit", "composition_sha256", "inventory_sha256", "entrypoints"}, "artifact")
        require(IMAGE.fullmatch(artifact["image"]) and artifact["image"].split("/")[1] == t["project"] and artifact["promoted_image"] == artifact["image"], "immutable built/promoted image identity mismatch")
        require(artifact["target_sha256"] == digest(t) and all(artifact[k] == config["source"][k] for k in ("source_commit", "composition_sha256", "inventory_sha256")), "artifact target/source identity mismatch")
        require(isinstance(artifact["entrypoints"], list) and artifact["entrypoints"] and len(set(artifact["entrypoints"])) == len(artifact["entrypoints"]), "artifact executable inventory required")
    require(set(config["workloads"]) == set(ENTRYPOINTS), "every command requires explicit deployment classification")
    names = set()
    for key, w in {**config["workloads"], "web": config["web"]}.items():
        exact(w, {"enabled", "reason", "name", "service_account", "image", "environment", "secret_refs", "min_instances", "max_instances", "arguments"}, "workload")
        require(type(w["enabled"]) is bool and isinstance(w["reason"], str) and (w["enabled"] or w["reason"].strip()), "disabled workload requires explicit reason")
        require(NAME.fullmatch(w["name"]) and w["name"] not in names, "unique resource name required"); names.add(w["name"])
        require(re.fullmatch(r"[a-z][a-z0-9-]{4,28}[a-z0-9]@" + re.escape(t["project"]) + r"\.iam\.gserviceaccount\.com", w["service_account"]), "project workload identity required")
        require(w["image"] in images["artifacts"] and type(w["min_instances"]) is int and type(w["max_instances"]) is int and 0 <= w["min_instances"] <= w["max_instances"] <= 10, "bounded workload instances required")
        environment(w["environment"], w["secret_refs"])
        require(not set(w["environment"]).intersection({"PORT", "K_SERVICE", "K_REVISION", "K_CONFIGURATION"}), "Cloud Run reserved environment key")
        kind = ENTRYPOINTS.get(key, "WEB")
        if kind in ("REFERENCE_ONLY", "OPERATOR_CLI", "FILE_IMPORT_CLI"):
            require(w["enabled"] is False, "operator/file CLI cannot become a server or unattended job")
        if w["enabled"]:
            command = "node" if key == "web" else "/app/bin/" + key
            require(command in images["artifacts"][w["image"]]["entrypoints"], "image does not contain selected executable")
            if kind in ("SERVICE_WITH_LOOPS", "WORKER_POOL", "ARCA_WORKER_POOL", "ARCA_SIDECAR"):
                require(w["min_instances"] >= 1, "background loops require persistent CPU/instances")
            if key != "web": require("DATABASE_URL" in w["secret_refs"], "database reference required")
            if key in WORKER_ENV: require(w["environment"].get(WORKER_ENV[key]) == key and w["min_instances"] == w["max_instances"] == 1, "stable worker identity requires exactly one instance")
        require(isinstance(w["arguments"], list) and all(isinstance(x, str) and x and "\x00" not in x for x in w["arguments"]), "explicit argv list required")
        if key != "customer-survey-retention": require(w["arguments"] == [], "unexpected executable arguments")
    require(config["workloads"]["electromobility-api"]["enabled"] and config["web"]["enabled"] and config["web"]["min_instances"] == 0, "API and scalable-to-zero web required")
    social = config["workloads"]["social-publishing"]
    if social["enabled"]:
        require(social["min_instances"] == social["max_instances"] == 1 and set(social["environment"]) >= {"SOCIAL_PROFILE_FILE", "SOCIAL_PROFILE_SHA256", "SOCIAL_PYTHON", "SOCIAL_BRIDGE_FILE", "SOCIAL_BRIDGE_SHA256", "META_APP_ID", "OIDC_ISSUER", "OIDC_AUDIENCE"} and set(social["secret_refs"]) >= {"SOCIAL_FENCE_HMAC_KEY_HEX", "META_APP_SECRET", "META_PAGE_ACCESS_TOKEN"}, "social publishing requires complete provider/process profile and fixed worker identity")
        require(all(HEX.fullmatch(social["environment"][k]) for k in ("SOCIAL_PROFILE_SHA256", "SOCIAL_BRIDGE_SHA256")), "social profile and bridge hashes required")
        require(social["environment"]["SOCIAL_PYTHON"] in images["artifacts"][social["image"]]["entrypoints"], "social subprocess runtime must exist in exact image")
    meta = config["workloads"]["meta-lead-webhook"]
    if meta["enabled"]:
        require(set(meta["secret_refs"]) >= {"META_APP_SECRET", "META_VERIFY_TOKEN"} and set(meta["environment"]) >= {"TENANT_ID", "ORGANIZATION_ID"}, "Meta webhook provider and tenant/org bindings required")
    require(config["web"]["max_instances"] >= 1 and config["web"]["environment"].get("OIDC_CLIENT_ID") and "OIDC_CLIENT_SECRET" in config["web"]["secret_refs"], "web OIDC client configuration required")
    require(all(re.fullmatch(r"[a-z][a-z0-9]*(?:-[a-z0-9]+)*", config["web"]["environment"].get(k, "")) for k in ("ENTERPRISE_TENANT_CODE", "ENTERPRISE_ORGANIZATION_CODE")), "explicit public catalog tenant/org configuration required")
    for key in ("return-effect-worker", "return-exchange-worker", "return-accounting-worker", "customer-survey-retention"):
        require(config["workloads"][key]["enabled"], "required durable workload omitted: " + key)
    retention = config["workloads"]["customer-survey-retention"]["arguments"]
    require(len(retention) == 6 and retention[::2] == ["--tenant", "--organization", "--limit"] and re.fullmatch(r"[0-9a-f-]{36}", retention[1]) and NAME.fullmatch(retention[3]) and retention[5].isdigit() and 1 <= int(retention[5]) <= 1000, "retention requires bounded explicit tenant/org")
    exact(config["arca"], {"enabled", "reason", "socket", "bootstrap_command", "bootstrap_sha256", "certificate_secret", "certificate_password_secret", "thumbprint", "store"}, "ARCA")
    a = config["arca"]
    require(type(a["enabled"]) is bool and (a["enabled"] or a["reason"].strip()) and a["socket"] == "/run/arca/wsfe.sock" and a["store"] == "CurrentUser", "ARCA explicit activation and UDS contract required")
    require(a["bootstrap_command"] == ["dotnet", "/app/Elite.Cloud.ArcaLauncher.dll"] and HEX.fullmatch(a["bootstrap_sha256"]) and SECRET.fullmatch(a["certificate_secret"]) and SECRET.fullmatch(a["certificate_password_secret"]) and re.fullmatch(r"[A-F0-9]{40}", a["thumbprint"]), "ARCA bootstrap/certificate references required")
    for key in ("arca-fiscal-worker", "arca-parameter-worker", "return-fiscal-worker"):
        require(config["workloads"][key]["enabled"] is a["enabled"], "ARCA processes must be absent together when deferred")
    if a["enabled"]:
        require(config["workloads"]["electromobility-api"]["max_instances"] == 1 and a["bootstrap_command"][0] in images["artifacts"]["arca"]["entrypoints"], "ARCA fixed worker identities require one API instance and bootstrap image")
        require(a["bootstrap_sha256"] == file_hash(contained(Path(__file__).parent, "cert_bridge/Launcher.cs")), "ARCA bootstrap source drift")
    d = config["database"]
    exact(d, {"instance", "database", "user", "tier", "backup_time", "retained_backups", "password_secret", "migration_baseline", "migration_last", "migration_timeout_seconds"}, "database")
    require(all(NAME.fullmatch(d[k]) for k in ("instance", "database", "user")) and re.fullmatch(r"db-custom-[1-9][0-9]*-[1-9][0-9]*", d["tier"]) and re.fullmatch(r"(?:[01][0-9]|2[0-3]):[0-5][0-9]", d["backup_time"]), "explicit managed database shape")
    require(type(d["retained_backups"]) is int and 7 <= d["retained_backups"] <= 365 and SECRET.fullmatch(d["password_secret"]), "retained database backups and exact password reference required")
    require(type(d["migration_baseline"]) is int and 0 <= d["migration_baseline"] < d["migration_last"] == len(inv["migration_paths"]) and 60 <= d["migration_timeout_seconds"] <= 3600, "verified ordered pending migration suffix required")
    require(set(config["storage"]) == {"quarantine", "retained", "controller"}, "all storage responsibilities required")
    for key, v in config["storage"].items():
        exact(v, {"bucket", "retention_days", "soft_delete_days"}, "bucket")
        require(NAME.fullmatch(v["bucket"]) and type(v["retention_days"]) is int and (v["retention_days"] == 0 if key == "controller" else 1 <= v["retention_days"] <= 36500) and 7 <= v["soft_delete_days"] <= 90, "bounded storage policy; mutable CAS controller must not have bucket retention")
    require(len({v["bucket"] for v in config["storage"].values()}) == 3, "quarantine, retained and controller buckets must differ")
    identity = config["identity"]
    exact(identity, {"issuer", "audience", "web_url", "api_url", "browser_access", "allowed_web_principals", "session_secret", "migration_service_account", "scheduler_service_account", "controller_service_account"}, "identity")
    require(identity["browser_access"] in ("IAP_DIRECT", "IAM_API_ONLY"), "explicit protected browser or API-only scope required")
    for k in ("issuer", "web_url", "api_url"):
        require(re.fullmatch(r"https://[a-zA-Z0-9.-]+(?:/[a-zA-Z0-9/_-]*)?", identity[k]), "HTTPS identity and service URL required")
    require(identity["api_url"] != identity["web_url"] and SECRET.fullmatch(identity["session_secret"]), "separate API/web URL and session reference required")
    for k, name in (("web_url", config["web"]["name"]), ("api_url", config["workloads"]["electromobility-api"]["name"])):
        require(identity[k] == "https://" + name + "-" + t["project_number"] + "." + t["region"] + ".run.app", "service URL must match exact target resource")
    require(identity["allowed_web_principals"] and all(re.fullmatch(r"(?:user|group):[A-Za-z0-9._+-]+@[A-Za-z0-9.-]+", x) for x in identity["allowed_web_principals"]), "explicit web principals; public invoker forbidden")
    for k in ("migration_service_account", "scheduler_service_account", "controller_service_account"):
        require(re.fullmatch(r"[a-z][a-z0-9-]{4,28}[a-z0-9]@" + re.escape(t["project"]) + r"\.iam\.gserviceaccount\.com", identity[k]), "separate bounded workload service accounts required")
    exact(config["scheduler"], {"retention_schedule", "timezone", "retry_count", "max_retry_seconds"}, "scheduler")
    require(re.fullmatch(r"[0-9*/,-]+(?: [0-9*/,-]+){4}", config["scheduler"]["retention_schedule"]) and config["scheduler"]["timezone"] == "Etc/UTC" and config["scheduler"]["retry_count"] == 0 and config["scheduler"]["max_retry_seconds"] == 0, "explicit bounded scheduler; no automatic retention retry")
    exact(config["documents"], {"mode", "automatic_persistence", "profile_ref", "profile_sha256"}, "documents")
    require(config["documents"]["mode"] == "REVIEW_ONLY" and config["documents"]["automatic_persistence"] is False and HEX.fullmatch(config["documents"]["profile_sha256"]), "document evidence remains review-only")
    fc = config["cost_controller"]
    exact(fc, {"state_object", "billing_table", "budget_policy_ref", "budget_policy_sha256", "job_name", "schedule", "timeout_seconds", "cycle", "cycle_owner_sha256", "mounts", "seed_receipt_ref", "seed_receipt_sha256"}, "cost controller")
    from finops_cycle import validate_config as validate_cycle
    validate_cycle(fc["cycle"])
    require(fc["state_object"] == "control/state.json" and fc["billing_table"] == fc["cycle"]["billing_table"] and HEX.fullmatch(fc["budget_policy_sha256"]), "bound durable billing configuration required")
    require(fc["cycle"]["project"] == t["project"] and fc["cycle"]["environment"] == config["environment"], "FinOps target mismatch")
    bucket = "gs://" + config["storage"]["controller"]["bucket"]
    require(fc["cycle"]["journal_uri"] == bucket + "/control/billing.sqlite" and fc["cycle"]["fence_uri"] == bucket + "/control/finops-intents", "FinOps journal/fence namespace must bind controller storage")
    require(NAME.fullmatch(fc["job_name"]) and fc["job_name"] not in names and fc["job_name"] != d["instance"] + "-migrations", "unique FinOps Job name required")
    other_accounts = {w["service_account"] for w in [*config["workloads"].values(), config["web"]]} | {identity["migration_service_account"], identity["scheduler_service_account"], config["ci"]["deploy_service_account"]}
    require(identity["controller_service_account"] not in other_accounts, "FinOps must use a separate least-privilege workload identity")
    interval = fc["cycle"]["interval_seconds"]
    require(interval % 60 == 0 and 1 <= interval // 60 <= 60 and 60 % (interval // 60) == 0, "FinOps scheduler requires a whole-minute divisor of one hour")
    require(fc["schedule"] == ("0 * * * *" if interval == 3600 else "*/" + str(interval // 60) + " * * * *"), "FinOps schedule must equal approved interval")
    require(type(fc["timeout_seconds"]) is int and 60 <= fc["timeout_seconds"] <= interval, "bounded FinOps task timeout must not exceed interval")
    require(fc["cycle_owner_sha256"] == file_hash(Path(__file__).with_name("finops_cycle.py")), "FinOps cycle owner drift")
    require({"python3", "/app/cloud/finops_cycle.py"} <= set(images["artifacts"]["controller"]["entrypoints"]), "controller image executable inventory incomplete")
    exact(fc["mounts"], {"config", "approval", "tools", "runner"}, "FinOps mounted references")
    for key, reference in fc["mounts"].items():
        exact(reference, {"secret_ref", "sha256"}, "FinOps mounted file")
        require(SECRET.fullmatch(reference["secret_ref"]) and HEX.fullmatch(reference["sha256"]), "versioned hash-bound FinOps mount required")
    require(fc["mounts"]["config"]["sha256"] == digest(fc["cycle"]), "mounted FinOps configuration differs from generated cycle")
    require(len({r["secret_ref"].split(":")[0] for r in fc["mounts"].values()}) == 4, "separate FinOps config/approval/tools/runner secrets required")
    require(re.fullmatch(r"target-evidence/[a-z0-9-]+\.json", fc["seed_receipt_ref"]) and HEX.fullmatch(fc["seed_receipt_sha256"]), "hash-bound seed receipt reference required")
    exact(config["ci"], {"library_repository", "library_commit", "project_repository", "wif_principal", "deploy_service_account"}, "CI")
    require(config["ci"]["library_repository"] == "https://github.com/etor6233/elite-engineering-library.git" and re.fullmatch(r"[a-f0-9]{40}", config["ci"]["library_commit"]) and re.fullmatch(r"https://github.com/[A-Za-z0-9_-]+/[A-Za-z0-9_.-]+", config["ci"]["project_repository"]), "two distinct pinned repositories required")
    require(config["ci"]["library_repository"] != config["ci"]["project_repository"] and re.fullmatch(r"principalSet://iam.googleapis.com/projects/[0-9]+/locations/global/workloadIdentityPools/[a-z0-9-]+/attribute.repository/[A-Za-z0-9_-]+/[A-Za-z0-9_.-]+", config["ci"]["wif_principal"]), "explicit repository-scoped WIF; no exported service account key")
    require(config["ci"]["wif_principal"].startswith("principalSet://iam.googleapis.com/projects/" + t["project_number"] + "/") and config["ci"]["wif_principal"].endswith("/attribute.repository/" + config["ci"]["project_repository"].removeprefix("https://github.com/")), "WIF repository or target project mismatch")
    require(re.fullmatch(r"[a-z][a-z0-9-]{4,28}[a-z0-9]@" + re.escape(t["project"]) + r"\.iam\.gserviceaccount\.com", config["ci"]["deploy_service_account"]), "bounded deployment identity required")
    return inv

def build_stack(config, images, *, source_root=None, composition=None):
    project = Path(__file__).resolve().parents[1]
    source_root = Path(source_root) if source_root is not None else Path(os.environ["ELITE_CLOUD_SOURCE_ROOT"]) if "ELITE_CLOUD_SOURCE_ROOT" in os.environ else project if (project / "cmd").is_dir() else project / "reference-composition"
    composition = composition if composition is not None else read(Path(__file__).with_name("composition.json"))
    inv = validate_config(config, images, source_root, composition)
    c = copy.deepcopy(config); t = c["target"]; d = c["database"]; identity = c["identity"]
    manifests, commands = {}, []
    common = ["--project", t["project"], "--region", t["region"], "--quiet"]
    conn = t["project"] + ":" + t["region"] + ":" + d["instance"]
    def step(key, argv, needs=(), gates=()):
        commands.append({"id": key, "effect": "WRITE", "argv": argv, "depends_on": list(needs), "required_target_receipts": list(gates), "execution": "NOT_RUN"})
    def add_manifest(key, kind, name, account, containers, minimum=1, maximum=1, volumes=None):
        annotations = {"run.googleapis.com/cloudsql-instances": conn, "run.googleapis.com/network-interfaces": __import__("json").dumps([{"network": t["network"], "subnetwork": t["subnetwork"]}], separators=(",", ":")), "run.googleapis.com/vpc-access-egress": "private-ranges-only"}
        spec = {"serviceAccountName": account, "containers": containers}
        if volumes: spec["volumes"] = volumes
        meta = {"name": name, "namespace": t["project_number"], "labels": {"cloud.googleapis.com/location": t["region"]}}
        if kind == "Service":
            meta["annotations"] = {"run.googleapis.com/ingress": "all", "run.googleapis.com/invoker-iam-disabled": "false"}
            if key == "web" and identity["browser_access"] == "IAP_DIRECT": meta["annotations"]["run.googleapis.com/iap-enabled"] = "true"
            annotations.update({"autoscaling.knative.dev/minScale": str(minimum), "autoscaling.knative.dev/maxScale": str(maximum), "run.googleapis.com/cpu-throttling": "false" if minimum else "true"})
            spec.update(containerConcurrency=20, timeoutSeconds=60)
            body = {"template": {"metadata": {"annotations": annotations}, "spec": spec}}
            api_version, resource = "serving.knative.dev/v1", "services"
        elif kind == "WorkerPool":
            meta["annotations"] = {"run.googleapis.com/manualInstanceCount": str(minimum)}
            body = {"template": {"metadata": {"annotations": annotations}, "spec": spec}}
            api_version, resource = "run.googleapis.com/v1", "worker-pools"
        else:
            spec.update(maxRetries=0, timeoutSeconds=d["migration_timeout_seconds"] if key == "migrations" else 60)
            body = {"template": {"metadata": {"annotations": annotations}, "spec": {"parallelism": 1, "taskCount": 1, "template": {"spec": spec}}}}
            api_version, resource = "run.googleapis.com/v1", "jobs"
        path = "manifests/" + key + ".yaml"
        manifests[path] = {"apiVersion": api_version, "kind": kind, "metadata": meta, "spec": body}
        step("configure-" + key, ["gcloud", "run", resource, "replace", path, *common], ("identity", "database"), ("artifact_same_digest", "network_access", "migration_baseline"))
    def container(key, w, ports=False):
        out = {"name": key, "image": images["artifacts"][w["image"]]["image"], "command": ["node"] if key == "web" else ["/app/bin/" + key], "args": ["/app/server.js"] if key == "web" else w["arguments"], "env": environment(w["environment"], w["secret_refs"]), "resources": {"limits": {"cpu": "1", "memory": "512Mi"}}}
        if ports: out["ports"] = [{"containerPort": 8080}]
        return out
    step("database", ["gcloud", "sql", "instances", "create", d["instance"], "--project", t["project"], "--region", t["region"], "--database-version=POSTGRES_18", "--tier", d["tier"], "--network=projects/" + t["project"] + "/global/networks/" + t["network"], "--no-assign-ip", "--enable-point-in-time-recovery", "--backup-start-time=" + d["backup_time"], "--retained-backups-count=" + str(d["retained_backups"]), "--deletion-protection", "--retain-backups-on-delete", "--ssl-mode=ENCRYPTED_ONLY", "--quiet"], gates=("budget", "private_service_access", "new_resource_absent", "retention_policy"))
    step("database-name", ["gcloud", "sql", "databases", "create", d["database"], "--instance", d["instance"], "--project", t["project"], "--quiet"], ("database",), ("database_user_and_secret_provisioned",))
    accounts = sorted({w["service_account"] for key, w in {**c["workloads"], "web": c["web"]}.items() if w["enabled"] and ENTRYPOINTS.get(key) != "ARCA_SIDECAR"} | {identity[k] for k in ("migration_service_account", "scheduler_service_account", "controller_service_account")})
    sql_accounts = {w["service_account"] for key, w in c["workloads"].items() if w["enabled"] and ENTRYPOINTS[key] != "ARCA_SIDECAR"} | {identity["migration_service_account"]}
    identity_dependencies = []
    for index, account in enumerate(accounts):
        step("identity-" + str(index), ["gcloud", "iam", "service-accounts", "create", account.split("@")[0], "--project", t["project"], "--quiet"], gates=("new_resource_absent", "least_privilege"))
        identity_dependencies.append("identity-" + str(index))
        if account in sql_accounts:
            step("sql-client-" + str(index), ["gcloud", "projects", "add-iam-policy-binding", t["project"], "--member=serviceAccount:" + account, "--role=roles/cloudsql.client", "--condition=None", "--quiet"], ("identity-" + str(index),))
            identity_dependencies.append("sql-client-" + str(index))
    commands.append({"id": "identity", "effect": "BARRIER", "argv": [], "depends_on": identity_dependencies, "required_target_receipts": ["least_privilege", "secret_accessor_bindings", "oidc_discovery", "wif_repository_condition"], "execution": "NOT_RUN"})
    for key, b in c["storage"].items():
        retention_flags = ["--retention-period=P" + str(b["retention_days"]) + "D"] if b["retention_days"] else []
        step("storage-" + key, ["gcloud", "storage", "buckets", "create", "gs://" + b["bucket"], "--project", t["project"], "--location", t["region"], "--uniform-bucket-level-access", "--public-access-prevention", *retention_flags, "--soft-delete-duration=" + str(b["soft_delete_days"]) + "d", "--quiet"], gates=("new_resource_absent", "retention_policy", "storage_iam"))
        step("storage-versioning-" + key, ["gcloud", "storage", "buckets", "update", "gs://" + b["bucket"], "--project", t["project"], "--versioning", "--quiet"], ("storage-" + key,), ("version_history_lifecycle_budget",))
    migration_files = inv["migration_paths"][d["migration_baseline"]:]
    migration = {"name": "migrations", "image": images["artifacts"]["migrations"]["image"], "command": ["psql"], "args": ["-X", "--set=ON_ERROR_STOP=1", *[arg for p in migration_files for arg in ("--file", "/migrations/" + Path(p).name)]], "env": environment({"PGHOST": "/cloudsql/" + conn, "PGDATABASE": d["database"], "PGUSER": d["user"]}, {"PGPASSWORD": d["password_secret"]})}
    require("psql" in images["artifacts"]["migrations"]["entrypoints"], "migration image must contain exact PostgreSQL client and locked SQL files")
    add_manifest("migrations", "Job", d["instance"] + "-migrations", identity["migration_service_account"], [migration])
    step("run-migrations", ["gcloud", "run", "jobs", "execute", d["instance"] + "-migrations", *common, "--wait"], ("configure-migrations", "database-name"), ("backup_restore", "empty_database_or_verified_pending_suffix", "no_concurrent_migrator"))
    for key, w in c["workloads"].items():
        if not w["enabled"] or ENTRYPOINTS[key] == "ARCA_SIDECAR": continue
        kind = ENTRYPOINTS[key]
        if key == "electromobility-api":
            w["environment"].update(HTTP_ADDRESS=":8080", OIDC_ISSUER=identity["issuer"], OIDC_AUDIENCE=identity["audience"], ARCA_ENABLED="true" if c["arca"]["enabled"] else "false")
            if c["arca"]["enabled"]: w["environment"]["ARCA_WSFE_SOCKET"] = c["arca"]["socket"]
        if key == "social-publishing": w["environment"]["SOCIAL_HTTP_ADDRESS"] = ":8080"
        if key == "meta-lead-webhook": w["environment"]["LISTEN_ADDR"] = ":8080"
        containers, volumes = [container(key, w, kind.startswith("SERVICE"))], []
        if key == "electromobility-api" and c["arca"]["enabled"]:
            a = c["arca"]; secret_name, secret_version = a["certificate_secret"].split(":"); password_name, password_version = a["certificate_password_secret"].split(":")
            volumes = [{"name": "arca-uds", "emptyDir": {"medium": "Memory", "sizeLimit": "16Mi"}}, {"name": "arca-certificate", "secret": {"secretName": secret_name, "items": [{"key": secret_version, "path": "cert.pfx"}]}}, {"name": "arca-password", "secret": {"secretName": password_name, "items": [{"key": password_version, "path": "password"}]}}]
            containers[0]["volumeMounts"] = [{"name": "arca-uds", "mountPath": "/run/arca"}]
            for worker in ("arca-fiscal-worker", "arca-parameter-worker"):
                row = copy.deepcopy(c["workloads"][worker]); row["environment"]["ARCA_WSFE_SOCKET"] = a["socket"]
                item = container(worker, row); item["volumeMounts"] = [{"name": "arca-uds", "mountPath": "/run/arca"}]; containers.append(item)
            dotnet_env = environment({"HOME": "/app/private-home", "ARCA_ENVIRONMENT": "homologation", "ARCA_WSFE_SOCKET": a["socket"], "ARCA_CERTIFICATE_STORE": "CurrentUser", "ARCA_CERTIFICATE_THUMBPRINT": a["thumbprint"], "ARCA_PFX_FILE": "/var/run/secrets/arca/certificate/cert.pfx"}, {})
            dotnet_env.append({"name": "ARCA_PFX_PASSWORD_FILE", "value": "/var/run/secrets/arca/password/password"})
            containers.append({"name": "arca-dotnet", "image": images["artifacts"]["arca"]["image"], "command": a["bootstrap_command"], "env": dotnet_env, "volumeMounts": [{"name": "arca-uds", "mountPath": "/run/arca"}, {"name": "arca-certificate", "mountPath": "/var/run/secrets/arca/certificate"}, {"name": "arca-password", "mountPath": "/var/run/secrets/arca/password"}], "resources": {"limits": {"cpu": "1", "memory": "512Mi"}}})
        resource_kind = "Service" if kind.startswith("SERVICE") else "Job" if kind == "FINITE_JOB" else "WorkerPool"
        add_manifest(key, resource_kind, w["name"], w["service_account"], containers, w["min_instances"], w["max_instances"], volumes)
        commands[-1]["depends_on"].append("run-migrations")
        if key in ("social-publishing", "meta-lead-webhook"):
            commands[-1]["required_target_receipts"].extend(["provider_profile_and_exact_process_files", "provider_ingress_identity_transport"])
    web = c["web"]; web["environment"].update(HOSTNAME="0.0.0.0", APP_BASE_URL=identity["web_url"], ENTERPRISE_API_BASE_URL=identity["api_url"], ELITE_CLOUD_RUN_IDENTITY="metadata", ELITE_CLOUD_RUN_AUDIENCE=identity["api_url"], OIDC_ISSUER=identity["issuer"], OIDC_AUDIENCE=identity["audience"]); web["secret_refs"]["AUTH_SESSION_SECRET"] = identity["session_secret"]
    if identity["browser_access"] == "IAP_DIRECT":
        step("iap-api", ["gcloud", "services", "enable", "iap.googleapis.com", "--project", t["project"], "--quiet"], gates=("budget", "iap_oauth_setup_allowed_users_and_app_oidc_callback"))
    add_manifest("web", "Service", web["name"], web["service_account"], [container("web", web, True)], 0, web["max_instances"])
    commands[-1]["depends_on"].append("configure-electromobility-api")
    if identity["browser_access"] == "IAP_DIRECT":
        commands[-1]["depends_on"].append("iap-api")
        commands[-1]["required_target_receipts"].extend(["iap_oauth_setup_allowed_users_and_app_oidc_callback", "iap_service_agent_exists"])
        step("iap-invoker", ["gcloud", "run", "services", "add-iam-policy-binding", web["name"], *common, "--member=serviceAccount:service-" + t["project_number"] + "@gcp-sa-iap.iam.gserviceaccount.com", "--role=roles/run.invoker"], ("configure-web",), ("iap_service_agent_exists",))
    for index, principal in enumerate(identity["allowed_web_principals"]):
        if identity["browser_access"] == "IAP_DIRECT":
            step("web-invoker-" + str(index), ["gcloud", "iap", "web", "add-iam-policy-binding", "--resource-type=cloud-run", "--service", web["name"], *common, "--member", principal, "--role=roles/iap.httpsResourceAccessor"], ("iap-invoker",), ("allowed_user_identity", "iap_oauth_setup_allowed_users_and_app_oidc_callback"))
        else:
            step("web-invoker-" + str(index), ["gcloud", "run", "services", "add-iam-policy-binding", web["name"], *common, "--member", principal, "--role=roles/run.invoker"], ("configure-web",), ("allowed_user_identity",))
    step("api-invoker", ["gcloud", "run", "services", "add-iam-policy-binding", c["workloads"]["electromobility-api"]["name"], *common, "--member=serviceAccount:" + web["service_account"], "--role=roles/run.invoker"], ("configure-electromobility-api",), ("dual_authorization_headers",))
    job = c["workloads"]["customer-survey-retention"]["name"]
    step("retention-invoker", ["gcloud", "run", "jobs", "add-iam-policy-binding", job, *common, "--member=serviceAccount:" + identity["scheduler_service_account"], "--role=roles/run.invoker"], ("configure-customer-survey-retention",))
    step("scheduler", ["gcloud", "scheduler", "jobs", "create", "http", job + "-schedule", "--project", t["project"], "--location", t["region"], "--schedule", c["scheduler"]["retention_schedule"], "--time-zone=Etc/UTC", "--uri=https://run.googleapis.com/v2/projects/" + t["project"] + "/locations/" + t["region"] + "/jobs/" + job + ":run", "--http-method=POST", "--oauth-service-account-email=" + identity["scheduler_service_account"], "--max-retry-attempts=0", "--max-retry-duration=0s", "--quiet"], ("retention-invoker",), ("retention_scope_and_schedule", "pause_resume_reconciliation"))
    # FinOps is a finite Python Job with no SQL connection, public endpoint or
    # cloud-management role. Its immutable config is mounted, never inline secrets.
    fc = c["cost_controller"]; cycle = fc["cycle"]; controller = identity["controller_service_account"]
    controller_identity = "identity-" + str(accounts.index(controller))
    mounts, volumes, args = [], [], ["/app/cloud/finops_cycle.py", "run"]
    for key, reference in sorted(fc["mounts"].items()):
        secret, version = reference["secret_ref"].split(":")
        directory = "/var/run/elite/finops/" + key
        mounts.append({"name": "finops-" + key, "mountPath": directory})
        volumes.append({"name": "finops-" + key, "secret": {"secretName": secret, "items": [{"key": version, "path": key + ".json"}], "defaultMode": 292}})
        args += ["--" + key, directory + "/" + key + ".json"]
    args += ["--workdir", "/tmp/elite-finops-cycle", "--authorize-execution"]
    fc_container = {"name": "finops-controller", "image": images["artifacts"]["controller"]["image"], "command": ["python3"], "args": args, "volumeMounts": mounts, "resources": {"limits": {"cpu": "1", "memory": "512Mi"}}}
    add_manifest("finops-controller", "Job", fc["job_name"], controller, [fc_container], volumes=volumes)
    template = manifests["manifests/finops-controller.yaml"]["spec"]["template"]
    template["metadata"] = {}; template["spec"]["template"]["spec"]["timeoutSeconds"] = fc["timeout_seconds"]
    configure_finops = next(x for x in commands if x["id"] == "configure-finops-controller")
    configure_finops["depends_on"] = ["identity", "finops-query-user", "finops-billing-reader", "finops-journal-access", "finops-intent-create", "finops-seed"]
    configure_finops["required_target_receipts"] = ["artifact_same_digest", "finops_runtime_and_mount_hashes"]
    step("finops-query-user", ["gcloud", "projects", "add-iam-policy-binding", t["project"], "--member=serviceAccount:" + controller, "--role=roles/bigquery.jobUser", "--condition=None", "--quiet"], (controller_identity,), ("finops_billing_scope_and_query_budget",))
    bp, dataset, table = cycle["billing_table"].split(".")
    expression = "resource.service == 'bigquery.googleapis.com' && resource.type == 'bigquery.googleapis.com/Table' && resource.name == 'projects/" + bp + "/datasets/" + dataset + "/tables/" + table + "'"
    step("finops-billing-reader", ["gcloud", "projects", "add-iam-policy-binding", bp, "--member=serviceAccount:" + controller, "--role=roles/bigquery.dataViewer", "--condition=title=finops_exact_billing_table,expression=" + expression, "--quiet"], (controller_identity,), ("finops_billing_scope_and_query_budget",))
    bucket = c["storage"]["controller"]["bucket"]
    step("finops-journal-role", ["gcloud", "iam", "roles", "create", "eliteFinopsJournal", "--project", t["project"], "--title=FinOps single journal access", "--permissions=storage.objects.get,storage.objects.create,storage.objects.delete", "--stage=GA", "--quiet"], gates=("new_resource_absent", "finops_namespace_iam"))
    expression = "resource.type == 'storage.googleapis.com/Object' && resource.name == 'projects/_/buckets/" + bucket + "/objects/control/billing.sqlite'"
    step("finops-journal-access", ["gcloud", "storage", "buckets", "add-iam-policy-binding", "gs://" + bucket, "--member=serviceAccount:" + controller, "--role=projects/" + t["project"] + "/roles/eliteFinopsJournal", "--condition=title=finops_exact_journal,expression=" + expression, "--quiet"], (controller_identity, "storage-versioning-controller", "finops-journal-role"), ("finops_namespace_iam",))
    expression = "resource.type == 'storage.googleapis.com/Object' && resource.name.startsWith('projects/_/buckets/" + bucket + "/objects/control/finops-intents/')"
    step("finops-intent-create", ["gcloud", "storage", "buckets", "add-iam-policy-binding", "gs://" + bucket, "--member=serviceAccount:" + controller, "--role=roles/storage.objectCreator", "--condition=title=finops_create_only_intents,expression=" + expression, "--quiet"], (controller_identity, "storage-versioning-controller"), ("finops_namespace_iam",))
    commands.append({"id": "finops-seed", "effect": "BARRIER", "argv": [], "depends_on": ["storage-versioning-controller"], "required_target_receipts": ["finops_seed_generation_zero_and_restore"], "execution": "NOT_RUN"})
    if cycle["notification"]["mode"] == "PUBSUB":
        step("finops-pubsub-publisher", ["gcloud", "pubsub", "topics", "add-iam-policy-binding", cycle["notification"]["topic"], "--project", t["project"], "--member=serviceAccount:" + controller, "--role=roles/pubsub.publisher", "--quiet"], (controller_identity,), ("finops_pubsub_topic_and_dispatch_approval",))
        configure_finops["depends_on"].append("finops-pubsub-publisher")
    step("finops-invoker", ["gcloud", "run", "jobs", "add-iam-policy-binding", fc["job_name"], *common, "--member=serviceAccount:" + identity["scheduler_service_account"], "--role=roles/run.invoker"], ("configure-finops-controller",))
    step("finops-scheduler", ["gcloud", "scheduler", "jobs", "create", "http", fc["job_name"] + "-schedule", "--project", t["project"], "--location", t["region"], "--schedule", fc["schedule"], "--time-zone=Etc/UTC", "--uri=https://run.googleapis.com/v2/projects/" + t["project"] + "/locations/" + t["region"] + "/jobs/" + fc["job_name"] + ":run", "--http-method=POST", "--oauth-service-account-email=" + identity["scheduler_service_account"], "--max-retry-attempts=0", "--max-retry-duration=0s", "--quiet"], ("finops-invoker", "finops-seed"), ("finops_schedule_approval_and_reconciliation",))
    fc["seed_initialization"] = {"execution": "NOT_RUN", "local_argv": ["python3", "/app/cloud/finops_cycle.py", "init", "--output", "seed/finops.sqlite"], "upload_argv": ["gcloud", "storage", "cp", "seed/finops.sqlite", cycle["journal_uri"], "--if-generation-match=0", "--project", t["project"], "--quiet"], "receipt_ref": fc["seed_receipt_ref"], "receipt_sha256": fc["seed_receipt_sha256"]}
    bindings = []
    for key, w in {**c["workloads"], "web": web}.items():
        if w["enabled"]:
            account = c["workloads"]["electromobility-api"]["service_account"] if ENTRYPOINTS.get(key) == "ARCA_SIDECAR" else w["service_account"]
            bindings += [(account, ref.split(":")[0]) for ref in w["secret_refs"].values()]
    bindings.append((identity["migration_service_account"], d["password_secret"].split(":")[0]))
    bindings += [(controller, reference["secret_ref"].split(":")[0]) for reference in fc["mounts"].values()]
    if c["arca"]["enabled"]: bindings += [(c["workloads"]["electromobility-api"]["service_account"], c["arca"][key].split(":")[0]) for key in ("certificate_secret", "certificate_password_secret")]
    for index, (account, name) in enumerate(sorted(set(bindings))):
        step("secret-binding-" + str(index), ["gcloud", "secrets", "add-iam-policy-binding", name, "--project", t["project"], "--member=serviceAccount:" + account, "--role=roles/secretmanager.secretAccessor", "--quiet"], ("identity-" + str(accounts.index(account)),), gates=("secret_versions_provisioned", "least_privilege"))
        next(x for x in commands if x["id"] == "identity")["depends_on"].append("secret-binding-" + str(index))
    configs = {"identity": identity, "database": d, "migrations": {"paths": migration_files, "hashes": {p: inv["files"][p] for p in migration_files}, "baseline": d["migration_baseline"], "last": d["migration_last"]}, "storage": c["storage"], "arca": c["arca"], "documents": c["documents"], "scheduler": c["scheduler"], "cost-controller": {**c["cost_controller"], "state_uri": "gs://" + c["storage"]["controller"]["bucket"] + "/" + c["cost_controller"]["state_object"], "owner": "cloud_control.BillingJournal + state_cas_plan", "hard_cap": False}, "ci": c["ci"], "api": c["workloads"]["electromobility-api"], "web": web, "go-workers": {k: v for k, v in c["workloads"].items() if "WORKER" in ENTRYPOINTS[k]}}
    plan = {"schema": "elite-cloud-stack-plan/v403.1", "scope": "LIBRARY_INFRASTRUCTURE", "environment": c["environment"], "execution": "NOT_RUN", "production_authorized": False, "target_sha256": digest(t), "configuration_sha256": digest(config), "images_sha256": digest(images), "source": c["source"], "image_evidence_method": images["method"], "components": {k: {"status": "REQUIRED", "config_ref": "configs/" + k + ".json"} for k in sorted(COMPONENTS)}, "configs": configs, "manifests": manifests, "commands": commands, "entrypoint_inventory": inv, "required_target_gates": ["exact_images_build_sca_licenses_and_runtime", "existing_network_private_service_access", "postgres18_pitr_restore_migration_baseline", "workload_identity_and_database_user_secret_values", "protected_url_dual_authorization_transport", "SIGTERM_10s_forced_stop_and_durable_recovery", "arca_pfx_bootstrap_uid_gid_socket_readiness_and_crash_recovery" if c["arca"]["enabled"] else "arca_all_processes_absent", "document_profile_corpus_review_acceptance", "billing_export_freshness_and_generation_cas", "whole_stack_provider_execution_and_semantic_step_receipts"], "limits": {"provider_executed": False, "config_is_runtime_proof": False, "fixture_artifacts_are_admitted": False, "operator_clis_are_servers": False}, "operator_access_argv": ["gcloud", "run", "services", "proxy", web["name"], *common]}
    plan["configs"]["cost-controller"]["owner"] = "finops_cycle.cycle + cloud_control.BillingJournal + execute_step.load_runner"
    plan["required_target_gates"] += ["finops_runtime_and_mount_hashes", "finops_billing_scope_and_query_budget", "finops_seed_generation_zero_and_restore", "finops_namespace_iam", "finops_schedule_approval_and_reconciliation"]
    if cycle["notification"]["mode"] == "PUBSUB": plan["required_target_gates"].append("finops_pubsub_topic_and_dispatch_approval")
    plan["browser_journey_status"] = "NOT_RUN" if identity["browser_access"] == "IAP_DIRECT" else "INFRASTRUCTURE_INCOMPLETE_PROTECTED_API_ONLY"
    if identity["browser_access"] == "IAP_DIRECT":
        plan["operator_access_argv"] = []; plan["operator_access_url"] = identity["web_url"]
        plan["required_target_gates"].append("iap_oauth_setup_allowed_users_and_app_oidc_callback")
    validate_stack(plan)
    return plan

def validate_stack(plan):
    require(plan["schema"] == "elite-cloud-stack-plan/v403.1" and plan["execution"] == "NOT_RUN" and plan["production_authorized"] is False and plan["environment"] in ("development", "staging"), "offline non-production stack required")
    require(set(plan["components"]) == set(plan["configs"]) == COMPONENTS and all(x["status"] == "REQUIRED" and x["config_ref"] == "configs/" + k + ".json" for k, x in plan["components"].items()), "required component omitted")
    known = {x["id"] for x in plan["commands"]}; require(len(known) == len(plan["commands"]), "duplicate step")
    graph = {x["id"]: set(x["depends_on"]) for x in plan["commands"]}
    require(all(v <= known for v in graph.values()), "missing dependency")
    while graph:
        ready = {k for k, v in graph.items() if not v}; require(ready, "dependency cycle")
        graph = {k: v - ready for k, v in graph.items() if k not in ready}
    for item in plan["commands"]:
        require(item["execution"] == "NOT_RUN" and item["effect"] in ("WRITE", "BARRIER") and (not item["argv"] if item["effect"] == "BARRIER" else item["argv"][0] == "gcloud"), "invalid offline step")
        require(not any(v in item["argv"] for v in ("--allow-unauthenticated", "--no-deletion-protection", "allUsers", "allAuthenticatedUsers")), "unsafe public or destructive plan")
    return {"result": "PASS", "claim": "offline stack configuration integrity only", "execution": "NOT_RUN", "production_authorized": False, "components": len(COMPONENTS), "manifests": len(plan["manifests"]), "plan_sha256": digest(plan)}

def emit(config_path, output, source_root=None, composition_path=None):
    bundle = read(config_path); exact(bundle, {"config", "images"}, "input bundle")
    plan = build_stack(bundle["config"], bundle["images"], source_root=source_root, composition=read(composition_path) if composition_path else None)
    target = Path(output).absolute()
    for p in (target, *target.parents): require(not p.is_symlink() and not (hasattr(p, "is_junction") and p.is_junction()), "linked output rejected")
    require(not target.exists(), "output must be absent; preserve earlier generation")
    target.mkdir(parents=True)
    for name, value in [("stack-plan.json", plan), ("stack-validation.json", validate_stack(plan)), *plan["manifests"].items(), *(("configs/" + k + ".json", v) for k, v in plan["configs"].items())]:
        p = contained(target, name); p.parent.mkdir(parents=True, exist_ok=True); p.write_bytes(raw(value))
    return validate_stack(plan)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(); sub = parser.add_subparsers(dest="command", required=True)
    p = sub.add_parser("emit"); p.add_argument("--config", required=True); p.add_argument("--output", required=True); p.add_argument("--source-root"); p.add_argument("--composition")
    p = sub.add_parser("validate"); p.add_argument("--plan", required=True)
    args = parser.parse_args()
    print(raw(emit(args.config, args.output, args.source_root, args.composition) if args.command == "emit" else validate_stack(read(args.plan))).decode(), end="")
