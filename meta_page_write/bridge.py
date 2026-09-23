"""AUTHORED bounded process bridge; the seven-file SDK adapter is unchanged."""
from dataclasses import asdict
import json
import hashlib
import importlib.metadata
import os
from pathlib import Path
import sys

# -I excludes cwd/user packages. Import only this materialized package root.
sys.path.insert(0, str(Path(__file__).resolve().parent.parent))
from meta_page_write.adapter import (
    ApprovedIntent, FacebookPageWriteAdapter, Scope, UnknownDelivery, create_api,
)


def preflight():
    # Source bytes already compared to the fixed official wheel by G0-G8.
    locked = {
        "facebook_business/adobjects/page.py": "0d08e317fe078de5ad2805ab5501423fe51fea39ef90ae4a0160122a5914ee51",
        "facebook_business/adobjects/pagepost.py": "b0b2009a0667ad57914423284877b104068dca1af4bdfb7f6543b0bec41f0b16",
        "facebook_business/api.py": "2f2e537354859758faba0d89a910b2212b478e0b2aaee771bb1c239abf4c2c72",
        "facebook_business/session.py": "72dd2e43efacc27728222d6c176202f6ff8919a85ba19811ac71fa357bc01587",
        "facebook_business/adobjects/objectparser.py": "9634dee1363b3930469b74942a5b5f58a53bfef5f21bd91519a2e9a25f1de9be",
    }
    distribution = importlib.metadata.distribution("facebook-business")
    for path, expected in locked.items():
        if hashlib.sha256(distribution.locate_file(path).read_bytes()).hexdigest() != expected:
            raise ValueError("installed SDK source differs")
    root = Path(__file__).resolve().parent
    if hashlib.sha256((root / "adapter.py").read_bytes()).hexdigest() != "a397d84a5f92fc8d9d3a44850ad568875da837e59800aeab324fcb5adc6a9e1e":
        raise ValueError("materialized adapter differs")
    requirements = (root / "requirements-windows-py314.lock").read_bytes()
    if hashlib.sha256(requirements).hexdigest() != "a4b4b65d2a19be4bb51892bcb25a4e8127f7b2c16b8f49a2f40b319797314e19":
        raise ValueError("dependency lock differs")
    for line in requirements.decode("utf-8").splitlines():
        if not line or line.startswith("#"):
            continue
        name, url = line.split(" @ ", 1)
        version = url.rsplit("/", 1)[1].split("-")[1]
        if importlib.metadata.version(name) != version:
            raise ValueError("dependency version differs")
    create_api(os.environ.get("META_APP_ID"), os.environ.get("META_APP_SECRET"), os.environ.get("META_PAGE_ACCESS_TOKEN"))


def main(api_factory=None):
    try:
        raw = sys.stdin.buffer.read(32769)
        if not raw or len(raw) > 32768:
            raise ValueError("protocol bounds")
        value = json.loads(raw.decode("utf-8"))
        if set(value) != {"request", "reconcile"} or type(value["reconcile"]) is not bool:
            raise ValueError("protocol fields")
        request = value["request"]
        if set(request) != {"intent", "scheduled_at", "expires_at", "original_approval_id", "provider_reference"}:
            raise ValueError("request fields")
        intent = ApprovedIntent(**request["intent"])
        factory = api_factory or (lambda: create_api(os.environ.get("META_APP_ID"), os.environ.get("META_APP_SECRET"), os.environ.get("META_PAGE_ACCESS_TOKEN")))
        adapter = FacebookPageWriteAdapter(Scope(intent.tenant_id, intent.page_id, intent.profile_sha256), factory())
        if value["reconcile"]:
            if intent.operation != "publish":
                # A missing/permission-denied post cannot prove a DELETE. Never
                # turn absence into success or issue another DELETE implicitly.
                raise UnknownDelivery("REVOKE_REQUIRES_PROVIDER_PROOF", request["provider_reference"])
            receipt = adapter.reconcile(intent, request["provider_reference"])
        elif intent.operation == "publish":
            receipt = adapter.publish_once(intent)
        else:
            receipt = adapter.revoke_once(intent, request["provider_reference"])
        result = {"receipt": asdict(receipt)}
    except UnknownDelivery as error:
        result = {"unknown": True, "code": error.code}
        if error.provider_reference:
            result["provider_reference"] = error.provider_reference
    except Exception:
        # The caller conservatively fences any process/protocol exception.
        # SDK exceptions and secrets are never printed.
        result = {"unknown": True, "code": "PROCESS_BOUNDARY_UNCONFIRMED"}
    sys.stdout.write(json.dumps(result, separators=(",", ":"), ensure_ascii=True))


if __name__ == "__main__":
    if sys.argv[1:] == ["--preflight"]:
        try:
            preflight()
            sys.stdout.write('{"state":"PASS"}')
        except Exception:
            sys.stdout.write('{"state":"FAIL"}')
            sys.exit(2)
    elif len(sys.argv) == 1:
        main()
    else:
        sys.exit(2)
