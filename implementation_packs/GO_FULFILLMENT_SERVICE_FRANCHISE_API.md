# Go Fulfillment, Service and Franchise API

## 1. Metadata

```yaml
pack_id: "GO-FULFILLMENT-SERVICE-FRANCHISE-API"
pack_version: "0.6.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Añade logística, servicio, recalls, red de franquicias y transporte conectado a un customer shipment: identidad carrier corregida, reportes provider inmutables/idempotentes, transición optimista y entrega del pedido sin aceptar automáticamente el handover comercial."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "pgx 5.10.0"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "GO-ELECTROMOBILITY-PUBLIC-CRM-API 0.2.x", "GO-SUPPLY-FACTORY-INVENTORY-API 0.16.x", "GO-COMMERCE-PRICING-PAYMENT-API 0.3.x", "PYTHON-AMAZON-SPAPI-SHIPPING-TRACKING-ADAPTER 0.6.x", "PYTHON-AMAZON-SPAPI-EASYSHIP-HANDOVER-ADAPTER 0.1.x", "PYTHON-AMAZON-SPAPI-FULFILLMENT-DELIVERY-EVIDENCE-ADAPTER 0.1.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner AND MIT dependency"
upstream_sources: ["https://go.dev", "https://www.postgresql.org/docs/18/", "https://github.com/jackc/pgx", "https://github.com/microsoft/BCApps/tree/2eae56d704a1fd035d104f333602aea7091b7749", "https://github.com/amzn/selling-partner-api-sdk/tree/8e792ae345a8d334ccdbdd03181f05f040e6a4fc", "https://github.com/amzn/selling-partner-api-models/tree/8e429486005c4ebdce5099e48cc48515a65359bb", "https://developer-docs.amazon.com/sp-api/reference/gettracking", "https://developer-docs.amazon.com/sp-api/reference/getscheduledpackage", "https://developer-docs.amazon.com/sp-api/docs/fulfillment-outbound-api-v2020-07-01-reference"]
verified_at: "2026-09-02"
```

## 2. Applicability

Use when the enterprise composition owns logistics, service/warranty cases, recalls, the organization/franchise network, franchise agreements, consented communications and a carrier boundary connected to posted customer shipments. Reject it when any domain is externally authoritative without an explicit adapter/synchronization contract. Carrier account/access, webhook authentication, polling, address/privacy policy, delivery evidence, customer acceptance and jurisdictional rules remain project conditions.

## 3. Architecture contract

Each command is tenant/organization scoped, versioned where concurrent mutation is possible and transactionally emits an outbox event. Network nodes form a guarded hierarchy and agreement activation is serialized per tenant/territory. A connected carrier shipment references one posted `sales.customer_shipment`, its order, source organization and customer. Request replay is exact; provider reports are immutable, version-bound and normalized only through an admitted provider schema. Multiple not-yet-assigned provider references remain legal while non-null references are unique. Remote carrier/message effects run through idempotent adapters outside the transaction. Provider pickup/delivery never creates or accepts `sales.delivery_handover`; customer checklist/acceptance is a distinct boundary. Invalid transitions, stale versions, unknown provider status/schema, malformed evidence, absent consent and cross-scope commands fail closed.

## 4. Exact file manifest

```text
CREATE internal/fulfillment/service.go
CREATE internal/fulfillment/service_test.go
CREATE internal/platform/postgres/fulfillment.go
CREATE internal/platform/postgres/fulfillment_integration_test.go
CREATE internal/platform/postgres/fulfillment_transport.go
CREATE internal/platform/postgres/fulfillment_transport_integration_test.go
CREATE internal/platform/httpapi/fulfillment.go
CREATE internal/platform/httpapi/fulfillment_test.go
CREATE db/migrations/0008_franchise_network_admin.up.sql
CREATE db/migrations/0008_franchise_network_admin.down.sql
CREATE db/tests/0008_franchise_network_admin.test.sql
CREATE db/migrations/0039_connected_carrier_delivery.up.sql
CREATE db/migrations/0039_connected_carrier_delivery.down.sql
CREATE db/tests/0039_connected_carrier_delivery.test.sql
CREATE docs/logistics/MICROSOFT_BC_AMAZON_CONNECTED_CARRIER_DERIVATION.md
```

## 5. Materialization blocks

### FILE: `internal/fulfillment/service.go`

```yaml
block_id: "GO-FULFILLMENT-API:internal-fulfillment-service-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "8294da7acc1b516ba3857a27ec4ba0e39115b4bd8e0e8dc60c8cea2d21c64bcf"
variables: []
secrets_allowed: false
```

````go
package fulfillment

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var ErrConflict = errors.New("fulfillment conflict")
var codePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

type Shipment struct {
	ID                        string `json:"id"`
	ProviderCode              string `json:"provider_code"`
	ProviderReference         string `json:"provider_reference,omitempty"`
	OriginOrganizationID      string `json:"origin_organization_id"`
	DestinationOrganizationID string `json:"destination_organization_id"`
	State                     string `json:"state"`
}

type CustomerTransportCommand struct {
	RequestID            string `json:"request_id"`
	CustomerShipmentID   string `json:"customer_shipment_id"`
	OriginOrganizationID string `json:"origin_organization_id"`
	ProviderCode         string `json:"provider_code"`
	ProviderServiceCode  string `json:"provider_service_code,omitempty"`
	ProviderReference    string `json:"provider_reference,omitempty"`
}

type CustomerTransport struct {
	ID                         string `json:"id"`
	RequestID                  string `json:"request_id"`
	CustomerShipmentID         string `json:"customer_shipment_id"`
	OrderID                    string `json:"order_id"`
	OriginOrganizationID       string `json:"origin_organization_id"`
	DestinationCustomerSubject string `json:"destination_customer_subject"`
	ProviderCode               string `json:"provider_code"`
	ProviderServiceCode        string `json:"provider_service_code,omitempty"`
	ProviderReference          string `json:"provider_reference,omitempty"`
	State                      string `json:"state"`
	Version                    int64  `json:"version"`
	OrderFulfillmentState      string `json:"order_fulfillment_state"`
}

type ProviderReportCommand struct {
	OriginOrganizationID                string    `json:"origin_organization_id"`
	ProviderEventID                     string    `json:"provider_event_id"`
	ProviderStatus                      string    `json:"provider_status"`
	ReportSchema                        string    `json:"report_schema"`
	EvidenceSHA256                      string    `json:"evidence_sha256"`
	OccurredAt                          time.Time `json:"occurred_at"`
	Version                             int64     `json:"version"`
	AutomaticBusinessDeliveryAcceptance bool      `json:"automatic_business_delivery_acceptance"`
}

type TransportIDs struct {
	ShipmentID string
	EventID    string
}
type ServiceCase struct {
	ID             string `json:"id"`
	StockUnitID    string `json:"stock_unit_id"`
	OrganizationID string `json:"organization_id"`
	State          string `json:"state"`
	Severity       string `json:"severity"`
	Description    string `json:"description"`
	Version        int64  `json:"version"`
}
type Recall struct {
	ID          string     `json:"id"`
	Code        string     `json:"code"`
	Title       string     `json:"title"`
	Severity    string     `json:"severity"`
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}
type Organization struct {
	ID                   string `json:"id"`
	ParentOrganizationID string `json:"parent_organization_id,omitempty"`
	Code                 string `json:"code"`
	DisplayName          string `json:"display_name"`
	Type                 string `json:"type"`
	Status               string `json:"status"`
	Version              int64  `json:"version"`
}
type Agreement struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	TerritoryCode  string     `json:"territory_code"`
	TermsVersion   string     `json:"terms_version"`
	StartsOn       time.Time  `json:"starts_on"`
	EndsOn         *time.Time `json:"ends_on,omitempty"`
	Status         string     `json:"status"`
	Version        int64      `json:"version"`
}
type Message struct {
	ID                   string `json:"id"`
	RecipientPrincipalID string `json:"recipient_principal_id"`
	Channel              string `json:"channel"`
	TemplateCode         string `json:"template_code"`
	TemplateVersion      string `json:"template_version"`
	State                string `json:"state"`
	ProviderReference    string `json:"provider_reference,omitempty"`
	LastErrorCode        string `json:"last_error_code,omitempty"`
}
type Repository interface {
	CreateShipment(context.Context, string, string, Shipment) error
	TransitionShipment(context.Context, string, string, string, string, string, string, string, string) error
	CreateCustomerTransport(context.Context, string, TransportIDs, CustomerTransportCommand) (CustomerTransport, error)
	RecordProviderReport(context.Context, string, string, string, ProviderReportCommand) (CustomerTransport, error)
	OpenServiceCase(context.Context, string, string, ServiceCase) error
	TransitionServiceCase(context.Context, string, string, string, string, string, int64, string) error
	CreateRecall(context.Context, string, string, Recall) error
	ActivateRecall(context.Context, string, string, string) error
	AddRecallUnit(context.Context, string, string, string, string, string) error
	CreateOrganization(context.Context, string, string, Organization) error
	TransitionOrganization(context.Context, string, string, string, string, int64, string) error
	CreateAgreement(context.Context, string, string, Agreement) error
	TransitionAgreement(context.Context, string, string, string, string, string, int64, string) error
	QueueMessage(context.Context, string, string, Message) error
	TransitionMessage(context.Context, string, string, string, string, string, string, string) error
}
type IDGenerator interface{ New() string }
type Service struct {
	repository Repository
	ids        IDGenerator
}

func NewService(r Repository, ids IDGenerator) *Service { return &Service{repository: r, ids: ids} }
func (s *Service) CreateShipment(ctx context.Context, tenant string, v Shipment) (Shipment, error) {
	if tenant == "" || !codePattern.MatchString(v.ProviderCode) || v.OriginOrganizationID == "" || v.DestinationOrganizationID == "" || v.OriginOrganizationID == v.DestinationOrganizationID {
		return Shipment{}, fmt.Errorf("invalid shipment")
	}
	v.ID = s.ids.New()
	v.State = "planned"
	if err := s.repository.CreateShipment(ctx, tenant, s.ids.New(), v); err != nil {
		return Shipment{}, err
	}
	return v, nil
}

var shipmentTransitions = map[string]map[string]bool{"planned": {"booked": true, "cancelled": true}, "booked": {"picked-up": true, "exception": true, "cancelled": true}, "picked-up": {"in-transit": true, "exception": true}, "in-transit": {"delivered": true, "exception": true}, "exception": {"booked": true, "in-transit": true, "cancelled": true}}

func (s *Service) TransitionShipment(ctx context.Context, tenant, originOrganization, destinationOrganization, id, current, target, reference string) error {
	if tenant == "" || originOrganization == "" || destinationOrganization == "" || originOrganization == destinationOrganization || id == "" || !shipmentTransitions[current][target] {
		return fmt.Errorf("%w: invalid shipment transition", ErrConflict)
	}
	return s.repository.TransitionShipment(ctx, tenant, originOrganization, destinationOrganization, id, current, target, reference, s.ids.New())
}

func (s *Service) CreateCustomerTransport(ctx context.Context, tenant string, command CustomerTransportCommand) (CustomerTransport, error) {
	if tenant == "" || !codePattern.MatchString(command.RequestID) || !codePattern.MatchString(command.CustomerShipmentID) || !codePattern.MatchString(command.OriginOrganizationID) || !codePattern.MatchString(command.ProviderCode) || (command.ProviderServiceCode != "" && !codePattern.MatchString(command.ProviderServiceCode)) || (command.ProviderReference != "" && !codePattern.MatchString(command.ProviderReference)) {
		return CustomerTransport{}, fmt.Errorf("%w: invalid customer transport", ErrConflict)
	}
	return s.repository.CreateCustomerTransport(ctx, tenant, TransportIDs{ShipmentID: s.ids.New(), EventID: s.ids.New()}, command)
}

