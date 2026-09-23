package postgres

// AUTHORED composition: use the existing outbound fence transaction as the
// command acceptance point. No network call occurs while database locks are held.
import (
	"context"
	"encoding/json"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/paymentbridge"
	"github.com/jackc/pgx/v5"
)

type PaymentDispatchFence struct {
	*OutboundDeliveryStore
	Scope paymentbridge.Scope
}

func (f *PaymentDispatchFence) Claim(ctx context.Context, m channels.Message, hash string) (outbounddelivery.Claim, error) {
	if f == nil || f.OutboundDeliveryStore == nil || f.Scope.Validate() != nil || m.TenantID != f.Scope.TenantID || m.ChannelCode != "payment_"+f.Scope.ProviderCode {
		return outbounddelivery.Claim{}, paymentbridge.ErrCheckoutConflict
	}
	computed, err := outbounddelivery.MessageSHA256(m)
	if err != nil || computed != hash {
		return outbounddelivery.Claim{}, paymentbridge.ErrCheckoutConflict
	}
	var requested paymentbridge.Request
	if json.Unmarshal([]byte(m.Text), &requested) != nil || requested.PaymentAttemptID != m.DeliveryKey || requested.OrderID != m.ThreadID || requested.CustomerSubject != m.ExternalID {
		return outbounddelivery.Claim{}, paymentbridge.ErrCheckoutConflict
	}
	return f.claimWithAdmission(ctx, m, hash, func(ctx context.Context, tx pgx.Tx) error {
		if err := checkoutConnection(ctx, tx, f.Scope); err != nil {
			return err
		}
		actual, err := lockCheckoutRequest(ctx, tx, f.Scope, requested.PaymentAttemptID)
		if err != nil {
			return err
		}
		if actual != requested {
			return paymentbridge.ErrCheckoutConflict
		}
		// Bind admission and the created→pending command acceptance in the same
		// transaction as the one allowed send. Later cancellation never creates a
		// second send; subsequent financial observations preserve terminal states.
		row, err := tx.Exec(ctx, `with accepted as (
   update payment.payment_attempt p set state='pending',version=p.version+1,updated_at=clock_timestamp()
   from sales.customer_order o,payment.provider_checkout c
   where p.tenant_id=$1 and p.payment_attempt_id=$2 and p.state='created'
     and o.tenant_id=p.tenant_id and o.order_id=p.order_id and o.organization_id=$3 and o.state in ('placed','confirmed','allocated')
     and c.tenant_id=p.tenant_id and c.payment_attempt_id=p.payment_attempt_id and c.connection_id=$4 and c.account_ref=$5 and c.request_sha256_hex=$6 and c.live_mode=$7 and c.session_id is null and c.expires_at>clock_timestamp()+interval '31 minutes'
   returning p.version)
   insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)
   select $1,gen_random_uuid(),'payment',$2,version,'payment.checkout-requested',1,clock_timestamp(),jsonb_build_object('connection_id',$4::text,'request_sha256',$6::text) from accepted`, f.Scope.TenantID, requested.PaymentAttemptID, f.Scope.OrganizationID, f.Scope.ConnectionID, f.Scope.AccountRef, hash, f.Scope.LiveMode)
		if err != nil {
			return err
		}
		if row.RowsAffected() != 1 {
			return paymentbridge.ErrCheckoutConflict
		}
		return nil
	})
}
