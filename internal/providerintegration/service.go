package providerintegration

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalid         = errors.New("invalid webhook")
	ErrUnauthenticated = errors.New("webhook authentication failed")
	ErrConnection      = errors.New("provider connection not found")
	ErrConflict        = errors.New("provider event conflict")
)

type Connection struct {
	TenantID, ConnectionID, ProviderCode, SecretRef string
	Secret                                          []byte
}

type Registry interface {
	Lookup(providerCode, connectionID string) (Connection, bool)
}

type Receipt struct {
	TenantID, ConnectionID, ProviderCode string
	ProviderEventID, EventType, BodyHash string
	Payload                              json.RawMessage
}

type Repository interface {
	AcceptWebhook(context.Context, Receipt) (replayed bool, err error)
}

type SignatureVerifier interface {
	Verify(timestamp, signature string, body, secret []byte) error
}

type Service struct {
	repository Repository
	registry   Registry
	verifier   SignatureVerifier
}

func NewService(repository Repository, registry Registry, verifier SignatureVerifier) *Service {
	return &Service{repository: repository, registry: registry, verifier: verifier}
}

func (s *Service) Receive(ctx context.Context, providerCode, connectionID, timestamp, signature string, body []byte) (bool, error) {
	if s == nil || s.repository == nil || s.registry == nil || s.verifier == nil || len(body) == 0 {
		return false, ErrInvalid
	}
	connection, ok := s.registry.Lookup(providerCode, connectionID)
	if !ok {
		return false, ErrConnection
	}
	if err := s.verifier.Verify(timestamp, signature, body, connection.Secret); err != nil {
		return false, err
	}
	var envelope struct {
		EventID string          `json:"event_id"`
		Type    string          `json:"type"`
		Data    json.RawMessage `json:"data"`
	}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&envelope); err != nil {
		return false, ErrInvalid
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) || !bounded(envelope.EventID, 200) || !bounded(envelope.Type, 200) || !jsonObject(envelope.Data) {
		return false, ErrInvalid
	}
	hash := sha256.Sum256(body)
	return s.repository.AcceptWebhook(ctx, Receipt{TenantID: connection.TenantID, ConnectionID: connection.ConnectionID, ProviderCode: connection.ProviderCode, ProviderEventID: envelope.EventID, EventType: envelope.Type, BodyHash: hex.EncodeToString(hash[:]), Payload: envelope.Data})
}

func bounded(value string, limit int) bool { return value != "" && len(value) <= limit }
func jsonObject(value json.RawMessage) bool {
	var object map[string]json.RawMessage
	return len(value) > 0 && json.Unmarshal(value, &object) == nil && object != nil
}

type HMACSHA256 struct {
	Now       func() time.Time
	Tolerance time.Duration
}

func (v HMACSHA256) Verify(timestamp, signature string, body, secret []byte) error {
	if len(secret) < 32 || v.Now == nil || v.Tolerance <= 0 || v.Tolerance > 15*time.Minute {
		return ErrUnauthenticated
	}
	seconds, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return ErrUnauthenticated
	}
	delta := v.Now().Sub(time.Unix(seconds, 0))
	if delta < -v.Tolerance || delta > v.Tolerance {
		return ErrUnauthenticated
	}
	if !strings.HasPrefix(signature, "v1=") || strings.Contains(signature[3:], ",") {
		return ErrUnauthenticated
	}
	provided, err := hex.DecodeString(signature[3:])
	if err != nil || len(provided) != sha256.Size {
		return ErrUnauthenticated
	}
	mac := hmac.New(sha256.New, secret)
	_, _ = fmt.Fprintf(mac, "%s.", timestamp)
	_, _ = mac.Write(body)
	if !hmac.Equal(provided, mac.Sum(nil)) {
		return ErrUnauthenticated
	}
	return nil
}
