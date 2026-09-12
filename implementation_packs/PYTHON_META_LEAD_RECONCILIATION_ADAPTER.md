# Python Meta Lead Reconciliation Adapter

## 1. Metadata

```yaml
pack_id: "PYTHON-META-LEAD-RECONCILIATION-ADAPTER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa recuperación y reconciliación read-only de leads mediante los edges oficiales LeadgenForm /leads y /test_leads del SDK Meta Business 26.0.1, preserva filas exactas y produce candidatos separados sin conceder contacto ni escritura CRM automática."
stacks: ["CPython 3.14 Windows x86-64", "facebook-business 26.0.1", "Meta Graph API v26.0"]
compatible_with: ["GO-OMNICHANNEL-LEAD-INGRESS 0.2.x", "GO-LEAD-CANDIDATE-PROMOTION 0.1.x", "OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x"]
incompatible_with: ["contacto automático", "write CRM directo", "webhook Meta no demostrado", "campos no aprobados", "secreto en CLI/Markdown", "uso fuera de APIs Meta"]
license_expression: "LicenseRef-Workspace-Owner AND LicenseRef-Meta-Platform AND third-party-wheel-licenses"
upstream_sources: ["https://github.com/facebook/facebook-python-business-sdk/blob/788f363d15b1269ab5efb7cd00fb5e3b133cd99b/facebook_business/adobjects/lead.py", "https://github.com/facebook/facebook-python-business-sdk/blob/788f363d15b1269ab5efb7cd00fb5e3b133cd99b/facebook_business/adobjects/leadgenform.py", "https://pypi.org/project/facebook-business/26.0.1/"]
verified_at: "2026-09-04"
```

## 2. Applicability

Use para recuperar y reconciliar leads de formularios Meta cuando el proyecto demuestre app/business, autorización de acceso a leads, ownership de Page/form, test lead oficial, controles de PII, retención, cuota/costo y reconciliación. No presenta un webhook Meta no demostrado, no concede consentimiento y no escribe directamente al CRM.

`fetch_leads.py` es glue `ADAPTED` sobre los métodos y campos públicos exactos de `LeadgenForm`/`Lead` del SDK oficial fijado. Meta gobierna el contrato del SDK; la normalización, los gates y la evidencia son código local y no se atribuyen a Meta. La licencia oficial limita el SDK a los servicios/APIs Facebook y debe conservarse.

## 3. Architecture contract

La plantilla nace bloqueada. `TEST` usa exclusivamente `LeadgenForm.get_test_leads`; `LIVE` usa exclusivamente `LeadgenForm.get_leads`. Ambos preservan `provider-response.json`, generan un lote de resultados candidato/rechazo separado y un receipt sin form ID ni PII. Los campos de formulario no aprobados quedan sólo en evidencia raw y fuerzan rechazo del candidato. Un resultado recuperado conserva `contact_eligibility=pending_policy`; nunca equivale a consentimiento, lead CRM, cita o venta.

El límite local de filas queda registrado. Si se alcanza, el resultado declara truncamiento y no prueba reconciliación completa. Un error de SDK/provider/serialización elimina staging. Un directorio de salida existente se rechaza para impedir sobrescritura.

## 4. Exact file manifest

```text
CREATE meta_lead_reconciliation/requirements-direct.in
CREATE meta_lead_reconciliation/requirements-windows-py314.lock
CREATE meta_lead_reconciliation/sdk-artifact.lock.json
CREATE meta_lead_reconciliation/provider-profile.template.json
CREATE meta_lead_reconciliation/retrieval.template.json
CREATE meta_lead_reconciliation/fetch_leads.py
CREATE meta_lead_reconciliation/test_fetch_leads.py
CREATE meta_lead_reconciliation/README.md
```

## 5. Materialization blocks

### FILE: `meta_lead_reconciliation/requirements-direct.in`
```yaml
block_id: "PY-META-LEAD-RECONCILIATION:requirements-direct:v1"
operation: CREATE
provenance: AUTHORED
source: "exact official PyPI wheel URL and SHA-256 inherited from the admitted Meta SDK artifact lock"
license: "LicenseRef-Workspace-Owner"
sha256: "b0fe2f0d83c899fd9986a9990b9102b613774bd912aaffd5219736b014051c88"
variables: []
secrets_allowed: false
```
````text
facebook-business @ https://files.pythonhosted.org/packages/bd/fb/038e02129ff35dbe1eac38bd7ab40234a08576359ef8806ca6d132e2d558/facebook_business-26.0.1-py3-none-any.whl#sha256=41555ce87ac105617ed678a489a630108ac162fbdccfbbf65fa20ae1de7cd0ca
````

