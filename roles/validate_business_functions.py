"""AUTHORED contract glue. Validates library scaffolding, never product permission or readiness."""
from __future__ import annotations

import argparse
import copy
import hashlib
import json
from pathlib import Path, PurePosixPath
import re
import sys
from urllib.parse import urlsplit

MAX_JSON_BYTES = 2 * 1024 * 1024
EXPECTED_TRACKS = {
    "grok-bot-101": "Grok Bot 101", "engineering": "Engineering",
    "product-managers": "Product Managers", "founders": "Founders",
    "sales-engineering": "Sales Engineering", "sales": "Sales", "sdrs": "SDRs",
    "customer-support": "Customer Support", "marketing-operations": "Marketing Operations",
    "post-sales": "Post-Sales", "marketing": "Marketing",
}
OWNER_KEYS = {"progress", "execution_state", "execution_events", "failures", "dependency", "freshness"}
OWNER_FIXTURE = "specs/library-extension/tasks.md"
MAINTENANCE_PREFIXES = ("elite-library", "library-extension-v403")
OFFICIAL_HOSTS = {"learn.microsoft.com", "dora.dev", "knowledge.hubspot.com", "docs.x.ai"}
SECTIONS = {"purpose", "systems", "tasks", "permissions", "allowed", "prohibited", "inputs",
            "outputs", "acceptance", "implementation", "sources", "status", "owners", "resume"}
ROOT_KEYS = {"schema_version", "contract_id", "scope", "claim_scope", "production_authorized",
             "provenance", "frontend_contract", "owner_refs", "permission_policy", "done_policy",
             "functions", "owner_binding", "actions_disabled_in_catalog",
             "catalog_restriction_scope", "consumer_authorization"}
FUNCTION_KEYS = {"id", "track_name", "function_kind", "auth_role", "purpose", "systems", "allowed",
                 "prohibited", "status", "source_refs", "galaxy_admission", "owner_refs", "resume",
                 "implementation_binding", "open_requirement", "tasks", "owner_binding_scope"}

def require(condition, code):
    if not condition:
        raise ValueError(code)

def closed(value, keys, code):
    require(type(value) is dict and set(value) == set(keys), code)

def strings(value, code, maximum=32):
    require(type(value) is list and 0 < len(value) <= maximum, code)
    require(all(type(v) is str and v.strip() and len(v) <= 4000 for v in value), code)

