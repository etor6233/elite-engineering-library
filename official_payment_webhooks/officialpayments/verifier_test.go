package officialpayments

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stripe/stripe-go/v86"
	stripewebhook "github.com/stripe/stripe-go/v86/webhook"
)

func TestVerifyStripeUsesOfficialSignatureAndAPIVersion(t *testing.T) {
	payload := []byte(fmt.Sprintf(`{"id":"evt_123","object":"event","api_version":%q,"type":"payment_intent.succeeded","data":{"object":{}}}`, stripe.APIVersion))
	signed := stripewebhook.GenerateTestSignedPayload(&stripewebhook.UnsignedPayload{Payload: payload, Secret: "whsec_test", Timestamp: time.Now()})
	event, err := VerifyStripe(payload, signed.Header, "whsec_test", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if event.Provider != "stripe" || event.ID != "evt_123" || event.Type != "payment_intent.succeeded" {
		t.Fatalf("unexpected event: %#v", event)
	}

	tampered := append([]byte(nil), payload...)
	tampered[len(tampered)-2] ^= 1
	if _, err := VerifyStripe(tampered, signed.Header, "whsec_test", 5*time.Minute); err == nil {
		t.Fatal("tampered Stripe payload accepted")
	}
}

func mercadoPagoHeader(dataID, requestID, secret string, timestamp time.Time) string {
	ts := fmt.Sprintf("%d", timestamp.UnixMilli())
	manifest := "id:" + dataID + ";request-id:" + requestID + ";ts:" + ts + ";"
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(manifest))
	return "ts=" + ts + ",v1=" + hex.EncodeToString(mac.Sum(nil))
}

func TestVerifyMercadoPagoUsesOfficialSignatureAndCrossChecksID(t *testing.T) {
	const dataID = "123456"
	const requestID = "request-123"
	const secret = "mp_secret"
	payload := []byte(`{"type":"payment","data":{"id":"123456"}}`)
	header := mercadoPagoHeader(dataID, requestID, secret, time.Now())
	event, err := VerifyMercadoPago(payload, header, requestID, dataID, secret, 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if event.Provider != "mercado_pago" || event.ID != dataID || event.Type != "payment" || event.RequestID != requestID {
		t.Fatalf("unexpected event: %#v", event)
	}

	mismatch := []byte(`{"type":"payment","data":{"id":"different"}}`)
	if _, err := VerifyMercadoPago(mismatch, header, requestID, dataID, secret, 5*time.Minute); !errors.Is(err, ErrPayloadIDMismatch) {
		t.Fatalf("expected ID mismatch, got %v", err)
	}
	if _, err := VerifyMercadoPago(payload, strings.Replace(header, "v1=", "v1=00", 1), requestID, dataID, secret, 5*time.Minute); err == nil {
		t.Fatal("tampered Mercado Pago signature accepted")
	}
}

func TestFailClosedInputAndSize(t *testing.T) {
	if _, err := VerifyStripe([]byte(`{}`), "", "secret", time.Minute); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
	tooLarge := make([]byte, MaxPayloadBytes+1)
	if _, err := VerifyStripe(tooLarge, "header", "secret", time.Minute); !errors.Is(err, ErrPayloadTooLarge) {
		t.Fatalf("expected too large, got %v", err)
	}
	if _, err := VerifyMercadoPago([]byte(`{}`), "header", "request", "", "secret", time.Minute); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}
