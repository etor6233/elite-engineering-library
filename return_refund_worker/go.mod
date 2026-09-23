module elite.local/return-refund-worker

go 1.26.0

require (
	github.com/jackc/pgx/v5 v5.10.0
	github.com/mercadopago/sdk-go v1.14.0
	github.com/stripe/stripe-go/v86 v86.3.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	// Security floor for selected upstream tools: GO-2026-6179/6180.
	// Preserve this constraint when tidying; validate the complete module graph.
	golang.org/x/mod v0.40.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/text v0.39.0 // indirect
)
