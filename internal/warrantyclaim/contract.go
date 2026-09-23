// AUTHORED typed configuration and calendar representation glue around the
// admitted BC CheckWarranty predicate. No legal period or coverage is inferred.
package warrantyclaim

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"strings"
	"time"
	_ "time/tzdata"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/warrantycoverage"
)

var (
	ErrInvalid  = errors.New("invalid warranty command")
	ErrConflict = errors.New("warranty conflict")
	ErrNotFound = errors.New("warranty resource unavailable")
	identifier  = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,119}$`)
	tenantID    = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

type ProfileDocument struct {
	Schema                 string   `json:"schema"`
	Scope                  string   `json:"scope"`
	Algorithm              string   `json:"algorithm"`
	AlgorithmRevision      int      `json:"algorithm_revision"`
	TenantID               string   `json:"tenant_id"`
	OrganizationID         string   `json:"organization_id"`
	FactoryOrganizationID  string   `json:"factory_organization_id"`
	PolicyID               string   `json:"policy_id"`
	TermsVersion           string   `json:"terms_version"`
	TermsText              string   `json:"terms_text"`
	BusinessTimeZone       string   `json:"business_time_zone"`
	PartsDurationDays      int      `json:"parts_duration_days"`
	LaborDurationDays      int      `json:"labor_duration_days"`
	WorkReservationSeconds int      `json:"work_reservation_seconds"`
	FaultExclusions        []string `json:"fault_exclusions"`
	Settlement             string   `json:"settlement"`
	AuthorityReference     string   `json:"authority_reference"`
	DecisionReference      string   `json:"decision_reference"`
}

// Profile is immutable after validation. Document returns a copy, including slices.
type Profile struct {
	document ProfileDocument
	raw      []byte
	hash     string
	zone     *time.Location
}

func (p Profile) Document() ProfileDocument {
	d := p.document
	d.FaultExclusions = append([]string(nil), d.FaultExclusions...)
	return d
}
func (p Profile) Bytes() []byte           { return append([]byte(nil), p.raw...) }
func (p Profile) Hash() string            { return p.hash }
func (p Profile) Scope() (string, string) { return p.document.TenantID, p.document.OrganizationID }
func (p Profile) Valid() bool             { return p.zone != nil && len(p.hash) == 64 }
func ValidID(value string) bool           { return identifier.MatchString(value) }
func SHA(raw []byte) string               { h := sha256.Sum256(raw); return hex.EncodeToString(h[:]) }

func LoadProfile(raw []byte, expectedSHA string) (Profile, error) {
	var empty Profile
	if len(raw) == 0 || len(raw) > 32768 || len(expectedSHA) != 64 || SHA(raw) != expectedSHA {
		return empty, ErrInvalid
	}
	if _, _, err := approval.CanonicalPayload(raw); err != nil {
		return empty, ErrInvalid
	}
	var d ProfileDocument
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&d) != nil || decoder.Decode(new(any)) != io.EOF {
		return empty, ErrInvalid
	}
	if d.Schema != "elite-warranty-profile/v1" || d.Scope != "MATERIALIZED_PROFILE" || d.Algorithm != "bc-inclusive-fixed-terms" || d.AlgorithmRevision != 1 || !tenantID.MatchString(d.TenantID) || !ValidID(d.OrganizationID) || !ValidID(d.FactoryOrganizationID) || !ValidID(d.PolicyID) || !ValidID(d.TermsVersion) || !ValidID(d.AuthorityReference) || !ValidID(d.DecisionReference) {
		return empty, ErrInvalid
	}
	if len(strings.TrimSpace(d.TermsText)) < 1 || len(d.TermsText) > 16000 || d.PartsDurationDays < 1 || d.LaborDurationDays < 1 || d.PartsDurationDays > 3652500 || d.LaborDurationDays > 3652500 || len(d.FaultExclusions) > 64 || d.WorkReservationSeconds < 60 || d.WorkReservationSeconds > 86400 || d.Settlement != "INTERNAL_WARRANTY_SERVICE_ACKNOWLEDGEMENT" {
		return empty, ErrInvalid
	}
	seen := map[string]bool{}
	for _, reason := range d.FaultExclusions {
		if !ValidID(reason) || seen[reason] {
			return empty, ErrInvalid
		}
		seen[reason] = true
	}
	zone, err := time.LoadLocation(d.BusinessTimeZone)
	if err != nil || d.BusinessTimeZone == "" || d.BusinessTimeZone == "Local" {
		return empty, ErrInvalid
	}
	return Profile{document: d, raw: append([]byte(nil), raw...), hash: expectedSHA, zone: zone}, nil
}

type Dates struct {
	PartsStart string `json:"parts_start"`
	PartsEnd   string `json:"parts_end"`
	LaborStart string `json:"labor_start"`
	LaborEnd   string `json:"labor_end"`
}

// Integer-day formulas are explicit profile data. Go AddDate represents the
// selected calendar-day operation; it does not choose or infer a warranty term.
// Date evaluation itself remains the separately admitted BC source predicate.
func (p Profile) DatesAt(acceptedAt time.Time) (Dates, error) {
	if !p.Valid() || acceptedAt.IsZero() {
		return Dates{}, ErrInvalid
	}
	local := acceptedAt.In(p.zone)
	start := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.UTC)
	partsEnd := start.AddDate(0, 0, p.document.PartsDurationDays)
	laborEnd := start.AddDate(0, 0, p.document.LaborDurationDays)
	if start.Year() < 1753 || start.Year() > 9999 || partsEnd.Year() > 9999 || laborEnd.Year() > 9999 {
		return Dates{}, ErrInvalid
	}
	return Dates{PartsStart: start.Format(time.DateOnly), PartsEnd: partsEnd.Format(time.DateOnly), LaborStart: start.Format(time.DateOnly), LaborEnd: laborEnd.Format(time.DateOnly)}, nil
}

func ordinal(value string) (int32, error) {
	day, err := time.Parse(time.DateOnly, value)
	if err != nil || len(value) != 10 || day.Format(time.DateOnly) != value || day.Year() < 1753 || day.Year() > 9999 {
		return 0, ErrInvalid
	}
	// Positive Gregorian ordinal, independent of local UTC offsets/DST.
	return int32(day.Unix()/86400 + 719163), nil
}
func CheckCoverage(date string, dates Dates) (warrantycoverage.Coverage, error) {
	values := []string{date, dates.PartsStart, dates.PartsEnd, dates.LaborStart, dates.LaborEnd}
	days := make([]int32, 5)
	for i, value := range values {
		parsed, err := ordinal(value)
		if err != nil {
			return warrantycoverage.Coverage{}, err
		}
		days[i] = parsed
	}
	return warrantycoverage.Evaluate(days[0], warrantycoverage.Period{Start: days[1], End: days[2]}, warrantycoverage.Period{Start: days[3], End: days[4]})
}
func (p Profile) FaultExcluded(reason string) bool {
	for _, code := range p.document.FaultExclusions {
		if code == reason {
			return true
		}
	}
	return false
}
