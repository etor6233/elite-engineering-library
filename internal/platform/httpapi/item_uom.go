package httpapi

import (
	"net/http"

	"elite.local/enterprise/internal/inventorycontrol"
)

func (a bulkInventoryAPI) configureUnitOfMeasure(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:write", true)
	if !ok {
		return
	}
	var input inventorycontrol.ItemUnitOfMeasure
	if !decodeStrict(w, r, &input) {
		return
	}
	value, err := a.service.ConfigureUnitOfMeasure(r.Context(), p.TenantID, input)
	if writeBulkError(w, err, "INVALID_ITEM_UNIT_OF_MEASURE") {
		return
	}
	writeJSON(w, 201, value)
}

func (a bulkInventoryAPI) convertUnitOfMeasure(w http.ResponseWriter, r *http.Request) {
	p, ok := a.authorize(w, r, "inventory:read", true)
	if !ok {
		return
	}
	var input struct {
		ItemID   string `json:"item_id"`
		Code     string `json:"code"`
		Quantity string `json:"quantity"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	value, err := a.service.ConvertUnitOfMeasure(r.Context(), p.TenantID, input.ItemID, input.Code, input.Quantity)
	if writeBulkError(w, err, "INVALID_UNIT_OF_MEASURE_CONVERSION") {
		return
	}
	writeJSON(w, 200, value)
}
