package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFeedbackUnconfiguredModuleHasNoRoutes(t *testing.T) {
	mux := http.NewServeMux()
	CustomerFeedbackModule{}.Register(mux, nil)
	for _, path := range []string{"/v1/customer/surveys/survey", "/v1/customer/surveys/survey/response", "/v1/admin/surveys/survey/summary"} {
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, httptest.NewRequest("GET", path, nil))
		if response.Code != 404 {
			t.Fatal(path, response.Code)
		}
	}
}

func TestFeedbackStrictSubmission(t *testing.T) {
	good := `{"score":0,"consent":true,"consent_version":"v1"}`
	got, e := readSurveySubmission(strings.NewReader(good))
	if e != nil || got.Score != 0 || !got.Consent || got.ConsentVersion != "v1" {
		t.Fatal(got, e)
	}
	for _, body := range []string{
		`{"consent":true,"consent_version":"v1"}`, `{"score":null,"consent":true,"consent_version":"v1"}`,
		`{"score":1.5,"consent":true,"consent_version":"v1"}`, `{"score":1,"score":2,"consent":true,"consent_version":"v1"}`,
		`{"score":1,"\u0073core":2,"consent":true,"consent_version":"v1"}`,
		good + `{}`, good + "false", `[]`, `null`, `{"score":1,"consent":null,"consent_version":"v1"}`,
		`{"score":1,"consent":true,"consent_version":null}`, `{"score":1,"consent":true,"consent_version":"v1","customer":"another"}`,
	} {
		t.Run(body, func(t *testing.T) {
			if _, e := readSurveySubmission(strings.NewReader(body)); e == nil {
				t.Fatal("accepted ambiguous body")
			}
		})
	}
}
