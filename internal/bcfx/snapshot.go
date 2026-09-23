package bcfx

// AUTHORED immutable input/identity boundary around the declared BC adaptation.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"time"
)

var ErrSnapshot = errors.New("invalid, expired or unsupported FX snapshot")

const ProfileSchema = "elite-fx-direct-base-profile/v1"
const SourceSchema = "elite-fx-direct-base-rates/v1"
const Algorithm = "BC_DIRECT_BASE_V1"

type CurrencyPrecision struct {
	Code                   string `json:"code"`
	MinorUnitDecimals      *uint8 `json:"minor_unit_decimals"`
	RoundingPrecisionMinor string `json:"rounding_precision_minor"`
}
type ProfileDocument struct {
	Schema             string              `json:"schema"`
	ID                 string              `json:"profile_id"`
	Revision           int                 `json:"revision"`
	Algorithm          string              `json:"algorithm"`
	Rounding           string              `json:"rounding"`
	TenantID           string              `json:"tenant_id"`
	OrganizationID     string              `json:"organization_id"`
	LocalCurrency      string              `json:"local_currency"`
	ValidFrom          string              `json:"valid_from"`
	ValidUntil         string              `json:"valid_until"`
	DateFrom           string              `json:"conversion_date_from"`
	DateThrough        string              `json:"conversion_date_through"`
	MaximumRateAgeDays int                 `json:"maximum_rate_age_days"`
	SourceID           string              `json:"source_id"`
	SourceRevision     string              `json:"source_revision"`
	SourceSHA256       string              `json:"source_sha256"`
	Currencies         []CurrencyPrecision `json:"currencies"`
}
type RateRow struct {
	ID                 string  `json:"rate_id"`
	Currency           string  `json:"currency"`
	StartingDate       string  `json:"starting_date"`
	Exchange           string  `json:"exchange_rate_amount"`
	Relational         string  `json:"relational_exchange_rate_amount"`
	RelationalCurrency *string `json:"relational_currency"`
}
type SourceDocument struct {
	Schema   string    `json:"schema"`
	ID       string    `json:"source_id"`
	Revision string    `json:"source_revision"`
	Kind     string    `json:"kind"`
	Rates    []RateRow `json:"rates"`
}
type Activation struct {
	Enabled        bool
	ProfileID      string
	Revision       int
	ProfileSHA256  string
	TenantID       string
	OrganizationID string
	LocalCurrency  string
}
type precision struct {
	decimals uint8
	quantum  int64
}
type Snapshot struct {
	profile                                      ProfileDocument
	profileRaw, sourceRaw                        []byte
	hash                                         string
	currencies                                   map[string]precision
	rows                                         []DatedRate
	identities                                   map[string]RateRow
	validFrom, validUntil, dateFrom, dateThrough time.Time
}
type SnapshotIdentity struct {
	ProfileID      string    `json:"profile_id"`
	Revision       int       `json:"profile_revision"`
	ProfileSHA256  string    `json:"profile_sha256"`
	SourceID       string    `json:"source_id"`
	SourceRevision string    `json:"source_revision"`
	SourceSHA256   string    `json:"source_sha256"`
	TenantID       string    `json:"tenant_id"`
	OrganizationID string    `json:"organization_id"`
	LocalCurrency  string    `json:"local_currency"`
	ValidFrom      time.Time `json:"valid_from"`
	ValidUntil     time.Time `json:"valid_until"`
}
type Conversion struct {
	FromCurrency           string `json:"from_currency"`
	ToCurrency             string `json:"to_currency"`
	InputMinor             int64  `json:"input_minor,string"`
	OutputMinor            int64  `json:"output_minor,string"`
	ConversionDate         string `json:"conversion_date"`
	FromDecimals           uint8  `json:"from_decimals"`
	ToDecimals             uint8  `json:"to_decimals"`
	RoundingPrecisionMinor int64  `json:"rounding_precision_minor,string"`
	RoundingApplied        bool   `json:"rounding_applied"`
	ExactMajorNumerator    string `json:"exact_major_numerator"`
	ExactMajorDenominator  string `json:"exact_major_denominator"`
	FromRateID             string `json:"from_rate_id,omitempty"`
	ToRateID               string `json:"to_rate_id,omitempty"`
	FromRateDate           string `json:"from_rate_date,omitempty"`
	ToRateDate             string `json:"to_rate_date,omitempty"`
}

