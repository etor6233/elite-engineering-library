package postgres_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	jose "github.com/go-jose/go-jose/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type portalFixture struct {
	mu                                                      sync.Mutex
	issuer                                                  *httptest.Server
	api                                                     *httptest.Server
	keys                                                    []*rsa.PrivateKey
	keyIndex                                                int
	mode                                                    string
	codes                                                   map[string]url.Values
	refresh                                                 map[string]bool
	counter                                                 int
	grants, refreshGrants, serviceGrants, jwks, revocations int
	tokenAuthValid, tokenClientValid, tokenBodyGrantValid   bool
	pool                                                    *pgxpool.Pool
	profile                                                 identity.PortalProfile
}

func (f *portalFixture) jwt(claims map[string]any) string {
	raw, _ := json.Marshal(claims)
	signer, _ := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: f.keys[f.keyIndex]}, (&jose.SignerOptions{}).WithHeader("kid", fmt.Sprintf("fixture-%d", f.keyIndex)))
	signed, _ := signer.Sign(raw)
	value, _ := signed.CompactSerialize()
	return value
}
func (f *portalFixture) issue(w http.ResponseWriter, nonce string, service bool) {
	now := time.Now()
	f.counter++
	claims := map[string]any{"iss": f.issuer.URL, "sub": "person", "aud": "portal-client", "iat": now.Unix(), "exp": now.Add(120 * time.Second).Unix(), "tenant_id": "tenant", "organization_ids": []string{"org"}, "permissions": []string{"customer:read"}}
	if nonce != "" {
		claims["nonce"] = nonce
	}
	if service {
		claims["sub"] = "portal-service"
		claims["aud"] = "api"
		claims["permissions"] = []string{"portal-session:manage"}
		claims["scope"] = "portal-session:manage"
	}
	if !service {
		switch f.mode {
		case "permission_withdrawn":
			claims["permissions"] = []string{}
		case "subject_changed":
			claims["sub"] = "different-person"
		case "wrong_issuer":
			claims["iss"] = "https://foreign.invalid"
		case "wrong_audience":
			claims["aud"] = "foreign"
		case "expired":
			claims["exp"] = now.Add(-time.Hour).Unix()
		case "wrong_nonce":
			claims["nonce"] = "foreign"
		case "foreign_tenant":
			claims["tenant_id"] = "foreign"
		case "extra_permission":
			claims["permissions"] = []string{"admin:write"}
		}
	}
	id := f.jwt(claims)
	claims["aud"] = "api"
	access := f.jwt(claims)
	response := map[string]any{"access_token": access, "token_type": "Bearer", "expires_in": 120}
	if !service {
		refresh := fmt.Sprintf("fixture-refresh-%d", f.counter)
		f.refresh[refresh] = true
		response["id_token"] = id
		response["refresh_token"] = refresh
		response["scope"] = "openid profile offline_access"
	}
	if f.mode == "grant_lost" && !service {
		_, _ = w.Write([]byte(`{"truncated":`))
		f.mode = ""
		return
	}
	_ = json.NewEncoder(w).Encode(response)
}
func (f *portalFixture) serveIssuer(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	switch r.URL.Path {
	case "/.well-known/openid-configuration":
		_ = json.NewEncoder(w).Encode(map[string]any{"issuer": f.issuer.URL, "authorization_endpoint": f.issuer.URL + "/authorize", "token_endpoint": f.issuer.URL + "/token", "jwks_uri": f.issuer.URL + "/jwks", "revocation_endpoint": f.issuer.URL + "/revoke", "end_session_endpoint": f.issuer.URL + "/logout", "response_types_supported": []string{"code"}, "subject_types_supported": []string{"public"}, "id_token_signing_alg_values_supported": []string{"RS256"}, "code_challenge_methods_supported": []string{"S256"}, "token_endpoint_auth_methods_supported": []string{"client_secret_basic"}, "grant_types_supported": []string{"authorization_code", "refresh_token", "client_credentials"}})
	case "/jwks":
		f.jwks++
		_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{Key: &f.keys[f.keyIndex].PublicKey, KeyID: fmt.Sprintf("fixture-%d", f.keyIndex), Algorithm: "RS256", Use: "sig"}}})
	case "/authorize":
		q := r.URL.Query()
		if q.Get("client_id") != "portal-client" || q.Get("code_challenge_method") != "S256" || q.Get("code_challenge") == "" || q.Get("state") == "" || q.Get("nonce") == "" || q.Get("redirect_uri") != f.profile.Document().CallbackURL {
			w.WriteHeader(400)
			return
		}
		f.counter++
		code := fmt.Sprintf("fixture-code-%d", f.counter)
		f.codes[code] = q
		callback, _ := url.Parse(q.Get("redirect_uri"))
		query := callback.Query()
		query.Set("code", code)
		query.Set("state", q.Get("state"))
		callback.RawQuery = query.Encode()
		http.Redirect(w, r, callback.String(), 303)
	case "/token":
		_ = r.ParseForm()
		client, secret, ok := r.BasicAuth()
		client, _ = url.QueryUnescape(client)
		secret, _ = url.QueryUnescape(secret)
		f.tokenAuthValid = ok && secret == "synthetic-oidc-only"
		f.tokenClientValid = client == "portal-client" || client == "portal-service-client"
		f.tokenBodyGrantValid = r.Form.Get("grant_type") == "authorization_code" || r.Form.Get("grant_type") == "refresh_token" || r.Form.Get("grant_type") == "client_credentials"
		if !ok || secret != "synthetic-oidc-only" {
			w.WriteHeader(401)
			return
		}
		if f.mode == "account_disabled" && client == "portal-client" {
			w.WriteHeader(400)
			_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
			return
		}
		switch r.Form.Get("grant_type") {
		case "client_credentials":
			if client != "portal-service-client" || r.Form.Get("scope") != "portal-session:manage" || r.Form.Get("audience") != "api" {
				w.WriteHeader(400)
				return
			}
			f.serviceGrants++
			f.issue(w, "", true)
		case "authorization_code":
			q, exists := f.codes[r.Form.Get("code")]
			sum := sha256.Sum256([]byte(r.Form.Get("code_verifier")))
			if client != "portal-client" || !exists || base64.RawURLEncoding.EncodeToString(sum[:]) != q.Get("code_challenge") || r.Form.Get("redirect_uri") != q.Get("redirect_uri") {
				w.WriteHeader(400)
				_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
				return
			}
			delete(f.codes, r.Form.Get("code"))
			f.grants++
			f.issue(w, q.Get("nonce"), false)
		case "refresh_token":
			if client != "portal-client" || !f.refresh[r.Form.Get("refresh_token")] {
				w.WriteHeader(400)
				_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
				return
			}
			f.refresh[r.Form.Get("refresh_token")] = false
			f.refreshGrants++
			f.issue(w, "", false)
		default:
			w.WriteHeader(400)
		}
	case "/revoke":
		_ = r.ParseForm()
		client, secret, ok := r.BasicAuth()
		client, _ = url.QueryUnescape(client)
		secret, _ = url.QueryUnescape(secret)
		if !ok || client != "portal-client" || secret != "synthetic-oidc-only" {
			w.WriteHeader(401)
			return
		}
		if f.mode == "revoke_unavailable" {
			w.WriteHeader(503)
			_, _ = w.Write([]byte(`{"error":"temporarily_unavailable"}`))
			return
		}
		f.revocations++
		if r.Form.Get("token_type_hint") == "refresh_token" {
			f.refresh[r.Form.Get("token")] = false
		}
		_, _ = w.Write([]byte(`{}`))
	case "/logout":
		http.Redirect(w, r, f.profile.Document().PostLogoutURL, 303)
	default:
		w.WriteHeader(404)
	}
}
func newPortalFixture(t *testing.T, pool *pgxpool.Pool) *portalFixture {
	f := &portalFixture{pool: pool, codes: map[string]url.Values{}, refresh: map[string]bool{}}
	for range 2 {
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}
		f.keys = append(f.keys, key)
	}
	f.issuer = httptest.NewServer(http.HandlerFunc(f.serveIssuer))
	t.Cleanup(f.issuer.Close)
	var handler http.Handler
	f.api = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/fixture/") {
			f.control(w, r)
			return
		}
		f.mu.Lock()
		lost := f.mode == "commit_lost" && strings.HasSuffix(r.URL.Path, "/commit")
		if lost {
			f.mode = ""
		}
		f.mu.Unlock()
		if lost {
			record := httptest.NewRecorder()
			handler.ServeHTTP(record, r)
			if record.Code != 200 {
				t.Error("commit fixture failed", record.Code)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"truncated":`))
			return
		}
		handler.ServeHTTP(w, r)
	}))
	t.Cleanup(f.api.Close)
	d := identity.PortalProfileDocument{Schema: "elite.oidc.portal-lifecycle.v1", ProfileID: "fixture-portal", Revision: 1, Transport: "LOOPBACK_FIXTURE", Issuer: f.issuer.URL, AuthorizationEndpoint: f.issuer.URL + "/authorize", TokenEndpoint: f.issuer.URL + "/token", JWKSEndpoint: f.issuer.URL + "/jwks", RevocationEndpoint: f.issuer.URL + "/revoke", EndSessionEndpoint: f.issuer.URL + "/logout", ClientID: "portal-client", CallbackURL: "http://127.0.0.1:4567/api/auth/callback", PostLogoutURL: "http://127.0.0.1:4567/", BridgeURL: f.api.URL, BackendAudience: "api", ServiceClientID: "portal-service-client", ServiceSubject: "portal-service", ServiceScopes: []string{"portal-session:manage"}, ServiceAudienceParameter: "api", TenantID: "tenant", Organizations: []string{"org"}, AllowedPermissions: []string{"customer:read"}, Scopes: []string{"openid", "profile", "offline_access"}, MaximumSessionSeconds: 3600, RefreshBeforeSeconds: 15, RetentionSeconds: 86400}
	raw, _ := json.Marshal(d)
	sum := sha256.Sum256(raw)
	p, err := identity.LoadPortalProfile(raw, hex.EncodeToString(sum[:]))
	if err != nil {
		t.Fatal(err)
	}
	f.profile = p
	verifier, err := identity.NewOIDCVerifier(context.Background(), f.issuer.URL, "api")
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	httpapi.PortalSessionModule{Profile: p, Store: &db.PortalSessions{Pool: pool}}.Register(mux, verifier)
	handler = mux
	return f
}
func (f *portalFixture) control(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" && r.URL.Path == "/fixture/control" {
		var input struct {
			Mode           string `json:"mode"`
			Expire         bool   `json:"expire"`
			Rotate         bool   `json:"rotate"`
			Swap           bool   `json:"swap"`
			ExpireAbsolute bool   `json:"expire_absolute"`
			Retention      bool   `json:"retention"`
		}
		if json.NewDecoder(r.Body).Decode(&input) != nil {
			w.WriteHeader(400)
			return
		}
		f.mode = input.Mode
		if input.Rotate {
			f.keyIndex = 1
		}
		if input.Expire {
			_, err := f.pool.Exec(r.Context(), `update platform.portal_session set access_expires_at=clock_timestamp()-interval '1 second' where profile_sha256=$1 and state='active'`, f.profile.SHA256())
			if err != nil {
				w.WriteHeader(500)
				return
			}
		}
		if input.Swap {
			tx, err := f.pool.Begin(r.Context())
			if err != nil {
				w.WriteHeader(500)
				return
			}
			defer tx.Rollback(r.Context())
			rows, err := tx.Query(r.Context(), `select session_id_sha256,ciphertext,ciphertext_sha256 from platform.portal_session where profile_sha256=$1 order by created_at desc limit 2 for update`, f.profile.SHA256())
			if err != nil {
				w.WriteHeader(500)
				return
			}
			var ids, hashes [][]byte
			var payloads []string
			for rows.Next() {
				var id, hash []byte
				var body string
				if rows.Scan(&id, &body, &hash) != nil {
					w.WriteHeader(500)
					return
				}
				ids = append(ids, id)
				hashes = append(hashes, hash)
				payloads = append(payloads, body)
			}
			rows.Close()
			if len(ids) != 2 {
				w.WriteHeader(500)
				return
			}
			for i := 0; i < 2; i++ {
				_, err = tx.Exec(r.Context(), `update platform.portal_session set ciphertext=$3,ciphertext_sha256=$4 where profile_sha256=$1 and session_id_sha256=$2`, f.profile.SHA256(), ids[i], payloads[1-i], hashes[1-i])
				if err != nil {
					w.WriteHeader(500)
					return
				}
			}
			if tx.Commit(r.Context()) != nil {
				w.WriteHeader(500)
				return
			}
		}
		if input.ExpireAbsolute || input.Retention {
			interval := "1 second"
			if input.Retention {
				interval = "2 days"
			}
			_, err := f.pool.Exec(r.Context(), `update platform.portal_session set absolute_expires_at=clock_timestamp()-$2::interval,access_expires_at=least(access_expires_at,clock_timestamp()-$2::interval) where profile_sha256=$1`, f.profile.SHA256(), interval)
			if err != nil {
				w.WriteHeader(500)
				return
			}
		}
		_, _ = w.Write([]byte(`{}`))
		return
	}
	if r.URL.Path == "/fixture/stats" {
		var active, revoked, reauth int
		_ = f.pool.QueryRow(r.Context(), `select count(*) filter(where state='active'),count(*) filter(where state='revoked'),count(*) filter(where state='reauth_required') from platform.portal_session where profile_sha256=$1`, f.profile.SHA256()).Scan(&active, &revoked, &reauth)
		_ = json.NewEncoder(w).Encode(map[string]any{"code_grants": f.grants, "refresh_grants": f.refreshGrants, "service_grants": f.serviceGrants, "jwks": f.jwks, "revocations": f.revocations, "active": active, "revoked": revoked, "reauth_required": reauth, "token_auth_valid": f.tokenAuthValid, "token_client_valid": f.tokenClientValid, "token_body_grant_valid": f.tokenBodyGrantValid})
		return
	}
	w.WriteHeader(404)
}
func TestPortalLifecycleOfficialSDKPostgres(t *testing.T) {
	dbURL := os.Getenv("PORTAL_SESSION_DB_URL")
	if dbURL == "" {
		t.Skip("explicit disposable portal PG fixture required")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	f := newPortalFixture(t, pool)
	directory := t.TempDir()
	doc, _ := json.Marshal(f.profile.Document())
	profile := filepath.Join(directory, "profile.json")
	if err = os.WriteFile(profile, doc, 0600); err != nil {
		t.Fatal(err)
	}
	secret := filepath.Join(directory, "client-secret")
	key := filepath.Join(directory, "session-key")
	if err = os.WriteFile(secret, []byte("synthetic-oidc-only"), 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(key, []byte("synthetic-portal-key-not-for-production-000000"), 0600); err != nil {
		t.Fatal(err)
	}
	node := os.Getenv("ELITE_NODE_EXACT")
	if node == "" {
		t.Fatal("pinned Node path required")
	}
	root, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.CommandContext(ctx, node, filepath.Join(root, "node_modules/vitest/vitest.mjs"), "run", "src/platform/auth/portal-lifecycle.connected.test.ts", "--reporter=verbose")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "OIDC_PORTAL_LIFECYCLE_ENABLED=true", "OIDC_PORTAL_PROFILE_FILE="+profile, "OIDC_PORTAL_PROFILE_SHA256="+f.profile.SHA256(), "OIDC_PORTAL_CLIENT_SECRET_FILE="+secret, "OIDC_PORTAL_SERVICE_CLIENT_SECRET_FILE="+secret, "OIDC_PORTAL_SESSION_KEY_FILE="+key, "PORTAL_FIXTURE_CONTROL_URL="+f.api.URL, "APP_BASE_URL=http://127.0.0.1:4567")
	output, err := cmd.CombinedOutput()
	t.Log(string(output))
	if err != nil {
		t.Fatal("real SDK BFF/PG fixture failed", err)
	}
	var plain int
	if err = pool.QueryRow(ctx, `select count(*) from platform.portal_session where profile_sha256=$1 and (ciphertext like '%fixture-refresh%' or ciphertext like '%accessToken%')`, f.profile.SHA256()).Scan(&plain); err != nil || plain != 0 {
		t.Fatal("plaintext token persisted", err)
	}
	t.Log("PORTAL_LIFECYCLE_SDK_PG_PASS actual_code_pkce=true actual_refresh_rotation=true durable_cas=true logout_replay_rejected=true plaintext_tokens=0")
}

func TestPortalSessionCASAndScope(t *testing.T) {
	dbURL := os.Getenv("PORTAL_SESSION_DB_URL")
	if dbURL == "" {
		t.Skip("explicit disposable portal PG fixture required")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	f := newPortalFixture(t, pool)
	store := &db.PortalSessions{Pool: pool}
	id := func() string {
		bytes := make([]byte, 32)
		_, _ = rand.Read(bytes)
		return base64.RawURLEncoding.EncodeToString(bytes)
	}
	sid := id()
	now := time.Now().UTC().Truncate(time.Millisecond)
	cipher := "header..iv.ciphertext-payload-fixture.tag"
	create := identity.PortalSessionCommand{SessionID: sid, Ciphertext: cipher, AccessExpiresAt: now.Add(time.Minute), AbsoluteExpiresAt: now.Add(time.Hour - time.Second)}
	row, err := store.Execute(ctx, f.profile, "create", create)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = store.Execute(ctx, f.profile, "create", create); err != nil {
		t.Fatal("create replay", err)
	}
	var wg sync.WaitGroup
	wins := make(chan string, 16)
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			op := id()
			if _, err := store.Execute(ctx, f.profile, "claim", identity.PortalSessionCommand{SessionID: sid, Version: row.Version, OperationID: op}); err == nil {
				wins <- op
			}
		}()
	}
	wg.Wait()
	close(wins)
	winner := ""
	count := 0
	for op := range wins {
		winner = op
		count++
	}
	if count != 1 {
		t.Fatal("multiple refresh owners", count)
	}
	if _, err = store.Execute(ctx, f.profile, "commit", identity.PortalSessionCommand{SessionID: sid, Version: 1, OperationID: id(), Ciphertext: cipher, AccessExpiresAt: now.Add(2 * time.Minute)}); err == nil {
		t.Fatal("foreign operation committed")
	}
	if _, err = store.Execute(ctx, f.profile, "revoke", identity.PortalSessionCommand{SessionID: sid}); err != nil {
		t.Fatal(err)
	}
	changed := cipher + "-rotated"
	commit := identity.PortalSessionCommand{SessionID: sid, Version: 1, OperationID: winner, Ciphertext: changed, AccessExpiresAt: now.Add(2 * time.Minute)}
	row, err = store.Execute(ctx, f.profile, "commit", commit)
	if err != nil || row.State != "revoked" || !row.RevocationPending || row.Version != 2 {
		t.Fatal("refresh reactivated logged out session", row.State, err)
	}
	if _, err = store.Execute(ctx, f.profile, "ack-revocation", identity.PortalSessionCommand{SessionID: sid, Version: 1, OperationID: winner}); err == nil {
		t.Fatal("stale revocation acknowledgement deleted new token")
	}
	if _, err = store.Execute(ctx, f.profile, "ack-revocation", identity.PortalSessionCommand{SessionID: sid, Version: 2, OperationID: winner}); err != nil {
		t.Fatal(err)
	}
	row, err = store.Execute(ctx, f.profile, "read", identity.PortalSessionCommand{SessionID: sid})
	if err != nil || row.Ciphertext != "" || row.State != "revoked" {
		t.Fatal("revocation cleanup", err)
	}
	d := f.profile.Document()
	d.ProfileID = "foreign-profile"
	raw, _ := json.Marshal(d)
	sum := sha256.Sum256(raw)
	foreign, err := identity.LoadPortalProfile(raw, hex.EncodeToString(sum[:]))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = store.Execute(ctx, foreign, "read", identity.PortalSessionCommand{SessionID: sid}); err != identity.ErrPortalSessionNotFound {
		t.Fatal("profile isolation", err)
	}
	t.Log("PORTAL_CAS_PASS concurrent_claims=16 grant_owner=1 logout_wins=true stale_ack_rejected=true profile_isolation=true")

	validClaims := func() map[string]any {
		return map[string]any{"iss": f.issuer.URL, "sub": "portal-service", "aud": "api", "iat": time.Now().Unix(), "exp": time.Now().Add(time.Minute).Unix(), "tenant_id": "tenant", "organization_ids": []string{"org"}, "permissions": []string{"portal-session:manage"}}
	}
	for _, tc := range []struct {
		name   string
		field  string
		value  any
		status int
	}{{"user", "sub", "person", 403}, {"tenant", "tenant_id", "foreign", 403}, {"organization", "organization_ids", []string{"foreign"}, 403}, {"wildcard", "permissions", []string{"*"}, 403}, {"extra_permission", "permissions", []string{"portal-session:manage", "admin:write"}, 403}, {"issuer", "iss", "https://foreign.invalid", 401}, {"audience", "aud", "foreign", 401}, {"expiry", "exp", time.Now().Add(-time.Minute).Unix(), 401}} {
		t.Run("http_"+tc.name, func(t *testing.T) {
			claims := validClaims()
			claims[tc.field] = tc.value
			request, _ := http.NewRequest("POST", f.api.URL+"/v1/private/portal-sessions/read", strings.NewReader(`{"session_id":"`+sid+`"}`))
			request.Header.Set("Authorization", "Bearer "+f.jwt(claims))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("X-Portal-Profile-SHA256", f.profile.SHA256())
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if response.StatusCode != tc.status {
				t.Fatal("authorization negative", response.StatusCode)
			}
		})
	}
	for _, tc := range []struct {
		name, body, profile string
		status              int
	}{{"profile", `{"session_id":"` + sid + `"}`, strings.Repeat("0", 64), 403}, {"duplicate", `{"session_id":"` + sid + `","session_id":"` + sid + `"}`, f.profile.SHA256(), 400}, {"unknown", `{"session_id":"` + sid + `","tenant_id":"foreign"}`, f.profile.SHA256(), 400}, {"too_large", strings.Repeat("x", 40001), f.profile.SHA256(), 413}} {
		t.Run("http_"+tc.name, func(t *testing.T) {
			request, _ := http.NewRequest("POST", f.api.URL+"/v1/private/portal-sessions/read", strings.NewReader(tc.body))
			request.Header.Set("Authorization", "Bearer "+f.jwt(validClaims()))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("X-Portal-Profile-SHA256", tc.profile)
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			_, _ = io.Copy(io.Discard, response.Body)
			if response.StatusCode != tc.status {
				t.Fatal("transport negative", response.StatusCode)
			}
		})
	}
	// Expiry while waiting on a row lock cannot be evaluated using pre-wait time.
	shortID := id()
	at := time.Now().UTC().Truncate(time.Millisecond)
	_, err = store.Execute(ctx, f.profile, "create", identity.PortalSessionCommand{SessionID: shortID, Ciphertext: cipher, AccessExpiresAt: at.Add(200 * time.Millisecond), AbsoluteExpiresAt: at.Add(350 * time.Millisecond)})
	if err != nil {
		t.Fatal(err)
	}
	lock, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	shortHash := sha256.Sum256([]byte(shortID))
	if _, err = lock.Exec(ctx, `select 1 from platform.portal_session where profile_sha256=$1 and session_id_sha256=$2 for update`, f.profile.SHA256(), shortHash[:]); err != nil {
		t.Fatal(err)
	}
	result := make(chan error, 1)
	go func() {
		_, e := store.Execute(ctx, f.profile, "claim", identity.PortalSessionCommand{SessionID: shortID, Version: 1, OperationID: id()})
		result <- e
	}()
	time.Sleep(450 * time.Millisecond)
	if err = lock.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	if err = <-result; err != identity.ErrPortalSessionConflict {
		t.Fatal("expired lock-wait accepted", err)
	}
}
