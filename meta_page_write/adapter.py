"""AUTHORED binding glue over the locked official Meta Business SDK.

This leaf does not authorize content, claim a durable delivery lease, schedule
posts or retry writes. Its caller must own those decisions before invoking it.
"""
from dataclasses import dataclass
import hashlib
import importlib.metadata
import json
import re

from facebook_business.adobjects.page import Page
from facebook_business.adobjects.pagepost import PagePost
from facebook_business.api import FacebookAdsApi
from facebook_business.session import FacebookSession


class InvalidBinding(ValueError):
    pass


class UnknownDelivery(RuntimeError):
    """No automatic retry is safe; retain the optional reference for GET."""

    def __init__(self, code, provider_reference=None):
        super().__init__(code)
        self.code = code
        self.provider_reference = provider_reference


def _hash(value):
    return isinstance(value, str) and re.fullmatch(r"[a-f0-9]{64}", value) is not None


def _identifier(value):
    return isinstance(value, str) and re.fullmatch(r"[A-Za-z0-9_-]{1,128}", value) is not None


def _page(value):
    return isinstance(value, str) and re.fullmatch(r"[1-9][0-9]{0,31}", value) is not None


@dataclass(frozen=True)
class Scope:
    tenant_id: str
    page_id: str
    profile_sha256: str

    def validate(self):
        if not _identifier(self.tenant_id) or not _page(self.page_id) or not _hash(self.profile_sha256):
            raise InvalidBinding("invalid configured page scope")


@dataclass(frozen=True)
class ApprovedIntent:
    tenant_id: str
    page_id: str
    profile_sha256: str
    approval_id: str
    delivery_key: str
    operation: str
    message: str
    content_sha256: str


@dataclass(frozen=True)
class Receipt:
    tenant_id: str
    page_id: str
    profile_sha256: str
    approval_id: str
    delivery_key: str
    provider_reference: str
    state: str
    content_sha256: str
    evidence_sha256: str


def create_api(app_id, app_secret, page_access_token):
    """Credentials stay in memory; no global default API or debug logger."""
    if importlib.metadata.version("facebook-business") != "26.0.1":
        raise InvalidBinding("SDK version differs from admitted lock")
    if not _page(app_id) or not isinstance(app_secret, str) or not app_secret or not isinstance(page_access_token, str) or not page_access_token:
        raise InvalidBinding("credentials are unavailable")
    session = FacebookSession(app_id, app_secret, page_access_token, timeout=10)
    session.requests.trust_env = False
    return FacebookAdsApi(session, api_version="v26.0", enable_debug_logger=False)


class FacebookPageWriteAdapter:
    def __init__(self, scope, api):
        if not isinstance(scope, Scope) or not isinstance(api, FacebookAdsApi):
            raise InvalidBinding("scope and official SDK API required")
        scope.validate()
        # This attribute is pinned by the exact official api.py source lock.
        if api._api_version != "v26.0" or api._enable_debug_logger:
            raise InvalidBinding("Graph API version differs from admitted lock")
        self.scope, self.api = scope, api

    def _validate(self, intent, operation):
        if not isinstance(intent, ApprovedIntent):
            raise InvalidBinding("bound intent required")
        if (intent.tenant_id, intent.page_id, intent.profile_sha256) != (self.scope.tenant_id, self.scope.page_id, self.scope.profile_sha256):
            raise InvalidBinding("intent does not match configured scope")
        if intent.operation != operation or not _identifier(intent.approval_id) or not _identifier(intent.delivery_key) or not 16 <= len(intent.delivery_key) <= 128:
            raise InvalidBinding("intent approval, action or delivery binding invalid")
        try:
            body = intent.message.encode("utf-8")
        except (AttributeError, UnicodeError):
            raise InvalidBinding("invalid UTF-8 content") from None
        # Local bounded-transport limit; not a claim about Meta's character limit.
        if not body.strip() or len(body) > 16384 or not _hash(intent.content_sha256) or hashlib.sha256(body).hexdigest() != intent.content_sha256:
            raise InvalidBinding("content differs from the approved bytes")

    def _reference(self, provider_reference):
        if not isinstance(provider_reference, str) or re.fullmatch(re.escape(self.scope.page_id) + r"_[1-9][0-9]{0,63}", provider_reference) is None:
            raise InvalidBinding("post reference does not belong to configured page")

    def _observe(self, intent, provider_reference):
        self._reference(provider_reference)
        try:
            post = PagePost(provider_reference, api=self.api).api_get(fields=["id", "from", "message", "is_published"])
            value = post.export_all_data()
            if value.get("id") != provider_reference or value.get("from", {}).get("id") != self.scope.page_id or value.get("message") != intent.message or value.get("is_published") is not True:
                raise ValueError("post binding differs")
        except Exception:
            # SDK exceptions can include request details; keep only a safe code.
            raise UnknownDelivery("POST_OBSERVATION_UNCONFIRMED", provider_reference) from None
        return {"id": provider_reference, "from": {"id": self.scope.page_id}, "message": intent.message, "is_published": True}

    def _receipt(self, intent, reference, state, observed):
        evidence = hashlib.sha256(json.dumps(observed, sort_keys=True, separators=(",", ":"), ensure_ascii=False).encode("utf-8")).hexdigest()
        return Receipt(intent.tenant_id, intent.page_id, intent.profile_sha256, intent.approval_id, intent.delivery_key, reference, state, intent.content_sha256, evidence)

    def publish_once(self, intent):
        """Caller must hold its durable outbound fence. One POST, then one GET."""
        self._validate(intent, "publish")
        reference = None
        try:
            result = Page(self.scope.page_id, api=self.api).create_feed(params={"message": intent.message, "published": True})
            reference = result.get("id")
            self._reference(reference)
        except Exception:
            raise UnknownDelivery("PUBLISH_OUTCOME_UNKNOWN") from None
        observed = self._observe(intent, reference)
        return self._receipt(intent, reference, "published", observed)

    def reconcile(self, intent, provider_reference):
        """GET only; requires a known post ID. Never republishes an unknown POST."""
        self._validate(intent, "publish")
        observed = self._observe(intent, provider_reference)
        return self._receipt(intent, provider_reference, "published", observed)

    def revoke_once(self, intent, provider_reference):
        """Separate caller-approved revoke; verifies page and original bytes first."""
        self._validate(intent, "revoke")
        observed = self._observe(intent, provider_reference)
        try:
            # PagePost.api_delete's ObjectParser discards the success field.
            # Keep the official SDK's raw response for this same node DELETE.
            result = self.api.call("DELETE", (provider_reference,)).json()
            if not isinstance(result, dict) or result.get("success") is not True:
                raise ValueError("delete not confirmed")
        except Exception:
            raise UnknownDelivery("REVOKE_OUTCOME_UNKNOWN", provider_reference) from None
        return self._receipt(intent, provider_reference, "revoked", {"observed_before_delete": observed, "delete_success": True})
