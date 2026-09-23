package main

// AUTHORED optional composition hook. This file belongs to the base API pack;
// it does not import WhatsApp/AI owners into a backend-only composition.
import (
	"context"
	"errors"

	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type whatsappRuntime interface {
	httpapi.EnterpriseModule
	run(context.Context) error
	close() error
}

var whatsappRuntimeFactory func(context.Context, *pgxpool.Pool, identity.Verifier, func(string) string) (whatsappRuntime, error)

func selectedWhatsAppHost(ctx context.Context, pool *pgxpool.Pool, verifier identity.Verifier, lookup func(string) string) (whatsappRuntime, error) {
	if lookup == nil {
		return nil, errors.New("WhatsApp activation configuration unavailable")
	}
	switch lookup("WHATSAPP_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
		if whatsappRuntimeFactory == nil {
			return nil, errors.New("WhatsApp host pack is not selected")
		}
		return whatsappRuntimeFactory(ctx, pool, verifier, lookup)
	default:
		return nil, errors.New("WHATSAPP_ENABLED must be true, false or absent")
	}
}
