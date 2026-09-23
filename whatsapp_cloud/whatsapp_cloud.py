"""
Copyright (c) Meta Platforms, Inc. and affiliates.
All rights reserved.

Adapted from the exact official examples listed in PROVENANCE.md under the
root upstream/LICENSE. Safety, atomic evidence and validation changes are
Copyright (c) the Elite Engineering Library owner.
"""

from __future__ import annotations

from datetime import datetime, timezone
import hashlib
import hmac
import json
import os
from pathlib import Path
import re
import shutil
import sys
from typing import Any, Callable
from urllib.error import HTTPError
from urllib.request import Request, urlopen
import uuid


SOURCE_ID = "meta-whatsapp-api-examples"
SOURCE_COMMIT = "de70ee908a67026e642aaee3703d20464e2a9466"
Transport = Callable[[str, str, dict[str, str], bytes, float], tuple[int, dict[str, str], bytes]]


def _sha256(value: bytes) -> str:
    return hashlib.sha256(value).hexdigest()


def _object(path: Path) -> dict[str, Any]:
    value = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(value, dict):
        raise ValueError(f"{path.name} must contain one JSON object")
    return value


def verify_packaged_official_source(root: Path, lock_path: Path) -> None:
    lock = _object(lock_path)
    if lock.get("source_id") != SOURCE_ID or lock.get("commit") != SOURCE_COMMIT or lock.get("commit_signature_verified") is not True:
        raise RuntimeError("official Meta source identity/signature mismatch")
    resolved = root.resolve()
    for entry in lock.get("files", []):
        relative = str(entry.get("packaged_path", ""))
        if not relative or Path(relative).is_absolute() or ".." in Path(relative).parts:
            raise RuntimeError("unsafe official source path")
        target = (resolved / relative).resolve()
        if resolved not in target.parents or not target.is_file() or _sha256(target.read_bytes()) != entry.get("sha256"):
            raise RuntimeError(f"packaged official source mismatch: {relative}")


def validate_profile(profile: dict[str, Any]) -> None:
    required = (
        "official_source_license_accepted", "platform_terms_accepted", "business_account_proven",
        "app_registration_proven", "phone_number_id_proven", "message_template_approval_proven",
        "webhook_subscription_proven", "test_recipient_consent_proven", "quota_and_cost_approved",
        "reconciliation_approved",
    )
    if profile.get("provider") != "Meta WhatsApp Cloud API" or profile.get("source_id") != SOURCE_ID or profile.get("source_commit") != SOURCE_COMMIT:
        raise ValueError("WhatsApp profile authority mismatch")
    if profile.get("decision") != "PROVEN" or any(profile.get(key) is not True for key in required):
        raise PermissionError("WhatsApp access, terms or reconciliation remain blocked")
    if profile.get("automatic_business_write") is not False or not str(profile.get("data_retention", "")).strip():
        raise ValueError("retention is required and automatic_business_write must remain false")
    if not re.fullmatch(r"v[0-9]{1,3}\.[0-9]{1,2}", str(profile.get("graph_api_version", ""))):
        raise ValueError("graph_api_version must be explicitly approved, for example vNN.0")
    for key in ("business_account_id", "phone_number_id"):
        if not isinstance(profile.get(key), str) or not re.fullmatch(r"[0-9]{5,30}", profile[key]):
            raise ValueError(f"{key} must be a configured numeric string")
    for key in ("access_token_environment_variable", "app_secret_environment_variable", "verify_token_environment_variable"):
        if not re.fullmatch(r"[A-Z][A-Z0-9_]{2,80}", str(profile.get(key, ""))):
            raise ValueError(f"{key} must name an environment variable")
    approved = profile.get("approved_templates")
    if not isinstance(approved, list) or not approved:
        raise PermissionError("at least one exact approved template is required")
    seen: set[tuple[str, str]] = set()
    for item in approved:
        if not isinstance(item, dict):
            raise ValueError("approved template entries must be objects")
        name, language, count = item.get("name"), item.get("language_code"), item.get("body_parameter_count")
        if not isinstance(name, str) or not re.fullmatch(r"[a-z0-9_]{1,512}", name):
            raise ValueError("approved template name is invalid")
        if not isinstance(language, str) or not re.fullmatch(r"[a-z]{2}(?:_[A-Z]{2})?", language):
            raise ValueError("approved template language is invalid")
        if not isinstance(count, int) or isinstance(count, bool) or not 0 <= count <= 20 or (name, language) in seen:
            raise ValueError("approved template parameter count/identity is invalid")
        seen.add((name, language))


