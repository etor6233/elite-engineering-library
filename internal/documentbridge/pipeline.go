// Package documentbridge is AUTHORED bounded orchestration glue around the
// original secure-file gate and pinned official AWS Textract SDK adapter.
package documentbridge

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode/utf8"

	runtime "example.com/elite/aws-textract-document-runtime/textractruntime"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
)

const MaxBytes = 2 * 1024 * 1024
const PublicFixtureSHA = "489f0c63b6a05e0ae0fcfbf29121299fc2ccacab35d5cd40a5d63d42764078cb"
const Class = "supplier-invoice"

var ErrContract = errors.New("document contract rejected")
var ErrSecurity = errors.New("document security gate rejected")
var ErrExtraction = errors.New("document extraction unavailable")

type Scope struct {
	TenantID       string `json:"tenant_id"`
	OrganizationID string `json:"organization_id"`
	ProfileSHA     string `json:"profile_sha256"`
	Mode           string `json:"mode"`
}
type Original struct {
	Name          string
	SHA256        string
	Bytes         []byte
	ClassID       string
	SchemaVersion string
}
type Fields struct {
	InvoiceNumber string `json:"invoice_number"`
	Vendor        string `json:"vendor"`
	Total         string `json:"total"`
	Currency      string `json:"currency"`
}
type Evidence struct {
	Security, Provider, Receipt []byte
	Suggested                   Fields
	TypedFields                 json.RawMessage
	Mode                        string
}
type Processor interface {
	Run(context.Context, Original) (Evidence, error)
}
type SecurityStage func(context.Context, string, string) error

func Hash(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func Text(v string, n int) bool {
	return v != "" && len(v) <= n && utf8.ValidString(v) && !strings.ContainsRune(v, 0) && strings.TrimSpace(v) == v
}
func Hex(v string) bool {
	b, e := hex.DecodeString(v)
	return e == nil && len(b) == 32 && strings.ToLower(v) == v
}
func (v Fields) Validate() error {
	// Values remain reviewed text, never amounts to post, tax calculations or OCR truth.
	if !Text(v.InvoiceNumber, 128) || !Text(v.Vendor, 512) || !Text(v.Total, 64) || len(v.Currency) != 3 {
		return ErrContract
	}
	for _, r := range v.Currency {
		if r < 'A' || r > 'Z' {
			return ErrContract
		}
	}
	return nil
}
func (v Original) Validate(mode string) error {
	if !Text(v.Name, 128) || strings.ContainsAny(v.Name, "/\\:") || v.Name == "." || v.Name == ".." || len(v.Bytes) == 0 || len(v.Bytes) > MaxBytes || Hash(v.Bytes) != v.SHA256 {
		return ErrContract
	}
	ext := strings.ToLower(filepath.Ext(v.Name))
	mime := http.DetectContentType(v.Bytes)
	if !(ext == ".pdf" && mime == "application/pdf" || (ext == ".jpg" || ext == ".jpeg") && mime == "image/jpeg") {
		return ErrContract
	}
	if mode == "TYPED_FIXTURE" {
		return validTypedOriginal(v)
	}
	if v.ClassID != "" && v.ClassID != Class || v.SchemaVersion != "" && v.SchemaVersion != "1" {
		return ErrContract
	}
	if mode != "PROVIDER" && mode != "FIXTURE" || mode == "FIXTURE" && v.SHA256 != PublicFixtureSHA {
		return ErrContract
	}
	return nil
}

type Pipeline struct {
	client                                  runtime.Client
	security                                SecurityStage
	profile, profileSHA, root, region, mode string
}

func NewPipeline(client runtime.Client, security SecurityStage, profile, expectedSHA, root, mode string) (*Pipeline, error) {
	if client == nil || security == nil || !Hex(expectedSHA) || (mode != "PROVIDER" && mode != "FIXTURE") {
		return nil, ErrContract
	}
	b, e := regular(profile, 32768)
	if e != nil || Hash(b) != expectedSHA {
		return nil, ErrContract
	}
	region, e := runtime.ApprovedRegion(profile, Class)
	if e != nil {
		return nil, ErrContract
	}
	info, e := os.Lstat(root)
	if e != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil, ErrContract
	}
	return &Pipeline{client, security, profile, expectedSHA, root, region, mode}, nil
}
func regular(path string, max int64) ([]byte, error) {
	info, e := os.Lstat(path)
	if e != nil || !info.Mode().IsRegular() || info.Size() > max {
		return nil, ErrContract
	}
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, max+1))
	if e != nil || int64(len(b)) > max {
		return nil, ErrContract
	}
	return b, nil
}
func (p *Pipeline) Run(ctx context.Context, v Original) (Evidence, error) {
	if p == nil || v.Validate(p.mode) != nil {
		return Evidence{}, ErrContract
	}
	profile, e := regular(p.profile, 32768)
	if e != nil || Hash(profile) != p.profileSHA {
		return Evidence{}, ErrContract
	}
	dir, e := os.MkdirTemp(p.root, "document-")
	if e != nil {
		return Evidence{}, e
	}
	// The caller persists all successful evidence in PostgreSQL. Temporary originals
	// are removed on return; interrupted attempts require host orphan cleanup policy.
	defer os.RemoveAll(dir)
	input := filepath.Join(dir, "original"+strings.ToLower(filepath.Ext(v.Name)))
	if e = os.WriteFile(input, v.Bytes, 0600); e != nil {
		return Evidence{}, e
	}
	securityDir := filepath.Join(dir, "security")
	if e = p.security(ctx, input, securityDir); e != nil {
		return Evidence{}, ErrSecurity
	}
	security, e := regular(filepath.Join(securityDir, "security-receipt.json"), 65536)
	if e != nil {
		return Evidence{}, ErrSecurity
	}
	receipt, e := runtime.AnalyzeToEvidence(ctx, p.client, runtime.Request{Region: p.region, DocumentClass: Class, InputPath: input, OutputDirectory: filepath.Join(dir, "extraction"), ProfilePath: p.profile, SecurityReceiptPath: filepath.Join(securityDir, "security-receipt.json"), MaxBytes: MaxBytes})
	if e != nil {
		return Evidence{}, ErrExtraction
	}
	raw, e := regular(filepath.Join(dir, "extraction", "provider-response.json"), 4*1024*1024)
	if e != nil {
		return Evidence{}, ErrContract
	}
	original, e := regular(input, MaxBytes)
	if e != nil || !bytes.Equal(original, v.Bytes) || receipt.InputSHA256 != v.SHA256 || receipt.DocumentProfileSHA256 != p.profileSHA || receipt.SecurityReceiptSHA256 != Hash(security) || receipt.ProviderResponseSHA256 != Hash(raw) || receipt.AutomaticStorageAuthorized {
		return Evidence{}, ErrContract
	}
	var output textract.AnalyzeExpenseOutput
	if json.Unmarshal(raw, &output) != nil || len(output.ExpenseDocuments) != 1 {
		return Evidence{}, ErrContract
	}
	fields := Fields{}
	seen := map[string]bool{}
	for _, f := range output.ExpenseDocuments[0].SummaryFields {
		if f.Type == nil || f.ValueDetection == nil {
			continue
		}
		key := aws.ToString(f.Type.Text)
		if key != "INVOICE_RECEIPT_ID" && key != "VENDOR_NAME" && key != "TOTAL" {
			continue
		}
		if seen[key] {
			return Evidence{}, ErrContract
		}
		seen[key] = true
		switch key {
		case "INVOICE_RECEIPT_ID":
			fields.InvoiceNumber = aws.ToString(f.ValueDetection.Text)
		case "VENDOR_NAME":
			fields.Vendor = aws.ToString(f.ValueDetection.Text)
		case "TOTAL":
			fields.Total = aws.ToString(f.ValueDetection.Text)
			if f.Currency != nil {
				fields.Currency = aws.ToString(f.Currency.Code)
			}
		}
	}
	// Missing extraction values are preserved for manual correction; validation
	// of all four required fields belongs to the immutable review proposal.
	if len(fields.InvoiceNumber) > 128 || len(fields.Vendor) > 512 || len(fields.Total) > 64 || len(fields.Currency) > 3 {
		return Evidence{}, ErrContract
	}
	receiptBytes, e := json.Marshal(receipt)
	if e != nil {
		return Evidence{}, e
	}
	return Evidence{Security: security, Provider: raw, Receipt: receiptBytes, Suggested: fields, Mode: p.mode}, nil
}

