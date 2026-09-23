package search

import (
	"context"
	"errors"
	"testing"
)

func TestSearchRankingAndScope(t *testing.T) {
	s := NewMemoryStore()
	a := doc("t-1", "org-1", "d1", "product", "Scooter Urbano Pro", "scooter electrico urbano 350W")
	a.Facets = map[string]string{"category": "movilidad"}
	b := doc("t-1", "org-1", "d2", "product", "Scooter Plegable", "scooter plegable ligero")
	c := doc("t-1", "org-2", "d3", "product", "Scooter Urbano", "scooter urbano") // other org
	d := doc("t-2", "org-1", "d4", "product", "Scooter Urbano", "scooter urbano") // other tenant
	for _, x := range []Document{a, b, c, d} {
		if err := s.Upsert(context.Background(), x); err != nil {
			t.Fatal(err)
		}
	}

	res, err := s.Search(context.Background(), Query{TenantID: "t-1", OrganizationID: "org-1", Text: "scooter urbano"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results scoped to org-1, got %d", len(res))
	}
	if res[0].Document.ID != "d1" {
		t.Fatalf("expected d1 ranked first (title hit), got %s", res[0].Document.ID)
	}
}

func TestSearchFacetAndKindFilter(t *testing.T) {
	s := NewMemoryStore()
	a := doc("t-1", "", "d1", "product", "Casco Integral", "casco seguridad")
	a.Facets = map[string]string{"category": "accesorio"}
	b := doc("t-1", "", "d2", "product", "Casco Abierto", "casco")
	b.Facets = map[string]string{"category": "accesorio"}
	c := doc("t-1", "", "d3", "faq", "Casco", "como elegir casco")
	for _, x := range []Document{a, b, c} {
		if err := s.Upsert(context.Background(), x); err != nil {
			t.Fatal(err)
		}
	}

	res, err := s.Search(context.Background(), Query{TenantID: "t-1", Kind: "product", Text: "casco"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("kind filter failed: got %d", len(res))
	}

	res, err = s.Search(context.Background(), Query{TenantID: "t-1", Text: "casco", Facets: map[string]string{"category": "accesorio"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("facet filter failed: got %d", len(res))
	}
}

func TestSearchDeniesEmptyTenantAndQuery(t *testing.T) {
	s := NewMemoryStore()
	if _, err := s.Search(context.Background(), Query{TenantID: "", Text: "x"}); !errors.Is(err, ErrDenied) {
		t.Fatalf("empty tenant accepted: %v", err)
	}
	if _, err := s.Search(context.Background(), Query{TenantID: "t-1", Text: "   "}); !errors.Is(err, ErrInvalidQuery) {
		t.Fatalf("empty query accepted: %v", err)
	}
}

func TestVectorSearchTopKAndThreshold(t *testing.T) {
	s := NewMemoryStore()
	a := doc("t-1", "", "d1", "product", "A", "x")
	a.Embedding = []float32{1, 0, 0}
	b := doc("t-1", "", "d2", "product", "B", "x")
	b.Embedding = []float32{0, 1, 0}
	c := doc("t-1", "", "d3", "product", "C", "x")
	c.Embedding = []float32{1, 0, 0}
	for _, x := range []Document{a, b, c} {
		if err := s.Upsert(context.Background(), x); err != nil {
			t.Fatal(err)
		}
	}

	res, err := s.VectorSearch(context.Background(), VectorQuery{
		TenantID: "t-1", Embedding: []float32{1, 0, 0}, TopK: 2, MinSimilarity: 0.99,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("expected top-2, got %d", len(res))
	}
	if res[0].Score < 0.999 || res[1].Score < 0.999 {
		t.Fatalf("expected near-1 similarities, got %v %v", res[0].Score, res[1].Score)
	}
}

func TestDelete(t *testing.T) {
	s := NewMemoryStore()
	a := doc("t-1", "", "d1", "product", "A", "x")
	if err := s.Upsert(context.Background(), a); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(context.Background(), "t-1", "d1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(context.Background(), "t-1", "d1"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected not found after delete, got %v", err)
	}
	if err := s.Delete(context.Background(), "", "d1"); !errors.Is(err, ErrDenied) {
		t.Fatalf("empty tenant delete accepted: %v", err)
	}
}
