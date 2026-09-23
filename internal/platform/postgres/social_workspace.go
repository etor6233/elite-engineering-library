package postgres

// AUTHORED selection/read extension of the existing social approval owner.
// No provider call, permission grant, content policy or publication occurs here.
import (
	"context"
	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/socialbridge"
	"encoding/json"
	"time"
)

type SocialDraft struct {
	ID          string    `json:"approval_id"`
	Operation   string    `json:"operation"`
	Message     string    `json:"message"`
	OriginalID  string    `json:"original_approval_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}
type SocialPrepared struct {
	Request socialbridge.Request `json:"request"`
	Hash    string               `json:"request_sha256"`
}
type SocialWorkspaceRow struct {
	ID              string    `json:"approval_id"`
	Message         string    `json:"message"`
	Operation       string    `json:"operation"`
	RequestedBySelf bool      `json:"requested_by_self"`
	State           string    `json:"approval_state"`
	ScheduledAt     time.Time `json:"scheduled_at"`
	ExpiresAt       time.Time `json:"expires_at"`
}
type SocialWorkspacePage struct {
	OrganizationID string               `json:"organization_id"`
	Items          []SocialWorkspaceRow `json:"items"`
	NextCursor     string               `json:"next_cursor,omitempty"`
	Publish        bool                 `json:"publish_enabled"`
	Revoke         bool                 `json:"revoke_enabled"`
}
type SocialWorkspaceDetail struct {
	OrganizationID  string       `json:"organization_id"`
	RequestedBySelf bool         `json:"requested_by_self"`
	Status          SocialStatus `json:"status"`
}

// The durable approval is immutable; each projection nevertheless validates the
// binding and canonical hash before displaying any of its content.
func (s *SocialPublishing) verifiedWorkspaceRequest(id string, raw []byte, hash string) (socialbridge.Request, error) {
	var r socialbridge.Request
	if socialbridge.Decode(raw, &r) != nil || s.profile.Validate(r) != nil || r.Intent.ApprovalID != id {
		return r, socialbridge.ErrBinding
	}
	_, actual, e := approval.CanonicalPayload(raw)
	if e != nil || actual != hash {
		return r, socialbridge.ErrBinding
	}
	return r, nil
}
func (s *SocialPublishing) WorkspaceStatus(ctx context.Context, p identity.Principal, id string) (SocialWorkspaceDetail, error) {
	var out SocialWorkspaceDetail
	status, e := s.Status(ctx, p, id)
	if e != nil {
		return out, e
	}
	c := s.profile.Config()
	var mine bool
	e = s.pool.QueryRow(ctx, `select requester=$4 from approval.request where tenant_id=$1 and organization_id=$2 and request_id=$3`, c.TenantID, c.OrganizationID, id, p.Subject).Scan(&mine)
	if e != nil {
		return out, e
	}
	return SocialWorkspaceDetail{c.OrganizationID, mine, status}, nil
}
func (s *SocialPublishing) Workspace(ctx context.Context, p identity.Principal, after string) (SocialWorkspacePage, error) {
	c := s.profile.Config()
	out := SocialWorkspacePage{OrganizationID: c.OrganizationID, Items: []SocialWorkspaceRow{}, Publish: c.Publish, Revoke: c.Revoke}
	if !s.profile.Authorize(p, "social:read") || after != "" && !cr.ValidID(after) {
		return out, socialbridge.ErrBinding
	}
	rows, e := s.pool.Query(ctx, `select request_id,payload,evidence_sha,requester=$6,state
 from approval.request where tenant_id=$1 and organization_id=$2 and kind in('social_publish','social_revoke')
 and payload#>>'{intent,profile_sha256}'=$3 and payload#>>'{intent,page_id}'=$4 and request_id>$5 order by request_id limit 26`, c.TenantID, c.OrganizationID, s.profile.SHA256(), c.PageID, after, p.Subject)
	if e != nil {
		return out, e
	}
	defer rows.Close()
	for rows.Next() {
		var row SocialWorkspaceRow
		var raw []byte
		var hash string
		if e = rows.Scan(&row.ID, &raw, &hash, &row.RequestedBySelf, &row.State); e != nil {
			return out, e
		}
		r, e := s.verifiedWorkspaceRequest(row.ID, raw, hash)
		if e != nil {
			return SocialWorkspacePage{}, e
		}
		row.Message = r.Intent.Message
		row.Operation = r.Intent.Operation
		row.ScheduledAt = r.ScheduledAt
		row.ExpiresAt = r.ExpiresAt
		out.Items = append(out.Items, row)
	}
	if e = rows.Err(); e != nil {
		return out, e
	}
	if len(out.Items) > 25 {
		out.NextCursor = out.Items[24].ID
		out.Items = out.Items[:25]
	}
	return out, nil
}
func (s *SocialPublishing) PrepareDraft(ctx context.Context, p identity.Principal, in SocialDraft) (SocialPrepared, error) {
	out := SocialPrepared{}
	if !s.profile.Authorize(p, "social:request") || !cr.ValidID(in.ID) || len(in.ID) < 16 || in.ExpiresAt.IsZero() || !in.ExpiresAt.After(in.ScheduledAt) {
		return out, socialbridge.ErrBinding
	}
	c := s.profile.Config()
	r := socialbridge.Request{Intent: socialbridge.Intent{TenantID: c.TenantID, PageID: c.PageID, ProfileSHA256: s.profile.SHA256(), ApprovalID: in.ID, Operation: in.Operation}, ScheduledAt: in.ScheduledAt, ExpiresAt: in.ExpiresAt}
	switch in.Operation {
	case "publish":
		if in.OriginalID != "" {
			return out, socialbridge.ErrBinding
		}
		r.Intent.Message = in.Message
		r.Intent.ContentSHA256 = socialbridge.Hash([]byte(in.Message))
	case "revoke":
		if in.Message != "" || !cr.ValidID(in.OriginalID) {
			return out, socialbridge.ErrBinding
		}
		original, e := s.Status(ctx, p, in.OriginalID)
		if e != nil {
			return out, e
		}
		var req socialbridge.Request
		var receipt socialbridge.Receipt
		if socialbridge.Decode(original.Payload, &req) != nil || socialbridge.Decode(original.Receipt, &receipt) != nil || req.Intent.Operation != "publish" || receipt.Validate(req.Intent) != nil {
			return out, socialbridge.ErrBinding
		}
		r.Intent.Message = req.Intent.Message
		r.Intent.ContentSHA256 = req.Intent.ContentSHA256
		r.OriginalApprovalID = in.OriginalID
		r.ProviderReference = receipt.ProviderReference
	default:
		return out, socialbridge.ErrBinding
	}
	r.Intent.DeliveryKey = socialbridge.DeliveryKey(r.Intent)
	if s.profile.Validate(r) != nil {
		return out, socialbridge.ErrBinding
	}
	tx, e := s.pool.Begin(ctx)
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	if e = s.guard(ctx, tx, r, false); e != nil {
		return out, e
	}
	raw, e := json.Marshal(r)
	if e != nil {
		return out, e
	}
	_, h, e := approval.CanonicalPayload(raw)
	if e != nil {
		return out, e
	}
	return SocialPrepared{Request: r, Hash: h}, tx.Commit(ctx)
}