### FILE: `meta_lead_reconciliation/requirements-windows-py314.lock`
```yaml
block_id: "PY-META-LEAD-RECONCILIATION:requirements-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "exact admitted PyPI wheel graph for CPython 3.14 Windows x86-64"
license: "LicenseRef-Workspace-Owner"
sha256: "a4b4b65d2a19be4bb51892bcb25a4e8127f7b2c16b8f49a2f40b319797314e19"
variables: []
secrets_allowed: false
```
````text
# Exact wheel graph resolved and verified for CPython 3.14 on Windows x86-64.
aiohappyeyeballs @ https://files.pythonhosted.org/packages/71/43/1947f06babed6b3f1d7f38b0c767f52df66bfb2bc10b468c4a7de9eceff2/aiohappyeyeballs-2.7.1-py3-none-any.whl#sha256=9243213661e29250eb41368e5daa826fc017156c3b8a11440826b2e3ed376472
aiohttp @ https://files.pythonhosted.org/packages/f5/8b/c7baa1ba1eda4db6989baefe5de6d99834921b84ebd7918624febcb9f290/aiohttp-3.14.3-cp314-cp314-win_amd64.whl#sha256=8b3b60de05f3dcb6f6a00f818bb2ec781cee4de0645f59ccaf99b1d1823b6100
aiosignal @ https://files.pythonhosted.org/packages/fb/76/641ae371508676492379f16e2fa48f4e2c11741bd63c48be4b12a6b09cba/aiosignal-1.4.0-py3-none-any.whl#sha256=053243f8b92b990551949e63930a839ff0cf0b0ebbe0597b0f3fb19e1a0fe82e
attrs @ https://files.pythonhosted.org/packages/64/b4/17d4b0b2a2dc85a6df63d1157e028ed19f90d4cd97c36717afef2bc2f395/attrs-26.1.0-py3-none-any.whl#sha256=c647aa4a12dfbad9333ca4e71fe62ddc36f4e63b2d260a37a8b83d2f043ac309
capi-param-builder-python @ https://files.pythonhosted.org/packages/4f/1d/463144a77228b0137b498f70eb4ac60808fddba7bc890308729b76197391/capi_param_builder_python-1.3.0-py3-none-any.whl#sha256=641fd8832378033670783a8841dc3168fd75f5ab9629779d11a6ce0a6687f86d
certifi @ https://files.pythonhosted.org/packages/0b/a7/71ac2cff56fec219ed242bb11b8efb69fcc4bec75db06fb7bfe35de520e6/certifi-2026.7.22-py3-none-any.whl#sha256=62f22742b58a1a33014a2b6b706588a8d7e2a88ae7bd1a6ebe8c992928483775
charset-normalizer @ https://files.pythonhosted.org/packages/7a/7c/4938c329b6a9d446f6a59aa2092ff7118f274209b5ed0e26893d1d30a63c/charset_normalizer-3.5.1-cp314-cp314-win_amd64.whl#sha256=c658c50ac0c98cd755a2dd50b7977d3bca7df401dcc47fbdfa87db53ef7d4e8b
curlify @ https://files.pythonhosted.org/packages/9e/f8/912ebddbff8ea603d4c90fa31557096f927b17efd30a166ce7ac1242910a/curlify-3.0.0-py3-none-any.whl#sha256=52060c0eb7a656b7bde6b668c32f337bed4d736ce230755767e3ada56a09c338
facebook-business @ https://files.pythonhosted.org/packages/bd/fb/038e02129ff35dbe1eac38bd7ab40234a08576359ef8806ca6d132e2d558/facebook_business-26.0.1-py3-none-any.whl#sha256=41555ce87ac105617ed678a489a630108ac162fbdccfbbf65fa20ae1de7cd0ca
frozenlist @ https://files.pythonhosted.org/packages/59/ad/9caa9b9c836d9ad6f067157a531ac48b7d36499f5036d4141ce78c230b1b/frozenlist-1.8.0-cp314-cp314-win_amd64.whl#sha256=3e0761f4d1a44f1d1a47996511752cf3dcec5bbdd9cc2b4fe595caf97754b7a0
idna @ https://files.pythonhosted.org/packages/57/b0/0e52c878c53f245edd3a11020f20979b3f490f245af532c7cae3027754b5/idna-3.19-py3-none-any.whl#sha256=815e7be7a7806d54abb586dc943addc79e8b2ee16915059658cbeff4b1b43bf4
multidict @ https://files.pythonhosted.org/packages/e0/bf/52f25716bbe93745595800f36fb17b73711f14da59ed0bb2eba141bc9f0f/multidict-6.7.1-cp314-cp314-win_amd64.whl#sha256=5e01429a929600e7dab7b166062d9bb54a5eed752384c7384c968c2afab8f50f
propcache @ https://files.pythonhosted.org/packages/61/d2/45c9defbaa1ea297035d9d4cce9e8f80daafbf19319c6007f157c6256ea9/propcache-0.5.2-cp314-cp314-win_amd64.whl#sha256=81e3a30b0bb60caa22033dd0f8a3618d1d67356212514f62c57db75cb0ef410c
pycountry @ https://files.pythonhosted.org/packages/9c/42/7703bd45b62fecd44cd7d3495423097e2f7d28bc2e99e7c1af68892ab157/pycountry-26.2.16-py3-none-any.whl#sha256=115c4baf7cceaa30f59a4694d79483c9167dbce7a9de4d3d571c5f3ea77c305a
requests @ https://files.pythonhosted.org/packages/a0/f4/c67b0b3f1b9245e8d266f0f112c500d50e5b4e83cb6f3b71b6528104182a/requests-2.34.2-py3-none-any.whl#sha256=2a0d60c172f83ac6ab31e4554906c0f3b3588d37b5cb939b1c061f4907e278e0
six @ https://files.pythonhosted.org/packages/b7/ce/149a00dd41f10bc29e5921b496af8b574d8413afcd5e30dfa0ed46c2cc5e/six-1.17.0-py2.py3-none-any.whl#sha256=4721f391ed90541fddacab5acf947aa0d3dc7d27b2e1e8eda2be8970586c3274
urllib3 @ https://files.pythonhosted.org/packages/7f/3e/5db95bcf282c52709639744ca2a8b149baccf648e39c8cc87553df9eae0c/urllib3-2.7.0-py3-none-any.whl#sha256=9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897
yarl @ https://files.pythonhosted.org/packages/cf/52/6daa2ee9d95e5c98b8128f8df91eb692eb423ab274b8cf08db52152fad26/yarl-1.24.5-cp314-cp314-win_amd64.whl#sha256=5ba4f78df2bcc19f764a4b26a8a4f5049c110090ad5825993aacb052bf8003ad
````

