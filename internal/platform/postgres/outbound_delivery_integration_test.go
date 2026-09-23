package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/leadstream"
	"elite.local/enterprise/internal/outbounddelivery"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestOutboundDeliveryFenceReplayUnknownAndReconciliation(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := leadstream.StableUUID(t.Name(), time.Now().UTC().Format(time.RFC3339Nano))
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Outbound','Outbound')`, tenant, "outbound-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	key := []byte("0123456789abcdef0123456789abcdef")
	store, err := NewOutboundDeliveryStore(pool, key, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	message := channels.Message{ChannelCode: "whatsapp", TenantID: tenant, ExternalID: "recipient-secret", ThreadID: "thread", DeliveryKey: strings.Repeat("d", 64), Direction: channels.DirectionOut, Text: "respuesta"}
	hash, err := outbounddelivery.MessageSHA256(message)
	if err != nil {
		t.Fatal(err)
	}
	claim, err := store.Claim(ctx, message, hash)
	if err != nil || claim.Replay {
		t.Fatalf("claim=%+v err=%v", claim, err)
	}
	if _, err = store.Claim(ctx, message, hash); !errors.Is(err, outbounddelivery.ErrInProgress) {
		t.Fatalf("second claim=%v", err)
	}
	receipt := outbounddelivery.Receipt{ProviderMessageID: "wamid.secret", EvidenceSHA256: strings.Repeat("a", 64), AcceptedAt: time.Now().UTC()}
	if err = store.Complete(ctx, message, hash, receipt); err != nil {
		t.Fatal(err)
	}
	claim, err = store.Claim(ctx, message, hash)
	if err != nil || !claim.Replay {
		t.Fatalf("replay=%+v err=%v", claim, err)
	}
	changed := message
	changed.Text = "different"
	changedHash, _ := outbounddelivery.MessageSHA256(changed)
	if _, err = store.Claim(ctx, changed, changedHash); !errors.Is(err, outbounddelivery.ErrConflict) {
		t.Fatalf("divergent=%v", err)
	}

	unknown := message
	unknown.DeliveryKey = strings.Repeat("e", 64)
	unknownHash, _ := outbounddelivery.MessageSHA256(unknown)
	if _, err = store.Claim(ctx, unknown, unknownHash); err != nil {
		t.Fatal(err)
	}
	if err = store.MarkUnknown(ctx, unknown, unknownHash, "PROVIDER_TIMEOUT"); err != nil {
		t.Fatal(err)
	}
	if _, err = store.Claim(ctx, unknown, unknownHash); !errors.Is(err, outbounddelivery.ErrUnknown) {
		t.Fatalf("unknown replay=%v", err)
	}
	if err = store.ReconcileAccepted(ctx, unknown, unknownHash, receipt); err != nil {
		t.Fatal(err)
	}
	claim, err = store.Claim(ctx, unknown, unknownHash)
	if err != nil || !claim.Replay {
		t.Fatalf("reconciled=%+v err=%v", claim, err)
	}
	var attempts, events int
	var leaked bool
	if err = pool.QueryRow(ctx, `select d.attempt_count,count(e.*),d.recipient_hmac like '%recipient-secret%' or d.provider_message_hmac like '%wamid%' from communication.outbound_delivery d join communication.outbound_delivery_event e using(tenant_id,channel_code,delivery_key) where d.tenant_id=$1 and d.delivery_key=$2 group by d.attempt_count,d.recipient_hmac,d.provider_message_hmac`, tenant, message.DeliveryKey).Scan(&attempts, &events, &leaked); err != nil {
		t.Fatal(err)
	}
	if attempts != 1 || events != 2 || leaked {
		t.Fatalf("attempts=%d events=%d leaked=%v", attempts, events, leaked)
	}
	expired := message
	expired.DeliveryKey = strings.Repeat("f", 64)
	expiredHash, _ := outbounddelivery.MessageSHA256(expired)
	if _, err = store.Claim(ctx, expired, expiredHash); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `update communication.outbound_delivery set locked_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and delivery_key=$2`, tenant, expired.DeliveryKey); err != nil {
		t.Fatal(err)
	}
	if _, err = store.Claim(ctx, expired, expiredHash); !errors.Is(err, outbounddelivery.ErrUnknown) {
		t.Fatalf("expired lease=%v", err)
	}
	if err = pool.QueryRow(ctx, `select attempt_count from communication.outbound_delivery where tenant_id=$1 and delivery_key=$2`, tenant, expired.DeliveryKey).Scan(&attempts); err != nil || attempts != 1 {
		t.Fatalf("expired attempts=%d err=%v", attempts, err)
	}
	failed := message
	failed.DeliveryKey = strings.Repeat("9", 64)
	failedHash, _ := outbounddelivery.MessageSHA256(failed)
	if _, err = store.Claim(ctx, failed, failedHash); err != nil {
		t.Fatal(err)
	}
	if err = store.MarkUnknown(ctx, failed, failedHash, "PROVIDER_TIMEOUT"); err != nil {
		t.Fatal(err)
	}
	if err = store.ReconcileFailed(ctx, failed, failedHash, strings.Repeat("8", 64), "PROVIDER_CONFIRMED_FAILED"); err != nil {
		t.Fatal(err)
	}
	if _, err = store.Claim(ctx, failed, failedHash); !errors.Is(err, outbounddelivery.ErrTerminal) {
		t.Fatalf("terminal=%v", err)
	}
}
