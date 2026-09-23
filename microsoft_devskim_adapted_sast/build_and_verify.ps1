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
