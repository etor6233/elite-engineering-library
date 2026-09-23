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
