package providerintegration

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"
)

type receiptRepository struct {
	receipt Receipt
	replay  bool
}

func (r *receiptRepository) AcceptWebhook(_ context.Context, value Receipt) (bool, error) {
	r.receipt = value
	return r.replay, nil
}

func signed(secret []byte, timestamp string, body []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(timestamp + "."))
	_, _ = mac.Write(body)
	return "v1=" + hex.EncodeToString(mac.Sum(nil))
}

func TestReceiveAuthenticatesAndNormalizes(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	secret := []byte("01234567890123456789012345678901")
	registry, err := LoadHMACConnections([]byte(`[{"tenant_id":"tenant","connection_id":"primary","provider_code":"sandbox","secret_env":"SANDBOX_SECRET"}]`), func(string) (string, bool) { return string(secret), true })
	if err != nil {
		t.Fatal(err)
	}
	repository := &receiptRepository{}
	service := NewService(repository, registry, HMACSHA256{Now: func() time.Time { return now }, Tolerance: 5 * time.Minute})
	body := []byte(`{"event_id":"evt-1","type":"payment.updated","data":{"state":"paid"}}`)
	replayed, err := service.Receive(context.Background(), "sandbox", "primary", "1800000000", signed(secret, "1800000000", body), body)
	if err != nil || replayed || repository.receipt.ProviderEventID != "evt-1" || repository.receipt.TenantID != "tenant" || len(repository.receipt.BodyHash) != 64 {
		t.Fatalf("replayed=%v receipt=%+v err=%v", replayed, repository.receipt, err)
	}
}

func TestReceiveRejectsTamperReplayWindowAndShape(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	secret := []byte("01234567890123456789012345678901")
	registry := &StaticRegistry{connections: map[string]Connection{"sandbox\x00primary": {TenantID: "tenant", ConnectionID: "primary", ProviderCode: "sandbox", Secret: secret}}}
	service := NewService(&receiptRepository{}, registry, HMACSHA256{Now: func() time.Time { return now }, Tolerance: 5 * time.Minute})
	body := []byte(`{"event_id":"evt-1","type":"payment.updated","data":{}}`)
	if _, err := service.Receive(context.Background(), "sandbox", "primary", "1799999000", signed(secret, "1799999000", body), body); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("old timestamp err=%v", err)
	}
	if _, err := service.Receive(context.Background(), "sandbox", "primary", "1800000000", signed(secret, "1800000000", body), append(body, ' ')); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("tamper err=%v", err)
	}
	invalid := []byte(`{"event_id":"evt-1","type":"payment.updated","data":[],"extra":1}`)
	if _, err := service.Receive(context.Background(), "sandbox", "primary", "1800000000", signed(secret, "1800000000", invalid), invalid); !errors.Is(err, ErrInvalid) {
		t.Fatalf("shape err=%v", err)
	}
}

func TestConnectionConfigurationRejectsInlineOrMissingSecrets(t *testing.T) {
	lookup := func(string) (string, bool) { return "", false }
	if _, err := LoadHMACConnections([]byte(`[{"tenant_id":"t","connection_id":"c","provider_code":"p","secret":"inline"}]`), lookup); err == nil {
		t.Fatal("inline secret accepted")
	}
	if _, err := LoadHMACConnections([]byte(`[{"tenant_id":"t","connection_id":"c","provider_code":"p","secret_env":"PROVIDER_SECRET"}]`), lookup); err == nil {
		t.Fatal("missing secret accepted")
	}
}
