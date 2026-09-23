package paymentbridge

// AUTHORED read contract. This exposes a provider-hosted redirect, never payer
// data or a browser-authorized payment transition.
import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
)

var ErrCheckoutNotAvailable = errors.New("customer checkout is not available")

type CustomerCheckout struct {
	OrderID      string    `json:"order_id"`
	ProviderCode string    `json:"provider_code"`
	URL          string    `json:"url"`
	ExpiresAt    time.Time `json:"expires_at"`
}
type CustomerCheckoutReader interface {
	CustomerCheckout(context.Context, string, string, string, string) (CustomerCheckout, error)
}

func (c CustomerCheckout) Valid(now time.Time) bool {
	u, err := url.Parse(c.URL)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Fragment != "" || (u.Port() != "" && u.Port() != "443") || len(c.URL) > 8192 || !utf8.ValidString(c.URL) || strings.ContainsAny(c.URL, "\r\n") || !c.ExpiresAt.After(now) || c.OrderID == "" {
		return false
	}
	switch c.ProviderCode {
	case "stripe":
		return strings.EqualFold(u.Hostname(), "checkout.stripe.com")
	case "mercadopago":
		return strings.EqualFold(u.Hostname(), "www.mercadopago.com.ar") || strings.EqualFold(u.Hostname(), "sandbox.mercadopago.com.ar")
	}
	return false
}
