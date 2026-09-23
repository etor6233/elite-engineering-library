package search

import "context"

// Store is the search index boundary. This helper supplies MemoryStore only.
// The catalog owner has a separate PostgreSQL publication/search implementation;
// this interface does not imply a generic durable indexer or global search API.
type Store interface {
	Upsert(ctx context.Context, doc Document) error
	Delete(ctx context.Context, tenantID, id string) error
	Search(ctx context.Context, q Query) ([]Result, error)
	VectorSearch(ctx context.Context, q VectorQuery) ([]Result, error)
}

// Query is a full-text search scoped by tenant (mandatory) and optionally by
// organization, kind and facets. OrganizationID is enforced deny-by-default by
// the caller: the store never widens scope.
type Query struct {
	TenantID       string
	OrganizationID string
	Kind           string
	Text           string
	Facets         map[string]string
	Limit          int
	MinRank        float64
}

// VectorQuery is a cosine similarity search over embeddings.
type VectorQuery struct {
	TenantID       string
	OrganizationID string
	Kind           string
	Embedding      []float32
	TopK           int
	MinSimilarity  float64
}

// Result pairs a document with its relevance score (rank or similarity).
type Result struct {
	Document Document
	Score    float64
}