const AmazonEasyShipReconciliationSchema = "elite-amazon-spapi-easyship-reconciliation-receipt/v1"

var amazonEasyShipStates = map[string]string{
	"PickedUp":         "picked-up",
	"AtOriginFC":       "in-transit",
	"AtDestinationFC":  "in-transit",
	"OutForDelivery":   "in-transit",
	"Delivered":        "delivered",
	"Rejected":         "exception",
	"Undeliverable":    "exception",
	"ReturnedToSeller": "exception",
	"LostInTransit":    "exception",
	"DamagedInTransit": "exception",
}

func NormalizeProviderReport(command ProviderReportCommand) (string, error) {
	if command.AutomaticBusinessDeliveryAcceptance || command.ReportSchema != AmazonEasyShipReconciliationSchema || command.OriginOrganizationID == "" || !codePattern.MatchString(command.ProviderEventID) || command.Version < 1 || command.OccurredAt.IsZero() || len(command.ProviderStatus) > 64 {
		return "", fmt.Errorf("%w: invalid provider report", ErrConflict)
	}
	decoded, err := hex.DecodeString(command.EvidenceSHA256)
	if err != nil || len(decoded) != 32 || command.EvidenceSHA256 != strings.ToLower(command.EvidenceSHA256) {
		return "", fmt.Errorf("%w: invalid provider evidence", ErrConflict)
	}
	target, ok := amazonEasyShipStates[command.ProviderStatus]
	if !ok {
		return "", fmt.Errorf("%w: unsupported provider status", ErrConflict)
	}
	return target, nil
}

func (s *Service) RecordProviderReport(ctx context.Context, tenant, shipmentID string, command ProviderReportCommand) (CustomerTransport, error) {
	if tenant == "" || !codePattern.MatchString(shipmentID) {
		return CustomerTransport{}, fmt.Errorf("%w: invalid shipment", ErrConflict)
	}
	if _, err := NormalizeProviderReport(command); err != nil {
		return CustomerTransport{}, err
	}
	return s.repository.RecordProviderReport(ctx, tenant, shipmentID, s.ids.New(), command)
}
func (s *Service) OpenServiceCase(ctx context.Context, tenant string, v ServiceCase) (ServiceCase, error) {
	if tenant == "" || v.StockUnitID == "" || v.OrganizationID == "" || !map[string]bool{"low": true, "medium": true, "high": true, "safety": true}[v.Severity] || len(v.Description) < 3 {
		return ServiceCase{}, fmt.Errorf("invalid service case")
	}
	v.ID = s.ids.New()
	v.State = "opened"
	v.Version = 1
	if err := s.repository.OpenServiceCase(ctx, tenant, s.ids.New(), v); err != nil {
		return ServiceCase{}, err
	}
	return v, nil
}

var serviceTransitions = map[string]map[string]bool{"opened": {"diagnosis": true, "cancelled": true}, "diagnosis": {"awaiting-parts": true, "repair": true, "cancelled": true}, "awaiting-parts": {"repair": true, "cancelled": true}, "repair": {"quality": true, "awaiting-parts": true}, "quality": {"repair": true, "closed": true}}

func (s *Service) TransitionServiceCase(ctx context.Context, tenant, organization, id, current, target string, version int64) error {
	if tenant == "" || organization == "" || id == "" || version < 1 || !serviceTransitions[current][target] {
		return fmt.Errorf("%w: invalid service transition", ErrConflict)
	}
	return s.repository.TransitionServiceCase(ctx, tenant, organization, id, current, target, version, s.ids.New())
}
func (s *Service) CreateRecall(ctx context.Context, tenant string, v Recall) (Recall, error) {
	if tenant == "" || !codePattern.MatchString(v.Code) || len(v.Title) < 3 || !map[string]bool{"service": true, "safety": true, "regulatory": true}[v.Severity] {
		return Recall{}, fmt.Errorf("invalid recall")
	}
	v.ID = s.ids.New()
	v.Status = "draft"
	v.PublishedAt = nil
	if err := s.repository.CreateRecall(ctx, tenant, s.ids.New(), v); err != nil {
		return Recall{}, err
	}
	return v, nil
}
func (s *Service) ActivateRecall(ctx context.Context, tenant, id string) error {
	if tenant == "" || id == "" {
		return fmt.Errorf("invalid recall")
	}
	return s.repository.ActivateRecall(ctx, tenant, id, s.ids.New())
}
func (s *Service) AddRecallUnit(ctx context.Context, tenant, organization, recallID, stockID string) error {
	if tenant == "" || organization == "" || recallID == "" || stockID == "" {
		return fmt.Errorf("invalid recall unit")
	}
	return s.repository.AddRecallUnit(ctx, tenant, organization, recallID, stockID, s.ids.New())
}
func (s *Service) CreateOrganization(ctx context.Context, tenant string, v Organization) (Organization, error) {
	types := map[string]bool{"enterprise": true, "franchisor": true, "franchisee": true, "factory": true, "warehouse": true, "store": true, "service_center": true}
	root := v.Type == "enterprise" || v.Type == "franchisor"
	if tenant == "" || !codePattern.MatchString(v.Code) || len(v.DisplayName) < 2 || len(v.DisplayName) > 100 || !types[v.Type] || (root && v.ParentOrganizationID != "") || (!root && v.ParentOrganizationID == "") {
		return Organization{}, fmt.Errorf("invalid organization")
	}
	v.ID = s.ids.New()
	v.Status = "provisioning"
	v.Version = 1
	if err := s.repository.CreateOrganization(ctx, tenant, s.ids.New(), v); err != nil {
		return Organization{}, err
	}
	return v, nil
}

var organizationTransitions = map[string]map[string]bool{"provisioning": {"active": true, "closed": true}, "active": {"suspended": true, "closed": true}, "suspended": {"active": true, "closed": true}}

func (s *Service) TransitionOrganization(ctx context.Context, tenant, id, current, target string, version int64) error {
	if tenant == "" || id == "" || version < 1 || !organizationTransitions[current][target] {
		return fmt.Errorf("%w: invalid organization transition", ErrConflict)
	}
	return s.repository.TransitionOrganization(ctx, tenant, id, current, target, version, s.ids.New())
}
func (s *Service) CreateAgreement(ctx context.Context, tenant string, v Agreement) (Agreement, error) {
	if tenant == "" || v.OrganizationID == "" || !codePattern.MatchString(v.TerritoryCode) || !codePattern.MatchString(v.TermsVersion) || v.StartsOn.IsZero() || (v.EndsOn != nil && !v.EndsOn.After(v.StartsOn)) {
		return Agreement{}, fmt.Errorf("invalid agreement")
	}
	v.ID = s.ids.New()
	v.Status = "draft"
	v.Version = 1
	if err := s.repository.CreateAgreement(ctx, tenant, s.ids.New(), v); err != nil {
		return Agreement{}, err
	}
	return v, nil
}

var agreementTransitions = map[string]map[string]bool{"draft": {"active": true, "terminated": true}, "active": {"suspended": true, "terminated": true, "expired": true}, "suspended": {"active": true, "terminated": true, "expired": true}}

func (s *Service) TransitionAgreement(ctx context.Context, tenant, organization, id, current, target string, version int64) error {
	if tenant == "" || organization == "" || id == "" || version < 1 || !agreementTransitions[current][target] {
		return fmt.Errorf("%w: invalid agreement transition", ErrConflict)
	}
	return s.repository.TransitionAgreement(ctx, tenant, organization, id, current, target, version, s.ids.New())
}
func (s *Service) QueueMessage(ctx context.Context, tenant string, v Message) (Message, error) {
	if tenant == "" || v.RecipientPrincipalID == "" || !map[string]bool{"email": true, "sms": true, "push": true, "chat": true}[v.Channel] || !codePattern.MatchString(v.TemplateCode) || !codePattern.MatchString(v.TemplateVersion) {
		return Message{}, fmt.Errorf("invalid message")
	}
	v.ID = s.ids.New()
	v.State = "queued"
	if err := s.repository.QueueMessage(ctx, tenant, s.ids.New(), v); err != nil {
		return Message{}, err
	}
	return v, nil
}

var messageTransitions = map[string]map[string]bool{"queued": {"sent": true, "failed": true, "suppressed": true}, "sent": {"delivered": true, "failed": true}}

func (s *Service) TransitionMessage(ctx context.Context, tenant, id, current, target, reference, errorCode string) error {
	if tenant == "" || id == "" || !messageTransitions[current][target] || (target == "sent" && reference == "") || (target == "failed" && errorCode == "") {
		return fmt.Errorf("%w: invalid message transition", ErrConflict)
	}
	return s.repository.TransitionMessage(ctx, tenant, id, current, target, reference, errorCode, s.ids.New())
}
````

### FILE: `internal/fulfillment/service_test.go`

```yaml
block_id: "GO-FULFILLMENT-API:internal-fulfillment-service_test-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "9d930a39921c216c9a713bfae7e2d86962afc6734e970230687ef1641146bfb4"
variables: []
secrets_allowed: false
```

````go
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
````

### FILE: `internal/platform/postgres/fulfillment.go`

