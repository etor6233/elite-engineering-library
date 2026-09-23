package main

import (
	"elite.local/enterprise/internal/customerfeedback"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func selectedCustomerFeedbackModule(pool *pgxpool.Pool, getenv func(string) string) (httpapi.EnterpriseModule, error) {
	switch getenv("CUSTOMER_SURVEYS_ENABLED") {
	case "", "0":
		return nil, nil
	case "1":
		if pool == nil {
			return nil, errors.New("customer survey database is missing")
		}
		return httpapi.CustomerFeedbackModule{Service: customerfeedback.NewService(postgres.NewCustomerFeedback(pool))}, nil
	default:
		return nil, errors.New("customer survey activation must be 0 or 1")
	}
}
