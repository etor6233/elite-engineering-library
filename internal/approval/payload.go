package approval

// AUTHORED serialization glue for durable approval of exact payloads.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"strings"
	"unicode/utf8"
)

func CanonicalPayload(raw []byte) (json.RawMessage, string, error) {
	if len(raw) == 0 || len(raw) > 32768 || !utf8.Valid(raw) {
		return nil, "", ErrInvalidRequest
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	if payloadKeys(d, 0) != nil {
		return nil, "", ErrInvalidRequest
	}
	if _, e := d.Token(); e != io.EOF {
		return nil, "", ErrInvalidRequest
	}
	d = json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	var value map[string]any
	if d.Decode(&value) != nil || value == nil {
		return nil, "", ErrInvalidRequest
	}
	canonical, e := json.Marshal(value)
	if e != nil || len(canonical) > 32768 {
		return nil, "", ErrInvalidRequest
	}
	h := sha256.Sum256(canonical)
	return canonical, hex.EncodeToString(h[:]), nil
}
func payloadKeys(d *json.Decoder, depth int) error {
	if depth > 16 {
		return ErrInvalidRequest
	}
	t, e := d.Token()
	if e != nil {
		return e
	}
	v, ok := t.(json.Delim)
	if !ok {
		return nil
	}
	switch v {
	case '{':
		seen := map[string]bool{}
		for d.More() {
			k, e := d.Token()
			if e != nil {
				return e
			}
			s, ok := k.(string)
			if !ok {
				return ErrInvalidRequest
			}
			s = strings.ToLower(s)
			if seen[s] {
				return ErrInvalidRequest
			}
			seen[s] = true
			if payloadKeys(d, depth+1) != nil {
				return ErrInvalidRequest
			}
		}
	case '[':
		for d.More() {
			if payloadKeys(d, depth+1) != nil {
				return ErrInvalidRequest
			}
		}
	default:
		return ErrInvalidRequest
	}
	_, e = d.Token()
	return e
}
