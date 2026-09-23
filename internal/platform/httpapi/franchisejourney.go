package httpapi

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/identity"
)

type FranchiseJourneyModule struct{ Service *franchisejourney.Service }

func (m FranchiseJourneyModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	a := franchiseJourneyAPI{service: m.Service, verifier: verifier}
	mux.HandleFunc("GET /v1/public/{tenantCode}/locations", a.publicLocations)
	mux.HandleFunc("GET /v1/public/{tenantCode}/{organizationCode}/appointment-slots", a.publicAppointmentSlots)
	mux.HandleFunc("POST /v1/public/{tenantCode}/{organizationCode}/appointments", a.requestAppointment)
	mux.HandleFunc("POST /v1/franchise/appointment-slots", a.createAppointmentSlot)
	mux.HandleFunc("GET /v1/franchise/appointment-slots/result", a.appointmentSlotCreationResult)
	mux.HandleFunc("POST /v1/franchise/resources", a.createServiceResource)
	mux.HandleFunc("GET /v1/franchise/resources/result", a.serviceResourceCreationResult)
	mux.HandleFunc("GET /v1/franchise/availability", a.availability)
	mux.HandleFunc("POST /v1/franchise/availability", a.createAvailability)
	mux.HandleFunc("GET /v1/franchise/availability/result", a.availabilityCreationResult)
	mux.HandleFunc("POST /v1/franchise/availability/{id}/cancel", a.cancelAvailability)
	mux.HandleFunc("POST /v1/franchise/appointments/{id}/resources", a.assignAppointmentResource)
	mux.HandleFunc("POST /v1/franchise/appointments/{id}/transitions", a.transitionAppointment)
	mux.HandleFunc("GET /v1/franchise/leads", a.leads)
	mux.HandleFunc("GET /v1/franchise/agenda", a.appointmentAgenda)
	mux.HandleFunc("POST /v1/franchise/leads/{id}/assign", a.assignLead)
	mux.HandleFunc("POST /v1/franchise/leads/{id}/transitions", a.transitionLead)
	mux.HandleFunc("POST /v1/franchise/quotes", a.createQuote)
	mux.HandleFunc("GET /v1/franchise/quotes/result", a.quoteResult)
	mux.HandleFunc("POST /v1/franchise/delivery-checklists", a.publishDeliveryChecklist)
	mux.HandleFunc("GET /v1/franchise/delivery-checklists/result", a.publishedDeliveryChecklist)
	mux.HandleFunc("POST /v1/franchise/handovers/{id}/complete-checklist", a.completeDeliveryChecklist)
	mux.HandleFunc("GET /v1/franchise/handovers/{id}/checklist-result", a.deliveryChecklistCompletion)
	mux.HandleFunc("GET /v1/franchise/delivery-exceptions", a.deliveryExceptions)
	mux.HandleFunc("POST /v1/franchise/delivery-exceptions/{id}/resolve", a.resolveDeliveryException)
	mux.HandleFunc("GET /v1/franchise/returns", a.returnCases)
	mux.HandleFunc("GET /v1/franchise/returns/result", a.returnCaseResult)
	mux.HandleFunc("GET /v1/franchise/returns/{id}/outcome", a.returnOutcome)
	mux.HandleFunc("POST /v1/franchise/return-authorizations/{id}/receive", a.receiveReturn)
	mux.HandleFunc("POST /v1/franchise/return-receipts/{id}/decide", a.decideReturn)
	mux.HandleFunc("GET /v1/customer/journey", a.customerJourney)
	mux.HandleFunc("POST /v1/customer/appointments/{id}/cancel", a.cancelCustomerAppointment)
	mux.HandleFunc("POST /v1/customer/quotes/{id}/accept", a.acceptQuote)
	mux.HandleFunc("POST /v1/customer/handovers/{id}/accept", a.acceptHandover)
	mux.HandleFunc("POST /v1/customer/handovers/{id}/reject", a.rejectHandover)
}

