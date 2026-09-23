package main

import "testing"

func TestCustomerFeedbackActivation(t *testing.T) {
	for _, value := range []string{"", "0", "1", "true", "false", " 1", "2"} {
		t.Run(value, func(t *testing.T) {
			module, err := selectedCustomerFeedbackModule(nil, func(key string) string {
				if key != "CUSTOMER_SURVEYS_ENABLED" {
					t.Fatal(key)
				}
				return value
			})
			if module != nil {
				t.Fatal("nil database mounted")
			}
			if (err != nil) != (value != "" && value != "0") {
				t.Fatalf("activation %q: %v", value, err)
			}
		})
	}
}
