package postgres

// AUTHORED composition of immutable catalog, existing ATP, human approvals and
// the shared outbound fence. No new monetary/inventory algorithm or ledger.
import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/channels"
	mb "elite.local/enterprise/internal/marketplacebridge"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
)

type MarketplacePublication struct {
	catalog   *CatalogRelease
	profile   mb.Profile
	approvals *HumanApprovals
	fence     *OutboundDeliveryStore
	client    *mb.Client
}

func NewMarketplacePublication(catalog *CatalogRelease, p mb.Profile, key []byte, tokens mb.TokenSource, http mb.Doer) (*MarketplacePublication, error) {
	if catalog == nil || p.Validate() != nil || p.TenantID != catalog.profile.TenantID || p.OrganizationID != catalog.profile.OrganizationID || p.Currency != catalog.profile.Currency {
		return nil, mb.ErrBinding
	}
	raw, _ := json.Marshal(p)
	var frozen mb.Profile
	if json.Unmarshal(raw, &frozen) != nil {
		return nil, mb.ErrBinding
	}
	client, e := mb.NewClient(frozen, tokens, http)
	if e != nil {
		return nil, e
	}
	fence, e := NewOutboundDeliveryStore(catalog.pool, key, 30*time.Second)
	if e != nil {
		return nil, e
	}
	return &MarketplacePublication{catalog, frozen, NewHumanApprovals(catalog.pool), fence, client}, nil
}
func (s *MarketplacePublication) authorized(p identity.Principal, permission string) bool {
	return s != nil && approvalPrincipal(p, s.profile.TenantID, s.profile.OrganizationID, permission)
}
func (s *MarketplacePublication) source(ctx context.Context, tx pgx.Tx, generation int64, variant string, horizon time.Time) (cr.PublicDocument, int64, error) {
	var publication cr.Publication
	var quantity int64
	if _, e := tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, "catalog-release:"+s.profile.TenantID); e != nil {
		return cr.PublicDocument{}, 0, e
	}
	e := tx.QueryRow(ctx, `select p.generation,d.snapshot_sha256,d.snapshot_canonical,p.effective_price_book_id from catalog.release_publication p
 join catalog.release_draft d using(tenant_id,draft_id)
 join catalog.release_public_current c on c.tenant_id=p.tenant_id and c.generation=p.generation
 join org.organization o on o.tenant_id=p.tenant_id and o.organization_id=$3
 join platform.tenant t on t.tenant_id=p.tenant_id
 where p.tenant_id=$1 and p.generation=$2 and o.status='active' and t.status='active'
 for share of o,t`, s.profile.TenantID, generation, s.profile.OrganizationID).Scan(&publication.Generation, &publication.SHA256, &publication.Snapshot, &publication.EffectivePriceBookID)
	if e != nil {
		return cr.PublicDocument{}, 0, mb.ErrBinding
	}
	e = tx.QueryRow(ctx, `select available_to_promise from inventory.serial_atp($1,$2,$3,$4)`, s.profile.TenantID, s.profile.OrganizationID, variant, horizon).Scan(&quantity)
	if e != nil || quantity < 0 {
		return cr.PublicDocument{}, 0, mb.ErrBinding
	}
	doc, e := cr.ProjectPublic(publication)
	return doc, quantity, e
}
func (s *MarketplacePublication) Prepare(ctx context.Context, p identity.Principal, r mb.PrepareRequest) (mb.Intent, error) {
	if !s.authorized(p, "marketplace:request") || r.Validate() != nil {
		return mb.Intent{}, mb.ErrBinding
	}
	// Remote preflight is read-only. Source publication and ATP remain database owners.
	before, e := s.initialObservation(ctx, r)
	if r.Operation != "MEDIA" && r.Operation != "CREATE" {
		before, e = s.client.Observe(ctx, r.VariantID, r.ItemID)
	}
	if e != nil {
		return mb.Intent{}, e
	}
	if r.Operation == "CONTENT" {
		if e = s.client.SingleUnsoldItem(ctx, before.Item); e != nil {
			return mb.Intent{}, e
		}
	}
	tx, e := s.catalog.pool.Begin(ctx)
	if e != nil {
		return mb.Intent{}, e
	}
	defer tx.Rollback(ctx)
	var now time.Time
	if e = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); e != nil {
		return mb.Intent{}, e
	}
	if !r.ExpiresAt.After(now) || r.ExpiresAt.After(now.Add(15*time.Minute)) {
		return mb.Intent{}, mb.ErrBinding
	}
	doc, quantity, e := s.source(ctx, tx, r.Generation, r.VariantID, now)
	if e != nil {
		return mb.Intent{}, e
	}
	var pictures []string
	if r.Operation == "CONTENT" {
		pic, e := s.acceptedMedia(ctx, tx, r)
		if e != nil {
			return mb.Intent{}, e
		}
		pictures = []string{pic}
	}
	out, e := mb.Build(s.profile, doc, r, quantity, now, before, pictures...)
	if e != nil {
		return out, e
	}
	return out, tx.Commit(ctx)
}
func (s *MarketplacePublication) guard(ctx context.Context, tx pgx.Tx, r mb.Intent, requireCurrent bool) error {
	if r.TenantID != s.profile.TenantID || r.OrganizationID != s.profile.OrganizationID || r.ProfileSHA256 != s.profile.SHA256() {
		return mb.ErrBinding
	}
	if !requireCurrent {
		return nil
	}
	var now time.Time
	if e := tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); e != nil {
		return e
	}
	if !r.ExpiresAt.After(now) || r.ExpiresAt.After(r.StockHorizon.Add(15*time.Minute)) || r.StockHorizon.After(now) {
		return mb.ErrBinding
	}
	doc, q, e := s.source(ctx, tx, r.Generation, r.VariantID, r.StockHorizon)
	if e != nil {
		return e
	}
	// Rebuild all source-derived fields; use the recorded observation only as a
	// before-state binding. The live client separately compares that observation.
	before := mb.Observation{Item: r.Expected, Exists: true, StockVersion: r.StockVersion}
	if r.Operation == "MEDIA" {
		before = mb.Observation{}
	}
	if r.Operation == "CREATE" {
		pic, e := s.acceptedMedia(ctx, tx, r.PrepareRequest)
		if e != nil {
			return e
		}
		before = mb.Observation{Item: mb.Item{Pictures: []mb.Picture{{ID: pic}}}}
	}
	var pictures []string
	if r.Operation == "CONTENT" {
		pic, e := s.acceptedMedia(ctx, tx, r.PrepareRequest)
		if e != nil {
			return e
		}
		pictures = []string{pic}
		zero := int64(0)
		before.Item.SoldQuantity = &zero
	}
	rebuilt, e := mb.Build(s.profile, doc, r.PrepareRequest, q, r.StockHorizon, before, pictures...)
	if e != nil {
		return e
	}
	rebuilt.BeforeSHA256 = r.BeforeSHA256
	_, a, e := cr.Canonical(rebuilt)
	if e != nil {
		return e
	}
	_, b, e := cr.Canonical(r)
	if e != nil || a != b {
		return mb.ErrBinding
	}
	return nil
}
func (s *MarketplacePublication) Submit(ctx context.Context, p identity.Principal, r mb.Intent) (string, bool, error) {
	if !s.authorized(p, "marketplace:request") || !cr.ValidSHA(r.BeforeSHA256) {
		return "", false, mb.ErrBinding
	}
	raw, hash, e := cr.Canonical(r)
	if e != nil {
		return "", false, e
	}
	spec := HumanApprovalSpec{Request: approval.Request{TenantID: p.TenantID, ID: r.ApprovalID, Kind: approval.KindMarketplaceMutation, SubjectID: r.SKU, Requester: p.Subject, EvidenceSHA: hash}, OrganizationID: s.profile.OrganizationID, Payload: raw}
	replay, e := s.approvals.Submit(ctx, p, spec, "marketplace:request", func(ctx context.Context, tx pgx.Tx) error { return s.guard(ctx, tx, r, true) })
	return hash, replay, e
}
func (s *MarketplacePublication) request(ctx context.Context, id string) (mb.Intent, string, error) {
	var r mb.Intent
	var raw []byte
	var hash string
	e := s.catalog.pool.QueryRow(ctx, `select payload,evidence_sha from approval.request where tenant_id=$1 and request_id=$2 and organization_id=$3 and kind='marketplace_mutation'`, s.profile.TenantID, id, s.profile.OrganizationID).Scan(&raw, &hash)
	if e != nil || mb.Decode(raw, &r) != nil || r.ApprovalID != id || r.ProfileSHA256 != s.profile.SHA256() {
		return r, "", mb.ErrBinding
	}
	_, h, e := cr.Canonical(r)
	if e != nil || h != hash {
		return r, "", mb.ErrBinding
	}
	return r, hash, nil
}
func (s *MarketplacePublication) Decide(ctx context.Context, p identity.Principal, id, hash string, approve bool, reason string) (approval.State, error) {
	if !s.authorized(p, "marketplace:approve") {
		return "", mb.ErrBinding
	}
	r, stored, e := s.request(ctx, id)
	if e != nil || stored != hash {
		return "", mb.ErrBinding
	}
	return s.approvals.Decide(ctx, p, s.profile.TenantID, id, s.profile.OrganizationID, hash, approve, reason, "marketplace:approve", func(ctx context.Context, tx pgx.Tx) error {
		if !approve {
			return nil
		}
		return s.guard(ctx, tx, r, true)
	})
}

