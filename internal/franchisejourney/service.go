package franchisejourney

import (
	"context"
	"elite.local/enterprise/internal/businesspolicy"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalid  = errors.New("invalid franchise journey input")
	ErrConflict = errors.New("franchise journey conflict")
	ErrNotFound = errors.New("franchise journey resource not found")
)

type Location struct {
	OrganizationID string `json:"organization_id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	City           string `json:"city"`
	Region         string `json:"region"`
	Country        string `json:"country"`
	ContactPhone   string `json:"contact_phone,omitempty"`
	ContactEmail   string `json:"contact_email,omitempty"`
}

type Appointment struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	LeadID         string    `json:"lead_id"`
	ModelID        string    `json:"model_id,omitempty"`
	Kind           string    `json:"kind"`
	StartsAt       time.Time `json:"starts_at"`
	State          string    `json:"state"`
	Version        int64     `json:"version"`
	SlotID         string    `json:"slot_id,omitempty"`
	EndsAt         time.Time `json:"ends_at,omitempty"`
	ResourceID     string    `json:"resource_id,omitempty"`
}

type ServiceResource struct {
	ID               string   `json:"id"`
	OrganizationID   string   `json:"organization_id"`
	PrincipalSubject string   `json:"principal_subject,omitempty"`
	DisplayName      string   `json:"display_name"`
	Kind             string   `json:"kind"`
	Status           string   `json:"status"`
	Version          int64    `json:"version"`
	Skills           []string `json:"skills"`
}

type AppointmentAgenda struct {
	Appointments []Appointment     `json:"appointments"`
	Resources    []ServiceResource `json:"resources"`
	Truncated    bool              `json:"truncated"`
}

type AppointmentSlot struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Kind           string    `json:"kind"`
	StartsAt       time.Time `json:"starts_at"`
	EndsAt         time.Time `json:"ends_at"`
	Capacity       int       `json:"capacity"`
	Booked         int64     `json:"booked"`
	State          string    `json:"state"`
	Version        int64     `json:"version"`
}

type Lead struct {
	ID              string    `json:"id"`
	OrganizationID  string    `json:"organization_id"`
	ModelID         string    `json:"model_id,omitempty"`
	State           string    `json:"state"`
	SourceCode      string    `json:"source_code"`
	AssignedSubject string    `json:"assigned_subject,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	Version         int64     `json:"version"`
}

type Quote struct {
	ID              string    `json:"id"`
	OrganizationID  string    `json:"organization_id"`
	LeadID          string    `json:"lead_id"`
	CustomerSubject string    `json:"customer_subject,omitempty"`
	VariantID       string    `json:"variant_id"`
	Currency        string    `json:"currency"`
	TotalMinorUnits int64     `json:"total_minor_units"`
	ValidUntil      time.Time `json:"valid_until"`
	State           string    `json:"state"`
	Version         int64     `json:"version"`
	PriceBookID     string    `json:"price_book_id"`
	OrderID         string    `json:"order_id,omitempty"`
}

type Handover struct {
	ID                   string          `json:"id"`
	OrganizationID       string          `json:"organization_id"`
	OrderID              string          `json:"order_id"`
	CustomerSubject      string          `json:"customer_subject"`
	StockUnitID          string          `json:"stock_unit_id"`
	State                string          `json:"state"`
	Version              int64           `json:"version"`
	CustomerAcceptedAt   *time.Time      `json:"customer_accepted_at,omitempty"`
	AcceptanceEvidence   string          `json:"acceptance_evidence_sha256,omitempty"`
	ChecklistID          string          `json:"checklist_id,omitempty"`
	ChecklistVersion     int64           `json:"checklist_version,omitempty"`
	ChecklistTitle       string          `json:"checklist_title,omitempty"`
	ChecklistItems       []ChecklistItem `json:"checklist_items"`
	ChecklistCompletedAt *time.Time      `json:"checklist_completed_at,omitempty"`
	SupersedesHandoverID string          `json:"supersedes_handover_id,omitempty"`
}

type ChecklistItem struct {
	ID           string `json:"id"`
	Ordinal      int    `json:"ordinal"`
	Prompt       string `json:"prompt"`
	ResponseType string `json:"response_type"`
	Required     bool   `json:"required"`
}

type DeliveryChecklist struct {
	ID             string          `json:"id"`
	OrganizationID string          `json:"organization_id"`
	Version        int64           `json:"version"`
	Title          string          `json:"title"`
	State          string          `json:"state"`
	Items          []ChecklistItem `json:"items"`
}

type ChecklistResponse struct {
	ItemID         string `json:"item_id"`
	ResponseText   string `json:"response_text"`
	EvidenceSHA256 string `json:"evidence_sha256,omitempty"`
}

