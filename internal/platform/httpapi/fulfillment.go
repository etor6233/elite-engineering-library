package httpapi

import (
	"elite.local/enterprise/internal/fulfillment"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"net/http"
)

type FulfillmentModule struct{ Service *fulfillment.Service }

func (m FulfillmentModule) Register(mux *http.ServeMux, v identity.Verifier) {
	a := fulfillmentAPI{service: m.Service, verifier: v}
	mux.HandleFunc("POST /v1/logistics/shipments", a.createShipment)
	mux.HandleFunc("POST /v1/logistics/shipments/{id}/transitions", a.transitionShipment)
	mux.HandleFunc("POST /v1/logistics/customer-transports", a.createCustomerTransport)
	mux.HandleFunc("POST /v1/logistics/shipments/{id}/provider-reports", a.recordProviderReport)
	mux.HandleFunc("POST /v1/service/cases", a.openCase)
	mux.HandleFunc("POST /v1/service/cases/{id}/transitions", a.transitionCase)
	mux.HandleFunc("POST /v1/service/recalls", a.createRecall)
	mux.HandleFunc("POST /v1/service/recalls/{id}/activate", a.activateRecall)
	mux.HandleFunc("POST /v1/service/recalls/{id}/units", a.addRecallUnit)
	mux.HandleFunc("POST /v1/franchise/network/organizations", a.createOrganization)
	mux.HandleFunc("POST /v1/franchise/network/organizations/{id}/transitions", a.transitionOrganization)
	mux.HandleFunc("POST /v1/franchise/agreements", a.createAgreement)
	mux.HandleFunc("POST /v1/franchise/agreements/{id}/transitions", a.transitionAgreement)
	mux.HandleFunc("POST /v1/communications/messages", a.queueMessage)
	mux.HandleFunc("POST /v1/communications/messages/{id}/transitions", a.transitionMessage)
}
func (a fulfillmentAPI) createOrganization(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "network:admin")
	if !ok {
		return
	}
	var v fulfillment.Organization
	if !decodeStrict(w, r, &v) {
		return
	}
	if (v.ParentOrganizationID == "" && !p.Allowed("*")) || (v.ParentOrganizationID != "" && !p.AllowedOrganization(v.ParentOrganizationID)) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for the parent organization")
		return
	}
	result, err := a.service.CreateOrganization(r.Context(), p.TenantID, v)
	if err != nil {
		writeFulfillment(w, err)
		return
	}
	writeJSON(w, 201, result)
}
func (a fulfillmentAPI) transitionOrganization(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "network:admin")
	if !ok {
		return
	}
	var v struct {
		Current string `json:"current"`
		Target  string `json:"target"`
		Version int64  `json:"version"`
	}
	if !decodeStrict(w, r, &v) {
		return
	}
	if !p.AllowedOrganization(r.PathValue("id")) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeFulfillment(w, a.service.TransitionOrganization(r.Context(), p.TenantID, r.PathValue("id"), v.Current, v.Target, v.Version))
}

type fulfillmentAPI struct {
	service  *fulfillment.Service
	verifier identity.Verifier
}

