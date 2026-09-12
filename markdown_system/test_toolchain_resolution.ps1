#requires -Version 7.0
# AUTHORED: tests the real runner declarations/functions without executing its audit.
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot
$tokens = $null
$errors = $null
$ast = [Management.Automation.Language.Parser]::ParseFile((Join-Path $root 'VERIFY_EXECUTABLE_LIBRARY.ps1'), [ref]$tokens, [ref]$errors)
if ($errors.Count) { throw 'Runner parse errors' }
$checks = 0
function Check([bool] $Condition, [string] $Message) {
  if (-not $Condition) { throw $Message }
  $script:checks++
}
$parameters = @($ast.ParamBlock.Parameters | ForEach-Object { $_.Name.VariablePath.UserPath })
foreach ($name in @('PythonExecutable','GoExecutable','DotNetExecutable','NodeExecutable','PnpmExecutable','DockerExecutable','PsqlExecutable')) {
  Check ($name -cin $parameters) "Missing explicit tool parameter: $name"
  Set-Variable -Name $name -Value ''
}
foreach ($name in @('Resolve-Tool','Resolve-PythonTool','Get-ToolVersion','Test-MinimumVersion','Get-ToolchainReport')) {
  $nodes = @($ast.FindAll({ param($n) $n -is [Management.Automation.Language.FunctionDefinitionAst] -and $n.Name -ceq $name }, $false))
  Check ($nodes.Count -eq 1) "Missing/ambiguous actual function: $name"
  . ([scriptblock]::Create($nodes[0].Extent.Text))
}
$hostPath = (Get-Process -Id $PID).Path
Check ((Resolve-Tool 'not-a-real-command-elite' $hostPath) -ceq $hostPath) 'Explicit path with spaces not preserved'
Check ($null -eq (Resolve-Tool 'pwsh' (Join-Path $root 'missing-tool-v286.exe'))) 'Invalid explicit tool fell back to PATH'
Check ($null -eq (Resolve-Tool 'pwsh' $root)) 'Directory accepted as executable'
Check ([bool](Resolve-Tool 'pwsh')) 'Default PATH resolution regressed'
Check ((Get-ToolVersion $hostPath @('--version')) -match '^PowerShell 7\.') 'Actual host invocation failed'
Check ((Get-ToolVersion $hostPath @('-NoProfile','-Command','exit 9')) -ceq '') 'Failed process accepted as version'
Check ((Get-ToolVersion '' @('--version')) -ceq '') 'Missing process accepted as version'
Check (Test-MinimumVersion 'go version go1.26.7 windows/amd64' ([version]'1.26.7')) 'Go version parsing failed'
Check (-not (Test-MinimumVersion 'go1.26.6' ([version]'1.26.7'))) 'Below minimum accepted'
Check (-not (Test-MinimumVersion 'not a version' ([version]'1.0'))) 'Malformed version accepted'
# Version invocations are stubbed only for routing tests: these files are NOT runtimes.
# Distinct paths detect accidental use of one parameter for another tool.
$mapping = [ordered]@{
  PythonExecutable=@('python-3.12+','AGENTS.md')
  GoExecutable=@('go-1.26.7','README.md')
  DotNetExecutable=@('dotnet-10+','CLAUDE.md')
  NodeExecutable=@('node-24+','VERIFY_LIBRARY.ps1')
  PnpmExecutable=@('pnpm-11.25.0','CREATE_PORTABLE_ARCHIVE.ps1')
  DockerExecutable=@('docker','INSTALL_AGENT_BRIDGE.ps1')
  PsqlExecutable=@('psql','VERIFY_EXECUTABLE_LIBRARY.ps1')
}
foreach ($entry in $mapping.GetEnumerator()) { Set-Variable -Name $entry.Key -Value (Join-Path $root $entry.Value[1]) }
function Get-ToolVersion([string] $Path, [string[]] $Arguments) { if ($Path) { '999.0.0' } else { '' } }
$report = @(Get-ToolchainReport)
Check ($report.Count -eq 8) 'Toolchain inventory changed unexpectedly'
foreach ($entry in $mapping.GetEnumerator()) {
  $row = @($report | Where-Object id -CEQ $entry.Value[0])
  Check ($row.Count -eq 1 -and $row[0].path -ceq (Join-Path $root $entry.Value[1])) "Incorrect explicit routing: $($entry.Key)"
}
$PythonExecutable = Join-Path $root 'missing-tool-v286.exe'
Check ($null -eq (Resolve-PythonTool)) 'Explicit Python selection silently fell back to python3'
Check (-not (@(Get-ToolchainReport | Where-Object id -CEQ 'python-3.12+')[0].available)) 'Invalid explicit Python reported available'
# All consumers, including audit-time probes, must use the same explicit paths.
$calls = @($ast.FindAll({ param($n) $n -is [Management.Automation.Language.CommandAst] -and $n.GetCommandName() -ceq 'Resolve-Tool' }, $true))
foreach ($tool in @(@('node','NodeExecutable'),@('pnpm','PnpmExecutable'),@('docker','DockerExecutable'),@('psql','PsqlExecutable'))) {
  $selected = @($calls | Where-Object { $_.CommandElements.Count -ge 2 -and $_.CommandElements[1].Extent.Text -ceq ("'" + $tool[0] + "'") })
  Check ($selected.Count -gt 0 -and @($selected | Where-Object { $_.CommandElements.Count -ne 3 -or $_.CommandElements[2].Extent.Text -cne ('$' + $tool[1]) }).Count -eq 0) "An audit consumer ignores $($tool[1])"
}
$probeRoot = Join-Path ([IO.Path]::GetTempPath()) ('elite-toolchain-test-' + [guid]::NewGuid().ToString('N'))
[void](New-Item -ItemType Directory -Path $probeRoot)
try {
  $probeEvidence = Join-Path $probeRoot 'probe.json'
  $missingPath = Join-Path $probeRoot 'absent.exe'
  $arguments = @('-NoProfile','-File',(Join-Path $root 'VERIFY_EXECUTABLE_LIBRARY.ps1'),'-Mode','Toolchain','-EvidencePath',$probeEvidence)
  foreach ($parameter in $mapping.Keys) { $arguments += @("-$parameter", $missingPath) }
  $output = @(& $hostPath @arguments 2>&1)
  Check ($LASTEXITCODE -eq 0) ("Informational toolchain probe failed: " + ($output -join "`n"))
  $receipt = Get-Content -LiteralPath $probeEvidence -Raw | ConvertFrom-Json
  Check ($receipt.mode -ceq 'Toolchain' -and $receipt.status -ceq 'TOOLS_MISSING') 'Incorrect probe status'
  Check ($receipt.library_audit_executed -eq $false -and $receipt.availability_is_admission -eq $false) 'Probe misrepresents audit/admission'
  Check ($receipt.steps.Count -eq 0) 'Probe ran audit steps'
  Check (@($receipt.toolchains | Where-Object { -not $_.available }).Count -eq 7) 'Explicit missing tools unexpectedly resolved'
  Check (@($output | Where-Object { $_ -match 'VERIFY_EXECUTABLE_LIBRARY_PASS|EXECUTABLE_PREFLIGHT_PASS' }).Count -eq 0) 'Probe claims an integrated PASS'
  $originalHash = (Get-FileHash -LiteralPath $probeEvidence -Algorithm SHA256).Hash
  $output = @(& $hostPath @arguments 2>&1)
  Check ($LASTEXITCODE -ne 0 -and ($output -join "`n") -match 'EvidencePath already exists') 'Existing evidence overwritten or wrong rejection'
  Check ((Get-FileHash -LiteralPath $probeEvidence -Algorithm SHA256).Hash -ceq $originalHash) 'Evidence bytes changed on rejection'
  $strictArguments = $arguments.Clone()
  $strictArguments[$strictArguments.IndexOf($probeEvidence)] = Join-Path $probeRoot 'strict.json'
  $output = @(& $hostPath @strictArguments -RequireAll 2>&1)
  Check ($LASTEXITCODE -ne 0 -and ($output -join "`n") -match 'missing toolchains:') 'RequireAll accepted missing tools'
  Check (Test-Path -LiteralPath (Join-Path $probeRoot 'strict.json') -PathType Leaf) 'RequireAll omitted negative evidence'
} finally {
  $resolved = [IO.Path]::GetFullPath($probeRoot)
  $expectedParent = [IO.Path]::GetFullPath([IO.Path]::GetTempPath()).TrimEnd([IO.Path]::DirectorySeparatorChar)
  if ((Split-Path -Parent $resolved) -cne $expectedParent -or (Split-Path -Leaf $resolved) -notlike 'elite-toolchain-test-*') { throw 'Unsafe test cleanup target' }
  Remove-Item -LiteralPath $resolved -Recurse -Force
}
Write-Output "PASS: toolchain resolution $checks checks; routing stubs are not runtime evidence"
exit 0
