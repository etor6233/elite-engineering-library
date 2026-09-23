package search

import (
	"context"
	"encoding/json"
	"math"
	"testing"
)

func reviewSearchStore(t *testing.T) *MemoryStore {
	t.Helper()
	s := NewMemoryStore()
	d := doc("t-1", "org-1", "d1", "product", "A", "a")
	d.Embedding = []float32{1, 0}
	d.Facets = map[string]string{"category": "a"}
	if err := s.Upsert(context.Background(), d); err != nil {
		t.Fatal(err)
	}
	return s
}
func TestReviewPositiveVectorScopeAndClone(t *testing.T) {
	s := reviewSearchStore(t)
	q := VectorQuery{TenantID: "t-1", OrganizationID: "org-1", Embedding: []float32{1, 0}, MinSimilarity: .99}
	r, e := s.VectorSearch(context.Background(), q)
	if e != nil || len(r) != 1 || r[0].Score < .99 {
		t.Fatalf("valid query %v %v", r, e)
	}
	r[0].Document.Embedding[0] = 0
	r[0].Document.Facets["category"] = "changed"
	again, e := s.VectorSearch(context.Background(), q)
	if e != nil || len(again) != 1 || again[0].Document.Embedding[0] != 1 || again[0].Document.Facets["category"] != "a" {
		t.Fatal("vector result aliased stored document")
	}
	q.TenantID = "other"
	if r, e = s.VectorSearch(context.Background(), q); e != nil || len(r) != 0 {
		t.Fatalf("scope invalid %v %v", r, e)
	}
}
func TestReviewPositiveStoredNonfiniteRejected(t *testing.T) {
	s := NewMemoryStore()
	d := doc("t-1", "", "d1", "product", "A", "a")
	d.Embedding = []float32{float32(math.NaN())}
	if s.Upsert(context.Background(), d) == nil {
		t.Fatal("stored NaN accepted")
	}
}
func TestReviewDefectNonfiniteVectorQuery(t *testing.T) {
	for name, value := range map[string]float32{"NaN": float32(math.NaN()), "PositiveInfinity": float32(math.Inf(1)), "NegativeInfinity": float32(math.Inf(-1))} {
		t.Run(name, func(t *testing.T) {
			s := reviewSearchStore(t)
			r, e := s.VectorSearch(context.Background(), VectorQuery{TenantID: "t-1", Embedding: []float32{value, 1}, MinSimilarity: .99})
			if e != nil {
				return
			}
			for _, row := range r {
				if math.IsNaN(row.Score) || math.IsInf(row.Score, 0) {
					_, marshalErr := json.Marshal(r)
					t.Fatalf("non-finite similarity returned through threshold: score=%v jsonError=%v", row.Score, marshalErr)
				}
			}
		})
	}
}

func TestNXVectorRejectsThresholdBeforeEmptyStore(t *testing.T) {
	for _, threshold := range []float64{math.NaN(), math.Inf(1), 2, -2} {
		if _, e := NewMemoryStore().VectorSearch(context.Background(), VectorQuery{TenantID: "t", Embedding: []float32{1, 0}, MinSimilarity: threshold}); e == nil {
			t.Fatal("invalid threshold accepted")
		}
	}
}
