package serialsupply

import (
	"encoding/json"
	"os"
	"testing"
)

func TestSupplyRoleCommandGoldens(t *testing.T) {
	raw, e := os.ReadFile("../../deploy/supply/role-goldens.json")
	if e != nil {
		t.Fatal(e)
	}
	var cases []struct {
		Name      string
		Payload   json.RawMessage
		Canonical string
		SHA256    string
	}
	if e = json.Unmarshal(raw, &cases); e != nil {
		t.Fatal(e)
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			var v any
			if c.Name == "create" {
				v = new(CreateRequest)
			} else {
				v = new(Command)
			}
			if e = json.Unmarshal(c.Payload, v); e != nil {
				t.Fatal(e)
			}
			switch r := v.(type) {
			case *CreateRequest:
				if !r.Valid() {
					t.Fatal("invalid create")
				}
			case *Command:
				if !r.Valid() {
					t.Fatal("invalid command")
				}
			}
			raw, hash, e := Canonical(v)
			if e != nil || string(raw) != c.Canonical || hash != c.SHA256 {
				t.Fatal(string(raw), hash, e)
			}
		})
	}
}
