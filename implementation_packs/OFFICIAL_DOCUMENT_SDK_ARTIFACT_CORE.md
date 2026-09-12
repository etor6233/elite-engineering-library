# Official Document SDK Artifact Core

## 1. Metadata

```yaml
pack_id: "OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE"
pack_version: "0.3.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un lock ejecutable y un adquiridor fail-closed para catorce artefactos oficiales exactos de SDKs documentales Microsoft, Google, OpenAI y Oracle: wheel+source estables de Azure Document Intelligence, Azure Content Understanding, Google Document AI, Google Cloud Storage, OpenAI y Oracle OCI, más wheel+sdist de Google Document AI Toolbox sólo como candidato condicionado, con aprobación enlazada al hash, verificación offline/network y recibo reproducible."
stacks: ["PowerShell 7", "Python 3.10+"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "DEPENDENCY-LICENSE-EVIDENCE-CORE 0.1.x"]
incompatible_with: ["versiones móviles", "artefactos sin hash", "credenciales embebidas", "instalación transitive sin lock del target"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://pypi.org/project/azure-ai-documentintelligence/1.0.2/", "https://pypi.org/project/azure-ai-contentunderstanding/1.1.0/", "https://pypi.org/project/google-cloud-documentai/3.15.0/", "https://pypi.org/project/google-cloud-storage/3.13.1/", "https://pypi.org/project/google-cloud-documentai-toolbox/0.17.3/", "https://pypi.org/project/openai/3.3.1/", "https://pypi.org/project/oci/2.185.0/", "https://github.com/oracle/oci-python-sdk/tree/e988c91dcc9963718454cb1e215bd44540524881"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use este pack cuando el blueprint haya seleccionado uno o más SDKs documentales Python oficiales y necesite adquirir los artefactos exactos sin Git. Requiere PowerShell 7, una selección explícita, aceptación de las licencias y reconocimiento de que cuentas, cuotas y llamadas de proveedor pueden tener costo. Rechácelo si el proyecto no ha elegido proveedor, si necesita otra versión/plataforma o si pretende tratar el artefacto como prueba de exactitud del negocio.

El adquiridor y el probe son `AUTHORED`; los artefactos descargados no se redistribuyen dentro del pack y conservan propietario/licencia upstream. El pack prueba identidad e interfaces importables, no calidad de extracción, seguridad productiva ni disponibilidad de una cuenta.

## 3. Architecture contract

`artifact-lock.json` es la única autoridad de selección. El runner valida esquema exacto, IDs únicos, host oficial `files.pythonhosted.org`, nombre derivado de URL, tamaño y SHA-256 antes de aceptar el lock. Fuera de `ValidateOnly`, exige un approval con el SHA del lock y el conjunto exacto de IDs; descarga o importa desde cache a staging, verifica bytes/hash, escribe un recibo y mueve atómicamente a un destino inexistente.

Un error de schema, approval, selección, red, tamaño o hash falla cerrado y elimina staging. El destino final no aparece ante una adquisición fallida. El probe comprueba versiones y símbolos públicos después de que el proyecto construya su virtual environment y su lock transitivo por plataforma. No procesa documentos, no recibe secretos y no llama endpoints.

Presupuesto: validación local menor a 2 segundos; adquisición limitada por bytes/red del artefacto; import cache limitada por I/O local. Rollback: borrar sólo el directorio nuevo de artefactos después de preservar receipt/evidencia. Actualizaciones se hacen mediante un nuevo lock y expediente de dependency update, nunca mutando una versión aprobada en sitio.

## 4. Exact file manifest

```text
CREATE official_document_sdks/artifact-lock.json
CREATE official_document_sdks/artifact-approval.template.json
CREATE official_document_sdks/acquire_document_sdk_artifacts.ps1
CREATE official_document_sdks/probe_document_sdks.py
CREATE official_document_sdks/test_artifact_acquisition.ps1
CREATE official_document_sdks/README.md
```

## 5. Materialization blocks

### FILE: `official_document_sdks/artifact-lock.json`
```yaml
block_id: "OFFICIAL-DOCUMENT-SDK-ARTIFACT:lock:v1"
operation: CREATE
provenance: AUTHORED
source: "local registry derived from exact official PyPI artifacts"
license: "LicenseRef-Workspace-Owner"
sha256: "5a3fd625518a20c8f9fee154b924552f9877a0b7fb7fcf209887215f2cfecfe6"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-official-document-sdk-artifact-lock/v1",
  "verified_at": "2026-08-28",
  "artifacts": [
    {
      "id": "azure-ai-documentintelligence-1.0.2-wheel",
      "owner": "Microsoft Azure",
      "classification": "INTEGRATION_ONLY",
      "package": "azure-ai-documentintelligence",
      "version": "1.0.2",
      "python_requires": ">=3.8",
      "filename": "azure_ai_documentintelligence-1.0.2-py3-none-any.whl",
      "kind": "wheel",
      "url": "https://files.pythonhosted.org/packages/d9/75/c9ec040f23082f54ffb1977ff8f364c2d21c79a640a13d1c1809e7fd6b1a/azure_ai_documentintelligence-1.0.2-py3-none-any.whl",
      "bytes": 106005,
      "sha256": "e1fb446abbdeccc9759d897898a0fe13141ed29f9ad11fc705f951925822ed59",
      "license_expression": "MIT"
    },
    {
      "id": "azure-ai-documentintelligence-1.0.2-sdist",
      "owner": "Microsoft Azure",
      "classification": "INTEGRATION_ONLY",
      "package": "azure-ai-documentintelligence",
      "version": "1.0.2",
      "python_requires": ">=3.8",
      "filename": "azure_ai_documentintelligence-1.0.2.tar.gz",
      "kind": "sdist",
      "url": "https://files.pythonhosted.org/packages/44/7b/8115cd713e2caa5e44def85f2b7ebd02a74ae74d7113ba20bdd41fd6dd80/azure_ai_documentintelligence-1.0.2.tar.gz",
      "bytes": 170940,
      "sha256": "4d75a2513f2839365ebabc0e0e1772f5601b3a8c9a71e75da12440da13b63484",
      "license_expression": "MIT"
    },
    {
      "id": "azure-ai-contentunderstanding-1.1.0-wheel",
      "owner": "Microsoft Azure",
      "classification": "INTEGRATION_ONLY",
      "package": "azure-ai-contentunderstanding",
      "version": "1.1.0",
      "python_requires": ">=3.9",
      "filename": "azure_ai_contentunderstanding-1.1.0-py3-none-any.whl",
      "kind": "wheel",
      "url": "https://files.pythonhosted.org/packages/67/d6/6f60ac39b73a0b8a2bc86cd6895e26eccc20d9b7bde86814fe6a36964dff/azure_ai_contentunderstanding-1.1.0-py3-none-any.whl",
      "bytes": 101987,
      "sha256": "d1d6bdeffe02f5c8cc5309f0e120a1338f51af04638605f9617b7b501e2d1dd6",
      "license_expression": "MIT"
    },
    {
      "id": "azure-ai-contentunderstanding-1.1.0-sdist",
      "owner": "Microsoft Azure",
      "classification": "PINNED_CANDIDATE",
      "package": "azure-ai-contentunderstanding",
      "version": "1.1.0",
      "python_requires": ">=3.9",
      "filename": "azure_ai_contentunderstanding-1.1.0.tar.gz",
      "kind": "sdist",
      "url": "https://files.pythonhosted.org/packages/b6/5a/6dbcb8278c8d3779fc03992d8a9f87a9ce5280437d76e72f8121f9a90d46/azure_ai_contentunderstanding-1.1.0.tar.gz",
      "bytes": 230330,
      "sha256": "00f39c7cf2ba50c8d4586315f4295a45fad0881daaff07d13d33aea2bafc7ec8",
      "license_expression": "MIT"
    },
    {
      "id": "google-cloud-documentai-3.15.0-wheel",
      "owner": "Google Cloud",
      "classification": "INTEGRATION_ONLY",
      "package": "google-cloud-documentai",
      "version": "3.15.0",
      "python_requires": ">=3.10",
      "filename": "google_cloud_documentai-3.15.0-py3-none-any.whl",
      "kind": "wheel",
      "url": "https://files.pythonhosted.org/packages/55/d1/2a873f97cb08bb592f5b944e39f040161d5f7d2d4edbfae4c7515869b12a/google_cloud_documentai-3.15.0-py3-none-any.whl",
      "bytes": 310694,
      "sha256": "f040f4f9db43411184197a808b11fde52b580723a9fca336f43ea7c6a885bfd8",
      "license_expression": "Apache-2.0"
    },
    {
      "id": "google-cloud-documentai-3.15.0-sdist",
      "owner": "Google Cloud",
      "classification": "INTEGRATION_ONLY",
      "package": "google-cloud-documentai",
      "version": "3.15.0",
      "python_requires": ">=3.10",
      "filename": "google_cloud_documentai-3.15.0.tar.gz",
      "kind": "sdist",
      "url": "https://files.pythonhosted.org/packages/e0/63/a775dd454e0fa523b999780665f3117c0026eaa604dd7c47ad3c45c7624f/google_cloud_documentai-3.15.0.tar.gz",
      "bytes": 362736,
      "sha256": "d50b69a8a62aaf803b0d926f93b6579df440f5f779f45fec2c5972768a1cdbbd",
      "license_expression": "Apache-2.0"
    },
    {
      "id": "google-cloud-storage-3.13.1-wheel",
      "owner": "Google Cloud",
      "classification": "INTEGRATION_ONLY",
      "package": "google-cloud-storage",
      "version": "3.13.1",
      "python_requires": ">=3.10",
      "filename": "google_cloud_storage-3.13.1-py3-none-any.whl",
      "kind": "wheel",
      "url": "https://files.pythonhosted.org/packages/06/6f/d69f0e185e08ddb58c323a0a935af2b492907b5de362bc08933b0a3b5644/google_cloud_storage-3.13.1-py3-none-any.whl",
      "bytes": 341486,
      "sha256": "98208de6c21e85cecd3eb44551894efff33d98365500e178867d4305854a770a",
      "license_expression": "Apache-2.0"
    },
    {
      "id": "google-cloud-storage-3.13.1-sdist",
      "owner": "Google Cloud",
      "classification": "INTEGRATION_ONLY",
      "package": "google-cloud-storage",
      "version": "3.13.1",
      "python_requires": ">=3.10",
      "filename": "google_cloud_storage-3.13.1.tar.gz",
      "kind": "sdist",
      "url": "https://files.pythonhosted.org/packages/ce/7e/73bb7512df1d1aad6ce3f9aed847cd40e0cd400ba4a85d86ab8eb412e9cc/google_cloud_storage-3.13.1.tar.gz",
      "bytes": 17341051,
      "sha256": "a80bf8cac2794808aa61c50c5f769ecbbe2d10331bacd0d69d30e59b14b346b2",
      "license_expression": "Apache-2.0"
    },
    {
      "id": "openai-3.3.1-wheel",
      "owner": "OpenAI",
      "classification": "INTEGRATION_ONLY",
      "package": "openai",
      "version": "3.3.1",
      "python_requires": ">=3.10",
      "filename": "openai-3.3.1-py3-none-any.whl",
      "kind": "wheel",
      "url": "https://files.pythonhosted.org/packages/6a/db/2b7a1b3de659bb82aef979116c74e809982b13e42c057759767552b5155f/openai-3.3.1-py3-none-any.whl",
      "bytes": 1690337,
      "sha256": "9652df7fdf8ee6f5bd58e0a12f2b1d414a18e0f06bb7a9a57c8643a5f5469bd3",
      "license_expression": "Apache-2.0"
    },
    {
      "id": "openai-3.3.1-sdist",
      "owner": "OpenAI",
      "classification": "INTEGRATION_ONLY",
      "package": "openai",
      "version": "3.3.1",
      "python_requires": ">=3.10",
      "filename": "openai-3.3.1.tar.gz",
      "kind": "sdist",
      "url": "https://files.pythonhosted.org/packages/7d/9c/ba0c292b4032ede74c249ca314ad64eb1bb5a03a843f6e01facb02f80cd8/openai-3.3.1.tar.gz",
      "bytes": 1282113,
      "sha256": "6f22807de1a976c932cecda620e8172a8c3fdbaeed29c7f21564e0c2410edf56",
      "license_expression": "Apache-2.0"
    },
    {
      "id": "oci-2.185.0-wheel",
      "owner": "Oracle",
      "classification": "INTEGRATION_ONLY",
      "package": "oci",
      "version": "2.185.0",
      "python_requires": "3.12 verified; package metadata omits Requires-Python",
      "filename": "oci-2.185.0-py3-none-any.whl",
      "kind": "wheel",
      "url": "https://files.pythonhosted.org/packages/b8/75/0f9264cac734696c209117421fd35d373aaeae7a57ea43796b135281e65d/oci-2.185.0-py3-none-any.whl",
      "bytes": 36716800,
      "sha256": "93556f08c6270e2e174d59d94fcd49fa794ba3f4d60630af173bc9b6204b085c",
      "license_expression": "UPL-1.0 OR Apache-2.0"
    },
    {
      "id": "oci-2.185.0-sdist",
      "owner": "Oracle",
      "classification": "INTEGRATION_ONLY",
      "package": "oci",
      "version": "2.185.0",
      "python_requires": "3.12 verified; package metadata omits Requires-Python",
      "filename": "oci-2.185.0.tar.gz",
      "kind": "sdist",
      "url": "https://files.pythonhosted.org/packages/a5/fd/5408374a5c732882fc7b6bc1b6e6a96f06984a325ed2b66d3a1f67cf8b8b/oci-2.185.0.tar.gz",
      "bytes": 18098448,
      "sha256": "375eec7c8568b0066bd5e6ff48cf406a4a1af34437c29abd7981f1ab78fac9c2",
      "license_expression": "UPL-1.0 OR Apache-2.0"
    },
    {
      "id": "google-cloud-documentai-toolbox-0.17.3-wheel",
      "owner": "Google Cloud",
      "classification": "PINNED_CANDIDATE",
      "package": "google-cloud-documentai-toolbox",
      "version": "0.17.3",
      "python_requires": ">=3.10",
      "filename": "google_cloud_documentai_toolbox-0.17.3-py3-none-any.whl",
      "kind": "wheel",
      "url": "https://files.pythonhosted.org/packages/66/89/54bda8e808c6dc09082997246645d95431146161c265b02496ff9e717d11/google_cloud_documentai_toolbox-0.17.3-py3-none-any.whl",
      "bytes": 43670,
      "sha256": "46f37a42ab98a1255ea078fe4ce85084cc7ecf127b6f5ccc3c546599262da014",
      "license_expression": "Apache-2.0"
    },
    {
      "id": "google-cloud-documentai-toolbox-0.17.3-sdist",
      "owner": "Google Cloud",
      "classification": "PINNED_CANDIDATE",
      "package": "google-cloud-documentai-toolbox",
      "version": "0.17.3",
      "python_requires": ">=3.10",
      "filename": "google_cloud_documentai_toolbox-0.17.3.tar.gz",
      "kind": "sdist",
      "url": "https://files.pythonhosted.org/packages/12/7d/6d255f5d8360a695e922987733fcaf84b59b9b9a7038cadf8e19c198134a/google_cloud_documentai_toolbox-0.17.3.tar.gz",
      "bytes": 24314548,
      "sha256": "4e19e04abc9fcb4d3b131697b91c3486d03c38cbf01d06581f2bda43050f0f46",
      "license_expression": "Apache-2.0"
    }
  ]
}
````

