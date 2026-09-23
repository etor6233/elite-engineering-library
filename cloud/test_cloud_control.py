"""Behavioral fixture tests; none of these constitute cloud target evidence."""
from __future__ import annotations
import copy
import hashlib
import io
import json
from pathlib import Path
import subprocess
import tempfile
import unittest
import zipfile
import tarfile
import os
from cloud_control import *
from bootstrap import checkout
from toolchain import install

HERE=Path(__file__).resolve().parent
NOW="2026-09-13T12:00:00Z"

def sample_config():
    return {"schema":"elite-cloud-target/v403.1","environment":"staging","platform":"linux/amd64","project":"fixture-project","service":"fixture-api","database_instance":"fixture-db","region":"us-central1","runtime_identity":"fixture-runtime@fixture-project.iam.gserviceaccount.com","private_url":True,"production_authorized":False,"max_instances":3,"min_instances":1,"cpu_always":True,"database_policy":"RETAIN_BACKUP_PITR","state_backend":"GCS_GENERATION_CAS","state_uri":"gs://fixture-retained/control/budget.json","secrets":{"DATABASE_URL":"fixture-database:1"},"environment_values":{"OIDC_ISSUER":"https://identity.example.test","OIDC_AUDIENCE":"fixture-api"}}

def snapshot(**changes):
    d={"schema":"elite-billing-snapshot/v403.1","kind":"PROVIDER_BILLING_OBSERVATION","provider":"gcp","account":"account-a","period":"2026-09","generation":1,"observed_at":NOW,"currency":"USD","complete_scope":True,"export_delay_possible":True,"rows":[{"project":"one","service":"run","amount_micros":600},{"project":"two","service":"sql","amount_micros":250}]}
    d.update(changes);d["source_proof"]={"method":"SIMULATED_PROVIDER","truncated":False,"row_count":len(d["rows"])};return d

def policy():
    return {"kind":"AGGREGATE_ALERT_POLICY","scope":[{"provider":"gcp","account":"account-a"}],"period":"2026-09","currency":"USD","limit_micros":1000,"max_age_seconds":3600}

