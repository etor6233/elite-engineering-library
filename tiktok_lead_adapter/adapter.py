"""Strict Lead API wrapper over TikTok's official generic ApiClient transport."""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Dict, List, Mapping, Optional, Sequence, Tuple
from urllib.parse import urlparse

from business_api_client.api_client import ApiClient


LEAD_GET_PATH = "/open_api/v1.3/lead/get/"
SUBSCRIBE_PATH = "/open_api/v1.3/subscription/subscribe/"
LEAD_SOURCES = frozenset({"INSTANT_FORM", "DIRECT_MESSAGE"})


class TikTokLeadConfigurationError(ValueError):
    pass


class TikTokLeadProviderError(RuntimeError):
    pass


def _require_nonempty(value: object, name: str, maximum: int = 512) -> str:
    if not isinstance(value, str) or not value.strip() or len(value) > maximum:
        raise TikTokLeadConfigurationError(f"{name} must be a non-empty bounded string")
    return value.strip()


def _require_identifier(value: object, name: str) -> str:
    result = _require_nonempty(value, name, 128)
    if not result.isdigit():
        raise TikTokLeadConfigurationError(f"{name} must contain decimal digits only")
    return result


def validate_profile(profile: Mapping[str, Any]) -> None:
    required_true = (
        "business_account_proven",
        "developer_app_proven",
        "lead_access_proven",
        "terms_and_consent_approved",
        "test_account_proven",
        "callback_control_and_tls_proven",
        "provider_delivery_test_proven",
        "durable_idempotency_proven",
        "retrieval_reconciliation_proven",
        "retention_and_deletion_approved",
        "quota_and_cost_approved",
    )
    if profile.get("decision") != "PROVEN":
        raise PermissionError("TikTok Lead profile decision is not PROVEN")
    missing = [name for name in required_true if profile.get(name) is not True]
    if missing:
        raise PermissionError("TikTok Lead profile lacks required proof")
    if profile.get("automatic_business_write") is not False:
        raise PermissionError("This adapter does not authorize automatic business writes")


def _selectors(
    lead_source: str,
    advertiser_id: Optional[str],
    library_id: Optional[str],
    page_id: Optional[str],
) -> List[Tuple[str, str]]:
    if lead_source not in LEAD_SOURCES:
        raise TikTokLeadConfigurationError("lead_source is not admitted")
    if bool(advertiser_id) == bool(library_id):
        raise TikTokLeadConfigurationError("exactly one of advertiser_id or library_id is required")
    pairs: List[Tuple[str, str]] = [("lead_source", lead_source)]
    if advertiser_id:
        pairs.append(("advertiser_id", _require_identifier(advertiser_id, "advertiser_id")))
    if library_id:
        pairs.append(("library_id", _require_identifier(library_id, "library_id")))
    if lead_source == "INSTANT_FORM":
        pairs.append(("page_id", _require_identifier(page_id, "page_id")))
    elif page_id is not None:
        raise TikTokLeadConfigurationError("page_id is unsupported for DIRECT_MESSAGE")
    return pairs


def _call(api_client: ApiClient, path: str, method: str, query: Sequence[Tuple[str, str]], headers: Dict[str, str], body: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
    try:
        result = api_client.call_api(
            path,
            method,
            {},
            list(query),
            headers,
            body=body,
            post_params=[],
            files={},
            response_type="InlineResponse200",
            auth_settings=[],
            async_req=False,
            _return_http_data_only=True,
            _preload_content=True,
            _request_timeout=(5, 30),
            collection_formats={},
        )
    except Exception:
        raise TikTokLeadProviderError("TikTok Lead API request failed") from None
    if not isinstance(result, dict) or not isinstance(result.get("data"), dict):
        raise TikTokLeadProviderError("TikTok Lead API returned an invalid response")
    request_id = result.get("request_id")
    if not isinstance(request_id, str) or not request_id:
        raise TikTokLeadProviderError("TikTok Lead API response lacks request_id")
    return result


@dataclass(frozen=True)
class TikTokLeadClient:
    api_client: ApiClient
    profile: Mapping[str, Any]

    def __post_init__(self) -> None:
        validate_profile(self.profile)

    def get_lead(self, access_token: str, lead_source: str, *, advertiser_id: Optional[str] = None, library_id: Optional[str] = None, page_id: Optional[str] = None) -> Dict[str, Any]:
        token = _require_nonempty(access_token, "access_token", 4096)
        query = _selectors(lead_source, advertiser_id, library_id, page_id)
        result = _call(self.api_client, LEAD_GET_PATH, "GET", query, {"Accept": "application/json", "Access-Token": token})
        data = result["data"]
        if not isinstance(data.get("lead_data"), dict) or not isinstance(data.get("meta_data"), dict):
            raise TikTokLeadProviderError("TikTok Lead response lacks lead_data or meta_data")
        return result

    def create_subscription(self, access_token: str, app_id: str, app_secret: str, callback_url: str, lead_source: str, *, advertiser_id: Optional[str] = None, library_id: Optional[str] = None, page_id: Optional[str] = None) -> Dict[str, Any]:
        token = _require_nonempty(access_token, "access_token", 4096)
        app = _require_nonempty(app_id, "app_id", 256)
        secret = _require_nonempty(app_secret, "app_secret", 4096)
        callback = _require_nonempty(callback_url, "callback_url", 2048)
        parsed = urlparse(callback)
        if parsed.scheme != "https" or not parsed.hostname or parsed.username or parsed.password or parsed.fragment:
            raise TikTokLeadConfigurationError("callback_url must be a credential-free HTTPS URL")
        selectors = dict(_selectors(lead_source, advertiser_id, library_id, page_id))
        detail: Dict[str, Any] = {"access_token": token, "lead_source": lead_source}
        detail.update({key: value for key, value in selectors.items() if key != "lead_source"})
        body = {
            "app_id": app,
            "secret": secret,
            "subscribe_entity": "LEAD",
            "callback_url": callback,
            "subscription_detail": detail,
        }
        result = _call(self.api_client, SUBSCRIBE_PATH, "POST", (), {"Accept": "application/json", "Content-Type": "application/json", "Access-Token": token}, body)
        subscription_id = result["data"].get("subscription_id")
        if not isinstance(subscription_id, str) or not subscription_id:
            raise TikTokLeadProviderError("TikTok subscription response lacks subscription_id")
        return result
