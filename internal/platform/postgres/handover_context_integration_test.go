package postgres

import (
	"context"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/randomid"
	"errors"
	"testing"
)

func TestHandoverOperatorContextPostgres(t *testing.T) {
	pool := initialHandoverPool(t)
	ctx := context.Background()
	tenant := initialHandoverFixture(t, pool)
	for _, sql := range []string{`insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref)values($1,'checkout','stripe','FIXTURE_ONLY')`, `insert into payment.provider_checkout(tenant_id,payment_attempt_id,provider_code,connection_id,account_ref,request_sha256_hex,live_mode,expires_at)values($1,'payment','stripe','checkout','acct_fixture',repeat('a',64),false,clock_timestamp()+interval '1 hour')`} {
		if _, err := pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	svc, err := franchisejourney.NewHandoverPreparationService(NewFranchiseJourney(pool), randomid.Generator{}, commercialProfile(t, tenant, 1))
	if err != nil {
		t.Fatal(err)
	}
	v, err := svc.OperatorContext(ctx, tenant, "store", "order")
	if err != nil || !v.CanPrepare || v.Handover != nil || v.PaymentAttemptID != "payment" || v.OrderLineID != "line" || v.ReleaseEffect != franchisejourney.CommercialReleaseEffect {
		t.Fatal("scoped derivation", v, err)
	}
	if _, err = svc.OperatorContext(ctx, tenant, "other", "order"); !errors.Is(err, franchisejourney.ErrReleaseConditioned) {
		t.Fatal("foreign org", err)
	}
	if _, err = svc.OperatorContext(ctx, randomid.Generator{}.New(), "store", "order"); !errors.Is(err, franchisejourney.ErrReleaseConditioned) {
		t.Fatal("foreign tenant", err)
	}
	prepared, _, err := svc.Prepare(ctx, tenant, "operator", initialHandoverCommand())
	if err != nil {
		t.Fatal(err)
	}
	v, err = svc.OperatorContext(ctx, tenant, "store", "order")
	if err != nil || v.CanPrepare || v.Handover == nil || v.Handover.ID != prepared.Handover.ID || v.Handover.Version != 1 {
		t.Fatal("prepared recovery view", v, err)
	}
	if _, err = pool.Exec(ctx, `update payment.provider_observation set hold=true,hold_reason='CALLBACK_PENDING',generation=generation+1 where tenant_id=$1`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.OperatorContext(ctx, tenant, "store", "order"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("held observation exposed as ready", err)
	}
	if _, err = svc.Result(ctx, tenant, "store", "order", initialHandoverCommand().IdempotencyKey); err != nil {
		t.Fatal("historical recovery must remain available", err)
	}
	t.Log("HANDOVER_CONTEXT_PG_PASS exact_order_line_payment_profile=true no_mutation=true foreign_scope_denied=true hold_denied=true historical_recovery_available=true")
}