### FILE: `meta_lead_reconciliation/sdk-artifact.lock.json`
```yaml
block_id: "PY-META-LEAD-RECONCILIATION:sdk-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "local receipt derived from official GitHub/PyPI metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "03eb5d7bd845cd02d42045e1d74ee307976cc31eadc5b8c4925b222f86c5a887"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-python-sdk-artifact-lock/v1",
  "provider": "Meta",
  "distribution": "facebook-business",
  "version": "26.0.1",
  "verified_platform_lock": "CPython 3.14 / Windows x86-64 / 18 wheels",
  "wheel": {
    "filename": "facebook_business-26.0.1-py3-none-any.whl",
    "bytes": 1554832,
    "sha256": "41555ce87ac105617ed678a489a630108ac162fbdccfbbf65fa20ae1de7cd0ca"
  },
  "source": {
    "repository": "facebook/facebook-python-business-sdk",
    "release": "26.0.1",
    "commit": "788f363d15b1269ab5efb7cd00fb5e3b133cd99b",
    "archive_bytes": 2070568,
    "archive_sha256": "bd3bc1f14072662eecb05fb6f43888eb652c9a1f039894f99a7e7609e4452db7",
    "commit_signature_verified": false
  },
  "capi_parameter_builder_source": {
    "repository": "facebook/capi-param-builder",
    "release": "v1.3.0-python",
    "commit": "66a7eb84e29a65aa999089366c1f7c3c0e49bdb7",
    "archive_bytes": 674665,
    "archive_sha256": "be8d21efab5a1ea16b1f6e84d4219d9e343436b8789ede407d0b94d0a54dabf5"
  },
  "license_expression": "LicenseRef-Meta-Platform",
  "license_sha256": "48d97b3c936203a750a3288c2b924327769d62327383a1060244b3d3c05ad01f",
  "distribution_notice": "capi-param-builder-python 1.3.0 wheel omits LICENSE; preserve the exact official source LICENSE with redistributions",
  "verified_at": "2026-08-26"
}
````

### FILE: `meta_lead_reconciliation/provider-profile.template.json`
```yaml
block_id: "PY-META-LEAD-RECONCILIATION:provider-profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed access, policy, PII, quota and reconciliation gate"
license: "LicenseRef-Workspace-Owner"
sha256: "f3417cf84c7b4865863cea224bb5e30283aa82bfc5285303bf67df3c7b396833"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-meta-lead-reconciliation-profile/v1",
  "provider": "Meta Marketing API",
  "sdk": "facebook-business==26.0.1",
  "graph_api_version": "v26.0",
  "decision": "BLOCKED_ACCESS_AND_POLICY_APPROVAL_REQUIRED",
  "app_registration_proven": false,
  "business_verification_proven": false,
  "lead_access_permission_proven": false,
  "page_and_form_ownership_proven": false,
  "test_lead_contract_proven": false,
  "app_id_environment_variable": "META_APP_ID",
  "app_secret_environment_variable": "META_APP_SECRET",
  "access_token_environment_variable": "META_ACCESS_TOKEN",
  "tenant_id": "",
  "organization_id": "",
  "form_id": "",
  "approved_provider_fields": [],
  "approved_form_field_names": [],
  "data_retention": "",
  "pii_storage_controls_proven": false,
  "quota_and_cost_approved": false,
  "reconciliation_approved": false,
  "automatic_business_write": false,
  "automatic_contact_eligibility": false
}
````

### FILE: `meta_lead_reconciliation/retrieval.template.json`
```yaml
block_id: "PY-META-LEAD-RECONCILIATION:retrieval:v1"
operation: CREATE
provenance: AUTHORED
source: "local bounded query template over exact official Lead field constants"
license: "LicenseRef-Workspace-Owner"
sha256: "c11aee0e9c8d59619a63d32bbf1d048114494582345375d9340bc80b3976e462"
variables: []
secrets_allowed: false
```
````json
{
  "retrieval_mode": "TEST",
  "fields": [
    "id",
    "created_time",
    "form_id",
    "campaign_id",
    "campaign_name",
    "adset_id",
    "adset_name",
    "ad_id",
    "ad_name",
    "field_data",
    "is_organic",
    "platform"
  ],
  "max_rows": 100
}
````

