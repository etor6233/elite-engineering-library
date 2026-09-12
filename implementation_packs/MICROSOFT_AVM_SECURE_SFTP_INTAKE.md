# Microsoft AVM Secure SFTP Intake

## 1. Metadata

```yaml
pack_id: "MICROSOFT-AVM-SECURE-SFTP-INTAKE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa una composición SFTP privada sobre Microsoft Azure Verified Modules Storage 0.33.0, dos fuentes oficiales AVM y su licencia, con SSH-only, create/write-only, Private Endpoint, soft delete, diagnósticos y verificador Bicep 0.46.1 hash-locked. No crea recursos ni afirma WORM o persistencia automática."
stacks: ["Microsoft Azure", "Bicep 0.46.1", "Azure Verified Modules", "Blob Storage SFTP"]
compatible_with: ["MICROSOFT-DURABLE-DOCUMENT-ORCHESTRATION 0.1.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x", "PYTHON-AZURE-BLOB-IMMUTABLE-EVIDENCE-ADAPTER 0.1.x"]
incompatible_with: ["SFTP público", "autenticación por contraseña", "Azure Storage Shared Key", "permisos read/delete/list del uploader", "WORM en la cuenta SFTP", "persistencia automática desde landing"]
license_expression: "LicenseRef-Workspace-Owner AND MIT"
upstream_sources: ["https://github.com/Azure/bicep-registry-modules/tree/32f1df57c75a48b5f0519083b7048f8fc00c33d8/avm/res/storage/storage-account", "https://github.com/Azure/bicep/releases/tag/v0.46.1", "https://learn.microsoft.com/azure/storage/blobs/secure-file-transfer-protocol-support-connect", "https://learn.microsoft.com/azure/storage/blobs/immutable-storage-overview"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use este pack cuando un proyecto autorizado elija Azure Blob Storage SFTP como puerta de recepción de archivos de proveedores o fábricas. Entrega infraestructura tangible y compilable sin inventar un servidor SFTP: el runtime y sus módulos son Microsoft. Rechácelo si el proyecto no dispone de Azure, red privada, DNS privado, Log Analytics, clave SSH administrada, costo aprobado o un retained-original store separado.

El pack no cubre correo, cámara/móvil ni persistencia empresarial. Los bytes entran a cuarentena; no se consideran completos ni confiables hasta reconciliar el cierre de transferencia, calcular SHA-256, ejecutar el gate de seguridad y copiar el objeto exacto a la retención seleccionada.

## 3. Architecture contract

El módulo oficial fijado crea StorageV2 ZRS con HNS/SFTP, Shared Key deshabilitado, acceso público deshabilitado, firewall deny, Private Endpoint Blob y DNS privado. Un solo usuario local usa SSH key, sin contraseña, con alcance al contenedor y permisos exclusivos create/write. Soft delete y diagnósticos protegen/observan el landing.

Microsoft documenta dos límites decisivos: blob versioning no funciona con HNS y immutable storage no funciona mientras SFTP está habilitado. Por eso este landing nunca sustituye la conservación WORM: el pipeline durable copia bytes y hash a otra cuenta retenida antes de extracción o storage empresarial. La cuenta, claves, red, costos y despliegue son decisiones explícitas del proyecto.

## 4. Exact file manifest

```text
CREATE azure_avm_secure_sftp_intake/README.md
CREATE azure_avm_secure_sftp_intake/main.bicep
CREATE azure_avm_secure_sftp_intake/main.bicepparam.example
CREATE azure_avm_secure_sftp_intake/UPSTREAM.lock.json
CREATE azure_avm_secure_sftp_intake/verify_contract.ps1
CREATE azure_avm_secure_sftp_intake/upstream/LICENSE
CREATE azure_avm_secure_sftp_intake/upstream/local-user/main.bicep
CREATE azure_avm_secure_sftp_intake/upstream/waf-aligned/main.test.bicep
```

## 5. Materialization blocks

### FILE: `azure_avm_secure_sftp_intake/README.md`

```yaml
block_id: "MICROSOFT-AVM-SFTP:01:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://azure-avm-secure-sftp-intake/README.md"
license: "LicenseRef-Workspace-Owner"
sha256: "ee594fc1bdf4f0436ab4b6be75fb24f2988eaafbfd4e7180303e30bd75c97365"
variables: []
secrets_allowed: false
```
````text
# Microsoft AVM secure SFTP landing

This materializes a fail-closed SFTP landing account from Microsoft Azure Verified Modules, not a custom SFTP server. The composition pins `avm/res/storage/storage-account:0.33.0`, enables hierarchical namespace and SFTP, accepts only an SSH public key, grants the uploader only create/write (`cw`) in one quarantine container, disables Azure Storage Shared Key and public blob access, denies the public network, and creates a Blob private endpoint with private DNS.

The landing account intentionally uses soft delete rather than blob versioning or WORM. Microsoft documents that versioning is unavailable with hierarchical namespace and immutable storage is unavailable while SFTP is enabled. A project must therefore copy the exact uploaded bytes, after close/reconciliation and SHA-256 calculation, into the separately selected retained-original store before extraction or persistence. Compose the existing Azure retained-evidence lane and durable document pipeline for that boundary.

Run `pwsh ./verify_contract.ps1 -AllowNetwork` to download the exact Microsoft Bicep CLI `0.46.1`, verify its published SHA-256, restore the pinned AVM module, and compile this source. No Azure resource is created by the verifier. A live deployment remains blocked until the project supplies an Azure subscription, region, private network/DNS, Log Analytics workspace, SSH key ownership/rotation, cost approval, downstream identity, Event Grid/durable reconciliation, malware gate, retained store, backup/restore and a real SFTP transfer test.

````

### FILE: `azure_avm_secure_sftp_intake/main.bicep`

```yaml
block_id: "MICROSOFT-AVM-SFTP:02:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://azure-avm-secure-sftp-intake/main.bicep"
license: "LicenseRef-Workspace-Owner"
sha256: "f26bf00cc8568921f00faf2ad61553e588c803937d3c41a95840115903c1983d"
variables: []
secrets_allowed: false
```
````text
targetScope = 'resourceGroup'

@description('Globally unique Azure Storage account name.')
@minLength(3)
@maxLength(24)
param storageAccountName string

@description('Azure region for the SFTP landing account.')
param location string = resourceGroup().location

@description('Private quarantine container used as the SFTP home directory.')
@minLength(3)
@maxLength(63)
param quarantineContainerName string = 'sftp-quarantine'

@description('Local SFTP user name scoped to the quarantine container.')
@minLength(3)
param sftpUserName string

@secure()
@description('OpenSSH public key. Password authentication is prohibited.')
param sshPublicKey string

@description('Subnet resource ID for the Blob private endpoint.')
param privateEndpointSubnetResourceId string

@description('Private DNS zone resource ID for privatelink.blob.core.windows.net.')
param privateDnsZoneResourceId string

@description('Log Analytics workspace resource ID for storage data-plane diagnostics.')
param logAnalyticsWorkspaceResourceId string

@description('Soft-delete retention for landing bytes. The separate retained evidence store remains mandatory.')
@minValue(7)
@maxValue(365)
param softDeleteRetentionDays int = 30

module storage 'br/public:avm/res/storage/storage-account:0.33.0' = {
  name: 'secure-sftp-${uniqueString(storageAccountName, resourceGroup().id)}'
  params: {
    name: storageAccountName
    location: location
    skuName: 'Standard_ZRS'
    kind: 'StorageV2'
    enableHierarchicalNamespace: true
    enableSftp: true
    isLocalUserEnabled: true
    allowSharedKeyAccess: false
    defaultToOAuthAuthentication: true
    allowBlobPublicAccess: false
    publicNetworkAccess: 'Disabled'
    supportsHttpsTrafficOnly: true
    minimumTlsVersion: 'TLS1_2'
    requireInfrastructureEncryption: true
    networkAcls: {
      bypass: 'None'
      defaultAction: 'Deny'
      ipRules: []
      virtualNetworkRules: []
    }
    privateEndpoints: [
      {
        service: 'blob'
        subnetResourceId: privateEndpointSubnetResourceId
        privateDnsZoneGroup: {
          privateDnsZoneGroupConfigs: [
            {
              privateDnsZoneResourceId: privateDnsZoneResourceId
            }
          ]
        }
      }
    ]
    blobServices: {
      containerDeleteRetentionPolicyEnabled: true
      containerDeleteRetentionPolicyDays: softDeleteRetentionDays
      deleteRetentionPolicyEnabled: true
      deleteRetentionPolicyDays: softDeleteRetentionDays
      isVersioningEnabled: false
      containers: [
        {
          name: quarantineContainerName
          publicAccess: 'None'
        }
      ]
      diagnosticSettings: [
        {
          name: 'blob-audit'
          workspaceResourceId: logAnalyticsWorkspaceResourceId
          logCategoriesAndGroups: [
            { category: 'StorageRead' }
            { category: 'StorageWrite' }
            { category: 'StorageDelete' }
          ]
          metricCategories: [
            { category: 'Transaction' }
          ]
        }
      ]
    }
    localUsers: [
      {
        name: sftpUserName
        hasSharedKey: false
        hasSshKey: true
        hasSshPassword: false
        homeDirectory: quarantineContainerName
        sshAuthorizedKeys: [
          {
            description: 'project-managed upload key'
            key: sshPublicKey
          }
        ]
        permissionScopes: [
          {
            permissions: 'cw'
            service: 'blob'
            resourceName: quarantineContainerName
          }
        ]
      }
    ]
  }
}

output storageAccountResourceId string = storage.outputs.resourceId
output privateBlobEndpoint string = storage.outputs.primaryBlobEndpoint
````

### FILE: `azure_avm_secure_sftp_intake/main.bicepparam.example`

```yaml
block_id: "MICROSOFT-AVM-SFTP:03:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://azure-avm-secure-sftp-intake/main.bicepparam.example"
license: "LicenseRef-Workspace-Owner"
sha256: "af43e7dc25b2e973c0f231296b1aed1856214408de336c504ee5bfb98850a014"
variables: []
secrets_allowed: false
```
````text
using './main.bicep'

param storageAccountName = 'replacewithgloballyunique'
param location = 'eastus2'
param quarantineContainerName = 'sftp-quarantine'
param sftpUserName = 'supplier-upload'
param sshPublicKey = 'ssh-ed25519 REPLACE_WITH_PROJECT_OWNED_PUBLIC_KEY'
param privateEndpointSubnetResourceId = '/subscriptions/REPLACE/resourceGroups/REPLACE/providers/Microsoft.Network/virtualNetworks/REPLACE/subnets/REPLACE'
param privateDnsZoneResourceId = '/subscriptions/REPLACE/resourceGroups/REPLACE/providers/Microsoft.Network/privateDnsZones/privatelink.blob.core.windows.net'
param logAnalyticsWorkspaceResourceId = '/subscriptions/REPLACE/resourceGroups/REPLACE/providers/Microsoft.OperationalInsights/workspaces/REPLACE'
param softDeleteRetentionDays = 30

````

### FILE: `azure_avm_secure_sftp_intake/UPSTREAM.lock.json`

```yaml
block_id: "MICROSOFT-AVM-SFTP:04:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://azure-avm-secure-sftp-intake/UPSTREAM.lock.json"
license: "LicenseRef-Workspace-Owner"
sha256: "4f94fc7b53abdd4ccd9bdb9a482592b280079aa4d0313acf8e4bfd77208f69dc"
variables: []
secrets_allowed: false
```
````text
{
  "schema_version": "1.0",
  "verified_at": "2026-08-28",
  "azure_verified_modules": {
    "repository": "Azure/bicep-registry-modules",
    "commit": "32f1df57c75a48b5f0519083b7048f8fc00c33d8",
    "commit_signature": "verified",
    "storage_account_module": "br/public:avm/res/storage/storage-account:0.33.0",
    "local_user_module": "br/public:avm/res/storage/storage-account/local-user:0.1.0",
    "license": "MIT"
  },
  "bicep_cli": {
    "version": "0.46.1",
    "asset": "bicep-win-x64.exe",
    "url": "https://github.com/Azure/bicep/releases/download/v0.46.1/bicep-win-x64.exe",
    "size": 117685696,
    "sha256": "441d3d6094513acaa8a9be0b96b754d174bc0c19064c7be04e35b87e67b66250"
  },
  "official_constraints": {
    "versioning_with_hierarchical_namespace": "unsupported",
    "immutable_storage_while_sftp_enabled": "unsupported",
    "private_endpoint_sftp": "supported",
    "sftp_port": 22,
    "sftp_uploaded_content_md5_default": null
  }
}

````

### FILE: `azure_avm_secure_sftp_intake/verify_contract.ps1`

```yaml
block_id: "MICROSOFT-AVM-SFTP:05:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://azure-avm-secure-sftp-intake/verify_contract.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "da64ec5b9443b6de38ba2c663fc085226fbea86d98cf73db890fcc3d877659c5"
variables: []
secrets_allowed: false
```
````text
#requires -Version 7.0
[CmdletBinding()]
param(
  [string]$SourceRoot = $PSScriptRoot,
  [string]$BicepExe = '',
  [switch]$AllowNetwork
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

function Assert-Present([string]$Text, [string]$Literal) {
  $count = ([regex]::Matches($Text, [regex]::Escape($Literal))).Count
  if ($count -lt 1) { throw "Required contract token missing: [$Literal]" }
}

$mainPath = Join-Path $SourceRoot 'main.bicep'
$lockPath = Join-Path $SourceRoot 'UPSTREAM.lock.json'
if (-not (Test-Path -LiteralPath $mainPath -PathType Leaf)) { throw 'main.bicep missing' }
if (-not (Test-Path -LiteralPath $lockPath -PathType Leaf)) { throw 'UPSTREAM.lock.json missing' }

$main = Get-Content -Raw -LiteralPath $mainPath
$lock = Get-Content -Raw -LiteralPath $lockPath | ConvertFrom-Json -Depth 20
if ($lock.azure_verified_modules.storage_account_module -cne 'br/public:avm/res/storage/storage-account:0.33.0') { throw 'AVM storage module lock drifted' }
if ($lock.bicep_cli.version -cne '0.46.1') { throw 'Bicep version lock drifted' }
if ($lock.bicep_cli.sha256 -cne '441d3d6094513acaa8a9be0b96b754d174bc0c19064c7be04e35b87e67b66250') { throw 'Bicep hash lock drifted' }

$required = @(
  "br/public:avm/res/storage/storage-account:0.33.0",
  'enableHierarchicalNamespace: true',
  'enableSftp: true',
  'isLocalUserEnabled: true',
  'allowSharedKeyAccess: false',
  "publicNetworkAccess: 'Disabled'",
  'allowBlobPublicAccess: false',
  'requireInfrastructureEncryption: true',
  "minimumTlsVersion: 'TLS1_2'",
  "bypass: 'None'",
  "defaultAction: 'Deny'",
  "service: 'blob'",
  'privateDnsZoneResourceId: privateDnsZoneResourceId',
  'containerDeleteRetentionPolicyEnabled: true',
  'deleteRetentionPolicyEnabled: true',
  'isVersioningEnabled: false',
  "publicAccess: 'None'",
  "{ category: 'StorageRead' }",
  "{ category: 'StorageWrite' }",
  "{ category: 'StorageDelete' }",
  'hasSharedKey: false',
  'hasSshKey: true',
  'hasSshPassword: false',
  "permissions: 'cw'",
  'key: sshPublicKey'
)
foreach ($literal in $required) { Assert-Present $main $literal }
foreach ($literal in @("br/public:avm/res/storage/storage-account:0.33.0", 'hasSshPassword: false', "permissions: 'cw'")) {
  $count = ([regex]::Matches($main, [regex]::Escape($literal))).Count
  if ($count -ne 1) { throw "Expected exactly one occurrence of [$literal], found $count" }
}

$forbidden = @(
  'hasSshPassword: true',
  'hasSharedKey: true',
  'allowSharedKeyAccess: true',
  "publicNetworkAccess: 'Enabled'",
  "permissions: 'r'",
  "permissions: 'rcwdl'",
  'isVersioningEnabled: true'
)
foreach ($literal in $forbidden) {
  if ($main.Contains($literal, [System.StringComparison]::Ordinal)) { throw "Forbidden contract token present: $literal" }
}

$compile = 'SKIPPED_NO_PINNED_BICEP'
if ([string]::IsNullOrWhiteSpace($BicepExe)) {
  $candidate = Get-Command bicep -ErrorAction SilentlyContinue
  if ($null -ne $candidate) { $BicepExe = $candidate.Source }
}
if ([string]::IsNullOrWhiteSpace($BicepExe) -and $AllowNetwork) {
  $toolRoot = Join-Path ([IO.Path]::GetTempPath()) 'elite-bicep-cli-0.46.1-441d3d60'
  if (-not (Test-Path -LiteralPath $toolRoot)) { New-Item -ItemType Directory -Path $toolRoot | Out-Null }
  $BicepExe = Join-Path $toolRoot 'bicep-win-x64.exe'
  if (-not (Test-Path -LiteralPath $BicepExe -PathType Leaf)) {
    Invoke-WebRequest -Uri $lock.bicep_cli.url -OutFile $BicepExe
  }
}

if (-not [string]::IsNullOrWhiteSpace($BicepExe)) {
  if (-not (Test-Path -LiteralPath $BicepExe -PathType Leaf)) { throw 'Specified Bicep executable missing' }
  $actualSize = (Get-Item -LiteralPath $BicepExe).Length
  $actualHash = (Get-FileHash -LiteralPath $BicepExe -Algorithm SHA256).Hash.ToLowerInvariant()
  if ($actualSize -ne [long]$lock.bicep_cli.size) { throw "Bicep size mismatch: $actualSize" }
  if ($actualHash -cne $lock.bicep_cli.sha256) { throw "Bicep SHA-256 mismatch: $actualHash" }
  $versionText = (& $BicepExe --version 2>&1) -join "`n"
  if ($LASTEXITCODE -ne 0 -or $versionText -notmatch '0\.46\.1') { throw "Unexpected Bicep version: $versionText" }
  $outputRoot = Join-Path ([IO.Path]::GetTempPath()) ('elite-sftp-arm-' + [guid]::NewGuid().ToString('N'))
  New-Item -ItemType Directory -Path $outputRoot | Out-Null
  $armPath = Join-Path $outputRoot 'main.json'
  & $BicepExe build $mainPath --outfile $armPath
  if ($LASTEXITCODE -ne 0 -or -not (Test-Path -LiteralPath $armPath -PathType Leaf)) { throw 'Bicep compilation failed' }
  $armText = Get-Content -Raw -LiteralPath $armPath
  foreach ($token in @('isSftpEnabled', 'hasSshPassword', 'allowSharedKeyAccess', 'publicNetworkAccess', 'StorageWrite', 'privateDnsZoneResourceId')) {
    if (-not $armText.Contains($token, [System.StringComparison]::Ordinal)) { throw "Compiled ARM missing token: $token" }
  }
  $compile = 'PASS_PINNED_0.46.1'
}

