package postgres

// SPDX-License-Identifier: MIT
// ADAPTED from Microsoft BCApps PriceCalculationBufferMgt.SetFiltersOnPriceListLine
// at 2eae56d704a1fd035d104f333602aea7091b7749, lines 258-269.
// Mapping: the local price table fixes sales/unit-price and exact variant/UOM;
// tenant, market, book ID and exact currency remain caller bindings. SQL NULL
// ending maps to BC's unbounded 0D. PostgreSQL timestamps have microsecond
// precision: local exclusive ending -> inclusive ending minus one microsecond.
// No fallback currency, UOM conversion, tax, discount or best-price policy.
const bcPriceEligibilitySQL = `b.status='active'
 and b.valid_from<=price_context.at
 and (b.valid_until is null or b.valid_until-interval '1 microsecond'>=price_context.at)`

// VerifySelectedLine:280-285 accepts the requested currency or BC's blank
// fallback. The existing local currency constraint excludes blank values, so
// that branch is unreachable. $3 is the caller's locked order currency.
const bcPriceCurrencyBindingSQL = `b.currency=$3`
