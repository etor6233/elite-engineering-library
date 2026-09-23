package leadpromotion

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func validCommand() Command {
	return Command{TenantID: "tenant", Provider: "google_ads", ProviderEventID: "event-1", LeadID: "lead-1", ConsentID: "consent-1", PurposeCode: "sales-contact", PolicyVersion: "policy-7", EvidenceSHA256: strings.Repeat("a", 64), DecisionAt: time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC), MappingVersion: "google-form-v1", FieldMapping: map[string]string{"EMAIL": "email"}, IdempotencyKey: "promote-0001"}
}

func TestCommandRequiresExplicitPolicyAndContactMapping(t *testing.T) {
	c := validCommand()
	c.PolicyVersion = ""
	if !errors.Is(c.Validate(), ErrInvalid) {
		t.Fatal("missing policy accepted")
	}
	c = validCommand()
	c.FieldMapping = map[string]string{"NAME": "name"}
	if !errors.Is(c.Validate(), ErrInvalid) {
		t.Fatal("mapping without contact accepted")
	}
}

func TestContactPayloadUsesOnlyConfiguredFields(t *testing.T) {
	raw, err := ContactPayload([]Field{{ID: "EMAIL", Value: "person@example.test"}, {ID: "UNMAPPED", Value: "secret"}}, map[string]string{"EMAIL": "email"})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != `{"email":"person@example.test"}` || strings.Contains(string(raw), "secret") {
		t.Fatalf("unexpected payload %s", raw)
	}
}

func TestContactPayloadRejectsMissingConfiguredField(t *testing.T) {
	_, err := ContactPayload([]Field{{ID: "PHONE", Value: "+541100000000"}}, map[string]string{"EMAIL": "email"})
	if !errors.Is(err, ErrFieldAbsent) {
		t.Fatalf("expected absent field, got %v", err)
	}
}

type fakeStore struct{ called bool }

func (f *fakeStore) Promote(_ context.Context, c Command) (Receipt, error) {
	f.called = true
	return Receipt{LeadID: c.LeadID}, nil
}

func TestServiceValidatesBeforeStore(t *testing.T) {
	f := &fakeStore{}
	s := Service{Store: f}
	c := validCommand()
	c.EvidenceSHA256 = "bad"
	if _, err := s.Promote(context.Background(), c); !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid, got %v", err)
	}
	if f.called {
		t.Fatal("store called for invalid command")
	}
}
