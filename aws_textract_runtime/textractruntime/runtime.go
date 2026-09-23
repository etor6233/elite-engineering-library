package textractruntime

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

const (
	ModeExpense           = "expense"
	ModeDocument          = "document"
	MaxSyncBytes          = 10 * 1024 * 1024
	profileSchema         = "elite-aws-textract-document-profile/v2"
	securityReceiptSchema = "elite-secure-local-file-receipt/v1"
	sdkModule             = "github.com/aws/aws-sdk-go-v2/service/textract"
	sdkVersion            = "v1.45.0"
)

var allowedMIMETypes = map[string]string{
	".pdf": "application/pdf", ".png": "image/png", ".jpg": "image/jpeg", ".jpeg": "image/jpeg",
	".tif": "image/tiff", ".tiff": "image/tiff",
}
var featureTypes = map[string]types.FeatureType{
	"TABLES": types.FeatureTypeTables, "FORMS": types.FeatureTypeForms, "QUERIES": types.FeatureTypeQueries,
	"SIGNATURES": types.FeatureTypeSignatures, "LAYOUT": types.FeatureTypeLayout,
}

type Client interface {
	AnalyzeExpense(context.Context, *textract.AnalyzeExpenseInput, ...func(*textract.Options)) (*textract.AnalyzeExpenseOutput, error)
	AnalyzeDocument(context.Context, *textract.AnalyzeDocumentInput, ...func(*textract.Options)) (*textract.AnalyzeDocumentOutput, error)
}

type Query struct {
	Text  string   `json:"text"`
	Alias string   `json:"alias"`
	Pages []string `json:"pages,omitempty"`
}

type Request struct {
	Region              string
	DocumentClass       string
	InputPath           string
	OutputDirectory     string
	ProfilePath         string
	SecurityReceiptPath string
	MaxBytes            int64
}

type Receipt struct {
	Schema                     string `json:"schema"`
	CreatedAt                  string `json:"created_at"`
	Provider                   string `json:"provider"`
	SDKModule                  string `json:"sdk_module"`
	SDKVersion                 string `json:"sdk_version"`
	Operation                  string `json:"operation"`
	Region                     string `json:"region"`
	DocumentClass              string `json:"document_class"`
	ResponseSchemaVersion      string `json:"response_schema_version"`
	DocumentProfileSHA256      string `json:"document_profile_sha256"`
	SecurityReceiptSHA256      string `json:"security_receipt_sha256"`
	ContentType                string `json:"content_type"`
	InputFilename              string `json:"input_filename"`
	InputBytes                 int    `json:"input_bytes"`
	InputSHA256                string `json:"input_sha256"`
	ProviderResponseSHA256     string `json:"provider_response_sha256"`
	PageCount                  int32  `json:"page_count"`
	ResultCount                int    `json:"result_count"`
	ModelVersion               string `json:"model_version,omitempty"`
	AdapterID                  string `json:"adapter_id,omitempty"`
	AdapterVersion             string `json:"adapter_version,omitempty"`
	AutomaticStorageAuthorized bool   `json:"automatic_storage_authorized"`
}

type documentProfile struct {
	Schema   string         `json:"schema"`
	Provider string         `json:"provider"`
	SDK      string         `json:"sdk"`
	Region   string         `json:"region"`
	Classes  []profileClass `json:"classes"`
}

type profileClass struct {
	ID                    string   `json:"id"`
	Decision              string   `json:"decision"`
	Operation             string   `json:"operation"`
	Features              []string `json:"features"`
	Queries               []Query  `json:"queries"`
	AdapterID             string   `json:"adapter_id"`
	AdapterVersion        string   `json:"adapter_version"`
	ResponseSchemaVersion string   `json:"response_schema_version"`
	AutomaticStorage      bool     `json:"automatic_storage"`
}

