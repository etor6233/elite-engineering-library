package capture

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func captureConfig(t *testing.T) Config {
	t.Helper()
	path := os.Getenv("CAPTURE_TEST_PROFILE")
	if path == "" {
		t.Skip("pinned capture runtime required")
	}
	b, e := os.ReadFile(path)
	if e != nil {
		t.Fatal(e)
	}
	var c Config
	if json.Unmarshal(b, &c) != nil {
		t.Fatal("config")
	}
	return c
}
func TestCaptureConfigAndBoundedResult(t *testing.T) {
	c := captureConfig(t)
	if c.Verify() != nil {
		t.Fatal("verified runtime refused")
	}
	wrong := c
	wrong.OS = "linux"
	if wrong.Verify() == nil {
		t.Fatal("OS mismatch")
	}
	wrong = c
	wrong.Arch = "arm64"
	if wrong.Verify() == nil {
		t.Fatal("arch mismatch")
	}
	wrong = c
	wrong.WheelTag = "cp312-abi3-manylinux_2_27_x86_64"
	if wrong.Verify() == nil {
		t.Fatal("ABI/platform mismatch")
	}
	b, e := os.ReadFile(c.Lock)
	if e != nil {
		t.Fatal(e)
	}
	var lock map[string]any
	if json.Unmarshal(b, &lock) != nil {
		t.Fatal("lock")
	}
	files := lock["files"].([]any)
	files[len(files)-1] = files[0]
	b, _ = json.Marshal(lock)
	wrong = c
	wrong.Lock = filepath.Join(t.TempDir(), "duplicate.json")
	wrong.LockSHA = digest(b)
	if os.WriteFile(wrong.Lock, b, 0600) != nil {
		t.Fatal("write")
	}
	if wrong.Verify() == nil {
		t.Fatal("duplicate lock paths admitted")
	}
	frame, e := os.ReadFile(filepath.Join(os.Getenv("CAPTURE_TEST_FRAMES"), "serial-code128.png.gray8"))
	if e != nil {
		t.Fatal(e)
	}
	for _, payload := range []string{
		`{"schema":"elite-qr-image-decode/v1","status":"DECODED","candidates":[{"text":"SKU-1","symbology":"CODE_128","points":[[-1,0],[1,1],[2,2],[3,3]]}],"authorized":false,"business_effect":false}`,
		`{"schema":"elite-qr-image-decode/v1","status":"DECODED","candidates":[{"text":"evil\ntext","symbology":"CODE_128","points":[[0,0],[1,1],[2,2],[3,3]]}],"authorized":false,"business_effect":false}`,
		`{"schema":"elite-qr-image-decode/v1","status":"DECODED","candidates":[],"authorized":true,"business_effect":false}`,
	} {
		p := filepath.Join(t.TempDir(), "worker.py")
		code := []byte("print(" + string(mustJSON(payload)) + ")\n")
		if os.WriteFile(p, code, 0600) != nil {
			t.Fatal("write")
		}
		cfg := c
		cfg.Worker = p
		cfg.WorkerSHA = digest(code)
		d, e := NewDecoder(cfg)
		if e != nil {
			t.Fatal(e)
		}
		r := d.Decode(context.Background(), frame)
		if r.Status != "FAILED" && r.Status != "REJECTED" {
			t.Fatal("invalid worker result admitted", r)
		}
	}
	d, e := NewDecoder(c)
	if e != nil {
		t.Fatal(e)
	}
	d.slots <- struct{}{}
	d.slots <- struct{}{}
	if r := d.Decode(context.Background(), frame); r.Status != "BUSY" {
		t.Fatal("bounded capacity")
	}
	<-d.slots
	<-d.slots
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if r := d.Decode(ctx, frame); r.Status != "TIMEOUT" {
		t.Fatal("cancel bound", r)
	}
}
func mustJSON(v any) []byte { b, _ := json.Marshal(v); return b }
