package marketplacebridge_test

import (
	"bytes"
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/electromobility"
	mb "elite.local/enterprise/internal/marketplacebridge"
	fixture "elite.local/enterprise/internal/marketplacebridge/testfixture"
	"elite.local/enterprise/internal/outbounddelivery"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestMarketplaceJSONBWhitespaceAndCrossProductBinding(t *testing.T) {
	p := profile("seller_warehouse")
	receiver := fixture.New(p)
	defer receiver.Close()
	client, _ := mb.NewClient(p, tokens{}, receiver.HTTP())
	intent := prepare(t, p, client, "PRICE")
	var spaced bytes.Buffer
	if e := json.Indent(&spaced, intent.Body, "", "  "); e != nil {
		t.Fatal(e)
	}
	intent.Body = spaced.Bytes()
	if _, e := client.Write(context.Background(), intent); e != nil {
		t.Fatal("semantic JSONB representation changed wire authorization", e)
	}
	if receiver.Writes[0].Body != `{"price":90071992547409.91}` {
		t.Fatal("wire is not canonical", receiver.Writes[0].Body)
	}
	intent = prepare(t, p, client, "STOCK")
	original := intent.Expected.UserProductID
	intent.Expected.UserProductID = "MLAU987654321"
	intent.Path = strings.ReplaceAll(intent.Path, original, intent.Expected.UserProductID)
	var terminal *outbounddelivery.TerminalFailure
	if _, e := client.Write(context.Background(), intent); !errors.As(e, &terminal) || len(receiver.Writes) != 1 {
		t.Fatal("observed item does not own write UP", e)
	}
}

type tokens struct{}

func (tokens) MercadoLibreAccessToken(context.Context) (string, error) { return fixture.Token, nil }
func profile(mode string) mb.Profile {
	b := mb.Binding{VariantID: "fixture-variant", SKU: "fixture-bicycle", CategoryID: "MLA3530", ListingTypeID: "gold_special", StockMode: mode, Attributes: []mb.Attribute{{ID: "ITEM_CONDITION", ValueID: "2230284"}}}
	if mode == "seller_warehouse" {
		b.StoreID = "123456"
		b.NetworkNodeID = "ARP12345"
	}
	return mb.Profile{Schema: "elite-marketplace-publication/v1", TenantID: "ca6c9cb7-e548-4f19-a2e6-27d2aa2d4893", OrganizationID: "j3-store", SellerID: "123456789", SiteID: "MLA", Currency: "ARS", CurrencyDigits: 2, MediaOrigin: "https://catalog.example.invalid", Bindings: []mb.Binding{b}}
}
func pub() cr.PublicDocument {
	return cr.PublicDocument{Generation: 1, SourceSHA256: strings.Repeat("a", 64), Currency: "ARS", Models: []cr.Model{{Model: electromobility.Model{ID: "model", DisplayName: "Fixture bicycle"}, Media: cr.Media{SHA256: strings.Repeat("b", 64)}}}, Variants: []cr.Variant{{ID: "fixture-variant", ModelID: "model", AmountMinorUnits: 9007199254740991}}}
}
func prepare(t *testing.T, p mb.Profile, c *mb.Client, op string) mb.Intent {
	t.Helper()
	before, e := c.Observe(context.Background(), p.Bindings[0].VariantID, "MLA123456789")
	if e != nil {
		t.Fatal(e)
	}
	in := mb.PrepareRequest{ApprovalID: "marketplace-unit-0001", Generation: 1, VariantID: p.Bindings[0].VariantID, Operation: op, ItemID: "MLA123456789", ExpiresAt: time.Now().Add(time.Minute)}
	out, e := mb.Build(p, pub(), in, 2, time.Now(), before)
	if e != nil {
		t.Fatal(e)
	}
	return out
}
func TestMarketplaceMutationContracts(t *testing.T) {
	for _, mode := range []string{"item", "selling_address", "seller_warehouse"} {
		t.Run(mode, func(t *testing.T) {
			p := profile(mode)
			receiver := fixture.New(p)
			defer receiver.Close()
			client, e := mb.NewClient(p, tokens{}, receiver.HTTP())
			if e != nil {
				t.Fatal(e)
			}
			for _, op := range []string{"PRICE", "STOCK", "PAUSE", "RESUME"} {
				intent := prepare(t, p, client, op)
				receipt, e := client.Write(context.Background(), intent)
				if e != nil || receipt.Validate() != nil {
					t.Fatal(op, e)
				}
			}
			if len(receiver.Writes) != 4 || receiver.Effects != 4 || receiver.PriceMinor != 9007199254740991 || receiver.Quantity != 2 || receiver.Item.Status != "active" {
				t.Fatal("request/effect mismatch", receiver.Writes)
			}
			if !strings.Contains(receiver.Writes[0].Body, "90071992547409.91") {
				t.Fatal("minor units lost precision")
			}
			if mode != "item" && receiver.Writes[1].Version != "7" {
				t.Fatal("stock version not preserved")
			}
			t.Log("one exact field per write; standard price exact decimal; configured stock location; pause/resume read back")
		})
	}
}
func TestMarketplaceFalseSuccessAndUnknownAreReadOnly(t *testing.T) {
	p := profile("seller_warehouse")
	receiver := fixture.New(p)
	defer receiver.Close()
	client, _ := mb.NewClient(p, tokens{}, receiver.HTTP())
	intent := prepare(t, p, client, "PRICE")
	receiver.IgnorePrice = true
	if _, e := client.Write(context.Background(), intent); !errors.Is(e, mb.ErrUnknown) {
		t.Fatal("200 with ignored price accepted", e)
	}
	if _, e := client.Reconcile(context.Background(), intent); !errors.Is(e, mb.ErrUnknown) {
		t.Fatal("wrong price reconciled", e)
	}
	receiver.IgnorePrice = false
	receiver.PriceMinor = intent.PriceMinorUnits
	if _, e := client.Reconcile(context.Background(), intent); e != nil {
		t.Fatal(e)
	}
	if len(receiver.Writes) != 1 {
		t.Fatal("recovery wrote")
	}
	intent = prepare(t, p, client, "STOCK")
	receiver.DropNext = true
	if _, e := client.Write(context.Background(), intent); !errors.Is(e, mb.ErrUnknown) {
		t.Fatal("lost response accepted", e)
	}
	if _, e := client.Reconcile(context.Background(), intent); e != nil {
		t.Fatal(e)
	}
	if len(receiver.Writes) != 2 || receiver.Quantity != 2 {
		t.Fatal("lost acknowledgement resent")
	}
}
func TestMarketplaceDriftApprovalBindingAndTerminal(t *testing.T) {
	p := profile("item")
	receiver := fixture.New(p)
	defer receiver.Close()
	client, _ := mb.NewClient(p, tokens{}, receiver.HTTP())
	intent := prepare(t, p, client, "PRICE")
	receiver.Item.Tags = []string{"dynamic_standard_price"}
	var terminal *outbounddelivery.TerminalFailure
	if _, e := client.Write(context.Background(), intent); !errors.As(e, &terminal) || len(receiver.Writes) != 0 {
		t.Fatal("automation guard", e)
	}
	receiver.Item.Tags = nil
	intent = prepare(t, p, client, "PRICE")
	tampered := intent
	tampered.Path = "/users/123456789"
	if _, e := client.Write(context.Background(), tampered); !errors.Is(e, mb.ErrBinding) || len(receiver.Writes) != 0 {
		t.Fatal("path injection")
	}
	tampered = intent
	tampered.Body = json.RawMessage(`{"price":1}`)
	if _, e := client.Write(context.Background(), tampered); !errors.Is(e, mb.ErrBinding) {
		t.Fatal("body injection")
	}
	receiver.RejectNext = true
	if _, e := client.Write(context.Background(), intent); !errors.As(e, &terminal) {
		t.Fatal("documented rejection", e)
	}
	receiver.Item.SellerID = 987654321
	if _, e := client.Observe(context.Background(), p.Bindings[0].VariantID, "MLA123456789"); !errors.Is(e, mb.ErrBinding) {
		t.Fatal("foreign seller")
	}
}