def build_template_payload(profile: dict[str, Any], request: dict[str, Any]) -> tuple[str, dict[str, Any]]:
    validate_profile(profile)
    recipient = str(request.get("recipient", ""))
    if not re.fullmatch(r"[1-9][0-9]{7,14}", recipient):
        raise ValueError("recipient must be an E.164 number without plus sign")
    name, language, parameters = request.get("template_name"), request.get("language_code"), request.get("body_parameters")
    if not isinstance(parameters, list) or len(parameters) > 20 or any(not isinstance(value, str) or not value or len(value) > 1024 for value in parameters):
        raise ValueError("body_parameters must be bounded non-empty strings")
    approved = {(item["name"], item["language_code"]): item["body_parameter_count"] for item in profile["approved_templates"]}
    if approved.get((name, language)) != len(parameters):
        raise PermissionError("template, language or parameter count is not approved")
    components = []
    if parameters:
        components.append({"type": "body", "parameters": [{"type": "text", "text": value} for value in parameters]})
    payload = {"messaging_product": "whatsapp", "recipient_type": "individual", "to": recipient, "type": "template", "template": {"name": name, "language": {"code": language}, "components": components}}
    return recipient, payload


def build_reply_payload(profile: dict[str, Any], request: dict[str, Any], *, observed_at: float | None = None) -> tuple[str, dict[str, Any]]:
    """ADAPTED get_text_message_input at the already-pinned Meta commit.

    Scope/window guards are local composition of a human-approved proposal;
    a text message is never interpreted as an approved Meta template.
    """
    validate_profile(profile)
    if set(request) != {"kind", "recipient", "text", "source_message_id", "last_inbound_at", "window_expires_at"} or request.get("kind") != "text_reply":
        raise ValueError("text reply contract mismatch")
    recipient, text = request.get("recipient"), request.get("text")
    if not isinstance(recipient, str) or not re.fullmatch(r"[1-9][0-9]{7,14}", recipient):
        raise ValueError("reply recipient must be exact E.164 digits")
    # Conservative local budget, not a statement of the provider maximum.
    if not isinstance(text, str) or not 0 < len(text) <= 1024 or len(text.encode("utf-8", "strict")) > 4096 or any(ord(c)<32 and c not in "\n\t" for c in text):
        raise ValueError("reply text exceeds local contract")
    _event_identity(request.get("source_message_id"), "source message")
    start, end = request.get("last_inbound_at"), request.get("window_expires_at")
    if type(start) is not int or type(end) is not int or start<1 or end!=start+86400:
        raise ValueError("reply service window is invalid")
    now=datetime.now(timezone.utc).timestamp() if observed_at is None else observed_at
    if now<start or now>=end:
        raise PermissionError("reply service window is closed")
    return recipient, {"messaging_product":"whatsapp", "recipient_type":"individual", "to":recipient, "type":"text", "text":{"body":text}}


def stdlib_transport(method: str, url: str, headers: dict[str, str], body: bytes, timeout: float) -> tuple[int, dict[str, str], bytes]:
    request = Request(url, data=body, headers=headers, method=method)
    try:
        with urlopen(request, timeout=timeout) as response:
            return response.status, dict(response.headers.items()), response.read()
    except HTTPError as error:
        return error.code, dict(error.headers.items()), error.read()


