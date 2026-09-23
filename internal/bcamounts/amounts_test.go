// SPDX-License-Identifier: LicenseRef-Workspace-Owner
// Local test harness for the declared BCApps translation; not a Microsoft suite.
package bcamounts

import (
	"errors"
	"math"
	"testing"
)

func TestSalesLineUpdateAmountsExactIntegerProfile(t *testing.T) {
	// Arithmetic oracle is the fixed SalesLine.UpdateAmounts expression. Inputs
	// are deterministic local fixtures, not purported Microsoft fixture values.
	cases := []struct {
		name                 string
		q, p, discount, want int64
	}{
		{"single", 1, 10001, 0, 10001},
		{"quantity", 3, 1200, 0, 3600},
		{"fixed discount", 3, 1200, 100, 3500},
		{"zero", 0, 1200, 0, 0},
		{"credit quantity", -2, 1200, 0, -2400},
		{"maximum", 1, math.MaxInt64, 0, math.MaxInt64},
		{"minimum", 1, math.MinInt64, 0, math.MinInt64},
		{"intermediate above range final exact", 2, math.MaxInt64, math.MaxInt64, math.MaxInt64},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := LineAmount(tc.q, tc.p, tc.discount)
			if err != nil || got != tc.want {
				t.Fatalf("got=%d err=%v want=%d", got, err, tc.want)
			}
		})
	}
	for _, args := range [][3]int64{{2, math.MaxInt64, 0}, {-1, math.MinInt64, 0}, {math.MaxInt64, math.MaxInt64, 0}, {1, math.MinInt64, 1}} {
		if _, err := LineAmount(args[0], args[1], args[2]); !errors.Is(err, ErrAmount) {
			t.Fatalf("accepted out-of-range %+v: %v", args, err)
		}
	}
}

func TestPostingBalanceMatchesSignedAmountConservation(t *testing.T) {
	for amount := int64(1); amount <= 200; amount++ {
		// Adapted invariant from GenJnlPostBatch.CheckBalance and the reverse
		// tests' out-of-balance oracle. Deterministic amounts replace random AL
		// setup and the BC runtime; no claim that the Microsoft suite ran.
		d, c, err := JournalTotals([]Entry{{Debit: amount}, {Credit: amount}})
		if err != nil || d != amount || c != amount {
			t.Fatalf("amount=%d totals=%d/%d err=%v", amount, d, c, err)
		}
		if err := CheckBalance(amount, amount+1); !errors.Is(err, ErrBalance) {
			t.Fatalf("unbalanced amount %d accepted", amount)
		}
	}
	if _, _, err := JournalTotals([]Entry{{Debit: math.MaxInt64}, {Credit: math.MaxInt64}}); err != nil {
		t.Fatal(err)
	}
	for _, entries := range [][]Entry{
		nil, {{Debit: 1, Credit: 1}}, {{Debit: -1}, {Credit: -1}},
		{{Debit: 10, Credit: -5}, {Debit: -5, Credit: 10}},
		{{Debit: math.MaxInt64}, {Debit: math.MaxInt64}, {Debit: 3}, {Credit: math.MaxInt64}, {Credit: math.MaxInt64}, {Credit: 3}},
	} {
		if _, _, err := JournalTotals(entries); err == nil {
			t.Fatalf("invalid representation accepted: %+v", entries)
		}
	}
}

func FuzzIntegralSalesLineConservation(f *testing.F) {
	for _, seed := range [][2]uint64{{0, 0}, {1, 10001}, {3, 1200}, {1, math.MaxInt64}, {2, math.MaxInt64}, {math.MaxInt64, math.MaxInt64}, {math.MaxInt64, 1}} {
		f.Add(seed[0], seed[1])
	}
	f.Fuzz(func(t *testing.T, rawQuantity, rawPrice uint64) {
		q, p := rawQuantity&math.MaxInt64, rawPrice&math.MaxInt64
		amount, err := LineAmount(int64(q), int64(p), 0)
		// Independent overflow oracle uses division; it does not repeat the
		// big-integer multiplication used by the translated implementation.
		if p != 0 && q > math.MaxInt64/p {
			if !errors.Is(err, ErrAmount) {
				t.Fatalf("overflow accepted: q=%d p=%d result=%d err=%v", q, p, amount, err)
			}
			return
		}
		if err != nil || amount != int64(q*p) {
			t.Fatalf("amount not conserved: q=%d p=%d result=%d err=%v", q, p, amount, err)
		}
	})
}
