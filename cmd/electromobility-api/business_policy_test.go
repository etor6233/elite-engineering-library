package main

// AUTHORED startup tests. No listener, database connection or provider is used.
import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/businesspolicy"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func policyEnvironment(values map[string]string) func(string) string {
	return func(key string) string { return values[key] }
}

func TestBusinessPolicyHostReferenceAndCompleteFile(t *testing.T) {
	reference, err := selectedBusinessPolicy(policyEnvironment(nil))
	if err != nil || reference.SHA256() != businesspolicy.ReferenceSHA256 || reference.LeadTime() != 30*time.Minute {
		t.Fatal("default reference changed", err)
	}
	raw := bytes.ReplaceAll(businesspolicy.ReferenceJSON(), []byte(`"lead_time_seconds": 1800`), []byte(`"lead_time_seconds": 0`))
	hash := fmt.Sprintf("%x", sha256.Sum256(raw))
	file := filepath.Join(t.TempDir(), "business-policy.json")
	if err = os.WriteFile(file, raw, 0600); err != nil {
		t.Fatal(err)
	}
	custom, err := selectedBusinessPolicy(policyEnvironment(map[string]string{"BUSINESS_POLICY_PROFILE_FILE": file, "BUSINESS_POLICY_PROFILE_SHA256": hash}))
	if err != nil || custom.SHA256() != hash || custom.LeadTime() != 0 {
		t.Fatal("custom profile not selected", err)
	}
	// A zero pool is a nonconnecting handle: constructors must not perform I/O.
	runtime, err := prepareBusinessPolicyRuntime(&pgxpool.Pool{}, randomid.Generator{}, systemClock{}, custom)
	if err != nil || runtime == nil || runtime.profile != custom || runtime.commerceService == nil || runtime.journeyService == nil {
		t.Fatal("owner construction", err)
	}
	if runtime.commerceRepository.BusinessPolicySHA256() != custom.SHA256() || runtime.journeyRepository.BusinessPolicySHA256() != custom.SHA256() {
		t.Fatal("owners received different profiles")
	}
}

func TestBusinessPolicyHostRejectsIncompleteInvalidAndMissingFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "business-policy.json")
	if err := os.WriteFile(file, businesspolicy.ReferenceJSON(), 0600); err != nil {
		t.Fatal(err)
	}
	cases := map[string]map[string]string{
		"path-only":       {"BUSINESS_POLICY_PROFILE_FILE": file},
		"hash-only":       {"BUSINESS_POLICY_PROFILE_SHA256": businesspolicy.ReferenceSHA256},
		"wrong-hash":      {"BUSINESS_POLICY_PROFILE_FILE": file, "BUSINESS_POLICY_PROFILE_SHA256": strings.Repeat("0", 64)},
		"missing-file":    {"BUSINESS_POLICY_PROFILE_FILE": file + ".absent", "BUSINESS_POLICY_PROFILE_SHA256": businesspolicy.ReferenceSHA256},
		"whitespace-path": {"BUSINESS_POLICY_PROFILE_FILE": " ", "BUSINESS_POLICY_PROFILE_SHA256": businesspolicy.ReferenceSHA256},
	}
	for name, values := range cases {
		t.Run(name, func(t *testing.T) {
			if p, err := selectedBusinessPolicy(policyEnvironment(values)); err == nil || p != nil {
				t.Fatal("invalid configuration fell back to reference")
			}
		})
	}
	invalid := bytes.ReplaceAll(businesspolicy.ReferenceJSON(), []byte(`"maximum_slot_capacity": 100`), []byte(`"maximum_slot_capacity": 101`))
	if err := os.WriteFile(file, invalid, 0600); err != nil {
		t.Fatal(err)
	}
	if p, err := selectedBusinessPolicy(policyEnvironment(map[string]string{"BUSINESS_POLICY_PROFILE_FILE": file, "BUSINESS_POLICY_PROFILE_SHA256": fmt.Sprintf("%x", sha256.Sum256(invalid))})); err == nil || p != nil {
		t.Fatal("hash-matching incompatible policy admitted")
	}
	if _, err := selectedBusinessPolicy(nil); err == nil {
		t.Fatal("nil environment accessor")
	}
	for _, p := range []*businesspolicy.Profile{nil, new(businesspolicy.Profile)} {
		if runtime, err := prepareBusinessPolicyRuntime(&pgxpool.Pool{}, randomid.Generator{}, systemClock{}, p); err == nil || runtime != nil {
			t.Fatal("invalid profile reached constructed owners")
		}
	}
	if _, err := prepareBusinessPolicyRuntime(nil, randomid.Generator{}, systemClock{}, businesspolicy.Reference()); err == nil {
		t.Fatal("nil persistence handle")
	}
	if _, err := prepareBusinessPolicyRuntime(&pgxpool.Pool{}, nil, systemClock{}, businesspolicy.Reference()); err == nil {
		t.Fatal("nil id generator")
	}
	if _, err := prepareBusinessPolicyRuntime(&pgxpool.Pool{}, randomid.Generator{}, nil, businesspolicy.Reference()); err == nil {
		t.Fatal("nil clock")
	}
}
