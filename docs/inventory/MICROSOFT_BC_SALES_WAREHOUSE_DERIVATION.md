# Microsoft BC sales-order to warehouse derivation

This connected demand lane is an `ADAPTED` Go/PostgreSQL implementation. It is not Microsoft-authored Go and it is not a Business Central runtime.

Exact authority is the MIT-licensed `microsoft/BCApps` commit `2eae56d704a1fd035d104f333602aea7091b7749`, tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`:

- `WhseCreateSourceDocument.Codeunit.al`, SHA-256 `df2a5f4767b58d4d307a223dcfe1c3d0176c9d336d6f9ebdd25ec414f40672f8`, sets quantity and base quantity, initializes outstanding quantities and checks the source document line and bin.
- `SalesWhsePostShipment.Codeunit.al`, SHA-256 `84121bd38555f3a7cba59dd97f985ae41af2da349359a60217bcab14ee7c3747`, resolves the sales document and line and reconciles quantity-to-ship/base quantity.
- `SalesLine.Table.al`, SHA-256 `ca65615dc05bc555a878ba2635c955897930332ef171e313a00103074fbc7278`, owns variant, UOM, quantity, outstanding quantity and their base forms.
- `SCMWarehousePick.Codeunit.al`, SHA-256 `f856e941d86a78a09d233aba6dfe5470dacbf08c853c26bc5e26512c3d6881c5`, includes released sales-order, partial-pick, registration and multiple-pick scenarios.
- `SalesOrderWhseValidateLine.Codeunit.al`, SHA-256 `12b4a3ce4fdef9592aad3054f5f174c2c4493d8377216cd4bea6890289d66deb`, verifies warehouse coupling when variant or UOM changes.

The portable contract therefore requires an immutable variant-to-item/sales-UOM binding, a placed operational sales line, exact organization/item/source identity, base-quantity conversion using the fixed item UOM, subtraction of every open or registered pick, atomic refusal beyond outstanding quantity, exact request replay, divergent replay refusal and durable pick evidence. Shipment and sale completion are not claimed by this lane.
