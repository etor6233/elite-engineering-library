package documentbridge

// AUTHORED fixture composition. Typed fields are reviewed document metadata,
// not OCR accuracy, a payment settlement, a legal conclusion or domain posting.
import (
	"bytes"
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"regexp"
	"time"
	"unicode"
)

//go:embed typed-fixtures.json
var typedCatalogBytes []byte

type FieldDefinition struct {
	Key       string `json:"key"`
	Label     string `json:"label_es"`
	Type      string `json:"type"`
	MaxLength int    `json:"max_length"`
}
type ClassDefinition struct {
	ID               string            `json:"class_id"`
	SchemaVersion    string            `json:"schema_version"`
	Label            string            `json:"label_es"`
	Task             string            `json:"task_id"`
	Fields           []FieldDefinition `json:"fields"`
	Method           string            `json:"extraction_method"`
	AutomaticStorage bool              `json:"automatic_storage_authorized"`
	FixtureSHA       string            `json:"fixture_sha256"`
	FixtureBase64    string            `json:"fixture_base64,omitempty"`
	Suggested        map[string]string `json:"suggested,omitempty"`
}
type ClassCatalog struct {
	Schema     string            `json:"schema"`
	Mode       string            `json:"mode"`
	ProfileSHA string            `json:"profile_sha256"`
	Classes    []ClassDefinition `json:"classes"`
}

func typedCatalog() ClassCatalog {
	var c ClassCatalog
	if json.Unmarshal(typedCatalogBytes, &c) != nil || len(c.Classes) != 21 {
		panic("invalid embedded typed fixture catalog")
	}
	return c
}
func TypedProfileSHA() string { return typedCatalog().ProfileSHA }
func ClassByID(id, version string) (ClassDefinition, error) {
	for _, c := range typedCatalog().Classes {
		if c.ID == id && c.SchemaVersion == version {
			return c, nil
		}
	}
	return ClassDefinition{}, ErrContract
}
func Catalog(mode, profile string) ClassCatalog {
	c := typedCatalog()
	c.Mode = mode
	c.ProfileSHA = profile
	if mode != "TYPED_FIXTURE" {
		c.Classes = c.Classes[:1]
		c.Classes[0].Method = "EXISTING_AWS_EXPENSE_REFERENCE"
		c.Classes[0].FixtureSHA = ""
		if mode == "FIXTURE" {
			c.Classes[0].FixtureSHA = PublicFixtureSHA
		}
		c.Classes[0].Fields[0].MaxLength = 128
		c.Classes[0].Fields[2].Type = "text"
		c.Classes[0].Fields[2].MaxLength = 64
		c.Classes[0].Fields[3].MaxLength = 3
	}
	for i := range c.Classes {
		c.Classes[i].FixtureBase64 = ""
		c.Classes[i].Suggested = nil
	}
	return c
}
func TypedFixture(id string) (Original, error) {
	c, e := ClassByID(id, "1")
	if e != nil {
		return Original{}, e
	}
	b, e := base64.StdEncoding.DecodeString(c.FixtureBase64)
	if e != nil || Hash(b) != c.FixtureSHA {
		return Original{}, ErrContract
	}
	return Original{Name: id + ".pdf", SHA256: c.FixtureSHA, Bytes: b, ClassID: id, SchemaVersion: "1"}, nil
}

var decimalPattern = regexp.MustCompile(`^(0|[1-9][0-9]{0,17})(\.[0-9]{1,8})?$`)
var integerPattern = regexp.MustCompile(`^(0|[1-9][0-9]{0,8})$`)
var currencyPattern = regexp.MustCompile(`^[A-Z]{3}$`)

func ValidateTypedFields(id, version string, raw []byte) error {
	if len(raw) == 0 || len(raw) > 32768 {
		return ErrContract
	}
	c, e := ClassByID(id, version)
	if e != nil {
		return e
	}
	var fields map[string]string
	d := json.NewDecoder(bytes.NewReader(raw))
	if d.Decode(&fields) != nil || len(fields) != len(c.Fields) {
		return ErrContract
	}
	for _, def := range c.Fields {
		v, ok := fields[def.Key]
		if !ok || !Text(v, def.MaxLength) {
			return ErrContract
		}
		for _, r := range v {
			if unicode.IsControl(r) {
				return ErrContract
			}
		}
		switch def.Type {
		case "date":
			date, e := time.Parse("2006-01-02", v)
			if e != nil || date.Year() < 1 || date.Format("2006-01-02") != v {
				return ErrContract
			}
		case "decimal":
			if !decimalPattern.MatchString(v) {
				return ErrContract
			}
		case "integer":
			if !integerPattern.MatchString(v) {
				return ErrContract
			}
		case "currency":
			if !currencyPattern.MatchString(v) {
				return ErrContract
			}
		case "text":
		default:
			return ErrContract
		}
	}
	return nil
}
func validTypedOriginal(o Original) error {
	c, e := ClassByID(o.ClassID, o.SchemaVersion)
	if e != nil || o.SHA256 != c.FixtureSHA {
		return ErrContract
	}
	return nil
}

