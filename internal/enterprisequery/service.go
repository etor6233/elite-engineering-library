package enterprisequery

import (
	"context"
	"errors"
	"fmt"
)

var ErrInvalid = errors.New("invalid query")

type Overview struct {
	OrganizationID  string `json:"organization_id"`
	Orders          int64  `json:"orders"`
	OpenLeads       int64  `json:"open_leads"`
	StockAvailable  int64  `json:"stock_available"`
	OpenCases       int64  `json:"open_cases"`
	ActiveShipments int64  `json:"active_shipments"`
}

type Order struct {
	ID              string `json:"id"`
	OrganizationID  string `json:"organization_id"`
	CustomerSubject string `json:"customer_subject"`
	State           string `json:"state"`
	Currency        string `json:"currency"`
	TotalMinorUnits int64  `json:"total_minor_units"`
	Version         int64  `json:"version"`
}

type FactoryUnit struct {
	ID              string `json:"id"`
	OrganizationID  string `json:"organization_id"`
	PurchaseOrderID string `json:"purchase_order_id"`
	VariantID       string `json:"variant_id"`
	SerialNumber    string `json:"serial_number"`
	State           string `json:"state"`
}

type ServiceCase struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	StockUnitID    string `json:"stock_unit_id"`
	State          string `json:"state"`
	Severity       string `json:"severity"`
	Description    string `json:"description"`
	Version        int64  `json:"version"`
}

type Page[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
}

type Repository interface {
	Overview(context.Context, string, string) (Overview, error)
	Orders(context.Context, string, string, string, int, string) (Page[Order], error)
	FactoryUnits(context.Context, string, string, int, string) (Page[FactoryUnit], error)
	ServiceCases(context.Context, string, string, string, int, string) (Page[ServiceCase], error)
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func validate(tenant, organization string, limit int) error {
	if tenant == "" || organization == "" || limit < 1 || limit > 100 {
		return fmt.Errorf("%w: scope or limit", ErrInvalid)
	}
	return nil
}

func (s *Service) Overview(ctx context.Context, tenant, organization string) (Overview, error) {
	if err := validate(tenant, organization, 1); err != nil {
		return Overview{}, err
	}
	return s.repository.Overview(ctx, tenant, organization)
}

func (s *Service) Orders(ctx context.Context, tenant, organization, customer string, limit int, after string) (Page[Order], error) {
	if err := validate(tenant, organization, limit); err != nil {
		return Page[Order]{}, err
	}
	return s.repository.Orders(ctx, tenant, organization, customer, limit, after)
}

func (s *Service) FactoryUnits(ctx context.Context, tenant, organization string, limit int, after string) (Page[FactoryUnit], error) {
	if err := validate(tenant, organization, limit); err != nil {
		return Page[FactoryUnit]{}, err
	}
	return s.repository.FactoryUnits(ctx, tenant, organization, limit, after)
}

func (s *Service) ServiceCases(ctx context.Context, tenant, organization, customer string, limit int, after string) (Page[ServiceCase], error) {
	if err := validate(tenant, organization, limit); err != nil {
		return Page[ServiceCase]{}, err
	}
	return s.repository.ServiceCases(ctx, tenant, organization, customer, limit, after)
}
