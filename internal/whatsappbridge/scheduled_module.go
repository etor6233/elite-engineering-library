package whatsappbridge

// AUTHORED bounded operator interface, mounted on the existing role API.
import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

func decodeSchedule(w http.ResponseWriter, r *http.Request, v any) error {
	if r.Header.Get("Content-Type") != "application/json" || r.URL.RawQuery != "" {
		return ErrApproval
	}
	if http.NewResponseController(w).SetReadDeadline(time.Now().Add(5*time.Second)) != nil {
		return ErrApproval
	}
	raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
	if e != nil {
		return ErrApproval
	}
	if _, _, e = approval.CanonicalPayload(raw); e != nil {
		return ErrApproval
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(v) != nil || d.Decode(new(any)) != io.EOF {
		return ErrApproval
	}
	return nil
}
func (m *ScheduledNotifications) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m == nil {
		return
	}
	route := func(method, path, permission string, action func(http.ResponseWriter, *http.Request, identity.Principal)) {
		mux.HandleFunc(method+" /v1/franchise/notifications/scheduled"+path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			header := r.Header.Values("Authorization")
			if verifier == nil || len(header) != 1 || !strings.HasPrefix(header[0], "Bearer ") || len(header[0]) > 16391 || len(strings.Fields(header[0])) != 2 {
				notificationProblem(w, 401, "UNAUTHENTICATED")
				return
			}
			p, e := verifier.Verify(r.Context(), strings.TrimPrefix(header[0], "Bearer "))
			if e != nil {
				notificationProblem(w, 401, "UNAUTHENTICATED")
				return
			}
			if !m.approvals.authorized(p, permission) {
				notificationProblem(w, 403, "FORBIDDEN")
				return
			}
			if r.URL.RawQuery != "" || r.PathValue("key") != "" && !validScheduleKey(r.PathValue("key")) {
				notificationProblem(w, 400, "INVALID_TARGET")
				return
			}
			ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
			defer cancel()
			action(w, r.WithContext(ctx), p)
		})
	}
	fail := func(w http.ResponseWriter, e error) bool {
		if e == nil {
			return false
		}
		notificationProblem(w, 409, "SCHEDULE_NOT_CURRENT_OR_DIVERGENT")
		return true
	}
	route("POST", "/prepare", "notification:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in ScheduleRequest
		if e := decodeSchedule(w, r, &in); e != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		out, e := m.approvals.Prepare(r.Context(), p, in)
		if !fail(w, e) {
			replyJSON(w, out)
		}
	})
	route("POST", "/requests", "notification:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in ScheduleContext
		if e := decodeSchedule(w, r, &in); e != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		h, replay, e := m.approvals.Submit(r.Context(), p, in)
		if !fail(w, e) {
			replyJSON(w, map[string]any{"delivery_key": in.DeliveryKey, "request_sha256": h, "replay": replay})
		}
	})
	route("GET", "/requests/{key}", "notification:read", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		out, e := m.approvals.Status(r.Context(), p, r.PathValue("key"))
		if !fail(w, e) {
			replyJSON(w, out)
		}
	})
	route("POST", "/requests/{key}/decision", "notification:approve", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in struct {
			Hash    string `json:"request_sha256"`
			Approve *bool  `json:"approve"`
			Reason  string `json:"reason"`
		}
		if decodeSchedule(w, r, &in) != nil || in.Approve == nil || !validDigest(in.Hash) {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		out, e := m.approvals.Decide(r.Context(), p, r.PathValue("key"), in.Hash, *in.Approve, in.Reason)
		if !fail(w, e) {
			replyJSON(w, map[string]any{"state": out})
		}
	})
	route("POST", "/requests/{key}/cancel", "notification:cancel", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in struct {
			Hash   string `json:"request_sha256"`
			Reason string `json:"reason"`
		}
		if decodeSchedule(w, r, &in) != nil || !validDigest(in.Hash) {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		if !fail(w, m.approvals.Cancel(r.Context(), p, r.PathValue("key"), in.Hash, in.Reason)) {
			replyJSON(w, map[string]bool{"cancelled": true})
		}
	})
	route("POST", "/requests/{key}/reconcile", "notification:reconcile", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in struct {
			Hash string `json:"request_sha256"`
		}
		if decodeSchedule(w, r, &in) != nil || !validDigest(in.Hash) {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		out, e := m.Recover(r.Context(), p, r.PathValue("key"), in.Hash)
		if !fail(w, e) {
			replyJSON(w, out)
		}
	})
}

func (m *ScheduledNotifications) Run(ctx context.Context, source StatusPrincipalSource, reporter StatusReporter, interval time.Duration) error {
	if m == nil || source == nil || reporter == nil || interval < time.Second || interval > time.Minute {
		return ErrStatusHost
	}
	process := func(ctx context.Context, p identity.Principal) (StatusWorkResult, error) {
		v, e := m.ProcessOnce(ctx, p)
		return StatusWorkResult{Claimed: v.Claimed, Completed: v.Outcome != "" && v.Outcome != "LEASE_EXHAUSTED", FailureRecorded: v.Outcome == "RECONCILIATION_REQUIRED" || v.Outcome == "FAILED_TERMINAL" || v.Outcome == "LEASE_EXHAUSTED", Terminal: v.Outcome == "FAILED_TERMINAL" || v.Outcome == "LEASE_EXHAUSTED"}, e
	}
	return runStatusLoop(ctx, process, source, reporter, interval, waitStatusPoll)
}
func (s *ScheduleApprovals) VerifyInfrastructure(ctx context.Context) error {
	var valid bool
	e := s.base.pool.QueryRow(ctx, `select
 (select count(*) from pg_trigger where not tgisinternal and tgenabled in('O','A') and
 (tgrelid=to_regclass('communication.whatsapp_schedule') and tgname='whatsapp_schedule_immutable'
 or tgrelid=to_regclass('communication.whatsapp_schedule_cancellation') and tgname='whatsapp_schedule_cancellation_immutable'
 or tgrelid=to_regclass('communication.whatsapp_schedule_result') and tgname='whatsapp_schedule_result_immutable'))=3
 and pg_get_functiondef(to_regprocedure('approval.guard_bound_request()')) like '%whatsapp_schedule%'
 and pg_get_functiondef(to_regprocedure('approval.guard_bound_decision()')) like '%whatsapp_schedule%'
 and pg_get_viewdef(to_regclass('communication.whatsapp_delivery_approval')) like '%whatsapp_schedule%'`).Scan(&valid)
	if e != nil || !valid {
		return ErrApproval
	}
	return nil
}
