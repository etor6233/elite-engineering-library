package postgres

import (
	"context"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/conversationruntime"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestConversationDurableBudgetAndIntent(t *testing.T) {
	u := os.Getenv("TEST_DATABASE_URL")
	if u == "" {
		t.Skip("explicit synthetic database required")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, u)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := uuid.NewString()
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,$2,'AI fixture','AI fixture')`, tenant, "ai-"+tenant); err != nil {
		t.Fatal(err)
	}
	defer func() {
		for _, table := range []string{"conversation_intent", "conversation_reservation", "conversation_budget", "conversation_turn"} {
			pool.Exec(ctx, "delete from communication."+table+" where tenant_id=$1", tenant)
		}
		pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	}()
	s, _ := NewConversationStore(pool, time.Minute, time.Hour, 3)
	hash := strings.Repeat("b", 64)
	intent := strings.Repeat("c", 64)
	m := channels.Message{TenantID: tenant, ChannelCode: "whatsapp", ExternalID: "fixture", ThreadID: "fixture", ProviderMessageID: "one", OccurredAt: time.Now().UTC(), Direction: channels.DirectionIn, Text: "fixture"}
	c, err := s.Claim(ctx, m, hash)
	if err != nil {
		t.Fatal(err)
	}
	if err = s.Reserve(ctx, m, hash, c.Generation, 1000, 400); err != nil {
		t.Fatal(err)
	}
	if err = s.Reserve(ctx, m, hash, c.Generation, 1000, 400); err != nil {
		t.Fatal("reservation recovery", err)
	}
	if err = s.Reserve(ctx, m, hash, c.Generation, 2000, 400); !errors.Is(err, conversationruntime.ErrBudgetConfiguration) {
		t.Fatal("cap reset", err)
	}
	if err = s.PinTool(ctx, m, hash, c.Generation, intent); err != nil {
		t.Fatal(err)
	}
	if err = s.PinTool(ctx, m, hash, c.Generation, strings.Repeat("d", 64)); !errors.Is(err, conversationruntime.ErrToolIntentConflict) {
		t.Fatal("changed tool", err)
	}
	if err = s.FailRetryable(ctx, m, hash, c.Generation, "SYNTHETIC_RESPONSE_LOSS"); err != nil {
		t.Fatal(err)
	}
	restarted, _ := NewConversationStore(pool, time.Minute, time.Hour, 3)
	again, err := restarted.Claim(ctx, m, hash)
	if err != nil || again.Generation != 2 {
		t.Fatal(again, err)
	}
	if err = restarted.Reserve(ctx, m, hash, again.Generation, 1000, 400); err != nil {
		t.Fatal(err)
	}
	if err = restarted.PinTool(ctx, m, hash, again.Generation, intent); err != nil {
		t.Fatal(err)
	}
	if err = s.PinTool(ctx, m, hash, c.Generation, intent); !errors.Is(err, conversationruntime.ErrReplayConflict) {
		t.Fatal("stale tool allowed", err)
	}
	var spent, count int64
	if err = pool.QueryRow(ctx, `select reserved,(select count(*) from communication.conversation_reservation where tenant_id=$1) from communication.conversation_budget where tenant_id=$1`, tenant).Scan(&spent, &count); err != nil || spent != 800 || count != 2 {
		t.Fatal(spent, count, err)
	}
	// Two independently claimed messages compete for the final 200 tokens.
	var wg sync.WaitGroup
	answers := make(chan error, 2)
	for _, id := range []string{"two", "three"} {
		n := m
		n.ProviderMessageID = id
		cl, err := s.Claim(ctx, n, hash)
		if err != nil {
			t.Fatal(err)
		}
		wg.Add(1)
		go func() { defer wg.Done(); answers <- s.Reserve(ctx, n, hash, cl.Generation, 1000, 200) }()
	}
	wg.Wait()
	close(answers)
	wins, blocked := 0, 0
	for err := range answers {
		if err == nil {
			wins++
		} else if errors.Is(err, conversationruntime.ErrBudgetExceeded) {
			blocked++
		} else {
			t.Fatal(err)
		}
	}
	if wins != 1 || blocked != 1 {
		t.Fatal("budget race", wins, blocked)
	}
	if _, err = pool.Exec(ctx, `update communication.conversation_intent set intent_sha256=$2 where tenant_id=$1`, tenant, hash); err == nil {
		t.Fatal("intent mutation accepted")
	}
	if _, err = pool.Exec(ctx, `update communication.conversation_reservation set tokens=1 where tenant_id=$1`, tenant); err == nil {
		t.Fatal("reservation mutation accepted")
	}
	t.Log("CONVERSATION_DURABLE_BUDGET_INTENT_PASS")
}
