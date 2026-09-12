# PostgreSQL Backup and Restore Drill Core

## 1. Metadata

```yaml
pack_id: "PG-BACKUP-RESTORE-CORE"
pack_version: "0.1.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa backup lógico, restore aislado, manifiesto verificable e invariantes para la fundación PostgreSQL."
stacks:
  - "PostgreSQL 18.6 client tools"
  - "PowerShell 7+ operational tooling"
compatible_with:
  - "PG-TX-FOUNDATION >=0.1.0 <1.0.0"
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources:
  - "https://www.postgresql.org/docs/18/app-pgdump.html"
  - "https://www.postgresql.org/docs/18/app-pgrestore.html"
  - "https://www.postgresql.org/docs/18/libpq-pgpass.html"
verified_at: "2026-08-25"
```

Todos los bloques son `AUTHORED`. El pack usa herramientas oficiales de PostgreSQL pero no incorpora código de PostgreSQL. No promete PITR, failover ni un RPO/RTO específico: prueba que un backup lógico puede restaurarse en una base distinta y conservar invariantes mínimas.

## 2. Applicability

Use to prove a logical PostgreSQL backup can be hashed, restored into isolation and checked for application invariants. Reject it as the sole production disaster-recovery design: PITR, managed-service snapshots, encryption, retention, geographic/account isolation, failover and organizational RPO/RTO still require target-specific implementation.

## 3. Architecture contract

- la contraseña sólo se entrega mediante `PGPASSWORD`, `.pgpass` o el mecanismo seguro del entorno; nunca por argumento o Markdown;
- el script no imprime secretos ni connection strings;
- el restore crea una base aislada cuyo nombre comienza con `restore_drill_`;
- la base de ensayo se elimina en `finally`, salvo petición explícita de conservarla;
- backup, manifiesto y SHA-256 quedan juntos;
- el hash se valida antes de restaurar;
- `pg_restore --exit-on-error --single-transaction` impide aceptar estados parciales;
- el ensayo ejecuta invariantes después del restore;
- producción aún exige backup cifrado, retención, copia fuera de cuenta/región, PITR, monitoreo y medición de RPO/RTO.

## 4. Exact file manifest

```text
CREATE ops/postgres/Backup-Postgres.ps1
CREATE ops/postgres/Invoke-RestoreDrill.ps1
CREATE ops/postgres/assert-platform-invariants.sql
CREATE ops/postgres/Test-BackupRestore.ps1
CREATE ops/postgres/README.md
```

## 5. Materialization blocks

### FILE: `ops/postgres/Backup-Postgres.ps1`

```yaml
block_id: "PG-BACKUP-RESTORE:backup:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "cc6948549fe1a52ccde48689c3d4f80649703751b5da79ff6509fd58dfe3ca08"
variables: []
secrets_allowed: false
```

````powershell
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
````

### FILE: `ops/postgres/Invoke-RestoreDrill.ps1`

```yaml
block_id: "PG-BACKUP-RESTORE:drill:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "50276e572422c6c97205859b0597f69a42b4bf71fa5e672030e16e01c833558c"
variables: []
secrets_allowed: false
```

````powershell
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
````

### FILE: `ops/postgres/assert-platform-invariants.sql`

```yaml
block_id: "PG-BACKUP-RESTORE:invariants:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "174cd615b57329372d25532b468ffd5e562aa8f8450755260539a88e370b05b4"
variables: []
secrets_allowed: false
```

````sql
begin transaction read only;

do $$
declare
  missing_relations text[];
  invalid_count bigint;
begin
  select array_agg(required_name order by required_name)
    into missing_relations
  from unnest(array[
    'platform.tenant',
    'platform.business_profile_version',
    'platform.idempotency_record',
    'platform.outbox_event',
    'platform.consumer_inbox',
    'platform.job',
    'audit.event'
  ]) as required(required_name)
  where to_regclass(required_name) is null;

  if missing_relations is not null then
    raise exception 'missing required relations: %', missing_relations;
  end if;

  select count(*) into invalid_count
  from platform.business_profile_version p
  left join platform.tenant t on t.tenant_id = p.tenant_id
  where t.tenant_id is null;
  if invalid_count <> 0 then raise exception 'orphan business profiles: %', invalid_count; end if;

  select count(*) into invalid_count
  from (
    select tenant_id
    from platform.business_profile_version
    where status = 'active'
    group by tenant_id
    having count(*) > 1
  ) duplicated_active;
  if invalid_count <> 0 then raise exception 'tenants with multiple active profiles: %', invalid_count; end if;

  select count(*) into invalid_count
  from platform.outbox_event
  where aggregate_version <= 0 or available_at < occurred_at;
  if invalid_count <> 0 then raise exception 'invalid outbox rows: %', invalid_count; end if;

  select count(*) into invalid_count
  from platform.job
  where attempts > max_attempts
     or ((claimed_by is null) <> (claimed_until is null))
     or (completed_at is not null and terminal_error_code is not null);
  if invalid_count <> 0 then raise exception 'invalid job rows: %', invalid_count; end if;