type SecurityCommand struct {
	Python, PythonSHA, Script, ScriptSHA, Policy, PolicySHA string
	Magika, ClamScan, ClamDatabase, Yara, Rules             string
}

func (c SecurityCommand) Stage() (SecurityStage, error) {
	for _, v := range []struct{ path, hash string }{{c.Python, c.PythonSHA}, {c.Script, c.ScriptSHA}, {c.Policy, c.PolicySHA}} {
		if !filepath.IsAbs(v.path) || !Hex(v.hash) {
			return nil, ErrContract
		}
		b, e := regular(v.path, 64*1024*1024)
		if e != nil || Hash(b) != v.hash {
			return nil, ErrContract
		}
	}
	for _, p := range []string{c.Magika, c.ClamScan, c.ClamDatabase, c.Yara, c.Rules} {
		if !filepath.IsAbs(p) {
			return nil, ErrContract
		}
	}
	return func(ctx context.Context, input, output string) error {
		// Recheck executable/script/policy hashes per attempt. Native binary hashes,
		// official database age, YARA rule hash and detector limits remain gate-owned.
		for _, v := range []struct{ path, hash string }{{c.Python, c.PythonSHA}, {c.Script, c.ScriptSHA}, {c.Policy, c.PolicySHA}} {
			b, e := regular(v.path, 64*1024*1024)
			if e != nil || Hash(b) != v.hash {
				return ErrSecurity
			}
		}
		cmd := exec.CommandContext(ctx, c.Python, "-I", "-X", "utf8", "-B", c.Script, "--input", input, "--output", output, "--policy", c.Policy, "--magika-exe", c.Magika, "--clamscan-exe", c.ClamScan, "--clamav-database", c.ClamDatabase, "--yara-exe", c.Yara, "--compiled-rules", c.Rules)
		cmd.Env = []string{"PYTHONIOENCODING=utf-8"}
		for _, k := range []string{"SystemRoot", "TEMP", "TMP"} {
			if v := os.Getenv(k); v != "" {
				cmd.Env = append(cmd.Env, k+"="+v)
			}
		}
		// No raw detector output/private paths are copied to application logs.
		if e := cmd.Run(); e != nil {
			return fmt.Errorf("%w: detector process", ErrSecurity)
		}
		return nil
	}, nil
}
