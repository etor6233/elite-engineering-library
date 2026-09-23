package whatsappbridge

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
)

// Same local RS256/discovery/JWKS method as the existing confirmation fixture.
// Uses the real identity verifier; the issuer and its identities are synthetic.
func notificationIssuer(t *testing.T) (identity.Verifier, func(identity.Principal) string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	enc := base64.RawURLEncoding.EncodeToString
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/.well-known/openid-configuration" {
			_ = json.NewEncoder(w).Encode(map[string]any{"issuer": server.URL, "jwks_uri": server.URL + "/keys", "id_token_signing_alg_values_supported": []string{"RS256"}})
		} else if r.URL.Path == "/keys" {
			_ = json.NewEncoder(w).Encode(map[string]any{"keys": []any{map[string]any{"kty": "RSA", "use": "sig", "alg": "RS256", "kid": "notification-fixture", "n": enc(key.N.Bytes()), "e": enc(big.NewInt(int64(key.E)).Bytes())}}})
		} else {
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)
	verifier, err := identity.NewOIDCVerifier(context.Background(), server.URL, "notification-fixture")
	if err != nil {
		t.Fatal(err)
	}
	sign := func(p identity.Principal) string {
		permissions, organizations := []string{}, []string{}
		for s := range p.Permissions {
			permissions = append(permissions, s)
		}
		for s := range p.Organizations {
			organizations = append(organizations, s)
		}
		header, _ := json.Marshal(map[string]string{"alg": "RS256", "kid": "notification-fixture"})
		claims, _ := json.Marshal(map[string]any{"iss": server.URL, "aud": "notification-fixture", "sub": p.Subject, "tenant_id": p.TenantID, "permissions": permissions, "organization_ids": organizations, "exp": time.Now().Add(time.Minute).Unix()})
		unsigned := enc(header) + "." + enc(claims)
		hash := sha256.Sum256([]byte(unsigned))
		signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
		if err != nil {
			t.Fatal(err)
		}
		return unsigned + "." + enc(signature)
	}
	return verifier, sign
}

func notificationRequest(c AppointmentApprovalCommand) AppointmentNotificationRequest {
	return AppointmentNotificationRequest{OrganizationID: c.OrganizationID, ConfirmationEventID: c.ConfirmationEventID, AppointmentVersion: c.AppointmentVersion, BindingVersion: c.BindingVersion, ConsentID: c.ConsentID, ConsentEvidenceSHA256: c.ConsentEvidenceSHA256, EvidenceSHA256: c.EvidenceSHA256, ExpiresAt: c.ExpiresAt, Recipient: c.Message.ExternalID, TemplateName: "appointment_confirmed", LanguageCode: "es_AR", BodyParameters: []string{"2026-09-07T10:00Z"}}
}

func TestNotificationRequestRejectsAmbiguity(t *testing.T) {
	valid, _ := json.Marshal(notificationRequest(AppointmentApprovalCommand{ExpiresAt: time.Now()}))
	if _, err := readNotificationRequest(bytes.NewReader(valid)); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{
		"null_parameter": strings.Replace(string(valid), `"body_parameters":["2026-09-07T10:00Z"]`, `"body_parameters":[null]`, 1),
		"duplicate":      strings.Replace(string(valid), `"recipient":`, `"recipient":"other","recipient":`, 1),
		"case_alias":     strings.Replace(string(valid), `"recipient":`, `"Recipient":`, 1),
		"unknown":        strings.Replace(string(valid), `"recipient":`, `"tenant_id":`, 1),
		"null":           strings.Replace(string(valid), `"recipient":""`, `"recipient":null`, 1),
		"trailing":       string(valid) + "{}",
		"array":          "[]", "empty": "{}", "number": "1", "malformed": "{",
		"bad_type": strings.Replace(string(valid), `"binding_version":0`, `"binding_version":"0"`, 1),
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := readNotificationRequest(strings.NewReader(body)); err == nil {
				t.Fatal("ambiguous request accepted")
			}
		})
	}
}

