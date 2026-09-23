package search

import (
	"context"
	"testing"
)

func TestNXStoreAndResultsDoNotAliasCallerMemory(t *testing.T) {
	m := NewMemoryStore()
	d := Document{ID: "d", ExternalID: "d", ContentSHA256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", TenantID: "t", OrganizationID: "o", Kind: "product", Title: "engine", Body: "engine", Facets: map[string]string{"safe": "yes"}, Embedding: []float32{1, 0}}
	if err := m.Upsert(context.Background(), d); err != nil {
		t.Fatal(err)
	}
	d.Facets["safe"] = "no"
	d.Embedding[0] = 0
	q := Query{TenantID: "t", OrganizationID: "o", Text: "engine"}
	r, err := m.Search(context.Background(), q)
	if err != nil || len(r) != 1 {
		t.Fatal(err, r)
	}
	if r[0].Document.Facets["safe"] != "yes" || r[0].Document.Embedding[0] != 1 {
		t.Fatal("write aliases caller")
	}
	r[0].Document.Facets["safe"] = "no"
	r[0].Document.Embedding[0] = 0
	r, _ = m.Search(context.Background(), q)
	if r[0].Document.Facets["safe"] != "yes" || r[0].Document.Embedding[0] != 1 {
		t.Fatal("read aliases store")
	}
}
