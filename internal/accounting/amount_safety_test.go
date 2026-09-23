package accounting

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"
)

// Local representation boundary regression. These are not Microsoft fixtures.
func TestJournalRejectsNegativeSidesAndWrappedTotalsBeforeRepository(t *testing.T) {
	cases := map[string][]Line{
		"negative opposite sides": {
			{LineNo: 1, AccountCode: "CASH", DebitMinorUnits: 10, CreditMinorUnits: -5},
			{LineNo: 2, AccountCode: "REVENUE", DebitMinorUnits: -5, CreditMinorUnits: 10},
		},
		"int64 totals wrap to positive": {
			{LineNo: 1, AccountCode: "CASH", DebitMinorUnits: math.MaxInt64},
			{LineNo: 2, AccountCode: "CASH", DebitMinorUnits: math.MaxInt64},
			{LineNo: 3, AccountCode: "CASH", DebitMinorUnits: 3},
			{LineNo: 4, AccountCode: "REVENUE", CreditMinorUnits: math.MaxInt64},
			{LineNo: 5, AccountCode: "REVENUE", CreditMinorUnits: math.MaxInt64},
			{LineNo: 6, AccountCode: "REVENUE", CreditMinorUnits: 3},
		},
	}
	for name, lines := range cases {
		t.Run(name, func(t *testing.T) {
			repo := &fakeRepo{}
			_, err := NewService(repo, &ids{}).CreateJournal(context.Background(), "tenant", Journal{
				OrganizationID: "org", PeriodID: "period", SourceType: "SALE", SourceID: "order",
				Currency: "ARS", PostingDate: time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC), Lines: lines,
			})
			if !errors.Is(err, ErrInvalid) || repo.journal.ID != "" {
				t.Fatalf("invalid journal reached repository: error=%v totalDebit=%d totalCredit=%d persisted=%t", err, repo.journal.TotalDebitMinorUnits, repo.journal.TotalCreditMinorUnits, repo.journal.ID != "")
			}
		})
	}
}
