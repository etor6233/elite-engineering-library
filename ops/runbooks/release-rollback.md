# Release and rollback runbook

Before release, record the source revision, immutable artifact digest, builder identity, SBOM, provenance, vulnerability/license/secret scans, migration compatibility, restore evidence, approver and planned health gates. Never deploy a mutable tag as the release identity.

Use expand-migrate-contract: deploy additive schema first, backfill with resumable jobs, shift reads/writes, observe, and remove old schema only in a later release. A rollback must not require reversing a destructive migration.

During canary, compare availability, tail latency, error-budget burn, saturation, queue/outbox lag and business invariants with the stable population. Stop promotion on any failed gate. Roll back to the previously recorded digest using the exact argv command in the readiness record; do not rebuild it.

After rollback, verify traffic, database compatibility, job ownership, event publication and external side effects. Preserve telemetry and an incident timeline. A rollback test is complete only when it proves the old artifact can serve against the current compatible schema.
