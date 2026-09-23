package whatsappbridge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
)

var ErrStatusObservation = errors.New("whatsappbridge: status observation unverified, unavailable or divergent")

type AppSecretSource interface {
	WhatsAppAppSecret(context.Context) (string, error)
}
type SignedStatusWebhook struct {
	Body      []byte `json:"body"`
	Signature string `json:"signature"`
}

// StatusObserver is a service integration, not a public HTTP endpoint. The
// caller supplies an actually verified Principal and durable inbox raw bytes.
// The database, never the caller, supplies the trusted outbound receipt anchor.
type StatusObserver struct {
	TenantID         string
	Profile          json.RawMessage
	Process          Process
	ReconcilerSHA256 string
	Secrets          AppSecretSource
	Approvals        *PostgresAppointmentApprovals
}

type statusEvent struct {
	Kind          string `json:"kind"`
	MessageHash   string `json:"id_sha256"`
	RecipientHash string `json:"recipient_sha256"`
	AccountHash   string `json:"business_account_id_sha256"`
	PhoneHash     string `json:"phone_number_id_sha256"`
	Timestamp     string `json:"timestamp"`
	Status        string `json:"status"`
	EventHash     string `json:"event_sha256"`
	EventKey      string `json:"event_key_sha256"`
}
type statusResult struct {
	Schema  string          `json:"schema"`
	Binding string          `json:"binding_sha256"`
	Receipt json.RawMessage `json:"receipt"`
	Events  []statusEvent   `json:"events"`
}
type statusReceipt struct {
	Schema           string `json:"schema"`
	Anchor           string `json:"send_receipt_sha256"`
	ObservationsHash string `json:"observations_sha256"`
	Matching         int    `json:"matching_unique_events"`
	BusinessWrite    bool   `json:"automatic_business_write"`
	Resend           bool   `json:"resend_authorized"`
}

type statusQuery interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func (o *StatusObserver) anchor(ctx context.Context, db statusQuery, p identity.Principal, organization, appointment, event string, locked bool) (string, string, error) {
	suffix := ""
	if locked {
		suffix = " for update of d for share of l"
	}
	var anchor, key string
	err := db.QueryRow(ctx, `select d.evidence_sha256_hex,g.delivery_key
 from communication.whatsapp_delivery_approval g
 join crm.lead l on l.tenant_id=g.tenant_id and l.lead_id=g.lead_id and l.organization_id=g.organization_id
 join communication.outbound_delivery d on d.tenant_id=g.tenant_id and d.delivery_key=g.delivery_key and d.channel_code=g.channel_code
 where g.tenant_id=$1 and g.organization_id=$2 and g.appointment_id=$3 and g.approval_event_id=$4
 and g.profile_sha256=$5 and d.state='accepted' and d.request_sha256_hex=g.message_sha256
 and d.recipient_hmac=g.external_id_hmac and d.accepted_at is not null`+suffix,
		p.TenantID, organization, appointment, event, digest(o.Profile)).Scan(&anchor, &key)
	if err != nil || !validDigest(anchor) {
		return "", "", ErrStatusObservation
	}
	return anchor, key, nil
}