### FILE: `meta_lead_reconciliation/fetch_leads.py`
```yaml
block_id: "PY-META-LEAD-RECONCILIATION:fetch:v1"
operation: CREATE
provenance: ADAPTED
source: "Meta LeadgenForm.get_leads/get_test_leads and Lead field contracts at exact commit 788f363d15b1269ab5efb7cd00fb5e3b133cd99b; local validation, evidence and normalization"
license: "LicenseRef-Meta-Platform AND LicenseRef-Workspace-Owner"
sha256: "5ab4269967b9d12065ea007a49206ef503b4c6e41638b19de0eb579920b71876"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
from datetime import datetime, timezone
import hashlib
import importlib.metadata
import json
import os
from pathlib import Path
import re
import shutil
import sys
from typing import Any, Callable, Iterable
import uuid

from facebook_business.adobjects.leadgenform import LeadgenForm
from facebook_business.api import FacebookAdsApi


SDK_VERSION = "26.0.1"
API_VERSION = "v26.0"
SOURCE_COMMIT = "788f363d15b1269ab5efb7cd00fb5e3b133cd99b"
ALLOWED_PROVIDER_FIELDS = {
    "id", "created_time", "form_id", "campaign_id", "campaign_name",
    "adset_id", "adset_name", "ad_id", "ad_name", "field_data",
    "is_organic", "platform", "custom_disclaimer_responses", "home_listing",
}


def sha256(value: bytes) -> str:
    return hashlib.sha256(value).hexdigest()


def _read_json(path: Path) -> dict[str, Any]:
    value = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(value, dict):
        raise ValueError(f"{path.name} must contain one JSON object")
    return value


def validate_configuration(profile: dict[str, Any], retrieval: dict[str, Any]) -> tuple[str, str, str, list[str], set[str], int]:
    proof_fields = (
        "app_registration_proven", "business_verification_proven",
        "lead_access_permission_proven", "page_and_form_ownership_proven",
        "test_lead_contract_proven", "pii_storage_controls_proven",
        "quota_and_cost_approved", "reconciliation_approved",
    )
    if (
        profile.get("provider") != "Meta Marketing API"
        or profile.get("sdk") != f"facebook-business=={SDK_VERSION}"
        or profile.get("graph_api_version") != API_VERSION
    ):
        raise ValueError("profile provider/SDK/API identity mismatch")
    if profile.get("decision") != "PROVEN" or any(profile.get(key) is not True for key in proof_fields):
        raise PermissionError("Meta lead access/policy profile remains blocked")
    if profile.get("automatic_business_write") is not False or profile.get("automatic_contact_eligibility") is not False:
        raise ValueError("automatic writes and contact eligibility must remain false")
    if not str(profile.get("data_retention", "")).strip():
        raise ValueError("data_retention is required")
    tenant_id = str(profile.get("tenant_id", "")).strip()
    organization_id = str(profile.get("organization_id", "")).strip()
    if not tenant_id or not organization_id:
        raise ValueError("tenant_id and organization_id are required")
    form_id = str(profile.get("form_id", ""))
    if not re.fullmatch(r"[0-9]{5,30}", form_id):
        raise ValueError("form_id must contain digits only")
    fields = retrieval.get("fields")
    approved_fields = profile.get("approved_provider_fields")
    if not isinstance(fields, list) or not fields or len(fields) > 20 or len(fields) != len(set(fields)):
        raise ValueError("fields must be a unique non-empty list with at most 20 entries")
    if not {"id", "created_time", "field_data"}.issubset(set(fields)):
        raise ValueError("id, created_time and field_data are mandatory")
    if not isinstance(approved_fields, list) or not set(fields).issubset(set(approved_fields)) or not set(fields).issubset(ALLOWED_PROVIDER_FIELDS):
        raise PermissionError("all provider fields must be locally allowed and explicitly approved")
    approved_form_fields = profile.get("approved_form_field_names")
    if not isinstance(approved_form_fields, list) or not approved_form_fields:
        raise ValueError("approved_form_field_names must be a non-empty list")
    approved_names: set[str] = set()
    for value in approved_form_fields:
        if not isinstance(value, str) or not re.fullmatch(r"[A-Za-z0-9_.-]{1,128}", value) or value in approved_names:
            raise ValueError("approved_form_field_names contains an invalid or duplicate name")
        approved_names.add(value)
    mode = retrieval.get("retrieval_mode")
    if mode not in {"TEST", "LIVE"}:
        raise ValueError("retrieval_mode must be TEST or LIVE")
    max_rows = retrieval.get("max_rows")
    if not isinstance(max_rows, int) or isinstance(max_rows, bool) or max_rows < 1 or max_rows > 1000:
        raise ValueError("max_rows must be within 1..1000")
    return tenant_id, organization_id, form_id, fields, approved_names, max_rows


def _row_to_dict(row: Any) -> dict[str, Any]:
    if isinstance(row, dict):
        return row
    exporter = getattr(row, "export_all_data", None)
    if not callable(exporter):
        raise RuntimeError("Meta lead row is not an official SDK object")
    value = exporter()
    if not isinstance(value, dict):
        raise RuntimeError("Meta SDK lead export must return an object")
    return value


def _normalize_row(row: dict[str, Any], tenant_id: str, organization_id: str, approved_names: set[str], is_test: bool) -> dict[str, Any]:
    candidate: dict[str, Any] = {
        "tenant_id": tenant_id,
        "organization_id": organization_id,
        "provider": "meta_lead_ads",
        "provider_lead_id": str(row.get("id", "")),
        "form_id": str(row.get("form_id", "")),
        "campaign_id": str(row.get("campaign_id", "")),
        "ad_group_id": str(row.get("adset_id", "")),
        "creative_id": str(row.get("ad_id", "")),
        "source_kind": str(row.get("platform", "")),
        "submitted_at": str(row.get("created_time", "")),
        "is_test": is_test,
        "fields": [],
        "contact_eligibility": "pending_policy",
    }
    errors: list[str] = []
    if not candidate["provider_lead_id"]:
        errors.append("MISSING_LEAD_ID")
    if not candidate["submitted_at"]:
        errors.append("MISSING_CREATED_TIME")
    field_data = row.get("field_data")
    if not isinstance(field_data, list):
        errors.append("INVALID_FIELD_DATA")
        field_data = []
    seen: set[str] = set()
    for item in field_data:
        if not isinstance(item, dict) or not isinstance(item.get("name"), str) or not isinstance(item.get("values"), list):
            errors.append("INVALID_FIELD_ENTRY")
            continue
        name = item["name"]
        values = item["values"]
        if name not in approved_names:
            errors.append("UNAPPROVED_FORM_FIELD")
            continue
        if name in seen:
            errors.append("DUPLICATE_FORM_FIELD")
            continue
        if len(values) != 1 or not isinstance(values[0], str) or not values[0].strip():
            errors.append("INVALID_FIELD_VALUE")
            continue
        seen.add(name)
        candidate["fields"].append({"id": name, "value": values[0]})
    return {
        "status": "CANDIDATE" if not errors else "REJECTED",
        "normalization_codes": sorted(set(errors)),
        "candidate": candidate if not errors else None,
        "provider_lead_id_sha256": sha256(candidate["provider_lead_id"].encode("utf-8")),
    }


def fetch_to_evidence(
    form_factory: Callable[[str], Any], tenant_id: str, organization_id: str,
    form_id: str, fields: list[str], approved_names: set[str], max_rows: int,
    retrieval_mode: str, output_directory: Path,
) -> dict[str, Any]:
    if importlib.metadata.version("facebook-business") != SDK_VERSION:
        raise RuntimeError(f"facebook-business {SDK_VERSION} is required")
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".meta-lead-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        form = form_factory(form_id)
        method = form.get_test_leads if retrieval_mode == "TEST" else form.get_leads
        rows: list[dict[str, Any]] = []
        truncated = False
        for raw in method(fields=fields, params={}):
            if len(rows) == max_rows:
                truncated = True
                break
            rows.append(_row_to_dict(raw))
        response = {
            "graph_api_version": API_VERSION,
            "retrieval_mode": retrieval_mode,
            "rows": rows,
            "truncated_by_local_bound": truncated,
        }
        response_bytes = (json.dumps(response, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "provider-response.json").write_bytes(response_bytes)
        outcomes = [_normalize_row(row, tenant_id, organization_id, approved_names, retrieval_mode == "TEST") for row in rows]
        candidate_document = {
            "schema": "elite-meta-lead-candidate-batch/v1",
            "provider": "meta_lead_ads",
            "source_commit": SOURCE_COMMIT,
            "tenant_id": tenant_id,
            "organization_id": organization_id,
            "provider_response_sha256": sha256(response_bytes),
            "outcomes": outcomes,
            "automatic_business_write": False,
            "automatic_contact_eligibility": False,
        }
        candidate_bytes = (json.dumps(candidate_document, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "lead-candidates.json").write_bytes(candidate_bytes)
        query_bytes = json.dumps({"fields": fields, "max_rows": max_rows, "retrieval_mode": retrieval_mode}, sort_keys=True, separators=(",", ":")).encode("utf-8")
        receipt = {
            "schema": "elite-meta-lead-reconciliation-receipt/v1",
            "created_at": datetime.now(timezone.utc).isoformat(),
            "provider": "Meta Marketing API",
            "sdk_version": SDK_VERSION,
            "graph_api_version": API_VERSION,
            "source_commit": SOURCE_COMMIT,
            "form_id_sha256": sha256(form_id.encode("ascii")),
            "query_sha256": sha256(query_bytes),
            "provider_response_sha256": sha256(response_bytes),
            "candidate_batch_sha256": sha256(candidate_bytes),
            "row_count": len(rows),
            "candidate_count": sum(item["status"] == "CANDIDATE" for item in outcomes),
            "rejected_count": sum(item["status"] == "REJECTED" for item in outcomes),
            "truncated_by_local_bound": truncated,
            "automatic_business_write": False,
            "automatic_contact_eligibility": False,
        }
        (stage / "RETRIEVAL_RECEIPT.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        stage.replace(output)
        return receipt
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", required=True, type=Path)
    parser.add_argument("--retrieval", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()
    profile, retrieval = _read_json(args.profile), _read_json(args.retrieval)
    tenant_id, organization_id, form_id, fields, approved_names, max_rows = validate_configuration(profile, retrieval)
    names = ("app_id_environment_variable", "app_secret_environment_variable", "access_token_environment_variable")
    environment_names = [str(profile.get(name, "")) for name in names]
    if any(not re.fullmatch(r"[A-Z][A-Z0-9_]{2,80}", name) for name in environment_names):
        raise ValueError("secret references must be canonical environment variable names")
    values = [os.environ.get(name, "") for name in environment_names]
    if any(not value for value in values):
        raise PermissionError("required Meta credential environment variables are unavailable")
    api = FacebookAdsApi.init(app_id=values[0], app_secret=values[1], access_token=values[2], api_version=API_VERSION)
    receipt = fetch_to_evidence(
        lambda identifier: LeadgenForm(identifier, api=api), tenant_id, organization_id,
        form_id, fields, approved_names, max_rows, retrieval["retrieval_mode"], args.output,
    )
    print(json.dumps({"status": "META_LEAD_RECONCILIATION_PASS", "receipt": receipt}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `meta_lead_reconciliation/test_fetch_leads.py`
```yaml
block_id: "PY-META-LEAD-RECONCILIATION:test:v1"
operation: CREATE
provenance: ADAPTED
source: "introspection of the exact official Meta SDK methods plus local fail-closed regressions"
license: "LicenseRef-Meta-Platform AND LicenseRef-Workspace-Owner"
sha256: "0fcb84c03bfc99c7fec1b484b2a19b822f56ac0e56e84547ad83b6ca99ed20b8"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import inspect
import json
from pathlib import Path
import tempfile
import unittest

