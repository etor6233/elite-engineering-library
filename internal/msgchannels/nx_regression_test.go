package msgchannels

import (
	"strings"
	"testing"
)

func TestNXWhatsAppWebhookShape(t *testing.T) {
	r := []byte(`{"object":"whatsapp_business_account","entry":[{"changes":[{"field":"messages","value":{"messaging_product":"whatsapp","metadata":{"phone_number_id":"p1"},"messages":[{"from":"u1","id":"wamid.1","type":"text","text":{"body":"hola"}}]}}]}]}`)
	m, e := ParseWebhook(r)
	if e != nil || len(m) != 1 || m[0].Channel != "whatsapp" || m[0].Text != "hola" {
		t.Fatalf("WA payload lost: %v %v", m, e)
	}
}
func TestNXUnknownWebhookEnvelopeRejected(t *testing.T) {
	for _, s := range []string{`null`, `{}`, `{"object":"unknown","entry":[]}`, `{"object":"page","entry":[]}` + strings.Repeat(" ", 1048577)} {
		if _, err := ParseWebhook([]byte(s)); err == nil {
			t.Error("unknown/unbounded envelope accepted")
		}
	}
}
