package whatsappbridge

// AUTHORED composition of the existing Meta-adapted verifier and SQL owners.
// Optional verified customer-message projection uses the existing runtime.
// No automatic reply sending, new queue, or principal fabrication.
import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
)

var ErrStatusRouting = errors.New("whatsappbridge: retained status routing unavailable or unverified")

type StatusRouter struct {
	observer              StatusObserver
	connection, retention string
	key                   []byte
	maxRoutes             int
	conversation          *ConversationRoute
	busy                  chan struct{}
}

// key must be the outbound owner's exact HMAC key (resolved server-side).
// A real verified principal is required per call; this is not a public handler.
func NewStatusRouter(o *StatusObserver, connection, retention string, key []byte, maxRoutes int) (*StatusRouter, error) {
	if o == nil || o.Approvals == nil || o.Approvals.pool == nil || o.Approvals.pool.Config().MaxConns < 2 || o.Secrets == nil || !notificationEventID.MatchString(o.TenantID) || !webhookConnection.MatchString(connection) || !validDigest(retention) || len(key) < 32 || maxRoutes < 1 || maxRoutes > 8 || len(o.Profile) > 32768 || !json.Valid(o.Profile) || !filepath.IsAbs(o.Process.EvidenceDirectory) {
		return nil, ErrStatusRouting
	}
	copyObserver := *o
	copyObserver.Profile = append(json.RawMessage(nil), o.Profile...)
	return &StatusRouter{observer: copyObserver, connection: connection, retention: retention, key: append([]byte(nil), key...), maxRoutes: maxRoutes, busy: make(chan struct{}, 1)}, nil
}

type retainedStatus struct {
	Schema        string `json:"schema"`
	Body          []byte `json:"body"`
	Signature     string `json:"signature"`
	ProfileHash   string `json:"profile_sha256"`
	RetentionHash string `json:"retention_approval_sha256"`
	Count         int    `json:"verified_event_count"`
}

type statusIdentity struct{ Message, Recipient string }
type resolvedStatus struct {
	organization, appointment, event, delivery, anchor string
	identity                                           statusIdentity
	receipt                                            []byte
}

// ObserveRetained verifies all routes before invoking the existing owners.
// It never marks inbox/job complete. A caller must retain/retry/reconcile on
// error; earlier observer commits may exist and are idempotently replayable.
func (r *StatusRouter) ObserveRetained(ctx context.Context, p identity.Principal, eventID string) (int, error) {
	return r.observeRetained(ctx, p, eventID, nil)
}

