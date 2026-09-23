package enterprisequery

// AUTHORED read port extension over the existing factory projection.
import (
	"context"
	"errors"
)

var ErrFactoryUnitNotFound = errors.New("factory unit not found")

type FactoryUnitReader interface {
	FactoryUnit(context.Context, string, string, string) (FactoryUnit, error)
}

func (s *Service) FactoryUnit(ctx context.Context, tenant, organization, id string) (FactoryUnit, error) {
	if validate(tenant, organization, 1) != nil || id == "" || len(id) > 200 {
		return FactoryUnit{}, ErrInvalid
	}
	reader, ok := s.repository.(FactoryUnitReader)
	if !ok {
		return FactoryUnit{}, ErrInvalid
	}
	return reader.FactoryUnit(ctx, tenant, organization, id)
}
