# Python Azure Blob Immutable Evidence Adapter

## 1. Metadata

```yaml
pack_id: "PYTHON-AZURE-BLOB-IMMUTABLE-EVIDENCE-ADAPTER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa glue Python sobre el SDK oficial Microsoft azure-storage-blob 12.30.0 para creación única de Block Blob con WORM Locked atómico, versionado, Content-MD5, cifrado por scope y verificación SHA-256 de la versión exacta descargada."
stacks: ["CPython 3.12.13", "azure-storage-blob 12.30.0", "Azure Blob version-level WORM"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "AZURE-CONTENT-UNDERSTANDING-DOCUMENT-RUNTIME 0.1.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x", "STRICT-DOCUMENT-FIELD-EVALUATION-GATE 0.1.x"]
incompatible_with: ["credenciales en archivos", "overwrite silencioso", "contenedor público", "versionado o WORM no verificado", "scope de cifrado reemplazable", "auto-storage documental", "habilitación o migración automática de inmutabilidad"]
license_expression: "LicenseRef-Workspace-Owner AND MIT"
upstream_sources: ["https://pypi.org/project/azure-storage-blob/12.30.0/", "https://github.com/Azure/azure-sdk-for-python/tree/245d2024e6764bb4dd5186b5cd30a776c46993cb/sdk/storage/azure-storage-blob", "https://learn.microsoft.com/azure/storage/blobs/immutable-storage-overview", "https://learn.microsoft.com/azure/storage/blobs/immutable-policy-configure-version-scope", "https://learn.microsoft.com/python/api/azure-storage-blob/azure.storage.blob.blobclient"]
verified_at: "2026-08-27"
```

## 2. Applicability

Use sólo cuando el blueprint seleccione Azure Blob Storage para evidencia documental y el usuario haya aprobado tenant/subscription, resource ID de la cuenta, RBAC, costo, red, contenedor, scope de cifrado, retención `Locked` y cualquier legal hold. Los bytes deben haber pasado ingreso seguro y evaluación estricta; este pack no extrae documentos ni convierte resultados del proveedor en hechos de negocio.

Rechazar si el contenedor es público, no reporta immutable storage with versioning, permite reemplazar el scope de cifrado, falta un receipt de autoridad fresco ligado por SHA-256 o se pretende que el runtime cree/migre/bloquee políticas. Microsoft advierte que habilitar/migrar version-level WORM puede ser irreversible; esa administración queda fuera del adapter. El SDK es Microsoft MIT; todos los archivos materializados son glue `AUTHORED` local y no se atribuyen a Microsoft.

## 3. Architecture contract

`EvidenceStore.put` valida configuración, approvals, receipt cerrado, límites, MIME y retención antes del provider. Consulta el estado real del contenedor. El backend oficial llama `BlobClient.upload_blob` con nombre determinista, `overwrite=False`, `validate_content=True`, Content-MD5 almacenado, metadata SHA-256, `ImmutabilityPolicy(..., LOCKED)` y legal hold opcional en la misma operación. La cuenta aplica un scope de cifrado default exacto y no reemplazable.

El éxito exige `version_id` y luego `get_blob_properties(version_id=...)` más `download_blob(version_id=..., validate_content=True, decompress=False)`. Se comparan versión, ETag, tamaño, metadata, MD5, cifrado servidor/scope, modo/expiración WORM, legal hold y SHA-256 de bytes descargados. Timeout/respuesta perdida produce receipt `UNKNOWN_EFFECT`; mismatch posterior produce `PARTIAL_VERIFICATION`. Ambos impiden retry ciego y conservan el nombre determinista; `reconcile` verifica una versión explícita sin crear otra.

El adapter nunca crea cuentas/contenedores, cambia RBAC/red/versionado/WORM/cifrado, elimina versiones, retira holds ni guarda credenciales. Disponibilidad, throughput, costo, cuotas, logs, key lifecycle, restore y cumplimiento se prueban en el target. Descargar para verificar agrega lectura/costo deliberados; el perfil debe aprobarlo.

## 4. Exact file manifest

```text
CREATE azure_blob_storage_adapter/requirements-direct.in
CREATE azure_blob_storage_adapter/requirements-windows-py312.lock
CREATE azure_blob_storage_adapter/provider-profile.template.json
CREATE azure_blob_storage_adapter/official-artifact-lock.json
CREATE azure_blob_storage_adapter/azureblob/__init__.py
CREATE azure_blob_storage_adapter/azureblob/storage.py
CREATE azure_blob_storage_adapter/azureblob/azure_backend.py
CREATE azure_blob_storage_adapter/test_storage.py
CREATE azure_blob_storage_adapter/README.md
```

## 5. Materialization blocks

### FILE: `azure_blob_storage_adapter/requirements-direct.in`
```yaml
block_id: "PYTHON-AZURE-BLOB-EVIDENCE:requirements-direct:v1"
operation: CREATE
provenance: AUTHORED
source: "local exact reference to the official Microsoft PyPI wheel"
license: "LicenseRef-Workspace-Owner"
sha256: "6b5a20c1db334ba9c03f477f5b8d28e7018bc74bd4624c13add2532e96a79bbb"
variables: []
secrets_allowed: false
```
````text
azure-storage-blob @ https://files.pythonhosted.org/packages/5e/0b/e106f0fd7fa785867d9ffcc47dc9e6237c0e58f51058473b777487a98edc/azure_storage_blob-12.30.0-py3-none-any.whl#sha256=d415ac50b67a8da6b3ae7e9f1014b1b55cd7aafa0b8d4ca9b380568dc7360423
````