Write-Output "MICROSOFT_AVM_SECURE_SFTP_INTAKE_PASS static=1 compile=$compile password=0 shared_key=0 public_network=0 upload_permissions=cw retained_store=SEPARATE_REQUIRED"
````

### FILE: `azure_avm_secure_sftp_intake/upstream/LICENSE`

```yaml
block_id: "MICROSOFT-AVM-SFTP:06:v1"
operation: CREATE
provenance: VERBATIM
source: "https://github.com/Azure/bicep-registry-modules/blob/32f1df57c75a48b5f0519083b7048f8fc00c33d8/LICENSE"
license: "MIT"
sha256: "c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383"
variables: []
secrets_allowed: false
```
````text
    MIT License

    Copyright (c) Microsoft Corporation.

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE
````

### FILE: `azure_avm_secure_sftp_intake/upstream/local-user/main.bicep`

```yaml
block_id: "MICROSOFT-AVM-SFTP:07:v1"
operation: CREATE
provenance: VERBATIM
source: "https://github.com/Azure/bicep-registry-modules/blob/32f1df57c75a48b5f0519083b7048f8fc00c33d8/avm/res/storage/storage-account/local-user/main.bicep"
license: "MIT"
sha256: "efd933ceeb703f6d71b6a8f185ff41f38e0734014f8f35930c1b2afb9c059adc"
variables: []
secrets_allowed: false
```
````text
metadata name = 'Storage Account Local Users'
metadata description = 'This module deploys a Storage Account Local User, which is used for SFTP authentication.'