func TestNotificationHTTPPostgresReplayAndRecovery(t *testing.T) {
	pool := approvalPool(t)
	verifier, sign := notificationIssuer(t)
	for _, mode := range []string{"accepted", "provider_lost", "http_lost"} {
		t.Run(mode, func(t *testing.T) {
			p, c, resolver := approvalFixtureData(t, pool)
			providerMode := "accepted"
			if mode == "provider_lost" {
				providerMode = "lost"
			}
			sender, _, marker := senderFixture(t, providerMode)
			sender.TenantID = p.TenantID
			store, err := postgres.NewOutboundDeliveryStore(pool, []byte("0123456789abcdef0123456789abcdef"), time.Minute)
			if err != nil {
				t.Fatal(err)
			}
			module, err := NewAppointmentNotificationModule(resolver, sender, store, receiverFixture{})
			if err != nil {
				t.Fatal(err)
			}
			mux := http.NewServeMux()
			module.Register(mux, verifier)
			var dropped atomic.Bool
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if mode == "http_lost" && dropped.CompareAndSwap(false, true) {
					record := httptest.NewRecorder()
					mux.ServeHTTP(record, r)
					if record.Code != 200 {
						t.Errorf("expected commit before lost HTTP response: %d %s", record.Code, record.Body.String())
					}
					panic(http.ErrAbortHandler)
				}
				mux.ServeHTTP(w, r)
			}))
			defer server.Close()
			body, _ := json.Marshal(notificationRequest(c))
			token := sign(p)
			call := func() (int, string, error) {
				req, err := http.NewRequest(http.MethodPost, server.URL+"/v1/franchise/appointments/appointment-1/whatsapp-confirmation", bytes.NewReader(body))
				if err != nil {
					return 0, "", err
				}
				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				response, err := server.Client().Do(req)
				if err != nil {
					return 0, "", err
				}
				defer response.Body.Close()
				raw, err := io.ReadAll(response.Body)
				return response.StatusCode, string(raw), err
			}
			if mode == "http_lost" {
				if _, _, err := call(); err == nil {
					t.Fatal("HTTP response loss was not observed")
				}
			}
			if mode == "accepted" {
				var wg sync.WaitGroup
				for i := 0; i < 8; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						status, _, err := call()
						if err != nil || (status != 200 && status != 409) {
							t.Errorf("concurrent call: %d %v", status, err)
						}
					}()
				}
				wg.Wait()
			}
			for i := 0; i < 2; i++ {
				status, response, err := call()
				want := 200
				if mode == "provider_lost" {
					want = 409
				}
				if err != nil || status != want {
					t.Fatalf("%d %s %v", status, response, err)
				}
				if want == 200 && (!strings.Contains(response, `"status":"accepted"`) || strings.Contains(response, "delivered")) {
					t.Fatal(response)
				}
				if want == 409 && !strings.Contains(response, "DELIVERY_RECONCILIATION_REQUIRED") {
					t.Fatal(response)
				}
				if strings.Contains(response, c.Message.ExternalID) || strings.Contains(response, "synthetic-token") {
					t.Fatal("response leaked sensitive data")
				}
			}
			calls, err := os.ReadFile(marker)
			if err != nil || string(calls) != "call\n" {
				t.Fatal("duplicate provider call", string(calls), err)
			}
			var grants, attempts int
			var state string
			if err = pool.QueryRow(context.Background(), `select count(*) from communication.whatsapp_appointment_approval where tenant_id=$1`, p.TenantID).Scan(&grants); err != nil || grants != 1 {
				t.Fatal(grants, err)
			}
			if err = pool.QueryRow(context.Background(), `select state,attempt_count from communication.outbound_delivery where tenant_id=$1`, p.TenantID).Scan(&state, &attempts); err != nil || attempts != 1 {
				t.Fatal(state, attempts, err)
			}
			wantState := "accepted"
			if mode == "provider_lost" {
				wantState = "unknown"
			}
			if state != wantState {
				t.Fatal(state)
			}
			// The identical event cannot silently change template parameters on replay.
			body = bytes.Replace(body, []byte("2026-09-07T10:00Z"), []byte("changed"), 1)
			if status, _, err := call(); err != nil || status != 409 {
				t.Fatal("divergent replay", status, err)
			}
		})
	}
}

