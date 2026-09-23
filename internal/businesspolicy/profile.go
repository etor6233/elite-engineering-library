package businesspolicy

// AUTHORED configuration-validation and contract-enforcement glue. No external
// algorithm or vendor authorship is claimed. SQL owner compatibility is explicit.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var ErrProfile = errors.New("business policy profile invalid, unbound or incompatible")

const MaximumProfileBytes = 16384

// These are compatibility limits of existing migration0007, not defaults.
const databaseMaximumSlotSeconds = 28800
const databaseMaximumSlotCapacity = 100

var profileID = regexp.MustCompile(`^[a-z][a-z0-9_-]{1,63}$`)

type document struct {
	Schema       string              `json:"schema"`
	Extends      string              `json:"extends"`
	ProfileID    string              `json:"profile_id"`
	Revision     int64               `json:"revision"`
	Appointments appointmentContract `json:"appointments"`
	Pricebooks   pricebookContract   `json:"pricebooks"`
}
type appointmentContract struct {
	LeadTimeSeconds     *int64   `json:"lead_time_seconds"`
	MaximumSlotSeconds  *int64   `json:"maximum_slot_seconds"`
	MinimumSlotCapacity *int     `json:"minimum_slot_capacity"`
	MaximumSlotCapacity *int     `json:"maximum_slot_capacity"`
	OccupyingStates     []string `json:"occupying_states"`
	WorkingWindow       string   `json:"working_window"`
	UnavailableWindow   string   `json:"unavailable_window"`
}
type pricebookContract struct {
	Selection string `json:"selection"`
	Validity  string `json:"validity"`
}

// Profile is immutable outside this package. Use Load with an independently
// selected SHA256, then give the same profile to both service and repository.
type Profile struct {
	hash  string
	value document
}

func validHash(value string) bool {
	decoded, err := hex.DecodeString(value)
	return err == nil && len(decoded) == sha256.Size && value == strings.ToLower(value)
}

func LoadFile(path, expectedSHA256 string) (*Profile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, ErrProfile
	}
	defer f.Close()
	raw, err := io.ReadAll(io.LimitReader(f, MaximumProfileBytes+1))
	if err != nil {
		return nil, ErrProfile
	}
	return Load(raw, expectedSHA256)
}

func Load(raw []byte, expectedSHA256 string) (*Profile, error) {
	if len(raw) == 0 || len(raw) > MaximumProfileBytes || !utf8.Valid(raw) || !validHash(expectedSHA256) {
		return nil, ErrProfile
	}
	sum := sha256.Sum256(raw)
	if hex.EncodeToString(sum[:]) != expectedSHA256 {
		return nil, ErrProfile
	}
	// encoding/json rejects unknown fields below but normally accepts duplicate
	// keys. Reject duplicates before decoding a contract that controls behavior.
	tokens := json.NewDecoder(bytes.NewReader(raw))
	if err := uniqueJSON(tokens, 0); err != nil {
		return nil, ErrProfile
	}
	if _, err := tokens.Token(); err != io.EOF {
		return nil, ErrProfile
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var value document
	if decoder.Decode(&value) != nil {
		return nil, ErrProfile
	}
	a := value.Appointments
	if value.Schema != "elite-business-policy/v1" || value.Extends != "PBC-CORE/1.0.0" || !profileID.MatchString(value.ProfileID) || value.Revision < 1 || a.LeadTimeSeconds == nil || a.MaximumSlotSeconds == nil || a.MinimumSlotCapacity == nil || a.MaximumSlotCapacity == nil {
		return nil, ErrProfile
	}
	if *a.LeadTimeSeconds < 0 || *a.LeadTimeSeconds > math.MaxInt64/int64(time.Second) || *a.MaximumSlotSeconds < 1 || *a.MaximumSlotSeconds > databaseMaximumSlotSeconds || *a.MinimumSlotCapacity < 1 || *a.MaximumSlotCapacity > databaseMaximumSlotCapacity || *a.MinimumSlotCapacity > *a.MaximumSlotCapacity {
		return nil, ErrProfile
	}
	// Changed modes would contradict existing lifecycle/triggers0007/0009/0015.
	// They require a separately admitted schema/owner change, not silent config.
	if len(a.OccupyingStates) != 2 || a.OccupyingStates[0] != "requested" || a.OccupyingStates[1] != "confirmed" || a.WorkingWindow != "contains_slot" || a.UnavailableWindow != "reject_overlap" || value.Pricebooks.Selection != "single_active_per_market_currency" || value.Pricebooks.Validity != "half_open" {
		return nil, ErrProfile
	}
	return &Profile{hash: expectedSHA256, value: value}, nil
}

func uniqueJSON(decoder *json.Decoder, depth int) error {
	if depth > 16 {
		return ErrProfile
	}
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		seen := map[string]bool{}
		for decoder.More() {
			key, err := decoder.Token()
			if err != nil {
				return err
			}
			name, ok := key.(string)
			if !ok || seen[name] || strings.ToLower(name) != name {
				return ErrProfile
			}
			seen[name] = true
			if err := uniqueJSON(decoder, depth+1); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil || end != json.Delim('}') {
			return ErrProfile
		}
	case '[':
		for decoder.More() {
			if err := uniqueJSON(decoder, depth+1); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil || end != json.Delim(']') {
			return ErrProfile
		}
	default:
		return ErrProfile
	}
	return nil
}

func (p *Profile) Valid() bool { return p != nil && validHash(p.hash) }
func (p *Profile) SHA256() string {
	if !p.Valid() {
		return ""
	}
	return p.hash
}
func (p *Profile) LeadTime() time.Duration {
	return time.Duration(*p.value.Appointments.LeadTimeSeconds) * time.Second
}
func (p *Profile) MaximumSlotDuration() time.Duration {
	return time.Duration(*p.value.Appointments.MaximumSlotSeconds) * time.Second
}
func (p *Profile) MinimumSlotCapacity() int { return *p.value.Appointments.MinimumSlotCapacity }
func (p *Profile) MaximumSlotCapacity() int { return *p.value.Appointments.MaximumSlotCapacity }

func (p *Profile) AllowsSlot(starts, ends, now time.Time, capacity int) bool {
	return p.Valid() && !starts.Before(now.Add(p.LeadTime())) && ends.After(starts) && ends.Sub(starts) <= p.MaximumSlotDuration() && capacity >= p.MinimumSlotCapacity() && capacity <= p.MaximumSlotCapacity()
}

// BindRequestHash preserves historical reference-profile idempotency receipts.
// Nonreference profiles domain-separate the same caller hash. A replay under a
// changed profile therefore conflicts instead of claiming the old decision was
// made using new policy. It does not rewrite any existing key or stored receipt.
func (p *Profile) BindRequestHash(requestHash string) (string, error) {
	if !p.Valid() || !validHash(requestHash) {
		return "", ErrProfile
	}
	if p.hash == ReferenceSHA256 {
		return requestHash, nil
	}
	sum := sha256.Sum256([]byte("business-policy-request/v1\x00" + p.hash + "\x00" + requestHash))
	return hex.EncodeToString(sum[:]), nil
}
