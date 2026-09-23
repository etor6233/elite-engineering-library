"""AUTHORED profile materialization glue; no credentials or commercial policy invention."""
import argparse
import hashlib
import json
from pathlib import Path
import re
import uuid


def main():
    p = argparse.ArgumentParser(description="Materialize the supported initial-handover algorithm and exact activation lock")
    p.add_argument("--output", required=True, type=Path)
    p.add_argument("--profile-id", required=True)
    p.add_argument("--revision", type=int, default=1)
    p.add_argument("--tenant-id", required=True)
    p.add_argument("--organization-id", required=True)
    p.add_argument("--payment-provider", choices=("stripe", "mercadopago"), required=True)
    p.add_argument("--payment-account-ref", required=True)
    p.add_argument("--payment-connection-id", required=True)
    p.add_argument("--expected-mode", choices=("sandbox", "live"), required=True)
    p.add_argument("--maximum-observation-age-seconds", type=int, default=300)
    p.add_argument("--authority-reference", default="docs/initial-handover-reference.md")
    p.add_argument("--decision-reference", default="handover-policy/DECISION.md")
    p.add_argument("--release-effect", choices=("read-only", "commercial-receipt"), default="read-only", help="Select the durable receipt explicitly; no physical shipment or money posting")
    p.add_argument("--funding", choices=("provider-only", "stored-value"), default="provider-only", help="Explicitly select revision3 with immutable approved stored-value funding")
    p.add_argument("--activate", action="store_true", help="Explicitly activate this exact supported profile; does not prove live readiness")
    a = p.parse_args()
    if not re.fullmatch(r"[a-z][a-z0-9._-]{0,79}", a.profile_id) or not 1 <= a.revision <= 1000000:
        p.error("unsupported profile id/revision")
    try:
        if str(uuid.UUID(a.tenant_id)) != a.tenant_id:
            raise ValueError()
    except ValueError:
        p.error("tenant id must be canonical UUID")
    if any(not x or len(x) > 128 or x.strip() != x for x in (a.organization_id, a.authority_reference, a.decision_reference, a.payment_account_ref, a.payment_connection_id)):
        p.error("invalid organization or documentary reference")
    if not 1 <= a.maximum_observation_age_seconds <= 900:
        p.error("maximum observation age must be 1..900 seconds")
    if a.funding == "stored-value" and a.release_effect != "commercial-receipt":
        p.error("stored-value funding requires the explicit commercial receipt profile")
    mode = a.expected_mode == "live"
    policy = {
        "schema": "elite-handover-profile/v1", "profile_id": a.profile_id,
        "revision": a.revision, "algorithm": "single-unit-full-observed-payment",
        "algorithm_revision": 3 if a.funding == "stored-value" else 2 if a.release_effect == "commercial-receipt" else 1, "scope": "MATERIALIZED_PROFILE",
        "tenant_id": a.tenant_id, "organization_id": a.organization_id,
        "provider_code": a.payment_provider, "provider_account_ref": a.payment_account_ref,
        "provider_connection_id": a.payment_connection_id,
        "expected_live_mode": mode,
        "maximum_observation_age_seconds": a.maximum_observation_age_seconds,
        "options": {"quantity": 1, "payment_coverage": "FULL_ORDER_WITH_STORED_VALUE" if a.funding == "stored-value" else "FULL_ORDER",
                    "stock_selection": "ALLOCATED_SERIALIZED_UNIT",
                    "reservation": "ACTIVE_MATCHING", "refunds": "ZERO",
                    "disputes": "DENY", "acceptance": "CUSTOMER_AND_REQUIRED_CHECKLIST",
                    "release_effect": "COMMIT_COMMERCIAL_RELEASE_RECEIPT" if a.release_effect == "commercial-receipt" else "READ_ONLY_ELIGIBILITY"},
        "authority_reference": a.authority_reference, "decision_reference": a.decision_reference,
    }
    encode = lambda value: (json.dumps(value, sort_keys=True, indent=2) + "\n").encode("utf-8")
    body = encode(policy)
    activation = {"enabled": a.activate, "profile_id": a.profile_id,
                  "profile_revision": a.revision,
                  "document_sha256": hashlib.sha256(body).hexdigest(),
                  "tenant_id": a.tenant_id, "organization_id": a.organization_id,
        "provider_code": a.payment_provider, "provider_account_ref": a.payment_account_ref,
        "provider_connection_id": a.payment_connection_id,
                  "expected_live_mode": mode}
    decision = ("# Supported handover profile selection\n\n"
                f"Profile: {a.profile_id}@{a.revision}. Expected provider mode: {a.expected_mode}.\n"
                "This selects the existing single-unit/full-observed-payment algorithm.\n"
                "The options do not grant credit, change pricing/tax rules or post shipping.\n"
                f"Funding selection: {a.funding}. Selected release effect: {a.release_effect}. A receipt records a committed checkpoint; current validity requires a new generation/freshness/acceptance check.\n"
                "The receipt never posts inventory or proves physical dispatch, taxation or production readiness.\n"
                "Authority/decision references are documentary, not secret or live-readiness proof.\n"
                f"Activation explicitly selected: {a.activate}. Production readiness remains unproven.\n")
    a.output.mkdir(parents=False, exist_ok=False)
    files = {"profile.json": body, "activation.json": encode(activation), "DECISION.md": decision.encode("utf-8")}
    for name, data in files.items():
        with (a.output / name).open("xb") as f:
            f.write(data)
    receipt = {"schema": "elite-handover-profile-materialization/v1",
               "scope": "SUPPORTED_ALGORITHM_CONFIGURATION_ONLY", "live_proven": False,
               "files": {name: hashlib.sha256(data).hexdigest() for name, data in files.items()}}
    (a.output / "receipt.json").write_bytes(encode(receipt))
    print(json.dumps({"output": str(a.output.resolve()), "document_sha256": activation["document_sha256"], "enabled": a.activate}))


if __name__ == "__main__":
    main()
