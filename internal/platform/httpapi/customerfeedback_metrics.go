package httpapi

// AUTHORED reporting wire. Existing Summary retains threshold/retention and calls the admitted PostHog NPS owner.
import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var metricSurveyID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

func (a feedbackAPI) registerMetrics(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/reporting/surveys/{survey}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "private, no-store")
		w.Header().Set("Vary", "Authorization")
		s, ok := a.scope(w, r, "surveys:read")
		if !ok {
			return
		}
		q, e := url.ParseQuery(r.URL.RawQuery)
		if e != nil || len(r.URL.RawQuery) > 256 || len(q) != 1 || !metricSurveyID.MatchString(s.Organization) || !metricSurveyID.MatchString(r.PathValue("survey")) {
			writeProblem(w, 400, "METRIC_INVALID", "bounded scope and survey required")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		v, e := a.service.Summary(ctx, s, r.PathValue("survey"))
		if e != nil {
			feedbackError(w, e)
			return
		}
		writeJSON(w, 200, map[string]any{"organization_id": s.Organization, "survey_id": r.PathValue("survey"), "source": "crm.survey_definition + crm.survey_response", "basis": "retained_survey_population_with_configured_minimum", "responses": strconv.FormatInt(v.Responses, 10), "available": v.Available, "nps": v.NPS, "observed_at": time.Now().UTC()})
	})
}