type DeliveryException struct {
	ID                    string     `json:"id"`
	OrganizationID        string     `json:"organization_id"`
	HandoverID            string     `json:"handover_id"`
	CustomerSubject       string     `json:"customer_subject"`
	ReasonCode            string     `json:"reason_code"`
	Details               string     `json:"details"`
	State                 string     `json:"state"`
	Version               int64      `json:"version"`
	CreatedAt             time.Time  `json:"created_at"`
	ResolvedAt            *time.Time `json:"resolved_at,omitempty"`
	ResolutionAction      string     `json:"resolution_action,omitempty"`
	SuccessorHandoverID   string     `json:"successor_handover_id,omitempty"`
	ReturnAuthorizationID string     `json:"return_authorization_id,omitempty"`
}

type DeliveryResolution struct {
	Exception             DeliveryException `json:"exception"`
	SuccessorHandover     *Handover         `json:"successor_handover,omitempty"`
	ReturnAuthorizationID string            `json:"return_authorization_id,omitempty"`
	Disposition           string            `json:"disposition,omitempty"`
}

type ReturnReceipt struct {
	ID                   string    `json:"id"`
	AuthorizationID      string    `json:"authorization_id"`
	OrganizationID       string    `json:"organization_id"`
	OrderID              string    `json:"order_id"`
	StockUnitID          string    `json:"stock_unit_id"`
	CustomerSubject      string    `json:"customer_subject"`
	ReceivedSerialNumber string    `json:"received_serial_number"`
	ConditionCode        string    `json:"condition_code"`
	Notes                string    `json:"notes"`
	EvidenceSHA256       string    `json:"evidence_sha256"`
	ReceivedBySubject    string    `json:"received_by_subject"`
	ReceivedAt           time.Time `json:"received_at"`
}

type ReturnEffectRequest struct {
	ID             string    `json:"id"`
	EffectKind     string    `json:"effect_kind"`
	OwnerContext   string    `json:"owner_context"`
	State          string    `json:"state"`
	IdempotencyKey string    `json:"idempotency_key"`
	RequestedAt    time.Time `json:"requested_at"`
}

type ReturnDisposition struct {
	ID               string                `json:"id"`
	ReceiptID        string                `json:"receipt_id"`
	InventoryAction  string                `json:"inventory_action"`
	CustomerRemedy   string                `json:"customer_remedy"`
	Notes            string                `json:"notes"`
	DecidedBySubject string                `json:"decided_by_subject"`
	DecidedAt        time.Time             `json:"decided_at"`
	Effects          []ReturnEffectRequest `json:"effect_requests"`
}

type ReturnCase struct {
	AuthorizationID  string             `json:"authorization_id"`
	OrganizationID   string             `json:"organization_id"`
	OrderID          string             `json:"order_id"`
	StockUnitID      string             `json:"stock_unit_id"`
	CustomerSubject  string             `json:"customer_subject"`
	AuthorizedAction string             `json:"authorized_action"`
	AuthorizedAt     time.Time          `json:"authorized_at"`
	Receipt          *ReturnReceipt     `json:"receipt,omitempty"`
	Disposition      *ReturnDisposition `json:"disposition,omitempty"`
}

type CustomerJourney struct {
	Appointments []Appointment       `json:"appointments"`
	Quotes       []Quote             `json:"quotes"`
	Handovers    []Handover          `json:"handovers"`
	Exceptions   []DeliveryException `json:"delivery_exceptions"`
}

type Page[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
}

type Repository interface {
	AppointmentAgenda(context.Context, string, string, time.Time, time.Time) (AppointmentAgenda, error)
	PublicLocations(context.Context, string) ([]Location, error)
	PublicAppointmentSlots(context.Context, string, string, string, time.Time, time.Time) ([]AppointmentSlot, error)
	CreateAppointmentSlot(context.Context, string, AppointmentSlot, string) (AppointmentSlot, error)
	RequestAppointment(context.Context, string, string, string, Appointment, string, string) (Appointment, bool, error)
	CreateServiceResource(context.Context, string, ServiceResource, string) (ServiceResource, error)
	AssignAppointmentResource(context.Context, string, string, string, string, int64, string) (Appointment, error)
	CreateAvailability(context.Context, string, string, AvailabilityEntry, string) (AvailabilityEntry, error)
	CancelAvailability(context.Context, string, string, string, int64, string, string, string) (AvailabilityEntry, error)
	Availability(context.Context, string, string, string, time.Time, time.Time) ([]AvailabilityEntry, error)
	TransitionAppointment(context.Context, string, string, string, string, string, int64, string, string, string) (Appointment, error)
	CancelCustomerAppointment(context.Context, string, string, string, string, int64, string, string, string) (Appointment, error)
	Leads(context.Context, string, string, int, string) (Page[Lead], error)
	AssignLead(context.Context, string, string, string, string, int64, string) (Lead, error)
	TransitionLead(context.Context, string, string, string, string, string, int64, string) (Lead, error)
	CreateQuote(context.Context, string, string, Quote, string, string) (Quote, bool, error)
	PublishDeliveryChecklist(context.Context, string, string, DeliveryChecklist, string) (DeliveryChecklist, error)
	CompleteDeliveryChecklist(context.Context, string, string, string, string, int64, string, int64, []ChecklistResponse, string) (Handover, error)
	RejectHandover(context.Context, string, string, string, string, int64, string, string, string, string, string) (DeliveryException, error)
	DeliveryExceptions(context.Context, string, string, int) ([]DeliveryException, error)
	ResolveDeliveryException(context.Context, string, string, string, string, int64, string, string, string, string, string, string) (DeliveryResolution, error)
	ReturnCases(context.Context, string, string, int) ([]ReturnCase, error)
	ReceiveReturn(context.Context, string, string, string, string, string, string, string, string, string, string) (ReturnReceipt, error)
	DecideReturn(context.Context, string, string, string, string, string, string, string, string, string, string, string, string) (ReturnDisposition, error)
	AcceptQuote(context.Context, string, string, string, string, int64, string, string, string, string, string) (Quote, error)
	CustomerJourney(context.Context, string, string, string) (CustomerJourney, error)
	AcceptHandover(context.Context, string, string, string, string, int64, string, string, int64, string, string) (Handover, error)
}

