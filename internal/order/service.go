package order

import (
	"context"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("order not found")

type Repository interface {
	Create(context.Context, Order, string, string) (Order, bool, error)
	Get(context.Context, string, string, string) (Order, error)
	SaveTransition(context.Context, Order, int64) error
}

type IDGenerator interface{ New() string }

type Service struct {
	repository Repository
	ids        IDGenerator
}

func NewService(repository Repository, ids IDGenerator) *Service {
	return &Service{repository: repository, ids: ids}
}

type Create struct {
	TenantID          string
	OrganizationID    string
	CustomerPrincipal string
	Currency          string
	TotalMinorUnits   int64
	IdempotencyKey    string
	RequestHash       string
}

func (s *Service) Create(ctx context.Context, command Create) (Order, bool, error) {
	if len(command.IdempotencyKey) < 16 || len(command.IdempotencyKey) > 128 || len(command.RequestHash) != 64 {
		return Order{}, false, fmt.Errorf("%w: invalid idempotency metadata", ErrConflict)
	}
	entity, err := New(s.ids.New(), command.TenantID, command.OrganizationID, command.CustomerPrincipal, command.Currency, command.TotalMinorUnits)
	if err != nil {
		return Order{}, false, err
	}
	return s.repository.Create(ctx, entity, command.IdempotencyKey, command.RequestHash)
}

func (s *Service) Transition(ctx context.Context, tenantID, organizationID, id string, target State, permissions map[string]struct{}) (Order, error) {
	current, err := s.repository.Get(ctx, tenantID, organizationID, id)
	if err != nil {
		return Order{}, err
	}
	next, err := current.Transition(target, permissions)
	if err != nil {
		return Order{}, err
	}
	if err := s.repository.SaveTransition(ctx, next, current.Version); err != nil {
		return Order{}, err
	}
	return next, nil
}
