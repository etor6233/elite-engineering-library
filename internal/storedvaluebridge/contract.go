package storedvaluebridge

// AUTHORED contracts and exact serialization glue. Commercial calculations are
// provided by the separately licensed source-derived Odoo process.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"elite.local/enterprise/internal/approval"
)

var ErrBinding = errors.New("stored value binding is invalid")
var identityRE = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}$`)
var hashRE = regexp.MustCompile(`^[0-9a-f]{64}$`)
var decimalRE = regexp.MustCompile(`^-?(0|[1-9][0-9]{0,17})(\.[0-9]{1,6})?$`)

func Hash(b []byte) string  { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func ValidID(v string) bool { return identityRE.MatchString(v) }
func Decode(raw []byte, target any) error {
	canonical, _, e := approval.CanonicalPayload(raw)
	if e != nil {
		return ErrBinding
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if e := d.Decode(target); e != nil {
		return ErrBinding
	}
	normalized, _, e := Canonical(target)
	if e != nil || !bytes.Equal(canonical, normalized) {
		return ErrBinding
	}
	return nil
}
func Canonical(v any) (json.RawMessage, string, error) {
	b, e := json.Marshal(v)
	if e != nil {
		return nil, "", e
	}
	return approval.CanonicalPayload(b)
}
func Minor(value string, digits int) (int64, error) {
	if !decimalRE.MatchString(value) || digits < 0 || digits > 6 {
		return 0, ErrBinding
	}
	r, ok := new(big.Rat).SetString(value)
	if !ok {
		return 0, ErrBinding
	}
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)
	r.Mul(r, new(big.Rat).SetInt(scale))
	if !r.IsInt() || !r.Num().IsInt64() {
		return 0, ErrBinding
	}
	return r.Num().Int64(), nil
}
func Major(value int64, digits int) string {
	n := big.NewInt(value)
	negative := n.Sign() < 0
	n.Abs(n)
	v := n.String()
	if digits > 0 {
		for len(v) <= digits {
			v = "0" + v
		}
		v = v[:len(v)-digits] + "." + v[len(v)-digits:]
	}
	if negative {
		v = "-" + v
	}
	return v
}

type Rule struct {
	ID            string   `json:"id"`
	Mode          string   `json:"mode"`
	Points        string   `json:"points"`
	Split         bool     `json:"split"`
	MinimumQty    string   `json:"minimum_qty"`
	MinimumAmount string   `json:"minimum_amount"`
	TaxMode       string   `json:"tax_mode"`
	ProductIDs    []string `json:"product_ids"`
	CodeRequired  bool     `json:"code_required"`
}
type Program struct {
	ID                string   `json:"id"`
	Kind              string   `json:"kind"`
	Currency          string   `json:"currency"`
	CurrencyDigits    int      `json:"currency_digits"`
	AppliesOn         string   `json:"applies_on"`
	Nominative        bool     `json:"nominative"`
	Trigger           string   `json:"trigger"`
	TriggerProductIDs []string `json:"trigger_product_ids"`
	Rules             []Rule   `json:"rules"`
}
type Reward struct {
	Mode           string `json:"mode"`
	Discount       string `json:"discount"`
	RequiredPoints string `json:"required_points"`
	MaxAmount      string `json:"max_amount"`
	ClearWallet    bool   `json:"clear_wallet"`
}
type PolicyProgram struct {
	Calculation Program `json:"calculation"`
	Reward      Reward  `json:"reward"`
}
type ProfileConfig struct {
	Schema             string          `json:"schema"`
	TenantID           string          `json:"tenant_id"`
	OrganizationID     string          `json:"organization_id"`
	AmountsMode        string          `json:"amounts_mode"`
	Review             string          `json:"review"`
	ReservationSeconds int             `json:"reservation_seconds"`
	IssuanceStates     []string        `json:"issuance_states"`
	RedemptionStates   []string        `json:"redemption_states"`
	Programs           []PolicyProgram `json:"programs"`
}
type Profile struct {
	config ProfileConfig
	raw    json.RawMessage
	hash   string
}

func Load(raw []byte, expected string) (*Profile, error) {
	var c ProfileConfig
	if Decode(raw, &c) != nil || !hashRE.MatchString(expected) || Hash(raw) != expected {
		return nil, ErrBinding
	}
	if c.Schema != "elite.stored-value-profile.v1" || !ValidID(c.TenantID) || !ValidID(c.OrganizationID) || c.AmountsMode != "bound_gross_inclusive" || c.Review != "one_distinct_human" || c.ReservationSeconds < 60 || c.ReservationSeconds > 86400 || len(c.Programs) == 0 || len(c.Programs) > 8 {
		return nil, ErrBinding
	}
	validStates := func(states []string, allowed string) bool {
		if len(states) == 0 {
			return false
		}
		seen := map[string]bool{}
		for _, s := range states {
			if seen[s] || !strings.Contains("|"+allowed+"|", "|"+s+"|") {
				return false
			}
			seen[s] = true
		}
		return true
	}
	if !validStates(c.IssuanceStates, "confirmed|paid|allocated|delivered") || !validStates(c.RedemptionStates, "placed|confirmed|allocated") {
		return nil, ErrBinding
	}
	seen := map[string]bool{}
	kinds := map[string]bool{}
	for _, p := range c.Programs {
		v := p.Calculation
		if !ValidID(v.ID) || seen[v.ID] || kinds[v.Kind] || (v.Kind != "gift_card" && v.Kind != "loyalty") || v.Currency != "ARS" || v.CurrencyDigits != 2 || len(v.Rules) == 0 || len(v.Rules) > 20 {
			return nil, ErrBinding
		}
		seen[v.ID] = true
		kinds[v.Kind] = true
		if (v.Kind == "gift_card" && (v.AppliesOn != "future" || v.Nominative)) || (v.Kind == "loyalty" && (v.AppliesOn != "both" || !v.Nominative)) {
			return nil, ErrBinding
		}
		if v.Trigger != "auto" || v.TriggerProductIDs == nil {
			return nil, ErrBinding
		}
		positive := func(v string, zero bool) bool {
			r, ok := new(big.Rat).SetString(v)
			return ok && decimalRE.MatchString(v) && (r.Sign() > 0 || (zero && r.Sign() == 0))
		}
		ruleIDs := map[string]bool{}
		for _, r := range v.Rules {
			if !ValidID(r.ID) || ruleIDs[r.ID] || r.TaxMode != "incl" || r.CodeRequired || r.ProductIDs == nil || !positive(r.Points, false) || !positive(r.MinimumQty, true) || !positive(r.MinimumAmount, true) || (r.Mode != "order" && r.Mode != "money" && r.Mode != "unit") || (r.Split && v.AppliesOn == "both") {
				return nil, ErrBinding
			}
			ruleIDs[r.ID] = true
		}
		if p.Reward.Mode != "per_point" && p.Reward.Mode != "per_order" && p.Reward.Mode != "percent" {
			return nil, ErrBinding
		}
		if !positive(p.Reward.Discount, false) || !positive(p.Reward.RequiredPoints, false) || !positive(p.Reward.MaxAmount, true) {
			return nil, ErrBinding
		}
	}
	return &Profile{c, append(json.RawMessage(nil), raw...), expected}, nil
}
func (p *Profile) Valid() bool {
	return p != nil && hashRE.MatchString(p.hash) && Hash(p.raw) == p.hash
}
func (p *Profile) SHA256() string {
	if p == nil {
		return ""
	}
	return p.hash
}
func (p *Profile) Scope() (string, string) {
	if !p.Valid() {
		return "", ""
	}
	return p.config.TenantID, p.config.OrganizationID
}
func (p *Profile) Program(id string) (PolicyProgram, error) {
	if !p.Valid() {
		return PolicyProgram{}, ErrBinding
	}
	for _, v := range p.config.Programs {
		if v.Calculation.ID == id {
			// Keep the validated hash-locked profile immutable across caller DTO changes.
			b, err := json.Marshal(v)
			if err != nil {
				return PolicyProgram{}, ErrBinding
			}
			var copied PolicyProgram
			if json.Unmarshal(b, &copied) != nil {
				return PolicyProgram{}, ErrBinding
			}
			return copied, nil
		}
	}
	return PolicyProgram{}, ErrBinding
}
func (p *Profile) ReservationSeconds() int {
	if !p.Valid() {
		return 0
	}
	return p.config.ReservationSeconds
}
func (p *Profile) StateAllowed(operation, state string) bool {
	if !p.Valid() {
		return false
	}
	states := p.config.RedemptionStates
	if operation == "issue" || operation == "accrue" {
		states = p.config.IssuanceStates
	}
	for _, v := range states {
		if v == state {
			return true
		}
	}
	return false
}

type Line struct {
	ID                string `json:"id"`
	ProductID         string `json:"product_id"`
	Quantity          string `json:"quantity"`
	Subtotal          string `json:"subtotal"`
	Tax               string `json:"tax"`
	Total             string `json:"total"`
	RewardProgramType string `json:"reward_program_type"`
	RewardProgramID   string `json:"reward_program_id"`
	RewardTrigger     string `json:"reward_trigger"`
	ThresholdExcluded bool   `json:"threshold_excluded"`
}
type Order struct {
	OrderID        string   `json:"order_id"`
	OrganizationID string   `json:"organization_id"`
	SubjectID      string   `json:"subject_id"`
	State          string   `json:"state"`
	PublicSubject  bool     `json:"public_subject"`
	Currency       string   `json:"currency"`
	Total          string   `json:"total"`
	EnabledRuleIDs []string `json:"enabled_rule_ids"`
	Lines          []Line   `json:"lines"`
}
type Calculation struct {
	Schema        string          `json:"schema"`
	Operation     string          `json:"operation"`
	Program       Program         `json:"program"`
	ProgramSHA256 string          `json:"program_sha256"`
	Order         Order           `json:"order"`
	Data          json.RawMessage `json:"data"`
}

func NewCalculation(operation string, p Program, o Order, data any) (Calculation, error) {
	_, hash, e := Canonical(p)
	if e != nil {
		return Calculation{}, e
	}
	raw, e := json.Marshal(data)
	return Calculation{"elite.odoo-loyalty-calc.v1", operation, p, hash, o, raw}, e
}

type Result struct {
	Schema         string          `json:"schema"`
	SourceRevision string          `json:"source_revision"`
	ProgramSHA256  string          `json:"program_sha256"`
	RequestSHA256  string          `json:"request_sha256"`
	Result         json.RawMessage `json:"result"`
}
type Request struct {
	OperationID          string `json:"operation_id"`
	Operation            string `json:"operation"`
	OrganizationID       string `json:"organization_id"`
	ProgramID            string `json:"program_id"`
	OrderID              string `json:"order_id"`
	ExpectedOrderVersion int64  `json:"expected_order_version"`
	AccountID            string `json:"account_id"`
	OriginalOperationID  string `json:"original_operation_id"`
	ProfileSHA256        string `json:"profile_sha256"`
}

func (r Request) Validate(p *Profile) error {
	if !p.Valid() || !ValidID(r.OperationID) || !ValidID(r.OrderID) || r.ExpectedOrderVersion < 1 || r.OrganizationID != p.config.OrganizationID || r.ProfileSHA256 != p.hash {
		return ErrBinding
	}
	if _, e := p.Program(r.ProgramID); e != nil {
		return e
	}
	switch r.Operation {
	case "issue", "accrue":
		if r.AccountID != "" || r.OriginalOperationID != "" {
			return ErrBinding
		}
	case "redeem":
		if !ValidID(r.AccountID) || r.OriginalOperationID != "" {
			return ErrBinding
		}
	case "reverse":
		if !ValidID(r.OriginalOperationID) || r.AccountID != "" {
			return ErrBinding
		}
	default:
		return ErrBinding
	}
	return nil
}
func StableID(kind string, parts ...string) string {
	return kind + "_" + Hash([]byte(strings.Join(parts, "\x00")))
}
func Number(v int64) string { return strconv.FormatInt(v, 10) }

// Document returns a copy of the selected immutable policy for authorized review.
func (p *Profile) Document() json.RawMessage {
	if !p.Valid() {
		return nil
	}
	return append(json.RawMessage(nil), p.raw...)
}
