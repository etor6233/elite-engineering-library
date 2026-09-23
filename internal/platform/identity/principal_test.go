package identity

import "testing"

func TestPrincipalOrganizationScope(t *testing.T) {
	p := Principal{
		Permissions:   map[string]struct{}{"order:create": {}},
		Organizations: map[string]struct{}{"org-a": {}},
	}
	if !p.AllowedOrganization("org-a") {
		t.Fatal("assigned organization rejected")
	}
	if p.AllowedOrganization("org-b") || p.AllowedOrganization("") {
		t.Fatal("unassigned or empty organization accepted")
	}
	admin := Principal{Permissions: map[string]struct{}{"*": {}}}
	if !admin.AllowedOrganization("org-any") {
		t.Fatal("wildcard principal did not receive global organization access")
	}
}
