package officialpayments

// AUTHORED request/response mapping for official hosted payment UIs. Payment
// credentials/card entry remain at the provider; a return URL is never receipt.
import (
	"context"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
	"github.com/mercadopago/sdk-go/pkg/requestoptions"
	"github.com/stripe/stripe-go/v86"
)

var checkoutSessionID = regexp.MustCompile(`^cs_(test|live)_[A-Za-z0-9]{1,200}$`)
var preferenceID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]{0,200}$`)

type CheckoutRequest struct {
	OrderID, PaymentAttemptID, DisplayName string
	AmountMinor                            int64
	Currency                               string
	MinorUnitExponent                      int
	SuccessURL, CancelURL, NotificationURL string
	IdempotencyKey                         string
	ExpiresAt                              time.Time
}
type CheckoutResult struct {
	ProviderCode     string    `json:"provider_code"`
	SessionID        string    `json:"session_id"`
	URL              string    `json:"url,omitempty"`
	PaymentReference string    `json:"payment_reference,omitempty"`
	OrderID          string    `json:"order_id"`
	PaymentAttemptID string    `json:"payment_attempt_id"`
	Currency         string    `json:"currency"`
	AmountMinor      int64     `json:"amount_minor"`
	Status           string    `json:"status"`
	PaymentStatus    string    `json:"payment_status,omitempty"`
	LiveMode         bool      `json:"live_mode"`
	CollectorID      string    `json:"collector_id,omitempty"`
	ExpiresAt        time.Time `json:"expires_at"`
}

func secureURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && len(raw) <= 8192 && u.Scheme == "https" && u.Hostname() != "" && u.User == nil && !strings.ContainsAny(raw, "\r\n") && (u.Port() == "" || u.Port() == "443")
}
func hostedURL(raw, host string) bool {
	if !secureURL(raw) {
		return false
	}
	u, _ := url.Parse(raw)
	return strings.EqualFold(u.Hostname(), host)
}
func validCheckout(r CheckoutRequest) bool {
	return r.AmountMinor > 0 && r.AmountMinor <= 1<<53-1 && currencyPattern.MatchString(strings.ToLower(r.Currency)) && r.OrderID != "" && len(r.OrderID) <= 200 && r.PaymentAttemptID != "" && len(r.PaymentAttemptID) <= 200 && strings.TrimSpace(r.DisplayName) != "" && len(r.DisplayName) <= 200 && idempotencyPattern.MatchString(r.IdempotencyKey) && secureURL(r.SuccessURL) && secureURL(r.CancelURL) && !r.ExpiresAt.IsZero() && r.ExpiresAt.After(time.Now().Add(30*time.Minute)) && r.ExpiresAt.Before(time.Now().Add(24*time.Hour))
}
func (c *StripeIntentClient) CreateCheckout(ctx context.Context, r CheckoutRequest) (CheckoutResult, error) {
	if c == nil || c.client == nil || ctx == nil || !validCheckout(r) {
		return CheckoutResult{}, ErrInvalidPaymentRequest
	}
	metadata := map[string]string{"order_id": r.OrderID, "payment_attempt_id": r.PaymentAttemptID}
	params := &stripe.CheckoutSessionCreateParams{Mode: stripe.String("payment"), SuccessURL: stripe.String(r.SuccessURL), CancelURL: stripe.String(r.CancelURL), ClientReferenceID: stripe.String(r.PaymentAttemptID), ExpiresAt: stripe.Int64(r.ExpiresAt.Unix()), Metadata: metadata, PaymentIntentData: &stripe.CheckoutSessionCreatePaymentIntentDataParams{Metadata: metadata}, LineItems: []*stripe.CheckoutSessionCreateLineItemParams{{Quantity: stripe.Int64(1), PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{Currency: stripe.String(strings.ToLower(r.Currency)), UnitAmount: stripe.Int64(r.AmountMinor), ProductData: &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{Name: stripe.String(r.DisplayName)}}}}}
	params.SetIdempotencyKey(r.IdempotencyKey)
	session, err := c.client.V1CheckoutSessions.Create(ctx, params)
	if err != nil {
		return CheckoutResult{}, err
	}
	result, err := stripeCheckoutResult(session)
	if err != nil || result.OrderID != r.OrderID || result.PaymentAttemptID != r.PaymentAttemptID || result.Currency != strings.ToUpper(r.Currency) || result.AmountMinor != r.AmountMinor || result.URL == "" {
		return CheckoutResult{}, ErrInvalidObservation
	}
	return result, nil
}
func (c *StripeIntentClient) RetrieveCheckout(ctx context.Context, id string) (CheckoutResult, error) {
	if c == nil || c.client == nil || ctx == nil || !checkoutSessionID.MatchString(id) {
		return CheckoutResult{}, ErrInvalidPaymentRequest
	}
	session, err := c.client.V1CheckoutSessions.Retrieve(ctx, id, &stripe.CheckoutSessionRetrieveParams{})
	if err != nil {
		return CheckoutResult{}, err
	}
	result, err := stripeCheckoutResult(session)
	if err == nil && result.SessionID != id {
		return CheckoutResult{}, ErrInvalidObservation
	}
	return result, err
}
func stripeCheckoutResult(s *stripe.CheckoutSession) (CheckoutResult, error) {
	if s == nil || !checkoutSessionID.MatchString(s.ID) || s.Mode != "payment" || s.AmountTotal <= 0 || !currencyPattern.MatchString(string(s.Currency)) || s.Metadata["order_id"] == "" || s.Metadata["payment_attempt_id"] == "" || s.ClientReferenceID != s.Metadata["payment_attempt_id"] || s.Status == "" || s.PaymentStatus == "" || (s.URL != "" && !hostedURL(s.URL, "checkout.stripe.com")) {
		return CheckoutResult{}, ErrInvalidObservation
	}
	ref := ""
	if s.PaymentIntent != nil {
		ref = s.PaymentIntent.ID
		if !stripeIntentID.MatchString(ref) {
			return CheckoutResult{}, ErrInvalidObservation
		}
	}
	return CheckoutResult{ProviderCode: "stripe", SessionID: s.ID, URL: s.URL, PaymentReference: ref, OrderID: s.Metadata["order_id"], PaymentAttemptID: s.Metadata["payment_attempt_id"], Currency: strings.ToUpper(string(s.Currency)), AmountMinor: s.AmountTotal, Status: string(s.Status), PaymentStatus: string(s.PaymentStatus), LiveMode: s.Livemode, ExpiresAt: time.Unix(s.ExpiresAt, 0).UTC()}, nil
}

// Argentina is the admitted Checkout Pro country lane for this reference.
// Other countries/custom checkout domains require their explicit provider lock.
type MercadoPagoCheckoutClient struct {
	client  preference.Client
	sandbox bool
}

func NewMercadoPagoCheckoutClient(token string, sandbox bool) (*MercadoPagoCheckoutClient, error) {
	return NewMercadoPagoCheckoutClientWithHTTPClient(token, sandbox, &http.Client{Timeout: 10 * time.Second})
}
func NewMercadoPagoCheckoutClientWithHTTPClient(token string, sandbox bool, httpClient *http.Client) (*MercadoPagoCheckoutClient, error) {
	if strings.TrimSpace(token) == "" || httpClient == nil {
		return nil, ErrInvalidPaymentRequest
	}
	bounded := *httpClient
	bounded.Timeout = 10 * time.Second
	bounded.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	cfg, err := config.New(token, config.WithHTTPClient(&bounded), config.WithMaxRetries(0))
	if err != nil {
		return nil, err
	}
	return &MercadoPagoCheckoutClient{client: preference.NewClient(cfg), sandbox: sandbox}, nil
}
func (c *MercadoPagoCheckoutClient) CreateCheckout(ctx context.Context, r CheckoutRequest) (CheckoutResult, error) {
	if c == nil || c.client == nil || ctx == nil || !validCheckout(r) || r.MinorUnitExponent < 0 || r.MinorUnitExponent > 3 || !secureURL(r.NotificationURL) {
		return CheckoutResult{}, ErrInvalidPaymentRequest
	}
	amount := float64(r.AmountMinor) / math.Pow10(r.MinorUnitExponent)
	roundTrip, err := decimalMinor(amount, r.MinorUnitExponent)
	if err != nil || roundTrip != r.AmountMinor {
		return CheckoutResult{}, ErrInvalidPaymentRequest
	}
	ctx = requestoptions.WithIdempotencyKey(ctx, r.IdempotencyKey)
	expires := r.ExpiresAt
	p, err := c.client.Create(ctx, preference.Request{Items: []preference.ItemRequest{{ID: r.OrderID, Title: r.DisplayName, Quantity: 1, CurrencyID: strings.ToUpper(r.Currency), UnitPrice: amount}}, ExternalReference: r.OrderID, Metadata: map[string]any{"order_id": r.OrderID, "payment_attempt_id": r.PaymentAttemptID}, NotificationURL: r.NotificationURL, BackURLs: &preference.BackURLsRequest{Success: r.SuccessURL, Pending: r.SuccessURL, Failure: r.CancelURL}, Expires: true, ExpirationDateTo: &expires})
	if err != nil {
		return CheckoutResult{}, err
	}
	result, err := c.preferenceResult(p, r.MinorUnitExponent)
	if err != nil || result.OrderID != r.OrderID || result.PaymentAttemptID != r.PaymentAttemptID || result.Currency != strings.ToUpper(r.Currency) || result.AmountMinor != r.AmountMinor {
		return CheckoutResult{}, ErrInvalidObservation
	}
	return result, nil
}
func (c *MercadoPagoCheckoutClient) RetrieveCheckout(ctx context.Context, id string, exponent int) (CheckoutResult, error) {
	if c == nil || c.client == nil || ctx == nil || !preferenceID.MatchString(id) {
		return CheckoutResult{}, ErrInvalidPaymentRequest
	}
	p, err := c.client.Get(ctx, id)
	if err != nil {
		return CheckoutResult{}, err
	}
	result, err := c.preferenceResult(p, exponent)
	if err == nil && result.SessionID != id {
		return CheckoutResult{}, ErrInvalidObservation
	}
	return result, err
}
func (c *MercadoPagoCheckoutClient) preferenceResult(p *preference.Response, exponent int) (CheckoutResult, error) {
	if p == nil || !preferenceID.MatchString(p.ID) || p.ExternalReference == "" || p.CollectorID <= 0 || len(p.Items) != 1 || p.Items[0].Quantity != 1 {
		return CheckoutResult{}, ErrInvalidObservation
	}
	order, ok := p.Metadata["order_id"].(string)
	attempt, ok2 := p.Metadata["payment_attempt_id"].(string)
	if !ok || !ok2 || order != p.ExternalReference || attempt == "" || !currencyPattern.MatchString(strings.ToLower(p.Items[0].CurrencyID)) {
		return CheckoutResult{}, ErrInvalidObservation
	}
	amount, err := decimalMinor(p.Items[0].UnitPrice, exponent)
	if err != nil || amount <= 0 {
		return CheckoutResult{}, ErrInvalidObservation
	}
	location, host := p.InitPoint, "www.mercadopago.com.ar"
	if c.sandbox {
		location, host = p.SandboxInitPoint, "sandbox.mercadopago.com.ar"
	}
	if !hostedURL(location, host) {
		return CheckoutResult{}, ErrInvalidObservation
	}
	return CheckoutResult{ProviderCode: "mercadopago", SessionID: p.ID, URL: location, OrderID: order, PaymentAttemptID: attempt, Currency: strings.ToUpper(p.Items[0].CurrencyID), AmountMinor: amount, Status: "preference-created", LiveMode: !c.sandbox, CollectorID: fmtInt(int(p.CollectorID)), ExpiresAt: p.ExpirationDateTo}, nil
}
