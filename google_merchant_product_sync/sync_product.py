from __future__ import annotations

import argparse
from datetime import datetime, timezone
import hashlib
import importlib.metadata
import json
from pathlib import Path
import re
import shutil
import sys
import time
import uuid
from typing import Any, Callable, Protocol
from urllib.parse import urlparse

from google.api_core.exceptions import NotFound
from google.protobuf.json_format import MessageToDict
from google.shopping import merchant_products_v1 as merchant
from google.shopping.type import Price


SDK_VERSION = "1.8.0"
PROFILE_SCHEMA = "elite-google-merchant-product-sync-profile/v1"
AVAILABILITY = {
    "IN_STOCK": merchant.Availability.IN_STOCK,
    "OUT_OF_STOCK": merchant.Availability.OUT_OF_STOCK,
    "PREORDER": merchant.Availability.PREORDER,
    "LIMITED_AVAILABILITY": merchant.Availability.LIMITED_AVAILABILITY,
    "BACKORDER": merchant.Availability.BACKORDER,
}
CONDITION = {
    "NEW": merchant.Condition.NEW,
    "USED": merchant.Condition.USED,
    "REFURBISHED": merchant.Condition.REFURBISHED,
}
PRODUCT_KEYS = {
    "offer_id", "content_language", "feed_label", "title", "description",
    "link", "image_link", "availability", "condition", "price", "brand", "gtins",
}


class ProductInputsClient(Protocol):
    def insert_product_input(self, request: Any, *, retry: Any, timeout: float) -> Any: ...


class ProductsClient(Protocol):
    def get_product(self, request: Any, *, retry: Any, timeout: float) -> Any: ...


def sha256(value: bytes) -> str:
    return hashlib.sha256(value).hexdigest()


