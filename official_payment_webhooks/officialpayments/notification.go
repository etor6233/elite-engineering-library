package officialpayments

// AUTHORED normalization only. SDK verifiers remain the sole provider signature
// implementation. Notifications schedule GET reconciliation; they never attest
// payment success merely because a callback body declares it.
import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
)

type PaymentNotification struct {
	ProviderCode      string `json:"provider_code"`
	ProviderEventID   string `json:"provider_event_id"`
	ProviderReference string `json:"provider_reference"`
	ResourceType      string `json:"resource_type"`
	Type              string `json:"type"`
	BodySHA256        string `json:"body_sha256"`
}

func bodyDigest(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

func VerifyStripePaymentNotification(payload []byte, signature, secret string, tolerance time.Duration, expectedLive bool, expectedAccount string) (PaymentNotification, error) {
	event, err := VerifyStripe(payload, signature, secret, tolerance)
	if err != nil {
		return PaymentNotification{}, err
	}
	var body struct {
		Account string `json:"account"`
		Live    *bool  `json:"livemode"`
		Data    struct {
			Object struct {
				ID            string          `json:"id"`
				Object        string          `json:"object"`
				PaymentIntent json.RawMessage `json:"payment_intent"`
			} `json:"object"`
		} `json:"data"`
	}
	if json.Unmarshal(payload, &body) != nil || event.ID == "" || len(event.ID) > 200 || body.Account != expectedAccount || body.Live == nil || *body.Live != expectedLive {
		return PaymentNotification{}, ErrInvalidObservation
	}
	resource := ""
	if strings.HasPrefix(event.Type, "payment_intent.") && body.Data.Object.Object == "payment_intent" && stripeIntentID.MatchString(body.Data.Object.ID) {
		resource = "payment_intent"
	}
	if strings.HasPrefix(event.Type, "checkout.session.") && body.Data.Object.Object == "checkout.session" && checkoutSessionID.MatchString(body.Data.Object.ID) {
		resource = "checkout_session"
	}
	if (strings.HasPrefix(event.Type, "charge.") || strings.HasPrefix(event.Type, "refund.")) && (body.Data.Object.Object == "charge" || body.Data.Object.Object == "dispute" || body.Data.Object.Object == "refund") {
		var intent string
		if json.Unmarshal(body.Data.Object.PaymentIntent, &intent) != nil {
			var object struct {
				ID string `json:"id"`
			}
			if json.Unmarshal(body.Data.Object.PaymentIntent, &object) == nil {
				intent = object.ID
			}
		}
		if stripeIntentID.MatchString(intent) {
			body.Data.Object.ID = intent
			resource = "payment_intent"
		}
	}
	if resource == "" {
		return PaymentNotification{}, ErrInvalidObservation
	}
	return PaymentNotification{ProviderCode: "stripe", ProviderEventID: event.ID, ProviderReference: body.Data.Object.ID, ResourceType: resource, Type: event.Type, BodySHA256: bodyDigest(payload)}, nil
}

func VerifyMercadoPagoPaymentNotification(payload []byte, signature, requestID, dataID, secret string, tolerance time.Duration) (PaymentNotification, error) {
	// Mercado Pago signs data.id/request-id/timestamp, not the full JSON body.
	// Only that authenticated resource ID selects GET; unsigned body status,
	// merchant and amounts cannot authorize domain effects.
	if strings.TrimSpace(requestID) == "" || len(requestID) > 128 || strings.ContainsAny(requestID, "\r\n") || !decimalPaymentID.MatchString(dataID) {
		return PaymentNotification{}, ErrInvalidInput
	}
	_, err := VerifyMercadoPago(payload, signature, requestID, dataID, secret, tolerance)
	if err != nil {
		return PaymentNotification{}, err
	}
	return PaymentNotification{ProviderCode: "mercadopago", ProviderEventID: "mp:" + requestID + ":" + dataID, ProviderReference: dataID, ResourceType: "payment", Type: "payment.reconciliation-requested", BodySHA256: bodyDigest(payload)}, nil
}