class CloudTests(unittest.TestCase):
    def setUp(self):
        self.tmp=tempfile.TemporaryDirectory();self.root=Path(self.tmp.name);self.config=sample_config()
        self.artifact={"schema":"elite-cloud-release/v403.1","target":"linux/amd64","environment_class":"TARGET_ARTIFACT","image":"us-central1-docker.pkg.dev/fixture-project/artifacts/api@sha256:"+"a"*64,"source_commit":"b"*40,"composition_sha256":"c"*64,"required_receipts":{}}
        for name in ("build","license","sca","secret_scan","migrations","restore","runtime","identity"):
            # Intentionally fabricated *inside the isolated test* to exercise
            # validation branches. Never exported as an admission receipt.
            r={"image":self.artifact["image"],"target":"linux/amd64","result":"PASS","method":"EXECUTED_TARGET"}
            p=self.root/(name+".json");p.write_bytes(raw(r));r.update(path=p.name,sha256=file_hash(p));self.artifact["required_receipts"][name]=r
        (self.root/"before.json").write_bytes(raw({"project":self.config["project"],"region":self.config["region"],"service":self.config["service"],"revision":"fixture-api-before","traffic":{"fixture-api-before":100},"image":self.artifact["image"].replace("a"*64,"d"*64),"method":"EXECUTED_TARGET"}))

    def tearDown(self):self.tmp.cleanup()
    def plan(self):return deploy_plan(self.config,self.artifact,"fixture-api-before",self.root)
    def test_api_background_cannot_scale_to_zero(self):
        self.config["min_instances"]=0
        with self.assertRaisesRegex(ValueError,"background"):target(self.config)
    def test_request_only_cpu_rejected(self):
        self.config["cpu_always"]=False
        with self.assertRaisesRegex(ValueError,"CPU"):target(self.config)
    def test_plan_is_private_digest_no_traffic_and_never_executes(self):
        p=self.plan();s=str(p);self.assertIn("--no-traffic",s);self.assertIn("--no-allow-unauthenticated",s);self.assertIn(self.artifact["image"],s);self.assertEqual(p["execution"],"NOT_RUN")
    def test_mutable_image_rejected(self):
        self.artifact["image"]="us-central1-docker.pkg.dev/fixture-project/artifacts/api:latest"
        with self.assertRaisesRegex(ValueError,"immutable"):self.plan()
    def test_windows_original_zip_cannot_be_promoted(self):
        self.artifact["target"]="windows/amd64"
        with self.assertRaisesRegex(ValueError,"Windows"):self.plan()
    def test_missing_receipt_fails(self):
        self.artifact["required_receipts"].pop("identity")
        with self.assertRaisesRegex(ValueError,"receipts"):self.plan()
    def test_fixture_receipt_cannot_authorize_deploy(self):
        self.artifact["required_receipts"]["build"]["method"]="SIMULATED_PROVIDER"
        with self.assertRaisesRegex(ValueError,"fixture"):self.plan()
    def test_tampered_receipt_fails(self):
        (self.root/"build.json").write_text("{}")
        with self.assertRaisesRegex(ValueError,"tampered"):self.plan()
    def test_prod_never_authorized(self):
        self.config["environment"]="production"
        with self.assertRaisesRegex(ValueError,"production"):self.plan()
    def test_plain_secret_and_mutable_version_rejected(self):
        self.config["secrets"]["DATABASE_URL"]="postgres://secret"
        with self.assertRaisesRegex(ValueError,"Secret Manager"):self.plan()
        self.config["secrets"]["DATABASE_URL"]="fixture-database:latest"
        with self.assertRaisesRegex(ValueError,"Secret Manager"):self.plan()
    def test_command_injection_rejected(self):
        self.config["project"]="fixture; touch /tmp/bad"
        with self.assertRaisesRegex(ValueError,"invalid target"):self.plan()
    def test_ephemeral_journal_rejected_on_target(self):
        self.config["state_backend"]="SQLITE_CONTAINER_DISK"
        with self.assertRaisesRegex(ValueError,"ephemeral"):self.plan()
    def test_approval_bound_and_expires(self):
        p=self.plan();a={"plan_sha256":digest(p),"scope":"CLOUD_REFERENCE_EXECUTION","environment":"staging","approved_by":"operator-record","approved_at":"2026-09-13T11:00:00Z","expires_at":"2026-09-13T13:00:00Z","provider_access_proven":True,"budget_authorized":True}
        self.assertEqual(authorize(p,a,NOW)["execution"],"NOT_RUN")
        a["expires_at"]="2026-09-13T11:59:59Z"
        with self.assertRaisesRegex(ValueError,"expired"):authorize(p,a,NOW)
    def test_approval_cannot_follow_changed_plan(self):
        p=self.plan();a={"plan_sha256":digest(p)};p["image"]=p["image"].replace("a"*64,"e"*64)
        with self.assertRaises(ValueError):authorize(p,a,NOW)
    def observation(self,op):
        p=self.plan();old=self.artifact["image"].replace("a"*64,"d"*64)
        return p,{"schema":"elite-cloud-observation/v403.1","plan_sha256":digest(p),"method":"SIMULATED_PROVIDER","operation":op,"traffic":{p["candidate_revision"] if op=="promote" else p["prior_revision"]:100},"ready":True,"identity_checked":True,"data_health":True,"pending_work_reconciled":True,"image":p["image"] if op=="promote" else old,"prior_image":old}
    def test_simulated_promotion_is_never_cloud_pass(self):
        p,o=self.observation("promote");r=reconcile_deployment(p,o);self.assertEqual(r["result"],"PASS");self.assertFalse(r["cloud_demonstrated"])
    def test_rollback_restores_prior_digest(self):
        p,o=self.observation("rollback");self.assertEqual(reconcile_deployment(p,o)["operation"],"rollback");o["image"]=p["image"]
        with self.assertRaisesRegex(ValueError,"not restored"):reconcile_deployment(p,o)
    def test_failed_health_blocks_reconciliation(self):
        p,o=self.observation("promote");o["pending_work_reconciled"]=False
        with self.assertRaisesRegex(ValueError,"health"):reconcile_deployment(p,o)
    def test_health_string_false_rejected(self):
        p,o=self.observation("promote");o["ready"]="false"
        with self.assertRaisesRegex(ValueError,"health"):reconcile_deployment(p,o)
    def test_rollback_cannot_self_assert_fake_prior_digest(self):
        p,o=self.observation("rollback");o["image"]=self.artifact["image"];o["prior_image"]=o["image"]
        with self.assertRaisesRegex(ValueError,"not restored"):reconcile_deployment(p,o)
    def test_string_false_is_not_account_authorization(self):
        p=self.plan();a={"plan_sha256":digest(p),"scope":"CLOUD_REFERENCE_EXECUTION","environment":"staging","approved_by":"operator-record","approved_at":"2026-09-13T11:00:00Z","expires_at":"2026-09-13T13:00:00Z","provider_access_proven":"false","budget_authorized":True}
        with self.assertRaisesRegex(ValueError,"account and budget"):authorize(p,a,NOW)
    def test_failed_suspension_cannot_delete_data(self):
        with self.assertRaisesRegex(ValueError,"incomplete drain"):suspend_plan(self.config,{})
    def test_suspension_retain_and_resume_explicit(self):
        proof={k:True for k in ("ingress_closed","schedules_paused","webhooks_buffered_or_provider_replay_proven","new_jobs_fenced","active_jobs_zero","outbox_reconciled","backup_verified","restore_proven","resume_recipe_verified")};proof.update(target_sha256=digest(self.config),method="SIMULATED_PROVIDER")
        p=suspend_plan(self.config,proof);self.assertFalse(p["database_deleted"]);self.assertFalse(p["zero_cost_claim"]);self.assertNotIn("delete",str(p["commands"]))
    def test_object_state_uses_generation_compare_and_swap(self):
        state=self.root/"state.json";state.write_text("{}")
        p=state_cas_plan(self.config,42,state);self.assertIn("--if-generation-match=42",p["argv"])

