// Package inventorycontrol is a narrow Go/PostgreSQL adaptation of the reservation,
// ATP and transfer invariants documented in docs/inventory/MICROSOFT_BC_DERIVATION.md.
// It is not verbatim Microsoft code or a Business Central replacement.
package inventorycontrol

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrConflict = errors.New("inventory control conflict")

type ATP struct {
	OrganizationID     string    `json:"organization_id"`
	VariantID          string    `json:"variant_id"`
	Horizon            time.Time `json:"horizon"`
	AvailableInventory int64     `json:"available_inventory"`
	ScheduledReceipt   int64     `json:"scheduled_receipt"`
	GrossRequirement   int64     `json:"gross_requirement"`
	AvailableToPromise int64     `json:"available_to_promise"`
}

type Reservation struct {
	ID                     string     `json:"id"`
	OrganizationID         string     `json:"organization_id"`
	StockUnitID            string     `json:"stock_unit_id"`
	VariantID              string     `json:"variant_id"`
	DemandKind             string     `json:"demand_kind"`
	DemandID               string     `json:"demand_id"`
	DemandLineID           string     `json:"demand_line_id"`
	Status                 string     `json:"status"`
	CancellationDisallowed bool       `json:"cancellation_disallowed"`
	ExpiresAt              *time.Time `json:"expires_at,omitempty"`
	Version                int64      `json:"version"`
}

type Transfer struct {
	ID                 string    `json:"id"`
	FromOrganizationID string    `json:"from_organization_id"`
	ToOrganizationID   string    `json:"to_organization_id"`
	StockUnitIDs       []string  `json:"stock_unit_ids"`
	ExpectedReceiptAt  time.Time `json:"expected_receipt_at"`
	State              string    `json:"state"`
	Version            int64     `json:"version"`
}

type Repository interface {
	AvailableToPromise(context.Context, string, string, string, time.Time) (ATP, error)
	Reserve(context.Context, string, string, string, Reservation, int64) (Reservation, error)
	Release(context.Context, string, string, string, int64, string) error
	CreateTransfer(context.Context, string, string, Transfer, map[string]int64, string) (Transfer, error)
	TransitionTransfer(context.Context, string, string, string, string, int64, string, string) error
}

type IDGenerator interface{ New() string }

type Service struct {
	repository Repository
	ids        IDGenerator
}

func NewService(repository Repository, ids IDGenerator) *Service {
	return &Service{repository: repository, ids: ids}
}

func (s *Service) AvailableToPromise(ctx context.Context, tenant, organization, variant string, horizon time.Time) (ATP, error) {
	if tenant == "" || organization == "" || variant == "" || horizon.IsZero() {
		return ATP{}, fmt.Errorf("invalid ATP query")
	}
	return s.repository.AvailableToPromise(ctx, tenant, organization, variant, horizon.UTC())
}

func (s *Service) Reserve(ctx context.Context, tenant string, value Reservation, stockVersion int64) (Reservation, error) {
	allowed := map[string]bool{"service": true, "manual": true}
	if tenant == "" || value.OrganizationID == "" || value.StockUnitID == "" || value.VariantID == "" || !allowed[value.DemandKind] || value.DemandID == "" || stockVersion < 1 {
		return Reservation{}, fmt.Errorf("invalid reservation")
	}
	value.ID, value.Status, value.Version = s.ids.New(), "reservation", 1
	return s.repository.Reserve(ctx, tenant, s.ids.New(), value.ID, value, stockVersion)
}

func (s *Service) Release(ctx context.Context, tenant, organization, reservationID string, version int64) error {
	if tenant == "" || organization == "" || reservationID == "" || version < 1 {
		return fmt.Errorf("invalid reservation release")
	}
	return s.repository.Release(ctx, tenant, organization, reservationID, version, s.ids.New())
}

func (s *Service) CreateTransfer(ctx context.Context, tenant string, value Transfer, stockVersions map[string]int64) (Transfer, error) {
	if tenant == "" || value.FromOrganizationID == "" || value.ToOrganizationID == "" || value.FromOrganizationID == value.ToOrganizationID || value.ExpectedReceiptAt.IsZero() || len(value.StockUnitIDs) == 0 || len(value.StockUnitIDs) != len(stockVersions) {
		return Transfer{}, fmt.Errorf("invalid transfer")
	}
	seen := map[string]bool{}
	for _, stockID := range value.StockUnitIDs {
		if stockID == "" || seen[stockID] || stockVersions[stockID] < 1 {
			return Transfer{}, fmt.Errorf("invalid transfer units")
		}
		seen[stockID] = true
	}
	value.ID, value.State, value.Version = s.ids.New(), "draft", 1
	value.ExpectedReceiptAt = value.ExpectedReceiptAt.UTC()
	return s.repository.CreateTransfer(ctx, tenant, value.ID, value, stockVersions, s.ids.New())
}

var transferTransitions = map[string]map[string]bool{
	"draft":      {"released": true, "cancelled": true},
	"released":   {"in-transit": true, "cancelled": true},
	"in-transit": {"received": true},
}

func (s *Service) TransitionTransfer(ctx context.Context, tenant, authorizedOrganization, transferID, current, target string, version int64) error {
	if tenant == "" || authorizedOrganization == "" || transferID == "" || !transferTransitions[current][target] || version < 1 {
		return fmt.Errorf("%w: invalid transfer transition", ErrConflict)
	}
	return s.repository.TransitionTransfer(ctx, tenant, authorizedOrganization, transferID, current, version, target, s.ids.New())
}
