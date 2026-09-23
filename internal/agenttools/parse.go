package agenttools

import (
	"regexp"
	"strings"
)

var kvRe = regexp.MustCompile(`(?i)([a-záéíóúñ]+):\s*([^,\n;]+)`)

// parseKV extracts deterministic "key: value" pairs (case-insensitive keys).
// This is the deterministic path; converting free text to structured input is
// the LLM's job (CONDITIONED).
func parseKV(text string) map[string]string {
	out := map[string]string{}
	for _, m := range kvRe.FindAllStringSubmatch(text, -1) {
		k := strings.ToLower(m[1])
		v := strings.TrimSpace(m[2])
		if v != "" {
			out[k] = v
		}
	}
	return out
}
