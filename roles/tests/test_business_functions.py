"""Focused adversarial tests of the AUTHORED contract validator and owner binder."""
import copy
import json
from pathlib import Path
import sys
import tempfile
import unittest

ROLES = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROLES))
from validate_business_functions import validate, load_json, confined, sha256
from bind_consumer_owners import bind_existing_owners

LIBRARY = Path(sys.argv.pop(sys.argv.index("--library-root") + 1)) if "--library-root" in sys.argv else Path(r"C:/Users/NL/Desktop/Public Elite Codes")
if "--library-root" in sys.argv:
    sys.argv.remove("--library-root")

class ContractsTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.original = load_json(ROLES / "business-functions.v403.json")
        cls.sources = load_json(ROLES / "method-sources.v403.json")
        cls.bindings = load_json(ROLES / "library-bindings.v403.json")["bindings"]
        cls.observations = load_json(ROLES / "evidence/source-observations.json")

    def setUp(self):
        self.contract = copy.deepcopy(self.original)
        self.temp = tempfile.TemporaryDirectory(prefix="contract-fixture-", dir=ROLES / "evidence")
        self.project = Path(self.temp.name).resolve()
        self.assertTrue(self.project.is_relative_to((ROLES / "evidence").resolve()))
        for key, ref in self.contract["owner_refs"].items():
            path = self.project / ref
            path.parent.mkdir(parents=True, exist_ok=True)
            if key == "execution_state":
                path.write_text(json.dumps({"project": {"id": "library-extension-v403"}}), encoding="utf-8")
            else:
                path.write_text("SYNTHETIC CONTRACT VALIDATION FIXTURE\n", encoding="utf-8")
        (self.project / "roles/tests").mkdir(parents=True)
        (self.project / "roles/tests/test_business_functions.py").write_text("# synthetic existing test reference\n", encoding="utf-8")
        self.path = self.project / "roles/business-functions.v403.json"
        self.persist()

    def tearDown(self):
        self.temp.cleanup()

    def persist(self):
        self.path.write_text(json.dumps(self.contract, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")

    def run_validation(self, sources=None, bindings=None):
        return validate(self.contract, self.sources if sources is None else sources, self.project, LIBRARY,
                        contract_path=self.path, bindings=self.bindings if bindings is None else bindings,
                        observations=self.observations)

    def rejected(self, code):
        with self.assertRaisesRegex(ValueError, code):
            self.run_validation()

    def test_actual_eleven_contracts_and_library_bindings(self):
        result = self.run_validation()
        self.assertEqual((result["function_count"], result["task_count"], result["acceptance_count"]), (11, 22, 22))
        self.assertFalse(result["business_operations_verified"])
        self.assertFalse(result["frontend_render_verified"])

    def test_track_cannot_be_auth_role(self):
        self.contract["functions"][3]["auth_role"] = True
        self.rejected("TRACK_IS_NOT_ROLE")

    def test_grok_is_learning(self):
        self.contract["functions"][0]["function_kind"] = "BUSINESS_FUNCTION"
        self.rejected("TRACK_IS_NOT_ROLE")

    def test_global_grants_forbidden(self):
        self.contract["permission_policy"]["effective_grants"] = ["admin:*"]
        self.rejected("IMPLICIT_PRIVILEGE")

    def test_inherited_privilege_forbidden(self):
        self.contract["permission_policy"]["inherited_privilege_allowed"] = True
        self.rejected("IMPLICIT_PRIVILEGE")

    def test_unknown_privilege_field_rejected(self):
        self.contract["functions"][3]["is_admin"] = True
        self.rejected("FUNCTION_FIELDS")

    def test_product_grant_forbidden(self):
        self.contract["functions"][5]["tasks"][0]["product_permission"]["granted"] = True
        self.rejected("PRODUCT_PERMISSION_GRANT")

    def test_permission_scope_and_wildcards_rejected(self):
        for key, value in (("organization_scope", "ALL"), ("resource_scope", "*"), ("id", "admin:*")):
            with self.subTest(key=key):
                self.contract = copy.deepcopy(self.original)
                self.contract["functions"][5]["tasks"][0]["product_permission"][key] = value
                self.rejected("PRODUCT_PERMISSION_GRANT")

    def test_observed_permission_must_exist_in_pack(self):
        self.contract["functions"][8]["tasks"][0]["product_permission"]["id"] = "root:admin"
        self.rejected("PERMISSION_NOT_IN_OWNER")

    def test_missing_frontend_binding_rejected(self):
        self.contract["functions"][2]["tasks"][0]["frontend"]["route"] = ""
        self.rejected("TASK_FRONTEND_BINDING")

    def test_wrong_task_frontend_binding_rejected(self):
        self.contract["functions"][2]["tasks"][0]["frontend"]["task_id"] = "sales.fake"
        self.rejected("TASK_FRONTEND_BINDING")

    def test_frontend_effect_cannot_be_enabled(self):
        self.contract["functions"][8]["tasks"][0]["frontend"]["effect_enabled"] = True
        self.rejected("TASK_FRONTEND_BINDING")

    def test_all_declared_frontend_sections_required(self):
        self.contract["frontend_contract"]["required_sections"].remove("permissions")
        self.rejected("FRONTEND_COVERAGE")

    def test_missing_acceptance_rejected(self):
        self.contract["functions"][9]["tasks"][0]["acceptance"] = []
        self.rejected("MISSING_ACCEPTANCE")

    def test_duplicate_function_rejected(self):
        self.contract["functions"][1] = copy.deepcopy(self.contract["functions"][0])
        self.rejected("FUNCTION_SET")

    def test_duplicate_task_rejected(self):
        self.contract["functions"][0]["tasks"][1] = copy.deepcopy(self.contract["functions"][0]["tasks"][0])
        self.rejected("TASK_ID")

    def test_proven_requires_evidence(self):
        self.contract["functions"][0]["status"]["verification_status"]["scaffold"] = "PROVEN"
        self.rejected("PROVEN_WITHOUT_EVIDENCE")

    def test_proven_rejects_wrong_hash_or_scope(self):
        verification = self.contract["functions"][0]["status"]["verification_status"]
        verification.update(scaffold="PROVEN", evidence_refs=["roles/proof.json"])
        self.persist()
        proof_path = self.project / "roles/proof.json"
        for proof in [
            {"result": "PASS", "claim_scope": "FUNCTION_CONTRACT_SCAFFOLD", "contract_sha256": "0"*64, "function_ids": ["grok-bot-101"]},
            {"result": "PASS", "claim_scope": "PRODUCTION", "contract_sha256": sha256(self.path), "function_ids": ["grok-bot-101"]},
            {"result": "CONDITIONED", "claim_scope": "FUNCTION_CONTRACT_SCAFFOLD", "contract_sha256": sha256(self.path), "function_ids": ["grok-bot-101"]},
        ]:
            proof_path.write_text(json.dumps(proof), encoding="utf-8")
            self.rejected("UNBOUND_PROOF")

    def test_proven_evidence_is_narrowly_bound(self):
        verification = self.contract["functions"][0]["status"]["verification_status"]
        verification.update(scaffold="PROVEN", evidence_refs=["roles/proof.json"])
        self.persist()
        proof = {"result": "PASS", "claim_scope": "FUNCTION_CONTRACT_SCAFFOLD", "contract_sha256": sha256(self.path),
                 "function_ids": ["grok-bot-101"]}
        (self.project / "roles/proof.json").write_text(json.dumps(proof), encoding="utf-8")
        self.assertEqual(self.run_validation()["result"], "PASS")
        self.contract["functions"][0]["purpose"] += " changed"
        self.persist()
        self.rejected("UNBOUND_PROOF")

    def test_target_proven_is_out_of_scope(self):
        self.contract["functions"][0]["status"]["verification_status"]["target"] = "PROVEN"
        self.rejected("TARGET_PROOF_OUT_OF_SCOPE")

    def test_future_galaxy_content_rejected(self):
        self.contract["functions"][0]["galaxy_admission"]["source_ref"] = "imagined-event-method"
        self.rejected("FUTURE_GALAXY_CLAIM")

    def test_current_methods_do_not_wait_for_event(self):
        result = self.run_validation()
        self.assertFalse(result["galaxy_event_admitted"])
        self.assertEqual(self.contract["functions"][1]["status"]["method_admission"]["status"], "OFFICIAL_METHOD_REFERENCED")

    def test_official_code_reattribution_rejected(self):
        self.contract["provenance"]["official_method_code_claim"] = True
        self.rejected("FALSE_OFFICIAL_CODE")

    def test_unofficial_method_rejected(self):
        sources = copy.deepcopy(self.sources)
        sources["sources"][0]["url"] = "https://example.invalid/method"
        with self.assertRaisesRegex(ValueError, "UNOFFICIAL_SOURCE"):
            self.run_validation(sources=sources)

    def test_unknown_source_rejected(self):
        self.contract["functions"][0]["source_refs"].append("FUTURE")
        self.rejected("UNKNOWN_METHOD_SOURCE")

    def test_drift_in_historical_pack_rejected(self):
        bindings = copy.deepcopy(self.bindings)
        bindings["engineering"]["pack_sha256"] = "0"*64
        with self.assertRaisesRegex(ValueError, "LIBRARY_BINDING_CHANGED"):
            self.run_validation(bindings=bindings)

    def test_missing_owner_and_path_escape_rejected(self):
        for reference, error in (("missing.md", "MISSING_REFERENCE"), ("../other.md", "REFERENCE_ESCAPE"),
                                 ("C:/outside.md", "REFERENCE_SYNTAX")):
            self.contract = copy.deepcopy(self.original)
            self.contract["owner_refs"]["failures"] = reference
            self.rejected(error)

    def test_resume_must_reference_actual_owners(self):
        self.contract["functions"][0]["resume"]["read_refs"].remove("PROJECT_FAILURE_LESSONS.md")
        self.rejected("RESUME_OWNERS")

    def test_consumer_cannot_inherit_maintenance_state(self):
        with self.assertRaisesRegex(ValueError, "MAINTENANCE_STATE_IN_CONSUMER"):
            bind_existing_owners(self.contract, self.contract["owner_refs"], self.project, "roles/consumer.json")

    def test_consumer_cannot_inherit_maintenance_tasks(self):
        state_path = self.project / self.contract["owner_refs"]["execution_state"]
        state_path.write_text(json.dumps({"project": {"id": "consumer-actual"}}), encoding="utf-8")
        with self.assertRaisesRegex(ValueError, "MAINTENANCE_OWNER_IN_CONSUMER"):
            bind_existing_owners(self.contract, self.contract["owner_refs"], self.project, "roles/consumer.json")

    def test_consumer_rebinds_existing_owners_and_resets_claims(self):
        owner_map = copy.deepcopy(self.contract["owner_refs"])
        owner_map["progress"] = "specs/customer-journey/tasks.md"
        (self.project / owner_map["progress"]).parent.mkdir(parents=True)
        (self.project / owner_map["progress"]).write_text("# Existing consumer task owner\n", encoding="utf-8")
        state_path = self.project / owner_map["execution_state"]
        state_path.write_text(json.dumps({"project": {"id": "consumer-actual"}}), encoding="utf-8")
        before = {ref: (self.project / ref).read_bytes() for ref in owner_map.values()}
        self.contract["functions"][0]["status"]["verification_status"]["evidence_refs"] = ["old-proof.json"]
        rebound = bind_existing_owners(self.contract, owner_map, self.project, "roles/consumer.json")
        self.assertEqual(rebound["owner_binding"]["scope"], "CONSUMER")
        self.assertEqual(rebound["functions"][0]["status"]["verification_status"]["evidence_refs"], [])
        self.assertEqual(before, {ref: (self.project / ref).read_bytes() for ref in owner_map.values()})
        self.contract = rebound
        self.path = self.project / "roles/consumer.json"
        self.persist()
        self.assertEqual(self.run_validation()["result"], "PASS")

    def test_duplicate_json_keys_rejected(self):
        duplicate = self.project / "duplicate.json"
        duplicate.write_text('{"scope":"A","scope":"B"}', encoding="utf-8")
        with self.assertRaisesRegex(ValueError, "DUPLICATE_JSON_KEY"):
            load_json(duplicate)

    def test_nonfinite_json_rejected(self):
        invalid = self.project / "nonfinite.json"
        invalid.write_text('{"metric":NaN}', encoding="utf-8")
        with self.assertRaisesRegex(ValueError, "NONFINITE_JSON"):
            load_json(invalid)

if __name__ == "__main__":
    unittest.main(verbosity=2)