### FILE: `official_document_sdks/artifact-approval.template.json`
```yaml
block_id: "OFFICIAL-DOCUMENT-SDK-ARTIFACT:approval:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "91d770c8487761ec3001b372540773fea00a8e9567bb67c67b6cf76fdd83f92a"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-official-document-sdk-artifact-approval/v1",
  "lock_sha256": "REPLACE_WITH_LOCK_SHA256",
  "approved_artifact_ids": [],
  "approved_by": "",
  "approved_at": "YYYY-MM-DD",
  "accept_artifact_licenses": false,
  "acknowledge_provider_accounts_quotas_and_costs": false
}
````

### FILE: `official_document_sdks/acquire_document_sdk_artifacts.ps1`
```yaml
block_id: "OFFICIAL-DOCUMENT-SDK-ARTIFACT:acquire:v1"
operation: CREATE
provenance: AUTHORED
source: "local verification wrapper around official immutable PyPI artifacts"
license: "LicenseRef-Workspace-Owner"
sha256: "8f2b8a6831ea39c882bdb54ee860524e1ae3506933fe85f1f183f0d456745657"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param(
  [string]$LockPath = (Join-Path $PSScriptRoot 'artifact-lock.json'),
  [Parameter(Mandatory = $true)][string[]]$ArtifactId,
  [string]$Destination,
  [string]$ApprovalPath,
  [string]$ImportDirectory,
  [switch]$ValidateOnly
)