from facebook_business.adobjects.leadgenform import LeadgenForm
from fetch_leads import fetch_to_evidence, validate_configuration


class FakeRow:
    def __init__(self, value): self.value = value
    def export_all_data(self): return self.value


class FakeForm:
    def __init__(self, identifier, rows=None, error=None):
        self.identifier = identifier
        self.rows = rows or []
        self.error = error
        self.calls = []

    def _run(self, mode, fields, params):
        self.calls.append((mode, fields, params))
        if self.error: raise self.error
        return self.rows

    def get_leads(self, *, fields, params): return self._run("LIVE", fields, params)
    def get_test_leads(self, *, fields, params): return self._run("TEST", fields, params)


def proven_profile():
    return {
        "provider": "Meta Marketing API", "sdk": "facebook-business==26.0.1", "graph_api_version": "v26.0",
        "decision": "PROVEN", "app_registration_proven": True, "business_verification_proven": True,
        "lead_access_permission_proven": True, "page_and_form_ownership_proven": True,
        "test_lead_contract_proven": True, "pii_storage_controls_proven": True,
        "quota_and_cost_approved": True, "reconciliation_approved": True,
        "automatic_business_write": False, "automatic_contact_eligibility": False,
        "data_retention": "30 days", "tenant_id": "tenant-1", "organization_id": "org-1",
        "form_id": "123456789", "approved_provider_fields": ["id", "created_time", "form_id", "field_data", "campaign_id", "adset_id", "ad_id", "platform"],
        "approved_form_field_names": ["email", "full_name"],
    }


