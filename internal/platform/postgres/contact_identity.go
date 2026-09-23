package postgres

import (
	"context"
	"errors"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/conversationruntime"
	"elite.local/enterprise/internal/domainbind"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactIdentityStore struct {
	pool    *pgxpool.Pool
	hmacKey []byte
}

func NewContactIdentityStore(pool *pgxpool.Pool, hmacKey []byte) (*ContactIdentityStore, error) {
	if pool == nil || len(hmacKey) < 32 {
		return nil, contactidentity.ErrInvalid
	}
	key := append([]byte(nil), hmacKey...)
	return &ContactIdentityStore{pool: pool, hmacKey: key}, nil
}

func (s *ContactIdentityStore) Resolve(ctx context.Context, message channels.Message) (conversationruntime.Contact, error) {
	if s == nil || s.pool == nil || message.Direction != channels.DirectionIn {
		return conversationruntime.Contact{}, contactidentity.ErrInvalid
	}
	digest, err := contactidentity.ExternalDigest(s.hmacKey, message.TenantID, message.ChannelCode, message.ExternalID)
	if err != nil {
		return conversationruntime.Contact{}, err
	}
	var leadID, organizationID, subjectID string
	var piiAllowed bool
	err = s.pool.QueryRow(ctx, `select b.lead_id,l.organization_id,b.subject_id,b.pii_allowed
from communication.contact_channel_binding b
join crm.lead l on l.tenant_id=b.tenant_id and l.lead_id=b.lead_id
where b.tenant_id=$1 and b.channel_code=$2 and b.external_id_hmac=$3 and b.state='active' and b.effective_at<=statement_timestamp()`, message.TenantID, message.ChannelCode, digest).Scan(&leadID, &organizationID, &subjectID, &piiAllowed)
	if errors.Is(err, pgx.ErrNoRows) {
		return conversationruntime.Contact{}, contactidentity.ErrNotFound
	}
	if err != nil {
		return conversationruntime.Contact{}, err
	}
	return conversationruntime.Contact{Scope: domainbind.Scope{OrganizationID: organizationID, LeadID: leadID}, SubjectID: subjectID, PIIAllowed: piiAllowed}, nil
}

func (s *ContactIdentityStore) Apply(ctx context.Context, command contactidentity.Command) (contactidentity.Receipt, error) {
	if s == nil || s.pool == nil {
		return contactidentity.Receipt{}, contactidentity.ErrInvalid
	}
	digest, requestHash, err := command.RequestSHA256(s.hmacKey)
	if err != nil {
		return contactidentity.Receipt{}, identityWriteError(err)
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return contactidentity.Receipt{}, identityWriteError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var replayHash string
	var replayVersion int64
	var replayState contactidentity.State
	err = tx.QueryRow(ctx, `select request_sha256_hex,binding_version,state from communication.contact_channel_binding_decision where tenant_id=$1 and request_id=$2`, command.TenantID, command.RequestID).Scan(&replayHash, &replayVersion, &replayState)
	if err == nil {
		if replayHash != requestHash {
			return contactidentity.Receipt{}, contactidentity.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return contactidentity.Receipt{}, err
		}
		return contactidentity.Receipt{Version: replayVersion, State: replayState, Replayed: true}, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return contactidentity.Receipt{}, err
	}
	var currentVersion int64
	var currentLead, currentSubject string
	err = tx.QueryRow(ctx, `select version,lead_id,subject_id from communication.contact_channel_binding where tenant_id=$1 and channel_code=$2 and external_id_hmac=$3 for update`, command.TenantID, command.ChannelCode, digest).Scan(&currentVersion, &currentLead, &currentSubject)
	exists := err == nil
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return contactidentity.Receipt{}, err
	}
	if command.ExpectedVersion == 0 {
		if exists || command.State != contactidentity.StateActive {
			return contactidentity.Receipt{}, contactidentity.ErrConflict
		}
		currentVersion = 1
		_, err = tx.Exec(ctx, `insert into communication.contact_channel_binding
(tenant_id,channel_code,external_id_hmac,lead_id,subject_id,pii_allowed,state,policy_version,evidence_sha256_hex,effective_at,version)
values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`, command.TenantID, command.ChannelCode, digest, command.LeadID, command.SubjectID, command.PIIAllowed, command.State, command.PolicyVersion, command.EvidenceSHA256, command.EffectiveAt.UTC(), currentVersion)
	} else {
		if !exists || currentVersion != command.ExpectedVersion || currentLead != command.LeadID || currentSubject != command.SubjectID {
			return contactidentity.Receipt{}, contactidentity.ErrConflict
		}
		currentVersion++
		_, err = tx.Exec(ctx, `update communication.contact_channel_binding set pii_allowed=$4,state=$5,policy_version=$6,evidence_sha256_hex=$7,effective_at=$8,version=$9,updated_at=clock_timestamp()
where tenant_id=$1 and channel_code=$2 and external_id_hmac=$3 and version=$10`, command.TenantID, command.ChannelCode, digest, command.PIIAllowed, command.State, command.PolicyVersion, command.EvidenceSHA256, command.EffectiveAt.UTC(), currentVersion, command.ExpectedVersion)
	}
	if err != nil {
		return contactidentity.Receipt{}, identityWriteError(err)
	}
	_, err = tx.Exec(ctx, `insert into communication.contact_channel_binding_decision
(tenant_id,request_id,request_sha256_hex,channel_code,external_id_hmac,binding_version,lead_id,subject_id,pii_allowed,state,policy_version,evidence_sha256_hex,effective_at)
values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`, command.TenantID, command.RequestID, requestHash, command.ChannelCode, digest, currentVersion, command.LeadID, command.SubjectID, command.PIIAllowed, command.State, command.PolicyVersion, command.EvidenceSHA256, command.EffectiveAt.UTC())
	if err != nil {
		return contactidentity.Receipt{}, identityWriteError(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return contactidentity.Receipt{}, identityWriteError(err)
	}
	return contactidentity.Receipt{Version: currentVersion, State: command.State}, nil
}

func identityWriteError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23505" || pgErr.Code == "40001") {
		return contactidentity.ErrConflict
	}
	return err
}
