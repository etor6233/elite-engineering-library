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
	"testing"
	"time"
)

func TestConversationLeaseAndRetentionRegression(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("explicit fixture database required")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := uuid.NewString()
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,$2,'Synthetic AI','Synthetic AI')`, tenant, "ai-"+tenant); err != nil {
		t.Fatal(err)
	}
	defer func() {
		pool.Exec(ctx, `delete from communication.conversation_turn where tenant_id=$1`, tenant)
		pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	}()
	s, _ := NewConversationStore(pool, time.Minute, time.Hour, 3)
	hash := strings.Repeat("a", 64)
	m := channels.Message{TenantID: tenant, ChannelCode: "whatsapp", ExternalID: "synthetic-contact", ThreadID: "synthetic-thread", ProviderMessageID: "lease", OccurredAt: time.Now().UTC(), Direction: channels.DirectionIn, Text: "consulta"}
	c := conversationruntime.Completion{State: conversationruntime.StateCompleted, ResponseText: "fixture"}
	for _, mode := range []string{"expired-complete", "expired-retry", "reclaimed-complete", "reclaimed-retry"} {
		t.Run(mode, func(t *testing.T) {
			m.ProviderMessageID = mode
			if _, err := s.Claim(ctx, m, hash); err != nil {
				t.Fatal(err)
			}
			if _, err := pool.Exec(ctx, `update communication.conversation_turn set locked_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and provider_message_id=$2`, tenant, mode); err != nil {
				t.Fatal(err)
			}
			if strings.HasPrefix(mode, "reclaimed") {
				if _, err := s.Claim(ctx, m, hash); err != nil {
					t.Fatal(err)
				}
			}
			var err error
			if strings.HasSuffix(mode, "complete") {
				err = s.Complete(ctx, m, hash, 1, c)
			} else {
				err = s.FailRetryable(ctx, m, hash, 1, "SYNTHETIC_FAILURE")
			}
			if !errors.Is(err, conversationruntime.ErrReplayConflict) {
				t.Fatalf("expired/old generation accepted: %v", err)
			}
		})
	}
	// A terminal row older than retention is a fixture of a database surviving
	// its application process; replay and model history must both reject it.
	m.ProviderMessageID = "expired-history"
	if _, err := pool.Exec(ctx, `insert into communication.conversation_turn(tenant_id,channel_code,provider_message_id,external_id,thread_id,occurred_at,request_sha256_hex,state,attempt_count,user_text,assistant_text,response_sha256_hex,completed_at,expires_at)values($1,'whatsapp',$2,$3,$4,$5,$6,'completed',1,'private old text','private old answer',$6,clock_timestamp()-interval '2 hours',clock_timestamp()-interval '1 hour')`, tenant, m.ProviderMessageID, m.ExternalID, m.ThreadID, m.OccurredAt, hash); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Claim(ctx, m, hash); !errors.Is(err, conversationruntime.ErrTerminal) {
		t.Errorf("expired replay accepted: %v", err)
	}
	m.ProviderMessageID = "new-history"
	h, err := s.History(ctx, m, 20)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range h {
		if v.UserText == "private old text" {
			t.Error("expired history returned")
		}
	}
	t.Log("CONVERSATION_LEASE_RETENTION_PASS")
}
