package whatsappbridge

import (
	"context"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"time"
)

func (m *ScheduledNotifications) Recover(ctx context.Context, p identity.Principal, key, hash string) (outbounddelivery.Receipt, error) {
	if m == nil || !m.approvals.authorized(p, "notification:reconcile") {
		return outbounddelivery.Receipt{}, ErrApproval
	}
	value, e := m.approvals.Status(ctx, p, key)
	if e != nil || value.RequestSHA256 != hash || value.ApprovalState != "approved" || value.Notification == nil {
		return outbounddelivery.Receipt{}, ErrApproval
	}
	state := value.Notification.FenceState
	if state != "accepted" && state != "sending" && state != "unknown" {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	var approvedAt time.Time
	var reviewer string
	s := m.approvals
	e = s.base.pool.QueryRow(ctx, `select r.decided_at,d.reviewer from approval.request r join approval.decision d on d.tenant_id=r.tenant_id and d.request_id=r.request_id and d.approved and d.reviewer<>r.requester where r.tenant_id=$1 and r.request_id=$2 and r.organization_id=$3 and r.kind='whatsapp_schedule' and r.state='approved' and (select count(*) from approval.decision d2 where d2.tenant_id=r.tenant_id and d2.request_id=r.request_id)=1`, s.tenant, key, s.org).Scan(&approvedAt, &reviewer)
	if e != nil || reviewer == "" {
		return outbounddelivery.Receipt{}, ErrApproval
	}
	c := value.Context
	if c.NotBefore.After(approvedAt) {
		approvedAt = c.NotBefore
	}
	return recoverBoundProviderReceipt(ctx, m.sender, m.store, s.base.pool, s.tenant, key, boundReceiptContext{c.Message, c.MessageSHA256, c.ExpiresAt}, approvedAt, state)
}
