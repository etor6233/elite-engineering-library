package main

// AUTHORED finance workspace host mount. Always registers routes; FX snapshot
// is optional and shared with selectedFXConversionModule when FX_ENABLED=true.
import (
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/royalty"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func selectedFinanceWorkspaceModule(pool *pgxpool.Pool, ids royalty.IDGenerator, lookup func(string) string) httpapi.EnterpriseModule {
	snapshot, _ := loadFXSnapshot(lookup, time.Now())
	return httpapi.FinanceWorkspaceModule{Store: postgres.NewFinanceStore(pool, ids, snapshot)}
}
