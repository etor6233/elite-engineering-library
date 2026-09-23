package marketplacebridge_test

import (
	"context"
	mb "elite.local/enterprise/internal/marketplacebridge"
	fixture "elite.local/enterprise/internal/marketplacebridge/testfixture"
	"elite.local/enterprise/internal/outbounddelivery"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestMarketplaceContentContracts(t *testing.T) {
	for _, scenario := range []string{"accepted", "lost", "multi", "sold", "missing-sales", "wrong-picture", "rejected", "tampered"} {
		t.Run(scenario, func(t *testing.T) {
			p := profile("item")
			server := fixture.NewContent(p)
			defer server.Close()
			client, _ := mb.NewClient(p, tokens{}, server.HTTP())
			before, e := client.Observe(context.Background(), p.Bindings[0].VariantID, server.Item.ID)
			if e != nil {
				t.Fatal(e)
			}
			request := mb.PrepareRequest{ApprovalID: "content-unit-approval", MediaApprovalID: "content-unit-media-id", Generation: 1, VariantID: p.Bindings[0].VariantID, Operation: "CONTENT", ItemID: server.Item.ID, ExpiresAt: time.Now().Add(time.Minute)}
			intent, e := mb.Build(p, pub(), request, 2, time.Now(), before, fixture.PictureID)
			if e != nil {
				t.Fatal(e)
			}
			if scenario == "multi" {
				server.Multiple = true
			}
			if scenario == "sold" {
				server.Sales = 1
			}
			if scenario == "missing-sales" {
				server.MissingSales = true
			}
			if scenario == "lost" {
				server.DropNext = true
			}
			if scenario == "rejected" {
				server.RejectNext = true
			}
			if scenario == "tampered" {
				var body map[string]any
				json.Unmarshal(intent.Body, &body)
				body["price"] = 1
				intent.Body, _ = json.Marshal(body)
			}
			_, e = client.Write(context.Background(), intent)
			switch scenario {
			case "accepted", "wrong-picture":
				if e != nil {
					t.Fatal("content effect", e)
				}
				if scenario == "wrong-picture" {
					server.Item.Pictures[0].ID = "654321-MLA123456789_092026"
					if _, e = client.Reconcile(context.Background(), intent); e == nil {
						t.Fatal("wrong picture accepted")
					}
				}
			case "lost":
				if !errors.Is(e, mb.ErrUnknown) {
					t.Fatal("lost write", e)
				}
				if _, e = client.Reconcile(context.Background(), intent); e != nil {
					t.Fatal("read-only recovery", e)
				}
			case "multi", "sold", "missing-sales", "rejected":
				var terminal *outbounddelivery.TerminalFailure
				if !errors.As(e, &terminal) {
					t.Fatal("unsafe scope/rejection", e)
				}
			case "tampered":
				if e == nil {
					t.Fatal("extra monetary field")
				}
			}
			expected := 1
			if scenario == "multi" || scenario == "sold" || scenario == "missing-sales" || scenario == "tampered" {
				expected = 0
			}
			if server.ContentWrites != expected {
				t.Fatal("write count", server.ContentWrites)
			}
			if server.PriceMinor != 9007199254740993 || server.Quantity != 4 {
				t.Fatal("content changed commerce fields")
			}
		})
	}
}
