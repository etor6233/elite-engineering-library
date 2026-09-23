// Adapted from PostHog/posthog, MIT, commit 6fafbb9081bd15e79448af5650e02a4f9ea435cc.
// Original: frontend/src/scenes/surveys/utils.ts, calculateNPSFromRawData.
// See docs/posthog-nps.md for the exact adaptation boundary and preserved notices.
package posthognps

// Breakdown uses the existing database's nonnegative int64 grouped counts.
type Breakdown struct {
	Total, Promoters, Passives, Detractors int64
	Score                                  float64
}

// Calculate translates the upstream aggregation after SQL has grouped 0..6,
// 7..8, and 9..10. Caller must validate total >= promoters + detractors >= 0.
// Presentation rounding is deliberately left to the consumer: the existing API
// returns an unrounded number, unlike upstream's one-decimal display string.
func Calculate(total, promoters, detractors int64) Breakdown {
	if total == 0 {
		return Breakdown{}
	}
	return Breakdown{
		Total: total, Promoters: promoters, Passives: total - promoters - detractors,
		Detractors: detractors,
		Score:      float64(promoters-detractors) / float64(total) * 100,
	}
}
