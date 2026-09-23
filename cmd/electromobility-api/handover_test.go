package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/paymentbridge"
)

func handoverHostFixture(t *testing.T) (*paymentRuntime, map[string]string) {
	t.Helper()
	scope := paymentbridge.Scope{TenantID: "018f4d4a-7b36-7a21-8d10-000000000001", OrganizationID: "store", ConnectionID: "connection", ProviderCode: "stripe", AccountRef: "acct_fixture", Currency: "ARS", MinorUnitExponent: 2}
	mode := false
	document := franchisejourney.HandoverProfileDocument{Schema: franchisejourney.HandoverProfileSchema, ProfileID: "franchise-host-fixture", Revision: 2, Algorithm: franchisejourney.HandoverSupportedAlgorithm, AlgorithmRevision: franchisejourney.HandoverSupportedAlgorithmRevision, Scope: "MATERIALIZED_PROFILE", TenantID: scope.TenantID, OrganizationID: scope.OrganizationID, ExpectedLiveMode: &mode, ProviderCode: scope.ProviderCode, ProviderAccountRef: scope.AccountRef, ProviderConnectionID: scope.ConnectionID, MaximumObservationAgeSeconds: 120, Options: franchisejourney.SupportedHandoverProfileOptions(), AuthorityReference: "docs/handover-profile.md", DecisionReference: "PROJECT_HANDOVER_POLICY_DECISION.md"}
	raw, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "profile.json")
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(raw)
	env := map[string]string{"HANDOVER_ENABLED": "true", "HANDOVER_PROFILE_FILE": path, "HANDOVER_PROFILE_ID": document.ProfileID, "HANDOVER_PROFILE_REVISION": "2", "HANDOVER_PROFILE_SHA256": hex.EncodeToString(digest[:])}
	return &paymentRuntime{worker: paymentbridge.Worker{Scope: scope}}, env
}
func TestHandoverHostDisabledDoesNotNeedPaymentOrProfile(t *testing.T) {
	for _, value := range []string{"", "false"} {
		lookup := func(key string) string {
			if key == "HANDOVER_ENABLED" {
				return value
			}
			return "missing-profile"
		}
		module, err := selectedInitialHandoverModule(nil, nil, lookup)
		if err != nil || module != nil {
			t.Fatal("disabled owner attempted activation")
		}
	}
}
func TestHandoverHostBindsMaterializedProfileToPaymentScope(t *testing.T) {
	runtime, env := handoverHostFixture(t)
	contract, err := loadInitialHandoverContract(runtime, func(k string) string { return env[k] })
	if err != nil || contract == nil || contract.Scope != "MATERIALIZED_PROFILE" || contract.ID != "franchise-host-fixture@2" || !contract.AllowsScope(runtime.worker.Scope.TenantID, runtime.worker.Scope.OrganizationID) {
		t.Fatal(contract, err)
	}
	if contract.AllowsScope(runtime.worker.Scope.TenantID, "other") {
		t.Fatal("foreign organization admitted")
	}
	// Missing payment runtime is rejected before any profile-backed route exists.
	if _, err = selectedInitialHandoverModule(nil, nil, func(k string) string { return env[k] }); err == nil {
		t.Fatal("handover activated without payment runtime")
	}
}
func TestHandoverHostRejectsIncompleteOrMismatchedActivation(t *testing.T) {
	runtime, env := handoverHostFixture(t)
	for key := range env {
		if key == "HANDOVER_ENABLED" {
			continue
		}
		t.Run("missing-"+key, func(t *testing.T) {
			lookup := func(k string) string {
				if k == key {
					return ""
				}
				return env[k]
			}
			if _, err := loadInitialHandoverContract(runtime, lookup); err == nil {
				t.Fatal("incomplete activation accepted")
			}
		})
	}
	for _, tc := range []struct{ key, value string }{
		{"HANDOVER_ENABLED", "yes"}, {"HANDOVER_PROFILE_ID", "other-profile"}, {"HANDOVER_PROFILE_REVISION", "3"}, {"HANDOVER_PROFILE_SHA256", env["HANDOVER_PROFILE_SHA256"][:63] + "x"}, {"HANDOVER_PROFILE_FILE", filepath.Join(t.TempDir(), "absent.json")},
	} {
		t.Run(tc.key+"-mismatch", func(t *testing.T) {
			lookup := func(k string) string {
				if k == tc.key {
					return tc.value
				}
				return env[k]
			}
			if _, err := loadInitialHandoverContract(runtime, lookup); err == nil {
				t.Fatal("mismatched activation accepted")
			}
		})
	}
	for _, field := range []string{"tenant", "organization", "provider", "account", "connection", "mode"} {
		t.Run("payment-"+field+"-mismatch", func(t *testing.T) {
			altered := *runtime
			switch field {
			case "tenant":
				altered.worker.Scope.TenantID = "018f4d4a-7b36-7a21-8d10-000000000002"
			case "organization":
				altered.worker.Scope.OrganizationID = "other"
			case "provider":
				altered.worker.Scope.ProviderCode = "mercadopago"
			case "account":
				altered.worker.Scope.AccountRef = "acct_other"
			case "connection":
				altered.worker.Scope.ConnectionID = "other"
			case "mode":
				altered.worker.Scope.LiveMode = true
			}
			if _, err := loadInitialHandoverContract(&altered, func(k string) string { return env[k] }); err == nil {
				t.Fatal("payment scope differs from locked profile")
			}
		})
	}
	// A valid activation never silently falls back to the old LOCAL_FIXTURES literal.
	if err := os.WriteFile(env["HANDOVER_PROFILE_FILE"], []byte(franchisejourney.ReferenceHandoverContractDocument), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadInitialHandoverContract(runtime, func(k string) string { return env[k] }); err == nil {
		t.Fatal("non-materialized fallback accepted")
	}
}
