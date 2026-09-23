// AUTHORED conformance harness against expressions extracted from fixed official AL.
// This is not execution of Microsoft's AL test suite.
package warrantycoverage

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestFixedOfficialDateExpressions(t *testing.T) {
	raw, err := os.ReadFile("testdata/official_date_vectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var vectors []struct {
		Date       int32 `json:"date"`
		PartsStart int32 `json:"parts_start"`
		PartsEnd   int32 `json:"parts_end"`
		LaborStart int32 `json:"labor_start"`
		LaborEnd   int32 `json:"labor_end"`
		Error      bool  `json:"error"`
		Parts      bool  `json:"parts"`
		Labor      bool  `json:"labor"`
		Any        bool  `json:"any"`
	}
	if err := json.Unmarshal(raw, &vectors); err != nil {
		t.Fatal(err)
	}
	if len(vectors) != 1280 {
		t.Fatalf("incomplete official projection corpus: %d", len(vectors))
	}
	for index, v := range vectors {
		got, err := Evaluate(v.Date, Period{v.PartsStart, v.PartsEnd}, Period{v.LaborStart, v.LaborEnd})
		if errors.Is(err, ErrPartsDateOrder) != v.Error || got != (Coverage{Any: v.Any, Parts: v.Parts, Labor: v.Labor}) {
			t.Fatalf("vector %d: got %+v/%v", index, got, err)
		}
	}
}