@maxLength(24)
@description('Conditional. The name of the parent Storage Account. Required if the template is used in a standalone deployment.')
param storageAccountName string

@description('Required. The name of the local user used for SFTP Authentication.')
param name string

@description('Optional. Indicates whether shared key exists. Set it to false to remove existing shared key.')
param hasSharedKey bool = false

@description('Required. Indicates whether SSH key exists. Set it to false to remove existing SSH key.')
param hasSshKey bool

@description('Required. Indicates whether SSH password exists. Set it to false to remove existing SSH password.')
param hasSshPassword bool

@description('Optional. The local user home directory.')
param homeDirectory string = ''

@description('Required. The permission scopes of the local user.')
param permissionScopes permissionScopeType[]

@description('Optional. The local user SSH authorized keys for SFTP.')
param sshAuthorizedKeys sshAuthorizedKeyType[]?

@description('Optional. Enable/Disable usage telemetry for module.')
param enableTelemetry bool = true

#disable-next-line no-deployments-resources
resource avmTelemetry 'Microsoft.Resources/deployments@2024-03-01' = if (enableTelemetry) {
  name: '46d3xbcp.res.storage-localuser.${replace('-..--..-', '.', '-')}.${substring(uniqueString(deployment().name), 0, 4)}'
  properties: {
    mode: 'Incremental'
    template: {
      '$schema': 'https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#'
      contentVersion: '1.0.0.0'
      resources: []
      outputs: {
        telemetry: {
          type: 'String'
          value: 'For more information, see https://aka.ms/avm/TelemetryInfo'
        }
      }
    }
  }
}

