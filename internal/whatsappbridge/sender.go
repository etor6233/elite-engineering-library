// Package whatsappbridge connects the admitted Python adapter to the existing
// outbound delivery fence. This is AUTHORED integration, not Meta source code.
package whatsappbridge

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
)

var ErrBridge = errors.New("whatsappbridge: unverified or uncertain provider result")

// Approval must be resolved from an authorized durable workflow, not message
// text. The bridge never generates consent, templates or business approvals.
type Approval struct {
	MessageSHA256  string
	ProfileSHA256  string
	EvidenceSHA256 string
	ApprovedBy     string
	NotBefore      time.Time
	ExpiresAt      time.Time
}

type ApprovalResolver interface {
	ResolveWhatsAppApproval(context.Context, string, string) (Approval, error)
}

type TokenSource interface {
	WhatsAppToken(context.Context) (string, error)
}

type Process struct {
	PythonExecutable  string
	PythonSHA256      string
	AdapterDirectory  string
	AdapterSHA256     string
	EvidenceDirectory string
}

// Sender is tenant-bound at composition time. Always wrap it in
// outbounddelivery.Channel with the shared PostgreSQL Store; never invoke a
// provider retry from an error. Accepted receipts do not assert delivery.
type Sender struct {
	TenantID  string
	Profile   json.RawMessage
	Approvals ApprovalResolver
	Tokens    TokenSource
	Process   Process
}

type frame struct {
	Schema          string          `json:"schema"`
	BindingSHA256   string          `json:"binding_sha256"`
	Profile         json.RawMessage `json:"profile"`
	Request         json.RawMessage `json:"request"`
	AccessToken     string          `json:"access_token"`
	OutputDirectory string          `json:"output_directory"`
}

type result struct {
	Schema            string    `json:"schema"`
	BindingSHA256     string    `json:"binding_sha256"`
	ProviderMessageID string    `json:"provider_message_id"`
	EvidenceSHA256    string    `json:"evidence_sha256"`
	AcceptedAt        time.Time `json:"accepted_at"`
}

func digest(value []byte) string { sum := sha256.Sum256(value); return hex.EncodeToString(sum[:]) }
func validDigest(value string) bool {
	b, e := hex.DecodeString(value)
	return e == nil && len(b) == 32 && value == strings.ToLower(value)
}

func exactFile(path, hash string) error {
	if !filepath.IsAbs(path) || !validDigest(hash) {
		return ErrBridge
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return ErrBridge
	}
	file, err := os.Open(path)
	if err != nil {
		return ErrBridge
	}
	defer file.Close()
	h := sha256.New()
	if _, err = io.Copy(h, file); err != nil || hex.EncodeToString(h.Sum(nil)) != hash {
		return ErrBridge
	}
	return nil
}

// Decode only the exact template shape; the Python owner validates policy and
// arity. The recipient must be the same one fenced in channels.Message.
func requestFor(message channels.Message) (json.RawMessage, error) {
	if len(message.Text) > 32768 {
		return nil, ErrBridge
	}
	var shape map[string]json.RawMessage
	if json.Unmarshal([]byte(message.Text), &shape) == nil {
		if _, ok := shape["kind"]; ok {
			return replyRequestFor(message)
		}
	}
	var body struct {
		Recipient      string   `json:"recipient"`
		TemplateName   string   `json:"template_name"`
		LanguageCode   string   `json:"language_code"`
		BodyParameters []string `json:"body_parameters"`
	}
	decoder := json.NewDecoder(strings.NewReader(message.Text))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&body) != nil || decoder.Decode(new(any)) != io.EOF || body.Recipient != message.ExternalID || body.TemplateName == "" || body.LanguageCode == "" || body.BodyParameters == nil {
		return nil, ErrBridge
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, ErrBridge
	}
	return raw, nil
}