type IDGenerator interface{ New() string }
type Clock interface{ Now() time.Time }

type Service struct {
	repository Repository
	ids        IDGenerator
	clock      Clock
	policy     *businesspolicy.Profile
}

func NewService(repository Repository, ids IDGenerator, clock Clock) *Service {
	policy := businesspolicy.Reference()
	if bound, ok := repository.(interface{ BusinessPolicySHA256() string }); ok && bound.BusinessPolicySHA256() != policy.SHA256() {
		panic("custom repository policy requires NewServiceWithProfile")
	}
	return &Service{repository: repository, ids: ids, clock: clock, policy: policy}
}

// NewServiceWithProfile binds the same validated policy used by persistence.
// Compatibility constructors retain the reference profile and historical keys.
func NewServiceWithProfile(repository Repository, ids IDGenerator, clock Clock, policy *businesspolicy.Profile) (*Service, error) {
	bound, ok := repository.(interface{ BusinessPolicySHA256() string })
	if !policy.Valid() || !ok || bound.BusinessPolicySHA256() != policy.SHA256() || ids == nil || clock == nil {
		return nil, businesspolicy.ErrProfile
	}
	return &Service{repository: repository, ids: ids, clock: clock, policy: policy}, nil
}

func (s *Service) PublicLocations(ctx context.Context, tenantCode string) ([]Location, error) {
	if !code(tenantCode) {
		return nil, ErrInvalid
	}
	return s.repository.PublicLocations(ctx, tenantCode)
}

func (s *Service) PublicAppointmentSlots(ctx context.Context, tenantCode, organizationCode, kind string, from, to time.Time) ([]AppointmentSlot, error) {
	if !code(tenantCode) || !code(organizationCode) || !appointmentKinds[kind] || from.Before(s.clock.Now()) || !to.After(from) || to.Sub(from) > 31*24*time.Hour {
		return nil, ErrInvalid
	}
	return s.repository.PublicAppointmentSlots(ctx, tenantCode, organizationCode, kind, from, to)
}

func (s *Service) CreateAppointmentSlot(ctx context.Context, tenant string, value AppointmentSlot) (AppointmentSlot, error) {
	value, err := s.prepareAppointmentSlot(tenant, value)
	if err != nil {
		return AppointmentSlot{}, err
	}
	return s.repository.CreateAppointmentSlot(ctx, tenant, value, s.ids.New())
}

func (s *Service) prepareAppointmentSlot(tenant string, value AppointmentSlot) (AppointmentSlot, error) {
	value.ID = s.ids.New()
	value.State = "open"
	value.Version = 1
	if tenant == "" || value.OrganizationID == "" || !appointmentKinds[value.Kind] || !s.policy.AllowsSlot(value.StartsAt, value.EndsAt, s.clock.Now(), value.Capacity) {
		return AppointmentSlot{}, ErrInvalid
	}
	return value, nil
}

func (s *Service) RequestAppointment(ctx context.Context, tenantCode, organizationCode, idempotencyKey, requestHash string, value Appointment) (Appointment, bool, error) {
	value.ID = s.ids.New()
	value.OrganizationID = ""
	value.State = "requested"
	value.Version = 1
	if !code(tenantCode) || !code(organizationCode) || len(idempotencyKey) < 16 || !sha256Hex(requestHash) || value.LeadID == "" || !appointmentKinds[value.Kind] || value.StartsAt.Before(s.clock.Now().Add(s.policy.LeadTime())) {
		return Appointment{}, false, ErrInvalid
	}
	return s.repository.RequestAppointment(ctx, tenantCode, organizationCode, idempotencyKey, value, requestHash, s.ids.New())
}

func (s *Service) CreateServiceResource(ctx context.Context, tenant string, value ServiceResource) (ServiceResource, error) {
	value, err := s.prepareServiceResource(tenant, value)
	if err != nil {
		return ServiceResource{}, err
	}
	return s.repository.CreateServiceResource(ctx, tenant, value, s.ids.New())
}

