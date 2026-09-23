// Package capture is AUTHORED bounded process glue around the pinned ZXing wheel.
package capture

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const MaxPixels = 4 * 1024 * 1024
const MaxFrameBytes = MaxPixels + 8

var ErrContract = errors.New("capture contract rejected")

type Candidate struct {
	Text      string   `json:"text"`
	Symbology string   `json:"symbology"`
	Points    [][2]int `json:"points"`
}
type Result struct {
	Schema         string      `json:"schema"`
	Status         string      `json:"status"`
	Reason         *string     `json:"reason"`
	Candidates     []Candidate `json:"candidates"`
	Authorized     bool        `json:"authorized"`
	BusinessEffect bool        `json:"business_effect"`
}

func result(status, reason string) Result {
	return Result{Schema: "elite-qr-image-decode/v1", Status: status, Reason: &reason, Candidates: []Candidate{}}
}

type Config struct {
	Schema        string `json:"schema"`
	OS            string `json:"os"`
	Arch          string `json:"arch"`
	PythonVersion string `json:"python_version"`
	WheelTag      string `json:"wheel_tag"`
	Python        string `json:"python"`
	PythonSHA     string `json:"python_sha256"`
	Worker        string `json:"worker"`
	WorkerSHA     string `json:"worker_sha256"`
	Runtime       string `json:"runtime"`
	Lock          string `json:"runtime_lock"`
	LockSHA       string `json:"runtime_lock_sha256"`
}
type Decoder struct {
	config Config
	slots  chan struct{}
}

func digest(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func exact(path, sha string, max int64) ([]byte, error) {
	if !filepath.IsAbs(path) || len(sha) != 64 {
		return nil, ErrContract
	}
	info, e := os.Lstat(path)
	if e != nil || !info.Mode().IsRegular() || info.Size() > max {
		return nil, ErrContract
	}
	b, e := os.ReadFile(path)
	if e != nil || digest(b) != sha {
		return nil, ErrContract
	}
	return b, nil
}
func (c Config) Verify() error {
	if c.Schema != "capture-runtime-profile/v1" || c.OS != runtime.GOOS || c.Arch != runtime.GOARCH || c.PythonVersion != "3.14" || c.WheelTag != "cp312-abi3-win_amd64" {
		return ErrContract
	}
	if c.OS != "windows" || c.Arch != "amd64" {
		return ErrContract
	}
	if _, e := exact(c.Python, c.PythonSHA, 32<<20); e != nil {
		return e
	}
	if _, e := exact(c.Worker, c.WorkerSHA, 1<<20); e != nil {
		return e
	}
	b, e := exact(c.Lock, c.LockSHA, 65536)
	if e != nil {
		return e
	}
	var lock struct {
		Schema string `json:"schema"`
		Files  []struct {
			Path string `json:"path"`
			SHA  string `json:"sha256"`
		} `json:"files"`
	}
	if json.Unmarshal(b, &lock) != nil || lock.Schema != "elite-exact-qr-wheel-files/v1" || len(lock.Files) != 7 || !filepath.IsAbs(c.Runtime) {
		return ErrContract
	}
	info, e := os.Lstat(c.Runtime)
	if e != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return ErrContract
	}
	expected := map[string]bool{}
	for _, f := range lock.Files {
		rel := filepath.FromSlash(f.Path)
		if filepath.IsAbs(rel) || strings.Contains(f.Path, "..") || strings.Contains(f.Path, "\\") {
			return ErrContract
		}
		if _, e = exact(filepath.Join(c.Runtime, rel), f.SHA, 32<<20); e != nil {
			return e
		}
		if expected[filepath.Clean(rel)] {
			return ErrContract
		}
		expected[filepath.Clean(rel)] = true
		if strings.HasSuffix(f.Path, ".dist-info/WHEEL") {
			raw, e := os.ReadFile(filepath.Join(c.Runtime, rel))
			if e != nil || !strings.Contains(string(raw), "Tag: "+c.WheelTag+"\n") {
				return ErrContract
			}
		}
	}
	n := 0
	return filepath.WalkDir(c.Runtime, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return ErrContract
		}
		n++
		if n > 64 || d.Type()&os.ModeSymlink != 0 {
			return ErrContract
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".py" || ext == ".pyc" || ext == ".pyd" || ext == ".dll" || ext == ".so" || ext == ".pth" {
			rel, e := filepath.Rel(c.Runtime, path)
			if e != nil || !expected[rel] {
				return ErrContract
			}
		}
		return nil
	})
}
func NewDecoder(c Config) (*Decoder, error) {
	if e := c.Verify(); e != nil {
		return nil, e
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, c.Python, "-I", "-B", "-c", "import sys;print('.'.join(map(str,sys.version_info[:2])))")
	configureHidden(cmd)
	raw, e := cmd.Output()
	if e != nil || strings.TrimSpace(string(raw)) != c.PythonVersion {
		return nil, ErrContract
	}
	return &Decoder{config: c, slots: make(chan struct{}, 2)}, nil
}
func ValidFrame(b []byte) bool {
	if len(b) < 9 || len(b) > MaxFrameBytes || !bytes.Equal(b[:4], []byte("EGY1")) {
		return false
	}
	w, h := int(binary.BigEndian.Uint16(b[4:6])), int(binary.BigEndian.Uint16(b[6:8]))
	return w > 20 && h > 20 && w <= 4096 && h <= 4096 && w*h <= MaxPixels && len(b) == 8+w*h
}

