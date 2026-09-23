package electromobility

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

var ErrNotFound = errors.New("resource not found")
var ErrConflict = errors.New("electromobility conflict")
var codePattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`)

type Model struct {
	ID            string          `json:"id"`
	Code          string          `json:"code"`
	DisplayName   string          `json:"displayName"`
	VehicleClass  string          `json:"vehicleClass"`
	Specification json.RawMessage `json:"specification"`
}
type Lead struct {
	ID, TenantID, OrganizationID, ModelID, SourceCode string
	Contact                                           json.RawMessage
	ConsentID, ConsentEvidenceHash                    string
}
type Repository interface {
	ListPublicModels(context.Context, string) ([]Model, error)
	CreateModel(context.Context, string, string, Model) error
	ResolvePublicOrganization(context.Context, string, string) (string, string, error)
	CreateLead(context.Context, Lead, string, string) (string, bool, error)
}
type IDGenerator interface{ New() string }
type Service struct {
	repository Repository
	ids        IDGenerator
}

func NewService(repository Repository, ids IDGenerator) *Service {
	return &Service{repository: repository, ids: ids}
}

func (s *Service) ListPublicModels(ctx context.Context, tenantCode string) ([]Model, error) {
	if !codePattern.MatchString(tenantCode) {
		return nil, fmt.Errorf("invalid tenant code")
	}
	return s.repository.ListPublicModels(ctx, tenantCode)
}
func (s *Service) CreateModel(ctx context.Context, tenantID string, input Model) (Model, error) {
	classes := map[string]bool{"motorcycle": true, "bicycle": true, "scooter": true, "utility": true, "other": true}
	if tenantID == "" || !codePattern.MatchString(input.Code) || len(input.DisplayName) < 2 || !classes[input.VehicleClass] || !validObject(input.Specification) {
		return Model{}, fmt.Errorf("invalid model")
	}
	input.ID = s.ids.New()
	eventID := s.ids.New()
	if err := s.repository.CreateModel(ctx, tenantID, eventID, input); err != nil {
		return Model{}, err
	}
	return input, nil
}
func (s *Service) CaptureLead(ctx context.Context, tenantCode, organizationCode, modelID, source string, contact json.RawMessage, consentHash, idempotencyKey string) (Lead, bool, error) {
	if !codePattern.MatchString(tenantCode) || !codePattern.MatchString(organizationCode) || !codePattern.MatchString(source) || !validObject(contact) || !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(consentHash) || len(idempotencyKey) < 16 || len(idempotencyKey) > 128 {
		return Lead{}, false, fmt.Errorf("invalid lead")
	}
	tenantID, organizationID, err := s.repository.ResolvePublicOrganization(ctx, tenantCode, organizationCode)
	if err != nil {
		return Lead{}, false, err
	}
	lead := Lead{ID: s.ids.New(), TenantID: tenantID, OrganizationID: organizationID, ModelID: modelID, SourceCode: source, Contact: contact, ConsentID: s.ids.New(), ConsentEvidenceHash: consentHash}
	resourceID, replayed, err := s.repository.CreateLead(ctx, lead, s.ids.New(), idempotencyKey)
	if err != nil {
		return Lead{}, false, err
	}
	lead.ID = resourceID
	return lead, replayed, nil
}
func validObject(value json.RawMessage) bool {
	var object map[string]any
	return len(value) > 1 && json.Unmarshal(value, &object) == nil && object != nil
}
