package main

// AUTHORED optional activation. No route is mounted until profile, source
// manifest and isolated calculation process pass their exact local bindings.
import (
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"path/filepath"
)

var errStoredValueConfiguration = errors.New("stored-value activation is invalid")

func init() { storedValueModuleFactory = selectedStoredValueModule }
func selectedStoredValueModule(ctx context.Context, pool *pgxpool.Pool, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	if lookup == nil {
		return nil, errStoredValueConfiguration
	}
	switch lookup("STORED_VALUE_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errStoredValueConfiguration
	}
	if pool == nil {
		return nil, errStoredValueConfiguration
	}
	path := lookup("STORED_VALUE_PROFILE_FILE")
	if !filepath.IsAbs(path) {
		return nil, errStoredValueConfiguration
	}
	f, e := os.Open(path)
	if e != nil {
		return nil, errStoredValueConfiguration
	}
	info, e := f.Stat()
	f.Close()
	if e != nil || !info.Mode().IsRegular() || info.Size() > 65536 {
		return nil, errStoredValueConfiguration
	}
	raw, e := os.ReadFile(path)
	if e != nil {
		return nil, errStoredValueConfiguration
	}
	profile, e := sv.Load(raw, lookup("STORED_VALUE_PROFILE_SHA256"))
	if e != nil {
		return nil, errStoredValueConfiguration
	}
	tenant, org := profile.Scope()
	if tenant != lookup("STORED_VALUE_TENANT_ID") || org != lookup("STORED_VALUE_ORGANIZATION_ID") {
		return nil, errStoredValueConfiguration
	}
	process := sv.Process{Python: lookup("STORED_VALUE_PYTHON"), Script: lookup("STORED_VALUE_SCRIPT"), ScriptSHA256: lookup("STORED_VALUE_SCRIPT_SHA256"), Manifest: lookup("STORED_VALUE_MANIFEST"), ManifestSHA256: lookup("STORED_VALUE_MANIFEST_SHA256")}
	if process.Preflight(ctx, profile) != nil {
		return nil, errStoredValueConfiguration
	}
	store, e := postgres.NewStoredValue(pool, profile, process)
	if e != nil {
		return nil, errStoredValueConfiguration
	}
	return httpapi.StoredValueModule{Service: store}, nil
}
