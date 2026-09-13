# Microsoft DevSkim Adapted SAST Gate

## 1. Metadata

```yaml
pack_id: "MICROSOFT-DEVSKIM-ADAPTED-SAST-GATE"
pack_version: "0.1.1"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Reconstruye Microsoft DevSkim desde un commit firmado y árbol canónico exactos, aplica una adaptación declarada para resolver SharpCompress 0.48.0, pasa 300 pruebas oficiales, SCA NuGet y fixtures Go/TypeScript; provee security linting local, no SAST interprocedural completo ni prueba de seguridad del producto."
stacks: ["Microsoft DevSkim", ".NET SDK 10.0.400", "SharpCompress 0.48.0", "PowerShell 7"]
compatible_with: ["GO-NATIVE-FUZZ-GATE 0.1.x", "SECURE-OPERATIONS-DELIVERY-CORE 0.4.x", "DEPENDENCY-LICENSE-EVIDENCE-CORE 0.1.x"]
incompatible_with: ["atribución VERBATIM posterior a la adaptación", "uso como único gate AppSec", "repositorio sin triage de findings", "red o toolchain no autorizados"]
license_expression: "LicenseRef-Workspace-Owner AND MIT"
upstream_sources: ["https://github.com/microsoft/DevSkim/tree/a452aa506f80928b611c4c7bfe306b8115821aae", "https://www.nuget.org/packages/SharpCompress/0.48.0"]
verified_at: "2026-09-05"
```

## 2. Applicability

Use como baseline local sin costo para security linting de Go y TypeScript cuando CodeQL privado no está licenciado. El agente debe materializar, ejecutar el verifier con red y SDK exacto y conservar el receipt. Los findings se trian en el proyecto; fuzzing, DAST, threat model, revisión y seguridad ofensiva siguen siendo gates separados.

No use el nombre Microsoft para ocultar la adaptación ni afirme que DevSkim realiza análisis interprocedural completo. Un advisory nuevo, cambio de .NET, nueva revisión Microsoft o cambio de claim reabre admisión y obliga a reconstruir.

## 3. Architecture contract

El runner rechaza destinos y receipts preexistentes; consulta la identidad y firma del commit mediante GitHub, descarga el archive fijado y valida un árbol canónico por archivo para tolerar únicamente reempaquetado de transporte. Verifica licencia/notices y proyectos originales antes de insertar dos referencias directas exactas a SharpCompress 0.48.0. Después usa sólo NuGet.org explícito, rechaza warnings NuGet, ejecuta build y 300 tests, exige cero vulnerabilidades, publica el CLI y ejecuta fixtures hostiles y limpio. La salida incluye licencias, provenance, lock y manifest SHA-256; el receipt es atómico y no sobrescribible.

## 4. Exact file manifest

```text
CREATE microsoft_devskim_adapted_sast/source-lock.json
CREATE microsoft_devskim_adapted_sast/build_and_verify.ps1
CREATE microsoft_devskim_adapted_sast/verify_pack.ps1
CREATE microsoft_devskim_adapted_sast/README.md
```

## 5. Materialization blocks

