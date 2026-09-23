package whatsappbridge

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAIActivityContract(t *testing.T) {
	base := AIActivityItem{RequestID: "wa-reply:" + strings.Repeat("a", 64), State: "completed", ReviewState: "pending", Tool: "create_quote", Attempts: 1, ReservedTokens: "100", InputTokens: "20", OutputTokens: "9", OccurredAt: time.Now(), UpdatedAt: time.Now()}
	if !validAIActivity(base) {
		t.Fatal("valid rejected")
	}
	for name, mutate := range map[string]func(*AIActivityItem){"short_id": func(x *AIActivityItem) { x.RequestID = "" }, "unbound_state": func(x *AIActivityItem) { x.State = "processing" }, "sent_is_not_review": func(x *AIActivityItem) { x.ReviewState = "sent" }, "unknown_tool": func(x *AIActivityItem) { x.Tool = "delete_customer" }, "bad_units": func(x *AIActivityItem) { x.ReservedTokens = "-1" }, "overflow": func(x *AIActivityItem) { x.InputTokens = "9223372036854775808" }, "time": func(x *AIActivityItem) { x.UpdatedAt = x.OccurredAt.Add(-time.Second) }} {
		t.Run(name, func(t *testing.T) {
			x := base
			mutate(&x)
			if validAIActivity(x) {
				t.Fatal("invalid accepted")
			}
		})
	}
}
func TestAIActivityPermissionBeforeStorage(t *testing.T) {
	p := identity.Principal{TenantID: "tenant", Subject: "user", Organizations: map[string]struct{}{"store": {}}, Permissions: map[string]struct{}{"whatsapp:approve": {}}}
	m := &ReplyModule{approvals: &PostgresReplyApprovals{tenant: "tenant", organization: "store"}}
	mux := http.NewServeMux()
	m.Register(mux, replyVerifier{"human": p, "worker": {TenantID: "tenant", Subject: "worker", Organizations: p.Organizations, Permissions: map[string]struct{}{"whatsapp:process": {}}}})
	for _, x := range []struct {
		token, query string
		want         int
	}{{"", "", 401}, {"worker", "", 403}, {"human", "?tenant_id=other", 400}, {"human", "?organization_id=other", 400}} {
		t.Run(x.token+x.query, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/v1/franchise/whatsapp/ai-activity"+x.query, nil)
			if x.token != "" {
				r.Header.Set("Authorization", "Bearer "+x.token)
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			if w.Code != x.want || w.Header().Get("Cache-Control") != "no-store" {
				t.Fatalf("status %d", w.Code)
			}
		})
	}
	wrong := p
	wrong.Organizations = map[string]struct{}{"other": {}}
	if _, e := m.approvals.AIActivity(context.Background(), wrong); e != ErrApproval {
		t.Fatal("scope queried storage")
	}
}