var snapshotID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
var snapshotTenant = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
var currencyCode = regexp.MustCompile(`^[A-Z]{3}$`)
var decimalRate = regexp.MustCompile(`^(0|[1-9][0-9]{0,17})(\.[0-9]{1,18})?$`)

func digest(raw []byte) string { sum := sha256.Sum256(raw); return hex.EncodeToString(sum[:]) }
func hexHash(value string) bool {
	if len(value) != 64 {
		return false
	}
	raw, e := hex.DecodeString(value)
	return e == nil && hex.EncodeToString(raw) == value
}
func exactDate(raw string) (time.Time, error) {
	v, e := time.Parse("2006-01-02", raw)
	if e != nil || len(raw) != 10 || v.IsZero() || v.Format("2006-01-02") != raw {
		return time.Time{}, ErrSnapshot
	}
	return v, nil
}
func exactInstant(raw string) (time.Time, error) {
	v, e := time.Parse(time.RFC3339, raw)
	if e != nil || v.IsZero() || v.UTC().Format(time.RFC3339) != raw {
		return time.Time{}, ErrSnapshot
	}
	return v, nil
}
func decodeSnapshot(raw []byte, value any) error {
	if len(raw) == 0 || len(raw) > 1048576 || !uniqueSnapshotJSON(raw) {
		return ErrSnapshot
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(value) != nil {
		return ErrSnapshot
	}
	if d.Decode(new(any)) != io.EOF {
		return ErrSnapshot
	}
	return nil
}

// DecodeExactJSON enforces the same unambiguous FX transport representation.
func DecodeExactJSON(raw []byte, value any) error { return decodeSnapshot(raw, value) }
func LoadSnapshot(profileRaw, sourceRaw []byte, a Activation) (*Snapshot, error) {
	if !a.Enabled || !hexHash(a.ProfileSHA256) || digest(profileRaw) != a.ProfileSHA256 {
		return nil, ErrSnapshot
	}
	var p ProfileDocument
	var source SourceDocument
	if decodeSnapshot(profileRaw, &p) != nil || decodeSnapshot(sourceRaw, &source) != nil {
		return nil, ErrSnapshot
	}
	if p.Schema != ProfileSchema || p.Algorithm != Algorithm || p.Rounding != "NEAREST_TIES_AWAY_FROM_ZERO" || !snapshotID.MatchString(p.ID) || p.Revision < 1 || p.ID != a.ProfileID || p.Revision != a.Revision || !snapshotTenant.MatchString(p.TenantID) || p.TenantID != a.TenantID || !snapshotID.MatchString(p.OrganizationID) || p.OrganizationID != a.OrganizationID || !currencyCode.MatchString(p.LocalCurrency) || p.LocalCurrency != a.LocalCurrency {
		return nil, ErrSnapshot
	}
	if source.Schema != SourceSchema || source.Kind != "SUPPLIED_SNAPSHOT" || !snapshotID.MatchString(source.ID) || !snapshotID.MatchString(source.Revision) || source.ID != p.SourceID || source.Revision != p.SourceRevision || !hexHash(p.SourceSHA256) || digest(sourceRaw) != p.SourceSHA256 || len(source.Rates) > 4096 || len(p.Currencies) < 1 || len(p.Currencies) > 64 || p.MaximumRateAgeDays < 1 || p.MaximumRateAgeDays > 3660 {
		return nil, ErrSnapshot
	}
	from, e1 := exactInstant(p.ValidFrom)
	until, e2 := exactInstant(p.ValidUntil)
	dateFrom, e3 := exactDate(p.DateFrom)
	dateThrough, e4 := exactDate(p.DateThrough)
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil || !until.After(from) || dateThrough.Before(dateFrom) {
		return nil, ErrSnapshot
	}
	s := &Snapshot{profile: p, hash: a.ProfileSHA256, profileRaw: bytes.Clone(profileRaw), sourceRaw: bytes.Clone(sourceRaw), currencies: map[string]precision{}, identities: map[string]RateRow{}, validFrom: from, validUntil: until, dateFrom: dateFrom, dateThrough: dateThrough}
	for _, v := range p.Currencies {
		q, e := strconv.ParseInt(v.RoundingPrecisionMinor, 10, 64)
		if !currencyCode.MatchString(v.Code) || v.MinorUnitDecimals == nil || *v.MinorUnitDecimals > 9 || e != nil || q <= 0 || strconv.FormatInt(q, 10) != v.RoundingPrecisionMinor {
			return nil, ErrSnapshot
		}
		if _, ok := s.currencies[v.Code]; ok {
			return nil, ErrSnapshot
		}
		s.currencies[v.Code] = precision{*v.MinorUnitDecimals, q}
	}
	if _, ok := s.currencies[p.LocalCurrency]; !ok {
		return nil, ErrSnapshot
	}
	seen := map[string]bool{}
	for _, r := range source.Rates {
		start, e := exactDate(r.StartingDate)
		if !snapshotID.MatchString(r.ID) || seen[r.ID] || r.Currency == p.LocalCurrency || r.RelationalCurrency == nil || *r.RelationalCurrency != "" || e != nil || !decimalRate.MatchString(r.Exchange) || !decimalRate.MatchString(r.Relational) {
			return nil, ErrSnapshot
		}
		if _, ok := s.currencies[r.Currency]; !ok {
			return nil, ErrSnapshot
		}
		key := r.Currency + "@" + r.StartingDate
		if _, ok := s.identities[key]; ok {
			return nil, ErrSnapshot
		}
		exchange, ok1 := new(big.Rat).SetString(r.Exchange)
		relational, ok2 := new(big.Rat).SetString(r.Relational)
		rate := DirectRate{exchange, relational}
		if !ok1 || !ok2 || !validRate(rate) {
			return nil, ErrSnapshot
		}
		seen[r.ID] = true
		s.identities[key] = r
		s.rows = append(s.rows, DatedRate{r.Currency, start, rate})
	}
	return s, nil
}
func LoadSnapshotFiles(profilePath, sourcePath string, a Activation) (*Snapshot, error) {
	read := func(path string) ([]byte, error) {
		f, e := os.Open(path)
		if e != nil {
			return nil, ErrSnapshot
		}
		defer f.Close()
		raw, e := io.ReadAll(io.LimitReader(f, 1048577))
		if e != nil || len(raw) > 1048576 {
			return nil, ErrSnapshot
		}
		return raw, nil
	}
	p, e := read(profilePath)
	if e != nil {
		return nil, e
	}
	r, e := read(sourcePath)
	if e != nil {
		return nil, e
	}
	return LoadSnapshot(p, r, a)
}
func (s *Snapshot) Allows(tenant, organization string) bool {
	return s != nil && s.hash != "" && s.profile.TenantID == tenant && s.profile.OrganizationID == organization
}
func (s *Snapshot) Current(now time.Time) bool {
	return s != nil && !now.IsZero() && !now.Before(s.validFrom) && now.Before(s.validUntil)
}
func (s *Snapshot) Identity() SnapshotIdentity {
	if s == nil {
		return SnapshotIdentity{}
	}
	return SnapshotIdentity{s.profile.ID, s.profile.Revision, s.hash, s.profile.SourceID, s.profile.SourceRevision, s.profile.SourceSHA256, s.profile.TenantID, s.profile.OrganizationID, s.profile.LocalCurrency, s.validFrom, s.validUntil}
}
func (s *Snapshot) Bytes() ([]byte, []byte) {
	if s == nil {
		return nil, nil
	}
	return bytes.Clone(s.profileRaw), bytes.Clone(s.sourceRaw)
}
func (s *Snapshot) ValidateRequest(from, to, date string) error {
	if s == nil {
		return ErrSnapshot
	}
	d, e := exactDate(date)
	_, f := s.currencies[from]
	_, t := s.currencies[to]
	if e != nil || !f || !t || d.Before(s.dateFrom) || d.After(s.dateThrough) {
		return ErrSnapshot
	}
	return nil
}
func (s *Snapshot) Convert(from, to, date string, minor int64, now time.Time) (Conversion, error) {
	var v Conversion
	if s.ValidateRequest(from, to, date) != nil || !s.Current(now) {
		return v, ErrSnapshot
	}
	f, t := s.currencies[from], s.currencies[to]
	day, _ := exactDate(date)
	factor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(f.decimals)), nil)
	major := new(big.Rat).SetFrac(big.NewInt(minor), factor)
	v = Conversion{FromCurrency: from, ToCurrency: to, InputMinor: minor, ConversionDate: date, FromDecimals: f.decimals, ToDecimals: t.decimals, RoundingPrecisionMinor: t.quantum}
	if from == to || minor == 0 {
		v.OutputMinor = minor
		v.ExactMajorNumerator = major.Num().String()
		v.ExactMajorDenominator = major.Denom().String()
		return v, nil
	}
	get := func(currency string) (*DirectRate, string, string, error) {
		if currency == s.profile.LocalCurrency {
			return nil, "", "", nil
		}
		row, e := FindLast(s.rows, currency, day)
		if e != nil || row.StartingDate.Before(day.AddDate(0, 0, -s.profile.MaximumRateAgeDays)) {
			return nil, "", "", ErrSnapshot
		}
		rate := s.identities[currency+"@"+row.StartingDate.Format("2006-01-02")]
		return &row.Amounts, rate.ID, rate.StartingDate, nil
	}
	fr, fid, fd, e := get(from)
	if e != nil {
		return Conversion{}, e
	}
	tr, tid, td, e := get(to)
	if e != nil {
		return Conversion{}, e
	}
	out, e := ExchangeExact(major, fr, tr)
	if e != nil {
		return Conversion{}, e
	}
	rounded, e := RoundMinor(out, t.decimals, t.quantum)
	if e != nil {
		return Conversion{}, e
	}
	v.OutputMinor = rounded
	v.RoundingApplied = true
	v.ExactMajorNumerator = out.Num().String()
	v.ExactMajorDenominator = out.Denom().String()
	v.FromRateID = fid
	v.ToRateID = tid
	v.FromRateDate = fd
	v.ToRateDate = td
	return v, nil
}

// Reuses the existing profile duplicate-key rejection; adds a depth bound.
func uniqueSnapshotJSON(raw []byte) bool {
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	var visit func(int) bool
	visit = func(depth int) bool {
		if depth > 16 {
			return false
		}
		token, e := d.Token()
		if e != nil {
			return false
		}
		delim, ok := token.(json.Delim)
		if !ok {
			return true
		}
		switch delim {
		case '{':
			seen := map[string]bool{}
			for d.More() {
				k, e := d.Token()
				key, ok := k.(string)
				if e != nil || !ok || seen[key] {
					return false
				}
				seen[key] = true
				if !visit(depth + 1) {
					return false
				}
			}
			end, e := d.Token()
			return e == nil && end == json.Delim('}')
		case '[':
			for d.More() {
				if !visit(depth + 1) {
					return false
				}
			}
			end, e := d.Token()
			return e == nil && end == json.Delim(']')
		default:
			return false
		}
	}
	if !visit(0) {
		return false
	}
	_, e := d.Token()
	return e == io.EOF
}
