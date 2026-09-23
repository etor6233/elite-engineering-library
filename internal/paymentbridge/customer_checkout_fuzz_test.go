package paymentbridge

// AUTHORED verification glue for the provider redirect boundary. These fixtures
// contain synthetic provider identifiers and do not contact either provider.
import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func FuzzCustomerCheckoutRedirect(f *testing.F) {
	for _, seed := range []struct{ url, provider string }{
		{"https://checkout.stripe.com/c/pay/cs_test_connected", "stripe"},
		{"https://sandbox.mercadopago.com.ar/checkout/v1/redirect?pref_id=pref_fixture", "mercadopago"},
		{"https://www.mercadopago.com.ar/checkout/v1/redirect?pref_id=pref_fixture", "mercadopago"},
		{"https://checkout.stripe.com.evil.invalid/c/pay/test", "stripe"},
		{"https://user@checkout.stripe.com/c/pay/test", "stripe"},
		{"http://checkout.stripe.com/c/pay/test", "stripe"},
		{"https://checkout.stripe.com/c/pay/test#other", "stripe"},
		{"https://checkout.stripe.com/c/pay/\xff", "stripe"},
		{"https://checkout.stripe.com/c/pay/test", "mercadopago"},
	} {
		f.Add(seed.url, seed.provider, int64(300))
	}
	f.Add("https://checkout.stripe.com/c/pay/expired", "stripe", int64(0))
	f.Fuzz(func(t *testing.T, rawURL, provider string, offset int64) {
		if len(rawURL) > 9000 || len(provider) > 128 {
			t.Skip()
		}
		now := time.Unix(1800000000, 0).UTC()
		c := CustomerCheckout{OrderID: "order-fixture", ProviderCode: provider, URL: rawURL, ExpiresAt: now.Add(time.Duration(offset%86400) * time.Second)}
		if !c.Valid(now) {
			return
		}
		// The HTTP consumer and JSON transport must agree with admission. A
		// parser discrepancy must not create a different redirect downstream.
		req, err := http.NewRequest(http.MethodGet, c.URL, nil)
		if err != nil {
			t.Fatalf("admitted redirect cannot be consumed: %v", err)
		}
		allowed := map[string]map[string]bool{
			"stripe":      {"checkout.stripe.com": true},
			"mercadopago": {"www.mercadopago.com.ar": true, "sandbox.mercadopago.com.ar": true},
		}
		if req.URL.Scheme != "https" || req.URL.Opaque != "" || req.URL.User != nil || req.URL.Fragment != "" || !allowed[provider][strings.ToLower(req.URL.Hostname())] {
			t.Fatal("admitted redirect escapes its provider origin")
		}
		wire, err := json.Marshal(c)
		if err != nil {
			t.Fatal(err)
		}
		var restored CustomerCheckout
		if json.Unmarshal(wire, &restored) != nil || restored.URL != c.URL || restored.ProviderCode != c.ProviderCode || restored.OrderID != c.OrderID || !restored.ExpiresAt.Equal(c.ExpiresAt) || !restored.Valid(now) {
			t.Fatal("JSON transport changes an admitted redirect")
		}
		expired := c
		expired.ExpiresAt = now
		if expired.Valid(now) {
			t.Fatal("expired redirect remains usable")
		}
		other := c
		if provider == "stripe" {
			other.ProviderCode = "mercadopago"
		} else {
			other.ProviderCode = "stripe"
		}
		if other.Valid(now) {
			t.Fatal("redirect can be transplanted to another provider")
		}
		for _, mutate := range []func(*url.URL){
			func(u *url.URL) { u.Host = "checkout.stripe.com.evil.invalid" },
			func(u *url.URL) { u.User = url.User("confused") },
			func(u *url.URL) { u.Scheme = "http" },
			func(u *url.URL) { u.Fragment = "changed" },
		} {
			u := *req.URL
			mutate(&u)
			changed := c
			changed.URL = u.String()
			if changed.Valid(now) {
				t.Fatal("redirect mutation crosses an admitted boundary")
			}
		}
	})
}
