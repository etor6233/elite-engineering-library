package merchantbridge

// AUTHORED bounded process binding following the admitted social/stored-value
// launchers. No shell, inherited proxies, Python paths or secrets in arguments.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Process struct {
	Python, PythonSHA256           string
	Script, ScriptSHA256           string
	OwnerSHA256, RuntimeLockSHA256 string
	Mode, FixtureOrigin            string
	CredentialsFile                string
}

func Hash(b []byte) string { s := sha256.Sum256(b); return hex.EncodeToString(s[:]) }
func exactFile(path, expected string, maximum int64) error {
	f, e := os.Open(path)
	if e != nil {
		return ErrBinding
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, maximum+1))
	if e != nil || int64(len(b)) > maximum || len(expected) != 64 || Hash(b) != expected {
		return ErrBinding
	}
	return nil
}
func (p Process) Validate() error {
	if !filepath.IsAbs(p.Python) || !filepath.IsAbs(p.Script) || exactFile(p.Python, p.PythonSHA256, 4*1024*1024) != nil || exactFile(p.Script, p.ScriptSHA256, 65536) != nil || exactFile(filepath.Join(filepath.Dir(p.Script), "sync_product.py"), p.OwnerSHA256, 65536) != nil || exactFile(filepath.Join(filepath.Dir(p.Script), "connected-runtime.lock.json"), p.RuntimeLockSHA256, 1048576) != nil {
		return ErrBinding
	}
	if p.Mode == "LOCAL_FIXTURES" {
		u, e := url.Parse(p.FixtureOrigin)
		if e != nil || u.Scheme != "http" || u.Hostname() != "127.0.0.1" || u.Port() == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || p.CredentialsFile != "" {
			return ErrBinding
		}
	} else if p.Mode == "CREDENTIALS" {
		if p.FixtureOrigin != "" || !filepath.IsAbs(p.CredentialsFile) {
			return ErrBinding
		}
		if info, e := os.Stat(p.CredentialsFile); e != nil || info.IsDir() {
			return ErrBinding
		}
	} else {
		return ErrBinding
	}
	return nil
}

type Result struct {
	State          string          `json:"state"`
	Operation      string          `json:"operation,omitempty"`
	Resource       string          `json:"resource,omitempty"`
	Response       json.RawMessage `json:"response,omitempty"`
	ResponseSHA256 string          `json:"response_sha256,omitempty"`
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
func (p Process) Execute(ctx context.Context, operation string, r Intent) (Result, error) {
	if p.Validate() != nil || operation != "INSERT" && operation != "GET" {
		return Result{}, ErrBinding
	}
	input, e := json.Marshal(struct {
		Operation string `json:"operation"`
		Intent    Intent `json:"intent"`
		Mode      string `json:"mode"`
		Origin    string `json:"fixture_origin"`
	}{operation, r, p.Mode, p.FixtureOrigin})
	if e != nil || len(input) > 32768 {
		return Result{}, ErrBinding
	}
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, p.Python, "-I", "-S", "-B", p.Script, p.OwnerSHA256, p.RuntimeLockSHA256)
	c.Stdin = bytes.NewReader(input)
	c.Dir = filepath.Dir(p.Script)
	for _, key := range []string{"SYSTEMROOT", "WINDIR", "TEMP", "TMP"} {
		if v := os.Getenv(key); v != "" {
			c.Env = append(c.Env, key+"="+v)
		}
	}
	if p.Mode == "CREDENTIALS" {
		if strings.ContainsAny(p.CredentialsFile, "\r\n\x00") {
			return Result{}, ErrBinding
		}
		c.Env = append(c.Env, "GOOGLE_APPLICATION_CREDENTIALS="+p.CredentialsFile)
	}
	var buffer output
	c.Stdout = &buffer
	c.Stderr = io.Discard
	if c.Run() != nil || buffer.exceeded {
		return Result{State: "UNKNOWN", Operation: operation}, nil
	}
	var result Result
	if Decode(buffer.Bytes(), &result) != nil {
		return Result{State: "UNKNOWN", Operation: operation}, nil
	}
	if result.State == "UNKNOWN" {
		return Result{State: "UNKNOWN", Operation: operation}, nil
	}
	if result.Operation != operation {
		return Result{}, ErrBinding
	}
	if result.State == "OBSERVED" {
		if result.Resource != r.Resource || len(result.Response) == 0 || Hash(result.Response) != result.ResponseSHA256 {
			return Result{}, ErrBinding
		}
	} else if result.State != "REJECTED" && result.State != "NOT_FOUND" {
		return Result{}, ErrBinding
	}
	return result, nil
}