type securityReceipt struct {
	Schema       string `json:"schema"`
	Decision     string `json:"decision"`
	Reason       string `json:"reason"`
	StartedAt    string `json:"started_at"`
	CompletedAt  string `json:"completed_at"`
	ApprovalID   string `json:"approval_id"`
	PolicySHA256 string `json:"policy_sha256"`
	Input        struct {
		SHA256 string `json:"sha256"`
		Bytes  int    `json:"bytes"`
	} `json:"input"`
	Tools map[string]struct {
		SHA256  string `json:"sha256"`
		Version string `json:"version"`
	} `json:"tools"`
	BusinessStorageAuthorized *bool  `json:"business_storage_authorized"`
	SecurityClaim             string `json:"security_claim"`
	ClamAV                    struct {
		ExitCode     int    `json:"exit_code"`
		OutputSHA256 string `json:"output_sha256"`
	} `json:"clamav"`
	ContentType struct {
		Label    string  `json:"label"`
		MIMEType string  `json:"mime_type"`
		Score    float64 `json:"score"`
	} `json:"content_type"`
	YaraX struct {
		CompiledRulesSHA256 string   `json:"compiled_rules_sha256"`
		MatchingRules       []string `json:"matching_rules"`
	} `json:"yara_x"`
}

type approvedSelection struct {
	Region                string
	DocumentClass         string
	Mode                  string
	Features              []types.FeatureType
	Queries               []types.Query
	Adapter               *types.Adapter
	AdapterID             string
	AdapterVersion        string
	ResponseSchemaVersion string
	ProfileSHA256         string
}

func ApprovedRegion(profilePath, documentClass string) (string, error) {
	selection, err := selectProfile(profilePath, documentClass, "")
	if err != nil {
		return "", err
	}
	return selection.Region, nil
}

