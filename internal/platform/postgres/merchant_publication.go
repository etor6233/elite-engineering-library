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
	mb "elite.local/enterprise/internal/merchantbridge"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
)

type MerchantPublication struct {
	catalog   *CatalogRelease
	profile   mb.Profile
	approvals *HumanApprovals
	fence     *OutboundDeliveryStore
	client    *mb.Client
}

func NewMerchantPublication(catalog *CatalogRelease, p mb.Profile, key []byte, executor mb.Executor) (*MerchantPublication, error) {
	if catalog == nil || p.Validate() != nil || p.TenantID != catalog.profile.TenantID || p.OrganizationID != catalog.profile.OrganizationID || p.Currency != catalog.profile.Currency {
		return nil, mb.ErrBinding
	}
	raw, _ := json.Marshal(p)
	var frozen mb.Profile
	if json.Unmarshal(raw, &frozen) != nil {
		return nil, mb.ErrBinding
	}
	client, e := mb.NewClient(frozen, executor)
	if e != nil {
		return nil, e
	}
	fence, e := NewOutboundDeliveryStore(catalog.pool, key, 30*time.Second)
	if e != nil {
		return nil, e
	}
	return &MerchantPublication{catalog, frozen, NewHumanApprovals(catalog.pool), fence, client}, nil
}
func (s *MerchantPublication) authorized(p identity.Principal, permission string) bool {
	return s != nil && approvalPrincipal(p, s.profile.TenantID, s.profile.OrganizationID, permission)
}
func (s *MerchantPublication) source(ctx context.Context, tx pgx.Tx, generation int64, variant string, horizon time.Time) (cr.PublicDocument, int64, error) {
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

func (s *MerchantPublication) Prepare(ctx context.Context, p identity.Principal, r mb.PrepareRequest) (mb.Intent, error) {
	if !s.authorized(p, "merchant:request") || r.Validate() != nil {
		return mb.Intent{}, mb.ErrBinding
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
	doc, q, e := s.source(ctx, tx, r.Generation, r.VariantID, now)
	if e != nil {
		return mb.Intent{}, e
	}
	intent, e := mb.Build(s.profile, doc, r, q, now)
	if e != nil {
		return mb.Intent{}, e
	}
	return intent, tx.Commit(ctx)
}
func (s *MerchantPublication) guard(ctx context.Context, tx pgx.Tx, r mb.Intent, requireCurrent bool) error {
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
	rebuilt, e := mb.Build(s.profile, doc, r.PrepareRequest, q, r.StockHorizon)
	if e != nil {
		return e
	}
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
func (s *MerchantPublication) Submit(ctx context.Context, p identity.Principal, r mb.Intent) (string, bool, error) {
	if !s.authorized(p, "merchant:request") {
		return "", false, mb.ErrBinding
	}
	raw, hash, e := cr.Canonical(r)
	if e != nil {
		return "", false, e
	}
	spec := HumanApprovalSpec{Request: approval.Request{TenantID: p.TenantID, ID: r.ApprovalID, Kind: approval.KindMarketplaceMutation, SubjectID: r.Product.OfferID, Requester: p.Subject, EvidenceSHA: hash}, OrganizationID: s.profile.OrganizationID, Payload: raw}
	replay, e := s.approvals.Submit(ctx, p, spec, "merchant:request", func(ctx context.Context, tx pgx.Tx) error { return s.guard(ctx, tx, r, true) })
	return hash, replay, e
}
func (s *MerchantPublication) request(ctx context.Context, id string) (mb.Intent, string, error) {
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
func (s *MerchantPublication) Decide(ctx context.Context, p identity.Principal, id, hash string, approve bool, reason string) (approval.State, error) {
	if !s.authorized(p, "merchant:approve") {
		return "", mb.ErrBinding
	}
	r, stored, e := s.request(ctx, id)
	if e != nil || stored != hash {
		return "", mb.ErrBinding
	}
	return s.approvals.Decide(ctx, p, s.profile.TenantID, id, s.profile.OrganizationID, hash, approve, reason, "merchant:approve", func(ctx context.Context, tx pgx.Tx) error {
		if !approve {
			return nil
		}
		return s.guard(ctx, tx, r, true)
	})
}

type MerchantStatus struct {
	InputAcknowledged   bool            `json:"input_acknowledged"`
	RefreshDueAt        *time.Time      `json:"refresh_due_at,omitempty"`
	ProcessingState     string          `json:"processing_state"`
	ProcessingConfirmed bool            `json:"processing_confirmed"`
	ProcessingResponse  json.RawMessage `json:"processing_response,omitempty"`
	ApprovalID          string          `json:"approval_id"`
	ApprovalState       string          `json:"approval_state"`
	RequestSHA256       string          `json:"request_sha256"`
	Payload             json.RawMessage `json:"payload"`
	DeliveryState       string          `json:"delivery_state"`
	FailureCode         string          `json:"failure_code"`
	EvidenceSHA256      string          `json:"evidence_sha256"`
}

func (s *MerchantPublication) Status(ctx context.Context, p identity.Principal, id string) (MerchantStatus, error) {
	var out MerchantStatus
	if !s.authorized(p, "merchant:read") || !cr.ValidID(id) {
		return out, mb.ErrBinding
	}
	e := s.catalog.pool.QueryRow(ctx, `select a.request_id,a.state,a.evidence_sha,a.payload,coalesce(d.state,''),coalesce(d.failure_code,''),coalesce(d.evidence_sha256_hex,'')
 from approval.request a left join communication.outbound_delivery d on d.tenant_id=a.tenant_id and d.channel_code='google_merchant' and d.delivery_key=a.request_id
 where a.tenant_id=$1 and a.request_id=$2 and a.organization_id=$3 and a.kind='marketplace_mutation' and a.payload->>'profile_sha256'=$4`, s.profile.TenantID, id, s.profile.OrganizationID, s.profile.SHA256()).Scan(&out.ApprovalID, &out.ApprovalState, &out.RequestSHA256, &out.Payload, &out.DeliveryState, &out.FailureCode, &out.EvidenceSHA256)
	if e != nil {
		return out, e
	}
	e = s.catalog.pool.QueryRow(ctx, `select exists(select 1 from catalog.merchant_observation where tenant_id=$1 and approval_id=$2 and operation='INSERT' and confirmed),
 (select max(observed_at)+make_interval(days=>$3) from catalog.merchant_observation where tenant_id=$1 and approval_id=$2 and operation='INSERT' and confirmed)`, s.profile.TenantID, id, s.profile.RefreshCadenceDays).Scan(&out.InputAcknowledged, &out.RefreshDueAt)
	if e != nil {
		return out, e
	}
	var raw []byte
	var hash string
	e = s.catalog.pool.QueryRow(ctx, `select result_state,confirmed,response,response_sha256 from catalog.merchant_observation
 where tenant_id=$1 and approval_id=$2 and operation='GET' order by observed_at desc,response_sha256 limit 1`, s.profile.TenantID, id).Scan(&out.ProcessingState, &out.ProcessingConfirmed, &raw, &hash)
	if errors.Is(e, pgx.ErrNoRows) {
		out.ProcessingState = "UNOBSERVED"
		return out, nil
	}
	if e != nil {
		return out, e
	}
	if mb.Hash(raw) != hash {
		return out, mb.ErrBinding
	}
	out.ProcessingResponse = raw
	return out, nil
}

type merchantBound struct {
	service *MerchantPublication
	intent  mb.Intent
	hash    string
}

func (b *merchantBound) Claim(ctx context.Context, m channels.Message, h string) (outbounddelivery.Claim, error) {
	return b.service.fence.claimWithAdmission(ctx, m, h, func(ctx context.Context, tx pgx.Tx) error {
		r := b.intent
		s := b.service
		// One active effect per tenant/provider SKU, including another approval id.
		if _, e := tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, "merchant:"+r.TenantID+":"+r.AccountID+":"+r.Product.OfferID); e != nil {
			return e
		}
		var busy bool

		e := tx.QueryRow(ctx, `select exists(select 1 from catalog.merchant_effect i
  join communication.outbound_delivery d using(tenant_id,channel_code,delivery_key)
  where i.tenant_id=$1 and i.account_id=$2 and i.offer_id=$3 and d.state in ('sending','unknown'))`, r.TenantID, r.AccountID, r.Product.OfferID).Scan(&busy)
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
		_, e = tx.Exec(ctx, `insert into catalog.merchant_effect(tenant_id,channel_code,delivery_key,account_id,offer_id,generation,source_sha256,profile_sha256,request_sha256) values($1,'google_merchant',$2,$3,$4,$5,$6,$7,$8)`, r.TenantID, r.ApprovalID, r.AccountID, r.Product.OfferID, r.Generation, r.SourceSHA256, r.ProfileSHA256, b.hash)
		return e
	})
}
func (b *merchantBound) Complete(c context.Context, m channels.Message, h string, r outbounddelivery.Receipt) error {
	return b.service.fence.Complete(c, m, h, r)
}
func (b *merchantBound) MarkUnknown(c context.Context, m channels.Message, h, code string) error {
	return b.service.fence.MarkUnknown(c, m, h, code)
}
func (b *merchantBound) MarkFailed(c context.Context, m channels.Message, h, sha, code string) error {
	return b.service.fence.MarkFailed(c, m, h, sha, code)
}

type merchantSender struct {
	client  *mb.Client
	intent  mb.Intent
	service *MerchantPublication
}

func (s merchantSender) SendWithReceipt(ctx context.Context, m channels.Message) (outbounddelivery.Receipt, error) {
	want, _ := outbounddelivery.MessageSHA256(s.intent.Message())
	got, e := outbounddelivery.MessageSHA256(m)
	if e != nil || got != want {
		return outbounddelivery.Receipt{}, mb.ErrBinding
	}

	return s.client.Insert(ctx, s.intent, func(ctx context.Context, result mb.Result) error {
		return s.service.recordMerchantObservation(ctx, s.intent, result, true)
	})
}

type merchantOutbound struct{}

func (merchantOutbound) Code() string { return "google_merchant" }
func (merchantOutbound) Receive(context.Context) ([]channels.Message, error) {
	return nil, mb.ErrBinding
}
func (merchantOutbound) Send(context.Context, channels.Message) error { return mb.ErrBinding }
func (s *MerchantPublication) Send(ctx context.Context, p identity.Principal, id string) error {
	if !s.authorized(p, "merchant:send") {
		return mb.ErrBinding
	}
	r, h, e := s.request(ctx, id)
	if e != nil {
		return e
	}
	channel := outbounddelivery.Channel{CodeValue: "google_merchant", Receiver: merchantOutbound{}, Sender: merchantSender{s.client, r, s}, Store: &merchantBound{s, r, h}}
	return channel.Send(ctx, r.Message())
}

func (s *MerchantPublication) Reconcile(ctx context.Context, p identity.Principal, id string) error {
	if !s.authorized(p, "merchant:reconcile") {
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
	if status.DeliveryState != "accepted" && status.DeliveryState != "unknown" && status.DeliveryState != "sending" {
		return mb.ErrBinding
	}
	message := r.Message()
	hash, e := outbounddelivery.MessageSHA256(message)
	if e != nil {
		return e
	}
	if status.DeliveryState != "accepted" {
		if _, e = s.fence.Claim(ctx, message, hash); !errors.Is(e, outbounddelivery.ErrUnknown) {
			return mb.ErrUnknown
		}
	}
	observation, observeErr := s.client.Observe(ctx, r)
	if e = s.recordMerchantObservation(ctx, r, observation, observeErr == nil); e != nil {
		return e
	}
	if observeErr != nil {
		return observeErr
	}
	if status.DeliveryState == "accepted" {
		return nil
	}
	receipt := outbounddelivery.Receipt{ProviderMessageID: r.Resource, EvidenceSHA256: observation.ResponseSHA256, AcceptedAt: time.Now().UTC()}
	return s.fence.ReconcileAccepted(ctx, message, hash, receipt)
}
