package whatsappbridge

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

// Native fuzz covers the bounded typed projection AFTER the separate signed
// Python trust boundary. This does not claim signature verification or DAST.
func FuzzWhatsAppConversationProjection(f *testing.F) {
	seed := fmt.Sprintf(`{"object":"whatsapp_business_account","entry":[{"id":"987654321","changes":[{"field":"messages","value":{"messaging_product":"whatsapp","metadata":{"phone_number_id":"123456789"},"messages":[{"id":"wamid.fixture","from":"5491112345678","timestamp":"%d","type":"text","text":{"body":"Cotizame un scooter"}}]}}]}]}`, time.Now().Add(-time.Minute).Unix())
	for _, v := range []string{seed, `{}`, `null`, strings.Replace(seed, `"type":"text"`, `"type":"image"`, 1), strings.Replace(seed, `"body":"Cotizame un scooter"`, `"body":null`, 1)} {
		f.Add([]byte(v))
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 16384 {
			return
		}
		messages, total, err := verifiedConversationMessages(raw, "fixture-tenant", "fixture-connection", 8)
		if err != nil {
			return
		}
		if len(messages) > 8 || total < len(messages) {
			t.Fatal("projection escaped route budget")
		}
		seen := map[string]bool{}
		for _, m := range messages {
			if m.TenantID != "fixture-tenant" || m.ChannelCode != "whatsapp" || m.ThreadID != conversationThread("fixture-connection") || m.Direction != "in" || !utf8.ValidString(m.Text) || m.Validate() != nil || seen[m.ProviderMessageID] {
				t.Fatal("projection escaped source scope or identity")
			}
			seen[m.ProviderMessageID] = true
			encoded, err := json.Marshal(m)
			if err != nil {
				t.Fatal(err)
			}
			var again map[string]any
			if json.Unmarshal(encoded, &again) != nil || again["Text"] != m.Text || again["ProviderMessageID"] != m.ProviderMessageID {
				t.Fatal("projection JSON changed visible text/identity")
			}
		}
	})
}