func AnalyzeToEvidence(ctx context.Context, client Client, request Request) (receipt Receipt, err error) {
	if client == nil {
		return receipt, errors.New("client is required")
	}
	if request.MaxBytes == 0 {
		request.MaxBytes = MaxSyncBytes
	}
	if request.MaxBytes < 1 || request.MaxBytes > MaxSyncBytes {
		return receipt, errors.New("max_bytes must be within the Textract synchronous 10 MiB limit")
	}
	selection, err := selectProfile(request.ProfilePath, request.DocumentClass, request.Region)
	if err != nil {
		return receipt, err
	}
	info, err := os.Lstat(request.InputPath)
	if err != nil {
		return receipt, err
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return receipt, errors.New("input must be a regular non-symlink file")
	}
	absInput, err := filepath.Abs(request.InputPath)
	if err != nil {
		return receipt, err
	}
	mimeType, ok := allowedMIMETypes[strings.ToLower(filepath.Ext(absInput))]
	if !ok {
		return receipt, errors.New("unsupported document extension")
	}
	payload, err := os.ReadFile(absInput)
	if err != nil {
		return receipt, err
	}
	if len(payload) == 0 || int64(len(payload)) > request.MaxBytes {
		return receipt, errors.New("document byte size is outside the allowed range")
	}
	payloadSHA := hash(payload)
	securitySHA, err := verifySecurityReceipt(request.SecurityReceiptPath, payloadSHA, len(payload), mimeType)
	if err != nil {
		return receipt, err
	}
	absOutput, err := filepath.Abs(request.OutputDirectory)
	if err != nil {
		return receipt, err
	}
	if _, statErr := os.Stat(absOutput); statErr == nil {
		return receipt, errors.New("output_directory must not exist")
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return receipt, statErr
	}
	if err = os.MkdirAll(filepath.Dir(absOutput), 0o750); err != nil {
		return receipt, err
	}
	stage, err := os.MkdirTemp(filepath.Dir(absOutput), ".textract-analysis-stage-")
	if err != nil {
		return receipt, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = os.RemoveAll(stage)
		}
	}()

	var output any
	var pages int32
	var resultCount int
	var modelVersion string
	switch selection.Mode {
	case ModeExpense:
		result, callErr := client.AnalyzeExpense(ctx, &textract.AnalyzeExpenseInput{Document: &types.Document{Bytes: payload}})
		if callErr != nil {
			return receipt, callErr
		}
		if result == nil || len(result.ExpenseDocuments) == 0 {
			return receipt, errors.New("provider response contains no expense documents")
		}
		output, resultCount = result, len(result.ExpenseDocuments)
		if result.DocumentMetadata != nil {
			pages = aws.ToInt32(result.DocumentMetadata.Pages)
		}
	case ModeDocument:
		input := &textract.AnalyzeDocumentInput{Document: &types.Document{Bytes: payload}, FeatureTypes: selection.Features}
		if len(selection.Queries) > 0 {
			input.QueriesConfig = &types.QueriesConfig{Queries: selection.Queries}
		}
		if selection.Adapter != nil {
			input.AdaptersConfig = &types.AdaptersConfig{Adapters: []types.Adapter{*selection.Adapter}}
		}
		result, callErr := client.AnalyzeDocument(ctx, input)
		if callErr != nil {
			return receipt, callErr
		}
		if result == nil || len(result.Blocks) == 0 {
			return receipt, errors.New("provider response contains no document blocks")
		}
		output, resultCount, modelVersion = result, len(result.Blocks), aws.ToString(result.AnalyzeDocumentModelVersion)
		if result.DocumentMetadata != nil {
			pages = aws.ToInt32(result.DocumentMetadata.Pages)
		}
	}
	if pages != 1 {
		return receipt, fmt.Errorf("synchronous provider response must prove exactly one page; got %d", pages)
	}
	responseBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return receipt, err
	}
	responseBytes = append(responseBytes, '\n')
	if err = os.WriteFile(filepath.Join(stage, "provider-response.json"), responseBytes, 0o600); err != nil {
		return receipt, err
	}
	receipt = Receipt{
		Schema: "elite-aws-textract-analysis-receipt/v2", CreatedAt: time.Now().UTC().Format(time.RFC3339Nano), Provider: "Amazon Textract",
		SDKModule: sdkModule, SDKVersion: sdkVersion, Operation: map[string]string{ModeExpense: "AnalyzeExpense", ModeDocument: "AnalyzeDocument"}[selection.Mode], Region: selection.Region,
		DocumentClass: selection.DocumentClass, ResponseSchemaVersion: selection.ResponseSchemaVersion, DocumentProfileSHA256: selection.ProfileSHA256, SecurityReceiptSHA256: securitySHA, ContentType: mimeType,
		InputFilename: filepath.Base(absInput), InputBytes: len(payload), InputSHA256: payloadSHA, ProviderResponseSHA256: hash(responseBytes), PageCount: pages, ResultCount: resultCount, ModelVersion: modelVersion,
		AdapterID: selection.AdapterID, AdapterVersion: selection.AdapterVersion, AutomaticStorageAuthorized: false,
	}
	receiptBytes, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		return receipt, err
	}
	receiptBytes = append(receiptBytes, '\n')
	if err = os.WriteFile(filepath.Join(stage, "ANALYSIS_RECEIPT.json"), receiptBytes, 0o600); err != nil {
		return receipt, err
	}
	if err = os.Rename(stage, absOutput); err != nil {
		return receipt, err
	}
	committed = true
	return receipt, nil
}

