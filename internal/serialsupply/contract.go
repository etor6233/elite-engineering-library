// AUTHORED bounded contracts for connecting existing procurement, factory and
// inventory owners. No new pricing, stock ledger or corporate implementation.
package serialsupply

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"elite.local/enterprise/internal/approval"
)

var ErrInvalid = errors.New("invalid serial supply request")
var ErrNotFound = errors.New("serial supply not found")
var ErrConflict = errors.New("serial supply conflict")
var idPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
var shaPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

const PolicyCode = "strict-serial-reference/v1"
const PolicyJSON = `{"schema":"elite-serial-supply-policy/v1","scope":"LIBRARY_INFRASTRUCTURE_REFERENCE","quantity":"whole-serialized-units","overreceipt":0,"receipt_state":"quarantine","release":"distinct-human-with-evidence","pricing":"existing-purchase-order-total-unchanged","production_authorized":false}`

func ValidID(s string) bool  { return idPattern.MatchString(s) }
func ValidSHA(s string) bool { return shaPattern.MatchString(s) }
func ValidText(s string, max int) bool {
	if s == "" || len(s) > max || !utf8.ValidString(s) || strings.TrimSpace(s) != s {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}
func Canonical(value any) (json.RawMessage, string, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, "", err
	}
	out, hash, err := approval.CanonicalPayload(raw)
	return out, hash, err
}

type Line struct {
	ID        string `json:"id"`
	VariantID string `json:"variant_id"`
	Quantity  int    `json:"quantity"`
}
type PlanRequest struct {
	PurchaseOrderID       string `json:"purchase_order_id"`
	CommandID             string `json:"command_id"`
	FactoryOrganizationID string `json:"factory_organization_id"`
	DemandReference       string `json:"demand_reference"`
	PolicyCode            string `json:"policy_code"`
	EvidenceSHA256        string `json:"evidence_sha256"`
	Lines                 []Line `json:"lines"`
}

func (r PlanRequest) Valid() bool {
	if !ValidID(r.PurchaseOrderID) || !ValidID(r.CommandID) || !ValidID(r.FactoryOrganizationID) || !ValidText(r.DemandReference, 1024) || r.PolicyCode != PolicyCode || !ValidSHA(r.EvidenceSHA256) || len(r.Lines) < 1 || len(r.Lines) > 32 {
		return false
	}
	ids := map[string]bool{}
	variants := map[string]bool{}
	total := 0
	for _, l := range r.Lines {
		if !ValidID(l.ID) || !ValidID(l.VariantID) || l.Quantity < 1 || l.Quantity > 1000 || ids[l.ID] || variants[l.VariantID] {
			return false
		}
		ids[l.ID] = true
		variants[l.VariantID] = true
		total += l.Quantity
	}
	return total <= 1000
}

type Command struct {
	PurchaseOrderID     string   `json:"purchase_order_id"`
	CommandID           string   `json:"command_id"`
	ExpectedVersion     int64    `json:"expected_version,string"`
	Kind                string   `json:"kind"`
	EvidenceSHA256      string   `json:"evidence_sha256"`
	LineID              string   `json:"line_id,omitempty"`
	UnitID              string   `json:"unit_id,omitempty"`
	SerialNumber        string   `json:"serial_number,omitempty"`
	VIN                 string   `json:"vin,omitempty"`
	BatterySerialNumber string   `json:"battery_serial_number,omitempty"`
	TargetState         string   `json:"target_state,omitempty"`
	ShipmentID          string   `json:"shipment_id,omitempty"`
	Units               []string `json:"units,omitempty"`
	Reason              string   `json:"reason,omitempty"`
}

