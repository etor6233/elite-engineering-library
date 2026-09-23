package merchantbridge

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func FuzzMerchantBoundedCommand(f *testing.F) {
	f.Add([]byte(`{"approval_id":"fixture-approval-0001","generation":"1","variant_id":"bicycle","expires_at":"2026-09-13T00:00:00Z"}`))
	f.Add([]byte(`{"generation":"1","Generation":"2"}`))
	f.Add([]byte(`{"price":1e9999}`))
	f.Add([]byte(`{"x":null,"x":true}`))
	f.Fuzz(func(t *testing.T, raw []byte) {
		var r PrepareRequest
		e := Decode(raw, &r)
		if len(raw) > 32768 && e == nil {
			t.Fatal("unbounded command")
		}
		if e == nil {
			b, e := json.Marshal(r)
			if e != nil {
				t.Fatal(e)
			}
			var again PrepareRequest
			if Decode(b, &again) != nil || again != r {
				t.Fatal("canonical roundtrip")
			}
			if r.Validate() == nil && (r.Generation < 1 || len(r.ApprovalID) < 16 || r.ExpiresAt == (time.Time{})) {
				t.Fatal("invalid admitted intent")
			}
		}
	})
}
func TestMerchantStrictCommandBoundary(t *testing.T) {
	for _, raw := range []string{`{"generation":"1","Generation":"2"}`, `{"generation":"1","generation":"2"}`, `{"generation":1}`, `{"generation":"1"} {}`, strings.Repeat(" ", 32769)} {
		if Decode([]byte(raw), new(PrepareRequest)) == nil {
			t.Fatal("invalid admitted", raw[:min(60, len(raw))])
		}
	}
}
