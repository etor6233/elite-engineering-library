// Copyright (c) Microsoft Corporation.
// SPDX-License-Identifier: MIT
// ADAPTED from BCApps 2eae56d704a1fd035d104f333602aea7091b7749,
// CurrencyExchangeRate.Table.al FindCurrency2 lines459–470.
// Slice representation, duplicate rejection and explicit date are local deltas.
package bcfx

import (
	"math/big"
	"time"
)

type DatedRate struct {
	Currency     string
	StartingDate time.Time
	Amounts      DirectRate
}

// FindLast translates SetRange(currency), SetRange(0D,date), FindLast().
// An explicit date is required: this boundary never substitutes WorkDate/Now.
// Duplicate source identities fail instead of depending on input slice order.
func FindLast(rows []DatedRate, currency string, date time.Time) (DatedRate, error) {
	var result DatedRate
	if currency == "" || !dayOnly(date) || len(rows) > 4096 {
		return result, ErrRate
	}
	seen := make(map[int64]struct{})
	for _, row := range rows {
		if row.Currency != currency {
			continue
		}
		if !dayOnly(row.StartingDate) || !validRate(row.Amounts) {
			return DatedRate{}, ErrRate
		}
		key := row.StartingDate.Unix()
		if _, ok := seen[key]; ok {
			return DatedRate{}, ErrRate
		}
		seen[key] = struct{}{}
		if !row.StartingDate.After(date) && (result.StartingDate.IsZero() || row.StartingDate.After(result.StartingDate)) {
			result = row
		}
	}
	if result.StartingDate.IsZero() {
		return DatedRate{}, ErrRate
	}
	result.Amounts = DirectRate{new(big.Rat).Set(result.Amounts.Exchange), new(big.Rat).Set(result.Amounts.Relational)}
	return result, nil
}

func dayOnly(t time.Time) bool {
	return !t.IsZero() && t.Equal(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC))
}