def sha256(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()

def pairs(items):
    result = {}
    for key, value in items:
        require(key not in result, "DUPLICATE_JSON_KEY")
        result[key] = value
    return result

def load_json(path):
    require(path.stat().st_size <= MAX_JSON_BYTES, "JSON_SIZE_LIMIT")
    return json.loads(path.read_text(encoding="utf-8"), object_pairs_hook=pairs,
                      parse_constant=lambda value: (_ for _ in ()).throw(ValueError("NONFINITE_JSON")))

def confined(root, reference, must_exist=True):
    require(type(reference) is str and reference.strip(), "MISSING_REFERENCE")
    require("\\" not in reference and ":" not in reference and "\x00" not in reference,
            "REFERENCE_SYNTAX")
    relative = PurePosixPath(reference)
    require(not relative.is_absolute() and ".." not in relative.parts and "." != reference,
            "REFERENCE_ESCAPE")
    base = root.resolve()
    path = (base / reference).resolve()
    require(path.is_relative_to(base), "REFERENCE_ESCAPE")
    if must_exist:
        require(path.is_file(), "MISSING_REFERENCE")
    return path

def check_owners(record, project_root, check_files):
    owners = record["owner_refs"]
    closed(owners, OWNER_KEYS, "OWNER_FIELDS")
    binding = record["owner_binding"]
    closed(binding, {"scope", "project_id", "require_existing_refs", "fixture_owners_must_not_be_inherited",
                     "resume_gate", "consumer_integration"}, "OWNER_BINDING_FIELDS")
    require(binding["scope"] in {"MAINTENANCE_FIXTURE", "CONSUMER"}, "OWNER_BINDING_SCOPE")
    require(binding["require_existing_refs"] is True and binding["fixture_owners_must_not_be_inherited"] is True,
            "OWNER_INHERITANCE_CONTROL")
    require(type(binding["project_id"]) is str and binding["project_id"].strip(), "OWNER_PROJECT_ID")
    for reference in owners.values():
        confined(project_root, reference, must_exist=check_files)
    require(len(set(owners.values())) == len(owners), "OWNER_COLLISION")
    if binding["scope"] == "MAINTENANCE_FIXTURE":
        require(binding["project_id"] == "library-extension-v403" and owners["progress"] == OWNER_FIXTURE,
                "FIXTURE_OWNER_DRIFT")
    else:
        forbidden = ("specs/library-extension/", "specs/library-maintenance/")
        require(not binding["project_id"].startswith(MAINTENANCE_PREFIXES), "MAINTENANCE_ID_IN_CONSUMER")
        require(not any(any(part in ref.lower() for part in forbidden) for ref in owners.values()),
                "MAINTENANCE_OWNER_IN_CONSUMER")
    if check_files:
        state = load_json(confined(project_root, owners["execution_state"]))
        require(state.get("project", {}).get("id") == binding["project_id"], "OWNER_STATE_ID_MISMATCH")
        if binding["scope"] == "CONSUMER":
            require(not state["project"]["id"].startswith(MAINTENANCE_PREFIXES), "MAINTENANCE_STATE_IN_CONSUMER")
    return owners

def validate(record, sources, project_root, library_root, check_files=True,
             contract_path=None, bindings=None, observations=None):
    closed(record, ROOT_KEYS, "ROOT_FIELDS")
    require(record["schema_version"] == "1.0.0" and record["contract_id"] == "BUSINESS-FUNCTIONS-V403",
            "SCHEMA_ID")
    require(record["scope"] == "LIBRARY_INFRASTRUCTURE" and
            record["claim_scope"] == "FUNCTION_CONTRACT_SCAFFOLD" and
            record["production_authorized"] is False, "SCOPE_ESCALATION")
    closed(record["provenance"], {"classification", "kind", "upstream_code_copied", "official_method_code_claim"},
           "PROVENANCE_FIELDS")
    require(record["provenance"] == {"classification": "AUTHORED", "kind": "CONTRACT_AND_VALIDATOR_GLUE",
                                    "upstream_code_copied": False, "official_method_code_claim": False},
            "FALSE_OFFICIAL_CODE")
    owners = check_owners(record, project_root, check_files)
    policy = record["permission_policy"]
    closed(policy, {"track_is_auth_role", "default", "effective_grants", "requires",
                    "unbound_permission_behavior", "inherited_privilege_allowed"}, "PERMISSION_POLICY_FIELDS")
    require(policy["track_is_auth_role"] is False and policy["default"] == "DENY" and
            policy["effective_grants"] == [] and policy["inherited_privilege_allowed"] is False,
            "IMPLICIT_PRIVILEGE")
    require(set(policy["requires"]) == {"authenticated_subject", "explicit_product_permission",
            "organization_scope", "resource_scope", "server_side_enforcement", "audit_evidence"},
            "AUTHORIZATION_DIMENSIONS")
    require(policy["unbound_permission_behavior"] == "BLOCK_PRODUCT_EFFECT", "UNBOUND_EFFECT")
    frontend = record["frontend_contract"]
    closed(frontend, {"route_pattern", "mode", "required_sections", "external_effects_enabled"}, "FRONTEND_FIELDS")
    require(frontend["route_pattern"] == "#/functions/{id}" and
            frontend["mode"] == "READ_ONLY_LIBRARY_CONSOLE" and frontend["external_effects_enabled"] is False,
            "CATALOG_EFFECT")
    require(set(frontend["required_sections"]) == SECTIONS, "FRONTEND_COVERAGE")
    require(record["catalog_restriction_scope"] == "READ_ONLY_LIBRARY_CONSOLE_ONLY", "CATALOG_SCOPE")
    require(set(record["actions_disabled_in_catalog"]) == {"SEND_MESSAGE", "PUBLISH", "SPEND", "MUTATE_CRM",
                                                         "RUN_PRODUCT_EFFECT"}, "CATALOG_EFFECT")
    require(type(record["consumer_authorization"]) is str and record["consumer_authorization"].strip(),
            "CONSUMER_AUTHORITY")
    closed(record["done_policy"], {"scaffold", "product", "galaxy"}, "DONE_POLICY")
    require(all(type(v) is str and v.strip() for v in record["done_policy"].values()), "DONE_POLICY")

    require(sources.get("scope") == "DOCUMENTARY_METHOD_REFERENCES" and sources.get("upstream_code_acquired") is False,
            "SOURCE_SCOPE")
    source_list = sources.get("sources", [])
    source_ids = [s["id"] for s in source_list]
    require(source_ids and len(source_ids) == len(set(source_ids)), "SOURCE_DUPLICATE")
    for source in source_list:
        parts = urlsplit(source["url"])
        require(parts.scheme == "https" and parts.hostname in OFFICIAL_HOSTS and not parts.username,
                "UNOFFICIAL_SOURCE")
        require(source.get("method_status") == "OFFICIAL_METHOD_REFERENCED" and
                source.get("upstream_code_copied") is False and
                source.get("architecture_code_admission") == "NOT_REQUESTED", "SOURCE_PROMOTION")
        require(all(type(source.get(k)) is str and source[k].strip() for k in
                    ("query", "claim", "limit", "publisher", "title", "freshness")), "SOURCE_PROVENANCE")
    if check_files:
        require(type(observations) is dict, "MISSING_SOURCE_OBSERVATIONS")
        obs = {o["id"]: o for o in observations.get("observations", [])}
        require(len(obs) == len(source_ids), "SOURCE_OBSERVATION_COUNT")
        for source in source_list:
            receipt = obs.get(source["id"], {})
            require(receipt.get("status") == "OBSERVED" and receipt.get("http_status") == 200 and
                    receipt.get("url") == source["url"] and receipt.get("body_stored") is False and
                    re.fullmatch("[0-9a-f]{64}", receipt.get("response_sha256", "")) is not None,
                    "SOURCE_OBSERVATION_MISMATCH")

    functions = record["functions"]
    require(type(functions) is list and len(functions) == 11, "FUNCTION_COUNT")
    require({f.get("id") for f in functions} == set(EXPECTED_TRACKS), "FUNCTION_SET")
    task_ids, acceptance_ids = set(), set()
    for function in functions:
        closed(function, FUNCTION_KEYS, "FUNCTION_FIELDS")
        fid = function["id"]
        require(function["track_name"] == EXPECTED_TRACKS[fid], "TRACK_NAME")
        kind = "LEARNING_TRACK" if fid == "grok-bot-101" else "BUSINESS_FUNCTION"
        require(function["function_kind"] == kind and function["auth_role"] is False, "TRACK_IS_NOT_ROLE")
        require(function["owner_binding_scope"] == record["owner_binding"]["scope"] and
                function["owner_refs"] == owners, "OWNER_BINDING_MISMATCH")
        require(type(function["purpose"]) is str and function["purpose"].strip(), "PURPOSE")
        for key in ("systems", "allowed", "prohibited", "source_refs"):
            strings(function[key], "FUNCTION_" + key.upper())
        require(set(function["source_refs"]).issubset(source_ids), "UNKNOWN_METHOD_SOURCE")
        closed(function["status"], {"specification_status", "method_admission", "verification_status"}, "STATUS_FIELDS")
        require(function["status"]["specification_status"] == "SPEC_READY", "SPECIFICATION_STATUS")
        method = function["status"]["method_admission"]
        require(method == {"status": "OFFICIAL_METHOD_REFERENCED", "claim_scope": "DOCUMENTARY_METHOD_ONLY",
                           "source_refs": function["source_refs"], "code_admission": "NONE"}, "METHOD_PROMOTION")
        verification = function["status"]["verification_status"]
        closed(verification, {"scaffold", "target", "evidence_refs"}, "VERIFICATION_FIELDS")
        require(verification["target"] == "NOT_TESTED", "TARGET_PROOF_OUT_OF_SCOPE")
        require(verification["scaffold"] in {"NOT_TESTED", "PROVEN"}, "VERIFICATION_STATUS")
        require(type(verification["evidence_refs"]) is list, "EVIDENCE_LIST")
        if verification["scaffold"] == "PROVEN":
            require(bool(verification["evidence_refs"]) and contract_path is not None, "PROVEN_WITHOUT_EVIDENCE")
            for reference in verification["evidence_refs"]:
                proof = load_json(confined(project_root, reference))
                require(proof.get("result") == "PASS" and proof.get("claim_scope") == "FUNCTION_CONTRACT_SCAFFOLD"
                        and proof.get("contract_sha256") == sha256(contract_path)
                        and fid in proof.get("function_ids", []), "UNBOUND_PROOF")
        galaxy = function["galaxy_admission"]
        closed(galaxy, {"status", "source_ref", "applies_only_to", "content_invented", "reopen_trigger"}, "GALAXY_FIELDS")
        require(galaxy["status"] == "WAITING_GALAXY_ADMISSION" and galaxy["source_ref"] is None and
                galaxy["applies_only_to"] == "FUTURE_EVENT_MATERIAL" and galaxy["content_invented"] is False
                and type(galaxy["reopen_trigger"]) is str and bool(galaxy["reopen_trigger"].strip()),
                "FUTURE_GALAXY_CLAIM")
        closed(function["resume"], {"read_refs", "rule"}, "RESUME_FIELDS")
        require(set(function["resume"]["read_refs"]) >= {owners["execution_state"], owners["execution_events"],
                owners["progress"], owners["failures"]}, "RESUME_OWNERS")
        for reference in function["resume"]["read_refs"]:
            confined(project_root, reference, must_exist=check_files)
        open_req = function["open_requirement"]
        closed(open_req, {"id", "scope", "description", "owner_ref"}, "OPEN_REQUIREMENT_FIELDS")
        require(open_req["scope"] == "FUTURE_CONSUMER_ONLY" and open_req["owner_ref"] == owners["progress"]
                and type(open_req["description"]) is str and open_req["description"].strip(), "OPEN_REQUIREMENT")
        impl = function["implementation_binding"]
        closed(impl, {"status", "pack", "evidence", "action", "kind", "target_verified", "v403_executes_binding"},
               "IMPLEMENTATION_FIELDS")
        require(impl["target_verified"] is False and impl["v403_executes_binding"] is False, "IMPLICIT_RUNTIME")
        require(impl["status"] in {"EXISTING_LIBRARY_REFERENCE", "DOCUMENTED_ONLY"}, "IMPLEMENTATION_STATUS")
        if impl["status"] == "DOCUMENTED_ONLY":
            require(impl["pack"] is None and impl["evidence"] is None and impl["action"] is None
                    and impl["kind"] == "NO_PRODUCT_RUNTIME", "UNDOCUMENTED_RUNTIME")
        else:
            require(type(impl["action"]) is str and impl["action"].strip() and
                    impl["kind"] in {"CLI_CONTRACT", "PACK_REFERENCE", "OBSERVED_API_ROUTE"}, "IMPLEMENTATION_ACTION")
            pack = confined(library_root, impl["pack"], must_exist=check_files)
            evidence = confined(library_root, impl["evidence"], must_exist=check_files)
            if check_files:
                require(type(bindings) is dict and fid in bindings, "MISSING_BINDING_RECEIPT")
                bound = bindings[fid]
                require(bound.get("pack") == impl["pack"] and bound.get("evidence") == impl["evidence"] and
                        bound.get("pack_sha256") == sha256(pack) and bound.get("evidence_sha256") == sha256(evidence),
                        "LIBRARY_BINDING_CHANGED")
                if impl["kind"] == "OBSERVED_API_ROUTE":
                    require("/v1/franchise/marketing/campaigns" in pack.read_text(encoding="utf-8"),
                            "ROUTE_NOT_IN_OWNER")
        require(type(function["tasks"]) is list and len(function["tasks"]) == 2, "TASK_COUNT")
        for task in function["tasks"]:
            closed(task, {"id", "title", "inputs", "outputs", "product_permission", "frontend", "acceptance",
                          "owner_ref"}, "TASK_FIELDS")
            tid = task["id"]
            require(type(tid) is str and tid.startswith(fid + ".") and tid not in task_ids, "TASK_ID")
            task_ids.add(tid)
            require(task["owner_ref"] == owners["progress"], "TASK_OWNER")
            require(type(task["title"]) is str and task["title"].strip(), "TASK_TITLE")
            strings(task["inputs"], "INPUT_CONTRACT")
            strings(task["outputs"], "OUTPUT_CONTRACT")
            permission = task["product_permission"]
            closed(permission, {"id", "definition_kind", "granted", "subject_binding", "organization_scope",
                                "resource_scope", "enforcement", "policy_ref"}, "PRODUCT_PERMISSION_FIELDS")
            require(re.fullmatch("[a-z][a-z-]*:[a-z][a-z-]*", permission["id"]) is not None and
                    permission["granted"] is False and
                    permission["subject_binding"] == "REQUIRES_TARGET_CONFIGURATION" and
                    permission["organization_scope"] == "REQUIRED" and permission["resource_scope"] == "REQUIRED",
                    "PRODUCT_PERMISSION_GRANT")
            require(permission["enforcement"] == "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP", "ENFORCEMENT")
            require(permission["definition_kind"] in {"PROPOSED_TARGET_REQUIREMENT", "OBSERVED_EXISTING_OWNER_NAME"},
                    "PERMISSION_DEFINITION")
            confined(library_root, permission["policy_ref"], must_exist=check_files)
            if check_files and permission["definition_kind"] == "OBSERVED_EXISTING_OWNER_NAME":
                require(impl["status"] == "EXISTING_LIBRARY_REFERENCE" and permission["id"] in
                        confined(library_root, impl["pack"]).read_text(encoding="utf-8"), "PERMISSION_NOT_IN_OWNER")
            ui = task["frontend"]
            require(ui == {"route": "#/functions/" + fid, "section": "tasks", "task_id": tid,
                           "interaction": "READ_CONTRACT", "effect_enabled": False}, "TASK_FRONTEND_BINDING")
            require(type(task["acceptance"]) is list and task["acceptance"], "MISSING_ACCEPTANCE")
            for acceptance in task["acceptance"]:
                closed(acceptance, {"id", "claim_scope", "criterion", "test_ref", "evidence_refs", "target_status"},
                       "ACCEPTANCE_FIELDS")
                aid = acceptance["id"]
                require(type(aid) is str and aid.startswith(fid + ".") and aid not in acceptance_ids, "ACCEPTANCE_ID")
                acceptance_ids.add(aid)
                require(acceptance["claim_scope"] == "CONTRACT_SPECIFICATION" and
                        acceptance["target_status"] == "NOT_TESTED" and acceptance["evidence_refs"] == [],
                        "ACCEPTANCE_SCOPE")
                require(type(acceptance["criterion"]) is str and acceptance["criterion"].strip(), "ACCEPTANCE_CRITERION")
                confined(project_root, acceptance["test_ref"], must_exist=check_files)
    return {"result": "PASS", "claim_scope": "FUNCTION_CONTRACT_SCAFFOLD", "function_ids": sorted(EXPECTED_TRACKS),
            "function_count": len(functions), "task_count": len(task_ids), "acceptance_count": len(acceptance_ids),
            "owner_binding_scope": record["owner_binding"]["scope"], "production_authorized": False,
            "business_operations_verified": False, "frontend_render_verified": False,
            "galaxy_event_admitted": False, "known_open_scaffold_defects": "SEE_EXISTING_FAILURE_OWNER",
            "scope_limit": "Structural/provenance/link controls only; frontend execution and human/business acceptance are separate."}

def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--project-root", type=Path, required=True)
    parser.add_argument("--library-root", type=Path, required=True)
    parser.add_argument("--contract", default="roles/business-functions.v403.json")
    parser.add_argument("--report", required=True)
    args = parser.parse_args()
    output = confined(args.project_root, args.report, must_exist=False)
    require(not output.exists(), "REPORT_ALREADY_EXISTS")
    try:
        contract_path = confined(args.project_root, args.contract)
        source_path = confined(args.project_root, "roles/method-sources.v403.json")
        result = validate(load_json(contract_path), load_json(source_path), args.project_root, args.library_root,
                          contract_path=contract_path,
                          bindings=load_json(confined(args.project_root, "roles/library-bindings.v403.json"))["bindings"],
                          observations=load_json(confined(args.project_root, "roles/evidence/source-observations.json")))
        result.update(contract_sha256=sha256(contract_path), sources_sha256=sha256(source_path),
                      validator_sha256=sha256(Path(__file__)), python_version=sys.version)
        exit_code = 0
    except (ValueError, KeyError, TypeError, OSError, json.JSONDecodeError) as exc:
        result = {"result": "FAIL", "claim_scope": "FUNCTION_CONTRACT_SCAFFOLD", "error": str(exc),
                  "production_authorized": False}
        exit_code = 2
    output.parent.mkdir(parents=True, exist_ok=True)
    with output.open("x", encoding="utf-8", newline="\n") as stream:
        stream.write(json.dumps(result, ensure_ascii=False, indent=2) + "\n")
    print(json.dumps(result, ensure_ascii=False))
    return exit_code

if __name__ == "__main__":
    raise SystemExit(main())
