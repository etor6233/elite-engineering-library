package main

import (
	"context"
	"crypto/sha256"
	tr "elite.local/enterprise/internal/trainingbridge"
	"encoding/hex"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"path/filepath"
	"testing"
)

func TestTrainingHostProfileAndGuards(t *testing.T) {
	ctx := context.Background()
	calls := 0
	off := func(k string) string {
		calls++
		if k != "TRAINING_ENABLED" {
			t.Fatal("disabled host read", k)
		}
		return "false"
	}
	if m, e := selectedTrainingModule(ctx, nil, off); e != nil || m != nil || calls != 1 {
		t.Fatal("disabled", e)
	}
	db := os.Getenv("ELITE_TRAINING_DATABASE_URL")
	if db == "" {
		t.Skip("owned fixture DB required")
	}
	pool, e := pgxpool.New(ctx, db)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	b, e := os.ReadFile("../../training_content/help.bundle.json")
	if e != nil {
		t.Fatal(e)
	}
	var bundle tr.ContentBundle
	if e = json.Unmarshal(b, &bundle); e != nil {
		t.Fatal(e)
	}
	sum := sha256.Sum256(b)
	d := tr.ProfileDocument{Schema: "elite-training-profile/v1", ID: "host-reference", Revision: 1, TenantID: "50f38793-8a22-4f6b-983f-81dd0fca8208", OrganizationID: "store-1", Method: tr.Method, ContentSHA256: hex.EncodeToString(sum[:]), SourceSHA256: bundle.SourceSHA256, Courses: []tr.Course{{ID: "resource-onboarding", Title: "Recuperar un registro", Role: "employee", Lessons: []string{"resource-create-view"}, Prompts: []tr.Prompt{{ID: "recovery", Text: "Explicá cómo consultarías el registro."}}}}}
	raw, _ := json.Marshal(d)
	sum = sha256.Sum256(raw)
	dir := t.TempDir()
	profile, content := filepath.Join(dir, "profile.json"), filepath.Join(dir, "content.json")
	if e = os.WriteFile(profile, raw, 0600); e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(content, b, 0600); e != nil {
		t.Fatal(e)
	}
	values := map[string]string{"TRAINING_ENABLED": "true", "TRAINING_PROFILE_FILE": profile, "TRAINING_CONTENT_FILE": content, "TRAINING_PROFILE_ID": d.ID, "TRAINING_PROFILE_REVISION": "1", "TRAINING_PROFILE_SHA256": hex.EncodeToString(sum[:]), "TRAINING_TENANT_ID": d.TenantID, "TRAINING_ORGANIZATION_ID": d.OrganizationID}
	lookup := func(k string) string { return values[k] }
	if m, e := selectedTrainingModule(ctx, pool, lookup); e != nil || m == nil {
		t.Fatal("valid activation", e)
	}
	values["TRAINING_PROFILE_FILE"] = "relative.json"
	if _, e = selectedTrainingModule(ctx, pool, lookup); e == nil {
		t.Fatal("relative profile")
	}
	values["TRAINING_PROFILE_FILE"] = dir
	if _, e = selectedTrainingModule(ctx, pool, lookup); e == nil {
		t.Fatal("directory profile")
	}
	values["TRAINING_PROFILE_FILE"] = profile
	values["TRAINING_PROFILE_SHA256"] = ""
	if _, e = selectedTrainingModule(ctx, pool, lookup); e == nil {
		t.Fatal("unbound profile")
	}
	values["TRAINING_PROFILE_SHA256"] = hex.EncodeToString(sum[:])
	if _, e = pool.Exec(ctx, `alter table approval.request disable trigger training_assessment_guard`); e != nil {
		t.Fatal(e)
	}
	defer func() {
		if _, e := pool.Exec(context.Background(), `alter table approval.request enable trigger training_assessment_guard`); e != nil {
			t.Error(e)
		}
	}()
	if _, e = selectedTrainingModule(ctx, pool, lookup); e == nil {
		t.Fatal("disabled guard")
	}
	t.Log("TRAINING_HOST_PASS opt-in no-read; exact profile/content scope; regular absolute paths; shared participation/review guards and indexes required")
}
