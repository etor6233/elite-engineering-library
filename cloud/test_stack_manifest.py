"""AUTHORED local negative fixtures. No provider, network, or credentials."""
import copy
import os
from pathlib import Path
import tempfile
import unittest
import stack_manifest as stack
from cloud_control import digest, file_hash, raw, read

ROOT = Path(__file__).resolve().parents[1]
SOURCE = Path(os.environ["ELITE_CLOUD_SOURCE_ROOT"]) if "ELITE_CLOUD_SOURCE_ROOT" in os.environ else ROOT if (ROOT / "cmd").is_dir() else ROOT / "reference-composition"

def fixture_bundle():
    """Explicit synthetic project/artifacts; hashes bind the selected local sources."""
    composition = read(ROOT / "cloud/composition.json")
    inv = stack.source_inventory(SOURCE)
    target = dict(project="fixture-project", project_number="123456789012", region="us-central1", platform="linux/amd64", network="fixture-network", subnetwork="fixture-subnet")
    source = dict(composition_sha256=digest(composition), inventory_sha256=digest(inv), source_commit="b" * 40, pack_count=composition["selection"]["pack_count"])
    account = lambda name: "fixture-" + name + "@fixture-project.iam.gserviceaccount.com"
    enabled = {"electromobility-api", "return-effect-worker", "return-exchange-worker", "return-accounting-worker", "customer-survey-retention"}
    workloads = {}
    for key, kind in stack.ENTRYPOINTS.items():
        workloads[key] = dict(enabled=key in enabled, reason="Required local composition" if key in enabled else "Deferred: provider credential/profile or explicit operator input remains absent", name="fixture-" + ("api" if key == "electromobility-api" else "reference-api" if key == "api" else key), service_account=account("api" if key == "electromobility-api" else "worker"), image="api" if key == "electromobility-api" else "workers", environment={stack.WORKER_ENV[key]: key} if key in stack.WORKER_ENV else {}, secret_refs={"DATABASE_URL": "fixture-database:1"}, min_instances=1, max_instances=1, arguments=[])
    workloads["customer-survey-retention"]["arguments"] = ["--tenant", "11111111-1111-4111-8111-111111111111", "--organization", "fixture-store", "--limit", "100"]
    web = dict(enabled=True, reason="Required web", name="fixture-web", service_account=account("web"), image="web", environment={"OIDC_CLIENT_ID": "fixture-web", "ENTERPRISE_TENANT_CODE": "fixture-tenant", "ENTERPRISE_ORGANIZATION_CODE": "fixture-store"}, secret_refs={"OIDC_CLIENT_SECRET": "fixture-oidc-client:1"}, min_instances=0, max_instances=3, arguments=[])
    config = dict(schema="elite-cloud-stack-config/v403.1", scope="LIBRARY_INFRASTRUCTURE", environment="development", production_authorized=False, target=target, source=source, workloads=workloads, web=web,
        arca=dict(enabled=False, reason="No fiscal account/certificate or target acceptance", socket="/run/arca/wsfe.sock", bootstrap_command=["dotnet", "/app/Elite.Cloud.ArcaLauncher.dll"], bootstrap_sha256=file_hash(ROOT / "cloud/cert_bridge/Launcher.cs"), certificate_secret="fixture-arca-certificate:1", certificate_password_secret="fixture-arca-password:1", thumbprint="A" * 40, store="CurrentUser"),
        database=dict(instance="fixture-db", database="fixture-app", user="fixture-app", tier="db-custom-2-7680", backup_time="03:00", retained_backups=7, password_secret="fixture-db-password:1", migration_baseline=0, migration_last=len(inv["migration_paths"]), migration_timeout_seconds=600),
        storage={key: dict(bucket="fixture-" + key, retention_days=365 if key == "retained" else 0 if key == "controller" else 30, soft_delete_days=7) for key in ("quarantine", "retained", "controller")},
        identity=dict(browser_access="IAP_DIRECT", issuer="https://identity.example.test", audience="fixture-api", web_url="https://fixture-web-123456789012.us-central1.run.app", api_url="https://fixture-api-123456789012.us-central1.run.app", allowed_web_principals=["group:developers@example.invalid"], session_secret="fixture-session:1", migration_service_account=account("migration"), scheduler_service_account=account("scheduler"), controller_service_account=account("controller")),
        scheduler=dict(retention_schedule="0 2 * * *", timezone="Etc/UTC", retry_count=0, max_retry_seconds=0), documents=dict(mode="REVIEW_ONLY", automatic_persistence=False, profile_ref="target-evidence/document-profile.json", profile_sha256="d" * 64),
        cost_controller=dict(state_object="control/state.json", billing_table="fixture-project.billing.export", budget_policy_ref="target-evidence/budget-policy.json", budget_policy_sha256="e" * 64),
        ci=dict(library_repository="https://github.com/etor6233/elite-engineering-library.git", library_commit=read(ROOT / "cloud/library-checkout.lock.json")["commit"], project_repository="https://github.com/fixture-org/fixture-project", wif_principal="principalSet://iam.googleapis.com/projects/123456789012/locations/global/workloadIdentityPools/fixture-pool/attribute.repository/fixture-org/fixture-project", deploy_service_account=account("deploy")))
    cycle = dict(schema="elite-finops-cycle/v403.1", production_authorized=False, environment="development", project=target["project"], location="US", billing_account="AAAAAA-BBBBBB-CCCCCC", billing_table="fixture-project.billing.export", currency="USD", period_policy="UTC_CURRENT_MONTH", maximum_bytes_billed=10000000, interval_seconds=900, max_state_bytes=1048576, max_age_seconds=1800, limit_micros=1000000, journal_uri="gs://fixture-controller/control/billing.sqlite", fence_uri="gs://fixture-controller/control/finops-intents", notification={"mode":"NONE"})
    config["cost_controller"].update(job_name="fixture-finops", schedule="*/15 * * * *", timeout_seconds=900, cycle=cycle, cycle_owner_sha256=file_hash(ROOT / "cloud/finops_cycle.py"), mounts={key:{"secret_ref":"fixture-finops-"+key+":1","sha256":digest(cycle) if key=="config" else "f"*64} for key in ("config","approval","tools","runner")}, seed_receipt_ref="target-evidence/finops-seed.json", seed_receipt_sha256="f"*64)
    artifacts = {}
    for index, family in enumerate(("api", "web", "workers", "migrations", "arca", "controller")):
        sha = ("a" if family == "workers" else "abcdef"[index]) * 64
        uri = "us-central1-docker.pkg.dev/fixture-project/runtime/" + ("api" if family == "workers" else family) + "@sha256:" + sha
        entries = ["node"] if family == "web" else ["psql"] if family == "migrations" else ["dotnet"] if family == "arca" else ["python3", "/app/cloud/finops_cycle.py"] if family == "controller" else ["/app/bin/" + name for name in stack.ENTRYPOINTS]
        artifacts[family] = dict(image=uri, promoted_image=uri, target_sha256=digest(target), source_commit=source["source_commit"], composition_sha256=source["composition_sha256"], inventory_sha256=source["inventory_sha256"], entrypoints=entries)
    return dict(config=config, images=dict(schema="elite-cloud-stack-images/v403.1", method="SIMULATED_ARTIFACTS", artifacts=artifacts))

