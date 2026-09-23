package httpapi

import (
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"net/http"
)

type CommerceModule struct {
	Service                                *commerce.Service
	PaymentProvider                        string
	PaymentTenantID, PaymentOrganizationID string
	ProviderObservedPayments               bool
}

func (m CommerceModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	api := commerceAPI{service: m.Service, verifier: verifier, paymentProvider: m.PaymentProvider, paymentTenantID: m.PaymentTenantID, paymentOrganizationID: m.PaymentOrganizationID, providerObservedPayments: m.ProviderObservedPayments}
	mux.HandleFunc("GET /v1/commerce/orders", api.operations)
	mux.HandleFunc("GET /v1/public/{tenantCode}/prices", api.publicPrice)
	mux.HandleFunc("POST /v1/pricing/books", api.createBook)
	mux.HandleFunc("POST /v1/pricing/books/{id}/activate", api.activateBook)
	mux.HandleFunc("POST /v1/commerce/orders/{id}/lines", api.addLine)
	mux.HandleFunc("POST /v1/commerce/orders/{id}/place", api.placeOrder)
	mux.HandleFunc("POST /v1/commerce/orders/{id}/allocations", api.allocate)
	mux.HandleFunc("POST /v1/commerce/orders/{id}/payments", api.createPayment)
	mux.HandleFunc("POST /v1/commerce/orders/{id}/payment-request", api.requestOrderPayment)
	mux.HandleFunc("POST /v1/payments/{id}/transitions", api.transitionPayment)
}

type commerceAPI struct {
	service                                *commerce.Service
	verifier                               identity.Verifier
	paymentProvider                        string
	paymentTenantID, paymentOrganizationID string
	providerObservedPayments               bool
}

func (a commerceAPI) operations(w http.ResponseWriter, r *http.Request) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return
	}
	if !p.Allowed("inventory:allocate") && !p.Allowed("payment:create") && !p.Allowed("handover:manage") && !p.Allowed("admin:read") {
		writeProblem(w, 403, "FORBIDDEN", "operational permission is required")
		return
	}
	organization := r.URL.Query().Get("organization_id")
	if organization == "" || !p.AllowedOrganization(organization) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "organization scope is required")
		return
	}
	value, err := a.service.Operations(r.Context(), p.TenantID, organization)
	if err != nil {
		writeProblem(w, 503, "OPERATIONS_UNAVAILABLE", "operation snapshot unavailable")
		return
	}
	if (a.paymentProvider == "stripe" || a.paymentProvider == "mercadopago") && (a.paymentTenantID == "" || (a.paymentTenantID == p.TenantID && a.paymentOrganizationID == organization)) {
		value.PaymentProvider = a.paymentProvider
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, 200, value)
}

