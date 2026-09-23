package main

import (
	"bytes"
	"testing"
)

func TestRetentionRequiresScopeAndBound(t *testing.T) {
	for _, args := range [][]string{nil, {"-tenant", "tenant"}, {"-tenant", "tenant", "-organization", "org"}, {"-tenant", "tenant", "-organization", "org", "-limit", "1001"}, {"-tenant", "tenant", "-organization", "org", "-limit", "1", "extra"}} {
		var out, err bytes.Buffer
		code := run(args, func(string) string { t.Fatal("database consulted before scope validation"); return "" }, &out, &err)
		if code != 2 || out.Len() != 0 {
			t.Fatalf("invalid invocation code=%d output=%s", code, out.String())
		}
	}
}
