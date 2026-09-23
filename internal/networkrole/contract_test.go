package networkrole

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestNetworkRoleWireGoldens(t *testing.T) {
	raw, e := os.ReadFile("../../deploy/network/role-goldens.json")
	if e != nil {
		t.Fatal(e)
	}
	var rows []struct {
		Command   Command `json:"command"`
		Canonical string  `json:"canonical"`
		SHA       string  `json:"sha256"`
	}
	if json.Unmarshal(raw, &rows) != nil || len(rows) != 4 {
		t.Fatal("four shared goldens required")
	}
	for _, r := range rows {
		if !r.Command.Valid() {
			t.Fatal("invalid golden", r.Command)
		}
		body, hash, e := Canonical(r.Command)
		if e != nil || string(body) != r.Canonical || hash != r.SHA {
			t.Fatal(r.Command.Action, string(body), e)
		}
	}
	for _, s := range []string{"0", "01", "-1", "9223372036854775807", "9223372036854775808"} {
		c := rows[2].Command
		c.Version = s
		if c.Valid() {
			t.Fatal("unsafe version", s)
		}
	}
	for _, s := range []string{"UPPER", "space code", strings.Repeat("x", 129)} {
		c := rows[0].Command
		c.Code = s
		if c.Valid() {
			t.Fatal("schema code", s)
		}
	}
}
func FuzzNetworkRoleWire(f *testing.F) {
	for _, s := range []string{`{"command_id":"x","action":"transition-organization","scope_organization_id":"org","entity_id":"org","current":"active","target":"closed","version":"1"}`, `{"action":"create-organization"}`, `{}`} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 32768 {
			return
		}
		var c Command
		if json.Unmarshal(raw, &c) != nil || !c.Valid() {
			return
		}
		body, hash, e := Canonical(c)
		if e != nil || len(hash) != 64 {
			t.Fatal("valid command is not canonical", e)
		}
		var again Command
		if json.Unmarshal(body, &again) != nil || again != c || !again.Valid() {
			t.Fatal("canonical round trip")
		}
	})
}
