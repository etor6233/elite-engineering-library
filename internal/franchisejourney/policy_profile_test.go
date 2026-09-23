package franchisejourney

// AUTHORED profile-to-owner binding regression tests.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/businesspolicy"
	"fmt"
	"testing"
	"time"
)

type policyBoundFake struct {
	fakeRepository
	hash string
}

func (r *policyBoundFake) BusinessPolicySHA256() string { return r.hash }
func TestJourneyPolicyConfigurationAndBinding(t *testing.T) {
	raw := bytes.ReplaceAll(businesspolicy.ReferenceJSON(), []byte(`"lead_time_seconds": 1800`), []byte(`"lead_time_seconds": 0`))
	p, err := businesspolicy.Load(raw, fmt.Sprintf("%x", sha256.Sum256(raw)))
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 9, 11, 12, 0, 0, 0, time.UTC)
	repo := &policyBoundFake{hash: p.SHA256()}
	service, err := NewServiceWithProfile(repo, &fixedIDs{}, fixedClock{now}, p)
	if err != nil {
		t.Fatal(err)
	}
	slot := AppointmentSlot{OrganizationID: "org", Kind: "consultation", StartsAt: now.Add(time.Minute), EndsAt: now.Add(61 * time.Minute), Capacity: 2}
	if _, err = NewService(&fakeRepository{}, &fixedIDs{}, fixedClock{now}).CreateAppointmentSlot(context.Background(), "tenant", slot); err == nil {
		t.Fatal("reference lead time bypass")
	}
	if _, err = service.CreateAppointmentSlot(context.Background(), "tenant", slot); err != nil {
		t.Fatal("configuration ignored", err)
	}
	t.Run("legacy-constructor-rejects-custom-policy", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("legacy service allowed a different repository policy")
			}
		}()
		NewService(repo, &fixedIDs{}, fixedClock{now})
	})
	repo.hash = businesspolicy.ReferenceSHA256
	if _, err = NewServiceWithProfile(repo, &fixedIDs{}, fixedClock{now}, p); err == nil {
		t.Fatal("mismatched service/repository profile")
	}
	if _, err = NewServiceWithProfile(&fakeRepository{}, &fixedIDs{}, fixedClock{now}, p); err == nil {
		t.Fatal("repository without binding")
	}
	if _, err = NewServiceWithProfile(repo, &fixedIDs{}, fixedClock{now}, nil); err == nil {
		t.Fatal("nil profile")
	}
}
