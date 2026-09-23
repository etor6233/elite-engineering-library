// Package qr provides a versioned reference for products, orders, customers
// and appointments. Its unkeyed checksum detects accidental changes, not
// forgery. Decode grants no authority. Verify requires the trusted session
// tenant and an actor/operation-bound object authorization resolver.
// This package contains no image encoder, camera or barcode decoder.
package qr

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Kind is the bounded set of QR target kinds.
type Kind string

const (
	KindProduct     Kind = "product"
	KindOrder       Kind = "order"
	KindCustomer    Kind = "customer"
	KindAppointment Kind = "appointment"
)

var (
	idRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

	ErrInvalidPayload = errors.New("qr: invalid payload")
	ErrChecksum       = errors.New("qr: checksum mismatch")
)

// Payload is a versioned, tenant-scoped QR identity.
type Payload struct {
	Version  int    `json:"v"`
	TenantID string `json:"t"`
	Kind     Kind   `json:"k"`
	ID       string `json:"id"`
	Checksum string `json:"c"`
}

// checksum computes an unkeyed digest. Anyone can recompute it; it is not a MAC.
func checksum(tenant string, kind Kind, id string) string {
	h := sha256.Sum256([]byte(tenant + "\x00" + string(kind) + "\x00" + id))
	return hex.EncodeToString(h[:])
}

// NewPayload builds an unsigned reference with a checksum, not a signed token.
func NewPayload(tenant string, kind Kind, id string) (Payload, error) {
	p := Payload{Version: 1, TenantID: tenant, Kind: kind, ID: id}
	if err := p.ValidateFields(); err != nil {
		return Payload{}, err
	}
	p.Checksum = checksum(tenant, kind, id)
	return p, nil
}

// ValidateFields checks identity fields (not the checksum).
func (p Payload) ValidateFields() error {
	if p.Version != 1 {
		return fmt.Errorf("%w: version", ErrInvalidPayload)
	}
	if strings.TrimSpace(p.TenantID) == "" || len(p.TenantID) > 64 {
		return fmt.Errorf("%w: tenant", ErrInvalidPayload)
	}
	switch p.Kind {
	case KindProduct, KindOrder, KindCustomer, KindAppointment:
	default:
		return fmt.Errorf("%w: kind", ErrInvalidPayload)
	}
	if !idRe.MatchString(p.ID) {
		return fmt.Errorf("%w: id", ErrInvalidPayload)
	}
	return nil
}

// Validate checks the checksum too.
func (p Payload) Validate() error {
	if err := p.ValidateFields(); err != nil {
		return err
	}
	if p.Checksum != checksum(p.TenantID, p.Kind, p.ID) {
		return ErrChecksum
	}
	return nil
}

// Encode returns the canonical JSON string for the QR.
func (p Payload) Encode() (string, error) {
	if err := p.Validate(); err != nil {
		return "", err
	}
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Decode parses fields and checksum only; success does not authenticate a sender.
func Decode(encoded string) (Payload, error) {
	var p Payload
	if err := json.Unmarshal([]byte(encoded), &p); err != nil {
		return Payload{}, fmt.Errorf("%w: %v", ErrInvalidPayload, err)
	}
	if err := p.Validate(); err != nil {
		return Payload{}, err
	}
	return p, nil
}
