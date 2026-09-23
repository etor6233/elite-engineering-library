package httpapi

import (
	"bytes"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type PortalSessionModule struct {
	Profile identity.PortalProfile
	Store   identity.PortalSessionStore
}

func (m PortalSessionModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Profile.SHA256() == "" || m.Store == nil {
		return
	}
	mux.HandleFunc("POST /v1/private/portal-sessions/sweep", func(w http.ResponseWriter, r *http.Request) {
		principal, err := authenticate(r.Context(), r.Header.Get("Authorization"), verifier)
		if err != nil {
			writeProblem(w, 401, "UNAUTHENTICATED", "service identity required")
			return
		}
		if !m.Profile.ServiceAllowed(principal) || r.Header.Get("X-Portal-Profile-SHA256") != m.Profile.SHA256() {
			writeProblem(w, 403, "FORBIDDEN", "portal service profile required")
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "JSON required")
			return
		}
		var input struct {
			OperationID string `json:"operation_id"`
		}
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 256))
		dec.DisallowUnknownFields()
		if dec.Decode(&input) != nil || dec.Decode(new(any)) != io.EOF {
			writeProblem(w, 400, "INVALID_BODY", "bounded operation required")
			return
		}
		result, err := m.Store.Sweep(r.Context(), m.Profile, input.OperationID)
		if err != nil {
			writeProblem(w, 503, "SESSION_UNAVAILABLE", "session maintenance unavailable")
			return
		}
		writeJSON(w, 200, result)
	})
	for _, action := range []string{"create", "read", "claim", "commit", "abort", "revoke", "ack-revocation"} {
		mux.HandleFunc("POST /v1/private/portal-sessions/"+action, func(w http.ResponseWriter, r *http.Request) {
			principal, err := authenticate(r.Context(), r.Header.Get("Authorization"), verifier)
			if err != nil {
				writeProblem(w, 401, "UNAUTHENTICATED", "service identity required")
				return
			}
			if !m.Profile.ServiceAllowed(principal) || r.Header.Get("X-Portal-Profile-SHA256") != m.Profile.SHA256() {
				writeProblem(w, 403, "FORBIDDEN", "portal service profile required")
				return
			}
			if r.Header.Get("Content-Type") != "application/json" {
				writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "JSON required")
				return
			}
			raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 40000))
			if err != nil {
				writeProblem(w, 413, "BODY_TOO_LARGE", "bounded JSON required")
				return
			}
			first := json.NewDecoder(bytes.NewReader(raw))
			token, err := first.Token()
			if err != nil || token != json.Delim('{') {
				writeProblem(w, 400, "INVALID_BODY", "object required")
				return
			}
			seen := map[string]bool{}
			for first.More() {
				key, e := first.Token()
				name, ok := key.(string)
				if e != nil || !ok || seen[name] {
					writeProblem(w, 400, "INVALID_BODY", "duplicate field")
					return
				}
				seen[name] = true
				var value json.RawMessage
				if first.Decode(&value) != nil {
					writeProblem(w, 400, "INVALID_BODY", "invalid field")
					return
				}
			}
			var c identity.PortalSessionCommand
			dec := json.NewDecoder(strings.NewReader(string(raw)))
			dec.DisallowUnknownFields()
			if dec.Decode(&c) != nil || dec.Decode(new(any)) != io.EOF {
				writeProblem(w, 400, "INVALID_BODY", "invalid session command")
				return
			}
			record, err := m.Store.Execute(r.Context(), m.Profile, action, c)
			if errors.Is(err, identity.ErrPortalSessionNotFound) {
				writeProblem(w, 404, "SESSION_NOT_FOUND", "session unavailable")
				return
			}
			if errors.Is(err, identity.ErrPortalSessionConflict) {
				writeProblem(w, 409, "SESSION_CONFLICT", "session state changed")
				return
			}
			if err != nil {
				writeProblem(w, 503, "SESSION_UNAVAILABLE", "session store unavailable")
				return
			}
			writeJSON(w, 200, record)
		})
	}
}
