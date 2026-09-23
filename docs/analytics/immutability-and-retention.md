# Analytics evidence guards and owner retention

Scope: additive migration0087 of GO-DATA-ANALYTICS-CORE. Historical0041 remains byte-identical. This is AUTHORED PostgreSQL glue; Go Fresh is a separate successor change.

Runtime contract:

- Landing identity, payload, source hash, event/received timestamps and lineage are immutable. Only state remains mutable under the original landed/validated/rejected enum. New fields default to immutable unless a future reviewed migration explicitly changes the allowlist.
- A metric definition is immutable per tenant/code/version. Insert a new version to change aggregation, source, dimensions or TTL. No-op updates are harmless and remain allowed.
- Normal INSERT/SELECT, bulk landing/state processing and metric_snapshot reconciliation remain available. This guard does not add tenant authorization, ingestion validation, a scheduler, an analytics warehouse or privacy policy.
- Runtime UPDATE of evidence/definition values, DELETE and TRUNCATE are rejected even if the role accidentally holds these DML privileges. Runtime must not own the schemas/tables/functions, inherit the owner role, have replication bypass or DDL rights.

Retention and trusted ownership:

The table/schema owner and database superuser are trusted control-plane principals. They can disable/drop guards and are outside the immutability claim. There is no custom GUC or caller-controlled bypass.

Retention is deliberately not an ordinary runtime DELETE. An authorized owner must first record the approved tenant, cutoff, record scope/count, retention policy/legal hold decision, archive/restore evidence and retention receipt location. This change neither invents a country retention period nor authorizes deleting live data.

The scoped procedure is: begin a maintenance transaction; acquire the table lock; disable only the relevant row-delete guard; delete the exact approved records using tenant and cutoff plus bounded identity selection; re-enable the same guard before commit; verify expected affected count and preserved data; commit; persist the external operational receipt and rerun the negative runtime deletion probe. Any failure rolls back the entire transaction. Never use a broad TRUNCATE for routine retention. Metric definitions referenced by snapshots retain their foreign-key protection and require a separately approved dependency-aware retention plan.

Fixture evidence exercises this owner procedure against one expired synthetic landing row, then proves the guard is enabled again. This is not a production retention rehearsal, privacy certification, or immutable protection against a privileged owner.

Migration and rollback:

1. Apply after0041 in the declared composition. 0087 is reserved after the observed maximum0086; missing0041 fails.
2. Normal downgrade: execute0087.down before0041.down. It preserves all evidence rows and removes the new guards; the immutability guarantee is withdrawn until reapplication.
3. Explicit DROP remains an owner operation. Guard functions live in analytics_guard so the historical0041.down can drop data/analytics without new functions obstructing those schemas;0087.down then removes remaining guard functions/schema and tolerates already absent tables.
4. Never claim DROP/owner/replication protection from row or TRUNCATE triggers. The target must preserve role separation and monitor unauthorized DDL via its existing operations owner.

Qualification:

`db/tests/0087_analytics_immutability.test.sql` is self-contained and transactional, intended only for a disposable fixture database after the migration. A fixture owner creates a temporary test role/seed, then tests under that non-superuser role. Roles, grants and rows roll back. Do not run this fixture against production. It covers20 protected mutations; the external qualification also tests load/state, new metric version/snapshot, owner retention, downgrade/reapply and historical DROP.
