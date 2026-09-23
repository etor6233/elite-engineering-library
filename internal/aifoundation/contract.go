package aifoundation

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"unicode/utf8"
)

// LLMProvider is the provider-agnostic generation boundary. Implementations
// (OpenAI, Anthropic, self-hosted vLLM, ...) are separate, admitted adapters.
type LLMProvider interface {
	Generate(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
}

// Message is a single role-tagged conversation turn.
type Message struct {
	Role    string // system | user | assistant | tool
	Content string
}

// CompletionRequest carries an exact model pin and a bounded generation budget.
type CompletionRequest struct {
	Model           ModelRef
	Messages        []Message
	Schema          json.RawMessage // closed object schema enforced on the response
	MaxOutputTokens int
	Temperature     *float64
	Seed            *int64
}

// CompletionResponse is the raw, unvalidated provider payload.
type CompletionResponse struct {
	Content      json.RawMessage
	FinishReason string
}

var (
	ErrInvalidSchema = errors.New("aifoundation: invalid schema")
	ErrInvalidOutput = errors.New("aifoundation: invalid structured output")
)

// This owner supports a bounded, closed subset; other JSON Schema keywords are
// rejected, never treated as if they had been enforced. Business owners still
// authorize and validate all effects independently.
type closedSchema struct {
	kind       string
	properties map[string]closedSchema
	required   []string
	items      *closedSchema
}

const maxStructuredBytes = 1 << 20

// strictJSON rejects duplicate keys, trailing values and excessive nesting.
func strictJSON(raw []byte) error {
	if len(raw) == 0 || len(raw) > maxStructuredBytes || !utf8.Valid(raw) {
		return errors.New("JSON size/encoding")
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	nodes := 0
	var walk func(int) error
	walk = func(depth int) error {
		nodes++
		if depth > 24 || nodes > 16384 {
			return errors.New("JSON complexity")
		}
		t, err := dec.Token()
		if err != nil {
			return err
		}
		if d, ok := t.(json.Delim); ok {
			switch d {
			case '{':
				seen := map[string]bool{}
				for dec.More() {
					k, err := dec.Token()
					if err != nil {
						return err
					}
					key, ok := k.(string)
					if !ok || seen[key] {
						return errors.New("duplicate/invalid key")
					}
					seen[key] = true
					if err = walk(depth + 1); err != nil {
						return err
					}
				}
				close, err := dec.Token()
				if err != nil || close != json.Delim('}') {
					return errors.New("object close")
				}
			case '[':
				for dec.More() {
					if err := walk(depth + 1); err != nil {
						return err
					}
				}
				close, err := dec.Token()
				if err != nil || close != json.Delim(']') {
					return errors.New("array close")
				}
			default:
				return errors.New("unexpected delimiter")
			}
		}
		return nil
	}
	if err := walk(0); err != nil {
		return err
	}
	if _, err := dec.Token(); err != io.EOF {
		return errors.New("trailing JSON")
	}
	return nil
}
func parseClosedSchema(raw []byte) (closedSchema, error) {
	if err := strictJSON(raw); err != nil {
		return closedSchema{}, fmt.Errorf("%w: %v", ErrInvalidSchema, err)
	}
	s, err := parseSchemaNode(raw, 0)
	if err != nil {
		return s, err
	}
	if s.kind != "object" {
		return s, fmt.Errorf("%w: root must be object", ErrInvalidSchema)
	}
	return s, nil
}
func parseSchemaNode(raw []byte, depth int) (closedSchema, error) {
	bad := func() (closedSchema, error) { return closedSchema{}, ErrInvalidSchema }
	if depth > 12 {
		return bad()
	}
	var obj map[string]json.RawMessage
	if json.Unmarshal(raw, &obj) != nil || obj == nil {
		return bad()
	}
	var s closedSchema
	if json.Unmarshal(obj["type"], &s.kind) != nil {
		return bad()
	}
	allowed := map[string]bool{"type": true, "description": true}
	switch s.kind {
	case "object":
		allowed["properties"] = true
		allowed["required"] = true
		allowed["additionalProperties"] = true
	case "array":
		allowed["items"] = true
	case "string", "number", "integer", "boolean", "null":
	default:
		return bad()
	}
	for k := range obj {
		if !allowed[k] {
			return closedSchema{}, fmt.Errorf("%w: unsupported keyword %q", ErrInvalidSchema, k)
		}
	}
	if desc, ok := obj["description"]; ok {
		var v string
		if json.Unmarshal(desc, &v) != nil || string(desc) == "null" {
			return bad()
		}
	}
	if s.kind == "object" {
		if string(bytes.TrimSpace(obj["additionalProperties"])) != "false" {
			return bad()
		}
		var props map[string]json.RawMessage
		if json.Unmarshal(obj["properties"], &props) != nil || props == nil || len(props) > 256 {
			return bad()
		}
		s.properties = map[string]closedSchema{}
		for k, v := range props {
			if len(k) == 0 || len(k) > 256 {
				return bad()
			}
			child, err := parseSchemaNode(v, depth+1)
			if err != nil {
				return s, err
			}
			s.properties[k] = child
		}
		if req, ok := obj["required"]; ok {
			if string(req) == "null" || json.Unmarshal(req, &s.required) != nil {
				return bad()
			}
		}
		seen := map[string]bool{}
		for _, k := range s.required {
			if _, ok := s.properties[k]; !ok || seen[k] {
				return bad()
			}
			seen[k] = true
		}
	}
	if s.kind == "array" {
		child, err := parseSchemaNode(obj["items"], depth+1)
		if err != nil {
			return s, err
		}
		s.items = &child
	}
	return s, nil
}

// ValidateStructuredSchema checks the entire supported schema before a provider call.
func ValidateStructuredSchema(raw []byte) error { _, err := parseClosedSchema(raw); return err }
func (s closedSchema) validate(raw json.RawMessage) error {
	if s.kind == "object" {
		var v map[string]json.RawMessage
		if json.Unmarshal(raw, &v) != nil || v == nil {
			return ErrInvalidOutput
		}
		for k := range v {
			if _, ok := s.properties[k]; !ok {
				return ErrInvalidOutput
			}
		}
		for _, k := range s.required {
			if _, ok := v[k]; !ok {
				return ErrInvalidOutput
			}
		}
		for k, value := range v {
			if err := s.properties[k].validate(value); err != nil {
				return err
			}
		}
		return nil
	}
	if s.kind == "array" {
		var a []json.RawMessage
		if json.Unmarshal(raw, &a) != nil || a == nil || s.items == nil {
			return ErrInvalidOutput
		}
		for _, v := range a {
			if err := s.items.validate(v); err != nil {
				return err
			}
		}
		return nil
	}
	var v any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		return ErrInvalidOutput
	}
	switch s.kind {
	case "string":
		if _, ok := v.(string); ok {
			return nil
		}
	case "boolean":
		if _, ok := v.(bool); ok {
			return nil
		}
	case "null":
		if v == nil {
			return nil
		}
	case "number":
		if _, ok := v.(json.Number); ok {
			return nil
		}
	case "integer":
		if n, ok := v.(json.Number); ok {
			if len(n) > 128 {
				return ErrInvalidOutput
			}
			f, err := n.Float64()
			if err != nil || math.IsInf(f, 0) {
				return ErrInvalidOutput
			}
			r, ok := new(big.Rat).SetString(string(n))
			if ok && r.IsInt() {
				return nil
			}
		}
	}
	return ErrInvalidOutput
}
func ValidateStructuredOutput(schema, raw []byte) error {
	s, err := parseClosedSchema(schema)
	if err != nil {
		return err
	}
	if err := strictJSON(raw); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidOutput, err)
	}
	return s.validate(raw)
}
