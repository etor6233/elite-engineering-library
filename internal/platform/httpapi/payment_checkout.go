package httpapi

// AUTHORED authorization/transport glue. Customer identity is always taken from
// the verified token; provider redirects come only from the scoped repository.
import (
	"elite.local/enterprise/internal/paymentbridge"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const PaymentProviderWebhookPath = "/v1/payment-provider/webhook"

type PaymentCheckoutModule struct {
	Reader  paymentbridge.CustomerCheckoutReader
	Webhook http.Handler
}

func (m PaymentCheckoutModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	mux.HandleFunc("GET /v1/customer/orders/{id}/checkout", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		if verifier == nil {
			writeProblem(w, 503, "CHECKOUT_UNAVAILABLE", "checkout is unavailable")
			return
		}
		p, err := authenticate(r.Context(), r.Header.Get("Authorization"), verifier)
		if err != nil || p.Subject == "" || p.TenantID == "" {
			writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
			return
		}
		if !p.Allowed("customer:self") {
			writeProblem(w, 403, "FORBIDDEN", "customer:self permission is required")
			return
		}
		q, err := url.ParseQuery(r.URL.RawQuery)
		organization := q.Get("organization_id")
		order := r.PathValue("id")
		if err != nil || len(r.URL.RawQuery) > 1024 || len(q) != 1 || len(q["organization_id"]) != 1 || order == "" || len(order) > 200 || strings.ContainsAny(order, "/\\\r\n") || strings.TrimSpace(order) != order {
			writeProblem(w, 400, "INVALID_CHECKOUT_QUERY", "one organization scope and order are required")
			return
		}
		if !p.AllowedOrganization(organization) {
			writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "organization scope is required")
			return
		}
		if m.Reader == nil {
			writeProblem(w, 503, "CHECKOUT_UNAVAILABLE", "checkout is unavailable")
			return
		}
		value, err := m.Reader.CustomerCheckout(r.Context(), p.TenantID, organization, p.Subject, order)
		if errors.Is(err, paymentbridge.ErrCheckoutNotAvailable) {
			writeProblem(w, 404, "CHECKOUT_NOT_AVAILABLE", "checkout is not available")
			return
		}
		if err != nil || value.OrderID != order || !value.Valid(time.Now()) {
			writeProblem(w, 503, "CHECKOUT_UNAVAILABLE", "checkout is unavailable")
			return
		}
		writeJSON(w, 200, value)
	})
	if m.Webhook != nil {
		mux.Handle(PaymentProviderWebhookPath, m.Webhook)
	}
}
