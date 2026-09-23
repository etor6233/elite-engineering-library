"""AUTHORED owner-binding glue. Never copies progress, grants, approvals, or execution state."""
from __future__ import annotations
import argparse
import copy
import json
from pathlib import Path
from validate_business_functions import MAINTENANCE_PREFIXES, check_owners, confined, load_json, require

def bind_existing_owners(source, owners, project_root, output_ref):
    record = copy.deepcopy(source)
    state = load_json(confined(project_root, owners["execution_state"]))
    project_id = state.get("project", {}).get("id")
    require(type(project_id) is str and project_id and not project_id.startswith(MAINTENANCE_PREFIXES),
            "MAINTENANCE_STATE_IN_CONSUMER")
    record["owner_refs"] = copy.deepcopy(owners)
    record["owner_binding"]["scope"] = "CONSUMER"
    record["owner_binding"]["project_id"] = project_id
    check_owners(record, project_root, check_files=True)
    confined(project_root, output_ref, must_exist=False)
    record["permission_policy"]["effective_grants"] = []
    for function in record["functions"]:
        function["owner_refs"] = copy.deepcopy(owners)
        function["owner_binding_scope"] = "CONSUMER"
        function["open_requirement"]["owner_ref"] = owners["progress"]
        function["resume"]["read_refs"] = [
            owners["execution_state"], owners["execution_events"], owners["progress"],
            owners["failures"], output_ref,
        ]
        function["status"]["verification_status"] = {
            "scaffold": "NOT_TESTED", "target": "NOT_TESTED", "evidence_refs": [],
        }
        for task in function["tasks"]:
            task["owner_ref"] = owners["progress"]
            task["product_permission"]["granted"] = False
            for acceptance in task["acceptance"]:
                acceptance["target_status"] = "NOT_TESTED"
                acceptance["evidence_refs"] = []
    return record

def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--source", type=Path, required=True)
    parser.add_argument("--project-root", type=Path, required=True)
    parser.add_argument("--owner-map", type=Path, required=True)
    parser.add_argument("--output", required=True)
    args = parser.parse_args()
    output = confined(args.project_root, args.output, must_exist=False)
    require(not output.exists(), "OUTPUT_ALREADY_EXISTS")
    record = bind_existing_owners(load_json(args.source), load_json(args.owner_map),
                                  args.project_root, args.output)
    output.parent.mkdir(parents=True, exist_ok=True)
    with output.open("x", encoding="utf-8", newline="\n") as stream:
        stream.write(json.dumps(record, ensure_ascii=False, indent=2) + "\n")
    print(json.dumps({"result": "OWNERS_BOUND_ONLY", "output": str(output),
                      "resume_authorized": False, "production_authorized": False,
                      "next": "Run the consumer execution-state validator and this contract validator."}))
    return 0

if __name__ == "__main__":
    raise SystemExit(main())
