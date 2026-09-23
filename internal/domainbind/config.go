// Package domainbind binds the conversational agent's business intents to the
// franchise backend HTTP endpoints. Business-specific mappings (service→kind,
// product→variant, price book) are injected as configuration, never invented;
// unknown values fail closed.
package domainbind

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Config carries the backend endpoint and the business mappings.
type Config struct {
	BaseURL         string            // backend root, e.g. http://localhost:8080
	TenantID        string            // exact tenant carried by the authenticated session
	TenantCode      string            // public tenant code
	OrganizationID  string            // optional legacy single-contact default
	LeadID          string            // optional legacy single-contact default
	ServiceKinds    map[string]string // "corte" -> "service"
	ProductVariants map[string]string // "scooter" -> "variant-1"
	PriceBookID     string
	ValidMinutes    int
	Timeout         time.Duration
	TokenProvider   AccessTokenProvider // current short-lived token source for protected routes
}

// AccessTokenProvider returns the current token at request time. Its
// implementation owns acquisition, refresh and secret storage.
type AccessTokenProvider interface {
	AccessToken(context.Context) (string, error)
}

// Validate enforces a safe, exact configuration.
func (c Config) Validate() error {
	u, err := url.Parse(c.BaseURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return errors.New("domainbind: base url must be absolute http(s)")
	}
	if strings.TrimSpace(c.TenantID) == "" || strings.TrimSpace(c.TenantCode) == "" {
		return errors.New("domainbind: tenant and tenant code required")
	}
	if (strings.TrimSpace(c.OrganizationID) == "") != (strings.TrimSpace(c.LeadID) == "") {
		return errors.New("domainbind: legacy organization and lead defaults must be both present or both absent")
	}
	if len(c.ServiceKinds) == 0 || len(c.ProductVariants) == 0 || strings.TrimSpace(c.PriceBookID) == "" {
		return errors.New("domainbind: service kinds, product variants and price book required")
	}
	if c.ValidMinutes <= 0 || c.Timeout <= 0 {
		return errors.New("domainbind: valid minutes and timeout required")
	}
	return nil
}

// ResolveService returns the appointment kind for a service name, fail-closed.
func (c Config) ResolveService(service string) (string, error) {
	kind, ok := c.ServiceKinds[strings.ToLower(strings.TrimSpace(service))]
	if !ok || kind == "" {
		return "", fmt.Errorf("domainbind: unknown service %q", service)
	}
	return kind, nil
}

// ResolveProduct returns the variant id for a product name, fail-closed.
func (c Config) ResolveProduct(product string) (string, error) {
	variant, ok := c.ProductVariants[strings.ToLower(strings.TrimSpace(product))]
	if !ok || variant == "" {
		return "", fmt.Errorf("domainbind: unknown product %q", product)
	}
	return variant, nil
}
