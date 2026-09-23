package fiscal

import (
	"context"
	"time"
)

var parameterKinds = map[string]bool{
	"voucher_type": true, "concept": true, "document_type": true, "vat_rate": true,
	"other_tax": true, "point_of_sale": true, "recipient_vat_condition": true,
}

type ParameterItem struct {
	Code           string  `json:"code"`
	Description    *string `json:"description"`
	ValidFrom      *string `json:"valid_from"`
	ValidUntil     *string `json:"valid_until"`
	VoucherClass   *string `json:"voucher_class"`
	EmissionType   *string `json:"emission_type"`
	Blocked        *string `json:"blocked"`
	DeregisteredOn *string `json:"deregistered_on"`
}

type ParameterSnapshot struct {
	TenantID       string          `json:"-"`
	ID             string          `json:"id"`
	OrganizationID string          `json:"organization_id"`
	TaxpayerCUIT   string          `json:"taxpayer_cuit"`
	Kind           string          `json:"kind"`
	VoucherClass   *string         `json:"voucher_class"`
	Items          []ParameterItem `json:"items"`
	ResponseHash   string          `json:"response_hash"`
	ProviderCodes  []string        `json:"provider_codes"`
	FetchedAt      time.Time       `json:"fetched_at"`
}

type ParameterSchedule struct {
	TenantID        string    `json:"-"`
	ID              string    `json:"id"`
	OrganizationID  string    `json:"organization_id"`
	TaxpayerCUIT    string    `json:"taxpayer_cuit"`
	Kind            string    `json:"kind"`
	VoucherClass    *string   `json:"voucher_class"`
	IntervalSeconds int64     `json:"interval_seconds"`
	NextRunAt       time.Time `json:"next_run_at"`
	Active          bool      `json:"active"`
	Version         int64     `json:"version"`
	LeaseOwner      string    `json:"-"`
}

type ParameterProvider interface {
	FetchParameters(context.Context, string, string, *string) (ParameterSnapshot, error)
}

type ParameterRepository interface {
	StoreParameterSnapshot(context.Context, ParameterSnapshot, string) (ParameterSnapshot, bool, error)
	GetParameterSnapshot(context.Context, string, string, string) (ParameterSnapshot, error)
	DecideParameterSnapshot(context.Context, string, string, string, string, bool, string, string) error
	ConfigureParameterSchedule(context.Context, ParameterSchedule, string, string) (ParameterSchedule, error)
}

type ParameterRegistry struct {
	repository ParameterRepository
	provider   ParameterProvider
	ids        IDGenerator
	now        func() time.Time
}

func NewParameterRegistry(repository ParameterRepository, provider ParameterProvider, ids IDGenerator, now func() time.Time) *ParameterRegistry {
	return &ParameterRegistry{repository: repository, provider: provider, ids: ids, now: now}
}

func (r *ParameterRegistry) Refresh(ctx context.Context, tenant, organization, taxpayerCUIT, kind string, voucherClass *string) (ParameterSnapshot, bool, error) {
	if tenant == "" || organization == "" || !ValidCUIT(taxpayerCUIT) || !validParameterScope(kind, voucherClass) {
		return ParameterSnapshot{}, false, ErrInvalid
	}
	value, err := r.provider.FetchParameters(ctx, taxpayerCUIT, kind, voucherClass)
	if err != nil {
		return ParameterSnapshot{}, false, err
	}
	if value.TaxpayerCUIT != taxpayerCUIT || value.Kind != kind || !sameOptional(value.VoucherClass, voucherClass) || !hashPattern.MatchString(value.ResponseHash) || len(value.Items) == 0 || !validParameterItems(value.Items, voucherClass) {
		return ParameterSnapshot{}, false, ErrInvalid
	}
	value.TenantID, value.OrganizationID, value.ID, value.FetchedAt = tenant, organization, r.ids.New(), r.now().UTC()
	return r.repository.StoreParameterSnapshot(ctx, value, r.ids.New())
}

func (r *ParameterRegistry) Snapshot(ctx context.Context, tenant, organization, snapshotID string) (ParameterSnapshot, error) {
	if tenant == "" || organization == "" || snapshotID == "" {
		return ParameterSnapshot{}, ErrInvalid
	}
	return r.repository.GetParameterSnapshot(ctx, tenant, organization, snapshotID)
}

func (r *ParameterRegistry) Decide(ctx context.Context, tenant, organization, snapshotID, subject string, approved bool, reason string) error {
	if tenant == "" || organization == "" || snapshotID == "" || subject == "" || len(reason) < 3 || len(reason) > 500 {
		return ErrInvalid
	}
	return r.repository.DecideParameterSnapshot(ctx, tenant, organization, snapshotID, subject, approved, reason, r.ids.New())
}

func (r *ParameterRegistry) ConfigureSchedule(ctx context.Context, tenant, subject string, value ParameterSchedule) (ParameterSchedule, error) {
	if tenant == "" || subject == "" || value.OrganizationID == "" || !ValidCUIT(value.TaxpayerCUIT) || !validParameterScope(value.Kind, value.VoucherClass) || value.IntervalSeconds < 900 || value.IntervalSeconds > 2592000 || value.NextRunAt.IsZero() {
		return value, ErrInvalid
	}
	value.TenantID, value.ID, value.Active, value.Version = tenant, r.ids.New(), true, 1
	value.NextRunAt = value.NextRunAt.UTC()
	return r.repository.ConfigureParameterSchedule(ctx, value, subject, r.ids.New())
}

func validParameterScope(kind string, voucherClass *string) bool {
	if !parameterKinds[kind] {
		return false
	}
	if kind != "recipient_vat_condition" {
		return voucherClass == nil
	}
	return voucherClass != nil && (*voucherClass == "A" || *voucherClass == "B" || *voucherClass == "C" || *voucherClass == "M")
}

func validParameterItems(items []ParameterItem, requestedClass *string) bool {
	seen := map[string]bool{}
	for _, item := range items {
		if item.Code == "" || len(item.Code) > 16 || seen[item.Code] || !validOptionalDate(item.ValidFrom) || !validOptionalDate(item.ValidUntil) || !validOptionalDate(item.DeregisteredOn) || (item.ValidFrom != nil && item.ValidUntil != nil && *item.ValidUntil < *item.ValidFrom) {
			return false
		}
		seen[item.Code] = true
		if requestedClass != nil && (item.VoucherClass == nil || *item.VoucherClass != *requestedClass) {
			return false
		}
	}
	return true
}

func validOptionalDate(value *string) bool {
	if value == nil {
		return true
	}
	_, err := time.Parse("2006-01-02", *value)
	return err == nil
}

func sameOptional(left, right *string) bool {
	return left == nil && right == nil || left != nil && right != nil && *left == *right
}
