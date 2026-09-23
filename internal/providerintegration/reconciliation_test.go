package providerintegration

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

type reconciliationRepository struct{ mappings, records, resolutions int }

func (r *reconciliationRepository) UpsertMapping(context.Context, Mapping) error {
	r.mappings++
	return nil
}
func (r *reconciliationRepository) RecordReconciliation(context.Context, Reconciliation) error {
	r.records++
	return nil
}
func (r *reconciliationRepository) ResolveReconciliation(context.Context, string, string, string) error {
	r.resolutions++
	return nil
}

func TestReconciliationContracts(t *testing.T) {
	repository := &reconciliationRepository{}
	service := NewReconciliationService(repository)
	if err := service.UpsertMapping(context.Background(), Mapping{TenantID: "tenant", ProviderCode: "provider", ResourceType: "offer", InternalID: "model-1", ExternalID: "external-1"}); err != nil {
		t.Fatal(err)
	}
	if err := service.Record(context.Background(), Reconciliation{TenantID: "tenant", ID: "recon-1", ProviderCode: "provider", ResourceType: "offer", InternalID: "model-1", State: "state-mismatch", Evidence: json.RawMessage(`{"internal":"active","external":"paused"}`), ObservedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if err := service.Resolve(context.Background(), "tenant", "provider", "recon-1"); err != nil {
		t.Fatal(err)
	}
	if repository.mappings != 1 || repository.records != 1 || repository.resolutions != 1 {
		t.Fatalf("repository=%+v", repository)
	}
	if err := service.Record(context.Background(), Reconciliation{TenantID: "tenant", ID: "bad", ProviderCode: "provider", ResourceType: "offer", State: "invented", Evidence: json.RawMessage(`{}`), ObservedAt: time.Now()}); err == nil {
		t.Fatal("invalid reconciliation accepted")
	}
}
