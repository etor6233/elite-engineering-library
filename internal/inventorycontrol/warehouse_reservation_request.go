package inventorycontrol

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
)

var ErrWarehouseRequestNotFound = errors.New("warehouse request not found")

type WarehouseReservationRequest struct {
	RequestID              string     `json:"request_id"`
	OrganizationID         string     `json:"organization_id"`
	BinID                  string     `json:"bin_id"`
	ItemID                 string     `json:"item_id"`
	LotID                  string     `json:"lot_id,omitempty"`
	DemandKind             string     `json:"demand_kind"`
	DemandID               string     `json:"demand_id"`
	DemandLineID           string     `json:"demand_line_id"`
	Quantity               string     `json:"quantity"`
	CancellationDisallowed bool       `json:"cancellation_disallowed"`
	ExpiresAt              *time.Time `json:"expires_at,omitempty"`
}
type WarehouseReservationReceipt struct {
	Schema        string          `json:"schema"`
	RequestID     string          `json:"request_id"`
	PayloadSHA256 string          `json:"payload_sha256"`
	Reservation   BulkReservation `json:"reservation"`
}
type WarehouseReservationStore interface {
	RequestWarehouseReservation(context.Context, string, string, WarehouseReservationRequest) (WarehouseReservationReceipt, error)
	ReadWarehouseReservationRequest(context.Context, string, string, string, string) (WarehouseReservationReceipt, error)
}

func (v WarehouseReservationRequest) Validate() error {
	if v.RequestID == "" || v.OrganizationID == "" || v.BinID == "" || v.ItemID == "" || v.DemandID == "" || v.DemandLineID == "" || !positiveDecimal(v.Quantity, quantityPattern) {
		return ErrWarehouseWorkspaceQuery
	}
	if v.DemandKind != "manual" && v.DemandKind != "service" && v.DemandKind != "transfer-outbound" {
		return ErrWarehouseWorkspaceQuery
	}
	for _, s := range []string{v.RequestID, v.OrganizationID, v.BinID, v.ItemID, v.LotID, v.DemandID, v.DemandLineID} {
		if !safeWorkspaceText(s, 128) {
			return ErrWarehouseWorkspaceQuery
		}
	}
	return nil
}
func (v WarehouseReservationRequest) Hash() string {
	b, _ := json.Marshal(v)
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
