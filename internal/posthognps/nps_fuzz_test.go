// AUTHORED domain oracle and boundary qualification for the declared Go adaptation.
package posthognps

import (
	"math"
	"math/big"
	"testing"
)

func checkOracle(t *testing.T, n, p, d int64) {
	t.Helper()
	g := Calculate(n, p, d)
	if g.Total != n || g.Promoters != p || g.Detractors != d || g.Passives != n-p-d {
		t.Fatalf("lost population: %+v", g)
	}
	if n == 0 {
		if g.Score != 0 {
			t.Fatal(g)
		}
		return
	}
	exact := new(big.Rat).SetFrac(big.NewInt(p-d), big.NewInt(n))
	exact.Mul(exact, big.NewRat(100, 1))
	want, _ := exact.Float64()
	if math.IsNaN(g.Score) || math.IsInf(g.Score, 0) || g.Score < -100 || g.Score > 100 || math.Abs(g.Score-want) > 3e-14 {
		t.Fatalf("score %.17g oracle %.17g", g.Score, want)
	}
}
func TestContractPrecisionAndInt64(t *testing.T) {
	for _, v := range [][3]int64{{3, 1, 0}, {3, 0, 1}, {math.MaxInt64, math.MaxInt64, 0}, {math.MaxInt64, 0, math.MaxInt64}, {math.MaxInt64, math.MaxInt64 / 2, math.MaxInt64 / 2}} {
		checkOracle(t, v[0], v[1], v[2])
	}
	if Calculate(3, 1, 0).Score == 33.3 {
		t.Fatal("presentation rounding changed API")
	}
}
func FuzzGroupedCounts(f *testing.F) {
	for _, v := range [][3]uint64{{0, 0, 0}, {17, 6, 7}, {3, 1, 0}, {10, 10, 0}, {14, 0, 14}, {math.MaxInt64, math.MaxInt64, 0}} {
		f.Add(v[0], v[1], v[2])
	}
	f.Fuzz(func(t *testing.T, a, b, c uint64) {
		n := int64(a & math.MaxInt64)
		p := int64(b % (uint64(n) + 1))
		d := int64(c % (uint64(n-p) + 1))
		checkOracle(t, n, p, d)
	})
}
