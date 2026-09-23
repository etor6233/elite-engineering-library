package whatsappbridge

// AUTHORED recovery of an existing local provider acceptance receipt. It never
// performs POST, invents a provider messageID or releases an unproved unknown.
import (
	"context"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"io"
	"os"
	"time"
)

func recoverEvidence(root *os.Root, path string) ([]byte, error) {
	f, err := root.Open(path)
	if err != nil {
		return nil, ErrBridge
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() < 1 || info.Size() > 65536 {
		return nil, ErrBridge
	}
	raw, err := io.ReadAll(io.LimitReader(f, 65537))
	if err != nil || len(raw) > 65536 {
		return nil, ErrBridge
	}
	return raw, nil
}
func (m *ReplyModule) Recover(ctx context.Context, p identity.Principal, key, hash string) (outbounddelivery.Receipt, error) {
	if m == nil || !m.approvals.validPrincipal(p, "whatsapp:send") {
		return outbounddelivery.Receipt{}, ErrApproval
	}
	value, err := m.approvals.Read(ctx, p, key)
	if err != nil || value.State != "approved" || value.PayloadSHA256 != hash {
		return outbounddelivery.Receipt{}, ErrApproval
	}
	if value.Status == nil || (value.Status.FenceState != "unknown" && value.Status.FenceState != "sending" && value.Status.FenceState != "accepted") {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	c := value.Context
	var approvalAt time.Time
	var reviewer string
	err = m.approvals.base.pool.QueryRow(ctx, `select r.decided_at,d.reviewer from approval.request r join approval.decision d on d.tenant_id=r.tenant_id and d.request_id=r.request_id and d.approved and d.reviewer<>r.requester where r.tenant_id=$1 and r.request_id=$2 and r.organization_id=$3 and r.kind='whatsapp_reply' and r.state='approved' and (select count(*) from approval.decision d2 where d2.tenant_id=r.tenant_id and d2.request_id=r.request_id)=1`, p.TenantID, key, m.approvals.organization).Scan(&approvalAt, &reviewer)
	if err != nil || reviewer == "" || c.ConnectionID != m.approvals.connection || c.ProfileSHA256 != m.approvals.profile {
		return outbounddelivery.Receipt{}, ErrApproval
	}
	return recoverBoundProviderReceipt(ctx, m.sender, m.store, m.approvals.base.pool, p.TenantID, key, boundReceiptContext{c.Message, c.MessageSHA256, c.ExpiresAt}, approvalAt, value.Status.FenceState)
}
