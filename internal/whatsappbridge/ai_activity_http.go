package whatsappbridge

// AUTHORED read-only endpoint; registration follows the already selected and
// activated WhatsApp reply owner. No automatic provider/model calls.
import (
	"elite.local/enterprise/internal/platform/identity"
	"net/http"
)

func (m *ReplyModule) registerAIActivity(mux *http.ServeMux, v identity.Verifier) {
	mux.HandleFunc("GET /v1/franchise/whatsapp/ai-activity", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		p, ok := m.principal(w, r, v, "whatsapp:approve")
		if !ok {
			return
		}
		if r.URL.RawQuery != "" {
			notificationProblem(w, 400, "INVALID_QUERY")
			return
		}
		value, e := m.approvals.AIActivity(r.Context(), p)
		if e != nil {
			notificationProblem(w, 503, "UNAVAILABLE")
			return
		}
		replyJSON(w, value)
	})
}
