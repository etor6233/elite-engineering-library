package identity

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jose "github.com/go-jose/go-jose/v4"
)

func TestServiceTokenTLSChainAndSignatureFailure(t *testing.T) {
	f := newServiceIssuerFixture(t)
	handler := f.server.Config.Handler
	tlsServer := httptest.NewUnstartedServer(handler)
	tlsServer.Config.ErrorLog = log.New(io.Discard, "", 0)
	tlsServer.StartTLS()
	t.Cleanup(tlsServer.Close)
	f.server = tlsServer
	p := serviceProfileForTest(t, f.document())
	secret := func(context.Context) (string, error) { return f.secret, nil }
	if _, err := NewServiceTokenBroker(context.Background(), p, secret); err != ErrServiceTokenConfiguration {
		t.Fatal("untrusted CA accepted", err)
	}
	b, err := newServiceTokenBroker(context.Background(), p, secret, tlsServer.Client().Transport)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = b.AccessToken(context.Background()); err != nil {
		t.Fatal("trusted fixture CA failed", err)
	}
	// A token signed by a different key with the same kid cannot pass a cached JWKS key.
	old := f.keys[0]
	replacement, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	f.keys[0] = replacement
	b.token.Expiry = time.Now()
	// Keep JWKS transport serving old key while token endpoint signs with replacement.
	original := tlsServer.Config.Handler
	tlsServer.Config.Handler = httpHandlerForJWKS(old, original)
	if token, err := b.AccessToken(context.Background()); err != ErrServiceTokenUnavailable || token != "" {
		t.Fatal("wrong signature accepted", err)
	}
}
func httpHandlerForJWKS(key *rsa.PrivateKey, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/jwks" {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{Key: &key.PublicKey, KeyID: "fixture-a", Use: "sig", Algorithm: "RS256"}}})
	})
}
