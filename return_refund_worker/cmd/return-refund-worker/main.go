package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"elite.local/return-refund-worker/internal/operationaltelemetry"
	"elite.local/return-refund-worker/internal/refundworker"
	"github.com/jackc/pgx/v5/pgxpool"
)

func claimToken() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}

func run() (result error) {
	databaseURL := os.Getenv("DATABASE_URL")
	workerID := os.Getenv("REFUND_WORKER_ID")
	if databaseURL == "" || workerID == "" {
		return errors.New("DATABASE_URL and REFUND_WORKER_ID are required")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cleanup, err := nativeStopContext(ctx, os.Getenv("ELITE_STOP_EVENT_HANDLE"))
	if err != nil {
		return err
	}
	defer func() {
		if cleanupErr := cleanup(); cleanupErr != nil && result == nil {
			result = cleanupErr
		}
	}()
	var reporter *operationaltelemetry.Reporter
	config, digest := os.Getenv("REFUND_TELEMETRY_CONFIG"), os.Getenv("REFUND_TELEMETRY_CONFIG_SHA256")
	if config != "" || digest != "" {
		reporter, err = operationaltelemetry.FromFile(config, digest)
		if err != nil {
			return err
		}
		defer reporter.Close()
	}
	observe := func(ctx context.Context, outcome string, elapsed time.Duration) error {
		if reporter == nil {
			return nil
		}
		return reporter.Observe(ctx, outcome, elapsed)
	}
	if ctx.Err() != nil {
		return nil
	}
	if err = observe(ctx, "started", 0); err != nil {
		return err
	}
	defer func() {
		outcome := "stopped"
		if result != nil {
			outcome = "error"
		}
		// A separate finite reporting budget records cooperative stop after parent
		// cancellation. It never retries an observation whose delivery is unknown.
		if !errors.Is(result, operationaltelemetry.ErrUnavailable) {
			if reportErr := observe(context.Background(), outcome, 0); reportErr != nil {
				result = reportErr
			}
		}
	}()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		return err
	}
	providers := map[string]refundworker.Provider{}
	if secret := os.Getenv("STRIPE_SECRET_KEY"); secret != "" {
		providers["stripe"], err = refundworker.NewStripeProvider(secret)
		if err != nil {
			return err
		}
	}
	if token := os.Getenv("MERCADO_PAGO_ACCESS_TOKEN"); token != "" {
		providers["mercado_pago"], err = refundworker.NewMercadoPagoProvider(token)
		if err != nil {
			return err
		}
	}
	processor, err := refundworker.NewProcessor(refundworker.NewPostgresStore(pool), providers, workerID, 2*time.Minute, 30*time.Second, 20)
	if err != nil {
		return err
	}
	if err = runObservedRefundLoop(ctx, processor.Step, claimToken, observe); err != nil {
		return err
	}
	return nativeStopFailure(ctx)
}

// The real host stops claiming work when its error report cannot be written.
// The existing durable processor remains responsible for reconciliation. This
// does not retry the uncertain report or change a refund result after commit.
func runRefundLoop(ctx context.Context, step func(context.Context, string) error, nextToken func() (string, error)) error {
	return runObservedRefundLoop(ctx, step, nextToken, nil)
}

// handled describes Step returning nil; the durable domain record, not this
// operational outcome, determines whether money moved or a refund succeeded.
func runObservedRefundLoop(ctx context.Context, step func(context.Context, string) error, nextToken func() (string, error), observe func(context.Context, string, time.Duration) error) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		if ctx.Err() != nil {
			return nil
		}
		token, tokenErr := nextToken()
		if tokenErr != nil {
			return tokenErr
		}
		began := time.Now()
		stepErr := step(ctx, token)
		if reportErr := logStepError(stepErr); reportErr != nil {
			return reportErr
		}
		if observe != nil {
			outcome := "handled"
			if errors.Is(stepErr, refundworker.ErrNoWork) {
				outcome = "idle"
			} else if stepErr != nil {
				outcome = "error"
			}
			if reportErr := observe(ctx, outcome, time.Since(began)); reportErr != nil {
				return reportErr
			}
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func main() {
	exitOnHostFailure(run())
}

func exitOnHostFailure(err error) {
	if err != nil {
		log.Fatal("return refund host failed: REFUND_HOST_FAILED")
	}
}

// Log only host-owned operational outcomes. Durable results own diagnostics.
var errRefundReport = errors.New("return refund report unavailable: REFUND_REPORT_FAILED")

type checkedRefundLogWriter struct{ writer io.Writer }

func (w checkedRefundLogWriter) Write(p []byte) (int, error) {
	if w.writer == nil {
		return 0, errRefundReport
	}
	n, err := w.writer.Write(p)
	if n < 0 || n > len(p) {
		return 0, errRefundReport
	}
	if err != nil || n != len(p) {
		return n, errRefundReport
	}
	return n, nil
}

// Preserve the host logger's configured prefix/flags and private fixed message,
// but surface delivery loss. The host loop emits serially. A blocking writer
// still requires target-owned process supervision; no I/O deadline is claimed.
// Full local Write acceptance does not prove log retention or alert delivery.
func logStepError(err error) error {
	if err != nil && !errors.Is(err, refundworker.ErrNoWork) {
		logger := log.New(checkedRefundLogWriter{log.Writer()}, log.Prefix(), log.Flags())
		if logger.Output(2, "return refund step failed: REFUND_STEP_FAILED") != nil {
			return errRefundReport
		}
	}
	return nil
}