func (s *Service) prepareServiceResource(tenant string, value ServiceResource) (ServiceResource, error) {
	value.ID = s.ids.New()
	value.Status = "active"
	value.Version = 1
	if tenant == "" || value.OrganizationID == "" || len(strings.TrimSpace(value.DisplayName)) < 1 || len(value.DisplayName) > 160 || !resourceKinds[value.Kind] || len(value.Skills) < 1 || len(value.Skills) > 4 || ((value.Kind == "employee" || value.Kind == "contractor") != (value.PrincipalSubject != "")) {
		return ServiceResource{}, ErrInvalid
	}
	seen := map[string]bool{}
	for _, skill := range value.Skills {
		if !appointmentKinds[skill] || seen[skill] {
			return ServiceResource{}, ErrInvalid
		}
		seen[skill] = true
	}
	return value, nil
}

func (s *Service) AssignAppointmentResource(ctx context.Context, tenant, organization, appointment, resource string, version int64) (Appointment, error) {
	if tenant == "" || organization == "" || appointment == "" || resource == "" || version < 1 {
		return Appointment{}, ErrInvalid
	}
	return s.repository.AssignAppointmentResource(ctx, tenant, organization, appointment, resource, version, s.ids.New())
}

var appointmentTransitions = map[string]map[string]bool{"requested": {"confirmed": true, "cancelled": true}, "confirmed": {"completed": true, "cancelled": true, "no-show": true}}

func (s *Service) TransitionAppointment(ctx context.Context, tenant, organization, appointment, current, target string, version int64, subject, reason string) (Appointment, error) {
	if tenant == "" || organization == "" || appointment == "" || version < 1 || subject == "" || !appointmentTransitions[current][target] || ((target == "cancelled" || target == "no-show") && !code(reason)) || ((target != "cancelled" && target != "no-show") && reason != "") {
		return Appointment{}, ErrInvalid
	}
	return s.repository.TransitionAppointment(ctx, tenant, organization, appointment, current, target, version, subject, reason, s.ids.New())
}

func (s *Service) Leads(ctx context.Context, tenant, organization string, limit int, after string) (Page[Lead], error) {
	if tenant == "" || organization == "" || limit < 1 || limit > 100 {
		return Page[Lead]{}, ErrInvalid
	}
	return s.repository.Leads(ctx, tenant, organization, limit, after)
}

func (s *Service) AssignLead(ctx context.Context, tenant, organization, lead, subject string, version int64) (Lead, error) {
	if tenant == "" || organization == "" || lead == "" || subject == "" || version < 1 {
		return Lead{}, ErrInvalid
	}
	return s.repository.AssignLead(ctx, tenant, organization, lead, subject, version, s.ids.New())
}

// HTTP commands require actor-aware persistence; legacy callers retain their API.
type auditedLeadRepository interface {
	AssignLeadAs(context.Context, string, string, string, string, int64, string, string) (Lead, error)
	TransitionLeadAs(context.Context, string, string, string, string, string, int64, string, string) (Lead, error)
}

func (s *Service) AssignLeadAs(ctx context.Context, tenant, organization, lead, subject string, version int64, actor string) (Lead, error) {
	repo, ok := s.repository.(auditedLeadRepository)
	if !ok || actor == "" || tenant == "" || organization == "" || lead == "" || subject == "" || version < 1 {
		return Lead{}, ErrInvalid
	}
	return repo.AssignLeadAs(ctx, tenant, organization, lead, subject, version, s.ids.New(), actor)
}

func (s *Service) TransitionLeadAs(ctx context.Context, tenant, organization, lead, current, target string, version int64, actor string) (Lead, error) {
	repo, ok := s.repository.(auditedLeadRepository)
	if !ok || actor == "" || tenant == "" || organization == "" || lead == "" || version < 1 || !leadTransitions[current][target] {
		return Lead{}, ErrInvalid
	}
	return repo.TransitionLeadAs(ctx, tenant, organization, lead, current, target, version, s.ids.New(), actor)
}

var leadTransitions = map[string]map[string]bool{
	"new":       {"contacted": true, "lost": true},
	"contacted": {"qualified": true, "lost": true},
	"qualified": {"converted": true, "lost": true},
}

func (s *Service) TransitionLead(ctx context.Context, tenant, organization, lead, current, target string, version int64) (Lead, error) {
	if tenant == "" || organization == "" || lead == "" || version < 1 || !leadTransitions[current][target] {
		return Lead{}, ErrInvalid
	}
	return s.repository.TransitionLead(ctx, tenant, organization, lead, current, target, version, s.ids.New())
}

