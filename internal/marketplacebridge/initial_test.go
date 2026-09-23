package marketplacebridge_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	mb "elite.local/enterprise/internal/marketplacebridge"
	fixture "elite.local/enterprise/internal/marketplacebridge/testfixture"
	"elite.local/enterprise/internal/outbounddelivery"
	"encoding/hex"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"testing"
	"time"
)

func TestMarketplaceInitialContracts(t *testing.T) {
	for _, scenario := range []string{"accepted", "lost-create", "lost-upload", "wrong-picture", "duplicate-search", "rejected", "persist-failed"} {
		t.Run(scenario, func(t *testing.T) {
			p := profile("item")
			server := fixture.NewInitial(p)
			defer server.Close()
			client, _ := mb.NewClient(p, tokens{}, server.HTTP())
			document := pub()
			var data bytes.Buffer
			if e := png.Encode(&data, image.NewNRGBA(image.Rect(0, 0, 2, 2))); e != nil {
				t.Fatal(e)
			}
			sum := sha256.Sum256(data.Bytes())
			document.Models[0].Media.SHA256 = hex.EncodeToString(sum[:])
			request := mb.PrepareRequest{ApprovalID: "initial-media-unit", Generation: 1, VariantID: p.Bindings[0].VariantID, Operation: "MEDIA", ExpiresAt: time.Now().Add(time.Minute)}
			media, e := mb.Build(p, document, request, 2, time.Now(), mb.Observation{})
			if e != nil {
				t.Fatal(e)
			}
			var stored mb.InitialObservation
			record := func(_ context.Context, o mb.InitialObservation) error { stored = o; return nil }
			if scenario == "lost-upload" {
				server.DropUpload = true
			}
			receipt, e := client.Upload(context.Background(), media, data.Bytes(), record)
			if scenario == "lost-upload" {
				if !errors.Is(e, mb.ErrUnknown) || stored.ProviderID != "" || server.Uploads != 1 {
					t.Fatal("unknown upload", e)
				}
				return
			}
			if e != nil || receipt.ProviderMessageID != fixture.PictureID || server.UploadedSHA != media.MediaSHA256 {
				t.Fatal("upload binding", e)
			}
			if _, e = client.AcceptedUpload(media, stored); e != nil {
				t.Fatal("durable response recovery", e)
			}
			tampered := append([]byte(nil), data.Bytes()...)
			tampered[len(tampered)-1] ^= 1
			if _, e = client.Upload(context.Background(), media, tampered, record); e == nil || server.Uploads != 1 {
				t.Fatal("source hash bypass")
			}
			request = mb.PrepareRequest{ApprovalID: "initial-create-unit", MediaApprovalID: media.ApprovalID, Generation: 1, VariantID: p.Bindings[0].VariantID, Operation: "CREATE", ExpiresAt: time.Now().Add(time.Minute)}
			creation, e := mb.Build(p, document, request, 2, time.Now(), mb.Observation{Item: mb.Item{Pictures: []mb.Picture{{ID: fixture.PictureID}}}})
			if e != nil {
				t.Fatal(e)
			}
			var pretty bytes.Buffer
			json.Indent(&pretty, creation.Body, "", " ")
			creation.Body = pretty.Bytes()
			if scenario == "lost-create" {
				server.DropCreate = true
			}
			if scenario == "wrong-picture" {
				server.BadPicture = true
			}
			if scenario == "duplicate-search" {
				server.Exists = true
				server.DuplicateSearch = true
			}
			if scenario == "rejected" {
				server.RejectNext = true
			}
			if scenario == "persist-failed" {
				record = func(context.Context, mb.InitialObservation) error { return errors.New("fixture store unavailable") }
			}
			receipt, e = client.Create(context.Background(), creation, record)
			switch scenario {
			case "accepted":
				if e != nil || receipt.ProviderMessageID != server.Item.ID {
					t.Fatal("create", e)
				}
			case "lost-create", "persist-failed":
				if !errors.Is(e, mb.ErrUnknown) {
					t.Fatal("must be unknown", e)
				}
				if _, e = client.ReconcileCreation(context.Background(), creation, ""); e != nil {
					t.Fatal("GET-only recovery", e)
				}
			case "wrong-picture":
				if !errors.Is(e, mb.ErrUnknown) {
					t.Fatal("false success", e)
				}
				if _, e = client.ReconcileCreation(context.Background(), creation, server.Item.ID); e == nil {
					t.Fatal("wrong provider image accepted")
				}
			case "duplicate-search":
				if e == nil || server.Creates != 0 {
					t.Fatal("duplicate precondition")
				}
			case "rejected":
				var terminal *outbounddelivery.TerminalFailure
				if !errors.As(e, &terminal) {
					t.Fatal("documented rejection", e)
				}
			}
			if scenario != "duplicate-search" && server.Creates != 1 {
				t.Fatal("provider POST count", server.Creates)
			}
		})
	}
}
