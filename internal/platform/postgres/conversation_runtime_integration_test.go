package postgres

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/conversationruntime"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestConversationStoreDurabilityIsolationAndFencing(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c20001"
	_, _ = pool.Exec(ctx, `delete from communication.conversation_turn where tenant_id=$1`, tenant)
	_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	defer func() {
		_, _ = pool.Exec(ctx, `delete from communication.conversation_turn where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	}()
	if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'conversation-test','Conversation Test','Conversation')`, tenant); err != nil {
		t.Fatal(err)
	}
	store, err := NewConversationStore(pool, time.Second, 24*time.Hour, 2)
	if err != nil {
		t.Fatal(err)
	}
	message := channels.Message{ChannelCode: "whatsapp", TenantID: tenant, ExternalID: "contact-a", ThreadID: "thread-a", ProviderMessageID: "wamid.pg.1", OccurredAt: time.Date(2026, 9, 4, 13, 0, 0, 0, time.UTC), Direction: channels.DirectionIn, Text: "Necesito una cotización"}
	hash := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	claim, err := store.Claim(ctx, message, hash)
	if err != nil || claim.Replay {
		t.Fatalf("initial claim=%#v err=%v", claim, err)
	}
	if _, err := store.Claim(ctx, message, hash); !errors.Is(err, conversationruntime.ErrInProgress) {
		t.Fatalf("active duplicate err=%v", err)
	}
	if _, err := store.Claim(ctx, message, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"); !errors.Is(err, conversationruntime.ErrReplayConflict) {
		t.Fatalf("divergent duplicate err=%v", err)
	}
	completion := conversationruntime.Completion{State: conversationruntime.StateCompleted, ResponseText: "Cotización creada", LLMResponseID: "resp-pg-1", ToolName: "create_quote", ToolArguments: []byte(`{"product":"bike","quantity":1}`), ToolResult: `{"quote_id":"q-1"}`, InputTokens: 10, OutputTokens: 4, ReservedTokens: 900}
	if err := store.Complete(ctx, message, hash, claim.Generation, completion); err != nil {
		t.Fatal(err)
	}
	replay, err := store.Claim(ctx, message, hash)
	if err != nil || !replay.Replay || replay.ResponseText != completion.ResponseText {
		t.Fatalf("replay=%#v err=%v", replay, err)
	}
	if err := store.Complete(ctx, message, hash, claim.Generation, completion); !errors.Is(err, conversationruntime.ErrReplayConflict) {
		t.Fatalf("terminal mutation err=%v", err)
	}

	second := message
	second.ProviderMessageID = "wamid.pg.2"
	second.Text = "Gracias"
	second.OccurredAt = second.OccurredAt.Add(time.Minute)
	if _, err := store.Claim(ctx, second, hash); err != nil {
		t.Fatal(err)
	}
	if err := store.Complete(ctx, second, hash, 1, conversationruntime.Completion{State: conversationruntime.StateCompleted, ResponseText: "De nada"}); err != nil {
		t.Fatal(err)
	}
	history, err := store.History(ctx, second, 10)
	if err != nil || len(history) != 1 || history[0].UserText != message.Text || history[0].AssistantText != completion.ResponseText {
		t.Fatalf("history=%#v err=%v", history, err)
	}
	isolated := second
	isolated.ProviderMessageID = "wamid.pg.isolated"
	isolated.ExternalID = "contact-b"
	history, err = store.History(ctx, isolated, 10)
	if err != nil || len(history) != 0 {
		t.Fatalf("cross-contact history=%#v err=%v", history, err)
	}

	retry := message
	retry.ProviderMessageID = "wamid.pg.retry"
	if _, err := store.Claim(ctx, retry, hash); err != nil {
		t.Fatal(err)
	}
	if err := store.FailRetryable(ctx, retry, hash, 1, "MODEL_START_FAILED"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Claim(ctx, retry, hash); err != nil {
		t.Fatal(err)
	}
	if err := store.FailRetryable(ctx, retry, hash, 2, "MODEL_START_FAILED"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Claim(ctx, retry, hash); !errors.Is(err, conversationruntime.ErrTerminal) {
		t.Fatalf("poison limit err=%v", err)
	}

	concurrent := message
	concurrent.ProviderMessageID = "wamid.pg.concurrent"
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := store.Claim(ctx, concurrent, hash)
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	wins, fenced := 0, 0
	for err := range errs {
		switch {
		case err == nil:
			wins++
		case errors.Is(err, conversationruntime.ErrInProgress):
			fenced++
		default:
			t.Fatalf("concurrent claim err=%v", err)
		}
	}
	if wins != 1 || fenced != 1 {
		t.Fatalf("concurrency winners=%d fenced=%d", wins, fenced)
	}

	if _, err := pool.Exec(ctx, `update communication.conversation_turn set external_id='attacker' where tenant_id=$1 and channel_code=$2 and provider_message_id=$3`, tenant, message.ChannelCode, message.ProviderMessageID); err == nil {
		t.Fatal("expected immutable terminal identity rejection")
	}
}
