package franchisejourney

// AUTHORED verification glue. The seed is the materialized sandbox profile
// used by the connected PostgreSQL journey, with synthetic account identifiers.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func fuzzHandoverProfileSeed(mode bool) ([]byte, HandoverProfileActivation) {
	d := HandoverProfileDocument{
		Schema: HandoverProfileSchema, ProfileID: "franchise-reference", Revision: 1,
		Algorithm: HandoverSupportedAlgorithm, AlgorithmRevision: 1,
		Scope: "MATERIALIZED_PROFILE", TenantID: "018f4d4a-7b36-7a21-8d10-000000000001", OrganizationID: "store",
		ProviderCode: "stripe", ProviderAccountRef: "acct_fixture", ProviderConnectionID: "checkout", ExpectedLiveMode: &mode,
		MaximumObservationAgeSeconds: 300, Options: SupportedHandoverProfileOptions(),
		AuthorityReference: "docs/initial-handover-reference.md", DecisionReference: "PROJECT_HANDOVER_POLICY_DECISION.md",
	}
	raw, _ := json.Marshal(d)
	a := HandoverProfileActivation{Enabled: true, ProfileID: d.ProfileID, ProfileRevision: d.Revision, TenantID: d.TenantID, OrganizationID: d.OrganizationID,
		ProviderCode: d.ProviderCode, ProviderAccountRef: d.ProviderAccountRef, ProviderConnectionID: d.ProviderConnectionID, ExpectedLiveMode: mode}
	return raw, a
}

func fuzzProfileHash(raw []byte) string {
	digest := sha256.Sum256(raw)
	return hex.EncodeToString(digest[:])
}

