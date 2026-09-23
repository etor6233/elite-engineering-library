package fiscal

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestParameterRegistryRefreshAndApprovalAreSeparated(t *testing.T) {
	description := "21%"
	provider := &parameterProvider{value: ParameterSnapshot{TaxpayerCUIT: "33693450239", Kind: "vat_rate", Items: []ParameterItem{{Code: "5", Description: &description}}, ResponseHash: string(make([]byte, 64))}}
	provider.value.ResponseHash = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	repository := &parameterRepository{}
	registry := NewParameterRegistry(repository, provider, sequenceIDs{}, func() time.Time { return time.Date(2026, 8, 30, 12, 0, 0, 0, time.UTC) })
	value, replay, err := registry.Refresh(context.Background(), "tenant", "organization", "33693450239", "vat_rate", nil)
	if err != nil || replay || value.ID == "" || repository.stored == nil || repository.decided {
		t.Fatalf("refresh=%+v replay=%v err=%v", value, replay, err)
	}
	if err = registry.Decide(context.Background(), "tenant", "organization", value.ID, "tax-owner", true, "validated against current ARCA configuration"); err != nil || !repository.decided {
		t.Fatalf("approval was not explicit: %v", err)
	}
}

func TestParameterRegistryRejectsInventedScopeAndDivergentResponse(t *testing.T) {
	registry := NewParameterRegistry(&parameterRepository{}, &parameterProvider{err: errors.New("must not call")}, sequenceIDs{}, time.Now)
	if _, _, err := registry.Refresh(context.Background(), "tenant", "organization", "33693450239", "invented", nil); !errors.Is(err, ErrInvalid) {
		t.Fatalf("invented kind accepted: %v", err)
	}
	class := "A"
	provider := &parameterProvider{value: ParameterSnapshot{TaxpayerCUIT: "33693450239", Kind: "recipient_vat_condition", VoucherClass: &class, Items: []ParameterItem{{Code: "1", VoucherClass: ptr("B")}}, ResponseHash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}}
	registry = NewParameterRegistry(&parameterRepository{}, provider, sequenceIDs{}, time.Now)
	if _, _, err := registry.Refresh(context.Background(), "tenant", "organization", "33693450239", "recipient_vat_condition", &class); !errors.Is(err, ErrInvalid) {
		t.Fatalf("divergent class accepted: %v", err)
	}
}

type parameterProvider struct {
	value ParameterSnapshot
	err   error
}

func (p *parameterProvider) FetchParameters(context.Context, string, string, *string) (ParameterSnapshot, error) {
	return p.value, p.err
}

type parameterRepository struct {
	stored  *ParameterSnapshot
	decided bool
}

func (r *parameterRepository) StoreParameterSnapshot(_ context.Context, value ParameterSnapshot, _ string) (ParameterSnapshot, bool, error) {
	r.stored = &value
	return value, false, nil
}
func (r *parameterRepository) GetParameterSnapshot(context.Context, string, string, string) (ParameterSnapshot, error) {
	return ParameterSnapshot{}, nil
}
func (r *parameterRepository) DecideParameterSnapshot(context.Context, string, string, string, string, bool, string, string) error {
	r.decided = true
	return nil
}
func (r *parameterRepository) ConfigureParameterSchedule(_ context.Context, value ParameterSchedule, _, _ string) (ParameterSchedule, error) {
	return value, nil
}

type sequenceIDs struct{}

func (sequenceIDs) New() string { return "018f4d4a-7b36-7a21-8d10-2f4c54c29f99" }
func ptr(value string) *string  { return &value }
