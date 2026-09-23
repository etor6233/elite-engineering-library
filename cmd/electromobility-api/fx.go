package main

// AUTHORED explicit FX host activation. Disabled means no FX route is mounted.
import (
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcfx"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"time"
)

var errFXConfiguration = errors.New("FX snapshot activation is invalid")

func loadFXSnapshot(lookup func(string) string, now time.Time) (*bcfx.Snapshot, error) {
	if lookup == nil {
		return nil, errFXConfiguration
	}
	switch lookup("FX_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errFXConfiguration
	}
	revision, e := strconv.Atoi(lookup("FX_PROFILE_REVISION"))
	if e != nil || revision < 1 {
		return nil, errFXConfiguration
	}
	s, e := bcfx.LoadSnapshotFiles(lookup("FX_PROFILE_FILE"), lookup("FX_RATES_FILE"), bcfx.Activation{Enabled: true, ProfileID: lookup("FX_PROFILE_ID"), Revision: revision, ProfileSHA256: lookup("FX_PROFILE_SHA256"), TenantID: lookup("FX_TENANT_ID"), OrganizationID: lookup("FX_ORGANIZATION_ID"), LocalCurrency: lookup("FX_LOCAL_CURRENCY")})
	if e != nil || !s.Current(now) {
		return nil, errFXConfiguration
	}
	return s, nil
}
func selectedFXConversionModule(pool *pgxpool.Pool, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	snapshot, e := loadFXSnapshot(lookup, time.Now())
	if e != nil || snapshot == nil {
		return nil, e
	}
	if pool == nil {
		return nil, errFXConfiguration
	}
	service, e := accounting.NewFXService(postgres.NewAccounting(pool), randomid.Generator{}, snapshot)
	if e != nil {
		return nil, errFXConfiguration
	}
	return httpapi.FXConversionModule{Service: service}, nil
}