func selectProfile(path, documentClass, expectedRegion string) (approvedSelection, error) {
	var profile documentProfile
	body, err := readStrictJSON(path, "document profile", &profile)
	if err != nil {
		return approvedSelection{}, err
	}
	if profile.Schema != profileSchema || profile.Provider != "Amazon Textract" || profile.SDK != sdkModule+"@"+sdkVersion {
		return approvedSelection{}, errors.New("document profile identity is not admitted")
	}
	if profile.Region == "" || profile.Region == "CONFIGURATION_REQUIRED" || strings.ContainsAny(profile.Region, " /\\") {
		return approvedSelection{}, errors.New("document profile requires an exact AWS region")
	}
	if expectedRegion != "" && expectedRegion != profile.Region {
		return approvedSelection{}, errors.New("runtime region does not match document profile")
	}
	documentClass = strings.TrimSpace(documentClass)
	if documentClass == "" {
		return approvedSelection{}, errors.New("document class is required")
	}
	matches := make([]profileClass, 0, 1)
	for _, item := range profile.Classes {
		if item.ID == documentClass {
			matches = append(matches, item)
		}
	}
	if len(matches) != 1 {
		return approvedSelection{}, errors.New("document class must appear exactly once in profile")
	}
	item := matches[0]
	if item.Decision != "REQUIRED" {
		return approvedSelection{}, fmt.Errorf("document class %s is not REQUIRED", documentClass)
	}
	if item.AutomaticStorage {
		return approvedSelection{}, errors.New("document profile must deny automatic storage")
	}
	if item.ResponseSchemaVersion == "" || item.ResponseSchemaVersion == "CONFIGURATION_REQUIRED" {
		return approvedSelection{}, errors.New("exact response schema version is required")
	}
	mode, features, queries, adapter, err := validateAnalysis(item)
	if err != nil {
		return approvedSelection{}, err
	}
	return approvedSelection{Region: profile.Region, DocumentClass: documentClass, Mode: mode, Features: features, Queries: queries, Adapter: adapter, AdapterID: item.AdapterID, AdapterVersion: item.AdapterVersion, ResponseSchemaVersion: item.ResponseSchemaVersion, ProfileSHA256: hash(body)}, nil
}

func validateAnalysis(item profileClass) (string, []types.FeatureType, []types.Query, *types.Adapter, error) {
	operation := strings.TrimSpace(item.Operation)
	if operation == "AnalyzeExpense" {
		if len(item.Features) > 0 || len(item.Queries) > 0 || item.AdapterID != "" || item.AdapterVersion != "" {
			return "", nil, nil, nil, errors.New("AnalyzeExpense profile cannot contain document features, queries, or adapters")
		}
		return ModeExpense, nil, nil, nil, nil
	}
	if operation != "AnalyzeDocument" {
		return "", nil, nil, nil, errors.New("profile operation must be AnalyzeExpense or AnalyzeDocument")
	}
	if len(item.Features) == 0 {
		return "", nil, nil, nil, errors.New("AnalyzeDocument profile requires at least one feature")
	}
	features := make([]types.FeatureType, 0, len(item.Features))
	seen := map[string]struct{}{}
	for _, raw := range item.Features {
		name := strings.ToUpper(strings.TrimSpace(raw))
		value, ok := featureTypes[name]
		if !ok {
			return "", nil, nil, nil, fmt.Errorf("unsupported feature %q", raw)
		}
		if _, dup := seen[name]; dup {
			return "", nil, nil, nil, fmt.Errorf("duplicate feature %q", name)
		}
		seen[name] = struct{}{}
		features = append(features, value)
	}
	hasQueries := slices.Contains(features, types.FeatureTypeQueries)
	if hasQueries != (len(item.Queries) > 0) {
		return "", nil, nil, nil, errors.New("QUERIES feature and non-empty profile queries must be supplied together")
	}
	if len(item.Queries) > 15 {
		return "", nil, nil, nil, errors.New("synchronous Textract profile exceeds 15 queries per page")
	}
	queries := make([]types.Query, 0, len(item.Queries))
	aliases := map[string]struct{}{}
	for _, query := range item.Queries {
		text, alias := strings.TrimSpace(query.Text), strings.TrimSpace(query.Alias)
		if text == "" || alias == "" {
			return "", nil, nil, nil, errors.New("query text and alias are required")
		}
		if _, dup := aliases[alias]; dup {
			return "", nil, nil, nil, fmt.Errorf("duplicate query alias %q", alias)
		}
		for _, page := range query.Pages {
			if page != "1" {
				return "", nil, nil, nil, errors.New("synchronous profile query pages may only contain page 1")
			}
		}
		aliases[alias] = struct{}{}
		queries = append(queries, types.Query{Text: aws.String(text), Alias: aws.String(alias), Pages: query.Pages})
	}
	if (item.AdapterID == "") != (item.AdapterVersion == "") {
		return "", nil, nil, nil, errors.New("adapter ID and exact version must be supplied together")
	}
	var adapter *types.Adapter
	if item.AdapterID != "" {
		if !hasQueries {
			return "", nil, nil, nil, errors.New("an adapter requires the QUERIES feature")
		}
		if item.AdapterID == "CONFIGURATION_REQUIRED" || item.AdapterVersion == "CONFIGURATION_REQUIRED" {
			return "", nil, nil, nil, errors.New("adapter identity is not configured")
		}
		adapter = &types.Adapter{AdapterId: aws.String(item.AdapterID), Version: aws.String(item.AdapterVersion)}
	}
	return ModeDocument, features, queries, adapter, nil
}