### FILE: `azure_blob_storage_adapter/requirements-windows-py312.lock`
```yaml
block_id: "PYTHON-AZURE-BLOB-EVIDENCE:requirements-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "dated exact PyPI resolution for CPython 3.12 Windows x86-64"
license: "LicenseRef-Workspace-Owner"
sha256: "2156fc146c023a35e110422f8c766948e84f68ef50a7e2231b8720db23eb3bf8"
variables: []
secrets_allowed: false
```
````text
azure-core==1.41.0 --hash=sha256:522b4011e8180b1a3dcd2024396a4e7fe9ac37fb8597db47163d230b5efe892d
azure-storage-blob==12.30.0 --hash=sha256:d415ac50b67a8da6b3ae7e9f1014b1b55cd7aafa0b8d4ca9b380568dc7360423
certifi==2026.7.22 --hash=sha256:62f22742b58a1a33014a2b6b706588a8d7e2a88ae7bd1a6ebe8c992928483775
cffi==2.1.1 --hash=sha256:f53e442b08449d42821fa4a4fba000095af9f62742a500f978a9f557ec44339a
charset-normalizer==3.5.1 --hash=sha256:3617ac3cfd8b9888f145ad89dd6e692285834b0201c6074a5eeaad3fd4d668c2
cryptography==50.0.1 --hash=sha256:aed8db4f6d71c51efb89530e12d9464e7bf2923d46c3205dc794a2a93f8c0648
idna==3.19 --hash=sha256:815e7be7a7806d54abb586dc943addc79e8b2ee16915059658cbeff4b1b43bf4
isodate==0.7.2 --hash=sha256:28009937d8031054830160fce6d409ed342816b543597cece116d966c6d99e15
pycparser==3.0 --hash=sha256:b727414169a36b7d524c1c3e31839a521725078d7b2ff038656844266160a992
requests==2.34.2 --hash=sha256:2a0d60c172f83ac6ab31e4554906c0f3b3588d37b5cb939b1c061f4907e278e0
typing-extensions==4.16.0 --hash=sha256:481caa481374e813c1b176ada14e97f1f67a4539ce9cfeb3f350d78d6370c2e8
urllib3==2.7.0 --hash=sha256:9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897
````

### FILE: `azure_blob_storage_adapter/provider-profile.template.json`
```yaml
block_id: "PYTHON-AZURE-BLOB-EVIDENCE:provider-profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed Azure provider profile"
license: "LicenseRef-Workspace-Owner"
sha256: "645e14a65074995eafebd3a759b3e4cbb6a861ce6c0caef16ec1fb306dbc9e52"
variables: []
secrets_allowed: false
```
````json
{
  "schema_version": "azure-blob-evidence-profile/v1",
  "state": "BLOCKED_CONFIGURATION",
  "account_url": "",
  "account_name": "",
  "storage_account_resource_id": "",
  "authority_receipt_sha256": "",
  "authority_receipt_max_age_seconds": 3600,
  "container": "",
  "default_encryption_scope": "",
  "max_bytes": 0,
  "allowed_content_types": [
    "application/pdf",
    "image/jpeg",
    "image/png",
    "image/tiff"
  ],
  "policy_mode": "Locked",
  "legal_hold_allowed": false,
  "approvals": {
    "access": false,
    "cost": false,
    "retention": false,
    "irreversible_lock": false,
    "legal_hold": false
  },
  "automatic_storage_authorized": false
}
````

### FILE: `azure_blob_storage_adapter/official-artifact-lock.json`
```yaml
block_id: "PYTHON-AZURE-BLOB-EVIDENCE:artifact-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "local verification record for exact official Microsoft artifacts"
license: "LicenseRef-Workspace-Owner"
sha256: "bdd491579586eb07026a497fc547467f5eaad4b2ca824b97298517a33770c6d2"
variables: []
secrets_allowed: false
```
````json
{
  "schema_version": "official-artifact-lock/v1",
  "package": "azure-storage-blob",
  "version": "12.30.0",
  "publisher": "Microsoft Corporation",
  "license": "MIT",
  "pypi": {
    "wheel": {
      "filename": "azure_storage_blob-12.30.0-py3-none-any.whl",
      "bytes": 435610,
      "sha256": "d415ac50b67a8da6b3ae7e9f1014b1b55cd7aafa0b8d4ca9b380568dc7360423"
    },
    "sdist": {
      "filename": "azure_storage_blob-12.30.0.tar.gz",
      "bytes": 618229,
      "sha256": "2cd74d4d5731e5eb6b8d5c5056ee115a5e88f8fdf22517b739836fda685018be"
    }
  },
  "source": {
    "repository": "Azure/azure-sdk-for-python",
    "tag": "azure-storage-blob_12.30.0",
    "commit": "245d2024e6764bb4dd5186b5cd30a776c46993cb",
    "archive_bytes": 205202391,
    "archive_sha256": "7c4743ccfe1cb81687ddc949df8343901c0b9a9e3c4e782ffaf4ad5f68c4a8fd",
    "signature_verified": false,
    "signature_reason": "unsigned"
  },
  "license_sha256": "fd532481d828e13a0b13ccb598e02338a3617740675a862ee6bdc1541b68e93d",
  "verified_at": "2026-08-27"
}
````

### FILE: `azure_blob_storage_adapter/azureblob/__init__.py`
```yaml
block_id: "PYTHON-AZURE-BLOB-EVIDENCE:init:v1"
operation: CREATE
provenance: AUTHORED
source: "local package exports"
license: "LicenseRef-Workspace-Owner"
sha256: "48425acd30bdf81d93b72a71150910ca0986294d4208e4ff1a4bdc4a8c3ef98c"
variables: []
secrets_allowed: false
```
````python
from .storage import (
    Approvals,
    EvidenceStore,
    PartialVerificationError,
    StorageConfig,
    StoreRequest,
    UnknownEffectError,
    ValidationError,
)
from .azure_backend import AzureBlobBackend

__all__ = [
    "Approvals",
    "AzureBlobBackend",
    "EvidenceStore",
    "PartialVerificationError",
    "StorageConfig",
    "StoreRequest",
    "UnknownEffectError",
    "ValidationError",
]
````

