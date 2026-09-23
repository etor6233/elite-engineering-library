package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/conversationruntime"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConversationStore struct {
	pool        *pgxpool.Pool
	lease       time.Duration
	retention   time.Duration
	maxAttempts int
}

func NewConversationStore(pool *pgxpool.Pool, lease, retention time.Duration, maxAttempts int) (*ConversationStore, error) {
	if pool == nil || lease < time.Second || lease > 10*time.Minute || retention < time.Hour || maxAttempts < 1 || maxAttempts > 100 {
		return nil, errors.New("postgres conversation store: invalid configuration")
	}
	return &ConversationStore{pool: pool, lease: lease, retention: retention, maxAttempts: maxAttempts}, nil
}

func (s *ConversationStore) Claim(ctx context.Context, message channels.Message, requestHash string) (conversationruntime.Claim, error) {
	if err := message.Validate(); err != nil || message.Direction != channels.DirectionIn || len(requestHash) != 64 {
		return conversationruntime.Claim{}, errors.New("postgres conversation store: invalid claim")
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return conversationruntime.Claim{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err := tx.Exec(ctx, `insert into communication.conversation_turn
(tenant_id,channel_code,provider_message_id,external_id,thread_id,occurred_at,request_sha256_hex,state,attempt_count,locked_until,expires_at)
values($1,$2,$3,$4,nullif($5,''),$6,$7,'processing',1,clock_timestamp()+$8*interval '1 millisecond',clock_timestamp()+$9*interval '1 millisecond')
on conflict do nothing`, message.TenantID, message.ChannelCode, message.ProviderMessageID, message.ExternalID, message.ThreadID, message.OccurredAt, requestHash, s.lease.Milliseconds(), s.retention.Milliseconds())
	if err != nil {
		return conversationruntime.Claim{}, err
	}
	if result.RowsAffected() == 1 {
		if err := tx.Commit(ctx); err != nil {
			return conversationruntime.Claim{}, err
		}
		return conversationruntime.Claim{Generation: 1}, nil
	}
	var storedHash, state string
	var response *string
	var leaseActive, expired bool
	var attempts int
	err = tx.QueryRow(ctx, `select request_sha256_hex,state,assistant_text,coalesce(locked_until>clock_timestamp(),false),attempt_count,expires_at<=clock_timestamp()
from communication.conversation_turn where tenant_id=$1 and channel_code=$2 and provider_message_id=$3 for update`, message.TenantID, message.ChannelCode, message.ProviderMessageID).Scan(&storedHash, &state, &response, &leaseActive, &attempts, &expired)
	if err != nil {
		return conversationruntime.Claim{}, err
	}
	if storedHash != requestHash {
		return conversationruntime.Claim{}, conversationruntime.ErrReplayConflict
	}
	if expired {
		return conversationruntime.Claim{}, conversationruntime.ErrTerminal
	}
	if state == string(conversationruntime.StateCompleted) || state == string(conversationruntime.StateHandedOff) {
		if response == nil {
			return conversationruntime.Claim{}, conversationruntime.ErrTerminal
		}
		if err := tx.Commit(ctx); err != nil {
			return conversationruntime.Claim{}, err
		}
		return conversationruntime.Claim{Replay: true, State: conversationruntime.TurnState(state), ResponseText: *response}, nil
	}
	if state == "failed_terminal" {
		return conversationruntime.Claim{}, conversationruntime.ErrTerminal
	}
	if state == "processing" && leaseActive {
		return conversationruntime.Claim{}, conversationruntime.ErrInProgress
	}
	if attempts >= s.maxAttempts {
		_, err = tx.Exec(ctx, `update communication.conversation_turn set state='failed_terminal',locked_until=null,failure_code='MAX_ATTEMPTS',completed_at=clock_timestamp(),updated_at=clock_timestamp()
where tenant_id=$1 and channel_code=$2 and provider_message_id=$3`, message.TenantID, message.ChannelCode, message.ProviderMessageID)
		if err != nil {
			return conversationruntime.Claim{}, err
		}
		if err := tx.Commit(ctx); err != nil {
			return conversationruntime.Claim{}, err
		}
		return conversationruntime.Claim{}, conversationruntime.ErrTerminal
	}
	_, err = tx.Exec(ctx, `update communication.conversation_turn set state='processing',attempt_count=attempt_count+1,locked_until=clock_timestamp()+$4*interval '1 millisecond',failure_code=null,updated_at=clock_timestamp()
where tenant_id=$1 and channel_code=$2 and provider_message_id=$3`, message.TenantID, message.ChannelCode, message.ProviderMessageID, s.lease.Milliseconds())
	if err != nil {
		return conversationruntime.Claim{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return conversationruntime.Claim{}, err
	}
	return conversationruntime.Claim{Generation: attempts + 1}, nil
}

func (s *ConversationStore) History(ctx context.Context, message channels.Message, limit int) ([]conversationruntime.HistoryMessage, error) {
	if limit == 0 {
		return nil, nil
	}
	if limit < 0 || limit > 20 {
		return nil, errors.New("postgres conversation store: invalid history limit")
	}
	rows, err := s.pool.Query(ctx, `select user_text,assistant_text from communication.conversation_turn
where tenant_id=$1 and channel_code=$2 and external_id=$3 and thread_id is not distinct from nullif($4,'')
and expires_at>clock_timestamp() and provider_message_id<>$5 and state='completed' and user_text is not null
order by completed_at desc limit $6`, message.TenantID, message.ChannelCode, message.ExternalID, message.ThreadID, message.ProviderMessageID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reversed := make([]conversationruntime.HistoryMessage, 0, limit)
	for rows.Next() {
		var item conversationruntime.HistoryMessage
		if err := rows.Scan(&item.UserText, &item.AssistantText); err != nil {
			return nil, err
		}
		reversed = append(reversed, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	result := make([]conversationruntime.HistoryMessage, len(reversed))
	for i := range reversed {
		result[len(reversed)-1-i] = reversed[i]
	}
	return result, nil
}

func (s *ConversationStore) Complete(ctx context.Context, message channels.Message, requestHash string, generation int, completion conversationruntime.Completion) error {
	if generation < 1 {
		return conversationruntime.ErrReplayConflict
	}
	if completion.State != conversationruntime.StateCompleted && completion.State != conversationruntime.StateHandedOff || completion.ResponseText == "" {
		return errors.New("postgres conversation store: invalid completion")
	}
	responseSum := sha256.Sum256([]byte(completion.ResponseText))
	responseHash := hex.EncodeToString(responseSum[:])
	var arguments any
	if len(completion.ToolArguments) > 0 {
		if !json.Valid(completion.ToolArguments) {
			return errors.New("postgres conversation store: invalid tool arguments")
		}
		arguments = completion.ToolArguments
	}
	result, err := s.pool.Exec(ctx, `update communication.conversation_turn set
state=$4,user_text=$5,assistant_text=$6,response_sha256_hex=$7,llm_response_id=nullif($8,''),tool_name=nullif($9,''),tool_arguments=$10,tool_result=nullif($11,''),input_tokens=$12,output_tokens=$13,reserved_tokens=$14,failure_code=nullif($15,''),locked_until=null,completed_at=clock_timestamp(),updated_at=clock_timestamp()
where tenant_id=$1 and channel_code=$2 and provider_message_id=$3 and request_sha256_hex=$16 and state='processing' and attempt_count=$17 and locked_until>clock_timestamp() and expires_at>clock_timestamp()`,
		message.TenantID, message.ChannelCode, message.ProviderMessageID, string(completion.State), message.Text, completion.ResponseText, responseHash,
		completion.LLMResponseID, completion.ToolName, arguments, completion.ToolResult, completion.InputTokens, completion.OutputTokens, completion.ReservedTokens, completion.FailureCode, requestHash, generation)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return conversationruntime.ErrReplayConflict
	}
	return nil
}

func (s *ConversationStore) FailRetryable(ctx context.Context, message channels.Message, requestHash string, generation int, code string) error {
	if code == "" {
		return errors.New("postgres conversation store: failure code required")
	}
	result, err := s.pool.Exec(ctx, `update communication.conversation_turn set state='retryable',locked_until=null,failure_code=$4,updated_at=clock_timestamp()
where tenant_id=$1 and channel_code=$2 and provider_message_id=$3 and request_sha256_hex=$5 and state='processing' and attempt_count=$6 and locked_until>clock_timestamp() and expires_at>clock_timestamp()`, message.TenantID, message.ChannelCode, message.ProviderMessageID, code, requestHash, generation)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("%w: retry transition lost", conversationruntime.ErrReplayConflict)
	}
	return nil
}