// Observe validates all signatures in an isolated, hash-locked Python process,
// then stores observations atomically. It never changes the send fence, creates
// a resend, or acknowledges a durable inbox before the caller sees success.
func (o *StatusObserver) Observe(ctx context.Context, p identity.Principal, organization, appointment, event string, receipt []byte, batches []SignedStatusWebhook) (int, error) {
	if o == nil || o.Approvals == nil || o.Approvals.pool == nil || o.Secrets == nil || p.TenantID == "" || p.TenantID != o.TenantID || p.Subject == "" || len(p.Subject) > 255 || !p.Allowed("appointment:manage") || !p.AllowedOrganization(organization) || !((appointment != "" && notificationEventID.MatchString(event)) || (appointment == "" && (strings.HasPrefix(event, "wa-reply:") && validDigest(strings.TrimPrefix(event, "wa-reply:")) || strings.HasPrefix(event, "wa-schedule:") && validDigest(strings.TrimPrefix(event, "wa-schedule:"))))) || len(o.Profile) > 32768 || !json.Valid(o.Profile) || len(receipt) == 0 || len(receipt) > 65536 || len(batches) == 0 || len(batches) > 8 {
		return 0, ErrStatusObservation
	}
	for _, b := range batches {
		if len(b.Body) == 0 || len(b.Body) > 1<<20 || len(b.Signature) != 71 {
			return 0, ErrStatusObservation
		}
	}
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	anchor, key, err := o.anchor(ctx, o.Approvals.pool, p, organization, appointment, event, false)
	if err != nil || digest(receipt) != anchor {
		return 0, ErrStatusObservation
	}
	proc := o.Process
	script := filepath.Join(proc.AdapterDirectory, "whatsapp_cloud.py")
	if !filepath.IsAbs(proc.AdapterDirectory) || !filepath.IsAbs(proc.EvidenceDirectory) || exactFile(proc.PythonExecutable, proc.PythonSHA256) != nil || exactFile(script, proc.AdapterSHA256) != nil || exactFile(filepath.Join(proc.AdapterDirectory, "status_reconciliation.py"), o.ReconcilerSHA256) != nil {
		return 0, ErrStatusObservation
	}
	secret, err := o.Secrets.WhatsAppAppSecret(ctx)
	if err != nil || secret == "" || len(secret) > 16384 {
		return 0, ErrStatusObservation
	}
	input, err := json.Marshal(struct {
		Schema    string                `json:"schema"`
		Profile   json.RawMessage       `json:"profile"`
		Receipt   []byte                `json:"send_receipt"`
		Anchor    string                `json:"expected_send_receipt_sha256"`
		Webhooks  []SignedStatusWebhook `json:"webhooks"`
		Secret    string                `json:"app_secret"`
		Directory string                `json:"evidence_directory"`
	}{"elite-whatsapp-status-bridge/v1", o.Profile, receipt, anchor, batches, secret, proc.EvidenceDirectory})
	if err != nil || len(input) > 12<<20 {
		return 0, ErrStatusObservation
	}
	cmd := exec.CommandContext(ctx, proc.PythonExecutable, "-I", "-B", script, "--status-bridge")
	cmd.Env = []string{}
	if root := os.Getenv("SystemRoot"); root != "" {
		cmd.Env = append(cmd.Env, "SystemRoot="+root)
	}
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stderr = io.Discard
	cmd.WaitDelay = 2 * time.Second
	stdout := statusOutput{}
	cmd.Stdout = &stdout
	if cmd.Run() != nil {
		return 0, ErrStatusObservation
	}
	var result statusResult
	decoder := json.NewDecoder(bytes.NewReader(stdout.Bytes()))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&result) != nil || decoder.Decode(new(any)) != io.EOF || result.Schema != "elite-whatsapp-status-result/v1" || result.Binding != digest(input) || len(result.Events) > 8000 {
		return 0, ErrStatusObservation
	}
	var evidence statusReceipt
	if len(result.Receipt) > 65536 || json.Unmarshal(result.Receipt, &evidence) != nil || evidence.Schema != "elite-whatsapp-status-reconciliation/v1" || evidence.Anchor != anchor || !validDigest(evidence.ObservationsHash) || evidence.Matching != len(result.Events) || evidence.BusinessWrite || evidence.Resend {
		return 0, ErrStatusObservation
	}
	for _, e := range result.Events {
		stamp, err := strconv.ParseInt(e.Timestamp, 10, 64)
		if err != nil || stamp < 1 || stamp > 999999999999 || strconv.FormatInt(stamp, 10) != e.Timestamp || !validDigest(e.EventKey) || !validDigest(e.EventHash) || e.Kind != "status" {
			return 0, ErrStatusObservation
		}
		switch e.Status {
		case "sent", "delivered", "read", "failed", "deleted":
		default:
			return 0, ErrStatusObservation
		}
	}
	tx, err := o.Approvals.pool.Begin(ctx)
	if err != nil {
		return 0, ErrStatusObservation
	}
	defer tx.Rollback(ctx)
	current, currentKey, err := o.anchor(ctx, tx, p, organization, appointment, event, true)
	if err != nil || current != anchor || currentKey != key {
		return 0, ErrStatusObservation
	}
	if len(result.Events) == 0 {
		return 0, nil
	}
	_, err = tx.Exec(ctx, `insert into communication.whatsapp_status_batch(tenant_id,delivery_key,observations_sha256,send_receipt_sha256,verification_receipt,observed_by) values($1,$2,$3,$4,$5,$6) on conflict do nothing`, p.TenantID, key, evidence.ObservationsHash, anchor, result.Receipt, p.Subject)
	if err != nil {
		return 0, ErrStatusObservation
	}
	inserted := 0
	for _, e := range result.Events {
		timestamp, _ := strconv.ParseInt(e.Timestamp, 10, 64)
		tag, err := tx.Exec(ctx, `insert into communication.whatsapp_status_observation(tenant_id,delivery_key,event_key_sha256,event_sha256,provider_timestamp,provider_status,observations_sha256) values($1,$2,$3,$4,$5,$6,$7) on conflict do nothing`, p.TenantID, key, e.EventKey, e.EventHash, timestamp, e.Status, evidence.ObservationsHash)
		if err != nil {
			return 0, ErrStatusObservation
		}
		if tag.RowsAffected() == 1 {
			inserted++
			continue
		}
		var hash, state string
		var storedTime int64
		if tx.QueryRow(ctx, `select event_sha256,provider_status,provider_timestamp from communication.whatsapp_status_observation where tenant_id=$1 and delivery_key=$2 and event_key_sha256=$3`, p.TenantID, key, e.EventKey).Scan(&hash, &state, &storedTime) != nil || hash != e.EventHash || state != e.Status || storedTime != timestamp {
			return 0, ErrStatusObservation
		}
	}
	if tx.Commit(ctx) != nil {
		return 0, ErrStatusObservation
	}
	return inserted, nil
}

type statusOutput struct{ bytes.Buffer }

func (b *statusOutput) Write(p []byte) (int, error) {
	if b.Len()+len(p) > 8<<20 {
		return 0, ErrStatusObservation
	}
	return b.Buffer.Write(p)
}