// finish is internal-only: it may atomically ACK the exact claimed job/inbox
// after every admitted route was observed. Public routing still performs no ACK.
func (r *StatusRouter) observeRetained(ctx context.Context, p identity.Principal, eventID string, finish func(context.Context, pgx.Tx) error) (int, error) {
	if r == nil || p.TenantID != r.observer.TenantID || p.Subject == "" || len(p.Subject) > 255 || !p.Allowed("appointment:manage") || !strings.HasPrefix(eventID, "wa:") || !validDigest(strings.TrimPrefix(eventID, "wa:")) {
		return 0, ErrStatusRouting
	}
	select {
	case r.busy <- struct{}{}:
		defer func() { <-r.busy }()
	default:
		return 0, ErrStatusRouting
	}
	timeout := 30 * time.Second
	if r.conversation != nil {
		timeout = 110 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	o := &r.observer
	tx, err := o.Approvals.pool.Begin(ctx)
	if err != nil {
		return 0, ErrStatusRouting
	}
	defer tx.Rollback(ctx)
	// Hold connection/inbox scope stable through verification and observation.
	// Observer needs another pool connection; this router admits one call at a time.
	var payload []byte
	var bodyHash string
	var connectionOrg *string
	err = tx.QueryRow(ctx, `select e.payload,e.body_sha256_hex,c.organization_id
 from integration.webhook_event e join integration.provider_connection c
 on c.tenant_id=e.tenant_id and c.connection_id=e.connection_id and c.provider_code=e.provider_code
 where e.tenant_id=$1 and e.connection_id=$2 and e.provider_event_id=$3
 and c.state='active' and e.provider_code='meta-whatsapp'
 and e.event_type='whatsapp.raw_webhook.received.v1' and e.state in ('received','processed')
 and octet_length(e.payload::text)<=2097152 for share of c,e`, p.TenantID, r.connection, eventID).Scan(&payload, &bodyHash, &connectionOrg)
	if err != nil {
		return 0, ErrStatusRouting
	}
	var retained retainedStatus
	d := json.NewDecoder(bytes.NewReader(payload))
	d.DisallowUnknownFields()
	if d.Decode(&retained) != nil || d.Decode(new(any)) != io.EOF || retained.Schema != "elite-whatsapp-retained-webhook/v1" || len(retained.Body) == 0 || len(retained.Body) > 1<<20 || retained.ProfileHash != digest(o.Profile) || retained.RetentionHash != r.retention || retained.Count < 1 || retained.Count > 1000 || digest(retained.Body) != bodyHash || eventID != "wa:"+bodyHash {
		return 0, ErrStatusRouting
	}
	secret, err := o.Secrets.WhatsAppAppSecret(ctx)
	if err != nil {
		return 0, ErrStatusRouting
	}
	verifier := &WebhookReceiver{config: WebhookReceiverConfig{Profile: o.Profile, Process: o.Process}}
	verified, err := verifier.verify(ctx, "receive", retained.Body, retained.Signature, secret, map[string]string{})
	if err != nil || verified.BodyHash != bodyHash || verified.Count != retained.Count || verified.Challenge != "" {
		return 0, ErrStatusRouting
	}
	// This is identity projection AFTER the shared verifier has checked raw HMAC,
	// unique JSON keys, complete envelope, WABA/phone, events, and status semantics.
	// It is not a replacement signature/scope/normalization implementation.
	messages, messageCount, err := verifiedConversationMessages(retained.Body, p.TenantID, r.connection, r.maxRoutes)
	if err != nil || (messageCount > 0 && (r.conversation == nil || connectionOrg == nil || *connectionOrg != r.conversation.OrganizationID || !p.AllowedOrganization(*connectionOrg))) {
		return 0, ErrStatusRouting
	}
	identities, err := verifiedStatusIdentities(retained.Body, verified.Count, r.maxRoutes)
	if err != nil {
		return 0, ErrStatusRouting
	}
	root, err := os.OpenRoot(o.Process.EvidenceDirectory)
	if err != nil {
		return 0, ErrStatusRouting
	}
	defer root.Close()
	routes := make([]resolvedStatus, 0, len(identities))
	for _, id := range identities {
		messageHMAC, e1 := contactidentity.ExternalDigest(r.key, p.TenantID, "whatsapp", id.Message)
		recipientHMAC, e2 := contactidentity.ExternalDigest(r.key, p.TenantID, "whatsapp", id.Recipient)
		if e1 != nil || e2 != nil {
			return 0, ErrStatusRouting
		}
		rows, err := tx.Query(ctx, `select g.organization_id,g.appointment_id,g.approval_event_id,g.delivery_key,d.evidence_sha256_hex
 from communication.outbound_delivery d join communication.whatsapp_delivery_approval g
 on g.tenant_id=d.tenant_id and g.channel_code=d.channel_code and g.delivery_key=d.delivery_key
 join crm.lead l on l.tenant_id=g.tenant_id and l.lead_id=g.lead_id and l.organization_id=g.organization_id
 where d.tenant_id=$1 and d.channel_code='whatsapp' and d.state='accepted' and d.accepted_at is not null
 and d.provider_message_hmac=$2 and d.recipient_hmac=$3 and g.external_id_hmac=d.recipient_hmac
 and g.profile_sha256=$4 and d.request_sha256_hex=g.message_sha256
 and ($5::text is null or g.organization_id=$5) limit 2`, p.TenantID, messageHMAC, recipientHMAC, digest(o.Profile), connectionOrg)
		if err != nil {
			return 0, ErrStatusRouting
		}
		matches := 0
		var route resolvedStatus
		for rows.Next() {
			matches++
			if rows.Scan(&route.organization, &route.appointment, &route.event, &route.delivery, &route.anchor) != nil {
				rows.Close()
				return 0, ErrStatusRouting
			}
		}
		queryErr := rows.Err()
		rows.Close()
		if queryErr != nil || matches != 1 || !p.AllowedOrganization(route.organization) || !validDigest(route.anchor) || !validApprovalRoute(p.TenantID, route) {
			return 0, ErrStatusRouting
		}
		route.identity = id
		route.receipt, err = routedReceipt(root, p.TenantID, route)
		if err != nil {
			return 0, ErrStatusRouting
		}
		routes = append(routes, route)
	}
	// All routes and message shapes are validated before effects. Runtime replay
	// handles a crash after proposal commit and before the shared job ACK.
	for _, message := range messages {
		if _, err := r.conversation.Runtime.Handle(ctx, message); err != nil {
			return 0, ErrStatusRouting
		}
		if _, err := r.conversation.Proposals.Propose(ctx, p, eventID, message); err != nil {
			return 0, ErrStatusRouting
		}
	}
	inserted := 0
	for _, route := range routes {
		n, err := o.Observe(ctx, p, route.organization, route.appointment, route.event, route.receipt, []SignedStatusWebhook{{Body: retained.Body, Signature: retained.Signature}})
		inserted += n
		if err != nil {
			return inserted, ErrStatusRouting
		}
	}
	if finish != nil {
		if err := finish(ctx, tx); err != nil {
			return inserted, ErrStatusRouting
		}
	}
	if tx.Commit(ctx) != nil {
		return inserted, ErrStatusRouting
	}
	return inserted, nil
}

func verifiedStatusIdentities(raw []byte, count, limit int) ([]statusIdentity, error) {
	// Exact map keys, not struct decoding: encoding/json matches struct fields
	// case-insensitively, unlike the already-verified Python provider parser.
	var envelope map[string]json.RawMessage
	if json.Unmarshal(raw, &envelope) != nil {
		return nil, ErrStatusRouting
	}
	var entries []map[string]json.RawMessage
	if json.Unmarshal(envelope["entry"], &entries) != nil {
		return nil, ErrStatusRouting
	}
	seen := map[statusIdentity]bool{}
	total := 0
	for _, entry := range entries {
		var changes []map[string]json.RawMessage
		if json.Unmarshal(entry["changes"], &changes) != nil {
			return nil, ErrStatusRouting
		}
		for _, change := range changes {
			var value map[string]json.RawMessage
			if json.Unmarshal(change["value"], &value) != nil {
				return nil, ErrStatusRouting
			}
			var messages []json.RawMessage
			if v, ok := value["messages"]; ok {
				if json.Unmarshal(v, &messages) != nil {
					return nil, ErrStatusRouting
				}
			}
			total += len(messages)
			var statuses []map[string]json.RawMessage
			if v, ok := value["statuses"]; ok && json.Unmarshal(v, &statuses) != nil {
				return nil, ErrStatusRouting
			}
			for _, status := range statuses {
				var message, recipient string
				if json.Unmarshal(status["id"], &message) != nil || json.Unmarshal(status["recipient_id"], &recipient) != nil {
					return nil, ErrStatusRouting
				}
				total++
				if message == "" || recipient == "" || len(message) > 256 || len(recipient) > 256 {
					return nil, ErrStatusRouting
				}
				seen[statusIdentity{message, recipient}] = true
				if len(seen) > limit {
					return nil, ErrStatusRouting
				}
			}
		}
	}
	if total != count || total == 0 {
		return nil, ErrStatusRouting
	}
	ids := make([]statusIdentity, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		if ids[i].Message == ids[j].Message {
			return ids[i].Recipient < ids[j].Recipient
		}
		return ids[i].Message < ids[j].Message
	})
	return ids, nil
}

