// AUTHORED read projection: all decisions remain in the existing owners.
package warrantyclaim

import "time"

type RoleContext struct {
	Description    string              `json:"description"`
	Severity       string              `json:"severity"`
	Diagnosis      *RoleDiagnosis      `json:"diagnosis,omitempty"`
	Plan           *RolePlan           `json:"plan,omitempty"`
	Work           *RoleWork           `json:"work,omitempty"`
	Quality        *RoleQuality        `json:"quality,omitempty"`
	Acceptance     *RoleAcceptance     `json:"acceptance,omitempty"`
	Reconciliation *RoleReconciliation `json:"reconciliation,omitempty"`
}
type RoleDiagnosis struct {
	FaultCode      string `json:"fault_code"`
	Description    string `json:"description"`
	Excluded       bool   `json:"excluded"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type RolePlan struct {
	Request       Plan      `json:"request"`
	Requester     string    `json:"requester"`
	PayloadSHA256 string    `json:"payload_sha256"`
	Excluded      bool      `json:"excluded"`
	ExpiresAt     time.Time `json:"expires_at"`
	ApprovalState string    `json:"approval_state"`
}
type RoleWork struct {
	Actor          string `json:"actor"`
	PayloadSHA256  string `json:"payload_sha256"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type RoleQuality struct {
	Passed         bool   `json:"passed"`
	PayloadSHA256  string `json:"payload_sha256"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type RoleAcceptance struct {
	Actor         string `json:"actor"`
	PayloadSHA256 string `json:"payload_sha256"`
}
type RoleReconciliation struct {
	Settlement            string `json:"settlement"`
	RecordedInventoryCost string `json:"recorded_inventory_cost"`
}

// QuoteForOffer avoids operator-supplied version numbers; BindOffer still locks
// and validates the real quote when a write is requested.
type QuoteForOffer struct {
	QuoteID         string `json:"quote_id"`
	OrganizationID  string `json:"organization_id"`
	QuoteVersion    int64  `json:"quote_version,string"`
	CustomerSubject string `json:"customer_subject"`
	State           string `json:"state"`
	Current         bool   `json:"current"`
	Currency        string `json:"currency"`
	TotalMinorUnits int64  `json:"total_minor_units,string"`
}
