package randomid

import (
	"regexp"
	"testing"
)

var uuidV4Pattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestGeneratorProducesUniqueVersion4Identifiers(t *testing.T) {
	generator := Generator{}
	seen := make(map[string]struct{}, 256)
	for range 256 {
		value := generator.New()
		if !uuidV4Pattern.MatchString(value) {
			t.Fatalf("invalid UUID v4: %q", value)
		}
		if _, exists := seen[value]; exists {
			t.Fatalf("duplicate UUID v4: %q", value)
		}
		seen[value] = struct{}{}
	}
}
