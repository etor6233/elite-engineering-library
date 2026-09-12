# Microsoft-generated ARCA WSAA and WSFEv1 Clients

## 1. Metadata

```yaml
pack_id: "MICROSOFT-ARCA-WSFE-GENERATED-CLIENT"
pack_version: "0.2.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Genera y compila proxies WSAA y WSFEv1 reproducibles desde WSDL públicos ARCA exactos usando un único tooling Microsoft fijado, sin redistribuir contratos ni inventar protocolos."
stacks: ["PowerShell 7", ".NET SDK 10.0.400", "Microsoft dotnet-svcutil 8.0.0", "WCF 10.0.652802"]
compatible_with: ["ARCA-WSAA-OFFICIAL-CLIENT-SAMPLES@0.1.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner; generated/upstream artifacts retain their own terms"
upstream_sources: ["https://www.arca.gob.ar/fe/ayuda/webservice.asp", "https://www.arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG.pdf", "https://learn.microsoft.com/dotnet/core/tools/dotnet-svcutil-guide", "https://github.com/dotnet/wcf"]
verified_at: "2026-08-30"
```

All six materialized files are local `AUTHORED` orchestration. The runner acquires but does not embed public ARCA WSDL bytes, and generates two namespace-isolated clients with the official Microsoft tool. Neither the runner nor either generated proxy is represented as ARCA-authored product code.

## 2. Applicability

Use when an Argentine project selects WSAA/WSFEv1 and needs reproducible strongly typed transport clients before implementing its credential and fiscal workflow. Do not use it to infer taxes, invoice validity, authorization, private-key controls or production admission.

## 3. Architecture contract

The lock separates four ARCA endpoints, exact bytes and four generated hashes. Homologation is the default; production metadata requires explicit acknowledgement. The runner prohibits DTDs, verifies required operations, generates WSAA and WSFEv1 into separate namespaces in one project, fixes the tool/runtime graph, builds with warnings as errors, runs NuGet vulnerability audit and compares both generated hashes. It never calls an authentication or fiscal operation.

## 4. Exact file manifest

```text
CREATE arca/fiscal/README.md
CREATE arca/fiscal/arca-ws-contracts.lock.json
CREATE arca/fiscal/dotnet-tools.json
CREATE arca/fiscal/Elite.Arca.Wsfe.Generated.csproj
CREATE tools/generate-arca-wsfe-client.ps1
CREATE tools/test-arca-wsfe-generator.ps1
```

## 5. Materialization blocks

### FILE: `arca/fiscal/README.md`

```yaml
block_id: "ARCA-WSFE-GEN:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local admission and usage boundary governed by ARCA/Microsoft authorities"
license: "LicenseRef-Workspace-Owner"
sha256: "60b8c482b309a75c7e50e4e373496b67f1c95d41033d3728db5224ed0c622256"
variables: []
secrets_allowed: false
```
````text
# ARCA WSAA and WSFEv1 generated clients

This component generates separate WCF clients from exact public ARCA WSAA and WSFEv1 WSDLs. It does not embed either WSDL or claim that generated code is an ARCA SDK.

Required inputs are an official .NET 10.0.400 SDK executable, a cache directory and an empty output directory. Homologation is the default and production metadata requires an explicit switch. The runner fixes `dotnet-svcutil` 8.0.0, WCF 10.0.652802 and patched cryptography 10.0.11, builds with warnings as errors and rejects vulnerable packages reported by current NuGet sources.

This is client generation only. It does not provide WSAA credentials, private-key custody, fiscal rules, authorization, idempotency, reconciliation or ARCA homologation.
````

### FILE: `arca/fiscal/arca-ws-contracts.lock.json`

