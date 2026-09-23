package postgres

// AUTHORED differential fixtures: source inclusive-ending predicate versus the
// established local half-open timestamp contract. This does not run AL suites.
import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestBCPriceEligibilityInclusiveSourceHalfOpenStorage(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires isolated loopback confirmation database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	at := time.Date(2026, 9, 11, 12, 0, 0, 123456000, time.UTC)
	count := 0
	for _, status := range []string{"draft", "active", "retired"} {
		for _, fromDelta := range []time.Duration{-24 * time.Hour, -time.Microsecond, 0, time.Microsecond, 24 * time.Hour} {
			from := at.Add(fromDelta)
			for _, untilDelta := range []time.Duration{-24 * time.Hour, -time.Microsecond, 0, time.Microsecond, 24 * time.Hour, 48 * time.Hour} {
				var until *time.Time
				if untilDelta != 48*time.Hour {
					v := at.Add(untilDelta)
					until = &v
				}
				for _, currency := range []string{"ARS", "USD"} {
					// $3 matches the exact caller binding used by AddOrderLine.
					query := `select ` + bcPriceEligibilitySQL + ` and ` + bcPriceCurrencyBindingSQL + `
					from (select $1::text status,$2::text currency,$4::timestamptz valid_from,$5::timestamptz valid_until) b
					cross join (select $6::timestamptz at) price_context`
					var got bool
					if err := pool.QueryRow(ctx, query, status, currency, "ARS", from, until, at).Scan(&got); err != nil {
						t.Fatal(err)
					}
					// Independent original local contract; includes exact ending,
					// next microsecond and T176's starting-after-observation refusal.
					want := status == "active" && !from.After(at) && (until == nil || at.Before(*until)) && currency == "ARS"
					if got != want {
						t.Fatalf("eligibility drift: status=%s from=%s until=%v currency=%s got=%t want=%t", status, from, until, currency, got, want)
					}
					count++
				}
			}
		}
	}
	if count != 180 {
		t.Fatalf("fixture matrix incomplete: %d", count)
	}
}
