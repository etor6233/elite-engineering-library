package commerce

import (
	"context"
	"testing"
	"time"
)

type fakeRepo struct{ calls int }

func (f *fakeRepo) CreatePriceBook(context.Context, string, string, PriceBook) error {
	f.calls++
	return nil
}
func (f *fakeRepo) ActivatePriceBook(context.Context, string, string, string) error { return nil }
func (f *fakeRepo) PublicPrice(context.Context, string, string, string) (PriceEntry, error) {
	return PriceEntry{}, nil
}
func (f *fakeRepo) AddOrderLine(_ context.Context, _ string, _ string, v OrderLine, version int64) (OrderLine, error) {
	v.OrderVersion = version + 1
	return v, nil
}
func (f *fakeRepo) PlaceOrder(context.Context, string, string, string, int64, string) error {
	return nil
}
func (f *fakeRepo) AllocateStock(context.Context, string, string, string, string, string, int64, int64, string) error {
	return nil
}
func (f *fakeRepo) CreatePaymentAttempt(context.Context, string, string, string, PaymentAttempt) error {
	return nil
}
func (f *fakeRepo) TransitionPayment(context.Context, string, string, string, string, string, int64, string, string) error {
	return nil
}

type ids struct{ n int }

func (i *ids) New() string {
	i.n++
	return []string{"018f4d4a-7b36-7a21-8d10-2f4c54c28301", "018f4d4a-7b36-7a21-8d10-2f4c54c28302", "018f4d4a-7b36-7a21-8d10-2f4c54c28303", "018f4d4a-7b36-7a21-8d10-2f4c54c28304"}[i.n-1]
}
func TestCommerceRejectsInvalidPolicies(t *testing.T) {
	service := NewService(&fakeRepo{}, &ids{})
	ctx := context.Background()
	_, err := service.CreatePriceBook(ctx, "tenant", PriceBook{Market: "AR", Currency: "ARS", ValidFrom: time.Now(), Entries: []PriceEntry{{VariantID: "v", AmountMinorUnits: 1, TaxMode: "inclusive"}, {VariantID: "v", AmountMinorUnits: 2, TaxMode: "inclusive"}}})
	if err == nil {
		t.Fatal("duplicate variant accepted")
	}
	if err := service.TransitionPayment(ctx, "tenant", "o", "p", "created", "captured", 1, ""); err == nil {
		t.Fatal("payment skipped authorization")
	}
	if err := service.AllocateStock(ctx, "tenant", "org", "o", "l", "s", 0, 1); err == nil {
		t.Fatal("allocation without order version")
	}
}

func TestPaymentRecorderAndActorFailClosed(t *testing.T) {
	s := NewService(&fakeRepo{}, &ids{})
	in := PaymentAttempt{OrderID: "order", OrganizationID: "org", ProviderCode: "synthetic", Currency: "ARS", AmountMinorUnits: 100}
	if _, err := s.CreatePaymentAttempt(context.Background(), "tenant", "valid-payment-key", in); err == nil {
		t.Fatal("legacy-only repository silently used")
	}
	if _, err := s.CreatePaymentAttemptAs(context.Background(), "tenant", "valid-payment-key", in, ""); err == nil {
		t.Fatal("missing authenticated actor accepted")
	}
}
