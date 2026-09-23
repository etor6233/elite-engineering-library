package serialsupply

import (
	"elite.local/enterprise/internal/approval"
	"encoding/json"
	"testing"
)

func FuzzSupplyCreateBounds(f *testing.F) {
	f.Add([]byte(`{"command_id":"golden-command","currency":"ARS","demand_reference":"Demanda á \u003c\u0026\u003e \u2028 interior","destination_organization_id":"store","evidence_sha256":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","factory_organization_id":"factory","lines":[{"id":"line","quantity":2,"variant_id":"variant"}],"policy_code":"strict-serial-reference/v1","purchase_order_id":"golden-po","supplier_id":"supplier","total_minor_units":"9007199254740993"}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"command_id":"a","command_id":"b"}`))
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 32768 {
			return
		}
		canon, _, e := approval.CanonicalPayload(raw)
		if e != nil {
			return
		}
		var value CreateRequest
		if json.Unmarshal(canon, &value) != nil {
			return
		}
		before, _ := json.Marshal(value)
		valid := value.Valid()
		after, _ := json.Marshal(value)
		if string(before) != string(after) {
			t.Fatal("validation mutated input")
		}
		if valid {
			if !value.PlanRequest.Valid() || value.TotalMinorUnits < 0 || !ValidID(value.SupplierID) || !ValidID(value.DestinationOrganizationID) {
				t.Fatal("unbounded source")
			}
			total := 0
			for _, line := range value.Lines {
				total += line.Quantity
			}
			if total < 1 || total > 1000 {
				t.Fatal("quantity scope")
			}
		}
	})
}