resource storageAccount 'Microsoft.Storage/storageAccounts@2025-06-01' existing = {
  name: storageAccountName
}

resource localUsers 'Microsoft.Storage/storageAccounts/localUsers@2025-06-01' = {
  name: name
  parent: storageAccount
  properties: {
    hasSharedKey: hasSharedKey
    hasSshKey: hasSshKey
    hasSshPassword: hasSshPassword
    homeDirectory: homeDirectory
    permissionScopes: permissionScopes
    sshAuthorizedKeys: sshAuthorizedKeys
  }
}

@description('The name of the deployed local user.')
output name string = localUsers.name

@description('The resource group of the deployed local user.')
output resourceGroupName string = resourceGroup().name

@description('The resource ID of the deployed local user.')
output resourceId string = localUsers.id

// =============== //
//   Definitions   //
// =============== //
@export()
type sshAuthorizedKeyType = {
  @description('Optional. Description used to store the function/usage of the key.')
  description: string?

  @secure()
  @description('Required. SSH public key base64 encoded. The format should be: \'{keyType} {keyData}\', e.g. ssh-rsa AAAABBBB.')
  key: string
}

@export()
type permissionScopeType = {
  @description('Required. The permissions for the local user. Possible values include: Read (r), Write (w), Delete (d), List (l), and Create (c).')
  permissions: string

  @description('Required. The name of resource, normally the container name or the file share name, used by the local user.')
  resourceName: string

  @description('Required. The service used by the local user, e.g. blob, file.')
  service: string
}
````