func TestNotificationHTTPRejectsBeforeFence(t *testing.T) {
	pool := approvalPool(t)
	verifier, sign := notificationIssuer(t)
	for _, mode := range []string{"no_token", "invalid_signature", "role", "organization", "tenant", "recipient", "consent", "revoked", "path", "oversized", "media", "injected_identity", "duplicate_auth"} {
		t.Run(mode, func(t *testing.T) {
			p, c, resolver := approvalFixtureData(t, pool)
			sender, _, marker := senderFixture(t, "accepted")
			sender.TenantID = p.TenantID
			store, err := postgres.NewOutboundDeliveryStore(pool, []byte("0123456789abcdef0123456789abcdef"), time.Minute)
			if err != nil {
				t.Fatal(err)
			}
			module, err := NewAppointmentNotificationModule(resolver, sender, store, receiverFixture{})
			if err != nil {
				t.Fatal(err)
			}
			mux := http.NewServeMux()
			module.Register(mux, verifier)
			input := notificationRequest(c)
			status := 409
			pathID := "appointment-1"
			switch mode {
			case "no_token", "invalid_signature", "duplicate_auth":
				status = 401
			case "role":
				p.Permissions = map[string]struct{}{}
				status = 403
			case "organization":
				p.Organizations = map[string]struct{}{}
				status = 403
			case "tenant":
				p.TenantID = "other-tenant"
				status = 403
			case "recipient":
				input.Recipient = "5491199999999"
			case "consent":
				input.ConsentEvidenceSHA256 = strings.Repeat("f", 64)
			case "path":
				pathID = "other-appointment"
			case "revoked":
				if _, err := pool.Exec(context.Background(), `update communication.contact_channel_binding set state='revoked',pii_allowed=false,version=version+1 where tenant_id=$1`, p.TenantID); err != nil {
					t.Fatal(err)
				}
			case "oversized", "injected_identity":
				status = 400
			case "media":
				status = 415
			}
			body, _ := json.Marshal(input)
			if mode == "oversized" {
				body = bytes.Replace(body, []byte(input.Recipient), bytes.Repeat([]byte("1"), 65536), 1)
			}
			if mode == "injected_identity" {
				body = append([]byte(`{"tenant_id":"injected",`), body[1:]...)
			}
			req := httptest.NewRequest(http.MethodPost, "/v1/franchise/appointments/"+pathID+"/whatsapp-confirmation", bytes.NewReader(body))
			token := sign(p)
			if mode == "invalid_signature" {
				token += "tamper"
			}
			if mode != "no_token" {
				req.Header.Set("Authorization", "Bearer "+token)
			}
			if mode == "duplicate_auth" {
				req.Header.Add("Authorization", "Bearer "+token)
			}
			req.Header.Set("Content-Type", "application/json")
			if mode == "media" {
				req.Header.Set("Content-Type", "application/jsonp")
			}
			record := httptest.NewRecorder()
			mux.ServeHTTP(record, req)
			if record.Code != status || record.Header().Get("Cache-Control") != "no-store" || record.Header().Get("Content-Type") != "application/problem+json" {
				t.Fatal(record.Code, record.Body.String())
			}
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatal("child invoked for rejected request")
			}
			for _, table := range []string{"communication.whatsapp_appointment_approval", "communication.outbound_delivery"} {
				var count int
				if err := pool.QueryRow(context.Background(), `select count(*) from `+table+` where tenant_id=$1`, c.Message.TenantID).Scan(&count); err != nil || count != 0 {
					t.Fatal(table, count, err)
				}
			}
		})
	}
}