$ErrorActionPreference = 'Stop'

function Fail([string]$Message) { throw "DOCUMENT_SDK_ACQUISITION_FAILED: $Message" }
function Sha256([string]$Path) { (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant() }
function ExactProperties($Object, [string[]]$Expected, [string]$Context) {
  $actual = @($Object.PSObject.Properties.Name | Sort-Object)
  $wanted = @($Expected | Sort-Object)
  if (($actual -join '|') -cne ($wanted -join '|')) { Fail "$Context properties differ: actual=$($actual -join ',')" }
}

if (-not (Test-Path -LiteralPath $LockPath -PathType Leaf)) { Fail "lock not found: $LockPath" }
$lock = Get-Content -LiteralPath $LockPath -Raw | ConvertFrom-Json -Depth 20
ExactProperties $lock @('schema','verified_at','artifacts') 'lock'
if ($lock.schema -cne 'elite-official-document-sdk-artifact-lock/v1') { Fail "unsupported lock schema" }
$verified = [datetime]::MinValue
if (-not [datetime]::TryParseExact([string]$lock.verified_at, 'yyyy-MM-dd', [Globalization.CultureInfo]::InvariantCulture, [Globalization.DateTimeStyles]::None, [ref]$verified)) { Fail 'verified_at is not an ISO calendar date' }
if (@($lock.artifacts).Count -eq 0) { Fail 'lock has no artifacts' }

$byId = @{}
foreach ($artifact in @($lock.artifacts)) {
  ExactProperties $artifact @('id','owner','classification','package','version','python_requires','filename','kind','url','bytes','sha256','license_expression') "artifact $($artifact.id)"
  if ([string]::IsNullOrWhiteSpace($artifact.id) -or $artifact.id -cnotmatch '^[a-z0-9][a-z0-9.-]+$') { Fail "invalid artifact id: $($artifact.id)" }
  if ($byId.ContainsKey($artifact.id)) { Fail "duplicate artifact id: $($artifact.id)" }
  if ($artifact.classification -cnotin @('INTEGRATION_ONLY','PINNED_CANDIDATE')) { Fail "invalid classification: $($artifact.id)" }
  if ($artifact.kind -cnotin @('wheel','sdist')) { Fail "invalid kind: $($artifact.id)" }
  if ([int64]$artifact.bytes -le 0) { Fail "invalid bytes: $($artifact.id)" }
  if ($artifact.sha256 -cnotmatch '^[0-9a-f]{64}$') { Fail "invalid sha256: $($artifact.id)" }
  $uri = [uri]$artifact.url
  if ($uri.Scheme -cne 'https' -or $uri.Host -cne 'files.pythonhosted.org') { Fail "non-official artifact host: $($artifact.id)" }
  if ([IO.Path]::GetFileName($uri.AbsolutePath) -cne $artifact.filename) { Fail "URL filename mismatch: $($artifact.id)" }
  if ($artifact.filename -match '[\\/]' -or $artifact.filename -in @('.','..')) { Fail "unsafe filename: $($artifact.id)" }
  $byId[$artifact.id] = $artifact
}

if (@($ArtifactId).Count -eq 0) { Fail 'at least one ArtifactId is required' }
if (@($ArtifactId | Select-Object -Unique).Count -ne @($ArtifactId).Count) { Fail 'ArtifactId contains duplicates' }
$selected = foreach ($id in $ArtifactId) {
  if (-not $byId.ContainsKey($id)) { Fail "unknown artifact id: $id" }
  $byId[$id]
}

if ($ValidateOnly) {
  Write-Output "DOCUMENT_SDK_LOCK_VALID artifacts=$($byId.Count) selected=$(@($selected).Count)"
  exit 0
}

if ([string]::IsNullOrWhiteSpace($Destination)) { Fail 'Destination is required outside ValidateOnly' }
if ([string]::IsNullOrWhiteSpace($ApprovalPath) -or -not (Test-Path -LiteralPath $ApprovalPath -PathType Leaf)) { Fail 'approval file is required before acquisition' }
$approval = Get-Content -LiteralPath $ApprovalPath -Raw | ConvertFrom-Json -Depth 10
ExactProperties $approval @('schema','lock_sha256','approved_artifact_ids','approved_by','approved_at','accept_artifact_licenses','acknowledge_provider_accounts_quotas_and_costs') 'approval'
if ($approval.schema -cne 'elite-official-document-sdk-artifact-approval/v1') { Fail 'unsupported approval schema' }
if ($approval.lock_sha256 -cne (Sha256 $LockPath)) { Fail 'approval lock_sha256 mismatch' }
$approvedIds = @($approval.approved_artifact_ids | Sort-Object)
$selectedIds = @($ArtifactId | Sort-Object)
if (($approvedIds -join '|') -cne ($selectedIds -join '|')) { Fail 'approval artifact selection mismatch' }
if ([string]::IsNullOrWhiteSpace($approval.approved_by)) { Fail 'approved_by is required' }
$approvedAt = [datetime]::MinValue
if (-not [datetime]::TryParseExact([string]$approval.approved_at, 'yyyy-MM-dd', [Globalization.CultureInfo]::InvariantCulture, [Globalization.DateTimeStyles]::None, [ref]$approvedAt)) { Fail 'approved_at is not an ISO calendar date' }
if ($approval.accept_artifact_licenses -cne $true) { Fail 'artifact licenses were not accepted' }
if ($approval.acknowledge_provider_accounts_quotas_and_costs -cne $true) { Fail 'provider accounts, quotas and possible usage costs were not acknowledged' }
if (Test-Path -LiteralPath $Destination) { Fail 'Destination must not exist' }
if ($ImportDirectory -and -not (Test-Path -LiteralPath $ImportDirectory -PathType Container)) { Fail 'ImportDirectory does not exist' }

$parent = Split-Path -Parent ([IO.Path]::GetFullPath($Destination))
if (-not $parent) { Fail 'Destination must have a parent directory' }
[IO.Directory]::CreateDirectory($parent) | Out-Null
$stage = Join-Path $parent ('.document-sdk-stage-' + [guid]::NewGuid().ToString('N'))
[IO.Directory]::CreateDirectory($stage) | Out-Null
try {
  $receipts = @()
  foreach ($artifact in $selected) {
    $target = Join-Path $stage $artifact.filename
    if ($ImportDirectory) {
      $source = Join-Path $ImportDirectory $artifact.filename
      if (-not (Test-Path -LiteralPath $source -PathType Leaf)) { Fail "cached artifact missing: $($artifact.filename)" }
      Copy-Item -LiteralPath $source -Destination $target
    } else {
      Invoke-WebRequest -Uri $artifact.url -OutFile $target -UseBasicParsing
    }
    $item = Get-Item -LiteralPath $target
    if ($item.Length -ne [int64]$artifact.bytes) { Fail "byte count mismatch: $($artifact.id)" }
    $actualHash = Sha256 $target
    if ($actualHash -cne $artifact.sha256) { Fail "sha256 mismatch: $($artifact.id)" }
    $receipts += [ordered]@{ id=$artifact.id; filename=$artifact.filename; bytes=$item.Length; sha256=$actualHash; source=if ($ImportDirectory) {'verified-offline-cache'} else {$artifact.url} }
  }
  [ordered]@{
    schema='elite-official-document-sdk-artifact-receipt/v1'
    acquired_at=(Get-Date).ToUniversalTime().ToString('o')
    lock_sha256=(Sha256 $LockPath)
    approved_by=[string]$approval.approved_by
    artifacts=$receipts
  } | ConvertTo-Json -Depth 10 | Set-Content -LiteralPath (Join-Path $stage 'ARTIFACT_RECEIPT.json') -Encoding utf8NoBOM
  Move-Item -LiteralPath $stage -Destination $Destination
  Write-Output "DOCUMENT_SDK_ACQUISITION_PASS artifacts=$(@($selected).Count) destination=$Destination"
} catch {
  if (Test-Path -LiteralPath $stage) { Remove-Item -LiteralPath $stage -Recurse -Force }
  throw
}
````

### FILE: `official_document_sdks/probe_document_sdks.py`
```yaml
block_id: "OFFICIAL-DOCUMENT-SDK-ARTIFACT:probe:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract probe using official public SDK interfaces"
license: "LicenseRef-Workspace-Owner"
sha256: "679a8a24a92723fd273df337c4d61696ee93f88039532b2d4b4727670faea9d2"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import importlib.metadata as metadata
import json
import sys


EXPECTED = {
    "azure-ai-documentintelligence": "1.0.2",
    "azure-ai-contentunderstanding": "1.1.0",
    "google-cloud-documentai": "3.15.0",
    "google-cloud-storage": "3.13.1",
    "google-cloud-documentai-toolbox": "0.17.3",
    "openai": "3.3.1",
    "oci": "2.185.0",
}


def main() -> int:
    observed = {package: metadata.version(package) for package in EXPECTED}
    if observed != EXPECTED:
        raise RuntimeError(f"document SDK version mismatch: {observed!r}")

    from azure.ai.contentunderstanding import ContentUnderstandingClient
    from azure.ai.contentunderstanding.models import AnalysisInput, ContentAnalyzer
    from azure.ai.documentintelligence import DocumentIntelligenceClient
    from google.cloud import documentai_v1
    from google.cloud import storage
    from google.cloud.documentai_toolbox import document as documentai_toolbox_document
    from openai import AsyncOpenAI, OpenAI
    from openai.types.responses import ResponseInputFileParam
    from oci.ai_document import AIServiceDocumentClient
    from oci.ai_document.models import (
        AnalyzeDocumentDetails,
        BoundingPolygon,
        DocumentKeyValueExtractionFeature,
        FieldValue,
        InvoiceProcessorConfig,
        KeyValueDetectionConfidenceEntry,
    )

    symbols = {
        "DocumentIntelligenceClient": DocumentIntelligenceClient.__name__,
        "ContentUnderstandingClient": ContentUnderstandingClient.__name__,
        "AnalysisInput": AnalysisInput.__name__,
        "ContentAnalyzer": ContentAnalyzer.__name__,
        "DocumentProcessorServiceClient": documentai_v1.DocumentProcessorServiceClient.__name__,
        "StorageClient": storage.Client.__name__,
        "DocumentAIToolboxDocument": documentai_toolbox_document.Document.__name__,
        "OpenAI": OpenAI.__name__,
        "AsyncOpenAI": AsyncOpenAI.__name__,
        "ResponseInputFileParam": ResponseInputFileParam.__name__,
        "AIServiceDocumentClient": AIServiceDocumentClient.__name__,
        "AnalyzeDocumentDetails": AnalyzeDocumentDetails.__name__,
        "BoundingPolygon": BoundingPolygon.__name__,
        "DocumentKeyValueExtractionFeature": DocumentKeyValueExtractionFeature.__name__,
        "FieldValue": FieldValue.__name__,
        "InvoiceProcessorConfig": InvoiceProcessorConfig.__name__,
        "KeyValueDetectionConfidenceEntry": KeyValueDetectionConfidenceEntry.__name__,
    }
    print(json.dumps({"status": "DOCUMENT_SDK_PROBE_PASS", "versions": observed, "symbols": symbols}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `official_document_sdks/test_artifact_acquisition.ps1`
```yaml
block_id: "OFFICIAL-DOCUMENT-SDK-ARTIFACT:test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "36a1300f584eba07ecfa21c419906fd0be1a281b968ab8233779918ef1d9293a"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
$ErrorActionPreference = 'Stop'

function Expect-Failure([scriptblock]$Action, [string]$Needle) {
  try { & $Action; throw "expected failure containing: $Needle" }
  catch { if ($_.Exception.Message -notlike "*$Needle*") { throw } }
}

$runner = Join-Path $PSScriptRoot 'acquire_document_sdk_artifacts.ps1'
$lock = Join-Path $PSScriptRoot 'artifact-lock.json'
& $runner -LockPath $lock -ArtifactId 'azure-ai-contentunderstanding-1.1.0-wheel' -ValidateOnly | Out-Null

Expect-Failure { & $runner -LockPath $lock -ArtifactId 'missing-artifact' -ValidateOnly } 'unknown artifact id'
Expect-Failure { & $runner -LockPath $lock -ArtifactId @('openai-3.3.1-wheel','openai-3.3.1-wheel') -ValidateOnly } 'duplicates'
Expect-Failure { & $runner -LockPath $lock -ArtifactId 'openai-3.3.1-wheel' -Destination (Join-Path ([IO.Path]::GetTempPath()) 'elite-document-sdk-never-created') } 'approval file is required'

$tempRoot = Join-Path ([IO.Path]::GetTempPath()) ('elite-document-sdk-test-' + [guid]::NewGuid().ToString('N'))
[IO.Directory]::CreateDirectory($tempRoot) | Out-Null
try {
  $cache = Join-Path $tempRoot 'cache'
  [IO.Directory]::CreateDirectory($cache) | Out-Null
  $fixtureName = 'official_fixture-1.0.0-py3-none-any.whl'
  $fixturePath = Join-Path $cache $fixtureName
  [IO.File]::WriteAllBytes($fixturePath, [Text.Encoding]::UTF8.GetBytes('official-fixture'))
  $fixtureHash = (Get-FileHash -LiteralPath $fixturePath -Algorithm SHA256).Hash.ToLowerInvariant()
  $fixtureBytes = (Get-Item -LiteralPath $fixturePath).Length
  $fixtureLockPath = Join-Path $tempRoot 'lock.json'
  [ordered]@{
    schema='elite-official-document-sdk-artifact-lock/v1'; verified_at='2026-08-26'; artifacts=@([ordered]@{
      id='official-fixture-1.0.0-wheel'; owner='Fixture'; classification='PINNED_CANDIDATE'; package='official-fixture'; version='1.0.0'; python_requires='>=3.10'; filename=$fixtureName; kind='wheel'; url="https://files.pythonhosted.org/packages/test/$fixtureName"; bytes=$fixtureBytes; sha256=$fixtureHash; license_expression='MIT'
    })
  } | ConvertTo-Json -Depth 10 | Set-Content -LiteralPath $fixtureLockPath -Encoding utf8NoBOM
  $approvalPath = Join-Path $tempRoot 'approval.json'
  [ordered]@{
    schema='elite-official-document-sdk-artifact-approval/v1'; lock_sha256=(Get-FileHash -LiteralPath $fixtureLockPath -Algorithm SHA256).Hash.ToLowerInvariant(); approved_artifact_ids=@('official-fixture-1.0.0-wheel'); approved_by='automated-regression'; approved_at='2026-08-26'; accept_artifact_licenses=$true; acknowledge_provider_accounts_quotas_and_costs=$true
  } | ConvertTo-Json -Depth 10 | Set-Content -LiteralPath $approvalPath -Encoding utf8NoBOM
  $destination = Join-Path $tempRoot 'verified'
  & $runner -LockPath $fixtureLockPath -ArtifactId 'official-fixture-1.0.0-wheel' -Destination $destination -ApprovalPath $approvalPath -ImportDirectory $cache | Out-Null
  if (-not (Test-Path -LiteralPath (Join-Path $destination $fixtureName))) { throw 'verified fixture missing' }
  if (-not (Test-Path -LiteralPath (Join-Path $destination 'ARTIFACT_RECEIPT.json'))) { throw 'receipt missing' }
  [IO.File]::AppendAllText($fixturePath, 'tamper')
  $tamperedDestination = Join-Path $tempRoot 'tampered'
  Expect-Failure { & $runner -LockPath $fixtureLockPath -ArtifactId 'official-fixture-1.0.0-wheel' -Destination $tamperedDestination -ApprovalPath $approvalPath -ImportDirectory $cache } 'byte count mismatch'
  if (Test-Path -LiteralPath $tamperedDestination) { throw 'failed acquisition left destination behind' }
} finally {
  if (Test-Path -LiteralPath $tempRoot) { Remove-Item -LiteralPath $tempRoot -Recurse -Force }
}

Write-Output 'DOCUMENT_SDK_ACQUISITION_TEST_PASS negatives=4 offline_verified=1'
````

### FILE: `official_document_sdks/README.md`
```yaml
block_id: "OFFICIAL-DOCUMENT-SDK-ARTIFACT:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "b4528c308d486d5268fc8b30ea0b663abb58bcf535e2a7192d583c89a370125f"
variables: []
secrets_allowed: false
```
````markdown
# Official document SDK artifacts

This directory acquires exact official PyPI artifacts without Git. The acquisition wrapper and probe are Elite-authored; downloaded wheels/source distributions remain official Microsoft, Google, OpenAI or Oracle artifacts under their recorded licenses.

The lock includes stable wheel+source pairs for Microsoft Azure Document Intelligence 1.0.2, Microsoft Azure Content Understanding 1.1.0, Google Cloud Document AI 3.15.0, Google Cloud Storage 3.13.1, OpenAI 3.3.1 and Oracle OCI SDK 2.185.0. The Azure Document Intelligence source distribution contains Microsoft's synchronous/async invoice, layout, classification, batch, result and conversion samples and its official tests. Google Cloud Storage supplies the official client required by Google's batch Document AI sample to list and download per-document JSON results. Oracle's wheel exposes the public AI Document client plus invoice, key-value, confidence, field-value and bounding-polygon model types; it is an integration SDK, not a universal document schema. These are upstream artifacts, not Elite-authored replacements. Google Document AI Toolbox remains the separately conditioned pair described below.

1. Run `acquire_document_sdk_artifacts.ps1 -ValidateOnly` for the intended IDs.
2. Copy `artifact-approval.template.json`, fill the exact lock SHA and selected IDs, and explicitly acknowledge licenses plus provider accounts, quotas and possible usage costs.
3. Acquire to a new destination from PyPI or from a previously downloaded offline cache.
4. Generate a platform-specific transitive dependency lock with hashes, SBOM, license report and vulnerability scan before installation.
5. After installing the selected wheels in an isolated environment, run `probe_document_sdks.py`.

When building a source distribution, bootstrap with a separately verified pip 26.2 or newer version that passes the project's current vulnerability scan. The Python 3.14 venv observed on 2026-08-27 inherited pip 26.0.1 with four advisory groups; pip 26.1.2 still retained CVE-2026-13346, while the exact official pip 26.2 wheel (`931c3036...f2b8aad`) scanned with zero known findings. Do not copy this dated tool version blindly: recheck PyPI identity, hash and advisories at project time.

The lock proves artifact identity, not extraction accuracy. Production remains blocked until the project records provider/resource/model versions, credentials by secret reference, representative ground truth, field-level evaluation, failure handling, retention/privacy and a rollback path. Never place API keys in approval files or receipts.

Oracle OCI SDK 2.185.0 is published as Production/Stable and licensed under UPL-1.0 or Apache-2.0. The exact wheel installed and imported on Python 3.12 with `pip check` passing; a complete 23-package manifest produced zero known findings and zero skips in the dated pip-audit run. This admits the official SDK artifact only. It does not admit the older invoice-only DevRel integration repository, create credentials, call OCI, or claim support for packing lists, bills of lading or customs documents.

Google Cloud Document AI Toolbox 0.17.3 is a conditioned Alpha/experimental candidate, not an admitted production component. Its official dependency range resolves PyArrow 22.0.0, while CVE-2026-25087 is fixed in 23.0.1 outside that range. Apache states the vulnerable C++ pre-buffer API is not exposed by Python bindings, but promotion still requires a later official compatible release or an explicitly owned, expiring reachability exception plus a fresh dependency scan. Do not override package metadata silently.
````

## 6. Configuration surface

| Variable/parámetro | Tipo | Default | Validación | Secreto | Mutabilidad/efecto |
|---|---|---|---|---|---|
| `LockPath` | path | lock junto al runner | archivo + schema exacto | no | cambiar requiere nueva aprobación por SHA |
| `ArtifactId` | lista | ninguno | IDs únicos y presentes | no | selecciona exclusivamente artefactos aprobados |
| `Destination` | path | ninguno | debe no existir | no | salida inmutable de adquisición |
| `ApprovalPath` | path | ninguno | schema/hash/IDs/aceptaciones exactos | no; nunca incluir keys | habilita adquisición |
| `ImportDirectory` | path | red | directorio existente | no | usa cache pero verifica igual |
| `ValidateOnly` | switch | false | sin I/O de artefactos | no | sólo valida lock/selección |

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| PowerShell | `>=7.0`, versión exacta fijada por target | validación/adquisición | MIT | build/ops | Microsoft PowerShell |
| Python | `>=3.10`, versión exacta fijada por target | probe | PSF-2.0 | test/runtime potencial | Python.org |
| Azure Document Intelligence | `1.0.2` wheel+sdist hashes fijados | SDK, 56 samples y 20 tests upstream en source | MIT | runtime/source potencial | PyPI Microsoft |
| Azure Content Understanding | `1.1.0` wheel+sdist hashes fijados | SDK/source | MIT | runtime potencial | PyPI Microsoft |
| Google Cloud Document AI | `3.15.0` wheel+sdist hashes fijados | SDK/source | Apache-2.0 | runtime/source potencial | PyPI Google |
| Google Cloud Storage | `3.13.1` wheel+sdist hashes fijados | cliente requerido por batch Document AI para listar/descargar resultados JSON | Apache-2.0 | runtime/source potencial | PyPI Google |
| Google Cloud Document AI Toolbox | `0.17.3` wheel+sdist hashes fijados | postprocesado condicionado; Alpha y PyArrow CVE pendiente | Apache-2.0 | referencia/candidato | PyPI Google |
| OpenAI | `3.3.1` wheel+sdist hashes fijados | SDK/source para file input opcional | Apache-2.0 | runtime/source potencial | PyPI OpenAI |
| Oracle OCI SDK | `2.185.0` wheel+sdist hashes fijados | AI Document client; invoice/KV/confidence/geometry model types | UPL-1.0 OR Apache-2.0 | runtime/source potencial | PyPI Oracle |

Las dependencias transitivas no están fingidamente bloqueadas: deben resolverse con hashes para el Python/OS del target y pasar SBOM, licencias y vulnerabilidades antes de deployment.

## 8. Apply order

1. Materializar en workspace vacío y ejecutar el test local.
2. Elegir proveedor/versiones en el expediente documental del proyecto.
3. Ejecutar `ValidateOnly` para los IDs exactos.
4. Completar el approval sin secretos y con SHA del lock.
5. Adquirir a un destino nuevo desde red o cache verificada.
6. Generar lock transitivo del target, SBOM/notices y scan; instalar en venv aislado.
7. Ejecutar el probe y luego contract/evaluation gates con cuenta sandbox y corpus aprobado.

En workspace existente, abortar ante rutas ocupadas o conflicto de versión. Rollback elimina únicamente venv/destino recién creados y conserva recibos, fallos y evidencia.

## 9. Verification

```powershell
pwsh -NoProfile -File .\materialize_markdown_pack.ps1 -PackFile .\implementation_packs\OFFICIAL_DOCUMENT_SDK_ARTIFACT_CORE.md -Destination <empty>
pwsh -NoProfile -File <empty>\official_document_sdks\test_artifact_acquisition.ps1
pwsh -NoProfile -File <empty>\official_document_sdks\acquire_document_sdk_artifacts.ps1 -ArtifactId azure-ai-contentunderstanding-1.1.0-sdist -ValidateOnly
```

Éxito esperado: reconstrucción hash-exacta; `DOCUMENT_SDK_ACQUISITION_TEST_PASS negatives=4 offline_verified=1`; lock de catorce artefactos válido. Los negativos cubren ID desconocido/duplicado, ausencia de aprobación y cache alterada sin dejar destino. El probe de imports requiere instalación aislada y se registra por separado; las llamadas reales requieren sandbox/cuenta/corpus y no forman parte del gate offline. Toolbox permanece candidato condicionado y no se selecciona automáticamente.

## 10. Reconstruction evidence

- snapshot vigente: `reconstruction_evidence/DOCUMENT_INTELLIGENCE_MULTI_VENDOR_REAUDIT_2026-08-28_V81.md`;
- entorno histórico limpio: `%LOCALAPPDATA%\Temp\elite-document-sdk-pack-verify-20260826`;
- toolchains: PowerShell 7; Python 3.14 para el probe histórico y Python 3.12.13 para Oracle OCI 2.185.0;
- materialización: seis archivos, manifest↔FILE y SHA-256 verificados por el materializador;
- test local: cinco positivos/negativos efectivos, sin red ni secretos; el lock vigente contiene catorce artefactos;
- fuente oficial: URLs inmutables `files.pythonhosted.org`, tamaños y hashes fijados;
- divergencia: no se afirma Trusted Publishing donde PyPI no lo publicó y no se afirma lock transitivo universal;
- fecha/revisor: 2026-08-28 / Codex; los expedientes históricos V1 permanecen como etapas anteriores.
