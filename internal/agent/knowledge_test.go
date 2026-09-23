package agent

import (
	"context"
	"strings"
	"testing"

	"elite.local/enterprise/internal/search"
)

func TestSearchKnowledgeRetrievesScoped(t *testing.T) {
	s := search.NewMemoryStore()
	d := search.Document{
		TenantID:      "t-1",
		ID:            "f1",
		Kind:          "faq",
		ExternalID:    "ext-f1",
		Title:         "Horarios",
		Body:          "Abrimos de 9 a 18",
		Facets:        map[string]string{},
		ContentSHA256: strings.Repeat("a", 64),
	}
	if err := s.Upsert(context.Background(), d); err != nil {
		t.Fatal(err)
	}

	k := SearchKnowledge{Store: s, Kind: "faq"}
	docs, err := k.Retrieve(context.Background(), "t-1", "horarios", 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(docs) != 1 || docs[0] != "Horarios: Abrimos de 9 a 18" {
		t.Fatalf("unexpected docs: %v", docs)
	}

	// other tenant sees nothing
	docs, err = k.Retrieve(context.Background(), "t-2", "horarios", 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(docs) != 0 {
		t.Fatalf("cross-tenant leak: %v", docs)
	}
}
