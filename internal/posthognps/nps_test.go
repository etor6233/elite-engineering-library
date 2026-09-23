// Adapted selected PostHog MIT calculateNpsBreakdown examples; see source-lock.
// Invalid/missing histogram parsing cases belong to the upstream TS harness,
// because the local caller receives typed, validated PostgreSQL counts.
package posthognps

import (
	"fmt"
	"testing"
)

func TestUpstreamBreakdownExamples(t *testing.T) {
	cases := []struct {
		name        string
		total, p, d int64
		want        Breakdown
		display     string
	}{
		{"zero total", 0, 0, 0, Breakdown{}, "0.0"},
		{"mixed", 17, 6, 7, Breakdown{17, 6, 4, 7, -100.0 / 17}, "-5.9"},
		{"all zero bins", 0, 0, 0, Breakdown{}, "0.0"},
		{"only promoters", 10, 10, 0, Breakdown{10, 10, 0, 0, 100}, "100.0"},
		{"only passives", 10, 0, 0, Breakdown{10, 0, 10, 0, 0}, "0.0"},
		{"only detractors", 14, 0, 14, Breakdown{14, 0, 0, 14, -100}, "-100.0"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := Calculate(c.total, c.p, c.d)
			if g.Total != c.want.Total || g.Promoters != c.want.Promoters || g.Passives != c.want.Passives || g.Detractors != c.want.Detractors || fmt.Sprintf("%.1f", g.Score) != c.display {
				t.Fatalf("got %+v want %+v display %s", g, c.want, c.display)
			}
		})
	}
}