### FILE: `azure_avm_secure_sftp_intake/upstream/waf-aligned/main.test.bicep`

```yaml
block_id: "MICROSOFT-AVM-SFTP:08:v1"
operation: CREATE
provenance: VERBATIM
source: "https://github.com/Azure/bicep-registry-modules/blob/32f1df57c75a48b5f0519083b7048f8fc00c33d8/avm/res/storage/storage-account/tests/e2e/waf-aligned/main.test.bicep"
license: "MIT"
sha256: "8f4e0c67787dfbb829824121cb813415cafd85d9ed2f15740c2b5ad9de4eab24"
variables: []
secrets_allowed: false
```
````text
targetScope = 'subscription'

metadata name = 'WAF-aligned'
metadata description = 'This instance deploys the module in alignment with the best-practices of the Azure Well-Architected Framework.'

// ========== //
// Parameters //
// ========== //

@description('Optional. The name of the resource group to deploy for testing purposes.')
@maxLength(90)
param resourceGroupName string = 'dep-${namePrefix}-storage.storageaccounts-${serviceShort}-rg'

@description('Optional. The location to deploy resources to.')
param resourceLocation string = deployment().location

@description('Optional. A short identifier for the kind of deployment. Should be kept short to not run into resource-name length-constraints.')
param serviceShort string = 'ssawaf'