func (a commerceAPI) auth(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
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
func (a commerceAPI) publicPrice(w http.ResponseWriter, r *http.Request) {
	value, err := a.service.PublicPrice(r.Context(), r.PathValue("tenantCode"), r.URL.Query().Get("market"), r.URL.Query().Get("variant_id"))
	if err != nil {
		writeProblem(w, 404, "PRICE_NOT_FOUND", "one current price was not found")
		return
	}
	writeJSON(w, 200, value)
}
func (a commerceAPI) createBook(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "pricing:write")
	if !ok {
		return
	}
	var input commerce.PriceBook
	if !decodeStrict(w, r, &input) {
		return
	}
	value, err := a.service.CreatePriceBook(r.Context(), p.TenantID, input)
	if err != nil {
		writeProblem(w, 400, "INVALID_PRICE_BOOK", "price book does not match contract")
		return
	}
	writeJSON(w, 201, value)
}
func (a commerceAPI) activateBook(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "pricing:write")
	if !ok {
		return
	}
	writeCommerceResult(w, a.service.ActivatePriceBook(r.Context(), p.TenantID, r.PathValue("id")))
}
func (a commerceAPI) addLine(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "order:write")
	if !ok {
		return
	}
	var input struct {
		OrganizationID  string `json:"organization_id"`
		PriceBookID     string `json:"price_book_id"`
		VariantID       string `json:"variant_id"`
		Quantity        int    `json:"quantity"`
		ExpectedVersion int64  `json:"expected_version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	value, err := a.service.AddOrderLine(r.Context(), p.TenantID, commerce.OrderLine{OrderID: r.PathValue("id"), OrganizationID: input.OrganizationID, PriceBookID: input.PriceBookID, VariantID: input.VariantID, Quantity: input.Quantity}, input.ExpectedVersion)
	if err != nil {
		writeCommerceResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a commerceAPI) placeOrder(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "order:write")
	if !ok {
		return
	}
	var input struct {
		OrganizationID  string `json:"organization_id"`
		ExpectedVersion int64  `json:"expected_version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeCommerceResult(w, a.service.PlaceOrder(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.ExpectedVersion))
}
func (a commerceAPI) allocate(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "inventory:allocate")
	if !ok {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		LineID         string `json:"line_id"`
		StockUnitID    string `json:"stock_unit_id"`
		OrderVersion   int64  `json:"order_version"`
		StockVersion   int64  `json:"stock_version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	writeCommerceResult(w, a.service.AllocateStockAs(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.LineID, input.StockUnitID, input.OrderVersion, input.StockVersion, p.Subject))
}
func (a commerceAPI) paymentEnabled(w http.ResponseWriter) bool {
	if a.paymentProvider != "stripe" && a.paymentProvider != "mercadopago" {
		writeProblem(w, 503, "PAYMENT_REQUEST_DISABLED", "payment provider selection is required; no payment effect enabled")
		return false
	}
	return true
}
func (a commerceAPI) paymentScope(w http.ResponseWriter, p identity.Principal, organization string) bool {
	if (a.paymentTenantID != "" || a.paymentOrganizationID != "") && (a.paymentTenantID != p.TenantID || a.paymentOrganizationID != organization) {
		writeProblem(w, 503, "PAYMENT_SCOPE_DISABLED", "payment checkout is not enabled for this scope")
		return false
	}
	return true
}
func (a commerceAPI) requestOrderPayment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "payment:create")
	if !ok {
		return
	}
	if !a.paymentEnabled(w) {
		return
	}
	var input struct {
		OrganizationID string `json:"organization_id"`
		RequestKey     string `json:"request_key"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "organization scope is required")
		return
	}
	if !a.paymentScope(w, p, input.OrganizationID) {
		return
	}
	v, err := a.service.RequestOrderPayment(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), a.paymentProvider, input.RequestKey, p.Subject)
	if err != nil {
		writeCommerceResult(w, err)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, 202, v)
}
func (a commerceAPI) createPayment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "payment:create")
	if !ok {
		return
	}
	if !a.paymentEnabled(w) {
		return
	}
	var input commerce.PaymentAttempt
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	input.OrderID = r.PathValue("id")
	if !a.paymentScope(w, p, input.OrganizationID) {
		return
	}
	if input.ProviderCode != a.paymentProvider {
		writeProblem(w, 400, "PAYMENT_PROVIDER_MISMATCH", "provider must match configured selection")
		return
	}
	value, err := a.service.CreatePaymentAttemptAs(r.Context(), p.TenantID, r.Header.Get("Idempotency-Key"), input, p.Subject)
	if err != nil {
		writeCommerceResult(w, err)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, 202, value)
}
func (a commerceAPI) transitionPayment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "payment:write")
	if !ok {
		return
	}
	if !a.paymentEnabled(w) {
		return
	}
	var input struct {
		OrganizationID    string `json:"organization_id"`
		Current           string `json:"current"`
		Target            string `json:"target"`
		Version           int64  `json:"version"`
		ProviderReference string `json:"provider_reference"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !p.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	if !a.paymentScope(w, p, input.OrganizationID) {
		return
	}
	// In the provider-observed profile, an operator's JSON cannot establish
	// authorization, capture, refund or dispute facts. Only a never-dispatched
	// created request may be cancelled here; the repository CAS rejects it if
	// the atomic send claim has already changed it to pending.
	if a.providerObservedPayments && (input.Current != "created" || input.Target != "failed" || input.ProviderReference != "") {
		writeProblem(w, 403, "PROVIDER_OBSERVATION_REQUIRED", "financial provider state requires authenticated reconciliation")
		return
	}
	writeCommerceResult(w, a.service.TransitionPayment(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Current, input.Target, input.Version, input.ProviderReference))
}
func writeCommerceResult(w http.ResponseWriter, err error) {
	if errors.Is(err, commerce.ErrPaymentUnavailable) {
		writeProblem(w, 503, "PAYMENT_INTENT_UNAVAILABLE", "result unavailable; consult persisted state before retrying")
		return
	}
	if errors.Is(err, commerce.ErrConflict) {
		writeProblem(w, 409, "COMMERCE_CONFLICT", "state, version or business invariant conflict")
		return
	}
	if err != nil {
		writeProblem(w, 400, "INVALID_COMMERCE_COMMAND", "command does not match contract")
		return
	}
	writeJSON(w, 200, map[string]string{"status": "accepted"})
}
