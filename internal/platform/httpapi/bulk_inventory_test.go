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
