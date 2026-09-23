// AUTHORED metamorphic date-translation invariant and actual warranty-date seeds.
package warrantycoverage

import "testing"

func FuzzWarrantyDateTranslation(f *testing.F) {
	// Ordinal dates: 2026-09-12, yearly parts and half-year labor coverage.
	f.Add(int32(739871), int32(739617), int32(739982), int32(739617), int32(739797), int32(1))
	f.Add(int32(0), int32(0), int32(0), int32(0), int32(0), int32(7))
	f.Add(int32(40), int32(50), int32(30), int32(10), int32(60), int32(-1))
	f.Fuzz(func(t *testing.T, date, ps, pe, ls, le, offset int32) {
		values := []int32{date, ps, pe, ls, le}
		shifted := make([]int32, len(values))
		for i, value := range values {
			sum := int64(value) + int64(offset)
			if sum < -2147483648 || sum > 2147483647 {
				return
			}
			shifted[i] = int32(sum)
		}
		before, err := Evaluate(date, Period{ps, pe}, Period{ls, le})
		after, shiftErr := Evaluate(shifted[0], Period{shifted[1], shifted[2]}, Period{shifted[3], shifted[4]})
		if before != after || (err == nil) != (shiftErr == nil) {
			t.Fatalf("date translation changed source decision: %+v/%v -> %+v/%v", before, err, after, shiftErr)
		}
	})
}