func (a franchiseJourneyAPI) publishDeliveryChecklist(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string                           `json:"organization_id"`
		ChecklistID    string                           `json:"checklist_id"`
		Version        int64                            `json:"version"`
		Title          string                           `json:"title"`
		Items          []franchisejourney.ChecklistItem `json:"items"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "handover:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.PublishDeliveryChecklist(r.Context(), p.TenantID, p.Subject, franchisejourney.DeliveryChecklist{ID: input.ChecklistID, OrganizationID: input.OrganizationID, Version: input.Version, Title: input.Title, Items: input.Items})
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_DELIVERY_CHECKLIST", "delivery checklist does not match contract")
		return
	}
	if errors.Is(err, franchisejourney.ErrConflict) {
		writeProblem(w, 409, "DELIVERY_CHECKLIST_CONFLICT", "delivery checklist version already exists or violates scope")
		return
	}
	if err != nil {
		writeProblem(w, 500, "DELIVERY_CHECKLIST_FAILED", "delivery checklist publication failed")
		return
	}
	writeJSON(w, 201, value)
}

func (a franchiseJourneyAPI) completeDeliveryChecklist(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID   string                               `json:"organization_id"`
		Version          int64                                `json:"version"`
		ChecklistID      string                               `json:"checklist_id"`
		ChecklistVersion int64                                `json:"checklist_version"`
		Responses        []franchisejourney.ChecklistResponse `json:"responses"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "handover:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.CompleteDeliveryChecklist(r.Context(), p.TenantID, input.OrganizationID, p.Subject, r.PathValue("id"), input.Version, input.ChecklistID, input.ChecklistVersion, input.Responses)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) deliveryExceptions(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "handover:manage", organization)
	if !ok {
		return
	}
	items, err := a.service.DeliveryExceptions(r.Context(), p.TenantID, organization, 100)
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "delivery exceptions query failed")
		return
	}
	writeJSON(w, 200, map[string]any{"items": items})
}

