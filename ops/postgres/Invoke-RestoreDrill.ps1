[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)][string] $BackupPath,
  [Parameter(Mandatory = $true)][string] $ManifestPath,
  [Parameter(Mandatory = $true)][string] $InvariantFile,
  [string] $HostName = '127.0.0.1',
  [ValidateRange(1, 65535)][int] $Port = 5432,
  [Parameter(Mandatory = $true)][ValidatePattern('^[A-Za-z0-9_.-]+$')][string] $Username,
  [ValidatePattern('^[A-Za-z0-9_.-]+$')][string] $AdminDatabase = 'postgres',
  [Parameter(Mandatory = $true)][string] $PgBin,
  [switch] $KeepRestoredDatabase
)

$ErrorActionPreference = 'Stop'
$createdb = Join-Path $PgBin 'createdb.exe'
$dropdb = Join-Path $PgBin 'dropdb.exe'
$pgRestore = Join-Path $PgBin 'pg_restore.exe'
$psql = Join-Path $PgBin 'psql.exe'
foreach ($tool in @($createdb, $dropdb, $pgRestore, $psql)) {
  if (-not (Test-Path -LiteralPath $tool -PathType Leaf)) { throw "Required PostgreSQL tool not found in PgBin" }
}
foreach ($file in @($BackupPath, $ManifestPath, $InvariantFile)) {
  if (-not (Test-Path -LiteralPath $file -PathType Leaf)) { throw "Required file not found" }
}

$manifest = Get-Content -LiteralPath $ManifestPath -Raw | ConvertFrom-Json
if ($manifest.schema_version -ne 'elite.postgres.backup-manifest.v1') { throw 'Unsupported backup manifest' }
$actualHash = (Get-FileHash -LiteralPath $BackupPath -Algorithm SHA256).Hash.ToLowerInvariant()
if ($actualHash -ne $manifest.sha256) { throw 'Backup SHA-256 does not match manifest' }
if ([IO.Path]::GetFileName($BackupPath) -ne $manifest.backup_file) { throw 'Backup filename does not match manifest' }

$drillDatabase = 'restore_drill_' + [DateTimeOffset]::UtcNow.ToString('yyyyMMddHHmmss') + '_' + [Guid]::NewGuid().ToString('N').Substring(0, 8)
if ($drillDatabase -notmatch '^restore_drill_[0-9]{14}_[0-9a-f]{8}$') { throw 'Unsafe drill database name' }
$startedAt = [DateTimeOffset]::UtcNow
$created = $false
$succeeded = $false

try {
  & $createdb --host $HostName --port $Port --username $Username --maintenance-db $AdminDatabase --encoding UTF8 --template template0 $drillDatabase
  if ($LASTEXITCODE -ne 0) { throw "createdb failed with exit code $LASTEXITCODE" }
  $created = $true

  & $pgRestore --host $HostName --port $Port --username $Username --dbname $drillDatabase --exit-on-error --single-transaction --no-owner --no-acl $BackupPath
  if ($LASTEXITCODE -ne 0) { throw "pg_restore failed with exit code $LASTEXITCODE" }

  & $psql --host $HostName --port $Port --username $Username --dbname $drillDatabase --no-psqlrc --set ON_ERROR_STOP=1 --file $InvariantFile
  if ($LASTEXITCODE -ne 0) { throw "Invariant validation failed with exit code $LASTEXITCODE" }
  $succeeded = $true
}
finally {
  if ($created -and -not $KeepRestoredDatabase) {
    & $dropdb --host $HostName --port $Port --username $Username --maintenance-db $AdminDatabase --if-exists $drillDatabase
    if ($LASTEXITCODE -ne 0 -and $succeeded) { throw "Restore passed but cleanup failed with exit code $LASTEXITCODE" }
  }
}

[pscustomobject]@{
  Status = if ($succeeded) { 'PASS' } else { 'FAIL' }
  DrillDatabase = $drillDatabase
  BackupSha256 = $actualHash
  StartedAtUtc = $startedAt.ToString('O')
  FinishedAtUtc = [DateTimeOffset]::UtcNow.ToString('O')
  DurationSeconds = [math]::Round(([DateTimeOffset]::UtcNow - $startedAt).TotalSeconds, 3)
  DatabaseRetained = [bool]$KeepRestoredDatabase
}
