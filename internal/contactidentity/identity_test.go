package contactidentity

import (
	"strings"
	"testing"
	"time"
)

func validCommand() Command {
	return Command{TenantID: "tenant-a", ChannelCode: "whatsapp", ExternalID: "+5491112345678", LeadID: "lead-1", SubjectID: "lead:lead-1", PIIAllowed: true, State: StateActive, PolicyVersion: "sales-contact-v1", EvidenceSHA256: strings.Repeat("a", 64), EffectiveAt: time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC), RequestID: "bind-1"}
}

func TestDigestIsScopedDeterministicAndDoesNotRevealExternalID(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	a, err := ExternalDigest(key, "tenant-a", "whatsapp", "+5491112345678")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := ExternalDigest(key, "tenant-a", "whatsapp", "+5491112345678")
	c, _ := ExternalDigest(key, "tenant-b", "whatsapp", "+5491112345678")
	if a != b || a == c || strings.Contains(a, "12345678") || len(a) != 64 {
		t.Fatalf("a=%q b=%q c=%q", a, b, c)
	}
}

func TestCommandValidationFailsClosed(t *testing.T) {
	c := validCommand()
	if err := c.Validate(); err != nil {
		t.Fatal(err)
	}
	c.State, c.PIIAllowed = StateRevoked, true
	if err := c.Validate(); err == nil {
		t.Fatal("revoked binding cannot allow PII")
	}
	c = validCommand()
	c.EvidenceSHA256 = "unproven"
	if err := c.Validate(); err == nil {
		t.Fatal("evidence hash required")
	}
	if _, err := ExternalDigest([]byte("short"), "tenant-a", "whatsapp", "id"); err == nil {
		t.Fatal("weak key accepted")
	}
}

func TestRequestHashChangesWithPolicyAndHidesRawIdentity(t *testing.T) {
	c := validCommand()
	key := []byte("0123456789abcdef0123456789abcdef")
	_, a, err := c.RequestSHA256(key)
	if err != nil {
		t.Fatal(err)
	}
	c.PolicyVersion = "sales-contact-v2"
	_, b, err := c.RequestSHA256(key)
	if err != nil {
		t.Fatal(err)
	}
	if a == b || strings.Contains(a, c.ExternalID) {
		t.Fatalf("a=%q b=%q", a, b)
	}
}