func (s *Service) CreateQuote(ctx context.Context, tenant, idempotencyKey, requestHash string, value Quote) (Quote, bool, error) {
	return s.createQuote(ctx, tenant, idempotencyKey, requestHash, value, "")
}
func (s *Service) CreateQuoteAs(ctx context.Context, tenant, subject, key, hash string, value Quote) (Quote, bool, error) {
	if subject == "" {
		return Quote{}, false, ErrInvalid
	}
	return s.createQuote(ctx, tenant, key, hash, value, subject)
}
func (s *Service) createQuote(ctx context.Context, tenant, idempotencyKey, requestHash string, value Quote, actor string) (Quote, bool, error) {
	value.ID = s.ids.New()
	value.State = "issued"
	value.Version = 1
	now := s.clock.Now()
	if tenant == "" || len(idempotencyKey) < 16 || !sha256Hex(requestHash) || value.OrganizationID == "" || value.LeadID == "" || value.VariantID == "" || value.PriceBookID == "" || !value.ValidUntil.After(now) {
		return Quote{}, false, ErrInvalid
	}
	if actor != "" {
		repository, ok := s.repository.(interface {
			CreateQuoteAs(context.Context, string, string, Quote, string, string, string) (Quote, bool, error)
		})
		if !ok {
			return Quote{}, false, ErrConflict
		}
		return repository.CreateQuoteAs(ctx, tenant, idempotencyKey, value, requestHash, s.ids.New(), actor)
	}
	return s.repository.CreateQuote(ctx, tenant, idempotencyKey, value, requestHash, s.ids.New())
}

func (s *Service) CustomerJourney(ctx context.Context, tenant, organization, customer string) (CustomerJourney, error) {
	if tenant == "" || organization == "" || customer == "" {
		return CustomerJourney{}, ErrInvalid
	}
	return s.repository.CustomerJourney(ctx, tenant, organization, customer)
}

func (s *Service) AppointmentAgenda(ctx context.Context, tenant, organization string, from, to time.Time) (AppointmentAgenda, error) {
	if tenant == "" || organization == "" || from.IsZero() || !to.After(from) || to.Sub(from) > 31*24*time.Hour {
		return AppointmentAgenda{}, ErrInvalid
	}
	return s.repository.AppointmentAgenda(ctx, tenant, organization, from, to)
}

var checklistResponseTypes = map[string]bool{"confirmation": true, "text": true, "serial": true, "evidence": true}

func (s *Service) PublishDeliveryChecklist(ctx context.Context, tenant, subject string, value DeliveryChecklist) (DeliveryChecklist, error) {
	value.State = "published"
	if tenant == "" || subject == "" || value.OrganizationID == "" || !code(value.ID) || value.Version < 1 || len(strings.TrimSpace(value.Title)) < 1 || len(value.Title) > 160 || len(value.Items) < 1 || len(value.Items) > 64 {
		return DeliveryChecklist{}, ErrInvalid
	}
	seen := map[string]bool{}
	for index := range value.Items {
		item := &value.Items[index]
		item.Ordinal = index + 1
		if !code(item.ID) || seen[item.ID] || len(strings.TrimSpace(item.Prompt)) < 1 || len(item.Prompt) > 500 || !checklistResponseTypes[item.ResponseType] {
			return DeliveryChecklist{}, ErrInvalid
		}
		seen[item.ID] = true
	}
	return s.repository.PublishDeliveryChecklist(ctx, tenant, subject, value, s.ids.New())
}

func (s *Service) CompleteDeliveryChecklist(ctx context.Context, tenant, organization, subject, handover string, version int64, checklistID string, checklistVersion int64, responses []ChecklistResponse) (Handover, error) {
	if tenant == "" || organization == "" || subject == "" || handover == "" || version < 1 || !code(checklistID) || checklistVersion < 1 || len(responses) < 1 || len(responses) > 64 {
		return Handover{}, ErrInvalid
	}
	seen := map[string]bool{}
	for _, response := range responses {
		if !code(response.ItemID) || seen[response.ItemID] || len(strings.TrimSpace(response.ResponseText)) < 1 || len(response.ResponseText) > 2048 || (response.EvidenceSHA256 != "" && !sha256Hex(response.EvidenceSHA256)) {
			return Handover{}, ErrInvalid
		}
		seen[response.ItemID] = true
	}
	return s.repository.CompleteDeliveryChecklist(ctx, tenant, organization, subject, handover, version, checklistID, checklistVersion, responses, s.ids.New())
}

func (s *Service) RejectHandover(ctx context.Context, tenant, organization, customer, handover string, version int64, reasonCode, details, evidence string) (DeliveryException, error) {
	if tenant == "" || organization == "" || customer == "" || handover == "" || version < 1 || !code(reasonCode) || len(strings.TrimSpace(details)) < 1 || len(details) > 1000 || !sha256Hex(evidence) {
		return DeliveryException{}, ErrInvalid
	}
	return s.repository.RejectHandover(ctx, tenant, organization, customer, handover, version, reasonCode, details, evidence, s.ids.New(), s.ids.New())
}

func (s *Service) DeliveryExceptions(ctx context.Context, tenant, organization string, limit int) ([]DeliveryException, error) {
	if tenant == "" || organization == "" || limit < 1 || limit > 100 {
		return nil, ErrInvalid
	}
	return s.repository.DeliveryExceptions(ctx, tenant, organization, limit)
}

var deliveryResolutionActions = map[string]bool{"correct-and-represent": true, "return": true, "exchange": true}