### FILE: `microsoft_devskim_adapted_sast/source-lock.json`
```yaml
block_id: "MICROSOFT-DEVSKIM-ADAPTED-SAST-GATE:source-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite lock derived from Microsoft GitHub API/source and exact NuGet package evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "c730eec041250a68bfe61aad33630233d473d9706050f48793e5d6ba00291b6d"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-microsoft-devskim-adapted-sast-lock/v1",
  "verified_at": "2026-09-05",
  "source": {
    "owner": "Microsoft",
    "repository": "microsoft/DevSkim",
    "commit": "a452aa506f80928b611c4c7bfe306b8115821aae",
    "tree": "f6be283e73f106afdb1345c4123539a7a2f9a136",
    "commit_signature_verified": true,
    "archive_url": "https://github.com/microsoft/DevSkim/archive/a452aa506f80928b611c4c7bfe306b8115821aae.zip",
    "commit_api_url": "https://api.github.com/repos/microsoft/DevSkim/commits/a452aa506f80928b611c4c7bfe306b8115821aae",
    "archive_transport_observations": [
      {"bytes": 935159, "sha256": "2e75a048dcf3e3db97e4551207b0a912d3bc504e2d2689d21d82e48dfcd89410"},
      {"bytes": 949833, "sha256": "c3b77fb88ff293e4d9f94ccb312a1bd65e65fa18847169d1217ece8fa04ec922"}
    ],
    "archive_file_entries": 268,
    "archive_root_observations": ["microsoft-DevSkim-a452aa5", "DevSkim-a452aa506f80928b611c4c7bfe306b8115821aae"],
    "canonical_file_tree_sha256": "4914ee46a553ee3c80b7e4849cbbcc349f1deab99951373229f177c88ad43ecb",
    "license_expression": "MIT",
    "license_path": "LICENSE.txt",
    "license_sha256": "ad0cf28f3381ca9bb0bf101d127402d44c17bfa0991e1a00bff7ae6679e9dada",
    "notice_path": "NOTICE.txt",
    "notice_sha256": "78535e301bab95ed0c86398020fbcc1654a4be55db2f48e4d1ba1eb8368b1947",
    "library_project_path": "DevSkim-DotNet/Microsoft.DevSkim/Microsoft.DevSkim.csproj",
    "library_project_sha256": "f2f557bfba82cb4dc605915136fb717d891d0ae4ecbb9d0a87740dbc4b55ad99",
    "cli_project_path": "DevSkim-DotNet/Microsoft.DevSkim.CLI/Microsoft.DevSkim.CLI.csproj",
    "cli_project_sha256": "7a44fb2969c7b42d0278f9cb33587e9401da889bb67c5288134b3477154c3a55",
    "test_project_path": "DevSkim-DotNet/Microsoft.DevSkim.Tests/Microsoft.DevSkim.Tests.csproj"
  },
  "adaptation": {
    "classification": "ADAPTED",
    "description": "Add a direct SharpCompress 0.48.0 reference to the Microsoft.DevSkim library and CLI projects so NuGet direct-dependency-wins removes the affected transitive 0.40.0 runtime.",
    "sharpcompress_version": "0.48.0",
    "fixed_minimum_for_ghsa_6c8g_7p36_r338": "0.48.0",
    "library_project_patched_sha256": "a5439c44f6d5dafae269185dfd090816c43f37c08f81d9283432f20b913b6016",
    "cli_project_patched_sha256": "9ffdc50a3fb44861d5cd6987ce805f3f4ef250b610ffe5339c47c3892e71a085",
    "package_url": "https://api.nuget.org/v3-flatcontainer/sharpcompress/0.48.0/sharpcompress.0.48.0.nupkg",
    "package_bytes": 9262981,
    "package_sha256": "d8c5da8a76d325eb81c1103a78953e025513f22ade36b5b11d8342324146f0b7",
    "package_repository_commit": "6e59c7d7bbf8c19a8a92c3c382599906684bb93d",
    "license_expression": "MIT",
    "license_url": "https://raw.githubusercontent.com/adamhathcock/sharpcompress/6e59c7d7bbf8c19a8a92c3c382599906684bb93d/LICENSE.txt",
    "license_bytes": 1081,
    "license_sha256": "b7ca2b6174cee11afe41d78b48527e3d7659a4435c25019a4a5e072299a2f8ed"
  },
  "verified_lane": {
    "os": "Windows x64",
    "dotnet_sdk": "10.0.400",
    "target_framework": "net10.0",
    "official_tests_passed": 300,
    "nuget_vulnerability_findings": 0,
    "typescript_fixture_rule": "DS189424",
    "go_fixture_rule": "DS112852",
    "clean_fixture_findings": 0
  }
}
````

### FILE: `microsoft_devskim_adapted_sast/build_and_verify.ps1`
```yaml
block_id: "MICROSOFT-DEVSKIM-ADAPTED-SAST-GATE:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite fail-closed reconstruction, adaptation, test, SCA and publication runner"
license: "LicenseRef-Workspace-Owner"
sha256: "46b531bcd4abb791092bdea7605d96dc4e13f2005cfff10630666cf7d16b74d7"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)] [string] $Destination,
  [Parameter(Mandatory = $true)] [string] $ReceiptPath,
  [Parameter(Mandatory = $true)] [string] $DotNetExecutable,
  [switch] $AllowNetwork
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$utf8 = [Text.UTF8Encoding]::new($false)
$destinationPath = [IO.Path]::GetFullPath($Destination)
$receipt = [IO.Path]::GetFullPath($ReceiptPath)
$dotnet = [IO.Path]::GetFullPath($DotNetExecutable)
if (Test-Path -LiteralPath $destinationPath) { throw "Destination already exists: $destinationPath" }
if (Test-Path -LiteralPath $receipt) { throw "Receipt already exists: $receipt" }
if (-not $AllowNetwork) { throw 'AllowNetwork is required for exact GitHub and NuGet acquisition' }
if (-not (Test-Path -LiteralPath $dotnet -PathType Leaf)) { throw "DotNetExecutable missing: $dotnet" }
$destinationParent = Split-Path -Parent $destinationPath
$receiptParent = Split-Path -Parent $receipt
if (-not $destinationParent -or -not $receiptParent) { throw 'Destination and ReceiptPath require parent directories' }
if (-not (Test-Path -LiteralPath $destinationParent -PathType Container)) { [void](New-Item -ItemType Directory -Path $destinationParent) }
if (-not (Test-Path -LiteralPath $receiptParent -PathType Container)) { [void](New-Item -ItemType Directory -Path $receiptParent) }

