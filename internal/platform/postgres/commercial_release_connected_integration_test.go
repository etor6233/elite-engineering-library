package postgres_test

// Reuses the already admitted official SDK fixture and actual callback/queue/
// reconciliation owners; adds only the durable commercial checkpoint delta.
import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"elite.local/enterprise/internal/franchisejourney"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
)

func TestCommercialReleaseOfficialSDKConnected(t *testing.T) {
	r := newConnectedRun(t, false)
	ctx := context.Background()
	if code := r.callback(t, "evt_commercial", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
		t.Fatal(code)
	}
	if n, err := r.processor.ProcessOnce(ctx); err != nil || n != 1 {
		t.Fatal("callback reconciliation", n, err)
	}
	var hash string
	if err := r.pool.QueryRow(ctx, `select evidence_sha256_hex from payment.provider_observation where tenant_id=$1`, r.tenant).Scan(&hash); err != nil {
		t.Fatal(err)
	}
	assertConnectedCommercialRelease(t, r, hash)
}

func assertConnectedCommercialRelease(t *testing.T, r *connectedRun, hash string, stored ...bool) {
	assertConnectedCommercialReleaseWithHandoverHook(t, r, hash, nil, stored...)
}

func assertConnectedCommercialReleaseWithHandoverHook(t *testing.T, r *connectedRun, hash string, hook func(*testing.T, *connectedRun, string), stored ...bool) {
	t.Helper()
	ctx := context.Background()
	policy, a := connectedCommercialProfile(t, r, len(stored) == 1 && stored[0])
	repo := db.NewFranchiseJourney(r.pool)
	ids := randomid.Generator{}
	svc, err := franchisejourney.NewHandoverPreparationService(repo, ids, policy)
	if err != nil {
		t.Fatal(err)
	}
	prepared, _, err := svc.Prepare(ctx, r.tenant, "operator", franchisejourney.PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: r.payment.ID, ObservationSHA256: hash, IdempotencyKey: "commercial-integral-prepare"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.PublishDeliveryChecklist(ctx, r.tenant, "operator", franchisejourney.DeliveryChecklist{ID: "commercial-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic delivery", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Read serial", ResponseType: "serial", Required: true}}}, ids.New()); err != nil {
		t.Fatal(err)
	}
	presented, err := repo.CompleteDeliveryChecklist(ctx, r.tenant, "store", "operator", prepared.Handover.ID, 1, "commercial-checklist", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SERIAL-SYNTHETIC"}}, ids.New())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.AcceptHandover(ctx, r.tenant, "store", "customer", prepared.Handover.ID, presented.Version, "SERIAL-SYNTHETIC", "commercial-checklist", 1, strings.Repeat("e", 64), ids.New()); err != nil {
		t.Fatal(err)
	}
	release := franchisejourney.CommitCommercialReleaseCommand{OrganizationID: "store", HandoverID: prepared.Handover.ID, ObservationSHA256: hash, IdempotencyKey: "commercial-integral-release"}
	receipt, replay, err := svc.CommitCommercialRelease(ctx, r.tenant, "operator", release)
	if err != nil || replay || receipt.ContractSHA256 != a.DocumentSHA256 {
		t.Fatal("commercial checkpoint", err)
	}
	current, err := svc.ValidateCommercialRelease(ctx, r.tenant, "store", prepared.Handover.ID)
	if err != nil || !current.Current {
		t.Fatal("checkpoint not current", err)
	}
	if hook != nil {
		hook(t, r, prepared.Handover.ID)
	}
	r.provider.mu.Lock()
	posts, sessionGets, paymentGets := r.provider.posts, r.provider.sessionGets, r.provider.paymentGets
	r.provider.refund = 1
	r.provider.mu.Unlock()
	if posts != 1 || sessionGets != 1 || paymentGets != 1 {
		t.Fatal("unexpected provider requests", posts, sessionGets, paymentGets)
	}
	if code := r.callback(t, "evt_commercial_refund", "charge.refunded", "ch_fixture", true); code != 200 {
		t.Fatal(code)
	}
	current, err = svc.ValidateCommercialRelease(ctx, r.tenant, "store", prepared.Handover.ID)
	if err != nil || current.Current {
		t.Fatal("durable callback hold did not invalidate before GET", err)
	}
	if n, e := r.processor.ProcessOnce(ctx); e != nil || n != 1 {
		t.Fatal("refund GET", n, e)
	}
	current, err = svc.ValidateCommercialRelease(ctx, r.tenant, "store", prepared.Handover.ID)
	if err != nil || current.Current {
		t.Fatal("refunded payment remained eligible", err)
	}
	recovered, err := svc.CommercialReleaseResult(ctx, r.tenant, "store", prepared.Handover.ID, release.IdempotencyKey)
	if err != nil || recovered.ID != receipt.ID {
		t.Fatal("historical recovery", err)
	}
	var receipts, events int
	if err = r.pool.QueryRow(ctx, `select (select count(*) from sales.commercial_release_receipt where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='commercial-release.recorded')`, r.tenant).Scan(&receipts, &events); err != nil || receipts != 1 || events != 1 {
		t.Fatal("duplicate or missing durable effect", err)
	}
	t.Logf("COMMERCIAL_RELEASE_OFFICIAL_SDK_PG_PASS quote_order_allocate_create_checkout_signed_callback_durable_job_get_payment_prepare_checklist_accept_commit_release=true post_callback_current=false refund_current=false profile_materialized=true receipts=1 events=1 provider_posts=1 live_proven=false")
}

func connectedCommercialProfile(t *testing.T, r *connectedRun, storedValue bool) (franchisejourney.HandoverReleaseContract, franchisejourney.HandoverProfileActivation) {
	t.Helper()
	ctx := context.Background()
	python := os.Getenv("HANDOVER_PROFILE_PYTHON")
	if python == "" {
		t.Fatal("explicit admitted Python runtime required")
	}
	out := filepath.Join(t.TempDir(), "policy")
	cmd := exec.CommandContext(ctx, python, "-X", "utf8", "-B", filepath.Join("..", "..", "..", "tools", "materialize_handover_profile.py"), "--output", out, "--profile-id", "franchise-commercial", "--tenant-id", r.tenant, "--organization-id", "store", "--expected-mode", "sandbox", "--payment-provider", "stripe", "--payment-account-ref", "acct_fixture", "--payment-connection-id", "checkout", "--release-effect", "commercial-receipt", "--activate")
	if storedValue {
		cmd.Args = append(cmd.Args, "--funding", "stored-value")
	}
	if body, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("materializer %v %s", err, body)
	}
	raw, err := os.ReadFile(filepath.Join(out, "activation.json"))
	if err != nil {
		t.Fatal(err)
	}
	var a franchisejourney.HandoverProfileActivation
	if err = json.Unmarshal(raw, &a); err != nil {
		t.Fatal(err)
	}
	policy, err := franchisejourney.LoadHandoverProfileFile(filepath.Join(out, "profile.json"), a)
	if err != nil || !policy.AllowsCommercialRelease(r.tenant, "store") {
		t.Fatal("hash-bound commercial profile", err)
	}
	return policy, a
}
