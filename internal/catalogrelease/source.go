// AUTHORED typed binding to the existing model and price services.
package catalogrelease

import (
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/electromobility"
	"encoding/json"
	"regexp"
	"time"
)

type SourceVariant struct {
	Code                 string          `json:"code"`
	DisplayName          string          `json:"display_name"`
	BatterySpecification json.RawMessage `json:"battery_specification"`
	AmountMinorUnits     int64           `json:"amount_minor_units,string"`
	TaxMode              string          `json:"tax_mode"`
}
type SourceRequest struct {
	CommandID  string                `json:"command_id"`
	Model      electromobility.Model `json:"model"`
	Variants   []SourceVariant       `json:"variants"`
	ValidFrom  time.Time             `json:"valid_from"`
	ValidUntil *time.Time            `json:"valid_until"`
}

func (s SourceRequest) Valid() bool {
	if !ValidID(s.CommandID) || s.Model.ID != "" || !ValidText(s.Model.DisplayName, 160) || len(s.Model.Code) > 64 || !sourceObject(s.Model.Specification) || s.ValidFrom.IsZero() || s.ValidUntil != nil && !s.ValidUntil.After(s.ValidFrom) || len(s.Variants) < 1 || len(s.Variants) > 16 {
		return false
	}
	seen := map[string]bool{}
	for _, v := range s.Variants {
		if len(v.Code) > 64 || !regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`).MatchString(v.Code) || seen[v.Code] || !ValidText(v.DisplayName, 160) || !sourceObject(v.BatterySpecification) || v.AmountMinorUnits < 0 {
			return false
		}
		seen[v.Code] = true
	}
	return true
}
func sourceObject(raw json.RawMessage) bool {
	if len(raw) > 8192 {
		return false
	}
	var v map[string]any
	_, _, e := approval.CanonicalPayload(raw)
	return e == nil && json.Unmarshal(raw, &v) == nil && v != nil
}
