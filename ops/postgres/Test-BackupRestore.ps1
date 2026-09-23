[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)][ValidatePattern('^[A-Za-z0-9_.-]+$')][string] $SourceDatabase,
  [string] $HostName = '127.0.0.1',
  [ValidateRange(1, 65535)][int] $Port = 5432,
  [Parameter(Mandatory = $true)][ValidatePattern('^[A-Za-z0-9_.-]+$')][string] $Username,
  [Parameter(Mandatory = $true)][string] $PgBin,
  [Parameter(Mandatory = $true)][string] $ArtifactDirectory
)

$ErrorActionPreference = 'Stop'
$scriptRoot = $PSScriptRoot
$backupCommand = Join-Path $scriptRoot 'Backup-Postgres.ps1'
$drillCommand = Join-Path $scriptRoot 'Invoke-RestoreDrill.ps1'
$invariants = Join-Path $scriptRoot 'assert-platform-invariants.sql'

$backup = & $backupCommand -Database $SourceDatabase -BackupDirectory $ArtifactDirectory -HostName $HostName -Port $Port -Username $Username -PgBin $PgBin -LogicalName 'platform-foundation'
if (-not $backup -or $backup.TocEntryCount -lt 1 -or $backup.Bytes -lt 1) { throw 'Backup smoke assertions failed' }

$drill = & $drillCommand -BackupPath $backup.BackupPath -ManifestPath $backup.ManifestPath -InvariantFile $invariants -HostName $HostName -Port $Port -Username $Username -PgBin $PgBin
if (-not $drill -or $drill.Status -ne 'PASS' -or $drill.DatabaseRetained) { throw 'Restore drill smoke assertions failed' }

[pscustomobject]@{
  Status = 'PASS'
  BackupPath = $backup.BackupPath
  ManifestPath = $backup.ManifestPath
  Sha256 = $backup.Sha256
  RestoreDurationSeconds = $drill.DurationSeconds
}
