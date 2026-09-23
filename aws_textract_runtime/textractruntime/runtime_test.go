package textractruntime

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

const testHex = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

type fakeClient struct {
	expense       *textract.AnalyzeExpenseOutput
	document      *textract.AnalyzeDocumentOutput
	err           error
	expenseInput  *textract.AnalyzeExpenseInput
	documentInput *textract.AnalyzeDocumentInput
	calls         int
}

func (f *fakeClient) AnalyzeExpense(_ context.Context, input *textract.AnalyzeExpenseInput, _ ...func(*textract.Options)) (*textract.AnalyzeExpenseOutput, error) {
	f.calls++
	f.expenseInput = input
	return f.expense, f.err
}
func (f *fakeClient) AnalyzeDocument(_ context.Context, input *textract.AnalyzeDocumentInput, _ ...func(*textract.Options)) (*textract.AnalyzeDocumentOutput, error) {
	f.calls++
	f.documentInput = input
	return f.document, f.err
}

type fixtureState struct{ root, input, profile, security string }

func newFixture(t *testing.T, item profileClass) fixtureState {
	t.Helper()
	root := t.TempDir()
	input := filepath.Join(root, "invoice.pdf")
	if err := os.WriteFile(input, []byte("%PDF-1.7 fixture"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := fixtureState{root: root, input: input, profile: filepath.Join(root, "profile.json"), security: filepath.Join(root, "security.json")}
	writeProfile(t, state.profile, item)
	writeSecurity(t, state.security, state.input, "ADMITTED", "application/pdf", true)
	return state
}

func expenseClass() profileClass {
	return profileClass{ID: "supplier-invoice", Decision: "REQUIRED", Operation: "AnalyzeExpense", Features: []string{}, Queries: []Query{}, ResponseSchemaVersion: "supplier-invoice/v1", AutomaticStorage: false}
}

func writeProfile(t *testing.T, path string, item profileClass) {
	t.Helper()
	profile := documentProfile{Schema: profileSchema, Provider: "Amazon Textract", SDK: sdkModule + "@" + sdkVersion, Region: "us-east-1", Classes: []profileClass{item}}
	body, err := json.Marshal(profile)
	if err != nil {
		t.Fatal(err)
	}
	body = append(body, '\n')
	if err = os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
}

func writeSecurity(t *testing.T, path, input, decision, mime string, includeTools bool) {
	t.Helper()
	payload, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(payload)
	receipt := map[string]any{
		"schema": securityReceiptSchema, "decision": decision, "reason": "ALL_SELECTED_GATES_PASSED", "started_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano), "completed_at": time.Now().UTC().Format(time.RFC3339Nano),
		"approval_id": "security-approval-test", "policy_sha256": testHex, "input": map[string]any{"sha256": hex.EncodeToString(sum[:]), "bytes": len(payload)}, "business_storage_authorized": false,
		"security_claim": "test fixture matching the secure gate contract", "clamav": map[string]any{"exit_code": 0, "output_sha256": testHex},
		"content_type": map[string]any{"label": "pdf", "mime_type": mime, "score": 1.0}, "yara_x": map[string]any{"compiled_rules_sha256": testHex, "matching_rules": []string{}},
	}
	if includeTools {
		receipt["tools"] = map[string]any{"clamscan": map[string]any{"sha256": testHex, "version": "ClamAV 1.5.4/test"}, "magika": map[string]any{"sha256": testHex, "version": "magika 1.1.0"}, "yara_x": map[string]any{"sha256": testHex, "version": "1.20.0"}}
	}
	body, err := json.Marshal(receipt)
	if err != nil {
		t.Fatal(err)
	}
	body = append(body, '\n')
	if err = os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
}

func request(state fixtureState, output string) Request {
	return Request{Region: "us-east-1", DocumentClass: "supplier-invoice", InputPath: state.input, OutputDirectory: filepath.Join(state.root, output), ProfilePath: state.profile, SecurityReceiptPath: state.security}
}

func TestAnalyzeExpensePreservesOfficialOutputAndEvidenceChain(t *testing.T) {
	state := newFixture(t, expenseClass())
	client := &fakeClient{expense: &textract.AnalyzeExpenseOutput{DocumentMetadata: &types.DocumentMetadata{Pages: aws.Int32(1)}, ExpenseDocuments: []types.ExpenseDocument{{ExpenseIndex: aws.Int32(1), SummaryFields: []types.ExpenseField{{Type: &types.ExpenseType{Text: aws.String("INVOICE_RECEIPT_ID")}, ValueDetection: &types.ExpenseDetection{Text: aws.String("INV-1"), Confidence: aws.Float32(99)}}}}}}}
	receipt, err := AnalyzeToEvidence(context.Background(), client, request(state, "expense"))
	if err != nil {
		t.Fatal(err)
	}
	if receipt.SDKVersion != "v1.45.0" || receipt.DocumentClass != "supplier-invoice" || receipt.ResponseSchemaVersion != "supplier-invoice/v1" || receipt.AutomaticStorageAuthorized || receipt.PageCount != 1 || receipt.ResultCount != 1 {
		t.Fatalf("unexpected receipt: %+v", receipt)
	}
	if !isSHA256(receipt.DocumentProfileSHA256) || !isSHA256(receipt.SecurityReceiptSHA256) || string(client.expenseInput.Document.Bytes) != "%PDF-1.7 fixture" {
		t.Fatal("official request or evidence hashes not preserved")
	}
	var raw map[string]any
	body, err := os.ReadFile(filepath.Join(state.root, "expense", "provider-response.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(body, &raw); err != nil {
		t.Fatal(err)
	}
	if _, ok := raw["ExpenseDocuments"]; !ok {
		t.Fatalf("complete modeled output missing: %s", body)
	}
}

func TestAnalyzeDocumentUsesOnlyProfileQueriesAndAdapter(t *testing.T) {
	item := profileClass{ID: "purchase-order", Decision: "REQUIRED", Operation: "AnalyzeDocument", Features: []string{"TABLES", "QUERIES"}, Queries: []Query{{Text: "What is the purchase order number?", Alias: "purchase_order_number", Pages: []string{"1"}}}, AdapterID: "adapter-1", AdapterVersion: "v7", ResponseSchemaVersion: "purchase-order/v3"}
	state := newFixture(t, item)
	client := &fakeClient{document: &textract.AnalyzeDocumentOutput{AnalyzeDocumentModelVersion: aws.String("1.0"), DocumentMetadata: &types.DocumentMetadata{Pages: aws.Int32(1)}, Blocks: []types.Block{{BlockType: types.BlockTypeQueryResult, Text: aws.String("PO-1"), Confidence: aws.Float32(98)}}}}
	req := request(state, "document")
	req.DocumentClass = "purchase-order"
	receipt, err := AnalyzeToEvidence(context.Background(), client, req)
	if err != nil {
		t.Fatal(err)
	}
	if receipt.AdapterVersion != "v7" || receipt.ModelVersion != "1.0" || client.documentInput.QueriesConfig == nil || len(client.documentInput.QueriesConfig.Queries) != 1 || client.documentInput.AdaptersConfig == nil || len(client.documentInput.AdaptersConfig.Adapters) != 1 {
		t.Fatalf("profile selection not preserved: %+v", receipt)
	}
}

func TestSecurityReceiptMustCompletelyMatchAndAdmit(t *testing.T) {
	state := newFixture(t, expenseClass())
	client := &fakeClient{expense: &textract.AnalyzeExpenseOutput{DocumentMetadata: &types.DocumentMetadata{Pages: aws.Int32(1)}, ExpenseDocuments: []types.ExpenseDocument{{}}}}
	cases := []struct {
		name, decision, mime string
		tools                bool
		want                 string
	}{{"rejected", "REJECTED", "application/pdf", true, "did not admit"}, {"wrong-mime", "ADMITTED", "image/png", true, "content type does not match"}, {"missing-tools", "ADMITTED", "application/pdf", false, "clamscan evidence is invalid"}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			writeSecurity(t, state.security, state.input, tc.decision, tc.mime, tc.tools)
			_, err := AnalyzeToEvidence(context.Background(), client, request(state, tc.name))
			if err == nil || !stringsContains(err.Error(), tc.want) {
				t.Fatalf("expected %q, got %v", tc.want, err)
			}
		})
	}
	if client.calls != 0 {
		t.Fatalf("provider called for rejected security evidence: %d", client.calls)
	}
}

func TestProfileBlocksUnapprovedDuplicateAndRegionMismatch(t *testing.T) {
	state := newFixture(t, expenseClass())
	client := &fakeClient{}
	item := expenseClass()
	item.Decision = "CONFIGURATION_REQUIRED"
	writeProfile(t, state.profile, item)
	if _, err := AnalyzeToEvidence(context.Background(), client, request(state, "blocked")); err == nil || !stringsContains(err.Error(), "not REQUIRED") {
		t.Fatalf("blocked class accepted: %v", err)
	}
	profile := documentProfile{Schema: profileSchema, Provider: "Amazon Textract", SDK: sdkModule + "@" + sdkVersion, Region: "us-east-1", Classes: []profileClass{expenseClass(), expenseClass()}}
	body, _ := json.Marshal(profile)
	_ = os.WriteFile(state.profile, append(body, '\n'), 0o600)
	if _, err := AnalyzeToEvidence(context.Background(), client, request(state, "duplicate")); err == nil || !stringsContains(err.Error(), "exactly once") {
		t.Fatalf("duplicate class accepted: %v", err)
	}
	writeProfile(t, state.profile, expenseClass())
	req := request(state, "region")
	req.Region = "eu-west-1"
	if _, err := AnalyzeToEvidence(context.Background(), client, req); err == nil || !stringsContains(err.Error(), "does not match") {
		t.Fatalf("region mismatch accepted: %v", err)
	}
	if client.calls != 0 {
		t.Fatal("provider called before profile admission")
	}
}

func TestTemplateEnablesNoDocumentClass(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("..", "document-profile.template.json"))
	if err != nil {
		t.Fatal(err)
	}
	var profile documentProfile
	if err = json.Unmarshal(body, &profile); err != nil {
		t.Fatal(err)
	}
	if profile.Region != "CONFIGURATION_REQUIRED" {
		t.Fatal("template region must require configuration")
	}
	for _, item := range profile.Classes {
		if item.Decision == "REQUIRED" || item.AutomaticStorage {
			t.Fatalf("template class enabled: %+v", item)
		}
	}
}