@description('Optional. A token to inject into the name of each resource.')
param namePrefix string = '#_namePrefix_#'

// ============ //
// Dependencies //
// ============ //

// General resources
// =================
resource resourceGroup 'Microsoft.Resources/resourceGroups@2025-04-01' = {
  name: resourceGroupName
  location: resourceLocation
}

module nestedDependencies 'dependencies.bicep' = {
  scope: resourceGroup
  name: '${uniqueString(deployment().name, resourceLocation)}-nestedDependencies'
  params: {
    location: resourceLocation
    virtualNetworkName: 'dep-${namePrefix}-vnet-${serviceShort}'
    managedIdentityName: 'dep-${namePrefix}-msi-${serviceShort}'
  }
}

// Diagnostics
// ===========
module diagnosticDependencies '../../../../../../../utilities/e2e-template-assets/templates/diagnostic.dependencies.bicep' = {
  scope: resourceGroup
  name: '${uniqueString(deployment().name, resourceLocation)}-diagnosticDependencies'
  params: {
    storageAccountName: 'dep${namePrefix}diasa${serviceShort}01'
    logAnalyticsWorkspaceName: 'dep-${namePrefix}-law-${serviceShort}'
    eventHubNamespaceEventHubName: 'dep-${namePrefix}-evh-${serviceShort}'
    eventHubNamespaceName: 'dep-${namePrefix}-evhns-${serviceShort}'
    location: resourceLocation
  }
}

// ============== //
// Test Execution //
// ============== //

