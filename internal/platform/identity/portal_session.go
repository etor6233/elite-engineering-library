package identity

import (
	"context"
	"errors"
	"time"
)

var ErrPortalSessionConflict = errors.New("portal session conflict")
var ErrPortalSessionNotFound = errors.New("portal session not found")

type PortalSessionCommand struct {
	SessionID         string    `json:"session_id"`
	OperationID       string    `json:"operation_id,omitempty"`
	Version           int64     `json:"version,omitempty"`
	Ciphertext        string    `json:"ciphertext,omitempty"`
	AccessExpiresAt   time.Time `json:"access_expires_at,omitempty"`
	AbsoluteExpiresAt time.Time `json:"absolute_expires_at,omitempty"`
}
type PortalSessionRecord struct {
	Version           int64     `json:"version"`
	State             string    `json:"state"`
	OperationID       string    `json:"operation_id"`
	Ciphertext        string    `json:"ciphertext"`
	AccessExpiresAt   time.Time `json:"access_expires_at"`
	AbsoluteExpiresAt time.Time `json:"absolute_expires_at"`
	RevocationPending bool      `json:"revocation_pending"`
}
type PortalSessionStore interface {
	Execute(context.Context, PortalProfile, string, PortalSessionCommand) (PortalSessionRecord, error)
	Sweep(context.Context, PortalProfile, string) (PortalSessionSweep, error)
}

type PortalSessionSweepItem struct {
	SessionIDSHA256 string `json:"session_id_sha256"`
	PortalSessionRecord
}
type PortalSessionSweep struct {
	Items             []PortalSessionSweepItem `json:"items"`
	Purged            int64                    `json:"purged"`
	UnconfirmedPurged int64                    `json:"unconfirmed_purged"`
}