```yaml
block_id: "ARCA-WSFE-GEN:lock:v1"
operation: CREATE
provenance: AUTHORED
source: "local lock of public ARCA endpoints and Microsoft toolchain"
license: "LicenseRef-Workspace-Owner"
sha256: "896d12457910a527305c7eb2dbb3dc67a1edba35e9b8d91db10cdd85d8daace8"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-arca-contract-lock/v1",
  "verifiedAt": "2026-08-30",
  "toolchain": {
    "dotnetSdk": "10.0.400",
    "dotnetSvcutil": "8.0.0",
    "systemServiceModelHttp": "10.0.652802",
    "systemSecurityCryptographyXml": "10.0.11",
    "systemSecurityCryptographyPkcs": "10.0.11"
  },
  "contracts": {
    "wsaaHomologation": {"url":"https://wsaahomo.afip.gov.ar/ws/services/LoginCms?WSDL","bytes":3222,"sha256":"c8910060052496f248b125e676085db689e0f21b2f32728417725c376b97b4b3","generatedSha256":"910a9e267718e663424972d2dd10d1f757e14f617e8a0d1afab8580bd6c450ea"},
    "wsaaProduction": {"url":"https://wsaa.afip.gov.ar/ws/services/LoginCms?WSDL","bytes":3198,"sha256":"2f2a3116ee6ae48b3d777f5a3c94219b106b05bb8174457e56be475bd7f9b192","generatedSha256":"71c36f4420b645365e06b35e76d4f856ac0094aff87aa80c615e26195bf3fb02"},
    "wsfeHomologation": {"url":"https://wswhomo.afip.gov.ar/wsfev1/service.asmx?WSDL","bytes":77307,"sha256":"16f5b8450c35a2cfc0187f07d6236162bb57b663954a1e1f0267e71afe90af24","generatedSha256":"e4608d9e113f293945a0fadf7a8c5f018ae0a2e19568ea450772be8bd1ec9d61"},
    "wsfeProduction": {"url":"https://servicios1.afip.gov.ar/wsfev1/service.asmx?WSDL","bytes":77313,"sha256":"10b78b714ea3bf95f5c12ce2fa05805bb21b9fd13f9a98b57fbb866ee9ebd43e","generatedSha256":"85c3fa927eec659e952b9a8982e2a21368fee29808c3e32fa1fd5704749fb021"}
  }
}
````

### FILE: `arca/fiscal/dotnet-tools.json`

```yaml
block_id: "ARCA-WSFE-GEN:tools:v1"
operation: CREATE
provenance: AUTHORED
source: "local exact Microsoft tool selection"
license: "LicenseRef-Workspace-Owner"
sha256: "a17261c17642d484c3f16b62e1c0be511df6d06923ed8f0604ef395f8b525048"
variables: []
secrets_allowed: false
```
````json
{
  "version": 1,
  "isRoot": true,
  "tools": {
    "dotnet-svcutil": {
      "version": "8.0.0",
      "commands": ["dotnet-svcutil"],
      "rollForward": false
    }
  }
}
````

### FILE: `arca/fiscal/Elite.Arca.Wsfe.Generated.csproj`

```yaml
block_id: "ARCA-WSFE-GEN:project:v1"
operation: CREATE
provenance: AUTHORED
source: "local exact Microsoft runtime and security package selection"
license: "LicenseRef-Workspace-Owner"
sha256: "ba6cab230772640178dca57783316b9e51adb8cfaa482e53ba5fd16a09ff1c45"
variables: []
secrets_allowed: false
```
````xml
<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <TargetFramework>net10.0</TargetFramework>
    <ImplicitUsings>enable</ImplicitUsings>
    <Nullable>enable</Nullable>
    <TreatWarningsAsErrors>true</TreatWarningsAsErrors>
    <RestorePackagesWithLockFile>true</RestorePackagesWithLockFile>
    <NuGetAudit>true</NuGetAudit>
    <NuGetAuditMode>all</NuGetAuditMode>
  </PropertyGroup>
  <ItemGroup>
    <PackageReference Include="System.Security.Cryptography.Pkcs" Version="10.0.11" />
    <PackageReference Include="System.Security.Cryptography.Xml" Version="10.0.11" />
    <PackageReference Include="System.ServiceModel.Http" Version="10.0.652802" />
  </ItemGroup>
</Project>
````

### FILE: `tools/generate-arca-wsfe-client.ps1`

```yaml
block_id: "ARCA-WSFE-GEN:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed WSDL acquisition, generation and verification"
license: "LicenseRef-Workspace-Owner"
sha256: "bffcc927ffc9a91e15c5c91d2eae078227f042b6a833c39f66a834e650af8caa"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param(
  [Parameter(Mandatory=$true)][string]$DotnetExecutable,
  [Parameter(Mandatory=$true)][string]$CacheDirectory,
  [Parameter(Mandatory=$true)][string]$OutputDirectory,
  [ValidateSet('Homologation','Production')][string]$Environment='Homologation',
  [switch]$AllowProductionMetadata,
  [switch]$Offline
)
$ErrorActionPreference='Stop'
function Hash([string]$Path){(Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()}
function Run([string]$Exe,[string[]]$Arguments){$lines=@(& $Exe @Arguments 2>&1);if($LASTEXITCODE-ne 0){throw "command failed ($LASTEXITCODE): $Exe $($Arguments-join' ')`n$($lines-join[Environment]::NewLine)"};$lines}
$dotnet=[IO.Path]::GetFullPath($DotnetExecutable)
if(-not(Test-Path -LiteralPath $dotnet -PathType Leaf)){throw 'dotnet executable missing'}
$assets=[IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\arca\fiscal'))
$lock=Get-Content -LiteralPath (Join-Path $assets 'arca-ws-contracts.lock.json') -Raw|ConvertFrom-Json -Depth 20
$actualSdk=(@(Run $dotnet @('--version'))[-1]).Trim()
if($actualSdk-ne$lock.toolchain.dotnetSdk){throw "dotnet SDK mismatch: expected=$($lock.toolchain.dotnetSdk) actual=$actualSdk"}
if($Environment-eq'Production'-and-not$AllowProductionMetadata){throw 'production metadata requires -AllowProductionMetadata'}
$wsaaContract=if($Environment-eq'Production'){$lock.contracts.wsaaProduction}else{$lock.contracts.wsaaHomologation}
$wsfeContract=if($Environment-eq'Production'){$lock.contracts.wsfeProduction}else{$lock.contracts.wsfeHomologation}
$cache=[IO.Path]::GetFullPath($CacheDirectory);$output=[IO.Path]::GetFullPath($OutputDirectory)
if($cache-eq$output){throw 'cache and output must differ'}
[IO.Directory]::CreateDirectory($cache)|Out-Null
if(Test-Path -LiteralPath $output){if((Get-ChildItem -Force -LiteralPath $output).Count-ne0){throw 'output must be absent or empty'}}else{[IO.Directory]::CreateDirectory($output)|Out-Null}
function Acquire([string]$Name,$Contract,[string[]]$Operations){
  $path=Join-Path $cache ("{0}-{1}.wsdl" -f $Name,$Environment.ToLowerInvariant())
  if(-not(Test-Path -LiteralPath $path -PathType Leaf)){if($Offline){throw "offline WSDL cache miss: $Name"};Invoke-WebRequest -UseBasicParsing -Uri $Contract.url -OutFile $path -TimeoutSec 30}
  if((Get-Item -LiteralPath $path).Length-ne[long]$Contract.bytes-or(Hash $path)-ne$Contract.sha256){throw "WSDL identity mismatch: $Name"}
  $settings=[Xml.XmlReaderSettings]::new();$settings.DtdProcessing=[Xml.DtdProcessing]::Prohibit;$settings.XmlResolver=$null
  $reader=[Xml.XmlReader]::Create($path,$settings);try{while($reader.Read()){}}finally{$reader.Dispose()}
  $text=[IO.File]::ReadAllText($path);foreach($operation in $Operations){if($text-notmatch[regex]::Escape($operation)){throw "required operation missing: $Name/$operation"}}
  $path
}
$wsaaWsdl=Acquire 'wsaa' $wsaaContract @('loginCms')
$wsfeWsdl=Acquire 'wsfe' $wsfeContract @('FECAESolicitar','FECompUltimoAutorizado','FEParamGetTiposCbte')
Copy-Item -LiteralPath (Join-Path $assets 'dotnet-tools.json') -Destination (Join-Path $output 'dotnet-tools.json')
$project=Join-Path $output 'Elite.Arca.Wsfe.Generated.csproj';Copy-Item -LiteralPath (Join-Path $assets 'Elite.Arca.Wsfe.Generated.csproj') -Destination $project
$oldPath=$env:PATH;$oldRoll=$env:DOTNET_ROLL_FORWARD;$oldCli=$env:DOTNET_CLI_TELEMETRY_OPTOUT;$oldSvc=$env:DOTNET_SVCUTIL_TELEMETRY_OPTOUT
try{
  $env:PATH="$([IO.Path]::GetDirectoryName($dotnet));$oldPath";$env:DOTNET_ROLL_FORWARD='Major';$env:DOTNET_CLI_TELEMETRY_OPTOUT='1';$env:DOTNET_SVCUTIL_TELEMETRY_OPTOUT='1'
  $manifest=Join-Path $output 'dotnet-tools.json';$restoreArgs=@('tool','restore','--tool-manifest',$manifest);if($Offline){$restoreArgs+='--ignore-failed-sources'};Run $dotnet $restoreArgs|Out-Null
  Push-Location $output
  try{
    Run $dotnet @('tool','run','dotnet-svcutil','--',$wsaaWsdl,'--projectFile',$project,'--outputDir',(Join-Path $output 'Generated'),'--outputFile','WsaaReference.cs','--namespace','*,Elite.Arca.Wsaa.Generated','--serializer','XmlSerializer','--noLogo','--verbosity','Minimal')|Out-Null
    Run $dotnet @('tool','run','dotnet-svcutil','--',$wsfeWsdl,'--projectFile',$project,'--outputDir',(Join-Path $output 'Generated'),'--outputFile','WsfeV1Reference.cs','--namespace','*,Elite.Arca.Wsfe.Generated','--serializer','XmlSerializer','--noLogo','--verbosity','Minimal')|Out-Null
  }finally{Pop-Location}
  Copy-Item -LiteralPath (Join-Path $assets 'Elite.Arca.Wsfe.Generated.csproj') -Destination $project -Force
  Run $dotnet @('restore',$project,'--force','--no-cache','--warnaserror')|Out-Null
  Run $dotnet @('build',$project,'--configuration','Release','--no-restore','--warnaserror')|Out-Null
  $audit=Run $dotnet @('list',$project,'package','--vulnerable','--include-transitive')
  if(($audit-join"`n")-notmatch'no vulnerable packages'){throw "NuGet vulnerability gate did not return clean`n$($audit-join[Environment]::NewLine)"}
}finally{$env:PATH=$oldPath;$env:DOTNET_ROLL_FORWARD=$oldRoll;$env:DOTNET_CLI_TELEMETRY_OPTOUT=$oldCli;$env:DOTNET_SVCUTIL_TELEMETRY_OPTOUT=$oldSvc}
$wsaaGenerated=Join-Path $output 'Generated\WsaaReference.cs';if((Hash $wsaaGenerated)-ne$wsaaContract.generatedSha256){throw 'generated proxy identity mismatch: WSAA'}
$wsfeGenerated=Join-Path $output 'Generated\WsfeV1Reference.cs';if((Hash $wsfeGenerated)-ne$wsfeContract.generatedSha256){throw 'generated proxy identity mismatch: WSFE'}
$receipt=[ordered]@{schema='elite-arca-generated-client/v2';environment=$Environment;contracts=@([ordered]@{service='wsaa';url=$wsaaContract.url;wsdl_sha256=$wsaaContract.sha256;generated_sha256=$wsaaContract.generatedSha256},[ordered]@{service='wsfev1';url=$wsfeContract.url;wsdl_sha256=$wsfeContract.sha256;generated_sha256=$wsfeContract.generatedSha256});dotnet_sdk=$actualSdk;dotnet_svcutil=$lock.toolchain.dotnetSvcutil;production_admitted=$false}
[IO.File]::WriteAllText((Join-Path $output 'GENERATION_RECEIPT.json'),($receipt|ConvertTo-Json -Depth 10),[Text.UTF8Encoding]::new($false))
"ARCA_GENERATION_PASS environment=$Environment generated_clients=2 production_admitted=false"
````

### FILE: `tools/test-arca-wsfe-generator.ps1`

```yaml
block_id: "ARCA-WSFE-GEN:test:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract regression"
license: "LicenseRef-Workspace-Owner"
sha256: "2ac2abbf1a7bb8b99e6a31ea31725fe2de482fb12b391caa0c48519aea864215"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]param([Parameter(Mandatory=$true)][string]$Root)
$ErrorActionPreference='Stop';$base=[IO.Path]::GetFullPath($Root)
$lock=Get-Content -LiteralPath (Join-Path $base 'arca\fiscal\arca-ws-contracts.lock.json') -Raw|ConvertFrom-Json -Depth 20
if($lock.toolchain.dotnetSdk-ne'10.0.400'-or$lock.toolchain.dotnetSvcutil-ne'8.0.0'){throw 'toolchain lock mismatch'}
foreach($name in @('wsaaHomologation','wsaaProduction','wsfeHomologation','wsfeProduction')){$c=$lock.contracts.$name;if(-not$c.url.StartsWith('https://',[StringComparison]::Ordinal)-or$c.bytes-lt1000-or$c.sha256-notmatch'^[0-9a-f]{64}$'){throw "invalid contract lock: $name"}}
foreach($name in @('wsaaHomologation','wsaaProduction','wsfeHomologation','wsfeProduction')){if($lock.contracts.$name.generatedSha256-notmatch'^[0-9a-f]{64}$'){throw "missing generated hash: $name"}}
$project=[IO.File]::ReadAllText((Join-Path $base 'arca\fiscal\Elite.Arca.Wsfe.Generated.csproj'));foreach($pin in @('10.0.652802','10.0.11','TreatWarningsAsErrors','NuGetAuditMode')){if($project-notmatch[regex]::Escape($pin)){throw "project contract missing: $pin"}}
$runner=[IO.File]::ReadAllText((Join-Path $base 'tools\generate-arca-wsfe-client.ps1'));foreach($gate in @('WSDL identity mismatch','DtdProcessing','production metadata requires','generated proxy identity mismatch','WsaaReference.cs','WsfeV1Reference.cs','--vulnerable','production_admitted=$false','DOTNET_ROLL_FORWARD')){if($runner-notmatch[regex]::Escape($gate)){throw "runner gate missing: $gate"}}
'ARCA_GENERATOR_CONTRACT_PASS contracts=4 generated_clients=2 toolchain=locked production_admitted=false'
````

## 6. Configuration surface

Inputs are the exact `dotnet` executable, cache, empty output, environment and offline mode. Production metadata needs `-AllowProductionMetadata`; this acknowledges only WSDL acquisition and never authorizes a production fiscal call.

## 7. Dependency bill

Free public dependencies: .NET SDK 10.0.400, `dotnet-svcutil` 8.0.0, System.ServiceModel.Http 10.0.652802, System.Security.Cryptography.Xml/Pkcs 10.0.11 and the selected ARCA WSDL. Actual ARCA credentials/services and business operations remain outside this pack.

## 8. Apply order

Materialize independently, run the structural regression, provide the official SDK and cache/network authority, then execute generation. Consume the generated assembly only from a separately reviewed WSAA/WSFE application adapter.

## 9. Verification

Require 6/6 reconstruction; four endpoint locks; structural regression; Homologation and Production-metadata generation of both clients; four exact generated hashes; zero build warnings/errors; and a clean current NuGet vulnerability result. Authentication and fiscal operation calls are prohibited.

## 10. Reconstruction evidence

Recorded in `reconstruction_evidence/MICROSOFT_ARCA_GENERATED_CLIENTS_2026-08-30_V124.md`; V123 remains the WSFE-only historical proof.