class BillingTests(unittest.TestCase):
    def setUp(self):self.tmp=tempfile.TemporaryDirectory();self.path=Path(self.tmp.name)/"journal.db";self.j=BillingJournal(self.path)
    def tearDown(self):self.j.close();self.tmp.cleanup()
    def test_duplicate_snapshot_does_not_double_bill(self):
        self.j.import_snapshot(snapshot());self.assertEqual(self.j.import_snapshot(snapshot()),"DUPLICATE_NO_EFFECT");self.assertEqual(self.j.budget(policy(),NOW)["observed_micros"],850)
    def test_reconciliation_replaces_prior_and_keeps_signed_credits(self):
        self.j.import_snapshot(snapshot());self.j.import_snapshot(snapshot(generation=2,rows=[{"project":"one","service":"run","amount_micros":-25}]));self.assertEqual(self.j.budget(policy(),NOW)["observed_micros"],-25)
    def test_conflicting_generation_rejected(self):
        self.j.import_snapshot(snapshot())
        with self.assertRaisesRegex(ValueError,"conflicting"):self.j.import_snapshot(snapshot(rows=[]))
    def test_restart_preserves_observed_budget(self):
        self.j.import_snapshot(snapshot());self.j.close();self.j=BillingJournal(self.path);self.assertEqual(self.j.budget(policy(),NOW)["observed_micros"],850)
    def test_global_multiple_providers_and_projects(self):
        self.j.import_snapshot(snapshot());self.j.import_snapshot(snapshot(provider="external-api",account="account-b",rows=[{"project":"three","service":"ai","amount_micros":300}]))
        p=policy();p["scope"].append({"provider":"external-api","account":"account-b"});r=self.j.budget(p,NOW);self.assertEqual(r["state"],"ALERT");self.assertEqual(r["observed_micros"],1150);self.assertFalse(r["hard_cap"])
    def test_missing_or_stale_observation_is_unknown(self):
        self.assertEqual(self.j.budget(policy(),NOW)["state"],"UNKNOWN");self.j.import_snapshot(snapshot(observed_at="2026-09-12T00:00:00Z"));self.assertEqual(self.j.budget(policy(),NOW)["state"],"UNKNOWN")
    def test_cross_currency_not_silently_converted(self):
        self.j.import_snapshot(snapshot(currency="ARS"))
        with self.assertRaisesRegex(ValueError,"FX"):self.j.budget(policy(),NOW)
    def test_internal_estimate_is_not_billing(self):
        self.j.record_internal({"kind":"INTERNAL_ESTIMATE","id":"request-one","currency":"USD","amount_micros":10000});self.j.import_snapshot(snapshot());r=self.j.budget(policy(),NOW);self.assertEqual(r["observed_micros"],850);self.assertFalse(r["internal_estimates_added_to_billing"])
    def test_internal_estimate_idempotency_conflict(self):
        r={"kind":"INTERNAL_ESTIMATE","id":"request-one","currency":"USD","amount_micros":2};self.j.record_internal(r);self.j.record_internal(r);r["amount_micros"]=3
        with self.assertRaisesRegex(ValueError,"payload conflict"):self.j.record_internal(r)
    def test_billing_query_all_costs_credits_exact_period(self):
        q=billing_query("fixture-project.billing.gcp_export","202609");self.assertIn("UNNEST(credits)",q[-1]);self.assertIn("_ACCOUNT_LEVEL",q[-1]);self.assertNotIn("LIMIT",q[-1])
        with self.assertRaises(ValueError):billing_query("anything`; DROP TABLE x;--","202609")
    def test_billing_truncation_cannot_claim_complete_scope(self):
        s=snapshot();s["source_proof"]["truncated"]=True
        with self.assertRaisesRegex(ValueError,"completeness"):self.j.import_snapshot(s)
        with self.assertRaisesRegex(ValueError,"row cap"):normalize_billing_rows([{}]*100000,{"currency":"USD"})

