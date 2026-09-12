# Microsoft Business Central Executable Platform

## 1. Metadata

```yaml
pack_id: "MICROSOFT-BUSINESS-CENTRAL-EXECUTABLE-PLATFORM"
pack_version: "0.1.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un lock y una ejecución fail-closed sin Git para adquirir BcContainerHelper 6.1.16 desde PowerShell Gallery y crear, sólo con aprobación del usuario, un entorno local de desarrollo Business Central 28.4.53241.0 mediante comandos oficiales Microsoft."
stacks: ["PowerShell 7", "Windows containers", "Microsoft Business Central 28.4", "BcContainerHelper 6.1.16", "AL"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "DEPENDENCY-LICENSE-EVIDENCE-CORE 0.1.x"]
incompatible_with: ["Linux container host", "production without Business Central license", "implicit EULA acceptance", "moving Business Central artifact", "Git-required bootstrap"]
license_expression: "LicenseRef-Workspace-Owner AND MIT AND LicenseRef-Microsoft-Product-EULA"
upstream_sources: ["https://github.com/microsoft/BCApps/tree/2eae56d704a1fd035d104f333602aea7091b7749", "https://github.com/microsoft/BCApps/releases/tag/releases%2F28.4%2FStrictMode", "https://www.powershellgallery.com/packages/BcContainerHelper/6.1.16", "https://learn.microsoft.com/dynamics365/business-central/dev-itpro/developer/devenv-running-container-development"]
verified_at: "2026-08-31"
```

## 2. Applicability

Use only when the project explicitly chooses Microsoft Business Central as its enterprise platform and accepts AL, Windows containers, Microsoft product terms and the difference between local demo development and production licensing. The pack is not the Go/PostgreSQL foundation and does not convert BCApps into a portable runtime.

The product code, artifact and container commands are Microsoft. Elite contributes only acquisition, validation, approval and receipt orchestration marked `AUTHORED`. No Microsoft code is relabeled as local and no local business rule is attributed to Microsoft.

## 3. Architecture contract

The lock pins the official Gallery package by byte length, SHA-256 and Gallery SHA-512; critical MIT license, manifest and nuspec files also have exact hashes. Acquisition rejects traversal, occupied destinations and any identity mismatch before module import. The project configuration contains only secret environment-variable names.

The runtime lane is development-only. `invoke_business_central_dev.ps1` requires an execution approval tied to the exact config and lock hashes, explicit network/container/EULA flags, administrator authority, Windows container mode and populated credential environment variables. It calls the official `New-BcContainer` command against the exact observed 28.4.53241.0 W1 artifact URL and records installed application count. It never creates Git metadata.

The signed `releases/28.4/StrictMode` source omits Base Application, APIV2 and full tests. The lock therefore retains both that signed release identity and the complete signed BCApps development snapshot `2eae56d704a1fd035d104f333602aea7091b7749` without claiming they are the same artifact. The snapshot archive, tree, MIT license, 36,673 AL files and the Inventory/Availability/Tracking/Warehouse/SCM-Reservation path manifests are fixed exactly. Production remains blocked by product license, localization, project data and target evidence.

## 4. Exact file manifest

```text
CREATE business_central/official-bc-platform-lock.json
CREATE business_central/project-bc-platform-config.template.json
CREATE business_central/acquire_bccontainerhelper.ps1
CREATE business_central/preflight_business_central.ps1
CREATE business_central/invoke_business_central_dev.ps1
CREATE business_central/test_business_central_runtime.ps1
```

## 5. Materialization blocks

