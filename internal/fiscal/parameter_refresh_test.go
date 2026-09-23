package fiscal

import (
	"context"
	"errors"
	"testing"
	"time"
)

type scheduleRepository struct {
	schedule  ParameterSchedule
	completed bool
	deferred  bool
}

func (r *scheduleRepository) ClaimParameterSchedule(context.Context, string, time.Duration) (ParameterSchedule, error) {
	return r.schedule, nil
}
func (r *scheduleRepository) CompleteParameterSchedule(_ context.Context, _ ParameterSchedule, _, snapshotID, _ string) error {
	r.completed = snapshotID != ""
	return nil
}
func (r *scheduleRepository) DeferParameterSchedule(context.Context, ParameterSchedule, string, string, time.Duration) error {
	r.deferred = true
	return nil
}

func TestParameterRefreshProcessorStoresButNeverApproves(t *testing.T) {
	repository := &parameterRepository{}
	provider := &parameterProvider{value: ParameterSnapshot{TaxpayerCUIT: "33693450239", Kind: "vat_rate", Items: []ParameterItem{{Code: "5"}}, ResponseHash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}}
	registry := NewParameterRegistry(repository, provider, sequenceIDs{}, time.Now)
	schedules := &scheduleRepository{schedule: ParameterSchedule{TenantID: "tenant", OrganizationID: "organization", TaxpayerCUIT: "33693450239", Kind: "vat_rate"}}
	processor, err := NewParameterRefreshProcessor(schedules, registry, sequenceIDs{}, "worker", time.Minute, 5*time.Minute)
	if err != nil || processor.ProcessOne(context.Background()) != nil || !schedules.completed || schedules.deferred || repository.decided {
		t.Fatalf("processor err=%v completed=%v deferred=%v approved=%v", err, schedules.completed, schedules.deferred, repository.decided)
	}
	provider.err = errors.New("provider unavailable")
	provider.value = ParameterSnapshot{}
	schedules.completed, schedules.deferred = false, false
	if err = processor.ProcessOne(context.Background()); err == nil || !schedules.deferred || schedules.completed {
		t.Fatalf("failure err=%v completed=%v deferred=%v", err, schedules.completed, schedules.deferred)
	}
}