class BootstrapTests(unittest.TestCase):
    def setUp(self):self.tmp=tempfile.TemporaryDirectory();self.root=Path(self.tmp.name)
    def tearDown(self):self.tmp.cleanup()
    def source(self):
        p=self.root/"source";p.mkdir();(p/"authority.md").write_text("fixture authority\n")
        def g(*a):return subprocess.check_output(["git",*a],cwd=p,stderr=subprocess.DEVNULL,text=True).strip()
        g("init");g("config","user.email","fixture@example.test");g("config","user.name","Fixture");g("config","core.autocrlf","false");g("add","authority.md");g("commit","-m","fixture only")
        return {"schema":"elite-library-checkout/v403.1","repository":str(p),"commit":g("rev-parse","HEAD"),"files":{"authority.md":file_hash(p/"authority.md")}}
    def test_real_distinct_git_checkout_without_desktop_path(self):
        lock=self.source();dest=self.root/"worker"/"library";r=checkout(lock,dest,fixture=True);self.assertEqual(r["commit"],lock["commit"]);self.assertEqual((dest/"authority.md").read_text(),"fixture authority\n");self.assertIn("DISABLED",subprocess.check_output(["git","remote","get-url","--push","origin"],cwd=dest,text=True))
    def test_changed_library_authority_rejected(self):
        lock=self.source();lock["files"]["authority.md"]="0"*64
        with self.assertRaisesRegex(ValueError,"locked file mismatch"):checkout(lock,self.root/"worker",fixture=True)
    def test_moving_branch_forbidden(self):
        lock=self.source();lock["commit"]="main"
        with self.assertRaisesRegex(ValueError,"exact Git"):checkout(lock,self.root/"worker",fixture=True)
    def test_tool_archive_safe_integrity_install(self):
        p=self.root/"tool.zip"
        with zipfile.ZipFile(p,"w") as z:z.writestr("bin/tool","fixture tool")
        row={"tool":"node","version":"fixture","platform":"linux/amd64","admission":"ADMITTED_EXACT_ARTIFACT","url":"https://nodejs.org/fixture.zip","sha256":file_hash(p),"format":"zip"}
        r=install(row,self.root/"installed",p);self.assertEqual(r["runtime_execution"],"NOT_RUN");self.assertTrue((self.root/"installed/bin/tool").is_file())
    def test_tool_archive_traversal_rejected_before_extract(self):
        p=self.root/"bad.zip"
        with zipfile.ZipFile(p,"w") as z:z.writestr("../escaped","bad")
        row={"tool":"node","version":"fixture","platform":"linux/amd64","admission":"ADMITTED_EXACT_ARTIFACT","url":"https://nodejs.org/fixture.zip","sha256":file_hash(p),"format":"zip"}
        with self.assertRaisesRegex(ValueError,"unsafe relative"):install(row,self.root/"installed",p)
        self.assertFalse((self.root/"escaped").exists())
    def test_official_dotnet_leading_dot_archive_is_normalized(self):
        p=self.root/"dot.tar.gz"
        with tarfile.open(p,"w:gz") as t:
            r=tarfile.TarInfo(".");r.type=tarfile.DIRTYPE;t.addfile(r)
            r=tarfile.TarInfo("./bin/tool");r.size=7;t.addfile(r,io.BytesIO(b"fixture"))
        row={"tool":"dotnet","version":"fixture","platform":"linux/amd64","admission":"ACQUISITION_ONLY","url":"https://builds.dotnet.microsoft.com/fixture.tar.gz","sha256":file_hash(p),"format":"tar.gz"}
        install(row,self.root/"installed",p);self.assertEqual((self.root/"installed/bin/tool").read_bytes(),b"fixture")
    def test_worker_parent_link_cannot_escape_project_separation(self):
        from bootstrap import bootstrap
        project=self.root/"project";project.mkdir();(project/".git").mkdir();alias=self.root/"alias"
        try:alias.symlink_to(project,target_is_directory=True)
        except OSError:self.skipTest("OS does not grant symlink creation; Linux lane executes this negative")
        with self.assertRaisesRegex(ValueError,"linked worker"):bootstrap(project,"unused","unused",alias/"worker",fixture=True)
        self.assertFalse((project/"worker").exists())

