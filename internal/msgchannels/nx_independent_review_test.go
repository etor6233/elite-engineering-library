package msgchannels

import (
	"testing"
)

const reviewWAText = `{"object":"whatsapp_business_account","entry":[{"changes":[{"field":"messages","value":{"messaging_product":"whatsapp","metadata":{"phone_number_id":"p1"},"messages":[{"from":"u1","id":"wamid.fixture","type":"text","text":{"body":"hello"}}]}}]}]}`

func TestReviewPositiveVerifiedTextRoute(t *testing.T) {
	for name, tc := range map[string]struct {
		channel, recipient string
		valid              bool
	}{"matching": {"whatsapp", "p1", true}, "wrongChannel": {"page", "p1", false}, "wrongRecipient": {"whatsapp", "p2", false}, "invalidConfig": {"typo", "p1", false}} {
		t.Run(name, func(t *testing.T) {
			raw := []byte(reviewWAText)
			rows, e := ParseVerifiedWebhook("fixture-only-secret", raw, sign("fixture-only-secret", raw), tc.channel, tc.recipient)
			if tc.valid {
				if e != nil || len(rows) != 1 || rows[0].RecipientID != "p1" {
					t.Fatalf("valid row %v %v", rows, e)
				}
			} else if e == nil {
				t.Fatalf("wrong route returned messages %v", rows)
			}
		})
	}
	raw := []byte(reviewWAText)
	if _, e := ParseVerifiedWebhook("fixture-only-secret", raw, sign("different-fixture", raw), "whatsapp", "p1"); e == nil {
		t.Fatal("invalid signature accepted")
	}
}
func TestReviewObservedNoTextRouteBinding(t *testing.T) {
	for name, tc := range map[string]struct{ body, channel, recipient string }{"WAStatusWrongRecipient": {`{"object":"whatsapp_business_account","entry":[{"changes":[{"field":"messages","value":{"messaging_product":"whatsapp","metadata":{"phone_number_id":"other"},"statuses":[{"id":"fixture-id","status":"delivered"}]}}]}]}`, "whatsapp", "p1"}, "PageEchoWrongChannel": {`{"object":"page","entry":[{"messaging":[{"sender":{"id":"u1"},"recipient":{"id":"p1"},"message":{"text":"hello","mid":"fixture-id","is_echo":true}}]}]}`, "whatsapp", "p1"}, "EmptyInvalidConfig": {`{"object":"page","entry":[]}`, "typo", "p1"}} {
		t.Run(name, func(t *testing.T) {
			raw := []byte(tc.body)
			rows, e := ParseVerifiedWebhook("fixture-only-secret", raw, sign("fixture-only-secret", raw), tc.channel, tc.recipient)
			t.Logf("OBSERVED route=%s recipient=%s rows=%d error=%v", tc.channel, tc.recipient, len(rows), e)
			if e != nil || len(rows) != 0 {
				t.Fatal("observed baseline changed; re-evaluate claim")
			}
		})
	}
}
