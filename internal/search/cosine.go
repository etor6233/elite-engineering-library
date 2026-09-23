package search

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrNonfiniteEmbedding = errors.New("search: non-finite embedding")
	ErrEmptyEmbedding     = errors.New("search: empty embedding")
	ErrDimMismatch        = errors.New("search: embedding dimension mismatch")
	ErrZeroVector         = errors.New("search: zero vector")
)

// CosineSimilarity returns the cosine of two float32 embeddings. Equal,
// non-empty dimensions are required; a zero vector yields an error (the
// similarity is undefined) rather than an invented score.
func CosineSimilarity(a, b []float32) (float64, error) {
	if len(a) == 0 || len(b) == 0 {
		return 0, ErrEmptyEmbedding
	}
	if len(a) != len(b) {
		return 0, fmt.Errorf("%w: %d vs %d", ErrDimMismatch, len(a), len(b))
	}
	var dot, na, nb float64
	for i := range a {
		av := float64(a[i])
		bv := float64(b[i])
		if math.IsNaN(av) || math.IsInf(av, 0) || math.IsNaN(bv) || math.IsInf(bv, 0) {
			return 0, ErrNonfiniteEmbedding
		}
		dot += av * bv
		na += av * av
		nb += bv * bv
	}
	if na == 0 || nb == 0 {
		return 0, ErrZeroVector
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb)), nil
}