```yaml
block_id: "GO-FULFILLMENT-API:internal-platform-postgres-fulfillment-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "eb457fc6310b1d68d91815480455b69790feacf40a348a36d876015b6a467f51"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"elite.local/enterprise/internal/fulfillment"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Fulfillment struct{ pool *pgxpool.Pool }

func NewFulfillment(pool *pgxpool.Pool) *Fulfillment { return &Fulfillment{pool: pool} }
func (r *Fulfillment) CreateShipment(ctx context.Context, tenant, eventID string, v fulfillment.Shipment) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into logistics.shipment(tenant_id,shipment_id,provider_code,provider_reference,origin_organization_id,destination_organization_id,state)values($1,$2,$3,nullif($4,''),$5,$6,'planned')`, tenant, v.ID, v.ProviderCode, v.ProviderReference, v.OriginOrganizationID, v.DestinationOrganizationID)
	if err != nil {
		return err
	}
	if err = outbox(ctx, tx, tenant, eventID, "shipment", v.ID, 1, "shipment.planned", `jsonb_build_object('provider_code',$7::text)`, v.ProviderCode); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Fulfillment) TransitionShipment(ctx context.Context, tenant, originOrganization, destinationOrganization, id, current, target, reference, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update logistics.shipment set state=$6,provider_reference=coalesce(nullif($7,''),provider_reference),last_event_at=clock_timestamp() where tenant_id=$1 and origin_organization_id=$2 and destination_organization_id=$3 and shipment_id=$4 and state=$5`, tenant, originOrganization, destinationOrganization, id, current, target, reference)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "shipment", id, 1, "shipment."+target, `'{}'::jsonb`); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Fulfillment) OpenServiceCase(ctx context.Context, tenant, eventID string, v fulfillment.ServiceCase) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = openServiceCaseInTx(ctx, tx, tenant, eventID, v); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED transaction composition only: original case writer SQL preserved.
func openServiceCaseInTx(ctx context.Context, tx pgx.Tx, tenant, eventID string, v fulfillment.ServiceCase) error {
	var err error
	result, err := tx.Exec(ctx, `insert into service_ops.service_case(tenant_id,service_case_id,stock_unit_id,organization_id,state,severity,description,version)select $1,$2,s.stock_unit_id,$4,'opened',$5,$6,1 from inventory.stock_unit s where s.tenant_id=$1 and s.stock_unit_id=$3 and s.organization_id=$4`, tenant, v.ID, v.StockUnitID, v.OrganizationID, v.Severity, v.Description)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "service-case", v.ID, 1, "service-case.opened", `jsonb_build_object('severity',$7::text)`, v.Severity); err != nil {
		return err
	}
	return nil
}
func (r *Fulfillment) TransitionServiceCase(ctx context.Context, tenant, organization, id, current, target string, version int64, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update service_ops.service_case set state=$6,version=version+1,closed_at=case when $6='closed' then clock_timestamp() else null end where tenant_id=$1 and organization_id=$2 and service_case_id=$3 and state=$4 and version=$5`, tenant, organization, id, current, version, target)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "service-case", id, version+1, "service-case."+target, `'{}'::jsonb`); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Fulfillment) CreateRecall(ctx context.Context, tenant, eventID string, v fulfillment.Recall) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into service_ops.recall(tenant_id,recall_id,recall_code,title,severity,status)values($1,$2,$3,$4,$5,'draft')`, tenant, v.ID, v.Code, v.Title, v.Severity)
	if err != nil {
		return err
	}
	if err = outbox(ctx, tx, tenant, eventID, "recall", v.ID, 1, "recall.created", `jsonb_build_object('recall_code',$7::text)`, v.Code); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Fulfillment) ActivateRecall(ctx context.Context, tenant, id, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update service_ops.recall set status='active',published_at=clock_timestamp() where tenant_id=$1 and recall_id=$2 and status='draft' and published_at is null`, tenant, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "recall", id, 2, "recall.activated", `'{}'::jsonb`); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Fulfillment) AddRecallUnit(ctx context.Context, tenant, organization, recallID, stockID, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into service_ops.recall_unit(tenant_id,recall_id,stock_unit_id,state)select r.tenant_id,r.recall_id,s.stock_unit_id,'identified' from service_ops.recall r join inventory.stock_unit s on s.tenant_id=r.tenant_id where r.tenant_id=$1 and s.organization_id=$2 and r.recall_id=$3 and r.status='active' and s.stock_unit_id=$4 on conflict do nothing`, tenant, organization, recallID, stockID)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "recall", recallID, 3, "recall.unit-identified", `jsonb_build_object('stock_unit_id',$7::text)`, stockID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Fulfillment) CreateOrganization(ctx context.Context, tenant, eventID string, v fulfillment.Organization) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = createOrganizationInTx(ctx, tx, tenant, eventID, v); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED transaction extraction; original owner SQL and guards preserved.
func createOrganizationInTx(ctx context.Context, tx pgx.Tx, tenant, eventID string, v fulfillment.Organization) error {
	result, err := tx.Exec(ctx, `insert into org.organization(tenant_id,organization_id,parent_organization_id,organization_code,display_name,organization_type,status,version)
		select t.tenant_id,$2,nullif($3,''),$4,$5,$6,'provisioning',1
		from platform.tenant t left join org.organization p on p.tenant_id=t.tenant_id and p.organization_id=nullif($3,'')
		where t.tenant_id=$1 and t.status='active' and (
		  ($3='' and $6 in ('enterprise','franchisor')) or
		  ($3<>'' and p.status='active' and (($6='franchisee' and p.organization_type in ('enterprise','franchisor')) or ($6 in ('store','service_center','warehouse') and p.organization_type in ('enterprise','franchisor','franchisee')) or ($6='factory' and p.organization_type in ('enterprise','franchisor'))))
		)`, tenant, v.ID, v.ParentOrganizationID, v.Code, v.DisplayName, v.Type)
	if err != nil {
		return fulfillmentConstraint(err)
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "organization", v.ID, 1, "organization.provisioning", `jsonb_build_object('parent_organization_id',nullif($7::text,''),'organization_type',$8::text)`, v.ParentOrganizationID, v.Type); err != nil {
		return err
	}
	return nil
}
func (r *Fulfillment) TransitionOrganization(ctx context.Context, tenant, id, current, target string, version int64, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = transitionOrganizationInTx(ctx, tx, tenant, id, current, target, version, eventID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED transaction extraction; original owner SQL and guards preserved.
func transitionOrganizationInTx(ctx context.Context, tx pgx.Tx, tenant, id, current, target string, version int64, eventID string) error {
	result, err := tx.Exec(ctx, `update org.organization o set status=$4,version=version+1,updated_at=clock_timestamp()
		where tenant_id=$1 and organization_id=$2 and status=$3 and version=$5
		and not ($4='closed' and (exists(select 1 from org.organization child where child.tenant_id=o.tenant_id and child.parent_organization_id=o.organization_id and child.status<>'closed') or exists(select 1 from franchise.agreement a where a.tenant_id=o.tenant_id and a.franchise_organization_id=o.organization_id and a.status in ('active','suspended'))))`, tenant, id, current, target, version)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "organization", id, version+1, "organization."+target, `'{}'::jsonb`); err != nil {
		return err
	}
	return nil
}
func (r *Fulfillment) CreateAgreement(ctx context.Context, tenant, eventID string, v fulfillment.Agreement) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = createAgreementInTx(ctx, tx, tenant, eventID, v); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED transaction extraction; original owner SQL and guards preserved.
func createAgreementInTx(ctx context.Context, tx pgx.Tx, tenant, eventID string, v fulfillment.Agreement) error {
	result, err := tx.Exec(ctx, `insert into franchise.agreement(tenant_id,agreement_id,franchise_organization_id,territory_code,terms_version,starts_on,ends_on,status,version)
		select $1,$2,o.organization_id,$4,$5,$6,$7,'draft',1 from org.organization o where o.tenant_id=$1 and o.organization_id=$3 and o.organization_type='franchisee' and o.status='active'`, tenant, v.ID, v.OrganizationID, v.TerritoryCode, v.TermsVersion, v.StartsOn, v.EndsOn)
	if err != nil {
		return fulfillmentConstraint(err)
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "franchise-agreement", v.ID, 1, "franchise-agreement.created", `jsonb_build_object('organization_id',$7::text)`, v.OrganizationID); err != nil {
		return err
	}
	return nil
}
func (r *Fulfillment) TransitionAgreement(ctx context.Context, tenant, organization, id, current, target string, version int64, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = transitionAgreementInTx(ctx, tx, tenant, organization, id, current, target, version, eventID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED transaction extraction; original owner SQL and guards preserved.
func transitionAgreementInTx(ctx context.Context, tx pgx.Tx, tenant, organization, id, current, target string, version int64, eventID string) error {
	result, err := tx.Exec(ctx, `update franchise.agreement a set status=$5,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and franchise_organization_id=$2 and agreement_id=$3 and status=$4 and version=$6 and ($5<>'active' or exists(select 1 from org.organization o where o.tenant_id=a.tenant_id and o.organization_id=a.franchise_organization_id and o.organization_type='franchisee' and o.status='active'))`, tenant, organization, id, current, target, version)
	if err != nil {
		return fulfillmentConstraint(err)
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "franchise-agreement", id, version+1, "franchise-agreement."+target, `'{}'::jsonb`); err != nil {
		return err
	}
	return nil
}
func (r *Fulfillment) QueueMessage(ctx context.Context, tenant, eventID string, v fulfillment.Message) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into communication.message(tenant_id,message_id,recipient_principal_id,channel,template_code,template_version,state)values($1,$2,$3,$4,$5,$6,'queued')`, tenant, v.ID, v.RecipientPrincipalID, v.Channel, v.TemplateCode, v.TemplateVersion)
	if err != nil {
		return err
	}
	if err = outbox(ctx, tx, tenant, eventID, "message", v.ID, 1, "message.queued", `jsonb_build_object('channel',$7::text)`, v.Channel); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Fulfillment) TransitionMessage(ctx context.Context, tenant, id, current, target, reference, errorCode, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update communication.message set state=$4,provider_reference=coalesce(nullif($5,''),provider_reference),last_error_code=nullif($6,''),updated_at=clock_timestamp() where tenant_id=$1 and message_id=$2 and state=$3`, tenant, id, current, target, reference, errorCode)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "message", id, 1, "message."+target, `'{}'::jsonb`); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func outbox(ctx context.Context, tx interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}, tenant, eventID, aggregateType, aggregateID string, version int64, eventType, payloadExpression string, payloadValues ...any) error {
	query := `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,$3,$4,$5,$6,1,clock_timestamp(),` + payloadExpression + `)`
	args := []any{tenant, eventID, aggregateType, aggregateID, version, eventType}
	args = append(args, payloadValues...)
	_, err := tx.Exec(ctx, query, args...)
	return err
}

func fulfillmentConstraint(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23503" || pgErr.Code == "23505" || pgErr.Code == "23514" || pgErr.Code == "23P01") {
		return fulfillment.ErrConflict
	}
	return err
}
````

### FILE: `internal/platform/postgres/fulfillment_integration_test.go`

```yaml
block_id: "GO-FULFILLMENT-API:internal-platform-postgres-fulfillment_integration_test-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "7521135d72d2261fe4f18d3c10ffa7555f474fecf096466c17f4f5b1c359b873"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"elite.local/enterprise/internal/fulfillment"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"sync"
	"testing"
	"time"
)

