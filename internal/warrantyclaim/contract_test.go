package warrantyclaim

import (
	"encoding/json"
	"testing"
	"time"
)

func fixtureProfileDocument() ProfileDocument {
	return ProfileDocument{
		Schema: "elite-warranty-profile/v1", Scope: "MATERIALIZED_PROFILE", Algorithm: "bc-inclusive-fixed-terms", AlgorithmRevision: 1,
		TenantID: "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", OrganizationID: "store", FactoryOrganizationID: "factory",
		PolicyID: "synthetic", TermsVersion: "fixture-v1", TermsText: "Synthetic fixture only; no legal warranty assertion.",
		BusinessTimeZone: "America/New_York", PartsDurationDays: 2, LaborDurationDays: 1, FaultExclusions: []string{"fixture-exclusion"},
		WorkReservationSeconds: 3600, Settlement: "INTERNAL_WARRANTY_SERVICE_ACKNOWLEDGEMENT", AuthorityReference: "fixture-authority", DecisionReference: "fixture-decision",
	}
}
func TestProfileRejectsImplicitOrMutablePolicy(t *testing.T) {
	d := fixtureProfileDocument()
	raw, _ := json.Marshal(d)
	p, err := LoadProfile(raw, SHA(raw))
	if err != nil {
		t.Fatal(err)
	}
	raw[0] = 'X'
	copy := p.Document()
	copy.FaultExclusions[0] = "changed"
	bytes := p.Bytes()
	bytes[0] = 'X'
	if !p.FaultExcluded("fixture-exclusion") || p.FaultExcluded("changed") || p.Bytes()[0] != '{' {
		t.Fatal("profile mutation escaped")
	}
	for _, mutate := range []func(*ProfileDocument){func(d *ProfileDocument) { d.BusinessTimeZone = "Local" }, func(d *ProfileDocument) { d.PartsDurationDays = 0 }, func(d *ProfileDocument) { d.Settlement = "automatic-refund" }, func(d *ProfileDocument) { d.FaultExclusions = []string{"same", "same"} }, func(d *ProfileDocument) { d.AuthorityReference = "" }} {
		d := fixtureProfileDocument()
		mutate(&d)
		raw, _ := json.Marshal(d)
		if _, err := LoadProfile(raw, SHA(raw)); err == nil {
			t.Fatal("invalid profile admitted", string(raw))
		}
	}
	d = fixtureProfileDocument()
	raw, _ = json.Marshal(d)
	duplicate := append([]byte(`{"Schema":"elite-warranty-profile/v1",`), raw[1:]...)
	if _, err := LoadProfile(duplicate, SHA(duplicate)); err == nil {
		t.Fatal("case-insensitive duplicate admitted")
	}
	if _, err := LoadProfile(raw, SHA([]byte("wrong"))); err == nil {
		t.Fatal("unlocked profile admitted")
	}
}
func TestSoldDatesUseBusinessCalendarAndSourceBoundaries(t *testing.T) {
	d := fixtureProfileDocument()
	raw, _ := json.Marshal(d)
	p, err := LoadProfile(raw, SHA(raw))
	if err != nil {
		t.Fatal(err)
	}
	// DST changes at 02:00 local: two calendar days remain March 10, not a
	// floating 48-hour instant. The source predicate includes the ending date.
	dates, err := p.DatesAt(time.Date(2026, 3, 8, 6, 30, 0, 0, time.UTC))
	if err != nil || dates != (Dates{"2026-03-08", "2026-03-10", "2026-03-08", "2026-03-09"}) {
		t.Fatal(dates, err)
	}
	for _, v := range []struct {
		day               string
		any, parts, labor bool
	}{{"2026-03-07", false, false, false}, {"2026-03-08", true, true, true}, {"2026-03-09", true, true, true}, {"2026-03-10", true, true, false}, {"2026-03-11", false, false, false}} {
		c, err := CheckCoverage(v.day, dates)
		if err != nil || c.Any != v.any || c.Parts != v.parts || c.Labor != v.labor {
			t.Fatal(v, c, err)
		}
	}
	if _, err = p.DatesAt(time.Date(9999, 12, 31, 23, 0, 0, 0, time.UTC)); err == nil {
		t.Fatal("out-of-range term admitted")
	}
	if _, err = CheckCoverage("2026-02-30", dates); err == nil {
		t.Fatal("impossible date admitted")
	}
	if _, err = (Profile{}).DatesAt(time.Now()); err == nil {
		t.Fatal("unconfigured profile admitted")
	}
}