type TypedPipeline struct{ profile string }

func NewTypedFixturePipeline(profile, sha string) (*TypedPipeline, error) {
	b, e := regular(profile, 32768)
	if e != nil || Hash(b) != sha || sha != TypedProfileSHA() {
		return nil, ErrContract
	}
	return &TypedPipeline{sha}, nil
}

type TypedReceipt struct {
	Schema           string `json:"schema"`
	Mode             string `json:"mode"`
	ClassID          string `json:"class_id"`
	SchemaVersion    string `json:"schema_version"`
	InputSHA         string `json:"input_sha256"`
	InputBytes       int    `json:"input_bytes"`
	ProfileSHA       string `json:"profile_sha256"`
	SecuritySHA      string `json:"security_sha256"`
	ProviderSHA      string `json:"provider_sha256"`
	FieldsSHA        string `json:"fields_sha256"`
	Method           string `json:"method"`
	AutomaticStorage bool   `json:"automatic_storage_authorized"`
	OCRExecuted      bool   `json:"ocr_executed"`
}

func (p *TypedPipeline) Run(ctx context.Context, o Original) (Evidence, error) {
	if ctx.Err() != nil {
		return Evidence{}, ctx.Err()
	}
	if o.Validate("TYPED_FIXTURE") != nil {
		return Evidence{}, ErrContract
	}
	c, _ := ClassByID(o.ClassID, o.SchemaVersion)
	fields, _ := json.Marshal(c.Suggested)
	security, _ := json.Marshal(map[string]any{"schema": "typed-fixture-allowlist/v1", "input_sha256": o.SHA256, "mode": "TYPED_FIXTURE", "method": "EXACT_SYNTHETIC_BYTE_ALLOWLIST", "malware_scanner_executed": false, "private_file_admission": false})
	provider, _ := json.Marshal(map[string]any{"schema": "typed-fixture-proposal/v1", "input_sha256": o.SHA256, "class_id": o.ClassID, "schema_version": o.SchemaVersion, "fields": c.Suggested, "ocr_executed": false})
	receipt, _ := json.Marshal(TypedReceipt{"document-typed-analysis/v1", "TYPED_FIXTURE", o.ClassID, o.SchemaVersion, o.SHA256, len(o.Bytes), p.profile, Hash(security), Hash(provider), Hash(fields), "HASH_LOCKED_TYPED_FIXTURE", false, false})
	return Evidence{Security: security, Provider: provider, Receipt: receipt, TypedFields: fields, Mode: "TYPED_FIXTURE"}, nil
}
func ValidateTypedEvidence(e Evidence, o Original, profile string) error {
	if len(e.Security) == 0 || len(e.Security) > 65536 || len(e.Provider) == 0 || len(e.Provider) > 4194304 || len(e.Receipt) > 32768 || e.Mode != "TYPED_FIXTURE" || o.Validate(e.Mode) != nil || ValidateTypedFields(o.ClassID, o.SchemaVersion, e.TypedFields) != nil {
		return ErrContract
	}
	var r TypedReceipt
	if json.Unmarshal(e.Receipt, &r) != nil || r.Schema != "document-typed-analysis/v1" || r.Mode != e.Mode || r.ClassID != o.ClassID || r.SchemaVersion != o.SchemaVersion || r.InputSHA != o.SHA256 || r.InputBytes != len(o.Bytes) || r.ProfileSHA != profile || profile != TypedProfileSHA() || r.SecuritySHA != Hash(e.Security) || r.ProviderSHA != Hash(e.Provider) || r.FieldsSHA != Hash(e.TypedFields) || r.Method != "HASH_LOCKED_TYPED_FIXTURE" || r.AutomaticStorage || r.OCRExecuted {
		return ErrContract
	}
	return nil
}