func (a franchiseJourneyAPI) resolveDeliveryException(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		Version        int64  `json:"version"`
		Action         string `json:"action"`
		Notes          string `json:"notes"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "handover:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.ResolveDeliveryException(r.Context(), p.TenantID, input.OrganizationID, p.Subject, r.PathValue("id"), input.Version, input.Action, input.Notes)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) returnCases(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "handover:manage", organization)
	if !ok {
		return
	}
	items, err := a.service.ReturnCases(r.Context(), p.TenantID, organization, 100)
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "return cases query failed")
		return
	}
	writeJSON(w, 200, map[string]any{"items": items})
}

func (a franchiseJourneyAPI) receiveReturn(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		SerialNumber   string `json:"serial_number"`
		ConditionCode  string `json:"condition_code"`
		Notes          string `json:"notes"`
		EvidenceSHA256 string `json:"evidence_sha256"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "handover:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.ReceiveReturn(r.Context(), p.TenantID, input.OrganizationID, p.Subject, r.PathValue("id"), input.SerialNumber, input.ConditionCode, input.Notes, input.EvidenceSHA256)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) decideReturn(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID  string `json:"organization_id"`
		InventoryAction string `json:"inventory_action"`
		Notes           string `json:"notes"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "handover:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.DecideReturn(r.Context(), p.TenantID, input.OrganizationID, p.Subject, r.PathValue("id"), input.InventoryAction, input.Notes)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) publicAppointmentSlots(w http.ResponseWriter, r *http.Request) {
	from, fromErr := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
	to, toErr := time.Parse(time.RFC3339, r.URL.Query().Get("to"))
	if fromErr != nil || toErr != nil {
		writeProblem(w, 400, "INVALID_SLOT_RANGE", "from and to must be RFC3339 timestamps")
		return
	}
	items, err := a.service.PublicAppointmentSlots(r.Context(), r.PathValue("tenantCode"), r.PathValue("organizationCode"), r.URL.Query().Get("kind"), from, to)
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_SLOT_QUERY", "slot query does not match contract")
		return
	}
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "appointment slot query failed")
		return
	}
	writeJSON(w, 200, map[string]any{"items": items})
}

func (a franchiseJourneyAPI) createAppointmentSlot(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string    `json:"organization_id"`
		Kind           string    `json:"kind"`
		StartsAt       time.Time `json:"starts_at"`
		EndsAt         time.Time `json:"ends_at"`
		Capacity       int       `json:"capacity"`
	}
	hash, decoded := decodeHashedJSON(w, r, &input)
	if !decoded {
		return
	}
	p, ok := a.protected(w, r, "appointment:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, replay, err := a.service.CreateAppointmentSlotOnce(r.Context(), p.TenantID, p.Subject, r.Header.Get("Idempotency-Key"), hash, franchisejourney.AppointmentSlot{OrganizationID: input.OrganizationID, Kind: input.Kind, StartsAt: input.StartsAt, EndsAt: input.EndsAt, Capacity: input.Capacity})
	if errors.Is(err, franchisejourney.ErrConflict) {
		writeProblem(w, 409, "SLOT_CONFLICT", "slot overlaps existing capacity or organization is unavailable")
		return
	}
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_SLOT", "slot does not match contract")
		return
	}
	if err != nil {
		writeProblem(w, 500, "SLOT_FAILED", "appointment slot could not be persisted")
		return
	}
	status := 201
	if replay {
		status = 200
		w.Header().Set("Idempotency-Replayed", "true")
	}
	writeJSON(w, status, value)
}

func (a franchiseJourneyAPI) createServiceResource(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID   string   `json:"organization_id"`
		PrincipalSubject string   `json:"principal_subject"`
		DisplayName      string   `json:"display_name"`
		Kind             string   `json:"kind"`
		Skills           []string `json:"skills"`
	}
	hash, decoded := decodeHashedJSON(w, r, &input)
	if !decoded {
		return
	}
	p, ok := a.protected(w, r, "resource:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, replay, err := a.service.CreateServiceResourceOnce(r.Context(), p.TenantID, p.Subject, r.Header.Get("Idempotency-Key"), hash, franchisejourney.ServiceResource{OrganizationID: input.OrganizationID, PrincipalSubject: input.PrincipalSubject, DisplayName: input.DisplayName, Kind: input.Kind, Skills: input.Skills})
	if errors.Is(err, franchisejourney.ErrConflict) {
		writeProblem(w, 409, "RESOURCE_CONFLICT", "resource identity or organization conflicts with this request")
		return
	}
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_RESOURCE", "resource does not match contract")
		return
	}
	if err != nil {
		writeProblem(w, 500, "RESOURCE_FAILED", "resource could not be persisted")
		return
	}
	status := 201
	if replay {
		status = 200
		w.Header().Set("Idempotency-Replayed", "true")
	}
	writeJSON(w, status, value)
}

func (a franchiseJourneyAPI) assignAppointmentResource(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		ResourceID     string `json:"resource_id"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "appointment:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.AssignAppointmentResource(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.ResourceID, input.Version)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) transitionAppointment(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		Current        string `json:"current"`
		Target         string `json:"target"`
		Version        int64  `json:"version"`
		ReasonCode     string `json:"reason_code"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "appointment:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.TransitionAppointment(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Current, input.Target, input.Version, p.Subject, input.ReasonCode)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) createAvailability(w http.ResponseWriter, r *http.Request) {
	var input franchisejourney.AvailabilityEntry
	hash, ok := decodeHashedJSON(w, r, &input)
	if !ok {
		return
	}
	p, ok := a.protected(w, r, "availability:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, replayed, err := a.service.CreateAvailabilityOnce(r.Context(), p.TenantID, p.Subject, r.Header.Get("Idempotency-Key"), hash, input)
	if errors.Is(err, franchisejourney.ErrConflict) {
		writeProblem(w, 409, "AVAILABILITY_CONFLICT", "availability overlaps or conflicts with active appointments")
		return
	}
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_AVAILABILITY", "availability does not match contract")
		return
	}
	if err != nil {
		writeProblem(w, 500, "AVAILABILITY_FAILED", "availability could not be persisted")
		return
	}
	status := 201
	if replayed {
		status = 200
		w.Header().Set("Idempotency-Replayed", "true")
	}
	writeJSON(w, status, value)
}

func (a franchiseJourneyAPI) cancelAvailability(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		Version        int64  `json:"version"`
		ReasonCode     string `json:"reason_code"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "availability:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.CancelAvailability(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Version, p.Subject, input.ReasonCode)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) availability(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "availability:read", organization)
	if !ok {
		return
	}
	from, fromErr := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
	to, toErr := time.Parse(time.RFC3339, r.URL.Query().Get("to"))
	if fromErr != nil || toErr != nil {
		writeProblem(w, 400, "INVALID_AVAILABILITY_RANGE", "from and to must be RFC3339 timestamps")
		return
	}
	items, err := a.service.Availability(r.Context(), p.TenantID, organization, r.URL.Query().Get("resource_id"), from, to)
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_AVAILABILITY_RANGE", "availability range does not match contract")
		return
	}
	if err != nil {
		writeProblem(w, 500, "AVAILABILITY_QUERY_FAILED", "availability query failed")
		return
	}
	writeJSON(w, 200, map[string]any{"items": items})
}