def valid_retrieval(mode="LIVE"):
    return {"retrieval_mode": mode, "fields": ["id", "created_time", "form_id", "field_data", "campaign_id"], "max_rows": 100}


def valid_row():
    return {
        "id": "lead-1", "created_time": "2026-09-04T12:00:00+0000", "form_id": "123456789",
        "campaign_id": "campaign-1", "field_data": [
            {"name": "email", "values": ["person@example.test"]},
            {"name": "full_name", "values": ["Test Person"]},
        ],
    }


class MetaLeadReconciliationTests(unittest.TestCase):
    def test_configuration_fails_closed(self):
        profile = proven_profile(); profile["decision"] = "BLOCKED_ACCESS_AND_POLICY_APPROVAL_REQUIRED"
        with self.assertRaises(PermissionError): validate_configuration(profile, valid_retrieval())
        profile = proven_profile(); profile["automatic_contact_eligibility"] = True
        with self.assertRaises(ValueError): validate_configuration(profile, valid_retrieval())

    def test_configuration_rejects_unapproved_fields_and_invalid_form(self):
        profile = proven_profile(); query = valid_retrieval(); query["fields"].append("ad_name")
        with self.assertRaises(PermissionError): validate_configuration(profile, query)
        profile = proven_profile(); profile["form_id"] = "form_123"
        with self.assertRaises(ValueError): validate_configuration(profile, valid_retrieval())

    def test_live_fetch_preserves_provider_and_emits_pending_candidate(self):
        form = FakeForm("", [FakeRow(valid_row())])
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = fetch_to_evidence(lambda identifier: self._bind(form, identifier), "tenant-1", "org-1", "123456789", valid_retrieval()["fields"], {"email", "full_name"}, 100, "LIVE", output)
            self.assertEqual(form.calls[0][0], "LIVE")
            self.assertEqual(receipt["candidate_count"], 1)
            self.assertNotIn("123456789", json.dumps(receipt))
            self.assertNotIn("person@example.test", json.dumps(receipt))
            provider = json.loads((output / "provider-response.json").read_text(encoding="utf-8"))
            candidate = json.loads((output / "lead-candidates.json").read_text(encoding="utf-8"))["outcomes"][0]["candidate"]
            batch = json.loads((output / "lead-candidates.json").read_text(encoding="utf-8"))
            self.assertEqual(provider["rows"][0]["field_data"][0]["values"][0], "person@example.test")
            self.assertEqual(batch["tenant_id"], "tenant-1")
            self.assertEqual(batch["organization_id"], "org-1")
            self.assertEqual(candidate["contact_eligibility"], "pending_policy")
            self.assertFalse(candidate["is_test"])
            self.assertEqual(candidate["fields"][0]["value"], "person@example.test")

    def test_test_fetch_marks_candidates_test(self):
        form = FakeForm("", [valid_row()])
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            fetch_to_evidence(lambda identifier: self._bind(form, identifier), "tenant-1", "org-1", "123456789", valid_retrieval("TEST")["fields"], {"email", "full_name"}, 100, "TEST", output)
            candidate = json.loads((output / "lead-candidates.json").read_text(encoding="utf-8"))["outcomes"][0]["candidate"]
            self.assertEqual(form.calls[0][0], "TEST")
            self.assertTrue(candidate["is_test"])

    def test_unknown_form_field_is_rejected_without_losing_raw(self):
        row = valid_row(); row["field_data"].append({"name": "unapproved", "values": ["secret"]})
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = fetch_to_evidence(lambda identifier: FakeForm(identifier, [row]), "tenant-1", "org-1", "123456789", valid_retrieval()["fields"], {"email", "full_name"}, 100, "LIVE", output)
            outcome = json.loads((output / "lead-candidates.json").read_text(encoding="utf-8"))["outcomes"][0]
            self.assertEqual(receipt["rejected_count"], 1)
            self.assertEqual(outcome["normalization_codes"], ["UNAPPROVED_FORM_FIELD"])
            self.assertIsNone(outcome["candidate"])
            self.assertIn("unapproved", (output / "provider-response.json").read_text(encoding="utf-8"))

    def test_multiple_values_are_retained_raw_and_rejected_as_ambiguous(self):
        row = valid_row(); row["field_data"][0]["values"] = ["first@example.test", "second@example.test"]
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = fetch_to_evidence(lambda identifier: FakeForm(identifier, [row]), "tenant-1", "org-1", "123456789", valid_retrieval()["fields"], {"email", "full_name"}, 100, "LIVE", output)
            outcome = json.loads((output / "lead-candidates.json").read_text(encoding="utf-8"))["outcomes"][0]
            self.assertEqual(receipt["rejected_count"], 1)
            self.assertEqual(outcome["normalization_codes"], ["INVALID_FIELD_VALUE"])
            self.assertIsNone(outcome["candidate"])
            self.assertIn("second@example.test", (output / "provider-response.json").read_text(encoding="utf-8"))

    def test_local_bound_is_explicit(self):
        rows = [valid_row(), {**valid_row(), "id": "lead-2"}]
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = fetch_to_evidence(lambda identifier: FakeForm(identifier, rows), "tenant-1", "org-1", "123456789", valid_retrieval()["fields"], {"email", "full_name"}, 1, "LIVE", output)
            self.assertEqual(receipt["row_count"], 1)
            self.assertTrue(receipt["truncated_by_local_bound"])

    def test_provider_failure_is_atomic(self):
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            with self.assertRaisesRegex(RuntimeError, "provider failed"):
                fetch_to_evidence(lambda identifier: FakeForm(identifier, error=RuntimeError("provider failed")), "tenant-1", "org-1", "123456789", valid_retrieval()["fields"], {"email"}, 1, "LIVE", output)
            self.assertFalse(output.exists())
            self.assertEqual(list(Path(temp).iterdir()), [])

    def test_existing_output_is_rejected(self):
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"; output.mkdir()
            with self.assertRaises(FileExistsError):
                fetch_to_evidence(lambda identifier: FakeForm(identifier), "tenant-1", "org-1", "123456789", valid_retrieval()["fields"], {"email"}, 1, "LIVE", output)

    def test_official_sdk_contract(self):
        for method_name in ("get_leads", "get_test_leads"):
            method = getattr(LeadgenForm, method_name)
            signature = inspect.signature(method)
            self.assertIn("fields", signature.parameters)
            self.assertIn("params", signature.parameters)
        self.assertIn("/leads", inspect.getsource(LeadgenForm.get_leads))
        self.assertIn("/test_leads", inspect.getsource(LeadgenForm.get_test_leads))

    @staticmethod
    def _bind(form, identifier): form.identifier = identifier; return form