### FILE: `business_central/official-bc-platform-lock.json`
```yaml
block_id: "MICROSOFT-BC-EXECUTABLE:lock:v1"
operation: CREATE
provenance: AUTHORED
source: "local exact lock derived from Microsoft Gallery, BCArtifacts and BCApps evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "000ffcaa6ba1e1ef0a3ec449d6e5f56e3899b7f81d8d9d2b33f9cd3b6233711f"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-official-microsoft-business-central-platform-lock/v1",
  "verified_at": "2026-08-26",
  "bccontainerhelper": {
    "owner": "Microsoft",
    "repository": "microsoft/navcontainerhelper",
    "version": "6.1.16",
    "published_at": "2026-08-21T10:56:10.073Z",
    "package_url": "https://www.powershellgallery.com/api/v2/package/BcContainerHelper/6.1.16",
    "package_bytes": 2886213,
    "package_sha256": "c20bae134a0cbd15927c062fd825234465ec8d189fe459688d282e11442bce10",
    "package_sha512_base64": "pGCpxgjcoUo8brljHY2ehsY0jyBUahhz4IXw97b9OVO5xE39Y6D74HQTWtCeDQoWrqewQJelN6f8Rzl/ouSfow==",
    "license_expression": "MIT",
    "license_path": "LICENSE",
    "license_sha256": "9906940f61b1f0b533fa7d99baf55178b2808fbe113ea51dfbfad8572ccd5f2b",
    "manifest_path": "BcContainerHelper.psd1",
    "manifest_sha256": "9667c6371b87d2c499607294a6c80dc18a20fd7e1964be17e2e5fdb327318892",
    "nuspec_path": "BcContainerHelper.nuspec",
    "nuspec_sha256": "e3fddfd506394a9bd8373a8a369e6c596415cd316a672a2621cae3df921a049c",
    "required_commands": [
      "Get-BCArtifactUrl",
      "New-BcContainer",
      "Get-BcContainerAppInfo",
      "Import-TestToolkitToBcContainer",
      "Run-TestsInBcContainer",
      "Run-AlPipeline"
    ]
  },
  "business_central_artifact": {
    "release": "28.4",
    "version": "28.4.53241.0",
    "country": "w1",
    "type": "OnPrem",
    "artifact_url": "https://bcartifacts-exdbf9fwegejdqak.b02.azurefd.net/onprem/28.4.53241.0/w1",
    "release_tag": "releases/28.4/StrictMode",
    "release_commit": "cb07eef12935e07dd4258c122d0e7b1920fc1223",
    "release_commit_signed": true,
    "release_source_contains_base_application": false,
    "eula_required": true,
    "local_demo_license_allowed": true
  },
  "bcapps_source": {
    "repository": "microsoft/BCApps",
    "commit": "2eae56d704a1fd035d104f333602aea7091b7749",
    "tree": "f9846fb1254c6311985c6131bb2af11c7f169e1b",
    "commit_signature_verified": true,
    "archive_url": "https://github.com/microsoft/BCApps/archive/2eae56d704a1fd035d104f333602aea7091b7749.zip",
    "archive_bytes": 226293141,
    "archive_sha256": "f7e984f2a1e9784a351068f42f0a0cefdc8317704b309f2a94fce5f0f91bc5b6",
    "license_expression": "MIT",
    "license_sha256": "c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383",
    "al_files": 36673,
    "inventory_al_files": 958,
    "warehouse_al_files": 355,
    "reservation_test_al_files": 22,
    "inventory_availability_manifest_sha256": "cf0b1c27267f341a912cf2bd69b0d3e8cc4f529cc667777d43d6e601b250189b",
    "inventory_tracking_manifest_sha256": "a1d432cb8b9a995712b57d1b4e2bc6839ae8a2b2d20a7bf5be4a32c4eba0f9f0",
    "warehouse_manifest_sha256": "95fab623e5d4c28935469fddbde128ce7c7cc0f6431113b7bfba8a16d1f08e62",
    "reservation_test_manifest_sha256": "86de5f6d08f4014c1905c371fc93ea3ac5ab47f8e138d5754a6555673ed1c0b2",
    "classification": "CONDITIONAL_PLATFORM_DEVELOPMENT_SNAPSHOT"
  },
  "production_blockers": [
    "Business Central runtime and product use remain subject to Microsoft EULA and licensing",
    "the signed 28.4 release source omits Base Application, APIV2 and full tests; the complete pinned BCApps source is a development snapshot",
    "the Business Central artifact itself is resolved by the official Microsoft module but is not independently mirrored or hash-locked by Elite",
    "local container development requires supported Windows, administrator authority, Windows containers, Docker-compatible runtime, capacity and explicit EULA acceptance",
    "CRONUS/demo licensing is development-only and does not authorize production",
    "country localization, fiscal rules, identities, secrets, data, external services, workload, security and operational evidence are project-specific"
  ]
}
````

### FILE: `business_central/project-bc-platform-config.template.json`
```yaml
block_id: "MICROSOFT-BC-EXECUTABLE:config-template:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed project decision template"
license: "LicenseRef-Workspace-Owner"
sha256: "8bd252fb6944ff39417a2fe9cf748b2d7af507a6df471db9849a9a6f72662fa5"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-project-microsoft-business-central-platform/v1",
  "project_id": "",
  "owner": "",
  "mode": "LOCAL_DEVELOPMENT_DEMO",
  "container_name": "",
  "artifact_version": "28.4.53241.0",
  "artifact_country": "w1",
  "artifact_url": "https://bcartifacts-exdbf9fwegejdqak.b02.azurefd.net/onprem/28.4.53241.0/w1",
  "auth": "NavUserPassword",
  "credential_user_env": "BC_DEV_USER",
  "credential_password_env": "BC_DEV_PASSWORD",
  "memory_limit": "8G",
  "accept_microsoft_eula": false,
  "allow_network": false,
  "allow_container_creation": false,
  "require_administrator": true,
  "require_windows_containers": true,
  "use_demo_license_only": true,
  "include_test_toolkit": true,
  "bcapps_source_receipt": "",
  "required_capabilities": [],
  "required_country_localization": "",
  "required_external_services": [],
  "required_data_migration": [],
  "required_roles_and_users": [],
  "production_blockers_acknowledged": []
}
````