### FILE: `azure_blob_storage_adapter/azureblob/storage.py`
```yaml
block_id: "PYTHON-AZURE-BLOB-EVIDENCE:storage:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed evidence storage contract over official provider boundary"
license: "LicenseRef-Workspace-Owner"
sha256: "913d55459f70b6d2e83ec57e11783f8c62d2e2cf9cdad2872f5b6025cbfa63ab"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import base64
import hashlib
import json
import re
from dataclasses import dataclass
from datetime import datetime, timedelta, timezone
from typing import Any, Mapping, Protocol
from urllib.parse import urlparse


_HEX64 = re.compile(r"^[0-9a-f]{64}$")
_SAFE_ID = re.compile(r"^[a-z0-9](?:[a-z0-9_-]{0,62})$")
_CONTAINER = re.compile(r"^[a-z0-9](?:[a-z0-9-]{1,61})[a-z0-9]$")
_SCOPE = re.compile(r"^[A-Za-z0-9](?:[A-Za-z0-9_-]{1,61})[A-Za-z0-9]$")
_ARM_ID = re.compile(
    r"^/subscriptions/[^/]+/resourceGroups/[^/]+/providers/Microsoft\.Storage/storageAccounts/([^/]+)$",
    re.IGNORECASE,
)
_AUTHORITY_KEYS = {
    "schema_version",
    "account_url",
    "account_name",
    "storage_account_resource_id",
    "container",
    "immutable_storage_with_versioning_enabled",
    "default_encryption_scope",
    "prevent_encryption_scope_override",
    "approval_reference",
    "approved_at",
}


class ValidationError(ValueError):
    pass


class UnknownEffectError(RuntimeError):
    def __init__(self, receipt: Mapping[str, Any]):
        super().__init__("provider create outcome is unknown; reconcile the deterministic blob name before retry")
        self.receipt = dict(receipt)


class PartialVerificationError(RuntimeError):
    def __init__(self, receipt: Mapping[str, Any], reason: str):
        super().__init__("created blob version did not pass complete verification: " + reason)
        self.receipt = dict(receipt)


@dataclass(frozen=True)
class Approvals:
    access: bool = False
    cost: bool = False
    retention: bool = False
    irreversible_lock: bool = False
    legal_hold: bool = False


@dataclass(frozen=True)
class StorageConfig:
    account_url: str
    account_name: str
    storage_account_resource_id: str
    authority_receipt_sha256: str
    authority_receipt_max_age_seconds: int
    container: str
    default_encryption_scope: str
    max_bytes: int
    allowed_content_types: frozenset[str]
    legal_hold_allowed: bool = False


@dataclass(frozen=True)
class StoreRequest:
    tenant_id: str
    document_id: str
    data: bytes
    content_type: str
    retain_until: datetime
    legal_hold: bool = False


@dataclass(frozen=True)
class TargetState:
    account_url: str
    account_name: str
    container: str
    public_access: str | None
    immutable_storage_with_versioning_enabled: bool
    default_encryption_scope: str | None
    prevent_encryption_scope_override: bool


@dataclass(frozen=True)
class CreatedVersion:
    version_id: str
    etag: str


@dataclass(frozen=True)
class VersionState:
    version_id: str
    etag: str
    size: int
    metadata: Mapping[str, str]
    content_md5: bytes | None
    server_encrypted: bool
    encryption_scope: str | None
    immutability_policy_mode: str | None
    immutability_expiry: datetime | None
    has_legal_hold: bool


class Backend(Protocol):
    def inspect_target(self, container: str) -> TargetState: ...

    def create_version(
        self,
        *,
        container: str,
        blob_name: str,
        data: bytes,
        content_type: str,
        content_md5: bytes,
        metadata: Mapping[str, str],
        retain_until: datetime,
        legal_hold: bool,
    ) -> CreatedVersion: ...

    def inspect_version(self, *, container: str, blob_name: str, version_id: str) -> VersionState: ...

    def download_version(self, *, container: str, blob_name: str, version_id: str) -> bytes: ...


class EvidenceStore:
    def __init__(self, config: StorageConfig, backend: Backend):
        self._config = config
        self._backend = backend

    def put(
        self,
        request: StoreRequest,
        approvals: Approvals,
        authority_receipt: Mapping[str, Any],
        *,
        now: datetime | None = None,
    ) -> dict[str, Any]:
        instant = _utc_seconds(now or datetime.now(timezone.utc))
        prepared = self._prepare(request, approvals, authority_receipt, instant)
        target = self._backend.inspect_target(self._config.container)
        self._verify_target(target)
        base = self._base_receipt(request, prepared, status="UNKNOWN_EFFECT")
        try:
            created = self._backend.create_version(
                container=self._config.container,
                blob_name=prepared["blob_name"],
                data=request.data,
                content_type=request.content_type,
                content_md5=prepared["content_md5"],
                metadata=prepared["metadata"],
                retain_until=prepared["retain_until"],
                legal_hold=request.legal_hold,
            )
        except Exception as exc:
            raise UnknownEffectError(base) from exc
        partial = dict(base, status="PARTIAL_VERIFICATION", created=True, version_id=created.version_id, etag=created.etag)
        if not created.version_id or not created.etag:
            raise PartialVerificationError(partial, "provider response omitted version_id or etag")
        try:
            return self._verify_version(request, prepared, created.version_id)
        except Exception as exc:
            if isinstance(exc, PartialVerificationError):
                raise
            raise PartialVerificationError(partial, str(exc)) from exc

    def reconcile(
        self,
        request: StoreRequest,
        version_id: str,
        approvals: Approvals,
        authority_receipt: Mapping[str, Any],
        *,
        now: datetime | None = None,
    ) -> dict[str, Any]:
        if not version_id or len(version_id) > 256 or any(ord(ch) < 32 for ch in version_id):
            raise ValidationError("version_id is required and must be a bounded printable value")
        instant = _utc_seconds(now or datetime.now(timezone.utc))
        prepared = self._prepare(request, approvals, authority_receipt, instant)
        self._verify_target(self._backend.inspect_target(self._config.container))
        return self._verify_version(request, prepared, version_id)

    def _prepare(
        self,
        request: StoreRequest,
        approvals: Approvals,
        authority_receipt: Mapping[str, Any],
        now: datetime,
    ) -> dict[str, Any]:
        self._validate_config()
        if not (approvals.access and approvals.cost and approvals.retention and approvals.irreversible_lock):
            raise ValidationError("access, cost, retention and irreversible-lock approvals are required")
        if request.legal_hold and (not self._config.legal_hold_allowed or not approvals.legal_hold):
            raise ValidationError("legal hold requires profile permission and explicit approval")
        if not _SAFE_ID.fullmatch(request.tenant_id) or not _SAFE_ID.fullmatch(request.document_id):
            raise ValidationError("tenant_id and document_id must use the bounded lowercase safe-id grammar")
        if not isinstance(request.data, bytes) or not (0 < len(request.data) <= self._config.max_bytes):
            raise ValidationError("data must be non-empty bytes within max_bytes")
        if request.content_type not in self._config.allowed_content_types:
            raise ValidationError("content_type is not admitted by the profile")
        retain_until = _utc_seconds(request.retain_until)
        delta = retain_until - now
        if delta < timedelta(days=1) or delta > timedelta(days=146000):
            raise ValidationError("Azure retention must be between 1 and 146000 days")
        self._verify_authority_receipt(authority_receipt, now)
        digest = hashlib.sha256(request.data).hexdigest()
        md5 = hashlib.md5(request.data, usedforsecurity=False).digest()
        blob_name = f"evidence/{request.tenant_id}/{request.document_id}/{digest}"
        metadata = {
            "sha256": digest,
            "bytes": str(len(request.data)),
            "authorityreceiptsha256": self._config.authority_receipt_sha256,
        }
        return {
            "blob_name": blob_name,
            "sha256": digest,
            "content_md5": md5,
            "metadata": metadata,
            "retain_until": retain_until,
        }

    def _validate_config(self) -> None:
        parsed = urlparse(self._config.account_url)
        if parsed.scheme != "https" or not parsed.hostname or parsed.path not in ("", "/") or parsed.query or parsed.fragment:
            raise ValidationError("account_url must be an exact HTTPS account endpoint without path, query or fragment")
        if parsed.hostname.split(".", 1)[0].lower() != self._config.account_name.lower():
            raise ValidationError("account_url hostname does not match account_name")
        arm = _ARM_ID.fullmatch(self._config.storage_account_resource_id)
        if not arm or arm.group(1).lower() != self._config.account_name.lower():
            raise ValidationError("storage_account_resource_id must be an exact Microsoft.Storage account resource ID")
        if not _HEX64.fullmatch(self._config.authority_receipt_sha256):
            raise ValidationError("authority_receipt_sha256 must be lowercase SHA-256")
        if not (60 <= self._config.authority_receipt_max_age_seconds <= 86400):
            raise ValidationError("authority receipt max age must be 60..86400 seconds")
        if not _CONTAINER.fullmatch(self._config.container):
            raise ValidationError("container name is invalid")
        if not _SCOPE.fullmatch(self._config.default_encryption_scope):
            raise ValidationError("default_encryption_scope is invalid")
        if self._config.max_bytes <= 0:
            raise ValidationError("max_bytes must be positive")
        if not self._config.allowed_content_types or any(not item or item != item.lower() for item in self._config.allowed_content_types):
            raise ValidationError("allowed_content_types must be a non-empty lowercase set")

    def _verify_authority_receipt(self, receipt: Mapping[str, Any], now: datetime) -> None:
        if set(receipt) != _AUTHORITY_KEYS:
            raise ValidationError("authority receipt schema is not closed")
        canonical = json.dumps(receipt, sort_keys=True, separators=(",", ":"), ensure_ascii=True).encode("utf-8")
        if hashlib.sha256(canonical).hexdigest() != self._config.authority_receipt_sha256:
            raise ValidationError("authority receipt hash does not match the approved profile")
        expected = {
            "schema_version": "azure-storage-authority-receipt/v1",
            "account_url": self._config.account_url,
            "account_name": self._config.account_name,
            "storage_account_resource_id": self._config.storage_account_resource_id,
            "container": self._config.container,
            "immutable_storage_with_versioning_enabled": True,
            "default_encryption_scope": self._config.default_encryption_scope,
            "prevent_encryption_scope_override": True,
        }
        for key, value in expected.items():
            if receipt.get(key) != value:
                raise ValidationError(f"authority receipt field {key} does not match the profile")
        approval_reference = receipt.get("approval_reference")
        if not isinstance(approval_reference, str) or not approval_reference.strip() or len(approval_reference) > 128:
            raise ValidationError("authority receipt approval_reference is invalid")
        approved_at = _parse_rfc3339(receipt.get("approved_at"))
        age = now - approved_at
        if age < timedelta(0) or age > timedelta(seconds=self._config.authority_receipt_max_age_seconds):
            raise ValidationError("authority receipt is future-dated or stale")

    def _verify_target(self, target: TargetState) -> None:
        if target.account_url.rstrip("/") != self._config.account_url.rstrip("/"):
            raise ValidationError("provider account URL differs from the approved target")
        if target.account_name.lower() != self._config.account_name.lower() or target.container != self._config.container:
            raise ValidationError("provider account or container differs from the approved target")
        if target.public_access is not None:
            raise ValidationError("public container access is forbidden")
        if not target.immutable_storage_with_versioning_enabled:
            raise ValidationError("container does not report immutable storage with versioning enabled")
        if target.default_encryption_scope != self._config.default_encryption_scope:
            raise ValidationError("container default encryption scope differs from the approved scope")
        if not target.prevent_encryption_scope_override:
            raise ValidationError("container permits encryption-scope override")

    def _verify_version(self, request: StoreRequest, prepared: Mapping[str, Any], version_id: str) -> dict[str, Any]:
        state = self._backend.inspect_version(
            container=self._config.container,
            blob_name=prepared["blob_name"],
            version_id=version_id,
        )
        checks = {
            "version_id": state.version_id == version_id,
            "etag": bool(state.etag),
            "size": state.size == len(request.data),
            "metadata_sha256": state.metadata.get("sha256") == prepared["sha256"],
            "metadata_bytes": state.metadata.get("bytes") == str(len(request.data)),
            "metadata_authority": state.metadata.get("authorityreceiptsha256") == self._config.authority_receipt_sha256,
            "content_md5": state.content_md5 == prepared["content_md5"],
            "server_encrypted": state.server_encrypted is True,
            "encryption_scope": state.encryption_scope == self._config.default_encryption_scope,
            "policy_mode": (state.immutability_policy_mode or "").lower() == "locked",
            "retain_until": state.immutability_expiry is not None and _utc_seconds(state.immutability_expiry) == prepared["retain_until"],
            "legal_hold": state.has_legal_hold is request.legal_hold,
        }
        failed = sorted(key for key, passed in checks.items() if not passed)
        if failed:
            raise ValidationError("post-write properties mismatch: " + ",".join(failed))
        downloaded = self._backend.download_version(
            container=self._config.container,
            blob_name=prepared["blob_name"],
            version_id=version_id,
        )
        if not isinstance(downloaded, bytes) or len(downloaded) != len(request.data):
            raise ValidationError("exact-version download length mismatch")
        if hashlib.sha256(downloaded).hexdigest() != prepared["sha256"]:
            raise ValidationError("exact-version download SHA-256 mismatch")
        return self._base_receipt(
            request,
            prepared,
            status="VERIFIED",
            created=True,
            version_id=version_id,
            etag=state.etag,
            retention_verified=True,
            bytes_verified_by_download=True,
        )

    def _base_receipt(self, request: StoreRequest, prepared: Mapping[str, Any], *, status: str, **extra: Any) -> dict[str, Any]:
        receipt = {
            "schema_version": "azure-blob-evidence-receipt/v1",
            "status": status,
            "provider": "azure-blob-storage",
            "account_url_sha256": hashlib.sha256(self._config.account_url.encode("utf-8")).hexdigest(),
            "storage_account_resource_id_sha256": hashlib.sha256(
                self._config.storage_account_resource_id.encode("utf-8")
            ).hexdigest(),
            "authority_receipt_sha256": self._config.authority_receipt_sha256,
            "container": self._config.container,
            "blob_name": prepared["blob_name"],
            "sha256": prepared["sha256"],
            "bytes": len(request.data),
            "content_md5_base64": base64.b64encode(prepared["content_md5"]).decode("ascii"),
            "policy_mode": "Locked",
            "retain_until": prepared["retain_until"].isoformat().replace("+00:00", "Z"),
            "legal_hold": request.legal_hold,
            "default_encryption_scope_sha256": hashlib.sha256(
                self._config.default_encryption_scope.encode("utf-8")
            ).hexdigest(),
        }
        receipt.update(extra)
        return receipt


def _utc_seconds(value: datetime) -> datetime:
    if value.tzinfo is None or value.utcoffset() is None:
        raise ValidationError("datetime must be timezone-aware")
    return value.astimezone(timezone.utc).replace(microsecond=0)


def _parse_rfc3339(value: Any) -> datetime:
    if not isinstance(value, str) or not value.endswith("Z"):
        raise ValidationError("approved_at must be RFC3339 UTC with Z")
    try:
        parsed = datetime.fromisoformat(value[:-1] + "+00:00")
    except ValueError as exc:
        raise ValidationError("approved_at is invalid") from exc
    return _utc_seconds(parsed)
````

### FILE: `azure_blob_storage_adapter/azureblob/azure_backend.py`
```yaml
block_id: "PYTHON-AZURE-BLOB-EVIDENCE:azure-backend:v1"
operation: CREATE
provenance: AUTHORED
source: "local glue calling official azure-storage-blob 12.30.0 APIs"
license: "LicenseRef-Workspace-Owner"
sha256: "ece9bdeeb3c4f3d467258f6354c3cae8a2a4c56a96e5dff0d5784a1ed418450e"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