type MarketplaceStatus struct {
	ApprovalID     string          `json:"approval_id"`
	ApprovalState  string          `json:"approval_state"`
	RequestSHA256  string          `json:"request_sha256"`
	Payload        json.RawMessage `json:"payload"`
	DeliveryState  string          `json:"delivery_state"`
	FailureCode    string          `json:"failure_code"`
	EvidenceSHA256 string          `json:"evidence_sha256"`
}

func (s *MarketplacePublication) Status(ctx context.Context, p identity.Principal, id string) (MarketplaceStatus, error) {
	var out MarketplaceStatus
	if !s.authorized(p, "marketplace:read") || !cr.ValidID(id) {
		return out, mb.ErrBinding
	}
	e := s.catalog.pool.QueryRow(ctx, `select a.request_id,a.state,a.evidence_sha,a.payload,coalesce(d.state,''),coalesce(d.failure_code,''),coalesce(d.evidence_sha256_hex,'')
 from approval.request a left join communication.outbound_delivery d on d.tenant_id=a.tenant_id and d.channel_code='mercadolibre_catalog' and d.delivery_key=a.request_id
 where a.tenant_id=$1 and a.request_id=$2 and a.organization_id=$3 and a.kind='marketplace_mutation' and a.payload->>'profile_sha256'=$4`, s.profile.TenantID, id, s.profile.OrganizationID, s.profile.SHA256()).Scan(&out.ApprovalID, &out.ApprovalState, &out.RequestSHA256, &out.Payload, &out.DeliveryState, &out.FailureCode, &out.EvidenceSHA256)
	return out, e
}

