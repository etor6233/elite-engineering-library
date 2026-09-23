package trainingbridge

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func profileFixture(t *testing.T) (ProfileDocument, []byte) {
	t.Helper()
	b, e := os.ReadFile("../../training_content/help.bundle.json")
	if e != nil {
		t.Fatal(e)
	}
	var content ContentBundle
	if e = json.Unmarshal(b, &content); e != nil {
		t.Fatal(e)
	}
	return ProfileDocument{Schema: "elite-training-profile/v1", ID: "reference-onboarding", Revision: 2, TenantID: "50f38793-8a22-4f6b-983f-81dd0fca8208", OrganizationID: "store-1", Method: Method, ContentSHA256: digest(b), SourceSHA256: content.SourceSHA256, Courses: []Course{{ID: "resource-onboarding", Title: "Registrar recursos y recuperar una respuesta", Role: "employee", Lessons: []string{"resource-create-view"}, Prompts: []Prompt{{ID: "recovery", Text: "Describí qué harías si se pierde la respuesta después de registrar un recurso."}}}, {ID: "admin-onboarding", Title: "Revisión de recursos", Role: "admin", Lessons: []string{"resource-create-view", "supply-role-view", "warranty-role-view", "network-role-view"}, Prompts: []Prompt{{ID: "review", Text: "Describí cómo recuperarías una operación incierta, revisarías contenido y separarías una evaluación de los permisos de acceso."}}}, {ID: "owner-onboarding", Title: "Supervisión de la operación", Role: "owner", Lessons: []string{"operation-sections-view", "network-role-view", "training-role-view"}, Prompts: []Prompt{{ID: "review", Text: "Describí cómo distinguís una sección no disponible de una lista vacía."}}}, {ID: "customer-onboarding", Title: "Revisar una cotización", Role: "customer", Lessons: []string{"quote-acceptance-view"}, Prompts: []Prompt{{ID: "review", Text: "Describí qué confirma aceptar una cotización y qué se debe consultar ante una respuesta incierta."}}}, {ID: "content-onboarding", Title: "Contenido, publicación y revisión humana", Role: "admin", Lessons: []string{"help-cms-view", "catalog-role-view", "training-role-view"}, Prompts: []Prompt{{ID: "review", Text: "Explicá cómo publicarías contenido revisado, recuperarías una respuesta incierta y preservarías cursos anteriores sin conceder accesos."}}}}}, b
}
func activatedFixture(t *testing.T) (*Profile, []byte, []byte, Activation) {
	t.Helper()
	d, b := profileFixture(t)
	raw, _ := json.Marshal(d)
	a := Activation{true, d.ID, d.Revision, digest(raw), d.TenantID, d.OrganizationID}
	p, e := ParseProfile(raw, b, a)
	if e != nil {
		t.Fatal(e)
	}
	return p, raw, b, a
}
func TestTrainingProfileFixedContentAndNoGrants(t *testing.T) {
	p, raw, b, a := activatedFixture(t)
	if len(p.Courses()) != 5 {
		t.Fatal("role views")
	}
	for name, mutate := range map[string]func([]byte, []byte, Activation) ([]byte, []byte, Activation){
		"disabled":     func(x, y []byte, a Activation) ([]byte, []byte, Activation) { a.Enabled = false; return x, y, a },
		"wrong-tenant": func(x, y []byte, a Activation) ([]byte, []byte, Activation) { a.TenantID = "other"; return x, y, a },
		"wrong-org": func(x, y []byte, a Activation) ([]byte, []byte, Activation) {
			a.OrganizationID = "other"
			return x, y, a
		},
		"wrong-revision":  func(x, y []byte, a Activation) ([]byte, []byte, Activation) { a.Revision++; return x, y, a },
		"content-altered": func(x, y []byte, a Activation) ([]byte, []byte, Activation) { return x, append(y, ' '), a },
		"policy-incompatible": func(x, y []byte, a Activation) ([]byte, []byte, Activation) {
			x = []byte(strings.Replace(string(x), Method, "AUTO_SCORE_AND_GRANT", 1))
			a.SHA256 = digest(x)
			return x, y, a
		},
		"duplicate-key": func(x, y []byte, a Activation) ([]byte, []byte, Activation) {
			x = []byte(strings.Replace(string(x), `"revision":2`, `"revision":2,"revision":2`, 1))
			a.SHA256 = digest(x)
			return x, y, a
		},
		"unknown-field": func(x, y []byte, a Activation) ([]byte, []byte, Activation) {
			x = append([]byte(`{"auto_approve":true,`), x[1:]...)
			a.SHA256 = digest(x)
			return x, y, a
		},
	} {
		t.Run(name, func(t *testing.T) {
			x, y, c := mutate(append([]byte(nil), raw...), append([]byte(nil), b...), a)
			if _, e := ParseProfile(x, y, c); e == nil {
				t.Fatal("unsafe activation accepted")
			}
		})
	}
}
