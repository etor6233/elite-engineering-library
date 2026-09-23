package postgres

// AUTHORED binding of the admitted catalog PNG, human approvals and outbound
// fence to provider observations. No credentials or raw bearer tokens persist.
import (
	"context"

	cr "elite.local/enterprise/internal/catalogrelease"
	mb "elite.local/enterprise/internal/marketplacebridge"
	"elite.local/enterprise/internal/outbounddelivery"
	"github.com/jackc/pgx/v5"
)

type marketplaceQuery interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func (s *MarketplacePublication) acceptedMedia(ctx context.Context, q marketplaceQuery, r mb.PrepareRequest) (string, error) {
	var payload, raw []byte
	var sha string
	e := q.QueryRow(ctx, `select a.payload,o.observation,o.observation_sha256
 from approval.request a join communication.outbound_delivery d on d.tenant_id=a.tenant_id and d.channel_code='mercadolibre_catalog' and d.delivery_key=a.request_id
 join catalog.marketplace_observation o on o.tenant_id=a.tenant_id and o.approval_id=a.request_id
 where a.tenant_id=$1 and a.request_id=$2 and a.organization_id=$3 and a.kind='marketplace_mutation' and a.state='approved' and d.state='accepted' and o.operation='MEDIA'`, s.profile.TenantID, r.MediaApprovalID, s.profile.OrganizationID).Scan(&payload, &raw, &sha)
	var media mb.Intent
	var observation mb.InitialObservation
	if e != nil || mb.Decode(payload, &media) != nil || mb.Decode(raw, &observation) != nil || media.Operation != "MEDIA" || media.ApprovalID != r.MediaApprovalID || media.Generation != r.Generation || media.VariantID != r.VariantID || media.ProfileSHA256 != s.profile.SHA256() || media.TenantID != s.profile.TenantID || media.OrganizationID != s.profile.OrganizationID || mb.ValidateUploadObservation(media, observation) != nil {
		return "", mb.ErrBinding
	}
	_, h, e := cr.Canonical(observation)
	if e != nil || h != sha {
		return "", mb.ErrBinding
	}
	return observation.ProviderID, nil
}
func (s *MarketplacePublication) initialObservation(ctx context.Context, r mb.PrepareRequest) (mb.Observation, error) {
	if r.Operation != "CREATE" {
		return mb.Observation{}, nil
	}
	pic, e := s.acceptedMedia(ctx, s.catalog.pool, r)
	if e != nil {
		return mb.Observation{}, e
	}
	if e = s.client.CreationPreflight(ctx, r.VariantID); e != nil {
		return mb.Observation{}, e
	}
	return mb.Observation{Item: mb.Item{Pictures: []mb.Picture{{ID: pic}}}}, nil
}
func (s *MarketplacePublication) recordInitial(ctx context.Context, r mb.Intent, o mb.InitialObservation) error {
	if o.Operation != r.Operation || o.MediaSHA256 != r.MediaSHA256 {
		return mb.ErrBinding
	}
	raw, h, e := cr.Canonical(o)
	if e != nil {
		return e
	}
	result, e := s.catalog.pool.Exec(ctx, `insert into catalog.marketplace_observation(tenant_id,approval_id,operation,provider_id,media_sha256,observation,observation_sha256)
 select i.tenant_id,i.delivery_key,$3,$4,$5,$6,$7 from catalog.marketplace_effect i
 join communication.outbound_delivery d using(tenant_id,channel_code,delivery_key)
 where i.tenant_id=$1 and i.delivery_key=$2 and i.channel_code='mercadolibre_catalog' and i.operation=$3
 and i.profile_sha256=$8 and i.source_sha256=$9 and d.state in ('sending','unknown')
 on conflict(tenant_id,approval_id) do nothing`, r.TenantID, r.ApprovalID, r.Operation, o.ProviderID, r.MediaSHA256, raw, h, r.ProfileSHA256, r.SourceSHA256)
	if e != nil {
		return e
	}
	if result.RowsAffected() == 1 {
		return nil
	}
	var stored string
	if e = s.catalog.pool.QueryRow(ctx, `select observation_sha256 from catalog.marketplace_observation where tenant_id=$1 and approval_id=$2`, r.TenantID, r.ApprovalID).Scan(&stored); e != nil || stored != h {
		return mb.ErrBinding
	}
	return nil
}
func (s *MarketplacePublication) sendInitial(ctx context.Context, r mb.Intent) (outbounddelivery.Receipt, error) {
	record := func(ctx context.Context, o mb.InitialObservation) error { return s.recordInitial(ctx, r, o) }
	if r.Operation == "CREATE" {
		return s.client.Create(ctx, r, record)
	}
	png, e := s.catalog.PublicMedia(ctx, r.MediaSHA256)
	if e != nil {
		return outbounddelivery.Receipt{}, e
	}
	return s.client.Upload(ctx, r, png, record)
}
func (s *MarketplacePublication) reconcileInitial(ctx context.Context, r mb.Intent) (outbounddelivery.Receipt, error) {
	var raw []byte
	var h string
	var o mb.InitialObservation
	e := s.catalog.pool.QueryRow(ctx, `select observation,observation_sha256 from catalog.marketplace_observation where tenant_id=$1 and approval_id=$2`, r.TenantID, r.ApprovalID).Scan(&raw, &h)
	if e == nil {
		if mb.Decode(raw, &o) != nil || o.Operation != r.Operation || o.MediaSHA256 != r.MediaSHA256 {
			return outbounddelivery.Receipt{}, mb.ErrBinding
		}
		_, hash, e := cr.Canonical(o)
		if e != nil || hash != h {
			return outbounddelivery.Receipt{}, mb.ErrBinding
		}
	} else if e != pgx.ErrNoRows {
		return outbounddelivery.Receipt{}, e
	}
	if r.Operation == "MEDIA" {
		if e != nil {
			return outbounddelivery.Receipt{}, mb.ErrUnknown
		}
		return s.client.AcceptedUpload(r, o)
	}
	if r.Operation == "CREATE" {
		return s.client.ReconcileCreation(ctx, r, o.ProviderID)
	}
	return outbounddelivery.Receipt{}, mb.ErrBinding
}