func FuzzHandoverProfileAdmission(f *testing.F) {
	base, _ := fuzzHandoverProfileSeed(false)
	live, _ := fuzzHandoverProfileSeed(true)
	f.Add(base, uint8(0))
	f.Add(live, uint8(1))
	var commercial HandoverProfileDocument
	_ = json.Unmarshal(base, &commercial)
	commercial.AlgorithmRevision = 2
	commercial.Options = SupportedCommercialReleaseOptions()
	commercialRaw, _ := json.Marshal(commercial)
	f.Add(commercialRaw, uint8(0))
	for i, seed := range [][]byte{base, live, commercialRaw} {
		_, activation := fuzzHandoverProfileSeed(i == 1)
		activation.DocumentSHA256 = fuzzProfileHash(seed)
		if contract, err := LoadHandoverProfile(seed, activation); err != nil || !contract.Valid() {
			f.Fatalf("real supported profile seed %d was rejected: %v", i, err)
		}
	}
	for variant := uint8(2); variant < 10; variant++ {
		f.Add(base, variant)
	}
	f.Add(append([]byte(`{"schema":"duplicate",`), base[1:]...), uint8(0))
	f.Add(bytes.Replace(base, []byte(`"expected_live_mode":false,`), nil, 1), uint8(0))
	f.Add(bytes.Replace(base, []byte(`"scope":"MATERIALIZED_PROFILE"`), []byte(`"scope":"LOCAL_FIXTURES"`), 1), uint8(0))
	f.Add([]byte(`{"policy":true}`), uint8(0))
	f.Add(bytes.Replace(base, []byte(`"algorithm_revision":1`), []byte(`"algorithm_revision":999`), 1), uint8(0))
	f.Fuzz(func(t *testing.T, raw []byte, variant uint8) {
		if len(raw) > 32769 {
			t.Skip()
		}
		variant %= 10
		_, a := fuzzHandoverProfileSeed(variant == 1)
		a.DocumentSHA256 = fuzzProfileHash(raw)
		switch variant {
		case 2:
			a.DocumentSHA256 = strings.Repeat("0", 64)
		case 3:
			a.TenantID = "018f4d4a-7b36-7a21-8d10-000000000099"
		case 4:
			a.OrganizationID = "another-store"
		case 5:
			a.ProviderAccountRef = "another-account"
		case 6:
			a.ProviderConnectionID = "another-connection"
		case 7:
			a.ProviderCode = "mercadopago"
		case 8:
			a.Enabled = false
		case 9:
			a.ProfileRevision = 0
		}
		contract, err := LoadHandoverProfile(raw, a)
		if err != nil {
			if contract.Valid() {
				t.Fatal("rejected profile leaked an admitted contract")
			}
			return
		}
		// Locks are external inputs. A mutated document may legitimately select
		// any supported revision, but must bind exactly to this activation.
		if !contract.Valid() || !contract.AllowsScope(a.TenantID, a.OrganizationID) || contract.AllowsScope(a.TenantID+"-foreign", a.OrganizationID) || contract.AllowsScope(a.TenantID, a.OrganizationID+"-foreign") || contract.ExpectedLiveMode != a.ExpectedLiveMode || contract.DocumentSHA256 != fuzzProfileHash(raw) {
			t.Fatal("profile escaped its explicit activation")
		}
		provider, account, connection, required := contract.PaymentBinding()
		if !required || provider != a.ProviderCode || account != a.ProviderAccountRef || connection != a.ProviderConnectionID {
			t.Fatal("profile lost its payment identity binding")
		}
		for _, mutate := range []func(*HandoverProfileActivation){
			func(p *HandoverProfileActivation) { p.Enabled = false },
			func(p *HandoverProfileActivation) { p.ExpectedLiveMode = !p.ExpectedLiveMode },
			func(p *HandoverProfileActivation) { p.TenantID += "-foreign" },
			func(p *HandoverProfileActivation) { p.OrganizationID += "-foreign" },
			func(p *HandoverProfileActivation) { p.ProviderAccountRef += "-foreign" },
			func(p *HandoverProfileActivation) { p.ProviderConnectionID += "-foreign" },
			func(p *HandoverProfileActivation) { p.ProfileID += "-foreign" },
			func(p *HandoverProfileActivation) { p.ProfileRevision++ },
		} {
			changed := a
			mutate(&changed)
			if _, err := LoadHandoverProfile(raw, changed); err == nil {
				t.Fatal("same document was admitted under another activation")
			}
		}
		// Admission is process-local: public JSON must not forge the private
		// binding, and changes to the admitted contract must invalidate it.
		wire, err := json.Marshal(contract)
		if err != nil {
			t.Fatal(err)
		}
		var forged HandoverReleaseContract
		if json.Unmarshal(wire, &forged) != nil || forged.Valid() {
			t.Fatal("serialized contract bypasses profile admission")
		}
		for _, mutate := range []func(*HandoverReleaseContract){
			func(p *HandoverReleaseContract) { p.ExpectedLiveMode = !p.ExpectedLiveMode },
			func(p *HandoverReleaseContract) { p.DocumentSHA256 = strings.Repeat("f", 64) },
			func(p *HandoverReleaseContract) { p.MaximumObservationAge += time.Second },
			func(p *HandoverReleaseContract) { p.Scope = "LOCAL_FIXTURES" },
		} {
			changed := contract
			mutate(&changed)
			if changed.Valid() {
				t.Fatal("public contract mutation retained admission")
			}
		}
		changedBytes := append(append([]byte(nil), raw...), ' ')
		if _, err := LoadHandoverProfile(changedBytes, a); err == nil {
			t.Fatal("profile lock accepted different bytes")
		}
		// A valid semantic roundtrip requires an explicitly updated byte lock.
		var document HandoverProfileDocument
		if json.Unmarshal(raw, &document) != nil {
			t.Fatal("accepted profile is not a decodable document")
		}
		normalized, err := json.Marshal(document)
		if err != nil {
			t.Fatal(err)
		}
		a.DocumentSHA256 = fuzzProfileHash(normalized)
		restored, err := LoadHandoverProfile(normalized, a)
		if err != nil || !restored.Valid() || !restored.AllowsScope(a.TenantID, a.OrganizationID) || restored.ID != contract.ID || restored.ExpectedLiveMode != contract.ExpectedLiveMode {
			t.Fatal("explicitly relocked semantic roundtrip changes admission")
		}
		// Duplicate properties must be rejected even when their bytes are
		// locked: two reviewers must not see different policy declarations.
		duplicate := append([]byte(`{"schema":"duplicate",`), normalized[1:]...)
		a.DocumentSHA256 = fuzzProfileHash(duplicate)
		if _, err := LoadHandoverProfile(duplicate, a); err == nil {
			t.Fatal("duplicate property accepted under an exact byte lock")
		}
	})
}
