"""AUTHORED bounded IPC around the exact admitted Google models/client.

Candidate implementation. The original build_insert_request stays unchanged.
No provider response, secret or traceback is printed outside the typed result.
"""
from __future__ import annotations
import hashlib
import importlib.util
import importlib.metadata
import importlib.machinery
import json
from pathlib import Path
from types import ModuleType
import re
import sys
from typing import Any
from urllib.parse import urlsplit

LIMIT = 32768
SDK_VERSION = "1.8.0"


def admit_installed_sources(expected: str) -> None:
    """Same verified-byte principle as the existing isolated source launcher.

    -I -S prevents site/.pth startup. Source imports compile the verified bytes;
    injected package initializers or bytecode cannot replace those bytes.
    """
    if not sys.flags.isolated or not sys.flags.no_site or sys.version_info[:2] != (3, 14) or sys.platform != "win32":
        raise ValueError("runtime lane")
    path = Path(__file__).resolve().with_name("connected-runtime.lock.json")
    raw = path.read_bytes()
    if len(raw) > 1048576 or not re.fullmatch("[0-9a-f]{64}", expected) or hashlib.sha256(raw).hexdigest() != expected:
        raise ValueError("runtime source lock")
    lock = json.loads(raw)
    if lock["schema"] != "elite-merchant-installed-source-lock/v1" or lock["sdk"] != SDK_VERSION:
        raise ValueError("runtime lock schema")
    site = (Path(sys.executable).resolve().parent.parent / "Lib" / "site-packages").resolve()
    sources = {}
    for relative, digest in lock["files"].items():
        target = (site / relative).resolve()
        if not target.is_relative_to(site) or target.is_symlink():
            raise ValueError("installed source path")
        data = target.read_bytes()
        if hashlib.sha256(data).hexdigest() != digest:
            raise ValueError("installed source mismatch")
        sources[str(target)] = data

    class VerifiedSource(importlib.machinery.SourceFileLoader):
        def get_code(self, fullname):
            data = sources[str(Path(self.path).resolve())]
            return compile(data, self.path, "exec", dont_inherit=True)

    class VerifiedFinder:
        @staticmethod
        def find_spec(fullname, path=None, target=None):
            spec = importlib.machinery.PathFinder.find_spec(fullname, path, target)
            if spec is None or spec.origin in (None, "built-in", "frozen"):
                return None
            origin = Path(spec.origin).resolve()
            if not origin.is_relative_to(site):
                return None
            if str(origin) not in sources:
                raise ImportError("unlocked installed module")
            if isinstance(spec.loader, importlib.machinery.SourceFileLoader):
                spec.loader = VerifiedSource(fullname, str(origin))
            return spec

    sys.meta_path.insert(0, VerifiedFinder())
    sys.path.append(str(site))


def strict_object(pairs: list[tuple[str, Any]]) -> dict[str, Any]:
    result = {}
    seen = set()
    for key, value in pairs:
        folded = key.casefold()
        if folded in seen:
            raise ValueError("duplicate key")
        seen.add(folded)
        result[key] = value
    return result


def decode(raw: bytes) -> dict[str, Any]:
    if not 0 < len(raw) <= LIMIT:
        raise ValueError("input bound")
    value = json.loads(raw.decode("utf-8"), object_pairs_hook=strict_object,
                       parse_constant=lambda _: (_ for _ in ()).throw(ValueError("nonfinite")))
    if not isinstance(value, dict):
        raise ValueError("object required")
    return value


def load_owner(expected: str):
    path = Path(__file__).resolve().with_name("sync_product.py")
    raw = path.read_bytes()
    if len(raw) > 65536 or not re.fullmatch("[0-9a-f]{64}", expected) or hashlib.sha256(raw).hexdigest() != expected:
        raise ValueError("source lock")
    module = ModuleType("elite_locked_merchant_owner")
    module.__file__ = str(path)
    sys.modules[module.__name__] = module
    exec(compile(raw, str(path), "exec", dont_inherit=True), module.__dict__)
    if importlib.metadata.version("google-shopping-merchant-products") != SDK_VERSION:
        raise ValueError("SDK pin")
    return module


def request_from_intent(owner, intent: dict[str, Any]):
    product = intent["product"]
    account, datasource = intent["account_id"], intent["data_source"]
    generation = intent["generation"]
    if not isinstance(account, str) or not re.fullmatch("[0-9]{3,20}", account):
        raise ValueError("account")
    if not isinstance(datasource, str) or not re.fullmatch("accounts/" + account + "/dataSources/[0-9]{3,20}", datasource):
        raise ValueError("data source")
    if not isinstance(generation, str) or not re.fullmatch("[1-9][0-9]{0,18}", generation) or int(generation) > 9223372036854775807:
        raise ValueError("generation")
    if not isinstance(product, dict) or set(product) != owner.PRODUCT_KEYS:
        raise ValueError("product shape")
    product = dict(product)
    if not isinstance(product["price"], dict) or set(product["price"]) != {"amount_micros", "currency_code"}:
        raise ValueError("price")
    micros = product["price"]["amount_micros"]
    if not isinstance(micros, str) or not re.fullmatch("[1-9][0-9]{0,18}", micros):
        raise ValueError("micros")
    product["price"] = dict(product["price"], amount_micros=int(micros))
    profile = {"merchant_account_id": account, "data_source_name": datasource,
               "approved_offer_ids": [product["offer_id"]]}
    request = owner.build_insert_request(profile, product)
    request.product_input.version_number = int(generation)
    resource = f"accounts/{account}/products/{product['content_language']}~{product['feed_label']}~{product['offer_id']}"
    if intent["resource"] != resource:
        raise ValueError("resource binding")
    return request, resource


