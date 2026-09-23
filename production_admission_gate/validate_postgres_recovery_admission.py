from __future__ import annotations

import argparse
from datetime import datetime
import json
from pathlib import Path
import re
import sys

from validate_k6_load_admission import atomic, canonical_sha, proof, read_json, safe, utc, verify_proof

KINDS = ("VERIFY_BACKUP", "RESTORE_START", "RESTORE_BOOT", "DATA_INTEGRITY", "PITR")
TOOLS = {"VERIFY_BACKUP":"PostgreSQL pg_verifybackup", "RESTORE_START":"PostgreSQL pg_ctl", "RESTORE_BOOT":"PostgreSQL pg_isready", "DATA_INTEGRITY":"PostgreSQL psql", "PITR":"PostgreSQL psql"}
SOURCE = "https://ftp.postgresql.org/pub/source/v18.6/"
PROFILE_FIELDS = {"schema","project_id","environment","release_digest","evaluated_at","expires_at","executor_ref","target","policy","drill","runs","output"}
RUN_FIELDS = {"kind","execution_receipt","expected_arguments","input_files"}
LSN = re.compile(r"[0-9A-F]+/[0-9A-F]+")


def seconds(later: datetime, earlier: datetime, label: str) -> float:
    value = (later - earlier).total_seconds()
    if value < 0: raise ValueError(f"{label} cannot be negative")
    return value


