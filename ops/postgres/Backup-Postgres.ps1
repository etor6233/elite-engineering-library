[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)][ValidatePattern('^[A-Za-z0-9_.-]+$')][string] $Database,
  [Parameter(Mandatory = $true)][string] $BackupDirectory,
  [string] $HostName = '127.0.0.1',
  [ValidateRange(1, 65535)][int] $Port = 5432,
  [Parameter(Mandatory = $true)][ValidatePattern('^[A-Za-z0-9_.-]+$')][string] $Username,
  [Parameter(Mandatory = $true)][string] $PgBin,
  [ValidatePattern('^[A-Za-z0-9_.-]+$')][string] $LogicalName = 'platform'
)

$ErrorActionPreference = 'Stop'
$pgDump = Join-Path $PgBin 'pg_dump.exe'
$pgRestore = Join-Path $PgBin 'pg_restore.exe'
if (-not (Test-Path -LiteralPath $pgDump -PathType Leaf)) { throw "pg_dump not found in PgBin" }
if (-not (Test-Path -LiteralPath $pgRestore -PathType Leaf)) { throw "pg_restore not found in PgBin" }

$backupRoot = [IO.Path]::GetFullPath($BackupDirectory)
if (-not (Test-Path -LiteralPath $backupRoot)) {
  [void](New-Item -ItemType Directory -Path $backupRoot)
}

$stamp = [DateTimeOffset]::UtcNow.ToString('yyyyMMddTHHmmssZ')
$nonce = [Guid]::NewGuid().ToString('N').Substring(0, 8)
$baseName = "${LogicalName}_${stamp}_${nonce}"
$backupPath = Join-Path $backupRoot "${baseName}.dump"
$manifestPath = Join-Path $backupRoot "${baseName}.manifest.json"

& $pgDump --host $HostName --port $Port --username $Username --dbname $Database --format custom --compress 'zstd:6' --no-owner --no-acl --file $backupPath
if ($LASTEXITCODE -ne 0) { throw "pg_dump failed with exit code $LASTEXITCODE" }
if (-not (Test-Path -LiteralPath $backupPath -PathType Leaf)) { throw 'pg_dump did not create the backup' }

$hash = (Get-FileHash -LiteralPath $backupPath -Algorithm SHA256).Hash.ToLowerInvariant()
$toc = & $pgRestore --list $backupPath
if ($LASTEXITCODE -ne 0) { throw "pg_restore --list failed with exit code $LASTEXITCODE" }
$serverVersion = (& $pgDump --version) -join "`n"

$manifest = [ordered]@{
  schema_version = 'elite.postgres.backup-manifest.v1'
  created_at_utc = [DateTimeOffset]::UtcNow.ToString('O')
  logical_name = $LogicalName
  source_database = $Database
  source_host = $HostName
  source_port = $Port
  tool_version = $serverVersion.Trim()
  format = 'custom'
  sha256 = $hash
  bytes = (Get-Item -LiteralPath $backupPath).Length
  toc_entry_count = @($toc | Where-Object { $_ -match '^\d+;' }).Count
  backup_file = [IO.Path]::GetFileName($backupPath)
}
$json = $manifest | ConvertTo-Json -Depth 5
[IO.File]::WriteAllText($manifestPath, $json + "`n", [Text.UTF8Encoding]::new($false))

[pscustomobject]@{
  BackupPath = $backupPath
  ManifestPath = $manifestPath
  Sha256 = $hash
  Bytes = $manifest.bytes
  TocEntryCount = $manifest.toc_entry_count
}
