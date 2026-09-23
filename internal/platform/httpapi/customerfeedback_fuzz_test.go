package httpapi

import (
	"encoding/json"
	"strings"
	"testing"
)

func FuzzFeedbackSubmission(f *testing.F) {
	for _, seed := range []string{
		`{"score":0,"consent":true,"consent_version":"v1"}`,
		`{"score":10,"consent":true,"consent_version":"v1"}`,
		`{"score":-1,"consent":true,"consent_version":"v1"}`,
		`{"score":0,"score":10,"consent":true,"consent_version":"v1"}`,
		`{"score":null,"consent":true,"consent_version":"v1"}`,
		`{"score":0,"consent":false,"consent_version":"v1"}`,
		`{}`, `[]`, `null`, `{"score":1e0}`, `{"scor\u0065":0,"score":1}`,
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		if len(raw) > 2048 {
			return
		}
		got, err := readSurveySubmission(strings.NewReader(raw))
		if err != nil {
			return
		}
		var wire struct {
			Score   int    `json:"score"`
			Consent bool   `json:"consent"`
			Version string `json:"consent_version"`
		}
		if e := json.Unmarshal([]byte(raw), &wire); e != nil {
			t.Fatalf("accepted non-JSON: %v", e)
		}
		if got.Score != wire.Score || got.Consent != wire.Consent || got.ConsentVersion != wire.Version {
			t.Fatal("wire and strict parser disagree")
		}
	})
}
