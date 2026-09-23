package whatsappbridge

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestAIActivityConnectedProjection(t *testing.T) {
	f := newConnectedReplyFixture(t)
	ctx := context.Background()
	read := func() AIActivity {
		t.Helper()
		req, _ := http.NewRequest("GET", f.api.URL+"/v1/franchise/whatsapp/ai-activity", nil)
		req.Header.Set("Authorization", "Bearer human")
		res, e := f.api.Client().Do(req)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		raw, _ := io.ReadAll(res.Body)
		if res.StatusCode != 200 {
			t.Fatalf("read status %d %s", res.StatusCode, raw)
		}
		for _, forbidden := range []string{"5491112345678", "Cotización preparada", "fixture-human", "synthetic", "tenant_id", "external_id", "user_text", "assistant_text", "tool_arguments"} {
			if strings.Contains(string(raw), forbidden) {
				t.Fatal("private or execution field emitted", forbidden)
			}
		}
		var out AIActivity
		if json.Unmarshal(raw, &out) != nil {
			t.Fatal("invalid projection")
		}
		return out
	}
	if len(read().Items) != 0 {
		t.Fatal("initial not empty")
	}
	f.ingest(t, f.inbound(t, "wamid.ai-activity", time.Now().Add(-time.Second), "5491112345678"))
	f.process(t)
	out := read()
	if len(out.Items) != 1 || out.Items[0].State != "completed" || out.Items[0].ReviewState != "pending" || out.Items[0].Tool != "create_quote" || out.Items[0].ReservedTokens == "0" || !out.Items[0].IntentPinned {
		t.Fatalf("projection %+v", out)
	}
	p := f.proposal(t, "wamid.ai-activity")
	// Fixture preload represents an exhausted internal allowance. Keep the cap
	// configuration unchanged: cap mismatch is retryable, exhaustion hands off.
	// This tests the real runtime handoff; it does not claim this preload was
	// spent by a provider or produced by the runtime.
	if _, e := f.pool.Exec(ctx, `update communication.conversation_budget set reserved=cap where tenant_id=$1`, f.service.TenantID); e != nil {
		t.Fatal(e)
	}
	f.ingest(t, f.inbound(t, "wamid.ai-handoff", time.Now().Add(-time.Second), "5491112345678"))
	f.process(t)
	out = read()
	handoff := false
	for _, x := range out.Items {
		if x.State == "handed_off" && x.ReservedTokens == "0" {
			handoff = true
		}
	}
	if len(out.Items) != 2 || !handoff || f.llmCalls.Load() != 2 || f.domainCalls.Load() != 1 || f.metaCalls.Load() != 0 {
		t.Fatal("handoff/read caused effects")
	}
	for i := 0; i < 3; i++ {
		read()
	}
	if f.llmCalls.Load() != 2 || f.domainCalls.Load() != 1 || f.metaCalls.Load() != 0 {
		t.Fatal("GET invoked a provider")
	}
	if f.command(t, p.RequestID, "send", "human", map[string]any{"payload_sha256": p.PayloadSHA256}) != 409 {
		t.Fatal("unapproved send")
	}
	if f.command(t, p.RequestID, "decision", "human", map[string]any{"payload_sha256": p.PayloadSHA256, "approved": true, "reason": "Fixture review"}) != 200 {
		t.Fatal("decision")
	}
	out = read()
	found := false
	for _, x := range out.Items {
		if x.RequestID == p.RequestID && x.ReviewState == "approved" {
			found = true
		}
	}
	if !found || f.metaCalls.Load() != 0 {
		t.Fatal("approval misrepresented")
	}
	for i := 0; i < 2; i++ {
		if f.command(t, p.RequestID, "send", "human", map[string]any{"payload_sha256": p.PayloadSHA256}) != 200 {
			t.Fatal("send")
		}
	}
	if f.metaCalls.Load() != 1 {
		t.Fatal("duplicate provider effect")
	}
	before := read()
	wrong := f.human
	wrong.TenantID = "00000000-0000-4000-8000-000000000000"
	if _, e := f.module.approvals.AIActivity(ctx, wrong); e != ErrApproval {
		t.Fatal("tenant leak")
	}
	wrong = f.human
	wrong.Organizations = map[string]struct{}{"other": {}}
	if _, e := f.module.approvals.AIActivity(ctx, wrong); e != ErrApproval {
		t.Fatal("org leak")
	}
	// Approval owner retains immutable org authority even if a connection changes.
	if _, e := f.pool.Exec(ctx, `update integration.provider_connection set state='disabled' where tenant_id=$1 and connection_id='wa-primary'`, f.service.TenantID); e != nil {
		t.Fatal(e)
	}
	after := read()
	if len(after.Items) != len(before.Items) {
		t.Fatal("historical outcomes hidden by deactivation")
	}
	t.Log("AI_ACTIVITY_CONNECTED_PASS runtime proposal + budget handoff + exact scoped metadata + no GET effects + existing human approval/send idempotency")
}
