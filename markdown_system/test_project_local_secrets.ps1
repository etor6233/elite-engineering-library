#requires -Version 7.0
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$tool = Join-Path $PSScriptRoot 'project_local_secrets.ps1'
$projectId = 'test-' + [guid]::NewGuid().ToString('N')
$root = Join-Path ([IO.Path]::GetTempPath()) $projectId
$null = New-Item -ItemType Directory -Path $root
$vault = Join-Path ([Environment]::GetFolderPath('LocalApplicationData')) "EliteEngineeringSecrets/$projectId/dev"
$checks = 0
function Assert($condition, $label) { if (-not $condition) { throw "TEST_FAILED $label" }; $script:checks++ }
function Expect-Failure([scriptblock]$Action) { $failed = $false; try { $null = & $Action } catch { $failed = $true }; Assert $failed 'negative accepted' }
function Read-Host { param($Prompt, [switch]$AsSecureString) return ConvertTo-SecureString 'synthetic-never-valid-provider-credential' -AsPlainText -Force }
try {
  Expect-Failure { & $tool -Mode Check -ProjectId $projectId -ProjectRoot $root -Names MISSING_SECRET }
  Expect-Failure { & $tool -Mode Generate -ProjectId $projectId -ProjectRoot $root -Names '../escape' }
  Expect-Failure { & $tool -Mode Generate -ProjectId $projectId -ProjectRoot $root -Names @('DUPLICATE','DUPLICATE') }
  Expect-Failure { & $tool -Mode Generate -ProjectId $projectId -ProjectRoot $root -Names NODE_OPTIONS }
  $result = & $tool -Mode Set -ProjectId $projectId -ProjectRoot $root -Names TEST_SECRET
  Assert ($result -eq 'TEST_SECRET STORED_DEV_DPAPI') 'set status'
  Assert (-not ([IO.File]::ReadAllText((Join-Path $vault 'TEST_SECRET.clixml')).Contains('synthetic-never-valid-provider-credential'))) 'plaintext on disk'
  $result = & $tool -Mode Check -ProjectId $projectId -ProjectRoot $root -Names TEST_SECRET
  Assert ($result -eq 'TEST_SECRET AVAILABLE_NOT_PROVIDER_VALIDATED') 'readability check'
  Expect-Failure { & $tool -Mode Set -ProjectId $projectId -ProjectRoot $root -Names TEST_SECRET }
  $null = & $tool -Mode Generate -ProjectId $projectId -ProjectRoot $root -Names AUTH_SESSION_SECRET
  $credential = Import-Clixml -LiteralPath (Join-Path $vault 'AUTH_SESSION_SECRET.clixml')
  Assert ([Convert]::FromBase64String($credential.GetNetworkCredential().Password).Length -eq 32) 'CSPRNG size'
  $credential.Password.Dispose()
  $public = Join-Path $root 'public.json'
  [IO.File]::WriteAllText($public, '{"ARCA_ENABLED":"false"}')
  $child = Join-Path $root 'child.ps1'
  [IO.File]::WriteAllText($child, 'if ($env:TEST_SECRET -ne "synthetic-never-valid-provider-credential" -or $env:ARCA_ENABLED -ne "false" -or $env:ELITE_MUST_NOT_INHERIT) { exit 1 }; Write-Output "CHILD_PASS"')
  $env:ELITE_MUST_NOT_INHERIT = 'synthetic'
  $result = & $tool -Mode Run -ProjectId $projectId -ProjectRoot $root -Names TEST_SECRET -Executable (Join-Path $PSHOME 'pwsh.exe') -Arguments @('-NoProfile','-File',$child) -PublicConfigFile $public
  Assert ($result -eq 'CHILD_EXIT_SUCCESS_NOT_PROVIDER_VALIDATED') 'child controlled injection'
  Assert ($null -eq [Environment]::GetEnvironmentVariable('TEST_SECRET','Process')) 'parent environment mutated'
  [IO.File]::WriteAllText($public, '{"TEST_SECRET":"do-not-accept"}')
  Expect-Failure { & $tool -Mode Run -ProjectId $projectId -ProjectRoot $root -Names TEST_SECRET -Executable (Join-Path $PSHOME 'pwsh.exe') -PublicConfigFile $public }
  [IO.File]::WriteAllText($public, '{"NODE_ENV":"production"}')
  Expect-Failure { & $tool -Mode Run -ProjectId $projectId -ProjectRoot $root -Names TEST_SECRET -Executable (Join-Path $PSHOME 'pwsh.exe') -PublicConfigFile $public }
  [IO.File]::WriteAllText((Join-Path $vault 'CORRUPT.clixml'), 'not a credential')
  Expect-Failure { & $tool -Mode Check -ProjectId $projectId -ProjectRoot $root -Names CORRUPT }
  'not secure' | Export-Clixml -LiteralPath (Join-Path $vault 'WRONG_TYPE.clixml')
  Expect-Failure { & $tool -Mode Check -ProjectId $projectId -ProjectRoot $root -Names WRONG_TYPE }
  $acl = Get-Acl -LiteralPath (Join-Path $vault 'TEST_SECRET.clixml')
  $acl.AddAccessRule([Security.AccessControl.FileSystemAccessRule]::new([Security.Principal.SecurityIdentifier]::new('S-1-1-0'),'Read','Allow'))
  [IO.FileSystemAclExtensions]::SetAccessControl([IO.FileInfo]::new((Join-Path $vault 'TEST_SECRET.clixml')), $acl)
  Expect-Failure { & $tool -Mode Check -ProjectId $projectId -ProjectRoot $root -Names TEST_SECRET }
  Write-Output "LOCAL_SECRETS_TEST_PASS checks=$checks synthetic_only=true"
} finally {
  [Environment]::SetEnvironmentVariable('ELITE_MUST_NOT_INHERIT',$null,'Process')
  # Only directories created by this invocation, with explicit prefix verification.
  $expected = [IO.Path]::GetFullPath((Join-Path ([Environment]::GetFolderPath('LocalApplicationData')) "EliteEngineeringSecrets/$projectId"))
  if ((Split-Path $vault -Parent) -eq $expected -and $projectId -match '^test-[a-f0-9]{32}$') {
    if (Test-Path -LiteralPath $expected) { Remove-Item -LiteralPath $expected -Recurse -Force }
    if (([IO.Path]::GetFullPath($root)) -eq (Join-Path ([IO.Path]::GetTempPath()) $projectId)) { Remove-Item -LiteralPath $root -Recurse -Force }
  }
}
