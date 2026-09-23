package order

import (
	"errors"
	"fmt"
)

type State string

const (
	Draft     State = "draft"
	Placed    State = "placed"
	Confirmed State = "confirmed"
	Paid      State = "paid"
	Allocated State = "allocated"
	Delivered State = "delivered"
	Cancelled State = "cancelled"
)

var ErrConflict = errors.New("order conflict")

type Order struct {
	ID                string `json:"id"`
	TenantID          string `json:"tenantId"`
	OrganizationID    string `json:"organizationId"`
	CustomerPrincipal string `json:"customerPrincipal"`
	State             State  `json:"state"`
	Currency          string `json:"currency"`
	TotalMinorUnits   int64  `json:"totalMinorUnits"`
	Version           int64  `json:"version"`
}

type transition struct {
	from       State
	to         State
	permission string
}

var transitions = []transition{
	{Draft, Placed, "order:create"},
	{Placed, Confirmed, "order:transition"},
	{Confirmed, Paid, "payment:reconcile"},
	{Paid, Allocated, "inventory:reserve"},
	{Allocated, Delivered, "order:transition"},
	{Draft, Cancelled, "order:transition"},
	{Placed, Cancelled, "order:transition"},
}

func New(id, tenantID, organizationID, customerPrincipal, currency string, totalMinorUnits int64) (Order, error) {
	if id == "" || tenantID == "" || organizationID == "" || customerPrincipal == "" {
		return Order{}, fmt.Errorf("%w: identifiers are required", ErrConflict)
	}
	if len(currency) != 3 || totalMinorUnits < 0 {
		return Order{}, fmt.Errorf("%w: invalid money", ErrConflict)
	}
	return Order{ID: id, TenantID: tenantID, OrganizationID: organizationID, CustomerPrincipal: customerPrincipal, State: Draft, Currency: currency, TotalMinorUnits: totalMinorUnits, Version: 1}, nil
}

func (o Order) Transition(target State, permissions map[string]struct{}) (Order, error) {
	for _, candidate := range transitions {
		if candidate.from != o.State || candidate.to != target {
			continue
		}
		if _, all := permissions["*"]; !all {
			if _, allowed := permissions[candidate.permission]; !allowed {
				return Order{}, fmt.Errorf("%w: permission %s required", ErrConflict, candidate.permission)
			}
		}
		o.State = target
		o.Version++
		return o, nil
	}
	return Order{}, fmt.Errorf("%w: transition %s -> %s is not configured", ErrConflict, o.State, target)
}