func (s *Sender) SendWithReceipt(ctx context.Context, message channels.Message) (outbounddelivery.Receipt, error) {
	empty := outbounddelivery.Receipt{}
	if s == nil || s.TenantID == "" || message.TenantID != s.TenantID || message.ChannelCode != "whatsapp" || message.Direction != channels.DirectionOut || s.Approvals == nil || s.Tokens == nil || len(s.Profile) > 32768 || !json.Valid(s.Profile) {
		return empty, ErrBridge
	}
	binding, err := outbounddelivery.MessageSHA256(message)
	if err != nil {
		return empty, ErrBridge
	}
	request, err := requestFor(message)
	if err != nil {
		return empty, err
	}
	approval, err := s.Approvals.ResolveWhatsAppApproval(ctx, message.TenantID, message.DeliveryKey)
	now := time.Now().UTC()
	if err != nil || approval.MessageSHA256 != binding || approval.ProfileSHA256 != digest(s.Profile) || !validDigest(approval.EvidenceSHA256) || strings.TrimSpace(approval.ApprovedBy) == "" || approval.NotBefore.IsZero() || now.Before(approval.NotBefore) || !now.Before(approval.ExpiresAt) {
		return empty, ErrBridge
	}
	p := s.Process
	script := filepath.Join(p.AdapterDirectory, "whatsapp_cloud.py")
	if !filepath.IsAbs(p.AdapterDirectory) || !filepath.IsAbs(p.EvidenceDirectory) || exactFile(p.PythonExecutable, p.PythonSHA256) != nil || exactFile(script, p.AdapterSHA256) != nil {
		return empty, ErrBridge
	}
	token, err := s.Tokens.WhatsAppToken(ctx)
	if err != nil || token == "" || len(token) > 16384 {
		return empty, ErrBridge
	}
	// Delivery identity determines one evidence destination. It does not replace
	// the database fence and must never be changed merely to force another send.
	output := filepath.Join(p.EvidenceDirectory, digest([]byte(message.TenantID+"\x00"+message.DeliveryKey)))
	input, err := json.Marshal(frame{"elite-whatsapp-send-bridge/v1", binding, s.Profile, request, token, output})
	if err != nil {
		return empty, ErrBridge
	}
	limitedCtx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()
	command := exec.CommandContext(limitedCtx, p.PythonExecutable, "-I", "-B", script, "--send-bridge")
	// No shell and no inherited API tokens, PYTHONPATH, proxies or startup hooks.
	command.Env = []string{}
	if root := os.Getenv("SystemRoot"); root != "" {
		command.Env = append(command.Env, "SystemRoot="+root)
	}
	command.Stdin = bytes.NewReader(input)
	var stdout boundedOutput
	command.Stdout = &stdout
	command.Stderr = io.Discard
	command.WaitDelay = 2 * time.Second
	if command.Run() != nil || stdout.overflow {
		return empty, ErrBridge
	}
	var response result
	decoder := json.NewDecoder(bytes.NewReader(stdout.Bytes()))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&response) != nil || decoder.Decode(new(any)) != io.EOF || response.Schema != "elite-whatsapp-send-result/v1" || response.BindingSHA256 != binding {
		return empty, ErrBridge
	}
	receipt := outbounddelivery.Receipt{ProviderMessageID: response.ProviderMessageID, EvidenceSHA256: response.EvidenceSHA256, AcceptedAt: response.AcceptedAt}
	if receipt.Validate() != nil || receipt.AcceptedAt.Before(now.Add(-time.Minute)) || receipt.AcceptedAt.After(time.Now().Add(time.Minute)) {
		return empty, ErrBridge
	}
	return receipt, nil
}

type boundedOutput struct {
	bytes.Buffer
	overflow bool
}

func (b *boundedOutput) Write(p []byte) (int, error) {
	if b.Len()+len(p) > 4096 {
		b.overflow = true
		return 0, ErrBridge
	}
	return b.Buffer.Write(p)
}
