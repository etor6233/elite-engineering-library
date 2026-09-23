"""Parser for TikTok Lead webhook signals; it never authenticates or persists them."""

from __future__ import annotations

import hashlib
import json
from dataclasses import dataclass
from typing import Any, Dict, List, Mapping, Optional, Tuple


MAX_WEBHOOK_BYTES = 1_048_576
MAX_ENTRIES = 100
MAX_CHANGES = 200


class TikTokWebhookError(ValueError):
    pass


def _bounded_string(value: object, name: str, maximum: int = 1024) -> str:
    if not isinstance(value, (str, int)):
        raise TikTokWebhookError(f"{name} must be a string-compatible identifier")
    result = str(value).strip()
    if not result or len(result) > maximum:
        raise TikTokWebhookError(f"{name} is empty or too long")
    return result


@dataclass(frozen=True)
class LeadWebhookEntry:
    lead_id: str
    page_id: Optional[str]
    create_time: int
    campaign_id: Optional[str]
    adgroup_id: Optional[str]
    ad_id: Optional[str]
    changes: Tuple[Tuple[str, Any], ...]


@dataclass(frozen=True)
class LeadWebhookSignal:
    request_id: Optional[str]
    notification_time: int
    raw_sha256: str
    delivery_key: str
    trust_state: str
    entries: Tuple[LeadWebhookEntry, ...]

    def redacted_receipt(self) -> Dict[str, Any]:
        return {
            "schema": "elite-tiktok-lead-webhook-receipt/v1",
            "trust_state": self.trust_state,
            "request_id_present": self.request_id is not None,
            "notification_time": self.notification_time,
            "raw_sha256": self.raw_sha256,
            "delivery_key": self.delivery_key,
            "entry_count": len(self.entries),
            "automatic_persistence_allowed": False,
        }


def _optional_identifier(entry: Mapping[str, Any], name: str) -> Optional[str]:
    value = entry.get(name)
    return None if value is None else _bounded_string(value, name, 256)


def parse_lead_webhook(raw_body: bytes) -> LeadWebhookSignal:
    if not isinstance(raw_body, bytes) or not raw_body or len(raw_body) > MAX_WEBHOOK_BYTES:
        raise TikTokWebhookError("webhook body is empty, non-bytes or too large")
    raw_sha = hashlib.sha256(raw_body).hexdigest()
    try:
        payload = json.loads(raw_body.decode("utf-8"))
    except (UnicodeDecodeError, json.JSONDecodeError):
        raise TikTokWebhookError("webhook body is not valid UTF-8 JSON") from None
    if not isinstance(payload, dict) or payload.get("object") != 1:
        raise TikTokWebhookError("webhook object is not TikTok Lead")
    notification_time = payload.get("time")
    if not isinstance(notification_time, int) or notification_time < 0:
        raise TikTokWebhookError("webhook time must be a non-negative integer")
    raw_entries = payload.get("entry")
    if not isinstance(raw_entries, list) or not raw_entries or len(raw_entries) > MAX_ENTRIES:
        raise TikTokWebhookError("webhook entry list is empty or exceeds the safety bound")
    entries: List[LeadWebhookEntry] = []
    for raw_entry in raw_entries:
        if not isinstance(raw_entry, dict):
            raise TikTokWebhookError("webhook entry must be an object")
        lead_id = _bounded_string(raw_entry.get("id"), "id", 256)
        create_time = raw_entry.get("create_time")
        if not isinstance(create_time, int) or create_time < 0:
            raise TikTokWebhookError("entry create_time must be a non-negative integer")
        raw_changes = raw_entry.get("changes")
        if not isinstance(raw_changes, list) or not raw_changes or len(raw_changes) > MAX_CHANGES:
            raise TikTokWebhookError("entry changes are empty or exceed the safety bound")
        changes: List[Tuple[str, Any]] = []
        seen = set()
        for change in raw_changes:
            if not isinstance(change, dict):
                raise TikTokWebhookError("change must be an object")
            field = _bounded_string(change.get("field"), "field", 256)
            if field in seen:
                raise TikTokWebhookError("duplicate change field")
            seen.add(field)
            changes.append((field, change.get("value")))
        entries.append(LeadWebhookEntry(lead_id, _optional_identifier(raw_entry, "page_id"), create_time, _optional_identifier(raw_entry, "campaign_id"), _optional_identifier(raw_entry, "adgroup_id"), _optional_identifier(raw_entry, "ad_id"), tuple(changes)))
    request_value = payload.get("request_id")
    request_id = None if request_value is None else _bounded_string(request_value, "request_id", 256)
    stable = json.dumps({"request_id": request_id, "time": notification_time, "entry_ids": [[entry.lead_id, entry.page_id, entry.create_time] for entry in entries]}, ensure_ascii=False, separators=(",", ":"), sort_keys=True).encode("utf-8")
    delivery_key = hashlib.sha256(stable).hexdigest()
    return LeadWebhookSignal(request_id, notification_time, raw_sha, delivery_key, "UNAUTHENTICATED_PROVIDER_SIGNAL", tuple(entries))
