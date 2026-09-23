from __future__ import annotations

import hashlib
import json
from pathlib import Path
import tempfile
import unittest

from validate_document_routing import PROVIDER_PLANS, REQUIRED_CLASS_IDS, validate_configuration, write_receipt


def complete_configuration() -> dict:
    classes = []
    for class_id in sorted(REQUIRED_CLASS_IDS):
        if class_id != "supplier-invoice":
            classes.append({"id": class_id, "decision": "NONE_WITH_REASON", "reason": "not required by test project"})
            continue
        classes.append({
            "id": class_id,
            "decision": "REQUIRED",
            "schema_version": "supplier-invoice/v1",
            "mime_types": ["application/pdf"],
            "max_bytes": 10485760,
            "max_pages": 1,
            "required_fields": ["invoice_number", "total_amount"],
            "candidate_lanes": [{
                "provider": "AWS",
                "role": "PRIMARY",
                "pack_plan": PROVIDER_PLANS["AWS"],
                "provider_profile_path": "config/aws-document-profile.json",
            }],
            "security": {
                "policy_path": "config/secure-file-policy.json",
                "required_receipt_schema": "elite-secure-local-file-receipt/v1",
            },
            "evaluation": {
                "strict_profile_path": "config/strict-field-evaluation.json",
                "status": "EVALUATION_REQUIRED",
            },
            "corpus": {
                "authorized": True,
                "ground_truth_documents": 10,
                "ground_truth_owner": "data-owner",
                "evidence_path": "evidence/corpus/supplier-invoice.json",
            },
            "access": {
                "probe_status": "PROVEN",
                "identity_reference": "secret-manager://aws-role-reference",
                "region": "us-east-1",
                "evidence_path": "evidence/access/aws-sandbox.json",
            },
            "automatic_storage": False,
        })
    return {
        "schema": "elite-document-pipeline-routing/v1",
        "project": "test-project",
        "owner": "business-owner",
        "security_owner": "security-owner",
        "data_owner": "data-owner",
        "review_operations_owner": "review-owner",
        "automatic_storage_authorized": False,
        "classes": classes,
    }


class DocumentRoutingTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary = tempfile.TemporaryDirectory()
        self.root = Path(self.temporary.name)
        self.configuration = self.root / "routing.json"

    def tearDown(self) -> None:
        self.temporary.cleanup()

    def write(self, value: dict) -> None:
        self.configuration.write_text(json.dumps(value, sort_keys=True) + "\n", encoding="utf-8", newline="\n")

    def test_complete_explicit_inventory_writes_hash_linked_receipt(self) -> None:
        self.write(complete_configuration())
        output = self.root / "evidence"
        receipt = write_receipt(self.configuration, output)
        persisted = json.loads((output / "document-routing-receipt.json").read_text(encoding="utf-8"))
        self.assertEqual(receipt, persisted)
        self.assertEqual(receipt["configuration_sha256"], hashlib.sha256(self.configuration.read_bytes()).hexdigest())
        self.assertEqual(receipt["required_classes"], ["supplier-invoice"])
        self.assertEqual(receipt["selected_pack_plans"], [PROVIDER_PLANS["AWS"]])
        self.assertEqual(receipt["schema"], "elite-document-pipeline-routing-receipt/v3")
        self.assertEqual(receipt["baseline_class_ids"], sorted(REQUIRED_CLASS_IDS))
        self.assertEqual(receipt["custom_class_ids"], [])
        self.assertEqual(receipt["selected_routes"], [{
            "class_id": "supplier-invoice",
            "class_decision": "REQUIRED",
            "provider": "AWS",
            "role": "PRIMARY",
            "pack_plan": PROVIDER_PLANS["AWS"],
            "provider_profile_path": "config/aws-document-profile.json",
        }])
        self.assertFalse(receipt["automatic_storage_authorized"])

    def test_receipt_preserves_every_class_lane_and_profile_path(self) -> None:
        value = complete_configuration()
        active = next(item for item in value["classes"] if item["decision"] == "REQUIRED")
        active["candidate_lanes"].append({
            "provider": "GOOGLE",
            "role": "EVALUATION",
            "pack_plan": PROVIDER_PLANS["GOOGLE"],
            "provider_profile_path": "config/google-document-profile.json",
        })
        self.write(value)
        receipt = write_receipt(self.configuration, self.root / "route-evidence")
        self.assertEqual(receipt["selected_pack_plans"], sorted([PROVIDER_PLANS["AWS"], PROVIDER_PLANS["GOOGLE"]]))
        self.assertEqual(
            [(route["provider"], route["role"], route["provider_profile_path"]) for route in receipt["selected_routes"]],
            [("GOOGLE", "EVALUATION", "config/google-document-profile.json"), ("AWS", "PRIMARY", "config/aws-document-profile.json")],
        )

    def test_additional_canonical_class_is_routed_with_full_contract(self) -> None:
        value = complete_configuration()
        source = next(item for item in value["classes"] if item["decision"] == "REQUIRED")
        custom = json.loads(json.dumps(source))
        custom["id"] = "letter-of-credit-amendment"
        custom["schema_version"] = "letter-of-credit-amendment/v1"
        custom["required_fields"] = ["credit_number", "amendment_number", "effective_date"]
        custom["candidate_lanes"] = [{
            "provider": "GOOGLE",
            "role": "PRIMARY",
            "pack_plan": PROVIDER_PLANS["GOOGLE"],
            "provider_profile_path": "config/google-letter-of-credit-profile.json",
        }]
        custom["corpus"]["evidence_path"] = "evidence/corpus/letter-of-credit-amendment.json"
        custom["access"]["identity_reference"] = "secret-manager://google-document-ai-reference"
        custom["access"]["region"] = "us"
        custom["access"]["evidence_path"] = "evidence/access/google-sandbox.json"
        value["classes"].append(custom)
        self.write(value)
        receipt = write_receipt(self.configuration, self.root / "custom-route-evidence")
        self.assertEqual(receipt["custom_class_ids"], ["letter-of-credit-amendment"])
        self.assertEqual(receipt["required_classes"], ["letter-of-credit-amendment", "supplier-invoice"])
        self.assertIn(PROVIDER_PLANS["GOOGLE"], receipt["selected_pack_plans"])
        self.assertEqual(
            [route["class_id"] for route in receipt["selected_routes"]],
            ["letter-of-credit-amendment", "supplier-invoice"],
        )

    def test_additional_class_id_must_be_canonical(self) -> None:
        value = complete_configuration()
        value["classes"].append({"id": "../../new class", "decision": "NONE_WITH_REASON", "reason": "invalid identifier"})
        self.write(value)
        with self.assertRaisesRegex(ValueError, "non-canonical class id"):
            validate_configuration(self.configuration)

    def test_distributed_template_is_intentionally_not_routable(self) -> None:
        template = Path(__file__).with_name("document-routing.template.json")
        with self.assertRaisesRegex(ValueError, "project is required"):
            validate_configuration(template)

    def test_missing_or_duplicate_class_fails(self) -> None:
        value = complete_configuration()
        value["classes"].pop()
        self.write(value)
        with self.assertRaisesRegex(ValueError, "inventory is incomplete"):
            validate_configuration(self.configuration)
        value = complete_configuration()
        value["classes"].append(value["classes"][0].copy())
        self.write(value)
        with self.assertRaisesRegex(ValueError, "duplicate class id"):
            validate_configuration(self.configuration)

    def test_unclassified_and_missing_reason_fail(self) -> None:
        value = complete_configuration()
        value["classes"][0] = {"id": value["classes"][0]["id"], "decision": "AWAITING_USER"}
        self.write(value)
        with self.assertRaisesRegex(ValueError, "decision is not complete"):
            validate_configuration(self.configuration)
        value = complete_configuration()
        inactive = next(item for item in value["classes"] if item["decision"] == "NONE_WITH_REASON")
        inactive["reason"] = ""
        self.write(value)
        with self.assertRaisesRegex(ValueError, "reason is required"):
            validate_configuration(self.configuration)

    def test_provider_plan_mapping_and_safe_paths_are_closed(self) -> None:
        value = complete_configuration()
        active = next(item for item in value["classes"] if item["decision"] == "REQUIRED")
        active["candidate_lanes"][0]["pack_plan"] = PROVIDER_PLANS["GOOGLE"]
        self.write(value)
        with self.assertRaisesRegex(ValueError, "does not match provider"):
            validate_configuration(self.configuration)
        active["candidate_lanes"][0]["pack_plan"] = PROVIDER_PLANS["AWS"]
        active["candidate_lanes"][0]["provider_profile_path"] = "../secret.json"
        self.write(value)
        with self.assertRaisesRegex(ValueError, "canonical relative path"):
            validate_configuration(self.configuration)

    def test_required_class_needs_security_corpus_access_and_evaluation(self) -> None:
        cases = (
            ("security", None, "security is required"),
            ("corpus", {"authorized": False}, "corpus must be explicitly authorized"),
            ("access", {"probe_status": "NOT_PROVIDED"}, "probe_status must be PROVEN"),
            ("evaluation", {"strict_profile_path": "config/eval.json", "status": "READY_FOR_AUTOMATIC_STORAGE"}, "cannot authorize storage"),
        )
        for field, replacement, message in cases:
            value = complete_configuration()
            active = next(item for item in value["classes"] if item["decision"] == "REQUIRED")
            active[field] = replacement
            self.write(value)
            with self.assertRaisesRegex(ValueError, message):
                validate_configuration(self.configuration)

    def test_storage_authorization_is_rejected_at_both_levels(self) -> None:
        value = complete_configuration()
        value["automatic_storage_authorized"] = True
        self.write(value)
        with self.assertRaisesRegex(ValueError, "deny automatic storage"):
            validate_configuration(self.configuration)
        value = complete_configuration()
        active = next(item for item in value["classes"] if item["decision"] == "REQUIRED")
        active["automatic_storage"] = True
        self.write(value)
        with self.assertRaisesRegex(ValueError, "automatic_storage must be false"):
            validate_configuration(self.configuration)

    def test_existing_output_is_not_overwritten(self) -> None:
        self.write(complete_configuration())
        output = self.root / "evidence"
        output.mkdir()
        with self.assertRaises(FileExistsError):
            write_receipt(self.configuration, output)


if __name__ == "__main__":
    unittest.main()
