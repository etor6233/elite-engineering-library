package paymentbridge

// AUTHORED composition only: provider HTTP serialization, authentication and
// resource GETs are performed by the exact official SDK dependency pins.
import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	officialpayments "example.com/elite/official-payment-webhooks/officialpayments"
	mpconfig "github.com/mercadopago/sdk-go/pkg/config"
	mpuser "github.com/mercadopago/sdk-go/pkg/user"
)

var ErrSDKDriverUnavailable = errors.New("payment SDK operation unavailable")

type SDKDriverConfig struct {
	Scope                                  Scope
	SuccessURL, CancelURL, NotificationURL string
	DisplayName                            string
}
type SDKDriver struct {
	config     SDKDriverConfig
	stripe     *officialpayments.StripeIntentClient
	mpPayment  *officialpayments.MercadoPagoPaymentClient
	mpCheckout *officialpayments.MercadoPagoCheckoutClient
	mpUser     mpuser.Client
}

func configuredHTTPS(value string) bool {
	u, err := url.Parse(value)
	return err == nil && u.Scheme == "https" && u.Hostname() != "" && u.User == nil && len(value) <= 2048 && !strings.ContainsAny(value, "\r\n") && (u.Port() == "" || u.Port() == "443")
}

// No network is used by construction. ValidateCredential must succeed at host
// startup before taking a durable send claim; CreateCheckout rechecks it.
func NewSDKDriver(c SDKDriverConfig, token string, client *http.Client) (*SDKDriver, error) {
	if c.Scope.Validate() != nil || strings.TrimSpace(token) == "" || len(token) > 16384 || !configuredHTTPS(c.SuccessURL) || !configuredHTTPS(c.CancelURL) || strings.TrimSpace(c.DisplayName) == "" || len(c.DisplayName) > 200 {
		return nil, ErrCheckoutConflict
	}
	if c.Scope.ProviderCode == "mercadopago" && (!configuredHTTPS(c.NotificationURL) || c.Scope.Currency != "ARS" || c.Scope.MinorUnitExponent != 2) {
		return nil, ErrCheckoutConflict
	}
	// Stripe's documented test/live secret and restricted key prefixes provide
	// an early configuration guard; the official account GET still proves identity.
	if c.Scope.ProviderCode == "stripe" {
		mode := "test_"
		if c.Scope.LiveMode {
			mode = "live_"
		}
		if (!strings.HasPrefix(token, "sk_"+mode) && !strings.HasPrefix(token, "rk_"+mode)) || len(token) <= len("sk_"+mode) || strings.TrimSpace(token) != token || strings.ContainsAny(token, "\r\n") {
			return nil, ErrPaymentMismatch
		}
	}
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	bounded := *client
	bounded.Timeout = 10 * time.Second
	bounded.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	d := &SDKDriver{config: c}
	var err error
	if c.Scope.ProviderCode == "stripe" {
		d.stripe, err = officialpayments.NewStripeIntentClientWithHTTPClient(token, &bounded)
	} else {
		d.mpPayment, err = officialpayments.NewMercadoPagoPaymentClientWithHTTPClient(token, &bounded)
		if err != nil {
			return nil, ErrCheckoutConflict
		}
		d.mpCheckout, err = officialpayments.NewMercadoPagoCheckoutClientWithHTTPClient(token, !c.Scope.LiveMode, &bounded)
		if err != nil {
			return nil, ErrCheckoutConflict
		}
		cfg, configErr := mpconfig.New(token, mpconfig.WithHTTPClient(&bounded), mpconfig.WithMaxRetries(0))
		if configErr != nil {
			return nil, ErrCheckoutConflict
		}
		d.mpUser = mpuser.NewClient(cfg)
	}
	if err != nil {
		return nil, ErrCheckoutConflict
	}
	return d, nil
}
func (d *SDKDriver) ValidateCredential(ctx context.Context) error {
	if d == nil || ctx == nil {
		return ErrCheckoutConflict
	}
	if d.stripe != nil {
		id, err := d.stripe.ProbeAccount(ctx)
		if err != nil {
			return ErrSDKDriverUnavailable
		}
		if id != d.config.Scope.AccountRef {
			return ErrPaymentMismatch
		}
		return nil
	}
	if d.mpUser == nil {
		return ErrCheckoutConflict
	}
	account, err := d.mpUser.Get(ctx)
	if err != nil {
		return ErrSDKDriverUnavailable
	}
	if account == nil || strconv.Itoa(account.ID) != d.config.Scope.AccountRef || account.CountryID != "AR" || account.SiteID != "MLA" {
		return ErrPaymentMismatch
	}
	return nil
}
func (d *SDKDriver) CreateCheckout(ctx context.Context, r Request) (Checkout, error) {
	if d == nil || ctx == nil || r.TenantID != d.config.Scope.TenantID || r.OrganizationID != d.config.Scope.OrganizationID || r.ProviderCode != d.config.Scope.ProviderCode || r.Currency != d.config.Scope.Currency || r.AmountMinor <= 0 || r.PaymentAttemptID == "" || r.OrderID == "" || r.CheckoutExpiresAt.IsZero() {
		return Checkout{}, ErrPaymentMismatch
	}
	if err := d.ValidateCredential(ctx); err != nil {
		return Checkout{}, err
	}
	request := officialpayments.CheckoutRequest{OrderID: r.OrderID, PaymentAttemptID: r.PaymentAttemptID, AmountMinor: r.AmountMinor, Currency: r.Currency, MinorUnitExponent: d.config.Scope.MinorUnitExponent, DisplayName: d.config.DisplayName, SuccessURL: d.config.SuccessURL, CancelURL: d.config.CancelURL, NotificationURL: d.config.NotificationURL, IdempotencyKey: r.PaymentAttemptID, ExpiresAt: r.CheckoutExpiresAt}
	var result officialpayments.CheckoutResult
	var err error
	if d.stripe != nil {
		result, err = d.stripe.CreateCheckout(ctx, request)
	} else {
		result, err = d.mpCheckout.CreateCheckout(ctx, request)
	}
	if err != nil {
		return Checkout{}, ErrSDKDriverUnavailable
	}
	checkout, err := d.checkout(result)
	if err != nil || checkout.OrderID != r.OrderID || checkout.PaymentAttemptID != r.PaymentAttemptID || checkout.AmountMinor != r.AmountMinor || checkout.ExpiresAt.Unix() != r.CheckoutExpiresAt.Unix() {
		return Checkout{}, ErrPaymentMismatch
	}
	return checkout, nil
}
func (d *SDKDriver) RetrieveCheckout(ctx context.Context, id string) (Checkout, error) {
	if d == nil || ctx == nil || id == "" {
		return Checkout{}, ErrCheckoutConflict
	}
	if err := d.ValidateCredential(ctx); err != nil {
		return Checkout{}, err
	}
	var result officialpayments.CheckoutResult
	var err error
	if d.stripe != nil {
		result, err = d.stripe.RetrieveCheckout(ctx, id)
	} else {
		result, err = d.mpCheckout.RetrieveCheckout(ctx, id, d.config.Scope.MinorUnitExponent)
	}
	if err != nil {
		return Checkout{}, ErrSDKDriverUnavailable
	}
	return d.checkout(result)
}
func (d *SDKDriver) checkout(r officialpayments.CheckoutResult) (Checkout, error) {
	if r.ProviderCode != d.config.Scope.ProviderCode || r.Currency != d.config.Scope.Currency || r.LiveMode != d.config.Scope.LiveMode || r.OrderID == "" || r.PaymentAttemptID == "" || r.AmountMinor <= 0 || r.ExpiresAt.IsZero() {
		return Checkout{}, ErrPaymentMismatch
	}
	if r.ProviderCode == "mercadopago" && r.CollectorID != d.config.Scope.AccountRef {
		return Checkout{}, ErrPaymentMismatch
	}
	return Checkout{ProviderCode: r.ProviderCode, SessionID: r.SessionID, URL: r.URL, PaymentReference: r.PaymentReference, Status: r.Status, OrderID: r.OrderID, PaymentAttemptID: r.PaymentAttemptID, Currency: r.Currency, AccountRef: d.config.Scope.AccountRef, AmountMinor: r.AmountMinor, LiveMode: r.LiveMode, ExpiresAt: r.ExpiresAt}, nil
}
func (d *SDKDriver) RetrievePayment(ctx context.Context, id string) (officialpayments.PaymentSnapshot, error) {
	if d == nil || ctx == nil || id == "" {
		return officialpayments.PaymentSnapshot{}, ErrCheckoutConflict
	}
	if err := d.ValidateCredential(ctx); err != nil {
		return officialpayments.PaymentSnapshot{}, err
	}
	var result officialpayments.PaymentSnapshot
	var err error
	if d.stripe != nil {
		result, err = d.stripe.RetrieveIntent(ctx, id)
	} else {
		result, err = d.mpPayment.RetrievePayment(ctx, id, d.config.Scope.MinorUnitExponent)
	}
	if err != nil {
		return officialpayments.PaymentSnapshot{}, ErrSDKDriverUnavailable
	}
	if result.ProviderReference != id || result.ProviderCode != d.config.Scope.ProviderCode || result.Currency != d.config.Scope.Currency || result.LiveMode != d.config.Scope.LiveMode || (result.ProviderCode == "mercadopago" && result.CollectorID != d.config.Scope.AccountRef) {
		return officialpayments.PaymentSnapshot{}, ErrPaymentMismatch
	}
	return result, nil
}
