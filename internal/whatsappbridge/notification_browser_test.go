package whatsappbridge

// AUTHORED opt-in browser fixture. Real Go/PG/Next/TLS; synthetic issuer,
// approval, receipt and provider webhook. Never a live Meta or IdP assertion.
import (
	"context"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type notificationClock struct{}

func (notificationClock) Now() time.Time { return time.Now() }

func TestNotificationStatusBrowserPostgres(t *testing.T) {
	if os.Getenv("ELITE_NOTIFICATION_BROWSER_E2E") != "1" {
		t.Skip("explicit isolated notification browser gate not selected")
	}
	web := os.Getenv("ELITE_WEB_ROOT")
	if !filepath.IsAbs(web) {
		t.Fatal("absolute materialized web root required")
	}
	for _, relative := range []string{"node_modules/next/dist/bin/next", ".next/BUILD_ID", "microsoft_playwright_browser_gate/node_modules/@playwright/test/cli.js"} {
		if _, err := os.Stat(filepath.Join(web, relative)); err != nil {
			t.Fatal("browser prerequisite missing; install gate with --ignore-workspace and build exact web first: " + relative)
		}
	}
	for _, project := range []string{"chromium-desktop", "chromium-mobile", "firefox-desktop", "webkit-desktop"} {
		t.Run(project, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			r, o, pool, p, c, _, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
			w, err := NewStatusWorker(r, "browser-worker", time.Minute, time.Second)
			if err != nil {
				t.Fatal(err)
			}
			if result, err := w.ProcessOnce(ctx, p); err != nil || !result.Completed {
				t.Fatal("worker prerequisite", err, result)
			}
			verifier, sign := notificationIssuer(t)
			module := &AppointmentNotificationModule{approvals: o.Approvals, sender: &Sender{TenantID: p.TenantID}}
			mux := http.NewServeMux()
			module.Register(mux, verifier)
			httpapi.FranchiseJourneyModule{Service: franchisejourney.NewService(postgres.NewFranchiseJourney(pool), randomid.Generator{}, notificationClock{})}.Register(mux, verifier)
			api := httptest.NewServer(mux)
			defer api.Close()
			target, _ := url.Parse("http://127.0.0.1:4173")
			edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
			defer edge.Close()
			var start time.Time
			if err := pool.QueryRow(ctx, `select starts_at from crm.appointment where tenant_id=$1 and appointment_id=$2`, p.TenantID, c.AppointmentID).Scan(&start); err != nil {
				t.Fatal(err)
			}
			env := []string{}
			for _, entry := range os.Environ() {
				name := strings.ToUpper(strings.SplitN(entry, "=", 2)[0])
				if name == "BUSINESS_CONFIG_FILE" || name == "TEST_DATABASE_URL" || name == "DATABASE_URL" || name == "ELITE_WHATSAPP_TEST_DATABASE_URL" || strings.HasPrefix(name, "ELITE_NOTIFICATION_") || strings.HasPrefix(name, "ELITE_CONFIRMATION_") || name == "ENTERPRISE_API_BASE_URL" || name == "AUTH_SESSION_SECRET" || name == "APP_BASE_URL" || name == "ELITE_BASE_URL" || name == "ELITE_RUNTIME_ONLY" || name == "ELITE_OPERATOR_AGENDA_E2E" || name == "ELITE_WEB_ROOT" {
					continue
				}
				env = append(env, entry)
			}
			token := sign(p)
			secret := "synthetic-" + (randomid.Generator{}).New() + (randomid.Generator{}).New()
			// Use the existing business configuration owner, never a new runtime flag.
			rawConfig, err := os.ReadFile(filepath.Join(web, "config", "business.example.json"))
			if err != nil {
				t.Fatal(err)
			}
			var config map[string]any
			if err := json.Unmarshal(rawConfig, &config); err != nil {
				t.Fatal(err)
			}
			features, ok := config["features"].(map[string]any)
			if !ok {
				t.Fatal("business fixture features missing")
			}
			features["whatsapp_status_history"] = true
			configBytes, err := json.Marshal(config)
			if err != nil {
				t.Fatal(err)
			}
			configName := "notification-fixture-" + (randomid.Generator{}).New() + ".json"
			configPath := filepath.Join(web, "config", configName)
			configFile, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := os.Remove(configPath); err != nil {
					t.Error(err)
				}
			})
			_, writeErr := configFile.Write(configBytes)
			closeErr := configFile.Close()
			if writeErr != nil || closeErr != nil {
				t.Fatal("business fixture could not be persisted")
			}
			env = append(env, "BUSINESS_CONFIG_FILE="+configName)
			env = append(env, "ENTERPRISE_API_BASE_URL="+api.URL, "APP_BASE_URL="+edge.URL, "ELITE_BASE_URL="+edge.URL, "ELITE_WEB_ROOT="+web, "ELITE_NOTIFICATION_BROWSER_E2E=1", "AUTH_SESSION_SECRET="+secret, "ELITE_NOTIFICATION_TOKEN="+token, "ELITE_NOTIFICATION_TENANT="+p.TenantID, "ELITE_NOTIFICATION_ORGANIZATION="+c.OrganizationID, "ELITE_NOTIFICATION_APPOINTMENT="+c.AppointmentID, "ELITE_NOTIFICATION_EVENT="+c.ConfirmationEventID, "ELITE_NOTIFICATION_DAY="+start.UTC().Format("2006-01-02"))
			server := exec.CommandContext(ctx, "node", filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", "4173")
			server.Dir = web
			server.Env = env
			if err := server.Start(); err != nil {
				t.Fatal(err)
			}
			defer func() { _ = server.Process.Kill(); _ = server.Wait() }()
			ready := false
			client := &http.Client{Timeout: time.Second}
			for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
				response, err := client.Get(target.String() + "/icon.svg")
				if err == nil {
					response.Body.Close()
					if response.StatusCode == 200 {
						ready = true
						break
					}
				}
				select {
				case <-ctx.Done():
					t.Fatal("web startup cancelled")
				case <-time.After(100 * time.Millisecond):
				}
			}
			if !ready {
				t.Fatal("web fixture did not start")
			}
			cmd := exec.CommandContext(ctx, "node", filepath.Join(web, "microsoft_playwright_browser_gate", "node_modules", "@playwright", "test", "cli.js"), "test", "tests/notification-status.spec.mjs", "--project", project, "--workers=1", "--retries=0", "--max-failures=1", "--output", filepath.Join(web, "microsoft_playwright_browser_gate", "test-results", "notifications-"+project))
			cmd.Dir = filepath.Join(web, "microsoft_playwright_browser_gate")
			cmd.Env = env
			output, err := cmd.CombinedOutput()
			clean := strings.ReplaceAll(strings.ReplaceAll(string(output), token, "[REDACTED_TEST_TOKEN]"), secret, "[REDACTED_TEST_SECRET]")
			t.Log(clean)
			if err != nil {
				t.Fatal("notification browser gate", err)
			}
			var jobs, attempts, observations int
			if err := pool.QueryRow(ctx, `select (select count(*) from platform.job where tenant_id=$1 and completed_at is not null),(select sum(attempt_count)::int from communication.outbound_delivery where tenant_id=$1),(select count(*) from communication.whatsapp_status_observation where tenant_id=$1)`, p.TenantID).Scan(&jobs, &attempts, &observations); err != nil || jobs != 1 || attempts != 1 || observations != 1 {
				t.Fatal("read changed durable effects", err, jobs, attempts, observations)
			}
			t.Log("NOTIFICATION_BROWSER_POSTGRES_PASS status=observed_delivered jobs=1 send_attempts=1 observations=1 no_resend=true")
		})
	}
}
