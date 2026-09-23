// AUTHORED operator command. Scope and retention policy belong to the project.
package main

import (
	"context"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"os"
	"strings"
	"time"
)

func run(args []string, getenv func(string) string, out, stderr io.Writer) int {
	flags := flag.NewFlagSet("customer-survey-retention", flag.ContinueOnError)
	flags.SetOutput(stderr)
	tenant := flags.String("tenant", "", "explicit tenant UUID")
	organization := flags.String("organization", "", "explicit organization ID")
	limit := flags.Int("limit", 0, "required maximum expired responses per invocation (1..1000)")
	if flags.Parse(args) != nil {
		return 2
	}
	if flags.NArg() != 0 || strings.TrimSpace(*tenant) != *tenant || *tenant == "" || strings.TrimSpace(*organization) != *organization || *organization == "" || *limit < 1 || *limit > 1000 {
		fmt.Fprintln(stderr, "tenant, organization and limit 1..1000 are required")
		return 2
	}
	database := getenv("DATABASE_URL")
	if database == "" {
		fmt.Fprintln(stderr, "database configuration is missing")
		return 2
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, database)
	if err != nil {
		fmt.Fprintln(stderr, "database configuration is invalid")
		return 2
	}
	defer pool.Close()
	count, err := postgres.NewCustomerFeedback(pool).PurgeExpired(ctx, *tenant, *organization, *limit)
	if err != nil {
		fmt.Fprintln(stderr, "retention operation failed; inspect authorized database diagnostics")
		return 1
	}
	if json.NewEncoder(out).Encode(map[string]int64{"deleted": count}) != nil {
		return 1
	}
	return 0
}
func main() { os.Exit(run(os.Args[1:], os.Getenv, os.Stdout, os.Stderr)) }
