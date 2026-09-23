package postgres_test

// AUTHORED connected test: current PG publication -> actual Next -> Chromium.
// Synthetic loopback only; the existing HTTP journey owns all catalog writes.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func catalogStorefrontFixtureRequested() bool { return os.Getenv("ELITE_CATALOG_STOREFRONT") == "1" }
func TestCatalogReleaseStorefrontReference(t *testing.T) {
	if !catalogStorefrontFixtureRequested() {
		t.Skip("explicit catalog storefront fixture required")
	}
	testCatalogReleaseHTTPReference(t, true)
}
func attachCatalogStorefrontFixture(t *testing.T, c *catalogHTTPFixture) {
	t.Helper()
	web, node := os.Getenv("ELITE_WEB_ROOT"), os.Getenv("ELITE_NODE_BIN")
	if !filepath.IsAbs(web) || !filepath.IsAbs(node) {
		t.Fatal("absolute exact web and Node paths required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	t.Cleanup(cancel)
	listener, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	address := listener.Addr().String()
	listener.Close()
	base := "http://" + address
	env := []string{}
	for _, v := range os.Environ() {
		k := strings.ToUpper(strings.SplitN(v, "=", 2)[0])
		if strings.HasPrefix(k, "ELITE_") || strings.HasPrefix(k, "CATALOG_") || strings.HasPrefix(k, "PUBLIC_") ||
			strings.HasPrefix(k, "ENTERPRISE_") || k == "DATABASE_URL" || k == "TEST_DATABASE_URL" || k == "PAYMENT_CONNECTED_DB_URL" ||
			k == "AUTH_SESSION_SECRET" || k == "APP_BASE_URL" {
			continue
		}
		env = append(env, v)
	}
	env = append(env, "CATALOG_RELEASE_ENABLED=true", "PUBLIC_SITE_ORIGIN=https://catalog.example.invalid", "PUBLIC_INDEXING_ENABLED=1",
		"ENTERPRISE_TENANT_CODE=catalog-j3", "ENTERPRISE_ORGANIZATION_CODE=j3-store", "ENTERPRISE_API_BASE_URL="+c.base,
		"APP_BASE_URL="+base, "AUTH_SESSION_SECRET=synthetic-catalog-next-reference-000000000000000",
		"ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+base, "ELITE_CATALOG_STOREFRONT=1", "NEXT_TELEMETRY_DISABLED=1", "BUSINESS_CONFIG_FILE=business.example.json")
	artifacts, e := os.MkdirTemp(web, "catalog-storefront-artifacts-")
	if e != nil {
		t.Fatal(e)
	}
	log, e := os.Create(filepath.Join(artifacts, "next.log"))
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { log.Close() })
	server := exec.CommandContext(ctx, node, filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", strings.Split(address, ":")[1])
	server.Dir, server.Env, server.Stdout, server.Stderr = web, env, log, log
	if e = server.Start(); e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { server.Process.Kill(); server.Wait() })
	client := &http.Client{Timeout: time.Second}
	ready := false
	for deadline := time.Now().Add(25 * time.Second); time.Now().Before(deadline); {
		r, e := client.Get(base + "/icon.svg")
		if e == nil {
			r.Body.Close()
			if r.StatusCode == 200 {
				ready = true
				break
			}
		}
		select {
		case <-ctx.Done():
			t.Fatal("startup cancelled")
		case <-time.After(100 * time.Millisecond):
		}
	}
	if !ready {
		t.Fatal("Next startup failed; see ", artifacts)
	}
	c.webProbe = func(document cr.PublicDocument) {
		raw, e := json.Marshal(document)
		if e != nil {
			t.Fatal(e)
		}
		generation := filepath.Join(artifacts, "generation-"+jsonNumber(document.Generation))
		if e = os.Mkdir(generation, 0700); e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(filepath.Join(generation, "publication.json"), raw, 0600); e != nil {
			t.Fatal(e)
		}
		gate := filepath.Join(web, "microsoft_playwright_browser_gate")
		command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/catalog-publication-connected.spec.mjs", "--project=chromium-desktop", "--timeout=30000", "--output="+filepath.Join(generation, "browser"))
		command.Dir, command.Env = gate, append(append([]string(nil), env...), "ELITE_CATALOG_EXPECTED="+string(raw))
		output, e := command.CombinedOutput()
		if x := os.WriteFile(filepath.Join(generation, "browser.log"), output, 0600); x != nil {
			t.Fatal(x)
		}
		if e != nil || !strings.Contains(string(output), "1 passed") {
			t.Fatalf("storefront generation%d: %v\n%s\nartifacts=%s", document.Generation, e, output, artifacts)
		}
	}
	t.Cleanup(func() {
		t.Log("CATALOG_STOREFRONT_PASS actual Next/Chromium/PG; published versions and rollback, canonical, sitemap, robots, normalized PNG and absent-model404; artifacts=", artifacts)
	})
}
func jsonNumber(value int64) string { raw, _ := json.Marshal(value); return string(raw) }
