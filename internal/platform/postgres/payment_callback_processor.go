package postgres

// AUTHORED composition glue. Signed callbacks are resource hints; the pinned
// SDK GET and existing checkout/observation owners establish payment facts.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/paymentbridge"
	officialpayments "example.com/elite/official-payment-webhooks/officialpayments"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrPaymentCallback = errors.New("payment callback scope or durable receipt invalid")

type PaymentCallbackProcessor struct {
	pool     *pgxpool.Pool
	worker   paymentbridge.Worker
	workerID string
}

func NewPaymentCallbackProcessor(pool *pgxpool.Pool, worker paymentbridge.Worker, workerID string) (*PaymentCallbackProcessor, error) {
	if pool == nil || worker.Scope.Validate() != nil || worker.Store == nil || worker.Driver == nil || worker.Fence == nil || strings.TrimSpace(workerID) == "" || len(workerID) > 128 {
		return nil, ErrPaymentCallback
	}
	return &PaymentCallbackProcessor{pool: pool, worker: worker, workerID: workerID}, nil
}

// Claim selects only this connection. A worker never exhausts another tenant's
// valid payment job merely because its configured credentials differ.
type PaymentCallbackJobs struct {
	*Jobs
	Scope paymentbridge.Scope
}

func (s *PaymentCallbackJobs) Claim(ctx context.Context, queue, worker string, lease time.Duration, limit int) ([]Job, error) {
	if s == nil || s.Jobs == nil || s.Scope.Validate() != nil || queue != "payment-provider-events" {
		return nil, ErrPaymentCallback
	}
	match, _ := json.Marshal(map[string]string{"connection_id": s.Scope.ConnectionID, "provider_code": s.Scope.ProviderCode, "event_type": "payment.reconciliation-requested.v1"})
	return s.Jobs.ClaimScoped(ctx, queue, worker, lease, limit, JobScope{TenantID: s.Scope.TenantID, JobType: "provider.webhook.received", SchemaVersion: 1, PayloadMatch: match})
}

type paymentCallbackJob struct {
	ConnectionID    string `json:"connection_id"`
	ProviderCode    string `json:"provider_code"`
	ProviderEventID string `json:"provider_event_id"`
	EventType       string `json:"event_type"`
}

