package bcsales

// AUTHORED source-conservation fixtures, not execution of Microsoft AL suites.
import (
	"errors"
	"reflect"
	"testing"
)

func TestQuoteTransferPreservesFieldsAndFiltersForeignDocuments(t *testing.T) {
	source := Header{DocumentType: Quote, ID: "quote-a", CustomerID: "customer-a", Currency: "ARS"}
	lines := []Line{
		{DocumentType: Quote, DocumentID: "quote-b", VariantID: "foreign", Quantity: 1, UnitPriceMinor: 9},
		{DocumentType: Order, DocumentID: "quote-a", VariantID: "wrong-type", Quantity: 1, UnitPriceMinor: 9},
		{DocumentType: Quote, DocumentID: "quote-a", VariantID: "one", Quantity: 1, UnitPriceMinor: 1001},
		{DocumentType: Quote, DocumentID: "quote-a", VariantID: "two", Quantity: 3, UnitPriceMinor: 2003},
	}
	before := append([]Line(nil), lines...)
	order, transferred, err := TransferQuoteToOrder(source, lines, "order-a")
	if err != nil || order != (Header{DocumentType: Order, ID: "order-a", QuoteID: "quote-a", CustomerID: "customer-a", Currency: "ARS"}) {
		t.Fatalf("header copy/retarget failed: %+v %v", order, err)
	}
	want := []Line{{DocumentType: Order, DocumentID: "order-a", VariantID: "one", Quantity: 1, UnitPriceMinor: 1001}, {DocumentType: Order, DocumentID: "order-a", VariantID: "two", Quantity: 3, UnitPriceMinor: 2003}}
	if !reflect.DeepEqual(transferred, want) || !reflect.DeepEqual(lines, before) {
		t.Fatalf("line conservation/scope failed: %+v", transferred)
	}
	transferred[0].UnitPriceMinor = 999
	if !reflect.DeepEqual(lines, before) {
		t.Fatal("output aliases the caller's immutable quote lines")
	}
	if _, _, err = TransferQuoteToOrder(Header{DocumentType: Order}, lines, "new"); !errors.Is(err, ErrQuoteType) {
		t.Fatal("non-quote source accepted")
	}
	_, empty, err := TransferQuoteToOrder(source, lines[:2], "order-empty")
	if err != nil || len(empty) != 0 {
		t.Fatal("foreign lines copied or a new empty-order policy invented")
	}
}

func FuzzQuoteTransferConservation(f *testing.F) {
	f.Add("quote-a", "order-a", "variant-a", int64(1), int64(1000), true)
	f.Add("q", "o", "v", int64(0), int64(-1), false)
	f.Add("q", "o", "v", int64(9223372036854775807), int64(9223372036854775807), true)
	f.Fuzz(func(t *testing.T, quoteID, orderID, variant string, quantity, price int64, belongs bool) {
		if len(quoteID)+len(orderID)+len(variant) > 1024 {
			t.Skip()
		}
		lineID := quoteID
		if !belongs {
			lineID += "different"
		}
		source := Header{DocumentType: Quote, ID: quoteID, Currency: "ARS", CustomerID: "customer"}
		line := Line{DocumentType: Quote, DocumentID: lineID, VariantID: variant, Quantity: quantity, UnitPriceMinor: price}
		order, lines, err := TransferQuoteToOrder(source, []Line{line}, orderID)
		if err != nil || order.ID != orderID || order.QuoteID != quoteID || order.Currency != source.Currency || order.CustomerID != source.CustomerID || order.DocumentType != Order {
			t.Fatal("header conservation failed")
		}
		if !belongs {
			if len(lines) != 0 {
				t.Fatal("foreign quote copied")
			}
			return
		}
		if len(lines) != 1 || lines[0].DocumentID != orderID || lines[0].DocumentType != Order || lines[0].Quantity != quantity || lines[0].UnitPriceMinor != price || lines[0].VariantID != variant {
			t.Fatal("line field changed beyond type/document")
		}
	})
}
