package helpcms

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestHelpCMSWireGoldens(t *testing.T) {
	raw, e := os.ReadFile("../../deploy/help/cms-goldens.json")
	if e != nil {
		t.Fatal(e)
	}
	var rows []struct {
		Command   Command `json:"command"`
		Canonical string  `json:"canonical"`
		SHA       string  `json:"sha256"`
	}
	if json.Unmarshal(raw, &rows) != nil || len(rows) != 4 {
		t.Fatal("four typed goldens required")
	}
	for _, r := range rows {
		body, hash, e := Canonical(r.Command)
		if !r.Command.Valid() || e != nil || string(body) != r.Canonical || hash != r.SHA {
			t.Fatal(r.Command.Action, e)
		}
	}
	for _, s := range []string{"0", "01", "-1", "9223372036854775807", "9223372036854775808"} {
		c := rows[1].Command
		c.Version = s
		if c.Valid() {
			t.Fatal("unsafe version", s)
		}
	}
	for _, body := range []string{strings.Repeat("a", 16385), strings.Repeat("<", 5000), "bad\x00body", string([]byte{0xff})} {
		c := rows[0].Command
		c.Body = body
		if c.Valid() {
			t.Fatal("unbounded/invalid body")
		}
	}
}
func FuzzHelpCMSWire(f *testing.F) {
	for _, s := range []string{`{"command_id":"create","action":"create","article_id":"guide","organization_id":"org","locale":"es","category":"operations","title":"Guía","body":"Contenido"}`, `{"command_id":"publish","action":"publish","article_id":"guide","organization_id":"org","version":"9007199254740993"}`, `{}`} {
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
			t.Fatal("valid command cannot canonicalize", e)
		}
		var again Command
		if json.Unmarshal(body, &again) != nil || again != c || !again.Valid() {
			t.Fatal("canonical round trip")
		}
	})
}