type boundedOutput struct {
	b bytes.Buffer
	n int
}

func (w *boundedOutput) Write(b []byte) (int, error) {
	if len(b) > w.n-w.b.Len() {
		return 0, ErrContract
	}
	return w.b.Write(b)
}
func (d *Decoder) Decode(ctx context.Context, b []byte) Result {
	if !ValidFrame(b) {
		return result("REJECTED", "FRAME_FORMAT_OR_SIZE")
	}
	select {
	case d.slots <- struct{}{}:
		defer func() { <-d.slots }()
	default:
		return result("BUSY", "CAPTURE_CAPACITY")
	}
	if e := d.config.Verify(); e != nil {
		return result("FAILED", "RUNTIME_MISMATCH")
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, d.config.Python, "-I", "-B", d.config.Worker, "--runtime-root", d.config.Runtime)
	configureHidden(cmd)
	cmd.Stdin = bytes.NewReader(b)
	out := boundedOutput{n: 65536}
	cmd.Stdout = &out
	cmd.Stderr = io.Discard
	if e := cmd.Run(); e != nil {
		if ctx.Err() != nil {
			return result("TIMEOUT", "DECODE_TIMEOUT")
		}
		return result("FAILED", "DECODER_PROCESS")
	}
	var r Result
	if json.Unmarshal(out.b.Bytes(), &r) != nil || r.Schema != "elite-qr-image-decode/v1" || r.Authorized || r.BusinessEffect || len(r.Candidates) > 8 {
		return result("FAILED", "DECODER_PROTOCOL")
	}
	switch r.Status {
	case "DECODED":
		if len(r.Candidates) != 1 {
			return result("FAILED", "DECODER_PROTOCOL")
		}
	case "MULTIPLE":
		if len(r.Candidates) < 2 {
			return result("FAILED", "DECODER_PROTOCOL")
		}
	case "NOT_FOUND", "REJECTED", "FAILED":
		if len(r.Candidates) != 0 {
			return result("FAILED", "DECODER_PROTOCOL")
		}
	default:
		return result("FAILED", "DECODER_PROTOCOL")
	}
	for _, c := range r.Candidates {
		if len(c.Text) == 0 || len(c.Text) > 1024 || len(c.Points) != 4 || c.Symbology != "QR_CODE" && c.Symbology != "CODE_128" && c.Symbology != "EAN_13" {
			return result("FAILED", "DECODER_PROTOCOL")
		}
		if !utf8.ValidString(c.Text) {
			return result("FAILED", "DECODER_PROTOCOL")
		}
		for _, r := range c.Text {
			if unicode.IsControl(r) {
				return result("REJECTED", "INVALID_DECODE_RESULT")
			}
		}
		w, h := int(binary.BigEndian.Uint16(b[4:6])), int(binary.BigEndian.Uint16(b[6:8]))
		for _, p := range c.Points {
			if p[0] < 0 || p[0] > w || p[1] < 0 || p[1] > h {
				return result("FAILED", "DECODER_PROTOCOL")
			}
		}
	}
	return r
}