func verifySecurityReceipt(path, payloadSHA string, payloadBytes int, mimeType string) (string, error) {
	var receipt securityReceipt
	body, err := readStrictJSON(path, "security receipt", &receipt)
	if err != nil {
		return "", err
	}
	if receipt.Schema != securityReceiptSchema {
		return "", errors.New("unsupported security receipt schema")
	}
	if receipt.Decision != "ADMITTED" || receipt.Reason != "ALL_SELECTED_GATES_PASSED" {
		return "", errors.New("security receipt did not admit the input")
	}
	if receipt.BusinessStorageAuthorized == nil || *receipt.BusinessStorageAuthorized {
		return "", errors.New("security receipt must explicitly deny business storage")
	}
	if receipt.ApprovalID == "" || !isSHA256(receipt.PolicySHA256) {
		return "", errors.New("security receipt approval or policy evidence is invalid")
	}
	if _, e := time.Parse(time.RFC3339Nano, receipt.StartedAt); e != nil {
		return "", errors.New("security receipt started_at is invalid")
	}
	if _, e := time.Parse(time.RFC3339Nano, receipt.CompletedAt); e != nil {
		return "", errors.New("security receipt completed_at is invalid")
	}
	if receipt.SecurityClaim == "" {
		return "", errors.New("security receipt claim is missing")
	}
	if receipt.Input.SHA256 != payloadSHA || receipt.Input.Bytes != payloadBytes {
		return "", errors.New("security receipt does not match input bytes")
	}
	if receipt.ContentType.MIMEType != mimeType || receipt.ContentType.Label == "" || receipt.ContentType.Score < 0 || receipt.ContentType.Score > 1 {
		return "", errors.New("security receipt content type does not match input")
	}
	for _, name := range []string{"clamscan", "magika", "yara_x"} {
		tool, ok := receipt.Tools[name]
		if !ok || !isSHA256(tool.SHA256) || tool.Version == "" {
			return "", fmt.Errorf("security receipt %s evidence is invalid", name)
		}
	}
	if receipt.ClamAV.ExitCode != 0 || !isSHA256(receipt.ClamAV.OutputSHA256) {
		return "", errors.New("security receipt ClamAV result is invalid")
	}
	if !isSHA256(receipt.YaraX.CompiledRulesSHA256) || receipt.YaraX.MatchingRules == nil || len(receipt.YaraX.MatchingRules) != 0 {
		return "", errors.New("security receipt YARA-X result is invalid")
	}
	return hash(body), nil
}

func readStrictJSON(path, label string, target any) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("%s must be a regular non-symlink file", label)
	}
	if info.Size() < 2 || info.Size() > 1024*1024 {
		return nil, fmt.Errorf("%s size is invalid", label)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(target); err != nil {
		return nil, fmt.Errorf("invalid %s: %w", label, err)
	}
	var trailing any
	if err = decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("invalid %s: trailing JSON value", label)
	}
	return body, nil
}

func isSHA256(value string) bool {
	if len(value) != 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil && value == strings.ToLower(value)
}
func hash(value []byte) string { sum := sha256.Sum256(value); return hex.EncodeToString(sum[:]) }
