package documentbridge

// AUTHORED simulator. No detector binaries or live OCR run here. FIXTURE mode is
// permanent on original/evidence/commit records, restricted to one public hash.
import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const FixtureResponse = `{"DocumentMetadata":{"Pages":1},"ExpenseDocuments":[{"ExpenseIndex":1,"SummaryFields":[{"Type":{"Text":"INVOICE_RECEIPT_ID"},"ValueDetection":{"Text":"FIXTURE-INVOICE-1"}},{"Type":{"Text":"VENDOR_NAME"},"ValueDetection":{"Text":"Fixture vendor"}},{"Type":{"Text":"TOTAL"},"ValueDetection":{"Text":"123.45"},"Currency":{"Code":"USD"}}],"LineItemGroups":[]}]}`

type fixtureHTTP struct{}

func (fixtureHTTP) Do(r *http.Request) (*http.Response, error) {
	if r.Method != "POST" || r.Header.Get("X-Amz-Target") != "Textract.AnalyzeExpense" || !strings.HasPrefix(r.Header.Get("Authorization"), "AWS4-HMAC-SHA256 Credential=AKIDFIXTURE/") {
		return nil, ErrContract
	}
	var wire struct{ Document struct{ Bytes []byte } }
	b, e := io.ReadAll(io.LimitReader(r.Body, 3*1024*1024))
	if e != nil || json.Unmarshal(b, &wire) != nil || Hash(wire.Document.Bytes) != PublicFixtureSHA {
		return nil, ErrContract
	}
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/x-amz-json-1.1"}, "X-Amzn-Requestid": []string{"simulated-public-document-only"}}, Body: io.NopCloser(strings.NewReader(FixtureResponse)), Request: r}, nil
}
func FixtureSecurity(ctx context.Context, input, output string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	b, e := regular(input, MaxBytes)
	if e != nil || Hash(b) != PublicFixtureSHA {
		return ErrSecurity
	}
	synthetic := Hash([]byte("SIMULATED DETECTOR CONTRACT; NO SCANNER EXECUTED"))
	now := time.Now().UTC().Format(time.RFC3339Nano)
	receipt := map[string]any{
		"schema": "elite-secure-local-file-receipt/v1", "decision": "ADMITTED", "reason": "ALL_SELECTED_GATES_PASSED", "started_at": now, "completed_at": now,
		"approval_id": "user-authorized-public-fixture-only", "policy_sha256": synthetic, "input": map[string]any{"sha256": PublicFixtureSHA, "bytes": len(b)}, "business_storage_authorized": false,
		"security_claim": "SIMULATED detector contract for exact public fixture; no scanner executed; no malware or private-file admission proof",
		"clamav":         map[string]any{"exit_code": 0, "output_sha256": synthetic}, "content_type": map[string]any{"label": "jpeg", "mime_type": "image/jpeg", "score": 1.0},
		"yara_x": map[string]any{"compiled_rules_sha256": synthetic, "matching_rules": []string{}},
		"tools":  map[string]any{"clamscan": map[string]any{"sha256": synthetic, "version": "ClamAV 1.5.4/SIMULATED"}, "magika": map[string]any{"sha256": synthetic, "version": "magika 1.1.0"}, "yara_x": map[string]any{"sha256": synthetic, "version": "1.20.0"}},
	}
	raw, e := json.Marshal(receipt)
	if e != nil {
		return e
	}
	if e = os.Mkdir(output, 0700); e != nil {
		return e
	}
	return os.WriteFile(filepath.Join(output, "security-receipt.json"), raw, 0600)
}
func NewFixturePipeline(profile, sha, root string) (*Pipeline, error) {
	client := textract.NewFromConfig(aws.Config{Region: "us-east-1", Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
		return aws.Credentials{AccessKeyID: "AKIDFIXTURE", SecretAccessKey: "fixture-only-no-live-authority"}, nil
	}), HTTPClient: fixtureHTTP{}}, func(o *textract.Options) { o.RetryMaxAttempts = 1 })
	return NewPipeline(client, FixtureSecurity, profile, sha, root, "FIXTURE")
}
