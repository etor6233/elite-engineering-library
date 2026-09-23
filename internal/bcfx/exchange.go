// Copyright (c) Microsoft Corporation.
// SPDX-License-Identifier: MIT
// ADAPTED from BCApps 2eae56d704a1fd035d104f333602aea7091b7749,
// CurrencyExchangeRate.Table.al, ExchangeAmtFCYToFCY and ExchangeAmount.
// Go representation, bounds and errors are declared local modifications.
// See docs/provenance/BC_FX_DERIVATION.md and licenses/Microsoft-BCApps-MIT.txt.
package bcfx

import (
	"errors"
	"math/big"
)

var ErrRate = errors.New("unsupported or invalid exact exchange rate")
var ErrAmount = errors.New("converted amount is outside int64 minor units")

// DirectRate is the supported BC profile: Relational Currency Code is empty,
// so Exchange/Relational describe one currency against the selected local base.
// Currency identity and dated selection belong to the snapshot boundary.
type DirectRate struct{ Exchange, Relational *big.Rat }

func validRate(r DirectRate) bool {
	return r.Exchange != nil && r.Relational != nil && r.Exchange.Sign() > 0 && r.Relational.Sign() > 0 &&
		r.Exchange.Num().BitLen() <= 128 && r.Exchange.Denom().BitLen() <= 128 &&
		r.Relational.Num().BitLen() <= 128 && r.Relational.Denom().BitLen() <= 128
}

// ExchangeExact translates the direct-base branches at source lines383–387,
// 401–405 and425–433. Nil means the explicitly configured local currency.
// Callers supply major units; no floats, quote spread, fees or tax are added.
func ExchangeExact(amount *big.Rat, from, to *DirectRate) (*big.Rat, error) {
	if amount == nil || amount.Num().BitLen() > 128 || amount.Denom().BitLen() > 128 {
		return nil, ErrAmount
	}
	if from != nil && !validRate(*from) || to != nil && !validRate(*to) {
		return nil, ErrRate
	}
	out := new(big.Rat).Set(amount)
	if from != nil {
		out.Quo(out, from.Exchange)
		out.Mul(out, from.Relational)
	}
	if to != nil {
		out.Quo(out, to.Relational)
		out.Mul(out, to.Exchange)
	}
	return out, nil
}
