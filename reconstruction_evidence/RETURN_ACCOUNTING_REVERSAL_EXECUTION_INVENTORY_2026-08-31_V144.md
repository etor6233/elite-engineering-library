# Return Accounting Reversal Execution — V144

Date: 2026-08-31  
Result: `REBUILD_VERIFIED / CONDITIONED`

## Scope and provenance

`GO-RETURN-ACCOUNTING-REVERSAL-WORKER` 0.1.0 adds eight `AUTHORED` files. It reuses the existing Microsoft-BCApps-governed ledger and never invents account codes or amounts. Current Microsoft Dynamics 365 documentation governs original-sale linkage, return credit/correction and equal/opposite journal reversal; PostgreSQL and AWS govern locking, durable retry and reconciliation.

## Executable evidence

- pack materializes 8/8 byte-exact files;
- author tree passes focused Go tests, vet and worker build;
- PostgreSQL 18.6 applies migrations 0001–0022 and the 0022 SQL gate;
- the first return reverses the unique posted `SALE` journal with two exactly swapped entries and closes posting/attempt/execution evidence;
- the second scenario commits the reversal first, simulates a worker crash, then a new lease reconciles the deterministic journal without creating another reversal;
- immutable return-accounting evidence rejects mutation.

## Admission boundary

The worker requires succeeded inventory and refund/exchange effects and one original posted journal. It proves exact reversal of what was posted; it does not prove the original journal contained legally approved revenue, cost, tax or currency treatment. Those mappings and the posting period require project/accountant approval before activation. ARCA credit note, statutory books, live provider settlement, audited close and production acceptance remain separate gates.

## Failure memory

`LIB-FAIL-1425` and `LIB-FAIL-1426` record that the first Go gate ran from the Markdown root and that the final grouped window ended after tests/vet but before observable builds. The same files passed from the confirmed backend target and all seven mains then built in a separate gate; no incomplete command was misreported as evidence.