type marketplaceBound struct {
	service *MarketplacePublication
	intent  mb.Intent
	hash    string
}

func (b *marketplaceBound) Claim(ctx context.Context, m channels.Message, h string) (outbounddelivery.Claim, error) {
	return b.service.fence.claimWithAdmission(ctx, m, h, func(ctx context.Context, tx pgx.Tx) error {
		r := b.intent
		s := b.service
		// One active effect per tenant/provider SKU, including another approval id.
		if _, e := tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, "marketplace:"+r.TenantID+":"+r.SellerID+":"+r.SKU); e != nil {
			return e
		}
		var busy bool
		e := tx.QueryRow(ctx, `select exists(select 1 from catalog.marketplace_effect i join communication.outbound_delivery d using(tenant_id,channel_code,delivery_key)
    where i.tenant_id=$1 and i.seller_id=$2 and i.sku=$3 and (d.state in ('sending','unknown') or (d.state='accepted' and i.operation='CREATE' and $4='CREATE')) and i.operation<>'MEDIA' and $4<>'MEDIA')`, r.TenantID, r.SellerID, r.SKU, r.Operation).Scan(&busy)
		if e != nil {
			return e
		}
		if busy {
			return mb.ErrUnknown
		}
		if e = s.guard(ctx, tx, r, true); e != nil {
			return e
		}
		if _, e = LookupHumanApproval(ctx, tx, r.TenantID, r.ApprovalID, r.OrganizationID, approval.KindMarketplaceMutation, b.hash); e != nil {
			return e
		}
		_, e = tx.Exec(ctx, `insert into catalog.marketplace_effect(tenant_id,channel_code,delivery_key,seller_id,sku,generation,source_sha256,profile_sha256,request_sha256,operation) values($1,'mercadolibre_catalog',$2,$3,$4,$5,$6,$7,$8,$9)`, r.TenantID, r.ApprovalID, r.SellerID, r.SKU, r.Generation, r.SourceSHA256, r.ProfileSHA256, b.hash, r.Operation)
		return e
	})
}
func (b *marketplaceBound) Complete(c context.Context, m channels.Message, h string, r outbounddelivery.Receipt) error {
	return b.service.fence.Complete(c, m, h, r)
}
func (b *marketplaceBound) MarkUnknown(c context.Context, m channels.Message, h, code string) error {
	return b.service.fence.MarkUnknown(c, m, h, code)
}
func (b *marketplaceBound) MarkFailed(c context.Context, m channels.Message, h, sha, code string) error {
	return b.service.fence.MarkFailed(c, m, h, sha, code)
}