class StackTests(unittest.TestCase):
    def setUp(self):
        # The checked-in receipt fixture belongs to its original qualification.
        # Explicit consumer/source overrides create new simulated fixtures;
        # they never rewrite the historical file or claim its old receipt.
        self.bundle = fixture_bundle() if "ELITE_CLOUD_SOURCE_ROOT" in os.environ or SOURCE == ROOT else read(Path(__file__).with_name("stack.fixture.json"))
        self.c, self.i = self.bundle["config"], self.bundle["images"]
        stack.build_stack(self.c, self.i)  # A broken positive control must not make negatives appear to pass.
    def build(self): return stack.build_stack(self.c, self.i)
    def reject(self):
        with self.assertRaises(ValueError): self.build()
    def test_deterministic_and_no_execution_claim(self):
        before = raw(self.bundle); one = self.build(); two = self.build()
        self.assertEqual(raw(one), raw(two)); self.assertEqual(before, raw(self.bundle))
        self.assertEqual(one["execution"], "NOT_RUN"); self.assertFalse(one["limits"]["config_is_runtime_proof"])
        self.assertEqual(set(one["configs"]), stack.COMPONENTS)
    def test_workload_semantics(self):
        p = self.build(); m = p["manifests"]
        api = m["manifests/electromobility-api.yaml"]
        self.assertEqual(api["kind"], "Service")
        self.assertEqual(api["spec"]["template"]["metadata"]["annotations"]["run.googleapis.com/cpu-throttling"], "false")
        self.assertEqual(m["manifests/web.yaml"]["spec"]["template"]["metadata"]["annotations"]["autoscaling.knative.dev/minScale"], "0")
        for key in ("return-effect-worker", "return-exchange-worker", "return-accounting-worker"):
            obj = m["manifests/" + key + ".yaml"]
            self.assertEqual(obj["kind"], "WorkerPool")
            self.assertNotIn("ports", obj["spec"]["template"]["spec"]["containers"][0])
        self.assertEqual(m["manifests/customer-survey-retention.yaml"]["kind"], "Job")
        self.assertFalse(any("meta-lead-import" in k or "warranty-profile" in k for k in m))
    def test_selected_migration_gaps_are_preserved(self):
        p = self.build(); selected = p["configs"]["migrations"]["paths"]
        self.assertEqual(len(selected), 84)
        self.assertEqual([int(Path(x).name[:4]) for x in selected], [*range(1, 63), *range(65, 87)])
        self.c["database"]["migration_baseline"] = 62
        pending = self.build()["configs"]["migrations"]["paths"]
        self.assertEqual(pending, selected[62:]); self.assertTrue(Path(pending[0]).name.startswith("0065_"))
    def test_migrations_do_not_retry_or_wrap_existing_transactions(self):
        obj = self.build()["manifests"]["manifests/migrations.yaml"]
        spec = obj["spec"]["template"]["spec"]["template"]["spec"]
        self.assertEqual(spec["maxRetries"], 0)
        self.assertNotIn("--single-transaction", spec["containers"][0]["args"])
        self.assertTrue(spec["containers"][0]["args"][3].startswith("/migrations/0001_"))
    def test_provider_scheduler_uses_google_api_oauth(self):
        p = self.build(); cmd = next(x for x in p["commands"] if x["id"] == "scheduler")
        self.assertTrue(any(x.startswith("--oauth-service-account-email=") for x in cmd["argv"]))
        self.assertFalse(any(x.startswith("--oidc") for x in cmd["argv"]))
        self.assertIn("retention-invoker", cmd["depends_on"])
    def test_finops_finite_job_seed_dependency_and_exact_mounts(self):
        p = self.build(); fc = self.c["cost_controller"]
        template = p["manifests"]["manifests/finops-controller.yaml"]["spec"]["template"]
        self.assertEqual(template["metadata"], {})
        self.assertEqual(template["spec"]["taskCount"], 1)
        self.assertEqual(template["spec"]["parallelism"], 1)
        spec = template["spec"]["template"]["spec"]
        self.assertEqual(spec["maxRetries"], 0)
        self.assertEqual(spec["timeoutSeconds"], 900)
        self.assertEqual(spec["serviceAccountName"], self.c["identity"]["controller_service_account"])
        self.assertEqual(spec["containers"][0]["command"], ["python3"])
        self.assertEqual(spec["containers"][0]["args"][:2], ["/app/cloud/finops_cycle.py", "run"])
        self.assertEqual(len(spec["volumes"]), 4)
        self.assertTrue(all(v["secret"]["items"][0]["key"] == "1" for v in spec["volumes"]))
        self.assertNotIn("ports", spec["containers"][0])
        configure = next(x for x in p["commands"] if x["id"] == "configure-finops-controller")
        self.assertIn("finops-seed", configure["depends_on"])
        self.assertNotIn("database", configure["depends_on"])
        seed = next(x for x in p["commands"] if x["id"] == "finops-seed")
        self.assertEqual(seed["effect"], "BARRIER"); self.assertEqual(seed["argv"], [])
        initialization = p["configs"]["cost-controller"]["seed_initialization"]
        self.assertIn("--if-generation-match=0", initialization["upload_argv"])
        self.assertIn(fc["cycle"]["journal_uri"], initialization["upload_argv"])
    def test_finops_scheduler_oauth_no_retry_and_scope_iam(self):
        p = self.build(); commands = {x["id"]:x for x in p["commands"]}
        schedule = commands["finops-scheduler"]["argv"]
        self.assertIn("--max-retry-attempts=0", schedule)
        self.assertIn("--max-retry-duration=0s", schedule)
        self.assertIn("*/15 * * * *", schedule)
        self.assertTrue(any(x.startswith("--oauth-service-account-email=") for x in schedule))
        self.assertIn("--permissions=storage.objects.get,storage.objects.create,storage.objects.delete", commands["finops-journal-role"]["argv"])
        journal = next(x for x in commands["finops-journal-access"]["argv"] if x.startswith("--condition="))
        self.assertIn("objects/control/billing.sqlite'", journal)
        self.assertNotIn("state.json", journal)
        self.assertIn("--role=roles/storage.objectCreator", commands["finops-intent-create"]["argv"])
        table = next(x for x in commands["finops-billing-reader"]["argv"] if x.startswith("--condition="))
        self.assertIn("resource.type == 'bigquery.googleapis.com/Table'", table)
        self.assertIn("projects/fixture-project/datasets/billing/tables/export", table)
        self.assertNotIn("finops-pubsub-publisher", commands)
    def test_finops_pubsub_only_explicit_notification(self):
        fc = self.c["cost_controller"]
        fc["cycle"]["notification"] = {"mode":"PUBSUB", "topic":"projects/fixture-project/topics/billing-alerts", "dedupe_seconds":3600}
        fc["mounts"]["config"]["sha256"] = digest(fc["cycle"])
        p = self.build(); commands = {x["id"]:x for x in p["commands"]}
        self.assertIn("finops-pubsub-publisher", commands["configure-finops-controller"]["depends_on"])
        self.assertIn("projects/fixture-project/topics/billing-alerts", commands["finops-pubsub-publisher"]["argv"])
        self.assertIn("--role=roles/pubsub.publisher", commands["finops-pubsub-publisher"]["argv"])
    def test_finops_mutable_mount_drift_scope_limit_and_schedule_rejected(self):
        original = copy.deepcopy(self.c["cost_controller"])
        for mutate in (lambda x:x["mounts"]["approval"].update(secret_ref="fixture-finops-approval:latest"), lambda x:x["mounts"]["config"].update(sha256="0"*64), lambda x:x["cycle"].update(journal_uri="gs://foreign-bucket/control/billing.sqlite"), lambda x:x["cycle"].update(maximum_bytes_billed=True), lambda x:x.update(schedule="* * * * *"), lambda x:x.update(timeout_seconds=901), lambda x:x.update(cycle_owner_sha256="0"*64)):
            self.c["cost_controller"] = copy.deepcopy(original); mutate(self.c["cost_controller"]); self.reject()
        self.c["cost_controller"] = original
        del self.i["artifacts"]["controller"]; self.reject()
    def test_finops_cannot_reuse_deployer_or_database_identity(self):
        old = self.c["identity"]["controller_service_account"]
        for account in (self.c["identity"]["scheduler_service_account"], self.c["ci"]["deploy_service_account"], self.c["workloads"]["electromobility-api"]["service_account"]):
            self.c["identity"]["controller_service_account"] = account; self.reject()
        self.c["identity"]["controller_service_account"] = old
    def test_least_privilege_and_secret_order(self):
        p = self.build(); identity = next(x for x in p["commands"] if x["id"] == "identity")
        sql = [x for x in p["commands"] if x["id"].startswith("sql-client-")]
        for x in sql:
            self.assertFalse(any("fixture-web@" in a or "fixture-scheduler@" in a or "fixture-controller@" in a for a in x["argv"]))
        for x in p["commands"]:
            if x["id"].startswith("secret-binding-"): self.assertIn(x["id"], identity["depends_on"])
    def test_arca_enabled_four_containers_and_private_secret_paths(self):
        self.c["arca"]["enabled"] = True
        for key in ("arca-fiscal-worker", "arca-parameter-worker", "return-fiscal-worker"): self.c["workloads"][key]["enabled"] = True
        p = self.build(); spec = p["manifests"]["manifests/electromobility-api.yaml"]["spec"]["template"]["spec"]
        self.assertEqual(len(spec["containers"]), 4)
        for container in spec["containers"]: self.assertIn({"name": "arca-uds", "mountPath": "/run/arca"}, container["volumeMounts"])
        dotnet = spec["containers"][-1]
        self.assertEqual(dotnet["command"], ["dotnet", "/app/Elite.Cloud.ArcaLauncher.dll"])
        env = {x["name"]: x.get("value") for x in dotnet["env"]}
        self.assertEqual(env["ARCA_PFX_FILE"], "/var/run/secrets/arca/certificate/cert.pfx")
        self.assertEqual(env["ARCA_PFX_PASSWORD_FILE"], "/var/run/secrets/arca/password/password")
        self.assertNotIn("ARCA_PFX_PASSWORD", env)
    def test_cas_bucket_stays_mutable_and_all_buckets_versioned(self):
        p = self.build()
        cmd = next(x for x in p["commands"] if x["id"] == "storage-controller")
        self.assertFalse(any(x.startswith("--retention-period") for x in cmd["argv"]))
        self.assertEqual(len([x for x in p["commands"] if x["id"].startswith("storage-versioning-")]), 3)
        self.c["storage"]["controller"]["retention_days"] = 30; self.reject()
    def test_cloud_transport_and_public_scope_configured(self):
        p = self.build(); env = p["configs"]["web"]["environment"]
        self.assertEqual(env["ELITE_CLOUD_RUN_IDENTITY"], "metadata")
        self.assertEqual(env["ELITE_CLOUD_RUN_AUDIENCE"], self.c["identity"]["api_url"])
        del self.c["web"]["environment"]["ENTERPRISE_TENANT_CODE"]; self.reject()
    def test_arca_partial_enable_rejected(self):
        self.c["arca"]["enabled"] = True; self.reject()
    def test_arca_bootstrap_drift_rejected(self):
        self.c["arca"]["enabled"] = True
        for key in ("arca-fiscal-worker", "arca-parameter-worker", "return-fiscal-worker"): self.c["workloads"][key]["enabled"] = True
        self.c["arca"]["bootstrap_sha256"] = "f" * 64; self.reject()
    def test_foreign_service_url_rejected(self):
        self.c["identity"]["api_url"] = "https://foreign.example.test"; self.reject()
    def test_provider_service_missing_profile_rejected(self):
        self.c["workloads"]["social-publishing"]["enabled"] = True; self.reject()
    def test_linked_source_root_rejected(self):
        import os
        evidence = ROOT / "cloud/evidence/stack"
        with tempfile.TemporaryDirectory(prefix="link-fixture-", dir=evidence) as temp:
            link = Path(temp) / "composition-link"
            try: os.symlink(SOURCE, link, target_is_directory=True)
            except OSError as error: self.skipTest("Symlink creation unavailable: " + str(error))
            with self.assertRaises(ValueError): stack.build_stack(self.c, self.i, source_root=link)
    def test_production_or_scope_rejected(self):
        for k, v in (("production_authorized", True), ("environment", "production"), ("scope", "FRANCHISE")):
            with self.subTest(k=k):
                old=self.c[k]; self.c[k]=v; self.reject(); self.c[k]=old
    def test_omitted_required_component_rejected(self):
        del self.c["documents"]; self.reject()
    def test_disabled_required_worker_rejected(self):
        self.c["workloads"]["return-accounting-worker"]["enabled"] = False; self.reject()
    def test_cli_cannot_be_daemon(self):
        self.c["workloads"]["meta-lead-import"]["enabled"] = True; self.reject()
    def test_duplicate_worker_identity_cannot_scale(self):
        self.c["workloads"]["return-effect-worker"]["max_instances"] = 2; self.reject()
    def test_background_api_scale_zero_rejected(self):
        self.c["workloads"]["electromobility-api"]["min_instances"] = 0; self.reject()
    def test_inline_secret_and_latest_rejected(self):
        self.c["web"]["environment"]["OIDC_CLIENT_SECRET"] = "not-a-secret"; self.reject()
        del self.c["web"]["environment"]["OIDC_CLIENT_SECRET"]
        self.c["web"]["secret_refs"]["OIDC_CLIENT_SECRET"] = "fixture-client:latest"; self.reject()
    def test_cloud_run_reserved_port_rejected(self):
        self.c["web"]["environment"]["PORT"] = "8080"; self.reject()
    def test_missing_oidc_client_rejected(self):
        del self.c["web"]["secret_refs"]["OIDC_CLIENT_SECRET"]; self.reject()
    def test_mutable_or_different_promoted_image_rejected(self):
        a = self.i["artifacts"]["api"]; old=a["image"]
        a["image"] = "image:latest"; self.reject(); a["image"] = old
        a["promoted_image"] = old[:-1] + "b"; self.reject()
    def test_wrong_target_or_source_digest_rejected(self):
        for key in ("target_sha256", "source_commit", "inventory_sha256", "composition_sha256"):
            with self.subTest(key=key):
                a=self.i["artifacts"]["api"]; old=a[key]; a[key]="f" * len(old); self.reject(); a[key]=old
    def test_undeclared_executable_rejected(self):
        self.i["artifacts"]["workers"]["entrypoints"] = ["/bin/sh"]; self.reject()
    def test_iap_preserves_browser_origin_and_api_workload_iam(self):
        p = self.build()
        web = p["manifests"]["manifests/web.yaml"]; api = p["manifests"]["manifests/electromobility-api.yaml"]
        self.assertEqual(web["metadata"]["annotations"]["run.googleapis.com/iap-enabled"], "true")
        self.assertNotIn("run.googleapis.com/iap-enabled", api["metadata"]["annotations"])
        self.assertEqual(p["operator_access_url"], self.c["identity"]["web_url"])
        self.assertEqual(p["operator_access_argv"], [])
        role = next(x for x in p["commands"] if x["id"] == "web-invoker-0")
        self.assertIn("--role=roles/iap.httpsResourceAccessor", role["argv"])
        self.assertIn("iap-invoker", role["depends_on"])
    def test_iam_only_does_not_claim_browser_infrastructure(self):
        self.c["identity"]["browser_access"] = "IAM_API_ONLY"
        p = self.build(); self.assertEqual(p["browser_journey_status"], "INFRASTRUCTURE_INCOMPLETE_PROTECTED_API_ONLY")
        self.assertNotIn("operator_access_url", p)
    def test_public_web_principal_rejected(self):
        self.c["identity"]["allowed_web_principals"] = ["allUsers"]; self.reject()
    def test_wrong_wif_repository_rejected(self):
        self.c["ci"]["wif_principal"] += "-different"; self.reject()
    def test_document_autopersist_rejected(self):
        self.c["documents"]["automatic_persistence"] = True; self.reject()
    def test_plan_missing_dependency_or_cycle_rejected(self):
        p = self.build(); p["commands"][0]["depends_on"] = ["absent"]
        with self.assertRaises(ValueError): stack.validate_stack(p)
        p = self.build(); p["commands"][0]["depends_on"] = [p["commands"][0]["id"]]
        with self.assertRaises(ValueError): stack.validate_stack(p)
    def test_emit_exact_files_and_preserve_existing_output(self):
        evidence = ROOT / "cloud/evidence/stack"; evidence.mkdir(parents=True, exist_ok=True)
        with tempfile.TemporaryDirectory(prefix="fixture-", dir=evidence) as temp:
            target = Path(temp) / "generated"
            fixture_path = Path(temp) / "input.json"; fixture_path.write_bytes(raw(self.bundle))
            result = stack.emit(fixture_path, target, source_root=SOURCE)
            self.assertEqual(result["result"], "PASS")
            p=read(target / "stack-plan.json")
            for key, value in p["manifests"].items(): self.assertEqual((target / key).read_bytes(), raw(value))
            before=file_hash(target / "stack-plan.json")
            with self.assertRaises(ValueError): stack.emit(fixture_path, target, source_root=SOURCE)
            self.assertEqual(before, file_hash(target / "stack-plan.json"))

if __name__ == "__main__": unittest.main(verbosity=2)
