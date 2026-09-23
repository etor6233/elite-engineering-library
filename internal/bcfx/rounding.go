package bcfx

// AUTHORED exact minor-unit representation/rounding glue. The explicitly selected
// profile is nearest, ties away from zero. This is not an AL built-in equivalence
// claim: the inspected Learn table has a contradictory negative example.
import "math/big"

func RoundMinor(major *big.Rat, decimals uint8, precisionMinor int64) (int64, error) {
	if major == nil || decimals > 9 || precisionMinor <= 0 || major.Num().BitLen() > 640 || major.Denom().BitLen() > 640 {
		return 0, ErrAmount
	}
	factor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	units := new(big.Rat).Mul(major, new(big.Rat).SetInt(factor))
	units.Quo(units, new(big.Rat).SetInt64(precisionMinor))
	n := new(big.Int).Abs(units.Num())
	q, rem := new(big.Int), new(big.Int)
	q.QuoRem(n, units.Denom(), rem)
	if rem.Lsh(rem, 1).Cmp(units.Denom()) >= 0 {
		q.Add(q, big.NewInt(1))
	}
	if units.Sign() < 0 {
		q.Neg(q)
	}
	q.Mul(q, big.NewInt(precisionMinor))
	if !q.IsInt64() {
		return 0, ErrAmount
	}
	return q.Int64(), nil
}
