package main

// AUTHORED optional hook owned by the base API pack. No frontend/portal import.
import (
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type portalRuntime interface {
	httpapi.EnterpriseModule
	run(context.Context)
}

var portalRuntimeFactory func(context.Context, *pgxpool.Pool, func(string) string) (portalRuntime, error)

func selectedPortalHost(ctx context.Context, pool *pgxpool.Pool, getenv func(string) string) (portalRuntime, error) {
	if getenv == nil {
		return nil, errors.New("portal activation unavailable")
	}
	switch getenv("OIDC_PORTAL_LIFECYCLE_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
		if portalRuntimeFactory == nil {
			return nil, errors.New("portal lifecycle pack is not selected")
		}
		return portalRuntimeFactory(ctx, pool, getenv)
	default:
		return nil, errors.New("portal lifecycle enabled must be true or false")
	}
}
