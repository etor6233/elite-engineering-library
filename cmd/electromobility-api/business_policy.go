package main

// AUTHORED configuration and owner-construction glue. No vendor attribution.
import (
	"elite.local/enterprise/internal/businesspolicy"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func selectedBusinessPolicy(lookup func(string) string) (*businesspolicy.Profile, error) {
	if lookup == nil {
		return nil, businesspolicy.ErrProfile
	}
	path, hash := lookup("BUSINESS_POLICY_PROFILE_FILE"), lookup("BUSINESS_POLICY_PROFILE_SHA256")
	if path == "" && hash == "" {
		return businesspolicy.Reference(), nil
	}
	if path == "" || hash == "" {
		return nil, businesspolicy.ErrProfile
	}
	return businesspolicy.LoadFile(path, hash)
}

type businessPolicyRuntime struct {
	profile            *businesspolicy.Profile
	commerceRepository *postgres.Commerce
	journeyRepository  *postgres.FranchiseJourney
	commerceService    *commerce.Service
	journeyService     *franchisejourney.Service
}

func prepareBusinessPolicyRuntime(pool *pgxpool.Pool, ids franchisejourney.IDGenerator, clock franchisejourney.Clock, profile *businesspolicy.Profile) (*businessPolicyRuntime, error) {
	commerceRepository, err := postgres.NewCommerceWithProfile(pool, profile)
	if err != nil {
		return nil, err
	}
	journeyRepository, err := postgres.NewFranchiseJourneyWithProfile(pool, profile)
	if err != nil {
		return nil, err
	}
	journeyService, err := franchisejourney.NewServiceWithProfile(journeyRepository, ids, clock, profile)
	if err != nil {
		return nil, err
	}
	return &businessPolicyRuntime{
		profile: profile, commerceRepository: commerceRepository, journeyRepository: journeyRepository,
		commerceService: commerce.NewService(commerceRepository, ids), journeyService: journeyService,
	}, nil
}
