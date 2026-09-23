package main

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"
)

func completePaymentEnvironment() map[string]string {
	return map[string]string{
		"PAYMENT_CHECKOUT_ENABLED": "true", "PAYMENT_TENANT_ID": "tenant", "PAYMENT_ORGANIZATION_ID": "store", "PAYMENT_CONNECTION_ID": "connection", "PAYMENT_REQUEST_PROVIDER": "stripe", "PAYMENT_ACCOUNT_REF": "acct_fixture", "PAYMENT_CURRENCY": "ARS", "PAYMENT_MINOR_UNIT_EXPONENT": "2", "PAYMENT_LIVE_MODE": "false", "PAYMENT_SUCCESS_URL": "https://portal.example.test/customer", "PAYMENT_CANCEL_URL": "https://portal.example.test/customer", "PAYMENT_DISPLAY_NAME": "Order", "PAYMENT_PROVIDER_SECRET": "sk_test_fixture", "PAYMENT_WEBHOOK_SECRET": "whsec_fixture_00000000000000", "PAYMENT_WORKER_ID": "fixture-worker", "PAYMENT_OUTBOUND_HMAC_KEY_BASE64": base64.StdEncoding.EncodeToString([]byte(strings.Repeat("f", 32))),
	}
}
func TestPaymentHostDisabledDoesNotRequireSecrets(t *testing.T) {
	for _, enabled := range []string{"", "false"} {
		lookup := func(key string) string {
			if key == "PAYMENT_CHECKOUT_ENABLED" {
				return enabled
			}
			return ""
		}
		value, err := preparePaymentRuntime(context.Background(), nil, lookup, nil)
		if err != nil || value != nil {
			t.Fatal(value, err)
		}
	}
}
func TestPaymentHostRejectsIncompleteConfiguration(t *testing.T) {
	base := completePaymentEnvironment()
	if _, err := loadPaymentConfiguration(func(k string) string { return base[k] }); err != nil {
		t.Fatal(err)
	}
	for key := range base {
		if key == "PAYMENT_CHECKOUT_ENABLED" {
			continue
		}
		t.Run(key, func(t *testing.T) {
			env := completePaymentEnvironment()
			delete(env, key)
			if _, err := loadPaymentConfiguration(func(k string) string { return env[k] }); err == nil {
				t.Fatal("missing input enabled checkout")
			}
		})
	}
	for _, tc := range []struct{ key, value string }{{"PAYMENT_CHECKOUT_ENABLED", "yes"}, {"PAYMENT_LIVE_MODE", ""}, {"PAYMENT_PROVIDER_SECRET", "sk_live_fixture"}, {"PAYMENT_MINOR_UNIT_EXPONENT", "-1"}, {"PAYMENT_SUCCESS_URL", "http://portal.example.test/"}, {"PAYMENT_ACCOUNT_REF", ""}, {"PAYMENT_OUTBOUND_HMAC_KEY_BASE64", "secret-not-base64"}} {
		env := completePaymentEnvironment()
		env[tc.key] = tc.value
		if _, err := loadPaymentConfiguration(func(k string) string { return env[k] }); err == nil {
			t.Fatal("invalid config accepted", tc.key)
		}
	}
}
