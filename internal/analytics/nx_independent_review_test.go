package analytics

import (
	"testing"
	"time"
)

func TestReviewPositiveFreshBoundaries(t *testing.T) {
	now := time.Unix(1000, 0)
	for name, tc := range map[string]struct {
		computed time.Time
		ttl      time.Duration
		want     bool
	}{"future": {now.Add(time.Nanosecond), time.Minute, false}, "now": {now, time.Minute, true}, "exactTTL": {now.Add(-time.Minute), time.Minute, true}, "expired": {now.Add(-time.Minute - time.Nanosecond), time.Minute, false}, "zeroDate": {time.Time{}, time.Hour, false}, "zeroTTL": {now, 0, false}} {
		t.Run(name, func(t *testing.T) {
			if got := Fresh(tc.computed, now, tc.ttl); got != tc.want {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}