type marketplaceSender struct {
	client  *mb.Client
	intent  mb.Intent
	service *MarketplacePublication
}

func (s marketplaceSender) SendWithReceipt(ctx context.Context, m channels.Message) (outbounddelivery.Receipt, error) {
	want, _ := outbounddelivery.MessageSHA256(s.intent.Message())
	got, e := outbounddelivery.MessageSHA256(m)
	if e != nil || got != want {
		return outbounddelivery.Receipt{}, mb.ErrBinding
	}
	if s.intent.Operation == "MEDIA" || s.intent.Operation == "CREATE" {
		return s.service.sendInitial(ctx, s.intent)
	}
	return s.client.Write(ctx, s.intent)
}

type marketplaceOutbound struct{}

func (marketplaceOutbound) Code() string { return "mercadolibre_catalog" }
func (marketplaceOutbound) Receive(context.Context) ([]channels.Message, error) {
	return nil, mb.ErrBinding
}
func (marketplaceOutbound) Send(context.Context, channels.Message) error { return mb.ErrBinding }
func (s *MarketplacePublication) Send(ctx context.Context, p identity.Principal, id string) error {
	if !s.authorized(p, "marketplace:send") {
		return mb.ErrBinding
	}
	r, h, e := s.request(ctx, id)
	if e != nil {
		return e
	}
	channel := outbounddelivery.Channel{CodeValue: "mercadolibre_catalog", Receiver: marketplaceOutbound{}, Sender: marketplaceSender{s.client, r, s}, Store: &marketplaceBound{s, r, h}}
	return channel.Send(ctx, r.Message())
}
func (s *MarketplacePublication) Reconcile(ctx context.Context, p identity.Principal, id string) error {
	if !s.authorized(p, "marketplace:reconcile") {
		return mb.ErrBinding
	}
	r, _, e := s.request(ctx, id)
	if e != nil {
		return e
	}
	status, e := s.Status(ctx, p, id)
	if e != nil {
		return e
	}
	if status.DeliveryState == "accepted" {
		return nil
	}
	if status.DeliveryState != "unknown" && status.DeliveryState != "sending" {
		return mb.ErrBinding
	}
	message := r.Message()
	hash, e := outbounddelivery.MessageSHA256(message)
	if e != nil {
		return e
	}
	// Existing lease expiry makes a crashed sender unknown. No new claim is made.
	if _, e = s.fence.Claim(ctx, message, hash); !errors.Is(e, outbounddelivery.ErrUnknown) {
		return mb.ErrUnknown
	}
	receipt, e := s.reconcileInitial(ctx, r)
	if r.Operation != "MEDIA" && r.Operation != "CREATE" {
		receipt, e = s.client.Reconcile(ctx, r)
	}
	if e != nil {
		return e
	}
	return s.fence.ReconcileAccepted(ctx, message, hash, receipt)
}
