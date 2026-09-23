package agent

import (
	"context"

	"elite.local/enterprise/internal/search"
)

// Knowledge retrieves tenant-scoped supporting text (FAQ, policies) for an
// answer. It never fabricates: empty results are reported, not guessed.
type Knowledge interface {
	Retrieve(ctx context.Context, tenantID, query string, topK int) ([]string, error)
}

// SearchKnowledge retrieves FAQ-style documents from a search.Store.
type SearchKnowledge struct {
	Store search.Store
	Kind  string
}

// Retrieve runs a tenant-scoped, kind-filtered full-text search and returns
// the top results as "title: body" strings.
func (k SearchKnowledge) Retrieve(ctx context.Context, tenantID, query string, topK int) ([]string, error) {
	res, err := k.Store.Search(ctx, search.Query{
		TenantID: tenantID,
		Kind:     k.Kind,
		Text:     query,
		Limit:    topK,
	})
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(res))
	for _, r := range res {
		out = append(out, r.Document.Title+": "+r.Document.Body)
	}
	return out, nil
}
