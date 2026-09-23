package paymentbridge

// AUTHORED optional recovery wiring. The pinned provider SDK performs both
// preference Search and Get; this path never calls CreateCheckout.
import (
	"context"
	officialpayments "example.com/elite/official-payment-webhooks/officialpayments"
)

type CheckoutRecoveryDriver interface {
	RecoverCheckout(context.Context, Request) (Checkout, error)
}

func (d *SDKDriver) RecoverCheckout(ctx context.Context, r Request) (Checkout, error) {
	if d == nil || ctx == nil || d.mpCheckout == nil || d.config.Scope.ProviderCode != "mercadopago" || r.ProviderCode != d.config.Scope.ProviderCode || r.TenantID != d.config.Scope.TenantID || r.OrganizationID != d.config.Scope.OrganizationID || r.Currency != d.config.Scope.Currency {
		return Checkout{}, ErrCheckoutConflict
	}
	if err := d.ValidateCredential(ctx); err != nil {
		return Checkout{}, err
	}
	result, err := d.mpCheckout.RecoverCheckout(ctx, officialpayments.PreferenceRecoveryRequest{OrderID: r.OrderID, PaymentAttemptID: r.PaymentAttemptID, Currency: r.Currency, CollectorID: d.config.Scope.AccountRef, AmountMinor: r.AmountMinor, MinorUnitExponent: d.config.Scope.MinorUnitExponent, LiveMode: d.config.Scope.LiveMode, ExpiresAt: r.CheckoutExpiresAt})
	if err != nil {
		return Checkout{}, ErrSDKDriverUnavailable
	}
	return d.checkout(result)
}