$lock = Get-Content -LiteralPath (Join-Path $PSScriptRoot 'source-lock.json') -Raw | ConvertFrom-Json -Depth 30
$sdkRoot = Split-Path -Parent $dotnet
$previousRoot = $env:DOTNET_ROOT
$previousRootX64 = $env:DOTNET_ROOT_X64
$env:DOTNET_ROOT = $sdkRoot
$env:DOTNET_ROOT_X64 = $sdkRoot
$stage = Join-Path ([IO.Path]::GetTempPath()) ('elite-devskim-build-' + [Guid]::NewGuid().ToString('N'))
$candidate = Join-Path $destinationParent ('.elite-devskim-candidate-' + [Guid]::NewGuid().ToString('N'))

function Get-Sha([string] $Path) { (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant() }
function Invoke-DotNet([string[]] $Arguments) {
  $lines = @(& $dotnet @Arguments 2>&1 | ForEach-Object { [string]$_ })
  if ($LASTEXITCODE -ne 0) { throw "dotnet failed ($LASTEXITCODE): $($Arguments -join ' ')`n$(($lines | Select-Object -Last 20) -join "`n")" }
  $lines
}
function Assert-NoNuGetWarning([string[]] $Lines, [string] $StageName) {
  $warnings = @($Lines | Where-Object { $_ -match '(?i)\bwarning\s+NU\d+' })
  if ($warnings.Count -ne 0) { throw "$StageName emitted NuGet warnings: $($warnings -join ' | ')" }
}
function Get-VulnerabilityCount([object] $Report) {
  $count = 0
  foreach ($project in @($Report.projects)) {
    if ($null -eq $project.PSObject.Properties['frameworks']) { continue }
    foreach ($framework in @($project.frameworks)) {
      $top = if ($null -ne $framework.PSObject.Properties['topLevelPackages']) { @($framework.topLevelPackages) } else { @() }
      $transitive = if ($null -ne $framework.PSObject.Properties['transitivePackages']) { @($framework.transitivePackages) } else { @() }
      foreach ($package in $top + $transitive) {
        if ($null -ne $package -and $null -ne $package.PSObject.Properties['vulnerabilities']) { $count += @($package.vulnerabilities).Count }
      }
    }
  }
  $count
}
function Get-PackageVersions([object] $Report, [string] $Name) {
  @(
    foreach ($project in @($Report.projects)) {
      if ($null -eq $project.PSObject.Properties['frameworks']) { continue }
      foreach ($framework in @($project.frameworks)) {
        $top = if ($null -ne $framework.PSObject.Properties['topLevelPackages']) { @($framework.topLevelPackages) } else { @() }
        $transitive = if ($null -ne $framework.PSObject.Properties['transitivePackages']) { @($framework.transitivePackages) } else { @() }
        foreach ($package in $top + $transitive) {
          if ($null -ne $package -and [string]$package.id -ceq $Name) { [string]$package.resolvedVersion }
        }
      }
    }
  )
}

try {
  $sdkVersion = ((Invoke-DotNet @('--version')) | Out-String).Trim()
  if ($sdkVersion -cne [string]$lock.verified_lane.dotnet_sdk) { throw "SDK mismatch: expected=$($lock.verified_lane.dotnet_sdk) actual=$sdkVersion" }
  [void](New-Item -ItemType Directory -Path $stage)
  $archive = Join-Path $stage 'devskim-source.zip'
  $headers = @{'User-Agent'='EliteEngineeringLibrary-DevSkimGate';'Accept'='application/vnd.github+json'}
  $commitRecord = Invoke-RestMethod -Uri $lock.source.commit_api_url -Headers $headers
  if ($commitRecord.sha -cne $lock.source.commit -or $commitRecord.commit.tree.sha -cne $lock.source.tree -or $commitRecord.commit.verification.verified -ne $true) { throw 'GitHub commit identity or signature mismatch' }
  Invoke-WebRequest -Uri $lock.source.archive_url -OutFile $archive
  $archiveBytes = (Get-Item -LiteralPath $archive).Length
  $archiveTransportSha = Get-Sha $archive

  Add-Type -AssemblyName System.IO.Compression.FileSystem
  $expanded = Join-Path $stage 'expanded'
  $zip = [IO.Compression.ZipFile]::OpenRead($archive)
  try {
    $files = @($zip.Entries | Where-Object { -not [string]::IsNullOrEmpty($_.Name) })
    if ($files.Count -ne [int]$lock.source.archive_file_entries) { throw "DevSkim archive file inventory drifted: $($files.Count)" }
    $roots = @($zip.Entries | ForEach-Object { $_.FullName.Replace('\','/').Split('/')[0] } | Where-Object { $_ } | Sort-Object -Unique)
    if ($roots.Count -ne 1 -or $roots[0] -match '(^\.|[\\/:])') { throw "DevSkim archive must have one safe root: $($roots -join ',')" }
    $archiveRoot = $roots[0]
    $canonicalRows = [Collections.Generic.List[string]]::new()
    foreach ($entry in $zip.Entries) {
      $name = $entry.FullName.Replace('\','/')
      if ([IO.Path]::IsPathRooted($name) -or $name.Split('/') -contains '..' -or -not $name.StartsWith($archiveRoot + '/', [StringComparison]::Ordinal)) { throw "Unsafe or unexpected archive entry: $name" }
      if (-not [string]::IsNullOrEmpty($entry.Name)) {
        $relative = $name.Substring($archiveRoot.Length + 1)
        $stream = $entry.Open()
        try { $entrySha = [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData($stream)).ToLowerInvariant() } finally { $stream.Dispose() }
        $canonicalRows.Add("$relative|$($entry.Length)|$entrySha`n")
      }
    }
    $canonicalText = (@($canonicalRows | Sort-Object) -join '')
    $canonicalHash = [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData($utf8.GetBytes($canonicalText))).ToLowerInvariant()
    if ($canonicalHash -cne $lock.source.canonical_file_tree_sha256) { throw "DevSkim canonical file tree drifted: $canonicalHash" }
  } finally { $zip.Dispose() }
  [IO.Compression.ZipFile]::ExtractToDirectory($archive, $expanded)
  $source = Join-Path $expanded $archiveRoot
  $license = Join-Path $source $lock.source.license_path
  $notice = Join-Path $source $lock.source.notice_path
  $libraryProject = Join-Path $source $lock.source.library_project_path
  $cliProject = Join-Path $source $lock.source.cli_project_path
  $testProject = Join-Path $source $lock.source.test_project_path
  if ((Get-Sha $license) -cne $lock.source.license_sha256 -or (Get-Sha $notice) -cne $lock.source.notice_sha256) { throw 'DevSkim license or notice drifted' }
  if ((Get-Sha $libraryProject) -cne $lock.source.library_project_sha256 -or (Get-Sha $cliProject) -cne $lock.source.cli_project_sha256) { throw 'DevSkim project source drifted' }

  $oldLibrary = '<PackageReference Include="Microsoft.CST.ApplicationInspector.RulesEngine" Version="1.9.50" />'
  $newLibrary = $oldLibrary + "`r`n      " + '<PackageReference Include="SharpCompress" Version="0.48.0" />'
  $oldCli = '<PackageReference Include="Microsoft.CST.ApplicationInspector.Logging" Version="1.9.50" />'
  $newCli = $oldCli + "`r`n      " + '<PackageReference Include="SharpCompress" Version="0.48.0" />'
  $libraryText = [IO.File]::ReadAllText($libraryProject)
  $cliText = [IO.File]::ReadAllText($cliProject)
  if (($libraryText.Split($oldLibrary).Count - 1) -ne 1 -or ($cliText.Split($oldCli).Count - 1) -ne 1) { throw 'DevSkim adaptation anchors drifted' }
  [IO.File]::WriteAllText($libraryProject, $libraryText.Replace($oldLibrary, $newLibrary), $utf8)
  [IO.File]::WriteAllText($cliProject, $cliText.Replace($oldCli, $newCli), $utf8)
  if ((Get-Sha $libraryProject) -cne $lock.adaptation.library_project_patched_sha256 -or (Get-Sha $cliProject) -cne $lock.adaptation.cli_project_patched_sha256) { throw 'Adapted project hashes mismatch' }

  $sharpPackage = Join-Path $stage 'sharpcompress.0.48.0.nupkg'
  Invoke-WebRequest -Uri $lock.adaptation.package_url -OutFile $sharpPackage
  if ((Get-Item -LiteralPath $sharpPackage).Length -ne [int64]$lock.adaptation.package_bytes -or (Get-Sha $sharpPackage) -cne $lock.adaptation.package_sha256) { throw 'SharpCompress package identity mismatch' }
  $sharpLicense = Join-Path $stage 'SharpCompress-LICENSE.txt'
  Invoke-WebRequest -Uri $lock.adaptation.license_url -OutFile $sharpLicense
  if ((Get-Item -LiteralPath $sharpLicense).Length -ne [int64]$lock.adaptation.license_bytes -or (Get-Sha $sharpLicense) -cne $lock.adaptation.license_sha256) { throw 'SharpCompress license identity mismatch' }

  $nugetConfig = Join-Path $stage 'NuGet.Public.config'
  [IO.File]::WriteAllText($nugetConfig, "<?xml version=`"1.0`" encoding=`"utf-8`"?>`n<configuration>`n  <packageSources>`n    <clear />`n    <add key=`"nuget.org`" value=`"https://api.nuget.org/v3/index.json`" protocolVersion=`"3`" />`n  </packageSources>`n</configuration>`n", $utf8)
  $restore = @(Invoke-DotNet @('restore',$testProject,'--configfile',$nugetConfig,'--property:TargetFramework=net10.0','--force','--no-cache'))
  Assert-NoNuGetWarning $restore 'restore'
  $build = @(Invoke-DotNet @('build',$testProject,'--framework','net10.0','--no-restore','--configuration','Release'))
  Assert-NoNuGetWarning $build 'build'
  $tests = @(Invoke-DotNet @('test',$testProject,'--framework','net10.0','--no-restore','--no-build','--configuration','Release'))
  $testText = $tests -join "`n"
  if ($testText -notmatch 'Failed:\s*0,\s*Passed:\s*300,\s*Skipped:\s*0,\s*Total:\s*300') { throw 'Official Microsoft test result count drifted' }

  $projects = @($libraryProject,$cliProject,$testProject)
  $vulnerabilityCount = 0
  foreach ($project in $projects) {
    $vulnerableLines = @(Invoke-DotNet @('list',$project,'package','--vulnerable','--include-transitive','--framework','net10.0','--configfile',$nugetConfig,'--no-restore','--format','json'))
    Assert-NoNuGetWarning $vulnerableLines 'vulnerability audit'
    $vulnerabilityCount += Get-VulnerabilityCount (($vulnerableLines -join "`n") | ConvertFrom-Json -Depth 50)
  }
  if ($vulnerabilityCount -ne 0) { throw "NuGet vulnerability findings: $vulnerabilityCount" }
  $inventoryLines = @(Invoke-DotNet @('list',$cliProject,'package','--include-transitive','--framework','net10.0','--no-restore','--format','json'))
  $inventory = ($inventoryLines -join "`n") | ConvertFrom-Json -Depth 50
  $sharpVersions = @(Get-PackageVersions $inventory 'SharpCompress' | Sort-Object -Unique)
  if ($sharpVersions.Count -ne 1 -or $sharpVersions[0] -cne $lock.adaptation.sharpcompress_version) { throw "SharpCompress runtime resolution drifted: $($sharpVersions -join ',')" }

  $publish = Join-Path $stage 'publish'
  [void](Invoke-DotNet @('publish',$cliProject,'--configuration','Release','--framework','net10.0','--no-restore','--output',$publish,'-p:UseAppHost=false'))
  $cliDll = Join-Path $publish 'devskim.dll'
  if (-not (Test-Path -LiteralPath $cliDll -PathType Leaf)) { throw 'Published DevSkim CLI missing' }
  $fixtures = Join-Path $stage 'fixtures'
  [void](New-Item -ItemType Directory -Path $fixtures)
  [IO.File]::WriteAllText((Join-Path $fixtures 'unsafe.ts'), "const value = eval(userControlledExpression);`n", $utf8)
  [IO.File]::WriteAllText((Join-Path $fixtures 'unsafe.go'), "package fixture`n`nimport `"crypto/tls`"`n`nvar minimum = tls.VersionTLS10`n", $utf8)
  [IO.File]::WriteAllText((Join-Path $fixtures 'safe.ts'), "export function normalize(value: string): string {`n  return value.trim();`n}`n", $utf8)
  $expected = [ordered]@{'unsafe.ts'='DS189424';'unsafe.go'='DS112852';'safe.ts'=''}
  foreach ($name in $expected.Keys) {
    $sarifPath = Join-Path $fixtures ($name + '.sarif')
    [void](Invoke-DotNet @($cliDll,'analyze','--source-code',(Join-Path $fixtures $name),'--output-file',$sarifPath,'--file-format','sarif','--disable-console','--skip-excerpts'))
    $sarif = Get-Content -LiteralPath $sarifPath -Raw | ConvertFrom-Json -Depth 100
    $ids = @($sarif.runs[0].results | ForEach-Object { [string]$_.ruleId })
    if ($expected[$name]) {
      if ($ids.Count -ne 1 -or $ids[0] -cne $expected[$name]) { throw "Fixture result drifted for ${name}: $($ids -join ',')" }
    } elseif ($ids.Count -ne 0) { throw "Clean fixture produced findings: $($ids -join ',')" }
  }

  [void](New-Item -ItemType Directory -Path $candidate)
  foreach ($item in Get-ChildItem -LiteralPath $publish -Force) { Copy-Item -LiteralPath $item.FullName -Destination $candidate -Recurse }
  Copy-Item -LiteralPath $license -Destination (Join-Path $candidate 'MICROSOFT_DEVSKIM_LICENSE.txt')
  Copy-Item -LiteralPath $notice -Destination (Join-Path $candidate 'MICROSOFT_DEVSKIM_NOTICE.txt')
  Copy-Item -LiteralPath $sharpLicense -Destination (Join-Path $candidate 'SHARPCOMPRESS_LICENSE.txt')
  Copy-Item -LiteralPath (Join-Path $PSScriptRoot 'source-lock.json') -Destination (Join-Path $candidate 'SOURCE_LOCK.json')
  $provenance = [ordered]@{schema='elite-adapted-devskim-provenance/v1';classification='ADAPTED';microsoft_commit=$lock.source.commit;adaptation=$lock.adaptation.description;sharpcompress_version=$sharpVersions[0];official_tests_passed=300;nuget_vulnerability_findings=0;non_claims=@('not unmodified Microsoft code','not comprehensive interprocedural SAST','not proof that analyzed code is secure','not a replacement for DAST, fuzzing, review or production security testing')}
  [IO.File]::WriteAllText((Join-Path $candidate 'PROVENANCE.json'), (($provenance | ConvertTo-Json -Depth 10) + "`n"), $utf8)
  $manifestEntries = @(
    Get-ChildItem -LiteralPath $candidate -Recurse -File | Sort-Object FullName | ForEach-Object {
      [ordered]@{path=[IO.Path]::GetRelativePath($candidate,$_.FullName).Replace('\','/');bytes=$_.Length;sha256=(Get-Sha $_.FullName)}
    }
  )
  $manifest = [ordered]@{schema='elite-adapted-devskim-manifest/v1';files=$manifestEntries}
  [IO.File]::WriteAllText((Join-Path $candidate 'MANIFEST.sha256.json'), (($manifest | ConvertTo-Json -Depth 10) + "`n"), $utf8)
  [IO.Directory]::Move($candidate, $destinationPath)
  $record = [ordered]@{schema='elite-adapted-devskim-build-receipt/v1';status='PASS';verified_at=[DateTimeOffset]::UtcNow.ToString('O');destination=$destinationPath;source_commit=$lock.source.commit;source_tree=$lock.source.tree;canonical_file_tree_sha256=$canonicalHash;archive_transport_root=$archiveRoot;archive_transport_bytes=$archiveBytes;archive_transport_sha256=$archiveTransportSha;classification='ADAPTED';dotnet_sdk=$sdkVersion;target_framework='net10.0';official_tests_passed=300;nuget_vulnerability_findings=0;sharpcompress_version=$sharpVersions[0];typescript_fixture_rule='DS189424';go_fixture_rule='DS112852';clean_fixture_findings=0;manifest_sha256=(Get-Sha (Join-Path $destinationPath 'MANIFEST.sha256.json'));limitations=@('security linting is not comprehensive interprocedural SAST','findings require project triage','project DAST, fuzzing, threat model and offensive tests remain required')}
  $receiptTemp = Join-Path $receiptParent ('.elite-devskim-receipt-' + [Guid]::NewGuid().ToString('N') + '.tmp')
  [IO.File]::WriteAllText($receiptTemp, (($record | ConvertTo-Json -Depth 10) + "`n"), $utf8)
  [IO.File]::Move($receiptTemp, $receipt)
  Write-Output "MICROSOFT_DEVSKIM_ADAPTED_SAST_PASS tests=300 vulnerabilities=0 typescript=DS189424 go=DS112852 clean=0"
} finally {
  $env:DOTNET_ROOT = $previousRoot
  $env:DOTNET_ROOT_X64 = $previousRootX64
  if (Test-Path -LiteralPath $candidate) { [IO.Directory]::Delete($candidate, $true) }
  if (Test-Path -LiteralPath $stage) { [IO.Directory]::Delete($stage, $true) }
}
````

### FILE: `microsoft_devskim_adapted_sast/verify_pack.ps1`
```yaml
block_id: "MICROSOFT-DEVSKIM-ADAPTED-SAST-GATE:verifier:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite static, negative and optional live reconstruction verifier"
license: "LicenseRef-Workspace-Owner"
sha256: "e62d3646dcf9a6fed4af336987f549ce25ae583e3ba39b728205b07dcc2a9ea2"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [string] $DotNetExecutable = '',
  [switch] $AllowNetwork
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$lock = Get-Content -LiteralPath (Join-Path $PSScriptRoot 'source-lock.json') -Raw | ConvertFrom-Json -Depth 30
if ($lock.source.owner -cne 'Microsoft' -or $lock.source.commit -notmatch '^[0-9a-f]{40}$' -or $lock.source.tree -notmatch '^[0-9a-f]{40}$' -or $lock.source.archive_url -match '/(main|master|latest)(/|\.|$)' -or $lock.source.canonical_file_tree_sha256 -notmatch '^[0-9a-f]{64}$') { throw 'Source identity is not exact Microsoft source' }
if ($lock.source.license_expression -cne 'MIT' -or $lock.adaptation.classification -cne 'ADAPTED') { throw 'License or adaptation classification drifted' }
if ($lock.adaptation.sharpcompress_version -cne '0.48.0' -or $lock.adaptation.license_expression -cne 'MIT') { throw 'SharpCompress lock drifted' }
if ($lock.verified_lane.official_tests_passed -ne 300 -or $lock.verified_lane.nuget_vulnerability_findings -ne 0 -or $lock.verified_lane.typescript_fixture_rule -cne 'DS189424' -or $lock.verified_lane.go_fixture_rule -cne 'DS112852') { throw 'Verified lane drifted' }
foreach ($name in @('build_and_verify.ps1','verify_pack.ps1')) {
  $tokens=$null; $errors=$null
  [void][Management.Automation.Language.Parser]::ParseFile((Join-Path $PSScriptRoot $name),[ref]$tokens,[ref]$errors)
  if (@($errors).Count -ne 0) { throw "PowerShell syntax failed: $name" }
}
$runner = Join-Path $PSScriptRoot 'build_and_verify.ps1'
$negativeRoot = Join-Path ([IO.Path]::GetTempPath()) ('elite-devskim-negative-' + [Guid]::NewGuid().ToString('N'))
try {
  [void](New-Item -ItemType Directory -Path $negativeRoot)
  $existingDestination = Join-Path $negativeRoot 'existing-destination'
  [void](New-Item -ItemType Directory -Path $existingDestination)
  $failed = $false
  try { & $runner -Destination $existingDestination -ReceiptPath (Join-Path $negativeRoot 'receipt.json') -DotNetExecutable 'C:\missing\dotnet.exe' -AllowNetwork } catch { $failed = $_.Exception.Message -like 'Destination already exists:*' }
  if (-not $failed) { throw 'Existing-destination negative did not fail closed' }
  $existingReceipt = Join-Path $negativeRoot 'existing-receipt.json'
  [IO.File]::WriteAllText($existingReceipt, '{}')
  $failed = $false
  try { & $runner -Destination (Join-Path $negativeRoot 'new-destination') -ReceiptPath $existingReceipt -DotNetExecutable 'C:\missing\dotnet.exe' -AllowNetwork } catch { $failed = $_.Exception.Message -like 'Receipt already exists:*' }
  if (-not $failed) { throw 'Existing-receipt negative did not fail closed' }
  if ($AllowNetwork) {
    if (-not $DotNetExecutable) { throw 'DotNetExecutable is required with AllowNetwork' }
    & $runner -Destination (Join-Path $negativeRoot 'runtime') -ReceiptPath (Join-Path $negativeRoot 'runtime-receipt.json') -DotNetExecutable $DotNetExecutable -AllowNetwork
    if ($LASTEXITCODE -ne 0) { throw "Runtime verification exited $LASTEXITCODE" }
  }
  Write-Output "MICROSOFT_DEVSKIM_ADAPTED_PACK_PASS static=1 negatives=2 runtime=$([int][bool]$AllowNetwork)"
} finally {
  if (Test-Path -LiteralPath $negativeRoot) { [IO.Directory]::Delete($negativeRoot,$true) }
}
````

### FILE: `microsoft_devskim_adapted_sast/README.md`
```yaml
block_id: "MICROSOFT-DEVSKIM-ADAPTED-SAST-GATE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite operating instructions and non-claims"
license: "LicenseRef-Workspace-Owner"
sha256: "1a9791b25bb7db555c84b2e874da8b8e3dac1dd4eefb75a8ce83362de463a84b"
variables: []
secrets_allowed: false
```
````markdown
# Microsoft DevSkim adapted SAST baseline

This pack reconstructs Microsoft DevSkim from the exact signed commit in `source-lock.json`, applies one transparent dependency adaptation, runs the 300 official tests, audits all NuGet graphs, publishes the local CLI and proves one TypeScript finding, one Go finding and one clean fixture.

Run from a newly materialized directory:

```powershell
pwsh -NoProfile -File ./verify_pack.ps1 -AllowNetwork -DotNetExecutable C:/absolute/path/to/dotnet.exe
```

Or build a retained tool and immutable receipt:

```powershell
pwsh -NoProfile -File ./build_and_verify.ps1 -AllowNetwork -DotNetExecutable C:/absolute/path/to/dotnet.exe -Destination C:/absent/devskim-runtime -ReceiptPath C:/absent/devskim-receipt.json
```

The published runtime is `ADAPTED`, not verbatim Microsoft code. It is a no-cost local security-lint baseline for Go and TypeScript. It is not comprehensive interprocedural SAST, does not prove an application secure and does not replace CodeQL entitlement, threat modeling, review, fuzzing, DAST or offensive testing. Re-run the dependency/freshness admission when Microsoft publishes a new release/commit, NuGet advisories change, .NET 10.0.400 leaves support or the project language/claim changes.
````

## 6. Configuration surface

- `Destination`: directorio nuevo para el runtime publicado.
- `ReceiptPath`: archivo nuevo para evidencia atómica.
- `DotNetExecutable`: SDK .NET 10.0.400 exacto observado.
- `AllowNetwork`: autorización explícita para GitHub y NuGet.org.
- El proyecto define scope, exclusiones, severidades, baseline, suppressions, owner y SLA de triage fuera del pack.

## 7. Dependency bill

- Microsoft DevSkim commit `a452aa506f80928b611c4c7bfe306b8115821aae`, MIT, source adquirido y adaptado durante reconstrucción.
- SharpCompress 0.48.0, MIT, nupkg/licencia/hash fijados.
- .NET SDK 10.0.400 como toolchain de construcción; no se redistribuye por el pack.
- PowerShell 7 y NuGet.org como runner/registry.

## 8. Apply order

1. Completar el perfil de seguridad del proyecto y autorizar red/toolchain.
2. Materializar los cuatro archivos en destino vacío.
3. Ejecutar el verifier completo y conservar receipt/runtime.
4. Escanear únicamente el scope aprobado y retener SARIF sin filtrar secretos/código indebidamente.
5. Triar findings, corregir, añadir regresiones y repetir hasta el policy PASS.
6. Ejecutar fuzzing, DAST, threat model, revisión y ofensiva como gates separados.

## 9. Verification

```powershell
pwsh -NoProfile -File ./materialize_markdown_pack.ps1 `
  -PackFile ./implementation_packs/MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md `
  -Destination C:/absent/devskim-pack

pwsh -NoProfile -File C:/absent/devskim-pack/microsoft_devskim_adapted_sast/verify_pack.ps1 `
  -AllowNetwork `
  -DotNetExecutable C:/absolute/dotnet.exe
```

PASS exige materialización 4/4, sintaxis PowerShell, dos negativos fail-closed, identidad/firma/tree del commit, adaptación/hash exactos, 300/300 tests, cero findings NuGet y resultados exactos de fixtures. Cualquier divergencia bloquea la incorporación.

## 10. Reconstruction evidence

Evidencia canónica: `reconstruction_evidence/MICROSOFT_DEVSKIM_ADAPTED_SAST_2026-09-05_V252.md`. La reconstrucción exacta pasó en Windows x64 con .NET SDK 10.0.400: 300 tests Microsoft, cero vulnerabilidades NuGet y fixtures Go/TypeScript/clean esperados.

- `ADAPTED`: los binarios publicados no son bytes oficiales Microsoft.
- El claim es security linting multi-lenguaje; no análisis dataflow/interprocedural completo.
- No reemplaza CodeQL autorizado, fuzzing, DAST, threat model, revisión ni ofensiva.
- No contiene GitHub entitlement, proyecto, baseline/suppressions ni aceptación de findings.
- El agente debe reabrir el gate ante drift de fuente, toolchain, advisory, regla o alcance.

V402 composed delta: COMPOSITION_SECURITY_RELEASE_V402.md/json: preserve admitted Go x/mod0.40 graph floor and original DevSkim short-error diagnostics; no new corporate authorship.