def json_bytes(value: Any) -> bytes:
    return (json.dumps(value, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")


def message_dict(value: Any) -> dict[str, Any]:
    if isinstance(value, dict):
        return value
    protobuf = getattr(value, "_pb", None)
    if protobuf is not None:
        return MessageToDict(protobuf, preserving_proto_field_name=True)
    to_dict = getattr(value, "to_dict", None)
    if callable(to_dict):
        result = to_dict()
        if isinstance(result, dict):
            return result
    raise TypeError(f"unsupported provider message type: {type(value).__name__}")


def write_json_atomic(path: Path, value: Any) -> bytes:
    data = json_bytes(value)
    temporary = path.with_name(f".{path.name}.{uuid.uuid4().hex}.tmp")
    temporary.write_bytes(data)
    temporary.replace(path)
    return data


def valid_https_url(value: Any) -> bool:
    if not isinstance(value, str) or len(value) > 2048:
        return False
    parsed = urlparse(value)
    return parsed.scheme == "https" and bool(parsed.hostname) and parsed.username is None and parsed.password is None


def read_approved_profile(path: Path) -> dict[str, Any]:
    profile = json.loads(path.read_text(encoding="utf-8"))
    if profile.get("schema") != PROFILE_SCHEMA:
        raise ValueError("unsupported profile schema")
    if profile.get("sdk") != f"google-shopping-merchant-products=={SDK_VERSION}":
        raise ValueError("profile SDK version mismatch")
    if profile.get("decision") not in {"APPROVED_FOR_TEST_ACCOUNT", "APPROVED_FOR_PRODUCTION"}:
        raise PermissionError("Google Merchant external write decision is not approved")
    account = str(profile.get("merchant_account_id", ""))
    if not re.fullmatch(r"[0-9]{3,20}", account):
        raise ValueError("merchant_account_id must contain 3..20 digits")
    expected_source = rf"accounts/{account}/dataSources/[0-9]{{1,20}}"
    if not re.fullmatch(expected_source, str(profile.get("data_source_name", ""))):
        raise ValueError("data_source_name must belong to the approved Merchant account")
    if profile.get("data_source_type") != "API_PRIMARY_PROVEN":
        raise PermissionError("an API primary data source must be proven")
    for key in (
        "merchant_api_enabled", "billing_approved", "account_access_proven",
        "test_account_contract_proven", "external_product_write_approved", "quota_and_cost_approved",
    ):
        if profile.get(key) is not True:
            raise PermissionError(f"{key} must be proven")
    if not str(profile.get("cloud_project_reference", "")).strip():
        raise ValueError("cloud_project_reference is required")
    if not str(profile.get("adc_identity_reference", "")).strip():
        raise ValueError("adc_identity_reference is required")
    if not str(profile.get("reconciliation_owner", "")).strip() or not str(profile.get("evidence_retention", "")).strip():
        raise ValueError("reconciliation_owner and evidence_retention are required")
    if profile.get("automatic_local_business_write") is not False:
        raise ValueError("automatic local business writes are forbidden")
    cadence = profile.get("refresh_cadence_days")
    if not isinstance(cadence, int) or isinstance(cadence, bool) or cadence < 1 or cadence > 30:
        raise ValueError("refresh_cadence_days must be within 1..30")
    approved = profile.get("approved_offer_ids")
    if not isinstance(approved, list) or not approved or any(not isinstance(value, str) for value in approved):
        raise ValueError("approved_offer_ids must be a non-empty string array")
    return profile


def build_insert_request(profile: dict[str, Any], product: dict[str, Any]) -> merchant.InsertProductInputRequest:
    if set(product) != PRODUCT_KEYS:
        raise ValueError("product keys must match the exact supported schema")
    offer_id = product["offer_id"]
    if not isinstance(offer_id, str) or not re.fullmatch(r"[A-Za-z0-9._-]{1,50}", offer_id):
        raise ValueError("offer_id must be a safe 1..50 character identifier")
    if offer_id not in profile["approved_offer_ids"]:
        raise PermissionError("offer_id is not approved by the profile")
    language = product["content_language"]
    label = product["feed_label"]
    if not isinstance(language, str) or not re.fullmatch(r"[a-z]{2}", language):
        raise ValueError("content_language must be a two-letter lowercase code")
    if not isinstance(label, str) or not re.fullmatch(r"[A-Z0-9-]{2,20}", label):
        raise ValueError("feed_label must be 2..20 uppercase alphanumeric/hyphen characters")
    title = product["title"]
    description = product["description"]
    if not isinstance(title, str) or not title.strip() or len(title) > 150:
        raise ValueError("title must contain 1..150 characters")
    if not isinstance(description, str) or not description.strip() or len(description) > 5000:
        raise ValueError("description must contain 1..5000 characters")
    if not valid_https_url(product["link"]) or not valid_https_url(product["image_link"]):
        raise ValueError("link and image_link must be credential-free HTTPS URLs")
    if product["availability"] not in AVAILABILITY or product["condition"] not in CONDITION:
        raise ValueError("availability or condition is unsupported")
    price = product["price"]
    if not isinstance(price, dict) or set(price) != {"amount_micros", "currency_code"}:
        raise ValueError("price must contain amount_micros and currency_code")
    amount = price["amount_micros"]
    currency = price["currency_code"]
    if not isinstance(amount, int) or isinstance(amount, bool) or amount <= 0 or amount > 10**18:
        raise ValueError("amount_micros must be a positive bounded integer")
    if not isinstance(currency, str) or not re.fullmatch(r"[A-Z]{3}", currency):
        raise ValueError("currency_code must contain three uppercase letters")
    brand = product["brand"]
    if not isinstance(brand, str) or len(brand) > 70:
        raise ValueError("brand must contain at most 70 characters")
    gtins = product["gtins"]
    if not isinstance(gtins, list) or len(gtins) > 10 or any(not isinstance(value, str) or not re.fullmatch(r"(?:[0-9]{8}|[0-9]{12,14})", value) for value in gtins):
        raise ValueError("gtins must contain at most 10 numeric GTIN-8/12/13/14 values")

    attributes = merchant.ProductAttributes(
        title=title.strip(),
        description=description.strip(),
        link=product["link"],
        image_link=product["image_link"],
        availability=AVAILABILITY[product["availability"]],
        condition=CONDITION[product["condition"]],
        price=Price(amount_micros=amount, currency_code=currency),
        brand=brand or None,
        gtins=gtins,
    )
    product_input = merchant.ProductInput(
        offer_id=offer_id,
        content_language=language,
        feed_label=label,
        product_attributes=attributes,
    )
    return merchant.InsertProductInputRequest(
        parent=f"accounts/{profile['merchant_account_id']}",
        data_source=profile["data_source_name"],
        product_input=product_input,
    )


def sync_product(
    inputs_client: ProductInputsClient,
    products_client: ProductsClient,
    request: merchant.InsertProductInputRequest,
    output_directory: Path,
    *,
    timeout_seconds: float = 30.0,
    status_attempts: int = 3,
    status_interval_seconds: float = 1.0,
    sleep_fn: Callable[[float], None] = time.sleep,
) -> dict[str, Any]:
    if importlib.metadata.version("google-shopping-merchant-products") != SDK_VERSION:
        raise RuntimeError(f"google-shopping-merchant-products {SDK_VERSION} is required")
    if timeout_seconds <= 0 or timeout_seconds > 120:
        raise ValueError("timeout must be within 1..120 seconds")
    if status_attempts < 1 or status_attempts > 12 or status_interval_seconds < 0 or status_interval_seconds > 60:
        raise ValueError("status polling bounds are invalid")
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".google-merchant-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    request_payload = message_dict(request)
    request_data = json_bytes(request_payload)
    (stage / "insert-request.json").write_bytes(request_data)
    try:
        inserted = inputs_client.insert_product_input(request, retry=None, timeout=timeout_seconds)
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise
    inserted_payload = message_dict(inserted)
    product_name = getattr(inserted, "product", None)
    if not isinstance(product_name, str) or not re.fullmatch(r"accounts/[0-9]{3,20}/products/.+", product_name):
        shutil.rmtree(stage, ignore_errors=True)
        raise RuntimeError("provider insert response omitted the canonical product resource name")
    insert_data = json_bytes(inserted_payload)
    (stage / "insert-response.json").write_bytes(insert_data)
    offer_id = request.product_input.offer_id
    receipt: dict[str, Any] = {
        "schema": "elite-google-merchant-product-sync-receipt/v1",
        "created_at": datetime.now(timezone.utc).isoformat(),
        "provider": "Google Merchant API",
        "sdk_version": SDK_VERSION,
        "api_version": "products/v1",
        "operation": "productInputs.insert + products.get",
        "processing_state": "INSERT_ACCEPTED_STATUS_UNCHECKED",
        "offer_id_sha256": sha256(offer_id.encode("utf-8")),
        "product_name_sha256": sha256(product_name.encode("utf-8")),
        "request_sha256": sha256(request_data),
        "insert_response_sha256": sha256(insert_data),
        "external_product_write": True,
        "automatic_local_business_write": False,
        "approval_inferred": False,
    }
    (stage / "SYNC_RECEIPT.json").write_bytes(json_bytes(receipt))
    stage.replace(output)

    request_get = merchant.GetProductRequest(name=product_name)
    for attempt in range(1, status_attempts + 1):
        try:
            product = products_client.get_product(request_get, retry=None, timeout=timeout_seconds)
            product_payload = message_dict(product)
            product_data = write_json_atomic(output / "product-response.json", product_payload)
            receipt["processing_state"] = "PROCESSED_PRODUCT_STATUS_OBSERVED"
            receipt["status_attempts"] = attempt
            receipt["product_response_sha256"] = sha256(product_data)
            write_json_atomic(output / "SYNC_RECEIPT.json", receipt)
            return receipt
        except NotFound:
            if attempt < status_attempts:
                sleep_fn(status_interval_seconds)
                continue
            receipt["processing_state"] = "INSERT_ACCEPTED_PRODUCT_NOT_YET_VISIBLE"
            receipt["status_attempts"] = attempt
            write_json_atomic(output / "SYNC_RECEIPT.json", receipt)
            return receipt
        except BaseException as error:
            receipt["processing_state"] = "INSERT_ACCEPTED_STATUS_CHECK_FAILED"
            receipt["status_attempts"] = attempt
            receipt["status_error_type"] = type(error).__name__
            write_json_atomic(output / "SYNC_RECEIPT.json", receipt)
            raise
    raise AssertionError("unreachable status loop")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--product", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--timeout-seconds", type=float, default=30.0)
    args = parser.parse_args()
    profile = read_approved_profile(args.profile)
    product = json.loads(args.product.read_text(encoding="utf-8"))
    request = build_insert_request(profile, product)
    inputs_client = merchant.ProductInputsServiceClient()
    products_client = merchant.ProductsServiceClient()
    receipt = sync_product(inputs_client, products_client, request, args.output, timeout_seconds=args.timeout_seconds)
    print(json.dumps({"status": "GOOGLE_MERCHANT_PRODUCT_SYNC_PASS", "receipt": receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
