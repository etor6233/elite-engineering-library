package providerintegration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
)

type StaticRegistry struct{ connections map[string]Connection }

func (r *StaticRegistry) Lookup(providerCode, connectionID string) (Connection, bool) {
	if r == nil {
		return Connection{}, false
	}
	value, ok := r.connections[providerCode+"\x00"+connectionID]
	return value, ok
}

var (
	codePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,127}$`)
	envPattern  = regexp.MustCompile(`^[A-Z][A-Z0-9_]{2,127}$`)
)

func LoadHMACConnections(configuration []byte, lookupEnv func(string) (string, bool)) (*StaticRegistry, error) {
	registry := &StaticRegistry{connections: map[string]Connection{}}
	if len(bytes.TrimSpace(configuration)) == 0 {
		return registry, nil
	}
	if len(configuration) > 64<<10 || lookupEnv == nil {
		return nil, fmt.Errorf("invalid provider connection configuration")
	}
	var entries []struct {
		TenantID     string `json:"tenant_id"`
		ConnectionID string `json:"connection_id"`
		ProviderCode string `json:"provider_code"`
		SecretEnv    string `json:"secret_env"`
	}
	decoder := json.NewDecoder(bytes.NewReader(configuration))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&entries); err != nil || len(entries) > 100 {
		return nil, fmt.Errorf("invalid provider connection configuration")
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return nil, fmt.Errorf("invalid provider connection configuration")
	}
	for _, entry := range entries {
		if entry.TenantID == "" || !codePattern.MatchString(entry.ConnectionID) || len(entry.ProviderCode) > 64 || !codePattern.MatchString(entry.ProviderCode) || !envPattern.MatchString(entry.SecretEnv) {
			return nil, fmt.Errorf("invalid provider connection metadata")
		}
		secret, ok := lookupEnv(entry.SecretEnv)
		if !ok || len(secret) < 32 {
			return nil, fmt.Errorf("provider secret is missing or too short")
		}
		key := entry.ProviderCode + "\x00" + entry.ConnectionID
		if _, duplicate := registry.connections[key]; duplicate {
			return nil, fmt.Errorf("duplicate provider connection")
		}
		registry.connections[key] = Connection{TenantID: entry.TenantID, ConnectionID: entry.ConnectionID, ProviderCode: entry.ProviderCode, SecretRef: entry.SecretEnv, Secret: []byte(secret)}
	}
	return registry, nil
}
