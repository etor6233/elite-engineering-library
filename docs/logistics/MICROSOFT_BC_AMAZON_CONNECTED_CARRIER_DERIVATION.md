# Connected carrier delivery — derivation record

## Narrow claim

This vertical connects an already posted `sales.customer_shipment` to one durable carrier shipment, records normalized provider reports idempotently, advances only through the allowed transport states, and marks the customer order delivered only when every posted customer shipment has a delivered transport. It never creates or accepts `sales.delivery_handover`.

The Go and SQL are local `ADAPTED` implementation. They are not copied from Microsoft or Amazon. The sources below govern only the stated boundaries.

## Fixed Microsoft authority

Repository: `microsoft/BCApps` at commit `2eae56d704a1fd035d104f333602aea7091b7749`, tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`, MIT.

| Exact path | Bytes | SHA-256 | Narrow authority |
|---|---:|---|---|
| `src/Layers/W1/BaseApp/Foundation/Shipping/ShippingAgent.Table.al` | 3,982 | `1f18ed96830b0bb8af68f22012661f886f48b49403b354eb5f923fa37aa54d73` | shipping agent and tracking URL are explicit provider-owned concepts |
| `src/Layers/W1/BaseApp/Foundation/Shipping/ShippingAgentServices.Table.al` | 2,473 | `60d79a7c71e3d32e03a9e275160f7680810f3c8e4849a3e51cf160899462c9ee` | provider service is distinct from provider identity |
| `src/Layers/W1/BaseApp/Sales/History/SalesShipmentHeader.Table.al` | 54,368 | `9ea6772c4e470ee65cbb033a8e99ca6154607851ec3d8c0704f4c25deaa1dd41` | posted sales shipment preserves agent, service and package tracking number |
| `src/Layers/W1/Tests/SMB/O365ShippingAgent.Codeunit.al` | 23,865 | `d99105e91b43b234cd5b18d0f90f78e1ec758b85317a2b79c6877038e4caeda2` | sales order shipping fields survive posting into the sales shipment |
| `src/Layers/W1/Tests/SCM-Reservation/SCMPackageTrackingSales.Codeunit.al` | 130,076 | `96ed3a37ace2a44f14d835d6a75676f264900d41750544ebc48cbabc7a5df1bb` | package/lot/serial reservation and partial shipment remain quantity-bound |

Microsoft does not govern this repository's Go API, SQL schema, state machine, authorization or order-completion rule.

## Fixed Amazon adapter authority

The provider mapping is admitted only through already verified library packs built from Amazon's official SDK commit `8e792ae345a8d334ccdbdd03181f05f040e6a4fc`, model commit `8e429486005c4ebdce5099e48cc48515a65359bb` and exact official reference pages:

- `PYTHON-AMAZON-SPAPI-SHIPPING-TRACKING-ADAPTER 0.6.x`: tracking receipts are hash-linked, redacted and non-authoritative for automatic delivery completion.
- `PYTHON-AMAZON-SPAPI-EASYSHIP-HANDOVER-ADAPTER 0.1.x`: exact provider statuses are reconciled; provider pickup/delivery never means internal handover acceptance.
- `PYTHON-AMAZON-SPAPI-FULFILLMENT-DELIVERY-EVIDENCE-ADAPTER 0.1.x`: approved delivery evidence excludes signed URLs and recipient identity and never invents business acceptance.

The admitted Easy Ship mapping is closed:

| Exact provider status | Internal transport state |
|---|---|
| `PickedUp` | `picked-up` |
| `AtOriginFC` | `in-transit` |
| `AtDestinationFC` | `in-transit` |
| `OutForDelivery` | `in-transit` |
| `Delivered` | `delivered` |
| `Rejected` | `exception` |
| `Undeliverable` | `exception` |
| `ReturnedToSeller` | `exception` |
| `LostInTransit` | `exception` |
| `DamagedInTransit` | `exception` |

Unknown status, unknown schema, malformed evidence hash or `automatic_business_acceptance=true` fails closed.

## Local invariants

1. A customer transport references one real posted customer shipment, its order, source organization and customer.
2. Request replay is exact; divergent reuse conflicts.
3. A provider reference is unique only after it exists. Multiple planned shipments with `NULL` references are valid.
4. Provider reports are immutable and unique per shipment version. A concurrent stale report loses.
5. State changes are monotonic under the explicit transition table. A report cannot skip from `planned` directly to `delivered`.
6. Order fulfillment becomes `delivered` only when no posted customer shipment lacks a delivered carrier record.
7. Provider delivery emits durable audit/outbox evidence but creates zero `sales.delivery_handover` rows.
8. Serial/specific allocation and customer checklist/acceptance remain a separate future vertical; bulk FIFO shipment does not pretend to prove them.

## Executed evidence required before promotion

- migrations `0001` through `0039` on fresh PostgreSQL 18.6;
- SQL regression proving two null provider references and rejection of a duplicated non-null reference;
- exact request replay and divergent replay rejection;
- concurrent reports with exactly one version winner;
- provider journey `PickedUp → AtOriginFC/AtDestinationFC → OutForDelivery → Delivered`;
- exact evidence replay and divergent evidence rejection;
- delivered order, four immutable provider reports and zero delivery handovers;
- complete Go tests, vet and build;
- migration `0039` down/up and focal SQL test.

## Production conditions

This reusable code does not prove a project-specific carrier account, marketplace eligibility, webhook authentication, polling schedule, address/privacy policy, evidence retention, SLA, load, outage reconciliation or customer acceptance workflow. Those remain project gates and cannot be inferred from the library PASS.
