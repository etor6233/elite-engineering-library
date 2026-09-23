package fiscal

import (
	"context"
	"errors"
	"time"
)

var ErrNoParameterWork = errors.New("no parameter refresh work")

type ParameterScheduleRepository interface {
	ClaimParameterSchedule(context.Context, string, time.Duration) (ParameterSchedule, error)
	CompleteParameterSchedule(context.Context, ParameterSchedule, string, string, string) error
	DeferParameterSchedule(context.Context, ParameterSchedule, string, string, time.Duration) error
}

type ParameterRefreshProcessor struct {
	repository ParameterScheduleRepository
	registry   *ParameterRegistry
	workerID   string
	lease      time.Duration
	retry      time.Duration
	ids        IDGenerator
}

func NewParameterRefreshProcessor(repository ParameterScheduleRepository, registry *ParameterRegistry, ids IDGenerator, workerID string, lease, retry time.Duration) (*ParameterRefreshProcessor, error) {
	if repository == nil || registry == nil || ids == nil || workerID == "" || lease < 30*time.Second || lease > 10*time.Minute || retry < time.Minute || retry > time.Hour {
		return nil, ErrInvalid
	}
	return &ParameterRefreshProcessor{repository: repository, registry: registry, workerID: workerID, lease: lease, retry: retry, ids: ids}, nil
}

func (p *ParameterRefreshProcessor) ProcessOne(ctx context.Context) error {
	schedule, err := p.repository.ClaimParameterSchedule(ctx, p.workerID, p.lease)
	if err != nil {
		return err
	}
	snapshot, _, err := p.registry.Refresh(ctx, schedule.TenantID, schedule.OrganizationID, schedule.TaxpayerCUIT, schedule.Kind, schedule.VoucherClass)
	if err != nil {
		_ = p.repository.DeferParameterSchedule(ctx, schedule, p.workerID, "PARAMETER_PROVIDER_FAILURE", p.retry)
		return err
	}
	return p.repository.CompleteParameterSchedule(ctx, schedule, p.workerID, snapshot.ID, p.ids.New())
}
