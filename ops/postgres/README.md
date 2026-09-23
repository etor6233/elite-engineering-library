# PostgreSQL backup and restore drill

Prerequisites: PowerShell 7+, PostgreSQL 18.6 client tools, connectivity to the target and credentials supplied through `PGPASSWORD`, `.pgpass` or the platform secret injector.

Run `Test-BackupRestore.ps1` against a non-production database containing the platform foundation. The command creates a custom-format dump and JSON manifest, verifies its hash, restores it into an isolated database, executes read-only invariants and removes the drill database.

Do not treat one passing logical restore as a complete disaster-recovery strategy. Before production, add encrypted immutable storage, retention and deletion policy, scheduled restore drills from the actual backup service, WAL/PITR coverage, cross-account or cross-region copies, access audit, alerting, a dependency-aware application recovery runbook, and measured RPO/RTO under representative data volume.
