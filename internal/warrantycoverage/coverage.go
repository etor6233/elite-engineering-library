// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License; see Microsoft-BCApps-MIT.txt.
// ADAPTED from BCApps 2eae56d704a1fd035d104f333602aea7091b7749,
// ServiceItemLine.Table.al CheckWarranty, lines 1894-1913.
// Only the pure date decision is translated. AL validation/event hooks and
// service-line writes are not implemented or claimed here. Ordinal day 0
// represents AL 0D; callers must validate their own date/terms representation.
package warrantycoverage

import "errors"

var ErrPartsDateOrder = errors.New("warranty parts start follows end")

type Period struct{ Start, End int32 }
type Coverage struct{ Any, Parts, Labor bool }

// Evaluate preserves the source's inclusive comparisons and parts-only error.
// It selects no duration, percentage, fault exclusion, authorization or payment.
func Evaluate(date int32, parts, labor Period) (Coverage, error) {
	if parts.Start > parts.End {
		return Coverage{}, ErrPartsDateOrder
	}
	warrantyParts := date >= parts.Start && date <= parts.End
	warrantyLabor := date >= labor.Start && date <= labor.End
	return Coverage{Any: warrantyParts || warrantyLabor, Parts: warrantyParts, Labor: warrantyLabor}, nil
}
