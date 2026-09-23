package warrantyclaim

import (
	"encoding/json"
	"testing"
	"time"
)

func FuzzWarrantyProfileCalendar(f *testing.F) {
	raw, _ := json.Marshal(fixtureProfileDocument())
	f.Add(raw)
	f.Add([]byte(`{"schema":"elite-warranty-profile/v1","Schema":"ambiguous"}`))
	f.Add([]byte{0xff, '{', '}'})
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 32769 {
			return
		}
		p, err := LoadProfile(raw, SHA(raw))
		if err != nil {
			return
		}
		if !p.Valid() || p.Hash() != SHA(raw) || SHA(p.Bytes()) != SHA(raw) {
			t.Fatal("accepted profile changed identity")
		}
		doc := p.Document()
		again, _ := json.Marshal(doc)
		round, err := LoadProfile(again, SHA(again))
		if err != nil {
			t.Fatal("accepted semantic profile cannot reload", err)
		}
		at := time.Date(2026, 3, 8, 6, 30, 0, 0, time.UTC)
		dates, e1 := p.DatesAt(at)
		reconstructed, e2 := round.DatesAt(at)
		if (e1 == nil) != (e2 == nil) || e1 == nil && dates != reconstructed {
			t.Fatal("JSON representation changed calendar")
		}
		if e1 == nil {
			first, err := CheckCoverage(dates.PartsStart, dates)
			if err != nil || !first.Parts || !first.Labor {
				t.Fatal("accepted term excludes first date", dates, first, err)
			}
			last, err := CheckCoverage(dates.PartsEnd, dates)
			if err != nil || !last.Parts {
				t.Fatal("accepted parts end is not inclusive", dates, last, err)
			}
		}
		copy := p.Bytes()
		if len(copy) > 0 {
			copy[0] ^= 255
		}
		if SHA(p.Bytes()) != p.Hash() {
			t.Fatal("mutable profile bytes escaped")
		}
	})
}
