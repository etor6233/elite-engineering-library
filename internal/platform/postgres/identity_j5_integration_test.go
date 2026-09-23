package postgres_test

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	nr "elite.local/enterprise/internal/networkrole"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func j5RoleIssuer(t *testing.T) (identity.Verifier, func(string, string, []string, []string, bool) string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	encode := base64.RawURLEncoding.EncodeToString
	var issuer *httptest.Server
	issuer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			_ = json.NewEncoder(w).Encode(map[string]any{"issuer": issuer.URL, "jwks_uri": issuer.URL + "/keys", "id_token_signing_alg_values_supported": []string{"RS256"}})
		case "/keys":
			_ = json.NewEncoder(w).Encode(map[string]any{"keys": []any{map[string]any{"kty": "RSA", "use": "sig", "alg": "RS256", "kid": "local-confirmation", "n": encode(key.N.Bytes()), "e": encode(big.NewInt(int64(key.E)).Bytes())}}})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(issuer.Close)
	verifier, err := identity.NewOIDCVerifier(context.Background(), issuer.URL, "confirmation-api")
	if err != nil {
		t.Fatal(err)
	}
	token := func(subject, tenant string, permissions, organizations []string, expired bool) string {
		header, err := json.Marshal(map[string]string{"alg": "RS256", "kid": "local-confirmation", "typ": "JWT"})
		if err != nil {
			t.Fatal(err)
		}
		issued, expires := time.Now(), time.Now().Add(120*time.Second)
		if expired {
			issued = time.Now().Add(-300 * time.Second)
			expires = time.Now().Add(-120 * time.Second)
		}
		claims, err := json.Marshal(map[string]any{"iss": issuer.URL, "aud": "confirmation-api", "sub": subject, "iat": issued.Unix(), "exp": expires.Unix(), "tenant_id": tenant, "permissions": permissions, "organization_ids": organizations})
		if err != nil {
			t.Fatal(err)
		}
		unsigned := encode(header) + "." + encode(claims)
		digest := sha256.Sum256([]byte(unsigned))
		signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
		if err != nil {
			t.Fatal(err)
		}
		return unsigned + "." + encode(signature)
	}
	return verifier, token
}

func TestJ5IdentityScopedBootstrap(t *testing.T) {
	dsn := os.Getenv("PORTAL_SESSION_DB_URL")
	if dsn == "" {
		t.Skip("explicit owned PG fixture required")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := uuid.NewString()
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,$2,'Synthetic','Synthetic')`, tenant, "j5-"+tenant); err != nil {
		t.Fatal(err)
	}
	verifier, mint := j5RoleIssuer(t)
	store, err := db.NewNetworkRole(pool)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	httpapi.NetworkRoleModule{Service: store}.Register(mux, verifier)
	api := httptest.NewServer(mux)
	defer api.Close()
	post := func(token string, command nr.Command, status int) nr.Receipt {
		t.Helper()
		raw, _ := json.Marshal(command)
		request, _ := http.NewRequest("POST", api.URL+"/v1/franchise/network/commands", bytes.NewReader(raw))
		request.Header.Set("Authorization", "Bearer "+token)
		request.Header.Set("Content-Type", "application/json")
		response, err := api.Client().Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		if response.StatusCode != status {
			t.Fatalf("%s expected%d got%d %s", command.CommandID, status, response.StatusCode, body)
		}
		var receipt nr.Receipt
		if status == 200 && json.Unmarshal(body, &receipt) != nil {
			t.Fatal("receipt")
		}
		return receipt
	}
	get := func(token, path string, status int) {
		t.Helper()
		request, _ := http.NewRequest("GET", api.URL+path, nil)
		request.Header.Set("Authorization", "Bearer "+token)
		response, err := api.Client().Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
		if response.StatusCode != status {
			t.Fatal("read authority", status, response.StatusCode)
		}
	}
	admin := mint("initial-admin", tenant, []string{"network:admin", "network:bootstrap"}, []string{"root"}, false)
	root := nr.Command{CommandID: "j5-root-create", Action: "create-organization", EntityID: "root", Code: "root", DisplayName: "Synthetic root", Type: "franchisor"}
	post(mint("unprivileged", tenant, []string{"network:bootstrap"}, []string{"root"}, false), root, 403)
	post(mint("normal-admin", tenant, []string{"network:admin"}, []string{"root"}, false), root, 403)
	receipt := post(admin, root, 200)
	if receipt.Actor != "initial-admin" || receipt.Entity.ID != "root" {
		t.Fatal("issuer actor not retained")
	}
	current := mint("initial-admin", tenant, []string{"network:admin"}, []string{"root"}, false)
	activation := nr.Command{CommandID: "j5-root-active", Action: "transition-organization", ScopeOrganizationID: "root", EntityID: "root", Current: "provisioning", Target: "active", Version: "1"}
	post(current, activation, 200)
	get(current, "/v1/franchise/network/entities/organization/root?scope_organization_id=root", 200)
	other := root
	other.CommandID = "forbidden-extra-root"
	other.EntityID = "other"
	other.Code = "other"
	post(current, other, 403)
	scoped := nr.Command{CommandID: "j5-branch-create", Action: "create-organization", ScopeOrganizationID: "root", EntityID: "branch", Code: "branch", DisplayName: "Synthetic branch", Type: "store"}
	post(current, scoped, 200)
	foreignScope := mint("initial-admin", tenant, []string{"network:admin"}, []string{"foreign"}, false)
	post(foreignScope, activation, 403)
	withdrawn := mint("initial-admin", tenant, []string{}, []string{"root"}, false)
	post(withdrawn, activation, 403)
	get(withdrawn, "/v1/franchise/network/entities/organization/root?scope_organization_id=root", 403)
	expired := mint("initial-admin", tenant, []string{"network:admin"}, []string{"root"}, true)
	post(expired, activation, 401)
	post(current+"bad", activation, 401)
	get(mint("other-admin", tenant, []string{"network:admin"}, []string{"root"}, false), "/v1/franchise/network/commands/j5-root-active?scope_organization_id=root", 404)
	get(mint("initial-admin", uuid.NewString(), []string{"network:admin", "network:bootstrap"}, []string{"root"}, false), "/v1/franchise/network/commands/j5-root-create", 404)
	var organizations, receipts, events int
	err = pool.QueryRow(ctx, `select (select count(*)from org.organization where tenant_id=$1),(select count(*)from franchise.network_command_receipt where tenant_id=$1),(select count(*)from platform.outbox_event where tenant_id=$1)`, tenant).Scan(&organizations, &receipts, &events)
	if err != nil || organizations != 2 || receipts != 3 || events != 3 {
		t.Fatal("denial mutated state", organizations, receipts, events, err)
	}
	t.Log("J5_IDENTITY_PASS signed_OIDC=true wildcard=false bootstrap_scoped=true explicit_role_removal=true expiry_and_tamper_rejected=true tenant_actor_isolated=true organizations=2 immutable_receipts=3 outbox=3; existing bearer remains valid until fixed expiry, no immediate live IdP revocation claim")
}