func (a franchiseJourneyAPI) cancelCustomerAppointment(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		Version        int64  `json:"version"`
		ReasonCode     string `json:"reason_code"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "customer:self", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.CancelCustomerAppointment(r.Context(), p.TenantID, input.OrganizationID, p.Subject, r.PathValue("id"), input.Version, input.ReasonCode)
	writeJourneyResult(w, value, err)
}

type franchiseJourneyAPI struct {
	service  *franchisejourney.Service
	verifier identity.Verifier
}

func (a franchiseJourneyAPI) protected(w http.ResponseWriter, r *http.Request, permission, organization string) (identity.Principal, bool) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return p, false
	}
	if !p.AllowedOrganization(organization) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return p, false
	}
	return p, true
}

func (a franchiseJourneyAPI) publicLocations(w http.ResponseWriter, r *http.Request) {
	items, err := a.service.PublicLocations(r.Context(), r.PathValue("tenantCode"))
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_TENANT", "tenant code is invalid")
		return
	}
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "locations query failed")
		return
	}
	writeJSON(w, 200, map[string]any{"items": items})
}

func (a franchiseJourneyAPI) appointmentAgenda(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "appointment:manage", organization)
	if !ok {
		return
	}
	from, fromErr := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
	to, toErr := time.Parse(time.RFC3339, r.URL.Query().Get("to"))
	if fromErr != nil || toErr != nil {
		writeProblem(w, 400, "INVALID_AGENDA_RANGE", "explicit RFC3339 range required")
		return
	}
	value, err := a.service.AppointmentAgenda(r.Context(), p.TenantID, organization, from, to)
	w.Header().Set("Cache-Control", "no-store")
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) requestAppointment(w http.ResponseWriter, r *http.Request) {
	var input franchisejourney.Appointment
	hash, ok := decodeHashedJSON(w, r, &input)
	if !ok {
		return
	}
	key := r.Header.Get("Idempotency-Key")
	value, replayed, err := a.service.RequestAppointment(r.Context(), r.PathValue("tenantCode"), r.PathValue("organizationCode"), key, hash, input)
	if errors.Is(err, franchisejourney.ErrNotFound) {
		writeProblem(w, 404, "LOCATION_NOT_FOUND", "public location not found")
		return
	}
	if errors.Is(err, franchisejourney.ErrConflict) {
		writeProblem(w, 409, "APPOINTMENT_CONFLICT", "idempotency key or appointment capacity conflicts with this request")
		return
	}
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_APPOINTMENT", "appointment does not match contract")
		return
	}
	if err != nil {
		writeProblem(w, 500, "APPOINTMENT_FAILED", "appointment could not be persisted")
		return
	}
	status := 202
	if replayed {
		status = 200
		w.Header().Set("Idempotency-Replayed", "true")
	}
	writeJSON(w, status, value)
}

func (a franchiseJourneyAPI) leads(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "lead:read", organization)
	if !ok {
		return
	}
	limit, after, ok := pageInput(w, r)
	if !ok {
		return
	}
	value, err := a.service.Leads(r.Context(), p.TenantID, organization, limit, after)
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "leads query failed")
		return
	}
	writeJSON(w, 200, value)
}

func (a franchiseJourneyAPI) assignLead(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID  string `json:"organization_id"`
		AssignedSubject string `json:"assigned_subject"`
		Version         int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "lead:assign", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.AssignLeadAs(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.AssignedSubject, input.Version, p.Subject)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) transitionLead(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		Current        string `json:"current"`
		Target         string `json:"target"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "lead:update", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.TransitionLeadAs(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Current, input.Target, input.Version, p.Subject)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) createQuote(w http.ResponseWriter, r *http.Request) {
	var input franchisejourney.Quote
	hash, ok := decodeHashedJSON(w, r, &input)
	if !ok {
		return
	}
	p, ok := a.protected(w, r, "quote:write", input.OrganizationID)
	if !ok {
		return
	}
	value, replayed, err := a.service.CreateQuoteAs(r.Context(), p.TenantID, p.Subject, r.Header.Get("Idempotency-Key"), hash, input)
	if errors.Is(err, franchisejourney.ErrConflict) {
		writeProblem(w, 409, "PRICE_OR_LEAD_CONFLICT", "active server price or lead not found")
		return
	}
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_QUOTE", "quote does not match contract")
		return
	}
	if err != nil {
		writeProblem(w, 500, "QUOTE_FAILED", "quote could not be persisted")
		return
	}
	status := 201
	if replayed {
		status = 200
	}
	writeJSON(w, status, value)
}

func (a franchiseJourneyAPI) customerJourney(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "customer:self", organization)
	if !ok {
		return
	}
	value, err := a.service.CustomerJourney(r.Context(), p.TenantID, organization, p.Subject)
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "customer journey query failed")
		return
	}
	writeJSON(w, 200, value)
}

