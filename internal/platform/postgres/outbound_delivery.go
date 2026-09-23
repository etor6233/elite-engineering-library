package postgres

import (
	"context"
	"errors"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/outbounddelivery"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboundDeliveryStore struct {
	pool    *pgxpool.Pool
	hmacKey []byte
	lease   time.Duration
}

func NewOutboundDeliveryStore(pool *pgxpool.Pool, hmacKey []byte, lease time.Duration) (*OutboundDeliveryStore, error) {
	if pool == nil || len(hmacKey) < 32 || lease < time.Second || lease > 10*time.Minute {
		return nil, outbounddelivery.ErrInvalid
	}
	return &OutboundDeliveryStore{pool: pool, hmacKey: append([]byte(nil), hmacKey...), lease: lease}, nil
}

func (s *OutboundDeliveryStore) Claim(ctx context.Context, message channels.Message, requestHash string) (outbounddelivery.Claim, error) {
	return s.claimWithAdmission(ctx, message, requestHash, nil)
}

func (s *OutboundDeliveryStore) claimWithAdmission(ctx context.Context, message channels.Message, requestHash string, admit func(context.Context, pgx.Tx) error) (outbounddelivery.Claim, error) {
	if s == nil || s.pool == nil || len(requestHash) != 64 {
		return outbounddelivery.Claim{}, outbounddelivery.ErrInvalid
	}
	recipientHash, err := contactidentity.ExternalDigest(s.hmacKey, message.TenantID, message.ChannelCode, message.ExternalID)
	if err != nil {
		return outbounddelivery.Claim{}, err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return outbounddelivery.Claim{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err := tx.Exec(ctx, `insert into communication.outbound_delivery
(tenant_id,channel_code,delivery_key,request_sha256_hex,recipient_hmac,state,attempt_count,locked_until)
values($1,$2,$3,$4,$5,'sending',1,clock_timestamp()+$6*interval '1 millisecond') on conflict do nothing`, message.TenantID, message.ChannelCode, message.DeliveryKey, requestHash, recipientHash, s.lease.Milliseconds())
	if err != nil {
		return outbounddelivery.Claim{}, outboundWriteError(err)
	}
	if result.RowsAffected() == 1 {
		if admit != nil {
			if err = admit(ctx, tx); err != nil {
				return outbounddelivery.Claim{}, err
			}
		}
		_, err = tx.Exec(ctx, `insert into communication.outbound_delivery_event(tenant_id,channel_code,delivery_key,sequence,state) values($1,$2,$3,1,'sending')`, message.TenantID, message.ChannelCode, message.DeliveryKey)
		if err != nil {
			return outbounddelivery.Claim{}, outboundWriteError(err)
		}
		if err = tx.Commit(ctx); err != nil {
			return outbounddelivery.Claim{}, outboundWriteError(err)
		}
		return outbounddelivery.Claim{}, nil
	}
	var storedHash, state string
	var lockedUntil *time.Time
	err = tx.QueryRow(ctx, `select request_sha256_hex,state,locked_until from communication.outbound_delivery where tenant_id=$1 and channel_code=$2 and delivery_key=$3 for update`, message.TenantID, message.ChannelCode, message.DeliveryKey).Scan(&storedHash, &state, &lockedUntil)
	if err != nil {
		return outbounddelivery.Claim{}, err
	}
	if storedHash != requestHash {
		return outbounddelivery.Claim{}, outbounddelivery.ErrConflict
	}
	switch state {
	case "accepted":
		if err = tx.Commit(ctx); err != nil {
			return outbounddelivery.Claim{}, err
		}
		return outbounddelivery.Claim{Replay: true}, nil
	case "unknown":
		return outbounddelivery.Claim{}, outbounddelivery.ErrUnknown
	case "failed_terminal":
		return outbounddelivery.Claim{}, outbounddelivery.ErrTerminal
	case "sending":
		if lockedUntil != nil && lockedUntil.After(time.Now()) {
			return outbounddelivery.Claim{}, outbounddelivery.ErrInProgress
		}
		_, err = tx.Exec(ctx, `update communication.outbound_delivery set state='unknown',locked_until=null,failure_code='LEASE_EXPIRED',updated_at=clock_timestamp() where tenant_id=$1 and channel_code=$2 and delivery_key=$3 and state='sending'`, message.TenantID, message.ChannelCode, message.DeliveryKey)
		if err != nil {
			return outbounddelivery.Claim{}, err
		}
		_, err = tx.Exec(ctx, `insert into communication.outbound_delivery_event(tenant_id,channel_code,delivery_key,sequence,state,failure_code) select tenant_id,channel_code,delivery_key,2,'unknown','LEASE_EXPIRED' from communication.outbound_delivery where tenant_id=$1 and channel_code=$2 and delivery_key=$3`, message.TenantID, message.ChannelCode, message.DeliveryKey)
		if err != nil {
			return outbounddelivery.Claim{}, outboundWriteError(err)
		}
		if err = tx.Commit(ctx); err != nil {
			return outbounddelivery.Claim{}, outboundWriteError(err)
		}
		return outbounddelivery.Claim{}, outbounddelivery.ErrUnknown
	default:
		return outbounddelivery.Claim{}, outbounddelivery.ErrInvalid
	}
}

func (s *OutboundDeliveryStore) Complete(ctx context.Context, message channels.Message, requestHash string, receipt outbounddelivery.Receipt) error {
	if err := receipt.Validate(); err != nil {
		return err
	}
	providerHash, err := contactidentity.ExternalDigest(s.hmacKey, message.TenantID, message.ChannelCode, receipt.ProviderMessageID)
	if err != nil {
		return err
	}
	return s.transition(ctx, message, requestHash, "sending", "accepted", providerHash, receipt.EvidenceSHA256, "", receipt.AcceptedAt.UTC())
}

func (s *OutboundDeliveryStore) MarkUnknown(ctx context.Context, message channels.Message, requestHash, code string) error {
	if code == "" {
		return outbounddelivery.ErrInvalid
	}
	return s.transition(ctx, message, requestHash, "sending", "unknown", "", "", code, time.Time{})
}

func (s *OutboundDeliveryStore) MarkFailed(ctx context.Context, message channels.Message, requestHash, evidenceHash, code string) error {
	if len(evidenceHash) != 64 || code == "" {
		return outbounddelivery.ErrInvalid
	}
	return s.transition(ctx, message, requestHash, "sending", "failed_terminal", "", evidenceHash, code, time.Time{})
}

func (s *OutboundDeliveryStore) ReconcileAccepted(ctx context.Context, message channels.Message, requestHash string, receipt outbounddelivery.Receipt) error {
	if err := receipt.Validate(); err != nil {
		return err
	}
	providerHash, err := contactidentity.ExternalDigest(s.hmacKey, message.TenantID, message.ChannelCode, receipt.ProviderMessageID)
	if err != nil {
		return err
	}
	return s.transition(ctx, message, requestHash, "unknown", "accepted", providerHash, receipt.EvidenceSHA256, "", receipt.AcceptedAt.UTC())
}

func (s *OutboundDeliveryStore) ReconcileFailed(ctx context.Context, message channels.Message, requestHash, evidenceHash, code string) error {
	if len(evidenceHash) != 64 || code == "" {
		return outbounddelivery.ErrInvalid
	}
	return s.transition(ctx, message, requestHash, "unknown", "failed_terminal", "", evidenceHash, code, time.Time{})
}

func (s *OutboundDeliveryStore) transition(ctx context.Context, message channels.Message, requestHash, from, to, providerHash, evidenceHash, code string, acceptedAt time.Time) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var storedHash, state string
	var sequence int64
	err = tx.QueryRow(ctx, `select request_sha256_hex,state from communication.outbound_delivery where tenant_id=$1 and channel_code=$2 and delivery_key=$3 for update`, message.TenantID, message.ChannelCode, message.DeliveryKey).Scan(&storedHash, &state)
	if errors.Is(err, pgx.ErrNoRows) {
		return outbounddelivery.ErrConflict
	}
	if err != nil {
		return err
	}
	if storedHash != requestHash || state != from {
		return outbounddelivery.ErrConflict
	}
	if err = tx.QueryRow(ctx, `select coalesce(max(sequence),0) from communication.outbound_delivery_event where tenant_id=$1 and channel_code=$2 and delivery_key=$3`, message.TenantID, message.ChannelCode, message.DeliveryKey).Scan(&sequence); err != nil {
		return err
	}
	result, err := tx.Exec(ctx, `update communication.outbound_delivery set state=$4,locked_until=null,provider_message_hmac=nullif($5,''),evidence_sha256_hex=nullif($6,''),failure_code=nullif($7,''),accepted_at=$8,updated_at=clock_timestamp() where tenant_id=$1 and channel_code=$2 and delivery_key=$3 and state=$9 and request_sha256_hex=$10`, message.TenantID, message.ChannelCode, message.DeliveryKey, to, providerHash, evidenceHash, code, nullableTime(acceptedAt), from, requestHash)
	if err != nil {
		return outboundWriteError(err)
	}
	if result.RowsAffected() != 1 {
		return outbounddelivery.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into communication.outbound_delivery_event(tenant_id,channel_code,delivery_key,sequence,state,evidence_sha256_hex,failure_code) values($1,$2,$3,$4,$5,nullif($6,''),nullif($7,''))`, message.TenantID, message.ChannelCode, message.DeliveryKey, sequence+1, to, evidenceHash, code)
	if err != nil {
		return outboundWriteError(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return outboundWriteError(err)
	}
	return nil
}

func nullableTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value
}

func outboundWriteError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23505" || pgErr.Code == "40001") {
		return outbounddelivery.ErrConflict
	}
	return err
}
