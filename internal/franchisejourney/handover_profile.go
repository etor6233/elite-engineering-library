package franchisejourney

// AUTHORED configuration/validation glue. This selects the existing supported
// owner algorithm; it does not introduce pricing, credit, shipping or tax rules.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"
)

const HandoverProfileSchema = "elite-handover-profile/v1"
const HandoverSupportedAlgorithm = "single-unit-full-observed-payment"
const HandoverSupportedAlgorithmRevision = 1

type HandoverProfileOptions struct {
	Quantity        int    `json:"quantity"`
	PaymentCoverage string `json:"payment_coverage"`
	StockSelection  string `json:"stock_selection"`
	Reservation     string `json:"reservation"`
	Refunds         string `json:"refunds"`
	Disputes        string `json:"disputes"`
	Acceptance      string `json:"acceptance"`
	ReleaseEffect   string `json:"release_effect"`
}

func SupportedHandoverProfileOptions() HandoverProfileOptions {
	return HandoverProfileOptions{Quantity: 1, PaymentCoverage: "FULL_ORDER", StockSelection: "ALLOCATED_SERIALIZED_UNIT", Reservation: "ACTIVE_MATCHING", Refunds: "ZERO", Disputes: "DENY", Acceptance: "CUSTOMER_AND_REQUIRED_CHECKLIST", ReleaseEffect: "READ_ONLY_ELIGIBILITY"}
}

type HandoverProfileDocument struct {
	Schema                       string                 `json:"schema"`
	ProfileID                    string                 `json:"profile_id"`
	Revision                     int                    `json:"revision"`
	Algorithm                    string                 `json:"algorithm"`
	AlgorithmRevision            int                    `json:"algorithm_revision"`
	Scope                        string                 `json:"scope"`
	TenantID                     string                 `json:"tenant_id"`
	OrganizationID               string                 `json:"organization_id"`
	ProviderCode                 string                 `json:"provider_code"`
	ProviderAccountRef           string                 `json:"provider_account_ref"`
	ProviderConnectionID         string                 `json:"provider_connection_id"`
	ExpectedLiveMode             *bool                  `json:"expected_live_mode"`
	MaximumObservationAgeSeconds int64                  `json:"maximum_observation_age_seconds"`
	Options                      HandoverProfileOptions `json:"options"`
	// Documentary references record who selected the supported profile. They are
	// not runtime credential proof, a signature or live readiness certification.
	AuthorityReference string `json:"authority_reference"`
	DecisionReference  string `json:"decision_reference"`
}
type HandoverProfileActivation struct {
	Enabled              bool   `json:"enabled"`
	ProfileID            string `json:"profile_id"`
	ProfileRevision      int    `json:"profile_revision"`
	DocumentSHA256       string `json:"document_sha256"`
	TenantID             string `json:"tenant_id"`
	OrganizationID       string `json:"organization_id"`
	ProviderCode         string `json:"provider_code"`
	ProviderAccountRef   string `json:"provider_account_ref"`
	ProviderConnectionID string `json:"provider_connection_id"`
	ExpectedLiveMode     bool   `json:"expected_live_mode"`
}
type handoverProfileBinding struct {
	id, sha, tenant, organization string
	provider, account, connection string
	releaseEffect                 string
	storedValue                   bool
	mode                          bool
	age                           time.Duration
}