func routedReceipt(root *os.Root, tenant string, route resolvedStatus) ([]byte, error) {
	// Same deterministic location as Sender, but root-confined and bounded.
	file, err := root.Open(filepath.Join(digest([]byte(tenant+"\x00"+route.delivery)), "SEND_RECEIPT.json"))
	if err != nil {
		return nil, ErrStatusRouting
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() < 1 || info.Size() > 65536 {
		return nil, ErrStatusRouting
	}
	raw, err := io.ReadAll(io.LimitReader(file, 65537))
	if err != nil || len(raw) > 65536 || digest(raw) != route.anchor {
		return nil, ErrStatusRouting
	}
	var receipt struct {
		Message   string `json:"message_id_sha256"`
		Recipient string `json:"recipient_sha256"`
	}
	if json.Unmarshal(raw, &receipt) != nil || receipt.Message != digest([]byte(route.identity.Message)) || receipt.Recipient != digest([]byte(route.identity.Recipient)) {
		return nil, ErrStatusRouting
	}
	return raw, nil
}

func validApprovalRoute(tenant string, r resolvedStatus) bool {
	if r.appointment == "" {
		return r.delivery == r.event && strings.HasPrefix(r.event, "wa-reply:") && validDigest(strings.TrimPrefix(r.event, "wa-reply:"))
	}
	return r.delivery == AppointmentConfirmationDeliveryKey(tenant, r.event)
}
