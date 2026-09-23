package franchisejourney

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func profileFixture(mode bool) (HandoverProfileDocument, HandoverProfileActivation) {
	d := HandoverProfileDocument{Schema: HandoverProfileSchema, ProfileID: "franchise-reference", Revision: 1, Algorithm: HandoverSupportedAlgorithm, AlgorithmRevision: 1, Scope: "MATERIALIZED_PROFILE", TenantID: "018f4d4a-7b36-7a21-8d10-000000000001", OrganizationID: "store", ProviderCode: "stripe", ProviderAccountRef: "acct_fixture", ProviderConnectionID: "checkout", ExpectedLiveMode: &mode, MaximumObservationAgeSeconds: 300, Options: SupportedHandoverProfileOptions(), AuthorityReference: "docs/initial-handover-reference.md", DecisionReference: "PROJECT_HANDOVER_POLICY_DECISION.md"}
	a := HandoverProfileActivation{Enabled: true, ProfileID: d.ProfileID, ProfileRevision: d.Revision, TenantID: d.TenantID, OrganizationID: d.OrganizationID, ProviderCode: d.ProviderCode, ProviderAccountRef: d.ProviderAccountRef, ProviderConnectionID: d.ProviderConnectionID, ExpectedLiveMode: mode}
	return d, a
}
func lockProfile(t *testing.T, d HandoverProfileDocument, a HandoverProfileActivation) ([]byte, HandoverProfileActivation) {
	t.Helper()
	raw, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	hash := sha256.Sum256(raw)
	a.DocumentSHA256 = hex.EncodeToString(hash[:])
	return raw, a
}
func TestHandoverProfileBindsAlgorithmRevisionModeAndScope(t *testing.T) {
	for _, mode := range []bool{false, true} {
		d, a := profileFixture(mode)
		raw, a := lockProfile(t, d, a)
		contract, err := LoadHandoverProfile(raw, a)
		if err != nil || !contract.Valid() || contract.ExpectedLiveMode != mode || !contract.AllowsScope(d.TenantID, "store") || contract.AllowsScope(d.TenantID, "foreign") {
			t.Fatal("profile", contract, err)
		}
		path := filepath.Join(t.TempDir(), "profile.json")
		if err = os.WriteFile(path, raw, 0600); err != nil {
			t.Fatal(err)
		}
		loaded, err := LoadHandoverProfileFile(path, a)
		if err != nil || loaded.DocumentSHA256 != contract.DocumentSHA256 {
			t.Fatal("file loader", err)
		}
		forged := contract
		forged.ExpectedLiveMode = !mode
		if forged.Valid() {
			t.Fatal("mode changed after admission")
		}
		repo := &preparationRepoProbe{}
		service, err := NewHandoverPreparationService(repo, preparationIDs{}, contract)
		if err != nil {
			t.Fatal(err)
		}
		command := PrepareHandoverCommand{OrganizationID: "foreign", OrderID: "order", OrderLineID: "line", PaymentAttemptID: "payment", ObservationSHA256: strings.Repeat("a", 64), IdempotencyKey: "profile-key"}
		if _, _, err = service.Prepare(context.Background(), d.TenantID, "operator", command); !errors.Is(err, ErrReleaseConditioned) || repo.calls != 0 {
			t.Fatal("scope bypass", err)
		}
		// JSON cannot manufacture the loader's private immutable admission binding.
		encoded, _ := json.Marshal(contract)
		var roundtrip HandoverReleaseContract
		json.Unmarshal(encoded, &roundtrip)
		if roundtrip.Valid() {
			t.Fatal("unadmitted reconstructed contract accepted")
		}
	}
}
func TestHandoverProfileRejectsAlteredMissingAndUnsupportedPolicy(t *testing.T) {
	edits := map[string]func(*HandoverProfileDocument, *HandoverProfileActivation){
		"provider-binding":   func(_ *HandoverProfileDocument, a *HandoverProfileActivation) { a.ProviderCode = "mercadopago" },
		"account-binding":    func(_ *HandoverProfileDocument, a *HandoverProfileActivation) { a.ProviderAccountRef = "acct_other" },
		"connection-binding": func(_ *HandoverProfileDocument, a *HandoverProfileActivation) { a.ProviderConnectionID = "other" },
		"missing-provider":   func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.ProviderCode = "" },
		"disabled":           func(_ *HandoverProfileDocument, a *HandoverProfileActivation) { a.Enabled = false },
		"schema":             func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.Schema = "v2" },
		"algorithm":          func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.Algorithm = "arbitrary-policy" },
		"algorithm-revision": func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.AlgorithmRevision = 2 },
		"profile-revision":   func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.Revision = 2 },
		"scope":              func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.Scope = "PRODUCTION_PROVEN" },
		"mode-absent":        func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.ExpectedLiveMode = nil },
		"mode-binding":       func(_ *HandoverProfileDocument, a *HandoverProfileActivation) { a.ExpectedLiveMode = true },
		"quantity":           func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.Options.Quantity = 2 },
		"credit":             func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.Options.PaymentCoverage = "PARTIAL" },
		"shipping": func(d *HandoverProfileDocument, _ *HandoverProfileActivation) {
			d.Options.ReleaseEffect = "POST_SHIPMENT"
		},
		"missing-decision": func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.DecisionReference = "" },
		"missing-tenant":   func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.TenantID = "" },
		"tenant-binding": func(_ *HandoverProfileDocument, a *HandoverProfileActivation) {
			a.TenantID = "018f4d4a-7b36-7a21-8d10-000000000002"
		},
		"stale-age": func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.MaximumObservationAgeSeconds = 901 },
	}
	for name, edit := range edits {
		t.Run(name, func(t *testing.T) {
			d, a := profileFixture(false)
			edit(&d, &a)
			raw, a := lockProfile(t, d, a)
			if _, err := LoadHandoverProfile(raw, a); !errors.Is(err, ErrReleaseConditioned) {
				t.Fatal("unsupported policy admitted", err)
			}
		})
	}
	d, a := profileFixture(false)
	raw, a := lockProfile(t, d, a)
	if _, err := LoadHandoverProfile(append(raw, ' '), a); err == nil {
		t.Fatal("byte lock ignored")
	}
	if _, err := LoadHandoverProfile(nil, a); err == nil {
		t.Fatal("absent policy admitted")
	}
	if _, err := LoadHandoverProfileFile(filepath.Join(t.TempDir(), "absent"), a); err == nil {
		t.Fatal("absent policy file admitted")
	}
	for _, changed := range [][]byte{[]byte(`{"policy":true}`), []byte(strings.Replace(string(raw), `"revision":1`, `"revision":1,"revision":1`, 1)), append(raw, []byte(`{}`)...)} {
		hash := sha256.Sum256(changed)
		b := a
		b.DocumentSHA256 = hex.EncodeToString(hash[:])
		if _, err := LoadHandoverProfile(changed, b); err == nil {
			t.Fatal("unknown/ambiguous JSON admitted")
		}
	}
}
