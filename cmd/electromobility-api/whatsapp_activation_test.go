package main

import (
	"context"
	"testing"
)

func TestWhatsAppActivationWithoutOptionalPack(t *testing.T) {
	saved := whatsappRuntimeFactory
	whatsappRuntimeFactory = nil
	defer func() { whatsappRuntimeFactory = saved }()
	if h, e := selectedWhatsAppHost(context.Background(), nil, nil, func(string) string { return "true" }); e == nil || h != nil {
		t.Fatal("activation without selected implementation")
	}
	if h, e := selectedWhatsAppHost(context.Background(), nil, nil, func(k string) string {
		if k != "WHATSAPP_ENABLED" {
			t.Fatal("read inactive secret/config")
		}
		return "false"
	}); e != nil || h != nil {
		t.Fatal("disabled optional owner")
	}
}
