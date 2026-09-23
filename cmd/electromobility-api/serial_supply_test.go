package main

// AUTHORED activation and active-guard boundary; no credential access.
import (
	"context"
	"crypto/sha256"
	sc "elite.local/enterprise/internal/serialsupply"
	"encoding/hex"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
)

func TestSerialSupplyHostDisabledDoesNotReadPolicy(t *testing.T) {
	calls := 0
	module, err := selectedSerialSupplyModule(context.Background(), nil, func(key string) string {
		calls++
		if key != "SERIAL_SUPPLY_ENABLED" {
			t.Fatal("disabled module read", key)
		}
		return "false"
	})
	if err != nil || module != nil || calls != 1 {
		t.Fatal(module, err, calls)
	}
	if _, err = selectedSerialSupplyModule(context.Background(), nil, func(string) string { return "TRUE" }); err == nil {
		t.Fatal("invalid flag enabled")
	}
}
func TestSerialSupplyHostRequiresPolicyAndGuards(t *testing.T) {
	url := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if url == "" {
		t.Skip("owned fixture DB not configured")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	hash := sha256.Sum256([]byte(sc.PolicyJSON))
	want := hex.EncodeToString(hash[:])
	config := map[string]string{"SERIAL_SUPPLY_ENABLED": "true", "SERIAL_SUPPLY_POLICY_SHA256": want}
	lookup := func(key string) string { return config[key] }
	if _, err = selectedSerialSupplyModule(ctx, nil, lookup); err == nil {
		t.Fatal("missing pool admitted")
	}
	config["SERIAL_SUPPLY_POLICY_SHA256"] = ""
	if _, err = selectedSerialSupplyModule(ctx, pool, lookup); err == nil {
		t.Fatal("missing policy hash admitted")
	}
	config["SERIAL_SUPPLY_POLICY_SHA256"] = want
	module, err := selectedSerialSupplyModule(ctx, pool, lookup)
	if err != nil || module == nil {
		t.Fatal("valid schema/profile rejected", err)
	}
	if _, err = pool.Exec(ctx, `alter table approval.request disable trigger serial_supply_approval_guard`); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := pool.Exec(context.Background(), `alter table approval.request enable trigger serial_supply_approval_guard`); err != nil {
			t.Error(err)
		}
	}()
	if _, err = selectedSerialSupplyModule(ctx, pool, lookup); err == nil {
		t.Fatal("disabled approval binding guard admitted")
	}
	t.Log("SERIAL_SUPPLY_HOST_PASS exact policy hash, pool and15active guards required; disabled mode reads no policy")
}