if __name__ == "__main__": unittest.main()
````

### FILE: `meta_lead_reconciliation/README.md`
```yaml
block_id: "PY-META-LEAD-RECONCILIATION:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local operational instructions constrained by the exact SDK and license"
license: "LicenseRef-Workspace-Owner"
sha256: "39d7bd2931409e3fcfacffec0509861e49931a72f3a5ec5a11a0ec535e95b791"
variables: []
secrets_allowed: false
```
````markdown
# Meta Lead Reconciliation Adapter

Read-only recovery and reconciliation over the official Meta Business SDK for Python 26.0.1 and Graph API v26.0. It calls only `LeadgenForm.get_leads` or `LeadgenForm.get_test_leads`, preserves the exact exported provider rows and emits a separate candidate batch plus hash-linked receipt. The wrapper is local `ADAPTED` glue; it is not Meta-authored code.

The template starts blocked. Before use, prove the app and business registration, lead-access authorization, Page/form ownership, an official test-lead contract, PII storage controls, retention, quota/cost and reconciliation. Secrets belong only in the named environment variables. A retrieved lead is neither permission to contact nor an automatic CRM write. Form-field names not explicitly approved, malformed values and multi-value fields are retained only in the provider evidence and the corresponding candidate is rejected; a versioned project mapping is required before multi-select data can be promoted.

