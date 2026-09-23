package postgres

// AUTHORED exported composition point for the existing transactional fence.
// A missing domain admission callback is never silently accepted.
import (
	"context"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"github.com/jackc/pgx/v5"
)

func (s *OutboundDeliveryStore) ClaimWithDomainAdmission(ctx context.Context, m channels.Message, h string, admit func(context.Context, pgx.Tx) error) (outbounddelivery.Claim, error) {
	if admit == nil {
		return outbounddelivery.Claim{}, outbounddelivery.ErrInvalid
	}
	return s.claimWithAdmission(ctx, m, h, admit)
}
