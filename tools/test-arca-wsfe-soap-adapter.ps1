#requires -Version 7.0
[CmdletBinding()]
param(
  [Parameter(Mandatory=$true)][string]$DotnetExecutable,
  [Parameter(Mandatory=$true)][string]$GeneratedRoot,
  [Parameter(Mandatory=$true)][string]$CredentialRoot,
  [Parameter(Mandatory=$true)][string]$WorkDirectory
)
$ErrorActionPreference='Stop'
$dotnet=[IO.Path]::GetFullPath($DotnetExecutable)
$generated=[IO.Path]::GetFullPath($GeneratedRoot)
$credentials=[IO.Path]::GetFullPath($CredentialRoot)
$work=[IO.Path]::GetFullPath($WorkDirectory)
if(-not(Test-Path $dotnet -PathType Leaf)){throw 'dotnet executable missing'}
if((@(& $dotnet --version 2>&1)[-1]).Trim()-ne'10.0.400'){throw 'dotnet SDK must equal 10.0.400'}
if(-not(Test-Path "$generated\Generated\WsfeV1Reference.cs" -PathType Leaf)-or-not(Test-Path "$generated\GENERATION_RECEIPT.json" -PathType Leaf)){throw 'generated homologation client/receipt missing'}
$receipt=Get-Content "$generated\GENERATION_RECEIPT.json" -Raw|ConvertFrom-Json
if($receipt.environment-ne'Homologation'-or$receipt.production_admitted-ne$false){throw 'only non-production homologation generation is admitted by this gate'}
$credentialProject=Join-Path $credentials 'arca\credentials\src\Elite.Arca.Credentials'
if(-not(Test-Path "$credentialProject\Elite.Arca.Credentials.csproj" -PathType Leaf)){throw 'credential core project missing'}
if(Test-Path $work){if((Get-ChildItem $work -Force).Count-ne0){throw 'work directory must be absent or empty'}}else{[IO.Directory]::CreateDirectory($work)|Out-Null}
$assets=[IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\arca\fiscal\bridge'))
$bridge=Join-Path $work 'arca\fiscal\bridge';$generatedOut=Join-Path $work 'arca\fiscal\generated';$credentialOut=Join-Path $work 'arca\credentials\src\Elite.Arca.Credentials'
$sourceOut=Join-Path $bridge 'src\Elite.Arca.Wsfe.Bridge';$testOut=Join-Path $bridge 'tests\Elite.Arca.Wsfe.Bridge.Tests';$parameterTestOut=Join-Path $bridge 'tests\Elite.Arca.Wsfe.Parameter.Tests'
foreach($directory in @($bridge,$sourceOut,$testOut,$parameterTestOut,$generatedOut,$credentialOut)){[IO.Directory]::CreateDirectory($directory)|Out-Null}
foreach($name in @('Elite.Arca.Wsfe.Bridge.csproj','packages.lock.json','Contracts.cs','GeneratedAdapters.cs','WsfeBridge.cs','WsfeParameterBridge.cs')){Copy-Item (Join-Path $assets "src\Elite.Arca.Wsfe.Bridge\$name") $sourceOut}
foreach($name in @('Elite.Arca.Wsfe.Bridge.Tests.csproj','packages.lock.json','Program.cs')){Copy-Item (Join-Path $assets "tests\Elite.Arca.Wsfe.Bridge.Tests\$name") $testOut}
foreach($name in @('Elite.Arca.Wsfe.Parameter.Tests.csproj','packages.lock.json','Program.cs')){Copy-Item (Join-Path $assets "tests\Elite.Arca.Wsfe.Parameter.Tests\$name") $parameterTestOut}
Copy-Item "$assets\README.md" $bridge
foreach($name in @('Elite.Arca.Wsfe.Generated.csproj','packages.lock.json','GENERATION_RECEIPT.json')){Copy-Item (Join-Path $generated $name) $generatedOut}
Copy-Item "$generated\Generated" $generatedOut -Recurse
foreach($name in @('Elite.Arca.Credentials.csproj','packages.lock.json','CertificateStoreLoader.cs','CmsRequestSigner.cs','LoginTicketRequestFactory.cs','WsaaContracts.cs','WsaaCredentialProvider.cs')){Copy-Item (Join-Path $credentialProject $name) $credentialOut}
$project=Join-Path $bridge 'tests\Elite.Arca.Wsfe.Bridge.Tests\Elite.Arca.Wsfe.Bridge.Tests.csproj';$parameterProject=Join-Path $bridge 'tests\Elite.Arca.Wsfe.Parameter.Tests\Elite.Arca.Wsfe.Parameter.Tests.csproj'
$nuget=Join-Path $work 'NuGet.Config';[IO.File]::WriteAllText($nuget,'<?xml version="1.0" encoding="utf-8"?><configuration><packageSources><clear /></packageSources></configuration>',[Text.UTF8Encoding]::new($false))
$old=$env:DOTNET_CLI_TELEMETRY_OPTOUT;try{$env:DOTNET_CLI_TELEMETRY_OPTOUT='1';foreach($candidate in @($project,$parameterProject)){& $dotnet restore $candidate --locked-mode --configfile $nuget;if($LASTEXITCODE){throw "locked offline restore failed: $candidate"};& $dotnet build $candidate -c Release --no-restore --warnaserror;if($LASTEXITCODE){throw "warning-free build failed: $candidate"}};$invoice=@(& $dotnet run --project $project -c Release --no-build 2>&1);if($LASTEXITCODE-ne0-or($invoice-join"`n")-notmatch'ARCA_WSFE_BRIDGE_TEST_PASS cases=15'){throw "invoice contract tests failed`n$($invoice-join[Environment]::NewLine)"};$parameters=@(& $dotnet run --project $parameterProject -c Release --no-build 2>&1);if($LASTEXITCODE-ne0-or($parameters-join"`n")-notmatch'ARCA_WSFE_PARAMETER_BRIDGE_PASS cases=15 secrets_returned=0 production_admitted=false'){throw "parameter contract tests failed`n$($parameters-join[Environment]::NewLine)"}}finally{$env:DOTNET_CLI_TELEMETRY_OPTOUT=$old}
"ARCA_WSFE_ADAPTER_GATE_PASS environment=Homologation invoice_cases=15 parameter_cases=15 production_admitted=false"
