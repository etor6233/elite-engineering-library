package franchisejourney

import (
	"context"
	"regexp"
	"time"
)

type AvailabilityEntry struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	ResourceID     string    `json:"resource_id,omitempty"`
	EntryType      string    `json:"entry_type"`
	ReasonCode     string    `json:"reason_code,omitempty"`
	StartsAt       time.Time `json:"starts_at"`
	EndsAt         time.Time `json:"ends_at"`
	State          string    `json:"state"`
	Version        int64     `json:"version"`
}

func (s *Service) CreateAvailability(ctx context.Context, tenant, subject string, value AvailabilityEntry) (AvailabilityEntry, error) {
	value, err := s.prepareAvailability(tenant, subject, value)
	if err != nil {
		return AvailabilityEntry{}, err
	}
	return s.repository.CreateAvailability(ctx, tenant, subject, value, s.ids.New())
}
func (s *Service) prepareAvailability(tenant, subject string, value AvailabilityEntry) (AvailabilityEntry, error) {
	value.ID, value.State, value.Version = s.ids.New(), "active", 1
	validType := value.EntryType == "working" || value.EntryType == "unavailable"
	validReason := value.EntryType == "working" && value.ReasonCode == "" || value.EntryType == "unavailable" && code(value.ReasonCode)
	if tenant == "" || subject == "" || value.OrganizationID == "" || !validType || !validReason || value.StartsAt.IsZero() || !value.EndsAt.After(value.StartsAt) || value.EndsAt.Sub(value.StartsAt) > 366*24*time.Hour {
		return AvailabilityEntry{}, ErrInvalid
	}
	return value, nil
}

func (s *Service) CancelAvailability(ctx context.Context, tenant, organization, entry string, version int64, subject, reason string) (AvailabilityEntry, error) {
	if tenant == "" || organization == "" || entry == "" || version < 1 || subject == "" || !code(reason) {
		return AvailabilityEntry{}, ErrInvalid
	}
	return s.repository.CancelAvailability(ctx, tenant, organization, entry, version, subject, reason, s.ids.New())
}

func (s *Service) Availability(ctx context.Context, tenant, organization, resource string, from, to time.Time) ([]AvailabilityEntry, error) {
	if tenant == "" || organization == "" || from.IsZero() || !to.After(from) || to.Sub(from) > 366*24*time.Hour {
		return nil, ErrInvalid
	}
	return s.repository.Availability(ctx, tenant, organization, resource, from, to)
}

func (s *Service) CancelCustomerAppointment(ctx context.Context, tenant, organization, customer, appointment string, version int64, reason string) (Appointment, error) {
	if tenant == "" || organization == "" || customer == "" || appointment == "" || version < 1 || !code(reason) {
		return Appointment{}, ErrInvalid
	}
	return s.repository.CancelCustomerAppointment(ctx, tenant, organization, customer, appointment, version, reason, s.ids.New(), s.ids.New())
}

var availabilityKeyPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{16,128}$`)

func (s *Service) CreateAvailabilityOnce(ctx context.Context, tenant, subject, key, hash string, value AvailabilityEntry) (AvailabilityEntry, bool, error) {
	if !availabilityKeyPattern.MatchString(key) || !sha256Hex(hash) {
		return AvailabilityEntry{}, false, ErrInvalid
	}
	value, err := s.prepareAvailability(tenant, subject, value)
	if err != nil {
		return AvailabilityEntry{}, false, err
	}
	r, ok := s.repository.(interface {
		CreateAvailabilityOnce(context.Context, string, string, string, string, AvailabilityEntry, string) (AvailabilityEntry, bool, error)
	})
	if !ok {
		return AvailabilityEntry{}, false, ErrConflict
	}
	return r.CreateAvailabilityOnce(ctx, tenant, subject, key, hash, value, s.ids.New())
}
func (s *Service) AvailabilityCreationResult(ctx context.Context, tenant, organization, key string) (AvailabilityEntry, error) {
	if tenant == "" || organization == "" || !availabilityKeyPattern.MatchString(key) {
		return AvailabilityEntry{}, ErrInvalid
	}
	r, ok := s.repository.(interface {
		AvailabilityCreationResult(context.Context, string, string, string) (AvailabilityEntry, error)
	})
	if !ok {
		return AvailabilityEntry{}, ErrConflict
	}
	return r.AvailabilityCreationResult(ctx, tenant, organization, key)
}
