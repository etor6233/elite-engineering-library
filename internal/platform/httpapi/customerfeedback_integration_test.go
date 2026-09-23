package httpapi

import (
	"context"
	"elite.local/enterprise/internal/customerfeedback"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const feedbackTenant = "018f4d4a-7b36-7a21-8d10-2f4c54c24001"
const feedbackOtherTenant = "018f4d4a-7b36-7a21-8d10-2f4c54c24002"

func feedbackPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is required")
	}
	config, e := pgxpool.ParseConfig(url)
	if e != nil {
		t.Fatal(e)
	}
	if config.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(config.ConnConfig.Database, "elite_feedback_") {
		t.Fatal("requires an owned disposable loopback elite_feedback_* database")
	}
	p, e := pgxpool.NewWithConfig(context.Background(), config)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(p.Close)
	return p
}
func TestFeedbackConnectedPostgres(t *testing.T) {
	pool := feedbackPool(t)
	ctx := context.Background()
	exec := func(q string, args ...any) {
		t.Helper()
		if _, e := pool.Exec(ctx, q, args...); e != nil {
			t.Fatal(e)
		}
	}
	for _, tenant := range []string{feedbackTenant, feedbackOtherTenant} {
		exec(`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Feedback fixture','Feedback fixture')`, tenant, "feedback-"+tenant)
		exec(`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
  values($1,'org-a','a','A','store'),($1,'org-b','b','B','store')`, tenant)
		for _, customer := range []string{"alice", "bob", "carol", "inactive"} {
			state := "active"
			if customer == "inactive" {
				state = "restricted"
			}
			exec(`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized,status)
   values($1,$2,$2,$2||'@example.test',$3)`, tenant, customer, state)
		}
		exec(`insert into crm.survey_definition(tenant_id,organization_id,survey_id,prompt,consent_version,consent_notice,opens_at,closes_at,retain_until,minimum_responses,active)
  values($1,'org-a','satisfaction','Synthetic NPS question','fixture-v1','Synthetic fixture acknowledgement; no real collection',clock_timestamp()-interval '1 hour',clock_timestamp()+interval '1 hour',clock_timestamp()+interval '1 day',2,true)`, tenant)
	}
	verifier, token := confirmationTestIssuer(t)
	service := customerfeedback.NewService(postgres.NewCustomerFeedback(pool))
	mux := http.NewServeMux()
	CustomerFeedbackModule{Service: service}.Register(mux, verifier)
	var drop atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && drop.Swap(false) {
			record := httptest.NewRecorder()
			mux.ServeHTTP(record, r)
			if record.Code != 201 {
				t.Errorf("loss fixture status=%d", record.Code)
			}
			conn, _, e := w.(http.Hijacker).Hijack()
			if e != nil {
				t.Error(e)
				return
			}
			conn.Close()
			return
		}
		mux.ServeHTTP(w, r)
	}))
	defer server.Close()
	client := &http.Client{Timeout: 5 * time.Second}
	customerToken := func(tenant, subject string) string {
		return token(subject, tenant, []string{"customer:self"}, []string{"org-a"})
	}
	alice := customerToken(feedbackTenant, "alice")
	admin := token("operator", feedbackTenant, []string{"surveys:read"}, []string{"org-a"})
	path := "/v1/customer/surveys/satisfaction"
	send := func(method, path, bearer, body string) (int, []byte, error) {
		request, e := http.NewRequest(method, server.URL+path, strings.NewReader(body))
		if e != nil {
			return 0, nil, e
		}
		if bearer != "" {
			request.Header.Set("Authorization", "Bearer "+bearer)
		}
		if method == "POST" {
			request.Header.Set("Content-Type", "application/json")
		}
		response, e := client.Do(request)
		if e != nil {
			return 0, nil, e
		}
		defer response.Body.Close()
		data, e := io.ReadAll(response.Body)
		return response.StatusCode, data, e
	}
	must := func(method, p, bearer, body string, status int) []byte {
		t.Helper()
		code, data, e := send(method, p, bearer, body)
		if e != nil || code != status {
			t.Fatalf("%s %s: status=%d wanted=%d body=%s error=%v", method, p, code, status, data, e)
		}
		return data
	}
	answer := func(score int) string {
		return fmt.Sprintf(`{"score":%d,"consent":true,"consent_version":"fixture-v1"}`, score)
	}
	own := path + "/response?organization_id=org-a"
	must("GET", path+"?organization_id=org-a", "", "", 401)
	must("GET", path+"?organization_id=org-b", alice, "", 403)
	must("GET", path+"?organization_id=org-a&customer_id=bob", alice, "", 400)
	must("GET", path+"?organization_id=org-a", admin, "", 403)
	must("GET", path+"?organization_id=org-a", alice, "", 200)
	must("POST", own, alice, `{"score":9,"consent":false,"consent_version":"fixture-v1"}`, 400)
	must("POST", own, alice, `{"score":9,"consent":true,"consent_version":"wrong"}`, 404)
	must("GET", own, alice, "", 404)
	first := must("POST", own, alice, answer(9), 201)
	var saved customerfeedback.Result
	if e := json.Unmarshal(first, &saved); e != nil {
		t.Fatal(e)
	}
	replay := must("POST", own, alice, answer(9), 200)
	var repeated customerfeedback.Result
	json.Unmarshal(replay, &repeated)
	if !repeated.Replay || !saved.Answer.ReceivedAt.Equal(repeated.Answer.ReceivedAt) {
		t.Fatal("replay changed original")
	}
	must("POST", own, alice, answer(0), 409)
	summaryPath := "/v1/admin/surveys/satisfaction/summary?organization_id=org-a"
	var summary customerfeedback.Summary
	json.Unmarshal(must("GET", summaryPath, admin, "", 200), &summary)
	if summary.Responses != 1 || summary.Available || summary.NPS != nil {
		t.Fatal(summary)
	}
	bob := customerToken(feedbackTenant, "bob")
	results := make(chan int, 24)
	var workers sync.WaitGroup
	for i := 0; i < 24; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			code, _, e := send("POST", own, bob, answer(0))
			if e != nil {
				results <- 0
			} else {
				results <- code
			}
		}()
	}
	workers.Wait()
	close(results)
	created := 0
	for code := range results {
		if code == 201 {
			created++
		} else if code != 200 {
			t.Fatal("concurrent status", code)
		}
	}
	if created != 1 {
		t.Fatal("created", created)
	}
	drop.Store(true)
	if _, _, e := send("POST", own, customerToken(feedbackTenant, "carol"), answer(10)); e == nil {
		t.Fatal("lost response was not lost")
	}
	var recovered customerfeedback.Answer
	json.Unmarshal(must("GET", own, customerToken(feedbackTenant, "carol"), "", 200), &recovered)
	if recovered.Score != 10 {
		t.Fatal(recovered)
	}
	json.Unmarshal(must("GET", summaryPath, admin, "", 200), &summary)
	if summary.Responses != 3 || summary.NPS == nil || *summary.NPS < 33.3333333 || *summary.NPS > 33.3333334 {
		t.Fatal(summary)
	}
	must("GET", summaryPath, alice, "", 403)
	must("POST", own, customerToken(feedbackOtherTenant, "alice"), answer(5), 201)
	must("POST", own, customerToken(feedbackTenant, "inactive"), answer(7), 404)
	if _, e := pool.Exec(ctx, `update crm.survey_definition set prompt='changed' where tenant_id=$1`, feedbackTenant); e == nil {
		t.Fatal("mutable historical question")
	}
	exec(`update crm.survey_definition set active=false where tenant_id=$1`, feedbackTenant)
	must("POST", own, alice, answer(9), 200)
	t.Run("close_window_while_waiting_for_definition_lock", func(t *testing.T) {
		exec(`insert into crm.survey_definition(tenant_id,organization_id,survey_id,prompt,consent_version,consent_notice,opens_at,closes_at,retain_until,minimum_responses,active)
  values($1,'org-a','closing','Synthetic','fixture-v1','Synthetic',clock_timestamp()-interval '1 hour',clock_timestamp()+interval '1500 milliseconds',clock_timestamp()+interval '1 day',1,true)`, feedbackTenant)
		tx, e := pool.Begin(ctx)
		if e != nil {
			t.Fatal(e)
		}
		defer tx.Rollback(ctx)
		if _, e = tx.Exec(ctx, `select 1 from crm.survey_definition where tenant_id=$1 and survey_id='closing' for update`, feedbackTenant); e != nil {
			t.Fatal(e)
		}
		done := make(chan int, 1)
		go func() {
			code, _, _ := send("POST", "/v1/customer/surveys/closing/response?organization_id=org-a", alice, answer(9))
			done <- code
		}()
		deadline := time.Now().Add(1200 * time.Millisecond)
		waiting := false
		for time.Now().Before(deadline) {
			if e = pool.QueryRow(ctx, `select exists(select 1 from pg_stat_activity where datname=current_database() and wait_event_type='Lock' and query like '%crm.survey_definition d%')`).Scan(&waiting); e != nil {
				t.Fatal(e)
			}
			if waiting {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if !waiting {
			t.Fatal("request never waited on definition")
		}
		for {
			var closed bool
			if e = pool.QueryRow(ctx, `select clock_timestamp()>=closes_at from crm.survey_definition where tenant_id=$1 and survey_id='closing'`, feedbackTenant).Scan(&closed); e != nil {
				t.Fatal(e)
			}
			if closed {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if e = tx.Commit(ctx); e != nil {
			t.Fatal(e)
		}
		if code := <-done; code != 404 {
			t.Fatalf("expired while lock held: got %d; want404", code)
		}
		var n int
		if e = pool.QueryRow(ctx, `select count(*) from crm.survey_response where tenant_id=$1 and survey_id='closing'`, feedbackTenant).Scan(&n); e != nil || n != 0 {
			t.Fatal(n, e)
		}
	})
	t.Run("unpublished_definition_is_not_disclosed", func(t *testing.T) {
		_, e := pool.Exec(ctx, `insert into crm.survey_definition(tenant_id,organization_id,survey_id,prompt,consent_version,consent_notice,opens_at,closes_at,retain_until,minimum_responses)values($1,'org-a','draft','Unpublished','v1','Synthetic',clock_timestamp()-interval '1 hour',clock_timestamp()+interval '1 hour',clock_timestamp()+interval '1 day',2)`, feedbackTenant)
		if e != nil {
			t.Fatal(e)
		}
		repository := postgres.NewCustomerFeedback(pool)
		s := customerfeedback.Scope{Tenant: feedbackTenant, Organization: "org-a", Customer: "alice"}
		if _, e = repository.Definition(ctx, s, "draft"); e != customerfeedback.ErrUnavailable {
			t.Fatalf("draft exposed: %v", e)
		}
		if d, e := repository.Definition(ctx, s, "satisfaction"); e != nil || d.Accepting {
			t.Fatalf("closed receipt context lost: %+v %v", d, e)
		}
	})
	t.Log("FEEDBACK_CONNECTED_HTTP_OIDC_POSTGRES: identity, consent, replay,24concurrent requests, lost response, aggregation and close-window fence checked")
}
func TestFeedbackAfterDatabaseRestart(t *testing.T) {
	pool := feedbackPool(t)
	service := customerfeedback.NewService(postgres.NewCustomerFeedback(pool))
	ctx := context.Background()
	s := customerfeedback.Scope{Tenant: feedbackTenant, Organization: "org-a", Customer: "carol"}
	got, e := service.OwnAnswer(ctx, s, "satisfaction")
	if e != nil || got.Score != 10 {
		t.Fatal(got, e)
	}
	stats, e := service.Summary(ctx, s, "satisfaction")
	if e != nil || stats.Responses != 3 || stats.NPS == nil {
		t.Fatal(stats, e)
	}
	for _, customer := range []string{"alice", "bob", "carol"} {
		s.Customer = customer
		if _, e := service.OwnAnswer(ctx, s, "satisfaction"); e != nil {
			t.Fatal(e)
		}
	}
	t.Log("FEEDBACK_POSTGRES_RESTART_PASS responses=3 no_mutations=true")
}
