#requires -Version 7.0
[CmdletBinding()]
param(
  [Parameter(Mandatory=$true)][string]$DotnetExecutable,
  [Parameter(Mandatory=$true)][string]$Root,
  [Parameter(Mandatory=$true)][string]$GeneratedClientRoot
)
$ErrorActionPreference='Stop'
$dotnet=[IO.Path]::GetFullPath($DotnetExecutable);$base=[IO.Path]::GetFullPath($Root)
if(-not(Test-Path -LiteralPath $dotnet -PathType Leaf)){throw 'dotnet executable missing'}
function Run([string[]]$Arguments){$lines=@(& $dotnet @Arguments 2>&1);if($LASTEXITCODE-ne0){throw "dotnet failed ($LASTEXITCODE): $($Arguments-join' ')`n$($lines-join[Environment]::NewLine)"};$lines}
$version=(@(Run @('--version'))[-1]).Trim();if($version-ne'10.0.400'){throw "dotnet SDK mismatch: $version"}
$project=Join-Path $base 'arca\credentials\tests\Elite.Arca.Credentials.Tests\Elite.Arca.Credentials.Tests.csproj'
$integrationProject=Join-Path $base 'arca\credentials\integration-tests\Elite.Arca.Wsaa.Transport.Tests\Elite.Arca.Wsaa.Transport.Tests.csproj'
$generated=[IO.Path]::GetFullPath($GeneratedClientRoot);if(-not(Test-Path -LiteralPath (Join-Path $generated 'Generated\WsaaReference.cs') -PathType Leaf)){throw 'generated WSAA client missing'}
$generatedProperty="-p:GeneratedClientRoot=$generated"
$oldTelemetry=$env:DOTNET_CLI_TELEMETRY_OPTOUT
try{
  $env:DOTNET_CLI_TELEMETRY_OPTOUT='1'
  Run @('restore',$project,'--locked-mode','--warnaserror')|Out-Null
  Run @('build',$project,'--configuration','Release','--no-restore','--warnaserror')|Out-Null
  $tests=Run @('run','--project',$project,'--configuration','Release','--no-build')
  if(($tests-join"`n")-notmatch'ARCA_WSAA_CREDENTIAL_TEST_PASS tests=9 secrets_logged=0'){throw 'test proof missing'}
  Run @('restore',$integrationProject,$generatedProperty,'--locked-mode','--warnaserror')|Out-Null
  Run @('build',$integrationProject,$generatedProperty,'--configuration','Release','--no-restore','--warnaserror')|Out-Null
  $integrationTests=Run @('run','--project',$integrationProject,$generatedProperty,'--configuration','Release','--no-build')
  if(($integrationTests-join"`n")-notmatch'ARCA_GENERATED_WSAA_TRANSPORT_TEST_PASS tests=3 manual_glue=0'){throw 'generated transport proof missing'}
  $audit=Run @('list',$project,'package','--vulnerable','--include-transitive')
  if(($audit-join"`n")-notmatch'no vulnerable packages'){throw "NuGet vulnerability gate did not return clean`n$($audit-join[Environment]::NewLine)"}
}finally{$env:DOTNET_CLI_TELEMETRY_OPTOUT=$oldTelemetry}
'ARCA_WSAA_CREDENTIAL_CORE_PASS tests=12 build_warnings=0 vulnerable_packages=0 secrets_logged=0 manual_glue=0 production_admitted=false'