@batchSize(1)
module testDeployment '../../../main.bicep' = [
  for iteration in ['init', 'idem']: {
    scope: resourceGroup
    name: '${uniqueString(deployment().name, resourceLocation)}-test-${serviceShort}-${iteration}'
    params: {
      name: '${namePrefix}${serviceShort}001'
      skuName: 'Standard_ZRS'
      allowBlobPublicAccess: false
      requireInfrastructureEncryption: true
      largeFileSharesState: 'Enabled'
      privateEndpoints: [
        {
          service: 'blob'
          subnetResourceId: nestedDependencies.outputs.subnetResourceId
          privateDnsZoneGroup: {
            privateDnsZoneGroupConfigs: [
              {
                privateDnsZoneResourceId: nestedDependencies.outputs.privateDNSZoneResourceId
              }
            ]
          }
          tags: {
            'hidden-title': 'This is visible in the resource name'
            Environment: 'Non-Prod'
            Role: 'DeploymentValidation'
          }
        }
      ]
      networkAcls: {
        bypass: 'AzureServices'
        defaultAction: 'Deny'
        virtualNetworkRules: [
          {
            action: 'Allow'
            id: nestedDependencies.outputs.subnetResourceId
          }
        ]
        ipRules: [
          {
            action: 'Allow'
            value: '1.1.1.1'
          }
        ]
      }
      localUsers: [
        {
          name: 'testuser'
          hasSharedKey: false
          hasSshKey: true
          hasSshPassword: false
          homeDirectory: 'avdscripts'
          permissionScopes: [
            {
              permissions: 'r'
              service: 'blob'
              resourceName: 'avdscripts'
            }
          ]
        }
      ]
      blobServices: {
        lastAccessTimeTrackingPolicyEnabled: true
        diagnosticSettings: [
          {
            name: 'customSetting'
            metricCategories: [
              {
                category: 'AllMetrics'
              }
            ]
            eventHubName: diagnosticDependencies.outputs.eventHubNamespaceEventHubName
            eventHubAuthorizationRuleResourceId: diagnosticDependencies.outputs.eventHubAuthorizationRuleId
            storageAccountResourceId: diagnosticDependencies.outputs.storageAccountResourceId
            workspaceResourceId: diagnosticDependencies.outputs.logAnalyticsWorkspaceResourceId
          }
        ]
        containers: [
          {
            name: 'avdscripts'
            publicAccess: 'None'
          }
          {
            name: 'archivecontainer'
            publicAccess: 'None'
            metadata: {
              testKey: 'testValue'
            }
            immutabilityPolicy: {
              immutabilityPeriodSinceCreationInDays: 7
              allowProtectedAppendWrites: false
              allowProtectedAppendWritesAll: true
            }
          }
        ]
        automaticSnapshotPolicyEnabled: true
        containerDeleteRetentionPolicyEnabled: true
        containerDeleteRetentionPolicyDays: 10
        deleteRetentionPolicyEnabled: true
        deleteRetentionPolicyDays: 9
        isVersioningEnabled: true
        versionDeletePolicyDays: 3
      }
      fileServices: {
        diagnosticSettings: [
          {
            name: 'customSetting'
            metricCategories: [
              {
                category: 'AllMetrics'
              }
            ]
            eventHubName: diagnosticDependencies.outputs.eventHubNamespaceEventHubName
            eventHubAuthorizationRuleResourceId: diagnosticDependencies.outputs.eventHubAuthorizationRuleId
            storageAccountResourceId: diagnosticDependencies.outputs.storageAccountResourceId
            workspaceResourceId: diagnosticDependencies.outputs.logAnalyticsWorkspaceResourceId
          }
        ]
        shares: [
          {
            name: 'avdprofiles'
            accessTier: 'Hot'
            shareQuota: 5120
          }
          {
            name: 'avdprofiles2'
            shareQuota: 102400
          }
        ]
      }
      tableServices: {
        diagnosticSettings: [
          {
            name: 'customSetting'
            metricCategories: [
              {
                category: 'AllMetrics'
              }
            ]
            eventHubName: diagnosticDependencies.outputs.eventHubNamespaceEventHubName
            eventHubAuthorizationRuleResourceId: diagnosticDependencies.outputs.eventHubAuthorizationRuleId
            storageAccountResourceId: diagnosticDependencies.outputs.storageAccountResourceId
            workspaceResourceId: diagnosticDependencies.outputs.logAnalyticsWorkspaceResourceId
          }
        ]
        tables: [
          {
            name: 'table1'
          }
          {
            name: 'table2'
          }
        ]
      }
      queueServices: {
        diagnosticSettings: [
          {
            name: 'customSetting'
            metricCategories: [
              {
                category: 'AllMetrics'
              }
            ]
            eventHubName: diagnosticDependencies.outputs.eventHubNamespaceEventHubName
            eventHubAuthorizationRuleResourceId: diagnosticDependencies.outputs.eventHubAuthorizationRuleId
            storageAccountResourceId: diagnosticDependencies.outputs.storageAccountResourceId
            workspaceResourceId: diagnosticDependencies.outputs.logAnalyticsWorkspaceResourceId
          }
        ]
        queues: [
          {
            name: 'queue1'
            metadata: {
              key1: 'value1'
              key2: 'value2'
            }
          }
          {
            name: 'queue2'
            metadata: {}
          }
        ]
      }
      sasExpirationPeriod: '180.00:00:00'
      managedIdentities: {
        systemAssigned: true
        userAssignedResourceIds: [
          nestedDependencies.outputs.managedIdentityResourceId
        ]
      }
      diagnosticSettings: [
        {
          name: 'customSetting'
          metricCategories: [
            {
              category: 'AllMetrics'
            }
          ]
          eventHubName: diagnosticDependencies.outputs.eventHubNamespaceEventHubName
          eventHubAuthorizationRuleResourceId: diagnosticDependencies.outputs.eventHubAuthorizationRuleId
          storageAccountResourceId: diagnosticDependencies.outputs.storageAccountResourceId
          workspaceResourceId: diagnosticDependencies.outputs.logAnalyticsWorkspaceResourceId
        }
      ]
      managementPolicyRules: [
        {
          enabled: true
          name: 'FirstRule'
          type: 'Lifecycle'
          definition: {
            actions: {
              baseBlob: {
                delete: {
                  daysAfterModificationGreaterThan: 30
                }
                tierToCool: {
                  daysAfterLastAccessTimeGreaterThan: 5
                }
              }
            }
            filters: {
              blobIndexMatch: [
                {
                  name: 'BlobIndex'
                  op: '=='
                  value: '1'
                }
              ]
              blobTypes: [
                'blockBlob'
              ]
              prefixMatch: [
                'sample-container/log'
              ]
            }
          }
        }
      ]
      tags: {
        'hidden-title': 'This is visible in the resource name'
        Environment: 'Non-Prod'
        Role: 'DeploymentValidation'
      }
    }
  }
]
````