func (a fulfillmentAPI) auth(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return p, false
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return p, false
	}
	return p, true
}
func (a fulfillmentAPI) createShipment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "logistics:write")
	if !ok {
		return
	}
	var v fulfillment.Shipment
	if !decodeStrict(w, r, &v) {
		return
	}
	if !p.AllowedOrganization(v.OriginOrganizationID) || !p.AllowedOrganization(v.DestinationOrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for every shipment organization")
		return
	}
	result, err := a.service.CreateShipment(r.Context(), p.TenantID, v)
	if err != nil {
		writeFulfillment(w, err)
		return
	}
	writeJSON(w, 201, result)
}
func (a fulfillmentAPI) transitionShipment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "logistics:write")
	if !ok {
		return
	}
	var v struct {
		OriginOrganizationID      string `json:"origin_organization_id"`
		DestinationOrganizationID string `json:"destination_organization_id"`
		Current                   string `json:"current"`
		Target                    string `json:"target"`
		ProviderReference         string `json:"provider_reference"`
	}
	if !decodeStrict(w, r, &v) {
		return
	}
	if !p.AllowedOrganization(v.OriginOrganizationID) || !p.AllowedOrganization(v.DestinationOrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for every shipment organization")
		return
	}
	writeFulfillment(w, a.service.TransitionShipment(r.Context(), p.TenantID, v.OriginOrganizationID, v.DestinationOrganizationID, r.PathValue("id"), v.Current, v.Target, v.ProviderReference))
}
func (a fulfillmentAPI) createCustomerTransport(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "logistics:write")
	if !ok {
		return
	}
	var value fulfillment.CustomerTransportCommand
	if !decodeStrict(w, r, &value) {
		return
	}
	if !p.AllowedOrganization(value.OriginOrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for the shipment origin")
		return
	}
	result, err := a.service.CreateCustomerTransport(r.Context(), p.TenantID, value)
	if err != nil {
		writeFulfillment(w, err)
		return
	}
	writeJSON(w, 201, result)
}
func (a fulfillmentAPI) recordProviderReport(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "logistics:write")
	if !ok {
		return
	}
	var value fulfillment.ProviderReportCommand
	if !decodeStrict(w, r, &value) {
		return
	}
	if !p.AllowedOrganization(value.OriginOrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for the shipment origin")
		return
	}
	result, err := a.service.RecordProviderReport(r.Context(), p.TenantID, r.PathValue("id"), value)
	if err != nil {
		writeFulfillment(w, err)
		return
	}
	writeJSON(w, 200, result)
}
func (a fulfillmentAPI) openCase(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "service:write")
	if !ok {
		return
	}
	var v fulfillment.ServiceCase
	if !decodeStrict(w, r, &v) {
		return
	}
	if !p.AllowedOrganization(v.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	result, err := a.service.OpenServiceCase(r.Context(), p.TenantID, v)
	if err != nil {
		writeFulfillment(w, err)
		return
	}
	writeJSON(w, 201, result)
}
func (a fulfillmentAPI) transitionCase(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "service:write")
	if !ok {
		return
	}
	var v struct {
		OrganizationID string `json:"organization_id"`
		Current        string `json:"current"`
		Target         string `json:"target"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &v) {
		return
	}
	if !p.AllowedOrganization(v.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeFulfillment(w, a.service.TransitionServiceCase(r.Context(), p.TenantID, v.OrganizationID, r.PathValue("id"), v.Current, v.Target, v.Version))
}
func (a fulfillmentAPI) createRecall(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "recall:write")
	if !ok {
		return
	}
	var v fulfillment.Recall
	if !decodeStrict(w, r, &v) {
		return
	}
	result, err := a.service.CreateRecall(r.Context(), p.TenantID, v)
	if err != nil {
		writeFulfillment(w, err)
		return
	}
	writeJSON(w, 201, result)
}
func (a fulfillmentAPI) activateRecall(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "recall:publish")
	if !ok {
		return
	}
	writeFulfillment(w, a.service.ActivateRecall(r.Context(), p.TenantID, r.PathValue("id")))
}
func (a fulfillmentAPI) addRecallUnit(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "recall:write")
	if !ok {
		return
	}
	var v struct {
		OrganizationID string `json:"organization_id"`
		StockUnitID    string `json:"stock_unit_id"`
	}
	if !decodeStrict(w, r, &v) {
		return
	}
	if !p.AllowedOrganization(v.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeFulfillment(w, a.service.AddRecallUnit(r.Context(), p.TenantID, v.OrganizationID, r.PathValue("id"), v.StockUnitID))
}
func (a fulfillmentAPI) createAgreement(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "franchise:write")
	if !ok {
		return
	}
	var v fulfillment.Agreement
	if !decodeStrict(w, r, &v) {
		return
	}
	if !p.AllowedOrganization(v.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	result, err := a.service.CreateAgreement(r.Context(), p.TenantID, v)
	if err != nil {
		writeFulfillment(w, err)
		return
	}
	writeJSON(w, 201, result)
}
func (a fulfillmentAPI) transitionAgreement(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "franchise:write")
	if !ok {
		return
	}
	var v struct {
		OrganizationID string `json:"organization_id"`
		Current        string `json:"current"`
		Target         string `json:"target"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &v) {
		return
	}
	if !p.AllowedOrganization(v.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeFulfillment(w, a.service.TransitionAgreement(r.Context(), p.TenantID, v.OrganizationID, r.PathValue("id"), v.Current, v.Target, v.Version))
}
func (a fulfillmentAPI) queueMessage(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "communication:write")
	if !ok {
		return
	}
	var v fulfillment.Message
	if !decodeStrict(w, r, &v) {
		return
	}
	result, err := a.service.QueueMessage(r.Context(), p.TenantID, v)
	if err != nil {
		writeFulfillment(w, err)
		return
	}
	writeJSON(w, 202, result)
}
func (a fulfillmentAPI) transitionMessage(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "communication:provider")
	if !ok {
		return
	}
	var v struct {
		Current           string `json:"current"`
		Target            string `json:"target"`
		ProviderReference string `json:"provider_reference"`
		ErrorCode         string `json:"error_code"`
	}
	if !decodeStrict(w, r, &v) {
		return
	}
	writeFulfillment(w, a.service.TransitionMessage(r.Context(), p.TenantID, r.PathValue("id"), v.Current, v.Target, v.ProviderReference, v.ErrorCode))
}
func writeFulfillment(w http.ResponseWriter, err error) {
	if errors.Is(err, fulfillment.ErrConflict) {
		writeProblem(w, 409, "FULFILLMENT_CONFLICT", "state, version or business invariant conflict")
		return
	}
	if err != nil {
		writeProblem(w, 400, "INVALID_FULFILLMENT_COMMAND", "command does not match contract")
		return
	}
	writeJSON(w, 200, map[string]string{"status": "accepted"})
}