func TestRequestSurfaceCannotOverrideProfileOperation(t *testing.T) {
	typeOf := reflect.TypeOf(Request{})
	for _, name := range []string{"Mode", "Features", "Queries", "AdapterID", "AdapterVersion"} {
		if _, ok := typeOf.FieldByName(name); ok {
			t.Fatalf("unsafe request override field exists: %s", name)
		}
	}
}

func TestProviderFailureOrUnprovedPageLeavesNoOutput(t *testing.T) {
	state := newFixture(t, expenseClass())
	req := request(state, "failure")
	client := &fakeClient{err: errors.New("provider unavailable")}
	if _, err := AnalyzeToEvidence(context.Background(), client, req); err == nil {
		t.Fatal("provider error accepted")
	}
	if _, err := os.Stat(req.OutputDirectory); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("failed output committed")
	}
	client = &fakeClient{expense: &textract.AnalyzeExpenseOutput{ExpenseDocuments: []types.ExpenseDocument{{}}}}
	if _, err := AnalyzeToEvidence(context.Background(), client, req); err == nil || !stringsContains(err.Error(), "exactly one page") {
		t.Fatalf("unproved page accepted: %v", err)
	}
	if _, err := os.Stat(req.OutputDirectory); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("unproved output committed")
	}
}

func stringsContains(value, fragment string) bool {
	return len(fragment) == 0 || len(value) >= len(fragment) && func() bool {
		for i := 0; i+len(fragment) <= len(value); i++ {
			if value[i:i+len(fragment)] == fragment {
				return true
			}
		}
		return false
	}()
}