func TestFulfillmentServiceFranchiseFlow(t *testing.T) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28701"
	cleanup := func() {
		for _, q := range []string{`delete from platform.outbox_event where tenant_id=$1`, `delete from communication.message where tenant_id=$1`, `delete from franchise.agreement where tenant_id=$1`, `delete from service_ops.recall_unit where tenant_id=$1`, `delete from service_ops.recall where tenant_id=$1`, `delete from service_ops.service_case where tenant_id=$1`, `delete from logistics.shipment where tenant_id=$1`, `delete from inventory.stock_unit where tenant_id=$1`, `delete from catalog.vehicle_variant where tenant_id=$1`, `delete from catalog.vehicle_model where tenant_id=$1`, `delete from crm.customer_profile where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`} {
			_, _ = pool.Exec(ctx, q, tenant)
		}
	}
	cleanup()
	defer cleanup()
	fixtures := []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'fulfillment-api','Fulfillment','Fulfillment')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'org-a','org-a','A','franchisee'),($1,'org-b','org-b','B','service_center')`, `insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,status)values($1,'customer','Customer','active')`, `insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Model','bicycle','active')`, `insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Variant','{}','active')`, `insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','org-b','variant','FULFILLMENT-SERIAL','available',1,clock_timestamp())`}
	for _, q := range fixtures {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	r := NewFulfillment(pool)
	n := 0
	event := func() string { n++; return "018f4d4a-7b36-7a21-8d10-2f4c54c28" + fmt.Sprintf("%03d", 700+n) }
	if err = r.CreateShipment(ctx, tenant, event(), fulfillment.Shipment{ID: "shipment", ProviderCode: "carrier", OriginOrganizationID: "org-a", DestinationOrganizationID: "org-b"}); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionShipment(ctx, tenant, "org-other", "org-b", "shipment", "planned", "booked", "carrier-ref", event()); err == nil {
		t.Fatal("shipment crossed organization scope")
	}
	if err = r.TransitionShipment(ctx, tenant, "org-a", "org-b", "shipment", "planned", "booked", "carrier-ref", event()); err != nil {
		t.Fatal(err)
	}
	c := fulfillment.ServiceCase{ID: "case", StockUnitID: "stock", OrganizationID: "org-b", Severity: "safety", Description: "brake inspection"}
	if err = r.OpenServiceCase(ctx, tenant, event(), c); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionServiceCase(ctx, tenant, "org-a", "case", "opened", "diagnosis", 1, event()); err == nil {
		t.Fatal("service case crossed organization scope")
	}
	if err = r.TransitionServiceCase(ctx, tenant, "org-b", "case", "opened", "diagnosis", 1, event()); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionServiceCase(ctx, tenant, "org-b", "case", "opened", "cancelled", 1, event()); err == nil {
		t.Fatal("stale service transition succeeded")
	}
	if err = r.CreateRecall(ctx, tenant, event(), fulfillment.Recall{ID: "recall", Code: "REC-1", Title: "Safety inspection", Severity: "safety"}); err != nil {
		t.Fatal(err)
	}
	if err = r.ActivateRecall(ctx, tenant, "recall", event()); err != nil {
		t.Fatal(err)
	}
	if err = r.AddRecallUnit(ctx, tenant, "org-a", "recall", "stock", event()); err == nil {
		t.Fatal("recall unit crossed organization scope")
	}
	if err = r.AddRecallUnit(ctx, tenant, "org-b", "recall", "stock", event()); err != nil {
		t.Fatal(err)
	}
	if err = r.AddRecallUnit(ctx, tenant, "org-b", "recall", "stock", event()); err == nil {
		t.Fatal("duplicate recall unit succeeded")
	}
	if err = r.CreateOrganization(ctx, tenant, event(), fulfillment.Organization{ID: "store-child", ParentOrganizationID: "org-a", Code: "store-child", DisplayName: "Store Child", Type: "store", Status: "provisioning", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionOrganization(ctx, tenant, "store-child", "provisioning", "active", 1, event()); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionOrganization(ctx, tenant, "store-child", "provisioning", "closed", 1, event()); err == nil {
		t.Fatal("stale organization transition succeeded")
	}
	if err = r.TransitionOrganization(ctx, tenant, "org-a", "active", "closed", 1, event()); err == nil {
		t.Fatal("parent with active child was closed")
	}
	starts := time.Now().UTC()
	if err = r.CreateAgreement(ctx, tenant, event(), fulfillment.Agreement{ID: "agreement", OrganizationID: "org-a", TerritoryCode: "AR-CBA", TermsVersion: "v1", StartsOn: starts}); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionAgreement(ctx, tenant, "org-b", "agreement", "draft", "active", 1, event()); err == nil {
		t.Fatal("agreement crossed organization scope")
	}
	if err = r.TransitionAgreement(ctx, tenant, "org-a", "agreement", "draft", "active", 1, event()); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionAgreement(ctx, tenant, "org-a", "agreement", "draft", "terminated", 1, event()); err == nil {
		t.Fatal("stale agreement transition succeeded")
	}
	if err = r.CreateAgreement(ctx, tenant, event(), fulfillment.Agreement{ID: "overlap-agreement", OrganizationID: "org-a", TerritoryCode: "AR-CBA", TermsVersion: "v2", StartsOn: starts.AddDate(0, 1, 0)}); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionAgreement(ctx, tenant, "org-a", "overlap-agreement", "draft", "active", 1, event()); err == nil {
		t.Fatal("overlapping active territory succeeded")
	}
	for _, agreement := range []fulfillment.Agreement{
		{ID: "concurrent-agreement-a", OrganizationID: "org-a", TerritoryCode: "AR-SFE", TermsVersion: "v1", StartsOn: starts},
		{ID: "concurrent-agreement-b", OrganizationID: "org-a", TerritoryCode: "AR-SFE", TermsVersion: "v2", StartsOn: starts.AddDate(0, 1, 0)},
	} {
		if err = r.CreateAgreement(ctx, tenant, event(), agreement); err != nil {
			t.Fatal(err)
		}
	}
	events := []string{event(), event()}
	start := make(chan struct{})
	results := make(chan error, 2)
	var workers sync.WaitGroup
	for index, agreementID := range []string{"concurrent-agreement-a", "concurrent-agreement-b"} {
		workers.Add(1)
		go func(id, eventID string) {
			defer workers.Done()
			<-start
			results <- r.TransitionAgreement(ctx, tenant, "org-a", id, "draft", "active", 1, eventID)
		}(agreementID, events[index])
	}
	close(start)
	workers.Wait()
	close(results)
	succeeded, conflicted := 0, 0
	for result := range results {
		if result == nil {
			succeeded++
		} else if errors.Is(result, fulfillment.ErrConflict) {
			conflicted++
		} else {
			t.Fatal(result)
		}
	}
	if succeeded != 1 || conflicted != 1 {
		t.Fatalf("concurrent territory activation success=%d conflict=%d", succeeded, conflicted)
	}
	if err = r.QueueMessage(ctx, tenant, event(), fulfillment.Message{ID: "message", RecipientPrincipalID: "customer", Channel: "email", TemplateCode: "recall", TemplateVersion: "v1"}); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionMessage(ctx, tenant, "message", "queued", "sent", "provider-ref", "", event()); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionMessage(ctx, tenant, "message", "sent", "delivered", "provider-ref", "", event()); err != nil {
		t.Fatal(err)
	}
	var outboxCount int
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1`, tenant).Scan(&outboxCount); err != nil || outboxCount != 18 {
		t.Fatalf("outbox=%d err=%v", outboxCount, err)
	}
}
````

### FILE: `internal/platform/httpapi/fulfillment.go`

```yaml
block_id: "GO-FULFILLMENT-API:internal-platform-httpapi-fulfillment-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "a74a2ac43e69665be91e5b9d68a20f370ca953bf0e3cfd6ba4161e134f2a60a5"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"elite.local/enterprise/internal/fulfillment"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"net/http"
)

type FulfillmentModule struct{ Service *fulfillment.Service }

func (m FulfillmentModule) Register(mux *http.ServeMux, v identity.Verifier) {
	a := fulfillmentAPI{service: m.Service, verifier: v}
	mux.HandleFunc("POST /v1/logistics/shipments", a.createShipment)
	mux.HandleFunc("POST /v1/logistics/shipments/{id}/transitions", a.transitionShipment)
	mux.HandleFunc("POST /v1/logistics/customer-transports", a.createCustomerTransport)
	mux.HandleFunc("POST /v1/logistics/shipments/{id}/provider-reports", a.recordProviderReport)
	mux.HandleFunc("POST /v1/service/cases", a.openCase)
	mux.HandleFunc("POST /v1/service/cases/{id}/transitions", a.transitionCase)
	mux.HandleFunc("POST /v1/service/recalls", a.createRecall)
	mux.HandleFunc("POST /v1/service/recalls/{id}/activate", a.activateRecall)
	mux.HandleFunc("POST /v1/service/recalls/{id}/units", a.addRecallUnit)
	mux.HandleFunc("POST /v1/franchise/network/organizations", a.createOrganization)
	mux.HandleFunc("POST /v1/franchise/network/organizations/{id}/transitions", a.transitionOrganization)
	mux.HandleFunc("POST /v1/franchise/agreements", a.createAgreement)
	mux.HandleFunc("POST /v1/franchise/agreements/{id}/transitions", a.transitionAgreement)
	mux.HandleFunc("POST /v1/communications/messages", a.queueMessage)
	mux.HandleFunc("POST /v1/communications/messages/{id}/transitions", a.transitionMessage)
}
func (a fulfillmentAPI) createOrganization(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "network:admin")
	if !ok {
		return
	}
	var v fulfillment.Organization
	if !decodeStrict(w, r, &v) {
		return
	}
	if (v.ParentOrganizationID == "" && !p.Allowed("*")) || (v.ParentOrganizationID != "" && !p.AllowedOrganization(v.ParentOrganizationID)) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for the parent organization")
		return
	}
	result, err := a.service.CreateOrganization(r.Context(), p.TenantID, v)
	if err != nil {
		writeFulfillment(w, err)
		return
	}
	writeJSON(w, 201, result)
}
func (a fulfillmentAPI) transitionOrganization(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "network:admin")
	if !ok {
		return
	}
	var v struct {
		Current string `json:"current"`
		Target  string `json:"target"`
		Version int64  `json:"version"`
	}
	if !decodeStrict(w, r, &v) {
		return
	}
	if !p.AllowedOrganization(r.PathValue("id")) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeFulfillment(w, a.service.TransitionOrganization(r.Context(), p.TenantID, r.PathValue("id"), v.Current, v.Target, v.Version))
}

type fulfillmentAPI struct {
	service  *fulfillment.Service
	verifier identity.Verifier
}

func (a fulfillmentAPI) auth(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
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
func (a fulfillmentAPI) createShipment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "logistics:write")
	if !ok {
		return
	}
	var v fulfillment.Shipment
	if !decodeStrict(w, r, &v) {
		return
	}
	if !p.AllowedOrganization(v.OriginOrganizationID) || !p.AllowedOrganization(v.DestinationOrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for every shipment organization")
		return
	}
	result, err := a.service.CreateShipment(r.Context(), p.TenantID, v)
	if err != nil {
		writeFulfillment(w, err)
		return
	}
	writeJSON(w, 201, result)
}
func (a fulfillmentAPI) transitionShipment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "logistics:write")
	if !ok {
		return
	}
	var v struct {
		OriginOrganizationID      string `json:"origin_organization_id"`
		DestinationOrganizationID string `json:"destination_organization_id"`
		Current                   string `json:"current"`
		Target                    string `json:"target"`
		ProviderReference         string `json:"provider_reference"`
	}
	if !decodeStrict(w, r, &v) {
		return
	}
	if !p.AllowedOrganization(v.OriginOrganizationID) || !p.AllowedOrganization(v.DestinationOrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for every shipment organization")
		return
	}
	writeFulfillment(w, a.service.TransitionShipment(r.Context(), p.TenantID, v.OriginOrganizationID, v.DestinationOrganizationID, r.PathValue("id"), v.Current, v.Target, v.ProviderReference))
}
func (a fulfillmentAPI) createCustomerTransport(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "logistics:write")
	if !ok {
		return
	}
	var value fulfillment.CustomerTransportCommand
	if !decodeStrict(w, r, &value) {
		return
	}
	if !p.AllowedOrganization(value.OriginOrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for the shipment origin")
		return
	}
	result, err := a.service.CreateCustomerTransport(r.Context(), p.TenantID, value)
	if err != nil {
		writeFulfillment(w, err)
		return
	}
	writeJSON(w, 201, result)
}
func (a fulfillmentAPI) recordProviderReport(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "logistics:write")
	if !ok {
		return
	}
	var value fulfillment.ProviderReportCommand
	if !decodeStrict(w, r, &value) {
		return
	}
	if !p.AllowedOrganization(value.OriginOrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for the shipment origin")
		return
	}
	result, err := a.service.RecordProviderReport(r.Context(), p.TenantID, r.PathValue("id"), value)
	if err != nil {
		writeFulfillment(w, err)
		return
	}
	writeJSON(w, 200, result)
}
func (a fulfillmentAPI) openCase(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "service:write")
	if !ok {
		return
	}
	var v fulfillment.ServiceCase
	if !decodeStrict(w, r, &v) {
		return
	}
	if !p.AllowedOrganization(v.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	result, err := a.service.OpenServiceCase(r.Context(), p.TenantID, v)
	if err != nil {
		writeFulfillment(w, err)
		return
	}
	writeJSON(w, 201, result)
}
func (a fulfillmentAPI) transitionCase(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "service:write")
	if !ok {
		return
	}
	var v struct {
		OrganizationID string `json:"organization_id"`
		Current        string `json:"current"`
		Target         string `json:"target"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &v) {
		return
	}
	if !p.AllowedOrganization(v.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeFulfillment(w, a.service.TransitionServiceCase(r.Context(), p.TenantID, v.OrganizationID, r.PathValue("id"), v.Current, v.Target, v.Version))
}
func (a fulfillmentAPI) createRecall(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "recall:write")
	if !ok {
		return
	}
	var v fulfillment.Recall
	if !decodeStrict(w, r, &v) {
		return
	}
	result, err := a.service.CreateRecall(r.Context(), p.TenantID, v)
	if err != nil {
		writeFulfillment(w, err)
		return
	}
	writeJSON(w, 201, result)
}
func (a fulfillmentAPI) activateRecall(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "recall:publish")
	if !ok {
		return
	}
	writeFulfillment(w, a.service.ActivateRecall(r.Context(), p.TenantID, r.PathValue("id")))
}
func (a fulfillmentAPI) addRecallUnit(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "recall:write")
	if !ok {
		return
	}
	var v struct {
		OrganizationID string `json:"organization_id"`
		StockUnitID    string `json:"stock_unit_id"`
	}
	if !decodeStrict(w, r, &v) {
		return
	}
	if !p.AllowedOrganization(v.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeFulfillment(w, a.service.AddRecallUnit(r.Context(), p.TenantID, v.OrganizationID, r.PathValue("id"), v.StockUnitID))
}
func (a fulfillmentAPI) createAgreement(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "franchise:write")
	if !ok {
		return
	}
	var v fulfillment.Agreement
	if !decodeStrict(w, r, &v) {
		return
	}
	if !p.AllowedOrganization(v.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	result, err := a.service.CreateAgreement(r.Context(), p.TenantID, v)
	if err != nil {
		writeFulfillment(w, err)
		return
	}
	writeJSON(w, 201, result)
}
func (a fulfillmentAPI) transitionAgreement(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "franchise:write")
	if !ok {
		return
	}
	var v struct {
		OrganizationID string `json:"organization_id"`
		Current        string `json:"current"`
		Target         string `json:"target"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &v) {
		return
	}
	if !p.AllowedOrganization(v.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeFulfillment(w, a.service.TransitionAgreement(r.Context(), p.TenantID, v.OrganizationID, r.PathValue("id"), v.Current, v.Target, v.Version))
}
func (a fulfillmentAPI) queueMessage(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "communication:write")
	if !ok {
		return
	}
	var v fulfillment.Message
	if !decodeStrict(w, r, &v) {
		return
	}
	result, err := a.service.QueueMessage(r.Context(), p.TenantID, v)
	if err != nil {
		writeFulfillment(w, err)
		return
	}
	writeJSON(w, 202, result)
}
func (a fulfillmentAPI) transitionMessage(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "communication:provider")
	if !ok {
		return
	}
	var v struct {
		Current           string `json:"current"`
		Target            string `json:"target"`
		ProviderReference string `json:"provider_reference"`
		ErrorCode         string `json:"error_code"`
	}
	if !decodeStrict(w, r, &v) {
		return
	}
	writeFulfillment(w, a.service.TransitionMessage(r.Context(), p.TenantID, r.PathValue("id"), v.Current, v.Target, v.ProviderReference, v.ErrorCode))
}
func writeFulfillment(w http.ResponseWriter, err error) {
	if errors.Is(err, fulfillment.ErrConflict) {
		writeProblem(w, 409, "FULFILLMENT_CONFLICT", "state, version or business invariant conflict")
		return
	}
	if err != nil {
		writeProblem(w, 400, "INVALID_FULFILLMENT_COMMAND", "command does not match contract")
		return
	}
	writeJSON(w, 200, map[string]string{"status": "accepted"})
}
````

