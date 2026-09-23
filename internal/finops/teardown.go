package finops

// TeardownPlan accounts for every tenant resource before teardown. Nothing may
// escape: teardown is safe only when the audit is clean (no untagged, no
// drift, no leak).
type TeardownPlan struct {
	Resources []Resource
	Audit     AuditResult
}

// PlanTeardown builds the teardown plan for a tenant from the actual cloud
// resource set. Complete means the audit is clean and every resource is
// accounted for.
func PlanTeardown(inv *Inventory, tenant string, actual []Resource) TeardownPlan {
	return TeardownPlan{
		Resources: inv.List(tenant),
		Audit:     inv.Audit(tenant, actual),
	}
}

// Complete reports whether teardown can proceed without escaping anything.
func (p TeardownPlan) Complete() bool {
	return p.Audit.Clean()
}
