package postgres

import (
	"context"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/conversationruntime"
	"encoding/hex"
	"errors"
	"github.com/jackc/pgx/v5"
)

func (s *ConversationStore) activeTurn(ctx context.Context, tx pgx.Tx, m channels.Message, hash string, generation int) error {
	var active bool
	err := tx.QueryRow(ctx, `select state='processing' and request_sha256_hex=$4 and attempt_count=$5 and locked_until>clock_timestamp() and expires_at>clock_timestamp() from communication.conversation_turn where tenant_id=$1 and channel_code=$2 and provider_message_id=$3 for update`, m.TenantID, m.ChannelCode, m.ProviderMessageID, hash, generation).Scan(&active)
	if errors.Is(err, pgx.ErrNoRows) || err == nil && !active {
		return conversationruntime.ErrReplayConflict
	}
	return err
}

// Reserve charges each attempt before any external model request. Retries after
// an uncertain provider response spend again. Restart never replenishes the cap.
// A cap change is rejected; changing the budget requires an explicit operator
// migration, never an automatic reset or an unbounded scheduled refill.
func (s *ConversationStore) Reserve(ctx context.Context, m channels.Message, hash string, generation int, cap, tokens int64) error {
	if cap < 1 || cap > 1000000000000 || tokens < 1 {
		return conversationruntime.ErrInvalidRuntime
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = s.activeTurn(ctx, tx, m, hash, generation); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `insert into communication.conversation_budget(tenant_id,cap)values($1,$2)on conflict do nothing`, m.TenantID, cap); err != nil {
		return err
	}
	var oldCap, spent int64
	if err = tx.QueryRow(ctx, `select cap,reserved from communication.conversation_budget where tenant_id=$1 for update`, m.TenantID).Scan(&oldCap, &spent); err != nil {
		return err
	}
	if oldCap != cap {
		return conversationruntime.ErrBudgetConfiguration
	}
	var previous int64
	err = tx.QueryRow(ctx, `select tokens from communication.conversation_reservation where tenant_id=$1 and channel_code=$2 and provider_message_id=$3 and generation=$4`, m.TenantID, m.ChannelCode, m.ProviderMessageID, generation).Scan(&previous)
	if err == nil {
		if previous != tokens {
			return conversationruntime.ErrReplayConflict
		}
		return tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if tokens > cap-spent {
		return conversationruntime.ErrBudgetExceeded
	}
	if _, err = tx.Exec(ctx, `update communication.conversation_budget set reserved=reserved+$2 where tenant_id=$1`, m.TenantID, tokens); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `insert into communication.conversation_reservation(tenant_id,channel_code,provider_message_id,generation,tokens)values($1,$2,$3,$4,$5)`, m.TenantID, m.ChannelCode, m.ProviderMessageID, generation, tokens); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// PinTool binds the first admitted tool, arguments, contact and policy. A later
// model attempt cannot reinterpret a turn into a different business effect.
// The domain owner still supplies effect idempotency across an expired lease.
func (s *ConversationStore) PinTool(ctx context.Context, m channels.Message, hash string, generation int, intent string) error {
	b, err := hex.DecodeString(intent)
	if err != nil || len(b) != 32 {
		return conversationruntime.ErrInvalidRuntime
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = s.activeTurn(ctx, tx, m, hash, generation); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `insert into communication.conversation_intent(tenant_id,channel_code,provider_message_id,intent_sha256)values($1,$2,$3,$4)on conflict do nothing`, m.TenantID, m.ChannelCode, m.ProviderMessageID, intent); err != nil {
		return err
	}
	var old string
	if err = tx.QueryRow(ctx, `select intent_sha256 from communication.conversation_intent where tenant_id=$1 and channel_code=$2 and provider_message_id=$3`, m.TenantID, m.ChannelCode, m.ProviderMessageID).Scan(&old); err != nil {
		return err
	}
	if old != intent {
		return conversationruntime.ErrToolIntentConflict
	}
	return tx.Commit(ctx)
}
