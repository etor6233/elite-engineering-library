package aifoundation

import (
	"strings"
	"testing"
)

func TestNXUnsupportedSchemaCannotBeSilentlyAccepted(t *testing.T) {
	for _, property := range []string{`{"type":"string","enum":["yes"]}`, `{"type":"integer","minimum":1}`, `{"type":"string","pattern":"^x$"}`, `{"type":"array","items":{"type":"integer"}}`, `{"type":"object","properties":{"n":{"type":"integer"}},"required":["n"],"additionalProperties":false}`} {
		schema := []byte(`{"type":"object","properties":{"v":` + property + `},"additionalProperties":false}`)
		for _, value := range []string{`"wrong"`, `0`, `["wrong"]`, `{"n":"wrong"}`} {
			if err := ValidateStructuredOutput(schema, []byte(`{"v":`+value+`}`)); err == nil {
				t.Errorf("silently accepted %s / %s", property, value)
			}
		}
	}
}
func TestNXDuplicateJSONKeysRejected(t *testing.T) {
	schema := []byte(`{"type":"object","properties":{"v":{"type":"string"}},"additionalProperties":false}`)
	if ValidateStructuredOutput(schema, []byte(`{"v":"one","v":"two"}`)) == nil {
		t.Error("duplicate instance key accepted")
	}
	if ValidateStructuredOutput([]byte(`{"type":"object","type":"object","properties":{},"additionalProperties":false}`), []byte(`{}`)) == nil {
		t.Error("duplicate schema key accepted")
	}
}
func TestNXSchemaRequiredAndKeywordValidation(t *testing.T) {
	for _, s := range []string{`{"type":"object","properties":{},"required":["ghost"],"additionalProperties":false}`, `{"type":"object","properties":{"v":{"type":"string"}},"required":["v","v"],"additionalProperties":false}`, `{"type":"object","properties":{},"additionalProperties":null}`, `{"Type":"object","properties":{},"additionalProperties":false}`, `{"type":"object","properties":{},"additionalProperties":false,"minimum":1}`} {
		if _, err := parseClosedSchema([]byte(s)); err == nil {
			t.Errorf("invalid schema accepted: %s", s)
		}
	}
}
func TestNXOutputByteLimit(t *testing.T) {
	s := []byte(`{"type":"object","properties":{"v":{"type":"string"}},"additionalProperties":false}`)
	if ValidateStructuredOutput(s, []byte(`{"v":"`+strings.Repeat("x", 1048577)+`"}`)) == nil {
		t.Error("unbounded output")
	}
}
