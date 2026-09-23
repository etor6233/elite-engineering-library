package socialbridge

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Process invokes only exact materialized code. Credentials are environment
// values, never command arguments, receipts, logs or approval payloads.
type Process struct {
	Python, Script, ScriptSHA256 string
	Environment                  []string
}

func (p Process) Preflight(ctx context.Context) error {
	if p.Validate() != nil {
		return ErrBinding
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, p.Python, "-I", "-B", p.Script, "--preflight")
	for _, k := range []string{"SYSTEMROOT", "WINDIR", "TEMP", "TMP"} {
		if v := os.Getenv(k); v != "" {
			c.Env = append(c.Env, k+"="+v)
		}
	}
	for _, v := range p.Environment {
		k, _, ok := strings.Cut(v, "=")
		if !ok || (k != "META_APP_ID" && k != "META_APP_SECRET" && k != "META_PAGE_ACCESS_TOKEN") {
			return ErrBinding
		}
		c.Env = append(c.Env, v)
	}
	var out boundedOutput
	c.Stdout = &out
	c.Stderr = io.Discard
	if c.Run() != nil || out.exceeded || out.String() != `{"state":"PASS"}` {
		return ErrBinding
	}
	return nil
}

func (p Process) Validate() error {
	if !filepath.IsAbs(p.Python) || !filepath.IsAbs(p.Script) {
		return ErrBinding
	}
	if _, e := os.Stat(p.Python); e != nil {
		return ErrBinding
	}
	b, e := os.ReadFile(p.Script)
	if e != nil || !hashRE.MatchString(p.ScriptSHA256) || Hash(b) != p.ScriptSHA256 {
		return ErrBinding
	}
	return nil
}

type boundedOutput struct {
	bytes.Buffer
	exceeded bool
}

func (b *boundedOutput) Write(v []byte) (int, error) {
	if b.Len()+len(v) > 32768 {
		b.exceeded = true
		return 0, io.ErrShortBuffer
	}
	return b.Buffer.Write(v)
}
func (p Process) Execute(ctx context.Context, r Request, reconcile bool) (Result, error) {
	if p.Validate() != nil {
		return Result{}, ErrBinding
	}
	input, _ := json.Marshal(struct {
		Request   Request `json:"request"`
		Reconcile bool    `json:"reconcile"`
	}{r, reconcile})
	ctx, cancel := context.WithTimeout(ctx, 35*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, p.Python, "-I", "-B", p.Script)
	c.Stdin = bytes.NewReader(input)
	// Explicit allowlist. Proxy/SSL overrides, PYTHONPATH, debug and arbitrary
	// inherited service credentials never cross this child process boundary.
	for _, k := range []string{"SYSTEMROOT", "WINDIR", "TEMP", "TMP"} {
		if v := os.Getenv(k); v != "" {
			c.Env = append(c.Env, k+"="+v)
		}
	}
	for _, v := range p.Environment {
		k, _, ok := strings.Cut(v, "=")
		if !ok || (k != "META_APP_ID" && k != "META_APP_SECRET" && k != "META_PAGE_ACCESS_TOKEN") {
			return Result{}, ErrBinding
		}
		c.Env = append(c.Env, v)
	}
	var out boundedOutput
	c.Stdout = &out
	c.Stderr = io.Discard
	if c.Run() != nil || out.exceeded {
		return Result{Unknown: true, Code: "PROCESS_RESULT_UNCONFIRMED"}, nil
	}
	var result Result
	if Decode(out.Bytes(), &result) != nil {
		return Result{Unknown: true, Code: "PROCESS_RESULT_INVALID"}, nil
	}
	if result.Receipt != nil {
		if result.Unknown || result.Receipt.Validate(r.Intent) != nil {
			return Result{Unknown: true, Code: "PROCESS_RECEIPT_INVALID"}, nil
		}
	} else if !result.Unknown {
		return Result{Unknown: true, Code: "PROCESS_RESULT_INVALID"}, nil
	}
	if result.ProviderReference != "" && !ValidReference(r.Intent.PageID, result.ProviderReference) {
		return Result{Unknown: true, Code: "PROCESS_REFERENCE_INVALID"}, nil
	}
	return result, nil
}