### FILE: `business_central/acquire_bccontainerhelper.ps1`
```yaml
block_id: "MICROSOFT-BC-EXECUTABLE:acquirer:v1"
operation: CREATE
provenance: AUTHORED
source: "local acquisition orchestration; official module remains byte-identical"
license: "LicenseRef-Workspace-Owner"
sha256: "99f92945d453c61fa9178d70ea2bd8b359bf026d513193acea448c0c5916acd1"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param(
  [Parameter(Mandatory=$true)][string]$LockPath,
  [Parameter(Mandatory=$true)][string]$Destination,
  [string]$PackageCachePath,
  [switch]$AllowNetwork
)
$ErrorActionPreference='Stop'
$utf8=[Text.UTF8Encoding]::new($false)
function Fail([string]$Message){ throw "BC_HELPER_ACQUISITION_FAILED: $Message" }
function Hash([string]$Path,[string]$Algorithm='SHA256'){ (Get-FileHash -LiteralPath $Path -Algorithm $Algorithm).Hash.ToLowerInvariant() }
$lockFile=(Resolve-Path -LiteralPath $LockPath).Path
$lock=[IO.File]::ReadAllText($lockFile,$utf8)|ConvertFrom-Json -Depth 20
if($lock.schema -cne 'elite-official-microsoft-business-central-platform-lock/v1'){Fail 'lock schema mismatch'}
$item=$lock.bccontainerhelper
if($item.version -notmatch '^\d+\.\d+\.\d+$' -or $item.package_sha256 -notmatch '^[0-9a-f]{64}$'){Fail 'invalid module identity'}
$dest=[IO.Path]::GetFullPath($Destination)
if(Test-Path -LiteralPath $dest){Fail 'destination must not exist'}
$scratch=Join-Path ([IO.Path]::GetTempPath()) ('elite-bc-helper-'+[guid]::NewGuid().ToString('N'))
try{
  [void](New-Item -ItemType Directory -Path $scratch)
  $package=Join-Path $scratch 'BcContainerHelper.nupkg'
  if($PackageCachePath){
    $cache=(Resolve-Path -LiteralPath $PackageCachePath).Path
    [IO.File]::Copy($cache,$package,$false)
  }else{
    if(-not $AllowNetwork){Fail 'AllowNetwork or PackageCachePath is required'}
    Invoke-WebRequest -Uri ([string]$item.package_url) -OutFile $package
  }
  if((Get-Item -LiteralPath $package).Length -ne [long]$item.package_bytes){Fail 'package size mismatch'}
  if((Hash $package) -cne [string]$item.package_sha256){Fail 'package SHA-256 mismatch'}
  $sha512=[Convert]::ToBase64String([Security.Cryptography.SHA512]::HashData([IO.File]::ReadAllBytes($package)))
  if($sha512 -cne [string]$item.package_sha512_base64){Fail 'package Gallery SHA-512 mismatch'}
  $archive=[IO.Compression.ZipFile]::OpenRead($package)
  try{
    foreach($entry in $archive.Entries){
      $relative=$entry.FullName.Replace('\','/')
      if([IO.Path]::IsPathRooted($relative) -or $relative.Split('/') -contains '..'){Fail "unsafe package entry: $relative"}
    }
  }finally{$archive.Dispose()}
  Expand-Archive -LiteralPath $package -DestinationPath $dest
  foreach($pair in @(
    @([string]$item.license_path,[string]$item.license_sha256),
    @([string]$item.manifest_path,[string]$item.manifest_sha256),
    @([string]$item.nuspec_path,[string]$item.nuspec_sha256)
  )){
    $path=Join-Path $dest $pair[0]
    if(-not(Test-Path -LiteralPath $path -PathType Leaf) -or (Hash $path) -cne $pair[1]){Fail "critical file mismatch: $($pair[0])"}
  }
  $manifest=Test-ModuleManifest -Path (Join-Path $dest $item.manifest_path)
  if($manifest.Version.ToString() -cne [string]$item.version){Fail 'module version mismatch'}
  $receipt=[ordered]@{schema='elite-official-bc-helper-receipt/v1';version=[string]$item.version;lock_sha256=Hash $lockFile;package_sha256=Hash $package;acquired_at=[DateTimeOffset]::UtcNow.ToString('O');production_status='BLOCKED_PENDING_PROJECT_PREFLIGHT'}
  [IO.File]::WriteAllText((Join-Path $dest 'ELITE_ACQUISITION_RECEIPT.json'),($receipt|ConvertTo-Json -Depth 10)+"`n",$utf8)
  Write-Output "BC_HELPER_ACQUIRED version=$($item.version) destination=$dest"
}catch{
  if(Test-Path -LiteralPath $dest){[IO.Directory]::Delete($dest,$true)}
  throw
}finally{
  if(Test-Path -LiteralPath $scratch){[IO.Directory]::Delete($scratch,$true)}
}
````

