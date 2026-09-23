package bcfx

import (
	"bytes"
	"encoding/json"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func fxFixture(t testing.TB) ([]byte, []byte, Activation) {
	t.Helper()
	p, e := os.ReadFile("../../config/fx/reference-profile.json")
	if e != nil {
		t.Fatal(e)
	}
	r, e := os.ReadFile("../../config/fx/reference-rates.json")
	if e != nil {
		t.Fatal(e)
	}
	var d ProfileDocument
	if json.Unmarshal(p, &d) != nil {
		t.Fatal("profile")
	}
	return p, r, Activation{true, d.ID, d.Revision, digest(p), d.TenantID, d.OrganizationID, d.LocalCurrency}
}
func TestSnapshotExactBindingsAndConversion(t *testing.T) {
	p, r, a := fxFixture(t)
	s, e := LoadSnapshot(p, r, a)
	if e != nil {
		t.Fatal(e)
	}
	now := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	for _, v := range []struct {
		from, to, date string
		amount, want   int64
		rate           string
		rounded        bool
	}{{"EUR", "USD", "2026-01-01", 1000, 1250, "eur-2026-01-01", true}, {"EUR", "USD", "2026-01-15", 1000, 1300, "eur-2026-01-02", true}, {"USD", "EUR", "2026-01-01", 100, 80, "", true}, {"EUR", "GBP", "2026-01-01", 1000, 8333, "eur-2026-01-01", true}, {"EUR", "EUR", "2026-01-01", 103, 103, "", false}, {"EUR", "GBP", "2026-01-01", 0, 0, "", false}, {"EUR", "USD", "2026-01-01", -1000, -1250, "eur-2026-01-01", true}} {
		got, e := s.Convert(v.from, v.to, v.date, v.amount, now)
		if e != nil || got.OutputMinor != v.want || got.FromRateID != v.rate || got.RoundingApplied != v.rounded {
			t.Fatal(v, got, e)
		}
	}
	p[0] = 'x'
	r[0] = 'x'
	raw, _ := s.Bytes()
	raw[0] = 'x'
	if !s.Allows(a.TenantID, a.OrganizationID) || s.Identity().ProfileSHA256 != a.ProfileSHA256 {
		t.Fatal("mutable")
	}
	if _, e = s.Convert("EUR", "USD", "2026-01-01", 1000, now); e != nil {
		t.Fatal(e)
	}
	for _, v := range []struct {
		from, to, date string
		now            time.Time
	}{{"ZZZ", "USD", "2026-01-01", now}, {"EUR", "ZZZ", "2026-01-01", now}, {"EUR", "USD", "2025-12-31", now}, {"EUR", "USD", "2027-01-01", now}, {"EUR", "USD", "2026-01-01", time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)}} {
		if _, e := s.Convert(v.from, v.to, v.date, 0, v.now); e == nil {
			t.Fatal("invalid admitted", v)
		}
	}
}
func TestSnapshotRejectsAmbiguousAndUnsupported(t *testing.T) {
	for _, mode := range []string{"hash", "scope", "duplicate-json", "unknown", "precision", "decimals-missing", "duplicate-currency", "relational", "duplicate-rate", "duplicate-id", "source-hash", "source-unknown", "bad-date", "zero-rate", "exponent", "missing-local", "bad-rounding", "stale-rate"} {
		t.Run(mode, func(t *testing.T) {
			p, r, a := fxFixture(t)
			var d ProfileDocument
			var src SourceDocument
			_ = json.Unmarshal(p, &d)
			_ = json.Unmarshal(r, &src)
			switch mode {
			case "hash":
				a.ProfileSHA256 = strings.Repeat("0", 64)
			case "scope":
				a.OrganizationID = "foreign"
			case "precision":
				d.Currencies[0].RoundingPrecisionMinor = "0"
			case "decimals-missing":
				d.Currencies[0].MinorUnitDecimals = nil
			case "duplicate-currency":
				d.Currencies = append(d.Currencies, d.Currencies[0])
			case "relational":
				value := "USD"
				src.Rates[0].RelationalCurrency = &value
			case "duplicate-rate":
				src.Rates = append(src.Rates, src.Rates[0])
			case "duplicate-id":
				src.Rates[1].ID = src.Rates[0].ID
			case "source-hash":
				d.SourceSHA256 = strings.Repeat("0", 64)
			case "source-unknown":
				src.Rates[0].Currency = "JPY"
			case "bad-date":
				src.Rates[0].StartingDate = "2026-01-01T00:00:00Z"
			case "zero-rate":
				src.Rates[0].Exchange = "0"
			case "exponent":
				src.Rates[0].Exchange = "1e2"
			case "missing-local":
				d.Currencies = d.Currencies[1:]
			case "bad-rounding":
				d.Rounding = "DEFAULT"
			case "stale-rate":
				d.MaximumRateAgeDays = 1
			}
			r, _ = json.Marshal(src)
			if mode != "source-hash" {
				d.SourceSHA256 = digest(r)
			}
			p, _ = json.Marshal(d)
			if mode == "duplicate-json" {
				p = append([]byte(`{"schema":"other",`), p[1:]...)
			}
			if mode == "unknown" {
				p = append([]byte(`{"unknown":true,`), p[1:]...)
			}
			if mode != "hash" {
				a.ProfileSHA256 = digest(p)
			}
			s, e := LoadSnapshot(p, r, a)
			if mode == "stale-rate" {
				if e != nil {
					t.Fatal(e)
				}
				_, e = s.Convert("EUR", "USD", "2026-01-15", 100, time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC))
			}
			if e == nil {
				t.Fatal("invalid snapshot admitted")
			}
		})
	}
}
func TestRoundMinorIndependentDecimalOracle(t *testing.T) {
	raw, e := os.ReadFile("testdata/rounding-oracle.json")
	if e != nil {
		t.Fatal(e)
	}
	var doc struct {
		Vectors []struct {
			Numerator, Denominator string
			Decimals               uint8
			Precision, Want        string
			Overflow               bool
		}
	}
	if json.Unmarshal(raw, &doc) != nil || len(doc.Vectors) < 500 {
		t.Fatal("oracle")
	}
	for i, v := range doc.Vectors {
		n, _ := new(big.Int).SetString(v.Numerator, 10)
		d, _ := new(big.Int).SetString(v.Denominator, 10)
		q, _ := strconv.ParseInt(v.Precision, 10, 64)
		got, e := RoundMinor(new(big.Rat).SetFrac(n, d), v.Decimals, q)
		if v.Overflow {
			if e == nil {
				t.Fatal(i, "overflow accepted")
			}
		} else if e != nil || strconv.FormatInt(got, 10) != v.Want {
			t.Fatal(i, got, e, v.Want)
		}
	}
	t.Logf("independent Decimal oracle vectors=%d", len(doc.Vectors))
}
func FuzzSnapshotExactByteBoundary(f *testing.F) {
	p, r, a := fxFixture(f)
	f.Add(p, uint8(0))
	f.Add(append([]byte(`{"schema":"duplicate",`), p[1:]...), uint8(1))
	f.Fuzz(func(t *testing.T, raw []byte, mode uint8) {
		if len(raw) > 1048576 {
			return
		}
		activation := a
		if mode%2 == 0 {
			activation.ProfileSHA256 = digest(raw)
		}
		s, e := LoadSnapshot(raw, r, activation)
		if e != nil {
			return
		}
		if !bytes.Equal(raw, s.profileRaw) || s.Identity().ProfileSHA256 != digest(raw) || !s.Allows(a.TenantID, a.OrganizationID) {
			t.Fatal("binding invariant")
		}
		mutated, _ := s.Bytes()
		if len(mutated) > 0 {
			mutated[0] ^= 1
		}
		if !bytes.Equal(raw, s.profileRaw) {
			t.Fatal("mutable bytes")
		}
	})
}
