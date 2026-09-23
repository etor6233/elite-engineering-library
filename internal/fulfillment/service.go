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
