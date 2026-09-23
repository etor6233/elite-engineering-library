package search

import (
	"errors"
	"math"
	"strings"
	"testing"
)

func doc(tenant, org, id, kind, title, body string) Document {
	return Document{
		TenantID:       tenant,
		OrganizationID: org,
		ID:             id,
		Kind:           kind,
		ExternalID:     "ext-" + id,
		Title:          title,
		Body:           body,
		Facets:         map[string]string{},
		ContentSHA256:  strings.Repeat("a", 64),
	}
}

func TestDocumentValidate(t *testing.T) {
	valid := doc("t-1", "org-1", "d1", "product", "Scooter Urbano", "movilidad electrica")
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid document rejected: %v", err)
	}

	bad := doc("t-1", "", "d1", "Product", "x", "y") // kind uppercase
	if err := bad.Validate(); !errors.Is(err, ErrInvalidDocument) {
		t.Fatalf("bad kind accepted: %v", err)
	}

	empty := doc("t-1", "", "d1", "product", "  ", "y")
	if err := empty.Validate(); !errors.Is(err, ErrInvalidDocument) {
		t.Fatalf("empty title accepted: %v", err)
	}

	badSHA := doc("t-1", "", "d1", "product", "x", "y")
	badSHA.ContentSHA256 = "zz"
	if err := badSHA.Validate(); !errors.Is(err, ErrInvalidDocument) {
		t.Fatalf("bad content sha accepted: %v", err)
	}

	nan := doc("t-1", "", "d1", "product", "x", "y")
	nan.Embedding = []float32{1.0, float32(math.NaN())} // NaN
	if err := nan.Validate(); !errors.Is(err, ErrInvalidDocument) {
		t.Fatalf("non-finite embedding accepted: %v", err)
	}
}

func TestDocumentValidateEmbedding(t *testing.T) {
	d := doc("t-1", "", "d1", "product", "x", "y")
	d.Embedding = []float32{}
	if err := d.Validate(); !errors.Is(err, ErrInvalidDocument) {
		t.Fatalf("empty embedding slice accepted: %v", err)
	}
	d.Embedding = []float32{1.0, 0.0}
	if err := d.Validate(); err != nil {
		t.Fatalf("finite embedding rejected: %v", err)
	}
}
