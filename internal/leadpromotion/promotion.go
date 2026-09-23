// Package leadpromotion converts a durable provider candidate into a CRM lead
// only after an explicit, evidenced contact-policy decision. Ingestion alone
// never grants permission to contact.
package leadpromotion

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	ErrInvalid     = errors.New("leadpromotion: invalid command")
	ErrNotFound    = errors.New("leadpromotion: candidate not found")
	ErrTestLead    = errors.New("leadpromotion: test lead cannot be promoted")
	ErrConflict    = errors.New("leadpromotion: conflicting replay")
	ErrFieldAbsent = errors.New("leadpromotion: configured field absent")
	providerRe     = regexp.MustCompile(`^[a-z][a-z0-9_]{0,31}$`)
	hex64Re        = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

// Command contains only decisions supplied by the project. It deliberately
// has no default policy, field mapping, purpose or consent version.
type Command struct {
	TenantID        string
	Provider        string
	ProviderEventID string
	LeadID          string
	ConsentID       string
	PurposeCode     string
	PolicyVersion   string
	EvidenceSHA256  string
	DecisionAt      time.Time
	MappingVersion  string
	FieldMapping    map[string]string // provider field ID -> email|phone|name
	IdempotencyKey  string
}

func (c Command) Validate() error {
	if strings.TrimSpace(c.TenantID) == "" || !providerRe.MatchString(c.Provider) || strings.TrimSpace(c.ProviderEventID) == "" {
		return fmt.Errorf("%w: source identity", ErrInvalid)
	}
	if strings.TrimSpace(c.LeadID) == "" || strings.TrimSpace(c.ConsentID) == "" || strings.TrimSpace(c.IdempotencyKey) == "" {
		return fmt.Errorf("%w: target identity", ErrInvalid)
	}
	if len(c.IdempotencyKey) < 8 || len(c.IdempotencyKey) > 200 {
		return fmt.Errorf("%w: idempotency key", ErrInvalid)
	}
	if strings.TrimSpace(c.PurposeCode) == "" || strings.TrimSpace(c.PolicyVersion) == "" || strings.TrimSpace(c.MappingVersion) == "" || c.DecisionAt.IsZero() {
		return fmt.Errorf("%w: policy evidence", ErrInvalid)
	}
	if !hex64Re.MatchString(strings.ToLower(c.EvidenceSHA256)) || len(c.FieldMapping) == 0 {
		return fmt.Errorf("%w: evidence or mapping", ErrInvalid)
	}
	seenTargets := map[string]bool{}
	for source, target := range c.FieldMapping {
		if strings.TrimSpace(source) == "" || (target != "email" && target != "phone" && target != "name") || seenTargets[target] {
			return fmt.Errorf("%w: field mapping", ErrInvalid)
		}
		seenTargets[target] = true
	}
	if !seenTargets["email"] && !seenTargets["phone"] {
		return fmt.Errorf("%w: email or phone mapping required", ErrInvalid)
	}
	return nil
}

type Field struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

// ContactPayload applies only the project's admitted mapping. Unknown fields
// remain in the immutable source candidate and are not guessed into CRM.
func ContactPayload(fields []Field, mapping map[string]string) (json.RawMessage, error) {
	values := map[string]string{}
	for _, field := range fields {
		if target, ok := mapping[field.ID]; ok {
			if strings.TrimSpace(field.Value) == "" || values[target] != "" {
				return nil, fmt.Errorf("%w: %s", ErrInvalid, target)
			}
			values[target] = strings.TrimSpace(field.Value)
		}
	}
	for source, target := range mapping {
		_ = source
		if (target == "email" || target == "phone") && values[target] == "" {
			return nil, fmt.Errorf("%w: %s", ErrFieldAbsent, target)
		}
	}
	return json.Marshal(values)
}

func (c Command) RequestSHA256() (string, error) {
	if err := c.Validate(); err != nil {
		return "", err
	}
	keys := make([]string, 0, len(c.FieldMapping))
	for key := range c.FieldMapping {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	type pair struct{ Source, Target string }
	pairs := make([]pair, 0, len(keys))
	for _, key := range keys {
		pairs = append(pairs, pair{key, c.FieldMapping[key]})
	}
	payload, err := json.Marshal(struct {
		TenantID, Provider, ProviderEventID, LeadID, ConsentID, PurposeCode, PolicyVersion, EvidenceSHA256, DecisionAt, MappingVersion string
		Mapping                                                                                                                        []pair
	}{c.TenantID, c.Provider, c.ProviderEventID, c.LeadID, c.ConsentID, c.PurposeCode, c.PolicyVersion, strings.ToLower(c.EvidenceSHA256), c.DecisionAt.UTC().Format(time.RFC3339Nano), c.MappingVersion, pairs})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

type Receipt struct {
	LeadID   string
	Replayed bool
}

type Store interface {
	Promote(context.Context, Command) (Receipt, error)
}

type Service struct{ Store Store }

func (s Service) Promote(ctx context.Context, command Command) (Receipt, error) {
	if s.Store == nil {
		return Receipt{}, errors.New("leadpromotion: nil store")
	}
	if err := command.Validate(); err != nil {
		return Receipt{}, err
	}
	return s.Store.Promote(ctx, command)
}
