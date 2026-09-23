package postgres

// AUTHORED query composition: one database-owned observation point for the
// source-derived eligibility predicate; no local/client clock or timezone rule.
const bcPriceTimeContextSQL = `cross join (select clock_timestamp() as at) as price_context`
