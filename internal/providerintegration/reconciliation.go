package providerintegration

import (
	"context"
	"encoding/json"
	"time"
)

type Mapping struct {
	TenantID, ProviderCode, ResourceType string
	InternalID, ExternalID, VersionToken string
}

type Reconciliation struct {
	TenantID, ID, ProviderCode, ResourceType string
	InternalID, ExternalID, State            string
	Evidence                                 json.RawMessage
	ObservedAt                               time.Time
}

type ReconciliationRepository interface {
	UpsertMapping(context.Context, Mapping) error
	RecordReconciliation(context.Context, Reconciliation) error
	ResolveReconciliation(context.Context, string, string, string) error
}

type ReconciliationService struct{ repository ReconciliationRepository }

func NewReconciliationService(repository ReconciliationRepository) *ReconciliationService {
	return &ReconciliationService{repository: repository}
}

func (s *ReconciliationService) UpsertMapping(ctx context.Context, value Mapping) error {
	if s == nil || s.repository == nil || !bounded(value.TenantID, 200) || !bounded(value.ProviderCode, 64) || !bounded(value.ResourceType, 100) || !bounded(value.InternalID, 200) || !bounded(value.ExternalID, 200) || len(value.VersionToken) > 500 {
		return ErrInvalid
	}
	return s.repository.UpsertMapping(ctx, value)
}

func (s *ReconciliationService) Record(ctx context.Context, value Reconciliation) error {
	allowedState := map[string]bool{"matched": true, "missing-internal": true, "missing-external": true, "amount-mismatch": true, "state-mismatch": true}
	if s == nil || s.repository == nil || !bounded(value.TenantID, 200) || !bounded(value.ID, 200) || !bounded(value.ProviderCode, 64) || !bounded(value.ResourceType, 100) || (!bounded(value.InternalID, 200) && !bounded(value.ExternalID, 200)) || !allowedState[value.State] || !jsonObject(value.Evidence) || value.ObservedAt.IsZero() || value.ObservedAt.After(time.Now().Add(24*time.Hour)) {
		return ErrInvalid
	}
	return s.repository.RecordReconciliation(ctx, value)
}

func (s *ReconciliationService) Resolve(ctx context.Context, tenantID, providerCode, reconciliationID string) error {
	if s == nil || s.repository == nil || !bounded(tenantID, 200) || !bounded(providerCode, 64) || !bounded(reconciliationID, 200) {
		return ErrInvalid
	}
	return s.repository.ResolveReconciliation(ctx, tenantID, providerCode, reconciliationID)
}
