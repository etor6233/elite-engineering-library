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
