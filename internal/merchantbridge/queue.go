package merchantbridge

// AUTHORED read-model classification over original approval/fence receipts.
import "time"

type QueueItem struct {
	VariantID            string     `json:"variant_id"`
	OfferID              string     `json:"offer_id"`
	Generation           int64      `json:"generation,string"`
	SourceSHA256         string     `json:"source_sha256"`
	DesiredProductSHA256 string     `json:"desired_product_sha256"`
	ApprovalID           string     `json:"approval_id"`
	ApprovalState        string     `json:"approval_state"`
	DeliveryState        string     `json:"delivery_state"`
	UnresolvedApprovalID string     `json:"unresolved_approval_id"`
	RefreshDueAt         *time.Time `json:"refresh_due_at"`
	Action               string     `json:"action"`
}

func QueueAction(now time.Time, item QueueItem, expires time.Time, changed bool) string {
	if item.UnresolvedApprovalID != "" {
		return "RECONCILE"
	}
	if item.ApprovalID == "" {
		return "INITIAL_APPROVAL_REQUIRED"
	}
	if item.DeliveryState == "failed_terminal" {
		return "PROVIDER_REJECTED"
	}
	if item.ApprovalState == "rejected" {
		return "REVIEW_REJECTED"
	}
	if changed {
		return "CHANGED_SOURCE_APPROVAL_REQUIRED"
	}
	if item.DeliveryState == "" {
		if !expires.After(now) {
			return "EXPIRED_APPROVAL_REQUIRED"
		}
		if item.ApprovalState == "approved" {
			return "SEND_APPROVED"
		}
		return "AWAITING_REVIEW"
	}
	if item.DeliveryState != "accepted" {
		return "RECONCILE"
	}
	if item.RefreshDueAt == nil {
		return "FRESHNESS_UNCONFIRMED"
	}
	if !item.RefreshDueAt.After(now) {
		return "REFRESH_APPROVAL_REQUIRED"
	}
	return "CURRENT_INPUT"
}
