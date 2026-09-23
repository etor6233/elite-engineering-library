package aifoundation

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// Assertion judges a response against a versioned expectation.
type Assertion func(CompletionResponse) bool

// EvalCase is a single golden test.
type EvalCase struct {
	Required bool
	ID       string
	Input    CompletionRequest
	Assert   Assertion
}

// EvalSuite is a versioned golden set plus a release threshold.
type EvalSuite struct {
	ID          string
	MinPassRate float64 // 0..1
	Cases       []EvalCase
}

// EvalResult summarizes a run and its release decision.
type EvalCaseResult struct {
	ID       string
	Required bool
	Passed   bool
}

type EvalResult struct {
	Cases      []EvalCaseResult
	Total      int
	Passed     int
	PassRate   float64
	GatePassed bool
}

// Run executes every case and computes the release decision.
func (s EvalSuite) Run(p LLMProvider) (EvalResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.RunContext(ctx, p)
}

func (s EvalSuite) RunContext(ctx context.Context, p LLMProvider) (EvalResult, error) {
	if ctx == nil || strings.TrimSpace(s.ID) == "" {
		return EvalResult{}, errors.New("aifoundation: eval identity/context required")
	}
	if p == nil {
		return EvalResult{}, errors.New("aifoundation: nil provider")
	}
	if len(s.Cases) == 0 {
		return EvalResult{}, errors.New("aifoundation: empty eval suite")
	}
	if math.IsNaN(s.MinPassRate) || math.IsInf(s.MinPassRate, 0) || s.MinPassRate < 0 || s.MinPassRate > 1 {
		return EvalResult{}, errors.New("aifoundation: min pass rate out of range")
	}
	seen := map[string]bool{}
	for _, c := range s.Cases {
		if c.Assert == nil || strings.TrimSpace(c.ID) == "" || seen[c.ID] {
			return EvalResult{}, errors.New("aifoundation: unique eval identity and assertion required")
		}
		seen[c.ID] = true
	}
	var res EvalResult
	requiredPassed := true
	res.Total = len(s.Cases)
	for _, c := range s.Cases {
		if err := ctx.Err(); err != nil {
			return res, err
		}
		resp, err := p.Generate(ctx, c.Input)
		if err != nil {
			return res, fmt.Errorf("aifoundation: case %s: %w", c.ID, err)
		}
		if err := ctx.Err(); err != nil {
			return res, err
		}
		passed := c.Assert(resp)
		res.Cases = append(res.Cases, EvalCaseResult{ID: c.ID, Required: c.Required, Passed: passed})
		if c.Required && !passed {
			requiredPassed = false
		}
		if passed {
			res.Passed++
		}
	}
	res.PassRate = float64(res.Passed) / float64(res.Total)
	res.GatePassed = requiredPassed && res.PassRate >= s.MinPassRate
	return res, nil
}
