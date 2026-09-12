# Microsoft BC bin replenishment — reconstruction evidence V165

## Scope and provenance

`GO-SUPPLY-FACTORY-INVENTORY-API` 0.9.0 adds a narrow bin-replenishment lane adapted from the MIT-licensed `microsoft/BCApps` commit `2eae56d704a1fd035d104f333602aea7091b7749`. The Go and SQL are local `ADAPTED` implementation, not Microsoft-authored code and not a claim of complete Business Central WMS equivalence.

The fixed official inputs were downloaded byte-for-byte:

| Microsoft file | SHA-256 |
|---|---|
| `src/Layers/W1/BaseApp/Warehouse/Structure/CalculateBinReplenishment.Report.al` | `4f961670cd1ba96a49a65d7e8585b4886675d65bf1b69841fe0421b505d9aca6` |
| `src/Layers/W1/BaseApp/Warehouse/Worksheet/Replenishment.Codeunit.al` | `e6d5e3873b4828db8dca49755670bd9815faa7ad1b0396189014af57380b7c4e` |
| `src/Layers/W1/BaseApp/Warehouse/Structure/BinContent.Table.al` | `841f172d7e92a0c581624f7e3b466e6e2ca59807d8d32201d82e1d38e3baa904` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMMovement.Codeunit.al` | `0e63996a0816d9379bdf76acc5712990f10624de81d875f2e499e6944b2955ed` |

Microsoft Learn references fixed in the pack cover Calculate Bin Replenishment, movement worksheets, advanced warehouse movements and bin-content setup.

## Executable contract

- target is a fixed, unblocked PICK/PUTPICK bin with explicit min/max;
- planning starts only below minimum and targets `max - current - open inbound movement`;
- candidates are same organization/item, lower rank, non-RECEIVE/SHIP and net of reservations;
- blocked bins, blocked lots and expired lots never become lines;
- FEFO is explicit; exact replay is idempotent and divergent replay conflicts;
- source quantities are reserved atomically; registration produces paired movement-out/in entries and cancellation releases every reservation;
- serialization by organization/item prevents concurrent over-planning.

Breakbulk/UOM conversion, cross-docking, cubage/weight, warehouse class, barcode/device execution, planning/CTP and external WMS reconciliation are not implemented and remain conditioned.

## Reconstruction and gates

- supply pack materialized 52/52 files; every byte matched the authoring tree;
- backend composed 344 files from 29 packs from Markdown only;
- official Go 1.26.7 archive: 74,955,002 bytes, SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`;
- official PostgreSQL 18.6 archive: 343,808,005 bytes, SHA-256 `fbe23da234ee31547bf8a36d29dfd81e82b849df2d2b78d2eecb43d360252f8c`; isolated loopback cluster with data checksums;
- migrations 0001–0031 applied cleanly; 0031 down/up passed after cleanup;
- focal PostgreSQL flow passed twice consecutively, then again after down/up;
- full `go test ./...`, `go vet ./...` and `go build ./cmd/...` passed;
- the focal corpus proves FEFO 5+3 to reach 10, exact/divergent replay, four paired ledger entries, consumed reservations, target-full rejection, multi-source cancellation, total reserved zero and exactly one concurrent winner.
- root `VERIFY_LIBRARY_PASS` reports 94 packs, 1,061 materialized files, 521 Markdown files and backend 344; `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` reports 121 governed sources and passed the available offline gates with Go 1.26.7 supplied explicitly.

No live WMS, operator device, load target or production environment was exercised; production admission remains project-specific and fail-closed.