### FILE: `internal/platform/httpapi/fulfillment_test.go`

```yaml
block_id: "GO-FULFILLMENT-API:internal-platform-httpapi-fulfillment_test-go:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "265a7caa54222242b4d43ebe615158508664f131a58285d29f2270c41304fed3"
variables: []
secrets_allowed: false
```

````go
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
````

### FILE: `db/migrations/0008_franchise_network_admin.up.sql`

```yaml
block_id: "GO-FULFILLMENT-API:db-migrations-0008-franchise-network-admin-up-sql:v3"
operation: CREATE
provenance: AUTHORED
source: "local verified composition governed by pinned Microsoft BCApps network/intercompany sources"
license: "LicenseRef-Workspace-Owner"
sha256: "4b27ef2e97a616458fe2d4f7151f58359db8f4ab12c142880e8b7078c0771922"
variables: []
secrets_allowed: false
```

````sql
begin;

alter table org.organization add column version bigint not null default 1 check (version > 0);
alter table franchise.agreement add column version bigint not null default 1 check (version > 0), add column updated_at timestamptz not null default clock_timestamp();

do $$
begin
  if exists (
    with recursive hierarchy as (
      select o.tenant_id,o.organization_id,o.parent_organization_id,array[o.organization_id]::text[] as path,false as cycle
      from org.organization o
      union all
      select p.tenant_id,p.organization_id,p.parent_organization_id,h.path || p.organization_id,p.organization_id=any(h.path)
      from hierarchy h
      join org.organization p on p.tenant_id=h.tenant_id and p.organization_id=h.parent_organization_id
      where not h.cycle
    )
    select 1 from hierarchy where cycle
  ) then
    raise exception using errcode='23514', message='existing organization hierarchy contains a cycle';
  end if;
end;
$$;

create or replace function org.enforce_organization_hierarchy()
returns trigger language plpgsql as $$
begin
  perform pg_advisory_xact_lock(hashtextextended(new.tenant_id::text || '|organization-hierarchy', 0));
  if new.parent_organization_id is not null then
    if not exists (select 1 from org.organization p where p.tenant_id=new.tenant_id and p.organization_id=new.parent_organization_id and p.status <> 'closed') then
      raise exception using errcode='23503', message='organization parent is unavailable';
    end if;
    if exists (
      with recursive ancestors as (
        select p.organization_id,p.parent_organization_id,array[p.organization_id]::text[] as path,false as cycle from org.organization p where p.tenant_id=new.tenant_id and p.organization_id=new.parent_organization_id
        union all
        select p.organization_id,p.parent_organization_id,a.path || p.organization_id,p.organization_id=any(a.path) from org.organization p join ancestors a on a.parent_organization_id=p.organization_id where p.tenant_id=new.tenant_id and not a.cycle
      ) select 1 from ancestors where organization_id=new.organization_id or cycle
    ) then
      raise exception using errcode='23514', message='organization hierarchy cycle';
    end if;
  end if;
  return new;
end;
$$;

create trigger organization_hierarchy_guard before insert or update of parent_organization_id on org.organization for each row execute function org.enforce_organization_hierarchy();

create or replace function franchise.enforce_active_territory_non_overlap()
returns trigger language plpgsql as $$
begin
  perform pg_advisory_xact_lock(hashtextextended(new.tenant_id::text || '|' || new.territory_code, 0));
  if new.status='active' and exists (
    select 1 from franchise.agreement a where a.tenant_id=new.tenant_id and a.territory_code=new.territory_code and a.agreement_id<>new.agreement_id and a.status='active'
      and daterange(a.starts_on,coalesce(a.ends_on,'infinity'::date),'[)') && daterange(new.starts_on,coalesce(new.ends_on,'infinity'::date),'[)')
  ) then
    raise exception using errcode='23P01', message='active franchise territories overlap';
  end if;
  return new;
end;
$$;

create trigger franchise_territory_non_overlap before insert or update of territory_code,starts_on,ends_on,status on franchise.agreement for each row execute function franchise.enforce_active_territory_non_overlap();

create index organization_parent_status_idx on org.organization(tenant_id,parent_organization_id,status,organization_id);
create index agreement_territory_status_idx on franchise.agreement(tenant_id,territory_code,status,starts_on,agreement_id);

commit;
````

### FILE: `db/migrations/0008_franchise_network_admin.down.sql`

```yaml
block_id: "GO-FULFILLMENT-API:db-migrations-0008-franchise-network-admin-down-sql:v3"
operation: CREATE
provenance: AUTHORED
source: "local verified composition governed by pinned Microsoft BCApps network/intercompany sources"
license: "LicenseRef-Workspace-Owner"
sha256: "5d87a453e6cb7094391f9b4e881c1927b73e6afebf58702e968800d2b1e4b0ae"
variables: []
secrets_allowed: false
```

````sql
begin;
drop index if exists franchise.agreement_territory_status_idx;
drop index if exists org.organization_parent_status_idx;
drop trigger if exists franchise_territory_non_overlap on franchise.agreement;
drop function if exists franchise.enforce_active_territory_non_overlap();
drop trigger if exists organization_hierarchy_guard on org.organization;
drop function if exists org.enforce_organization_hierarchy();
alter table franchise.agreement drop column if exists updated_at, drop column if exists version;
alter table org.organization drop column if exists version;
commit;
````

### FILE: `db/tests/0008_franchise_network_admin.test.sql`

```yaml
block_id: "GO-FULFILLMENT-API:db-tests-0008-franchise-network-admin-test-sql:v3"
operation: CREATE
provenance: AUTHORED
source: "local verified composition governed by pinned Microsoft BCApps network/intercompany sources"
license: "LicenseRef-Workspace-Owner"
sha256: "08f21443ef470d45f17997f329f2195c9e82ef2bc5b02733f628d8bb7b08757e"
variables: []
secrets_allowed: false
```

````sql
begin;

do $$
begin
  if not exists(select 1 from information_schema.columns where table_schema='org' and table_name='organization' and column_name='version') then raise exception 'organization version missing'; end if;
  if not exists(select 1 from information_schema.columns where table_schema='franchise' and table_name='agreement' and column_name='version') then raise exception 'agreement version missing'; end if;
  if to_regprocedure('org.enforce_organization_hierarchy()') is null then raise exception 'hierarchy guard missing'; end if;
  if to_regprocedure('franchise.enforce_active_territory_non_overlap()') is null then raise exception 'territory guard missing'; end if;
end;
$$;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values('018f4d4a-7b36-7a21-8d10-2f4c54c28900','network-test','Network Test','Network Test');

insert into org.organization(tenant_id,organization_id,parent_organization_id,organization_code,display_name,organization_type)
values
('018f4d4a-7b36-7a21-8d10-2f4c54c28900','root',null,'root','Root','enterprise'),
('018f4d4a-7b36-7a21-8d10-2f4c54c28900','franchisee','root','franchisee','Franchisee','franchisee'),
('018f4d4a-7b36-7a21-8d10-2f4c54c28900','store','franchisee','store','Store','store');

do $$
begin
  begin
    update org.organization set parent_organization_id='store'
    where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c28900' and organization_id='root';
    raise exception 'organization hierarchy accepted a cycle';
  exception when check_violation then
    null;
  end;
end;
$$;

insert into franchise.agreement(tenant_id,agreement_id,franchise_organization_id,territory_code,terms_version,starts_on,status)
values
('018f4d4a-7b36-7a21-8d10-2f4c54c28900','active-agreement','franchisee','AR-CBA','v1','2026-01-01','active'),
('018f4d4a-7b36-7a21-8d10-2f4c54c28900','draft-agreement','franchisee','AR-CBA','v2','2026-06-01','draft');

do $$
begin
  begin
    update franchise.agreement set status='active'
    where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c28900' and agreement_id='draft-agreement';
    raise exception 'overlapping active territory accepted';
  exception when exclusion_violation then
    null;
  end;
end;
$$;

rollback;
````

### FILE: `internal/platform/postgres/fulfillment_transport.go`

