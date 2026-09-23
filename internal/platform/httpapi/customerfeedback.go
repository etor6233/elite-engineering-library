package httpapi

import (
	"context"
	"elite.local/enterprise/internal/customerfeedback"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"time"
)

// CustomerFeedbackModule is opt-in. The existing verifier remains the identity authority.
type CustomerFeedbackModule struct{ Service *customerfeedback.Service }

func (m CustomerFeedbackModule) Register(mux *http.ServeMux, v identity.Verifier) {
	if m.Service == nil {
		return
	}
	a := feedbackAPI{service: m.Service, verifier: v}
	a.registerMetrics(mux)
	mux.HandleFunc("GET /v1/customer/surveys/{survey}", a.definition)
	mux.HandleFunc("POST /v1/customer/surveys/{survey}/response", a.submit)
	mux.HandleFunc("GET /v1/customer/surveys/{survey}/response", a.ownAnswer)
	mux.HandleFunc("GET /v1/admin/surveys/{survey}/summary", a.summary)
}

type feedbackAPI struct {
	service  *customerfeedback.Service
	verifier identity.Verifier
}

func (a feedbackAPI) scope(w http.ResponseWriter, r *http.Request, permission string) (customerfeedback.Scope, bool) {
	p, org, ok := (enterpriseQueryAPI{verifier: a.verifier}).scope(w, r, permission)
	if !ok {
		return customerfeedback.Scope{}, false
	}
	query := r.URL.Query()
	if len(query) != 1 || len(query["organization_id"]) != 1 {
		writeProblem(w, 400, "INVALID_QUERY", "one organization_id is required")
		return customerfeedback.Scope{}, false
	}
	return customerfeedback.Scope{Tenant: p.TenantID, Organization: org, Customer: p.Subject}, true
}
func feedbackError(w http.ResponseWriter, e error) {
	switch {
	case errors.Is(e, customerfeedback.ErrInvalid):
		writeProblem(w, 400, "INVALID_SURVEY_REQUEST", "request does not match the survey contract")
	case errors.Is(e, customerfeedback.ErrConflict):
		writeProblem(w, 409, "SURVEY_RESPONSE_CONFLICT", "a different response is already stored; read your response")
	case errors.Is(e, customerfeedback.ErrUnavailable), errors.Is(e, customerfeedback.ErrNotFound):
		writeProblem(w, 404, "SURVEY_NOT_AVAILABLE", "survey or response is not available")
	default:
		writeProblem(w, 503, "SURVEY_UNAVAILABLE", "survey operation is unavailable; recover with a read")
	}
}
func (a feedbackAPI) definition(w http.ResponseWriter, r *http.Request) {
	s, ok := a.scope(w, r, "customer:self")
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	d, e := a.service.Definition(ctx, s, r.PathValue("survey"))
	if e != nil {
		feedbackError(w, e)
		return
	}
	writeJSON(w, 200, d)
}
func (a feedbackAPI) ownAnswer(w http.ResponseWriter, r *http.Request) {
	s, ok := a.scope(w, r, "customer:self")
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	d, e := a.service.OwnAnswer(ctx, s, r.PathValue("survey"))
	if e != nil {
		feedbackError(w, e)
		return
	}
	writeJSON(w, 200, d)
}
func (a feedbackAPI) summary(w http.ResponseWriter, r *http.Request) {
	s, ok := a.scope(w, r, "surveys:read")
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	d, e := a.service.Summary(ctx, s, r.PathValue("survey"))
	if e != nil {
		feedbackError(w, e)
		return
	}
	writeJSON(w, 200, d)
}
func readSurveySubmission(reader io.Reader) (customerfeedback.Submission, error) {
	var input customerfeedback.Submission
	d := json.NewDecoder(reader)
	token, e := d.Token()
	if e != nil || token != json.Delim('{') {
		return input, customerfeedback.ErrInvalid
	}
	seen := map[string]bool{}
	for d.More() {
		token, e = d.Token()
		key, ok := token.(string)
		if e != nil || !ok || seen[key] {
			return input, customerfeedback.ErrInvalid
		}
		seen[key] = true
		switch key {
		case "score":
			var value *int
			if e = d.Decode(&value); e != nil || value == nil {
				return input, customerfeedback.ErrInvalid
			}
			input.Score = *value
		case "consent":
			var value *bool
			if e = d.Decode(&value); e != nil || value == nil {
				return input, customerfeedback.ErrInvalid
			}
			input.Consent = *value
		case "consent_version":
			var value *string
			if e = d.Decode(&value); e != nil || value == nil {
				return input, customerfeedback.ErrInvalid
			}
			input.ConsentVersion = *value
		default:
			return input, customerfeedback.ErrInvalid
		}
	}
	token, e = d.Token()
	if e != nil || token != json.Delim('}') || len(seen) != 3 {
		return input, customerfeedback.ErrInvalid
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		return input, customerfeedback.ErrInvalid
	}
	return input, nil
}
func (a feedbackAPI) submit(w http.ResponseWriter, r *http.Request) {
	s, ok := a.scope(w, r, "customer:self")
	if !ok {
		return
	}
	kind, _, e := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if e != nil || kind != "application/json" || (r.Header.Get("Content-Encoding") != "" && r.Header.Get("Content-Encoding") != "identity") {
		writeProblem(w, 415, "INVALID_CONTENT_TYPE", "application/json is required")
		return
	}
	input, e := readSurveySubmission(http.MaxBytesReader(w, r.Body, 2048))
	if e != nil {
		feedbackError(w, e)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	result, e := a.service.Submit(ctx, s, r.PathValue("survey"), input)
	if e != nil {
		feedbackError(w, e)
		return
	}
	status := 201
	if result.Replay {
		status = 200
		w.Header().Set("Idempotency-Replayed", "true")
	}
	writeJSON(w, status, result)
}