func (r Command) Valid() bool {
	if !ValidID(r.PurchaseOrderID) || !ValidID(r.CommandID) || r.ExpectedVersion < 1 || !ValidSHA(r.EvidenceSHA256) {
		return false
	}
	allowed := map[string]bool{}
	switch r.Kind {
	case "submit", "confirm", "start", "cancel":
	case "register":
		if !ValidID(r.LineID) || !ValidText(r.SerialNumber, 128) || r.VIN != "" && !ValidText(r.VIN, 128) || r.BatterySerialNumber != "" && !ValidText(r.BatterySerialNumber, 128) {
			return false
		}
		allowed["line_id"] = true
		allowed["serial_number"] = true
		allowed["vin"] = true
		allowed["battery_serial_number"] = true
	case "milestone":
		if !ValidID(r.UnitID) || !map[string]bool{"assembly": true, "quality": true, "released": true, "rejected": true}[r.TargetState] {
			return false
		}
		allowed["unit_id"] = true
		allowed["target_state"] = true
		if r.TargetState == "released" || r.TargetState == "rejected" {
			allowed["reason"] = true
			if !ValidText(r.Reason, 2000) {
				return false
			}
		}
	case "ship", "receive":
		if !ValidID(r.ShipmentID) || len(r.Units) < 1 || len(r.Units) > 100 {
			return false
		}
		seen := map[string]bool{}
		for _, u := range r.Units {
			if !ValidID(u) || seen[u] {
				return false
			}
			seen[u] = true
		}
		allowed["shipment_id"] = true
		allowed["units"] = true
	case "quality", "quality-reject", "reinspect":
		if !ValidID(r.UnitID) {
			return false
		}
		allowed["unit_id"] = true
		allowed["reason"] = true
		if !ValidText(r.Reason, 2000) {
			return false
		}
	default:
		return false
	}
	values := map[string]string{"line_id": r.LineID, "unit_id": r.UnitID, "serial_number": r.SerialNumber, "vin": r.VIN, "battery_serial_number": r.BatterySerialNumber, "target_state": r.TargetState, "shipment_id": r.ShipmentID, "reason": r.Reason}
	for k, v := range values {
		if v != "" && !allowed[k] {
			return false
		}
	}
	return len(r.Units) == 0 || allowed["units"]
}

type Receipt struct {
	PurchaseOrderID string          `json:"purchase_order_id"`
	CommandID       string          `json:"command_id"`
	Version         int64           `json:"version,string"`
	Kind            string          `json:"kind"`
	Actor           string          `json:"actor"`
	RequestSHA256   string          `json:"request_sha256"`
	PayloadSHA256   string          `json:"payload_sha256"`
	Payload         json.RawMessage `json:"payload"`
	RecordedAt      time.Time       `json:"recorded_at"`
	Replay          bool            `json:"replay"`
}
type Plan struct {
	PurchaseOrderID           string   `json:"purchase_order_id"`
	DestinationOrganizationID string   `json:"destination_organization_id"`
	FactoryOrganizationID     string   `json:"factory_organization_id"`
	SupplierID                string   `json:"supplier_id"`
	State                     string   `json:"state"`
	PurchaseVersion           int64    `json:"purchase_version,string"`
	Currency                  string   `json:"currency"`
	TotalMinorUnits           int64    `json:"total_minor_units,string"`
	Version                   int64    `json:"version,string"`
	PolicyCode                string   `json:"policy_code"`
	DemandReference           string   `json:"demand_reference"`
	Lines                     []Line   `json:"lines"`
	Units                     []Unit   `json:"units"`
	Latest                    *Receipt `json:"latest,omitempty"`
	NextUnitID                string   `json:"next_unit_id,omitempty"`
}
type Unit struct {
	FactoryReviewState string `json:"factory_review_state,omitempty"`
	FactoryRequester   string `json:"factory_requester,omitempty"`
	ReceiptReviewState string `json:"receipt_review_state,omitempty"`
	ReceiptRequester   string `json:"receipt_requester,omitempty"`
	ID                 string `json:"id"`
	LineID             string `json:"line_id"`
	State              string `json:"state"`
	SerialNumber       string `json:"serial_number"`
	StockUnitID        string `json:"stock_unit_id,omitempty"`
	StockState         string `json:"stock_state,omitempty"`
	ShipmentID         string `json:"shipment_id,omitempty"`
}