```yaml
block_id: "GO-FULFILLMENT-API:internal-platform-postgres-fulfillment-transport-go:v1"
operation: CREATE
provenance: ADAPTED
path: "internal/platform/postgres/fulfillment_transport.go"
language: "go"
source: "local adapted composition governed by fixed Microsoft BCApps shipment authority and admitted Amazon provider adapters"
license: "LicenseRef-Workspace-Owner"
sha256: "2433323b464b2accff6af6e1fe35067401ad612b4ed66678dc525962024e7ac6"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"elite.local/enterprise/internal/fulfillment"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func scanCustomerTransport(row pgx.Row) (fulfillment.CustomerTransport, error) {
	var value fulfillment.CustomerTransport
	err := row.Scan(&value.ID, &value.RequestID, &value.CustomerShipmentID, &value.OrderID, &value.OriginOrganizationID, &value.DestinationCustomerSubject, &value.ProviderCode, &value.ProviderServiceCode, &value.ProviderReference, &value.State, &value.Version, &value.OrderFulfillmentState)
	return value, err
}

func (r *Fulfillment) readCustomerTransportReplay(ctx context.Context, tx pgx.Tx, tenant string, command fulfillment.CustomerTransportCommand) (fulfillment.CustomerTransport, bool, error) {
	value, err := scanCustomerTransport(tx.QueryRow(ctx, `select s.shipment_id,coalesce(s.request_id,''),coalesce(s.customer_shipment_id,''),coalesce(s.order_id,''),s.origin_organization_id,coalesce(s.destination_customer_principal_id,''),s.provider_code,coalesce(s.provider_service_code,''),coalesce(s.provider_reference,''),s.state,s.version,o.fulfillment_state from logistics.shipment s join sales.customer_order o on o.tenant_id=s.tenant_id and o.order_id=s.order_id where s.tenant_id=$1 and s.request_id=$2 for share of s,o`, tenant, command.RequestID))
	if errors.Is(err, pgx.ErrNoRows) {
		return fulfillment.CustomerTransport{}, false, nil
	}
	if err != nil {
		return fulfillment.CustomerTransport{}, false, err
	}
	if value.CustomerShipmentID != command.CustomerShipmentID || value.OriginOrganizationID != command.OriginOrganizationID || value.ProviderCode != command.ProviderCode || value.ProviderServiceCode != command.ProviderServiceCode || value.ProviderReference != command.ProviderReference {
		return fulfillment.CustomerTransport{}, false, fulfillment.ErrConflict
	}
	return value, true, nil
}

func (r *Fulfillment) CreateCustomerTransport(ctx context.Context, tenant string, ids fulfillment.TransportIDs, command fulfillment.CustomerTransportCommand) (fulfillment.CustomerTransport, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fulfillment.CustomerTransport{}, err
	}
	defer tx.Rollback(ctx)
	if replay, ok, err := r.readCustomerTransportReplay(ctx, tx, tenant, command); err != nil || ok {
		return replay, err
	}
	state := "planned"
	if command.ProviderReference != "" {
		state = "booked"
	}
	value, err := scanCustomerTransport(tx.QueryRow(ctx, `insert into logistics.shipment(tenant_id,shipment_id,provider_code,provider_service_code,provider_reference,origin_organization_id,destination_organization_id,destination_customer_principal_id,state,version,request_id,customer_shipment_id,order_id) select s.tenant_id,$3,$4,nullif($5,''),nullif($6,''),s.organization_id,null,o.customer_principal_id,$7,1,$8,s.shipment_id,s.order_id from sales.customer_shipment s join sales.customer_order o on o.tenant_id=s.tenant_id and o.order_id=s.order_id where s.tenant_id=$1 and s.shipment_id=$2 and s.organization_id=$9 returning shipment_id,request_id,customer_shipment_id,order_id,origin_organization_id,destination_customer_principal_id,provider_code,coalesce(provider_service_code,''),coalesce(provider_reference,''),state,version,(select fulfillment_state from sales.customer_order where tenant_id=$1 and order_id=logistics.shipment.order_id)`, tenant, command.CustomerShipmentID, ids.ShipmentID, command.ProviderCode, command.ProviderServiceCode, command.ProviderReference, state, command.RequestID, command.OriginOrganizationID))
	if postgresConflict(err) {
		if replay, ok, replayErr := r.readCustomerTransportReplay(ctx, tx, tenant, command); replayErr != nil || ok {
			return replay, replayErr
		}
		return fulfillment.CustomerTransport{}, fulfillment.ErrConflict
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return fulfillment.CustomerTransport{}, fulfillment.ErrConflict
	}
	if err != nil {
		return fulfillment.CustomerTransport{}, err
	}
	if err = outbox(ctx, tx, tenant, ids.EventID, "shipment", value.ID, value.Version, "shipment."+state, `jsonb_build_object('customer_shipment_id',$7::text,'order_id',$8::text,'automatic_business_delivery_acceptance',false)`, value.CustomerShipmentID, value.OrderID); err != nil {
		return fulfillment.CustomerTransport{}, err
	}
	return value, tx.Commit(ctx)
}

func providerTransitionAllowed(current, target string) bool {
	return (current == "booked" && (target == "picked-up" || target == "exception")) ||
		(current == "picked-up" && (target == "in-transit" || target == "exception")) ||
		(current == "in-transit" && (target == "in-transit" || target == "delivered" || target == "exception"))
}

func (r *Fulfillment) readProviderReportReplay(ctx context.Context, tx pgx.Tx, tenant, shipmentID string, command fulfillment.ProviderReportCommand) (fulfillment.CustomerTransport, bool, error) {
	var storedStatus, storedSchema, storedEvidence string
	var storedOccurred time.Time
	var storedAutomatic bool
	err := tx.QueryRow(ctx, `select provider_status,report_schema,evidence_sha256_hex,occurred_at,automatic_business_delivery_acceptance from logistics.shipment_provider_report where tenant_id=$1 and shipment_id=$2 and provider_event_id=$3`, tenant, shipmentID, command.ProviderEventID).Scan(&storedStatus, &storedSchema, &storedEvidence, &storedOccurred, &storedAutomatic)
	if errors.Is(err, pgx.ErrNoRows) {
		return fulfillment.CustomerTransport{}, false, nil
	}
	if err != nil {
		return fulfillment.CustomerTransport{}, false, err
	}
	if storedStatus != command.ProviderStatus || storedSchema != command.ReportSchema || storedEvidence != command.EvidenceSHA256 || !storedOccurred.Equal(command.OccurredAt) || storedAutomatic != command.AutomaticBusinessDeliveryAcceptance {
		return fulfillment.CustomerTransport{}, false, fulfillment.ErrConflict
	}
	value, err := scanCustomerTransport(tx.QueryRow(ctx, `select s.shipment_id,s.request_id,s.customer_shipment_id,s.order_id,s.origin_organization_id,s.destination_customer_principal_id,s.provider_code,coalesce(s.provider_service_code,''),coalesce(s.provider_reference,''),s.state,s.version,o.fulfillment_state from logistics.shipment s join sales.customer_order o on o.tenant_id=s.tenant_id and o.order_id=s.order_id where s.tenant_id=$1 and s.shipment_id=$2`, tenant, shipmentID))
	return value, true, err
}

func (r *Fulfillment) RecordProviderReport(ctx context.Context, tenant, shipmentID, eventID string, command fulfillment.ProviderReportCommand) (fulfillment.CustomerTransport, error) {
	target, err := fulfillment.NormalizeProviderReport(command)
	if err != nil {
		return fulfillment.CustomerTransport{}, err
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fulfillment.CustomerTransport{}, err
	}
	defer tx.Rollback(ctx)
	if replay, ok, err := r.readProviderReportReplay(ctx, tx, tenant, shipmentID, command); err != nil || ok {
		return replay, err
	}
	var current, orderID string
	var version int64
	err = tx.QueryRow(ctx, `select state,version,order_id from logistics.shipment where tenant_id=$1 and shipment_id=$2 and origin_organization_id=$3 and customer_shipment_id is not null for update`, tenant, shipmentID, command.OriginOrganizationID).Scan(&current, &version, &orderID)
	if errors.Is(err, pgx.ErrNoRows) || version != command.Version || !providerTransitionAllowed(current, target) {
		return fulfillment.CustomerTransport{}, fulfillment.ErrConflict
	}
	if err != nil {
		return fulfillment.CustomerTransport{}, err
	}
	_, err = tx.Exec(ctx, `insert into logistics.shipment_provider_report(tenant_id,shipment_id,provider_event_id,provider_status,normalized_state,report_schema,evidence_sha256_hex,occurred_at,automatic_business_delivery_acceptance,shipment_version) values($1,$2,$3,$4,$5,$6,$7,$8,false,$9)`, tenant, shipmentID, command.ProviderEventID, command.ProviderStatus, target, command.ReportSchema, command.EvidenceSHA256, command.OccurredAt, version+1)
	if postgresConflict(err) {
		return fulfillment.CustomerTransport{}, fulfillment.ErrConflict
	}
	if err != nil {
		return fulfillment.CustomerTransport{}, err
	}
	result, err := tx.Exec(ctx, `update logistics.shipment set state=$4,version=version+1,last_event_at=$5 where tenant_id=$1 and shipment_id=$2 and version=$3`, tenant, shipmentID, version, target, command.OccurredAt)
	if err != nil || result.RowsAffected() != 1 {
		if err != nil {
			return fulfillment.CustomerTransport{}, err
		}
		return fulfillment.CustomerTransport{}, fulfillment.ErrConflict
	}
	if target == "delivered" {
		_, err = tx.Exec(ctx, `update sales.customer_order o set fulfillment_state='delivered',version=version+1,updated_at=clock_timestamp() where o.tenant_id=$1 and o.order_id=$2 and o.fulfillment_state='shipped' and not exists (select 1 from sales.customer_shipment cs left join logistics.shipment ls on ls.tenant_id=cs.tenant_id and ls.customer_shipment_id=cs.shipment_id where cs.tenant_id=o.tenant_id and cs.order_id=o.order_id and (ls.shipment_id is null or ls.state<>'delivered'))`, tenant, orderID)
		if err != nil {
			return fulfillment.CustomerTransport{}, err
		}
	}
	if err = outbox(ctx, tx, tenant, eventID, "shipment", shipmentID, version+1, "shipment.provider-report-recorded", `jsonb_build_object('provider_status',$7::text,'normalized_state',$8::text,'evidence_sha256',$9::text,'automatic_business_delivery_acceptance',false)`, command.ProviderStatus, target, command.EvidenceSHA256); err != nil {
		return fulfillment.CustomerTransport{}, err
	}
	value, err := scanCustomerTransport(tx.QueryRow(ctx, `select s.shipment_id,s.request_id,s.customer_shipment_id,s.order_id,s.origin_organization_id,s.destination_customer_principal_id,s.provider_code,coalesce(s.provider_service_code,''),coalesce(s.provider_reference,''),s.state,s.version,o.fulfillment_state from logistics.shipment s join sales.customer_order o on o.tenant_id=s.tenant_id and o.order_id=s.order_id where s.tenant_id=$1 and s.shipment_id=$2`, tenant, shipmentID))
	if err != nil {
		return fulfillment.CustomerTransport{}, fmt.Errorf("read connected transport: %w", err)
	}
	return value, tx.Commit(ctx)
}
````

### FILE: `internal/platform/postgres/fulfillment_transport_integration_test.go`

```yaml
block_id: "GO-FULFILLMENT-API:internal-platform-postgres-fulfillment-transport-integration-test-go:v1"
operation: CREATE
provenance: ADAPTED
path: "internal/platform/postgres/fulfillment_transport_integration_test.go"
language: "go"
source: "local adapted regression composition"
license: "LicenseRef-Workspace-Owner"
sha256: "eb335329965bfb6fac143669742984a2055b0808ad1cdf93aea7bc55f9f7ed54"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"elite.local/enterprise/internal/fulfillment"
	"elite.local/enterprise/internal/inventorycontrol"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestConnectedCarrierDeliveryIsIdempotentConcurrentAndNeverAcceptsHandover(t *testing.T) {
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
	defer cleanupConnectedTransport(t, pool, tenant)
	for _, command := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Connected Carrier','Connected Carrier')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'warehouse','warehouse','Warehouse','warehouse')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,status) values($1,'customer','Customer','active')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'model','MODEL','Model','other','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,homologation_state,lifecycle_state) values($1,'variant','model','VARIANT','Variant','{}','approved','active')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,'order','warehouse','customer','placed','USD',100,1)`,
		`insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units) values($1,'order','line','variant',1,100)`,
	} {
		args := []any{tenant}
		if strings.Contains(command, "$2") {
			args = append(args, "carrier-"+tenant[:8])
		}
		if _, err = pool.Exec(ctx, command, args...); err != nil {
			t.Fatal(err)
		}
	}
	inventory := NewInventoryControl(pool)
	if _, err = inventory.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "item", Code: "ITEM", Description: "Item", BaseUOM: "EA", BaseRoundingPrecision: "1", TrackingMode: "none", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	for _, bin := range []inventorycontrol.WarehouseBin{{ID: "receive", OrganizationID: "warehouse", Code: "RECEIVE", Type: "receive", Version: 1}, {ID: "pick", OrganizationID: "warehouse", Code: "PICK", Type: "putpick", Ranking: 100, Version: 1}, {ID: "ship", OrganizationID: "warehouse", Code: "SHIP", Type: "ship", Version: 1}} {
		if _, err = inventory.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatal(err)
		}
	}
	for _, policy := range []inventorycontrol.ItemBinPolicy{{OrganizationID: "warehouse", ItemID: "item", BinID: "receive", Fixed: true, MinQuantity: "0", MaxQuantity: "1", Version: 1}, {OrganizationID: "warehouse", ItemID: "item", BinID: "pick", Fixed: true, Default: true, MinQuantity: "0", MaxQuantity: "1", Version: 1}, {OrganizationID: "warehouse", ItemID: "item", BinID: "ship", Fixed: true, MinQuantity: "0", MaxQuantity: "1", Version: 1}} {
		if _, err = inventory.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), policy); err != nil {
			t.Fatal(err)
		}
	}
	posting := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	receipt, err := inventory.PostWarehouseReceipt(ctx, tenant, warehousePackagingIDs(t, "carrier-source"), inventorycontrol.WarehouseReceiptCommand{RequestID: "carrier-source-receipt", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "item", HandlingUOM: "EA", HandlingQuantity: "1", UnitCost: "10", PostingDate: posting, SourceKind: "purchase", SourceID: "carrier-source-po"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = inventory.RegisterWarehouseActivity(ctx, tenant, "warehouse", receipt.PutAway.ID, 1, posting.AddDate(0, 0, 1), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	if _, err = inventory.ConfigureSalesWarehouseBinding(ctx, tenant, bulkTestUUID(t), inventorycontrol.SalesWarehouseBinding{RequestID: "carrier-binding", VariantID: "variant", ItemID: "item", SalesUOMCode: "EA", Version: 1}); err != nil {
		t.Fatal(err)
	}
	pick, err := inventory.CreateWarehousePick(ctx, tenant, warehousePackagingIDs(t, "carrier-pick"), inventorycontrol.WarehousePickCommand{RequestID: "carrier-pick", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "item", DemandKind: "customer-order", DemandID: "order", DemandLineID: "line", Quantity: "1", HandlingUOM: "EA"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = inventory.RegisterWarehouseActivity(ctx, tenant, "warehouse", pick.ID, 1, posting.AddDate(0, 0, 2), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	if _, err = inventory.PostCustomerShipment(ctx, tenant, inventorycontrol.CustomerShipmentIDs{ShipmentID: "sales-shipment", LineID: "sales-shipment-line", EventID: bulkTestUUID(t)}, inventorycontrol.CustomerShipmentCommand{RequestID: "sales-shipment-request", OrganizationID: "warehouse", OrderID: "order", WarehouseActivityID: pick.ID, PostingDate: posting.AddDate(0, 0, 3)}); err != nil {
		t.Fatal(err)
	}
	transportRepo := NewFulfillment(pool)
	command := fulfillment.CustomerTransportCommand{RequestID: "transport-request", CustomerShipmentID: "sales-shipment", OriginOrganizationID: "warehouse", ProviderCode: "amazon-easyship", ProviderServiceCode: "standard", ProviderReference: "tracking-1"}
	transport, err := transportRepo.CreateCustomerTransport(ctx, tenant, fulfillment.TransportIDs{ShipmentID: "transport", EventID: bulkTestUUID(t)}, command)
	if err != nil || transport.State != "booked" || transport.OrderFulfillmentState != "shipped" {
		t.Fatalf("transport=%+v err=%v", transport, err)
	}
	replay, err := transportRepo.CreateCustomerTransport(ctx, tenant, fulfillment.TransportIDs{ShipmentID: "ignored", EventID: bulkTestUUID(t)}, command)
	if err != nil || replay.ID != transport.ID {
		t.Fatalf("transport replay=%+v err=%v", replay, err)
	}
	divergent := command
	divergent.ProviderReference = "tracking-other"
	if _, err = transportRepo.CreateCustomerTransport(ctx, tenant, fulfillment.TransportIDs{ShipmentID: "ignored-divergent", EventID: bulkTestUUID(t)}, divergent); !errors.Is(err, fulfillment.ErrConflict) {
		t.Fatalf("divergent transport error=%v", err)
	}
	baseReport := fulfillment.ProviderReportCommand{OriginOrganizationID: "warehouse", ProviderEventID: "pickup-event", ProviderStatus: "PickedUp", ReportSchema: fulfillment.AmazonEasyShipReconciliationSchema, EvidenceSHA256: strings.Repeat("a", 64), OccurredAt: posting.AddDate(0, 0, 4), Version: 1}
	picked, err := transportRepo.RecordProviderReport(ctx, tenant, transport.ID, bulkTestUUID(t), baseReport)
	if err != nil || picked.State != "picked-up" || picked.Version != 2 {
		t.Fatalf("picked=%+v err=%v", picked, err)
	}
	start := make(chan struct{})
	results := make(chan error, 2)
	var workers sync.WaitGroup
	for index, status := range []string{"AtOriginFC", "AtDestinationFC"} {
		workers.Add(1)
		go func(index int, status string) {
			defer workers.Done()
			<-start
			report := fulfillment.ProviderReportCommand{OriginOrganizationID: "warehouse", ProviderEventID: "transit-event-" + string(rune('a'+index)), ProviderStatus: status, ReportSchema: fulfillment.AmazonEasyShipReconciliationSchema, EvidenceSHA256: strings.Repeat(string(rune('b'+index)), 64), OccurredAt: posting.AddDate(0, 0, 5), Version: 2}
			_, reportErr := transportRepo.RecordProviderReport(ctx, tenant, transport.ID, bulkTestUUID(t), report)
			results <- reportErr
		}(index, status)
	}
	close(start)
	workers.Wait()
	close(results)
	succeeded, conflicted := 0, 0
	for result := range results {
		if result == nil {
			succeeded++
		} else if errors.Is(result, fulfillment.ErrConflict) {
			conflicted++
		} else {
			t.Fatal(result)
		}
	}
	if succeeded != 1 || conflicted != 1 {
		t.Fatalf("provider concurrency success=%d conflict=%d", succeeded, conflicted)
	}
	outForDelivery := fulfillment.ProviderReportCommand{OriginOrganizationID: "warehouse", ProviderEventID: "out-for-delivery-event", ProviderStatus: "OutForDelivery", ReportSchema: fulfillment.AmazonEasyShipReconciliationSchema, EvidenceSHA256: strings.Repeat("d", 64), OccurredAt: posting.AddDate(0, 0, 6), Version: 3}
	if _, err = transportRepo.RecordProviderReport(ctx, tenant, transport.ID, bulkTestUUID(t), outForDelivery); err != nil {
		t.Fatal(err)
	}
	delivery := fulfillment.ProviderReportCommand{OriginOrganizationID: "warehouse", ProviderEventID: "delivery-event", ProviderStatus: "Delivered", ReportSchema: fulfillment.AmazonEasyShipReconciliationSchema, EvidenceSHA256: strings.Repeat("e", 64), OccurredAt: posting.AddDate(0, 0, 7), Version: 4}
	delivered, err := transportRepo.RecordProviderReport(ctx, tenant, transport.ID, bulkTestUUID(t), delivery)
	if err != nil || delivered.State != "delivered" || delivered.Version != 5 || delivered.OrderFulfillmentState != "delivered" {
		t.Fatalf("delivered=%+v err=%v", delivered, err)
	}
	replayedDelivery, err := transportRepo.RecordProviderReport(ctx, tenant, transport.ID, bulkTestUUID(t), delivery)
	if err != nil || replayedDelivery.Version != 5 {
		t.Fatalf("delivery replay=%+v err=%v", replayedDelivery, err)
	}
	divergentDelivery := delivery
	divergentDelivery.EvidenceSHA256 = strings.Repeat("f", 64)
	if _, err = transportRepo.RecordProviderReport(ctx, tenant, transport.ID, bulkTestUUID(t), divergentDelivery); !errors.Is(err, fulfillment.ErrConflict) {
		t.Fatalf("divergent delivery error=%v", err)
	}
	var handovers, reports int
	if err = pool.QueryRow(ctx, `select count(*) from sales.delivery_handover where tenant_id=$1`, tenant).Scan(&handovers); err != nil || handovers != 0 {
		t.Fatalf("provider created business handover=%d err=%v", handovers, err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from logistics.shipment_provider_report where tenant_id=$1`, tenant).Scan(&reports); err != nil || reports != 4 {
		t.Fatalf("provider reports=%d err=%v", reports, err)
	}
}