### FILE: `business_central/preflight_business_central.ps1`
```yaml
block_id: "MICROSOFT-BC-EXECUTABLE:preflight:v1"
operation: CREATE
provenance: AUTHORED
source: "local admission gate invoking official module validation"
license: "LicenseRef-Workspace-Owner"
sha256: "f6db964259a19a468b0361a77c79c3432e5971627afd05cc4a64ee4950e8840f"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param(
  [Parameter(Mandatory=$true)][string]$ConfigPath,
  [Parameter(Mandatory=$true)][string]$LockPath,
  [Parameter(Mandatory=$true)][string]$ModuleRoot,
  [switch]$ValidateConfigOnly
)
$ErrorActionPreference='Stop'
$utf8=[Text.UTF8Encoding]::new($false)
function Fail([string]$Message){throw "BC_PLATFORM_PREFLIGHT_FAILED: $Message"}
$configFile=(Resolve-Path -LiteralPath $ConfigPath).Path
$lockFile=(Resolve-Path -LiteralPath $LockPath).Path
$moduleDir=(Resolve-Path -LiteralPath $ModuleRoot).Path
$config=[IO.File]::ReadAllText($configFile,$utf8)|ConvertFrom-Json -Depth 20
$lock=[IO.File]::ReadAllText($lockFile,$utf8)|ConvertFrom-Json -Depth 20
if($config.schema -cne 'elite-project-microsoft-business-central-platform/v1'){Fail 'config schema mismatch'}
if($lock.schema -cne 'elite-official-microsoft-business-central-platform-lock/v1'){Fail 'lock schema mismatch'}
if([string]::IsNullOrWhiteSpace($config.project_id) -or [string]::IsNullOrWhiteSpace($config.owner)){Fail 'project_id and owner are required'}
if($config.mode -cne 'LOCAL_DEVELOPMENT_DEMO'){Fail 'only LOCAL_DEVELOPMENT_DEMO is admitted'}
if($config.artifact_version -cne $lock.business_central_artifact.version -or $config.artifact_url -cne $lock.business_central_artifact.artifact_url){Fail 'artifact must match exact lock'}
if($config.artifact_country -cne $lock.business_central_artifact.country){Fail 'artifact country must match exact lock'}
if($config.auth -cne 'NavUserPassword' -or $config.credential_user_env -notmatch '^[A-Z][A-Z0-9_]+$' -or $config.credential_password_env -notmatch '^[A-Z][A-Z0-9_]+$'){Fail 'credential environment references are invalid'}
if($config.container_name -notmatch '^[a-z][a-z0-9-]{2,30}$'){Fail 'container_name is invalid'}
if(-not $config.use_demo_license_only){Fail 'this pack admits only the local demo license path'}
if($config.production_blockers_acknowledged.Count -ne $lock.production_blockers.Count){Fail 'all exact production blockers must be acknowledged'}
for($i=0;$i -lt $lock.production_blockers.Count;$i++){if($config.production_blockers_acknowledged[$i] -cne $lock.production_blockers[$i]){Fail "production blocker mismatch at index $i"}}
if($ValidateConfigOnly){Write-Output "BC_PLATFORM_CONFIG_VALID version=$($config.artifact_version)";exit 0}
$manifest=Join-Path $moduleDir ([string]$lock.bccontainerhelper.manifest_path)
if(-not(Test-Path -LiteralPath $manifest -PathType Leaf)){Fail 'BcContainerHelper manifest is missing'}
$tested=Test-ModuleManifest -Path $manifest
if($tested.Version.ToString() -cne [string]$lock.bccontainerhelper.version){Fail 'BcContainerHelper version mismatch'}
if(-not $IsWindows){Fail 'Business Central local container lane requires Windows'}
if(-not $config.accept_microsoft_eula){Fail 'Microsoft EULA acceptance is required from the user'}
if(-not $config.allow_network){Fail 'network authority is required to obtain the Microsoft artifact'}
if(-not $config.allow_container_creation){Fail 'container creation authority is required'}
$principal=[Security.Principal.WindowsPrincipal]::new([Security.Principal.WindowsIdentity]::GetCurrent())
if($config.require_administrator -and -not $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)){Fail 'administrator authority is required'}
$docker=Get-Command docker -ErrorAction SilentlyContinue
if(-not $docker){Fail 'docker-compatible runtime is missing'}
$osType=(& docker info --format '{{.OSType}}' 2>$null).Trim()
if($LASTEXITCODE -ne 0){Fail 'docker daemon is unavailable'}
if($config.require_windows_containers -and $osType -cne 'windows'){Fail 'docker daemon is not in Windows container mode'}
if([string]::IsNullOrWhiteSpace([Environment]::GetEnvironmentVariable([string]$config.credential_user_env)) -or [string]::IsNullOrWhiteSpace([Environment]::GetEnvironmentVariable([string]$config.credential_password_env))){Fail 'credential environment variables are not populated'}
Import-Module $manifest -Force -DisableNameChecking
foreach($command in $lock.bccontainerhelper.required_commands){if(-not(Get-Command $command -ErrorAction SilentlyContinue)){Fail "required official command missing: $command"}}
Write-Output "BC_PLATFORM_PREFLIGHT_PASS version=$($config.artifact_version) container=$($config.container_name)"
````

### FILE: `business_central/invoke_business_central_dev.ps1`
```yaml
block_id: "MICROSOFT-BC-EXECUTABLE:invoke-dev:v1"
operation: CREATE
provenance: AUTHORED
source: "local approval wrapper around official Microsoft BcContainerHelper commands"
license: "LicenseRef-Workspace-Owner"
sha256: "4069b7f4def5b403ce5e7299225335254bce28546587d980c7e2e8bf8b850196"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param(
  [Parameter(Mandatory=$true)][string]$ConfigPath,
  [Parameter(Mandatory=$true)][string]$LockPath,
  [Parameter(Mandatory=$true)][string]$ModuleRoot,
  [Parameter(Mandatory=$true)][string]$ApprovalPath,
  [switch]$Execute
)
$ErrorActionPreference='Stop'
$utf8=[Text.UTF8Encoding]::new($false)
function Fail([string]$Message){throw "BC_PLATFORM_INVOCATION_FAILED: $Message"}
function Hash([string]$Path){(Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()}
if(-not $Execute){Fail 'Execute is required; dry invocation cannot create a container'}
$configFile=(Resolve-Path -LiteralPath $ConfigPath).Path
$lockFile=(Resolve-Path -LiteralPath $LockPath).Path
$approvalFile=(Resolve-Path -LiteralPath $ApprovalPath).Path
$moduleDir=(Resolve-Path -LiteralPath $ModuleRoot).Path
$config=[IO.File]::ReadAllText($configFile,$utf8)|ConvertFrom-Json -Depth 20
$lock=[IO.File]::ReadAllText($lockFile,$utf8)|ConvertFrom-Json -Depth 20
$approval=[IO.File]::ReadAllText($approvalFile,$utf8)|ConvertFrom-Json -Depth 20
if($approval.schema -cne 'elite-project-microsoft-business-central-execution-approval/v1'){Fail 'approval schema mismatch'}
if($approval.config_sha256 -cne (Hash $configFile) -or $approval.lock_sha256 -cne (Hash $lockFile)){Fail 'approval hash mismatch'}
foreach($flag in @('accept_microsoft_eula','allow_network','allow_container_creation')){if($approval.$flag -ne $true){Fail "approval flag is not true: $flag"}}
if($approval.production_use -ne $false){Fail 'this pack does not authorize production use'}
& (Join-Path $PSScriptRoot 'preflight_business_central.ps1') -ConfigPath $configFile -LockPath $lockFile -ModuleRoot $moduleDir
if(-not $?){Fail 'preflight failed'}
$manifest=Join-Path $moduleDir ([string]$lock.bccontainerhelper.manifest_path)
Import-Module $manifest -Force -DisableNameChecking
$user=[Environment]::GetEnvironmentVariable([string]$config.credential_user_env)
$password=[Environment]::GetEnvironmentVariable([string]$config.credential_password_env)|ConvertTo-SecureString -AsPlainText -Force
$credential=[Management.Automation.PSCredential]::new($user,$password)
New-BcContainer -accept_eula -containerName ([string]$config.container_name) -artifactUrl ([string]$config.artifact_url) -auth NavUserPassword -Credential $credential -memoryLimit ([string]$config.memory_limit) -includeTestToolkit -includeTestLibrariesOnly -shortcuts None -updateHosts:$false
$apps=@(Get-BcContainerAppInfo -containerName ([string]$config.container_name))
if($apps.Count -eq 0){Fail 'official container returned no installed applications'}
$receipt=[ordered]@{schema='elite-project-microsoft-business-central-dev-receipt/v1';project_id=[string]$config.project_id;container_name=[string]$config.container_name;artifact_version=[string]$config.artifact_version;artifact_url=[string]$config.artifact_url;bccontainerhelper_version=[string]$lock.bccontainerhelper.version;installed_app_count=$apps.Count;created_at=[DateTimeOffset]::UtcNow.ToString('O');production_status='DEVELOPMENT_ONLY_NOT_PRODUCTION'}
$receiptPath=Join-Path ([IO.Path]::GetDirectoryName($configFile)) 'BUSINESS_CENTRAL_DEV_RECEIPT.json'
[IO.File]::WriteAllText($receiptPath,($receipt|ConvertTo-Json -Depth 10)+"`n",$utf8)
Write-Output "BC_PLATFORM_DEV_CREATED container=$($config.container_name) apps=$($apps.Count) receipt=$receiptPath"
````

### FILE: `business_central/test_business_central_runtime.ps1`
```yaml
block_id: "MICROSOFT-BC-EXECUTABLE:test:v1"
operation: CREATE
provenance: AUTHORED
source: "local offline regression suite"
license: "LicenseRef-Workspace-Owner"
sha256: "ccca355770c15a192926a92961f75545319df1b582a96741a97844c6e936a73b"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$utf8=[Text.UTF8Encoding]::new($false)
$root=[IO.Path]::GetFullPath($PSScriptRoot)
$scratch=Join-Path ([IO.Path]::GetTempPath()) ('elite-bc-runtime-test-'+[guid]::NewGuid().ToString('N'))
function Expect-Failure([scriptblock]$Action,[string]$Expected){try{& $Action;throw "expected failure: $Expected"}catch{if($_.Exception.Message -notlike "*$Expected*"){throw}}}
try{
  [void](New-Item -ItemType Directory -Path $scratch)
  $lockPath=Join-Path $root 'official-bc-platform-lock.json'
  $templatePath=Join-Path $root 'project-bc-platform-config.template.json'
  $lock=[IO.File]::ReadAllText($lockPath,$utf8)|ConvertFrom-Json -Depth 20
  if($lock.bccontainerhelper.version -cne '6.1.16' -or $lock.business_central_artifact.version -cne '28.4.53241.0'){throw 'exact platform lock drift'}
  if($lock.bcapps_source.commit -cne '2eae56d704a1fd035d104f333602aea7091b7749' -or $lock.bcapps_source.tree -cne 'f9846fb1254c6311985c6131bb2af11c7f169e1b' -or $lock.bcapps_source.archive_sha256 -cne 'f7e984f2a1e9784a351068f42f0a0cefdc8317704b309f2a94fce5f0f91bc5b6' -or $lock.bcapps_source.commit_signature_verified -ne $true){throw 'exact BCApps source lock drift'}
  if($lock.bcapps_source.inventory_al_files -ne 958 -or $lock.bcapps_source.warehouse_al_files -ne 355 -or $lock.bcapps_source.reservation_test_al_files -ne 22 -or $lock.bcapps_source.inventory_availability_manifest_sha256 -cne 'cf0b1c27267f341a912cf2bd69b0d3e8cc4f529cc667777d43d6e601b250189b' -or $lock.bcapps_source.inventory_tracking_manifest_sha256 -cne 'a1d432cb8b9a995712b57d1b4e2bc6839ae8a2b2d20a7bf5be4a32c4eba0f9f0' -or $lock.bcapps_source.warehouse_manifest_sha256 -cne '95fab623e5d4c28935469fddbde128ce7c7cc0f6431113b7bfba8a16d1f08e62' -or $lock.bcapps_source.reservation_test_manifest_sha256 -cne '86de5f6d08f4014c1905c371fc93ea3ac5ab47f8e138d5754a6555673ed1c0b2'){throw 'BCApps inventory scope drift'}
  if($lock.bccontainerhelper.required_commands.Count -ne 6 -or $lock.production_blockers.Count -lt 6){throw 'lock gates are incomplete'}
  $config=[IO.File]::ReadAllText($templatePath,$utf8)|ConvertFrom-Json -Depth 20
  $config.project_id='test-project';$config.owner='test-owner';$config.container_name='elite-bc-test'
  $config.production_blockers_acknowledged=@($lock.production_blockers)
  $configPath=Join-Path $scratch 'config.json'
  [IO.File]::WriteAllText($configPath,($config|ConvertTo-Json -Depth 20)+"`n",$utf8)
  & (Join-Path $root 'preflight_business_central.ps1') -ConfigPath $configPath -LockPath $lockPath -ModuleRoot $root -ValidateConfigOnly
  Expect-Failure { & (Join-Path $root 'invoke_business_central_dev.ps1') -ConfigPath $configPath -LockPath $lockPath -ModuleRoot $root -ApprovalPath $configPath } 'Execute is required'
  Expect-Failure { & (Join-Path $root 'acquire_bccontainerhelper.ps1') -LockPath $lockPath -Destination (Join-Path $scratch 'module') } 'AllowNetwork or PackageCachePath is required'
  foreach($script in @('acquire_bccontainerhelper.ps1','preflight_business_central.ps1','invoke_business_central_dev.ps1')){
    $text=[IO.File]::ReadAllText((Join-Path $root $script),$utf8)
    if($text -match '(?m)^\s*foreach\s*\([^\r\n]+\)\s*\{[^\r\n]*\}\s*\|'){throw "recurring foreach pipeline pattern: $script"}
    if($text -match '(?i)\bgit\s+(clone|init|checkout)\b'){throw "Git command forbidden: $script"}
  }
  Write-Output 'BC_PLATFORM_RUNTIME_TEST_PASS positives=1 negatives=2 regression=1'
}finally{if(Test-Path -LiteralPath $scratch){[IO.Directory]::Delete($scratch,$true)}}
````

