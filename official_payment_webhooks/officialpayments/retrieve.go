package officialpayments

// AUTHORED boundary mapping over the exact official SDK pins. This adapter
// observes provider resources; it never decides order release or edits a ledger.
import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/stripe/stripe-go/v86"
)

var ErrInvalidObservation = errors.New("payment observation does not match contract")
var stripeIntentID = regexp.MustCompile(`^pi_[A-Za-z0-9]{1,200}$`)
var decimalPaymentID = regexp.MustCompile(`^[1-9][0-9]{0,18}$`)
var stripeAccountID = regexp.MustCompile(`^acct_[A-Za-z0-9]{1,200}$`)

// ProbeAccount observes the account associated with this exact credential.
// It uses /v1/account, never an ID chosen from request input or /accounts/{id}.
func (c *StripeIntentClient) ProbeAccount(ctx context.Context) (string, error) {
	if c == nil || c.client == nil || ctx == nil {
		return "", ErrInvalidPaymentRequest
	}
	account, err := c.client.V1Accounts.Retrieve(ctx, &stripe.AccountRetrieveParams{})
	if err != nil {
		return "", err
	}
	if account == nil || !stripeAccountID.MatchString(account.ID) {
		return "", ErrInvalidObservation
	}
	return account.ID, nil
}

// PaymentSnapshot contains only reconciliation fields, never client secrets,
// payer details or card data. Status retains the provider's original vocabulary.
type PaymentSnapshot struct {
	ProviderCode      string `json:"provider_code"`
	ProviderReference string `json:"provider_reference"`
	OrderID           string `json:"order_id"`
	PaymentAttemptID  string `json:"payment_attempt_id,omitempty"`
	Currency          string `json:"currency"`
	AmountMinor       int64  `json:"amount_minor"`
	ReceivedMinor     int64  `json:"received_minor"`
	RefundedMinor     int64  `json:"refunded_minor"`
	CapturableMinor   int64  `json:"capturable_minor"`
	Status            string `json:"status"`
	LiveMode          bool   `json:"live_mode"`
	Disputed          bool   `json:"disputed"`
	CollectorID       string `json:"collector_id,omitempty"`
}

// SHA256 binds the normalized observation, not a provider-signed receipt.
func (s PaymentSnapshot) SHA256() string {
	data, _ := json.Marshal(s)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func (c *StripeIntentClient) RetrieveIntent(ctx context.Context, id string) (PaymentSnapshot, error) {
	if c == nil || c.client == nil || ctx == nil || !stripeIntentID.MatchString(id) {
		return PaymentSnapshot{}, ErrInvalidPaymentRequest
	}
	params := &stripe.PaymentIntentRetrieveParams{}
	params.AddExpand("latest_charge")
	p, err := c.client.V1PaymentIntents.Retrieve(ctx, id, params)
	if err != nil {
		return PaymentSnapshot{}, err
	}
	if p == nil || p.ID != id || p.Amount <= 0 || p.AmountReceived < 0 || p.AmountReceived > p.Amount || p.AmountCapturable < 0 || p.AmountCapturable > p.Amount || !currencyPattern.MatchString(string(p.Currency)) || strings.TrimSpace(p.Metadata["order_id"]) == "" || len(p.Metadata["order_id"]) > 200 || p.Status == "" {
		return PaymentSnapshot{}, ErrInvalidObservation
	}
	refunded := int64(0)
	disputed := false
	if p.AmountReceived > 0 {
		charge := p.LatestCharge
		if charge == nil || charge.ID == "" || charge.Amount != p.Amount || charge.AmountCaptured != p.AmountReceived || charge.AmountRefunded < 0 || charge.AmountRefunded > charge.AmountCaptured || charge.Currency != p.Currency || !charge.Captured || charge.Livemode != p.Livemode {
			return PaymentSnapshot{}, ErrInvalidObservation
		}
		refunded = charge.AmountRefunded
		disputed = charge.Disputed
	}
	return PaymentSnapshot{ProviderCode: "stripe", ProviderReference: p.ID, OrderID: p.Metadata["order_id"], PaymentAttemptID: p.Metadata["payment_attempt_id"], Currency: strings.ToUpper(string(p.Currency)), AmountMinor: p.Amount, ReceivedMinor: p.AmountReceived, RefundedMinor: refunded, CapturableMinor: p.AmountCapturable, Status: string(p.Status), LiveMode: p.Livemode, Disputed: disputed}, nil
}

// decimalMinor converts the SDK's decimal JSON float representation without
// rounding a fractional minor unit. The bounded exact integer is checked again
// by the calling owner against the durable order currency/amount.
func decimalMinor(value float64, exponent int) (int64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || exponent < 0 || exponent > 3 {
		return 0, ErrInvalidObservation
	}
	amount, ok := new(big.Rat).SetString(strconv.FormatFloat(value, 'f', -1, 64))
	if !ok {
		return 0, ErrInvalidObservation
	}
	amount.Mul(amount, new(big.Rat).SetInt64(int64(math.Pow10(exponent))))
	if !amount.IsInt() || !amount.Num().IsInt64() {
		return 0, ErrInvalidObservation
	}
	minor := amount.Num().Int64()
	if minor > 1<<53-1 {
		return 0, ErrInvalidObservation
	}
	return minor, nil
}

// RetrievePayment requires the currency exponent from the selected provider
// contract; no currency table or merchant financial policy is invented here.
func (c *MercadoPagoPaymentClient) RetrievePayment(ctx context.Context, id string, exponent int) (PaymentSnapshot, error) {
	if c == nil || c.client == nil || ctx == nil || !decimalPaymentID.MatchString(id) || exponent < 0 || exponent > 3 {
		return PaymentSnapshot{}, ErrInvalidPaymentRequest
	}
	numeric, err := strconv.Atoi(id)
	if err != nil || numeric <= 0 {
		return PaymentSnapshot{}, ErrInvalidPaymentRequest
	}
	p, err := c.client.Get(ctx, numeric)
	if err != nil {
		return PaymentSnapshot{}, err
	}
	if p == nil || p.ID != numeric || !currencyPattern.MatchString(strings.ToLower(p.CurrencyID)) || strings.TrimSpace(p.ExternalReference) == "" || len(p.ExternalReference) > 200 || p.Status == "" || p.CollectorID <= 0 {
		return PaymentSnapshot{}, ErrInvalidObservation
	}
	amount, err := decimalMinor(p.TransactionAmount, exponent)
	if err != nil || amount <= 0 {
		return PaymentSnapshot{}, ErrInvalidObservation
	}
	refunded, err := decimalMinor(p.TransactionAmountRefunded, exponent)
	if err != nil || refunded > amount {
		return PaymentSnapshot{}, ErrInvalidObservation
	}
	received := int64(0)
	if p.Captured {
		received = amount
	}
	attempt, _ := p.Metadata["payment_attempt_id"].(string)
	return PaymentSnapshot{ProviderCode: "mercadopago", ProviderReference: id, OrderID: p.ExternalReference, PaymentAttemptID: attempt, Currency: strings.ToUpper(p.CurrencyID), AmountMinor: amount, ReceivedMinor: received, RefundedMinor: refunded, Status: p.Status, LiveMode: p.LiveMode, Disputed: p.Status == "charged_back", CollectorID: strconv.FormatInt(p.CollectorID, 10)}, nil
}
