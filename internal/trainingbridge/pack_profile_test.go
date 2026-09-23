package trainingbridge

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestTrainingMaterializedProfile(t *testing.T) {
	raw, e := os.ReadFile("../../deploy/training/reference.profile.json")
	if e != nil {
		t.Fatal(e)
	}
	var doc ProfileDocument
	if e = json.Unmarshal(raw, &doc); e != nil {
		t.Fatal(e)
	}
	expected, b := profileFixture(t)
	if !reflect.DeepEqual(doc, expected) {
		t.Fatal("materialized profile differs from exercised curriculum")
	}
	path, _ := filepath.Abs("../../deploy/training/reference.profile.json")
	content, _ := filepath.Abs("../../training_content/help.bundle.json")
	p, e := LoadProfile(path, content, Activation{true, doc.ID, doc.Revision, digest(raw), doc.TenantID, doc.OrganizationID})
	if e != nil || len(p.Courses()) != 5 || digest(b) != doc.ContentSHA256 {
		t.Fatalf("materialized profile: %v", e)
	}
}
