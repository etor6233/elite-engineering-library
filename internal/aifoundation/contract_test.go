package aifoundation

import "testing"

const closedSchemaJSON = `{
  "type": "object",
  "properties": {
    "answer": {"type": "string"},
    "count": {"type": "integer"},
    "score": {"type": "number"},
    "ok": {"type": "boolean"}
  },
  "required": ["answer"],
  "additionalProperties": false
}`

func TestValidateStructuredOutputValid(t *testing.T) {
	raw := []byte(`{"answer":"hi","count":3,"score":0.5,"ok":true}`)
	if err := ValidateStructuredOutput([]byte(closedSchemaJSON), raw); err != nil {
		t.Fatalf("valid output rejected: %v", err)
	}
}

func TestValidateStructuredOutputMissingRequired(t *testing.T) {
	raw := []byte(`{"count":1}`)
	if err := ValidateStructuredOutput([]byte(closedSchemaJSON), raw); err == nil {
		t.Fatal("missing required key accepted")
	}
}

func TestValidateStructuredOutputUnknownKey(t *testing.T) {
	raw := []byte(`{"answer":"hi","extra":"x"}`)
	if err := ValidateStructuredOutput([]byte(closedSchemaJSON), raw); err == nil {
		t.Fatal("unknown key accepted")
	}
}

func TestValidateStructuredOutputWrongType(t *testing.T) {
	raw := []byte(`{"answer":123}`)
	if err := ValidateStructuredOutput([]byte(closedSchemaJSON), raw); err == nil {
		t.Fatal("wrong type accepted")
	}
}

func TestValidateStructuredOutputOpenSchemaRejected(t *testing.T) {
	open := []byte(`{"type":"object","properties":{"answer":{"type":"string"}}}`)
	raw := []byte(`{"answer":"hi"}`)
	if err := ValidateStructuredOutput(open, raw); err == nil {
		t.Fatal("open schema (additionalProperties not false) accepted")
	}
}

func TestValidateStructuredOutputNonObject(t *testing.T) {
	raw := []byte(`[1,2,3]`)
	if err := ValidateStructuredOutput([]byte(closedSchemaJSON), raw); err == nil {
		t.Fatal("array root accepted")
	}
}
