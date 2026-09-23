package qr

import (
	"context"
	"errors"
	"testing"
)

type fakeResolver struct {
	ok  map[string]bool
	err error
}

func (f *fakeResolver) Resolve(_ context.Context, tenant string, kind Kind, id string) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	return f.ok[tenant+"\x00"+string(kind)+"\x00"+id], nil
}

func TestPayloadRoundTrip(t *testing.T) {
	p, err := NewPayload("tenant-a", KindProduct, "SKU-123")
	if err != nil {
		t.Fatal(err)
	}
	enc, err := p.Encode()
	if err != nil {
		t.Fatal(err)
	}
	got, err := Decode(enc)
	if err != nil {
		t.Fatal(err)
	}
	if got != p {
		t.Fatalf("round-trip mismatch: %+v vs %+v", got, p)
	}
}

func TestPayloadChecksumTamper(t *testing.T) {
	p, _ := NewPayload("tenant-a", KindOrder, "ORD-1")
	p.ID = "ORD-2" // tamper without recomputing checksum
	if err := p.Validate(); !errors.Is(err, ErrChecksum) {
		t.Fatalf("expected checksum mismatch, got %v", err)
	}
}

func TestPayloadInvalidFields(t *testing.T) {
	if _, err := NewPayload("", KindProduct, "SKU"); !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("empty tenant accepted: %v", err)
	}
	if _, err := NewPayload("t", Kind("bad"), "x"); !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("bad kind accepted: %v", err)
	}
	if _, err := NewPayload("t", KindProduct, "bad id!"); !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("bad id accepted: %v", err)
	}
}

func TestVerifyResolvesAndFailsClosed(t *testing.T) {
	p, _ := NewPayload("tenant-a", KindProduct, "SKU-1")
	enc, _ := p.Encode()

	r := &fakeResolver{ok: map[string]bool{"tenant-a\x00product\x00SKU-1": true}}
	if _, err := Verify(context.Background(), enc, "tenant-a", r); err != nil {
		t.Fatalf("valid QR rejected: %v", err)
	}

	missing := &fakeResolver{ok: map[string]bool{}}
	if _, err := Verify(context.Background(), enc, "tenant-a", missing); !errors.Is(err, ErrUnknown) {
		t.Fatalf("expected ErrUnknown, got %v", err)
	}

	// An absent tenant-a object remains unknown even if another tenant has it.
	other := &fakeResolver{ok: map[string]bool{"tenant-b\x00product\x00SKU-1": true}}
	if _, err := Verify(context.Background(), enc, "tenant-a", other); !errors.Is(err, ErrUnknown) {
		t.Fatalf("cross-tenant should fail closed, got %v", err)
	}
}

func TestVerifyRejectsTampered(t *testing.T) {
	p, _ := NewPayload("tenant-a", KindOrder, "ORD-1")
	enc, _ := p.Encode()
	// tamper the encoded JSON's id (the checksum won't match)
	tampered := enc[:len(enc)-2] + `"}` // corrupts JSON → decode fails
	if _, err := Verify(context.Background(), tampered, "tenant-a", &fakeResolver{}); err == nil {
		t.Fatal("tampered payload accepted")
	}
}
