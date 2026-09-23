package accounting

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepo struct{ journal Journal }

func (*fakeRepo) CreateAccount(context.Context, string, string, Account) error { return nil }
func (*fakeRepo) OpenPeriod(context.Context, string, string, Period) error     { return nil }
func (f *fakeRepo) CreateJournal(_ context.Context, _, _ string, v Journal) error {
	f.journal = v
	return nil
}
func (*fakeRepo) PostJournal(context.Context, string, string, string, string, int64, string, string) (Journal, error) {
	return Journal{Status: "posted"}, nil
}
func (*fakeRepo) ReverseJournal(context.Context, string, string, string, string, string, string, int64, string, string) (Journal, error) {
	return Journal{Status: "posted", ReversalOf: "j"}, nil
}
func (*fakeRepo) ClosePeriod(context.Context, string, string, int64, string) (Period, error) {
	return Period{Status: "closed"}, nil
}
func (*fakeRepo) TrialBalance(context.Context, string, string, string) ([]Balance, error) {
	return []Balance{{AccountCode: "CASH"}}, nil
}

type ids struct{ n int }

func (i *ids) New() string { i.n++; return "id" + string(rune('0'+i.n)) }
func TestBalancedJournal(t *testing.T) {
	r := &fakeRepo{}
	s := NewService(r, &ids{})
	j, err := s.CreateJournal(context.Background(), "tenant", Journal{OrganizationID: "org", PeriodID: "p", SourceType: "SALE", SourceID: "order", Currency: "ARS", PostingDate: time.Now(), Lines: []Line{{LineNo: 1, AccountCode: "CASH", DebitMinorUnits: 100}, {LineNo: 2, AccountCode: "REVENUE", CreditMinorUnits: 100}}})
	if err != nil || j.TotalDebitMinorUnits != 100 || r.journal.ID == "" {
		t.Fatalf("journal=%+v err=%v", j, err)
	}
}
func TestRejectsUnbalancedOrDualSided(t *testing.T) {
	s := NewService(&fakeRepo{}, &ids{})
	for _, lines := range [][]Line{{{LineNo: 1, AccountCode: "CASH", DebitMinorUnits: 100}, {LineNo: 2, AccountCode: "REVENUE", CreditMinorUnits: 99}}, {{LineNo: 1, AccountCode: "CASH", DebitMinorUnits: 100, CreditMinorUnits: 1}, {LineNo: 2, AccountCode: "REVENUE", CreditMinorUnits: 100}}} {
		_, err := s.CreateJournal(context.Background(), "tenant", Journal{OrganizationID: "org", PeriodID: "p", SourceType: "SALE", SourceID: "order", Currency: "ARS", PostingDate: time.Now(), Lines: lines})
		if !errors.Is(err, ErrInvalid) {
			t.Fatalf("accepted lines=%+v err=%v", lines, err)
		}
	}
}
