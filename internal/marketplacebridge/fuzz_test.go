package marketplacebridge

import (
	"bytes"
	"encoding/json"
	"testing"
)

func FuzzMarketplaceCommandBoundary(f *testing.F) {
	f.Add([]byte(`{"approval_id":"marketplace-content-01","media_approval_id":"marketplace-media-0001","generation":"1","variant_id":"bike-standard","operation":"CONTENT","item_id":"MLA123456789","expires_at":"2026-09-13T00:00:00Z"}`))
	f.Add([]byte(`{"approval_id":"marketplace-media-0001","generation":"1","variant_id":"bike-standard","operation":"MEDIA","expires_at":"2026-09-13T00:00:00Z"}`))
	f.Add([]byte(`{"approval_id":"marketplace-create-001","media_approval_id":"marketplace-media-0001","generation":"1","variant_id":"bike-standard","operation":"CREATE","expires_at":"2026-09-13T00:00:00Z"}`))
	f.Add([]byte(`{"approval_id":"marketplace-price-0001","generation":"1","variant_id":"bike-standard","operation":"PRICE","item_id":"MLA123456789","expires_at":"2026-09-13T00:00:00Z"}`))
	f.Add([]byte(`{"approval_id":"one","Approval_ID":"two","operation":"STOCK"}`))
	f.Add([]byte(`{"operation":"PRICE","item_id":"../../users"}`))
	f.Fuzz(func(t *testing.T, raw []byte) {
		var r PrepareRequest
		if Decode(raw, &r) != nil {
			return
		}
		before, _ := json.Marshal(r)
		err := r.Validate()
		after, _ := json.Marshal(r)
		if !bytes.Equal(before, after) {
			t.Fatal("validation mutated command")
		}
		var round PrepareRequest
		if Decode(after, &round) != nil {
			t.Fatal("accepted typed serialization cannot decode")
		}
		if err == nil {
			if len(r.ApprovalID) < 16 || r.Generation < 1 || r.ExpiresAt.IsZero() {
				t.Fatal("command boundary lost")
			}
			switch r.Operation {
			case "CONTENT":
				if !itemID.MatchString(r.ItemID) || len(r.MediaApprovalID) < 16 || r.MediaApprovalID == r.ApprovalID {
					t.Fatal("unbound content")
				}
			case "MEDIA":
				if r.ItemID != "" || r.MediaApprovalID != "" {
					t.Fatal("media target injection")
				}
			case "CREATE":
				if r.ItemID != "" || len(r.MediaApprovalID) < 16 || r.MediaApprovalID == r.ApprovalID {
					t.Fatal("unbound create")
				}
			case "PRICE", "STOCK", "PAUSE", "RESUME":
				if !itemID.MatchString(r.ItemID) || r.MediaApprovalID != "" {
					t.Fatal("mutation target injection")
				}
			default:
				t.Fatal("unadmitted write")
			}
		}
	})
}