func (s *Service) ResolveDeliveryException(ctx context.Context, tenant, organization, subject, exceptionID string, version int64, action, notes string) (DeliveryResolution, error) {
	if tenant == "" || organization == "" || subject == "" || exceptionID == "" || version < 1 || !deliveryResolutionActions[action] || len(strings.TrimSpace(notes)) < 1 || len(notes) > 1000 {
		return DeliveryResolution{}, ErrInvalid
	}
	return s.repository.ResolveDeliveryException(ctx, tenant, organization, subject, exceptionID, version, action, notes, s.ids.New(), s.ids.New(), s.ids.New(), s.ids.New())
}

func (s *Service) ReturnCases(ctx context.Context, tenant, organization string, limit int) ([]ReturnCase, error) {
	if tenant == "" || organization == "" || limit < 1 || limit > 100 {
		return nil, ErrInvalid
	}
	return s.repository.ReturnCases(ctx, tenant, organization, limit)
}

var returnConditionCodes = map[string]bool{"sealed": true, "opened": true, "damaged": true, "incomplete": true}

func (s *Service) ReceiveReturn(ctx context.Context, tenant, organization, subject, authorizationID, serialNumber, conditionCode, notes, evidence string) (ReturnReceipt, error) {
	if tenant == "" || organization == "" || subject == "" || authorizationID == "" || len(serialNumber) < 1 || len(serialNumber) > 128 || !returnConditionCodes[conditionCode] || len(strings.TrimSpace(notes)) < 1 || len(notes) > 1000 || !sha256Hex(evidence) {
		return ReturnReceipt{}, ErrInvalid
	}
	return s.repository.ReceiveReturn(ctx, tenant, organization, subject, authorizationID, serialNumber, conditionCode, notes, evidence, s.ids.New(), s.ids.New())
}

var returnInventoryActions = map[string]bool{"quarantine": true, "restock": true, "repair": true, "scrap": true}

func (s *Service) DecideReturn(ctx context.Context, tenant, organization, subject, receiptID, inventoryAction, notes string) (ReturnDisposition, error) {
	if tenant == "" || organization == "" || subject == "" || receiptID == "" || !returnInventoryActions[inventoryAction] || len(strings.TrimSpace(notes)) < 1 || len(notes) > 1000 {
		return ReturnDisposition{}, ErrInvalid
	}
	return s.repository.DecideReturn(ctx, tenant, organization, subject, receiptID, inventoryAction, notes, s.ids.New(), s.ids.New(), s.ids.New(), s.ids.New(), s.ids.New(), s.ids.New())
}

func (s *Service) AcceptQuote(ctx context.Context, tenant, organization, customer, quote string, version int64, evidence string) (Quote, error) {
	if tenant == "" || organization == "" || customer == "" || quote == "" || version < 1 || !sha256Hex(evidence) {
		return Quote{}, ErrInvalid
	}
	return s.repository.AcceptQuote(ctx, tenant, organization, customer, quote, version, evidence, s.ids.New(), s.ids.New(), s.ids.New(), s.ids.New())
}

func (s *Service) AcceptHandover(ctx context.Context, tenant, organization, customer, handover string, version int64, serialNumber, checklistID string, checklistVersion int64, evidence string) (Handover, error) {
	if tenant == "" || organization == "" || customer == "" || handover == "" || version < 1 || serialNumber == "" || len(serialNumber) > 128 || !code(checklistID) || checklistVersion < 1 || !sha256Hex(evidence) {
		return Handover{}, ErrInvalid
	}
	return s.repository.AcceptHandover(ctx, tenant, organization, customer, handover, version, serialNumber, checklistID, checklistVersion, evidence, s.ids.New())
}

func sha256Hex(value string) bool {
	decoded, err := hex.DecodeString(value)
	return err == nil && len(decoded) == 32 && value == strings.ToLower(value)
}

var appointmentKinds = map[string]bool{"consultation": true, "test-drive": true, "delivery": true, "service": true}
var resourceKinds = map[string]bool{"employee": true, "contractor": true, "service-bay": true, "vehicle": true, "equipment": true}

func code(value string) bool {
	if len(value) < 1 || len(value) > 64 || value[0] < 'a' || value[0] > 'z' {
		return false
	}
	for _, r := range value {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			return false
		}
	}
	return !strings.Contains(value, "--") && !strings.HasSuffix(value, "-")
}

func WrapConflict(operation string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %s", ErrConflict, operation)
}

// QuoteResult only reads the durable result; absence never authorizes a new command.
func (s *Service) QuoteResult(ctx context.Context, tenant, organization, lead, key string) (Quote, error) {
	if tenant == "" || organization == "" || lead == "" || len(key) < 16 || len(key) > 128 {
		return Quote{}, ErrInvalid
	}
	repository, ok := s.repository.(interface {
		QuoteResult(context.Context, string, string, string, string) (Quote, error)
	})
	if !ok {
		return Quote{}, ErrConflict
	}
	return repository.QuoteResult(ctx, tenant, organization, lead, key)
}

var resourceCreationKeyPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{16,128}$`)

func (s *Service) CreateServiceResourceOnce(ctx context.Context, tenant, subject, key, hash string, value ServiceResource) (ServiceResource, bool, error) {
	if subject == "" || !resourceCreationKeyPattern.MatchString(key) || !sha256Hex(hash) {
		return ServiceResource{}, false, ErrInvalid
	}
	value, err := s.prepareServiceResource(tenant, value)
	if err != nil {
		return ServiceResource{}, false, err
	}
	r, ok := s.repository.(interface {
		CreateServiceResourceOnce(context.Context, string, string, string, string, ServiceResource, string) (ServiceResource, bool, error)
	})
	if !ok {
		return ServiceResource{}, false, ErrConflict
	}
	return r.CreateServiceResourceOnce(ctx, tenant, subject, key, hash, value, s.ids.New())
}
func (s *Service) ServiceResourceCreationResult(ctx context.Context, tenant, organization, key string) (ServiceResource, error) {
	if tenant == "" || organization == "" || !resourceCreationKeyPattern.MatchString(key) {
		return ServiceResource{}, ErrInvalid
	}
	r, ok := s.repository.(interface {
		ServiceResourceCreationResult(context.Context, string, string, string) (ServiceResource, error)
	})
	if !ok {
		return ServiceResource{}, ErrConflict
	}
	return r.ServiceResourceCreationResult(ctx, tenant, organization, key)
}

var slotCreationKeyPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{16,128}$`)

func (s *Service) CreateAppointmentSlotOnce(ctx context.Context, tenant, subject, key, hash string, value AppointmentSlot) (AppointmentSlot, bool, error) {
	if subject == "" || !slotCreationKeyPattern.MatchString(key) || !sha256Hex(hash) {
		return AppointmentSlot{}, false, ErrInvalid
	}
	value, err := s.prepareAppointmentSlot(tenant, value)
	if err != nil {
		return AppointmentSlot{}, false, err
	}
	r, ok := s.repository.(interface {
		CreateAppointmentSlotOnce(context.Context, string, string, string, string, AppointmentSlot, string) (AppointmentSlot, bool, error)
	})
	if !ok {
		return AppointmentSlot{}, false, ErrConflict
	}
	return r.CreateAppointmentSlotOnce(ctx, tenant, subject, key, hash, value, s.ids.New())
}
func (s *Service) AppointmentSlotCreationResult(ctx context.Context, tenant, organization, key string) (AppointmentSlot, error) {
	if tenant == "" || organization == "" || !slotCreationKeyPattern.MatchString(key) {
		return AppointmentSlot{}, ErrInvalid
	}
	r, ok := s.repository.(interface {
		AppointmentSlotCreationResult(context.Context, string, string, string) (AppointmentSlot, error)
	})
	if !ok {
		return AppointmentSlot{}, ErrConflict
	}
	return r.AppointmentSlotCreationResult(ctx, tenant, organization, key)
}

func (s *Service) PublishedDeliveryChecklist(ctx context.Context, tenant, organization, id string, version int64) (DeliveryChecklist, error) {
	if tenant == "" || organization == "" || !code(id) || version < 1 {
		return DeliveryChecklist{}, ErrInvalid
	}
	repo, ok := s.repository.(interface {
		PublishedDeliveryChecklist(context.Context, string, string, string, int64) (DeliveryChecklist, error)
	})
	if !ok {
		return DeliveryChecklist{}, ErrConflict
	}
	v, err := repo.PublishedDeliveryChecklist(ctx, tenant, organization, id, version)
	if err != nil {
		return DeliveryChecklist{}, err
	}
	if v.ID != id || v.OrganizationID != organization || v.Version != version || v.State != "published" || len(strings.TrimSpace(v.Title)) < 1 || len(v.Title) > 160 || len(v.Items) < 1 || len(v.Items) > 64 {
		return DeliveryChecklist{}, ErrConflict
	}
	seen := map[string]bool{}
	for index, item := range v.Items {
		if item.Ordinal != index+1 || !code(item.ID) || seen[item.ID] || len(strings.TrimSpace(item.Prompt)) < 1 || len(item.Prompt) > 500 || !checklistResponseTypes[item.ResponseType] {
			return DeliveryChecklist{}, ErrConflict
		}
		seen[item.ID] = true
	}
	return v, nil
}

type ChecklistCompletion struct {
	HandoverID       string              `json:"handover_id"`
	OrganizationID   string              `json:"organization_id"`
	State            string              `json:"state"`
	Version          int64               `json:"version"`
	ChecklistID      string              `json:"checklist_id"`
	ChecklistVersion int64               `json:"checklist_version"`
	CompletedAt      time.Time           `json:"completed_at"`
	ActorSubject     string              `json:"actor_subject"`
	Responses        []ChecklistResponse `json:"responses"`
}

