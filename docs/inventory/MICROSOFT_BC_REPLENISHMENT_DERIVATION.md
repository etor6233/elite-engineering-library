# Microsoft Business Central bin-replenishment derivation

This implementation is `ADAPTED`, not Microsoft-authored Go or SQL. It is governed by the MIT-licensed Microsoft BCApps files fixed at commit `2eae56d704a1fd035d104f333602aea7091b7749`: `CalculateBinReplenishment.Report.al`, `Replenishment.Codeunit.al`, `BinContent.Table.al` and the official `SCMMovement.Codeunit.al` tests, together with Microsoft Learn warehouse-movement and bin-content documentation.

The admitted contract is narrow and executable. A destination must be a fixed PICK/PUTPICK bin with explicit minimum and maximum quantities. Replenishment is needed only when current quantity plus already planned inbound movement is below the minimum; its desired quantity is maximum minus current minus planned inbound. Sources belong to the same organization/item, have lower bin ranking, are not RECEIVE/SHIP, are not movement-blocked and expose quantity net of reservations. The plan reserves exact source lot quantities and creates one open movement instruction. Registration atomically moves those quantities and cancellation atomically releases them. Exact request replay returns the same instruction; a divergent replay or a concurrent second plan fails closed.

When `use_fefo=true`, source lots are ordered by earliest non-null expiration before bin rank; otherwise candidates follow the upstream higher-lower-bin-rank rule. This flag is explicit because the portable schema does not invent a Business Central Location card. The target UOM is resolved before planning. Exact target-UOM composition is selected first; a larger source UOM is considered only under explicit `allow_breakbulk=true`, following the fixed BCApps `AllowBreakbulk` branch. Each activity line retains From/To UOM quantities and factors, reserves physical and source-composition quantities together, records immutable Take/Place conversion evidence on registration, preserves the exact remainder and releases both reservations on cancellation. Cubage/weight, warehouse class, barcode/device execution, planning/CTP and external WMS reconciliation remain conditioned; cross-docking is governed by its separate admitted derivation.

Official source identities and packaged SHA-256 values:

- `CalculateBinReplenishment.Report.al`: `4f961670cd1ba96a49a65d7e8585b4886675d65bf1b69841fe0421b505d9aca6`;
- `Replenishment.Codeunit.al`: `e6d5e3873b4828db8dca49755670bd9815faa7ad1b0396189014af57380b7c4e`;
- `BinContent.Table.al`: `841f172d7e92a0c581624f7e3b466e6e2ca59807d8d32201d82e1d38e3baa904`;
- `SCMMovement.Codeunit.al`: `0e63996a0816d9379bdf76acc5712990f10624de81d875f2e499e6944b2955ed`.

Migration 0036 and `warehouse_packaging_flow_integration_test.go` prove exact target-UOM preference, fail-closed implicit breakbulk, authorized BOX→EA replenishment, base/composition conservation, immutable conversion provenance, registration and cancellation. This is a local verified adaptation of those Microsoft invariants, not Microsoft-authored Go or SQL.
