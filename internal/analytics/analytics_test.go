package analytics

import (
	"errors"
	"testing"
	"time"
)

func validMetric() MetricDefinition {
	return MetricDefinition{
		TenantID: "tenant-a", Code: "sales.revenue", Version: 1,
		Aggregation: AggSum, Source: "sales.customer_order",
		Dimensions: []string{"organization", "month"}, FreshnessTTL: time.Hour,
	}
}

func TestMetricValidate(t *testing.T) {
	if err := validMetric().Validate(); err != nil {
		t.Fatalf("valid metric rejected: %v", err)
	}
	bad := validMetric()
	bad.Aggregation = "median"
	if err := bad.Validate(); !errors.Is(err, ErrInvalidMetric) {
		t.Fatalf("invalid aggregation accepted: %v", err)
	}
	dup := validMetric()
	dup.Dimensions = []string{"a", "a"}
	if err := dup.Validate(); !errors.Is(err, ErrInvalidMetric) {
		t.Fatalf("duplicate dimension accepted: %v", err)
	}
	zero := validMetric()
	zero.FreshnessTTL = 0
	if err := zero.Validate(); !errors.Is(err, ErrInvalidMetric) {
		t.Fatalf("zero freshness accepted: %v", err)
	}
}

func TestFresh(t *testing.T) {
	now := time.Now()
	if !Fresh(now.Add(-30*time.Minute), now, time.Hour) {
		t.Fatal("fresh value reported stale")
	}
	if Fresh(now.Add(-2*time.Hour), now, time.Hour) {
		t.Fatal("stale value reported fresh")
	}
	if Fresh(now, now, 0) {
		t.Fatal("zero ttl reported fresh (should fail closed)")
	}
}

func TestReconciles(t *testing.T) {
	if !Reconciles([]float64{1, 2, 3}, 6, 0) {
		t.Fatal("exact reconcile failed")
	}
	if !Reconciles([]float64{1, 2, 3}, 6.05, 0.1) {
		t.Fatal("tolerated reconcile failed")
	}
	if Reconciles([]float64{1, 2}, 6, 0.1) {
		t.Fatal("mismatched reconcile passed")
	}
	if Reconciles([]float64{1}, 1, -0.1) {
		t.Fatal("negative tolerance passed")
	}
}

func validRecord() IngestRecord {
	return IngestRecord{
		SourceID: "pos", ExternalID: "rec-1",
		EventTime: time.Now().Add(-time.Minute), ReceivedAt: time.Now(),
		PayloadSHA256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		MaxLateness:   time.Hour,
	}
}

func TestIngestValidate(t *testing.T) {
	if err := validRecord().Validate(); err != nil {
		t.Fatalf("valid record rejected: %v", err)
	}
	late := validRecord()
	late.EventTime = time.Now().Add(-3 * time.Hour) // exceeds MaxLateness (1h)
	if err := late.Validate(); !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("late record accepted: %v", err)
	}
	future := validRecord()
	future.ReceivedAt = future.EventTime.Add(-time.Minute) // received before event
	if err := future.Validate(); !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("received-before-event accepted: %v", err)
	}
	badSHA := validRecord()
	badSHA.PayloadSHA256 = "zz"
	if err := badSHA.Validate(); !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("bad sha accepted: %v", err)
	}
}

func TestIngestReplayKey(t *testing.T) {
	a := validRecord()
	b := validRecord()
	if a.ReplayKey() != b.ReplayKey() {
		t.Fatal("identical records produced different replay keys")
	}
	b.ExternalID = "rec-2"
	if a.ReplayKey() == b.ReplayKey() {
		t.Fatal("different records produced same replay key")
	}
}
