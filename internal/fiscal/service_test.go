package fiscal

import (
	"context"
	"errors"
	"testing"
	"time"
)

type serviceRepo struct{ requested int }

func (*serviceRepo) ConfigurePointOfSale(context.Context, string, string, PointOfSale) error {
	return nil
}
func (r *serviceRepo) RequestInvoice(_ context.Context, _, _, _, _ string, value Invoice) (Invoice, bool, error) {
	r.requested++
	value.TotalMinorUnits = value.NetMinorUnits + value.VATMinorUnits + value.ExemptMinorUnits + value.NonTaxedMinorUnits + value.OtherTaxMinorUnits
	return value, false, nil
}
func (*serviceRepo) GetInvoice(context.Context, string, string, string) (Invoice, error) {
	return Invoice{ID: "invoice"}, nil
}

type testIDs struct{ n int }

func (i *testIDs) New() string { i.n++; return "id" + string(rune('0'+i.n)) }

func TestValidCUITAndInvoiceBoundary(t *testing.T) {
	if !ValidCUIT("30715117564") || ValidCUIT("30715117565") || ValidCUIT("not-a-cuit") {
		t.Fatal("CUIT checksum contract failed")
	}
	repo := &serviceRepo{}
	service := NewService(repo, &testIDs{})
	base := Invoice{OrganizationID: "franchise", OrderID: "order", PaymentAttemptID: "payment", PointOfSaleID: "pos", VoucherType: 6, Concept: 1, RecipientDocumentType: 99, RecipientDocument: "0", RecipientVATConditionID: 5, NetMinorUnits: 8264, VATMinorUnits: 1736, VATLines: []VATLine{{ID: 5, BaseMinorUnits: 8264, AmountMinorUnits: 1736}}, IssuedOn: time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC)}
	if _, _, err := service.RequestInvoice(context.Background(), "tenant", "1234567890abcdef", string(make([]byte, 64)), base); !errors.Is(err, ErrInvalid) {
		t.Fatalf("non-hex hash accepted: %v", err)
	}
	value, replay, err := service.RequestInvoice(context.Background(), "tenant", "1234567890abcdef", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", base)
	if err != nil || replay || value.TotalMinorUnits != 10000 || repo.requested != 1 {
		t.Fatalf("value=%+v replay=%v err=%v count=%d", value, replay, err, repo.requested)
	}
	base.Concept = 2
	if _, _, err = service.RequestInvoice(context.Background(), "tenant", "1234567890abcdef", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", base); !errors.Is(err, ErrInvalid) {
		t.Fatalf("service dates omission accepted: %v", err)
	}
	base.Concept = 1
	base.VATLines[0].AmountMinorUnits = 1735
	if _, _, err = service.RequestInvoice(context.Background(), "tenant", "1234567890abcdef", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", base); !errors.Is(err, ErrInvalid) {
		t.Fatalf("VAT mismatch accepted: %v", err)
	}
	base.VATLines[0].AmountMinorUnits = 1736
	base.RecipientVATConditionID = 0
	if _, _, err = service.RequestInvoice(context.Background(), "tenant", "1234567890abcdef", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", base); !errors.Is(err, ErrInvalid) {
		t.Fatalf("missing recipient VAT condition accepted: %v", err)
	}
	base.RecipientVATConditionID = 5
	base.VoucherType = 8
	if _, _, err = service.RequestInvoice(context.Background(), "tenant", "1234567890abcdef", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", base); !errors.Is(err, ErrInvalid) {
		t.Fatalf("credit without associated invoice accepted: %v", err)
	}
	base.AssociatedVouchers = []AssociatedVoucher{{InvoiceID: "original-invoice"}}
	if _, _, err = service.RequestInvoice(context.Background(), "tenant", "1234567890abcdef", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", base); err != nil {
		t.Fatalf("valid credit association rejected: %v", err)
	}
	base.AssociatedVouchers[0].Number = 7
	if _, _, err = service.RequestInvoice(context.Background(), "tenant", "1234567890abcdef", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", base); !errors.Is(err, ErrInvalid) {
		t.Fatalf("client-supplied associated identity accepted: %v", err)
	}
}
