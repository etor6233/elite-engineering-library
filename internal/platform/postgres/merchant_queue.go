package postgres

// AUTHORED bounded read-only projection. No queue read creates approvals,
// provider writes or jobs; each refresh follows the existing manual workflow.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	mb "elite.local/enterprise/internal/merchantbridge"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"
)

type MerchantQueue struct {
	ObservedAt    time.Time      `json:"observed_at"`
	ProfileSHA256 string         `json:"profile_sha256"`
	Items         []mb.QueueItem `json:"items"`
}

func (s *MerchantPublication) Queue(ctx context.Context, p identity.Principal) (MerchantQueue, error) {
	var out MerchantQueue
	if !s.authorized(p, "merchant:read") {
		return out, mb.ErrBinding
	}
	tx, e := s.catalog.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	var generation int64
	if e = tx.QueryRow(ctx, `select clock_timestamp(),generation from catalog.release_public_current where tenant_id=$1`, s.profile.TenantID).Scan(&out.ObservedAt, &generation); e != nil {
		return out, mb.ErrBinding
	}
	out.ProfileSHA256 = s.profile.SHA256()
	out.Items = make([]mb.QueueItem, 0, len(s.profile.Bindings))
	for _, b := range s.profile.Bindings {
		horizon := out.ObservedAt.Add(24 * time.Hour)
		doc, q, e := s.source(ctx, tx, generation, b.VariantID, horizon)
		if e != nil {
			return out, e
		}
		desired, e := mb.Build(s.profile, doc, mb.PrepareRequest{ApprovalID: "queue-projection-only", Generation: generation, VariantID: b.VariantID, ExpiresAt: out.ObservedAt.Add(time.Minute)}, q, horizon)
		if e != nil {
			return out, e
		}
		_, desiredSHA, e := cr.Canonical(desired.Product)
		if e != nil {
			return out, e
		}
		item := mb.QueueItem{VariantID: b.VariantID, OfferID: b.OfferID, Generation: generation, SourceSHA256: doc.SourceSHA256, DesiredProductSHA256: desiredSHA}
		var raw []byte
		var latest mb.Intent
		e = tx.QueryRow(ctx, `select a.request_id,a.state,coalesce(d.state,''),a.payload from approval.request a
 left join communication.outbound_delivery d on d.tenant_id=a.tenant_id and d.delivery_key=a.request_id and d.channel_code='google_merchant'
 where a.tenant_id=$1 and a.organization_id=$2 and a.kind='marketplace_mutation' and a.payload->>'profile_sha256'=$3
 and a.payload->'product'->>'offer_id'=$4 order by a.created_at desc,a.request_id desc limit 1`,
			s.profile.TenantID, s.profile.OrganizationID, s.profile.SHA256(), b.OfferID).Scan(&item.ApprovalID, &item.ApprovalState, &item.DeliveryState, &raw)
		if e != nil && !errors.Is(e, pgx.ErrNoRows) {
			return out, e
		}
		changed := false
		if e == nil {
			if mb.Decode(raw, &latest) != nil {
				return out, mb.ErrBinding
			}
			_, h, e := cr.Canonical(latest.Product)
			if e != nil {
				return out, e
			}
			changed = latest.SourceSHA256 != doc.SourceSHA256 || latest.Generation != generation || h != desiredSHA
		}
		e = tx.QueryRow(ctx, `select coalesce((select i.delivery_key from catalog.merchant_effect i
 join communication.outbound_delivery d using(tenant_id,channel_code,delivery_key)
 where i.tenant_id=$1 and i.account_id=$2 and i.offer_id=$3 and d.state in('sending','unknown')
 order by i.delivery_key limit 1),''),
 (select max(o.observed_at)+make_interval(days=>$5) from catalog.merchant_observation o
 join catalog.merchant_effect i on i.tenant_id=o.tenant_id and i.delivery_key=o.approval_id
 where i.tenant_id=$1 and i.account_id=$2 and i.offer_id=$3 and i.profile_sha256=$4
 and o.operation='INSERT' and o.confirmed)`,
			s.profile.TenantID, s.profile.AccountID, b.OfferID, s.profile.SHA256(), s.profile.RefreshCadenceDays).Scan(&item.UnresolvedApprovalID, &item.RefreshDueAt)
		if e != nil {
			return out, e
		}
		item.Action = mb.QueueAction(out.ObservedAt, item, latest.ExpiresAt, changed)
		out.Items = append(out.Items, item)
	}
	return out, tx.Commit(ctx)
}
