package postgres

// AUTHORED composition around the existing callback/checkout/fence owners. The
// financial observation and preference identity remain separate receipts.
import (
	"context"
	"errors"

	"elite.local/enterprise/internal/paymentbridge"
)

var ErrCheckoutRecovery = errors.New("original checkout identity unresolved; no resend permitted")

type PaymentCallbackRecoveryProcessor struct{ Base *PaymentCallbackProcessor }

func (p PaymentCallbackRecoveryProcessor) Handle(ctx context.Context, j Job) error {
	if p.Base == nil {
		return ErrPaymentCallback
	}
	// The signed callback's payment resource is reconciled through the unchanged
	// processor first. Missing preference identity never fabricates payment facts.
	if err := p.Base.Handle(ctx, j); err != nil {
		return err
	}
	base := p.Base
	c := base.worker.Scope
	if c.ProviderCode != "mercadopago" {
		return nil
	}
	_, notice, err := base.notice(ctx, j)
	if err != nil {
		return err
	}
	var attempt, session, fence string
	err = base.pool.QueryRow(ctx, `select b.payment_attempt_id,coalesce(b.session_id,''),f.state from payment.provider_checkout b
 join payment.provider_observation o using(tenant_id,payment_attempt_id)
 join communication.outbound_delivery f on f.tenant_id=b.tenant_id and f.delivery_key=b.payment_attempt_id and f.channel_code='payment_'||b.provider_code and f.request_sha256_hex=b.request_sha256_hex
 where b.tenant_id=$1 and b.connection_id=$2 and b.provider_code=$3 and b.account_ref=$4 and b.live_mode=$5 and o.provider_reference=$6 and o.observed_at is not null and o.evidence_sha256_hex is not null`, c.TenantID, c.ConnectionID, c.ProviderCode, c.AccountRef, c.LiveMode, notice.ProviderReference).Scan(&attempt, &session, &fence)
	if err != nil {
		return errors.Join(ErrCheckoutRecovery, err)
	}
	if fence == "accepted" && session != "" {
		return nil
	}
	if fence != "unknown" && fence != "sending" {
		return ErrCheckoutRecovery
	}
	r, hash, known, err := base.binding(ctx, attempt)
	if err != nil {
		return err
	}
	var checkout paymentbridge.Checkout
	if known != "" {
		// Recover a crash between SaveCheckout and fence receipt without another
		// search or POST. Even a persisted session is re-observed through SDK GET.
		checkout, err = base.worker.Driver.RetrieveCheckout(ctx, known)
		if err == nil && checkout.SessionID != known {
			return ErrCheckoutRecovery
		}
	} else {
		driver, ok := base.worker.Driver.(paymentbridge.CheckoutRecoveryDriver)
		if !ok {
			return ErrCheckoutRecovery
		}
		checkout, err = driver.RecoverCheckout(ctx, r)
	}
	if err != nil {
		return errors.Join(ErrCheckoutRecovery, err)
	}
	if !checkoutMatches(c, r, checkout) {
		return ErrCheckoutRecovery
	}
	if _, _, err = base.notice(ctx, j); err != nil {
		return err
	}
	if err = base.worker.Store.SaveCheckout(ctx, c, r, hash, checkout); err != nil {
		return err
	}
	return base.recoverFence(ctx, r, hash, checkout)
}
