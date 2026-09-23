package officialpayments

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	mpwebhook "github.com/mercadopago/sdk-go/pkg/webhook"
	stripewebhook "github.com/stripe/stripe-go/v86/webhook"
)

const MaxPayloadBytes = 1 << 20

var (
	ErrInvalidInput      = errors.New("invalid payment webhook input")
	ErrPayloadTooLarge   = errors.New("payment webhook payload too large")
	ErrPayloadIDMismatch = errors.New("payment webhook payload id mismatch")
)

type Event struct {
	Provider  string
	ID        string
	Type      string
	RequestID string
	Payload   []byte
}

func VerifyStripe(payload []byte, signatureHeader, signingSecret string, tolerance time.Duration) (Event, error) {
	if len(payload) == 0 || strings.TrimSpace(signatureHeader) == "" || strings.TrimSpace(signingSecret) == "" || tolerance <= 0 {
		return Event{}, ErrInvalidInput
	}
	if len(payload) > MaxPayloadBytes {
		return Event{}, ErrPayloadTooLarge
	}
	event, err := stripewebhook.ConstructEventWithTolerance(payload, signatureHeader, signingSecret, tolerance)
	if err != nil {
		return Event{}, fmt.Errorf("stripe official webhook verification: %w", err)
	}
	return Event{Provider: "stripe", ID: event.ID, Type: string(event.Type), Payload: append([]byte(nil), payload...)}, nil
}

func VerifyMercadoPago(payload []byte, signatureHeader, requestID, dataID, signingSecret string, tolerance time.Duration) (Event, error) {
	if len(payload) == 0 || strings.TrimSpace(signatureHeader) == "" || strings.TrimSpace(signingSecret) == "" || strings.TrimSpace(dataID) == "" || tolerance <= 0 {
		return Event{}, ErrInvalidInput
	}
	if len(payload) > MaxPayloadBytes {
		return Event{}, ErrPayloadTooLarge
	}
	if err := mpwebhook.ValidateSignature(signatureHeader, requestID, dataID, signingSecret, mpwebhook.WithTolerance(tolerance)); err != nil {
		return Event{}, fmt.Errorf("mercado pago official webhook verification: %w", err)
	}
	var body struct {
		Type string `json:"type"`
		Data struct {
			ID json.RawMessage `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(payload, &body); err != nil {
		return Event{}, fmt.Errorf("mercado pago payload JSON: %w", err)
	}
	bodyID := strings.Trim(string(body.Data.ID), `"`)
	if bodyID != "" && !strings.EqualFold(bodyID, dataID) {
		return Event{}, ErrPayloadIDMismatch
	}
	return Event{Provider: "mercado_pago", ID: dataID, Type: body.Type, RequestID: requestID, Payload: append([]byte(nil), payload...)}, nil
}
