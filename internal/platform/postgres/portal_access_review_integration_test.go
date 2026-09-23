package postgres_test

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestPortalAccessReviewOfficialSDKPostgres(t *testing.T) {
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
	cmd := exec.CommandContext(ctx, node, filepath.Join(root, "node_modules/vitest/vitest.mjs"), "run", "src/platform/auth/portal-access-review.connected.test.ts", "--reporter=verbose")
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
	t.Log("PORTAL_ACCESS_REVIEW_PASS official_sdk_refresh=true permission_withdrawal=true changed_subject_rejected=true account_disabled_reauth=true plaintext_tokens=0")
}
