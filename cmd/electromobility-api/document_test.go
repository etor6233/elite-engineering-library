package main

import (
	"context"
	doc "elite.local/enterprise/internal/documentbridge"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	urlpkg "net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDocumentHostActivation(t *testing.T) {
	ctx := context.Background()
	none := func(string) string { return "" }
	if m, e := selectedDocumentModule(ctx, nil, none); e != nil || m != nil {
		t.Fatal("default activation")
	}
	url := os.Getenv("DOCUMENT_CONNECTED_DB_URL")
	if url == "" {
		t.Skip("owned document fixture database required")
	}
	parsed, parseErr := urlpkg.Parse(url)
	if parseErr != nil || (parsed.Hostname() != "127.0.0.1" && parsed.Hostname() != "localhost") || (!strings.HasPrefix(parsed.Path, "/elite_document_reference_") && !strings.HasPrefix(parsed.Path, "/elite_payment_connected_")) {
		t.Fatal("owned loopback fixture database required")
	}
	pool, e := pgxpool.New(ctx, url)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	profile := filepath.Join(os.Getenv("DOCUMENT_TEST_ROOT"), "config", "documents", "reference-profile.json")
	b, e := os.ReadFile(profile)
	if e != nil {
		t.Fatal(e)
	}
	env := map[string]string{"DOCUMENTS_ENABLED": "true", "DOCUMENTS_MODE": "FIXTURE", "DOCUMENTS_PROFILE_FILE": profile, "DOCUMENTS_PROFILE_SHA256": doc.Hash(b), "DOCUMENTS_WORK_ROOT": t.TempDir(), "DOCUMENTS_TENANT_ID": uuid.NewString(), "DOCUMENTS_ORGANIZATION_ID": "fixture-org"}
	get := func(k string) string { return env[k] }
	if m, e := selectedDocumentModule(ctx, pool, get); e != nil || m == nil {
		t.Fatal("selected fixture module unavailable", e)
	}
	env["DOCUMENTS_MODE"] = "PROVIDER"
	if _, e = selectedDocumentModule(ctx, pool, get); e == nil {
		t.Fatal("provider mode accepted missing security runtime")
	}
	env["DOCUMENTS_MODE"] = "FIXTURE"
	env["DOCUMENTS_PROFILE_SHA256"] = "bad"
	if _, e = selectedDocumentModule(ctx, pool, get); e == nil {
		t.Fatal("unbound profile")
	}
	env["DOCUMENTS_PROFILE_SHA256"] = doc.Hash(b)
	_, e = pool.Exec(ctx, `alter table document.committed disable trigger document_commit_guard`)
	if e != nil {
		t.Fatal(e)
	}
	_, e = selectedDocumentModule(ctx, pool, get)
	_, restore := pool.Exec(ctx, `alter table document.committed enable trigger document_commit_guard`)
	if restore != nil {
		t.Fatal(restore)
	}
	if e == nil {
		t.Fatal("host accepted absent approval/commit guard")
	}
	t.Log("DOCUMENT_HOST_PASS: opt-in fixture module; provider runtime required; profile hash; database guard refusal")
}