func decodePaymentCallback(raw []byte, value any) error {
	if len(raw) == 0 || len(raw) > 16384 {
		return ErrPaymentCallback
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(value) != nil || d.Decode(new(any)) != io.EOF {
		return ErrPaymentCallback
	}
	return nil
}

func (p *PaymentCallbackProcessor) notice(ctx context.Context, j Job) (paymentCallbackJob, officialpayments.PaymentNotification, error) {
	var payload paymentCallbackJob
	var n officialpayments.PaymentNotification
	c := p.worker.Scope
	if j.TenantID != c.TenantID || j.Queue != "payment-provider-events" || j.JobType != "provider.webhook.received" || j.SchemaVersion != 1 || j.Attempts < 1 || decodePaymentCallback(j.Payload, &payload) != nil || payload.ConnectionID != c.ConnectionID || payload.ProviderCode != c.ProviderCode || payload.EventType != "payment.reconciliation-requested.v1" || payload.ProviderEventID == "" {
		return payload, n, ErrPaymentCallback
	}
	var raw []byte
	var digest string
	err := p.pool.QueryRow(ctx, `select e.payload,e.body_sha256_hex from platform.job j
 join integration.webhook_event e on e.tenant_id=j.tenant_id and e.connection_id=$9 and e.provider_event_id=$10 and e.provider_code=$11 and e.event_type=$12
 join integration.provider_connection c on c.tenant_id=e.tenant_id and c.connection_id=e.connection_id and c.provider_code=e.provider_code and c.state='active' and (c.organization_id is null or c.organization_id=$13)
 where j.tenant_id=$1 and j.job_id=$2 and j.claimed_by=$3 and j.attempts=$4 and j.queue=$5 and j.job_type=$6 and j.schema_version=$7 and j.payload=$8::jsonb
 and j.completed_at is null and j.terminal_error_code is null and j.claimed_until>clock_timestamp()`, j.TenantID, j.JobID, p.workerID, j.Attempts, j.Queue, j.JobType, j.SchemaVersion, j.Payload, c.ConnectionID, payload.ProviderEventID, c.ProviderCode, payload.EventType, c.OrganizationID).Scan(&raw, &digest)
	if err != nil {
		return payload, n, errors.Join(ErrPaymentCallback, err)
	}
	var inbox paymentbridge.InboxPayload
	if decodePaymentCallback(raw, &inbox) != nil || inbox.Schema != "elite-payment-notification/v1" {
		return payload, n, ErrPaymentCallback
	}
	n = inbox.Notification
	if n.ProviderCode != c.ProviderCode || n.ProviderEventID != payload.ProviderEventID || n.BodySHA256 != digest || len(digest) != 64 || n.ProviderReference == "" || n.Type == "" {
		return payload, n, ErrPaymentCallback
	}
	if _, err = hex.DecodeString(digest); err != nil {
		return payload, n, ErrPaymentCallback
	}
	if (c.ProviderCode == "stripe" && n.ResourceType != "checkout_session" && n.ResourceType != "payment_intent") || (c.ProviderCode == "mercadopago" && n.ResourceType != "payment") {
		return payload, n, ErrPaymentCallback
	}
	return payload, n, nil
}

// binding also accepts a durable pre-dispatch skeleton. This permits recovery
// after a lost POST response without issuing a second provider POST.
func (p *PaymentCallbackProcessor) binding(ctx context.Context, attempt string) (paymentbridge.Request, string, string, error) {
	var r paymentbridge.Request
	var hash, session string
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return r, hash, session, err
	}
	defer tx.Rollback(ctx)
	c := p.worker.Scope
	if err = checkoutConnection(ctx, tx, c); err != nil {
		return r, hash, session, err
	}
	r, err = lockCheckoutRequest(ctx, tx, c, attempt)
	if err != nil {
		return r, hash, session, err
	}
	err = tx.QueryRow(ctx, `select b.request_sha256_hex,coalesce(b.session_id,'') from payment.provider_checkout b
 join communication.outbound_delivery f on f.tenant_id=b.tenant_id and f.delivery_key=b.payment_attempt_id and f.channel_code='payment_'||b.provider_code and f.request_sha256_hex=b.request_sha256_hex and f.state in ('sending','unknown','accepted')
 where b.tenant_id=$1 and b.payment_attempt_id=$2 and b.provider_code=$3 and b.connection_id=$4 and b.account_ref=$5 and b.live_mode=$6 and b.expires_at=$7`, c.TenantID, attempt, c.ProviderCode, c.ConnectionID, c.AccountRef, c.LiveMode, r.CheckoutExpiresAt).Scan(&hash, &session)
	if err != nil {
		return r, hash, session, errors.Join(ErrPaymentCallback, err)
	}
	_, actualHash, err := paymentCallbackMessage(r)
	if err != nil || actualHash != hash {
		return r, hash, session, ErrPaymentCallback
	}
	return r, hash, session, tx.Commit(ctx)
}

func paymentCallbackMessage(r paymentbridge.Request) (channels.Message, string, error) {
	raw, err := json.Marshal(r)
	if err != nil {
		return channels.Message{}, "", err
	}
	m := channels.Message{TenantID: r.TenantID, ChannelCode: "payment_" + r.ProviderCode, ExternalID: r.CustomerSubject, ThreadID: r.OrderID, DeliveryKey: r.PaymentAttemptID, Direction: channels.DirectionOut, Text: string(raw)}
	hash, err := outbounddelivery.MessageSHA256(m)
	return m, hash, err
}

func checkoutMatches(c paymentbridge.Scope, r paymentbridge.Request, v paymentbridge.Checkout) bool {
	return v.ProviderCode == c.ProviderCode && v.AccountRef == c.AccountRef && v.LiveMode == c.LiveMode && v.OrderID == r.OrderID && v.PaymentAttemptID == r.PaymentAttemptID && v.Currency == r.Currency && v.AmountMinor == r.AmountMinor && v.ExpiresAt.Equal(r.CheckoutExpiresAt) && v.SessionID != ""
}

