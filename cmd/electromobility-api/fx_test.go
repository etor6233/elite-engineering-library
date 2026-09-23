package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFXHostActivationExactSnapshot(t *testing.T) {
	p := filepath.Join("..", "..", "config", "fx", "reference-profile.json")
	r := filepath.Join("..", "..", "config", "fx", "reference-rates.json")
	raw, e := os.ReadFile(p)
	if e != nil {
		t.Fatal(e)
	}
	sum := sha256.Sum256(raw)
	config := map[string]string{"FX_ENABLED": "true", "FX_PROFILE_FILE": p, "FX_RATES_FILE": r, "FX_PROFILE_ID": "fx-reference", "FX_PROFILE_REVISION": "1", "FX_PROFILE_SHA256": hex.EncodeToString(sum[:]), "FX_TENANT_ID": "018f4d4a-7b36-7a21-8d10-2f4c54c28101", "FX_ORGANIZATION_ID": "store", "FX_LOCAL_CURRENCY": "USD"}
	now := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	lookup := func(k string) string { return config[k] }
	if s, e := loadFXSnapshot(lookup, now); e != nil || s == nil {
		t.Fatal(e)
	}
	for _, key := range []string{"FX_PROFILE_FILE", "FX_RATES_FILE", "FX_PROFILE_ID", "FX_PROFILE_REVISION", "FX_PROFILE_SHA256", "FX_TENANT_ID", "FX_ORGANIZATION_ID", "FX_LOCAL_CURRENCY"} {
		saved := config[key]
		config[key] = ""
		if _, e := loadFXSnapshot(lookup, now); e == nil {
			t.Fatal("incomplete", key)
		}
		config[key] = saved
	}
	if _, e := loadFXSnapshot(lookup, time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)); e == nil {
		t.Fatal("expired")
	}
	config["FX_ENABLED"] = "false"
	if v, e := selectedFXConversionModule(nil, lookup); e != nil || v != nil {
		t.Fatal("disabled mounted")
	}
	config["FX_ENABLED"] = "invalid"
	if _, e := loadFXSnapshot(lookup, now); e == nil {
		t.Fatal("flag")
	}
}
