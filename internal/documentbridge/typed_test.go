package documentbridge

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestTypedFixtureContracts(t *testing.T) {
	profile := filepath.Join(t.TempDir(), "profile.json")
	source := os.Getenv("DOCUMENT_TEST_ROOT")
	if source == "" {
		source = "../.."
	}
	b, e := os.ReadFile(filepath.Join(source, "config/documents/typed-profile.json"))
	if e != nil {
		t.Fatal(e)
	}
	if os.WriteFile(profile, b, 0600) != nil {
		t.Fatal("profile")
	}
	p, e := NewTypedFixturePipeline(profile, Hash(b))
	if e != nil {
		t.Fatal(e)
	}
	for _, c := range typedCatalog().Classes {
		t.Run(c.ID, func(t *testing.T) {
			o, e := TypedFixture(c.ID)
			if e != nil {
				t.Fatal(e)
			}
			ev, e := p.Run(context.Background(), o)
			if e != nil || ValidateTypedEvidence(ev, o, TypedProfileSHA()) != nil {
				t.Fatal(e)
			}
			var fields map[string]string
			if json.Unmarshal(ev.TypedFields, &fields) != nil {
				t.Fatal("fields")
			}
			delete(fields, c.Fields[0].Key)
			raw, _ := json.Marshal(fields)
			if ValidateTypedFields(c.ID, "1", raw) == nil {
				t.Fatal("missing field admitted")
			}
			fields = c.Suggested
			fields["untrusted_extra"] = "x"
			raw, _ = json.Marshal(fields)
			if ValidateTypedFields(c.ID, "1", raw) == nil {
				t.Fatal("extra field admitted")
			}
			o.Bytes = append(o.Bytes, 'x')
			o.SHA256 = Hash(o.Bytes)
			if _, e = p.Run(context.Background(), o); e == nil {
				t.Fatal("unknown bytes admitted")
			}
		})
	}
	o, _ := TypedFixture("receipt")
	o.ClassID = "supplier-invoice"
	if o.Validate("TYPED_FIXTURE") == nil {
		t.Fatal("cross class hash")
	}
	o.ClassID = "receipt"
	o.SchemaVersion = "2"
	if o.Validate("TYPED_FIXTURE") == nil {
		t.Fatal("unknown schema")
	}
	c, _ := ClassByID("proforma-invoice", "1")
	for _, bad := range []string{"2026-02-30", "0000-01-01", "2026-1-01"} {
		c.Suggested["valid_until"] = bad
		raw, _ := json.Marshal(c.Suggested)
		if ValidateTypedFields(c.ID, "1", raw) == nil {
			t.Fatal("bad date", bad)
		}
	}
	c, _ = ClassByID("supplier-invoice", "1")
	for _, bad := range []string{"1e10", "-1", "01", "NaN", " 1"} {
		c.Suggested["total"] = bad
		raw, _ := json.Marshal(c.Suggested)
		if ValidateTypedFields(c.ID, "1", raw) == nil {
			t.Fatal("bad decimal", bad)
		}
	}
	c, _ = ClassByID("receipt", "1")
	c.Suggested["merchant"] = "line\ncommand"
	raw, _ := json.Marshal(c.Suggested)
	if ValidateTypedFields(c.ID, "1", raw) == nil {
		t.Fatal("control char")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	o, _ = TypedFixture("receipt")
	if _, e = p.Run(ctx, o); e == nil {
		t.Fatal("canceled pipeline")
	}
}