`TEST` calls the SDK's official `/test_leads` edge and marks every candidate `is_test=true`. `LIVE` calls `/leads`. The candidate batch carries tenant and organization scope independently so rejected raw evidence remains isolated even when `candidate=null`. The locally enforced row bound is recorded; a truncated retrieval is not a complete reconciliation and must be continued or rerun under the project's provider procedure.

The frozen lock is for CPython 3.14 on Windows x86-64. Meta's platform license restricts use of the SDK to Facebook services/APIs. Preserve the exact official LICENSE and the source/artifact lock.

```powershell
python -m venv .venv
.venv\Scripts\python -m pip install --require-hashes -r .\meta_lead_reconciliation\requirements-windows-py314.lock
Push-Location .\meta_lead_reconciliation
..\.venv\Scripts\python -m unittest -v test_fetch_leads.py
Pop-Location
```

No fixture proves live access, completeness, provider retention, rate limits, cost, consent or downstream delivery. The output remains candidate evidence until the central ingress and explicit promotion gates accept it.
````

## 6. Configuration surface

- App ID, app secret y access token sólo por referencias a variables de entorno.
- Identidad tenant/organization y form ID explícitos.
- Campos provider y nombres `field_data` aprobados de forma separada.
- Modo `TEST` o `LIVE`, retención PII y límite local 1..1000.
- Todas las pruebas de acceso, ownership, test contract, controles PII, cuota/costo y reconciliación deben estar `true` y `decision=PROVEN`.

## 7. Dependency bill

El lock reutiliza exactamente el grafo de 18 wheels ya admitido para `facebook-business==26.0.1` en CPython 3.14/Windows x86-64. El wheel oficial Meta tiene SHA-256 `41555ce87ac105617ed678a489a630108ac162fbdccfbbf65fa20ae1de7cd0ca`; source commit `788f363d15b1269ab5efb7cd00fb5e3b133cd99b`. El grafo conserva las licencias registradas, incluido LGPL-2.1-only de `pycountry`, y el source LICENSE de CAPI Parameter Builder porque su wheel lo omite.

## 8. Apply order

1. Materializar en destino vacío.
2. Crear el source profile del proyecto y validar cuenta/acceso/políticas.
3. Instalar con `--require-hashes` usando el lock exacto del target.
4. Ejecutar tests, `compileall` y `pip check`.
5. Probar `/test_leads` con formulario oficial test y reconciliar cantidad/contenido.
6. Sólo entonces habilitar `LIVE`; enviar candidatos al owner durable central y mantener promoción CRM separada.

## 9. Verification

```powershell
python -m venv .venv
.venv\Scripts\python -m pip install --require-hashes -r .\meta_lead_reconciliation\requirements-windows-py314.lock
Push-Location .\meta_lead_reconciliation
..\.venv\Scripts\python -m unittest -v test_fetch_leads.py
..\.venv\Scripts\python -m compileall -q .
..\.venv\Scripts\python -m pip check
Pop-Location
```

### Production blockers

- Credenciales/cuenta/Page/form reales y permiso de lectura no están disponibles en la biblioteca.
- La prueba local no demuestra retención de Meta, exhaustividad histórica, rate limits, billing ni disponibilidad regional.
- El lote candidato todavía requiere ingesta durable y promoción explícita; no prueba cita, venta ni outcome feedback.
- El webhook Meta no está implementado en este pack porque su contrato público exacto no fue admitido.

## 10. Reconstruction evidence

`reconstruction_evidence/PYTHON_META_LEAD_RECONCILIATION_ADAPTER_2026-09-04_V226.md`