def validate(profile: object, root: Path) -> dict[str, object]:
    if not isinstance(profile, dict) or set(profile) != PROFILE_FIELDS: raise ValueError("profile fields must be exact")
    if profile["schema"] != "elite-postgres-recovery-admission/v1" or profile["environment"] != "production": raise ValueError("profile identity is invalid")
    if not isinstance(profile["project_id"], str) or not profile["project_id"].strip() or not re.fullmatch(r"sha256:[0-9a-f]{64}", str(profile["release_digest"])): raise ValueError("project/release identity is invalid")
    evaluated = utc(profile["evaluated_at"], "evaluated_at"); expires = utc(profile["expires_at"], "expires_at")
    if not evaluated < expires or seconds(expires, evaluated, "validity") > 30 * 86400: raise ValueError("admission validity must be positive and at most 30 days")
    if not isinstance(profile["executor_ref"], str) or not profile["executor_ref"].strip() or not isinstance(profile["target"], dict) or set(profile["target"]) != {"cluster","backup_id","recovery_target_time","restore_cluster"}: raise ValueError("exact real recovery target is required")
    target = profile["target"]
    if any(not isinstance(target[key], str) or not target[key].strip() for key in target): raise ValueError("recovery target values are required")
    recovery_target = utc(target["recovery_target_time"], "target.recovery_target_time")

    policy_path = safe(root, profile["policy"], "policy"); policy = read_json(policy_path, "policy")
    if set(policy) != {"schema","project_id","environment","release_digest","approved_at","approvers","maximum_rpo_seconds","maximum_rto_seconds"} or policy.get("schema") != "elite-postgres-recovery-policy/v1": raise ValueError("recovery policy fields/schema must be exact")
    for key in ("project_id","environment","release_digest"):
        if policy.get(key) != profile[key]: raise ValueError(f"policy {key} mismatch")
    utc(policy["approved_at"], "policy.approved_at")
    approvers = policy["approvers"]
    if not isinstance(approvers, list) or len(approvers) < 2 or len(set(approvers)) != len(approvers) or any(not isinstance(x, str) or not x.strip() for x in approvers): raise ValueError("policy requires two distinct approvers")
    max_rpo = policy["maximum_rpo_seconds"]; max_rto = policy["maximum_rto_seconds"]
    if isinstance(max_rpo, bool) or isinstance(max_rto, bool) or not isinstance(max_rpo, (int,float)) or not isinstance(max_rto, (int,float)) or max_rpo < 0 or max_rto <= 0: raise ValueError("RPO/RTO policy values are invalid")

    drill_path = safe(root, profile["drill"], "drill"); drill = read_json(drill_path, "drill")
    drill_fields = {"schema","project_id","environment","release_digest","backup_id","recovery_target_time","recovered_through","restore_started_at","restore_ready_at","wal_archive","invariants"}
    if set(drill) != drill_fields or drill.get("schema") != "elite-postgres-recovery-drill/v1": raise ValueError("drill fields/schema must be exact")
    for key in ("project_id","environment","release_digest"):
        if drill.get(key) != profile[key]: raise ValueError(f"drill {key} mismatch")
    if drill.get("backup_id") != target["backup_id"] or drill.get("recovery_target_time") != target["recovery_target_time"]: raise ValueError("drill target identity mismatch")
    recovered = utc(drill["recovered_through"], "drill.recovered_through"); started = utc(drill["restore_started_at"], "drill.restore_started_at"); ready = utc(drill["restore_ready_at"], "drill.restore_ready_at")
    rpo = seconds(recovery_target, recovered, "achieved RPO"); rto = seconds(ready, started, "achieved RTO")
    if rpo > max_rpo: raise ValueError("achieved RPO exceeds approved maximum")
    if rto > max_rto: raise ValueError("achieved RTO exceeds approved maximum")
    wal = drill["wal_archive"]
    if not isinstance(wal, dict) or set(wal) != {"timeline","start_lsn","end_lsn","location_digest"} or not isinstance(wal["timeline"], int) or wal["timeline"] < 1 or not LSN.fullmatch(str(wal["start_lsn"])) or not LSN.fullmatch(str(wal["end_lsn"])) or not re.fullmatch(r"sha256:[0-9a-f]{64}", str(wal["location_digest"])): raise ValueError("WAL archive identity is invalid")
    invariants = drill["invariants"]
    if not isinstance(invariants, list) or not invariants or any(not isinstance(x, dict) or set(x) != {"name","pass"} or not isinstance(x["name"], str) or not x["name"].strip() or x["pass"] is not True for x in invariants): raise ValueError("all named data invariants must pass")

    runs = profile["runs"]
    if not isinstance(runs, list) or len(runs) != 5 or any(not isinstance(x, dict) or set(x) != RUN_FIELDS for x in runs) or [x["kind"] for x in runs] != list(KINDS): raise ValueError("five ordered recovery runs are required")
    evidence = [proof(root, policy_path), proof(root, drill_path)]
    for run in runs:
        kind = run["kind"]; receipt_path = safe(root, run["execution_receipt"], f"{kind}.receipt"); receipt = read_json(receipt_path, f"{kind}.receipt")
        if receipt.get("schema") != "elite-official-tool-execution/v1" or receipt.get("project_id") != profile["project_id"] or receipt.get("environment") != "production" or receipt.get("release_digest") != profile["release_digest"] or receipt.get("control_id") != "POSTGRES_RECOVERY" or receipt.get("exit_code") != 0: raise ValueError(f"{kind} execution identity/exit mismatch")
        tool = receipt.get("tool")
        if not isinstance(tool, dict) or set(tool) != {"name","version","source","sha256"} or tool.get("name") != TOOLS[kind] or tool.get("version") != "18.6" or tool.get("source") != SOURCE or not re.fullmatch(r"[0-9a-f]{64}", str(tool.get("sha256"))): raise ValueError(f"{kind} exact PostgreSQL tool identity mismatch")
        args = run["expected_arguments"]
        if not isinstance(args, list) or any(not isinstance(x, str) or not x for x in args) or receipt.get("arguments_sha256") != canonical_sha(args): raise ValueError(f"{kind} arguments are not hash-bound")
        inputs = run["input_files"]
        if not isinstance(inputs, list) or len(inputs) != len(set(inputs)): raise ValueError(f"{kind} input_files are invalid")
        input_proofs = [proof(root, safe(root, item, f"{kind}.input")) for item in inputs]
        expected_target = dict(target); expected_target.update({"run_kind":kind, "input_sha256s":[item["sha256"] for item in input_proofs]})
        if receipt.get("target") != expected_target: raise ValueError(f"{kind} target/input binding mismatch")
        for channel in ("stdout","stderr"): verify_proof(root, receipt.get(channel), f"{kind}.{channel}")
        evidence.extend([proof(root, receipt_path), receipt["stdout"], receipt["stderr"], *input_proofs])

    return {"id":"POSTGRES_RECOVERY", "project_id":profile["project_id"], "environment":"production", "release_digest":profile["release_digest"], "result":"PASS", "executed_at":profile["evaluated_at"], "expires_at":profile["expires_at"], "executor_ref":profile["executor_ref"], "tool":{"name":"PostgreSQL recovery toolchain", "version":"18.6", "source":SOURCE, "digest":"sha256:" + canonical_sha([run["kind"] for run in runs])}, "target":target, "assertions":{"pg_verifybackup_pass":True,"restore_boot_pass":True,"data_integrity_pass":True,"pitr_pass":True,"rpo_met":True,"rto_met":True}, "evidence":evidence}


def main() -> int:
    parser=argparse.ArgumentParser(); parser.add_argument("--project-root",required=True,type=Path); parser.add_argument("--profile",required=True); args=parser.parse_args()
    try:
        root=args.project_root.resolve(strict=True); profile=read_json(safe(root,args.profile,"profile"),"profile"); output=safe(root,profile.get("output"),"output",exists=False); output.parent.resolve(strict=True).relative_to(root); receipt=validate(profile,root); atomic(output,(json.dumps(receipt,indent=2,sort_keys=True)+"\n").encode()); print(f"POSTGRES_RECOVERY_ADMISSION_PASS project={receipt['project_id']} release={receipt['release_digest']}"); return 0
    except (ValueError,OSError) as error: print(f"POSTGRES_RECOVERY_ADMISSION_FAILED: {error}",file=sys.stderr); return 1


if __name__ == "__main__": sys.exit(main())