end $$;

select 'platform_restore_invariants' as check_name, 'PASS' as result;
rollback;
````

### FILE: `ops/postgres/Test-BackupRestore.ps1`

```yaml
block_id: "PG-BACKUP-RESTORE:test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c1aabf6c78ad96f8c38beda1f5fd9840d3c8c1566390938499de71fdb4b5c0bd"
variables: []
secrets_allowed: false
```

````powershell
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
````

### FILE: `ops/postgres/README.md`

```yaml
block_id: "PG-BACKUP-RESTORE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "daa77d24cedebcab9cfea4ea66f5086e4bf0e1b7bf570154d75947c15560d3bd"
variables: []
secrets_allowed: false
```

````markdown
# PostgreSQL backup and restore drill

Prerequisites: PowerShell 7+, PostgreSQL 18.6 client tools, connectivity to the target and credentials supplied through `PGPASSWORD`, `.pgpass` or the platform secret injector.

Run `Test-BackupRestore.ps1` against a non-production database containing the platform foundation. The command creates a custom-format dump and JSON manifest, verifies its hash, restores it into an isolated database, executes read-only invariants and removes the drill database.

Do not treat one passing logical restore as a complete disaster-recovery strategy. Before production, add encrypted immutable storage, retention and deletion policy, scheduled restore drills from the actual backup service, WAL/PITR coverage, cross-account or cross-region copies, access audit, alerting, a dependency-aware application recovery runbook, and measured RPO/RTO under representative data volume.
````

## 6. Configuration surface

| Parameter/input | Type/default | Secret | Validation/effect |
|---|---|---|---|
| database/host/port/user | validated CLI values / local host where stated | user may be sensitive | identifies source; password is never an argument |
| `PGPASSWORD`/`.pgpass`/injector | external credential | yes | consumed by official clients; never logged |
| `PgBin` | absolute tool directory / none | no | required `pg_dump`, `pg_restore`, `psql`, `createdb`, `dropdb` |
| backup directory/logical name | path/safe token | none | no | output kept together and traversal avoided |
| keep drill database | explicit switch / false | no | default cleanup in `finally` |

## 7. Dependency bill

| Tool | Pin | Use | License | Scope | Source |
|---|---|---|---|---|---|
| PostgreSQL client tools | `18.6` verified baseline | dump/restore/query/database lifecycle | PostgreSQL | operations | `postgresql.org` |
| PowerShell | `7+` | orchestration/tests | MIT | operations | `github.com/PowerShell/PowerShell` |

## 8. Apply order

Materialize after `PG-TX-FOUNDATION`, provision a non-production source with external credentials, run backup, verify manifest/hash, restore to a generated `restore_drill_*` database, execute invariants and prove cleanup. Existing systems must replace/add domain invariants. Rollback means deleting only verified drill artifacts/databases; never delete the source or production backup set.

## 9. Verification

Required reconstruction gates:

1. materializer validates every embedded SHA-256;
2. run platform migration and SQL tests on PostgreSQL 18.6;
3. create a logical backup and manifest;
4. restore into a separately named database;
5. execute invariant SQL;
6. prove the drill database is removed;
7. corrupt a copy of the dump and prove hash validation fails before restore.

Los siete gates fueron ejecutados localmente en PostgreSQL 18.6 y quedaron registrados en `reconstruction_evidence/POSTGRES_BACKUP_RESTORE_CORE_2026-08-24_V1.md`; por eso la implementación es `REBUILD_VERIFIED`. La admisión sigue `CONDITIONED`: producción requiere almacenamiento, cifrado, scheduling, PITR, failover y RPO/RTO específicos de la organización.

## 10. Reconstruction evidence

Commands, hashes, isolated restore, invariants, corruption rejection and cleanup are recorded in `reconstruction_evidence/POSTGRES_BACKUP_RESTORE_CORE_2026-08-24_V1.md`; the final audit rechecks the current pack.