func (s *Service) DeliveryChecklistCompletion(ctx context.Context, tenant, organization, handover string) (ChecklistCompletion, error) {
	if tenant == "" || organization == "" || handover == "" || len(handover) > 128 {
		return ChecklistCompletion{}, ErrInvalid
	}
	repo, ok := s.repository.(interface {
		DeliveryChecklistCompletion(context.Context, string, string, string) (ChecklistCompletion, error)
	})
	if !ok {
		return ChecklistCompletion{}, ErrConflict
	}
	v, err := repo.DeliveryChecklistCompletion(ctx, tenant, organization, handover)
	if err != nil {
		return ChecklistCompletion{}, err
	}
	if v.HandoverID != handover || v.OrganizationID != organization || v.Version < 2 || (v.State != "presented" && v.State != "accepted" && v.State != "rejected") || !code(v.ChecklistID) || v.ChecklistVersion < 1 || v.CompletedAt.IsZero() || strings.TrimSpace(v.ActorSubject) == "" || len(v.ActorSubject) > 255 || len(v.Responses) < 1 || len(v.Responses) > 64 {
		return ChecklistCompletion{}, ErrConflict
	}
	previous := ""
	for _, response := range v.Responses {
		if !code(response.ItemID) || response.ItemID <= previous || len(strings.TrimSpace(response.ResponseText)) < 1 || len(response.ResponseText) > 2048 || (response.EvidenceSHA256 != "" && !sha256Hex(response.EvidenceSHA256)) {
			return ChecklistCompletion{}, ErrConflict
		}
		previous = response.ItemID
	}
	return v, nil
}

func (s *Service) ReturnCaseResult(ctx context.Context, tenant, organization, authorization string) (ReturnCase, error) {
	if tenant == "" || organization == "" || authorization == "" || len(authorization) > 128 {
		return ReturnCase{}, ErrInvalid
	}
	repo, ok := s.repository.(interface {
		ReturnCaseResult(context.Context, string, string, string) (ReturnCase, error)
	})
	if !ok {
		return ReturnCase{}, ErrConflict
	}
	v, err := repo.ReturnCaseResult(ctx, tenant, organization, authorization)
	if err != nil {
		return ReturnCase{}, err
	}
	bounded := func(v string, max int) bool { return len(strings.TrimSpace(v)) > 0 && len(v) <= max }
	if v.AuthorizationID != authorization || v.OrganizationID != organization || !bounded(v.OrderID, 128) || !bounded(v.StockUnitID, 128) || !bounded(v.CustomerSubject, 255) || (v.AuthorizedAction != "return" && v.AuthorizedAction != "exchange") || v.AuthorizedAt.IsZero() {
		return ReturnCase{}, ErrConflict
	}
	if v.Receipt == nil {
		if v.Disposition != nil {
			return ReturnCase{}, ErrConflict
		}
		return v, nil
	}
	receipt := v.Receipt
	if !bounded(receipt.ID, 128) || receipt.AuthorizationID != authorization || receipt.OrganizationID != organization || receipt.OrderID != v.OrderID || receipt.StockUnitID != v.StockUnitID || receipt.CustomerSubject != v.CustomerSubject || (len(receipt.ReceivedSerialNumber) < 1 || len(receipt.ReceivedSerialNumber) > 128) || !returnConditionCodes[receipt.ConditionCode] || !bounded(receipt.Notes, 1000) || !sha256Hex(receipt.EvidenceSHA256) || !bounded(receipt.ReceivedBySubject, 255) || receipt.ReceivedAt.IsZero() {
		return ReturnCase{}, ErrConflict
	}
	if v.Disposition == nil {
		return v, nil
	}
	d := v.Disposition
	remedy := "refund"
	if v.AuthorizedAction == "exchange" {
		remedy = "exchange"
	}
	if !bounded(d.ID, 128) || d.ReceiptID != receipt.ID || !returnInventoryActions[d.InventoryAction] || d.CustomerRemedy != remedy || !bounded(d.Notes, 1000) || !bounded(d.DecidedBySubject, 255) || d.DecidedAt.IsZero() {
		return ReturnCase{}, ErrConflict
	}
	owners := map[string]string{"inventory": "inventory", "accounting": "accounting", remedy: map[string]string{"refund": "payment", "exchange": "fulfillment"}[remedy]}
	if remedy == "refund" {
		owners["fiscal"] = "fiscal"
	}
	if len(d.Effects) != len(owners) {
		return ReturnCase{}, ErrConflict
	}
	seenKinds, seenIDs := map[string]bool{}, map[string]bool{}
	for _, e := range d.Effects {
		owner, ok := owners[e.EffectKind]
		if !ok || owner != e.OwnerContext || seenKinds[e.EffectKind] || seenIDs[e.ID] || !bounded(e.ID, 128) || e.State != "requested" || len(e.IdempotencyKey) < 16 || len(e.IdempotencyKey) > 128 || e.RequestedAt.IsZero() {
			return ReturnCase{}, ErrConflict
		}
		seenKinds[e.EffectKind] = true
		seenIDs[e.ID] = true
	}
	return v, nil
}