from typing import Any, Mapping

from azure.storage.blob import (
    BlobImmutabilityPolicyMode,
    BlobServiceClient,
    ContentSettings,
    ImmutabilityPolicy,
)

from .storage import CreatedVersion, TargetState, VersionState


class AzureBlobBackend:
    def __init__(self, service_client: BlobServiceClient):
        self._service = service_client

    def inspect_target(self, container: str) -> TargetState:
        props = self._service.get_container_client(container).get_container_properties()
        scope = props.encryption_scope
        return TargetState(
            account_url=self._service.url.rstrip("/"),
            account_name=self._service.account_name,
            container=props.name,
            public_access=props.public_access,
            immutable_storage_with_versioning_enabled=bool(
                props.immutable_storage_with_versioning_enabled
            ),
            default_encryption_scope=(scope.default_encryption_scope if scope else None),
            prevent_encryption_scope_override=bool(
                scope and scope.prevent_encryption_scope_override
            ),
        )

    def create_version(
        self,
        *,
        container: str,
        blob_name: str,
        data: bytes,
        content_type: str,
        content_md5: bytes,
        metadata: Mapping[str, str],
        retain_until,
        legal_hold: bool,
    ) -> CreatedVersion:
        blob = self._service.get_blob_client(container=container, blob=blob_name)
        response = blob.upload_blob(
            data,
            length=len(data),
            overwrite=False,
            validate_content=True,
            content_settings=ContentSettings(content_type=content_type, content_md5=content_md5),
            metadata=dict(metadata),
            immutability_policy=ImmutabilityPolicy(
                expiry_time=retain_until,
                policy_mode=BlobImmutabilityPolicyMode.LOCKED,
            ),
            legal_hold=legal_hold,
        )
        return CreatedVersion(
            version_id=str(response.get("version_id") or ""),
            etag=str(response.get("etag") or ""),
        )

    def inspect_version(self, *, container: str, blob_name: str, version_id: str) -> VersionState:
        props = self._service.get_blob_client(container=container, blob=blob_name).get_blob_properties(
            version_id=version_id
        )
        content_md5 = props.content_settings.content_md5
        policy = props.immutability_policy
        return VersionState(
            version_id=str(props.version_id or ""),
            etag=str(props.etag or ""),
            size=int(props.size),
            metadata=dict(props.metadata or {}),
            content_md5=(bytes(content_md5) if content_md5 is not None else None),
            server_encrypted=bool(props.server_encrypted),
            encryption_scope=props.encryption_scope,
            immutability_policy_mode=(
                str(getattr(policy.policy_mode, "value", policy.policy_mode)) if policy else None
            ),
            immutability_expiry=(policy.expiry_time if policy else None),
            has_legal_hold=bool(props.has_legal_hold),
        )

    def download_version(self, *, container: str, blob_name: str, version_id: str) -> bytes:
        downloader = self._service.get_blob_client(container=container, blob=blob_name).download_blob(
            version_id=version_id,
            validate_content=True,
            decompress=False,
        )
        data: Any = downloader.readall()
        if not isinstance(data, bytes):
            raise TypeError("Azure SDK returned non-bytes for an exact-version download")
        return data