func (a franchiseJourneyAPI) acceptHandover(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID    string `json:"organization_id"`
		Version           int64  `json:"version"`
		ConfirmedReceived bool   `json:"confirmed_received"`
		SerialNumber      string `json:"serial_number"`
		ChecklistID       string `json:"checklist_id"`
		ChecklistVersion  int64  `json:"checklist_version"`
	}
	evidence, decoded := decodeHashedJSON(w, r, &input)
	if !decoded {
		return
	}
	p, ok := a.protected(w, r, "customer:self", input.OrganizationID)
	if !ok {
		return
	}
	if !input.ConfirmedReceived {
		writeProblem(w, 400, "HANDOVER_CONFIRMATION_REQUIRED", "confirmed_received must be true")
		return
	}
	value, err := a.service.AcceptHandover(r.Context(), p.TenantID, input.OrganizationID, p.Subject, r.PathValue("id"), input.Version, input.SerialNumber, input.ChecklistID, input.ChecklistVersion, evidence)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) rejectHandover(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		Version        int64  `json:"version"`
		ReasonCode     string `json:"reason_code"`
		Details        string `json:"details"`
	}
	evidence, decoded := decodeHashedJSON(w, r, &input)
	if !decoded {
		return
	}
	p, ok := a.protected(w, r, "customer:self", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.RejectHandover(r.Context(), p.TenantID, input.OrganizationID, p.Subject, r.PathValue("id"), input.Version, input.ReasonCode, input.Details, evidence)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) acceptQuote(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		Version        int64  `json:"version"`
	}
	evidence, ok := decodeHashedJSON(w, r, &input)
	if !ok {
		return
	}
	p, ok := a.protected(w, r, "customer:self", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.AcceptQuote(r.Context(), p.TenantID, input.OrganizationID, p.Subject, r.PathValue("id"), input.Version, evidence)
	writeJourneyResult(w, value, err)
}

func decodeHashedJSON(w http.ResponseWriter, r *http.Request, destination any) (string, bool) {
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return "", false
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeProblem(w, 400, "INVALID_JSON", "body is invalid or too large")
		return "", false
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	if !decodeStrict(w, r, destination) {
		return "", false
	}
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:]), true
}

func writeJourneyResult(w http.ResponseWriter, value any, err error) {
	if errors.Is(err, franchisejourney.ErrConflict) {
		writeProblem(w, 409, "JOURNEY_CONFLICT", "state or version conflict")
		return
	}
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_JOURNEY_COMMAND", "command does not match contract")
		return
	}
	if err != nil {
		writeProblem(w, 500, "JOURNEY_FAILED", "command failed")
		return
	}
	writeJSON(w, 200, value)
}

func (a franchiseJourneyAPI) quoteResult(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "quote:write", organization)
	if !ok {
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	value, err := a.service.QuoteResult(r.Context(), p.TenantID, organization, r.URL.Query().Get("lead_id"), r.URL.Query().Get("request_key"))
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) availabilityCreationResult(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "availability:manage", organization)
	if !ok {
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	value, err := a.service.AvailabilityCreationResult(r.Context(), p.TenantID, organization, r.URL.Query().Get("request_key"))
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) serviceResourceCreationResult(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "resource:manage", organization)
	if !ok {
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	value, err := a.service.ServiceResourceCreationResult(r.Context(), p.TenantID, organization, r.URL.Query().Get("request_key"))
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) appointmentSlotCreationResult(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "appointment:manage", organization)
	if !ok {
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	value, err := a.service.AppointmentSlotCreationResult(r.Context(), p.TenantID, organization, r.URL.Query().Get("request_key"))
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) publishedDeliveryChecklist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "handover:manage", organization)
	if !ok {
		return
	}
	version, err := strconv.ParseInt(r.URL.Query().Get("version"), 10, 64)
	if err != nil || version < 1 {
		writeProblem(w, 400, "INVALID_DELIVERY_CHECKLIST_QUERY", "checklist version is invalid")
		return
	}
	value, err := a.service.PublishedDeliveryChecklist(r.Context(), p.TenantID, organization, r.URL.Query().Get("checklist_id"), version)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) deliveryChecklistCompletion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "handover:manage", organization)
	if !ok {
		return
	}
	value, err := a.service.DeliveryChecklistCompletion(r.Context(), p.TenantID, organization, r.PathValue("id"))
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) returnCaseResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "handover:manage", organization)
	if !ok {
		return
	}
	value, err := a.service.ReturnCaseResult(r.Context(), p.TenantID, organization, r.URL.Query().Get("authorization_id"))
	writeJourneyResult(w, value, err)
}
