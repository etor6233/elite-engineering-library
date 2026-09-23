package economy

import "errors"

// CompressedPrompt is a bounded, prioritized prompt.
type CompressedPrompt struct {
	System  string
	Context []string // highest-priority context, truncated to the budget
	User    string
	Err     error // mandatory input did not fit: no runnable prompt is returned
	Tokens  int
}

// Compress builds a prompt within a token budget. The system message and the
// user's last message are always kept; context items are added only while the
// budget allows. This minimizes tokens for the eventual model call.
func Compress(system, user string, context []string, budgetTokens int) CompressedPrompt {
	used := EstimateTokens(system) + EstimateTokens(user)
	if budgetTokens <= 0 || used > budgetTokens {
		return CompressedPrompt{Err: errors.New("economy: mandatory prompt exceeds estimate budget")}
	}
	out := CompressedPrompt{System: system, User: user}
	for _, ctx := range context {
		t := EstimateTokens(ctx)
		if used+t > budgetTokens {
			break
		}
		out.Context = append(out.Context, ctx)
		used += t
	}
	out.Tokens = used
	return out
}