func cleanupConnectedTransport(t *testing.T, pool *pgxpool.Pool, tenant string) {
	t.Helper()
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Errorf("connected transport cleanup begin: %v", err)
		return
	}
	defer tx.Rollback(ctx)
	commands := []string{
		`alter table logistics.shipment_provider_report disable trigger shipment_provider_report_immutable`,
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
		`delete from logistics.shipment_provider_report where tenant_id=$1`,
		`delete from logistics.shipment where tenant_id=$1`,
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
		`delete from crm.customer_profile where tenant_id=$1`,
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
		`alter table logistics.shipment_provider_report enable trigger shipment_provider_report_immutable`,
	}
	for _, command := range commands {
		args := []any{}
		if strings.Contains(command, "$1") {
			args = append(args, tenant)
		}
		if _, err = tx.Exec(ctx, command, args...); err != nil {
			t.Errorf("connected transport cleanup: %v", err)
			return
		}
	}
	if err = tx.Commit(ctx); err != nil {
		t.Errorf("connected transport cleanup commit: %v", err)
	}
}
````

### FILE: `db/migrations/0039_connected_carrier_delivery.up.sql`

```yaml
block_id: "GO-FULFILLMENT-API:db-migrations-0039-connected-carrier-delivery-up-sql:v1"
operation: CREATE
provenance: ADAPTED
path: "db/migrations/0039_connected_carrier_delivery.up.sql"
language: "sql"
source: "local adapted PostgreSQL composition"
license: "LicenseRef-Workspace-Owner"
sha256: "7a2a2812ef8f1e7fca4e12f5768cfe9453c031f427cd7198a22f5ec371050e74"
variables: []
secrets_allowed: false
```

````sql
begin;

alter table logistics.shipment
  drop constraint shipment_tenant_id_provider_code_provider_reference_key,
  alter column destination_organization_id drop not null,
  add column destination_customer_principal_id text,
  add column provider_service_code text,
  add column request_id text,
  add column customer_shipment_id text,
  add column order_id text,
  add column version bigint not null default 1 check (version > 0),
  add constraint shipment_destination_exactly_one_ck check (
    (destination_organization_id is not null and destination_customer_principal_id is null)
    or (destination_organization_id is null and destination_customer_principal_id is not null)
  ),
  add constraint shipment_customer_source_pair_ck check ((customer_shipment_id is null) = (order_id is null)),
  add constraint shipment_customer_destination_fk foreign key (tenant_id,destination_customer_principal_id)
    references crm.customer_profile(tenant_id,customer_principal_id),
  add constraint shipment_customer_shipment_fk foreign key (tenant_id,customer_shipment_id)
    references sales.customer_shipment(tenant_id,shipment_id),
  add constraint shipment_customer_order_fk foreign key (tenant_id,order_id)
    references sales.customer_order(tenant_id,order_id),
  add constraint shipment_request_id_ck check (request_id is null or length(request_id) between 1 and 128),
  add constraint shipment_provider_service_code_ck check (provider_service_code is null or length(provider_service_code) between 1 and 128);

