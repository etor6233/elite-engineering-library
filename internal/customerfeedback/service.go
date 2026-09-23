// AUTHORED customer feedback boundary. Policy texts and survey activation are project-owned.
package customerfeedback

import (
	"context"
	"elite.local/enterprise/internal/posthognps"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalid     = errors.New("invalid survey request")
	ErrUnavailable = errors.New("survey unavailable")
	ErrConflict    = errors.New("survey response conflicts with stored response")
	ErrNotFound    = errors.New("survey response not found")
)

type Scope struct{ Tenant, Organization, Customer string }
type Definition struct {
	ID             string `json:"id"`
	Prompt         string `json:"prompt"`
	ConsentVersion string `json:"consent_version"`
	ConsentNotice  string `json:"consent_notice"`
	Accepting      bool   `json:"accepting"`
}
type Answer struct {
	Score          int       `json:"score"`
	ConsentVersion string    `json:"consent_version"`
	ReceivedAt     time.Time `json:"received_at"`
}
type Submission struct {
	Score          int    `json:"score"`
	Consent        bool   `json:"consent"`
	ConsentVersion string `json:"consent_version"`
}
type Result struct {
	Answer Answer `json:"answer"`
	Replay bool   `json:"replay"`
}
type Summary struct {
	Responses int64    `json:"responses"`
	Available bool     `json:"available"`
	NPS       *float64 `json:"nps"`
}
type Repository interface {
	Definition(context.Context, Scope, string) (Definition, error)
	Submit(context.Context, Scope, string, Submission) (Result, error)
	OwnAnswer(context.Context, Scope, string) (Answer, error)
	Summary(context.Context, Scope, string) (Summary, error)
}
type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }
func valid(value string) bool {
	return strings.TrimSpace(value) == value && value != "" && len(value) <= 128 && !strings.ContainsAny(value, "\x00\r\n")
}
func validate(scope Scope, id string, customer bool) error {
	if !valid(scope.Tenant) || !valid(scope.Organization) || !valid(id) || (customer && !valid(scope.Customer)) {
		return ErrInvalid
	}
	return nil
}
func (s *Service) Definition(ctx context.Context, scope Scope, id string) (Definition, error) {
	if e := validate(scope, id, true); e != nil {
		return Definition{}, e
	}
	return s.repository.Definition(ctx, scope, id)
}
func ValidateSubmission(scope Scope, id string, input Submission) error {
	if e := validate(scope, id, true); e != nil {
		return e
	}
	if input.Score < 0 || input.Score > 10 || !input.Consent || !valid(input.ConsentVersion) {
		return ErrInvalid
	}
	return nil
}
func (s *Service) Submit(ctx context.Context, scope Scope, id string, input Submission) (Result, error) {
	if e := ValidateSubmission(scope, id, input); e != nil {
		return Result{}, e
	}
	return s.repository.Submit(ctx, scope, id, input)
}
func (s *Service) OwnAnswer(ctx context.Context, scope Scope, id string) (Answer, error) {
	if e := validate(scope, id, true); e != nil {
		return Answer{}, e
	}
	return s.repository.OwnAnswer(ctx, scope, id)
}
func (s *Service) Summary(ctx context.Context, scope Scope, id string) (Summary, error) {
	if e := validate(scope, id, false); e != nil {
		return Summary{}, e
	}
	return s.repository.Summary(ctx, scope, id)
}

// NPS requires a nonempty population and an explicit project-owned reporting threshold.
// The selected PostHog MIT owner calculates the score; this caller preserves validation
// and the project-owned disclosure threshold. No representative-sampling claim.
func NPS(count, promoters, detractors, minimum int64) (Summary, error) {
	if count < 0 || promoters < 0 || detractors < 0 || minimum < 1 || promoters > count || detractors > count-promoters {
		return Summary{}, ErrInvalid
	}
	result := Summary{Responses: count}
	if count >= minimum {
		value := posthognps.Calculate(count, promoters, detractors).Score
		result.Available = true
		result.NPS = &value
	}
	return result, nil
}
