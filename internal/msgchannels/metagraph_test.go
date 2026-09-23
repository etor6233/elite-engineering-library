package msgchannels

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func sign(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestVerifyWebhookSignature(t *testing.T) {
	body := []byte(`{"object":"page"}`)
	secret := "app-secret"
	if !VerifyWebhookSignature(secret, body, sign(secret, body)) {
		t.Fatal("valid signature rejected")
	}
	if VerifyWebhookSignature(secret, body, "sha256=deadbeef") {
		t.Fatal("wrong signature accepted")
	}
	if VerifyWebhookSignature(secret, body, "sha1=deadbeef") {
		t.Fatal("non-sha256 accepted")
	}
	if VerifyWebhookSignature("", body, sign(secret, body)) {
		t.Fatal("empty secret accepted")
	}
}

func TestParseWebhook(t *testing.T) {
	raw := []byte(`{"object":"page","entry":[{"messaging":[{"sender":{"id":"55119999"},"recipient":{"id":"p1"},"message":{"mid":"m1","text":"hola"}}]}]}`)
	msgs, err := ParseWebhook(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 1 || msgs[0].From != "55119999" || msgs[0].Text != "hola" {
		t.Fatalf("unexpected parse: %+v", msgs)
	}
}

func TestSendTextRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v22.0/phone-1/messages" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("authorization") != "Bearer tok" {
			t.Errorf("missing bearer")
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["to"] != "55119999" || body["type"] != "text" {
			t.Errorf("unexpected body: %+v", body)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message_id":"m1"}`))
	}))
	defer srv.Close()

	cfg := GraphConfig{GraphURL: srv.URL, Version: "v22.0", SenderID: "phone-1", AccessToken: "tok", Timeout: 5 * time.Second}
	if err := SendText(context.Background(), cfg, "55119999", "hola"); err != nil {
		t.Fatal(err)
	}
}

func TestGraphConfigValidate(t *testing.T) {
	cfg := GraphConfig{GraphURL: "https://graph.facebook.com", Version: "v22.0", SenderID: "p1", AccessToken: "t", Timeout: time.Second}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}
	bad := cfg
	bad.AccessToken = ""
	if err := bad.Validate(); err == nil {
		t.Fatal("missing token accepted")
	}
}
