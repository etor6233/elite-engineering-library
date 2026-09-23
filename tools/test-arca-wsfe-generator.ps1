#requires -Version 7.0
[CmdletBinding()]param([Parameter(Mandatory=$true)][string]$Root)
$ErrorActionPreference='Stop';$base=[IO.Path]::GetFullPath($Root)
$lock=Get-Content -LiteralPath (Join-Path $base 'arca\fiscal\arca-ws-contracts.lock.json') -Raw|ConvertFrom-Json -Depth 20
if($lock.toolchain.dotnetSdk-ne'10.0.400'-or$lock.toolchain.dotnetSvcutil-ne'8.0.0'){throw 'toolchain lock mismatch'}
foreach($name in @('wsaaHomologation','wsaaProduction','wsfeHomologation','wsfeProduction')){$c=$lock.contracts.$name;if(-not$c.url.StartsWith('https://',[StringComparison]::Ordinal)-or$c.bytes-lt1000-or$c.sha256-notmatch'^[0-9a-f]{64}$'){throw "invalid contract lock: $name"}}
foreach($name in @('wsaaHomologation','wsaaProduction','wsfeHomologation','wsfeProduction')){if($lock.contracts.$name.generatedSha256-notmatch'^[0-9a-f]{64}$'){throw "missing generated hash: $name"}}
$project=[IO.File]::ReadAllText((Join-Path $base 'arca\fiscal\Elite.Arca.Wsfe.Generated.csproj'));foreach($pin in @('10.0.652802','10.0.11','TreatWarningsAsErrors','NuGetAuditMode')){if($project-notmatch[regex]::Escape($pin)){throw "project contract missing: $pin"}}
$runner=[IO.File]::ReadAllText((Join-Path $base 'tools\generate-arca-wsfe-client.ps1'));foreach($gate in @('WSDL identity mismatch','DtdProcessing','supports only pinned offline homologation','generated proxy identity mismatch','WsaaReference.cs','WsfeV1Reference.cs','NOT_RUN_OFFLINE_RELEASE_SCA_REQUIRED','--locked-mode','tool receipt file drift','offline package pin mismatch','production_admitted=$false','DOTNET_ROLL_FORWARD')){if($runner-notmatch[regex]::Escape($gate)){throw "runner gate missing: $gate"}}
'ARCA_GENERATOR_STRUCTURAL_CONTRACT_PASS historical_contracts=4 local_generation_contracts=2 adapted_toolchain=locked actual_generation_and_SCA_separate=true production_admitted=false'
