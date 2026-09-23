package merchantbridge

import (
	"testing"
	"time"
)

func TestMerchantQueueFreshnessAndApproval(t *testing.T) {
	now := time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC)
	past := now.Add(-time.Second)
	future := now.Add(time.Second)
	cases := []struct {
		name    string
		item    QueueItem
		expires time.Time
		changed bool
		want    string
	}{
		{"initial", QueueItem{}, future, false, "INITIAL_APPROVAL_REQUIRED"},
		{"ambiguous-prior-profile", QueueItem{UnresolvedApprovalID: "old"}, past, true, "RECONCILE"},
		{"review", QueueItem{ApprovalID: "a", ApprovalState: "pending"}, future, false, "AWAITING_REVIEW"},
		{"expired", QueueItem{ApprovalID: "a", ApprovalState: "approved"}, now, false, "EXPIRED_APPROVAL_REQUIRED"},
		{"approved", QueueItem{ApprovalID: "a", ApprovalState: "approved"}, future, false, "SEND_APPROVED"},
		{"source-changed", QueueItem{ApprovalID: "a", DeliveryState: "accepted"}, future, true, "CHANGED_SOURCE_APPROVAL_REQUIRED"},
		{"get-is-not-refresh", QueueItem{ApprovalID: "a", DeliveryState: "accepted"}, future, false, "FRESHNESS_UNCONFIRMED"},
		{"refresh-deadline", QueueItem{ApprovalID: "a", DeliveryState: "accepted", RefreshDueAt: &now}, past, false, "REFRESH_APPROVAL_REQUIRED"},
		{"current", QueueItem{ApprovalID: "a", DeliveryState: "accepted", RefreshDueAt: &future}, past, false, "CURRENT_INPUT"},
		{"provider-reject", QueueItem{ApprovalID: "a", DeliveryState: "failed_terminal", RefreshDueAt: &future}, future, false, "PROVIDER_REJECTED"},
		{"review-reject", QueueItem{ApprovalID: "a", ApprovalState: "rejected"}, future, false, "REVIEW_REJECTED"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := QueueAction(now, c.item, c.expires, c.changed); got != c.want {
				t.Fatalf("%s != %s", got, c.want)
			}
		})
	}
}
