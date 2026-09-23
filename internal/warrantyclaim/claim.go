// AUTHORED typed service/approval/stock orchestration contracts. No new stock
// costing algorithm, reimbursement, payroll or financial ledger is defined.
package warrantyclaim

import (
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/inventorycontrol"
	"encoding/json"
	"regexp"
	"time"
)

type OpenClaim struct {
	CaseID         string `json:"case_id"`
	HandoverID     string `json:"handover_id"`
	AppointmentID  string `json:"appointment_id"`
	Severity       string `json:"severity"`
	Description    string `json:"description"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type Command struct {
	CaseID          string `json:"case_id"`
	CommandID       string `json:"command_id"`
	ExpectedVersion int64  `json:"expected_version,string"`
}

func (r Command) Valid() bool {
	return ValidID(r.CaseID) && ValidID(r.CommandID) && r.ExpectedVersion > 0
}

type Diagnose struct {
	Command
	FaultCode      string `json:"fault_code"`
	Description    string `json:"description"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type Part struct {
	LineID               string `json:"line_id"`
	ItemID               string `json:"item_id"`
	BinID                string `json:"bin_id"`
	LotID                string `json:"lot_id,omitempty"`
	Quantity             string `json:"quantity"`
	SpecificReceiptEntry string `json:"specific_receipt_entry,omitempty"`
}
type Plan struct {
	Command
	LaborWork string `json:"labor_work"`
	Parts     []Part `json:"parts"`
}
type DecideRepair struct {
	Command
	PayloadSHA256 string `json:"payload_sha256"`
	Approved      bool   `json:"approved"`
	Reason        string `json:"reason"`
}
type CompleteWork struct {
	Command
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type Quality struct {
	Command
	Passed                   bool   `json:"passed"`
	WorkEvidenceSHA256       string `json:"work_evidence_sha256"`
	EvidenceSHA256           string `json:"evidence_sha256"`
	CorrectionEvidenceSHA256 string `json:"correction_evidence_sha256,omitempty"`
}
type AcceptRepair struct {
	Command
	QualitySHA256  string `json:"quality_sha256"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type ReconcileRepair struct {
	Command
	AcceptanceSHA256 string `json:"acceptance_sha256"`
	EvidenceSHA256   string `json:"evidence_sha256"`
}
type CancelRepair struct {
	Command
	Reason         string `json:"reason"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type IssuedPart struct {
	LineID        string                           `json:"line_id"`
	ReservationID string                           `json:"reservation_id"`
	Issue         inventorycontrol.BulkIssueResult `json:"issue"`
}
type WorkReceipt struct {
	Request        CompleteWork `json:"request"`
	ApprovalID     string       `json:"approval_id"`
	ApprovedSHA256 string       `json:"approved_sha256"`
	Parts          []IssuedPart `json:"parts"`
}
type Claim struct {
	Context               *RoleContext `json:"context,omitempty"`
	CaseID                string       `json:"case_id"`
	WarrantyID            string       `json:"warranty_id"`
	HandoverID            string       `json:"handover_id"`
	AppointmentID         string       `json:"appointment_id"`
	OrganizationID        string       `json:"organization_id"`
	FactoryOrganizationID string       `json:"factory_organization_id"`
	CustomerSubject       string       `json:"customer_subject"`
	StockUnitID           string       `json:"stock_unit_id"`
	ServiceDate           string       `json:"service_date"`
	State                 string       `json:"state"`
	Version               int64        `json:"version,string"`
	ProfileSHA256         string       `json:"profile_sha256"`
	PartsCovered          bool         `json:"parts_covered"`
	LaborCovered          bool         `json:"labor_covered"`
	Latest                *Step        `json:"latest,omitempty"`
}
type Step struct {
	CaseID        string          `json:"case_id"`
	CommandID     string          `json:"command_id"`
	Version       int64           `json:"version,string"`
	Kind          string          `json:"kind"`
	State         string          `json:"state"`
	Actor         string          `json:"actor"`
	RequestSHA256 string          `json:"request_sha256"`
	PayloadSHA256 string          `json:"payload_sha256"`
	Payload       json.RawMessage `json:"payload"`
	RecordedAt    time.Time       `json:"recorded_at"`
	Replay        bool            `json:"replay"`
}
type RepairPlan struct {
	Schema          string                             `json:"schema"`
	ClaimID         string                             `json:"claim_id"`
	ProfileSHA256   string                             `json:"profile_sha256"`
	DiagnosisSHA256 string                             `json:"diagnosis_sha256"`
	Requester       string                             `json:"requester"`
	Request         Plan                               `json:"request"`
	PartsCovered    bool                               `json:"parts_covered"`
	LaborCovered    bool                               `json:"labor_covered"`
	Excluded        bool                               `json:"excluded"`
	ExpiresAt       time.Time                          `json:"expires_at"`
	Reservations    []inventorycontrol.BulkReservation `json:"reservations"`
}

func Canonical(value any) ([]byte, string, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, "", err
	}
	return approval.CanonicalPayload(raw)
}

var partQuantity = regexp.MustCompile(`^[1-9][0-9]{0,5}$`)

func (r Plan) Valid() bool {
	if !r.Command.Valid() || len(r.LaborWork) > 4000 || len(r.Parts) > 16 || len(r.Parts) == 0 && len(r.LaborWork) == 0 {
		return false
	}
	seen := map[string]bool{}
	for _, p := range r.Parts {
		if !ValidID(p.LineID) || seen[p.LineID] || !ValidID(p.ItemID) || !ValidID(p.BinID) || p.LotID != "" && !ValidID(p.LotID) || !partQuantity.MatchString(p.Quantity) || p.SpecificReceiptEntry != "" && !ValidID(p.SpecificReceiptEntry) {
			return false
		}
		seen[p.LineID] = true
	}
	return true
}
