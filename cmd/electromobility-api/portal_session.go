package main

import (
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func selectedPortalSessionModule(pool *pgxpool.Pool, getenv func(string) string) (*httpapi.PortalSessionModule, error) {
	enabled := getenv("OIDC_PORTAL_LIFECYCLE_ENABLED")
	if enabled == "" || enabled == "false" {
		return nil, nil
	}
	bad := errors.New("portal session configuration rejected")
	if enabled != "true" || pool == nil {
		return nil, bad
	}
	profile, err := identity.LoadPortalProfileFile(getenv("OIDC_PORTAL_PROFILE_FILE"), getenv("OIDC_PORTAL_PROFILE_SHA256"))
	if err != nil {
		return nil, bad
	}
	d := profile.Document()
	if d.Issuer != getenv("OIDC_ISSUER") || d.BackendAudience != getenv("OIDC_AUDIENCE") || (d.Transport != "TLS" && getenv("OIDC_PORTAL_ALLOW_LOOPBACK_FIXTURE") != "true") {
		return nil, bad
	}
	return &httpapi.PortalSessionModule{Profile: profile, Store: &postgres.PortalSessions{Pool: pool}}, nil
}
