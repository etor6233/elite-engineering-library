package catalogrelease

import (
	"elite.local/enterprise/internal/approval"
	"encoding/json"
	"testing"
)

func FuzzCatalogSourceBounds(f *testing.F) {
	f.Add([]byte(`{"command_id":"source","model":{"id":"","code":"bike","displayName":"Fixture","vehicleClass":"bicycle","specification":{"description":"bounded"}},"variants":[{"code":"bike-one","display_name":"Fixture","battery_specification":{},"amount_minor_units":"123456","tax_mode":"not-applicable"}],"valid_from":"2026-09-12T00:00:00Z","valid_until":null}`))
	f.Add([]byte(`{"command_id":"a","Command_ID":"b"}`))
	f.Add([]byte(`{}`))
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 32768 {
			return
		}
		canon, _, e := approval.CanonicalPayload(raw)
		if e != nil {
			return
		}
		var value SourceRequest
		if json.Unmarshal(canon, &value) != nil {
			return
		}
		before, _ := json.Marshal(value)
		valid := value.Valid()
		after, _ := json.Marshal(value)
		if string(before) != string(after) {
			t.Fatal("validation mutated source")
		}
		if valid && (value.Model.ID != "" || len(value.Variants) < 1 || len(value.Variants) > 16 || value.ValidFrom.IsZero() || value.ValidUntil != nil && !value.ValidUntil.After(value.ValidFrom)) {
			t.Fatal("unbounded source admitted")
		}
		for _, v := range value.Variants {
			if valid && v.AmountMinorUnits < 0 {
				t.Fatal("negative amount")
			}
		}
	})
}
