package storedvaluebridge

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func fixture(t *testing.T) (*Profile, Process) {
	t.Helper()
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	raw, e := os.ReadFile(filepath.Join(root, "deploy", "stored-value", "profile.reference.json"))
	if e != nil {
		t.Fatal(e)
	}
	p, e := Load(raw, Hash(raw))
	if e != nil {
		t.Fatal(e)
	}
	script := filepath.Join(root, "odoo_loyalty", "run.py")
	manifest := filepath.Join(root, "odoo_loyalty", "engine-lock.json")
	scriptBytes, _ := os.ReadFile(script)
	manifestBytes, _ := os.ReadFile(manifest)
	return p, Process{Python: `C:\Python314\python.exe`, Script: script, ScriptSHA256: Hash(scriptBytes), Manifest: manifest, ManifestSHA256: Hash(manifestBytes)}
}
func TestExactAmountsAndProfileContracts(t *testing.T) {
	for _, n := range []int64{0, 1, 100, -25, 9223372036854775807, -9223372036854775808} {
		v, e := Minor(Major(n, 2), 2)
		if e != nil || v != n {
			t.Fatalf("roundtrip %d: %d %v", n, v, e)
		}
	}
	for _, v := range []string{"0.001", "NaN", "1e2", "92233720368547758.08"} {
		if _, e := Minor(v, 2); e == nil {
			t.Fatalf("accepted %s", v)
		}
	}
	p, _ := fixture(t)
	raw := p.raw
	for _, bad := range [][]byte{[]byte(strings.Replace(string(raw), `"currency_digits": 2`, `"currency_digits": null`, 1)), []byte(strings.Replace(string(raw), `"points": "1"`, `"points": "0"`, 1)), []byte(strings.Replace(string(raw), `"schema"`, `"Schema"`, 1)), []byte(strings.Replace(string(raw), `"tax_mode": "incl"`, `"tax_mode": "excl"`, 1))} {
		if _, e := Load(bad, Hash(bad)); e == nil {
			t.Fatal("accepted ambiguous or unsupported profile")
		}
	}
	if _, e := Load(raw, strings.Repeat("0", 64)); e == nil {
		t.Fatal("profile hash ignored")
	}
}
func TestActualIsolatedProcessAndSourceBinding(t *testing.T) {
	p, process := fixture(t)
	if _, e := os.Stat(process.Python); e != nil {
		t.Fatal("qualified Python runtime missing")
	}
	if e := process.Preflight(context.Background(), p); e != nil {
		t.Fatal(e)
	}
	entry, _ := p.Program("reference-gift_card")
	order := Order{OrderID: "order-1", OrganizationID: "franchise-1", SubjectID: "customer-1", State: "draft", Currency: "ARS", Total: "100", EnabledRuleIDs: []string{}, Lines: []Line{{ID: "line-1", ProductID: "gift-50", Quantity: "2", Subtotal: "100", Tax: "0", Total: "100"}}}
	request, e := NewCalculation("evaluate", entry.Calculation, order, map[string]any{})
	if e != nil {
		t.Fatal(e)
	}
	result, e := process.Execute(context.Background(), request)
	if e != nil {
		t.Fatal(e)
	}
	points, e := ReadPoints(result)
	if e != nil || strings.Join(points, ",") != "50.00,50.00" {
		t.Fatalf("%s %v", result.Result, e)
	}
	bad := process
	bad.ManifestSHA256 = strings.Repeat("0", 64)
	if _, e = bad.Execute(context.Background(), request); e == nil {
		t.Fatal("changed manifest binding accepted")
	}
	// Same manifest bytes with a changed module in an isolated copied payload must fail before calculation.
	temp := t.TempDir()
	base := filepath.Dir(process.Script)
	e = filepath.WalkDir(base, func(path string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(base, path)
		dest := filepath.Join(temp, rel)
		if e = os.MkdirAll(filepath.Dir(dest), 0700); e != nil {
			return e
		}
		body, e := os.ReadFile(path)
		if e != nil {
			return e
		}
		return os.WriteFile(dest, body, 0600)
	})
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(temp, "engine.py"), []byte("raise RuntimeError('tampered')\n"), 0600); e != nil {
		t.Fatal(e)
	}
	bad = process
	bad.Script = filepath.Join(temp, "run.py")
	bad.Manifest = filepath.Join(temp, "engine-lock.json")
	if _, e = bad.Execute(context.Background(), request); e == nil {
		t.Fatal("changed source bytes accepted")
	}
}
func TestSourceProfileChangeChangesCalculatedResult(t *testing.T) {
	p, process := fixture(t)
	entry, _ := p.Program("reference-loyalty")
	o := Order{OrderID: "order-1", OrganizationID: "franchise-1", SubjectID: "customer-1", State: "draft", Currency: "ARS", Total: "100", EnabledRuleIDs: []string{}, Lines: []Line{{ID: "line-1", ProductID: "sku", Quantity: "2", Subtotal: "100", Tax: "0", Total: "100"}}}
	first, _ := NewCalculation("evaluate", entry.Calculation, o, map[string]any{})
	a, e := process.Execute(context.Background(), first)
	if e != nil {
		t.Fatal(e)
	}
	entry.Calculation.Rules[0].Points = "2"
	unchanged, _ := p.Program("reference-loyalty")
	if unchanged.Calculation.Rules[0].Points != "1" {
		t.Fatal("mutated immutable profile")
	}
	second, _ := NewCalculation("evaluate", entry.Calculation, o, map[string]any{})
	b, e := process.Execute(context.Background(), second)
	if e != nil {
		t.Fatal(e)
	}
	if string(a.Result) == string(b.Result) || a.ProgramSHA256 == b.ProgramSHA256 {
		t.Fatal("configuration change not reflected")
	}
	var result struct{ Points []string }
	if json.Unmarshal(b.Result, &result) != nil || len(result.Points) != 1 || result.Points[0] != "200.00" {
		t.Fatalf("%s", b.Result)
	}
}

func TestPointsRepresentationAndReceiptIntegerBinding(t *testing.T) {
	for _, v := range []string{"50.00", "-0.250000", "200.000001"} {
		if _, e := Minor(v, 6); e != nil {
			t.Fatal(v, e)
		}
	}
	var c Calculation
	raw := []byte(`{"schema":"x","operation":"reverse","program":{"id":"","kind":"","currency":"","currency_digits":0,"applies_on":"","nominative":false,"trigger":"","trigger_product_ids":null,"rules":null},"program_sha256":"","order":{"order_id":"","organization_id":"","subject_id":"","state":"","public_subject":false,"currency":"","total":"","enabled_rule_ids":null,"lines":null},"data":{"applied_minor_units":9007199254740993}}`)
	if Decode(raw, &c) != nil || !strings.Contains(string(c.Data), "9007199254740993") {
		t.Fatal("approval data lost integer precision")
	}
}
