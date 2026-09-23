package storedvaluebridge

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Process verifies the materialized source manifest before each bounded pure
// calculation. No credentials, inherited Python path, HTTP or DB are exposed.
type Process struct{ Python, Script, ScriptSHA256, Manifest, ManifestSHA256 string }

func (p Process) Validate() error {
	if !filepath.IsAbs(p.Python) || !filepath.IsAbs(p.Script) || !filepath.IsAbs(p.Manifest) || !hashRE.MatchString(p.ScriptSHA256) || !hashRE.MatchString(p.ManifestSHA256) {
		return ErrBinding
	}
	if info, e := os.Stat(p.Python); e != nil || info.IsDir() {
		return ErrBinding
	}
	for path, expected := range map[string]string{p.Script: p.ScriptSHA256, p.Manifest: p.ManifestSHA256} {
		b, e := os.ReadFile(path)
		if e != nil || len(b) > 65536 || Hash(b) != expected {
			return ErrBinding
		}
	}
	return nil
}

type output struct {
	bytes.Buffer
	exceeded bool
}

func (o *output) Write(b []byte) (int, error) {
	if o.Len()+len(b) > 32768 {
		o.exceeded = true
		return 0, io.ErrShortBuffer
	}
	return o.Buffer.Write(b)
}
func (p Process) Execute(ctx context.Context, request Calculation) (Result, error) {
	var result Result
	if p.Validate() != nil {
		return result, ErrBinding
	}
	input, expected, e := Canonical(request)
	if e != nil {
		return result, e
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, p.Python, "-I", "-S", "-B", p.Script, p.ManifestSHA256)
	c.Env = []string{}
	for _, key := range []string{"SYSTEMROOT", "WINDIR", "TEMP", "TMP"} {
		if v := os.Getenv(key); v != "" {
			c.Env = append(c.Env, key+"="+v)
		}
	}
	c.Stdin = bytes.NewReader(input)
	var out output
	c.Stdout = &out
	c.Stderr = io.Discard
	if c.Run() != nil || out.exceeded || Decode(out.Bytes(), &result) != nil {
		return result, ErrBinding
	}
	if result.Schema != "elite.odoo-loyalty-result.v1" || result.SourceRevision != "99edb6dd82b7b560930c00b03b694ba700785370" || result.ProgramSHA256 != request.ProgramSHA256 || result.RequestSHA256 != expected {
		return Result{}, ErrBinding
	}
	return result, nil
}
func (p Process) Preflight(ctx context.Context, profile *Profile) error {
	if p.Validate() != nil || !profile.Valid() {
		return ErrBinding
	}
	for _, v := range profile.config.Programs {
		o := Order{OrderID: "preflight", OrganizationID: profile.config.OrganizationID, SubjectID: "preflight", State: "draft", Currency: v.Calculation.Currency, Total: "0", EnabledRuleIDs: []string{}, Lines: []Line{{ID: "preflight", ProductID: "preflight", Quantity: "1", Subtotal: "0", Tax: "0", Total: "0"}}}
		r, e := NewCalculation("evaluate", v.Calculation, o, map[string]any{})
		if e != nil {
			return e
		}
		if _, e = p.Execute(ctx, r); e != nil {
			return e
		}
	}
	return nil
}
func ReadPoints(result Result) ([]string, error) {
	var v struct {
		Points []string `json:"points"`
		Error  string   `json:"error"`
	}
	if json.Unmarshal(result.Result, &v) != nil || v.Error != "" || v.Points == nil {
		return nil, ErrBinding
	}
	return v.Points, nil
}
