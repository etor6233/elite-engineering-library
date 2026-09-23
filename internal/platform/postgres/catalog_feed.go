package postgres

// AUTHORED feed binding to the existing outbound delivery state machine/fence.
import (
	"context"
	"crypto/sha256"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type CatalogFeed struct {
	catalog  *CatalogRelease
	fence    *OutboundDeliveryStore
	receiver *cr.FeedHTTP
}

func NewCatalogFeed(catalog *CatalogRelease, hmacKey []byte, receiver *cr.FeedHTTP) (*CatalogFeed, error) {
	if catalog == nil || receiver == nil {
		return nil, cr.ErrInvalid
	}
	fence, e := NewOutboundDeliveryStore(catalog.pool, hmacKey, 15*time.Second)
	if e != nil {
		return nil, e
	}
	return &CatalogFeed{catalog, fence, receiver}, nil
}

type catalogFeedPlan struct {
	message               channels.Message
	payload               []byte
	payloadSHA, sourceSHA string
	generation            int64
}

func (s *CatalogFeed) plan(ctx context.Context, p identity.Principal, generation int64, permission string) (catalogFeedPlan, error) {
	var plan catalogFeedPlan
	if generation < 1 || !s.catalog.authorized(p, permission) {
		return plan, cr.ErrNotFound
	}
	var publication cr.Publication
	e := s.catalog.pool.QueryRow(ctx, `select p.generation,p.draft_id,d.snapshot_sha256,d.snapshot_canonical,p.effective_price_book_id
 from catalog.release_publication p join catalog.release_draft d using(tenant_id,draft_id)
 where p.tenant_id=$1 and p.generation=$2`, p.TenantID, generation).Scan(&publication.Generation, &publication.DraftID, &publication.SHA256, &publication.Snapshot, &publication.EffectivePriceBookID)
	if errors.Is(e, pgx.ErrNoRows) {
		return plan, cr.ErrNotFound
	}
	if e != nil {
		return plan, e
	}
	doc, e := cr.ProjectPublic(publication)
	if e != nil {
		return plan, e
	}
	raw, e := json.Marshal(doc)
	if e != nil {
		return plan, e
	}
	raw = append(raw, '\n')
	digest := sha256.Sum256(raw)
	plan = catalogFeedPlan{payload: raw, payloadSHA: hex.EncodeToString(digest[:]), sourceSHA: publication.SHA256, generation: generation}
	receiver := s.receiver.Profile()
	metadata, _, e := cr.Canonical(map[string]string{"source_sha256": plan.sourceSHA, "payload_sha256": plan.payloadSHA, "receiver_profile_sha256": s.receiver.ProfileSHA256()})
	if e != nil {
		return plan, e
	}
	plan.message = channels.Message{ChannelCode: "catalog_feed", TenantID: p.TenantID, ExternalID: receiver.ReceiverID, ThreadID: fmt.Sprint(generation), DeliveryKey: fmt.Sprintf("catalog-feed:%d:%s", generation, receiver.ReceiverID), Direction: channels.DirectionOut, Text: string(metadata)}
	if e = plan.message.Validate(); e != nil {
		return plan, e
	}
	return plan, nil
}

type catalogFeedBoundStore struct {
	service *CatalogFeed
	actor   identity.Principal
	plan    catalogFeedPlan
}

func (s *catalogFeedBoundStore) checkIntent(ctx context.Context) error {
	var actor, profile string
	e := s.service.catalog.pool.QueryRow(ctx, `select actor,receiver_profile_sha256 from catalog.release_feed_intent where tenant_id=$1 and channel_code='catalog_feed' and delivery_key=$2`, s.actor.TenantID, s.plan.message.DeliveryKey).Scan(&actor, &profile)
	if errors.Is(e, pgx.ErrNoRows) {
		return nil
	}
	if e != nil {
		return e
	}
	if actor != s.actor.Subject || profile != s.service.receiver.ProfileSHA256() {
		return cr.ErrConflict
	}
	return nil
}
func (s *catalogFeedBoundStore) Claim(ctx context.Context, m channels.Message, h string) (outbounddelivery.Claim, error) {
	expected, e := outbounddelivery.MessageSHA256(s.plan.message)
	if e != nil || expected != h {
		return outbounddelivery.Claim{}, cr.ErrInvalid
	}
	if e = s.checkIntent(ctx); e != nil {
		return outbounddelivery.Claim{}, e
	}
	claim, e := s.service.fence.claimWithAdmission(ctx, m, h, func(ctx context.Context, tx pgx.Tx) error {
		var current bool
		e := tx.QueryRow(ctx, `select exists(select 1 from catalog.release_public_current where tenant_id=$1 and generation=$2 and snapshot_sha256=$3)`, m.TenantID, s.plan.generation, s.plan.sourceSHA).Scan(&current)
		if e != nil {
			return e
		}
		if !current {
			return cr.ErrConflict
		}
		_, e = tx.Exec(ctx, `insert into catalog.release_feed_intent(tenant_id,channel_code,delivery_key,receiver_id,receiver_profile_sha256,generation,actor,request_sha256,payload_sha256,source_sha256,payload)
   values($1,'catalog_feed',$2,$3,$4,$5,$6,$7,$8,$9,$10)`, m.TenantID, m.DeliveryKey, m.ExternalID, s.service.receiver.ProfileSHA256(), s.plan.generation, s.actor.Subject, h, s.plan.payloadSHA, s.plan.sourceSHA, s.plan.payload)
		return e
	})
	if e != nil {
		return claim, e
	}
	return claim, s.checkIntent(ctx)
}
func (s *catalogFeedBoundStore) Complete(c context.Context, m channels.Message, h string, r outbounddelivery.Receipt) error {
	return s.service.fence.Complete(c, m, h, r)
}
func (s *catalogFeedBoundStore) MarkUnknown(c context.Context, m channels.Message, h, code string) error {
	return s.service.fence.MarkUnknown(c, m, h, code)
}
func (s *catalogFeedBoundStore) MarkFailed(c context.Context, m channels.Message, h, sha, code string) error {
	return s.service.fence.MarkFailed(c, m, h, sha, code)
}

type catalogFeedSender struct {
	receiver *cr.FeedHTTP
	plan     catalogFeedPlan
}

func (s catalogFeedSender) SendWithReceipt(ctx context.Context, m channels.Message) (outbounddelivery.Receipt, error) {
	got, e := outbounddelivery.MessageSHA256(m)
	expected, x := outbounddelivery.MessageSHA256(s.plan.message)
	if e != nil || x != nil || got != expected {
		return outbounddelivery.Receipt{}, cr.ErrInvalid
	}
	return s.receiver.Send(ctx, m.DeliveryKey, s.plan.generation, s.plan.payload, s.plan.payloadSHA)
}

type catalogOutboundOnly struct{}

func (catalogOutboundOnly) Code() string { return "catalog_feed" }
func (catalogOutboundOnly) Receive(context.Context) ([]channels.Message, error) {
	return nil, outbounddelivery.ErrInvalid
}
func (catalogOutboundOnly) Send(context.Context, channels.Message) error {
	return outbounddelivery.ErrInvalid
}
func (s *CatalogFeed) Send(ctx context.Context, p identity.Principal, generation int64) (cr.FeedStatus, error) {
	plan, e := s.plan(ctx, p, generation, "catalog:feed")
	if e != nil {
		return cr.FeedStatus{}, e
	}
	bound := &catalogFeedBoundStore{s, p, plan}
	channel := outbounddelivery.Channel{CodeValue: "catalog_feed", Receiver: catalogOutboundOnly{}, Sender: catalogFeedSender{s.receiver, plan}, Store: bound}
	e = channel.Send(ctx, plan.message)
	status, readErr := s.readStatus(ctx, p, plan)
	if readErr != nil {
		return status, errors.Join(e, readErr)
	}
	return status, e
}
func (s *CatalogFeed) readStatus(ctx context.Context, p identity.Principal, plan catalogFeedPlan) (cr.FeedStatus, error) {
	var out cr.FeedStatus
	e := s.catalog.pool.QueryRow(ctx, `select i.delivery_key,i.generation,d.state,i.payload_sha256,i.source_sha256
 from catalog.release_feed_intent i join communication.outbound_delivery d using(tenant_id,channel_code,delivery_key)
 where i.tenant_id=$1 and i.delivery_key=$2 and i.channel_code='catalog_feed' and i.actor=$3 and i.receiver_profile_sha256=$4`,
		p.TenantID, plan.message.DeliveryKey, p.Subject, s.receiver.ProfileSHA256()).Scan(&out.DeliveryKey, &out.Generation, &out.State, &out.PayloadSHA256, &out.SourceSHA256)
	if errors.Is(e, pgx.ErrNoRows) {
		return out, cr.ErrNotFound
	}
	return out, e
}
func (s *CatalogFeed) Status(ctx context.Context, p identity.Principal, generation int64) (cr.FeedStatus, error) {
	plan, e := s.plan(ctx, p, generation, "catalog:read")
	if e != nil {
		return cr.FeedStatus{}, e
	}
	return s.readStatus(ctx, p, plan)
}
func (s *CatalogFeed) Reconcile(ctx context.Context, p identity.Principal, generation int64) (cr.FeedStatus, error) {
	plan, e := s.plan(ctx, p, generation, "catalog:feed")
	if e != nil {
		return cr.FeedStatus{}, e
	}
	status, e := s.readStatus(ctx, p, plan)
	if e != nil {
		return status, e
	}
	if status.State == "accepted" || status.State == "failed_terminal" {
		return status, nil
	}
	if status.State != "unknown" {
		return status, outbounddelivery.ErrInProgress
	}
	receipt, e := s.receiver.Reconcile(ctx, plan.message.DeliveryKey, generation, plan.payloadSHA)
	hash, hashErr := outbounddelivery.MessageSHA256(plan.message)
	if hashErr != nil {
		return status, hashErr
	}
	if e != nil {
		var terminal *outbounddelivery.TerminalFailure
		if !errors.As(e, &terminal) {
			return status, e
		}
		if e = s.fence.ReconcileFailed(ctx, plan.message, hash, terminal.EvidenceSHA256, terminal.Code); e != nil {
			return status, e
		}
	} else if e = s.fence.ReconcileAccepted(ctx, plan.message, hash, receipt); e != nil {
		return status, e
	}
	return s.readStatus(ctx, p, plan)
}
