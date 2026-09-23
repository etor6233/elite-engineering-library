// Copyright (c) Microsoft Corporation.
// SPDX-License-Identifier: MIT
// Adapted from microsoft/BCApps at 2eae56d704a1fd035d104f333602aea7091b7749.
// This Go translation and its representation guards are local modifications.
// See docs/provenance/BC_AMOUNT_DERIVATION.md and licenses/Microsoft-BCApps-MIT.txt.
package bcamounts

import (
	"errors"
	"math/big"
)

var ErrAmount = errors.New("amount is outside the supported exact representation")
var ErrBalance = errors.New("journal is out of balance")

// LineAmount translates SalesLine.Table.al UpdateAmounts, line 5883:
// Round(Quantity * UnitPrice, AmountRoundingPrecision) - LineDiscountAmount.
// This profile uses integral quantities, minor units and rounding precision 1;
// the product is already exact, so rounding is an identity. Tax, FX, fractional
// quantities and percentage discounts are outside this function's contract.
// big.Int preserves AL's exact arithmetic before the explicit int64 boundary.
func LineAmount(quantity, unitPriceMinor, lineDiscountMinor int64) (int64, error) {
	var amount, right big.Int
	amount.SetInt64(quantity)
	right.SetInt64(unitPriceMinor)
	amount.Mul(&amount, &right)
	right.SetInt64(lineDiscountMinor)
	amount.Sub(&amount, &right)
	if !amount.IsInt64() {
		return 0, ErrAmount
	}
	return amount.Int64(), nil
}

// Entry is the existing library representation, not a Business Central type.
// Exactly one nonnegative side is positive; source AL signed amounts are mapped
// to Debit-Credit. The representation validation is a declared local delta.
type Entry struct{ Debit, Credit int64 }

// JournalTotals translates GenJnlPostBatch.ProcessBalanceOfLines' accumulation
// and CheckBalance into the library's single currency, single journal boundary.
// Document/currency grouping is supplied and authorized by the caller.
func JournalTotals(entries []Entry) (int64, int64, error) {
	var debit, credit, value big.Int
	for _, entry := range entries {
		if entry.Debit < 0 || entry.Credit < 0 || (entry.Debit > 0) == (entry.Credit > 0) {
			return 0, 0, ErrAmount
		}
		value.SetInt64(entry.Debit)
		debit.Add(&debit, &value)
		value.SetInt64(entry.Credit)
		credit.Add(&credit, &value)
	}
	if !debit.IsInt64() || !credit.IsInt64() {
		return 0, 0, ErrAmount
	}
	d, c := debit.Int64(), credit.Int64()
	if err := CheckBalance(d, c); err != nil {
		return 0, 0, err
	}
	return d, c, nil
}

// CheckBalance is the restricted equivalent of GenJnlPostBatch.CheckBalance:
// nonzero signed balance fails. Positive turnover is an existing library rule;
// it is explicit here and is not attributed to Microsoft.
func CheckBalance(debit, credit int64) error {
	if debit <= 0 || credit <= 0 {
		return ErrAmount
	}
	if debit != credit {
		return ErrBalance
	}
	return nil
}