func (p *PaymentCallbackProcessor) Handle(ctx context.Context, j Job) error {
	if p == nil || p.pool == nil {
		return ErrPaymentCallback
	}
	payload, n, err := p.notice(ctx, j)
	if err != nil {
		return err
	}
	c := p.worker.Scope
	attempt, reference := "", n.ProviderReference
	if n.ResourceType == "checkout_session" {
		observed, err := p.worker.Driver.RetrieveCheckout(ctx, reference)
		if err != nil {
			return err
		}
		if observed.SessionID != reference || observed.PaymentReference == "" {
			return ErrPaymentCallback
		}
		attempt = observed.PaymentAttemptID
		r, hash, session, err := p.binding(ctx, attempt)
		if err != nil {
			return err
		}
		if !checkoutMatches(c, r, observed) || (session != "" && session != observed.SessionID) {
			return ErrPaymentCallback
		}
		// Refresh hold before persisting newly learned session/payment linkage.
		// The payment observation itself still comes exclusively from Reconcile GET.
		if _, _, err = p.worker.Store.BeginObservation(ctx, c, attempt, observed.PaymentReference); err != nil {
			return err
		}
		if err = p.worker.Store.SaveCheckout(ctx, c, r, hash, observed); err != nil {
			return err
		}
		if err = p.recoverFence(ctx, r, hash, observed); err != nil {
			return err
		}
		reference = observed.PaymentReference
	} else {
		// The common refund/status path resolves the already bound provider ID.
		var matches int
		err = p.pool.QueryRow(ctx, `select count(*),coalesce(min(b.payment_attempt_id),'') from payment.provider_checkout b left join payment.provider_observation o using(tenant_id,payment_attempt_id)
 where b.tenant_id=$1 and b.connection_id=$2 and b.provider_code=$3 and b.account_ref=$4 and b.live_mode=$5 and (b.payment_reference=$6 or o.provider_reference=$6)`, c.TenantID, c.ConnectionID, c.ProviderCode, c.AccountRef, c.LiveMode, reference).Scan(&matches, &attempt)
		if err != nil {
			return err
		}
		if matches > 1 {
			return ErrPaymentCallback
		}
		if matches == 0 {
			// First payment callback may precede checkout completion. Metadata from an
			// authenticated provider GET locates the attempt; callback body cannot.
			snapshot, getErr := p.worker.Driver.RetrievePayment(ctx, reference)
			if getErr != nil {
				return getErr
			}
			attempt = snapshot.PaymentAttemptID
			// Mercado Pago preferences carry our order external_reference; payment
			// metadata is not guaranteed to inherit the preference metadata. Only
			// the authenticated GET may select a unique previously dispatched
			// whole-order request. Zero/ambiguous matches never pick an attempt.
			if attempt == "" && c.ProviderCode == "mercadopago" {
				var candidates int
				err = p.pool.QueryRow(ctx, `select count(*),coalesce(min(b.payment_attempt_id),'') from payment.provider_checkout b
 join payment.payment_attempt p using(tenant_id,payment_attempt_id)
 join sales.customer_order o using(tenant_id,order_id)
 join communication.outbound_delivery f on f.tenant_id=b.tenant_id and f.delivery_key=b.payment_attempt_id and f.channel_code='payment_'||b.provider_code and f.request_sha256_hex=b.request_sha256_hex and f.state in ('sending','unknown','accepted')
 where b.tenant_id=$1 and b.connection_id=$2 and b.provider_code=$3 and b.account_ref=$4 and b.live_mode=$5 and p.order_id=$6 and o.organization_id=$7 and p.currency=$8 and p.amount_minor_units=$9
 and (b.payment_reference is null or b.payment_reference=$10)`, c.TenantID, c.ConnectionID, c.ProviderCode, c.AccountRef, c.LiveMode, snapshot.OrderID, c.OrganizationID, snapshot.Currency, snapshot.AmountMinor, reference).Scan(&candidates, &attempt)
				if err != nil {
					return err
				}
				if candidates != 1 {
					return ErrPaymentCallback
				}
			}
			r, _, _, bindErr := p.binding(ctx, attempt)
			if bindErr != nil {
				return bindErr
			}
			if snapshot.ProviderCode != c.ProviderCode || snapshot.ProviderReference != reference || snapshot.OrderID != r.OrderID || snapshot.Currency != r.Currency || snapshot.AmountMinor != r.AmountMinor || snapshot.LiveMode != c.LiveMode || (c.ProviderCode == "mercadopago" && snapshot.CollectorID != c.AccountRef) {
				return ErrPaymentCallback
			}
		}
		if _, _, _, err = p.binding(ctx, attempt); err != nil {
			return err
		}
		var referenceMatches bool
		err = p.pool.QueryRow(ctx, `select (b.payment_reference is null or b.payment_reference=$3) and (o.provider_reference is null or o.provider_reference=$3) from payment.provider_checkout b left join payment.provider_observation o using(tenant_id,payment_attempt_id) where b.tenant_id=$1 and b.payment_attempt_id=$2`, c.TenantID, attempt, reference).Scan(&referenceMatches)
		if err != nil || !referenceMatches {
			return ErrPaymentCallback
		}
	}
	if _, _, err = p.notice(ctx, j); err != nil {
		return err
	}
	if err = p.worker.Reconcile(ctx, attempt, reference); err != nil {
		return err
	}
	// A crash before this receipt or the processor's subsequent ACK retries only
	// GETs. Existing observation/outbox owners suppress duplicate money effects.
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var owned bool
	err = tx.QueryRow(ctx, lockedJob+`select claimed_by=$3 and attempts=$4 and claimed_until>clock_timestamp() and payload=$5::jsonb from owned`, j.TenantID, j.JobID, p.workerID, j.Attempts, j.Payload).Scan(&owned)
	if err != nil || !owned {
		return ErrJobClaimLost
	}
	_, err = tx.Exec(ctx, `update integration.webhook_event set state='processed',processed_at=clock_timestamp(),last_error_code=null where tenant_id=$1 and connection_id=$2 and provider_event_id=$3 and provider_code=$4 and body_sha256_hex=$5`, c.TenantID, c.ConnectionID, payload.ProviderEventID, c.ProviderCode, n.BodySHA256)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (p *PaymentCallbackProcessor) recoverFence(ctx context.Context, r paymentbridge.Request, hash string, v paymentbridge.Checkout) error {
	var state string
	if err := p.pool.QueryRow(ctx, `select state from communication.outbound_delivery where tenant_id=$1 and channel_code=$2 and delivery_key=$3 and request_sha256_hex=$4`, r.TenantID, "payment_"+r.ProviderCode, r.PaymentAttemptID, hash).Scan(&state); err != nil {
		return err
	}
	if state == "accepted" {
		return nil
	}
	m, actual, err := paymentCallbackMessage(r)
	if err != nil || actual != hash {
		return ErrPaymentCallback
	}
	evidence, _ := json.Marshal(struct{ Provider, Session, Payment, Status string }{v.ProviderCode, v.SessionID, v.PaymentReference, v.Status})
	sum := sha256.Sum256(evidence)
	receipt := outbounddelivery.Receipt{ProviderMessageID: v.SessionID, EvidenceSHA256: hex.EncodeToString(sum[:]), AcceptedAt: time.Now().UTC()}
	if state == "sending" {
		return p.worker.Fence.Complete(ctx, m, hash, receipt)
	}
	reconciler, ok := p.worker.Fence.(interface {
		ReconcileAccepted(context.Context, channels.Message, string, outbounddelivery.Receipt) error
	})
	if state != "unknown" || !ok {
		return ErrPaymentCallback
	}
	return reconciler.ReconcileAccepted(ctx, m, hash, receipt)
}