def send_template_to_evidence(profile: dict[str, Any], request: dict[str, Any], access_token: str, output_directory: Path, transport: Transport = stdlib_transport) -> dict[str, Any]:
    if not access_token:
        raise PermissionError("WhatsApp access token is unavailable")
    recipient, payload = (build_reply_payload(profile, request) if request.get("kind") == "text_reply" else build_template_payload(profile, request))
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".whatsapp-send-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        request_bytes = (json.dumps(payload, ensure_ascii=False, sort_keys=True, separators=(",", ":")) + "\n").encode("utf-8")
        version, phone_id = profile["graph_api_version"], profile["phone_number_id"]
        url = f"https://graph.facebook.com/{version}/{phone_id}/messages"
        status, _, response_bytes = transport("POST", url, {"Authorization": f"Bearer {access_token}", "Content-Type": "application/json"}, request_bytes, 30.0)
        try:
            response = json.loads(response_bytes.decode("utf-8"))
        except (UnicodeDecodeError, json.JSONDecodeError) as error:
            raise RuntimeError("WhatsApp response is not UTF-8 JSON") from error
        if status < 200 or status >= 300 or not isinstance(response, dict):
            raise RuntimeError(f"WhatsApp provider rejected message with HTTP {status}")
        messages = response.get("messages")
        message_id = messages[0].get("id") if isinstance(messages, list) and messages and isinstance(messages[0], dict) else None
        if not isinstance(message_id, str) or not message_id:
            raise RuntimeError("WhatsApp success response is missing message id")
        (stage / "provider-response.json").write_bytes(response_bytes)
        receipt = {"schema": "elite-whatsapp-cloud-send-receipt/v1", "created_at": datetime.now(timezone.utc).isoformat(), "source_commit": SOURCE_COMMIT, "graph_api_version": version, "phone_number_id_sha256": _sha256(phone_id.encode("ascii")), "recipient_sha256": _sha256(recipient.encode("ascii")), "request_sha256": _sha256(request_bytes), "response_sha256": _sha256(response_bytes), "message_id_sha256": _sha256(message_id.encode("utf-8")), "automatic_business_write": False}
        (stage / "SEND_RECEIPT.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        stage.replace(output)
        return receipt
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def verify_webhook_signature(raw_body: bytes, signature_header: str, app_secret: str) -> None:
    if not app_secret:
        raise PermissionError("Meta app secret is unavailable")
    if not re.fullmatch(r"sha256=[0-9a-f]{64}", signature_header or ""):
        raise PermissionError("WhatsApp webhook signature header is missing or malformed")
    expected = "sha256=" + hmac.new(app_secret.encode("utf-8"), raw_body, hashlib.sha256).hexdigest()
    if not hmac.compare_digest(signature_header, expected):
        raise PermissionError("WhatsApp webhook signature mismatch")


def verify_subscription(query: dict[str, str], verify_token: str) -> str:
    challenge = query.get("hub.challenge", "")
    if not verify_token or query.get("hub.mode") != "subscribe" or not challenge or not hmac.compare_digest(query.get("hub.verify_token", ""), verify_token):
        raise PermissionError("WhatsApp webhook subscription verification failed")
    return challenge


def _unique_object(pairs: list[tuple[str, Any]]) -> dict[str, Any]:
    result: dict[str, Any] = {}
    for key, value in pairs:
        if key in result:
            raise ValueError("WhatsApp webhook contains duplicate JSON keys")
        result[key] = value
    return result


def _event_array(value: Any) -> list[Any]:
    if not isinstance(value, list) or len(value) > 1000:
        raise ValueError("WhatsApp webhook arrays must be bounded lists")
    return value


def _event_identity(value: Any, label: str) -> str:
    if not isinstance(value, str) or not value or len(value) > 512 or any(ord(c) < 33 or ord(c) > 126 for c in value):
        raise ValueError(f"WhatsApp {label} must be a bounded non-empty identity")
    return value


def normalize_verified_webhook(raw_body: bytes, signature_header: str, app_secret: str, *, profile: dict[str, Any]) -> tuple[dict[str, str], list[dict[str, Any]]]:
    """The shared signature/scope/parser boundary; no disk or business effects."""
    validate_profile(profile)
    # Local resource budget, not a claimed Meta API limit. Reject before parsing.
    if not isinstance(raw_body, bytes) or not 0 < len(raw_body) <= 1024 * 1024:
        raise ValueError("WhatsApp webhook exceeds local body budget")
    verify_webhook_signature(raw_body, signature_header, app_secret)
    try:
        body = json.loads(raw_body.decode("utf-8"), object_pairs_hook=_unique_object)
    except (UnicodeDecodeError, json.JSONDecodeError, RecursionError) as error:
        raise ValueError("WhatsApp webhook must be UTF-8 JSON") from error
    if not isinstance(body, dict) or body.get("object") != "whatsapp_business_account" or not isinstance(body.get("entry"), list):
        raise ValueError("WhatsApp webhook envelope is invalid")
    events: list[dict[str, Any]] = []
    scope = {"business_account_id_sha256": _sha256(profile["business_account_id"].encode("ascii")), "phone_number_id_sha256": _sha256(profile["phone_number_id"].encode("ascii"))}
    for entry in _event_array(body["entry"]):
        if not isinstance(entry, dict):
            raise ValueError("WhatsApp entry is invalid")
        if entry.get("id") != profile["business_account_id"]:
            raise PermissionError("WhatsApp business account scope mismatch")
        for change in _event_array(entry.get("changes")):
            if not isinstance(change, dict) or change.get("field") != "messages" or not isinstance(change.get("value"), dict):
                raise ValueError("WhatsApp change is not an admitted messages envelope")
            value = change["value"]
            metadata = value.get("metadata")
            if value.get("messaging_product") != "whatsapp" or not isinstance(metadata, dict):
                raise ValueError("WhatsApp product/metadata is missing")
            if metadata.get("phone_number_id") != profile["phone_number_id"]:
                raise PermissionError("WhatsApp phone number scope mismatch")
            for field, kind, contact_key in (("messages", "message", "from"), ("statuses", "status", "recipient_id")):
                for event in _event_array(value.get(field, [])):
                    if not isinstance(event, dict):
                        raise ValueError("WhatsApp event must be an object")
                    identity = _event_identity(event.get("id"), "message id")
                    contact = _event_identity(event.get(contact_key), "contact id")
                    timestamp = event.get("timestamp")
                    if not isinstance(timestamp, str) or not re.fullmatch(r"[1-9][0-9]{0,11}", timestamp):
                        raise ValueError("WhatsApp timestamp must be canonical positive Unix seconds")
                    normalized_event = {**scope, "kind": kind, "id_sha256": _sha256(identity.encode("ascii")), "timestamp": timestamp, "event_sha256": _sha256(json.dumps(event, sort_keys=True, separators=(",", ":")).encode("utf-8"))}
                    normalized_event["sender_sha256" if kind == "message" else "recipient_sha256"] = _sha256(contact.encode("ascii"))
                    if kind == "status":
                        status = event.get("status")
                        if not isinstance(status, str) or status not in {"sent", "delivered", "read", "failed", "deleted"}:
                            raise ValueError("WhatsApp status is not admitted")
                        normalized_event["status"] = status
                    else:
                        message_type = event.get("type")
                        if not isinstance(message_type, str) or not re.fullmatch(r"[a-z][a-z0-9_]{0,63}", message_type):
                            raise ValueError("WhatsApp message type is invalid")
                        if message_type == "text":
                            text = event.get("text")
                            if not isinstance(text, dict) or not isinstance(text.get("body"), str) or not 0 < len(text["body"].encode("utf-8", "strict")) <= 16384:
                                raise ValueError("text message body is invalid")
                        normalized_event["type"] = message_type
                    # Stable provider-event identity independent of batch order and JSON layout.
                    identity_fields = {key: val for key, val in normalized_event.items() if key != "event_sha256"}
                    normalized_event["event_key_sha256"] = _sha256(json.dumps(identity_fields, sort_keys=True, separators=(",", ":")).encode("ascii"))
                    events.append(normalized_event)
                    if len(events) > 1000:
                        raise ValueError("WhatsApp webhook exceeds local event budget")
    if not events:
        raise ValueError("WhatsApp webhook contains no admitted message/status events")
    return scope, events


def webhook_to_evidence(raw_body: bytes, signature_header: str, app_secret: str, output_directory: Path, *, profile: dict[str, Any]) -> dict[str, Any]:
    scope, events = normalize_verified_webhook(raw_body, signature_header, app_secret, profile=profile)
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".whatsapp-webhook-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        normalized = (json.dumps({"schema": "elite-whatsapp-cloud-events/v2", "events": events}, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "normalized-events.json").write_bytes(normalized)
        receipt = {"schema": "elite-whatsapp-cloud-webhook-receipt/v2", **scope, "profile_sha256": _sha256(json.dumps(profile, sort_keys=True, separators=(",", ":")).encode("utf-8")), "created_at": datetime.now(timezone.utc).isoformat(), "source_commit": SOURCE_COMMIT, "raw_body_sha256": _sha256(raw_body), "normalized_sha256": _sha256(normalized), "event_count": len(events), "raw_payload_persisted": False, "automatic_business_write": False}
        (stage / "WEBHOOK_RECEIPT.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        stage.replace(output)
        return receipt
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def send_from_files(profile_path: Path, request_path: Path, output: Path) -> dict[str, Any]:
    profile, message = _object(profile_path), _object(request_path)
    secret_name = str(profile.get("access_token_environment_variable", ""))
    return send_template_to_evidence(profile, message, os.environ.get(secret_name, ""), output)


def send_bridge(raw: bytes, transport: Transport = stdlib_transport) -> dict[str, Any]:
    """One approved send, invoked only behind the existing durable Go fence.

    The frame contains secrets in memory; neither the frame nor exceptions may
    be printed/persisted. This bridge adds no retry or delivery-state owner.
    """
    if not 0 < len(raw) <= 131072:
        raise ValueError("bridge frame exceeds local budget")
    frame = json.loads(raw.decode("utf-8"), object_pairs_hook=_unique_object)
    fields = {"schema", "binding_sha256", "profile", "request", "access_token", "output_directory"}
    if not isinstance(frame, dict) or set(frame) != fields or frame["schema"] != "elite-whatsapp-send-bridge/v1":
        raise ValueError("bridge frame contract mismatch")
    binding = frame["binding_sha256"]
    if not isinstance(binding, str) or not re.fullmatch(r"[0-9a-f]{64}", binding):
        raise ValueError("bridge binding is invalid")
    if not isinstance(frame["profile"], dict) or not isinstance(frame["request"], dict):
        raise ValueError("bridge profile/request must be objects")
    token, destination = frame["access_token"], frame["output_directory"]
    if not isinstance(token, str) or not 0 < len(token) <= 16384 or not isinstance(destination, str) or not Path(destination).is_absolute():
        raise ValueError("bridge secret or output unavailable")
    root = Path(__file__).resolve().parent
    verify_packaged_official_source(root, root / "official-source.lock.json")
    receipt = send_template_to_evidence(frame["profile"], frame["request"], token, Path(destination), transport)
    evidence = Path(destination)
    provider_bytes = (evidence / "provider-response.json").read_bytes()
    if _sha256(provider_bytes) != receipt["response_sha256"]:
        raise ValueError("provider evidence changed")
    response = json.loads(provider_bytes.decode("utf-8"), object_pairs_hook=_unique_object)
    messages = response.get("messages")
    if not isinstance(messages, list) or len(messages) != 1 or not isinstance(messages[0], dict):
        raise ValueError("provider message identity is ambiguous")
    identity = _event_identity(messages[0].get("id"), "provider message id")
    if len(identity) > 256 or _sha256(identity.encode("ascii")) != receipt["message_id_sha256"]:
        raise ValueError("provider message identity mismatch")
    return {"schema": "elite-whatsapp-send-result/v1", "binding_sha256": binding,
            "provider_message_id": identity, "evidence_sha256": _sha256((evidence / "SEND_RECEIPT.json").read_bytes()),
            "accepted_at": receipt["created_at"]}


def ingress_bridge(raw: bytes) -> dict[str, Any]:
    """AUTHORED stdio boundary reusing the existing Meta verifier; no I/O effects."""
    import base64
    if not 0 < len(raw) <= 2 * 1024 * 1024:
        raise ValueError("ingress frame exceeds local budget")
    frame = json.loads(raw.decode("utf-8"), object_pairs_hook=_unique_object)
    if not isinstance(frame, dict) or set(frame) != {"schema", "mode", "profile", "body", "signature", "secret", "query"} or frame["schema"] != "elite-whatsapp-ingress-bridge/v1":
        raise ValueError("ingress frame contract mismatch")
    if not isinstance(frame["secret"], str) or not 0 < len(frame["secret"]) <= 16384:
        raise PermissionError("ingress secret unavailable")
    root = Path(__file__).resolve().parent
    verify_packaged_official_source(root, root / "official-source.lock.json")
    reply = {"schema": "elite-whatsapp-ingress-result/v1", "binding_sha256": _sha256(raw), "body_sha256": "", "event_count": 0, "challenge": ""}
    if frame["mode"] == "subscribe":
        query = frame["query"]
        if not isinstance(query, dict) or set(query) != {"hub.mode", "hub.verify_token", "hub.challenge"} or any(not isinstance(v, str) for v in query.values()):
            raise ValueError("subscription query mismatch")
        if not re.fullmatch(r"[0-9]{1,256}", query["hub.challenge"]):
            raise ValueError("challenge exceeds local contract")
        if frame["body"] or frame["signature"]:
            raise ValueError("subscription cannot contain body or signature")
        reply["challenge"] = verify_subscription(query, frame["secret"])
    elif frame["mode"] == "receive":
        if frame["query"] or not isinstance(frame["body"], str):
            raise ValueError("ingress body contract mismatch")
        body = base64.b64decode(frame["body"], validate=True)
        _, events = normalize_verified_webhook(body, frame["signature"], frame["secret"], profile=frame["profile"])
        reply.update(body_sha256=_sha256(body), event_count=len(events))
    else:
        raise ValueError("ingress mode rejected")
    return reply


def recover_send_bridge(raw: bytes) -> dict[str, Any]:
    """Pure evidence validation; never calls transport or changes delivery state."""
    import base64
    if not 0<len(raw)<=262144:
        raise ValueError("recovery frame exceeds local budget")
    frame=json.loads(raw.decode("utf-8"),object_pairs_hook=_unique_object)
    if not isinstance(frame,dict) or set(frame)!={"schema","binding_sha256","profile","request","receipt","response"} or frame["schema"]!="elite-whatsapp-recover-send/v1":
        raise ValueError("recovery frame contract mismatch")
    if not re.fullmatch(r"[0-9a-f]{64}",str(frame["binding_sha256"])):
        raise ValueError("recovery binding invalid")
    root=Path(__file__).resolve().parent
    verify_packaged_official_source(root,root/"official-source.lock.json")
    receipt_raw=base64.b64decode(frame["receipt"],validate=True)
    response_raw=base64.b64decode(frame["response"],validate=True)
    if not 0<len(receipt_raw)<=65536 or not 0<len(response_raw)<=65536:
        raise ValueError("recovery evidence exceeds local budget")
    receipt=json.loads(receipt_raw.decode("utf-8"),object_pairs_hook=_unique_object)
    response=json.loads(response_raw.decode("utf-8"),object_pairs_hook=_unique_object)
    at=datetime.fromisoformat(receipt["created_at"])
    if at.tzinfo is None or at.timestamp()>datetime.now(timezone.utc).timestamp()+60:
        raise ValueError("recovery acceptance time invalid")
    request,profile=frame["request"],frame["profile"]
    if not isinstance(request,dict):
        raise ValueError("recovery request must be an object")
    if request.get("kind")=="text_reply":
        recipient,payload=build_reply_payload(profile,request,observed_at=at.timestamp())
    elif set(request)=={"recipient","template_name","language_code","body_parameters"}:
        recipient,payload=build_template_payload(profile,request)
    else:
        raise ValueError("recovery request kind or shape is not supported")
    expected=(json.dumps(payload,ensure_ascii=False,sort_keys=True,separators=(",",":"))+"\n").encode("utf-8")
    if receipt.get("schema")!="elite-whatsapp-cloud-send-receipt/v1" or receipt.get("source_commit")!=SOURCE_COMMIT or receipt.get("automatic_business_write") is not False or receipt.get("graph_api_version")!=profile["graph_api_version"] or receipt.get("phone_number_id_sha256")!=_sha256(profile["phone_number_id"].encode("ascii")) or receipt.get("recipient_sha256")!=_sha256(recipient.encode("ascii")) or receipt.get("request_sha256")!=_sha256(expected) or receipt.get("response_sha256")!=_sha256(response_raw):
        raise ValueError("recovery evidence binding mismatch")
    messages=response.get("messages")
    if not isinstance(messages,list) or len(messages)!=1 or not isinstance(messages[0],dict):
        raise ValueError("recovery provider identity ambiguous")
    identity=_event_identity(messages[0].get("id"),"provider message id")
    if len(identity)>256 or receipt.get("message_id_sha256")!=_sha256(identity.encode("ascii")):
        raise ValueError("recovery provider identity mismatch")
    return {"schema":"elite-whatsapp-send-result/v1","binding_sha256":frame["binding_sha256"],"provider_message_id":identity,"evidence_sha256":_sha256(receipt_raw),"accepted_at":receipt["created_at"]}


def validate_profile_bridge(raw: bytes) -> dict[str, str]:
    """Bounded startup validation of the exact configuration bytes; no I/O."""
    import base64
    if not 0 < len(raw) <= 65536:
        raise ValueError("profile validation frame exceeds local budget")
    frame=json.loads(raw.decode("utf-8"),object_pairs_hook=_unique_object)
    if not isinstance(frame,dict) or set(frame)!={"schema","profile"} or frame["schema"]!="elite-whatsapp-validate-profile/v1" or not isinstance(frame["profile"],str):
        raise ValueError("profile validation contract mismatch")
    profile_raw=base64.b64decode(frame["profile"],validate=True)
    if not 0<len(profile_raw)<=32768:
        raise ValueError("profile exceeds local budget")
    profile=json.loads(profile_raw.decode("utf-8"),object_pairs_hook=_unique_object)
    if not isinstance(profile,dict):
        raise ValueError("profile must be an object")
    validate_profile(profile)
    return {"schema":"elite-whatsapp-profile-validation/v1","profile_sha256":_sha256(profile_raw)}


if __name__ == "__main__":
    if sys.argv[1:] not in (["--send-bridge"], ["--status-bridge"], ["--ingress-bridge"], ["--recover-send-bridge"], ["--validate-profile-bridge"]):
        sys.exit(2)
    try:
        if sys.argv[1:] == ["--send-bridge"]:
            reply = send_bridge(sys.stdin.buffer.read(131073))
        elif sys.argv[1:] == ["--validate-profile-bridge"]:
            reply = validate_profile_bridge(sys.stdin.buffer.read(65537))
        elif sys.argv[1:] == ["--recover-send-bridge"]:
            reply = recover_send_bridge(sys.stdin.buffer.read(262145))
        elif sys.argv[1:] == ["--ingress-bridge"]:
            reply = ingress_bridge(sys.stdin.buffer.read(2 * 1024 * 1024 + 1))
        else:
            # Fixed local module; Go verifies both source hashes before -I/-B.
            import importlib.util
            sys.modules["whatsapp_cloud"] = sys.modules[__name__]
            spec = importlib.util.spec_from_file_location("status_reconciliation", Path(__file__).resolve().parent / "status_reconciliation.py")
            module = importlib.util.module_from_spec(spec)
            spec.loader.exec_module(module)
            reply = module.status_bridge(sys.stdin.buffer.read(12 * 1024 * 1024 + 1))
        sys.stdout.buffer.write((json.dumps(reply, sort_keys=True, separators=(",", ":")) + "\n").encode("utf-8"))
    except Exception:
        # No provider body, token, recipient or exception text crosses stderr.
        sys.stderr.buffer.write(b"WHATSAPP_INGRESS_UNVERIFIED\n" if sys.argv[1:] == ["--ingress-bridge"] else b"WHATSAPP_STATUS_UNVERIFIED\n" if sys.argv[1:] == ["--status-bridge"] else b"WHATSAPP_SEND_UNCERTAIN_OR_REJECTED\n")
        sys.exit(2)
