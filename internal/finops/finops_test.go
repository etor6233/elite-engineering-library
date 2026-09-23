package finops

import (
	"errors"
	"testing"
)

func res(tenant, provider, id string, tags map[string]string) Resource {
	if tags == nil {
		tags = map[string]string{"tenant": tenant, "environment": "prod", "owner": "platform"}
	}
	return Resource{
		TenantID: tenant, Provider: provider, ID: id, Kind: "cloud_run",
		Region: "us-central1", Tags: tags, EstMonthlyMinorUnits: 100,
	}
}

func TestResourceValidateTags(t *testing.T) {
	if err := res("t", "gcp", "r1", nil).Validate(); err != nil {
		t.Fatalf("valid resource rejected: %v", err)
	}
	bad := res("t", "gcp", "r1", map[string]string{"tenant": "t"}) // missing env/owner
	if err := bad.Validate(); !errors.Is(err, ErrMissingTag) {
		t.Fatalf("missing tag accepted: %v", err)
	}
	neg := res("t", "gcp", "r1", nil)
	neg.EstMonthlyMinorUnits = -1
	if err := neg.Validate(); !errors.Is(err, ErrInvalidResource) {
		t.Fatalf("negative cost accepted: %v", err)
	}
}

func TestInventoryDedupAndCost(t *testing.T) {
	inv := NewInventory()
	if err := inv.Register(res("t", "gcp", "r1", nil)); err != nil {
		t.Fatal(err)
	}
	if err := inv.Register(res("t", "gcp", "r1", nil)); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate accepted: %v", err)
	}
	_ = inv.Register(res("t", "aws", "r2", nil))
	if got := inv.MonthlyCost("t"); got != 200 {
		t.Fatalf("expected 200, got %d", got)
	}
	if len(inv.List("t")) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(inv.List("t")))
	}
}

func TestInventoryRejectsUntagged(t *testing.T) {
	inv := NewInventory()
	// Register rejects untagged fail-closed
	if err := inv.Register(res("t", "aws", "r2", map[string]string{"tenant": "t"})); !errors.Is(err, ErrMissingTag) {
		t.Fatalf("untagged resource accepted: %v", err)
	}
}

func TestAuditDetectsEscapes(t *testing.T) {
	inv := NewInventory()
	_ = inv.Register(res("t", "gcp", "r1", nil))

	// actual: r1 present (good), r2 untagged, r3 drifted-in; r1 expected but a
	// leak scenario is tested separately.
	actual := []Resource{
		res("t", "gcp", "r1", nil),
		Resource{TenantID: "t", Provider: "aws", ID: "r2", Kind: "cloud_run", Region: "us-east-1", Tags: map[string]string{"tenant": "t"}}, // untagged
		res("t", "gcp", "r3", nil), // drifted-in (not inventoried)
	}
	a := inv.Audit("t", actual)
	if len(a.Untagged) != 1 || len(a.Unknown) != 1 {
		t.Fatalf("expected 1 untagged + 1 unknown, got %+v", a)
	}
	if a.Clean() {
		t.Fatal("audit should not be clean")
	}
}

func TestAuditDetectsLeak(t *testing.T) {
	inv := NewInventory()
	_ = inv.Register(res("t", "gcp", "r1", nil))
	_ = inv.Register(res("t", "aws", "r2", nil))

	// actual only has r1 → r2 leaked out of inventory
	a := inv.Audit("t", []Resource{res("t", "gcp", "r1", nil)})
	if len(a.Missing) != 1 || a.Missing[0].ID != "r2" {
		t.Fatalf("expected r2 missing, got %+v", a.Missing)
	}
}

func TestBudgetFailClosed(t *testing.T) {
	b := NewBudget()
	b.SetCap("t", 100)
	if err := b.Record("t", 60); err != nil {
		t.Fatal(err)
	}
	if b.Remaining("t") != 40 {
		t.Fatalf("expected 40 remaining, got %d", b.Remaining("t"))
	}
	if err := b.Record("t", 41); !errors.Is(err, ErrOverBudget) {
		t.Fatalf("expected ErrOverBudget, got %v", err)
	}
	if b.OverBudget("t") {
		t.Fatal("should not be over budget at 60/100")
	}
	if err := b.Record("t", 40); err != nil {
		t.Fatal(err)
	}
	if !b.OverBudget("t") {
		t.Fatal("expected over-budget when spent == cap")
	}
}

func TestTeardownPlanAccountsAll(t *testing.T) {
	inv := NewInventory()
	_ = inv.Register(res("t", "gcp", "r1", nil))

	// actual matches inventory → clean teardown
	actual := []Resource{res("t", "gcp", "r1", nil)}
	if !PlanTeardown(inv, "t", actual).Complete() {
		t.Fatal("plan should be complete when actual == inventory")
	}
	// actual has an untagged drift → not complete
	drift := []Resource{
		res("t", "gcp", "r1", nil),
		Resource{TenantID: "t", Provider: "aws", ID: "x", Kind: "compute", Region: "us-east-1", Tags: map[string]string{"tenant": "t"}},
	}
	if PlanTeardown(inv, "t", drift).Complete() {
		t.Fatal("plan should be incomplete when something escaped")
	}
}