## 6. Configuration surface

| Campo | Tipo | Default | Validación | Secreto | Efecto |
|---|---|---|---|---|---|
| `mode` | enum | `LOCAL_DEVELOPMENT_DEMO` | único modo admitido | no | impide claim productivo |
| `artifact_version/url/country` | exactos | lock 28.4 W1 | coincidencia byte de strings | no | runtime Microsoft |
| `container_name` | slug | vacío | `^[a-z][a-z0-9-]{2,30}$` | no | identidad local |
| `credential_*_env` | env refs | `BC_DEV_*` | nombre, valor sólo en entorno | referencia | usuario demo |
| `accept_microsoft_eula` | bool | false | usuario explícito | no | habilita invocación |
| `allow_network/container_creation` | bool | false | approval exacto | no | mutación externa/local |
| `production_blockers_acknowledged` | array | vacío | igualdad y orden con lock | no | visibilidad de condiciones |

## 7. Dependency bill

| Artefacto | Pin | Licencia/términos | Uso | Estado |
|---|---|---|---|---|
| BcContainerHelper | 6.1.16, package SHA-256 `c20bae13...ce10`, Gallery SHA-512 fijado | MIT | resolver/crear/probar contenedor | package/import/309 comandos y artifact resolution verificados |
| Business Central artifact | 28.4.53241.0 W1 OnPrem, URL exacta observada | Microsoft EULA; demo sólo desarrollo | producto ejecutable | no descargado ni creado por esta auditoría |
| BCApps | commit firmado `2eae56d...`, tree/archive/path manifests fijados | MIT por source; runtime condicionado | 36.673 AL, incluidos Inventory 958, Warehouse 355 y 22 tests SCM-Reservation | source auditado, no release consolidada ni runtime portable |
| Windows container runtime | selección del usuario | términos del runtime | host local | administrador/Docker no disponibles en esta auditoría |

