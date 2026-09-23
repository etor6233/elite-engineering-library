// AUTHORED typed composition input. Operations remains the purchase owner.
package serialsupply

import "regexp"

type CreateRequest struct {
	PlanRequest
	DestinationOrganizationID string `json:"destination_organization_id"`
	SupplierID                string `json:"supplier_id"`
	Currency                  string `json:"currency"`
	TotalMinorUnits           int64  `json:"total_minor_units,string"`
}

func (r CreateRequest) Valid() bool {
	return r.PlanRequest.Valid() && ValidID(r.DestinationOrganizationID) && ValidID(r.SupplierID) && regexp.MustCompile(`^[A-Z]{3}$`).MatchString(r.Currency) && r.TotalMinorUnits >= 0
}
