# Go Supply, Factory and Inventory API

## 1. Metadata

```yaml
pack_id: "GO-SUPPLY-FACTORY-INVENTORY-API"
pack_version: "0.16.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Añade compras, fábrica, stock serial y bulk, lotes, bins, UOM base/alternativa, composición física exacta, conversión explícita breakbulk/gather, recepción→put-away y picking/reposición con empaque preservado o breakbulk autorizado, reservas físicas/de composición, Take/Place y cancelación; conecta pedido comercial colocado→binding variante/item/UOM→pendiente base→pick registrado→customer shipment inmutable parcial/completo, con consumo atómico de stock/UOM/costo, replay y concurrencia, y conserva picking FEFO, reposición min/max, cross-docking de transferencia, ATP, FIFO/específico y transferencias bulk cross-organization con posting parcial, tránsito durable, lote y costo específico exactos."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "pgx 5.10.0"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "GO-ELECTROMOBILITY-PUBLIC-CRM-API >=0.2.0", "ELECTROMOBILITY-FRANCHISE-MODULES 0.1.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner AND MIT"
upstream_sources: ["https://go.dev", "https://www.postgresql.org/docs/18/", "https://github.com/jackc/pgx", "https://github.com/microsoft/BCApps/tree/2eae56d704a1fd035d104f333602aea7091b7749", "https://learn.microsoft.com/en-us/dynamics365/business-central/design-details-item-tracking", "https://learn.microsoft.com/en-us/dynamics365/business-central/design-details-item-tracking-in-the-warehouse", "https://learn.microsoft.com/en-us/dynamics365/business-central/design-details-warehouse-management", "https://learn.microsoft.com/en-us/dynamics365/business-central/design-details-inbound-warehouse-flow", "https://learn.microsoft.com/en-us/dynamics365/business-central/warehouse-how-to-put-items-away-with-warehouse-put-aways", "https://learn.microsoft.com/en-us/dynamics365/business-central/warehouse-how-to-pick-items-for-warehouse-shipment", "https://learn.microsoft.com/en-ca/dynamics365/business-central/warehouse-how-to-plan-warehouse-movements-in-worksheets", "https://learn.microsoft.com/en-ca/dynamics365/business-central/warehouse-how-to-calculate-bin-replenishment", "https://learn.microsoft.com/en-us/dynamics365/business-central/warehouse-how-to-move-items-in-advanced-warehousing", "https://learn.microsoft.com/en-us/dynamics365/business-central/warehouse-how-to-set-up-bin-contents", "https://learn.microsoft.com/en-gb/dynamics365/business-central/warehouse-how-to-cross-dock-items", "https://learn.microsoft.com/en-gb/dynamics365/business-central/inventory-how-setup-units-of-measure", "https://learn.microsoft.com/en-ca/dynamics365/business-central/warehouse-enable-automatic-breaking-bulk-with-directed-put-away-and-pick", "https://learn.microsoft.com/en-us/dynamics365/business-central/design-details-inventory-costing", "https://learn.microsoft.com/en-gb/dynamics365/business-central/design-details-reservation-order-tracking-and-action-messaging", "https://learn.microsoft.com/en-gb/dynamics365/business-central/design-details-availability-in-the-warehouse", "https://learn.microsoft.com/en-us/dynamics365/business-central/inventory-how-transfer-between-locations", "https://learn.microsoft.com/en-us/dynamics365/business-central/inventory-how-work-item-tracking", "https://learn.microsoft.com/en-us/dynamics365/business-central/design-details-transfers-in-planning"]
verified_at: "2026-09-01"
```

Los bloques de reserva/ATP, bulk/lot/bin/FIFO-específico, recepción/put-away/pick/FEFO, reposición, cross-docking y transferencia cross-organization son `ADAPTED` desde invariantes y suites Microsoft BCApps MIT fijadas; HTTP, documentación y wiring local permanecen `AUTHORED`. No se presentan como código Go escrito o soportado por Microsoft. Métodos LIFO/Average/Standard y capacidades WMS no implementadas fallan cerrados en vez de simularse. El pack no inventa reglas fiscales, Incoterms, homologación, QC, política contable ni propiedad legal: esas políticas siguen condicionadas al proyecto.

## 2. Applicability

Use when the platform owns purchase orders, factory-unit progress, serialized vehicles and bulk/lot warehouse inventory including directed receipt, put-away and pick. Reject it where ERP/MES/WMS is authoritative unless commands/events are adapted through a defined ownership contract. Fiscal, Incoterm, QC, homologation, landed-cost, accounting-period and legal-title rules remain project-specific conditions.

## 3. Architecture contract

Procurement, factory and inventory are separate boundaries over one tenant-scoped PostgreSQL source of truth. Serial, bulk, warehouse and transfer activities share one inventory owner but retain explicit identities. State transitions use optimistic versions and allowed graphs; serial identities and item/lot/bin balances remain unique. Quantity reservation uses exact PostgreSQL numeric arithmetic. Receipt posts into RECEIVE before a capacity/ranking put-away. Picks reserve atomically from PICK/PUTPICK, exclude non-pickable/blocked/expired stock and optionally enforce FEFO; registration moves reservations to SHIP. Bin replenishment selects a fixed PICK/PUTPICK target below minimum, subtracts open inbound movement, plans only toward its explicit maximum, reserves eligible lower-rank sources and atomically registers or cancels the movement. A bulk transfer posts an exact registered pick as one durable shipment, exposes shipped-not-received cost components in transit and admits explicit partial receipts in separate versioned transitions. Every shipment/receipt has an idempotency key and posting identity; divergent replay, reused pick, quantity beyond the cumulative remainder and concurrent stale posting fail closed. Lot identity and every FIFO cost component cross the boundary without value invention; a specific-cost item must bind the transfer line to one exact source receipt and cannot silently fall back to FIFO. Each receipt creates its own destination put-away atomically. Both organizations are authorized at HTTP and matched again in SQL. Each accepted change and outbox event commits atomically. Unsupported costing/WMS behavior, invalid transitions, overselling, expired/blocked lots, stale writes and cross-organization access fail closed.

Cross-docking is explicit configuration, never an inferred optimization. A fixed, unblocked PICK/PUTPICK bin marked `cross_dock` receives only the minimum of receipt quantity, remaining released transfer demand inside the configured horizon and bin capacity. Existing allocations and exact generic pick quantities reduce demand; unallocated receipt quantity follows ordinary put-away. The allocation becomes pickable only after its put-away line is registered. Generic demand cannot consume cross-dock stock; the matching transfer consumes its allocation first, cannot exceed outstanding demand and links every reservation to its exact activity line. Cancellation releases allocation and balance together; registration converts reserved allocation into picked evidence atomically.

Item UOM configuration has one immutable base code with factor `1` and explicit base rounding precision. Alternate factors are positive, unique, immutable and exactly aligned with that precision. Conversion multiplies the requested quantity by the fixed factor in PostgreSQL numeric arithmetic and rejects any residual rather than rounding it. Item creation, base UOM and outbox evidence commit atomically. Physical quantity remains singular in `bulk_balance`; packaging rows are an exactly conserved composition. Explicit serializable `breakbulk` and `gather` commands preserve base quantity, reject residuals and concurrent double consumption, and persist immutable Take/Place plus outbox evidence. Ordinary quantity changes use base composition and reject implicit unpacking. Warehouse receipt accepts exactly one quantity representation: base or handling UOM/quantity. Put-away preserves that UOM only when every placement is exact; otherwise it rolls back unless automatic breakbulk is explicitly allowed, in which case Take/Place, receipt binding and outbox commit together. Open put-away reserves physical and composition quantities; registration preserves packages across bins o consume la composición base autorizada, y cancellation releases both reservations. Picking and replenishment resolve the requested target UOM first, prefer exact target composition and consider a larger source UOM only under explicit `allow_breakbulk`; each line stores From/To UOM quantities/factors, reserves physical and exact source composition together, records automatic Take/Place provenance on registration and preserves the remainder. Cancellation releases both layers. Transfer-receipt put-away follows the same base-UOM reservation contract. `customer-order` exige pedido/línea colocados, organización exacta y binding inmutable variante→item/UOM comercial; convierte cantidad vendida a base, resta picks abiertos o registrados y rechaza excedentes bajo lock serializable. El request exacto es idempotente y el replay divergente falla; cancelación libera pendiente y registro lo conserva como manejado. V173 postea cada pick registrado como un customer shipment inmutable, consume balance físico, composición UOM, reserva y FIFO cost en una transacción serializable, incrementa `shipped_quantity`, separa `fulfillment_state` del estado comercial y cierra replay divergente o doble despacho concurrente. Transporte, entrega/handover, facturación, pago/crédito, costo específico/serial customer shipment, service/production demand y working-day calendars permanecen owners y gates separados.

## 4. Exact file manifest

```text
CREATE internal/operations/service.go
CREATE internal/operations/service_test.go
CREATE internal/platform/postgres/operations.go
CREATE internal/platform/postgres/operations_integration_test.go
CREATE internal/platform/httpapi/operations.go
CREATE internal/platform/httpapi/operations_test.go
CREATE db/migrations/0025_serial_inventory_reservation_transfer.up.sql
CREATE db/migrations/0025_serial_inventory_reservation_transfer.down.sql
CREATE docs/inventory/MICROSOFT_BC_DERIVATION.md
CREATE internal/inventorycontrol/service.go
CREATE internal/inventorycontrol/service_test.go
CREATE internal/platform/postgres/inventorycontrol.go
CREATE internal/platform/postgres/inventorycontrol_integration_test.go
CREATE internal/platform/httpapi/inventorycontrol.go
CREATE internal/platform/httpapi/inventorycontrol_test.go
CREATE db/migrations/0026_bulk_lot_bin_costing.up.sql
CREATE db/migrations/0026_bulk_lot_bin_costing.down.sql
CREATE docs/inventory/MICROSOFT_BC_BULK_DERIVATION.md
CREATE internal/inventorycontrol/bulk.go
CREATE internal/inventorycontrol/bulk_test.go
CREATE internal/platform/postgres/bulk_inventory.go
CREATE internal/platform/postgres/bulk_inventory_integration_test.go
CREATE internal/platform/httpapi/bulk_inventory.go
CREATE internal/platform/httpapi/bulk_inventory_test.go
CREATE db/migrations/0027_operational_warehouse_fefo.up.sql
CREATE db/migrations/0027_operational_warehouse_fefo.down.sql
CREATE docs/inventory/MICROSOFT_BC_WAREHOUSE_DERIVATION.md
CREATE internal/inventorycontrol/warehouse.go
CREATE internal/inventorycontrol/warehouse_test.go
CREATE internal/platform/postgres/warehouse.go
CREATE internal/platform/postgres/warehouse_integration_test.go
CREATE internal/platform/httpapi/warehouse.go
CREATE internal/platform/httpapi/warehouse_test.go
CREATE db/migrations/0028_bulk_transfer_in_transit.up.sql
CREATE db/migrations/0028_bulk_transfer_in_transit.down.sql
CREATE db/migrations/0029_bulk_transfer_specific_cost.up.sql
CREATE db/migrations/0029_bulk_transfer_specific_cost.down.sql
CREATE db/migrations/0030_bulk_transfer_partial_posting.up.sql
CREATE db/migrations/0030_bulk_transfer_partial_posting.down.sql
CREATE docs/inventory/MICROSOFT_BC_TRANSFER_DERIVATION.md
CREATE internal/inventorycontrol/bulk_transfer.go
CREATE internal/inventorycontrol/bulk_transfer_test.go
CREATE internal/platform/postgres/bulk_transfer.go
CREATE internal/platform/postgres/bulk_transfer_integration_test.go
CREATE internal/platform/httpapi/bulk_transfer.go
CREATE internal/platform/httpapi/bulk_transfer_test.go
CREATE db/migrations/0031_warehouse_bin_replenishment.up.sql
CREATE db/migrations/0031_warehouse_bin_replenishment.down.sql
CREATE docs/inventory/MICROSOFT_BC_REPLENISHMENT_DERIVATION.md
CREATE internal/inventorycontrol/warehouse_replenishment.go
CREATE internal/platform/postgres/warehouse_replenishment.go
CREATE internal/platform/postgres/warehouse_replenishment_integration_test.go
CREATE db/migrations/0032_warehouse_transfer_crossdock.up.sql
CREATE db/migrations/0032_warehouse_transfer_crossdock.down.sql
CREATE docs/inventory/MICROSOFT_BC_CROSSDOCK_DERIVATION.md
CREATE internal/inventorycontrol/warehouse_crossdock.go
CREATE internal/platform/postgres/warehouse_crossdock.go
CREATE internal/platform/postgres/warehouse_crossdock_integration_test.go
CREATE db/migrations/0033_item_unit_of_measure.up.sql
CREATE db/migrations/0033_item_unit_of_measure.down.sql
CREATE docs/inventory/MICROSOFT_BC_ITEM_UOM_DERIVATION.md
CREATE internal/inventorycontrol/item_uom.go
CREATE internal/inventorycontrol/item_uom_test.go
CREATE internal/platform/postgres/item_uom.go
CREATE internal/platform/postgres/item_uom_integration_test.go
CREATE internal/platform/httpapi/item_uom.go
CREATE db/migrations/0034_bulk_uom_packaging.up.sql
CREATE db/migrations/0034_bulk_uom_packaging.down.sql
CREATE docs/inventory/MICROSOFT_BC_BREAKBULK_DERIVATION.md
CREATE internal/inventorycontrol/warehouse_breakbulk.go
CREATE internal/platform/postgres/warehouse_breakbulk.go
CREATE internal/platform/postgres/warehouse_breakbulk_integration_test.go
CREATE db/migrations/0035_warehouse_handling_uom.up.sql
CREATE db/migrations/0035_warehouse_handling_uom.down.sql
CREATE internal/platform/postgres/warehouse_packaging_flow_integration_test.go
CREATE db/migrations/0036_warehouse_pick_replenishment_uom.up.sql
CREATE db/migrations/0036_warehouse_pick_replenishment_uom.down.sql
CREATE internal/platform/postgres/warehouse_packaging.go
CREATE internal/platform/postgres/sales_warehouse.go
CREATE internal/platform/postgres/sales_warehouse_integration_test.go
CREATE db/migrations/0037_sales_order_warehouse_demand.up.sql
CREATE db/migrations/0037_sales_order_warehouse_demand.down.sql
CREATE docs/inventory/MICROSOFT_BC_SALES_WAREHOUSE_DERIVATION.md
CREATE internal/platform/postgres/sales_shipment.go
CREATE db/migrations/0038_customer_sales_shipment.up.sql
CREATE db/migrations/0038_customer_sales_shipment.down.sql
CREATE db/tests/0038_customer_sales_shipment.test.sql
CREATE docs/inventory/MICROSOFT_BC_CUSTOMER_SHIPMENT_DERIVATION.md
```

## 5. Materialization blocks

### FILE: `internal/operations/service.go`

```yaml
block_id: "GO-OPS-API:internal-operations-service-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "26ffcd62aad45d010c96806161ccf8fa9a4198bc9120e3294277c21cb58b2d80"
variables: []
secrets_allowed: false
```

````go
package operations

import (
	"context"
	"errors"
	"fmt"
	"regexp"
)

var ErrConflict = errors.New("operation conflict")
var code = regexp.MustCompile(`^[A-Z]{3}$`)

type PurchaseOrder struct {
	ID                        string `json:"id"`
	SupplierID                string `json:"supplier_id"`
	DestinationOrganizationID string `json:"destination_organization_id"`
	State                     string `json:"state"`
	Currency                  string `json:"currency"`
	TotalMinorUnits           int64  `json:"total_minor_units"`
	Version                   int64  `json:"version"`
}
type ProductionUnit struct {
	ID                  string `json:"id"`
	OrganizationID      string `json:"organization_id"`
	PurchaseOrderID     string `json:"purchase_order_id"`
	VariantID           string `json:"variant_id"`
	SerialNumber        string `json:"serial_number"`
	VIN                 string `json:"vin,omitempty"`
	BatterySerialNumber string `json:"battery_serial_number,omitempty"`
	State               string `json:"state"`
}
type StockUnit struct {
	ID                  string `json:"id"`
	OrganizationID      string `json:"organization_id"`
	VariantID           string `json:"variant_id"`
	ProductionUnitID    string `json:"production_unit_id,omitempty"`
	SerialNumber        string `json:"serial_number"`
	VIN                 string `json:"vin,omitempty"`
	BatterySerialNumber string `json:"battery_serial_number,omitempty"`
	State               string `json:"state"`
	Version             int64  `json:"version"`
}
type Repository interface {
	CreatePurchaseOrder(context.Context, string, string, PurchaseOrder) error
	TransitionPurchaseOrder(context.Context, string, string, string, string, int64, string, string) error
	CreateProductionUnit(context.Context, string, string, ProductionUnit) error
	TransitionProductionUnit(context.Context, string, string, string, string, string, string) error
	CreateStockUnit(context.Context, string, string, StockUnit) error
	TransitionStockUnit(context.Context, string, string, string, string, int64, string, string) error
}
type IDGenerator interface{ New() string }
type Service struct {
	repository Repository
	ids        IDGenerator
}

func NewService(r Repository, ids IDGenerator) *Service { return &Service{repository: r, ids: ids} }

var poTransitions = map[string]map[string]bool{"draft": {"submitted": true, "cancelled": true}, "submitted": {"accepted": true, "cancelled": true}, "accepted": {"in-production": true, "cancelled": true}, "in-production": {"shipped": true}, "shipped": {"received": true}}
var factoryTransitions = map[string]map[string]bool{"planned": {"assembly": true, "rejected": true}, "assembly": {"quality": true, "rejected": true}, "quality": {"released": true, "rejected": true}, "released": {"shipped": true}, "shipped": {"received": true}}
var stockTransitions = map[string]map[string]bool{"in-transit": {"available": true, "quarantine": true}, "available": {"reserved": true, "service": true, "quarantine": true}, "reserved": {"available": true, "sold": true}, "sold": {"service": true}, "service": {"available": true, "retired": true}, "quarantine": {"available": true, "retired": true}}

func (s *Service) CreatePurchaseOrder(ctx context.Context, tenant string, input PurchaseOrder) (PurchaseOrder, error) {
	if tenant == "" || input.SupplierID == "" || input.DestinationOrganizationID == "" || !code.MatchString(input.Currency) || input.TotalMinorUnits < 0 {
		return PurchaseOrder{}, fmt.Errorf("invalid purchase order")
	}
	input.ID = s.ids.New()
	input.State = "draft"
	input.Version = 1
	if err := s.repository.CreatePurchaseOrder(ctx, tenant, s.ids.New(), input); err != nil {
		return PurchaseOrder{}, err
	}
	return input, nil
}
func (s *Service) TransitionPurchaseOrder(ctx context.Context, tenant, organization, id, current, target string, version int64) error {
	if organization == "" || !poTransitions[current][target] || version < 1 {
		return fmt.Errorf("%w: invalid purchase order transition", ErrConflict)
	}
	return s.repository.TransitionPurchaseOrder(ctx, tenant, organization, id, current, version, target, s.ids.New())
}
func (s *Service) RegisterProductionUnit(ctx context.Context, tenant string, input ProductionUnit) (ProductionUnit, error) {
	if tenant == "" || input.OrganizationID == "" || input.PurchaseOrderID == "" || input.VariantID == "" || input.SerialNumber == "" {
		return ProductionUnit{}, fmt.Errorf("invalid production unit")
	}
	input.ID = s.ids.New()
	input.State = "planned"
	if err := s.repository.CreateProductionUnit(ctx, tenant, s.ids.New(), input); err != nil {
		return ProductionUnit{}, err
	}
	return input, nil
}
func (s *Service) TransitionProductionUnit(ctx context.Context, tenant, organization, id, current, target string) error {
	if organization == "" || !factoryTransitions[current][target] {
		return fmt.Errorf("%w: invalid production transition", ErrConflict)
	}
	return s.repository.TransitionProductionUnit(ctx, tenant, organization, id, current, target, s.ids.New())
}
func (s *Service) ReceiveStockUnit(ctx context.Context, tenant string, input StockUnit) (StockUnit, error) {
	if tenant == "" || input.OrganizationID == "" || input.VariantID == "" || input.SerialNumber == "" {
		return StockUnit{}, fmt.Errorf("invalid stock unit")
	}
	input.ID = s.ids.New()
	input.State = "in-transit"
	input.Version = 1
	if err := s.repository.CreateStockUnit(ctx, tenant, s.ids.New(), input); err != nil {
		return StockUnit{}, err
	}
	return input, nil
}
func (s *Service) TransitionStockUnit(ctx context.Context, tenant, organization, id, current, target string, version int64) error {
	if organization == "" || !stockTransitions[current][target] || version < 1 {
		return fmt.Errorf("%w: invalid stock transition", ErrConflict)
	}
	return s.repository.TransitionStockUnit(ctx, tenant, organization, id, current, version, target, s.ids.New())
}
````

### FILE: `internal/operations/service_test.go`

```yaml
block_id: "GO-OPS-API:internal-operations-service_test-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "9b988a6646d47c00f348a7104f969cf8e0bfb4728971be20eefffeca3ce22e53"
variables: []
secrets_allowed: false
```

````go
package operations

import (
	"context"
	"testing"
)

type fakeRepo struct{ po, factory, stock int }

func (f *fakeRepo) CreatePurchaseOrder(context.Context, string, string, PurchaseOrder) error {
	f.po++
	return nil
}
func (f *fakeRepo) TransitionPurchaseOrder(context.Context, string, string, string, string, int64, string, string) error {
	f.po++
	return nil
}
func (f *fakeRepo) CreateProductionUnit(context.Context, string, string, ProductionUnit) error {
	f.factory++
	return nil
}
func (f *fakeRepo) TransitionProductionUnit(context.Context, string, string, string, string, string, string) error {
	f.factory++
	return nil
}
func (f *fakeRepo) CreateStockUnit(context.Context, string, string, StockUnit) error {
	f.stock++
	return nil
}
func (f *fakeRepo) TransitionStockUnit(context.Context, string, string, string, string, int64, string, string) error {
	f.stock++
	return nil
}

type ids struct{ n int }

func (i *ids) New() string {
	i.n++
	return []string{"018f4d4a-7b36-7a21-8d10-2f4c54c28001", "018f4d4a-7b36-7a21-8d10-2f4c54c28002", "018f4d4a-7b36-7a21-8d10-2f4c54c28003", "018f4d4a-7b36-7a21-8d10-2f4c54c28004", "018f4d4a-7b36-7a21-8d10-2f4c54c28005", "018f4d4a-7b36-7a21-8d10-2f4c54c28006"}[i.n-1]
}
func TestStateMachinesRejectSkippedTransitions(t *testing.T) {
	repo := &fakeRepo{}
	service := NewService(repo, &ids{})
	ctx := context.Background()
	if _, err := service.CreatePurchaseOrder(ctx, "tenant", PurchaseOrder{SupplierID: "s", DestinationOrganizationID: "o", Currency: "USD", TotalMinorUnits: 1}); err != nil {
		t.Fatal(err)
	}
	if err := service.TransitionPurchaseOrder(ctx, "tenant", "o", "po", "draft", "received", 1); err == nil {
		t.Fatal("purchase order skipped states")
	}
	if err := service.TransitionProductionUnit(ctx, "tenant", "o", "u", "planned", "released"); err == nil {
		t.Fatal("factory skipped states")
	}
	if err := service.TransitionStockUnit(ctx, "tenant", "o", "s", "available", "sold", 1); err == nil {
		t.Fatal("stock skipped reservation")
	}
	if repo.po != 1 || repo.factory != 0 || repo.stock != 0 {
		t.Fatalf("unexpected repository calls %+v", repo)
	}
}
````

### FILE: `internal/platform/postgres/operations.go`

```yaml
block_id: "GO-OPS-API:internal-platform-postgres-operations-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "44245e35f8244a99f269aa64ef2d5bd988d29dd50ed96a8a7a0bba793ac2b17b"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Operations struct{ pool *pgxpool.Pool }

func NewOperations(pool *pgxpool.Pool) *Operations { return &Operations{pool: pool} }
func outboxPayload(value any) json.RawMessage      { payload, _ := json.Marshal(value); return payload }
func (r *Operations) CreatePurchaseOrder(ctx context.Context, tenant, eventID string, value operations.PurchaseOrder) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into procurement.purchase_order(tenant_id,purchase_order_id,supplier_id,destination_organization_id,state,currency,total_minor_units,version)values($1,$2,$3,$4,$5,$6,$7,$8)`, tenant, value.ID, value.SupplierID, value.DestinationOrganizationID, value.State, value.Currency, value.TotalMinorUnits, value.Version)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'purchase-order',$3,$4,'purchase-order.created',1,clock_timestamp(),$5)`, tenant, eventID, value.ID, value.Version, outboxPayload(value))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Operations) TransitionPurchaseOrder(ctx context.Context, tenant, organization, id, current string, version int64, target, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update procurement.purchase_order set state=$6,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and destination_organization_id=$2 and purchase_order_id=$3 and state=$4 and version=$5`, tenant, organization, id, current, version, target)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return operations.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'purchase-order',$3,$4,$5,1,clock_timestamp(),$6)`, tenant, eventID, id, version+1, "purchase-order."+target, outboxPayload(map[string]string{"state": target}))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Operations) CreateProductionUnit(ctx context.Context, tenant, eventID string, value operations.ProductionUnit) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into factory.production_unit(tenant_id,production_unit_id,purchase_order_id,variant_id,serial_number,vin,battery_serial_number,state) select $1,$2,p.purchase_order_id,$4,$5,nullif($6,''),nullif($7,''),$8 from procurement.purchase_order p where p.tenant_id=$1 and p.purchase_order_id=$3 and p.destination_organization_id=$9`, tenant, value.ID, value.PurchaseOrderID, value.VariantID, value.SerialNumber, value.VIN, value.BatterySerialNumber, value.State, value.OrganizationID)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return operations.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'production-unit',$3,1,'production-unit.registered',1,clock_timestamp(),$4)`, tenant, eventID, value.ID, outboxPayload(value))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Operations) TransitionProductionUnit(ctx context.Context, tenant, organization, id, current, target, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update factory.production_unit u set state=$5,updated_at=clock_timestamp() from procurement.purchase_order p where u.tenant_id=$1 and u.production_unit_id=$3 and u.state=$4 and p.tenant_id=u.tenant_id and p.purchase_order_id=u.purchase_order_id and p.destination_organization_id=$2`, tenant, organization, id, current, target)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return operations.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'production-unit',$3,extract(epoch from clock_timestamp())::bigint,$4,1,clock_timestamp(),$5)`, tenant, eventID, id, "production-unit."+target, outboxPayload(map[string]string{"state": target}))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Operations) CreateStockUnit(ctx context.Context, tenant, eventID string, value operations.StockUnit) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,production_unit_id,serial_number,vin,battery_serial_number,state,version)values($1,$2,$3,$4,nullif($5,''),$6,nullif($7,''),nullif($8,''),$9,$10)`, tenant, value.ID, value.OrganizationID, value.VariantID, value.ProductionUnitID, value.SerialNumber, value.VIN, value.BatterySerialNumber, value.State, value.Version)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'stock-unit',$3,$4,'stock-unit.received-in-transit',1,clock_timestamp(),$5)`, tenant, eventID, value.ID, value.Version, outboxPayload(value))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Operations) TransitionStockUnit(ctx context.Context, tenant, organization, id, current string, version int64, target, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update inventory.stock_unit set state=$6,version=version+1,received_at=case when $6 in ('available','quarantine') then coalesce(received_at,clock_timestamp()) else received_at end,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and state=$4 and version=$5`, tenant, organization, id, current, version, target)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return operations.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'stock-unit',$3,$4,$5,1,clock_timestamp(),$6)`, tenant, eventID, id, version+1, "stock-unit."+target, outboxPayload(map[string]string{"state": target}))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
````

### FILE: `internal/platform/postgres/operations_integration_test.go`

```yaml
block_id: "GO-OPS-API:internal-platform-postgres-operations_integration_test-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "f64f3fc4fa648ceb7d94e69fa0e35a3685baf8184b0f9264ad4ae4ef9e5149cb"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
)

func TestOperationsFlowAndConcurrency(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28101"
	cleanup := func() {
		for _, q := range []string{`delete from platform.outbox_event where tenant_id=$1`, `delete from inventory.stock_unit where tenant_id=$1`, `delete from factory.production_unit where tenant_id=$1`, `delete from procurement.purchase_order where tenant_id=$1`, `delete from partner.supplier where tenant_id=$1`, `delete from catalog.vehicle_variant where tenant_id=$1`, `delete from catalog.vehicle_model where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`} {
			_, _ = pool.Exec(ctx, q, tenant)
		}
	}
	cleanup()
	defer cleanup()
	fixtures := []struct {
		q    string
		args []any
	}{{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'ops-api','Ops','Ops')`, []any{tenant}}, {`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'org','central','Central','warehouse')`, []any{tenant}}, {`insert into partner.supplier(tenant_id,supplier_id,supplier_code,legal_name,status)values($1,'supplier','supplier','Supplier','active')`, []any{tenant}}, {`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Model','bicycle','active')`, []any{tenant}}, {`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Variant','{}','active')`, []any{tenant}}}
	for _, f := range fixtures {
		if _, err := pool.Exec(ctx, f.q, f.args...); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewOperations(pool)
	po := operations.PurchaseOrder{ID: "po", SupplierID: "supplier", DestinationOrganizationID: "org", State: "draft", Currency: "USD", TotalMinorUnits: 100, Version: 1}
	if err := repo.CreatePurchaseOrder(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28102", po); err != nil {
		t.Fatal(err)
	}
	if err := repo.TransitionPurchaseOrder(ctx, tenant, "org-other", "po", "draft", 1, "submitted", "018f4d4a-7b36-7a21-8d10-2f4c54c28110"); err == nil {
		t.Fatal("purchase order crossed organization scope")
	}
	if err := repo.TransitionPurchaseOrder(ctx, tenant, "org", "po", "draft", 1, "submitted", "018f4d4a-7b36-7a21-8d10-2f4c54c28106"); err != nil {
		t.Fatal(err)
	}
	if err := repo.TransitionPurchaseOrder(ctx, tenant, "org", "po", "draft", 1, "cancelled", "018f4d4a-7b36-7a21-8d10-2f4c54c28107"); err == nil {
		t.Fatal("stale purchase order transition succeeded")
	}
	unit := operations.ProductionUnit{ID: "unit", OrganizationID: "org", PurchaseOrderID: "po", VariantID: "variant", SerialNumber: "SERIAL", State: "planned"}
	if err := repo.CreateProductionUnit(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28103", unit); err != nil {
		t.Fatal(err)
	}
	if err := repo.TransitionProductionUnit(ctx, tenant, "org-other", "unit", "planned", "assembly", "018f4d4a-7b36-7a21-8d10-2f4c54c28111"); err == nil {
		t.Fatal("production unit crossed organization scope")
	}
	if err := repo.TransitionProductionUnit(ctx, tenant, "org", "unit", "planned", "assembly", "018f4d4a-7b36-7a21-8d10-2f4c54c28104"); err != nil {
		t.Fatal(err)
	}
	stock := operations.StockUnit{ID: "stock", OrganizationID: "org", VariantID: "variant", ProductionUnitID: "unit", SerialNumber: "SERIAL", State: "in-transit", Version: 1}
	if err := repo.CreateStockUnit(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28105", stock); err != nil {
		t.Fatal(err)
	}
	if err := repo.TransitionStockUnit(ctx, tenant, "org-other", "stock", "in-transit", 1, "available", "018f4d4a-7b36-7a21-8d10-2f4c54c28112"); err == nil {
		t.Fatal("stock unit crossed organization scope")
	}
	if err := repo.TransitionStockUnit(ctx, tenant, "org", "stock", "in-transit", 1, "available", "018f4d4a-7b36-7a21-8d10-2f4c54c28108"); err != nil {
		t.Fatal(err)
	}
	if err := repo.TransitionStockUnit(ctx, tenant, "org", "stock", "in-transit", 1, "quarantine", "018f4d4a-7b36-7a21-8d10-2f4c54c28109"); err == nil {
		t.Fatal("stale stock transition succeeded")
	}
	var events int
	if err := pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1`, tenant).Scan(&events); err != nil || events != 6 {
		t.Fatalf("events=%d err=%v", events, err)
	}
}
````

### FILE: `internal/platform/httpapi/operations.go`

```yaml
block_id: "GO-OPS-API:internal-platform-httpapi-operations-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "f7d8fa60640549c31d74f403c2886f19cad7a428d4b62a12e1be261dc4bbd4a7"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"net/http"
)

type OperationsModule struct{ Service *operations.Service }

func (m OperationsModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	api := operationsAPI{service: m.Service, verifier: verifier}
	mux.HandleFunc("POST /v1/procurement/purchase-orders", api.createPO)
	mux.HandleFunc("POST /v1/procurement/purchase-orders/{id}/transitions", api.transitionPO)
	mux.HandleFunc("POST /v1/factory/units", api.createUnit)
	mux.HandleFunc("POST /v1/factory/units/{id}/transitions", api.transitionUnit)
	mux.HandleFunc("POST /v1/inventory/stock", api.createStock)
	mux.HandleFunc("POST /v1/inventory/stock/{id}/transitions", api.transitionStock)
}

type operationsAPI struct {
	service  *operations.Service
	verifier identity.Verifier
}

func (a operationsAPI) principal(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return p, false
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return p, false
	}
	return p, true
}
func (a operationsAPI) createPO(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "procurement:write")
	if !ok {
		return
	}
	var input operations.PurchaseOrder
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.DestinationOrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.CreatePurchaseOrder(r.Context(), p.TenantID, input)
	if err != nil {
		writeProblem(w, 400, "INVALID_PURCHASE_ORDER", "purchase order does not match contract")
		return
	}
	writeJSON(w, 201, value)
}
func (a operationsAPI) transitionPO(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "procurement:write")
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		Current        string `json:"current"`
		Target         string `json:"target"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeTransition(w, a.service.TransitionPurchaseOrder(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Current, input.Target, input.Version))
}
func (a operationsAPI) createUnit(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "factory:write")
	if !ok {
		return
	}
	var input operations.ProductionUnit
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.RegisterProductionUnit(r.Context(), p.TenantID, input)
	if err != nil {
		writeProblem(w, 400, "INVALID_PRODUCTION_UNIT", "unit does not match contract")
		return
	}
	writeJSON(w, 201, value)
}
func (a operationsAPI) transitionUnit(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "factory:write")
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		Current        string `json:"current"`
		Target         string `json:"target"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeTransition(w, a.service.TransitionProductionUnit(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Current, input.Target))
}
func (a operationsAPI) createStock(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "inventory:write")
	if !ok {
		return
	}
	var input operations.StockUnit
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.ReceiveStockUnit(r.Context(), p.TenantID, input)
	if err != nil {
		writeProblem(w, 400, "INVALID_STOCK_UNIT", "stock unit does not match contract")
		return
	}
	writeJSON(w, 201, value)
}
func (a operationsAPI) transitionStock(w http.ResponseWriter, r *http.Request) {
	p, ok := a.principal(w, r, "inventory:write")
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		Current        string `json:"current"`
		Target         string `json:"target"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeTransition(w, a.service.TransitionStockUnit(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Current, input.Target, input.Version))
}
func writeTransition(w http.ResponseWriter, err error) {
	if errors.Is(err, operations.ErrConflict) {
		writeProblem(w, 409, "TRANSITION_CONFLICT", "state or version conflict")
		return
	}
	if err != nil {
		writeProblem(w, 500, "INTERNAL_ERROR", "transition failed")
		return
	}
	writeJSON(w, 200, map[string]string{"status": "accepted"})
}
````

### FILE: `internal/platform/httpapi/operations_test.go`

```yaml
block_id: "GO-OPS-API:internal-platform-httpapi-operations_test-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "17b22c111231bc8a315fa799e2b61a1495ba5b93c63f0a2cf45ce5ab74ea2c47"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/identity"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type opsRepo struct{ created int }

func (o *opsRepo) CreatePurchaseOrder(context.Context, string, string, operations.PurchaseOrder) error {
	o.created++
	return nil
}
func (o *opsRepo) TransitionPurchaseOrder(context.Context, string, string, string, string, int64, string, string) error {
	return operations.ErrConflict
}
func (o *opsRepo) CreateProductionUnit(context.Context, string, string, operations.ProductionUnit) error {
	return nil
}
func (o *opsRepo) TransitionProductionUnit(context.Context, string, string, string, string, string, string) error {
	return nil
}
func (o *opsRepo) CreateStockUnit(context.Context, string, string, operations.StockUnit) error {
	return nil
}
func (o *opsRepo) TransitionStockUnit(context.Context, string, string, string, string, int64, string, string) error {
	return nil
}

type opsIDs struct{ n int }

func (i *opsIDs) New() string {
	i.n++
	return []string{"018f4d4a-7b36-7a21-8d10-2f4c54c28201", "018f4d4a-7b36-7a21-8d10-2f4c54c28202", "018f4d4a-7b36-7a21-8d10-2f4c54c28203"}[i.n-1]
}

type opsVerifier struct{}

func (opsVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{Subject: "operator", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c28200", Permissions: map[string]struct{}{"procurement:write": {}}, Organizations: map[string]struct{}{"o": {}}}, nil
}
func TestOperationsHTTPAuthorizationAndConflict(t *testing.T) {
	repo := &opsRepo{}
	service := operations.NewService(repo, &opsIDs{})
	module := OperationsModule{Service: service}
	mux := http.NewServeMux()
	module.Register(mux, opsVerifier{})
	request := httptest.NewRequest("POST", "/v1/procurement/purchase-orders", strings.NewReader(`{"supplier_id":"s","destination_organization_id":"o","currency":"USD","total_minor_units":100}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.created != 1 {
		t.Fatalf("create status=%d calls=%d body=%s", response.Code, repo.created, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/procurement/purchase-orders/po/transitions", strings.NewReader(`{"organization_id":"o","current":"draft","target":"submitted","version":1}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 409 {
		t.Fatalf("conflict status=%d body=%s", response.Code, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/inventory/stock", strings.NewReader(`{}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 {
		t.Fatalf("permission status=%d", response.Code)
	}
}
````


### FILE: `db/migrations/0025_serial_inventory_reservation_transfer.up.sql`

```yaml
block_id: "GO-OPS-API:serial-inventory-migration-up:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps 2eae56d: AvailabletoPromise, ReservationEngineMgt, CreateReservEntry, TransferLineReserve and SCMTransferReservation; portable PostgreSQL scope"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "f65ce0b66c25360985c21f8d99f918b712f0f384af694644e19c7756500e0079"
variables: []
secrets_allowed: false
```

````sql
create table inventory.serial_reservation (
  tenant_id uuid not null,
  reservation_id text not null,
  organization_id text not null,
  stock_unit_id text not null,
  variant_id text not null,
  demand_kind text not null check (demand_kind in ('customer-order','service','transfer-outbound','manual')),
  demand_id text not null,
  demand_line_id text not null default '',
  status text not null check (status in ('reservation','released','consumed')),
  cancellation_disallowed boolean not null default false,
  expires_at timestamptz,
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,reservation_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,stock_unit_id) references inventory.stock_unit(tenant_id,stock_unit_id),
  foreign key (tenant_id,variant_id) references catalog.vehicle_variant(tenant_id,variant_id),
  check (expires_at is null or expires_at > created_at)
);

create unique index serial_reservation_active_stock_uidx
  on inventory.serial_reservation(tenant_id,stock_unit_id)
  where status='reservation';

create unique index serial_reservation_active_demand_uidx
  on inventory.serial_reservation(tenant_id,demand_kind,demand_id,demand_line_id)
  where status='reservation';

create index serial_reservation_scope_idx
  on inventory.serial_reservation(tenant_id,organization_id,variant_id,status,expires_at);

create table inventory.serial_transfer (
  tenant_id uuid not null,
  transfer_id text not null,
  from_organization_id text not null,
  to_organization_id text not null,
  expected_receipt_at timestamptz not null,
  state text not null check (state in ('draft','released','in-transit','received','cancelled')),
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,transfer_id),
  foreign key (tenant_id,from_organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,to_organization_id) references org.organization(tenant_id,organization_id),
  check (from_organization_id<>to_organization_id)
);

create table inventory.serial_transfer_unit (
  tenant_id uuid not null,
  transfer_id text not null,
  stock_unit_id text not null,
  variant_id text not null,
  primary key (tenant_id,transfer_id,stock_unit_id),
  foreign key (tenant_id,transfer_id) references inventory.serial_transfer(tenant_id,transfer_id) on delete restrict,
  foreign key (tenant_id,stock_unit_id) references inventory.stock_unit(tenant_id,stock_unit_id),
  foreign key (tenant_id,variant_id) references catalog.vehicle_variant(tenant_id,variant_id)
);

create index serial_transfer_inbound_idx
  on inventory.serial_transfer(tenant_id,to_organization_id,state,expected_receipt_at);

create function inventory.serial_atp(
  p_tenant uuid,
  p_organization text,
  p_variant text,
  p_horizon timestamptz
) returns table(
  available_inventory bigint,
  scheduled_receipt bigint,
  gross_requirement bigint,
  available_to_promise bigint
) language sql stable as $$
with available as (
  select count(*)::bigint quantity
  from inventory.stock_unit
  where tenant_id=p_tenant and organization_id=p_organization and variant_id=p_variant and state='available'
), inbound as (
  select count(*)::bigint quantity
  from inventory.serial_transfer t
  join inventory.serial_transfer_unit u using(tenant_id,transfer_id)
  where t.tenant_id=p_tenant and t.to_organization_id=p_organization and u.variant_id=p_variant
    and t.state='in-transit' and t.expected_receipt_at<=p_horizon
), demand as (
  select coalesce(sum(l.quantity-case when l.allocated_stock_unit_id is null then 0 else 1 end),0)::bigint quantity
  from sales.customer_order o
  join sales.customer_order_line l using(tenant_id,order_id)
  where o.tenant_id=p_tenant and o.organization_id=p_organization and l.variant_id=p_variant
    and o.state in ('placed','confirmed')
)
select a.quantity,i.quantity,d.quantity,greatest(0::bigint,a.quantity+i.quantity-d.quantity)
from available a cross join inbound i cross join demand d;
$$;

comment on function inventory.serial_atp(uuid,text,text,timestamptz) is
  'Portable serial-stock ATP adaptation governed by Microsoft BC Available to Promise: available inventory plus due inbound transfer receipts minus unallocated placed demand. It is not full Business Central planning/CTP.';
````

### FILE: `db/migrations/0025_serial_inventory_reservation_transfer.down.sql`

```yaml
block_id: "GO-OPS-API:serial-inventory-migration-down:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps 2eae56d portable rollback counterpart; local schema lifecycle"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "4329814cef1aa29e7ed1c4fd6ab5ad50eb825a740744eb76132ef3116d5c2cb3"
variables: []
secrets_allowed: false
```

````sql
drop function if exists inventory.serial_atp(uuid,text,text,timestamptz);
drop table if exists inventory.serial_transfer_unit;
drop table if exists inventory.serial_transfer;
drop table if exists inventory.serial_reservation;
````

### FILE: `docs/inventory/MICROSOFT_BC_DERIVATION.md`

```yaml
block_id: "GO-OPS-API:microsoft-bc-derivation:v1"
operation: CREATE
provenance: AUTHORED
source: "Pinned Microsoft BCApps source and Microsoft Learn authorities enumerated in the document"
license: "LicenseRef-Workspace-Owner"
sha256: "6151c41c6db1eff26458bc5c6d33e623c900ca6839461d752212cdbb6442df57"
variables: []
secrets_allowed: false
```

````markdown
# Portable serial inventory — Microsoft Business Central derivation record

## Authority and license

The portable module is an **ADAPTED** implementation, not verbatim Microsoft Go code and not a Business Central runtime. Its authority is Microsoft `BCApps` at GitHub-verified commit `2eae56d704a1fd035d104f333602aea7091b7749` (tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`), licensed under MIT by the root `License.txt`.

| Official Microsoft source | Git blob | Bytes | SHA-256 |
|---|---:|---:|---|
| `src/Layers/W1/BaseApp/Inventory/Availability/AvailabletoPromise.Codeunit.al` | `440b6ed4dfc6c9297dd3115d9fee682aab6aaf17` | 36,958 | `76d582db57a228559e2458920e11f70776597af25fd976863dbbccac1079d29f` |
| `src/Layers/W1/BaseApp/Inventory/Tracking/ReservationEngineMgt.Codeunit.al` | `b4943be474b5c2023ecc1c5590ae0ce4e21e6158` | 57,538 | `7b0b4c4d1f58b74587e7582a763208342db6376684371f0e6b6448a5f78358db` |
| `src/Layers/W1/BaseApp/Inventory/Tracking/CreateReservEntry.Codeunit.al` | `9fe395c7ac00b36ea628eaac98b4af2b839a9529` | 64,760 | `363c0aefc1350867cf3891ae34ff46967905f2493774492af7a481da35836998` |
| `src/Layers/W1/BaseApp/Inventory/Transfer/TransferLineReserve.Codeunit.al` | `0fd198a2c51219ee76eacee26c9c76bde5b58483` | 65,044 | `81f33a5d8cfe079420bbfde13f4fb8b589422fc4e9ea010d4ab409b301974a9e` |
| `src/Layers/W1/BaseApp/Warehouse/Availability/WarehouseAvailabilityMgt.Codeunit.al` | `946f3604a69e9bbaeab514a2d1fab42f2dc01e6f` | 62,348 | `1744161a3d50d013c0ca07f1f4e2e0dd657d81b4fab1f733cf9867325a04c096` |
| `src/Layers/W1/Tests/SCM-Reservation/SCMTransferReservation.Codeunit.al` | `ec88292b6cfe26db06e19f93598fc3cf48db121d` | 161,402 | `b41c605cd18bca08bdafae7680f04434c14853667be2fb44f05f7e390a846607` |

## Narrow rules transferred

1. ATP follows the Microsoft shape: available inventory plus scheduled receipts minus gross requirement, while avoiding subtraction of demand already bound by reservation. The portable function applies this only to serialized sellable units, due inbound transfers and unallocated placed/confirmed order lines.
2. A reservation binds one supply unit to one demand and cannot be silently duplicated. Cancellation is versioned; customer-order allocations are cancellation-disallowed through the generic endpoint.
3. Transfer reservation is direction-aware. Draft creation reserves outbound units, shipment moves them to `in-transit`, only the destination organization can receive them, and receipt changes ownership before returning the unit to `available`.
4. Every mutation is tenant- and organization-scoped, optimistic-versioned, transactional with its outbox event and tested for conflict/replay boundaries.

## Explicit non-claims

This serial-inventory boundary does not itself implement full Business Central planning, CTP, warehouse picks/put-aways, replenishment, calendars, substitutions, costing adjustment or AL extension events. The composed pack's separate `MICROSOFT_BC_BULK_DERIVATION.md` governs the admitted lot/bin/FIFO-specific boundary without changing these serial rules. Neither boundary claims Microsoft authored or supports the Go/PostgreSQL translation.

## Executed regression mapping

- concurrent attempts to reserve one stock unit produce exactly one winner;
- customer allocation creates the same canonical reservation row;
- ATP reports inbound stock only after shipment and within the requested horizon;
- source organization cannot receive its own outbound transfer;
- receipt changes organization and permits a later reverse transfer;
- all migrations 0001–0025, full Go tests, vet and both API builds execute with the pinned toolchains.
````

### FILE: `internal/inventorycontrol/service.go`

```yaml
block_id: "GO-OPS-API:inventory-control-service:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps 2eae56d reservation/ATP/transfer invariants; local Go port"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "e8fbdbb495604a4e56a4945b492847a429550d0b78e6f6513231a43b3c3dd60f"
variables: []
secrets_allowed: false
```

````go
// Package inventorycontrol is a narrow Go/PostgreSQL adaptation of the reservation,
// ATP and transfer invariants documented in docs/inventory/MICROSOFT_BC_DERIVATION.md.
// It is not verbatim Microsoft code or a Business Central replacement.
package inventorycontrol

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrConflict = errors.New("inventory control conflict")

type ATP struct {
	OrganizationID     string    `json:"organization_id"`
	VariantID          string    `json:"variant_id"`
	Horizon            time.Time `json:"horizon"`
	AvailableInventory int64     `json:"available_inventory"`
	ScheduledReceipt   int64     `json:"scheduled_receipt"`
	GrossRequirement   int64     `json:"gross_requirement"`
	AvailableToPromise int64     `json:"available_to_promise"`
}

type Reservation struct {
	ID                     string     `json:"id"`
	OrganizationID         string     `json:"organization_id"`
	StockUnitID            string     `json:"stock_unit_id"`
	VariantID              string     `json:"variant_id"`
	DemandKind             string     `json:"demand_kind"`
	DemandID               string     `json:"demand_id"`
	DemandLineID           string     `json:"demand_line_id"`
	Status                 string     `json:"status"`
	CancellationDisallowed bool       `json:"cancellation_disallowed"`
	ExpiresAt              *time.Time `json:"expires_at,omitempty"`
	Version                int64      `json:"version"`
}

type Transfer struct {
	ID                 string    `json:"id"`
	FromOrganizationID string    `json:"from_organization_id"`
	ToOrganizationID   string    `json:"to_organization_id"`
	StockUnitIDs       []string  `json:"stock_unit_ids"`
	ExpectedReceiptAt  time.Time `json:"expected_receipt_at"`
	State              string    `json:"state"`
	Version            int64     `json:"version"`
}

type Repository interface {
	AvailableToPromise(context.Context, string, string, string, time.Time) (ATP, error)
	Reserve(context.Context, string, string, string, Reservation, int64) (Reservation, error)
	Release(context.Context, string, string, string, int64, string) error
	CreateTransfer(context.Context, string, string, Transfer, map[string]int64, string) (Transfer, error)
	TransitionTransfer(context.Context, string, string, string, string, int64, string, string) error
}

type IDGenerator interface{ New() string }

type Service struct {
	repository Repository
	ids        IDGenerator
}

func NewService(repository Repository, ids IDGenerator) *Service {
	return &Service{repository: repository, ids: ids}
}

func (s *Service) AvailableToPromise(ctx context.Context, tenant, organization, variant string, horizon time.Time) (ATP, error) {
	if tenant == "" || organization == "" || variant == "" || horizon.IsZero() {
		return ATP{}, fmt.Errorf("invalid ATP query")
	}
	return s.repository.AvailableToPromise(ctx, tenant, organization, variant, horizon.UTC())
}

func (s *Service) Reserve(ctx context.Context, tenant string, value Reservation, stockVersion int64) (Reservation, error) {
	allowed := map[string]bool{"service": true, "manual": true}
	if tenant == "" || value.OrganizationID == "" || value.StockUnitID == "" || value.VariantID == "" || !allowed[value.DemandKind] || value.DemandID == "" || stockVersion < 1 {
		return Reservation{}, fmt.Errorf("invalid reservation")
	}
	value.ID, value.Status, value.Version = s.ids.New(), "reservation", 1
	return s.repository.Reserve(ctx, tenant, s.ids.New(), value.ID, value, stockVersion)
}

func (s *Service) Release(ctx context.Context, tenant, organization, reservationID string, version int64) error {
	if tenant == "" || organization == "" || reservationID == "" || version < 1 {
		return fmt.Errorf("invalid reservation release")
	}
	return s.repository.Release(ctx, tenant, organization, reservationID, version, s.ids.New())
}

func (s *Service) CreateTransfer(ctx context.Context, tenant string, value Transfer, stockVersions map[string]int64) (Transfer, error) {
	if tenant == "" || value.FromOrganizationID == "" || value.ToOrganizationID == "" || value.FromOrganizationID == value.ToOrganizationID || value.ExpectedReceiptAt.IsZero() || len(value.StockUnitIDs) == 0 || len(value.StockUnitIDs) != len(stockVersions) {
		return Transfer{}, fmt.Errorf("invalid transfer")
	}
	seen := map[string]bool{}
	for _, stockID := range value.StockUnitIDs {
		if stockID == "" || seen[stockID] || stockVersions[stockID] < 1 {
			return Transfer{}, fmt.Errorf("invalid transfer units")
		}
		seen[stockID] = true
	}
	value.ID, value.State, value.Version = s.ids.New(), "draft", 1
	value.ExpectedReceiptAt = value.ExpectedReceiptAt.UTC()
	return s.repository.CreateTransfer(ctx, tenant, value.ID, value, stockVersions, s.ids.New())
}

var transferTransitions = map[string]map[string]bool{
	"draft":      {"released": true, "cancelled": true},
	"released":   {"in-transit": true, "cancelled": true},
	"in-transit": {"received": true},
}

func (s *Service) TransitionTransfer(ctx context.Context, tenant, authorizedOrganization, transferID, current, target string, version int64) error {
	if tenant == "" || authorizedOrganization == "" || transferID == "" || !transferTransitions[current][target] || version < 1 {
		return fmt.Errorf("%w: invalid transfer transition", ErrConflict)
	}
	return s.repository.TransitionTransfer(ctx, tenant, authorizedOrganization, transferID, current, version, target, s.ids.New())
}
````

### FILE: `internal/inventorycontrol/service_test.go`

```yaml
block_id: "GO-OPS-API:inventory-control-service-test:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps 2eae56d SCM-Reservation test surface; local Go contract tests"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "8ea59f5b31e8ec453383b08fd8ae4a0b7f5808a3e6dfd31aceb75b03bbe87183"
variables: []
secrets_allowed: false
```

````go
package inventorycontrol

import (
	"context"
	"testing"
	"time"
)

type fakeRepo struct{ reserve, release, create, transition int }

func (f *fakeRepo) AvailableToPromise(context.Context, string, string, string, time.Time) (ATP, error) {
	return ATP{AvailableToPromise: 2}, nil
}
func (f *fakeRepo) Reserve(_ context.Context, _, _, _ string, value Reservation, _ int64) (Reservation, error) {
	f.reserve++
	return value, nil
}
func (f *fakeRepo) Release(context.Context, string, string, string, int64, string) error {
	f.release++
	return nil
}
func (f *fakeRepo) CreateTransfer(_ context.Context, _, _ string, value Transfer, _ map[string]int64, _ string) (Transfer, error) {
	f.create++
	return value, nil
}
func (f *fakeRepo) TransitionTransfer(context.Context, string, string, string, string, int64, string, string) error {
	f.transition++
	return nil
}

type ids struct{ n int }

func (i *ids) New() string { i.n++; return "id" }

func TestServiceRejectsInvalidReservationAndTransfer(t *testing.T) {
	repo, generator := &fakeRepo{}, &ids{}
	service := NewService(repo, generator)
	if _, err := service.Reserve(context.Background(), "tenant", Reservation{}, 1); err == nil {
		t.Fatal("invalid reservation accepted")
	}
	if _, err := service.CreateTransfer(context.Background(), "tenant", Transfer{FromOrganizationID: "a", ToOrganizationID: "a", StockUnitIDs: []string{"s"}, ExpectedReceiptAt: time.Now()}, map[string]int64{"s": 1}); err == nil {
		t.Fatal("same-organization transfer accepted")
	}
	if err := service.TransitionTransfer(context.Background(), "tenant", "a", "t", "draft", "received", 1); err == nil {
		t.Fatal("skipped transfer transition accepted")
	}
	if repo.reserve+repo.create+repo.transition != 0 {
		t.Fatal("repository called for invalid input")
	}
}

func TestServiceBuildsGovernedCommands(t *testing.T) {
	repo, generator := &fakeRepo{}, &ids{}
	service := NewService(repo, generator)
	reservation, err := service.Reserve(context.Background(), "tenant", Reservation{OrganizationID: "a", StockUnitID: "s", VariantID: "v", DemandKind: "service", DemandID: "d"}, 2)
	if err != nil || reservation.Status != "reservation" || reservation.Version != 1 || repo.reserve != 1 {
		t.Fatalf("reservation not created: %#v %v", reservation, err)
	}
	transfer, err := service.CreateTransfer(context.Background(), "tenant", Transfer{FromOrganizationID: "a", ToOrganizationID: "b", StockUnitIDs: []string{"s"}, ExpectedReceiptAt: time.Now().Add(time.Hour)}, map[string]int64{"s": 3})
	if err != nil || transfer.State != "draft" || transfer.Version != 1 || repo.create != 1 {
		t.Fatalf("transfer not created: %#v %v", transfer, err)
	}
}
````

### FILE: `internal/platform/postgres/inventorycontrol.go`

```yaml
block_id: "GO-OPS-API:inventory-control-postgres:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps 2eae56d reservation/ATP/transfer invariants; local PostgreSQL transaction port"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "18a9f8e45ed0f088feea05d787560a829fca2602ccda364f82b404b82b16fc33"
variables: []
secrets_allowed: false
```

````go
// Inventory control persistence implements the Microsoft BC-derived portable
// contract recorded in docs/inventory/MICROSOFT_BC_DERIVATION.md.
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryControl struct{ pool *pgxpool.Pool }

func NewInventoryControl(pool *pgxpool.Pool) *InventoryControl { return &InventoryControl{pool: pool} }

func (r *InventoryControl) AvailableToPromise(ctx context.Context, tenant, organization, variant string, horizon time.Time) (inventorycontrol.ATP, error) {
	value := inventorycontrol.ATP{OrganizationID: organization, VariantID: variant, Horizon: horizon}
	err := r.pool.QueryRow(ctx, `select available_inventory,scheduled_receipt,gross_requirement,available_to_promise from inventory.serial_atp($1,$2,$3,$4)`, tenant, organization, variant, horizon).Scan(&value.AvailableInventory, &value.ScheduledReceipt, &value.GrossRequirement, &value.AvailableToPromise)
	return value, err
}

func (r *InventoryControl) Reserve(ctx context.Context, tenant, eventID, reservationID string, value inventorycontrol.Reservation, stockVersion int64) (inventorycontrol.Reservation, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update inventory.stock_unit set state='reserved',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and variant_id=$4 and state='available' and version=$5`, tenant, value.OrganizationID, value.StockUnitID, value.VariantID, stockVersion)
	if err != nil {
		return value, err
	}
	if result.RowsAffected() != 1 {
		return value, inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into inventory.serial_reservation(tenant_id,reservation_id,organization_id,stock_unit_id,variant_id,demand_kind,demand_id,demand_line_id,status,cancellation_disallowed,expires_at,version) values($1,$2,$3,$4,$5,$6,$7,$8,'reservation',$9,$10,1)`, tenant, reservationID, value.OrganizationID, value.StockUnitID, value.VariantID, value.DemandKind, value.DemandID, value.DemandLineID, value.CancellationDisallowed, value.ExpiresAt)
	if err != nil {
		return value, err
	}
	payload, _ := json.Marshal(map[string]string{"reservation_id": reservationID, "stock_unit_id": value.StockUnitID, "demand_kind": value.DemandKind, "demand_id": value.DemandID})
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'serial-reservation',$3,1,'serial-reservation.created',1,clock_timestamp(),$4)`, tenant, eventID, reservationID, payload)
	if err != nil {
		return value, err
	}
	if err = tx.Commit(ctx); err != nil {
		return value, err
	}
	return value, nil
}

func (r *InventoryControl) Release(ctx context.Context, tenant, organization, reservationID string, version int64, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var stockID string
	err = tx.QueryRow(ctx, `update inventory.serial_reservation set status='released',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and reservation_id=$3 and status='reservation' and version=$4 and not cancellation_disallowed returning stock_unit_id`, tenant, organization, reservationID, version).Scan(&stockID)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.ErrConflict
	}
	if err != nil {
		return err
	}
	result, err := tx.Exec(ctx, `update inventory.stock_unit set state='available',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and state='reserved'`, tenant, organization, stockID)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'serial-reservation',$3,$4,'serial-reservation.released',1,clock_timestamp(),jsonb_build_object('stock_unit_id',$5::text))`, tenant, eventID, reservationID, version+1, stockID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *InventoryControl) CreateTransfer(ctx context.Context, tenant, transferID string, value inventorycontrol.Transfer, stockVersions map[string]int64, eventID string) (inventorycontrol.Transfer, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into inventory.serial_transfer(tenant_id,transfer_id,from_organization_id,to_organization_id,expected_receipt_at,state,version) values($1,$2,$3,$4,$5,'draft',1)`, tenant, transferID, value.FromOrganizationID, value.ToOrganizationID, value.ExpectedReceiptAt)
	if err != nil {
		return value, err
	}
	for _, stockID := range value.StockUnitIDs {
		var variant string
		err = tx.QueryRow(ctx, `update inventory.stock_unit set state='reserved',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and state='available' and version=$4 returning variant_id`, tenant, value.FromOrganizationID, stockID, stockVersions[stockID]).Scan(&variant)
		if errors.Is(err, pgx.ErrNoRows) {
			return value, inventorycontrol.ErrConflict
		}
		if err != nil {
			return value, err
		}
		_, err = tx.Exec(ctx, `insert into inventory.serial_transfer_unit(tenant_id,transfer_id,stock_unit_id,variant_id) values($1,$2,$3,$4)`, tenant, transferID, stockID, variant)
		if err != nil {
			return value, err
		}
		_, err = tx.Exec(ctx, `insert into inventory.serial_reservation(tenant_id,reservation_id,organization_id,stock_unit_id,variant_id,demand_kind,demand_id,demand_line_id,status,version) values($1,$2,$3,$4,$5,'transfer-outbound',$6,$4,'reservation',1)`, tenant, transferID+":"+stockID, value.FromOrganizationID, stockID, variant, transferID)
		if err != nil {
			return value, err
		}
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'serial-transfer',$3,1,'serial-transfer.created',1,clock_timestamp(),jsonb_build_object('from_organization_id',$4::text,'to_organization_id',$5::text,'unit_count',$6::int))`, tenant, eventID, transferID, value.FromOrganizationID, value.ToOrganizationID, len(value.StockUnitIDs))
	if err != nil {
		return value, err
	}
	if err = tx.Commit(ctx); err != nil {
		return value, err
	}
	return value, nil
}

func (r *InventoryControl) TransitionTransfer(ctx context.Context, tenant, authorizedOrganization, transferID, current string, version int64, target, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var fromOrganization, toOrganization string
	err = tx.QueryRow(ctx, `select from_organization_id,to_organization_id from inventory.serial_transfer where tenant_id=$1 and transfer_id=$2 and state=$3 and version=$4 for update`, tenant, transferID, current, version).Scan(&fromOrganization, &toOrganization)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.ErrConflict
	}
	if err != nil {
		return err
	}
	if (target == "received" && authorizedOrganization != toOrganization) || (target != "received" && authorizedOrganization != fromOrganization) {
		return inventorycontrol.ErrConflict
	}
	rows, err := tx.Query(ctx, `select stock_unit_id from inventory.serial_transfer_unit where tenant_id=$1 and transfer_id=$2 order by stock_unit_id for update`, tenant, transferID)
	if err != nil {
		return err
	}
	var stockIDs []string
	for rows.Next() {
		var stockID string
		if err = rows.Scan(&stockID); err != nil {
			rows.Close()
			return err
		}
		stockIDs = append(stockIDs, stockID)
	}
	rows.Close()
	if err = rows.Err(); err != nil || len(stockIDs) == 0 {
		if err != nil {
			return err
		}
		return inventorycontrol.ErrConflict
	}
	for _, stockID := range stockIDs {
		var result pgconn.CommandTag
		switch target {
		case "in-transit":
			result, err = tx.Exec(ctx, `update inventory.stock_unit set state='in-transit',version=version+1,received_at=null,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and state='reserved'`, tenant, fromOrganization, stockID)
		case "received":
			result, err = tx.Exec(ctx, `update inventory.stock_unit set organization_id=$3,state='available',version=version+1,received_at=clock_timestamp(),updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$4 and state='in-transit'`, tenant, fromOrganization, toOrganization, stockID)
		case "cancelled":
			result, err = tx.Exec(ctx, `update inventory.stock_unit set state='available',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and state='reserved'`, tenant, fromOrganization, stockID)
		default:
			continue
		}
		if err != nil {
			return err
		}
		if result.RowsAffected() != 1 {
			return inventorycontrol.ErrConflict
		}
	}
	if target == "in-transit" {
		_, err = tx.Exec(ctx, `update inventory.serial_reservation set status='consumed',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and demand_kind='transfer-outbound' and demand_id=$2 and status='reservation'`, tenant, transferID)
	} else if target == "cancelled" {
		_, err = tx.Exec(ctx, `update inventory.serial_reservation set status='released',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and demand_kind='transfer-outbound' and demand_id=$2 and status='reservation'`, tenant, transferID)
	}
	if err != nil {
		return err
	}
	result, err := tx.Exec(ctx, `update inventory.serial_transfer set state=$5,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and transfer_id=$2 and state=$3 and version=$4`, tenant, transferID, current, version, target)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'serial-transfer',$3,$4,$5,1,clock_timestamp(),jsonb_build_object('unit_count',$6::int))`, tenant, eventID, transferID, version+1, "serial-transfer."+target, len(stockIDs))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
````

### FILE: `internal/platform/postgres/inventorycontrol_integration_test.go`

```yaml
block_id: "GO-OPS-API:inventory-control-postgres-test:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps 2eae56d SCMTransferReservation tests; local concurrency/scope/retransfer regression"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "c1c5602970c701ece86267e390758088d570d812319ffb56e17743904536c365"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestInventoryReservationATPTransferAndConcurrency(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c29901"
	cleanup := func() {
		for _, q := range []string{`delete from platform.outbox_event where tenant_id=$1`, `delete from inventory.serial_transfer_unit where tenant_id=$1`, `delete from inventory.serial_transfer where tenant_id=$1`, `delete from inventory.serial_reservation where tenant_id=$1`, `delete from inventory.stock_unit where tenant_id=$1`, `delete from catalog.vehicle_variant where tenant_id=$1`, `delete from catalog.vehicle_model where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`} {
			_, _ = pool.Exec(ctx, q, tenant)
		}
	}
	cleanup()
	defer cleanup()
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'inventory-v159','Inventory','Inventory')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'a','a','A','warehouse'),($1,'b','b','B','franchisee')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'model','model','Model','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state) values($1,'variant','model','variant','Variant','{}','active')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) values($1,'s1','a','variant','SERIAL-1','VIN-1','BATTERY-1','available',1,clock_timestamp()),($1,'s2','a','variant','SERIAL-2','VIN-2','BATTERY-2','available',1,clock_timestamp()),($1,'s3','a','variant','SERIAL-3','VIN-3','BATTERY-3','available',1,clock_timestamp())`,
	}
	for _, q := range fixtures {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewInventoryControl(pool)
	horizon := time.Now().Add(24 * time.Hour).UTC()
	atp, err := repo.AvailableToPromise(ctx, tenant, "a", "variant", horizon)
	if err != nil || atp.AvailableInventory != 3 || atp.AvailableToPromise != 3 {
		t.Fatalf("initial ATP=%+v err=%v", atp, err)
	}
	values := []inventorycontrol.Reservation{
		{ID: "r1", OrganizationID: "a", StockUnitID: "s1", VariantID: "variant", DemandKind: "service", DemandID: "service-1", Status: "reservation", Version: 1},
		{ID: "r2", OrganizationID: "a", StockUnitID: "s1", VariantID: "variant", DemandKind: "service", DemandID: "service-2", Status: "reservation", Version: 1},
	}
	errs := make([]error, 2)
	var wg sync.WaitGroup
	for index := range values {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, errs[i] = repo.Reserve(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c2991"+string(rune('0'+i)), values[i].ID, values[i], 1)
		}(index)
	}
	wg.Wait()
	winner, success := "", 0
	for i, callErr := range errs {
		if callErr == nil {
			winner, success = values[i].ID, success+1
		} else if !errors.Is(callErr, inventorycontrol.ErrConflict) {
			t.Fatalf("unexpected concurrent error: %v", callErr)
		}
	}
	if success != 1 {
		t.Fatalf("concurrent reservation successes=%d errors=%v", success, errs)
	}
	if err = repo.Release(ctx, tenant, "a", winner, 1, "018f4d4a-7b36-7a21-8d10-2f4c54c29912"); err != nil {
		t.Fatal(err)
	}
	transfer := inventorycontrol.Transfer{ID: "t1", FromOrganizationID: "a", ToOrganizationID: "b", StockUnitIDs: []string{"s2"}, ExpectedReceiptAt: time.Now().Add(time.Hour).UTC(), State: "draft", Version: 1}
	if _, err = repo.CreateTransfer(ctx, tenant, transfer.ID, transfer, map[string]int64{"s2": 1}, "018f4d4a-7b36-7a21-8d10-2f4c54c29913"); err != nil {
		t.Fatal(err)
	}
	if err = repo.TransitionTransfer(ctx, tenant, "a", "t1", "draft", 1, "released", "018f4d4a-7b36-7a21-8d10-2f4c54c29914"); err != nil {
		t.Fatal(err)
	}
	if err = repo.TransitionTransfer(ctx, tenant, "a", "t1", "released", 2, "in-transit", "018f4d4a-7b36-7a21-8d10-2f4c54c29915"); err != nil {
		t.Fatal(err)
	}
	atp, err = repo.AvailableToPromise(ctx, tenant, "b", "variant", horizon)
	if err != nil || atp.ScheduledReceipt != 1 || atp.AvailableToPromise != 1 {
		t.Fatalf("inbound ATP=%+v err=%v", atp, err)
	}
	if err = repo.TransitionTransfer(ctx, tenant, "a", "t1", "in-transit", 3, "received", "018f4d4a-7b36-7a21-8d10-2f4c54c29916"); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("wrong receiver accepted: %v", err)
	}
	if err = repo.TransitionTransfer(ctx, tenant, "b", "t1", "in-transit", 3, "received", "018f4d4a-7b36-7a21-8d10-2f4c54c29917"); err != nil {
		t.Fatal(err)
	}
	var stockVersion int64
	if err = pool.QueryRow(ctx, `select version from inventory.stock_unit where tenant_id=$1 and stock_unit_id='s2' and organization_id='b' and state='available'`, tenant).Scan(&stockVersion); err != nil || stockVersion != 4 {
		t.Fatalf("received stock version=%d err=%v", stockVersion, err)
	}
	returnTransfer := inventorycontrol.Transfer{ID: "t2", FromOrganizationID: "b", ToOrganizationID: "a", StockUnitIDs: []string{"s2"}, ExpectedReceiptAt: time.Now().Add(2 * time.Hour).UTC(), State: "draft", Version: 1}
	if _, err = repo.CreateTransfer(ctx, tenant, returnTransfer.ID, returnTransfer, map[string]int64{"s2": stockVersion}, "018f4d4a-7b36-7a21-8d10-2f4c54c29918"); err != nil {
		t.Fatalf("received unit was not transferable again: %v", err)
	}
	if err = repo.TransitionTransfer(ctx, tenant, "b", "t2", "draft", 1, "cancelled", "018f4d4a-7b36-7a21-8d10-2f4c54c29919"); err != nil {
		t.Fatal(err)
	}
}
````

### FILE: `internal/platform/httpapi/inventorycontrol.go`

```yaml
block_id: "GO-OPS-API:inventory-control-http:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified HTTP composition over the admitted inventory control port"
license: "LicenseRef-Workspace-Owner"
sha256: "22db7ef2fd6d500c5a19acc2467418433fa952ece2903d074f0248fe1f2c210c"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"errors"
	"net/http"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
)

type InventoryControlModule struct {
	Service   *inventorycontrol.Service
	Bulk      *inventorycontrol.BulkService
	Warehouse *inventorycontrol.WarehouseService
	Transfer  *inventorycontrol.BulkTransferService
}

func (m InventoryControlModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	api := inventoryControlAPI{service: m.Service, verifier: verifier}
	mux.HandleFunc("GET /v1/inventory/atp", api.atp)
	mux.HandleFunc("POST /v1/inventory/reservations", api.reserve)
	mux.HandleFunc("POST /v1/inventory/reservations/{id}/release", api.release)
	mux.HandleFunc("POST /v1/inventory/transfers", api.createTransfer)
	mux.HandleFunc("POST /v1/inventory/transfers/{id}/transitions", api.transitionTransfer)
	if m.Bulk != nil {
		registerBulkInventory(mux, m.Bulk, verifier)
	}
	if m.Warehouse != nil {
		registerWarehouse(mux, m.Warehouse, verifier)
	}
	if m.Transfer != nil {
		registerBulkTransfer(mux, m.Transfer, verifier)
	}
}

type inventoryControlAPI struct {
	service  *inventorycontrol.Service
	verifier identity.Verifier
}

func (a inventoryControlAPI) authorize(w http.ResponseWriter, r *http.Request, permission string, requireJSON bool) (identity.Principal, bool) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return p, false
	}
	if requireJSON && r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return p, false
	}
	return p, true
}

func (a inventoryControlAPI) atp(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:read", false)
	if !ok {
		return
	}
	organization, variant := r.URL.Query().Get("organization_id"), r.URL.Query().Get("variant_id")
	if !p.AllowedOrganization(organization) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	horizon, err := time.Parse(time.RFC3339, r.URL.Query().Get("horizon"))
	if err != nil {
		writeProblem(w, 400, "INVALID_ATP_QUERY", "horizon must be RFC3339")
		return
	}
	value, err := a.service.AvailableToPromise(r.Context(), p.TenantID, organization, variant, horizon)
	if err != nil {
		writeProblem(w, 500, "ATP_FAILED", "availability could not be calculated")
		return
	}
	writeJSON(w, 200, value)
}

func (a inventoryControlAPI) reserve(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input struct {
		inventorycontrol.Reservation
		StockVersion int64 `json:"stock_version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.Reserve(r.Context(), p.TenantID, input.Reservation, input.StockVersion)
	if errors.Is(err, inventorycontrol.ErrConflict) {
		writeProblem(w, 409, "RESERVATION_CONFLICT", "stock is no longer available")
		return
	}
	if err != nil {
		writeProblem(w, 400, "INVALID_RESERVATION", "reservation does not match contract")
		return
	}
	writeJSON(w, 201, value)
}

func (a inventoryControlAPI) release(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeInventoryTransition(w, a.service.Release(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Version))
}

func (a inventoryControlAPI) createTransfer(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input struct {
		inventorycontrol.Transfer
		StockVersions map[string]int64 `json:"stock_versions"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.FromOrganizationID) || !p.AllowedOrganization(input.ToOrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for both organizations")
		return
	}
	value, err := a.service.CreateTransfer(r.Context(), p.TenantID, input.Transfer, input.StockVersions)
	if errors.Is(err, inventorycontrol.ErrConflict) {
		writeProblem(w, 409, "TRANSFER_CONFLICT", "one or more stock units are unavailable")
		return
	}
	if err != nil {
		writeProblem(w, 400, "INVALID_TRANSFER", "transfer does not match contract")
		return
	}
	writeJSON(w, 201, value)
}

func (a inventoryControlAPI) transitionTransfer(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		Current        string `json:"current"`
		Target         string `json:"target"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeInventoryTransition(w, a.service.TransitionTransfer(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Current, input.Target, input.Version))
}

func writeInventoryTransition(w http.ResponseWriter, err error) {
	if errors.Is(err, inventorycontrol.ErrConflict) {
		writeProblem(w, 409, "INVENTORY_CONFLICT", "state, scope or version conflict")
		return
	}
	if err != nil {
		writeProblem(w, 500, "INTERNAL_ERROR", "inventory transition failed")
		return
	}
	writeJSON(w, 200, map[string]string{"status": "accepted"})
}
````

### FILE: `internal/platform/httpapi/inventorycontrol_test.go`

```yaml
block_id: "GO-OPS-API:inventory-control-http-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local negative authorization and contract tests"
license: "LicenseRef-Workspace-Owner"
sha256: "d72d31c2f0a958129748fa5652cea53bc02ebd5b6ca6eaa48090ffff1741c49d"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
)

type inventoryRepo struct{ atp, reserve int }

func (f *inventoryRepo) AvailableToPromise(context.Context, string, string, string, time.Time) (inventorycontrol.ATP, error) {
	f.atp++
	return inventorycontrol.ATP{OrganizationID: "a", VariantID: "v", AvailableToPromise: 3}, nil
}
func (f *inventoryRepo) Reserve(_ context.Context, _, _, _ string, value inventorycontrol.Reservation, _ int64) (inventorycontrol.Reservation, error) {
	f.reserve++
	return value, nil
}
func (f *inventoryRepo) Release(context.Context, string, string, string, int64, string) error {
	return inventorycontrol.ErrConflict
}
func (f *inventoryRepo) CreateTransfer(_ context.Context, _, _ string, value inventorycontrol.Transfer, _ map[string]int64, _ string) (inventorycontrol.Transfer, error) {
	return value, nil
}
func (f *inventoryRepo) TransitionTransfer(context.Context, string, string, string, string, int64, string, string) error {
	return nil
}

type inventoryIDs struct{}

func (inventoryIDs) New() string { return "018f4d4a-7b36-7a21-8d10-2f4c54c29999" }

type inventoryVerifier struct{}

func (inventoryVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{Subject: "operator", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c29990", Permissions: map[string]struct{}{"inventory:read": {}, "inventory:write": {}}, Organizations: map[string]struct{}{"a": {}, "b": {}}}, nil
}

func TestInventoryControlHTTPATPAndScope(t *testing.T) {
	repo := &inventoryRepo{}
	mux := http.NewServeMux()
	InventoryControlModule{Service: inventorycontrol.NewService(repo, inventoryIDs{})}.Register(mux, inventoryVerifier{})
	request := httptest.NewRequest("GET", "/v1/inventory/atp?organization_id=a&variant_id=v&horizon=2030-01-01T00:00:00Z", nil)
	request.Header.Set("Authorization", "Bearer valid")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 200 || repo.atp != 1 || !strings.Contains(response.Body.String(), `"available_to_promise":3`) {
		t.Fatalf("ATP status=%d calls=%d body=%s", response.Code, repo.atp, response.Body.String())
	}
	request = httptest.NewRequest("GET", "/v1/inventory/atp?organization_id=forbidden&variant_id=v&horizon=2030-01-01T00:00:00Z", nil)
	request.Header.Set("Authorization", "Bearer valid")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 || repo.atp != 1 {
		t.Fatalf("scope status=%d calls=%d", response.Code, repo.atp)
	}
}

func TestInventoryControlHTTPRejectsCustomerOrderAndMapsConflict(t *testing.T) {
	repo := &inventoryRepo{}
	mux := http.NewServeMux()
	InventoryControlModule{Service: inventorycontrol.NewService(repo, inventoryIDs{})}.Register(mux, inventoryVerifier{})
	request := httptest.NewRequest("POST", "/v1/inventory/reservations", strings.NewReader(`{"organization_id":"a","stock_unit_id":"s","variant_id":"v","demand_kind":"customer-order","demand_id":"o","stock_version":1}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 400 || repo.reserve != 0 {
		t.Fatalf("customer-order status=%d calls=%d", response.Code, repo.reserve)
	}
	request = httptest.NewRequest("POST", "/v1/inventory/reservations/r/release", strings.NewReader(`{"organization_id":"a","version":1}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 409 {
		t.Fatalf("release conflict status=%d body=%s", response.Code, response.Body.String())
	}
}
````

### FILE: `db/migrations/0026_bulk_lot_bin_costing.up.sql`

```yaml
block_id: "GO-OPS-API:bulk-lot-bin-costing-migration-up:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps 2eae56d Tracking, BinContent, WarehouseAvailability, CostingMethod and ItemApplication invariants; portable PostgreSQL scope"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "fd53f303e0eb1d014ff769447021eb75764d2c28aac6f99e4e8293cd06b290da"
variables: []
secrets_allowed: false
```

````sql
begin;

create table inventory.stock_item (
  tenant_id uuid not null,
  item_id text not null,
  item_code text not null,
  description text not null,
  base_uom text not null,
  tracking_mode text not null check (tracking_mode in ('none', 'lot')),
  costing_method text not null check (costing_method in ('fifo', 'specific')),
  version bigint not null default 1 check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, item_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  unique (tenant_id, item_code),
  check (length(item_id) between 1 and 128),
  check (item_code ~ '^[A-Za-z0-9][A-Za-z0-9._/-]{0,63}$'),
  check (base_uom ~ '^[A-Z0-9][A-Z0-9._/-]{0,15}$')
);

create table inventory.warehouse_bin (
  tenant_id uuid not null,
  organization_id text not null,
  bin_id text not null,
  bin_code text not null,
  bin_type text not null check (bin_type in ('receive', 'ship', 'put-away', 'pick', 'putpick')),
  movement_blocked boolean not null default false,
  version bigint not null default 1 check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, organization_id, bin_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  unique (tenant_id, organization_id, bin_code),
  check (length(bin_id) between 1 and 128),
  check (bin_code ~ '^[A-Za-z0-9][A-Za-z0-9._/-]{0,63}$')
);

create table inventory.item_bin_policy (
  tenant_id uuid not null,
  organization_id text not null,
  item_id text not null,
  bin_id text not null,
  fixed boolean not null default false,
  dedicated boolean not null default false,
  is_default boolean not null default false,
  min_quantity numeric(20,6) not null default 0 check (min_quantity >= 0),
  max_quantity numeric(20,6),
  version bigint not null default 1 check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, organization_id, item_id, bin_id),
  foreign key (tenant_id, organization_id, bin_id) references inventory.warehouse_bin (tenant_id, organization_id, bin_id),
  foreign key (tenant_id, item_id) references inventory.stock_item (tenant_id, item_id),
  check (max_quantity is null or max_quantity > min_quantity)
);

create unique index item_bin_policy_one_default_uidx
  on inventory.item_bin_policy(tenant_id, organization_id, item_id)
  where is_default;

create table inventory.inventory_lot (
  tenant_id uuid not null,
  lot_id text not null,
  item_id text not null,
  lot_no text not null,
  expiration_date date,
  warranty_date date,
  blocked boolean not null default false,
  version bigint not null default 1 check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, lot_id),
  foreign key (tenant_id, item_id) references inventory.stock_item (tenant_id, item_id),
  unique (tenant_id, item_id, lot_no),
  check (length(lot_id) between 1 and 128),
  check (length(lot_no) between 1 and 50)
);

create table inventory.bulk_balance (
  tenant_id uuid not null,
  organization_id text not null,
  bin_id text not null,
  item_id text not null,
  lot_id text,
  quantity numeric(20,6) not null default 0 check (quantity >= 0),
  reserved_quantity numeric(20,6) not null default 0 check (reserved_quantity >= 0 and reserved_quantity <= quantity),
  version bigint not null default 1 check (version > 0),
  updated_at timestamptz not null default clock_timestamp(),
  balance_id bigint generated always as identity primary key,
  foreign key (tenant_id, organization_id, item_id, bin_id) references inventory.item_bin_policy (tenant_id, organization_id, item_id, bin_id),
  foreign key (tenant_id, lot_id) references inventory.inventory_lot (tenant_id, lot_id),
  unique nulls not distinct (tenant_id, organization_id, bin_id, item_id, lot_id)
);

create table inventory.bulk_reservation (
  tenant_id uuid not null,
  reservation_id text not null,
  organization_id text not null,
  bin_id text not null,
  item_id text not null,
  lot_id text,
  demand_kind text not null check (demand_kind in ('customer-order', 'service', 'transfer-outbound', 'manual')),
  demand_id text not null,
  demand_line_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  status text not null check (status in ('reservation', 'released', 'consumed', 'expired')),
  cancellation_disallowed boolean not null default false,
  expires_at timestamptz,
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, reservation_id),
  foreign key (tenant_id, organization_id, item_id, bin_id) references inventory.item_bin_policy (tenant_id, organization_id, item_id, bin_id),
  foreign key (tenant_id, lot_id) references inventory.inventory_lot (tenant_id, lot_id),
  check (expires_at is null or expires_at > created_at)
);

create unique index bulk_reservation_active_demand_uidx
  on inventory.bulk_reservation(tenant_id, demand_kind, demand_id, demand_line_id)
  where status='reservation';

create index bulk_reservation_scope_idx
  on inventory.bulk_reservation(tenant_id, organization_id, item_id, status, expires_at);

create table inventory.bulk_inventory_entry (
  tenant_id uuid not null,
  entry_id text not null,
  organization_id text not null,
  bin_id text not null,
  item_id text not null,
  lot_id text,
  entry_type text not null check (entry_type in ('receipt', 'issue', 'movement-in', 'movement-out', 'adjustment')),
  quantity numeric(20,6) not null check (quantity <> 0),
  unit_cost numeric(20,4) not null check (unit_cost >= 0),
  cost_amount numeric(24,4) not null,
  posting_date date not null,
  source_kind text not null,
  source_id text not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, entry_id),
  foreign key (tenant_id, organization_id, item_id, bin_id) references inventory.item_bin_policy (tenant_id, organization_id, item_id, bin_id),
  foreign key (tenant_id, lot_id) references inventory.inventory_lot (tenant_id, lot_id),
  unique (tenant_id, source_kind, source_id, entry_type, organization_id, bin_id),
  check ((quantity > 0 and entry_type in ('receipt','movement-in','adjustment')) or (quantity < 0 and entry_type in ('issue','movement-out','adjustment')))
);

create table inventory.bulk_cost_layer (
  tenant_id uuid not null,
  layer_id text not null,
  receipt_entry_id text not null,
  organization_id text not null,
  item_id text not null,
  lot_id text,
  posting_date date not null,
  original_quantity numeric(20,6) not null check (original_quantity > 0),
  remaining_quantity numeric(20,6) not null check (remaining_quantity >= 0 and remaining_quantity <= original_quantity),
  unit_cost numeric(20,4) not null check (unit_cost >= 0),
  primary key (tenant_id, layer_id),
  foreign key (tenant_id, receipt_entry_id) references inventory.bulk_inventory_entry (tenant_id, entry_id),
  foreign key (tenant_id, item_id) references inventory.stock_item (tenant_id, item_id),
  foreign key (tenant_id, lot_id) references inventory.inventory_lot (tenant_id, lot_id),
  unique (tenant_id, receipt_entry_id)
);

create index bulk_cost_layer_fifo_idx
  on inventory.bulk_cost_layer(tenant_id, organization_id, item_id, posting_date, receipt_entry_id)
  where remaining_quantity > 0;

create table inventory.bulk_cost_application (
  tenant_id uuid not null,
  outbound_entry_id text not null,
  inbound_entry_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  cost_amount numeric(24,4) not null check (cost_amount >= 0),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, outbound_entry_id, inbound_entry_id),
  foreign key (tenant_id, outbound_entry_id) references inventory.bulk_inventory_entry (tenant_id, entry_id),
  foreign key (tenant_id, inbound_entry_id) references inventory.bulk_inventory_entry (tenant_id, entry_id)
);

create function inventory.reject_bulk_ledger_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='bulk inventory ledger is append-only';
end;
$function$;

create trigger bulk_inventory_entry_immutable
before update or delete on inventory.bulk_inventory_entry
for each row execute function inventory.reject_bulk_ledger_mutation();

create trigger bulk_cost_application_immutable
before update or delete on inventory.bulk_cost_application
for each row execute function inventory.reject_bulk_ledger_mutation();

create view inventory.bulk_available as
select b.tenant_id,b.organization_id,b.bin_id,b.item_id,b.lot_id,
       b.quantity,b.reserved_quantity,(b.quantity-b.reserved_quantity) available_quantity,
       l.lot_no,l.expiration_date,p.fixed,p.dedicated,p.is_default,w.bin_code,w.bin_type
from inventory.bulk_balance b
join inventory.item_bin_policy p using(tenant_id,organization_id,item_id,bin_id)
join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id)
left join inventory.inventory_lot l using(tenant_id,lot_id)
where not w.movement_blocked
  and not p.dedicated
  and coalesce(l.blocked,false)=false
  and (l.expiration_date is null or l.expiration_date >= current_date);

comment on view inventory.bulk_available is
  'Portable BC-derived warehouse availability: physical bin quantity minus reserved quantity, excluding movement-blocked/dedicated bins and blocked/expired lots.';

commit;
````

### FILE: `db/migrations/0026_bulk_lot_bin_costing.down.sql`

```yaml
block_id: "GO-OPS-API:bulk-lot-bin-costing-migration-down:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps-derived portable rollback counterpart; local schema lifecycle"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "297323f34bc7a91b01e50723073fef5f27dc3cb21c39373f7e67130653ddbe35"
variables: []
secrets_allowed: false
```

````sql
begin;
drop view if exists inventory.bulk_available;
drop trigger if exists bulk_cost_application_immutable on inventory.bulk_cost_application;
drop trigger if exists bulk_inventory_entry_immutable on inventory.bulk_inventory_entry;
drop function if exists inventory.reject_bulk_ledger_mutation();
drop table if exists inventory.bulk_cost_application;
drop table if exists inventory.bulk_cost_layer;
drop table if exists inventory.bulk_inventory_entry;
drop table if exists inventory.bulk_reservation;
drop table if exists inventory.bulk_balance;
drop table if exists inventory.inventory_lot;
drop table if exists inventory.item_bin_policy;
drop table if exists inventory.warehouse_bin;
drop table if exists inventory.stock_item;
commit;
````

### FILE: `docs/inventory/MICROSOFT_BC_BULK_DERIVATION.md`

```yaml
block_id: "GO-OPS-API:microsoft-bc-bulk-derivation:v1"
operation: CREATE
provenance: AUTHORED
source: "Pinned Microsoft BCApps source and Microsoft Learn authorities enumerated in the document"
license: "LicenseRef-Workspace-Owner"
sha256: "4a3ba733da1902be82ae9b9f7ddfb77f4a45c224920dc0fccc1f10c7b267a554"
variables: []
secrets_allowed: false
```

````markdown
# Microsoft BC-derived bulk, lot, bin and costing boundary

This portable implementation is `ADAPTED`; it is not verbatim Microsoft Go code and is not a Business Central runtime.

Authority is Microsoft `BCApps` commit `2eae56d704a1fd035d104f333602aea7091b7749`, root MIT license. The audited files are `TrackingSpecification.Table.al`, `LotNoInformation.Table.al`, `Bin.Table.al`, `BinContent.Table.al`, `WarehouseAvailabilityMgt.Codeunit.al`, `CostingMethod.Enum.al`, `AvgCostAdjmtEntryPoint.Table.al` and `ItemApplicationEntry.Table.al` under `src/Layers/W1/BaseApp`.

The portable contract preserves these narrow invariants:

1. lot identity is item-scoped and carries blocked/expiration state;
2. a bin is the smallest physical storage unit, while item/bin policy owns fixed, dedicated and default behavior;
3. sellable availability is physical quantity less reservations and excludes blocked movement, dedicated generic stock and blocked/expired lots;
4. quantity reservations are atomic and cannot exceed balance;
5. every issue is linked to one or more inbound cost entries through immutable application rows;
6. FIFO consumes the oldest remaining inbound layer; specific costing requires the exact inbound entry;
7. movements inside one organization preserve valuation and do not manufacture a cost;
8. serial-tracked vehicles remain in the existing serial owner and cannot enter the bulk API.

Business Central also defines LIFO, Average and Standard costing and a much larger adjustment engine. This portable boundary deliberately rejects those methods. It does not claim periodic average adjustment, manufacturing variance, expected cost, G/L posting, warehouse picks/put-aways, cross-organization bulk transfer or full CTP/planning.

The project must select accounting policy, currency/rounding, landed-cost treatment, closed periods, expiry policy and warehouse roles before production admission.
````

### FILE: `internal/inventorycontrol/bulk.go`

```yaml
block_id: "GO-OPS-API:bulk-inventory-service:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps 2eae56d lot/bin/availability/cost-application invariants; local Go domain port"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "7498fbcbc558f6d35eeb8e5de88d5954a0858b95e584f7c61227944420c9fcac"
variables: []
secrets_allowed: false
```

````go
package inventorycontrol

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"time"
)

var (
	quantityPattern = regexp.MustCompile(`^(0|[1-9][0-9]{0,13})(\.[0-9]{1,6})?$`)
	costPattern     = regexp.MustCompile(`^(0|[1-9][0-9]{0,15})(\.[0-9]{1,4})?$`)
)

type BulkItem struct {
	ID                    string `json:"id"`
	Code                  string `json:"code"`
	Description           string `json:"description"`
	BaseUOM               string `json:"base_uom"`
	BaseRoundingPrecision string `json:"base_rounding_precision"`
	TrackingMode          string `json:"tracking_mode"`
	CostingMethod         string `json:"costing_method"`
	Version               int64  `json:"version"`
}

type WarehouseBin struct {
	ID              string `json:"id"`
	OrganizationID  string `json:"organization_id"`
	Code            string `json:"code"`
	Type            string `json:"type"`
	Ranking         int    `json:"ranking"`
	MovementBlocked bool   `json:"movement_blocked"`
	CrossDock       bool   `json:"cross_dock"`
	Version         int64  `json:"version"`
}

type ItemBinPolicy struct {
	OrganizationID string `json:"organization_id"`
	ItemID         string `json:"item_id"`
	BinID          string `json:"bin_id"`
	Fixed          bool   `json:"fixed"`
	Dedicated      bool   `json:"dedicated"`
	Default        bool   `json:"default"`
	MinQuantity    string `json:"min_quantity"`
	MaxQuantity    string `json:"max_quantity,omitempty"`
	Version        int64  `json:"version"`
}

type BulkReceipt struct {
	OrganizationID string     `json:"organization_id"`
	BinID          string     `json:"bin_id"`
	ItemID         string     `json:"item_id"`
	LotNo          string     `json:"lot_no,omitempty"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty"`
	WarrantyDate   *time.Time `json:"warranty_date,omitempty"`
	Quantity       string     `json:"quantity"`
	UnitCost       string     `json:"unit_cost"`
	PostingDate    time.Time  `json:"posting_date"`
	SourceKind     string     `json:"source_kind"`
	SourceID       string     `json:"source_id"`
}

type BulkReceiptResult struct {
	EntryID    string `json:"entry_id"`
	LotID      string `json:"lot_id,omitempty"`
	Quantity   string `json:"quantity"`
	UnitCost   string `json:"unit_cost"`
	CostAmount string `json:"cost_amount"`
}

type BulkAvailability struct {
	OrganizationID    string     `json:"organization_id"`
	BinID             string     `json:"bin_id"`
	BinCode           string     `json:"bin_code"`
	ItemID            string     `json:"item_id"`
	LotID             string     `json:"lot_id,omitempty"`
	LotNo             string     `json:"lot_no,omitempty"`
	ExpirationDate    *time.Time `json:"expiration_date,omitempty"`
	Quantity          string     `json:"quantity"`
	ReservedQuantity  string     `json:"reserved_quantity"`
	AvailableQuantity string     `json:"available_quantity"`
}

type BulkReservation struct {
	ID                     string     `json:"id"`
	OrganizationID         string     `json:"organization_id"`
	BinID                  string     `json:"bin_id"`
	ItemID                 string     `json:"item_id"`
	LotID                  string     `json:"lot_id,omitempty"`
	DemandKind             string     `json:"demand_kind"`
	DemandID               string     `json:"demand_id"`
	DemandLineID           string     `json:"demand_line_id"`
	Quantity               string     `json:"quantity"`
	CancellationDisallowed bool       `json:"cancellation_disallowed"`
	ExpiresAt              *time.Time `json:"expires_at,omitempty"`
	Status                 string     `json:"status"`
	Version                int64      `json:"version"`
}

type BulkMovement struct {
	OrganizationID string    `json:"organization_id"`
	FromBinID      string    `json:"from_bin_id"`
	ToBinID        string    `json:"to_bin_id"`
	ItemID         string    `json:"item_id"`
	LotID          string    `json:"lot_id,omitempty"`
	Quantity       string    `json:"quantity"`
	PostingDate    time.Time `json:"posting_date"`
	SourceID       string    `json:"source_id"`
}

type BulkIssue struct {
	OrganizationID       string    `json:"organization_id"`
	ReservationID        string    `json:"reservation_id"`
	ReservationVersion   int64     `json:"reservation_version"`
	SpecificReceiptEntry string    `json:"specific_receipt_entry,omitempty"`
	PostingDate          time.Time `json:"posting_date"`
	SourceKind           string    `json:"source_kind"`
	SourceID             string    `json:"source_id"`
}

type BulkIssueResult struct {
	EntryID      string `json:"entry_id"`
	Quantity     string `json:"quantity"`
	CostAmount   string `json:"cost_amount"`
	Applications int    `json:"applications"`
}

type BulkRepository interface {
	CreateBulkItem(context.Context, string, string, BulkItem) (BulkItem, error)
	ConfigureItemUnitOfMeasure(context.Context, string, string, ItemUnitOfMeasure) (ItemUnitOfMeasure, error)
	ConvertItemUnitOfMeasure(context.Context, string, string, string, string) (UnitOfMeasureConversion, error)
	ConvertBulkHandlingUnits(context.Context, string, PackagingConversionIDs, PackagingConversionCommand) (PackagingConversionResult, error)
	CreateWarehouseBin(context.Context, string, string, WarehouseBin) (WarehouseBin, error)
	ConfigureItemBin(context.Context, string, string, ItemBinPolicy) (ItemBinPolicy, error)
	ReceiveBulk(context.Context, string, string, string, string, string, BulkReceipt) (BulkReceiptResult, error)
	BulkAvailability(context.Context, string, string, string) ([]BulkAvailability, error)
	ReserveBulk(context.Context, string, string, BulkReservation) (BulkReservation, error)
	ReleaseBulk(context.Context, string, string, string, int64, string) error
	MoveBulk(context.Context, string, string, string, BulkMovement) error
	IssueBulk(context.Context, string, string, string, BulkIssue) (BulkIssueResult, error)
}

type BulkService struct {
	repository BulkRepository
	ids        IDGenerator
}

func NewBulkService(repository BulkRepository, ids IDGenerator) *BulkService {
	return &BulkService{repository: repository, ids: ids}
}

func positiveDecimal(value string, pattern *regexp.Regexp) bool {
	if !pattern.MatchString(value) {
		return false
	}
	n, ok := new(big.Rat).SetString(value)
	return ok && n.Sign() > 0
}

func nonNegativeDecimal(value string) bool {
	if !quantityPattern.MatchString(value) {
		return false
	}
	n, ok := new(big.Rat).SetString(value)
	return ok && n.Sign() >= 0
}

func (s *BulkService) CreateItem(ctx context.Context, tenant string, value BulkItem) (BulkItem, error) {
	tracking := map[string]bool{"none": true, "lot": true}
	costing := map[string]bool{"fifo": true, "specific": true}
	if tenant == "" || value.Code == "" || value.Description == "" || value.BaseUOM == "" || !positiveDecimal(value.BaseRoundingPrecision, quantityPattern) || !tracking[value.TrackingMode] || !costing[value.CostingMethod] {
		return BulkItem{}, fmt.Errorf("invalid bulk item")
	}
	value.ID, value.Version = s.ids.New(), 1
	return s.repository.CreateBulkItem(ctx, tenant, s.ids.New(), value)
}

func (s *BulkService) CreateBin(ctx context.Context, tenant string, value WarehouseBin) (WarehouseBin, error) {
	types := map[string]bool{"receive": true, "ship": true, "put-away": true, "pick": true, "putpick": true, "qc": true}
	if tenant == "" || value.OrganizationID == "" || value.Code == "" || !types[value.Type] || value.Ranking < 0 || value.Ranking > 1_000_000 {
		return WarehouseBin{}, fmt.Errorf("invalid warehouse bin")
	}
	value.ID, value.Version = s.ids.New(), 1
	return s.repository.CreateWarehouseBin(ctx, tenant, s.ids.New(), value)
}

func (s *BulkService) ConfigureBin(ctx context.Context, tenant string, value ItemBinPolicy) (ItemBinPolicy, error) {
	if tenant == "" || value.OrganizationID == "" || value.ItemID == "" || value.BinID == "" || !nonNegativeDecimal(value.MinQuantity) || (value.MaxQuantity != "" && !positiveDecimal(value.MaxQuantity, quantityPattern)) {
		return ItemBinPolicy{}, fmt.Errorf("invalid item bin policy")
	}
	value.Version = 1
	return s.repository.ConfigureItemBin(ctx, tenant, s.ids.New(), value)
}

func (s *BulkService) Receive(ctx context.Context, tenant string, value BulkReceipt) (BulkReceiptResult, error) {
	if tenant == "" || value.OrganizationID == "" || value.BinID == "" || value.ItemID == "" || !positiveDecimal(value.Quantity, quantityPattern) || !positiveDecimal(value.UnitCost, costPattern) || value.PostingDate.IsZero() || value.SourceKind == "" || value.SourceID == "" {
		return BulkReceiptResult{}, fmt.Errorf("invalid bulk receipt")
	}
	value.PostingDate = value.PostingDate.UTC()
	return s.repository.ReceiveBulk(ctx, tenant, s.ids.New(), s.ids.New(), s.ids.New(), s.ids.New(), value)
}

func (s *BulkService) Availability(ctx context.Context, tenant, organization, item string) ([]BulkAvailability, error) {
	if tenant == "" || organization == "" || item == "" {
		return nil, fmt.Errorf("invalid bulk availability query")
	}
	return s.repository.BulkAvailability(ctx, tenant, organization, item)
}

func (s *BulkService) Reserve(ctx context.Context, tenant string, value BulkReservation) (BulkReservation, error) {
	allowed := map[string]bool{"service": true, "manual": true, "transfer-outbound": true}
	if tenant == "" || value.OrganizationID == "" || value.BinID == "" || value.ItemID == "" || !allowed[value.DemandKind] || value.DemandID == "" || value.DemandLineID == "" || !positiveDecimal(value.Quantity, quantityPattern) {
		return BulkReservation{}, fmt.Errorf("invalid bulk reservation")
	}
	value.ID, value.Status, value.Version = s.ids.New(), "reservation", 1
	return s.repository.ReserveBulk(ctx, tenant, s.ids.New(), value)
}

func (s *BulkService) Release(ctx context.Context, tenant, organization, reservation string, version int64) error {
	if tenant == "" || organization == "" || reservation == "" || version < 1 {
		return fmt.Errorf("invalid bulk release")
	}
	return s.repository.ReleaseBulk(ctx, tenant, organization, reservation, version, s.ids.New())
}

func (s *BulkService) Move(ctx context.Context, tenant string, value BulkMovement) error {
	if tenant == "" || value.OrganizationID == "" || value.FromBinID == "" || value.ToBinID == "" || value.FromBinID == value.ToBinID || value.ItemID == "" || !positiveDecimal(value.Quantity, quantityPattern) || value.PostingDate.IsZero() || value.SourceID == "" {
		return fmt.Errorf("invalid bulk movement")
	}
	value.PostingDate = value.PostingDate.UTC()
	return s.repository.MoveBulk(ctx, tenant, s.ids.New(), s.ids.New(), value)
}

func (s *BulkService) Issue(ctx context.Context, tenant string, value BulkIssue) (BulkIssueResult, error) {
	if tenant == "" || value.OrganizationID == "" || value.ReservationID == "" || value.ReservationVersion < 1 || value.PostingDate.IsZero() || value.SourceKind == "" || value.SourceID == "" {
		return BulkIssueResult{}, fmt.Errorf("invalid bulk issue")
	}
	value.PostingDate = value.PostingDate.UTC()
	return s.repository.IssueBulk(ctx, tenant, s.ids.New(), s.ids.New(), value)
}
````

### FILE: `internal/inventorycontrol/bulk_test.go`

```yaml
block_id: "GO-OPS-API:bulk-inventory-service-test:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps-supported costing/tracking surface; local fail-closed contract regression"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "a2f68a08e2af19793f51b21ee1966c07acd2aad62cbea8f4eff648586af691f2"
variables: []
secrets_allowed: false
```

````go
package inventorycontrol

import (
	"context"
	"testing"
	"time"
)

type bulkFake struct{ created int }

func (*bulkFake) ConfigureItemUnitOfMeasure(context.Context, string, string, ItemUnitOfMeasure) (ItemUnitOfMeasure, error) {
	return ItemUnitOfMeasure{}, nil
}
func (*bulkFake) ConvertItemUnitOfMeasure(context.Context, string, string, string, string) (UnitOfMeasureConversion, error) {
	return UnitOfMeasureConversion{}, nil
}
func (*bulkFake) ConvertBulkHandlingUnits(context.Context, string, PackagingConversionIDs, PackagingConversionCommand) (PackagingConversionResult, error) {
	return PackagingConversionResult{}, nil
}

func (f *bulkFake) CreateBulkItem(_ context.Context, _ string, _ string, v BulkItem) (BulkItem, error) {
	f.created++
	return v, nil
}
func (*bulkFake) CreateWarehouseBin(context.Context, string, string, WarehouseBin) (WarehouseBin, error) {
	return WarehouseBin{}, nil
}
func (*bulkFake) ConfigureItemBin(context.Context, string, string, ItemBinPolicy) (ItemBinPolicy, error) {
	return ItemBinPolicy{}, nil
}
func (*bulkFake) ReceiveBulk(context.Context, string, string, string, string, string, BulkReceipt) (BulkReceiptResult, error) {
	return BulkReceiptResult{}, nil
}
func (*bulkFake) BulkAvailability(context.Context, string, string, string) ([]BulkAvailability, error) {
	return nil, nil
}
func (*bulkFake) ReserveBulk(context.Context, string, string, BulkReservation) (BulkReservation, error) {
	return BulkReservation{}, nil
}
func (*bulkFake) ReleaseBulk(context.Context, string, string, string, int64, string) error {
	return nil
}
func (*bulkFake) MoveBulk(context.Context, string, string, string, BulkMovement) error { return nil }
func (*bulkFake) IssueBulk(context.Context, string, string, string, BulkIssue) (BulkIssueResult, error) {
	return BulkIssueResult{}, nil
}

func TestBulkServiceRejectsUnsupportedCostAndSerialDuplication(t *testing.T) {
	fake := &bulkFake{}
	service := NewBulkService(fake, inventoryIDsForBulk{})
	for _, value := range []BulkItem{
		{Code: "PART", Description: "Part", BaseUOM: "EA", TrackingMode: "serial", CostingMethod: "specific"},
		{Code: "PART", Description: "Part", BaseUOM: "EA", TrackingMode: "lot", CostingMethod: "average"},
		{Code: "PART", Description: "Part", BaseUOM: "EA", TrackingMode: "lot", CostingMethod: "lifo"},
	} {
		if _, err := service.CreateItem(context.Background(), "tenant", value); err == nil {
			t.Fatalf("unsupported item admitted: %+v", value)
		}
	}
	if fake.created != 0 {
		t.Fatalf("repository called %d times", fake.created)
	}
	if _, err := service.Receive(context.Background(), "tenant", BulkReceipt{OrganizationID: "o", BinID: "b", ItemID: "i", Quantity: "1.0000001", UnitCost: "1", PostingDate: time.Now(), SourceKind: "purchase", SourceID: "p"}); err == nil {
		t.Fatal("over-precision quantity admitted")
	}
}

func TestBulkServiceValidatesPackagingConversionBeforeRepository(t *testing.T) {
	service := NewBulkService(&bulkFake{}, inventoryIDsForBulk{})
	valid := PackagingConversionCommand{RequestID: "request-1", OrganizationID: "warehouse", BinID: "pick", ItemID: "part", Operation: "breakbulk", FromUOM: "BOX", ToUOM: "EA", FromQuantity: "1"}
	if _, err := service.ConvertHandlingUnits(context.Background(), "tenant", valid); err != nil {
		t.Fatal(err)
	}
	for _, invalid := range []PackagingConversionCommand{
		{},
		{RequestID: "r", OrganizationID: "o", BinID: "b", ItemID: "i", Operation: "invented", FromUOM: "BOX", ToUOM: "EA", FromQuantity: "1"},
		{RequestID: "r", OrganizationID: "o", BinID: "b", ItemID: "i", Operation: "gather", FromUOM: "EA", ToUOM: "EA", FromQuantity: "1"},
		{RequestID: "r", OrganizationID: "o", BinID: "b", ItemID: "i", Operation: "gather", FromUOM: "EA", ToUOM: "BOX", FromQuantity: "0"},
	} {
		if _, err := service.ConvertHandlingUnits(context.Background(), "tenant", invalid); err == nil {
			t.Fatalf("invalid packaging conversion admitted: %+v", invalid)
		}
	}
}

type inventoryIDsForBulk struct{ n int }

func (g inventoryIDsForBulk) New() string { return "018f4d4a-7b36-7a21-8d10-2f4c54c28888" }
````

### FILE: `internal/platform/postgres/bulk_inventory.go`

```yaml
block_id: "GO-OPS-API:bulk-inventory-postgres:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps 2eae56d warehouse availability and item application invariants; local PostgreSQL transaction port"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "d4147e29c8f2c141de02624a07493d3fea8ddd50f932c85171f49cbe7eb4e754"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func bulkConflict(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23505" || pgErr.Code == "23514" || pgErr.Code == "23503" || pgErr.Code == "55000" || pgErr.Code == "40001") {
		return inventorycontrol.ErrConflict
	}
	return err
}

func recordBulkEvent(ctx context.Context, tx pgx.Tx, tenant, eventID, aggregateType, aggregateID, eventType string, version int64, payload any) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,$3,$4,$5,$6,1,clock_timestamp(),$7)`, tenant, eventID, aggregateType, aggregateID, version, eventType, encoded)
	return err
}

func (r *InventoryControl) CreateBulkItem(ctx context.Context, tenant, eventID string, value inventorycontrol.BulkItem) (inventorycontrol.BulkItem, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into inventory.stock_item(tenant_id,item_id,item_code,description,base_uom,tracking_mode,costing_method,version) values($1,$2,$3,$4,$5,$6,$7,1)`, tenant, value.ID, value.Code, value.Description, value.BaseUOM, value.TrackingMode, value.CostingMethod)
	if err != nil {
		return value, bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.item_unit_of_measure(tenant_id,item_id,uom_code,qty_per_uom,rounding_precision,is_base,version) values($1,$2,$3,1,$4::numeric,true,1)`, tenant, value.ID, value.BaseUOM, value.BaseRoundingPrecision)
	if err != nil {
		return value, bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-item", value.ID, "bulk-item.created", 1, map[string]string{"item_code": value.Code, "base_uom": value.BaseUOM, "base_rounding_precision": value.BaseRoundingPrecision, "tracking_mode": value.TrackingMode, "costing_method": value.CostingMethod}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func (r *InventoryControl) CreateWarehouseBin(ctx context.Context, tenant, eventID string, value inventorycontrol.WarehouseBin) (inventorycontrol.WarehouseBin, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_bin(tenant_id,organization_id,bin_id,bin_code,bin_type,bin_rank,movement_blocked,cross_dock,version) values($1,$2,$3,$4,$5,$6,$7,$8,1)`, tenant, value.OrganizationID, value.ID, value.Code, value.Type, value.Ranking, value.MovementBlocked, value.CrossDock)
	if err != nil {
		return value, bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "warehouse-bin", value.ID, "warehouse-bin.created", 1, map[string]any{"organization_id": value.OrganizationID, "bin_code": value.Code, "bin_type": value.Type, "bin_rank": value.Ranking, "cross_dock": value.CrossDock}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func (r *InventoryControl) ConfigureItemBin(ctx context.Context, tenant, eventID string, value inventorycontrol.ItemBinPolicy) (inventorycontrol.ItemBinPolicy, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into inventory.item_bin_policy(tenant_id,organization_id,item_id,bin_id,fixed,dedicated,is_default,min_quantity,max_quantity,version) values($1,$2,$3,$4,$5,$6,$7,$8::numeric,nullif($9,'')::numeric,1)`, tenant, value.OrganizationID, value.ItemID, value.BinID, value.Fixed, value.Dedicated, value.Default, value.MinQuantity, value.MaxQuantity)
	if err != nil {
		return value, bulkConflict(err)
	}
	aggregate := value.OrganizationID + ":" + value.ItemID + ":" + value.BinID
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "item-bin-policy", aggregate, "item-bin-policy.created", 1, map[string]any{"organization_id": value.OrganizationID, "item_id": value.ItemID, "bin_id": value.BinID, "fixed": value.Fixed, "dedicated": value.Dedicated, "default": value.Default}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func (r *InventoryControl) ReceiveBulk(ctx context.Context, tenant, eventID, entryID, layerID, proposedLotID string, value inventorycontrol.BulkReceipt) (inventorycontrol.BulkReceiptResult, error) {
	result := inventorycontrol.BulkReceiptResult{EntryID: entryID, Quantity: value.Quantity, UnitCost: value.UnitCost}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return result, err
	}
	defer tx.Rollback(ctx)
	var tracking string
	var movementBlocked bool
	err = tx.QueryRow(ctx, `select i.tracking_mode,w.movement_blocked from inventory.stock_item i join inventory.item_bin_policy p on p.tenant_id=i.tenant_id and p.item_id=i.item_id join inventory.warehouse_bin w on w.tenant_id=p.tenant_id and w.organization_id=p.organization_id and w.bin_id=p.bin_id where i.tenant_id=$1 and i.item_id=$2 and p.organization_id=$3 and p.bin_id=$4 for share of i,w`, tenant, value.ItemID, value.OrganizationID, value.BinID).Scan(&tracking, &movementBlocked)
	if errors.Is(err, pgx.ErrNoRows) || movementBlocked {
		return result, inventorycontrol.ErrConflict
	}
	if err != nil {
		return result, err
	}
	lotID := ""
	if tracking == "lot" {
		if value.LotNo == "" {
			return result, inventorycontrol.ErrConflict
		}
		var existingExpiration, existingWarranty string
		err = tx.QueryRow(ctx, `select lot_id,coalesce(expiration_date::text,''),coalesce(warranty_date::text,'') from inventory.inventory_lot where tenant_id=$1 and item_id=$2 and lot_no=$3 for update`, tenant, value.ItemID, value.LotNo).Scan(&lotID, &existingExpiration, &existingWarranty)
		expiration, warranty := "", ""
		if value.ExpirationDate != nil {
			expiration = value.ExpirationDate.UTC().Format("2006-01-02")
		}
		if value.WarrantyDate != nil {
			warranty = value.WarrantyDate.UTC().Format("2006-01-02")
		}
		if errors.Is(err, pgx.ErrNoRows) {
			lotID = proposedLotID
			_, err = tx.Exec(ctx, `insert into inventory.inventory_lot(tenant_id,lot_id,item_id,lot_no,expiration_date,warranty_date) values($1,$2,$3,$4,nullif($5,'')::date,nullif($6,'')::date)`, tenant, lotID, value.ItemID, value.LotNo, expiration, warranty)
		} else if err == nil && (existingExpiration != expiration || existingWarranty != warranty) {
			return result, inventorycontrol.ErrConflict
		}
		if err != nil {
			return result, bulkConflict(err)
		}
	} else if value.LotNo != "" || value.ExpirationDate != nil || value.WarrantyDate != nil {
		return result, inventorycontrol.ErrConflict
	}
	result.LotID = lotID
	var balanceVersion int64
	err = tx.QueryRow(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity,version) values($1,$2,$3,$4,nullif($5,''),$6::numeric,0,1) on conflict (tenant_id,organization_id,bin_id,item_id,lot_id) do update set quantity=inventory.bulk_balance.quantity+excluded.quantity,version=inventory.bulk_balance.version+1,updated_at=clock_timestamp() where (select max_quantity is null or inventory.bulk_balance.quantity+excluded.quantity<=max_quantity from inventory.item_bin_policy where tenant_id=$1 and organization_id=$2 and item_id=$4 and bin_id=$3) returning version`, tenant, value.OrganizationID, value.BinID, value.ItemID, lotID, value.Quantity).Scan(&balanceVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, inventorycontrol.ErrConflict
	}
	if err != nil {
		return result, bulkConflict(err)
	}
	err = tx.QueryRow(ctx, `insert into inventory.bulk_inventory_entry(tenant_id,entry_id,organization_id,bin_id,item_id,lot_id,entry_type,quantity,unit_cost,cost_amount,posting_date,source_kind,source_id) values($1,$2,$3,$4,$5,nullif($6,''),'receipt',$7::numeric,$8::numeric,round($7::numeric*$8::numeric,4),$9::date,$10,$11) returning cost_amount::text`, tenant, entryID, value.OrganizationID, value.BinID, value.ItemID, lotID, value.Quantity, value.UnitCost, value.PostingDate, value.SourceKind, value.SourceID).Scan(&result.CostAmount)
	if err != nil {
		return result, bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_cost_layer(tenant_id,layer_id,receipt_entry_id,organization_id,item_id,lot_id,posting_date,original_quantity,remaining_quantity,unit_cost) values($1,$2,$3,$4,$5,nullif($6,''),$7::date,$8::numeric,$8::numeric,$9::numeric)`, tenant, layerID, entryID, value.OrganizationID, value.ItemID, lotID, value.PostingDate, value.Quantity, value.UnitCost)
	if err != nil {
		return result, bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-inventory-entry", entryID, "bulk-inventory.received", 1, map[string]string{"organization_id": value.OrganizationID, "bin_id": value.BinID, "item_id": value.ItemID, "lot_id": lotID, "quantity": value.Quantity, "cost_amount": result.CostAmount}); err != nil {
		return result, err
	}
	return result, tx.Commit(ctx)
}

func (r *InventoryControl) BulkAvailability(ctx context.Context, tenant, organization, item string) ([]inventorycontrol.BulkAvailability, error) {
	rows, err := r.pool.Query(ctx, `select organization_id,bin_id,bin_code,item_id,coalesce(lot_id,''),coalesce(lot_no,''),expiration_date,quantity::text,reserved_quantity::text,available_quantity::text from inventory.bulk_available where tenant_id=$1 and organization_id=$2 and item_id=$3 and available_quantity>0 order by expiration_date nulls last,lot_no,bin_code`, tenant, organization, item)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	values := []inventorycontrol.BulkAvailability{}
	for rows.Next() {
		var value inventorycontrol.BulkAvailability
		if err = rows.Scan(&value.OrganizationID, &value.BinID, &value.BinCode, &value.ItemID, &value.LotID, &value.LotNo, &value.ExpirationDate, &value.Quantity, &value.ReservedQuantity, &value.AvailableQuantity); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, rows.Err()
}

func (r *InventoryControl) ReserveBulk(ctx context.Context, tenant, eventID string, value inventorycontrol.BulkReservation) (inventorycontrol.BulkReservation, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	var balanceVersion int64
	err = tx.QueryRow(ctx, `update inventory.bulk_balance b set reserved_quantity=reserved_quantity+$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and quantity-reserved_quantity >= $6::numeric and exists(select 1 from inventory.item_bin_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) left join inventory.inventory_lot l on l.tenant_id=b.tenant_id and l.lot_id=b.lot_id where p.tenant_id=b.tenant_id and p.organization_id=b.organization_id and p.item_id=b.item_id and p.bin_id=b.bin_id and not p.dedicated and not w.movement_blocked and coalesce(l.blocked,false)=false and (l.expiration_date is null or l.expiration_date>=current_date)) returning version`, tenant, value.OrganizationID, value.BinID, value.ItemID, value.LotID, value.Quantity).Scan(&balanceVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, inventorycontrol.ErrConflict
	}
	if err != nil {
		return value, err
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_reservation(tenant_id,reservation_id,organization_id,bin_id,item_id,lot_id,demand_kind,demand_id,demand_line_id,quantity,status,cancellation_disallowed,expires_at,version) values($1,$2,$3,$4,$5,nullif($6,''),$7,$8,$9,$10::numeric,'reservation',$11,$12,1)`, tenant, value.ID, value.OrganizationID, value.BinID, value.ItemID, value.LotID, value.DemandKind, value.DemandID, value.DemandLineID, value.Quantity, value.CancellationDisallowed, value.ExpiresAt)
	if err != nil {
		return value, bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-reservation", value.ID, "bulk-reservation.created", 1, map[string]string{"organization_id": value.OrganizationID, "bin_id": value.BinID, "item_id": value.ItemID, "lot_id": value.LotID, "quantity": value.Quantity}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func (r *InventoryControl) ReleaseBulk(ctx context.Context, tenant, organization, reservationID string, version int64, eventID string) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var binID, itemID, lotID, quantity string
	err = tx.QueryRow(ctx, `update inventory.bulk_reservation set status='released',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and reservation_id=$3 and status='reservation' and version=$4 and not cancellation_disallowed returning bin_id,item_id,coalesce(lot_id,''),quantity::text`, tenant, organization, reservationID, version).Scan(&binID, &itemID, &lotID, &quantity)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.ErrConflict
	}
	if err != nil {
		return err
	}
	result, err := tx.Exec(ctx, `update inventory.bulk_balance set reserved_quantity=reserved_quantity-$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and reserved_quantity >= $6::numeric`, tenant, organization, binID, itemID, lotID, quantity)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return inventorycontrol.ErrConflict
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-reservation", reservationID, "bulk-reservation.released", version+1, map[string]string{"quantity": quantity}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *InventoryControl) MoveBulk(ctx context.Context, tenant, eventOutID, eventInID string, value inventorycontrol.BulkMovement) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update inventory.bulk_balance set quantity=quantity-$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and quantity-reserved_quantity >= $6::numeric and exists(select 1 from inventory.warehouse_bin where tenant_id=$1 and organization_id=$2 and bin_id=$3 and not movement_blocked)`, tenant, value.OrganizationID, value.FromBinID, value.ItemID, value.LotID, value.Quantity)
	if err != nil {
		return bulkConflict(err)
	}
	if result.RowsAffected() != 1 {
		return inventorycontrol.ErrConflict
	}
	var targetVersion int64
	err = tx.QueryRow(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity,version) values($1,$2,$3,$4,nullif($5,''),$6::numeric,0,1) on conflict (tenant_id,organization_id,bin_id,item_id,lot_id) do update set quantity=inventory.bulk_balance.quantity+excluded.quantity,version=inventory.bulk_balance.version+1,updated_at=clock_timestamp() where exists(select 1 from inventory.item_bin_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) where p.tenant_id=$1 and p.organization_id=$2 and p.item_id=$4 and p.bin_id=$3 and not w.movement_blocked and (p.max_quantity is null or inventory.bulk_balance.quantity+excluded.quantity<=p.max_quantity)) returning version`, tenant, value.OrganizationID, value.ToBinID, value.ItemID, value.LotID, value.Quantity).Scan(&targetVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.ErrConflict
	}
	if err != nil {
		return bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_inventory_entry(tenant_id,entry_id,organization_id,bin_id,item_id,lot_id,entry_type,quantity,unit_cost,cost_amount,posting_date,source_kind,source_id) values($1,$2,$3,$4,$5,nullif($6,''),'movement-out',-$7::numeric,0,0,$8::date,'bin-movement',$9),($1,$10,$3,$11,$5,nullif($6,''),'movement-in',$7::numeric,0,0,$8::date,'bin-movement',$9)`, tenant, eventOutID, value.OrganizationID, value.FromBinID, value.ItemID, value.LotID, value.Quantity, value.PostingDate, value.SourceID, eventInID, value.ToBinID)
	if err != nil {
		return bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventOutID, "bulk-bin-movement", value.SourceID, "bulk-inventory.moved", 1, map[string]string{"organization_id": value.OrganizationID, "from_bin_id": value.FromBinID, "to_bin_id": value.ToBinID, "item_id": value.ItemID, "lot_id": value.LotID, "quantity": value.Quantity}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

type bulkAllocation struct {
	layerID, inboundEntryID, quantity, unitCost, costAmount string
}

func parseRat(value string) (*big.Rat, bool) { return new(big.Rat).SetString(value) }

func minRat(a, b *big.Rat) *big.Rat {
	if a.Cmp(b) <= 0 {
		return new(big.Rat).Set(a)
	}
	return new(big.Rat).Set(b)
}

func formatRat(value *big.Rat, scale int) string {
	factor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(scale)), nil)
	numerator := new(big.Int).Mul(value.Num(), factor)
	quotient, remainder := new(big.Int), new(big.Int)
	quotient.QuoRem(numerator, value.Denom(), remainder)
	if new(big.Int).Mul(remainder, big.NewInt(2)).Cmp(value.Denom()) >= 0 {
		quotient.Add(quotient, big.NewInt(1))
	}
	digits := quotient.String()
	if scale == 0 {
		return digits
	}
	for len(digits) <= scale {
		digits = "0" + digits
	}
	whole, fraction := digits[:len(digits)-scale], strings.TrimRight(digits[len(digits)-scale:], "0")
	if fraction == "" {
		return whole
	}
	return whole + "." + fraction
}

func (r *InventoryControl) IssueBulk(ctx context.Context, tenant, eventID, outboundEntryID string, value inventorycontrol.BulkIssue) (inventorycontrol.BulkIssueResult, error) {
	result := inventorycontrol.BulkIssueResult{EntryID: outboundEntryID}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return result, err
	}
	defer tx.Rollback(ctx)
	var binID, itemID, lotID, quantity, costing string
	err = tx.QueryRow(ctx, `select r.bin_id,r.item_id,coalesce(r.lot_id,''),r.quantity::text,i.costing_method from inventory.bulk_reservation r join inventory.stock_item i using(tenant_id,item_id) left join inventory.inventory_lot l on l.tenant_id=r.tenant_id and l.lot_id=r.lot_id where r.tenant_id=$1 and r.organization_id=$2 and r.reservation_id=$3 and r.status='reservation' and r.version=$4 and (r.expires_at is null or r.expires_at>clock_timestamp()) and coalesce(l.blocked,false)=false and (l.expiration_date is null or l.expiration_date>=current_date) for update of r`, tenant, value.OrganizationID, value.ReservationID, value.ReservationVersion).Scan(&binID, &itemID, &lotID, &quantity, &costing)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, inventorycontrol.ErrConflict
	}
	if err != nil {
		return result, err
	}
	if (costing == "specific") != (value.SpecificReceiptEntry != "") {
		return result, inventorycontrol.ErrConflict
	}
	wanted, ok := parseRat(quantity)
	if !ok {
		return result, inventorycontrol.ErrConflict
	}
	query := `select layer_id,receipt_entry_id,remaining_quantity::text,unit_cost::text from inventory.bulk_cost_layer where tenant_id=$1 and organization_id=$2 and item_id=$3 and lot_id is not distinct from nullif($4,'') and remaining_quantity>0`
	args := []any{tenant, value.OrganizationID, itemID, lotID}
	if costing == "specific" {
		query += ` and receipt_entry_id=$5`
		args = append(args, value.SpecificReceiptEntry)
	}
	query += ` order by posting_date,receipt_entry_id for update`
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return result, err
	}
	remaining := new(big.Rat).Set(wanted)
	applications := []bulkAllocation{}
	totalCost := new(big.Rat)
	for rows.Next() && remaining.Sign() > 0 {
		var layerID, inbound, availableText, unitCost string
		if err = rows.Scan(&layerID, &inbound, &availableText, &unitCost); err != nil {
			rows.Close()
			return result, err
		}
		available, validAvailable := parseRat(availableText)
		unit, validUnit := parseRat(unitCost)
		if !validAvailable || !validUnit {
			rows.Close()
			return result, inventorycontrol.ErrConflict
		}
		take := minRat(available, remaining)
		cost := new(big.Rat).Mul(take, unit)
		costText := formatRat(cost, 4)
		rounded, _ := parseRat(costText)
		totalCost.Add(totalCost, rounded)
		applications = append(applications, bulkAllocation{layerID: layerID, inboundEntryID: inbound, quantity: formatRat(take, 6), unitCost: unitCost, costAmount: costText})
		remaining.Sub(remaining, take)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return result, err
	}
	if remaining.Sign() != 0 || len(applications) == 0 {
		return result, inventorycontrol.ErrConflict
	}
	result.Quantity, result.CostAmount, result.Applications = quantity, formatRat(totalCost, 4), len(applications)
	unitCost := formatRat(new(big.Rat).Quo(totalCost, wanted), 4)
	_, err = tx.Exec(ctx, `insert into inventory.bulk_inventory_entry(tenant_id,entry_id,organization_id,bin_id,item_id,lot_id,entry_type,quantity,unit_cost,cost_amount,posting_date,source_kind,source_id) values($1,$2,$3,$4,$5,nullif($6,''),'issue',-$7::numeric,$8::numeric,-$9::numeric,$10::date,$11,$12)`, tenant, outboundEntryID, value.OrganizationID, binID, itemID, lotID, quantity, unitCost, result.CostAmount, value.PostingDate, value.SourceKind, value.SourceID)
	if err != nil {
		return result, bulkConflict(err)
	}
	for _, application := range applications {
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_cost_layer set remaining_quantity=remaining_quantity-$3::numeric where tenant_id=$1 and layer_id=$2 and remaining_quantity >= $3::numeric`, tenant, application.layerID, application.quantity)
		if updateErr != nil {
			return result, updateErr
		}
		if updated.RowsAffected() != 1 {
			return result, inventorycontrol.ErrConflict
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_cost_application(tenant_id,outbound_entry_id,inbound_entry_id,quantity,cost_amount) values($1,$2,$3,$4::numeric,$5::numeric)`, tenant, outboundEntryID, application.inboundEntryID, application.quantity, application.costAmount)
		if err != nil {
			return result, bulkConflict(err)
		}
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_balance set quantity=quantity-$6::numeric,reserved_quantity=reserved_quantity-$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and quantity >= $6::numeric and reserved_quantity >= $6::numeric`, tenant, value.OrganizationID, binID, itemID, lotID, quantity)
	if err != nil {
		return result, err
	}
	if updated.RowsAffected() != 1 {
		return result, inventorycontrol.ErrConflict
	}
	updated, err = tx.Exec(ctx, `update inventory.bulk_reservation set status='consumed',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and reservation_id=$2 and status='reservation' and version=$3`, tenant, value.ReservationID, value.ReservationVersion)
	if err != nil {
		return result, err
	}
	if updated.RowsAffected() != 1 {
		return result, inventorycontrol.ErrConflict
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-inventory-entry", outboundEntryID, "bulk-inventory.issued", 1, map[string]string{"organization_id": value.OrganizationID, "item_id": itemID, "lot_id": lotID, "quantity": quantity, "cost_amount": result.CostAmount, "reservation_id": value.ReservationID}); err != nil {
		return result, err
	}
	return result, tx.Commit(ctx)
}
````

### FILE: `internal/platform/postgres/bulk_inventory_integration_test.go`

```yaml
block_id: "GO-OPS-API:bulk-inventory-postgres-integration-test:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps 2eae56d tracking/warehouse/cost application regression surface; local concurrency and immutability tests"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "31911e5fa2705ff97acb9a7aca9975ee87f116b864c21e094f1e1e821ab8f433"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func bulkTestUUID(t *testing.T) string {
	t.Helper()
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		t.Fatal(err)
	}
	value[6] = (value[6] & 0x0f) | 0x40
	value[8] = (value[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", value[0:4], value[4:6], value[6:8], value[8:10], value[10:16])
}

func TestBulkLotBinFIFOConcurrencyAndImmutableApplications(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := bulkTestUUID(t)
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Bulk Inventory','Bulk Inventory')`, tenant, "bulk-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	defer func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			t.Errorf("bulk fixture cleanup begin failed: %v", cleanupErr)
			return
		}
		defer tx.Rollback(ctx)
		commands := []struct {
			sql  string
			args []any
		}{
			{`alter table inventory.bulk_cost_application disable trigger bulk_cost_application_immutable`, nil},
			{`alter table inventory.bulk_inventory_entry disable trigger bulk_inventory_entry_immutable`, nil},
			{`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`, nil},
			{`delete from inventory.bulk_cost_application where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_cost_layer where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_inventory_entry where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_reservation where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_balance where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.inventory_lot where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.item_bin_policy where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.warehouse_bin where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.item_unit_of_measure where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.stock_item where tenant_id=$1`, []any{tenant}},
			{`delete from platform.outbox_event where tenant_id=$1`, []any{tenant}},
			{`delete from org.organization where tenant_id=$1`, []any{tenant}},
			{`delete from platform.tenant where tenant_id=$1`, []any{tenant}},
			{`alter table inventory.bulk_cost_application enable trigger bulk_cost_application_immutable`, nil},
			{`alter table inventory.bulk_inventory_entry enable trigger bulk_inventory_entry_immutable`, nil},
			{`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`, nil},
		}
		for _, command := range commands {
			if _, cleanupErr = tx.Exec(ctx, command.sql, command.args...); cleanupErr != nil {
				t.Errorf("bulk fixture cleanup failed: %v", cleanupErr)
				return
			}
		}
		if cleanupErr = tx.Commit(ctx); cleanupErr != nil {
			t.Errorf("bulk fixture cleanup commit failed: %v", cleanupErr)
		}
	}()
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'warehouse','warehouse','Warehouse','warehouse')`, tenant); err != nil {
		t.Fatal(err)
	}
	repo := NewInventoryControl(pool)
	item := inventorycontrol.BulkItem{ID: "part", Code: "PART-1", Description: "Brake pad", BaseUOM: "EA", BaseRoundingPrecision: "0.000001", TrackingMode: "lot", CostingMethod: "fifo", Version: 1}
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), item); err != nil {
		t.Fatal(err)
	}
	for _, bin := range []inventorycontrol.WarehouseBin{
		{ID: "pick", OrganizationID: "warehouse", Code: "PICK-01", Type: "pick", Version: 1},
		{ID: "ship", OrganizationID: "warehouse", Code: "SHIP-01", Type: "ship", Version: 1},
	} {
		if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatal(err)
		}
		policy := inventorycontrol.ItemBinPolicy{OrganizationID: "warehouse", ItemID: "part", BinID: bin.ID, Fixed: true, Default: bin.ID == "pick", MinQuantity: "0", MaxQuantity: "100", Version: 1}
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), policy); err != nil {
			t.Fatal(err)
		}
	}
	expires := time.Date(2035, 1, 1, 0, 0, 0, 0, time.UTC)
	receipts := []struct {
		entry, layer, lot, quantity, cost string
		date                              time.Time
	}{
		{"receipt-1", "layer-1", "lot-1", "5", "100", time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"receipt-2", "layer-2", "ignored-existing-lot", "5", "120", time.Date(2030, 1, 2, 0, 0, 0, 0, time.UTC)},
	}
	for index, receipt := range receipts {
		value := inventorycontrol.BulkReceipt{OrganizationID: "warehouse", BinID: "pick", ItemID: "part", LotNo: "LOT-2035", ExpirationDate: &expires, Quantity: receipt.quantity, UnitCost: receipt.cost, PostingDate: receipt.date, SourceKind: "purchase-receipt", SourceID: fmt.Sprintf("purchase-%d", index+1)}
		got, receiveErr := repo.ReceiveBulk(ctx, tenant, bulkTestUUID(t), receipt.entry, receipt.layer, receipt.lot, value)
		if receiveErr != nil || got.LotID != "lot-1" {
			t.Fatalf("receipt %d result=%+v err=%v", index, got, receiveErr)
		}
	}
	availability, err := repo.BulkAvailability(ctx, tenant, "warehouse", "part")
	if err != nil || len(availability) != 1 || availability[0].AvailableQuantity != "10.000000" {
		t.Fatalf("availability=%+v err=%v", availability, err)
	}
	reservations := []inventorycontrol.BulkReservation{
		{ID: "reserve-a", OrganizationID: "warehouse", BinID: "pick", ItemID: "part", LotID: "lot-1", DemandKind: "service", DemandID: "service-a", DemandLineID: "line", Quantity: "7", Status: "reservation", Version: 1},
		{ID: "reserve-b", OrganizationID: "warehouse", BinID: "pick", ItemID: "part", LotID: "lot-1", DemandKind: "service", DemandID: "service-b", DemandLineID: "line", Quantity: "7", Status: "reservation", Version: 1},
	}
	errs := make([]error, 2)
	var wg sync.WaitGroup
	for index := range reservations {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, errs[i] = repo.ReserveBulk(ctx, tenant, bulkTestUUID(t), reservations[i])
		}(index)
	}
	wg.Wait()
	winner, successes := "", 0
	for index, reserveErr := range errs {
		if reserveErr == nil {
			winner, successes = reservations[index].ID, successes+1
		} else if !errors.Is(reserveErr, inventorycontrol.ErrConflict) {
			t.Fatalf("unexpected reserve error: %v", reserveErr)
		}
	}
	if successes != 1 {
		t.Fatalf("reservation successes=%d errors=%v", successes, errs)
	}
	issue, err := repo.IssueBulk(ctx, tenant, bulkTestUUID(t), "issue-1", inventorycontrol.BulkIssue{OrganizationID: "warehouse", ReservationID: winner, ReservationVersion: 1, PostingDate: time.Date(2030, 1, 3, 0, 0, 0, 0, time.UTC), SourceKind: "service-use", SourceID: "service-use-1"})
	if err != nil || issue.Quantity != "7.000000" || issue.CostAmount != "740" || issue.Applications != 2 {
		t.Fatalf("issue=%+v err=%v", issue, err)
	}
	var applicationQuantity, applicationCost string
	if err = pool.QueryRow(ctx, `select sum(quantity)::text,sum(cost_amount)::text from inventory.bulk_cost_application where tenant_id=$1 and outbound_entry_id='issue-1'`, tenant).Scan(&applicationQuantity, &applicationCost); err != nil || applicationQuantity != "7.000000" || applicationCost != "740.0000" {
		t.Fatalf("applications quantity=%s cost=%s err=%v", applicationQuantity, applicationCost, err)
	}
	if _, err = pool.Exec(ctx, `update inventory.bulk_inventory_entry set cost_amount=0 where tenant_id=$1 and entry_id='issue-1'`, tenant); err == nil {
		t.Fatal("immutable inventory entry accepted update")
	} else {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || pgErr.Code != "55000" {
			t.Fatalf("unexpected immutable error: %v", err)
		}
	}
	if err = repo.MoveBulk(ctx, tenant, bulkTestUUID(t), bulkTestUUID(t), inventorycontrol.BulkMovement{OrganizationID: "warehouse", FromBinID: "pick", ToBinID: "ship", ItemID: "part", LotID: "lot-1", Quantity: "2", PostingDate: time.Date(2030, 1, 4, 0, 0, 0, 0, time.UTC), SourceID: "movement-1"}); err != nil {
		t.Fatal(err)
	}
	var pickQuantity, shipQuantity string
	if err = pool.QueryRow(ctx, `select quantity::text from inventory.bulk_balance where tenant_id=$1 and bin_id='pick'`, tenant).Scan(&pickQuantity); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select quantity::text from inventory.bulk_balance where tenant_id=$1 and bin_id='ship'`, tenant).Scan(&shipQuantity); err != nil {
		t.Fatal(err)
	}
	if pickQuantity != "1.000000" || shipQuantity != "2.000000" {
		t.Fatalf("movement pick=%s ship=%s", pickQuantity, shipQuantity)
	}

	specific := inventorycontrol.BulkItem{ID: "specific-part", Code: "SPECIFIC-1", Description: "Specific assembly", BaseUOM: "EA", BaseRoundingPrecision: "0.000001", TrackingMode: "none", CostingMethod: "specific", Version: 1}
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), specific); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: "warehouse", ItemID: "specific-part", BinID: "pick", Fixed: true, MinQuantity: "0", MaxQuantity: "10", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.ReceiveBulk(ctx, tenant, bulkTestUUID(t), "specific-receipt", "specific-layer", "unused-lot", inventorycontrol.BulkReceipt{OrganizationID: "warehouse", BinID: "pick", ItemID: "specific-part", Quantity: "1", UnitCost: "999.99", PostingDate: time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC), SourceKind: "purchase-receipt", SourceID: "specific-purchase"}); err != nil {
		t.Fatal(err)
	}
	specificReservation := inventorycontrol.BulkReservation{ID: "specific-reservation", OrganizationID: "warehouse", BinID: "pick", ItemID: "specific-part", DemandKind: "service", DemandID: "specific-service", DemandLineID: "line", Quantity: "1", Status: "reservation", Version: 1}
	if _, err = repo.ReserveBulk(ctx, tenant, bulkTestUUID(t), specificReservation); err != nil {
		t.Fatal(err)
	}
	withoutSpecific := inventorycontrol.BulkIssue{OrganizationID: "warehouse", ReservationID: specificReservation.ID, ReservationVersion: 1, PostingDate: time.Date(2030, 1, 2, 0, 0, 0, 0, time.UTC), SourceKind: "service-use", SourceID: "specific-use"}
	if _, err = repo.IssueBulk(ctx, tenant, bulkTestUUID(t), "specific-issue-rejected", withoutSpecific); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("specific issue without source=%v", err)
	}
	withoutSpecific.SpecificReceiptEntry = "specific-receipt"
	specificIssue, err := repo.IssueBulk(ctx, tenant, bulkTestUUID(t), "specific-issue", withoutSpecific)
	if err != nil || specificIssue.CostAmount != "999.99" || specificIssue.Applications != 1 {
		t.Fatalf("specific issue=%+v err=%v", specificIssue, err)
	}
}
````

### FILE: `internal/platform/httpapi/bulk_inventory.go`

```yaml
block_id: "GO-OPS-API:bulk-inventory-http:v1"
operation: CREATE
provenance: AUTHORED
source: "local authorized HTTP boundary over the admitted inventory owner"
license: "LicenseRef-Workspace-Owner"
sha256: "7d6a0aa7c334ade3885a6061fe279f6d70581339f18dcd10e807a9ac42ba0765"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"errors"
	"net/http"

	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
)

type bulkInventoryAPI struct {
	service  *inventorycontrol.BulkService
	verifier identity.Verifier
}

func registerBulkInventory(mux *http.ServeMux, service *inventorycontrol.BulkService, verifier identity.Verifier) {
	api := bulkInventoryAPI{service: service, verifier: verifier}
	mux.HandleFunc("POST /v1/inventory/bulk/items", api.createItem)
	mux.HandleFunc("POST /v1/inventory/bulk/item-units", api.configureUnitOfMeasure)
	mux.HandleFunc("POST /v1/inventory/bulk/uom-conversions", api.convertUnitOfMeasure)
	mux.HandleFunc("POST /v1/inventory/bulk/handling-unit-conversions", api.convertHandlingUnits)
	mux.HandleFunc("POST /v1/inventory/bulk/bins", api.createBin)
	mux.HandleFunc("POST /v1/inventory/bulk/bin-policies", api.configureBin)
	mux.HandleFunc("POST /v1/inventory/bulk/receipts", api.receive)
	mux.HandleFunc("GET /v1/inventory/bulk/availability", api.availability)
	mux.HandleFunc("POST /v1/inventory/bulk/reservations", api.reserve)
	mux.HandleFunc("POST /v1/inventory/bulk/reservations/{id}/release", api.release)
	mux.HandleFunc("POST /v1/inventory/bulk/movements", api.move)
	mux.HandleFunc("POST /v1/inventory/bulk/issues", api.issue)
}

func (a bulkInventoryAPI) convertHandlingUnits(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.PackagingConversionCommand
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.ConvertHandlingUnits(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_HANDLING_UNIT_CONVERSION") {
		return
	}
	writeJSON(w, 201, value)
}

func (a bulkInventoryAPI) authorize(w http.ResponseWriter, r *http.Request, permission string, requireJSON bool) (identity.Principal, bool) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return p, false
	}
	if requireJSON && r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return p, false
	}
	return p, true
}

func writeBulkError(w http.ResponseWriter, err error, invalidCode string) bool {
	if errors.Is(err, inventorycontrol.ErrConflict) {
		writeProblem(w, 409, "BULK_INVENTORY_CONFLICT", "inventory state, identity, capacity or costing conflict")
		return true
	}
	if err != nil {
		writeProblem(w, 400, invalidCode, "request does not match the bulk inventory contract")
		return true
	}
	return false
}

func (a bulkInventoryAPI) createItem(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.BulkItem
	if !decodeStrict(w, r, &input) {
		return
	}
	value, err := a.service.CreateItem(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_BULK_ITEM") {
		return
	}
	writeJSON(w, 201, value)
}

func (a bulkInventoryAPI) createBin(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.WarehouseBin
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.CreateBin(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_BIN") {
		return
	}
	writeJSON(w, 201, value)
}

func (a bulkInventoryAPI) configureBin(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.ItemBinPolicy
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.ConfigureBin(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_ITEM_BIN_POLICY") {
		return
	}
	writeJSON(w, 201, value)
}

func (a bulkInventoryAPI) receive(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.BulkReceipt
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.Receive(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_BULK_RECEIPT") {
		return
	}
	writeJSON(w, 201, value)
}

func (a bulkInventoryAPI) availability(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:read", false)
	if !ok {
		return
	}
	organization, item := r.URL.Query().Get("organization_id"), r.URL.Query().Get("item_id")
	if !p.AllowedOrganization(organization) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	values, err := a.service.Availability(r.Context(), p.TenantID, organization, item)
	if err != nil {
		writeProblem(w, 400, "INVALID_BULK_AVAILABILITY_QUERY", "query does not match the bulk inventory contract")
		return
	}
	writeJSON(w, 200, map[string]any{"availability": values})
}

func (a bulkInventoryAPI) reserve(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.BulkReservation
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.Reserve(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_BULK_RESERVATION") {
		return
	}
	writeJSON(w, 201, value)
}

func (a bulkInventoryAPI) release(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	err := a.service.Release(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Version)
	if writeBulkError(w, err, "INVALID_BULK_RELEASE") {
		return
	}
	writeJSON(w, 200, map[string]string{"status": "released"})
}

func (a bulkInventoryAPI) move(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.BulkMovement
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	err := a.service.Move(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_BULK_MOVEMENT") {
		return
	}
	writeJSON(w, 200, map[string]string{"status": "moved"})
}

func (a bulkInventoryAPI) issue(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.BulkIssue
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.Issue(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_BULK_ISSUE") {
		return
	}
	writeJSON(w, 201, value)
}
````

### FILE: `internal/platform/httpapi/bulk_inventory_test.go`

```yaml
block_id: "GO-OPS-API:bulk-inventory-http-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local HTTP authorization and unsupported-costing regression"
license: "LicenseRef-Workspace-Owner"
sha256: "689007ed7bfa34d1bfa336911284cc5a107328d1a8f5eb4fd0f940c62075a264"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"elite.local/enterprise/internal/inventorycontrol"
)

type bulkHTTPRepo struct{ items int }

func (*bulkHTTPRepo) ConfigureItemUnitOfMeasure(_ context.Context, _ string, _ string, value inventorycontrol.ItemUnitOfMeasure) (inventorycontrol.ItemUnitOfMeasure, error) {
	return value, nil
}
func (*bulkHTTPRepo) ConvertItemUnitOfMeasure(_ context.Context, _ string, item, code, quantity string) (inventorycontrol.UnitOfMeasureConversion, error) {
	return inventorycontrol.UnitOfMeasureConversion{ItemID: item, Code: code, Quantity: quantity, QuantityPerUnit: "12", BaseCode: "EA", BaseQuantity: "30", RoundingPrecision: "1"}, nil
}
func (*bulkHTTPRepo) ConvertBulkHandlingUnits(_ context.Context, _ string, ids inventorycontrol.PackagingConversionIDs, value inventorycontrol.PackagingConversionCommand) (inventorycontrol.PackagingConversionResult, error) {
	return inventorycontrol.PackagingConversionResult{ID: ids.ConversionID, RequestID: value.RequestID, OrganizationID: value.OrganizationID, Operation: value.Operation, FromUOM: value.FromUOM, ToUOM: value.ToUOM, FromQuantity: value.FromQuantity, ToQuantity: "12.000000", BaseQuantity: "12.000000", Version: 1}, nil
}

func (f *bulkHTTPRepo) CreateBulkItem(_ context.Context, _ string, _ string, v inventorycontrol.BulkItem) (inventorycontrol.BulkItem, error) {
	f.items++
	return v, nil
}
func (*bulkHTTPRepo) CreateWarehouseBin(context.Context, string, string, inventorycontrol.WarehouseBin) (inventorycontrol.WarehouseBin, error) {
	return inventorycontrol.WarehouseBin{}, nil
}
func (*bulkHTTPRepo) ConfigureItemBin(context.Context, string, string, inventorycontrol.ItemBinPolicy) (inventorycontrol.ItemBinPolicy, error) {
	return inventorycontrol.ItemBinPolicy{}, nil
}
func (*bulkHTTPRepo) ReceiveBulk(context.Context, string, string, string, string, string, inventorycontrol.BulkReceipt) (inventorycontrol.BulkReceiptResult, error) {
	return inventorycontrol.BulkReceiptResult{}, nil
}
func (*bulkHTTPRepo) BulkAvailability(context.Context, string, string, string) ([]inventorycontrol.BulkAvailability, error) {
	return []inventorycontrol.BulkAvailability{{OrganizationID: "a", ItemID: "i", AvailableQuantity: "2"}}, nil
}
func (*bulkHTTPRepo) ReserveBulk(context.Context, string, string, inventorycontrol.BulkReservation) (inventorycontrol.BulkReservation, error) {
	return inventorycontrol.BulkReservation{}, nil
}
func (*bulkHTTPRepo) ReleaseBulk(context.Context, string, string, string, int64, string) error {
	return nil
}
func (*bulkHTTPRepo) MoveBulk(context.Context, string, string, string, inventorycontrol.BulkMovement) error {
	return nil
}
func (*bulkHTTPRepo) IssueBulk(context.Context, string, string, string, inventorycontrol.BulkIssue) (inventorycontrol.BulkIssueResult, error) {
	return inventorycontrol.BulkIssueResult{}, nil
}

func TestBulkInventoryHTTPRejectsUnsupportedMethodAndEnforcesOrganization(t *testing.T) {
	repo := &bulkHTTPRepo{}
	mux := http.NewServeMux()
	InventoryControlModule{Service: inventorycontrol.NewService(&inventoryRepo{}, inventoryIDs{}), Bulk: inventorycontrol.NewBulkService(repo, inventoryIDs{})}.Register(mux, inventoryVerifier{})
	request := httptest.NewRequest("POST", "/v1/inventory/bulk/items", strings.NewReader(`{"code":"P","description":"Part","base_uom":"EA","tracking_mode":"lot","costing_method":"average"}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 400 || repo.items != 0 {
		t.Fatalf("unsupported method status=%d calls=%d", response.Code, repo.items)
	}
	request = httptest.NewRequest("GET", "/v1/inventory/bulk/availability?organization_id=forbidden&item_id=i", nil)
	request.Header.Set("Authorization", "Bearer valid")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 {
		t.Fatalf("organization scope status=%d body=%s", response.Code, response.Body.String())
	}
}

func TestBulkInventoryHTTPExposesStrictUnitOfMeasureConfigurationAndConversion(t *testing.T) {
	repo := &bulkHTTPRepo{}
	mux := http.NewServeMux()
	InventoryControlModule{Service: inventorycontrol.NewService(&inventoryRepo{}, inventoryIDs{}), Bulk: inventorycontrol.NewBulkService(repo, inventoryIDs{})}.Register(mux, inventoryVerifier{})
	request := httptest.NewRequest("POST", "/v1/inventory/bulk/item-units", strings.NewReader(`{"item_id":"part","code":"BOX","quantity_per_unit":"12","rounding_precision":"1"}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || !strings.Contains(response.Body.String(), `"code":"BOX"`) {
		t.Fatalf("configure status=%d body=%s", response.Code, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/inventory/bulk/uom-conversions", strings.NewReader(`{"item_id":"part","code":"BOX","quantity":"2.5"}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 200 || !strings.Contains(response.Body.String(), `"base_quantity":"30"`) {
		t.Fatalf("convert status=%d body=%s", response.Code, response.Body.String())
	}
}

func TestBulkInventoryHTTPExposesScopedHandlingUnitConversion(t *testing.T) {
	repo := &bulkHTTPRepo{}
	mux := http.NewServeMux()
	InventoryControlModule{Service: inventorycontrol.NewService(&inventoryRepo{}, inventoryIDs{}), Bulk: inventorycontrol.NewBulkService(repo, inventoryIDs{})}.Register(mux, inventoryVerifier{})
	request := httptest.NewRequest("POST", "/v1/inventory/bulk/handling-unit-conversions", strings.NewReader(`{"request_id":"request-1","organization_id":"a","bin_id":"pick","item_id":"part","operation":"breakbulk","from_uom":"BOX","to_uom":"EA","from_quantity":"1"}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || !strings.Contains(response.Body.String(), `"operation":"breakbulk"`) {
		t.Fatalf("conversion status=%d body=%s", response.Code, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/inventory/bulk/handling-unit-conversions", strings.NewReader(`{"request_id":"request-2","organization_id":"forbidden","bin_id":"pick","item_id":"part","operation":"breakbulk","from_uom":"BOX","to_uom":"EA","from_quantity":"1"}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 {
		t.Fatalf("forbidden conversion status=%d body=%s", response.Code, response.Body.String())
	}
}
````

### FILE: `db/migrations/0027_operational_warehouse_fefo.up.sql`

```yaml
block_id: "GO-OPS-API:operational-warehouse-fefo-migration-up:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps 2eae56d WhsePostReceipt, CreatePutaway, CreatePick, FEFO and WarehouseActivity invariants; portable PostgreSQL scope"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "edfb67dcb3fba6e6de48ce226e5b685b8d8d7498b590de99fe7861e95610adf9"
variables: []
secrets_allowed: false
```

````sql
begin;

alter table inventory.warehouse_bin
  drop constraint warehouse_bin_bin_type_check;

alter table inventory.warehouse_bin
  add column bin_rank integer not null default 0 check (bin_rank between 0 and 1000000),
  add constraint warehouse_bin_bin_type_check
    check (bin_type in ('receive','ship','put-away','pick','putpick','qc'));

create table inventory.warehouse_receipt (
  tenant_id uuid not null,
  receipt_id text not null,
  organization_id text not null,
  receive_bin_id text not null,
  source_kind text not null,
  source_id text not null,
  posting_date date not null,
  status text not null check (status = 'posted'),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, receipt_id),
  foreign key (tenant_id, organization_id, receive_bin_id)
    references inventory.warehouse_bin (tenant_id, organization_id, bin_id),
  unique (tenant_id, source_kind, source_id),
  check (length(receipt_id) between 1 and 128),
  check (length(source_kind) between 1 and 64),
  check (length(source_id) between 1 and 128)
);

create table inventory.warehouse_receipt_line (
  tenant_id uuid not null,
  receipt_id text not null,
  line_id text not null,
  item_id text not null,
  lot_id text,
  quantity numeric(20,6) not null check (quantity > 0),
  unit_cost numeric(20,4) not null check (unit_cost >= 0),
  inventory_entry_id text not null,
  primary key (tenant_id, receipt_id, line_id),
  foreign key (tenant_id, receipt_id)
    references inventory.warehouse_receipt (tenant_id, receipt_id),
  foreign key (tenant_id, item_id)
    references inventory.stock_item (tenant_id, item_id),
  foreign key (tenant_id, lot_id)
    references inventory.inventory_lot (tenant_id, lot_id),
  foreign key (tenant_id, inventory_entry_id)
    references inventory.bulk_inventory_entry (tenant_id, entry_id),
  unique (tenant_id, inventory_entry_id)
);

create table inventory.warehouse_activity (
  tenant_id uuid not null,
  activity_id text not null,
  organization_id text not null,
  activity_type text not null check (activity_type in ('put-away','pick')),
  status text not null check (status in ('open','registered','cancelled')),
  source_kind text not null,
  source_id text not null,
  source_line_id text not null default '',
  request_id text not null,
  assigned_to text not null default '',
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  registered_at timestamptz,
  primary key (tenant_id, activity_id),
  foreign key (tenant_id, organization_id)
    references org.organization (tenant_id, organization_id),
  unique (tenant_id, request_id),
  unique (tenant_id, activity_id, organization_id),
  check (length(activity_id) between 1 and 128),
  check (length(request_id) between 1 and 128),
  check ((status = 'registered' and registered_at is not null) or
         (status <> 'registered' and registered_at is null))
);

create table inventory.warehouse_activity_line (
  tenant_id uuid not null,
  activity_id text not null,
  organization_id text not null,
  line_id text not null,
  sequence_no integer not null check (sequence_no > 0),
  from_bin_id text not null,
  to_bin_id text not null,
  item_id text not null,
  lot_id text,
  reservation_id text,
  quantity numeric(20,6) not null check (quantity > 0),
  expiration_date date,
  source_entry_id text,
  primary key (tenant_id, activity_id, line_id),
  foreign key (tenant_id, activity_id)
    references inventory.warehouse_activity (tenant_id, activity_id),
  foreign key (tenant_id, activity_id, organization_id)
    references inventory.warehouse_activity (tenant_id, activity_id, organization_id),
  foreign key (tenant_id, organization_id, from_bin_id)
    references inventory.warehouse_bin (tenant_id, organization_id, bin_id),
  foreign key (tenant_id, organization_id, to_bin_id)
    references inventory.warehouse_bin (tenant_id, organization_id, bin_id),
  foreign key (tenant_id, item_id)
    references inventory.stock_item (tenant_id, item_id),
  foreign key (tenant_id, lot_id)
    references inventory.inventory_lot (tenant_id, lot_id),
  foreign key (tenant_id, reservation_id)
    references inventory.bulk_reservation (tenant_id, reservation_id),
  foreign key (tenant_id, source_entry_id)
    references inventory.bulk_inventory_entry (tenant_id, entry_id),
  unique (tenant_id, activity_id, sequence_no),
  check (from_bin_id <> to_bin_id)
);

create index warehouse_activity_open_scope_idx
  on inventory.warehouse_activity(tenant_id, organization_id, activity_type, created_at)
  where status = 'open';

create function inventory.reject_warehouse_evidence_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='posted warehouse evidence is append-only';
end;
$function$;

create trigger warehouse_receipt_immutable
before update or delete on inventory.warehouse_receipt
for each row execute function inventory.reject_warehouse_evidence_mutation();

create trigger warehouse_receipt_line_immutable
before update or delete on inventory.warehouse_receipt_line
for each row execute function inventory.reject_warehouse_evidence_mutation();

create trigger warehouse_activity_line_immutable
before update or delete on inventory.warehouse_activity_line
for each row execute function inventory.reject_warehouse_evidence_mutation();

create function inventory.enforce_warehouse_activity_transition()
returns trigger language plpgsql as $function$
begin
  if old.status <> 'open' or new.status not in ('registered','cancelled') or
     new.version <> old.version + 1 or
     new.activity_id <> old.activity_id or new.organization_id <> old.organization_id or
     new.activity_type <> old.activity_type or new.source_kind <> old.source_kind or
     new.source_id <> old.source_id or new.source_line_id <> old.source_line_id or
     new.request_id <> old.request_id or new.created_at <> old.created_at then
    raise exception using errcode='55000', message='invalid warehouse activity transition';
  end if;
  return new;
end;
$function$;

create trigger warehouse_activity_transition_guard
before update on inventory.warehouse_activity
for each row execute function inventory.enforce_warehouse_activity_transition();

commit;
````

### FILE: `db/migrations/0027_operational_warehouse_fefo.down.sql`

```yaml
block_id: "GO-OPS-API:operational-warehouse-fefo-migration-down:v1"
operation: CREATE
provenance: ADAPTED
source: "Portable rollback for the Microsoft-governed V161 warehouse boundary"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "bc08e6a74a4fbfd4330cd6985f9c77f4403c345bbe3112ffb0409560410dc7da"
variables: []
secrets_allowed: false
```

````sql
begin;
drop trigger if exists warehouse_activity_transition_guard on inventory.warehouse_activity;
drop function if exists inventory.enforce_warehouse_activity_transition();
drop trigger if exists warehouse_activity_line_immutable on inventory.warehouse_activity_line;
drop trigger if exists warehouse_receipt_line_immutable on inventory.warehouse_receipt_line;
drop trigger if exists warehouse_receipt_immutable on inventory.warehouse_receipt;
drop function if exists inventory.reject_warehouse_evidence_mutation();
drop table if exists inventory.warehouse_activity_line;
drop table if exists inventory.warehouse_activity;
drop table if exists inventory.warehouse_receipt_line;
drop table if exists inventory.warehouse_receipt;
alter table inventory.warehouse_bin drop constraint warehouse_bin_bin_type_check;
alter table inventory.warehouse_bin drop column bin_rank;
alter table inventory.warehouse_bin add constraint warehouse_bin_bin_type_check
  check (bin_type in ('receive','ship','put-away','pick','putpick'));
commit;
````

### FILE: `docs/inventory/MICROSOFT_BC_WAREHOUSE_DERIVATION.md`

```yaml
block_id: "GO-OPS-API:microsoft-bc-warehouse-derivation:v1"
operation: CREATE
provenance: AUTHORED
source: "local evidence record derived from pinned Microsoft BCApps MIT files and Microsoft Learn"
license: "LicenseRef-Workspace-Owner"
sha256: "f7760ce16310f7c272ea595f845c891c060625e76f5d52406f0cdf4e7cdf05e9"
variables: []
secrets_allowed: false
```

````markdown
# Microsoft BC operational warehouse derivation

This portable Go/PostgreSQL boundary is `ADAPTED`; it is not Microsoft-authored Go and is not a replacement for Business Central. Its behavioral authority is Microsoft BCApps at commit `2eae56d704a1fd035d104f333602aea7091b7749`, tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`, under the repository MIT license.

## Exact official source identities

| BCApps path under `src/Layers/W1/BaseApp` | Git blob | Bytes | SHA-256 |
|---|---|---:|---|
| `Warehouse/Document/WhsePostReceipt.Codeunit.al` | `dbd6327a935a69c564d112203eaf42ee845ba393` | 105803 | `9aae2fa85c20e12e70b88b3214fd4b4defcb13b6e767c225251d413eda626989` |
| `Warehouse/Activity/CreatePutaway.Codeunit.al` | `df3011a0a74ace99e3d1788996ba750a7732a1e6` | 77370 | `b004c0663d237585d8b57a06356431acb2302cdfbbd30f7b82154d282f325b2b` |
| `Warehouse/Activity/CreatePick.Codeunit.al` | `07be15a7effeeda1ac749ef033d20d817bb5b218` | 285538 | `6ca4acba3eebb46c719267ed788d8efc6ce097005db0d1a560489797e03b35e5` |
| `Warehouse/Tracking/WhseItemTrackingFEFO.Codeunit.al` | `9568f0f687384a65ac6dc4135273cc44b7256c4e` | 17061 | `24f0b95de17d6d2a042b97fe2cefb0ded952df8db2f0ebfefdb75c47aca09965` |
| `Warehouse/Activity/WarehouseActivityHeader.Table.al` | `4e3e2cf15ed00cfdc627898a1fafb54a8c8eca29` | 47747 | `34f6f3b633aed6c77c1cf286e9528d449fb255dcc219531946b43b59ce862030` |
| `Warehouse/Activity/WarehouseActivityLine.Table.al` | `83d6dc89fa2b088e05d9ac224fe0ab9e3c1921b9` | 181260 | `f4ab727d452c158b1f397e83708150f04b9fac22262e378bbca0bd9a8ed6836b` |

Microsoft Learn authorities checked on 2026-08-31:

- `https://learn.microsoft.com/en-us/dynamics365/business-central/design-details-warehouse-management`
- `https://learn.microsoft.com/en-us/dynamics365/business-central/design-details-inbound-warehouse-flow`
- `https://learn.microsoft.com/en-us/dynamics365/business-central/warehouse-how-to-put-items-away-with-warehouse-put-aways`
- `https://learn.microsoft.com/en-us/dynamics365/business-central/warehouse-how-to-pick-items-for-warehouse-shipment`

## Narrow invariants transferred

- Posting a warehouse receipt creates inventory in a RECEIVE bin before a separate put-away is registered.
- Put-away planning selects only eligible PUT AWAY or PUTPICK bins, accounts for maximum quantity and already-open placements, and prefers default/high-ranking bins.
- Registration records paired Take/Place movement evidence; an open activity can transition once to registered or cancelled under optimistic version control.
- Picks use only PICK or PUTPICK bins. RECEIVE, SHIP, PUT AWAY, QC, blocked and unapproved dedicated bins are excluded.
- FEFO excludes expired/blocked tracking identities and existing reservations. It orders by expiration, lot identity, registration age and then bin ranking; selection and reservation commit atomically.
- A registered pick moves both quantity and its reservation to the SHIP bin, keeping issuance separate. Cancellation is allowed only while the pick is open and releases every reservation atomically.
- A pick resolves its requested target UOM first. It consumes exact target-UOM composition before considering a larger source UOM, and that fallback is admitted only when `allow_breakbulk=true`, matching the explicit `AllowBreakbulk` branch in the fixed BCApps replenishment/pick behavior.
- Every packaging-aware line persists From/To UOM, quantity and factor. Open work reserves both physical base quantity and the exact source composition; registration consumes the source package, records immutable breakbulk Take/Place evidence when conversion is required, moves the requested target composition and retains any package remainder as exact base composition. Cancellation releases both reservations.
- Receipt evidence, activity lines and inventory entries are append-only. Request IDs and source identities prevent silent duplicate effects.

## Deliberate limits

The portable lane admits only the packaging-aware customer/transfer picking and bin-replenishment behavior demonstrated by migrations 0033–0036 and its PostgreSQL tests. Cubage/weight, warehouse classes, barcode/device flows, partial line handling, serial tracking, picks for production/assembly/projects, planning/CTP and the full Business Central extension-event surface remain rejected until their own official derivation and executable evidence exist. Cross-docking and cross-organization transfer are governed separately by their fixed derivations and tests.

## Regression mapping

Migration 0027 enforces bin types, ranking, immutable receipt/activity evidence and one-way state transitions. Migration 0036 adds immutable From/To UOM evidence and stateful physical/composition reservations without binding history to consumable composition rows. Domain tests reject malformed quantities and unsupported demands before persistence. PostgreSQL integration proves clean receipt-to-put-away, capacity/ranking, FEFO across two lots and bins, one winner under concurrent pick pressure, exact BOX→BOX picking, explicit rejection of implicit BOX→EA, authorized automatic breakbulk, registration conservation, cancellation release and SQLSTATE 55000 immutability. HTTP tests prove strict JSON, organization scope and executable registration routing.
````

### FILE: `internal/inventorycontrol/warehouse.go`

```yaml
block_id: "GO-OPS-API:operational-warehouse-domain:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps 2eae56d warehouse receipt/activity/pick/FEFO invariants adapted to the existing Go inventory owner"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "05c0ab474f08f4d04742dc001827498a557b1da2e112e1ddd3ed0715e4e0b2c9"
variables: []
secrets_allowed: false
```

````go
package inventorycontrol

import (
	"context"
	"fmt"
	"time"
)

type WarehouseReceiptCommand struct {
	RequestID        string     `json:"request_id"`
	OrganizationID   string     `json:"organization_id"`
	ReceiveBinID     string     `json:"receive_bin_id"`
	ItemID           string     `json:"item_id"`
	LotNo            string     `json:"lot_no,omitempty"`
	ExpirationDate   *time.Time `json:"expiration_date,omitempty"`
	WarrantyDate     *time.Time `json:"warranty_date,omitempty"`
	Quantity         string     `json:"quantity,omitempty"`
	HandlingUOM      string     `json:"handling_uom,omitempty"`
	HandlingQuantity string     `json:"handling_quantity,omitempty"`
	AllowBreakbulk   bool       `json:"allow_breakbulk"`
	UnitCost         string     `json:"unit_cost"`
	PostingDate      time.Time  `json:"posting_date"`
	SourceKind       string     `json:"source_kind"`
	SourceID         string     `json:"source_id"`
}

type WarehousePickCommand struct {
	RequestID      string `json:"request_id"`
	OrganizationID string `json:"organization_id"`
	ShipBinID      string `json:"ship_bin_id"`
	ItemID         string `json:"item_id"`
	DemandKind     string `json:"demand_kind"`
	DemandID       string `json:"demand_id"`
	DemandLineID   string `json:"demand_line_id"`
	Quantity       string `json:"quantity"`
	HandlingUOM    string `json:"handling_uom,omitempty"`
	AllowBreakbulk bool   `json:"allow_breakbulk"`
	UseFEFO        bool   `json:"use_fefo"`
	AllowDedicated bool   `json:"allow_dedicated"`
	AssignedTo     string `json:"assigned_to,omitempty"`
}

type SalesWarehouseBinding struct {
	RequestID    string `json:"request_id"`
	VariantID    string `json:"variant_id"`
	ItemID       string `json:"item_id"`
	SalesUOMCode string `json:"sales_uom_code"`
	Version      int64  `json:"version"`
}

type CustomerShipmentCommand struct {
	RequestID           string    `json:"request_id"`
	OrganizationID      string    `json:"organization_id"`
	OrderID             string    `json:"order_id"`
	WarehouseActivityID string    `json:"warehouse_activity_id"`
	PostingDate         time.Time `json:"posting_date"`
}

type CustomerShipmentLine struct {
	ID              string `json:"id"`
	OrderLineID     string `json:"order_line_id"`
	VariantID       string `json:"variant_id"`
	ItemID          string `json:"item_id"`
	SalesUOMCode    string `json:"sales_uom_code"`
	Quantity        string `json:"quantity"`
	QuantityBase    string `json:"quantity_base"`
	CostAmount      string `json:"cost_amount"`
	AllocationCount int    `json:"allocation_count"`
}

type CustomerShipment struct {
	ID                  string               `json:"id"`
	RequestID           string               `json:"request_id"`
	OrganizationID      string               `json:"organization_id"`
	OrderID             string               `json:"order_id"`
	WarehouseActivityID string               `json:"warehouse_activity_id"`
	PostingDate         time.Time            `json:"posting_date"`
	OrderVersion        int64                `json:"order_version"`
	FulfillmentState    string               `json:"fulfillment_state"`
	Line                CustomerShipmentLine `json:"line"`
}

type CustomerShipmentIDs struct {
	ShipmentID string
	LineID     string
	EventID    string
}

type WarehouseActivityLine struct {
	ID                 string     `json:"id"`
	Sequence           int        `json:"sequence"`
	FromBinID          string     `json:"from_bin_id"`
	ToBinID            string     `json:"to_bin_id"`
	ItemID             string     `json:"item_id"`
	LotID              string     `json:"lot_id,omitempty"`
	LotNo              string     `json:"lot_no,omitempty"`
	ReservationID      string     `json:"reservation_id,omitempty"`
	ReservationVersion int64      `json:"reservation_version,omitempty"`
	Quantity           string     `json:"quantity"`
	FromUOMCode        string     `json:"from_uom_code"`
	FromUOMQuantity    string     `json:"from_uom_quantity"`
	FromQuantityPerUOM string     `json:"from_quantity_per_uom"`
	UOMCode            string     `json:"uom_code"`
	UOMQuantity        string     `json:"uom_quantity"`
	QuantityPerUOM     string     `json:"quantity_per_uom"`
	ExpirationDate     *time.Time `json:"expiration_date,omitempty"`
	SourceEntryID      string     `json:"source_entry_id,omitempty"`
}

type WarehouseActivity struct {
	ID             string                  `json:"id"`
	OrganizationID string                  `json:"organization_id"`
	Type           string                  `json:"type"`
	Status         string                  `json:"status"`
	SourceKind     string                  `json:"source_kind"`
	SourceID       string                  `json:"source_id"`
	SourceLineID   string                  `json:"source_line_id,omitempty"`
	RequestID      string                  `json:"request_id"`
	AssignedTo     string                  `json:"assigned_to,omitempty"`
	Version        int64                   `json:"version"`
	Lines          []WarehouseActivityLine `json:"lines"`
}

type WarehouseReceiptResult struct {
	ReceiptID        string            `json:"receipt_id"`
	EntryID          string            `json:"entry_id"`
	LotID            string            `json:"lot_id,omitempty"`
	Quantity         string            `json:"quantity"`
	HandlingUOM      string            `json:"handling_uom"`
	HandlingQuantity string            `json:"handling_quantity"`
	CostAmount       string            `json:"cost_amount"`
	PutAway          WarehouseActivity `json:"put_away"`
}

type WarehouseIDs struct {
	ActivityID        string
	ReceiptID         string
	LineID            string
	EntryID           string
	LayerID           string
	LotID             string
	EventID           string
	ConversionID      string
	TakeLineID        string
	PlaceLineID       string
	ConversionEventID string
}

type WarehouseRepository interface {
	ConfigureSalesWarehouseBinding(context.Context, string, string, SalesWarehouseBinding) (SalesWarehouseBinding, error)
	PostCustomerShipment(context.Context, string, CustomerShipmentIDs, CustomerShipmentCommand) (CustomerShipment, error)
	PostWarehouseReceipt(context.Context, string, WarehouseIDs, WarehouseReceiptCommand) (WarehouseReceiptResult, error)
	CreateWarehousePick(context.Context, string, WarehouseIDs, WarehousePickCommand) (WarehouseActivity, error)
	CreateWarehouseReplenishment(context.Context, string, WarehouseIDs, WarehouseReplenishmentCommand) (WarehouseActivity, error)
	ConfigureWarehouseCrossDock(context.Context, string, string, WarehouseCrossDockPolicy) (WarehouseCrossDockPolicy, error)
	RegisterWarehouseActivity(context.Context, string, string, string, int64, time.Time, string) (WarehouseActivity, error)
	CancelWarehousePick(context.Context, string, string, string, int64, string) (WarehouseActivity, error)
	CancelWarehousePutAway(context.Context, string, string, string, int64, string) (WarehouseActivity, error)
	CancelWarehouseReplenishment(context.Context, string, string, string, int64, string) (WarehouseActivity, error)
}

type WarehouseService struct {
	repository WarehouseRepository
	ids        IDGenerator
}

func NewWarehouseService(repository WarehouseRepository, ids IDGenerator) *WarehouseService {
	return &WarehouseService{repository: repository, ids: ids}
}

func (s *WarehouseService) newIDs() WarehouseIDs {
	return WarehouseIDs{ActivityID: s.ids.New(), ReceiptID: s.ids.New(), LineID: s.ids.New(), EntryID: s.ids.New(), LayerID: s.ids.New(), LotID: s.ids.New(), EventID: s.ids.New(), ConversionID: s.ids.New(), TakeLineID: s.ids.New(), PlaceLineID: s.ids.New(), ConversionEventID: s.ids.New()}
}

func (s *WarehouseService) ConfigureSalesBinding(ctx context.Context, tenant string, value SalesWarehouseBinding) (SalesWarehouseBinding, error) {
	if tenant == "" || value.RequestID == "" || value.VariantID == "" || value.ItemID == "" || value.SalesUOMCode == "" {
		return SalesWarehouseBinding{}, fmt.Errorf("invalid sales warehouse binding")
	}
	value.Version = 1
	return s.repository.ConfigureSalesWarehouseBinding(ctx, tenant, s.ids.New(), value)
}

func (s *WarehouseService) PostCustomerShipment(ctx context.Context, tenant string, value CustomerShipmentCommand) (CustomerShipment, error) {
	if tenant == "" || value.RequestID == "" || value.OrganizationID == "" || value.OrderID == "" || value.WarehouseActivityID == "" || value.PostingDate.IsZero() {
		return CustomerShipment{}, fmt.Errorf("invalid customer shipment")
	}
	value.PostingDate = value.PostingDate.UTC()
	ids := CustomerShipmentIDs{ShipmentID: s.ids.New(), LineID: s.ids.New(), EventID: s.ids.New()}
	return s.repository.PostCustomerShipment(ctx, tenant, ids, value)
}

func (s *WarehouseService) PostReceipt(ctx context.Context, tenant string, value WarehouseReceiptCommand) (WarehouseReceiptResult, error) {
	legacyQuantity := positiveDecimal(value.Quantity, quantityPattern) && value.HandlingUOM == "" && value.HandlingQuantity == ""
	handlingQuantity := value.Quantity == "" && value.HandlingUOM != "" && positiveDecimal(value.HandlingQuantity, quantityPattern)
	if tenant == "" || value.RequestID == "" || value.OrganizationID == "" || value.ReceiveBinID == "" || value.ItemID == "" || (!legacyQuantity && !handlingQuantity) || !positiveDecimal(value.UnitCost, costPattern) || value.PostingDate.IsZero() || value.SourceKind == "" || value.SourceID == "" {
		return WarehouseReceiptResult{}, fmt.Errorf("invalid warehouse receipt")
	}
	value.PostingDate = value.PostingDate.UTC()
	return s.repository.PostWarehouseReceipt(ctx, tenant, s.newIDs(), value)
}

func (s *WarehouseService) CancelPutAway(ctx context.Context, tenant, organization, activity string, version int64) (WarehouseActivity, error) {
	if tenant == "" || organization == "" || activity == "" || version < 1 {
		return WarehouseActivity{}, fmt.Errorf("invalid warehouse put-away cancellation")
	}
	return s.repository.CancelWarehousePutAway(ctx, tenant, organization, activity, version, s.ids.New())
}

func (s *WarehouseService) CreatePick(ctx context.Context, tenant string, value WarehousePickCommand) (WarehouseActivity, error) {
	allowed := map[string]bool{"customer-order": true, "service": true, "transfer-outbound": true, "manual": true}
	if tenant == "" || value.RequestID == "" || value.OrganizationID == "" || value.ShipBinID == "" || value.ItemID == "" || !allowed[value.DemandKind] || value.DemandID == "" || value.DemandLineID == "" || !positiveDecimal(value.Quantity, quantityPattern) {
		return WarehouseActivity{}, fmt.Errorf("invalid warehouse pick")
	}
	return s.repository.CreateWarehousePick(ctx, tenant, s.newIDs(), value)
}

func (s *WarehouseService) Register(ctx context.Context, tenant, organization, activity string, version int64, postingDate time.Time) (WarehouseActivity, error) {
	if tenant == "" || organization == "" || activity == "" || version < 1 || postingDate.IsZero() {
		return WarehouseActivity{}, fmt.Errorf("invalid warehouse registration")
	}
	return s.repository.RegisterWarehouseActivity(ctx, tenant, organization, activity, version, postingDate.UTC(), s.ids.New())
}

func (s *WarehouseService) CancelPick(ctx context.Context, tenant, organization, activity string, version int64) (WarehouseActivity, error) {
	if tenant == "" || organization == "" || activity == "" || version < 1 {
		return WarehouseActivity{}, fmt.Errorf("invalid warehouse pick cancellation")
	}
	return s.repository.CancelWarehousePick(ctx, tenant, organization, activity, version, s.ids.New())
}
````

### FILE: `internal/inventorycontrol/warehouse_test.go`

```yaml
block_id: "GO-OPS-API:operational-warehouse-domain-test:v1"
operation: CREATE
provenance: ADAPTED
source: "Regression of the Microsoft-governed portable warehouse command boundary"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "b9f3429c9453ea50271fe611610423d5ff74e81cacee6b5fd33d02d8acc1f58a"
variables: []
secrets_allowed: false
```

````go
package inventorycontrol

import (
	"context"
	"testing"
	"time"
)

type warehouseIDs struct{ n int }

func (g *warehouseIDs) New() string {
	g.n++
	return "id-" + time.Unix(int64(g.n), 0).UTC().Format("150405")
}

type warehouseRepo struct {
	bindings             int
	shipments            int
	receipts             int
	picks                int
	replenishments       int
	cancellations        int
	putAwayCancellations int
	crossDock            int
}

func (r *warehouseRepo) ConfigureSalesWarehouseBinding(_ context.Context, _ string, _ string, value SalesWarehouseBinding) (SalesWarehouseBinding, error) {
	r.bindings++
	return value, nil
}
func (r *warehouseRepo) PostCustomerShipment(_ context.Context, _ string, ids CustomerShipmentIDs, value CustomerShipmentCommand) (CustomerShipment, error) {
	r.shipments++
	return CustomerShipment{ID: ids.ShipmentID, RequestID: value.RequestID, OrganizationID: value.OrganizationID, OrderID: value.OrderID, WarehouseActivityID: value.WarehouseActivityID, PostingDate: value.PostingDate, Line: CustomerShipmentLine{ID: ids.LineID}}, nil
}

func (r *warehouseRepo) PostWarehouseReceipt(_ context.Context, _ string, ids WarehouseIDs, value WarehouseReceiptCommand) (WarehouseReceiptResult, error) {
	r.receipts++
	return WarehouseReceiptResult{ReceiptID: ids.ReceiptID, Quantity: value.Quantity}, nil
}
func (r *warehouseRepo) CreateWarehousePick(_ context.Context, _ string, ids WarehouseIDs, value WarehousePickCommand) (WarehouseActivity, error) {
	r.picks++
	return WarehouseActivity{ID: ids.ActivityID, OrganizationID: value.OrganizationID, Type: "pick", Version: 1}, nil
}
func (r *warehouseRepo) CreateWarehouseReplenishment(_ context.Context, _ string, ids WarehouseIDs, value WarehouseReplenishmentCommand) (WarehouseActivity, error) {
	r.replenishments++
	return WarehouseActivity{ID: ids.ActivityID, OrganizationID: value.OrganizationID, Type: "movement", Version: 1}, nil
}
func (r *warehouseRepo) ConfigureWarehouseCrossDock(_ context.Context, _ string, _ string, value WarehouseCrossDockPolicy) (WarehouseCrossDockPolicy, error) {
	r.crossDock++
	return value, nil
}
func (*warehouseRepo) RegisterWarehouseActivity(context.Context, string, string, string, int64, time.Time, string) (WarehouseActivity, error) {
	return WarehouseActivity{}, nil
}
func (*warehouseRepo) CancelWarehousePick(context.Context, string, string, string, int64, string) (WarehouseActivity, error) {
	return WarehouseActivity{}, nil
}
func (r *warehouseRepo) CancelWarehousePutAway(context.Context, string, string, string, int64, string) (WarehouseActivity, error) {
	r.putAwayCancellations++
	return WarehouseActivity{}, nil
}
func (r *warehouseRepo) CancelWarehouseReplenishment(context.Context, string, string, string, int64, string) (WarehouseActivity, error) {
	r.cancellations++
	return WarehouseActivity{}, nil
}

func TestWarehouseServiceRejectsAmbiguityBeforeRepository(t *testing.T) {
	repo := &warehouseRepo{}
	service := NewWarehouseService(repo, &warehouseIDs{})
	if _, err := service.ConfigureSalesBinding(context.Background(), "tenant", SalesWarehouseBinding{RequestID: "binding-request", VariantID: "variant", ItemID: "part"}); err == nil || repo.bindings != 0 {
		t.Fatalf("invalid sales binding reached repository: err=%v calls=%d", err, repo.bindings)
	}
	binding, err := service.ConfigureSalesBinding(context.Background(), "tenant", SalesWarehouseBinding{RequestID: "binding-request", VariantID: "variant", ItemID: "part", SalesUOMCode: "BOX"})
	if err != nil || repo.bindings != 1 || binding.Version != 1 {
		t.Fatalf("valid sales binding failed: value=%+v err=%v calls=%d", binding, err, repo.bindings)
	}
	shipment := CustomerShipmentCommand{RequestID: "shipment-request", OrganizationID: "warehouse", OrderID: "order", WarehouseActivityID: "pick", PostingDate: time.Now()}
	invalidShipment := shipment
	invalidShipment.WarehouseActivityID = ""
	if _, err := service.PostCustomerShipment(context.Background(), "tenant", invalidShipment); err == nil || repo.shipments != 0 {
		t.Fatalf("invalid customer shipment reached repository: err=%v calls=%d", err, repo.shipments)
	}
	posted, err := service.PostCustomerShipment(context.Background(), "tenant", shipment)
	if err != nil || repo.shipments != 1 || posted.ID == "" || posted.Line.ID == "" || !posted.PostingDate.Equal(shipment.PostingDate.UTC()) {
		t.Fatalf("valid customer shipment failed: value=%+v err=%v calls=%d", posted, err, repo.shipments)
	}
	validReceipt := WarehouseReceiptCommand{RequestID: "receipt-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "part", Quantity: "1", UnitCost: "10", PostingDate: time.Now(), SourceKind: "purchase", SourceID: "po-1"}
	invalidReceipt := validReceipt
	invalidReceipt.Quantity = "1e3"
	if _, err := service.PostReceipt(context.Background(), "tenant", invalidReceipt); err == nil || repo.receipts != 0 {
		t.Fatalf("invalid receipt reached repository: err=%v calls=%d", err, repo.receipts)
	}
	if _, err := service.PostReceipt(context.Background(), "tenant", validReceipt); err != nil || repo.receipts != 1 {
		t.Fatalf("valid receipt failed: err=%v calls=%d", err, repo.receipts)
	}
	packagedReceipt := validReceipt
	packagedReceipt.Quantity = ""
	packagedReceipt.HandlingUOM = "BOX"
	packagedReceipt.HandlingQuantity = "2"
	if _, err := service.PostReceipt(context.Background(), "tenant", packagedReceipt); err != nil || repo.receipts != 2 {
		t.Fatalf("valid packaged receipt failed: err=%v calls=%d", err, repo.receipts)
	}
	ambiguousReceipt := packagedReceipt
	ambiguousReceipt.Quantity = "24"
	if _, err := service.PostReceipt(context.Background(), "tenant", ambiguousReceipt); err == nil || repo.receipts != 2 {
		t.Fatalf("ambiguous receipt reached repository: err=%v calls=%d", err, repo.receipts)
	}
	if _, err := service.CancelPutAway(context.Background(), "tenant", "warehouse", "put-away", 0); err == nil || repo.putAwayCancellations != 0 {
		t.Fatalf("invalid put-away cancellation reached repository: err=%v calls=%d", err, repo.putAwayCancellations)
	}
	if _, err := service.CancelPutAway(context.Background(), "tenant", "warehouse", "put-away", 1); err != nil || repo.putAwayCancellations != 1 {
		t.Fatalf("valid put-away cancellation failed: err=%v calls=%d", err, repo.putAwayCancellations)
	}
	invalidPick := WarehousePickCommand{RequestID: "pick-request", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "invented", DemandID: "order", DemandLineID: "line", Quantity: "1", UseFEFO: true}
	if _, err := service.CreatePick(context.Background(), "tenant", invalidPick); err == nil || repo.picks != 0 {
		t.Fatalf("invalid pick reached repository: err=%v calls=%d", err, repo.picks)
	}
	invalidPick.DemandKind = "customer-order"
	if _, err := service.CreatePick(context.Background(), "tenant", invalidPick); err != nil || repo.picks != 1 {
		t.Fatalf("valid pick failed: err=%v calls=%d", err, repo.picks)
	}
}

func TestWarehouseReplenishmentRejectsAmbiguityBeforeRepository(t *testing.T) {
	repo := &warehouseRepo{}
	service := NewWarehouseService(repo, &warehouseIDs{})
	value := WarehouseReplenishmentCommand{RequestID: "replenish-request", OrganizationID: "warehouse", ToBinID: "pick", ItemID: "part", UseFEFO: true}
	invalid := value
	invalid.ToBinID = ""
	if _, err := service.CreateReplenishment(context.Background(), "tenant", invalid); err == nil || repo.replenishments != 0 {
		t.Fatalf("invalid replenishment reached repository: err=%v calls=%d", err, repo.replenishments)
	}
	if _, err := service.CreateReplenishment(context.Background(), "tenant", value); err != nil || repo.replenishments != 1 {
		t.Fatalf("valid replenishment failed: err=%v calls=%d", err, repo.replenishments)
	}
	if _, err := service.CancelReplenishment(context.Background(), "tenant", "warehouse", "activity", 0); err == nil || repo.cancellations != 0 {
		t.Fatalf("invalid cancellation reached repository: err=%v calls=%d", err, repo.cancellations)
	}
	if _, err := service.CancelReplenishment(context.Background(), "tenant", "warehouse", "activity", 1); err != nil || repo.cancellations != 1 {
		t.Fatalf("valid cancellation failed: err=%v calls=%d", err, repo.cancellations)
	}
}

func TestWarehouseCrossDockRejectsAmbiguityBeforeRepository(t *testing.T) {
	repo := &warehouseRepo{}
	service := NewWarehouseService(repo, &warehouseIDs{})
	value := WarehouseCrossDockPolicy{OrganizationID: "warehouse", ItemID: "part", BinID: "cross-dock", DueDateDays: 3, Enabled: true, Version: 1}
	invalid := value
	invalid.DueDateDays = 366
	if _, err := service.ConfigureCrossDock(context.Background(), "tenant", invalid); err == nil || repo.crossDock != 0 {
		t.Fatalf("invalid cross-dock policy reached repository: err=%v calls=%d", err, repo.crossDock)
	}
	if _, err := service.ConfigureCrossDock(context.Background(), "tenant", value); err != nil || repo.crossDock != 1 {
		t.Fatalf("valid cross-dock policy failed: err=%v calls=%d", err, repo.crossDock)
	}
}
````

### FILE: `internal/platform/postgres/warehouse.go`

```yaml
block_id: "GO-OPS-API:operational-warehouse-postgres:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps 2eae56d receipt, put-away, activity, availability, bin ranking and FEFO invariants adapted to PostgreSQL 18.6"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "6980dc86f5cbbf5d9d2057a0bb1ea9bf766399dee8defddbdc9cd4790bf1c0f6"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
)

type warehousePlacement struct {
	binID                 string
	quantity              string
	uomCode               string
	uomQuantity           string
	quantityPerUOM        string
	crossDockAllocationID string
	transferID            string
	transferLineID        string
	dueDate               time.Time
	lotID                 string
}

func warehouseLock(ctx context.Context, tx pgx.Tx, tenant, organization, item string) error {
	_, err := tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1 || ':' || $2 || ':' || $3,0))`, tenant, organization, item)
	return err
}

func decimalCapacity(maximum, onHand, planned string) (*big.Rat, bool) {
	if maximum == "" {
		return nil, true
	}
	max, okMax := parseRat(maximum)
	stock, okStock := parseRat(onHand)
	pending, okPending := parseRat(planned)
	if !okMax || !okStock || !okPending {
		return nil, false
	}
	available := new(big.Rat).Sub(max, stock)
	available.Sub(available, pending)
	if available.Sign() < 0 {
		available.SetInt64(0)
	}
	return available, true
}

func exactUOMQuantity(baseQuantity, quantityPerUOM, roundingPrecision string) (string, bool) {
	base, baseOK := parseRat(baseQuantity)
	factor, factorOK := parseRat(quantityPerUOM)
	precision, precisionOK := parseRat(roundingPrecision)
	if !baseOK || !factorOK || !precisionOK || factor.Sign() <= 0 || precision.Sign() <= 0 {
		return "", false
	}
	quantity := new(big.Rat).Quo(base, factor)
	steps := new(big.Rat).Quo(quantity, precision)
	if !steps.IsInt() {
		return "", false
	}
	return formatRat(quantity, 6), true
}

func (r *InventoryControl) PostWarehouseReceipt(ctx context.Context, tenant string, ids inventorycontrol.WarehouseIDs, value inventorycontrol.WarehouseReceiptCommand) (inventorycontrol.WarehouseReceiptResult, error) {
	result := inventorycontrol.WarehouseReceiptResult{ReceiptID: ids.ReceiptID, EntryID: ids.EntryID}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return result, err
	}
	defer tx.Rollback(ctx)
	if err = warehouseLock(ctx, tx, tenant, value.OrganizationID, value.ItemID); err != nil {
		return result, err
	}
	var tracking, receiveType string
	err = tx.QueryRow(ctx, `select i.tracking_mode,w.bin_type from inventory.stock_item i join inventory.item_bin_policy p on p.tenant_id=i.tenant_id and p.item_id=i.item_id join inventory.warehouse_bin w on w.tenant_id=p.tenant_id and w.organization_id=p.organization_id and w.bin_id=p.bin_id where i.tenant_id=$1 and i.item_id=$2 and p.organization_id=$3 and p.bin_id=$4 and not w.movement_blocked for share of i,p,w`, tenant, value.ItemID, value.OrganizationID, value.ReceiveBinID).Scan(&tracking, &receiveType)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, fmt.Errorf("warehouse receipt receive-bin contract: %w", inventorycontrol.ErrConflict)
	}
	if err != nil {
		return result, err
	}
	if receiveType != "receive" {
		return result, fmt.Errorf("warehouse receipt receive-bin contract: %w", inventorycontrol.ErrConflict)
	}
	handlingUOM, handlingQuantity := value.HandlingUOM, value.HandlingQuantity
	if handlingUOM == "" {
		handlingQuantity = value.Quantity
	}
	var baseUOM, baseQuantityPerUOM, basePrecision, quantityPerUOM, handlingPrecision, baseQuantity string
	err = tx.QueryRow(ctx, `select u.uom_code,u.qty_per_uom::text,u.rounding_precision::text,b.uom_code,b.qty_per_uom::text,b.rounding_precision::text,($4::numeric*u.qty_per_uom)::numeric(20,6)::text
from inventory.item_unit_of_measure u
join inventory.item_unit_of_measure b on b.tenant_id=u.tenant_id and b.item_id=u.item_id and b.is_base
where u.tenant_id=$1 and u.item_id=$2 and u.uom_code=coalesce(nullif($3,''),b.uom_code)
  and mod($4::numeric,u.rounding_precision)=0
  and mod($4::numeric*u.qty_per_uom,b.rounding_precision)=0
for share of u,b`, tenant, value.ItemID, handlingUOM, handlingQuantity).Scan(&handlingUOM, &quantityPerUOM, &handlingPrecision, &baseUOM, &baseQuantityPerUOM, &basePrecision, &baseQuantity)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, fmt.Errorf("warehouse receipt handling UOM contract: %w", inventorycontrol.ErrConflict)
	}
	if err != nil {
		return result, bulkConflict(err)
	}
	result.Quantity, result.HandlingUOM, result.HandlingQuantity = baseQuantity, handlingUOM, handlingQuantity
	lotID := ""
	if tracking == "lot" {
		if value.LotNo == "" {
			return result, fmt.Errorf("warehouse receipt lot required: %w", inventorycontrol.ErrConflict)
		}
		var existingExpiration, existingWarranty string
		err = tx.QueryRow(ctx, `select lot_id,coalesce(expiration_date::text,''),coalesce(warranty_date::text,'') from inventory.inventory_lot where tenant_id=$1 and item_id=$2 and lot_no=$3 for update`, tenant, value.ItemID, value.LotNo).Scan(&lotID, &existingExpiration, &existingWarranty)
		expiration, warranty := "", ""
		if value.ExpirationDate != nil {
			expiration = value.ExpirationDate.UTC().Format("2006-01-02")
		}
		if value.WarrantyDate != nil {
			warranty = value.WarrantyDate.UTC().Format("2006-01-02")
		}
		if errors.Is(err, pgx.ErrNoRows) {
			lotID = ids.LotID
			_, err = tx.Exec(ctx, `insert into inventory.inventory_lot(tenant_id,lot_id,item_id,lot_no,expiration_date,warranty_date) values($1,$2,$3,$4,nullif($5,'')::date,nullif($6,'')::date)`, tenant, lotID, value.ItemID, value.LotNo, expiration, warranty)
		} else if err == nil && (existingExpiration != expiration || existingWarranty != warranty) {
			return result, fmt.Errorf("warehouse receipt lot metadata mismatch: %w", inventorycontrol.ErrConflict)
		}
		if err != nil {
			return result, fmt.Errorf("warehouse receipt lot: %w", bulkConflict(err))
		}
	} else if value.LotNo != "" || value.ExpirationDate != nil || value.WarrantyDate != nil {
		return result, fmt.Errorf("warehouse receipt unexpected tracking: %w", inventorycontrol.ErrConflict)
	}
	result.LotID = lotID
	var balanceID, balanceVersion int64
	err = tx.QueryRow(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity,version) values($1,$2,$3,$4,nullif($5,''),$6::numeric,0,1) on conflict (tenant_id,organization_id,bin_id,item_id,lot_id) do update set quantity=inventory.bulk_balance.quantity+excluded.quantity,version=inventory.bulk_balance.version+1,updated_at=clock_timestamp() where (select max_quantity is null or inventory.bulk_balance.quantity+excluded.quantity<=max_quantity from inventory.item_bin_policy where tenant_id=$1 and organization_id=$2 and item_id=$4 and bin_id=$3) returning balance_id,version`, tenant, value.OrganizationID, value.ReceiveBinID, value.ItemID, lotID, baseQuantity).Scan(&balanceID, &balanceVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, fmt.Errorf("warehouse receipt receive capacity: %w", inventorycontrol.ErrConflict)
	}
	if err != nil {
		return result, fmt.Errorf("warehouse receipt balance: %w", bulkConflict(err))
	}
	err = tx.QueryRow(ctx, `insert into inventory.bulk_inventory_entry(tenant_id,entry_id,organization_id,bin_id,item_id,lot_id,entry_type,quantity,unit_cost,cost_amount,posting_date,source_kind,source_id) values($1,$2,$3,$4,$5,nullif($6,''),'receipt',$7::numeric,$8::numeric,round($7::numeric*$8::numeric,4),$9::date,$10,$11) returning cost_amount::text`, tenant, ids.EntryID, value.OrganizationID, value.ReceiveBinID, value.ItemID, lotID, baseQuantity, value.UnitCost, value.PostingDate, value.SourceKind, value.SourceID).Scan(&result.CostAmount)
	if err != nil {
		return result, fmt.Errorf("warehouse receipt entry: %w", bulkConflict(err))
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_cost_layer(tenant_id,layer_id,receipt_entry_id,organization_id,item_id,lot_id,posting_date,original_quantity,remaining_quantity,unit_cost) values($1,$2,$3,$4,$5,nullif($6,''),$7::date,$8::numeric,$8::numeric,$9::numeric)`, tenant, ids.LayerID, ids.EntryID, value.OrganizationID, value.ItemID, lotID, value.PostingDate, baseQuantity, value.UnitCost)
	if err != nil {
		return result, fmt.Errorf("warehouse receipt cost layer: %w", bulkConflict(err))
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_receipt(tenant_id,receipt_id,organization_id,receive_bin_id,source_kind,source_id,posting_date,status,allow_breakbulk) values($1,$2,$3,$4,$5,$6,$7::date,'posted',$8)`, tenant, ids.ReceiptID, value.OrganizationID, value.ReceiveBinID, value.SourceKind, value.SourceID, value.PostingDate, value.AllowBreakbulk)
	if err != nil {
		return result, fmt.Errorf("warehouse receipt header: %w", bulkConflict(err))
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_receipt_line(tenant_id,receipt_id,line_id,item_id,lot_id,quantity,unit_cost,inventory_entry_id,uom_code,uom_quantity,qty_per_uom) values($1,$2,$3,$4,nullif($5,''),$6::numeric,$7::numeric,$8,$9,$10::numeric,$11::numeric)`, tenant, ids.ReceiptID, ids.LineID, value.ItemID, lotID, baseQuantity, value.UnitCost, ids.EntryID, handlingUOM, handlingQuantity, quantityPerUOM)
	if err != nil {
		return result, fmt.Errorf("warehouse receipt line: %w", bulkConflict(err))
	}
	planningValue := value
	planningValue.Quantity = baseQuantity
	crossDockPlacements, remaining, err := planWarehouseCrossDock(ctx, tx, tenant, ids, planningValue, lotID)
	if err != nil {
		return result, fmt.Errorf("warehouse receipt cross-dock plan: %w", err)
	}
	rows, err := tx.Query(ctx, `select p.bin_id,coalesce(p.max_quantity::text,''),coalesce(b.quantity::text,'0'),coalesce((select sum(al.quantity)::text from inventory.warehouse_activity_line al join inventory.warehouse_activity a using(tenant_id,activity_id) where al.tenant_id=p.tenant_id and a.organization_id=p.organization_id and a.activity_type='put-away' and a.status='open' and al.to_bin_id=p.bin_id and al.item_id=p.item_id),'0') from inventory.item_bin_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) left join inventory.bulk_balance b on b.tenant_id=p.tenant_id and b.organization_id=p.organization_id and b.bin_id=p.bin_id and b.item_id=p.item_id and b.lot_id is not distinct from nullif($4,'') where p.tenant_id=$1 and p.organization_id=$2 and p.item_id=$3 and p.bin_id<>$5 and w.bin_type in ('put-away','putpick') and not w.cross_dock and not w.movement_blocked order by p.is_default desc,w.bin_rank desc,p.fixed desc,w.bin_code for share of p,w`, tenant, value.OrganizationID, value.ItemID, lotID, value.ReceiveBinID)
	if err != nil {
		return result, err
	}
	placements := crossDockPlacements
	for rows.Next() && remaining.Sign() > 0 {
		var binID, maximum, onHand, planned string
		if err = rows.Scan(&binID, &maximum, &onHand, &planned); err != nil {
			rows.Close()
			return result, err
		}
		capacity, valid := decimalCapacity(maximum, onHand, planned)
		if !valid {
			rows.Close()
			return result, fmt.Errorf("warehouse receipt target capacity parse: %w", inventorycontrol.ErrConflict)
		}
		take := new(big.Rat).Set(remaining)
		if capacity != nil {
			take = minRat(take, capacity)
		}
		if take.Sign() > 0 {
			placements = append(placements, warehousePlacement{binID: binID, quantity: formatRat(take, 6)})
			remaining.Sub(remaining, take)
		}
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return result, err
	}
	if remaining.Sign() != 0 || len(placements) == 0 {
		return result, fmt.Errorf("warehouse receipt put-away capacity: %w", inventorycontrol.ErrConflict)
	}
	requiresBreakbulk := false
	for index := range placements {
		quantity, exact := exactUOMQuantity(placements[index].quantity, quantityPerUOM, handlingPrecision)
		if !exact {
			requiresBreakbulk = handlingUOM != baseUOM
			break
		}
		placements[index].uomCode = handlingUOM
		placements[index].uomQuantity = quantity
		placements[index].quantityPerUOM = quantityPerUOM
	}
	if requiresBreakbulk && !value.AllowBreakbulk {
		return result, fmt.Errorf("warehouse receipt requires explicit breakbulk authorization: %w", inventorycontrol.ErrConflict)
	}
	if requiresBreakbulk {
		for index := range placements {
			quantity, exact := exactUOMQuantity(placements[index].quantity, baseQuantityPerUOM, basePrecision)
			if !exact {
				return result, fmt.Errorf("warehouse receipt base UOM placement: %w", inventorycontrol.ErrConflict)
			}
			placements[index].uomCode = baseUOM
			placements[index].uomQuantity = quantity
			placements[index].quantityPerUOM = baseQuantityPerUOM
		}
	}
	if handlingUOM != baseUOM && !requiresBreakbulk {
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_uom_balance set quantity=quantity-($3::numeric/qty_per_uom),quantity_base=quantity_base-$3::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and quantity-reserved_quantity >= ($3::numeric/qty_per_uom) and quantity_base >= $3::numeric`, balanceID, baseUOM, baseQuantity)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return result, bulkConflict(updateErr)
			}
			return result, inventorycontrol.ErrConflict
		}
		_, err = tx.Exec(ctx, `delete from inventory.bulk_uom_balance where balance_id=$1 and uom_code=$2 and quantity_base=0 and reserved_quantity=0`, balanceID, baseUOM)
		if err != nil {
			return result, err
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base) values($1,$2,$3::numeric,$4::numeric,$5::numeric) on conflict(balance_id,uom_code) do update set quantity=inventory.bulk_uom_balance.quantity+excluded.quantity,quantity_base=inventory.bulk_uom_balance.quantity_base+excluded.quantity_base,version=inventory.bulk_uom_balance.version+1,updated_at=clock_timestamp()`, balanceID, handlingUOM, handlingQuantity, quantityPerUOM, baseQuantity)
		if err != nil {
			return result, bulkConflict(err)
		}
	}
	if requiresBreakbulk {
		baseHandlingQuantity, exact := exactUOMQuantity(baseQuantity, baseQuantityPerUOM, basePrecision)
		if !exact {
			return result, inventorycontrol.ErrConflict
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_conversion(tenant_id,conversion_id,request_id,organization_id,bin_id,item_id,lot_id,operation,from_uom,to_uom,from_quantity,to_quantity,base_quantity,from_qty_per_uom,to_qty_per_uom,source_kind,source_id) values($1,$2,$3,$4,$5,$6,nullif($7,''),'breakbulk',$8,$9,$10::numeric,$11::numeric,$12::numeric,$13::numeric,$14::numeric,'warehouse-receipt',$15)`, tenant, ids.ConversionID, value.RequestID+"/automatic-breakbulk", value.OrganizationID, value.ReceiveBinID, value.ItemID, lotID, handlingUOM, baseUOM, handlingQuantity, baseHandlingQuantity, baseQuantity, quantityPerUOM, baseQuantityPerUOM, ids.ReceiptID)
		if err != nil {
			return result, fmt.Errorf("warehouse receipt automatic breakbulk: %w", bulkConflict(err))
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_conversion_line(tenant_id,conversion_id,line_id,sequence_no,action_type,uom_code,quantity,quantity_base) values($1,$2,$3,1,'take',$5,$6::numeric,$7::numeric),($1,$2,$4,2,'place',$8,$9::numeric,$7::numeric)`, tenant, ids.ConversionID, ids.TakeLineID, ids.PlaceLineID, handlingUOM, handlingQuantity, baseQuantity, baseUOM, baseHandlingQuantity)
		if err != nil {
			return result, fmt.Errorf("warehouse receipt automatic breakbulk lines: %w", bulkConflict(err))
		}
		if err = recordBulkEvent(ctx, tx, tenant, ids.ConversionEventID, "bulk-uom-conversion", ids.ConversionID, "bulk-uom-conversion.completed", 1, map[string]string{"request_id": value.RequestID + "/automatic-breakbulk", "organization_id": value.OrganizationID, "bin_id": value.ReceiveBinID, "item_id": value.ItemID, "lot_id": lotID, "operation": "breakbulk", "from_uom": handlingUOM, "to_uom": baseUOM, "from_quantity": handlingQuantity, "to_quantity": baseHandlingQuantity, "base_quantity": baseQuantity, "source_kind": "warehouse-receipt", "source_id": ids.ReceiptID}); err != nil {
			return result, err
		}
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity(tenant_id,activity_id,organization_id,activity_type,status,source_kind,source_id,source_line_id,request_id,version) values($1,$2,$3,'put-away','open','warehouse-receipt',$4,$5,$6,1)`, tenant, ids.ActivityID, value.OrganizationID, ids.ReceiptID, ids.LineID, value.RequestID)
	if err != nil {
		return result, fmt.Errorf("warehouse put-away header: %w", bulkConflict(err))
	}
	activity := inventorycontrol.WarehouseActivity{ID: ids.ActivityID, OrganizationID: value.OrganizationID, Type: "put-away", Status: "open", SourceKind: "warehouse-receipt", SourceID: ids.ReceiptID, SourceLineID: ids.LineID, RequestID: value.RequestID, Version: 1, Lines: []inventorycontrol.WarehouseActivityLine{}}
	for index, placement := range placements {
		lineID := fmt.Sprintf("%s-%06d", ids.ActivityID, index+1)
		_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity_line(tenant_id,activity_id,organization_id,line_id,sequence_no,from_bin_id,to_bin_id,item_id,lot_id,quantity,expiration_date,source_entry_id,uom_code,uom_quantity,qty_per_uom,from_uom_code,from_uom_quantity,from_qty_per_uom) values($1,$2,$3,$4,$5,$6,$7,$8,nullif($9,''),$10::numeric,$11,$12,$13,$14::numeric,$15::numeric,$13,$14::numeric,$15::numeric)`, tenant, ids.ActivityID, value.OrganizationID, lineID, index+1, value.ReceiveBinID, placement.binID, value.ItemID, lotID, placement.quantity, value.ExpirationDate, ids.EntryID, placement.uomCode, placement.uomQuantity, placement.quantityPerUOM)
		if err != nil {
			return result, fmt.Errorf("warehouse put-away line: %w", bulkConflict(err))
		}
		if placement.crossDockAllocationID != "" {
			_, err = tx.Exec(ctx, `insert into inventory.warehouse_crossdock_allocation(tenant_id,allocation_id,receipt_id,receipt_line_id,put_away_activity_id,put_away_line_id,organization_id,bin_id,item_id,lot_id,transfer_id,transfer_line_id,quantity,due_date,available) values($1,$2,$3,$4,$5,$6,$7,$8,$9,nullif($10,''),$11,$12,$13::numeric,$14::date,false)`, tenant, placement.crossDockAllocationID, ids.ReceiptID, ids.LineID, ids.ActivityID, lineID, value.OrganizationID, placement.binID, value.ItemID, lotID, placement.transferID, placement.transferLineID, placement.quantity, placement.dueDate)
			if err != nil {
				return result, fmt.Errorf("warehouse cross-dock allocation: %w", bulkConflict(err))
			}
		}
		activity.Lines = append(activity.Lines, inventorycontrol.WarehouseActivityLine{ID: lineID, Sequence: index + 1, FromBinID: value.ReceiveBinID, ToBinID: placement.binID, ItemID: value.ItemID, LotID: lotID, LotNo: value.LotNo, Quantity: placement.quantity, FromUOMCode: placement.uomCode, FromUOMQuantity: placement.uomQuantity, FromQuantityPerUOM: placement.quantityPerUOM, UOMCode: placement.uomCode, UOMQuantity: placement.uomQuantity, QuantityPerUOM: placement.quantityPerUOM, ExpirationDate: value.ExpirationDate, SourceEntryID: ids.EntryID})
	}
	activityUOM := placements[0].uomCode
	activityUOMQuantity := baseQuantity
	if activityUOM == handlingUOM {
		activityUOMQuantity = handlingQuantity
	} else {
		activityUOMQuantity, _ = exactUOMQuantity(baseQuantity, baseQuantityPerUOM, basePrecision)
	}
	updated, updateErr := tx.Exec(ctx, `update inventory.bulk_balance set reserved_quantity=reserved_quantity+$2::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and quantity-reserved_quantity >= $2::numeric`, balanceID, baseQuantity)
	if updateErr != nil || updated.RowsAffected() != 1 {
		if updateErr != nil {
			return result, updateErr
		}
		return result, inventorycontrol.ErrConflict
	}
	updated, updateErr = tx.Exec(ctx, `update inventory.bulk_uom_balance set reserved_quantity=reserved_quantity+$3::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and quantity-reserved_quantity >= $3::numeric`, balanceID, activityUOM, activityUOMQuantity)
	if updateErr != nil || updated.RowsAffected() != 1 {
		if updateErr != nil {
			return result, bulkConflict(updateErr)
		}
		return result, inventorycontrol.ErrConflict
	}
	result.PutAway = activity
	if err = recordBulkEvent(ctx, tx, tenant, ids.EventID, "warehouse-receipt", ids.ReceiptID, "warehouse-receipt.posted", 1, map[string]any{"organization_id": value.OrganizationID, "item_id": value.ItemID, "lot_id": lotID, "quantity": baseQuantity, "handling_uom": handlingUOM, "handling_quantity": handlingQuantity, "allow_breakbulk": value.AllowBreakbulk, "automatic_breakbulk": requiresBreakbulk, "put_away_id": ids.ActivityID, "put_away_lines": len(activity.Lines), "cross_dock_allocations": len(crossDockPlacements)}); err != nil {
		return result, err
	}
	return result, tx.Commit(ctx)
}

func (r *InventoryControl) CreateWarehousePick(ctx context.Context, tenant string, ids inventorycontrol.WarehouseIDs, value inventorycontrol.WarehousePickCommand) (inventorycontrol.WarehouseActivity, error) {
	activity := inventorycontrol.WarehouseActivity{ID: ids.ActivityID, OrganizationID: value.OrganizationID, Type: "pick", Status: "open", SourceKind: value.DemandKind, SourceID: value.DemandID, SourceLineID: value.DemandLineID, RequestID: value.RequestID, AssignedTo: value.AssignedTo, Version: 1, Lines: []inventorycontrol.WarehouseActivityLine{}}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return activity, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+"\x1fwarehouse-pick-request\x1f"+value.RequestID); err != nil {
		return activity, err
	}
	if replay, found, replayErr := readWarehousePickReplay(ctx, tx, tenant, value); replayErr != nil {
		return activity, replayErr
	} else if found {
		if err = tx.Commit(ctx); err != nil {
			return activity, err
		}
		return replay, nil
	}
	if err = warehouseLock(ctx, tx, tenant, value.OrganizationID, value.ItemID); err != nil {
		return activity, err
	}
	var tracking, targetType string
	err = tx.QueryRow(ctx, `select i.tracking_mode,w.bin_type from inventory.stock_item i join inventory.item_bin_policy p on p.tenant_id=i.tenant_id and p.item_id=i.item_id join inventory.warehouse_bin w on w.tenant_id=p.tenant_id and w.organization_id=p.organization_id and w.bin_id=p.bin_id where i.tenant_id=$1 and i.item_id=$2 and p.organization_id=$3 and p.bin_id=$4 and not w.movement_blocked for share of i,p,w`, tenant, value.ItemID, value.OrganizationID, value.ShipBinID).Scan(&tracking, &targetType)
	if errors.Is(err, pgx.ErrNoRows) {
		return activity, fmt.Errorf("warehouse pick ship-bin contract: %w", inventorycontrol.ErrConflict)
	}
	if err != nil {
		return activity, err
	}
	if targetType != "ship" {
		return activity, fmt.Errorf("warehouse pick ship-bin type: %w", inventorycontrol.ErrConflict)
	}
	demandSalesUOM, demandTotalBase, demandOutstandingBase := "", "", ""
	if value.DemandKind == "customer-order" {
		demandSalesUOM, demandTotalBase, demandOutstandingBase, err = validateCustomerOrderWarehouseDemand(ctx, tx, tenant, value)
		if err != nil {
			return activity, err
		}
	}
	if value.DemandKind == "transfer-outbound" {
		var outstandingText string
		err = tx.QueryRow(ctx, `select (l.quantity-l.shipped_quantity-coalesce((select sum(al.quantity) from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=t.tenant_id and a.organization_id=t.from_organization_id and a.activity_type='pick' and a.source_kind='transfer-outbound' and a.source_id=t.transfer_id and a.source_line_id=l.line_id and (a.status='open' or (a.status='registered' and not exists(select 1 from inventory.bulk_transfer_shipment s where s.tenant_id=a.tenant_id and s.warehouse_activity_id=a.activity_id)))),0))::text from inventory.bulk_transfer t join inventory.bulk_transfer_line l using(tenant_id,transfer_id) where t.tenant_id=$1 and t.transfer_id=$2 and l.line_id=$3 and t.from_organization_id=$4 and l.item_id=$5 and t.status in ('released','partially-shipped','partially-received') for share of t,l`, tenant, value.DemandID, value.DemandLineID, value.OrganizationID, value.ItemID).Scan(&outstandingText)
		if errors.Is(err, pgx.ErrNoRows) {
			return activity, fmt.Errorf("warehouse pick transfer demand: %w", inventorycontrol.ErrConflict)
		}
		if err != nil {
			return activity, err
		}
		outstanding, validOutstanding := parseRat(outstandingText)
		requested, validRequested := parseRat(value.Quantity)
		if !validOutstanding || !validRequested || requested.Cmp(outstanding) > 0 {
			return activity, fmt.Errorf("warehouse pick transfer outstanding: %w", inventorycontrol.ErrConflict)
		}
	}
	remaining, ok := parseRat(value.Quantity)
	if !ok {
		return activity, fmt.Errorf("warehouse pick quantity parse: %w", inventorycontrol.ErrConflict)
	}
	targetUOM, err := resolveWarehouseUOM(ctx, tx, tenant, value.ItemID, value.HandlingUOM)
	if err != nil {
		return activity, fmt.Errorf("warehouse pick handling UOM: %w", err)
	}
	if _, exact := exactUOMQuantity(value.Quantity, targetUOM.factor, targetUOM.precision); !exact {
		return activity, fmt.Errorf("warehouse pick handling UOM quantity: %w", inventorycontrol.ErrConflict)
	}
	type candidate struct {
		binID, lotID, lotNo, available, allocationID string
		balanceID                                    int64
		expiration                                   *time.Time
		plan                                         warehouseUOMPlan
	}
	candidates := []candidate{}
	appendCandidates := func(rows pgx.Rows) error {
		raw := []candidate{}
		for rows.Next() {
			var c candidate
			if scanErr := rows.Scan(&c.balanceID, &c.binID, &c.lotID, &c.lotNo, &c.expiration, &c.available, &c.allocationID); scanErr != nil {
				rows.Close()
				return scanErr
			}
			raw = append(raw, c)
		}
		rows.Close()
		if rows.Err() != nil {
			return rows.Err()
		}
		for _, c := range raw {
			if remaining.Sign() <= 0 {
				break
			}
			available, valid := parseRat(c.available)
			if !valid || available.Sign() <= 0 {
				return fmt.Errorf("warehouse pick availability parse: %w", inventorycontrol.ErrConflict)
			}
			limit := minRat(available, remaining)
			plans, leftover, planErr := planWarehousePackaging(ctx, tx, c.balanceID, targetUOM, value.AllowBreakbulk, limit)
			if planErr != nil {
				return planErr
			}
			moved := new(big.Rat).Sub(limit, leftover)
			for _, plan := range plans {
				planned := c
				planned.plan = plan
				candidates = append(candidates, planned)
			}
			remaining.Sub(remaining, moved)
		}
		return nil
	}
	if value.DemandKind == "transfer-outbound" {
		rows, queryErr := tx.Query(ctx, `select b.balance_id,x.bin_id,coalesce(x.lot_id,''),coalesce(l.lot_no,''),l.expiration_date,least(x.quantity-x.reserved_quantity-x.picked_quantity,b.quantity-b.reserved_quantity)::text,x.allocation_id from inventory.warehouse_crossdock_allocation x join inventory.bulk_balance b on b.tenant_id=x.tenant_id and b.organization_id=x.organization_id and b.bin_id=x.bin_id and b.item_id=x.item_id and b.lot_id is not distinct from x.lot_id join inventory.warehouse_bin w on w.tenant_id=x.tenant_id and w.organization_id=x.organization_id and w.bin_id=x.bin_id left join inventory.inventory_lot l on l.tenant_id=x.tenant_id and l.lot_id=x.lot_id where x.tenant_id=$1 and x.organization_id=$2 and x.item_id=$3 and x.transfer_id=$4 and x.transfer_line_id=$5 and x.available and x.quantity-x.reserved_quantity-x.picked_quantity>0 and b.quantity-b.reserved_quantity>0 and w.cross_dock and not w.movement_blocked and coalesce(l.blocked,false)=false and (l.expiration_date is null or l.expiration_date>=current_date) order by x.due_date,w.bin_rank desc,w.bin_code,l.expiration_date nulls last,l.lot_no for update of x,b`, tenant, value.OrganizationID, value.ItemID, value.DemandID, value.DemandLineID)
		if queryErr != nil {
			return activity, queryErr
		}
		if err = appendCandidates(rows); err != nil {
			return activity, err
		}
	}
	if remaining.Sign() > 0 {
		query := `select b.balance_id,b.bin_id,coalesce(b.lot_id,''),coalesce(l.lot_no,''),l.expiration_date,(b.quantity-b.reserved_quantity)::text,'' from inventory.bulk_balance b join inventory.item_bin_policy p using(tenant_id,organization_id,item_id,bin_id) join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) left join inventory.inventory_lot l on l.tenant_id=b.tenant_id and l.lot_id=b.lot_id where b.tenant_id=$1 and b.organization_id=$2 and b.item_id=$3 and b.quantity-b.reserved_quantity>0 and w.bin_type in ('pick','putpick') and not w.cross_dock and not w.movement_blocked and ($4 or not p.dedicated) and coalesce(l.blocked,false)=false and (l.expiration_date is null or l.expiration_date>=current_date)`
		if tracking == "lot" && value.UseFEFO {
			query += ` order by l.expiration_date nulls last,l.lot_no,l.created_at,w.bin_rank desc,w.bin_code`
		} else {
			query += ` order by w.bin_rank desc,w.bin_code,l.expiration_date nulls last,l.lot_no`
		}
		query += ` for update of b`
		rows, queryErr := tx.Query(ctx, query, tenant, value.OrganizationID, value.ItemID, value.AllowDedicated)
		if queryErr != nil {
			return activity, queryErr
		}
		if err = appendCandidates(rows); err != nil {
			return activity, err
		}
	}
	if remaining.Sign() != 0 || len(candidates) == 0 {
		return activity, fmt.Errorf("warehouse pick insufficient FEFO availability: %w", inventorycontrol.ErrConflict)
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity(tenant_id,activity_id,organization_id,activity_type,status,source_kind,source_id,source_line_id,request_id,assigned_to,version) values($1,$2,$3,'pick','open',$4,$5,$6,$7,$8,1)`, tenant, ids.ActivityID, value.OrganizationID, value.DemandKind, value.DemandID, value.DemandLineID, value.RequestID, value.AssignedTo)
	if err != nil {
		return activity, bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_pick_request(tenant_id,request_id,activity_id,organization_id,ship_bin_id,item_id,demand_kind,demand_id,demand_line_id,quantity,requested_handling_uom,resolved_handling_uom,allow_breakbulk,use_fefo,allow_dedicated,assigned_to) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10::numeric,$11,$12,$13,$14,$15,$16)`, tenant, value.RequestID, ids.ActivityID, value.OrganizationID, value.ShipBinID, value.ItemID, value.DemandKind, value.DemandID, value.DemandLineID, value.Quantity, value.HandlingUOM, targetUOM.code, value.AllowBreakbulk, value.UseFEFO, value.AllowDedicated, value.AssignedTo)
	if err != nil {
		return activity, bulkConflict(err)
	}
	for index, c := range candidates {
		lineID := fmt.Sprintf("%s-%06d", ids.ActivityID, index+1)
		reservationID := fmt.Sprintf("%s-r-%06d", ids.ActivityID, index+1)
		_, err = tx.Exec(ctx, `insert into inventory.bulk_reservation(tenant_id,reservation_id,organization_id,bin_id,item_id,lot_id,demand_kind,demand_id,demand_line_id,quantity,status,cancellation_disallowed,version) values($1,$2,$3,$4,$5,nullif($6,''),$7,$8,$9,$10::numeric,'reservation',true,1)`, tenant, reservationID, value.OrganizationID, c.binID, value.ItemID, c.lotID, value.DemandKind, value.DemandID, fmt.Sprintf("%s:%s:%06d", value.DemandLineID, ids.ActivityID, index+1), c.plan.movedBase)
		if err != nil {
			return activity, bulkConflict(err)
		}
		_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity_line(tenant_id,activity_id,organization_id,line_id,sequence_no,from_bin_id,to_bin_id,item_id,lot_id,reservation_id,quantity,expiration_date,uom_code,uom_quantity,qty_per_uom,from_uom_code,from_uom_quantity,from_qty_per_uom) values($1,$2,$3,$4,$5,$6,$7,$8,nullif($9,''),$10,$11::numeric,$12,$13,$14::numeric,$15::numeric,$16,$17::numeric,$18::numeric)`, tenant, ids.ActivityID, value.OrganizationID, lineID, index+1, c.binID, value.ShipBinID, value.ItemID, c.lotID, reservationID, c.plan.movedBase, c.expiration, c.plan.targetUOM, c.plan.targetQuantity, c.plan.targetFactor, c.plan.sourceUOM, c.plan.sourceQuantity, c.plan.sourceFactor)
		if err != nil {
			return activity, bulkConflict(err)
		}
		line := inventorycontrol.WarehouseActivityLine{ID: lineID, Sequence: index + 1, FromBinID: c.binID, ToBinID: value.ShipBinID, ItemID: value.ItemID, LotID: c.lotID, LotNo: c.lotNo, ReservationID: reservationID, ReservationVersion: 1, Quantity: c.plan.movedBase, FromUOMCode: c.plan.sourceUOM, FromUOMQuantity: c.plan.sourceQuantity, FromQuantityPerUOM: c.plan.sourceFactor, UOMCode: c.plan.targetUOM, UOMQuantity: c.plan.targetQuantity, QuantityPerUOM: c.plan.targetFactor, ExpirationDate: c.expiration}
		if err = reserveWarehousePackaging(ctx, tx, tenant, ids.ActivityID, value.OrganizationID, c.plan, line); err != nil {
			return activity, err
		}
		if c.allocationID != "" {
			updated, err := tx.Exec(ctx, `update inventory.warehouse_crossdock_allocation set reserved_quantity=reserved_quantity+$4::numeric where tenant_id=$1 and allocation_id=$2 and available and transfer_id=$3 and quantity-reserved_quantity-picked_quantity >= $4::numeric`, tenant, c.allocationID, value.DemandID, c.plan.movedBase)
			if err != nil || updated.RowsAffected() != 1 {
				if err != nil {
					return activity, err
				}
				return activity, inventorycontrol.ErrConflict
			}
			_, err = tx.Exec(ctx, `insert into inventory.warehouse_crossdock_pick_link(tenant_id,pick_activity_id,pick_line_id,allocation_id,quantity,status) values($1,$2,$3,$4,$5::numeric,'open')`, tenant, ids.ActivityID, lineID, c.allocationID, c.plan.movedBase)
			if err != nil {
				return activity, bulkConflict(err)
			}
		}
		activity.Lines = append(activity.Lines, line)
	}
	if err = recordBulkEvent(ctx, tx, tenant, ids.EventID, "warehouse-pick", ids.ActivityID, "warehouse-pick.created", 1, map[string]any{"organization_id": value.OrganizationID, "item_id": value.ItemID, "quantity": value.Quantity, "handling_uom": targetUOM.code, "allow_breakbulk": value.AllowBreakbulk, "use_fefo": value.UseFEFO, "lines": len(activity.Lines), "demand_sales_uom": demandSalesUOM, "demand_total_base": demandTotalBase, "demand_outstanding_before_base": demandOutstandingBase}); err != nil {
		return activity, err
	}
	return activity, tx.Commit(ctx)
}

func readWarehouseActivity(ctx context.Context, tx pgx.Tx, tenant, organization, activityID string) (inventorycontrol.WarehouseActivity, error) {
	var value inventorycontrol.WarehouseActivity
	err := tx.QueryRow(ctx, `select activity_id,organization_id,activity_type,status,source_kind,source_id,source_line_id,request_id,assigned_to,version from inventory.warehouse_activity where tenant_id=$1 and organization_id=$2 and activity_id=$3`, tenant, organization, activityID).Scan(&value.ID, &value.OrganizationID, &value.Type, &value.Status, &value.SourceKind, &value.SourceID, &value.SourceLineID, &value.RequestID, &value.AssignedTo, &value.Version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return value, inventorycontrol.ErrConflict
		}
		return value, err
	}
	rows, err := tx.Query(ctx, `select al.line_id,al.sequence_no,al.from_bin_id,al.to_bin_id,al.item_id,coalesce(al.lot_id,''),coalesce(l.lot_no,''),coalesce(al.reservation_id,''),al.quantity::text,al.from_uom_code,al.from_uom_quantity::text,al.from_qty_per_uom::text,al.uom_code,al.uom_quantity::text,al.qty_per_uom::text,al.expiration_date,coalesce(al.source_entry_id,'') from inventory.warehouse_activity_line al left join inventory.inventory_lot l on l.tenant_id=al.tenant_id and l.lot_id=al.lot_id where al.tenant_id=$1 and al.activity_id=$2 order by al.sequence_no`, tenant, activityID)
	if err != nil {
		return value, err
	}
	defer rows.Close()
	value.Lines = []inventorycontrol.WarehouseActivityLine{}
	for rows.Next() {
		var line inventorycontrol.WarehouseActivityLine
		if err = rows.Scan(&line.ID, &line.Sequence, &line.FromBinID, &line.ToBinID, &line.ItemID, &line.LotID, &line.LotNo, &line.ReservationID, &line.Quantity, &line.FromUOMCode, &line.FromUOMQuantity, &line.FromQuantityPerUOM, &line.UOMCode, &line.UOMQuantity, &line.QuantityPerUOM, &line.ExpirationDate, &line.SourceEntryID); err != nil {
			return value, err
		}
		value.Lines = append(value.Lines, line)
	}
	return value, rows.Err()
}

func releasePutAwaySourceComposition(ctx context.Context, tx pgx.Tx, tenant, organization string, line inventorycontrol.WarehouseActivityLine) error {
	var balanceID int64
	var isBase bool
	err := tx.QueryRow(ctx, `select b.balance_id,u.is_base from inventory.bulk_balance b join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.uom_code=$7 where b.tenant_id=$1 and b.organization_id=$2 and b.bin_id=$3 and b.item_id=$4 and b.lot_id is not distinct from nullif($5,'') and b.quantity >= $6::numeric and b.reserved_quantity >= $6::numeric for update of b,u`, tenant, organization, line.FromBinID, line.ItemID, line.LotID, line.Quantity, line.UOMCode).Scan(&balanceID, &isBase)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.ErrConflict
	}
	if err != nil {
		return err
	}
	if isBase {
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_uom_balance set reserved_quantity=reserved_quantity-$3::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and reserved_quantity >= $3::numeric`, balanceID, line.UOMCode, line.UOMQuantity)
		if updateErr != nil {
			return bulkConflict(updateErr)
		}
		if updated.RowsAffected() != 1 {
			return inventorycontrol.ErrConflict
		}
		return nil
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_uom_balance set quantity=quantity-$3::numeric,reserved_quantity=reserved_quantity-$3::numeric,quantity_base=quantity_base-$4::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and quantity >= $3::numeric and reserved_quantity >= $3::numeric and quantity_base >= $4::numeric`, balanceID, line.UOMCode, line.UOMQuantity, line.Quantity)
	if err != nil {
		return bulkConflict(err)
	}
	if updated.RowsAffected() != 1 {
		return inventorycontrol.ErrConflict
	}
	if _, err = tx.Exec(ctx, `delete from inventory.bulk_uom_balance where balance_id=$1 and uom_code=$2 and quantity_base=0 and reserved_quantity=0`, balanceID, line.UOMCode); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base) select $1,u.uom_code,$2::numeric/u.qty_per_uom,u.qty_per_uom,$2::numeric from inventory.bulk_balance b join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.is_base where b.balance_id=$1 on conflict(balance_id,uom_code) do update set quantity=inventory.bulk_uom_balance.quantity+excluded.quantity,quantity_base=inventory.bulk_uom_balance.quantity_base+excluded.quantity_base,version=inventory.bulk_uom_balance.version+1,updated_at=clock_timestamp()`, balanceID, line.Quantity)
	return bulkConflict(err)
}

func preservePutAwayTargetComposition(ctx context.Context, tx pgx.Tx, balanceID int64, line inventorycontrol.WarehouseActivityLine) error {
	var isBase bool
	err := tx.QueryRow(ctx, `select is_base from inventory.item_unit_of_measure u join inventory.bulk_balance b on b.tenant_id=u.tenant_id and b.item_id=u.item_id where b.balance_id=$1 and u.uom_code=$2`, balanceID, line.UOMCode).Scan(&isBase)
	if err != nil || isBase {
		return err
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_uom_balance target set quantity=target.quantity-($2::numeric/target.qty_per_uom),quantity_base=target.quantity_base-$2::numeric,version=target.version+1,updated_at=clock_timestamp() where target.balance_id=$1 and target.uom_code=(select u.uom_code from inventory.bulk_balance b join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.is_base where b.balance_id=$1) and target.quantity-target.reserved_quantity >= ($2::numeric/target.qty_per_uom) and target.quantity_base >= $2::numeric`, balanceID, line.Quantity)
	if err != nil {
		return bulkConflict(err)
	}
	if updated.RowsAffected() != 1 {
		return inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `delete from inventory.bulk_uom_balance target where target.balance_id=$1 and target.quantity_base=0 and target.reserved_quantity=0 and target.uom_code=(select u.uom_code from inventory.bulk_balance b join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.is_base where b.balance_id=$1)`, balanceID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base) values($1,$2,$3::numeric,$4::numeric,$5::numeric) on conflict(balance_id,uom_code) do update set quantity=inventory.bulk_uom_balance.quantity+excluded.quantity,quantity_base=inventory.bulk_uom_balance.quantity_base+excluded.quantity_base,version=inventory.bulk_uom_balance.version+1,updated_at=clock_timestamp()`, balanceID, line.UOMCode, line.UOMQuantity, line.QuantityPerUOM, line.Quantity)
	return bulkConflict(err)
}

func (r *InventoryControl) RegisterWarehouseActivity(ctx context.Context, tenant, organization, activityID string, version int64, postingDate time.Time, eventID string) (inventorycontrol.WarehouseActivity, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	defer tx.Rollback(ctx)
	var activityType, itemID string
	err = tx.QueryRow(ctx, `select a.activity_type,al.item_id from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=$1 and a.organization_id=$2 and a.activity_id=$3 and a.status='open' and a.version=$4 order by al.sequence_no limit 1 for update of a`, tenant, organization, activityID, version).Scan(&activityType, &itemID)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	if err = warehouseLock(ctx, tx, tenant, organization, itemID); err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	activity, err := readWarehouseActivity(ctx, tx, tenant, organization, activityID)
	if err != nil {
		return activity, err
	}
	for index, line := range activity.Lines {
		if activityType == "put-away" {
			if err = releasePutAwaySourceComposition(ctx, tx, tenant, organization, line); err != nil {
				return activity, err
			}
			updated, updateErr := tx.Exec(ctx, `update inventory.bulk_balance set quantity=quantity-$6::numeric,reserved_quantity=reserved_quantity-$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and quantity >= $6::numeric and reserved_quantity >= $6::numeric`, tenant, organization, line.FromBinID, line.ItemID, line.LotID, line.Quantity)
			if updateErr != nil {
				return activity, updateErr
			}
			if updated.RowsAffected() != 1 {
				return activity, inventorycontrol.ErrConflict
			}
			var targetBalanceID, targetVersion int64
			err = tx.QueryRow(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity,version) values($1,$2,$3,$4,nullif($5,''),$6::numeric,0,1) on conflict (tenant_id,organization_id,bin_id,item_id,lot_id) do update set quantity=inventory.bulk_balance.quantity+excluded.quantity,version=inventory.bulk_balance.version+1,updated_at=clock_timestamp() where exists(select 1 from inventory.item_bin_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) where p.tenant_id=$1 and p.organization_id=$2 and p.item_id=$4 and p.bin_id=$3 and w.bin_type in ('put-away','putpick') and not w.movement_blocked and (p.max_quantity is null or inventory.bulk_balance.quantity+excluded.quantity<=p.max_quantity)) returning balance_id,version`, tenant, organization, line.ToBinID, line.ItemID, line.LotID, line.Quantity).Scan(&targetBalanceID, &targetVersion)
			if errors.Is(err, pgx.ErrNoRows) {
				return activity, inventorycontrol.ErrConflict
			}
			if err != nil {
				return activity, bulkConflict(err)
			}
			if err = preservePutAwayTargetComposition(ctx, tx, targetBalanceID, line); err != nil {
				return activity, err
			}
			_, err = tx.Exec(ctx, `update inventory.warehouse_crossdock_allocation set available=true where tenant_id=$1 and put_away_activity_id=$2 and put_away_line_id=$3 and not available`, tenant, activityID, line.ID)
			if err != nil {
				return activity, err
			}
		} else if activityType == "pick" {
			if err = moveWarehousePackagingSource(ctx, tx, tenant, organization, activityID, "warehouse-pick", line); err != nil {
				return activity, err
			}
			activity.Lines[index].ReservationVersion = 2
			var targetBalanceID, targetVersion int64
			err = tx.QueryRow(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity,version) values($1,$2,$3,$4,nullif($5,''),$6::numeric,$6::numeric,1) on conflict (tenant_id,organization_id,bin_id,item_id,lot_id) do update set quantity=inventory.bulk_balance.quantity+excluded.quantity,reserved_quantity=inventory.bulk_balance.reserved_quantity+excluded.reserved_quantity,version=inventory.bulk_balance.version+1,updated_at=clock_timestamp() where exists(select 1 from inventory.item_bin_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) where p.tenant_id=$1 and p.organization_id=$2 and p.item_id=$4 and p.bin_id=$3 and w.bin_type='ship' and not w.movement_blocked and (p.max_quantity is null or inventory.bulk_balance.quantity+excluded.quantity<=p.max_quantity)) returning balance_id,version`, tenant, organization, line.ToBinID, line.ItemID, line.LotID, line.Quantity).Scan(&targetBalanceID, &targetVersion)
			if errors.Is(err, pgx.ErrNoRows) {
				return activity, inventorycontrol.ErrConflict
			}
			if err != nil {
				return activity, bulkConflict(err)
			}
			if err = preservePutAwayTargetComposition(ctx, tx, targetBalanceID, line); err != nil {
				return activity, err
			}
			updated, err := tx.Exec(ctx, `update inventory.bulk_reservation set bin_id=$4,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and reservation_id=$3 and status='reservation' and version=1`, tenant, organization, line.ReservationID, line.ToBinID)
			if err != nil {
				return activity, err
			}
			if updated.RowsAffected() != 1 {
				return activity, inventorycontrol.ErrConflict
			}
			var allocationID string
			linkErr := tx.QueryRow(ctx, `select allocation_id from inventory.warehouse_crossdock_pick_link where tenant_id=$1 and pick_activity_id=$2 and pick_line_id=$3 and status='open' for update`, tenant, activityID, line.ID).Scan(&allocationID)
			if linkErr == nil {
				updated, err = tx.Exec(ctx, `update inventory.warehouse_crossdock_allocation set reserved_quantity=reserved_quantity-$3::numeric,picked_quantity=picked_quantity+$3::numeric where tenant_id=$1 and allocation_id=$2 and reserved_quantity >= $3::numeric and reserved_quantity+picked_quantity <= quantity`, tenant, allocationID, line.Quantity)
				if err != nil || updated.RowsAffected() != 1 {
					if err != nil {
						return activity, err
					}
					return activity, inventorycontrol.ErrConflict
				}
				updated, err = tx.Exec(ctx, `update inventory.warehouse_crossdock_pick_link set status='picked' where tenant_id=$1 and pick_activity_id=$2 and pick_line_id=$3 and status='open'`, tenant, activityID, line.ID)
				if err != nil || updated.RowsAffected() != 1 {
					if err != nil {
						return activity, err
					}
					return activity, inventorycontrol.ErrConflict
				}
			} else if !errors.Is(linkErr, pgx.ErrNoRows) {
				return activity, linkErr
			}
		} else if activityType == "movement" {
			if err = moveWarehousePackagingSource(ctx, tx, tenant, organization, activityID, "warehouse-replenishment", line); err != nil {
				return activity, err
			}
			var targetBalanceID, targetVersion int64
			err = tx.QueryRow(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity,version) values($1,$2,$3,$4,nullif($5,''),$6::numeric,0,1) on conflict (tenant_id,organization_id,bin_id,item_id,lot_id) do update set quantity=inventory.bulk_balance.quantity+excluded.quantity,version=inventory.bulk_balance.version+1,updated_at=clock_timestamp() where exists(select 1 from inventory.item_bin_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) where p.tenant_id=$1 and p.organization_id=$2 and p.item_id=$4 and p.bin_id=$3 and p.fixed and w.bin_type in ('pick','putpick') and not w.movement_blocked and p.max_quantity is not null and (select coalesce(sum(quantity),0) from inventory.bulk_balance total where total.tenant_id=$1 and total.organization_id=$2 and total.bin_id=$3 and total.item_id=$4)+excluded.quantity<=p.max_quantity) returning balance_id,version`, tenant, organization, line.ToBinID, line.ItemID, line.LotID, line.Quantity).Scan(&targetBalanceID, &targetVersion)
			if errors.Is(err, pgx.ErrNoRows) {
				return activity, inventorycontrol.ErrConflict
			}
			if err != nil {
				return activity, bulkConflict(err)
			}
			if err = preservePutAwayTargetComposition(ctx, tx, targetBalanceID, line); err != nil {
				return activity, err
			}
			updated, err := tx.Exec(ctx, `update inventory.bulk_reservation set status='consumed',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and reservation_id=$3 and status='reservation' and version=1`, tenant, organization, line.ReservationID)
			if err != nil || updated.RowsAffected() != 1 {
				if err != nil {
					return activity, err
				}
				return activity, inventorycontrol.ErrConflict
			}
		} else {
			return activity, inventorycontrol.ErrConflict
		}
		outID := fmt.Sprintf("%s-out-%06d", eventID, index+1)
		inID := fmt.Sprintf("%s-in-%06d", eventID, index+1)
		_, err = tx.Exec(ctx, `insert into inventory.bulk_inventory_entry(tenant_id,entry_id,organization_id,bin_id,item_id,lot_id,entry_type,quantity,unit_cost,cost_amount,posting_date,source_kind,source_id) values($1,$2,$3,$4,$5,nullif($6,''),'movement-out',-$7::numeric,0,0,$8::date,$9,$10),($1,$11,$3,$12,$5,nullif($6,''),'movement-in',$7::numeric,0,0,$8::date,$9,$10)`, tenant, outID, organization, line.FromBinID, line.ItemID, line.LotID, line.Quantity, postingDate, "warehouse-"+activityType, line.ID, inID, line.ToBinID)
		if err != nil {
			return activity, bulkConflict(err)
		}
	}
	updated, err := tx.Exec(ctx, `update inventory.warehouse_activity set status='registered',version=version+1,registered_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and activity_id=$3 and status='open' and version=$4`, tenant, organization, activityID, version)
	if err != nil {
		return activity, err
	}
	if updated.RowsAffected() != 1 {
		return activity, inventorycontrol.ErrConflict
	}
	activity.Status, activity.Version = "registered", version+1
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "warehouse-activity", activityID, "warehouse-activity.registered", activity.Version, map[string]any{"organization_id": organization, "activity_type": activityType, "lines": len(activity.Lines)}); err != nil {
		return activity, err
	}
	return activity, tx.Commit(ctx)
}

func (r *InventoryControl) CancelWarehousePick(ctx context.Context, tenant, organization, activityID string, version int64, eventID string) (inventorycontrol.WarehouseActivity, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	defer tx.Rollback(ctx)
	var itemID string
	err = tx.QueryRow(ctx, `select al.item_id from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=$1 and a.organization_id=$2 and a.activity_id=$3 and a.activity_type='pick' and a.status='open' and a.version=$4 order by al.sequence_no limit 1 for update of a`, tenant, organization, activityID, version).Scan(&itemID)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	if err = warehouseLock(ctx, tx, tenant, organization, itemID); err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	activity, err := readWarehouseActivity(ctx, tx, tenant, organization, activityID)
	if err != nil {
		return activity, err
	}
	for _, line := range activity.Lines {
		if err = releaseWarehousePackagingReservation(ctx, tx, tenant, organization, activityID, line); err != nil {
			return activity, err
		}
		updated, err := tx.Exec(ctx, `update inventory.bulk_reservation set status='released',cancellation_disallowed=false,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and reservation_id=$3 and status='reservation' and version=1`, tenant, organization, line.ReservationID)
		if err != nil {
			return activity, err
		}
		if updated.RowsAffected() != 1 {
			return activity, inventorycontrol.ErrConflict
		}
		var allocationID string
		linkErr := tx.QueryRow(ctx, `select allocation_id from inventory.warehouse_crossdock_pick_link where tenant_id=$1 and pick_activity_id=$2 and pick_line_id=$3 and status='open' for update`, tenant, activityID, line.ID).Scan(&allocationID)
		if linkErr == nil {
			updated, err = tx.Exec(ctx, `update inventory.warehouse_crossdock_allocation set reserved_quantity=reserved_quantity-$3::numeric where tenant_id=$1 and allocation_id=$2 and reserved_quantity >= $3::numeric`, tenant, allocationID, line.Quantity)
			if err != nil || updated.RowsAffected() != 1 {
				if err != nil {
					return activity, err
				}
				return activity, inventorycontrol.ErrConflict
			}
			updated, err = tx.Exec(ctx, `update inventory.warehouse_crossdock_pick_link set status='released' where tenant_id=$1 and pick_activity_id=$2 and pick_line_id=$3 and status='open'`, tenant, activityID, line.ID)
			if err != nil || updated.RowsAffected() != 1 {
				if err != nil {
					return activity, err
				}
				return activity, inventorycontrol.ErrConflict
			}
		} else if !errors.Is(linkErr, pgx.ErrNoRows) {
			return activity, linkErr
		}
	}
	updated, err := tx.Exec(ctx, `update inventory.warehouse_activity set status='cancelled',version=version+1 where tenant_id=$1 and organization_id=$2 and activity_id=$3 and status='open' and version=$4`, tenant, organization, activityID, version)
	if err != nil {
		return activity, err
	}
	if updated.RowsAffected() != 1 {
		return activity, inventorycontrol.ErrConflict
	}
	activity.Status, activity.Version = "cancelled", version+1
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "warehouse-pick", activityID, "warehouse-pick.cancelled", activity.Version, map[string]any{"organization_id": organization, "lines": len(activity.Lines)}); err != nil {
		return activity, err
	}
	return activity, tx.Commit(ctx)
}

func (r *InventoryControl) CancelWarehousePutAway(ctx context.Context, tenant, organization, activityID string, version int64, eventID string) (inventorycontrol.WarehouseActivity, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	defer tx.Rollback(ctx)
	var itemID string
	err = tx.QueryRow(ctx, `select al.item_id from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=$1 and a.organization_id=$2 and a.activity_id=$3 and a.activity_type='put-away' and a.status='open' and a.version=$4 order by al.sequence_no limit 1 for update of a`, tenant, organization, activityID, version).Scan(&itemID)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	if err = warehouseLock(ctx, tx, tenant, organization, itemID); err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	activity, err := readWarehouseActivity(ctx, tx, tenant, organization, activityID)
	if err != nil {
		return activity, err
	}
	for _, line := range activity.Lines {
		var balanceID int64
		err = tx.QueryRow(ctx, `update inventory.bulk_balance set reserved_quantity=reserved_quantity-$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and reserved_quantity >= $6::numeric returning balance_id`, tenant, organization, line.FromBinID, line.ItemID, line.LotID, line.Quantity).Scan(&balanceID)
		if errors.Is(err, pgx.ErrNoRows) {
			return activity, inventorycontrol.ErrConflict
		}
		if err != nil {
			return activity, err
		}
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_uom_balance set reserved_quantity=reserved_quantity-$3::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and reserved_quantity >= $3::numeric`, balanceID, line.UOMCode, line.UOMQuantity)
		if updateErr != nil {
			return activity, bulkConflict(updateErr)
		}
		if updated.RowsAffected() != 1 {
			return activity, inventorycontrol.ErrConflict
		}
	}
	updated, err := tx.Exec(ctx, `update inventory.warehouse_activity set status='cancelled',version=version+1 where tenant_id=$1 and organization_id=$2 and activity_id=$3 and activity_type='put-away' and status='open' and version=$4`, tenant, organization, activityID, version)
	if err != nil {
		return activity, err
	}
	if updated.RowsAffected() != 1 {
		return activity, inventorycontrol.ErrConflict
	}
	activity.Status, activity.Version = "cancelled", version+1
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "warehouse-put-away", activityID, "warehouse-put-away.cancelled", activity.Version, map[string]any{"organization_id": organization, "lines": len(activity.Lines)}); err != nil {
		return activity, err
	}
	return activity, tx.Commit(ctx)
}
````

### FILE: `internal/platform/postgres/warehouse_integration_test.go`

```yaml
block_id: "GO-OPS-API:operational-warehouse-postgres-test:v1"
operation: CREATE
provenance: ADAPTED
source: "Executable regression for Microsoft-governed ranking, capacity, FEFO, Take/Place, reservation and immutability invariants"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "f566c8f6d188d3c2de2d432922b59e8359f3df6e741faf5fa407527a5e290d91"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func warehouseTestIDs(t *testing.T, prefix string) inventorycontrol.WarehouseIDs {
	t.Helper()
	return inventorycontrol.WarehouseIDs{ActivityID: prefix + "-activity", ReceiptID: prefix + "-receipt", LineID: prefix + "-line", EntryID: prefix + "-entry", LayerID: prefix + "-layer", LotID: prefix + "-lot", EventID: bulkTestUUID(t)}
}

func TestWarehouseReceiptPutAwayFEFOPickAndCancellation(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := bulkTestUUID(t)
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Warehouse V161','Warehouse V161')`, tenant, "wh-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	defer func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			t.Errorf("warehouse cleanup begin: %v", cleanupErr)
			return
		}
		defer tx.Rollback(ctx)
		commands := []struct {
			sql  string
			args []any
		}{
			{`alter table inventory.warehouse_pick_request disable trigger warehouse_pick_request_immutable`, nil},
			{`alter table inventory.warehouse_activity_uom_reservation disable trigger warehouse_uom_reservation_delete_guard`, nil},
			{`alter table inventory.warehouse_activity_line disable trigger warehouse_activity_line_immutable`, nil},
			{`alter table inventory.warehouse_receipt_line disable trigger warehouse_receipt_line_immutable`, nil},
			{`alter table inventory.warehouse_receipt disable trigger warehouse_receipt_immutable`, nil},
			{`alter table inventory.bulk_cost_application disable trigger bulk_cost_application_immutable`, nil},
			{`alter table inventory.bulk_inventory_entry disable trigger bulk_inventory_entry_immutable`, nil},
			{`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`, nil},
			{`delete from inventory.warehouse_activity_uom_reservation where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.warehouse_pick_request where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.warehouse_activity_line where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.warehouse_activity where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.warehouse_receipt_line where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.warehouse_receipt where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_cost_application where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_cost_layer where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_inventory_entry where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_reservation where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_balance where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.inventory_lot where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.item_bin_policy where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.warehouse_bin where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.item_unit_of_measure where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.stock_item where tenant_id=$1`, []any{tenant}},
			{`delete from platform.outbox_event where tenant_id=$1`, []any{tenant}},
			{`delete from org.organization where tenant_id=$1`, []any{tenant}},
			{`delete from platform.tenant where tenant_id=$1`, []any{tenant}},
			{`alter table inventory.warehouse_activity_line enable trigger warehouse_activity_line_immutable`, nil},
			{`alter table inventory.warehouse_receipt_line enable trigger warehouse_receipt_line_immutable`, nil},
			{`alter table inventory.warehouse_receipt enable trigger warehouse_receipt_immutable`, nil},
			{`alter table inventory.bulk_cost_application enable trigger bulk_cost_application_immutable`, nil},
			{`alter table inventory.bulk_inventory_entry enable trigger bulk_inventory_entry_immutable`, nil},
			{`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`, nil},
			{`alter table inventory.warehouse_activity_uom_reservation enable trigger warehouse_uom_reservation_delete_guard`, nil},
			{`alter table inventory.warehouse_pick_request enable trigger warehouse_pick_request_immutable`, nil},
		}
		for _, command := range commands {
			if _, cleanupErr = tx.Exec(ctx, command.sql, command.args...); cleanupErr != nil {
				t.Errorf("warehouse cleanup failed: %v", cleanupErr)
				return
			}
		}
		if cleanupErr = tx.Commit(ctx); cleanupErr != nil {
			t.Errorf("warehouse cleanup commit: %v", cleanupErr)
		}
	}()
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'warehouse','warehouse','Warehouse','warehouse')`, tenant); err != nil {
		t.Fatal(err)
	}
	repo := NewInventoryControl(pool)
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "part", Code: "PART-FEFO", Description: "Tracked part", BaseUOM: "EA", BaseRoundingPrecision: "0.000001", TrackingMode: "lot", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	bins := []inventorycontrol.WarehouseBin{
		{ID: "receive", OrganizationID: "warehouse", Code: "RECEIVE", Type: "receive", Ranking: 0, Version: 1},
		{ID: "bulk", OrganizationID: "warehouse", Code: "BULK", Type: "putpick", Ranking: 10, Version: 1},
		{ID: "forward", OrganizationID: "warehouse", Code: "FORWARD", Type: "putpick", Ranking: 100, Version: 1},
		{ID: "ship", OrganizationID: "warehouse", Code: "SHIP", Type: "ship", Ranking: 0, Version: 1},
		{ID: "qc", OrganizationID: "warehouse", Code: "QC", Type: "qc", Ranking: 1000, Version: 1},
	}
	for _, bin := range bins {
		if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatal(err)
		}
		maximum := "100"
		if bin.ID == "forward" {
			maximum = "6"
		}
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: "warehouse", ItemID: "part", BinID: bin.ID, Fixed: true, Default: bin.ID == "forward", MinQuantity: "0", MaxQuantity: maximum, Version: 1}); err != nil {
			t.Fatal(err)
		}
	}
	var configuredTracking, configuredBinType string
	if err = pool.QueryRow(ctx, `select i.tracking_mode,w.bin_type from inventory.stock_item i join inventory.item_bin_policy p on p.tenant_id=i.tenant_id and p.item_id=i.item_id join inventory.warehouse_bin w on w.tenant_id=p.tenant_id and w.organization_id=p.organization_id and w.bin_id=p.bin_id where i.tenant_id=$1 and i.item_id='part' and p.organization_id='warehouse' and p.bin_id='receive'`, tenant).Scan(&configuredTracking, &configuredBinType); err != nil || configuredTracking != "lot" || configuredBinType != "receive" {
		t.Fatalf("configured receipt boundary tracking=%s bin=%s err=%v", configuredTracking, configuredBinType, err)
	}
	late := time.Date(2035, 12, 31, 0, 0, 0, 0, time.UTC)
	early := time.Date(2034, 1, 1, 0, 0, 0, 0, time.UTC)
	receipts := []struct {
		prefix, lot, quantity, cost string
		expiry                      *time.Time
	}{
		{"late", "LOT-LATE", "6", "100", &late},
		{"early", "LOT-EARLY", "4", "120", &early},
	}
	activities := make([]inventorycontrol.WarehouseActivity, 0, len(receipts))
	for index, receipt := range receipts {
		result, receiveErr := repo.PostWarehouseReceipt(ctx, tenant, warehouseTestIDs(t, receipt.prefix), inventorycontrol.WarehouseReceiptCommand{RequestID: receipt.prefix + "-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "part", LotNo: receipt.lot, ExpirationDate: receipt.expiry, Quantity: receipt.quantity, UnitCost: receipt.cost, PostingDate: time.Date(2030, 1, index+1, 0, 0, 0, 0, time.UTC), SourceKind: "purchase", SourceID: receipt.prefix + "-po"})
		if receiveErr != nil {
			t.Fatalf("receipt %s: %v", receipt.prefix, receiveErr)
		}
		if len(result.PutAway.Lines) != 1 {
			t.Fatalf("put-away lines=%+v", result.PutAway.Lines)
		}
		activities = append(activities, result.PutAway)
	}
	if activities[0].Lines[0].ToBinID != "forward" || activities[1].Lines[0].ToBinID != "bulk" {
		t.Fatalf("ranking/capacity placements=%s,%s", activities[0].Lines[0].ToBinID, activities[1].Lines[0].ToBinID)
	}
	for index, activity := range activities {
		registered, registerErr := repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", activity.ID, 1, time.Date(2030, 1, index+3, 0, 0, 0, 0, time.UTC), bulkTestUUID(t))
		if registerErr != nil || registered.Status != "registered" || registered.Version != 2 {
			t.Fatalf("register=%+v err=%v", registered, registerErr)
		}
	}
	var pickableQuantity string
	if err = pool.QueryRow(ctx, `select coalesce(sum(b.quantity-b.reserved_quantity),0)::text from inventory.bulk_balance b join inventory.item_bin_policy p using(tenant_id,organization_id,item_id,bin_id) join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) left join inventory.inventory_lot l on l.tenant_id=b.tenant_id and l.lot_id=b.lot_id where b.tenant_id=$1 and b.organization_id='warehouse' and b.item_id='part' and w.bin_type in ('pick','putpick') and not w.movement_blocked and not p.dedicated and coalesce(l.blocked,false)=false and (l.expiration_date is null or l.expiration_date>=current_date)`, tenant).Scan(&pickableQuantity); err != nil || pickableQuantity != "10.000000" {
		t.Fatalf("pickable quantity=%s err=%v", pickableQuantity, err)
	}
	pick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "pick"), inventorycontrol.WarehousePickCommand{RequestID: "pick-request", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "manual", DemandID: "order-1", DemandLineID: "line-1", Quantity: "7", UseFEFO: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(pick.Lines) != 2 || pick.Lines[0].LotNo != "LOT-EARLY" || pick.Lines[0].Quantity != "4" || pick.Lines[1].LotNo != "LOT-LATE" || pick.Lines[1].Quantity != "3" {
		t.Fatalf("FEFO lines=%+v", pick.Lines)
	}
	registeredPick, err := repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", pick.ID, 1, time.Date(2030, 1, 5, 0, 0, 0, 0, time.UTC), bulkTestUUID(t))
	if err != nil || registeredPick.Status != "registered" || registeredPick.Lines[0].ReservationVersion != 2 {
		t.Fatalf("registered pick=%+v err=%v", registeredPick, err)
	}
	var shipQuantity, shipReserved string
	if err = pool.QueryRow(ctx, `select sum(quantity)::text,sum(reserved_quantity)::text from inventory.bulk_balance where tenant_id=$1 and organization_id='warehouse' and bin_id='ship'`, tenant).Scan(&shipQuantity, &shipReserved); err != nil || shipQuantity != "7.000000" || shipReserved != "7.000000" {
		t.Fatalf("ship quantity=%s reserved=%s err=%v", shipQuantity, shipReserved, err)
	}
	concurrentIDs := []inventorycontrol.WarehouseIDs{warehouseTestIDs(t, "race-a"), warehouseTestIDs(t, "race-b")}
	concurrentPicks := make([]inventorycontrol.WarehouseActivity, 2)
	concurrentErrors := make([]error, 2)
	var group sync.WaitGroup
	for index := range concurrentIDs {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			concurrentPicks[i], concurrentErrors[i] = repo.CreateWarehousePick(ctx, tenant, concurrentIDs[i], inventorycontrol.WarehousePickCommand{RequestID: fmt.Sprintf("race-request-%d", i), OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "service", DemandID: fmt.Sprintf("service-%d", i), DemandLineID: "line-1", Quantity: "3", UseFEFO: true})
		}(index)
	}
	group.Wait()
	winner, successes := inventorycontrol.WarehouseActivity{}, 0
	for index, concurrentErr := range concurrentErrors {
		if concurrentErr == nil {
			winner, successes = concurrentPicks[index], successes+1
		} else if !errors.Is(concurrentErr, inventorycontrol.ErrConflict) {
			t.Fatalf("unexpected concurrent pick error: %v", concurrentErr)
		}
	}
	if successes != 1 {
		t.Fatalf("concurrent pick successes=%d errors=%v", successes, concurrentErrors)
	}
	cancelled, err := repo.CancelWarehousePick(ctx, tenant, "warehouse", winner.ID, 1, bulkTestUUID(t))
	if err != nil || cancelled.Status != "cancelled" || cancelled.Version != 2 {
		t.Fatalf("cancelled=%+v err=%v", cancelled, err)
	}
	var remainingReserved string
	if err = pool.QueryRow(ctx, `select reserved_quantity::text from inventory.bulk_balance where tenant_id=$1 and bin_id='forward' and lot_id='late-lot'`, tenant).Scan(&remainingReserved); err != nil || remainingReserved != "0.000000" {
		t.Fatalf("cancel released=%s err=%v", remainingReserved, err)
	}
	if _, err = pool.Exec(ctx, `update inventory.warehouse_receipt set source_id='tampered' where tenant_id=$1 and receipt_id='late-receipt'`, tenant); err == nil {
		t.Fatal("immutable receipt accepted update")
	} else {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || pgErr.Code != "55000" {
			t.Fatalf("unexpected immutability error: %v", err)
		}
	}
}
````

### FILE: `internal/platform/httpapi/warehouse.go`

```yaml
block_id: "GO-OPS-API:operational-warehouse-http:v1"
operation: CREATE
provenance: AUTHORED
source: "local authenticated HTTP boundary over the Microsoft-governed warehouse service"
license: "LicenseRef-Workspace-Owner"
sha256: "70af08b4c0857b849f652a5f44fc10b90ef5599d5fe8bc8c0a58b82ec21cfb4f"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"net/http"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
)

type warehouseAPI struct {
	service  *inventorycontrol.WarehouseService
	verifier identity.Verifier
}

func registerWarehouse(mux *http.ServeMux, service *inventorycontrol.WarehouseService, verifier identity.Verifier) {
	api := warehouseAPI{service: service, verifier: verifier}
	mux.HandleFunc("POST /v1/inventory/warehouse/receipts", api.postReceipt)
	mux.HandleFunc("POST /v1/inventory/warehouse/picks", api.createPick)
	mux.HandleFunc("POST /v1/inventory/warehouse/activities/{id}/register", api.register)
	mux.HandleFunc("POST /v1/inventory/warehouse/picks/{id}/cancel", api.cancelPick)
	mux.HandleFunc("POST /v1/inventory/warehouse/put-aways/{id}/cancel", api.cancelPutAway)
	mux.HandleFunc("POST /v1/inventory/warehouse/replenishments", api.createReplenishment)
	mux.HandleFunc("POST /v1/inventory/warehouse/replenishments/{id}/cancel", api.cancelReplenishment)
	mux.HandleFunc("POST /v1/inventory/warehouse/cross-dock-policies", api.configureCrossDock)
	mux.HandleFunc("POST /v1/inventory/warehouse/sales-bindings", api.configureSalesBinding)
	mux.HandleFunc("POST /v1/inventory/warehouse/customer-shipments", api.postCustomerShipment)
}

func (a warehouseAPI) postCustomerShipment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input inventorycontrol.CustomerShipmentCommand
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.PostCustomerShipment(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_CUSTOMER_SHIPMENT") {
		return
	}
	writeJSON(w, 201, value)
}

func (a warehouseAPI) configureSalesBinding(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input inventorycontrol.SalesWarehouseBinding
	if !decodeStrict(w, r, &input) {
		return
	}
	value, err := a.service.ConfigureSalesBinding(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_SALES_WAREHOUSE_BINDING") {
		return
	}
	writeJSON(w, 201, value)
}

func (a warehouseAPI) configureCrossDock(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input inventorycontrol.WarehouseCrossDockPolicy
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.ConfigureCrossDock(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_CROSS_DOCK_POLICY") {
		return
	}
	writeJSON(w, 201, value)
}

func (a warehouseAPI) createReplenishment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input inventorycontrol.WarehouseReplenishmentCommand
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.CreateReplenishment(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_REPLENISHMENT") {
		return
	}
	writeJSON(w, 201, value)
}

func (a warehouseAPI) cancelReplenishment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.CancelReplenishment(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Version)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_REPLENISHMENT_CANCELLATION") {
		return
	}
	writeJSON(w, 200, value)
}

func (a warehouseAPI) authorize(w http.ResponseWriter, r *http.Request) (identity.Principal, bool) {
	return (bulkInventoryAPI{verifier: a.verifier}).authorize(w, r, "inventory:write", true)
}

func (a warehouseAPI) postReceipt(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input inventorycontrol.WarehouseReceiptCommand
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.PostReceipt(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_RECEIPT") {
		return
	}
	writeJSON(w, 201, value)
}

func (a warehouseAPI) createPick(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input inventorycontrol.WarehousePickCommand
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.CreatePick(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_PICK") {
		return
	}
	writeJSON(w, 201, value)
}

func (a warehouseAPI) register(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input struct {
		OrganizationID string    `json:"organization_id"`
		Version        int64     `json:"version"`
		PostingDate    time.Time `json:"posting_date"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.Register(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Version, input.PostingDate)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_REGISTRATION") {
		return
	}
	writeJSON(w, 200, value)
}

func (a warehouseAPI) cancelPick(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.CancelPick(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Version)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_PICK_CANCELLATION") {
		return
	}
	writeJSON(w, 200, value)
}

func (a warehouseAPI) cancelPutAway(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.CancelPutAway(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Version)
	if writeBulkError(w, err, "INVALID_WAREHOUSE_PUT_AWAY_CANCELLATION") {
		return
	}
	writeJSON(w, 200, value)
}
````

### FILE: `internal/platform/httpapi/warehouse_test.go`

```yaml
block_id: "GO-OPS-API:operational-warehouse-http-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local strict JSON, authorization, organization-scope and route regression"
license: "LicenseRef-Workspace-Owner"
sha256: "7916cca7f7f8e4ec0496250fbe318f9160ec01532116a87c2ccbe4804c1c9cae"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
)

type warehouseHTTPRepo struct {
	bindings            int
	shipments           int
	receipts            int
	picks               int
	register            int
	cancel              int
	cancelPutAway       int
	replenish           int
	cancelReplenishment int
	crossDock           int
}

func (r *warehouseHTTPRepo) ConfigureSalesWarehouseBinding(_ context.Context, _ string, _ string, value inventorycontrol.SalesWarehouseBinding) (inventorycontrol.SalesWarehouseBinding, error) {
	r.bindings++
	return value, nil
}
func (r *warehouseHTTPRepo) PostCustomerShipment(_ context.Context, _ string, ids inventorycontrol.CustomerShipmentIDs, value inventorycontrol.CustomerShipmentCommand) (inventorycontrol.CustomerShipment, error) {
	r.shipments++
	return inventorycontrol.CustomerShipment{ID: ids.ShipmentID, RequestID: value.RequestID, OrganizationID: value.OrganizationID, OrderID: value.OrderID, WarehouseActivityID: value.WarehouseActivityID, PostingDate: value.PostingDate, Line: inventorycontrol.CustomerShipmentLine{ID: ids.LineID}}, nil
}

func (r *warehouseHTTPRepo) PostWarehouseReceipt(_ context.Context, _ string, ids inventorycontrol.WarehouseIDs, value inventorycontrol.WarehouseReceiptCommand) (inventorycontrol.WarehouseReceiptResult, error) {
	r.receipts++
	return inventorycontrol.WarehouseReceiptResult{ReceiptID: ids.ReceiptID, Quantity: value.Quantity}, nil
}
func (r *warehouseHTTPRepo) CreateWarehousePick(_ context.Context, _ string, ids inventorycontrol.WarehouseIDs, value inventorycontrol.WarehousePickCommand) (inventorycontrol.WarehouseActivity, error) {
	r.picks++
	return inventorycontrol.WarehouseActivity{ID: ids.ActivityID, OrganizationID: value.OrganizationID, Type: "pick", Status: "open", Version: 1}, nil
}
func (r *warehouseHTTPRepo) CreateWarehouseReplenishment(_ context.Context, _ string, ids inventorycontrol.WarehouseIDs, value inventorycontrol.WarehouseReplenishmentCommand) (inventorycontrol.WarehouseActivity, error) {
	r.replenish++
	return inventorycontrol.WarehouseActivity{ID: ids.ActivityID, OrganizationID: value.OrganizationID, Type: "movement", Status: "open", Version: 1}, nil
}
func (r *warehouseHTTPRepo) ConfigureWarehouseCrossDock(_ context.Context, _ string, _ string, value inventorycontrol.WarehouseCrossDockPolicy) (inventorycontrol.WarehouseCrossDockPolicy, error) {
	r.crossDock++
	return value, nil
}
func (r *warehouseHTTPRepo) RegisterWarehouseActivity(_ context.Context, _, organization, activity string, version int64, _ time.Time, _ string) (inventorycontrol.WarehouseActivity, error) {
	r.register++
	return inventorycontrol.WarehouseActivity{ID: activity, OrganizationID: organization, Status: "registered", Version: version + 1}, nil
}
func (r *warehouseHTTPRepo) CancelWarehousePick(_ context.Context, _, organization, activity string, version int64, _ string) (inventorycontrol.WarehouseActivity, error) {
	r.cancel++
	return inventorycontrol.WarehouseActivity{ID: activity, OrganizationID: organization, Status: "cancelled", Version: version + 1}, nil
}
func (r *warehouseHTTPRepo) CancelWarehousePutAway(_ context.Context, _, organization, activity string, version int64, _ string) (inventorycontrol.WarehouseActivity, error) {
	r.cancelPutAway++
	return inventorycontrol.WarehouseActivity{ID: activity, OrganizationID: organization, Type: "put-away", Status: "cancelled", Version: version + 1}, nil
}
func (r *warehouseHTTPRepo) CancelWarehouseReplenishment(_ context.Context, _, organization, activity string, version int64, _ string) (inventorycontrol.WarehouseActivity, error) {
	r.cancelReplenishment++
	return inventorycontrol.WarehouseActivity{ID: activity, OrganizationID: organization, Type: "movement", Status: "cancelled", Version: version + 1}, nil
}

func TestWarehouseHTTPIsStrictScopedAndExecutable(t *testing.T) {
	repo := &warehouseHTTPRepo{}
	mux := http.NewServeMux()
	InventoryControlModule{Service: inventorycontrol.NewService(&inventoryRepo{}, inventoryIDs{}), Warehouse: inventorycontrol.NewWarehouseService(repo, inventoryIDs{})}.Register(mux, inventoryVerifier{})

	binding := `{"request_id":"binding-1","variant_id":"variant","item_id":"part","sales_uom_code":"BOX"}`
	request := httptest.NewRequest("POST", "/v1/inventory/warehouse/sales-bindings", strings.NewReader(binding))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.bindings != 1 {
		t.Fatalf("binding status=%d calls=%d body=%s", response.Code, repo.bindings, response.Body.String())
	}

	shipment := `{"request_id":"shipment-1","organization_id":"a","order_id":"order-1","warehouse_activity_id":"pick-1","posting_date":"2030-01-03T00:00:00Z"}`
	request = httptest.NewRequest("POST", "/v1/inventory/warehouse/customer-shipments", strings.NewReader(shipment))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.shipments != 1 {
		t.Fatalf("shipment status=%d calls=%d body=%s", response.Code, repo.shipments, response.Body.String())
	}
	forbiddenShipment := strings.Replace(shipment, `"a"`, `"forbidden"`, 1)
	request = httptest.NewRequest("POST", "/v1/inventory/warehouse/customer-shipments", strings.NewReader(forbiddenShipment))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 || repo.shipments != 1 {
		t.Fatalf("forbidden shipment status=%d calls=%d", response.Code, repo.shipments)
	}

	receipt := `{"request_id":"receipt-1","organization_id":"a","receive_bin_id":"receive","item_id":"part","lot_no":"LOT-1","quantity":"2","unit_cost":"10","posting_date":"2030-01-01T00:00:00Z","source_kind":"purchase","source_id":"po-1"}`
	request = httptest.NewRequest("POST", "/v1/inventory/warehouse/receipts", strings.NewReader(receipt))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.receipts != 1 {
		t.Fatalf("receipt status=%d calls=%d body=%s", response.Code, repo.receipts, response.Body.String())
	}

	forbidden := strings.Replace(receipt, `"a"`, `"forbidden"`, 1)
	request = httptest.NewRequest("POST", "/v1/inventory/warehouse/receipts", strings.NewReader(forbidden))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 || repo.receipts != 1 {
		t.Fatalf("forbidden status=%d calls=%d", response.Code, repo.receipts)
	}

	pick := `{"request_id":"pick-1","organization_id":"allowed","ship_bin_id":"ship","item_id":"part","demand_kind":"customer-order","demand_id":"order-1","demand_line_id":"line-1","quantity":"1","use_fefo":true,"unknown":1}`
	request = httptest.NewRequest("POST", "/v1/inventory/warehouse/picks", strings.NewReader(pick))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 400 || repo.picks != 0 {
		t.Fatalf("unknown field status=%d calls=%d", response.Code, repo.picks)
	}

	request = httptest.NewRequest("POST", "/v1/inventory/warehouse/activities/activity-1/register", strings.NewReader(`{"organization_id":"a","version":1,"posting_date":"2030-01-02T00:00:00Z"}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 200 || repo.register != 1 {
		t.Fatalf("register status=%d calls=%d body=%s", response.Code, repo.register, response.Body.String())
	}

	request = httptest.NewRequest("POST", "/v1/inventory/warehouse/put-aways/activity-1/cancel", strings.NewReader(`{"organization_id":"a","version":1}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 200 || repo.cancelPutAway != 1 {
		t.Fatalf("put-away cancellation status=%d calls=%d body=%s", response.Code, repo.cancelPutAway, response.Body.String())
	}

	request = httptest.NewRequest("POST", "/v1/inventory/warehouse/replenishments", strings.NewReader(`{"request_id":"replenish-1","organization_id":"a","to_bin_id":"pick","item_id":"part","use_fefo":true}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.replenish != 1 {
		t.Fatalf("replenishment status=%d calls=%d body=%s", response.Code, repo.replenish, response.Body.String())
	}

	request = httptest.NewRequest("POST", "/v1/inventory/warehouse/replenishments/activity-1/cancel", strings.NewReader(`{"organization_id":"a","version":1}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 200 || repo.cancelReplenishment != 1 {
		t.Fatalf("replenishment cancellation status=%d calls=%d body=%s", response.Code, repo.cancelReplenishment, response.Body.String())
	}

	request = httptest.NewRequest("POST", "/v1/inventory/warehouse/replenishments", strings.NewReader(`{"request_id":"replenish-2","organization_id":"a","to_bin_id":"pick","item_id":"part","unexpected":true}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 400 || repo.replenish != 1 {
		t.Fatalf("strict replenishment status=%d calls=%d body=%s", response.Code, repo.replenish, response.Body.String())
	}

	request = httptest.NewRequest("POST", "/v1/inventory/warehouse/cross-dock-policies", strings.NewReader(`{"organization_id":"a","item_id":"part","bin_id":"cross-dock","due_date_days":3,"enabled":true,"version":1}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.crossDock != 1 {
		t.Fatalf("cross-dock policy status=%d calls=%d body=%s", response.Code, repo.crossDock, response.Body.String())
	}
}
````

### FILE: `db/migrations/0028_bulk_transfer_in_transit.up.sql`

```yaml
block_id: "GO-OPS-API:bulk-transfer-in-transit-up-sql:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d TransferHeader/TransferLine and shipment/receipt posting invariants"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "83d92d563bc5243607c87169c946081835c433b79631dad79d2fafcb5b4ecffd"
variables: []
secrets_allowed: false
```

````sql
begin;

create table inventory.bulk_transfer (
  tenant_id uuid not null,
  transfer_id text not null,
  request_id text not null,
  from_organization_id text not null,
  to_organization_id text not null,
  receive_bin_id text not null,
  in_transit_code text not null,
  status text not null check (status in ('released','shipped','received','cancelled')),
  posting_date date not null,
  source_kind text not null,
  source_id text not null,
  put_away_activity_id text,
  version bigint not null check (version > 0),
  shipped_at timestamptz,
  received_at timestamptz,
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, transfer_id),
  foreign key (tenant_id, from_organization_id) references org.organization (tenant_id, organization_id),
  foreign key (tenant_id, to_organization_id, receive_bin_id) references inventory.warehouse_bin (tenant_id, organization_id, bin_id),
  foreign key (tenant_id, put_away_activity_id) references inventory.warehouse_activity (tenant_id, activity_id),
  unique (tenant_id, request_id),
  check (from_organization_id <> to_organization_id),
  check (length(in_transit_code) between 1 and 50),
  check ((status='released' and shipped_at is null and received_at is null and put_away_activity_id is null) or
         (status='shipped' and shipped_at is not null and received_at is null and put_away_activity_id is null) or
         (status='received' and shipped_at is not null and received_at is not null and put_away_activity_id is not null) or
         (status='cancelled' and received_at is null and put_away_activity_id is null))
);

create table inventory.bulk_transfer_line (
  tenant_id uuid not null,
  transfer_id text not null,
  line_id text not null,
  item_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  shipped_quantity numeric(20,6) not null default 0 check (shipped_quantity >= 0 and shipped_quantity <= quantity),
  received_quantity numeric(20,6) not null default 0 check (received_quantity >= 0 and received_quantity <= shipped_quantity),
  cost_amount numeric(24,4) not null default 0 check (cost_amount >= 0),
  primary key (tenant_id, transfer_id, line_id),
  foreign key (tenant_id, transfer_id) references inventory.bulk_transfer (tenant_id, transfer_id),
  foreign key (tenant_id, item_id) references inventory.stock_item (tenant_id, item_id),
  unique (tenant_id, transfer_id, item_id)
);

create table inventory.bulk_transfer_allocation (
  tenant_id uuid not null,
  transfer_id text not null,
  line_id text not null,
  allocation_id text not null,
  reservation_id text not null,
  from_bin_id text not null,
  item_id text not null,
  lot_id text,
  quantity numeric(20,6) not null check (quantity > 0),
  outbound_entry_id text not null,
  primary key (tenant_id, transfer_id, allocation_id),
  foreign key (tenant_id, transfer_id, line_id) references inventory.bulk_transfer_line (tenant_id, transfer_id, line_id),
  foreign key (tenant_id, reservation_id) references inventory.bulk_reservation (tenant_id, reservation_id),
  foreign key (tenant_id, outbound_entry_id) references inventory.bulk_inventory_entry (tenant_id, entry_id),
  foreign key (tenant_id, lot_id) references inventory.inventory_lot (tenant_id, lot_id),
  unique (tenant_id, reservation_id)
);

create table inventory.bulk_transfer_cost_component (
  tenant_id uuid not null,
  transfer_id text not null,
  allocation_id text not null,
  component_id text not null,
  source_receipt_entry_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  unit_cost numeric(20,4) not null check (unit_cost >= 0),
  cost_amount numeric(24,4) not null check (cost_amount >= 0),
  receipt_entry_id text,
  primary key (tenant_id, transfer_id, component_id),
  foreign key (tenant_id, transfer_id, allocation_id) references inventory.bulk_transfer_allocation (tenant_id, transfer_id, allocation_id),
  foreign key (tenant_id, source_receipt_entry_id) references inventory.bulk_inventory_entry (tenant_id, entry_id),
  foreign key (tenant_id, receipt_entry_id) references inventory.bulk_inventory_entry (tenant_id, entry_id),
  check (round(quantity * unit_cost,4)=cost_amount)
);

create index bulk_transfer_in_transit_idx
  on inventory.bulk_transfer(tenant_id,to_organization_id,status,posting_date)
  where status='shipped';

create view inventory.bulk_inventory_in_transit as
select t.tenant_id,t.transfer_id,t.from_organization_id,t.to_organization_id,t.in_transit_code,
       l.line_id,l.item_id,a.lot_id,sum(c.quantity) quantity,sum(c.cost_amount) cost_amount,
       t.shipped_at,t.version
from inventory.bulk_transfer t
join inventory.bulk_transfer_line l using(tenant_id,transfer_id)
join inventory.bulk_transfer_allocation a using(tenant_id,transfer_id,line_id)
join inventory.bulk_transfer_cost_component c using(tenant_id,transfer_id,allocation_id)
where t.status='shipped'
group by t.tenant_id,t.transfer_id,t.from_organization_id,t.to_organization_id,t.in_transit_code,
         l.line_id,l.item_id,a.lot_id,t.shipped_at,t.version;

comment on view inventory.bulk_inventory_in_transit is
  'BC-derived transfer order evidence: shipped demand remains unavailable at source and visible in transit until destination receipt.';

commit;
````

### FILE: `db/migrations/0028_bulk_transfer_in_transit.down.sql`

```yaml
block_id: "GO-OPS-API:bulk-transfer-in-transit-down-sql:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d TransferHeader/TransferLine and shipment/receipt posting invariants"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "5d86318e64df95a97e781d0a1c9bd4ff4cbfaf3eb434c1c17753f9c964db2265"
variables: []
secrets_allowed: false
```

````sql
begin;
drop view if exists inventory.bulk_inventory_in_transit;
drop table if exists inventory.bulk_transfer_cost_component;
drop table if exists inventory.bulk_transfer_allocation;
drop table if exists inventory.bulk_transfer_line;
drop table if exists inventory.bulk_transfer;
commit;
````

### FILE: `docs/inventory/MICROSOFT_BC_TRANSFER_DERIVATION.md`

```yaml
block_id: "GO-OPS-API:microsoft-bc-transfer-derivation:v1"
operation: CREATE
provenance: AUTHORED
source: "local derivation and provenance boundary"
license: "LicenseRef-Workspace-Owner"
sha256: "98d9c3900d362da5a3cecbdfc6862f037fcf45c82248c1320dd6e9aef9185f11"
variables: []
secrets_allowed: false
```

````markdown
# Microsoft Business Central transfer-order derivation

This implementation is `ADAPTED`, not Microsoft-authored Go. It is governed by the MIT-licensed BCApps transfer header, line, shipment posting and receipt posting sources fixed at commit `2eae56d704a1fd035d104f333602aea7091b7749`, together with Microsoft's current transfer-order and item-tracking documentation fixed in the pack metadata.

The portable contract preserves the upstream invariants: one transfer line represents outbound demand and inbound supply; shipment precedes receipt through an explicit in-transit code; quantities in transit are unavailable at the source; lot identity is unchanged; posted shipment cost components become destination receipt cost layers without value invention; and concurrent posting is serialized with optimistic versions and PostgreSQL locks. `TransferLine.Table.al` validates explicit quantities to ship and receive and rejects receiving beyond quantity in transit. `TransferOrderPostShipment.Codeunit.al` adds quantity to ship to cumulative quantity shipped and creates posted shipment evidence; `TransferOrderPostReceipt.Codeunit.al` does the equivalent for quantity received and posted receipt evidence. Microsoft's transfer-mode matrix expressly permits partial posting when an in-transit location is used. The Microsoft item-tracking contract also requires the tracking identity shipped from one location to be received unchanged at the other.

For an item configured with specific costing, the portable line therefore requires one exact source receipt entry. That immutable selector participates in idempotent replay, shipment locks only the selected cost layer, insufficient selected quantity fails closed, and the destination receives the same cost component. FIFO lines reject the selector and continue consuming ordered open layers. This is a declared Go/PostgreSQL translation of the fixed Microsoft contracts, not copied AL and not a claim of Microsoft support.

Workflow: create released transfer -> create/register an exact `transfer-outbound` warehouse pick -> post that pick as one shipment -> observe only shipped-not-received components in `bulk_inventory_in_transit` -> post an explicit receipt quantity into the destination `receive` bin -> execute the receipt-specific put-away under the existing warehouse workflow -> repeat until cumulative quantities close the line. Every shipment and receipt has its own request identity and durable posting identity. Exact idempotent replay returns that posting even with the caller's stale aggregate version; a changed quantity, pick or date under the same request fails closed. A registered pick can fund only one shipment, and concurrent attempts against the same transfer version/pick admit at most one commit. Receipt slices the oldest remaining shipped cost components without changing lot or unit cost and cannot exceed shipped minus received. A released transfer can be cancelled; a transfer with any shipment cannot be cancelled or edited.

Admission remains fail closed: both organizations must belong to the authenticated tenant and principal, the destination bin must be a configured receive bin, each posting quantity must be positive and cannot exceed its cumulative remainder, a shipment must name one registered exact pick whose reservations equal that posting quantity, FIFO/specific cost applications must exist, the selector must agree with the configured costing method, and every database transition plus outbox event commits atomically. The integration proof covers repeated partial shipment/receipt, request replay, divergent replay rejection, concurrent shipment contention, in-transit quantity/value after each receipt, lot preservation, FIFO and specific cost, destination put-away, cancellation before shipment and conservation. It does not claim break-bulk, cross-docking, replenishment planning, freight allocation or a live Business Central integration.
````

### FILE: `internal/inventorycontrol/bulk_transfer.go`

```yaml
block_id: "GO-OPS-API:bulk-transfer-domain:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d transfer order lifecycle invariants"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "9c94409787af8e313bbf381363d98226fac837ff36651372539e032cf4a034f5"
variables: []
secrets_allowed: false
```

````go
package inventorycontrol

import (
	"context"
	"fmt"
	"time"
)

type BulkTransferCommand struct {
	RequestID            string    `json:"request_id"`
	FromOrganizationID   string    `json:"from_organization_id"`
	ToOrganizationID     string    `json:"to_organization_id"`
	ReceiveBinID         string    `json:"receive_bin_id"`
	InTransitCode        string    `json:"in_transit_code"`
	ItemID               string    `json:"item_id"`
	Quantity             string    `json:"quantity"`
	SpecificReceiptEntry string    `json:"specific_receipt_entry,omitempty"`
	PostingDate          time.Time `json:"posting_date"`
	SourceKind           string    `json:"source_kind"`
	SourceID             string    `json:"source_id"`
}

type BulkTransfer struct {
	ID                   string    `json:"id"`
	LineID               string    `json:"line_id"`
	RequestID            string    `json:"request_id"`
	FromOrganizationID   string    `json:"from_organization_id"`
	ToOrganizationID     string    `json:"to_organization_id"`
	ReceiveBinID         string    `json:"receive_bin_id"`
	InTransitCode        string    `json:"in_transit_code"`
	ItemID               string    `json:"item_id"`
	Quantity             string    `json:"quantity"`
	SpecificReceiptEntry string    `json:"specific_receipt_entry,omitempty"`
	ShippedQuantity      string    `json:"shipped_quantity"`
	ReceivedQuantity     string    `json:"received_quantity"`
	CostAmount           string    `json:"cost_amount"`
	PostingID            string    `json:"posting_id,omitempty"`
	PutAwayID            string    `json:"put_away_id,omitempty"`
	Status               string    `json:"status"`
	Version              int64     `json:"version"`
	PostingDate          time.Time `json:"posting_date"`
}

type BulkTransferPostingCommand struct {
	RequestID           string    `json:"request_id"`
	Quantity            string    `json:"quantity"`
	WarehouseActivityID string    `json:"warehouse_activity_id,omitempty"`
	PostingDate         time.Time `json:"posting_date"`
}

type BulkTransferRepository interface {
	CreateBulkTransfer(context.Context, string, string, string, string, BulkTransferCommand) (BulkTransfer, error)
	ShipBulkTransfer(context.Context, string, string, string, string, int64, string, string, BulkTransferPostingCommand) (BulkTransfer, error)
	ReceiveBulkTransfer(context.Context, string, string, string, string, int64, string, string, string, BulkTransferPostingCommand) (BulkTransfer, error)
	CancelBulkTransfer(context.Context, string, string, string, string, int64, string) (BulkTransfer, error)
}

type BulkTransferService struct {
	repository BulkTransferRepository
	ids        IDGenerator
}

func NewBulkTransferService(repository BulkTransferRepository, ids IDGenerator) *BulkTransferService {
	return &BulkTransferService{repository: repository, ids: ids}
}

func (s *BulkTransferService) Create(ctx context.Context, tenant string, value BulkTransferCommand) (BulkTransfer, error) {
	if tenant == "" || value.RequestID == "" || value.FromOrganizationID == "" || value.ToOrganizationID == "" || value.FromOrganizationID == value.ToOrganizationID || value.ReceiveBinID == "" || value.InTransitCode == "" || value.ItemID == "" || !positiveDecimal(value.Quantity, quantityPattern) || value.PostingDate.IsZero() || value.SourceKind == "" || value.SourceID == "" {
		return BulkTransfer{}, fmt.Errorf("invalid bulk transfer")
	}
	value.PostingDate = value.PostingDate.UTC()
	return s.repository.CreateBulkTransfer(ctx, tenant, s.ids.New(), s.ids.New(), s.ids.New(), value)
}

func (s *BulkTransferService) Ship(ctx context.Context, tenant, transfer, fromOrganization, toOrganization string, version int64, command BulkTransferPostingCommand) (BulkTransfer, error) {
	if tenant == "" || transfer == "" || fromOrganization == "" || toOrganization == "" || fromOrganization == toOrganization || version < 1 || command.RequestID == "" || !positiveDecimal(command.Quantity, quantityPattern) || command.WarehouseActivityID == "" || command.PostingDate.IsZero() {
		return BulkTransfer{}, fmt.Errorf("invalid bulk transfer shipment")
	}
	command.PostingDate = command.PostingDate.UTC()
	return s.repository.ShipBulkTransfer(ctx, tenant, transfer, fromOrganization, toOrganization, version, s.ids.New(), s.ids.New(), command)
}

func (s *BulkTransferService) Receive(ctx context.Context, tenant, transfer, fromOrganization, toOrganization string, version int64, command BulkTransferPostingCommand) (BulkTransfer, error) {
	if tenant == "" || transfer == "" || fromOrganization == "" || toOrganization == "" || fromOrganization == toOrganization || version < 1 || command.RequestID == "" || !positiveDecimal(command.Quantity, quantityPattern) || command.WarehouseActivityID != "" || command.PostingDate.IsZero() {
		return BulkTransfer{}, fmt.Errorf("invalid bulk transfer receipt")
	}
	command.PostingDate = command.PostingDate.UTC()
	return s.repository.ReceiveBulkTransfer(ctx, tenant, transfer, fromOrganization, toOrganization, version, s.ids.New(), s.ids.New(), s.ids.New(), command)
}

func (s *BulkTransferService) Cancel(ctx context.Context, tenant, transfer, fromOrganization, toOrganization string, version int64) (BulkTransfer, error) {
	if tenant == "" || transfer == "" || fromOrganization == "" || toOrganization == "" || fromOrganization == toOrganization || version < 1 {
		return BulkTransfer{}, fmt.Errorf("invalid bulk transfer cancellation")
	}
	return s.repository.CancelBulkTransfer(ctx, tenant, transfer, fromOrganization, toOrganization, version, s.ids.New())
}
````

### FILE: `internal/inventorycontrol/bulk_transfer_test.go`

```yaml
block_id: "GO-OPS-API:bulk-transfer-domain-test:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d transfer validation invariants"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "55993c0ade1bc85d3b1f95be1127b53259dc1f7271ceb37e8a12b458fdd914a4"
variables: []
secrets_allowed: false
```

````go
package inventorycontrol

import (
	"context"
	"testing"
	"time"
)

type transferRepo struct{ calls int }

func (r *transferRepo) CreateBulkTransfer(context.Context, string, string, string, string, BulkTransferCommand) (BulkTransfer, error) {
	r.calls++
	return BulkTransfer{Status: "released"}, nil
}
func (r *transferRepo) ShipBulkTransfer(context.Context, string, string, string, string, int64, string, string, BulkTransferPostingCommand) (BulkTransfer, error) {
	r.calls++
	return BulkTransfer{Status: "shipped"}, nil
}
func (r *transferRepo) ReceiveBulkTransfer(context.Context, string, string, string, string, int64, string, string, string, BulkTransferPostingCommand) (BulkTransfer, error) {
	r.calls++
	return BulkTransfer{Status: "received"}, nil
}
func (r *transferRepo) CancelBulkTransfer(context.Context, string, string, string, string, int64, string) (BulkTransfer, error) {
	r.calls++
	return BulkTransfer{Status: "cancelled"}, nil
}

type transferIDs struct{ n int }

func (g *transferIDs) New() string { g.n++; return "id" }

func TestBulkTransferServiceRejectsAmbiguityBeforeRepository(t *testing.T) {
	repo := &transferRepo{}
	service := NewBulkTransferService(repo, &transferIDs{})
	_, err := service.Create(context.Background(), "tenant", BulkTransferCommand{RequestID: "request", FromOrganizationID: "same", ToOrganizationID: "same", ReceiveBinID: "receive", InTransitCode: "TRANSIT", ItemID: "item", Quantity: "1", PostingDate: time.Now(), SourceKind: "plan", SourceID: "source"})
	if err == nil || repo.calls != 0 {
		t.Fatalf("expected fail-closed validation, err=%v calls=%d", err, repo.calls)
	}
	_, err = service.Create(context.Background(), "tenant", BulkTransferCommand{RequestID: "request", FromOrganizationID: "from", ToOrganizationID: "to", ReceiveBinID: "receive", InTransitCode: "TRANSIT", ItemID: "item", Quantity: "1.250000", PostingDate: time.Now(), SourceKind: "plan", SourceID: "source"})
	if err != nil || repo.calls != 1 {
		t.Fatalf("expected valid transfer, err=%v calls=%d", err, repo.calls)
	}
	_, err = service.Ship(context.Background(), "tenant", "transfer", "from", "to", 1, BulkTransferPostingCommand{RequestID: "ship", Quantity: "1", PostingDate: time.Now()})
	if err == nil || repo.calls != 1 {
		t.Fatalf("shipment without exact pick must fail, err=%v calls=%d", err, repo.calls)
	}
	_, err = service.Ship(context.Background(), "tenant", "transfer", "from", "to", 1, BulkTransferPostingCommand{RequestID: "ship", Quantity: "1", WarehouseActivityID: "pick", PostingDate: time.Now()})
	if err != nil || repo.calls != 2 {
		t.Fatalf("valid shipment command, err=%v calls=%d", err, repo.calls)
	}
	_, err = service.Receive(context.Background(), "tenant", "transfer", "from", "to", 2, BulkTransferPostingCommand{RequestID: "receipt", Quantity: "1", WarehouseActivityID: "unexpected", PostingDate: time.Now()})
	if err == nil || repo.calls != 2 {
		t.Fatalf("receipt with warehouse activity must fail, err=%v calls=%d", err, repo.calls)
	}
}
````

### FILE: `internal/platform/postgres/bulk_transfer.go`

```yaml
block_id: "GO-OPS-API:bulk-transfer-postgres:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d shipment/in-transit/receipt and cost-conservation invariants"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "ae73df066209ec0d5beb4daf8947d1a685a410f31b070dab81661fce1b76011b"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
)

type transferCost struct {
	layerID, inboundEntryID, quantity, unitCost, costAmount string
}

func readBulkTransfer(ctx context.Context, tx pgx.Tx, tenant, transfer string) (inventorycontrol.BulkTransfer, error) {
	var value inventorycontrol.BulkTransfer
	err := tx.QueryRow(ctx, `select t.transfer_id,l.line_id,t.request_id,t.from_organization_id,t.to_organization_id,t.receive_bin_id,t.in_transit_code,l.item_id,l.quantity::text,coalesce(l.specific_receipt_entry_id,''),l.shipped_quantity::text,l.received_quantity::text,l.cost_amount::text,coalesce(t.put_away_activity_id,''),t.status,t.version,t.posting_date from inventory.bulk_transfer t join inventory.bulk_transfer_line l using(tenant_id,transfer_id) where t.tenant_id=$1 and t.transfer_id=$2`, tenant, transfer).Scan(&value.ID, &value.LineID, &value.RequestID, &value.FromOrganizationID, &value.ToOrganizationID, &value.ReceiveBinID, &value.InTransitCode, &value.ItemID, &value.Quantity, &value.SpecificReceiptEntry, &value.ShippedQuantity, &value.ReceivedQuantity, &value.CostAmount, &value.PutAwayID, &value.Status, &value.Version, &value.PostingDate)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, inventorycontrol.ErrConflict
	}
	return value, err
}

func (r *InventoryControl) CreateBulkTransfer(ctx context.Context, tenant, transferID, lineID, eventID string, command inventorycontrol.BulkTransferCommand) (inventorycontrol.BulkTransfer, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	defer tx.Rollback(ctx)
	var binType, costing string
	err = tx.QueryRow(ctx, `select w.bin_type,i.costing_method from inventory.warehouse_bin w join inventory.item_bin_policy p using(tenant_id,organization_id,bin_id) join inventory.stock_item i on i.tenant_id=p.tenant_id and i.item_id=p.item_id where w.tenant_id=$1 and w.organization_id=$2 and w.bin_id=$3 and p.item_id=$4 and not w.movement_blocked for share of w,p,i`, tenant, command.ToOrganizationID, command.ReceiveBinID, command.ItemID).Scan(&binType, &costing)
	if errors.Is(err, pgx.ErrNoRows) || binType != "receive" {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	if (costing == "specific") != (command.SpecificReceiptEntry != "") || (costing != "fifo" && costing != "specific") {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	inserted, err := tx.Exec(ctx, `insert into inventory.bulk_transfer(tenant_id,transfer_id,request_id,from_organization_id,to_organization_id,receive_bin_id,in_transit_code,status,posting_date,source_kind,source_id,version) values($1,$2,$3,$4,$5,$6,$7,'released',$8::date,$9,$10,1) on conflict (tenant_id,request_id) do nothing`, tenant, transferID, command.RequestID, command.FromOrganizationID, command.ToOrganizationID, command.ReceiveBinID, command.InTransitCode, command.PostingDate, command.SourceKind, command.SourceID)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, bulkConflict(err)
	}
	if inserted.RowsAffected() == 0 {
		err = tx.QueryRow(ctx, `select transfer_id from inventory.bulk_transfer where tenant_id=$1 and request_id=$2 and from_organization_id=$3 and to_organization_id=$4 and receive_bin_id=$5 and in_transit_code=$6 and posting_date=$7::date and source_kind=$8 and source_id=$9`, tenant, command.RequestID, command.FromOrganizationID, command.ToOrganizationID, command.ReceiveBinID, command.InTransitCode, command.PostingDate, command.SourceKind, command.SourceID).Scan(&transferID)
		if errors.Is(err, pgx.ErrNoRows) {
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		if err != nil {
			return inventorycontrol.BulkTransfer{}, err
		}
		value, readErr := readBulkTransfer(ctx, tx, tenant, transferID)
		storedQuantity, storedOK := parseRat(value.Quantity)
		requestedQuantity, requestedOK := parseRat(command.Quantity)
		if readErr != nil || !storedOK || !requestedOK || storedQuantity.Cmp(requestedQuantity) != 0 || value.ItemID != command.ItemID || value.SpecificReceiptEntry != command.SpecificReceiptEntry {
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		return value, tx.Commit(ctx)
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_transfer_line(tenant_id,transfer_id,line_id,item_id,quantity,specific_receipt_entry_id) values($1,$2,$3,$4,$5::numeric,nullif($6,''))`, tenant, transferID, lineID, command.ItemID, command.Quantity, command.SpecificReceiptEntry)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-transfer", transferID, "bulk-transfer.released", 1, map[string]any{"from_organization_id": command.FromOrganizationID, "to_organization_id": command.ToOrganizationID, "item_id": command.ItemID, "quantity": command.Quantity, "in_transit_code": command.InTransitCode}); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	value, err := readBulkTransfer(ctx, tx, tenant, transferID)
	if err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func transferCosts(ctx context.Context, tx pgx.Tx, tenant, organization, item, lot, quantity, specificReceiptEntry string) ([]transferCost, *big.Rat, error) {
	wanted, ok := parseRat(quantity)
	if !ok {
		return nil, nil, inventorycontrol.ErrConflict
	}
	query := `select layer_id,receipt_entry_id,remaining_quantity::text,unit_cost::text from inventory.bulk_cost_layer where tenant_id=$1 and organization_id=$2 and item_id=$3 and lot_id is not distinct from nullif($4,'') and remaining_quantity>0`
	args := []any{tenant, organization, item, lot}
	if specificReceiptEntry != "" {
		query += ` and receipt_entry_id=$5`
		args = append(args, specificReceiptEntry)
	}
	query += ` order by posting_date,receipt_entry_id for update`
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	remaining, total := new(big.Rat).Set(wanted), new(big.Rat)
	values := []transferCost{}
	for rows.Next() && remaining.Sign() > 0 {
		var value transferCost
		var availableText string
		if err = rows.Scan(&value.layerID, &value.inboundEntryID, &availableText, &value.unitCost); err != nil {
			return nil, nil, err
		}
		available, validAvailable := parseRat(availableText)
		unit, validUnit := parseRat(value.unitCost)
		if !validAvailable || !validUnit {
			return nil, nil, inventorycontrol.ErrConflict
		}
		take := minRat(available, remaining)
		value.quantity = formatRat(take, 6)
		value.costAmount = formatRat(new(big.Rat).Mul(take, unit), 4)
		rounded, _ := parseRat(value.costAmount)
		total.Add(total, rounded)
		values = append(values, value)
		remaining.Sub(remaining, take)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, err
	}
	if remaining.Sign() != 0 || len(values) == 0 {
		return nil, nil, inventorycontrol.ErrConflict
	}
	return values, total, nil
}

func (r *InventoryControl) ShipBulkTransfer(ctx context.Context, tenant, transferID, fromOrganization, toOrganization string, version int64, shipmentID, eventID string, command inventorycontrol.BulkTransferPostingCommand) (inventorycontrol.BulkTransfer, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	defer tx.Rollback(ctx)
	var replayID, replayActivity, replayQuantity string
	var replayDate time.Time
	err = tx.QueryRow(ctx, `select s.shipment_id,s.warehouse_activity_id,s.quantity::text,s.posting_date from inventory.bulk_transfer_shipment s join inventory.bulk_transfer t using(tenant_id,transfer_id) where s.tenant_id=$1 and s.transfer_id=$2 and s.request_id=$3 and t.from_organization_id=$4 and t.to_organization_id=$5`, tenant, transferID, command.RequestID, fromOrganization, toOrganization).Scan(&replayID, &replayActivity, &replayQuantity, &replayDate)
	if err == nil {
		stored, storedOK := parseRat(replayQuantity)
		requested, requestedOK := parseRat(command.Quantity)
		if !storedOK || !requestedOK || stored.Cmp(requested) != 0 || replayActivity != command.WarehouseActivityID || !replayDate.Equal(command.PostingDate) {
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		value, readErr := readBulkTransfer(ctx, tx, tenant, transferID)
		if readErr != nil {
			return value, readErr
		}
		value.PostingID = replayID
		return value, tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.BulkTransfer{}, err
	}
	var lineID, organization, itemID, totalText, shippedText, receivedText, costing, specificReceiptEntry string
	err = tx.QueryRow(ctx, `select l.line_id,t.from_organization_id,l.item_id,l.quantity::text,l.shipped_quantity::text,l.received_quantity::text,i.costing_method,coalesce(l.specific_receipt_entry_id,'') from inventory.bulk_transfer t join inventory.bulk_transfer_line l using(tenant_id,transfer_id) join inventory.stock_item i on i.tenant_id=t.tenant_id and i.item_id=l.item_id where t.tenant_id=$1 and t.transfer_id=$2 and t.from_organization_id=$3 and t.to_organization_id=$4 and t.status in ('released','partially-shipped','partially-received') and t.version=$5 for update of t,l`, tenant, transferID, fromOrganization, toOrganization, version).Scan(&lineID, &organization, &itemID, &totalText, &shippedText, &receivedText, &costing, &specificReceiptEntry)
	if errors.Is(err, pgx.ErrNoRows) || (costing == "specific") != (specificReceiptEntry != "") || (costing != "fifo" && costing != "specific") {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	if err = warehouseLock(ctx, tx, tenant, organization, itemID); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	total, totalOK := parseRat(totalText)
	shippedBefore, shippedOK := parseRat(shippedText)
	receivedBefore, receivedOK := parseRat(receivedText)
	postingQuantity, postingOK := parseRat(command.Quantity)
	if !totalOK || !shippedOK || !receivedOK || !postingOK || postingQuantity.Cmp(new(big.Rat).Sub(total, shippedBefore)) > 0 {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	rows, err := tx.Query(ctx, `select r.reservation_id,r.bin_id,coalesce(r.lot_id,''),r.quantity::text,r.version from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) join inventory.bulk_reservation r on r.tenant_id=al.tenant_id and r.reservation_id=al.reservation_id join inventory.warehouse_bin w on w.tenant_id=r.tenant_id and w.organization_id=r.organization_id and w.bin_id=r.bin_id where a.tenant_id=$1 and a.organization_id=$2 and a.activity_id=$3 and a.activity_type='pick' and a.status='registered' and a.source_kind='transfer-outbound' and a.source_id=$4 and a.source_line_id=$5 and r.item_id=$6 and r.demand_kind='transfer-outbound' and r.demand_id=$4 and r.demand_line_id=a.source_line_id||':'||a.activity_id||':'||lpad(al.sequence_no::text,6,'0') and r.status='reservation' and (r.expires_at is null or r.expires_at>clock_timestamp()) and w.bin_type='ship' and not w.movement_blocked order by al.sequence_no for update of a,r`, tenant, organization, command.WarehouseActivityID, transferID, lineID, itemID)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	type allocation struct {
		reservation, bin, lot, quantity string
		reservationVersion              int64
	}
	allocations, shipped := []allocation{}, new(big.Rat)
	for rows.Next() {
		var a allocation
		if err = rows.Scan(&a.reservation, &a.bin, &a.lot, &a.quantity, &a.reservationVersion); err != nil {
			rows.Close()
			return inventorycontrol.BulkTransfer{}, err
		}
		q, ok := parseRat(a.quantity)
		if !ok {
			rows.Close()
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		shipped.Add(shipped, q)
		allocations = append(allocations, a)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	if shipped.Cmp(postingQuantity) != 0 || len(allocations) == 0 {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_transfer_shipment(tenant_id,shipment_id,transfer_id,request_id,warehouse_activity_id,quantity,cost_amount,posting_date) values($1,$2,$3,$4,$5,$6::numeric,0,$7::date)`, tenant, shipmentID, transferID, command.RequestID, command.WarehouseActivityID, command.Quantity, command.PostingDate)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, bulkConflict(err)
	}
	totalCost := new(big.Rat)
	for index, a := range allocations {
		outboundID := fmt.Sprintf("%s-entry-%06d", shipmentID, index+1)
		allocationID := fmt.Sprintf("%s-allocation-%06d", shipmentID, index+1)
		costs, allocationCost, costErr := transferCosts(ctx, tx, tenant, organization, itemID, a.lot, a.quantity, specificReceiptEntry)
		if costErr != nil {
			return inventorycontrol.BulkTransfer{}, costErr
		}
		totalCost.Add(totalCost, allocationCost)
		unitCost := formatRat(new(big.Rat).Quo(allocationCost, mustRat(a.quantity)), 4)
		_, err = tx.Exec(ctx, `insert into inventory.bulk_inventory_entry(tenant_id,entry_id,organization_id,bin_id,item_id,lot_id,entry_type,quantity,unit_cost,cost_amount,posting_date,source_kind,source_id) values($1,$2,$3,$4,$5,nullif($6,''),'issue',-$7::numeric,$8::numeric,-$9::numeric,$10::date,'bulk-transfer-shipment',$11)`, tenant, outboundID, organization, a.bin, itemID, a.lot, a.quantity, unitCost, formatRat(allocationCost, 4), command.PostingDate, fmt.Sprintf("%s:%06d", shipmentID, index+1))
		if err != nil {
			return inventorycontrol.BulkTransfer{}, bulkConflict(err)
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_transfer_allocation(tenant_id,transfer_id,line_id,allocation_id,reservation_id,from_bin_id,item_id,lot_id,quantity,outbound_entry_id,shipment_id) values($1,$2,$3,$4,$5,$6,$7,nullif($8,''),$9::numeric,$10,$11)`, tenant, transferID, lineID, allocationID, a.reservation, a.bin, itemID, a.lot, a.quantity, outboundID, shipmentID)
		if err != nil {
			return inventorycontrol.BulkTransfer{}, bulkConflict(err)
		}
		for componentIndex, cost := range costs {
			updated, updateErr := tx.Exec(ctx, `update inventory.bulk_cost_layer set remaining_quantity=remaining_quantity-$3::numeric where tenant_id=$1 and layer_id=$2 and remaining_quantity >= $3::numeric`, tenant, cost.layerID, cost.quantity)
			if updateErr != nil || updated.RowsAffected() != 1 {
				if updateErr != nil {
					return inventorycontrol.BulkTransfer{}, updateErr
				}
				return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
			}
			_, err = tx.Exec(ctx, `insert into inventory.bulk_cost_application(tenant_id,outbound_entry_id,inbound_entry_id,quantity,cost_amount) values($1,$2,$3,$4::numeric,$5::numeric)`, tenant, outboundID, cost.inboundEntryID, cost.quantity, cost.costAmount)
			if err != nil {
				return inventorycontrol.BulkTransfer{}, bulkConflict(err)
			}
			componentID := fmt.Sprintf("%s-cost-%06d-%06d", shipmentID, index+1, componentIndex+1)
			_, err = tx.Exec(ctx, `insert into inventory.bulk_transfer_cost_component(tenant_id,transfer_id,allocation_id,component_id,source_receipt_entry_id,quantity,unit_cost,cost_amount) values($1,$2,$3,$4,$5,$6::numeric,$7::numeric,$8::numeric)`, tenant, transferID, allocationID, componentID, cost.inboundEntryID, cost.quantity, cost.unitCost, cost.costAmount)
			if err != nil {
				return inventorycontrol.BulkTransfer{}, bulkConflict(err)
			}
		}
		if err = consumeRegisteredPickPackaging(ctx, tx, tenant, organization, command.WarehouseActivityID, a.reservation); err != nil {
			return inventorycontrol.BulkTransfer{}, err
		}
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_balance set quantity=quantity-$6::numeric,reserved_quantity=reserved_quantity-$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and quantity >= $6::numeric and reserved_quantity >= $6::numeric`, tenant, organization, a.bin, itemID, a.lot, a.quantity)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return inventorycontrol.BulkTransfer{}, updateErr
			}
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		updated, err = tx.Exec(ctx, `update inventory.bulk_reservation set status='consumed',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and reservation_id=$2 and status='reservation' and version=$3`, tenant, a.reservation, a.reservationVersion)
		if err != nil || updated.RowsAffected() != 1 {
			if err != nil {
				return inventorycontrol.BulkTransfer{}, err
			}
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
	}
	shippedAfter := new(big.Rat).Add(shippedBefore, postingQuantity)
	status := "partially-shipped"
	if receivedBefore.Sign() > 0 {
		status = "partially-received"
	} else if shippedAfter.Cmp(total) == 0 {
		status = "shipped"
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_transfer set status=$4,version=version+1,shipped_at=coalesce(shipped_at,clock_timestamp()),updated_at=clock_timestamp() where tenant_id=$1 and transfer_id=$2 and status in ('released','partially-shipped','partially-received') and version=$3`, tenant, transferID, version, status)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return inventorycontrol.BulkTransfer{}, err
		}
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `update inventory.bulk_transfer_line set shipped_quantity=shipped_quantity+$3::numeric,cost_amount=cost_amount+$4::numeric where tenant_id=$1 and transfer_id=$2`, tenant, transferID, command.Quantity, formatRat(totalCost, 4))
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	_, err = tx.Exec(ctx, `update inventory.bulk_transfer_shipment set cost_amount=$3::numeric where tenant_id=$1 and shipment_id=$2`, tenant, shipmentID, formatRat(totalCost, 4))
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-transfer", transferID, "bulk-transfer.shipped", version+1, map[string]any{"shipment_id": shipmentID, "request_id": command.RequestID, "organization_id": organization, "item_id": itemID, "quantity": command.Quantity, "cost_amount": formatRat(totalCost, 4), "allocations": len(allocations)}); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	value, err := readBulkTransfer(ctx, tx, tenant, transferID)
	if err != nil {
		return value, err
	}
	value.PostingID = shipmentID
	return value, tx.Commit(ctx)
}

func mustRat(value string) *big.Rat {
	parsed, ok := parseRat(value)
	if !ok {
		panic("validated decimal became invalid")
	}
	return parsed
}

func (r *InventoryControl) ReceiveBulkTransfer(ctx context.Context, tenant, transferID, fromOrganization, toOrganization string, version int64, receiptID, putAwayID, eventID string, command inventorycontrol.BulkTransferPostingCommand) (inventorycontrol.BulkTransfer, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	defer tx.Rollback(ctx)
	var replayID, replayPutAway, replayQuantity string
	var replayDate time.Time
	err = tx.QueryRow(ctx, `select r.receipt_id,r.put_away_activity_id,r.quantity::text,r.posting_date from inventory.bulk_transfer_receipt r join inventory.bulk_transfer t using(tenant_id,transfer_id) where r.tenant_id=$1 and r.transfer_id=$2 and r.request_id=$3 and t.from_organization_id=$4 and t.to_organization_id=$5`, tenant, transferID, command.RequestID, fromOrganization, toOrganization).Scan(&replayID, &replayPutAway, &replayQuantity, &replayDate)
	if err == nil {
		stored, storedOK := parseRat(replayQuantity)
		requested, requestedOK := parseRat(command.Quantity)
		if !storedOK || !requestedOK || stored.Cmp(requested) != 0 || !replayDate.Equal(command.PostingDate) {
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		value, readErr := readBulkTransfer(ctx, tx, tenant, transferID)
		if readErr != nil {
			return value, readErr
		}
		value.PostingID, value.PutAwayID = replayID, replayPutAway
		return value, tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.BulkTransfer{}, err
	}
	var lineID, organization, receiveBin, itemID, totalText, shippedText, receivedText string
	err = tx.QueryRow(ctx, `select l.line_id,t.to_organization_id,t.receive_bin_id,l.item_id,l.quantity::text,l.shipped_quantity::text,l.received_quantity::text from inventory.bulk_transfer t join inventory.bulk_transfer_line l using(tenant_id,transfer_id) join inventory.warehouse_bin w on w.tenant_id=t.tenant_id and w.organization_id=t.to_organization_id and w.bin_id=t.receive_bin_id where t.tenant_id=$1 and t.transfer_id=$2 and t.from_organization_id=$3 and t.to_organization_id=$4 and t.status in ('partially-shipped','shipped','partially-received') and t.version=$5 and w.bin_type='receive' and not w.movement_blocked for update of t,l,w`, tenant, transferID, fromOrganization, toOrganization, version).Scan(&lineID, &organization, &receiveBin, &itemID, &totalText, &shippedText, &receivedText)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	if err = warehouseLock(ctx, tx, tenant, organization, itemID); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	total, totalOK := parseRat(totalText)
	shipped, shippedOK := parseRat(shippedText)
	receivedBefore, receivedOK := parseRat(receivedText)
	postingQuantity, postingOK := parseRat(command.Quantity)
	if !totalOK || !shippedOK || !receivedOK || !postingOK || postingQuantity.Cmp(new(big.Rat).Sub(shipped, receivedBefore)) > 0 {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	rows, err := tx.Query(ctx, `select c.component_id,coalesce(a.lot_id,''),(c.quantity-c.received_quantity)::text,c.unit_cost::text from inventory.bulk_transfer_cost_component c join inventory.bulk_transfer_allocation a using(tenant_id,transfer_id,allocation_id) join inventory.bulk_transfer_shipment s on s.tenant_id=a.tenant_id and s.shipment_id=a.shipment_id where c.tenant_id=$1 and c.transfer_id=$2 and c.received_quantity<c.quantity order by s.created_at,c.component_id for update of c`, tenant, transferID)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	type component struct{ id, lot, quantity, unit, cost string }
	components := []component{}
	remaining := new(big.Rat).Set(postingQuantity)
	totalCost := new(big.Rat)
	for rows.Next() && remaining.Sign() > 0 {
		var c component
		var availableText string
		if err = rows.Scan(&c.id, &c.lot, &availableText, &c.unit); err != nil {
			rows.Close()
			return inventorycontrol.BulkTransfer{}, err
		}
		available, availableOK := parseRat(availableText)
		unit, unitOK := parseRat(c.unit)
		if !availableOK || !unitOK {
			rows.Close()
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		take := minRat(available, remaining)
		c.quantity = formatRat(take, 6)
		c.cost = formatRat(new(big.Rat).Mul(take, unit), 4)
		rounded, _ := parseRat(c.cost)
		totalCost.Add(totalCost, rounded)
		components = append(components, c)
		remaining.Sub(remaining, take)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	if len(components) == 0 || remaining.Sign() != 0 {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity(tenant_id,activity_id,organization_id,activity_type,status,source_kind,source_id,source_line_id,request_id,version) values($1,$2,$3,'put-away','open','bulk-transfer-receipt',$4,$5,$6,1)`, tenant, putAwayID, organization, transferID, lineID, command.RequestID)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_transfer_receipt(tenant_id,receipt_id,transfer_id,request_id,quantity,cost_amount,posting_date,put_away_activity_id) values($1,$2,$3,$4,$5::numeric,$6::numeric,$7::date,$8)`, tenant, receiptID, transferID, command.RequestID, command.Quantity, formatRat(totalCost, 4), command.PostingDate, putAwayID)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, bulkConflict(err)
	}
	for index, c := range components {
		receiptEntryID := fmt.Sprintf("%s-entry-%06d", receiptID, index+1)
		layerID := fmt.Sprintf("%s-layer-%06d", receiptID, index+1)
		var balanceVersion int64
		err = tx.QueryRow(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity,version) values($1,$2,$3,$4,nullif($5,''),$6::numeric,0,1) on conflict (tenant_id,organization_id,bin_id,item_id,lot_id) do update set quantity=inventory.bulk_balance.quantity+excluded.quantity,version=inventory.bulk_balance.version+1,updated_at=clock_timestamp() where (select max_quantity is null or inventory.bulk_balance.quantity+excluded.quantity<=max_quantity from inventory.item_bin_policy where tenant_id=$1 and organization_id=$2 and item_id=$4 and bin_id=$3) returning version`, tenant, organization, receiveBin, itemID, c.lot, c.quantity).Scan(&balanceVersion)
		if errors.Is(err, pgx.ErrNoRows) {
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		if err != nil {
			return inventorycontrol.BulkTransfer{}, bulkConflict(err)
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_inventory_entry(tenant_id,entry_id,organization_id,bin_id,item_id,lot_id,entry_type,quantity,unit_cost,cost_amount,posting_date,source_kind,source_id) values($1,$2,$3,$4,$5,nullif($6,''),'receipt',$7::numeric,$8::numeric,$9::numeric,$10::date,'bulk-transfer-receipt',$11)`, tenant, receiptEntryID, organization, receiveBin, itemID, c.lot, c.quantity, c.unit, c.cost, command.PostingDate, receiptID+":"+c.id)
		if err != nil {
			return inventorycontrol.BulkTransfer{}, bulkConflict(err)
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_cost_layer(tenant_id,layer_id,receipt_entry_id,organization_id,item_id,lot_id,posting_date,original_quantity,remaining_quantity,unit_cost) values($1,$2,$3,$4,$5,nullif($6,''),$7::date,$8::numeric,$8::numeric,$9::numeric)`, tenant, layerID, receiptEntryID, organization, itemID, c.lot, command.PostingDate, c.quantity, c.unit)
		if err != nil {
			return inventorycontrol.BulkTransfer{}, bulkConflict(err)
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_transfer_receipt_component(tenant_id,receipt_id,transfer_id,component_id,receipt_entry_id,quantity,cost_amount) values($1,$2,$3,$4,$5,$6::numeric,$7::numeric)`, tenant, receiptID, transferID, c.id, receiptEntryID, c.quantity, c.cost)
		if err != nil {
			return inventorycontrol.BulkTransfer{}, bulkConflict(err)
		}
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_transfer_cost_component set receipt_entry_id=coalesce(receipt_entry_id,$4),received_quantity=received_quantity+$5::numeric where tenant_id=$1 and transfer_id=$2 and component_id=$3 and received_quantity+$5::numeric<=quantity`, tenant, transferID, c.id, receiptEntryID, c.quantity)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return inventorycontrol.BulkTransfer{}, updateErr
			}
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
	}
	type lotTotal struct{ lot, quantity, sourceEntry string }
	lotOrder := []string{}
	lotTotals := map[string]*lotTotal{}
	for index, c := range components {
		current, exists := lotTotals[c.lot]
		if !exists {
			current = &lotTotal{lot: c.lot, quantity: "0", sourceEntry: fmt.Sprintf("%s-entry-%06d", receiptID, index+1)}
			lotTotals[c.lot] = current
			lotOrder = append(lotOrder, c.lot)
		}
		total, _ := parseRat(current.quantity)
		part, _ := parseRat(c.quantity)
		current.quantity = formatRat(total.Add(total, part), 6)
	}
	sequence := 0
	for _, lotID := range lotOrder {
		lot := lotTotals[lotID]
		rows, queryErr := tx.Query(ctx, `select p.bin_id,coalesce(p.max_quantity::text,''),coalesce(b.quantity::text,'0'),coalesce((select sum(al.quantity)::text from inventory.warehouse_activity_line al join inventory.warehouse_activity a using(tenant_id,activity_id) where al.tenant_id=p.tenant_id and a.organization_id=p.organization_id and a.activity_type='put-away' and a.status='open' and al.to_bin_id=p.bin_id and al.item_id=p.item_id),'0') from inventory.item_bin_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) left join inventory.bulk_balance b on b.tenant_id=p.tenant_id and b.organization_id=p.organization_id and b.bin_id=p.bin_id and b.item_id=p.item_id and b.lot_id is not distinct from nullif($4,'') where p.tenant_id=$1 and p.organization_id=$2 and p.item_id=$3 and p.bin_id<>$5 and w.bin_type in ('put-away','putpick') and not w.movement_blocked order by p.is_default desc,w.bin_rank desc,p.fixed desc,w.bin_code for share of p,w`, tenant, organization, itemID, lotID, receiveBin)
		if queryErr != nil {
			return inventorycontrol.BulkTransfer{}, queryErr
		}
		remaining, ok := parseRat(lot.quantity)
		if !ok {
			rows.Close()
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		placements := []warehousePlacement{}
		for rows.Next() && remaining.Sign() > 0 {
			var binID, maximum, onHand, planned string
			if err = rows.Scan(&binID, &maximum, &onHand, &planned); err != nil {
				rows.Close()
				return inventorycontrol.BulkTransfer{}, err
			}
			capacity, valid := decimalCapacity(maximum, onHand, planned)
			if !valid {
				rows.Close()
				return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
			}
			take := new(big.Rat).Set(remaining)
			if capacity != nil {
				take = minRat(take, capacity)
			}
			if take.Sign() > 0 {
				placements = append(placements, warehousePlacement{binID: binID, quantity: formatRat(take, 6)})
				remaining.Sub(remaining, take)
			}
		}
		rows.Close()
		if err = rows.Err(); err != nil {
			return inventorycontrol.BulkTransfer{}, err
		}
		if remaining.Sign() != 0 || len(placements) == 0 {
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		var expiration *time.Time
		if lotID != "" {
			if err = tx.QueryRow(ctx, `select expiration_date from inventory.inventory_lot where tenant_id=$1 and lot_id=$2`, tenant, lotID).Scan(&expiration); err != nil {
				return inventorycontrol.BulkTransfer{}, err
			}
		}
		for _, placement := range placements {
			sequence++
			activityLineID := fmt.Sprintf("%s-%06d", putAwayID, sequence)
			_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity_line(tenant_id,activity_id,organization_id,line_id,sequence_no,from_bin_id,to_bin_id,item_id,lot_id,quantity,expiration_date,source_entry_id,uom_code,uom_quantity,qty_per_uom,from_uom_code,from_uom_quantity,from_qty_per_uom) select $1,$2,$3,$4,$5,$6,$7,$8,nullif($9,''),$10::numeric,$11,$12,u.uom_code,$10::numeric/u.qty_per_uom,u.qty_per_uom,u.uom_code,$10::numeric/u.qty_per_uom,u.qty_per_uom from inventory.item_unit_of_measure u where u.tenant_id=$1 and u.item_id=$8 and u.is_base`, tenant, putAwayID, organization, activityLineID, sequence, receiveBin, placement.binID, itemID, lotID, placement.quantity, expiration, lot.sourceEntry)
			if err != nil {
				return inventorycontrol.BulkTransfer{}, bulkConflict(err)
			}
			var receiveBalanceID int64
			err = tx.QueryRow(ctx, `update inventory.bulk_balance set reserved_quantity=reserved_quantity+$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and quantity-reserved_quantity >= $6::numeric returning balance_id`, tenant, organization, receiveBin, itemID, lotID, placement.quantity).Scan(&receiveBalanceID)
			if errors.Is(err, pgx.ErrNoRows) {
				return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
			}
			if err != nil {
				return inventorycontrol.BulkTransfer{}, err
			}
			reserved, reserveErr := tx.Exec(ctx, `update inventory.bulk_uom_balance composition set reserved_quantity=composition.reserved_quantity+($3::numeric/u.qty_per_uom),version=composition.version+1,updated_at=clock_timestamp() from inventory.item_unit_of_measure u where composition.balance_id=$1 and u.tenant_id=$2 and u.item_id=$4 and u.is_base and composition.uom_code=u.uom_code and composition.quantity-composition.reserved_quantity >= ($3::numeric/u.qty_per_uom)`, receiveBalanceID, tenant, placement.quantity, itemID)
			if reserveErr != nil {
				return inventorycontrol.BulkTransfer{}, bulkConflict(reserveErr)
			}
			if reserved.RowsAffected() != 1 {
				return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
			}
		}
	}
	receivedAfter := new(big.Rat).Add(receivedBefore, postingQuantity)
	status := "partially-received"
	finalReceipt := receivedAfter.Cmp(total) == 0
	if finalReceipt {
		status = "received"
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_transfer set status=$4,version=version+1,received_at=case when $5 then clock_timestamp() else null end,put_away_activity_id=$6,updated_at=clock_timestamp() where tenant_id=$1 and transfer_id=$2 and status in ('partially-shipped','shipped','partially-received') and version=$3`, tenant, transferID, version, status, finalReceipt, putAwayID)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return inventorycontrol.BulkTransfer{}, err
		}
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `update inventory.bulk_transfer_line set received_quantity=received_quantity+$3::numeric where tenant_id=$1 and transfer_id=$2`, tenant, transferID, command.Quantity)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-transfer", transferID, "bulk-transfer.received", version+1, map[string]any{"receipt_id": receiptID, "request_id": command.RequestID, "organization_id": organization, "item_id": itemID, "quantity": command.Quantity, "cost_amount": formatRat(totalCost, 4), "cost_components": len(components), "receive_bin_id": receiveBin, "put_away_id": putAwayID, "put_away_lines": sequence}); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	value, err := readBulkTransfer(ctx, tx, tenant, transferID)
	if err != nil {
		return value, err
	}
	value.PostingID, value.PutAwayID = receiptID, putAwayID
	return value, tx.Commit(ctx)
}

func (r *InventoryControl) CancelBulkTransfer(ctx context.Context, tenant, transferID, fromOrganization, toOrganization string, version int64, eventID string) (inventorycontrol.BulkTransfer, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	defer tx.Rollback(ctx)
	var organization, itemID string
	err = tx.QueryRow(ctx, `select t.from_organization_id,l.item_id from inventory.bulk_transfer t join inventory.bulk_transfer_line l using(tenant_id,transfer_id) where t.tenant_id=$1 and t.transfer_id=$2 and t.from_organization_id=$3 and t.to_organization_id=$4 and t.status='released' and t.version=$5 for update of t,l`, tenant, transferID, fromOrganization, toOrganization, version).Scan(&organization, &itemID)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	if err = warehouseLock(ctx, tx, tenant, organization, itemID); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	rows, err := tx.Query(ctx, `select a.activity_id,al.reservation_id,al.from_bin_id,coalesce(al.lot_id,''),al.quantity::text,r.version from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) join inventory.bulk_reservation r on r.tenant_id=al.tenant_id and r.reservation_id=al.reservation_id where a.tenant_id=$1 and a.organization_id=$2 and a.activity_type='pick' and a.status='open' and a.source_kind='transfer-outbound' and a.source_id=$3 and r.status='reservation' order by a.activity_id,al.sequence_no for update of a,r`, tenant, organization, transferID)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	activities := map[string]bool{}
	type transferRelease struct {
		reservation, bin, lot, quantity string
		reservationVersion              int64
	}
	releases := []transferRelease{}
	for rows.Next() {
		var activity, reservation, bin, lot, quantity string
		var reservationVersion int64
		if err = rows.Scan(&activity, &reservation, &bin, &lot, &quantity, &reservationVersion); err != nil {
			rows.Close()
			return inventorycontrol.BulkTransfer{}, err
		}
		activities[activity] = true
		releases = append(releases, transferRelease{reservation: reservation, bin: bin, lot: lot, quantity: quantity, reservationVersion: reservationVersion})
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	for _, release := range releases {
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_balance set reserved_quantity=reserved_quantity-$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and reserved_quantity >= $6::numeric`, tenant, organization, release.bin, itemID, release.lot, release.quantity)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return inventorycontrol.BulkTransfer{}, updateErr
			}
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		updated, updateErr = tx.Exec(ctx, `update inventory.bulk_reservation set status='released',cancellation_disallowed=false,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and reservation_id=$2 and status='reservation' and version=$3`, tenant, release.reservation, release.reservationVersion)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return inventorycontrol.BulkTransfer{}, updateErr
			}
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
	}
	for activity := range activities {
		updated, updateErr := tx.Exec(ctx, `update inventory.warehouse_activity set status='cancelled',version=version+1 where tenant_id=$1 and activity_id=$2 and status='open'`, tenant, activity)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return inventorycontrol.BulkTransfer{}, updateErr
			}
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_transfer set status='cancelled',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and transfer_id=$2 and status='released' and version=$3`, tenant, transferID, version)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return inventorycontrol.BulkTransfer{}, err
		}
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-transfer", transferID, "bulk-transfer.cancelled", version+1, map[string]any{"organization_id": organization, "cancelled_pick_activities": len(activities)}); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	value, err := readBulkTransfer(ctx, tx, tenant, transferID)
	if err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}
````

### FILE: `internal/platform/postgres/bulk_transfer_integration_test.go`

```yaml
block_id: "GO-OPS-API:bulk-transfer-postgres-integration-test:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d transfer posting and item tracking invariants"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "581255ea8eca55fcd252085e37104548a44939bc8b19b8fa0322217ff75c599f"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestBulkTransferShipTransitReceiveConservesQuantityLotAndCost(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := bulkTestUUID(t)
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Transfer V164','Transfer V164')`, tenant, "tr-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	defer func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			t.Errorf("cleanup begin: %v", cleanupErr)
			return
		}
		defer tx.Rollback(ctx)
		commands := []string{
			`alter table inventory.warehouse_pick_request disable trigger warehouse_pick_request_immutable`,
			`alter table inventory.warehouse_activity_uom_reservation disable trigger warehouse_uom_reservation_delete_guard`,
			`alter table inventory.warehouse_activity_line disable trigger warehouse_activity_line_immutable`,
			`alter table inventory.warehouse_receipt_line disable trigger warehouse_receipt_line_immutable`,
			`alter table inventory.warehouse_receipt disable trigger warehouse_receipt_immutable`,
			`alter table inventory.bulk_cost_application disable trigger bulk_cost_application_immutable`,
			`alter table inventory.bulk_inventory_entry disable trigger bulk_inventory_entry_immutable`,
			`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`,
			`alter table inventory.bulk_uom_conversion disable trigger bulk_uom_conversion_immutable`,
			`alter table inventory.bulk_uom_conversion_line disable trigger bulk_uom_conversion_line_immutable`,
			`delete from inventory.warehouse_activity_uom_reservation where tenant_id=$1`,
			`delete from inventory.bulk_transfer_receipt_component where tenant_id=$1`,
			`delete from inventory.bulk_transfer_cost_component where tenant_id=$1`,
			`delete from inventory.bulk_transfer_allocation where tenant_id=$1`,
			`delete from inventory.bulk_transfer_receipt where tenant_id=$1`,
			`delete from inventory.bulk_transfer_shipment where tenant_id=$1`,
			`delete from inventory.bulk_transfer_line where tenant_id=$1`,
			`delete from inventory.bulk_transfer where tenant_id=$1`,
			`delete from inventory.warehouse_pick_request where tenant_id=$1`,
			`delete from inventory.warehouse_activity_line where tenant_id=$1`,
			`delete from inventory.warehouse_activity where tenant_id=$1`,
			`delete from inventory.warehouse_receipt_line where tenant_id=$1`,
			`delete from inventory.warehouse_receipt where tenant_id=$1`,
			`delete from inventory.bulk_uom_conversion_line where tenant_id=$1`,
			`delete from inventory.bulk_uom_conversion where tenant_id=$1`,
			`delete from inventory.bulk_cost_application where tenant_id=$1`,
			`delete from inventory.bulk_cost_layer where tenant_id=$1`,
			`delete from inventory.bulk_inventory_entry where tenant_id=$1`,
			`delete from inventory.bulk_reservation where tenant_id=$1`,
			`delete from inventory.bulk_balance where tenant_id=$1`,
			`delete from inventory.inventory_lot where tenant_id=$1`,
			`delete from inventory.item_bin_policy where tenant_id=$1`,
			`delete from inventory.warehouse_bin where tenant_id=$1`,
			`delete from inventory.item_unit_of_measure where tenant_id=$1`,
			`delete from inventory.stock_item where tenant_id=$1`,
			`delete from platform.outbox_event where tenant_id=$1`,
			`delete from org.organization where tenant_id=$1`,
			`delete from platform.tenant where tenant_id=$1`,
			`alter table inventory.warehouse_activity_line enable trigger warehouse_activity_line_immutable`,
			`alter table inventory.warehouse_receipt_line enable trigger warehouse_receipt_line_immutable`,
			`alter table inventory.warehouse_receipt enable trigger warehouse_receipt_immutable`,
			`alter table inventory.bulk_cost_application enable trigger bulk_cost_application_immutable`,
			`alter table inventory.bulk_inventory_entry enable trigger bulk_inventory_entry_immutable`,
			`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`,
			`alter table inventory.bulk_uom_conversion_line enable trigger bulk_uom_conversion_line_immutable`,
			`alter table inventory.bulk_uom_conversion enable trigger bulk_uom_conversion_immutable`,
			`alter table inventory.warehouse_activity_uom_reservation enable trigger warehouse_uom_reservation_delete_guard`,
			`alter table inventory.warehouse_pick_request enable trigger warehouse_pick_request_immutable`,
		}
		for _, command := range commands {
			args := []any{}
			if len(command) >= 6 && command[:6] == "delete" {
				args = []any{tenant}
			}
			if _, cleanupErr = tx.Exec(ctx, command, args...); cleanupErr != nil {
				t.Errorf("cleanup %s: %v", command, cleanupErr)
				return
			}
		}
		if cleanupErr = tx.Commit(ctx); cleanupErr != nil {
			t.Errorf("cleanup commit: %v", cleanupErr)
		}
	}()
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'source','source','Source','warehouse'),($1,'destination','destination','Destination','warehouse')`, tenant); err != nil {
		t.Fatal(err)
	}
	repo := NewInventoryControl(pool)
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "part", Code: "TRANSFER-PART", Description: "Transfer part", BaseUOM: "EA", BaseRoundingPrecision: "0.000001", TrackingMode: "lot", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.ConfigureItemUnitOfMeasure(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemUnitOfMeasure{ItemID: "part", Code: "BOX", QuantityPerUnit: "2", RoundingPrecision: "1", Version: 1}); err != nil {
		t.Fatal(err)
	}
	bins := []inventorycontrol.WarehouseBin{{ID: "source-receive", OrganizationID: "source", Code: "RECEIVE", Type: "receive"}, {ID: "source-pick", OrganizationID: "source", Code: "PICK", Type: "putpick", Ranking: 100}, {ID: "source-ship", OrganizationID: "source", Code: "SHIP", Type: "ship"}, {ID: "destination-receive", OrganizationID: "destination", Code: "RECEIVE", Type: "receive"}, {ID: "destination-pick", OrganizationID: "destination", Code: "PICK", Type: "putpick", Ranking: 100}}
	for _, bin := range bins {
		if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatal(err)
		}
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: bin.OrganizationID, ItemID: "part", BinID: bin.ID, Fixed: true, Default: bin.Type == "putpick", MinQuantity: "0", MaxQuantity: "100", Version: 1}); err != nil {
			t.Fatal(err)
		}
	}
	expiry := time.Date(2035, 1, 1, 0, 0, 0, 0, time.UTC)
	receipt, err := repo.PostWarehouseReceipt(ctx, tenant, warehouseTestIDs(t, "transfer-source"), inventorycontrol.WarehouseReceiptCommand{RequestID: "source-receipt", OrganizationID: "source", ReceiveBinID: "source-receive", ItemID: "part", LotNo: "LOT-TRANSFER", ExpirationDate: &expiry, HandlingUOM: "BOX", HandlingQuantity: "3", UnitCost: "10", PostingDate: time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC), SourceKind: "purchase", SourceID: "po-transfer"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "source", receipt.PutAway.ID, 1, time.Date(2030, 1, 2, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	transfer, err := repo.CreateBulkTransfer(ctx, tenant, "transfer-1", "transfer-line-1", bulkTestUUID(t), inventorycontrol.BulkTransferCommand{RequestID: "transfer-request", FromOrganizationID: "source", ToOrganizationID: "destination", ReceiveBinID: "destination-receive", InTransitCode: "OWN-LOG", ItemID: "part", Quantity: "3", PostingDate: time.Date(2030, 1, 3, 0, 0, 0, 0, time.UTC), SourceKind: "replenishment", SourceID: "plan-1"})
	if err != nil || transfer.Status != "released" {
		t.Fatalf("create=%+v err=%v", transfer, err)
	}
	replay, err := repo.CreateBulkTransfer(ctx, tenant, "ignored-transfer", "ignored-line", bulkTestUUID(t), inventorycontrol.BulkTransferCommand{RequestID: "transfer-request", FromOrganizationID: "source", ToOrganizationID: "destination", ReceiveBinID: "destination-receive", InTransitCode: "OWN-LOG", ItemID: "part", Quantity: "3", PostingDate: time.Date(2030, 1, 3, 0, 0, 0, 0, time.UTC), SourceKind: "replenishment", SourceID: "plan-1"})
	if err != nil || replay.ID != transfer.ID {
		t.Fatalf("idempotent replay=%+v err=%v", replay, err)
	}
	firstPick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "transfer-pick-first"), inventorycontrol.WarehousePickCommand{RequestID: "transfer-pick-first-request", OrganizationID: "source", ShipBinID: "source-ship", ItemID: "part", DemandKind: "transfer-outbound", DemandID: transfer.ID, DemandLineID: transfer.LineID, Quantity: "1", AllowBreakbulk: true, UseFEFO: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "source", firstPick.ID, 1, time.Date(2030, 1, 4, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	firstShipCommand := inventorycontrol.BulkTransferPostingCommand{RequestID: "ship-first-request", Quantity: "1", WarehouseActivityID: firstPick.ID, PostingDate: time.Date(2030, 1, 5, 0, 0, 0, 0, time.UTC)}
	firstShipped, err := repo.ShipBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 1, "shipment-first", bulkTestUUID(t), firstShipCommand)
	if err != nil || firstShipped.Status != "partially-shipped" || firstShipped.Version != 2 || firstShipped.ShippedQuantity != "1.000000" || firstShipped.CostAmount != "10.0000" || firstShipped.PostingID != "shipment-first" {
		t.Fatalf("first partial shipment=%+v err=%v", firstShipped, err)
	}
	firstShipReplay, err := repo.ShipBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 1, "ignored-shipment", bulkTestUUID(t), firstShipCommand)
	if err != nil || firstShipReplay.PostingID != firstShipped.PostingID || firstShipReplay.Version != 2 {
		t.Fatalf("shipment replay=%+v err=%v", firstShipReplay, err)
	}
	divergentShip := firstShipCommand
	divergentShip.Quantity = "2"
	if _, err = repo.ShipBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 1, "ignored-divergent", bulkTestUUID(t), divergentShip); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("divergent shipment replay: %v", err)
	}
	firstReceiptCommand := inventorycontrol.BulkTransferPostingCommand{RequestID: "receive-first-request", Quantity: "1", PostingDate: time.Date(2030, 1, 6, 0, 0, 0, 0, time.UTC)}
	firstReceived, err := repo.ReceiveBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 2, "receipt-first", "put-away-first", bulkTestUUID(t), firstReceiptCommand)
	if err != nil || firstReceived.Status != "partially-received" || firstReceived.Version != 3 || firstReceived.ReceivedQuantity != "1.000000" || firstReceived.PostingID != "receipt-first" || firstReceived.PutAwayID != "put-away-first" {
		t.Fatalf("first partial receipt=%+v err=%v", firstReceived, err)
	}
	firstReceiptReplay, err := repo.ReceiveBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 2, "ignored-receipt", "ignored-put-away", bulkTestUUID(t), firstReceiptCommand)
	if err != nil || firstReceiptReplay.PostingID != firstReceived.PostingID || firstReceiptReplay.PutAwayID != firstReceived.PutAwayID {
		t.Fatalf("receipt replay=%+v err=%v", firstReceiptReplay, err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "destination", firstReceived.PutAwayID, 1, time.Date(2030, 1, 7, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	secondPick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "transfer-pick-second"), inventorycontrol.WarehousePickCommand{RequestID: "transfer-pick-second-request", OrganizationID: "source", ShipBinID: "source-ship", ItemID: "part", DemandKind: "transfer-outbound", DemandID: transfer.ID, DemandLineID: transfer.LineID, Quantity: "2", HandlingUOM: "BOX", UseFEFO: true})
	if err != nil || len(secondPick.Lines) != 1 || secondPick.Lines[0].UOMCode != "BOX" || !exactDecimalEqual(secondPick.Lines[0].UOMQuantity, "1") {
		t.Fatalf("second packaged pick=%+v err=%v", secondPick, err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "source", secondPick.ID, 1, time.Date(2030, 1, 8, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	results := make([]inventorycontrol.BulkTransfer, 2)
	failures := make([]error, 2)
	shipmentIDs := []string{"shipment-second-a", "shipment-second-b"}
	eventIDs := []string{bulkTestUUID(t), bulkTestUUID(t)}
	commands := []inventorycontrol.BulkTransferPostingCommand{{RequestID: "ship-second-a", Quantity: "2", WarehouseActivityID: secondPick.ID, PostingDate: time.Date(2030, 1, 9, 0, 0, 0, 0, time.UTC)}, {RequestID: "ship-second-b", Quantity: "2", WarehouseActivityID: secondPick.ID, PostingDate: time.Date(2030, 1, 9, 0, 0, 0, 0, time.UTC)}}
	var group sync.WaitGroup
	for index := range results {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			results[i], failures[i] = repo.ShipBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 3, shipmentIDs[i], eventIDs[i], commands[i])
		}(index)
	}
	group.Wait()
	successes := 0
	for index, failure := range failures {
		if failure == nil {
			successes++
			if results[index].Status != "partially-received" || results[index].Version != 4 || results[index].CostAmount != "30.0000" {
				t.Fatalf("shipment=%+v", results[index])
			}
		} else if !errors.Is(failure, inventorycontrol.ErrConflict) {
			t.Fatalf("shipment error=%v", failure)
		}
	}
	if successes != 1 {
		t.Fatalf("shipment successes=%d errors=%v", successes, failures)
	}
	var winningShipment inventorycontrol.BulkTransfer
	var winningCommand inventorycontrol.BulkTransferPostingCommand
	for index := range results {
		if failures[index] == nil {
			winningShipment, winningCommand = results[index], commands[index]
		}
	}
	winningReplay, err := repo.ShipBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 3, "ignored-winning-replay", bulkTestUUID(t), winningCommand)
	if err != nil || winningReplay.PostingID != winningShipment.PostingID || winningReplay.Version != 4 {
		t.Fatalf("winning shipment replay=%+v err=%v", winningReplay, err)
	}
	var transitQuantity, transitCost, transitLot string
	if err = pool.QueryRow(ctx, `select quantity::text,cost_amount::text,lot_id from inventory.bulk_inventory_in_transit where tenant_id=$1 and transfer_id=$2`, tenant, transfer.ID).Scan(&transitQuantity, &transitCost, &transitLot); err != nil || transitQuantity != "2.000000" || transitCost != "20.0000" || transitLot != "transfer-source-lot" {
		t.Fatalf("transit quantity=%s cost=%s lot=%s err=%v", transitQuantity, transitCost, transitLot, err)
	}
	secondReceiptCommand := inventorycontrol.BulkTransferPostingCommand{RequestID: "receive-second-request", Quantity: "1", PostingDate: time.Date(2030, 1, 10, 0, 0, 0, 0, time.UTC)}
	secondReceived, err := repo.ReceiveBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 4, "receipt-second", "put-away-second", bulkTestUUID(t), secondReceiptCommand)
	if err != nil || secondReceived.Status != "partially-received" || secondReceived.Version != 5 || secondReceived.ReceivedQuantity != "2.000000" {
		t.Fatalf("second partial receipt=%+v err=%v", secondReceived, err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "destination", secondReceived.PutAwayID, 1, time.Date(2030, 1, 11, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select quantity::text,cost_amount::text,lot_id from inventory.bulk_inventory_in_transit where tenant_id=$1 and transfer_id=$2`, tenant, transfer.ID).Scan(&transitQuantity, &transitCost, &transitLot); err != nil || transitQuantity != "1.000000" || transitCost != "10.0000" {
		t.Fatalf("partial transit quantity=%s cost=%s err=%v", transitQuantity, transitCost, err)
	}
	finalReceiptCommand := inventorycontrol.BulkTransferPostingCommand{RequestID: "receive-final-request", Quantity: "1", PostingDate: time.Date(2030, 1, 12, 0, 0, 0, 0, time.UTC)}
	received, err := repo.ReceiveBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 5, "receipt-final", "put-away-final", bulkTestUUID(t), finalReceiptCommand)
	if err != nil || received.Status != "received" || received.Version != 6 || received.ReceivedQuantity != "3.000000" || received.CostAmount != "30.0000" || received.PutAwayID != "put-away-final" {
		t.Fatalf("final receive=%+v err=%v", received, err)
	}
	registeredPutAway, err := repo.RegisterWarehouseActivity(ctx, tenant, "destination", received.PutAwayID, 1, time.Date(2030, 1, 13, 0, 0, 0, 0, time.UTC), bulkTestUUID(t))
	if err != nil || registeredPutAway.Status != "registered" || len(registeredPutAway.Lines) != 1 || registeredPutAway.Lines[0].ToBinID != "destination-pick" {
		t.Fatalf("destination put-away=%+v err=%v", registeredPutAway, err)
	}
	var transitRows, shipmentRows, receiptRows int
	if err = pool.QueryRow(ctx, `select (select count(*) from inventory.bulk_inventory_in_transit where tenant_id=$1 and transfer_id=$2),(select count(*) from inventory.bulk_transfer_shipment where tenant_id=$1 and transfer_id=$2),(select count(*) from inventory.bulk_transfer_receipt where tenant_id=$1 and transfer_id=$2)`, tenant, transfer.ID).Scan(&transitRows, &shipmentRows, &receiptRows); err != nil || transitRows != 0 || shipmentRows != 2 || receiptRows != 3 {
		t.Fatalf("posting evidence transit=%d shipments=%d receipts=%d err=%v", transitRows, shipmentRows, receiptRows, err)
	}
	var sourceQuantity, destinationQuantity, destinationCost string
	if err = pool.QueryRow(ctx, `select (select sum(quantity)::text from inventory.bulk_balance where tenant_id=$1 and organization_id='source'),(select sum(quantity)::text from inventory.bulk_balance where tenant_id=$1 and organization_id='destination'),(select sum(cost_amount)::text from inventory.bulk_inventory_entry where tenant_id=$1 and organization_id='destination' and entry_type='receipt' and source_kind='bulk-transfer-receipt')`, tenant).Scan(&sourceQuantity, &destinationQuantity, &destinationCost); err != nil || sourceQuantity != "3.000000" || destinationQuantity != "3.000000" || destinationCost != "30.0000" {
		t.Fatalf("conservation source=%s destination=%s cost=%s err=%v", sourceQuantity, destinationQuantity, destinationCost, err)
	}
	cancellable, err := repo.CreateBulkTransfer(ctx, tenant, "transfer-cancel", "transfer-cancel-line", bulkTestUUID(t), inventorycontrol.BulkTransferCommand{RequestID: "transfer-cancel-request", FromOrganizationID: "source", ToOrganizationID: "destination", ReceiveBinID: "destination-receive", InTransitCode: "OWN-LOG", ItemID: "part", Quantity: "1", PostingDate: time.Date(2030, 1, 8, 0, 0, 0, 0, time.UTC), SourceKind: "replenishment", SourceID: "plan-cancel"})
	if err != nil {
		t.Fatal(err)
	}
	cancelPick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "transfer-cancel-pick"), inventorycontrol.WarehousePickCommand{RequestID: "transfer-cancel-pick-request", OrganizationID: "source", ShipBinID: "source-ship", ItemID: "part", DemandKind: "transfer-outbound", DemandID: cancellable.ID, DemandLineID: cancellable.LineID, Quantity: "1", UseFEFO: true})
	if err != nil {
		t.Fatal(err)
	}
	cancelled, err := repo.CancelBulkTransfer(ctx, tenant, cancellable.ID, "source", "destination", 1, bulkTestUUID(t))
	if err != nil || cancelled.Status != "cancelled" || cancelled.Version != 2 {
		t.Fatalf("cancelled transfer=%+v err=%v", cancelled, err)
	}
	var cancelledActivity, releasedReservation, reservedAfterCancel string
	if err = pool.QueryRow(ctx, `select a.status,r.status,b.reserved_quantity::text from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) join inventory.bulk_reservation r on r.tenant_id=al.tenant_id and r.reservation_id=al.reservation_id join inventory.bulk_balance b on b.tenant_id=r.tenant_id and b.organization_id=r.organization_id and b.bin_id=r.bin_id and b.item_id=r.item_id and b.lot_id is not distinct from r.lot_id where a.tenant_id=$1 and a.activity_id=$2`, tenant, cancelPick.ID).Scan(&cancelledActivity, &releasedReservation, &reservedAfterCancel); err != nil || cancelledActivity != "cancelled" || releasedReservation != "released" || reservedAfterCancel != "0.000000" {
		t.Fatalf("cancel evidence activity=%s reservation=%s reserved=%s err=%v", cancelledActivity, releasedReservation, reservedAfterCancel, err)
	}
	if _, err = repo.ShipBulkTransfer(ctx, tenant, cancellable.ID, "source", "destination", 1, "cancelled-shipment", bulkTestUUID(t), inventorycontrol.BulkTransferPostingCommand{RequestID: "cancelled-ship-request", Quantity: "1", WarehouseActivityID: cancelPick.ID, PostingDate: time.Date(2030, 1, 14, 0, 0, 0, 0, time.UTC)}); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("cancelled transfer shipped: %v", err)
	}

	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "specific-part", Code: "TRANSFER-SPECIFIC", Description: "Specific-cost transfer part", BaseUOM: "EA", BaseRoundingPrecision: "0.000001", TrackingMode: "none", CostingMethod: "specific", Version: 1}); err != nil {
		t.Fatal(err)
	}
	for _, bin := range bins {
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: bin.OrganizationID, ItemID: "specific-part", BinID: bin.ID, Fixed: true, Default: bin.Type == "putpick", MinQuantity: "0", MaxQuantity: "100", Version: 1}); err != nil {
			t.Fatal(err)
		}
	}
	firstSpecific, err := repo.PostWarehouseReceipt(ctx, tenant, warehouseTestIDs(t, "specific-first"), inventorycontrol.WarehouseReceiptCommand{RequestID: "specific-first-request", OrganizationID: "source", ReceiveBinID: "source-receive", ItemID: "specific-part", Quantity: "1", UnitCost: "11", PostingDate: time.Date(2030, 2, 1, 0, 0, 0, 0, time.UTC), SourceKind: "purchase", SourceID: "specific-po-first"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "source", firstSpecific.PutAway.ID, 1, time.Date(2030, 2, 2, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	secondSpecific, err := repo.PostWarehouseReceipt(ctx, tenant, warehouseTestIDs(t, "specific-second"), inventorycontrol.WarehouseReceiptCommand{RequestID: "specific-second-request", OrganizationID: "source", ReceiveBinID: "source-receive", ItemID: "specific-part", Quantity: "1", UnitCost: "17", PostingDate: time.Date(2030, 2, 3, 0, 0, 0, 0, time.UTC), SourceKind: "purchase", SourceID: "specific-po-second"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "source", secondSpecific.PutAway.ID, 1, time.Date(2030, 2, 4, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	specificCommand := inventorycontrol.BulkTransferCommand{RequestID: "specific-transfer-request", FromOrganizationID: "source", ToOrganizationID: "destination", ReceiveBinID: "destination-receive", InTransitCode: "OWN-LOG", ItemID: "specific-part", Quantity: "1", PostingDate: time.Date(2030, 2, 5, 0, 0, 0, 0, time.UTC), SourceKind: "replenishment", SourceID: "specific-plan"}
	if _, err = repo.CreateBulkTransfer(ctx, tenant, "specific-missing", "specific-missing-line", bulkTestUUID(t), specificCommand); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("specific transfer without receipt selector: %v", err)
	}
	specificCommand.SpecificReceiptEntry = secondSpecific.EntryID
	specificTransfer, err := repo.CreateBulkTransfer(ctx, tenant, "specific-transfer", "specific-transfer-line", bulkTestUUID(t), specificCommand)
	if err != nil || specificTransfer.SpecificReceiptEntry != secondSpecific.EntryID {
		t.Fatalf("specific create=%+v err=%v", specificTransfer, err)
	}
	mismatchedReplay := specificCommand
	mismatchedReplay.SpecificReceiptEntry = firstSpecific.EntryID
	if _, err = repo.CreateBulkTransfer(ctx, tenant, "ignored-specific", "ignored-specific-line", bulkTestUUID(t), mismatchedReplay); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("specific replay changed selector: %v", err)
	}
	specificPick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "specific-transfer-pick"), inventorycontrol.WarehousePickCommand{RequestID: "specific-transfer-pick-request", OrganizationID: "source", ShipBinID: "source-ship", ItemID: "specific-part", DemandKind: "transfer-outbound", DemandID: specificTransfer.ID, DemandLineID: specificTransfer.LineID, Quantity: "1"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "source", specificPick.ID, 1, time.Date(2030, 2, 6, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	specificShipped, err := repo.ShipBulkTransfer(ctx, tenant, specificTransfer.ID, "source", "destination", 1, "specific-shipment", bulkTestUUID(t), inventorycontrol.BulkTransferPostingCommand{RequestID: "specific-ship-request", Quantity: "1", WarehouseActivityID: specificPick.ID, PostingDate: time.Date(2030, 2, 7, 0, 0, 0, 0, time.UTC)})
	if err != nil || specificShipped.Status != "shipped" || specificShipped.CostAmount != "17.0000" {
		t.Fatalf("specific shipment=%+v err=%v", specificShipped, err)
	}
	specificReceived, err := repo.ReceiveBulkTransfer(ctx, tenant, specificTransfer.ID, "source", "destination", 2, "specific-receipt-transfer", "specific-putaway", bulkTestUUID(t), inventorycontrol.BulkTransferPostingCommand{RequestID: "specific-receive-request", Quantity: "1", PostingDate: time.Date(2030, 2, 8, 0, 0, 0, 0, time.UTC)})
	if err != nil || specificReceived.Status != "received" || specificReceived.CostAmount != "17.0000" {
		t.Fatalf("specific receipt=%+v err=%v", specificReceived, err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "destination", specificReceived.PutAwayID, 1, time.Date(2030, 2, 9, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	var firstRemaining, secondRemaining, transferredCost string
	if err = pool.QueryRow(ctx, `select (select remaining_quantity::text from inventory.bulk_cost_layer where tenant_id=$1 and receipt_entry_id=$2),(select remaining_quantity::text from inventory.bulk_cost_layer where tenant_id=$1 and receipt_entry_id=$3),(select sum(e.cost_amount)::text from inventory.bulk_transfer_receipt_component rc join inventory.bulk_inventory_entry e on e.tenant_id=rc.tenant_id and e.entry_id=rc.receipt_entry_id where rc.tenant_id=$1 and rc.receipt_id=$4)`, tenant, firstSpecific.EntryID, secondSpecific.EntryID, specificReceived.PostingID).Scan(&firstRemaining, &secondRemaining, &transferredCost); err != nil || firstRemaining != "1.000000" || secondRemaining != "0.000000" || transferredCost != "17.0000" {
		t.Fatalf("specific cost evidence first=%s second=%s transferred=%s err=%v", firstRemaining, secondRemaining, transferredCost, err)
	}
}
````

### FILE: `internal/platform/httpapi/bulk_transfer.go`

```yaml
block_id: "GO-OPS-API:bulk-transfer-http:v1"
operation: CREATE
provenance: AUTHORED
source: "local strict JSON and dual-organization authorization boundary"
license: "LicenseRef-Workspace-Owner"
sha256: "0d7955bd5748d70771db0c5a8dfddd960d0a9accd213ed3274331ceb98c55b62"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"net/http"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
)

type bulkTransferAPI struct {
	service  *inventorycontrol.BulkTransferService
	verifier identity.Verifier
}

func registerBulkTransfer(mux *http.ServeMux, service *inventorycontrol.BulkTransferService, verifier identity.Verifier) {
	api := bulkTransferAPI{service: service, verifier: verifier}
	mux.HandleFunc("POST /v1/inventory/bulk-transfers", api.create)
	mux.HandleFunc("POST /v1/inventory/bulk-transfers/{id}/ship", api.ship)
	mux.HandleFunc("POST /v1/inventory/bulk-transfers/{id}/receive", api.receive)
	mux.HandleFunc("POST /v1/inventory/bulk-transfers/{id}/cancel", api.cancel)
}

func (a bulkTransferAPI) authorize(w http.ResponseWriter, r *http.Request) (identity.Principal, bool) {
	return (bulkInventoryAPI{verifier: a.verifier}).authorize(w, r, "inventory:write", true)
}
func allowedTransfer(p identity.Principal, from, to string) bool {
	return from != "" && to != "" && from != to && p.AllowedOrganization(from) && p.AllowedOrganization(to)
}

func (a bulkTransferAPI) create(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input inventorycontrol.BulkTransferCommand
	if !decodeStrict(w, r, &input) {
		return
	}
	if !allowedTransfer(p, input.FromOrganizationID, input.ToOrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token must be authorized for both transfer organizations")
		return
	}
	value, err := a.service.Create(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_BULK_TRANSFER") {
		return
	}
	writeJSON(w, 201, value)
}

type bulkTransferAction struct {
	FromOrganizationID  string    `json:"from_organization_id"`
	ToOrganizationID    string    `json:"to_organization_id"`
	Version             int64     `json:"version"`
	RequestID           string    `json:"request_id,omitempty"`
	Quantity            string    `json:"quantity,omitempty"`
	WarehouseActivityID string    `json:"warehouse_activity_id,omitempty"`
	PostingDate         time.Time `json:"posting_date"`
}

func (a bulkTransferAPI) action(w http.ResponseWriter, r *http.Request, kind string) {
	p, ok := a.authorize(w, r)
	if !ok {
		return
	}
	var input bulkTransferAction
	if !decodeStrict(w, r, &input) {
		return
	}
	if !allowedTransfer(p, input.FromOrganizationID, input.ToOrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token must be authorized for both transfer organizations")
		return
	}
	var value inventorycontrol.BulkTransfer
	var err error
	switch kind {
	case "ship":
		value, err = a.service.Ship(r.Context(), p.TenantID, r.PathValue("id"), input.FromOrganizationID, input.ToOrganizationID, input.Version, inventorycontrol.BulkTransferPostingCommand{RequestID: input.RequestID, Quantity: input.Quantity, WarehouseActivityID: input.WarehouseActivityID, PostingDate: input.PostingDate})
	case "receive":
		value, err = a.service.Receive(r.Context(), p.TenantID, r.PathValue("id"), input.FromOrganizationID, input.ToOrganizationID, input.Version, inventorycontrol.BulkTransferPostingCommand{RequestID: input.RequestID, Quantity: input.Quantity, WarehouseActivityID: input.WarehouseActivityID, PostingDate: input.PostingDate})
	default:
		value, err = a.service.Cancel(r.Context(), p.TenantID, r.PathValue("id"), input.FromOrganizationID, input.ToOrganizationID, input.Version)
	}
	if writeBulkError(w, err, "INVALID_BULK_TRANSFER_TRANSITION") {
		return
	}
	writeJSON(w, 200, value)
}
func (a bulkTransferAPI) ship(w http.ResponseWriter, r *http.Request)    { a.action(w, r, "ship") }
func (a bulkTransferAPI) receive(w http.ResponseWriter, r *http.Request) { a.action(w, r, "receive") }
func (a bulkTransferAPI) cancel(w http.ResponseWriter, r *http.Request)  { a.action(w, r, "cancel") }
````

### FILE: `internal/platform/httpapi/bulk_transfer_test.go`

```yaml
block_id: "GO-OPS-API:bulk-transfer-http-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local route, strict JSON and least-privilege regression"
license: "LicenseRef-Workspace-Owner"
sha256: "b845fee8515e49950c2462db1518b88f2eb01fd671b88fdcdd6c9826c80b80bc"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"elite.local/enterprise/internal/inventorycontrol"
)

type transferHTTPRepo struct{ creates, ships, receives, cancels int }

func (r *transferHTTPRepo) CreateBulkTransfer(_ context.Context, _ string, transfer, line, event string, value inventorycontrol.BulkTransferCommand) (inventorycontrol.BulkTransfer, error) {
	r.creates++
	return inventorycontrol.BulkTransfer{ID: transfer, LineID: line, FromOrganizationID: value.FromOrganizationID, ToOrganizationID: value.ToOrganizationID, Status: "released", Version: 1}, nil
}
func (r *transferHTTPRepo) ShipBulkTransfer(context.Context, string, string, string, string, int64, string, string, inventorycontrol.BulkTransferPostingCommand) (inventorycontrol.BulkTransfer, error) {
	r.ships++
	return inventorycontrol.BulkTransfer{Status: "shipped", Version: 2}, nil
}
func (r *transferHTTPRepo) ReceiveBulkTransfer(context.Context, string, string, string, string, int64, string, string, string, inventorycontrol.BulkTransferPostingCommand) (inventorycontrol.BulkTransfer, error) {
	r.receives++
	return inventorycontrol.BulkTransfer{Status: "received", Version: 3}, nil
}
func (r *transferHTTPRepo) CancelBulkTransfer(context.Context, string, string, string, string, int64, string) (inventorycontrol.BulkTransfer, error) {
	r.cancels++
	return inventorycontrol.BulkTransfer{Status: "cancelled", Version: 2}, nil
}

func TestBulkTransferHTTPRequiresBothOrganizationScopesAndStrictJSON(t *testing.T) {
	repo := &transferHTTPRepo{}
	mux := http.NewServeMux()
	InventoryControlModule{Service: inventorycontrol.NewService(&inventoryRepo{}, inventoryIDs{}), Transfer: inventorycontrol.NewBulkTransferService(repo, inventoryIDs{})}.Register(mux, inventoryVerifier{})
	create := `{"request_id":"request-1","from_organization_id":"a","to_organization_id":"b","receive_bin_id":"receive","in_transit_code":"OWN-LOG","item_id":"part","quantity":"2","posting_date":"2030-01-01T00:00:00Z","source_kind":"replenishment","source_id":"plan-1"}`
	request := httptest.NewRequest("POST", "/v1/inventory/bulk-transfers", strings.NewReader(create))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.creates != 1 {
		t.Fatalf("create status=%d calls=%d body=%s", response.Code, repo.creates, response.Body.String())
	}
	forbidden := strings.Replace(create, `"to_organization_id":"b"`, `"to_organization_id":"forbidden"`, 1)
	request = httptest.NewRequest("POST", "/v1/inventory/bulk-transfers", strings.NewReader(forbidden))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 || repo.creates != 1 {
		t.Fatalf("scope status=%d calls=%d", response.Code, repo.creates)
	}
	action := `{"from_organization_id":"a","to_organization_id":"b","version":1,"posting_date":"2030-01-02T00:00:00Z","unknown":true}`
	request = httptest.NewRequest("POST", "/v1/inventory/bulk-transfers/transfer-1/ship", strings.NewReader(action))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 400 || repo.ships != 0 {
		t.Fatalf("strict status=%d calls=%d", response.Code, repo.ships)
	}
	action = `{"from_organization_id":"a","to_organization_id":"b","version":1,"request_id":"ship-1","quantity":"1","warehouse_activity_id":"pick-1","posting_date":"2030-01-02T00:00:00Z"}`
	request = httptest.NewRequest("POST", "/v1/inventory/bulk-transfers/transfer-1/ship", strings.NewReader(action))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 200 || repo.ships != 1 {
		t.Fatalf("ship status=%d calls=%d body=%s", response.Code, repo.ships, response.Body.String())
	}
}
````

### FILE: `db/migrations/0029_bulk_transfer_specific_cost.up.sql`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:0029-bulk-transfer-specific-cost-up:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749 transfer and item-application contracts"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "0539c74ccf92c9fc6022e1b4516a2248c929b08ad34f480ec3a692663585d1cf"
variables: []
secrets_allowed: false
```

````sql
alter table inventory.bulk_transfer_line
  add column specific_receipt_entry_id text;

alter table inventory.bulk_transfer_line
  add constraint bulk_transfer_line_specific_receipt_fk
  foreign key (tenant_id, specific_receipt_entry_id)
  references inventory.bulk_inventory_entry (tenant_id, entry_id);
````

### FILE: `db/migrations/0029_bulk_transfer_specific_cost.down.sql`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:0029-bulk-transfer-specific-cost-down:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749 transfer and item-application contracts"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "6dd581ea3b38e8edaf592d95524975f822a9eb52d86ed60ccae2a83010c913ae"
variables: []
secrets_allowed: false
```

````sql
alter table inventory.bulk_transfer_line
  drop constraint if exists bulk_transfer_line_specific_receipt_fk;

alter table inventory.bulk_transfer_line
  drop column if exists specific_receipt_entry_id;
````

### FILE: `db/migrations/0030_bulk_transfer_partial_posting.up.sql`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:0030-bulk-transfer-partial-posting-up:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749 explicit Qty. to Ship/Receive, posted shipment/receipt and in-transit partial-posting contracts"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "384f39117ba2020babf5d273ba02e7c94c37740d105d20405ee1ab868b843c5d"
variables: []
secrets_allowed: false
```

````sql
begin;

do $migration$
declare
  constraint_name text;
begin
  for constraint_name in
    select conname
    from pg_constraint
    where conrelid='inventory.bulk_transfer'::regclass
      and contype='c'
      and pg_get_constraintdef(oid) ilike '%status%'
  loop
    execute format('alter table inventory.bulk_transfer drop constraint %I', constraint_name);
  end loop;
end;
$migration$;

alter table inventory.bulk_transfer
  add constraint bulk_transfer_status_v164_check
  check (status in ('released','partially-shipped','shipped','partially-received','received','cancelled')),
  add constraint bulk_transfer_state_v164_check
  check (
    (status='released' and shipped_at is null and received_at is null and put_away_activity_id is null) or
    (status='partially-shipped' and shipped_at is not null and received_at is null and put_away_activity_id is null) or
    (status='shipped' and shipped_at is not null and received_at is null and put_away_activity_id is null) or
    (status='partially-received' and shipped_at is not null and received_at is null and put_away_activity_id is not null) or
    (status='received' and shipped_at is not null and received_at is not null and put_away_activity_id is not null) or
    (status='cancelled' and received_at is null and put_away_activity_id is null)
  );

create table inventory.bulk_transfer_shipment (
  tenant_id uuid not null,
  shipment_id text not null,
  transfer_id text not null,
  request_id text not null,
  warehouse_activity_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  cost_amount numeric(24,4) not null check (cost_amount >= 0),
  posting_date date not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, shipment_id),
  foreign key (tenant_id, transfer_id) references inventory.bulk_transfer (tenant_id, transfer_id),
  foreign key (tenant_id, warehouse_activity_id) references inventory.warehouse_activity (tenant_id, activity_id),
  unique (tenant_id, transfer_id, request_id),
  unique (tenant_id, warehouse_activity_id)
);

create table inventory.bulk_transfer_receipt (
  tenant_id uuid not null,
  receipt_id text not null,
  transfer_id text not null,
  request_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  cost_amount numeric(24,4) not null check (cost_amount >= 0),
  posting_date date not null,
  put_away_activity_id text not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, receipt_id),
  foreign key (tenant_id, transfer_id) references inventory.bulk_transfer (tenant_id, transfer_id),
  foreign key (tenant_id, put_away_activity_id) references inventory.warehouse_activity (tenant_id, activity_id),
  unique (tenant_id, transfer_id, request_id),
  unique (tenant_id, put_away_activity_id)
);

alter table inventory.bulk_transfer_allocation
  add column shipment_id text,
  add constraint bulk_transfer_allocation_shipment_fk
    foreign key (tenant_id, shipment_id)
    references inventory.bulk_transfer_shipment (tenant_id, shipment_id);

alter table inventory.bulk_transfer_cost_component
  add column received_quantity numeric(20,6) not null default 0,
  add constraint bulk_transfer_cost_component_received_check
    check (received_quantity >= 0 and received_quantity <= quantity);

update inventory.bulk_transfer_cost_component
set received_quantity=quantity
where receipt_entry_id is not null;

create table inventory.bulk_transfer_receipt_component (
  tenant_id uuid not null,
  receipt_id text not null,
  transfer_id text not null,
  component_id text not null,
  receipt_entry_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  cost_amount numeric(24,4) not null check (cost_amount >= 0),
  primary key (tenant_id, receipt_id, component_id),
  foreign key (tenant_id, receipt_id) references inventory.bulk_transfer_receipt (tenant_id, receipt_id),
  foreign key (tenant_id, transfer_id, component_id) references inventory.bulk_transfer_cost_component (tenant_id, transfer_id, component_id),
  foreign key (tenant_id, receipt_entry_id) references inventory.bulk_inventory_entry (tenant_id, entry_id)
);

drop view inventory.bulk_inventory_in_transit;

create view inventory.bulk_inventory_in_transit as
select t.tenant_id,t.transfer_id,t.from_organization_id,t.to_organization_id,t.in_transit_code,
       l.line_id,l.item_id,a.lot_id,
       sum(c.quantity-c.received_quantity) quantity,
       sum(round((c.quantity-c.received_quantity)*c.unit_cost,4)) cost_amount,
       t.shipped_at,t.version
from inventory.bulk_transfer t
join inventory.bulk_transfer_line l using(tenant_id,transfer_id)
join inventory.bulk_transfer_allocation a using(tenant_id,transfer_id,line_id)
join inventory.bulk_transfer_cost_component c using(tenant_id,transfer_id,allocation_id)
where c.received_quantity<c.quantity
group by t.tenant_id,t.transfer_id,t.from_organization_id,t.to_organization_id,t.in_transit_code,
         l.line_id,l.item_id,a.lot_id,t.shipped_at,t.version;

comment on view inventory.bulk_inventory_in_transit is
  'BC-derived transfer evidence: only shipped cost components not yet received remain visible in transit.';

drop index inventory.bulk_transfer_in_transit_idx;
create index bulk_transfer_in_transit_idx
  on inventory.bulk_transfer(tenant_id,to_organization_id,status,posting_date)
  where status in ('partially-shipped','shipped','partially-received');

commit;
````

### FILE: `db/migrations/0030_bulk_transfer_partial_posting.down.sql`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:0030-bulk-transfer-partial-posting-down:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749 explicit Qty. to Ship/Receive, posted shipment/receipt and in-transit partial-posting contracts"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "12d92f5b9f0e9870c6239763f63d4785e794aa7841db1b63fb37cc4be46088c3"
variables: []
secrets_allowed: false
```

````sql
begin;

do $migration$
begin
  if exists (select 1 from inventory.bulk_transfer_shipment)
     or exists (select 1 from inventory.bulk_transfer_receipt) then
    raise exception using errcode='55000', message='cannot remove partial transfer posting schema after postings exist';
  end if;
end;
$migration$;

drop index inventory.bulk_transfer_in_transit_idx;
drop view inventory.bulk_inventory_in_transit;
drop table inventory.bulk_transfer_receipt_component;

alter table inventory.bulk_transfer_cost_component
  drop constraint bulk_transfer_cost_component_received_check,
  drop column received_quantity;

alter table inventory.bulk_transfer_allocation
  drop constraint bulk_transfer_allocation_shipment_fk,
  drop column shipment_id;

drop table inventory.bulk_transfer_receipt;
drop table inventory.bulk_transfer_shipment;

alter table inventory.bulk_transfer
  drop constraint bulk_transfer_state_v164_check,
  drop constraint bulk_transfer_status_v164_check,
  add constraint bulk_transfer_status_check
    check (status in ('released','shipped','received','cancelled')),
  add constraint bulk_transfer_state_check
    check (
      (status='released' and shipped_at is null and received_at is null and put_away_activity_id is null) or
      (status='shipped' and shipped_at is not null and received_at is null and put_away_activity_id is null) or
      (status='received' and shipped_at is not null and received_at is not null and put_away_activity_id is not null) or
      (status='cancelled' and received_at is null and put_away_activity_id is null)
    );

create index bulk_transfer_in_transit_idx
  on inventory.bulk_transfer(tenant_id,to_organization_id,status,posting_date)
  where status='shipped';

create view inventory.bulk_inventory_in_transit as
select t.tenant_id,t.transfer_id,t.from_organization_id,t.to_organization_id,t.in_transit_code,
       l.line_id,l.item_id,a.lot_id,sum(c.quantity) quantity,sum(c.cost_amount) cost_amount,
       t.shipped_at,t.version
from inventory.bulk_transfer t
join inventory.bulk_transfer_line l using(tenant_id,transfer_id)
join inventory.bulk_transfer_allocation a using(tenant_id,transfer_id,line_id)
join inventory.bulk_transfer_cost_component c using(tenant_id,transfer_id,allocation_id)
where t.status='shipped'
group by t.tenant_id,t.transfer_id,t.from_organization_id,t.to_organization_id,t.in_transit_code,
         l.line_id,l.item_id,a.lot_id,t.shipped_at,t.version;

comment on view inventory.bulk_inventory_in_transit is
  'BC-derived transfer order evidence: shipped demand remains unavailable at source and visible in transit until destination receipt.';

commit;
````

### FILE: `db/migrations/0031_warehouse_bin_replenishment.up.sql`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:0031-warehouse-bin-replenishment-up:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "ca23309dbc8395b183188b9f543015faf7a45cf8dbfdca09855111fdc351cf15"
variables: []
secrets_allowed: false
```

````sql
begin;

alter table inventory.warehouse_activity
  drop constraint warehouse_activity_activity_type_check,
  add constraint warehouse_activity_activity_type_v165_check
    check (activity_type in ('put-away','pick','movement'));

alter table inventory.bulk_reservation
  drop constraint bulk_reservation_demand_kind_check,
  add constraint bulk_reservation_demand_kind_v165_check
    check (demand_kind in ('customer-order','service','transfer-outbound','manual','warehouse-replenishment'));

create table inventory.warehouse_replenishment_request (
  tenant_id uuid not null,
  request_id text not null,
  activity_id text not null,
  organization_id text not null,
  to_bin_id text not null,
  item_id text not null,
  use_fefo boolean not null,
  planned_quantity numeric(20,6) not null check (planned_quantity > 0),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, request_id),
  foreign key (tenant_id, activity_id)
    references inventory.warehouse_activity (tenant_id, activity_id),
  foreign key (tenant_id, organization_id)
    references org.organization (tenant_id, organization_id),
  foreign key (tenant_id, organization_id, to_bin_id)
    references inventory.warehouse_bin (tenant_id, organization_id, bin_id),
  foreign key (tenant_id, item_id)
    references inventory.stock_item (tenant_id, item_id),
  unique (tenant_id, activity_id),
  check (length(request_id) between 1 and 128)
);

commit;
````

### FILE: `db/migrations/0031_warehouse_bin_replenishment.down.sql`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:0031-warehouse-bin-replenishment-down:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "df89384cc5a483d8878955ece953aacea3f319662b5a33e7e029212317e32fd4"
variables: []
secrets_allowed: false
```

````sql
begin;

do $migration$
begin
  if exists (select 1 from inventory.warehouse_replenishment_request)
     or exists (select 1 from inventory.warehouse_activity where activity_type='movement') then
    raise exception using errcode='55000', message='cannot remove replenishment schema after movement evidence exists';
  end if;
end;
$migration$;

drop table inventory.warehouse_replenishment_request;

alter table inventory.bulk_reservation
  drop constraint bulk_reservation_demand_kind_v165_check,
  add constraint bulk_reservation_demand_kind_check
    check (demand_kind in ('customer-order','service','transfer-outbound','manual'));

alter table inventory.warehouse_activity
  drop constraint warehouse_activity_activity_type_v165_check,
  add constraint warehouse_activity_activity_type_check
    check (activity_type in ('put-away','pick'));

commit;
````

### FILE: `docs/inventory/MICROSOFT_BC_REPLENISHMENT_DERIVATION.md`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:microsoft-bc-replenishment-derivation:v1"
operation: CREATE
provenance: AUTHORED
source: "local derivation record for microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749"
license: "LicenseRef-Workspace-Owner"
sha256: "568dc0c1a584dbf4368650b28d8b5223640a6294a30716c53db26dc11c47c861"
variables: []
secrets_allowed: false
```

````markdown
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
````

### FILE: `internal/inventorycontrol/warehouse_replenishment.go`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:warehouse-replenishment-domain:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "e6f8d7e18c1011cbddc3f9176c3297db3113e929cf1381f4b2962e563a9a28c3"
variables: []
secrets_allowed: false
```

````go
package inventorycontrol

import (
	"context"
	"fmt"
)

type WarehouseReplenishmentCommand struct {
	RequestID      string `json:"request_id"`
	OrganizationID string `json:"organization_id"`
	ToBinID        string `json:"to_bin_id"`
	ItemID         string `json:"item_id"`
	TargetUOM      string `json:"target_uom,omitempty"`
	AllowBreakbulk bool   `json:"allow_breakbulk"`
	UseFEFO        bool   `json:"use_fefo"`
	AssignedTo     string `json:"assigned_to,omitempty"`
}

func (s *WarehouseService) CreateReplenishment(ctx context.Context, tenant string, value WarehouseReplenishmentCommand) (WarehouseActivity, error) {
	if tenant == "" || value.RequestID == "" || value.OrganizationID == "" || value.ToBinID == "" || value.ItemID == "" {
		return WarehouseActivity{}, fmt.Errorf("invalid warehouse replenishment")
	}
	return s.repository.CreateWarehouseReplenishment(ctx, tenant, s.newIDs(), value)
}

func (s *WarehouseService) CancelReplenishment(ctx context.Context, tenant, organization, activity string, version int64) (WarehouseActivity, error) {
	if tenant == "" || organization == "" || activity == "" || version < 1 {
		return WarehouseActivity{}, fmt.Errorf("invalid warehouse replenishment cancellation")
	}
	return s.repository.CancelWarehouseReplenishment(ctx, tenant, organization, activity, version, s.ids.New())
}
````

### FILE: `internal/platform/postgres/warehouse_replenishment.go`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:warehouse-replenishment-postgres:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "00ad2f176a8644fd408902fd9fd725d4391775b9a9702d01bff9e5aa1d1210ae"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
)

func (r *InventoryControl) CreateWarehouseReplenishment(ctx context.Context, tenant string, ids inventorycontrol.WarehouseIDs, command inventorycontrol.WarehouseReplenishmentCommand) (inventorycontrol.WarehouseActivity, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	defer tx.Rollback(ctx)

	var replayActivity, replayOrganization, replayBin, replayItem, replayAssigned, replayTargetUOM, replayBaseUOM string
	var replayFEFO, replayBreakbulk bool
	err = tx.QueryRow(ctx, `select r.activity_id,r.organization_id,r.to_bin_id,r.item_id,r.target_uom,r.allow_breakbulk,r.use_fefo,a.assigned_to,(select u.uom_code from inventory.item_unit_of_measure u where u.tenant_id=r.tenant_id and u.item_id=r.item_id and u.is_base) from inventory.warehouse_replenishment_request r join inventory.warehouse_activity a using(tenant_id,activity_id) where r.tenant_id=$1 and r.request_id=$2`, tenant, command.RequestID).Scan(&replayActivity, &replayOrganization, &replayBin, &replayItem, &replayTargetUOM, &replayBreakbulk, &replayFEFO, &replayAssigned, &replayBaseUOM)
	if err == nil {
		requestedTarget := command.TargetUOM
		if requestedTarget == "" {
			requestedTarget = replayBaseUOM
		}
		if replayOrganization != command.OrganizationID || replayBin != command.ToBinID || replayItem != command.ItemID || replayTargetUOM != requestedTarget || replayBreakbulk != command.AllowBreakbulk || replayFEFO != command.UseFEFO || replayAssigned != command.AssignedTo {
			return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
		}
		value, readErr := readWarehouseActivity(ctx, tx, tenant, command.OrganizationID, replayActivity)
		if readErr != nil {
			return value, readErr
		}
		return value, tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.WarehouseActivity{}, err
	}
	if err = warehouseLock(ctx, tx, tenant, command.OrganizationID, command.ItemID); err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	targetUOM, err := resolveWarehouseUOM(ctx, tx, tenant, command.ItemID, command.TargetUOM)
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}

	var minimumText, maximumText, currentText, pendingText, targetType string
	var targetRank int
	var targetFixed, targetBlocked bool
	err = tx.QueryRow(ctx, `select p.min_quantity::text,coalesce(p.max_quantity::text,''),w.bin_rank,w.bin_type,p.fixed,w.movement_blocked,coalesce((select sum(b.quantity)::text from inventory.bulk_balance b where b.tenant_id=p.tenant_id and b.organization_id=p.organization_id and b.bin_id=p.bin_id and b.item_id=p.item_id),'0'),coalesce((select sum(al.quantity)::text from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=p.tenant_id and a.organization_id=p.organization_id and a.activity_type='movement' and a.status='open' and al.to_bin_id=p.bin_id and al.item_id=p.item_id),'0') from inventory.item_bin_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) where p.tenant_id=$1 and p.organization_id=$2 and p.bin_id=$3 and p.item_id=$4 for share of p,w`, tenant, command.OrganizationID, command.ToBinID, command.ItemID).Scan(&minimumText, &maximumText, &targetRank, &targetType, &targetFixed, &targetBlocked, &currentText, &pendingText)
	if errors.Is(err, pgx.ErrNoRows) || !targetFixed || targetBlocked || (targetType != "pick" && targetType != "putpick") || maximumText == "" {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	minimum, minimumOK := parseRat(minimumText)
	maximum, maximumOK := parseRat(maximumText)
	current, currentOK := parseRat(currentText)
	pending, pendingOK := parseRat(pendingText)
	if !minimumOK || !maximumOK || !currentOK || !pendingOK {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}
	effective := new(big.Rat).Add(current, pending)
	if effective.Cmp(minimum) >= 0 {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}
	remaining := new(big.Rat).Sub(maximum, effective)
	if remaining.Sign() <= 0 {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}

	order := `w.bin_rank desc,w.bin_code,coalesce(l.expiration_date,'infinity'::date),coalesce(b.lot_id,'')`
	if command.UseFEFO {
		order = `coalesce(l.expiration_date,'infinity'::date),w.bin_rank desc,w.bin_code,coalesce(b.lot_id,'')`
	}
	query := `select b.balance_id,b.bin_id,coalesce(b.lot_id,''),(b.quantity-b.reserved_quantity)::text,l.expiration_date from inventory.bulk_balance b join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) join inventory.item_bin_policy p on p.tenant_id=b.tenant_id and p.organization_id=b.organization_id and p.bin_id=b.bin_id and p.item_id=b.item_id left join inventory.inventory_lot l on l.tenant_id=b.tenant_id and l.lot_id=b.lot_id where b.tenant_id=$1 and b.organization_id=$2 and b.item_id=$3 and b.bin_id<>$4 and w.bin_rank<$5 and w.bin_type not in ('receive','ship') and not w.movement_blocked and coalesce(l.blocked,false)=false and (l.expiration_date is null or l.expiration_date>=current_date) and b.quantity>b.reserved_quantity order by ` + order + ` for update of b`
	rows, err := tx.Query(ctx, query, tenant, command.OrganizationID, command.ItemID, command.ToBinID, targetRank)
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	type candidate struct {
		balanceID           int64
		bin, lot, available string
		expiration          *time.Time
		plan                warehouseUOMPlan
	}
	rawCandidates := []candidate{}
	for rows.Next() {
		var candidate candidate
		if err = rows.Scan(&candidate.balanceID, &candidate.bin, &candidate.lot, &candidate.available, &candidate.expiration); err != nil {
			rows.Close()
			return inventorycontrol.WarehouseActivity{}, err
		}
		rawCandidates = append(rawCandidates, candidate)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	candidates := []candidate{}
	for _, candidate := range rawCandidates {
		if remaining.Sign() <= 0 {
			break
		}
		available, ok := parseRat(candidate.available)
		if !ok || available.Sign() <= 0 {
			return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
		}
		limit := minRat(available, remaining)
		plans, leftover, planErr := planWarehousePackaging(ctx, tx, candidate.balanceID, targetUOM, command.AllowBreakbulk, limit)
		if planErr != nil {
			return inventorycontrol.WarehouseActivity{}, planErr
		}
		moved := new(big.Rat).Sub(limit, leftover)
		for _, plan := range plans {
			plannedCandidate := candidate
			plannedCandidate.plan = plan
			candidates = append(candidates, plannedCandidate)
		}
		remaining.Sub(remaining, moved)
	}
	if len(candidates) == 0 {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}
	planned := new(big.Rat).Sub(new(big.Rat).Sub(maximum, effective), remaining)

	_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity(tenant_id,activity_id,organization_id,activity_type,status,source_kind,source_id,source_line_id,request_id,assigned_to,version) values($1,$2,$3,'movement','open','bin-replenishment',$4,$5,$6,$7,1)`, tenant, ids.ActivityID, command.OrganizationID, command.ToBinID, command.ItemID, command.RequestID, command.AssignedTo)
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_replenishment_request(tenant_id,request_id,activity_id,organization_id,to_bin_id,item_id,target_uom,allow_breakbulk,use_fefo,planned_quantity) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10::numeric)`, tenant, command.RequestID, ids.ActivityID, command.OrganizationID, command.ToBinID, command.ItemID, targetUOM.code, command.AllowBreakbulk, command.UseFEFO, formatRat(planned, 6))
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, bulkConflict(err)
	}
	for index, candidate := range candidates {
		lineID := fmt.Sprintf("%s-%06d", ids.ActivityID, index+1)
		reservationID := fmt.Sprintf("%s-r-%06d", ids.ActivityID, index+1)
		_, err = tx.Exec(ctx, `insert into inventory.bulk_reservation(tenant_id,reservation_id,organization_id,bin_id,item_id,lot_id,demand_kind,demand_id,demand_line_id,quantity,status,cancellation_disallowed,version) values($1,$2,$3,$4,$5,nullif($6,''),'warehouse-replenishment',$7,$8,$9::numeric,'reservation',true,1)`, tenant, reservationID, command.OrganizationID, candidate.bin, command.ItemID, candidate.lot, ids.ActivityID, lineID, candidate.plan.movedBase)
		if err != nil {
			return inventorycontrol.WarehouseActivity{}, bulkConflict(err)
		}
		_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity_line(tenant_id,activity_id,organization_id,line_id,sequence_no,from_bin_id,to_bin_id,item_id,lot_id,reservation_id,quantity,expiration_date,uom_code,uom_quantity,qty_per_uom,from_uom_code,from_uom_quantity,from_qty_per_uom) values($1,$2,$3,$4,$5,$6,$7,$8,nullif($9,''),$10,$11::numeric,$12,$13,$14::numeric,$15::numeric,$16,$17::numeric,$18::numeric)`, tenant, ids.ActivityID, command.OrganizationID, lineID, index+1, candidate.bin, command.ToBinID, command.ItemID, candidate.lot, reservationID, candidate.plan.movedBase, candidate.expiration, candidate.plan.targetUOM, candidate.plan.targetQuantity, candidate.plan.targetFactor, candidate.plan.sourceUOM, candidate.plan.sourceQuantity, candidate.plan.sourceFactor)
		if err != nil {
			return inventorycontrol.WarehouseActivity{}, bulkConflict(err)
		}
		line := inventorycontrol.WarehouseActivityLine{ID: lineID, Sequence: index + 1, FromBinID: candidate.bin, ToBinID: command.ToBinID, ItemID: command.ItemID, LotID: candidate.lot, ReservationID: reservationID, ReservationVersion: 1, Quantity: candidate.plan.movedBase, FromUOMCode: candidate.plan.sourceUOM, FromUOMQuantity: candidate.plan.sourceQuantity, FromQuantityPerUOM: candidate.plan.sourceFactor, UOMCode: candidate.plan.targetUOM, UOMQuantity: candidate.plan.targetQuantity, QuantityPerUOM: candidate.plan.targetFactor, ExpirationDate: candidate.expiration}
		if err = reserveWarehousePackaging(ctx, tx, tenant, ids.ActivityID, command.OrganizationID, candidate.plan, line); err != nil {
			return inventorycontrol.WarehouseActivity{}, err
		}
	}
	if err = recordBulkEvent(ctx, tx, tenant, ids.EventID, "warehouse-replenishment", ids.ActivityID, "warehouse-replenishment.created", 1, map[string]any{"organization_id": command.OrganizationID, "item_id": command.ItemID, "to_bin_id": command.ToBinID, "quantity": formatRat(planned, 6), "target_uom": targetUOM.code, "allow_breakbulk": command.AllowBreakbulk, "use_fefo": command.UseFEFO, "lines": len(candidates)}); err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	value, err := readWarehouseActivity(ctx, tx, tenant, command.OrganizationID, ids.ActivityID)
	if err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func (r *InventoryControl) CancelWarehouseReplenishment(ctx context.Context, tenant, organization, activityID string, version int64, eventID string) (inventorycontrol.WarehouseActivity, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	defer tx.Rollback(ctx)
	var itemID string
	err = tx.QueryRow(ctx, `select al.item_id from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=$1 and a.organization_id=$2 and a.activity_id=$3 and a.activity_type='movement' and a.source_kind='bin-replenishment' and a.status='open' and a.version=$4 order by al.sequence_no limit 1 for update of a`, tenant, organization, activityID, version).Scan(&itemID)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	if err = warehouseLock(ctx, tx, tenant, organization, itemID); err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	activity, err := readWarehouseActivity(ctx, tx, tenant, organization, activityID)
	if err != nil {
		return activity, err
	}
	for _, line := range activity.Lines {
		if err = releaseWarehousePackagingReservation(ctx, tx, tenant, organization, activityID, line); err != nil {
			return activity, err
		}
		updated, err := tx.Exec(ctx, `update inventory.bulk_reservation set status='released',cancellation_disallowed=false,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and reservation_id=$3 and status='reservation' and version=1`, tenant, organization, line.ReservationID)
		if err != nil || updated.RowsAffected() != 1 {
			if err != nil {
				return activity, err
			}
			return activity, inventorycontrol.ErrConflict
		}
	}
	updated, err := tx.Exec(ctx, `update inventory.warehouse_activity set status='cancelled',version=version+1 where tenant_id=$1 and organization_id=$2 and activity_id=$3 and activity_type='movement' and status='open' and version=$4`, tenant, organization, activityID, version)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return activity, err
		}
		return activity, inventorycontrol.ErrConflict
	}
	activity.Status, activity.Version = "cancelled", version+1
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "warehouse-replenishment", activityID, "warehouse-replenishment.cancelled", activity.Version, map[string]any{"organization_id": organization, "item_id": itemID, "lines": len(activity.Lines)}); err != nil {
		return activity, err
	}
	return activity, tx.Commit(ctx)
}
````

### FILE: `internal/platform/postgres/warehouse_replenishment_integration_test.go`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:warehouse-replenishment-postgres-integration:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "be0117fbafdc68f28baa2d0c428d7035093c913cf907d35c52a2e2011f1456ab"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestWarehouseReplenishmentFEFORegisterAndCancel(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := bulkTestUUID(t)
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Replenishment V165','Replenishment V165')`, tenant, "rp-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'warehouse','warehouse','Warehouse','warehouse')`, tenant); err != nil {
		t.Fatal(err)
	}
	defer cleanupWarehouseReplenishment(t, pool, tenant)

	repo := NewInventoryControl(pool)
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "part", Code: "PART-REPLENISH", Description: "Replenished part", BaseUOM: "EA", BaseRoundingPrecision: "0.000001", TrackingMode: "lot", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	bins := []inventorycontrol.WarehouseBin{
		{ID: "source-early", OrganizationID: "warehouse", Code: "SOURCE-EARLY", Type: "putpick", Ranking: 30, Version: 1},
		{ID: "source-late", OrganizationID: "warehouse", Code: "SOURCE-LATE", Type: "putpick", Ranking: 40, Version: 1},
		{ID: "source-blocked", OrganizationID: "warehouse", Code: "SOURCE-BLOCKED", Type: "putpick", Ranking: 50, MovementBlocked: true, Version: 1},
		{ID: "source-receive", OrganizationID: "warehouse", Code: "SOURCE-RECEIVE", Type: "receive", Ranking: 60, Version: 1},
		{ID: "target", OrganizationID: "warehouse", Code: "TARGET", Type: "pick", Ranking: 100, Version: 1},
		{ID: "target-cancel", OrganizationID: "warehouse", Code: "TARGET-CANCEL", Type: "pick", Ranking: 110, Version: 1},
		{ID: "target-race", OrganizationID: "warehouse", Code: "TARGET-RACE", Type: "pick", Ranking: 120, Version: 1},
	}
	for _, bin := range bins {
		if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatalf("create bin %s: %v", bin.ID, err)
		}
		minimum, maximum := "0", "100"
		if bin.ID == "target" || bin.ID == "target-cancel" || bin.ID == "target-race" {
			minimum, maximum = "5", "10"
		}
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: "warehouse", ItemID: "part", BinID: bin.ID, Fixed: true, MinQuantity: minimum, MaxQuantity: maximum, Version: 1}); err != nil {
			t.Fatalf("configure bin %s: %v", bin.ID, err)
		}
	}

	if _, err = pool.Exec(ctx, `
insert into inventory.inventory_lot(tenant_id,lot_id,item_id,lot_no,expiration_date,blocked) values
 ($1,'lot-target','part','LOT-TARGET','2036-01-01',false),
 ($1,'lot-early','part','LOT-EARLY','2034-01-01',false),
 ($1,'lot-late','part','LOT-LATE','2035-01-01',false),
		 ($1,'lot-blocked','part','LOT-BLOCKED','2033-01-01',true),
		 ($1,'lot-receive','part','LOT-RECEIVE','2032-01-01',false)`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity) values
		 ($1,'warehouse','target','part','lot-target',2,0),
		 ($1,'warehouse','source-early','part','lot-early',5,0),
		 ($1,'warehouse','source-late','part','lot-late',6,0),
 ($1,'warehouse','source-blocked','part','lot-blocked',20,0),
 ($1,'warehouse','source-receive','part','lot-receive',20,0)`, tenant); err != nil {
		t.Fatal(err)
	}

	ids := warehouseTestIDs(t, "replenish")
	command := inventorycontrol.WarehouseReplenishmentCommand{RequestID: "replenish-request", OrganizationID: "warehouse", ToBinID: "target", ItemID: "part", UseFEFO: true, AssignedTo: "worker-1"}
	activity, err := repo.CreateWarehouseReplenishment(ctx, tenant, ids, command)
	if err != nil {
		t.Fatal(err)
	}
	if activity.Type != "movement" || activity.Status != "open" || len(activity.Lines) != 2 || activity.Lines[0].LotNo != "LOT-EARLY" || activity.Lines[0].Quantity != "5.000000" || activity.Lines[1].LotNo != "LOT-LATE" || activity.Lines[1].Quantity != "3.000000" {
		t.Fatalf("unexpected FEFO activity: %+v", activity)
	}
	replay, err := repo.CreateWarehouseReplenishment(ctx, tenant, warehouseTestIDs(t, "replay"), command)
	if err != nil || replay.ID != activity.ID || len(replay.Lines) != 2 {
		t.Fatalf("idempotent replay=%+v err=%v", replay, err)
	}
	divergent := command
	divergent.ToBinID = "target-cancel"
	if _, err = repo.CreateWarehouseReplenishment(ctx, tenant, warehouseTestIDs(t, "divergent"), divergent); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("divergent replay error=%v", err)
	}

	registered, err := repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", activity.ID, 1, time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC), bulkTestUUID(t))
	if err != nil || registered.Status != "registered" || registered.Version != 2 {
		t.Fatalf("register=%+v err=%v", registered, err)
	}
	var targetQuantity, sourceEarly, sourceLate string
	if err = pool.QueryRow(ctx, `select (select sum(quantity)::text from inventory.bulk_balance where tenant_id=$1 and bin_id='target'),(select quantity::text from inventory.bulk_balance where tenant_id=$1 and bin_id='source-early'),(select quantity::text from inventory.bulk_balance where tenant_id=$1 and bin_id='source-late')`, tenant).Scan(&targetQuantity, &sourceEarly, &sourceLate); err != nil || targetQuantity != "10.000000" || sourceEarly != "0.000000" || sourceLate != "3.000000" {
		t.Fatalf("balances target=%s early=%s late=%s err=%v", targetQuantity, sourceEarly, sourceLate, err)
	}
	var consumed int
	if err = pool.QueryRow(ctx, `select count(*) from inventory.bulk_reservation where tenant_id=$1 and demand_kind='warehouse-replenishment' and status='consumed'`, tenant).Scan(&consumed); err != nil || consumed != 2 {
		t.Fatalf("consumed reservations=%d err=%v", consumed, err)
	}
	var movementEntries int
	if err = pool.QueryRow(ctx, `select count(*) from inventory.bulk_inventory_entry where tenant_id=$1 and source_kind='warehouse-movement'`, tenant).Scan(&movementEntries); err != nil || movementEntries != 4 {
		t.Fatalf("movement entries=%d err=%v", movementEntries, err)
	}
	if _, err = repo.CreateWarehouseReplenishment(ctx, tenant, warehouseTestIDs(t, "full"), inventorycontrol.WarehouseReplenishmentCommand{RequestID: "full-request", OrganizationID: "warehouse", ToBinID: "target", ItemID: "part", UseFEFO: true}); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("full target accepted: %v", err)
	}

	cancelActivity, err := repo.CreateWarehouseReplenishment(ctx, tenant, warehouseTestIDs(t, "cancel"), inventorycontrol.WarehouseReplenishmentCommand{RequestID: "cancel-request", OrganizationID: "warehouse", ToBinID: "target-cancel", ItemID: "part", UseFEFO: true})
	if err != nil || len(cancelActivity.Lines) != 3 || cancelActivity.Lines[0].Quantity != "5.000000" || cancelActivity.Lines[1].Quantity != "3.000000" || cancelActivity.Lines[2].Quantity != "2.000000" {
		t.Fatalf("cancel activity=%+v err=%v", cancelActivity, err)
	}
	cancelled, err := repo.CancelWarehouseReplenishment(ctx, tenant, "warehouse", cancelActivity.ID, 1, bulkTestUUID(t))
	if err != nil || cancelled.Status != "cancelled" || cancelled.Version != 2 {
		t.Fatalf("cancelled=%+v err=%v", cancelled, err)
	}
	var released int
	var reserved string
	if err = pool.QueryRow(ctx, `select (select count(*) from inventory.bulk_reservation where tenant_id=$1 and demand_id=$2 and status='released'),(select coalesce(sum(reserved_quantity),0)::text from inventory.bulk_balance where tenant_id=$1)`, tenant, cancelActivity.ID).Scan(&released, &reserved); err != nil || released != 3 || reserved != "0.000000" {
		t.Fatalf("cancel released=%d reserved=%s err=%v", released, reserved, err)
	}

	activities := make([]inventorycontrol.WarehouseActivity, 2)
	errorsByCall := make([]error, 2)
	raceIDs := []inventorycontrol.WarehouseIDs{warehouseTestIDs(t, "race-0"), warehouseTestIDs(t, "race-1")}
	var group sync.WaitGroup
	for index := range activities {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			activities[i], errorsByCall[i] = repo.CreateWarehouseReplenishment(ctx, tenant, raceIDs[i], inventorycontrol.WarehouseReplenishmentCommand{RequestID: fmt.Sprintf("race-request-%d", i), OrganizationID: "warehouse", ToBinID: "target-race", ItemID: "part", UseFEFO: true})
		}(index)
	}
	group.Wait()
	winner, successes := inventorycontrol.WarehouseActivity{}, 0
	for index, callErr := range errorsByCall {
		if callErr == nil {
			winner, successes = activities[index], successes+1
		} else if !errors.Is(callErr, inventorycontrol.ErrConflict) {
			t.Fatalf("unexpected concurrent error: %v", callErr)
		}
	}
	if successes != 1 || len(winner.Lines) == 0 {
		t.Fatalf("concurrent successes=%d errors=%v winner=%+v", successes, errorsByCall, winner)
	}
	if _, err = repo.CancelWarehouseReplenishment(ctx, tenant, "warehouse", winner.ID, 1, bulkTestUUID(t)); err != nil {
		t.Fatalf("cancel concurrent winner: %v", err)
	}
}

func cleanupWarehouseReplenishment(t *testing.T, pool *pgxpool.Pool, tenant string) {
	t.Helper()
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Errorf("replenishment cleanup begin: %v", err)
		return
	}
	defer tx.Rollback(ctx)
	commands := []struct {
		sql  string
		args []any
	}{
		{`alter table inventory.warehouse_pick_request disable trigger warehouse_pick_request_immutable`, nil},
		{`alter table inventory.warehouse_activity_uom_reservation disable trigger warehouse_uom_reservation_delete_guard`, nil},
		{`alter table inventory.warehouse_activity_line disable trigger warehouse_activity_line_immutable`, nil},
		{`alter table inventory.bulk_inventory_entry disable trigger bulk_inventory_entry_immutable`, nil},
		{`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`, nil},
		{`delete from inventory.warehouse_replenishment_request where tenant_id=$1`, []any{tenant}},
		{`delete from inventory.warehouse_activity_uom_reservation where tenant_id=$1`, []any{tenant}},
		{`delete from inventory.warehouse_pick_request where tenant_id=$1`, []any{tenant}},
		{`delete from inventory.warehouse_activity_line where tenant_id=$1`, []any{tenant}},
		{`delete from inventory.warehouse_activity where tenant_id=$1`, []any{tenant}},
		{`delete from inventory.bulk_inventory_entry where tenant_id=$1`, []any{tenant}},
		{`delete from inventory.bulk_reservation where tenant_id=$1`, []any{tenant}},
		{`delete from inventory.bulk_balance where tenant_id=$1`, []any{tenant}},
		{`delete from inventory.inventory_lot where tenant_id=$1`, []any{tenant}},
		{`delete from inventory.item_bin_policy where tenant_id=$1`, []any{tenant}},
		{`delete from inventory.warehouse_bin where tenant_id=$1`, []any{tenant}},
		{`delete from inventory.item_unit_of_measure where tenant_id=$1`, []any{tenant}},
		{`delete from inventory.stock_item where tenant_id=$1`, []any{tenant}},
		{`delete from platform.outbox_event where tenant_id=$1`, []any{tenant}},
		{`delete from org.organization where tenant_id=$1`, []any{tenant}},
		{`delete from platform.tenant where tenant_id=$1`, []any{tenant}},
		{`alter table inventory.warehouse_activity_line enable trigger warehouse_activity_line_immutable`, nil},
		{`alter table inventory.bulk_inventory_entry enable trigger bulk_inventory_entry_immutable`, nil},
		{`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`, nil},
		{`alter table inventory.warehouse_activity_uom_reservation enable trigger warehouse_uom_reservation_delete_guard`, nil},
		{`alter table inventory.warehouse_pick_request enable trigger warehouse_pick_request_immutable`, nil},
	}
	for _, command := range commands {
		if _, err = tx.Exec(ctx, command.sql, command.args...); err != nil {
			t.Errorf("replenishment cleanup failed: %v", err)
			return
		}
	}
	if err = tx.Commit(ctx); err != nil {
		t.Errorf("replenishment cleanup commit: %v", err)
	}
}
````

### FILE: `db/migrations/0032_warehouse_transfer_crossdock.up.sql`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:warehouse-crossdock-migration-up:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "2ad4922d66423a9e7d03f514b49f6cfbba6021e6e446022242cfdd8cb7028a83"
variables: []
secrets_allowed: false
```

````sql
begin;

alter table inventory.warehouse_bin
  add column cross_dock boolean not null default false,
  add constraint warehouse_bin_cross_dock_type_check
    check (not cross_dock or bin_type in ('pick','putpick'));

create table inventory.warehouse_crossdock_policy (
  tenant_id uuid not null,
  organization_id text not null,
  item_id text not null,
  bin_id text not null,
  due_date_days integer not null check (due_date_days between 0 and 365),
  enabled boolean not null,
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, organization_id, item_id),
  foreign key (tenant_id, organization_id, bin_id)
    references inventory.warehouse_bin (tenant_id, organization_id, bin_id),
  foreign key (tenant_id, item_id)
    references inventory.stock_item (tenant_id, item_id)
);

create table inventory.warehouse_crossdock_allocation (
  tenant_id uuid not null,
  allocation_id text not null,
  receipt_id text not null,
  receipt_line_id text not null,
  put_away_activity_id text not null,
  put_away_line_id text not null,
  organization_id text not null,
  bin_id text not null,
  item_id text not null,
  lot_id text,
  transfer_id text not null,
  transfer_line_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  reserved_quantity numeric(20,6) not null default 0,
  picked_quantity numeric(20,6) not null default 0,
  due_date date not null,
  available boolean not null default false,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, allocation_id),
  foreign key (tenant_id, receipt_id, receipt_line_id)
    references inventory.warehouse_receipt_line (tenant_id, receipt_id, line_id),
  foreign key (tenant_id, put_away_activity_id, put_away_line_id)
    references inventory.warehouse_activity_line (tenant_id, activity_id, line_id),
  foreign key (tenant_id, organization_id, bin_id)
    references inventory.warehouse_bin (tenant_id, organization_id, bin_id),
  foreign key (tenant_id, item_id)
    references inventory.stock_item (tenant_id, item_id),
  foreign key (tenant_id, lot_id)
    references inventory.inventory_lot (tenant_id, lot_id),
  foreign key (tenant_id, transfer_id, transfer_line_id)
    references inventory.bulk_transfer_line (tenant_id, transfer_id, line_id),
  unique (tenant_id, put_away_activity_id, put_away_line_id),
  check (reserved_quantity >= 0 and picked_quantity >= 0 and reserved_quantity + picked_quantity <= quantity)
);

create index warehouse_crossdock_allocation_demand_idx
  on inventory.warehouse_crossdock_allocation(tenant_id,organization_id,transfer_id,transfer_line_id,item_id,available,due_date);

create table inventory.warehouse_crossdock_pick_link (
  tenant_id uuid not null,
  pick_activity_id text not null,
  pick_line_id text not null,
  allocation_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  status text not null check (status in ('open','picked','released')),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, pick_activity_id, pick_line_id),
  foreign key (tenant_id, pick_activity_id, pick_line_id)
    references inventory.warehouse_activity_line (tenant_id, activity_id, line_id),
  foreign key (tenant_id, allocation_id)
    references inventory.warehouse_crossdock_allocation (tenant_id, allocation_id)
);

commit;
````

### FILE: `db/migrations/0032_warehouse_transfer_crossdock.down.sql`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:warehouse-crossdock-migration-down:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "5c6949987fc09d2b4a2810d1351396f93b9ebf1415b7441c0985b7e51f2ea0c5"
variables: []
secrets_allowed: false
```

````sql
begin;

do $migration$
begin
  if exists (select 1 from inventory.warehouse_crossdock_policy)
     or exists (select 1 from inventory.warehouse_crossdock_allocation)
     or exists (select 1 from inventory.warehouse_crossdock_pick_link)
     or exists (select 1 from inventory.warehouse_bin where cross_dock) then
    raise exception using errcode='55000', message='cannot remove cross-dock schema after configuration or evidence exists';
  end if;
end;
$migration$;

drop table inventory.warehouse_crossdock_pick_link;
drop table inventory.warehouse_crossdock_allocation;
drop table inventory.warehouse_crossdock_policy;

alter table inventory.warehouse_bin
  drop constraint warehouse_bin_cross_dock_type_check,
  drop column cross_dock;

commit;
````

### FILE: `docs/inventory/MICROSOFT_BC_CROSSDOCK_DERIVATION.md`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:warehouse-crossdock-derivation:v1"
operation: CREATE
provenance: AUTHORED
source: "fixed upstream evidence and local executable contract"
license: "LicenseRef-Workspace-Owner"
sha256: "475a317ca231c5e953ca9586c236b760d8bc2c3e3fbaef5f6cf51ba5eb259ed5"
variables: []
secrets_allowed: false
```

````markdown
# Microsoft Business Central warehouse cross-dock derivation

This portable implementation is `ADAPTED`; its Go, SQL and tests are not Microsoft-authored code. It is governed by the MIT-licensed `microsoft/BCApps` commit `2eae56d704a1fd035d104f333602aea7091b7749` (tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`) and by Microsoft Learn, [Cross-dock items](https://learn.microsoft.com/en-gb/dynamics365/business-central/warehouse-how-to-cross-dock-items). The upstream code and documentation remain the authority; this file states the smaller contract that the portable implementation proves.

## Fixed official inputs

| BCApps path | Bytes | SHA-256 |
|---|---:|---|
| `src/Layers/W1/BaseApp/Warehouse/Activity/WhseCrossDockManagement.Codeunit.al` | 47,115 | `85a1eede4b2132f8ef22c853be1bd66a23a29b5176875a0432f500be150609fd` |
| `src/Layers/W1/BaseApp/Warehouse/Activity/WhseCrossDockOpportunity.Table.al` | 18,740 | `778cc8833ea7c28262a9ae695136281c66002d190845785e6f41e100db274ecb` |
| `src/Layers/W1/BaseApp/Warehouse/Document/WarehouseReceiptHeader.Table.al` | 24,887 | `d02b7f1fa0be0847e10e688d98125fdf34e1f157089a99d4632024d16fd2c0fe` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMWarehouseV.Codeunit.al` | 306,725 | `faee2bf1ef5bae41211f1374b8e1b7d49e9e7e75f2672e70ed0c5b181136f839` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMWarehouseOrders.Codeunit.al` | 238,014 | `79f2e97c278763aa7c140dcc6dffaa53f7e6e6fe2d55100b7f42bad2ccd0a1da` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMWarehouseUnitTests.Codeunit.al` | 170,444 | `9af357c3eddaa6f6acaa50623aec43753d212efe7d412e6a948f35c87521ac3f` |

## Admitted executable contract

- an organization/item has at most one enabled policy, an explicit due-date horizon and an explicit fixed, unblocked PICK/PUTPICK bin marked `cross_dock`;
- only released outbound bulk-transfer demand from the same organization and item is eligible, ordered by due date and stable identities;
- demand beyond the horizon, specific-cost transfer lines and demand with an existing open or registered pick are excluded;
- required quantity is outstanding transfer quantity minus every existing cross-dock allocation, whether or not its receipt has already been put away;
- receipt allocation is the minimum of received quantity, remaining demand and available cross-dock-bin capacity; the unallocated receipt quantity follows normal put-away rules;
- an allocation is not pickable until its exact put-away line is registered; generic picks exclude cross-dock bins;
- a matching transfer pick consumes its own available allocation before normal storage, atomically reserves balance and allocation quantity, and records the exact pick-line link;
- pick cancellation releases both reservations; pick registration moves stock to the ship bin and converts reserved allocation quantity to picked quantity;
- organization/item advisory locking and serializable transactions prevent concurrent double allocation and double reservation.

The PostgreSQL corpus proves partial cross-dock plus ordinary storage, existing-opportunity subtraction, due-date exclusion, unavailable-before-put-away rejection, exact transfer prioritization, cancellation/retry, final pick registration, ledger balances and zero leaked tenant data.

## Explicit non-claims

This version does not claim complete Business Central WMS equivalence. Sales, service and production demand; base-calendar working-day calculation; UOM conversion; breakbulk; cubage/weight; warehouse class; handheld/barcode execution; cross-dock time-zone policy; carrier routing; and external-WMS reconciliation remain conditioned until their authoritative mappings and executable tests exist. No production readiness is inferred from pack admission.
````

### FILE: `internal/inventorycontrol/warehouse_crossdock.go`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:warehouse-crossdock-domain:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "270632118d3c87a85e2591cab4e51d528ebc1d001914260f1c761e9f0ef0c6e4"
variables: []
secrets_allowed: false
```

````go
package inventorycontrol

import (
	"context"
	"fmt"
)

type WarehouseCrossDockPolicy struct {
	OrganizationID string `json:"organization_id"`
	ItemID         string `json:"item_id"`
	BinID          string `json:"bin_id"`
	DueDateDays    int    `json:"due_date_days"`
	Enabled        bool   `json:"enabled"`
	Version        int64  `json:"version"`
}

func (s *WarehouseService) ConfigureCrossDock(ctx context.Context, tenant string, value WarehouseCrossDockPolicy) (WarehouseCrossDockPolicy, error) {
	if tenant == "" || value.OrganizationID == "" || value.ItemID == "" || value.BinID == "" || value.DueDateDays < 0 || value.DueDateDays > 365 || value.Version != 1 {
		return value, fmt.Errorf("invalid warehouse cross-dock policy")
	}
	return s.repository.ConfigureWarehouseCrossDock(ctx, tenant, s.ids.New(), value)
}
````

### FILE: `internal/platform/postgres/warehouse_crossdock.go`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:warehouse-crossdock-postgres:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "95eec888716868b7a6e20fb94b9cf92765138807f7f9fd1694959af211c3da51"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
)

func (r *InventoryControl) ConfigureWarehouseCrossDock(ctx context.Context, tenant, eventID string, value inventorycontrol.WarehouseCrossDockPolicy) (inventorycontrol.WarehouseCrossDockPolicy, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	var crossDock, blocked, fixed bool
	var binType string
	err = tx.QueryRow(ctx, `select w.cross_dock,w.movement_blocked,w.bin_type,p.fixed from inventory.warehouse_bin w join inventory.item_bin_policy p using(tenant_id,organization_id,bin_id) where w.tenant_id=$1 and w.organization_id=$2 and w.bin_id=$3 and p.item_id=$4 for share of w,p`, tenant, value.OrganizationID, value.BinID, value.ItemID).Scan(&crossDock, &blocked, &binType, &fixed)
	if errors.Is(err, pgx.ErrNoRows) || !crossDock || blocked || !fixed || (binType != "pick" && binType != "putpick") {
		return value, inventorycontrol.ErrConflict
	}
	if err != nil {
		return value, err
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_crossdock_policy(tenant_id,organization_id,item_id,bin_id,due_date_days,enabled,version) values($1,$2,$3,$4,$5,$6,1)`, tenant, value.OrganizationID, value.ItemID, value.BinID, value.DueDateDays, value.Enabled)
	if err != nil {
		return value, bulkConflict(err)
	}
	aggregate := value.OrganizationID + ":" + value.ItemID
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "warehouse-crossdock-policy", aggregate, "warehouse-crossdock-policy.created", 1, map[string]any{"organization_id": value.OrganizationID, "item_id": value.ItemID, "bin_id": value.BinID, "due_date_days": value.DueDateDays, "enabled": value.Enabled}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func planWarehouseCrossDock(ctx context.Context, tx pgx.Tx, tenant string, ids inventorycontrol.WarehouseIDs, value inventorycontrol.WarehouseReceiptCommand, lotID string) ([]warehousePlacement, *big.Rat, error) {
	receiptQuantity, ok := parseRat(value.Quantity)
	if !ok {
		return nil, nil, inventorycontrol.ErrConflict
	}
	var binID, maximumText, onHandText, plannedText string
	var dueDays int
	err := tx.QueryRow(ctx, `select p.bin_id,p.due_date_days,coalesce(bp.max_quantity::text,''),coalesce((select sum(b.quantity)::text from inventory.bulk_balance b where b.tenant_id=p.tenant_id and b.organization_id=p.organization_id and b.bin_id=p.bin_id and b.item_id=p.item_id),'0'),coalesce((select sum(al.quantity)::text from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=p.tenant_id and a.organization_id=p.organization_id and a.activity_type='put-away' and a.status='open' and al.to_bin_id=p.bin_id and al.item_id=p.item_id),'0') from inventory.warehouse_crossdock_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) join inventory.item_bin_policy bp on bp.tenant_id=p.tenant_id and bp.organization_id=p.organization_id and bp.item_id=p.item_id and bp.bin_id=p.bin_id where p.tenant_id=$1 and p.organization_id=$2 and p.item_id=$3 and p.enabled and w.cross_dock and not w.movement_blocked and w.bin_type in ('pick','putpick') and bp.fixed for share of p,w,bp`, tenant, value.OrganizationID, value.ItemID).Scan(&binID, &dueDays, &maximumText, &onHandText, &plannedText)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, receiptQuantity, nil
	}
	if err != nil {
		return nil, nil, err
	}
	capacity, valid := decimalCapacity(maximumText, onHandText, plannedText)
	if !valid {
		return nil, nil, inventorycontrol.ErrConflict
	}
	crossDockBudget := new(big.Rat).Set(receiptQuantity)
	if capacity != nil {
		crossDockBudget = minRat(crossDockBudget, capacity)
	}
	if crossDockBudget.Sign() <= 0 {
		return nil, receiptQuantity, nil
	}
	cutoff := value.PostingDate.UTC().AddDate(0, 0, dueDays)
	rows, err := tx.Query(ctx, `select t.transfer_id,l.line_id,t.posting_date,(l.quantity-l.shipped_quantity-coalesce(x.quantity,0)-coalesce(k.quantity,0))::text from inventory.bulk_transfer t join inventory.bulk_transfer_line l using(tenant_id,transfer_id) left join lateral (select sum(a.quantity) quantity from inventory.warehouse_crossdock_allocation a where a.tenant_id=t.tenant_id and a.transfer_id=t.transfer_id and a.transfer_line_id=l.line_id) x on true left join lateral (select sum(al.quantity) quantity from inventory.warehouse_activity wa join inventory.warehouse_activity_line al using(tenant_id,activity_id) left join inventory.warehouse_crossdock_pick_link pl on pl.tenant_id=al.tenant_id and pl.pick_activity_id=al.activity_id and pl.pick_line_id=al.line_id where wa.tenant_id=t.tenant_id and wa.organization_id=t.from_organization_id and wa.activity_type='pick' and wa.source_kind='transfer-outbound' and wa.source_id=t.transfer_id and wa.source_line_id=l.line_id and pl.pick_activity_id is null and (wa.status='open' or (wa.status='registered' and not exists(select 1 from inventory.bulk_transfer_shipment s where s.tenant_id=wa.tenant_id and s.warehouse_activity_id=wa.activity_id)))) k on true where t.tenant_id=$1 and t.from_organization_id=$2 and l.item_id=$3 and t.status='released' and t.posting_date<=$4::date and l.specific_receipt_entry_id is null and l.quantity-l.shipped_quantity-coalesce(x.quantity,0)-coalesce(k.quantity,0)>0 order by t.posting_date,t.created_at,t.transfer_id for update of t,l`, tenant, value.OrganizationID, value.ItemID, cutoff)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	placements := []warehousePlacement{}
	allocated := new(big.Rat)
	for rows.Next() && crossDockBudget.Sign() > 0 {
		var transferID, lineID, neededText string
		var dueDate time.Time
		if err = rows.Scan(&transferID, &lineID, &dueDate, &neededText); err != nil {
			return nil, nil, err
		}
		needed, valid := parseRat(neededText)
		if !valid || needed.Sign() <= 0 {
			return nil, nil, inventorycontrol.ErrConflict
		}
		take := minRat(needed, crossDockBudget)
		index := len(placements) + 1
		placements = append(placements, warehousePlacement{binID: binID, quantity: formatRat(take, 6), crossDockAllocationID: fmt.Sprintf("%s-crossdock-%06d", ids.ActivityID, index), transferID: transferID, transferLineID: lineID, dueDate: dueDate, lotID: lotID})
		crossDockBudget.Sub(crossDockBudget, take)
		allocated.Add(allocated, take)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, err
	}
	if capacity != nil && allocated.Cmp(capacity) > 0 {
		return nil, nil, inventorycontrol.ErrConflict
	}
	remaining := new(big.Rat).Sub(receiptQuantity, allocated)
	return placements, remaining, nil
}
````

### FILE: `internal/platform/postgres/warehouse_crossdock_integration_test.go`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:warehouse-crossdock-postgres-integration:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "0a0ccd8c9eae94f34650641a7829ffde42bc6bfa4400961238e96e6e43cb4fd8"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestWarehouseCrossDockTransferDemandEndToEnd(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := bulkTestUUID(t)
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Cross Dock V166','Cross Dock V166')`, tenant, "xd-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	defer cleanupWarehouseCrossDock(t, pool, tenant)
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'warehouse','warehouse','Warehouse','warehouse')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'franchise','franchise','Franchise','franchisee')`, tenant); err != nil {
		t.Fatal(err)
	}

	repo := NewInventoryControl(pool)
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "part", Code: "PART-XD", Description: "Cross-docked part", BaseUOM: "EA", BaseRoundingPrecision: "0.000001", TrackingMode: "lot", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.WarehouseBin{ID: "invalid-cross-dock", OrganizationID: "warehouse", Code: "INVALID-XD", Type: "receive", CrossDock: true, Version: 1}); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("cross-dock receive bin accepted: %v", err)
	}
	bins := []inventorycontrol.WarehouseBin{
		{ID: "receive", OrganizationID: "warehouse", Code: "RECEIVE", Type: "receive", Version: 1},
		{ID: "storage", OrganizationID: "warehouse", Code: "STORAGE", Type: "putpick", Ranking: 20, Version: 1},
		{ID: "cross-dock", OrganizationID: "warehouse", Code: "CROSS-DOCK", Type: "putpick", Ranking: 100, CrossDock: true, Version: 1},
		{ID: "ship", OrganizationID: "warehouse", Code: "SHIP", Type: "ship", Version: 1},
		{ID: "destination-receive", OrganizationID: "franchise", Code: "RECEIVE", Type: "receive", Version: 1},
	}
	for _, bin := range bins {
		if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatalf("create bin %s: %v", bin.ID, err)
		}
		maximum := "100"
		if bin.CrossDock {
			maximum = "5"
		}
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: bin.OrganizationID, ItemID: "part", BinID: bin.ID, Fixed: true, Default: bin.ID == "storage", MinQuantity: "0", MaxQuantity: maximum, Version: 1}); err != nil {
			t.Fatalf("configure bin %s: %v", bin.ID, err)
		}
	}
	policy := inventorycontrol.WarehouseCrossDockPolicy{OrganizationID: "warehouse", ItemID: "part", BinID: "cross-dock", DueDateDays: 5, Enabled: true, Version: 1}
	invalidPolicy := policy
	invalidPolicy.BinID = "storage"
	if _, err = repo.ConfigureWarehouseCrossDock(ctx, tenant, bulkTestUUID(t), invalidPolicy); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("ordinary storage accepted as cross-dock policy: %v", err)
	}
	if _, err = repo.ConfigureWarehouseCrossDock(ctx, tenant, bulkTestUUID(t), policy); err != nil {
		t.Fatal(err)
	}

	eligible, err := repo.CreateBulkTransfer(ctx, tenant, "eligible-transfer", "eligible-line", bulkTestUUID(t), inventorycontrol.BulkTransferCommand{RequestID: "eligible-request", FromOrganizationID: "warehouse", ToOrganizationID: "franchise", ReceiveBinID: "destination-receive", InTransitCode: "ROAD", ItemID: "part", Quantity: "5", PostingDate: time.Date(2030, 1, 3, 0, 0, 0, 0, time.UTC), SourceKind: "replenishment", SourceID: "eligible-source"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.CreateBulkTransfer(ctx, tenant, "late-transfer", "late-line", bulkTestUUID(t), inventorycontrol.BulkTransferCommand{RequestID: "late-request", FromOrganizationID: "warehouse", ToOrganizationID: "franchise", ReceiveBinID: "destination-receive", InTransitCode: "ROAD", ItemID: "part", Quantity: "10", PostingDate: time.Date(2030, 1, 20, 0, 0, 0, 0, time.UTC), SourceKind: "replenishment", SourceID: "late-source"}); err != nil {
		t.Fatal(err)
	}
	expires := time.Date(2035, 1, 1, 0, 0, 0, 0, time.UTC)
	if _, err = pool.Exec(ctx, `insert into inventory.inventory_lot(tenant_id,lot_id,item_id,lot_no,expiration_date) values($1,'initial-lot','part','LOT-INITIAL',$2::date)`, tenant, expires); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity) values($1,'warehouse','storage','part','initial-lot',2,0)`, tenant); err != nil {
		t.Fatal(err)
	}
	genericPick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "generic-transfer-pick"), inventorycontrol.WarehousePickCommand{RequestID: "generic-transfer-pick-request", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "transfer-outbound", DemandID: eligible.ID, DemandLineID: eligible.LineID, Quantity: "2", UseFEFO: true})
	if err != nil || len(genericPick.Lines) != 1 || genericPick.Lines[0].FromBinID != "storage" {
		t.Fatalf("generic transfer pick=%+v err=%v", genericPick, err)
	}
	first, err := repo.PostWarehouseReceipt(ctx, tenant, warehouseTestIDs(t, "cross-first"), inventorycontrol.WarehouseReceiptCommand{RequestID: "cross-first-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "part", LotNo: "LOT-XD-1", ExpirationDate: &expires, Quantity: "7", UnitCost: "10", PostingDate: time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC), SourceKind: "purchase", SourceID: "po-cross-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(first.PutAway.Lines) != 2 || first.PutAway.Lines[0].ToBinID != "cross-dock" || first.PutAway.Lines[0].Quantity != "3" || first.PutAway.Lines[1].ToBinID != "storage" || first.PutAway.Lines[1].Quantity != "4" {
		t.Fatalf("partial cross-dock placement=%+v", first.PutAway.Lines)
	}
	second, err := repo.PostWarehouseReceipt(ctx, tenant, warehouseTestIDs(t, "cross-second"), inventorycontrol.WarehouseReceiptCommand{RequestID: "cross-second-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "part", LotNo: "LOT-XD-2", ExpirationDate: &expires, Quantity: "3", UnitCost: "11", PostingDate: time.Date(2030, 1, 2, 0, 0, 0, 0, time.UTC), SourceKind: "purchase", SourceID: "po-cross-2"})
	if err != nil {
		t.Fatal(err)
	}
	if len(second.PutAway.Lines) != 1 || second.PutAway.Lines[0].ToBinID != "storage" || second.PutAway.Lines[0].Quantity != "3" {
		t.Fatalf("existing opportunity was duplicated: %+v", second.PutAway.Lines)
	}
	var allocations int
	var allocated string
	if err = pool.QueryRow(ctx, `select count(*),sum(quantity)::text from inventory.warehouse_crossdock_allocation where tenant_id=$1`, tenant).Scan(&allocations, &allocated); err != nil || allocations != 1 || allocated != "3.000000" {
		t.Fatalf("allocations=%d quantity=%s err=%v", allocations, allocated, err)
	}
	if _, err = repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "unavailable-pick"), inventorycontrol.WarehousePickCommand{RequestID: "unavailable-pick-request", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "transfer-outbound", DemandID: eligible.ID, DemandLineID: eligible.LineID, Quantity: "1", UseFEFO: true}); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("unregistered cross-dock allocation was pickable: %v", err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", first.PutAway.ID, 1, time.Date(2030, 1, 2, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "cross-over-demand"), inventorycontrol.WarehousePickCommand{RequestID: "cross-over-demand-request", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "transfer-outbound", DemandID: eligible.ID, DemandLineID: eligible.LineID, Quantity: "4", UseFEFO: true}); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("over-demand transfer pick accepted: %v", err)
	}

	cancelPick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "cross-cancel"), inventorycontrol.WarehousePickCommand{RequestID: "cross-cancel-request", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "transfer-outbound", DemandID: eligible.ID, DemandLineID: eligible.LineID, Quantity: "1", UseFEFO: true})
	if err != nil || len(cancelPick.Lines) != 1 || cancelPick.Lines[0].FromBinID != "cross-dock" {
		t.Fatalf("cross-dock cancellation pick=%+v err=%v", cancelPick, err)
	}
	if _, err = repo.CancelWarehousePick(ctx, tenant, "warehouse", cancelPick.ID, 1, bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	var reserved string
	var releasedLinks int
	if err = pool.QueryRow(ctx, `select (select reserved_quantity::text from inventory.warehouse_crossdock_allocation where tenant_id=$1),(select count(*) from inventory.warehouse_crossdock_pick_link where tenant_id=$1 and status='released')`, tenant).Scan(&reserved, &releasedLinks); err != nil || reserved != "0.000000" || releasedLinks != 1 {
		t.Fatalf("cancel reserved=%s released_links=%d err=%v", reserved, releasedLinks, err)
	}

	pick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "cross-pick"), inventorycontrol.WarehousePickCommand{RequestID: "cross-pick-request", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "transfer-outbound", DemandID: eligible.ID, DemandLineID: eligible.LineID, Quantity: "3", UseFEFO: true})
	if err != nil || len(pick.Lines) != 1 || pick.Lines[0].FromBinID != "cross-dock" || pick.Lines[0].Quantity != "3" {
		t.Fatalf("cross-dock pick=%+v err=%v", pick, err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", pick.ID, 1, time.Date(2030, 1, 3, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	var picked, crossQuantity, storageQuantity, shipQuantity string
	var pickedLinks, lateAllocations int
	if err = pool.QueryRow(ctx, `select x.reserved_quantity::text,x.picked_quantity::text,(select sum(quantity)::text from inventory.bulk_balance where tenant_id=$1 and organization_id='warehouse' and bin_id='cross-dock'),(select sum(quantity)::text from inventory.bulk_balance where tenant_id=$1 and organization_id='warehouse' and bin_id='storage'),(select sum(quantity)::text from inventory.bulk_balance where tenant_id=$1 and organization_id='warehouse' and bin_id='ship'),(select count(*) from inventory.warehouse_crossdock_pick_link where tenant_id=$1 and status='picked'),(select count(*) from inventory.warehouse_crossdock_allocation where tenant_id=$1 and transfer_id='late-transfer') from inventory.warehouse_crossdock_allocation x where x.tenant_id=$1`, tenant).Scan(&reserved, &picked, &crossQuantity, &storageQuantity, &shipQuantity, &pickedLinks, &lateAllocations); err != nil {
		t.Fatal(err)
	}
	if reserved != "0.000000" || picked != "3.000000" || crossQuantity != "0.000000" || storageQuantity != "6.000000" || shipQuantity != "3.000000" || pickedLinks != 1 || lateAllocations != 0 {
		t.Fatalf("final reserved=%s picked=%s cross=%s storage=%s ship=%s picked_links=%d late=%d", reserved, picked, crossQuantity, storageQuantity, shipQuantity, pickedLinks, lateAllocations)
	}
}

func cleanupWarehouseCrossDock(t *testing.T, pool *pgxpool.Pool, tenant string) {
	t.Helper()
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Errorf("cross-dock cleanup begin: %v", err)
		return
	}
	defer tx.Rollback(ctx)
	commands := []string{
		`alter table inventory.warehouse_pick_request disable trigger warehouse_pick_request_immutable`,
		`alter table inventory.warehouse_activity_uom_reservation disable trigger warehouse_uom_reservation_delete_guard`,
		`alter table inventory.warehouse_activity_line disable trigger warehouse_activity_line_immutable`,
		`alter table inventory.warehouse_receipt_line disable trigger warehouse_receipt_line_immutable`,
		`alter table inventory.warehouse_receipt disable trigger warehouse_receipt_immutable`,
		`alter table inventory.bulk_inventory_entry disable trigger bulk_inventory_entry_immutable`,
		`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`,
		`delete from inventory.warehouse_activity_uom_reservation where tenant_id=$1`,
		`delete from inventory.warehouse_crossdock_pick_link where tenant_id=$1`,
		`delete from inventory.warehouse_crossdock_allocation where tenant_id=$1`,
		`delete from inventory.warehouse_crossdock_policy where tenant_id=$1`,
		`delete from inventory.warehouse_pick_request where tenant_id=$1`,
		`delete from inventory.warehouse_activity_line where tenant_id=$1`,
		`delete from inventory.warehouse_activity where tenant_id=$1`,
		`delete from inventory.warehouse_receipt_line where tenant_id=$1`,
		`delete from inventory.warehouse_receipt where tenant_id=$1`,
		`delete from inventory.bulk_transfer_line where tenant_id=$1`,
		`delete from inventory.bulk_transfer where tenant_id=$1`,
		`delete from inventory.bulk_cost_layer where tenant_id=$1`,
		`delete from inventory.bulk_inventory_entry where tenant_id=$1`,
		`delete from inventory.bulk_reservation where tenant_id=$1`,
		`delete from inventory.bulk_balance where tenant_id=$1`,
		`delete from inventory.inventory_lot where tenant_id=$1`,
		`delete from inventory.item_bin_policy where tenant_id=$1`,
		`delete from inventory.warehouse_bin where tenant_id=$1`,
		`delete from inventory.item_unit_of_measure where tenant_id=$1`,
		`delete from inventory.stock_item where tenant_id=$1`,
		`delete from platform.outbox_event where tenant_id=$1`,
		`delete from org.organization where tenant_id=$1`,
		`delete from platform.tenant where tenant_id=$1`,
		`alter table inventory.warehouse_activity_line enable trigger warehouse_activity_line_immutable`,
		`alter table inventory.warehouse_receipt_line enable trigger warehouse_receipt_line_immutable`,
		`alter table inventory.warehouse_receipt enable trigger warehouse_receipt_immutable`,
		`alter table inventory.bulk_inventory_entry enable trigger bulk_inventory_entry_immutable`,
		`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`,
		`alter table inventory.warehouse_activity_uom_reservation enable trigger warehouse_uom_reservation_delete_guard`,
		`alter table inventory.warehouse_pick_request enable trigger warehouse_pick_request_immutable`,
	}
	for _, command := range commands {
		args := []any{}
		if strings.Contains(command, "$1") {
			args = append(args, tenant)
		}
		if _, err = tx.Exec(ctx, command, args...); err != nil {
			t.Errorf("cross-dock cleanup failed: %v", err)
			return
		}
	}
	if err = tx.Commit(ctx); err != nil {
		t.Errorf("cross-dock cleanup commit: %v", err)
	}
}
````

### FILE: `db/migrations/0033_item_unit_of_measure.up.sql`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:item-uom-migration-up:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d ItemUnitofMeasure.Table.al narrow invariant adaptation"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "55d3d01ab5474be6ecd57501c398ef380d3c558e30e440290dcda5bade6289c9"
variables: []
secrets_allowed: false
```

````sql
begin;

create table inventory.item_unit_of_measure (
  tenant_id uuid not null,
  item_id text not null,
  uom_code text not null,
  qty_per_uom numeric(20,6) not null check (qty_per_uom > 0),
  rounding_precision numeric(20,6) not null check (rounding_precision > 0 and rounding_precision <= 1),
  is_base boolean not null default false,
  version bigint not null default 1 check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, item_id, uom_code),
  foreign key (tenant_id, item_id) references inventory.stock_item (tenant_id, item_id),
  check (uom_code ~ '^[A-Z0-9][A-Z0-9._/-]{0,15}$'),
  check (not is_base or qty_per_uom = 1),
  check (mod(qty_per_uom, rounding_precision) = 0)
);

create unique index item_unit_of_measure_one_base_uidx
  on inventory.item_unit_of_measure(tenant_id, item_id)
  where is_base;

create function inventory.reject_item_unit_of_measure_mutation()
returns trigger language plpgsql as $$
begin
  raise exception using errcode='55000', message='item unit of measure is immutable; create a new item or approved code instead';
end;
$$;

create trigger item_unit_of_measure_immutable
before update or delete on inventory.item_unit_of_measure
for each row execute function inventory.reject_item_unit_of_measure_mutation();

commit;
````

### FILE: `db/migrations/0033_item_unit_of_measure.down.sql`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:item-uom-migration-down:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d ItemUnitofMeasure.Table.al narrow invariant adaptation"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "2e25a62d9943f1252bffb5f537705c19402a0685640097c9d841a174b49288fb"
variables: []
secrets_allowed: false
```

````sql
begin;

drop trigger if exists item_unit_of_measure_immutable on inventory.item_unit_of_measure;
drop function if exists inventory.reject_item_unit_of_measure_mutation();
drop table if exists inventory.item_unit_of_measure;

commit;
````

### FILE: `docs/inventory/MICROSOFT_BC_ITEM_UOM_DERIVATION.md`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:item-uom-derivation:v1"
operation: CREATE
provenance: AUTHORED
source: "fixed Microsoft BCApps commit and Microsoft Learn authorities"
license: "LicenseRef-Workspace-Owner"
sha256: "91aaa0489e0f7708cdcb000a27754a72bc4b46bbb568fe0535ff839399f7f848"
variables: []
secrets_allowed: false
```

````markdown
# Microsoft Business Central item UOM derivation

## Fixed authority

- Repository: `microsoft/BCApps`
- Commit: `2eae56d704a1fd035d104f333602aea7091b7749`
- Tree: `f9846fb1254c6311985c6131bb2af11c7f169e1b`
- License: MIT
- Product guidance:
  - <https://learn.microsoft.com/en-gb/dynamics365/business-central/inventory-how-setup-units-of-measure>
  - <https://learn.microsoft.com/en-ca/dynamics365/business-central/warehouse-enable-automatic-breaking-bulk-with-directed-put-away-and-pick>
  - <https://learn.microsoft.com/en-ca/dynamics365/business-central/design-details-warehouse-management>

| Official file | Bytes | SHA-256 |
|---|---:|---|
| `src/Layers/W1/BaseApp/Inventory/Item/ItemUnitofMeasure.Table.al` | 20,689 | `ab785d08d4bbc9160c2e068b5660ada732f752d35ae7c94a6677ef0fb79b2009` |
| `src/Layers/W1/BaseApp/Warehouse/Activity/CreatePutaway.Codeunit.al` | 77,370 | `b004c0663d237585d8b57a06356431acb2302cdfbbd30f7b82154d282f325b2b` |
| `src/Layers/W1/BaseApp/Warehouse/Activity/WhseChangeUnitofMeasure.Report.al` | 8,008 | `d1fca83ae830a8b615c558e15cb6ec3734560614ba66e86d801f3e9794b06186` |
| `src/Layers/W1/BaseApp/Warehouse/Worksheet/Replenishment.Codeunit.al` | 19,064 | `e6d5e3873b4828db8dca49755670bd9815faa7ad1b0396189014af57380b7c4e` |
| `src/Layers/W1/BaseApp/Warehouse/Worksheet/WhseCreatePick.Codeunit.al` | 1,584 | `43ffaded138adb24d0cfc927c276d63ef3dfb2760f83f518bcbc8bbc217163bc` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMMovement.Codeunit.al` | 89,432 | `0e63996a0816d9379bdf76acc5712990f10624de81d875f2e499e6944b2955ed` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMWMSItemUnitofMeasure.Codeunit.al` | 125,924 | `8faec658fad625a12c83b251bab88553319bc1cf564e7671104f6a19cadf4713` |

## Narrow portable contract

This pack does not copy AL into Go and does not claim that Microsoft supports the adaptation. It transfers only the following demonstrated invariants:

1. Inventory authority is expressed in one base UOM. Its factor is exactly `1`.
2. Every alternate factor is positive and must align exactly with the base quantity rounding precision.
3. A conversion is `requested quantity × quantity per UOM`; it is admitted only when the result has no residual against base precision.
4. UOM identity, factor and precision are immutable in the portable owner. This is stricter than changing only when no warehouse or open-document evidence exists, and prevents historical reinterpretation.
5. Duplicate codes, unknown codes, misaligned factors and residual quantities fail closed before any inventory effect.

Microsoft BCApps further demonstrates same-UOM-first selection, explicitly enabled breakbulk, larger-package fallback, smaller-UOM gathering and paired warehouse Take/Place evidence. V168 adds an authoritative composition plus explicit breakbulk/gather with paired evidence; automatic selection inside receipt, movement, pick, registration and cancellation remains unclaimed until connected and tested end to end.

## Reconstruction and verification

Migration `0033` owns the immutable item-UOM table and one-base constraint. Item creation inserts its base record in the same transaction as the item and outbox event. The service exposes strict alternate configuration and exact conversion. PostgreSQL integration proves factor `12`, `2.5 BOX = 30 EA`, rejection of `0.1 BOX` for an integer base, rejection of factor `2.5`, duplicate rejection and SQLSTATE `55000` on mutation. Migration `0034` supplies the separately documented physical-composition owner and reversible legacy backfill. HTTP tests prove strict authenticated configuration and conversion surfaces.
````

### FILE: `internal/inventorycontrol/item_uom.go`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:item-uom-domain:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d ItemUnitofMeasure.Table.al narrow invariant adaptation"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "a9932a19ca1187a5f57d186ccf868b30d04457c70cfd9c5c804482655aa55ef3"
variables: []
secrets_allowed: false
```

````go
package inventorycontrol

import (
	"context"
	"fmt"
)

type ItemUnitOfMeasure struct {
	ItemID            string `json:"item_id"`
	Code              string `json:"code"`
	QuantityPerUnit   string `json:"quantity_per_unit"`
	RoundingPrecision string `json:"rounding_precision"`
	Base              bool   `json:"base"`
	Version           int64  `json:"version"`
}

type UnitOfMeasureConversion struct {
	ItemID            string `json:"item_id"`
	Code              string `json:"code"`
	Quantity          string `json:"quantity"`
	QuantityPerUnit   string `json:"quantity_per_unit"`
	BaseCode          string `json:"base_code"`
	BaseQuantity      string `json:"base_quantity"`
	RoundingPrecision string `json:"rounding_precision"`
}

func (s *BulkService) ConfigureUnitOfMeasure(ctx context.Context, tenant string, value ItemUnitOfMeasure) (ItemUnitOfMeasure, error) {
	if tenant == "" || value.ItemID == "" || value.Code == "" || value.Base || !positiveDecimal(value.QuantityPerUnit, quantityPattern) || !positiveDecimal(value.RoundingPrecision, quantityPattern) {
		return ItemUnitOfMeasure{}, fmt.Errorf("invalid item unit of measure")
	}
	value.Version = 1
	return s.repository.ConfigureItemUnitOfMeasure(ctx, tenant, s.ids.New(), value)
}

func (s *BulkService) ConvertUnitOfMeasure(ctx context.Context, tenant, itemID, code, quantity string) (UnitOfMeasureConversion, error) {
	if tenant == "" || itemID == "" || code == "" || !positiveDecimal(quantity, quantityPattern) {
		return UnitOfMeasureConversion{}, fmt.Errorf("invalid unit of measure conversion")
	}
	return s.repository.ConvertItemUnitOfMeasure(ctx, tenant, itemID, code, quantity)
}
````

### FILE: `internal/inventorycontrol/item_uom_test.go`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:item-uom-domain-test:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d SCMWMSItemUnitofMeasure.Codeunit.al negative invariant adaptation"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "49d5bea0bf87644d999e9f645c385173a88813497f9581b1e7e3282a52bbf040"
variables: []
secrets_allowed: false
```

````go
package inventorycontrol

import (
	"context"
	"testing"
)

func TestUnitOfMeasureServiceFailsClosedBeforeRepository(t *testing.T) {
	repo := &bulkFake{}
	service := NewBulkService(repo, inventoryIDsForBulk{})
	invalid := []ItemUnitOfMeasure{
		{},
		{ItemID: "part", Code: "BOX", QuantityPerUnit: "12", RoundingPrecision: "1", Base: true},
		{ItemID: "part", Code: "BOX", QuantityPerUnit: "0", RoundingPrecision: "1"},
		{ItemID: "part", Code: "BOX", QuantityPerUnit: "12", RoundingPrecision: "0"},
	}
	for _, value := range invalid {
		if _, err := service.ConfigureUnitOfMeasure(context.Background(), "tenant", value); err == nil {
			t.Fatalf("accepted invalid UOM: %+v", value)
		}
	}
	if _, err := service.ConvertUnitOfMeasure(context.Background(), "tenant", "part", "BOX", "0"); err == nil {
		t.Fatal("accepted zero conversion")
	}
}
````

### FILE: `internal/platform/postgres/item_uom.go`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:item-uom-postgres:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d UOM quantity-base conversion invariant adaptation"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "f872b285d8d9b9068164995486b7cb6256b9eb35b3ddfff3dc5322ecd92e9e17"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
)

func (r *InventoryControl) ConfigureItemUnitOfMeasure(ctx context.Context, tenant, eventID string, value inventorycontrol.ItemUnitOfMeasure) (inventorycontrol.ItemUnitOfMeasure, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into inventory.item_unit_of_measure(tenant_id,item_id,uom_code,qty_per_uom,rounding_precision,is_base,version) values($1,$2,$3,$4::numeric,$5::numeric,false,1)`, tenant, value.ItemID, value.Code, value.QuantityPerUnit, value.RoundingPrecision)
	if err != nil {
		return value, bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "item-unit-of-measure", value.ItemID+":"+value.Code, "item-unit-of-measure.created", 1, map[string]string{"item_id": value.ItemID, "code": value.Code, "quantity_per_unit": value.QuantityPerUnit, "rounding_precision": value.RoundingPrecision}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func (r *InventoryControl) ConvertItemUnitOfMeasure(ctx context.Context, tenant, itemID, code, quantity string) (inventorycontrol.UnitOfMeasureConversion, error) {
	result := inventorycontrol.UnitOfMeasureConversion{ItemID: itemID, Code: code, Quantity: quantity}
	err := r.pool.QueryRow(ctx, `select u.qty_per_uom::text,i.base_uom,(($4::numeric)*u.qty_per_uom)::numeric(20,6)::text,b.rounding_precision::text from inventory.item_unit_of_measure u join inventory.stock_item i using(tenant_id,item_id) join inventory.item_unit_of_measure b on b.tenant_id=u.tenant_id and b.item_id=u.item_id and b.is_base where u.tenant_id=$1 and u.item_id=$2 and u.uom_code=$3 and mod(($4::numeric)*u.qty_per_uom,b.rounding_precision)=0`, tenant, itemID, code, quantity).Scan(&result.QuantityPerUnit, &result.BaseCode, &result.BaseQuantity, &result.RoundingPrecision)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.UnitOfMeasureConversion{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.UnitOfMeasureConversion{}, err
	}
	return result, nil
}
````

### FILE: `internal/platform/postgres/item_uom_integration_test.go`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:item-uom-postgres-integration:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d item UOM and breakbulk residual scenarios narrow adaptation"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "e38d0154ce8a3c7279171c996fd5097f10376fc3c9cb5b26593b18002114e911"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"os"
	"testing"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestItemUnitOfMeasureExactConversionAndImmutability(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := bulkTestUUID(t)
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'UOM V167','UOM V167')`, tenant, "uom-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	defer func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			t.Errorf("UOM cleanup begin: %v", cleanupErr)
			return
		}
		defer tx.Rollback(ctx)
		for _, command := range []string{
			`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`,
			`delete from inventory.item_unit_of_measure where tenant_id=$1`,
			`delete from platform.outbox_event where tenant_id=$1`,
			`delete from inventory.stock_item where tenant_id=$1`,
			`delete from platform.tenant where tenant_id=$1`,
			`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`,
		} {
			args := []any{}
			if command[0:6] == "delete" {
				args = append(args, tenant)
			}
			if _, cleanupErr = tx.Exec(ctx, command, args...); cleanupErr != nil {
				t.Errorf("UOM cleanup: %v", cleanupErr)
				return
			}
		}
		if cleanupErr = tx.Commit(ctx); cleanupErr != nil {
			t.Errorf("UOM cleanup commit: %v", cleanupErr)
		}
	}()

	repo := NewInventoryControl(pool)
	item := inventorycontrol.BulkItem{ID: "part", Code: "UOM-PART", Description: "Discrete part", BaseUOM: "EA", BaseRoundingPrecision: "1", TrackingMode: "none", CostingMethod: "fifo", Version: 1}
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), item); err != nil {
		t.Fatal(err)
	}
	box := inventorycontrol.ItemUnitOfMeasure{ItemID: "part", Code: "BOX", QuantityPerUnit: "12", RoundingPrecision: "1", Version: 1}
	if _, err = repo.ConfigureItemUnitOfMeasure(ctx, tenant, bulkTestUUID(t), box); err != nil {
		t.Fatal(err)
	}
	converted, err := repo.ConvertItemUnitOfMeasure(ctx, tenant, "part", "BOX", "2.5")
	if err != nil || converted.BaseCode != "EA" || converted.BaseQuantity != "30.000000" || converted.QuantityPerUnit != "12.000000" || converted.RoundingPrecision != "1.000000" {
		t.Fatalf("conversion=%+v err=%v", converted, err)
	}
	if _, err = repo.ConvertItemUnitOfMeasure(ctx, tenant, "part", "BOX", "0.1"); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("residual conversion error=%v", err)
	}
	bad := inventorycontrol.ItemUnitOfMeasure{ItemID: "part", Code: "INNER", QuantityPerUnit: "2.5", RoundingPrecision: "1", Version: 1}
	if _, err = repo.ConfigureItemUnitOfMeasure(ctx, tenant, bulkTestUUID(t), bad); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("misaligned factor error=%v", err)
	}
	if _, err = repo.ConfigureItemUnitOfMeasure(ctx, tenant, bulkTestUUID(t), box); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("duplicate UOM error=%v", err)
	}
	_, err = pool.Exec(ctx, `update inventory.item_unit_of_measure set qty_per_uom=24 where tenant_id=$1 and item_id='part' and uom_code='BOX'`, tenant)
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "55000" {
		t.Fatalf("immutable factor error=%v", err)
	}
}
````

### FILE: `internal/platform/httpapi/item_uom.go`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:item-uom-http:v1"
operation: CREATE
provenance: AUTHORED
source: "local authenticated transport over the verified item UOM owner"
license: "LicenseRef-Workspace-Owner"
sha256: "3183beeec1c9232488ede88290b84c2f385d7f2534e8278ff2b547cf68f4586e"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"net/http"

	"elite.local/enterprise/internal/inventorycontrol"
)

func (a bulkInventoryAPI) configureUnitOfMeasure(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.ItemUnitOfMeasure
	if !decodeStrict(w, r, &input) {
		return
	}
	value, err := a.service.ConfigureUnitOfMeasure(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_ITEM_UNIT_OF_MEASURE") {
		return
	}
	writeJSON(w, 201, value)
}

func (a bulkInventoryAPI) convertUnitOfMeasure(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:read", true)
	if !ok {
		return
	}
	var input struct {
		ItemID   string `json:"item_id"`
		Code     string `json:"code"`
		Quantity string `json:"quantity"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	value, err := a.service.ConvertUnitOfMeasure(r.Context(), p.TenantID, input.ItemID, input.Code, input.Quantity)
	if writeBulkError(w, err, "INVALID_UNIT_OF_MEASURE_CONVERSION") {
		return
	}
	writeJSON(w, 200, value)
}
````

### FILE: `db/migrations/0034_bulk_uom_packaging.up.sql`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:bulk-uom-packaging-migration-up:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d breakbulk/gather and Take/Place invariant adaptation"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "d98f11843b6bc186a3da34daeac6eb85372dd1ab58d13467902246c6170d6d7f"
variables: []
secrets_allowed: false
```

````sql
begin;

create table inventory.bulk_uom_migration_backfill (
  tenant_id uuid not null,
  item_id text not null,
  uom_code text not null,
  primary key (tenant_id,item_id),
  foreign key (tenant_id,item_id) references inventory.stock_item(tenant_id,item_id)
);

insert into inventory.bulk_uom_migration_backfill(tenant_id,item_id,uom_code)
select i.tenant_id,i.item_id,i.base_uom
  from inventory.stock_item i
 where not exists (select 1 from inventory.item_unit_of_measure u where u.tenant_id=i.tenant_id and u.item_id=i.item_id and u.is_base);

insert into inventory.item_unit_of_measure(tenant_id,item_id,uom_code,qty_per_uom,rounding_precision,is_base,version)
select tenant_id,item_id,uom_code,1,0.000001,true,1
  from inventory.bulk_uom_migration_backfill;

create table inventory.bulk_uom_balance (
  balance_id bigint not null references inventory.bulk_balance(balance_id) on delete cascade,
  uom_code text not null,
  quantity numeric(20,6) not null check (quantity >= 0),
  qty_per_uom numeric(20,6) not null check (qty_per_uom > 0),
  quantity_base numeric(20,6) not null check (quantity_base >= 0),
  version bigint not null default 1 check (version > 0),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (balance_id, uom_code),
  check (quantity * qty_per_uom = quantity_base)
);

create table inventory.bulk_uom_conversion (
  tenant_id uuid not null,
  conversion_id text not null,
  request_id text not null,
  organization_id text not null,
  bin_id text not null,
  item_id text not null,
  lot_id text,
  operation text not null check (operation in ('breakbulk','gather')),
  from_uom text not null,
  to_uom text not null,
  from_quantity numeric(20,6) not null check (from_quantity > 0),
  to_quantity numeric(20,6) not null check (to_quantity > 0),
  base_quantity numeric(20,6) not null check (base_quantity > 0),
  from_qty_per_uom numeric(20,6) not null check (from_qty_per_uom > 0),
  to_qty_per_uom numeric(20,6) not null check (to_qty_per_uom > 0),
  version bigint not null default 1 check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, conversion_id),
  unique (tenant_id, request_id),
  foreign key (tenant_id, item_id, from_uom) references inventory.item_unit_of_measure(tenant_id,item_id,uom_code),
  foreign key (tenant_id, item_id, to_uom) references inventory.item_unit_of_measure(tenant_id,item_id,uom_code),
  check (from_uom <> to_uom),
  check (from_quantity * from_qty_per_uom = base_quantity),
  check (to_quantity * to_qty_per_uom = base_quantity),
  check ((operation = 'breakbulk' and from_qty_per_uom > to_qty_per_uom) or (operation = 'gather' and from_qty_per_uom < to_qty_per_uom))
);

create table inventory.bulk_uom_conversion_line (
  tenant_id uuid not null,
  conversion_id text not null,
  line_id text not null,
  sequence_no smallint not null check (sequence_no in (1,2)),
  action_type text not null check (action_type in ('take','place')),
  uom_code text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  quantity_base numeric(20,6) not null check (quantity_base > 0),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, conversion_id, line_id),
  unique (tenant_id, conversion_id, sequence_no),
  foreign key (tenant_id, conversion_id) references inventory.bulk_uom_conversion(tenant_id,conversion_id),
  check ((sequence_no=1 and action_type='take') or (sequence_no=2 and action_type='place'))
);

create function inventory.validate_bulk_uom_balance()
returns trigger language plpgsql as $$
declare
  expected_factor numeric(20,6);
  expected_precision numeric(20,6);
begin
  select u.qty_per_uom,u.rounding_precision into expected_factor,expected_precision
    from inventory.bulk_balance b
    join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.uom_code=new.uom_code
   where b.balance_id=new.balance_id;
  if expected_factor is null or new.qty_per_uom<>expected_factor or mod(new.quantity,expected_precision)<>0 then
    raise exception using errcode='55000',message='bulk UOM composition does not match the immutable item UOM contract';
  end if;
  return new;
end;
$$;

create trigger bulk_uom_balance_validate
before insert or update on inventory.bulk_uom_balance
for each row execute function inventory.validate_bulk_uom_balance();

create function inventory.sync_bulk_base_uom()
returns trigger language plpgsql as $$
declare
  delta numeric(20,6);
  base_code text;
  base_factor numeric(20,6);
  base_precision numeric(20,6);
begin
  delta := new.quantity - case when tg_op='INSERT' then 0 else old.quantity end;
  if delta=0 then return new; end if;
  select u.uom_code,u.qty_per_uom,u.rounding_precision into base_code,base_factor,base_precision
    from inventory.item_unit_of_measure u
   where u.tenant_id=new.tenant_id and u.item_id=new.item_id and u.is_base;
  if base_code is null then
    raise exception using errcode='55000',message='base UOM missing for bulk balance';
  end if;
  if delta>0 then
    insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base)
    values(new.balance_id,base_code,delta/base_factor,base_factor,delta)
    on conflict(balance_id,uom_code) do update
      set quantity=inventory.bulk_uom_balance.quantity+excluded.quantity,
          quantity_base=inventory.bulk_uom_balance.quantity_base+excluded.quantity_base,
          version=inventory.bulk_uom_balance.version+1,updated_at=clock_timestamp();
  else
    update inventory.bulk_uom_balance
       set quantity=quantity-((-delta)/qty_per_uom),quantity_base=quantity_base+delta,
           version=version+1,updated_at=clock_timestamp()
     where balance_id=new.balance_id and uom_code=base_code and quantity_base>=-delta
       and mod((-delta)/qty_per_uom,base_precision)=0;
    if not found then
      raise exception using errcode='55000',message='bulk operation requires explicit handling-unit conversion';
    end if;
    delete from inventory.bulk_uom_balance where balance_id=new.balance_id and uom_code=base_code and quantity_base=0;
  end if;
  return new;
end;
$$;

create trigger bulk_balance_sync_base_uom
after insert or update of quantity on inventory.bulk_balance
for each row execute function inventory.sync_bulk_base_uom();

insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base)
select b.balance_id,u.uom_code,b.quantity/u.qty_per_uom,u.qty_per_uom,b.quantity
  from inventory.bulk_balance b
  join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.is_base
 where b.quantity>0;

create function inventory.assert_bulk_uom_conservation()
returns trigger language plpgsql as $$
declare
  affected bigint;
  physical numeric(20,6);
  composed numeric(20,6);
begin
  affected := case when tg_table_name='bulk_balance' then coalesce(new.balance_id,old.balance_id) else coalesce(new.balance_id,old.balance_id) end;
  select quantity into physical from inventory.bulk_balance where balance_id=affected;
  if physical is null then return null; end if;
  select coalesce(sum(quantity_base),0) into composed from inventory.bulk_uom_balance where balance_id=affected;
  if physical<>composed then
    raise exception using errcode='55000',message='bulk physical quantity and UOM composition diverged';
  end if;
  return null;
end;
$$;

create constraint trigger bulk_balance_uom_conservation
after insert or update of quantity or delete on inventory.bulk_balance
deferrable initially deferred for each row execute function inventory.assert_bulk_uom_conservation();

create constraint trigger bulk_uom_balance_conservation
after insert or update or delete on inventory.bulk_uom_balance
deferrable initially deferred for each row execute function inventory.assert_bulk_uom_conservation();

create function inventory.reject_bulk_uom_conversion_mutation()
returns trigger language plpgsql as $$
begin
  raise exception using errcode='55000',message='bulk UOM conversion evidence is immutable';
end;
$$;

create trigger bulk_uom_conversion_immutable
before update or delete on inventory.bulk_uom_conversion
for each row execute function inventory.reject_bulk_uom_conversion_mutation();

create trigger bulk_uom_conversion_line_immutable
before update or delete on inventory.bulk_uom_conversion_line
for each row execute function inventory.reject_bulk_uom_conversion_mutation();

commit;
````

### FILE: `db/migrations/0034_bulk_uom_packaging.down.sql`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:bulk-uom-packaging-migration-down:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d breakbulk/gather reversible adaptation"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "649a416e9c8ebcfd029b0562842342a120a9ce69d81c743091c4499c616bca83"
variables: []
secrets_allowed: false
```

````sql
begin;

drop trigger if exists bulk_uom_conversion_line_immutable on inventory.bulk_uom_conversion_line;
drop trigger if exists bulk_uom_conversion_immutable on inventory.bulk_uom_conversion;
drop function if exists inventory.reject_bulk_uom_conversion_mutation();
drop trigger if exists bulk_uom_balance_conservation on inventory.bulk_uom_balance;
drop trigger if exists bulk_balance_uom_conservation on inventory.bulk_balance;
drop function if exists inventory.assert_bulk_uom_conservation();
drop trigger if exists bulk_balance_sync_base_uom on inventory.bulk_balance;
drop function if exists inventory.sync_bulk_base_uom();
drop trigger if exists bulk_uom_balance_validate on inventory.bulk_uom_balance;
drop function if exists inventory.validate_bulk_uom_balance();
drop table if exists inventory.bulk_uom_conversion_line;
drop table if exists inventory.bulk_uom_conversion;
drop table if exists inventory.bulk_uom_balance;
alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable;
delete from inventory.item_unit_of_measure u using inventory.bulk_uom_migration_backfill b
 where u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.uom_code=b.uom_code and u.is_base;
alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable;
drop table if exists inventory.bulk_uom_migration_backfill;

commit;
````

### FILE: `docs/inventory/MICROSOFT_BC_BREAKBULK_DERIVATION.md`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:bulk-uom-packaging-derivation:v1"
operation: CREATE
provenance: AUTHORED
source: "fixed Microsoft BCApps commit and Microsoft Learn authorities"
license: "LicenseRef-Workspace-Owner"
sha256: "4c45aaf05a2f3851256829f4960d2a4b10cc71aa7b4940a225c3dec2dad62dfc"
variables: []
secrets_allowed: false
```

````markdown
# Microsoft Business Central handling-unit composition derivation

This portable implementation is `ADAPTED`. Its Go, SQL and tests are not Microsoft-authored code and Microsoft does not support this adaptation. The governing source is the MIT-licensed `microsoft/BCApps` commit `2eae56d704a1fd035d104f333602aea7091b7749` (tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`) plus the linked Microsoft Learn product guidance.

## Fixed official inputs

- [Set up units of measure](https://learn.microsoft.com/en-gb/dynamics365/business-central/inventory-how-setup-units-of-measure)
- [Enable automatic breaking bulk](https://learn.microsoft.com/en-ca/dynamics365/business-central/warehouse-enable-automatic-breaking-bulk-with-directed-put-away-and-pick)
- [Warehouse management design details](https://learn.microsoft.com/en-ca/dynamics365/business-central/design-details-warehouse-management)

| BCApps path | Bytes | SHA-256 |
|---|---:|---|
| `src/Layers/W1/BaseApp/Inventory/Item/ItemUnitofMeasure.Table.al` | 20,689 | `ab785d08d4bbc9160c2e068b5660ada732f752d35ae7c94a6677ef0fb79b2009` |
| `src/Layers/W1/BaseApp/Warehouse/Activity/CreatePutaway.Codeunit.al` | 77,370 | `b004c0663d237585d8b57a06356431acb2302cdfbbd30f7b82154d282f325b2b` |
| `src/Layers/W1/BaseApp/Warehouse/Activity/WhseChangeUnitofMeasure.Report.al` | 8,008 | `d1fca83ae830a8b615c558e15cb6ec3734560614ba66e86d801f3e9794b06186` |
| `src/Layers/W1/BaseApp/Warehouse/Worksheet/Replenishment.Codeunit.al` | 19,064 | `e6d5e3873b4828db8dca49755670bd9815faa7ad1b0396189014af57380b7c4e` |
| `src/Layers/W1/BaseApp/Warehouse/Worksheet/WhseCreatePick.Codeunit.al` | 1,584 | `43ffaded138adb24d0cfc927c276d63ef3dfb2760f83f518bcbc8bbc217163bc` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMMovement.Codeunit.al` | 89,432 | `0e63996a0816d9379bdf76acc5712990f10624de81d875f2e499e6944b2955ed` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMWMSItemUnitofMeasure.Codeunit.al` | 125,924 | `8faec658fad625a12c83b251bab88553319bc1cf564e7671104f6a19cadf4713` |

## Official invariants transferred

The fixed replenishment owner searches the requested UOM first. Its breakbulk fallback is explicit, considers a larger package factor, retains base quantity as the authority and writes source and destination UOM/factor evidence. The fixed Microsoft tests aggregate Take and Place activity lines in base quantity, verify warehouse entries by UOM in base quantity, cover lot tracking, expose breakbulk lines during pick registration and also exercise gathering from a smaller UOM into a larger UOM.

The portable contract therefore admits only these narrower behaviors:

1. `inventory.bulk_balance.quantity` is the sole physical quantity in base UOM. Packaging rows are a composition, never a second inventory authority.
2. Every composition row carries an immutable UOM factor snapshot and exact base quantity. A deferred database invariant requires the sum of composition base quantities to equal the physical balance at commit.
3. `breakbulk` converts a larger factor to a smaller factor; `gather` converts a smaller factor to a larger factor. Both require exact source and target precision and preserve base quantity without rounding.
4. Each accepted command is serializable and idempotent by tenant/request. Divergent replay conflicts. Source composition is conditionally decremented so concurrent commands cannot consume the same package twice.
5. The immutable conversion header and its exactly two lines record Take from the source UOM and Place into the target UOM with equal base quantity. The outbox event commits in the same transaction.
6. Ordinary balance increases enter the base composition. Ordinary decreases consume only base composition. If stock is packaged only in a larger UOM, implicit unpacking is rejected and the caller must execute an explicit conversion first.
7. The upgrade migrates a legacy article lacking UOM evidence to factor `1` at the existing six-decimal physical storage precision, records the exact backfill set and reverses only that set on migration rollback.
8. A warehouse receipt accepts either an unambiguous base quantity or an unambiguous handling UOM/quantity pair. PostgreSQL resolves the immutable factor, verifies handling and base precision exactly, and keeps unit cost denominated in base UOM.
9. Put-away planning retains the received handling UOM only when every activity line is an exact multiple. Otherwise it fails closed unless the caller explicitly authorizes automatic breakbulk; the authorized transaction records equal-base Take/Place evidence linked to the posted receipt.
10. Every open put-away reserves both the singular physical balance and its matching composition. Registration releases/consumes both reservations, preserves the handling UOM across bins when possible, and commits movement entries in the same transaction. Cancellation releases both reservations without deleting receipt evidence or stock.
11. A manual conversion may consume only physical and composition quantities not reserved by an open activity. Transfer receipts use the same base-UOM reservation contract, so they do not bypass warehouse registration.

## Executable proof

PostgreSQL 18.6 from migrations `0001` through `0035` retains the V168 proof and adds an end-to-end V169 corpus. Two received BOX are preserved as one BOX on each of two 12-EA put-away lines; physical `24` and composition `2 BOX` are both fully reserved, manual conversion conflicts, registration leaves one unreserved BOX in each target and conservation remains exact. A one-BOX receipt split into two 6-EA targets rolls back without residue when authorization is absent; with authorization it commits two base-UOM lines, immutable Take/Place, receipt source binding and one outbox event. Cancelling an open one-BOX put-away leaves `12 EA = 1 BOX`, releases both reservations and permits the subsequent explicit conversion.

The full PostgreSQL suite also proves transfer-receipt put-away under the same reservation contract. Go tests, vet and build pass. Migration proof creates legacy receipt/activity rows under `0034`, applies `0035`, verifies exact `EA / 12 / factor 1` snapshots and zero packaging reservation, rolls back while preserving physical/receipt/activity quantities, reapplies and reproduces the same snapshots. A post-cycle V169 corpus passes again.

## Explicit non-claims

V169 proves automatic handling-UOM selection only for warehouse receipt → put-away → registration/cancellation and base-UOM transfer-receipt put-away. It does not claim packaging-aware pick or replenishment selection, sales/service/production demand parity, handheld/barcode execution, cubage/weight constraints, external WMS reconciliation or Business Central feature parity. Those lanes remain fail-closed until independently connected and tested. No target production readiness is inferred.
````

### FILE: `internal/inventorycontrol/warehouse_breakbulk.go`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:bulk-uom-packaging-domain:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d explicit breakbulk/gather contract adaptation"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "ae007c5e4883aec28b7b81c206a2c68258409771dba82e5ed20b5e2838c4940c"
variables: []
secrets_allowed: false
```

````go
package inventorycontrol

import (
	"context"
	"fmt"
)

type PackagingConversionCommand struct {
	RequestID      string `json:"request_id"`
	OrganizationID string `json:"organization_id"`
	BinID          string `json:"bin_id"`
	ItemID         string `json:"item_id"`
	LotID          string `json:"lot_id,omitempty"`
	Operation      string `json:"operation"`
	FromUOM        string `json:"from_uom"`
	ToUOM          string `json:"to_uom"`
	FromQuantity   string `json:"from_quantity"`
}

type PackagingConversionIDs struct {
	ConversionID string
	TakeLineID   string
	PlaceLineID  string
	EventID      string
}

type PackagingConversionLine struct {
	Sequence     int    `json:"sequence"`
	Action       string `json:"action"`
	UOM          string `json:"uom"`
	Quantity     string `json:"quantity"`
	BaseQuantity string `json:"base_quantity"`
}

type PackagingConversionResult struct {
	ID              string                    `json:"id"`
	RequestID       string                    `json:"request_id"`
	OrganizationID  string                    `json:"organization_id"`
	BinID           string                    `json:"bin_id"`
	ItemID          string                    `json:"item_id"`
	LotID           string                    `json:"lot_id,omitempty"`
	Operation       string                    `json:"operation"`
	FromUOM         string                    `json:"from_uom"`
	ToUOM           string                    `json:"to_uom"`
	FromQuantity    string                    `json:"from_quantity"`
	ToQuantity      string                    `json:"to_quantity"`
	BaseQuantity    string                    `json:"base_quantity"`
	FromQuantityPer string                    `json:"from_quantity_per_uom"`
	ToQuantityPer   string                    `json:"to_quantity_per_uom"`
	Version         int64                     `json:"version"`
	Lines           []PackagingConversionLine `json:"lines"`
}

func (s *BulkService) ConvertHandlingUnits(ctx context.Context, tenant string, value PackagingConversionCommand) (PackagingConversionResult, error) {
	operations := map[string]bool{"breakbulk": true, "gather": true}
	if tenant == "" || value.RequestID == "" || value.OrganizationID == "" || value.BinID == "" || value.ItemID == "" || !operations[value.Operation] || value.FromUOM == "" || value.ToUOM == "" || value.FromUOM == value.ToUOM || !positiveDecimal(value.FromQuantity, quantityPattern) {
		return PackagingConversionResult{}, fmt.Errorf("invalid packaging conversion")
	}
	ids := PackagingConversionIDs{ConversionID: s.ids.New(), TakeLineID: s.ids.New(), PlaceLineID: s.ids.New(), EventID: s.ids.New()}
	return s.repository.ConvertBulkHandlingUnits(ctx, tenant, ids, value)
}
````

### FILE: `internal/platform/postgres/warehouse_breakbulk.go`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:bulk-uom-packaging-postgres:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d base quantity and Take/Place invariant adaptation"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "f6b442a6c75249dfe1e2d79d98afa4a523a48445dd552691962973fa2b38aa77"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
)

func loadPackagingConversion(ctx context.Context, tx pgx.Tx, tenant, request string) (inventorycontrol.PackagingConversionResult, bool, error) {
	var value inventorycontrol.PackagingConversionResult
	err := tx.QueryRow(ctx, `select conversion_id,request_id,organization_id,bin_id,item_id,coalesce(lot_id,''),operation,from_uom,to_uom,from_quantity::text,to_quantity::text,base_quantity::text,from_qty_per_uom::text,to_qty_per_uom::text,version from inventory.bulk_uom_conversion where tenant_id=$1 and request_id=$2`, tenant, request).Scan(&value.ID, &value.RequestID, &value.OrganizationID, &value.BinID, &value.ItemID, &value.LotID, &value.Operation, &value.FromUOM, &value.ToUOM, &value.FromQuantity, &value.ToQuantity, &value.BaseQuantity, &value.FromQuantityPer, &value.ToQuantityPer, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, false, nil
	}
	if err != nil {
		return value, false, err
	}
	rows, err := tx.Query(ctx, `select sequence_no,action_type,uom_code,quantity::text,quantity_base::text from inventory.bulk_uom_conversion_line where tenant_id=$1 and conversion_id=$2 order by sequence_no`, tenant, value.ID)
	if err != nil {
		return value, false, err
	}
	defer rows.Close()
	for rows.Next() {
		var line inventorycontrol.PackagingConversionLine
		if err = rows.Scan(&line.Sequence, &line.Action, &line.UOM, &line.Quantity, &line.BaseQuantity); err != nil {
			return value, false, err
		}
		value.Lines = append(value.Lines, line)
	}
	return value, true, rows.Err()
}

func samePackagingRequest(value inventorycontrol.PackagingConversionResult, command inventorycontrol.PackagingConversionCommand) bool {
	return value.RequestID == command.RequestID && value.OrganizationID == command.OrganizationID && value.BinID == command.BinID && value.ItemID == command.ItemID && value.LotID == command.LotID && value.Operation == command.Operation && value.FromUOM == command.FromUOM && value.ToUOM == command.ToUOM
}

func (r *InventoryControl) ConvertBulkHandlingUnits(ctx context.Context, tenant string, ids inventorycontrol.PackagingConversionIDs, command inventorycontrol.PackagingConversionCommand) (inventorycontrol.PackagingConversionResult, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.PackagingConversionResult{}, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":"+command.RequestID); err != nil {
		return inventorycontrol.PackagingConversionResult{}, err
	}
	if existing, found, loadErr := loadPackagingConversion(ctx, tx, tenant, command.RequestID); loadErr != nil {
		return existing, loadErr
	} else if found {
		var quantityMatches bool
		if err = tx.QueryRow(ctx, `select from_quantity=$3::numeric from inventory.bulk_uom_conversion where tenant_id=$1 and request_id=$2`, tenant, command.RequestID, command.FromQuantity).Scan(&quantityMatches); err != nil {
			return existing, err
		}
		if !samePackagingRequest(existing, command) || !quantityMatches {
			return inventorycontrol.PackagingConversionResult{}, inventorycontrol.ErrConflict
		}
		return existing, tx.Commit(ctx)
	}

	var balanceID int64
	var fromFactor, toFactor, baseQuantity, toQuantity string
	err = tx.QueryRow(ctx, `select b.balance_id,fu.qty_per_uom::text,tu.qty_per_uom::text,($9::numeric*fu.qty_per_uom)::numeric(20,6)::text,(($9::numeric*fu.qty_per_uom)/tu.qty_per_uom)::numeric(20,6)::text
from inventory.bulk_balance b
join inventory.item_unit_of_measure fu on fu.tenant_id=b.tenant_id and fu.item_id=b.item_id and fu.uom_code=$7
join inventory.item_unit_of_measure tu on tu.tenant_id=b.tenant_id and tu.item_id=b.item_id and tu.uom_code=$8
where b.tenant_id=$1 and b.organization_id=$2 and b.bin_id=$3 and b.item_id=$4 and b.lot_id is not distinct from nullif($5,'')
  and (($6='breakbulk' and fu.qty_per_uom>tu.qty_per_uom) or ($6='gather' and fu.qty_per_uom<tu.qty_per_uom))
  and mod($9::numeric,fu.rounding_precision)=0
  and mod(($9::numeric*fu.qty_per_uom)/tu.qty_per_uom,tu.rounding_precision)=0
  and b.quantity-b.reserved_quantity >= $9::numeric*fu.qty_per_uom
for update of b`, tenant, command.OrganizationID, command.BinID, command.ItemID, command.LotID, command.Operation, command.FromUOM, command.ToUOM, command.FromQuantity).Scan(&balanceID, &fromFactor, &toFactor, &baseQuantity, &toQuantity)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.PackagingConversionResult{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.PackagingConversionResult{}, bulkConflict(err)
	}

	result, err := tx.Exec(ctx, `update inventory.bulk_uom_balance set quantity=quantity-$3::numeric,quantity_base=quantity_base-$4::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and quantity-reserved_quantity >= $3::numeric and quantity_base>=$4::numeric`, balanceID, command.FromUOM, command.FromQuantity, baseQuantity)
	if err != nil {
		return inventorycontrol.PackagingConversionResult{}, bulkConflict(err)
	}
	if result.RowsAffected() != 1 {
		return inventorycontrol.PackagingConversionResult{}, inventorycontrol.ErrConflict
	}
	if _, err = tx.Exec(ctx, `delete from inventory.bulk_uom_balance where balance_id=$1 and uom_code=$2 and quantity_base=0 and reserved_quantity=0`, balanceID, command.FromUOM); err != nil {
		return inventorycontrol.PackagingConversionResult{}, err
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base) values($1,$2,$3::numeric,$4::numeric,$5::numeric) on conflict(balance_id,uom_code) do update set quantity=inventory.bulk_uom_balance.quantity+excluded.quantity,quantity_base=inventory.bulk_uom_balance.quantity_base+excluded.quantity_base,version=inventory.bulk_uom_balance.version+1,updated_at=clock_timestamp()`, balanceID, command.ToUOM, toQuantity, toFactor, baseQuantity)
	if err != nil {
		return inventorycontrol.PackagingConversionResult{}, bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_conversion(tenant_id,conversion_id,request_id,organization_id,bin_id,item_id,lot_id,operation,from_uom,to_uom,from_quantity,to_quantity,base_quantity,from_qty_per_uom,to_qty_per_uom) values($1,$2,$3,$4,$5,$6,nullif($7,''),$8,$9,$10,$11::numeric,$12::numeric,$13::numeric,$14::numeric,$15::numeric)`, tenant, ids.ConversionID, command.RequestID, command.OrganizationID, command.BinID, command.ItemID, command.LotID, command.Operation, command.FromUOM, command.ToUOM, command.FromQuantity, toQuantity, baseQuantity, fromFactor, toFactor)
	if err != nil {
		return inventorycontrol.PackagingConversionResult{}, bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_conversion_line(tenant_id,conversion_id,line_id,sequence_no,action_type,uom_code,quantity,quantity_base) values($1,$2,$3,1,'take',$5,$6::numeric,$7::numeric),($1,$2,$4,2,'place',$8,$9::numeric,$7::numeric)`, tenant, ids.ConversionID, ids.TakeLineID, ids.PlaceLineID, command.FromUOM, command.FromQuantity, baseQuantity, command.ToUOM, toQuantity)
	if err != nil {
		return inventorycontrol.PackagingConversionResult{}, bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, ids.EventID, "bulk-uom-conversion", ids.ConversionID, "bulk-uom-conversion.completed", 1, map[string]string{"request_id": command.RequestID, "organization_id": command.OrganizationID, "bin_id": command.BinID, "item_id": command.ItemID, "lot_id": command.LotID, "operation": command.Operation, "from_uom": command.FromUOM, "to_uom": command.ToUOM, "from_quantity": command.FromQuantity, "to_quantity": toQuantity, "base_quantity": baseQuantity}); err != nil {
		return inventorycontrol.PackagingConversionResult{}, err
	}
	value, found, err := loadPackagingConversion(ctx, tx, tenant, command.RequestID)
	if err != nil || !found {
		return value, err
	}
	return value, tx.Commit(ctx)
}
````

### FILE: `internal/platform/postgres/warehouse_breakbulk_integration_test.go`

```yaml
block_id: "GO-SUPPLY-FACTORY-INVENTORY-API:bulk-uom-packaging-postgres-integration:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps@2eae56d SCM Movement and WMS Item UOM scenario adaptation"
license: "MIT AND LicenseRef-Workspace-Owner"
sha256: "cab6d7757dd3d08100bff50ade2a57c5b301d2e4fc935390b35b6282c18beffe"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestWarehouseHandlingUnitConservationIdempotencyAndConcurrency(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := bulkTestUUID(t)
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Packaging V168','Packaging V168')`, tenant, "pack-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	defer func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			t.Errorf("packaging cleanup begin: %v", cleanupErr)
			return
		}
		defer tx.Rollback(ctx)
		commands := []struct {
			sql  string
			args []any
		}{
			{`alter table inventory.bulk_uom_conversion_line disable trigger bulk_uom_conversion_line_immutable`, nil},
			{`alter table inventory.bulk_uom_conversion disable trigger bulk_uom_conversion_immutable`, nil},
			{`alter table inventory.bulk_inventory_entry disable trigger bulk_inventory_entry_immutable`, nil},
			{`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`, nil},
			{`delete from inventory.bulk_uom_conversion_line where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_uom_conversion where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_cost_application where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_cost_layer where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_inventory_entry where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_balance where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.item_bin_policy where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.warehouse_bin where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.item_unit_of_measure where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.stock_item where tenant_id=$1`, []any{tenant}},
			{`delete from platform.outbox_event where tenant_id=$1`, []any{tenant}},
			{`delete from org.organization where tenant_id=$1`, []any{tenant}},
			{`delete from platform.tenant where tenant_id=$1`, []any{tenant}},
			{`alter table inventory.bulk_uom_conversion_line enable trigger bulk_uom_conversion_line_immutable`, nil},
			{`alter table inventory.bulk_uom_conversion enable trigger bulk_uom_conversion_immutable`, nil},
			{`alter table inventory.bulk_inventory_entry enable trigger bulk_inventory_entry_immutable`, nil},
			{`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`, nil},
		}
		for _, command := range commands {
			if _, cleanupErr = tx.Exec(ctx, command.sql, command.args...); cleanupErr != nil {
				t.Errorf("packaging cleanup: %v", cleanupErr)
				return
			}
		}
		if cleanupErr = tx.Commit(ctx); cleanupErr != nil {
			t.Errorf("packaging cleanup commit: %v", cleanupErr)
		}
	}()
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'warehouse','warehouse','Warehouse','warehouse')`, tenant); err != nil {
		t.Fatal(err)
	}
	repo := NewInventoryControl(pool)
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "part", Code: "PACK-PART", Description: "Packaged part", BaseUOM: "EA", BaseRoundingPrecision: "1", TrackingMode: "none", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	for _, unit := range []inventorycontrol.ItemUnitOfMeasure{{ItemID: "part", Code: "BOX", QuantityPerUnit: "12", RoundingPrecision: "1", Version: 1}, {ItemID: "part", Code: "PALLET", QuantityPerUnit: "120", RoundingPrecision: "1", Version: 1}} {
		if _, err = repo.ConfigureItemUnitOfMeasure(ctx, tenant, bulkTestUUID(t), unit); err != nil {
			t.Fatal(err)
		}
	}
	for _, bin := range []inventorycontrol.WarehouseBin{{ID: "pick", OrganizationID: "warehouse", Code: "PICK", Type: "pick", Version: 1}, {ID: "target", OrganizationID: "warehouse", Code: "TARGET", Type: "putpick", Version: 1}} {
		if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatal(err)
		}
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: "warehouse", ItemID: "part", BinID: bin.ID, Fixed: true, Default: bin.ID == "pick", MinQuantity: "0", MaxQuantity: "1000", Version: 1}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err = repo.ReceiveBulk(ctx, tenant, bulkTestUUID(t), bulkTestUUID(t), bulkTestUUID(t), bulkTestUUID(t), inventorycontrol.BulkReceipt{OrganizationID: "warehouse", BinID: "pick", ItemID: "part", Quantity: "24", UnitCost: "1", PostingDate: time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC), SourceKind: "purchase-receipt", SourceID: "receipt-1"}); err != nil {
		t.Fatal(err)
	}
	ids := func() inventorycontrol.PackagingConversionIDs {
		return inventorycontrol.PackagingConversionIDs{ConversionID: bulkTestUUID(t), TakeLineID: bulkTestUUID(t), PlaceLineID: bulkTestUUID(t), EventID: bulkTestUUID(t)}
	}
	gather := inventorycontrol.PackagingConversionCommand{RequestID: "gather-1", OrganizationID: "warehouse", BinID: "pick", ItemID: "part", Operation: "gather", FromUOM: "EA", ToUOM: "BOX", FromQuantity: "24"}
	first, err := repo.ConvertBulkHandlingUnits(ctx, tenant, ids(), gather)
	if err != nil || first.ToQuantity != "2.000000" || first.BaseQuantity != "24.000000" || len(first.Lines) != 2 || first.Lines[0].Action != "take" || first.Lines[1].Action != "place" {
		t.Fatalf("gather=%+v err=%v", first, err)
	}
	replay, err := repo.ConvertBulkHandlingUnits(ctx, tenant, ids(), gather)
	if err != nil || replay.ID != first.ID {
		t.Fatalf("replay=%+v err=%v", replay, err)
	}
	divergent := gather
	divergent.FromQuantity = "12"
	if _, err = repo.ConvertBulkHandlingUnits(ctx, tenant, ids(), divergent); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("divergent replay error=%v", err)
	}
	breakOne := inventorycontrol.PackagingConversionCommand{RequestID: "break-1", OrganizationID: "warehouse", BinID: "pick", ItemID: "part", Operation: "breakbulk", FromUOM: "BOX", ToUOM: "EA", FromQuantity: "1"}
	if _, err = repo.ConvertBulkHandlingUnits(ctx, tenant, ids(), breakOne); err != nil {
		t.Fatal(err)
	}

	commands := []inventorycontrol.PackagingConversionCommand{breakOne, breakOne}
	commands[0].RequestID = "race-1"
	commands[1].RequestID = "race-2"
	errs := make([]error, 2)
	var wait sync.WaitGroup
	for index := range commands {
		wait.Add(1)
		go func(i int) {
			defer wait.Done()
			_, errs[i] = repo.ConvertBulkHandlingUnits(ctx, tenant, ids(), commands[i])
		}(index)
	}
	wait.Wait()
	successes, conflicts := 0, 0
	for _, raceErr := range errs {
		if raceErr == nil {
			successes++
		} else if errors.Is(raceErr, inventorycontrol.ErrConflict) {
			conflicts++
		} else {
			t.Fatalf("unexpected race error=%v", raceErr)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("race successes=%d conflicts=%d errors=%v", successes, conflicts, errs)
	}
	var physical, composed string
	var conversions, lines, outbox int
	if err = pool.QueryRow(ctx, `select b.quantity::text,coalesce(sum(u.quantity_base),0)::text,(select count(*) from inventory.bulk_uom_conversion where tenant_id=$1),(select count(*) from inventory.bulk_uom_conversion_line where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='bulk-uom-conversion.completed') from inventory.bulk_balance b left join inventory.bulk_uom_balance u using(balance_id) where b.tenant_id=$1 and b.organization_id='warehouse' and b.bin_id='pick' and b.item_id='part' group by b.balance_id`, tenant).Scan(&physical, &composed, &conversions, &lines, &outbox); err != nil || physical != "24.000000" || composed != "24.000000" || conversions != 3 || lines != 6 || outbox != 3 {
		t.Fatalf("physical=%s composed=%s conversions=%d lines=%d outbox=%d err=%v", physical, composed, conversions, lines, outbox, err)
	}
	if _, err = pool.Exec(ctx, `update inventory.bulk_uom_conversion set operation='gather' where tenant_id=$1 and conversion_id=$2`, tenant, first.ID); err == nil {
		t.Fatal("immutable conversion mutated")
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base) select balance_id,'BOX',1,12,12 from inventory.bulk_balance where tenant_id=$1 and bin_id='pick' and item_id='part'`, tenant)
	if err != nil {
		tx.Rollback(ctx)
		t.Fatal(err)
	}
	err = tx.Commit(ctx)
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "55000" {
		t.Fatalf("conservation commit error=%v", err)
	}

	if _, err = repo.ConvertBulkHandlingUnits(ctx, tenant, ids(), inventorycontrol.PackagingConversionCommand{RequestID: "gather-2", OrganizationID: "warehouse", BinID: "pick", ItemID: "part", Operation: "gather", FromUOM: "EA", ToUOM: "BOX", FromQuantity: "24"}); err != nil {
		t.Fatal(err)
	}
	err = repo.MoveBulk(ctx, tenant, bulkTestUUID(t), bulkTestUUID(t), inventorycontrol.BulkMovement{OrganizationID: "warehouse", FromBinID: "pick", ToBinID: "target", ItemID: "part", Quantity: "1", PostingDate: time.Date(2030, 1, 2, 0, 0, 0, 0, time.UTC), SourceID: "move-packaged"})
	if !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("implicit breakbulk error=%v", err)
	}

	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "weight", Code: "WEIGHT", Description: "Fractional material", BaseUOM: "KG", BaseRoundingPrecision: "0.001", TrackingMode: "none", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	for _, bin := range []string{"pick", "target"} {
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: "warehouse", ItemID: "weight", BinID: bin, Fixed: true, Default: bin == "pick", MinQuantity: "0", MaxQuantity: "1000", Version: 1}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err = repo.ReceiveBulk(ctx, tenant, bulkTestUUID(t), bulkTestUUID(t), bulkTestUUID(t), bulkTestUUID(t), inventorycontrol.BulkReceipt{OrganizationID: "warehouse", BinID: "pick", ItemID: "weight", Quantity: "0.125", UnitCost: "1", PostingDate: time.Date(2030, 1, 3, 0, 0, 0, 0, time.UTC), SourceKind: "purchase-receipt", SourceID: "weight-1"}); err != nil {
		t.Fatal(err)
	}
	if err = repo.MoveBulk(ctx, tenant, bulkTestUUID(t), bulkTestUUID(t), inventorycontrol.BulkMovement{OrganizationID: "warehouse", FromBinID: "pick", ToBinID: "target", ItemID: "weight", Quantity: "0.025", PostingDate: time.Date(2030, 1, 4, 0, 0, 0, 0, time.UTC), SourceID: "weight-move"}); err != nil {
		t.Fatalf("fractional base move: %v", err)
	}
}
````

### FILE: `db/migrations/0035_warehouse_handling_uom.up.sql`

```yaml
block_id: "GO-SUPPLY-WAREHOUSE-HANDLING-UOM-MIGRATION-UP:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps 2eae56d704a1fd035d104f333602aea7091b7749 invariants plus local verified adaptation"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "f76e6c75084e0631149ba8533404d7e747161470ef879a2b4ac178a4b7398fa5"
variables: []
secrets_allowed: false
```

````sql
begin;

alter table inventory.bulk_uom_balance
  add column reserved_quantity numeric(20,6) not null default 0,
  add constraint bulk_uom_balance_reserved_quantity_check
    check (reserved_quantity >= 0 and reserved_quantity <= quantity);

alter table inventory.bulk_uom_conversion
  add column source_kind text not null default 'manual',
  add column source_id text not null default '',
  add constraint bulk_uom_conversion_source_kind_check
    check (length(source_kind) between 1 and 64),
  add constraint bulk_uom_conversion_source_id_check
    check (length(source_id) <= 128);

alter table inventory.warehouse_receipt
  add column allow_breakbulk boolean not null default false;

alter table inventory.warehouse_receipt_line
  add column uom_code text,
  add column uom_quantity numeric(20,6),
  add column qty_per_uom numeric(20,6);

alter table inventory.warehouse_activity_line
  add column uom_code text,
  add column uom_quantity numeric(20,6),
  add column qty_per_uom numeric(20,6);

alter table inventory.warehouse_receipt_line disable trigger warehouse_receipt_line_immutable;
update inventory.warehouse_receipt_line l
   set uom_code=u.uom_code,
       uom_quantity=l.quantity/u.qty_per_uom,
       qty_per_uom=u.qty_per_uom
  from inventory.item_unit_of_measure u
 where u.tenant_id=l.tenant_id and u.item_id=l.item_id and u.is_base;
alter table inventory.warehouse_receipt_line enable trigger warehouse_receipt_line_immutable;

alter table inventory.warehouse_activity_line disable trigger warehouse_activity_line_immutable;
update inventory.warehouse_activity_line l
   set uom_code=u.uom_code,
       uom_quantity=l.quantity/u.qty_per_uom,
       qty_per_uom=u.qty_per_uom
  from inventory.item_unit_of_measure u
 where u.tenant_id=l.tenant_id and u.item_id=l.item_id and u.is_base;
alter table inventory.warehouse_activity_line enable trigger warehouse_activity_line_immutable;

alter table inventory.warehouse_receipt_line
  alter column uom_code set not null,
  alter column uom_quantity set not null,
  alter column qty_per_uom set not null,
  add constraint warehouse_receipt_line_uom_quantity_check check (uom_quantity > 0),
  add constraint warehouse_receipt_line_qty_per_uom_check check (qty_per_uom > 0),
  add constraint warehouse_receipt_line_uom_exact_check check (uom_quantity * qty_per_uom = quantity),
  add constraint warehouse_receipt_line_uom_fk foreign key (tenant_id,item_id,uom_code)
    references inventory.item_unit_of_measure(tenant_id,item_id,uom_code);

alter table inventory.warehouse_activity_line
  alter column uom_code set not null,
  alter column uom_quantity set not null,
  alter column qty_per_uom set not null,
  add constraint warehouse_activity_line_uom_quantity_check check (uom_quantity > 0),
  add constraint warehouse_activity_line_qty_per_uom_check check (qty_per_uom > 0),
  add constraint warehouse_activity_line_uom_exact_check check (uom_quantity * qty_per_uom = quantity),
  add constraint warehouse_activity_line_uom_fk foreign key (tenant_id,item_id,uom_code)
    references inventory.item_unit_of_measure(tenant_id,item_id,uom_code);

create function inventory.validate_warehouse_line_uom()
returns trigger language plpgsql as $function$
declare
  expected_factor numeric(20,6);
  expected_precision numeric(20,6);
begin
  select u.qty_per_uom,u.rounding_precision
    into expected_factor,expected_precision
    from inventory.item_unit_of_measure u
   where u.tenant_id=new.tenant_id and u.item_id=new.item_id and u.uom_code=new.uom_code;
  if not found or new.qty_per_uom<>expected_factor or mod(new.uom_quantity,expected_precision)<>0 then
    raise exception using errcode='55000',message='warehouse line UOM does not match the immutable item UOM contract';
  end if;
  return new;
end;
$function$;

create trigger warehouse_receipt_line_uom_validate
before insert or update on inventory.warehouse_receipt_line
for each row execute function inventory.validate_warehouse_line_uom();

create trigger warehouse_activity_line_uom_validate
before insert or update on inventory.warehouse_activity_line
for each row execute function inventory.validate_warehouse_line_uom();

commit;
````

### FILE: `db/migrations/0035_warehouse_handling_uom.down.sql`

```yaml
block_id: "GO-SUPPLY-WAREHOUSE-HANDLING-UOM-MIGRATION-DOWN:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps 2eae56d704a1fd035d104f333602aea7091b7749 invariants plus local verified adaptation"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "725c06b79a4851577406a28604ef42caaf4b442c749e74efe3c60cc4853cae45"
variables: []
secrets_allowed: false
```

````sql
begin;

drop trigger if exists warehouse_activity_line_uom_validate on inventory.warehouse_activity_line;
drop trigger if exists warehouse_receipt_line_uom_validate on inventory.warehouse_receipt_line;
drop function if exists inventory.validate_warehouse_line_uom();

alter table inventory.warehouse_activity_line
  drop constraint if exists warehouse_activity_line_uom_fk,
  drop constraint if exists warehouse_activity_line_uom_exact_check,
  drop constraint if exists warehouse_activity_line_qty_per_uom_check,
  drop constraint if exists warehouse_activity_line_uom_quantity_check,
  drop column if exists qty_per_uom,
  drop column if exists uom_quantity,
  drop column if exists uom_code;

alter table inventory.warehouse_receipt_line
  drop constraint if exists warehouse_receipt_line_uom_fk,
  drop constraint if exists warehouse_receipt_line_uom_exact_check,
  drop constraint if exists warehouse_receipt_line_qty_per_uom_check,
  drop constraint if exists warehouse_receipt_line_uom_quantity_check,
  drop column if exists qty_per_uom,
  drop column if exists uom_quantity,
  drop column if exists uom_code;

alter table inventory.warehouse_receipt drop column if exists allow_breakbulk;

alter table inventory.bulk_uom_conversion
  drop constraint if exists bulk_uom_conversion_source_id_check,
  drop constraint if exists bulk_uom_conversion_source_kind_check,
  drop column if exists source_id,
  drop column if exists source_kind;

alter table inventory.bulk_uom_balance
  drop constraint if exists bulk_uom_balance_reserved_quantity_check,
  drop column if exists reserved_quantity;

commit;
````

### FILE: `internal/platform/postgres/warehouse_packaging_flow_integration_test.go`

```yaml
block_id: "GO-SUPPLY-WAREHOUSE-PACKAGING-FLOW-INTEGRATION:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps 2eae56d704a1fd035d104f333602aea7091b7749 invariants plus local verified adaptation"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "dfe2ad1f7dca87344fd3f07710386a3c96f67865168e02650e7df775bb028e1a"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5/pgxpool"
)

func warehousePackagingIDs(t *testing.T, prefix string) inventorycontrol.WarehouseIDs {
	t.Helper()
	return inventorycontrol.WarehouseIDs{
		ActivityID:        prefix + "-activity",
		ReceiptID:         prefix + "-receipt",
		LineID:            prefix + "-line",
		EntryID:           prefix + "-entry",
		LayerID:           prefix + "-layer",
		LotID:             prefix + "-lot",
		EventID:           bulkTestUUID(t),
		ConversionID:      prefix + "-conversion",
		TakeLineID:        prefix + "-take",
		PlaceLineID:       prefix + "-place",
		ConversionEventID: bulkTestUUID(t),
	}
}

func exactDecimalEqual(left, right string) bool {
	leftValue, leftOK := parseRat(left)
	rightValue, rightOK := parseRat(right)
	return leftOK && rightOK && leftValue.Cmp(rightValue) == 0
}

func TestWarehouseReceiptPackagingPreservationBreakbulkAndCancellation(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := bulkTestUUID(t)
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Warehouse Packaging V169','Warehouse Packaging V169')`, tenant, "whp-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'warehouse','warehouse','Warehouse','warehouse')`, tenant); err != nil {
		t.Fatal(err)
	}
	defer func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			t.Errorf("warehouse packaging cleanup begin: %v", cleanupErr)
			return
		}
		defer tx.Rollback(ctx)
		commands := []string{
			`alter table inventory.warehouse_pick_request disable trigger warehouse_pick_request_immutable`,
			`alter table inventory.warehouse_activity_uom_reservation disable trigger warehouse_uom_reservation_delete_guard`,
			`alter table inventory.bulk_uom_conversion_line disable trigger bulk_uom_conversion_line_immutable`,
			`alter table inventory.bulk_uom_conversion disable trigger bulk_uom_conversion_immutable`,
			`alter table inventory.warehouse_activity_line disable trigger warehouse_activity_line_immutable`,
			`alter table inventory.warehouse_receipt_line disable trigger warehouse_receipt_line_immutable`,
			`alter table inventory.warehouse_receipt disable trigger warehouse_receipt_immutable`,
			`alter table inventory.bulk_inventory_entry disable trigger bulk_inventory_entry_immutable`,
			`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`,
			`delete from inventory.warehouse_replenishment_request where tenant_id=$1`,
			`delete from inventory.warehouse_activity_uom_reservation where tenant_id=$1`,
			`delete from inventory.bulk_uom_conversion_line where tenant_id=$1`,
			`delete from inventory.bulk_uom_conversion where tenant_id=$1`,
			`delete from inventory.warehouse_pick_request where tenant_id=$1`,
			`delete from inventory.warehouse_activity_line where tenant_id=$1`,
			`delete from inventory.warehouse_activity where tenant_id=$1`,
			`delete from inventory.warehouse_receipt_line where tenant_id=$1`,
			`delete from inventory.warehouse_receipt where tenant_id=$1`,
			`delete from inventory.bulk_cost_application where tenant_id=$1`,
			`delete from inventory.bulk_cost_layer where tenant_id=$1`,
			`delete from inventory.bulk_inventory_entry where tenant_id=$1`,
			`delete from inventory.bulk_reservation where tenant_id=$1`,
			`delete from inventory.bulk_balance where tenant_id=$1`,
			`delete from inventory.item_bin_policy where tenant_id=$1`,
			`delete from inventory.warehouse_bin where tenant_id=$1`,
			`delete from inventory.item_unit_of_measure where tenant_id=$1`,
			`delete from inventory.stock_item where tenant_id=$1`,
			`delete from platform.outbox_event where tenant_id=$1`,
			`delete from org.organization where tenant_id=$1`,
			`delete from platform.tenant where tenant_id=$1`,
			`alter table inventory.bulk_uom_conversion_line enable trigger bulk_uom_conversion_line_immutable`,
			`alter table inventory.bulk_uom_conversion enable trigger bulk_uom_conversion_immutable`,
			`alter table inventory.warehouse_activity_line enable trigger warehouse_activity_line_immutable`,
			`alter table inventory.warehouse_receipt_line enable trigger warehouse_receipt_line_immutable`,
			`alter table inventory.warehouse_receipt enable trigger warehouse_receipt_immutable`,
			`alter table inventory.bulk_inventory_entry enable trigger bulk_inventory_entry_immutable`,
			`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`,
			`alter table inventory.warehouse_activity_uom_reservation enable trigger warehouse_uom_reservation_delete_guard`,
			`alter table inventory.warehouse_pick_request enable trigger warehouse_pick_request_immutable`,
		}
		for _, command := range commands {
			arguments := []any{}
			if strings.Contains(command, "$1") {
				arguments = append(arguments, tenant)
			}
			if _, cleanupErr = tx.Exec(ctx, command, arguments...); cleanupErr != nil {
				t.Errorf("warehouse packaging cleanup: %v", cleanupErr)
				return
			}
		}
		if cleanupErr = tx.Commit(ctx); cleanupErr != nil {
			t.Errorf("warehouse packaging cleanup commit: %v", cleanupErr)
		}
	}()

	repo := NewInventoryControl(pool)
	for _, item := range []string{"preserve", "split", "cancel", "pick-split", "replenish-split"} {
		if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: item, Code: "PACK-" + item, Description: item, BaseUOM: "EA", BaseRoundingPrecision: "1", TrackingMode: "none", CostingMethod: "fifo", Version: 1}); err != nil {
			t.Fatal(err)
		}
		if _, err = repo.ConfigureItemUnitOfMeasure(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemUnitOfMeasure{ItemID: item, Code: "BOX", QuantityPerUnit: "12", RoundingPrecision: "1", Version: 1}); err != nil {
			t.Fatal(err)
		}
	}
	bins := []inventorycontrol.WarehouseBin{
		{ID: "receive", OrganizationID: "warehouse", Code: "RECEIVE", Type: "receive", Version: 1},
		{ID: "target-a", OrganizationID: "warehouse", Code: "TARGET-A", Type: "putpick", Ranking: 20, Version: 1},
		{ID: "target-b", OrganizationID: "warehouse", Code: "TARGET-B", Type: "putpick", Ranking: 10, Version: 1},
		{ID: "half-a", OrganizationID: "warehouse", Code: "HALF-A", Type: "put-away", Ranking: 20, Version: 1},
		{ID: "half-b", OrganizationID: "warehouse", Code: "HALF-B", Type: "put-away", Ranking: 10, Version: 1},
		{ID: "target-c", OrganizationID: "warehouse", Code: "TARGET-C", Type: "put-away", Ranking: 20, Version: 1},
		{ID: "pick-source", OrganizationID: "warehouse", Code: "PICK-SOURCE", Type: "putpick", Ranking: 20, Version: 1},
		{ID: "replenish-source", OrganizationID: "warehouse", Code: "REPLENISH-SOURCE", Type: "putpick", Ranking: 30, Version: 1},
		{ID: "replenish-target", OrganizationID: "warehouse", Code: "REPLENISH-TARGET", Type: "pick", Ranking: 100, Version: 1},
		{ID: "ship", OrganizationID: "warehouse", Code: "SHIP", Type: "ship", Version: 1},
	}
	for _, bin := range bins {
		if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatal(err)
		}
	}
	policies := []struct{ item, bin, maximum string }{
		{"preserve", "receive", "100"}, {"preserve", "target-a", "12"}, {"preserve", "target-b", "12"}, {"preserve", "ship", "100"},
		{"split", "receive", "100"}, {"split", "half-a", "6"}, {"split", "half-b", "6"},
		{"cancel", "receive", "100"}, {"cancel", "target-c", "12"},
		{"pick-split", "receive", "100"}, {"pick-split", "pick-source", "12"}, {"pick-split", "ship", "100"},
		{"replenish-split", "receive", "100"}, {"replenish-split", "replenish-source", "12"}, {"replenish-split", "replenish-target", "6"},
	}
	for _, policy := range policies {
		isDefault := policy.bin == "target-a" || policy.bin == "half-a" || policy.bin == "target-c" || policy.bin == "pick-source" || policy.bin == "replenish-source"
		minimum := "0"
		if policy.bin == "replenish-target" {
			minimum = "5"
		}
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: "warehouse", ItemID: policy.item, BinID: policy.bin, Fixed: true, Default: isDefault, MinQuantity: minimum, MaxQuantity: policy.maximum, Version: 1}); err != nil {
			t.Fatal(err)
		}
	}
	posting := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	preserve, err := repo.PostWarehouseReceipt(ctx, tenant, warehousePackagingIDs(t, "preserve"), inventorycontrol.WarehouseReceiptCommand{RequestID: "preserve-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "preserve", HandlingUOM: "BOX", HandlingQuantity: "2", UnitCost: "10", PostingDate: posting, SourceKind: "purchase", SourceID: "preserve-po"})
	if err != nil || !exactDecimalEqual(preserve.Quantity, "24") || preserve.HandlingUOM != "BOX" || len(preserve.PutAway.Lines) != 2 {
		t.Fatalf("preserve receipt=%+v err=%v", preserve, err)
	}
	for _, line := range preserve.PutAway.Lines {
		if !exactDecimalEqual(line.Quantity, "12") || line.UOMCode != "BOX" || !exactDecimalEqual(line.UOMQuantity, "1") || !exactDecimalEqual(line.QuantityPerUOM, "12") {
			t.Fatalf("preserved line=%+v", line)
		}
	}
	var physical, physicalReserved, composed, composedReserved string
	if err = pool.QueryRow(ctx, `select b.quantity::text,b.reserved_quantity::text,u.quantity_base::text,u.reserved_quantity::text from inventory.bulk_balance b join inventory.bulk_uom_balance u using(balance_id) where b.tenant_id=$1 and b.item_id='preserve' and b.bin_id='receive' and u.uom_code='BOX'`, tenant).Scan(&physical, &physicalReserved, &composed, &composedReserved); err != nil || physical != "24.000000" || physicalReserved != "24.000000" || composed != "24.000000" || composedReserved != "2.000000" {
		t.Fatalf("preserve reserve physical=%s/%s composed=%s/%s err=%v", physical, physicalReserved, composed, composedReserved, err)
	}
	_, err = repo.ConvertBulkHandlingUnits(ctx, tenant, inventorycontrol.PackagingConversionIDs{ConversionID: "reserved-conversion", TakeLineID: "reserved-take", PlaceLineID: "reserved-place", EventID: bulkTestUUID(t)}, inventorycontrol.PackagingConversionCommand{RequestID: "reserved-conversion-request", OrganizationID: "warehouse", BinID: "receive", ItemID: "preserve", Operation: "breakbulk", FromUOM: "BOX", ToUOM: "EA", FromQuantity: "1"})
	if !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("reserved package conversion error=%v", err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", preserve.PutAway.ID, 1, posting.AddDate(0, 0, 1), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	var targetPackages, targetReserved string
	if err = pool.QueryRow(ctx, `select sum(u.quantity)::text,sum(u.reserved_quantity)::text from inventory.bulk_balance b join inventory.bulk_uom_balance u using(balance_id) where b.tenant_id=$1 and b.item_id='preserve' and b.bin_id in ('target-a','target-b') and u.uom_code='BOX'`, tenant).Scan(&targetPackages, &targetReserved); err != nil || targetPackages != "2.000000" || targetReserved != "0.000000" {
		t.Fatalf("target packages=%s reserved=%s err=%v", targetPackages, targetReserved, err)
	}

	splitCommand := inventorycontrol.WarehouseReceiptCommand{RequestID: "split-rejected-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "split", HandlingUOM: "BOX", HandlingQuantity: "1", UnitCost: "10", PostingDate: posting, SourceKind: "purchase", SourceID: "split-rejected-po"}
	if _, err = repo.PostWarehouseReceipt(ctx, tenant, warehousePackagingIDs(t, "split-rejected"), splitCommand); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("breakbulk without authorization error=%v", err)
	}
	var rejectedResidue int
	if err = pool.QueryRow(ctx, `select count(*) from inventory.warehouse_receipt where tenant_id=$1 and receipt_id='split-rejected-receipt'`, tenant).Scan(&rejectedResidue); err != nil || rejectedResidue != 0 {
		t.Fatalf("rejected receipt residue=%d err=%v", rejectedResidue, err)
	}
	splitCommand.RequestID, splitCommand.SourceID, splitCommand.AllowBreakbulk = "split-allowed-request", "split-allowed-po", true
	split, err := repo.PostWarehouseReceipt(ctx, tenant, warehousePackagingIDs(t, "split-allowed"), splitCommand)
	if err != nil || len(split.PutAway.Lines) != 2 {
		t.Fatalf("authorized split=%+v err=%v", split, err)
	}
	for _, line := range split.PutAway.Lines {
		if !exactDecimalEqual(line.Quantity, "6") || line.UOMCode != "EA" || !exactDecimalEqual(line.UOMQuantity, "6") {
			t.Fatalf("breakbulk line=%+v", line)
		}
	}
	var sourceKind, sourceID string
	var conversionLines, conversionEvents int
	if err = pool.QueryRow(ctx, `select c.source_kind,c.source_id,(select count(*) from inventory.bulk_uom_conversion_line l where l.tenant_id=c.tenant_id and l.conversion_id=c.conversion_id),(select count(*) from platform.outbox_event e where e.tenant_id=c.tenant_id and e.aggregate_id=c.conversion_id and e.event_type='bulk-uom-conversion.completed') from inventory.bulk_uom_conversion c where c.tenant_id=$1 and c.conversion_id='split-allowed-conversion'`, tenant).Scan(&sourceKind, &sourceID, &conversionLines, &conversionEvents); err != nil || sourceKind != "warehouse-receipt" || sourceID != "split-allowed-receipt" || conversionLines != 2 || conversionEvents != 1 {
		t.Fatalf("automatic conversion source=%s/%s lines=%d events=%d err=%v", sourceKind, sourceID, conversionLines, conversionEvents, err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", split.PutAway.ID, 1, posting.AddDate(0, 0, 2), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}

	cancelReceipt, err := repo.PostWarehouseReceipt(ctx, tenant, warehousePackagingIDs(t, "cancel"), inventorycontrol.WarehouseReceiptCommand{RequestID: "cancel-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "cancel", HandlingUOM: "BOX", HandlingQuantity: "1", UnitCost: "10", PostingDate: posting, SourceKind: "purchase", SourceID: "cancel-po"})
	if err != nil {
		t.Fatal(err)
	}
	cancelled, err := repo.CancelWarehousePutAway(ctx, tenant, "warehouse", cancelReceipt.PutAway.ID, 1, bulkTestUUID(t))
	if err != nil || cancelled.Status != "cancelled" || cancelled.Version != 2 {
		t.Fatalf("cancelled put-away=%+v err=%v", cancelled, err)
	}
	if err = pool.QueryRow(ctx, `select b.quantity::text,b.reserved_quantity::text,u.quantity::text,u.reserved_quantity::text from inventory.bulk_balance b join inventory.bulk_uom_balance u using(balance_id) where b.tenant_id=$1 and b.item_id='cancel' and b.bin_id='receive' and u.uom_code='BOX'`, tenant).Scan(&physical, &physicalReserved, &composed, &composedReserved); err != nil || physical != "12.000000" || physicalReserved != "0.000000" || composed != "1.000000" || composedReserved != "0.000000" {
		t.Fatalf("cancel release physical=%s/%s composed=%s/%s err=%v", physical, physicalReserved, composed, composedReserved, err)
	}
	if _, err = repo.ConvertBulkHandlingUnits(ctx, tenant, inventorycontrol.PackagingConversionIDs{ConversionID: "cancel-conversion", TakeLineID: "cancel-take", PlaceLineID: "cancel-place", EventID: bulkTestUUID(t)}, inventorycontrol.PackagingConversionCommand{RequestID: "cancel-conversion-request", OrganizationID: "warehouse", BinID: "receive", ItemID: "cancel", Operation: "breakbulk", FromUOM: "BOX", ToUOM: "EA", FromQuantity: "1"}); err != nil {
		t.Fatalf("conversion after cancellation: %v", err)
	}

	pickReceipt, err := repo.PostWarehouseReceipt(ctx, tenant, warehousePackagingIDs(t, "pick-split"), inventorycontrol.WarehouseReceiptCommand{RequestID: "pick-split-receipt-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "pick-split", HandlingUOM: "BOX", HandlingQuantity: "1", UnitCost: "10", PostingDate: posting, SourceKind: "purchase", SourceID: "pick-split-po"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", pickReceipt.PutAway.ID, 1, posting.AddDate(0, 0, 3), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	pickCommand := inventorycontrol.WarehousePickCommand{RequestID: "pick-no-breakbulk", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "pick-split", DemandKind: "manual", DemandID: "order-pack", DemandLineID: "line-pack", Quantity: "6", HandlingUOM: "EA"}
	if _, err = repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "pick-no-breakbulk"), pickCommand); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("partial package pick without authorization error=%v", err)
	}
	pickCommand.RequestID, pickCommand.AllowBreakbulk = "pick-cancel-request", true
	pickToCancel, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "pick-cancel"), pickCommand)
	if err != nil || len(pickToCancel.Lines) != 1 || pickToCancel.Lines[0].FromUOMCode != "BOX" || !exactDecimalEqual(pickToCancel.Lines[0].FromUOMQuantity, "1") || pickToCancel.Lines[0].UOMCode != "EA" || !exactDecimalEqual(pickToCancel.Lines[0].UOMQuantity, "6") {
		t.Fatalf("packaging-aware pick=%+v err=%v", pickToCancel, err)
	}
	if err = pool.QueryRow(ctx, `select b.reserved_quantity::text,u.reserved_quantity::text from inventory.bulk_balance b join inventory.bulk_uom_balance u using(balance_id) where b.tenant_id=$1 and b.item_id='pick-split' and b.bin_id='pick-source' and u.uom_code='BOX'`, tenant).Scan(&physicalReserved, &composedReserved); err != nil || physicalReserved != "6.000000" || composedReserved != "1.000000" {
		t.Fatalf("pick reservations physical=%s package=%s err=%v", physicalReserved, composedReserved, err)
	}
	if _, err = repo.ConvertBulkHandlingUnits(ctx, tenant, inventorycontrol.PackagingConversionIDs{ConversionID: "pick-reserved-conversion", TakeLineID: "pick-reserved-take", PlaceLineID: "pick-reserved-place", EventID: bulkTestUUID(t)}, inventorycontrol.PackagingConversionCommand{RequestID: "pick-reserved-conversion-request", OrganizationID: "warehouse", BinID: "pick-source", ItemID: "pick-split", Operation: "breakbulk", FromUOM: "BOX", ToUOM: "EA", FromQuantity: "1"}); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("reserved pick package conversion error=%v", err)
	}
	if _, err = repo.CancelWarehousePick(ctx, tenant, "warehouse", pickToCancel.ID, 1, bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select b.reserved_quantity::text,u.reserved_quantity::text from inventory.bulk_balance b join inventory.bulk_uom_balance u using(balance_id) where b.tenant_id=$1 and b.item_id='pick-split' and b.bin_id='pick-source' and u.uom_code='BOX'`, tenant).Scan(&physicalReserved, &composedReserved); err != nil || physicalReserved != "0.000000" || composedReserved != "0.000000" {
		t.Fatalf("cancelled pick reservations physical=%s package=%s err=%v", physicalReserved, composedReserved, err)
	}
	pickCommand.RequestID = "pick-register-request"
	pick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "pick-register"), pickCommand)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", pick.ID, 1, posting.AddDate(0, 0, 4), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	var pickSourceBase, pickSourceEA, pickShipBase, pickShipEA string
	if err = pool.QueryRow(ctx, `select (select quantity::text from inventory.bulk_balance where tenant_id=$1 and item_id='pick-split' and bin_id='pick-source'),(select u.quantity::text from inventory.bulk_uom_balance u join inventory.bulk_balance b using(balance_id) where b.tenant_id=$1 and b.item_id='pick-split' and b.bin_id='pick-source' and u.uom_code='EA'),(select quantity::text from inventory.bulk_balance where tenant_id=$1 and item_id='pick-split' and bin_id='ship'),(select u.quantity::text from inventory.bulk_uom_balance u join inventory.bulk_balance b using(balance_id) where b.tenant_id=$1 and b.item_id='pick-split' and b.bin_id='ship' and u.uom_code='EA')`, tenant).Scan(&pickSourceBase, &pickSourceEA, &pickShipBase, &pickShipEA); err != nil || pickSourceBase != "6.000000" || pickSourceEA != "6.000000" || pickShipBase != "6.000000" || pickShipEA != "6.000000" {
		t.Fatalf("registered pick source=%s/%s ship=%s/%s err=%v", pickSourceBase, pickSourceEA, pickShipBase, pickShipEA, err)
	}
	if err = pool.QueryRow(ctx, `select source_kind,source_id from inventory.bulk_uom_conversion where tenant_id=$1 and source_kind='warehouse-pick' and source_id=$2`, tenant, pick.ID).Scan(&sourceKind, &sourceID); err != nil || sourceKind != "warehouse-pick" || sourceID != pick.ID {
		t.Fatalf("pick conversion source=%s/%s err=%v", sourceKind, sourceID, err)
	}

	replenishReceipt, err := repo.PostWarehouseReceipt(ctx, tenant, warehousePackagingIDs(t, "replenish-split"), inventorycontrol.WarehouseReceiptCommand{RequestID: "replenish-split-receipt-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "replenish-split", HandlingUOM: "BOX", HandlingQuantity: "1", UnitCost: "10", PostingDate: posting, SourceKind: "purchase", SourceID: "replenish-split-po"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", replenishReceipt.PutAway.ID, 1, posting.AddDate(0, 0, 5), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	replenishCommand := inventorycontrol.WarehouseReplenishmentCommand{RequestID: "replenish-no-breakbulk", OrganizationID: "warehouse", ToBinID: "replenish-target", ItemID: "replenish-split", TargetUOM: "EA"}
	if _, err = repo.CreateWarehouseReplenishment(ctx, tenant, warehouseTestIDs(t, "replenish-no-breakbulk"), replenishCommand); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("partial replenishment without authorization error=%v", err)
	}
	replenishCommand.RequestID, replenishCommand.AllowBreakbulk = "replenish-breakbulk", true
	replenishment, err := repo.CreateWarehouseReplenishment(ctx, tenant, warehouseTestIDs(t, "replenish-breakbulk"), replenishCommand)
	if err != nil || len(replenishment.Lines) != 1 || replenishment.Lines[0].FromUOMCode != "BOX" || replenishment.Lines[0].UOMCode != "EA" || !exactDecimalEqual(replenishment.Lines[0].Quantity, "6") {
		t.Fatalf("packaging-aware replenishment=%+v err=%v", replenishment, err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", replenishment.ID, 1, posting.AddDate(0, 0, 6), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	var replenished string
	if err = pool.QueryRow(ctx, `select u.quantity::text from inventory.bulk_uom_balance u join inventory.bulk_balance b using(balance_id) where b.tenant_id=$1 and b.item_id='replenish-split' and b.bin_id='replenish-target' and u.uom_code='EA'`, tenant).Scan(&replenished); err != nil || replenished != "6.000000" {
		t.Fatalf("replenished EA=%s err=%v", replenished, err)
	}

	exactPick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "pick-exact-box"), inventorycontrol.WarehousePickCommand{RequestID: "pick-exact-box-request", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "preserve", DemandKind: "manual", DemandID: "order-box", DemandLineID: "line-box", Quantity: "12", HandlingUOM: "BOX"})
	if err != nil || len(exactPick.Lines) != 1 || exactPick.Lines[0].FromUOMCode != "BOX" || exactPick.Lines[0].UOMCode != "BOX" || !exactDecimalEqual(exactPick.Lines[0].FromUOMQuantity, "1") || !exactDecimalEqual(exactPick.Lines[0].UOMQuantity, "1") {
		t.Fatalf("exact-package pick=%+v err=%v", exactPick, err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", exactPick.ID, 1, posting.AddDate(0, 0, 7), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	var exactShipPackages string
	if err = pool.QueryRow(ctx, `select u.quantity::text from inventory.bulk_uom_balance u join inventory.bulk_balance b using(balance_id) where b.tenant_id=$1 and b.item_id='preserve' and b.bin_id='ship' and u.uom_code='BOX'`, tenant).Scan(&exactShipPackages); err != nil || exactShipPackages != "1.000000" {
		t.Fatalf("exact package at ship=%s err=%v", exactShipPackages, err)
	}
}
````

### FILE: `db/migrations/0036_warehouse_pick_replenishment_uom.up.sql`

```yaml
block_id: "GO-SUPPLY-WAREHOUSE-PICK-REPLENISHMENT-UOM-MIGRATION-UP:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps 2eae56d704a1fd035d104f333602aea7091b7749 invariants plus local verified adaptation"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "f34d1ad2e733e8c7917fcf4e3901181dd6a5813c5322946810cd9ebbe63eaedd"
variables: []
secrets_allowed: false
```

````sql
begin;

do $block$
begin
  if exists (
    select 1 from inventory.warehouse_activity
     where status='open' and activity_type in ('pick','movement')
  ) then
    raise exception using errcode='55000',
      message='drain open pick and movement activities before packaging-aware migration';
  end if;
end;
$block$;

alter table inventory.warehouse_activity_line
  add column from_uom_code text,
  add column from_uom_quantity numeric(20,6),
  add column from_qty_per_uom numeric(20,6);

alter table inventory.warehouse_activity_line disable trigger warehouse_activity_line_immutable;
update inventory.warehouse_activity_line
   set from_uom_code=uom_code,
       from_uom_quantity=uom_quantity,
       from_qty_per_uom=qty_per_uom;
alter table inventory.warehouse_activity_line enable trigger warehouse_activity_line_immutable;

alter table inventory.warehouse_activity_line
  alter column from_uom_code set not null,
  alter column from_uom_quantity set not null,
  alter column from_qty_per_uom set not null,
  add constraint warehouse_activity_line_from_uom_quantity_check check (from_uom_quantity > 0),
  add constraint warehouse_activity_line_from_qty_per_uom_check check (from_qty_per_uom > 0),
  add constraint warehouse_activity_line_from_uom_fk foreign key (tenant_id,item_id,from_uom_code)
    references inventory.item_unit_of_measure(tenant_id,item_id,uom_code);

create function inventory.validate_warehouse_line_from_uom()
returns trigger language plpgsql as $function$
declare
  expected_factor numeric(20,6);
  expected_precision numeric(20,6);
begin
  select u.qty_per_uom,u.rounding_precision
    into expected_factor,expected_precision
    from inventory.item_unit_of_measure u
   where u.tenant_id=new.tenant_id and u.item_id=new.item_id and u.uom_code=new.from_uom_code;
  if not found or new.from_qty_per_uom<>expected_factor or
     mod(new.from_uom_quantity,expected_precision)<>0 or
     new.from_uom_quantity*new.from_qty_per_uom<new.quantity then
    raise exception using errcode='55000',message='warehouse line source UOM does not match the immutable item UOM contract';
  end if;
  return new;
end;
$function$;

create trigger warehouse_activity_line_from_uom_validate
before insert or update on inventory.warehouse_activity_line
for each row execute function inventory.validate_warehouse_line_from_uom();

alter table inventory.warehouse_replenishment_request
  add column target_uom text,
  add column allow_breakbulk boolean not null default false;

update inventory.warehouse_replenishment_request r
   set target_uom=u.uom_code
  from inventory.item_unit_of_measure u
 where u.tenant_id=r.tenant_id and u.item_id=r.item_id and u.is_base;

alter table inventory.warehouse_replenishment_request
  alter column target_uom set not null,
  add constraint warehouse_replenishment_target_uom_fk foreign key (tenant_id,item_id,target_uom)
    references inventory.item_unit_of_measure(tenant_id,item_id,uom_code);

create table inventory.warehouse_activity_uom_reservation (
  tenant_id uuid not null,
  activity_id text not null,
  line_id text not null,
  organization_id text not null,
  source_balance_id bigint not null,
  source_uom_code text not null,
  source_uom_quantity numeric(20,6) not null check (source_uom_quantity > 0),
  source_qty_per_uom numeric(20,6) not null check (source_qty_per_uom > 0),
  source_base_quantity numeric(20,6) not null check (source_base_quantity > 0),
  target_uom_code text not null,
  target_uom_quantity numeric(20,6) not null check (target_uom_quantity > 0),
  target_qty_per_uom numeric(20,6) not null check (target_qty_per_uom > 0),
  moved_base_quantity numeric(20,6) not null check (moved_base_quantity > 0),
  conversion_required boolean not null,
  status text not null default 'open' check (status in ('open','consumed','released')),
  version bigint not null default 1 check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,activity_id,line_id),
  foreign key (tenant_id,activity_id,line_id)
    references inventory.warehouse_activity_line(tenant_id,activity_id,line_id),
  foreign key (tenant_id,organization_id)
    references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,activity_id,organization_id)
    references inventory.warehouse_activity(tenant_id,activity_id,organization_id),
  check (source_uom_quantity*source_qty_per_uom=source_base_quantity),
  check (target_uom_quantity*target_qty_per_uom=moved_base_quantity),
  check (source_base_quantity>=moved_base_quantity),
  check (conversion_required=(source_uom_code<>target_uom_code))
);

create function inventory.enforce_warehouse_uom_reservation_transition()
returns trigger language plpgsql as $function$
begin
  if old.status<>'open' or new.status not in ('consumed','released') or
     new.version<>old.version+1 or
     new.tenant_id<>old.tenant_id or new.activity_id<>old.activity_id or
     new.line_id<>old.line_id or new.organization_id<>old.organization_id or
     new.source_balance_id<>old.source_balance_id or
     new.source_uom_code<>old.source_uom_code or
     new.source_uom_quantity<>old.source_uom_quantity or
     new.source_qty_per_uom<>old.source_qty_per_uom or
     new.source_base_quantity<>old.source_base_quantity or
     new.target_uom_code<>old.target_uom_code or
     new.target_uom_quantity<>old.target_uom_quantity or
     new.target_qty_per_uom<>old.target_qty_per_uom or
     new.moved_base_quantity<>old.moved_base_quantity or
     new.conversion_required<>old.conversion_required or
     new.created_at<>old.created_at then
    raise exception using errcode='55000',message='invalid warehouse UOM reservation transition';
  end if;
  new.updated_at:=clock_timestamp();
  return new;
end;
$function$;

create trigger warehouse_uom_reservation_transition_guard
before update on inventory.warehouse_activity_uom_reservation
for each row execute function inventory.enforce_warehouse_uom_reservation_transition();

create function inventory.reject_warehouse_uom_reservation_delete()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000',message='warehouse UOM reservation evidence is immutable';
end;
$function$;

create trigger warehouse_uom_reservation_delete_guard
before delete on inventory.warehouse_activity_uom_reservation
for each row execute function inventory.reject_warehouse_uom_reservation_delete();

commit;
````

### FILE: `db/migrations/0036_warehouse_pick_replenishment_uom.down.sql`

```yaml
block_id: "GO-SUPPLY-WAREHOUSE-PICK-REPLENISHMENT-UOM-MIGRATION-DOWN:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps 2eae56d704a1fd035d104f333602aea7091b7749 invariants plus local verified adaptation"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "63bb3cce88f40f0b03c2e6de8ce87019770407f3af9e11e38fab2861d85f4e87"
variables: []
secrets_allowed: false
```

````sql
begin;

do $block$
begin
  if exists (select 1 from inventory.warehouse_activity_uom_reservation) or
     exists (
       select 1 from inventory.warehouse_activity_line
        where from_uom_code<>uom_code or
              from_uom_quantity<>uom_quantity or
              from_qty_per_uom<>qty_per_uom
     ) then
    raise exception using errcode='55000',
      message='cannot remove packaging-aware warehouse schema after evidence exists';
  end if;
end;
$block$;

drop trigger warehouse_uom_reservation_delete_guard on inventory.warehouse_activity_uom_reservation;
drop function inventory.reject_warehouse_uom_reservation_delete();
drop trigger warehouse_uom_reservation_transition_guard on inventory.warehouse_activity_uom_reservation;
drop function inventory.enforce_warehouse_uom_reservation_transition();
drop table inventory.warehouse_activity_uom_reservation;

alter table inventory.warehouse_replenishment_request
  drop constraint warehouse_replenishment_target_uom_fk,
  drop column allow_breakbulk,
  drop column target_uom;

drop trigger warehouse_activity_line_from_uom_validate on inventory.warehouse_activity_line;
drop function inventory.validate_warehouse_line_from_uom();
alter table inventory.warehouse_activity_line
  drop constraint warehouse_activity_line_from_uom_fk,
  drop constraint warehouse_activity_line_from_qty_per_uom_check,
  drop constraint warehouse_activity_line_from_uom_quantity_check,
  drop column from_qty_per_uom,
  drop column from_uom_quantity,
  drop column from_uom_code;

commit;
````

### FILE: `internal/platform/postgres/warehouse_packaging.go`

```yaml
block_id: "GO-SUPPLY-WAREHOUSE-PICK-REPLENISHMENT-PACKAGING:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps 2eae56d704a1fd035d104f333602aea7091b7749 invariants plus local verified adaptation"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "b66cdafe3e5f138a84bd6200d97e2d3d261af108f4df0f97364551f5eef5c1b1"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
)

type warehouseUOMContract struct {
	code      string
	factor    string
	precision string
}

type warehouseUOMPlan struct {
	balanceID          int64
	sourceUOM          string
	sourceQuantity     string
	sourceFactor       string
	sourceBase         string
	targetUOM          string
	targetQuantity     string
	targetFactor       string
	movedBase          string
	requiresConversion bool
}

func resolveWarehouseUOM(ctx context.Context, tx pgx.Tx, tenant, item, requested string) (warehouseUOMContract, error) {
	var value warehouseUOMContract
	err := tx.QueryRow(ctx, `select u.uom_code,u.qty_per_uom::text,u.rounding_precision::text from inventory.item_unit_of_measure u where u.tenant_id=$1 and u.item_id=$2 and u.uom_code=coalesce(nullif($3,''),(select b.uom_code from inventory.item_unit_of_measure b where b.tenant_id=$1 and b.item_id=$2 and b.is_base)) for share of u`, tenant, item, requested).Scan(&value.code, &value.factor, &value.precision)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, inventorycontrol.ErrConflict
	}
	return value, err
}

func floorToQuantum(value, quantum *big.Rat) *big.Rat {
	if value.Sign() <= 0 || quantum.Sign() <= 0 {
		return new(big.Rat)
	}
	steps := new(big.Rat).Quo(value, quantum)
	whole := new(big.Int).Quo(steps.Num(), steps.Denom())
	return new(big.Rat).Mul(new(big.Rat).SetInt(whole), quantum)
}

func ceilToQuantum(value, quantum *big.Rat) *big.Rat {
	if value.Sign() <= 0 || quantum.Sign() <= 0 {
		return new(big.Rat)
	}
	steps := new(big.Rat).Quo(value, quantum)
	whole, remainder := new(big.Int), new(big.Int)
	whole.QuoRem(steps.Num(), steps.Denom(), remainder)
	if remainder.Sign() > 0 {
		whole.Add(whole, big.NewInt(1))
	}
	return new(big.Rat).Mul(new(big.Rat).SetInt(whole), quantum)
}

func planWarehousePackaging(ctx context.Context, tx pgx.Tx, balanceID int64, target warehouseUOMContract, allowBreakbulk bool, maximumBase *big.Rat) ([]warehouseUOMPlan, *big.Rat, error) {
	targetFactor, targetFactorOK := parseRat(target.factor)
	targetPrecision, targetPrecisionOK := parseRat(target.precision)
	if !targetFactorOK || !targetPrecisionOK {
		return nil, new(big.Rat), inventorycontrol.ErrConflict
	}
	targetQuantum := new(big.Rat).Mul(targetFactor, targetPrecision)
	remaining := new(big.Rat).Set(maximumBase)
	rows, err := tx.Query(ctx, `select c.uom_code,c.qty_per_uom::text,u.rounding_precision::text,(c.quantity-c.reserved_quantity)::text from inventory.bulk_uom_balance c join inventory.item_unit_of_measure u on u.tenant_id=(select b.tenant_id from inventory.bulk_balance b where b.balance_id=c.balance_id) and u.item_id=(select b.item_id from inventory.bulk_balance b where b.balance_id=c.balance_id) and u.uom_code=c.uom_code where c.balance_id=$1 and c.quantity>c.reserved_quantity and (c.uom_code=$2 or ($3 and c.qty_per_uom>$4::numeric)) order by (c.uom_code=$2) desc,c.qty_per_uom,c.uom_code for update of c`, balanceID, target.code, allowBreakbulk, target.factor)
	if err != nil {
		return nil, remaining, err
	}
	type composition struct{ code, factor, precision, available string }
	available := []composition{}
	for rows.Next() {
		var value composition
		if err = rows.Scan(&value.code, &value.factor, &value.precision, &value.available); err != nil {
			rows.Close()
			return nil, remaining, err
		}
		available = append(available, value)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return nil, remaining, err
	}
	plans := []warehouseUOMPlan{}
	for _, source := range available {
		if remaining.Sign() <= 0 {
			break
		}
		sourceFactor, factorOK := parseRat(source.factor)
		sourcePrecision, precisionOK := parseRat(source.precision)
		sourceAvailable, availableOK := parseRat(source.available)
		if !factorOK || !precisionOK || !availableOK {
			return nil, remaining, inventorycontrol.ErrConflict
		}
		sourceAvailableBase := new(big.Rat).Mul(sourceAvailable, sourceFactor)
		move := floorToQuantum(minRat(remaining, sourceAvailableBase), targetQuantum)
		if move.Sign() <= 0 {
			continue
		}
		sourceQuantity := new(big.Rat).Quo(move, sourceFactor)
		if source.code != target.code {
			sourceQuantity = ceilToQuantum(sourceQuantity, sourcePrecision)
		}
		if sourceQuantity.Cmp(sourceAvailable) > 0 {
			continue
		}
		sourceBase := new(big.Rat).Mul(sourceQuantity, sourceFactor)
		targetQuantity := new(big.Rat).Quo(move, targetFactor)
		plans = append(plans, warehouseUOMPlan{
			balanceID: balanceID, sourceUOM: source.code, sourceQuantity: formatRat(sourceQuantity, 6),
			sourceFactor: source.factor, sourceBase: formatRat(sourceBase, 6), targetUOM: target.code,
			targetQuantity: formatRat(targetQuantity, 6), targetFactor: target.factor,
			movedBase: formatRat(move, 6), requiresConversion: source.code != target.code,
		})
		remaining.Sub(remaining, move)
	}
	return plans, remaining, nil
}

func reserveWarehousePackaging(ctx context.Context, tx pgx.Tx, tenant, activityID, organization string, plan warehouseUOMPlan, line inventorycontrol.WarehouseActivityLine) error {
	updated, err := tx.Exec(ctx, `update inventory.bulk_balance set reserved_quantity=reserved_quantity+$2::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and quantity-reserved_quantity >= $2::numeric`, plan.balanceID, plan.movedBase)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return bulkConflict(err)
		}
		return inventorycontrol.ErrConflict
	}
	updated, err = tx.Exec(ctx, `update inventory.bulk_uom_balance set reserved_quantity=reserved_quantity+$3::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and quantity-reserved_quantity >= $3::numeric`, plan.balanceID, plan.sourceUOM, plan.sourceQuantity)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return bulkConflict(err)
		}
		return inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity_uom_reservation(tenant_id,activity_id,line_id,organization_id,source_balance_id,source_uom_code,source_uom_quantity,source_qty_per_uom,source_base_quantity,target_uom_code,target_uom_quantity,target_qty_per_uom,moved_base_quantity,conversion_required) values($1,$2,$3,$4,$5,$6,$7::numeric,$8::numeric,$9::numeric,$10,$11::numeric,$12::numeric,$13::numeric,$14)`, tenant, activityID, line.ID, organization, plan.balanceID, plan.sourceUOM, plan.sourceQuantity, plan.sourceFactor, plan.sourceBase, plan.targetUOM, plan.targetQuantity, plan.targetFactor, plan.movedBase, plan.requiresConversion)
	return bulkConflict(err)
}

func deterministicWarehouseUUID(parts ...string) string {
	sum := sha256.Sum256([]byte(fmt.Sprint(parts)))
	sum[6] = (sum[6] & 0x0f) | 0x50
	sum[8] = (sum[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", sum[0:4], sum[4:6], sum[6:8], sum[8:10], sum[10:16])
}

func moveWarehousePackagingSource(ctx context.Context, tx pgx.Tx, tenant, organization, activityID, sourceKind string, line inventorycontrol.WarehouseActivityLine) error {
	var balanceID int64
	var sourceUOM, sourceQuantity, sourceFactor, sourceBase, targetUOM, targetQuantity, targetFactor, movedBase string
	var conversion bool
	err := tx.QueryRow(ctx, `select source_balance_id,source_uom_code,source_uom_quantity::text,source_qty_per_uom::text,source_base_quantity::text,target_uom_code,target_uom_quantity::text,target_qty_per_uom::text,moved_base_quantity::text,conversion_required from inventory.warehouse_activity_uom_reservation where tenant_id=$1 and activity_id=$2 and line_id=$3 and organization_id=$4 and status='open' and version=1 for update`, tenant, activityID, line.ID, organization).Scan(&balanceID, &sourceUOM, &sourceQuantity, &sourceFactor, &sourceBase, &targetUOM, &targetQuantity, &targetFactor, &movedBase, &conversion)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.ErrConflict
	}
	if err != nil {
		return err
	}
	if movedBase != line.Quantity || sourceUOM != line.FromUOMCode || sourceQuantity != line.FromUOMQuantity || sourceFactor != line.FromQuantityPerUOM || targetUOM != line.UOMCode || targetQuantity != line.UOMQuantity || targetFactor != line.QuantityPerUOM {
		return inventorycontrol.ErrConflict
	}
	var sourceIsBase, targetIsBase bool
	var targetPrecision string
	err = tx.QueryRow(ctx, `select s.is_base,t.is_base,t.rounding_precision::text from inventory.bulk_balance b join inventory.item_unit_of_measure s on s.tenant_id=b.tenant_id and s.item_id=b.item_id and s.uom_code=$2 join inventory.item_unit_of_measure t on t.tenant_id=b.tenant_id and t.item_id=b.item_id and t.uom_code=$3 where b.balance_id=$1 for share of b,s,t`, balanceID, sourceUOM, targetUOM).Scan(&sourceIsBase, &targetIsBase, &targetPrecision)
	if err != nil {
		return err
	}
	if sourceIsBase {
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_uom_balance set reserved_quantity=reserved_quantity-$3::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and reserved_quantity >= $3::numeric`, balanceID, sourceUOM, sourceQuantity)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return bulkConflict(updateErr)
			}
			return inventorycontrol.ErrConflict
		}
	} else {
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_uom_balance set quantity=quantity-$3::numeric,reserved_quantity=reserved_quantity-$3::numeric,quantity_base=quantity_base-$4::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and quantity >= $3::numeric and reserved_quantity >= $3::numeric and quantity_base >= $4::numeric`, balanceID, sourceUOM, sourceQuantity, sourceBase)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return bulkConflict(updateErr)
			}
			return inventorycontrol.ErrConflict
		}
		if _, err = tx.Exec(ctx, `delete from inventory.bulk_uom_balance where balance_id=$1 and uom_code=$2 and quantity_base=0 and reserved_quantity=0`, balanceID, sourceUOM); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base) select $1,u.uom_code,$2::numeric/u.qty_per_uom,u.qty_per_uom,$2::numeric from inventory.bulk_balance b join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.is_base where b.balance_id=$1 on conflict(balance_id,uom_code) do update set quantity=inventory.bulk_uom_balance.quantity+excluded.quantity,quantity_base=inventory.bulk_uom_balance.quantity_base+excluded.quantity_base,version=inventory.bulk_uom_balance.version+1,updated_at=clock_timestamp()`, balanceID, sourceBase)
		if err != nil {
			return bulkConflict(err)
		}
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_balance set quantity=quantity-$2::numeric,reserved_quantity=reserved_quantity-$2::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and quantity >= $2::numeric and reserved_quantity >= $2::numeric`, balanceID, movedBase)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return err
		}
		return inventorycontrol.ErrConflict
	}
	if conversion {
		toTotal, exact := exactUOMQuantity(sourceBase, targetFactor, targetPrecision)
		if !exact {
			return inventorycontrol.ErrConflict
		}
		conversionID := activityID + ":" + line.ID + ":breakbulk"
		requestID := activityID + ":" + line.ID + ":automatic-breakbulk"
		_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_conversion(tenant_id,conversion_id,request_id,organization_id,bin_id,item_id,lot_id,operation,from_uom,to_uom,from_quantity,to_quantity,base_quantity,from_qty_per_uom,to_qty_per_uom,source_kind,source_id) values($1,$2,$3,$4,$5,$6,nullif($7,''),'breakbulk',$8,$9,$10::numeric,$11::numeric,$12::numeric,$13::numeric,$14::numeric,$15,$16)`, tenant, conversionID, requestID, organization, line.FromBinID, line.ItemID, line.LotID, sourceUOM, targetUOM, sourceQuantity, toTotal, sourceBase, sourceFactor, targetFactor, sourceKind, activityID)
		if err != nil {
			return bulkConflict(err)
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_conversion_line(tenant_id,conversion_id,line_id,sequence_no,action_type,uom_code,quantity,quantity_base) values($1,$2,$3,1,'take',$5,$6::numeric,$7::numeric),($1,$2,$4,2,'place',$8,$9::numeric,$7::numeric)`, tenant, conversionID, conversionID+":take", conversionID+":place", sourceUOM, sourceQuantity, sourceBase, targetUOM, toTotal)
		if err != nil {
			return bulkConflict(err)
		}
		if err = recordBulkEvent(ctx, tx, tenant, deterministicWarehouseUUID(tenant, activityID, line.ID, "breakbulk"), "bulk-uom-conversion", conversionID, "bulk-uom-conversion.completed", 1, map[string]any{"organization_id": organization, "activity_id": activityID, "line_id": line.ID, "from_uom": sourceUOM, "to_uom": targetUOM, "source_base_quantity": sourceBase, "moved_base_quantity": movedBase}); err != nil {
			return err
		}
	}
	leftover, leftOK := parseRat(sourceBase)
	moved, movedOK := parseRat(movedBase)
	if !leftOK || !movedOK {
		return inventorycontrol.ErrConflict
	}
	leftover.Sub(leftover, moved)
	if leftover.Sign() > 0 && !targetIsBase {
		leftQuantity, exact := exactUOMQuantity(formatRat(leftover, 6), targetFactor, targetPrecision)
		if !exact {
			return inventorycontrol.ErrConflict
		}
		var baseUOM string
		if err = tx.QueryRow(ctx, `select uom_code from inventory.item_unit_of_measure u join inventory.bulk_balance b on b.tenant_id=u.tenant_id and b.item_id=u.item_id where b.balance_id=$1 and u.is_base`, balanceID).Scan(&baseUOM); err != nil {
			return err
		}
		updated, err = tx.Exec(ctx, `update inventory.bulk_uom_balance set quantity=quantity-($2::numeric/qty_per_uom),quantity_base=quantity_base-$2::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$3 and quantity_base >= $2::numeric`, balanceID, formatRat(leftover, 6), baseUOM)
		if err != nil || updated.RowsAffected() != 1 {
			if err != nil {
				return bulkConflict(err)
			}
			return inventorycontrol.ErrConflict
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base) values($1,$2,$3::numeric,$4::numeric,$5::numeric) on conflict(balance_id,uom_code) do update set quantity=inventory.bulk_uom_balance.quantity+excluded.quantity,quantity_base=inventory.bulk_uom_balance.quantity_base+excluded.quantity_base,version=inventory.bulk_uom_balance.version+1,updated_at=clock_timestamp()`, balanceID, targetUOM, leftQuantity, targetFactor, formatRat(leftover, 6))
		if err != nil {
			return bulkConflict(err)
		}
	}
	updated, err = tx.Exec(ctx, `update inventory.warehouse_activity_uom_reservation set status='consumed',version=version+1 where tenant_id=$1 and activity_id=$2 and line_id=$3 and status='open' and version=1`, tenant, activityID, line.ID)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return err
		}
		return inventorycontrol.ErrConflict
	}
	return nil
}

func releaseWarehousePackagingReservation(ctx context.Context, tx pgx.Tx, tenant, organization, activityID string, line inventorycontrol.WarehouseActivityLine) error {
	var balanceID int64
	var sourceUOM, sourceQuantity, movedBase string
	err := tx.QueryRow(ctx, `select source_balance_id,source_uom_code,source_uom_quantity::text,moved_base_quantity::text from inventory.warehouse_activity_uom_reservation where tenant_id=$1 and activity_id=$2 and line_id=$3 and organization_id=$4 and status='open' and version=1 for update`, tenant, activityID, line.ID, organization).Scan(&balanceID, &sourceUOM, &sourceQuantity, &movedBase)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.ErrConflict
	}
	if err != nil {
		return err
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_balance set reserved_quantity=reserved_quantity-$2::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and reserved_quantity >= $2::numeric`, balanceID, movedBase)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return err
		}
		return inventorycontrol.ErrConflict
	}
	updated, err = tx.Exec(ctx, `update inventory.bulk_uom_balance set reserved_quantity=reserved_quantity-$3::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and reserved_quantity >= $3::numeric`, balanceID, sourceUOM, sourceQuantity)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return err
		}
		return inventorycontrol.ErrConflict
	}
	updated, err = tx.Exec(ctx, `update inventory.warehouse_activity_uom_reservation set status='released',version=version+1 where tenant_id=$1 and activity_id=$2 and line_id=$3 and status='open' and version=1`, tenant, activityID, line.ID)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return err
		}
		return inventorycontrol.ErrConflict
	}
	return nil
}

func consumeRegisteredPickPackaging(ctx context.Context, tx pgx.Tx, tenant, organization, activityID, reservationID string) error {
	var balanceID int64
	var uomCode, uomQuantity, baseQuantity, baseUOM, baseFactor string
	var isBase bool
	err := tx.QueryRow(ctx, `select b.balance_id,al.uom_code,al.uom_quantity::text,al.quantity::text,u.is_base,base.uom_code,base.qty_per_uom::text from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) join inventory.bulk_reservation r on r.tenant_id=al.tenant_id and r.reservation_id=al.reservation_id join inventory.bulk_balance b on b.tenant_id=r.tenant_id and b.organization_id=r.organization_id and b.bin_id=r.bin_id and b.item_id=r.item_id and b.lot_id is not distinct from r.lot_id join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.uom_code=al.uom_code join inventory.item_unit_of_measure base on base.tenant_id=b.tenant_id and base.item_id=b.item_id and base.is_base where a.tenant_id=$1 and a.organization_id=$2 and a.activity_id=$3 and a.activity_type='pick' and a.status='registered' and r.reservation_id=$4 and r.status='reservation' and b.quantity>=al.quantity and b.reserved_quantity>=al.quantity for update of b,u,base`, tenant, organization, activityID, reservationID).Scan(&balanceID, &uomCode, &uomQuantity, &baseQuantity, &isBase, &baseUOM, &baseFactor)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.ErrConflict
	}
	if err != nil {
		return err
	}
	if isBase {
		return nil
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_uom_balance set quantity=quantity-$3::numeric,quantity_base=quantity_base-$4::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and quantity-reserved_quantity >= $3::numeric and quantity_base >= $4::numeric`, balanceID, uomCode, uomQuantity, baseQuantity)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return bulkConflict(err)
		}
		return inventorycontrol.ErrConflict
	}
	if _, err = tx.Exec(ctx, `delete from inventory.bulk_uom_balance where balance_id=$1 and uom_code=$2 and quantity_base=0 and reserved_quantity=0`, balanceID, uomCode); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base) values($1,$2,$3::numeric/$4::numeric,$4::numeric,$3::numeric) on conflict(balance_id,uom_code) do update set quantity=inventory.bulk_uom_balance.quantity+excluded.quantity,quantity_base=inventory.bulk_uom_balance.quantity_base+excluded.quantity_base,version=inventory.bulk_uom_balance.version+1,updated_at=clock_timestamp()`, balanceID, baseUOM, baseQuantity, baseFactor)
	return bulkConflict(err)
}
````

### FILE: `internal/platform/postgres/sales_warehouse.go`

```yaml
block_id: "GO-SUPPLY-SALES-WAREHOUSE-POSTGRES:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps 2eae56d704a1fd035d104f333602aea7091b7749 sales line, source document and warehouse pick invariants"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "8e32d5380a41084b86c20f989ba3980d5f72348068fc1f9042e0251a8a3f2d82"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
)

func (r *InventoryControl) ConfigureSalesWarehouseBinding(ctx context.Context, tenant, eventID string, value inventorycontrol.SalesWarehouseBinding) (inventorycontrol.SalesWarehouseBinding, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	var existing inventorycontrol.SalesWarehouseBinding
	err = tx.QueryRow(ctx, `select request_id,variant_id,item_id,sales_uom_code,version from inventory.sales_warehouse_binding where tenant_id=$1 and request_id=$2 for update`, tenant, value.RequestID).Scan(&existing.RequestID, &existing.VariantID, &existing.ItemID, &existing.SalesUOMCode, &existing.Version)
	if err == nil {
		if existing.VariantID != value.VariantID || existing.ItemID != value.ItemID || existing.SalesUOMCode != value.SalesUOMCode {
			return value, inventorycontrol.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return value, err
		}
		return existing, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return value, err
	}
	_, err = tx.Exec(ctx, `insert into inventory.sales_warehouse_binding(tenant_id,request_id,variant_id,item_id,sales_uom_code,version) values($1,$2,$3,$4,$5,1)`, tenant, value.RequestID, value.VariantID, value.ItemID, value.SalesUOMCode)
	if err != nil {
		return value, bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "sales-warehouse-binding", value.VariantID, "sales-warehouse-binding.configured", 1, map[string]any{"item_id": value.ItemID, "sales_uom_code": value.SalesUOMCode, "request_id": value.RequestID}); err != nil {
		return value, err
	}
	value.Version = 1
	return value, tx.Commit(ctx)
}

func readWarehousePickReplay(ctx context.Context, tx pgx.Tx, tenant string, value inventorycontrol.WarehousePickCommand) (inventorycontrol.WarehouseActivity, bool, error) {
	var activityID, organization, shipBin, itemID, demandKind, demandID, demandLineID, quantity, requestedUOM, assignedTo string
	var allowBreakbulk, useFEFO, allowDedicated bool
	err := tx.QueryRow(ctx, `select activity_id,organization_id,ship_bin_id,item_id,demand_kind,demand_id,demand_line_id,quantity::text,requested_handling_uom,allow_breakbulk,use_fefo,allow_dedicated,assigned_to from inventory.warehouse_pick_request where tenant_id=$1 and request_id=$2 for update`, tenant, value.RequestID).Scan(&activityID, &organization, &shipBin, &itemID, &demandKind, &demandID, &demandLineID, &quantity, &requestedUOM, &allowBreakbulk, &useFEFO, &allowDedicated, &assignedTo)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.WarehouseActivity{}, false, nil
	}
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, false, err
	}
	storedQuantity, storedOK := parseRat(quantity)
	requestedQuantity, requestedOK := parseRat(value.Quantity)
	if !storedOK || !requestedOK || storedQuantity.Cmp(requestedQuantity) != 0 || organization != value.OrganizationID || shipBin != value.ShipBinID || itemID != value.ItemID || demandKind != value.DemandKind || demandID != value.DemandID || demandLineID != value.DemandLineID || requestedUOM != value.HandlingUOM || allowBreakbulk != value.AllowBreakbulk || useFEFO != value.UseFEFO || allowDedicated != value.AllowDedicated || assignedTo != value.AssignedTo {
		return inventorycontrol.WarehouseActivity{}, false, fmt.Errorf("warehouse pick divergent request replay: %w", inventorycontrol.ErrConflict)
	}
	activity, err := readWarehouseActivity(ctx, tx, tenant, organization, activityID)
	return activity, true, err
}

func validateCustomerOrderWarehouseDemand(ctx context.Context, tx pgx.Tx, tenant string, value inventorycontrol.WarehousePickCommand) (string, string, string, error) {
	var salesUOM, totalBase string
	err := tx.QueryRow(ctx, `select b.sales_uom_code,(l.quantity::numeric*u.qty_per_uom)::text from sales.customer_order o join sales.customer_order_line l using(tenant_id,order_id) join inventory.sales_warehouse_binding b on b.tenant_id=l.tenant_id and b.variant_id=l.variant_id join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.uom_code=b.sales_uom_code where o.tenant_id=$1 and o.organization_id=$2 and o.order_id=$3 and l.line_id=$4 and b.item_id=$5 and o.state in ('placed','confirmed','paid','allocated') for share of o,l,b,u`, tenant, value.OrganizationID, value.DemandID, value.DemandLineID, value.ItemID).Scan(&salesUOM, &totalBase)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", "", fmt.Errorf("warehouse pick customer-order source contract: %w", inventorycontrol.ErrConflict)
	}
	if err != nil {
		return "", "", "", err
	}
	var handledBase string
	if err = tx.QueryRow(ctx, `select coalesce(sum(al.quantity),0)::text from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=$1 and a.organization_id=$2 and a.activity_type='pick' and a.source_kind='customer-order' and a.source_id=$3 and a.source_line_id=$4 and a.status in ('open','registered')`, tenant, value.OrganizationID, value.DemandID, value.DemandLineID).Scan(&handledBase); err != nil {
		return "", "", "", err
	}
	total, totalOK := parseRat(totalBase)
	handled, handledOK := parseRat(handledBase)
	requested, requestedOK := parseRat(value.Quantity)
	if !totalOK || !handledOK || !requestedOK {
		return "", "", "", fmt.Errorf("warehouse pick customer-order quantity parse: %w", inventorycontrol.ErrConflict)
	}
	outstanding := new(big.Rat).Sub(total, handled)
	if outstanding.Sign() <= 0 || requested.Cmp(outstanding) > 0 {
		return "", "", "", fmt.Errorf("warehouse pick customer-order outstanding: %w", inventorycontrol.ErrConflict)
	}
	return salesUOM, formatRat(total, 6), formatRat(outstanding, 6), nil
}
````

### FILE: `internal/platform/postgres/sales_warehouse_integration_test.go`

```yaml
block_id: "GO-SUPPLY-SALES-WAREHOUSE-POSTGRES-TEST:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps 2eae56d704a1fd035d104f333602aea7091b7749 SCM warehouse sales-order scenarios"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "c800a485be25357716290619cdc29538e7a74a780b93c40f27c087ce8b56013e"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestSalesOrderWarehouseDemandIsBoundOutstandingIdempotentAndConcurrent(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := bulkTestUUID(t)
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Connected Sales Warehouse','Connected Sales Warehouse')`, tenant, "swh-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'warehouse','warehouse','Warehouse','warehouse')`, tenant); err != nil {
		t.Fatal(err)
	}
	defer cleanupSalesWarehouse(t, pool, tenant)
	for _, command := range []string{
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'model','MODEL','Model','other','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,homologation_state,lifecycle_state) values($1,'variant','model','VARIANT','Variant','{}','approved','active')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,'order','warehouse','customer','draft','USD',300,1)`,
		`insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units) values($1,'order','line','variant',3,100)`,
	} {
		if _, err = pool.Exec(ctx, command, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewInventoryControl(pool)
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "part", Code: "SALE-PART", Description: "Sale part", BaseUOM: "EA", BaseRoundingPrecision: "1", TrackingMode: "none", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.ConfigureItemUnitOfMeasure(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemUnitOfMeasure{ItemID: "part", Code: "BOX", QuantityPerUnit: "12", RoundingPrecision: "1", Version: 1}); err != nil {
		t.Fatal(err)
	}
	for _, bin := range []inventorycontrol.WarehouseBin{{ID: "receive", OrganizationID: "warehouse", Code: "RECEIVE", Type: "receive", Version: 1}, {ID: "pick", OrganizationID: "warehouse", Code: "PICK", Type: "putpick", Ranking: 100, Version: 1}, {ID: "ship", OrganizationID: "warehouse", Code: "SHIP", Type: "ship", Version: 1}} {
		if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatal(err)
		}
	}
	for _, policy := range []inventorycontrol.ItemBinPolicy{{OrganizationID: "warehouse", ItemID: "part", BinID: "receive", Fixed: true, MinQuantity: "0", MaxQuantity: "36", Version: 1}, {OrganizationID: "warehouse", ItemID: "part", BinID: "pick", Fixed: true, Default: true, MinQuantity: "0", MaxQuantity: "36", Version: 1}, {OrganizationID: "warehouse", ItemID: "part", BinID: "ship", Fixed: true, MinQuantity: "0", MaxQuantity: "36", Version: 1}} {
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), policy); err != nil {
			t.Fatal(err)
		}
	}
	posting := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	receipt, err := repo.PostWarehouseReceipt(ctx, tenant, warehousePackagingIDs(t, "sales-source"), inventorycontrol.WarehouseReceiptCommand{RequestID: "sales-source-receipt", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "part", HandlingUOM: "BOX", HandlingQuantity: "3", UnitCost: "10", PostingDate: posting, SourceKind: "purchase", SourceID: "sales-source-po"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", receipt.PutAway.ID, 1, posting.AddDate(0, 0, 1), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	binding := inventorycontrol.SalesWarehouseBinding{RequestID: "binding-request", VariantID: "variant", ItemID: "part", SalesUOMCode: "BOX", Version: 1}
	configured, err := repo.ConfigureSalesWarehouseBinding(ctx, tenant, bulkTestUUID(t), binding)
	if err != nil || configured != binding {
		t.Fatalf("binding=%+v err=%v", configured, err)
	}
	replayedBinding, err := repo.ConfigureSalesWarehouseBinding(ctx, tenant, bulkTestUUID(t), binding)
	if err != nil || replayedBinding != binding {
		t.Fatalf("binding replay=%+v err=%v", replayedBinding, err)
	}
	divergentBinding := binding
	divergentBinding.SalesUOMCode = "EA"
	if _, err = repo.ConfigureSalesWarehouseBinding(ctx, tenant, bulkTestUUID(t), divergentBinding); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("divergent binding error=%v", err)
	}
	draftPick := inventorycontrol.WarehousePickCommand{RequestID: "draft-pick", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "customer-order", DemandID: "order", DemandLineID: "line", Quantity: "12", HandlingUOM: "BOX"}
	if _, err = repo.CreateWarehousePick(ctx, tenant, warehousePackagingIDs(t, "draft-pick"), draftPick); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("draft order pick error=%v", err)
	}
	if _, err = pool.Exec(ctx, `update sales.customer_order set state='placed',version=2,updated_at=clock_timestamp() where tenant_id=$1 and order_id='order'`, tenant); err != nil {
		t.Fatal(err)
	}
	cancelCommand := draftPick
	cancelCommand.RequestID = "cancel-pick"
	cancellable, err := repo.CreateWarehousePick(ctx, tenant, warehousePackagingIDs(t, "cancel-pick"), cancelCommand)
	if err != nil || len(cancellable.Lines) != 1 || cancellable.Lines[0].FromUOMCode != "BOX" || cancellable.Lines[0].UOMCode != "BOX" {
		t.Fatalf("cancellable=%+v err=%v", cancellable, err)
	}
	replayed, err := repo.CreateWarehousePick(ctx, tenant, warehousePackagingIDs(t, "ignored-replay"), cancelCommand)
	if err != nil || replayed.ID != cancellable.ID || replayed.Status != "open" {
		t.Fatalf("pick replay=%+v err=%v", replayed, err)
	}
	divergentPick := cancelCommand
	divergentPick.Quantity = "6"
	if _, err = repo.CreateWarehousePick(ctx, tenant, warehousePackagingIDs(t, "divergent-replay"), divergentPick); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("divergent pick replay error=%v", err)
	}
	cancelled, err := repo.CancelWarehousePick(ctx, tenant, "warehouse", cancellable.ID, 1, bulkTestUUID(t))
	if err != nil || cancelled.Status != "cancelled" {
		t.Fatalf("cancelled=%+v err=%v", cancelled, err)
	}
	registeredCommand := cancelCommand
	registeredCommand.RequestID = "registered-pick"
	registeredPick, err := repo.CreateWarehousePick(ctx, tenant, warehousePackagingIDs(t, "registered-pick"), registeredCommand)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", registeredPick.ID, 1, posting.AddDate(0, 0, 2), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	firstShipmentCommand := inventorycontrol.CustomerShipmentCommand{RequestID: "shipment-first-request", OrganizationID: "warehouse", OrderID: "order", WarehouseActivityID: registeredPick.ID, PostingDate: posting.AddDate(0, 0, 3)}
	firstShipment, err := repo.PostCustomerShipment(ctx, tenant, inventorycontrol.CustomerShipmentIDs{ShipmentID: "shipment-first", LineID: "shipment-first-line", EventID: bulkTestUUID(t)}, firstShipmentCommand)
	if err != nil || firstShipment.ID != "shipment-first" || firstShipment.FulfillmentState != "partially-shipped" || firstShipment.OrderVersion != 3 || firstShipment.Line.SalesUOMCode != "BOX" || !exactDecimalEqual(firstShipment.Line.Quantity, "1") || !exactDecimalEqual(firstShipment.Line.QuantityBase, "12") || !exactDecimalEqual(firstShipment.Line.CostAmount, "120") || firstShipment.Line.AllocationCount != 1 {
		t.Fatalf("first shipment=%+v err=%v", firstShipment, err)
	}
	firstReplay, err := repo.PostCustomerShipment(ctx, tenant, inventorycontrol.CustomerShipmentIDs{ShipmentID: "ignored-replay", LineID: "ignored-replay-line", EventID: bulkTestUUID(t)}, firstShipmentCommand)
	if err != nil || firstReplay.ID != firstShipment.ID || firstReplay.Line.ID != firstShipment.Line.ID || firstReplay.OrderVersion != firstShipment.OrderVersion {
		t.Fatalf("first shipment replay=%+v err=%v", firstReplay, err)
	}
	divergentShipment := firstShipmentCommand
	divergentShipment.PostingDate = divergentShipment.PostingDate.AddDate(0, 0, 1)
	if _, err = repo.PostCustomerShipment(ctx, tenant, inventorycontrol.CustomerShipmentIDs{ShipmentID: "ignored-divergent", LineID: "ignored-divergent-line", EventID: bulkTestUUID(t)}, divergentShipment); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("divergent shipment replay error=%v", err)
	}
	commands := []inventorycontrol.WarehousePickCommand{registeredCommand, registeredCommand}
	commands[0].RequestID, commands[0].Quantity = "concurrent-a", "24"
	commands[1].RequestID, commands[1].Quantity = "concurrent-b", "24"
	results := make([]inventorycontrol.WarehouseActivity, 2)
	errs := make([]error, 2)
	var wg sync.WaitGroup
	for index := range commands {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i], errs[i] = repo.CreateWarehousePick(ctx, tenant, warehousePackagingIDs(t, commands[i].RequestID), commands[i])
		}(index)
	}
	wg.Wait()
	winners := 0
	for index := range errs {
		if errs[index] == nil {
			winners++
			if len(results[index].Lines) != 1 || results[index].Lines[0].UOMCode != "BOX" || !exactDecimalEqual(results[index].Lines[0].UOMQuantity, "2") {
				t.Fatalf("concurrent winner=%+v", results[index])
			}
		} else if !errors.Is(errs[index], inventorycontrol.ErrConflict) {
			t.Fatalf("concurrent error[%d]=%v", index, errs[index])
		}
	}
	if winners != 1 {
		t.Fatalf("concurrent winners=%d errors=%v", winners, errs)
	}
	winningPick := inventorycontrol.WarehouseActivity{}
	for index := range errs {
		if errs[index] == nil {
			winningPick = results[index]
		}
	}
	registeredRemainder, err := repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", winningPick.ID, 1, posting.AddDate(0, 0, 4), bulkTestUUID(t))
	if err != nil || registeredRemainder.Status != "registered" {
		t.Fatalf("registered remainder=%+v err=%v", registeredRemainder, err)
	}
	shipmentCommands := []inventorycontrol.CustomerShipmentCommand{
		{RequestID: "shipment-race-a", OrganizationID: "warehouse", OrderID: "order", WarehouseActivityID: winningPick.ID, PostingDate: posting.AddDate(0, 0, 5)},
		{RequestID: "shipment-race-b", OrganizationID: "warehouse", OrderID: "order", WarehouseActivityID: winningPick.ID, PostingDate: posting.AddDate(0, 0, 5)},
	}
	shipmentResults := make([]inventorycontrol.CustomerShipment, 2)
	shipmentErrors := make([]error, 2)
	for index := range shipmentCommands {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			shipmentResults[i], shipmentErrors[i] = repo.PostCustomerShipment(ctx, tenant, inventorycontrol.CustomerShipmentIDs{ShipmentID: "shipment-race-" + string(rune('a'+i)), LineID: "shipment-race-line-" + string(rune('a'+i)), EventID: bulkTestUUID(t)}, shipmentCommands[i])
		}(index)
	}
	wg.Wait()
	shipmentWinners := 0
	winningShipmentIndex := -1
	for index := range shipmentErrors {
		if shipmentErrors[index] == nil {
			shipmentWinners++
			winningShipmentIndex = index
			if shipmentResults[index].FulfillmentState != "shipped" || shipmentResults[index].OrderVersion != 4 || !exactDecimalEqual(shipmentResults[index].Line.Quantity, "2") || !exactDecimalEqual(shipmentResults[index].Line.QuantityBase, "24") || !exactDecimalEqual(shipmentResults[index].Line.CostAmount, "240") {
				t.Fatalf("shipment winner=%+v", shipmentResults[index])
			}
		} else if !errors.Is(shipmentErrors[index], inventorycontrol.ErrConflict) {
			t.Fatalf("shipment error[%d]=%v", index, shipmentErrors[index])
		}
	}
	if shipmentWinners != 1 {
		t.Fatalf("shipment winners=%d errors=%v", shipmentWinners, shipmentErrors)
	}
	winningShipmentReplay, err := repo.PostCustomerShipment(ctx, tenant, inventorycontrol.CustomerShipmentIDs{ShipmentID: "ignored-winning-shipment", LineID: "ignored-winning-line", EventID: bulkTestUUID(t)}, shipmentCommands[winningShipmentIndex])
	if err != nil || winningShipmentReplay.ID != shipmentResults[winningShipmentIndex].ID || winningShipmentReplay.FulfillmentState != "shipped" {
		t.Fatalf("winning shipment replay=%+v err=%v", winningShipmentReplay, err)
	}
	var handled string
	if err = pool.QueryRow(ctx, `select coalesce(sum(al.quantity),0)::text from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=$1 and a.source_kind='customer-order' and a.source_id='order' and a.source_line_id='line' and a.status in ('open','registered')`, tenant).Scan(&handled); err != nil || !exactDecimalEqual(handled, "36") {
		t.Fatalf("handled=%s err=%v", handled, err)
	}
	blocked := registeredCommand
	blocked.RequestID = "beyond-outstanding"
	if _, err = repo.CreateWarehousePick(ctx, tenant, warehousePackagingIDs(t, "beyond-outstanding"), blocked); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("beyond outstanding error=%v", err)
	}
	var orderFulfillment, shippedQuantity, physicalQuantity, composedQuantity, remainingCost, shipmentCost string
	var orderVersionAfter, shipmentCount, consumedReservations int64
	if err = pool.QueryRow(ctx, `select o.fulfillment_state,o.version,l.shipped_quantity::text from sales.customer_order o join sales.customer_order_line l using(tenant_id,order_id) where o.tenant_id=$1 and o.order_id='order'`, tenant).Scan(&orderFulfillment, &orderVersionAfter, &shippedQuantity); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select coalesce(sum(b.quantity),0)::text,coalesce(sum(c.quantity_base),0)::text from inventory.bulk_balance b left join inventory.bulk_uom_balance c on c.balance_id=b.balance_id where b.tenant_id=$1 and b.item_id='part'`, tenant).Scan(&physicalQuantity, &composedQuantity); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select coalesce(sum(remaining_quantity),0)::text from inventory.bulk_cost_layer where tenant_id=$1 and item_id='part'`, tenant).Scan(&remainingCost); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*),coalesce(sum(total_cost_amount),0)::text from sales.customer_shipment where tenant_id=$1 and order_id='order'`, tenant).Scan(&shipmentCount, &shipmentCost); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from inventory.bulk_reservation where tenant_id=$1 and demand_kind='customer-order' and status='consumed'`, tenant).Scan(&consumedReservations); err != nil {
		t.Fatal(err)
	}
	if orderFulfillment != "shipped" || orderVersionAfter != 4 || !exactDecimalEqual(shippedQuantity, "3") || !exactDecimalEqual(physicalQuantity, "0") || !exactDecimalEqual(composedQuantity, "0") || !exactDecimalEqual(remainingCost, "0") || shipmentCount != 2 || !exactDecimalEqual(shipmentCost, "360") || consumedReservations != 2 {
		t.Fatalf("fulfillment=%s version=%d shipped=%s physical=%s composed=%s cost-remaining=%s shipments=%d shipment-cost=%s consumed=%d", orderFulfillment, orderVersionAfter, shippedQuantity, physicalQuantity, composedQuantity, remainingCost, shipmentCount, shipmentCost, consumedReservations)
	}
}

func cleanupSalesWarehouse(t *testing.T, pool *pgxpool.Pool, tenant string) {
	t.Helper()
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Errorf("sales warehouse cleanup begin: %v", err)
		return
	}
	defer tx.Rollback(ctx)
	commands := []string{
		`alter table inventory.customer_shipment_allocation disable trigger customer_shipment_allocation_immutable`,
		`alter table sales.customer_shipment_line disable trigger customer_shipment_line_immutable`,
		`alter table sales.customer_shipment disable trigger customer_shipment_immutable`,
		`alter table inventory.warehouse_pick_request disable trigger warehouse_pick_request_immutable`,
		`alter table inventory.sales_warehouse_binding disable trigger sales_warehouse_binding_immutable`,
		`alter table inventory.warehouse_activity_uom_reservation disable trigger warehouse_uom_reservation_delete_guard`,
		`alter table inventory.warehouse_activity_line disable trigger warehouse_activity_line_immutable`,
		`alter table inventory.warehouse_receipt_line disable trigger warehouse_receipt_line_immutable`,
		`alter table inventory.warehouse_receipt disable trigger warehouse_receipt_immutable`,
		`alter table inventory.bulk_inventory_entry disable trigger bulk_inventory_entry_immutable`,
		`alter table inventory.bulk_cost_application disable trigger bulk_cost_application_immutable`,
		`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`,
		`delete from inventory.warehouse_pick_request where tenant_id=$1`,
		`delete from inventory.customer_shipment_allocation where tenant_id=$1`,
		`delete from sales.customer_shipment_line where tenant_id=$1`,
		`delete from sales.customer_shipment where tenant_id=$1`,
		`delete from inventory.warehouse_activity_uom_reservation where tenant_id=$1`,
		`delete from inventory.warehouse_activity_line where tenant_id=$1`,
		`delete from inventory.warehouse_activity where tenant_id=$1`,
		`delete from inventory.warehouse_receipt_line where tenant_id=$1`,
		`delete from inventory.warehouse_receipt where tenant_id=$1`,
		`delete from inventory.bulk_uom_conversion_line where tenant_id=$1`,
		`delete from inventory.bulk_uom_conversion where tenant_id=$1`,
		`delete from inventory.bulk_cost_application where tenant_id=$1`,
		`delete from inventory.bulk_cost_layer where tenant_id=$1`,
		`delete from inventory.bulk_inventory_entry where tenant_id=$1`,
		`delete from inventory.bulk_reservation where tenant_id=$1`,
		`delete from inventory.bulk_balance where tenant_id=$1`,
		`delete from inventory.sales_warehouse_binding where tenant_id=$1`,
		`delete from inventory.item_bin_policy where tenant_id=$1`,
		`delete from inventory.warehouse_bin where tenant_id=$1`,
		`delete from inventory.item_unit_of_measure where tenant_id=$1`,
		`delete from inventory.stock_item where tenant_id=$1`,
		`delete from sales.customer_order_line where tenant_id=$1`,
		`delete from sales.customer_order where tenant_id=$1`,
		`delete from catalog.vehicle_variant where tenant_id=$1`,
		`delete from catalog.vehicle_model where tenant_id=$1`,
		`delete from platform.outbox_event where tenant_id=$1`,
		`delete from org.organization where tenant_id=$1`,
		`delete from platform.tenant where tenant_id=$1`,
		`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`,
		`alter table inventory.bulk_inventory_entry enable trigger bulk_inventory_entry_immutable`,
		`alter table inventory.bulk_cost_application enable trigger bulk_cost_application_immutable`,
		`alter table inventory.warehouse_receipt enable trigger warehouse_receipt_immutable`,
		`alter table inventory.warehouse_receipt_line enable trigger warehouse_receipt_line_immutable`,
		`alter table inventory.warehouse_activity_line enable trigger warehouse_activity_line_immutable`,
		`alter table inventory.warehouse_activity_uom_reservation enable trigger warehouse_uom_reservation_delete_guard`,
		`alter table inventory.sales_warehouse_binding enable trigger sales_warehouse_binding_immutable`,
		`alter table inventory.warehouse_pick_request enable trigger warehouse_pick_request_immutable`,
		`alter table sales.customer_shipment enable trigger customer_shipment_immutable`,
		`alter table sales.customer_shipment_line enable trigger customer_shipment_line_immutable`,
		`alter table inventory.customer_shipment_allocation enable trigger customer_shipment_allocation_immutable`,
	}
	for _, command := range commands {
		arguments := []any{}
		if strings.Contains(command, "$1") {
			arguments = append(arguments, tenant)
		}
		if _, err = tx.Exec(ctx, command, arguments...); err != nil {
			t.Errorf("sales warehouse cleanup: %v", err)
			return
		}
	}
	if err = tx.Commit(ctx); err != nil {
		t.Errorf("sales warehouse cleanup commit: %v", err)
	}
}
````

### FILE: `db/migrations/0037_sales_order_warehouse_demand.up.sql`

```yaml
block_id: "GO-SUPPLY-SALES-WAREHOUSE-MIGRATION-UP:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps 2eae56d704a1fd035d104f333602aea7091b7749 source document, sales line and outstanding quantity invariants"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "e2f6be2618bd1970f84363938306dbf7224b9c568912db1e0e76d09380cc5eba"
variables: []
secrets_allowed: false
```

````sql
begin;

create table inventory.sales_warehouse_binding (
  tenant_id uuid not null,
  request_id text not null,
  variant_id text not null,
  item_id text not null,
  sales_uom_code text not null,
  version bigint not null default 1 check (version = 1),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, variant_id),
  foreign key (tenant_id, variant_id)
    references catalog.vehicle_variant(tenant_id, variant_id),
  foreign key (tenant_id, item_id, sales_uom_code)
    references inventory.item_unit_of_measure(tenant_id, item_id, uom_code),
  unique (tenant_id, request_id),
  check (length(request_id) between 1 and 128)
);

create index sales_warehouse_binding_item_idx
  on inventory.sales_warehouse_binding(tenant_id, item_id, sales_uom_code);

create function inventory.reject_sales_warehouse_binding_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000',
    message='sales warehouse binding is immutable; create a new sellable variant for a changed item or unit contract';
end;
$function$;

create trigger sales_warehouse_binding_immutable
before update or delete on inventory.sales_warehouse_binding
for each row execute function inventory.reject_sales_warehouse_binding_mutation();

create table inventory.warehouse_pick_request (
  tenant_id uuid not null,
  request_id text not null,
  activity_id text not null,
  organization_id text not null,
  ship_bin_id text not null,
  item_id text not null,
  demand_kind text not null,
  demand_id text not null,
  demand_line_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  requested_handling_uom text not null default '',
  resolved_handling_uom text not null,
  allow_breakbulk boolean not null,
  use_fefo boolean not null,
  allow_dedicated boolean not null,
  assigned_to text not null default '',
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, request_id),
  foreign key (tenant_id, activity_id, organization_id)
    references inventory.warehouse_activity(tenant_id, activity_id, organization_id),
  foreign key (tenant_id, organization_id, ship_bin_id)
    references inventory.warehouse_bin(tenant_id, organization_id, bin_id),
  foreign key (tenant_id, item_id, resolved_handling_uom)
    references inventory.item_unit_of_measure(tenant_id, item_id, uom_code),
  unique (tenant_id, activity_id),
  check (length(request_id) between 1 and 128)
);

create function inventory.reject_warehouse_pick_request_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='warehouse pick request evidence is immutable';
end;
$function$;

create trigger warehouse_pick_request_immutable
before update or delete on inventory.warehouse_pick_request
for each row execute function inventory.reject_warehouse_pick_request_mutation();

commit;
````

### FILE: `db/migrations/0037_sales_order_warehouse_demand.down.sql`

```yaml
block_id: "GO-SUPPLY-SALES-WAREHOUSE-MIGRATION-DOWN:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps 2eae56d704a1fd035d104f333602aea7091b7749 source document, sales line and outstanding quantity invariants"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "5b99dce01ee7ca7c906c08ef522a7c437024803b98e0c35a15d8572e1704d7c1"
variables: []
secrets_allowed: false
```

````sql
begin;

do $block$
begin
  if exists (select 1 from inventory.sales_warehouse_binding) or
     exists (select 1 from inventory.warehouse_pick_request) then
    raise exception using errcode='55000',
      message='cannot remove connected sales warehouse schema after binding or request evidence exists';
  end if;
end;
$block$;

drop trigger warehouse_pick_request_immutable on inventory.warehouse_pick_request;
drop function inventory.reject_warehouse_pick_request_mutation();
drop table inventory.warehouse_pick_request;

drop trigger sales_warehouse_binding_immutable on inventory.sales_warehouse_binding;
drop function inventory.reject_sales_warehouse_binding_mutation();
drop table inventory.sales_warehouse_binding;

commit;
````

### FILE: `docs/inventory/MICROSOFT_BC_SALES_WAREHOUSE_DERIVATION.md`

```yaml
block_id: "GO-SUPPLY-SALES-WAREHOUSE-DERIVATION:v1"
operation: CREATE
provenance: AUTHORED
source: "exact local derivation record for microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749"
license: "LicenseRef-Workspace-Owner"
sha256: "98750c5e40c72ef804a95a22dd9f5771af03b12924be75284999c8f18d5de3d9"
variables: []
secrets_allowed: false
```

````markdown
# Microsoft BC sales-order to warehouse derivation

This connected demand lane is an `ADAPTED` Go/PostgreSQL implementation. It is not Microsoft-authored Go and it is not a Business Central runtime.

Exact authority is the MIT-licensed `microsoft/BCApps` commit `2eae56d704a1fd035d104f333602aea7091b7749`, tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`:

- `WhseCreateSourceDocument.Codeunit.al`, SHA-256 `df2a5f4767b58d4d307a223dcfe1c3d0176c9d336d6f9ebdd25ec414f40672f8`, sets quantity and base quantity, initializes outstanding quantities and checks the source document line and bin.
- `SalesWhsePostShipment.Codeunit.al`, SHA-256 `84121bd38555f3a7cba59dd97f985ae41af2da349359a60217bcab14ee7c3747`, resolves the sales document and line and reconciles quantity-to-ship/base quantity.
- `SalesLine.Table.al`, SHA-256 `ca65615dc05bc555a878ba2635c955897930332ef171e313a00103074fbc7278`, owns variant, UOM, quantity, outstanding quantity and their base forms.
- `SCMWarehousePick.Codeunit.al`, SHA-256 `f856e941d86a78a09d233aba6dfe5470dacbf08c853c26bc5e26512c3d6881c5`, includes released sales-order, partial-pick, registration and multiple-pick scenarios.
- `SalesOrderWhseValidateLine.Codeunit.al`, SHA-256 `12b4a3ce4fdef9592aad3054f5f174c2c4493d8377216cd4bea6890289d66deb`, verifies warehouse coupling when variant or UOM changes.

The portable contract therefore requires an immutable variant-to-item/sales-UOM binding, a placed operational sales line, exact organization/item/source identity, base-quantity conversion using the fixed item UOM, subtraction of every open or registered pick, atomic refusal beyond outstanding quantity, exact request replay, divergent replay refusal and durable pick evidence. Shipment and sale completion are not claimed by this lane.
````

### FILE: `internal/platform/postgres/sales_shipment.go`

```yaml
block_id: "GO-SUPPLY-CUSTOMER-SHIPMENT-POSTGRES:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps 2eae56d704a1fd035d104f333602aea7091b7749 posted sales shipment and warehouse posting invariants"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "e865843d1d03f3cb5926547f100f4e98a021e6d79b68a19ae07cfd6a528f7a30"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
)

type customerShipmentAllocation struct {
	reservationID      string
	binID              string
	lotID              string
	quantityBase       string
	reservationVersion int64
	costs              []transferCost
	costAmount         string
}

func readCustomerShipmentReplay(ctx context.Context, tx pgx.Tx, tenant string, command inventorycontrol.CustomerShipmentCommand) (inventorycontrol.CustomerShipment, bool, error) {
	var value inventorycontrol.CustomerShipment
	var postingDate time.Time
	err := tx.QueryRow(ctx, `select s.shipment_id,s.request_id,s.organization_id,s.order_id,s.warehouse_activity_id,s.posting_date,s.order_version,s.fulfillment_state,l.shipment_line_id,l.order_line_id,l.variant_id,l.item_id,l.sales_uom_code,l.quantity::text,l.quantity_base::text,l.cost_amount::text,(select count(*) from inventory.customer_shipment_allocation a where a.tenant_id=s.tenant_id and a.shipment_id=s.shipment_id) from sales.customer_shipment s join sales.customer_shipment_line l using(tenant_id,shipment_id) where s.tenant_id=$1 and s.request_id=$2 for share of s,l`, tenant, command.RequestID).Scan(&value.ID, &value.RequestID, &value.OrganizationID, &value.OrderID, &value.WarehouseActivityID, &postingDate, &value.OrderVersion, &value.FulfillmentState, &value.Line.ID, &value.Line.OrderLineID, &value.Line.VariantID, &value.Line.ItemID, &value.Line.SalesUOMCode, &value.Line.Quantity, &value.Line.QuantityBase, &value.Line.CostAmount, &value.Line.AllocationCount)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, false, nil
	}
	if err != nil {
		return value, false, err
	}
	if value.OrganizationID != command.OrganizationID || value.OrderID != command.OrderID || value.WarehouseActivityID != command.WarehouseActivityID || !postingDate.Equal(command.PostingDate) {
		return inventorycontrol.CustomerShipment{}, false, inventorycontrol.ErrConflict
	}
	value.PostingDate = postingDate
	return value, true, nil
}

func (r *InventoryControl) PostCustomerShipment(ctx context.Context, tenant string, ids inventorycontrol.CustomerShipmentIDs, command inventorycontrol.CustomerShipmentCommand) (inventorycontrol.CustomerShipment, error) {
	value := inventorycontrol.CustomerShipment{ID: ids.ShipmentID, RequestID: command.RequestID, OrganizationID: command.OrganizationID, OrderID: command.OrderID, WarehouseActivityID: command.WarehouseActivityID, PostingDate: command.PostingDate, Line: inventorycontrol.CustomerShipmentLine{ID: ids.LineID}}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+"\x1fcustomer-shipment-request\x1f"+command.RequestID); err != nil {
		return value, err
	}
	if replay, found, replayErr := readCustomerShipmentReplay(ctx, tx, tenant, command); replayErr != nil {
		return value, replayErr
	} else if found {
		if err = tx.Commit(ctx); err != nil {
			return value, err
		}
		return replay, nil
	}
	var orderLineID, variantID, itemID, salesUOM, salesFactorText, orderedText, shippedBeforeText, costing string
	var orderVersion int64
	err = tx.QueryRow(ctx, `select l.line_id,l.variant_id,b.item_id,b.sales_uom_code,u.qty_per_uom::text,l.quantity::numeric::text,l.shipped_quantity::text,i.costing_method,o.version from sales.customer_order o join inventory.warehouse_activity a on a.tenant_id=o.tenant_id and a.organization_id=o.organization_id and a.source_kind='customer-order' and a.source_id=o.order_id join sales.customer_order_line l on l.tenant_id=o.tenant_id and l.order_id=o.order_id and l.line_id=a.source_line_id join inventory.sales_warehouse_binding b on b.tenant_id=l.tenant_id and b.variant_id=l.variant_id join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.uom_code=b.sales_uom_code join inventory.stock_item i on i.tenant_id=b.tenant_id and i.item_id=b.item_id where o.tenant_id=$1 and o.organization_id=$2 and o.order_id=$3 and a.activity_id=$4 and a.activity_type='pick' and a.status='registered' and o.state in ('placed','confirmed','paid','allocated') and o.fulfillment_state in ('unfulfilled','partially-shipped') for update of o,l,a`, tenant, command.OrganizationID, command.OrderID, command.WarehouseActivityID).Scan(&orderLineID, &variantID, &itemID, &salesUOM, &salesFactorText, &orderedText, &shippedBeforeText, &costing, &orderVersion)
	if errors.Is(err, pgx.ErrNoRows) || costing != "fifo" {
		return value, inventorycontrol.ErrConflict
	}
	if err != nil {
		return value, err
	}
	if err = warehouseLock(ctx, tx, tenant, command.OrganizationID, itemID); err != nil {
		return value, err
	}
	rows, err := tx.Query(ctx, `select r.reservation_id,r.bin_id,coalesce(r.lot_id,''),r.quantity::text,r.version from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) join inventory.bulk_reservation r on r.tenant_id=al.tenant_id and r.reservation_id=al.reservation_id join inventory.warehouse_bin w on w.tenant_id=r.tenant_id and w.organization_id=r.organization_id and w.bin_id=r.bin_id where a.tenant_id=$1 and a.organization_id=$2 and a.activity_id=$3 and a.activity_type='pick' and a.status='registered' and a.source_kind='customer-order' and a.source_id=$4 and a.source_line_id=$5 and r.item_id=$6 and r.demand_kind='customer-order' and r.demand_id=$4 and r.demand_line_id=a.source_line_id||':'||a.activity_id||':'||lpad(al.sequence_no::text,6,'0') and r.status='reservation' and w.bin_type='ship' and not w.movement_blocked order by al.sequence_no for update of r`, tenant, command.OrganizationID, command.WarehouseActivityID, command.OrderID, orderLineID, itemID)
	if err != nil {
		return value, err
	}
	allocations := []customerShipmentAllocation{}
	shipmentBase := new(big.Rat)
	for rows.Next() {
		var allocation customerShipmentAllocation
		if err = rows.Scan(&allocation.reservationID, &allocation.binID, &allocation.lotID, &allocation.quantityBase, &allocation.reservationVersion); err != nil {
			rows.Close()
			return value, err
		}
		quantity, ok := parseRat(allocation.quantityBase)
		if !ok {
			rows.Close()
			return value, inventorycontrol.ErrConflict
		}
		shipmentBase.Add(shipmentBase, quantity)
		allocations = append(allocations, allocation)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return value, err
	}
	if len(allocations) == 0 || shipmentBase.Sign() <= 0 {
		return value, inventorycontrol.ErrConflict
	}
	salesFactor, factorOK := parseRat(salesFactorText)
	ordered, orderedOK := parseRat(orderedText)
	shippedBefore, shippedOK := parseRat(shippedBeforeText)
	if !factorOK || !orderedOK || !shippedOK || salesFactor.Sign() <= 0 {
		return value, inventorycontrol.ErrConflict
	}
	shipmentSalesRaw := new(big.Rat).Quo(shipmentBase, salesFactor)
	shipmentSalesText := formatRat(shipmentSalesRaw, 6)
	shipmentSales, salesOK := parseRat(shipmentSalesText)
	if !salesOK || new(big.Rat).Mul(shipmentSales, salesFactor).Cmp(shipmentBase) != 0 || new(big.Rat).Add(shippedBefore, shipmentSales).Cmp(ordered) > 0 {
		return value, inventorycontrol.ErrConflict
	}
	totalCost := new(big.Rat)
	for index := range allocations {
		costs, cost, costErr := transferCosts(ctx, tx, tenant, command.OrganizationID, itemID, allocations[index].lotID, allocations[index].quantityBase, "")
		if costErr != nil {
			return value, costErr
		}
		allocations[index].costs = costs
		allocations[index].costAmount = formatRat(cost, 4)
		totalCost.Add(totalCost, cost)
		for _, component := range costs {
			updated, updateErr := tx.Exec(ctx, `update inventory.bulk_cost_layer set remaining_quantity=remaining_quantity-$3::numeric where tenant_id=$1 and layer_id=$2 and remaining_quantity >= $3::numeric`, tenant, component.layerID, component.quantity)
			if updateErr != nil || updated.RowsAffected() != 1 {
				if updateErr != nil {
					return value, updateErr
				}
				return value, inventorycontrol.ErrConflict
			}
		}
	}
	for index, allocation := range allocations {
		outboundID := fmt.Sprintf("%s-entry-%06d", ids.ShipmentID, index+1)
		unitCost := formatRat(new(big.Rat).Quo(mustRat(allocation.costAmount), mustRat(allocation.quantityBase)), 4)
		_, err = tx.Exec(ctx, `insert into inventory.bulk_inventory_entry(tenant_id,entry_id,organization_id,bin_id,item_id,lot_id,entry_type,quantity,unit_cost,cost_amount,posting_date,source_kind,source_id) values($1,$2,$3,$4,$5,nullif($6,''),'issue',-$7::numeric,$8::numeric,-$9::numeric,$10::date,'customer-shipment',$11)`, tenant, outboundID, command.OrganizationID, allocation.binID, itemID, allocation.lotID, allocation.quantityBase, unitCost, allocation.costAmount, command.PostingDate, fmt.Sprintf("%s:%06d", ids.ShipmentID, index+1))
		if err != nil {
			return value, bulkConflict(err)
		}
		for _, component := range allocation.costs {
			if _, err = tx.Exec(ctx, `insert into inventory.bulk_cost_application(tenant_id,outbound_entry_id,inbound_entry_id,quantity,cost_amount) values($1,$2,$3,$4::numeric,$5::numeric)`, tenant, outboundID, component.inboundEntryID, component.quantity, component.costAmount); err != nil {
				return value, bulkConflict(err)
			}
		}
		if err = consumeRegisteredPickPackaging(ctx, tx, tenant, command.OrganizationID, command.WarehouseActivityID, allocation.reservationID); err != nil {
			return value, err
		}
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_balance set quantity=quantity-$6::numeric,reserved_quantity=reserved_quantity-$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and quantity >= $6::numeric and reserved_quantity >= $6::numeric`, tenant, command.OrganizationID, allocation.binID, itemID, allocation.lotID, allocation.quantityBase)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return value, updateErr
			}
			return value, inventorycontrol.ErrConflict
		}
		updated, updateErr = tx.Exec(ctx, `update inventory.bulk_reservation set status='consumed',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and reservation_id=$2 and status='reservation' and version=$3`, tenant, allocation.reservationID, allocation.reservationVersion)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return value, updateErr
			}
			return value, inventorycontrol.ErrConflict
		}
	}
	shippedAfter := new(big.Rat).Add(shippedBefore, shipmentSales)
	updated, err := tx.Exec(ctx, `update sales.customer_order_line set shipped_quantity=shipped_quantity+$4::numeric where tenant_id=$1 and order_id=$2 and line_id=$3 and shipped_quantity+$4::numeric<=quantity`, tenant, command.OrderID, orderLineID, shipmentSalesText)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return value, err
		}
		return value, inventorycontrol.ErrConflict
	}
	fulfillmentState := "partially-shipped"
	var incomplete bool
	if err = tx.QueryRow(ctx, `select exists(select 1 from sales.customer_order_line where tenant_id=$1 and order_id=$2 and shipped_quantity<quantity)`, tenant, command.OrderID).Scan(&incomplete); err != nil {
		return value, err
	}
	if !incomplete && shippedAfter.Cmp(ordered) == 0 {
		fulfillmentState = "shipped"
	}
	updated, err = tx.Exec(ctx, `update sales.customer_order set fulfillment_state=$4,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and order_id=$2 and version=$3 and fulfillment_state in ('unfulfilled','partially-shipped')`, tenant, command.OrderID, orderVersion, fulfillmentState)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return value, err
		}
		return value, inventorycontrol.ErrConflict
	}
	value.OrderVersion = orderVersion + 1
	value.FulfillmentState = fulfillmentState
	value.Line = inventorycontrol.CustomerShipmentLine{ID: ids.LineID, OrderLineID: orderLineID, VariantID: variantID, ItemID: itemID, SalesUOMCode: salesUOM, Quantity: shipmentSalesText, QuantityBase: formatRat(shipmentBase, 6), CostAmount: formatRat(totalCost, 4), AllocationCount: len(allocations)}
	_, err = tx.Exec(ctx, `insert into sales.customer_shipment(tenant_id,shipment_id,request_id,organization_id,order_id,warehouse_activity_id,posting_date,total_cost_amount,order_version,fulfillment_state) values($1,$2,$3,$4,$5,$6,$7::date,$8::numeric,$9,$10)`, tenant, value.ID, value.RequestID, value.OrganizationID, value.OrderID, value.WarehouseActivityID, value.PostingDate, value.Line.CostAmount, value.OrderVersion, value.FulfillmentState)
	if err != nil {
		return value, bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into sales.customer_shipment_line(tenant_id,shipment_id,shipment_line_id,order_id,order_line_id,variant_id,item_id,sales_uom_code,quantity,quantity_base,cost_amount) values($1,$2,$3,$4,$5,$6,$7,$8,$9::numeric,$10::numeric,$11::numeric)`, tenant, value.ID, value.Line.ID, value.OrderID, value.Line.OrderLineID, value.Line.VariantID, value.Line.ItemID, value.Line.SalesUOMCode, value.Line.Quantity, value.Line.QuantityBase, value.Line.CostAmount)
	if err != nil {
		return value, bulkConflict(err)
	}
	for index, allocation := range allocations {
		outboundID := fmt.Sprintf("%s-entry-%06d", ids.ShipmentID, index+1)
		allocationID := fmt.Sprintf("%s-allocation-%06d", ids.ShipmentID, index+1)
		_, err = tx.Exec(ctx, `insert into inventory.customer_shipment_allocation(tenant_id,shipment_id,shipment_line_id,allocation_id,reservation_id,outbound_entry_id,organization_id,from_bin_id,item_id,lot_id,quantity_base,cost_amount) values($1,$2,$3,$4,$5,$6,$7,$8,$9,nullif($10,''),$11::numeric,$12::numeric)`, tenant, value.ID, value.Line.ID, allocationID, allocation.reservationID, outboundID, value.OrganizationID, allocation.binID, itemID, allocation.lotID, allocation.quantityBase, allocation.costAmount)
		if err != nil {
			return value, bulkConflict(err)
		}
	}
	if err = recordBulkEvent(ctx, tx, tenant, ids.EventID, "customer-shipment", value.ID, "customer-shipment.posted", value.OrderVersion, map[string]any{"request_id": value.RequestID, "organization_id": value.OrganizationID, "order_id": value.OrderID, "order_line_id": value.Line.OrderLineID, "warehouse_activity_id": value.WarehouseActivityID, "quantity": value.Line.Quantity, "quantity_base": value.Line.QuantityBase, "sales_uom_code": value.Line.SalesUOMCode, "cost_amount": value.Line.CostAmount, "fulfillment_state": value.FulfillmentState, "allocations": value.Line.AllocationCount}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}
````

### FILE: `db/migrations/0038_customer_sales_shipment.up.sql`

```yaml
block_id: "GO-SUPPLY-CUSTOMER-SHIPMENT-MIGRATION-UP:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps 2eae56d704a1fd035d104f333602aea7091b7749 posted sales shipment header, line and shipped-quantity invariants"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "98b02cebe2f8bb9e01870cbe87a5bf5e4ee3a11c94b5a1f7cc6668e88c035da8"
variables: []
secrets_allowed: false
```

````sql
begin;

alter table sales.customer_order
  add column fulfillment_state text not null default 'unfulfilled'
    check (fulfillment_state in ('unfulfilled','partially-shipped','shipped','delivered'));

alter table sales.customer_order_line
  add column shipped_quantity numeric(20,6) not null default 0
    check (shipped_quantity >= 0 and shipped_quantity <= quantity);

create table sales.customer_shipment (
  tenant_id uuid not null,
  shipment_id text not null,
  request_id text not null,
  organization_id text not null,
  order_id text not null,
  warehouse_activity_id text not null,
  posting_date date not null,
  total_cost_amount numeric(20,4) not null check (total_cost_amount >= 0),
  order_version bigint not null check (order_version > 0),
  fulfillment_state text not null check (fulfillment_state in ('partially-shipped','shipped')),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, shipment_id),
  unique (tenant_id, request_id),
  unique (tenant_id, warehouse_activity_id),
  foreign key (tenant_id, organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id, order_id) references sales.customer_order(tenant_id,order_id),
  foreign key (tenant_id, warehouse_activity_id, organization_id)
    references inventory.warehouse_activity(tenant_id,activity_id,organization_id),
  check (length(request_id) between 1 and 128)
);

create table sales.customer_shipment_line (
  tenant_id uuid not null,
  shipment_id text not null,
  shipment_line_id text not null,
  order_id text not null,
  order_line_id text not null,
  variant_id text not null,
  item_id text not null,
  sales_uom_code text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  quantity_base numeric(20,6) not null check (quantity_base > 0),
  cost_amount numeric(20,4) not null check (cost_amount >= 0),
  primary key (tenant_id, shipment_id, shipment_line_id),
  foreign key (tenant_id, shipment_id) references sales.customer_shipment(tenant_id,shipment_id),
  foreign key (tenant_id, order_id, order_line_id) references sales.customer_order_line(tenant_id,order_id,line_id),
  foreign key (tenant_id, variant_id) references catalog.vehicle_variant(tenant_id,variant_id),
  foreign key (tenant_id, item_id, sales_uom_code) references inventory.item_unit_of_measure(tenant_id,item_id,uom_code),
  check (quantity_base > 0)
);

create table inventory.customer_shipment_allocation (
  tenant_id uuid not null,
  shipment_id text not null,
  shipment_line_id text not null,
  allocation_id text not null,
  reservation_id text not null,
  outbound_entry_id text not null,
  organization_id text not null,
  from_bin_id text not null,
  item_id text not null,
  lot_id text,
  quantity_base numeric(20,6) not null check (quantity_base > 0),
  cost_amount numeric(20,4) not null check (cost_amount >= 0),
  primary key (tenant_id, shipment_id, allocation_id),
  unique (tenant_id, reservation_id),
  unique (tenant_id, outbound_entry_id),
  foreign key (tenant_id, shipment_id, shipment_line_id)
    references sales.customer_shipment_line(tenant_id,shipment_id,shipment_line_id),
  foreign key (tenant_id, reservation_id) references inventory.bulk_reservation(tenant_id,reservation_id),
  foreign key (tenant_id, outbound_entry_id) references inventory.bulk_inventory_entry(tenant_id,entry_id),
  foreign key (tenant_id, organization_id, from_bin_id)
    references inventory.warehouse_bin(tenant_id,organization_id,bin_id),
  foreign key (tenant_id, lot_id) references inventory.inventory_lot(tenant_id,lot_id)
);

create function sales.reject_customer_shipment_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='posted customer shipment evidence is immutable';
end;
$function$;

create trigger customer_shipment_immutable
before update or delete on sales.customer_shipment
for each row execute function sales.reject_customer_shipment_mutation();

create trigger customer_shipment_line_immutable
before update or delete on sales.customer_shipment_line
for each row execute function sales.reject_customer_shipment_mutation();

create trigger customer_shipment_allocation_immutable
before update or delete on inventory.customer_shipment_allocation
for each row execute function sales.reject_customer_shipment_mutation();

create index customer_shipment_order_idx
  on sales.customer_shipment(tenant_id,order_id,posting_date,shipment_id);

commit;
````

### FILE: `db/migrations/0038_customer_sales_shipment.down.sql`

```yaml
block_id: "GO-SUPPLY-CUSTOMER-SHIPMENT-MIGRATION-DOWN:v1"
operation: CREATE
provenance: ADAPTED
source: "microsoft/BCApps 2eae56d704a1fd035d104f333602aea7091b7749 posted sales shipment header, line and shipped-quantity invariants"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "4f0ff8f7459bfe46075ee47dd56881b07ad25710403ce449698b392c44cf421e"
variables: []
secrets_allowed: false
```

````sql
begin;

drop index if exists sales.customer_shipment_order_idx;
drop trigger if exists customer_shipment_allocation_immutable on inventory.customer_shipment_allocation;
drop trigger if exists customer_shipment_line_immutable on sales.customer_shipment_line;
drop trigger if exists customer_shipment_immutable on sales.customer_shipment;
drop function if exists sales.reject_customer_shipment_mutation();
drop table if exists inventory.customer_shipment_allocation;
drop table if exists sales.customer_shipment_line;
drop table if exists sales.customer_shipment;
alter table sales.customer_order_line drop column if exists shipped_quantity;
alter table sales.customer_order drop column if exists fulfillment_state;

commit;
````

### FILE: `db/tests/0038_customer_sales_shipment.test.sql`

```yaml
block_id: "GO-SUPPLY-CUSTOMER-SHIPMENT-SCHEMA-TEST:v1"
operation: CREATE
provenance: ADAPTED
source: "executable local verification of the fixed microsoft/BCApps shipment adaptation"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "0e5e830cdb38825b6c5464e9b75aaaf80703bd2d43ac3c71758e503074d66215"
variables: []
secrets_allowed: false
```

````sql
do $$
begin
  if to_regclass('sales.customer_shipment') is null then raise exception 'customer shipment missing'; end if;
  if to_regclass('sales.customer_shipment_line') is null then raise exception 'customer shipment line missing'; end if;
  if to_regclass('inventory.customer_shipment_allocation') is null then raise exception 'customer shipment allocation missing'; end if;
  if not exists(select 1 from information_schema.columns where table_schema='sales' and table_name='customer_order' and column_name='fulfillment_state') then raise exception 'order fulfillment state missing'; end if;
  if not exists(select 1 from information_schema.columns where table_schema='sales' and table_name='customer_order_line' and column_name='shipped_quantity') then raise exception 'order shipped quantity missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='customer_shipment_immutable') then raise exception 'customer shipment immutability missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='customer_shipment_line_immutable') then raise exception 'customer shipment line immutability missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='customer_shipment_allocation_immutable') then raise exception 'customer shipment allocation immutability missing'; end if;
  if not exists(select 1 from pg_indexes where schemaname='sales' and indexname='customer_shipment_order_idx') then raise exception 'customer shipment order index missing'; end if;
end $$;
````

### FILE: `docs/inventory/MICROSOFT_BC_CUSTOMER_SHIPMENT_DERIVATION.md`

```yaml
block_id: "GO-SUPPLY-CUSTOMER-SHIPMENT-DERIVATION:v1"
operation: CREATE
provenance: AUTHORED
source: "exact local derivation record for microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749"
license: "LicenseRef-Workspace-Owner"
sha256: "93a5c53d1d620aa7bec0401c0115cd1c80c323b0825924b84a7cd3216509cff2"
variables: []
secrets_allowed: false
```

````markdown
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
````

## 6. Configuration surface

No independent environment variables are introduced. The module consumes the composed database, ID source and verifier. Organization, supplier/factory, state/version, serial identity, item code/base UOM/base rounding precision/tracking mode, immutable alternate UOM code/factor/precision, explicit handling-unit operation/request/source/target UOM, receipt/pick/replenishment handling UOM and `allow_breakbulk`, bin type/ranking/capacity/cross-dock marker, cross-dock due-date horizon, lot/expiry, request/source identity, activity assignment/version, reservation demand, FEFO choice, transfer route/in-transit code, dual-organization scope, quantity and posting date are strict domain inputs. La conexión comercial agrega como configuración obligatoria el binding inmutable variante→item/UOM comercial y exige identidades exactas de pedido/línea/request. Residual UOM and packaging conversions fail closed. RECEIVE/SHIP/PUT AWAY/PICK/PUTPICK/QC behavior is closed in code; invalid combinations fail before effect. Each shipment additionally requires an exact registered pick ID and each shipment/receipt requires a posting request ID plus positive explicit quantity; these are domain inputs, never environment defaults. Ordinary bulk issues and cross-organization transfers admit `fifo` or explicit `specific`; a specific-cost transfer requires `specific_receipt_entry`, persists that selector on the line, rejects a changed idempotent replay and consumes only that exact source layer. FIFO rejects a selector. `lifo`, `average` and `standard` fail closed. Accounting policy, currency/rounding, landed cost, closed periods, QC disposition, title, Incoterm, cross-dock working-day calendar and external ERP/MES/WMS mappings require project-owned configuration and approval. Portable ATP remains limited to serialized sellable stock, due inbound transfers and unallocated placed demand; full BC planning/CTP remains outside this pack.

V173 añade como inputs estrictos `customer-shipment.request_id`, organización, pedido, pick registrado y fecha de posting. La política de crédito/pago, el carrier, la referencia logística, la entrega al cliente y el evento fiscal no reciben defaults: se configuran o integran en sus owners posteriores. El shipment bulk comercial portable admite FIFO; costo específico y stock serial fallan cerrados hasta poseer su vínculo exacto de fuente/activo.

## 7. Dependency bill

| Dependency | Pin | Use | License | Scope | Source |
|---|---|---|---|---|---|
| Go | `1.26.7` | domain/HTTP/tests | BSD-3-Clause | build/runtime | `go.dev` |
| PostgreSQL | `18.6` verified baseline | durable operations/outbox | PostgreSQL | runtime/test | `postgresql.org` |
| `pgx` | `5.10.0` | PostgreSQL adapter | MIT | build/runtime | `github.com/jackc/pgx` |
| Microsoft BCApps | commit `2eae56d704a1fd035d104f333602aea7091b7749` | reservation/ATP, tracking/bin/cost, receipt/activity/ranking/put-away/pick/FEFO, replenishment, cross-dock and posted sales shipment authority | MIT | design/adaptation | `github.com/microsoft/BCApps` |

## 8. Apply order

Compose after foundation, electromobility schema, commerce order lines and CRM types; then wire the single serial+bulk+warehouse+transfer inventory owner. Apply migrations through 0038 and execute state, cumulative quantity, posting identity/idempotency, cost-layer slicing, receipt, bin eligibility/ranking/capacity, FEFO, Take/Place, concurrent pick/shipment/replenishment, min/max and pending movement, transfer create/replay/partial ship/in-transit/partial receive/put-away, cross-dock horizon/opportunity/capacity/availability/prioritization/cancellation, base/alternate UOM uniqueness/precision/residual/immutability, packaging conservation/replay/divergence/concurrency/evidence/legacy upgrade, receipt/pick/replenishment UOM preservation or authorized breakbulk, physical+composition reservation, customer-order source/binding/outstanding/replay/cancel/register/concurrency, registered-pick→customer-shipment partial/final posting, exact/divergent replay, one-winner concurrency, FIFO cost/UOM conservation, order-line shipped quantity and fulfillment state, duplicate identity, expiry, dual-organization scope, immutable-ledger and PostgreSQL integration tests. Existing ERP/MES/WMS systems require an ownership/cutover map. Rollback disables commands/workers while preserving operational and valuation history; migrations 0030–0038 down are development-only unless posting, movement, cross-dock, UOM, packaging, sales-demand, customer-shipment and open-activity reservation evidence preservation or recovery is proven.

## 9. Verification

V173 reconstruye el pack 0.16.0 en 88/88 y, desde PostgreSQL 18.6 limpio, pasa 38 migraciones, 24 pruebas SQL, ciclo 0038 down/up y toda la suite Go. El journey cliente postea primero 1 BOX/12 EA y después 2 BOX/24 EA con costo FIFO total 360, exact replay, rechazo divergente y un solo ganador concurrente; termina con `shipped_quantity=3`, fulfillment `shipped`, reservas consumidas y balance físico/UOM/costo restante cero. La regresión transfer recibe BOX reales, autoriza breakbulk de una EA, conserva otra BOX intacta y demuestra que shipment consume también composición no-base. La composición canónica 29 packs / 380 archivos, vet, build y gates raíz se registran en la evidencia V173 antes de promoción. Producción aún exige carrier/transporte conectado, entrega/handover, facturación/fiscal/pago-crédito, customer shipment serial/específico, demanda service/production, carga representativa, políticas contables/landed-cost/períodos aprobadas, calendarios, planning/CTP, reconciliación ERP/WMS, operator tooling y aceptación real.

## 10. Reconstruction evidence

La evidencia histórica es `reconstruction_evidence/GO_SUPPLY_FACTORY_INVENTORY_API_2026-08-24_V1.md`; V159–V167 registran la evolución hasta UOM 0.11.0. `reconstruction_evidence/MICROSOFT_BC_HANDLING_UNIT_COMPOSITION_2026-09-01_V168.md` gobierna 0.12.0, V169 gobierna 0.13.0 y `reconstruction_evidence/MICROSOFT_BC_PICK_REPLENISHMENT_PACKAGING_2026-09-01_V171.md` gobierna 0.14.0. `reconstruction_evidence/MICROSOFT_BC_SALES_ORDER_WAREHOUSE_DEMAND_2026-09-01_V172.md` gobierna 0.15.0. `reconstruction_evidence/MICROSOFT_BC_CUSTOMER_SALES_SHIPMENT_2026-09-01_V173.md` gobierna 0.16.0 sólo después de reconstrucción Markdown, PostgreSQL limpio, 0038 down/up, Go/vet/build, composición y gates raíz.
