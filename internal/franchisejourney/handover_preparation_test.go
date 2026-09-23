package franchisejourney

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
	"time"
)

type preparationRepoProbe struct {
	calls       int
	actor, hash string
	command     PrepareHandoverCommand
}

func (p *preparationRepoProbe) PrepareInitialHandover(_ context.Context, _, actor, _, _ string, command PrepareHandoverCommand, _ HandoverReleaseContract, hash string) (HandoverPreparation, bool, error) {
	p.calls++
	p.actor = actor
	p.hash = hash
	p.command = command
	return HandoverPreparation{}, false, nil
}
func (p *preparationRepoProbe) InitialHandoverResult(context.Context, string, string, string, string) (HandoverPreparation, error) {
	p.calls++
	return HandoverPreparation{}, nil
}

type preparationIDs struct{}

func (preparationIDs) New() string { return "generated-id" }
func referencePreparationContract() HandoverReleaseContract {
	hash := sha256.Sum256([]byte(ReferenceHandoverContractDocument))
	return HandoverReleaseContract{ID: "reference-single-unit-observed-payment-v1", DocumentSHA256: hex.EncodeToString(hash[:]), Scope: "LOCAL_FIXTURES", MaximumObservationAge: 5 * time.Minute}
}
func TestHandoverPreparationRequiresExplicitExactReferenceContract(t *testing.T) {
	cases := []HandoverReleaseContract{{}, referencePreparationContract(), referencePreparationContract(), referencePreparationContract(), referencePreparationContract()}
	cases[1].Scope = "PRODUCTION"
	cases[2].DocumentSHA256 = strings.Repeat("a", 64)
	cases[3].MaximumObservationAge = 0
	cases[4].MaximumObservationAge = time.Hour
	for _, contract := range cases {
		if _, err := NewHandoverPreparationService(&preparationRepoProbe{}, preparationIDs{}, contract); !errors.Is(err, ErrReleaseConditioned) {
			t.Fatalf("unsafe contract accepted: %v", err)
		}
	}
}
func TestHandoverPreparationRequestBindsActorAndObservation(t *testing.T) {
	repo := &preparationRepoProbe{}
	service, err := NewHandoverPreparationService(repo, preparationIDs{}, referencePreparationContract())
	if err != nil {
		t.Fatal(err)
	}
	command := PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: "payment", ObservationSHA256: strings.Repeat("a", 64), IdempotencyKey: "prepare-key"}
	if _, _, err = service.Prepare(context.Background(), "tenant", "operator", command); err != nil {
		t.Fatal(err)
	}
	original := repo.hash
	if repo.actor != "operator" || repo.command.OrderID != "order" {
		t.Fatal("principal or scope lost")
	}
	if _, _, err = service.Prepare(context.Background(), "tenant", "other", command); err != nil || repo.hash == original {
		t.Fatal("actor not bound")
	}
	command.ObservationSHA256 = strings.Repeat("b", 64)
	if _, _, err = service.Prepare(context.Background(), "tenant", "operator", command); err != nil || repo.hash == original {
		t.Fatal("observation not bound")
	}
}
func TestHandoverPreparationRejectsInvalidInputsBeforeRepository(t *testing.T) {
	repo := &preparationRepoProbe{}
	service, _ := NewHandoverPreparationService(repo, preparationIDs{}, referencePreparationContract())
	command := PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: "payment", ObservationSHA256: strings.Repeat("a", 64), IdempotencyKey: "prepare-key"}
	cases := []PrepareHandoverCommand{command, command, command, command, command}
	cases[0].OrganizationID = ""
	cases[1].OrderLineID = " "
	cases[2].PaymentAttemptID = ""
	cases[3].ObservationSHA256 = "captured"
	cases[4].IdempotencyKey = "short"
	for _, bad := range cases {
		if _, _, err := service.Prepare(context.Background(), "tenant", "operator", bad); !errors.Is(err, ErrInvalid) {
			t.Fatal(err)
		}
	}
	if repo.calls != 0 {
		t.Fatal("invalid request reached repository")
	}
}
