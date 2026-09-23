package franchisejourney

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestAvailabilityRequiresExplicitIntervalsAndReasons(t *testing.T) {
	service, repository, _ := newService()
	from := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	value, err := service.CreateAvailability(context.Background(), "tenant", "scheduler", AvailabilityEntry{OrganizationID: "store", ResourceID: "employee", EntryType: "unavailable", ReasonCode: "annual-leave", StartsAt: from, EndsAt: from.Add(8 * time.Hour)})
	if err != nil || value.State != "active" || repository.availability.EntryType != "unavailable" {
		t.Fatalf("value=%+v err=%v", value, err)
	}
	if _, err = service.CreateAvailability(context.Background(), "tenant", "scheduler", AvailabilityEntry{OrganizationID: "store", EntryType: "unavailable", StartsAt: from, EndsAt: from.Add(time.Hour)}); !errors.Is(err, ErrInvalid) {
		t.Fatalf("reasonless absence accepted: %v", err)
	}
	if _, err = service.TransitionAppointment(context.Background(), "tenant", "store", "appointment", "confirmed", "no-show", 2, "operator", ""); !errors.Is(err, ErrInvalid) {
		t.Fatalf("reasonless no-show accepted: %v", err)
	}
	if _, err = service.TransitionAppointment(context.Background(), "tenant", "store", "appointment", "confirmed", "completed", 2, "operator", "invented-reason"); !errors.Is(err, ErrInvalid) {
		t.Fatalf("irrelevant reason accepted: %v", err)
	}
}