class CompositionTests(unittest.TestCase):
    def test_duplicate_or_omitted_pack_identity_cannot_pass_same_count(self):
        with tempfile.TemporaryDirectory() as tmp:
            root=Path(tmp);packs=[];locks=[]
            for name in ("one","two"):
                f=root/(name+".md");f.write_text(name)
                packs.append({"packId":name,"path":f.name,"version":"1"});locks.append(dict(packs[-1],sha256=file_hash(f)))
            plan=root/"plan.md";plan.write_text("```json\n"+json.dumps({"packs":packs})+"\n```\n")
            profile={"schema":"elite-cloud-composition/v403.1","scope":"LIBRARY_INFRASTRUCTURE","production_authorized":False,"selection":{"path":"plan.md","sha256":file_hash(plan),"pack_count":2},"pack_locks":locks,"authority_locks":{}}
            self.assertEqual(validate_profile(profile,root)["pack_count"],2)
            profile["pack_locks"]=[locks[0],locks[0]]
            with self.assertRaisesRegex(ValueError,"identity set"):validate_profile(profile,root)

class ExecutionTests(unittest.TestCase):
    setUp=CloudTests.setUp
    tearDown=CloudTests.tearDown
    def plan(self):
        p=CloudTests.plan(self);p["execution_fence"]={"backend":"LOCAL_FIXTURE","directory":str(self.root/"intents")};return p
    observation=CloudTests.observation
    def prerequisite(self,p,step):
        from execute_step import PREREQUISITES
        if step not in PREREQUISITES:return None
        image=p["prior_image"] if step=="rollback" else p["image"]
        record={"schema":"elite-cloud-step-prerequisite/v403.1","plan_sha256":digest(p),"result":"PASS","for_step":step,"method":"SIMULATED_PROVIDER","target_sha256":p["target_sha256"],"image":image,"observed_at":"2026-09-13T11:59:00Z","expires_at":"2026-09-13T12:05:00Z","checks":{k:True for k in PREREQUISITES[step]},"evidence":[]}
        for check in sorted(PREREQUISITES[step]):
            path=self.root/(step+"-"+check+".json");path.write_bytes(raw({"check":check,"result":"PASS","method":"SIMULATED_PROVIDER","target_sha256":p["target_sha256"],"image":image,"observed_at":"2026-09-13T11:59:00Z","expires_at":"2026-09-13T12:05:00Z"}));record["evidence"].append({"check":check,"path":path.name,"sha256":file_hash(path)})
        path=self.root/(step+"-prerequisite.json");path.write_bytes(raw(record));return {"path":path.name,"sha256":file_hash(path)}
    def test_generated_commands_run_existing_runner_fake_provider_and_reconcile(self):
        import sys
        from execute_step import execute
        p=self.plan();state=self.root/"provider-state.json"
        old=self.artifact["image"].replace("a"*64,"d"*64)
        state.write_bytes(raw({"revisions":{"fixture-api-before":old},"traffic":{"fixture-api-before":100}}))
        approval={"plan_sha256":digest(p),"scope":"CLOUD_REFERENCE_EXECUTION","environment":"staging","approved_by":"test-operator-record","approved_at":"2026-09-13T11:00:00Z","expires_at":"2026-09-13T13:00:00Z","provider_access_proven":True,"budget_authorized":True}
        owner=Path(os.environ["ELITE_CLOUD_OWNER_ROOT"])/"run_official_tool.py" if os.environ.get("ELITE_CLOUD_OWNER_ROOT") else HERE/"owners/secure-ops/production_admission_gate/run_official_tool.py"
        self.assertTrue(owner.is_file(),"materialize the unchanged SECURE-OPS owner before integration test")
        runner={"path":str(owner),"sha256":file_hash(owner)}
        journal=self.root/"intents";journal.mkdir()
        tool={"binary":sys.executable,"binary_sha256":file_hash(sys.executable),"method":"SIMULATED_PROVIDER","admission":"FIXTURE_ONLY","prefix_arguments":[str(HERE/"fake_provider.py"),str(state)],"journal_directory":str(journal)}
        for step in ("describe-before","deploy-no-traffic","describe-candidate","promote-after-health","rollback","describe-after"):
            pre=self.prerequisite(p,step)
            r=execute(p,step,approval,NOW,tool,runner,self.root/(step+".receipt.json"),fixture=True,prerequisite=pre)
            self.assertEqual(r["exit_code"],0);self.assertFalse(r["cloud_demonstrated"])
            if step=="deploy-no-traffic":self.assertEqual(read(state)["traffic"],{"fixture-api-before":100})
            if step in ("promote-after-health","rollback"):
                operation="promote" if step=="promote-after-health" else "rollback"
                pp,ob=self.observation(operation);current=read(state);ob["traffic"]=current["traffic"];rev=next(iter(current["traffic"]));ob["image"]=current["revisions"][rev]
                self.assertEqual(reconcile_deployment(pp,ob)["result"],"PASS")
        with self.assertRaisesRegex(ValueError,"already exists"):
            execute(p,"rollback",approval,NOW,tool,runner,self.root/"rollback.receipt.json",fixture=True,prerequisite=self.prerequisite(p,"rollback"))
        with self.assertRaises(FileExistsError):
            execute(p,"rollback",approval,NOW,tool,runner,self.root/"different-receipt.json",fixture=True,prerequisite=self.prerequisite(p,"rollback"))
    def test_effect_then_crash_has_durable_intent_and_rejects_second_execution(self):
        import sys
        from unittest.mock import patch
        from execute_step import execute
        p=self.plan();journal=self.root/"intents";journal.mkdir();count=[]
        a={"plan_sha256":digest(p),"scope":"CLOUD_REFERENCE_EXECUTION","environment":"staging","approved_by":"test-operator-record","approved_at":"2026-09-13T11:00:00Z","expires_at":"2026-09-13T13:00:00Z","provider_access_proven":True,"budget_authorized":True}
        tool={"binary":sys.executable,"binary_sha256":file_hash(sys.executable),"method":"SIMULATED_PROVIDER","admission":"FIXTURE_ONLY","journal_directory":str(journal)}
        def effect_then_crash(*args):
            count.append(1);(self.root/"provider-effect").write_text("committed");raise ValueError("simulated interrupted result")
        with patch("execute_step.load_runner",return_value=effect_then_crash):
            with self.assertRaises(ValueError):execute(p,"describe-before",a,NOW,tool,{"path":"unused","sha256":"unused"},self.root/"first.json",fixture=True)
            with self.assertRaises(FileExistsError):execute(p,"describe-before",a,NOW,tool,{"path":"unused","sha256":"unused"},self.root/"second.json",fixture=True)
        self.assertEqual(len(count),1);self.assertEqual(read(self.root/"first.json")["result"],"UNKNOWN_RECONCILE_BEFORE_RETRY")
    def test_invalid_production_plan_rejected_before_tool_load(self):
        from execute_step import execute
        p=self.plan();p["schema"]="wrong";p["production_authorized"]=True
        with self.assertRaisesRegex(ValueError,"non-production"):execute(p,"describe-before",{"plan_sha256":digest(p)},NOW,{}, {},self.root/"never.json",fixture=True)
        self.assertFalse((self.root/"never.json").exists())
    def test_write_effect_cannot_be_relabelled_as_read(self):
        p=self.plan();p["commands"][1]["effect"]="READ"
        with self.assertRaisesRegex(ValueError,"commands/effects"):authorize(p,{"plan_sha256":digest(p)},NOW)
    def test_fresh_wrapper_cannot_reuse_expired_individual_observation(self):
        from execute_step import validate_prerequisite
        p=self.plan();ref=self.prerequisite(p,"promote-after-health");wrapper=self.root/ref["path"];v=read(wrapper);e=v["evidence"][0];path=self.root/e["path"];ob=read(path);ob["observed_at"]="2020-01-01T00:00:00Z";ob["expires_at"]="2020-01-02T00:00:00Z";path.write_bytes(raw(ob));e["sha256"]=file_hash(path);wrapper.write_bytes(raw(v));ref["sha256"]=file_hash(wrapper)
        with self.assertRaisesRegex(ValueError,"individual.*stale"):validate_prerequisite(p,"promote-after-health",ref,"SIMULATED_PROVIDER",NOW)
    def test_prerequisite_target_freshness_and_evidence_required(self):
        from execute_step import validate_prerequisite
        p=self.plan();ref=self.prerequisite(p,"promote-after-health")
        validate_prerequisite(p,"promote-after-health",ref,"SIMULATED_PROVIDER",NOW)
        path=self.root/ref["path"];v=read(path);v["target_sha256"]="0"*64;path.write_bytes(raw(v));ref["sha256"]=file_hash(path)
        with self.assertRaisesRegex(ValueError,"target/image"):validate_prerequisite(p,"promote-after-health",ref,"SIMULATED_PROVIDER",NOW)
        ref=self.prerequisite(p,"promote-after-health");path=self.root/ref["path"];v=read(path);v["expires_at"]="2026-09-13T11:59:59Z";path.write_bytes(raw(v));ref["sha256"]=file_hash(path)
        with self.assertRaisesRegex(ValueError,"stale"):validate_prerequisite(p,"promote-after-health",ref,"SIMULATED_PROVIDER",NOW)
        ref=self.prerequisite(p,"promote-after-health");v=read(self.root/ref["path"]);(self.root/v["evidence"][0]["path"]).write_text("{}")
        with self.assertRaisesRegex(ValueError,"tampered"):validate_prerequisite(p,"promote-after-health",ref,"SIMULATED_PROVIDER",NOW)
    def test_concurrent_attempts_create_only_one_effect(self):
        import sys
        import concurrent.futures
        from unittest.mock import patch
        from execute_step import execute
        p=self.plan();journal=self.root/"intents";journal.mkdir();count=[]
        a={"plan_sha256":digest(p),"scope":"CLOUD_REFERENCE_EXECUTION","environment":"staging","approved_by":"test-operator-record","approved_at":"2026-09-13T11:00:00Z","expires_at":"2026-09-13T13:00:00Z","provider_access_proven":True,"budget_authorized":True}
        tool={"binary":sys.executable,"binary_sha256":file_hash(sys.executable),"method":"SIMULATED_PROVIDER","admission":"FIXTURE_ONLY","journal_directory":str(journal)}
        def effect(*args):count.append(1);return subprocess.CompletedProcess([],0,b"{}",b"")
        def attempt(i):
            try:return execute(p,"describe-before",a,NOW,tool,{"path":"unused","sha256":"unused"},self.root/(str(i)+".json"),fixture=True)["exit_code"]
            except FileExistsError:return "REJECTED"
        with patch("execute_step.load_runner",return_value=effect),concurrent.futures.ThreadPoolExecutor(max_workers=2) as pool:results=list(pool.map(attempt,[1,2]))
        self.assertEqual(len(count),1);self.assertEqual(sorted(map(str,results)),["0","REJECTED"])
    def test_tool_cannot_change_approved_fence_namespace(self):
        import sys
        from unittest.mock import patch
        from execute_step import execute
        p=self.plan();journal=self.root/"intents";journal.mkdir();other=self.root/"different";other.mkdir()
        a={"plan_sha256":digest(p),"scope":"CLOUD_REFERENCE_EXECUTION","environment":"staging","approved_by":"test-operator-record","approved_at":"2026-09-13T11:00:00Z","expires_at":"2026-09-13T13:00:00Z","provider_access_proven":True,"budget_authorized":True}
        tool={"binary":sys.executable,"binary_sha256":file_hash(sys.executable),"method":"SIMULATED_PROVIDER","admission":"FIXTURE_ONLY","journal_directory":str(other)}
        with patch("execute_step.load_runner") as runner:
            with self.assertRaisesRegex(ValueError,"namespace override"):execute(p,"describe-before",a,NOW,tool,{"path":"unused","sha256":"unused"},self.root/"out.json",fixture=True)
            runner.return_value.assert_not_called()
        self.assertFalse(list(other.iterdir()))

if __name__=="__main__":unittest.main(verbosity=2)
