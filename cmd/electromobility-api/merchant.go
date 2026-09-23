package main

// AUTHORED activation binding. Account credentials are supplied only in the
// future process environment; no provider call is made during module creation.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	mb "elite.local/enterprise/internal/merchantbridge"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/base64"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"os"
)

var errMerchantConfiguration = errors.New("merchant activation invalid")

func init() { merchantModuleFactory = selectedMerchantModule }

type merchantRuntime struct {
	Schema            string `json:"schema"`
	Python            string `json:"python"`
	PythonSHA256      string `json:"python_sha256"`
	Script            string `json:"script"`
	ScriptSHA256      string `json:"script_sha256"`
	OwnerSHA256       string `json:"owner_sha256"`
	RuntimeLockSHA256 string `json:"runtime_lock_sha256"`
	Mode              string `json:"mode"`
	FixtureOrigin     string `json:"fixture_origin"`
}

func merchantConfig(path, hash string, v any) error {
	f, e := os.Open(path)
	if e != nil {
		return errMerchantConfiguration
	}
	defer f.Close()
	raw, e := io.ReadAll(io.LimitReader(f, 32769))
	if e != nil || len(raw) > 32768 || mb.Hash(raw) != hash || mb.Decode(raw, v) != nil {
		return errMerchantConfiguration
	}
	return nil
}
func selectedMerchantModule(ctx context.Context, pool *pgxpool.Pool, price *postgres.Commerce, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	if lookup == nil {
		return nil, errMerchantConfiguration
	}
	switch lookup("MERCHANT_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errMerchantConfiguration
	}
	if lookup("CATALOG_RELEASE_ENABLED") != "true" || pool == nil || price == nil {
		return nil, errMerchantConfiguration
	}
	var profile mb.Profile
	var runtime merchantRuntime
	if merchantConfig(lookup("MERCHANT_PROFILE_FILE"), lookup("MERCHANT_PROFILE_SHA256"), &profile) != nil || profile.Validate() != nil || merchantConfig(lookup("MERCHANT_RUNTIME_FILE"), lookup("MERCHANT_RUNTIME_SHA256"), &runtime) != nil || runtime.Schema != "elite-merchant-runtime/v1" {
		return nil, errMerchantConfiguration
	}
	cp, e := cr.LoadProfile(lookup("CATALOG_RELEASE_PROFILE_FILE"), lookup("CATALOG_RELEASE_PROFILE_SHA256"))
	if e != nil || cp.Origin != profile.Origin {
		return nil, errMerchantConfiguration
	}
	catalog, e := postgres.NewCatalogRelease(pool, cp, price)
	if e != nil {
		return nil, errMerchantConfiguration
	}
	key, e := base64.StdEncoding.DecodeString(lookup("MERCHANT_HMAC_KEY"))
	if e != nil || len(key) != 32 {
		return nil, errMerchantConfiguration
	}
	process := mb.Process{Python: runtime.Python, PythonSHA256: runtime.PythonSHA256, Script: runtime.Script, ScriptSHA256: runtime.ScriptSHA256, OwnerSHA256: runtime.OwnerSHA256, RuntimeLockSHA256: runtime.RuntimeLockSHA256, Mode: runtime.Mode, FixtureOrigin: runtime.FixtureOrigin, CredentialsFile: lookup("GOOGLE_APPLICATION_CREDENTIALS")}
	if process.Validate() != nil {
		return nil, errMerchantConfiguration
	}
	var ready bool
	e = pool.QueryRow(ctx, `select
 exists(select 1 from pg_trigger where tgrelid=to_regclass('catalog.merchant_effect') and tgname='merchant_effect_immutable' and tgenabled in('O','A'))
 and exists(select 1 from pg_trigger where tgrelid=to_regclass('catalog.merchant_observation') and tgname='merchant_observation_immutable' and tgenabled in('O','A'))
 and (select count(*) from pg_trigger where not tgisinternal and tgenabled in('O','A') and
 (tgrelid=to_regclass('approval.request') and tgname='bound_request_immutable' or tgrelid=to_regclass('approval.decision') and tgname='bound_decision_immutable'))=2
 and pg_get_functiondef(to_regprocedure('approval.guard_bound_request()')) like '%marketplace_mutation%'
 and pg_get_functiondef(to_regprocedure('approval.guard_bound_decision()')) like '%marketplace_mutation%'`).Scan(&ready)
	if e != nil || !ready {
		return nil, errMerchantConfiguration
	}
	service, e := postgres.NewMerchantPublication(catalog, profile, key, process)
	if e != nil {
		return nil, errMerchantConfiguration
	}
	return httpapi.MerchantModule{Service: service, TenantID: profile.TenantID, OrganizationID: profile.OrganizationID}, nil
}
