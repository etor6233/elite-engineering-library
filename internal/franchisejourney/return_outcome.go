// AUTHORED read model over existing return owners. No command or new financial rule.
package franchisejourney

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ReturnOwnerOutcome struct {
	State            string `json:"state"`
	Reference        string `json:"reference"`
	AmountMinorUnits string `json:"amount_minor_units,omitempty"`
	Currency         string `json:"currency,omitempty"`
	HandoverState    string `json:"handover_state,omitempty"`
}
type ReturnOutcomeStage struct {
	RequestID    string              `json:"request_id"`
	Kind         string              `json:"kind"`
	Status       string              `json:"status"`
	Attempts     int                 `json:"attempts"`
	UpdatedAt    time.Time           `json:"updated_at"`
	ErrorCode    string              `json:"error_code,omitempty"`
	ResultSHA256 string              `json:"result_sha256,omitempty"`
	Owner        *ReturnOwnerOutcome `json:"owner,omitempty"`
}
type ReturnOutcome struct {
	AuthorizationID string               `json:"authorization_id"`
	OrganizationID  string               `json:"organization_id"`
	OrderID         string               `json:"order_id"`
	DispositionID   string               `json:"disposition_id,omitempty"`
	Remedy          string               `json:"remedy,omitempty"`
	ObservedAt      time.Time            `json:"observed_at"`
	Stages          []ReturnOutcomeStage `json:"stages"`
}

var returnPositiveMinor = regexp.MustCompile(`^[1-9][0-9]{0,18}$`)
var returnCurrency = regexp.MustCompile(`^[A-Z]{3}$`)
var returnErrorCode = regexp.MustCompile(`^[A-Z][A-Z0-9_]{1,63}$`)

func (v ReturnOutcome) Validate(organization, authorization string) error {
	bounded := func(v string, max int) bool { return len(strings.TrimSpace(v)) > 0 && len(v) <= max }
	if v.AuthorizationID != authorization || v.OrganizationID != organization || !bounded(v.OrderID, 128) || v.ObservedAt.IsZero() {
		return ErrConflict
	}
	if v.DispositionID == "" {
		if len(v.Stages) != 0 || v.Remedy != "" {
			return ErrConflict
		}
		return nil
	}
	if !bounded(v.DispositionID, 128) || (v.Remedy != "refund" && v.Remedy != "exchange") {
		return ErrConflict
	}
	expected := map[string]bool{"inventory": true, "accounting": true, v.Remedy: true}
	if v.Remedy == "refund" {
		expected["fiscal"] = true
	}
	if len(v.Stages) != len(expected) {
		return ErrConflict
	}
	ids := map[string]bool{}
	for _, x := range v.Stages {
		if !expected[x.Kind] || ids[x.RequestID] || !bounded(x.RequestID, 128) || x.Attempts < 0 || x.UpdatedAt.IsZero() {
			return ErrConflict
		}
		delete(expected, x.Kind)
		ids[x.RequestID] = true
		switch x.Status {
		case "requested", "claimed", "retry", "blocked", "failed", "succeeded":
		default:
			return ErrConflict
		}
		if x.ErrorCode != "" && !returnErrorCode.MatchString(x.ErrorCode) {
			return ErrConflict
		}
		if x.ResultSHA256 != "" && !sha256Hex(x.ResultSHA256) {
			return ErrConflict
		}
		if x.Status == "succeeded" && (x.ResultSHA256 == "" || x.ErrorCode != "") {
			return ErrConflict
		}
		if (x.Status == "retry" || x.Status == "blocked" || x.Status == "failed") && x.ErrorCode == "" {
			return ErrConflict
		}
		if x.Owner == nil {
			if x.Status == "succeeded" && x.Kind != "inventory" {
				return ErrConflict
			}
			continue
		}
		o := x.Owner
		if x.Kind == "inventory" || !bounded(o.State, 40) || len(o.Reference) > 255 {
			return ErrConflict
		}
		amount, amountErr := strconv.ParseInt(o.AmountMinorUnits, 10, 64)
		if amountErr != nil || amount <= 0 || !returnPositiveMinor.MatchString(o.AmountMinorUnits) || !returnCurrency.MatchString(o.Currency) {
			return ErrConflict
		}
		final := ""
		switch x.Kind {
		case "refund":
			switch o.State {
			case "prepared", "pending", "requires_action", "blocked", "succeeded", "failed", "canceled":
			default:
				return ErrConflict
			}
			final = "succeeded"
			if o.State != "prepared" && o.Reference == "" {
				return ErrConflict
			}
		case "exchange":
			if o.State != "prepared" || o.Reference == "" || !bounded(o.HandoverState, 40) {
				return ErrConflict
			}
			final = "prepared"
		case "accounting":
			if o.State != "posted" || o.Reference == "" {
				return ErrConflict
			}
			final = "posted"
		case "fiscal":
			switch o.State {
			case "queued", "authorizing", "reconcile_required", "authorized", "rejected":
			default:
				return ErrConflict
			}
			if o.Reference == "" {
				return ErrConflict
			}
			final = "authorized"
		}
		if x.Status == "succeeded" && o.State != final {
			return ErrConflict
		}
	}
	return nil
}
func (s *Service) ReturnOutcome(ctx context.Context, tenant, organization, authorization string) (ReturnOutcome, error) {
	if tenant == "" || organization == "" || authorization == "" || len(authorization) > 128 {
		return ReturnOutcome{}, ErrInvalid
	}
	repo, ok := s.repository.(interface {
		ReturnOutcome(context.Context, string, string, string) (ReturnOutcome, error)
	})
	if !ok {
		return ReturnOutcome{}, ErrConflict
	}
	v, e := repo.ReturnOutcome(ctx, tenant, organization, authorization)
	if e != nil {
		return ReturnOutcome{}, e
	}
	if e = v.Validate(organization, authorization); e != nil {
		return ReturnOutcome{}, e
	}
	return v, nil
}
