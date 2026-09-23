package officialpayments

// AUTHORED bounded recovery mapping over preference.Search/Get in the pinned
// Mercado Pago SDK. Search summaries locate identity; complete GET proves it.
import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/mercadopago/sdk-go/pkg/preference"
)

type PreferenceRecoveryRequest struct {
	OrderID, PaymentAttemptID, Currency, CollectorID string
	AmountMinor                                      int64
	MinorUnitExponent                                int
	LiveMode                                         bool
	ExpiresAt                                        time.Time
}

// Exactly one total result is required; ambiguous/truncated or older-than-search
// history cannot authorize a guessed preference. No mutation API is called.
func (c *MercadoPagoCheckoutClient) RecoverCheckout(ctx context.Context, r PreferenceRecoveryRequest) (CheckoutResult, error) {
	if c == nil || c.client == nil || ctx == nil || r.OrderID == "" || len(r.OrderID) > 200 || r.PaymentAttemptID == "" || len(r.PaymentAttemptID) > 200 || r.AmountMinor <= 0 || !currencyPattern.MatchString(strings.ToLower(r.Currency)) || r.MinorUnitExponent < 0 || r.MinorUnitExponent > 3 || r.ExpiresAt.IsZero() || r.LiveMode == c.sandbox {
		return CheckoutResult{}, ErrInvalidPaymentRequest
	}
	collector, err := strconv.ParseInt(r.CollectorID, 10, 64)
	if err != nil || collector <= 0 || strconv.FormatInt(collector, 10) != r.CollectorID {
		return CheckoutResult{}, ErrInvalidPaymentRequest
	}
	page, err := c.client.Search(ctx, preference.SearchRequest{Limit: 2, Offset: 0, Filters: map[string]string{"external_reference": r.OrderID, "site_id": "MLA"}})
	if err != nil {
		return CheckoutResult{}, err
	}
	if page == nil || page.Total != 1 || len(page.Elements) != 1 || page.NextOffset < 0 || page.NextOffset > 1 {
		return CheckoutResult{}, ErrInvalidObservation
	}
	item := page.Elements[0]
	if !preferenceID.MatchString(item.ID) || item.ExternalReference != r.OrderID || item.CollectorID != collector || item.SiteID != "MLA" || item.LiveMode != r.LiveMode || !item.Expires || !item.ExpirationDateTo.Equal(r.ExpiresAt) {
		return CheckoutResult{}, ErrInvalidObservation
	}
	complete, err := c.client.Get(ctx, item.ID)
	if err != nil {
		return CheckoutResult{}, err
	}
	result, err := c.preferenceResult(complete, r.MinorUnitExponent)
	if err != nil {
		return CheckoutResult{}, err
	}
	if result.SessionID != item.ID || result.OrderID != r.OrderID || result.PaymentAttemptID != r.PaymentAttemptID || result.Currency != r.Currency || result.AmountMinor != r.AmountMinor || result.CollectorID != r.CollectorID || result.LiveMode != r.LiveMode || !result.ExpiresAt.Equal(r.ExpiresAt) || !complete.Expires || complete.SiteID != "MLA" {
		return CheckoutResult{}, ErrInvalidObservation
	}
	return result, nil
}
