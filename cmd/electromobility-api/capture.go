package main

// AUTHORED opt-in host adapter. The profile contains deployment paths and hashes,
// never browser-supplied executable names or arguments. Root composition adds
// selectedCaptureModule to the existing EnterpriseModule list.
import (
	"context"
	cap "elite.local/enterprise/internal/capture"
	doc "elite.local/enterprise/internal/documentbridge"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func selectedCaptureModule(ctx context.Context, pool *pgxpool.Pool, getenv func(string) string) (httpapi.EnterpriseModule, error) {
	if getenv == nil {
		return nil, cap.ErrContract
	}
	switch getenv("CAPTURE_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, cap.ErrContract
	}
	if pool == nil || ctx.Err() != nil {
		return nil, cap.ErrContract
	}
	path, sha := getenv("CAPTURE_PROFILE_FILE"), getenv("CAPTURE_PROFILE_SHA256")
	info, e := os.Lstat(path)
	if e != nil || !info.Mode().IsRegular() || info.Size() > 32768 {
		return nil, cap.ErrContract
	}
	b, e := os.ReadFile(path)
	if e != nil || doc.Hash(b) != sha {
		return nil, cap.ErrContract
	}
	var c cap.Config
	if json.Unmarshal(b, &c) != nil {
		return nil, cap.ErrContract
	}
	decoder, e := cap.NewDecoder(c)
	if e != nil {
		return nil, e
	}
	return httpapi.CaptureModule{Decoder: decoder, Store: postgres.NewCaptureStore(pool)}, nil
}
