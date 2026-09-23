package storedvaluebridge

// AUTHORED domain-bound representation invariants, no external effects.
import (
	"math"
	"os"
	"path/filepath"
	"testing"
)

func FuzzStoredValueMinorRoundTrip(f *testing.F) {
	for _, v := range []int64{0, 1, -1, 50, 5000, 123456, 9007199254740993, math.MaxInt64, math.MinInt64} {
		f.Add(v, uint8(2))
	}
	f.Fuzz(func(t *testing.T, v int64, d uint8) {
		digits := int(d % 7)
		raw := Major(v, digits)
		got, e := Minor(raw, digits)
		// The wire decimal contract admits at most18whole digits, so int64 extrema
		// at scale0 are correctly rejected; all representable values round-trip.
		if digits == 0 && (v > 999999999999999999 || v < -999999999999999999) {
			if e == nil {
				t.Fatal("unbounded decimal")
			}
			return
		}
		if e != nil || got != v {
			t.Fatal("amount representation changed", raw, digits, got, e)
		}
	})
}
func FuzzStoredValueProfileBinding(f *testing.F) {
	raw, e := os.ReadFile(filepath.Join("..", "..", "deploy", "stored-value", "profile.reference.json"))
	if e != nil {
		f.Fatal(e)
	}
	f.Add(raw)
	f.Add([]byte(`{"schema":"elite.stored-value-profile.v1","schema":"duplicate"}`))
	f.Add([]byte("null"))
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 8192 {
			return
		}
		p, e := Load(raw, Hash(raw))
		if e != nil {
			return
		}
		if !p.Valid() || Hash(p.Document()) != p.SHA256() {
			t.Fatal("mutable or unbound profile")
		}
		copy := p.Document()
		if len(copy) > 0 {
			copy[0] ^= 1
		}
		if !p.Valid() {
			t.Fatal("document aliases live policy")
		}
	})
}