## 8. Apply order

1. El usuario elige `BUSINESS_CENTRAL_PLATFORM` y completa la configuración sin secretos.
2. Ejecutar `test_business_central_runtime.ps1`.
3. Adquirir BcContainerHelper por cache verificada o `-AllowNetwork`.
4. Ejecutar preflight completo; no corregir automáticamente permisos, Docker ni hosts.
5. El usuario aprueba config/lock hashes, EULA, red y creación de contenedor.
6. Sólo entonces invocar con `-Execute`; conservar receipt y ejecutar suites oficiales seleccionadas.
7. Production permanece bloqueada hasta licencia, localización, seguridad, datos, restore, carga y operación reales.

## 9. Verification

```text
pwsh -NoProfile -File business_central/test_business_central_runtime.ps1
expected: BC_PLATFORM_RUNTIME_TEST_PASS positives=1 negatives=2 regression=1

pwsh -NoProfile -File business_central/acquire_bccontainerhelper.ps1 ... -PackageCachePath <exact-nupkg>
expected: BC_HELPER_ACQUIRED version=6.1.16

pwsh -NoProfile -File business_central/preflight_business_central.ps1 ... -ValidateConfigOnly
expected: BC_PLATFORM_CONFIG_VALID version=28.4.53241.0
```

El preflight host completo debe fallar si no hay Windows, administrador, Docker en modo Windows, secretos de desarrollo, red o autorización. La creación real no fue ejecutada en esta evidencia.

## 10. Reconstruction evidence

`reconstruction_evidence/MICROSOFT_BUSINESS_CENTRAL_EXECUTABLE_PLATFORM_2026-08-26_V1.md` fija Gallery metadata, hashes, import, comandos, artifact resolution, límites de la release y fallos observados. `reconstruction_evidence/ALL_IMPLEMENTATION_PACKS_MATERIALIZATION_2026-08-26_V16.md` gobierna conteos y reconstrucción. Estado: `REBUILD_VERIFIED / CONDITIONED`, nunca `REUSABLE_PACK` productivo.