def bounded_adapter(fixture_origin: str = ""):
    from requests.adapters import HTTPAdapter

    if fixture_origin:
        target = urlsplit(fixture_origin)
        if target.scheme != "http" or target.hostname != "127.0.0.1" or not target.port or target.path or target.query or target.fragment or target.username:
            raise ValueError("fixture origin")

    class Adapter(HTTPAdapter):
        def send(self, request, **kwargs):
            origin = urlsplit(request.url)
            if origin.scheme != "https" or origin.netloc != "merchantapi.googleapis.com" or not origin.path.startswith("/products/v1/accounts/") or request.method not in {"GET", "POST"}:
                raise ValueError("provider origin/path")
            if fixture_origin:
                request.url = fixture_origin + origin.path + ("?" + origin.query if origin.query else "")
                request.headers["Authorization"] = "Bearer synthetic-merchant-sdk-fixture"
            kwargs.update(stream=True, timeout=8, proxies={})
            response = super().send(request, **kwargs)
            if 300 <= response.status_code < 400:
                response.close()
                raise ValueError("redirect denied")
            body = bytearray()
            try:
                for part in response.iter_content(chunk_size=4096):
                    if len(body) + len(part) > LIMIT:
                        raise ValueError("response bound")
                    body.extend(part)
            finally:
                response.close()
            response._content = bytes(body)
            response._content_consumed = True
            return response
    return Adapter(max_retries=0)


def clients(owner, mode: str, fixture_origin: str):
    if mode not in {"CREDENTIALS", "LOCAL_FIXTURES"} or (mode == "LOCAL_FIXTURES") != bool(fixture_origin):
        raise ValueError("runtime scope")
    kwargs = {"transport": "rest", "client_options": {"api_endpoint": "merchantapi.googleapis.com"}}
    if mode == "LOCAL_FIXTURES":
        from google.auth.credentials import AnonymousCredentials
        kwargs["credentials"] = AnonymousCredentials()
    inputs = owner.merchant.ProductInputsServiceClient(**kwargs)
    products = owner.merchant.ProductsServiceClient(**kwargs)
    for client in (inputs, products):
        # This is a pinned transport wrapper, not a modification to Google source.
        session = client.transport._session
        session.trust_env = False
        session.max_redirects = 0
        session._max_refresh_attempts = 0
        session.mount("https://", bounded_adapter(fixture_origin))
    return inputs, products


def execute(owner, command: dict[str, Any], inputs, products) -> dict[str, Any]:
    from google.api_core import exceptions
    request, resource = request_from_intent(owner, command["intent"])
    operation = command["operation"]
    if operation not in {"INSERT", "GET"}:
        raise ValueError("operation")
    try:
        if operation == "INSERT":
            result = inputs.insert_product_input(request, retry=None, timeout=8)
        else:
            result = products.get_product(owner.merchant.GetProductRequest(name=resource), retry=None, timeout=8)
        value = owner.message_dict(result)
        raw = json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(",", ":")).encode("utf-8")
        if len(raw) > LIMIT:
            return {"state": "UNKNOWN", "operation": operation}
        return {"state": "OBSERVED", "operation": operation, "resource": resource,
                "response": value, "response_sha256": hashlib.sha256(raw).hexdigest()}
    except exceptions.NotFound:
        return {"state": "NOT_FOUND" if operation == "GET" else "UNKNOWN", "operation": operation}
    except (exceptions.BadRequest, exceptions.Unauthorized, exceptions.Forbidden,
            exceptions.Conflict, exceptions.Aborted):
        return {"state": "REJECTED" if operation == "INSERT" else "UNKNOWN", "operation": operation}
    except Exception:
        return {"state": "UNKNOWN", "operation": operation}


def main() -> int:
    result = {"state": "UNKNOWN"}
    try:
        if len(sys.argv) != 3:
            raise ValueError("source hash argument")
        admit_installed_sources(sys.argv[2])
        owner = load_owner(sys.argv[1])
        command = decode(sys.stdin.buffer.read(LIMIT + 1))
        if set(command) != {"operation", "intent", "mode", "fixture_origin"}:
            raise ValueError("command keys")
        # Validate the request before credential lookup or SDK client construction.
        request_from_intent(owner, command["intent"])
        inputs, products = clients(owner, command["mode"], command["fixture_origin"])
        try:
            result = execute(owner, command, inputs, products)
        finally:
            inputs.transport.close()
            products.transport.close()
    except Exception:
        result = {"state": "UNKNOWN"}
    raw = json.dumps(result, ensure_ascii=False, sort_keys=True, separators=(",", ":")).encode("utf-8")
    if len(raw) > LIMIT:
        raw = b'{"state":"UNKNOWN"}'
    sys.stdout.buffer.write(raw)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