````

### FILE: `azure_blob_storage_adapter/test_storage.py`
```yaml
block_id: "PYTHON-AZURE-BLOB-EVIDENCE:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local unit and official SDK surface contract regressions"
license: "LicenseRef-Workspace-Owner"
sha256: "088f90d6ecc37037560a39f13c9c21c625f5badf08c73d65c128472255ec1853"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import hashlib
import json
import unittest
from dataclasses import replace
from datetime import datetime, timedelta, timezone
from types import SimpleNamespace

from azure.storage.blob import BlobImmutabilityPolicyMode

from azureblob.azure_backend import AzureBlobBackend
from azureblob.storage import (
    Approvals,
    CreatedVersion,
    EvidenceStore,
    PartialVerificationError,
    StorageConfig,
    StoreRequest,
    TargetState,
    UnknownEffectError,
    ValidationError,
    VersionState,
)


NOW = datetime(2026, 8, 27, 12, 0, tzinfo=timezone.utc)
DATA = b"invoice-evidence-v1"
SCOPE = "elite-documents-cmk"
ACCOUNT = "elitestore001"
URL = f"https://{ACCOUNT}.blob.core.windows.net"
RESOURCE_ID = (
    "/subscriptions/00000000-0000-0000-0000-000000000001/"
    "resourceGroups/elite-rg/providers/Microsoft.Storage/storageAccounts/elitestore001"
)


def authority_receipt(**changes):
    value = {
        "schema_version": "azure-storage-authority-receipt/v1",
        "account_url": URL,
        "account_name": ACCOUNT,
        "storage_account_resource_id": RESOURCE_ID,
        "container": "evidence",
        "immutable_storage_with_versioning_enabled": True,
        "default_encryption_scope": SCOPE,
        "prevent_encryption_scope_override": True,
        "approval_reference": "change-2026-08-27-001",
        "approved_at": "2026-08-27T11:55:00Z",
    }
    value.update(changes)
    return value


def receipt_hash(receipt):
    canonical = json.dumps(receipt, sort_keys=True, separators=(",", ":"), ensure_ascii=True).encode()
    return hashlib.sha256(canonical).hexdigest()


def config(receipt=None, **changes):
    auth = receipt or authority_receipt()
    value = StorageConfig(
        account_url=URL,
        account_name=ACCOUNT,
        storage_account_resource_id=RESOURCE_ID,
        authority_receipt_sha256=receipt_hash(auth),
        authority_receipt_max_age_seconds=3600,
        container="evidence",
        default_encryption_scope=SCOPE,
        max_bytes=1024 * 1024,
        allowed_content_types=frozenset({"application/pdf", "image/png"}),
        legal_hold_allowed=False,
    )
    return replace(value, **changes)


def request(**changes):
    value = StoreRequest(
        tenant_id="tenant-1",
        document_id="invoice-1",
        data=DATA,
        content_type="application/pdf",
        retain_until=NOW + timedelta(days=30),
        legal_hold=False,
    )
    return replace(value, **changes)


APPROVED = Approvals(access=True, cost=True, retention=True, irreversible_lock=True)


class FakeBackend:
    def __init__(self):
        self.calls = []
        self.target = TargetState(URL, ACCOUNT, "evidence", None, True, SCOPE, True)
        self.create_error = None
        self.downloaded = DATA
        digest = hashlib.sha256(DATA).hexdigest()
        self.state = VersionState(
            version_id="2026-08-27T12:00:01.0000000Z",
            etag='"etag-1"',
            size=len(DATA),
            metadata={
                "sha256": digest,
                "bytes": str(len(DATA)),
                "authorityreceiptsha256": receipt_hash(authority_receipt()),
            },
            content_md5=hashlib.md5(DATA, usedforsecurity=False).digest(),
            server_encrypted=True,
            encryption_scope=SCOPE,
            immutability_policy_mode="Locked",
            immutability_expiry=NOW + timedelta(days=30),
            has_legal_hold=False,
        )

    def inspect_target(self, container):
        self.calls.append(("inspect_target", container))
        return self.target

    def create_version(self, **kwargs):
        self.calls.append(("create_version", kwargs))
        if self.create_error:
            raise self.create_error
        return CreatedVersion(self.state.version_id, self.state.etag)

    def inspect_version(self, **kwargs):
        self.calls.append(("inspect_version", kwargs))
        return self.state

    def download_version(self, **kwargs):
        self.calls.append(("download_version", kwargs))
        return self.downloaded


class EvidenceStoreTests(unittest.TestCase):
    def test_verified_create_uses_deterministic_name_and_exact_download(self):
        backend = FakeBackend()
        receipt = EvidenceStore(config(), backend).put(request(), APPROVED, authority_receipt(), now=NOW)
        digest = hashlib.sha256(DATA).hexdigest()
        self.assertEqual("VERIFIED", receipt["status"])
        self.assertEqual(f"evidence/tenant-1/invoice-1/{digest}", receipt["blob_name"])
        self.assertTrue(receipt["retention_verified"])
        self.assertTrue(receipt["bytes_verified_by_download"])
        create = next(call for name, call in backend.calls if name == "create_version")
        self.assertEqual(hashlib.md5(DATA, usedforsecurity=False).digest(), create["content_md5"])

    def test_missing_approvals_fail_before_provider(self):
        backend = FakeBackend()
        with self.assertRaises(ValidationError):
            EvidenceStore(config(), backend).put(request(), Approvals(), authority_receipt(), now=NOW)
        self.assertEqual([], backend.calls)

    def test_stale_authority_receipt_fails_before_provider(self):
        auth = authority_receipt(approved_at="2026-08-26T11:00:00Z")
        backend = FakeBackend()
        with self.assertRaisesRegex(ValidationError, "stale"):
            EvidenceStore(config(auth), backend).put(request(), APPROVED, auth, now=NOW)
        self.assertEqual([], backend.calls)

    def test_authority_receipt_is_closed_and_hash_bound(self):
        auth = authority_receipt()
        cfg = config(auth)
        auth["unexpected"] = True
        with self.assertRaisesRegex(ValidationError, "not closed"):
            EvidenceStore(cfg, FakeBackend()).put(request(), APPROVED, auth, now=NOW)

    def test_target_requires_private_versioned_immutable_container(self):
        backend = FakeBackend()
        backend.target = replace(backend.target, immutable_storage_with_versioning_enabled=False)
        with self.assertRaisesRegex(ValidationError, "versioning"):
            EvidenceStore(config(), backend).put(request(), APPROVED, authority_receipt(), now=NOW)
        self.assertFalse(any(name == "create_version" for name, _ in backend.calls))

    def test_target_requires_exact_non_overridable_encryption_scope(self):
        backend = FakeBackend()
        backend.target = replace(backend.target, prevent_encryption_scope_override=False)
        with self.assertRaisesRegex(ValidationError, "override"):
            EvidenceStore(config(), backend).put(request(), APPROVED, authority_receipt(), now=NOW)

    def test_locked_retention_has_official_bounds(self):
        for delta in (timedelta(hours=23), timedelta(days=146001)):
            with self.subTest(delta=delta):
                with self.assertRaisesRegex(ValidationError, "146000"):
                    EvidenceStore(config(), FakeBackend()).put(
                        request(retain_until=NOW + delta), APPROVED, authority_receipt(), now=NOW
                    )

    def test_legal_hold_needs_profile_and_approval(self):
        with self.assertRaisesRegex(ValidationError, "legal hold"):
            EvidenceStore(config(), FakeBackend()).put(
                request(legal_hold=True), APPROVED, authority_receipt(), now=NOW
            )

    def test_create_exception_is_unknown_effect_not_blind_retry(self):
        backend = FakeBackend()
        backend.create_error = TimeoutError("simulated response loss")
        with self.assertRaises(UnknownEffectError) as caught:
            EvidenceStore(config(), backend).put(request(), APPROVED, authority_receipt(), now=NOW)
        self.assertEqual("UNKNOWN_EFFECT", caught.exception.receipt["status"])
        self.assertIn("blob_name", caught.exception.receipt)
        self.assertNotIn("simulated response loss", json.dumps(caught.exception.receipt))

    def test_property_mismatch_returns_partial_receipt(self):
        backend = FakeBackend()
        backend.state = replace(backend.state, encryption_scope="wrong")
        with self.assertRaises(PartialVerificationError) as caught:
            EvidenceStore(config(), backend).put(request(), APPROVED, authority_receipt(), now=NOW)
        self.assertEqual("PARTIAL_VERIFICATION", caught.exception.receipt["status"])
        self.assertEqual(backend.state.version_id, caught.exception.receipt["version_id"])

    def test_download_hash_mismatch_returns_partial_receipt(self):
        backend = FakeBackend()
        backend.downloaded = b"changed"
        with self.assertRaises(PartialVerificationError):
            EvidenceStore(config(), backend).put(request(), APPROVED, authority_receipt(), now=NOW)

    def test_reconcile_verifies_supplied_exact_version_without_create(self):
        backend = FakeBackend()
        receipt = EvidenceStore(config(), backend).reconcile(
            request(), backend.state.version_id, APPROVED, authority_receipt(), now=NOW
        )
        self.assertEqual("VERIFIED", receipt["status"])
        self.assertFalse(any(name == "create_version" for name, _ in backend.calls))


class FakeDownload:
    def __init__(self, data):
        self._data = data

    def readall(self):
        return self._data


class FakeBlobClient:
    def __init__(self):
        self.upload = None
        self.property_version = None
        self.download = None

    def upload_blob(self, data, **kwargs):
        self.upload = (data, kwargs)
        return {"version_id": "v1", "etag": '"e1"'}

    def get_blob_properties(self, **kwargs):
        self.property_version = kwargs["version_id"]
        return SimpleNamespace(
            version_id="v1",
            etag='"e1"',
            size=len(DATA),
            metadata={"sha256": hashlib.sha256(DATA).hexdigest()},
            content_settings=SimpleNamespace(content_md5=bytearray(hashlib.md5(DATA, usedforsecurity=False).digest())),
            server_encrypted=True,
            encryption_scope=SCOPE,
            immutability_policy=SimpleNamespace(
                policy_mode=BlobImmutabilityPolicyMode.LOCKED,
                expiry_time=NOW + timedelta(days=30),
            ),
            has_legal_hold=True,
        )

    def download_blob(self, **kwargs):
        self.download = kwargs
        return FakeDownload(DATA)


class FakeContainerClient:
    def get_container_properties(self):
        return SimpleNamespace(
            name="evidence",
            public_access=None,
            immutable_storage_with_versioning_enabled=True,
            encryption_scope=SimpleNamespace(
                default_encryption_scope=SCOPE,
                prevent_encryption_scope_override=True,
            ),
        )


class FakeServiceClient:
    url = URL
    account_name = ACCOUNT

    def __init__(self):
        self.blob = FakeBlobClient()

    def get_container_client(self, container):
        self.container = container
        return FakeContainerClient()

    def get_blob_client(self, **kwargs):
        self.blob_request = kwargs
        return self.blob


class AzureBackendContractTests(unittest.TestCase):
    def test_official_sdk_contract_is_create_only_locked_and_checksummed(self):
        service = FakeServiceClient()
        backend = AzureBlobBackend(service)
        created = backend.create_version(
            container="evidence",
            blob_name="evidence/t/d/hash",
            data=DATA,
            content_type="application/pdf",
            content_md5=hashlib.md5(DATA, usedforsecurity=False).digest(),
            metadata={"sha256": hashlib.sha256(DATA).hexdigest()},
            retain_until=NOW + timedelta(days=30),
            legal_hold=True,
        )
        self.assertEqual("v1", created.version_id)
        _, kwargs = service.blob.upload
        self.assertIs(kwargs["overwrite"], False)
        self.assertIs(kwargs["validate_content"], True)
        self.assertEqual(BlobImmutabilityPolicyMode.LOCKED, kwargs["immutability_policy"].policy_mode)
        self.assertIs(kwargs["legal_hold"], True)
        self.assertEqual(hashlib.md5(DATA, usedforsecurity=False).digest(), kwargs["content_settings"].content_md5)

    def test_official_sdk_contract_reads_properties_and_bytes_by_exact_version(self):
        service = FakeServiceClient()
        backend = AzureBlobBackend(service)
        state = backend.inspect_version(container="evidence", blob_name="x", version_id="v1")
        data = backend.download_version(container="evidence", blob_name="x", version_id="v1")
        self.assertEqual("v1", state.version_id)
        self.assertEqual("Locked", state.immutability_policy_mode)
        self.assertEqual("v1", service.blob.property_version)
        self.assertEqual(DATA, data)
        self.assertEqual("v1", service.blob.download["version_id"])
        self.assertIs(service.blob.download["validate_content"], True)
        self.assertIs(service.blob.download["decompress"], False)

    def test_official_sdk_contract_reads_versioned_immutable_container_state(self):
        target = AzureBlobBackend(FakeServiceClient()).inspect_target("evidence")
        self.assertTrue(target.immutable_storage_with_versioning_enabled)
        self.assertTrue(target.prevent_encryption_scope_override)
        self.assertEqual(SCOPE, target.default_encryption_scope)


if __name__ == "__main__":
    unittest.main()
````

### FILE: `azure_blob_storage_adapter/README.md`
```yaml
block_id: "PYTHON-AZURE-BLOB-EVIDENCE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local usage and non-claim boundary"
license: "LicenseRef-Workspace-Owner"
sha256: "a07c02e4564ca8ff8ae7c65f4547871eabf4de9717d3244088e72cb04c645b51"
variables: []
secrets_allowed: false
```
````markdown
# Azure Blob immutable evidence adapter

This directory is local `AUTHORED` integration glue that calls Microsoft's official `azure-storage-blob==12.30.0` SDK. It is not Microsoft-authored application code.

The adapter admits only a private container that reports immutable storage with versioning enabled and an exact, non-overridable default encryption scope. It requires a fresh, hash-bound authority receipt for the exact Azure account resource ID and explicit access, cost, retention and irreversible-lock approvals. It uploads a deterministic Block Blob name with `overwrite=False`, wire validation, stored Content-MD5, SHA-256 metadata, a `Locked` immutability policy and an optional separately approved legal hold in the same SDK call.

Success is not inferred from the upload response. The adapter reads properties and downloads the exact returned `version_id`, then checks ETag, length, SHA-256 metadata, Content-MD5, server encryption, encryption scope, locked retention expiry, legal hold and the downloaded bytes' SHA-256. A response-loss error produces an `UNKNOWN_EFFECT` receipt; a post-write mismatch produces a `PARTIAL_VERIFICATION` receipt. Both preserve the deterministic blob name and prohibit blind retry. `reconcile` verifies an explicitly supplied version without creating another object.

The adapter never creates accounts or containers, enables/migrates versioning, enables/locks/changes immutability policies, changes RBAC or networking, removes legal holds, deletes versions, stores credentials or authorizes storage merely because document extraction returned fields. Microsoft documents version-level WORM enablement/migration as an administrative action that can be irreversible; it remains outside this runtime.

From a materialized directory with CPython 3.12 and the exact downloaded wheels:

```text
python -m pip install --no-index --find-links <wheelhouse> --require-hashes -r requirements-windows-py312.lock
python -m unittest -v
python -m compileall -q azureblob test_storage.py
```

Offline tests do not prove Azure tenant/subscription authority, RBAC, policy locks, cost, network controls, target immutability, audit logs, concurrency, recovery or legal compliance. A project must demonstrate those gates in its own sandbox before setting `automatic_storage_authorized=true`.
````

## 6. Configuration surface

| Variable | Tipo | Default seguro | Validación | Secreto | Efecto |
|---|---|---|---|---|---|
| account URL/name/resource ID | string | vacío/bloqueado | HTTPS y ARM ID exactos; receipt y target deben coincidir | no | autoridad Azure |
| authority receipt/hash/max age | JSON/SHA/segundos | vacío/bloqueado | schema cerrado, SHA exacto, 60..86400 s | no | evidencia humana fresca |
| container | string | vacío/bloqueado | nombre Azure; privado y version-level WORM observado | no | destino |
| default encryption scope | string | vacío/bloqueado | scope exacto y override prohibido por contenedor | no | cifrado servidor |
| max bytes/MIME | int/set | 0/vacío | positivo y allowlist antes de provider | no | frontera de input |
| retain until | UTC datetime | ausente | 1..146000 días; microsegundos normalizados | no | WORM Locked irreversible |
| legal hold | boolean | false | perfil más aprobación separada | no | retención indefinida |
| access/cost/retention/lock approvals | booleans | false | todos true pre-provider | no | autoridad humana |
| Azure credential | TokenCredential externo | ausente | inyectado al SDK; nunca serializado | sí | autenticación |

## 7. Dependency bill

| Package/tool | Pin exacto | Uso | Licencia | Fuente oficial |
|---|---|---|---|---|
| azure-storage-blob | `12.30.0`, tag `azure-storage-blob_12.30.0`, commit `245d2024…` | Blob/version/WORM | MIT | Microsoft GitHub + PyPI |
| azure-core | `1.41.0` | pipeline SDK | MIT | PyPI/Microsoft |
| cryptography | `50.0.1` | dependencia SDK | Apache-2.0 OR BSD-3-Clause | PyPI |
| CPython | `3.12.13 windows/amd64` | runtime verificado | PSF-2.0 | Python/uv managed runtime |
| grafo Python | 12 wheels con SHA-256 | runtime | licencias individuales | `requirements-windows-py312.lock` |

El tag exacto es unsigned según GitHub; el lock conserva esa condición y la compensa con tag→commit, archive, wheel, sdist, licencia y comparación byte a byte, sin afirmar firma Microsoft. Una actualización abre un expediente nuevo y repite identidad, licencia, source compare, lock, OSV y tests.

## 8. Apply order

1. Completar blueprint, decisión documental Azure y approvals; no auto-seleccionar.
2. Adquirir archive/wheel/sdist/licencia exactos, verificar SHA-256 y lock.
3. Materializar los nueve archivos en destino vacío o sin colisiones.
4. Instalar los doce wheels con `--require-hashes`; inyectar un `BlobServiceClient` con credencial externa.
5. Completar receipt de autoridad y profile; mantener `automatic_storage_authorized=false`.
6. Ejecutar compile/tests/OSV y sandbox Azure sobre recurso aprobado; comprobar timeout/efecto parcial/audit/costo/concurrencia.
7. Habilitar por clase sólo después de corpus, seguridad, recovery y política legal. Rollback detiene nuevas escrituras y reconcilia versiones; nunca afloja o borra retención.

## 9. Verification

- artefactos oficiales: wheel/sdist hashes PyPI; 79/79 archivos Python wheel↔sdist y 193 archivos seleccionados sdist↔commit byte-idénticos;
- `pip install --require-hashes` y `pip check` sobre CPython 3.12.13 → PASS;
- `python -m unittest -v` → 15 pruebas PASS, incluidas precondiciones, SDK exacto, create-only, Locked, checksum, version readback, mismatch, efecto desconocido y reconcile;
- `python -m compileall -q azureblob test_storage.py` → PASS;
- OSV batch oficial fechado sobre 12 distribuciones runtime → 0 findings conocidos;
- materializador/compositor/verificador global deben confirmar manifest, SHA, profile y ejecución desde Markdown;
- pendiente del proyecto: llamada Azure real, RBAC/red/scope, política bloqueada, costo, logs, carga, restore y compliance. El pack permanece `CONDITIONED`.

## 10. Reconstruction evidence

Reconstrucción limpia Windows 11 / PowerShell 7 / CPython 3.12.13 el 2026-08-27. PyPI `12.30.0`: wheel 435610 bytes SHA-256 `d415ac50…`, sdist 618229 bytes SHA-256 `2cd74d4d…`; tag oficial resuelve al commit `245d2024…`, archive 205202391 bytes SHA-256 `7c4743cc…`, licencia MIT SHA-256 `fd532481…`, tag unsigned retenido. Evidencia canónica: `reconstruction_evidence/PYTHON_AZURE_BLOB_IMMUTABLE_EVIDENCE_ADAPTER_2026-08-27_V1.md`. No existe claim de llamada cloud ni autorización automática.
