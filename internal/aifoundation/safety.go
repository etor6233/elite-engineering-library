package aifoundation

import "strings"

// SafetyDecision is the outcome of a safety gate.
type SafetyDecision string

const (
	SafetyAllow SafetyDecision = "allow"
	SafetyBlock SafetyDecision = "block"
)

// SafetyInput carries the untrusted prompt context.
type SafetyInput struct {
	UserText    string
	ContainsPII bool
	Allowlisted bool // an explicit, versioned allowlist decision
}

// SafetyOutput carries the model response before it reaches the caller.
type SafetyOutput struct {
	Content     string
	Refusal     bool
	Allowlisted bool
}

// SafetyGate decides, fail-closed, whether a request or response may proceed.
type SafetyGate interface {
	PreFlight(SafetyInput) SafetyDecision
	PostFlight(SafetyOutput) SafetyDecision
}

// HeuristicSafetyGate is an AUTHORED, deterministic safety gate. Its detectors
// are heuristics, not a guarantee; production must add an admitted content
// classifier and PII recognizer. It fails closed: the zero value blocks all
// traffic until an explicit policy is configured.
type HeuristicSafetyGate struct {
	enabled              bool
	blockedInputPhrases  []string
	blockedOutputPhrases []string
}

// NewSafetyGate returns a gate with a minimal, explicit allowlist.
func NewSafetyGate() HeuristicSafetyGate {
	return HeuristicSafetyGate{
		enabled: true,
		blockedInputPhrases: []string{
			"ignore previous instructions",
			"ignore all prior instructions",
			"disregard your instructions",
			"reveal your system prompt",
		},
		blockedOutputPhrases: []string{
			"system prompt:",
		},
	}
}

// PreFlight fails closed: PII without an allowlist or any injection signal
// blocks; the zero value blocks everything.
func (g HeuristicSafetyGate) PreFlight(in SafetyInput) SafetyDecision {
	if !g.enabled {
		return SafetyBlock
	}
	lower := strings.ToLower(in.UserText)
	for _, p := range g.blockedInputPhrases {
		if strings.Contains(lower, p) {
			return SafetyBlock
		}
	}
	if in.ContainsPII && !in.Allowlisted {
		return SafetyBlock
	}
	return SafetyAllow
}

// PostFlight fails closed: an explicit refusal is a safe final answer, but a
// blocked output phrase blocks; the zero value blocks everything.
func (g HeuristicSafetyGate) PostFlight(out SafetyOutput) SafetyDecision {
	if !g.enabled {
		return SafetyBlock
	}
	lower := strings.ToLower(out.Content)
	for _, p := range g.blockedOutputPhrases {
		if strings.Contains(lower, p) {
			return SafetyBlock
		}
	}
	if out.Refusal {
		return SafetyAllow
	}
	return SafetyAllow
}
