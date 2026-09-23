package refundworker

// AUTHORED compatibility fixtures: real pinned Mercado Pago SDK calls use an
// in-process requester. No network, account, secret or provider capture occurs.
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestPostgresMercadoPagoProviderAliases(t *testing.T) {
	for _, sourceProvider := range []string{"mercado_pago", "mercadopago"} {
		t.Run(sourceProvider, func(t *testing.T) {
			creates := make(map[string]int)
			paymentReads, refundReads := 0, 0
			transport := requesterFunc(func(request *http.Request) (*http.Response, error) {
				if request.Header.Get("Authorization") != "Bearer TEST-local" || request.URL.Host != "api.mercadopago.com" {
					return nil, fmt.Errorf("unexpected SDK authority")
				}
				var response string
				switch {
				case request.Method == http.MethodGet && request.URL.Path == "/v1/payments/7186040733":
					paymentReads++
					response = `{"id":7186040733,"status":"approved","captured":true,"currency_id":"ARS","transaction_amount":20.00}`
				case request.Method == http.MethodPost && request.URL.Path == "/v1/payments/7186040733/refunds":
					key := request.Header.Get("X-Idempotency-Key")
					if key != "return-refund-key-0001" && key != "return-refund-key-0002" {
						return nil, fmt.Errorf("historical refund idempotency key changed")
					}
					var body map[string]any
					if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
						return nil, err
					}
					if body["amount"] != float64(10) {
						return nil, fmt.Errorf("refund amount changed: %v", body["amount"])
					}
					creates[key]++
					id := 1622029221
					if key == "return-refund-key-0002" {
						id++
					}
					response = fmt.Sprintf(`{"id":%d,"payment_id":7186040733,"status":"approved","amount":10.00}`, id)
				case request.Method == http.MethodGet && request.URL.Path == "/v1/payments/7186040733/refunds/1622029221":
					refundReads++
					response = `{"id":1622029221,"payment_id":7186040733,"status":"approved","amount":10.00}`
				default:
					return nil, fmt.Errorf("unexpected SDK request %s %s", request.Method, request.URL.Path)
				}
				return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(response)), Request: request}, nil
			})
			provider, err := newMercadoPagoProvider("TEST-local", transport)
			if err != nil {
				t.Fatal(err)
			}
			testRefundLineAllocationPartialThenFullAndAmbiguity(t, sourceProvider, provider)
			if len(creates) != 2 || creates["return-refund-key-0001"] != 1 || creates["return-refund-key-0002"] != 1 || paymentReads != 3 || refundReads != 1 {
				t.Fatalf("SDK effects/reconciliation changed: creates=%v payment_reads=%d refund_reads=%d", creates, paymentReads, refundReads)
			}
		})
	}
}
