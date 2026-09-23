package trainingbridge

import (
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type Module struct{ Store *Store }

func trainingJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func trainingError(w http.ResponseWriter, e error) {
	code, status := "UNAVAILABLE", 503
	switch {
	case errors.Is(e, ErrScope):
		code, status = "NOT_FOUND", 404
	case errors.Is(e, ErrContract):
		code, status = "INVALID_CONTRACT", 400
	case errors.Is(e, ErrConflict), errors.Is(e, approval.ErrDuplicate), errors.Is(e, approval.ErrNotPending), errors.Is(e, approval.ErrSeparation):
		code, status = "CONSULT_RECORDED_STATE", 409
	}
	trainingJSON(w, status, map[string]string{"code": code})
}
func (m Module) principal(w http.ResponseWriter, r *http.Request, v identity.Verifier, permission string) (identity.Principal, bool) {
	values := r.Header.Values("Authorization")
	if v == nil || len(values) != 1 || !strings.HasPrefix(values[0], "Bearer ") || len(values[0]) > 16391 || len(strings.Fields(values[0])) != 2 {
		trainingJSON(w, 401, map[string]string{"code": "UNAUTHENTICATED"})
		return identity.Principal{}, false
	}
	p, e := v.Verify(r.Context(), strings.TrimPrefix(values[0], "Bearer "))
	if e != nil {
		trainingJSON(w, 401, map[string]string{"code": "UNAUTHENTICATED"})
		return p, false
	}
	allowed := m.Store != nil && (m.Store.allowed(p, permission) || (permission == "training:view" && (m.Store.allowed(p, "training:learn") || m.Store.allowed(p, "training:review"))))
	if !allowed {
		trainingJSON(w, 403, map[string]string{"code": "FORBIDDEN"})
		return p, false
	}
	if r.URL.RawQuery != "" {
		trainingJSON(w, 400, map[string]string{"code": "INVALID_QUERY"})
		return p, false
	}
	return p, true
}
func trainingBody(w http.ResponseWriter, r *http.Request, out any) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return ErrContract
	}
	b, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
	if e != nil {
		return ErrContract
	}
	return strict(b, out)
}
func (m Module) Register(mux *http.ServeMux, v identity.Verifier) {
	mux.HandleFunc("GET /v1/training/courses", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "training:view")
		if !ok {
			return
		}
		value, e := m.Store.Courses(p)
		if e != nil {
			trainingError(w, e)
			return
		}
		trainingJSON(w, 200, value)
	})
	mux.HandleFunc("GET /v1/training/assessments", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "training:view")
		if !ok {
			return
		}
		value, e := m.Store.Assessments(r.Context(), p)
		if e != nil {
			trainingError(w, e)
			return
		}
		trainingJSON(w, 200, value)
	})
	mux.HandleFunc("GET /v1/training/assessments/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "training:view")
		if !ok {
			return
		}
		value, e := m.Store.Assessment(r.Context(), p, r.PathValue("id"))
		if e != nil {
			trainingError(w, e)
			return
		}
		trainingJSON(w, 200, value)
	})
	mux.HandleFunc("GET /v1/training/attempts/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "training:view")
		if !ok {
			return
		}
		value, e := m.Store.Read(r.Context(), p, r.PathValue("id"))
		if e != nil {
			trainingError(w, e)
			return
		}
		trainingJSON(w, 200, value)
	})
	mux.HandleFunc("POST /v1/training/attempts", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "training:learn")
		if !ok {
			return
		}
		var body struct {
			ID         string `json:"attempt_id"`
			Course     string `json:"course_id"`
			ProfileSHA string `json:"profile_sha256"`
		}
		if e := trainingBody(w, r, &body); e != nil {
			trainingError(w, e)
			return
		}
		value, e := m.Store.Start(r.Context(), p, body.ID, body.Course, body.ProfileSHA)
		if e != nil {
			trainingError(w, e)
			return
		}
		trainingJSON(w, 201, value)
	})
	mux.HandleFunc("POST /v1/training/attempts/{id}/acknowledgements", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "training:learn")
		if !ok {
			return
		}
		var body struct {
			Lesson     string `json:"lesson_id"`
			ProfileSHA string `json:"profile_sha256"`
		}
		if e := trainingBody(w, r, &body); e != nil {
			trainingError(w, e)
			return
		}
		value, e := m.Store.Acknowledge(r.Context(), p, r.PathValue("id"), body.Lesson, body.ProfileSHA)
		if e != nil {
			trainingError(w, e)
			return
		}
		trainingJSON(w, 200, value)
	})
	mux.HandleFunc("POST /v1/training/attempts/{id}/submissions", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "training:learn")
		if !ok {
			return
		}
		var body struct {
			Answers    map[string]string `json:"answers"`
			ProfileSHA string            `json:"profile_sha256"`
		}
		if e := trainingBody(w, r, &body); e != nil {
			trainingError(w, e)
			return
		}
		value, e := m.Store.Submit(r.Context(), p, r.PathValue("id"), body.ProfileSHA, body.Answers)
		if e != nil {
			trainingError(w, e)
			return
		}
		trainingJSON(w, 200, value)
	})
	mux.HandleFunc("POST /v1/training/assessments/{id}/decision", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "training:review")
		if !ok {
			return
		}
		var body struct {
			SHA      string `json:"payload_sha256"`
			Approved *bool  `json:"approved"`
			Reason   string `json:"reason"`
		}
		if e := trainingBody(w, r, &body); e != nil || body.Approved == nil {
			trainingError(w, ErrContract)
			return
		}
		value, e := m.Store.Assess(r.Context(), p, r.PathValue("id"), body.SHA, *body.Approved, body.Reason)
		if e != nil {
			trainingError(w, e)
			return
		}
		trainingJSON(w, 200, value)
	})
}
