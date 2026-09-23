package search

import (
	"errors"
	"math"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	if s, err := CosineSimilarity([]float32{1, 0}, []float32{1, 0}); err != nil || math.Abs(s-1.0) > 1e-6 {
		t.Fatalf("identical vectors: got %v, %v", s, err)
	}
	if s, err := CosineSimilarity([]float32{1, 0}, []float32{0, 1}); err != nil || math.Abs(s-0.0) > 1e-6 {
		t.Fatalf("orthogonal vectors: got %v, %v", s, err)
	}
	if s, err := CosineSimilarity([]float32{1, 0}, []float32{-1, 0}); err != nil || math.Abs(s+1.0) > 1e-6 {
		t.Fatalf("opposite vectors: got %v, %v", s, err)
	}
}

func TestCosineDimMismatch(t *testing.T) {
	if _, err := CosineSimilarity([]float32{1, 0, 0}, []float32{1, 0}); !errors.Is(err, ErrDimMismatch) {
		t.Fatalf("dim mismatch accepted: %v", err)
	}
}

func TestCosineZeroVector(t *testing.T) {
	if _, err := CosineSimilarity([]float32{0, 0}, []float32{1, 0}); !errors.Is(err, ErrZeroVector) {
		t.Fatalf("zero vector accepted: %v", err)
	}
}

func TestCosineEmpty(t *testing.T) {
	if _, err := CosineSimilarity(nil, []float32{1}); !errors.Is(err, ErrEmptyEmbedding) {
		t.Fatalf("empty embedding accepted: %v", err)
	}
}