create unique index shipment_provider_reference_uq
  on logistics.shipment(tenant_id,provider_code,provider_reference)
  where provider_reference is not null;

create unique index shipment_request_uq
  on logistics.shipment(tenant_id,request_id)
  where request_id is not null;

create unique index shipment_customer_shipment_uq
  on logistics.shipment(tenant_id,customer_shipment_id)
  where customer_shipment_id is not null;

create table logistics.shipment_provider_report (
  tenant_id uuid not null,
  shipment_id text not null,
  provider_event_id text not null,
  provider_status text not null check (length(provider_status) between 1 and 64),
  normalized_state text not null check (normalized_state in ('picked-up','in-transit','delivered','exception')),
  report_schema text not null check (report_schema='elite-amazon-spapi-easyship-reconciliation-receipt/v1'),
  evidence_sha256_hex text not null check (evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  occurred_at timestamptz not null,
  received_at timestamptz not null default clock_timestamp(),
  automatic_business_delivery_acceptance boolean not null default false check (not automatic_business_delivery_acceptance),
  shipment_version bigint not null check (shipment_version > 1),
  primary key (tenant_id,shipment_id,provider_event_id),
  foreign key (tenant_id,shipment_id) references logistics.shipment(tenant_id,shipment_id),
  unique (tenant_id,shipment_id,shipment_version)
);

create function logistics.reject_shipment_provider_report_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='shipment provider report evidence is immutable';
end;
$function$;

create trigger shipment_provider_report_immutable
before update or delete on logistics.shipment_provider_report
for each row execute function logistics.reject_shipment_provider_report_mutation();

create index shipment_provider_report_time_idx
  on logistics.shipment_provider_report(tenant_id,shipment_id,occurred_at,provider_event_id);

commit;
````

### FILE: `db/migrations/0039_connected_carrier_delivery.down.sql`

```yaml
block_id: "GO-FULFILLMENT-API:db-migrations-0039-connected-carrier-delivery-down-sql:v1"
operation: CREATE
provenance: ADAPTED
path: "db/migrations/0039_connected_carrier_delivery.down.sql"
language: "sql"
source: "local adapted PostgreSQL rollback composition"
license: "LicenseRef-Workspace-Owner"
sha256: "2fadfc8c4e3e84a9a3c261811faf8a12c30541382fea805169c9be2cd47b5e3f"
variables: []
secrets_allowed: false
```

````sql
begin;

drop index if exists logistics.shipment_provider_report_time_idx;
drop trigger if exists shipment_provider_report_immutable on logistics.shipment_provider_report;
drop function if exists logistics.reject_shipment_provider_report_mutation();
drop table if exists logistics.shipment_provider_report;

drop index if exists logistics.shipment_customer_shipment_uq;
drop index if exists logistics.shipment_request_uq;
drop index if exists logistics.shipment_provider_reference_uq;

alter table logistics.shipment
  drop constraint if exists shipment_provider_service_code_ck,
  drop constraint if exists shipment_request_id_ck,
  drop constraint if exists shipment_customer_order_fk,
  drop constraint if exists shipment_customer_shipment_fk,
  drop constraint if exists shipment_customer_destination_fk,
  drop constraint if exists shipment_customer_source_pair_ck,
  drop constraint if exists shipment_destination_exactly_one_ck,
  drop column if exists version,
  drop column if exists order_id,
  drop column if exists customer_shipment_id,
  drop column if exists request_id,
  drop column if exists provider_service_code,
  drop column if exists destination_customer_principal_id,
  alter column destination_organization_id set not null,
  add unique nulls not distinct (tenant_id,provider_code,provider_reference);

commit;
````

### FILE: `db/tests/0039_connected_carrier_delivery.test.sql`

```yaml
block_id: "GO-FULFILLMENT-API:db-tests-0039-connected-carrier-delivery-test-sql:v1"
operation: CREATE
provenance: ADAPTED
path: "db/tests/0039_connected_carrier_delivery.test.sql"
language: "sql"
source: "local adapted PostgreSQL regression composition"
license: "LicenseRef-Workspace-Owner"
sha256: "593693d2f89f1a1292a31683c04216322b5d5437cabe9391630a84beb5d64727"
variables: []
secrets_allowed: false
```

````sql
begin;

do $$
begin
  if to_regclass('logistics.shipment_provider_report') is null then raise exception 'provider report table missing'; end if;
  if not exists(select 1 from information_schema.columns where table_schema='logistics' and table_name='shipment' and column_name='customer_shipment_id') then raise exception 'customer shipment link missing'; end if;
  if not exists(select 1 from pg_indexes where schemaname='logistics' and indexname='shipment_provider_reference_uq' and indexdef like '%WHERE (provider_reference IS NOT NULL)%') then raise exception 'partial provider reference identity missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='shipment_provider_report_immutable') then raise exception 'provider report immutability missing'; end if;
end;
$$;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values('018f4d4a-7b36-7a21-8d10-2f4c54c29939','carrier-v174','Carrier V174','Carrier V174');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
values
('018f4d4a-7b36-7a21-8d10-2f4c54c29939','origin','origin','Origin','store'),
('018f4d4a-7b36-7a21-8d10-2f4c54c29939','destination','destination','Destination','store');

insert into logistics.shipment(tenant_id,shipment_id,provider_code,origin_organization_id,destination_organization_id,state)
values
('018f4d4a-7b36-7a21-8d10-2f4c54c29939','null-ref-a','carrier','origin','destination','planned'),
('018f4d4a-7b36-7a21-8d10-2f4c54c29939','null-ref-b','carrier','origin','destination','planned');

insert into logistics.shipment(tenant_id,shipment_id,provider_code,provider_reference,origin_organization_id,destination_organization_id,state)
values('018f4d4a-7b36-7a21-8d10-2f4c54c29939','ref-a','carrier','provider-1','origin','destination','booked');

do $$
begin
  begin
    insert into logistics.shipment(tenant_id,shipment_id,provider_code,provider_reference,origin_organization_id,destination_organization_id,state)
    values('018f4d4a-7b36-7a21-8d10-2f4c54c29939','ref-b','carrier','provider-1','origin','destination','booked');
    raise exception 'duplicate provider reference accepted';
  exception when unique_violation then null;
  end;
end;
$$;

rollback;
````

### FILE: `docs/logistics/MICROSOFT_BC_AMAZON_CONNECTED_CARRIER_DERIVATION.md`

```yaml
block_id: "GO-FULFILLMENT-API:docs-logistics-microsoft-bc-amazon-connected-carrier-derivation-md:v1"
operation: CREATE
provenance: AUTHORED
path: "docs/logistics/MICROSOFT_BC_AMAZON_CONNECTED_CARRIER_DERIVATION.md"
language: "markdown"
source: "local authored provenance and non-claim record"
license: "LicenseRef-Workspace-Owner"
sha256: "cd2338a09fce64142736f1a89998daef3a18cb4f25289d234b6bd427dd47170a"
variables: []
secrets_allowed: false
```

````markdown
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
````

## 6. Configuration surface

No standalone environment variables are introduced. The module consumes the composed database, ID source and verified authorization. `provider_code`, optional `provider_service_code`, exact provider schema, provider status, event time, evidence SHA-256, shipment/order/customer identity and optimistic version are strict request inputs. Secrets, webhook verification, polling, provider account/region, retention and customer acceptance remain selected project adapter/configuration inputs; they are never inferred from a carrier report.

## 7. Dependency bill

| Dependency | Pin | Use | License | Scope | Source |
|---|---|---|---|---|---|
| Go | `1.26.7` | domain/HTTP/tests | BSD-3-Clause | build/runtime | `go.dev` |
| PostgreSQL | `18.6` verified baseline | durable state/outbox | PostgreSQL | runtime/test | `postgresql.org` |
| `pgx` | `5.10.0` | PostgreSQL adapter | MIT | build/runtime | `github.com/jackc/pgx` |
| Microsoft BCApps | commit `2eae56d704a1fd035d104f333602aea7091b7749`, tree `f9846fb1254c6311985c6131bb2af11c7f169e1b` | autoridad estrecha de shipping agent/service, package tracking y persistencia al posted sales shipment; también conserva la autoridad previa de red | MIT | diseño/test | `github.com/microsoft/BCApps` |
| Amazon SP-API SDK/models | SDK `8e792ae345a8d334ccdbdd03181f05f040e6a4fc`, models `8e429486005c4ebdce5099e48cc48515a65359bb` | estados Easy Ship y límites de tracking/evidencia mediante tres adapters ya admitidos; no gobierna aceptación comercial | Apache-2.0/avisos del snapshot | adapter condicionado | `github.com/amzn` + referencias oficiales SP-API |

## 8. Apply order

Compose after foundation, electromobility, supply `0.16.x` and commerce schemas/modules. Apply migrations in numeric order through `0039`, then wire the application root. A customer shipment must already exist before creating its transport. Execute transition, replay, provider-schema/status, evidence, concurrency, scope, hierarchy, territory and PostgreSQL integration tests before enabling an adapter. Rollback disables report intake/routes/workers first; `0039` down requires no connected transports/reports and restores the previous null-reference constraint only as a development reversal. Production migrates forward compatibly and retains durable reports.

## 9. Verification

Compose with the Go application and migrations `0001`–`0039`. Run format, complete unit/HTTP/PostgreSQL integration, SQL tests, vet and build. Revision `0.4.0` preserves network/service boundaries and proves: two null provider references coexist; a repeated non-null reference fails; transport request replay is exact; divergent replay fails; concurrent provider reports have one version winner; exact report replay converges; divergent evidence fails; the order becomes delivered only after the connected transport reaches delivered; four reports remain immutable; zero delivery handovers are created. Execute `0039` down/up and its focal test. Production still requires a chosen carrier account/sandbox, authenticated webhook or polling lane, reconciliation/outage policy, address/privacy/retention policy, load/resilience/security, serial/specific customer shipment where required, and the separate checklist/customer-acceptance workflow.

## 10. Reconstruction evidence

Clean reconstruction, authorization and PostgreSQL integration are recorded in `reconstruction_evidence/GO_FULFILLMENT_SERVICE_FRANCHISE_API_2026-08-24_V1.md`; network hierarchy and exact source authority for version `0.3.0` are recorded in `reconstruction_evidence/FRANCHISE_NETWORK_ADMINISTRATION_2026-08-30_V118.md`; fixed Microsoft/Amazon authority, migration `0039`, connected transport/report concurrency and the no-automatic-handover invariant for version `0.4.0` are recorded in `reconstruction_evidence/MICROSOFT_BC_AMAZON_CONNECTED_CARRIER_DELIVERY_2026-09-02_V174.md`.

V402 composed delta: Connected warranty reuses existing transaction/approval/stock/service owners; SQL ordering and public wrapper behavior retained. Optional host factory fails closed. Exact source tested in WARRANTY_INTERFACE_AND_PORTABILITY_V402.md; no new dependency or corporate attribution.

V402 composed delta: T2804 network role: original four fulfillment SQL bodies extracted unchanged into one transaction with immutable result; optional host, forms and GET recovery. Migration0077, no dependency/domain-rule change. NETWORK_ROLE_RELEASE_V402.md.
