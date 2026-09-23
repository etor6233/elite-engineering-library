package warrantyclaim

import (
	"encoding/json"
	"os"
	"testing"
)

func TestWarrantyRoleWireGoldens(t *testing.T) {
	raw, e := os.ReadFile("../../deploy/warranty/role-goldens.json")
	if e != nil {
		t.Fatal(e)
	}
	var cases []struct {
		Name      string
		Payload   json.RawMessage
		Canonical string
		SHA256    string
	}
	if e = json.Unmarshal(raw, &cases); e != nil {
		t.Fatal(e)
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			var v any
			switch c.Name {
			case "offer":
				v = new(OfferRequest)
			case "acknowledge":
				v = new(AcknowledgeRequest)
			case "activate":
				v = new(struct {
					HandoverID string `json:"handover_id"`
				})
			case "open":
				v = new(OpenClaim)
			case "diagnose":
				v = new(Diagnose)
			case "plan":
				v = new(Plan)
			case "decide":
				v = new(DecideRepair)
			case "work":
				v = new(CompleteWork)
			case "quality":
				v = new(Quality)
			case "accept":
				v = new(AcceptRepair)
			case "reconcile":
				v = new(ReconcileRepair)
			case "cancel":
				v = new(CancelRepair)
			default:
				t.Fatal("unknown")
			}
			if e = json.Unmarshal(c.Payload, v); e != nil {
				t.Fatal(e)
			}
			raw, h, e := Canonical(v)
			if e != nil || string(raw) != c.Canonical || h != c.SHA256 {
				t.Fatal(string(raw), h, e)
			}
		})
	}
}
