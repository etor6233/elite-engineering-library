package httpapi

import (
	"context"
	"elite.local/enterprise/internal/fulfillment"
	"elite.local/enterprise/internal/platform/identity"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fulfillmentRepository struct{ shipments, transports, reports, organizations int }

func (f *fulfillmentRepository) CreateShipment(context.Context, string, string, fulfillment.Shipment) error {
	f.shipments++
	return nil
}
func (*fulfillmentRepository) TransitionShipment(context.Context, string, string, string, string, string, string, string, string) error {
	return nil
}
func (f *fulfillmentRepository) CreateCustomerTransport(_ context.Context, _ string, ids fulfillment.TransportIDs, value fulfillment.CustomerTransportCommand) (fulfillment.CustomerTransport, error) {
	f.transports++
	return fulfillment.CustomerTransport{ID: ids.ShipmentID, RequestID: value.RequestID, State: "booked", Version: 1}, nil
}
func (f *fulfillmentRepository) RecordProviderReport(_ context.Context, _ string, shipment, _ string, value fulfillment.ProviderReportCommand) (fulfillment.CustomerTransport, error) {
	f.reports++
	return fulfillment.CustomerTransport{ID: shipment, State: "delivered", Version: value.Version + 1}, nil
}
func (*fulfillmentRepository) OpenServiceCase(context.Context, string, string, fulfillment.ServiceCase) error {
	return nil
}
func (*fulfillmentRepository) TransitionServiceCase(context.Context, string, string, string, string, string, int64, string) error {
	return nil
}
func (*fulfillmentRepository) CreateRecall(context.Context, string, string, fulfillment.Recall) error {
	return nil
}
func (*fulfillmentRepository) ActivateRecall(context.Context, string, string, string) error {
	return nil
}
func (*fulfillmentRepository) AddRecallUnit(context.Context, string, string, string, string, string) error {
	return nil
}
func (f *fulfillmentRepository) CreateOrganization(context.Context, string, string, fulfillment.Organization) error {
	f.organizations++
	return nil
}
func (*fulfillmentRepository) TransitionOrganization(context.Context, string, string, string, string, int64, string) error {
	return nil
}
func (*fulfillmentRepository) CreateAgreement(context.Context, string, string, fulfillment.Agreement) error {
	return nil
}
func (*fulfillmentRepository) TransitionAgreement(context.Context, string, string, string, string, string, int64, string) error {
	return nil
}
func (*fulfillmentRepository) QueueMessage(context.Context, string, string, fulfillment.Message) error {
	return nil
}
func (*fulfillmentRepository) TransitionMessage(context.Context, string, string, string, string, string, string, string) error {
	return nil
}

type fulfillmentIDs struct{ n int }

func (i *fulfillmentIDs) New() string {
	i.n++
	return fmt.Sprintf("018f4d4a-7b36-7a21-8d10-%012x", i.n)
}

type fulfillmentVerifier struct{}

func (fulfillmentVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{Subject: "operator", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c28800", Permissions: map[string]struct{}{"logistics:write": {}, "network:admin": {}}, Organizations: map[string]struct{}{"a": {}, "b": {}}}, nil
}
func TestFulfillmentHTTPPermissionBoundaries(t *testing.T) {
	repo := &fulfillmentRepository{}
	service := fulfillment.NewService(repo, &fulfillmentIDs{})
	mux := http.NewServeMux()
	FulfillmentModule{Service: service}.Register(mux, fulfillmentVerifier{})
	request := httptest.NewRequest("POST", "/v1/logistics/shipments", strings.NewReader(`{"provider_code":"carrier","origin_organization_id":"a","destination_organization_id":"b"}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.shipments != 1 {
		t.Fatalf("shipment status=%d body=%s", response.Code, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/logistics/customer-transports", strings.NewReader(`{"request_id":"request-1","customer_shipment_id":"sales-shipment-1","origin_organization_id":"a","provider_code":"amazon-easyship","provider_reference":"tracking-1"}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.transports != 1 {
		t.Fatalf("transport status=%d body=%s", response.Code, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/logistics/shipments/transport-1/provider-reports", strings.NewReader(`{"origin_organization_id":"a","provider_event_id":"event-1","provider_status":"Delivered","report_schema":"elite-amazon-spapi-easyship-reconciliation-receipt/v1","evidence_sha256":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","occurred_at":"2026-09-02T12:00:00Z","version":1,"automatic_business_delivery_acceptance":false}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 200 || repo.reports != 1 {
		t.Fatalf("provider report status=%d body=%s", response.Code, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/franchise/network/organizations", strings.NewReader(`{"parent_organization_id":"a","code":"store-cordoba","display_name":"Store Cordoba","type":"store"}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.organizations != 1 {
		t.Fatalf("organization status=%d body=%s", response.Code, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/communications/messages", strings.NewReader(`{}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 {
		t.Fatalf("message permission status=%d", response.Code)
	}
}
