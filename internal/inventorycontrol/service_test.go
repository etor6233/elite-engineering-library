package inventorycontrol

import (
	"context"
	"testing"
	"time"
)

type fakeRepo struct{ reserve, release, create, transition int }

func (f *fakeRepo) AvailableToPromise(context.Context, string, string, string, time.Time) (ATP, error) {
	return ATP{AvailableToPromise: 2}, nil
}
func (f *fakeRepo) Reserve(_ context.Context, _, _, _ string, value Reservation, _ int64) (Reservation, error) {
	f.reserve++
	return value, nil
}
func (f *fakeRepo) Release(context.Context, string, string, string, int64, string) error {
	f.release++
	return nil
}
func (f *fakeRepo) CreateTransfer(_ context.Context, _, _ string, value Transfer, _ map[string]int64, _ string) (Transfer, error) {
	f.create++
	return value, nil
}
func (f *fakeRepo) TransitionTransfer(context.Context, string, string, string, string, int64, string, string) error {
	f.transition++
	return nil
}

type ids struct{ n int }

func (i *ids) New() string { i.n++; return "id" }

func TestServiceRejectsInvalidReservationAndTransfer(t *testing.T) {
	repo, generator := &fakeRepo{}, &ids{}
	service := NewService(repo, generator)
	if _, err := service.Reserve(context.Background(), "tenant", Reservation{}, 1); err == nil {
		t.Fatal("invalid reservation accepted")
	}
	if _, err := service.CreateTransfer(context.Background(), "tenant", Transfer{FromOrganizationID: "a", ToOrganizationID: "a", StockUnitIDs: []string{"s"}, ExpectedReceiptAt: time.Now()}, map[string]int64{"s": 1}); err == nil {
		t.Fatal("same-organization transfer accepted")
	}
	if err := service.TransitionTransfer(context.Background(), "tenant", "a", "t", "draft", "received", 1); err == nil {
		t.Fatal("skipped transfer transition accepted")
	}
	if repo.reserve+repo.create+repo.transition != 0 {
		t.Fatal("repository called for invalid input")
	}
}

func TestServiceBuildsGovernedCommands(t *testing.T) {
	repo, generator := &fakeRepo{}, &ids{}
	service := NewService(repo, generator)
	reservation, err := service.Reserve(context.Background(), "tenant", Reservation{OrganizationID: "a", StockUnitID: "s", VariantID: "v", DemandKind: "service", DemandID: "d"}, 2)
	if err != nil || reservation.Status != "reservation" || reservation.Version != 1 || repo.reserve != 1 {
		t.Fatalf("reservation not created: %#v %v", reservation, err)
	}
	transfer, err := service.CreateTransfer(context.Background(), "tenant", Transfer{FromOrganizationID: "a", ToOrganizationID: "b", StockUnitIDs: []string{"s"}, ExpectedReceiptAt: time.Now().Add(time.Hour)}, map[string]int64{"s": 3})
	if err != nil || transfer.State != "draft" || transfer.Version != 1 || repo.create != 1 {
		t.Fatalf("transfer not created: %#v %v", transfer, err)
	}
}
