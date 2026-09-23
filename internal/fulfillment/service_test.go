package fulfillment

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type fakeRepository struct{}

func (*fakeRepository) CreateShipment(context.Context, string, string, Shipment) error { return nil }
func (*fakeRepository) TransitionShipment(context.Context, string, string, string, string, string, string, string, string) error {
	return nil
}
func (*fakeRepository) CreateCustomerTransport(_ context.Context, _ string, ids TransportIDs, value CustomerTransportCommand) (CustomerTransport, error) {
	return CustomerTransport{ID: ids.ShipmentID, RequestID: value.RequestID, CustomerShipmentID: value.CustomerShipmentID, OriginOrganizationID: value.OriginOrganizationID, ProviderCode: value.ProviderCode, ProviderServiceCode: value.ProviderServiceCode, ProviderReference: value.ProviderReference, State: "booked", Version: 1}, nil
}
func (*fakeRepository) RecordProviderReport(_ context.Context, _ string, shipment, _ string, value ProviderReportCommand) (CustomerTransport, error) {
	state, err := NormalizeProviderReport(value)
	return CustomerTransport{ID: shipment, State: state, Version: value.Version + 1}, err
}
func (*fakeRepository) OpenServiceCase(context.Context, string, string, ServiceCase) error {
	return nil
}
func (*fakeRepository) TransitionServiceCase(context.Context, string, string, string, string, string, int64, string) error {
	return nil
}
func (*fakeRepository) CreateRecall(context.Context, string, string, Recall) error   { return nil }
func (*fakeRepository) ActivateRecall(context.Context, string, string, string) error { return nil }
func (*fakeRepository) AddRecallUnit(context.Context, string, string, string, string, string) error {
	return nil
}
func (*fakeRepository) CreateOrganization(context.Context, string, string, Organization) error {
	return nil
}
func (*fakeRepository) TransitionOrganization(context.Context, string, string, string, string, int64, string) error {
	return nil
}
func (*fakeRepository) CreateAgreement(context.Context, string, string, Agreement) error { return nil }
func (*fakeRepository) TransitionAgreement(context.Context, string, string, string, string, string, int64, string) error {
	return nil
}
func (*fakeRepository) QueueMessage(context.Context, string, string, Message) error { return nil }
func (*fakeRepository) TransitionMessage(context.Context, string, string, string, string, string, string, string) error {
	return nil
}

type testIDs struct{ n int }

func (i *testIDs) New() string {
	i.n++
	return fmt.Sprintf("018f4d4a-7b36-7a21-8d10-%012x", i.n)
}
func TestFulfillmentPoliciesRejectSkippedStates(t *testing.T) {
	s := NewService(&fakeRepository{}, &testIDs{})
	ctx := context.Background()
	if err := s.TransitionShipment(ctx, "tenant", "a", "b", "shipment", "planned", "delivered", ""); err == nil {
		t.Fatal("shipment skipped states")
	}
	if err := s.TransitionServiceCase(ctx, "tenant", "o", "case", "opened", "closed", 1); err == nil {
		t.Fatal("service skipped states")
	}
	if err := s.TransitionAgreement(ctx, "tenant", "o", "agreement", "draft", "suspended", 1); err == nil {
		t.Fatal("agreement skipped states")
	}
	if err := s.TransitionMessage(ctx, "tenant", "message", "queued", "sent", "", ""); err == nil {
		t.Fatal("sent message without provider reference")
	}
}
func TestOrganizationNetworkContracts(t *testing.T) {
	s := NewService(&fakeRepository{}, &testIDs{})
	created, err := s.CreateOrganization(context.Background(), "tenant", Organization{ParentOrganizationID: "franchisor", Code: "store-cordoba", DisplayName: "Store Córdoba", Type: "store"})
	if err != nil || created.Status != "provisioning" || created.Version != 1 || created.ID == "" {
		t.Fatalf("organization=%+v err=%v", created, err)
	}
	if _, err = s.CreateOrganization(context.Background(), "tenant", Organization{Code: "orphan-store", DisplayName: "Orphan", Type: "store"}); err == nil {
		t.Fatal("non-root organization without parent accepted")
	}
	if err = s.TransitionOrganization(context.Background(), "tenant", "store", "provisioning", "active", 1); err != nil {
		t.Fatal(err)
	}
	if err = s.TransitionOrganization(context.Background(), "tenant", "store", "provisioning", "closed", 0); err == nil {
		t.Fatal("organization transition without version accepted")
	}
}

func TestConnectedTransportRejectsInventedProviderCompletion(t *testing.T) {
	s := NewService(&fakeRepository{}, &testIDs{})
	created, err := s.CreateCustomerTransport(context.Background(), "tenant", CustomerTransportCommand{RequestID: "request-1", CustomerShipmentID: "sales-shipment-1", OriginOrganizationID: "store", ProviderCode: "amazon-easyship", ProviderServiceCode: "standard", ProviderReference: "tracking-1"})
	if err != nil || created.ID == "" || created.State != "booked" {
		t.Fatalf("transport=%+v err=%v", created, err)
	}
	command := ProviderReportCommand{OriginOrganizationID: "store", ProviderEventID: "provider-event-1", ProviderStatus: "Delivered", ReportSchema: AmazonEasyShipReconciliationSchema, EvidenceSHA256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", OccurredAt: time.Date(2026, 9, 2, 12, 0, 0, 0, time.UTC), Version: 1}
	delivered, err := s.RecordProviderReport(context.Background(), "tenant", created.ID, command)
	if err != nil || delivered.State != "delivered" || delivered.Version != 2 {
		t.Fatalf("delivered=%+v err=%v", delivered, err)
	}
	command.AutomaticBusinessDeliveryAcceptance = true
	if _, err = s.RecordProviderReport(context.Background(), "tenant", created.ID, command); err == nil {
		t.Fatal("provider report accepted automatic business delivery")
	}
	command.AutomaticBusinessDeliveryAcceptance = false
	command.ProviderStatus = "InventedStatus"
	if _, err = s.RecordProviderReport(context.Background(), "tenant", created.ID, command); err == nil {
		t.Fatal("unknown provider status accepted")
	}
}
