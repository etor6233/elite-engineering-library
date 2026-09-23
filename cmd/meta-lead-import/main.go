package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"elite.local/enterprise/internal/leadstream"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxArtifactBytes int64 = 32 << 20

func main() {
	providerPath := flag.String("provider-response", "", "V226 provider-response.json")
	batchPath := flag.String("candidate-batch", "", "V226 lead-candidates.json")
	receiptPath := flag.String("retrieval-receipt", "", "V226 RETRIEVAL_RECEIPT.json")
	flag.Parse()
	if *providerPath == "" || *batchPath == "" || *receiptPath == "" || os.Getenv("DATABASE_URL") == "" {
		slog.Error("DATABASE_URL and all three V226 artifact paths are required")
		os.Exit(2)
	}
	providerBytes, err := readBounded(*providerPath)
	if err != nil {
		fail("provider response unavailable", err)
	}
	batchBytes, err := readBounded(*batchPath)
	if err != nil {
		fail("candidate batch unavailable", err)
	}
	receiptBytes, err := readBounded(*receiptPath)
	if err != nil {
		fail("retrieval receipt unavailable", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		fail("database configuration invalid", err)
	}
	config.MaxConns = 4
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		fail("database pool unavailable", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		fail("database unavailable", err)
	}
	result, err := leadstream.ImportMetaEvidence(ctx, postgres.NewLeadIngress(pool), providerBytes, batchBytes, receiptBytes, time.Now().UTC())
	encoded, _ := json.Marshal(result)
	if err != nil {
		fmt.Fprintln(os.Stderr, string(encoded))
		fail("Meta lead import incomplete; replay is safe", err)
	}
	fmt.Println(string(encoded))
}

func readBounded(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := io.ReadAll(io.LimitReader(file, maxArtifactBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(bytes)) > maxArtifactBytes {
		return nil, fmt.Errorf("artifact exceeds %d bytes", maxArtifactBytes)
	}
	if len(bytes) == 0 {
		return nil, fmt.Errorf("artifact is empty")
	}
	return bytes, nil
}

func fail(message string, err error) {
	slog.Error(message, "error", err)
	os.Exit(1)
}
