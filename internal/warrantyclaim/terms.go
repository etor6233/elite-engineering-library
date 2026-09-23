// AUTHORED wire and receipt binding around the admitted coverage predicate.
package warrantyclaim

import (
	"encoding/json"
	"strings"
	"time"
)

type OfferRequest struct {
	QuoteID       string `json:"quote_id"`
	QuoteVersion  int64  `json:"quote_version,string"`
	ProfileSHA256 string `json:"profile_sha256"`
}
type AcknowledgeRequest struct {
	QuoteID        string `json:"quote_id"`
	QuoteVersion   int64  `json:"quote_version,string"`
	ProfileSHA256  string `json:"profile_sha256"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type Offer struct {
	QuoteID         string          `json:"quote_id"`
	QuoteVersion    int64           `json:"quote_version,string"`
	CustomerSubject string          `json:"customer_subject"`
	ProfileSHA256   string          `json:"profile_sha256"`
	Profile         json.RawMessage `json:"profile"`
	OfferedBy       string          `json:"offered_by"`
	OfferedAt       time.Time       `json:"offered_at"`
	Acknowledged    bool            `json:"acknowledged"`
	EvidenceSHA256  string          `json:"evidence_sha256,omitempty"`
	Replay          bool            `json:"replay"`
}
type Activation struct {
	WarrantyID      string    `json:"warranty_id"`
	HandoverID      string    `json:"handover_id"`
	QuoteID         string    `json:"quote_id"`
	OrderID         string    `json:"order_id"`
	StockUnitID     string    `json:"stock_unit_id"`
	OrganizationID  string    `json:"organization_id"`
	CustomerSubject string    `json:"customer_subject"`
	ProfileSHA256   string    `json:"profile_sha256"`
	TermsVersion    string    `json:"terms_version"`
	Dates           Dates     `json:"dates"`
	AcceptedAt      time.Time `json:"accepted_at"`
	ActivatedBy     string    `json:"activated_by"`
	ActivatedAt     time.Time `json:"activated_at"`
	Replay          bool      `json:"replay"`
}

func ValidSHA(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, c := range value {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}
func StableID(kind string, values ...string) string {
	return kind + "-" + SHA([]byte(strings.Join(values, "\x00")))[:40]
}
