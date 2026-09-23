package whatsappbridge

// AUTHORED extraction of the existing exact local provider receipt verifier.
// Callers must first authorize and bind the immutable reviewed context.
import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type boundReceiptContext struct {
	Message       channels.Message
	MessageSHA256 string
	ExpiresAt     time.Time
}

func recoverBoundProviderReceipt(ctx context.Context, sender *Sender, store *postgres.OutboundDeliveryStore, pool *pgxpool.Pool, tenant, key string, c boundReceiptContext, approvalAt time.Time, fenceState string) (outbounddelivery.Receipt, error) {
	request := json.RawMessage(c.Message.Text)
	actual, err := outbounddelivery.MessageSHA256(c.Message)
	if err != nil || actual != c.MessageSHA256 {
		return outbounddelivery.Receipt{}, ErrApproval
	}
	proc := sender.Process
	script := filepath.Join(proc.AdapterDirectory, "whatsapp_cloud.py")
	if !filepath.IsAbs(proc.EvidenceDirectory) || exactFile(proc.PythonExecutable, proc.PythonSHA256) != nil || exactFile(script, proc.AdapterSHA256) != nil {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	root, err := os.OpenRoot(proc.EvidenceDirectory)
	if err != nil {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	defer root.Close()
	dir := digest([]byte(tenant + "\x00" + key))
	receiptRaw, err := recoverEvidence(root, filepath.Join(dir, "SEND_RECEIPT.json"))
	if err != nil {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	responseRaw, err := recoverEvidence(root, filepath.Join(dir, "provider-response.json"))
	if err != nil {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	raw, err := json.Marshal(struct {
		Schema   string          `json:"schema"`
		Binding  string          `json:"binding_sha256"`
		Profile  json.RawMessage `json:"profile"`
		Request  json.RawMessage `json:"request"`
		Receipt  []byte          `json:"receipt"`
		Response []byte          `json:"response"`
	}{"elite-whatsapp-recover-send/v1", actual, sender.Profile, request, receiptRaw, responseRaw})
	if err != nil {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, proc.PythonExecutable, "-I", "-B", script, "--recover-send-bridge")
	cmd.Stdin = bytes.NewReader(raw)
	cmd.Stderr = io.Discard
	cmd.Env = []string{}
	cmd.WaitDelay = 2 * time.Second
	if root := os.Getenv("SystemRoot"); root != "" {
		cmd.Env = append(cmd.Env, "SystemRoot="+root)
	}
	var stdout boundedOutput
	cmd.Stdout = &stdout
	if cmd.Run() != nil || stdout.overflow {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	var result result
	d := json.NewDecoder(bytes.NewReader(stdout.Bytes()))
	d.DisallowUnknownFields()
	if d.Decode(&result) != nil || d.Decode(new(any)) != io.EOF || result.Schema != "elite-whatsapp-send-result/v1" || result.BindingSHA256 != actual || result.EvidenceSHA256 != digest(receiptRaw) {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	receipt := outbounddelivery.Receipt{ProviderMessageID: result.ProviderMessageID, EvidenceSHA256: result.EvidenceSHA256, AcceptedAt: result.AcceptedAt}
	if receipt.Validate() != nil || receipt.AcceptedAt.Before(approvalAt) || !receipt.AcceptedAt.Before(c.ExpiresAt) {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	// Use the same owner's expired-lease transition; it never calls the provider.
	// The existing row was read above, and this path cannot create a new delivery.
	if fenceState == "sending" {
		claim, claimErr := store.Claim(ctx, c.Message, actual)
		if claim.Replay {
			fenceState = "accepted"
		} else if !errors.Is(claimErr, outbounddelivery.ErrUnknown) {
			return outbounddelivery.Receipt{}, ErrBridge
		}
	}
	if fenceState == "accepted" {
		var same bool
		err = pool.QueryRow(ctx, `select evidence_sha256_hex=$3 and request_sha256_hex=$4 from communication.outbound_delivery where tenant_id=$1 and channel_code='whatsapp' and delivery_key=$2 and state='accepted'`, tenant, key, receipt.EvidenceSHA256, actual).Scan(&same)
		if err != nil || !same {
			return outbounddelivery.Receipt{}, ErrBridge
		}
		return receipt, nil
	}
	if err = store.ReconcileAccepted(ctx, c.Message, actual, receipt); err != nil {
		return outbounddelivery.Receipt{}, err
	}
	return receipt, nil
}
