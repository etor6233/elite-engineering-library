package trainingbridge

import (
	"encoding/json"
	"os"
	"testing"
)

func FuzzTrainingProfileIsolation(f *testing.F) {
	content, e := os.ReadFile("../../training_content/help.bundle.json")
	if e != nil {
		f.Fatal(e)
	}
	var bundle ContentBundle
	if e = json.Unmarshal(content, &bundle); e != nil {
		f.Fatal(e)
	}
	d := ProfileDocument{Schema: "elite-training-profile/v1", ID: "reference-onboarding", Revision: 1, TenantID: "50f38793-8a22-4f6b-983f-81dd0fca8208", OrganizationID: "store-1", Method: Method, ContentSHA256: digest(content), SourceSHA256: bundle.SourceSHA256, Courses: []Course{{ID: "resource-onboarding", Title: "Revisar referencia", Role: "employee", Lessons: []string{"resource-create-view"}, Prompts: []Prompt{{ID: "recovery", Text: "Explicá cómo consultarías el registro."}}}}}
	raw, _ := json.Marshal(d)
	f.Add(raw)
	f.Add([]byte("{}"))
	f.Add([]byte(`{"revision":1,"revision":2}`))
	f.Fuzz(func(t *testing.T, candidate []byte) {
		if len(candidate) > 65536 {
			return
		}
		p, e := ParseProfile(candidate, content, Activation{true, d.ID, d.Revision, digest(candidate), d.TenantID, d.OrganizationID})
		if e != nil {
			return
		}
		views := p.Courses()
		tenant, org := p.Scope()
		if p.Hash() != digest(candidate) || tenant != d.TenantID || org != d.OrganizationID || len(views) < 1 || len(views) > 8 {
			t.Fatal("accepted profile escaped activation")
		}
		for _, v := range views {
			if v.Method != Method || v.ProfileSHA256 != digest(candidate) || v.ContentSHA256 != digest(content) || len(v.Articles) < 1 || len(v.Articles) > 4 {
				t.Fatal("view binding changed")
			}
			before, _ := json.Marshal(v)
			v.Articles[0].Paragraphs[0] = "caller mutation"
			again, e := p.Course(v.Course.ID)
			if e != nil {
				t.Fatal(e)
			}
			actual, _ := json.Marshal(again)
			if string(before) != string(actual) {
				t.Fatal("consumer mutated frozen content")
			}
		}
	})
}
