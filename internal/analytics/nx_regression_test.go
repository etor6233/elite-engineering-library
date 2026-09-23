package analytics

import (
	"testing"
	"time"
)

func TestNXFutureMetricIsNotFresh(t *testing.T) {
	n := time.Unix(1000, 0)
	if Fresh(n.Add(time.Second), n, time.Hour) {
		t.Fatal("future metric accepted as fresh")
	}
}
