// AUTHORED verification vectors for the declared narrow adaptation.
// These do not claim execution of the Microsoft AL suite.
package bcfx

import (
	"math"
	"math/big"
	"testing"
	"time"
)

func rat(s string) *big.Rat {
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		panic(s)
	}
	return r
}
func direct(e, r string) *DirectRate { return &DirectRate{rat(e), rat(r)} }
func TestExactDirectBaseConversions(t *testing.T) {
	for _, c := range []struct {
		name, amount string
		from, to     *DirectRate
		want         string
	}{
		{"local-to-local", "123.45", nil, nil, "2469/20"},
		{"foreign-to-local", "10", direct("100", "125"), nil, "25/2"},
		{"local-to-foreign", "10", nil, direct("100", "125"), "8"},
		{"foreign-cross", "10", direct("100", "125"), direct("200", "300"), "25/3"},
		{"negative", "-10", direct("100", "125"), direct("200", "300"), "-25/3"},
		{"zero", "0", direct("1", "3"), direct("1", "7"), "0"},
	} {
		t.Run(c.name, func(t *testing.T) {
			a := rat(c.amount)
			before := a.RatString()
			got, err := ExchangeExact(a, c.from, c.to)
			if err != nil || got.Cmp(rat(c.want)) != 0 {
				t.Fatalf("got %v/%v", got, err)
			}
			if a.RatString() != before {
				t.Fatal("mutated input")
			}
		})
	}
	for _, r := range []*DirectRate{direct("0", "1"), direct("1", "0"), direct("-1", "1"), {nil, rat("1")}} {
		if _, err := ExchangeExact(rat("1"), r, nil); err == nil {
			t.Fatal("invalid rate admitted")
		}
	}
}
func TestDestinationRoundingAndBounds(t *testing.T) {
	for _, c := range []struct {
		amount  string
		dec     uint8
		p, want int64
	}{
		{"1.234", 2, 1, 123}, {"1.235", 2, 1, 124}, {"-1.235", 2, 1, -124},
		{"0.025", 2, 5, 5}, {"-0.025", 2, 5, -5}, {"1.024", 2, 5, 100},
		{"9223372036854775807", 0, 1, math.MaxInt64}, {"-9223372036854775808", 0, 1, math.MinInt64},
	} {
		got, err := RoundMinor(rat(c.amount), c.dec, c.p)
		if err != nil || got != c.want {
			t.Fatalf("%s got %d %v", c.amount, got, err)
		}
	}
	for _, s := range []string{"9223372036854775807.5", "-9223372036854775808.5"} {
		if _, err := RoundMinor(rat(s), 0, 1); err == nil {
			t.Fatal("overflow")
		}
	}
	if _, err := RoundMinor(rat("1"), 10, 1); err == nil {
		t.Fatal("scale")
	}
	if _, err := RoundMinor(rat("1"), 2, 0); err == nil {
		t.Fatal("precision")
	}
}
func TestFindLastEffectiveDate(t *testing.T) {
	day := func(s string) time.Time {
		v, e := time.Parse("2006-01-02", s)
		if e != nil {
			t.Fatal(e)
		}
		return v
	}
	a := DatedRate{"USD", day("2026-01-01"), *direct("1", "10")}
	b := DatedRate{"USD", day("2026-02-01"), *direct("1", "12")}
	rows := []DatedRate{b, a, {"EUR", day("2026-01-20"), *direct("1", "20")}}
	for _, c := range []struct{ date, want string }{{"2026-01-01", "10"}, {"2026-01-31", "10"}, {"2026-02-01", "12"}} {
		got, e := FindLast(rows, "USD", day(c.date))
		if e != nil || got.Amounts.Relational.Cmp(rat(c.want)) != 0 {
			t.Fatalf("%v %v", got, e)
		}
		got.Amounts.Relational.SetInt64(99)
	}
	for _, c := range []struct {
		rows     []DatedRate
		currency string
		date     time.Time
	}{{rows, "USD", day("2025-12-31")}, {rows, "GBP", day("2026-02-01")}, {append(rows, a), "USD", day("2026-02-01")}, {rows, "USD", day("2026-02-01").Add(time.Hour)}} {
		if _, e := FindLast(c.rows, c.currency, c.date); e == nil {
			t.Fatal("ambiguous/missing/invalid date")
		}
	}
}
func FuzzDirectConversionExactInverse(f *testing.F) {
	f.Add(int64(123456), uint32(100), uint32(127))
	f.Add(int64(-1), uint32(1), uint32(3))
	f.Add(int64(math.MinInt64), uint32(1), uint32(1))
	f.Fuzz(func(t *testing.T, amount int64, exchange, relational uint32) {
		if exchange == 0 || relational == 0 {
			return
		}
		r := &DirectRate{new(big.Rat).SetInt64(int64(exchange)), new(big.Rat).SetInt64(int64(relational))}
		a := new(big.Rat).SetInt64(amount)
		converted, e := ExchangeExact(a, r, nil)
		if e != nil {
			t.Fatal(e)
		}
		back, e := ExchangeExact(converted, nil, r)
		if e != nil || a.Cmp(back) != 0 {
			t.Fatalf("inverse %v %v", back, e)
		}
		n, e := RoundMinor(back, 0, 1)
		if e != nil || n != amount {
			t.Fatalf("roundtrip %d %v", n, e)
		}
	})
}