var handoverProfileID = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,79}$`)
var handoverTenantID = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// LoadHandoverProfile requires an external exact-byte lock plus an explicit
// activation matching the deployment's identity and provider mode. A JSON
// policy:true or an unrecognized algorithm cannot admit arbitrary behavior.
func LoadHandoverProfile(raw []byte, a HandoverProfileActivation) (HandoverReleaseContract, error) {
	var empty HandoverReleaseContract
	if len(raw) == 0 || len(raw) > 32768 || !a.Enabled || !sha256Hex(a.DocumentSHA256) {
		return empty, ErrReleaseConditioned
	}
	digest := sha256.Sum256(raw)
	if hex.EncodeToString(digest[:]) != a.DocumentSHA256 || !uniqueProfileJSON(raw) {
		return empty, ErrReleaseConditioned
	}
	var d HandoverProfileDocument
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&d) != nil || decoder.Decode(new(any)) != io.EOF {
		return empty, ErrReleaseConditioned
	}
	if d.Schema != HandoverProfileSchema || !handoverProfileID.MatchString(d.ProfileID) || d.Revision < 1 || d.Revision > 1000000 || d.Algorithm != HandoverSupportedAlgorithm || !supportedHandoverAlgorithm(d.AlgorithmRevision, d.Options) || d.Scope != "MATERIALIZED_PROFILE" || !handoverTenantID.MatchString(d.TenantID) || !validPreparationID(d.OrganizationID) || d.ExpectedLiveMode == nil || d.MaximumObservationAgeSeconds < 1 || d.MaximumObservationAgeSeconds > 900 || !validPreparationID(d.AuthorityReference) || !validPreparationID(d.DecisionReference) {
		return empty, ErrReleaseConditioned
	}
	if (d.ProviderCode != "stripe" && d.ProviderCode != "mercadopago") || !validPreparationID(d.ProviderAccountRef) || !validPreparationID(d.ProviderConnectionID) {
		return empty, ErrReleaseConditioned
	}
	if a.ProfileID != d.ProfileID || a.ProfileRevision != d.Revision || a.TenantID != d.TenantID || a.OrganizationID != d.OrganizationID || a.ExpectedLiveMode != *d.ExpectedLiveMode || a.ProviderCode != d.ProviderCode || a.ProviderAccountRef != d.ProviderAccountRef || a.ProviderConnectionID != d.ProviderConnectionID {
		return empty, ErrReleaseConditioned
	}
	id := d.ProfileID + "@" + strconv.Itoa(d.Revision)
	age := time.Duration(d.MaximumObservationAgeSeconds) * time.Second
	bound := &handoverProfileBinding{id: id, sha: a.DocumentSHA256, tenant: d.TenantID, organization: d.OrganizationID, mode: *d.ExpectedLiveMode, age: age, provider: d.ProviderCode, account: d.ProviderAccountRef, connection: d.ProviderConnectionID, releaseEffect: d.Options.ReleaseEffect, storedValue: d.AlgorithmRevision == 3 && d.Options == SupportedStoredValueReleaseOptions()}
	return HandoverReleaseContract{ID: id, DocumentSHA256: a.DocumentSHA256, Scope: d.Scope, ExpectedLiveMode: *d.ExpectedLiveMode, MaximumObservationAge: age, profile: bound}, nil
}

func LoadHandoverProfileFile(path string, a HandoverProfileActivation) (HandoverReleaseContract, error) {
	if path == "" {
		return HandoverReleaseContract{}, ErrReleaseConditioned
	}
	f, err := os.Open(path)
	if err != nil {
		return HandoverReleaseContract{}, ErrReleaseConditioned
	}
	defer f.Close()
	raw, err := io.ReadAll(io.LimitReader(f, 32769))
	if err != nil {
		return HandoverReleaseContract{}, ErrReleaseConditioned
	}
	return LoadHandoverProfile(raw, a)
}

// Duplicate keys have ambiguous human review semantics even with a byte lock.
func uniqueProfileJSON(raw []byte) bool {
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	var visit func() bool
	visit = func() bool {
		token, err := d.Token()
		if err != nil {
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
				k, err := d.Token()
				key, ok := k.(string)
				if err != nil || !ok || seen[key] {
					return false
				}
				seen[key] = true
				if !visit() {
					return false
				}
			}
			end, err := d.Token()
			return err == nil && end == json.Delim('}')
		case '[':
			for d.More() {
				if !visit() {
					return false
				}
			}
			end, err := d.Token()
			return err == nil && end == json.Delim(']')
		default:
			return false
		}
	}
	if !visit() {
		return false
	}
	_, err := d.Token()
	return err == io.EOF
}

func (p HandoverReleaseContract) AllowsScope(tenant, organization string) bool {
	if !p.Valid() {
		return false
	}
	return p.profile == nil || p.profile.tenant == tenant && p.profile.organization == organization
}

func (p HandoverReleaseContract) PaymentBinding() (provider, account, connection string, required bool) {
	if p.profile == nil {
		return "", "", "", false
	}
	return p.profile.provider, p.profile.account, p.profile.connection, true
}