## 6. Configuration surface

Entradas obligatorias: nombre global de Storage, región, contenedor, usuario, clave pública OpenSSH, subnet de Private Endpoint, zona DNS privada y workspace Log Analytics. El ejemplo sólo contiene placeholders. Ninguna credencial privada se guarda en Markdown. El proyecto registra owner, rotación/revocación de clave, IP/VPN/ExpressRoute o Firewall, presupuesto, retención del landing y retained store separado.

## 7. Dependency bill

- Microsoft AVM Storage Account `0.33.0`, publicado en el Bicep Public Registry.
- Microsoft AVM Local User `0.1.0`, fuente exacta incluida como evidencia.
- Microsoft Bicep CLI `0.46.1`, asset Windows x64 de 117.685.696 bytes y SHA-256 `441d3d6094513acaa8a9be0b96b754d174bc0c19064c7be04e35b87e67b66250`.
- Azure Storage SFTP, Private Link/Private DNS y Log Analytics; pueden generar costo y no se activan por materializar.

## 8. Apply order

1. Completar readiness de cuenta, región, costo, red, DNS, Log Analytics, owner y rotación de clave.
2. Materializar en destino vacío y ejecutar `verify_contract.ps1 -AllowNetwork`; esto sólo compila.
3. Componer retained-original Azure y pipeline durable en destinos no colisionantes.
4. Revisar el ARM compilado y obtener aprobación para crear recursos externos.
5. Desplegar primero en DEV, realizar transferencia SFTP real, reconciliar tamaño/hash, malware, copia retenida y replay.
6. Recién después habilitar el routing documental; nunca almacenar campos automáticamente por éxito de transferencia.

## 9. Verification

Gate ejecutado el 2026-08-28: Bicep oficial 0.46.1, tamaño/hash exactos, restore de AVM 0.33.0 y compilación PASS. El verificador exige los controles estáticos y el ARM generado; las mutaciones password=true y permissions=rcwdl fueron rechazadas. La fuente oficial AVM conserva su test WAF E2E, pero la biblioteca no afirma un deployment Azure propio sin cuenta/costo autorizados.

## 10. Reconstruction evidence

Gobierna `reconstruction_evidence/MICROSOFT_AVM_SECURE_SFTP_INTAKE_2026-08-28_V87.md`. Procedencia: tres bloques VERBATIM Microsoft (LICENSE, módulo Local User y test WAF) y cinco bloques AUTHORED (composición, lock, ejemplo, documentación y verificador). El código authored no se atribuye a Microsoft; su admisión se basa en compile exacto y negativos reproducibles.
