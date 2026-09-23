package socialbridge

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func testProfile(t *testing.T) *Profile {
	t.Helper()
	raw, _ := json.Marshal(Config{Schema: "elite.meta-page-publishing.v1", TenantID: "tenant", OrganizationID: "organization", PageID: "123", Queue: "social", Publish: true, Revoke: true, Review: "one_distinct_human", LeaseSeconds: 60, RetrySeconds: 1, PollSeconds: 1, MaxAttempts: 2})
	p, e := Load(raw, Hash(raw))
	if e != nil {
		t.Fatal(e)
	}
	return p
}
func TestSocialProfileAdmissionAndChange(t *testing.T) {
	p := testProfile(t)
	c := p.Config()
	c.PollSeconds = 5
	c.Queue = "different_queue"
	raw, _ := json.Marshal(c)
	changed, e := Load(raw, Hash(raw))
	if e != nil || changed.SHA256() == p.SHA256() || changed.Config().PollSeconds != 5 {
		t.Fatal("profile configuration not admitted")
	}
	if _, e = Load(raw, p.SHA256()); e == nil {
		t.Fatal("hash drift")
	}
	c.Review = "automatic"
	raw, _ = json.Marshal(c)
	if _, e = Load(raw, Hash(raw)); e == nil {
		t.Fatal("unimplemented policy admitted")
	}
	for _, raw := range []string{`{"schema":"a","schema":"b"}`, `{"queue":"x","QUEUE":"y"}`, `{} {}`, `{"unknown":true}`, `null`} {
		if _, e = Load([]byte(raw), Hash([]byte(raw))); e == nil {
			t.Fatal("bad profile", raw)
		}
	}
}
func TestSocialProcessPreflight(t *testing.T) {
	python := os.Getenv("SOCIAL_TEST_PYTHON")
	if python == "" {
		t.Skip("locked Python runtime required")
	}
	script, e := filepath.Abs("../../meta_page_write/bridge.py")
	if e != nil {
		t.Fatal(e)
	}
	raw, e := os.ReadFile(script)
	if e != nil {
		t.Fatal(e)
	}
	p := Process{Python: python, Script: script, ScriptSHA256: Hash(raw), Environment: []string{"META_APP_ID=321", "META_APP_SECRET=fixture-secret", "META_PAGE_ACCESS_TOKEN=fixture-token"}}
	if e = p.Preflight(context.Background()); e != nil {
		t.Fatal("offline exact SDK preflight", e)
	}
	// A changed lock must not remove the installed-version checks silently.
	copyRoot := filepath.Join(t.TempDir(), "meta_page_write")
	if e = os.Mkdir(copyRoot, 0700); e != nil {
		t.Fatal(e)
	}
	for _, name := range []string{"bridge.py", "adapter.py", "requirements-windows-py314.lock"} {
		b, e := os.ReadFile(filepath.Join(filepath.Dir(script), name))
		if e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(filepath.Join(copyRoot, name), b, 0600); e != nil {
			t.Fatal(e)
		}
	}
	changed := p
	changed.Script = filepath.Join(copyRoot, "bridge.py")
	if e = os.WriteFile(filepath.Join(copyRoot, "requirements-windows-py314.lock"), []byte("# removed dependencies\n"), 0600); e != nil {
		t.Fatal(e)
	}
	if e = changed.Preflight(context.Background()); e == nil {
		t.Fatal("changed dependency lock disabled preflight checks")
	}
	p.Environment = nil
	if e = p.Preflight(context.Background()); e == nil {
		t.Fatal("missing credentials accepted")
	}
	p.ScriptSHA256 = Hash([]byte("different"))
	if e = p.Validate(); e == nil {
		t.Fatal("script drift accepted")
	}
}
func FuzzSocialProfileAndPayload(f *testing.F) {
	for _, s := range []string{`{"a":1}`, `{"a":{"b":[1,null,true]}}`, `{"a":1,"A":2}`, `null`, `{} {}`} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		var v any
		if Decode(raw, &v) == nil {
			canonical, e := json.Marshal(v)
			if e != nil {
				t.Fatal(e)
			}
			var again any
			if Decode(canonical, &again) != nil {
				t.Fatal("admitted JSON failed canonical decode")
			}
		}
		if p, e := Load(raw, Hash(raw)); e == nil && p.SHA256() != Hash(raw) {
			t.Fatal("unbound profile")
		}
	})
}
