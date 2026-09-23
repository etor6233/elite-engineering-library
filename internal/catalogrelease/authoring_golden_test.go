package catalogrelease

import (
	"encoding/json"
	"os"
	"testing"
)

func TestCatalogRoleCommandGolden(t *testing.T) {
	raw, e := os.ReadFile("../../deploy/catalog/authoring-goldens.json")
	if e != nil {
		t.Fatal(e)
	}
	var cases []struct {
		Name      string          `json:"name"`
		Payload   json.RawMessage `json:"payload"`
		Canonical string          `json:"canonical"`
		SHA256    string          `json:"sha256"`
	}
	if e = json.Unmarshal(raw, &cases); e != nil {
		t.Fatal(e)
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			var value any
			switch c.Name {
			case "source":
				v := new(SourceRequest)
				e = json.Unmarshal(c.Payload, v)
				value = v
			case "review":
				v := new(ReviewRequest)
				e = json.Unmarshal(c.Payload, v)
				value = v
			case "publish":
				v := new(PublishRequest)
				e = json.Unmarshal(c.Payload, v)
				value = v
			case "media":
				v := map[string]string{}
				e = json.Unmarshal(c.Payload, &v)
				value = v
			default:
				t.Fatal("unknown golden")
			}
			if e != nil {
				t.Fatal(e)
			}
			got, h, e := Canonical(value)
			if e != nil || string(got) != c.Canonical || h != c.SHA256 {
				t.Fatalf("typed Go command differs: %v %s %s", e, h, got)
			}
		})
	}
}
