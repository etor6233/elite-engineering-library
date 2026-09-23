package main

import (
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"net/url"
	"time"
)

type portalSessionRuntime struct {
	*httpapi.PortalSessionModule
	broker interface {
		AccessToken(context.Context) (string, error)
	}
	target string
	client *http.Client
}

func init() { portalRuntimeFactory = preparePortalSessionRuntime }
func preparePortalSessionRuntime(ctx context.Context, pool *pgxpool.Pool, getenv func(string) string) (portalRuntime, error) {
	module, err := selectedPortalSessionModule(pool, getenv)
	if err != nil || module == nil {
		return nil, errors.New("portal session runtime unavailable")
	}
	secret := getenv("OIDC_PORTAL_SERVICE_CLIENT_SECRET_FILE")
	if secret == "" {
		return nil, errors.New("portal service credential reference required")
	}
	broker, err := identity.NewPortalServiceTokenBroker(ctx, module.Profile, identity.ServiceSecretFile(secret))
	if err != nil {
		return nil, err
	}
	base, err := url.Parse(module.Profile.Document().PostLogoutURL)
	if err != nil {
		return nil, err
	}
	target := base.ResolveReference(&url.URL{Path: "/api/internal/oidc/session-maintenance"}).String()
	return &portalSessionRuntime{PortalSessionModule: module, broker: broker, target: target, client: &http.Client{Timeout: 45 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("portal maintenance redirect rejected") }}}, nil
}
func (p *portalSessionRuntime) poll(ctx context.Context) error {
	token, err := p.broker.AccessToken(ctx)
	if err != nil {
		return errors.New("portal service grant unavailable")
	}
	request, err := http.NewRequestWithContext(ctx, "POST", p.target, nil)
	if err != nil {
		return errors.New("portal maintenance request unavailable")
	}
	request.Header.Set("Authorization", "Bearer "+token)
	response, err := p.client.Do(request)
	if err != nil {
		return errors.New("portal maintenance transport unavailable")
	}
	defer response.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(response.Body, 4097))
	if err != nil || len(raw) > 4096 {
		return errors.New("portal maintenance response unavailable")
	}
	var result struct {
		Claimed           int `json:"claimed"`
		Confirmed         int `json:"confirmed"`
		Pending           int `json:"pending"`
		Purged            int `json:"purged"`
		UnconfirmedPurged int `json:"unconfirmed_purged"`
	}
	var fields map[string]json.RawMessage
	media, _, mediaErr := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if mediaErr != nil || media != "application/json" || json.Unmarshal(raw, &fields) != nil || len(fields) != 5 || json.Unmarshal(raw, &result) != nil || response.StatusCode != 200 && response.StatusCode != 503 {
		return errors.New("portal maintenance result unavailable")
	}
	for _, name := range []string{"claimed", "confirmed", "pending", "purged", "unconfirmed_purged"} {
		if _, ok := fields[name]; !ok {
			return errors.New("portal maintenance result incomplete")
		}
	}
	if result.Claimed < 0 || result.Claimed > 2 || result.Confirmed < 0 || result.Pending < 0 || result.Confirmed+result.Pending != result.Claimed || result.Purged < 0 || result.Purged > 100 || result.UnconfirmedPurged < 0 || result.UnconfirmedPurged > result.Purged {
		return errors.New("portal maintenance counters invalid")
	}
	if result.UnconfirmedPurged > 0 {
		slog.Error("portal credentials removed at retention deadline without confirmed provider revocation", "count", result.UnconfirmedPurged)
	}
	if response.StatusCode != 200 || result.Pending > 0 {
		return errors.New("portal provider revocation pending")
	}
	return nil
}
func (p *portalSessionRuntime) run(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		if ctx.Err() != nil {
			return
		}
		if err := p.poll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("portal session maintenance deferred")
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
