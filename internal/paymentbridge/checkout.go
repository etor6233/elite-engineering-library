package paymentbridge

// AUTHORED orchestration glue. Provider effects belong to the pinned SDK adapter;
// the existing outbound fence prevents an uncertain call from being sent twice.
import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	officialpayments "example.com/elite/official-payment-webhooks/officialpayments"
)

var ErrCheckoutConflict = errors.New("payment checkout binding conflict")
var ErrObservationStale = errors.New("payment observation superseded")
var ErrPaymentMismatch = errors.New("provider observation differs from payment contract")

type Scope struct {
	TenantID, OrganizationID, ConnectionID, ProviderCode, AccountRef string
	Currency                                                         string
	MinorUnitExponent                                                int
	LiveMode                                                         bool
}

func (s Scope) Validate() error {
	for _, v := range []string{s.TenantID, s.OrganizationID, s.ConnectionID, s.AccountRef} {
		if strings.TrimSpace(v) != v || v == "" || len(v) > 200 || strings.ContainsAny(v, "\r\n") {
			return ErrCheckoutConflict
		}
	}
	if (s.ProviderCode != "stripe" && s.ProviderCode != "mercadopago") || len(s.Currency) != 3 || s.Currency != strings.ToUpper(s.Currency) || s.MinorUnitExponent < 0 || s.MinorUnitExponent > 3 {
		return ErrCheckoutConflict
	}
	return nil
}

type Request struct {
	CheckoutExpiresAt                                                                            time.Time
	TenantID, OrganizationID, PaymentAttemptID, OrderID, CustomerSubject, ProviderCode, Currency string
	AmountMinor                                                                                  int64
}

type Checkout struct {
	ProviderCode, SessionID, URL, PaymentReference, Status string
	OrderID, PaymentAttemptID, Currency, AccountRef        string
	AmountMinor                                            int64
	LiveMode                                               bool
	ExpiresAt                                              time.Time
}
type Driver interface {
	CreateCheckout(context.Context, Request) (Checkout, error)
	RetrieveCheckout(context.Context, string) (Checkout, error)
	RetrievePayment(context.Context, string) (officialpayments.PaymentSnapshot, error)
}

type CheckoutStore interface {
	PendingRequests(context.Context, Scope, int) ([]Request, error)
	BindRequest(context.Context, Scope, Request, string) error
	SaveCheckout(context.Context, Scope, Request, string, Checkout) error
	Checkout(context.Context, Scope, string) (Request, Checkout, error)
	BeginObservation(context.Context, Scope, string, string) (Request, int64, error)
	RecordObservation(context.Context, Scope, Request, int64, officialpayments.PaymentSnapshot) error
}

type Worker struct {
	Scope  Scope
	Store  CheckoutStore
	Fence  outbounddelivery.Store
	Driver Driver
}

func (w Worker) valid() error {
	if w.Scope.Validate() != nil || w.Store == nil || w.Fence == nil || w.Driver == nil {
		return ErrCheckoutConflict
	}
	return nil
}

// message binds immutable database values, not caller-supplied prices or URLs.
func (r Request) message() (channels.Message, string, error) {
	if r.AmountMinor <= 0 || r.PaymentAttemptID == "" || r.OrderID == "" || r.CustomerSubject == "" {
		return channels.Message{}, "", ErrCheckoutConflict
	}
	data, err := json.Marshal(r)
	if err != nil {
		return channels.Message{}, "", err
	}
	m := channels.Message{ChannelCode: "payment_" + r.ProviderCode, TenantID: r.TenantID, ExternalID: r.CustomerSubject, ThreadID: r.OrderID, DeliveryKey: r.PaymentAttemptID, Direction: channels.DirectionOut, Text: string(data)}
	hash, err := outbounddelivery.MessageSHA256(m)
	return m, hash, err
}

func (w Worker) ProcessOnce(ctx context.Context, limit int) (int, error) {
	if w.valid() != nil || limit < 1 || limit > 100 {
		return 0, ErrCheckoutConflict
	}
	requests, err := w.Store.PendingRequests(ctx, w.Scope, limit)
	if err != nil {
		return 0, err
	}
	completed := 0
	for _, r := range requests {
		if ctx.Err() != nil {
			return completed, ctx.Err()
		}
		m, hash, err := r.message()
		if err != nil {
			return completed, err
		}
		if err = w.Store.BindRequest(ctx, w.Scope, r, hash); err != nil {
			return completed, err
		}
		claim, err := w.Fence.Claim(ctx, m, hash)
		if errors.Is(err, outbounddelivery.ErrUnknown) || errors.Is(err, outbounddelivery.ErrTerminal) || errors.Is(err, outbounddelivery.ErrInProgress) {
			continue
		}
		if err != nil {
			return completed, err
		}
		if claim.Replay {
			continue
		}
		checkout, err := w.Driver.CreateCheckout(ctx, r)
		if err != nil {
			mark := w.Fence.MarkUnknown(ctx, m, hash, "CHECKOUT_CALL_UNCERTAIN")
			if mark != nil {
				return completed, mark
			}
			continue
		}
		if checkout.ProviderCode != r.ProviderCode || checkout.SessionID == "" || checkout.URL == "" || checkout.Status == "" {
			mark := w.Fence.MarkUnknown(ctx, m, hash, "CHECKOUT_RESPONSE_INVALID")
			if mark != nil {
				return completed, mark
			}
			continue
		}
		if err = w.Store.SaveCheckout(ctx, w.Scope, r, hash, checkout); err != nil {
			_ = w.Fence.MarkUnknown(ctx, m, hash, "CHECKOUT_PERSISTENCE_UNCERTAIN")
			return completed, err
		}
		// URLs may contain opaque session tokens; only normalized identity is hashed.
		proof, _ := json.Marshal(struct{ Provider, Session, Payment, Status string }{checkout.ProviderCode, checkout.SessionID, checkout.PaymentReference, checkout.Status})
		sum := sha256.Sum256(proof)
		receipt := outbounddelivery.Receipt{ProviderMessageID: checkout.SessionID, EvidenceSHA256: hex.EncodeToString(sum[:]), AcceptedAt: time.Now().UTC()}
		if err = w.Fence.Complete(ctx, m, hash, receipt); err != nil {
			return completed, err
		}
		completed++
	}
	return completed, nil
}

// Reconcile starts a generation before GET. A later callback invalidates that
// generation, so a delayed older response cannot clear the handover hold.
func (w Worker) Reconcile(ctx context.Context, attempt, providerReference string) error {
	if w.valid() != nil || attempt == "" || providerReference == "" {
		return ErrCheckoutConflict
	}
	r, generation, err := w.Store.BeginObservation(ctx, w.Scope, attempt, providerReference)
	if err != nil {
		return err
	}
	observed, err := w.Driver.RetrievePayment(ctx, providerReference)
	if err != nil {
		return err
	}
	return w.Store.RecordObservation(ctx, w.Scope, r, generation, observed)
}
