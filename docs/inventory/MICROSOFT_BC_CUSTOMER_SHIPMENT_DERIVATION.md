# Microsoft BC customer shipment derivation

## Admitted claim

This implementation posts one immutable customer-sales shipment from one registered `customer-order` warehouse pick. In one serializable transaction it verifies the exact order, line, organization, sellable-variant binding, item, sales UOM and SHIP-bin reservations; consumes physical quantity, handling-unit composition and FIFO cost; increments line-level shipped quantity; records partial or complete order fulfillment; emits durable allocation and outbox evidence; and supports exact replay while rejecting divergent replay and competing posts.

It does **not** claim carrier pickup/delivery, customer handover, invoicing, tax authorization, payment capture, serialized-unit shipment or project-specific credit policy. Those remain separate owners and gates.

## Exact official authority

- Repository: `microsoft/BCApps`.
- Commit: `2eae56d704a1fd035d104f333602aea7091b7749`.
- Tree: `f9846fb1254c6311985c6131bb2af11c7f169e1b`.
- Root license: MIT.
- `src/Layers/W1/BaseApp/Sales/Posting/SalesPost.Codeunit.al`: 780,050 bytes, SHA-256 `f0f3a8a4c914e2257a8186abe0c80a56b060b05a694d2e3cfabdfbf711103341`.
- `src/Layers/W1/BaseApp/Sales/History/SalesShipmentHeader.Table.al`: 54,368 bytes, SHA-256 `9ea6772c4e470ee65cbb033a8e99ca6154607851ec3d8c0704f4c25deaa1dd41`.
- `src/Layers/W1/BaseApp/Sales/History/SalesShipmentLine.Table.al`: 78,739 bytes, SHA-256 `607cc7697cd019fbecf82bfa7f0c8b287b5a444186184b2c4f44a0ee3eb1550d`.
- `src/Layers/W1/BaseApp/Sales/SalesWhsePostShipment.Codeunit.al`: 46,608 bytes, SHA-256 `84121bd38555f3a7cba59dd97f985ae41af2da349359a60217bcab14ee7c3747`.
- `src/Layers/W1/Tests/ERM-Sales/ERMSalesOrder.Codeunit.al`: 376,840 bytes, SHA-256 `817bc8bf7f3efbd46b6e0e095a79f4e7fb3dffe501fc86400ce44ee04c62829b`.

The admitted authority creates posted shipment header/line records, initializes shipment lines from sales lines, transfers `Qty. to Ship` and base quantity, updates shipped quantities, and verifies the sales-line/posted-shipment relationship. No AL source was copied into this Go/PostgreSQL implementation.

## Local adaptation

- `sales.customer_shipment` and `sales.customer_shipment_line` are immutable posted-document evidence.
- `inventory.customer_shipment_allocation` links each posted line to the exact warehouse reservation and outbound inventory entry.
- `sales.customer_order_line.shipped_quantity` is independent from the commercial order state.
- `sales.customer_order.fulfillment_state` preserves `unfulfilled`, `partially-shipped`, `shipped` and the later delivery boundary without overwriting payment/order state.
- `consumeRegisteredPickPackaging` normalizes a non-base handling unit to its exact base quantity before the physical balance trigger consumes it. This same owner repairs transfer shipment after packaging support was added.
- Only FIFO bulk items are admitted by this vertical. Specific-cost and serialized units require their own exact source binding instead of a guessed fallback.

## Executable gates

- Open/draft/unrelated activities are rejected.
- Exact request replay returns the same posted shipment; changed fields are rejected.
- A registered pick can be posted once; concurrent distinct requests have one winner.
- Partial shipment increments line quantity and leaves fulfillment partial.
- Final shipment completes fulfillment only when every order line is fully shipped.
- Physical balance, UOM composition, reservations, cost layers, applications, posted cost and outbox evidence reconcile in the same transaction.
- Transfer shipment proves non-base BOX composition consumption after explicit breakbulk.
- Migration `0038` passes clean up, down and up paths.

## Residual production conditions

Project promotion still requires project-specific payment/credit policy, carrier adapter and webhook contracts, serialized vehicle flow, customer handover, fiscal posting, authorization matrix, load, recovery, observability and operational acceptance. This narrow proof must not be described as production readiness for an entire franchise.
